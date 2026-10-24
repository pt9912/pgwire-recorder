# Slice slice-harness-gegenproben-grenze-rest: Die übrigen Gegenproben werden gemessen, die unter zehn Minuten aufgenommen, die anderen geschnitten

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Entscheidung des Planners, 2026-10-21). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist, und keine Anforderung hängt an ihm; er steht neben `slice-harness-gegenproben-schneller`, aus dessen Rückgabe 4 er hervorgeht, und ist wie jener nicht in `welle-v1-nachtrag` aufgenommen.

**Bezug:** kein Lastenheft-Bezug — der Slice ändert Werkzeuge des Harness, keinen Anspruch des Produkts. Bindung an das Register: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Laufzeit der Gegenproben), `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` und `BEO-REPO/lastabhaengige-zaehlung` (die Zeit hängt an der Last des Hosts). Entscheidungen: [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Obergrenze der Zeit für jede Gegenprobe; Accepted, 2026-10-21) und [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (die Gegenproben stehen nicht in `GATE_CHECKS`; hier unverändert).

**Berührte Spec-Stellen:** `SPEC-060` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-060-gegenproben.md); Punkt 16, Tabelle der Aufnahme) · `SPEC-058` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-058-container-start.md); `make container-start-gegenprobe`, Messung am Host). Die Zusagen zu Zeit und Platz stehen in `SPEC-060` Punkt 16; die Verträge der übrigen Gegenproben nennt die Zeile je Gegenprobe in `harness/README.md` §Sensors und der Architect je aufzunehmender Gegenprobe in §3.

**Verantwortlich:** —
**Autor:** Planner. **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Jede der 18 Gegenproben, die `slice-harness-gegenproben-schneller` nicht aufnimmt, hat eine gemessene Zeit: die 17 der Liste `GEGENPROBEN_CONTAINER` ohne `release-gegenprobe` und `test-container-gegenprobe` (`lint-gegenprobe`, `stop-hook-gegenprobe`, `gegenproben-nach-einlieferung-gegenprobe`, `ci-gegenprobe`, `coverage-gegenprobe`, `gates-schnell-gegenprobe`, `gegenproben-gegenprobe`, `workflows-check-gegenprobe`, `range-konfiguration-gegenprobe`, `integration-volumes-gegenprobe`, `test-postgres-versionen-gegenprobe`, `archiv-bezug-gegenprobe`, `abdeckung-gegenprobe`, `kopf-check-gegenprobe`, `meldungskatalog-check-gegenprobe`, `hook-gegenprobe`, `a-check-negativ`) und `container-start-gegenprobe` der Liste `GEGENPROBEN_HOST`. Jede, deren voller Lauf auf ruhigem Host höchstens **600 s** dauert ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md), Entscheidung des Nutzers vom 2026-10-21: „die gegenproben werden zeitlich begrenzt. höchstens 10 minuten“), ist in der Tabelle von `SPEC-060` Punkt 16 `aufgenommen` und trägt die Fälle von Punkt 16 (9); jede, die darüber liegt, bekommt einen eigenen Slice in `open/` und in der Tabelle dessen Kennung statt der dieses Slice. Damit steht keine Gegenprobe mehr im Zustand `nicht gemessen`, und die Freigabe eines Tags kann die Punkte 2 und 3 der Checkliste wieder voll fahren (Auflösung von `CO-001`).

**Sendung aus dem Carveout `CO-001`** (Geber: [`CO-001`](../../carveouts/CO-001-gegenproben-am-getaggten-commit.md), Aktiv seit 2026-10-21; `AGENTS.md` §3.13): Dieser Slice ist die Adresse dafür, dass die übrigen Gegenproben aufgenommen werden; ohne ihn bliebe die Pflicht zum vollen Lauf am getaggten Commit über `v0.1.2` hinaus ausgesetzt. Trigger von `CO-001` ist, dass die Tabelle jede Gegenprobe als `aufgenommen` führt; die DoD unten trägt ihn. `grep` nach `CO-001` findet die Sendung in diesem §1 und in der DoD.

