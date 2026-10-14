# Verifikation: slice-v1-abschluss-anmeldung-doku

**Rolle:** Verifier. **Datum:** 2026-10-14. **Spanne:** `ea5019b` (Move nach `in-progress/`) bis `b529096` (HEAD). **Gegenstand:** [slice-v1-abschluss-anmeldung-doku](../plan/planning/in-progress/slice-v1-abschluss-anmeldung-doku.md) gegen DoD (§2), den geschlossenen Code-Slice [slice-v1-abschluss-anmeldung](../plan/planning/done/slice-v1-abschluss-anmeldung.md), [`LH-FA-05`](../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) und in der Spezifikation [`LH-FA-05.b`](../../spec/spezifikation.md#lh-fa-05b--authentifizierung) (nach `6a44529` und `0aadcac`), `AGENTS.md` §3.9, §3.11, §3.13, §3.14 und Modul 11. Review: [Review-Report](2026-10-14-review-slice-v1-abschluss-anmeldung-doku.md).

Alle Belege unten sind in diesem Lauf selbst erzeugt, nicht aus §7 übernommen. Aufbau: Binary `pgwire-recorder:dev` (`make build` auf HEAD; `git diff ea5019b..HEAD -- internal cmd test harness .claude tools Makefile` ist leer, der Code ist seit dem Code-Slice unverändert), `postgres:17-alpine@sha256:b0f9560a…` aus `harness/mk/integration.mk` (17.11), `psql` 17.11 aus demselben Image, `python:3.13-alpine` für Rohclient und Ersatz-Server, OpenSSL 3.3.7 (`alpine:3.20`) für das Zertifikat. Eigenes Docker-Netz `vrf-net`, ein Lauf zur Zeit, kein `pkill`; danach Container und Netz entfernt (`docker ps -a` leer, `docker network ls` nur `bridge`, `host`, `none`). Scratch-Dateien liegen außerhalb des Repos; `git status` war nach den Läufen sauber.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren) |
|---|---|---|
| 1 Benutzerhandbuch beschreibt, was der Code-Slice liefert; Optionen gegen `--help`, Codes gegen Katalog und ausgelöst, Beispiele als Datei, keine Verweise (`grep`) | bestätigt, mit V-1 (LOW) | Abschnitte *Beispiele als Datei*, *Verhaltensaussagen*, *Befehlsfolge*, *Ist-Zustand und Verweise*. V-1: ein Satz ist weiter gefasst als das Binary. |
| 2 `README.md` nennt, was der Code-Slice liefert, im Ist-Zustand; Verweise nach oben ohne Stand | bestätigt | Diff `README.md` (4 + 3): `record` vermittelt Klartext, MD5, SCRAM-SHA-256 „geprüft mit PostgreSQL 17“ (P1), zeichnet das Passwort nicht auf (P1, P2). Der alte Satz „Datenbank muss den Benutzer ohne Passwort anmelden“ steht nirgends mehr. Verweise auf `spec/` und `docs/plan/` (Zeilen 13, 96, 98) unverändert im Bestand, ohne Aussage über einen Stand. |
| 3 `make gates` grün | bestätigt mit Einschränkung der Kette | `echo '{}' \| bash .claude/hooks/stop-require-gates.sh` antwortet auf HEAD `{"decision":"approve"}` (vorher und nachher); der Hook verlangt keinen vollen Lauf. Selbst gefahren: `make gates-schnell` auf HEAD Exit 0, `make docs-check` `d-check: 745 Datei(en) geprüft, 0 Befund(e)`, `make kopf-check` Exit 0. Der volle Lauf auf `30b9b49` (§7) ist von mir nicht neu gefahren. |
| 4 Review durchgeführt | bestätigt | Report liegt vor (0 HIGH, 1 MEDIUM, 2 LOW, 4 INFO); Behebung siehe *Review-Findings*. |
| 5 bis 8 Closure-Notiz, Register, Risiko-Ausgänge, Paarungen, Häkchen | offen (erwartet) | Alle Häkchen leer; die drei Risiken in §6 stehen auf „offen bis Closure“; §7 trägt die Belege des Implementers und keine Closure-Notiz. Vorschläge zu den Ausgängen unten. |

