# Verifikation: slice-v1-abschluss-einspielen-extended-doku — 2026-10-10

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-10-review-slice-v1-abschluss-einspielen-extended-doku.md`, F-587 bis F-590). Ich prüfe die Übergaben nach (Abschnitt 2), ordne F-588 und F-589 ein und trage keine Reparatur bei.

**Gegenstand:** Slice-Plan `slice-v1-abschluss-einspielen-extended-doku` (Kopf, §1 bis §8) gegen den Gesamt-Diff `93382d1..8892421`. Architect: `2e010b5` (Randformen und Probe-Matrix in §6). Implementer: `880522e` (Handbuch und README), `8e72827` (Gate-Lauf), `2cb0519` (Nacharbeit zu F-587, Probe 6), `8892421` (Gate-Lauf der Nacharbeit). Review: `bec21ad`. Rahmen: [MR-000](../../harness/conventions.md#mr-000--baseline-aussage).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingang:**

- der Plan ganz am Stand `8892421`, besonders §1 (*Ausdrücklich NICHT*), §2 (DoD), §6 (*Randformen*, *Probe-Matrix*, *Risiken*) und §7 (*Belege des Implementers*)
- `docs/user/benutzerhandbuch.md` und `README.md` ganz am Stand `8892421`; der Diff beider Dateien gegen `93382d1`
- `AGENTS.md` §3.9, §3.11 (auch *Handbuch und README beschreiben den Ist-Zustand*), §3.13; [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler) (beschriebenes Verhalten, nicht geändert)
- `slice-doku-ist-stand` §6 und §7 (Randformen und Befehlsfolge für Doku-Slices); der Katalog `internal/hexagon/model/fehler.go`; `--help` des gebauten Binaries
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen, nur die Belege in Plan §7

Beim Start war der Arbeitsbaum sauber, HEAD `8892421`. Der Commit dieses Berichts enthält nur diese Datei.

**So wurden die Proben gebaut.**

- **Binary:** frische Kopie per `git archive HEAD | tar -x` in das Scratchpad, daraus die Stufe `runtime` des `Dockerfile` mit eigenem Tag (`ver-extdoku-bin:dev`). Das Image ist `sha256:b8abba88772c…`, dieselbe Kennung wie im Beleg des Implementers (der Code ist seit `2e010b5` unverändert).
- **Datenbank:** das gepinnte Image aus `harness/mk/integration.mk` (`postgres:17-alpine@sha256:b0f9560a…`, Anmeldung `trust`) in einem eigenen Netz `ver-extdoku-net`; je Probe eine eigene Datenbank mit `t (n int PRIMARY KEY, s text)`. Client für Probe 1 ist `psql` (`\bind`) aus demselben Image hinter `record`; die Aufzeichnungen der Proben 2 bis 6 habe ich selbst von Hand geschrieben (nicht die des Implementers), `play` lief als Container mit `timeout 60`, der Signal-Lauf mit `docker wait` unter `timeout 20`.
- **Danach:** Container, Netz, Image und die Scratch-Verzeichnisse `ver-extdoku-w` und `ver-extdoku-src` (nur diese benannten Pfade) entfernt; `docker ps -a`, `docker network ls`, `docker images` und `ls` des Scratchpads mit `ver-extdoku` danach je 0 Treffer. Die Log-Datei des Gate-Laufs (`ver-extdoku-gates.log`) bleibt bis zum Ende der Sitzung im Scratchpad liegen.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-einspielen-extended` liefert (§4 `play`, §7 `PGR-E6001`; Sätze zu `--continue-on-error`, `--allow-recorded-errors`, Abbruchsignal). **Bestätigt.**

Probe-Matrix aus Plan §6, Zeile für Zeile, gegen das frische Binary und PostgreSQL 17:

