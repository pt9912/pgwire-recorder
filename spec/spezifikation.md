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

**Hilfe vor Prüfung.** Eine *Hilfe-Angabe* ist ein Argument, das genau eine der
vier Formen `-h`, `--h`, `-help` oder `--help` lautet, ohne oder mit `=` und
beliebigem Wert dahinter, auch leerem oder `false` (`--help=`, `-h=x`,
`-help=false` sind Hilfe-Angaben); andere Schreibweisen (`---help`, `-H`,
`--Help`) sind es nicht. Die Liste ist abgeschlossen. Steht unter den Argumenten
vor einem `--` eine Hilfe-Angabe, an welcher Stelle auch immer, gibt das
Programm die Hilfe auf `stdout` aus und endet mit Exit-Code `0`; das gilt auch
vor dem Kommando (`--listen x --help`), ohne Kommando und nach einem unbekannten
Kommando (`bogus --help`). Ist das erste Argument ein bekanntes Kommando, ist es
dessen Hilfe, sonst die globale. Es prüft dann nichts weiter: kein Kommando,
keine Option, auch keine unbekannte oder ungültige und keine fehlende
Pflichtoption, gleich ob vor oder nach der Hilfe-Angabe, keine
Umgebungsvariable und keine Konfigurationsdatei (LH-FA-17.a). Eine Hilfe-Angabe
an der Stelle eines Optionswerts ist die Hilfe, nicht der Wert (`--input -help`
gibt die Hilfe aus); ein solcher Wert geht nur mit `=` (`--input=--help`), denn
eine Hilfe-Angabe ist nur ein ganzes Argument. Nach `--` ist sie ein
gewöhnliches Argument.

**Ende der Optionen.** Das erste Argument, das genau `--` lautet, beendet die
Optionen, an jeder Stelle, auch an der Stelle eines Optionswerts; für die
Hilfe-Suche und das Lesen der Optionen gilt dieselbe Lesart. Was danach steht,
ist ein gewöhnliches Argument; `record` und `replay` nehmen keines an
(`PGR-E2001`, unerwartetes Argument). Fehlt einer Option vor `--` dadurch ihr
Wert (`--input --`), ist das `PGR-E2001` (Option ohne Wert). Der Wert `--` geht
nur mit `=` (`--input=--`). Kein Fehlertext nennt die Hilfe-Anforderung des
Parsers; eine Hilfe-Angabe ergibt immer die Hilfe, nie einen Fehler.
`version` liest weder Umgebungsvariablen noch Konfigurationsdatei; eine Option
`--version` gibt es nicht (`PGR-E2001`).

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
LH-FA-21.a), `--record-empty-sessions` (siehe LH-FA-12.a) sowie `--tls-cert`,
`--tls-key` und `--allow-plaintext` (siehe LH-FA-23.a).

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
das, ohne Fehlerantwort davor, ein unerwartetes Verbindungsende (Verbindungsfehler, `PGR-E4003`); die
unvollständige Interaktion wird nicht in das Recording übernommen, die
vorherigen Interaktionen der Session bleiben erhalten.

**Fehlerantwort vor dem Abbruch.** Endet die Upstream-Verbindung vor dem
`ReadyForQuery` einer laufenden Interaktion, nachdem der Upstream in dieser
Interaktion eine `ErrorResponse` gesendet hat, ist die Interaktion nicht
unterstützt (`PGR-E6001`, Exit-Code `6`, LH-FA-13.b):

* *Gelesen.* Maßgeblich ist eine `ErrorResponse` der laufenden Interaktion, die
  der Recorder gelesen hat, bevor er das Verbindungsende bemerkt. Gleich ist,
  ob er das Ende beim Lesen vom Upstream oder beim Senden an ihn bemerkt, und bei
  einer Extended-Interaktion, welche der beiden Richtungen es zuerst bemerkt:
  Beide stufen nach demselben Stand ein, und die Session wird bei `PGR-E6001`
  in jedem Fall nicht übernommen. Eine `ErrorResponse`, die der Recorder bis
  dahin nicht gelesen hat, zählt nicht; das Verbindungsende ist dann
  `PGR-E4003`. Das gilt etwa, wenn der Upstream nach einer `ErrorResponse`
  zwischen zwei Interaktionen mit einem TCP-Reset schließt und schon das Senden
  der nächsten Anfrage scheitert: Nach einem Sendefehler liest der Recorder
  nicht weiter.
* *Verbindungsende* ist das Ende des Datenstroms vom Upstream, auch mitten in
  einer Nachricht, oder ein Fehler der Verbindung selbst (etwa ein Reset). Kann
  der Recorder eine Nachricht des Upstreams nicht lesen, ohne dass die
  Verbindung endet — die Bibliothek kennt den Nachrichtentyp nicht oder kann
  die Nachricht nicht dekodieren —, ist das eine nicht unterstützte
  Serverantwort (`PGR-E6001`, LH-FA-05.a). Eine vorher gelesene `ErrorResponse`
  ändert daran nichts, weder die Klasse noch den Meldungstext: Der nennt die
  nicht lesbare Nachricht, nicht die `ErrorResponse`.

* *Schweregrad.* Er zählt nicht: `FATAL`, `PANIC` und `ERROR` gelten gleich.
* *Stelle.* Die `ErrorResponse` darf an jeder Stelle der Antworten stehen, auch
  nach Ergebnissen einer anderen Anweisung derselben Anfrage oder als erste
  Antwort, wenn der Upstream sie zwischen zwei Interaktionen gesendet hat (sie
  gehört dann zur nächsten, LH-FA-05.a). Eine `ErrorResponse`, auf die ein
  `ReadyForQuery` folgt, gehört zu ihrer abgeschlossenen Interaktion und zählt
  für eine spätere nicht. Sendet der Client nach einer solchen Nachricht keine
  Anfrage mehr, endet die Session wie nach einem `ReadyForQuery` (oben).
* *Diagnose.* Meldung und Fehlerantwort an den Client nennen SQLSTATE (`C`) und
  Meldung (`M`) der letzten `ErrorResponse` des Upstreams vor dem Abbruch.
* *Client.* Von einer einfachen Anfrage erhält der Client keine der Antworten,
  nur die `ErrorResponse` mit `PGR-E6001` (Weitergabe, unten). Bei einer
  Extended-Interaktion bleiben die Server-Nachrichten, die schon weitergegeben
  waren, beim Client, auch die `ErrorResponse` des Upstreams (LH-FA-18.a);
  danach folgt die mit `PGR-E6001`, soweit die Verbindung sie noch annimmt.
* *Aufzeichnung.* Die Session wird nicht übernommen, auch nicht ihre vorherigen
  Interaktionen (LH-FA-05.a) — anders als beim Abbruch ohne Fehlerantwort.
* *Replay.* Eine Aufzeichnung enthält keine solche Interaktion. Eine Interaktion,
  die nicht mit `ready_for_query` endet, macht die Aufzeichnung zu
  einer beschädigten (`PGR-E3003`).

**Weitergabe an den Client.** Im Record-Modus gehen die Serverantworten einer
einfachen Anfrage in der Reihenfolge an den Client, in der der Upstream sie
sendet, auch `NoticeResponse` und `ParameterStatus`. Ist eine davon nicht
unterstützt (LH-FA-05.a), erhält der Client keine von ihnen, sondern die
`ErrorResponse` mit `PGR-E6001`, auch wenn vorher Ergebnisse einer anderen
Anweisung derselben Anfrage kamen. v1 sagt nicht zu, dass eine Antwort vor dem
`ReadyForQuery` ihrer Interaktion beim Client ankommt; bis dahin hält der
Recorder die Antworten im Speicher, auch eine große Ergebnismenge (`SPEC-032`).

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
werten, Default aus; siehe LH-FA-03.b), `--session-assignment` (siehe
LH-FA-12.a) sowie `--tls-cert`, `--tls-key` und `--allow-plaintext` (siehe
LH-FA-23.a).

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
   mit Exit-Code `3`; ein Recording ohne Session mit Interaktion ist nicht
   verwendbar (Exit-Code `3`, `PGR-E3004`), auch eines, dessen Sessions nur
   aufgezeichnete Lebendprüfungen enthalten (LH-FA-09.a).
3. Jede eingehende Client-Verbindung erhält einen eigenen Replay-Cursor und, je
   nach `--session-assignment`, eine Session (LH-FA-12.a).
4. Der Cursor zeigt auf die nächste erwartete Interaktion der zugeordneten
   Session; er beginnt am Anfang der Session. Stellt eine Verbindung eine
   Anfrage, zu der es keine nicht zugeordnete Session mehr gibt, ist das ein
   Replay-Mismatch (`PGR-E5003`); eine Lebendprüfung (LH-FA-09.a) ist dafür
   keine Anfrage.

---

### LH-FA-03.b — Nicht verbrauchte Interaktionen

Wird eine Replay-Session beendet, bevor alle ihr zugeordneten Interaktionen
verbraucht wurden, oder wird eine aufgezeichnete Session nie verbunden, wird dies mindestens als Warnung mit dem Meldungscode `PGR-W2001` protokolliert (`SPEC-034`).

Aufgezeichnete Lebendprüfungen (LH-FA-09.a) zählen dabei nicht zu den
Interaktionen einer Session: Bleiben sie in der Wiedergabe aus, ist das keine
nicht verbrauchte Interaktion, und die Warnung zählt nur die übrigen.

Die Option `--fail-on-unconsumed` (Umgebungsvariable
`PGWIRE_RECORDER_FAIL_ON_UNCONSUMED`) wertet dies als Fehler: Die Verbindung
zählt als fehlerhaft beendet (`PGR-E5002`, Exit-Code `5`, Fehlerebenen:
LH-FA-13.b). Der Default ist aus (`SPEC-012`); Query-Mismatches sind unabhängig
davon immer Fehler.

**Verbraucht.** Eine einfache Interaktion ist verbraucht, sobald ihre
aufgezeichneten Antworten gesendet sind; eine Extended-Interaktion erst mit den
Antworten ihrer letzten Gruppe (LH-FA-18.a). Eine Extended-Interaktion, die vor
dem `Sync` ihrer letzten Gruppe endet, ist nicht verbraucht, auch wenn schon
Gruppen beantwortet sind. Eine Session ohne Interaktion, auch eine nur aus
Lebendprüfungen (LH-FA-12.a), hat nichts zu verbrauchen: Sie ist weder nicht
verbraucht noch nie zugeordnet, auch nicht bei `connection`.

**Zeitpunkte.** Geprüft wird an zwei Stellen, mit und ohne Option gleich:

1. *Ende einer Verbindung*, gleich aus welchem Grund: der Client schließt sie,
   ein Verbindungsfehler beendet sie, das Herunterfahren beendet sie nach der
   laufenden Interaktion, oder die Frist `--shutdown-timeout` beziehungsweise
   ein weiteres Signal beendet sie zwangsweise (LH-FA-13.a). Hat die Verbindung
   eine Session und ist darin eine Interaktion nicht verbraucht, entsteht eine
   Meldung für diese Session. Eine Verbindung ohne zugeordnete Session meldet
   nichts.
2. *Prozessende* nach dem kontrollierten Herunterfahren, wenn alle
   Verbindungen beendet sind: Sessions mit Interaktionen, die nie zugeordnet
   wurden, ergeben zusammen eine Meldung. Nach einem Startfehler (LH-FA-13.b)
   wird nicht geprüft.

**Meldung.** Mit und ohne Option dieselbe Stelle und derselbe Text; die
Option ändert nur Code und Stufe (`PGR-W2001` als `warn`, `PGR-E5002` als
`error` mit dem Kopf aus `SPEC-034`). Die Meldung am Ende einer Verbindung nennt
die Session, die Zahl der nicht verbrauchten und aller Interaktionen der Session
(Lebendprüfungen nicht gezählt) und die aufgezeichnete Nummer (`sequence`) der
ersten nicht verbrauchten. Die Meldung beim Prozessende nennt die Zahl der nie
zugeordneten Sessions und die `id` der ersten.

**Fehlerebene.** `PGR-E5002` ist ein Verbindungsfehler (LH-FA-13.b) und wird
gemerkt, wenn er entsteht; der Prozess läuft weiter. Er wird dem Client nicht
zugestellt: Er entsteht erst, wenn die Verbindung schon endet. Endet eine
Verbindung durch einen anderen Verbindungsfehler (etwa `PGR-E5001`, `PGR-E6001`
oder `PGR-E4006` durch die Frist), wird dieser zuerst gemerkt und danach die
Meldung über die nicht verbrauchten Interaktionen geschrieben; beide stehen im
Log, und der Exit-Code ist der des zuerst gemerkten (LH-FA-13.b). Nie
zugeordnete Sessions werden zuletzt gemerkt und bestimmen den Exit-Code nur,
wenn bis dahin kein Verbindungsfehler auftrat.

**Andere Kommandos.** Die Option gibt es nur bei `replay`. Bei `record` und
`play` ist sie eine unbekannte Option (`PGR-E2001`, Exit-Code `2`); ihre
Umgebungsvariable lassen sie unbeachtet (LH-FA-17.a).

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

**Nachrichten zwischen Interaktionen.** Eine `NoticeResponse` oder ein
`ParameterStatus`, die der Upstream außerhalb einer laufenden Interaktion sendet,
gehören zur nächsten Interaktion der Session und stehen vor deren übrigen
Serverantworten; welcher Gruppe einer Extended-Interaktion sie zugeordnet sind,
regelt LH-FA-18.a. Das Replay gibt sie mit dieser Interaktion wieder; ist die
Interaktion eine Lebendprüfung, gibt es sie nicht wieder (LH-FA-09.a). Was der Upstream nach der letzten Interaktion
einer Session sendet, wird nicht aufgezeichnet. `NotificationResponse` ist nicht
unterstützt (LH-FA-05.e).

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

TLS-Passthrough ohne Protokolleinsicht reicht für Recording nicht aus. Der
Recorder terminiert TLS nur zum Client und nur auf Wunsch (`LH-FA-23.a`); die
Verbindung zum Upstream im Record-Modus ist immer unverschlüsselt.

