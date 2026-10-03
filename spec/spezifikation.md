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

**Fehlermodi:** Listen-Port nicht zu öffnen und Upstream nicht erreichbar →
Exit-Code `4` (`SPEC-017`).

---

### LH-FA-02.b — Abschluss einer Interaktion

Eine Simple-Query-Interaktion gilt als abgeschlossen, sobald das zugehörige
`ReadyForQuery` des Servers verarbeitet wurde.

---

### LH-FA-03.a — Replay: Aufruf, Upstream-Freiheit, Session-Cursor

**Eingabe:** Client-Verbindungen auf der `--listen`-Adresse und ein Recording.
**Ausgabe:** aufgezeichnete Backend-Nachrichten.

Verpflichtende Optionen von `replay`:

| Option | Bedeutung |
|---|---|
| `--listen` | Adresse, auf der der Replay-Server Clients annimmt |
| `--input` | Zu verwendendes Recording |

Beispiel:

```bash
pgwire-recorder replay \
  --listen 0.0.0.0:15432 \
  --input ./recordings/users.yaml
```

**Schritte:**

1. Im Replay-Modus wird keine Verbindung zu einem realen PostgreSQL-Server
   benötigt.
2. Jede eingehende Client-Verbindung erhält einen eigenen Replay-Cursor.
3. Der Cursor zeigt auf die nächste erwartete Interaktion der zugeordneten
   Recording-Session.

---

### LH-FA-03.b — Nicht verbrauchte Interaktionen

Wird eine Replay-Session beendet, bevor alle ihr zugeordneten Interaktionen
verbraucht wurden, wird dies mindestens als Warnung mit dem Meldungscode `PGR-W2001` protokolliert (`SPEC-034`).

Für CI ist optional ein Strictness-Schalter vorgesehen, mit dem dies als Fehler
gewertet wird. Der Default dieses Zusatzschalters wird bei der Implementierung
festgelegt (`SPEC-012`); Query-Mismatches selbst bleiben immer Fehler.

---

### LH-FA-05.a — Unterstützter Umfang: Simple Query Protocol

Startup, Simple Query Protocol und Serverantworten bilden den unterstützten
Umfang von v1.

**Startup.** Der Recorder unterstützt die für typische PostgreSQL-Clients
notwendige Startup-Sequenz, soweit sie für die definierten Record-/Replay-
Szenarien erforderlich ist. Startup-Parameter werden als Teil einer Session
erfasst, soweit sie für das Replay relevant sind.

**Client-Nachrichten.** Für den Kern von v1 sind relevant:

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

Die konkrete Menge ist nicht auf diese Liste beschränkt.

**Nicht unterstützte Protokollnachrichten.** Trifft v1 auf eine
Client-Interaktion, die als nicht unterstützt klassifiziert ist, schlägt das Tool
kontrolliert fehl beziehungsweise beendet die betroffene Verbindung mit einer
verständlichen Diagnose (Exit-Code `6`, `SPEC-019`). Es wird keine scheinbar
erfolgreiche Aufzeichnung erzeugt, wenn für die betreffende Session wesentliche
Nachrichten nicht korrekt verarbeitet wurden.

**Extended Query Protocol.** Nicht Teil des zugesicherten v1-Scopes sind
insbesondere `Parse`, `Bind`, `Describe`, `Execute`, `Sync` und Prepared
Statements über diesen Ablauf. Eine spätere Erweiterung muss möglich bleiben.

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

Versucht ein Client eine in v1 nicht unterstützte SSL/TLS-Aushandlung, ist das
Verhalten eindeutig und dokumentiert. Für den initialen v1-Umfang darf der
Recorder TLS ablehnen und einen unverschlüsselten PGWire-Verbindungsaufbau
verlangen. TLS-Unterstützung ist eine spätere Erweiterung.

---

### LH-FA-05.d — Transaktionen

