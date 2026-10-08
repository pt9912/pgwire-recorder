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
  mehrfache Option: die letzte (`LH-FA-17.a`, allgemein). Führende Nullen (Rückgabe des
  Implementers, bestätigt vom Architect am 2026-10-08, `LH-FA-17.a` *Dauer*): vor einer
  Einheit erlaubt, ohne Einheit nur genau `0`. Fälle für `TestParseShutdownTimeoutWerte`:
  `05s` → 5 s, `00s` → `0`, `00` → `PGR-E2001`; Mutation, die rot werden muss: führende
  Nullen vor einer Einheit abgelehnt.
- **Zweites Signal** — die Frist läuft sofort ab, auch beim Wert `0`; entschieden in
  `LH-FA-13.a` *Zweites Signal*. Vor dem ersten Signal gibt es kein zweites: Die Frist
  beginnt mit dem ersten.
- **Signal ohne offene Verbindung** (auch vor der ersten Verbindung) — der Prozess endet
  ohne Warten, `record` schreibt; Info-Zeile mit `sessions=0`; `LH-FA-13.a` *Ohne offene
  Verbindung* und Absatz zur Info-Zeile.
- **Signal in der Startphase vor dem Lauschen** — Startprüfungen laufen zu Ende, ein
  Startfehler geht vor, sonst wie ohne offene Verbindung; `LH-FA-13.a` *Startphase*. Ein
  Signal, bevor der Prozess Signale behandelt, ist eine Grenze (`LH-FA-13.a` *Grenze*).
  Eine Verbindung, die das Betriebssystem zwischen `Listen` und dem Schließen des Ports
  schon angenommen hat, kann `Accept` noch liefern (Review F-489): Grenze, entschieden vom
  Architect am 2026-10-08 in `LH-FA-13.a` *Startphase*. Die Verbindung zählt in `sessions`
  und endet ohne Interaktion. Eine Prüfung von `ctx` vor `Serve` verkleinerte das Fenster
  nur, schlösse es nicht (Signal zwischen Prüfung und `Accept`), und kein Test fängt sie
  wettlauffrei (Mutant G01, Probe 0 von 400); kein Code dafür. G01 in §7 ist damit als
  *Fenster offen, als Grenze entschieden* zu lesen, nicht als *nicht erreichbar*.
- **Ablauf der Frist während des Schreibens der Aufzeichnung** — die Frist begrenzt nur das
  Warten auf die Verbindungen; ein laufendes Schreiben wird nicht abgebrochen, eine schon
  beendete Session endet nicht zwangsweise; `LH-FA-13.a` *Was die Frist begrenzt*.
- **Zwangsende im Verbindungsaufbau** (Startnachricht, Verbindungsaufbau zum Upstream) —
  die Verbindung endet ohne Meldung; `LH-FA-13.a` *Zwangsende*. Wartet der Aufbau noch auf
  die Antwort des Upstreams, schließt der Adapter die Client-Verbindung und bricht den
  Kontext ab; der Postgres-Adapter liest ohne Kontext weiter, und die Upstream-Verbindung
  schließt erst das Prozessende. Rückgabe des Implementers, entschieden vom Architect am
  2026-10-08 als Grenze in `LH-FA-13.a` *Zwangsende*: Das Prozessende folgt dem Zwangsende
  unmittelbar nach dem Schreiben, die Verbindung trägt keine Interaktion, und der Upstream
  sieht dasselbe Verbindungsende wie bei einem Schließen. Ein Schließen im Driven-Adapter
  wäre eine dritte Schicht (§1, §4) ohne Unterschied, den ein Client oder die Aufzeichnung
  sieht; kein Folge-Slice.
- **Lese- oder Schreibfehler, der das Schließen durch das Zwangsende bemerkt** — kein
  `PGR-E4003`; `LH-FA-13.a` *Zwangsende*.
