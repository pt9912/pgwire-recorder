# Welle welle-erster-release — Erster Release: CI auf GitHub, Veröffentlichung von Tag, Binaries und Image, Scan, Freigabe — Closure-Notiz

**Welle:** welle-erster-release
**Abschluss:** 2026-10-19
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3.

- Der erste stabile Tag `v0.1.0` ist veröffentlicht (`LH-QA-03`, `LH-FA-16`): sechs Binaries mit `SHA256SUMS` am GitHub-Release, das Image `pt9912/pgwire-recorder` für `linux/amd64` und `linux/arm64` in `ghcr.io` und `docker.io` mit gleichem Index-Digest `sha256:b45c4044fe4313934da6fc36095692642b65d21a9945231c2aa3162436c27e05`, `latest` in beiden, Smoke aus Abnahmeszenario 9 aus beiden Registries bestanden. Der Eintrag der Freigabe-Checkliste steht unten.
- 6 Slices in `done/`, alle als Datei dort (Abzählung `ls` und `grep` auf `Welle:`, Lauf des Planners am 2026-10-19): `slice-harness-ci-workflows` (Gate der Workflows), `slice-harness-ci-range-konfiguration` (die Range-Prüfungen wirken), `slice-harness-ci` (`ci.yml` fährt `make gates` und prüft die Range), `slice-erster-release-veroeffentlichung` (Tag und Binaries), `slice-erster-release-veroeffentlichung-image` (Image, Veröffentlichung mit gleichem Digest, `latest`), `slice-erster-release-freigabe` (Scan, Freigabe-Checkliste, erster Tag, Handbuch zum Bezug).
- Die CI auf GitHub: `ci-gates` fährt `make gates` und die Range-Prüfung auf jedem Push und Pull-Request (grün auf `35778a2f`, Lauf `39298463448`, und auf `1f501a6c`, Lauf `39328402488`); die Workflow-Dateien bestehen ein Gate (`make workflows-check`).
- Nicht in dieser Welle (Entscheidung des Nutzers): die Homebrew-Formel und Abnahmeszenario 11 in `welle-v1-nachtrag`, die Nachtläufe in `welle-ci-nachtlaeufe`.

## Eintrag der Freigabe-Checkliste des ersten Releases (`v0.1.0`)

Form: Tabelle aus `docs/maintainer/releasing.md` §10 (elf Punkte) mit ausgefüllter Spalte *Beleg*; was ein Lauf nicht zeigt, steht als Grenze. **Getaggter Commit:** `35778a2f4059237d309b2c71e194392b9c565a0e` (trägt `VERSION` 0.1.0 und die Zeile 0.1.0 der Änderungshistorie des Handbuchs), Tag `v0.1.0`, Datum 2026-10-19. Quellen: vom Planner gelesen (`gh run view`, `gh release view`, `git show`, `skopeo inspect --config`, Download von `SHA256SUMS`) oder Angabe der Koordination aus der Sitzung des Tags; die Angaben der Koordination sind mit „Angabe“ markiert, der Planner hat sie nicht nachgefahren.

