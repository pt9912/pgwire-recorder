# Harness

---

## Purpose

Dieser Harness verbindet bestehende Spezifikationen, ADRs,
Planning-Dokumente und Gates. Er ist **kein Ersatz** für `spec/` oder
`docs/`, sondern ein **Einstiegspunkt** für Menschen und AI-Code-Agenten.

Wenn diese Datei einer kanonischen Quelle widerspricht, **gewinnt die
kanonische Quelle**, und diese Datei wird angepasst.

Strukturregeln (Verzeichniskonvention, ID-Schemata, Modus-Deklarationen
pro Sub-Area, Zusatzklassen für Sensors-Bindung) sowie Adaptionen ggü.
der adoptierten Baseline leben in [`conventions.md`](conventions.md).
Diese Datei dupliziert sie nicht.

## Source precedence

| Rang | Datei | Charakter |
|---|---|---|
| 1 | [`spec/lastenheft.md`](../spec/lastenheft.md) | vertraglich abnahmebindend |
| 2 | [`spec/spezifikation.md`](../spec/spezifikation.md) | technisch fortschreibbar |
| 3 | [`spec/architecture.md`](../spec/architecture.md) | Komponenten/Sequenzen, meilensteinfrei |
| 4 | [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| 5 | [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md) | Wellen-Sequenz |
| 6 | `docs/user/*` *(falls vorhanden)* | Operations, Quality, Releasing | <!-- d-check:ignore (Verzeichnis optional; entlinkt, da im frischen Repo selten vorhanden) -->
| 7 | [`README.md`](../README.md) | Projekt-Überblick |
| 8 | [`AGENTS.md`](../AGENTS.md) | Agent-Briefing |
| 9 | diese Datei | Harness-Einstieg |

> Die Ränge 1–3 sind die **drei Spec-Straten** — Vertrag, Technik, Sicht —,
> und alle drei sind obligatorisch (Baseline-Regelwerk
> `grundlagen-referenz-richtung.md` §Spec-Straten). **Adaption ist die
> Zwei-Straten-Form, nicht die Drei-Straten-Form**: Wer Rang 2 streicht,
> deklariert das als `MR-<NNN>` in [`conventions.md`](conventions.md) und
> nummeriert neu (dann acht Ränge).

## Guides (Feedforward-Quellen)



| Quelle | Inhalt |
|---|---|
| [`spec/lastenheft.md`](../spec/lastenheft.md) | Anforderungen, IDs, Akzeptanzkriterien |
| [`spec/spezifikation.md`](../spec/spezifikation.md) | technische Details, Defaults |
| [`spec/architecture.md`](../spec/architecture.md) | Komponenten, Schichten, Constraints |
| [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| [`docs/plan/planning/`](../docs/plan/planning/) | Slice-Pläne und Roadmap |
| [`AGENTS.md`](../AGENTS.md) | Hard Rules, Source Precedence, Workflow |
| [`conventions.md`](conventions.md) | repo-lokale Strukturregeln, Adaptions-Block (`MR-*`), Modus-Deklarationen |
| `.harness/skills/reviewer.md` | Reviewer-Skill: HIGH-Liste, Kategorien-Regeln, Negativbefund-Pflicht, Output-Schema (Modul 10) — nächste Rolle nach Schritt 8 des Minimal Agent Workflow, nicht Teil der Implementer-Eingabe |
| `.harness/baseline/<tag>/regelwerk/` (vendored; `README.md` = Index) | adoptiertes Betriebsregelwerk in Agenten-Kurzform — **präsente nachschlagbare Vertiefung**, pro Entscheidung abschnittsweise (siehe [`AGENTS.md`](../AGENTS.md) §1); derivativ, Stand/Tag siehe [`conventions.md`](conventions.md) §Baseline |
| `.harness/baseline/<tag>/templates/` (vendored, parallel) | Referenz-Form der Skelette, auf die das Regelwerk mit `../templates/…` als „Ziel-Form" verweist (netzlos, weil parallel zu `regelwerk/`); Vorlagen zum Kopieren-und-Ausfüllen |

## Sensors (Feedback-Gates)



| Target | Vertrag | Bindung |
|---|---|---|
| `make docs-check` | Doku-Referenzen (d-check); das Gate, das die Vorlage mitbringt | — |
| `make build` | baut das Binary über das Multistage-`Dockerfile` (Stufe `build`, netzlos außer der Download-Stufe) | [ADR-0009](../docs/plan/adr/0009-implementierungssprache-go.md), [ADR-0026](../docs/plan/adr/0026-build-und-test-im-multistage-dockerfile.md) |
| `make test` | `go vet` und Unit-Tests über das `Dockerfile` (Stufe `test`, netzlos außer der Download-Stufe) | [ADR-0026](../docs/plan/adr/0026-build-und-test-im-multistage-dockerfile.md) |
| `make test-integration` | Integrationstests des Binaries gegen ein gepinntes PostgreSQL-Image in einem eigenen Docker-Netz; schreibt `docs/user/e2e-abdeckung.md` | [ADR-0026](../docs/plan/adr/0026-build-und-test-im-multistage-dockerfile.md) |
| `make a-check` | Architektur-Regeln des Hexagons gemäß `.a-check.yml` (a-check, netzlos, schreibgeschützt): Import-Richtungen der Packages und Bibliotheken je Adapter | [ADR-0001](../docs/plan/adr/0001-hexagonale-architektur.md) |
| `make a-check-negativ` | Gegenprobe des Architektur-Gates in Kopien des Arbeitsbaums: `pgproto3` und `crypto/tls` werden im Domain Model und im Recording-Adapter abgelehnt, in beiden PGWire-Adaptern zugelassen | [ADR-0001](../docs/plan/adr/0001-hexagonale-architektur.md), [ADR-0010](../docs/plan/adr/0010-verwendung-von-pgproto3.md) |
| `make baseline-verify` | vendored Baseline gegen `SHA256SUMS` (Integrität und Vollständigkeit, netzlos) | — |
| `make hook-gegenprobe` | Gegenprobe des Commit-Trägers `.githooks/commit-msg`: Kennung eines vorhandenen Slice angenommen, erfundene und fehlende Kennung abgelehnt (bash, ohne Docker) | [ADR-0025](../docs/plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md) |
| `make gates` | alle inneren Gates | — |

**Werkzeuge — genannt, weil der Lauf sie braucht, aber kein Gate:**

| Target | Tut was | Bindung |
|---|---|---|
| `make slice-mv` | bewegt einen Slice zwischen den Lifecycle-Verzeichnissen und zieht die Verweise nach, prüft nichts | kein Gate |
| `make archive-welle` | archiviert die Zeitdokumente einer geschlossenen Welle | kein Gate |
| `make schema-validate` | prüft das neutrale Schema der SQLite-Aufzeichnung mit d-migrate (netzlos), prüft sonst nichts | kein Gate |
| `make schema-generate` | gibt das SQL für SQLite aus dem neutralen Schema auf stdout aus | kein Gate |
| `make hooks-install` | aktiviert den git-eigenen `commit-msg`-Träger im Klon | kein Gate |
| `make go-mod-tidy` | aktualisiert `go.mod` und `go.sum` mit Netz im gepinnten Go-Image, prüft nichts | kein Gate |
| `make a-check-graph` | gibt den Architektur-Graphen (Mermaid) aus `.a-check.yml` aus, prüft nichts | kein Gate |
| `make doc-trace` | gibt die Requirements Traceability Matrix aus (Anforderung, Entscheidungen, Slices), prüft nichts | kein Gate |

**Aktueller Lauf-Status:** CI-Badge bzw. lokal `make help` / `make gates`.
**Rote Gates:** Begründung in einem Carveout (Modul 7); bisher keiner.
**Nicht behauptet:** Lint und Testabdeckung (keine Schwelle).



## Traceability rules

- PRs/Commits **müssen** mindestens eine `LH-*`-, `ADR-*`-, `MR-*`- oder `slice-*`-Kennung nennen; eine benannte Slice-Kennung zählt, wenn sie als eigenes Wort in der Message (ohne Kommentarzeilen, vor dem Diff von `git commit -v`) steht und der Slice als Datei unter `docs/plan/planning/<lifecycle>/` im Index liegt ([ADR-0025](../docs/plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md)).
- Neue oder geänderte Anforderungen brauchen einen Beleg: Test, Gate, Demo oder ADR.
- Neue ADRs müssen im ADR-Index ergänzt werden.
- Änderungen an Planning-Dokumenten müssen die Lifecycle-Regeln beachten (open → next → in-progress → done; reine `git mv`-Commits siehe AGENTS.md §3.3).

## Safety and scope boundaries



- Docker-only: kein lokaler Toolchain-Install, alles läuft über `make` (`AGENTS.md` §3.1).
- Recordings können sensible Daten enthalten; das Produkt maskiert nichts (`LH-RB-01`).

## Minimal agent workflow

1. Diese Datei lesen.
2. Relevante kanonische Quelle lesen.
3. Betroffene IDs identifizieren.
4. Kleinste Änderung planen.
5. Engsten nützlichen Sensor laufen lassen.
6. Repo-weiten Gate-Lauf vor Handoff (`make gates`).
7. Doku/Indizes aktualisieren, falls ein öffentlicher Vertrag berührt.
8. Ausgeführte Sensors und verbleibende Risiken berichten.

Dieser Workflow deckt ausschließlich die Implementer-Rolle ab. Schritt 8
ist der Rollenwechsel, kein Abschluss: Bericht → Handoff an Reviewer
(`.harness/skills/reviewer.md`, siehe §Guides) → Verifier. Kein
Self-Review — anderer Kontext findet andere Findings, derselbe Kontext
dieselben blinden Flecken (Baseline-Regelwerk `modul-08-agentenrollen.md`).

## Leseordnung

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-harness-dateien.md`
§harness/README.md als Einstiegspunkt — die Menschen-Hälfte des Einstiegs:
drei bis fünf **geordnete** Zeiger, was ein neuer Mensch zuerst liest und was
bei Bedarf; eine Leseordnung, die alles nennt, ist keine.

1. `AGENTS.md` §Harte Regeln
2. `spec/lastenheft.md`, danach `spec/spezifikation.md` und `spec/architecture.md`
3. bei Bedarf `harness/conventions.md`