**Sendung aus dem Carveout `CO-003`** (Geber: [`CO-003`](../../carveouts/CO-003-gegenproben-am-wellenende.md), Aktiv seit 2026-10-24; Pflicht am Wellenende): `CO-003` hat denselben Auflösungs-Trigger wie `CO-001` (Tabelle jede Gegenprobe `aufgenommen`, jeder Nehmer in `done/`) und nennt diesen Slice als Folge-Slice. Die DoD trägt die Zeile; `grep` nach `CO-003` findet die Sendung in §1, DoD und §8. Keine neue Arbeit: Löst `CO-001` auf, löst `CO-003` mit auf.

**Sendung aus `slice-harness-container-start-u-boot-0-13`** (Geber: `open/`, eingetragen vom Planner am 2026-10-21, `AGENTS.md` §3.13; Nachzählung in §8). Dieser Slice **nimmt an**, für `container-start-gegenprobe` (Liste für den Host, in §1 oben und in der Messung der DoD schon geführt):

- **Die Zeit der Gruppe Docker unter `u-boot` 0.13.0**, die jener Slice misst und in seinem §7 belegt, als **Vormessung** — sie ersetzt den vollen Lauf dieser DoD nicht (ohne `GP_MUT`, ohne Auswahl, auf ruhigem Host).
- **Der gemessene Stand des vollen Laufs: 19 min 38 s** (264 Mutationen, Gruppe Docker eingeschlossen, am Host am 2026-10-18; Beleg `../observations/BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/evidence/slice-harness-devcontainer-einrichten-start.md`). Daraus folgt für die Aufnahme: `container-start-gegenprobe` liegt mit mehr als dem Dreifachen über 600 s, wird also **geschnitten** statt aufgenommen (Liefer-Punkt 3), solange die eigene Messung dieses Slice es nicht widerlegt.
- **Die Messung läuft unter der Version, die jener Slice in den Kopf von `SPEC-058` einträgt** (`u-boot` 0.13.0 statt 0.10.0); liegt am Host eine andere, nennt der Bericht sie.