Ohne diesen Wunsch gilt: Versucht ein Client eine SSL/TLS-Aushandlung
(`SSLRequest`), lehnt der Recorder sie mit dem Einzelbyte `N` ab und erwartet
einen unverschlüsselten PGWire-Verbindungsaufbau auf derselben Verbindung.
Besteht der Client auf TLS und bricht ab, ist das kein Fehler des Recorders.

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
| `SSLRequest` | mit `N` abgelehnt (LH-FA-05.c); mit Zertifikat und Schlüssel angenommen (LH-FA-23.a) |
| `GSSENCRequest` | mit `N` abgelehnt |
| `CancelRequest` | Verbindung wird geschlossen, nicht an den Upstream weitergeleitet; kein Verbindungsfehler, Warnung `PGR-W3001` |
| Verbindung ohne erste Nachricht (etwa eine TCP-Probe) | Verbindung wird geschlossen; kein Verbindungsfehler, keine Warnung |
| erste Nachricht, die keine PGWire-Startnachricht ist: Längenfeld unter 8 oder über `SPEC-045` (etwa eine HTTP-Anfrage) | Verbindung wird ohne Antwort geschlossen; kein Verbindungsfehler, Warnung `PGR-W3003` |
| Startnachricht mit anderer Protokollversion als 3.0 | `ErrorResponse` mit `PGR-E6002` (siehe oben) |
| Sonderanfrage mit unbekanntem Code (Hauptnummer 1234) | nicht unterstützt (`PGR-E6001`) |
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

1. Ist die eingehende Anfrage eine Lebendprüfung (unten) und steht der Cursor
   zwischen zwei Interaktionen, wird sie außerhalb der Reihe beantwortet, und
   die Schritte 2 und 3 entfallen. Sonst wird der eingehende SQL-Text mit dem
   aufgezeichneten SQL-Text der Interaktion am Cursor verglichen:
   `eingehender SQL-Text == aufgezeichneter SQL-Text`.
2. Bei Gleichheit werden die für die Interaktion aufgezeichneten
   Backend-Nachrichten in der gespeicherten Reihenfolge an den Client gesendet
   und der Cursor rückt vor.
3. Bei Ungleichheit liegt ein Mismatch vor (LH-FA-10.a).

Der Vergleich erfolgt byte- beziehungsweise stringgenau nach der
PGWire-Dekodierung. v1 normalisiert SQL nicht. Insbesondere werden nicht
automatisch ignoriert: Whitespace-Unterschiede, Kommentare,
Groß-/Kleinschreibung, Literalwerte und semantisch äquivalente SQL-Varianten.
Die einzige Ausnahme ist die Lebendprüfung.

Mehrfach identische Queries werden über die Position am Cursor zugeordnet; jede
aufgezeichnete Interaktion wird höchstens einmal verbraucht.

**Lebendprüfung.** Connection-Pools und Treiber prüfen eine ruhende Verbindung
mit einer Anfrage ohne Anweisung, etwa `-- ping`.

* *Definition.* Eine Lebendprüfung ist eine `Query`, deren SQL-Text nach den
  lexikalischen Regeln von PostgreSQL nur aus Leerraum, Zeilenkommentaren und
  geschlossenen Blockkommentaren besteht; der leere Text zählt dazu. Leerraum
  sind die Zeichen Leerzeichen, Tabulator (`\t`), Zeilenvorschub (`\n`),
  Wagenrücklauf (`\r`) und Seitenvorschub (`\f`), ab PostgreSQL 17 auch der
  vertikale Tabulator (`\v`); PostgreSQL 14 bis 16 beantworten eine Anfrage mit
  `\v` außerhalb eines Kommentars mit einem Syntaxfehler.
  Ein Zeilenkommentar beginnt mit `--` und reicht bis zum nächsten `\n` oder
  `\r` oder bis zum Ende des Texts. Ein Blockkommentar beginnt mit `/*` und
  endet mit dem dazu passenden `*/`; Blockkommentare sind verschachtelt, jedes
  `/*` darin öffnet eine weitere Ebene, und innerhalb eines Blockkommentars hat
  `--` keine Bedeutung. Keine Lebendprüfung ist ein Text mit einem nicht
  geschlossenen Blockkommentar oder mit einem anderen Zeichen, auch einem
  einzelnen `;` oder einem Leerraum außerhalb von ASCII. Die Klasse ist
  lexikalisch und gilt für jede solche Anfrage, gleich wozu sie gesendet wird.
* *Serverversion.* Maßgeblich ist die Hauptversion im Serverparameter
  `server_version` einer Session, also seine führenden Ziffern: für eine
  aufgezeichnete Lebendprüfung die ihrer Session, für eine eingehende die der
  Session, die der Verbindung zugeordnet ist. Vor der Zuordnung, und wenn
  `server_version` fehlt, nicht mit einer Ziffer beginnt oder keine lesbare Zahl
  ergibt, ist `\v` kein Leerraum; eine Anfrage mit `\v` außerhalb eines
  Kommentars wird dann wie jede andere Anfrage behandelt.
* *Antwort.* Zwischen zwei Interaktionen, auch vor der ersten und nach der
  letzten der Session, beantwortet das Replay eine Lebendprüfung mit
  `EmptyQueryResponse` und `ReadyForQuery`, wie PostgreSQL. Der Cursor bleibt
  stehen, und die Lebendprüfung löst keine Session-Zuordnung aus
  (LH-FA-12.a). Das `ReadyForQuery` trägt den Transaktionsstatus des letzten
  `ReadyForQuery`, das das Replay auf dieser Verbindung gesendet hat, nach dem
  Handshake `I`. Mitten in einer Extended-Interaktion ist eine Lebendprüfung
  eine Abweichung (LH-FA-18.a).
* *Aufgezeichnete Lebendprüfungen.* Das Replay liest eine Aufzeichnung so, als
  stünden die Interaktionen mit einer Lebendprüfung als Anfrage nicht darin: Der
  Cursor überspringt sie, für die Session-Zuordnung (LH-FA-12.a) und die nicht
  verbrauchten Interaktionen (LH-FA-03.b) zählen sie nicht, und ihre
  aufgezeichneten Antworten verwendet es nicht. Eine Aufzeichnung, deren
  Sessions nur Lebendprüfungen enthalten, enthält damit keine Session mit
  Interaktion (`PGR-E3004`, LH-FA-03.a). Diagnosen nennen weiter die
  aufgezeichnete Interaktionsnummer (`sequence`). Record und Einspielen
  behandeln Lebendprüfungen wie jede andere Anfrage.

---

### LH-FA-10.a — Mismatch

Bei einem Mismatch enthält die Diagnose mindestens:

- Session beziehungsweise Verbindungskontext,
- erwartete Interaktionsnummer,
- erwartete Query,
- tatsächlich empfangene Query.

Der Recorder springt nicht zur nächsten Interaktion und führt keine Fuzzy-Suche
durch. Sendet ein Client weitere Anfragen, obwohl keine aufgezeichnete
Interaktion mehr verfügbar ist, wird dies als Replay-Mismatch behandelt; eine
Lebendprüfung beantwortet das Replay auch dann (LH-FA-09.a).

**Fehlermodi:** Replay-Mismatch → Exit-Code `5` (`SPEC-018`).

---

### LH-FA-11.a — Aufzeichnung von Fehlerantworten

`ErrorResponse` wird als Teil der geordneten Serverantworten der Interaktion
aufgezeichnet (siehe LH-FA-05.a) und im Replay in gespeicherter Reihenfolge,
einschließlich des abschließenden `ReadyForQuery`, wiedergegeben.

**Diagnosefelder.** `ErrorResponse` und `NoticeResponse` werden je Feldcode
aufgezeichnet; die Reihenfolge der Felder innerhalb der Nachricht ist nicht Teil
der Aufzeichnung. Für die Feldcodes, die die PGWire-Bibliothek benennt
(`S`, `V`, `C`, `M`, `D`, `H`, `P`, `p`, `q`, `W`, `s`, `t`, `c`, `d`, `n`, `F`,
`L`, `R`), unterscheidet sie ein Feld mit leerem Wert nicht von einem fehlenden,
ebenso ein Zahlenfeld (`P`, `p`, `L`) mit dem Wert `0` oder ohne lesbare Zahl;
solche Felder gelten als fehlend und gehen weder in die Aufzeichnung noch an den
Client, im Record- wie im Replay-Modus. Ein anderer Feldcode wird mit seinem Wert
aufgezeichnet und gesendet, auch wenn der Wert leer ist. Der Client erhält damit
in beiden Modi dieselbe Nachricht.

**Nicht aufzeichenbare Fehlerantwort.** Eine `ErrorResponse`, nach der der
Upstream die Verbindung ohne `ReadyForQuery` beendet, ist nicht verlustfrei
aufzuzeichnen und nicht unterstützt (`PGR-E6001`); die Regeln stehen in
LH-FA-02.b §Fehlerantwort vor dem Abbruch.

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

Eine Lebendprüfung (LH-FA-09.a) ist für die Zuordnung im Replay keine Anfrage:
Sie löst keine Zuordnung aus, auch nicht, wenn keine Session mehr frei ist, und
eine Verbindung, die nur Lebendprüfungen stellt, ist eine Verbindung ohne
Anfrage. Für das Replay ist eine Session, deren Interaktionen alle
Lebendprüfungen sind, eine Session ohne Interaktion: Bei `first-request` wird
sie übersprungen; bei `connection` erhält die n-te Verbindung sie wie jede
Session ohne Interaktion. Das Einspielen behandelt Lebendprüfungen wie jede
andere Anfrage und führt eine solche Session aus.

---

### LH-FA-13.a — Signalbehandlung

Auf `SIGINT` und `SIGTERM` fährt der Prozess kontrolliert herunter. Unter
Windows entspricht der Konsolenabbruch (`Strg+C`, `Strg+Break`) einem `SIGINT`:

- keine neuen Verbindungen annehmen,
- laufende Schreiboperationen soweit möglich abschließen,
- Sessions, die noch laufen, nach Abschluss ihrer laufenden Interaktion
  beenden; eine nicht abgeschlossene Interaktion wird nicht übernommen
  (LH-FA-02.b). Eine Client-Nachricht, die um das Signal herum eintrifft, wird
  entweder noch verarbeitet oder nicht weitergeleitet; aufgezeichnet wird sie
  nur mit ihrer abgeschlossenen Interaktion. Danach beginnt keine
  neue Interaktion: Eine Client-Nachricht, die eine neue begänne (eine `Query`
  oder eine Nachricht nach dem `Sync` einer Extended-Interaktion), wird nicht
  weitergeleitet, und die Session endet regulär nach der laufenden Interaktion;
  diese nimmt bis zu ihrem `Sync` noch ihre Client-Nachrichten an,
- das Recording schreiben (LH-FA-07.a); schlägt das Schreiben fehl, bleibt die
  Zieldatei der letzte vollständig geschriebene Stand.

Beim Beginn des Herunterfahrens schreibt der Prozess eine Zeile der Stufe `info`
mit dem Attribut `sessions`: der Zahl der angenommenen Verbindungen, die zu diesem
Zeitpunkt noch nicht beendet sind, auch `0`. Die Frist `--shutdown-timeout` (`SPEC-046`) zählt ab dem ersten Signal. Ist bei
ihrem Ablauf eine Session nicht beendet, endet sie zwangsweise wie bei einem
Abbruch (LH-FA-02.b): ihre laufende Interaktion wird nicht übernommen, die
abgeschlossenen bleiben, und die Verbindungen zu Client und Server werden
geschlossen. Endet dabei eine Interaktion unvollständig, ist das `PGR-E4006`
(Verbindungsfehler der Klasse 4, LH-FA-13.b); danach wird das Recording
geschrieben. Der Wert `0` schaltet die Frist ab. Ein weiteres Signal lässt die
Frist sofort ablaufen.

**Randfälle des Herunterfahrens.**

* *Ohne offene Verbindung.* Ist beim Signal keine Verbindung offen, endet der
  Prozess ohne zu warten; `record` schreibt das Recording, auch eines ohne
  Sessions (LH-FA-07.a).
* *Startphase.* Ein Signal, bevor der Prozess lauscht, bricht den Start nicht ab:
  Die Startprüfungen laufen zu Ende, und ein Startfehler beendet den Prozess mit
  dem Exit-Code seiner Klasse (LH-FA-13.b). Sonst nimmt der Prozess keine
  Verbindung an und endet wie ohne offene Verbindung.
* *Was die Frist begrenzt.* Die Frist begrenzt das Warten auf die Verbindungen.
  Das Zwangsende und das Schreiben des Recordings danach zählen nicht zu ihr. Ein
  Schreiben des Recordings, das bei ihrem Ablauf läuft, wird nicht abgebrochen;
  eine Verbindung, deren Ende der Prozess vor dem Ablauf bemerkt hat, endet nicht
  zwangsweise.
* *Zwangsende.* Es trifft jede angenommene Verbindung, die bei Ablauf nicht
  beendet ist, auch eine im Aufbau (Startnachricht, Verbindungsaufbau zum
  Upstream). `PGR-E4006` entsteht genau für eine Verbindung mit einer begonnenen,
  nicht abgeschlossenen Interaktion: im Record vor deren `ReadyForQuery`, im
  Replay einer begonnenen, nicht verbrauchten (LH-FA-03.b *Verbraucht*). Eine
  Verbindung ohne eine solche Interaktion, auch eine im Aufbau, endet ohne
  Meldung. Bemerkt der Prozess das Schließen durch das Zwangsende beim Lesen oder
  Schreiben, ist das kein weiteres Verbindungsende (kein `PGR-E4003`).
* *Meldung.* Je Verbindung mit unvollständiger Interaktion eine Meldung
  `PGR-E4006` als Log-Zeile der Stufe `error`. Dem Client wird sie nach LH-FA-13.b
  zugestellt; das Schreiben dauert höchstens `SPEC-051`, auch an einen Client, der
  nicht liest. Sie nennt im Record die `id`, unter der die Session geschrieben
  wird, oder dass sie ohne abgeschlossene Interaktion nicht geschrieben wird, und
  die Nummer (`sequence`), die die verworfene Interaktion getragen hätte; im Replay
  die `id` der zugeordneten Session und die `sequence` der unvollständigen
  Interaktion. Bei mehreren Meldungen wird die erste gemerkt (LH-FA-13.b).
* *Zweites Signal.* Es lässt die Frist sofort ablaufen, auch beim Wert `0`.
  Spätestens wenn der Prozess keine Verbindung mehr annimmt, hat er das erste
  Signal behandelt; zwei Signale, die davor eintreffen, darf das Betriebssystem
  zu einem zusammenfassen.
