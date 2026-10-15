# pgwire-recorder

> **Zeichnet die Kommunikation zwischen Ihrer Anwendung und PostgreSQL auf, spielt sie ohne Datenbank wieder ab oder führt ihre Anfragen erneut gegen eine Datenbank aus.**

## Was ist pgwire-recorder?

`pgwire-recorder` ist ein Kommandozeilenwerkzeug, das zwischen einer Anwendung
und PostgreSQL vermittelt. Im Modus `record` leitet es die Kommunikation an eine
echte Datenbank weiter und speichert sie in einer Datei, im Modus `replay`
beantwortet es dieselben Anfragen aus dieser Datei, ohne Datenbank, und im Modus
`play` führt es die aufgezeichneten Anfragen gegen eine Datenbank aus. Es richtet
sich an Entwickler, Tester und CI-Systeme. Die Anforderungen stehen im
[Lastenheft](spec/lastenheft.md).

## Was kann ich heute tun?

`pgwire-recorder record` vermittelt einfache Anfragen und das erweiterte Protokoll
(vorbereitete Anweisungen, etwa von pgx im Standardmodus) zwischen einem Client und
PostgreSQL und schreibt eine Aufzeichnung, als YAML-Datei oder mit `--format sqlite` als
SQLite-Datei; `pgwire-recorder replay` beantwortet einfache wie vorbereitete Anfragen
daraus ohne Datenbank und erkennt das Format der Datei selbst, ebenso `play`. Lebendprüfungen von Verbindungspools (`pgxpool`,
`database/sql`) beantwortet `replay` unabhängig davon, ob sie in der Aufzeichnung
stehen; mit `--fail-on-unconsumed` endet `replay` mit einem Fehler, wenn
aufgezeichnete Interaktionen nicht abgerufen wurden. `pgwire-recorder play` spielt
die Anfragen einer Aufzeichnung, einfache wie vorbereitete Anweisungen, gegen eine
Datenbank ein; ohne `--compare-responses` bricht es nach einer Fehlerantwort der
Datenbank ab, mit `--continue-on-error` läuft es weiter, und mit
`--finish-session-on-interrupt` endet es nach einem Abbruchsignal erst nach der laufenden
Sitzung. Mit `--compare-responses` vergleicht `play` die Antworten der Datenbank mit der
Aufzeichnung (Art und Reihenfolge der Antworten, Spalten, Befehl, SQLSTATE eines Fehlers,
Transaktionsstatus; nicht Zeilenwerte und Hinweise) und meldet eine Abweichung als
`PGR-E5004` mit Exit-Code 5. Ohne `--compare-responses` endet `play`, wenn die Datenbank
mit einem COPY-Datenstrom antwortet, mit `PGR-E6001`. `record` und `replay` nehmen
Verbindungen der Anwendung ohne TLS an; mit `--tls-cert` und `--tls-key` (PEM-Dateien) nehmen
sie TLS an (mindestens TLS 1.2, ohne ein Zertifikat der Anwendung zu verlangen) und weisen eine
Anwendung ohne TLS mit `PGR-E6003` ab, es sei denn, `--allow-plaintext` lässt beides zu;
`record` verbindet ohne TLS zur Datenbank; `play`
verbindet auf Wunsch (`--upstream-tls` oder `sslmode=require` der Verbindung) mit TLS,
prüft das Zertifikat des Servers, auch gegen eine eigene Zertifizierungsstelle
(`--upstream-ca`), und meldet die Ablehnung von TLS und ein nicht passendes Zertifikat als `PGR-E4005`.
`record` vermittelt die Anmeldung der Anwendung (Klartext-Passwort, MD5, SCRAM-SHA-256)
und zeichnet das Passwort nicht auf (SCRAM-SHA-256 geprüft mit PostgreSQL 14 bis 18,
Klartext-Passwort und MD5 mit PostgreSQL 17); `play` meldet sich mit
Klartext-Passwort, MD5 und SCRAM-SHA-256 an, mit dem Passwort aus dem Platzhalter der
Verbindung oder aus `PGWIRE_RECORDER_PASSWORD`.
Mit `--record-timing` trägt `record` an jeder Anfrage die Millisekunden seit der
Startup-Nachricht der Verbindung ein (`offset_ms`); `play --keep-timing` wartet vor den
Anfragen auf diese Zeit (beim Standard `relative` ohne Wartezeit vor der ersten Anfrage
einer Verbindung), gemessen ab der vorigen Anfrage (`relative`, Standard) oder ab einem Bezugspunkt
(`--timing-mode absolute` mit `--timing-reference`), und `replay` wartet nie. Trägt eine Anfrage kein
`offset_ms`, endet `play --keep-timing` mit `PGR-E2003`; ein ungültiges oder absteigendes
`offset_ms` lehnen `replay` und `play` mit `PGR-E3003` ab.