## Beispiele als Datei gegen das gebaute Binary (AGENTS.md §3.11)

Die zwei Blöcke des Abschnitts *Eine Anwendung mit Passwort-Anmeldung aufzeichnen* per `awk` unverändert aus dem Handbuch gezogen. Block 1 mit `pgwire-recorder` als Funktion (`docker run -d …`), `postgres:5432` → `vrf-pg:5432`, `./recordings/` → `/recordings/`; Block 2 mit `localhost` → `vrf-rec` und `psql` als Funktion (Container, `PGPASSWORD=geheim` aus dem Block). Ergebnis: Zeile `1`, Exit 0; `login.yaml` trägt `user: app`, `database: app`, `sql: select 1`; `grep -c -E 'geheim|SCRAM-SHA|n,,n=|biws'` ergibt 0.

## Verhaltensaussagen der neuen und geänderten Stellen

| Aussage | Lauf | Befund |
|---|---|---|
| Klartext, MD5, SCRAM-SHA-256 und `trust` mit demselben Aufruf von `record`; 4 Sitzungen | P1 (`--log-level debug`) | vier Zeilen, `sql: select current_user` 4-mal |
| Passwort und Nachrichten der Anmeldung weder in Datei noch Log, auch `debug` | P1, P11 | Muster `GEHEIM\|md5[0-9a-f]{32}\|SCRAM-SHA\|n,,n=\|biws`: 0 Treffer in Datei und Log; Gegenprobe des Musters auf fünf Spuren: 5 Treffer |
| `replay` nimmt jedes Passwort an | P3 (frische Datei) | ohne Passwort `1`; mit `anders` `1`; anderer Benutzer und Datenbank `1` |
| `record` liest weder Benutzer/Passwort/Datenbank der Verbindung aus `--upstream` noch `PGWIRE_RECORDER_PASSWORD` | P5 (Konfigurationsdatei mit `anderer`/`andere`, Variablen gesetzt) | beide Anmeldungen gelingen, Aufzeichnung trägt `pw_scram`/`pw_md5` und `anmd`; `record --help` nennt keine Option dazu |
| Fehlerantwort der Datenbank unverändert (falsches Passwort in drei Verfahren, `reject`, unbekannter Benutzer) | P4 gegen direkten Lauf, P12 | Texte gleich dem direkten Lauf; Klartext-Fehlerantwort byte-gleich (`cmp`) |
| keine Meldung, Exit-Code unverändert; keine Sitzung, auch mit `--record-empty-sessions` | P4 | Log nur „gestartet“, „Herunterfahren“, „beendet“, Exit 0, `empty_sessions: true`, `sessions: []` |
| Schließen, Rücksetzen, `Terminate` in der Anmeldung | P8 (Rohclient) | drei Wege, Log ohne Fehlerzeile, Exit 0, in der Datei nur die folgende Sitzung (`id: 1`, `select 'danach'`) |
| Channel Binding: `require` bricht ab mit und ohne `--tls-cert`; `prefer` gelingt | P7 | ohne TLS „channel binding required, but SSL not in use“; mit TLS „server did not offer an authentication method that supports channel binding“; `prefer` je `1`; `pw_md5` mit `sslmode=require` `1`; 2 Sitzungen |
| Anforderung eines anderen Verfahrens: `PGR-E6001`, Anwendung erhält die Meldung, nicht aufgezeichnet, Exit 6 | P6 (Ersatz-Server, Codes 2, 6, 7, 9, 42; Code 7 auch mit `--record-empty-sessions`) | je „Anmeldeverfahren (Code N), das record nicht vermittelt“, Log-Zeile `PGR-E6001`, Exit 6, `sessions: []` |
| Nachricht außerhalb des Ablaufs `PGR-E4002`, Ende in der Anmeldung `PGR-E4003`, Exit 4 | P6b | Code 11 als erste Nachricht und `AuthenticationOk` dann Code 3: `PGR-E4002` mit Code in der Meldung; Server schließt nach der Antwort: `PGR-E4003`; je Exit 4, `sessions: []` |
| Zeit der Anmeldung zählt zu `offset_ms` | P10 | 1,5 s Wartezeit vor dem Passwort: `offset_ms: 1504`; ohne Wartezeit `5` |
| Kerberos/GSSAPI/SSPI nur als Reaktion des Werkzeugs, mit Ersatz-Server | P6 | Handbuch sagt „Ersatz-Server“ und „mit einem Server, der Kerberos, GSSAPI oder SSPI spricht, ist sie nicht geprüft“; Codes 2, 7, 9 stammen vom Python-Ersatz. Trägt. |
| PostgreSQL 17 / libpq 17 | Aufbau | `postgres --version` 17.11, `psql --version` 17.11; Handbuch nennt „PostgreSQL 17“ an allen Stellen zur Anmeldung (§1, §4, §6, neuer Abschnitt, README) und sagt „andere Versionen und andere Treiber nicht geprüft“ |
| Zusätzliche Probe (Randform F-7, nicht im Plan) | Ersatz-Server `garbage` (R mit Länge 2) und Rohclient mit `Q` in der Anmeldung | beide `PGR-E6001` („Nachricht des Upstreams im Verbindungsaufbau nicht lesbar“; „Client-Nachricht … in der Anmeldung wird nicht unterstützt“), Exit 6; siehe V-1 |

