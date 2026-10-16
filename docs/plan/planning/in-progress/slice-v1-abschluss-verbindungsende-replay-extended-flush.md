# Slice slice-v1-abschluss-verbindungsende-replay-extended-flush: Verbindungsende im Replay (Extended-Interaktion, Gruppe mit Flush)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-03.a` · `LH-FA-03.b` · `LH-FA-18.a` · `ARC-002` · `ARC-006`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `replay` gibt das aufgezeichnete Ende einer Extended-Interaktion wieder, deren Gruppe mit dem Ende mit `Flush` abschließt: Der Replay-Service liefert die Server-Nachrichten dieser Gruppe bis zur `error_response` des Endes bei der abschließenden `Flush`-Nachricht, danach schließt der Recorder die Verbindung zum Client, ohne `ReadyForQuery`, ohne Meldung und ohne Fehler; die Interaktion ist mit dem Senden verbraucht, die Session beendet; eine Abweichung vor dem `Flush` ist `PGR-E5001` ohne Ende (`LH-FA-25.a` *Replay*, *Form*, `LH-FA-03.a` Schritt 5, `LH-FA-03.b` *Verbraucht*, `LH-FA-18.a`). Nachweis am Binary mit handgeschriebener Aufzeichnung, ohne PostgreSQL (Abnahmeszenario 18, Replay-Teil, für die Extended-Interaktion mit `Flush`).

**Geschnitten aus `slice-v1-abschluss-verbindungsende-replay-extended`** (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Drift-Log der Roadmap): Jener Slice (Schritt 37b) liefert den Pfad `ClientMessage` für die Gruppe mit dem Ende, die mit `Sync` abschließt, samt Gruppen davor, Verbrauch, Abweichung, Zuordnung, Adapter-Test `sync` und E2E; dieser Slice (Schritt 37c) setzt ihn voraus und fügt die Gruppe mit `Flush` hinzu.

**Übernommen aus `slice-v1-abschluss-verbindungsende-replay-extended`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-16): die Gruppe mit dem Ende, die mit `Flush` abschließt.

**Gebunden durch `slice-v1-abschluss-verbindungsende-replay-extended`** (Sendung des Planners, 2026-10-16): Die Randformen aus dessen §6 gelten hier unverändert (Stelle von `connection_closed` in der Gruppe, Client-Nachrichten der Gruppe, `Terminate` und Abbruch, Gruppen davor, Herunterfahren, Fixtures je Fall isoliert); dieser Slice entscheidet keine davon neu. Er benutzt den Pfad, den jener Slice liefert; zeigt sich, dass die `Flush`-Gruppe eine Änderung dieses Pfads, des Ports `Replayer` oder der Abbildung im Adapter (aus `slice-v1-abschluss-verbindungsende-replay`) verlangt, hält der Implementer an und gibt es dem Architect zurück. Die Zuordnung bei `first-request` und `connection` bleibt die aus `slice-v1-abschluss-sessions` (Nahtstelle 1 und 2, über jene Slices); der eine Schritt `ohneLebendpruefungen` bleibt die Entscheidung, was eine Interaktion ist.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Gruppe mit dem Ende, die mit `Sync` abschließt, die Gruppen davor, Verbrauch, Zuordnung, Herunterfahren und der Pfad `ClientMessage` selbst — `slice-v1-abschluss-verbindungsende-replay-extended` (Klasse 1, Voraussetzung dieses Slice): sonst läge der Plan bei etwa 87 % der Marke 700 und darüber.
- Die einfache Interaktion mit Ende, Lebendprüfung mit Ende, Meldung im Port und Abbildung im Adapter — `slice-v1-abschluss-verbindungsende-replay` (Klasse 1, in `done/`): geliefert, dieser Slice ändert sie nicht.
- Der Record-Pfad (Zustellen der Gruppe an den Client, Schließen beider Verbindungen) und der Use Case `RecordService` — `slice-v1-abschluss-verbindungsende-record`: anderer Code-Pfad im Adapter (`richtungen`) und andere Schicht des Kerns.
- Das Einspielen einer Extended-Aufzeichnung mit Ende und der Vergleich — `slice-v1-abschluss-verbindungsende-einspielen`, `-vergleich`: `play` benutzt den Replay-Service nicht.
- Eine Ende-Antwort ohne `error_response` im Replay — Out-of-Scope von `LH-FA-25`; das Format lässt sie nicht zu (`PGR-E3003`).
- Weitere Anmeldeverfahren und `NotificationResponse` im Replay — Entscheidung: nein, wie in `slice-v1-abschluss-verbindungsende-replay` §1 begründet; die Spezifikation trägt sie.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Der Replay-Service liefert bei einer Extended-Interaktion, deren Gruppe mit dem Ende mit `Flush` abschließt, die Server-Nachrichten dieser Gruppe bis zur `error_response` bei `Flush` (ohne `ReadyForQuery`) und meldet danach, dass die Verbindung endet; die Interaktion ist mit dem Senden verbraucht, die Session beendet; eine Abweichung vor dem abschließenden `Flush` ist `PGR-E5001` ohne Ende, die Interaktion bleibt unverbraucht; `Terminate` oder Abbruch vor dieser Gruppe lässt sie unverbraucht (`PGR-W2001`, mit `--fail-on-unconsumed` `PGR-E5002`); die Randformen aus §6 von `slice-v1-abschluss-verbindungsende-replay-extended` gelten (Test je Fall, Fixture je Fall).
- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Am Binary: Ein Client mit Extended Query erhält für die Gruppe mit dem Ende, die mit `Flush` abschließt, die aufgezeichneten Antworten und danach das Ende der Verbindung; die Gruppen davor sind unverändert (Integrationstest mit handgeschriebener Aufzeichnung, eine Datei je Fall; Test des PGWire-Adapters mit Fake-Replayer, Lauf `flush`, zugesagt). Die Zeilen je E2E-Fall mit Testdaten stehen in §7 (Messung).
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
| `internal/hexagon/services/replay.go` | update (Kommentar) | Pfad `ClientMessage` für die Gruppe mit Ende, die mit `Flush` abschließt; der Pfad selbst liefert `slice-v1-abschluss-verbindungsende-replay-extended`; `ClientMessage` liest die Art der abschließenden Nachricht nicht, keine Logik geändert; der Kommentar von `mitten()` sagt „vor ihrem Abschluss“ (F-1 des Reviews) |
| `internal/hexagon/ports/driving/replay.go` | update (Kommentar) | `Shutdown` sagt „deren Abschluss noch aussteht“ (V-1 der Verifikation, `LH-FA-13.a` seit `7dcbd12`); keine Logik, Port unverändert |
| `internal/adapters/driving/pgwire/server.go` | update (Kommentar) | Senden der Gruppe vor dem Schließen im Pfad `replayAntwort`; der Lauf `flush` ist ein Test (`server_ende_test.go`), kein Code; der Kopf von `replaySitzung` sagt „bis zu ihrem Abschluss (Sync, bei einem Ende des Servers das abschließende Flush)“ (F-2 des Reviews, `LH-FA-13.a` seit `7dcbd12`) |
| `internal/hexagon/services/replay_ende_extended_flush_test.go` (neu), `internal/hexagon/services/replay_ende_extended_test.go` (Parameter `abschluss` statt Kopie, `dupl`), `internal/adapters/driving/pgwire/server_ende_test.go` | create, update | Randformen aus §6 für die `Flush`-Gruppe (Service); Lauf `flush` im Adapter, Fake `gruppenReplayer` mit Modus `flush` |
| `test/integration` (`replay_ende_extended_e2e_test.go`, `testdata`, eine Datei je Fall) | update | Fall `TestE2EReplayEndeExtendedFlush` in der Datei, die `slice-v1-abschluss-verbindungsende-replay-extended` anlegt, Fixture `testdata/ende-extended-flush.yaml`; Binary mit handgeschriebener Extended-Aufzeichnung, ohne PostgreSQL |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-replay-extended` liegt in `done/` (Voraussetzung: Pfad `ClientMessage` für die Gruppe mit Ende, Verbrauch, Zuordnung, Datei des E2E-Tests); dieser Slice ist Schritt 37c der Reihenfolge in §5 von welle-v1-abschluss, geschnitten am 2026-10-16. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis:** `e982c7be4073c4e54be8c04a50125389af5ee660`, danach nicht mehr geändert); **über 500** geht der Slice zurück zur Zerlegung. **Meldeschwelle 400:** Erreicht die Summe 400, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik beim Schnitt vor dem Anspruch (Faktor 10, gemessen 10,2): 180 bis 335, Mitte 260 und damit 52 % der Marke (§8), das obere Ende bei 67 %; beim Anspruch ersetzt durch das gemessene Maß des Gebers (532, §7 dort): **160 bis 330, Mitte 245, 49 % der Marke, oberes Ende 66 %** (Nachzählung in §8). Zuschnitt eines weiteren Schnitts: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`).
- `in-progress` → `next`: Die `Flush`-Gruppe verlangt eine Änderung des Pfads aus `slice-v1-abschluss-verbindungsende-replay-extended`, des Ports `Replayer` oder der Abbildung im Adapter — zurück zur Zerlegung (Architect und Planner), nicht still im Code.
- `in-progress` → `open`: Eine Entscheidung des Nutzers ändert die Wiedergabe — Change Request an die Spezifikation.

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

