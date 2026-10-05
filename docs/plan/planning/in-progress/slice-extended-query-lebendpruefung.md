# Slice slice-extended-query-lebendpruefung: Lebendprüfungen im Replay

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-extended-query.

**Bezug:** [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-QA-02`](../../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung), [ADR-0007](../../adr/0007-strict-replay.md), [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md)

**Berührte Spec-Stellen:** `LH-FA-09` · `LH-FA-09.a` · `LH-FA-18.a` · `LH-FA-12.a` · `LH-FA-03.a` · `LH-FA-03.b` · `LH-FA-10.a` · `SPEC-011` · `ARC-002`

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Eine Anwendung mit `pgxpool` oder `database/sql` läuft im Replay unabhängig von ihren Pausen: Reine Kommentar- oder Leer-Anfragen, wie sie Pools nach Leerlauf als Lebendprüfung senden (`-- ping`), beantwortet das Replay außerhalb der Reihe, ohne den Cursor zu bewegen.

**Herkunft:** Validierung von `slice-extended-query-replay` (`docs/reviews/2026-10-05-validierung-slice-extended-query-replay.md`, Frage 2, Sonden S4 bis S7). Ob die Lebendprüfung in der Aufzeichnung oder in der Wiedergabe steht, hängt an Pausen über 1 s; das strenge Replay meldet dann `PGR-E5001`. Der Nutzer hat entschieden, solche Anfragen zu tolerieren, statt die Zusage „nur Host und Port“ aus [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) zu verengen.

- Eine einfache Anfrage, deren Text nur aus Leerraum und Kommentaren besteht, beantwortet das Replay zwischen zwei Interaktionen mit `EmptyQueryResponse` und `ReadyForQuery`, wie PostgreSQL es tut; der Cursor bleibt stehen.
- Eine aufgezeichnete Lebendprüfung, die in der Wiedergabe ausbleibt, führt nicht zur Abweichung: Das Replay liest die Aufzeichnung für Cursor, Session-Zuordnung und nicht verbrauchte Interaktionen so, als stünde sie nicht darin ([ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md)); vorhandene Aufzeichnungen mit Lebendprüfungen bleiben ohne Umwandlung verwendbar. Das `ReadyForQuery` der Antwort außerhalb der Reihe trägt den Transaktionsstatus des letzten `ReadyForQuery` auf dieser Verbindung.
- Jede andere einfache Anfrage außer der Reihe bleibt eine Abweichung, auch eine, die mitten in einer Extended-Interaktion kommt (`PGR-E5001`). Diesen Fall nennt `LH-FA-18.a` §Replay heute nicht ausdrücklich (Review F-327 zu `slice-extended-query-replay`); der Satz, der die Ausnahme für Lebendprüfungen schreibt, grenzt ihn mit ab.
- **ADR:** Die Lockerung widerspricht [ADR-0007](../../adr/0007-strict-replay.md) und `LH-FA-09.a` („Kommentare werden nicht ignoriert“). [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) ergänzt [ADR-0007](../../adr/0007-strict-replay.md) um die Ausnahme für Lebendprüfungen (kein `Supersedes`); sie ist angenommen. Die Spezifikation (`LH-FA-09.a`, `LH-FA-18.a` §Replay, `LH-FA-12.a`, `LH-FA-03.a`, `LH-FA-03.b`, `LH-FA-10.a`, `SPEC-011`) folgt ihr im selben Slice; `LH-FA-03.a` (Schritt 4) und `LH-FA-10.a` (Anfrage nach der letzten Interaktion) nennen die Ausnahme, weil sie sonst der Antwort vor der ersten und nach der letzten Interaktion widersprächen. `LH-FA-06.a` bleibt unverändert, weil der Record unverändert bleibt.
- **Serverversion:** Ob der vertikale Tabulator (`\v`) Leerraum einer Lebendprüfung ist, hängt an der Hauptversion in `server_version` der Session: ab PostgreSQL 17 ja, davor nicht (Sonde, §6). Für eine eingehende Lebendprüfung gilt die Version der zugeordneten Session; vor der Zuordnung und ohne lesbare Version ist `\v` kein Leerraum (Entscheidung des Nutzers nach Review F-350, 2026-10-05; die Regel vor der Zuordnung legt `LH-FA-09.a` fest).
- **Startfehler:** Eine Aufzeichnung, deren Sessions nur Lebendprüfungen enthalten, ist `PGR-E3004` (Entscheidung des Nutzers nach Review F-355, 2026-10-05); `LH-FA-03.a` und das Handbuch nennen den Fall.
- **Lastenheft:** Die Negative von [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) nimmt Anfragen aus Leerraum und Kommentaren aus, eingehend wie aufgezeichnet und ausbleibend (Bestätigung des Nutzers, 2026-10-05; die zweite Richtung nach Review F-357); das Lastenheft steht auf Draft, kein Versionssprung.
- **Handbuch:** Der Abschnitt „Mit einem Datenbanktreiber arbeiten“ nennt das Verhalten bei Lebendprüfungen von Pools und `database/sql` (`ShouldPing`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Record und `play` — bleiben bewusst unverändert: Sie zeichnen Lebendprüfungen weiter auf beziehungsweise senden sie, damit Aufzeichnung und Einspielen wiedergeben, was die Anwendung tat; das Replay allein trägt die Toleranz, auch für vorhandene Aufzeichnungen ([ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md), Option B verworfen). Keine Änderung am Recording-Format.
- Anfragen mit Anweisung, auch ein einzelnes `;` — bleiben strikt ([ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md), Option D verworfen); kein beobachteter Treiber sendet sie als Lebendprüfung.
- Toleranz für andere Anfragen (Normalisierung, Kommentare in fachlichen Anfragen, Reihenfolge) — [ADR-0007](../../adr/0007-strict-replay.md) bleibt für alles außer Lebendprüfungen bestehen; Out-of-Scope der Welle.
- Gleichzeitig genutzte Verbindungen eines Pools (Validierung, Sonde S8) — `slice-v1-abschluss-sessions` §6 führt sie.
- Frist beim Herunterfahren im Replay — `slice-v1-abschluss-betrieb`.
- Passwort-Anmeldung beim Aufzeichnen gegen ein Standard-Image (Validierung, Sonde S12) — `slice-v1-abschluss-anmeldung`.
- Lebendprüfungen über das Extended Query Protocol — kein beobachteter Treiber sendet sie; pgx pingt mit einer einfachen Anfrage (Validierung, pgx-Quelltext). Ein solcher Fall wäre ein anderer Vorgang.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay): Eine reine Kommentar- oder Leer-Anfrage zwischen zwei Interaktionen beantwortet das Replay mit `EmptyQueryResponse` und `ReadyForQuery`, ohne den Cursor zu bewegen; eine aufgezeichnete, die ausbleibt, ist keine Abweichung; jede andere einfache Anfrage außer der Reihe, auch mitten in einer Extended-Interaktion, bleibt `PGR-E5001` (Unit-Tests).
- [x] [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-02`](../../../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung): `pgxpool` und `database/sql` mit pgx in der Default-Konfiguration laufen mit Pausen über 1 s beim Aufzeichnen, beim Wiedergeben oder bei beiden ohne Abweichung durch (Integrationstest nach den Lagen S4 bis S7 der Validierung).
- [x] Die neue ADR zu Lebendprüfungen ist angenommen und im ADR-Index; `LH-FA-09.a` und `LH-FA-18.a` §Replay nennen die Ausnahme und die falsche Protokollart; das Benutzerhandbuch nennt das Verhalten bei Lebendprüfungen.
- [x] `make gates` grün — an `aca39ce` (Verifikation, Punkt 4); die Closure ändert nur Planungsdokumente.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8) — Review bis `0841891`, Folge-Review bis `c9ba569`; `aca39ce` ohne Review, von der Verifikation mit Mutationen geprüft, siehe §7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — bei Closure geprüft (§7); die Closure von `welle-extended-query` prüft sie erneut.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/` (neue ADR, Index) | new, update | Einschränkung von [ADR-0007](../../adr/0007-strict-replay.md) für Lebendprüfungen; Architect |
| `spec/spezifikation.md` (`LH-FA-09.a`, `LH-FA-18.a` §Replay, `LH-FA-12.a`, `LH-FA-03.a`, `LH-FA-03.b`, `LH-FA-10.a`, `SPEC-011`, Historie) | update | Ausnahme für Lebendprüfungen (Definition mit den Leerraum-Zeichen des Scanners von PostgreSQL, `\v` nach `server_version` der Session, Antwort, Überspringen aufgezeichneter), ihre Wirkung auf Session-Zuordnung bei `first-request` und `connection`, Startfehler und nicht verbrauchte Interaktionen; Einspielen führt Sessions nur aus Lebendprüfungen aus; falsche Protokollart im Replay (F-327) |
| `spec/lastenheft.md` (`LH-FA-09`) | update | vom Nutzer bestätigt (§6): Negative nimmt eingehende und ausbleibende aufgezeichnete Lebendprüfungen aus |
| `docs/plan/planning/open/slice-v1-abschluss-sessions.md`, `docs/plan/planning/open/slice-v1-abschluss-postgres-versionen.md` | update | Folge-Slices: Lebendprüfungen bei der Zuordnung `connection`; gefundene Versionsgrenze für `\v` |
| `internal/hexagon/services` (Replay-Service; `lebendpruefung.go` neu) | update, new | Erkennung als reine Funktion (`istLebendpruefung`, Tabellentest je Variante mit und ohne `\v`); Version aus `server_version` (`vtLeerraum`, Grenze `vtAbVersion`), je aufgezeichneter Session und je Verbindung nach der Zuordnung; Antwort außerhalb der Reihe mit dem Status des letzten `ReadyForQuery` der Verbindung; aufgezeichnete Lebendprüfungen beim Laden aus den Sessions genommen (Cursor, Zuordnung `first-request`, nicht verbrauchte Interaktionen, `PGR-E3004`); die Meldung nach der letzten erwarteten Interaktion nennt deren aufgezeichnete Nummer, in `Query` und `ClientMessage` (F-356), ohne erwartete Interaktion 0 (F-362); `TestReplayExtendedAmCursor` erwartet für die leere Anfrage `EmptyQueryResponse` |
| `internal/hexagon/ports/driving` | update | nur der Kommentar von `Query` nennt die Antwort außerhalb der Reihe |
| `internal/adapters/driving/pgwire` | keine Änderung | Die Antwort außerhalb der Reihe geht als gewöhnliche Antwortliste über den Port; der Adapter übersetzt `EmptyQueryResponse` schon |
| `test/integration` (`lebendpruefung_e2e_test.go`), `go.mod` | new, update | `pgxpool` und `database/sql` mit pgx: ohne und mit Pausen über 1 s aufgezeichnet, je ohne und mit Pausen wiedergegeben (Lagen S4 bis S7); jede Lebendprüfung der Tabelle gegen PostgreSQL und Replay, die Texte mit `\v` nach `server_version` des Servers, ohne und mit ihnen aufgezeichnet (unter 17 der aufgezeichnete Syntaxfehler, F-360); `go.mod` nimmt `puddle` für `pgxpool` auf, dazu `golang.org/x/sync` als dessen indirekte Abhängigkeit (V-34) |
| `docs/user/benutzerhandbuch.md`, `docs/user/abdeckung-*.md` | update | Verhalten bei Lebendprüfungen („Mit einem Datenbanktreiber arbeiten“, Wiedergabe, Fragen); Abdeckungstabellen |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-extended-query-replay` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die ADR verlangt eine Änderung am Recording-Format (etwa eine Markierung aufgezeichneter Lebendprüfungen) — das ist eine dritte Schicht; Format und Replay werden getrennt geschnitten.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer verwirft die Einschränkung von [ADR-0007](../../adr/0007-strict-replay.md) im Review der ADR — dann gilt der zweite Weg der Validierung (Zusage verengen, Änderung am Lastenheft), und dieser Slice geht mit Grund zurück.

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

