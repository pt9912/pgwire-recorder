# Slice slice-harness-ci-workflows: Gate der Workflows — Pins, Rechte und Zeitgrenzen der Dateien unter `.github/workflows/`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-erster-release, erste Stelle der Reihenfolge, vor `slice-harness-ci` und vor `slice-erster-release-veroeffentlichung` (Auftrag des Nutzers vom 2026-10-18: die CI gehört in welle-erster-release, vor die Veröffentlichung; **der Schnitt in fünf Slices und die Reihenfolge darin sind Empfehlung des Planners, nicht entschieden**, Frage 1). Start-Trigger weicht von der Welle ab: `slice-harness-gegenproben-aus-gates` in `done/` (die Namens- und Zuordnungsregel von `SPEC-060` Teil 1 legt die Gegenprobe dieses Slice in `GEGENPROBEN_CONTAINER`; ohne die Listen hinge sie an `GATE_CHECKS` und müsste später umziehen), dazu der Start-Trigger der Welle (§4).

**Bezug:** Kein Anspruch des Produkts (`LH-*`): der Slice ändert Gates und Prüfungen des Harness, nicht das Binary. Aktive ADRs: [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Docker-only, Prüfungen als Stufen und Werkzeug-Images), [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Gate-Kette, Ausnahmeliste der schnellen Kette), [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) (Listen der Gegenproben, Host-Ziel). [ADR-0050](../../adr/0050-ci-auf-github-actions.md) (CI auf GitHub Actions, Accepted 2026-10-18) trägt Entscheidung und Gründe dieses Slice.

**Berührte Spec-Stellen:** `SPEC-063` (neu; Vertrag in `harness/vertraege/`, vor dem ersten Code-Commit geschrieben) · `SPEC-060` Teil 1 (gelesen: Namens- und Zuordnungsregel, Listen; Vertrag vom Architect am 2026-10-18 geschrieben) · [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) Punkte 1 bis 4 (gelesen: Aufrufarten des Treibers, Ausnahmeliste).