**Randformen** (`AGENTS.md` §3.12), vor dem Code in der Spezifikation entschieden oder vom Architect bei der Prüfung entschieden. Ort: `LH-FA-25.a` *Replay*, `LH-FA-03.a` Schritt 5, `LH-FA-03.b`, `LH-FA-18.a`. Die Randformen aus §6 von `slice-v1-abschluss-verbindungsende-replay-extended` gelten für beide Gruppenarten; hier stehen die, die der `Flush` ändert.

- *Stelle der Nachricht `connection_closed`, Gruppe mit dem Ende, abgeschlossen mit `Flush`* (`AGENTS.md` §3.12, Stelle einer Nachricht): (a) *Letzte Server-Nachricht der Gruppe:* die einzige zulässige Stelle (`PGR-E3003`, `Session.ValidateEnd`). (b) *Letzte Gruppe der Interaktion und der Session:* immer. (c) *`Flush`-Gruppe wie `Sync`-Gruppe:* die Server-Nachrichten vor dem Ende (etwa `ParseComplete`, `BindComplete`) werden gesendet, `ReadyForQuery` gibt es nicht, danach schließt die Verbindung. Entschieden, `LH-FA-25.a` *Replay*, *Form* und *Stelle von `connection_closed` im Replay*: „in einer `Sync`- wie in einer `Flush`-Gruppe die letzte“.
- *Client-Nachrichten der Gruppe:* geliefert wird beim abschließenden `Flush` (`LH-FA-18.a`); eine Abweichung davor ist `PGR-E5001`, es gibt kein Ende, die Interaktion bleibt unverbraucht. Entschieden, `LH-FA-25.a` *Reihenfolge der Gründe im Replay*.
- *Ein `Sync` nach dem `Flush` mit Ende:* gehört nicht zur Aufzeichnung, das Ende ist das letzte Ereignis der Session (`LH-FA-25.a` *Form*); der Adapter liest nach dem Ende nichts mehr.
- *Fixtures je Fall isoliert und Zeilen je E2E-Fall:* wie in §6 von `slice-v1-abschluss-verbindungsende-replay-extended`.
- *Ergänzungen des Architect aus `slice-v1-abschluss-verbindungsende-replay-extended` §6* (Sendung des Architect, 2026-10-16, `AGENTS.md` §3.13; Ort `LH-FA-25.a`): Sie gelten für die `Flush`-Gruppe gespiegelt, je ein Fall und eine Mutation: Abweichung bei der abschließenden Nachricht (`Sync` statt `Flush`), eine einfache Anfrage und eine Lebendprüfung mitten in der Gruppe, `Terminate` oder Abbruch vor dem `Flush`, Ende-Gruppe als einzige Gruppe, Verbrauch beider Seiten. Der Service liest die Art der abschließenden Nachricht nicht anders; der Adapter-Test fährt den Lauf `flush`. Der E2E-Fall ist ein Fall in der Datei des Gebers.
- *Wortlaut von `LH-FA-18.a` gegenüber einer Interaktion ohne `Sync`* (Sendung aus V-5 von `slice-v1-abschluss-verbindungsende-replay-extended`, Planner, 2026-10-16; **entschieden vom Architect vor dem Code**, 2026-10-16): Die Spezifikation definiert den *Abschluss* einer Interaktion (`LH-FA-18.a` *Interaktion*): ihr `Sync`, bei einer Interaktion mit Ende, die mit `Flush` schließt, dieses `Flush`. „Vor deren `Sync`“ heißt an den drei Stellen (`LH-FA-18.a` *Nicht unterstützt* und *Replay*, `LH-FA-25.a` *Reihenfolge der Gründe im Replay*) jetzt „vor ihrem Abschluss“. Folge: Eine einfache Anfrage und eine Lebendprüfung mitten in der `Flush`-Gruppe und zwischen zwei Gruppen derselben Interaktion sind `PGR-E5001`, unverbraucht; zwischen zwei Interaktionen bleibt die Lebendprüfung erlaubt (`LH-FA-09.a`). Der Code (`mitten`: `gruppe > 0 || nachricht > 0`) entscheidet das schon so und bleibt unverändert; Fall und Mutation dieses Slice. [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) sagt „vor deren `Sync`“ und bleibt unverändert; die Spezifikation schärft sie (`AGENTS.md` §3.5).
- *Hinweise des Architect an den Implementer* (2026-10-16): (1) Der Service liest die Art der abschließenden Nachricht nicht; Produktcode 0 ist die Erwartung, die Tests belegen sie. (2) Spiegelt die Fälle des Gebers, wo er Tabellen führt, als Zeilen (Parameter `Abschluss` Sync/Flush), sonst eigene Tests; die Hilfen des Gebers nutzen. (3) Zusagen mit Aufzählung je Element ein Fall und eine Mutation, die nur ihn bricht (`AGENTS.md` §3.10): Verbrauch je Seite (Senden gelingt, Senden scheitert), `Terminate` und Abbruch je ein Fall, Lebendprüfung mitten in der Gruppe und zwischen zwei Gruppen je ein Fall, Modus `connection` und `first-request` je ein Fall. (4) Nicht darstellbar (das Format lehnt es mit `PGR-E3003` ab, `Group.validate`): `Sync` am Ende einer Gruppe vor der Gruppe mit Ende, `ready_for_query` in der Ende-Gruppe, `connection_closed` nicht als letzte Nachricht; dazu keine Zusage und kein Fall. Dargestellt: Gruppen davor mit `Flush`, Ende-Gruppe mit `Flush`. (5) Der Adapter-Test fährt den Lauf `flush` mit einem Fake, der das Ende beim `Flush` der letzten Gruppe liefert und bei dem `Flush` davor nicht (Zähler oder Modus in `gruppenReplayer`), Ablauf wie beim Lauf `sync`; er gehört `TestReplayEndeGruppeDavorSchliesstNicht` als zweiter Lauf. (6) E2E: ein Fall `TestE2EReplayEndeExtendedFlush` in `test/integration/replay_ende_extended_e2e_test.go` mit eigener Fixture `testdata/ende-extended-flush.yaml` (eine Gruppe davor mit `Flush`, Ende-Gruppe mit `Flush`; `parseSync(..., 'H')` der Hilfe des Gebers), nicht mit den anderen Fällen geteilt. (7) Abdeckungs-Deklarationen nur so weit, wie der Test es prüft (`AGENTS.md` §3.11); kein akzeptiertes Negativ ohne benannten Test, sonst als ungeprüft benennen.

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`): die Rückführung in §4 benennt Marke und Meldeschwelle — **Ausgang:** offen bis Closure.
- Der Pfad aus `slice-v1-abschluss-verbindungsende-replay-extended` trägt die `Flush`-Gruppe nicht (er unterscheidet die Gruppenart im Service oder im Adapter): die zweite Rückführung in §4 — **Ausgang:** offen bis Closure.
- Handgeschriebene Extended-Aufzeichnungen bilden nur ab, was die Spezifikation sagt, nicht was PostgreSQL liefert; die Entsprechung belegt der E2E-Test in `slice-v1-abschluss-verbindungsende-upstream` — **Ausgang:** offen bis Closure (erwartet: weiter offen oder entfallen mit jenem Beleg).

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

**Stand.** Code und Tests an `245d245`; keine Logik geändert (an `245d245` Produktcode 0 Zeilen, `git diff --stat e982c7b..245d245 -- internal ':!*_test.go'` leer; danach zwei Kommentare in `replay.go` und `server.go`, Nacharbeit unten). Geändert sind Tests: neue Datei `replay_ende_extended_flush_test.go` (sieben Tests), `replay_ende_extended_test.go` des Gebers (Parameter `abschluss` für Verbrauch, `Terminate`, Herunterfahren, `first-request`, `connection`), `server_ende_test.go` (Fake `gruppenReplayer` mit Modus `flush`, Test mit den Unterfällen `sync` und `flush`), `replay_ende_extended_e2e_test.go` mit Fixture `ende-extended-flush.yaml`, die zwei Abdeckungstabellen.

**Läufe.** `make test-integration` an `3d02a29`: grün, 2 min 27 s, mit `TestE2EReplayEndeExtendedFlush` (der E2E-Code ist danach unverändert). `make abdeckung` und `make abdeckung-check` grün (die Tabellen stehen an `eaec1f3`). Ein erster voller `make gates` an `eaec1f3` war rot an `make lint` (`dupl`: die Kopien von Verbrauch, `first-request` mit `connection` und `Terminate` gleichen den Tests des Gebers; keine Ausnahme im Profil, `AGENTS.md` §3.2), behoben in `245d245` mit dem Parameter `abschluss` (auch für Herunterfahren, derselben Form); danach `make lint` und `make test` einzeln grün. **Voller `make gates` ein Lauf an `245d245`: Exit 0, 11 min 24 s** (10:00:24 bis 10:11:48). Danach nur Plan (Pfade der Positivliste).

**Wege der Mutanten.** Service und Adapter: Kopie von `go.mod`, `go.sum`, `internal/` in ein frisches Verzeichnis je Mutant (kein mtime-Problem), ein Perl-Textersatz auf der Kopie, `go test -count=1` im Image `deps` (`docker build --target deps`) mit `--network=none`, Bind-Mount der Kopie und eines Go-Caches; die Ausgabe nennt den roten Test, bei Tabellen die Zeile (Fehlertext). Der Lauf der 55 Mutanten des Service steht an `245d245`; die sechs des Adapters an `eaec1f3` mit `-run TestReplayEndeGruppeDavorSchliesstNicht/flush` (der Test ist danach unverändert). E2E: Mutant im Arbeitsbaum, `touch`, `make test-integration`, danach `git checkout -- internal`. Mutanten liegen nicht im Repo. Namen ohne Präfix `TestReplayEndeExtendedFlush`: Gruppe = `GruppeMitEnde`, Verbrauch, Abweichung, Terminate = `TerminateUnverbraucht`, Herunterfahren, First = `FirstRequest`, Conn = `Connection`; der Geber heißt „Bestand“.

| Zusage | Mutation | roter Test |
|---|---|---|
| Server-Nachrichten der Gruppe in aufgezeichneter Reihenfolge | letzte Nachricht an die erste Stelle (M1) | Gruppe, Abweichung, First, Conn |
| `connection_closed` als letzte Antwort | letzte Antwort gestrichen (M2) | Gruppe, Abweichung, First, Conn |
| … auch als einzige Gruppe (PANIC) | nur bei genau einer Gruppe gestrichen (M7) | Gruppe, First, Conn |
| ohne `ready_for_query` (verneinend) | `ready_for_query` vor `connection_closed` eingefügt (M4) | Gruppe, Abweichung, First, Conn |
| Gruppe davor mit `Flush` liefert kein `connection_closed` | an die Gruppe davor gehängt (M5b) | Gruppe (nur dieser Test des Slice; Bestand `TestReplayExtendedGruppen`, Sync-`Gruppe`) |
| Gruppe mit `Sync` davor (frühere Interaktion) liefert kein `connection_closed` | an eine Gruppe mit `ready_for_query` gehängt (M5a) | Conn (Session 2); Bestand |
| Lieferung erst mit dem abschließenden `Flush` | Lieferung bei der ersten Nachricht der Gruppe mit Ende (M3) | Gruppe, Abweichung, Verbrauch, Terminate, Herunterfahren, First, Conn |
| … nicht schon bei der mittleren Nachricht | Lieferung bei der zweiten Nachricht (M3b) | dieselben sieben |
| Verbrauch, Seite „Senden scheitert“: unverbraucht ohne `Sent` | `ungesendet` nicht gesetzt (M8) | Verbrauch (Zwangsende vor `Sent`), Herunterfahren |
| Verbrauch, Seite „Senden gelingt“: mit `Sent` verbraucht, keine Meldung | `Sent` löscht nichts (M9) | Verbrauch (Zwangsende nach `Sent`), Terminate, Herunterfahren, First, Conn |
| … ohne `Sent` mit Option `PGR-E5002` | Option ignoriert (M10) | Verbrauch, Terminate |
| … `PGR-E4006` mit der Nummer | `Sequence + 1` (M11) | Verbrauch, Herunterfahren |
| … Zwangsende, Zweig „geliefert, nicht gesendet“ | Zweig entfernt (M20a) | Verbrauch, Herunterfahren |
| … Zwangsende, Zweig „mitten in der Interaktion“ | Zweig entfernt (M20b) | Herunterfahren |
| `Terminate`/Abbruch (für den Service `CloseConnection`): unverbraucht, vor der Gruppe | Cursor vor der Gruppe zählt als verbraucht (T0) | Terminate (Fall „nach 0 Nachrichten“), Verbrauch |
| … mitten in der Gruppe davor | Cursor dort zählt als verbraucht (T1) | Terminate (Fall „nach 1“) |
| … nach der Gruppe davor | dto. (T2) | Terminate (Fall „nach 2“) |
| … mitten in der Gruppe mit Ende, nach der ersten Nachricht | dto. (T3) | Terminate (Fall „nach 3“) |
| … nach der zweiten Nachricht (vor dem `Flush`) | dto. (T4) | Terminate (Fall „nach 4“) |
| … jede Stelle mitten in einer Interaktion | Cursor mitten in der Interaktion zählt als verbraucht (M19) | Terminate; Bestand |
| `Shutdown` erlaubt nur außerhalb der Interaktion (Glied `c.gruppe == 0`) | auf `c.nachricht == 0` verkürzt (M18a) | Herunterfahren; Bestand |
| … (Glied `c.nachricht == 0`) | auf `c.gruppe == 0` verkürzt (M18b) | Herunterfahren; Bestand |
| `Shutdown` je Stand 0 bis 6 (davor, mitten in der Gruppe davor, nach ihr, zwei Stellen in der Gruppe mit Ende, geliefert ohne `Sent`, nach `Sent`) | je Stand ein Mutant, der dort das Gegenteil liefert (S0 bis S6) | Herunterfahren, je Mutant der Fall „Stand i“ |
| Zwangsende je Stand 0 bis 6 | je Stand ein Mutant, der dort das Gegenteil liefert (F0b für Stand 0, F1 bis F6) | Herunterfahren, je Mutant der Fall „Stand i“ |
| Abweichung: erste Nachricht der Interaktion `PGR-E5001` | Vergleich bei Gruppe 1, Nachricht 1 ausgesetzt (M16a) | Abweichung (Zeile „erste Nachricht der Interaktion“) |
| … `Sync` statt `Flush` in der Gruppe davor | Vergleich bei Gruppe 1, Nachricht 2 ausgesetzt (M16e) | Abweichung (nur diese Zeile) |
| … erste Nachricht der Gruppe mit Ende | Gruppe 2, Nachricht 1 (M16b) | Abweichung (Zeile) |
| … mittlere Nachricht der Gruppe mit Ende | Gruppe 2, Nachricht 2 (M16c) | Abweichung (Zeile) |
| … `Sync` statt `Flush` als abschließende | Gruppe 2, Nachricht 3 (M16d) | Abweichung (nur diese Zeile) |
| … einfache Anfrage vor den Gruppen | Art-Prüfung für Gruppe 1, Nachricht 1 und SQL-Vergleich ausgesetzt (M-q0); Glied `istLebendpruefung` entfernt (M17b) | Abweichung (Zeile „vor den Gruppen“) |
| … einfache Anfrage nach der Gruppe davor | dto. für Gruppe 2, Nachricht 1 (M-q2) | Abweichung (Zeile) |
| … einfache Anfrage mitten in der Gruppe mit Ende | dto. für Gruppe 2, Nachricht 2 (M-q3) | Abweichung (Zeile) |
| … Lebendprüfung mitten in der Gruppe davor (Glied `c.nachricht > 0` von `mitten`) | `mitten` auf `c.gruppe > 0` verkürzt (M17f) | Abweichung (Zeile), Herunterfahren; Bestand |
| … Lebendprüfung nach der Gruppe davor, zwischen zwei Gruppen (Glied `c.gruppe > 0`) | `mitten` auf `c.nachricht > 0` verkürzt (M17e) | Abweichung (Zeile), Herunterfahren; Bestand |
| … Lebendprüfung mitten in der Gruppe mit Ende (beide Glieder) | `!c.mitten()` dort aufgehoben (M17g) | Abweichung (nur diese Zeile) |
| Abweichung liefert kein Ende (verneinend) | Abweichung der Art liefert die Server-Nachrichten der letzten Gruppe (M22) | Abweichung; Bestand |
| Cursor bleibt nach der Abweichung stehen | Cursor rückt bei Abweichung der Art vor (M23) | Abweichung (rot durch Laufzeitfehler „Index außerhalb“, wie beim Geber; kein Fehlschlag der Prüfung) |
| `first-request`: Zuordnung ohne freie Session ist `PGR-E5003` vor Abweichung und Ende | ohne freie Session wird die letzte zugeordnet (M12) | First |
| … nächste freie Session | letzte statt erster freier Session (M13) | First |
| `connection`: Verbindung ohne Session ist `PGR-E5003` | Zweig `s.nachVerbindung` ausgeschaltet (M14) | Conn |
| … n-te Verbindung erhält die Session n | Session 1 für jede Verbindung (M14b) | Conn |
| … Schließen ändert die Zuordnung nicht | `s.next--` beim Schließen (M15) | Conn |
| Adapter: Gruppe davor wird gesendet und gemeldet, schließt nicht | jede gelieferte Folge schließt (MA1) | `TestReplayEndeGruppeDavorSchliesstNicht/flush` (auch `/sync`) |
| Adapter: Gruppe mit Ende gesendet, dann `Sent`, dann Schließen (Lauf `flush`) | Schließen vor `Sent` (MA2) / `Sent` fehlt (MA3) | `…/flush` (auch `/sync`) |
| Adapter: ohne `ReadyForQuery` nach dem Ende (verneinend) | `ReadyForQuery` nach dem Ende gesendet (MA4) | `…/flush` (auch `/sync`) |
| Adapter: kein Byte für `connection_closed` | Auslassung entfernt (MA5) | `…/flush` (auch `/sync`) |
| Adapter: `Flush` des Clients erreicht den Service als `Flush` | `pgproto3.Flush` wird als `sync` abgebildet (MA6) | `…/flush` (auch `/sync`, drei Bestandstests) |
| Binary: nach dem Ende schließt die Verbindung | Adapter schließt nicht (E2a) | E2E `TestE2EReplayEndeExtendedFlush` (Zeitüberschreitung 10 s), Sync, Fatal, Panic, Lebendprüfung |
| Binary: die Gruppe davor schließt nicht | `connection_closed` an die Gruppe davor gehängt (E2b) | E2E `TestE2EReplayEndeExtendedFlush`, `…Sync` und fünf Bestandstests |
| Binary: `Flush` des Clients wird als `Flush` verglichen | Abbildung `Flush` auf `sync` (E2c) | E2E `TestE2EReplayEndeExtendedFlush` (und 16 weitere Tests, darunter `…Sync`) |

**Grüne Mutanten.** Keiner der 55 Mutanten des Service und der 6 des Adapters blieb grün. **Nicht gespiegelt, über die Schnittstelle fangbar:** M17c (`endetEinfach` aus `anfrageMitEnde` entfernt) färbt nur Tests des Gebers (`…IstKeineEinfacheAnfrage`, `…LebendpruefungVorDerGruppe`, `TestReplayExtendedAmCursor`), keinen Test dieses Slice; Test-Idee: der Fall „Lebendprüfung, auch die leere, vor der Interaktion mit Ende ist außer der Reihe“ mit `extendedEndeFlush`. Er ist hier nicht geschrieben, weil `anfrageMitEnde` die Antworten einer Extended-Interaktion liest, die keine trägt (`Interaction.Validate`), die Art der abschließenden Nachricht also nie erreicht; der Geber-Test deckt die Zeile. Ebenso nicht gespiegelt: „Interaktion mit `Sync` vor der Interaktion mit Ende“ als eigener Test (Geber `SyncDavor`); M5a färbt hier `Connection` (Session 2 ist eine Interaktion mit `Sync`).

**Akzeptierte Negative, nicht getestet:** (1) `Terminate` und Abbruch mitten in einer Gruppe, die mit `Flush` schließt, am Adapter: kein Fall im Adapter-Test. `TestReplayEndeTerminateMittenInDerGruppe` (Geber) sendet die abschließende Nachricht nie, `replayAntwort` behandelt `Terminate` unabhängig vom Stand der Gruppe; ein Mutant, der nur die `Flush`-Gruppe bricht, existiert nicht. Die Mutanten des Gebers (`Terminate` ignoriert, Lesefehler beendet nicht) stehen in dessen §7. (2) Wie beim Geber: teilweises Senden der Gruppe mit Ende (`replayZustellen`, `TestReplayEndeSchreibfehlerBleibtUnverbraucht`), Abweichung am Binary, ein Fall je Schweregrad (Unit-Fall `PANIC` als einzige Gruppe, E2E `FATAL`).

**Messung (§3.14) an `245d245`, vor der Nacharbeit.** `git diff --numstat e982c7be4073c4e54be8c04a50125389af5ee660..245d245 -- internal test docs/user`: **342** (Einfügungen 306, Löschungen 36; Marke 500 nicht erreicht, Meldeschwelle 400 nicht erreicht; Höchststand 377 an `eaec1f3`, vor der Parametrisierung des Gebers). Aufschlüsselung: Produktcode keine Logik (an `245d245` 0 Zeilen); Service-Tests 222 (neue Datei 162 für sieben Tests, Parameter `abschluss` im Geber 41 Einfügungen und 19 Löschungen; ohne die Parametrisierung hatte die Datei 257); Adapter-Test 50 (36 und 14); E2E 27, Fixture 20, zusammen 47; Abdeckungstabellen 23 (20 und 3). Schätzung 160 bis 330, Mitte 245: gemessen 342 liegt 12 darüber, bei 140 % der Mitte und 68 % der Marke; die Service-Tests (222) lagen über der Schätzung des Architect (bis 175), weil die Parametrisierung Zeilen im Geber kostet und die Abweichungstabelle eine Zeile mehr trägt (einfache Anfrage nach der Gruppe davor). **Zeilen je E2E-Fall:** der Fall mit Einbettung der Fixture 27, Fixture 20, Tabellenzeilen 3, zusammen 50 (Geber: 63); Hilfen (`rohNachricht`, `parseSync`) wiederverwendet. 55 Mutanten in Unit-Läufen im Service, 6 im Adapter, 3 E2E-Läufe.

**Beobachtungen.** (1) `dupl` im Lint-Profil (Schwelle 150) lehnt kopierte Spiegel-Tests ab, sobald ein Test über rund 16 Zeilen dem des Gebers gleicht; der Hinweis (2) des Architect („Zeilen mit Parameter, sonst eigene Tests“) trägt das nicht ohne Ausnahme: drei Funde (vier Tests, dazu Herunterfahren derselben Form) mussten den Parameter `abschluss` bekommen, der Plan zählte sie als eigene Tests. (2) `git checkout -- internal` nach einem E2E-Mutanten setzt auch nicht eingecheckte Änderungen unter `internal` zurück; eine Änderung am Adapter-Test ging so einmal verloren und wurde neu gemacht. Vor einem E2E-Mutanten ist der Arbeitsbaum committet. (3) Der Kommentar von `mitten()` in `replay.go` (Zeile 238) und der Kopf der Sitzung in `server.go` (Zeile 333) sagen „vor ihrem Sync“ bzw. „bis zu ihrem Sync“; mit `Flush` als Abschluss ist `Flush` gemeint (`LH-FA-18.a` *Interaktion*). Keine Logik geändert; zu entscheiden ist, ob ein Wort im Kommentar zu ändern ist (Architect oder Planner). (4) Die Zeile der Abdeckung für `TestReplayEndeExtendedFlushAbweichung` nennt „an jeder Stelle“; die Tabelle der Fälle führt alle fünf Nachrichtenstellen der Interaktion als Zeilen. (5) Keine Randform war offen; kein Halt, keine Rückgabe an den Architect.

### Nacharbeit nach dem Review (`docs/reviews/2026-10-16-review-slice-v1-abschluss-verbindungsende-replay-extended-flush.md`, 0 HIGH, 0 MEDIUM, 2 LOW, 2 INFO)

Nur Kommentare, keine Logik; Architect: `LH-FA-13.a` „bis zu ihrem Abschluss“ mit `7dcbd12`.

| Befund | Ausgang |
|---|---|
| F-1 (LOW): Kommentar von `mitten()` sagt „vor ihrem Sync“ | Kommentar in `replay.go` auf „vor ihrem Abschluss (`LH-FA-18.a`)“ gesetzt |
| F-2 (LOW): Kopf von `replaySitzung` sagt „bis zu ihrem Sync“ | Kommentar in `server.go` auf „bis zu ihrem Abschluss (Sync, bei einem Ende des Servers das abschließende Flush; `LH-FA-13.a`, `LH-FA-18.a`)“ gesetzt; die Testkommentare in `server_extended_test.go` unberührt (Auftrag) |
| F-3 (INFO): M17c färbt nur Tests des Gebers | verlangt nichts; Test-Idee und Grenze stehen oben |
| F-4 (INFO): akzeptiertes Negativ zu `Terminate` am Adapter trägt | verlangt nichts |
| V-1 (LOW, Verifikation): Kommentar von `Shutdown` im Port sagt „deren Sync noch aussteht“ | Kommentar in `ports/driving/replay.go` auf „deren Abschluss noch aussteht“ gesetzt; Suche nach weiteren Kommentaren im Replay-Pfad (`services/replay*.go`, `ports/driving/replay.go`, `server.go`): `replay.go` Z. 384 („bis zum Sync verworfen“) meint das Verwerfen bis zum Sync des Protokolls, nicht den Abschluss der Interaktion, und bleibt; sonst keiner |

**Messung nach der Nacharbeit (§3.14).** `git diff --numstat e982c7be4073c4e54be8c04a50125389af5ee660..<letzter Commit> -- internal test docs/user`: **351 nach V-1** (349 nach F-1/F-2, 342 davor; Kommentarzeilen: `replay.go` 1 und 1, `server.go` 3 und 2, Port 1 und 1), 70 % der Marke 500, Meldeschwelle 400 nicht erreicht. Produktcode: keine Logik, zwei Kommentare. Voller `make gates` ein Lauf an `07faa02`: Exit 0, 13 min 2 s. Nach V-1 voller `make gates` ein Lauf an `6b04534`: Exit 0, 12 min 59 s.

### Closure des Planners (2026-10-16)

**Gate-Nachweis.** Voller `make gates` ein Lauf an `6b04534` grün (Exit 0, 12 min 59 s, §7 oben); die Commits danach (`e36263a` und die Commits dieser Closure) ändern nur Pfade der Positivliste (`docs/plan/planning/`, `docs/reviews/`). `make kopf-check` und `make docs-check` laufen nach der Closure nacheinander.

**Gegenstand:** geliefert (kein Fall von „übernommen“ oder „entfallen“).

**Größe.** 351 in der Metrik von §4 (Messung oben), 70 % der Marke 500, Meldeschwelle 400 nicht erreicht; Schätzung 160 bis 330, Mitte 245: die Obergrenze um 21 (6 %) überschritten, 143 % der Mitte. Datenpunkt als Vermerk bei `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (kein Auftreten: Marke nicht erreicht, ein Zehntel über der Obergrenze, wie bei `slice-v1-abschluss-sessions`).

