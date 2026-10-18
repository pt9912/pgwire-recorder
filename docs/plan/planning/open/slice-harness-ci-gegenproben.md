# Slice slice-harness-ci-gegenproben: Der nächtliche Lauf der Gegenproben auf GitHub — `gegenproben.yml`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-ci-nachtlaeufe (Entscheidung des Nutzers vom 2026-10-18), erste Stelle der Reihenfolge ([welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md) §5; bis zur Entscheidung vom 2026-10-18 dritte Stelle von welle-erster-release, die Welle verlässt dieser Slice). Start-Trigger weicht von der Welle ab: `slice-harness-ci` und `slice-harness-gegenproben-aus-gates-sammelziele` in `done/` (§4); `gegenproben` und `gegenproben-host` entstehen erst im zweiten; dazu der Start-Trigger der Welle, `welle-erster-release` in `done/`.

**Bezug:** Kein Anspruch des Produkts (`LH-*`): der Slice lässt die Gegenproben des Harness regelmäßig an einem Ort laufen, den kein Agentenlauf stört. Aktive ADRs: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Gate-Kette ohne Gegenproben), [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md) und [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) (Host-Voraussetzungen und Host-Liste). Die ADR dieses Slice (Beschaffung der Host-Werkzeuge auf dem Runner) schreibt der Architect vor dem Code (Entscheidung des Nutzers).

**Berührte Spec-Stellen:** `SPEC-065` (neu; Nummer vergibt der Architect, nächste freie; Vertrag in `harness/vertraege/`, vor dem ersten Code-Commit) · `SPEC-060` Teil 2 (gelesen: Namen und Verhalten der Sammelziele) · `SPEC-063` und `SPEC-064` (gelesen: Gate der Workflows, Rechte und Pins der CI) · [`SPEC-058`](../../../../harness/vertraege/SPEC-058-container-start.md) Punkt 3 (gelesen: Host-Voraussetzungen von `u-boot`, `devcontainer`, `node`, `ssh`, `ssh-keygen`).

**Verantwortlich:** —
**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `.github/workflows/gegenproben.yml` fährt nachts (`schedule`) und von Hand (`workflow_dispatch`) zwei Jobs auf dem Standard-Zweig bzw. dem gewählten Zweig: `make gegenproben` (Container-Liste) und `make gegenproben-host` (Host-Liste, braucht `u-boot` ab 0.10.0, die `devcontainer`-CLI 0.80.2, `node`, `ssh` und `ssh-keygen`), beide mit `permissions: {}` oben und `contents: read` im Job, Zeitgrenze, Pins nach `SPEC-063`. Die Werkzeuge, die der Runner nicht mitbringt, installieren **Schritte der Workflow-Datei** mit festen Versionen und Prüfsummen; das Repo-Skript ruft keinen Host-Paketmanager (`AGENTS.md` §3.1). Rot ist ein roter Lauf in der Actions-Übersicht; der Slice legt weder Issue noch Nachricht an.

