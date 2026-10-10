# Verifikation: slice-v1-abschluss-einspielen-anmeldung-doku — 2026-10-10

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-10-review-slice-v1-abschluss-einspielen-anmeldung-doku.md`, F-597 bis F-600). Ich prüfe die Übergaben nach (Abschnitt 2) und repariere nichts.

**Gegenstand:** Slice-Plan `slice-v1-abschluss-einspielen-anmeldung-doku` (Kopf, §1 bis §8) gegen den Gesamt-Diff `e8fd0b2..465f09b`. Architect: `36480b8`. Implementer: `e70ebac`, `8cc0a48`, `1898740`. Review: `05fafac`. Nacharbeit: `f060f1a`, `24d8f85`, `465f09b`. Rahmen: [MR-000](../../harness/conventions.md#mr-000--baseline-aussage); Entscheidung: [ADR-0016](../plan/adr/0016-einspielen-anmeldung-und-tls.md).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingang:**

- der Plan ganz am Stand `465f09b`, besonders §1 (*Ausdrücklich NICHT*), §2 (DoD), §3 (Ansatz, benannte Ausnahme der Befehlsfolge), §6 (Randformen-Tabelle, Risiken) und §7 (Belege des Implementers)
- `docs/user/benutzerhandbuch.md` und `README.md` ganz an den geänderten Stellen und die Abschnitte §1, §4, §5, §6, §7 vollständig; der Diff beider Dateien gegen `e8fd0b2`
- `AGENTS.md` §3.9, §3.11 (auch Handbuch-/README-Regel), §3.13; [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung) (.a Passwort, Verfahren, SCRAM), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-RB-01`](../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten) (beschriebenes Verhalten, nicht geändert)
- `slice-doku-ist-stand` §6 und §7 (Randformen und Befehlsfolge für Doku-Slices); der Katalog `internal/hexagon/model/fehler.go`; `--help` des gebauten Binaries
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen, nur die Belege in Plan §7

Beim Start war der Arbeitsbaum sauber, HEAD `465f09b`. Der Commit dieses Berichts enthält nur diese Datei.

**So wurden die Proben gebaut.**

