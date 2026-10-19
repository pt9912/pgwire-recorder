# Slice slice-harness-ci: Die CI auf GitHub — `ci.yml` fährt `make gates` und prüft die Commit-Range

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-erster-release, dritte Stelle der Reihenfolge, direkt hinter `slice-harness-ci-range-konfiguration` (Schnitt vom 2026-10-18, §4) und vor `slice-erster-release-veroeffentlichung` (Auftrag des Nutzers vom 2026-10-18; Entscheidung des Nutzers vom 2026-10-18: die Nachtläufe `slice-harness-ci-gegenproben`, `-ci-freshness` und `-ci-freshness-images` folgen nach dem Release in [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md)). Start-Trigger weicht von der Welle ab: `slice-harness-ci-workflows`, `slice-harness-ci-range-konfiguration` und `slice-harness-gegenproben-aus-gates` in `done/` (§4); `slice-harness-gegenproben-aus-gates-fahrstelle-umgebung` ist keine Bedingung (Entscheidung des Nutzers vom 2026-10-18).

**Bezug:** Kein Anspruch des Produkts (`LH-*`): der Slice ändert die Prüfung der Integration, nicht das Binary. Aktive ADRs: [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Docker-only), [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Gate-Kette; der Stop-Hook bleibt), [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) und [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md) (die lokale Strenge der Kennung im Commit, die die CI nicht wiederholt, Randform 11). Die ADR dieses Slice schreibt der Architect vor dem Code (Entscheidung des Nutzers).

**Berührte Spec-Stellen:** `SPEC-064` Teile 1 und 2 und Punkt 14 für `ci.yml` und `ci-range.sh` (Vertrag des Architect in `harness/vertraege/`; Teil 3, die Konfiguration der Range-Prüfungen mit ihrer Gegenprobe, trägt `slice-harness-ci-range-konfiguration`) · `SPEC-063` (gelesen: Gate der Workflows; `ci.yml` besteht es) · [`SPEC-050`](../../../../harness/vertraege/SPEC-050-struktur-kennungen-im-commit.md) (gelesen: Kennungen im Commit) · [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) Punkte 1 bis 4 und 21 (gelesen: Nachweis des vollen Laufs, Aufrufarten).

**Verantwortlich:** pt9912 (Zweig impl-ci)
**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Jede Integration (Pull-Request, Push auf einen Zweig) läuft auf GitHub durch `make gates` und durch die Prüfung ihrer Commit-Range, unabhängig vom Klon und vom lokalen Hook: `.github/workflows/ci.yml` (nach dem Vorbild von d-check) mit `pull_request` und `push` auf Zweigen, `tags-ignore` (Tags gehören `release.yml`), `permissions: {}` oben und `contents: read` allein im Job, `timeout-minutes`, `ubuntu-latest`, Actions per SHA mit Tag-Kommentar gepinnt, einem Job mit festem Namen (der Name ist der Status, den ein Branch-Schutz verlangt, Randform 8). Die Range-Prüfung ruft `make doc-commits` (jede Commit-Message nennt eine Kennung) und `make doc-immutable` (eine angenommene ADR bleibt unverändert) mit einer `RANGE`, die ein Skript `tools/harness/ci-range.sh` aus dem Ereignis berechnet (nicht im YAML, wo sie ungeprüft bliebe); eine leere oder unauflösbare Range ist rot, nicht grün (`history-range-guard`). Gegenproben prüfen die Zusagen der Datei und die Berechnung der Range. **Übernommen von `slice-harness-ci-workflows` (Sendung des Architect, 2026-10-18):** der Eintrag für die Pins in `.d-check.yml` (`workflows` in `modules:` und `workflows.dir`, `SPEC-063` Punkt 8): das Modul meldet ohne Datei und ohne `uses:` einen Befund, der Eintrag kommt darum mit `ci.yml`. **Abgegeben an `slice-harness-ci-range-konfiguration` (Schnitt vom 2026-10-18, §4):** die Blöcke `commits` und `vcs` in `.d-check.yml` (`SPEC-064` Punkte 10 bis 12) mit ihrer Gegenprobe; ohne sie prüfen `make doc-commits` und `make doc-immutable` in diesem Repo nichts. Dieser Slice setzt sie voraus.

