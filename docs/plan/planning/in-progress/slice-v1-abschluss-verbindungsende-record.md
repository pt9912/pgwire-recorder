# Slice slice-v1-abschluss-verbindungsende-record: Verbindungsende im Record: Use Case und Zustellung an den Client

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-02.b` · `LH-FA-12.a` · `LH-FA-13.a` · `LH-FA-18.a` · `LH-FA-21.a` · `SPEC-001` · `SPEC-002` · `ARC-002` · `ARC-006`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Record-Service übernimmt eine Interaktion, deren Antworten ein Ende des Servers zeigen (`connection_closed` als letzte Antwort, aus dem Upstream-Port), als letzte Interaktion der Session — bei einer einfachen Anfrage wie bei einer Extended-Interaktion — und der PGWire-Adapter stellt dem Client die Antworten bis einschließlich der Fehlerantwort zu und schließt dann beide Verbindungen, ohne `ReadyForQuery`, ohne Meldung und mit unverändertem Exit-Code (`LH-FA-25.a` *Aufzeichnen*, *Form*, *Stelle*). Ein Ende in einer Gruppe, die der Client nicht mit `Sync` oder `Flush` abgeschlossen hat, ist `PGR-E6001` (die Session wird nicht übernommen). Nachweis mit Fake-Upstream und Fake-Recorder; der echte Upstream und Abnahmeszenario 18 folgen in `slice-v1-abschluss-verbindungsende-upstream`.

**Übernommen aus `slice-v1-abschluss-verbindungsende-format`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): das Aufzeichnen des Endes im Use Case und im PGWire-Adapter. Jener Slice liefert `connection_closed` und seine Formregeln im Modell (`Validate` nimmt die Interaktion mit Ende an); dieser benutzt sie.

**Bindung aus `slice-v1-abschluss-verbindungsende-format`** (Closure des Gebers, Planner, 2026-10-15, `AGENTS.md` §3.13; dort §6 *Ort der Formregel einer Session*): Die Formregel einer Session (höchstens ein Ende, in der letzten Interaktion) liegt als `Session.ValidateEnd() error` im Modell, damit dieser Slice dieselbe Regel aufruft statt sie nachzubauen; `Interaction.Validate` allein nimmt ein Ende an jeder Interaktion an. Ob der Use Case sie beim Schreiben der Session ruft oder die Reihenfolge der Interaktionen sie von selbst erfüllt, entscheidet der Architect vor dem Code; ein eigener Nachbau der Regel im Service ist es nicht.

**Bindung aus `slice-v1-abschluss-verbindungsende-format-schreiber`** (Architect, 2026-10-15, `AGENTS.md` §3.13; dort §6 *Aufrufer mit verletzter Form*, `SPEC-001` *Schreiber in `yaml`*; für `sqlite` bestätigt vom Architect am 2026-10-15 mit `slice-v1-abschluss-verbindungsende-sqlite-schreiber`, `SPEC-001` *Schreiber in `sqlite`*): Der Schreiber prüft die Form des Endes nicht, in `yaml` wie in `sqlite`; er schreibt eine verletzte Form als Datei, die der Leser mit `PGR-E3003` ablehnt. Der Use Case garantiert darum, dass keine Session mit verletzter Form (Ende nicht in der letzten Interaktion, zweites Ende) in die Aufzeichnung gelangt — durch den Aufruf von `Session.ValidateEnd` oder durch eine Reihenfolge, die sie erfüllt, mit einem Test samt Mutation; welcher Weg, entscheidet der Architect vor dem Code (siehe die Bindung aus `-format` darüber).

**Übernommen aus `slice-v1-abschluss-verbindungsende-replay`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12): der Record-Pfad des PGWire-Adapters. Die Abbildung von `connection_closed` auf „keine PGWire-Nachricht“ (`toMessage`, `send`) liefert jener Slice; dieser ruft sie auf und schließt in `richtungen` beide Verbindungen.

**Bindung aus `slice-v1-abschluss-verbindungsende-replay`** (Closure des Gebers, Planner, 2026-10-16, `AGENTS.md` §3.13; dort §6, Risiko *aufgeteilte Abbildung*, Entscheidung des Architect vom 2026-10-16): `send` im PGWire-Adapter lässt `connection_closed` als letzte Antwort aus und schließt nicht selbst; das Schließen steht beim Aufrufer (Replay: `replayZustellen`, nach `Sent`; Record: `richtungen`). Dieser Slice ändert diese Aufteilung nicht; ändert er sie, zieht er den Replay-Pfad und dessen Tests (`TestReplayEndeSendetDavorUndSchliesst`, `TestReplayEndeNurAlsLetzteAntwort`, `TestReplayEndeSentVorSchliessen`) im selben Commit nach. Der Record-Pfad hat für die Auslassung keinen eigenen Test (Review F-7 des Gebers): Liefert der Use Case `connection_closed` an `send`, trägt dieser Slice den Fall im Record-Pfad mit einer Mutation; liefert er sie nicht, nennt §7, dass `send` sie im Record-Pfad nie erhält. Entschieden (Architect, 2026-10-16): Er liefert sie (`LH-FA-25.a` *Meldung des Endes beim Aufzeichnen*), der Fall mit Mutation gehört in diesen Slice.

**Übernommen aus `slice-v1-abschluss-verbindungsende-replay-extended-flush`** (Verifikation V-2 dort, Sendung des Planners, 2026-10-16, `AGENTS.md` §3.13): Der Abschluss einer Extended-Interaktion ist ihr `Sync`, bei einer Interaktion mit Ende, deren letzte Gruppe mit `Flush` schließt, dieses `Flush` (`LH-FA-18.a` *Interaktion*, `LH-FA-13.a` „bis zu ihrem Abschluss“, Architect, `7dcbd12`). Die Kommentare im Record-Pfad sagen noch „Sync“: `internal/hexagon/services/record.go` (Feld `offen`: „solange ihr Sync aussteht“; `Query`: „vor deren Sync“; `ClientMessage`: „vor deren Sync“) und `internal/hexagon/ports/driving/record.go` (`Query`: „deren Sync schon gesendet ist“). Dieser Slice setzt sie auf „Abschluss“ (Kommentar, keine Logik), weil er in diesen Pfaden die Interaktion mit Ende und `Flush` einführt (DoD-Punkt 2). Zusätzlich entscheidet der Architect vor dem Code (§6), ob mit einem Ende des Servers beim Herunterfahren etwas zu entscheiden ist: Die Spezifikation sagt in `LH-FA-18.a` *Abbruch*, eine Session ende beim kontrollierten Herunterfahren „erst nach dem `ReadyForQuery` ihrer laufenden Extended-Interaktion, auch wenn deren `Sync` noch aussteht (`LH-FA-13.a`)“; eine Interaktion mit Ende hat kein `ReadyForQuery`. Änderungen an der Spezifikation macht der Architect, nicht dieser Slice.

**Übernommen aus `slice-v1-abschluss-sessions`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12): Das Ende des Servers als Ereignis einer Session. `slice-v1-abschluss-sessions` regelt Verbindungsende und Verbindungsfehler ohne das Ende des Servers; die Reaktion des Use Case auf ein Ende nach `LH-FA-25.a` liegt hier (die Einstufung von `PGR-E4003` und `PGR-E6001` im Adapter in `slice-v1-abschluss-verbindungsende-upstream`), und `CloseSession`, `--record-empty-sessions` und die Kennung der Session (Nummer in der Reihenfolge der Sessionenden) sind nach jenem Slice (Schritt 26) unverändert zu benutzen.

**Gebunden durch `slice-v1-abschluss-sessions`** (dort §1, *Nahtstelle zu den Verbindungsende-Slices*, Punkt 3 und 4; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): Dieser Slice hält ein, was `slice-v1-abschluss-sessions` im Record-Service liefert:

3. *`CloseSession` entscheidet allein, ob und unter welcher Kennung eine Session geschrieben wird* (Kennung = Zahl der geschriebenen Sessions + 1, unter dem Sperrzustand des Services). Dieser Slice fügt die Art des Endes (Ende des Servers) hinzu und ändert nicht, dass eine Session mit nicht unterstützter Interaktion und das Zwangsende im Aufbau (`EndForced` ohne laufende und ohne abgeschlossene Interaktion, `wirdVerworfen` in `internal/hexagon/services/record.go`) nie geschrieben werden und jede andere Session mit `--record-empty-sessions` geschrieben wird. Eine Session mit Ende trägt mindestens die Interaktion mit dem Ende und ist deshalb nie eine leere Session.
4. *Eine Verbindung wird erst mit eröffneter Session eine Session:* Ein Ende des Servers vor dem ersten `ReadyForQuery` (`LH-FA-25.a` *Kein Ende*) bleibt ein Aufbaufehler und erzeugt keine Session, auch mit `--record-empty-sessions`; `OpenSession` liefert dann die Kennung 0.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Upstream-Adapter (Erkennen der Fehlerantwort `FATAL`/`PANIC`, nicht auf das Verbindungsende warten, Ende im Leerlauf mit nachfolgender Anfrage, Lesestand bei Sendefehler, Einstufung `PGR-E4003`/`PGR-E6001` im Adapter) und Abnahmeszenario 18 — `slice-v1-abschluss-verbindungsende-upstream` (dort §1 und DoD): andere Schicht (Upstream-Adapter); mit ihm läge dieser Plan über zwei Schichten. Dieser Slice behandelt, was der Port liefert.
- SQLite-Aufzeichnung eines echten Endes — `slice-v1-abschluss-verbindungsende-upstream` (E2E); der Schreiber liegt in `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Schnitt von `slice-v1-abschluss-verbindungsende-sqlite`, 2026-10-15).
- Die Wiedergabe im Replay — `slice-v1-abschluss-verbindungsende-replay`.
- Ein Netzwerkabbruch ohne Fehlerantwort, ein Ende im Verbindungsaufbau und `FATAL` im Aufbau (etwa eine abgelehnte Anmeldung) als Ende — Out-of-Scope von `LH-FA-25`; sie bleiben wie bisher (`PGR-E4003`, `PGR-E4002`, `PGR-E4005`; `slice-v1-abschluss-anmeldung`, Schritt 30, bleibt unberührt).
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus): Der Record-Service übernimmt eine einfache Interaktion, deren Antworten mit `error_response` `FATAL`/`PANIC` und `connection_closed` enden, als letzte Interaktion der Session und gibt dem Aufrufer die Antworten der Interaktion mit `connection_closed` als letzter Antwort zurück (Architect, 2026-10-16, `LH-FA-25.a` *Meldung des Endes beim Aufzeichnen*; der Client erhält davon die Antworten bis einschließlich der Fehlerantwort); die Interaktion ist abgeschlossen, `CloseSession` meldet für das Ende der Verbindung weder `PGR-E4003` noch `PGR-E6001`; eine Nachricht nach dem Ende gehört nicht zur Aufzeichnung; mit `--record-timing` trägt die Interaktion `offset_ms`, mit `--record-empty-sessions` ändert das Ende nichts; `wirdVerworfen` und die Entscheidung in `CloseSession` bleiben die eine Stelle, die über das Schreiben entscheidet (Bindung aus `slice-v1-abschluss-sessions`; Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Bei einer Extended-Interaktion steht das Ende in der letzten Gruppe, die mit `Sync` oder mit `Flush` abgeschlossen ist; Interaktionen, die der Client nach dieser Gruppe in die Pipeline gestellt hat, gehören nicht zur Aufzeichnung; ein Ende in einer nicht abgeschlossenen Gruppe ist `PGR-E6001`, die Session wird nicht übernommen, auch nicht ihre vorherigen Interaktionen; `ERROR` mit anschließendem Abbruch (kein Ende) bleibt `PGR-E6001` (Test). Die Kommentare im Record-Pfad nennen den Abschluss statt des `Sync` (`record.go` dreimal, `ports/driving/record.go` einmal; Sendung aus `slice-v1-abschluss-verbindungsende-replay-extended-flush`, §1), und `grep -n "Sync"` über beide Dateien findet danach keinen Kommentar, der den Abschluss einer Interaktion mit Ende meint.
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Der PGWire-Adapter schreibt die Antworten bis einschließlich der Fehlerantwort an den Client und schließt dann Client- und Upstream-Verbindung (Session-Ende über den Use Case), ohne `ReadyForQuery`, ohne Meldung, ohne den Exit-Code zu berühren; ein Schreibfehler davor bleibt `EndWriteFailed` (Test mit Fake-Recorder); die Aufteilung `send` lässt `connection_closed` aus, der Aufrufer schließt, bleibt (Bindung aus `slice-v1-abschluss-verbindungsende-replay`), und der Record-Pfad hat dafür einen Fall mit Mutation (Architect, 2026-10-16: der Use Case liefert `connection_closed` an `send`).
- [ ] Je Zusage des neuen Vertrags eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services/record.go` | update | Übernahme einer Interaktion mit Ende (`Query`, `AwaitServer`, `uebernehmen`), Session als letzte, Stelle des Endes in der Gruppe, `CloseSession`; Kommentare „Sync“ auf „Abschluss“ (Feld `offen`, `Query`, `ClientMessage`; Sendung aus `slice-v1-abschluss-verbindungsende-replay-extended-flush`) |
| `internal/hexagon/ports/driving/record.go` | update | Port-Vertrag: `connection_closed` als letzte Antwort von `Query` und `AwaitServer` (Ende des Servers, Adapter sendet davor, meldet `Delivered`, beendet mit `EndServer`); Kommentar von `Query` „Sync“ auf „Abschluss“ (Sendung aus `slice-v1-abschluss-verbindungsende-replay-extended-flush`). `ports/driven/upstream.go` bleibt unverändert: der Upstream-Port liefert `connection_closed` schon in den Antworten |
| `internal/hexagon/model/recording.go` | update | neues Ereignis `EndServer` (kein Fehler, keine Meldung); keine Logik |
| `internal/adapters/driving/pgwire/server.go` | update | Record-Pfad: `richtungen.schreibe` (nach `Delivered` `EndServer`, beide Richtungen beenden); `serverRichtung` und `clientRichtung` bleiben, sie rufen `schreibe` |
| `internal/hexagon/services`, `internal/adapters/driving/pgwire` (Tests) | update | neue Dateien `record_ende_test.go` und `server_record_ende_test.go` zu den Randformen aus §6; die Fakes `fakeUpstream`/`fakeSession` (`record_test.go`) und `fakeRecorder` (`server_test.go`) tragen konfigurierbare Antworten. `TestRecordExtendedFehlerantwortVorDemAbbruch` und `TestRecordExtendedSendefehlerNachFehlerantwort` bleiben unverändert: sie prüfen PGR-E6001 aus dem Empfang des Upstream-Ports, kein Ende mit `connection_closed` |
| `docs/user/abdeckung-unit.md` | update (`make abdeckung`) | Tabelle aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` und `slice-v1-abschluss-verbindungsende-replay` liegen in `done/` (Voraussetzung: Formregeln im Modell, Abbildung im Adapter); `slice-v1-abschluss-sessions` und `slice-v1-abschluss-anmeldung` liegen davor in der Reihenfolge, weil sie `record.go` und den Aufbau im Upstream-Adapter berühren. Dieser Slice ist Schritt 40 der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen; die Spezifikation und die Pläne liegen außerhalb der Pfade. **Basis:** `9c19eba6` (Stand des Hauptzweigs nach den Plan-Commits des Anspruchs; die Plan-Commits berühren keinen der Pfade, die Basis ändert sich danach nicht). **Marke: über 1400** geht der Slice zurück zur Zerlegung. **Meldeschwelle: 1000** — erreicht die Summe 1000, meldet der Implementer sofort an den Planner und arbeitet nicht still weiter; dann gilt der vorgemerkte Schnitt des Architect als Rückführungsbedingung: **A** einfache Interaktion mit Ende im Use Case und Zustellung im PGWire-Adapter (DoD-Punkt 1 und 3) | **B** Extended-Interaktion mit Ende (DoD-Punkt 2), B hinter A.
  - *Begründung der Marke (Entscheidung des Orchestrators auf Vorgabe des Nutzers: Konzentration auf das Release, weniger Slices; der Schnitt vor dem Anspruch entfällt):* Schätzung des Architect vor dem Code (2026-10-16, in dieser Metrik, aus den gemessenen Gebern `-replay` 968 und `-replay-extended` 532, Tests zu Produktcode 4,3 bis 9,9): 750 bis 1350, Mitte 1050; die Marke 1400 ist die Obergrenze plus Spielraum. Die Geber mit 968 (`-replay`, ungeteilt in einer Review-Sitzung geprüft) und 732 zeigen, dass diese Größe in einer Sitzung prüfbar ist; jeder Schnitt kostet eine volle Kette (Architect, Implementer, Review, Verifikation, Closure) und fügte bei den Schritten 37, 38 und 39 jeweils mehr Gesamtzeilen hinzu. Die frühere Marke 1000 lag unter der Obergrenze der Schätzung und hätte den Schnitt vor dem Anspruch erzwungen. Liegt die Prüfung des Architect oder des Reviews über der Grenze der Sitzung, wird mindestens einer von `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` zum dritten Mal belegt.
- `in-progress` → `next`: Der Use Case braucht das Ende vom Adapter statt aus den Antworten des Upstream-Ports, oder der Upstream-Adapter muss sich ändern — dritte Schicht, zurück zur Zerlegung.
- `in-progress` → `open`: `slice-v1-abschluss-sessions` oder `slice-v1-abschluss-anmeldung` ändert die Einstufung des Verbindungsendes in `record.go` so, dass dieser Plan nicht mehr trägt — Planner prüft neu.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Setzungen des Architect; die Regel „aufgezeichnet wird genau: der Server beendet in laufender Interaktion mit `FATAL`/`PANIC`“ bestätigte der Nutzer am 2026-10-12). Ort: `LH-FA-25.a` *Aufzeichnen*, *Form*, *Stelle*, *Kein Ende*, `LH-FA-02.b`, `LH-FA-12.a`, `LH-FA-18.a`.