**Beschaffung der Host-Werkzeuge** (Empfehlung des Planners, **Entscheidung beim Architect**, Randform 2): (A) Schritte im Workflow — `actions/setup-node` mit fester Version (per SHA gepinnt) und `npm install -g @devcontainers/cli@0.80.2`, `u-boot` als Release-Asset mit fester Version und gegen eine im Workflow stehende SHA256-Prüfsumme verifiziert, abgelegt in einem Verzeichnis im `PATH` des Jobs. (B) ein Make-Ziel `runner-werkzeuge` im Repo — ausgeschlossen: Es wäre ein Host-Paketmanager im Repo-Skript. (C) ein Image mit den Werkzeugen, in dem der Host-Job läuft — verfehlt den Zweck, denn die Host-Gegenprobe soll am Host laufen. Empfehlung: (A).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Sammelziele `make gegenproben` und `make gegenproben-host` und die Listen — `slice-harness-gegenproben-aus-gates-sammelziele` (Geber, vor diesem Slice); dieser Slice ruft sie und ändert sie nicht.
- Die Pflicht-Fahrstelle (der eigene Dev-Container, am Ende jeder Welle und vor dem Release) — `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (Pflicht in den Anweisungen, Abschnitt im Dokument, erste Fahrt; der Hook für die Einlieferung liegt in `slice-harness-gegenproben-aus-gates-fahrstelle`); der Nachtlauf ergänzt sie, er ersetzt sie nicht (Frage 1), und der Wortlaut von `AGENTS.md` §6 ist Entscheidung des Nutzers vom 2026-10-18.
- `ci.yml`, `make gates`, die Range-Prüfung — `slice-harness-ci`; ein roter Nachtlauf sperrt keinen Merge.
- Aktualität der Pins — `slice-harness-ci-freshness`; ein Nachtlauf über Fremdbestände ist ein anderer Workflow mit anderer Fehlersemantik.
- Die Prüfung der Datei auf Pins, Rechte, Zeitgrenze und Syntax — `make workflows-check` aus `slice-harness-ci-workflows`; dieser Slice baut keinen zweiten Prüfer.
- Ein Eintrag in `GATE_CHECKS`, ein Hook oder ein Stand, der dem Nachtlauf glaubt — die Sammelziele sind kein Gate (`SPEC-060`), und ein Nachweis, den ein Lauf auf GitHub schriebe, würde der lokale Hook nicht lesen.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Container-Job):** `.github/workflows/gegenproben.yml` mit `schedule` (Uhrzeit abseits der vollen Stunde, vom Vertrag festgelegt) und `workflow_dispatch` ruft in einem Job `make gegenproben`, hält die Zusagen von `SPEC-065` (Rechte, Zeitgrenze, `concurrency`, Checkout, Pins) und besteht `make workflows-check` und `make docs-check`; jede Entscheidung je Randform aus §6 steht vor dem ersten Code-Commit in `SPEC-065`. Belegt in §7 mit der Nummer eines grünen Laufs auf GitHub und der Dauer.
- [ ] **Liefer-Punkt 2 (Host-Job):** Ein zweiter Job installiert `u-boot`, `node` und die `devcontainer`-CLI nach dem im Vertrag entschiedenen Weg (Versionen fest, Prüfsumme verifiziert, Abbruch bei Abweichung vor dem ersten Aufruf), prüft die Versionen gegen `SPEC-058` Punkt 3 und ruft `make gegenproben-host`; läuft die Gegenprobe des Container-Starts auf dem Runner nicht (Randform 3), steht das im Vertrag mit der Folge. Belegt in §7 mit dem Lauf oder mit dem benannten Befund.
- [ ] **Liefer-Punkt 3 (Gegenprobe):** `make ci-gegenproben-gegenprobe` (Name: Architect) zeigt an Kopien von `gegenproben.yml` je Zusage von `SPEC-065` und je Element ihrer Aufzählungen (Auslöser, Rechte, Zeitgrenze, beide Ziele, Pin der Werkzeug-Versionen, Prüfsumme vor dem Aufruf) einen Fall und eine Mutation, die nur ihn färbt (`AGENTS.md` §3.10); die unveränderte Kopie grün; `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` je ein Fall. Das Ziel steht in `GEGENPROBEN_CONTAINER`. `harness/README.md` §Sensors führt die Zeile und nennt den Nachtlauf als Ergänzung der Pflicht aus `AGENTS.md` §6, nicht als Ersatz.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `harness/` und `tools/` ändern) und `make gates-schnell-gegenprobe` grün; `make ci-gegenproben-gegenprobe` grün auf dem Endstand.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/gegenproben.yml` | neu | zwei Jobs, Auslöser, Rechte, Installation der Host-Werkzeuge; Vorbild d-check `image-scan.yml` (Nachtlauf mit `workflow_dispatch`) |
| `harness/mk/…`, `tools/harness/…` (Gegenprobe der Datei; Namen: Architect) | neu | Gegenprobe der Zusagen an Kopien; Eintrag in `GEGENPROBEN_CONTAINER` |
| `harness/vertraege/SPEC-065-…md` | neu (Architect) | Vertrag; vor dem Code |
| `harness/README.md` §Sensors | update | Zeile der Gegenprobe; Verhältnis zur Pflicht in `AGENTS.md` §6 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-ci` und `slice-harness-gegenproben-aus-gates-sammelziele` liegen in `done/`, dazu der Start-Trigger von [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md) (`welle-erster-release` in `done/`), und der Architect hat `SPEC-065` und die ADR geschrieben, jede Randform aus §6 entschieden (`AGENTS.md` §3.12). Schritt 1 der Reihenfolge in §5 der Welle.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: mehr als drei Liefer-Punkte, oder der Diff ist nicht in einer Review-Sitzung prüfbar — zurück zur Zerlegung. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- .github tools harness Makefile`, Tests und Gegenprobe eingeschlossen, liegt über 900; Basis: beim Anspruch (HEAD vor dem Anspruch, hier einzutragen). Schätzung des Planners in dieser Metrik (2026-10-18, ohne Basis): `gegenproben.yml` 110 bis 160, Gegenprobe der Zusagen 200 bis 320, `SPEC-065` 100 bis 160, `harness/README.md` 10 bis 20, zusammen 420 bis 660, unter der Marke.
- `in-progress` → `open`: die Gegenprobe des Container-Starts läuft auf dem Runner nicht (Randform 3) und der Nutzer entscheidet, den Host-Job nicht zu führen, oder die Quelle von `u-boot` ist ohne Prüfsumme nicht zu beziehen (Randform 2).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, ein von Hand angestoßener Lauf (`workflow_dispatch`) beider Jobs auf GitHub für den Endstand als Beleg in §7, `make gates` grün, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die Gegenprobe des Container-Starts startet einen Dev-Container mit eigenem Docker-Daemon; ob das auf dem Runner läuft** (verschachtelte Laufzeit, Rechte, Ports, `sshd`), zeigt nur ein Lauf (`BEO-REPO/nested-runtime-podman-oder-docker-offen`). — **Ausgang (vorgesehen):** entfallen, wenn der Lauf in §7 grün endet; sonst eingetreten → Rückführung nach `open` (§4) und Beleg bei `BEO-REPO/nested-runtime-podman-oder-docker-offen`.
- **Ein Nachtlauf, in den niemand sieht, ist ein verwaister Sensor; GitHub schaltet `schedule` in einem Repo ohne Aktivität nach 60 Tagen ab, ein abgeschalteter Lauf sieht aus wie ein grüner.** — **Ausgang (vorgesehen):** weiter offen → neues `BEO-REPO/nachtlauf-ohne-leser` bei der Closure, wenn ein roter Lauf unbemerkt bleibt; der Wortlaut in `harness/README.md` nennt die Grenze.
- **Die Versionen und Prüfsummen von `u-boot` und der `devcontainer`-CLI altern im Workflow still.** — **Ausgang (vorgesehen):** weiter offen → neues `BEO-REPO/ci-werkzeug-versionen-altern-still` bei der Closure; kein Slice nimmt die Achse an (Frage 3), die Grenze steht in `SPEC-065`.
- **Der Nachtlauf beweist den Stand des Standard-Zweigs, nicht den eines Releases.** — **Ausgang (vorgesehen):** entfallen, weil die Pflicht vor dem Release (`AGENTS.md` §6) beim Gegenproben-Container bleibt und `slice-erster-release-freigabe` sie in die Checkliste trägt; der Nachtlauf wird nicht als Beleg des Release-Stands behauptet.

