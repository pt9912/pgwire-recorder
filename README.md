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
PostgreSQL und schreibt eine YAML-Aufzeichnung; `pgwire-recorder replay` beantwortet
beide daraus ohne Datenbank. Lebendprüfungen von Verbindungspools (`pgxpool`,
`database/sql`) beantwortet `replay` unabhängig davon, ob sie in der Aufzeichnung
stehen; mit `--fail-on-unconsumed` endet `replay` mit einem Fehler, wenn
aufgezeichnete Interaktionen nicht abgerufen wurden. `pgwire-recorder play` spielt
die Anfragen einer Aufzeichnung, einfache wie vorbereitete Anweisungen, gegen eine
Datenbank ein; nach einer Fehlerantwort der Datenbank bricht es ab, mit
`--continue-on-error` läuft es weiter, und mit `--finish-session-on-interrupt` endet
es nach einem Abbruchsignal erst nach der laufenden Sitzung. Antwortet die Datenbank
mit einem COPY-Datenstrom, endet `play` mit `PGR-E6001`. Alle Verbindungen laufen
unverschlüsselt. `record` verlangt, dass die Datenbank den Benutzer ohne Passwort
anmeldet; `play` meldet sich mit Klartext-Passwort, MD5 und SCRAM-SHA-256 an, mit dem
Passwort aus dem Platzhalter der Verbindung oder aus `PGWIRE_RECORDER_PASSWORD`. Beim Beenden warten `record` und `replay` höchstens `--shutdown-timeout`
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
des Servers.

## Was macht es vertrauenswürdig?

- **Prozess:** [`AGENTS.md`](AGENTS.md) (Hard Rules: die Regeln, die jede Änderung einhalten muss), [`harness/README.md`](harness/README.md) (Source Precedence: welche Quelle bei Konflikt gewinnt, und die Gates).
- **Verträge:** [`spec/lastenheft.md`](spec/lastenheft.md) (`LH-*`-IDs mit Akzeptanzkriterien), danach [`spec/spezifikation.md`](spec/spezifikation.md) und [`spec/architecture.md`](spec/architecture.md).
- **Gates:** `make gates` führt die Gates aus; welche es sind und was jedes prüft, steht in [`harness/README.md` §Sensors](harness/README.md#sensors-feedback-gates).
- **Auditierbarkeit:** Entscheidungen in [`docs/plan/adr/`](docs/plan/adr/), Planung in [`docs/plan/planning/`](docs/plan/planning/), Reviews in [`docs/reviews/`](docs/reviews/).
- **Rückverfolgbarkeit:** `make doc-trace` gibt die Requirements Traceability Matrix aus (Anforderung, Entscheidungen, Slices); sie ist ein Bericht und kein Gate.