- *Ende bei `PANIC` und `FATAL`:* beide wie `FATAL` — `LH-FA-25.a`.
- *`ERROR`, dann `FATAL` in derselben Interaktion:* die Interaktion trägt beide Fehlerantworten in der Reihenfolge des Upstreams und das Ende; die letzte Fehlerantwort ist die des Endes — `LH-FA-25.a` *Aufzeichnen* (Antworten in Reihenfolge des Upstreams bis einschließlich der Fehlerantwort).
- *`ERROR` und dann Abbruch bleibt `PGR-E6001`:* Das Ende verlangt `FATAL`/`PANIC` — `LH-FA-25.a` *Kein Ende*, `LH-FA-02.b`.
- *Fehlerantwort gelesen oder nicht (Lesestand):* maßgeblich ist, was der Recorder gelesen hat; eine nicht gelesene zählt nicht — `LH-FA-25.a` *Aufzeichnen*, `LH-FA-02.b` *Gelesen*. Die Einstufung im Adapter liegt in `slice-v1-abschluss-verbindungsende-upstream`; hier die Reaktion des Use Case auf das, was der Port liefert.
- *Ende in einer `Sync`-Gruppe, in einer `Flush`-Gruppe, in einer unvollständigen Gruppe:* `Sync` und `Flush` schließen die Gruppe ab, die unvollständige ist `PGR-E6001` — `LH-FA-25.a` *Form*, *Stelle*.
- *Nachrichten des Clients nach der Gruppe mit dem Ende, schon in der Pipeline stehende Interaktionen:* gehören nicht zur Aufzeichnung — `LH-FA-25.a` *Form*. Entschieden (Architect, 2026-10-16): gilt auch für Interaktionen, die schon vollständig an den Upstream gegangen sind; der Use Case nimmt nach der Interaktion mit Ende nichts mehr an, auch nichts Wartendes — `LH-FA-25.a` *Stelle von `connection_closed` beim Aufzeichnen*.
- *Ende mit `--record-timing` und `--record-empty-sessions`:* die Interaktion trägt `offset_ms` wie jede; das Ende ändert die Regeln der leeren Session nicht (die Session hat die Interaktion mit Ende) — `LH-FA-21.a`, `LH-FA-12.a`.
- *Schreibfehler an den Client beim Zustellen der Antworten des Endes:* `EndWriteFailed` wie bei jeder Interaktion; die Interaktion zählt als nicht zugestellt — `LH-FA-02.b`, `LH-FA-13.a`. Entschieden (Architect, 2026-10-16): Die Interaktion entfällt wie jede nicht zugestellte, die Session wird mit den vorherigen übernommen, `PGR-E4003`, die Aufzeichnung trägt kein Ende — `LH-FA-25.a` *Meldung des Endes beim Aufzeichnen*.
- *Schreiben der Session, Zwangsende im Aufbau und Aufbaufehler* (Bindung aus `slice-v1-abschluss-sessions`, Nahtstelle 3 und 4): `CloseSession` entscheidet allein über das Schreiben und die Kennung; das Zwangsende im Aufbau und eine nicht unterstützte Interaktion verwirft die Session weiter, auch mit `--record-empty-sessions`; eine Session mit Ende ist nie leer. Ein Ende des Servers vor dem ersten `ReadyForQuery` erzeugt keine Session (Kennung 0, `CloseSession` mit 0 ist wirkungslos). Prüfung: `TestRecordZwangsendeLeereSession`, `TestRecordLeereSessionNichtUnterstuetztVerworfen`, `TestRecordAufbauFehlerKeineSession` und am Binary `TestE2ERecordKeineVerbindung` bleiben grün, dazu ein Fall `EndForced` nach einer Interaktion mit Ende (geschrieben, nicht „im Aufbau“) — `LH-FA-12.a` *Angenommene Verbindung*, `LH-FA-25.a` *Kein Ende*.
- *Abschluss einer Interaktion mit Ende, Kommentare und Herunterfahren* (Sendung aus `slice-v1-abschluss-verbindungsende-replay-extended-flush`, Planner, 2026-10-16, `AGENTS.md` §3.13; **entschieden vom Architect am 2026-10-16**): (a) Die Kommentare im Record-Pfad folgen dem Wortlaut „Abschluss“ (DoD-Punkt 2). (b) Das Ende ersetzt das `ReadyForQuery`: Die Session endet nach dem Zustellen der Antworten bis einschließlich der Fehlerantwort; vor dem Ablauf der Frist gelesen: keine Meldung, Session geschrieben, nicht gelesen: `PGR-E4006` — `LH-FA-25.a` *Herunterfahren beim Aufzeichnen*, `LH-FA-18.a` *Abbruch*; Code fällt nicht an, es braucht einen Fall mit Mutation im Service (`Shutdown` mitten in der Interaktion, dann das Ende; `EndForced` nach der Interaktion mit Ende). Gestellt war die Frage: Eine laufende Extended-Interaktion, die ein Ende des Servers trägt, hat kein `ReadyForQuery`; der Wortlaut in `LH-FA-18.a` *Abbruch* („erst nach dem `ReadyForQuery` … auch wenn deren `Sync` noch aussteht“, `LH-FA-13.a`) nennt das Ende nicht. Zu entscheiden: endet die Session dann mit dem Ende wie im Replay (`LH-FA-13.a`, `Shutdown` im Replay-Pfad) oder gilt der Wortlaut unverändert; ein Fall und eine Mutation je Entscheidung, falls Code betroffen ist. Die Änderung der Spezifikation macht der Architect.
- *Aufrufe der Client-Richtung einer Session nacheinander:* Bindung aus `slice-v1-abschluss-zeitangaben-aufzeichnen` (`LH-FA-21.a` *Messpunkt*), unverändert — hier genannt.
- *Exit-Code:* ein Ende ist weder Meldung noch Fehler und ändert den Exit-Code nicht — `LH-FA-25.a` *Aufzeichnen*.
- *Formregel der Session* (Bindung aus `slice-v1-abschluss-verbindungsende-format`): `Session.ValidateEnd` ist die eine Regel; der Service baut sie nicht nach. Entschieden (Architect, 2026-10-16): Der Use Case ruft sie nicht; die Reihenfolge erfüllt sie von selbst, weil er nach der Interaktion mit Ende nichts mehr annimmt (zweite Zusicherung wäre ein Zweig, den kein Test erreicht). Test samt Mutation: nach dem Ende wird keine weitere Interaktion übernommen, die Session trägt höchstens ein Ende in der letzten Interaktion — `LH-FA-25.a` *Stelle von `connection_closed` beim Aufzeichnen*, `SPEC-002`.
- *Form beim Schreiben* (Bindung aus `slice-v1-abschluss-verbindungsende-format-schreiber` und `slice-v1-abschluss-verbindungsende-sqlite-schreiber`): Der Schreiber prüft die Form des Endes nicht, in `yaml` wie in `sqlite`; der Use Case sichert sie, sonst ist die Datei für den Leser beschädigt — `SPEC-001` *Schreiber in `yaml`* und *Schreiber in `sqlite`*.
- *Aufteilung der Abbildung im Adapter* (Bindung aus `slice-v1-abschluss-verbindungsende-replay`, `AGENTS.md` §3.13): `send` lässt `connection_closed` als letzte Antwort aus und schließt nicht selbst, der Aufrufer schließt (`richtungen`). Randformen: (a) *`connection_closed` nicht als letzte Antwort:* gibt es nicht, das Format lässt die Stelle nicht zu; `send` lässt nur die letzte aus (`TestReplayEndeNurAlsLetzteAntwort` im Replay-Pfad). (b) *Schreibfehler vor dem Ende:* bleibt `EndWriteFailed`, die Interaktion zählt als nicht zugestellt. (c) *Der Use Case liefert `connection_closed` an `send` oder nicht:* entschieden (Architect, 2026-10-16): er liefert sie, als letzte Antwort; nach `Delivered` beendet der Adapter die Session mit einem eigenen Ereignis `EndServer` (kein Fehler, keine Meldung), `send` lässt sie aus; ein Test mit Mutation im Record-Pfad. Ort: `LH-FA-25.a` *Aufzeichnen*, *Meldung des Endes im Replay* (die Abbildung); die Spezifikation trägt den Replay-Teil, der Record-Teil folgt aus *Aufzeichnen*.