* *Grenze.* Ein Signal, bevor der Prozess Signale behandelt, beendet ihn nach der
  Voreinstellung des Betriebssystems; dann gilt kein Exit-Code der Spezifikation
  (LH-FA-13.b).

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
nicht zugeordnete Session) beenden nur die betroffene Verbindung. Eine Verbindung,
die der Listener nicht annehmen kann, ist ebenso ein Verbindungsfehler
(`PGR-E4000`); der Prozess nimmt danach weiter an. Entstehen beim Ende einer
Verbindung mehrere Fehler, ist jeder eine eigene Meldung, und der erste wird
gemerkt (`SPEC-034` §Ausgabe). Dem Client wird, wo das Protokoll es erlaubt, eine
`ErrorResponse` mit dem Meldungscode im Meldungstext zugestellt; `PGR-E5002`
entsteht erst beim Ende der Verbindung und wird nicht zugestellt (LH-FA-03.b). Der Prozess
läuft weiter und merkt sich die Klasse des ersten aufgetretenen
Verbindungsfehlers.

**Prozessende.** Sessions, die bis zum kontrollierten Herunterfahren nie
zugeordnet wurden, zählen bei `--fail-on-unconsumed` als Fehler der Klasse 5
(`PGR-E5002`), sonst als Warnung (`PGR-W2001`); sie werden nach allen
Verbindungsfehlern gemerkt (LH-FA-03.b). Nach einem kontrollierten
Herunterfahren ist der Exit-Code der
der gemerkten Klasse, sonst `0`. Schlägt dabei das Schreiben des Recordings
fehl, ist er `3` und hat Vorrang vor der gemerkten Klasse. Beendet der Prozess nicht kontrolliert (zum Beispiel durch
`SIGKILL`), gilt kein Exit-Code der Spezifikation.

---

### LH-FA-14.a — Logging und Diagnose

Logs werden nach `stderr` geschrieben (`SPEC-006`); Nutzdaten beziehungsweise
maschinenlesbare Ausgaben auf `stdout` werden dadurch nicht verunreinigt. Auf
`stdout` stehen nur die Hilfe (LH-FA-01.a), die Ausgabe von `version` und die von
`config show` (LH-FA-17.a); sie schreiben nichts auf `stderr`.

**Log-Level.** Es gibt genau die Stufen `error`, `warn`, `info` und `debug`;
Standard ist `info` (`SPEC-005`). Eine Stufe zeigt ihre Zeilen und die aller
strengeren Stufen. Die Detailstufe wird über `--log-level` (Umgebungsvariable
`PGWIRE_RECORDER_LOG_LEVEL`) gesetzt, bei `record`, `replay` und `play` und wie
jede Option nach dem Kommando. Der Wert ist eine Aufzählung nach LH-FA-17.a: genau
einer der vier Namen in Kleinbuchstaben; jeder andere Wert, auch der leere
(`--log-level=`), ein großgeschriebener (`INFO`) und ein anderer Name (`warning`,
`trace`, `off`), ist `PGR-E2001`; eine leere Umgebungsvariable gilt als nicht
gesetzt. In die Stufen gehen:

* `error`: jeder Verbindungsfehler (LH-FA-13.b) als Log-Zeile,
* `warn`: jede Warnung; jede Zeile dieser Stufe trägt einen Meldungscode
  (`PGR-W…`), ein Hinweis ohne Maßnahme für den Anwender steht unter `info`
  oder `debug`,
* `info`: Start und Ende von `record` und `replay` und der Beginn des
  Herunterfahrens (LH-FA-13.a), keine Zeile je Verbindung oder Interaktion,
* `debug`: Ereignisse je Verbindung; welche, ist nicht Vertrag.

Ein Startfehler (LH-FA-13.b) und ein Fehler beim Schreiben des Recordings am Ende
erscheinen als Zeile beim Prozessende auf jeder Stufe, auch wenn er vor dem Lesen
von `--log-level` entsteht oder dessen Wert betrifft; vorher schreibt der Prozess
nichts anderes.

**Zeilenform.** Eine Log-Zeile ist eine Zeile im Format `logfmt`: zuerst `time`
(RFC 3339 mit Millisekunden und Zonenversatz), `level` (`DEBUG`, `INFO`, `WARN`,
`ERROR`) und `msg`, danach die Attribute; ein Wert mit Leerraum, `=`, `"` oder
einem Steuerzeichen steht in Anführungszeichen mit Escapes. Ein Verbindungsfehler
trägt die Attribute `code` und `error` (Fehlertext, `SPEC-034`), eine Warnung das
Attribut `code`. Vertrag sind die Schlüssel `level`, `code` und `error`, an der
Zeile beim Beginn des Herunterfahrens auch `sessions` (LH-FA-13.a); der Text
von `msg`, weitere Attribute und deren Reihenfolge sind es nicht. Das Attribut
`error` steht nur an einer Zeile der Stufe `error` und trägt immer einen Fehlertext
mit Kopf; der Text einer Bibliothek an einer Zeile einer anderen Stufe steht unter
dem Schlüssel `grund`. Die Zeile beim
Prozessende ist keine Log-Zeile: Sie ist genau der Fehlertext.

Passwörter aus Verbindungsdaten werden nicht absichtlich in Logs ausgegeben.
Parameterwerte erscheinen auf keiner Stufe (`SPEC-033`).

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
Werte lauten `true` oder `false`. Eine boolesche Option ohne Wert ist `true`; mit
`=` nimmt sie `true` oder `false`, jeder andere Wert, auch der leere und `1`, ist
`PGR-E2001`. Für die Umgebungsvariable gilt dieselbe Wertemenge; eine leere
Umgebungsvariable gilt als nicht gesetzt. Nennt die Kommandozeile eine Option
mehrfach, gilt die letzte Angabe. Die Umgebungsvariable einer Option, die das
Kommando nicht kennt, bleibt unbeachtet. Geprüft wird jeder gesetzte Wert einer
Option des Kommandos, unabhängig von der Priorität: Eine gesetzte
Umgebungsvariable mit ungültigem Wert ist `PGR-E2001`, auch wenn die Kommandozeile
dieselbe Option setzt und damit vorgeht. Fordert der Aufruf die Hilfe an, wird
keine Quelle gelesen oder geprüft (LH-FA-01.a). Das **Passwort** hat eine eigene Regel: Es kommt
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

**Dauer.** Der Wert von `--shutdown-timeout` ist `0` oder eine ganze Zahl ohne
Vorzeichen mit genau einer Einheit `ms`, `s` oder `m` in Kleinbuchstaben; `0` mit
Einheit ist ebenfalls `0`. Jeder andere Wert ist ungültig, auch der leere, ein
negativer, einer ohne Einheit außer `0`, einer mit Nachkommastellen, Leerraum,
großgeschriebener oder zusammengesetzter Einheit (`1m30s`) und einer, der länger
ist als die längste Dauer, die die Implementierung darstellt: als Option oder
Umgebungsvariable `PGR-E2001`, als Schlüssel der Konfigurationsdatei `PGR-E2004`.

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
| `--shutdown-timeout` | `record`, `replay` | `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT` | `5s` (`SPEC-046`; Dauer mit Einheit `ms`, `s` oder `m`, `0` ohne Frist) |
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
| `--tls-cert` | `record`, `replay` | `PGWIRE_RECORDER_TLS_CERT` | — (Pfad einer PEM-Datei; verlangt `--tls-key`) |
| `--tls-key` | `record`, `replay` | `PGWIRE_RECORDER_TLS_KEY` | — (Pfad einer PEM-Datei; verlangt `--tls-cert`) |
| `--allow-plaintext` | `record`, `replay` | `PGWIRE_RECORDER_ALLOW_PLAINTEXT` | `false` (verlangt `--tls-cert`) |
| `--upstream-ca` | `play` | `PGWIRE_RECORDER_UPSTREAM_CA` | — (Pfad einer PEM-Datei; verlangt TLS zum Server) |
| `--compare-responses` | `play` | `PGWIRE_RECORDER_COMPARE_RESPONSES` | `false` (`LH-FA-24.a`) |
| — (nur Umgebung) | `play` | `PGWIRE_RECORDER_PASSWORD` | — |
| `--config` | `record`, `replay`, `play`, `config show` | `PGWIRE_RECORDER_CONFIG` | `.pgwire-recorder.yaml` im aktuellen Verzeichnis, sofern vorhanden |
| `--log-level` | `record`, `replay`, `play` | `PGWIRE_RECORDER_LOG_LEVEL` | `info` (`SPEC-005`) |

---

### LH-FA-18.a — Extended Query: Ablauf, Aufzeichnung, Matching

**Nachrichten.** Client: `Parse`, `Bind`, `Describe`, `Execute`, `Close`,
`Flush`, `Sync`. Server: `ParseComplete`, `BindComplete`, `CloseComplete`,
`ParameterDescription`, `RowDescription`, `NoData`, `DataRow`,
`CommandComplete`, `EmptyQueryResponse`, `PortalSuspended`, `ErrorResponse`,
`NoticeResponse`, `ParameterStatus` (etwa nach einem `SET`), `ReadyForQuery`.

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
der Session bleiben erhalten. Beim kontrollierten Herunterfahren endet eine
Session erst nach dem `ReadyForQuery` ihrer laufenden Extended-Interaktion, auch
wenn deren `Sync` noch aussteht (LH-FA-13.a). Endet die Client-Verbindung,
während der Recorder an den Server sendet oder auf ihn wartet, endet die Session
ohne weiteres Warten, sobald der Recorder das Ende beim Lesen vom oder Schreiben
zum Client bemerkt, auch beim Herunterfahren. Endet eine Session, schließt der
Recorder die Client-Verbindung; ein blockiertes Schreiben an den Client endet
damit. Eine Fehlerantwort an den Client schreibt er davor höchstens `SPEC-051` lang,
auch an einen Client, der nicht liest.

**Record.** Der Recorder leitet alle Nachrichten unverändert und in
Ankunftsreihenfolge weiter, auch wenn der Client mehrere Nachrichten sendet,
ohne auf Antworten zu warten (Pipelining). Er vermittelt beide Richtungen
unabhängig: Das Weiterleiten an den Server wartet nicht auf das Lesen von
Antworten, das Weiterleiten an den Client nicht auf weitere Client-Nachrichten;
ein Gegenüber, das nicht liest, hält nur die Richtung zu ihm an. Zeichnet er die Interaktion auf,
gruppiert er sie (`SPEC-041`); Zeitangaben enthält die Aufzeichnung nur mit
`--record-timing` (LH-FA-21.a).

Nicht unterstützt (`PGR-E6001`, die Session wird nicht übernommen, LH-FA-05.a)
sind im Record:

* eine `Query` während einer laufenden Extended-Interaktion, also nach deren
  erster Nachricht und vor deren `Sync`,
* ein `Describe` oder `Close` mit einer anderen Zielart als Statement (`S`) oder
  Portal (`P`),
* eine Interaktion, einfach oder Extended, deren Nachrichten nicht die Form nach
  `SPEC-002` beziehungsweise `SPEC-041` ergeben, etwa ein `ReadyForQuery` ohne
  `Sync` oder eine Server-Nachricht des Extended Query Protocol in der Antwort auf
  eine `Query`.

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
  nächsten Client-Nachricht eintreffen; eine danach eintreffende gehört zur
  folgenden Gruppe. Wartet der Client nach dem `Flush` nicht
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

Eine Client-Nachricht der falschen Protokollart am Cursor ist ein Mismatch
(`PGR-E5001`): eine Extended-Nachricht, wo eine einfache Anfrage erwartet wird,
und eine `Query`, wo eine Extended-Interaktion erwartet wird. Dazu gehört jede
`Query` mitten in einer Extended-Interaktion, also nach deren erster Nachricht
und vor deren `Sync`, auch eine Lebendprüfung. Ausgenommen ist nur eine
Lebendprüfung zwischen zwei Interaktionen (LH-FA-09.a), auch wenn am Cursor eine
Extended-Interaktion steht. Ein `Parse` mit leerem SQL-Text ist keine
Lebendprüfung und wird verglichen wie jede andere Nachricht.

**Fehler.** Nach einer `ErrorResponse` verwirft ein Server Nachrichten bis zum
nächsten `Sync`. Die Aufzeichnung enthält die empfangenen Client-Nachrichten
und die Server-Nachrichten der Gruppe; Replay reproduziert sie, ohne eigene
Fehlerlogik.

**Mismatch.** Die Diagnose nach LH-FA-10.a nennt zusätzlich den Index der
Gruppe und Nachricht in der Interaktion sowie erwarteten und empfangenen Nachrichtentyp.
Weichen Parameterwerte ab, nennt sie die Nummer des ersten abweichenden Parameters,
gezählt ab 1 wie `$1`; weicht die Zahl der Parameter ab, nennt sie stattdessen beide
Anzahlen, die erwartete und die empfangene. Parameterwerte erscheinen nie in der
Diagnose, weder im Log noch in der `ErrorResponse` an den Client, auf keinem
Log-Level, auch nicht bei `debug` (`SPEC-033`).

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
| `--allow-recorded-errors` | ein Serverfehler in einer Interaktion, deren Aufzeichnung ebenfalls eine `ErrorResponse` enthält, gilt als erwartet; mit `--compare-responses` ohne Wirkung, ohne dass die Kombination ein Fehler ist |
| `--upstream-tls` | die Verbindung zum Server wird mit TLS aufgebaut |
| `--upstream-ca` | eine zusätzliche Zertifizierungsstelle für die Prüfung des Serverzertifikats (nur mit TLS) |
| `--compare-responses` | die Serverantworten werden mit der Aufzeichnung verglichen (`LH-FA-24.a`) |
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
   Systems und die Zertifikate aus `--upstream-ca` geprüft, und der Host der
   Verbindung muss zum Zertifikat passen. Ein Überspringen der Prüfung gibt es nicht.
   Der Recorder authentifiziert sich als Client mit Klartext-Passwort, MD5 oder
   SCRAM-SHA-256.

   `--upstream-ca` nennt eine PEM-Datei mit einem oder mehreren Zertifikaten; sie
   ergänzt den Zertifikatsspeicher des Systems und ersetzt ihn nicht. Ohne TLS ist
   die Option eine ungültige Verwendung (`PGR-E2001`). Die Datei wird beim Start
   gelesen, vor der ersten Verbindung; eine Datei, die nicht lesbar ist oder kein
   gültiges PEM-Zertifikat enthält, ist ein Startfehler `PGR-E2007` (Exit-Code `2`).
   Der Ablauf eines Zertifikats der Datei oder des Servers wird nicht beim Start,
   sondern bei jedem Verbindungsaufbau geprüft und ist dort `PGR-E4005`. Ein relativer
   Pfad gilt ab dem aktuellen Verzeichnis, auch wenn er in der Konfigurationsdatei
   steht.
