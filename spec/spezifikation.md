# Spezifikation — pgwire-recorder

**Bezug zum Lastenheft:** Diese Spezifikation präzisiert die in
`spec/lastenheft.md` formulierten Anforderungen (`LH-*`-IDs). Bei
Konflikt gewinnt das Lastenheft — präzisieren ja, erweitern nie.

**Rolle:** Technik-Stratum — fortschreibbar ohne Change Request; eine ADR darf
sie schärfen, das Lastenheft nicht. Regeln: Baseline-Regelwerk
`modul-03-spec.md` §Ziel-Form: Spezifikation.

**Zielversion des Produkts:** v1

**Zweck:** `pgwire-recorder` wird als Kommandozeilenwerkzeug implementiert, das
PostgreSQL-PGWire-Kommunikation im Record-Modus zu einem realen PostgreSQL-Server
weiterleitet und die relevanten Interaktionen persistiert. Im Replay-Modus tritt
das Werkzeug selbst als PostgreSQL-Endpunkt auf und reproduziert aufgezeichnete
Antworten ohne Upstream-Datenbank.

---

## 1. Algorithmen und Datenflüsse

Regeln dieser Sektion: Tatsächlicher Code gehört in `src/`, nicht hierher.
ID-Schema `LH-FA-<NN>.<Buchstabe>` für Verfeinerungen einzelner
Lastenheft-IDs (Baseline-Regelwerk `grundlagen-source-precedence.md`
§ID-Schema als Klammer). Was **keine** einzelne Lastenheft-ID verfeinert,
trägt eine `SPEC-<NNN>` — siehe §2 bis §7.

### LH-FA-01.a — Kommandos und Hilfe

**Eingabe:** Aufruf des Programms `pgwire-recorder` mit Kommando und Optionen.
**Ausgabe:** Ausführung des Kommandos; Hilfetext bei `--help` beziehungsweise
`-h`.

v1 stellt mindestens folgende Kommandos bereit:

```text
pgwire-recorder record
pgwire-recorder replay
pgwire-recorder play
pgwire-recorder config show
pgwire-recorder version
```

Globale Hilfetexte sind über `--help` beziehungsweise `-h` verfügbar.

**Fehlermodi:** ungültige CLI-Verwendung → Exit-Code `2` (`SPEC-014`).

---

### LH-FA-02.a — Record: Verbindungsannahme, Upstream, Weiterleitung

**Eingabe:** Client-Verbindungen auf der `--listen`-Adresse. **Ausgabe:**
weitergeleitete Kommunikation und eine Aufzeichnung.

Verpflichtende Optionen von `record`:

| Option | Bedeutung |
|---|---|
| `--listen` | Adresse, auf der der Recorder PostgreSQL-Clients annimmt |
| `--upstream` | Adresse des realen PostgreSQL-Servers |
| `--output` | Zieldatei des Recordings |

Optional: `--force` (Zieldatei überschreiben, Default aus; siehe LH-FA-07.a),
`--format` (`yaml` oder `sqlite`, siehe LH-FA-22.a), `--record-timing` (siehe
LH-FA-21.a) und `--record-empty-sessions` (siehe LH-FA-12.a).

Beispiel:

```bash
pgwire-recorder record \
  --listen 0.0.0.0:15432 \
  --upstream postgres:5432 \
  --output ./recordings/users.yaml
```

**Schritte:**

1. Der Recorder öffnet den mit `--listen` angegebenen TCP-Endpunkt.
2. Für jede akzeptierte Client-Verbindung wird eine Recorder-Session angelegt.
3. Für jede Recorder-Session wird eine zugehörige Verbindung zum konfigurierten
   PostgreSQL-Upstream hergestellt. Die Kommunikation darf nicht über eine
   gemeinsam genutzte PostgreSQL-Verbindung verschiedener Client-Sessions
   gemultiplext werden.
4. Unterstützte Client-Nachrichten werden an den Upstream weitergeleitet,
   Antworten des Upstreams an den Client. Der Recorder verändert fachliche
   Inhalte nicht.

**Fehlermodi:** Listen-Port nicht zu öffnen → Startfehler, Exit-Code `4`
(`PGR-E4001`); Upstream nicht erreichbar → Verbindungsfehler, `PGR-E4002`
(Fehlerebenen: LH-FA-13.b).

---

### LH-FA-02.b — Abschluss einer Interaktion

Eine Simple-Query-Interaktion gilt als abgeschlossen, sobald das zugehörige
`ReadyForQuery` des Servers verarbeitet wurde.

**Verbindungsende.** Endet die Client-Verbindung nach einem `ReadyForQuery`,
mit oder ohne `Terminate`, ist das regulär; die Session wird mit ihren
abgeschlossenen Interaktionen übernommen. Bricht die Client- oder die
Upstream-Verbindung vor dem `ReadyForQuery` einer laufenden Interaktion ab, ist
das ein unerwartetes Verbindungsende (Verbindungsfehler, `PGR-E4003`); die
unvollständige Interaktion wird nicht in das Recording übernommen, die
vorherigen Interaktionen der Session bleiben erhalten.

---

### LH-FA-03.a — Replay: Aufruf, Upstream-Freiheit, Session-Cursor

**Eingabe:** Client-Verbindungen auf der `--listen`-Adresse und ein Recording.
**Ausgabe:** aufgezeichnete Backend-Nachrichten.

Verpflichtende Optionen von `replay`:

| Option | Bedeutung |
|---|---|
| `--listen` | Adresse, auf der der Replay-Server Clients annimmt |
| `--input` | Zu verwendendes Recording |

Optional: `--fail-on-unconsumed` (nicht verbrauchte Interaktionen als Fehler
werten, Default aus; siehe LH-FA-03.b) und `--session-assignment` (siehe
LH-FA-12.a).

Beispiel:

```bash
pgwire-recorder replay \
  --listen 0.0.0.0:15432 \
  --input ./recordings/users.yaml
```

**Schritte:**

1. Im Replay-Modus wird keine Verbindung zu einem realen PostgreSQL-Server
   benötigt.
2. Beim Start wird das Recording geladen und geprüft. Ein nicht lesbares oder
   beschädigtes Recording oder eine unbekannte Version endet als Startfehler
   mit Exit-Code `3`; ein Recording ohne Session ist nicht verwendbar (Exit-Code
   `3`, `PGR-E3004`).
3. Jede eingehende Client-Verbindung erhält einen eigenen Replay-Cursor und, je
   nach `--session-assignment`, eine Session (LH-FA-12.a).
4. Der Cursor zeigt auf die nächste erwartete Interaktion der zugeordneten
   Session; er beginnt am Anfang der Session. Stellt eine Verbindung eine
   Anfrage, zu der es keine nicht zugeordnete Session mehr gibt, ist das ein
   Replay-Mismatch (`PGR-E5003`).

---

### LH-FA-03.b — Nicht verbrauchte Interaktionen

Wird eine Replay-Session beendet, bevor alle ihr zugeordneten Interaktionen
verbraucht wurden, oder wird eine aufgezeichnete Session nie verbunden, wird dies mindestens als Warnung mit dem Meldungscode `PGR-W2001` protokolliert (`SPEC-034`).

Die Option `--fail-on-unconsumed` (Umgebungsvariable
`PGWIRE_RECORDER_FAIL_ON_UNCONSUMED`) wertet dies als Fehler: Die Verbindung
zählt als fehlerhaft beendet (`PGR-E5002`, Exit-Code `5`, Fehlerebenen:
LH-FA-13.b). Der Default ist aus (`SPEC-012`); Query-Mismatches sind unabhängig
davon immer Fehler.

---

### LH-FA-05.a — Unterstützter Umfang: Startup, Simple Query und Extended Query

Startup, Simple Query Protocol, Extended Query Protocol (LH-FA-18.a) und
Serverantworten bilden den unterstützten Umfang von v1.

**Startup.** Der Recorder unterstützt die für typische PostgreSQL-Clients
notwendige Startup-Sequenz, soweit sie für die definierten Record-/Replay-
Szenarien erforderlich ist. Startup-Parameter werden als Teil einer Session
erfasst, soweit sie für das Replay relevant sind.

**Client-Nachrichten.** Für das Simple Query Protocol sind relevant:

- `Query`
- `Terminate`

Die `Query`-Nachricht wird vollständig einschließlich SQL-Text aufgezeichnet.

**Serverantworten.** Alle Servernachrichten, die während einer unterstützten
Simple-Query-Interaktion bis zum zugehörigen `ReadyForQuery` auftreten und von
der verwendeten PGWire-Bibliothek verlustfrei repräsentiert werden können,
werden in ihrer Reihenfolge aufgezeichnet. Dazu gehören insbesondere:

- `RowDescription`
- `DataRow`
- `CommandComplete`
- `EmptyQueryResponse`
- `ErrorResponse`
- `NoticeResponse`
- `ParameterStatus`
- `ReadyForQuery`

Die konkrete Menge ist nicht auf diese Liste beschränkt. Eine Serverantwort,
die die PGWire-Bibliothek nicht verlustfrei repräsentiert, ist eine nicht
unterstützte Interaktion (`PGR-E6001`).

**Nicht unterstützte Protokollnachrichten.** Trifft v1 auf eine
Client-Interaktion, die als nicht unterstützt klassifiziert ist, schlägt das Tool
kontrolliert fehl beziehungsweise beendet die betroffene Verbindung mit einer
verständlichen Diagnose (Exit-Code `6`, `SPEC-019`). Es wird keine scheinbar
erfolgreiche Aufzeichnung erzeugt: Eine Session, in der eine nicht unterstützte
Interaktion auftrat, wird nicht in das Recording übernommen.

**Extended Query Protocol.** `Parse`, `Bind`, `Describe`, `Execute`, `Close`,
`Flush`, `Sync` und Prepared Statements sind unterstützt; Ablauf, Aufzeichnung
und Matching legt LH-FA-18.a fest.

