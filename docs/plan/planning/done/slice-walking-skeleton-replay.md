# Slice slice-walking-skeleton-replay: Replay-Pfad für SELECT 1

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-walking-skeleton.

**Bezug:** [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), [ADR-0007](../../adr/0007-strict-replay.md), [ADR-0008](../../adr/0008-kein-sql-parser-in-v1.md)

**Berührte Spec-Stellen:** `ARC-002` · `ARC-003` · `ARC-004` · `ARC-005` · `ARC-006` · `ARC-008` · `ARC-009` · `SPEC-011` · `SPEC-018` · `SPEC-027` · `SPEC-034` · `LH-FA-03.a` · `LH-FA-03.b` · `LH-FA-09.a` · `LH-FA-10.a` · `LH-FA-12.a`

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

**Ziel:** `pgwire-recorder replay` beantwortet `SELECT 1;` aus dem Recording des Record-Slice, ohne PostgreSQL, mit strict sequential matching.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Mismatch bei Extended-Nachrichten — `slice-replay-semantik-mismatch` (mit `slice-extended-query-replay`); hier liefert der Replay-Service für einfache Anfragen schon den Mismatch mit Diagnose nach LH-FA-10.a (`PGR-E5001`) und Exit-Code 5.
- Fehlerreplay — welle-replay-semantik.
- `--session-assignment connection` und `--fail-on-unconsumed` — `slice-v1-abschluss-sessions` und `slice-replay-semantik-mismatch`; hier gilt die Zuordnung `first-request`, nicht verbrauchte Interaktionen sind die Warnung `PGR-W2001`.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay): Mit gestoppter PostgreSQL liefert Replay dem Client für `SELECT 1;` dasselbe Ergebnis wie im Record (End-to-End-Test).
- [x] [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus): Zwei aufeinanderfolgende Replay-Läufe mit demselben Recording zeigen identisches Verhalten.
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
| `internal/hexagon/ports/driving` (Port `Replayer`), `internal/hexagon/services` (Replay-Service) | neu | Cursor, exakter SQL-Vergleich, Zuordnung `first-request`, Warnung `PGR-W2001` |
| `internal/adapters/driving/pgwire` | update | Server für Record und Replay über `NewRecordServer`/`NewReplayServer`; Handshake ohne Upstream; Warnungen als `model.Warning` |
| `internal/adapters/driving/cli` | update | Kommando `replay` |
| `internal/adapters/driven/recording` | update | Laden; die Felder `offset_ms` und `empty_sessions` der Version 1 werden gelesen |
| `internal/bootstrap` | update | Verdrahtung Replay |
| `internal/hexagon/services/replay_test.go`, `internal/adapters/driving/pgwire/server_test.go`, `test/integration/replay_e2e_test.go`, `tools/test/run-integration-tests.sh` | neu / update | Unit: Matching, Diagnose, Zuordnung, Handshake, Startfehler, Replay-Modus des Adapters; E2E: Record → zehn Replay-Läufe mit der Sicht beim Aufzeichnen, Abweichung mit Exit-Code 5, beschädigte Aufzeichnung; zweite Runner-Phase mit gestopptem PostgreSQL über ein gemeinsames Volume |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-walking-skeleton-record` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: der Replay-Handshake braucht mehr als die geplanten Authentifizierungsnachrichten — zurück zur Zerlegung.
- `in-progress` → `open`: das Recording enthält Nachrichten, die Replay nicht reproduzieren kann — Carveout.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben; der Welle-Smoke (Record, PostgreSQL stoppen, Replay) ist durchlaufen.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Authentifizierungsnachrichten, die typische Treiber im Replay akzeptieren, sind noch nicht festgelegt — **Ausgang:** entfallen: AuthenticationOk, sortierte Serverparameter der Session und ReadyForQuery tragen `pgconn` (E2E) und `psql` 17 (Review-Sonde); weitere Treiber prüft der Validator der Welle.

- Der Leser des Recording-Adapters lehnt die für Version 1 spezifizierten Felder `offset_ms`, `empty_sessions` und `type: extended` noch ab (`KnownFields`); Replay muss sie lesen oder gezielt ablehnen — **Ausgang:** eingetreten: `offset_ms` und `empty_sessions` liest der Leser jetzt; `type: extended` übernimmt `slice-extended-query-modell`.

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

- **Was hat funktioniert:** `replay` beantwortet Anfragen aus einer von `record` geschriebenen Aufzeichnung mit derselben Sicht des Clients wie beim Aufzeichnen, auch bei gestopptem PostgreSQL; strict sequential matching, Diagnose nach LH-FA-10.a und Exit-Code 5 laufen. Die zweite Runner-Phase belegt den DoD-Punkt „mit gestoppter PostgreSQL“ wörtlich.
- **Was ging anders als geplant:** Der Slice lieferte mehr als geplant (Mismatch einfacher Anfragen, Exit-Code 5, Zuordnung mehrerer Verbindungen); Plan, Welle und Folge-Slices wurden nachgezogen. Die handgeschriebene Fixture wich von einer echten Aufzeichnung ab und ist ersetzt.
- **Steering-Loop-Eintrag:** Sensor ergänzt: zweite Phase des Integrations-Runners mit gestopptem PostgreSQL, die die Aufzeichnung der ersten Phase gegen die Sicht beim Aufzeichnen prüft — liegt in `harness/mk/integration.mk` (Herkunft `· seit slice-walking-skeleton-replay`). Auslöser: Review F-280, Verifikation V-9 und V-11; `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (2×).
- **Beobachtungs-Register (`../observations/`):** `BEO-REPO/plan-folgt-korrektur-nicht/` Beleg `evidence/slice-walking-skeleton-replay.md` ergänzt — Zähler 3×, geht in die Closure von `welle-walking-skeleton`; `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` neu mit Belegen aus record und replay — 2×.
- **Folge-Slices:** `slice-replay-semantik-mismatch` (Extended-Mismatch, `--fail-on-unconsumed`), `slice-v1-abschluss-sessions` (`--session-assignment connection`, `--record-empty-sessions`), `slice-extended-query-modell` (Leser für `type: extended`) — alle Dateien in `open/`.
- **Risiken aus §6:** jedes mit genau einem Ausgang — siehe §6.
- **Drei Paarungen:** Repo mit Wellen-Betrieb — geprüft von der Closure von `welle-walking-skeleton`.
- **Belege:** Review `docs/reviews/2026-10-04-review-slice-walking-skeleton-replay.md`, Verifikation `docs/reviews/2026-10-04-verifikation-slice-walking-skeleton-replay.md`; die Korrekturen nach der Verifikation (Phase 2 mit echter Aufzeichnung, volle Sicht, Rücknahme von LH-QA-01) sind durch `make gates` belegt, kein weiteres Review.

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