**Was im Repo vorliegt und was fehlt** (Bestand vom Planner am 2026-10-18 gelesen): Es liegen `make doc-commits` und `make doc-immutable` (`d-check.mk`, mit `RANGE`, über das Fragment `harness/mk/doc-gate.mk` mit dem Vorlauf `history-range-guard`), `tools/harness/commit-msg-traceability.sh` (eine Message-Datei, Mengen `ADR-`, `LH-`, `MR-`, `slice-`) und der lokale Träger `.githooks/commit-msg` (prüft zusätzlich, dass die Slice- oder Welle-Datei im Index liegt, [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md)). **Es fehlen:** jeder Workflow (`.github/` gibt es nicht), die Berechnung der Range aus dem Ereignis, ein Ziel, das beide Range-Prüfungen bündelt (nicht nötig, zwei Schritte genügen), und ein Sensor, dass `ci.yml` hält, was sie zusagt. Ein ADR-Immutable-Skript gibt es nicht; es ist das Modul `vcs` von d-check, das `make doc-immutable` fährt (Randform 11).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Blöcke `commits` und `vcs` in `.d-check.yml` und ihre Gegenprobe (`SPEC-064` Punkte 10 bis 12 und der Teil von Punkt 14, der sie beschreibt) — `slice-harness-ci-range-konfiguration` (Geber, vor diesem Slice; Schnitt vom 2026-10-18 aus Größe, §4): Die Konfiguration hängt weder an `ci.yml` noch an `ci-range.sh` und liefert allein Wert; `ci.yml` ruft `make doc-commits` und `make doc-immutable` und setzt voraus, dass sie etwas prüfen.
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
- [ ] **Liefer-Punkt 2 (Range):** `tools/harness/ci-range.sh` berechnet aus den Eingaben des Ereignisses (Ereignisart, Basis des Pull-Requests, Vorgänger-Stand des Pushes, Stand, Standard-Zweig) die `RANGE` für `make doc-commits` und `make doc-immutable`, je Randform aus §6 (Pull-Request, Push, neuer Zweig, Vorgänger nicht mehr im Klon, leere Range, Tiefe des Klons) mit der im Vertrag entschiedenen Antwort; ein Fehler endet mit Exit ungleich 0 und einer Meldung, nie mit einer Range, die nichts prüft. Die Blöcke `commits` und `vcs` liegen mit `slice-harness-ci-range-konfiguration` vor (Start-Trigger); dieser Slice ändert sie nicht.
- [ ] **Liefer-Punkt 3 (Gegenprobe):** `make ci-gegenprobe` (Name: Architect; ohne die Fälle der Konfiguration, die `slice-harness-ci-range-konfiguration` in einem eigenen Ziel führt) zeigt je Zusage von `SPEC-064` Teile 1 und 2 und je Element ihrer Aufzählungen einen Fall und eine Mutation, die nur ihn färbt (`AGENTS.md` §3.10): an Kopien von `ci.yml` (Auslöser, Rechte, Zeitgrenze, `fetch-depth`, `persist-credentials`, Jobname, kein Geheimnis, Schritte in Reihenfolge) und an `ci-range.sh` (jede Form der Range); die unveränderte Kopie grün; Aufrufarten von `make` haben keinen Fall (`SPEC-064` Punkt 14: kein Hook und kein Stand liest das Ergebnis, `ci-range.sh` läuft nie unter `make`). Das Ziel steht in `GEGENPROBEN_CONTAINER`. `harness/README.md` §Sensors führt die Zeile; sie nennt Branch-Schutz als Handlung des Betreibers und die CI als Ergänzung des Stop-Hooks, nicht als Ersatz.
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
| `.github/workflows/ci.yml` | neu | Text wörtlich nach `SPEC-064` Punkt 1 (Auslöser, Rechte, Zeitgrenze, Range-Prüfung, `make gates`); `actions/checkout` auf den Commit des Tags v7.0.1 gepinnt |
| `tools/harness/ci-range.sh` | neu | `RANGE` aus dem Ereignis; fail-closed (Exit 0, 1, 2 nach `SPEC-064` Punkte 6 bis 9) |
| `harness/mk/ci.mk`, `tools/harness/ci-gegenprobe.sh` | neu | Gegenprobe der Zusagen von `ci.yml` und der Range: 148 Fälle (`p1` 99, `p6` 26, `p7` 10, `p8` 4, `p9` 6, `p18` 3), 144 Mutationen (99 am Text von `ci.yml`, 38 an `ci-range.sh`, 7 am Skript); Eintrag in `GEGENPROBEN_CONTAINER` |
| `.d-check.yml` | update | Eintrag der Pins (`workflows` in `modules:` und `workflows.dir`, `SPEC-063` Punkt 8); Sendung von `slice-harness-ci-workflows`. Die Blöcke `commits` und `vcs` sind bei `slice-harness-ci-range-konfiguration` |
| `harness/vertraege/SPEC-064-…md` | neu (Architect) | Vertrag; vor dem Code |
| `harness/README.md` §Sensors | update | Zeile von `make ci-gegenprobe`; Grenze Branch-Schutz (Job `ci-gates`); Verhältnis zum Stop-Hook |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-ci-workflows`, `slice-harness-ci-range-konfiguration` und `slice-harness-gegenproben-aus-gates` liegen in `done/` (der erste für das Gate der Workflows, der zweite für die Blöcke `commits` und `vcs` in `.d-check.yml`, der dritte, weil `make gates` ohne ihn die Gegenproben fährt und die Zeitgrenze des Jobs nicht trägt), dazu der Start-Trigger von [welle-erster-release](../welle-erster-release.md), und der Architect hat `SPEC-064` und die ADR geschrieben, jede Randform aus §6 entschieden (`AGENTS.md` §3.12). Schritt 2 der Reihenfolge in §5 der Welle.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: mehr als drei Liefer-Punkte, oder der Diff ist nicht in einer Review-Sitzung prüfbar — zurück zur Zerlegung. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- .github tools harness .d-check.yml Makefile`, Tests und Gegenprobe eingeschlossen, liegt über 900; Basis: `39f1ce8a` (Elternteil des Commits, der diese Datei nach `in-progress/` bewegt). Schätzung nach dem Schnitt vom 2026-10-18 (Planner, ohne Basis): Die Schätzung des Architect für den ungeteilten Slice (535 bis 735, §8) lag bei `slice-harness-ci-workflows` um den Faktor 1,3 bis 1,65 unter der Messung (899 gegen 545 bis 695; die Gegenprobe trug zwei Drittel). Mit dem Faktor liegt der ungeteilte Slice bei 695 bis 1210, Mitte 920, über der Marke; die Konfiguration der Range-Prüfungen (Blöcke `commits` und `vcs` mit ihrer Gegenprobe, 27 Prozent der Schätzung, 145 bis 200) ging deshalb nach `slice-harness-ci-range-konfiguration`. Dieser Slice behält 390 bis 540 (Schätzung des Architect abzüglich 27 Prozent), mit dem Faktor 510 bis 890, Mitte 670, unter der Marke. Die Spanne reicht mit ihrer Obergrenze an die Marke; Schnitt, falls sie eintritt: *Range* (Liefer-Punkt 2 mit seinem Teil der Gegenprobe) und *Workflow* (Punkt 1 und 3).
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
- **Der Pin-Check läuft zum ersten Mal an einer echten Datei und hängt am gepinnten d-check** (Version des Moduls `workflows`, Verhalten bei `.yaml`, Grenze der Zerlegung der Blockform; Risiko 2 von `slice-harness-ci-workflows`, dessen Gate nur die Bindung in `.d-check.yml` prüft, nicht die Pins). — **Ausgang (vorgesehen):** entfallen, wenn `make docs-check` mit dem Eintrag in `.d-check.yml` an `ci.yml` 0 Befunde meldet und ein Pin ohne SHA oder ohne Tag-Kommentar in der Probe (`SPEC-064` Punkt 1, Fälle `p1-pin-*`) rot wird; sonst eingetreten → Rückgabe an den Architect. Die offene Beobachtung führt `BEO-REPO/workflow-pin-pruefung-am-gepinnten-d-check`.
- **Doppelte Läufe** (derselbe Commit löst `push` und `pull_request` aus). — **Ausgang (vorgesehen):** entfallen, wenn der Vertrag Auslöser und `concurrency` entscheidet (Randform 1 und 3) und die Gegenprobe es zeigt.
- **Die Größe liegt nach dem Schnitt unter der Marke, mit der Obergrenze an ihr** (510 bis 890 mit dem gemessenen Faktor, Marke 900; ungeteilt 695 bis 1210). — **Ausgang (vorgesehen):** entfallen, wenn die Messung nach der Metrik unter 900 bleibt; darüber meldet der Implementer (§4 Rückführung *zu groß*).

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
11. **Range und Kennung.** Berechnung: Pull-Request `Basis..Kopf` (der Kopf ist der letzte Commit des Zweigs, nicht der Merge-Commit von GitHub, der mit `Merge` beginnt und den Betreff-Test des Skripts übergeht); Push `Vorgänger..Stand`; neuer Zweig (Vorgänger nur Nullen) gegen den Standard-Zweig; Vorgänger nicht mehr im Klon (Force-Push); leere Range; erster Commit des Repos. Klon mit voller Historie (`fetch-depth: 0`), sonst bricht `history-range-guard` zu Recht. Empfehlung: jede Form, die nichts prüft, ist rot. Kennung: die CI nutzt `make doc-commits` (Anwesenheit einer Kennung, Muster nach `.d-check.yml`/Modul `commits`); die Strenge des Hooks (die Datei des Slice liegt im Index) bleibt lokal — akzeptiertes Negativ, Grund: ein Commit in der Range kann einen Slice nennen, der in einem späteren Commit der Range umbenannt wurde; die Prüfung am Index der Spitze würde alte Commits falsch ablehnen. Ob das Modul `commits` im gepinnten d-check `v0.86.1` dieselben Mengen (`ADR-`, `LH-`, `MR-`, `slice-`, `welle-`) liest wie der Hook: zu prüfen und in `SPEC-064` zu nennen (Lücke, wenn nicht). **Entschieden:** `SPEC-064` Punkte 6 bis 9 und 11 (Range; Anwesenheit der Kennung, Muster in `.d-check.yml`); die Muster von `commits` und die Blöcke `vcs` fehlen heute in `.d-check.yml` (`make doc-commits` und `make doc-immutable` prüfen darum nichts): Punkte 10 bis 12 legen sie fest, Gegenprobe Punkt 14; sie setzt `slice-harness-ci-range-konfiguration` um (Schnitt vom 2026-10-18), dieser Slice setzt sie voraus.
12. **Gebündelte Prüfung.** Zwei Schritte im Workflow (`make doc-commits`, `make doc-immutable`) mit `if: always()` auf dem zweiten, damit eine rote erste die zweite nicht verdeckt; ein Fall hält es. Ein gemeinsames Make-Ziel entfällt. **Entschieden:** `SPEC-064` Punkt 2 (`!cancelled()` am zweiten Schritt, `Gates` nur bei Grün der beiden).
13. **Netz und Docker Hub.** Gepinnte Images werden anonym gezogen; ein Rate-Limit-Fehler ist rot und bleibt es (kein stilles Überspringen). Empfehlung: keine Anmeldung in dieser Datei (keine Geheimnisse, Randform 2); ein Fehler, den das Limit verursacht, wird als solcher in §7 festgehalten. **Entschieden:** `SPEC-064` Punkt 5 (keine Anmeldung, roter Lauf bleibt rot).

