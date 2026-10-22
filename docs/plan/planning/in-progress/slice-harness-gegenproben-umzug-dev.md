# Slice slice-harness-gegenproben-umzug-dev: Die Gegenproben ziehen nach `tools/dev/` um, das Skript unter Test bleibt, wo es läuft

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Entscheidung des Planners, 2026-10-22). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist. `welle-harness-regeln` und `welle-gegenproben-und-mutanten` führen ihn nicht in ihrem §4, und ihre Closure-Trigger hängen nicht an ihm; ihn aufzunehmen änderte deren Closure (ein Slice mehr in `done/`), ohne dass eine Bedingung dazukäme. Angelegt auf Anweisung des Koordinators nach der Entscheidung des Nutzers vom 2026-10-22 (§1) und direkt in `next/`: Der Nutzer hat den Start entschieden, `Verantwortlich:` ist gesetzt.

**Bezug:** kein Lastenheft-Bezug — der Slice ändert den Ort von Werkzeugen des Harness, keinen Anspruch des Produkts. Entscheidungen: Entscheidungen des Nutzers vom 2026-10-22 (§1, wörtlich); [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (die Gegenproben stehen nicht in `GATE_CHECKS`; hier unverändert). Bindung an das Register (§8): `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`, `BEO-REPO/slice-waechst-durch-uebernahmen`, `BEO-REPO/gegenprobe-variable-kollidiert-mit-skript`.

**Berührte Spec-Stellen:** `SPEC-060` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-060-gegenproben.md); Punkte 1 bis 3 (Name, Listen, Zuordnung) und 16 (Kopf jeder Gegenprobe) bleiben unverändert; Punkt 17 trägt einen Ort) · `SPEC-048` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-048-abdeckung.md); Punkt 7 liest `-gegenprobe.sh` in jeder Tiefe unter `tools/`, unverändert) · nur als Ortsangabe: `SPEC-057` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-057-integration-aufraeumen.md)), `SPEC-061` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md)), `SPEC-062` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md)), `SPEC-064` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-064-ci.md)), `SPEC-067` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-067-release.md)), `SPEC-070` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md)). Die Ortsangaben schreibt der Architect vor dem ersten Code-Commit (§3).

**Verantwortlich:** pt9912

**Autor:** Planner. **Datum:** 2026-10-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die 21 Gegenproben (die 16 `*-gegenprobe.sh` unter `tools/harness/`, die 3 unter `tools/test/`, `tools/hook/commit-msg-gegenprobe.sh` und `tools/arch/a-check-negativ.sh`; zusammen 14 702 Zeilen, gezählt mit `wc -l`) liegen nach einem reinen `git mv` unter `tools/dev/harness/`, `tools/dev/test/`, `tools/dev/arch/` und `tools/dev/hook/`; jede Zeile im Repo, die ihren Ort nennt und stimmen soll, nennt den neuen, und jede der 21 läuft am neuen Ort im ungemutierten Fall grün. Name, Make-Ziel, Liste (`GEGENPROBEN_CONTAINER` oder `GEGENPROBEN_HOST`), Fälle und Mutationen bleiben, und das Skript unter Test bleibt, wo es läuft.

**Entscheidungen des Nutzers vom 2026-10-22** (wörtlich, so vom Koordinator übermittelt):

1. „Man müsste die Entwicklung der Skripte (incl. Gegenproben Fixtures) trennen vom Nutzen/Ausführen“
2. „Wir können die Entwicklung erst einmal hier in einem eigenen Bereich lassen“
3. „Dev Ordner: tools/dev/harness, tools/dev/test, ...“
4. „erst tools/dev/... anlegen und dann die script kopieren/gegeproben umziehen, pinnen (sha, ...)“
5. „Aber für das Pinnen braucht man kein eigenes Repo“
6. „die scripte halten wir jetzt doppelt, einmal zum Nutzen im aktuellen Ordner und den dev-Ordner zum Entwickeln“

Anlass war die Aussage des Nutzers „Wir können so nicht weitermachen“: seit dem 2026-10-18 vier Commits am Produkt gegen 169 an Harness und Tools (Angabe des Koordinators, vom Planner nicht nachgezählt). **Reihenfolge nach (4):** Schritt 1 ist dieser Slice (`tools/dev/` anlegen, Gegenproben umziehen), Schritt 2 die Kopie der Skripte als Quelle und der Pin mit SHA (`slice-harness-dev-kopie-und-pin`, `open/`). **Folge aus (6):** Die Skripte liegen künftig zweimal vor, `tools/harness/…` zum Nutzen (gepinnt) und `tools/dev/…` zum Entwickeln. Dieser Slice legt keine zweite Kopie an; deshalb prüfen die Gegenproben nach dem Umzug weiter das Skript unter `tools/harness/` und `tools/test/`, und kein Verweis einer Gegenprobe auf ihr Skript unter Test ändert sich (Anweisung des Koordinators nach (6)). Geändert wird, wo die Gegenprobe selbst liegt und wie sie von dort ihre Wurzel und das Skript unter Test findet.