3. Die Client-Nachrichten der Interaktionen werden in der aufgezeichneten
   Reihenfolge gesendet: bei einer einfachen Anfrage die `Query`, bei einer
   Extended-Interaktion die Nachrichten jeder Gruppe (`SPEC-041`). Nach jeder
   Interaktion wartet der Recorder auf das `ReadyForQuery`, bei einer aufgezeichneten
   Interaktion ohne `ReadyForQuery` auf das aufgezeichnete Verbindungsende.
4. Die Serverantworten werden gelesen und verworfen, mit Ausnahme einer
   `ErrorResponse`. Ein Vergleich mit den aufgezeichneten Antworten findet nur mit
   `--compare-responses` statt (`LH-FA-24.a`); die Antworten werden dann bis zum
   Ende der Interaktion gehalten.
5. **Fehler und Abweichungen** folgen der Tabelle unten.
6. **Fortsetzung und Abbruch.** Läuft das Einspielen nach einer Fehlerantwort oder
   einer Abweichung weiter (`--continue-on-error`, erwarteter Fehler), sendet der
   Recorder den Rest der Interaktion wie aufgezeichnet und macht mit der nächsten
   Interaktion weiter; bei einer Extended-Interaktion verwirft der Server nach einem
   Fehler bis zum `Sync`. Bricht das Einspielen ab, sendet der Recorder keine weiteren
   Nachrichten der Interaktion (auch keine weiteren Gruppen) und schließt die
   Verbindung, soweit möglich, mit `Terminate`. Mit Vergleich wird eine Abweichung
   am Ende der Interaktion festgestellt (Schritt 4); der Abbruch erfolgt dann vor der
   nächsten Interaktion.
7. **Abbruchsignal:** Bei `SIGINT` oder `SIGTERM` endet das Einspielen nach der
   laufenden Interaktion; eine Wartezeit (`LH-FA-21.a`) wird abgebrochen. Mit
   `--finish-session-on-interrupt` zuvor nach der laufenden Session. Ein zweites
   Signal beendet den Prozess sofort. Die Verbindung wird mit `Terminate`
   geschlossen. Der Exit-Code folgt der Tabelle.

**Fehlerregeln beim Einspielen.** Die Tabelle und der Absatz darunter legen Code,
Exit-Code und Abbruch je Ursache fest; die Schritte beschreiben den Ablauf und
wiederholen sie nicht.

| Ursache | Ohne `--compare-responses` | Mit `--compare-responses` | Abbruch |
|---|---|---|---|
| Server nicht erreichbar; `FATAL` beim Aufbau, der nicht die Anmeldung betrifft (zum Beispiel fehlende Datenbank) | `PGR-E4002`, Exit-Code `4` | gleich | immer sofort |
| Anmeldung fehlgeschlagen (SQLSTATE-Klasse 28), Verfahren nicht unterstützt, TLS abgelehnt (der Server antwortet mit `N`), Zertifikat ungültig oder abgelaufen, Server lehnt unverschlüsselte Verbindung ab | `PGR-E4005`, Exit-Code `4` | gleich | immer sofort |
| Verbindung endet nach dem ersten `ReadyForQuery`, auch durch einen `FATAL` | `PGR-E4003`, Exit-Code `4`, immer | **Aufgezeichnetes Ende, das der Server genauso verursacht** (gleicher SQLSTATE; ohne aufgezeichnete `ErrorResponse` ohne SQLSTATE-Vergleich): erwartet, die Session endet, das Einspielen macht mit der nächsten Session weiter. **`FATAL` ohne Vorbild oder mit anderem SQLSTATE, oder ein `ReadyForQuery`, wo die Aufzeichnung das Ende zeigt**: Abweichung `PGR-E5004`, Exit-Code `5`; die Verbindung ist beendet, die Session endet, mit `--continue-on-error` läuft die nächste Session. **Verbindungsverlust ohne Fehlerantwort des Servers** (zum Beispiel Netzwerkabbruch): `PGR-E4003`, Exit-Code `4` | `PGR-E4003`: immer sofort; erwartetes Ende: nein; `PGR-E5004`: ohne `--continue-on-error` sofort |
| `ErrorResponse` des Servers auf eine Anfrage | `PGR-E4004`, Exit-Code `4` am Ende; zählt nicht, wenn `--allow-recorded-errors` gesetzt ist und die aufgezeichnete Interaktion eine `ErrorResponse` enthält | der Vergleich entscheidet: gleicher SQLSTATE an der aufgezeichneten Stelle gilt als erwartet, jeder andere Fehler ist eine Abweichung; `PGR-E4004` entfällt | ohne `--continue-on-error` sofort, außer bei erwartetem Fehler |
| Abweichung der Antwort (`LH-FA-24.a`) | kein Vergleich | `PGR-E5004`, Exit-Code `5` am Ende; je Interaktion die erste Abweichung | ohne `--continue-on-error` sofort |
| Abbruchsignal (Schritt 7) | Exit-Code `0`, wenn bis dahin kein Fehler auftrat, sonst `4` | `0`, nach einer Abweichung `5` | nach der laufenden Interaktion (oder Session) |

Ein Verbindungsverlust ohne Fehlerantwort (`PGR-E4002`, `PGR-E4003`, `PGR-E4005`)
endet immer mit Exit-Code `4`, auch nach einer früheren Abweichung. Wo die
Aufzeichnung ein Verbindungsende zeigt, wartet der Recorder wie bei jeder Antwort
ohne eigene Frist darauf. Ohne Abbruch ist der Exit-Code am Ende `0`, `4` (ohne
Vergleich nach einer Fehlerantwort) oder `5` (mit Vergleich nach einer Abweichung).

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

### LH-FA-23.a — TLS zum Client

Gilt für `record` und `replay`. Mit `--tls-cert` und `--tls-key` (PEM-Dateien) nimmt
der Recorder TLS von Clients an; die Optionen stehen nur zusammen, jede allein ist
eine ungültige Verwendung (`PGR-E2001`). Beim Start werden Zertifikat und Schlüssel
geladen und geprüft; eine nicht lesbare Datei, ein Inhalt ohne gültiges PEM, ein
Zertifikat, das abgelaufen oder noch nicht gültig ist, ein nicht zum Zertifikat
passender Schlüssel und ein verschlüsselter Schlüssel sind Startfehler `PGR-E2007`
(Exit-Code `2`), bevor ein Listen-Port geöffnet wird. Die Zertifikatskette steht in
der Zertifikatsdatei. Die Namen des Zertifikats prüft der Recorder nicht; das ist
Sache des Clients. Ein relativer Pfad gilt ab dem aktuellen Verzeichnis, auch wenn
er in der Konfigurationsdatei steht. Zertifikat und Schlüssel werden nur beim Start
geladen; ein Ablauf danach ändert laufende und neue Verbindungen nicht.

**Aushandlung.** Auf einen `SSLRequest` antwortet der Recorder mit `S` und führt den
TLS-Handshake aus, mindestens TLS 1.2; danach folgt der PGWire-Verbindungsaufbau auf
der verschlüsselten Verbindung. Scheitert der Handshake, endet die Verbindung mit der
Warnung `PGR-W3002`; der Prozessausgang bleibt unverändert. `GSSENCRequest` wird wie
bisher mit `N` abgelehnt.

**Unverschlüsselte Clients.** Ist TLS konfiguriert, beantwortet der Recorder einen
Verbindungsaufbau mit einer `StartupMessage` ohne vorherigen `SSLRequest` (auch nach
einem mit `N` abgelehnten `GSSENCRequest`) mit einer `ErrorResponse` und beendet die
Verbindung (`PGR-E6003`, Klasse „nicht unterstützt", ein Verbindungsfehler mit
Exit-Code `6` nach `LH-FA-13.b`). Ein `CancelRequest` trägt keine Anfrage- und
keine Antwortdaten und gilt wie sonst (`LH-FA-05.e`, `PGR-W3001`), auch ohne
Verschlüsselung. Mit `--allow-plaintext` (nur zusammen mit `--tls-cert`, sonst
`PGR-E2001`) bleibt die unverschlüsselte Verbindung zulässig; der Recorder bietet
dann beides an. Ohne TLS-Konfiguration gilt `LH-FA-05.c`.

**Wirkung auf die Aufzeichnung.** Die Verschlüsselung der Client-Verbindung steht
nicht in der Aufzeichnung. Eine Aufzeichnung, die über eine verschlüsselte
Verbindung entstand, ist von einer unverschlüsselt entstandenen nicht zu
unterscheiden, und `replay` liefert sie über beide Verbindungsarten gleich aus.

---

### LH-FA-24.a — Vergleich der Antworten

Gilt für `play` mit `--compare-responses`. Nach jeder Interaktion vergleicht der
Recorder die Serverantworten mit den aufgezeichneten. Bei einer Extended-Interaktion
ist das die Folge aller Server-Nachrichten ihrer Gruppen, ohne die Gruppengrenzen,
weil die Zuordnung zu einer Gruppe ohne `Sync` zeitabhängig ist (`LH-FA-18.a`). Beide
Seiten werden zuvor normalisiert: Nachrichten vom Typ `data_row`, `notice_response`
und `parameter_status` entfallen. Die übrigen Nachrichten werden in ihrer Reihenfolge
Stück für Stück verglichen:

| Nachricht | Verglichen wird |
|---|---|
| jede | der Nachrichtentyp; eine zusätzliche oder fehlende Nachricht ist eine Abweichung |
| `row_description` | die Anzahl der Spalten und je Spalte Name und Typ-OID |
| `parameter_description` | die Liste der Typ-OIDs |
| `command_complete` | der Tag ohne die abschließenden Zahlen (`INSERT 0 1` gilt als `INSERT`, `CREATE TABLE` bleibt `CREATE TABLE`), also der Befehl, nicht die Zeilenzahl |
| `error_response` | der SQLSTATE (Feld `C`), nicht der Text |
| `ready_for_query` | der Transaktionsstatus |

Zeilenwerte und Zeilenzahlen werden nicht verglichen. Eine `ErrorResponse` des
Servers in einer Extended-Interaktion ändert den Vergleich nicht: die Folge ab dem
Fehler bis zum `ReadyForQuery` wird wie aufgezeichnet verglichen.

**Unvollständige Aufzeichnung.** Endet eine aufgezeichnete Interaktion nicht mit
`ready_for_query` (Verbindungsende, `FATAL`) oder trägt sie keine Server-Nachrichten,
wird nach der Normalisierung verglichen, soweit die Aufzeichnung reicht: Die
aufgezeichneten Nachrichten müssen übereinstimmen; was der Server danach sendet
(zum Beispiel das `ReadyForQuery` nach einem aufgezeichneten `FATAL`), wird nicht
verglichen und ist keine Abweichung, mit einer Ausnahme: Eine `ErrorResponse` ohne
Vorbild in der Aufzeichnung ist stets eine Abweichung, auch nach dem aufgezeichneten
Ende und bei einer Interaktion ohne aufgezeichnete Server-Nachrichten. Ein
aufgezeichnetes Verbindungsende gilt als erwartet, wenn der Server es genauso
verursacht: mit gleichem SQLSTATE, bei einem Ende ohne aufgezeichnete
`ErrorResponse` ohne SQLSTATE-Vergleich (Fehlerregeln beim Einspielen).

**Meldung.** Die Abweichung wird als `PGR-E5004` gemeldet (Exit-Code `5`) und nennt
Session, Sequenznummer, die Position der Nachricht in der normalisierten Folge und
die Art: Nachrichtenart, Spaltenbeschreibung, Befehl, Fehler oder Transaktionsstatus.
Werte des Servers erscheinen in der Meldung nicht, nur Typen und Namen von Spalten.
Je Interaktion wird die erste Abweichung gemeldet. Bei jedem Ende des Laufs (auch
Abbruch und Signal) schreibt der Recorder eine Log-Zeile der Stufe `info` auf
`stderr` (`SPEC-006`) mit den Attributen `eingespielt`, `verglichen` und
`abweichend`, ohne Meldungscode. `verglichen` zählt jede Interaktion, deren Antworten
mit Vergleich geprüft wurden, auch eine unvollständig aufgezeichnete. Ohne
`--compare-responses` gibt es keine Vergleiche, `verglichen` und `abweichend` sind
`0`. Die Zeile steht auf der Stufe `info` und erscheint daher nicht bei
`--log-level warn` oder `error`; sie erscheint auch, wenn der Lauf endet, bevor eine
Interaktion eingespielt wurde (alle Zähler `0`). Die Folgen einer Abweichung
beschreibt `LH-FA-20.a` (Fehlerregeln beim Einspielen). Ohne `--compare-responses`
findet kein Vergleich statt.

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

`version` ist eine einzelne ganze Zahl und zählt Änderungen des Formats. Ein
Leser lehnt jedes Feld ab, das er für seine Version nicht kennt (`PGR-E3003`);
darum erhöht auch ein neues optionales Feld `version`.
Version 1 umfasst einfache und Extended-Interaktionen: Bei einer einfachen steht
`type: query` in `request` (`SPEC-002`), bei einer Extended-Interaktion steht
`type: extended` an der Interaktion selbst (`SPEC-041`); eine einfache trägt dort
kein `type`. Ein anderer Wert von `type` an einer dieser Stellen, auch ein leerer
oder `null`, macht das Recording zu einem beschädigten (`PGR-E3003`).

Ein Leser liest jede `version` von 1 bis zu seiner eigenen; jede andere lehnt er ab
(`PGR-E3002`). Ein Recording ohne oder mit abweichender `format`-Kennung ist
beschädigt (`PGR-E3003`). Beides endet mit Exit-Code `3` (`SPEC-016`).

Anker, Aliase und Merge-Keys von YAML (`&`, `*`, `<<`) sind nicht Teil des
Formats; ein Recording, das einen davon enthält, ist beschädigt (`PGR-E3003`).
Der Recorder schreibt keine.

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
`notice_response`, `parameter_status`, `ready_for_query`).

Jede Client-Nachricht trägt neben `type` genau die Felder ihres Typs, auch wenn sie
leer sind (`portal: ""`, `param_types: []`, `max_rows: 0`):

| Typ | Felder |
|---|---|
| `parse` | `statement`, `sql`, `param_types` |
| `bind` | `portal`, `statement`, `param_formats`, `params`, `result_formats` |
| `describe`, `close` | `target` (`statement` oder `portal`), `name` |
| `execute` | `portal`, `max_rows` |
| `flush`, `sync` | — |