`record` zeichnet eine Verbindung, die eine Anfrage abgeschlossen hat, als eigene Sitzung auf,
auch mehrere gleichzeitig; antwortet die Datenbank ihr mit einem COPY-Datenstrom oder sendet sie
ihr eine Benachrichtigung (`LISTEN` mit `NOTIFY`), zeichnet `record` sie nicht auf (`PGR-E6001`).
Mit `--record-empty-sessions` wird auch eine Verbindung ohne Anfrage eine Sitzung ohne
Interaktionen, soweit sie keine der Ausnahmen ist, die das Benutzerhandbuch nennt. `replay`
ordnet beim Standard `--session-assignment first-request` die n-te Verbindung mit einer
Anfrage der n-ten Sitzung mit Interaktionen zu; mit `--session-assignment connection` erhält
die n-te Verbindung die Sitzung mit der `id` n, was eine mit `--record-empty-sessions`
aufgezeichnete Datei verlangt (sonst `PGR-E2001`). Eine Anfrage ohne zugeordnete Sitzung ist
`PGR-E5003`, eine Anfrage an eine Sitzung ohne Interaktionen `PGR-E5001`.

Am Rand des Protokolls gilt in `record` und `replay`: Ein Parameter `_pq_.` in der
Startnachricht und eine Nachricht `p` nach dem Aufbau sind `PGR-E6001`, eine andere Protokollversion
als 3.0 ist `PGR-E6002`. Eine Verbindung, deren Startnachricht das Werkzeug ablehnt, ist keine
Sitzung und verbraucht im `replay` keine. `replay` beginnt jede angenommene Verbindung mit `AuthenticationOk`,
ohne ein Passwort zu verlangen, und prüft Benutzer und Datenbank nicht; ein Client, der mit `require_auth` ein
Anmeldeverfahren verlangt, bricht deshalb selbst ab.

Beim Beenden warten `record` und `replay` höchstens `--shutdown-timeout`
(Standard 5 Sekunden) auf laufende Anfragen; was das im Container bedeutet, sagt das
Benutzerhandbuch.

- `make build` baut das Image `pgwire-recorder:dev` mit dem Binary; `make gates`
  führt die Gates aus, `make help` zeigt alle Targets.
- Das [Benutzerhandbuch](docs/user/benutzerhandbuch.md) beschreibt Bau, Aufruf und
  Verhalten des Binaries.

## Warum pgwire-recorder?

Automatisierte Tests von Anwendungen mit PostgreSQL-Anbindung brauchen häufig
eine echte oder eigens bereitgestellte Datenbank. Das erhöht den Aufwand für
lokale Umgebungen und CI-Systeme und kann Tests langsamer, komplexer und weniger
deterministisch machen. Für viele dieser Tests ist nicht die Datenbank selbst der
Gegenstand, sondern das aus Sicht der Anwendung beobachtbare Verhalten der
PostgreSQL-Kommunikation. Genau dieses Verhalten zeichnet `pgwire-recorder` auf.
Dieselbe Aufzeichnung lässt sich außerdem in eine Datenbank einspielen, um die
Wirkung der aufgezeichneten Anfragen dort herzustellen, zum Beispiel um eine
Komponente zu testen, die Änderungen der Datenbank verarbeitet (Change Data
Capture).

## Kerngedanke

Im Modus `replay` gibt der Recorder nur wieder, was er aufgezeichnet hat, strikt
in der aufgezeichneten Reihenfolge. Eine Anfrage, die nicht zur Aufzeichnung
passt, ist ein Fehler und bekommt nie eine geratene Antwort. Im Modus `play`
führt er die aufgezeichneten Anfragen nacheinander aus und wertet Fehlerantworten
des Servers; mit `--compare-responses` vergleicht er die Antworten mit der Aufzeichnung.

## Was macht es vertrauenswürdig?

- **Prozess:** [`AGENTS.md`](AGENTS.md) (Hard Rules: die Regeln, die jede Änderung einhalten muss), [`harness/README.md`](harness/README.md) (Source Precedence: welche Quelle bei Konflikt gewinnt, und die Gates).
- **Verträge:** [`spec/lastenheft.md`](spec/lastenheft.md) (`LH-*`-IDs mit Akzeptanzkriterien), danach [`spec/spezifikation.md`](spec/spezifikation.md) und [`spec/architecture.md`](spec/architecture.md).
- **Gates:** `make gates` führt die Gates aus; welche es sind und was jedes prüft, steht in [`harness/README.md` §Sensors](harness/README.md#sensors-feedback-gates).
- **Auditierbarkeit:** Entscheidungen in [`docs/plan/adr/`](docs/plan/adr/), Planung in [`docs/plan/planning/`](docs/plan/planning/), Reviews in [`docs/reviews/`](docs/reviews/).
- **Rückverfolgbarkeit:** `make doc-trace` gibt die Requirements Traceability Matrix aus (Anforderung, Entscheidungen, Slices); sie ist ein Bericht und kein Gate.