**Regel des Umzugs:** Die Gegenprobe zieht um, das Skript unter Test bleibt. Daraus folgt für die drei Hilfsskripte unter `tools/harness/` (Entscheidung des Planners, begründet): `gegenproben.sh` (Rezept von `make gegenproben` und `make gegenproben-host`, `SPEC-060` Punkt 4) und `gegenproben-nach-einlieferung.sh` (Inhalt des `post-receive`-Hooks im Bare-Repo des Gegenproben-Containers) **laufen** und bleiben Nutzer-Seite; `gegenproben-regel.sh` wird allein von `gates-schnell-gegenprobe.sh` geladen und von ihr mit Fällen und Mutationen geprüft, ist also Skript unter Test wie die beiden anderen und bleibt aus demselben Grund. Ob sie, wie jedes Skript unter Test, in Schritt 2 eine dev-Kopie bekommen, entscheidet jener Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Das Kopieren der Skripte nach `tools/dev/` als Quelle und das Pinnen mit SHA (Schritt 2)** — Klasse 1, Adresse `slice-harness-dev-kopie-und-pin` (`open/`; der Plan nimmt die Sendung in §1 und DoD an und führt die offenen Fragen (a) bis (g) und (i)). Grund: Wer Quelle ist, wie freigegeben wird und wie ein Gate die Summe prüft, hat der Nutzer nicht entschieden; „die scripte halten wir jetzt doppelt“ (6) beschreibt das Ziel, nicht den Weg. Bis dahin gibt es keine zweite Kopie, und ein Umzug, der schon eine anlegte, entschiede diese Fragen im Code.
- **Verhalten der Gegenproben** (Fälle, Mutationen, Fall-Listen, Auswahl-Variablen, Obergrenze nach `SPEC-060` Punkt 16, Meldungen) — Klasse 3, anderer Vorgang. Grund: Ein Umzug, der Verhalten ändert, ist nicht mehr mechanisch prüfbar; die Prüfung der Wirkung (DoD 3) vergleicht gegen den Zustand davor.
- **Die Skripte der Nutzer-Seite** (`tools/harness/*` und `tools/test/*` außer den 21; darunter `gegenproben.sh`, `gegenproben-nach-einlieferung.sh`, `gegenproben-regel.sh`) — Klasse 2, Bestand bleibt bewusst stehen. Grund: Sie laufen oder sind Skript unter Test (Regel des Umzugs); nur ihr Kommentar „GEPRÜFT DURCH“ folgt dem Ort (17 Zeilen in 16 Dateien, DoD 2).
- **Eingefrorene und chronische Texte, die einen alten Ort nennen** (für die 21 Pfade gezählt): `docs/reviews/` (67 Dateien, 212 Zeilen), `docs/plan/planning/done/` (9 Dateien, 39 Zeilen; Volltexte und Stubs geschlossener Slices), [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (1 Zeile; angenommene ADR, `AGENTS.md` §3.5), Belege und Stand im Beobachtungs-Register (`observations/`, 5 Dateien, 6 Zeilen), das Drift-Log der Roadmap (2 Zeilen; es trägt Umplanungen, `AGENTS.md` §3.7) und `docs/maintainer/freigaben/0.1.1.md` (1 Zeile; ein `git diff --stat` mit Hashes, nur an jenen Commits reproduzierbar) — Klasse 2. Grund: Ein Zeitdokument nennt den Ort zu seiner Zeit, `git log --follow` verbindet ihn mit dem neuen, und ein umgeschriebener Beleg behauptete einen Zustand, den er nicht beschrieb.
- **Die Pläne anderer Slices in `open/` und `next/`** (18 Slice-Pläne und der Welle-Plan `welle-gegenproben-und-mutanten`, Liste in §3) — Klasse 4 (Schicht: Planung), mit Sendung. Grund: Der Implementer ändert fremde Pläne nicht; die Ortsangabe trägt der Planner bei der Closure in jeden Nehmer ein, mit Nachzählung in dessen §8 (`AGENTS.md` §3.13, *Auch eine Ortsangabe wird nachgezählt*), weil sie erst dann stimmt.
- **Produkt-Code** (`cmd/`, `internal/`, `test/`) und das `Dockerfile` außer einer Kommentarzeile — Klasse 4. Grund: Der Slice berührt allein Werkzeug-Skripte, Fragmente und Kommentare.
- **Ein neuer Sensor oder eine neue Regel für den Ort** (etwa ein Gate, das eine Gegenprobe unter `tools/harness/` ablehnt) — Klasse 3. Grund: Der Nutzer hat den Umzug verlangt, kein Gate (Anweisung des Koordinators: keine neue Regel, kein neuer Sensor); der eine Satz in `harness/README.md` (DoD 2) beschreibt den Ist-Zustand und prüft nichts (`AGENTS.md` §3.11).
- **Die Frage, welche Gegenproben generisch sind** — Klasse 3, anderer Vorgang, keine Adresse. Grund: Der Nutzer hat sie nicht gestellt; sie ist eine Einordnung der Gegenproben nach Wiederverwendbarkeit und braucht keinen Umzug.

**Sendung an die Nehmer (`AGENTS.md` §3.13):** Die Ortsangabe „Gegenproben liegen unter `tools/dev/<bereich>/`“ geht bei der Closure dieses Slice an die 18 Slice-Pläne und den Welle-Plan der Liste in §3 (sie nennen eine vorhandene oder geplante Gegenprobe unter `tools/harness/`, `tools/test/`, `tools/hook/` oder `tools/arch/`). Der Eintrag steht nicht jetzt, weil er vor der Lieferung falsch wäre; **Reihenfolge-Empfehlung des Planners:** Dieser Slice liegt vor `slice-harness-gegenproben-schneller` und vor jedem Slice, der eine Gegenprobe anlegt (§6, Risiko 4).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 — Umzug:** Die 21 Dateien liegen unter `tools/dev/harness/` (16), `tools/dev/test/` (3: `abdeckung-gegenprobe.sh`, `integration-volumes-gegenprobe.sh`, `postgres-versionen-gegenprobe.sh`), `tools/dev/arch/` (`a-check-negativ.sh`) und `tools/dev/hook/` (`commit-msg-gegenprobe.sh`), in **einem Commit, der nur `git mv` enthält** (Beleg in §7: `git show --name-status --find-renames=100% <Commit>` nennt 21 Zeilen `R100` und sonst nichts). `find tools -name '*gegenprobe*.sh' -o -name 'a-check-negativ.sh'` nennt danach außerhalb von `tools/dev/` nur `gegenproben.sh`, `gegenproben-nach-einlieferung.sh` und `gegenproben-regel.sh`.
- [ ] **Liefer-Punkt 2 — Pfadnachzug** (mehrere Commits hinter Liefer-Punkt 1 sind erlaubt, Verfahren in §3, Randformen in §6): (a) in den 24 Fragmenten unter `harness/mk/` die 21 Rezeptzeilen und die 9 Kommentarzeilen „GEPRUEFT DURCH“; (b) in jeder der 21 Gegenproben die Wurzel (drei statt zwei Ebenen), der Verweis auf das Skript unter Test (über die Wurzel, nicht über das Nachbarverzeichnis), Selbstkopien und Attrappen am neuen Ort und die Kommentare; (c) die 17 Zeilen „GEPRÜFT DURCH“ in den 16 Skripten unter Test, eine Kommentarzeile im `Dockerfile` (Zeile 61), eine in `.claude/hooks/stop-require-gates.sh` (Zeile 9, danach `bash -n` ohne Befund) und eine in `harness/conventions.md` (Zeile 113); (d) die Tabellen unter `docs/user/` mit `make abdeckung` neu geschrieben (nicht von Hand); (e) ein Satz in `harness/README.md` §Sensors, der den Ort der Gegenproben nennt (Ist-Zustand, prüft nichts); (f) die Ortsangaben in den Verträgen (7 Dateien, 8 Zeilen, §3; schreibt der Architect vor dem ersten Code-Commit). **Beleg in §7:** `git grep -n -E 'tools/(harness|test|hook|arch)/[a-z0-9-]*-gegenprobe\.sh|tools/arch/a-check-negativ\.sh' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/adr' ':!docs/maintainer' ':!.harness'` nennt keine Datei unter `harness/`, `docs/user/` und `.claude/` und kein `Dockerfile`, unter `tools/` allein die erfundene Fixture `tools/test/x-gegenprobe.sh` in `abdeckung-gegenprobe.sh` (Zeile 434, ein Pfad im Temp-Baum; bleibt, §6 (h)); was sonst bleibt, sind die Pläne der Nehmer-Liste (§3) und das Drift-Log der Roadmap.
- [ ] **Liefer-Punkt 3 — Wirkung belegt:** In §7 steht eine Tabelle mit 21 Zeilen (Gegenprobe · Aufrufzeile · Ausgang · Zeit): Jede läuft am neuen Ort im ungemutierten Fall grün. Hat das Skript eine Auswahl (`GP_MUT`, `GP_NUR`; der Kopf nennt sie, §3 Schritt 7), läuft es mit der engsten Auswahl, die den ungemutierten Fall und die Stellen trifft, die Liefer-Punkt 2 berührt hat; `release-gegenprobe` und `test-container-gegenprobe` laufen **nur so, nie voll** (Zustand `nicht gemessen` nach `SPEC-060` Punkt 16, `AGENTS.md` §6 Punkt 6); Gegenproben ohne Auswahl laufen voll. `container-start-gegenprobe` läuft nur am Host (Host-Ziel) und steht in §5.
- [ ] `make gates-schnell-gegenprobe`, `make gegenproben-gegenprobe`, `make abdeckung-check`, `make docs-check` und `make kopf-check` grün, danach `make gates` grün als letzter Lauf (der Move ändert die Dateiliste, `tools/` steht nicht auf der Positivliste des Stop-Hooks).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).
- [ ] **Sendung an die Nehmer (Closure, Planner, `AGENTS.md` §3.13):** Die Ortsangabe steht in jedem Plan der Liste in §3 (18 Slice-Pläne, der Welle-Plan), mit Nachzählung in dessen §8 im selben Commit; `grep` nach `slice-harness-gegenproben-umzug-dev` findet jeden. Die Adresse für Schritt 2, `slice-harness-dev-kopie-und-pin`, nimmt die Sendung in §1 und DoD an (angelegt im Commit mit diesem Plan).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| die 21 Gegenproben: `tools/harness/` (16), `tools/test/` (3), `tools/hook/` (1), `tools/arch/` (1) | `git mv` nach `tools/dev/<bereich>/`, rein | Liefer-Punkt 1; ein Commit ohne Inhaltsänderung (`AGENTS.md` §3.3) |
| dieselben 21 am neuen Ort | update | Liefer-Punkt 2 (b): Wurzel 2 → 3 Ebenen in 16 (Vorgabe von `GP_QUELLE` in `container-einrichten-gegenprobe.sh` und `container-start-gegenprobe.sh` eingeschlossen); Skript unter Test über die Wurzel statt über das Nachbarverzeichnis in 8; Selbstkopie, Attrappe und Rezepttext, die die Gegenprobe selbst am Ort nennen, in 10 (28 Zeilen) |
| `harness/mk/` (24 Fragmente) | update | 21 Rezeptzeilen `@bash tools/…-gegenprobe.sh` und 9 Kommentarzeilen „GEPRUEFT DURCH“; Name des Ziels, Liste und Rezept sonst gleich |
| 16 Skripte unter Test (`tools/harness/`: `archiv-bezug.sh`, `ci-range.sh`, `container-einrichten.sh`, `container-start.sh`, `coverage.sh`, `dateiliste.sh`, `gegenproben.sh`, `gegenproben-nach-einlieferung.sh`, `gegenproben-regel.sh` (2 Zeilen), `kopf-check.sh`, `lint.sh`, `meldungskatalog-check.sh`, `release.sh`, `workflows-check.sh`; `tools/test/`: `abdeckung.sh`, `integration-aufraeumen.sh`) | update (Kommentar) | 17 Zeilen „GEPRÜFT DURCH <Gegenprobe>“ nennen den neuen Ort; sonst nichts (Skript unter Test bleibt, §1) |
| `Dockerfile` (Zeile 61), `.claude/hooks/stop-require-gates.sh` (Zeile 9), `harness/conventions.md` (Zeile 113) | update (Kommentar, Verweis) | je ein Ort; der Hook-Kommentar ändert kein Verhalten, `bash -n` belegt die Syntax |
| `docs/user/abdeckung-gate.md` (3 Zeilen) und die übrigen vier Tabellen | `make abdeckung` | `abdeckung.sh` liest `*-gegenprobe.sh` in jeder Tiefe unter `tools/` (Zeile 148, `SPEC-048` Punkt 7) und bleibt unverändert; die Tabelle nennt den gelesenen Pfad |
| `harness/README.md` §Sensors | update | ein Satz: Die Gegenproben liegen unter `tools/dev/{harness,test,arch,hook}/`, eine neue legt der Schreiber dort an; das Skript unter Test liegt unter `tools/` (kein Verweis dort wird falsch, §1) |
| `harness/vertraege/`: `SPEC-057` (Z. 101), `SPEC-060` (Z. 291), `SPEC-061` (Z. 19), `SPEC-062` (Z. 36), `SPEC-064` (Z. 12, 14), `SPEC-067` (Z. 22), `SPEC-070` (Z. 307) | durch den Architect, vor dem ersten Code-Commit | nur der Ort, kein Inhalt; `SPEC-062` und `SPEC-070` nennen Gegenproben, die es noch nicht gibt (Ort ihrer künftigen Datei); `SPEC-060` Punkt 17 nennt `bash tools/harness/<name>.sh` für ein Gegenprobe-Skript **und** das Skript der Sammelziele, den Wortlaut der beiden Orte wählt der Architect |
| `.d-check.yml`, `.golangci.yml`, `.a-check.yml`, `.dockerignore`, `Makefile`, `.github/workflows/`, `.githooks/` | unverändert | gelesen am Stand `1ad628a7`: keine Datei nennt eine der 21 Gegenproben; `.dockerignore` ist eine Allowlist (`lint.sh`, `coverage.sh`), `tools/dev/` bleibt wie `tools/harness/` im Übrigen draußen |