| Punkt | Beleg |
|---|---|
| 1. `make gates` grün auf dem getaggten Commit | `35778a2f`, 2026-10-19. Angabe: im Haupt-Checkout Exit 0 im ersten Versuch; im Integrations-Worktree Exit 0 erst im dritten Versuch, die Versuche 1 und 2 rot durch Zeitgrenzen unter Last (`TestE2EPlayVergleichAbweichung`, `DROP DATABASE`; Host-Last 15 bis 22 durch fremde Container; Register `BEO-REPO/test-faellt-unter-cpu-last-aus`). Dazu `ci.yml` auf dem Stand: Lauf `39298463448`, Job `ci-gates` success, 9 min 53 s. |
| 2. Gegenproben-Container grün für diesen Hash | **Das Ergebnis des Hooks war rot:** Ablage `gegenproben-main-35778a2f…`, Beginn 17:57, Ende 19:36 UTC, rot waren `a-check-negativ`, `ci-gegenprobe`, `test-container-gegenprobe`, `gegenproben-gegenprobe`, `gegenproben-nach-einlieferung-gegenprobe`, `integration-volumes-gegenprobe`, `test-postgres-versionen-gegenprobe`, `range-konfiguration-gegenprobe` (Angabe). Deshalb fuhr die Freigabe, wie der Punkt es vorsieht, `make gegenproben` selbst per `ssh` im Container `pgr-gp` auf `35778a2f`: Exit 0, alle Ziele grün, Log `/tmp/gp-manuell.log` im Container, Lauf von etwa 19:40 bis etwa 21:00 UTC (Angabe). **Grenze:** Die Datei der Ablage für diesen Hash trägt `ergebnis=rot`; der Beleg ist der Handlauf, dessen Log im Container liegt und vom Planner nicht gelesen ist. Die Ursache des roten Hook-Laufs ist nicht geklärt (Umgebung des Hintergrundlaufs ignorierte `SIGINT`, `a-check` meldete „permission denied“ auf `/src`; Register `BEO-REPO/gegenproben-hook-lauf-rot-handlauf-gruen`; Befund bei `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`). |
| 3. `make gegenproben-host` am Host grün auf demselben Hash | 2026-10-19, `35778a2f`, Exit 0, `container-start-gegenprobe` grün (Angabe). |
| 4. `VERSION` und die Zeile der Änderungshistorie im getaggten Commit | `git show 35778a2f:VERSION` ergibt `0.1.0`; `git show 35778a2f:docs/user/benutzerhandbuch.md` führt in Zeile 2113 `\| 0.1.0 \| 2026-10-19 \| Erste veröffentlichte Version: …` (Planner, 2026-10-19). |
| 5. Trockenlauf auf dem Stand | Lauf `39324550440` (`workflow_dispatch`, `main` = `35778a2f`, 2026-10-19 21:04 UTC), success, 3 min 4 s; Ausgabe `release: eingaben: grün für v0.1.0 (Trockenlauf)`, die Zeile der Änderungshistorie gefunden (`release: hinweis: Änderungshistorie 0.1.0: …`), kein Hinweis auf eine fehlende Zeile. Vorher: `39263455592` rot (Trivy-Annahme zum Layout, siehe unten), `39284378634` rot (Scan fand drei HIGH-Befunde im Go-Builder 1.27.1), `39290074657` grün. |
| 6. Scan grün | `release: scan: 0 Befund(e) über der Schwelle, 0 ohne Behebung` im Log des Trockenlaufs `39324550440` und im Log des Laufs des Tags `39325091059` (Planner). **Grenze:** Ein grüner Scan sagt „nichts Bekanntes in diesem Raum“ (Schwelle `HIGH` und `CRITICAL` mit Behebung, [ADR-0052](../../adr/0052-scan-des-images-vor-dem-push.md)); `MEDIUM` und Befunde ohne Behebung sperren nicht. |
| 7. Ausdrückliches Ja des Nutzers zum Tag | Nachtrag vom 2026-10-18: „ja für den tag und publishing“; am 2026-10-19: „bitte ohne anzuhalten das Release erstellen“ (Wortlaute nach Angabe der Koordination). Ein Folge-Tag `v0.1.1` war erlaubt und wurde nicht gebraucht; ein Probe-Tag wurde nie benutzt. |
| 8. Lauf des Tags | `39325091059`, Tag-Push `v0.1.0` 2026-10-19 21:08:16Z, Ende 21:12:49Z (Job 4 min 30 s), alle Schritte grün. Index-Digest `sha256:b45c4044fe4313934da6fc36095692642b65d21a9945231c2aa3162436c27e05` gleich in `ghcr.io` und `docker.io` (`release: vergleichen: … in ghcr.io und docker.io`), `release: smoke: grün aus ghcr.io`, `release: smoke: grün aus docker.io`, `release: latest: ghcr.io zeigt auf sha256:b45c4044…` und dasselbe für `docker.io`. |
| 9. Anhänge am GitHub-Release | https://github.com/pt9912/pgwire-recorder/releases/tag/v0.1.0, kein Entwurf, kein Vorab-Release, veröffentlicht 2026-10-19 21:12:46Z, sieben Dateien: `pgwire-recorder_0.1.0_darwin_amd64`, `…_darwin_arm64`, `…_linux_amd64`, `…_linux_arm64`, `…_windows_amd64.exe`, `…_windows_arm64.exe` und `SHA256SUMS` (Planner, `gh release view`). `SHA256SUMS` nennt für `linux_amd64` die Summe `fe53dabf4c48eac30221c47c99d40cf7839bd6d852a345b664d9078e5d69b5d4` (Planner, heruntergeladen). |
| 10. Beobachtungen des ersten Laufs | Siehe unten, je Beobachtung Lauf-Kennung und Zeile oder die Grenze. |
| 11. Bezug nach dem Tag | Commit `1f501a6c` (Kopf *Software-Version* 0.1.0, §2 *Installation* mit Bezug der Binaries und der Images, `README.md`). Gelaufen gegen die veröffentlichten Artefakte (Angabe): `curl` von `SHA256SUMS` und `pgwire-recorder_0.1.0_linux_amd64`, Summe `fe53dabf…b5d4`, `sha256sum -c --ignore-missing SHA256SUMS` ohne Abweichung, `./pgwire-recorder_0.1.0_linux_amd64 version` gibt `pgwire-recorder 0.1.0` aus; `docker pull` und `docker run … version` aus `ghcr.io` und `docker.io` (Tags `0.1.0` und `latest`) geben dieselbe Ausgabe und denselben Digest. **Grenzen:** `darwin`, `windows` und `linux/arm64` sind gebaut, nicht ausgeführt (Handbuch §2 sagt nur das); das Limit anonymer Pulls bei `docker.io` ist nicht gemessen. |