- **Ende durch die Frist mitten in einer Interaktion** — im Record vor dem `ReadyForQuery`,
  im Replay eine begonnene, nicht verbrauchte (auch eine einfache, deren Antworten nicht
  gesendet sind): `PGR-E4006`, eingestuft vom Use Case; `LH-FA-13.a` *Zwangsende*; mit
  `--fail-on-unconsumed` `PGR-E4006` vor `PGR-E5002`, entschieden in `LH-FA-03.b`.
- **Inhalt und Zustellung der Meldung `PGR-E4006`** — Log-Zeile `error` je Verbindung,
  Zustellung an den Client nach `LH-FA-13.b` höchstens `SPEC-051` lang, Session und
  Interaktion je Modus; `LH-FA-13.a` *Meldung*. Grund einer nicht geschriebenen Session im
  Record (Review F-487, entschieden vom Architect am 2026-10-08): ohne abgeschlossene
  Interaktion, oder wegen einer nicht unterstützten Interaktion (`PGR-E6001`), auch mit
  abgeschlossenen; eine dritte Form gibt es nicht. Vorgabe: Unit-Test in
  `internal/hexagon/services` ohne Wettlauf — die Session hat abgeschlossene Interaktionen,
  ist als nicht unterstützt markiert und hat eine laufende Interaktion; `CloseSession` mit
  `EndForced` liefert `PGR-E4006` mit dem Grund *nicht unterstützte Interaktion*. Mutant M1
  (Text des Falls ohne abgeschlossene Interaktion) muss rot werden. Der Text nach dem Kopf ist nicht Vertrag
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
- **Mutant X03** (Ziffer 1 aus dem Bereich in `model.exitCode` genommen; Rückgabe des
  Implementers) — äquivalent, akzeptiertes Negativ, entschieden vom Architect am
  2026-10-08: Ein Code der Klasse 1 (`PGR-E1xxx`) und ein Code ohne Klasse fallen nach
  `SPEC-034` *Fehler* beide auf Exit-Code 1, Ziffer und Rückfall liefern dasselbe. Der Code
  bleibt so; ein Test kann den Mutanten nicht unterscheiden.
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

- **Belege zur DoD (Implementer):** siehe *Belege des Implementers* unten.
- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

**Belege des Implementers** (Commits `17e4925` DoD 1, `64f48b3` DoD 2, `61ad648` DoD 3,
`5e20a87` Integrationstest ohne `pgproto3` nach dem Befund von `make a-check`):

- **Läufe am Stand `5e20a87`:** `make gates` grün (Exit 0; darin `make test-integration`
  grün, `make a-check` 0 Befunde, `make docs-check` 0 Befunde, `make abdeckung-check`
  grün); `make lint` 0 Befunde; `make test` grün. Am Stand `61ad648` war `make a-check`
  rot (`tech-leak`: `pgproto3` in `test/integration/herunterfahren_e2e_test.go`), behoben
  in `5e20a87`; R38 dort erneut rot gesehen.
- **Wiederholungen:** die neuen Unit-Tests (`-run 'Zwang|Frist|Offen|Startphase|ZweitesSignal|OhneOffene|Herunterfahren|LogLevel'`
  in `pgwire`, `bootstrap`, `services`) 20 Läufe in einem Wegwerf-Container ohne Netz,
  alle grün; `TestE2ERecordFristLaeuftAb`, `TestE2ERecordZweitesSignal`,
  `TestE2ERecordDrittesSignal` je 5 Läufe, `TestE2EReplayFristLaeuftAb`,
  `TestE2EReplayFristVorNichtVerbraucht` je 3 Läufe, alle grün.