- **Binary:** frische Kopie per `git archive HEAD | tar -x` in das Scratchpad, daraus die Stufe `runtime` des `Dockerfile` mit eigenem Tag (`ver-anmdoku-bin:dev`). Das Image ist `sha256:c2871ccf3ceb…`, dieselbe Kennung wie im Beleg des Implementers (der Produktcode ist seit `36480b8` unverändert).
- **Datenbank:** das gepinnte Image aus `harness/mk/integration.mk` (`postgres:17-alpine@sha256:b0f9560a…`) in einem eigenen internen Netz `ver-anmdoku-net`, Alias `postgres`. `pg_hba.conf` wie im Runner: `play_scram` (`scram-sha-256`, Passwort `GEHEIM scram $%41`), `play_md5` (`md5`, `GEHEIMmd5`), `play_pw` (`password`, `GEHEIMklartext`), alle übrigen `trust`; dazu von mir `play_trust` (`trust`), `play_nbsp` (SCRAM, Passwort `GEHEIM` U+00A0 `x`) und `play_sonder` (SCRAM, `a@b:c/d?e#f`). Die Aufzeichnung (`SELECT 1`) habe ich selbst mit `record` und `psql` als `postgres` angelegt; `play` und `record` liefen als Container mit `timeout 60`.
- **Danach:** Container, Netz, Image und die Scratch-Pfade `ver-anmdoku-src` und `ver-anmdoku-w` (nur diese benannten) entfernt; `docker ps -a`, `docker network ls`, `docker images` und `ls` des Scratchpads mit `ver-anmdoku` danach je 0 Treffer. Die Log-Datei `ver-anmdoku-gates.log` des Gate-Laufs bleibt bis zum Ende der Sitzung im Scratchpad. Nur Linux im Container (`BEO-REPO/verhalten-nur-unter-linux-geprueft`).

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-einspielen-anmeldung` liefert (§1, §4, §5, §7). **Bestätigt.**

Jede geänderte Aussage zur Anmeldung, einzeln gegen Binary und PostgreSQL 17:

| Aussage (Handbuch, README) | Probe, selbst beobachtet | Urteil |
|---|---|---|
| §4 *Voraussetzung*, §1, README: Klartext-Passwort, MD5 und SCRAM-SHA-256 melden an | `--user play_pw`, `play_md5`, `play_scram`, `host:port`, Passwort aus der Variable; je Exit 0 (SCRAM mit Leerzeichen, `$`, `%` im Wert) | bestätigt |
| §5: Variable gilt bei `host:port` und bei einer Verbindung ohne Passwortteil | `--upstream postgres:5432` und Verbindung `ohne` (`play_scram@…`) mit richtiger Variable: Exit 0 | bestätigt |
| §5: Schreibt die Verbindung einen Passwortteil, gilt nur dessen Wert, die Variable bleibt unbeachtet | Platzhalter richtig, Variable falsch: Exit 0; Platzhalter falsch, Variable richtig: `PGR-E4005`, `28P01`, Exit 4; Platzhalter richtig, Variable ungesetzt: Exit 0 | bestätigt |
| §5: eine leere Variable gilt als nicht gesetzt | `PGWIRE_RECORDER_PASSWORD=` gegen `play_scram` (host:port) und gegen Verbindung `ohne`: je `PGR-E4005` „der Server verlangt ein Passwort, play hat keines“ | bestätigt |
| §5: bei `trust` bleibt eine gesetzte Variable ohne Wirkung, die Anmeldung gelingt | `play_trust` mit und ohne Variable: je Exit 0 | bestätigt |
| §5, §4, §7: fehlendes Passwort endet mit `PGR-E4005`; die Meldung sagt, dass `play` keines hat | `play_pw`, `play_md5`, `play_scram` ohne Variable: je `PGR-E4005`, Exit 4, Text „… play hat keines“ | bestätigt |
| §4, §7: falsches Passwort oder fehlender Benutzer: `PGR-E4005`, Meldung nennt SQLSTATE und Meldung der Datenbank, Exit 4 | falsches Passwort je Verfahren: `28P01`; `--user niemand`: `28000` „role … does not exist“; alle Exit 4 | bestätigt |
| §5: keine Meldung, keine Log-Zeile von `play`, keine Ausgabe von `config show` nennt das Passwort; `config show` nennt den Namen der Variable, nie den Wert | Marker `GEHEIM`: `--log-level debug` und `info`, je Verfahren richtiges und falsches Passwort, Platzhalter richtig und falsch, Verbindung ohne Passwortteil, `PGR-E2006` mit Klartext-Passwort `GEHEIMklar` und Parameter `password=GEHEIMparam`, `config show` mit gesetzter Variable `GEHEIMenv`: 65 Ausgabezeilen, `grep -c GEHEIM` = 0; Gegenkontrolle der Aufnahme: 8 Zeilen mit `play_`, `level=DEBUG` kommt vor (die Aufnahme enthielt also Debug-Zeilen); `config show` endet mit der Zeile `PGWIRE_RECORDER_PASSWORD` | bestätigt |
| §5: `play` sendet das Passwort unverändert; ein Passwort mit Zeichen, die SASLprep ändert, kann bei SCRAM scheitern | `play_nbsp` mit `GEHEIM` U+00A0 `x`: `PGR-E4005`, `28P01`, Exit 4; dasselbe Passwort mit `psql`: Anmeldung gelingt; mit ASCII-Leerzeichen statt U+00A0: Exit 0 (die Datenbank speichert das SASLprep-Ergebnis) | bestätigt (Beleg nur Probe, siehe Abschnitt 2, F-600) |
| §5 Beispielverbindung `test` (`${DB_PASSWORD}`) und Satz „Mit `--upstream test` meldet sich `play` als `tester` mit dem Wert von `DB_PASSWORD` an“ | Beispieldatei aus dem Handbuch extrahiert: `config show --config` mit `CI_DB_HOST=x` gibt die Datei zeilengleich aus, Exit 0. Mit `sed` nur Benutzer, Host und Datenbank der Zeile `test` auf `play_scram`, `postgres:5432`, `postgres` gesetzt: `--upstream test` ohne `--user`, richtiges `DB_PASSWORD`: Exit 0; falsches `DB_PASSWORD`: `PGR-E4005` `28P01` (der Benutzer der URL wird also benutzt; der Aufzeichnungs-Benutzer `postgres` wäre mit `trust` durchgekommen); ohne `DB_PASSWORD`: `PGR-E2005`, Exit 2 | bestätigt |
| §5: „Bei `play` ist der eingesetzte Wert eines Platzhalters im Passwort das Passwort der Anmeldung“; Wert unverändert, auch mit Sonderzeichen | `play_sonder` mit `PW=a@b:c/d?e#f`: Exit 0; `PW=GEHEIM scram $%41`: Exit 0; `PW=` oder ungesetzt: `PGR-E2005`, Exit 2 | bestätigt |
| §4 *Hinweise*: „Verlangt die Datenbank ein Klartext-Passwort, sendet `play` es über die Verbindung, wie sie ist, also ohne TLS unverschlüsselt.“ | `play_pw` meldet sich über die unverschlüsselte Verbindung an (Exit 0); `sslmode=require` bei `play`: `PGR-E2004` „play verbindet ohne TLS zum Upstream“ | bestätigt, Reichweite siehe V-166 |
| §4 *Vorgehen*, Beispiel `play --upstream postgres:5432 --input …` | gegen den `trust`-Benutzer `postgres`: Exit 0 | bestätigt |
| §7 Zeile `PGR-E4005`: nur `SCRAM-SHA-256-PLUS`, Kerberos, GSSAPI, SSPI enden mit `PGR-E4005` | kein Server löst es aus; Test `TestAnmeldungNichtUnterstuetzt` (Fälle `nur PLUS`, `Kerberos`, `GSS`, `SSPI`) existiert in `internal/adapters/driven/postgres/einspielen_anmeldung_test.go` und lief im `make gates` grün; die Mutanten des Implementers habe ich nicht wiederholt | bestätigt |
| §1, §7 `PGR-E6001`: nur `record`, Datenbank verlangt Passwort | `record` gegen SCRAM, MD5 und Klartext: je `PGR-E6001` „der Upstream verlangt ein Anmeldeverfahren …“, Exit 6 | bestätigt |
| §7 `PGR-E2005`, `PGR-E2006` | siehe Zeilen oben (E2005: Platzhalter leer oder ungesetzt; E2006: Klartext im Passwort, leerer Passwortteil `user:@host` und Parameter `password`, auch bei `config show`) | bestätigt |
| §6 (Nachbarsatz): `record` und `replay` wie vorher, `play` meldet sich mit Benutzer, Datenbank, Passwort an | `replay` nimmt `psql` mit Benutzer `jemand` und Passwort `irgendwas` an (Antwort `1`); `record` ohne Passwort-Weitergabe siehe oben; `play` siehe oben | bestätigt |