**Randformen** (`AGENTS.md` §3.12) — **alle offen für den Architect**, zu entscheiden vor dem ersten Code-Commit in `SPEC-065` (Technik-Stratum, `harness/vertraege/`), mit ADR für Entscheidung und Gründe. Vorbereitet vom Planner am 2026-10-18 (Empfehlungen, kein Beschluss):

1. **Auslöser.** `schedule` (Uhrzeit abseits der vollen Stunde, zwei Stunden versetzt zu den übrigen Nachtläufen) läuft nur auf dem Standard-Zweig und nicht in Forks; `workflow_dispatch` darf jeden Zweig wählen und prüft dann diesen. Ob `workflow_dispatch` Pflicht ist: ja, wie bei den Nachtläufen von d-check (ein Workflow, der nur an einem seltenen Ereignis hängt, wird erst im Ernstfall geprobt). Ein Fall je Auslöser.
2. **Beschaffung der Host-Werkzeuge** (siehe §1): Quelle und Form von `u-boot` (Release-Asset, Archiv oder Binary, Plattform `linux/amd64`), Prüfsumme (im Workflow oder aus dem Release), Ablage im `PATH`, Versionen fest (`SPEC-058` Punkt 3: `u-boot` ab 0.10.0, `devcontainer` 0.80.2), `node`-Version (die des Runners oder gepinnt), Abbruch bei fehlender oder abweichender Prüfsumme vor dem ersten Aufruf; Aktionen von Dritten nur per SHA. Ein Fall je Element.
3. **Läuft die Gegenprobe des Container-Starts auf dem Runner?** Sie ist ein Host-Ziel ([ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md)) und startet einen Dev-Container mit eigenem Docker-Daemon; der „Host“ ist hier die Runner-VM. Fälle: `ssh-keygen` und Schlüsselablage im Runner-Home, Ports, verschachtelte Laufzeit, Dauer des Erststarts. Empfehlung: erst der Lauf, dann die Entscheidung; trägt er nicht, nennt der Vertrag den Host-Job als nicht geführt und `make gegenproben-host` bleibt Handlung am Host.
4. **Zwei Jobs oder einer.** Zwei Jobs: getrennte Zuordnung des Rots, getrennte Zeitgrenzen, parallel. Ein Job: eine Installation, nacheinander (`SPEC-052` Punkt 8, Aufruf-Voraussetzung nacheinander). Empfehlung: zwei Jobs; die Sammelziele fahren je Liste nacheinander.
5. **Zeitgrenze.** Je Job; die Dauer der Container-Liste ist erst nach `slice-harness-gegenproben-aus-gates-sammelziele` gemessen (die Zahlen stehen in dessen Bericht). Empfehlung: Grenze aus der gemessenen Dauer mit Faktor 2, nicht aus der Annahme; die Messung steht in §7.
6. **Gleichzeitige Läufe.** `concurrency` mit fester Gruppe, `cancel-in-progress: false` (ein Lauf zu Ende, der nächste wartet); ein von Hand angestoßener Lauf kollidiert mit dem nächtlichen.
7. **Rechte, Geheimnisse, Pins.** `permissions: {}` oben, im Job `contents: read`; keine Geheimnisse; `persist-credentials: false`; nur Aktionen per SHA mit Tag-Kommentar. Kein Docker-Hub-Login (Randform 13 in `slice-harness-ci`).
8. **Wer sieht das Rot.** Kein Issue, keine Nachricht im Workflow (`issues: write` wäre ein Recht mehr). Empfehlung: ein roter Lauf in der Actions-Übersicht und die Benachrichtigung, die GitHub dem Auslöser schickt; eine eigene Meldung ist ein Slice nach Beobachtung `BEO-REPO/nachtlauf-ohne-leser`.
9. **Ausgabe und Exit.** Der Workflow wertet nur den Exit der Sammelziele (`SPEC-060`: 1 bei mindestens einem roten) und reicht ihn durch; er übersetzt keine Meldung. Die Zusammenfassung der Sammelziele steht im Log. Ein Fall: Exit des Ziels wird nicht überschrieben (`continue-on-error` fehlt).
10. **Aufrufarten des Treibers.** Der Workflow ruft `make gegenproben` ohne weitere Flags; `-j` und `-k` stehen nicht in der Datei, `MAKEFLAGS` setzt der Runner nicht. Ein Fall, dass `ci-gegenproben` keine Flags enthält.