SQL-Transaktionsbefehle, die über Simple Query übertragen werden, werden wie
andere Query-Nachrichten aufgezeichnet und strikt sequenziell wiedergegeben.

Der Recorder implementiert im Replay-Modus keine eigene Transaktionslogik. Der
beobachtbare Transaktionsstatus wird durch die aufgezeichneten Serverantworten,
insbesondere `ReadyForQuery`, reproduziert.

---

### LH-FA-06.a — Inhalt der Aufzeichnung

Für jede unterstützte Query wird mindestens gespeichert:

- Session-Zuordnung,
- fortlaufende Interaktionsnummer,
- Anfrageart,
- SQL-Text,
- geordnete Folge der zugehörigen Serverantworten,
- für das Replay notwendige PGWire-Felder.

Zeitstempel dürfen zur Diagnose gespeichert werden, sind aber keine
Voraussetzung für ein erfolgreiches Replay.

---

### LH-FA-07.a — Sicheres Schreiben des Recordings

**Eingabe:** abgeschlossene oder abgebrochene Record-Läufe. **Ausgabe:** Datei
unter dem `--output`-Pfad.

**Schritte:**

1. Die Implementierung schreibt zunächst in eine temporäre Datei.
2. Beim erfolgreichen Abschluss wird sie atomar beziehungsweise bestmöglich
   atomar auf die Zieldatei verschoben.

Ein erfolgreich beendeter Record-Lauf hinterlässt keine syntaktisch
unvollständige Recording-Datei. Bei einem Abbruch darf eine temporäre
beziehungsweise als unvollständig erkennbare Datei verbleiben; sie wird nicht
stillschweigend als gültiges Recording behandelt.

**Fehlermodi:** Recording nicht lesbar/schreibbar, unbekannte Version oder
beschädigt → Exit-Code `3` (`SPEC-016`).

---

### LH-FA-08.a — Auswahl des Recordings

Das Recording wird über `--output` (Record) beziehungsweise `--input` (Replay)
bestimmt (siehe LH-FA-02.a und LH-FA-03.a).

---

### LH-FA-09.a — Strict sequential matching

v1 verwendet **strict sequential matching**; Replay ist standardmäßig strict
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
durch. Sendet ein Client weitere Queries, obwohl keine aufgezeichnete
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
Recording (`SPEC-002`) erhalten die Reihenfolge der Interaktionen.

**Recording.** Mehrere Client-Verbindungen dürfen parallel angenommen werden.
Jede Verbindung wird als separate Session im Recording geführt.

**Replay.** v1 garantiert deterministisches Replay für Recordings mit einer
einzelnen Session. Mehrere aufgezeichnete Sessions dürfen im Format
repräsentiert werden, gelten aber erst dann als vollständig unterstützter
Replay-Anwendungsfall, wenn eine eindeutige Session-Zuordnungsstrategie
implementiert und getestet ist. Damit wird verhindert, dass v1 bei parallelen
Verbindungen nichtdeterministisch eine falsche Session auswählt.

---

### LH-FA-13.a — Signalbehandlung

Auf `SIGINT` und `SIGTERM` fährt der Prozess kontrolliert herunter:

- keine neuen Verbindungen annehmen,
- laufende Schreiboperationen soweit möglich abschließen,
- Recording konsistent abschließen, wenn der Zustand dies erlaubt,
- andernfalls das Recording als unvollständig behandeln (siehe LH-FA-07.a).

---

### LH-FA-14.a — Logging und Diagnose

Logs werden nach `stderr` geschrieben (`SPEC-006`); Nutzdaten beziehungsweise
maschinenlesbare Ausgaben auf `stdout` werden dadurch nicht verunreinigt.
Unterstützte Log-Level sind mindestens `error`, `warn`, `info`, `debug`;
Standard ist `info` (`SPEC-005`).

Passwörter aus Verbindungsdaten werden nicht absichtlich in Logs ausgegeben.

Fehler und Warnungen tragen einen Meldungscode und erzeugen einen für Entwickler verständlichen Text (`SPEC-034`).

