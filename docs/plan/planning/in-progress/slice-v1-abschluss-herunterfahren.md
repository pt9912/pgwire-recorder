# Slice slice-v1-abschluss-herunterfahren: Herunterfahren mit Frist in record und replay

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-15`](../../../../spec/lastenheft.md#lh-fa-15--ci-eignung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0003](../../adr/0003-pgwire-server-ist-driving-adapter.md), [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md)

**Berührte Spec-Stellen:** `LH-FA-13.a` · `LH-FA-13.b` · `LH-FA-02.b` · `LH-FA-03.b` · `LH-FA-14.a` · `LH-FA-17.a` · `LH-FA-18.a` · `SPEC-013` bis `SPEC-019` · `SPEC-034` · `SPEC-038` · `SPEC-046` · `SPEC-051`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `record` und `replay` fahren auf `SIGINT`/`SIGTERM` kontrolliert herunter, warten dabei höchstens `--shutdown-timeout` lang auf laufende Sessions und enden mit dem Exit-Code der gemerkten Klasse; eine durch die Frist abgebrochene Interaktion ist `PGR-E4006`, eingestuft vom Use Case
([ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md)).

**Übernimmt:** `slice-v1-abschluss-betrieb` — dessen Teil *Herunterfahren* (Entscheidung des
Nutzers vom 2026-10-08, Schnitt nach F-345 in drei Slices; dort §7 `Gegenstand:`). Im
Einzelnen, je mit dem ursprünglichen Geber:

- **Aus `slice-extended-query-record`** (Risiko F-309/F-317, Validierung
  `docs/reviews/2026-10-05-validierung-slice-extended-query-record.md`, Frage 1): Das
  Herunterfahren bekommt eine Obergrenze. Heute wartet `record` ohne Frist auf die
  laufende Interaktion jeder Session; im Container beendet der `SIGKILL` nach der
  Stopp-Frist den Prozess, und die Aufzeichnung jeder noch wartenden Session fehlt ganz.
  Frist gezählt ab dem ersten Signal, Default 5 s, `0` heißt ohne Frist, setzbar per
  Option und Umgebungsvariable ([`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)). Nach Ablauf endet jede noch laufende Session
  zwangsweise wie bei einem Abbruch nach `LH-FA-02.b`: die laufende Interaktion wird nicht
  übernommen, die abgeschlossenen bleiben; danach wird die Aufzeichnung geschrieben.
  `PGR-E4006` der Klasse 4 (vergeben in `SPEC-034`), Exit-Code `4`; das Log nennt je
  abgebrochener Session die Session und die verworfene Interaktion. Ein zweites Signal
  lässt die Frist sofort ablaufen. Beim Beginn eine Info-Zeile mit der Zahl der Sessions,
  auf die gewartet wird.
- **Aus `slice-extended-query-replay`** (Folge-Review F-330): `replay` beantwortet beim
  Herunterfahren eine begonnene Extended-Interaktion bis zu ihrem `Sync` und wartet dabei
  ohne Frist auf einen pausierenden Client. `replay` setzt die Frist durch Schließen der
  Client-Verbindung durch (`replaySitzung`); dieses Ende behandelt `replaySitzung` heute
  als regulär. Der Adapter meldet dem Replay-Use-Case das Zwangsende als Ereignis, und der
  Use Case stuft eine unvollständige Interaktion als `PGR-E4006` ein, nicht der Adapter
  (Architektur-Sicht §4.5 *Adapter-Verantwortung*: die Einstufung eines Verbindungsendes
  liegt nicht beim PGWire-Server).
- **Aus `slice-replay-semantik-mismatch`:** der Test, dass eine durch die Frist
  zwangsweise beendete Replay-Session mit `--fail-on-unconsumed` `PGR-E4006` vor
  `PGR-E5002` merkt und mit Exit-Code `4` endet (`LH-FA-03.b`).