## Befehlsfolge des Plans §3

Selbst nachgefahren (`--help` aus dem gebauten Binary für Hauptkommando, `record`, `replay`, `play`, `config show`; Katalog `internal/hexagon/model/fehler.go`; Variablen-Zeile mit `grep -vx PGWIRE_RECORDER_PASSWORD`), dazu die drei Gegenrichtungen:

| Stand | Ergebnis |
|---|---|
| Handbuch `ea5019b` (vor dem Slice) | Optionen nur `--h`, `--help`, `--name`, `--rm`, `--version`; Variablen, Codes und alle Gegenrichtungen leer. **Nicht rot**, anders als in der Vorgabe der Aufgabe erwartet: Der Slice führt keine Option, Variable und keinen Code ein, den das Handbuch vorher nicht nannte (`PGR-E4002`, `E4003`, `E6001` standen im Katalog und in der Tabelle). Die Befehlsfolge kann diesen Slice weder rot noch grün von vorher unterscheiden; der Plan führt auch keine „rot vorher“-Zeile (V-3). |
| Handbuch HEAD | wie oben, leer bis auf die benannten Nicht-Optionen |
| Mutant (Kopie mit angehängter Zeile `--tls-ca PGWIRE_RECORDER_TLS_CA PGR-E2099`) | je rot: `--tls-ca`, `PGWIRE_RECORDER_TLS_CA`, `PGR-E2099` |

Die Befehlsfolge belegt hier nur, dass der Slice keinen Namen einführt, den Hilfe und Katalog nicht kennen. Die inhaltlichen Aussagen tragen die Läufe oben.

## Ist-Zustand und Verweise (grep)

- `grep -n -i -E "slice|welle|adr-|LH-|SPEC-|ARC-|spec/|lastenheft|spezifikation|review|noch nicht|kommt|geplant|später|bisher|künftig|vorerst"` über die hinzugefügten Zeilen von Handbuch und README: kein Treffer. Über das Handbuch insgesamt auf Slice, Welle, ADR-, LH-, SPEC-, ARC-, `spec/`, Lastenheft, Spezifikation, Review: kein Treffer.
- README: Treffer nur Zeilen 13, 96, 98 (Verweise nach oben auf `spec/…`, `docs/plan/…`, `docs/reviews/`), im Bestand, nicht im Diff.
- Der neue Abschnitt verweist nur auf Handbuch-Anker (`#5-einstellungen`, `#eine-anwendung-mit-passwort-anmeldung-aufzeichnen`, `#mehrere-verbindungen-aufzeichnen-und-wiedergeben`).
- Versionszeile `Version: 0.8` (vorher 0.7), `Stand: 14.10.2026`.