- Transaktionsstatus im `ReadyForQuery` der Antwort außerhalb der Reihe: Er muss dem Stand nach der letzten beantworteten Interaktion entsprechen, auch innerhalb einer Transaktion oder nach einem Fehler in ihr. Entschieden in [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md): der Status des letzten `ReadyForQuery` auf der Verbindung, nach dem Handshake `I`. `TestReplayLebendpruefungAusserDerReihe` und `TestReplayLebendpruefungExtended` prüfen `I`, `T` und `E` nach einfachen und Extended-Interaktionen; die Mutationen „Status immer `I`“, „Status nicht aus Extended-Antworten“ und „Status nach dem Handshake leer“ färben sie rot — **Ausgang:** entfallen: Die Regel ist in [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) entschieden und in `cursor.status` umgesetzt; jede der Mutationen „immer `I`“, „nicht aus Extended-Antworten“, „nicht aus einfachen Antworten“ und „nach dem Handshake leer“ färbt einen Unit-Test rot (Verifikation U04, U05, U06, U24), „immer `I`“ auch den Integrationstest gegen 17 (Transaktionsstatus `E` nach `SELECT 1/0`). Ein falscher Status kann damit nicht unbemerkt entstehen.
- Eine fachliche Anfrage, die nur aus Kommentar besteht, ist von einer Lebendprüfung nicht zu unterscheiden. PostgreSQL antwortet auf beide gleich (`EmptyQueryResponse`); ob die Ausnahme damit für jede solche Anfrage gelten darf, entscheidet die ADR. Entschieden in [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md): ja, als akzeptiertes Negativ (keine Zustandsänderung, gleiche Antwort wie PostgreSQL) — **Ausgang:** entfallen: [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) nimmt den Fall als akzeptiertes Negativ an, und das Replay antwortet darauf wie der aufgezeichnete Server: leer und ohne Zustandsänderung. Wo der Server die Anfrage nicht als leer liest (`\v` unter 17), ist sie seit F-350 keine Lebendprüfung und wird verglichen (Risiko 5). Eine Antwort, die nicht die des Servers ist, kann aus der Ununterscheidbarkeit nicht mehr folgen.
- Vorhandene Aufzeichnungen enthalten Lebendprüfungen als Interaktionen (Validierung, Sonde S6). Bleibt die Wiedergabe einer solchen Aufzeichnung verträglich? Entschieden in [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md): Der Record zeichnet weiter auf, das Replay überspringt aufgezeichnete Lebendprüfungen; keine Formatänderung, die Rückführung aus §4 tritt nicht ein — **Ausgang:** entfallen: Das Recording-Format ist unverändert, die Rückführung aus §4 ist nicht eingetreten (Verifikation §2, §4). Dass eine vorhandene Aufzeichnung mit Lebendprüfungen ohne Umwandlung wiedergegeben wird, belegt Lage S6 des Integrationstests; ohne das Überspringen ist sie rot (Verifikation U02, E0). Einzige Ausnahme ist eine Aufzeichnung nur aus Lebendprüfungen, die jetzt mit `PGR-E3004` endet; das hat der Nutzer nach F-355 entschieden, und `LH-FA-03.a` und das Handbuch nennen es.
- Ob [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) im Lastenheft einen Satz braucht oder die Spezifikation genügt, bestätigt der Nutzer vor dem ersten Code-Commit. Empfehlung des Architect: ein Satz in der Negative von `LH-FA-09`, weil eine tolerierte Lebendprüfung sonst „eine Anfrage ohne passende Aufzeichnung“ ist, für die das Lastenheft `LH-FA-10` verlangt, und eine ADR das Lastenheft nicht schärfen darf. Der Nutzer hat den Satz am 2026-10-05 bestätigt; er steht in der Negative von `LH-FA-09` — **Ausgang:** entfallen: Der Satz steht in der Negative von [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), für die eingehende (`eb61782`) und nach F-357 auch für die ausbleibende aufgezeichnete Lebendprüfung (`e2871a3`); das Lastenheft steht auf Draft, kein Versionssprung. Die Verifikation hat ihn bestätigt (Punkt 3). Die Frage ist beantwortet.
- Die Leerraum-Zeichen in `LH-FA-09.a` hängen an der Serverversion. Eingetreten im Review (F-350): PostgreSQL 16 beantwortet eine Anfrage mit `\v` außerhalb eines Kommentars mit einem Syntaxfehler. Sonde vom 2026-10-05 mit `psql -c` gegen `postgres:<n>-alpine`, je per Digest: 14.24 (`sha256:4ea9e5ed…`), 15.19 (`sha256:f7d23353…`) und 16.15 (`sha256:721873c3…`) lehnen `\v` und `\v-- ping\v` mit `syntax error` ab; 17.11 (`sha256:b0f9560a…`, das Image der Integrationstests) und 18.6 (`sha256:77f58511…`) antworten leer. Die übrigen Fälle (`\f`, `\v` im Zeilenkommentar, verschachtelte Blockkommentare, `\r` als Ende eines Zeilenkommentars) verhalten sich in allen fünf gleich. `LH-FA-09.a` bindet `\v` deshalb an `server_version` der Session, ab 17. `TestE2EReplayLebendpruefungWiePostgres` prüft die Grenze am jeweiligen Server; gelaufen ist die volle Integrationssuite gegen 14, 16, 17 und 18 — **Ausgang:** eingetreten → `slice-v1-abschluss-postgres-versionen`. Im Slice behoben: `LH-FA-09.a` bindet `\v` an die Hauptversion in `server_version`, ab 17, vor der Zuordnung nie (`4149080`, `98f4140`); der Integrationstest prüft die Grenze seit `aca39ce` auch von unten (F-360), Ende zu Ende gegen 16 und 17 mit den Mutanten „Grenze 16“ und „Grenze 18“ rot (Verifikation 1e). Den Rest trägt der Folge-Slice: Das Gate fährt nur gegen 17, eine zu niedrige Grenze fängt dort nur der Unit-Test (V-33); die Matrix über alle unterstützten Versionen und die Prüfung der übrigen Leerraum-Zeichen je Version stehen in dessen §1 und §6.

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