Ein fehlendes Feld des Typs, ein Feld mit dem Wert `null`, ein anderes Feld oder
ein anderer Typ macht das Recording zu einem beschädigten (`PGR-E3003`).
`parameter_description` trägt die Typ-OIDs der Parameter als `param_types`, auch
leer (`param_types: []`); fehlt das Feld dort oder steht es an einer anderen
Server-Nachricht, auch mit dem Wert `null`, ist das Recording beschädigt. Eine
Interaktion mit `request`, `responses` oder `groups` mit dem Wert `null` ist
beschädigt. Jede Gruppe trägt `client` und `server`; fehlt einer der beiden oder
steht er mit `null`, ist das Recording beschädigt. `server: []` ist zulässig (eine
`Flush`-Gruppe, auf die der Server nichts gesendet hat). Ein NULL-Parameter steht
als `null: true`; ein Leser nimmt den Schlüssel `null` gequotet (`"null": true`)
und ungequotet an. Binärwerte folgen `SPEC-003`.

Jede Interaktion, einfach oder Extended, kann das Feld `offset_ms` tragen (ganze
Zahl ≥ 0, `LH-FA-21.a`); es steht neben `sequence`, bei einer Extended-Interaktion
auch neben `type`. Simple-Interaktionen behalten `request`/`responses`
(`SPEC-002`). Formatkennung und Version (`SPEC-001`) gelten für beide Arten;
Binärdaten folgen `SPEC-003`.

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
| `SPEC-011` | Replay-Matching | strict sequential (einziges Verfahren in v1); Lebendprüfungen außerhalb der Reihe (LH-FA-09.a) | Determinismus (LH-FA-09, LH-QA-01); eine Lebendprüfung ändert den Zustand der Session nicht |
| `SPEC-012` | `--fail-on-unconsumed` | `false` | nicht verbrauchte Interaktionen sind standardmäßig eine Warnung; Query-Mismatches bleiben immer Fehler (LH-FA-10) |
| `SPEC-045` | Höchstlänge der ersten Client-Nachricht | 10000 Bytes | eine längere erste Nachricht ist keine PGWire-Startnachricht (LH-FA-05.e) |
| `SPEC-046` | `--shutdown-timeout` | `5s` | liegt unter der üblichen Stopp-Frist von Containern (10 s vor `SIGKILL`), sodass das Recording geschrieben wird (LH-FA-13) |
| `SPEC-051` | Schreibfrist der Fehlerantwort beim Ende einer Session | `1s`, nicht einstellbar | reicht für eine kurze Nachricht an einen lesenden Client; greift nur bei vollem Puffer zum Client, wo auch eine längere nicht hilft; Sessions enden parallel, die Frist addiert sich nicht und bleibt im Abstand zwischen `SPEC-046` und der Stopp-Frist von Containern (LH-FA-13) |

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
| `SPEC-017` | 4 | Netzwerk-/Upstream-Fehler, beim Einspielen auch Anmeldefehler (`PGR-E4005`) und eine Fehlerantwort des Servers (`PGR-E4004`, ohne `--compare-responses`) | Startfehler: Prozess endet; Verbindungsfehler: Verbindung endet, Prozessende nach LH-FA-13.b; Einspielen: Abbruch beziehungsweise Exit-Code 4 am Ende (LH-FA-20.a) |
| `SPEC-018` | 5 | Replay-Mismatch (`PGR-E5001`), Anfrage ohne nicht zugeordnete Session (`PGR-E5003`), bei `--fail-on-unconsumed` auch nicht verbrauchte Interaktionen oder Sessions (`PGR-E5002`), beim Einspielen mit `--compare-responses` eine Abweichung der Serverantwort (`PGR-E5004`) | Diagnose nach LH-FA-10.a, Verbindung endet, Prozessende nach LH-FA-13.b; Einspielen: Abbruch beziehungsweise Exit-Code 5 am Ende (LH-FA-20.a) |
| `SPEC-019` | 6 | nicht unterstützte PGWire-Funktion, auch eine unverschlüsselte Verbindung, die nicht zugelassen ist (`PGR-E6003`) | Diagnose, Verbindung endet, Prozessende nach LH-FA-13.b |

**Fehlerklassen.** Mindestens folgende Klassen werden unterschieden und
erzeugen einen für Entwickler verständlichen Text:

| ID | Fehlerklasse | Exit Code | Meldungscode |
|---|---|---|---|
| `SPEC-020` | ungültige CLI-Verwendung oder Konfiguration | 2 | `PGR-E2001`, `PGR-E2004`, `PGR-E2005`, `PGR-E2006`, `PGR-E2007` |
| `SPEC-021` | Listen-Port kann nicht geöffnet werden | 4 | `PGR-E4001` |
| `SPEC-022` | Upstream nicht erreichbar oder lehnt ab | 4 | `PGR-E4002`, `PGR-E4004` (ohne `--compare-responses`), `PGR-E4005` |
| `SPEC-023` | Recording kann nicht gelesen/geschrieben werden | 3 | `PGR-E3001` |
| `SPEC-024` | unbekannte Recording-Version | 3 | `PGR-E3002` |
| `SPEC-025` | beschädigtes Recording | 3 | `PGR-E3003` |
| `SPEC-026` | nicht unterstützte PGWire-Nachricht | 6 | `PGR-E6001` |
| `SPEC-027` | Replay-Mismatch, Abweichung beim Einspielen | 5 | `PGR-E5001`, `PGR-E5003`, `PGR-E5004` |
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
| `PGR-E2007` | Konfiguration (Exit 2) | Zertifikat, Schlüssel oder Zertifizierungsstelle nicht lesbar oder kein gültiges PEM; eigenes Zertifikat abgelaufen oder noch nicht gültig, Schlüssel verschlüsselt oder nicht zum Zertifikat passend (LH-FA-20.a, LH-FA-23.a) |
| `PGR-E3000` | Recording (Exit 3) | Rückfall |
| `PGR-E3001` | Recording (Exit 3) | Recording nicht lesbar oder nicht schreibbar (`SPEC-023`) |
| `PGR-E3002` | Recording (Exit 3) | unbekannte Recording-Version (`SPEC-024`) |
| `PGR-E3003` | Recording (Exit 3) | beschädigtes Recording (`SPEC-025`) |
| `PGR-E3004` | Recording (Exit 3) | Recording ohne Session, die sich im Replay zuordnen lässt (LH-FA-03.a) |
| `PGR-E4000` | Netzwerk (Exit 4) | Rückfall |
| `PGR-E4001` | Netzwerk (Exit 4) | Listen-Port nicht zu öffnen (`SPEC-021`) |
| `PGR-E4002` | Netzwerk (Exit 4) | Upstream nicht erreichbar (`SPEC-022`) |
| `PGR-E4003` | Netzwerk (Exit 4) | unerwartetes Verbindungsende (`SPEC-028`) |
| `PGR-E4004` | Netzwerk (Exit 4) | Server beantwortet eine eingespielte Anfrage mit einem Fehler, nur ohne `--compare-responses` (LH-FA-20.a) |
| `PGR-E4005` | Netzwerk (Exit 4) | Anmeldung am Server fehlgeschlagen, Verfahren nicht unterstützt, TLS abgelehnt, ungültiges oder abgelaufenes Zertifikat des Servers oder der Zertifizierungsstelle, oder der Server lehnt unverschlüsselte Verbindungen ab (LH-FA-20.a) |
| `PGR-E4006` | Netzwerk (Exit 4) | Frist beim Herunterfahren abgelaufen, eine Interaktion endete unvollständig (LH-FA-13.a) |
| `PGR-E5000` | Replay (Exit 5) | Rückfall |
| `PGR-E5001` | Replay (Exit 5) | Replay-Mismatch (`SPEC-027`) |
| `PGR-E5002` | Replay (Exit 5) | nicht verbrauchte Interaktionen oder Sessions bei `--fail-on-unconsumed` (LH-FA-03.b) |
| `PGR-E5003` | Replay (Exit 5) | Anfrage einer Verbindung, zu der keine nicht zugeordnete Session mehr existiert (LH-FA-12.a, `SPEC-027`) |
| `PGR-E5004` | Replay (Exit 5) | Serverantwort beim Einspielen weicht von der Aufzeichnung ab, auch ein Serverfehler ohne Vorbild in der Aufzeichnung, ebenso nach deren Ende (`SPEC-027`, LH-FA-24.a) |
| `PGR-E6000` | nicht unterstützt (Exit 6) | Rückfall |
| `PGR-E6001` | nicht unterstützt (Exit 6) | nicht unterstützte PGWire-Nachricht (`SPEC-026`) |
| `PGR-E6002` | nicht unterstützt (Exit 6) | nicht unterstützte PGWire-Protokollversion (LH-FA-05.e) |
| `PGR-E6003` | nicht unterstützt (Exit 6) | unverschlüsselte Verbindung, obwohl TLS konfiguriert und `--allow-plaintext` nicht gesetzt ist (LH-FA-23.a) |
| `PGR-W2001` | Replay | Sitzung endet vor Verbrauch aller Interaktionen, oder Sessions nie zugeordnet (LH-FA-03.b) |
| `PGR-W3001` | Protokollrand | `CancelRequest` empfangen und nicht weitergeleitet (LH-FA-05.e) |
| `PGR-W3002` | Protokollrand | TLS-Aushandlung eines Clients gescheitert (LH-FA-23.a) |
| `PGR-W3003` | Protokollrand | erste Nachricht einer Verbindung ist keine PGWire-Startnachricht (LH-FA-05.e) |

**Ausgabe.** Der Fehlertext (Fehlerwert, Zeile beim Prozessende, Attribut
`error` einer Log-Zeile) beginnt mit dem Kopf `<klasse> [<code>]: <Ursache>`
auf `stderr`. `<klasse>` ist der Name der Klasse aus der Spalte *Klasse /
Bereich* ohne den Exit Code: `sonstiger Fehler`, `Konfiguration`, `Recording`,
`Netzwerk`, `Replay`, `nicht unterstützt`. Eine Warnung steht als eigenes
Log-Attribut `code`; der Meldungstext trägt keinen Code. Der Prozessausgang
bleibt unverändert, kein Code ändert ihn.

* *Eine Zeile.* Der Fehlertext ist eine Zeile; Kontinuationszeilen gibt es nicht.
  Ein Zeilenumbruch in der Ursache wird mit dem Leerraum danach durch ein
  Leerzeichen ersetzt, auch in einem Text aus einer Bibliothek. Zeilenumbruch ist
  LF, CR LF und ein einzelnes CR; Leerraum danach sind Leerzeichen, Tabulatoren und
  weitere Zeilenumbrüche. VT, FF, U+0085, U+2028 und U+2029 sind keine
  Zeilenumbrüche und bleiben stehen.
* *Fehlerkette.* Der Kopf steht genau einmal, am Anfang. Code und Klasse sind die
  des äußersten klassifizierten Fehlers; ein innerer klassifizierter Fehler trägt
  nur seine Ursache bei, ohne eigenen Kopf. Eine Kette ist auch ein Fehler mit
  eigenem Text um mehrere Ursachen: Er bleibt eine Meldung mit seinem ganzen Text,
  Code und Klasse sind die des ersten klassifizierten Fehlers unter ihm in
  Tiefensuche: jede Ursache wird ganz durchsucht, bis in ihre innersten Fehler,
  bevor die nächste Ursache an die Reihe kommt; ein klassifizierter Fehler tief in
  der ersten Ursache geht einem flachen in der zweiten vor. Kein innerer Fehler
  trägt einen Kopf; ohne klassifizierten Fehler darunter ist er `PGR-E1000`.
* *Gleichrangige Fehler.* Entstehen bei einem Ereignis mehrere Fehler
  nebeneinander, etwa beim Ende einer Session das Verbindungsende, das Schreiben
  des Recordings und das Schließen der Verbindung zum Upstream, ist jeder
  klassifizierte eine eigene Meldung mit eigenem Kopf, in der Reihenfolge ihres
  Entstehens; gemerkt wird der erste (LH-FA-13.b). Ein nicht klassifizierter Fehler
  daneben ist keine eigene Meldung: Sein Text folgt als Ursache der ersten
  klassifizierten. Sind alle nicht klassifiziert, ist es eine Meldung `PGR-E1000`,
  ihre Texte durch `; ` getrennt. Gleichrangig sind nur Fehler, die ohne eigenen
  Text nebeneinander stehen (eine reine Zusammenfassung); eine Hülle mit eigenem
  Text um mehrere Ursachen ist eine Kette (*Fehlerkette*).
* *Fremde Fehler.* Ein Fehler aus einer Bibliothek oder dem Betriebssystem trägt
  den Code der Stelle, an der der Recorder ihn einordnet, und sein Text folgt als
  Ursache; ordnet der Recorder ihn nicht ein, ist er `PGR-E1000` mit Kopf und
  Exit Code 1. Ein Fehlertext ohne Kopf kommt nicht vor.
* *Zustellung an den Client.* Die `ErrorResponse` eines Verbindungsfehlers
  (LH-FA-13.b) trägt den Schweregrad `FATAL`, als Meldungstext denselben
  Fehlertext wie das Attribut `error` der Log-Zeile, unabhängig vom Log-Level, und
  als SQLSTATE nach der Klasse `0A000` (nicht unterstützt), `08006` (Netzwerk),
  sonst `XX000`; weitere Felder trägt sie nicht. Ergibt der Fehler mehrere
  Meldungen (*Gleichrangige Fehler*), wird genau eine `ErrorResponse` zugestellt,
  mit Text und SQLSTATE der ersten, der gemerkten; die übrigen stehen nur im Log.
  Eine `ErrorResponse` des Servers,
  die `record` weiterleitet oder `replay` wiedergibt, ist keine Meldung des
  Recorders und trägt keinen Meldungscode (LH-FA-11.a).

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
- Der private Schlüssel (`--tls-key`) und sein Inhalt erscheinen weder in Logs noch
  in Fehlertexten; `config show` nennt nur den Pfad. Die Dateirechte des Schlüssels
  liegen in der Verantwortung des Anwenders.
- Diagnosen nennen SQL-Text im Klartext, auch mit Literalen darin; verborgen
  werden nur die Parameterwerte der Extended-Interaktion, und zwar immer: Sie
  erscheinen auf keinem Log-Level, auch nicht bei `debug`, und in keiner
  `ErrorResponse` an den Client (LH-FA-18.a, LH-FA-14.a).

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

**Warten in Tests.** Ein Test, der auf ein Ereignis wartet, das der Prüfling
herbeiführt (ein Wert auf einem Kanal oder dessen Schließen, das Ende eines
gestarteten Prozesses), wartet mit eigener Frist. Die Frist steht als Literal im
Test und ist höchstens 60 s lang. Läuft sie ab, wird der Test rot, und die Meldung
nennt das ausgebliebene Ereignis; die Zeitgrenze von `go test` ist keine Frist.
Auf das Ende eines gestarteten Prozesses wartet je Prozess genau eine Stelle, jeder
weitere Leser liest ihr Ergebnis. Ein Test-Double, das auf eine Freigabe durch den
Test wartet, fällt nicht unter diese Regel; der Test, der danach auf den Prüfling
wartet, schon.

