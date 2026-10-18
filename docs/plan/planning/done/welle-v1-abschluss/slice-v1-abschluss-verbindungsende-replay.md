# Slice slice-v1-abschluss-verbindungsende-replay: Verbindungsende im Replay (einfache Interaktion)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-03.a` · `LH-FA-03.b` · `LH-FA-09.a` · `LH-FA-12.a` · `SPEC-002` · `ARC-002` · `ARC-006`

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

**Ziel:** `replay` gibt ein aufgezeichnetes Ende des Servers wieder: Nach den aufgezeichneten Server-Nachrichten einer einfachen Interaktion (`Query`), deren letzte die `error_response` des Endes ist, schließt der Recorder die Verbindung zum Client, ohne `ReadyForQuery`, ohne Meldung und ohne Fehler; die Session ist damit beendet und die Interaktion verbraucht. Eine aufgezeichnete Lebendprüfung mit Ende ist eine Interaktion wie jede andere (`LH-FA-25.a` *Replay*, `LH-FA-09.a` *Lebendprüfung mit Ende*, `LH-FA-03.a` Schritt 5, `LH-FA-03.b` *Verbraucht*). Nachweis am Binary mit handgeschriebenen Aufzeichnungen, ohne PostgreSQL (Abnahmeszenario 18, Replay-Teil, für die einfache Interaktion).

**Geschnitten vor dem Anspruch** (Planner, 2026-10-16, `AGENTS.md` §3.14; Drift-Log der Roadmap): Die Extended-Interaktion mit Ende (die Gruppe mit dem Ende, in einer `Sync`- wie in einer `Flush`-Gruppe, die Gruppen davor) liefert `slice-v1-abschluss-verbindungsende-replay-extended` (Schritt 37b, neu); dieser Slice bleibt Schritt 37a. Ungeteilt schätzte der Planner 750 bis 1250, Mitte 1000 und damit 100 % der Marke (§8); die Rückführung in §4 war vorab benannt. Dieser Slice liefert die Abbildung im Adapter, die jener voraussetzt.

**Übernommen aus `slice-v1-abschluss-verbindungsende-format`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): die Wiedergabe einer Aufzeichnung mit Ende. Das Laden liefert jener Slice; dieser setzt ihn voraus (Fixtures: handgeschriebene Aufzeichnungen in YAML, `test/integration/testdata`).

**Gebunden durch `slice-v1-abschluss-sessions`** (dort §1, *Nahtstelle zu den Verbindungsende-Slices*, Punkt 1 und 2; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): Dieser Slice hält ein, was `slice-v1-abschluss-sessions` im Replay-Service liefert:

1. *Eine Stelle entscheidet, was für das Replay eine Interaktion ist:* der Schritt `ohneLebendpruefungen` in `NewReplayService`. Eine Lebendprüfung mit Ende bleibt in den Interaktionen (`LH-FA-09.a`); dieser Slice ändert diesen einen Schritt und baut keine zweite Entscheidung daneben. `first-request`, `connection`, `Unassigned` und der Handshake lesen sein Ergebnis.
2. *Bei `connection` ist die Session mit der Annahme zugeordnet:* `OpenConnection` setzt sie (`LH-FA-12.a`), nicht die erste Anfrage. Das Schließen der Verbindung nach der Interaktion mit Ende ändert die Zuordnung nicht, und eine Session, deren letzte Interaktion das Ende ist, zählt wie jede andere für `Unassigned` und für die nicht verbrauchten Interaktionen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Record-Pfad des PGWire-Adapters (Antworten bis zur Fehlerantwort zustellen, beide Verbindungen schließen) und der Use Case `RecordService` — `slice-v1-abschluss-verbindungsende-record` (dort §1 und DoD): eigener Code-Pfad im Adapter (`richtungen`, nicht der Replay-Pfad) und eine andere Schicht des Kerns; mit ihm läge dieser Plan über der Grenze einer Review-Sitzung. Die Abbildung von `connection_closed` auf „keine PGWire-Nachricht“ entsteht hier und wird dort benutzt.
- Das Einlesen einer SQLite-Datei mit Ende als Fixture — `slice-v1-abschluss-verbindungsende-sqlite` (der Leser, liegt davor, Schritt 36a): Der E2E-Test dieses Slice benutzt YAML; die Lesbarkeit der SQLite-Datei belegt jener Slice.
- Einspielen und Vergleich — `-einspielen`, `-vergleich`: `play` benutzt den Replay-Service nicht.
- Eine Ende-Antwort ohne `error_response` im Replay — Out-of-Scope von `LH-FA-25` („Wiedergabe eines Verbindungsendes ohne Fehlerantwort“); das Format lässt sie nicht zu (`PGR-E3003`).
- Die Extended-Interaktion mit Ende (`Parse`, `Bind`, `Describe`, `Execute`, `Close`, `Flush`, `Sync`; Ende in der Gruppe, Gruppen davor) — `slice-v1-abschluss-verbindungsende-replay-extended` (dort §1 und DoD; Klasse 1): Das Wiedergeben von Gruppen trägt einen großen Teil der Tests des Replays (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`); mit ihr läge dieser Plan bei etwa 100 % der Marke. Bis jener Slice schließt, ist die Extended-Interaktion mit Ende im Replay nicht zugesagt und hat hier keinen Test (§6, Zwischenstand).
- Weitere Anmeldeverfahren im `replay` (Klartext-Passwort, MD5, SCRAM) — Entscheidung: nein (Default der Planner-Frage; die Spezifikation trägt sie, *Auswahl im Replay*): `replay` sendet `AuthenticationOk` ohne Anforderung, fordert kein Passwort und bietet kein Verfahren an; ein weiteres Verfahren ist dort ausdrücklich nicht zugesagt. Ein Ende ändert den Aufbau nicht: Es steht in einer Interaktion nach dem Aufbau, und ein Ende vor dem Abschluss des Verbindungsaufbaus wird nicht aufgezeichnet (`LH-FA-25`, *Negative*). Ein Anmeldeverfahren im Replay wäre ein anderer Vorgang als das Ende.
- `NotificationResponse` im Replay — Entscheidung: nein (Default der Planner-Frage): Sie gibt es nur in `record`, denn `replay` liest keinen Server und keine Aufzeichnung enthält sie; `LISTEN` und `NOTIFY` sind aufgezeichnete Anfragen wie jede andere (Spezifikation, *Randformen des Protokollrands*, Modus). Ein Ende ändert daran nichts.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus): Der Replay-Service liefert die aufgezeichneten Antworten einer einfachen Interaktion mit Ende und meldet danach, dass die Verbindung endet; die Interaktion ist mit dem Senden verbraucht; die Session ist beendet. Eine aufgezeichnete Lebendprüfung mit Ende wird nicht übersprungen, zählt für die Zuordnung (`first-request` und `connection`) und für die nicht verbrauchten Interaktionen, und eine eingehende Lebendprüfung mit demselben SQL-Text ist diese Anfrage; die Entscheidung, was eine Interaktion ist, bleibt der eine Schritt `ohneLebendpruefungen`, und bei `--session-assignment connection` bleibt die Session mit der Annahme zugeordnet, auch nachdem die Verbindung mit dem Ende geschlossen hat (Bindung aus `slice-v1-abschluss-sessions`; Test).
- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server): Der PGWire-Adapter bildet `connection_closed` nicht auf eine PGWire-Nachricht ab; er sendet die Nachrichten davor und schließt die Verbindung zum Client, ohne `ReadyForQuery`, ohne Meldung, ohne Fehler und ohne den Exit-Code zu berühren; ein Schreibfehler davor bleibt wie bisher (Test mit Fake-Replayer).
- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus): Am Binary: Ein Client erhält für die Anfrage mit Ende die aufgezeichneten Antworten und danach das Ende der Verbindung; die nächste Verbindung erhält die nächste Session; hat der Client vor der Ende-Interaktion `Terminate` gesendet, ist sie nicht verbraucht (`PGR-W2001`, mit `--fail-on-unconsumed` `PGR-E5002`) (Integrationstest mit handgeschriebener Aufzeichnung).
- [x] Je Zusage des neuen Vertrags eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services/replay.go` | update | Ende der einfachen Interaktion liefern (die Antworten tragen `connection_closed`, `Query` ändert dafür nichts); `ohneLebendpruefungen` behält die Lebendprüfung mit Ende; `Query` beantwortet eine eingehende Lebendprüfung, die diese Anfrage ist (`anfrageMitEnde`), nicht außer der Reihe, bei `first-request` auch vor der Zuordnung; der Pfad der Extended-Interaktion (`ClientMessage`) bleibt Gegenstand von `-replay-extended` |
| `internal/hexagon/ports/driving/replay.go` | update (Kommentar) | Port-Vertrag: Das Ende steht als `connection_closed` letzte Antwort der Lieferung von `Query` wie von `ClientMessage`; keine neue Methode, keine neue Signatur (Architect, 2026-10-16, `LH-FA-25.a` *Meldung des Endes im Replay*) |
| `internal/adapters/driving/pgwire/server.go` | update | `connection_closed` ist keine PGWire-Nachricht (`send` lässt sie als letzte Antwort aus, `toMessage` bleibt unverändert); Replay-Pfad schließt die Verbindung (`replayZustellen`, nach `Sent`) |
| `internal/hexagon/services/replay_ende_test.go`, `internal/adapters/driving/pgwire/server_ende_test.go`, `server_test.go` (Tests) | neu, update | Randformen aus §6; Fake-Replayer im Adapter (`endeReplayer` über dem vorhandenen `fakeReplayer`); `replayVerbindung` nimmt den Port `driving.Replayer` statt des Fakes |
| `test/integration` (`replay_ende_e2e_test.go`, `testdata/ende-fatal.yaml`, `ende-terminate.yaml`, `ende-panic.yaml`, `ende-lebend.yaml`) | neu | Binary mit handgeschriebener Aufzeichnung (Version 2, je Test eine Datei), ohne PostgreSQL; der Client ist roh (ohne Treiber), liest bis zum Ende der Verbindung und sendet nach der Anfrage mit Ende nichts |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` liegt in `done/` (Voraussetzung; das Laden einer Aufzeichnung mit Ende); dieser Slice ist Schritt 37a der Reihenfolge in §5 von welle-v1-abschluss (geschnitten am 2026-10-16, Schritt 37b ist `slice-v1-abschluss-verbindungsende-replay-extended`) und steht nach `slice-v1-abschluss-verbindungsende-sqlite` und `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Schritt 36a und 36b), damit seine Fixtures beide Formate wählen können. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis: `13c23a4559fb0df9e5718bf078c32053cdf10a5b`**, eingetragen vom Planner am 2026-10-16; der Übergang nach `in-progress/` und die Pläne ändern nichts unter den genannten Pfaden); **über 1000** geht der Slice zurück zur Zerlegung. **Meldeschwelle 800:** Erreicht die Summe 800, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung in dieser Metrik nach der Prüfung des Architect (2026-10-16, Faktor Tests zu Produktcode 8 bis 11): 500 bis 850, Mitte 680 und damit 68 % der Marke (§8), das obere Ende 85 %; Marke und Meldeschwelle bleiben. Erreicht die Summe 800, ist der Schnitt Lebendprüfung mit Ende und Zuordnung (`ohneLebendpruefungen`, `first-request`, `connection`) | Rest, Zuschnitt: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`; `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`).
- `in-progress` → `next`: Die Meldung des Endes verlangt eine Änderung des Ports, die eine dritte Schicht trägt (etwa ein Eingriff in `SessionEnd` oder in den Recorder-Port; der Architect hat 2026-10-16 entschieden: keine neue Methode, keine neue Signatur) — zurück zur Zerlegung; der Port `Replayer` liegt in dieser Zählung beim Kern.
- `in-progress` → `open`: Eine Entscheidung des Nutzers ändert die Wiedergabe (etwa „Replay sendet `ReadyForQuery` statt zu schließen“) — Change Request an die Spezifikation.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Setzungen des Architect, geprüft am 2026-10-16; „Lebendprüfung mit Ende wird im Replay als normale Interaktion wiedergegeben“ bestätigte der Nutzer am 2026-10-12). Ort: `LH-FA-25.a` *Replay*, `LH-FA-03.a`, `LH-FA-03.b`, `LH-FA-09.a`, `LH-FA-12.a`.

- *Ende bei `PANIC` und `FATAL`:* beide werden wie aufgezeichnet gesendet, der Schweregrad ändert die Wiedergabe nicht — `LH-FA-25.a` *Replay*.
- *Ende in der Extended-Interaktion (Gruppe mit Ende, `Sync`- wie `Flush`-Gruppe, Gruppen davor):* nicht Gegenstand dieses Slice, entschieden in `LH-FA-25.a` *Replay*, *Form*, `LH-FA-03.a` Schritt 5 und geprüft in `slice-v1-abschluss-verbindungsende-replay-extended`. Zwischenstand: siehe Risiken.
- *Stelle der Nachricht `connection_closed`* (`AGENTS.md` §3.12): (a) *Anfang oder Mitte der Antworten einer Interaktion:* kommt nicht vor, das Format lehnt sie beim Laden ab (`PGR-E3003`, `SPEC-001`, `Session.ValidateEnd`); der Replay-Service prüft sie nicht noch einmal. (b) *Direkt vor dem Ende der Antworten:* die einzige zulässige Stelle; der Service liefert die Antworten bis zur `error_response`, der Adapter sendet sie und schließt, die Nachricht selbst wird nicht gesendet. (c) *Letzte Interaktion der Session:* Ende ist immer die letzte Interaktion ihrer Session (`ValidateEnd`); eine weitere Anfrage des Clients auf dieser Verbindung erreicht den Service nicht, weil der Adapter nach dem Ende nichts mehr liest. (d) *Zwei Enden in einer Aufzeichnung (zwei Sessions mit Ende):* jede Session hat ihr eigenes; die Zuordnung der nächsten Verbindung ist die der Session-Zuordnung (`LH-FA-12.a`). (e) *Lebendprüfung mit Ende:* siehe oben, sie hat dieselbe Stelle. Reihenfolge der Gründe bei einer Anfrage: Abweichung vom erwarteten SQL-Text (`PGR-E5001`) geht dem Ende voraus, ein Ende wird nur bei erwartetem Text gesendet — `LH-FA-25.a`, `LH-FA-03.a` Schritt 1 und 5. **Entschieden** (Architect, 2026-10-16): (a) bis (e) tragen, gebunden in `LH-FA-25.a` *Stelle von `connection_closed` im Replay*; die Reihenfolge der Gründe ist Zuordnung ohne freie Session (`PGR-E5003`), dann Abweichung (`PGR-E5001`), dann das Ende (`LH-FA-25.a` *Reihenfolge der Gründe im Replay*). Die Form der Meldung: `connection_closed` ist letzte Antwort der Lieferung, der Adapter lässt sie beim Senden aus, liest danach nichts mehr und schließt (`LH-FA-25.a` *Meldung des Endes im Replay*).
- *`Terminate` des Clients vor der Ende-Interaktion:* die Interaktion ist nicht verbraucht, `PGR-W2001`; mit `--fail-on-unconsumed` `PGR-E5002` — `LH-FA-03.b` *Zeitpunkte* 1.
- *`--fail-on-unconsumed` bei verbrauchter Interaktion mit Ende:* verbraucht ist sie, sobald ihre Antworten gesendet sind; die Verbindung schließt danach ohne Meldung — `LH-FA-03.b` *Verbraucht*.
- *Lebendprüfung mit Ende:* wird nicht übersprungen, zählt für Zuordnung und nicht verbrauchte Interaktionen, Antworten werden verwendet; eine eingehende Lebendprüfung mit gleichem SQL-Text ist diese Anfrage, mit anderem Text gilt Schritt 1 — `LH-FA-09.a` *Lebendprüfung mit Ende*.
- *Session-Zuordnung:* eine Session mit Ende wird wie jede zugeordnet (`first-request`, `connection`); das Ende ist die letzte Interaktion ihrer Session — `LH-FA-12.a`.
- *Ende mit `--record-timing`:* `offset_ms` ist für das Replay keine Voraussetzung und ändert die Wiedergabe nicht — `LH-FA-21.a`, `SPEC-002`.
- *Entscheidung, was eine Interaktion ist, und Zuordnung bei `connection`* (Bindung aus `slice-v1-abschluss-sessions`, Nahtstelle 1 und 2): `ohneLebendpruefungen` bleibt die eine Stelle; die Lebendprüfung mit Ende wird dort behalten, nicht an anderer Stelle (`first-request`, `connection`, `Unassigned`, Handshake) eigens behandelt. `OpenConnection` setzt die Session bei `connection`; das Schließen nach dem Ende ändert sie nicht. Prüfung: die Tests von `slice-v1-abschluss-sessions` (`TestReplayConnectionNteVerbindungNteSession`, `TestReplayConnectionZuordnungMitDerAnnahme`, `TestReplayConnectionSessionOhneInteraktionNichtZuMelden`) bleiben grün, dazu ein Fall mit Ende bei `connection` — `LH-FA-12.a`, `LH-FA-09.a`.
- *Schreibfehler an den Client beim Senden der Antworten des Endes:* bleibt wie bei jeder Interaktion (Interaktion nicht verbraucht, `PGR-W2001`); kein neuer Fall — `LH-FA-03.b` *Verbraucht*; bestätigt (Architect, 2026-10-16): `LH-FA-25.a` *Meldung des Endes im Replay* sagt es.
- *Eingehende Lebendprüfung mit Ende vor der Zuordnung (`first-request`):* Ist die erste Interaktion der nächsten Session mit Interaktion eine Lebendprüfung mit Ende und gleichem SQL-Text, ist die eingehende Lebendprüfung diese Anfrage und ordnet zu; mit anderem Text gilt Schritt 1 (außer der Reihe, ohne Zuordnung). Bei `connection` ist die Session zugeordnet und der Cursor steht auf ihr — `LH-FA-09.a` *Lebendprüfung mit Ende*, `LH-FA-12.a`. Der Fall braucht einen Test mit Mutation (der eine Schritt `ohneLebendpruefungen` und die Prüfung in `Query` bleiben die einzigen Stellen).
- *Verbindung ohne Session bei `connection` (Nummer über der Zahl der Sessions), eingehende Lebendprüfung mit dem Text der ersten Interaktion einer Session mit Ende:* keine Interaktion am Cursor, also keine Anfrage; sie wird außer der Reihe beantwortet, eine andere Anfrage ist `PGR-E5003` — `LH-FA-12.a` *Mehr Verbindungen als Sessions* ("eine Lebendprüfung wird beantwortet"), `LH-FA-09.a` *Lebendprüfung mit Ende* (die Anfrage vor der Zuordnung gilt nur bei `first-request`). In der Spezifikation entschieden (Commit 94a3a24 und Sessions-Slice), im Review F-1 in §6 nachgetragen; Test `TestReplayConnectionOhneSessionLebendpruefungMitEnde`.
- *Extended-Interaktion mit Ende am Cursor, eingehende leere Lebendprüfung:* keine einfache Anfrage mit Ende, sie wird außer der Reihe beantwortet — `LH-FA-09.a`; `Interaction.Validate` lässt Extended-Interaktionen keine Antworten außerhalb der Gruppen tragen; Test `TestReplayEndeExtendedIstKeineEinfacheAnfrage`.
- *Client hat nach der Anfrage mit Ende schon weitere Nachrichten gesendet:* nicht zugesichert, dass er die Antworten vor dem Ende der Verbindung liest (ungelesene Bytes beim Schließen können die Verbindung zurücksetzen); **akzeptiertes Negativ**, ohne Test: Ein Client, der die Aufzeichnung wiederholt, sendet nach der Anfrage mit Ende nichts (`LH-FA-25.a` *Form*: Nachrichten danach gehören nicht zur Aufzeichnung), und der E2E-Test sendet nichts nach ihr — `LH-FA-25.a` *Meldung des Endes im Replay*.

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×): vor dem Anspruch geschnitten (Extended-Interaktion in `slice-v1-abschluss-verbindungsende-replay-extended`), die Rückführung in §4 benennt den weiteren Schnitt — **Ausgang: entfallen.** Begründung: gemessen 968, unter der Marke 1000 (97 %), die Rückführung trat nicht ein; das Review lief in einer Sitzung ohne Größen-Finding (F-5 ist ein Befund zur Meldung, nicht zur Prüfbarkeit); die Meldeschwelle 800 und die Obergrenze der Schätzung (850) lagen darunter und wurden überschritten, der Datenpunkt steht im Register bei `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Vermerk, kein Zählerzuwachs).
- Zwischenstand: Nach diesem Slice gibt `replay` das Ende einer einfachen Interaktion wieder; eine Extended-Interaktion mit Ende ist nicht zugesagt und nicht getestet, der Adapter bildet `connection_closed` aber für beide Pfade gleich ab — **Ausgang: eingetreten**, Folge-Slice `slice-v1-abschluss-verbindungsende-replay-extended` (nimmt an: dessen §1 *Übernommen aus* und *Gebunden durch* nennen die Extended-Interaktion mit Ende und die Entscheidung des Architect, DoD-Punkt 1 und 2 tragen sie; Closure: `grep` nach der Kennung des Gebers findet sie in §1, §4, §6 und §8 des Nehmers). **Entschieden** (Architect, 2026-10-16): Der Service reicht durch, ohne Code und ohne Meldungscode: `ClientMessage` liefert die Server-Nachrichten der Gruppe wie aufgezeichnet, sie enden mit `connection_closed`, und der Adapter schließt für beide Pfade gleich; eine Ablehnung wäre neuer Code mit neuem Meldungscode für einen Zustand, den `-replay-extended` in einem Slice später mit Tests belegt. Zugesagt und getestet wird sie erst dort.
- Der Schnitt zu `slice-v1-abschluss-verbindungsende-record` lässt die Abbildung von `connection_closed` im Adapter in diesem Slice und die Benutzung im Record-Pfad dort: entschieden (Architect, 2026-10-16): `send` lässt `connection_closed` aus und schließt nicht selbst; das Schließen steht beim Aufrufer (Replay: `replayZustellen`, nach `Sent`; Record: `richtungen`, `slice-v1-abschluss-verbindungsende-record`, dessen §1 die Abbildung bereits als gegeben nennt). Eine spätere Änderung dieser Aufteilung träfe beide — **Ausgang: weiter offen** → Register `BEO-REPO/abbildung-im-adapter-fuer-zwei-pfade-test-nur-in-einem` (neu, 1×); Adresse `slice-v1-abschluss-verbindungsende-record` (Bindung in dessen §1, DoD-Punkt 3, §6 und §8 mit der Closure eingetragen, dort nachgezählt).
- Der E2E-Test stützt sich auf handgeschriebene Aufzeichnungen, bevor `record` ein Ende erzeugen kann: Sie bilden nur ab, was die Spezifikation sagt, nicht was PostgreSQL liefert; die Entsprechung belegt der E2E-Test in `slice-v1-abschluss-verbindungsende-upstream` — **Ausgang: weiter offen** → Register `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (neu, 1×); Adresse `slice-v1-abschluss-verbindungsende-upstream` (dessen §1 *Übernommen aus* `slice-v1-abschluss-verbindungsende-replay` nennt die Entsprechung, sein DoD-Punkt 3 zeichnet ein echtes Ende auf und spielt es ab).

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

Die Belege des Implementers stehen zuerst, die Closure des Planners danach.

### Belege des Implementers

**Wege der Mutanten.** Service und Adapter: Kopie von `go.mod`, `go.sum`, `internal/` außerhalb des Repos, je Mutant ein Textersatz in der Kopie, Testbinary im Image `deps` (`docker build --target deps`) mit `--network=none`, Bind-Mount der Kopie und eines Go-Caches, kein Build-Kontext (kein mtime-Problem), Lauf von `go test -count=1 -run 'Ende|Connection|Lebend|ReplaySent|ReplayModus|FirstRequest'` über `internal/hexagon/services` und `internal/adapters/driving/pgwire`; ohne Mutant beide Pakete `ok`. E2E: Mutant im Arbeitsbaum, `touch` auf die Datei, `make test-integration`, danach `git checkout -- internal`; ohne Mutant grün (Stand 9826736). Mutanten liegen nicht im Repo.

| Zusage | Mutation | roter Test |
|---|---|---|
| Lebendprüfung mit Ende bleibt Interaktion (`ohneLebendpruefungen` die eine Stelle) | `&& !endetEinfach(in)` entfernt | `TestReplayEndeLebendpruefungMitEndeBleibtInteraktion`, `…FirstRequestLebendpruefungVorZuordnung`, `TestReplayConnectionEndeBleibtZugeordnet`; E2E `TestE2EReplayEndeLebendpruefung` |
| eingehende Lebendprüfung mit gleichem Text ist die Anfrage | `&& !s.anfrageMitEnde(c, sql)` entfernt | dieselben drei Unit-Tests |
| … nur bei gleichem Text, anderer Text außer der Reihe | Textvergleich in `anfrageMitEnde` entfernt | `TestReplayEndeLebendpruefungAndererTextAusserDerReihe`, `…FirstRequestLebendpruefungVorZuordnung` |
| first-request: vor der Zuordnung gilt die erste Interaktion der nächsten freien Session | Zweig auf `false` | `…AndererTextAusserDerReihe`, `…FirstRequestLebendpruefungVorZuordnung` |
| … die erste, nicht eine spätere Session | Suche über alle freien Sessions | `…MitEndeBleibtInteraktion`, `…FirstRequestLebendpruefungVorZuordnung` |
| … nur die Interaktion am Cursor (zugeordnete Session) | letzte Interaktion der Session statt Cursor | `…MitEndeBleibtInteraktion` (Lebendprüfung mitten in der Session außer der Reihe) |
| connection: Session mit der Annahme, Cursor auf ihr | zugeordnete Session ignoriert | `…MitEndeBleibtInteraktion`, `TestReplayConnectionEndeBleibtZugeordnet` |
| Schließen nach Ende ändert die Zuordnung nicht | `s.next--` beim Schließen nach Ende | `TestReplayConnectionEndeBleibtZugeordnet` |
| Session mit Ende zählt für Unassigned | Sessions mit Ende aus `Unassigned` | `TestReplayEndeMehrereSessionsMitEnde`, `…FirstRequestLebendpruefungVorZuordnung`, `TestReplayConnectionEndeBleibtZugeordnet` |
| nächste Verbindung erhält die nächste Session; Session mit Ende wird nicht wieder frei | Session nach Ende wieder vorn in `frei` | `TestReplayEndeMehrereSessionsMitEnde` |
| Reihenfolge: PGR-E5001 vor dem Ende | Textvergleich für Interaktion mit Ende übersprungen | `TestReplayEndeAblehnungVorDemEnde` |
| Reihenfolge: PGR-E5003 vor PGR-E5001 und vor dem Ende | ohne freie Session die letzte zuordnen | `TestReplayEndeAblehnungVorDemEnde`, `TestReplayLebendpruefungAusserDerReihe` |
| Antworten unverändert, letzte `connection_closed`, FATAL | Sonderfall im Service, der bei Schweregrad `FATAL` die Nachricht `connection_closed` nicht liefert (der Service verzweigt im Produktcode nicht nach dem Schweregrad; der Eingriff ist ein Nachbau der Verifikation, M23, nicht der Lauf des Implementers) | `TestReplayEndeAntwortenMitConnectionClosed` (Lauf FATAL), E2E `…EndeFatal`, `…EndeLebendpruefung` (Zeitüberschreitung) |
| … PANIC | derselbe Sonderfall bei Schweregrad `PANIC` (Nachbau der Verifikation, M23) | `TestReplayEndeAntwortenMitConnectionClosed` (Lauf PANIC), `…MehrereSessionsMitEnde` |
| ohne `ready_for_query` nach dem Ende (Service) | `ready_for_query` angehängt | `TestReplayEndeAntwortenMitConnectionClosed` u. a. |
| ohne `ready_for_query` nach dem Ende (Adapter) | Adapter sendet `ReadyForQuery` | `TestReplayEndeSendetDavorUndSchliesst`; E2E alle drei Fälle mit Ende |
| Interaktion verbraucht mit dem Senden, Cursor rückt vor | Cursor rückt bei Ende nicht vor | `TestReplayEndeVerbrauchtMitSent` |
| … nicht vor `Sent` | Ende gilt ohne `Sent` als gesendet | `TestReplayEndeVerbrauchtMitSent` (Teil Schreibfehler) |
| … nach `Sent` verbraucht | Ende nach `Sent` nie verbraucht | `TestReplayEndeVerbrauchtMitSent` |
| `Terminate` vor der Ende-Interaktion: unverbraucht (W2001/E5002) | Ende-Interaktion beim Schließen nicht gezählt | `TestReplayEndeTerminateDavorLaesstUnverbraucht`; E2E `TestE2EReplayEndeTerminateDavor` |
| Adapter lässt `connection_closed` als letzte aus | Auslassung entfernt | `TestReplayEndeSendetDavorUndSchliesst`, `…SentVorSchliessen`; E2E `…EndeFatal`, `…EndePanic`, `…EndeLebendpruefung` |
| … nur `connection_closed` | letzte Antwort gleich welchen Typs ausgelassen | `TestReplayEndeSendetDavorUndSchliesst` (u. a., über Zeitüberschreitung) |
| … nur als letzte | überall ausgelassen | `TestReplayEndeNurAlsLetzteAntwort` |
| Adapter schließt nach dem Ende | Rückgabe `folgeWeiter` | `TestReplayEndeSendetDavorUndSchliesst`, `…SentVorSchliessen` (Zeitüberschreitung), E2E (Zeitüberschreitung) |
| Adapter liest nichts mehr | wie vorige | `TestReplayEndeLiestNichtsMehr` (Use Case 2-mal gefragt) |
| `Sent` vor dem Schließen | Schließen vor `Sent` | `TestReplayEndeSentVorSchliessen` |
| Schreibfehler: nicht als gesendet gemeldet | `Sent` vor dem Senden | `TestReplayEndeSchreibfehlerBleibtUnverbraucht`, `TestReplaySent` |
| ohne Meldung und ohne Fehler am Ende | Adapter merkt PGR-E4003 | `TestReplayEndeSendetDavorUndSchliesst` |
| Schließen auch im Pfad der Extended-Nachricht (Fake, nicht zugesagt als Wiedergabe der Extended-Interaktion) | Pfad sendet ohne Schließen | `TestReplayEndeSendetDavorUndSchliesst` (Lauf `sync`) |

Rot „über Zeitüberschreitung“ heißt: der Test liest bis zum Ende der Verbindung und scheitert nach 5 s (Pipe) bzw. 10 s (E2E), weil sie offen bleibt. E2E-Läufe: Mutanten Lebendprüfung übersprungen (`…EndeLebendpruefung`), Terminate-Fall (`…EndeTerminateDavor`, zusätzlich vier Bestandstests), Adapter ohne Auslassung und ohne Schließen und mit `ReadyForQuery` (`…EndeFatal`, `…EndePanic`, `…EndeLebendpruefung`), FATAL ohne Ende (`…EndeFatal`, `…EndeLebendpruefung`; `…EndePanic` bleibt grün, wie es sein muss). Der Mutant „Schließen nur im Query-Pfad“ färbt nur den Sync-Lauf des Adapter-Tests.

**Grüne Mutanten.** Keiner im Endlauf. Zwei frühe Läufe waren grün und wurden behoben: „letzte Interaktion statt Cursor“ (Test hatte die Lebendprüfung nicht mitten in der Session; Fall ergänzt) und „Zuordnung bei connection löschen“ (äquivalent: `alle` wird nur beim Annehmen gelesen; ersetzt durch den Mutanten `s.next--`).

**Messung (§3.14).** `git diff --numstat 13c23a4559fb0df9e5718bf078c32053cdf10a5b..eb5e8fe -- internal test docs/user`: 880 Einfügungen, 17 Löschungen, Summe 897 (Stand 9826736: 898; Marke 1000, Meldeschwelle 800 überschritten und gemeldet, der Planner entschied: weiter, kein Schnitt). Aufschlüsselung auf dem Stand 3c07f28, nachgerechnet von der Verifikation (Einfügungen plus Löschungen): Produktcode 76 (`replay.go` 51, `server.go` 19, Port-Kommentar 6), Unit-Tests 575 (`replay_ende_test.go` 376, `server_ende_test.go` 196, `server_test.go` 3), E2E 202, Testdaten 71, Abdeckungstabellen 44, Summe 968. Faktor Tests (Unit und E2E, 777) zu Produktcode (76) rund 10,2 (Stand eb5e8fe: 897 mit 76, 533, 199, 49 und 40, Faktor 9,6). Die Schätzung des Architect, 500 bis 850, Mitte 680, lag unter der Messung: Datenpunkt zu `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`. Voller `make gates` an eb5e8fe grün. Nach der Nacharbeit zum Review (af7fba9): 951 Einfügungen, 17 Löschungen, Summe 968 (unter der Marke 1000, 97 %; Rückführung tritt nicht ein, die Nacharbeit fügte 71 Zeilen hinzu, davon der größere Teil Tests); voller `make gates` an af7fba9 grün.

**Beobachtungen.** (1) Das Schließen nach dem Ende hängt am Schreiben der Antworten vor ihm; der Adapter-Test über `net.Pipe` sieht das Ende des Stroms als `io.EOF` oder `io.ErrUnexpectedEOF`, beide nimmt der Test an. (2) `send` lässt `connection_closed` auch im Record-Pfad aus, wenn es dort je letzte Antwort wäre; der Record-Pfad liefert sie bis `slice-v1-abschluss-verbindungsende-record` nicht (Adresse `-record`, dessen §1 die Abbildung als gegeben nennt). (3) Ein Client, der nach einer Lebendprüfung mit Ende weiter sendet, ist wie jede Nachricht nach der Anfrage mit Ende akzeptiertes Negativ (§6); kein Test, keine Zusage im Code. (4) Die Zeilen der Schätzung lagen wieder unter der Messung (siehe oben und die Closure unten). (5) `replayVerbindung` nimmt den Port statt des Fakes (3 Zeilen, nötig für `endeReplayer`).

**Nacharbeit nach dem Review (F-1 bis F-7).** F-1 (MEDIUM): Die Spezifikation entscheidet die Randform (`LH-FA-12.a` *Mehr Verbindungen als Sessions*: eine Lebendprüfung wird beantwortet; `LH-FA-09.a`: die Anfrage vor der Zuordnung gilt bei `first-request`); §6 nennt sie jetzt, Test `TestReplayConnectionOhneSessionLebendpruefungMitEnde`; Mutant `case len(s.frei) > 0:` statt `case !s.nachVerbindung && len(s.frei) > 0:` rot, nur dort. F-2: Der Mutant ohne `Request.Type == RequestQuery` war äquivalent (`endetEinfach` liest `Responses`, die eine Extended-Interaktion nicht trägt, `Interaction.Validate`); die Bedingung ist entfernt, der Kommentar nennt den Grund, `TestReplayEndeExtendedIstKeineEinfacheAnfrage` belegt den Fall; Mutant „Gruppe mit Ende gilt als Ende der Anfrage“ rot, nur dort. F-3: `ende-terminate.yaml` für `TestE2EReplayEndeTerminateDavor`, `ende-fatal.yaml` nur für `…EndeFatal`; §3 und Kommentar abgeglichen. F-4: `TestReplayEndeNurAlsLetzteAntwort` prüft, dass keine Nachricht ankommt; Mutant „Rest nach dem Abbildungsfehler senden“ rot in `…NurAlsLetzteAntwort`, `…SchreibfehlerBleibtUnverbraucht`, `TestReplaySent`; der Satz „ein Terminate danach ändert nichts“ in `TestReplayEndeLiestNichtsMehr` ist gestrichen. F-5: Die Messung steht oben; die Entscheidung des Planners „weiter, kein Schnitt“ (Meldeschwelle 800 erreicht, Marke 1000 nicht) gilt als Entscheidung dieses Laufs. F-6: kein Handlungsbedarf (äquivalente Mutanten). F-7: Hinweis an `slice-v1-abschluss-verbindungsende-record`; die Adresse nimmt an: dessen §1 (*Übernommen aus `slice-v1-abschluss-verbindungsende-replay`*) nennt die Abbildung in `send`, seine DoD den Adapter-Test.

### Closure (Planner, 2026-10-16)

**Review und Verifikation.** Review (0 HIGH, 1 MEDIUM F-1, 3 LOW F-2 bis F-4, 3 INFO F-5 bis F-7): F-1 bis F-4 in `af7fba9` nachgearbeitet; F-5 (Meldeschwelle ohne Artefakt) ist hier und im Register behandelt, F-6 äquivalente Mutanten, F-7 Hinweis an `slice-v1-abschluss-verbindungsende-record` (dort eingetragen, siehe Paarungen). Verifikation: PASS, 0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO (V-1 bis V-4); 25 selbst gebaute Mutanten (davon 2 E2E) färbten Tests aus dem erwarteten Grund, die Verifikation prüfte F-1 bis F-4 unabhängig nach. V-1 (Aufschlüsselung der Messung auf dem Stand eb5e8fe) und V-2 (Eingriff der Zeilen FATAL und PANIC) sind in §7 oben nachgezogen; V-3 steht beim Register-Eintrag `meldeschwelle-ohne-artefakt`; V-4 ist ein Vermerk zu `.harness/state/`.

**Gate-Nachweis.** Der volle `make gates` lief an af7fba9 grün; seitdem änderten sich nur Dateien auf der Positivliste (`docs/plan/planning/`, `docs/reviews/`). `make kopf-check` und `make docs-check` laufen nach der Closure.

**Gegenstand:** geliefert (kein Fall von „übernommen“ oder „entfallen“).

**Meldeschwelle (Artefakt).** Die Meldeschwelle 800 war in §4 mit Metrik genannt und wurde bei 897 überschritten (Stand 9826736: 898); der Implementer meldete, und „weiter, kein Schnitt“ fiel in diesem Lauf ohne eigenen Plan-Commit (nur als Satz in §7 des Implementers, Review F-5, Verifikation V-3). Hier festgehalten, nicht nachträglich getroffen: Der Schnitt, den §4 für den Fall nannte (Lebendprüfung mit Ende und Zuordnung), unterblieb; die Marke 1000 wurde nicht erreicht. Die Lücke (Ablage der Entscheidung nach der Meldung) steht als Frage beim Register-Eintrag `BEO-REPO/meldeschwelle-ohne-artefakt`.

**Größe.** 968 in der Metrik von §4 (`git diff --numstat 13c23a4559fb0df9e5718bf078c32053cdf10a5b..3c07f28 -- internal test docs/user`): Produktcode 76, Unit-Tests 575, E2E 202, Testdaten 71, Tabellen 44; Faktor Tests (Unit und E2E) zu Produktcode rund 10,2. Schätzung des Architect 500 bis 850, Mitte 680 (68 % der Marke); gemessen 97 % der Marke, 14 % über der Obergrenze. Unit-Tests trugen (Service 376 in 250 bis 400, Adapter 196 in 100 bis 200); der E2E-Teil mit Testdaten (273 gegen angesetzte 100 bis 180, rund 68 Zeilen je Fall bei vier Fällen) und die Nacharbeit (71 Zeilen) brachten die Summe über die Spanne. Datenpunkt im Register bei `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Vermerk, kein Zählerzuwachs; dieselbe Lesart wie `slice-v1-abschluss-sessions`: Marke nicht erreicht). Folge für `slice-v1-abschluss-verbindungsende-replay-extended`: Gerechnet mit dem Faktor 10,2 liegt dessen Mitte bei rund 610 (87 % der Marke 700, über der Meldeschwelle 550); Vermerk in dessen §8, der Schnitt vor dem Anspruch ist Sache des Planners beim Anspruch.