14. **Reihenfolge der Gründe von `ci-range.sh`.** **Entschieden:** `SPEC-064` Punkt 9 (Ereignis, Variablen, flacher Klon, Auflösung, leere Range).
15. **Zeilenende und Kodierung** (Workflow-Datei, Messages). **Entschieden:** `SPEC-064` Punkt 17 und `SPEC-063` Punkt 11 (Grenze, kein Code).
16. **Plattform und Aufrufarten von `make`.** **Entschieden:** `SPEC-064` Kopf, Plattform (`linux/amd64`, `ubuntu-24.04`; andere nicht zugesagt) und Punkt 14 (kein Fall: kein Hook und kein Stand liest das Ergebnis).
17. **Teilung der Gegenprobe.** **Entschieden:** `SPEC-064` Kopf („Zwei Slices, zwei Ziele“): dieser Slice liefert `ci-gegenprobe` (Fragment `harness/mk/ci.mk`, Skript `tools/harness/ci-gegenprobe.sh`) mit den Fällen `p1-*` bis `p9-*`, ohne Docker und ohne `make`; der Eintrag `workflows` in `.d-check.yml` (Punkt 13) gehört diesem Slice.

18. **Abbruch und Herunterfahren von `ci-gegenprobe`.** **Entschieden:** `SPEC-064` Punkt 18 (Aufräumen bei Ende, `SIGINT`, `SIGTERM`; Fälle `p18-sigint`, `p18-sigterm`; die Mutation „Unterprozesse nicht beendet“ gilt auch hier).

19. **Meldungstexte, Reihenfolge innerhalb einer Stufe, Zusatzzeile aus Leerraum, Tag-Kommentar, Hinweis bei Fehler** (Rückgabe des Implementers). **Entschieden:** `SPEC-064` Punkt 8 (Wortlaut nicht zugesagt außer Variablenname und Hinweis auf `fetch-depth: 0`), Punkt 9 (Reihenfolge innerhalb einer Stufe nicht zugesagt), Punkt 1 (leere Zeilen und Kommentarzeilen frei; Leerraum-Zeile nicht zugesagt; Tag beginnt mit einem Zeichen, das kein Leerzeichen ist, Fall `p1-pin-tag-leer`).

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

### Belege des Implementers

Stand: Code-Commits `4190a5c3` und `a31f4a7f` (die Belege gelten für diesen Stand, wo sie nichts anderes nennen); Zweig `impl-ci`, Basis der Metrik `39f1ce8a`. Die Läufe stehen im Dev-Container pgr-dev3 (Benutzer `vscode`, Klon `/workspaces/pgr-dev3/impl-ci-klon`, Docker-in-Docker, d-check v0.86.1, Digest `3e0b9779…`, git 2.55.0, bash 5.2.37). **Nicht belegt in diesem Abschnitt:** ein Lauf von `ci.yml` auf GitHub (Nummer, Dauer, ein roter Lauf mit Verstoß gegen die Range): `ci.yml` ging nicht nach GitHub; die Probe ist eine Handlung des Koordinators mit Zustimmung des Nutzers (§5), DoD Liefer-Punkt 1 ist damit zur Hälfte belegt (Datei, Gates), zur Hälfte offen (Lauf).

**Pin (Risiko 4).** `actions/checkout`: `git ls-remote https://github.com/actions/checkout` aus dem Container nennt `refs/tags/v7.0.1` mit `3d3c42e5aac5ba805825da76410c181273ba90b1` (kein `^{}`-Eintrag, der Tag zeigt unmittelbar auf einen Commit; `v7` zeigt auf denselben); `git fetch --depth 1 … 3d3c42e5…` ergibt Typ `commit`, Betreff `prep v7.0.1 release (#2531)`, Datum 2026-07-17; `v6.1.0` ist die ältere Linie. Die Zeile in `ci.yml` ist `uses: actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`. Die Ermittlung lief im Container mit Netz zu github.com; eine zweite Ermittlung am Host ist eine Handlung des Koordinators.