## 10. Nicht zugesichert in v1

### SPEC-039 — Abgrenzung

Nicht zugesichert sind insbesondere:

- COPY,
- TLS zum Upstream im Record-Modus,
- Prüfung von Client-Zertifikaten,
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

## 11. Harness-Werkzeuge

Dieser Abschnitt legt die Werkzeuge der Prüfumgebung des Projekts fest, die
[`LH-QA-07`](lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) in Messmethode (4)
nennt: was ihre Prüfung zusagt, woran sie rot wird und was sie nicht prüft. Keine
Kennung dieses Abschnitts ist eine Zusage des Produkts, und keine erweitert eine
Anforderung des Lastenhefts. Die Schnittstelle eines Werkzeugs ist sein
`make`-Target. Die Punkte einer Kennung sind nummeriert; die Nummer ist ihre Adresse.

### SPEC-047 — Kopf der Pläne (`kopf-check`)

`make kopf-check` gleicht zwei Mengen von Kennungen ab, ohne über Inhalt zu urteilen:
Was §1 oder §2 eines lebenden Slice-Plans an Kennungen nennt, führt der Kopf des
Plans.

1. **Gegenstand.** Geprüft werden die Dateien, deren Name mit dem Präfix `slice` und
   einem Bindestrich beginnt und auf `.md` endet, flach in
   `docs/plan/planning/open/`, `docs/plan/planning/next/` und
   `docs/plan/planning/in-progress/`. Nicht geprüft werden Unterverzeichnisse dieser
   drei, `docs/plan/planning/done/` (flach und archiviert), Welle-Pläne, Roadmap,
   Register und andere Dateien.
2. **Kennungen.** Gezählt werden `LH-XX-NN`, `LH-XX-NN.x`, `SPEC-NNN` und `ARC-NNN`
   (X ein Großbuchstabe, N eine Ziffer, x ein Kleinbuchstabe) in genau dieser
   Schreibweise und als ganzes Wort; Wortzeichen sind Buchstabe, Ziffer und
   Unterstrich. Ein Punkt beendet das Wort: Folgt einer Hauptkennung nach dem Punkt
   keine Unterkennung (etwa `.ab`), zählt die Hauptkennung. Auszeichnung ist ohne
   Belang: Code-Span, Linktext, Hervorhebung und Codeblock zählen gleich. Link-Anker sind
   klein geschrieben und darum keine Nennung. Eine andere Schreibweise (klein, ohne
   führende Null, mit zusätzlicher Ziffer, ohne Bindestrich) ist keine Kennung, weder
   im Kopf noch in §1 oder §2.
3. **Bereich.** `SPEC-NNN bis SPEC-MMM`, ebenso mit `ARC`, steht im Kopf wie in §1 und
   §2 für jede Kennung von NNN bis MMM. Backticks um die beiden Enden sind erlaubt, und
   ein Zeilenumbruch innerhalb eines Absatzes zählt als Leerraum; in §1 und §2 reicht
   ein Bereich auch über Absatzgrenzen (siehe *Grenze*). Für die Kennungen
   dazwischen steht der Bereich nur, wenn beide Enden ganze Wörter derselben Klasse
   sind und NNN kleiner als MMM ist. Andere Bereichsformen (Gedankenstrich,
   Bindestrich) gibt es nicht.
4. **Gleichheit ist exakt.** Eine Unterkennung im Kopf deckt ihre Hauptkennung nicht,
   eine Hauptkennung keine Unterkennung.
5. **Kopf.** Der Kopf sind die Absätze, die mit `**Bezug:**` und
   `**Berührte Spec-Stellen:**` beginnen, vor der ersten Zeile `## `, je bis zur
   nächsten Leerzeile. Eine Leerzeile ist eine Zeile, die leer ist oder nur
   Leerzeichen und Tabs trägt. Ein Absatz beginnt nur in der ersten Zeile der Datei
   oder nach einer Leerzeile; eine Feldmarke mitten in einer Zeile oder in einer
   Folgezeile (direkt nach `# …`, nach `---` oder nach einem anderen Feld) ist kein
   Feld. Geprüft wird gegen die Vereinigung beider Felder; in welchem Feld eine
   Kennung steht, ist ohne Belang, und ein anderes Feld des Kopfs zählt nicht. `—`
   ist die leere Menge, Mehrfachnennung ist ohne Belang, eine Kopf-Kennung ohne
   Nennung in §1 oder §2 ist kein Befund.
6. **§1 und §2.** Ein Abschnitt reicht von der Zeile `## 1.` bzw. `## 2.`, die
   mitzählt, bis zur nächsten Zeile `## `; erkannt wird er an der Nummer, nicht am
   Titel. Unterüberschriften beenden ihn nicht. Ausgenommen ist der Absatz, der mit
   `Regeln dieser Sektion` beginnt, bis zur nächsten Leerzeile; dieselben Worte
   mitten in einem Absatz nehmen nichts aus. Alles übrige zählt, auch Abgrenzung,
   Herkunft und Bereits-Geliefertes; eine Markierung für eine Nennung ohne Anspruch
   gibt es nicht.
7. **Formfehler.** Fehlt einem geprüften Plan eines der beiden Kopf-Felder oder einer
   der beiden Abschnitte, ist das ein Befund.
8. **Ausgabe und Ausgang.** Je Befund eine Zeile auf stderr,
   `kopf-check: <pfad>: <abschnitt>: <befund>`, mit dem Pfad relativ zur Wurzel des
   Repos; der Abschnitt ist `Kopf`, `§1`, `§2` oder `Datei`. Der Befund lautet
   `<Kennung> fehlt im Kopf`, `Feld Bezug fehlt`, `Feld Berührte Spec-Stellen fehlt`,
   `Abschnitt ## 1. fehlt`, `Abschnitt ## 2. fehlt` oder `nicht lesbar`. Eine Kennung erscheint
   je Abschnitt einmal. Sortiert wird nach Pfad, Abschnitt und Befund in Bytefolge
   (`Kopf` vor `§1` und `§2`); stdout bleibt leer. Ein Plan, der nicht lesbar ist
   oder an dessen Auswertung die Prüfung scheitert, ist genau ein Befund
   `<pfad>: Datei: nicht lesbar`, gleich an welcher Position; die übrigen Pläne
   werden weiter gelesen. Die Prüfung liest alle Pläne und endet dann mit Ausgang 1
   bei mindestens einem Befund, mit 0 ohne Befund und ohne Ausgabe, mit 2 und einer
   Zeile, die `docs/plan/planning fehlt` nennt, wenn die Ablage
   `docs/plan/planning/` fehlt. Ein fehlendes Lifecycle-Verzeichnis enthält keinen
   Plan und ist kein Abbruch.
9. **Start ohne Stufung.** Die Prüfung ist vom ersten Lauf an voll scharf:
   `make kopf-check` hängt an der Gate-Kette von `make gates` und endet mit
   Fehlerstatus, wenn die Prüfung rot ist.

**Grenze.** Ob eine Kennung existiert, in welchem Kopf-Feld sie steht und was
außerhalb von §1 und §2 steht, prüft das Werkzeug nicht. Codeblöcke werden nicht
gesondert verfolgt: Eine Zeile `## ` in einem Codeblock beendet den Abschnitt. In §1
und §2 werden die gezählten Zeilen eines Abschnitts aneinandergehängt; ein Bereich
reicht dort auch über eine Leerzeile und über einen ausgenommenen Absatz
`Regeln dieser Sektion` hinweg. Das macht §1 und §2 nur strenger, nie still grün.

### SPEC-048 — Abdeckung je Anforderung und Pfad (`abdeckung`)

`make abdeckung` schreibt die Abdeckungstabellen `docs/user/abdeckung-*.md` aus den
Abdeckungs-Deklarationen der Tests; `make abdeckung-check` prüft, dass sie dem
erzeugten Stand entsprechen, und schreibt nichts. Beide wenden dieselben Regeln an.

1. **Deklaration.** Gelesen werden die Testdateien (Name auf `_test.go`) im Repo außer
   unter den Verzeichnissen `.git/` und `.harness/` an der Wurzel des Repos; ein
   gleichnamiges Verzeichnis tiefer im Baum wird gelesen. In einer gelesenen
   Testdatei steht direkt über `func Test…`, ohne Leerzeile und am Zeilenanfang:

   ```text
   // Abdeckung: <Anforderung>/<Pfad>, … — <Kurzbeschreibung>
   ```

   Folgezeilen mit `//` setzen sie fort und werden mit einem Leerzeichen angefügt.
   Fehlformen sind: eine eingerückte Deklaration, eine Deklaration ohne „ — “ und eine
   Deklaration, auf die nicht unmittelbar ein Test folgt (Leerzeile, andere Funktion,
   weitere Deklaration, Dateiende).
2. **Anforderung und Pfad.** Der Pfad ist bei funktionalen Anforderungen (`LH-FA`)
   eines der Akzeptanzkriterien `Happy`, `Boundary` oder `Negative`, bei
   Qualitätsanforderungen (`LH-QA`) und Randbedingungen (`LH-RB`) `Messung`; jede
   andere Paarung ist eine Fehlform. Eine Anforderung, die das Lastenheft nicht als
   Überschrift dritter Ebene führt (`### <Anforderung> …`), ist eine Fehlform.
3. **Nachweisart nach Ort.** Ein Test in einer Datei unter `test/integration/` ist ein
   E2E-Nachweis, jeder andere ein Unit-Nachweis. Jeder gelesene Test, dessen Name mit
   `TestE2E` beginnt, trägt eine Deklaration, gleich in welchem Verzeichnis er liegt;
   fehlt sie, ist das eine Fehlform.
4. **Tabellen.** Unter `docs/user/` entstehen vier Tabellen:
   - `docs/user/abdeckung-e2e.md` und `docs/user/abdeckung-unit.md`: je Deklaration
     und Paar aus Anforderung und Pfad eine Zeile mit Anforderung (Link auf das
     Lastenheft), Pfad, Kurzbeschreibung, Test und Datei, sortiert nach Anforderung,
     Pfad und Test; ein `|` in der Kurzbeschreibung ist maskiert.
   - `docs/user/abdeckung-gesamt.md`: je Anforderung mit mindestens einer
     Deklaration eine Zeile mit den Spalten `Happy`, `Boundary`, `Negative` und
     `Messung` — je die Nachweisarten, die den Pfad belegen, durch Komma getrennt,
     `—` für unbelegt, `n/a` für nicht anwendbar — und dem Stand `vollständig` oder
     `teilweise`.
   - `docs/user/abdeckung-vollstaendig.md`: nur die vollständig belegten
     Anforderungen, je mit den beteiligten Nachweisarten. Vollständig ist eine
     `LH-FA` mit `Happy`, `Boundary` und `Negative`, gleich aus welcher Nachweisart,
     eine `LH-QA` oder `LH-RB` mit `Messung`; Teilabdeckung steht hier nicht.
5. **Schreiben und Prüfen.** `make abdeckung` schreibt nur die Tabellen, die vom
   erzeugten Stand abweichen oder fehlen, setzt deren Rechte auf 0644 und meldet je
   geschriebener Tabelle eine Zeile `abdeckung: <pfad> geschrieben` auf stdout; eine
   aktuelle Tabelle bleibt unberührt, auch in ihren Rechten. `make abdeckung-check`
   vergleicht nur und meldet je veralteter Tabelle eine Zeile mit ihrem Pfad auf
   stderr.
6. **Ausgang.** Die Prüfung endet mit 0 ohne Fehlform und, beim Prüfen, ohne
   veraltete Tabelle. Sie endet mit 1 bei mindestens einer Fehlform, je mit einer
   Zeile auf stderr, die Datei und Zeile nennt (bei einem Test `TestE2E…` ohne
   Deklaration die Zeile des Tests), oder beim Prüfen bei mindestens einer veralteten
   Tabelle. Fehlt das Lastenheft, endet sie mit 2 und einer Meldung auf stderr.

**Grenze.** Ob eine Deklaration ihren Pfad tatsächlich belegt, prüft das Werkzeug
nicht: Es liest Deklarationen, es führt keine Tests aus. Führt das Lastenheft keine
Überschrift `### LH-…`, endet die Prüfung mit 1 ohne eine Zeile auf stderr.

### SPEC-049 — Lint-Profil (`lint`)

`make lint` prüft den Go-Code des Moduls statisch nach dem Profil in `.golangci.yml`
und mit drei eigenen Prüfungen; es schreibt nichts in den Arbeitsbaum.

1. **Gegenstand.** Geprüft werden alle Pakete des Moduls (`./...`) mit dem Build-Tag
   `integration`, also `cmd/`, `internal/` und `test/` einschließlich der
   Integrationstests und aller Testdateien. Das Profil ist die Datei `.golangci.yml`
   an der Wurzel des Repos; Pfade in ihr gelten relativ zu ihr
   (`run.relative-path-mode: cfg`) und beginnen mit `^`.
2. **Werkzeug und Umgebung.** golangci-lint `v2.14.0` aus dem Image
   `golangci/golangci-lint`, per Digest gepinnt, als Stufe `lint` des `Dockerfile`; der
   Pin steht nur dort. Die Stufe läuft auf der Plattform des Bau-Hosts, mit den Modulen
   der Stufe `deps`, ohne Netz, mit `GOTOOLCHAIN=local` und `GOFLAGS=-mod=readonly`.
   Trägt die Go-Version des Images die `go`-Zeile von `go.mod` nicht, scheitert die
   Stufe.
3. **Linter.** Aktiv sind genau diese (`linters.default: none`): `errcheck`, `govet`,
   `ineffassign`, `staticcheck`, `unused`, `containedctx`, `contextcheck`, `cyclop`,
   `dupl`, `fatcontext`, `forbidigo`, `funlen`, `gochecknoglobals`, `gochecknoinits`,
   `gocognit`, `gocyclo`, `gomodguard_v2`, `iface`, `inamedparam`, `interfacebloat`,
   `ireturn`, `maintidx`, `nestif`, `noctx`, `reassign`, `revive`, `testpackage`,
   `unparam`.
4. **Schwellen.** `cyclop` 15, `funlen` 100 Zeilen und 60 Anweisungen, `gocognit` 20,
   `gocyclo` 15, `nestif` 5, `maintidx` unter 20, `dupl` 150, `interfacebloat` 10.
