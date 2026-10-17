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
| 2 | [`spec/spezifikation.md`](../spec/spezifikation.md) (Produkt) und [`harness/vertraege/`](vertraege/README.md) (Werkzeuge) | technisch fortschreibbar |
| 3 | [`spec/architecture.md`](../spec/architecture.md) | Komponenten/Sequenzen, meilensteinfrei |
| 4 | [`docs/plan/adr/`](../docs/plan/adr/) | Architekturentscheidungen |
| 5 | [`docs/plan/planning/in-progress/roadmap.md`](../docs/plan/planning/in-progress/roadmap.md) | Wellen-Sequenz |
| 6 | `docs/user/*` *(falls vorhanden)* | Operations, Quality, Releasing | <!-- d-check:ignore (Verzeichnis optional; entlinkt, da im frischen Repo selten vorhanden) -->
| 7 | [`README.md`](../README.md) | Projekt-Überblick |
| 8 | [`AGENTS.md`](../AGENTS.md) | Agent-Briefing |
| 9 | diese Datei | Harness-Einstieg |

> Die Ränge 1–3 sind die **drei Spec-Straten** — Vertrag, Technik, Sicht —,
> und alle drei sind obligatorisch (Baseline-Regelwerk
> `grundlagen-referenz-richtung.md` §Spec-Straten). Rang 2 führt den Ordner `harness/vertraege/`
> außerhalb der Straten mit ([`MR-002`](conventions.md#mr-002)). **Adaption ist die
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
| `make build` | baut Binary und Produkt-Image über das Multistage-`Dockerfile` (Stufe `runtime`, netzlos außer der Download-Stufe) | [ADR-0009](../docs/plan/adr/0009-implementierungssprache-go.md), [ADR-0026](../docs/plan/adr/0026-build-und-test-im-multistage-dockerfile.md) |
| `make test` | Formatierung (`gofmt -l` leer, sonst die Dateien), `go vet` und Unit-Tests über das `Dockerfile` (Stufe `test`, netzlos außer der Download-Stufe) | [ADR-0026](../docs/plan/adr/0026-build-und-test-im-multistage-dockerfile.md) |
| `make test-integration` | Integrationstests des Binaries gegen ein gepinntes PostgreSQL-Image in einem eigenen Docker-Netz; schreibt nichts in den Arbeitsbaum | [ADR-0026](../docs/plan/adr/0026-build-und-test-im-multistage-dockerfile.md) |
| `make test-postgres-versionen` | die Teilmenge der Integrationstests, deren Aussage vom Server abhängt (Abnahmeszenarien 1, 2 und 7, Fall SCRAM-SHA-256 der Anmeldung, Abgleich der Lebendprüfung), gegen jede Version der Liste in `harness/mk/postgres-versionen.mk` außer der Referenz, je Image per Digest gepinnt, nacheinander in einem eigenen Docker-Netz; die Referenz fährt `make test-integration` mit der vollen Suite; Vertrag [`SPEC-053`](vertraege/SPEC-053-postgres-versionen.md) | [ADR-0042](../docs/plan/adr/0042-versionsmatrix-der-integrationstests.md) · seit slice-v1-abschluss-postgres-versionen-matrix |
| `make test-postgres-versionen-gegenprobe` | Gegenprobe der Versionsmatrix in Kopien des Arbeitsbaums mit einer Attrappe von `docker` (über `DOCKER`): die unveränderte Kopie grün, je Punkt 1 bis 9 und 12 von `SPEC-053` Fälle und Mutationen, die nur die Fälle ihrer Gruppe rot färben, darunter Eingabeformen, Teilmenge je Test, Server, Ablauf, Aufräumen nach Ende, Fehler, `SIGINT` und `SIGTERM`, Aufrufarten von `make`; hängt an `GATE_CHECKS` und steht in der Ausnahmeliste der schnellen Kette (bash, make, setsid, ohne Docker) | [ADR-0042](../docs/plan/adr/0042-versionsmatrix-der-integrationstests.md) · seit slice-v1-abschluss-postgres-versionen-matrix |
| `make test-container` | zwei Builds der sechs Binaries (`linux`, `darwin`, `windows`, je `amd64` und `arm64`) in zwei Verzeichnissen außerhalb des Arbeitsbaums, der zweite ohne Cache der Stufen `build` und `dist`, mit gleichen Summen; Kopf jedes Binaries (ELF, Mach-O, PE mit Architektur); Image je Linux-Plattform mit demselben Binary; Smoke auf dem Image der Plattform des Docker-Servers in einem `--internal`-Netz (`version`, `record` gegen PostgreSQL, `replay` ohne PostgreSQL, gleiche Anfrage, Exit 0 nach `SIGTERM`); Vertrag [`SPEC-055`](vertraege/SPEC-055-container-und-binaries.md) | [ADR-0043](../docs/plan/adr/0043-binaries-und-container-gate.md), [ADR-0026](../docs/plan/adr/0026-build-und-test-im-multistage-dockerfile.md) · seit slice-v1-abschluss-container |
| `make test-container-gegenprobe` | Gegenprobe von `binaries` und `test-container` in Kopien des Arbeitsbaums mit einer Attrappe von `docker` (über `DOCKER`): die unveränderte Kopie grün, je Punkt 1 bis 9, 11 und 12 von `SPEC-055` Fälle und Mutationen, die nur die Fälle ihrer Gruppe rot färben (dazu, dass ein Lauf keine neue Datei im Arbeitsbaum anlegt, Punkt 10), darunter Plattformlisten, `VERSION`, Summen, Kopf je Format, Image, Smoke, Aufräumen nach Ende, `SIGINT` und `SIGTERM`, Aufrufarten von `make`; hängt an `GATE_CHECKS` und steht in der Ausnahmeliste der schnellen Kette (bash, make, setsid, ohne Docker) | [ADR-0043](../docs/plan/adr/0043-binaries-und-container-gate.md) · seit slice-v1-abschluss-container |
| `make a-check` | Architektur-Regeln des Hexagons gemäß `.a-check.yml` (a-check, netzlos, schreibgeschützt): Import-Richtungen der Packages und Bibliotheken je Adapter | [ADR-0001](../docs/plan/adr/0001-hexagonale-architektur.md) |
| `make a-check-negativ` | Gegenprobe des Architektur-Gates in Kopien des Arbeitsbaums: `pgproto3` und `crypto/tls` werden im Domain Model und im Recording-Adapter abgelehnt, in beiden PGWire-Adaptern zugelassen; jeder der beiden YAML-Modulpfade wird im Domain Model, im Postgres-Adapter und im PGWire-Adapter abgelehnt, im CLI-Adapter zugelassen; `modernc.org/sqlite` wird im Domain Model, im Postgres-Adapter, im PGWire-Adapter und im CLI-Adapter abgelehnt, im Recording-Adapter zugelassen; `time` wird im Domain Model, in den Driven Ports und in den Services abgelehnt, im PGWire-Adapter zugelassen und im Composition Root nicht abgelehnt (von der Regel ausgenommen), `runtime` im Domain Model zugelassen | [ADR-0001](../docs/plan/adr/0001-hexagonale-architektur.md), [ADR-0010](../docs/plan/adr/0010-verwendung-von-pgproto3.md), [ADR-0024](../docs/plan/adr/0024-sqlite-bibliothek.md), [ADR-0027](../docs/plan/adr/0027-yaml-bibliothek.md), [ADR-0036](../docs/plan/adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), [ADR-0015](../docs/plan/adr/0015-uhr-port.md) |
| `make baseline-verify` | vendored Baseline gegen `SHA256SUMS` (Integrität und Vollständigkeit, netzlos) | — |
| `make hook-gegenprobe` | Gegenprobe des Commit-Trägers `.githooks/commit-msg`: Kennung eines vorhandenen Slice oder einer vorhandenen Welle angenommen, erfundene Kennung, Name einer Closure-Notiz und fehlende Kennung abgelehnt (bash, ohne Docker) | [ADR-0025](../docs/plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md), [ADR-0029](../docs/plan/adr/0029-benannte-welle-kennungen-im-commit-hook.md) |
| [`make abdeckung-check`](sensors/abdeckung-check.md) | die Abdeckungstabellen `docs/user/abdeckung-*.md` entsprechen den Abdeckungs-Deklarationen der Tests (Anforderung und Pfad), jeder E2E-Test trägt eine; Vertrag [`SPEC-048`](vertraege/SPEC-048-abdeckung.md#spec-048--abdeckung-je-anforderung-und-pfad-abdeckung) (bash, ohne Docker) | [ADR-0028](../docs/plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) |
| `make abdeckung-gegenprobe` | Gegenprobe des Abdeckungs-Skripts in Temp-Bäumen: Fehlformen abgelehnt, Inhalt der vier Tabellen geprüft, `--check` schreibt nicht (bash, ohne Docker) | [ADR-0028](../docs/plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) |
| [`make kopf-check`](sensors/kopf-check.md) | jeder lebende Slice-Plan führt im Kopf jede Lastenheft-, Spezifikations- und Sicht-Kennung, die §1 oder §2 nennt; Vertrag [`SPEC-047`](vertraege/SPEC-047-kopf-der-plaene.md#spec-047--kopf-der-pläne-kopf-check) (bash, ohne Docker) | [ADR-0032](../docs/plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) · seit slice-harness-kopf-sensor |
| `make kopf-check-gegenprobe` | Gegenprobe des Kopf-Skripts in Temp-Bäumen: je Nummer der ADR Fälle, die abgelehnt bzw. angenommen werden müssen, die akzeptierten Negative angenommen; dazu in einem Temp-Baum mit Makefile und Fragment, dass beide Ziele an `GATE_CHECKS` hängen und über `make` bei rotem Skript mit Exit ungleich 0 enden (bash, ohne Docker) | [ADR-0032](../docs/plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) |
| [`make meldungskatalog-check`](sensors/meldungskatalog-check.md) | die Meldungscodes `PGR-…` der Konstanten `Code…` in `internal/hexagon/model/fehler.go`, der Tabelle in `SPEC-034` und der Tabellen `### Fehlercodes` und `### Warnungen` des Benutzerhandbuchs stimmen überein (Rückfälle `…000` stehen nur im Katalog), jeder Code hat die Form, kein Literal `"PGR-` steht in einer anderen Nicht-Test-Datei unter `cmd/` und `internal/`; Vertrag [`SPEC-054`](vertraege/SPEC-054-meldungskatalog.md#spec-054--abgleich-von-code-tabelle-und-meldungskatalog-meldungskatalog-check); fährt in der schnellen Kette (bash und awk, ohne Docker) | [ADR-0011](../docs/plan/adr/0011-meldungscodes-praefix-pgr.md) · seit slice-harness-meldungskatalog-gate |
| `make meldungskatalog-check-gegenprobe` | Gegenprobe des Abgleichs in Temp-Bäumen: je Punkt 1 bis 6 und 8 bis 10 von `SPEC-054` und je Element seiner Aufzählungen Fälle, die abgelehnt bzw. angenommen werden müssen, darunter die Reihenfolge der Sperren; dazu in einem Temp-Baum mit Makefile und Fragmenten, dass beide Ziele an `GATE_CHECKS` hängen, `make` bei rotem Skript mit Exit ungleich 0 endet, die Prüfung nicht und die Gegenprobe in `GATE_SCHNELL_AUSNAHME` steht und `-n`, `-q`, `-t`, `-i` nichts hinterlassen; steht in der Ausnahmeliste der schnellen Kette (bash, make, ohne Docker) | [ADR-0011](../docs/plan/adr/0011-meldungscodes-praefix-pgr.md) |
| [`make lint`](sensors/lint.md) | der Go-Code unter `cmd/`, `internal/` und `test/` ist ohne Befund nach dem Lint-Profil `.golangci.yml` und den eigenen Prüfungen, trägt kein `//nolint`, und jede Ausnahme im Profil steht mit `# Why:` in fester Form und blendet etwas aus; Vertrag [`SPEC-049`](vertraege/SPEC-049-lint-profil.md#spec-049--lint-profil-lint) (Stufe `lint` des `Dockerfile`; dass sie ohne Netz läuft, prüft die Gegenprobe nicht) | [ADR-0034](../docs/plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) · seit slice-harness-lint |
| `make lint-gegenprobe` | Gegenprobe des Lint-Gates in Kopien des Arbeitsbaums: die unveränderte Kopie grün, je Punkt von `SPEC-049` Mutanten, die rot werden und ihre eigene Zeile nennen, die freien Fälle ohne Befund; dazu in einer Kopie mit nur dem Fragment, dass beide Ziele an `GATE_CHECKS` hängen und über `make` bei rotem Lauf mit Exit ungleich 0 enden (bash und docker) | [ADR-0034](../docs/plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) |
| `make gates-schnell` | die schnelle Kette: `GATE_CHECKS` ohne die Ausnahmeliste im Fragment `harness/mk/gates-schnell.mk` (jedes Gate mit `-gegenprobe` im Namen, `a-check-negativ`, `test-integration`, `test-postgres-versionen`); schreibt bei grünem Lauf den Stand `.harness/state/gates-schnell-passed.diffsha`, bei rotem oder abgebrochenem keinen (ebenso bei `-i` in `MAKEFLAGS`; bei `-n`, `-q` und `-t` bleibt ein vorhandener Stand unberührt), und ersetzt den Nachweis des vollen Laufs nie; `make gates` fährt es nicht; Vertrag [`SPEC-052`](vertraege/SPEC-052-gate-kette.md#spec-052--gate-kette-und-wahl-der-kette-gates-schnell) Teil 1 (bash, make, git, Docker für die Gates selbst) | [ADR-0038](../docs/plan/adr/0038-gate-kette-voll-und-schnell.md) · seit slice-harness-gate-kette-schnell |
| `make gates-schnell-gegenprobe` | Gegenprobe der schnellen Kette in Kopien des Arbeitsbaums mit Attrappen als Gates: die unveränderte Kopie grün, je Punkt 1 bis 7 von `SPEC-052` Mutanten, die rot werden; die Ausnahmeliste deckt sich mit `GATE_CHECKS` auch im Arbeitsbaum selbst (bash, setsid, ohne Docker); hängt an `GATE_CHECKS` und steht selbst in der Ausnahmeliste | [ADR-0038](../docs/plan/adr/0038-gate-kette-voll-und-schnell.md) · seit slice-harness-gate-kette-schnell |
| `make stop-hook-gegenprobe` | Gegenprobe des Stop-Hooks `.claude/hooks/stop-require-gates.sh` und der Dateiliste des vollen Laufs (`.harness/state/gates-passed.dateiliste`, Funktionen in `tools/harness/dateiliste.sh`) in Kopien eines Temp-Repos: die unveränderte Kopie grün, je Fall (Punkt 1, 9 bis 22 von `SPEC-052` im Namen) Mutanten, die rot werden: der Hook gibt bei identischer Liste frei oder wenn jeder Unterschied auf der Positivliste liegt (`docs/plan/planning/`, `docs/reviews/`, `docs/user/`, `docs/developer/`, Wurzel-`README.md`; `harness/README.md` nicht; Symlinks mit Ziel und Spiegelregel) und der Stand der schnellen Kette der Hash des Arbeitsbaums ist, ein Merge-Konflikt, ein fehlender, nicht lesbarer oder beschädigter Stand (auch bei sauberem Baum), ein scheiternder Schritt und ein Abbruch ergeben die Antwort `block`, der Schleifenschutz bleibt, `make gates` schreibt die Liste nur bei grünem Lauf und ohne `-i` in `MAKEFLAGS` (ein älterer Nachweis entfällt dann), `-n`, `-q`, `-t` lassen ihn stehen, der Hook urteilt über das Repo seines Arbeitsverzeichnisses auch aus einem Unterverzeichnis und in leerer Umgebung; hängt an `GATE_CHECKS` und steht in der Ausnahmeliste der schnellen Kette; die Umgebung von Claude Code prüft sie nicht (bash, git, make, ohne Docker) | [ADR-0038](../docs/plan/adr/0038-gate-kette-voll-und-schnell.md) · seit slice-harness-gate-kette-stop-hook-liste |
| `make gates` | alle inneren Gates | — |

**Werkzeuge — genannt, weil der Lauf sie braucht, aber kein Gate:**

| Target | Tut was | Bindung |
|---|---|---|
| `make binaries` | baut die sechs Binaries `pgwire-recorder_<VERSION>_<os>_<arch>` (unter Windows `.exe`) mit `SHA256SUMS` nach `DIST` (Vorgabe `dist`, muss fehlen oder leer sein; steht in `.gitignore`); Vertrag [`SPEC-055`](vertraege/SPEC-055-container-und-binaries.md) Punkt 3 | kein Gate |
| `make slice-mv` | bewegt einen Slice zwischen den Lifecycle-Verzeichnissen und zieht die Verweise nach, prüft nichts | kein Gate |
| `make archive-welle` | archiviert die Zeitdokumente einer geschlossenen Welle | kein Gate |
| `make schema-validate` | prüft das neutrale Schema der SQLite-Aufzeichnung mit d-migrate (netzlos), prüft sonst nichts | kein Gate |
| `make schema-generate` | gibt das SQL für SQLite aus dem neutralen Schema auf stdout aus | kein Gate |
| `make hooks-install` | aktiviert den git-eigenen `commit-msg`-Träger im Klon | kein Gate |
| `make abdeckung` | schreibt die Abdeckungstabellen aus den Test-Deklarationen, prüft nichts | kein Gate |
| `make go-mod-tidy` | aktualisiert `go.mod` und `go.sum` mit Netz im gepinnten Go-Image, prüft nichts | kein Gate |
| `make a-check-graph` | gibt den Architektur-Graphen (Mermaid) aus `.a-check.yml` aus, prüft nichts | kein Gate |
| `make doc-trace` | gibt die Requirements Traceability Matrix aus (Anforderung, Entscheidungen, Slices), prüft nichts; als Slice zählt nur eine Datei flach in `done/`, ein archivierter Slice unter `done/<welle-id>/` nicht | kein Gate |

**Aktueller Lauf-Status:** CI-Badge bzw. lokal `make help` / `make gates`.
**Rote Gates:** Begründung in einem Carveout (Modul 7); bisher keiner.
**Nicht behauptet:** Testabdeckung (keine Schwelle).



## Traceability rules

- PRs/Commits **müssen** mindestens eine `LH-*`-, `ADR-*`-, `MR-*`- oder `slice-*`-Kennung nennen; eine benannte Slice-Kennung zählt, wenn sie als eigenes Wort in der Message (ohne Kommentarzeilen, vor dem Diff von `git commit -v`) steht und der Slice als Datei unter `docs/plan/planning/<lifecycle>/` im Index liegt ([ADR-0025](../docs/plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md)); ebenso eine Welle-Kennung, deren Welle-Datei unter `docs/plan/planning/` im Index liegt, flach, in `done/` oder im Archiv-Ordner darunter, nicht aber der Name einer Closure-Notiz ([ADR-0029](../docs/plan/adr/0029-benannte-welle-kennungen-im-commit-hook.md)).
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