**Risiken aus §6 — Ausgang je Risiko** (geschlossene Menge):

| Risiko | Ausgang | Folge-Slice bzw. Begründung |
|---|---|---|
| Der Diff liegt über einer Review-Sitzung | entfallen | gemessen 351, 70 % der Marke 500 und unter der Meldeschwelle 400, die Rückführung trat nicht ein; Review in einer Sitzung ohne Größen-Finding (0 HIGH, 0 MEDIUM). Dass die Schätzung um 21 überschritten ist, ist eine Beobachtung, kein offenes Risiko (Register, Vermerk) |
| Der Pfad aus `slice-v1-abschluss-verbindungsende-replay-extended` trägt die `Flush`-Gruppe nicht | entfallen | Produktcode 0 Logik (an `245d245` leer, danach drei Kommentare), Port und Abbildung unverändert; der Service liefert die `Flush`-Gruppe ohne Codeänderung, 55 Mutanten im Service und 6 im Adapter rot, die Verifikation fuhr eigene; die zweite Rückführung in §4 trat nicht ein |
| Handgeschriebene Extended-Aufzeichnungen bilden nur ab, was die Spezifikation sagt | weiter offen | Register `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (drittes Auftreten); `slice-v1-abschluss-verbindungsende-upstream` nimmt die Extended-Fixtures nicht an (§1 nennt die Fixtures von `-replay`, DoD 3 fährt eine einfache Anfrage), deshalb keine Adresse; Frage an den Nutzer im Register und unten |

Ob die Ausgänge tragen: *Entfallen* (Größe) trägt, weil 351 gemessen ist und die Marke eine Metrik nennt (`AGENTS.md` §3.14); nach der Nacharbeit lag die Messung bei 351 statt 342, beide unter der Meldeschwelle. *Entfallen* (Pfad) trägt am Diff, nicht am Bericht: Der Review bestätigte Produktcode 0, die Verifikation fuhr die Mutanten selbst; eine Änderung des Pfads hätte die Mutanten des Gebers und dieses Slice gefärbt. *Weiter offen* nennt bewusst keinen Nehmer: Eine Sendung an `-upstream` hätte dessen DoD-Punkt 3 um einen Extended-Fall gegen PostgreSQL wachsen lassen (`AGENTS.md` §3.13, `BEO-REPO/slice-waechst-durch-uebernahmen`); die Entscheidung liegt beim Nutzer.

**Sendung V-2 (Planner, `AGENTS.md` §3.13).** Die Kommentare im Record-Pfad (`services/record.go` Z. 48, 198, 262, `ports/driving/record.go` Z. 40) sagen „Sync“, der Abschluss einer Interaktion mit Ende kann `Flush` sein. Nehmer ist `slice-v1-abschluss-verbindungsende-record` (Schritt 40, `next/`, nicht geschlossen). Geprüft vor dem Eintrag: Sein §1 *Ausdrücklich NICHT* schließt den Gegenstand nicht aus (der Upstream-Adapter, die SQLite-Aufzeichnung eines echten Endes, die Wiedergabe im Replay und das Handbuch sind ausgeschlossen, nicht Kommentare im Use Case); seine DoD-Punkt 2 führt die Interaktion mit Ende, die mit `Flush` abgeschlossen ist, und ändert diese Pfade ohnehin. Er nimmt an. Eingetragen im Commit `354017d`: §1 (*Übernommen aus*), DoD-Punkt 2, §3 (zwei Zeilen), §6 (Randform für den Architect: Stelle des Herunterfahrens bei einer Interaktion mit Ende, die Spezifikation `LH-FA-18.a` *Abbruch* und `LH-FA-13.a` sagen „erst nach dem `ReadyForQuery` … auch wenn deren `Sync` noch aussteht“, ein Ende hat keins), Kopf (`LH-FA-13.a`) und §8 mit Nachzählung: Liefer-Punkte 3, Schichten zwei (drei konservativ mit `docs/user/`, Nutzerentscheidung steht aus), nicht darüber; `grep -n "replay-extended-flush"` findet die Kennung des Gebers dort. Die Änderung der Spezifikation (Z. 1538: „auch wenn deren `Sync` noch aussteht“) macht der Architect vor dem Code von Schritt 40, nicht der Planner.

**Register** (Zähler = Dateien unter `evidence/`; Vermerke zählen nicht):

| Eintrag | Zähler | Fund | Stand |
|---|---|---|---|
| `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | 65 | F-1, F-2, V-1: drei Kommentare nannten „Sync“ als Abschluss einer Interaktion, die mit `Flush` schließen kann; die Suche der Nacharbeit übersah den Port | verkörpert, Regel bleibt |
| `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` | 3 | Risiko aus §6, Extended-Fixture mit `Flush` | offen, 3× ohne Zuweisung: Frage an den Nutzer |
| `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` | 5 | kein Auftreten: 351 gegen 160 bis 330, Marke 500 nicht erreicht | geplant, Vermerk |
| `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` | 2 | kein Auftreten: Rückführung nicht eingetreten, 70 % der Marke | offen, Vermerk |