- *Schweregrad, Reihenfolge der Gründe, Stelle (a bis e), Meldung und Zähler:* entschieden in `LH-FA-25.a` (*Stelle von `connection_closed` beim Aufzeichnen*, *Reihenfolge der Gründe beim Aufzeichnen*) und `SPEC-002` (`V` leer wie fehlend: der Use Case ruft nur die Prüfung der Form des Modells, kein eigener Fall). Die Aufzeichnung kennt keinen Zähler des Endes; das Ende erzeugt weder Meldung noch Fehler (`LH-FA-25.a` *Aufzeichnen*).

**Risiken:**

- Der Diff liegt an der Grenze einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×; `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`, 2×): die Rückführung in §4 benennt den Schnitt — **Ausgang:** offen bis Closure.
- `slice-v1-abschluss-sessions` (`CloseSession`, Nebenläufigkeit) und `slice-v1-abschluss-anmeldung` (`OpenSession`, Port) ändern `record.go` und die Ports vor diesem Slice; der Architect liest ihren Stand vor dem Code — **Ausgang:** offen bis Closure.
- Bis `slice-v1-abschluss-verbindungsende-upstream` liefert kein echter Upstream ein Ende; die Tests dieses Slice stützen sich auf Fakes, die nur abbilden, was die Spezifikation sagt — **Ausgang:** offen bis Closure (erwartet: eingetreten, `slice-v1-abschluss-verbindungsende-upstream`, dessen E2E belegt die Entsprechung).

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