| Matrix | Aufruf (kurz) | Ergebnis, selbst beobachtet | Urteil |
|---|---|---|---|
| 1 | `record` + `psql \bind \g` (zwei `INSERT`), `play --database` einer frischen Datenbank | Exit 0, Tabelle `1:eins,2:zwei` | bestätigt |
| 2 | Interaktion 1: Flush-Gruppe `INSERT (1)` (Schlüssel vorhanden), Sync-Gruppe `INSERT (2)`; Interaktion 2 `INSERT (3)` | ohne Option: Exit 4, ein `PGR-E4004` („Session 1, Interaktion 1“, 23505), Tabelle `1:vorhanden`. Mit `--continue-on-error`: Exit 4, ein `PGR-E4004`, Tabelle `1:vorhanden,3:naechste` (Zeile 2 fehlt). Kontrolle ohne Vorbelegung: Exit 0, `1,2,3` | bestätigt |
| 3 | `--allow-recorded-errors`, Fehlerantwort (23505) nur in der Sync-Gruppe aufgezeichnet, Fehler tritt in der Flush-Gruppe auf | Exit 0, kein `level=ERROR`, Tabelle `1:vorhanden,3:naechste`; dieselbe Aufzeichnung ohne Option: Exit 4, `PGR-E4004`; Aufzeichnung ohne Fehlerantwort mit der Option: Exit 4, `PGR-E4004`; aufgezeichnet `XX000` statt `23505`: Exit 0. Zusätzlich von mir: Fehlerantwort nur in der Flush-Gruppe aufgezeichnet, Fehler tritt in der Sync-Gruppe auf: Exit 0; Fehlerantwort nur in einer *anderen* Interaktion aufgezeichnet: Exit 4 (der Satz sagt „der aufgezeichneten Folge“) | bestätigt |
| 4 | `SIGTERM` nach 1,5 s während `pg_sleep(4)` in `Execute` | 4a: Exit 0, Log `Abbruchsignal, play endet vorzeitig`, `play` endet nach der Interaktion (2,6 s nach dem Signal, Docker-Aufschlag eingerechnet), Tabelle `10:warten` (weder Folge 2 noch Sitzung 2). 4b: Sleep in der Flush-Gruppe: `10:warten,20:gruppe-b`, die Sync-Gruppe läuft noch, `21` fehlt. 4c: `--finish-session-on-interrupt`: `10:warten,11:zweite`, `12` fehlt. 4d: zweites Signal 0,5 s später: Exit 0, Ende 0,16 s danach, kein `level=ERROR`, mit und ohne Option | bestätigt (auch der unveränderte Satz zum zweiten Signal gilt für Extended) |
| 5 | `COPY t FROM STDIN` und `COPY t TO STDOUT` in der zweiten von drei Interaktionen | je Exit 6, `PGR-E6001` („Session 1, Interaktion 2“, `CopyInResponse` bzw. `CopyOutResponse`), Tabelle `5:vorher` (Folge 3 nicht gelaufen), auch mit `--continue-on-error` | bestätigt |
| 6 (Nacharbeit) | Aufzeichnung von Hand: `select pg_terminate_backend(pg_backend_pid())` in Interaktion 1, Folge 2 danach | 6a (`FATAL` 57P01 aufgezeichnet): ohne Option, mit `--continue-on-error`, mit `--allow-recorded-errors`, mit beiden je Exit 4, genau ein `PGR-E4003` („Session 1, Interaktion 1“, 57P01), kein `PGR-E4004`, Folge 2 nicht gelaufen. 6b (ohne Fehlerantwort aufgezeichnet): mit und ohne `--allow-recorded-errors` Exit 4, `PGR-E4003`. Die Angabe des Implementers, `record` ließe sich für diese Anfrage nicht fahren, stimmt: `record` endet mit `PGR-E6001` („Verbindung zum Upstream nach der Fehlerantwort 57P01 … vor ReadyForQuery beendet“), die Aufzeichnung trägt keine Fehlerantwort | bestätigt |

Sätze des Diffs, einzeln (Handbuch §4 und §7):