*Einordnung der Finding-Klasse „Kommentar nennt Sync, Abschluss kann Flush sein“ (Urteil des Planners):* Sie ist keine eigene Klasse, sie ist `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: Der Kommentar sagt eine Allgemeinaussage über jede Extended-Interaktion zu („bis zu ihrem Sync“), die der Flush-Lauf nicht trägt, und kein Gate prüft Kommentarprosa. Drei Auftreten in einem Slice zählen einmal (ein Vorgang, eine Datei). Der Auslöser (eine Wortlaut-Entscheidung des Architect nach dem Plan zieht die Kommentare nicht nach) steht als Teilmerkmal im Beleg; ein eigener Eintrag entstünde beim nächsten Auftreten allein dieses Merkmals. `BEO-REPO/plan-folgt-korrektur-nicht` trägt es nicht mit: Seine Regel (`AGENTS.md` §3.9) gilt dem Plan und den Harness-Artefakten, nicht Kommentaren im Code, und der Plan folgte dem Diff (Verifikation: kein Finding).

Schwellen 3×: `handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` erreicht sie mit diesem Slice (3 Dateien) ohne Verkörperung und ohne Zuweisung. Der Planner trägt weder `geplant` noch `verkörpert` ein, weil beides eine Entscheidung des Nutzers verlangt (welcher Slice, welcher Zielort); die Welle-Closure liest den Eintrag. `zusage-im-kommentar-weiter-als-pruefung` (65) ist verkörpert. Kein anderer Eintrag erreicht neu 3× ohne Regel.

**Lerneintrag — benannte Lücke und Vorschläge, keine geänderte Regel.** Geändert wird keine Regel, kein `liegt in` (`AGENTS.md` und `.claude/` ändert der Planner nicht). (1) *Wortlaut-Entscheidung ohne Liste der Stellen im Code (benannte Lücke).* Der Architect entschied mit V-5 den Wortlaut „Abschluss“ statt „Sync“ in der Spezifikation; die Stellen im Code, die den alten Wortlaut tragen, fanden Review (F-1, F-2), Verifikation (V-1, Port) und ein zweites Mal die Verifikation (V-2, Record-Pfad): jedes Mal in einem anderen Schritt, weil die Suche der Nacharbeit die Dateien nannte, die der Befund nannte, nicht den Gegenstand. Was funktionierte: Die Verifikation suchte „Sync“ über Produktcode und Tests und entschied je Treffer (richtig, zu ändern, Nehmer); der Slice trug damit die Kommentare in drei Pfaden und gab den vierten an den Nehmer. Vorschlag (keine Entscheidung, der Nutzer kann ihn streichen): Nach einer Wortlaut-Entscheidung der Spezifikation sucht der Implementer das alte Wort über das gesamte `internal/` einschließlich `ports/` und nennt jeden Treffer mit Entscheidung in §7 (`.claude/commands/implement-slice.md` Schritt 17, dort steht die Durchsicht nach Quantoren für das Handbuch); die Sendung an einen Nehmer gehört dazu. (2) *Neuer Sensor: keiner.* Der Zähler von `zusage-im-kommentar-weiter-als-pruefung` verweist auf `slice-harness-mutation`, der Sensor für Kommentarprosa bleibt ohne Zuweisung. (3) *Hinweis des Architect an den Implementer trägt nicht ohne Ausnahme (Beobachtung, kein Eintrag):* „Zeilen mit Parameter, sonst eigene Tests“ ließ `dupl` im Lint-Profil (Schwelle 150) an drei Funden scheitern; die Spiegeltests bekamen den Parameter `abschluss`, was Zeilen im Geber kostete (41 und 19) und die Schätzung um 21 übersteigen ließ. Ein Plan, der einen Slice als Spiegel eines Gebers schneidet, rechnet die Parametrisierung ein; das Maß steht als Vermerk beim Eintrag der Schätzung, für `slice-harness-schaetzung-planner-faktor` (Datenpunkt, keine Sendung: der Nehmer trägt nichts davon). (4) *`git checkout -- internal` nach einem E2E-Mutanten verwirft nicht eingecheckte Änderungen (Beobachtung des Implementers, ein Auftreten, kein Eintrag):* Eine Änderung am Adapter-Test ging so verloren. Vorschlag: Vor einem E2E-Mutanten ist der Arbeitsbaum committet (`.claude/commands/implement-slice.md` Schritt 19, Weg der Mutanten); der Nutzer entscheidet, ob die Zeile aufgenommen wird, bei einem weiteren Auftreten wird ein Register-Eintrag angelegt. (5) Was funktionierte: Produktcode 0 trug zum zweiten Mal; der Schnitt vor dem Anspruch hielt den Slice bei 70 % der Marke; der Review fand nur Kommentare (0 HIGH, 0 MEDIUM).

**Fragen an den Nutzer.** (a) `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` steht bei 3× ohne Verkörperung: Soll `slice-v1-abschluss-verbindungsende-upstream` einen Extended-Fall gegen PostgreSQL aufnehmen (Szenario 18 mit Extended-Client; der Slice nähme dann Extended-Fixtures mit `Sync` und `Flush` auf, sein DoD-Punkt 3 wüchse) oder bekommt die Extended-Hälfte einen eigenen Slice? Die Frage steht seit der Closure von `slice-v1-abschluss-verbindungsende-replay-extended` offen. (b) Zählung von `docs/user/` als dritte Schicht (Nutzerentscheidung zu 37a steht aus): Dieser Slice zählt konservativ mit Markierung, zwei ohne. (c) Die Vorschläge (1) und (4) des Lerneintrags: aufnehmen oder streichen. (d) Die Zählung von Slices mit überschrittener Obergrenze der Schätzung bei nicht erreichter Marke in `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (dieser Slice, `-replay`, `slice-v1-abschluss-sessions`): kein Auftreten, wie bisher; zählte der Nutzer sie mit, wäre die Klasse weit über der Schwelle und bereits `geplant`.