**Risiken aus §6 — Ausgang je Risiko** (geschlossene Menge):

| Risiko | Ausgang | Folge-Slice bzw. Begründung |
|---|---|---|
| Der Diff liegt über einer Review-Sitzung | entfallen | gemessen 968, 97 % der Marke 1000, Rückführung nicht eingetreten; Review in einer Sitzung ohne Größen-Finding; Meldeschwelle und Obergrenze überschritten, Datenpunkt im Register |
| Zwischenstand: Extended-Interaktion mit Ende nicht zugesagt | eingetreten | `slice-v1-abschluss-verbindungsende-replay-extended` (nimmt an: §1 und DoD-Punkt 1 und 2) |
| Aufteilung `send` lässt aus / Aufrufer schließt trifft beide Pfade | weiter offen | Register `BEO-REPO/abbildung-im-adapter-fuer-zwei-pfade-test-nur-in-einem` (neu); Adresse `slice-v1-abschluss-verbindungsende-record` (Bindung dort eingetragen) |
| E2E stützt sich auf handgeschriebene Aufzeichnungen | weiter offen | Register `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (neu); Adresse `slice-v1-abschluss-verbindungsende-upstream` (§1 nennt die Entsprechung) |

**Register** (Zähler = Dateien unter `evidence/`; Vermerke zählen nicht):

| Eintrag | Zähler | Fund | Stand |
|---|---|---|---|
| `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` | 39 | F-1, F-2: ein Schritt einer Ortsangabe und ein Randzweig ohne Fall und Mutation | verkörpert, Regel bleibt (Retirement-Check: wieder aufgetreten, nicht gelockert) |
| `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | 63 | F-4, V-2: zwei Testkommentare weiter als der Test, zwei Zeilen der Mutationstabelle mit einem Eingriff, der nicht im Produktcode liegt | verkörpert, Regel bleibt |
| `BEO-REPO/plan-folgt-korrektur-nicht` | 51 | F-3, V-1: eine Plan-Zeile weiter als der Diff, Aufschlüsselung auf dem Stand vor der Nacharbeit | verkörpert, Regel bleibt |
| `BEO-REPO/spec-randform-erst-im-review-entschieden` | 32 | F-1: Fall der Zuordnung im Code entschieden, §6 nannte ihn nicht; der Absatz zur Stelle einer Nachricht hat hier gegriffen (kein Befund zu Stellen und Reihenfolge der Gründe) | verkörpert, Regel bleibt |
| `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` | 5 | kein Auftreten: 968 in 500 bis 850 überschritten, Marke nicht erreicht, Faktor 10,2, E2E-Anteil unterschätzt | geplant (`slice-harness-schaetzung-planner-faktor`), Vermerk |
| `BEO-REPO/meldeschwelle-ohne-artefakt` | 1 | Gegenbeleg: Schwelle stand in §4; die Entscheidung nach der Meldung hat keine Ablage (Frage an den Nutzer) | offen, Vermerk |
| `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` | 2 | kein Auftreten: Rückführung trat nicht ein, Teil bei 97 % | offen, Vermerk |
| `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` | 4 | kein Auftreten: Größe vor dem Anspruch vom Planner sichtbar gemacht und geschnitten | geplant, Vermerk |
| `BEO-REPO/abbildung-im-adapter-fuer-zwei-pfade-test-nur-in-einem` | 1 (neu) | Review F-7 und Risiko aus §6 | offen, Adresse `-record` |
| `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` | 1 (neu) | Risiko aus §6 | offen, Adresse `-upstream` |