**Nehmer der Ortsangabe (Closure, Planner; Befund des Planners vom 2026-10-22, bei der Closure mit demselben `git grep` neu zu ziehen):** `next/`: `slice-harness-gegenprobe-matrix-p7-flanke`, `slice-harness-gegenproben-schneller`, `slice-harness-slice-mv-exakter-treffer`; `open/`: `slice-harness-commit-struktur-id`, `slice-harness-devcontainer-einrichten-zugang`, `slice-harness-gegenprobe-baumkopie`, `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`, `slice-harness-gegenproben-grenze-rest`, `slice-harness-kommentare-beleg-ort-und-traeger-adr`, `slice-harness-lint-warten-ohne-frist`, `slice-harness-mutation`, `slice-harness-nachweis-uebernehmen`, `slice-harness-nachweis-uebernehmen-aufrufarten`, `slice-harness-regeln-implementer-uebergabe`, `slice-harness-rollenlauf-wiederholungsgrenze`, `slice-harness-rollenlauf-zeitgrenze`, `slice-harness-rollenlauf-zeitgrenze-haertung`, `slice-test-unter-last-ursache` (18); dazu `welle-gegenproben-und-mutanten` (2 Zeilen). Befund-Befehl: `git grep -l -E 'tools/(harness|test|hook|arch)/[A-Za-z0-9<>*_-]*(-gegenprobe|gegenprobe>)\.sh|tools/arch/a-check-negativ\.sh' -- docs/plan/planning/open docs/plan/planning/next docs/plan/planning/welle-*.md` (ohne diesen Plan und `slice-harness-dev-kopie-und-pin`, die den Ort als Pfad nennen). Nicht Nehmer: die Roadmap (zwei Zeilen im Drift-Log, Chronik, §1).