**Fragen an den Nutzer** (jede mit Empfehlung des Planners; eine Empfehlung ist keine Entscheidung):

1. **Ersetzt ein grüner Nachtlauf des Release-Stands das Fahren im Gegenproben-Container vor dem Release?** Empfehlung: nein, er ergänzt; der Wortlaut von `AGENTS.md` §6 bleibt (deine Entscheidung vom 2026-10-18).
2. **Quelle von `u-boot`** auf dem Runner: Release-Asset im öffentlichen Repo? Ohne öffentliche Quelle entfällt der Host-Job (§4, Rückführung). Empfehlung: die Quelle benennen, bevor der Architect schreibt.
3. **Aktualität der Werkzeug-Versionen** (`u-boot`, `devcontainer`) im Workflow. Kein Slice nimmt sie an; `slice-harness-ci-freshness-images` prüft Images, nicht Werkzeug-Versionen eines Workflows. Empfehlung: nach der ersten Beobachtung (Risiko 3) einen Folge-Slice anlegen, nicht vorab.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (`AGENTS.md` §3.10) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (`AGENTS.md` §3.12) — die Datei ist ein neuer Vertrag, Randformen stehen in §6 vor dem Code; `BEO-REPO/nested-runtime-podman-oder-docker-offen` (Risiko 1); `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` — dieser Slice ist die Fahrstelle der Gegenproben in der Cloud, nicht der Kette; `BEO-REPO/verhalten-nur-unter-linux-geprueft` — der Runner ist Linux, die Aussage bleibt auf `linux/amd64` begrenzt. Kein Eintrag erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Der Slice trägt keine Sendung eines Gebers (`SPEC-060`, `SPEC-063`, `SPEC-064`, `SPEC-058` sind gelesen). Liefer-Punkte 3 (Container-Job | Host-Job | Gegenprobe samt Gate-Index). Schichten 2: Workflows und Harness; Planung, Vertrag und ADR sind keine Schicht; Wartungs-Doku und Handbuch bleiben unberührt. Schätzung 420 bis 660 (§4), unter der Marke.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit, der die Welle wechselt; Entscheidung des Nutzers vom 2026-10-18: der erste Release hat Vorrang vor den Nachtläufen): Der Slice wechselt von welle-erster-release in welle-ci-nachtlaeufe und startet nach `welle-erster-release` in `done/`; Inhalt, Sendungen, Liefer-Punkte und Schichten bleiben unverändert (Liefer-Punkte und Schichten wie in der Zählung oben, Marke in §4 unverändert).