Kein Eintrag erreicht mit diesem Slice neu 3× ohne Regel: Die vier Einträge mit Zuwachs sind verkörpert, die beiden neuen stehen bei 1×, `BEO-REPO/meldeschwelle-ohne-artefakt` bleibt bei 1×. Zählte der Nutzer die Entscheidung nach der Meldung mit (siehe Frage), wären es 2×.

**Lerneintrag — benannte Lücken und Vorschläge, keine geänderte Regel.** Geändert wird keine Regel, kein `liegt in`. (1) Das Glied einer Konjunktion in einer Bedingung, die eine Zusage trägt (`!s.nachVerbindung && len(s.frei) > 0`), hatte keinen Fall; der Mutant, der es entfernt, blieb grün, bis das Review ihn baute (F-1). `AGENTS.md` §3.10 nennt Pfadformen einer Aufzählung und Schritte einer Ortsangabe, nicht die Glieder einer Bedingung; benannte Lücke: das Glied einer Bedingung als Element. Empfehlung (keine Entscheidung): abwarten bis zum zweiten eigenen Beleg, dann die Glieder in §3.10 nennen. (2) Die Schätzung nach dem Faktor Tests zu Produktcode trug für Unit-Tests und lag für den E2E-Teil mit handgeschriebenen Aufzeichnungen und Roh-Client darunter (273 gegen 100 bis 180); ein Maß je E2E-Fall (rund 68 Zeilen mit Testdaten) fehlt als eigene Zeile. Vorschlag an den Nutzer für `slice-harness-schaetzung-planner-faktor`, dessen Wortlaut er noch bestätigt; diese Closure weist dem Nehmer nichts zu (keine Sendung, `AGENTS.md` §3.13). (3) Die Meldung der Meldeschwelle hat keine Ablage; die Entscheidung des Planners stand nur als Satz im Bericht des Implementers (F-5, V-3). Frage an den Nutzer beim Register-Eintrag. Was funktionierte: Die Stellen von `connection_closed` und die Reihenfolge der Gründe standen vor dem Code in §6 und `LH-FA-25.a` (Prüfung des Architect), und kein Befund betraf sie — der Absatz zur Stelle einer Nachricht in `AGENTS.md` §3.12 hat gegriffen. Der Schnitt vor dem Anspruch hielt den Slice bei 97 % der Marke statt bei geschätzten 100 % und mehr des ungeteilten Plans; die Verifikation baute die Mutanten selbst (25) und fand keine Lücke in den Zusagen, nur im Beleg (V-2).

