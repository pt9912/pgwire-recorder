# Slice slice-harness-ci-workflows: Gate der Workflows — Pins, Rechte und Zeitgrenzen der Dateien unter `.github/workflows/`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-erster-release, erste Stelle der Reihenfolge, vor `slice-harness-ci` und vor `slice-erster-release-veroeffentlichung` (Auftrag des Nutzers vom 2026-10-18: die CI gehört in welle-erster-release, vor die Veröffentlichung; Entscheidung des Nutzers vom 2026-10-18: der kritische Pfad zum ersten Release besteht nur aus diesem Slice, `slice-harness-ci`, `slice-erster-release-veroeffentlichung` und `slice-erster-release-freigabe`, die drei Nachtläufe folgen in [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md)). Start-Trigger weicht von der Welle ab: `slice-harness-gegenproben-aus-gates` in `done/` (die Namens- und Zuordnungsregel von `SPEC-060` Teil 1 legt die Gegenprobe dieses Slice in `GEGENPROBEN_CONTAINER`; ohne die Listen hinge sie an `GATE_CHECKS` und müsste später umziehen), dazu der Start-Trigger der Welle (§4); `slice-harness-gegenproben-aus-gates-fahrstelle-umgebung` ist keine Bedingung (Entscheidung des Nutzers vom 2026-10-18).

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
| `tools/harness/workflows-check-gegenprobe.sh` (Ziel im Fragment `harness/mk/workflows.mk`) | neu | Gegenprobe in Temp-Bäumen und Kopien, ohne Abdeckungs-Deklaration (`SPEC-063` Punkt 12); Eintrag in `GEGENPROBEN_CONTAINER` |
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