---

### LH-FA-15.a — Nicht-interaktiver Betrieb

Alle Kommandos sind vollständig nicht-interaktiv; sie fragen keine Eingaben ab
(siehe LH-FA-17.a).

---

### LH-FA-16.a — Containerbetrieb

Das Projekt kann ein Container-Image bereitstellen (`SPEC-031`). Das Binary hat
keine Laufzeitabhängigkeit auf eine lokale PostgreSQL-Installation.

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

Alle für CI erforderlichen Einstellungen sind über CLI-Argumente verfügbar.
Umgebungsvariablen dürfen ergänzend unterstützt werden; ihre Namen tragen das
Präfix `PGWIRE_RECORDER_` (`SPEC-008`). Wird dieselbe Einstellung mehrfach
angegeben, gilt (`SPEC-007`):

```text
CLI-Argument > Umgebungsvariable > Standardwert
```

---

## 2. Datenstrukturen und Schemas

Regeln dieser Sektion: Jede Struktur trägt eine `SPEC-<NNN>` — eine Adresse,
keine Anforderung (Baseline-Regelwerk `grundlagen-source-precedence.md`
§ID-Schema als Klammer). Gezählt wird fortlaufend je Datei, nicht je Sektion
(Baseline-Regelwerk `grundlagen-source-precedence.md` §Vergabe).

### SPEC-001 — Recording: Serialisierung und Versionierung

Das Recording ist ein **versionsbehaftetes, textbasiertes und diff-freundliches**
Format. Für v1 wird **YAML** als menschenlesbare Serialisierung verwendet; die
Dateiendung ist standardmäßig `.yaml`.

Jedes Recording enthält auf oberster Ebene eine Formatkennung:

```yaml
format: pgwire-recorder
version: 1
```

Unbekannte Major-Versionen werden abgelehnt (Exit-Code `3`, `SPEC-016`).

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

Das endgültige Schema wird durch Go-Datentypen und Schema-Tests verbindlich
definiert.

### SPEC-003 — Recording: Binärdaten

PGWire-Felder, die nicht verlustfrei als normaler YAML-String repräsentiert
werden können, werden mit einer eindeutig gekennzeichneten binären
Repräsentation gespeichert, beispielsweise Base64. Replay rekonstruiert die
ursprünglichen Bytes.

### SPEC-004 — Recording: stabile Ausgabe

Die Ausgabe ist deterministisch, soweit keine fachlich notwendige Information
dagegen spricht. Das reduziert unnötige Änderungen in Versionskontrollsystemen.

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
| `SPEC-007` | Konfigurationspriorität | CLI-Argument > Umgebungsvariable > Standardwert | eindeutige Auflösung (LH-FA-17) |
| `SPEC-008` | Präfix der Umgebungsvariablen | `PGWIRE_RECORDER_` | eindeutiger Namensraum (LH-FA-17) |
| `SPEC-009` | Dateiendung des Recordings | `.yaml` | menschenlesbar, diff-freundlich (LH-FA-07) |
| `SPEC-010` | Formatkennung / Formatversion | `pgwire-recorder` / `1` | Erkennbarkeit inkompatibler Änderungen (LH-QA-06) |
| `SPEC-011` | Replay-Matching | strict sequential | Determinismus (LH-FA-09, LH-QA-01) |
| `SPEC-012` | Strictness-Schalter für nicht verbrauchte Interaktionen | Default bei Implementierung festzulegen | CI-Eignung; Query-Mismatches bleiben immer Fehler (LH-FA-10) |

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
| `SPEC-013` | 0 | erfolgreicher Programmabschluss | Prozess endet mit Erfolg |
| `SPEC-014` | 2 | ungültige CLI-Verwendung oder Konfiguration | Fehlertext auf `stderr`, Prozess endet |
| `SPEC-015` | 1 | sonstiger Fehler | Fehlertext auf `stderr`, Prozess endet |
| `SPEC-016` | 3 | Recording-Datei ungültig oder nicht zugreifbar | Fehlertext auf `stderr`, Prozess endet |
| `SPEC-017` | 4 | Netzwerk-/Upstream-Fehler | Fehlertext auf `stderr`, Prozess endet |
| `SPEC-018` | 5 | Replay-Mismatch | Diagnose nach LH-FA-10.a, Prozess endet |
| `SPEC-019` | 6 | nicht unterstützte PGWire-Funktion | Diagnose, Verbindung wird beendet beziehungsweise Prozess endet |