**Verfahren des Pfadnachzugs (Liefer-Punkt 2; der Implementer greppt selbst, die Zeilen oben sind der Befund des Planners am Stand `1ad628a7`, nicht abschließend):**

1. **Kein Ersetzen im Block.** Ein `sed` über alle 21 Pfade ist falsch: Die Probe des Planners in einem Wegwerf-Klon (Moves plus blankes Ersetzen der 21 Pfade in `tools harness Dockerfile .claude docs/user`) ergibt 58 Dateien und 170 Zeilen, setzt in `lint-gegenprobe.sh` aber das Ziel eines `cp` auf den neuen Ort und lässt das zugehörige `mkdir -p` auf dem alten (Zeilen 340 und 341). Jede Fundstelle wird gelesen und nach drei Klassen entschieden: (a) die Gegenprobe selbst (Attrappe, Selbstkopie, Rezepttext, Kommentar) → neuer Ort; (b) Skript unter Test oder anderes Werkzeug → bleibt; (c) Text, den eine Mutation in einem Fragment oder einem Skript unter Test ersetzt → folgt der Datei, deren Text ersetzt wird (steht dort der neue Ort, steht er in der Mutation; steht im Skript unter Test `…/../..`, bleibt es, weil es nicht umzieht).
2. **Wurzel:** drei Ebenen (`../../..`), bestimmt aus `BASH_SOURCE`, nicht aus dem Arbeitsverzeichnis, wie heute. `commit-msg-gegenprobe.sh` und `a-check-negativ.sh` bestimmen keine Wurzel (sie laufen im Arbeitsverzeichnis, das `make` setzt) und bleiben so.
3. **Skript unter Test:** `$wurzel/tools/<bereich>/<name>.sh`. Wo ein Skript bisher nur `hier` kennt (`kopf-check-`, `meldungskatalog-check-`, `gegenproben-nach-einlieferung-`, `integration-volumes-`, `abdeckung-gegenprobe.sh`), führt eine neue Variable die Wurzel mit dem Präfix `gp_` (`gp_wurzel`); vorher `grep -n gp_wurzel <Skript>` (`BEO-REPO/gegenprobe-variable-kollidiert-mit-skript`).
4. **Selbstkopie und Attrappe im Temp-Baum** liegen am neuen Ort (`tools/dev/<bereich>/`), mit dem zugehörigen `mkdir -p`; sonst bestimmt die Kopie ihre Wurzel falsch.
5. **Fragmente und Kommentare:** Rezept und „GEPRUEFT DURCH“ mit dem neuen Ort; jede Mutation, die Fragmenttext ersetzt (`mut … $'\t@bash tools/…'`), folgt dem Rezept.
6. **Dann** `make abdeckung`, die Verträge sind vom Architect geschrieben, `harness/README.md` bekommt den einen Satz.
7. **Wirkung (Liefer-Punkt 3):** Der Implementer liest den Kopf jeder Gegenprobe (Zeile „Zum Arbeiten“), nimmt die engste Auswahl und fährt sie mit `GP_ENTDECKE=1`, wo es das Skript kennt. Befund des Planners, Zeilen im Kopf: `GP_MUT=<Muster>` fährt nur Läufe, deren Name passt, die unveränderte Kopie nur, wenn das Muster `ok` trifft, in `archiv-bezug-` (Z. 5), `ci-` (Z. 14), `container-einrichten-` (Z. 28), `gegenproben-gegenprobe.sh` (Z. 11), `gegenproben-nach-einlieferung-` (Z. 7), `range-konfiguration-` (Z. 15), `release-` (Z. 16) und `workflows-check-gegenprobe.sh` (Z. 9); in `container-gegenprobe.sh` (Z. 21) und `integration-volumes-gegenprobe.sh` (Z. 24) fährt `GP_MUT` nur Mutationen, die unveränderte Kopie läuft immer; `container-einrichten-` und `container-start-gegenprobe.sh` kennen zusätzlich `GP_NUR=<Muster>` für die Fälle; `gates-schnell-gegenprobe.sh` kennt nur `GP_NUR_ARBEITSBAUM=1`, `lint-gegenprobe.sh` nur `LINT_GEGENPROBE_PARALLEL`; keine Auswahl haben `coverage-`, `kopf-check-`, `meldungskatalog-check-`, `stop-hook-`, `postgres-versionen-`, `abdeckung-gegenprobe.sh`, `commit-msg-gegenprobe.sh` und `a-check-negativ.sh` (sie laufen voll; Zeiten aus dem Register, ruhiger Host: `coverage` 89 bis 106 s, `lint` 150 bis 284 s, `stop-hook` 129 bis 164 s, `postgres-versionen` 137 s, `a-check-negativ` 11 bis 13 s, `abdeckung` 5 s, `hook` 1 s, `gates-schnell` 17 bis 19 s, `integration-volumes` 42 bis 67 s, `gegenproben-gegenprobe` 42 s; `kopf-check` und `meldungskatalog-check` stehen dort nicht einzeln, der Implementer nennt die Zeit; die Summe der Gegenproben ohne Auswahl liegt nach diesen Zahlen bei rund 13 Minuten). Die Stellen, an denen eine Gegenprobe ihren Ort selbst nennt, muss die Auswahl treffen: bei `ci-` und `range-konfiguration-gegenprobe.sh` die Läufe `m-p18*` (Zeilen 501 bis 534 und 278 bis 314, Kopie des Skripts), bei `release-gegenprobe.sh` die Läufe `m-p13*` (Zeilen 1981 bis 2014); der Implementer liest die Namen dort, er rät sie nicht. `container-start-gegenprobe.sh` ruft `docker`, `u-boot` und `devcontainer` des Hosts und läuft im Dev-Container nicht (`SPEC-060` Punkt 3); der Kopf nennt eine Auswahl der Gruppe Attrappen, die am Host zu fahren ist (§5).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Nutzer hat den Start entschieden (Entscheidungen vom 2026-10-22, §1). Beim Anlegen liegt kein Slice in `in-progress/` (WIP-Limit 1 je Rolleninhaber); der Architect hat die Ortsangaben in den Verträgen (§3) geschrieben und die Randformen (§6) geprüft, bevor der erste Code-Commit entsteht. **Basis der Metrik:** — (den Hash des Commits, der diese Datei nach `in-progress/` bewegt, trägt der Planner beim Übergang hier ein, `AGENTS.md` §3.14).