1. **Kein Verzeichnis, leeres Verzeichnis, keine Datei.** Empfehlung: Exit 0 mit einer Zeile „keine Workflows“, nicht still; ein Fall je Form. Gate-Wahrheit „grün über leerem Prüfbereich“ benennt `SPEC-063` als Grenze (`LH-QA-01`-Muster wie `history-range-guard`). **Entschieden:** `SPEC-063` Punkt 2 (Exit 0 mit Meldung, ohne Docker; ein Verzeichnis nur mit Symlinks ist Befund `symlink`, Exit 1, Fall `p2-nur-symlink`; Grenze in Punkt 11).
2. **Dateiendungen und Ort.** GitHub liest `*.yml` und `*.yaml` unmittelbar unter `.github/workflows/`, keine Unterverzeichnisse. Empfehlung: beide Endungen, Unterverzeichnisse und andere Endungen nicht geprüft, aber mit Hinweis genannt; ein Fall je Endung und je Tiefe. **Entschieden:** `SPEC-063` Punkt 1 (beide Endungen, kein Unterverzeichnis, Hinweis).
3. **Weitere Orte.** `.github/actions/*/action.yml` (zusammengesetzte Actions), `.github/dependabot.yml`: nicht geprüft; nennt der Vertrag als Grenze, weil das Repo keine führt. **Entschieden:** `SPEC-063` Punkt 11 (nicht gelesen, Grenze).
4. **Formen von `permissions`.** Oberste Ebene: fehlt; `permissions: {}`; leerer Block; `read-all`/`write-all`; Block mit Eintrag; Anker/Alias. Empfehlung: nur `{}` besteht; alles andere ist ein Befund mit der Form im Text. Jobebene: gelockert wird nur mit benannten Scopes (kein `write-all`); ein Fall je Form. **Entschieden:** `SPEC-063` Punkt 6 (oben nur `{}`, am Job `{}` oder Block mit `read`, `write`, `none`).
5. **Einstufung der Befunde.** Wann ist ein Befund von `actionlint` (inklusive `shellcheck`) ein Fehler: alle, oder eine Teilmenge (`-ignore` ist ein Weg an der Suppression-Regel vorbei, `AGENTS.md` §3.2 gilt dem Go-Code, das Muster `# Why:` ist der Maßstab). Empfehlung: alle sind Befund; eine Ausnahme nur in einer Konfiguration mit Begründung. **Entschieden:** `SPEC-063` Punkt 4 (jede Zeile ist ein Befund; keine Option, keine Konfigurationsdatei von `actionlint`).
6. **`timeout-minutes`.** Pflicht je Job mit `runs-on`; Obergrenze (Empfehlung: kleiner als der GitHub-Standard 360; die Zahl legt der Vertrag fest); ein Wert als Ausdruck (`${{ … }}`) oder 0; ein Job, der einen wiederverwendbaren Workflow aufruft (`uses:` auf Jobebene), darf `timeout-minutes` nicht tragen (GitHub lehnt es ab): zwei Fälle. **Entschieden:** `SPEC-063` Punkt 7 (1 bis 60 bei `runs-on:`; Aufruf-Job ausgenommen).
7. **Reihenfolge der Gründe.** Eine Datei kann mehrfach abgelehnt werden (kaputtes YAML, unpinned `uses:`, fehlende `permissions`). Empfehlung: Syntax zuerst; danach laufen die übrigen Gruppen alle, damit ein Lauf alle Befunde nennt; die Ausgabe ist je Datei sortiert. Ein Fall je Paar. **Entschieden:** `SPEC-063` Punkt 9 (Umgebung vor Befunden, sonst alle Gruppen in einem Lauf, feste Folge).
8. **Aufrufarten des Treibers** (`AGENTS.md` §3.12): `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` aus der Umgebung, einzeln und kombiniert; das Gate liefert ein Ergebnis, das der Stop-Hook über den Nachweis glaubt. Empfehlung: das Rezept enthält nur einen Docker-Aufruf und kein `$(MAKE)`, `-n`/`-q`/`-t` führen es nicht aus; `-i` schlägt durch (kein verschluckter Befund); je Flag ein Fall. **Entschieden:** `SPEC-063` Punkt 10 (Fälle `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS`; der Nachweis der Kette liegt bei `SPEC-052` Punkt 21).
9. **Aufruf aus einem Unterverzeichnis; Quelle read-only.** Das Gate urteilt über das Repo seines Arbeitsverzeichnisses; der Mount ist schreibgeschützt; Symlinks unter `.github/workflows/` (Ziel außerhalb): nicht gefolgt, als Befund genannt. Ein Fall je Form. **Entschieden:** `SPEC-063` Punkte 1, 3 und 10 (Symlink ist Befund, Mount `:ro`, Wurzel aus dem Skriptpfad).
10. **Gleichzeitige Läufe und Parallelität** (`make -j`): `GATE_CHECKS` läuft parallel; Container-Namen sind eindeutig oder ohne Namen (`--rm`, `docker run` ohne `--name`); Aufruf-Voraussetzung nach `SPEC-052` Punkt 8 benennen. **Entschieden:** `SPEC-063` Punkte 3 und 10 (kein `--name`, kein Hilfsverzeichnis mit festem Namen).
11. **Version und Digest von `actionlint`** und Herkunft des Images (gepinnt per Digest, `docker.io` oder `ghcr.io`); wo der Pin steht (`harness/mk/workflows.mk`); Verhalten bei fehlendem Image (Pull-Fehler: Exit ungleich 0, nicht „keine Befunde“). **Entschieden:** `SPEC-063` Punkt 3 (`rhysd/actionlint@sha256:<Index-Digest>` vom Docker Hub (ghcr.io ist anonym nicht ziehbar); Pull-Fehler Exit 2; Exit 3 mit Rechtehinweis; Tag und Digest ermittelt der Implementer).
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

### Belege des Implementers

Stand: Zweig `impl-ci-workflows`, letzter Code-Commit `0c30b767` (Gate `7a66cc39`, Gegenprobe `a207e17c`, Gate-Index `0c30b767`), im Dev-Container pgr-dev3 (Docker-in-Docker, bash 5, mawk). Basis der Größenmetrik: `a98d0b9e`.

**Randform-Rückgabe vor dem Code** (Image-Herkunft und Rechte des Nicht-root-Images): vom Architect in `SPEC-063` Punkte 3 und 12 entschieden (Commit `2890e31f`); das Fragment trägt `rhysd/actionlint@sha256:887a259a…147e9` (Index-Digest von Tag 1.7.7, Docker Hub).

**Läufe am Endstand (alle im Container):**