**Gates am Arbeitsbaum.** `make workflows-check`: `1 Datei(en) geprüft, kein Befund`; `make docs-check`: `1196 Datei(en) geprüft, 0 Befund(e)` (mit `workflows` in `modules:` und `workflows.dir`). Rot gesehen (Risiko 4): `checkout@v7` ⇒ `uses-pin-missing` (`kein voller 40-stelliger Commit-SHA`), der SHA ohne ` # v7.0.1` ⇒ `uses-pin-untagged`, jeweils `1 Befund(e)`, Exit 2 von `make`; `timeout-minutes: 0` ⇒ `workflows-check: 2 Befund(e)`.

**Die Schritte von `ci.yml` lokal gefahren** (nicht ein Lauf auf GitHub): `ci-range.sh` mit `CI_EVENT=push`, `CI_PUSH_BEFORE` aus Nullen, `CI_DEFAULT_BRANCH=main` ergibt `5e51b0c8…..4190a5c3…` (Exit 0), `pull_request` mit `CI_PR_BASE=origin/main` dieselbe Range; `make doc-commits RANGE=39f1ce8a…..4190a5c3…` und `make doc-immutable RANGE=…`: je `3 Commit(s) — OK` und `d-check: 1196 Datei(en) geprüft, 0 Befund(e)`.

**Läufe nach dem Review** (Code-Commit `7d80d9f2`): `make ci-gegenprobe` `grün (147 Läufe)` (69 s; 144 Mutationen, die unveränderte Kopie, `rekursion-schutz`, `fall-ohne-funktion`), `make gates-schnell-gegenprobe` grün (63 s), `make docs-check` 0 Befunde. `make gates` lief nicht erneut: Seit `a31f4a7f` änderten sich außer Gegenprobe, Plan und `SPEC-064` nur Kommentarzeilen in `ci.mk` und eine Zeile in `harness/README.md`; die Kette liest sie nicht anders als `docs-check` (grün).

**Läufe nach der Entscheidung des Architect** (Code-Commit `875f4276`): `make ci-gegenprobe` `grün (144 Läufe)` (58 s), `make gates-schnell-gegenprobe` grün (61 s), `make docs-check` 0 Befunde. `make gates` lief hier nicht erneut: seit dem Lauf unten ändern nur die Gegenprobe (`ci-gegenprobe.sh`, in keiner Kette), `SPEC-064` und der Plan; die Gates dieser Kette lesen davon nur Vertrag und Plan (`docs-check`, `kopf-check`, grün). Der Nachweis von `make gates` gilt für den Stand `a31f4a7f`.

**Läufe am Endstand** (Code-Commit `a31f4a7f`, nacheinander, ohne parallele volle Läufe): `make ci-gegenprobe` `grün (144 Läufe)` (58 s); `make gates-schnell-gegenprobe` grün (59 s); `make gates` Exit 0 (261 s = 4 min 21 s, Nachweis geschrieben); `make docs-check` 0 Befunde. **Dauer von `make gates` gegen die Zeitgrenze (Risiko 2):** 261 s im Dev-Container gegen 1800 s (`timeout-minutes: 30`), 15 Prozent; auf dem Runner (2 vCPU, kalte Docker-Layer, Pull der gepinnten Images) ist sie nicht gemessen, das zeigt nur der Lauf auf GitHub.

**Fälle (148).** `p1-*` (99): `gleich` (der Arbeitsbaum), `text-gleich` (die Kopie aus dem Text), je Zeile n des Textes (n = 1 bis 39) `fehlt-n` und `geaendert-n` (erstes Zeichen der Zeile durch `Z` ersetzt), `zusatz-<Wort>` (`concurrency`, `workflow-dispatch`, `pull-request-target`, `workflow-run`, `strategy`, `continue-on-error`, `secrets`, `schritt-cache`, `schritt-artefakt`, jeweils eine Zeile am Ende der Datei), `pin-kurz` (39 Ziffern), `pin-gross`, `pin-ohne-tag`, `pin-tag-leer` (`# ` ohne Text), `pin-andere-action`, `reihenfolge-schritte` (ADR-Schritt vor Commit-Schritt), `reihenfolge-gates` (Gates vor beiden), `leerzeile`, `kommentarzeile`, `kommentarzeile-eingerueckt` (die letzten drei grün erwartet). `p6-*` (26): `default-vorgaenger-fehlt-fehlt|ungueltig` (Push mit unbekanntem Vorgänger und fehlendem bzw. ungültigem `CI_DEFAULT_BRANCH`), `event-fehlt|leer|unbekannt|schedule`, je Variable `CI_HEAD`, `CI_PR_BASE`, `CI_PUSH_BEFORE` fünf Formen (`fehlt`, `ungueltig` mit `g`, `kurz` mit 39, `lang` mit 41, `gross` mit `A`), `default-fehlt|ungueltig`, `ungenutzt`, `ungenutzt-prbase`, `ungenutzt-default`. `p7-*` (10): `pr`, `pr-merge-kopf`, `pr-basis-fehlt-im-klon`, `pr-kopf-fehlt-im-klon`, `push`, `push-neu`, `push-neu-ohne-standard`, `push-vorgaenger-fehlt`, `leer`, `flach`. `p8-*` (4): `rev-list-fehler` (Basis und Kopf auflösbar, ein Commit dazwischen fehlt: das Objekt von `c2` ist gelöscht), `stdout-eine-zeile`, `fehler-stdout-leer` (acht Formen mit Exit ungleich 0: Stdout leer, jede Zeile von stderr mit `ci-range:`), `git-fehler`. `p9-*` (6): `ereignis-vor-head`, `head-vor-prbase`, `head-vor-before`, `eingabe-vor-flach`, `flach-vor-aufloesung`, `aufloesung-vor-leer`. `p18-*` (3): `sigint`, `sigterm`, `ende` (eine Kopie mit einem einzigen schnellen Lauf endet mit Exit 0, ihr `TMPDIR` ist leer). Ein Fall von `ci-range.sh` besteht, wenn Exit, stdout und stderr vollständig dem erwarteten Text entsprechen (jede Zeile mit Zeilenende); wo der Vertrag den Wortlaut nicht zusagt (Punkt 8), nur, wenn Exit-Code, leere stdout, das Präfix `ci-range:` in jeder Zeile von stderr und der Ort (Variable bzw. `fetch-depth: 0`) stimmen (Basis/Kopf nicht im Klon, Standard-Zweig fehlt, git-Fehler: nur Exit und Präfix); bei `pruefe_ci`, wenn `pruefe_ci` das erwartete Urteil fällt. Die Köder: jeder Fall von `ci-range.sh` setzt alle fünf Variablen gültig (die vom Ereignis nicht gelesenen als Köder mit anderer Range), dann die Abweichung des Falls. Der Name eines Falls ohne Funktion zählt als roter Fall (Lauf `fall-ohne-funktion`).