**Sendung aus `slice-harness-gegenproben-schneller`** (Geber: dieser Slice entsteht aus dessen §6 Rückgabe 4 und aus `SPEC-060` Punkt 16): Der Slice misst jede der 18, nimmt die unter 600 s auf und schneidet die übrigen nach Gegenprobe. Der Eintrag in der Tabelle von `SPEC-060` Punkt 16 nennt ihn als Nehmer; bis zu seinem Anlegen war der Eintrag ein Fehler des Vertrags (Punkt 16 „Aufnahme“).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Aufnahme von `release-gegenprobe` und `test-container-gegenprobe` — `slice-harness-gegenproben-schneller` (Klasse 1; nimmt in §1 und DoD an). Grund: Beide sind gemessen und gelesen, und ihre Beschleunigung ist dort der Liefergegenstand.
- Das Beschleunigen einer Gegenprobe, die über 600 s liegt — ein eigener Slice je Gegenprobe, den der Planner nach der Messung anlegt (Klasse 3, anderer Vorgang: Arbeit am Werkzeug, nicht Messung und Aufnahme). Grund: Ein Slice vor der Messung wäre ein toter Slice (`modul-05-planning-harness.md` §Regeln gegen typische Fehlannahmen); Zahl und Zuschnitt hängen an dem, was die Messung zeigt. Die Tabelle nennt dann dessen Kennung.
- Eine Mutation oder einen Fall streichen oder zusammenlegen, die Obergrenze erhöhen oder die roten Fälle nur unter einer Auswahl prüfen — eine Lockerung, die eine ADR verlangt (Klasse 3; `AGENTS.md` §3.6, `SPEC-060` Punkt 16 (9)). Grund: Aufnahme heißt gleiches Ergebnis in 600 s, nicht weniger Prüfung.
- Ein Nachtlauf der Gegenproben in der CI — `welle-ci-nachtlaeufe` und `slice-harness-ci-gegenproben` (Klasse 1, `open/`). Grund: neue Außenwirkung des Harness, die eine eigene ADR braucht.
- Die Aufnahme einer Gegenprobe in `GATE_CHECKS` oder `make gates` — Klasse 3 ([ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)). Grund: Aufnahme in die Zeitgrenze macht eine Gegenprobe nicht zum Gate der Kette.
- Der Gegenstand `slice-harness-gegenprobe-baumkopie` (Kopie des Arbeitsbaums durch einen gemeinsamen Helfer, `open/`) — Klasse 1 des anderen Slice. Grund: Berührt eine Aufnahme dieselbe Stelle, nennt der Architect das und gibt es an den Planner.
- Das Auflösen von `CO-001` vor dem Trigger und das Setzen eines Tags — Klasse 3 (anderer Vorgang). Grund: Der Carveout wandert erst nach `done/`, wenn die Tabelle jede Gegenprobe als `aufgenommen` führt; ein Tag ist Handlung der Freigabe.
- Produkt-Code (`cmd/`, `internal/`) und `Dockerfile` — Schicht-Abgrenzung (Klasse 4): Der Slice berührt allein Skripte unter `tools/harness/`, die Fragmente unter `harness/mk/`, die Zeilen in `harness/README.md` und (durch den Architect) die Tabelle in `SPEC-060`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Messung.** Der volle Lauf jeder der 18 Gegenproben (17 im Gegenproben-Container, `container-start-gegenprobe` am Host) ist gemessen, ohne `GP_MUT` oder Auswahl, auf ruhigem Host (Last unter einem Fünftel der Kernzahl, `SPEC-060` Punkt 16 (1)), je mit Wanduhr in Sekunden, Last (`/proc/loadavg`, erster Wert), Hash und Datum im Bericht; die Tabelle in `SPEC-060` Punkt 16 trägt je Gegenprobe die Zahl (durch den Architect). Ein Lauf unter Last gilt nicht und wird auf ruhigem Host wiederholt.
- [ ] **Aufnahme der Gegenproben mit höchstens 600 s.** Jede trägt die Fälle und Mutationen von `SPEC-060` Punkt 16 (9) in der Form, die der Architect für sie festlegt (§6, Rückgabe 1), mit Ergebnisgleichheit der Deklarationszeilen zu `<Basis>` und einem vollen Lauf am Endstand; ihre Zeile in der Tabelle steht auf `aufgenommen`, ihre Zeile in `harness/README.md` §Sensors nennt die Obergrenze erst mit der Lieferung (`AGENTS.md` §3.11).
- [ ] **Schnitt der übrigen und Stand von `CO-001`.** Für jede Gegenprobe über 600 s liegt ein Slice in `open/` (angelegt vom Planner nach der Messung, Kennung `slice-harness-gegenproben-schneller-<name>`), der sie aufnimmt und die Sendung `CO-001` in §1 und DoD trägt; die Tabelle nennt seine Kennung, und das Feld `Folge-Slice` von `CO-001` führt sie (`grep` nach `CO-001` findet sie in jedem). **Sendung `CO-001`:** Die Closure nennt, welche Zeilen der Tabelle noch nicht `aufgenommen` sind; sind alle es, ist `CO-001` aufgelöst (reiner `git mv` nach `docs/plan/carveouts/done/`, Index im Commit davor), sonst bleibt er Aktiv mit aktualisierter `Letzte Prüfung`. **Sendung `CO-003`:** dasselbe für `docs/plan/carveouts/CO-003-*.md` (gleicher Trigger, derselbe `git mv`-Weg, Index im Commit davor); `grep` nach `CO-003` findet die Zeile.
- [ ] `make gates` grün; jede geänderte Gegenprobe am Endstand in einem Zug (`AGENTS.md` §6); die Gegenproben, die der Slice nicht ändert, bleiben unberührt.
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