- **Aus `slice-v1-abschluss-sessions`** (dort §1, Abgrenzung): die Signalbehandlung.
- **Aus `slice-harness-integration-wait`** (dort §1, Abgrenzung): die Tests mit Signal und
  Frist je Modus, geschrieben über die Helfer, die jener umbaut; sie sind die Tests in
  DoD-Punkt 1 und 2. Dazu aus dessen Verifikation (V-105) die Randform *`stderr` zur
  Laufzeit lesen*, entschieden in §6: Kein Test liest `stderr`, solange der Prozess läuft.
- **Aus `slice-harness-blackbox-pgwire`** (Review F-451, Verifikation V-87; Abgrenzung in
  `slice-tests-ueberlebende-mutanten`): der Wert der Schreibfrist der Fehlerantwort beim
  Session-Ende (`meldeFrist`), entschieden als `SPEC-051` (1 s), und ein Test mit der
  Schranke als Literal; Randform in §6.
  Er gehört zu DoD-Punkt 1, weil die Frist einen Teil des Budgets von `--shutdown-timeout`
  verbraucht.
- **Aus `slice-tests-ueberlebende-mutanten-driving`** (dort §1, Abgrenzung; Verifikation
  V-95 zu `slice-lint-bestand-driving`): das Ende einer Replay-Sitzung durch die Frist und
  `PGR-E4006`. Der Test dazu entsteht hier im Paket `pgwire_test` über `pgwire.Handle` mit
  Replay-Fake (DoD-Punkt 2); jener Slice schreibt später den Test zu G2 (ein Ende, das der
  Client auslöst) daneben, ohne diesen zu ändern.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Atomares Schreiben der Aufzeichnung, `--output` und `--force` — `slice-v1-abschluss-schreiben`;
  dieser Slice schreibt nach dem Zwangsende so, wie der Recording-Adapter heute schreibt.
- Der Schlüssel der Frist in der Konfigurationsdatei und der allgemeine Leser der
  Umgebungsvariablen — `slice-v1-abschluss-konfiguration`; hier wird
  `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT` gelesen wie die Umgebungsvariablen der übrigen
  Optionen heute, und jener Slice ersetzt den Leser.
- Ein Ende der Replay-Sitzung, das der Client auslöst (G2) —
  `slice-tests-ueberlebende-mutanten-driving`.
- Container-Image — `slice-v1-abschluss-container`.
- Neue Meldungscodes über `PGR-E4006` hinaus — Bestand bleibt stehen: `SPEC-034` vergibt
  `PGR-E4006`, die Tabelle liegt seit `slice-replay-semantik-meldungscodes`; hier kommt nur
  die Konstante in die Code-Tabelle des Quelltexts.
- Eine Schreibfrist für die Fehlerantwort anderer Verbindungsfehler im Replay (etwa
  `PGR-E5001` an einen Client, der nicht liest) — Bestand bleibt stehen: `SPEC-051` gilt im
  Record für jedes Session-Ende, im Replay nur für `PGR-E4006`; ein blockiertes Schreiben im
  Replay endet spätestens mit dem Zwangsende dieses Slice.
