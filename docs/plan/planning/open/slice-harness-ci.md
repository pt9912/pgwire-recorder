# Slice slice-harness-ci: Die CI auf GitHub — `ci.yml` fährt `make gates` und prüft die Commit-Range

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-erster-release, zweite Stelle der Reihenfolge, direkt hinter `slice-harness-ci-workflows` und vor `slice-erster-release-veroeffentlichung` (Auftrag des Nutzers vom 2026-10-18; Entscheidung des Nutzers vom 2026-10-18: die Nachtläufe `slice-harness-ci-gegenproben`, `-ci-freshness` und `-ci-freshness-images` folgen nach dem Release in [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md)). Start-Trigger weicht von der Welle ab: `slice-harness-ci-workflows` und `slice-harness-gegenproben-aus-gates` in `done/` (§4); `slice-harness-gegenproben-aus-gates-fahrstelle-umgebung` ist keine Bedingung (Entscheidung des Nutzers vom 2026-10-18).

**Bezug:** Kein Anspruch des Produkts (`LH-*`): der Slice ändert die Prüfung der Integration, nicht das Binary. Aktive ADRs: [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Docker-only), [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Gate-Kette; der Stop-Hook bleibt), [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) und [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md) (die lokale Strenge der Kennung im Commit, die die CI nicht wiederholt, Randform 11). Die ADR dieses Slice schreibt der Architect vor dem Code (Entscheidung des Nutzers).

**Berührte Spec-Stellen:** `SPEC-064` (neu; Nummer vergibt der Architect, nächste freie; Vertrag in `harness/vertraege/`, vor dem ersten Code-Commit) · `SPEC-063` (gelesen: Gate der Workflows; `ci.yml` besteht es) · [`SPEC-050`](../../../../harness/vertraege/SPEC-050-struktur-kennungen-im-commit.md) (gelesen: Kennungen im Commit) · [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) Punkte 1 bis 4 und 21 (gelesen: Nachweis des vollen Laufs, Aufrufarten).

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

**Ziel:** Jede Integration (Pull-Request, Push auf einen Zweig) läuft auf GitHub durch `make gates` und durch die Prüfung ihrer Commit-Range, unabhängig vom Klon und vom lokalen Hook: `.github/workflows/ci.yml` (nach dem Vorbild von d-check) mit `pull_request` und `push` auf Zweigen, `tags-ignore` (Tags gehören `release.yml`), `permissions: {}` oben und `contents: read` allein im Job, `timeout-minutes`, `ubuntu-latest`, Actions per SHA mit Tag-Kommentar gepinnt, einem Job mit festem Namen (der Name ist der Status, den ein Branch-Schutz verlangt, Randform 8). Die Range-Prüfung ruft `make doc-commits` (jede Commit-Message nennt eine Kennung) und `make doc-immutable` (eine angenommene ADR bleibt unverändert) mit einer `RANGE`, die ein Skript `tools/harness/ci-range.sh` aus dem Ereignis berechnet (nicht im YAML, wo sie ungeprüft bliebe); eine leere oder unauflösbare Range ist rot, nicht grün (`history-range-guard`). Gegenproben prüfen die Zusagen der Datei und die Berechnung der Range. **Übernommen von `slice-harness-ci-workflows` (Sendung des Architect, 2026-10-18):** (a) der Eintrag für die Pins in `.d-check.yml` (`workflows` in `modules:` und `workflows.dir`, `SPEC-063` Punkt 8): das Modul meldet ohne Datei und ohne `uses:` einen Befund, der Eintrag kommt darum mit `ci.yml`; (b) die Blöcke `commits` und `vcs` in `.d-check.yml` (`SPEC-064` Punkte 10 bis 12): ohne sie prüfen `make doc-commits` und `make doc-immutable` in diesem Repo nichts, und die Gegenprobe zeigt, dass sie wirken (`SPEC-064` Punkt 14).