**Mutationen (144): Zusage · Mutation · roter Test.** Jede Mutation ist ein Textersatz in `pruefe_ci` (Quelltext aus `declare -f`, ausgewertet im Lauf) bzw. in einer Kopie von `ci-range.sh` bzw. im Skript (`m-p18-*`); rot sind genau die benannten Fälle (die Gegenprobe verlangt diese Menge); mit `GP_ENTDECKE=1` gesehen, dann der volle Lauf grün (144 Läufe: 141 Mutationen, die unveränderte Kopie, `rekursion-schutz`, `fall-ohne-funktion`). Fälle, die eine Mutation außerhalb der Aufzählung färbt, sind gemessen und je Zeile begründet.

Text von `ci.yml` (99):

| Mutation | Zusage / Element | rote Fälle |
|---|---|---|
| `m-p1-zeile-<n>-optional` (n = 1 bis 39, 39 Mutationen) | Zeile n ist Pflicht (fehlt sie, rückt die Prüfung nicht weiter) | `p1-fehlt-<n>` |
| `m-p1-zeile-<n>-nicht-verglichen` (39) | der Inhalt von Zeile n ist Teil des Textes | `p1-geaendert-<n>`; für n = 25 (die Zeile mit dem Pin) dazu `p1-pin-kurz`, `-gross`, `-ohne-tag`, `-tag-leer`, `-andere-action` (alle ändern diese Zeile) |
| `m-p1-zusatz-<Wort>-erlaubt` (9) | die Zeile mit `<Wort>` ist verboten | `p1-zusatz-<Wort>` |
| `m-p1-zusatzzeilen-erlaubt` | keine Zeile über den Text hinaus | alle neun `p1-zusatz-*` |
| `m-p1-pin-grossbuchstaben` | SHA in Kleinschrift | `p1-pin-gross` |
| `m-p1-pin-39-ziffern` | SHA mit 40 Ziffern | `p1-pin-kurz` |
| `m-p1-pin-ohne-kommentar` | Tag-Kommentar hinter dem SHA | `p1-pin-ohne-tag` |
| `m-p1-pin-kommentar-ohne-text` | der Tag-Kommentar hat Text (`# ` allein ist keiner, Punkt 1) | `p1-pin-tag-leer` |
| `m-p1-pin-andere-action` | nur `actions/checkout` | `p1-pin-andere-action` |
| `m-p1-kommentarzeilen-nicht-ignoriert` | Kommentarzeilen zählen nicht | `p1-kommentarzeile`, `p1-kommentarzeile-eingerueckt` |
| `m-p1-leerzeilen-nicht-ignoriert` | Leerzeilen zählen nicht | `p1-leerzeile` |
| `m-p1-kommentar-nur-am-zeilenanfang` | erstes Zeichen außer Leerzeichen `#` (eingerückt) | `p1-kommentarzeile-eingerueckt` |
| `m-p1-kommentar-irgendwo-in-der-zeile` | `#` in der Zeile macht sie nicht zur Kommentarzeile (die Zeile des Pins trägt `# v7.0.1`) | `p1-gleich`, `p1-text-gleich`, `p1-leerzeile`, `p1-kommentarzeile`, `p1-kommentarzeile-eingerueckt` (jede grün erwartete Datei hat die Pin-Zeile) |
| `m-p1-reihenfolge-commit-adr-tauschbar` | `Commit-Messages` vor `ADR-Immutabilität` | `p1-reihenfolge-schritte` |
| `m-p1-reihenfolge-gates-zuerst` | `Gates` nach den beiden Range-Schritten | `p1-reihenfolge-gates` |

`ci-range.sh` (38):