## Quantoren-Durchsicht

Neu geprüft am Binary, soweit der Quantor eine Allgemeinheit behauptet:

- „weder in der Datei noch im Log, auch nicht mit `--log-level debug`“: bestätigt (P1, Gegenprobe des Musters 5 Treffer).
- „jedes Passwort“ bei `replay`: bestätigt (P3; Aussage über die Aufzeichnung einer Sitzung).
- „liest weder Benutzer, Passwort und Datenbank … noch `PGWIRE_RECORDER_PASSWORD`“: bestätigt (P5).
- „erhält Ihre Anwendung deren Fehlerantwort unverändert“: byte-gleich für Klartext (P12), Text gleich für SCRAM, MD5, `reject`, unbekannten Benutzer (P4). „Unverändert“ für SCRAM und MD5 ist nicht byte-genau belegt; die Spezifikation sagt „mit allen Feldern“ (Info).
- „Beides ist keine Sitzung, auch nicht mit `--record-empty-sessions`“: bestätigt (P4, P8, P6b).
- „Sendet die Datenbank in der Anmeldung eine Nachricht, die der Ablauf nicht vorsieht, ist das `PGR-E4002`“ (Hinweis und Tabellenzeile): **zu weit**, siehe V-1.
- Bestandssätze, die der Abschnitt einschränkt (§4 *Mehrere Verbindungen*, *Zeitangaben*, §5, §6, §1 vertrauliche Daten): Entscheidungen der Durchsicht in §7 stimmen mit den Läufen. Gesucht mit `grep -n -i -E "ohne passwort|passwort|anmeld"` über das Handbuch; keine übrig gebliebene Aussage „record verlangt eine Datenbank ohne Passwort“. Die Zeilen 1498 ff. (`Bei record zählen nur Host und Port`) bleiben wahr (P5).

## Review-Findings

- **F-1 (MEDIUM, Plan folgt der Spec-Änderung nicht):** behoben in `30b9b49`. Der Kopf nennt `LH-FA-05.b` (*Aufzeichnung*) mit `6a44529` und `0aadcac`, §1 sagt, wer die Änderung machte, §7 führt die Probe „die nicht trug“ und sagt, dass die Spezifikation zum Server nichts mehr zusagt. `make kopf-check` Exit 0. Beleg im Diff: `git diff ea5019b..HEAD -- spec/spezifikation.md` ändert genau den Satz zum Startup-Parameter.
- **F-2 (LOW, Satz zum Server in der Spezifikation):** behoben in `0aadcac`. Der Satz sagt jetzt „Der Recorder reicht sie unverändert an den Server und entfernt keinen (`LH-RB-01`)“; kein Verhalten des Servers. Die Zusage „entfernt keinen“ habe ich nicht mit einem zusätzlichen Startup-Parameter am Binary gefahren; sie steht in der Spezifikation (Architect), das Handbuch wiederholt sie nicht.
- **F-3 (LOW, Version nicht genannt):** behoben in `1291bcd`. Die Sätze in §1, §4, §6 nennen „geprüft mit PostgreSQL 17“, der Satz zu anderen Verfahren spricht vom Werkzeug und nennt den Ersatz-Server; README ebenso. Offen bleibt als Risiko 3 die Matrix weiterer Versionen (`slice-v1-abschluss-postgres-versionen`), korrekt in §6 geführt.
- **F-4 (INFO, „liegt in `next/`“ im Verifikations-Report des Code-Slice):** unverändert, Mechanik von `make slice-mv`; der Report ist Lauf-Beleg. Keine Aktion.
- **F-5 (INFO, DoD-Wortlaut):** Die Zählung der DoD (neun Stellen) stimmt mit dem Diff überein, wenn „§4 zwei Hinweise“ die zwei Hinweise *Eine Anwendung aufzeichnen* und *Mehrere Verbindungen* meint; keine Aktion.
- **F-6 (INFO, lange Zeile):** unverändert, keine Konvention verlangt es. Keine Aktion.
- **F-7 (INFO, Lücke bei nicht lesbarer Nachricht und fremder Clientnachricht):** als bewusste Lücke im Plan geführt (§1 *Ausdrücklich NICHT*, §7). Das trägt für die Clientseite (das Handbuch sagt zur Clientnachricht außer `p` und `Terminate` nichts; es steht auch kein Satz da, der ihr widerspricht). Es trägt **nicht** für die Serverseite: siehe V-1.