**Größe** (Metrik und Befehl: `AGENTS.md` §3.14, hier nicht wiederholt). Pfade des Plans: `tools harness Dockerfile .claude docs/user`. Ein reiner Move zählt mit `--find-renames=10%` null. **Faktor Tests zu Produktcode: keine Messung.** Der Slice liefert weder Produktcode noch neue Tests; die geänderten Zeilen sind Pfadangaben in Skripten, die selbst Gegenproben sind (vgl. `slice-harness-gegenproben-schneller` §3, dort ebenfalls ohne Faktor); gerechnet wird mit Faktor 1 auf die Zeilen. **Schätzung des Planners:** Mitte 320, Spanne 250 bis 400. Beleg: Eine Probe in einem Wegwerf-Klon (21 Moves und blankes Ersetzen der 21 Pfade) ergibt 170 Zeilen in 58 Dateien (Fragmente 60, Skripte unter Test 34, übrige Kommentare und Tabellen rund 20, die Gegenproben selbst rund 55); von Hand kommen hinzu: Wurzel 2 → 3 Ebenen in 16 Gegenproben (rund 40), Skript unter Test über die Wurzel in 8 (rund 30), `mkdir`-Zeilen und Selbstkopien (rund 30), der Satz in `harness/README.md` (rund 4). **Marke: 600. Meldeschwelle: 450.** Die Mitte liegt bei 53 % der Marke, unter den 70 % des Schnitts vor dem Anspruch. Meldet der Implementer die Meldeschwelle, trägt der Planner `weiter` oder `schneiden` mit dem gemessenen Stand als datierte Zeile in §7 ein (`AGENTS.md` §3.14).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt zu irgendeiner Zeit am committeten Stand **über 600**. Die Naht für den Schnitt ist der Bereich: erst `tools/dev/test/`, `tools/dev/arch/` und `tools/dev/hook/` (5 Gegenproben), danach `tools/dev/harness/` (16); die Naht ersetzt keinen Schnitt beim Anlegen (die Schätzung liegt unter 70 %).
- `in-progress` → `open` (blockiert — Carveout?): Der ungemutierte Fall einer Gegenprobe ist rot, und die Ursache ist nicht der Pfad (Vergleich: derselbe Fall auf `1ad628a7`, dem Stand vor dem Umzug). Der Implementer repariert sie nicht (§1, Verhalten); der Planner entscheidet zwischen einem Carveout mit Folge-Slice (Modul 7) und der Rückführung. Ebenso, wenn der Nutzer die Entscheidungen aus §1 ändert.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben. **Dazu:** (1) `container-start-gegenprobe` ist am Host gefahren (Handlung des Nutzers oder des Koordinators; Host-Ziel, im Dev-Container nicht belegbar, `SPEC-060` Punkt 3): die Aufrufzeile, die Auswahl (aus dem Kopf des Skripts, Gruppe Attrappen) und der Ausgang stehen in §7; fehlt der Host, steht dort „nicht belegt“ mit Datum, und der Nutzer entscheidet über Verzicht oder Aufschub, kein stilles Grün. (2) Die Sendung an die Nehmer ist eingetragen (DoD, Planner). **Der Verifier** fährt die Wirkung nach wie der Implementer, `release-gegenprobe` und `test-container-gegenprobe` mit der engsten Auswahl, nicht voll (Anweisung des Koordinators; der Wortlaut von `AGENTS.md` §6 Punkt 6, der Verifier fahre den vollen Lauf einmal, widerspricht dem; den Widerspruch löst der Nutzer).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12). Der Slice liefert keinen neuen Vertrag und sagt nichts Neues zu; die Randformen des Umzugs stehen trotzdem hier, damit der Implementer keine im Code entscheidet. Die Spalte *Vorschlag* ist der Vorschlag des Planners; der Architect prüft ihn vor dem ersten Code-Commit, und was er als Zusage haben will, trägt er in `SPEC-060` ein (Empfehlung des Planners: nur die Ortsangaben aus §3). Eine nicht genannte oder offene Randform entscheidet der Implementer nicht, er gibt sie zurück.