**Paarungen.** (a) Anker: Der Lerneintrag trägt kein Feld `liegt in`, keine Verkörperung; nichts zu prüfen. (b) Folge-Slice: Genannt ist `slice-v1-abschluss-verbindungsende-record` (Sendung V-2); `grep -n "replay-extended-flush" docs/plan/planning/next/slice-v1-abschluss-verbindungsende-record.md` findet die Kennung des Gebers dort (§1, DoD 2, §3, §6, §8); er existiert in `next/`, ist nicht geschlossen und nimmt an. (c) Register: Alle genannten Kennungen `BEO-REPO/…` bestehen als Verzeichnis, jedes mit nicht leerem `evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` — der Plan nennt seine Schätzung mit Tests und den Schnitt vorab (§4); tritt die Rückführung ein, ist es ein weiteres Auftreten von `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`. Keine weiteren Treffer. Der Zähler-Stand ist beim Anspruch aus dem Register zu lesen.

Nachgezählt beim Anlegen (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Schnitt von `slice-v1-abschluss-verbindungsende-replay-extended`): Liefer-Punkte: 2 (Replay-Service für die `Flush`-Gruppe mit Randformen; Binary: Adapter-Test und E2E). Schichten: zwei nach der Auffassung des Architect (Kern mit Service, `ARC-002`; PGWire-Adapter, `ARC-006`, vor allem Tests), Abdeckungstabellen in `docs/user` keine Schicht; **konservativ gezählt mit `docs/user` als dritter Schicht: drei** — markiert, bis der Nutzer entscheidet. Schätzung in der Metrik von §4 (Faktor 10, gemessen 10,2): Produktcode 0 bis 10 (der Pfad besteht, der Architect fand ihn beim Geber durchreichend; nur die Gruppenart ist neu), Tests im Service 50 bis 110, im Adapter 40 bis 70, E2E mit Testdaten 70 bis 110 (rund 68 Zeilen je Fall in `-replay`), Tabellen 10 bis 20: **180 bis 335, Mitte 260, 52 % der Marke 500**, oberes Ende 67 %; Meldeschwelle 400. Die Summe aus beiden Teilen (Mitte 410 und 260) liegt bei 670 gegen 610 ungeteilt; der Mehraufwand ist die zweite Fixture und der zweite Adapter-Lauf.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).