**Option und Katalog.** `play --help` nennt `--upstream` mit der Passwortquelle (Platzhalter, sonst `PGWIRE_RECORDER_PASSWORD`), die Verfahren und `PGR-E4005`; das Handbuch widerspricht dem nicht.

**Handbuch ohne Chronik, Zielstand, Verweise.** `grep -n -i` über Handbuch und README auf `bisher|künftig|noch nicht|geplant|später|zunächst|derzeit|vorerst|dieser Stand|inzwischen|seit |früher|kommt|folgt|spezifikation|lastenheft|ADR|slice|welle|review|SPEC-|LH-|ARC-|spec/|docs/plan|roadmap|MVP`: im Handbuch nur Alltagsworte (`später ohne die Datenbank` §1 Zweck, `noch nicht` §3 Zieldatei, `früheren Laufs`, `vorher`, `kommt` in „Ob sie kommt“ und „kommt nur zum Einsatz“, `aktuellen Verzeichnis`), kein Verweis auf Spezifikation, Lastenheft, ADR, Slice, Welle oder Review und keine Aussage über einen künftigen Stand; die README nennt `spec/` und `docs/plan/` nur als Verweise in §Auditierbarkeit und §Verträge ohne Aussage über einen Stand. Alle Treffer von `ohne Passwort` (README Zeile 29, Handbuch §6, §7 `PGR-E6001`, §7 *Die Anwendung kann sich nicht verbinden*) betreffen `record`.

**Nachbarsätze.** §5 *Konfigurationsdatei*, Satz „Bei `record` zählen nur Host und Port der URL; Benutzer, Passwort und Datenbank vermittelt die Anwendung selbst“, §5 „Die Meldung nennt die Stelle (Schlüssel oder Verbindung), nie den Wert“, §4 *config show* (Hinweis „Die Ausgabe enthält nie einen aufgelösten Wert“, `PGR-E2006` bei `config show`), §4 *Aufzeichnen* („Verlangt die Datenbank ein Passwort …“, Exit 6) und §10 („Aufruf ohne Zugangsdaten“) reichen nicht weiter als die Proben; der Satz „Die Log-Zeile beim Start nennt auch hier nur `host:port`“ gilt auch für die Verbindung `mit` (Startzeile `upstream=postgres:5432`).

