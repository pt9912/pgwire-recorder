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

**Bezug:** [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-FA-06`](../../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md)

**Berührte Spec-Stellen:** `LH-FA-18.a` · `SPEC-041` · `SPEC-002` · `ARC-006` · `ARC-007`

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

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Replay — `slice-extended-query-replay`.
- Domain-Typen, Formregeln (`Interaction.Validate`) sowie Schreiber und Leser des Formats `yaml` für `type: extended` — geliefert von `slice-extended-query-modell`; dieser Slice bildet die PGWire-Nachrichten auf diese Typen ab.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Ein Go-Client mit Prepared Statements läuft über den Recorder, und die Aufzeichnung enthält die Nachrichtenfolge (Integrationstest).
- [ ] Der Record-Service gruppiert die Nachrichten nach `LH-FA-18.a` (Flush- und Sync-Gruppen, Pipelining, Abbruch), und jede aufgezeichnete Interaktion besteht `Interaction.Validate` (Unit-Tests).
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
| `internal/adapters/driving/pgwire`, `…/driven/postgres` | update | Neue Nachrichten |
| `internal/hexagon/services` | update | Record-Service |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-extended-query-modell` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: der Slice verlangt mehr als drei Liefer-Punkte — zurück zur Zerlegung.
- `in-progress` → `open`: Der Upstream-Adapter braucht eine neue Port-Operation — Entscheidung klären.


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

- Asynchrone Nachrichtenfolgen (Pipelining) sind in der Aufzeichnung nicht eindeutig geordnet — **Ausgang:** offen bis Closure.

- `Interaction.Validate` prüft nicht, ob eine Client-Nachricht nur die Felder ihres Typs trägt; der YAML-Schreiber gibt nur diese aus, ein fremd belegtes Feld fiele beim Schreiben still weg. Das Mapping aus den PGWire-Nachrichten belegt daher nur die Felder des Typs (aus `slice-extended-query-modell`) — **Ausgang:** offen bis Closure.

- Der YAML-Schreiber prüft nicht mit `Interaction.Validate`: eine fehlerhaft gruppierte Interaktion (etwa zwei `sync` in einer Gruppe) wird geschrieben und fällt erst beim Laden im Replay als beschädigt auf. Der Record-Service prüft daher jede Interaktion mit `Validate`, bevor er sie übernimmt (aus `slice-extended-query-modell`, Review F-292) — **Ausgang:** offen bis Closure.

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

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
