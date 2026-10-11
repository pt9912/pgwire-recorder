# Verifikation: slice-v1-abschluss-einspielen-tls-doku — 2026-10-11

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-11-review-slice-v1-abschluss-einspielen-tls-doku.md`, F-609 bis F-613); die Übergaben prüfe ich nach. Ich repariere nichts.

**Review-Art:** Verifikation (DoD- und Entscheidungs-Konformität, Plan gegen Diff) in frischem Kontext.

**Gegenstand:** Slice-Plan `slice-v1-abschluss-einspielen-tls-doku` (Kopf und §1 bis §8) gegen den Gesamt-Diff `0f9daf6..810ee49` (4 Dateien: `README.md`, `docs/user/benutzerhandbuch.md`, der Plan, der Review-Report; +242 −19). Architect `15a789a`; Implementer `eea9569`, `7b2663a`; Review `aab7ef2`; Nacharbeit `be0a8aa`, `810ee49`.

**Skill:** `.harness/skills/reviewer.md` (nicht angewandt; Rolle ist Verifier)
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- Plan ganz am Stand `810ee49`, besonders DoD (§2), §3, §6 (Randformen, Entscheidung vom 2026-10-11 Option A, Probe-Matrix, Risiken) und §7 (*Belege des Implementers*, Probe 16, Befehlsfolge, *Stand der Übergabe*).
- `AGENTS.md` §3.9, §3.11 (inklusive *Handbuch und README beschreiben den Ist-Zustand*), §3.13; `LH-FA-20.a` *TLS* (über den Kopf des Plans, nicht neu bewertet); `slice-doku-ist-stand` §6 und §7 (Befehlsfolge).
- Review-Report F-609 bis F-613.

## 1. DoD-Liefer-Punkte

**Methode.** Das Binary stammt aus einer frischen Kopie (`git archive HEAD`, `810ee49`, Stufe `runtime`, Image `ver-tlsdoku-bin:dev`). Gegenstellen im eigenen Docker-Netz (IPv4 und IPv6), alles Docker-only, Scratch-Präfix `ver-tlsdoku-`: `postgres:17-alpine` mit dem Digest aus `harness/mk/integration.mk` als `pgssl` (`ssl=on`, `hostnossl … reject`), `pgexp` (Blatt gültig 2020-01-01 bis 2020-01-02), `pgplain` (`ssl=off`), `pgpw` (`ssl=on`, `password`, `hostnossl … reject`), `pgpwplain` (`ssl=off`, `password`), `pgip` (Blatt ohne die Adresse des Servers), `pgchain` (sendet Blatt und Zwischenzertifikat), `pgleaf` (sendet nur das Blatt); Zertifikate mit `openssl` 3.5 im alpine-Hilfscontainer (zwei CAs, `Root`/`Mid`/`Other`); Python-Fake mit `N`, `Z`, Ende ohne Antwort und `S` mit Schweigen. Aufzeichnung neu mit `record` gegen `pgplain` (`DROP TABLE IF EXISTS probe`, `CREATE TABLE probe(id int)`, `INSERT`), „Tabelle“ heißt `select count(*) from probe` gleich 1 nach `DROP` vor dem Lauf. Das Beispiel `tls.yaml` ist wörtlich per `awk` aus dem Handbuch geschnitten, nicht abgetippt.

### Punkt 1: Handbuch (und README) beschreiben, was `slice-v1-abschluss-einspielen-tls` liefert — **bestätigt**

Jede geänderte Aussage zu TLS einzeln (Handbuch §1, §4 *Einspielen*, §5, §7; README):

| Aussage | Probe | Ergebnis |
|---|---|---|
| Beispiel `tls.yaml` (§4) mit `sslmode=require` und `--upstream-ca ./ca.pem` | (1), Host `localhost` im Netzwerk-Namensraum von `pgssl` | Exit 0, Tabelle; ohne `--upstream-ca`: Exit 4, `PGR-E4005`, „unknown authority“ |
| TLS-Wunsch: `--upstream-tls` (Option, Umgebung, Schlüssel `upstream_tls` in `play:`), `sslmode=require` | (1), (2), (12) | Exit 0 und Tabelle gegen `pgssl` (das Unverschlüsseltes ablehnt, also mit TLS); Schlüssel `upstream_ca` und Umgebung ebenso; `sslmode=disable` plus `--upstream-tls` Exit 0 |
| Gesetztes `--upstream-tls` geht `sslmode` vor, auch `false`, aus jeder Quelle | (12) gegen `pgplain` und `pgssl` | `false` als Schlüssel, Option, Umgebung: je Exit 0, Tabelle (kein TLS); Umgebung `false` und Option `=true`: Exit 4, „lehnt TLS ab“ (Option vor Umgebung); `--upstream-tls=false` mit `require` gegen `pgssl`: Exit 4, „no encryption“ |
| Ohne beides kein TLS; ohne Wunsch kein TLS | (12), (9) | `host:port` gegen `pgplain`: Exit 0, Tabelle; gegen `pgssl` ohne Wunsch: Exit 4, `PGR-E4005`, `pg_hba … no encryption` (auch `sslmode=disable`) |
| Jede Sitzung baut TLS auf | Aufzeichnung mit zwei Sitzungen gegen `pgssl` | Exit 0; im Log des Servers genau zwei Zeilen `SSL enabled` mehr |
| Kein Überspringen; unbekannte CA → `PGR-E4005` | (4) | `ca2.pem` und ohne CA: Exit 4, `PGR-E4005`, keine Tabelle |
| Name gegen den Host, wie eingesetzt (DNS, IPv4, IPv6); nicht im Blatt → `PGR-E4005` | (3), (5), Platzhalter `${DBH}` | `pgssl`, `172.30.99.10`, `[fd00:99::10]`, `${DBH}=pgssl`: Exit 0; `pgip` per Name, IPv4 und IPv6, `${DBH}` auf einen Namen außerhalb des Blatts: Exit 4 „valid for … not …“ |
| Abgelaufen → `PGR-E4005` | (6) per Name und per IP | Exit 4, „certificate has expired“ |
| `N` → `PGR-E4005`, kein Rückfall; Server lehnt Unverschlüsseltes ab → `PGR-E4005` | (7), (9) | Fake `N` und `pgplain` mit `require` und mit `--upstream-tls`: Exit 4 „lehnt TLS ab“, **keine Tabelle** (ein Rückfall hätte sie angelegt: `pgplain` nimmt Unverschlüsseltes an, Gegenprobe 9c Exit 0, Tabelle) |
| Weder `S` noch `N` / Ende davor → `PGR-E4002` | (8) | Fake `Z`: Exit 4, `PGR-E4002`; Ende ohne Antwort: Exit 4, `PGR-E4002` |
| `--upstream-ca` ergänzt den Speicher, mehrere Zertifikate | (13), Bündel | `SSL_CERT_FILE=ca2` plus `--upstream-ca ca`: Exit 0; umgekehrt Exit 0; beide `ca2`: Exit 4; Bündel `ca2`+`ca`: Exit 0 |
| `--upstream-ca` ohne TLS → `PGR-E2001` | (10) | `host:port`, `sslmode=disable`, und `--upstream-tls=false` mit `require`: je Exit 2, `PGR-E2001` |
| Datei nicht lesbar / keine reguläre Datei / kein Zertifikat → `PGR-E2007`, Exit 2, ohne Pfad | (11) | fehlende Datei, Rechte 000, Verzeichnis, leere Datei, Textdatei, PEM nur mit Schlüssel, ungültiger Wert über die Umgebung: je Exit 2, `PGR-E2007`, Meldung nennt `--upstream-ca` und Grund |
| `SSL_CERT_FILE` zählt zum Speicher des Systems | (13) am Binary der Stufe `runtime` | `SSL_CERT_FILE=ca.pem` ohne `--upstream-ca`: Exit 0, Tabelle, auch mit Datei und `require`; ohne die Variable Exit 4 |
| Klartext-Passwort: mit TLS auf der verschlüsselten, ohne TLS auf der unverschlüsselten Verbindung | (14) | `pgpw` mit TLS: Exit 0, Server-Log `method=password` mit `SSL enabled (protocol=TLSv1.3 …)`; `pgpwplain`: Exit 0, `method=password` ohne `SSL enabled`; `pgpw` ohne TLS: Exit 4 „no encryption“; falsches und fehlendes Passwort: Exit 4, `PGR-E4005` |
| Meldungen ohne Pfad, Inhalt, Passwort | `grep -n -i -E` über alle 181 Ausgabezeilen auf Dateinamen, `/ex`, `BEGIN`, `MIIB`, beide Passwörter; `grep -n -F` mit je zwei Zeilen Inhalt von `ca.pem` und `ca2.pem` | kein Treffer, außer dem Wort „PEM-Block“ (kein Pfad); Gegenprobe des Musters auf einen Text mit Pfad: 1 Treffer, Gegenprobe des `-F`-Musters auf seine eigene Quelle: 2 Treffer |
| `record` + `sslmode=require` → `PGR-E2004`; `record` und `replay` nehmen TLS nicht an; `replay` hat keine Verbindung | (15) | `record`: Exit 2, `PGR-E2004`, „sslmode=require ist bei record ungültig“; `replay` mit ungenutzter `require`-Verbindung startet; `psql sslmode=require` gegen `replay` und `record`: „server does not support SSL“; `sslmode=prefer` gegen `record`: Verbindung ohne TLS zu `pgssl` abgelehnt („no encryption“); `replay.upstream`: „unbekannter Schlüssel“ (`PGR-E2004`) |
| Ungenutzte `require`-Verbindung bleibt gültig | `play` mit `a` (require) und `b` (disable), benutzt `b` | Exit 0, Tabelle |

**Probe 16, der neue Satz zur Kette** („eine Zertifizierungsstelle aus seiner Kette … Zwischenzertifikate, die der Server mitsendet“), 7 Fälle selbst gefahren, Name per `--add-host` auf das Blatt:

| Server sendet | `--upstream-ca` | Ergebnis (selbst) | Tabelle des Implementers |
|---|---|---|---|
| Blatt und Zwischenzertifikat | Wurzel | Exit 0, Tabelle | gleich (a) |
| Blatt und Zwischenzertifikat | nur Zwischenzertifikat | Exit 0, Tabelle | gleich (b) |
| Blatt und Zwischenzertifikat | fremde CA | Exit 4, `PGR-E4005`, „unknown authority“ | gleich (c) |
| Blatt und Zwischenzertifikat | keine | Exit 4, `PGR-E4005` | gleich (d) |
| nur Blatt | Wurzel | Exit 4, `PGR-E4005`, „unknown authority“ | gleich (e) |
| nur Blatt | nur Zwischenzertifikat | Exit 0, Tabelle | gleich (f) |
| nur Blatt | Bündel Wurzel und Zwischenzertifikat | Exit 0, Tabelle | gleich (g) |

Mutanten des Satzes: „die Zertifizierungsstelle, die es ausgestellt hat“ (Vorfassung) wäre durch die Zeilen 2 und 6 widerlegt; „nur die Wurzel genügt immer“ durch Zeile 5. Zusatzprobe: Das Blatt selbst als `--upstream-ca` wird angenommen (Exit 0); der Satz nennt das Zertifikat des Servers ausdrücklich als Teil der Kette, die Probe widerspricht ihm nicht.

**Befehlsfolge (Plan §3, mit benannter Ausnahme)** selbst gefahren gegen `--help` des frisch gebauten Binaries und `internal/hexagon/model/fehler.go`: Optionen des Handbuchs, die `--help` nicht kennt: `--h --help --name --rm --version` (Nicht-Optionen bzw. `docker run`); Variablen mit `grep -vx PGWIRE_RECORDER_PASSWORD`: leer (ohne die Ausnahme: genau `PGWIRE_RECORDER_PASSWORD`); Codes außerhalb des Katalogs: leer; beide Gegenrichtungen (Hilfe → Handbuch, Katalog → Handbuch): leer. Mutanten an Kopien des Handbuchs: angehängtes `--upstream-cert` → in der Options-Ausgabe; `PGWIRE_RECORDER_UPSTREAM_TLSS` (Tippfehler) → in der Variablen-Ausgabe; `PGR-E2008` → in der Code-Ausgabe; `PGR-E2007` durch `PGR-E2099` ersetzt → `PGR-E2099` in der Code-Ausgabe. Je rot. Zusätzlich: `--upstream-tls` und `--upstream-ca` stehen nur in `play --help`, nicht in `record`, `replay`, `config show` (Betriebsart `play` der zwei Tabellenzeilen stimmt).

**Aktive Suche** (Handbuch §1, §4, §5, §6, §7 ganz gelesen, README ganz): `grep -n -i` im Handbuch nach Slice, Welle, ADR-, LH-, SPEC-, ARC-, `spec/`, Lastenheft, Spezifikation, Review, „noch nicht“, „kommt“, „geplant“, „später“, „bisher“, „künftig“, „vorerst“, „derzeit“, „bis zu“: Treffer nur Bestand (Zeilen 14, 153, 595) und Zeile 371 („gleich ob es als Option … kommt“, Herkunft des Werts; die Berichtigung in §7 trifft zu). Im Diff (`+`-Zeilen) kein Treffer auf Chronik- oder Zielstand-Wörter; README ohne Treffer auf „sicher“, „geschützt“, „noch nicht“, „bisher“. Keine Verweise auf Spezifikation, Lastenheft, ADRs, Slices, Wellen, Reviews im Handbuch. Nachbarsätze: §1 *Voraussetzungen* (Zeile zur Verbindung zum Werkzeug: „läuft ohne Verschlüsselung“), §4 *Mit einem Datenbanktreiber arbeiten* (`sslmode=prefer`/`disable`), §5 Passwort-Absatz, §6 *Rollen und Rechte*, §7 *Die Anwendung kann sich nicht verbinden* sind unberührt und stimmen mit den Proben (15) und (9) überein. Zwei Randbefunde: V-169, V-170.

**Urteil Punkt 1: bestätigt.** Jede geänderte TLS-Aussage trägt eine Probe; kein Satz musste entfallen.

### Punkt 2: §5 *Konfigurationsdatei* beschreibt TLS wie geliefert; Grenze *Zertifikatsspeicher nicht ladbar* nicht genannt; kein Test geändert — **bestätigt**

- §5-Absatz zu `sslmode` (`require` wünscht TLS mit Prüfung, `--upstream-tls` geht vor, auch `false`, Schlüssel `upstream_tls`, `upstream_ca` in `play:`): Proben (12), (12f).
- `--upstream-ca` und der beobachtbare Fall (weder Speicher noch Datei → `PGR-E4005`): Proben (4), (13), (16); das Handbuch nennt die Grenze *Zertifikatsspeicher nicht ladbar* an keiner Stelle (`grep -n -i 'nicht ladbar\|nicht geladen\|nicht lad\|gilt als leer'` in Handbuch und README: kein Treffer).
- Plan-vs-Diff: DoD-Zeile 2 trägt jetzt die Fassung „… den beobachtbaren Fall … (`PGR-E4005`); die Grenze *Zertifikatsspeicher nicht ladbar* nennt das Handbuch nicht (§6, Entscheidung vom 2026-10-11). Kein Test ist geändert, `make abdeckung-check` ist grün.“ (Diff `7b2663a..HEAD`, eine Zeile). §6 trägt die Entscheidung als eigenen Absatz (**Entscheidung des Nutzers vom 2026-10-11: Option A**) mit Gründen und der Verwerfung von Option B; der *Ungenannt*-Eintrag in §6 verweist auf die angeglichene Zeile. F-609 (Entscheidung ohne Artefakt, DoD folgt nicht) und F-613 (Angleichung) sind damit behoben.
- „Kein Test ist geändert“: `git diff --stat 0f9daf6..HEAD` zeigt vier Dateien, keine unter `test/`, `internal/`, `cmd/`; `make abdeckung-check` ist Teil von `make gates` (grün, Abschnitt 4).
- §7-Berichtigung der früheren Fehlbezeichnung: §7 *Abweichungen und Funde* nennt jetzt, dass der frühere Wortlaut „Entscheidung des Nutzers“ schon vor der Entscheidung falsch war („Es war eine Empfehlung des Architect; berichtigt“), und die DoD-Zeile ist angeglichen. Das ist die Berichtigung, die die Frage verlangt, und sie steht im Plan, nicht nur im Bericht. **Grenze meiner Prüfung:** Die Entscheidung selbst (Option A am 2026-10-11) kann ich nur gegen den Plan lesen; die Nutzer-Mitteilung liegt nicht im Repo (V-171).

**Urteil Punkt 2: bestätigt.**

### Punkt 3: README nennt, was geliefert wird, im Ist-Zustand — **bestätigt**

Der Satz „Alle Verbindungen laufen unverschlüsselt“ ist ersetzt; die neuen Teile sind je belegt: `record` und `replay` nehmen Verbindungen nur ohne TLS an (Probe 15: `psql sslmode=require` gegen beide „server does not support SSL“), `record` verbindet ohne TLS (Probe 15: `record` gegen `pgssl` mit `sslmode=prefer` des Clients: „no encryption“), `play` verbindet auf Wunsch (`--upstream-tls` oder `sslmode=require`) mit TLS (1, 2), prüft das Zertifikat auch gegen eine eigene CA (1, 4, 13), meldet Ablehnung von TLS und nicht passendes Zertifikat als `PGR-E4005` (4, 5, 6, 7, 9). Keine Chronik, kein Zielstand, kein „sicher“ oder „geschützt“; die Verweise auf `spec/`, `docs/plan/` und `docs/reviews/` bleiben Verweise nach oben ohne Aussage über einen Stand. Eine Randbeobachtung zur Formatierung: die neue README-Zeile „(`--upstream-ca`), und meldet …“ ist länger als ihre Nachbarn; kein Inhaltsfehler.

**Urteil Punkt 3: bestätigt.**

### Übrige DoD-Punkte

| Punkt | Urteil | Beleg |
|---|---|---|
| `make gates` grün | **bestätigt** | eigener Lauf am sauberen Baum `810ee49` (`git status` leer davor und danach): Ausgang 0, alle Fragmente bis `lint-gegenprobe` grün |
| Review durchgeführt, Report unter `docs/reviews/` | **bestätigt** | `docs/reviews/2026-10-11-review-slice-v1-abschluss-einspielen-tls-doku.md`; kein Self-Review (anderer Lauf, `aab7ef2`) |
| Closure-Notiz mit Lerneintrag | **offen, Pflicht des Planners** | §7 trägt *Belege des Implementers*, noch keine Closure-Notiz; das ist vor dem `git mv` nach `done/` fällig |
| Beobachtungs-Register fortgeschrieben | **offen, bei Closure** | Plan §7: „bei Closure“; keine Datei unter `evidence/` für diesen Slice (`ls`) |
| Jedes Risiko aus §6 trägt einen Ausgang | **offen, bei Closure** | beide Risiken stehen mit „**Ausgang:** offen bis Closure“ |
| Drei Paarungen (Anker · Folge-Slice · Register) | **offen, bei Closure** | nicht Teil der Liefer-Punkte; sind bei der Closure zu prüfen |

Die vier Closure-Pflichten sind keine Abweichung vom Plan: Der Plan schiebt sie ausdrücklich auf die Closure. Sie sind Eingang für den Planner (Abschnitt 5).

## 2. Übergaben des Reviews (F-609 bis F-613)

| ID | Übergabe | Stand am HEAD | Urteil |
|---|---|---|---|
| F-609 (MEDIUM) | Entscheidung zur DoD-Zeile 2 ohne Artefakt; DoD folgt nicht | §6 trägt die Entscheidung (Option A, Gründe, Option B verworfen), DoD-Zeile 2 ist angeglichen, §7 berichtigt die Fehlbezeichnung | behoben; Rest siehe V-171 |
| F-610 (LOW) | „Zertifizierungsstelle, die es ausgestellt hat“ weiter als die Probe | Handbuch §4 und §7 nennen die Kette; Probe 16 (7 Fälle) selbst gefahren, deckt sich | behoben |
| F-611 (INFO) | Beleg „keiner in den neuen Sätzen“ ungenau | §7 berichtigt: Zeile 371 ist neu, im Sinn der Herkunft des Werts | behoben |
| F-612 (INFO) | Reihenfolge `--config`/`--log-level` in der Tabelle gegen die Hilfe | Bestand, unverändert; §7 führt ihn als „Bestand (F-612)“ | Übergabe an den Planner bleibt offen (keine Aktion im Diff) |
| F-613 (INFO) | DoD-Angleichung bei Closure | die Zeile ist bereits angeglichen (siehe Punkt 2) | behoben |

## 3. Hexagon und Architektur-Gate

Nicht berührt: Der Diff enthält keinen Produkt-Code, keinen Test, kein Gate und keine Spec-Stelle (`git diff --stat 0f9daf6..HEAD`: `README.md`, `docs/user/benutzerhandbuch.md`, Plan, Review-Report). `make a-check` und `make a-check-negativ` sind Teil des grünen `make gates`.

## 4. Plan gegen Code

| Prüfung | Ergebnis |
|---|---|
| Kopf: `Berührte Spec-Stellen: —` | stimmt; Diff ändert keine Spec-Stelle; `make kopf-check` grün |
| §1 Ziel und Abgrenzung: nur Handbuch, README, Plan | stimmt; `docs/user/benutzerhandbuch-standard.md` und `docs/user/abdeckung-*.md` unverändert |
| §3 Plan: zwei Dateien (`benutzerhandbuch.md`, `README.md`) | stimmt; Plan und Report sind Slice-Artefakte |
| §6 Randformen: *Genannt* | jede Randform trägt eine Zeile in der Probe-Matrix und ist im Handbuch genannt (TLS-Wunsch, Prüfung, Name, `--upstream-ca`, `SSL_CERT_FILE`, Klartext-Passwort, Meldungscodes, Voraussetzungen); neu in §6: Zeile „Kette“ (Probe 16) ist im Handbuch |
| §6 Randformen: *Ungenannt* | im Handbuch nicht zu finden (TLS-Version, Handshake-Frist, Zone der IPv6-Adresse, `S` mit Schweigen, MD5/SCRAM mit TLS, Grenze *Speicher nicht ladbar*): `grep` ohne Treffer |
| Entscheidung vom 2026-10-11 Option A | umgesetzt wie entschieden; im Handbuch nur der beobachtbare Fall |
| §7 Namensabgleich | `TestE2EPlayTLSZertifikatsspeicherDesSystems` in `test/integration/play_tls_e2e_test.go` (Zeile 116), `harness/mk/integration.mk`, `internal/hexagon/model/fehler.go`, `slice-doku-ist-stand`, `slice-v1-abschluss-einspielen-anmeldung-doku`, `slice-v1-abschluss-einspielen-tls`, `welle-v1-abschluss` vorhanden; `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, `…/verhalten-nur-unter-linux-geprueft`, `…/zusage-im-kommentar-weiter-als-pruefung`, `…/schichtteilung-je-plan-verschieden` als Verzeichnisse vorhanden; Zeilen 371 und 595 im Handbuch stimmen mit den Zitaten in §7 |
| Liefer-Punkte ≤ 3, Schichten ≤ 2 | 3 Punkte (Handbuch Verhalten, Handbuch §5 TLS, README), eine Schicht (Doku) |
| `make docs-check` | grün (Ende des Laufs, siehe Verdikt) |
| `make gates` am HEAD `810ee49` | grün, Ausgang 0 |
| Beleg „`make gates` auf `be0a8aa` / `eea9569` grün“ (Implementer, §7) | nicht nachfahrbar; meine Messung am HEAD ersetzt sie, der Beleg war vorhanden |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-169 | LOW | Der Kopf des Handbuchs trägt weiter `Stand: 10.10.2026`, obwohl dieser Slice das Handbuch am 2026-10-11 ändert; frühere Änderungen an vier verschiedenen Tagen (04., 05., 06., 10.10.) haben die Zeile nachgezogen (`git log -L5,5`). | `AGENTS.md` §3.7 (Zustandsfeld nennt den Zustand), Maintainability | `docs/user/benutzerhandbuch.md` · „Stand: 10.10.2026“ | ja (`git log -L5,5:docs/user/benutzerhandbuch.md`) | Stand-Datum im Handbuch-Kopf nicht nachgezogen |
| V-170 | INFO | Die Klammer in `PGR-E2001` („weder `--upstream-tls` noch `sslmode=require` der Verbindung“) lässt den Fall `--upstream-tls=false` mit `sslmode=require` und `--upstream-ca` aus, den das Binary mit `PGR-E2001` ablehnt (Probe 10, Exit 2); der Vorrang steht in §4, der Satz ist als Ganzes wahr. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „oder `--upstream-ca` steht ohne TLS (weder `--upstream-tls` noch `sslmode=require` der Verbindung)“ | ja (Probe 10c) | Klammer enger als das beobachtbare Verhalten |
| V-171 | INFO | Die Entscheidung des Nutzers vom 2026-10-11 (Option A) liegt im Repo nur als Absatz in Plan §6; die Mitteilung selbst ist kein Artefakt. Der Verifier kann Option A gegen den Plan lesen, nicht gegen die Quelle. Zuständig: Planner (Quelle bei der Closure nennen, falls sie ein Artefakt hat). | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-tls-doku.md` · „Entscheidung des Nutzers vom 2026-10-11: Option A.“ | nein (Urteil) | Entscheidung des Nutzers nur im Plan, Quelle nicht prüfbar |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| DoD-Punkte 1 bis 3 (Liefer-Punkte), jede geänderte TLS-Aussage von Handbuch und README einzeln gegen das gebaute Binary und echte Gegenstellen | geprüft, ohne Befund außer V-169 und V-170 |
| Beispiel `tls.yaml` wörtlich aus dem Handbuch | geprüft, läuft (Exit 0, Tabelle) |
| Kettensatz gegen die Probe-16-Tabelle (7 Fälle) | geprüft, alle 7 gleich der Tabelle des Implementers |
| Befehlsfolge Plan §3 mit benannter Ausnahme, drei Mutanten plus Ersatz-Mutant | geprüft, ohne Befund; je Mutant rot |
| Chronik, Zielstand, Verweise auf Spezifikation, Lastenheft, ADRs, Slices, Wellen, Reviews (Handbuch ganz, README ganz) | geprüft, ohne Befund |
| Meldungen ohne Pfad, Inhalt, Passwort (mit Gegenprobe des Musters) | geprüft, ohne Befund |
| Rückfall nach `N` auf Unverschlüsselt | geprüft, ohne Befund (Gegenprobe zeigt, dass ein Rückfall eine Tabelle anlegte) |
| Review-Übergaben F-609 bis F-613 | geprüft, F-612 offen als Bestand beim Planner |
| Namensabgleich aller in §7 genannten Tests, Dateien, Pläne, Beobachtungs-Verzeichnisse | geprüft, ohne Befund |
| Produkt-Code, Tests, Gates, Spec | nicht im Diff, geprüft, ohne Befund |
| `make gates` am HEAD | grün (Ausgang 0) |
| Nicht geprüft | Zone einer IPv6-Adresse, `S` mit Schweigen, MD5/SCRAM mit TLS, der Standardpfad des Systems ohne `SSL_CERT_FILE` (kein Zertifikat einer öffentlichen CA gegen das Image), Windows und macOS; die Sätze dazu stehen nicht im Handbuch (`BEO-REPO/verhalten-nur-unter-linux-geprueft`: alle Proben unter Linux im Container) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Stand-Datum im Handbuch-Kopf nicht nachgezogen · Klammer enger als das beobachtbare Verhalten · Entscheidung des Nutzers nur im Plan, Quelle nicht prüfbar

## Verdikt

**Merge-blockierend:** nein. Die drei Liefer-Punkte der DoD sind bestätigt, `make gates` und `make docs-check` sind grün. Offen sind nur die Closure-Pflichten, die der Plan ausdrücklich auf die Closure schiebt.

**Übergabe an den Planner:**

- **V-169** (LOW) an den Implementer: Stand-Datum im Handbuch-Kopf nachziehen; keine Blockade der Closure.
- **V-170** (INFO) an den Implementer, ohne Pflicht; **V-171** (INFO) an den Planner (Quelle der Entscheidung bei der Closure nennen); **F-612** (Reihenfolge der Tabelle gegen die Hilfe, Bestand) bleibt beim Planner.
- Für die Closure: Closure-Notiz mit Lerneintrag; Register fortschreiben (Klassen der Review: Entscheidung des Nutzers ohne Artefakt im Plan · Zusage über die Zertifizierungsstelle weiter als die Probe · Beleg ungenauer als der Lauf · Reihenfolge der Tabelle gegen die Hilfe; Klassen dieses Laufs oben); Ausgänge der zwei Risiken aus §6: `BEO-REPO/verhalten-nur-unter-linux-geprueft` (alle Proben unter Linux im Container, auch meine) bleibt *weiter offen* und gehört ins Register; `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` ist an diesem Slice nicht eingetreten (der Handbuch- und README-Teil ist geliefert und hier bestätigt), der Ausgang ist Urteil des Planners.
- Dieser Report ist ein Lauf-Beleg; Aufbau und Scratch (Container, Netz, Image, Verzeichnisse mit Präfix `ver-tlsdoku-`) sind entfernt.
