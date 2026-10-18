# Review slice-v1-abschluss-tls-client-doku

Rolle: Reviewer. Datum: 2026-10-14.
Gegenstand: Diff `59e14db..77efd06` (Move-Commit bis HEAD, drei Commits) gegen den Plan [`slice-v1-abschluss-tls-client-doku`](../plan/planning/done/welle-v1-abschluss/slice-v1-abschluss-tls-client-doku.md), den Code-Slice [`slice-v1-abschluss-tls-client`](../plan/planning/done/welle-v1-abschluss/slice-v1-abschluss-tls-client.md), [LH-FA-23](../../spec/lastenheft.md#lh-fa-23--verschlüsselung-zum-client), `LH-FA-23.a`, `LH-FA-17.a` und `AGENTS.md` §3.7, §3.9, §3.11, §3.13, §3.14.

Der Report wurde vom Reviewer als Text geliefert, weil dessen Schreibzugriff abgelehnt wurde; der Orchestrator hat ihn unverändert abgelegt.

**Summary: 0 HIGH · 0 MEDIUM · 1 LOW · 4 INFO.**

## Findings

### F-1 LOW
- `quelle`: `AGENTS.md` §3.9 (Plan folgt dem Diff)
- `pfad`: `docs/plan/planning/in-progress/slice-v1-abschluss-tls-client-doku.md` · "Der Slice ersetzt im Handbuch die Grenz-Sätze zu TLS zum Client"
- `befund`: Der Diff ändert im Handbuch die Sitzungs-Sätze (*Mehrere Verbindungen aufzeichnen und wiedergeben*) und den Satz zu `config show`. Die Sitzungs-Sätze sagen: `PGR-W3002` und `PGR-E6003` zählen bei `record` und bei `replay` mit `first-request` nicht als Sitzung. Weder §1 noch §3 noch die DoD nennen diese Stellen; sie stehen nur in §7. Der Inhalt ist belegt (Probe P10, `TestConfigShowLiestTLSNicht`), der Plan trägt ihn nicht.
- `verifizierbar`: nein
- `klasse`: Plan nennt eine Handbuch-Stelle nicht, die der Diff ändert

### F-2 INFO
- `quelle`: `AGENTS.md` §3.9
- `pfad`: Plan · "Der Wortlaut der DoD ist veraltet."
- `befund`: DoD-Zeile 2 nennt den README-Satz „Alle Verbindungen laufen unverschlüsselt“. Er steht im README nicht mehr. Der Diff ersetzt stattdessen „`record` und `replay` nehmen Verbindungen nur ohne TLS an …“. Handbuch und README folgen dem Binary, nur die DoD ist veraltet. Der Planner entscheidet über den Wortlaut. Kein Befund gegen den Diff.
- `klasse`: DoD-Wortlaut hinter dem Bestand

### F-3 INFO
- `quelle`: `AGENTS.md` §3.14
- `pfad`: Plan · "ohne den Plan 231, mit dem Plan 337"
- `befund`: Die Messung stimmt für `59e14db..dbe6232`: Handbuch 196+30, README 4+1, Plan 105+1. Am letzten Commit `77efd06` stehen die Pläne bei 112+1 und die Summe bei 344. §7 nennt den Messbereich ausdrücklich, die Zahl ist also wahr. Sie gilt aber nicht für den letzten Commit des Slice, den §3.14 meint. Keine Aktion nötig, denn §4 hat keine Zeilenmarke.
- `klasse`: Größenzahl nicht am letzten Commit

### F-4 INFO
- `quelle`: `AGENTS.md` §3.11
- `pfad`: `docs/user/benutzerhandbuch.md` · "bei libpq `sslnegotiation=direct`"
- `befund`: Die Aussagen zu direktem TLS und zur Kette („Mit libpq gelingt `sslmode=verify-full` …“) sind nur mit libpq 17 belegt (Implementer P4, P6; Review mit `psql` 17). Das Handbuch nennt die Version nicht. Windows, macOS, pgx und JDBC erscheinen im Handbuch nicht als geprüft. „Der Aufruf lief mit OpenSSL 3.3“ nennt nur das Geprüfte.
- `klasse`: Belegte Umgebung nicht im Satz genannt

### F-5 INFO
- `quelle`: Maintainability
- `pfad`: `docs/user/benutzerhandbuch.md` · "Der erste Block ist das Zertifikat des Servers, die weiteren sind seine Kette;"
- `befund`: Zwei Absätze in *Hinweise* sind ohne Umbruch zu langen Zeilen geworden (zweite Stelle: „Die Namen, die Verwendung und den Aussteller …“). Es gibt keinen Konventions-Anker, also keine Aktion.
- `klasse`: Zeilenumbruch uneinheitlich

## Eigene Prüfung (nicht auf Zusage)

**Beispiele als Datei gegen das gebaute Binary** (`make build`, Image `pgwire-recorder:dev`, Code seit `59e14db` unverändert; gepinntes `postgres:17-alpine`, eigenes Docker-Netz, ein Lauf zur Zeit, Container und Netz danach entfernt, kein `pkill`):
- Der `openssl`-Aufruf aus dem Handbuch lief (OpenSSL 3.0.13). `server.pem` beginnt mit `BEGIN CERTIFICATE`, `server.key` mit `BEGIN PRIVATE KEY`.
- Das `record`-Beispiel mit beiden Dateien und `psql sslmode=require` liefert die Zeile `1`. Die Aufzeichnung enthält `select 1 as eins`.
- `psql sslmode=verify-full sslrootcert=./server.pem` mit `host=localhost` liefert `1`. Mit `host=recorder` scheitert es am Namen im Zertifikat, wie das Handbuch es sagt.
- `replay --config` mit `tls_cert` und `tls_key` im Abschnitt `replay:` liefert `1` unter `verify-full`.

**Meldungen gegen das Binary:**
- `--tls-cert` allein, `--allow-plaintext` ohne Dateien und `--allow-plaintext=true` ohne Dateien: `PGR-E2001`, Exit 2. `--allow-plaintext=false` ohne Dateien ist gültig.
- Schlüsseldatei als Zertifikat und umgekehrt: `PGR-E2007`, Option und Ursache, ohne Pfad.
- `play --tls-cert`: `PGR-E2001`. `PGWIRE_RECORDER_TLS_CERT` bei `play`: `PGR-E2001` mit dem Variablennamen. `play.tls_cert` in der Konfigurationsdatei: `PGR-E2004`.
- `sslmode=disable` gegen `record` und `replay`: `PGR-E6003`. Bei `record` ist der Exit-Code beim Beenden 6, die Verbindung wird nicht aufgezeichnet.
- `replay` mit `--session-assignment connection`: Die abgewiesene Klartext-Verbindung verbraucht die Sitzung nicht, die folgende TLS-Verbindung erhält Sitzung 1.
- `sslnegotiation=direct`: Log `PGR-W3003`, Client „EOF“, auch bei `replay`.
- `verify-full` ohne Wurzel: Log `PGR-W3002`, Exit 0 beim Beenden.

**Fehlercode-Tabelle:** Alle Codes des Katalogs `internal/hexagon/model/fehler.go` stehen in der Tabelle. `PGR-E6003` und `PGR-W3002` sind neu. `PGR-E2001`, `PGR-E2004`, `PGR-E2007`, `PGR-E6001` und `PGR-W3003` sind um die neuen Fälle ergänzt. Jede Ergänzung ist durch Test oder Probe gedeckt.

**Tests vorhanden** (Namen gegen die Quelle gesucht und gefunden): `TestTLSZweiterAufbau`, `TestTLSEndeOhneAlert`, `TestTLSDateiKeineRegulaere`, `TestTLSZertifikatUngeprueft`, `TestConfigShowLiestTLSNicht`, `TestTLSKetteWirdAusgeliefert`, `TestTLSRangfolge`, `TestRunTLSLogOhneSchluessel`, `TestTLSDirekt`, `TestE2ETLSStartfehler`, `TestTLSHilfe`. Nicht nachgefahren: Ablauf nach dem Start (P7), Zeit nach dem Handshake (P5), TLS 1.1, `close_notify` und Schlüssel im Log; dort wurden Tests und Proben des Implementers nur gelesen.

**Quantoren:**
- „jede allein ist ungültig“, „jede Datei … regulär“, „in keiner Log-Zeile“, „kein Zertifikat der Anwendung“ und „ruhende Verbindung hält das Beenden nicht auf“ haben je einen Test.
- Die Bestandssätze „Das Werkzeug fragt nie nach Eingaben“ und „nimmt jede Anmeldung an“ bleiben wahr. Ein verschlüsselter Schlüssel endet mit Exit 2 ohne Rückfrage, und es wird kein Client-Zertifikat angefordert.
- „Eine Verbindung zählt erst ab ihrer Startnachricht“ bleibt wahr.

**Ist-Zustand und Verweise:**
- `grep -n -i` im Handbuch auf Slice, Welle, ADR-, LH-, SPEC-, ARC-, `spec/`, Lastenheft, Spezifikation, Review, „noch nicht“, „geplant“, „künftig“, „bisher“: drei Treffer, alle Bestand und Zustandsbeschreibung.
- Das README hat keine Chronik, die Verweise nach oben sind unverändert.
- Windows, macOS und JDBC kommen im Handbuch nicht vor. pgx steht nur im Bestand zur Lebendprüfung.

**Plan:**
- §1, §3, §6 und §8 sind im Diff unverändert, §7 ist neu. Die Zählung in §8 (zwei Liefer-Punkte, eine Schicht) passt zum Diff.
- Die Übernahme aus `slice-v1-abschluss-umgebung-anderes-kommando` ist geliefert. Der Code-Slice liegt in `done/`.
- Der Move-Commit `59e14db` ist rein, der Inhalt steht in eigenen Commits (§3.3). Alle drei Commits nennen die Slice-Kennung.
- Die DoD-Häkchen sind noch offen, wie vor der Closure.

## Geprüft, ohne Befund
- `docs/user/benutzerhandbuch.md`: außer F-1, F-4, F-5 ohne Befund. Version 0.7 und Stand 14.10.2026 stimmen mit dem Diff.
- `README.md`: ohne Befund (ein Satz, Ist-Zustand, kein Zielstand).
- `docs/plan/planning/in-progress/slice-v1-abschluss-tls-client-doku.md`: außer F-1, F-2, F-3 ohne Befund.
- `spec/`, `internal/`, `cmd/`, `test/`, `harness/`, `.claude/`: im Diff unverändert, ohne Befund.