- Code in den Driven-Adaptern (`recording`, `postgres`) — Schicht-Abgrenzung: Der Slice
  ändert zwei Schichten, die Driving-Seite mit Composition Root (`pgwire`, `cli`,
  `bootstrap`, `cmd/pgwire-recorder`) und den Kern (`model`, `ports/driving`, `services`):
  Das Zwangsende ist ein Ereignis, das der Adapter meldet, seine Einstufung als
  `PGR-E4006` trifft der Use Case ([ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md), Architektur-Sicht §4.5).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-15`](../../../../spec/lastenheft.md#lh-fa-15--ci-eignung): Nach `SIGINT` oder `SIGTERM` wartet `record`
      höchstens `--shutdown-timeout` (ab dem ersten Signal, `0` ohne Frist) auf laufende
      Sessions; danach endet jede noch laufende Session wie bei einem Abbruch nach
      `LH-FA-02.b`, die Aufzeichnung wird geschrieben, das Log nennt je Session die
      verworfene Interaktion, der Prozess merkt `PGR-E4006` und endet mit Exit-Code 4; ein
      zweites Signal lässt die Frist sofort ablaufen; beim Beginn steht eine Info-Zeile mit
      der Zahl der Sessions (`sessions`); die Randfälle aus `LH-FA-13.a` und die Form der
      Dauer aus `LH-FA-17.a` gelten; `meldeFrist` hat den Wert von `SPEC-051` (Test mit
      Signal über die Helfer aus `slice-harness-integration-wait`; Unit-Test im Kern für die
      Einstufung; Test mit der Schranke als Literal); ein drittes Signal bleibt ohne
      Wirkung (Test mit drei Signalen).
- [ ] [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Die Frist begrenzt das Warten in `replay`; bei Ablauf
      schließt `replay` die Client-Verbindung, eine unvollständige Interaktion ist
      `PGR-E4006`, Exit-Code 4; mit `--fail-on-unconsumed` merkt die Session `PGR-E4006`
      vor `PGR-E5002` und der Prozess endet mit Exit-Code 4 (Test mit Signal; Test über
      `pgwire.Handle` mit Replay-Fake, V-95; Unit-Test im Kern für die Einstufung).
- [ ] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Das kontrollierte Herunterfahren liefert den Exit-Code nach `LH-FA-13.b`:
      `0` ohne zuvor aufgetretenen Verbindungsfehler, sonst die Klasse des ersten gemerkten
      (4, 5, 6, auch 1 für `PGR-E1000`), `3` bei einem Schreibfehler am Ende mit Vorrang;
      eine Klasse 2 entsteht nur als Startfehler und gehört nicht hierher. Beleg: die
      Zuordnung Code → Exit-Code als Unit-Test für jede Klasse aus `SPEC-013` bis `SPEC-019`
      (`SPEC-038`); nach Signal je Exit-Code 0, 3, 4, 5, 6 ein Integrationstest, die
      vorhandenen in §7 genannt, neu nur der mit `PGR-E4006`. Das Benutzerhandbuch
      beschreibt das Herunterfahren mit Frist in `record` und `replay`, auch `replay` im
      Container. Beleg in §7 für alle drei Punkte: je Zusage Zusage · Mutation · roter Test
      (`AGENTS.md` §3.10).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `cmd/pgwire-recorder` | update | nach dem ersten Signal bleibt die Signalbehandlung bestehen; das zweite beendet den Prozess nicht mehr, sondern lässt die Frist ablaufen (`LH-FA-13.a` *Zweites Signal*) |
| `internal/bootstrap` | update | Frist ab dem ersten Signal, zweites Signal, Info-Zeile mit `sessions`, erst Frist und zweites Signal bereit, dann Listener schließen (`LH-FA-13.a`), Exit-Code der gemerkten Klasse |
| `internal/adapters/driving/pgwire` | update | Zwangsende in `recordSitzung`, `replaySitzung` und im Verbindungsaufbau: Client-Verbindung schließen und das Ereignis dem Use Case melden; Lese- und Schreibfehler nach dem Zwangsende ohne `PGR-E4003`; Zustellung von `PGR-E4006` mit `meldeFrist` = `SPEC-051` auch im Replay; Zahl der offenen Verbindungen für `sessions` |
| `internal/hexagon/model` | update | Konstante `PGR-E4006` in der Code-Tabelle; Ereignis des Zwangsendes als `SessionEnd` |
| `internal/hexagon/ports/driving`, `internal/hexagon/services` | update | Record (`CloseSession`) und Replay stufen das Zwangsende ein: `PGR-E4006` genau bei begonnener, nicht abgeschlossener Interaktion, Meldung mit Session und Interaktion nach `LH-FA-13.a` *Meldung*; Unit-Tests der Einstufung |
| `internal/adapters/driving/pgwire` (Unit-Tests, V-95) | update | Ende einer Replay-Sitzung durch die Frist über `pgwire.Handle` mit Replay-Fake |
| `internal/adapters/driving/cli` | update | Option `--shutdown-timeout` und `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT` für `record` und `replay` (Optionstabelle, `SPEC-046`), Form der Dauer nach `LH-FA-17.a` *Dauer* |
| `internal/adapters/driving/pgwire` (Unit-Tests) | update | `meldeFrist` mit der Schranke `1s` als Literal (`SPEC-051`), rot bei `3 * time.Second` |
| `internal/bootstrap` (Unit-Tests) | update | Zuordnung Code → Exit-Code für jede Klasse aus `SPEC-013` bis `SPEC-019` |
| `test/integration` | update | Signal und Frist je Modus, Happy/Boundary/Negative nach `LH-FA-13.a`, über die Helfer aus `slice-harness-integration-wait`; vor dem zweiten Signal wartet der Test, bis der Port Verbindungen ablehnt; Log-Zeilen liest er nach dem Ende des Prozesses (§6 *`stderr` zur Laufzeit lesen*) |
| `docs/user/benutzerhandbuch.md` | update | Herunterfahren mit Frist in `record` und `replay`; `replay` im Container (heute `SIGKILL` nach der Stopp-Frist, Exit-Code `137`, keine Log-Zeile; Validierung zu `slice-extended-query-replay`, Frage 5) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-integration-wait` liegt in `done/` (die Tests mit Signal und
Frist entstehen über dessen Helfer). Schritt 2 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)
(Entscheidung des Nutzers vom 2026-10-08). Vor dem ersten Code-Commit prüft der Architect
die Randformen aus §6 und entscheidet die offenen (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Das Zwangsende verlangt eine
  Änderung in einem Driven-Adapter (`recording`, `postgres`), also eine dritte Schicht,
  oder DoD-Punkt 1 und 2 passen zusammen nicht in eine Review-Sitzung; dann geht die
  Frist in `replay` (DoD-Punkt 2) als eigener Slice der Welle heraus. Der Kern ist keine
  dritte Schicht: Die Einstufung des Zwangsendes liegt nach
  [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md) dort, und §3 führt ihn.
- `in-progress` → `open` (blockiert — Carveout?): Signalbehandlung ist im Container nicht
  prüfbar — Carveout.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — je Randform, wo sie entschieden ist. Der Architect hat
die Liste am 2026-10-08 vor dem Code geprüft; keine ist offen.

- **Frist `0`** — ohne Frist; entschieden in `LH-FA-13.a` und der Optionstabelle (`SPEC-046`).
- **Form des Werts** (Einheit, ohne Einheit, negativ, leer, Nachkommastellen, zusammengesetzt,
  Überlauf) — `LH-FA-17.a` *Dauer*: `0` oder ganze Zahl mit genau einer Einheit `ms`, `s`,
  `m`; sonst `PGR-E2001` (Option, Umgebungsvariable), `PGR-E2004` (Schlüssel der Datei, erst
  mit `slice-v1-abschluss-konfiguration`). Leere Umgebungsvariable gilt als nicht gesetzt,
  mehrfache Option: die letzte (`LH-FA-17.a`, allgemein).
- **Zweites Signal** — die Frist läuft sofort ab, auch beim Wert `0`; entschieden in
  `LH-FA-13.a` *Zweites Signal*. Vor dem ersten Signal gibt es kein zweites: Die Frist
  beginnt mit dem ersten.
- **Signal ohne offene Verbindung** (auch vor der ersten Verbindung) — der Prozess endet
  ohne Warten, `record` schreibt; Info-Zeile mit `sessions=0`; `LH-FA-13.a` *Ohne offene
  Verbindung* und Absatz zur Info-Zeile.
- **Signal in der Startphase vor dem Lauschen** — Startprüfungen laufen zu Ende, ein
  Startfehler geht vor, sonst wie ohne offene Verbindung; `LH-FA-13.a` *Startphase*. Ein
  Signal, bevor der Prozess Signale behandelt, ist eine Grenze (`LH-FA-13.a` *Grenze*).
- **Ablauf der Frist während des Schreibens der Aufzeichnung** — die Frist begrenzt nur das
  Warten auf die Verbindungen; ein laufendes Schreiben wird nicht abgebrochen, eine schon
  beendete Session endet nicht zwangsweise; `LH-FA-13.a` *Was die Frist begrenzt*.
- **Zwangsende im Verbindungsaufbau** (Startnachricht, Verbindungsaufbau zum Upstream) —
  die Verbindung endet ohne Meldung; `LH-FA-13.a` *Zwangsende*.
- **Lese- oder Schreibfehler, der das Schließen durch das Zwangsende bemerkt** — kein
  `PGR-E4003`; `LH-FA-13.a` *Zwangsende*.
- **Ende durch die Frist mitten in einer Interaktion** — im Record vor dem `ReadyForQuery`,
  im Replay eine begonnene, nicht verbrauchte (auch eine einfache, deren Antworten nicht
  gesendet sind): `PGR-E4006`, eingestuft vom Use Case; `LH-FA-13.a` *Zwangsende*; mit
  `--fail-on-unconsumed` `PGR-E4006` vor `PGR-E5002`, entschieden in `LH-FA-03.b`.
- **Inhalt und Zustellung der Meldung `PGR-E4006`** — Log-Zeile `error` je Verbindung,
  Zustellung an den Client nach `LH-FA-13.b` höchstens `SPEC-051` lang, Session und
  Interaktion je Modus; `LH-FA-13.a` *Meldung*. Der Text nach dem Kopf ist nicht Vertrag
  (`SPEC-034` *Stabilität*), sein Inhalt ist zugesagt und wird geprüft (`AGENTS.md` §3.11).
- **Info-Zeile: Text und Ort** — Stufe `info` auf `stderr` (`LH-FA-14.a`), Attribut
  `sessions` = angenommene, noch nicht beendete Verbindungen, auch `0`; Schlüssel und Wert
  sind Vertrag, `msg` nicht; `LH-FA-13.a`, `LH-FA-14.a` *Zeilenform*. Unter
  `--log-level warn` oder `error` steht sie nicht.
- **Exit-Code bei mehreren gemerkten Klassen** — die erste gemerkte, ein Schreibfehler am
  Ende (`3`) mit Vorrang, nie zugeordnete Sessions zuletzt; entschieden in `LH-FA-13.b`
  und `LH-FA-03.b` *Fehlerebene*.
- **Wert der Schreibfrist der Fehlerantwort beim Session-Ende** (seit slice-harness-blackbox-pgwire,
  Review F-451, Verifikation V-87; übernommen aus `slice-v1-abschluss-betrieb`) — `1s`,
  nicht einstellbar; entschieden in `SPEC-051` und `LH-FA-18.a` *Abbruch*, nach der
  Empfehlung der Validierung zu `slice-extended-query-record` (Frage 3). Sie gilt im Record
  für jedes Session-Ende und im Replay für `PGR-E4006`. Test mit der Schranke als Literal,
  rot bei `meldeFrist = 3 * time.Second`; `TestFehlerantwortMitFrist` liest heute die
  Konstante über die Brücke und fängt nur das Entfernen der Frist.
- **`stderr` zur Laufzeit lesen** (seit slice-harness-integration-wait, Verifikation V-105)
  — entschieden: kein Test liest `stderr`, solange der Prozess läuft; das Geschirr bleibt
  unverändert. Vor dem zweiten Signal wartet ein Test mit eigener Frist (`SPEC-038`
  *Warten in Tests*), bis der Port Verbindungen ablehnt: Dann hat der Prozess das erste
  Signal behandelt (`LH-FA-13.a` *Zweites Signal*). Info-Zeile und Meldungen prüft er nach
  dem Ende, wenn `Wait` die Ausgabe abgeschlossen hat. Ein synchronisierter Puffer im
  Geschirr wäre ein zweiter Lesepfad für ein Ereignis, das der Port schon zeigt.
- **Windows-Konsolenabbruch** (`Strg+C`, `Strg+Break`) — `LH-FA-13.a`; akzeptiertes
  Negativ: Die Tests laufen unter Linux im Container, Go liefert beide Ereignisse als
  `os.Interrupt` an dieselbe Behandlung; ein eigener Test entfällt.
- **Signale nach dem zweiten** — ohne Wirkung: Der Prozess führt das Zwangsende zu Ende,
  schreibt das Recording und endet mit dem Exit-Code nach `LH-FA-13.b`. Entscheidung des
  Nutzers vom 2026-10-08, festgehalten in `LH-FA-13.a` *Weitere Signale*. Mutation, die rot
  werden muss: Ein drittes Signal beendet den Prozess hart; der Test (Record mit laufender
  Interaktion, drei Signale, das dritte nach dem zweiten) erwartet danach die geschriebene
  Aufzeichnung mit den abgeschlossenen Interaktionen und den Exit-Code `4`.
  Grenze: Zwischen zweitem und drittem Signal zeigt der Prozess kein Ereignis; fasst das
  Betriebssystem beide zusammen, sieht der Test die Mutation in diesem Lauf nicht. Der
  Beleg in §7 nennt darum, in wie vielen Läufen die Mutation rot war.

**Risiken:**

- Herunterfahren ohne Obergrenze (aus `slice-extended-query-record`, Review F-309,
  Folge-Review F-317, Validierung Frage 1; übernommen aus `slice-v1-abschluss-betrieb`):
  Eine Session mit laufender Interaktion — einfache Anfrage oder Extended-Interaktion ohne
  `Sync` — hält das Herunterfahren beliebig lange; der Container-Stopp verliert dann ihre
  Aufzeichnung ganz, ohne Log-Zeile und mit Exit-Code `137` — **Ausgang:** offen bis
  Closure.
- Herunterfahren ohne Obergrenze im Replay (aus `slice-extended-query-replay`, Folge-Review
  F-330; übernommen aus `slice-v1-abschluss-betrieb`): Ein Client, der mitten in einer
  Extended-Interaktion pausiert, hält `replay` nach `SIGTERM` beliebig lange; endet die
  Session durch Schließen der Verbindung, behandelt `replaySitzung` das heute als
  reguläres Ende (Exit-Code 0 statt 4) — **Ausgang:** offen bis Closure.
- Ein Test mit Signal, dessen Prozess unter einer Mutation nicht endet, hängt bis zum
  Zeitlimit von `go test` (`BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`, 2×); jeder
  neue Test wartet darum über die Helfer aus `slice-harness-integration-wait` mit eigener
  Frist — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Belege zur DoD (Implementer):** <je Zusage: Zusage · Mutation · roter Test (`AGENTS.md` §3.10)>
- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `3864e44` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` (2×) — die Tests mit Signal; ein Risiko in
  §6. Hinge ein Test dieses Slice bis zum Zeitlimit, erreichte der Eintrag 3×.
- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (1×) — die Frist im Replay hängt an
  einem pausierenden Client; die Randformen in §6 nennen die Fälle.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (12×, verkörpert in `AGENTS.md`
  §3.12) — darum nennt §6 die Randformen, und der Architect prüft sie vor dem Code.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (14×, `AGENTS.md` §3.10) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (19×, §3.11) — je Zusage eine
  Mutation, Beleg in §7.
- `BEO-REPO/plan-folgt-korrektur-nicht` (17×, §3.9) — dieser Plan entsteht aus einem
  Schnitt; jede Korrektur zieht §1, §3 und §6 im selben Commit nach.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — jede Übernahme oben steht
  mit der Kennung des Gebers; die Geber zeigen im selben Commit hierher.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu; keine neue Lücke vor
dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