- **Was hat funktioniert:** Das Replay beantwortet Lebendprüfungen von `pgxpool` und `database/sql` außerhalb der Reihe und überspringt aufgezeichnete, die ausbleiben; Aufzeichnung und Wiedergabe hängen für diese Anfragen nicht mehr an Pausen des Pools. Die Randformen, die §6 vor dem Code nannte (Transaktionsstatus, fachliche Kommentar-Anfrage, vorhandene Aufzeichnungen, Satz im Lastenheft), hat [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) vor der ersten Code-Zeile entschieden; keine davon hat ein Review wieder geöffnet. Die Verifikation hat alle drei Liefer-Punkte mit eigenen Läufen bestätigt: 24 von 24 Unit-Mutationen rot, Ende zu Ende gegen 14.24, 16.15, 17.11 und 18.6, `make gates` an `aca39ce` mit Exit 0. Bewusst gebrochen (V-35): Im Stand vor dem Slice sind die Lagen S5 und S6 für beide Treiber rot, und die abgeschaltete Antwort außerhalb der Reihe färbt S5, das abgeschaltete Überspringen S6; diese beiden Lagen belegen den Fix. Lage S7 („Pausen bei beidem“) ist ohne den Fix grün, weil sich Lebendprüfungen in Aufzeichnung und Wiedergabe dort decken; sie ist Regressionsschutz, kein Beleg, und wer die Lagen kürzt, darf S5 und S6 nicht streichen.
- **Was ging anders als geplant:** Das Review fand eine versionsabhängige Eigenschaft des Servers (F-350, HIGH): `\v` ist erst ab PostgreSQL 17 Leerraum. Die Spezifikation zählte ihn ohne Version, der Kommentar an `leerraum` sagte „die Zeichen, die der Scanner von PostgreSQL liest“ zu, und das Repo testete nur gegen 17; gegen 16 lieferte das Replay eine Antwort, die weder aus der Aufzeichnung stammte noch der des Servers entsprach. Der Nutzer hat entschieden, `\v` an `server_version` der Session zu binden. Zwei Prüfrunden und eine Verifikation statt einer Runde. Erneut folgte ein Folge-Slice der Korrektur nicht vollständig (F-351, MEDIUM; Rest F-362), obwohl `AGENTS.md` §3.9 es verlangt; §3 nannte `golang.org/x/sync` nicht und ist mit dieser Closure nachgezogen (V-34). `aca39ce` (F-360 bis F-362) hat kein Review; die Verifikation hat ihn mit eigenen Mutationen geprüft (U18, U19, U07 Ende zu Ende gegen 16), DoD-Punkte 1 und 2 hängen an diesen Läufen.
  - **F-352, begründet nicht umgesetzt:** `Schärft:` von [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) nennt `LH-FA-03.a` und `LH-FA-10.a` nicht, obwohl beide unter der ADR geändert wurden, `LH-FA-03.a` zweimal (`7ffd029`, `4149080`). Das Feld ist die Aufwärts-Deklaration der ADR und damit nach der Annahme Inhalt, nicht Metadatum wie die Geschichtszeile (F-353); es zu ergänzen verstieße gegen `AGENTS.md` §3.5. Eine Folge-ADR nur für die Spur wäre unverhältnismäßig: Beide Stellen folgen aus Sätzen der ADR („löst keine Session-Zuordnung aus“, „auch vor der ersten und nach der letzten“), ohne sie zu erweitern. Folge: `make doc-trace` führt von `LH-FA-03.a` und `LH-FA-10.a` nicht zu [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md).
  - **F-363, Begründung der Regel vor der Zuordnung:** Vor der Zuordnung gilt `\v` konservativ nicht als Leerraum, statt die Version aus dem Handshake zu nehmen. Die Version des Handshakes ist die der Session `frei[0]` beim Verbindungsaufbau; die Session, die die Verbindung bei ihrer ersten Anfrage erhält, kann eine andere sein. Bei verschiedenen Versionen in einer Aufzeichnung gäbe die Handshake-Regel eine leere Antwort, wo der aufgezeichnete Server einen Fehler sendete — die Klasse von F-350, eine verschluckte Fehlerantwort. Die konservative Regel meldet den Fehler stattdessen laut (`PGR-E5001`), wie [ADR-0007](../../adr/0007-strict-replay.md) es für jede andere Abweichung tut. Ihr Preis (Folge-Review, Sonde S3): Eine Lebendprüfung mit `\v` vor der ersten Anfrage gegen eine Aufzeichnung ab 17 endet mit `PGR-E5001`, wo PostgreSQL leer antwortet; kein beobachteter Treiber sendet sie, und das Handbuch nennt den Fall (F-361). Gilt für `slice-v1-abschluss-sessions`, damit die Zuordnung `connection` die Frage nicht neu entscheidet.