| Satz | Probe | Urteil |
|---|---|---|
| *Voraussetzung*: „Eine Aufzeichnung liegt vor, und die Zieldatenbank …“ („mit einfachen Anfragen“ entfällt) | 1 | bestätigt |
| „Antwortet die Datenbank auf eine Anfrage des erweiterten Protokolls mit einem Fehler, führt sie die übrigen Nachrichten dieser Folge bis zu deren `Sync` nicht aus.“ | 2 (Zeile 2 fehlt) | bestätigt |
| „Mit `--continue-on-error` macht das Einspielen danach mit der nächsten Folge weiter, der Exit-Code ist dann 4.“ | 2 | bestätigt |
| „Bei einer Folge des erweiterten Protokolls genügt eine Fehlerantwort an irgendeiner Stelle der aufgezeichneten Folge.“ | 3 (beide Richtungen, anderer SQLSTATE) | bestätigt |
| Unveränderter Satz „Ein Fehler mit dem Schweregrad `FATAL` beendet die Verbindung und bricht das Einspielen auch dann ab (`PGR-E4003`).“, gelesen auch für Extended | 6a, 6b | bestätigt (Beleg nur Probe, kein automatisierter Test, V-161) |
| „Antwortet die Datenbank mit einem COPY-Datenstrom (`COPY … FROM STDIN`, `COPY … TO STDOUT`), kann `play` ihn nicht verarbeiten und endet mit `PGR-E6001` (Exit-Code 6); die Meldung nennt Sitzung und Nummer der Anfrage.“ | 5; für einfache Anfragen trägt es `TestE2EPlayFortsetzungAbbruch` (`COPY t FROM STDIN`, Exit 6) | bestätigt, Begriff siehe V-160 |
| „Bei `Strg+C` oder `SIGTERM` endet das Einspielen nach der laufenden Anfrage, bei einer Folge des erweiterten Protokolls nach deren `Sync`; mit `--finish-session-on-interrupt` erst nach der laufenden Sitzung.“ | 4a bis 4c | bestätigt |
| §7 Zeile `PGR-E6001`: „… oder die Datenbank antwortet bei `play` mit einem COPY-Datenstrom.“ | 5 | bestätigt |

Handbuch-Regeln, selbst gefahren:

- **Befehlsfolge aus `slice-doku-ist-stand` §7** (am frischen Binary): Optionen des Handbuchs, die `--help` nicht nennt: nur `--h`, `--help`, `--version` (im Handbuch als Nicht-Optionen genannt) und `--rm`, `--name` (`docker run`). Variablen ohne Option: keine. Codes außerhalb des Katalogs: keine. Umgekehrt (Katalog ohne Handbuch): leer, `PGR-E1000` steht jetzt im Handbuch.
- **`grep -n -i`** über das Handbuch auf `spec/|spezifikation|lastenheft|slice|welle|review|\bADR|LH-|SPEC-|noch nicht|kommt|geplant|künftig|bald|später|demnächst|derzeit|zunächst|vorerst|bisher|ab Version|seit`: nur Bestand ohne Bezug (`später` in §1 „spielt sie später ohne die Datenbank wieder ab“, „existiert noch nicht“ in §3, „kommt“ in der Ruhefrist, „ab Version 17“ bei der Leerraum-Regel von PostgreSQL). Kein Verweis auf Spezifikation, Lastenheft, ADR, Slice, Welle oder Review.
- **Beispiel:** Das Handbuch trägt keinen neuen Codeblock. Der Block in §4 (`play --upstream postgres:5432 --input ./recordings/users.yaml`) hat die Form meiner Probe 1 (`--upstream`, `--input`, hier mit `pg:5432` und `--database`); er lief in der Verifikation von `slice-doku-ist-stand` wörtlich als Datei.
- **Optionen gegen `--help`:** `play --help` ist unverändert zum Diff; Tabelle §5 (Zeilen `--continue-on-error`, `--allow-recorded-errors`, `--finish-session-on-interrupt`) nennt dieselben Variablen und Standardwerte.

### Punkt 2: `README.md` nennt, was geliefert ist, im Ist-Zustand; überholte Sätze ersetzt (auch der Untertitel); Verweise ohne Aussage über einen Stand. **Bestätigt.**

