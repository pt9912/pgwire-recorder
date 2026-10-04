# pgwire-recorder

> **Zeichnet die Kommunikation zwischen Ihrer Anwendung und PostgreSQL auf und spielt sie später ohne Datenbank wieder ab.**

## Was ist pgwire-recorder?

`pgwire-recorder` ist ein Kommandozeilenwerkzeug, das zwischen einer Anwendung
und PostgreSQL vermittelt. Im Modus `record` leitet es die Kommunikation an eine
echte Datenbank weiter und speichert sie in einer Datei, im Modus `replay`
beantwortet es dieselben Anfragen aus dieser Datei, ohne Datenbank, und im Modus
`play` führt es die aufgezeichneten Anfragen gegen eine Datenbank aus und
vergleicht auf Wunsch deren Antworten mit der Aufzeichnung. `record` und `replay`
nehmen auf Wunsch verschlüsselte Verbindungen (TLS) von Anwendungen an. Es richtet
sich an Entwickler, Tester und CI-Systeme. Die vollständige Beschreibung steht
im [Lastenheft](spec/lastenheft.md).

## Was kann ich heute tun?

Das Projekt steht am Beginn der Umsetzung. Es gibt ein Go-Gerüst mit der
Package-Struktur des Hexagons, aber noch keine Funktion.

- `make gates` läuft grün (Build und Test des Go-Moduls in Docker,
  Architekturregeln samt Gegenprobe, Doku-Referenzen und vendored Baseline);
  `make help` zeigt alle Targets.
- Lastenheft, Spezifikation, Architektur, die angenommenen Entscheidungen und die
  Planung (Roadmap, Wellen, Slices) liegen vor.
- Das [Benutzerhandbuch](docs/user/benutzerhandbuch.md) beschreibt das Verhalten
  der ersten Version.

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
führt er die aufgezeichneten Anfragen nacheinander aus, wertet Fehlerantworten
des Servers und vergleicht auf Wunsch die Struktur der Antworten nach festgelegten
Regeln.

## Was macht es vertrauenswürdig?

- **Prozess:** [`AGENTS.md`](AGENTS.md) (Hard Rules: die Regeln, die jede Änderung einhalten muss), [`harness/README.md`](harness/README.md) (Source Precedence: welche Quelle bei Konflikt gewinnt, und die Gates).
- **Verträge:** [`spec/lastenheft.md`](spec/lastenheft.md) (`LH-*`-IDs mit Akzeptanzkriterien), danach [`spec/spezifikation.md`](spec/spezifikation.md) und [`spec/architecture.md`](spec/architecture.md).
- **Gates:** `make gates` führt `make build` und `make test` (Go-Modul in Docker), `make a-check` (Architekturregeln) mit der Gegenprobe `make a-check-negativ`, `make docs-check` (Doku-Referenzen) und `make baseline-verify` (Integrität der vendored Baseline) aus.
- **Auditierbarkeit:** Entscheidungen in [`docs/plan/adr/`](docs/plan/adr/), Planung in [`docs/plan/planning/`](docs/plan/planning/), Reviews in [`docs/reviews/`](docs/reviews/).
- **Rückverfolgbarkeit:** `make doc-trace` gibt die Requirements Traceability Matrix aus (Anforderung, Entscheidungen, Slices); sie ist ein Bericht und kein Gate.