**Verantwortlich:** pt9912 (Zweig impl-ci-workflows)
**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Gate `make workflows-check` (`SPEC-063`) prüft jede Workflow-Datei unter `.github/workflows/` (`*.yml` und `*.yaml`), bevor es eine gibt, auf die ein Mensch oder `make gates` sich verlassen müsste: (a) **Syntax und Ausdrücke** mit `actionlint` (Image per Digest gepinnt, netzlos, schreibgeschützt); (b) **Rechte und Zeitgrenzen** durch ein kleines eigenes Skript: `permissions: {}` auf oberster Ebene, jeder Job mit `timeout-minutes`; (c) **Bindung der Pin-Prüfung:** Die Pins (voller Commit-SHA mit Tag-Kommentar je fremde `uses:`-Referenz) prüft allein `make docs-check` über das Modul `workflows` von d-check; das Gate prüft nicht die Pins, sondern nur, dass `.d-check.yml` das Modul mit `dir: .github/workflows` führt, sobald eine Workflow-Datei existiert (`SPEC-063` Punkt 8). Der Eintrag selbst wandert in `slice-harness-ci` zusammen mit `ci.yml`: Das Modul meldet ohne Datei einen Befund (leere Prüfmenge), er kann vor der ersten Datei nicht stehen. Das Gate hängt an `GATE_CHECKS`; seine Gegenprobe `make workflows-check-gegenprobe` zeigt je Regel und je Element ihrer Aufzählungen einen Fall und eine Mutation und steht in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS`.

Der Slice liefert das Gate **vor** der ersten Workflow-Datei. Er legt keine an; ohne Datei meldet er das und endet mit Exit 0 (Randform 1), und jede Datei der folgenden Slices (`slice-harness-ci`, `-ci-gegenproben`, `-ci-freshness`, `slice-erster-release-veroeffentlichung`) fällt vom ersten Commit an unter ihn.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Eintrag des Moduls `workflows` in `.d-check.yml` — kommt mit der ersten Datei, `slice-harness-ci` (Sendung, dort in §1 und DoD genannt); davor meldete `make docs-check` einen Befund über leerer Prüfmenge.
- Die Workflow-Dateien selbst (`ci.yml` und die übrigen) — `slice-harness-ci`, `slice-harness-ci-gegenproben`, `slice-harness-ci-freshness`; hier gibt es nichts, das ein Workflow ausführte, und das Gate ist ohne Datei prüfbar (Gegenprobe in Temp-Bäumen).
- Die Gültigkeit und Aktualität der SHA-Pins (existiert der Commit, bezeichnet er den Tag im Kommentar, gibt es einen neueren) — Netz, außerhalb der Gates; das Modul `workflows` von d-check nennt dieselbe Grenze. Kein Slice nimmt sie an (Frage 3); dieser Slice behauptet sie nicht (`AGENTS.md` §3.11).
- Die Lauffähigkeit auf GitHub (Runner, Rechte der Organisation, `startup_failure`) — nur dort messbar; das Gate prüft die Deklarationsklasse, nicht den Lauf.
- Branch-Schutz und „required status check“ — Einstellung auf GitHub, aus dem Klon nicht auditierbar; Grenze, die `slice-harness-ci` benennt.
- Eigene Actions unter `.github/actions/` und wiederverwendbare Workflows — das Repo hat keine; das Gate liest sie nicht (Randform 3), und ein Zusatz wäre eine Änderung des Plans, kein Zubehör.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Gate):** `make workflows-check` prüft nach `SPEC-063` die drei Gruppen (Syntax und Ausdrücke, Rechte und Zeitgrenzen, Bindung der Pin-Prüfung) über jede Datei unter `.github/workflows/`, nennt je Befund Datei, Job und Regel, endet mit Exit ungleich 0 bei einem Befund, schreibt nichts in den Arbeitsbaum, braucht kein Netz zur Laufzeit und hängt an `GATE_CHECKS`; jede Entscheidung je Randform aus §6 steht vor dem ersten Code-Commit in `SPEC-063`. Fragment `harness/mk/workflows.mk`, Skript `tools/harness/workflows-check.sh`; der Eintrag in `.d-check.yml` gehört nicht dazu (`slice-harness-ci`).
- [ ] **Liefer-Punkt 2 (Gegenprobe):** `make workflows-check-gegenprobe` zeigt je Punkt von `SPEC-063` und je Element seiner Aufzählungen (Dateiendung, Form von `permissions`, Job-Art, Bindung der Pin-Prüfung) einen Fall, der abgelehnt, und einen, der angenommen wird, dazu je Zusage eine Mutation des Skripts, die nur den Fall ihrer Gruppe färbt (`AGENTS.md` §3.10); die unveränderte Kopie grün; Aufrufarten `-n`, `-q`, `-t`, `-i` und `MAKEFLAGS` aus der Umgebung je ein Fall (`AGENTS.md` §3.12). Das Ziel steht in `GEGENPROBEN_CONTAINER` (Namens- und Zuordnungsregel von `SPEC-060`).
- [ ] **Liefer-Punkt 3 (Gate-Index):** `harness/README.md` §Sensors führt die Zeilen von Gate und Gegenprobe mit Vertrag und Bindung; die Prosa über das Gate nennt die Regel an den Schreiber einer Workflow-Datei und zeigt auf `SPEC-063` und `.d-check.yml` als Quelle, nicht auf eine nachgeschriebene Aufzählung (`AGENTS.md` §3.11).
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `harness/` und `tools/` ändern) und `make gates-schnell-gegenprobe` grün (die Deckung zählt die neue Gegenprobe); `make workflows-check-gegenprobe` grün auf dem Endstand.
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
| `harness/mk/workflows.mk`, `tools/harness/workflows-check.sh` | neu | Gate: `actionlint` im gepinnten Image, eigene Regeln für `permissions` und `timeout-minutes`, Bindung der Pin-Prüfung in `.d-check.yml` (lesend); hängt an `GATE_CHECKS` |
| `harness/mk/workflows-gegenprobe.mk`, `tools/harness/workflows-check-gegenprobe.sh` | neu | Gegenprobe in Temp-Bäumen und Kopien; Kopf mit Abdeckungs-Deklaration, falls der Architect eine Anforderung nennt; Eintrag in `GEGENPROBEN_CONTAINER` |
| `harness/vertraege/SPEC-063-…md` | neu (Architect) | Vertrag des Gates; vor dem Code |
| `harness/README.md` §Sensors | update | zwei Zeilen und die Prosa nach §3.11 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-gegenproben-aus-gates` liegt in `done/` (Listen und Namensregel), dazu der Start-Trigger von [welle-erster-release](../welle-erster-release.md) (§2 dort), und der Architect hat `SPEC-063` und die ADR geschrieben, jede Randform aus §6 entschieden (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: das Gate verlangt mehr als drei Liefer-Punkte, oder der Diff ist nicht in einer Review-Sitzung prüfbar — zurück zur Zerlegung. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- .github tools harness Makefile`, Tests und Gegenprobe eingeschlossen, liegt über 900; Basis: Elternteil des Commits, der diese Datei nach `in-progress/` bewegt, `a98d0b9e`. Schätzung des Planners in dieser Metrik (2026-10-18, ohne Basis): Fragment und Skript des Gates 90 bis 130, Gegenprobe das Doppelte bis Dreifache des Skripts (220 bis 340), `SPEC-063` 120 bis 180, `harness/README.md` und `.d-check.yml` 25 bis 40, zusammen 455 bis 690 (die Zeile `.d-check.yml` entfällt mit Randform 12; maßgeblich ist die Schätzung des Architect in §8). Schnitt, falls die Marke eintritt: *Syntax* (Gruppe a) und *Rechte, Zeitgrenzen, Bindung* (Gruppen b und c).
- `in-progress` → `open`: `slice-harness-gegenproben-aus-gates` liegt nicht in `done/` — extern klären.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` und `make workflows-check-gegenprobe` grün auf dem Endstand, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Das Gate prüft eine Klasse, nicht die Lauffähigkeit** (die Grenze des Moduls `workflows` und von `actionlint`): ein grünes Gate sagt nichts über den Lauf auf GitHub. — **Ausgang (vorgesehen):** weiter offen → neues `BEO-REPO/workflow-gate-gruen-lauf-rot` bei der Closure, wenn `slice-harness-ci` einen Lauf trifft, den das Gate hätte fangen sollen; sonst entfallen mit dem Hinweis in `SPEC-063`.
- **Der Pin-Check hängt am gepinnten d-check** (Version des Moduls, Verhalten bei `.yaml`, Grenze der Zerlegung der Blockform); er läuft in `make docs-check`, nicht in diesem Gate, und die Bindung prüft dieses Gate nach `SPEC-063` Punkt 8. — **Ausgang (vorgesehen):** weiter offen → an `slice-harness-ci`, der den Eintrag setzt und die Pins mit der ersten Datei zum ersten Mal wirklich prüft (`SPEC-064` Punkt 12 zeigt es an echtem d-check); sonst entfallen mit Grund bei der Closure.
- **`actionlint` ruft `shellcheck` auf `run:`-Blöcke auf**; Befunde dort sind Stil, nicht Fehler, und könnten das Gate dauerhaft rot machen. — **Ausgang (vorgesehen):** entfallen, wenn `SPEC-063` die Einstufung festlegt (Randform 5) und die Gegenprobe sie zeigt.
- **Die Gegenprobe braucht das Image von `actionlint`** (Netz beim ersten Pull, wie jedes gepinnte Werkzeug-Image). — **Ausgang (vorgesehen):** entfallen, weil die Gegenprobe in `GEGENPROBEN_CONTAINER` steht und kein Teil von `make gates` ist; der Lauf des Gates selbst ist netzlos nach dem Pull.

**Randformen** (`AGENTS.md` §3.12) — **alle offen für den Architect**, zu entscheiden vor dem ersten Code-Commit in `SPEC-063` (Technik-Stratum, `harness/vertraege/`), mit ADR für Entscheidung und Gründe. Vorbereitet vom Planner am 2026-10-18 (Empfehlungen, kein Beschluss):

1. **Kein Verzeichnis, leeres Verzeichnis, keine Datei.** Empfehlung: Exit 0 mit einer Zeile „keine Workflows“, nicht still; ein Fall je Form. Gate-Wahrheit „grün über leerem Prüfbereich“ benennt `SPEC-063` als Grenze (`LH-QA-01`-Muster wie `history-range-guard`). **Entschieden:** `SPEC-063` Punkt 2 (Exit 0 mit Meldung, ohne Docker; Grenze in Punkt 11).
2. **Dateiendungen und Ort.** GitHub liest `*.yml` und `*.yaml` unmittelbar unter `.github/workflows/`, keine Unterverzeichnisse. Empfehlung: beide Endungen, Unterverzeichnisse und andere Endungen nicht geprüft, aber mit Hinweis genannt; ein Fall je Endung und je Tiefe. **Entschieden:** `SPEC-063` Punkt 1 (beide Endungen, kein Unterverzeichnis, Hinweis).
3. **Weitere Orte.** `.github/actions/*/action.yml` (zusammengesetzte Actions), `.github/dependabot.yml`: nicht geprüft; nennt der Vertrag als Grenze, weil das Repo keine führt. **Entschieden:** `SPEC-063` Punkt 11 (nicht gelesen, Grenze).
4. **Formen von `permissions`.** Oberste Ebene: fehlt; `permissions: {}`; leerer Block; `read-all`/`write-all`; Block mit Eintrag; Anker/Alias. Empfehlung: nur `{}` besteht; alles andere ist ein Befund mit der Form im Text. Jobebene: gelockert wird nur mit benannten Scopes (kein `write-all`); ein Fall je Form. **Entschieden:** `SPEC-063` Punkt 6 (oben nur `{}`, am Job `{}` oder Block mit `read`, `write`, `none`).
5. **Einstufung der Befunde.** Wann ist ein Befund von `actionlint` (inklusive `shellcheck`) ein Fehler: alle, oder eine Teilmenge (`-ignore` ist ein Weg an der Suppression-Regel vorbei, `AGENTS.md` §3.2 gilt dem Go-Code, das Muster `# Why:` ist der Maßstab). Empfehlung: alle sind Befund; eine Ausnahme nur in einer Konfiguration mit Begründung. **Entschieden:** `SPEC-063` Punkt 4 (jede Zeile ist ein Befund; keine Option, keine Konfigurationsdatei von `actionlint`).
6. **`timeout-minutes`.** Pflicht je Job mit `runs-on`; Obergrenze (Empfehlung: kleiner als der GitHub-Standard 360; die Zahl legt der Vertrag fest); ein Wert als Ausdruck (`${{ … }}`) oder 0; ein Job, der einen wiederverwendbaren Workflow aufruft (`uses:` auf Jobebene), darf `timeout-minutes` nicht tragen (GitHub lehnt es ab): zwei Fälle. **Entschieden:** `SPEC-063` Punkt 7 (1 bis 60 bei `runs-on:`; Aufruf-Job ausgenommen).
7. **Reihenfolge der Gründe.** Eine Datei kann mehrfach abgelehnt werden (kaputtes YAML, unpinned `uses:`, fehlende `permissions`). Empfehlung: Syntax zuerst; danach laufen die übrigen Gruppen alle, damit ein Lauf alle Befunde nennt; die Ausgabe ist je Datei sortiert. Ein Fall je Paar. **Entschieden:** `SPEC-063` Punkt 9 (Umgebung vor Befunden, sonst alle Gruppen in einem Lauf, feste Folge).
8. **Aufrufarten des Treibers** (`AGENTS.md` §3.12): `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` aus der Umgebung, einzeln und kombiniert; das Gate liefert ein Ergebnis, das der Stop-Hook über den Nachweis glaubt. Empfehlung: das Rezept enthält nur einen Docker-Aufruf und kein `$(MAKE)`, `-n`/`-q`/`-t` führen es nicht aus; `-i` schlägt durch (kein verschluckter Befund); je Flag ein Fall. **Entschieden:** `SPEC-063` Punkt 10 (Fälle `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS`; der Nachweis der Kette liegt bei `SPEC-052` Punkt 21).
9. **Aufruf aus einem Unterverzeichnis; Quelle read-only.** Das Gate urteilt über das Repo seines Arbeitsverzeichnisses; der Mount ist schreibgeschützt; Symlinks unter `.github/workflows/` (Ziel außerhalb): nicht gefolgt, als Befund genannt. Ein Fall je Form. **Entschieden:** `SPEC-063` Punkte 1, 3 und 10 (Symlink ist Befund, Mount `:ro`, Wurzel aus dem Skriptpfad).
10. **Gleichzeitige Läufe und Parallelität** (`make -j`): `GATE_CHECKS` läuft parallel; Container-Namen sind eindeutig oder ohne Namen (`--rm`, `docker run` ohne `--name`); Aufruf-Voraussetzung nach `SPEC-052` Punkt 8 benennen. **Entschieden:** `SPEC-063` Punkte 3 und 10 (kein `--name`, kein Hilfsverzeichnis mit festem Namen).
11. **Version und Digest von `actionlint`** und Herkunft des Images (gepinnt per Digest, `docker.io` oder `ghcr.io`); wo der Pin steht (`harness/mk/workflows.mk`); Verhalten bei fehlendem Image (Pull-Fehler: Exit ungleich 0, nicht „keine Befunde“). **Entschieden:** `SPEC-063` Punkt 3 (`ghcr.io/rhysd/actionlint@sha256:…`; Pull-Fehler Exit 2; Tag und Digest ermittelt der Implementer).
12. **Zusammenspiel mit `make docs-check`.** Geben beide denselben Pin-Befund, ist es ein Befund in zwei Gates; Empfehlung: die Pins prüft allein `docs-check` (Modul), `workflows-check` prüft sie nicht ein zweites Mal; `SPEC-063` nennt es, die Gegenprobe zeigt es mit einer `docs-check`-Kopie. **Entschieden:** `SPEC-063` Punkt 8. **Abweichung vom ursprünglichen Plan (in §1, §2 und §3 nachgezogen):** das Modul `workflows` meldet ohne Datei und ohne `uses:` einen Befund (leere Prüfmenge); der Eintrag in `.d-check.yml` kann deshalb nicht vor der ersten Workflow-Datei stehen. Das Gate prüft nur, dass er steht, sobald eine Datei existiert; der Eintrag kommt mit `ci.yml` im Slice `slice-harness-ci`.

**Fragen an den Nutzer** (jede mit Empfehlung des Planners; eine Empfehlung ist keine Entscheidung):

1. **Schnitt und Reihenfolge.** Beauftragt waren zwei Slices (`slice-harness-ci`, `slice-harness-ci-freshness`) mit je Marke 900; die Schätzung ungeteilt liegt bei rund 1050 bis 1650 (CI) und 900 bis 1400 (Freshness). Nach der Regel „schneiden, wenn die Schätzung die Marke überschreitet“ sind es fünf: `slice-harness-ci-workflows` (dieser), `slice-harness-ci`, `slice-harness-ci-gegenproben`, `slice-harness-ci-freshness`, `slice-harness-ci-freshness-images`. Empfehlung: so, in dieser Reihenfolge. Alternative: dieser Slice zusammen mit `slice-harness-ci` (Schätzung 1000 bis 1400, darüber).
2. **Pins durch das Modul `workflows` von d-check oder durch eigenes Skript.** **Entschieden:** Modul, mit [ADR-0050](../../adr/0050-ci-auf-github-actions.md) (Accepted 2026-10-18) und `SPEC-063` Punkt 8; kein zweiter Parser.
3. **Aktualität und Gültigkeit der Action-SHA** (Netz). Empfehlung: als dritte Achse in `slice-harness-ci-freshness-images` oder als eigener Slice danach; heute nimmt keiner an.

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

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (`AGENTS.md` §3.10) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (`AGENTS.md` §3.12) — das Gate ist ein neuer Vertrag, die Randformen stehen in §6 vor dem Code; `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` — die Gegenprobe steht in der Container-Liste, nicht in der Kette; `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` — der Pin-Befund hat einen Ort (Randform 12). Kein Eintrag erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Der Slice trägt keine Sendung eines Gebers. Liefer-Punkte 3 (Gate | Gegenprobe | Gate-Index). Schichten 1 (Harness: `harness/mk`, `tools/harness`, `harness/vertraege`, `harness/README.md`); Planung, Vertrag und ADR sind keine Schicht, Wartungs-Doku und Handbuch bleiben unberührt. Schätzung 455 bis 690 (§4), unter der Marke.

**Schätzung des Architect nach dem Vertrag (2026-10-18; Metrik von §4, `SPEC-063` und ADR vor dem Anspruch geschrieben und damit nicht gezählt):** `workflows-check.sh` 150 bis 190 (Prüfbereich, Aufruf von `actionlint`, Lesegrammatik, Rechte, Zeitgrenzen, Bindung, Ausgabe), `workflows.mk` 25, Gegenprobe das Doppelte des Produktcodes, 350 bis 450 (rund 75 Fälle, 40 Mutationen, drei echte Fälle), `harness/README.md` 20 bis 30, zusammen 545 bis 695, unter der Marke 900. Die Spanne des Planners (455 bis 690) enthielt den Vertrag; ohne ihn und mit der Gegenprobe nach den Punkten von `SPEC-063` liegt sie ähnlich. Die Gruppe (c) der Pins (`.d-check.yml`) entfällt aus diesem Slice (siehe §6, Randform 12) und mindert die Schätzung um rund 25.

**Nachzählung nach dem Vertrag (Planner, 2026-10-18; `AGENTS.md` §3.13):** Liefer-Punkte 3 (Gate | Gegenprobe | Gate-Index), Schichten 1; die Zeile `.d-check.yml` entfällt (Sendung an `slice-harness-ci`, dort gezählt), die Schätzung des Architect (545 bis 695, abzüglich rund 25) liegt unter der Marke 900.