**Was im Repo vorliegt und was fehlt** (Bestand vom Planner am 2026-10-18 gelesen): Es liegen `make doc-commits` und `make doc-immutable` (`d-check.mk`, mit `RANGE`, über das Fragment `harness/mk/doc-gate.mk` mit dem Vorlauf `history-range-guard`), `tools/harness/commit-msg-traceability.sh` (eine Message-Datei, Mengen `ADR-`, `LH-`, `MR-`, `slice-`) und der lokale Träger `.githooks/commit-msg` (prüft zusätzlich, dass die Slice- oder Welle-Datei im Index liegt, [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md)). **Es fehlen:** jeder Workflow (`.github/` gibt es nicht), die Berechnung der Range aus dem Ereignis, ein Ziel, das beide Range-Prüfungen bündelt (nicht nötig, zwei Schritte genügen), und ein Sensor, dass `ci.yml` hält, was sie zusagt. Ein ADR-Immutable-Skript gibt es nicht; es ist das Modul `vcs` von d-check, das `make doc-immutable` fährt (Randform 11).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Gate der Workflows (Pins, Rechte, Zeitgrenzen, Syntax) — `slice-harness-ci-workflows` (Geber, vor diesem Slice); `ci.yml` besteht es und wird dadurch geprüft, dieser Slice baut keinen zweiten Prüfer.
- Der nächtliche Lauf der Gegenproben und die Installation von `u-boot` und `devcontainer` auf dem Runner — `slice-harness-ci-gegenproben`; `ci.yml` fährt `make gates`, das seit `slice-harness-gegenproben-aus-gates` keine Gegenprobe mehr fährt.
- Aktualität der Pins — `slice-harness-ci-freshness` und `-images`; ein Netzausfall bei Fremdbeständen darf diese CI nicht rot färben (d-check trennt aus demselben Grund).
- `release.yml`, Registry-Anmeldung, Veröffentlichung, `arm64` — `slice-erster-release-veroeffentlichung`; diese CI trägt keine Geheimnisse und schreibt nichts.
- Branch-Schutz, „required status check“ und die Sichtbarkeit des Repos — Einstellungen auf GitHub, aus dem Klon nicht auditierbar; der Slice nennt sie als Grenze und als Handlung des Betreibers in `harness/README.md`, ohne sie zu behaupten (`AGENTS.md` §3.11).
- Die Strenge des lokalen Hooks (Slice- oder Welle-Datei liegt im Index) in der Range — der Hook liest den Index des Arbeitsbaums, die Range liest Commits; die CI prüft die Anwesenheit einer Kennung, wie das Modul `commits` sie prüft. Bestand bleibt bewusst stehen, die Lücke steht in `SPEC-064` als akzeptiertes Negativ (Randform 11).
- Der Nachweis `.harness/state/` — die CI schreibt, übernimmt und prüft ihn nicht (Randform 9); `slice-harness-nachweis-uebernehmen` ist ein anderer Vorgang.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Workflow):** `.github/workflows/ci.yml` hält die Zusagen von `SPEC-064` (Auslöser, `tags-ignore`, Rechte, Zeitgrenze, `concurrency`, Checkout mit voller Historie und ohne gespeicherte Zugangsdaten, Schritte `make gates` und die Range-Prüfung, kein Geheimnis, keine Hochladung, fester Jobname) und besteht `make workflows-check` und `make docs-check`, und `.d-check.yml` trägt den Eintrag für die Pins (Sendung von `slice-harness-ci-workflows`, `SPEC-063` Punkt 8); jede Entscheidung je Randform aus §6 steht vor dem ersten Code-Commit in `SPEC-064`. Belegt in §7 mit der Nummer eines grünen Laufs auf GitHub und der gemessenen Dauer.
- [ ] **Liefer-Punkt 2 (Range):** `tools/harness/ci-range.sh` berechnet aus den Eingaben des Ereignisses (Ereignisart, Basis des Pull-Requests, Vorgänger-Stand des Pushes, Stand, Standard-Zweig) die `RANGE` für `make doc-commits` und `make doc-immutable`, je Randform aus §6 (Pull-Request, Push, neuer Zweig, Vorgänger nicht mehr im Klon, leere Range, Tiefe des Klons) mit der im Vertrag entschiedenen Antwort; ein Fehler endet mit Exit ungleich 0 und einer Meldung, nie mit einer Range, die nichts prüft; `.d-check.yml` trägt die Blöcke `commits` und `vcs` nach `SPEC-064` Punkte 11 und 12 (Sendung von `slice-harness-ci-workflows`).
- [ ] **Liefer-Punkt 3 (Gegenprobe):** `make ci-gegenprobe` (Name: Architect) zeigt je Zusage von `SPEC-064` und je Element ihrer Aufzählungen einen Fall und eine Mutation, die nur ihn färbt (`AGENTS.md` §3.10): an Kopien von `ci.yml` (Auslöser, Rechte, Zeitgrenze, `fetch-depth`, `persist-credentials`, Jobname, kein Geheimnis, Schritte in Reihenfolge) und an `ci-range.sh` (jede Form der Range); die unveränderte Kopie grün; `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` je ein Fall. Das Ziel steht in `GEGENPROBEN_CONTAINER`. `harness/README.md` §Sensors führt die Zeile; sie nennt Branch-Schutz als Handlung des Betreibers und die CI als Ergänzung des Stop-Hooks, nicht als Ersatz.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `harness/` und `tools/` ändern) und `make gates-schnell-gegenprobe` grün; `make ci-gegenprobe` grün auf dem Endstand.
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
| `.github/workflows/ci.yml` | neu | Auslöser, Rechte, Zeitgrenze, `make gates`, Range-Prüfung; Vorbild d-check `ci.yml` |
| `tools/harness/ci-range.sh` | neu | `RANGE` aus dem Ereignis; fail-closed |
| `harness/mk/ci-gegenprobe.mk`, `tools/harness/ci-gegenprobe.sh` (Namen: Architect) | neu | Gegenprobe der Zusagen von `ci.yml` und der Range; Eintrag in `GEGENPROBEN_CONTAINER` |
| `.d-check.yml` | update | Eintrag der Pins (`SPEC-063` Punkt 8) und Blöcke `commits`, `vcs` (`SPEC-064` Punkte 11, 12); Sendung von `slice-harness-ci-workflows` |
| `harness/vertraege/SPEC-064-…md` | neu (Architect) | Vertrag; vor dem Code |
| `harness/README.md` §Sensors | update | Zeile der Gegenprobe; Grenze Branch-Schutz; Verhältnis zum Stop-Hook |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-ci-workflows` und `slice-harness-gegenproben-aus-gates` liegen in `done/` (jener für das Gate der Workflows, dieser, weil `make gates` ohne ihn die Gegenproben fährt und die Zeitgrenze des Jobs nicht trägt), dazu der Start-Trigger von [welle-erster-release](../welle-erster-release.md), und der Architect hat `SPEC-064` und die ADR geschrieben, jede Randform aus §6 entschieden (`AGENTS.md` §3.12). Schritt 2 der Reihenfolge in §5 der Welle.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: mehr als drei Liefer-Punkte, oder der Diff ist nicht in einer Review-Sitzung prüfbar — zurück zur Zerlegung. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- .github tools harness .d-check.yml Makefile`, Tests und Gegenprobe eingeschlossen, liegt über 900; Basis: beim Anspruch (HEAD vor dem Anspruch, hier einzutragen). Schätzung des Planners in dieser Metrik (2026-10-18, ohne Basis): `ci.yml` 90 bis 130, `ci-range.sh` 60 bis 90, Gegenprobe das Doppelte des Produktcodes (300 bis 440), `SPEC-064` 140 bis 200, `harness/README.md` 15 bis 30, zusammen 605 bis 890. Die Spanne reicht an die Marke; Schnitt, falls sie eintritt: *Range* (Liefer-Punkt 2 mit seinem Teil der Gegenprobe) und *Workflow* (Punkt 1 und 3).
- `in-progress` → `open`: der Architect entscheidet, dass die CI ein anderer Träger als GitHub Actions trägt, oder das Repo hat keinen GitHub-Remote, auf dem ein Lauf belegt werden kann (§6, Risiko).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, ein grüner Lauf von `ci.yml` auf GitHub für den Endstand und ein roter Lauf mit einem Verstoß gegen die Range (Probe-Zweig, danach gelöscht) als Beleg in §7, `make gates` grün, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Ob `make gates` auf dem Runner läuft, zeigt nur ein Lauf.** Der Runner (Docker vorhanden, kein `u-boot`, kein `devcontainer`) muss alles tragen, was die Kette fährt (Docker-Netz der Integrationstests, `setsid`, Speicher, Zeit); `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert` und `BEO-REPO/nested-runtime-podman-oder-docker-offen` stehen dazu im Register. — **Ausgang (vorgesehen):** entfallen, wenn ein Lauf auf GitHub grün endet und die Dauer in §7 steht; sonst eingetreten → Rückführung nach `open` (§4), und `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert` erhält einen Beleg.
- **Die Zeitgrenze von 30 Minuten setzt voraus, dass der Kettenumbau die Laufzeit senkt.** — **Ausgang (vorgesehen):** entfallen, wenn die gemessene Dauer in §7 unter der Hälfte der Grenze liegt; sonst weiter offen → `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Vermerk) und die Grenze wandert per Vertrag, nicht still.
- **Rate-Limit von Docker Hub auf geteilten Runner-Adressen** (anonyme Pulls der gepinnten Images, darunter `postgres`). — **Ausgang (vorgesehen):** weiter offen → neues `BEO-REPO/ci-pull-rate-limit` bei der ersten Beobachtung; ohne Beobachtung entfallen mit der Probe-Messung in §7.
- **Die CI ist ohne Branch-Schutz ein Hinweis, kein Sperrriegel.** — **Ausgang (vorgesehen):** weiter offen → Beobachtung `BEO-REPO/ci-ohne-branch-schutz` bei der Closure (Handlung des Betreibers, aus dem Klon nicht prüfbar); der Wortlaut in `harness/README.md` nennt es.
- **Doppelte Läufe** (derselbe Commit löst `push` und `pull_request` aus). — **Ausgang (vorgesehen):** entfallen, wenn der Vertrag Auslöser und `concurrency` entscheidet (Randform 1 und 3) und die Gegenprobe es zeigt.
- **Die Größe liegt an der Marke** (605 bis 890, Marke 900). — **Ausgang (vorgesehen):** entfallen, wenn die Messung nach der Metrik unter 900 bleibt; darüber meldet der Implementer (§4 Rückführung *zu groß*).

**Randformen** (`AGENTS.md` §3.12) — **alle offen für den Architect**, zu entscheiden vor dem ersten Code-Commit in `SPEC-064` (Technik-Stratum, `harness/vertraege/`), mit ADR für Entscheidung und Gründe. Vorbereitet vom Planner am 2026-10-18 (Empfehlungen, kein Beschluss):

1. **Auslöser.** `pull_request` (Standardtypen) und `push` auf Zweige (`'**'`) mit `tags-ignore`. Ein Pull-Request aus demselben Repo löst beides aus, zwei Läufe je Commit. Optionen: (a) wie d-check, beides; (b) `push` nur auf den Standard-Zweig, `pull_request` für alle. Empfehlung: (a), solange ein Lauf nicht die Kosten dominiert; Messung in §7. Ob `workflow_dispatch` hinzukommt (Wiederanlauf ohne Commit): Empfehlung ja, ohne Eingaben; **der Nutzer hat am 2026-10-18 den Verzicht bestätigt.** **Entschieden:** `SPEC-064` Punkte 1 und 2 (beides, kein `workflow_dispatch`; Doppelläufe hingenommen).
2. **Pull-Requests aus Forks.** `pull_request` (nie `pull_request_target`), keine Geheimnisse in dieser Datei, `contents: read`: der Fork-Code läuft ohne Rechte. Ein Fall, dass `ci.yml` weder `pull_request_target` noch `secrets.` nennt. **Entschieden:** `SPEC-064` Punkt 3 (kein `pull_request_target`, keine Geheimnisse; der Text ist wörtlich festgelegt).
3. **Abbruch älterer Läufe.** `concurrency` mit Gruppe aus Workflow und Ref; Empfehlung: ältere Läufe abbrechen, außer auf dem Standard-Zweig (dort bleibt jeder Lauf als Befund). Der Wortlaut mit Ausdruck für den Standard-Zweig ist eine Zusage und braucht einen Fall. **Entschieden:** `SPEC-064` Punkt 2 (**kein** `concurrency`: ein Abbruch verwirft die Range eines Push; Abweichung von der Empfehlung des Planners, vom Nutzer am 2026-10-18 bestätigt).
4. **Zeitgrenze.** Job 30 Minuten (Vorgabe des Auftrags); kein Zeitlimit je Schritt. Überschreitung ist rot (GitHub bricht ab); die gemessene Dauer steht in §7. **Entschieden:** `SPEC-064` Punkt 1 (30 Minuten am Job).
5. **Cache der Docker-Layer.** `make` ruft `docker build` über das `Dockerfile` ([ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)); ein Cache über BuildKit-Aktionen verlangte einen anderen Treiber oder Aktionen von Dritten. Empfehlung: ohne Cache starten, Dauer messen; ein Cache ist ein eigener Slice mit Messung. Kein Cache heißt: jeder Lauf kalt, das Image-Pull zählt in die Zeit. **Entschieden:** `SPEC-064` Punkt 2 (kein Cache).
6. **Matrix.** Keine: ein Job, `ubuntu-latest`, `linux/amd64`. Die Binaries für andere Plattformen baut `make test-container` per Cross-Kompilierung; `linux/arm64` als Image prüft `release.yml` (`slice-erster-release-veroeffentlichung`). Ob `ubuntu-latest` ein fester Name (`ubuntu-24.04`) wird, damit ein Runner-Wechsel nicht still die Kette ändert: Empfehlung fest, Wechsel per Vertragsänderung. **Entschieden:** `SPEC-064` Punkte 1 und 2 (`ubuntu-24.04` fest, keine Matrix; Abweichung vom Auftrag, der `ubuntu-latest` nannte, vom Nutzer am 2026-10-18 bestätigt).
7. **Rechte.** `permissions: {}` oben; im Job nur `contents: read`; `persist-credentials: false` am Checkout (kein Token im `.git`, das ein Gate-Skript lesen könnte); Aktionen ausschließlich per SHA mit Tag-Kommentar (`make workflows-check`, `make docs-check`). Fremde Aktionen: nur `actions/checkout`; Empfehlung kein weiteres. **Entschieden:** `SPEC-064` Punkte 1 und 3 (Text wörtlich; nur `actions/checkout`).
8. **Was ist „grün“ für den Merge.** Der Jobname ist der Status-Check. Ändert er sich, verliert ein Branch-Schutz seine Bedingung still; Empfehlung: der Name steht im Vertrag, ein Fall der Gegenprobe hält ihn, und `harness/README.md` nennt ihn als den Namen für den Branch-Schutz. Die Einstellung selbst ist keine Aussage des Repos. **Entschieden:** `SPEC-064` Punkt 4 (Job `ci-gates`).
9. **Nachweis und Stop-Hook.** Die CI ergänzt den lokalen Stop-Hook (klonunabhängig, Range, `--no-verify`), sie ersetzt ihn nicht. `make gates` schreibt als letzten Schritt `.harness/state/gates-passed.dateiliste` in den Arbeitsbaum des Runners; das ist folgenlos, wird nicht hochgeladen, nicht zwischengespeichert, nicht committet und ersetzt keinen lokalen Nachweis. Empfehlung: so lassen, `make gates` unverändert rufen (kein zweites Gate-Ziel nur für die CI); der Vertrag nennt es, und ein Fall zeigt, dass `ci.yml` keinen Upload- oder Cache-Schritt trägt. **Entschieden:** `SPEC-064` Punkt 5.
10. **Ressourcen und Aufruf.** `make gates` mit oder ohne `-j`; Speicher und CPU des Runners (öffentliches Repo: größerer Runner als privates); `PGR_MEM` und `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert`. Empfehlung: `make gates` ohne `-j` starten, die Dauer messen; ein `-j` ist eine Messung, keine Annahme. Kein Wiederholen eines roten Laufs im Workflow (ein Wiederanlauf maskiert Flackern); Wiederanlauf von Hand über die Oberfläche von GitHub. **Entschieden:** `SPEC-064` Punkte 2 und 5 (`make gates` ohne `-j`, kein Wiederholen).
11. **Range und Kennung.** Berechnung: Pull-Request `Basis..Kopf` (der Kopf ist der letzte Commit des Zweigs, nicht der Merge-Commit von GitHub, der mit `Merge` beginnt und den Betreff-Test des Skripts übergeht); Push `Vorgänger..Stand`; neuer Zweig (Vorgänger nur Nullen) gegen den Standard-Zweig; Vorgänger nicht mehr im Klon (Force-Push); leere Range; erster Commit des Repos. Klon mit voller Historie (`fetch-depth: 0`), sonst bricht `history-range-guard` zu Recht. Empfehlung: jede Form, die nichts prüft, ist rot. Kennung: die CI nutzt `make doc-commits` (Anwesenheit einer Kennung, Muster nach `.d-check.yml`/Modul `commits`); die Strenge des Hooks (die Datei des Slice liegt im Index) bleibt lokal — akzeptiertes Negativ, Grund: ein Commit in der Range kann einen Slice nennen, der in einem späteren Commit der Range umbenannt wurde; die Prüfung am Index der Spitze würde alte Commits falsch ablehnen. Ob das Modul `commits` im gepinnten d-check `v0.86.1` dieselben Mengen (`ADR-`, `LH-`, `MR-`, `slice-`, `welle-`) liest wie der Hook: zu prüfen und in `SPEC-064` zu nennen (Lücke, wenn nicht). **Entschieden:** `SPEC-064` Punkte 6 bis 9 und 11 (Range; Anwesenheit der Kennung, Muster in `.d-check.yml`); die Muster von `commits` und die Blöcke `vcs` fehlen heute in `.d-check.yml` (`make doc-commits` und `make doc-immutable` prüfen darum nichts): Punkte 10 bis 12 legen sie fest, Gegenprobe Punkt 14.
12. **Gebündelte Prüfung.** Zwei Schritte im Workflow (`make doc-commits`, `make doc-immutable`) mit `if: always()` auf dem zweiten, damit eine rote erste die zweite nicht verdeckt; ein Fall hält es. Ein gemeinsames Make-Ziel entfällt. **Entschieden:** `SPEC-064` Punkt 2 (`!cancelled()` am zweiten Schritt, `Gates` nur bei Grün der beiden).
13. **Netz und Docker Hub.** Gepinnte Images werden anonym gezogen; ein Rate-Limit-Fehler ist rot und bleibt es (kein stilles Überspringen). Empfehlung: keine Anmeldung in dieser Datei (keine Geheimnisse, Randform 2); ein Fehler, den das Limit verursacht, wird als solcher in §7 festgehalten. **Entschieden:** `SPEC-064` Punkt 5 (keine Anmeldung, roter Lauf bleibt rot).

**Fragen an den Nutzer** (jede mit Empfehlung des Planners; eine Empfehlung ist keine Entscheidung):

1. **Repo öffentlich oder privat?** Es bestimmt Runner-Größe, Minutenbudget und die Wahl in Randform 10. Empfehlung: den Vertrag mit der Annahme „privat, 2 vCPU, 7 GB“ schreiben, denn sie ist die knappere.
2. **Auslöser** (Randform 1): beides wie d-check oder nur `pull_request` plus Standard-Zweig. **Entschieden:** beides, ohne `workflow_dispatch` und ohne `concurrency` (Nutzer, 2026-10-18; `SPEC-064` Punkte 1 und 2).
3. **Probe auf GitHub.** Der Beleg verlangt einen Push auf den Remote und einen Probe-Zweig mit einem Verstoß gegen die Range. Das ist eine Handlung mit Außenwirkung; die Ausführung braucht dein ausdrückliches Ja. Empfehlung: ja, danach den Zweig löschen.
4. **Branch-Schutz** (Handlung des Betreibers): wann setzt du den „required status check“? Empfehlung: nach dem ersten grünen Lauf, vor `slice-erster-release-veroeffentlichung`.

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

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (`AGENTS.md` §3.10) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (`AGENTS.md` §3.12) — `ci.yml` und `ci-range.sh` sind neue Verträge, Randformen stehen in §6 vor dem Code, je Zusage eine Mutation; `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert` und `BEO-REPO/nested-runtime-podman-oder-docker-offen` — der Runner ist eine neue Umgebung für `make gates` (Risiko 1); `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Risiko 2); `BEO-REPO/werkzeug-commit-ohne-zugelassene-kennung` und `BEO-REPO/commit-nennt-struktur-kennung` — die Range prüft Kennungen (Randform 11). Kein Eintrag erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Der Slice trägt keine Sendung eines Gebers (`SPEC-063`, `SPEC-050`, `SPEC-052` sind gelesen). Liefer-Punkte 3 (Workflow | Range | Gegenprobe samt Gate-Index). Schichten 2: Workflows (`.github/workflows/`) und Harness (`tools/harness`, `harness/mk`, `harness/vertraege`, `harness/README.md`); Planung, Vertrag und ADR sind keine Schicht; Wartungs-Doku und Handbuch bleiben unberührt. Schätzung 605 bis 890 (§4), an der Marke.

**Zählung der Sendung von `slice-harness-ci-workflows` (Architect, 2026-10-18, im Commit, der sie einträgt; `AGENTS.md` §3.13):** Die Sendung ist `.d-check.yml` (drei Blöcke, rund 25 Zeilen) und geht in die Liefer-Punkte 1 (Eintrag der Pins) und 2 (`commits`, `vcs`); die Zahl der Liefer-Punkte bleibt 3, die Schichten bleiben 2 (`.d-check.yml` ist Harness, keine dritte Schicht). Schätzung in der Metrik von §4 (Summe aus Einfügungen und Löschungen, Gegenprobe mit dem Doppelten des Produktcodes, `SPEC-064` und ADR vor dem Anspruch geschrieben und damit nicht gezählt): `ci.yml` 45 bis 70, `ci-range.sh` 80 bis 110, `.d-check.yml` 25 bis 35, `ci.mk` 20, Gegenprobe 350 bis 470 (Text von `ci.yml` mit Mutationen je Zeile, Formen der Range, Konfiguration mit echtem d-check), `harness/README.md` 15 bis 30, zusammen 535 bis 735, unter der Marke 900.