5. **Einstellungen.**
   - `errcheck`: Ein ungeprüfter Fehler ist nur bei `fmt.Fprint`, `fmt.Fprintf`,
     `fmt.Fprintln`, `(net.Conn).Close` und `(net.Listener).Close` zulässig.
   - `ireturn`: Eine Funktion darf ein Interface nur liefern, wenn es `error`, leer,
     anonym, aus der Standardbibliothek oder generisch ist, unter
     `internal/hexagon/ports/` liegt oder aus dem Paket `pgproto3` stammt.
   - `forbidigo`: Verboten sind `fmt.Print`, `fmt.Printf`, `fmt.Println`, die
     eingebauten `print` und `println` sowie `os.Stdout` und `os.Stderr`; die beiden
     letzten sind unter `cmd/` und `test/` erlaubt.
   - `gomodguard_v2`: Erlaubt sind Importe aus der Standardbibliothek, dem Modul selbst
     und den Modulen `github.com/jackc/pgx/v5` und `go.yaml.in/yaml/v3`. Jedes andere
     Modul ist ein Befund, auch eines, das `go.mod` als indirekt führt.
   - `revive`: die Regeln `blank-imports`, `context-as-argument`, `context-keys-type`,
     `dot-imports`, `empty-block`, `error-naming`, `error-return`, `error-strings`,
     `errorf`, `exported`, `if-return`, `increment-decrement`, `indent-error-flow`,
     `package-comments`, `range`, `receiver-naming`, `redefines-builtin-id`,
     `superfluous-else`, `time-naming`, `unexported-return`, `unreachable-code`,
     `unused-parameter`, `var-declaration`, `var-naming` und `unused-receiver`, auf jedem
     Pfad, auch unter `internal/`. `exported` prüft die Form des Doc-Kommentars, nicht
     seine Sprache.
   - `testpackage`: Übersprungen wird nur eine Datei namens `export_test.go`
     (`skip-regexp` `(^|/)export_test\.go$`).
   - `contextcheck`: ohne Einstellung und ohne Ausnahme, auch in Testdateien. Meldet
     es in einem Test einen neuen Kontext, kommt der Kontext aus `t.Context()` des
     laufenden Tests oder Subtests; in einer Funktion für `t.Cleanup`, die nach dessen
     Abbruch läuft, aus `context.WithoutCancel(t.Context())`. `context.Background()` an
     einer Stelle ohne Befund bleibt zulässig.
   - `gochecknoglobals`: ohne Einstellung; die Variable `version`, die Leerstelle `_`
     und Fehlerwerte mit Präfix `Err` lässt der Linter selbst zu.
6. **Kein `//nolint`.** Eine eigene Prüfung meldet jede Zeile einer Datei `*.go` unter
   `cmd/`, `internal/` und `test/`, in der auf `//` oder `/*`, gefolgt von beliebig
   vielen Leerzeichen, Tabs oder `/`, das Wort `nolint` folgt — gleich in welcher
   Schreibweise, an welcher Stelle der Zeile und mit welchem Zusatz (`:linter`,
   Begründung). Maßgeblich ist jedes `//` und jedes `/*` der Zeile, nicht nur das, mit
   dem der Kommentar beginnt: `// x //nolint` ist ein Befund. Folgt `nolint` erst nach
   einem anderen Wort auf das letzte Kommentarzeichen davor (`// siehe nolint`), ist
   das kein Befund.
7. **Export-Test-Brücke.** Unit-Tests liegen im Paket `<name>_test`. Auf Unexportiertes
   greifen sie nur über die Datei `export_test.go` im Verzeichnis des Pakets zu; sie
   gehört zum Paket `<name>` und ist dort die einzige Testdatei. Sie enthält nur
   Typ-Aliase, Konstanten und Funktionen oder Methoden, die an Unexportiertes
   weiterreichen. Zustand setzt sie nur an einem Wert, den sie übergeben bekommt, nie
   auf Paketebene. Einen Wert eines unexportierten Typs legt nur der Produkt-Code an:
   Weder die Brücke noch ein Test erzeugt ihn selbst, auch nicht als lokalen Wert, als
   Literal oder über einen Typ-Alias. Die Brücke bekommt ihn übergeben oder reicht an
   eine Funktion des Pakets weiter, die ihn erzeugt; so prüft ein Test nur Zustände,
   die das Produkt erzeugt. Diese Funktion ruft auch der Produkt-Code, um denselben
   Wert zu erzeugen; eine Funktion, die nur die Brücke ruft, gibt es nicht. Werte
   exportierter Typen, die sie annimmt, etwa eine Verbindung, stellt der Test. Was nur
   an einem Wert eines unexportierten Typs zu sehen ist, den die Brücke weder
   übergeben bekommt noch über eine solche Funktion erzeugen lässt, prüft der Test
   über die exportierte Schnittstelle. Jede andere Testdatei im Paket `<name>` ist ein
   Befund von `testpackage`; eine Variable auf Paketebene in der Brücke ist ein Befund von
   `gochecknoglobals`; eine Funktion, deren Name mit `Test`, `Benchmark`, `Example`
   oder `Fuzz` beginnt, meldet in der Brücke eine eigene Prüfung.
8. **Ausnahmen.** Eine Ausnahme steht nur in `.golangci.yml`, als Einstellung nach
   Punkt 5 oder als Regel unter `linters.exclusions.rules`. Unmittelbar über jeder
   Regel steht ein Kommentarblock, dessen erste Zeile mit `# Why:` beginnt; fehlt er,
   meldet das eine eigene Prüfung. Das Profil hat dafür eine feste Form: `linters:`
   ohne Einzug, darunter `exclusions:` mit zwei Leerzeichen, darunter `rules:` mit vier,
   jede Regel als Eintrag `- ` mit sechs, alle drei Schlüssel in Blockform (nach dem
   Doppelpunkt höchstens Leerraum und ein Kommentar). Ein Eintrag ist auch ein `-` allein
   auf der Zeile, dessen Inhalt auf den Fortsetzungszeilen steht; über ihm steht der
   Kommentarblock wie über jedem Eintrag. Ein Schlüssel `exclusions` an anderer Stelle
   oder in Flussform und ein Schlüssel `rules` unter `exclusions` mit anderem Einzug
   oder in Flussform ist ein Befund `Form nicht erkannt` mit seiner Zeile. Ebenso ist
   unter `rules` jede Zeile ein solcher Befund, die weder Kommentar noch Leerzeile ist,
   weder ein Eintrag `- ` mit sechs Leerzeichen noch dessen Fortsetzung mit mindestens
   acht unter einem solchen Eintrag; ein Eintrag `- ` mit anderem Einzug als sechs ist
   es immer. Unter `exclusions` stehen nur `warn-unused: true` und `rules`. Zulässig sind nur diese Regeln, jede dauerhaft und
   mit einem Grund, der auch für neuen Code gilt: für Testdateien `cyclop`,
   `gocognit`, `gocyclo`, `nestif`, `funlen`, `noctx`, `unparam` und `revive` mit
   `unused-parameter` und `unused-receiver`; `staticcheck` mit `ST1005`, weil die
   Fehlertexte deutsch sind und mit einem Substantiv beginnen dürfen;
   `gochecknoglobals` für eine benannte Nachschlage-Tabelle oder einen benannten
   Sentinel-Wert, die nach der Initialisierung nur gelesen werden, je Regel mit Datei
   und Namen; `forbidigo` für `os.Stdout` und `os.Stderr` nach Punkt 5. Eine Regel,
   die einen Befund nur deshalb ausblendet, weil der Bestand ihn trägt, gibt es nicht.

   Eine Regel, die im Lauf keinen Befund ausblendet (`exclusions.warn-unused: true`), ist
   ein Befund.
9. **Ausgabe und Ausgang.** golangci-lint liest das Profil nur aus `.golangci.yml`
   (`-c`), nie aus einer Default-Suche. Die Stufe prüft in dieser Reihenfolge:
   - **Profil vorhanden.** Fehlt `.golangci.yml`, schreibt sie
     `lint: .golangci.yml: fehlt`; Schema-Prüfung und golangci-lint laufen dann nicht,
     die eigenen Prüfungen nach Punkt 6 und 7 schon.
   - **Schema.** `golangci-lint config verify` prüft das Profil gegen das eingebettete
     Schema, ohne Netz. Lehnt es das Profil ab, etwa wegen eines unbekannten Schlüssels,
     schreibt die Stufe `lint: .golangci.yml: von golangci-lint abgelehnt` und dazu die
     Meldung von `config verify`; golangci-lint läuft trotzdem.
   - **Befunde.** golangci-lint gibt seine Befunde im Textformat ungekürzt aus
     (`max-issues-per-linter: 0`, `max-same-issues: 0`, `uniq-by-line: false`). Jede
     Meldung ist ein Befund, auch mehrere auf derselben Zeile.
   - **Ungenutzte Regel.** Je Warnung von `warn-unused` schreibt die Stufe eine Zeile
     `lint: .golangci.yml: Regel ohne Befund: <feld>: <wert>, …` ohne Zeilennummer. Die
     Felder stehen in dieser Reihenfolge und nur, wenn die Warnung sie nennt: `Linter`,
     `Pfad`, `Pfad außer`, `Text`, `Quelle`, je mit dem Wert aus der Warnung. Nennt sie
     keines davon, lautet die Zeile `lint: .golangci.yml: Regel ohne Befund: Felder
     nicht erkannt`; die Warnung selbst steht daneben auf der Ausgabe.

   Die eigenen Prüfungen nach Punkt 6 bis 8 schreiben je Befund eine Zeile
   `lint: <pfad>:<zeile>: <befund>`. Jede Prüfung läuft, auch wenn eine andere einen
   Befund hat. Die Stufe endet mit einem Ausgang ungleich 0 bei mindestens einer
   `lint:`-Zeile oder einem Befund von golangci-lint.
10. **Werkzeug, dann Gate.** Solange der Bestand außerhalb der Regeln nach Punkt 8
    Befunde hat, ist `make lint` ein Werkzeug ohne Gate: Es meldet alle Befunde des
    Moduls mit Pfad und endet nach Punkt 9, hängt aber nicht an der Gate-Kette von
    `make gates`. Ohne Befund im Bestand hängt es an der Gate-Kette und endet `make
    gates` mit Fehlerstatus, wenn es rot ist.

**Grenze.** Ob ein `Why:` zutrifft, ob die Brücke nur weiterreicht und ob ein Wert
eines unexportierten Typs nur aus dem Produkt-Code stammt, prüft das Werkzeug nicht. Dateien mit der Markierung für generierten Code nimmt golangci-lint nach
seinem Default aus; das Modul hat keine. Ebenso wenig prüft es, ob eine Einstellung
nach Punkt 5 ihren Grund als Kommentar trägt und ob unter `exclusions` nur die Schlüssel
nach Punkt 8 stehen. Ob eine Regel unter `rules` zu den zulässigen nach Punkt 8 gehört
und ob sie mehr als Bestand ausblendet, prüft es ebenfalls nicht: Eine neue Regel mit
`# Why:`, die einen Befund ausblendet, lässt die Stufe grün; das ist Urteil des Review. Eine Regel, die `config verify` annimmt und golangci-lint erst
beim Laden ablehnt, macht die Stufe nach Punkt 9 rot, aber ohne `lint:`-Zeile; zu
sehen ist dann nur die Meldung von golangci-lint. Ein `nolint` nach Punkt 6 in einem String-Literal ist ebenfalls ein Befund. Testdateien im Paket `main` lässt
`testpackage` zu; unter `cmd/` gibt es keine.
Ist das Profil durch einen Eintrag unter `rules` mit falschem Einzug kein gültiges YAML
mehr, sagt Punkt 8 nur den Befund `Form nicht erkannt` an der ersten solchen
Eintragszeile zu; wie die Prüfung die Zeilen danach einordnet, ist nicht zugesagt. Rot
ist die Stufe durch diesen Befund nach Punkt 9. Sind Image, Module und Inhalt des
Build-Kontexts dieselben wie bei einem früheren Lauf mit Ausgang 0, nimmt der Build das
Ergebnis der Stufe aus dem Cache, ohne neu zu prüfen; ein Lauf mit Ausgang ungleich 0
liegt nie im Cache. Eine Datei, deren Größe und Änderungszeit seit dem letzten Lauf
gleich blieben, überträgt der Build nicht neu; ihre Änderung sieht der Lauf dann nicht.

### SPEC-050 — Struktur-Kennungen im Commit-Träger (`commit-msg`)

`.githooks/commit-msg` ist der Träger, den git vor jedem Commit mit dem Pfad der Datei
der vorgeschlagenen Message aufruft; Exit ungleich 0 bricht den Commit ab. Aktiviert
wird er mit `make hooks-install`, seine Gegenprobe ist `make hook-gegenprobe`. Diese
Kennung legt fest, wann er eine Message wegen einer Struktur-Kennung ablehnt. Was er
sonst annimmt und was er an die mitgelieferte Prüfung weiterreicht, legt sie nicht
fest.

1. **Ort.** Die Ablehnung liegt im Träger `.githooks/commit-msg` selbst, nicht in der
   mitgelieferten Prüfung `tools/harness/commit-msg-traceability.sh`.
2. **Lesebereich.** Gelesen wird die Message-Datei bis vor die erste Zeile, die genau
   `# ------------------------ >8 ------------------------` lautet (je 24 Bindestriche),
   die Scissors-Zeile von `git commit -v`; eine Zeile in anderer Form beendet den
   Bereich nicht. Nicht gelesen werden Kommentarzeilen: Zeilen, deren erstes Zeichen
   nach führenden Leerzeichen und Tabs `#` ist.
3. **Struktur-Kennung.** `SPEC-NNN` und `ARC-NNN` mit genau drei Ziffern, in genau
   dieser Schreibweise und als ganzes Wort; Wortzeichen sind Buchstabe, Ziffer und
   Unterstrich, jedes andere Zeichen grenzt ab, auch `-`, `/`, `#`, `.`, Backtick,
   Anführungszeichen, `>` und Wagenrücklauf. Die Kennung zählt darum auch in Backticks,
   in einem Zitat, in Klammern, als Linktext und als Teil eines Pfads oder
   Dateinamens. Keine Struktur-Kennung ist eine andere Schreibweise (`spec-049`,
   `Spec-049`, `arc-003`), eine andere Ziffernzahl (`SPEC-49`, `SPEC-0491`) oder ein
   Platzhalter (`SPEC-<NNN>`, `SPEC-NNN`, `ARC-*`); ein Link-Anker wie
   `spezifikation.md#spec-049--lint-profil-lint` ist klein geschrieben und darum keine
   Nennung. Ob die Kennung in Spezifikation oder Sicht existiert, ist ohne Belang. In
   `SPEC-035 bis SPEC-040` sind die beiden Enden je eine Nennung.