---

### LH-FA-05.b — Authentifizierung

**Record.** Der Recorder vermittelt die Startup-/Authentifizierungskommunikation
transparent zwischen Client und Upstream, soweit dies mit dem gewählten
Betriebsmodell möglich ist. Credentials sind keine eigene Recorder-Konfiguration,
wenn der Client sie selbst im PGWire-Handshake liefert.

**Replay.** Replay stellt einen für Testclients verwendbaren
Startup-/Authentifizierungsablauf bereit. Die konkret emulierten
Authentifizierungsnachrichten werden in der Implementierung so gewählt, dass
typische PostgreSQL-Treiber ohne reale Datenbankverbindung eine Session aufbauen
können. v1 ist nicht als Sicherheitsgrenze gedacht; Replay läuft in einer
kontrollierten Testumgebung.

---

### LH-FA-05.c — TLS

TLS-Passthrough ohne Protokolleinsicht reicht für Recording nicht aus. TLS-
Terminierung beziehungsweise TLS-Unterstützung ist daher **kein Muss für v1**.

Versucht ein Client eine SSL/TLS-Aushandlung (`SSLRequest`), lehnt der Recorder
sie mit dem Einzelbyte `N` ab und erwartet einen unverschlüsselten
PGWire-Verbindungsaufbau auf derselben Verbindung. Besteht der Client auf TLS
und bricht ab, ist das kein Fehler des Recorders. TLS-Unterstützung ist eine
spätere Erweiterung.

---

### LH-FA-05.d — Transaktionen

SQL-Transaktionsbefehle, die über Simple Query oder Extended Query übertragen
werden, werden wie andere Anfragen aufgezeichnet und strikt sequenziell wiedergegeben.

Der Recorder implementiert im Replay-Modus keine eigene Transaktionslogik. Der
beobachtbare Transaktionsstatus wird durch die aufgezeichneten Serverantworten,
insbesondere `ReadyForQuery`, reproduziert.

---

### LH-FA-05.e — Protokollversion und Protokollrand

