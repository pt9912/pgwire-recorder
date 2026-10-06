# Slice slice-replay-semantik-meldungscodes: Meldungscodes, Fehlertext und Logging

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-replay-semantik.

**Bezug:** [`LH-FA-14`](../../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md)

**Berührte Spec-Stellen:** `SPEC-034` · `SPEC-005` · `SPEC-006` · `SPEC-013` bis `SPEC-028` · `LH-FA-18.a` · `SPEC-033`

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

**Ziel:** Jeder Fehler und jede Warnung trägt einen `PGR-…`-Meldungscode; der Fehlertext beginnt mit `<klasse> [<code>]: `, Logs gehen nach `stderr` mit einstellbarem Level.

**Übernommen aus `slice-extended-query-replay`** (Review F-328, Validierung `docs/reviews/2026-10-05-validierung-slice-extended-query-replay.md`, Befund 3): Mit `--log-level` wird die Regel zu Parameterwerten in der Diagnose wirksam, und die Spezifikation fasst sie zweimal verschieden — `LH-FA-18.a` §Mismatch verbietet Klartext-Werte nur, „wenn der Log-Level nicht `debug` ist“, `SPEC-033` verbirgt sie ohne Bedingung. Die Diagnose geht zudem als `ErrorResponse` an den Client, ein Log-Level wirkte also auch dort. Dazu nennt die Diagnose einer Parameter-Abweichung den Index des abweichenden Parameters nicht, obwohl er keinen Wert verrät. Beides entscheidet der Nutzer als Änderung von `LH-FA-18.a` §Mismatch und `SPEC-033`, bevor dieser Slice beginnt (§4); die Umsetzung liefert DoD-Punkt 3.

**Übernommen aus `slice-replay-semantik-fehlerreplay`** (Verifikation `docs/reviews/2026-10-06-verifikation-slice-replay-semantik-fehlerreplay.md`, V-52): Abnahmeszenario 4 ist für Simple Query nur zur Hälfte nachgewiesen. `TestE2EReplayAbweichung` prüft `PGR-E5001` und Exit-Code 5, deklariert auch „keine Antwort“, verwirft aber die Ergebnisse von `ReadAll`; `TestReplayMismatch` deklariert `LH-FA-10/Negative` und prüft nur Code und Cursor. Die Mutation „bei abweichender Anfrage liefert das Replay die Antworten der aufgezeichneten Anfrage, der Driving-Adapter sendet sie vor `PGR-E5001`“ bleibt grün (VF24). Der Code ist richtig, es fehlt nur der Nachweis für `LH-FA-10` Negative („statt eine unpassende aufgezeichnete Antwort zu verwenden“). Ohne ihn schließt welle-replay-semantik nicht (Closure-Trigger, Abnahmeszenario 4). Der Nachweis gehört in DoD-Punkt 3, weil er dieselbe Anforderung und dieselbe Diagnose einer Abweichung betrifft; er braucht keinen Produktionscode.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Gate, das Code-Tabelle und Katalog abgleicht — Folge-Slice, sobald der Katalog in der Betriebsdokumentation (welle-v1-abschluss) steht.
- Warnungen für nicht verbrauchte Interaktionen und `CancelRequest` — welle-v1-abschluss.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-14`](../../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben): Fehlerklassen aus `SPEC-020` bis `SPEC-028` liefern ihren Code im Kopf des Fehlertexts und als Log-Attribut (Test je Klasse).
- [ ] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Die Klasse eines Fehlers entspricht `SPEC-013` bis `SPEC-019` (Test je Klasse, deren Auslöser in dieser Welle existiert; die Abbildung auf den Exit-Code beim Herunterfahren und die übrigen Klassen prüft `slice-v1-abschluss-betrieb`).
- [ ] [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage): Die Diagnose einer Parameter-Abweichung nennt den Index des abweichenden Parameters ohne seinen Wert, und Parameterwerte erscheinen in Diagnose und `ErrorResponse` so, wie die vor Beginn bestätigte Fassung von `LH-FA-18.a` §Mismatch und `SPEC-033` es sagt, auf jedem Log-Level (Test). Bei einer abweichenden einfachen Anfrage erhält der Client keine aufgezeichnete Antwort, nur `PGR-E5001` (Abnahmeszenario 4, `LH-FA-10` Negative, übernommen aus `slice-replay-semantik-fehlerreplay`): `TestE2EReplayAbweichung` verlangt, dass `ReadAll` kein Ergebnis liefert, `TestReplayMismatch`, dass `Query` bei Abweichung keine Antworten liefert; beide werden rot, wenn das Replay die Antworten der aufgezeichneten Anfrage vor `PGR-E5001` liefert (E2E und Unit, Abdeckung `LH-FA-10/Negative`).
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
| Code-Tabelle (Paket im Core) | neu | Quelle der Wahrheit der Codes |
| `internal/adapters/driving/cli` | update | Kopf, Exit-Code-Abbildung, `--log-level` |
| Tests je Fehlerklasse | neu | Happy/Negative nach `SPEC-034` |
| `test/integration/replay_e2e_test.go` | update | DoD 3, Abnahmeszenario 4: `TestE2EReplayAbweichung` prüft die Ergebnisse von `ReadAll` (keine), Deklaration `LH-FA-10/Negative` ergänzt |
| `internal/hexagon/services/replay_test.go` | update | DoD 3: `TestReplayMismatch` prüft, dass `Query` bei Abweichung keine Antworten liefert |
| `docs/user/abdeckung-*.md` | update | von `make abdeckung` aus den geänderten Deklarationen geschrieben |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-replay-semantik-mismatch` ist `done`, und die Fassung von `LH-FA-18.a` §Mismatch und `SPEC-033` zu Parameterwerten und Parameter-Index ist vom Nutzer bestätigt (§1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: der Katalog verlangt eigene Doku-Arbeit — zurück zur Zerlegung.
- `in-progress` → `open`: Eine Klasse lässt sich nicht eindeutig zuordnen — Carveout.


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

- Vorrangfolge bei Fehlerketten mit mehreren Codes ist nicht festgelegt — **Ausgang:** offen bis Closure.

- Parameterwerte bei `debug` (aus `slice-extended-query-replay`, Review F-328): Liest man `LH-FA-18.a` §Mismatch als Erlaubnis für `debug`, verletzt das `SPEC-033`, und die Werte gingen auch an den Client — **Ausgang:** offen bis Closure.

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