Die Belege des Implementers stehen zuerst, die Closure des Planners danach (bei der Closure gefüllt, vor dem `git mv` nach `done/`).

### Belege des Implementers

**Stand.** Code an `a4daf4da`, Tests nachgezogen an `a16c27f0`; danach zwei Kommentare und eine Zusicherung in `TestRecordEndeEnden` (Upstream geschlossen) mit diesem Abschnitt. Geändert sind `services/record.go` (Flag `serverEnde`, `serverBeendet`, `zuordnen`, `bisEnde`, `endetMitEnde`, Zweige in `Query`, `ClientMessage`, `AwaitServer`), `model/recording.go` (`EndServer`), `ports/driving/record.go` (Kommentare), `pgwire/server.go` (`richtungen.schreibe`), neue Tests `services/record_ende_test.go` und `pgwire/server_record_ende_test.go`, die Fakes `fakeUpstream`/`fakeSession` und `fakeRecorder` (konfigurierbare Antworten), `docs/user/abdeckung-unit.md` (`make abdeckung`). Handbuch und README: nicht berührt (§1, `AGENTS.md` §3.11); die Durchsicht nach Quantoren im Handbuch entfällt.

**Umsetzung der Entscheidungen des Architect, soweit sie eine Form im Code haben.** (a) *Hinter `connection_closed`:* `bisEnde` schneidet die Antworten (`Query`) und die Server-Nachrichten (`AwaitServer`) hinter der ersten `connection_closed` ab, bevor `Interaction.Validate` sie sieht (`LH-FA-25.a` *Meldung des Endes beim Aufzeichnen*: „was der Upstream dahinter liefert, übernimmt der Use Case nicht“). Folge: Eine `connection_closed` *vor der letzten Nachricht* erreicht `Validate` über den Use Case nie; die Ablehnung „Mitte“ der *Stelle* trifft eine `connection_closed` hinter einer anderen Nachricht als einer Fehlerantwort `FATAL`/`PANIC` (`TestRecordEndeStelle`, Zeilen „hinter einer anderen Antwort“, „Nachricht dazwischen“). Invariante und Ort: die Folge endet nach `bisEnde` mit der ersten `connection_closed`, `bisEnde` in `record.go`; ohne Kürzung hängt die Ablehnung an der Prüfung `pruefeEnde` des Modells (Mutant M15 rot). (b) *`connection_closed` ohne laufende Interaktion* (Extended, nach einem `ready_for_query`, das die Interaktion schon geschlossen hat): `PGR-E6001` mit der Session als nicht übernehmbar, das ist „daneben ein `ready_for_query`“ der *Stelle* im Extended-Pfad, in dem das `ready_for_query` die Interaktion vor dem Ende schließt (`TestRecordEndeStelle`, Zeile „ReadyForQuery dazwischen“, Lauf `erweitert`; Mutant M9b). Jede andere Server-Nachricht ohne laufende Interaktion bleibt `PGR-E1000` wie bisher; darum läuft die Zeile „ReadyForQuery davor“ nur in der einfachen Anfrage, in der `Validate` es entscheidet. (c) *`Session.ValidateEnd`:* der Use Case ruft sie nicht; `endeAngenommen` ruft sie auf der geschriebenen Session, und das Flag `serverEnde` verhindert jede Interaktion danach (Zeilen „Anfrage“, „wartende Anfrage“, „Pipeline“ in `TestRecordEndeDanachNichts`).