### Punkt 2: §5 beschreibt das Passwort bei `play` wie geliefert (Platzhalter, sonst `PGWIRE_RECORDER_PASSWORD`; Befund F-533) und nennt die Grenzen *SASLprep* und *Klartext ohne TLS*. **Bestätigt, mit Hinweis auf den Ort (V-165).**

Quelle, Vorrang, leere Variable, SASLprep: siehe Punkt 1. Die Grenze *Klartext ohne TLS* steht in §4 *Hinweise* statt in §5; das entspricht der Entscheidung des Architects in §6 des Plans (Zeile *Klartext ohne TLS*, Ort §4), nicht dem Wortlaut der DoD (V-165).

### Punkt 3: `README.md` nennt, was der Code-Slice liefert, im Ist-Zustand. **Bestätigt.**

Der überholte Satz („die Datenbank muss den Benutzer ohne Passwort anmelden“) ist ersetzt: `record` verlangt es weiter (Probe oben: `PGR-E6001` bei SCRAM, MD5, Klartext), `play` meldet sich mit Klartext, MD5 und SCRAM-SHA-256 an, mit dem Passwort aus dem Platzhalter oder aus `PGWIRE_RECORDER_PASSWORD` (Proben oben). Kein Stand, keine Chronik; die Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand. Die README nennt keine Aussage über SASLprep oder Klartext, die Proben bräuchten.

### Weitere DoD-Punkte (Gates, Closure)

- **`make gates` grün: bestätigt.** Eigener Lauf am Stand `465f09b` auf sauberem Baum: Exit 0, einschließlich `docs-check`, `kopf-check`, `lint`, `run-integration-tests: gruen`, `lint-gegenprobe: gruen`. Der Lauf ändert den Arbeitsbaum nicht (`git status --short` leer).
- **Review durchgeführt: bestätigt.** Report liegt vor (F-597 bis F-600), Abschnitt 2.
- **Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen: offen, noch nicht fällig.** Die Häkchen der DoD stehen leer, der Slice liegt in `in-progress/`, §7 trägt Belege, aber noch keine Closure-Notiz, keinen Lerneintrag und keine Risiko-Ausgänge; beide Risiken in §6 stehen auf „offen bis Closure“. Das ist Sache der Closure; Hinweise für die Ausgänge in Abschnitt 3.

## 2. Plan-vs-Diff und Review-Übergaben

**Plan gegen Diff.** Der Diff berührt `docs/user/benutzerhandbuch.md`, `README.md` und den Plan; kein Produkt-Code, kein Hilfetext, kein Test, kein Gate, nicht `benutzerhandbuch-standard.md` oder `abdeckung-*.md` (§1 *Ausdrücklich NICHT*). Jede Zeile der §3-Tabelle ist im Diff getragen: Handbuch §1, §4 (Voraussetzung, Hinweise), §5 (Absatz hinter der Optionstabelle, Beispiel, Platzhalter-Absatz), §6 (Nachbarsatz) und §7 (Zeilen `PGR-E4005`, `PGR-E6001`). Jede Zeile der Randform-Tabelle in §6 hat ihren Ort im Handbuch; die zwei bedingten Zeilen (*SCRAM ohne Channel Binding*, *Verfahren, das `play` nicht kann*) sind durch den Test belegt und stehen in §7. Die Zeilen *Frist, Iterationen* und *Wortlaut der Datenbank* sind ungenannt, wie entschieden (`grep -i` auf `iteration` und `10 000 000` in Handbuch und README: kein Treffer). §6 *Zuweisungen* sagt „keine“; der Diff weist keinem Slice etwas zu (`AGENTS.md` §3.13).

**Befehlsfolge (Plan §3 mit benannter Ausnahme), selbst gefahren** am Binary `ver-anmdoku-bin:dev` und am Handbuch `465f09b`:

| Prüfung | Ausgabe |
|---|---|
| Optionen des Handbuchs gegen `--help` aller vier Kommandos | nur `--h`, `--help`, `--version`, `--name`, `--rm` (wie in §3 genannt) |
| Variablen gegen Optionen, roh | `PGWIRE_RECORDER_PASSWORD` |
| Variablen mit `grep -vx PGWIRE_RECORDER_PASSWORD` | leer |
| Codes gegen den Katalog | leer |