| Lauf | Ergebnis |
|---|---|
| `make workflows-check-gegenprobe` | grün, 81 Läufe (unveränderte Kopie und 80 Mutationen), 86 Fälle, davon 3 echte gegen die Images von `actionlint` und d-check; rund 35 s mit `GP_PARALLEL=6` |
| `make gates-schnell-gegenprobe` | grün (Ausnahmeliste deckt `GATE_CHECKS`, Namensregel von `SPEC-060` Punkt 2 am Arbeitsbaum mit der neuen Gegenprobe) |
| `make gates` | Exit 0; `workflows-check` meldet im Arbeitsbaum „keine Workflow-Dateien unter .github/workflows“ (Grenze des Vertrags, Punkt 2 und 11) |
| `make docs-check` | 0 Befund(e) nach dem Eintrag in `harness/README.md` |
| echter Lauf des Gates | an einer Datei mit Ausdrucksfehler, `permissions: read-all`, `timeout-minutes: 99`, einem Job ohne Zeitgrenze und `contents: admin`: 10 Befunde (4 von `actionlint`, 5 eigene, 1 Bindung), Exit 1 |
| d-check (`real-pins`) | `uses: actions/checkout@v4` ⇒ `uses-pin-missing`, voller SHA ohne Kommentar ⇒ `uses-pin-untagged`, mit Tag-Kommentar ⇒ 0 Befunde (Image `DCHECK_REF`, v0.86.1) |

**Zusage · Mutation · roter Test.** Jede Mutation ist ein Textersatz in einer Kopie von Skript oder Fragment (`mut` in `tools/harness/workflows-check-gegenprobe.sh`, Name `m-…` ist der Schlüssel; `GP_MUT=<Name>` fährt sie einzeln, `GP_ENTDECKE=1` nennt die roten Fälle). Rot sind genau die Fälle der letzten Spalte; die Mutation ist von der Gegenprobe selbst erwartet, ein Unterschied zwischen erwarteter und gefundener Menge macht die Gegenprobe rot. „breit“ heißt: die Mutation färbt jeden Fall, der die Zusage nebenbei braucht; die Gegenprobe prüft dort, dass der benannte Kern rot ist und die ausdrücklich genannten Fälle grün bleiben. Fälle anderer Punkte stehen in der Spalte, wo sie die Regel als Auslöser eines Befunds benutzen (zum Beispiel `p1-yml` benutzt `permissions-fehlt`, wie der Vertrag es in Punkt 1 vorgibt).