**Fehlerklassen.** Mindestens folgende Klassen werden unterschieden und
erzeugen einen für Entwickler verständlichen Text:

| ID | Fehlerklasse | Exit Code | Meldungscode |
|---|---|---|---|
| `SPEC-020` | ungültige CLI-Konfiguration | 2 | `PGR-E2001` |
| `SPEC-021` | Listen-Port kann nicht geöffnet werden | 4 | `PGR-E4001` |
| `SPEC-022` | Upstream nicht erreichbar | 4 | `PGR-E4002` |
| `SPEC-023` | Recording kann nicht gelesen/geschrieben werden | 3 | `PGR-E3001` |
| `SPEC-024` | unbekannte Recording-Version | 3 | `PGR-E3002` |
| `SPEC-025` | beschädigtes Recording | 3 | `PGR-E3003` |
| `SPEC-026` | nicht unterstützte PGWire-Nachricht | 6 | `PGR-E6001` |
| `SPEC-027` | Replay-Mismatch | 5 | `PGR-E5001` |
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
5 Konfiguration und Start, 9 reserviert. Eine Warnung trägt keine Klasse und
ändert den Exit Code nicht.

| Code | Klasse / Bereich | Bedeutung |
|---|---|---|
| `PGR-E1000` | sonstiger Fehler (Exit 1) | Rückfall; unerwarteter interner Fehler |
| `PGR-E2000` | Konfiguration (Exit 2) | Rückfall |
| `PGR-E2001` | Konfiguration (Exit 2) | ungültige CLI-Verwendung oder Konfiguration (`SPEC-020`) |
| `PGR-E3000` | Recording (Exit 3) | Rückfall |
| `PGR-E3001` | Recording (Exit 3) | Recording nicht lesbar oder nicht schreibbar (`SPEC-023`) |
| `PGR-E3002` | Recording (Exit 3) | unbekannte Recording-Version (`SPEC-024`) |
| `PGR-E3003` | Recording (Exit 3) | beschädigtes Recording (`SPEC-025`) |
| `PGR-E4000` | Netzwerk (Exit 4) | Rückfall |
| `PGR-E4001` | Netzwerk (Exit 4) | Listen-Port nicht zu öffnen (`SPEC-021`) |
| `PGR-E4002` | Netzwerk (Exit 4) | Upstream nicht erreichbar (`SPEC-022`) |
| `PGR-E4003` | Netzwerk (Exit 4) | unerwartetes Verbindungsende (`SPEC-028`) |
| `PGR-E5000` | Replay (Exit 5) | Rückfall |
| `PGR-E5001` | Replay (Exit 5) | Replay-Mismatch (`SPEC-027`) |
| `PGR-E6000` | nicht unterstützt (Exit 6) | Rückfall |
| `PGR-E6001` | nicht unterstützt (Exit 6) | nicht unterstützte PGWire-Nachricht (`SPEC-026`) |
| `PGR-W2001` | Replay | Sitzung endet vor Verbrauch aller Interaktionen (LH-FA-03.b) |

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
| `SPEC-029` | PostgreSQL (Upstream im Record-Modus), PGWire | unterstützte Versionen offen | — |
| `SPEC-030` | `github.com/jackc/pgx/v5/pgproto3` (Verarbeitung von PGWire-Nachrichten) | Major 5 | — |
| `SPEC-031` | Container-Image `pgwire-recorder` | — | — |