**Paarungen.** (a) Anker: Der Lerneintrag trägt kein Feld `liegt in`, keine Verkörperung; nichts zu prüfen. (b) Folge-Slice: Genannt sind `slice-v1-abschluss-verbindungsende-replay-extended` (Risiko *eingetreten*, Vermerk zur Schätzung), `slice-v1-abschluss-verbindungsende-record` (Bindung der Aufteilung von `send`, Review F-7) und `slice-v1-abschluss-verbindungsende-upstream` (Entsprechung der Fixtures); `grep -n "slice-v1-abschluss-verbindungsende-replay"` findet die Kennung des Gebers in `-replay-extended` (§1 Zeilen 35 bis 45, §4, §6, §8), in `-record` (§1 Zeilen 41 und 43, DoD-Punkt 3, §6 Zeile 138, §8 Zeile 189) und in `-upstream` (§1 Zeile 39); alle drei existieren in `next/` und sind nicht geschlossen, die Eintragung in `-record` geschah im Commit der Closure mit Nachzählung in dessen §8 (Liefer-Punkte 3, Schichten zwei). (c) Register: Alle genannten Kennungen `BEO-REPO/…` bestehen als Verzeichnis, jedes mit nicht leerem `evidence/` (die zwei neuen mit je einer Datei).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×) und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×) — der Plan nennt seine Schätzung mit Tests und den Schnitt vorab (§4); tritt die Rückführung ein, ist es das dritte Auftreten von `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` und der Eintrag eine Lücke. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13): Liefer-Punkte: 3 (Replay-Service mit Randformen, Abbildung im PGWire-Adapter, E2E am Binary). Schichten: zwei (Kern mit Service und Port, `ARC-002`, `ARC-003`; PGWire-Adapter, `ARC-006`). Schätzung des Planners: 600 bis 1000 Zeilen mit Tests.

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-sessions` (Planner, 2026-10-12, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 1 (Replay-Service) und fügt keinen Liefer-Punkt und keine Schicht hinzu: eine Bindung an zwei vorhandene Schritte (`ohneLebendpruefungen`, `OpenConnection`) und ein Testfall mit Ende bei `connection`. Liefer-Punkte weiter 3, Schichten weiter zwei (Kern, PGWire-Adapter); die Schätzung bleibt bei 600 bis 1000 Zeilen, der Fall kommt mit einigen Dutzend Zeilen hinzu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).

Nachgezählt beim Nachziehen der Kennungen aus dem Schnitt von `slice-v1-abschluss-verbindungsende-sqlite` (Planner, 2026-10-15, `AGENTS.md` §3.13): Eine Ortsangabe in Abgrenzung und Start (Leser 36a, Schreiber 36b liegen davor), kein Code und kein Test; sie fügt weder Liefer-Punkt noch Schicht hinzu: Liefer-Punkte weiter 3, Schichten weiter zwei.

Nachgezählt beim Schnitt vor dem Anspruch (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14): Ungeteilt Produktcode 85 bis 170 (Service, Port, Adapter), Tests mit dem Faktor Tests zu Produktcode von mindestens 5 (gemessene Datenpunkte 8 bis 11; `slice-v1-abschluss-verbindungsende-format`, `-format-schreiber`, `-sqlite`, `-sqlite-schreiber`: 9,4 zuletzt) 550 bis 1000, Testdaten und Abdeckungstabellen 50 bis 100: **750 bis 1250, Mitte 1000, 100 % der Marke 1000**, über der Grenze von 70 %. Geschnitten nach Lieferwert (einfache Interaktion | Extended-Interaktion), nicht nach Schichten. **Dieser Teil (37a):** Produktcode 55 bis 100 (Ende der einfachen Interaktion, Meldung des Endes im Port, Abbildung im Adapter), Tests 300 bis 520 (Service: Ende, `PANIC` und `FATAL`, `Terminate` davor, `--fail-on-unconsumed`, Lebendprüfung mit Ende, `first-request` und `connection`; Adapter: Abbildung, Schreibfehler, Fake-Replayer), Testdaten, E2E und Tabellen 45 bis 80: **400 bis 700, Mitte 550, 55 % der Marke 1000**, oberes Ende 70 %, Meldeschwelle 800. Liefer-Punkte: 3 (Replay-Service mit Randformen, Abbildung im PGWire-Adapter, E2E am Binary). Schichten: zwei (Kern mit Service und Port, `ARC-002`, `ARC-003`; PGWire-Adapter, `ARC-006`); die Abdeckungstabellen in `docs/user` sind von `make abdeckung` erzeugte Ausgabe des Gates und zählen hier wie in `slice-v1-abschluss-protokollrand` (Metrik über `docs/user`, zwei Schichten) nicht als Dokumentations-Schicht, kein Handbuch- oder README-Satz kommt hinzu (Doku: `slice-v1-abschluss-verbindungsende-doku`); Vorbehalt: `slice-v1-abschluss-verbindungsende-sqlite-schreiber` zählte `docs/user` mit, weil sonst nur eine Schicht blieb — der Architect bestätigt die Zählung. Die Schätzung des Planners von 600 bis 1000 (2026-10-12) galt dem ungeteilten Slice und liegt unter der Schätzung mit dem Faktor.

Nachgezählt bei der Prüfung des Architect vor dem Code (2026-10-16, `AGENTS.md` §3.13, §3.14): Liefer-Punkte 3, Schichten zwei (Kern, PGWire-Adapter). Die Abdeckungstabellen zählen nicht: Sie sind die erzeugte Ausgabe von `make abdeckung`, `make abdeckung-check` verlangt sie im Commit der Tests, sie lassen sich also nicht in einen eigenen Slice schneiden; die Zählung von `slice-v1-abschluss-verbindungsende-sqlite` und `-sqlite-schreiber` mit `docs/user` als Schicht wird damit für Abdeckungstabellen berichtigt, die Empfehlung an den Nutzer, `AGENTS.md` §3.13 um diesen Satz zu ergänzen, steht im Bericht des Architect (bis zu seiner Antwort ist es eine Empfehlung). Der Port bekommt keine Methode, der Produktcode sinkt auf etwa 50 bis 90 (Adapter: `send` lässt aus, `replayZustellen` beendet; Service: ein Schritt und eine Bedingung). Tests mit dem Faktor 8 bis 11 des Produktcodes statt 5: Service 250 bis 400, Adapter 100 bis 200, E2E mit Testdaten 100 bis 180, Tabellen 20 bis 40: **Schätzung in der Metrik von §4: 500 bis 850, Mitte 680, 68 % der Marke 1000** (statt 400 bis 700); die Marke 1000 und die Meldeschwelle 800 bleiben, das obere Ende liegt mit 85 % über der Grenze von 70 %: Erreicht die Summe 800, meldet der Implementer, und der Schnitt ist der in §4 benannte (Lebendprüfung mit Ende und Zuordnung).