| Zusage | Mutation | Name | rote Fälle |
|---|---|---|---|
| P1 Endung `.yml` wird geprüft | Prüfbereich nur `.yaml` (breit: jeder Fall mit `.yml`-Datei) | `m-p1-nur-yaml` | `p1-symlink` `p1-unterverzeichnis` `p1-yml` `p10-i` `p10-rezept` `p10-unterverzeichnis` … (60 Fälle) |
| P1 Endung `.yaml` wird geprüft | Prüfbereich nur `.yml` | `m-p1-nur-yml` | `p1-yaml` `p3-e6` `p9-gruen` |
| P1/P3/P9 Reihenfolge der Namen | `sort -rz` statt `sort -z` | `m-p1-reihenfolge` | `p3-e6` `p9-reihenfolge` |
| P1 kein Unterverzeichnis | rekursive Suche | `m-p1-rekursiv` | `p1-unterverzeichnis` |
| P1 Symlink ist Befund, nicht gefolgt | Symlink wie reguläre Datei | `m-p1-symlink-gefolgt` | `p1-symlink` |
| P1 Endung ist das Ende des Namens | Endung als Teilstring (`*.yml*`) | `m-p1-teilstring` | `p1-andere-endung` |
| P2 F1 kein Verzeichnis: Exit 0 mit Meldung | F1 endet mit Exit 2 | `m-p2-f1-exit2` | `p2-kein-verzeichnis` |
| P2 F2 leeres Verzeichnis: Meldung | F2 endet still | `m-p2-f2-still` | `p2-leer` `p2-ohne-4` |
| P2 F3 nur andere Dateien: kein `docker` | F3 geht weiter und ruft `docker` | `m-p2-f3-docker` | `p1-andere-endung` `p2-nur-andere` `p2-ohne-8` |
| P2 ohne die Prüfung von Punkt 4 | Punkt 4 läuft vor der Meldung | `m-p2-punkt4` | `p2-ohne-4` |
| P2 ohne die Prüfung von Punkt 8 | Punkt 8 läuft vor der Meldung | `m-p2-punkt8` | `p1-andere-endung` `p2-kein-verzeichnis` `p2-leer` `p2-nur-andere` `p2-ohne-4` `p2-ohne-8` |
| P3 Digest hat 64 Hexziffern | 63 Hexziffern | `m-p3-digest-kurz` | `p3-digest-form` |
| P3 Repository `rhysd/actionlint` | `ghcr.io/rhysd/actionlint` | `m-p3-digest-repository` | `p3-digest-form` |
| P3 Digest-Form `@sha256:` | Tag statt Digest | `m-p3-digest-tag` | `p3-digest-form` |
| P3 fehlendes `docker` ist Exit 2 | Exit 127 gilt als keine Befunde | `m-p3-docker-fehlt-als-0` | `p3-kein-docker` `p9-umgebung-vor-befund` |
| P3 E1 `--rm` | `--rm` fehlt | `m-p3-e1` | `p3-e1` |
| P3 E2 `--network none` | `--network none` fehlt | `m-p3-e2` | `p3-e2` |
| P3 E3 Mount `:ro` | `:ro` fehlt | `m-p3-e3` | `p3-e3` |
| P3 E4 Image aus `WORKFLOWS_ACTIONLINT` | festes Image statt der Variable | `m-p3-e4` | `p3-e4` |
| P3 E5 `-oneline` | `-oneline` fehlt | `m-p3-e5` | `p3-e5` |
| P3 E6 Dateiliste vollständig | nur die erste Datei | `m-p3-e6` | `p3-e6` |
| P3 E7 kein `--name` | `--name` kommt hinzu | `m-p3-e7` | `p3-e7` |
| P3 E8 kein `-ignore` | `-ignore` kommt hinzu | `m-p3-e8-ignore` | `p3-e8` |
| P3 E8 kein `-pyflakes` | `-pyflakes=` kommt hinzu | `m-p3-e8-pyflakes` | `p3-e8` |
| P3 E8 kein `-shellcheck` | `-shellcheck=` kommt hinzu | `m-p3-e8-shellcheck` | `p3-e8` |
| P3 Exit 1 heißt Befunde | Exit 1 gilt als 0 | `m-p3-exit1-als-0` | `p10-i` `p10-unterverzeichnis` `p3-exit1` `p3-exit1-leer` `p9-reihenfolge` `real-ausdruck` |
| P3 Exit 1 ohne Zeile ist Befund | Zweig entfällt (Exit 1 ohne Zeile gilt als grün) | `m-p3-exit1-leer-angenommen` | `p3-exit1-leer` |
| P3 Hinweis nur bei Exit 3 | Hinweis bei jedem Exit mit `permission denied` | `m-p3-hinweis-bei-jedem-exit` | `p3-exit2` |
| P3 Hinweis nur mit `permission denied` | Hinweis bei jedem Exit 3 | `m-p3-hinweis-bei-jedem-exit3` | `p3-exit3` |
| P3 Hinweis bei Exit 3 mit `permission denied` | Hinweis entfällt | `m-p3-hinweis-fehlt` | `p3-exit3-rechte` |
| P3 leeres Image: Exit 2 | leeres Image wird gerufen | `m-p3-image-leer` | `p3-image-leer` |
| P3 anderer Exit heißt Exit 2 | Exit 2, 3, 125, 127 gelten als 1 | `m-p3-jeder-exit-als-1` | `p3-exit125` `p3-exit2` `p3-exit3` `p3-exit3-rechte` `p3-kein-docker` `p9-umgebung-vor-befund` |
| P3 genau ein Aufruf | zweiter Aufruf von `docker` | `m-p3-zweimal` | `p10-i` `p10-rezept` `p3-einmal` |
| P4 `actionlint.yaml` ist Befund | Endung `.yaml` nicht geprüft | `m-p4-yaml-nicht` | `p4-yaml` |
| P4 `actionlint.yml` ist Befund | Endung `.yml` nicht geprüft | `m-p4-yml-nicht` | `p4-yml` |
| P5 E2 erster Job-Schlüssel mit Einzug 2 | Einzug nicht geprüft | `m-p5-einzug` | `p5-einzug` |
| P5 E2 späterer Job-Schlüssel mit Einzug 1 oder 3 | Einzug 1/3 nicht geprüft | `m-p5-einzug-spaeter` | `p5-einzug-spaeter` |
| P5 Schritt-Eigenschaft ist keine Job-Eigenschaft | Suche auf jedem Einzug ab 4 | `m-p5-jeder-einzug` | `p5-schritt-ist-kein-job` `p6-job-wert` |
| P5 E3 `jobs:` allein in der Zeile | Wert hinter `jobs:` angenommen | `m-p5-jobs-flow` | `p5-jobs-flow` |
| P5 E1 Tabulator im Einzug | Tabulator still angenommen | `m-p5-tab` | `p5-tab` |
| P6 oben fehlt `permissions`: Befund | fehlendes `permissions` angenommen | `m-p6-fehlt-angenommen` | `p1-yaml` `p1-yml` `p6-fehlt` `p9-alle-befunde` `p9-reihenfolge` |
| P6 Flow-Mapping mit Einträgen am Job ist Befund | Flow-Mapping am Job angenommen | `m-p6-flow-am-job` | `p6-job-flow-eintrag` |
| P6 oben nur `{}` | jeder Wert oben angenommen | `m-p6-jeder-wert-oben` | `p6-block` `p6-flow-eintrag` `p6-read-all` `p6-write-all` |
| P6 `permissions` am Job nicht verlangt | `permissions` am Job verlangt (breit) | `m-p6-job-pflicht` | `p1-symlink` `p1-unterverzeichnis` `p1-yaml` `p1-yml` `p10-rezept` `p3-e1` … (47 Fälle) |
| P6 Skalar am Job ist Befund | Skalar am Job angenommen | `m-p6-skalar-am-job` | `p6-job-read-all` `p6-job-write-all` `p9-alle-befunde` |
| P6 Wert im Block nur read, write, none | unbekannter Wert angenommen | `m-p6-unbekannter-wert` | `p6-job-wert` |
| P7 Ausdruck ist Befund | Ausdruck angenommen | `m-p7-ausdruck-angenommen` | `p7-ausdruck` |
| P7 Job ohne `runs-on` ausgenommen | Job ohne `runs-on` geprüft | `m-p7-job-ohne-runs-on` | `p7-uses-job-ok` |
| P7 `30.5` ist Befund | Nicht-Ganzzahl angenommen | `m-p7-nicht-ganzzahl-angenommen` | `p7-nicht-ganzzahl` |
| P7 0 ist Befund | 0 angenommen | `m-p7-null-angenommen` | `p7-0` |
| P7 Obergrenze unter 360 | Obergrenze 360 | `m-p7-obergrenze-360` | `p7-61` |
| P7 Obergrenze 60 | Obergrenze 61 | `m-p7-obergrenze-61` | `p7-61` |
| P7 Zeitgrenze Pflicht | Pflicht entfällt | `m-p7-pflicht-entfaellt` | `p5-schritt-ist-kein-job` `p7-fehlt` `p9-alle-befunde` `p9-reihenfolge` |
| P7 Untergrenze 1 | Untergrenze 2 | `m-p7-untergrenze` | `p7-1-ok` |
| P7 Wort ist Befund | Wort angenommen | `m-p7-wort-angenommen` | `p7-wort` |
| P8 (a) `modules:` führt `workflows` | (a) nicht geprüft | `m-p8-a` | `p8-modules` |
| P8 `"…"` zulässig | doppelte Anführungszeichen nicht zugelassen | `m-p8-anfuehrung-doppelt` | `p8-ok-doppelt` |
| P8 `'…'` zulässig | einfache Anführungszeichen nicht zugelassen | `m-p8-anfuehrung-einfach` | `p8-ok-einfach` |
| P8 (b) Zeile `workflows:` | (b) nicht geprüft | `m-p8-b` | `p8-block` |
| P8 (c) `dir: .github/workflows` | (c) nicht geprüft | `m-p8-c` | `p8-dir` `p8-dir-schraeg` |
| P8 fehlende `.d-check.yml` ist Befund | fehlende Datei angenommen | `m-p8-datei-fehlt-angenommen` | `p8-ohne-datei` `p9-reihenfolge` |
| P8 das Gate meldet Pins nicht selbst | Skript meldet unpinnte Referenz | `m-p8-pins-selbst` | `p8-pins-nicht-hier` |
| P8 ohne schließenden Schrägstrich | Schrägstrich am Ende zugelassen | `m-p8-schraegstrich` | `p8-dir-schraeg` |
| P9 `actionlint` zuerst | Datei-Befunde vor `actionlint` | `m-p9-datei-vor-al` | `p9-reihenfolge` |
| P9 Umgebungsfehler vor Befunden | Gruppen (b), (c) laufen vor dem Umgebungsfehler | `m-p9-eigene-vor-umgebung` | `p9-umgebung-vor-befund` |
| P9 alle Befunde in einem Lauf | Ausgabe endet beim ersten Befund | `m-p9-erster-befund` | `p3-exit1` `p9-alle-befunde` `p9-reihenfolge` |
| P9 Befunde ohne Datei zuletzt | Befunde ohne Datei vor den Datei-Befunden | `m-p9-global-vor-datei` | `p9-reihenfolge` |
| P9 nicht lesbare Datei: Exit 2 | Datei still übergangen | `m-p9-nicht-lesbar` | `p9-nicht-lesbar` |
| P9 Befunde nach Zeile | Sortierung entfällt | `m-p9-ohne-sortierung` | `p9-alle-befunde` |
| P9 Summe beider Quellen | Summe zählt nur `actionlint` (breit) | `m-p9-summe-nur-al` | `p1-symlink` `p1-yaml` `p1-yml` `p4-yaml` `p4-yml` `p5-einzug` … (32 Fälle) |
| P9 Summe beider Quellen | Summe ohne `actionlint` | `m-p9-summe-ohne-al` | `p10-unterverzeichnis` `p3-exit1` `p3-exit1-leer` `p9-reihenfolge` |
| P9 Umgebungsfehler ist Exit 2 | Exit 1 | `m-p9-umgebung-exit1` | `p3-exit125` `p3-exit2` `p3-exit3` `p3-exit3-rechte` `p3-kein-docker` `p9-umgebung-vor-befund` |
| Gate hängt an `GATE_CHECKS` | Zeile entfällt | `m-p10-kette-gate` | `p10-kette-gate` |
| Gegenprobe steht in `GEGENPROBEN_CONTAINER` | Zeile entfällt | `m-p10-kette-gegenprobe` | `p10-kette-gegenprobe` |
| P10 Rezept ohne `+` (`-n`, `-q`, `-t`, `MAKEFLAGS`) | `+` vor dem Rezept | `m-p10-plus` | `p10-makeflags` `p10-n` `p10-q` `p10-t` |
| P10 schreibt nichts in den Baum | Skript legt eine Datei im Baum an | `m-p10-schreibt` | `p10-schreibt-nichts` |
| P10 Wurzel aus Skriptpfad | Wurzel ist das Arbeitsverzeichnis | `m-p10-wurzel-ist-arbeitsverzeichnis` | `p10-unterverzeichnis` |
| P10 Rezept ruft das Skript genau einmal | Rezept doppelt | `m-p10-zweimal` | `p10-i` `p10-rezept` |
| P12 echtes Image findet den Ausdrucksfehler | Image durch Attrappe mit Exit 0 ersetzt (Umgebung) | `m-real-attrappe-image` | `real-ausdruck` |
| P12 d-check meldet Pins mit `workflows` in `modules:` | `modules:` ohne `workflows` (Umgebung) | `m-real-modules` | `real-pins` |