| Satz | Probe | Urteil |
|---|---|---|
| Untertitel: „… führt ihre Anfragen erneut gegen eine Datenbank aus“ („einfachen“ entfällt) | 1 (Extended), einfache Anfragen tragen `TestE2EPlayFortsetzungAbbruch` und die Verifikation von `slice-doku-ist-stand` | bestätigt |
| „`play` spielt die Anfragen einer Aufzeichnung, einfache wie vorbereitete Anweisungen, gegen eine Datenbank ein“ | 1 | bestätigt |
| „Antwortet die Datenbank mit einem COPY-Datenstrom, endet `play` mit `PGR-E6001`.“ | 5 | bestätigt |
| „nach einer Fehlerantwort der Datenbank bricht es ab, mit `--continue-on-error` läuft es weiter, und mit `--finish-session-on-interrupt` …“ (Bestand, umgebrochen) | 2, 4c | bestätigt |
| Satz „Eine Aufzeichnung mit vorbereiteten Anweisungen lehnt `play` ab“ | entfernt (`grep -n "lehnt" README.md`: 0 Treffer) | bestätigt |

`grep -n -i -E "slice|welle|noch|künftig|geplant|später|bisher|derzeit|vorerst|seit |kommt|aktuell|Version"` auf `README.md`: zwei Treffer, „bekommt“ in „bekommt nie eine geratene Antwort“ (Bestand ohne Stand-Aussage) und „Slices“ in der Zeile zu `make doc-trace` (Bestand, F-589). Die Verweise auf `spec/`, `docs/plan/` und `docs/reviews/` stehen ohne Aussage über einen Stand. Eine Zeile des Satzes „Alle Verbindungen laufen unverschlüsselt …“ ist beim Umbrechen nicht neu umbrochen (V-159).

### Punkt 3: `make gates` grün. **Bestätigt durch eigenen Lauf.**

Eigener Lauf im Arbeitsbaum am Stand `8892421`, sauber vor dem Lauf: Exit 0. Darin `d-check: 487 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.0 OK — 54 Dateien`, `gesamt: 0 Befund(e)` (a-check), `a-check-negativ: gruen`, `abdeckung-gegenprobe: gruen`, `commit-msg-gegenprobe: gruen`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `run-integration-tests: gruen`. `git status` nach dem Lauf: sauber bis auf diesen Bericht. Der Beleg des Implementers (`make gates` am Stand `2cb0519`, Exit 0, in Plan §7) deckt sich; zwischen `2cb0519` und `8892421` ändert sich nur Plan §7.

### Übrige DoD-Punkte (konstant je Slice)

- **Review:** Der Report liegt vor (`bec21ad`), aus einem anderen Lauf. Die Nacharbeit `2cb0519` folgt ihm. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen:** noch offen; die DoD-Kästchen stehen alle auf `[ ]`, §6 trägt bei beiden Risiken „offen bis Closure“. Das ist die Reihenfolge, kein Befund. Für die Closure aus meiner Sicht: Risiko 1 (Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut, `BEO-REPO/verhalten-nur-unter-linux-geprueft`): kein neuer Codeblock im Diff, die Sätze laufen als Proben im Container gegen das gepinnte Image; der Ausgang ist *weiter offen* (alle Proben, auch meine, laufen nur unter Linux im Container). Risiko 2 (Teil nicht vollständig geliefert, `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`): nicht eingetreten; die Stellen aus DoD 1 und 2 (Voraussetzung, Hinweise, `PGR-E6001`-Zeile, README-Sätze, Untertitel) sind ersetzt, und die Suche nach „einfache Anfragen“, „Folge des erweiterten Protokolls“ als Ablehnungsgrund und „lehnt“ findet nichts Verbliebenes; Ausgang *entfallen* mit dieser Begründung oder *weiter offen* als Evidenz. Beide Entscheidungen gehören dem Planner.

---

## 2. Übergaben des Reviews