- **Weg der Mutanten:** je Mutant eine frische Kopie des Arbeitsbaums unter eigenem Pfad
  außerhalb des Repos (`cp -r` ohne `-p`, also neue mtime, neuer Pfad je Mutant), genau eine
  Ersetzung per Skript, die vorher genau einmal gefunden sein muss; Unit-Mutanten über die
  Stufe `source` der Kopie und `go test -run` in einem Container mit `--network none`,
  Integrations-Mutanten über die Stufe `integration` der Kopie mit eigenem Image-Tag,
  eigenem `--internal`-Netz, eigenem PostgreSQL-Container (Image wie
  `make test-integration`) und eigenem Volume; danach Container, Netz, Volume, Image und
  Kopie entfernt. Kein Mutant im Repo. Stand der Mutanten: Kopien des Arbeitsbaums vor dem
  jeweiligen Commit (R01 bis R42 vor `17e4925`, P01 bis P19 vor `64f48b3`, X01 bis X05 vor
  `61ad648`); was sich danach bis zum Commit änderte (Entfernen der Prüfung aus G01, Kommentar
  am Port, `liestWieder` im Test zu P10), berührt die Zeilen der übrigen Mutanten nicht; P07,
  P08 und P10 liefen nach `liestWieder` erneut rot.
- **Zusage · Mutation · roter Test**, DoD 1 (`record`):

  | ID | Zusage | Mutation | roter Test |
  |---|---|---|---|
  | R01 | `PGR-E4006` nur bei laufender Interaktion | Bedingung `l.laeuft()` entfernt | `TestRecordZwangsende` (ohne laufende Interaktion) |
  | R02 | Meldung nennt die Nummer der verworfenen Interaktion | Nummer `+2` statt `+1` | `TestRecordZwangsende`, `TestRecordZwangsendeEinfacheAnfrage` |
  | R03 | Meldung nennt die `id`, unter der die Session geschrieben wird | `session.ID+1` | `TestRecordZwangsende`, `TestRecordZwangsendeEinfacheAnfrage` |
  | R04 | Klasse 4 (`PGR-E4006`) | `CodeConnectionLost` | `TestRecordZwangsende`, `TestRecordZwangsendeEinfacheAnfrage` |
  | R05 | Meldung: ohne abgeschlossene Interaktion nicht geschrieben | Fall `n == 0` nie | `TestRecordZwangsende` (ohne abgeschlossene Interaktion) |
  | R06 | Code `PGR-E4006` (`SPEC-034`) | Konstante `PGR-E4007` | `TestRecordZwangsende` |
  | R07 | `PGR-E4006` wird dem Client zugestellt | Zustellung abgeschaltet | `TestRecordZwangsendeZustellung` |
  | R08 | Zustellung höchstens `meldeFrist`, auch an einen Client, der nicht liest | Schreibfrist entfernt | `TestRecordZwangsendeClientLiestNicht` |
  | R09 | eine Verbindung, deren Ende schon bemerkt ist, endet nicht zwangsweise | `CloseSession(EndForced)` vor dem `Once` | `TestRecordZwangsendeNachBemerktemEnde` |
  | R10 | Schließen durch das Zwangsende ist kein weiteres Verbindungsende | Verbindung vor dem `Once` geschlossen | `TestRecordZwangsendeOhneMeldung` |
  | R11 | Zwangsende wartet nicht auf den Aufbau zum Upstream | Zweig `zwangCh` in `oeffne` nie | `TestRecordZwangsendeImAufbau` |
  | R12 | eine nach dem Zwangsende entstandene Session endet mit `EndForced` | Nachzug nie | `TestRecordZwangsendeImAufbau` |
  | R13 | Schreibfehler des Aufbaus nach dem Zwangsende ist `EndForced`, nicht `EndWriteFailed` | Prüfung von `zwang` umgangen | `TestRecordZwangsendeBeimSendenDesAufbaus` |
  | R14 | Fehlerantwort des Servers im Aufbau nach dem Zwangsende: kein `PGR-E4003` | Prüfung von `zwang` entfernt | `TestRecordZwangsendeBeimSendenDesAufbaus` |
  | R15 | Sitzung kehrt erst zurück, wenn `PGR-E4006` gemerkt ist | `<-r.ende` entfernt | `TestRecordZwangsendeWaehrendAnfrage` |
  | R16 | `sessions` zählt nur nicht beendete Verbindungen | kein Abzug am Ende | `TestOffeneVerbindungen` |
  | R17 | `meldeFrist` = 1 s (`SPEC-051`), Schranke als Literal | `3 * time.Second` | `TestFehlerantwortMitFrist` |
  | R18 | Frist begrenzt das Warten | Zeitgeber nie gesetzt | `TestRunRecordFristLaeuftAb` |
  | R19 | `0` heißt ohne Frist | Zeitgeber auch bei `0` | `TestRunRecordZweitesSignalOhneFrist` |
  | R20 | zweites Signal lässt die Frist sofort ablaufen | `ablauf` nie gelesen | `TestRunRecordZweitesSignalOhneFrist`, `TestRunRecordZweitesSignalLangeFrist` |
  | R21 | ohne offene Verbindung endet der Lauf ohne zu warten | nach `done` auf die Frist gewartet | `TestRunRecordOhneOffeneVerbindung` |
  | R22 | Info-Zeile nennt `sessions` = offene Verbindungen | `sessions` fest `0` | `TestRunRecordFristLaeuftAb` |
  | R23 | Info-Zeile auf Stufe `info` | als `debug` | `TestRunRecordLogLevel` |
  | R24 | Info-Zeile nicht unter `warn` | als `warn` | `TestRunRecordLogLevel` |
  | R25 | Startphase: Startprüfungen laufen zu Ende, Startfehler geht vor | bei beendetem ctx sofort 0 | `TestRunRecordSignalInDerStartphase` |
  | R26 | `--shutdown-timeout` wirkt bei `record` | Frist `0` statt der Option | `TestRunRecordFristLaeuftAb` |
  | R27 | nur Einheiten `ms`, `s`, `m` | `h` zugelassen | `TestParseShutdownTimeoutWerte` |
  | R28 | Überlauf ist `PGR-E2001` | Grenze ohne Einheit | `TestParseShutdownTimeoutWerte` |
  | R29 | `0` ohne Einheit ist `0` | nur `00` | `TestParseShutdownTimeout`, `…Werte`, `…Umgebung` |
  | R30 | Standard `5s` (`SPEC-046`) | `4s` | `TestParseShutdownTimeout`, `…Umgebung` |
  | R31 | Umgebungsvariable gilt, ungültig ist `PGR-E2001` auch neben der Option | Umgebung nie gelesen | `TestParseShutdownTimeoutUmgebung`, `…UmgebungUngueltig` |
  | R32 | Kommandozeile setzt den Wert, die letzte Angabe gilt | Option schreibt ins Leere | `TestParseShutdownTimeout`, `…Umgebung` |
  | R33 | Einheit `m` ist Minute | `m` als Sekunde | `TestParseShutdownTimeout`, `…Umgebung` |
  | R34 | drittes Signal ohne Wirkung | `signal.Stop` nach dem zweiten | `TestE2ERecordDrittesSignal` — 3 von 3 Läufen rot (Exit-Code −1, `signal: terminated`) |
  | R35 | zweites Signal beendet den Prozess nicht hart | `signal.Stop` nach dem ersten | `TestE2ERecordZweitesSignal` (Exit-Code −1) |
  | R36 | zweites Signal lässt die Frist ablaufen (Binary) | `ablauf` nie geschlossen | `TestE2ERecordZweitesSignal` (endet nicht binnen 15 s) |
  | R37 | Frist begrenzt das Warten (Binary) | Frist `0` statt der Option | `TestE2ERecordFristLaeuftAb` (endet nicht binnen 15 s) |
  | R38 | Client erhält `PGR-E4006` (Binary) | Zustellung abgeschaltet | `TestE2ERecordFristLaeuftAb` |
  | R39 | laufende Interaktion fehlt, Exit-Code 4 (Binary) | Bedingung umgekehrt | `TestE2ERecordFristLaeuftAb` (Exit-Code 0) |
  | R40 | Verbindungen enden beim Zwangsende nebeneinander | `ausloesen` nacheinander | `TestZwangsendeNebeneinander` (3,01 s) |
  | R41 | ein weiterer Aufruf von `Zwangsende` tut nichts | Rücksprung entfernt | `TestZwangsendeZweimal` |
  | R42 | Zwangsende bricht den Kontext des Aufbaus ab | `abbrechen()` im Zweig `zwangCh` entfernt | `TestRecordZwangsendeBrichtAufbauAb` |