**Protokollversion.** v1 unterstützt PGWire-Protokollversion 3.0. Enthält die
`StartupMessage` eine andere Protokollversion, beantwortet der Recorder sie mit
einer `ErrorResponse` und beendet die Verbindung (`PGR-E6002`, Klasse
„nicht unterstützt", Exit-Code `6`). Im Record-Modus ist jeder PostgreSQL-Server
zulässig, der PGWire 3.0 spricht (`SPEC-029`).

**Protokollrand.** Folgende Interaktionen sind in v1 eindeutig behandelt:

| Interaktion | Verhalten |
|---|---|
| `SSLRequest` | mit `N` abgelehnt (LH-FA-05.c) |
| `GSSENCRequest` | mit `N` abgelehnt |
| `CancelRequest` | Verbindung wird geschlossen, nicht an den Upstream weitergeleitet; kein Verbindungsfehler, Warnung `PGR-W3001` |
| Extended-Query-Nachrichten (`Parse`, `Bind`, `Describe`, `Execute`, `Close`, `Flush`, `Sync`) | unterstützt (LH-FA-18.a) |
| `COPY`-Nachrichten | nicht unterstützt (`PGR-E6001`) |
| `FunctionCall` | nicht unterstützt (`PGR-E6001`) |
| `NotificationResponse` (asynchron, z. B. nach `LISTEN`) | nicht unterstützt (`PGR-E6001`) |
| SQL-Befehle `LISTEN`/`NOTIFY` als `Query` | wie jede andere Query aufgezeichnet |

Bei einer nicht unterstützten Interaktion wird die betroffene Verbindung mit
einer `ErrorResponse` beendet; Fehlerebenen: LH-FA-13.b.

---

### LH-FA-06.a — Inhalt der Aufzeichnung

Für jede unterstützte Interaktion (einfache Query oder Extended-Interaktion) wird
mindestens gespeichert:

- Session-Zuordnung,
- fortlaufende Interaktionsnummer,
- Anfrageart (einfach oder Extended),
- SQL-Text beziehungsweise die Client-Nachrichten der Extended-Interaktion,
- geordnete Folge der zugehörigen Serverantworten,
- für das Replay notwendige PGWire-Felder.

Zeitangaben enthält die Aufzeichnung nur auf Wunsch (`--record-timing`,
LH-FA-21.a); sie sind keine Voraussetzung für ein erfolgreiches Replay.

---

### LH-FA-07.a — Sicheres Schreiben des Recordings

Dieser Abschnitt beschreibt das Format `yaml`; für `sqlite` gilt `LH-FA-22.a`.

**Eingabe:** Sessions eines Record-Laufs. **Ausgabe:** Datei unter dem
`--output`-Pfad.

**Zielpfad beim Start.** Existiert der `--output`-Pfad bereits und ist `--force`
nicht gesetzt, endet `record` als Startfehler mit Exit-Code `2` (`PGR-E2002`),
bevor eine Verbindung angenommen wird. Mit `--force` ersetzt bei `yaml` der erste
Schreibvorgang die vorhandene Datei.

**Schreibzeitpunkt.** Das Recording wird nach dem Ende jeder Session und beim
kontrollierten Beenden (LH-FA-13.a) als Ganzes neu geschrieben; es enthält alle
bis dahin beendeten Sessions. Eine Verbindung ohne Anfrage (zum Beispiel eine
Probe-Verbindung eines Connection-Pools) wird nicht aufgezeichnet, außer mit
`--record-empty-sessions` (LH-FA-12.a). Ein Lauf ohne aufgezeichnete Session
schreibt ein gültiges Recording ohne Sessions.

**Schritte je Schreibvorgang:**

1. Die Implementierung schreibt zunächst in eine temporäre Datei im Verzeichnis
   der Zieldatei.
2. Sie verschiebt die temporäre Datei atomar beziehungsweise bestmöglich atomar
   auf die Zieldatei.

Die Zieldatei ist damit zu jedem Zeitpunkt entweder nicht vorhanden oder ein
vollständiges Recording; sie ist nie syntaktisch unvollständig. Bei einem
Abbruch darf eine temporäre Datei verbleiben; sie wird nicht stillschweigend als
gültiges Recording behandelt.

**Fehlermodi:** vorhandenes `--output` ohne `--force` → Exit-Code `2`
(`PGR-E2002`); Recording nicht lesbar/schreibbar, unbekannte Version oder
beschädigt → Exit-Code `3` (`SPEC-016`).

---

### LH-FA-08.a — Auswahl des Recordings

Das Recording wird über `--output` (Record) beziehungsweise `--input` (Replay)
bestimmt (siehe LH-FA-02.a und LH-FA-03.a).

---

### LH-FA-09.a — Strict sequential matching

v1 verwendet als einziges Matching-Verfahren **strict sequential matching**
(`SPEC-011`): Anfragen müssen in der erwarteten Reihenfolge zur Aufzeichnung
passen.

**Eingabe:** eingehende `Query` und der Replay-Cursor der Session.
**Ausgabe:** die aufgezeichneten Backend-Nachrichten oder ein Mismatch.

**Schritte:**

1. Der eingehende SQL-Text wird mit dem aufgezeichneten SQL-Text der
   Interaktion am Cursor verglichen: `eingehender SQL-Text == aufgezeichneter
   SQL-Text`.
2. Bei Gleichheit werden die für die Interaktion aufgezeichneten
   Backend-Nachrichten in der gespeicherten Reihenfolge an den Client gesendet
   und der Cursor rückt vor.
3. Bei Ungleichheit liegt ein Mismatch vor (LH-FA-10.a).

Der Vergleich erfolgt byte- beziehungsweise stringgenau nach der
PGWire-Dekodierung. v1 normalisiert SQL nicht. Insbesondere werden nicht
automatisch ignoriert: Whitespace-Unterschiede, Kommentare,
Groß-/Kleinschreibung, Literalwerte und semantisch äquivalente SQL-Varianten.

Mehrfach identische Queries werden über die Position am Cursor zugeordnet; jede
aufgezeichnete Interaktion wird höchstens einmal verbraucht.

---

### LH-FA-10.a — Mismatch

Bei einem Mismatch enthält die Diagnose mindestens:

- Session beziehungsweise Verbindungskontext,
- erwartete Interaktionsnummer,
- erwartete Query,
- tatsächlich empfangene Query.

Der Recorder springt nicht zur nächsten Interaktion und führt keine Fuzzy-Suche
durch. Sendet ein Client weitere Anfragen, obwohl keine aufgezeichnete
Interaktion mehr verfügbar ist, wird dies als Replay-Mismatch behandelt.

**Fehlermodi:** Replay-Mismatch → Exit-Code `5` (`SPEC-018`).

---

### LH-FA-11.a — Aufzeichnung von Fehlerantworten

`ErrorResponse` wird als Teil der geordneten Serverantworten der Interaktion
aufgezeichnet (siehe LH-FA-05.a) und im Replay in gespeicherter Reihenfolge,
einschließlich des abschließenden `ReadyForQuery`, wiedergegeben.

---

### LH-FA-12.a — Reihenfolge, Sessions und Parallelität

**Interaktionsreihenfolge.** `sequence` und die geordnete Response-Liste im
Recording (`SPEC-002`, `SPEC-041`) erhalten die Reihenfolge der Interaktionen.

**Recording.** Mehrere Client-Verbindungen dürfen parallel angenommen werden. Jede
Verbindung mit mindestens einer abgeschlossenen Interaktion wird als eigene Session
geführt; mit `--record-empty-sessions` jede angenommene Verbindung, auch eine ohne
Anfrage (als Session ohne Interaktionen). Die `id` ist eine fortlaufende Zahl ab 1
ohne Lücken und wird beim Schreiben der Session vergeben, in der Reihenfolge der
Sessionenden; eine verworfene Session (LH-FA-05.a) erhält keine `id`. Werden die
Verbindungen nacheinander genutzt, ist das die Reihenfolge der ersten Anfragen. Mit
`--record-empty-sessions` trägt das Recording die Kennzeichnung `empty_sessions:
true` (YAML: Feld auf der obersten Ebene neben `format` und `version`; SQLite: Zeile
in `meta`, `SPEC-043`).

**Replay.** `--session-assignment` wählt die Zuordnung der Sessions zu den
Verbindungen in der Reihenfolge der `id`:

* `first-request` (Default): Die n-te Verbindung mit einer Anfrage, in der
  Reihenfolge der ersten Anfragen, erhält die n-te Session mit mindestens einer
  Interaktion. Sessions ohne Interaktion werden übersprungen; Verbindungen ohne
  Anfrage zählen nicht.
* `connection`: Die n-te angenommene Verbindung erhält die Session mit der `id` n,
  auch wenn sie keine Anfrage stellt. Das setzt ein Recording mit
  `empty_sessions: true` voraus; fehlt die Kennzeichnung, ist das ein Startfehler
  (`PGR-E2001`).

Der Handshake nutzt bei `connection` die Startup-Daten der zugeordneten Session; bei
`first-request` vor der ersten Anfrage die der nächsten noch nicht zugeordneten Session
mit Interaktion, gibt es keine mehr, die der letzten. Die Zuordnung ist deterministisch,
wenn die Clients ihre Verbindungen nacheinander aufbauen und nutzen; bei gleichzeitiger
Nutzung hängt sie von der Reihenfolge der Verbindungen beziehungsweise der ersten
Anfragen ab und ist nicht zugesichert. Eine Anfrage einer Verbindung, zu der es keine
nicht zugeordnete Session mehr gibt, ist ein Replay-Mismatch (`PGR-E5003`). Eine
Session mit Interaktionen, die nie zugeordnet wird, gilt als nicht verbraucht
(LH-FA-03.b). Das Einspielen (`play`) führt Sessions ohne Interaktion nicht aus.

---

### LH-FA-13.a — Signalbehandlung

Auf `SIGINT` und `SIGTERM` fährt der Prozess kontrolliert herunter. Unter
Windows entspricht der Konsolenabbruch (`Strg+C`, `Strg+Break`) einem `SIGINT`:

- keine neuen Verbindungen annehmen,
- laufende Schreiboperationen soweit möglich abschließen,
- Sessions, die noch laufen, nach Abschluss ihrer laufenden Interaktion
  beenden; eine nicht abgeschlossene Interaktion wird nicht übernommen
  (LH-FA-02.b),
- das Recording schreiben (LH-FA-07.a); schlägt das Schreiben fehl, bleibt die
  Zieldatei der letzte vollständig geschriebene Stand.

Der Exit-Code nach einem kontrollierten Herunterfahren folgt LH-FA-13.b.

---

### LH-FA-13.b — Fehlerebenen und Prozessstatus

Diese Fehlerebenen gelten für `record` und `replay`; `play` nimmt keine
Verbindungen an und endet am Ende seines Laufs (`LH-FA-20.a`).

**Startfehler** (ungültige Konfiguration, Listen-Port nicht zu öffnen,
Recording nicht ladbar, vorhandenes `--output` ohne `--force`) beenden den Prozess sofort mit dem Exit-Code
der Klasse (`SPEC-013` bis `SPEC-019`), bevor eine Verbindung angenommen wird.

**Verbindungsfehler** (Upstream nicht erreichbar, Replay-Mismatch, nicht
unterstützte Interaktion, unerwartetes Verbindungsende, bei
`--fail-on-unconsumed` nicht verbrauchte Interaktionen, Anfrage ohne
nicht zugeordnete Session) beenden nur die betroffene Verbindung. Dem Client wird, wo das Protokoll es erlaubt, eine
`ErrorResponse` mit dem Meldungscode im Meldungstext zugestellt. Der Prozess
läuft weiter und merkt sich die Klasse des ersten aufgetretenen
Verbindungsfehlers.

**Prozessende.** Sessions, die bis zum kontrollierten Herunterfahren nie
zugeordnet wurden, zählen bei `--fail-on-unconsumed` als Fehler der Klasse 5
(`PGR-E5002`), sonst als Warnung (`PGR-W2001`). Nach einem kontrollierten
Herunterfahren ist der Exit-Code der
der gemerkten Klasse, sonst `0`. Schlägt dabei das Schreiben des Recordings
fehl, ist er `3` und hat Vorrang vor der gemerkten Klasse. Beendet der Prozess nicht kontrolliert (zum Beispiel durch
`SIGKILL`), gilt kein Exit-Code der Spezifikation.

---

### LH-FA-14.a — Logging und Diagnose

Logs werden nach `stderr` geschrieben (`SPEC-006`); Nutzdaten beziehungsweise
maschinenlesbare Ausgaben auf `stdout` werden dadurch nicht verunreinigt.
Unterstützte Log-Level sind mindestens `error`, `warn`, `info`, `debug`;
Standard ist `info` (`SPEC-005`). Die Detailstufe wird über `--log-level`
(Umgebungsvariable `PGWIRE_RECORDER_LOG_LEVEL`) gesetzt.

Passwörter aus Verbindungsdaten werden nicht absichtlich in Logs ausgegeben.

Fehler und Warnungen tragen einen Meldungscode und erzeugen einen für Entwickler verständlichen Text (`SPEC-034`).

---

### LH-FA-15.a — Nicht-interaktiver Betrieb

Alle Kommandos sind vollständig nicht-interaktiv; sie fragen keine Eingaben ab
(siehe LH-FA-17.a).

---

### LH-FA-16.a — Containerbetrieb

Das Projekt stellt ein Docker/OCI-konformes Container-Image bereit
(`SPEC-031`): ein Linux-Image für `linux/amd64` und `linux/arm64` in einer
einzigen Manifestliste, das das Binary ohne weitere Laufzeitabhängigkeit
enthält. Das Binary hat keine Laufzeitabhängigkeit auf eine lokale
PostgreSQL-Installation. Das Image läuft mit jedem OCI-kompatiblen
Container-Laufzeitsystem. Es wird in der GitHub Container Registry
(`ghcr.io/pt9912/pgwire-recorder`) und auf Docker Hub
(`docker.io/pt9912/pgwire-recorder`) veröffentlicht.

Beispiel:

```yaml
services:
  recorder:
    image: pgwire-recorder:latest
    command:
      - replay
      - --listen=0.0.0.0:5432
      - --input=/recordings/test.yaml
    volumes:
      - ./recordings:/recordings:ro
```

---

### LH-FA-17.a — Konfiguration

Alle für CI erforderlichen Einstellungen sind über CLI-Argumente verfügbar, mit zwei
Ausnahmen: das Passwort und die benannten Verbindungen. Umgebungsvariablen und eine
Konfigurationsdatei dürfen ergänzend genutzt werden. Wird dieselbe Einstellung
mehrfach angegeben, gilt (`SPEC-007`):

```text
CLI-Argument > Umgebungsvariable > Konfigurationsdatei > Standardwert
```

Der Name der Umgebungsvariablen einer Option ist das Präfix `PGWIRE_RECORDER_`
(`SPEC-008`), gefolgt vom Optionsnamen in Großbuchstaben mit `_` statt `-`; boolesche
Werte lauten `true` oder `false`. Das **Passwort** hat eine eigene Regel: Es kommt
aus dem Passwort der benutzten benannten Verbindung (nur als Platzhalter `${VAR}`),
sonst aus `PGWIRE_RECORDER_PASSWORD`; eine Option dafür gibt es nicht.

**Konfigurationsdatei.** Es gilt genau **eine** Datei, in dieser Reihenfolge:

1. die Datei aus `--config`,
2. der Pfad aus der Umgebungsvariable `PGWIRE_RECORDER_CONFIG`,
3. die Datei `.pgwire-recorder.yaml` im aktuellen Verzeichnis, sofern sie existiert.

Findet sich keine Datei, wird keine gelesen. Eine mit `--config` oder
`PGWIRE_RECORDER_CONFIG` genannte Datei, die fehlt, ist ein Konfigurationsfehler. Die
Schlüssel heißen wie die Optionen, mit `_` statt `-` und ohne die führenden `--`
(`keep_timing` für `--keep-timing`). `log_level` und die benannten Verbindungen stehen
auf der obersten Ebene, die übrigen Schlüssel in einem Abschnitt je Kommando
(`record:`, `replay:`, `play:`):

```yaml
log_level: info
connections:
  lokal: "postgresql://dev@localhost:5432/myapp"
  staging: "postgresql://app:${STAGING_PASSWORD}@staging.example.com:5432/myapp?sslmode=require"
play:
  upstream: staging
  input: ./recordings/users.yaml
  keep_timing: true
  timing_mode: relative
```

**Benannte Verbindungen.** Der Wert unter `connections:` ist eine URL der Form
`postgresql://[benutzer[:passwort]@]host[:port]/datenbank[?parameter]`. Der
einzige Parameter ist `sslmode` mit den Werten `disable` (Default) und `require`;
`require` verbindet mit TLS und prüft das Serverzertifikat gegen den
Zertifikatsspeicher des Systems (strenger als bei libpq); jeder andere Parameter
oder Wert ist ein Konfigurationsfehler. `--upstream` und der Schlüssel `upstream`
nehmen den Namen einer Verbindung oder `host:port`; ein Name hat Vorrang vor
`host:port`. Die Wirkung einer URL:

* `play`: Host und Port, Benutzer und Datenbank. Benutzer und Datenbank aus
  `--user` und `--database` gehen vor denen der URL, diese vor den Startup-Daten der
  Session. TLS gilt, wenn `--upstream-tls` (Option, Umgebungsvariable oder
  Schlüssel) gesetzt ist; ist es nicht gesetzt, entscheidet `sslmode`. Ein
  ausdrücklich gesetztes `--upstream-tls` geht dem `sslmode` vor.
* `record`: nur Host und Port. Benutzer, Passwort und Datenbank der URL werden
  ignoriert, weil `record` die Anmeldung des Clients vermittelt; `sslmode=require`
  ist ein Konfigurationsfehler, weil `record` kein TLS zum Upstream kennt
  (`LH-FA-05.c`).

**Geheimnisse.** Der Platzhalter `${VAR}` ist nur in der URL einer benannten
Verbindung erlaubt und wird für die benutzte Verbindung aus der gleichnamigen
Umgebungsvariable ersetzt (`$${VAR}` bleibt wörtlich); die Variablen nicht benutzter
Verbindungen bleiben unbeachtet. Ein **Klartext-Passwort** ist ein Passwortteil
hinter dem `:` im Benutzerteil einer URL, der nicht genau ein `${VAR}` ist, oder ein
Parameter `password`.

**Fehler.** Jede Ursache trägt einen eigenen Code und nennt in der Meldung die
Stelle (Schlüssel oder Verbindungsname), nie einen Wert:

| Ursache | Code |
|---|---|
| Datei nicht lesbar oder nicht vorhanden, ungültiges YAML, unbekannter Schlüssel, Schlüssel im falschen Abschnitt, ungültiger Wert, ungültiger `sslmode`, Platzhalter außerhalb einer URL | `PGR-E2004` |
| nicht gesetzte Umgebungsvariable eines Platzhalters der benutzten Verbindung | `PGR-E2005` |
| Klartext-Passwort in der Datei | `PGR-E2006` |

Alle drei sind Startfehler mit Exit-Code `2`. Ein ungültiger Wert einer Option oder
Umgebungsvariable (auch ein Wert außerhalb einer Aufzählung) ist `PGR-E2001`.

**Anzeige.** `pgwire-recorder config show` gibt den Inhalt der gewählten Datei als
eingerückten Baum auf `stdout` aus und nennt die Datei; Exit-Code `0`. Findet sich
keine Datei, meldet der Befehl das und endet mit Exit-Code `0`. Ist die Datei
ungültig, endet er mit dem Fehlercode des Ladens (`PGR-E2004` bis `PGR-E2006`) und
zeigt nichts. Platzhalter erscheinen unaufgelöst; die Ausgabe enthält nie einen
aufgelösten Wert, aber Hosts, Benutzer und Pfade der Datei. Aktive
`PGWIRE_RECORDER_*`-Umgebungsvariablen listet der Befehl am Ende mit Namen und ohne
Werte. Er verbindet sich nicht und liest keine Aufzeichnung.

| Option | Kommando | Umgebungsvariable | Default |
|---|---|---|---|
| `--listen` | `record`, `replay` | `PGWIRE_RECORDER_LISTEN` | — (Pflicht) |
| `--upstream` | `record`, `play` | `PGWIRE_RECORDER_UPSTREAM` | — (Pflicht; `host:port` oder Name einer Verbindung) |
| `--output` | `record` | `PGWIRE_RECORDER_OUTPUT` | — (Pflicht) |
| `--force` | `record` | `PGWIRE_RECORDER_FORCE` | `false` |
| `--format` | `record` | `PGWIRE_RECORDER_FORMAT` | `yaml` (Werte `yaml`, `sqlite`) |
| `--record-timing` | `record` | `PGWIRE_RECORDER_RECORD_TIMING` | `false` |
| `--record-empty-sessions` | `record` | `PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS` | `false` |
| `--input` | `replay`, `play` | `PGWIRE_RECORDER_INPUT` | — (Pflicht) |
| `--fail-on-unconsumed` | `replay` | `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` | `false` (`SPEC-012`) |
| `--session-assignment` | `replay` | `PGWIRE_RECORDER_SESSION_ASSIGNMENT` | `first-request` (Werte `first-request`, `connection`) |
| `--user` | `play` | `PGWIRE_RECORDER_USER` | Startup-Daten der Session |
| `--database` | `play` | `PGWIRE_RECORDER_DATABASE` | Startup-Daten der Session |
| `--continue-on-error` | `play` | `PGWIRE_RECORDER_CONTINUE_ON_ERROR` | `false` |
| `--allow-recorded-errors` | `play` | `PGWIRE_RECORDER_ALLOW_RECORDED_ERRORS` | `false` |
| `--upstream-tls` | `play` | `PGWIRE_RECORDER_UPSTREAM_TLS` | `false` |
| `--finish-session-on-interrupt` | `play` | `PGWIRE_RECORDER_FINISH_SESSION_ON_INTERRUPT` | `false` |
| `--keep-timing` | `play` | `PGWIRE_RECORDER_KEEP_TIMING` | `false` |
| `--timing-mode` | `play` | `PGWIRE_RECORDER_TIMING_MODE` | `relative` (Werte `relative`, `absolute`) |
| `--timing-reference` | `play` | `PGWIRE_RECORDER_TIMING_REFERENCE` | `connect` (Werte `connect`, `first-request`) |
| — (nur Umgebung) | `play` | `PGWIRE_RECORDER_PASSWORD` | — |
| `--config` | `record`, `replay`, `play`, `config show` | `PGWIRE_RECORDER_CONFIG` | `.pgwire-recorder.yaml` im aktuellen Verzeichnis, sofern vorhanden |
| `--log-level` | `record`, `replay`, `play` | `PGWIRE_RECORDER_LOG_LEVEL` | `info` (`SPEC-005`) |

---

### LH-FA-18.a — Extended Query: Ablauf, Aufzeichnung, Matching

**Nachrichten.** Client: `Parse`, `Bind`, `Describe`, `Execute`, `Close`,
`Flush`, `Sync`. Server: `ParseComplete`, `BindComplete`, `CloseComplete`,
`ParameterDescription`, `RowDescription`, `NoData`, `DataRow`,
`CommandComplete`, `EmptyQueryResponse`, `PortalSuspended`, `ErrorResponse`,
`NoticeResponse`, `ReadyForQuery`.

**Interaktion.** Eine Extended-Interaktion beginnt mit der ersten
Extended-Nachricht nach dem Sessionbeginn oder nach dem `Sync` der vorherigen
Interaktion und endet mit dem `ReadyForQuery`, das auf ihr `Sync` folgt.
Client-Nachrichten, die nach einem `Sync` eintreffen, bevor dessen
`ReadyForQuery` eingetroffen ist, gehören zur nächsten Interaktion. Eine Session
mischt Simple- und Extended-Interaktionen in der aufgezeichneten Reihenfolge.

**Abbruch.** Endet die Verbindung, bevor das `ReadyForQuery` einer Interaktion
verarbeitet wurde, auch nach einem `Terminate` oder nach einem `Flush` ohne
`Sync`, ist das ein unerwartetes Verbindungsende (`PGR-E4003`, LH-FA-02.b): die
unvollständige Interaktion wird nicht übernommen, die vorherigen Interaktionen
der Session bleiben erhalten.

**Record.** Der Recorder leitet alle Nachrichten unverändert und in
Ankunftsreihenfolge weiter, auch wenn der Client mehrere Nachrichten sendet,
ohne auf Antworten zu warten (Pipelining). Zeichnet er die Interaktion auf,
gruppiert er sie (`SPEC-041`); Zeitangaben enthält die Aufzeichnung nur mit
`--record-timing` (LH-FA-21.a).

**Gruppen.** Eine Extended-Interaktion besteht aus einer oder mehreren Gruppen.
Eine Gruppe endet mit einem `Sync` oder einem `Flush`. Innerhalb einer Gruppe
stehen in der Aufzeichnung zuerst alle Client-Nachrichten in Sendereihenfolge,
danach alle Server-Nachrichten, die auf die Gruppe antworten, in
Empfangsreihenfolge. Welche Server-Nachrichten zu welcher Gruppe gehören:

* Endet die Gruppe mit `Sync`, sind es alle Server-Nachrichten bis
  einschließlich des zugehörigen `ReadyForQuery`; jedes `Sync` erzeugt genau
  ein `ReadyForQuery`. Die Zuordnung hängt nicht vom Zeitverhalten ab, auch
  nicht, wenn der Client weitere Gruppen sendet, ohne zu warten.
* Endet die Gruppe mit `Flush`, sind es die Server-Nachrichten, die vor der
  nächsten Client-Nachricht eintreffen. Wartet der Client nach dem `Flush` nicht
  auf die Antwort, ist die Zuordnung zeitabhängig; das ist eine bekannte Grenze
  von v1. Für diesen Fall ist die Aufzeichnung nicht reproduzierbar, die
  Wiedergabe einer vorhandenen Aufzeichnung bleibt deterministisch.

Die Aufzeichnung derselben Anwendung ist damit für alle Interaktionen mit
`Sync`-Gruppen und für `Flush`-Gruppen mit wartendem Client unabhängig vom
Zeitverhalten gleich.

**Replay.** Der Replay-Cursor zeigt auf die nächste erwartete Gruppe. Jede
eingehende Client-Nachricht muss in allen Feldern der erwarteten Client-Nachricht
der Gruppe entsprechen: Nachrichtentyp, Statement- und Portalname, SQL-Text,
Parametertypen, Format-Codes für Parameter und Ergebnis, Parameterwerte
(bytegenau), `max_rows`, bei
`Describe` und `Close` auch die Zielart (Statement oder Portal). Nachdem der
Recorder die letzte Client-Nachricht der Gruppe (`Flush` oder `Sync`) empfangen
und verglichen hat, sendet er alle aufgezeichneten Server-Nachrichten dieser
Gruppe. Vorher sendet er keine; die Freigabe hängt damit nur von der
Aufzeichnung ab, nicht vom Zeitverhalten des Clients. Das
Matching ist strict sequential wie in LH-FA-09.a; Namen werden nicht
normalisiert. Ein Client, der nichtdeterministische Statement-Namen erzeugt,
passt deshalb nicht zur Aufzeichnung (Mismatch, `PGR-E5001`); das ist eine
bekannte Grenze von v1.

**Fehler.** Nach einer `ErrorResponse` verwirft ein Server Nachrichten bis zum
nächsten `Sync`. Die Aufzeichnung enthält die empfangenen Client-Nachrichten
und die Server-Nachrichten der Gruppe; Replay reproduziert sie, ohne eigene
Fehlerlogik.

**Mismatch.** Die Diagnose nach LH-FA-10.a nennt zusätzlich den Index des
Gruppe und Nachricht in der Interaktion sowie erwarteten und empfangenen Nachrichtentyp;
Parameterwerte erscheinen nicht im Klartext der Diagnose, wenn der Log-Level
nicht `debug` ist (`SPEC-033`).

---

### LH-FA-19.a — Bereitstellung über Homebrew

Die Bereitstellung erfolgt über einen eigenen Tap `pt9912/homebrew-pgwire-recorder`
(Tap `pt9912/pgwire-recorder`, Formel `pgwire-recorder`). Die Formel installiert
das veröffentlichte Binary des Release-Tags für macOS und Linux, jeweils
`amd64` und `arm64`, und prüft es gegen die SHA-256-Summe desselben Releases;
sie baut nicht aus dem Quelltext und hat keine Abhängigkeiten. Die Formel
entsteht bei jedem stabilen Release. Eine Vorabversion ändert den Tap nicht.
Windows wird über Homebrew nicht bedient (`SPEC-042`).

---

### LH-FA-20.a — Einspielen

Das Kommando `play` führt die Client-Anfragen einer Aufzeichnung gegen einen
PostgreSQL-Server aus. Es nimmt keine Client-Verbindungen an; sein Prozess endet am
Ende des Laufs und folgt nicht den Fehlerebenen von `LH-FA-13.b`.

Verpflichtende Optionen: `--upstream` (Adresse des Servers oder Name einer
benannten Verbindung, `LH-FA-17.a`) und `--input` (Aufzeichnung). Optional (Defaults
in `LH-FA-17.a`):

| Option | Wirkung |
|---|---|
| `--user`, `--database` | überschreiben Benutzer und Datenbank der Verbindung und die Startup-Daten der jeweiligen Session |
| `--continue-on-error` | ein Serverfehler einer Anfrage beendet das Einspielen nicht sofort |
| `--allow-recorded-errors` | ein Serverfehler in einer Interaktion, deren Aufzeichnung ebenfalls eine `ErrorResponse` enthält, gilt als erwartet |
| `--upstream-tls` | die Verbindung zum Server wird mit TLS aufgebaut |
| `--finish-session-on-interrupt` | ein Abbruchsignal beendet zuvor die laufende Session |
| `--keep-timing`, `--timing-mode`, `--timing-reference` | zeitgetreues Einspielen (`LH-FA-21.a`) |

Beispiel:

```bash
pgwire-recorder play \
  --upstream postgres:5432 \
  --input ./recordings/users.yaml
```

**Schritte:**

1. Die Aufzeichnung wird geladen und geprüft (Fehler wie bei `replay`, Exit-Code
   `3`).
2. Für jede Session mit Interaktionen, in der Reihenfolge ihrer `id`, baut der
   Recorder eine eigene Verbindung zum Server auf. TLS gilt wie in `LH-FA-17.a`
   beschrieben; mit TLS wird das Serverzertifikat gegen den Zertifikatsspeicher des
   Systems geprüft, ein Überspringen der Prüfung gibt es nicht. Er authentifiziert
   sich als Client mit Klartext-Passwort, MD5 oder SCRAM-SHA-256.
3. Die Client-Nachrichten der Interaktionen werden in der aufgezeichneten
   Reihenfolge gesendet: bei einer einfachen Anfrage die `Query`, bei einer
   Extended-Interaktion die Nachrichten jeder Gruppe (`SPEC-041`). Nach jeder
   Interaktion wartet der Recorder auf das `ReadyForQuery`.
4. Die Serverantworten werden gelesen und verworfen, mit Ausnahme einer
   `ErrorResponse`. Ein Vergleich mit den aufgezeichneten Antworten findet nicht
   statt.
5. **Fehler beim Aufbau** (vor dem ersten `ReadyForQuery`): Ein nicht erreichbarer
   Server und jeder `FATAL`, der nicht die Anmeldung betrifft (zum Beispiel eine
   fehlende Datenbank), ist `PGR-E4002`. Eine fehlgeschlagene Anmeldung
   (SQLSTATE-Klasse 28), ein nicht unterstütztes Verfahren, ein abgelehntes TLS
   (der Server antwortet auf die TLS-Anfrage mit `N`), ein ungültiges Zertifikat und
   ein Server, der unverschlüsselte Verbindungen ablehnt, sind `PGR-E4005`.
   **Fehler danach:** Bricht die Verbindung nach dem ersten `ReadyForQuery` ab, auch
   durch einen `FATAL`, ist das `PGR-E4003`. Alle drei beenden das Einspielen immer,
   auch mit `--continue-on-error`.
6. **Fehlerantwort einer Anfrage:** Antwortet der Server mit einer `ErrorResponse`,
   endet das Einspielen mit Exit-Code `4` (`PGR-E4004`). Mit `--continue-on-error`
   läuft es mit der nächsten Interaktion weiter und endet am Ende mit Exit-Code
   `4`. Mit `--allow-recorded-errors` zählt eine `ErrorResponse` nicht, wenn die
   aufgezeichnete Interaktion mindestens eine `ErrorResponse` enthält; bei einer
   Extended-Interaktion verwirft der Server danach bis zum `Sync`, und der Recorder
   sendet weiter wie aufgezeichnet. Bei jedem Abbruch schließt der Recorder die
   Verbindung, soweit möglich, mit `Terminate`.
7. **Abbruchsignal:** Bei `SIGINT` oder `SIGTERM` endet das Einspielen nach der
   laufenden Interaktion; eine Wartezeit (`LH-FA-21.a`) wird abgebrochen. Mit
   `--finish-session-on-interrupt` zuvor nach der laufenden Session. Ein zweites
   Signal beendet den Prozess sofort. Die Verbindung wird mit `Terminate`
   geschlossen. Der Exit-Code ist `0`, wenn bis dahin kein Fehler auftrat, sonst
   `4`.

Das Einspielen ist sequenziell: die Anfragen einer Session nacheinander und die
Sessions nacheinander. Ohne `--keep-timing` gibt es keine Wartezeiten; mit
`--keep-timing` gilt `LH-FA-21.a`.

---

### LH-FA-21.a — Zeitangaben

**Aufzeichnen.** Mit `--record-timing` (Default `false`) trägt der Recorder an
jeder Interaktion das Feld `offset_ms` ein: eine ganze Zahl größer oder gleich
null, die Millisekunden zwischen der Annahme der Verbindung und dem Eintreffen der
ersten Client-Nachricht der Interaktion, gemessen mit einer monotonen Uhr. Das
Feld steht auf der Ebene der Interaktion (`SPEC-002`, `SPEC-041`, `SPEC-043`) und
fehlt ohne die Option. Eine Aufzeichnung mit Zeitangaben ist nicht reproduzierbar
im Sinne von `SPEC-004`: zwei Läufe derselben Anwendung unterscheiden sich in
`offset_ms`. Ein Wert, der keine ganze Zahl ≥ 0 ist, oder Werte, die innerhalb einer
Session absteigen, machen das Recording zu einem beschädigten (`PGR-E3003`).

**Einspielen.** Mit `--keep-timing` (Default `false`) stellt `play` den zeitlichen
Abstand her; `--timing-mode` wählt:

* `relative` (Default): vor jeder Interaktion außer der ersten einer Session wartet
  `play`, bis seit dem Beginn der vorigen Interaktion mindestens der Unterschied
  ihrer beiden `offset_ms` vergangen ist. Der Abstand zwischen zwei aufeinander
  folgenden Interaktionen ist damit nie kürzer als aufgezeichnet; ist der Server
  langsamer, wird sofort gesendet, und die Verspätung wird nicht aufgeholt.
* `absolute`: vor jeder Interaktion wartet `play`, bis seit dem Bezugspunkt
  mindestens der aufgezeichnete Abstand zum aufgezeichneten Bezugspunkt vergangen
  ist. Keine Interaktion liegt damit früher als aufgezeichnet, gemessen am
  Bezugspunkt; eine Verspätung wird aufgeholt, der Abstand zwischen zwei
  aufeinander folgenden Interaktionen kann dabei kürzer sein als aufgezeichnet.
  `--timing-reference` wählt den Bezugspunkt: `connect` (Default) ist der Beginn des
  Verbindungsaufbaus zum Server, vor TLS und Anmeldung; aufgezeichnet gilt
  `offset_ms` = 0. Die Dauer von Aufbau und Anmeldung zählt als Verspätung.
  `first-request` ist der Beginn der ersten Interaktion der Session; aufgezeichnet
  gilt der `offset_ms` der ersten Interaktion.

Jede Session beginnt mit eigener Uhr, die Sessions laufen weiterhin nacheinander.
Ungültige Verwendung sind (`PGR-E2001`, Exit-Code `2`, Startfehler): `--timing-mode`
oder `--timing-reference` ohne `--keep-timing`, `--timing-reference` zusammen mit
`relative`, ein Wert außerhalb der Aufzählungen. Fehlt irgendeiner Interaktion der
Aufzeichnung `offset_ms`, ist die Verwendung ebenfalls ungültig (`PGR-E2003`,
Exit-Code `2`, Startfehler). Ohne `--keep-timing` werden vorhandene Zeitangaben
ignoriert; der Replay-Modus ignoriert sie immer.

---

### LH-FA-22.a — Aufzeichnungsformat

**Wahl und Erkennung.** `record` wählt das Format mit `--format` (Umgebungsvariable
`PGWIRE_RECORDER_FORMAT`, Werte `yaml` und `sqlite`, Default `yaml`). Die
Dateiendung hat keine Bedeutung; üblich sind `.yaml` und `.sqlite` (`SPEC-009`).
`replay` und `play` erkennen das Format der Eingabedatei selbst: eine Datei, die mit
der SQLite-Kopfzeile `SQLite format 3` beginnt, ist eine SQLite-Aufzeichnung, jede
andere wird als YAML gelesen. Eine Datei, die in ihrem Format keine gültige
Aufzeichnung ist (falsche Formatkennung, fehlende Tabellen, unlesbarer oder
widersprüchlicher Inhalt), ist beschädigt (`PGR-E3003`); eine unbekannte `version`
ist `PGR-E3002`.

**Lesen von `sqlite`.** Die Datei wird nur lesend geöffnet, auch auf einem
schreibgeschützten Mount. Ist sie gesperrt oder nicht lesbar, ist das `PGR-E3001`.

**Schreiben von `sqlite`.** Beim Start prüft `record` den Zielpfad wie in
`LH-FA-07.a`. Danach legt es in einer temporären Datei im Verzeichnis der Zieldatei
eine gültige Aufzeichnung ohne Sessions an (Schema aus `SPEC-043`, `meta` gefüllt) und
verschiebt sie atomar auf den Zielpfad; mit `--force` ersetzt dieser Schritt die
vorhandene Datei. Jede beendete Session wird danach in einer Transaktion ergänzt; die
Datei enthält zu jedem Zeitpunkt nur vollständige Sessions und wird nie als Ganzes
neu geschrieben. Der Journalmodus ist das Rollback-Journal (`journal_mode=DELETE`),
es entstehen keine WAL-Nebendateien, und Fremdschlüssel sind eingeschaltet
(`foreign_keys=ON`).

**Beide Formate** tragen dasselbe logische Modell (`SPEC-002`, `SPEC-041`,
`SPEC-043`); eine Aufzeichnung in einem Format und dieselbe Aufzeichnung im anderen
liefern im Replay und beim Einspielen dasselbe Verhalten. Die Zusagen von `SPEC-001`
und `SPEC-004` (textbasiert, diff-freundlich, deterministisch) gelten für `yaml`; eine
SQLite-Datei ist binär, und zwei Läufe derselben Anwendung liefern logisch, aber nicht
byte-gleiche Dateien.

---

## 2. Datenstrukturen und Schemas

Regeln dieser Sektion: Jede Struktur trägt eine `SPEC-<NNN>` — eine Adresse,
keine Anforderung (Baseline-Regelwerk `grundlagen-source-precedence.md`
§ID-Schema als Klammer). Gezählt wird fortlaufend je Datei, nicht je Sektion
(Baseline-Regelwerk `grundlagen-source-precedence.md` §Vergabe).

### SPEC-001 — Recording: Serialisierung und Versionierung

Das Recording ist ein **versionsbehaftetes** Format; `yaml` ist textbasiert und
diff-freundlich. Standard ist **YAML** als menschenlesbare Serialisierung; die
Dateiendung ist standardmäßig `.yaml`. Wahlweise speichert der Recorder eine
SQLite-Datei (`LH-FA-22.a`, `SPEC-043`); beide Formate tragen dasselbe
logische Modell.

Jedes Recording enthält auf oberster Ebene eine Formatkennung:

```yaml
format: pgwire-recorder
version: 1
```

`version` ist eine einzelne ganze Zahl und zählt inkompatible Änderungen des
Formats; abwärtskompatible Ergänzungen (neue optionale Felder) ändern sie nicht.
Version 1 umfasst einfache und Extended-Interaktionen (`type: query`,
`type: extended`); ein Wert von `type`, den ein Leser nicht kennt, macht das
Recording zu einem beschädigten (`PGR-E3003`). Ein Leser lehnt jede `version`
ab, die er nicht kennt (`PGR-E3002`), und ein
Recording ohne oder mit abweichender `format`-Kennung als beschädigt
(`PGR-E3003`); beides endet mit Exit-Code `3` (`SPEC-016`).

### SPEC-002 — Recording: logisches Modell

Beispielhaftes Format:

```yaml
format: pgwire-recorder
version: 1
sessions:
  - id: 1
    startup:
      user: app
      database: app
    interactions:
      - sequence: 1
        request:
          type: query
          sql: "SELECT id, name FROM users ORDER BY id"
        responses:
          - type: row_description
            # protokollrelevante Felder
          - type: data_row
            # protokollrelevante Felder
          - type: command_complete
            tag: "SELECT 2"
          - type: ready_for_query
            tx_status: "I"
```

Das endgültige Schema des Formats `yaml` wird durch Go-Datentypen und Schema-Tests
verbindlich definiert, das des Formats `sqlite` durch `tools/schema/schema.yaml`
(`SPEC-043`). `offset_ms` (`LH-FA-21.a`) und `empty_sessions` (`LH-FA-12.a`)
ergänzen das Beispiel.

### SPEC-041 — Recording: Extended-Interaktion

Eine Extended-Interaktion trägt `type: extended` und eine geordnete Liste
`groups`; jede Gruppe hat die Client-Nachrichten (`client`) und die
Server-Nachrichten (`server`), die auf sie antworten (`LH-FA-18.a`):

```yaml
- sequence: 2
  type: extended
  groups:
    - client:
        - type: parse
          statement: "s1"
          sql: "SELECT name FROM users WHERE id = $1"
          param_types: [23]
        - type: bind
          portal: ""
          statement: "s1"
          param_formats: [0]
          params:
            - text: "1"
          result_formats: [0]
        - type: describe
          target: portal
          name: ""
        - type: execute
          portal: ""
          max_rows: 0
        - type: sync
      server:
        - type: parse_complete
        - type: bind_complete
        - type: row_description
          # protokollrelevante Felder
        - type: data_row
          values:
            - text: "alice"
        - type: command_complete
          tag: "SELECT 1"
        - type: ready_for_query
          tx_status: "I"
```

Die Nachrichtentypen heißen wie die PGWire-Nachrichten in Kleinbuchstaben mit
Unterstrich (`parse`, `bind`, `describe`, `execute`, `close`, `flush`, `sync`,
`parse_complete`, `bind_complete`, `close_complete`, `parameter_description`,
`row_description`, `no_data`, `data_row`, `command_complete`,
`empty_query_response`, `portal_suspended`, `error_response`,
`notice_response`, `ready_for_query`). `describe` und `close` tragen `target`
(`statement` oder `portal`) und `name`; ein NULL-Parameter steht als
`null: true`, Binärwerte folgen `SPEC-003`.

Jede Interaktion, einfach oder Extended, kann das Feld `offset_ms` tragen (ganze
Zahl ≥ 0, `LH-FA-21.a`); es steht neben `sequence` und `type`. Simple-Interaktionen
behalten `request`/`responses` (`SPEC-002`). Formatkennung
und Version (`SPEC-001`) gelten für beide Arten; Binärdaten folgen `SPEC-003`.

### SPEC-043 — Recording: SQLite-Format

Die Tabellenform (Tabellen, Spalten, Typen, Schlüssel) steht ausschließlich im
neutralen Schema-Format von d-migrate in `tools/schema/schema.yaml`
(`schema_format: "1.0"`); das SQL für SQLite entsteht daraus mit `d-migrate schema
generate --target sqlite`, nicht von Hand. Diese Spezifikation legt die Bedeutung
fest:

* `meta` trägt `format` = `pgwire-recorder`, `version` (die Formatversion nach
  `SPEC-001`, unabhängig von der `version` der Schema-Datei) und bei
  `--record-empty-sessions` die Zeile `empty_sessions` = `true`.
* `session.startup` ist der JSON-Text der Startup-Parameter. `session.id` ist die
  fortlaufende Nummer nach `LH-FA-12.a`.
* `interaction.type` ist `query` oder `extended`; `interaction.sql` trägt den
  SQL-Text bei `query` und ist sonst `NULL`; `interaction.offset_ms` ist `NULL` ohne
  Zeitangaben, sonst eine ganze Zahl ≥ 0 (`LH-FA-21.a`).
* `message` trägt die Nachrichten einer Interaktion in `position`-Reihenfolge. Bei
  `query` stehen nur die Server-Nachrichten dort (`direction` = `server`), die
  Anfrage steht in `interaction.sql`; bei `extended` stehen Client- und
  Server-Nachrichten dort mit `group_no` (`SPEC-041`). `kind` nennt die
  Nachrichtenart wie in `SPEC-041`. `bytes` ist die vollständige Nachricht im
  Wire-Format (Binärwerte also ohne Base64), `fields` dieselbe Nachricht als
  JSON-Text zur Prüfung. Beim Lesen gilt `bytes`; widersprechen sich beide, ist die
  Aufzeichnung beschädigt (`PGR-E3003`).
* Die Zugehörigkeit einer Nachricht zu einer Interaktion (`session_id`, `sequence`)
  lässt sich im neutralen Modell nicht als zusammengesetzter Fremdschlüssel
  ausdrücken; der Adapter stellt sie sicher. Die Fremdschlüssel auf `session` sind
  eingeschaltet (`LH-FA-22.a`).

### SPEC-003 — Recording: Binärdaten

PGWire-Felder, die nicht verlustfrei als normaler YAML-String repräsentiert
werden können, werden mit einer eindeutig gekennzeichneten binären
Repräsentation gespeichert, beispielsweise Base64. Replay rekonstruiert die
ursprünglichen Bytes.

### SPEC-004 — Recording: stabile Ausgabe

Die Ausgabe des Formats `yaml` ist deterministisch, soweit keine fachlich notwendige Information
dagegen spricht. Das reduziert unnötige Änderungen in Versionskontrollsystemen.
Ohne `--record-timing` ist die Ausgabe unabhängig vom Zeitverhalten des Laufs, ausgenommen `Flush`-Gruppen ohne wartenden Client (`SPEC-041`, LH-FA-18.a, LH-FA-21.a).
Der Recorder trägt keine rechnerspezifischen Angaben (Dateipfade, Hostnamen,
Adressen) in das Recording ein; vom Client oder Server gelieferte Werte
(Startup-Parameter, `ParameterStatus`) werden unverändert aufgezeichnet. Das
Recording ist damit auf einen anderen Rechner kopierbar.

## 3. Defaults und Konstanten

Regeln dieser Sektion: Die ADR, die einen Wert festlegt, deklariert das
aufwärts in ihrem `Schärft:`-Feld — kein ADR-Rückzeiger hier
(Baseline-Regelwerk `grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP)).
Die `SPEC-<NNN>` ist das, was ihr `Schärft:`-Feld **benennt**; ohne sie kann
eine ADR nur den ganzen Abschnitt nennen. Der Link zeigt weiter auf die
Sektion — Kennungen in Tabellenzellen haben keinen eigenen Anker, und kein
Sensor bemerkt, wenn eine umbenannt wird.

| ID | Name | Wert | Begründung |
|---|---|---|---|
| `SPEC-005` | Standard-Log-Level | `info` | Diagnose ohne Debug-Rauschen (LH-FA-14) |
| `SPEC-006` | Log-Ziel | `stderr` | `stdout` bleibt frei von Logs (LH-FA-14) |
| `SPEC-007` | Konfigurationspriorität | CLI-Argument > Umgebungsvariable > Konfigurationsdatei > Standardwert (Passwort: `LH-FA-17.a`) | eindeutige Auflösung (LH-FA-17) |
| `SPEC-008` | Präfix der Umgebungsvariablen | `PGWIRE_RECORDER_` | eindeutiger Namensraum (LH-FA-17) |
| `SPEC-009` | Dateiendung des Recordings | `.yaml` (Format `yaml`), `.sqlite` üblich (Format `sqlite`); die Endung hat keine Bedeutung für die Erkennung | menschenlesbar, diff-freundlich (LH-FA-07) |
| `SPEC-010` | Formatkennung / Formatversion | `pgwire-recorder` / `1` | Erkennbarkeit inkompatibler Änderungen (LH-QA-06) |
| `SPEC-011` | Replay-Matching | strict sequential (einziges Verfahren in v1) | Determinismus (LH-FA-09, LH-QA-01) |
| `SPEC-012` | `--fail-on-unconsumed` | `false` | nicht verbrauchte Interaktionen sind standardmäßig eine Warnung; Query-Mismatches bleiben immer Fehler (LH-FA-10) |

## 4. Fehler-Codes und Logging-Felder

Regeln dieser Sektion: Der Fehler-Code ist ein Laufzeit-Symbol, die
`SPEC-<NNN>` benennt die *Festlegung* darüber — beide stehen nebeneinander,
sonst hätte eine ADR, die die Fehlerbehandlung schärft, kein Ziel
(Baseline-Regelwerk `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Exit Codes.** Mindestens folgende Semantik wird zugesichert. Die
Implementierung darf intern detailliertere Fehler unterscheiden, solange diese
öffentliche Semantik erhalten bleibt.

| ID | Code | Bedingung | Aktion |
|---|---|---|---|
| `SPEC-013` | 0 | erfolgreicher Programmabschluss, auch kontrolliertes Herunterfahren ohne zuvor aufgetretenen Verbindungsfehler | Prozess endet mit Erfolg |
| `SPEC-014` | 2 | ungültige CLI-Verwendung oder Konfiguration, vorhandenes `--output` ohne `--force` | Fehlertext auf `stderr`, Prozess endet |
| `SPEC-015` | 1 | sonstiger Fehler | Fehlertext auf `stderr`, Prozess endet |
| `SPEC-016` | 3 | Recording-Datei ungültig oder nicht zugreifbar | Fehlertext auf `stderr`, Prozess endet |
| `SPEC-017` | 4 | Netzwerk-/Upstream-Fehler, beim Einspielen auch Anmeldefehler (`PGR-E4005`) und eine Fehlerantwort des Servers (`PGR-E4004`) | Startfehler: Prozess endet; Verbindungsfehler: Verbindung endet, Prozessende nach LH-FA-13.b; Einspielen: Abbruch beziehungsweise Exit-Code 4 am Ende (LH-FA-20.a) |
| `SPEC-018` | 5 | Replay-Mismatch (`PGR-E5001`), Anfrage ohne nicht zugeordnete Session (`PGR-E5003`), bei `--fail-on-unconsumed` auch nicht verbrauchte Interaktionen oder Sessions (`PGR-E5002`) | Diagnose nach LH-FA-10.a, Verbindung endet, Prozessende nach LH-FA-13.b |
| `SPEC-019` | 6 | nicht unterstützte PGWire-Funktion | Diagnose, Verbindung endet, Prozessende nach LH-FA-13.b |

**Fehlerklassen.** Mindestens folgende Klassen werden unterschieden und
erzeugen einen für Entwickler verständlichen Text:

| ID | Fehlerklasse | Exit Code | Meldungscode |
|---|---|---|---|
| `SPEC-020` | ungültige CLI-Verwendung oder Konfiguration | 2 | `PGR-E2001`, `PGR-E2004`, `PGR-E2005`, `PGR-E2006` |
| `SPEC-021` | Listen-Port kann nicht geöffnet werden | 4 | `PGR-E4001` |
| `SPEC-022` | Upstream nicht erreichbar oder lehnt ab | 4 | `PGR-E4002`, `PGR-E4004`, `PGR-E4005` |
| `SPEC-023` | Recording kann nicht gelesen/geschrieben werden | 3 | `PGR-E3001` |
| `SPEC-024` | unbekannte Recording-Version | 3 | `PGR-E3002` |
| `SPEC-025` | beschädigtes Recording | 3 | `PGR-E3003` |
| `SPEC-026` | nicht unterstützte PGWire-Nachricht | 6 | `PGR-E6001` |
| `SPEC-027` | Replay-Mismatch | 5 | `PGR-E5001`, `PGR-E5003` |
| `SPEC-028` | unerwartetes Verbindungsende | 4 | `PGR-E4003` |

**Logging-Felder.** Log-Level, Ziel und Geheimnisschutz stehen in LH-FA-14.a
und `SPEC-005`/`SPEC-006`. Fehler und Warnungen tragen zusätzlich einen
Meldungscode (`SPEC-034`).

### SPEC-034 — Meldungscodes

Jede klassifizierte Fehlerursache und jede Warnung mit einer Maßnahme für den
Anwender trägt einen Meldungscode der Form `PGR-<S><NNNN>` (ERE
`PGR-[EWI][0-9]{4}`). `S` ist die Schwere: `E` Fehler, `W` Warnung, `I` ist
reserviert.

**Fehler.** Die erste Ziffer ist der Exit Code der Fehlerklasse (`SPEC-013` bis
`SPEC-019`); die Klasse folgt aus dem Code. Die Endung `000` ist der Rückfall
der Klasse: ein klassifizierter Fehler ohne Einzelursache trägt den Rückfall
seiner Klasse, nie keinen Code. Ein Fehler ohne Klasse ist `PGR-E1000`.

**Warnungen.** Die erste Ziffer ist der Bereich: 1 Record, 2 Replay,
3 Protokollrand (beide Modi), 5 Konfiguration und Start, 9 reserviert. Eine Warnung trägt keine Klasse und
ändert den Exit Code nicht.

| Code | Klasse / Bereich | Bedeutung |
|---|---|---|
| `PGR-E1000` | sonstiger Fehler (Exit 1) | Rückfall; unerwarteter interner Fehler |
| `PGR-E2000` | Konfiguration (Exit 2) | Rückfall |
| `PGR-E2001` | Konfiguration (Exit 2) | ungültige CLI-Verwendung: unbekannte Option, fehlende Pflichtoption, ungültiger Wert oder ungültige Kombination (`SPEC-020`) |
| `PGR-E2002` | Konfiguration (Exit 2) | `--output` existiert bereits und `--force` ist nicht gesetzt (LH-FA-07.a) |
| `PGR-E2003` | Konfiguration (Exit 2) | zeitgetreues Einspielen verlangt, aber mindestens einer Interaktion der Aufzeichnung fehlt `offset_ms` (LH-FA-21.a) |
| `PGR-E2004` | Konfiguration (Exit 2) | Konfigurationsdatei nicht lesbar oder ungültig: YAML, Schlüssel, Abschnitt, Wert, `sslmode`, Platzhalter außerhalb einer URL (LH-FA-17.a) |
| `PGR-E2005` | Konfiguration (Exit 2) | Umgebungsvariable eines Platzhalters der benutzten Verbindung nicht gesetzt (LH-FA-17.a) |
| `PGR-E2006` | Konfiguration (Exit 2) | Klartext-Passwort in der Konfigurationsdatei (LH-FA-17.a) |
| `PGR-E3000` | Recording (Exit 3) | Rückfall |
| `PGR-E3001` | Recording (Exit 3) | Recording nicht lesbar oder nicht schreibbar (`SPEC-023`) |
| `PGR-E3002` | Recording (Exit 3) | unbekannte Recording-Version (`SPEC-024`) |
| `PGR-E3003` | Recording (Exit 3) | beschädigtes Recording (`SPEC-025`) |
| `PGR-E3004` | Recording (Exit 3) | Recording ohne Session, die sich im Replay zuordnen lässt (LH-FA-03.a) |
| `PGR-E4000` | Netzwerk (Exit 4) | Rückfall |
| `PGR-E4001` | Netzwerk (Exit 4) | Listen-Port nicht zu öffnen (`SPEC-021`) |
| `PGR-E4002` | Netzwerk (Exit 4) | Upstream nicht erreichbar (`SPEC-022`) |
| `PGR-E4003` | Netzwerk (Exit 4) | unerwartetes Verbindungsende (`SPEC-028`) |
| `PGR-E4004` | Netzwerk (Exit 4) | Server beantwortet eine eingespielte Anfrage mit einem Fehler (LH-FA-20.a) |
| `PGR-E4005` | Netzwerk (Exit 4) | Anmeldung am Server fehlgeschlagen, Verfahren nicht unterstützt, TLS abgelehnt oder ungültiges Zertifikat, oder der Server lehnt unverschlüsselte Verbindungen ab (LH-FA-20.a) |
| `PGR-E5000` | Replay (Exit 5) | Rückfall |
| `PGR-E5001` | Replay (Exit 5) | Replay-Mismatch (`SPEC-027`) |
| `PGR-E5002` | Replay (Exit 5) | nicht verbrauchte Interaktionen oder Sessions bei `--fail-on-unconsumed` (LH-FA-03.b) |
| `PGR-E5003` | Replay (Exit 5) | Anfrage einer Verbindung, zu der keine nicht zugeordnete Session mehr existiert (LH-FA-12.a, `SPEC-027`) |
| `PGR-E6000` | nicht unterstützt (Exit 6) | Rückfall |
| `PGR-E6001` | nicht unterstützt (Exit 6) | nicht unterstützte PGWire-Nachricht (`SPEC-026`) |
| `PGR-E6002` | nicht unterstützt (Exit 6) | nicht unterstützte PGWire-Protokollversion (LH-FA-05.e) |
| `PGR-W2001` | Replay | Sitzung endet vor Verbrauch aller Interaktionen (LH-FA-03.b) |
| `PGR-W3001` | Protokollrand | `CancelRequest` empfangen und nicht weitergeleitet (LH-FA-05.e) |

**Ausgabe.** Der Fehlertext (Fehlerwert, Zeile beim Prozessende, Attribut
`error` einer Log-Zeile) beginnt mit dem Kopf `<klasse> [<code>]: <Ursache>`
auf `stderr`. Eine Warnung steht als eigenes Log-Attribut `code`; der
Meldungstext trägt keinen Code. Der Prozessausgang bleibt unverändert, kein Code
ändert ihn.

**Stabilität.** Ein Code wird nie neu belegt; ein entfallener Code wird
zurückgezogen und bleibt vergeben; die Klasse eines Codes ändert sich nie.
Stabil sind Code, Klasse und Prozessausgang; der Text nach dem Kopf ist nicht
Vertrag. Die Code-Tabelle liegt im Quelltext, der Katalog (Code, Klasse,
Bedeutung, Maßnahme) in der Betriebsdokumentation.

## 5. Metriken und Tracing-Felder

Regeln dieser Sektion: verbindliche OTel-Felder pro Span
(Baseline-Regelwerk `modul-15-observability.md`).

Für v1 sind keine Metriken und Tracing-Felder festgelegt.

## 6. Externe Verträge

| ID | System | Version | Vertrag-Datei |
|---|---|---|---|
| `SPEC-029` | PostgreSQL (Upstream im Record-Modus), PGWire | PGWire-Protokollversion 3.0; jeder Server, der sie spricht | — |
| `SPEC-030` | `github.com/jackc/pgx/v5/pgproto3` (Verarbeitung von PGWire-Nachrichten) | Major 5 | — |
| `SPEC-031` | Docker/OCI-Image `ghcr.io/pt9912/pgwire-recorder` und `docker.io/pt9912/pgwire-recorder` | OCI-Image-Format, `linux/amd64` und `linux/arm64` | — |
| `SPEC-042` | Homebrew-Tap `pt9912/homebrew-pgwire-recorder` | eine Formel je stabilem Release | — |
| `SPEC-044` | SQLite (Dateiformat der Aufzeichnung) | SQLite 3 | — |

## 7. Querschnittsvorgaben

### SPEC-032 — Performance

v1 ist ein Testwerkzeug und kein Produktions-Connection-Pooler. Es gelten
folgende Ziele:

- Streaming statt vollständigem Puffern großer Resultsets, soweit mit
  Recording-Format und Konsistenz vereinbar,
- keine zusätzliche Kopie von Resultsets neben der Aufzeichnung; eine Session
  wird bis zum Schreiben des Recordings im Speicher gehalten,
- zusätzlicher Proxy-Overhead bleibt für typische Integrationstests praktisch
  vertretbar.

Für v1 werden keine harten Latenz- oder Durchsatz-SLAs zugesichert.

### SPEC-033 — Sicherheit

- Das Tool wird nicht als Produktions-Sicherheitskomponente positioniert.
- Recordings können sensible Daten enthalten (LH-RB-01); es gibt keine
  automatische Maskierung.
- Dateirechte und Repository-Regeln liegen in der Verantwortung des Anwenders.
- Logs erzeugen keine zusätzlichen Geheimnisse, die nicht für Diagnosezwecke
  erforderlich sind.
- Replay lauscht standardmäßig nur auf der explizit angegebenen Listen-Adresse.
- Diagnosen nennen SQL-Text im Klartext, auch mit Literalen darin; verborgen
  werden nur die Parameterwerte der Extended-Interaktion (LH-FA-18.a).

### SPEC-035 — Zielplattformen

Das Binary ist für Linux, macOS und Windows bereitgestellt, jeweils für `amd64`
und `arm64`. Das Container-Image ist ein Linux-Image für `amd64` und `arm64`.
Weitere Umgebungen sind nicht ausgeschlossen, aber nicht zugesichert.

## 8. Technische Leitentscheidungen und Architekturvorgabe

### SPEC-036 — Technische Leitentscheidungen

Für v1 gelten folgende Festlegungen:

- Implementierung in **Go**.
- Verwendung von `github.com/jackc/pgx/v5/pgproto3` für die Verarbeitung von
  PGWire-Nachrichten, soweit die Bibliothek die benötigten Nachrichtentypen
  abdeckt (`SPEC-030`).
- Unterstützung des **Simple Query Protocol** und des **Extended Query
  Protocol**.
- Das einzige Replay-Matching ist **strict sequential** (`SPEC-011`).
- Keine automatische Maskierung oder Redaktion von Recording-Inhalten.
- Das Tool ist für lokale Entwicklung, automatisierte Tests, CI und
  Containerbetrieb ausgelegt.

### SPEC-037 — Architekturvorgabe

Die Implementierung folgt einer **hexagonalen Architektur
(Ports & Adapters)** mit folgender Terminologie:

```text
Driving Adapter -> Inbound Port -> Application Core
Application Core -> Outbound Port -> Driven Adapter
```

Konkrete Infrastrukturabhängigkeiten wie `pgproto3`, YAML, TCP und Dateisystem
sind nicht Bestandteil des Domain Models. Detaillierte Regeln und
Package-Grenzen stehen in `spec/architecture.md`.

**Architektur-Gate.** Die hexagonalen Abhängigkeitsregeln sind die Vorgabe für
ein Architektur-Gate (a-check, Konfiguration `.a-check.yml`), das Verletzungen
als fehlgeschlagenen Check meldet. Driving Adapter verwenden Inbound Ports;
Driven Adapter implementieren beziehungsweise verwenden Outbound Ports.
Infrastrukturtechnologien wie `pgproto3` und YAML sind auf die dafür
vorgesehenen Adapter begrenzt.

## 9. Testanforderungen

### SPEC-038 — Unit-, Integrations- und Kompatibilitätstests

**Unit Tests** decken mindestens ab:

- Recording-Serialisierung und -Deserialisierung,
- Versionsprüfung,
- Query-Matching,
- Mismatch-Diagnosen,
- Abbildung unterstützter PGWire-Nachrichten,
- Exit-Code-Zuordnung.

**Integrationstests** starten einen realen PostgreSQL-Server und prüfen
mindestens:

1. Client → Recorder → PostgreSQL im Record-Modus,
2. erzeugtes Recording,
3. PostgreSQL wird beendet,
4. Client → Recorder im Replay-Modus,
5. identisches beobachtbares Ergebnis,
6. absichtlicher Query-Mismatch führt zum erwarteten Fehler.

**Kompatibilitätstests.** Mindestens ein verbreiteter PostgreSQL-Go-Client wird
in seinem Standardmodus (Extended Query) und im Simple-Query-Modus für
End-to-End-Tests verwendet. Zusätzliche Clients anderer Sprachen sind
erwünscht, aber kein v1-Muss.

## 10. Nicht zugesichert in v1

### SPEC-039 — Abgrenzung

Nicht zugesichert sind insbesondere:

- COPY,
- TLS-Terminierung,
- CancelRequest,
- LISTEN/NOTIFY als speziell getesteter Anwendungsfall,
- Replikationsprotokoll,
- Connection Pooling,
- Query-Normalisierung,
- Fuzzy Matching,
- dynamische Mock-Regeln,
- Recording-Editor,
- automatische Secret-/PII-Maskierung,
- garantierter Multi-Session-Replay bei parallelen Clients.

Diese Punkte benötigen vor Aufnahme in den Produktumfang eigene Anforderungen
und Tests. „Nicht zugesichert" heißt: Die Funktion fehlt; Nachrichten und
Interaktionen daraus werden nach LH-FA-05.e abgelehnt beziehungsweise
behandelt, nicht still ignoriert.

## 11. Historie

Regeln dieser Sektion: **kein ADR- und kein Slice-Verweis.** Die Decken-Regel
gilt für alle drei Spec-Straten, auch hier — welche ADR eine Festlegung
schärft, deklariert die ADR aufwärts in ihrem `Schärft:`-Feld
(Baseline-Regelwerk `modul-03-spec.md` §Ziel-Form: Spezifikation).

| Datum | Änderung |
|---|---|
| 2026-10-03 | Initial |