| Befund | Prüfung | Urteil |
|---|---|---|
| F-587 Satz zu `FATAL` bei Extended ohne Beleg | Probe 6 (Matrix) selbst gefahren: `PGR-E4003`, Exit 4, in allen sechs Läufen, unabhängig von Option und davon, ob die Aufzeichnung eine Fehlerantwort trägt; Handbuch bleibt unverändert. Die Angabe, `record` könne diese Aufzeichnung nicht erzeugen, ist gegengeprüft (`PGR-E6001` beim Aufzeichnen). Der Beleg ist eine Probe, kein Test; die E2E-Tests tragen `FATAL` nur für einfache Anfragen (V-161) | behoben, Probe trägt |
| F-588 „Interaktion“ in der Meldung gegen „Nummer der Anfrage“ im Handbuch | Eingeordnet: Die Meldungen nennen „Session 1, Interaktion 2“ (Probe 2, 5, 6; auch bei einfachen Anfragen), das Handbuch „Nummer der Anfrage“ (§4, §7 Zeile `PGR-E4004`). Das Handbuch ist am Binary nicht falsch, die Nummer ist dieselbe; die Begriffe weichen ab (V-160) | eingeordnet, Hand-off an den Planner |
| F-589 README „Slices“ in der Zeile zu `make doc-trace` | Eingeordnet wie im Review: Bestand, `AGENTS.md` §3.11 verbietet den Verweis auf Slices im Handbuch, dem README verbietet es Chronik und Zielstand; die Zeile erklärt die Ausgabe eines Targets | eingeordnet, Hand-off an den Planner |
| F-590 Fund „spielt nichts ein“ war falsch | Probe 5: Folgen vor dem Copy laufen, danach nichts, Exit 6, auch mit `--continue-on-error`; die neue Fassung („kann `play` ihn nicht verarbeiten und endet mit `PGR-E6001`“) stimmt | bestätigt |

---

## 3. Plan gegen Code