| Randform | Frage | Vorschlag |
|---|---|---|
| (a) Pfadbestimmung der Gegenprobe | Tiefe zwei auf drei; Aufruf aus anderem Arbeitsverzeichnis oder über einen Symlink; die Vorgabe `GP_QUELLE` in `container-einrichten-` und `container-start-gegenprobe.sh` | Wurzel aus `BASH_SOURCE`, drei Ebenen, wie heute; `GP_QUELLE` bleibt die Überschreibung, nur die Vorgabe zählt drei Ebenen; `commit-msg-gegenprobe.sh` und `a-check-negativ.sh` bestimmen keine Wurzel und bleiben an das Arbeitsverzeichnis gebunden, das `make` setzt; kein neuer Fall (dieselbe Zusage) |
| (b) Skript unter Test | Ein Verweis `$hier/<skript>.sh` zeigt nach dem Umzug auf eine Datei, die es dort nicht gibt; ein Fall, der eine fehlende Datei überspringt, bliebe grün | `$wurzel/tools/<bereich>/<skript>.sh`; der Implementer prüft je Stelle, dass die Datei gelesen wird (kein stilles Überspringen, kein `cmp` zweier leerer Eingaben) |
| (c) Selbstkopie und Attrappe im Temp-Baum | `ci-`, `range-konfiguration-`, `release-gegenprobe.sh` (Mutation der eigenen Kopie), `gates-schnell-` und `lint-gegenprobe.sh` (Kopie läuft im Temp-Baum), `container-` und `container-start-gegenprobe.sh` (Attrappe des Ziels `…-gegenprobe`), `kopf-check-`, `meldungskatalog-check-` und `stop-hook-gegenprobe.sh` (Attrappe) | die Kopie liegt im Temp-Baum am neuen Ort `tools/dev/<bereich>/` mit dem zugehörigen `mkdir -p`, damit ihre Wurzelbestimmung dieselbe ist wie die des Originals |
| (d) Mutation mit Pfad als Text | Eine Mutation ersetzt Text in einem Fragment oder Skript unter Test, der einen Ort nennt (`$'\t@bash tools/…'`, `wurzel=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." …)` in `container-start.sh`, `"$hier/../.."` in `workflows-check.sh`) | der Text folgt der Datei, deren Text ersetzt wird (§3 Schritt 1); die Mutationen auf `…/../..` im Skript unter Test bleiben, weil es nicht umzieht; trifft eine Mutation ihre Textstelle nicht mehr, ist das ein Befund (`ci-gegenprobe.sh` meldet „Textstelle nicht gefunden“; wo ein `sed -i` ohne Prüfung steht, etwa `gates-schnell-gegenprobe.sh` Zeile 259, prüft der Implementer von Hand) |
| (e) Name, Liste und Zuordnung | `SPEC-060` Punkte 1 bis 3 und 16 | unverändert: Ziel-Name, Liste (`GEGENPROBEN_CONTAINER` oder `GEGENPROBEN_HOST`) und Kopf jeder Gegenprobe bleiben; die Namensregel liest Ziele und Listen aus `make`, nicht Pfade; die Fälle `p2-*` laufen unverändert |
| (f) Aufrufarten von `make` | `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` aus der Umgebung, auch in Kombination (`AGENTS.md` §3.12) | unverändert: Jedes Rezept bleibt eine Zeile `@bash <Ort>`; kein neuer Fall, weil keine neue Zusage; die Fälle der Aufrufarten in den Gegenproben (`p10-*` in `gegenproben-gegenprobe.sh`, die Fälle zu den Flags von `make` in `gates-schnell-gegenprobe.sh` und `stop-hook-gegenprobe.sh`) laufen im Nachweis mit |
| (g) Hook im Bare-Repo | Der Hook des Gegenproben-Containers ruft `make gegenproben` im Arbeitsbaum auf dem eingelieferten Hash | liest keinen Pfad einer Gegenprobe (`tools/` steht in `gegenproben-nach-einlieferung.sh` nur im Kommentar, Zeile 7); der Kommentar folgt dem Ort; eine schon eingespielte Kopie weicht danach vom Skript ab (Risiko 5) |
| (h) Lesen von `tools/` | Gegenproben und Gates, die Dateien unter `tools/` zählen, suchen oder gegen eine Liste prüfen | nur `abdeckung.sh` (Zeile 148, `find tools -type f -name '*-gegenprobe.sh'`, jede Tiefe, `SPEC-048` Punkt 7): findet `tools/dev/…` ohne Änderung, die Tabelle nennt den neuen Pfad; die Fixture `tools/test/x-gegenprobe.sh` (`abdeckung-gegenprobe.sh` Zeile 434) ist ein erfundener Pfad im Temp-Baum und bleibt; `stop-hook-gegenprobe.sh` Zeile 1062 (`aufnahme 'tools/*'`) prüft den Pfadnamen der Positivliste, nicht den Bestand. Befund des Planners mit `git grep -n -E 'find tools\|ls tools\|tools/\*' -- tools harness/mk`; `.d-check.yml`, `.golangci.yml`, `.a-check.yml`, `Dockerfile` (außer Zeile 61), `Makefile`, `.github/` und `.githooks/` nennen keine der 21 |
| (i) Stop-Hook und Nachweis | Der Move ändert die Dateiliste des Arbeitsbaums; `tools/` steht nicht auf der Positivliste des Stop-Hooks | `make gates` läuft als letzter Lauf nach dem letzten Commit; ein älterer Nachweis entfällt |
| (j) Plattform | die drei Plattformen des Lastenhefts: Linux, macOS, Windows | geprüft ist Linux (Dev-Container, Gegenproben-Container); macOS und Windows sind nicht zugesagt (Kopf von `SPEC-060`); `git mv` und die neuen Pfade hängen an keiner Plattform (nur ASCII, kein Leerraum); der Modus (Ausführbit) bleibt beim Move |
| (k) Zwischenstand | Nach dem Move-Commit zeigen alle 21 Rezepte ins Leere | geliefert wird erst nach Liefer-Punkt 2; der Zwischenstand wird nicht einzeln geprüft |
| (l) Offene Zweige | `git branch` führt alte Zweige (`rev-fahr`, `rev-start`, 2026-10-18) mit früheren Fassungen von `gegenproben-nach-einlieferung-gegenprobe.sh` und `container-start-gegenprobe.sh` | nicht Gegenstand; wer sie mergt, trifft einen Rename-Konflikt (git erkennt den Umzug, weil der Pfadnachzug wenige Zeilen je Datei ändert) |

