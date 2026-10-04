# Verifikation: slice-walking-skeleton-build-gates — 2026-10-04

**Rolle:** Verifier (Modul 11). Geprüft wird, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator; Diff gegen Entscheidungen und Hard Rules prüft der Reviewer.

**Gegenstand:** Slice-Plan `docs/plan/planning/done/slice-walking-skeleton-build-gates.md` (§1 bis §3, §6) gegen `git diff 759fc4f HEAD` bei HEAD `a834db1` (Commits `1703db5`, `554bd0b`, `d04e43c`, `8143a47`, `a834db1`).

**Eingang:** DoD-Bestätigung und Sensor-Angaben des Implementers; Review-Report `docs/reviews/2026-10-04-review-slice-walking-skeleton-build-gates.md` (Stand `8143a47`).

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-04

Alle Sensoren unten habe ich in diesem Kontext selbst gefahren. Behauptungen des Implementers dienen nicht als Beleg. Mutationen liefen nur in Kopien des Arbeitsbaums im Scratchpad. `git status --porcelain` im Repo war nach jedem Lauf leer.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `make build` und `make test` in Docker ohne Host-Toolchain. **Bestätigt.**

| Kommando | Ausgabe (gekürzt) | Exit |
|---|---|---|
| `which go` (Host) | keine Ausgabe, kein Go auf dem Host | — |
| `make build` | `docker run --rm --network none -u "$(id -u):$(id -g)" -v ".../pgwire-recorder":/src:ro … golang:1.27-alpine@sha256:8a5910f3… go build -trimpath -buildvcs=false ./...` | 0 |
| `make test` | derselbe Container mit `go test …`; zehn Packages, jeweils `[no test files]` | 0 |

