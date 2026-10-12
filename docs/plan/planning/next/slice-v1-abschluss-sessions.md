# Slice slice-v1-abschluss-sessions: Mehrere Sessions und Verbindungsfehler

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md)

**Berührte Spec-Stellen:** `LH-FA-12.a` · `LH-FA-17.a` · `LH-FA-09.a` · `LH-FA-03.b` · `LH-FA-13.b` · `LH-FA-02.b` · `LH-FA-21.a` · `SPEC-017` · `SPEC-028` · `SPEC-034` · `SPEC-043` · `ARC-002`

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Bereits geliefert** vom Walking Skeleton: gleichzeitige Verbindungen im Record, Übernahme in der Reihenfolge der Sessionenden, Zuordnung `first-request` im Replay. Dieser Slice ergänzt `--record-empty-sessions`, `--session-assignment connection` und belegt LH-FA-12 vollständig. Von `slice-extended-query-lebendpruefung` ([ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md), `LH-FA-09.a`): Lebendprüfungen (Anfragen nur aus Leerraum und Kommentaren) beantwortet das Replay außerhalb der Reihe; sie lösen keine Zuordnung aus, auch nicht, wenn keine Session mehr frei ist, und eine Verbindung, die nur Lebendprüfungen stellt, ist eine Verbindung ohne Anfrage. Aufgezeichnete Lebendprüfungen überspringt das Replay; eine Session nur aus ihnen ist für das Replay eine Session ohne Interaktion und bei `first-request` übersprungen.

**Kopplung an `--session-assignment connection`:** `NewReplayService` nimmt Sessions ohne Interaktion, auch solche nur aus Lebendprüfungen, schon beim Laden aus der Liste der freien Sessions (`frei`). `LH-FA-12.a` verlangt für `connection` dagegen, dass die n-te Verbindung die Session mit der `id` n erhält, auch eine Session ohne Interaktion und eine nur aus Lebendprüfungen. Die Zuordnung nach `id` muss diese Sessions deshalb weiter führen; die Liste `frei` reicht dafür nicht. Für eingehende Lebendprüfungen gilt bei `connection` vom Verbindungsaufbau an die Serverversion der zugeordneten Session (`server_version`, `LH-FA-09.a`).

**Übernommen aus `slice-replay-semantik-mismatch`:** Eine Session ohne Interaktion, auch eine nur aus Lebendprüfungen, ist mit und ohne `--fail-on-unconsumed` weder nicht verbraucht noch nie zugeordnet, auch bei `connection` (`LH-FA-03.b` §Verbraucht). Für `first-request` prüft das `TestReplayNichtsZuMelden`; für `connection` gehört der Test hierher, weil die Zuordnung nach `id` diese Sessions weiter führt und `Unassigned` sie dann nicht als nie zugeordnet zählen darf.

**Übernommen von `slice-v1-abschluss-sqlite-format`** (dort §1, *Abgegeben*, und DoD; Verifikation V-174): Die Wirkung von `--record-empty-sessions` mit `--format sqlite` am Binary. Die Option gibt es erst mit diesem Slice, und `meta.empty_sessions` des Schreibens von `sqlite` (`SPEC-043`) ist bis dahin nur im Adapter-Test erreichbar. Der Test steht in DoD-Punkt 1; geschrieben wird mit `--format sqlite`, gelesen über `replay` (`slice-v1-abschluss-sqlite-format-lesen` liegt in der Reihenfolge davor).

**Gebunden durch `slice-v1-abschluss-zeitangaben-aufzeichnen`** (dort §6 *Nicht absteigend* und §7; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13):

1. Die Option `--record-empty-sessions` und ihr Schlüssel `record_empty_sessions` entstehen mit diesem Slice. Der Fall „Option, die der Stand nicht kennt“ in `TestDateiUngueltig` (`internal/adapters/driving/cli/datei_test.go`) benutzt heute den Schlüssel `record_empty_sessions` als Platzhalter einer unbekannten Option und wird mit der Lieferung laut rot. Dieser Slice stellt den Fall auf einen Schlüssel um, den auch spätere Slices nicht liefern, und führt `record_empty_sessions` mit eigenen Fällen (gültiger Wert, Ebene `replay:`/`play:`/oben, ungültiger Wert nach der Optionstabelle), wie `slice-v1-abschluss-zeitangaben-aufzeichnen` es für `record_timing` tat. Zum Test in DoD-Punkt 1.
2. Die Zusage „Die Werte einer Session steigen nicht ab“ (`LH-FA-21.a` *Messpunkt*, `--record-timing`) gilt nur, wenn die Aufrufe der Client-Richtung einer Session (`Query`, `ClientMessage`) nacheinander kommen. Nebenläufig werden die Verbindungen (je eine Session mit eigener Schleife), nicht die Aufrufe einer Session. Dieser Slice hält das ein; kein Test prüft es (akzeptiertes Negativ im Geber, ein Adapter-Test wäre dort eine dritte Schicht). Zeigt sich im Code, dass die Nebenläufigkeit Aufrufe einer Session überlappen lässt, ist das ein Plan-Wechsel (§4).