**Läufe.** Unit-Tests von `model`, `services` und `pgwire` im Image `deps` ohne Netz: grün (`go test -count=1`). `make lint` grün nach einer Nacharbeit (`dupl` zwischen einer neuen Verbindungshilfe und `replayVerbindung`; behoben mit der Hilfe `verbindeMitLog` und einem Feld `antworten` in `fakeRecorder`, keine Ausnahme im Profil, `AGENTS.md` §3.2; `errcheck` an einem Aufruf). `make abdeckung` geschrieben. GATES_PLATZHALTER

**Wege der Mutanten.** Je Mutant ein frisches Verzeichnis mit Kopie von `go.mod`, `go.sum`, `internal/` (kein mtime-Problem), ein Textersatz auf der Kopie (`assert`: genau eine Fundstelle), `go test -count=1 -run TestRecordEnde` im Image `deps` (`docker build --target deps`) mit `--network=none`, Bind-Mount der Kopie und eines Go-Caches; Paket `services` für M1 bis M32, `pgwire` für A1 bis A8. Die Ausgabe nennt den roten Test, bei Tabellen die Zeile. Der Lauf aller Mutanten steht an `a16c27f0` (Tests), M32 an dem Stand mit diesem Abschnitt. Mutanten liegen nicht im Repo; E2E-Mutanten entfallen (keine E2E-Fälle in diesem Slice, §1). Kurznamen der Tests ohne Präfix `TestRecordEnde`: `Schweregrad`, `Stelle`, `Dahinter` = `NichtsDahinter`, `Gruppe`, `Danach` = `DanachNichts`, `Enden`, `Einzige` = `EinzigeInteraktion`, `Schreibfehler` (Kern und Adapter getrennt genannt), `Zeit` = `ZeitUndLeereSessions`, `Sendet` = `SendetDavorUndSchliesst` (Adapter). Läufe der Tabellen: `einf`/`erw` = einfache Anfrage / Extended.