**Beobachtungen des ersten Laufs (Punkt 10)**, Quelle das Log von `39325091059` (Planner), sonst genannt:

- Wortlaute von `skopeo inspect` für einen nicht vorhandenen Tag: `ghcr.io`: `reading manifest 0.1.0 in ghcr.io/<Eigentümer>/pgwire-recorder: manifest unknown` (Zeile `release: skopeo: ` im Schritt `Veröffentlichen`, 21:11:46Z); `docker.io`: `reading manifest 0.1.0 in docker.io/<Eigentümer>/pgwire-recorder: manifest unknown` (21:12:04Z). Beide fallen unter das Muster `manifest unknown`; `name unknown` und `not found` blieben unbenutzt. `skopeo` und `gh` waren auf dem Runner verfügbar.
- `gh release view` für einen nicht vorhandenen Release: Das Muster `release not found` griff, der Schritt legte den Release an (`release: github-release: v0.1.0 angelegt`, 21:12:46Z). Der Wortlaut steht nicht im Log, das Skript gibt ihn nicht aus.
- Schreibrecht der Token in beiden Registries (`release: veroeffentlichen: 0.1.0 nach ghcr.io kopiert`, `… nach docker.io kopiert`), Sichtbarkeit des Pakets in `ghcr.io` (anonymer Zugriff im Schritt `Vergleichen`; Angabe: auch von Hand mit leerem `DOCKER_CONFIG` gezogen), das Repository auf Docker Hub besteht. Smoke aus beiden Registries mit der Referenz `<Registry>/pt9912/pgwire-recorder@<D>`, Digest je Registry gleich `D`.
- Zeitstempel (Beobachtung zur Wirkung von `SOURCE_DATE_EPOCH`, keine Zusage): `created` der Konfiguration von `linux/amd64` im veröffentlichten Image ist `2026-10-19T17:36:26Z`, die Committer-Zeit von `35778a2f` (Planner, `skopeo inspect --config` und `git show`). Der Index-Digest des Trockenlaufs `39324550440` (Bau 21:07) und des Laufs des Tags (Bau 21:10) ist derselbe: zwei Baue desselben Commits mit gleichem Ergebnis. **Grenze:** Das Image bleibt nicht als reproduzierbar zugesagt (`SPEC-067` Punkte 5 und 15); zwei Baue sind keine Zusage.
- Wiederanstoß eines roten Laufs nach `releasing.md` §11: **nicht belegt**, der Lauf des Tags war grün; Register `BEO-REPO/wiederanstoss-eines-roten-release-laufs-unbelegt`.
- Auffälligkeiten im Log von `39325091059` (Suche nach `warn`, `##[error]`, `denied`, `rate limit`, `deprecat`): keine.