- **DoD 2 (`replay`):**

  | ID | Zusage | Mutation | roter Test |
  |---|---|---|---|
  | P01 | `sequence` der nicht gesendeten Interaktion | anderer Index | `TestReplayZwangsende` (Kern) |
  | P02 | Extended vor dem letzten `Sync` ist `PGR-E4006` | Fall `mitten` nie | `TestReplayZwangsende`, `TestReplayZwangsendeNichtVerbraucht` |
  | P03 | nicht gesendete Antworten sind `PGR-E4006` | Fall `ungesendet` nie | `TestReplayZwangsende` |
  | P04 | Klasse 4 | `CodeConnectionLost` | `TestReplayZwangsende`, `TestReplayZwangsendeNichtVerbraucht` |
  | P05 | Meldung nennt die `id` der Session | `ID+1` | `TestReplayZwangsende` |
  | P06 | ohne begonnene Interaktion kein Fehler | immer Fehler | `TestReplayZwangsende` |
  | P07 | Adapter meldet das Zwangsende dem Use Case (V-95) | `Forced` nicht gerufen | `TestReplayZwangsende` (Adapter) |
  | P08 | Adapter stellt `PGR-E4006` zu | nur gemerkt | `TestReplayZwangsende` (Adapter) |
  | P09 | Schreibfrist beim Zwangsende | entfernt | `TestReplayZwangsendeClientLiestNicht` |
  | P10 | Lesefrist beim Zwangsende | entfernt | `TestReplayZwangsende`, `…VorNichtVerbraucht` (zuerst grün, siehe unten) |
  | P11 | Schreibfehler nach dem Zwangsende ist kein `PGR-E4003` | Prüfung entfernt | `TestReplayZwangsendeClientLiestNicht` |
  | P12 | Lesefehler nach dem Zwangsende führt zu `Forced` | Prüfung entfernt | `TestReplayZwangsendeOhneMeldung` |
  | P13 | Zurücksetzen der Lesefrist hebt das Zwangsende nicht auf | Nachprüfung entfernt | `TestReplayZwangsendeWaehrendShutdown` |
  | P14 | Schreibfehler des Aufbaus nach dem Zwangsende: kein `PGR-E4003` | Prüfung entfernt | `TestReplayZwangsendeBeimSendenDesAufbaus` |
  | P15 | abgebrochenes Senden führt zu `Forced` | als reguläres Ende | `TestReplayZwangsendeClientLiestNicht` |
  | P16 | `replay` nimmt `--shutdown-timeout` und die Umgebungsvariable | Option nicht angemeldet | `TestParseShutdownTimeout`, `…Umgebung`, `…UmgebungUngueltig` |
  | P17 | Info-Zeile auch bei `replay` | Logger verworfen | `TestRunLogLevel` |
  | P18 | Frist begrenzt das Warten (Binary) | Frist `0` | `TestE2EReplayFristLaeuftAb` (endet nicht binnen 15 s) |
  | P19 | `PGR-E4006` vor `PGR-E5002`, Exit-Code 4 (Binary) | `Forced` nach `CloseConnection` | `TestE2EReplayFristLaeuftAb` (Exit-Code 0), `TestE2EReplayFristVorNichtVerbraucht` (Exit-Code 5) |