Nachgezählt beim Eintragen der Ergänzungen des Architect (2026-10-16, `AGENTS.md` §3.13, Quelle `slice-v1-abschluss-verbindungsende-replay-extended`): Liefer-Punkte weiter 2, Schichten zwei (drei konservativ mit `docs/user/`, Nutzerentscheidung steht aus). Die Ergänzungen sind Fälle derselben Tests, keine neue Datei und kein Produktcode: Schätzung 170 bis 325 (Produktcode 0 bis 10 statt 10 bis 25), Mitte 245, 49 % der Marke 500, oberes Ende 65 %; Meldeschwelle 400 bleibt.

Nachgezählt beim Eintragen der Sendung V-5 (Planner, 2026-10-16, `AGENTS.md` §3.13, Quelle `slice-v1-abschluss-verbindungsende-replay-extended`): Sie ist eine Frage an den Architect zum Wortlaut der Spezifikation und verlangt keinen Code, keinen neuen Test und keine Datei; der Fall „Lebendprüfung mitten in der Gruppe“ stand in den Ergänzungen schon. Liefer-Punkte weiter 2, Schichten zwei (drei konservativ mit `docs/user/`, Nutzerentscheidung steht aus), Schätzung unverändert.

Vermerk bei der Closure des Gebers (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Schätzung und Plan bleiben unverändert): `slice-v1-abschluss-verbindungsende-replay-extended` maß 532 in der Metrik von §4 (Schätzung des Architect 360 bis 545, Produktcode 0, Service-Tests 322 für neun Tests, Adapter-Tests 87, E2E 61 und Fixture 36, Tabellen 26; Nacharbeit nach dem Review +30). Das Maß je E2E-Fall war 27 Zeilen Fall und 36 Fixture, rund 66 mit der Tabellenzeile; die Hilfen des Gebers (33) nutzt dieser Slice wieder. Gerechnet für diesen Slice, der den Pfad nur spiegelt: Service-Tests als weitere Zeilen derselben Tabellen und Fälle 60 bis 130, Adapter-Test Lauf `flush` 30 bis 60, E2E ein Fall 63 bis 70, Tabellen 6 bis 12, Nacharbeit nach dem Review bis 10 vom Hundert: **rund 160 bis 300, Mitte 230, 46 % der Marke 500**, oberes Ende 60 %; das liegt in der Spanne 170 bis 325 oben und in deren unterer Hälfte. Marke 500 und Meldeschwelle 400 bleiben.