4. **Merge und Revert.** Beginnt der Betreff mit `Merge ` oder `Revert `, sucht der
   Träger keine Struktur-Kennung. Betreff ist die erste Zeile des Lesebereichs, die
   weder leer ist noch nur Leerraum trägt noch Kommentarzeile ist, ohne führende
   Leerzeichen und Tabs.
5. **Vorrang.** Die Suche nach Struktur-Kennungen läuft vor jeder Annahme: Eine
   Message mit Struktur-Kennung wird abgelehnt, auch wenn sie daneben eine Kennung
   trägt, die der Träger oder die mitgelieferte Prüfung annimmt. Eine Message ohne
   Struktur-Kennung im Lesebereich, eine nach Punkt 4 ausgenommene und ein Aufruf ohne
   lesbare Message-Datei gehen den Weg, den der Träger ohne diese Kennung geht.
6. **Ausgabe und Ausgang.** Je Nennung eine Zeile auf stderr,
   `commit-msg: Zeile <n>: Struktur-Kennung <Kennung>`, mit `<n>` der Zeilennummer in
   der Message-Datei ab 1, Kommentarzeilen mitgezählt. Eine Kennung, die in derselben
   Zeile mehrmals steht, erscheint für diese Zeile einmal; geordnet wird nach Zeile,
   innerhalb einer Zeile nach erstem Auftreten. Danach folgt genau eine Zeile
   `commit-msg: Struktur-Kennungen (SPEC-NNN, ARC-NNN) gehoeren nicht in die Commit-Message`.
   stdout bleibt leer, der Ausgang ist 1.

**Grenze.** Ein Klon ohne `make hooks-install` und ein Commit mit
`git commit --no-verify` bleiben ungeprüft. Ein anderes Kommentarzeichen als `#`
(`core.commentChar`) kennt der Träger nicht: Eine solche Zeile wird gelesen. Eine selbst
geschriebene Message, deren Betreff mit `Merge ` oder `Revert ` beginnt, geht nach
Punkt 4 ungeprüft durch. Eine Kennung, die ein Zeilenumbruch teilt, ist keine Nennung.
Gelesen wird nur die Message-Datei; was der Commit ändert, liest der Träger nicht.

## 12. Historie

Regeln dieser Sektion: **kein ADR- und kein Slice-Verweis.** Die Decken-Regel
gilt für alle drei Spec-Straten, auch hier — welche ADR eine Festlegung
schärft, deklariert die ADR aufwärts in ihrem `Schärft:`-Feld
(Baseline-Regelwerk `modul-03-spec.md` §Ziel-Form: Spezifikation).

| Datum | Änderung |
|---|---|
| 2026-10-03 | Initial |
| 2026-10-03 | TLS zum Client (`LH-FA-23.a`), Antwortvergleich beim Einspielen (`LH-FA-24.a`), eigene Zertifizierungsstelle (`LH-FA-20.a`), Codes `PGR-E2007`, `PGR-E5004`, `PGR-E6003`, `PGR-W3002` |
| 2026-10-04 | Protokollrand: Verbindung ohne erste Nachricht und erste Nachricht ohne PGWire-Form (`PGR-W3003`, `SPEC-045`) |
| 2026-10-04 | Extended Query: `ParameterStatus` unter den Server-Nachrichten (`LH-FA-18.a`, `SPEC-041`, wie schon in `LH-FA-24.a` vorausgesetzt); Felder je Client-Nachricht und `param_types` der `parameter_description` (`SPEC-041`) |
| 2026-10-04 | Recording: Stelle von `type` bei einfacher und Extended-Interaktion (`SPEC-001`, `SPEC-041`); fehlendes oder mit `null` belegtes Feld einer Client-Nachricht, `param_types` nur und stets an `parameter_description`, `request`/`responses`/`groups` mit `null` als beschädigt (`SPEC-041`) |
| 2026-10-04 | Recording: leerer oder `null`-Wert von `type` (`SPEC-001`), `param_types` mit `null` an einer anderen Server-Nachricht, `client` und `server` je Gruppe Pflicht, `server: []` zulässig (`SPEC-041`) |
| 2026-10-04 | Recording: Anker, Aliase und Merge-Keys von YAML sind beschädigt (`SPEC-001`) |
| 2026-10-05 | Recording: jedes neue Feld erhöht `version`; ein Leser lehnt unbekannte Felder ab und liest jede `version` bis zu seiner eigenen (`SPEC-001`) |
| 2026-10-05 | Extended Query im Record: Server-Nachrichten nach der nächsten Client-Nachricht gehören zur folgenden Gruppe; Herunterfahren nach dem `ReadyForQuery` der laufenden Extended-Interaktion; `Query` während einer laufenden Extended-Interaktion, Zielart außer `S`/`P` und Interaktionen ohne die Form der Aufzeichnung sind nicht unterstützt (`LH-FA-18.a`) |
| 2026-10-05 | Record: beide Richtungen unabhängig vermittelt, Ende der Client-Verbindung beim Senden an oder Warten auf den Server (`LH-FA-18.a`); beim Herunterfahren schon gelesene Client-Nachricht wird verarbeitet (`LH-FA-13.a`) |
| 2026-10-05 | Herunterfahren: danach beginnt keine neue Interaktion, die Session endet nach der laufenden (`LH-FA-13.a`); Ende einer Session schließt die Client-Verbindung (`LH-FA-18.a`) |
| 2026-10-05 | Herunterfahren: eine Client-Nachricht um das Signal herum wird entweder noch verarbeitet oder nicht weitergeleitet, statt an den Lesestand beim Beginn gebunden (`LH-FA-13.a`) |
| 2026-10-05 | Herunterfahren: Frist `--shutdown-timeout` ab dem ersten Signal, Zwangsende mit `PGR-E4006`, weiteres Signal lässt die Frist ablaufen, Info-Zeile beim Beginn (`LH-FA-13.a`, `SPEC-046`) |
| 2026-10-05 | Replay: Lebendprüfungen außerhalb der Reihe beantwortet, aufgezeichnete übersprungen (`LH-FA-09.a`, `LH-FA-03.a`, `LH-FA-03.b`, `LH-FA-10.a`, `LH-FA-12.a`, `SPEC-011`); falsche Protokollart am Cursor, auch mitten in einer Extended-Interaktion, ist ein Mismatch (`LH-FA-18.a`) |
| 2026-10-05 | Replay: vertikaler Tabulator als Leerraum einer Lebendprüfung erst ab PostgreSQL 17, nach `server_version` der Session, vor der Zuordnung nicht (`LH-FA-09.a`); Startfehler bei Sessions nur aus Lebendprüfungen (`LH-FA-03.a`); Einspielen und `connection` bei Sessions nur aus Lebendprüfungen (`LH-FA-12.a`) |
| 2026-10-05 | Replay: verbraucht, Zeitpunkte der Prüfung, Meldung je Session und für nie zugeordnete Sessions, Rangfolge und keine Zustellung von `PGR-E5002`, Option bei anderen Kommandos (`LH-FA-03.b`, `LH-FA-13.b`); Konfiguration: Werte boolescher Optionen, leere Umgebungsvariable, Mehrfachangabe, Umgebungsvariable einer fremden Option (`LH-FA-17.a`); `PGR-W2001` auch für nie zugeordnete Sessions (`SPEC-034`) |
| 2026-10-05 | Konfiguration: ungültige Umgebungsvariable ist ein Fehler, auch wenn die Kommandozeile dieselbe Option setzt (`LH-FA-17.a`) |
| 2026-10-05 | Hilfe geht jeder Prüfung von Optionen, Umgebungsvariablen und Konfigurationsdatei vor; `version` liest keine Konfiguration (`LH-FA-01.a`, `LH-FA-17.a`) |
| 2026-10-05 | Hilfe-Angabe abgeschlossen: vier Formen mit beliebigem `=`-Wert, auch an Wertstelle, vor dem Kommando und nach unbekanntem Kommando (`LH-FA-01.a`) |
| 2026-10-05 | `--` beendet die Optionen auch an der Stelle eines Optionswerts; der Wert `--` nur mit `=` (`LH-FA-01.a`) |
| 2026-10-05 | Fehlerantwort vor dem Abbruch: maßgeblich ist der Lesestand beim Bemerken des Verbindungsendes, gleich in welcher Richtung; Sendefehler mit gelesener `ErrorResponse` ist `PGR-E6001`; nicht lesbare Nachricht ohne Verbindungsende ist `PGR-E6001` unabhängig von der Vorgeschichte (`LH-FA-02.b`) |
| 2026-10-05 | Diagnosefelder: leerer Wert gilt nur bei den von der Bibliothek benannten Feldcodes als fehlend, ein anderer Code wird auch leer aufgezeichnet (`LH-FA-11.a`); Fehlerantwort vor dem Abbruch zählt nur gelesen, Sendefehler der Anfrage bleibt `PGR-E4003` (`LH-FA-02.b`) |
| 2026-10-05 | Record: Fehlerantwort vor dem Abbruch der Upstream-Verbindung ohne `ReadyForQuery` ist nicht unterstützt (`PGR-E6001`), gleich welcher Schweregrad und welche Stelle; Diagnose mit SQLSTATE des Servers, Session nicht übernommen, Aufzeichnung ohne `ready_for_query` am Ende beschädigt (`LH-FA-02.b`, `LH-FA-11.a`) |
| 2026-10-05 | Record: Weitergabe der Antworten einer einfachen Anfrage an den Client, auch bei einer nicht unterstützten Antwort nach Ergebnissen (`LH-FA-02.b`); Nachrichten des Upstreams zwischen Interaktionen gehören zur nächsten (`LH-FA-05.a`); Diagnosefelder mit leerem Wert oder Zahl `0` gelten als fehlend, Feldreihenfolge nicht aufgezeichnet (`LH-FA-11.a`) |
| 2026-10-06 | Diagnose: Parameterwerte auf keinem Log-Level und in keiner `ErrorResponse`, Nummer des ersten abweichenden Parameters, bei abweichender Zahl beide Anzahlen (`LH-FA-18.a`, `SPEC-033`); Log-Level genau vier, Strenge des Werts, Inhalt der Stufen, Zeilenform, Zeile beim Prozessende auf jeder Stufe, `stdout` nur für Hilfe, `version` und `config show` (`LH-FA-14.a`); Fehlertext einzeilig, Kopf einmal mit dem Code des äußersten Fehlers, fremde Fehler, Felder der `ErrorResponse` (`SPEC-034`) |
| 2026-10-06 | Gleichrangige Fehler als eigene Meldungen, erster gemerkt, nicht klassifizierter als Ursache; Zeilenumbruch LF, CR LF, CR; Attribut `error` nur mit Kopf, `grund` für Bibliothekstexte (`SPEC-034`, `LH-FA-14.a`); nicht annehmbare Verbindung ist Verbindungsfehler `PGR-E4000` (`LH-FA-13.b`) |
| 2026-10-06 | Mehrere Meldungen: die `ErrorResponse` trägt die erste; Hülle mit eigenem Text um mehrere Ursachen ist eine Kette mit dem ersten klassifizierten Code, gleichrangig nur eine reine Zusammenfassung (`SPEC-034`) |
| 2026-10-06 | Harness-Werkzeuge: Lint-Profil mit Linter, Schwellen, Einstellungen, Verbot von `//nolint`, Export-Test-Brücke, dauerhafte Ausnahmen ohne Stufen, Ausgabe und Ausgang, Werkzeug vor dem Gate (`SPEC-049`) |
| 2026-10-06 | Lint-Profil: fehlendes Profil, Prüfung gegen das Schema, ungenutzte Regel als `lint:`-Zeile, jede Meldung ein Befund (`SPEC-049`) |
| 2026-10-06 | Lint-Profil: feste Form des Ausnahme-Abschnitts, zweites Kommentarzeichen vor `nolint`, Zeile der ungenutzten Regel, generierter Code nach Default (`SPEC-049`) |
| 2026-10-06 | Lint-Profil: Einzug der Regel-Einträge unter `rules`, Ablehnung beim Laden als Grenze (`SPEC-049`) |
| 2026-10-06 | Lint-Profil: Brücke und Test legen keinen Wert eines unexportierten Typs an; Kontext im Test bei einem Befund von `contextcheck` (`SPEC-049`) |
| 2026-10-07 | Lint-Profil: die Funktion, an die die Brücke zum Erzeugen weiterreicht, ruft auch der Produkt-Code; Eingaben exportierter Typen stellt der Test (`SPEC-049`) |
| 2026-10-07 | Lint-Profil: ungültiges YAML im Ausnahme-Abschnitt und Ergebnis aus dem Cache des Builds als Grenze (`SPEC-049`) |
| 2026-10-07 | Lint-Profil: Leerraum nach `exclusions:` und `rules:` ist Blockform, ein `-` allein ist ein Eintrag (`SPEC-049`) |
| 2026-10-08 | Lint-Profil: ob eine Regel zu den zulässigen gehört, als Grenze (`SPEC-049`) |
| 2026-10-06 | Kette mit mehreren Ursachen: der erste klassifizierte Fehler in Tiefensuche (`SPEC-034`) |
| 2026-10-06 | Harness-Werkzeuge: Abschnitt angelegt; Prüfung des Kopfs lebender Pläne (`SPEC-047`) und Abdeckung je Anforderung und Pfad (`SPEC-048`) mit ihrem heutigen Vertrag übertragen (`LH-QA-07`, Messmethode 4) |
| 2026-10-08 | Harness-Werkzeuge: Commit-Träger lehnt Struktur-Kennungen in der Commit-Message ab; Lesebereich, Schreibweise und Wortgrenze, Merge und Revert, Vorrang vor der Annahme, Ausgabe und Ausgang (`SPEC-050`) |
| 2026-10-08 | Herunterfahren: Attribut `sessions` der Zeile beim Beginn, ohne offene Verbindung, Startphase, was die Frist begrenzt, Zwangsende auch im Aufbau, `PGR-E4006` nur bei unvollständiger Interaktion und ohne `PGR-E4003`, Inhalt der Meldung, zweites Signal auch bei `0` (`LH-FA-13.a`, `LH-FA-14.a`); Form der Dauer von `--shutdown-timeout` (`LH-FA-17.a`); Schreibfrist der Fehlerantwort beim Ende einer Session (`SPEC-051`, `LH-FA-18.a`) |