Mutanten auf Kopien des Handbuchs (je ein angehängter Satz): `--upstream-tls` → Optionen melden ihn, rot; `PGWIRE_RECORDER_PASSWORT` → Variablen mit Abzug melden ihn, rot; `PGWIRE_RECORDER_FORCE_X` → rot; `PGWIRE_RECORDER_PASSWORD_FILE` → rot (der Abzug greift nur auf den vollen Namen); `PGR-E4007` und `PGR-W9999` → Codes melden sie, rot. Der Mutant „Zeile `PGWIRE_RECORDER_PASSWORD` zurück“ entfällt für diese Variable, wie §3 und §7 sagen.

**Review-Übergaben.**

| Befund | Stand im Plan und im Handbuch | Urteil |
|---|---|---|
| F-597 (MEDIUM, Befehlsfolge ohne benannte Ausnahme) | §3 *Ansatz* nennt die Ausnahme, die Variablen-Zeile mit `grep -vx` und die Folge „Variablen leer“; §6 trägt die Zeile *Befehlsfolge … benannte Ausnahme*; `done/slice-doku-ist-stand.md` ist unverändert (`git diff e8fd0b2..HEAD` auf die Datei: leer). Der Lauf oben ergibt genau die Ausgabe, die §3 nennt | behoben |
| F-598 (LOW, Nachbarsatz §6 als Randform) | Zeile *Nachbarsatz in §6 Rollen und Rechte* in §6 mit Wortlaut und Beleg; der Wortlaut im Handbuch stimmt zeichengleich überein; Proben oben tragen ihn | behoben |
| F-599 (INFO, Verweisziel) | §1 und §4 verweisen auf `#5-einstellungen`, in dem der Absatz mit Quelle und Vorrang steht; der Anker löst auf (`make docs-check` grün); der Absatz verweist zurück auf `#konfigurationsdatei` für den Passwortteil | behoben |
| F-600 (INFO, SASLprep nur Probe) | Satz mit „kann“; die Probe habe ich wiederholt (Punkt 1); kein Test, wie §7 sagt; `AGENTS.md` §3.11 gebrochen wird nicht (der Satz sagt nicht mehr zu, als die Probe zeigt) | wie eingestuft, offen als Beobachtung |

**Namensabgleich der in §7 genannten Tests und Dateien (`grep`).** `TestAnmeldungNichtUnterstuetzt` (`einspielen_anmeldung_test.go` Zeile 601), `TestAnmeldungWeitereAnforderungArt` (Zeile 678), `TestEinspielAufbauFehler` (`einspielen_test.go` Zeile 178), `envPassword` (`internal/adapters/driving/cli/cli.go` Zeile 86), `internal/hexagon/model/fehler.go`, `done/slice-doku-ist-stand.md` §7, `tools/test/run-integration-tests.sh` (Benutzer `play_scram`, `play_md5`, `play_pw`, Passwörter mit `GEHEIM`), `harness/mk/integration.mk`, der E2E-Test des Code-Slice (`test/integration/play_anmeldung_e2e_test.go`) und die Register-Einträge `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` und `BEO-REPO/verhalten-nur-unter-linux-geprueft`: alle vorhanden. Die Kennung `slice-v1-abschluss-einspielen-anmeldung` liegt in `done/`, die Abgrenzung dort (Zeilen 53 und 75) nennt diesen Slice als Adresse (`AGENTS.md` §3.13: die Adresse nimmt an, `slice-v1-abschluss-einspielen-anmeldung-doku` §1 und DoD führen den Teil).

## Findings