| Zusage | Mutation | roter Test |
|---|---|---|
| Schweregrad `FATAL` | `FATAL` aus `beendetVerbindung` gestrichen (M1) | Schweregrad (Zeilen `V FATAL, S FATAL`, `V fehlt, S FATAL`, `V FATAL, S ERROR`, je `einf`/`erw`), Dahinter, Gruppe, Danach, Enden, Schreibfehler, Zeit, Einzige |
| Schweregrad `PANIC` | `PANIC` gestrichen (M2) | Schweregrad (nur `V PANIC, S PANIC`, `V leer, S PANIC`, je `einf`/`erw`) |
| `V` geht vor `S` | nur `S` gelesen (M3) | Schweregrad (`V FATAL, S ERROR` und `V ERROR, S FATAL`) |
| ohne `V` zählt `S` (fehlendes und leeres `V`) | `S` nie gelesen (M4) | Schweregrad (`V fehlt`, `V leer`) |
| leeres `V` wie fehlendes | nur ein fehlender Schlüssel fällt auf `S` zurück (M5) | Schweregrad (nur `V leer, S PANIC`) |
| Ende am Anfang der Antworten nicht aufzeichenbar | `idx == 0` aus der Prüfung entfernt (M6) | Stelle (nur `Anfang`, `einf`/`erw`) |
| … hinter einer anderen Antwort (`command_complete`) | `command_complete` davor zugelassen (M7a) | Stelle (nur diese Zeile), Dahinter (zweiter Teil) |
| … mit einer Nachricht zwischen Fehlerantwort und Ende (`notice_response`) | `notice_response` davor zugelassen (M7b) | Stelle (nur diese Zeile) |
| … nach `ERROR` statt `FATAL`/`PANIC` | jede Fehlerantwort davor zugelassen (M7c) | Stelle (`nach ERROR`), Schweregrad (`V ERROR, S FATAL`) |
| … `ready_for_query` davor (einfach) | `ohneReadyForQuery` prüft nichts (M8) | Stelle (nur `ReadyForQuery davor`) |
| … `ready_for_query` dazwischen (einfach; zwei Prüfungen tragen es) | beide Prüfungen aufgehoben (M9); jede allein bleibt grün (äquivalent, doppelte Absicherung in `pruefeEnde` und `ohneReadyForQuery`) | Stelle (`ReadyForQuery davor`, `dazwischen`/`einf`) |
| … `ready_for_query` dazwischen (Extended) | Zweig „`connection_closed` ohne laufende Interaktion“ entfernt, es bleibt `PGR-E1000` (M9b) | Stelle (nur `dazwischen`/`erw`) |
| Was der Upstream dahinter liefert, übernimmt der Use Case nicht | `bisEnde` kürzt nicht (M15) | Dahinter (`einf`, `erw`) |
| `Sync`-Gruppe mit Ende | `Sync` mit Ende abgelehnt (M10) | Gruppe (`Sync-Gruppe`), alle Läufe `erw` |
| `Flush`-Gruppe mit Ende | `Flush` schließt die Gruppe mit Ende nicht ab (M11) | Gruppe (`Flush-Gruppe`, `Gruppe davor`, `ERROR vor FATAL`) |
| Gruppe davor (`Flush`) bleibt, das Ende liegt in der letzten Gruppe | Server-Nachrichten in die erste Gruppe (M12) | Gruppe (nur `Gruppe davor`) |
| `ERROR` vor `FATAL` in derselben Interaktion | Ende nach einer früheren Fehlerantwort abgelehnt (M13) | Gruppe (nur `ERROR vor FATAL`) |
| nicht abgeschlossene Gruppe ist `PGR-E6001`, Session nicht übernommen | Prüfung „`flush`/`sync` als letzte Client-Nachricht“ der Gruppe entfernt (M14) | Gruppe (nur `nicht abgeschlossene Gruppe`) |
| Fehler des Upstream-Ports vor der Form | Form vor dem Fehler geprüft (M16) | PortfehlerVorForm |
| nach dem Ende keine Anfrage (nicht an den Upstream) | Prüfung am Anfang von `Query` entfernt (M17) | Danach (`Anfrage, herunterfahren true`; ohne Shutdown trägt die Prüfung nach dem Warten, äquivalent, doppelte Absicherung) |
| … eine bei dem Ende wartende Anfrage | Prüfung nach dem Warten entfernt (M18) | Danach (`wartende Anfrage`) |
| … keine Client-Nachricht | Prüfung in `ClientMessage` entfernt (M19) | Danach (`Client-Nachricht`, beide Läufe) |
| … kein Lesen des Upstreams (die Folge dort bleibt ungelesen) | Wartebedingung von `AwaitServer` ohne `serverEnde` (M20) | Danach (`Lesen des Upstreams`, `Pipeline`) |
| … dasselbe, Prüfung nach dem Warten | Prüfung entfernt (M21) | Danach (`Lesen des Upstreams`, `Pipeline`) |
| … keine Interaktion der Pipeline, auch nicht vollständig bis zum `Sync` | `offen` wird beim Ende nicht geleert (M22) | Enden (alle drei Enden mit Pipeline); Danach/`Pipeline` trägt M20/M21 |
| Nach dem Ende läuft keine Interaktion (`laeuft`) | `laeuft` ist nach dem Ende wahr (M23) | Enden (alle), Einzige (`EndForced`) |
| `EndServer`: kein Fehler | `EndServer` liefert `PGR-E4003` (M24) | Schweregrad, Dahinter, Gruppe, Enden (`EndServer`), Zeit, Einzige, Danach/`Pipeline` |
| `EndForced` nach der Interaktion mit Ende: geschrieben, kein `PGR-E4006` | Zwangsende zählt nach dem Ende als laufend (M29) | Enden (`EndForced`), Einzige (`EndForced`) |
| … kein Zwangsende im Aufbau, auch mit nur einer Interaktion | Aufbau ohne die Bedingung „keine abgeschlossene Interaktion“ (M30) | Enden (`EndForced`), Einzige (`EndForced`) |
| … nicht gelesen: `PGR-E4006` | `laeuft` ignoriert `offen` (M26) | Enden (`Zwangsende vor dem Lesen`, alle) |
| Eine Session mit Ende ist nie leer (`LeereSessions` ändert das Ende nicht) | Session mit genau einer Interaktion mit `LeereSessions` verworfen (M28) | Einzige (`leere true`) |
| `Shutdown` mitten in der Interaktion gibt das Ende nicht frei, `Delivered` danach schon | `Delivered` ignoriert das Ende (M31) | Enden (alle, auch `einfache Anfrage`) |
| Upstream wird beim Ende geschlossen | `Close` bei `EndServer` ausgelassen (M32) | Enden (`EndServer`) |
| Schreibfehler vor dem Ende: Interaktion entfällt (einfach) | Ende bleibt „zugestellt“ (M25a) | Schreibfehler (Kern, `einf`) |
| … (Extended) | dto. im Pfad `zuordnen` (M25b) | Schreibfehler (Kern, `erw`) |
| `offset_ms` der Interaktion mit Ende (einfach) | `OffsetMS` bei Ende nil (M27a) | Zeit (`einf`) |
| … (Extended) | dto. (M27b) | Zeit (`erw`) |
| Adapter: `send` lässt `connection_closed` aus | Auslassung entfernt (A1) | Sendet (alle vier), Schreibfehler (Adapter) |
| Adapter: schließt die Verbindung nach dem Ende | kein `beende` bei Ende (A2) | Sendet (alle vier; der Client liest bis zum Ablauf der Frist) |
| Adapter: `Delivered` vor dem Schließen | `Delivered` nach `beende` (A3) | Sendet (alle vier) |
| Adapter: Ereignis `EndServer` | `EndShutdown` gemeldet (A4) | Sendet (alle vier) |
| Adapter: keine Meldung beim Ende | `beende` mit Meldung (A5) | Sendet (alle vier) |
| Adapter: Schreibfehler bleibt `EndWriteFailed` | Fehler bei Ende ignoriert (A6) | Schreibfehler (Adapter) |
| Adapter: `Delivered` meldet das Herunterfahren, das Ende gilt trotzdem | Rückkehr nach `Delivered` mit „endet“ (A8) | Sendet (Läufe `herunterfahren true`, `einf`/`erw`) |
| Adapter: liest nach dem Ende nichts mehr | `schreibe` liefert wahr und die Wache `beendet` am Anfang von `clientRichtung` entfällt (A7b) | Sendet (`einfach`, beide Läufe `herunterfahren`; die zweite Anfrage erreicht den Use Case) |

**Grüne Mutanten.** A7 (`schreibe` liefert nach dem Ende wahr, die Wache bleibt): äquivalent. Invariante: `beende` setzt `r.beendet` vor `CloseSession`, und `clientRichtung` prüft `r.beendet` vor jedem Lesen; Ort: `clientRichtung` in `server.go`. Fällt `return false` allein, liest die Richtung nicht weiter; A7b (beides) ist rot. Die Doppelabsicherungen von M9 und M17 stehen in der Tabelle. Keine weiteren grünen Mutanten.

