# Slice slice-harness-blackbox-driven: Black-Box-Tests der getriebenen Adapter (PostgreSQL und Recording)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Er trägt zum
Abnahmeszenario 17 (M3) bei, das `slice-lastenheft-pruefbarkeit` anlegt und das mit
dem letzten Slice der Reihe nachweisbar wird. Eingesammelt wird er von der
nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-05: nach `slice-replay-semantik-fehlerreplay` und vor dem nächsten großen
Slice (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** — (Harness-Arbeit; keine Produkt-Anforderung). Die neue Qualitätsanforderung an die Prüfbarkeit des Quellcodes trägt `slice-lastenheft-pruefbarkeit` hier nach, sobald sie im Lastenheft steht. Bindung: die ADR von `slice-harness-lint` (Stufung von `testpackage` mit Hochschalt-Trigger auf diesen Slice), [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (Abdeckungs-Deklarationen an den Tests bleiben unverändert).

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Unit-Tests unter `internal/adapters/driven/postgres` und `internal/adapters/driven/recording` laufen als Black-Box-Pakete
(`package <name>_test`) und prüfen über die exportierte Schnittstelle des Pakets;
`testpackage` ist für diese Pfade scharf. Wo ein Test einen unexportierten Teil
braucht, geht er über eine Export-Test-Brücke (`export_test.go`) in der Form, die die
ADR von `slice-harness-lint` festlegt, oder wird gegen die exportierte Schnittstelle
umgeschrieben. Die Fälle und ihre Prüfungen bleiben erhalten; der Slice ist ein Umbau
der Tests, keine neue Prüfung.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-05: Die Tests werden auf
Black-Box-Pakete umgestellt; bis dahin führt `slice-harness-lint` `testpackage`
gestuft, mit Hochschalt-Trigger auf die vier Umstellungs-Slices
(`slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
`slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`). Geschnitten ist
nach Paketgruppe, je in einer Schicht des Hexagons, damit jeder Schnitt einzeln
lieferbar und in einer Review-Sitzung prüfbar bleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Tests der übrigen Pakete — sie übernehmen `slice-harness-blackbox-kern`, `slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`.
- Neue Exporte im Produkt-Code nur für Tests — sie verbreiterten die Schnittstelle des
  Pakets für einen Zweck außerhalb des Produkts; der Weg für unexportierte Teile ist
  die Brücke in einer `_test.go`-Datei (§6). Braucht ein Test mehr, ist das ein Befund
  für den Architect, keine Entscheidung im Umbau (§4).
- Neue Fälle oder geänderte Erwartungen — ein anderer Vorgang; der Umbau soll an
  derselben Testliste messbar sein (§2). Die Schwelle der Testabdeckung übernimmt
  `slice-harness-coverage`, gemessen erst nach allen vier Umstellungs-Slices.
- Produkt-Verhalten, Spezifikation, Lastenheft — Schicht-Abgrenzung: Der Slice ändert
  Testdateien, `.golangci.yml` und die Gegenprobe des Lint-Gates.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Die Testdateien unter `internal/adapters/driven/postgres` und `internal/adapters/driven/recording` gehören zu `_test`-Paketen; die Liste der Tests
      (`go test -list .` je Paket, im gepinnten Go-Image) ist vor und nach dem Umbau gleich, keine Prüfung ist
      entfallen, die Abdeckungs-Deklarationen nach [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) stehen unverändert und
      `make abdeckung-check` ist grün. Unexportierte Teile erreicht ein Test nur über
      die Brücke nach der ADR des Lint-Gates.
- [ ] `testpackage` ist für `internal/adapters/driven/postgres` und `internal/adapters/driven/recording` scharf: Die Stufe aus `.golangci.yml` ist für diese
      Pfade aufgehoben (Hochschalt-Trigger nach der ADR), `make lint` ist grün, und die
      Gegenprobe des Lint-Gates führt einen Fall, in dem ein White-Box-Test in einem
      dieser Pfade rot wird; die Mutation ist gesehen (`AGENTS.md` §3.10).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/postgres/upstream_test.go`, `upstream_extended_test.go` | refactor | `package postgres_test` |
| `internal/adapters/driven/recording/yaml_test.go` | refactor | `package recording_test` |
| `export_test.go` je Paket, nur wo nötig | neu | Brücke zu unexportierten Teilen nach der ADR des Lint-Gates |
| `.golangci.yml` | update | Stufe von `testpackage` für diese Pfade aufheben |
| `tools/harness/lint-gegenprobe.sh` | update | Fall: White-Box-Test in einem dieser Pfade wird rot |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-blackbox-kern` liegt in `done/` (WIP-Limit 1,
Reihenfolge nach Entscheidung des Nutzers vom 2026-10-05 und 2026-10-06: Lastenheft, Lint,
die vier Umstellungs-Slices, `slice-harness-abdeckung-gate`, Coverage). Vor dem ersten Code-Commit prüft der Architect §6 gegen
die ADR des Lint-Gates (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Umbau ist nicht in einer
  Review-Sitzung prüfbar, etwa weil mehr als eine Handvoll Tests umgeschrieben statt
  umgestellt werden muss; dann je Paket ein eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): Ein Test erreicht einen Teil, den er
  prüfen muss, weder über die exportierte Schnittstelle noch über die Brücke, die die
  ADR zulässt; dann zuerst eine Entscheidung des Architect (Brücke erweitern oder
  Schnittstelle des Pakets ändern).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün mit scharfem `testpackage` für diese Pfade,
Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — offen, entschieden in der ADR von
`slice-harness-lint` vor dessen Code; was dort nicht steht, gibt der Implementer an
den Architect zurück.

- **Export-Test-Brücke** — ob sie zulässig ist und in welcher Form: `export_test.go`
  im Paket selbst (`package <name>`), das unexportierte Funktionen, Typen oder
  Konstanten unter exportiertem Namen für die `_test`-Dateien sichtbar macht; Name der
  Datei, Benennung der Aliase, und ob sie Zustand setzen darf (etwa eine Frist für
  einen Test verkürzen) oder nur lesen.
- **White-Box-Zugriffe im Bestand** — Namensabgleich per Suche am Stand `ce50a10`
  (ungemessen, kann Fehltreffer enthalten, wo ein Testhelfer gleich heißt): in `recording` keine; in `postgres` u. a. `session`, `toFrontendMessage`, `lies`.
  Je Zugriff: über die exportierte Schnittstelle prüfbar, über die Brücke, oder
  Befund. `toFrontendMessage` ist eine Übersetzung in `pgproto3`-Nachrichten; ob ein Test sie direkt oder über `Upstream` prüft, entscheidet die Brücke.
- **Fakes der Ports** — die Ports sind exportiert (`internal/hexagon/ports/...`); ein
  Fake in einem `_test`-Paket implementiert sie unverändert. Ein Fake, der auf
  unexportierte Felder eines Produkt-Typs greift, fällt unter die Brücke.
- **Mutationstests auf unexportierte Teile** — eine Mutation, die bisher ein Test auf
  eine unexportierte Funktion fing, muss auch nach dem Umbau fangen; sonst fehlt eine
  Prüfung, obwohl die Testliste gleich ist.

**Risiken:**

- **Prüfung geht im Umbau verloren** — ein Test bleibt in der Liste, prüft aber
  weniger, weil ein unexportierter Vergleich wegfiel
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). Je umgeschriebener (nicht nur
  umgestellter) Test eine Mutation, die er weiter fängt. — **Ausgang:** — (bei Closure)
- **Abdeckung sinkt** — Black-Box-Tests erreichen unexportierte Pfade seltener; die
  Zahl misst erst `slice-harness-coverage` nach allen vier Umstellungs-Slices.
  — **Ausgang:** — (bei Closure)
- **Gegenprobe sieht den Mutanten nicht** — Build-Kontext
  (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`, 1×); die Gegenprobe baut aus einer
  Kopie unter eigenem Temp-Pfad. — **Ausgang:** — (bei Closure)

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `ce50a10` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (8×, verkörpert in `AGENTS.md`
  §3.10) — das Scharfschalten von `testpackage` für diese Pfade ist eine neue Zusage
  des Lint-Gates und bekommt ihren roten Fall (§2).
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) und
  `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — je ein Risiko in §6.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — wellenlos; §4 nennt den
  Architect vor dem Code.
- `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (2×) — die Brücke ist
  genau eine Randform, die ein Umbau still entscheiden würde; sie steht in §6.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