- **§1:** Ziel und Abgrenzung folgen dem Diff. Der Diff ändert genau drei Dateien außer dem Review-Report: `docs/user/benutzerhandbuch.md`, `README.md` und diesen Plan; kein Produkt-Code, kein Hilfetext, kein Test, kein Gate, kein `docs/user/benutzerhandbuch-standard.md`, keine `abdeckung-*.md`.
- **§2 / DoD:** Punkt 1 und 2 sind im Diff erfüllt (oben). DoD 1 nennt als Sätze: Voraussetzung, Hinweis zu Fehler und `--continue-on-error`, `--allow-recorded-errors`, Abbruchsignal, `PGR-E6001`-Hinweis, §7-Zeile; alle stehen im Diff. Die DoD verlangt „jeder Meldungscode im Katalog … und ausgelöst durch einen Test oder eine Probe“: `PGR-E4003`, `E4004`, `E6001` sind ausgelöst (Proben 2, 5, 6).
- **§3:** Die Tabelle führt zwei Dateien; der Diff ändert zwei (plus Plan und Report).
- **§6:** Die Randformen folgen dem Diff. *Wortlaut der Grenzen*: Der Wortlaut des COPY-Satzes ist wörtlich der aus §6. *Ungenannt, weil nicht belegbar*: Warten je Antwort, Gegendruck, `CopyBothResponse`, nicht lesbare Antwort, Aufzeichnung mitten in einer Interaktion: keine dieser Aussagen steht in Handbuch oder README. *Serverversionen*: keine Version genannt. *Anmeldung und TLS*: Sätze zu `PGR-E4005` und `sslmode=require` im Diff unverändert. *Mutation*: entfällt wie entschieden. *Größe*: zwei Liefer-Punkte, eine Schicht.
- **§7:** Die Belege nennen Stand, Aufbau, Probe 1 bis 6, neue und entfernte Aussagen, Ungenanntes, Befehlsfolge, `grep`, Gate-Läufe. Die Probe-Tabelle stimmt mit meinen Beobachtungen überein (Exit-Codes, Zeilen, Codes, Zeiten bis auf den Docker-Aufschlag).
- **§8:** Liefer-Punkte 2, Schicht eine; stimmt mit dem Diff.
- **Commits (`AGENTS.md` §3.3):** sechs Commits, kein Move, jeder nennt `slice-v1-abschluss-einspielen-extended-doku` und `LH-FA-20`.
- **Adressen (`AGENTS.md` §3.13):** Der Diff weist keinem Slice etwas Neues zu; `slice-v1-abschluss-einspielen-anmeldung-doku` und `slice-v1-abschluss-einspielen-tls-doku` bleiben die Adressen für Passwort und `sslmode`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-159 | INFO | Zwei Absätze sind beim Umformulieren nicht neu umbrochen: In der README steht der Satz „… ohne Passwort anmelden. Beim Beenden warten `record` und `replay` höchstens `--shutdown-timeout`“ auf einer überlangen Zeile; im Handbuch §4 steht nach „Schweregrad `FATAL` beendet die Verbindung“ ein Zeilenumbruch mitten im Satz. Der gerenderte Text ist richtig. | `AGENTS.md` §3.11 (Form, kein Inhalt) | README.md · „ohne Passwort anmelden. Beim Beenden warten“; docs/user/benutzerhandbuch.md · „Schweregrad `FATAL` beendet die Verbindung“ | ja — Sicht auf die Zeilenlänge | Umbruch nach Satzersetzung nicht nachgezogen |
| V-160 | INFO | Zu F-588: Die Meldungen von `PGR-E4004`, `PGR-E4003` und `PGR-E6001` nennen „Session 1, Interaktion 2“ (am Binary beobachtet), das Handbuch sagt „Nummer der Anfrage“ (§4 Hinweis zu `PGR-E6001`, §7 Zeile `PGR-E4004`). Das Handbuch führt „Interaktion“ nicht als Begriff. Wer die Zahl in der Meldung sucht, findet sie ohne Erläuterung. | Maintainability | docs/user/benutzerhandbuch.md · „die Meldung nennt Sitzung und Nummer der Anfrage“ | nein | Begriff der Meldung weicht vom Begriff des Handbuchs ab |
| V-161 | INFO | Der Handbuch-Satz zu `FATAL` und `PGR-E4003` gilt für Extended laut Probe 6 (selbst bestätigt), aber kein automatisierter Test trägt ihn: `grep` nach `FATAL`/`57P01` in `test/integration/play_extended_e2e_test.go` und `internal/hexagon/services/play_extended_test.go` findet nichts, `play_laufsteuerung_e2e_test.go` trägt `FATAL` nur für einfache Anfragen. Plan §6 nimmt das ausdrücklich an (Mutation entfällt, Probe-Matrix ersetzt sie); wird der Code geändert, fängt nichts diese Zusage. | `AGENTS.md` §3.11 (Handbuch sagt nur zu, was ein Test oder eine Probe prüft; hier: Probe) | docs/user/benutzerhandbuch.md · „Schweregrad `FATAL` beendet die Verbindung“ | ja — E2E-Test mit `FATAL` in einer Extended-Interaktion | Zusage im Handbuch ruht auf einer einmaligen Probe |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Probe-Matrix Plan §6, Zeilen 1 bis 5 und Probe 6 (Nacharbeit), eigene Aufzeichnungen, frisches Binary, PostgreSQL 17 | geprüft, ohne Befund; jede neue Aussage trägt eine bestandene Probe |
| Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen, Codes) | geprüft, ohne Befund (Optionen nur `--h`, `--help`, `--version`, `--rm`, `--name`; Variablen und Codes leer) |
| Handbuch §4 `play` ganz (Voraussetzung, Vorgehen, Ergebnis, alle Hinweise) | geprüft, ohne Befund; zusätzlich Gegenproben zu Richtung und Stelle der Fehlerantwort (3f, 3g) bestätigen die Sätze |
| Handbuch §1 (Betriebsarten, Voraussetzungen), §5 (Tabelle `play`, `--user`, `--database`, Verbindung, Passwort), §6 (Log-Stufen bei `play`), §7 (Katalog, alle Zeilen zu `play`) | geprüft, ohne Befund; Katalog und Handbuch decken sich in beide Richtungen |
| README ganz: Untertitel, „Was kann ich heute tun?“, „Kerngedanke“, Verweise | geprüft, ohne Befund außer V-159 (Form) und F-589 (Einordnung) |
| Aktive Suche nach Falschem, Zielstand, Chronik und Verweisen (`grep` Handbuch und README, Wörter wie „noch“, „kommt“, „geplant“, „später“, „bisher“, „derzeit“, „Slice“, „Welle“, `ADR`, `LH-`) | geprüft, ohne Befund (Treffer nur Bestand ohne Bezug, siehe Punkt 1 und 2) |
| Nicht belegbare Aussagen (Warten je Antwort, Gegendruck, Serverversion, „übrige Gruppen ohne Warten“, `CopyBothResponse`, Aufzeichnung endet mitten in der Interaktion) | geprüft, ohne Befund (keine steht in Handbuch oder README) |
| Plan folgt dem Diff (`AGENTS.md` §3.9), Randformen vor dem Code (`AGENTS.md` §3.12: `2e010b5` vor `880522e`) | geprüft, ohne Befund |
| Adressen (`AGENTS.md` §3.13), Commit-Reihenfolge und Moves (`AGENTS.md` §3.3) | geprüft, ohne Befund (kein Move, keine neue Zuweisung) |
| Mutation (`AGENTS.md` §3.10) | entfällt wie in Plan §6 entschieden (kein neuer Vertrag, nur Beschreibung) |
| Grenzen Anmeldung und TLS (Passwort bei `play`, `sslmode=require`) | geprüft, ohne Befund; Sätze unverändert, Adressen `slice-v1-abschluss-einspielen-anmeldung-doku` und `slice-v1-abschluss-einspielen-tls-doku` |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Umbruch nach Satzersetzung nicht nachgezogen · Begriff der Meldung weicht vom Begriff des Handbuchs ab · Zusage im Handbuch ruht auf einer einmaligen Probe