- **Steering-Loop-Eintrag:** Benannte Spec-Lücke, adressiert: Seit [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) bildet das Replay eine Eigenschaft des PostgreSQL-Scanners nach, aber Lastenheft und Spezifikation nennen keine unterstützte Serverversion (`LH-FA-05.e`: „jeder Server, der PGWire 3.0 spricht“), und das Gate prüft nur gegen 17; eine versionsabhängige Eigenschaft fiel deshalb erst im Review als HIGH auf (F-350), und eine zu niedrige Grenze bliebe im Gate weiter Ende zu Ende grün (V-33). Adresse: `slice-v1-abschluss-postgres-versionen` (Satz in [`LH-RB-02`](../../../../spec/lastenheft.md#lh-rb-02--einsatzumgebung), feste Versionsliste in der Spezifikation, Integrationsmatrix). **Warum dieser Eintrag und nicht Kandidat b:** F-351 ist die sechste Evidenz für `BEO-REPO/plan-folgt-korrektur-nicht`, eine Beobachtung, die schon `verkörpert` ist (`AGENTS.md` §3.9). Ob die Regel geschärft oder durch einen Sensor ergänzt wird, entscheidet die Closure von `welle-extended-query` über alle sechs Vorgänge; „geschärfte Regel“ als Eintrag dieses Slice behauptete eine Schärfung, die nicht stattgefunden hat. Die Versionslücke dagegen ist Originalinformation dieses Slice, eine Lücke im Vertrag und nicht im Ablauf, und sie hat eine Adresse, die sie annimmt (dort §1, §6 Risiko *Versionsgrenze für `\v`*). Kandidat b steht im Register.
- **Beobachtungs-Register (`../observations/`):** Zähler = Dateien unter `evidence/`; je Eintrag Beleg `evidence/slice-extended-query-lebendpruefung.md`.
  - `BEO-REPO/plan-folgt-korrektur-nicht/` (verkörpert in `AGENTS.md` §3.9 seit welle-walking-skeleton) — Folge-Slice nach Spec-Ergänzung nicht nachgezogen (F-351), Kopplung im Folge-Slice ohne die gebrochene Invariante (F-362), §3 ohne eine geänderte Datei (V-34), **6×**. Retirement-Check: wieder aufgetreten, die Regel bleibt; nach `slice-extended-query-record` und `slice-extended-query-replay` der dritte Slice dieser Welle in Folge.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` — Nummer nach dem Ende in `ClientMessage` umgesetzt, aber ungeprüft (F-356, Mutation M9 grün), **6×**.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` — Kommentar an `leerraum` sagte die Zeichen des Scanners ohne Version zu (F-350); Abdeckungs-Deklaration und Folge-Slice sagten für die Versionsgrenze mehr zu, als der Integrationstest prüfte (F-360), **6×**.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden/` — `\v` je Serverversion (F-350) und Startfehler bei einer Aufzeichnung nur aus Lebendprüfungen (F-355; die Ableitung meldete der Implementer, Spezifikation und Handbuch fand erst das Review), beide vom Nutzer nach dem Review entschieden; die Ausnahme im Lastenheft nur in eine Richtung (F-357), **4×**. Die vor dem Code genannten Randformen aus §6 hielten, die nicht genannten nicht.
  - `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft/` — neu: F-350, V-33, **1×**.
  - `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients/` — kein neuer Fund, **1×**, `offen`. Dieser Slice behebt den ersten Fund (Lebendprüfungen); F-358 (Bedingung `nMit > nOhne` unter Last) betrifft den Test, nicht das Replay, und ist nicht beobachtet. Der zweite Fund (Flush-Gruppe mit später Antwort) bleibt offen; der Eintrag kann deshalb nicht `gestrichen` werden.

  Über der Schwelle stehen `plan-folgt-korrektur-nicht` (6×, verkörpert), `negativtests-fehlen-bei-neuem-vertrag` (6×), `zusage-im-kommentar-weiter-als-pruefung` (6×) und `spec-randform-erst-im-review-entschieden` (4×). Ihre Ausgänge vergibt die Closure von `welle-extended-query`; die `state.md` bleibt bis dahin unverändert. Einmalig und nicht eingetragen: F-352 (Schärft-Feld, begründet oben), F-353 (Geschichtszeile der ADR), F-354 (Satz durch Umstellung in neuen Bezug), F-361 (Handbuch ohne Ausnahme), F-358, F-359 und F-363 (kein Fehlermuster).
- **Folge-Slices:** `slice-v1-abschluss-postgres-versionen` (Matrix über alle unterstützten Versionen, Grenze für `\v` je Version, V-33; Risiko 5), `slice-v1-abschluss-sessions` (Lebendprüfungen bei `--session-assignment connection`, Kopplung an `frei` und `letzteNummer`, F-351, F-362; Regel vor der Zuordnung nach F-363) — beide Dateien in `open/`.
- **Risiken aus §6:** fünf, jedes mit genau einem Ausgang — eines eingetreten (→ `slice-v1-abschluss-postgres-versionen`), vier entfallen; siehe §6.
- **Drei Paarungen:** Anker — kein `liegt in`, mit diesem Slice ist nichts verkörpert. Folge-Slice — die zwei genannten Dateien liegen in `open/`, und beide nennen diesen Slice als Herkunft. Register — die sechs genannten Kennungen bestehen als Verzeichnis mit nicht leerem `evidence/`. Die Closure von `welle-extended-query` prüft erneut.
- **Belege:** Review `docs/reviews/2026-10-05-review-slice-extended-query-lebendpruefung.md` (bis `0841891`; F-350 bis F-359), Folge-Review `…-folge-review-slice-extended-query-lebendpruefung.md` (bis `c9ba569`; F-360 bis F-363), Verifikation `docs/reviews/2026-10-05-verifikation-slice-extended-query-lebendpruefung.md` (bis `aca39ce`; V-33 bis V-35; DoD 1 bis 4 bestätigt). `make gates` grün an `aca39ce`; die Closure ändert nur Planungsdokumente.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Stand der Closure von `slice-extended-query-replay`, 2026-10-05); alle Einträge liegen in `REPO`. Einschlägig:

- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` — 1×, `offen`; der Gegenstand dieses Slice ist sein erster Fund.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` — 3×, über der Schwelle; die Randformen in §6 gehören vor den Code in die ADR, nicht in eine Prüfrunde.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — je 5×; die Ausnahme ist ein neuer Vertrag am Replay, ihre Grenze (welche Anfrage, welcher Cursor) braucht Tests, die eine Mutation fangen.
- `BEO-REPO/plan-folgt-korrektur-nicht` — 5×, `verkörpert` (`AGENTS.md` §3.9); gilt als Regel, auch für den Kopf.

Die Einträge über der Schwelle bekommen ihren Ausgang bei der Closure von `welle-extended-query`, nicht hier.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Spezifikation führt, der Code folgt ihr; `harness/conventions.md` deklariert `*` als Greenfield).