**Gebunden durch `slice-v1-abschluss-umgebung-anderes-kommando`** (dort §7, Paarungen, offene Frage zu V-3; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): `--session-assignment` (bei `replay`) und `--record-empty-sessions` (bei `record`) liest der Stand noch nicht; `nichtGelesen()` in `internal/adapters/driving/cli/fremd_test.go` führt sie. Dieser Slice liefert den Leser, kürzt `nichtGelesen()` um diese Option(en) und bindet ihn in `liesUmgebung` (`internal/adapters/driving/cli/cli.go`) an; die Zuordnung `zuordnung()` nennt sie für `replay` (`--session-assignment`) und `record` (`--record-empty-sessions`) (Spalte *Kommando* der Optionstabelle, `LH-FA-17.a`). Sonst wird `TestLeserUmgebungNichtGelesenerOptionen` rot, weil die Variable beim eigenen Kommando dann geprüft wird (`TestLeserZuordnungGleichTabelle` bleibt grün). Beim anderen Kommando (bei `--session-assignment` `record` und `play`, bei `--record-empty-sessions` `replay` und `play`) ist die Variable schon jetzt `PGR-E2001`.

**Ziel:** Mehrere Client-Verbindungen werden im Record parallel als eigene Sessions aufgezeichnet (Verbindungen ohne Anfrage nicht); im Replay erhält die n-te Verbindung mit einer Anfrage die n-te Session; Verbindungsende und Verbindungsfehler verhalten sich wie spezifiziert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-sessions-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- Eine deterministische Zuordnung bei gleichzeitigem Verbindungsaufbau — Out-of-Scope von LH-FA-12.
- Signalbehandlung — `slice-v1-abschluss-herunterfahren` (aus `slice-v1-abschluss-betrieb` hervorgegangen).
- Das Ende des Servers (Fehlerantwort `FATAL` oder `PANIC`, danach das Verbindungsende) als Ereignis einer Session — `slice-v1-abschluss-verbindungsende-record` (dort §1, *Übernommen aus `slice-v1-abschluss-sessions`*) für das Aufzeichnen und `slice-v1-abschluss-verbindungsende-replay` für die Wiedergabe und die Zuordnung (Entscheidung des Nutzers vom 2026-10-12, Option B, [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md); Sendung des Planners, `AGENTS.md` §3.13): „Verbindungsende und Verbindungsfehler“ im Ziel dieses Slice meint das Ende der Client-Verbindung und die Verbindungsfehler nach `LH-FA-13.b`, nicht ein Ende, das der Server nach einer Fehlerantwort herbeiführt. Dieser Slice liefert die Verwaltung der Sessions (Kennung nach Sessionende, `--record-empty-sessions`, `connection`), die jene Slices unverändert benutzen; sie stehen in der Reihenfolge dahinter.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen): Zwei parallele Verbindungen erzeugen zwei Sessions mit je geordneten Interaktionen, eine Verbindung ohne Anfrage keine; im Replay erhält die n-te Verbindung mit einer Anfrage die n-te Session, eine Anfrage darüber hinaus ist ein Mismatch, eine Lebendprüfung nicht (`LH-FA-09.a`); bei `--session-assignment connection` erhält die n-te Verbindung die Session mit der `id` n auch dann, wenn diese nur Lebendprüfungen enthält; `--record-empty-sessions` und `--session-assignment connection` verhalten sich wie spezifiziert; mit `--fail-on-unconsumed` ist eine solche Session bei `connection` weder nicht verbraucht noch nie zugeordnet (`LH-FA-03.b`); `--record-empty-sessions` mit `--format sqlite` schreibt `meta.empty_sessions` mit dem Wert `true` und die leere Session, und `replay` liest die Datei wie dieselbe Aufzeichnung als YAML (am Binary; aus `slice-v1-abschluss-sqlite-format`, V-174); die Leser beider Optionen (auch ihrer Umgebungsvariablen) sind angebunden und `nichtGelesen()` in `fremd_test.go` führt sie nicht mehr (aus `slice-v1-abschluss-umgebung-anderes-kommando`) (Test).
- [ ] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Ein Verbindungsfehler beendet nur die Verbindung, der Prozess merkt sich die Klasse (Test).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services` | update | Session-Verwaltung, Recording-Zustand mit Synchronisierung; Zuordnung `connection` nach `id`, auch für Sessions ohne Interaktion und solche nur aus Lebendprüfungen, die `NewReplayService` heute aus `frei` nimmt |
| `internal/adapters/driving/cli` | update | Optionen `--record-empty-sessions` und `--session-assignment` mit Umgebungsvariable, Schlüssel, `config show`, Hilfe; Fälle von `record_empty_sessions` in `TestDateiUngueltig` (Sendung aus `slice-v1-abschluss-zeitangaben-aufzeichnen`) |
| `internal/adapters/driving/pgwire` | update | Verbindungen nebenläufig; die Client-Richtung je Session bleibt eine Schleife (Bindung aus `slice-v1-abschluss-zeitangaben-aufzeichnen`) |
| `test/integration` | update | Happy/Boundary/Negative |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-replay-semantik` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Parallelität verlangt eine Änderung des Recording-Formats — zurück zur Zerlegung.
- `in-progress` → `next`: Der PGWire-Adapter ändert sich (dritte Schicht neben Kern und CLI-Adapter, siehe §8) oder die Nebenläufigkeit lässt Aufrufe einer Session überlappen (Bindung aus `slice-v1-abschluss-zeitangaben-aufzeichnen`) — zurück zur Zerlegung, Schnitt vor dem Anspruch.
- `in-progress` → `open`: Die Zuordnungsregel (n-te Verbindung mit Anfrage, n-te Session) ändert sich durch eine Entscheidung des Auftraggebers — Carveout.

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