- **DoD 3 (Exit-Code):**

  | ID | Zusage | Mutation | roter Test |
  |---|---|---|---|
  | X01 | ohne gemerkten Fehler Exit-Code 0 (`SPEC-013`) | `1` | `TestExitCodeJeKlasse` |
  | X02 | Klasse 6 (`SPEC-019`) | Ziffer 6 außerhalb | `TestExitCodeJeKlasse` |
  | X04 | Schreibfehler am Ende: 3 mit Vorrang | Vorrang nur ohne gemerkten Fehler | `TestRunRecordSchreibfehlerVorGemerkterKlasse` (Exit-Code 4) |
  | X05 | nach dem Signal Exit-Code der gemerkten Klasse | gemerkter Code verworfen | `TestRunRecordSchreibfehlerVorGemerkterKlasse` (Exit-Code 0) |

  Nach Signal je Exit-Code ein Integrationstest, vorhanden: 0 `TestE2ERecordEndeOhneTerminate`,
  3 `TestE2ERecordSchreibfehlerAmEnde`, 4 `TestE2ERecordUpstreamNichtErreichbar`
  (`PGR-E4002`) und `TestE2ERecordExtendedAbbruch` (`PGR-E4003`), 5
  `TestE2EReplayAbweichung` (`PGR-E5001`) und `TestE2EReplayNichtVerbraucht`
  (`PGR-E5002`), 6 `TestE2ERecordNichtUnterstuetzt`; neu mit `PGR-E4006`:
  `TestE2ERecordFristLaeuftAb`, `TestE2EReplayFristLaeuftAb`. 1 (`PGR-E1000`) entsteht
  über das Binary nicht steuerbar und ist nur im Unit-Test belegt; 2 entsteht nur als
  Startfehler.