## Verdikt

**Merge-blockierend:** nein — kein HIGH, MEDIUM oder LOW. DoD-Liefer-Punkte 1 bis 3 sind bestätigt: Handbuch, README und `make gates`. Die Review-Übergaben F-587 und F-590 sind belegt, F-588 und F-589 eingeordnet.

**Falsches gefunden?** Ich habe aktiv gesucht (§4 `play` ganz, §7 Meldungscodes, §5 Tabelle `play`, §1, Log-Stufen, README Untertitel und `play`-Sätze, jeweils gegen das Binary und mit zusätzlichen Gegenproben zu Stelle und Richtung der aufgezeichneten Fehlerantwort, Fehlerantwort einer anderen Interaktion, Signal in der Flush-Gruppe, zweites Signal mit und ohne Option). Keine Aussage in Handbuch oder README war am gebauten Binary falsch. Die Stellen mit Schwäche stehen in V-159 bis V-161.

**Übergabe:** Dieser Bericht geht an den Planner. V-159 geht zur Kenntnis an den Implementer (Form, vor der Closure nachziehbar, ohne Auftrag). V-160 geht mit F-588 an den Planner (Entscheidung, ob das Handbuch „Interaktion“ einführt oder die Meldung „Anfrage“ sagt; keine Änderung an diesem Slice). V-161 geht an den Planner (Entscheidung, ob ein E2E-Test für `FATAL` in Extended in einen Code-Slice gehört oder die Probe genügt; ohne Änderung an diesem Slice). F-589 geht ebenfalls an den Planner. Für die Closure: die Ausgänge der beiden Risiken aus Plan §6 (oben, Abschnitt 1, *Übrige DoD-Punkte*) und die Klassen von Review und Verifikation in §7. Der Bericht ersetzt weder das Review noch die Validierung.