**Grüne Mutanten:** keiner. Alle 80 Mutationen färben mindestens einen Fall; es gibt deshalb keinen Mutanten, der als äquivalent einzuordnen wäre.

**Befundzeilen vorher und nachher:** kein Bestand (das Gate liegt vor der ersten Workflow-Datei; `.github/workflows/` fehlt im Arbeitsbaum).

**Metrik gegen die Marke von §4** (Summe aus Einfügungen und Löschungen von `git diff --numstat a98d0b9e..0c30b767 -- .github tools harness Makefile`): 830 (818 Einfügungen, 12 Löschungen), davon 22 aus den Vertragsdateien des Architect (`harness/vertraege`); ohne sie 808 (`workflows-check-gegenprobe.sh` 542, `workflows-check.sh` 240, `workflows.mk` 24, `harness/README.md` 2). Meldeschwelle 700 überschritten, Marke 900 nicht. Die Schätzung des Architect (545 bis 695) liegt unter der Messung; die Gegenprobe trägt zwei Drittel des Diffs.

**Offener Punkt für den Architect** (nicht entschieden, nicht getestet): Besteht der Prüfbereich nur aus einem Symlink mit Endung `.yml`/`.yaml`, nennt Punkt 1 einen Befund `symlink`, Punkt 2 aber „nur Dateien, die Punkt 1 nicht prüft“ als Fall ohne Prüfbereich. Das Skript meldet den Symlink als Befund (Exit 1), ruft `docker` nicht (keine Dateiliste) und prüft Punkt 8 nicht (keine reguläre Datei).

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