- **Grüne Mutanten und ihre Einordnung:**
  - G01 (die Prüfung `ctx.Err()` vor `Serve` in `betreiben` umgangen): grün. Sie änderte das
    Verhalten nur zwischen `Listen` und dem sofortigen `Close`, über die Schnittstelle nicht
    erreichbar; die Prüfung ist vor `17e4925` entfernt, `Serve` läuft ab dem Lauschen, und
    der Listener schließt mit dem ersten Signal (`TestRunRecordSignalInDerStartphase`).
  - `context.WithoutCancel` beim Anlegen von Record- und Replay-Service und eine Prüfung von
    `zwang` am Anfang von `replayLauf`: äquivalent (`Prepare` und `Load` beachten den
    Kontext nicht; den Fall am Schleifenanfang decken der Lesepfad und die Nachprüfung nach
    dem Zurücksetzen); beide vor `64f48b3` entfernt, kein Code ohne Test.
  - P10 war zuerst grün: `TestReplayZwangsende` rief `Zwangsende`, bevor die Sitzung wieder
    las. Der Test wartet seit `64f48b3` (`liestWieder`), bis die Sitzung den Use Case gefragt
    hat; danach rot.
  - X03 (Ziffer 1 aus dem Bereich genommen, `model.exitCode`): grün und äquivalent, der
    Rückfall liefert für die Ziffer 1 ebenfalls 1. Bestand, in diesem Slice unverändert.
- **Lesart der Dauer:** die Werte in `TestParseShutdownTimeout` und
  `TestParseShutdownTimeoutWerte` folgen dem Wortlaut von `LH-FA-17.a` *Dauer*; `00` ohne
  Einheit ist nicht `0` und abgelehnt, führende Nullen mit Einheit (`05s`) nimmt der Code an,
  ohne dass ein Test es festhält (im Bericht an den Architect zur Bestätigung).
- **`stderr` zur Laufzeit:** kein neuer Test liest `stderr`, solange der Prozess läuft; die
  Log-Zeilen prüfen die Integrationstests nach `warteEnde`, die Bootstrap-Tests nach dem Ende
  von `Run`; vor dem zweiten Signal wartet `TestE2ERecordZweitesSignal` und
  `TestE2ERecordDrittesSignal` mit Frist, bis der Port ablehnt (`warteAbgelehnt`).

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