## Plan-vs-Code-Diff

`git diff ea5019b..HEAD --stat`: `README.md`, `docs/user/benutzerhandbuch.md`, der Plan, der Review-Report, `spec/spezifikation.md` (`LH-FA-05.b`, 4 Zeilen, Architect), und der Link im Verifikations-Report des Code-Slice (`make slice-mv`). Nichts unter `internal/`, `cmd/`, `test/`, `harness/`, `.claude/`, `docs/user/abdeckung-*.md`, `docs/user/benutzerhandbuch-standard.md`: die Ausschlüsse aus §1 sind eingehalten. Der Move `ea5019b` ist rein; Inhalt steht in eigenen Commits (`AGENTS.md` §3.3); die Inhalts-Commits seither nennen die Slice-Kennung. Jede Handbuch-Stelle des Diffs steht in §3 des Plans (neuer Abschnitt, §1, §4 mit zwei Hinweisen und Zeitangaben, §6, §7 mit `PGR-E4002`, `E4003`, `E6001` und *Die Anwendung kann sich nicht verbinden*, Versionszeile). Die Spezifikationsänderung steht im Kopf und in §1 (§3.9); §3 trägt `spec/spezifikation.md` nicht als Zeile, weil nicht der Implementer sie änderte (INFO). Der Code-Slice liegt in `done/`; die Voraussetzung aus §4 ist erfüllt. Liefer-Punkte 2 (Handbuch, README), Schichten eine (Dokumentation): bestätigt.

## Größe (AGENTS.md §3.14)

`git diff --numstat ea5019b..HEAD -- docs/user README.md`: Handbuch 83 + 15 = 98, README 4 + 3 = 7, Summe **105** (Marke 300, Meldeschwelle 220 nicht erreicht, Schätzung des Architect 120 bis 200 unterschritten). Entspricht §7 nach dem Review. Der Plan, der Review-Report und die Spezifikation liegen außerhalb der Metrik (Pfade `docs/user README.md`).

## Risiko-Ausgänge (Vorschläge für die Closure, kein Eintrag)