Nachgezählt beim Anspruch (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14; Quelle: Messung 532 und Maß des Gebers in dessen §7): Der Geber schrieb neun Service-Tests mit 322 Zeilen (rund 36 je Test), nicht als Tabellenzeilen. Die `Flush`-Gruppe braucht einen eigenen Aufbau der Interaktion (Gegenstück zu `extendedEnde`) und spiegelt die Fälle der Ergänzungen (Hauptfall, `Sync` statt `Flush`, einfache Anfrage und Lebendprüfung mitten in der Gruppe, `Terminate` und Abbruch, einzige Gruppe, Verbrauch); wo der Geber Tabellen führt (`TestReplayEndeExtendedAbweichung`), sind es Zeilen, sonst eigene Tests. Service-Tests deshalb 60 bis 150 statt 60 bis 130, Adapter-Test Lauf `flush` 30 bis 60, E2E ein Fall mit Fixture 63 bis 70, Tabellen 6 bis 12, Nacharbeit bis 10 vom Hundert: **rund 160 bis 330, Mitte 245, 49 % der Marke 500**, oberes Ende 66 %, unter 70 %: nicht geschnitten. Marke 500 und Meldeschwelle 400 bleiben. Liefer-Punkte 2, Schichten zwei (drei konservativ mit `docs/user/`, Nutzerentscheidung steht aus). §1 (Ausschlüsse) und §6 unverändert; der Punkt „Wortlaut von `LH-FA-18.a`“ bleibt beim Architect.