**Freshness-Audit:** nicht durchgeführt; kein Hindernis. Es ist optional und kein Punkt der Checkliste (Entscheidung des Nutzers vom 2026-10-19, `SPEC-067` Punkt 17), ein Ergebnis sperrt den Tag nie. Hinweis des Nutzers: Der neuere Stand der Baseline (v6.18.1, 2026-10-15) ist keine Sicherheitskorrektur. Eine Tabelle als Beleg entfällt, weil nichts durchgeführt wurde.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3.

- Der Trockenlauf statt eines Probe-Tags: Drei Trockenläufe fanden zwei Fehler (Layout für Trivy, Go-Builder 1.27.1) ohne Außenwirkung; der erste Lauf mit Tag war grün, ein Folge-Tag nicht nötig.
- Der Scan vor dem Push hat seine Wirkung belegt: Er lehnte den Builder mit drei echten HIGH-Befunden ab, bevor eine Registry etwas erhielt.
- Die Schnitte von `slice-erster-release-veroeffentlichung` (Tag und Binaries gegen Image) und die Sendungen an die Freigabe: Jeder Beleg, den der Trockenlauf nicht erzeugt, war als Sendung eingetragen und wurde vom Lauf des Tags erzeugt oder als Grenze genannt.
- Die Checkliste mit dem Fall „Hook-Ergebnis rot, Handlauf“ hat den roten Hook-Lauf getragen, ohne den Tag zu stoppen und ohne das Rot zu verschweigen.

## Was ging anders als geplant?

- Eine Annahme über ein Fremdwerkzeug im Vertrag (Trivy liest das Layout als Tar mit `--platform`) war falsch und zeigte sich erst im ersten Trockenlauf; eine Behebungsrunde und drei Review-Runden folgten (Register `BEO-REPO/vertrag-nennt-aufruf-eines-fremdwerkzeugs-ungemessen`). Konsequenz: Vorschlag in der Closure von `slice-erster-release-freigabe` (Lücke 1), nicht eingetragen.
- Das Freshness-Audit wechselte von einem Punkt der Checkliste (zwölf Punkte) zu einer optionalen Handlung (elf Punkte); Entscheidung des Nutzers vom 2026-10-19, `SPEC-067` Punkt 17.
- `make gates` auf dem getaggten Hash war im Worktree erst im dritten Versuch grün (Last), der Hook-Lauf der Gegenproben war rot, der Handlauf grün; beides steht im Register.
- Ein Push von `main` erfolgte vor dem grünen `make gates` für den Stand (`80abcc06`); die Gates waren danach grün. Befund im Register (`BEO-REPO/push-auf-hauptzweig-vor-gruenem-gates-lauf`), Vorschlag in der Closure des Slice.
- Der Slice `slice-erster-release-veroeffentlichung` wurde am 2026-10-19 in Tag und Binaries gegen Image geschnitten (Drift-Log der Roadmap); die Größe der Freigabe lag am Ende bei 805 von 1200.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3 (nur 3× oder mehr). Diese Closure verkörpert keine Regel: Änderungen in `AGENTS.md`, `.claude/` und `.harness/skills/` brauchen das Ja des Nutzers; was der Lese-Schritt findet, steht als Vorschlag.