- *Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut* (`BEO-REPO/verhalten-nur-unter-linux-geprueft`): Beide Blöcke liefen in diesem Lauf unter Linux im Container; Windows und macOS sind nicht gefahren, das Handbuch nennt keine. Ausgang *weiter offen* ins Register (zweiter Beleg in `evidence/`, kein Zähler setzen); *entfallen* wäre falsch.
- *Slice liefert Handbuch-/README-Teil nicht vollständig* (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`): nicht eingetreten; die Befehlsfolge und die Stellenzählung (neun) stimmen, F-7 ist eine benannte Lücke. Ausgang *entfallen* mit dieser Begründung trägt, solange V-1 behoben oder bewusst angenommen ist.
- *Sätze über den Server nennen keine Version* (`slice-v1-abschluss-postgres-versionen`): „PostgreSQL 17“ steht an allen Stellen; die Matrix bleibt Gegenstand des Folge-Slice. Ausgang *eingetreten* mit der Kennung `slice-v1-abschluss-postgres-versionen` oder *weiter offen*; der Planner entscheidet (Nehmer-Prüfung nach `AGENTS.md` §3.13 vor dem Eintrag).

## Befunde

- **HIGH:** keine.
- **MEDIUM:** keine.
- **LOW (V-1):** Handbuch, Hinweis im neuen Abschnitt und Zeile `PGR-E4002` der Tabelle: „Sendet die Datenbank in der Anmeldung eine Nachricht, die der Ablauf nicht vorsieht, ist das `PGR-E4002`“. Eine nicht lesbare Nachricht des Servers im Aufbau ist nach `LH-FA-05.b` `PGR-E6001`, und das Binary zeigt es (Ersatz-Server mit Nachricht der Länge 2: `PGR-E6001`, „Nachricht des Upstreams im Verbindungsaufbau nicht lesbar“, Exit 6, nicht `PGR-E4002`). Der Satz mit dem Quantor „eine Nachricht, die der Ablauf nicht vorsieht“ schließt diesen Fall ein und ist für ihn falsch. Der Plan führt F-7 als „bewusste Lücke, das Handbuch sagt nichts dazu“; für die Serverseite stimmt das nicht, denn das Handbuch sagt etwas und es trifft hier nicht zu. Abhilfe in der Hand des Planners oder Implementers: den Satz auf die Fälle eingrenzen, die P6b belegt („eine Nachricht der Anmeldung außerhalb der Reihenfolge, deren Code die Meldung nennt“), oder den nicht lesbaren Fall mit einer Probe aufnehmen und `PGR-E6001` nennen. Die Clientseite (Nachricht außer `p` und `Terminate` → `PGR-E6001`, selbst gefahren) bleibt eine tragfähige, benannte Lücke. Kein Blocker für die Zusage zu Verfahren, Ablehnung und Channel Binding, aber vor der Closure zu klären.
- **LOW (V-2):** Die Aussage zu Kerberos, GSSAPI und SSPI ist nur als Reaktion des Werkzeugs auf Anforderungen belegt, die ein Python-Ersatz sendet (Codes 2, 6, 7, 9, 42); ein echter Server ist nicht gefahren. Das Handbuch sagt das genau so (V-2 verlangt keine Änderung, nur die Kenntnisnahme, dass der Satz „Ersatz-Server“ tragen muss).
- **INFO (V-3):** Die Befehlsfolge ist für diesen Slice nicht „vorher rot“ (siehe Tabelle); die Vorgabe der Aufgabe trifft nicht zu, und der Plan behauptet es nicht. Nur der Mutant ist rot.
- **INFO (V-4):** Nicht neu gefahren: voller `make gates` (Hook approve, `gates-schnell` und `docs-check` grün), Windows, macOS, andere Treiber (pgx, JDBC, `psycopg`), andere libpq- und PostgreSQL-Versionen, ein echter Kerberos-Server, Herunterfahren während der Anmeldung, gleichzeitige Verbindungen in der Anmeldung, die Zusage „Startup-Parameter unverändert durchgereicht“ mit einem zusätzlichen Parameter. Das Handbuch behauptet keines davon.
- **INFO (V-5):** Offen für die Closure wie in §2: DoD-Häkchen, Closure-Notiz mit Lerneintrag, Register, Ausgänge der drei Risiken, Paarungen.

## Fazit

Beide Liefer-Punkte sind bestätigt bis auf V-1: Die zwei Beispielblöcke laufen als Datei gegen das gebaute Binary; jede Verhaltensaussage zu Verfahren (Klartext, MD5, SCRAM, `trust`), Ablehnung, Abbruch durch den Client, Channel Binding, `PGR-E4002`, `E4003`, `E6001` für Verfahren und `offset_ms` stimmt am Binary; Versionen (PostgreSQL 17.11, libpq 17.11) sind belegt und genannt; das Handbuch enthält keine Verweise auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und nichts Künftiges; das README nennt den Ist-Zustand mit Verweisen nach oben; die Versionszeile steht auf 0.8; die Größe ist 105 gegen Marke 300; `make gates-schnell`, `make docs-check` (745 Dateien, 0 Befunde) und `make kopf-check` sind auf HEAD grün, der Hook gibt frei. Die Review-Findings F-1, F-2, F-3 sind behoben, F-4 bis F-6 sind INFO ohne Handbuch-Verstoß. F-7 trägt als Lücke auf der Clientseite, nicht auf der Serverseite, wo der Satz zu `PGR-E4002` für eine nicht lesbare Nachricht falsch ist (V-1). Freigabe zur Closure unter der Bedingung, dass V-1 vorher behoben oder vom Planner ausdrücklich als akzeptiertes Negativ geführt wird; danach DoD-Häkchen, Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge und Paarungen (V-5).
