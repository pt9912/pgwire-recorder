# Slice slice-walking-skeleton-record: Record-Pfad für SELECT 1

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-walking-skeleton.

**Bezug:** [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-06`](../../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), [ADR-0003](../../adr/0003-pgwire-server-ist-driving-adapter.md), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0005](../../adr/0005-recording-store-ist-driven-adapter.md), [ADR-0006](../../adr/0006-kanonisches-domain-model.md)

**Berührte Spec-Stellen:** `ARC-001` · `ARC-002` · `ARC-003` · `ARC-004` · `ARC-005` · `ARC-006` · `ARC-007` · `ARC-008` · `ARC-009` · `ARC-013` · `SPEC-001` · `SPEC-041` · `SPEC-045` · `LH-FA-05.e` · `SPEC-002` · `SPEC-003` · `SPEC-004` · `SPEC-034` · `LH-FA-02.a` · `LH-FA-02.b` · `LH-FA-05.c` · `LH-FA-05.e` · `LH-FA-07.a` · `LH-FA-13.b`

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

**Ziel:** `pgwire-recorder record` vermittelt `SELECT 1;` zwischen einem Client und einer realen PostgreSQL-Instanz und schreibt ein gültiges, versioniertes YAML-Recording.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Replay — `slice-walking-skeleton-replay`; hier genügt ein Roundtrip-Test des Recording-Adapters.
- Fehlerreplay — welle-replay-semantik; Fehlerantworten des Servers werden hier nur vermittelt und aufgezeichnet.
- Zuordnung und Kennungen paralleler Sessions nach LH-FA-12.a — welle-v1-abschluss; hier werden gleichzeitige Verbindungen vermittelt und in der Reihenfolge ihres Endes übernommen, ohne dass LH-FA-12 als belegt gilt.
- Signalbehandlung nach LH-FA-13.a (zweites Signal, Abschluss der laufenden Session) — welle-v1-abschluss; hier endet der Lauf auf SIGINT oder SIGTERM nach der laufenden Interaktion, ein zweites Signal beendet sofort.
- Anmeldeverfahren des Upstreams — hier nur ein Upstream ohne Passwort-Anmeldung (siehe §6).


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-06`](../../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten): Ein Client führt `SELECT 1;` über den Recorder aus und erhält das Ergebnis der realen Instanz; die Interaktion steht geordnet im Recording (Integrationstest).
- [x] [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings): Das Recording ist persistent, trägt Formatkennung und Version und lässt sich per Roundtrip laden (Adapter-Contract-Test).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der nächsten Welle-Closure geprüft.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/model` | neu | Recording, Session, Interaction, Request, Response, Value |
| `internal/hexagon/services` (Record-Service) | neu | Interaktion aus Request und Responses bilden |
| `internal/hexagon/ports/driving`, `…/driven` | neu | Record-Use-Case, Recording-Repository, PostgreSQL-Upstream |
| `internal/adapters/driving/pgwire`, `…/cli` | neu | Startup samt Protokollrand, `Query`, `Terminate`, Verbindungsende; Kommandos `record` und `version` |
| `internal/adapters/driven/postgres`, `…/recording` | neu | Upstream-Session, YAML-Schreiben |
| `internal/bootstrap` | neu | Verdrahtung |
| `test/integration` | neu | Black-Box-E2E gegen PostgreSQL: `SELECT 1;`, mehrere Interaktionen, SSL-Ablehnung, Verbindung ohne Anfrage, Ende ohne Terminate, nicht unterstützte Interaktion, Upstream nicht erreichbar, Beenden mit offener Verbindung, Schreibfehler am Ende |
| `tools/test/run-integration-tests.sh`, `harness/mk/integration.mk` | neu | Der Integrationstest-Runner startet PostgreSQL und das Testimage in einem eigenen Docker-Netz, räumt alles auf und schreibt nichts in den Arbeitsbaum |
| `Dockerfile`, `.dockerignore`, `harness/mk/build.mk`, `go.mod`, `go.sum` | neu / update | Multistage-Build nach [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md): Download-Stufe mit Netz, Build-, Test- und Integrationsstufe netzlos, Produkt-Image als `runtime`; `make go-mod-tidy` als Werkzeug |
| `.a-check.yml` (`tech`) | update | YAML-Bibliothek `go.yaml.in/yaml/v3` nach [ADR-0027](../../adr/0027-yaml-bibliothek.md); beide Modulpfade auf den Recording-Adapter beschränkt |
| `tools/test/abdeckung.sh`, `tools/test/abdeckung-gegenprobe.sh`, `harness/mk/abdeckung.mk`, `docs/user/abdeckung-*.md`, `.d-check.yml` (`trace.coverage`) | neu | Abdeckung je Anforderung und Pfad nach [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md); `trace.coverage` liest nur die vollständig belegten Anforderungen; die Gegenprobe hält das Skript |
| `internal/adapters/*/…_test.go`, `internal/hexagon/services/record_test.go` | neu / update | Protokollrand, Upstream-Fälle, Session-Ende, Nebenläufigkeit, Leser-Prüfungen |
| `harness/mk/vorgaben.mk` | neu | Verweis-Nachzug von `make slice-mv` lässt die Review-Reports aus |
| `cmd/pgwire-recorder/main.go` | update | Signal-Kontext, zweites Signal mit Standardverhalten, Version für `version` |
| `spec/spezifikation.md` (`LH-FA-05.e`, `SPEC-041`, `SPEC-045`, `SPEC-034`) | update | Protokollrand für fremde erste Nachrichten (`PGR-W3003`), beide Formen von `null` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-walking-skeleton-build-gates` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: der Slice braucht mehr als drei Liefer-Punkte (etwa weil der Startup-Handshake mehr Auth-Verfahren verlangt) — zurück zur Zerlegung.
- `in-progress` → `open`: eine PGWire-Nachricht der Startup-Sequenz ist mit `pgproto3` nicht verlustfrei abbildbar — Carveout.

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

- Startup-/Authentifizierungsverfahren des Upstreams (z. B. SCRAM) sind für die Weiterleitung noch nicht eingegrenzt; dieser Stand vermittelt nur einen Upstream ohne Passwort-Anmeldung und meldet jedes Anmeldeverfahren mit `PGR-E6001` — **Ausgang:** eingetreten: übernimmt `slice-v1-abschluss-anmeldung`.
- Eine erste Nachricht, die keine PGWire-Startnachricht ist (HTTP-Probe, Port-Scan), war in der Spezifikation nicht festgelegt — **Ausgang:** entfallen: `LH-FA-05.e` legt sie als Warnung `PGR-W3003` ohne Wirkung auf den Exit-Code fest (`SPEC-045`), umgesetzt und getestet.
- Die Aufzeichnung wird nach jeder beendeten Session und am Ende des Laufs geschrieben; ein Abbruch per SIGKILL verliert die laufenden Sessions — bewusst, Behandlung in welle-v1-abschluss — **Ausgang:** eingetreten: übernimmt `slice-v1-abschluss-betrieb` (Signalbehandlung, LH-FA-13.a).
- `make build` und `make test` laufen netzlos; mit `pgproto3` braucht das Modul eine netzlose Quelle der Abhängigkeiten (zum Beispiel ein vendored Verzeichnis), sonst scheitert der Build — **Ausgang:** entfallen: [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) lädt die Module in der Download-Stufe des Dockerfiles; Build und Tests laufen danach netzlos.

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

- **Was hat funktioniert:** `record` vermittelt einfache Anfragen gegen eine reale PostgreSQL-Instanz und schreibt eine geprüfte YAML-Aufzeichnung; zehn Black-Box-E2E-Tests und Unit-Tests je Adapter tragen den Protokollrand. Das Multistage-Dockerfile baut Binary, Tests und Produkt-Image netzlos. Drei Review-Runden und die Verifikation fanden je echte Mängel (verworfene Session wurde aufgezeichnet, Exit-Code 0 nach Fehlern, abgeschnittene Datei als gültig).
- **Was ging anders als geplant:** Der Record-Port braucht ein Session-Ende mit Grund; die Abdeckung wurde von einer Nennung je Anforderung auf Pfade umgestellt ([ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md)); drei ADRs kamen hinzu ([ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md), [ADR-0027](../../adr/0027-yaml-bibliothek.md), [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md)); die Spezifikation regelt jetzt fremde erste Nachrichten. Die Anmeldung beim Upstream ist nicht vermittelt (Folge-Slice).
- **Steering-Loop-Eintrag:** Sensor ergänzt: Gegenprobe des Abdeckungs-Skripts — liegt in `harness/mk/abdeckung.mk` (Herkunft `· seit slice-walking-skeleton-record`). Auslöser: Folge-Review F-259 und `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, keine Registerklasse mit 3×.
- **Beobachtungs-Register (`../observations/`):** `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` Beleg `evidence/slice-walking-skeleton-record.md` ergänzt — Zähler 2×; `BEO-REPO/plan-folgt-korrektur-nicht/` neu mit zwei Belegen (build-gates, record) — 2×; `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt/` neu — 1×; `BEO-REPO/record-fehlerantwort-im-aufbau-ungeregelt/` neu — 1×.
- **Folge-Slices:** `slice-v1-abschluss-anmeldung` (Anmeldung des Clients im Record-Modus vermitteln) — ist eine Datei in `open/`; `slice-walking-skeleton-replay` trägt das Risiko der noch abgelehnten Felder der Version 1.
- **Risiken aus §6:** jedes mit genau einem Ausgang — siehe §6.
- **Drei Paarungen:** Repo mit Wellen-Betrieb — geprüft von der Closure von `welle-walking-skeleton`.
- **Belege:** Review `docs/reviews/2026-10-04-review-slice-walking-skeleton-record.md`, Folge-Review `docs/reviews/2026-10-04-review-slice-walking-skeleton-record-folge.md`, Verifikation `docs/reviews/2026-10-04-verifikation-slice-walking-skeleton-record.md`; die Korrekturen nach dem Folge-Review (fremde erste Nachrichten, abgeschnittene Dateien, ehrliche Deklarationen) hat die Verifikation geprüft, kein drittes Review.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert in `harness/conventions.md` noch keine Sub-Areas; berührt wird das gesamte Repo als eine Sub-Area. Sub-Area-Wahl: eine Sub-Area für das gesamte Repo.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