Lese-Schritt über das Register (Zählung der `evidence/`-Dateien, 2026-10-19). Einträge, in denen Slices dieser Welle Evidence tragen:

- Bereits verkörpert, bei 5× oder mehr, ohne neuen Urteilsbedarf gegenüber der Closure von welle-v1-abschluss: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (87), `BEO-REPO/plan-folgt-korrektur-nicht` (74), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (62), `BEO-REPO/spec-randform-erst-im-review-entschieden` (45), `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (26), `BEO-REPO/gruener-mutant-faelschlich-aequivalent` (9), `BEO-REPO/groessenprobe-ohne-messbaren-stand` (5). Die Sensoren und ihre Grenzen stehen in der Closure von welle-v1-abschluss (`make kopf-check` für die Kopf-Hälfte von `plan-folgt-korrektur-nicht`, `slice-harness-mutation` für die Hälfte „Test fängt die Zusage“, sonst Urteil); diese Welle fügt Auftreten derselben Formen hinzu, keine neue Form. Aufgetaucht ist die Form „Zusage ‚nur X‘ ohne Grenzmutation“ (drei Auftreten in drei Slices: `slice-erster-release-veroeffentlichung` und `slice-erster-release-veroeffentlichung-image` aus dieser Welle, `slice-harness-doc-trace-archiv` wellenlos; Vorschlag an den Nutzer für `AGENTS.md` §3.10 in deren Closure-Notizen, nicht eingetragen).
- Ausgang `geplant` besteht: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (18) → `slice-harness-schaetzung-planner-faktor` (Messwert Faktor 5,0 der Gegenprobe des Scans im Register vermerkt), `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (3) → `slice-harness-regeln-gate-konfiguration`.
- **Schwelle 3× erreicht, ohne Ausgang (Urteil des Nutzers)** — Vorschlag des Planners, keine Zuweisung:
  - `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (3, erreicht mit `slice-erster-release-freigabe`; die ersten zwei Vorgänge waren Umbauten ohne Verhaltensänderung, dieser liefert einen neuen Vertrag). Auslöser: `slice-lint-bestand-driving`, `slice-lint-bestand-kern-driven`, `slice-erster-release-freigabe` — 3×. Vorschlag: `geplant`, der Sensor `slice-harness-mutation` (steht in `open/`) fängt „Wächter ohne Fall“ mechanisch; bis dahin kein Satz, der „keine grünen Mutanten“ ohne die Liste der Wächter der geänderten Funktionen zulässt. Zielort einer Prosa-Regel wäre `AGENTS.md` §3.10, Herkunfts-Anker `seit welle-erster-release`; ob verkörpert oder geplant, entscheidet der Nutzer.
  - `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` (3, seit `slice-erster-release-veroeffentlichung`, Ausgang beim Nutzer): In dieser Welle folgten drei Gegenbelege ohne Auftreten (`slice-harness-doc-trace-archiv`, `slice-erster-release-veroeffentlichung-image`, `slice-erster-release-freigabe`; die Verifikation brach jede Nacharbeit von Hand). Das ist ein Hinweis für die Entscheidung des Nutzers, kein Grund für `gestrichen` (die Begründung, dass es nicht mehr auftreten kann, trägt er nicht); die Zeile bleibt `offen`.
  - `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (8, Ausgang beim Nutzer und beim Lese-Schritt von `welle-gegenproben-und-mutanten`): In dieser Welle zwei Auftreten (`slice-harness-ci-workflows`, `slice-harness-ci-range-konfiguration`); der rote Hook-Lauf dieser Freigabe steht getrennt als `BEO-REPO/gegenproben-hook-lauf-rot-handlauf-gruen` (1×), Zuordnung nach der Ursache.
