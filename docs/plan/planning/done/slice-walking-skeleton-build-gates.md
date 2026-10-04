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
- `test/integration/` aus der Package-Struktur — entsteht mit dem ersten Integrationstest in `slice-walking-skeleton-record`; ein leeres Verzeichnis trägt git nicht.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `make build` und `make test` bauen und testen das Go-Modul in Docker ohne Host-Toolchain ([`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit)).
- [x] Das Architektur-Gate (`make a-check`, an `make gates` gehängt) prüft die angelegten Packages gegen `.a-check.yml`; eine absichtliche Verletzung (Import von `pgproto3` in `internal/hexagon/model`) lässt es fehlschlagen.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `go.mod`, `cmd/pgwire-recorder/`, `internal/…` (leere Packages gemäß Package-Struktur) | neu | Gerüst für das Hexagon; Gate braucht Code-Pfade |
| `harness/mk/build.mk` (gepinntes `golang`-Image, kein eigenes `Dockerfile`) | neu | Docker-only Build/Test (`AGENTS.md` §3.1) |
| `a-check.mk`, `harness/mk/arch-gate.mk` | vorhanden | Architektur-Gate gemäß `.a-check.yml`, Image per Digest gepinnt; keine Änderung nötig |
| `.a-check.yml` (`tech`) | update | Eine Regel je Bibliothek mit Liste beider PGWire-Adapter; zwei Einträge mit demselben Muster wertet a-check nur einmal aus |
| `harness/README.md` §Sensors | update | Bindung des neuen Gates deklarieren |
| `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk` | neu | Gegenprobe als Gate in fünf Fällen: rot, wenn a-check `pgproto3` oder `crypto/tls` außerhalb der PGWire-Adapter durchlässt oder in einem der beiden ablehnt |

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

- Ob a-check die Go-Importe der Packages wie in `.a-check.yml` modelliert erkennt, ist erst mit Code belegbar — **Ausgang:** entfallen: belegt durch `make a-check-negativ` und Mutationsproben des Reviewers und des Verifiers; der dabei gefundene Fehlschnitt der `tech`-Regel (F-204) ist im Slice behoben.
- Die Gate-Konfiguration ist Go-spezifisch; Schwellen dürfen nicht ohne ADR gesenkt werden (`AGENTS.md` §3.6) — **Ausgang:** entfallen: keine Schwelle gesenkt; die Zulassung des Upstream-Adapters deckt [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md).
- Der netzlose Build trägt nur ohne externe Abhängigkeiten; mit `pgproto3` braucht der nächste Slice eine netzlose Quelle der Module (zum Beispiel ein vendored Verzeichnis) — **Ausgang:** eingetreten: übernimmt `slice-walking-skeleton-record` (Risiko in dessen §6).
- Das Architektur-Gate erkennt einen Zugriff des Domain Models auf das Dateisystem (`os`) nicht — **Ausgang:** weiter offen: → `BEO-REPO/kern-fremdimporte-nur-ueber-tech-liste` im Register.

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

- **Was hat funktioniert:** Build und Test laufen netzlos und schreibgeschützt im per Digest gepinnten Image; die Package-Struktur entspricht `spec/architecture.md` §2.1. Rollen-Trennung hat gewirkt: der Reviewer fand den Fehlschnitt der `tech`-Regel, den ein grünes Gate verdeckte.
- **Was ging anders als geplant:** Plan §3 sah `a-check.mk` und `.a-check.yml` als unverändert; die `tech`-Regel musste als Liste beider PGWire-Adapter neu gefasst werden (F-204). Die Lifecycle-Moves liefen von Hand, bis [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) den Commit-Hook für benannte Slices öffnete. Die Gegenprobe wuchs von einem auf fünf Fälle.
- **Steering-Loop-Eintrag:** Sensor ergänzt: Gegenprobe des Architektur-Gates mit erlaubten und verbotenen Fällen je Bibliothek — liegt in `harness/mk/arch-negativ.mk`. Sensor ergänzt: Gegenprobe des Commit-Trägers — liegt in `harness/mk/hook-gegenprobe.mk`. Auslöser: Review-Finding F-204 und `BEO-REPO/werkzeug-commit-ohne-zugelassene-kennung` (je 1×).
- **Beobachtungs-Register (`../observations/`):** `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen/` neu angelegt, Beleg `evidence/slice-walking-skeleton-build-gates.md` — Zähler 1×; `BEO-REPO/kern-fremdimporte-nur-ueber-tech-liste/` neu angelegt, Zähler 1×; `BEO-REPO/werkzeug-commit-ohne-zugelassene-kennung/` neu angelegt und gestrichen (Ursache durch [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) entfallen).
- **Folge-Slices:** `slice-walking-skeleton-record` (Walking Skeleton: Record-Modus) — ist eine Datei in `open/`; trägt das Risiko der netzlosen Abhängigkeiten.
- **Risiken aus §6:** jedes mit genau einem Ausgang — siehe §6.
- **Drei Paarungen:** Repo mit Wellen-Betrieb — geprüft von der Closure von `welle-walking-skeleton`.
- **Belege:** Review `docs/reviews/2026-10-04-review-slice-walking-skeleton-build-gates.md`, Folge-Review `docs/reviews/2026-10-04-review-slice-walking-skeleton-build-gates-folge.md`, Verifikation `docs/reviews/2026-10-04-verifikation-slice-walking-skeleton-build-gates.md`.

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