| Mutation | Zusage / Element | rote Fälle |
|---|---|---|
| `m-p6-event-ungeprueft` | `CI_EVENT` ist `pull_request` oder `push` | `p6-event-fehlt`, `-leer`, `-unbekannt`, `-schedule`, `p8-fehler-stdout-leer` (Form `event`), `p9-ereignis-vor-head` |
| `m-p6-event-leer-erlaubt` | leeres/fehlendes `CI_EVENT` ist ungültig | `p6-event-fehlt`, `p6-event-leer`, `p8-fehler-stdout-leer`, `p9-ereignis-vor-head` (der Code unterscheidet fehlend und leer nicht) |
| `m-p6-event-dispatch-erlaubt` / `m-p6-event-schedule-erlaubt` | `workflow_dispatch` / `schedule` sind ungültig | `p6-event-unbekannt` / `p6-event-schedule` |
| `m-p6-head-ungeprueft` | `CI_HEAD` ist 40 Hexziffern | die fünf `p6-head-*`, `p9-head-vor-prbase`, `p9-head-vor-before`, `p9-eingabe-vor-flach` (diese drei haben `CI_HEAD` fehlend) |
| `m-p6-prbase-ungeprueft` / `m-p6-before-ungeprueft` | `CI_PR_BASE` / `CI_PUSH_BEFORE` ist 40 Hexziffern | die fünf `p6-prbase-*` / die fünf `p6-before-*` |
| `m-p6-default-ungeprueft` | `CI_DEFAULT_BRANCH` ist ein Zweigname | `p6-default-fehlt`, `p6-default-ungueltig`, `p6-default-vorgaenger-fehlt-fehlt`, `p6-default-vorgaenger-fehlt-ungueltig` |
| `m-p6-default-leer-erlaubt` | leerer Zweigname ist ungültig | `p6-default-fehlt`, `p6-default-vorgaenger-fehlt-fehlt` |
| `m-p6-default-format-ungeprueft` | Zweigname nach `git check-ref-format --branch` | `p6-default-ungueltig`, `p6-default-vorgaenger-fehlt-ungueltig` |
| `m-p6-default-vorgaenger-ungeprueft` | `CI_DEFAULT_BRANCH` wird auch bei unbekanntem Vorgänger geprüft (Punkt 6) | `p6-default-vorgaenger-fehlt-fehlt`, `p6-default-vorgaenger-fehlt-ungueltig`, `p7-push-vorgaenger-fehlt`, `p8-stdout-eine-zeile` (die Mutante liest `$standard` ungesetzt, `set -u` ⇒ Exit 1, auch bei gültigem Wert) |
| `m-p6-pr-liest-before` / `m-p6-push-liest-prbase` / `m-p6-push-liest-standard` | eine vom Ereignis nicht gelesene Variable wird nicht geprüft | `p6-ungenutzt` / `p6-ungenutzt-prbase` / `p6-ungenutzt-default` |
| `m-p6-form-39-ziffern` / `-41-ziffern` / `-grossbuchstaben` / `-nicht-hex` | 40 Ziffern genau, Kleinschrift, hexadezimal | `p6-{head,prbase,before}-kurz` / `-lang` / `-gross` / `-ungueltig` (je drei) |
| `m-p7-pr-nimmt-vorgaenger` | Pull-Request: Basis ist `CI_PR_BASE` | `p7-pr`, `p7-pr-merge-kopf`, `p7-pr-basis-fehlt-im-klon`, `p6-ungenutzt`, die fünf `p6-prbase-*`, `p8-fehler-stdout-leer` (Form `basis`) (`p9-aufloesung-vor-leer` bleibt grün: die Mutante endet mit Exit 2 und Präfix, der Wortlaut ist nicht zugesagt) |
| `m-p7-push-nimmt-pr-basis` | Push: Basis ist `CI_PUSH_BEFORE` | `p7-push`, `p7-push-vorgaenger-fehlt`, `p7-push-neu-ohne-standard`, `p6-ungenutzt-prbase`, `p6-ungenutzt-default`, die fünf `p6-before-*`, `p6-default-fehlt`, `p6-default-ungueltig`, `p8-fehler-stdout-leer`, `p8-stdout-eine-zeile`, `p6-default-vorgaenger-fehlt-fehlt`, `p6-default-vorgaenger-fehlt-ungueltig` (Köder `CI_PR_BASE` ist gleich dem Stand des Standard-Zweigs, darum bleiben `p7-push-neu` und `p7-leer` grün) |
| `m-p7-nullen-kein-neuer-zweig` | Nullen als Vorgänger: neuer Zweig, Basis ist der Standard-Zweig | `p7-push-neu`, `p7-leer` (die Mutante nimmt den Weg „Vorgänger fehlt“ und schreibt dessen Hinweis zusätzlich; `p7-push-neu-ohne-standard`, `p6-default-fehlt`, `p6-default-ungueltig` bleiben grün: Exit 2 mit Präfix, der Wortlaut ist nicht zugesagt) |
| `m-p7-standard-fehlt-ohne-exit` | fehlt der Standard-Zweig, Exit 2 | `p7-push-neu-ohne-standard` |
| `m-p7-vorgaenger-fehlt-exit-2` | Vorgänger nicht im Klon: Basis ist der Standard-Zweig | `p7-push-vorgaenger-fehlt`, `p8-stdout-eine-zeile`, `p6-default-vorgaenger-fehlt-fehlt`, `p6-default-vorgaenger-fehlt-ungueltig` (nennen `CI_DEFAULT_BRANCH` nicht mehr) |
| `m-p7-basis-fehlt-ohne-exit` | Basis eines Pull-Requests nicht im Klon: Exit 2 | `p7-pr-basis-fehlt-im-klon`, `p8-fehler-stdout-leer` (`p9-aufloesung-vor-leer` bleibt grün, Grund wie oben) |
| `m-p7-kopf-fehlt-ohne-exit` | Kopf nicht im Klon: Exit 2 | `p7-pr-kopf-fehlt-im-klon`, `p8-fehler-stdout-leer` |
| `m-p7-leer-exit-0` | leere Range: Exit 1 | `p7-leer`, `p8-fehler-stdout-leer` |
| `m-p7-flach-nicht-erkannt` | flacher Klon: Exit 2 | `p7-flach`, `p9-flach-vor-aufloesung` (`p8-fehler-stdout-leer` bleibt grün: die Mutante endet im Fall `flach` mit Exit 1, leere Range, ohne Ausgabe auf stdout) |
| `m-p7-merge-kopf-abgelehnt` | der Kopf darf ein Merge-Commit sein | `p7-pr-merge-kopf`, `p8-rev-list-fehler` (`git` schreibt dort ohne Präfix auf stderr) |
| `m-p7-namen-statt-shas` | Ausgabe mit aufgelösten SHAs | `p7-push-neu`, `p7-push-vorgaenger-fehlt`, `p7-leer`, `p8-stdout-eine-zeile` |
| `m-p8-hinweis-auf-stdout` | Hinweis auf stderr | `p7-push-vorgaenger-fehlt`, `p8-stdout-eine-zeile`, `p6-default-vorgaenger-fehlt-fehlt`, `p6-default-vorgaenger-fehlt-ungueltig` (stdout nicht leer) |
| `m-p8-teil-range-auf-stdout` | Fehler: stdout leer | `p7-pr-kopf-fehlt-im-klon`, `p8-fehler-stdout-leer` |
| `m-p8-rev-list-fehler-ignoriert` | ein Fehler von `git rev-list --count` ist Exit 2 (Punkt 8, zweite Stelle eines `git`-Fehlers) | `p8-rev-list-fehler` (die Mutante gibt bei fehlendem Zwischen-Commit die Range mit Exit 0 aus) |
| `m-p9-ereignis-nach-head` / `m-p9-head-nach-prbase` / `m-p9-head-nach-before` | Reihenfolge Ereignis, `CI_HEAD`, `CI_PR_BASE` bzw. `CI_PUSH_BEFORE` | `p9-ereignis-vor-head` / `p9-head-vor-prbase` / `p9-head-vor-before` |
| `m-p9-flach-vor-eingabe` | Variablen vor flachem Klon | `p9-eingabe-vor-flach`, `p8-fehler-stdout-leer`, `p8-git-fehler` (die Mutante ruft `git` ohne Meldung mit Präfix) |
| `m-p9-aufloesung-vor-flach` | flacher Klon vor Auflösung | `p9-flach-vor-aufloesung` (`p8-git-fehler` bleibt grün: Exit 2 mit Präfix) |
| `m-p9-leer-vor-aufloesung` | Auflösung vor leerer Range | `p9-aufloesung-vor-leer` |