Ausgewiesene Befunde ab V-165; die Quelle ist die Regel, gegen die geprüft wurde.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-165 | INFO | Die DoD (Punkt 2) verlangt, dass §5 *Konfigurationsdatei* die Grenzen *SASLprep* und *Klartext ohne TLS* nennt; der Klartext-Satz steht nach der Entscheidung in §6 des Plans in §4 *Hinweise*, SASLprep in §5 *Einstellungen*. Der Handbuch-Teil entspricht §6, der Wortlaut der DoD nicht. | `AGENTS.md` §3.9 | Plan · „und nennt die Grenzen *SASLprep* und *Klartext ohne TLS*“ (DoD Punkt 2) gegen Handbuch §4 *Hinweise* | nein — Urteil | DoD-Ort weicht von der Randform-Entscheidung ab |
| V-166 | INFO | Der Klartext-Satz („sendet `play` es über die Verbindung, wie sie ist, also ohne TLS unverschlüsselt“) ist durch das Gelingen der Anmeldung über die unverschlüsselte Verbindung und die Ablehnung von `sslmode=require` belegt, nicht durch einen Mitschnitt des gesendeten Passworts; der Satz sagt nicht weiter, als diese Proben zeigen. Er steht auf einer einzigen langen Zeile neben umbrochenen Nachbarn, ebenso der neue README-Satz hinter „Beim Beenden warten …“ (Zeile 31). | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „Verlangt die Datenbank ein Klartext-Passwort, sendet `play` es …“ | nein | Grenze nur durch Probe belegt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Handbuch §1 *Voraussetzungen*, §4 *Voraussetzung*, *Hinweise*, *Vorgehen* | geprüft, ohne Befund (außer V-165, V-166) |
| Handbuch §5 Absatz zum Passwort, Optionstabelle, Beispiel `test`, Platzhalter-Absatz | geprüft, ohne Befund |
| Handbuch §6 *Rollen und Rechte* für `record`, `replay`, `play` | geprüft, ohne Befund |
| Handbuch §7 Zeilen `PGR-E2005`, `PGR-E2006`, `PGR-E4005`, `PGR-E6001` und *Die Anwendung kann sich nicht verbinden*, FAQ, Glossar, §10 | geprüft, ohne Befund |
| Handbuch gesamt auf Falsches, Chronik, Zielstand, Verweise (`grep` s. o.) | geprüft, ohne Befund |
| `README.md` | geprüft, ohne Befund |
| Meldungscodes gegen Katalog; Optionen gegen `--help`; Variablen gegen Optionen mit benannter Ausnahme | geprüft, ohne Befund |
| `GEHEIM` in Ausgaben (`debug`, `info`, `config show`, `PGR-E2006`) | geprüft, ohne Befund (0 Treffer, Aufnahme mit Gegenkontrolle) |
| Plan §1, §3, §6 gegen den Diff (`AGENTS.md` §3.9); Adressen (`AGENTS.md` §3.13) | geprüft, ohne Befund |
| Hard Rules §3.10, §3.12: kein neuer Vertrag, kein Code im Diff | geprüft, ohne Befund |
| Offen und nicht von mir gefahren | die Mutanten des Implementers am Test `TestAnmeldungNichtUnterstuetzt` (verlassen mich auf den grünen Gate-Lauf und die Namen); Windows, macOS und andere Serverversionen als PostgreSQL 17; `SCRAM-SHA-256-PLUS`-Server, Kerberos, GSSAPI, SSPI gegen einen echten Server (nur Test) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** DoD-Ort weicht von der Randform-Entscheidung ab · Grenze nur durch Probe belegt

## 3. Hinweise für die Closure (Risiken aus §6)

- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`: der Teil von `slice-v1-abschluss-einspielen-anmeldung` ist hier vollständig geliefert (alle in dessen DoD genannten Stellen: §1, §4, §5, §7, README, jeweils belegt); als Ausgang bietet sich *entfallen* mit dieser Begründung an, der Planner urteilt.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft`: alle Proben liefen nur auf Linux im Container, auch meine; *weiter offen* (Register).
- Finding-Klasse „Grenze nur durch Probe belegt“ (F-600, V-166) geht als Zeile in das Register der Closure; ob ein Test für SASLprep nötig ist, entscheidet der Planner (§1 schließt Produkt-Tests hier aus).

## Verdikt

**Merge-blockierend:** nein. Alle drei Liefer-Punkte der DoD sind bestätigt, `make gates` ist am Stand `465f09b` grün, F-597 bis F-600 sind behoben oder wie eingestuft.

**Übergabe:** V-165 und V-166 gehen an den Planner (V-165: Wortlaut der DoD an die Randform-Entscheidung angleichen oder den Ort bewusst hinnehmen; V-166 nur zur Kenntnis und für das Register). Kein Befund geht an den Implementer. Offen bleiben die Closure-Pflichten (Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen) und der Lifecycle-Übergang nach `done/`. Dieser Report ist ein Lauf-Beleg; er ersetzt weder das Review noch die Validierung.
