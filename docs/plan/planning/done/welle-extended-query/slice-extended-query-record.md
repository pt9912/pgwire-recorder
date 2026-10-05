# Slice slice-extended-query-record: Extended Query im Record

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-extended-query.

**Bezug:** [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-06`](../../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0012](../../adr/0012-extended-query-gruppen.md), [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md)

**Berührte Spec-Stellen:** `LH-FA-18.a` · `LH-FA-13.a` · `SPEC-041` · `SPEC-002` · `ARC-002` · `ARC-003` · `ARC-004` · `ARC-006` · `ARC-007`

**Verantwortlich:** —
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Recorder vermittelt `Parse`, `Bind`, `Describe`, `Execute`, `Close`, `Flush` und `Sync` zwischen Client und Upstream und zeichnet sie geordnet auf.

Geliefert nach [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md): Der Record-Pfad vermittelt je Session in zwei unabhängigen Richtungen, und der Record-Service führt den Interaktionszustand allein.

- Upstream-Port: `Query` (einfach, synchron), `Send` (Client-Nachrichten einer Gruppe), `Receive` (nächste Server-Nachrichten), `Close`. Am Port steht der Vertrag: `Send` und `Receive` laufen gleichzeitig und blockieren nur an ihrer Richtung; `Close` blockiert nicht und beendet jedes wartende `Send`, `Receive` und `Query` mit einem Fehler.
- Record-Use-Case: `Query`, `ClientMessage`, `AwaitServer`, `Delivered`, `Shutdown`, `CloseSession`. Beide Richtungen rufen gleichzeitig; `CloseSession` beendet jeden wartenden Aufruf mit `ErrSessionEnded`. Der Service gruppiert nach `LH-FA-18.a`, ordnet Server-Nachrichten der ältesten laufenden Interaktion zu (Pipelining über mehrere `Sync`), übernimmt jede Interaktion erst nach `Interaction.Validate`, entscheidet die Einstufung jedes gemeldeten Endes (`PGR-E4003`, verworfene nicht zugestellte Interaktion) und ob das Herunterfahren wartet. Eine einfache Anfrage nach einem `Sync` wartet, bis dessen Antworten zugestellt sind; sie läuft nie gleichzeitig mit dem Empfang.
- Herunterfahren: Die Client-Richtung meldet es vor dem Lesen der nächsten Nachricht, eine schon gelesene wird also verarbeitet. Danach beginnt keine neue Interaktion (`ErrShutdown` aus `Query` und `ClientMessage`); die laufende nimmt bis zu ihrem `Sync` noch ihre Nachrichten an, und die Session endet nach deren `ReadyForQuery`. Wecken und Zurücksetzen der Lesefrist sind unter einer Sperre geordnet.
- PGWire-Adapter: im Record `recordSitzung` mit `clientRichtung` und `serverRichtung` als Transport ohne Interaktionszustand; er meldet Client-Nachricht, Verbindungsende, `Terminate`, Schreibfehler zum Client und Herunterfahren. Endet die Session, schließt er die Client-Verbindung; eine Fehlerantwort dazu schreibt er höchstens eine Sekunde lang. Im Replay `replaySitzung`, synchron.
- Gate: Die Stufe `test` des `Dockerfile` prüft die Formatierung (`gofmt -l`).
- Spezifikation: Randformen in `LH-FA-18.a` (späte Server-Nachrichten einer Flush-Gruppe, Herunterfahren mit laufender Interaktion, `Query` in laufender Interaktion, Zielart außer `S`/`P`, Form der Interaktion, beide Richtungen unabhängig, Ende der Client-Verbindung beim Senden oder Warten) und in `LH-FA-13.a` (schon gelesene Client-Nachricht beim Herunterfahren, danach keine neue Interaktion); Ende einer Session schließt die Client-Verbindung; Architektur-Sicht §2.3, §4, §4.5, §5.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Replay — `slice-extended-query-replay`. Im Replay-Modus lehnt `replaySitzung` Extended-Nachrichten mit `PGR-E6001` ab.
- `CancelRequest` ohne Schlüssel (Review F-310) — `slice-v1-abschluss-cancel-ohne-schluessel`; betrifft den Protokollrand, nicht Extended Query.
- Domain-Typen, Formregeln (`Interaction.Validate`) sowie Schreiber und Leser des Formats `yaml` für `type: extended` — geliefert von `slice-extended-query-modell`; dieser Slice bildet die PGWire-Nachrichten auf diese Typen ab.
- Asynchrone Server-Nachrichten zwischen zwei Interaktionen (etwa `NoticeResponse` ohne Anfrage) — unverändert wie bei einfachen Anfragen: der Recorder liest sie mit der nächsten Interaktion; keine Anforderung regelt sie anders.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Ein Go-Client mit Prepared Statements läuft über den Recorder, und die Aufzeichnung enthält die Nachrichtenfolge (Integrationstest).
- [x] Der Record-Service gruppiert die Nachrichten nach `LH-FA-18.a` (Flush- und Sync-Gruppen, Pipelining, Abbruch), und jede aufgezeichnete Interaktion besteht `Interaction.Validate` (Unit-Tests).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8) — Review zu `5d666db`, Folge-Review zu `0dbcf9d`; `09cc0e9` und `22b535a` ohne eigenes Review, siehe §7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der Closure von `welle-extended-query` geprüft.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/pgwire` | update | Neue Nachrichten; `recordSitzung` mit zwei Richtungen, `replaySitzung` |
| `internal/adapters/driven/postgres` | update | Neue Nachrichten; `Send`, `Receive`, nicht blockierendes `Close` |
| `internal/hexagon/ports/driving`, `…/ports/driven` | update | Ereignis-Operationen des Record-Use-Case und `Send`/`Receive` mit Gleichzeitigkeits- und Abbruchvertrag ([ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md)) |
| `internal/hexagon/model` | update | `SessionEnd` als Ereignis der Verbindung, `ErrSessionEnded` |
| `internal/hexagon/services` | update | Record-Service: Zustand je Session unter eigener Sperre, Gruppieren, Pipelining, `Validate` vor dem Übernehmen, Einstufung des Endes, Herunterfahren |
| Tests je Schicht | update | Fakes mit begrenztem Puffer (Gegendruck), Abbruchvertrag beider Ports |
| `test/integration` | update | pgx im Standardmodus, Pipeline mit Flush- und Sync-Gruppen, Abbruch ohne `Sync`, Gegendruck mit großer Ausgabe und großem Parameter, SIGTERM nach Blockade |
| `spec/spezifikation.md`, `spec/architecture.md` | update | Randformen in `LH-FA-18.a` und `LH-FA-13.a` (bei Closure beobachtbar gefasst, Validierung Frage 2); Ports, Sequenz, Adapter-Verantwortung, Nebenläufigkeit |
| `docs/plan/planning/open/slice-v1-abschluss-cancel-ohne-schluessel.md`, `docs/plan/planning/welle-v1-abschluss.md` | new, update | Folge-Slice für F-310, eingetragen in die Welle (Slice-Tabelle, Abhängigkeit) |
| `Dockerfile` (Stufe `test`), `harness/README.md` §Sensors | update | Formatierungsprüfung als Teil von `make test` (Folge-Review F-314), Herkunfts-Anker bei Closure |
| `docs/user/abdeckung-*.md` | update | `make abdeckung` |
| `docs/user/benutzerhandbuch.md` | update | Hinweise: einfache Anfrage in laufender Extended-Folge, Beenden mit laufender Folge; Warten ohne Frist, Stopp-Frist im Container, zweites Signal (Validierung Frage 1) |
| `docs/reviews/2026-10-05-mutationen-slice-extended-query-record.md` | new | Mutationstabelle (Runden eins bis drei und Nachtrag aus der Verifikation) |
| `docs/reviews/2026-10-05-review-slice-extended-query-record.md`, `docs/reviews/2026-10-05-folge-review-slice-extended-query-record.md`, `docs/reviews/2026-10-05-verifikation-slice-extended-query-record.md` | new | Review, Folge-Review und Verifikation, mitcommittet |
| `docs/plan/adr/0030-full-duplex-im-record-pfad.md`, `docs/plan/adr/README.md` | new, update | [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md) und ADR-Index, mitcommittet |
| `docs/plan/planning/done/slice-extended-query-replay.md` | update | Folge-Slice: §1 (geliefert, `replaySitzung`), §6 (Risiko aus F-308) |
| `docs/plan/planning/observations/BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/evidence/slice-extended-query-record.md`, `…/spec-randform-erst-im-review-entschieden/evidence/slice-extended-query-record.md`, `…/zusage-im-kommentar-weiter-als-pruefung/evidence/slice-extended-query-record.md`, `…/plan-folgt-korrektur-nicht/evidence/slice-extended-query-record.md` | new | Register-Belege (F-301/F-303/F-313/V-22, F-305/F-311/F-317/F-318/Validierung, F-311/F-319/V-22, F-307/F-315/V-23) |
| `docs/reviews/2026-10-05-validierung-slice-extended-query-record.md` | new | Validierung (Herunterfahren), mitcommittet bei Closure |
| `docs/plan/planning/open/slice-v1-abschluss-betrieb.md` | update | Folge-Slice: Frist für das Herunterfahren in §1 und §6 (F-309, F-317, Validierung Frage 1) |
| `README.md` | update | Stand „Was kann ich heute tun?“: Extended Query im Record (Validierung, Nebenbefund) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-extended-query-modell` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: der Slice verlangt mehr als drei Liefer-Punkte — zurück zur Zerlegung.
- `in-progress` → `open`: Der Upstream-Adapter braucht eine neue Port-Operation — Entscheidung klären. Die Bedingung trat in der Sache ein: Der Upstream-Port bekam einen neuen Gleichzeitigkeitsvertrag (Review F-302). Vollzogen wurde die Rückführung nicht, weil [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md) vor dem nächsten Code-Commit angenommen wurde; der Slice blieb in `in-progress/`.


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

- Asynchrone Nachrichtenfolgen (Pipelining) sind in der Aufzeichnung nicht eindeutig geordnet. Stand: Der Service ordnet Server-Nachrichten der ältesten laufenden Interaktion zu, bis deren `ReadyForQuery` eintrifft, auch wenn Client-Nachrichten der nächsten schon übergeben sind (`TestRecordExtendedPipelining`, `TestE2ERecordExtendedPipeline`); zeitabhängig bleibt nur die Flush-Gruppe ohne wartenden Client, die `LH-FA-18.a` als Grenze von v1 nennt — **Ausgang:** eingetreten, begrenzt: Eindeutig geordnet ist jede Sync-Gruppe und jede Flush-Gruppe mit wartendem Client; die Flush-Gruppe ohne wartenden Client ist die in `LH-FA-18.a` benannte Grenze von v1, ihre Folge für die Wiedergabe trägt `slice-extended-query-replay` §6 (Review F-308).

- `Interaction.Validate` prüft nicht, ob eine Client-Nachricht nur die Felder ihres Typs trägt; der YAML-Schreiber gibt nur diese aus, ein fremd belegtes Feld fiele beim Schreiben still weg. Ebenso schreibt er `Response.ParamTypes` nur an `parameter_description`; an jeder anderen Antwort fiele es still weg. Stand: Die Abbildungen `toClientMessage` (PGWire-Adapter) und `toResponse` (Upstream-Adapter) belegen nur die Felder des Typs, geprüft per Gleichheit des ganzen Modellwerts (`TestToClientMessage`, `TestExtendedSyncGruppe`, `TestSendUndReceive`); `Validate` selbst prüft es weiterhin nicht (aus `slice-extended-query-modell`) — **Ausgang:** entfallen: Die beiden Stellen, die Modellwerte aus PGWire-Nachrichten erzeugen, belegen nur die Felder des Typs (Belege oben), und der YAML-Leser lehnt fremde Schlüssel ab; einen Weg, auf dem ein fremd belegtes Feld den Schreiber erreicht, gibt es nicht.

- Der YAML-Schreiber prüft nicht mit `Interaction.Validate`: eine fehlerhaft gruppierte Interaktion (etwa zwei `sync` in einer Gruppe) wird geschrieben und fällt erst beim Laden im Replay als beschädigt auf. Stand: Der Record-Service prüft jede Interaktion, einfach und Extended, mit `Validate`, bevor er sie übernimmt; eine abgelehnte ist `PGR-E6001` (`TestRecordExtendedNichtDarstellbar`). Der Schreiber selbst prüft weiterhin nicht (aus `slice-extended-query-modell`, Review F-292) — **Ausgang:** entfallen: Jede Interaktion, die der Schreiber erhält, kommt aus dem Record-Service und hat `Validate` bestanden; die Mutation, die das Ergebnis von `Validate` ignoriert, ist rot (M3, von der Verifikation nachgefahren).

- Selbstblockade bei großer Gruppe mit großer Ausgabe (Review F-301, HIGH): Das synchrone Senden einer Gruppe hielt den Empfang an; Server und Recorder blockierten sich, und die Session endete weder beim Schließen des Clients noch nach `SIGTERM`. Stand: behoben nach [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md) mit zwei unabhängigen Richtungen; Belege `TestRecordGegendruck`, `TestExtendedZweiRichtungen`, `TestSendUndReceiveGleichzeitig`, `TestE2ERecordExtendedGegendruck` und `TestE2ERecordExtendedSigtermNachBlockade` (die beiden E2E-Tests sind gegen `5d666db` rot, jetzt grün) — **Ausgang:** eingetreten: behoben nach [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md), vom Folge-Review und von der Verifikation bestätigt (dazu Sonden V-S1 und V-S5).

- Schon gelesene Client-Nachricht beim Herunterfahren (Review F-305): Entschieden in `LH-FA-13.a`, sie wird verarbeitet; `clientRichtung` verarbeitet jede gelesene Nachricht und prüft das Ende erst beim nächsten Lesen (`TestExtendedHerunterfahren`, „gelesene Anfrage“) — **Ausgang:** eingetreten: entschieden in `LH-FA-13.a`; bei Closure beobachtbar gefasst — eine Client-Nachricht um das Signal herum wird entweder noch verarbeitet oder nicht weitergeleitet (Validierung, Frage 2; Verifikation V-24).

- Herunterfahren ohne Obergrenze (Review F-309): Eine Session mit laufender Extended-Interaktion wartet beim Herunterfahren auf deren `ReadyForQuery` (`LH-FA-18.a`, `LH-FA-13.a`); ein Client, der nach einem `Flush` ruht, etwa eine Verbindung in einem Pool, hält den Prozess beliebig lange. Frage an den Validator: Trifft unbegrenztes Warten den Bedarf in CI ([`LH-FA-15`](../../../../spec/lastenheft.md#lh-fa-15--ci-eignung))? Dazu (Folge-Review F-317): Die Sätze in `LH-FA-18.a` („sobald der Recorder das Ende bemerkt“) und `LH-FA-13.a` („schon vollständig gelesen“) binden das Verhalten an innere Zeitpunkte ohne Obergrenze; steht die Client-Richtung in `Send` an einem Server, der weder liest noch antwortet, bemerkt niemand einen geschlossenen Client, bis der Server weitermacht. Geprüft ist nur der Fall mit Datenfluss (`TestE2ERecordExtendedSigtermNachBlockade`). Weiteres Material (Verifikation V-24): Der „Beginn des Herunterfahrens“ ist im Code der erste `Shutdown`-Aufruf der Client-Richtung, nicht das Signal. Eine Nachricht, die zwischen Signal und Wecken geliefert wird, wird noch verarbeitet, auch eine neue `Query`; eine Nachricht, die vollständig im Lesepuffer liegt, aber noch nicht geliefert ist, wird verworfen, wenn keine Interaktion läuft. Validierung (Frage 1): trifft den Bedarf nicht; im Container endet das Warten mit `SIGKILL`, und die Aufzeichnung jeder noch wartenden Session fehlt ganz, ohne Log-Zeile — **Ausgang:** weiter offen → `slice-v1-abschluss-betrieb` §1 und §6 (Frist ab dem Signal, Zwangsende wie Abbruch nach `LH-FA-02.b`, eigener Meldungscode der Klasse 4); die Ergänzung von [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) Boundary wartet auf die Bestätigung des Nutzers. Bis dahin nennt das Benutzerhandbuch (§4, „Eine Anwendung aufzeichnen“) das heutige Verhalten.

- Herunterfahren mit pipelinendem Client (Folge-Review F-311): Die Session nahm nach dem Beginn des Herunterfahrens weiter neue Interaktionen an und endete erst in einer Pause des Clients. Stand: entschieden in `LH-FA-13.a`, nach dem Beginn beginnt keine neue Interaktion; Belege `TestRecordHerunterfahren`, `TestExtendedHerunterfahren` („neue Interaktion nach dem Beginn“), `TestE2ERecordExtendedSigtermBeimPipelining` — **Ausgang:** eingetreten: behoben, von der Verifikation mit Mutation X5 und den Sonden V-S4 und V-S6 bestätigt.

- Verlorenes Aufwecken beim Herunterfahren (Folge-Review F-312): Das Zurücksetzen der Lesefrist konnte ein gleichzeitiges Wecken löschen. Stand: Wecken (Signal und Frist) und Zurücksetzen laufen unter derselben Sperre; wer zurücksetzt, sieht ein Wecken und setzt nicht zurück (`TestWeiterlesenNachWecken`, `TestHerunterfahrenWeckenNichtVerloren`) — **Ausgang:** eingetreten: behoben, von der Verifikation mit den Mutationen W1 und X6 bestätigt.

- Ende einer Session bei einem Client, der nicht liest (Folge-Review F-318): Stand: entschieden in `LH-FA-18.a`, das Ende schließt die Client-Verbindung; ein blockiertes Schreiben endet damit, belegt mit `TestEndeSchliesstVerbindung` (Ende über `Terminate`). Die Fehlerantwort dazu hat eine Frist von einer Sekunde, belegt mit `TestFehlerantwortMitFrist` (nicht unterstützte Nachricht bei nicht lesendem Client; Verifikation V-22, Mutation X1). Die Fehlerantwort kann dabei hinter einem noch laufenden Stapel der Server-Richtung stehen oder bei nicht lesendem Client ausfallen — **Ausgang:** entfallen: vom Bedarf getragen; Klasse und Exit-Code sind gemerkt, bevor die Fehlerantwort geschrieben wird, und eine Sekunde reicht für eine kurze Nachricht an einen lesenden Client (Validierung, Frage 3).

- Gleichzeitigkeit vertraglich am Port: `Send` und `Receive` laufen auf derselben Upstream-Verbindung gleichzeitig, und `Close` beendet beide; der Vertrag steht an `UpstreamSession` und am Record-Use-Case und wird mit Fakes geprüft (`TestSendUndReceiveGleichzeitig`, `TestCloseBeendetWartende`, `TestCloseMitSchreibfrist`, `TestRecordCloseBeendetWartende`, `TestRecordEndeNachErfolgreichemUpstream`). Er trägt, weil `pgproto3.Frontend` Schreib- und Lesepuffer getrennt hält (Race-Detector-Lauf grün, kein Gate); eine andere Bibliotheksversion kann das ändern — **Ausgang:** entfallen für den gepinnten Stand: pgx v5.11.0 (`go.mod`), Race-Detector 20-fach grün (Verifikation), `TestSendUndReceiveGleichzeitig` läuft gegen die echte Bibliothek in jedem `make test`. Dass kein Gate den Race-Detector fährt, steht als Frage an die Closure von `welle-extended-query` in §7.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** Der Recorder vermittelt `Parse`, `Bind`, `Describe`, `Execute`, `Close`, `Flush` und `Sync` und zeichnet sie geordnet auf; pgx im Standardmodus läuft über `record`, und die Aufzeichnung lädt im Replay (`TestE2ERecordExtendedPgx`). Nach [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md) vermittelt der Record-Pfad je Session in zwei unabhängigen Richtungen, und der Record-Service führt den Interaktionszustand allein. Die Verifikation hat alle drei Liefer-Punkte mit eigenen Läufen bestätigt: `make gates`, die Stufe `test` ohne Cache, Race-Detector 20-fach, 18 Mutationen (17 rot, X1 grün → V-22, behoben), sechs Sonden gegen PostgreSQL 17, darunter ein Batch mit 400 Parametern zu 64 KiB und 24 MB Ausgabe.
- **Was ging anders als geplant:** Das Review fand trotz Mutationstabelle ein HIGH (F-301): Das synchrone Senden einer großen Gruppe hielt den Empfang an, Server und Recorder blockierten sich, und der Prozess endete nach `SIGTERM` nicht. Die Tabelle prüfte die Zusagen, die der Plan nannte; Gegendruck war keine davon, und kein Test schickte eine große Gruppe gegen große Ausgabe. Die Behebung verlangte einen neuen Gleichzeitigkeitsvertrag am Upstream-Port. Damit trat die Rückführung `in-progress` → `open` aus §4 in der Sache ein (F-302); aufgelöst wurde sie durch [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md), angenommen vor dem nächsten Code-Commit, und der Slice blieb in `in-progress/` (F-315). Der Umbau brachte das Herunterfahren als neuen Gegenstand: Folge-Review und Verifikation fanden je weitere Randformen (F-305, F-311, F-312, F-317, F-318, V-24), jede ist in `LH-FA-13.a` oder `LH-FA-18.a` entschieden. Die Validierung urteilte rot bei grüner Verifikation: Das Herunterfahren hat keine Obergrenze, und im Container verliert der `SIGKILL` nach der Stopp-Frist die Aufzeichnung jeder noch wartenden Session ganz — eine Probe mit `docker stop -t 10` zeigte es; die Ursache liegt schon im Pfad der einfachen Anfrage. Der Plan war erfüllt, der Bedarf nicht; der Gegenstand geht an `slice-v1-abschluss-betrieb`. `09cc0e9` hat kein eigenes Review; die Verifikation hat ihn mit eigenen Mutationen (X1 bis X7) und Sonden (V-S1 bis V-S6) geprüft. `22b535a` (Test zu V-22, Satz der Sicht zu F-319, Plan) hat weder Review noch Verifikation; ihn belegen die Mutation X1, rot in `TestFehlerantwortMitFrist` (Mutationstabelle, Implementer-Rolle), und `make gates`.
- **Steering-Loop-Eintrag:** Neuer Sensor: Die Stufe `test` des `Dockerfile` prüft die Formatierung aller Go-Dateien, auch der mit Build-Tag (`gofmt -l`; Folge-Review F-314), und läuft in `make test` und `make gates` — liegt in `Dockerfile §Stufe test`, Anker `· seit slice-extended-query-record` im Kommentar der Stufe. Dazu benannte Spec-Lücke, offen: [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) Boundary sagt für das kontrollierte Beenden keine Obergrenze zu (Validierung, Frage 1); die Ergänzung wartet auf die Bestätigung des Nutzers, die Umsetzung trägt `slice-v1-abschluss-betrieb`. `LH-FA-13.a` ist bei Closure beobachtbar gefasst (Validierung, Frage 2; Historienzeile der Spezifikation). Frage an die Closure von `welle-extended-query`: Der Gleichzeitigkeitsvertrag am Upstream-Port hängt an der Bibliothek, und kein Gate fährt den Race-Detector.
- **Beobachtungs-Register (`../observations/`):** Zähler = Dateien unter `evidence/`.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` — Beleg `evidence/slice-extended-query-record.md` (F-301/F-303, F-313, V-22), **4×**.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` — Beleg `evidence/slice-extended-query-record.md` (F-311, F-319, V-22), **4×**.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden/` — Beleg `evidence/slice-extended-query-record.md` (F-305, F-311, F-317/V-24, F-318), **2×**, `offen`.
  - `BEO-REPO/plan-folgt-korrektur-nicht/` (verkörpert in `AGENTS.md` §3.9 seit welle-walking-skeleton) — Wiederauftreten trotz Regel: Plan beschrieb abwesenden Code (F-307), Rückführung ohne vollzogenen Übergang (F-315), §3 ohne alle geänderten Dateien (V-23); neuer Beleg `evidence/slice-extended-query-record.md`, **4×**. Das ist die Frage des Retirement-Checks der Regel.
  
  Die Einträge ab 3× gehen an die Closure von `welle-extended-query`, die über Verkörperung oder Schärfung entscheidet; ihre `state.md` bleibt bis dahin unverändert.
- **Folge-Slices:** `slice-extended-query-replay` (Replay; Folge der Flush-Grenze, F-308), `slice-v1-abschluss-cancel-ohne-schluessel` (F-310), `slice-v1-abschluss-betrieb` (Frist für das Herunterfahren, F-309/F-317, Validierung Frage 1) — alle drei Dateien in `open/`.
- **Risiken aus §6:** zehn, jedes mit genau einem Ausgang — fünf eingetreten, vier entfallen, eines weiter offen mit Folge-Slice; siehe §6.
- **Drei Paarungen:** Repo mit Wellen-Betrieb — geprüft von der Closure von `welle-extended-query`.
- **Belege:** Review `docs/reviews/2026-10-05-review-slice-extended-query-record.md` (`5d666db`; F-301 bis F-310), Folge-Review `docs/reviews/2026-10-05-folge-review-slice-extended-query-record.md` (`0dbcf9d`; F-311 bis F-320), Verifikation `docs/reviews/2026-10-05-verifikation-slice-extended-query-record.md` (`09cc0e9`; DoD 1 bis 3 bestätigt, V-22 bis V-24), Validierung `docs/reviews/2026-10-05-validierung-slice-extended-query-record.md` (`22b535a`; Frage 1 rot, Fragen 2 und 3 grün), Mutationstabelle `docs/reviews/2026-10-05-mutationen-slice-extended-query-record.md`, [ADR-0030](../../adr/0030-full-duplex-im-record-pfad.md). Die Änderungen der Closure (Spezifikation, Handbuch, README, Folge-Slice) belegt `make gates`; der Code ist unverändert.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