Nachgezählt bei der Prüfung des Architect vor dem Code (2026-10-16, `AGENTS.md` §3.13, §3.14): Liefer-Punkte 2, Schichten zwei (Kern, PGWire-Adapter; `docs/user/` konservativ als dritte, Nutzerentscheidung steht aus). Die Entscheidung zu V-5 ist ein Satz in der Spezifikation und verlangt keinen Code. Der Service liefert die Gruppe mit Ende und `Flush` schon durch (`ClientMessage` gibt `g.Server` unabhängig von der Art zurück, `replayZustellen` schließt bei `connection_closed`): Produktcode 0, die zweite Rückführung in §4 (Pfad, Port, Abbildung) wird dadurch nicht berührt. Schätzung 160 bis 330, Mitte 245, Marke 500, Meldeschwelle 400 bestätigt; die Service-Tests liegen bei getrennten Tests am oberen Ende (bis 175), die Summe bleibt unter 400.

Nachgezählt bei der Übergabe (Implementer, 2026-10-16, `AGENTS.md` §3.13, §3.14): Liefer-Punkte 2, Schichten zwei (Kern, PGWire-Adapter; `docs/user/` konservativ als dritte, Nutzerentscheidung steht aus). Gemessen 342 in der Metrik von §4 (Messung in §7), 68 % der Marke 500, Meldeschwelle 400 nicht erreicht; Schätzung 160 bis 330, Mitte 245: 12 darüber. keine Logik geändert, die zweite Rückführung in §4 trat nicht ein; die Marke ebenso nicht. §1 (Ausschlüsse) und §6 unverändert; die Tests des Gebers sind um den Parameter `abschluss` geändert (§3), ohne Änderung ihrer Aussagen.