Abbruch (7): `m-p18-int-falle-fehlt`, `m-p18-int-kinder-bleiben`, `m-p18-int-verzeichnis-zu-frueh` ⇒ `p18-sigint`; `m-p18-term-*` ebenso ⇒ `p18-sigterm`; `m-p18-ende-falle-fehlt` (ohne `trap raeume_auf EXIT` bleibt `tmp.*` nach dem regulären Ende liegen) ⇒ `p18-ende`. Die Mutation „Nachkommen nicht beendet“ steht trotz des Satzes in `SPEC-064` Punkt 18 in der Gegenprobe, weil `ci-gegenprobe.sh` Nachkommen hat (Hintergrund-Unterprozesse); der Satz „entfällt“ in Punkt 18 und Randform 18 widerspricht dem Skript (Rückgabe an den Architect, F-4 des Reviews).

Kein Mutant blieb grün. `p8-git-fehler` ist ein Beleg ohne eigene Mutation: Eine Mutante, die nur den Fehler von `git rev-parse --is-shallow-repository` nicht meldet (`|| fehler …` gestrichen), ist für die Zusage äquivalent. Der Geltungsbereich ist dieser eine Aufruf; für `git rev-list --count` gilt die Invariante nicht (Review F-1) und der Aufruf hat Fall und Mutation (`p8-rev-list-fehler`, `m-p8-rev-list-fehler-ignoriert`). Invariante: außerhalb eines Klons endet jeder weitere `git`-Aufruf des Skripts (Auflösung von Basis und Kopf) mit Fehler und Exit 2 mit Präfix; Ort: die Zweige `case "$flach"` und die Auflösung von Basis und Kopf; Sonde: die Kopie ohne `|| fehler …` in einem Verzeichnis ohne Klon ergibt `ci-range: git rev-parse --is-shallow-repository ergab: ` und Exit 2, das Original `… schlug fehl` und Exit 2 (nur der Wortlaut trennt sie, und der ist nicht zugesagt). Sonst: Äquivalente Mutanten keine behauptet.

**Weitere Sonden (rot gesehen, nicht im Lauf).** (1) Die Zeile `declare -F …` der Fallliste in einer Kopie gestrichen, Lauf `fall-ohne-funktion`: `ROT — Fall ohne Funktion: erwartet rot: p1-tippfehler; rot sind: nichts`. (2) Die Zeile `[ -z "${GP_INNEN-}" ] || return 0` in `mut18` einer Kopie ersetzt: der Lauf `rekursion-schutz` meldet `Rekursionsschutz …`, mit der Zeile grün. (3) `timeout-minutes: 31` im Arbeitsbaum: `p1-gleich` rot (`ci.yml des Arbeitsbaums weicht von Punkt 1 ab`), die unveränderte Kopie ist nicht grün. (4) `exit 1` der leeren Range im Arbeitsbaum zu `exit 0`: `unveränderte Kopie ist nicht grün: p7-leer p8-fehler-stdout-leer`. (5) Vor `a31f4a7f` (damals 143 Läufe): als root (`sudo`), volle Gegenprobe `grün (143 Läufe)`; als Benutzer mit `umask 077`, `LC_ALL=en_US.utf8`, `TZ=Asia/Tokyo` und `GP_PARALLEL=8` ebenfalls `grün (143 Läufe)`; am Endstand mit dem Vorgabewert 4 `grün (144 Läufe)`. (6) Drei Läufe der sechs Mutationen des Abbruchs hintereinander mit derselben roten Menge. Kein Fall hängt von Dateirechten ab.

**Messung Zeilenende und Leerzeichen (SPEC-064 Punkt 17; keine Zusage, kein Fall).** `pruefe_ci` an Kopien des Arbeitsbaums mit der Funktion aus dem Skript: jede Zeile mit `\r\n` ⇒ rot; eine Zeile nur aus einem Tabulator zwischen Zeile 3 und 4 ⇒ rot; eine Zeile nur aus drei Leerzeichen ⇒ grün (Leerzeile).

**Kommentare (Schritt 20).** Neu geschrieben: Kopf und Zeilen von `ci-range.sh`, `ci-gegenprobe.sh`, `ci.mk`, Zeile in `harness/README.md`, Kommentar in `.d-check.yml`. Zusagen und Test: „Bei Exit ungleich 0 bleibt stdout leer“ ⇒ `p8-fehler-stdout-leer`; „Exit 0 nicht leer, 1 leer, 2 …“ ⇒ `p7-*`, `p8-*`; Reihenfolge der Gründe ⇒ `p9-*`; „GEPRÜFT DURCH …“ ⇒ der Lauf von `make ci-gegenprobe`; „bei `SIGINT` und `SIGTERM` entfernt sie ihr Temp-Verzeichnis …“ ⇒ `p18-*` mit den Mutationen; „Kein Fall ändert Dateirechte“ (Kopf von `ci-gegenprobe.sh`; beschreibt den Code, keine Zusage über Plattformen; die Sonde (5) als root ist ein Beleg, kein Test) und „sie ruft weder Docker noch make auf“ (`ci.mk`, `harness/README.md`; beschreibt das Skript, `grep` findet keinen Aufruf) ersetzen die weiteren Sätze (Review F-5); „Aufrufarten von `make` haben keinen Fall“ ist die Grenze aus Punkt 14, kein Satz über einen Test. `ci-range.sh` nennt „Braucht bash und git; eine Version prüft das Skript nicht“; die Werkzeugliste im Kopf von `ci-gegenprobe.sh` und die Schleife der Vorprüfung nennen dieselben Werkzeuge.

**Metrik gegen die Marke von §4** (`git diff --numstat 39f1ce8a..7d80d9f2 -- .github tools harness .d-check.yml Makefile`): Summe aus Einfügungen und Löschungen **771** (Gegenprobe 624, `ci-range.sh` 75, `ci.yml` 39, Fragment 11, `.d-check.yml` 7, `harness/README.md` 1, `SPEC-064` 14 = die Änderung des Architect `0bd5d47e`, die mitzählt; ohne sie 757) gegen die Marke 900: unter der Marke, die Meldeschwelle 700 ist überschritten, keine Rückführung. Die Gegenprobe trägt 81 Prozent.