**Nicht belegt, Grenze.** (1) Der Schreibfehler in der Extended-Richtung des Adapters (`serverRichtung`) hat keinen Test: Ein geschlossener Client lässt beide Richtungen zugleich scheitern (`EndClosed` oder `EndWriteFailed`, je nachdem, welche zuerst ankommt); der Pfad ist derselbe `schreibe` wie der der einfachen Anfrage, den `TestRecordEndeSchreibfehler` (Adapter) belegt. Test-Idee: eine Verbindung über TCP, deren Schreiben scheitert, während die Leserichtung blockiert; Adresse: `slice-v1-abschluss-verbindungsende-upstream` (E2E). (2) Der Fall „Schreibfehler mitten in den Nachrichten vor `connection_closed`“ ist derselbe Fehler aus `Flush` wie der am Ende der Nachrichten; `pgproto3` puffert und schreibt in einem Aufruf. (3) Ob der echte Upstream-Adapter `connection_closed` so liefert, wie die Fakes es abbilden, belegt `slice-v1-abschluss-verbindungsende-upstream` (Risiko aus §6).

**Messung** (Metrik aus §4: `git diff --numstat 9c19eba6..<Commit> -- internal test docs/user`, Einfügungen plus Löschungen). An `a16c27f0`: **820** (795 Einfügungen, 25 Löschungen), unter der Meldeschwelle 1000 und der Marke 1400. Aufschlüsselung: Produktcode (`internal` ohne `*_test.go`) 129, davon `record.go` 103, `server.go` 9, `record.go` des Ports 12, `recording.go` 5; Tests 670 (`record_ende_test.go` 537, `server_record_ende_test.go` 108, Fakes 25 in `record_test.go` und `server_test.go`), `abdeckung-unit.md` 21. Faktor Tests zu Produktcode 5,2 (670 zu 129; ohne Kommentare im Produktcode liegt er höher). Die Schätzung des Architect (750 bis 1350, Mitte 1050) lag über der Messung; die Tabellen (`jeWeg`, `Schweregrad`, `Stelle`) und die Hilfen des Gebers hielten die Tests unter dem Faktor 8 bis 14. Die Nacharbeit nach `a16c27f0` fügt unter 30 Zeilen hinzu.

**Kommentare** (Schritt 20, `AGENTS.md` §3.7, §3.11). Jeder neue oder geänderte Kommentar beschreibt den Ist-Zustand im Indikativ; zu jeder Zusage darin steht ein Test: `serverEnde`, `serverBeendet` (Danach, Enden), `Query`/`AwaitServer`/`ClientMessage` (Dahinter, Danach), `CloseSession` (Enden: nach `EndClosed` und `EndForced` läuft nichts mehr; „Ende des Clients“ allgemein ist auf `EndClosed` und `EndForced` gekürzt, `EndTerminate` ist nicht geprüft), `EndServer` (Enden, Sendet), `schreibe` (Sendet; „den Upstream schließt der Use Case“: Enden, M32), Port-Kommentare (Sendet, Danach). „Sync“ in den vier Kommentaren auf „Abschluss“ gesetzt (`ports/driving/record.go`: `Query`; `record.go`: Feld `offen`, `Query`, `ClientMessage`); `grep -n "Sync"` über beide Dateien findet danach nur Kommentare, die das `Sync` selbst meinen: `ports/driving/record.go:49` (Gruppe mit `Flush` oder `Sync`), `record.go:73` (Gruppe ohne Abschluss, Flush oder Sync), `:81` (Feld `sync` der gesendeten `Sync`-Nachricht), `:207` (Query wartet bei gesendetem `Sync`, das Verhalten ist an das Feld `sync` gebunden), `:279` und `:282` (nach einem `Sync` beginnt die nächste Interaktion; Gruppe mit `Flush` oder `Sync`). Suche des alten Wortes über `internal/` (`grep -rn "Sync" internal --include=*.go`, ohne Tests): weitere Treffer in `replay.go:384` („bis zum Sync verworfen“ meint das Verhalten des Servers nach einem Fehler), `model/extended.go:62` (Gruppe bis `Flush` oder `Sync`), `ports/driving/replay.go:32` (Gruppe schließt mit `Flush` oder `Sync`), `play.go:312` und `:315` (Gruppe, die mit `Sync` endet; die Gruppe mit `Flush` bei Ende steht im Satz danach), `recording/yaml.go:685` (die `Sync`-Gruppe; „bei einem Ende des Servers mit connection_closed“ steht im Satz) und `pgwire/server.go:334` („Sync, bei einem Ende des Servers das abschließende Flush“): unverändert, der Wortlaut meint dort das `Sync` selbst oder nennt den Abschluss mit Ende schon.

**Beobachtungen.** (1) Meine ersten Tests der Zusage „danach nichts mehr“ schlossen die Session mit `EndServer` vor dem Aufruf und prüften damit `l.beendet`, nicht `serverEnde`; erst die Mutanten M17 bis M21 zeigten drei grün (eine Zusage mit zwei Wächtern, einer maskiert den anderen; Zuordnung zu einer Klasse des Registers, etwa `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, beim Planner). Behoben: die Tests lassen die Session offen und fahren die Fälle mit und ohne `Shutdown`. (2) `dupl` im Lint-Profil traf wieder eine Hilfe, die der des Gebers glich (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, Parametrisierung statt Spiegel). (3) Die Spezifikation nennt „Mitte“ und „vor der letzten Nachricht“ der *Stelle* und „was der Upstream dahinter liefert, übernimmt der Use Case nicht“ in derselben Aufzählung; die Umsetzung (a) oben vereint beide, der Architect prüft sie.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×), `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×), `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×) und `BEO-REPO/record-fehlerantwort-im-aufbau-ungeregelt` (1×; die Fehlerantwort im Aufbau bleibt ausgenommen, §1) — der Plan nennt Schätzung und Schnitt vorab (§4). Tritt die Rückführung ein oder misst die Prüfung des Architect über der Grenze, wird mindestens einer der drei Einträge zum dritten Mal belegt und ist eine Lücke. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13), einschließlich der Sendungen aus `slice-v1-abschluss-verbindungsende-format`, `-replay` und `slice-v1-abschluss-sessions` (sie fügen keinen Liefer-Punkt hinzu: sie gehören zu DoD-Punkt 1 und 3): Liefer-Punkte: 3 (einfache Interaktion im Use Case, Extended-Interaktion im Use Case, Zustellung im PGWire-Adapter). Schichten: zwei (Kern mit Service und Ports, `ARC-002`, `ARC-003`, `ARC-004`; PGWire-Adapter, `ARC-006`). Schätzung des Planners: 700 bis 1300 Zeilen mit Tests.

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-sessions` (Nahtstelle 3 und 4, Planner, 2026-10-12, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 1 (Record-Service) und fügt keinen Liefer-Punkt und keine Schicht hinzu: eine Bindung an vorhandenen Code (`CloseSession`, `wirdVerworfen`, `OpenSession`) und ein Testfall (Zwangsende nach einer Interaktion mit Ende). Liefer-Punkte weiter 3, Schichten weiter zwei (Kern, PGWire-Adapter); die Schätzung bleibt bei 700 bis 1300 Zeilen, der Fall kommt mit einigen Dutzend Zeilen hinzu; der Slice liegt weiter an der Grenze aus §4.

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-verbindungsende-format` (Planner, 2026-10-15, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 1 (Record-Service) und fügt keinen Liefer-Punkt und keine Schicht hinzu: ein Aufruf vorhandenen Codes im Modell (`Session.ValidateEnd`, vom Geber geliefert) und ein Testfall, falls der Architect den Aufruf verlangt. Liefer-Punkte weiter 3, Schichten weiter zwei (Kern, PGWire-Adapter); die Schätzung bleibt bei 700 bis 1300 Zeilen, der Slice liegt weiter an der Grenze aus §4.

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-verbindungsende-format-schreiber` (Architect, 2026-10-15, `AGENTS.md` §3.13): Sie liegt in derselben Zusage wie `Session.ValidateEnd` (DoD-Punkt 1, Record-Service) und fügt weder Liefer-Punkt noch Schicht hinzu; Liefer-Punkte weiter 3, Schichten weiter zwei, die Schätzung bleibt bei 700 bis 1300 Zeilen (ein Fall mit Mutation mehr, rund 30 Zeilen); der Slice liegt weiter an der Grenze aus §4, der vorgemerkte Schnitt gilt unverändert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).