## 7. Querschnittsvorgaben

### SPEC-032 — Performance

v1 ist ein Testwerkzeug und kein Produktions-Connection-Pooler. Es gelten
folgende Ziele:

- Streaming statt vollständigem Puffern großer Resultsets, soweit mit
  Recording-Format und Konsistenz vereinbar,
- keine künstliche Speicherung kompletter Datenbanksitzungen im RAM,
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

## 8. Technische Leitentscheidungen und Architekturvorgabe

Für v1 gelten folgende Festlegungen:

- Implementierung in **Go**.
- Verwendung von `github.com/jackc/pgx/v5/pgproto3` für die Verarbeitung von
  PGWire-Nachrichten, soweit die Bibliothek die benötigten Nachrichtentypen
  abdeckt (`SPEC-030`).
- Unterstützung des **Simple Query Protocol** als fachlicher Kern; das Extended
  Query Protocol ist keine zugesicherte v1-Funktion.
- Replay ist standardmäßig **strict** (`SPEC-011`).
- Keine automatische Maskierung oder Redaktion von Recording-Inhalten.
- Das Tool ist für lokale Entwicklung, automatisierte Tests, CI und
  Containerbetrieb ausgelegt.

**Architektur.** Die Implementierung folgt einer **hexagonalen Architektur
(Ports & Adapters)** mit folgender Terminologie:

```text
Driving Adapter -> Inbound Port -> Application Core
Application Core -> Outbound Port -> Driven Adapter
```

Konkrete Infrastrukturabhängigkeiten wie `pgproto3`, YAML, TCP und Dateisystem
sind nicht Bestandteil des Domain Models. Detaillierte Regeln und
Package-Grenzen stehen in `spec/architecture.md`.

**Architektur-Gate.** Die hexagonalen Abhängigkeitsregeln werden über eine
`.a-check.yml` mit **a-check** automatisiert geprüft. Das Gate ist für CI
vorgesehen und meldet Architekturverletzungen als fehlgeschlagenen Check.
Driving Adapter verwenden Inbound Ports; Driven Adapter implementieren
beziehungsweise verwenden Outbound Ports. Infrastrukturtechnologien wie
`pgproto3` und YAML sind auf die dafür vorgesehenen Adapter begrenzt.

## 9. Testanforderungen

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
für End-to-End-Tests verwendet. Zusätzliche Clients anderer Sprachen sind
erwünscht, aber kein v1-Muss.

## 10. Nicht zugesichert in v1

Nicht zugesichert sind insbesondere:

- Extended Query Protocol,
- Prepared Statements,
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
und Tests.

## 11. Fertigstellungskriterien für v1

v1 ist technisch fertig, wenn:

1. alle Muss-Abnahmekriterien des Lastenhefts automatisiert getestet sind,
2. Record und Replay für mindestens einen realen PostgreSQL-End-to-End-Test
   funktionieren,
3. Replay ohne laufenden PostgreSQL-Server funktioniert,
4. Simple Queries mit Resultsets, Commands ohne Rows und PostgreSQL-Fehlern
   abgedeckt sind,
5. Mismatches deterministisch erkannt werden,
6. das Recording-Format versioniert und dokumentiert ist,
7. CLI-Hilfe, Exit Codes und grundlegende Betriebsdokumentation vorhanden sind,
8. das Binary und ein Container-Image reproduzierbar gebaut werden können.

## 12. Historie

Regeln dieser Sektion: **kein ADR- und kein Slice-Verweis.** Die Decken-Regel
gilt für alle drei Spec-Straten, auch hier — welche ADR eine Festlegung
schärft, deklariert die ADR aufwärts in ihrem `Schärft:`-Feld
(Baseline-Regelwerk `modul-03-spec.md` §Ziel-Form: Spezifikation).

| Datum | Änderung |
|---|---|
| 2026-10-03 | Initial |