**Antworten des Architect auf die Rückgaben des Implementers (`SPEC-064` Punkte 1, 8, 9; §6 Randform 19; Commit `0bd5d47e`).** (1) Wortlaut der Meldungen von `ci-range.sh`: nicht zugesagt; zugesagt sind Präfix, Ausgang und der Hinweis aus Punkt 7 — **Grenze**; die Fälle vergleichen für Basis/Kopf nicht im Klon, Standard-Zweig fehlt und git-Fehler nur Exit und Präfix, für Variablen und flachen Klon dazu den Ort (Variable bzw. `fetch-depth: 0`), der Hinweis zum Vorgänger (`p7-push-vorgaenger-fehlt`, `p8-stdout-eine-zeile`) und die leere Range (`p7-leer`) bleiben wörtlich. (2) Basis und Kopf beide nicht aufzulösen, (5) Hinweis vor einem ungültigen Standard-Zweig: Reihenfolge innerhalb einer Stufe nicht zugesagt — **Grenze**, kein Fall; `p9-aufloesung-vor-leer` prüft nur Exit 2 mit Präfix (Basis = Kopf, nicht auflösbar) und färbt noch `m-p9-leer-vor-aufloesung`, nicht mehr die Mutationen, die die Meldung ändern. (3) Zeilen aus Leerzeichen oder Tabulatoren in `ci.yml`: nicht zugesagt — **Grenze**; kein Fall setzt die Leerzeichen-Zeile als leer voraus (`p1-leerzeile` trägt eine leere Zeile), die Messung oben bleibt Beleg ohne Zusage. (4) Tag-Kommentar: **Fall** `p1-pin-tag-leer` mit der Mutation `m-p1-pin-kommentar-ohne-text`.


### Nacharbeit nach dem Review (`docs/reviews/2026-10-19-review-slice-harness-ci.md`)

| Finding | Ausgang |
|---|---|
| F-1 MEDIUM `\|\| fehler` an `git rev-list --count` ohne Fall | behoben: Fall `p8-rev-list-fehler` (Objekt des mittleren Commits gelöscht), Mutation `m-p8-rev-list-fehler-ignoriert` (rot gesehen: `p8-rev-list-fehler`); Invariante zu `p8-git-fehler` auf den Aufruf `--is-shallow-repository` begrenzt (oben) |
| F-2 MEDIUM Aufräumen bei normalem Ende | behoben: Fall `p18-ende`, Mutation `m-p18-ende-falle-fehlt` (rot gesehen) |
| F-3 MEDIUM `CI_DEFAULT_BRANCH` bei unbekanntem Vorgänger | behoben: Fälle `p6-default-vorgaenger-fehlt-fehlt`, `-ungueltig`, Mutation `m-p6-default-vorgaenger-ungeprueft` (rot gesehen); die Mutanten der Variable und der Basis färben die neuen Fälle mit (Tabelle) |
| F-4 LOW Satz „entfällt die Mutation“ in `SPEC-064` Punkt 18 und Randform 18 | **Rückgabe an den Architect**: Vertragstext und §6 ändert der Implementer nicht; das Skript führt die Mutationen `m-p18-*-kinder-bleiben` (Nachkommen: Hintergrund-Unterprozesse) |
| F-5 LOW Sätze weiter als der Test | behoben: Kopf von `ci-gegenprobe.sh` („Kein Fall ändert Dateirechte“), `ci.mk` und Zeile in `harness/README.md` („ruft weder Docker noch `make` auf“) beschreiben den Code |
| F-6 INFO Spannung in `SPEC-064` Punkte 6, 8, 9 | geht an den Architect, keine Arbeit hier |
| F-7 INFO Ausnahme `make gates` | bestätigt; hier erneut ohne `make gates`, Begründung oben |

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

**Nachzählung nach dem Schnitt (Planner, 2026-10-18, im Commit, der die Blöcke `commits` und `vcs` an `slice-harness-ci-range-konfiguration` abgibt; `AGENTS.md` §3.13):** Die Sendung `.d-check.yml` schrumpft auf den Eintrag der Pins (`workflows` in `modules:` und `workflows.dir`, rund 5 bis 10 Zeilen). Liefer-Punkte 3 (Workflow | Range | Gegenprobe samt Gate-Index), Schichten 2 (Workflows, Harness), unverändert. Schätzung in der Metrik von §4: 390 bis 540 vor dem gemessenen Faktor, 510 bis 890 danach, unter der Marke 900. Der Nehmer der abgegebenen Blöcke ist `slice-harness-ci-range-konfiguration`; er nennt sie in §1, DoD Liefer-Punkt 1 und §3 mit der Kennung des Gebers (`grep -n slice-harness-ci` findet sie dort). Die Teilung von `SPEC-064` Punkt 14 (Gegenprobe der Konfiguration als eigenes Ziel) liegt beim Architect.

**Schätzung des Architect nach dem Schnitt (2026-10-18; Metrik von §4; `SPEC-064` vor dem Anspruch geschrieben, nicht gezählt; Faktor 1,3 bis 1,65 aus `slice-harness-ci-workflows`):** `ci.yml` 45 bis 70, `ci-range.sh` 80 bis 110, Eintrag der Pins 5 bis 10, `ci.mk` 20, `harness/README.md` 15 bis 30, Gegenprobe 240 bis 330 (Text von `ci.yml` mit je einer Mutation je Zeile, rund 30 Formen von `ci-range.sh` mit Git-Repositories); roh zusammen 405 bis 570, mit dem Faktor 530 bis 930, Mitte rund 730. Die Obergrenze liegt über der Marke 900; tritt sie ein, ist der Schnitt *Range* (Liefer-Punkt 2 mit seinen Fällen `p6-*` bis `p9-*`) und *Workflow* (Punkt 1, Fälle `p1-*`).

**Nachschätzung des Architect (2026-10-18):** Punkt 18 von `SPEC-064` fügt dieser Gegenprobe zwei Fälle mit vier Mutationen und eine Falle hinzu, rund 30 bis 45 Zeilen; die Obergrenze 930 steigt auf rund 970, die Mitte auf rund 770. Das verschärft die Rückführung *zu groß* (Schnitt in §4); der Planner entscheidet, ob er ihn vorab setzt.

**Nachrechnung beim Anspruch (Planner, 2026-10-18; Empfehlung, kein Beschluss des Nutzers):** Die Schätzung des Architect (530 bis rund 970, Mitte rund 770) trägt den Faktor 1,3 bis 1,65 schon in sich (roh 435 bis 615 nach Punkt 18). Mit dem gemessenen Faktor der Nachbarn (1,1 bis 1,65; `slice-harness-ci-range-konfiguration` 445 gegen 300 bis 415, `slice-harness-ci-workflows` 899 gegen 545 bis 695) auf den rohen Wert liegt der Slice bei 480 bis 1015, Mitte rund 750, unter der Marke 900; die Obergrenze liegt darüber. Die Mitte liegt unter der Marke mit rund 150 Zeilen Spielraum, daher kein Schnitt vorab: ein Slice, der Schnitt *Range* gegen *Workflow* bleibt die Rückführung nach §4. Liefer-Punkte 3, Schichten 2, unverändert (kein neuer Eintrag einer Sendung).
