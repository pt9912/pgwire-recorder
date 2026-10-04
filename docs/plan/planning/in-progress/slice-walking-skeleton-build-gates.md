# Slice slice-walking-skeleton-build-gates: Docker-only Build, Test und Architektur-Gate

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-walking-skeleton.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [ADR-0001](../../adr/0001-hexagonale-architektur.md), [ADR-0009](../../adr/0009-implementierungssprache-go.md), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md)

**Berührte Spec-Stellen:** `ARC-009` · Architektur-Gate (architecture.md §6)

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

**Ziel:** Ein Go-Modul mit der Package-Struktur des Hexagons existiert, und Build und Test laufen ausschließlich über `make` in Docker; das bereits verdrahtete Architektur-Gate (`make a-check`) prüft die Packages, und `make gates` führt alles mit.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Fachlogik, PGWire- und Recording-Code — gehören zu `slice-walking-skeleton-record` und `slice-walking-skeleton-replay`; hier entstehen nur leere Package-Gerüste, an denen das Gate seine Regeln prüft.
- Container-Image des Produkts — Gegenstand von welle-v1-abschluss.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] `make build` und `make test` bauen und testen das Go-Modul in Docker ohne Host-Toolchain ([`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit)).
- [ ] Das Architektur-Gate (`make a-check`, an `make gates` gehängt) prüft die angelegten Packages gegen `.a-check.yml`; eine absichtliche Verletzung (Import von `pgproto3` in `internal/hexagon/model`) lässt es fehlschlagen.
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
| `go.mod`, `cmd/pgwire-recorder/`, `internal/…` (leere Packages gemäß Package-Struktur) | neu | Gerüst für das Hexagon; Gate braucht Code-Pfade |
| `Dockerfile` bzw. Build-Image, `harness/mk/build.mk` | neu | Docker-only Build/Test (`AGENTS.md` §3.1) |
| `a-check.mk`, `harness/mk/arch-gate.mk` | vorhanden | Architektur-Gate gemäß `.a-check.yml`, Image per Digest gepinnt; keine Änderung nötig |
| `harness/README.md` §Sensors | update | Bindung des neuen Gates deklarieren |
| Test: Negativfall des Gates | neu | Happy/Negative — Gate grün bei sauberer Struktur, rot bei verbotenem Import |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die ADRs des Welle-Triggers sind `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: das Gate verlangt mehr als einen Slice (z. B. zusätzlich Lint und Coverage) — zurück zur Zerlegung.
- `in-progress` → `open`: das gepinnte a-check-Image prüft die Go-Packages nicht wie erwartet — als Carveout klären.

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

- Ob a-check die Go-Importe der Packages wie in `.a-check.yml` modelliert erkennt, ist erst mit Code belegbar — **Ausgang:** offen bis Closure.
- Die Gate-Konfiguration ist Go-spezifisch; Schwellen dürfen nicht ohne ADR gesenkt werden (`AGENTS.md` §3.6) — **Ausgang:** offen bis Closure.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert in `harness/conventions.md` noch keine Sub-Areas; berührt wird das gesamte Repo als eine Sub-Area. Die Schwelle ≥ 2 von 3 Achsen ist für eine Ein-Sub-Area-Planung nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