Nachgezählt beim Erweitern der Bindung um `sqlite` aus `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Architect, 2026-10-15, `AGENTS.md` §3.13): Dieselbe Zusage (`Session.ValidateEnd` im Use Case, DoD-Punkt 1), kein zusätzlicher Fall: Der Use Case kennt das Format nicht, ein Fall deckt beide. Liefer-Punkte weiter 3, Schichten weiter zwei, die Schätzung bleibt bei 700 bis 1300 Zeilen.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Planner, 2026-10-15, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-sqlite`): Eine Kennung in einer Abgrenzung (Klasse 1), kein Code und kein Test; sie fügt weder Liefer-Punkt noch Schicht hinzu: Liefer-Punkte weiter 3, Schichten weiter zwei.

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-verbindungsende-replay` (Planner, 2026-10-16, `AGENTS.md` §3.13, Closure des Gebers): Sie gehört zu DoD-Punkt 3 (PGWire-Adapter) und fügt weder Liefer-Punkt noch Schicht hinzu: Liefer-Punkte weiter 3, Schichten weiter zwei (Kern, PGWire-Adapter). Hinzu kommt höchstens ein Fall mit Mutation im Record-Pfad (Auslassung), rund 30 bis 60 Zeilen mit Test in der Metrik von §4; die Schätzung von 700 bis 1300 bleibt, der Slice liegt weiter an der Grenze aus §4, der vorgemerkte Schnitt gilt unverändert. Der gemessene Faktor Tests zu Produktcode des Gebers (rund 10,2, dort §7) gilt für die Schätzung des Architect vor dem Code in der Metrik von §4.

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-verbindungsende-replay-extended-flush` (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Verifikation V-2 dort): Sie hat zwei Teile. Die Kommentarkorrektur („Abschluss“ statt „Sync“ an vier Stellen in `services/record.go` und `ports/driving/record.go`) gehört zu DoD-Punkt 2, ist Kommentar ohne Logik und fügt weder Liefer-Punkt noch Schicht hinzu. Die Frage zum Herunterfahren mit einem Ende ist eine Entscheidung des Architect vor dem Code (§6); fällt Code an, ist es ein Fall mit Mutation im Service (Kern, `Shutdown` beim Ende) und bleibt in DoD-Punkt 2. Liefer-Punkte weiter 3 (einfache Interaktion im Use Case, Extended-Interaktion im Use Case, Zustellung im PGWire-Adapter), Schichten weiter zwei (Kern, PGWire-Adapter); `docs/user/` (Abdeckungstabelle) zählt konservativ als dritte, die Entscheidung des Nutzers steht aus. Schätzung: Kommentare 4 bis 8 Zeilen, höchstens ein Fall mit Mutation 30 bis 60 Zeilen in der Metrik von §4; die Schätzung von 700 bis 1300 bleibt, der Slice liegt weiter an der Grenze aus §4, der vorgemerkte Schnitt gilt unverändert; der Nehmer liegt nicht darüber, die Sendung nimmt er an. Datenpunkt des Gebers (Planner): `slice-v1-abschluss-verbindungsende-replay-extended-flush` maß 351 gegen 160 bis 330, die Obergrenze um 21 überschritten, Tests dort 222 für Service-Tests bei Produktcode 0; die Parametrisierung von Spiegeltests wegen `dupl` im Lint-Profil kostet Zeilen, die ein Plan mit Spiegeltests einrechnet.

Nachgezählt bei der Prüfung des Architect vor dem Code (2026-10-16, `AGENTS.md` §3.13, §3.14): Liefer-Punkte 3, Schichten zwei (Kern, PGWire-Adapter); `docs/user/abdeckung-unit.md` ist die aus den Deklarationen der Tests erzeugte Tabelle ohne Prosa — Empfehlung des Architect: keine Schicht; die Entscheidung des Nutzers steht aus, konservativ gezählt sind es drei. Schätzung in der Metrik von §4: Produktcode rund 60 bis 90 Zeilen (Flag und Zweige im Service, ein Ereignis im Modell, ein Zweig in `schreibe`, Kommentare), Tests 650 bis 1250 (Faktor 8 bis 14 wegen zweier Pfade, Service und Adapter), Summe 750 bis 1350, Mitte 1050, Obergrenze über 1000. Empfohlener Schnitt, Linie nach Lieferwert: **A** `-record` — einfache Interaktion mit Ende im Use Case und Zustellung im PGWire-Adapter (DoD-Punkt 1 und 3, `Meldung des Endes beim Aufzeichnen`, `EndServer`, Zwangsende, Schreibfehler, Zeit, leere Session, Aufbaufehler; 2 Liefer-Punkte, Kern und Adapter), Schätzung 450 bis 750, Mitte 600; **B** `-record-extended` — Extended-Interaktion mit Ende (DoD-Punkt 2: Gruppe mit `Sync` oder `Flush`, unvollständige Gruppe, Pipeline, Herunterfahren, Kommentare „Abschluss“; 1 bis 2 Liefer-Punkte, Kern und der Pfad `serverRichtung` des Adapters als Test), Schätzung 400 bis 700, Mitte 550. B steht hinter A, weil es die Übernahme und das Ereignis von A benutzt. Die Sendungen (Kommentare „Abschluss“, Herunterfahren, Gruppe, Pipeline) gehen mit ihren Kennungen an B, die Bindungen aus `-format`, `-format-schreiber`, `-sqlite-schreiber`, `-sessions` und `-replay` an A; der Planner trägt sie im selben Commit ein und zählt beide nach (`AGENTS.md` §3.13).

Entscheidung vor dem Anspruch (Planner, 2026-10-16, Orchestrator auf Vorgabe des Nutzers): Der Slice wird nicht geschnitten; Marke 1400, Meldeschwelle 1000 (§4). Zählung des Architect unverändert geprüft: Liefer-Punkte 3, Schichten zwei (Kern, PGWire-Adapter), `docs/user/abdeckung-unit.md` konservativ als dritte; keine Grenze der Baseline (höchstens 3 Liefer-Punkte, höchstens 2 Schichten) ist ohne `docs/user` verletzt. Konservativ mit `docs/user` sind es drei Schichten, die Entscheidung des Nutzers über die Abdeckungstabelle steht aus; ein Schnitt dafür entfällt, weil die Tabelle ohne Prosa aus den Deklarationen erzeugt wird.