Der Architect füllt diese Tabelle nach dem Lesen der Skripte und vor dem ersten Code-Commit (Messzeiten, Zuschnitt je Gegenprobe, Schätzung in der Metrik von §4). Der Planner nennt nur den Rahmen:

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/dev/harness/*-gegenprobe.sh` und `tools/harness/<name>.sh` der aufzunehmenden Gegenproben (Liste: Messung) | update | Zeitbudget, Ausgabe und Fälle nach `SPEC-060` Punkt 16 (9) in der Form des Architect; nur für Gegenproben unter 600 s |
| `harness/mk/*.mk` der aufzunehmenden Gegenproben | update, falls das Fragment `GP_BUDGET` o. ä. durchreicht | nur wo die Form des Architect es verlangt |
| `harness/README.md` (§Sensors, Zeilen der aufgenommenen Gegenproben) | update, erst mit der Lieferung | nennt Obergrenze und Schalter; vorher nicht (`AGENTS.md` §3.11) |
| `harness/vertraege/SPEC-060-gegenproben.md` (Tabelle Punkt 16) | durch den Architect | Zahl je Gegenprobe, `aufgenommen`, Kennung der Folge-Slices |
| `docs/plan/planning/open/` (Slices für Gegenproben über 600 s), `docs/plan/carveouts/CO-001-*.md` | durch den Planner | Schnitt und Stand des Carveouts |

Nicht in dieser Liste: Produkt-Code, `Dockerfile`, `release-gegenprobe`, `test-container-gegenprobe`, CI (§1).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next`): `slice-harness-gegenproben-schneller` liegt in `done/` (er liefert die Form der Fälle von `SPEC-060` Punkt 16 (9), die dieser Slice für die übrigen übernimmt, und die Erfahrung, welche Mittel die Zeit senken), und der Nutzer priorisiert den Slice. Reihenfolge (Empfehlung des Planners; der Nutzer will zuerst den Release, Entscheidung des Nutzers vom 2026-10-21: „ich gebe alles frei, damit wir jetzt endlich diesen brew release machen können“): erst `slice-v1-abschluss-release-v0-1-1`, dann `slice-harness-gegenproben-schneller`, dann dieser Slice. Der Slice setzt den Release nicht voraus. Die Messung am Host (`container-start-gegenprobe`) braucht einen Host mit Docker, `u-boot`, `devcontainer`, `node` und `ssh` und läuft nie gleichzeitig mit dem Lauf eines anderen Slice auf demselben Daemon (`harness/README.md` §Sensors).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness/mk harness/README.md harness/vertraege` (Tests und Gegenproben eingeschlossen, `AGENTS.md` §3.14) **über 1000** liegt; `<Basis>` ist der Hash von `main` beim Übergang nach `in-progress/`, den der Architect in §3 einträgt. Die Schätzung des Architect in dieser Metrik (Faktor von Tabellen- und Gegenprobe-Code zum Produktcode rund zwei) steht vor dem Code in §3. **Vorgemerkter Schnitt, weil die Zahl der aufzunehmenden Gegenproben erst die Messung nennt:** Liegt die Schätzung nach der Messung über der Marke, schneidet der Planner nach Gegenprobe-Gruppen (Skripte mit gemeinsamem Treiber) in Slices `slice-harness-gegenproben-grenze-rest-<gruppe>`; dieser Slice behält Messung, Tabelle und Schnitt (Liefer-Punkte 1 und 3), die Aufnahme geht je Gruppe weiter.
- `in-progress` → `open`: Der Host für `container-start-gegenprobe` steht nicht zur Verfügung (Docker, `u-boot`, `devcontainer`, Zeit am Host): Die Gegenprobe bleibt `nicht gemessen` mit eigenem Nehmer, den der Planner anlegt, statt den Slice zu halten; die übrigen 17 laufen weiter.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben; die Tabelle von `SPEC-060` Punkt 16 führt keine Gegenprobe mehr im Zustand `nicht gemessen`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12; der Planner nennt sie, der Architect entscheidet sie vor dem ersten Code-Commit im Vertrag):

| Randform | Frage | Entscheidung |
|---|---|---|
| (a) Form der Fälle für Gegenproben ohne `GP_*` | `SPEC-060` Punkt 16 (9) beschreibt Fälle in Begriffen von `release-gegenprobe` (`GP_BUDGET`, `GP_MUT`, Mutationen mit roten Fällen); die übrigen Gegenproben haben diese Schalter und Deklarationen nicht oder anders | **offen, Architect** (Rückgabe 1) |
| (b) Messung unter Last | Zahl nur auf ruhigem Host gültig (`SPEC-060` Punkt 16 (1)); wie wird die Last der Messung belegt und was gilt, wenn keine Fahrt ruhig ist | **entschieden im Vertrag:** die Zeile nach (2) trägt die Last; eine Messung unter Last gilt nicht (DoD 1) |
| (c) Messung am Host | `container-start-gegenprobe` läuft nur am Host (Docker, `u-boot`, `devcontainer`) und nie gleichzeitig mit anderen Läufen auf demselben Daemon | **entschieden:** Messung am Host, ein Lauf zur Zeit; fehlt der Host, Rückführung nach `open` (§4) |
| (d) Gegenproben mit Docker oder Netz | `integration-volumes-gegenprobe`, `test-postgres-versionen-gegenprobe`, `workflows-check-gegenprobe`, `lint-gegenprobe` hängen an echtem Docker oder Images; Zeit hängt an Cache und Pull | **offen, Architect:** ob die Messung mit warmem Cache gilt und wie das im Bericht steht (Rückgabe 2) |
| (e) Plattform | Zeit und Platz sind auf dem Gegenproben-Container (Linux) gemessen; macOS und Windows nicht zugesagt | **entschieden:** Grenze aus dem Kopf von `SPEC-060`; `SPEC-060` Punkt 16 (6) |
| (f) Aufrufarten, Abbruch | `make`-Flags `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS`, `SIGINT`, `SIGTERM` je Gegenprobe | **entschieden:** wie `SPEC-060` Punkt 16 (7) und das Verhalten der jeweiligen Gegenprobe (Bestand); der Architect nennt je aufzunehmender Gegenprobe Abweichungen |

**Rückgaben an den Architect** (nach dem Lesen, vor dem Code):

1. Welche Form von Punkt 16 (9) gilt für jede der übrigen Gegenproben (Budget, Zeile nach (2), Arbeitsordner), wo sie die Schalter von `release-gegenprobe` nicht haben (Randform (a))?
2. Randform (d) (Cache und Pull bei Gegenproben mit Docker).
3. Schätzung in der Metrik von §4 nach der Messung, und ob ein Schnitt nach Gruppen nötig ist (§4).

**Risiken:**

- Eine oder mehrere Gegenproben liegen über 600 s (`lint-gegenprobe`, `stop-hook-gegenprobe`, `gegenproben-nach-einlieferung-gegenprobe`, `container-start-gegenprobe` sind die größten nach Zeilen) und brauchen eigene Slices; bis diese in `done/` liegen, ruht `CO-001` weiter — **Ausgang:** offen bis Closure (eingetreten: Kennungen der Slices `slice-harness-gegenproben-schneller-<name>` und `CO-001` bleibt Aktiv; entfallen: jede Gegenprobe liegt unter 600 s). **Für `container-start-gegenprobe` ist der Fall schon gemessen** (Sendung aus `slice-harness-container-start-u-boot-0-13`, §1): 19 min 38 s für den vollen Lauf am Host; der Ausgang ist damit *eingetreten*, sobald die eigene Messung ihn bestätigt, und die Tabelle von `SPEC-060` Punkt 16 nennt dann die Kennung des schneidenden Slice statt die dieses.
- Die Zahl der aufzunehmenden Gegenproben macht den Slice größer als die Marke (§4) — **Ausgang:** offen bis Closure (eingetreten: Schnitt nach Gruppen mit Kennung; entfallen: Messung der Diff-Summe am Endstand unter 1000).
- Die Messung hängt an der Last des Hosts (`BEO-REPO/lastabhaengige-zaehlung`) und gleichzeitige Läufe stören einander (`BEO-REPO/gegenprobe-temp-praefix-geteilt`) — **Ausgang:** offen bis Closure (weiter offen, wenn eine Zahl der Tabelle an der Last hängt: `BEO-REPO/lastabhaengige-zaehlung`).
- `CO-001` wird aufgelöst, ohne dass der volle Lauf an einem getaggten Commit je unter zehn Minuten belegt ist — **Ausgang:** offen bis Closure (entfallen, wenn die Tabelle jede Gegenprobe `aufgenommen` führt und der Trigger von `CO-001` in der Closure Zeile für Zeile belegt ist).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen. Treffer: `gate-laufzeit-waechst-mit-gegenprobe` (Stand verkörpert; die Laufzeit der Gegenproben außerhalb von `make gates` ist hier der Gegenstand), `gegenprobe-haengt-an-umgebung-des-hosts` (die Messung hängt am Host), `lastabhaengige-zaehlung` (Zahl der Tabelle an der Last, §6), `gegenprobe-temp-praefix-geteilt` (Messung am Host, ein Lauf zur Zeit), `test-faellt-unter-cpu-last-aus`. Der Stand der Zähler ist beim Start zu lesen; erreicht einer mit diesem Slice drei, braucht er einen Folge-Slice.

**Zählung** (Planner, 2026-10-21, `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 3** (1: Messung der 18 Gegenproben; 2: Aufnahme der Gegenproben mit höchstens 600 s; 3: Schnitt der übrigen und Stand von `CO-001`). **Schichten: eine** — Harness-Werkzeug (Skripte, Fragmente, Zeilen in `harness/README.md`); der Vertrag liegt beim Architect und zählt nicht, die Planung (Slices, Carveout) zählt nicht. Mit den Sendungen aus `CO-001` und `slice-harness-gegenproben-schneller` nachgezählt: Die Sendung aus `CO-001` liegt im dritten Liefer-Punkt (Stand des Carveouts in der Closure), die aus `slice-harness-gegenproben-schneller` in den Punkten 1 und 2; kein vierter Punkt.

**Nachzählung mit der Sendung aus `slice-harness-container-start-u-boot-0-13`** (Planner, 2026-10-21, im selben Commit, der sie in §1 und §6 einträgt; `AGENTS.md` §3.13): **Liefer-Punkte bleiben 3, Schichten bleibt 1.** Die Sendung trägt eine Zahl und eine Version in die vorhandenen Punkte 1 (Messung: `container-start-gegenprobe` steht dort schon, die Vormessung ist eine Eingabe, kein neuer Gegenstand) und 3 (Schnitt: für diese Gegenprobe ist der Ausgang damit vorgezeichnet, die Zahl der zu schneidenden Gegenproben wächst nicht); sie bringt kein Skript, keinen Vertragstext und keine Datei außerhalb der Planung. Damit bleibt der Slice unter drei Liefer-Punkten und unter zwei Schichten. Der vorgemerkte Schnitt nach Gruppen (§4) bleibt davon unberührt: Er hängt an der Zahl der **aufzunehmenden** Gegenproben, und `container-start-gegenprobe` wird nach dieser Messung nicht aufgenommen. **Größe über der Marke nicht ausgeschlossen:** Die Zahl der aufzunehmenden Gegenproben steht erst nach der Messung fest; der vorgemerkte Schnitt steht in §4. Größe: Marke und Schätzung in §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Verträge und Konventionen des Harness gehen den Skripten voraus).

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der die Ortsangabe aus `slice-harness-gegenproben-umzug-dev` einträgt): Die Gegenproben liegen seit jenem Slice unter `tools/dev/<bereich>/`, die Pfade der Gegenproben in diesem Plan sind nachgezogen (`tools/dev/` ist wie `tools/harness/` Harness-Werkzeug, keine neue Schicht). Ortsangabe der Gegenprobe geändert, Liefer-Punkte (3) und Schichten (1) unverändert. Die Zeile in §3 nennt als Glob nur `tools/dev/harness/`; die Gegenproben unter `tools/dev/test/`, `tools/dev/arch/` und `tools/dev/hook/` gehören zur Messung ebenso.

**Nachzählung der Sendung aus `CO-003`** (Planner, 2026-10-24, `AGENTS.md` §3.13, im Commit, der sie einträgt): Die Sendung liegt im dritten Liefer-Punkt (Stand der Carveouts in der Closure), wie die aus `CO-001`; kein Skript, kein Vertragstext, keine Datei außerhalb der Planung. **Liefer-Punkte bleiben 3, Schichten 1.** Größe über der Marke bleibt wie oben vorgemerkt (Schnitt nach Gruppen, §4).