- Unter der Schwelle, dieser Welle zugeordnet: `BEO-REPO/test-faellt-unter-cpu-last-aus` (2), `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder` (2), `BEO-REPO/zeilenende-crlf-im-gelesenen-wert` (2), `BEO-REPO/gegenprobe-startet-sich-ueber-einen-fall-selbst` (2); neu mit 1×: `BEO-REPO/vertrag-nennt-aufruf-eines-fremdwerkzeugs-ungemessen`, `BEO-REPO/push-auf-hauptzweig-vor-gruenem-gates-lauf`, `BEO-REPO/gegenproben-hook-lauf-rot-handlauf-gruen`, `BEO-REPO/wiederanstoss-eines-roten-release-laufs-unbelegt`.

## Trigger-Audit

- Carveouts: 0 offen (`docs/plan/carveouts/` trägt nur `.gitkeep`; Lauf des Planners 2026-10-19: `ls -la`).
- Bootstrap-aware Gates: keine Reifestufe im Repo, 0 offen (`harness/README.md` §Sensors führt keine Stufe).
- ADR-Re-Evaluierungs-Trigger (52 ADRs, die Abschnitte *Re-Evaluierungs-Trigger* der berührten gelesen). Urteil zu den berührten:
  - [ADR-0050](../../adr/0050-ci-auf-github-actions.md) (CI auf GitHub Actions): kein Lauf auf GitHub endete rot, den ein Gate hätte fangen sollen (`ci.yml` grün auf `35778a2f` und `1f501a6c`); die Zeitgrenze trägt `make gates` (9 min 53 s); ein Limit des Docker Hub färbte keinen Lauf — nicht eingetreten, bestätigt.
  - [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md) (Veröffentlichung ohne Probe-Tag): Der erste Lauf mit Tag scheiterte nicht; der Neubau des Images verhinderte den Wiederanlauf nicht (gleicher Digest zweier Baue); die Registries verlangten keine Tags `<MAJOR>.<MINOR>` — nicht eingetreten, bestätigt. Der Wiederanlauf eines roten Laufs blieb ungeprüft (Register).
  - [ADR-0052](../../adr/0052-scan-des-images-vor-dem-push.md) (Scan mit Trivy): Der rote Lauf `39284378634` war ein echter Befund, kein Fehlalarm; der Scan scheiterte einmal an einer falschen Annahme über das Layout (`39263455592`), nie am Netz oder an der Datenbank — nicht eingetreten, bestätigt.
  - [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (volle und schnelle Kette): Der dritte Trigger („wenn eine CI die volle Kette je Änderung fährt und der Hook nur noch lokal schützt“) ist eingetreten, weil `ci.yml` die volle Kette je Push fährt, und durch [ADR-0050](../../adr/0050-ci-auf-github-actions.md) beantwortet: Die CI ergänzt den Stop-Hook, sie ersetzt ihn nicht. Der Push vor dem grünen Lauf (Register) zeigt, dass die CI nachträglich prüft, nicht freigibt.
  - [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (Gegenproben verlassen die Kette): Die Closure fand ein Ergebnis des Gegenproben-Containers für den Release-Hash (rot) und fuhr danach selbst; „kein Ergebnis gefunden“ trat nicht ein, eine Regression blieb nicht unbemerkt — nicht eingetreten, bestätigt. Die erste Fahrt in `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` bleibt der Beleg der Fahrstelle.
  - [ADR-0043](../../adr/0043-binaries-und-container-gate.md) (Binaries und Container): Reproduzierbar verlangt der Vertrag das Binary, nicht das Image (`SPEC-067` Punkt 5, `LH-QA-01`); keine Plattform kam hinzu oder fiel weg — nicht eingetreten.
  - Übrige ADRs: Trigger nicht berührt (Urteil nach Lesen der Bedingung für die genannten, Stichprobe für den Rest).
- Hard Rules (`AGENTS.md` §3.2 bis §3.14): alle tragen den Retirement-Check, keine hat einen Auflösungs-Trigger, der eintrat.
- MR: `MR-002` (permanent bis das Kurs-Repo einen Ort nennt) — Trigger nicht eingetreten.

## Bestand in `open/` und `next/` (Vorschlag, Urteil des Nutzers)

27 Dateien in `open/` (18) und `next/` (9), Gruppierung vor dem Einzelurteil (`modul-06-roadmap.md` Schritt 3). Seit der Closure von welle-v1-abschluss sind neun neu (siehe unten), 18 haben diese Closure und die davor unverändert überstanden:

- Gruppe Mutanten (Auslöser: grüne Mutanten): `slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`, `slice-harness-mutation` — bestätigt; nicht konsolidiert, weil Test-Lieferung und Gate verschiedene Liefer-Punkte sind. `slice-harness-mutation` ist jetzt auch der benannte Sensor für `BEO-REPO/liste-gruener-mutanten-unvollstaendig`.
- Gruppe Regeln: `slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`, `slice-harness-schaetzung-planner-faktor` — bestätigt (Register-Ausgänge `geplant`); neu `slice-harness-regeln-gate-konfiguration` (Ausgang `geplant` für `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`) — bestätigt.
- Gruppe Doku-Gate: `slice-harness-lh-links-bestand`, `slice-harness-lh-links-pflicht` — bestätigt; Bestand zuerst.
- Gruppe Nachtläufe (`welle-ci-nachtlaeufe`, Start ist diese Welle in `done/`): `slice-harness-ci-gegenproben`, `slice-harness-ci-freshness`, `slice-harness-ci-freshness-images` — bestätigt, startbar. `slice-harness-ci-freshness-images` und `slice-harness-ci-freshness` haben jetzt einen Messwert: der Pin des Go-Builders alterte, der Scan der Pipeline war der einzige Melder (`BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder`); ob die Nachtläufe Pins in Workflows und `Dockerfile` annehmen, steht in deren §1 und ist ihre Entscheidung.
- Gruppe Gegenproben-Betrieb (`welle-gegenproben-und-mutanten`): `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (trägt den Befund zum roten Hook-Lauf), `slice-harness-nachweis-uebernehmen`, `slice-harness-nachweis-uebernehmen-aufrufarten`, `slice-harness-gegenprobe-baumkopie` — bestätigt.
- Gruppe Dev-Container: `slice-harness-devcontainer-einrichten`, `slice-harness-devcontainer-einrichten-zugang` — bestätigt (`welle-dev-container`).
- Gruppe Nachtrag (`welle-v1-nachtrag`): `slice-erster-release-homebrew-nachweis` — bestätigt; seine Start-Bedingung (`slice-erster-release-freigabe` in `done/`) ist mit dieser Closure erfüllt, der Release mit `SHA256SUMS` liegt vor. `slice-v1-abschluss-homebrew`, `slice-v1-abschluss-cancel-ohne-schluessel`, `slice-v1-abschluss-verbindungsende-upstream-extended` — bestätigt.
- Einzeln bestätigt, je mit Register-Anker oder Trigger: `slice-harness-commit-struktur-id`, `slice-harness-lint-warten-ohne-frist`, `slice-harness-slice-mv-exakter-treffer`, `slice-harness-gegenprobe-matrix-p7-flanke`.
- `slice-harness-gate-index-werkzeug-teil` — bestätigt, Start-Bedingung (§4) nicht eingetreten, geprüft am 2026-10-19: (1) Unter `harness/mk/` liegt kein Teil des Gate-Index (`harness/conventions.md` führt „Teil des Gate-Index eines Werkzeugs“ als *kommt nicht mit*); (2) `welle-erster-release` ist mit dieser Closure `done`, `welle-v1-nachtrag` nicht. Der Slice bleibt in `open/`.
- Konsolidiert: keiner. Entfallen: keiner. Die neun neuen Dateien sind nicht Gegenstand dieses Urteils, sie haben keine Closure überstanden; sie stehen in den Gruppen oben.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register.

Der Zähler steht in [`../observations/`](../observations/). Was in dieser Welle **3×** erreicht hat, steht oben unter *Steering-Loop-Einträge*.

## Folge-Slices

- `slice-erster-release-homebrew-nachweis` (Start-Bedingung erfüllt; `welle-v1-nachtrag`), `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (trägt den Befund zum Hook-Lauf), `slice-harness-schaetzung-planner-faktor` (Messwert), `slice-harness-regeln-gate-konfiguration` (Ausgang `geplant`), `slice-harness-mutation` (benannter Sensor); jeder existiert als Datei im Planning-Lifecycle (`next/` oder `open/`).

## Archivierung

- `make archive-welle WELLE=welle-erster-release` (Träger `.harness/state/bin/ai-harness-init`) hat 7 Slice-Dateien (6 Mitglieder, 1 wellenloser seit der letzten Closure: `slice-harness-doc-trace-archiv`) und den Welle-Plan nach `done/welle-erster-release/archiv.zip` (171792 Bytes) archiviert, Stubs daneben; zwei Commits (reiner Move, dann Archiv, Stubs und Verweis-Nachzug über 12 Dateien). Review-Reports archivierte der Lauf keine (Zählung des Werkzeugs: 0); `docs/reviews/` bleibt unverändert. Die Zeile `Bezug` der Stubs schrieb `make archiv-bezug WELLE=welle-erster-release` (7 Stubs, eigener Commit).
- Geltungsbereich der Sensoren nach dem Lauf: `make docs-check` ist grün (0 Befunde); der Commit-Träger sucht eine Slice-Kennung mit `docs/plan/planning/*/<kennung>.md` und findet die Stubs; `make doc-trace` liest die Stubs über die Zeile `Bezug` (die Zeilen zu `LH-FA-16` und `LH-QA-03` nennen `slice-erster-release-freigabe` weiter).
- Der Lauf hat keinen Pfad in `.claude/` oder `AGENTS.md` nachgezogen (`git diff` über beide leer).

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 1.

- 6 von 6 Slices in `done/`: belegt (`ls docs/plan/planning/done/<kennung>.md` für jede der sechs Kennungen, 2026-10-19; kein Slice mit `Welle: welle-erster-release` in `open/`, `next/` oder `in-progress/`).
- `make gates` grün auf `35778a2f` (siehe Eintrag der Checkliste, Punkt 1) und `ci.yml` grün auf `35778a2f` (`39298463448`) und auf `1f501a6c` (`39328402488`); die Commits dieser Closure ändern nur Pfade der Positivliste (`docs/plan/planning/`), `make docs-check` ist grün (siehe Archivierung für den Endstand).
- Der erste stabile Tag ist nach der Freigabe-Checkliste veröffentlicht: Binaries mit `SHA256SUMS` am Release, Image in beiden Registries mit gleichem Index-Digest, Smoke aus Abnahmeszenario 9 bestanden (Eintrag oben, Punkte 8 und 9); die Pipeline-Läufe sind nachgewiesen (`ci.yml` auf dem Stand des Tags und der Lauf der Veröffentlichung `39325091059`).
- Replay-Lauf: kein eigener Lauf dieser Welle; die Welle ändert kein Verhalten des Binaries (die Abnahmeszenarien 1 bis 10, 12 bis 16 und 18 verifizierte die Closure von welle-v1-abschluss). Ein Verifikations-Bericht der Abnahmeszenarien dieser Welle liegt nicht vor; der Beleg für Abnahmeszenario 9 ist der Smoke des Laufs des Tags.