**Randformen** (`AGENTS.md` §3.12), soweit aus Bindungen anderer Slices bekannt: *Aufrufe einer Session* — nacheinander, nie gleichzeitig; entschieden im Geber `slice-v1-abschluss-zeitangaben-aufzeichnen` (Spezifikation `LH-FA-21.a` *Messpunkt*, Sicht Zeile *Record-Use-Case*), hier als Bindung genannt, ohne Test. Die übrigen Randformen dieses Slice trägt der Architect vor dem Code ein.

**Randform** (`AGENTS.md` §3.12; Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando`, §1): *Variable beim eigenen und beim anderen Kommando* — `--session-assignment` (nur `replay`) und `--record-empty-sessions` (nur `record`): beim eigenen Kommando liest der Leser die Variable, beim anderen Kommando (`record`/`play` bzw. `replay`/`play`) bleibt sie `PGR-E2001` (entschieden in der Spezifikation, `LH-FA-17.a`; Test: `TestLeserUmgebungNichtGelesenerOptionen` mit gekürztem `nichtGelesen()`, `TestLeserUmgebungAndererKommandos`).

- Ein Connection-Pool baut Verbindungen gleichzeitig auf; die Zuordnung ist dann nicht zugesichert — **Ausgang:** offen bis Closure.
- Die Zuordnung `connection` übersieht Sessions nur aus Lebendprüfungen, weil das Replay sie beim Laden aus `frei` nimmt; die n-te Verbindung erhielte dann eine falsche Session. Ein Test mit einer solchen Session an einer mittleren `id` fängt das — **Ausgang:** offen bis Closure.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-sessions-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 2. Schichten: zwei Schichten (Kern, PGWire-Adapter).

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-sqlite-format` (2026-10-11, `AGENTS.md` §3.13): Der Test mit `--format sqlite` liegt in DoD-Punkt 1 und in `test/integration`, ohne Code im Recording-Adapter; Liefer-Punkte weiter 2, Schichten weiter zwei (Kern, PGWire-Adapter).

Nachgezählt beim Eintragen der Sendungen aus `slice-v1-abschluss-zeitangaben-aufzeichnen` (Planner, 2026-10-12, `AGENTS.md` §3.13): Beide Sendungen gehören zu DoD-Punkt 1 (Test mit `--record-empty-sessions`, Fälle der Option) und fügen keinen Liefer-Punkt hinzu: Liefer-Punkte 2. Schichten: Die Optionen `--record-empty-sessions` und `--session-assignment` gibt es im CLI-Adapter erst mit diesem Slice (nur die Tests der Optionstabelle nennen sie); der CLI-Adapter ist damit eine Schicht, die §3 bisher nicht führte. Kern und CLI-Adapter sind zwei; der PGWire-Adapter ist die dritte, wenn er sich ändert. Die gleichzeitigen Verbindungen lieferte der Walking Skeleton (§1, *Bereits geliefert*), die Zeile des PGWire-Adapters in §3 trägt damit höchstens Tests; der Architect prüft das vor dem Code. Ändert er den Adapter, liegt der Plan über der Grenze, und der Slice wird vor dem Anspruch geschnitten (§4). Eine Sendung schiebt ihn nicht darüber: Der Fall in `TestDateiUngueltig` liegt im CLI-Adapter, den die Option ohnehin berührt.

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando` (Planner, 2026-10-12, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 1 (Leser der Option und Kürzung von `nichtGelesen()` im Test des CLI-Adapters) und fügt keinen Liefer-Punkt hinzu: Liefer-Punkte 2. Schichten: unverändert (Kern und CLI-Adapter; der PGWire-Adapter als dritte nur, wenn er sich ändert, §4); der CLI-Adapter war schon gezählt, `fremd_test.go` liegt in ihm.

Nachgezählt beim Eintragen der Abgrenzung zum Ende des Servers (Planner, 2026-10-12, `AGENTS.md` §3.13): Die Abgrenzung nimmt dem Slice nichts und fügt nichts hinzu; Liefer-Punkte weiter 2, Schichten weiter zwei (Kern und CLI-Adapter; der PGWire-Adapter als dritte nur, wenn er sich ändert, §4).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