**Risiken:**

1. Eine Gegenprobe wird nach dem Umzug grün, prüft aber nicht mehr dasselbe: ein Verweis auf das Skript unter Test zeigt ins Leere und der Fall überspringt, oder eine Mutation trifft ihre Textstelle nicht (`BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`) — **Ausgang:** offen bis Closure (entfallen: DoD 3 zeigt je Gegenprobe den Lauf, bei den Skripten mit `GP_ENTDECKE` die roten Fälle je Mutation unverändert; eingetreten: Rückführung nach §4).
2. `release-gegenprobe` und `test-container-gegenprobe` laufen nur mit der engsten Auswahl; ein Pfadfehler an einer Stelle, die sie nicht trifft, bleibt bis zum ersten vollen Lauf unbemerkt — **Ausgang:** offen bis Closure (entfallen: Die Auswahl trifft den ungemutierten Fall und jede Stelle, die §3 Schritt 7 nennt; weiter offen → `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`; den vollen Lauf am Endstand verlangt DoD 1 von `slice-harness-gegenproben-schneller`).
3. `container-start-gegenprobe` ist im Dev-Container nicht belegbar — **Ausgang:** offen bis Closure (entfallen: der Lauf am Host steht in §7; weiter offen → `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`).
4. Ein anderer Slice startet vor dem Umzug oder vor dem Eintrag in die Nehmer und legt eine Gegenprobe am alten Ort an oder ändert eine dort — **Ausgang:** offen bis Closure (Gegenmittel: Dieser Slice liegt vor `slice-harness-gegenproben-schneller` und vor jedem Slice, der eine Gegenprobe anlegt; eingetreten: ein Rename-Konflikt, den der Slice löst, der ihn trifft).
5. Die Kopie des Hooks im Bare-Repo weicht nach der Kommentaränderung (Zeile 7) vom Skript ab, wenn sie schon eingespielt ist; `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` vergleicht die Summe von Kopie und Skript in der ersten Fahrt — **Ausgang:** offen bis Closure (entfallen: der Hook ist noch nicht eingespielt, jener Slice liegt in `open/`; eingetreten: jener Slice spielt neu ein, die Closure trägt es dort ein).
6. Der Slice wächst durch die Sendung an 19 Nehmer, jeder mit Nachzählung (`BEO-REPO/slice-waechst-durch-uebernahmen`, 7×) — **Ausgang:** offen bis Closure (entfallen: Jede Zählung steht in §8 des Nehmers, und kein Nehmer liegt über drei Liefer-Punkten oder zwei Schichten, oder die Zählung nennt den vorgemerkten Schnitt).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). **Ort der Belege, die die DoD verlangt:** hier — die Aufrufzeile und das Ergebnis des Move-Commits (DoD 1), der `git grep` des Nachzugs (DoD 2), die Tabelle der Wirkung mit 21 Zeilen (DoD 3) und der Lauf am Host für `container-start-gegenprobe` (§5); Entscheidungen des Planners zur Meldeschwelle stehen hier als datierte Zeilen (§4).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`, Greenfield); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist für eine Sub-Area, die das ganze Repo deckt, nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen (gemergter Stand, `1ad628a7`). Treffer: `gegenprobe-nicht-in-kette-regression-unbemerkt` (5×, Stand geplant, Adresse `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`; Risiko 1 und 2 sind die Berührung, ein weiteres Auftreten entstünde nur, wenn eine nicht gefahrene Gegenprobe nach dem Umzug bricht), `slice-waechst-durch-uebernahmen` (7×, verkörpert in `AGENTS.md` §3.13; die Sendung an 19 Nehmer verlangt die Nachzählung je Nehmer, DoD und Risiko 6), `folge-slice-adresse-nimmt-nicht-an` (13×; die Adresse für Schritt 2 nimmt an: `slice-harness-dev-kopie-und-pin` trägt die Sendung in §1 und DoD, angelegt im Commit mit diesem Plan), `zusage-im-kommentar-weiter-als-pruefung` (93×; jeder Ort in einem Kommentar, den der Umzug berührt, ist Zusage, DoD 2 belegt ihn mit `git grep`), `gegenprobe-variable-kollidiert-mit-skript` (1×, offen; §3 Schritt 3 wählt das Präfix `gp_` für jede neue Variable), `gate-laufzeit-waechst-mit-gegenprobe` (14×, verkörpert; die Kette wächst nicht), `gegenprobe-haengt-an-umgebung-des-hosts` (10×; Risiko 3), `schaetzung-in-zeilen-unterschaetzt-tests` (19×; §4 nennt Spanne, Marke und Meldeschwelle; ein Auftreten wäre ein Stand über 400). Erreicht keiner mit diesem Slice 3×, der nicht schon darüber liegt.

**Zählung** (Planner, 2026-10-22, `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 3** (1: Umzug der 21; 2: Pfadnachzug; 3: Wirkung belegt). **Schichten: 2** — Harness-Werkzeug (Skripte, Fragmente, Kommentare in `Dockerfile` und `.claude/hooks/`, `harness/conventions.md`, `harness/README.md`) und Nutzer-Doku (`docs/user/abdeckung-*.md`, mit `make abdeckung` regeneriert; `AGENTS.md` §3.13 nimmt sie nicht aus, die Ausnahme „Abdeckungstabellen zählen nicht als Schicht“ steht erst im Plan `slice-harness-regeln-nach-v1-abschluss` in `next/`, ist nicht angenommen, und die Tabellen sind deshalb als Schicht gezählt). Planung (Pläne, Register) ist keine Schicht; die Ortsangaben in den Verträgen schreibt der Architect und sind kein Liefer-Punkt. Sendungen anderer Slices an diesen: keine (er ist der Geber). Größe: Marke in §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Verträge und Konventionen des Harness gehen den Skripten voraus).