Der Container läuft netzlos (`--network none`), der Quellbaum ist schreibgeschützt eingehängt (`:ro`) und das Image per Digest gepinnt. Nach dem Lauf ist der Arbeitsbaum sauber. Grenze: `make test` belegt, dass die Test-Kette läuft. Weil es noch keine Tests gibt, prüft er inhaltlich nichts. Die Sensors-Tabelle sagt das ehrlich („Nicht behauptet: Lint und Testabdeckung").

### Punkt 2: Das Architektur-Gate prüft die Packages gegen `.a-check.yml` und schlägt bei `pgproto3` in `internal/hexagon/model` fehl. **Bestätigt.**

Verdrahtung: `harness/mk/arch-gate.mk` bindet `a-check.mk` ein und hängt `a-check` an `GATE_CHECKS`. `harness/mk/arch-negativ.mk` hängt `a-check-negativ` an.

Läufe im Repo:

| Kommando | Ausgabe | Exit |
|---|---|---|
| `make a-check` | `gesamt: 0 Befund(e)` | 0 |
| `make a-check-negativ` | `a-check-negativ: gruen — Domain Model und Recording-Adapter abgelehnt, Upstream-Adapter zugelassen` | 0 |

Eigene Mutationsprobe (Kopie per `tar --exclude=./.git` im Scratchpad, dort `make -C <kopie> …`):

| # | Mutation | Kommando | Ergebnis | Erwartung |
|---|---|---|---|---|
| M1 | `import _ "github.com/jackc/pgx/v5/pgproto3"` in `internal/hexagon/model/mut.go` | `make a-check` | `make: *** [a-check.mk:16: a-check] Fehler 1`, Exit 2 | rot ✔ |
| M1b | wie M1 | `make gates` | `internal/hexagon/model/mut.go:3: core-impurity: Kern importiert github.com/jackc/pgx/v5/pgproto3`, `gesamt: 1 Befund(e)`, Exit 2 | `gates` rot ✔ |
| M2 | `pgproto3` und `crypto/tls` in `internal/adapters/driven/postgres` | `make a-check` | Exit 0 | grün ✔ ([ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md)) |
| M3 | `pgproto3` in `internal/adapters/driving/pgwire` | `make a-check` | Exit 0 | grün ✔ |
| M4 | `pgproto3` in `internal/adapters/driving/cli` | `make a-check` | `tech-leak: Tech github.com/jackc/pgx/v5/pgproto3 außerhalb internal/adapters/driving/pgwire\|internal/adapters/driven/postgres`, Exit 2 | rot ✔ |
| M5 | Kante `services` → `internal/adapters/driven/recording` | `make a-check` | `app-impurity: Application importiert …/internal/adapters/driven/recording`, Exit 2 | rot ✔ |
| M6 | `pgproto3`-Zeile aus `tech` in `.a-check.yml` entfernt | `make a-check-negativ` | `ROT in Fall 'model' — erwartet rot, a-check war gruen`, ebenso Fall `recording`, Exit 2 | Gegenprobe erkennt die gelockerte Konfiguration ✔ |
| M6b | wie M6, zusätzlich `pgproto3` im Model | `make a-check` | Exit 0 | zeigt eine Grenze, siehe unten |

Der DoD-Fall (M1) färbt sowohl `make a-check` als auch `make gates` rot. Der Punkt trägt.

**Grenze (M6b):** a-check meldet einen Fremdimport im Domain Model nur dann als core-impurity, wenn die Bibliothek in der `tech`-Liste steht. Eine Drittbibliothek, die dort fehlt, oder `os` (Review F-212) lässt das Gate im Kern zu. Die Gegenprobe `make a-check-negativ` hält die `pgproto3`-Zeile fest (M6). Für eine Bibliothek, die nicht in der Liste steht, gibt es keinen Wächter. Die DoD verlangt nur den `pgproto3`-Fall, deshalb ist das keine DoD-Verletzung. Es ist aber eine Reichweitengrenze von [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) („Infrastrukturabhängigkeiten liegen außerhalb des Application Core"), die der Planner kennen sollte. Für `os` ist sie im Plan §6 schon als Risiko vermerkt.

### Punkt 3: `make gates` grün. **Bestätigt.**

| Kommando | Ausgabe (Schlusszeilen) | Exit |
|---|---|---|
| `make gates` | `gesamt: 0 Befund(e)` · `a-check-negativ: gruen — …` · `baseline-verify: v6.13.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` · `d-check: 88 Datei(en) geprüft, 0 Befund(e)`; davor `build` und `test` ohne Fehler | 0 |

Gelaufen bei HEAD `a834db1` mit sauberem Arbeitsbaum und vor dem Anlegen dieser Datei. Danach lief `make docs-check` noch einmal, damit diese Datei mitgeprüft ist (siehe §5).

### Prozess-Punkte der DoD (nicht Teil des Auftrags, nur vermerkt)

- **Review:** Der Report liegt vor, prüft aber den Stand `8143a47`. Der Commit `a834db1` (`.a-check.yml` `tech`, neu geschriebene Gegenprobe, Sensors-Bindung, Plan §3/§6) war nicht Gegenstand eines Reviews. Inhaltlich decken meine Mutationen M2 bis M6 die Behebung von F-204 ab. F-206, F-208 und F-210 sind im Diff sichtbar adressiert, F-207 teilweise: Das Muster ist jetzt der Dateipfad bzw. `tech-leak` statt nur `pgproto3`. Ein Folge-Review ist trotzdem nicht ersetzt. Das entscheidet der Planner.
- Closure-Notiz, Register, Risiko-Ausgänge und Paarungen: noch offen (§7 leer). Das ist für `in-progress` erwartbar.

## 2. Plan §3 gegen Code-Diff

| Plan §3 | Diff | Befund |
|---|---|---|
| `go.mod`, `cmd/pgwire-recorder/`, `internal/…` — neu | `go.mod` (`module github.com/pt9912/pgwire-recorder`, `go 1.27`, kein `require`), `cmd/pgwire-recorder/main.go`, neun `doc.go` | konform |
| `harness/mk/build.mk` (gepinntes `golang`-Image, kein `Dockerfile`) — neu | `harness/mk/build.mk`, kein `Dockerfile` | konform |
| `a-check.mk`, `harness/mk/arch-gate.mk` — vorhanden, keine Änderung | beide nicht im Diff | konform |
| `.a-check.yml` (`tech`) — update | zwei Doppel-Einträge zu je einem Eintrag mit Adapter-Liste zusammengeführt | konform; Wirkung durch M2 bis M4 belegt |
| `harness/README.md` §Sensors — update | vier Zeilen neu bzw. geändert (`build`, `test`, `a-check`, `a-check-negativ`), „Nicht behauptet" angepasst | konform |
| `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk` — neu | beide neu, an `GATE_CHECKS` gehängt | konform |
| — | `README.md` (Projektstand, Gate-Aufzählung) | **nicht im Plan.** Folgeänderung nach Schritt 7 des Minimal Agent Workflow (öffentliche Beschreibung der Gates). Kleine Abweichung ohne Konformitätsbruch; die Aufzählung deckt sich mit `GATE_CHECKS`. |
| — | Slice-Plan selbst (§1, §3, §6), Review-Report | Planungs- und Prozessartefakte, keine Code-Abweichung |

Plan §6 nennt zwei neue Risiken (netzloser Build mit Abhängigkeiten, `os` im Model). Beide habe ich bestätigt, siehe Punkt 2 (Grenze) und `build.mk` (`--network none`, `GOFLAGS=-mod=readonly`). Ihr Ausgang wird erst bei der Closure festgelegt.

## 3. Package-Struktur gegen `spec/architecture.md` §2.1

`find cmd internal -type f` ergibt `cmd/pgwire-recorder/main.go` und je ein `doc.go` in `internal/hexagon/{model,services,ports/driving,ports/driven}`, `internal/adapters/driving/{cli,pgwire}`, `internal/adapters/driven/{postgres,recording}` und `internal/bootstrap`. Das deckt den Baum in §2.1 vollständig ab, mit gleichen Namen. Die `ARC-*`-Angaben in den Package-Kommentaren stimmen mit der Tabelle in §2.1 überein (`ARC-001` bis `ARC-009`).

**Abweichung:** `test/integration/` fehlt. Sie ist im Plan §1 („Ausdrücklich NICHT") begründet: Das Verzeichnis entsteht mit dem ersten Integrationstest in `slice-walking-skeleton-record`. Damit ist sie gedeckt. Review F-213 ist damit für diesen Slice erledigt. Ob die offenen Slices `test/integration` weiter als „update" statt „neu" führen, ist eine Frage an den Planner.

## 4. ADR-Konformität

| ADR | Entscheidung | Befund |
|---|---|---|
| [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) | Ports & Adapters verbindlich; Folgepflicht: Abhängigkeitsregeln maschinell prüfen | **konform.** Die Package-Struktur bildet das Hexagon ab. `make a-check` hängt an `make gates` und meldet Kantenverstöße (M5) und Fremdimporte im Kern (M1). Grenze: Kern-Reinheit wird nur für Bibliotheken in der `tech`-Liste geprüft (M6b, F-212). |
| [ADR-0009](../plan/adr/0009-implementierungssprache-go.md) | v1 in Go; `.a-check.yml` ist Go-spezifisch | **konform.** Das Go-Modul baut mit `golang:1.27-alpine` (Digest gepinnt), `CGO_ENABLED=0`. `.a-check.yml` führt `languages: go`. |
| [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) | `pgproto3` nur in den beiden PGWire-Adaptern; Fitness Function: `tech`-Regel | **konform** (seit `a834db1`). `pgproto3` und `crypto/tls` sind im Upstream-Adapter erlaubt (M2) und im Server-Adapter erlaubt (M3). Im CLI-Adapter (M4), im Recording-Adapter (Gegenprobe Fall 2) und im Model (M1) werden sie abgelehnt. Review F-204 ist damit behoben. Die Sensors-Bindung von `a-check-negativ` nennt jetzt [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) (F-206). |

Keine ADR wurde im Diff geändert (Hard Rule 3.5).

## 5. Gegenprobe dieser Datei

`make docs-check` nach dem Anlegen dieser Datei: `d-check: 89 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

## 6. Gesamturteil

**Die DoD-Liefer-Punkte 1 bis 3 sind bestätigt.** Grundlage sind eigene Läufe und eine eigene Mutationsprobe, nicht die Angaben des Implementers. Plan §3 und der Diff stimmen bis auf die nicht geplante, sachlich gedeckte Änderung an `README.md` überein. Die Package-Struktur entspricht §2.1, die Auslassung `test/integration/` ist im Plan begründet. [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0009](../plan/adr/0009-implementierungssprache-go.md) und [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) sind eingehalten.

**Übergabe an den Planner:**

1. Der Commit `a834db1` ist von keinem Review-Report erfasst. Der Report deckt `8143a47` ab. Vor der Closure ist zu entscheiden, ob ein Folge-Review nötig ist.
2. Reichweitengrenze des Architektur-Gates: Kern-Reinheit gilt nur für Bibliotheken in der `tech`-Liste (M6b). Beim Risiko „`os` im Model" aus §6 sollte der Ausgang diese allgemeinere Form nennen.
3. Der netzlose Build trägt nur ohne Abhängigkeiten (F-211, §6). Das betrifft `slice-walking-skeleton-record`.
4. F-205 (Commit-Hook gegen benannte Slices) liegt außerhalb dieses Slice und wird hier nicht bewertet.
