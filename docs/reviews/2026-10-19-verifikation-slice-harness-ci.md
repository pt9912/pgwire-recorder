# Verifikation: slice-harness-ci

**Rolle:** Verifier. **Datum:** 2026-10-19. **Spanne:** Basis `39f1ce8a` bis `fc63d5ad` (Zweig `impl-ci`). **Gegenstand:** [slice-harness-ci](../plan/planning/done/slice-harness-ci.md) gegen DoD, [SPEC-063](../../harness/vertraege/SPEC-063-workflows-check.md), [SPEC-064](../../harness/vertraege/SPEC-064-ci.md), [ADR-0050](../plan/adr/0050-ci-auf-github-actions.md) und `AGENTS.md` §3.9 bis §3.14.

Die Review-Reports (Zweig `rev-ci`, Runde 1 und 2) sind nicht als Beleg gewertet; alles unten ist in diesem Lauf selbst erzeugt. Der Code ist seit `7d80d9f2` unverändert (`git diff 7d80d9f2..fc63d5ad -- tools harness .github .d-check.yml` zeigt nur `SPEC-064`, 4 + 4 Zeilen), die Belege des Implementers gelten damit für den Endstand.

**Engster Sensor.** Gefahren: `bash -n` (beide Skripte), `make docs-check` (1196 Dateien, 0 Befunde), `make kopf-check` (Exit 0), `make workflows-check` (1 Datei, kein Befund), `make abdeckung-check`, `make baseline-verify`, `make meldungskatalog-check` (je Exit 0), `make gates-schnell-gegenprobe` (grün, Namensregel von `SPEC-060` Punkt 2 mit `ci-gegenprobe` in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS`; `namensregel` am Arbeitsbaum ohne Verletzung), Teilläufe von `ci-gegenprobe` (`GP_MUT=ok`: 1 Lauf grün; `^m-p[6-9]-`: 38 Läufe grün; `^m-p18-`: 7 Läufe grün), dazu eigene Mutanten (unten). Nicht wiederholt: volles `make gates`, volle `ci-gegenprobe`.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1 Workflow: `ci.yml` hält `SPEC-064`, besteht `workflows-check` und `docs-check`, Eintrag `workflows` in `.d-check.yml`, Randformen vor dem Code entschieden, Lauf auf GitHub | Code-Hälfte bestätigt; Lauf auf GitHub **offen (erwartet, §5)** | `ci.yml` Zeile für Zeile gegen den Text von Punkt 1: gleich (`p1-gleich` grün, 39 Zeilen, LF, kein Tab). `.d-check.yml`: `workflows` in `modules:` und `workflows.dir: .github/workflows`. Pin `actions/checkout@3d3c42e5…ba90b1 # v7.0.1`: `git ls-remote` am Host (zweite Ermittlung, unabhängig vom Container des Implementers) nennt `refs/tags/v7.0.1` mit demselben SHA. Randformen: 1 bis 18 vor dem ersten Code-Commit `4190a5c3` entschieden; Randform 19 (Meldungstexte, Reihenfolge innerhalb der Stufe, Leerraumzeile, Tag-Kommentar) kam über die Rückgabe des Implementers und steht in `0bd5d47e`, nach dem ersten Code-Commit (Randform-Rückgabe nach `AGENTS.md` §3.12, kein Verstoß; der Wortlaut der DoD „vor dem ersten Code-Commit“ trifft für 19 nicht, siehe V-4). Grüner Lauf und Dauer auf GitHub: fehlen, siehe Risiken. |
| 2 Range: `ci-range.sh` je Randform, Fehler nie mit nichtleerer Range | bestätigt | Skript gegen Punkte 6 bis 9 gelesen: Stufen (1) Ereignis, (2) Variablen in der Folge, (3) flacher Klon, (4) Auflösung, (5) leere Range; `set -uo pipefail`, keine Pipeline (kein `pipefail`-Effekt); stdout bei Erfolg eine Zeile, jede Meldung `ci-range:` auf stderr. Eigene Mutanten, siehe unten: alle trennscharf rot außer zwei äquivalenten. Die Zuweisung `range="$(…)"` steht allein in der Zeile (Punkt 2), `run:` läuft unter `bash -e`, ein Exit 1 oder 2 beendet den Schritt. |
| 3 Gegenprobe: je Zusage von Teil 1 und 2 und je Element ein Fall und eine Mutation; unveränderte Kopie grün; Ziel in `GEGENPROBEN_CONTAINER`; Zeile in `harness/README.md` | bestätigt, mit V-1 (LOW) | Zählung 148 Fälle / 144 Mutationen: 38 und 7 am Lauf bestätigt (`m-p[6-9]` und `m-p18`), 99 aus dem Skript (12 + 87), Summe wie §3 und §7. Zeile in `harness/README.md` nennt Branch-Schutz als Handlung des Betreibers, die CI als Ergänzung des Stop-Hooks, den Namen `ci-gates`. Aufrufarten von `make`: kein Fall (Punkt 14). |
| 4 `make gates` und `make gates-schnell-gegenprobe` grün, `ci-gegenprobe` auf dem Endstand | bestätigt für Gegenprobe und schnelle Gegenprobe; `make gates` Nachweis gilt für `a31f4a7f` | siehe Bewertung zu Punkt 6 der Aufgabe unten. |
| 5 Review durchgeführt | bestätigt | Report auf `rev-ci` (`bc57440c`); drei MEDIUM der Runde 1 in Runde 2 behoben; ich habe F-1 bis F-3 über die Fälle und Mutanten gegengeprüft (`p8-rev-list-fehler`, `p18-ende`, `p6-default-vorgaenger-fehlt-*` laufen grün in `ok`). Der Report liegt noch nicht auf `impl-ci`. |
| 6 bis 9 Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen (erwartet) | Alle Häkchen leer; §7 trägt nur Belege des Implementers und die Nacharbeit; die Ausgänge der Risiken stehen als „vorgesehen“. Closure folgt nach dem Probe-Push. |

## Mutationen: Test ohne den Fix aus dem richtigen Grund rot

Eigene Mutanten, in Kopien (Skripte unter dem Scratchpad), jeweils `GP_MUT=ok`. Ein Mutant von `ci-range.sh` muss die unveränderte Kopie färben:

| Mutant (Textersatz in `ci-range.sh`) | rote Fälle | Urteil |
|---|---|---|
| `HEX` ohne `^` | `p6-{before,head,prbase}-lang` | richtig |
| `HEX` mit `[0-9a-fA-F]` | `p6-*-gross` | richtig |
| Ausgabe `basis...kopf` | `p7-pr`, `p7-push`, `p7-push-neu`, `p8-stdout-eine-zeile` u. a. | richtig |
| leere Range Exit 3 statt 1 | `p7-leer` | richtig |
| `pruefe_standard` bei Nullen entfernt | `p6-default-*`, `p7-push-neu*`, `p8-fehler-stdout-leer` | richtig |
| flacher Klon nicht erkannt | `p7-flach`, `p9-flach-vor-aufloesung` | richtig |
| Basis des Pull-Requests nicht im Klon ohne Exit 2 | `p7-pr-basis-fehlt-im-klon`, `p8-fehler-stdout-leer` | richtig |
| `^{commit}` entfernt | `p6-default-vorgaenger-fehlt-*`, `p7-push-vorgaenger-fehlt`, `p8-stdout-eine-zeile` | richtig (ein 40stelliger Hex-String „löst“ sonst ohne Objekt auf) |
| `-eq 0` zu `-lt 0` | `p7-leer`, `p8-fehler-stdout-leer` | richtig |
| `set -o pipefail` entfernt | keiner | äquivalent: Invariante „das Skript enthält keine Pipeline“; Ort: der Quelltext (`grep ' | '` findet nur die Alternative `pull_request | push` eines `case`); Sonde: `bash -n` und die Fälle unverändert |
| `export LC_ALL=C` entfernt | keiner | Zusage nicht gegeben (der Bereich `[0-9a-f]` hängt nur in Locales mit Kollationsbereichen davon ab); kein Fall, siehe V-3 |

Mutanten der Prüfung des Textes (`pruefe_ci`):

| Mutant | Ergebnis |
|---|---|
| Längenprüfung am Ende (`${#ist[@]} -eq $j`) entfernt | rot: alle neun `p1-zusatz-*` |
| Kommentarzeile nur am Zeilenanfang | rot: `p1-kommentarzeile-eingerueckt` |
| Tag-Kommentar `# .*` statt `# [^ ].*` | rot: `p1-pin-tag-leer` |
| SHA mit `{40,}` (41 Ziffern erlaubt) | **grün** |
| SHA mit `[0-9a-z]` (Nicht-Hex erlaubt) | **grün** |
| Leerzeichen um `#` optional (`#v7.0.1`) | **grün** |

Die drei grünen Mutanten sind V-1.

## Pins am gepinnten d-check (Risiko 4)

`ci.yml` mit gekürztem `.d-check.yml` (`modules: [workflows]`, `workflows.dir`) im Image mit Digest `3e0b9779…`, jeweils eine Änderung an der Pin-Zeile: unverändert 0 Befunde; 41 Ziffern und Großbuchstaben `uses-pin-missing` (kein voller 40-stelliger Commit-SHA), ohne Tag-Kommentar `uses-pin-untagged`, `@main` `uses-pin-missing`. Die grünen Mutanten von V-1 (41 Ziffern, Großschreibung ist gefangen) färbt also `make docs-check` in der Kette, nicht die Gegenprobe; V-1 bleibt LOW.

## Plan-vs-Code (`AGENTS.md` §3.9)

§3 stimmt mit dem Diff: sechs Dateien plus `SPEC-064`; 148 Fälle, 144 Mutationen, Aufteilung wie im Skript. §7 Fälle und Mutationen stimmen mit dem Skript überein (Stichproben: `m-p7-*`, `m-p18-*`, `p8-rev-list-fehler`). §6 Randformen verweisen auf Punkte, die es gibt. Bekannt und dem Planner in der Closure zugewiesen, kein Befund: §1 `ubuntu-latest`, DoD-Punkt `concurrency`, Einleitung §6 „offen“. Nicht nachgezogen: V-2.

## Metrik (`AGENTS.md` §3.14)

`git diff --numstat 39f1ce8a..fc63d5ad -- .github tools harness .d-check.yml Makefile`: 765 Einfügungen + 10 Löschungen = **775** (Implementer/Reviewer: 775, bestätigt); `ci-gegenprobe.sh` 624, `ci-range.sh` 75, `ci.yml` 39, `ci.mk` 11, `SPEC-064` 18, `.d-check.yml` 7, `harness/README.md` 1. Unter der Marke 900, über der Meldeschwelle 700: keine Rückführung. §7 nennt die Messung bei `7d80d9f2` (771); der Endstand 775 steht dort nicht (V-2).

## Findings

| # | Schwere | Befund |
|---|---|---|
| V-1 | LOW | Zusage „`<SHA>` sind 40 Hexziffern in Kleinschrift“ (`SPEC-064` Punkt 1) ist je Element nicht vollständig gefangen: Fälle gibt es für 39 Ziffern, Großbuchstaben, fehlenden und leeren Tag, nicht für 41 Ziffern, ein Nicht-Hex-Zeichen und `#` ohne Leerzeichen. Drei Mutanten von `PIN_RE` bleiben grün (gemessen). Gemildert: `make docs-check` (Modul `workflows`) lehnt die Formen 41 Ziffern und Großschrift ab (gemessen), die Prüfung der Pins gehört nach Punkt 1 `SPEC-063`. Nach `AGENTS.md` §3.10 (Aufzählung zählt je Element) ein Fall und eine Mutation je Form oder eine Eingrenzung der Zusage auf `SPEC-063`. Beheben ist klein (rund 3 Fälle, 3 Mutationen); nicht blockierend. |
| V-2 | LOW | §7 nennt die Metrik des Standes `7d80d9f2` (771, `SPEC-064` 14); der Endstand `fc63d5ad` hat 775 (`SPEC-064` 18 durch `eaa9f0a3`). Gleichzeitig nennt `SPEC-064` an zwei Stellen einen Fall `p1-zusatz-schritt`, den das Skript nicht führt (`p1-zusatz-schritt-cache`, `p1-zusatz-schritt-artefakt`). Plan und Vertrag folgen dem Diff nicht ganz (§3.9). |
| V-3 | INFO | `export LC_ALL=C` in `ci-range.sh` hat weder Fall noch Mutation. Das Skript sagt es nicht zu (nur Plattform); kein Befund, nur dokumentiert. |
| V-4 | INFO | DoD Punkt 1 sagt „jede Entscheidung je Randform aus §6 steht vor dem ersten Code-Commit in `SPEC-064`“. Randform 19 wurde nach `4190a5c3` entschieden (Rückgabe des Implementers, `0bd5d47e`), der vorgesehene Weg. Beim Abhaken die Formulierung lesen: „vor dem ersten Code-Commit, außer den Rückgaben des Implementers (Randform 19)“. |
| V-5 | MEDIUM | **Neuer Zweig auf älterem Stand des Standard-Zweigs gibt ein falsches `core-drift-vcs`.** `ci-range.sh` nimmt bei einem neuen Zweig und bei einem Vorgänger, der nicht im Klon liegt, als Basis die Spitze von `refs/remotes/origin/<Standard-Zweig>` (`SPEC-064` Punkt 7, entschieden), nicht den Merge-Base. d-check vergleicht `--range Basis..Kopf` als Baumunterschied. Gemessen im Image mit dem Digest des Repos: `main` hat nach dem Abzweig eine ADR auf `Accepted` gesetzt und eine zweite angelegt, der Zweig davon einen Commit mit Kennung; `make doc-immutable` mit `RANGE=<Spitze main>..<Zweig>` ergibt 2 Befunde (`unzulässiger Status-Übergang einer immutablen Datei`, `immutable Datei gelöscht`), mit `RANGE=<Merge-Base>..<Zweig>` 0. `doc-commits` ist nicht betroffen (Mengenunterschied). Betrifft jeden ersten Push eines Zweigs, der nicht auf dem aktuellen Stand von `main` liegt, in einem Repo mit Worktrees (dieses Repo) wahrscheinlich; Pull-Requests nicht (der Kopf ist der Merge-Commit). Entscheidung des Architect nötig: Basis = `git merge-base` (Vertrag, Skript, Fälle `p7-push-neu`/`p7-push-vorgaenger-fehlt`, Mutation), oder als Grenze mit Grund festhalten. Der erste Probe-Push (Zweig vom aktuellen `main`) zeigt es nicht. |

Kein HIGH. V-5 ist die einzige Abweichung, die auf GitHub rot färbte, ohne dass die CI falsch läuft; sie ist eine Entscheidung im Vertrag, kein Fehler des Implementers gegen den Vertrag.

## Risiken aus §6: Vorschlag für den Ausgang (Notiz für den Planner)

Bis zum Probe-Push bleiben die Punkte 1 bis 3 ausdrücklich **offen**; sie sind kein Befund des Codes, sondern nur auf GitHub messbar.

| Risiko | Vorschlag | Stand jetzt |
|---|---|---|
| 1 `make gates` auf dem Runner | entfallen, wenn der Probe-Lauf grün endet; sonst eingetreten (Rückführung nach `open`, Beleg an `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert`) | **offen bis Probe.** Nur ein Lauf zeigt: Benutzer `runner` (nicht root; die Sonde des Implementers lief als root und als Benutzer) mit Bind-Mount von `:ro` für d-check und `actionlint`; anonymer Pull von `ghcr.io/pt9912/d-check` (Sichtbarkeit des Pakets); `setsid`, Docker-Netz der Integrationstests, Speicher; `doc-commits`/`doc-immutable` auf dem Runner; erste Range alter Commits (`SPEC-064` Punkt 15 nennt sie rot ohne Fehler der CI) |
| 2 Zeitgrenze 30 Minuten | entfallen, wenn die Dauer unter 15 Minuten liegt; sonst weiter offen mit Vermerk an `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, Grenze per Vertrag | **offen bis Probe.** 261 s im Dev-Container; Runner mit 2 vCPU und kalten Layern ungemessen (`make gates` zieht mehrere Postgres-Images) |
| 3 Rate-Limit Docker Hub | entfallen mit der Probe-Messung ohne Fehler, sonst weiter offen → `BEO-REPO/ci-pull-rate-limit` | **offen bis Probe** |
| 4 Branch-Schutz | weiter offen → neue Beobachtung `BEO-REPO/ci-ohne-branch-schutz` bei der Closure (existiert noch nicht im Register) | Wortlaut in `harness/README.md` vorhanden |
| 5 Pin-Prüfung am gepinnten d-check | **entfallen** für diesen Slice: `make docs-check` mit dem Eintrag an `ci.yml` meldet 0 Befunde, Pin ohne SHA, mit Großschrift, mit 41 Ziffern, auf `@main` und ohne Tag-Kommentar werden im gepinnten Image rot (in diesem Lauf gemessen, Digest `3e0b9779…`; Fälle `p1-pin-*` der Gegenprobe zusätzlich); der Register-Eintrag `workflow-pin-pruefung-am-gepinnten-d-check` bleibt für Verhalten bei `.yaml` und Blockform offen (Planner) | belegt |
| 6 Doppelte Läufe | **entfallen** mit Begründung: hingenommen (`SPEC-064` Punkt 2, [ADR-0050](../plan/adr/0050-ci-auf-github-actions.md)); `p1-zusatz-concurrency`, `p1-zusatz-workflow-dispatch` zeigen Auslöser und Verzicht | belegt |
| 7 Größe | **entfallen:** 775 unter 900 | belegt |
| neu (V-5) | eingetreten → Folge-Slice oder Änderung von `SPEC-064` Punkt 7 durch den Architect (Basis Merge-Base); sonst weiter offen → Register | Entscheidung offen |

## Verdikt

**Closure-fähig: nein, noch nicht.** Die Code-Seite (`ci.yml`, `ci-range.sh`, Gegenprobe, Gate-Eintrag, Ketten) erfüllt DoD und `SPEC-064` Punkte 1 bis 9, 14, 17, 18 nach eigener Messung. Offen bleiben: (a) Closure-Trigger §5, der grüne Lauf und der rote Lauf mit Verstoß gegen die Range auf GitHub samt Dauer, (b) die Entscheidung des Architect zu V-5 (MEDIUM, vor dem ersten Push eines abgezweigten Zweigs auf GitHub zu entscheiden; sie ändert Vertrag, Skript und Fälle), (c) V-1 und V-2 (LOW, mit der Nacharbeit zu erledigen oder in der Closure zuzuweisen), (d) Closure-Notiz, Lerneintrag, Register, Ausgänge.

**Volles `make gates` für die Merge-Entscheidung:** nicht nötig für das Urteil über den Slice. Seit `a31f4a7f` änderten sich nur die Gegenprobe (in keiner Kette, `GEGENPROBEN_CONTAINER`), Kommentarzeilen in `ci.mk`, eine Zeile in `harness/README.md`, `SPEC-064` und der Plan; die Gates, die diese Dateien lesen (`docs-check`, `kopf-check`, `abdeckung-check`, `workflows-check`, `baseline-verify`, `meldungskatalog-check`, Namensregel in `gates-schnell-gegenprobe`), habe ich in diesem Lauf einzeln gefahren: grün. Der Merge fährt `make gates` ohnehin in der Integrations-Worktree (Sache des Orchestrators); wird V-5 geändert, ist `ci-gegenprobe` voll neu zu fahren.

## Runde 2

Stand `d3818443` (Zweig `impl-ci`; Änderung seit `fc63d5ad`: `SPEC-064` Punkte 1, 7, 9, `ci-range.sh`, `ci-gegenprobe.sh`, Plan §3/§6/§7). Alles unten selbst gemessen, nicht aus §7 übernommen.

### V-5 (Merge-Base) — behoben

Temp-Repo (`main` lief nach dem Abzweig weiter, `refs/remotes/origin/main` auf der Spitze), `ci-range.sh` mit Push und Nullen als Vorgänger:

| Fall | Ergebnis |
|---|---|
| Zweig abgezweigt von `c1`, `main` bei `c3` | Exit 0, Basis = `c1` (Merge-Base), nicht die Spitze |
| Vorgänger nicht im Klon (`1111…`) | Hinweis auf stderr, Exit 0, gleiche Basis `c1` |
| Zweig ohne gemeinsamen Vorfahr (`--orphan`) | Exit 2, Meldung „keine Merge-Base …“ |
| flacher Klon | Exit 2 `fetch-depth: 0` |
| Kopf = Stand des Standard-Zweigs | Exit 1, leere Range |
| Criss-Cross (zwei Merge-Basen, `merge-base --all` liefert 2) | Exit 0, zweimal dieselbe Basis (git wählt deterministisch); **kein Fall** in der Gegenprobe |

Mehrere Merge-Basen: `SPEC-064` Punkt 7 sagt „die Wahl von `git merge-base` ist maßgeblich“. Das ist keine eigene Zusage des Skripts, sondern ein Verweis; nach §3.10/§3.12 vertretbar (Randform genannt und entschieden, Zusage nicht über git hinaus). Befund V-6 (INFO): kein Fall, kein Bedarf, solange der Vertrag so formuliert bleibt.

### Mutanten (`GP_ENTDECKE=1`, Teilmenge, danach grün)

| Mutant | Rot |
|---|---|
| `m-p7-spitze-statt-merge-base` | `p7-leer`, `p7-push-neu`, `p7-push-neu-ohne-vorfahr`, `p7-push-vorgaenger-fehlt`, `p8-stdout-eine-zeile` (genau die benannten) |
| `m-p7-merge-base-fehler-ohne-exit` | `p7-push-neu-ohne-standard`, `p7-push-neu-ohne-vorfahr` |
| `m-p7-namen-statt-shas` | `p7-pr-tag-objekt` (nur dieser) |
| `m-p7-pr-nimmt-vorgaenger` | wie benannt, jetzt mit `p7-pr-tag-objekt` |

Eigener Mutant außerhalb der Liste: fehlender Vorfahr fällt auf die Spitze von `origin/<Standard>` zurück (`merge-base … || rev-parse origin/main`). Im Temp-Repo liefert er mit dem Orphan-Zweig Exit 0 mit Range; `p7-push-neu-ohne-vorfahr` erwartet Exit 2 und färbt ihn also rot (direkt am Skript gesehen, nicht über `mutr`).

`p7-pr-tag-objekt` prüft die Zusage „aufgelöste Commit-SHAs“ (Punkt 7): ein Pull-Request-Kopf, der ein annotiertes Tag-Objekt ist, wird zum Commit geschält (`^{commit}`). Einzige Mutation: `aufloesen` gibt die Eingabe zurück (`m-p7-namen-statt-shas`); die alte Zuordnung hing an Fällen, in denen die Eingabe ohnehin ein Commit-SHA war, der Mutant war dort faktisch äquivalent, jetzt ist er gefangen. Gut.

### Zählung

Skript: 99 `p1` (+ `mut1`-Schleife), 39 `mutr`, 7 `mut18`; Fälle p1 99, p6 26, p7 12, p8 4, p9 6, p18 3 = **150**; Mutationen 99 + 39 + 7 = **145**; Läufe 145 + unveränderte Kopie + `rekursion-schutz` + `fall-ohne-funktion` = **148**. Das stimmt mit §3 (150/145, `p7` 12, `ci-range.sh` 39), §7 (148 Läufe, 150 Fälle) und dem Vertrag (Fallnamen von Punkt 7 und 9). Volle `make ci-gegenprobe` am Endstand: `grün (148 Läufe)`, 79 s.

Metrik: `git diff --numstat 39f1ce8a..d3818443 -- .github tools harness .d-check.yml Makefile` = 785 + 17 = **802** (Gegenprobe 638, `ci-range.sh` 75, `ci.yml` 39, `ci.mk` 11, `SPEC-064` 31, `.d-check.yml` 7, README 1); bestätigt, unter der Marke 900.

### V-1, V-2

- V-1: als Grenze entschieden (`SPEC-064` Punkt 1, Satz „Grenze der Fälle zum Pin“, Grund: das Modul `workflows` hält die Form über `make docs-check`). Der Satz ist in meiner Runde-1-Messung belegt (41 Ziffern, Großschrift, `@main`, fehlender Tag werden im gepinnten Image rot). Offene Form: „Nicht-Hex-Zeichen“ und `#v7` ohne Leerzeichen habe ich in Runde 1 nicht gegen d-check gefahren (nur die vier genannten); die Grenze nennt sie als von `docs-check` gefangen. **V-1 geschlossen unter dem Vorbehalt**, dass die Aussage für diese beiden Formen ungemessen ist (Vertragssatz reicht weiter als der Beleg, `AGENTS.md` §3.11; LOW, der Planner kann sie in der Closure mitnehmen oder der Implementer fährt zwei Fälle gegen das Image).
- V-2: §7-Metrik **behoben** (802 bei `c55bf42f`, gleich dem Endstand). `SPEC-064` nennt Punkt 1 weiter `p1-zusatz-schritt` (Zeilen 104 und 134), das Skript führt `p1-zusatz-schritt-cache` und `-artefakt`; **offen**, der Implementer hat es richtig an den Architect zurückgegeben (Vertragstext, Rückgabe-Regel). Außerdem steht in §7 im Absatz „Mutationen (145)“ noch „(147 Läufe: 144 Mutationen …)“ aus dem Stand davor (Zeile 187); der Absatz zu den Läufen sagt es („gelten für den Stand davor“), die Zahl im Absatz selbst ist veraltet (LOW, V-7).

### Sensoren dieser Runde

`bash -n` beider Skripte, `make docs-check` 0 Befunde, `make kopf-check` Exit 0, `make gates-schnell-gegenprobe` grün (65 s), `make ci-gegenprobe` `grün (148 Läufe)`. Kein `make gates`: geändert sind `ci-range.sh`, die Gegenprobe (in keiner Kette), `SPEC-064`, Plan.

### Echte GitHub-Läufe (nur gelesen)

| Lauf | Stand | Befund |
|---|---|---|
| `39189651083` (`probe/ci-2026-10-19`, `fc63d5ad`) | `success` | `ci-gates` 10 min 35 s (Job; 10m38s im Lauf); Schritte Set up, Checkout, Commit-Messages, ADR-Immutabilität, Gates, Post alle grün |
| `39190523004` (`probe/ci-rot-2026-10-19`, `16c2f563`, Commit ohne Kennung) | `failure` nach 13 s | Schritt „Commit-Messages der Range“ rot: `history-range-guard … 16 Commit(s) — OK`, dann d-check `commit-untraceable 16c2f56 Probe: Commit ohne Kennung`, `make: doc-commits Error 1`; ADR-Immutabilität grün, **Gates übersprungen** (`-`). Das ist der rote Fall des Closure-Triggers §5, wie vorgesehen |
| `39190523174` (`probe/ci-2026-10-19`, `d3818443`) | bei Abschluss dieser Runde `in_progress`, Schritte Commit-Messages und ADR-Immutabilität (Merge-Base-Pfad: Push auf bestehenden Zweig, Vorgänger bekannt) grün, `Gates` läuft (rund 10 min zu erwarten) | nicht abgewartet; **offen** |

Risiko 1 bis 3 (§6) bleiben bis zum Ende von `39190523174` offen; der erste grüne Lauf zeigt bereits: `make gates` läuft auf dem Runner als `runner`, anonymer Pull von `ghcr.io/pt9912/d-check` gelingt (sichtbar im roten Lauf), Dauer 10m35s liegt unter 15 min (Risiko 2 entfällt nach Vorschlag), kein Docker-Hub-Fehler im grünen Lauf. Der neue Lauf prüft zusätzlich den Pfad mit bekanntem Vorgänger auf dem neuen Skript; der Pfad „neuer Zweig“ mit Merge-Base ist auf GitHub noch nicht gefahren worden (der rote Lauf ist ein neuer Zweig, aber seine Range `afd4518d..16c2f563` hat nicht den V-5-Fall: er zweigt vom aktuellen `main` ab).

### Findings Runde 2

| # | Schwere | Befund |
|---|---|---|
| V-5 | geschlossen | Merge-Base gemessen, Orphan Exit 2, flach Exit 2, gleicher Stand Exit 1; Mutanten rot wie genannt |
| V-1 | geschlossen (LOW-Vorbehalt) | Grenze im Vertrag; Nicht-Hex und `#` ohne Leerzeichen gegen `docs-check` ungemessen |
| V-2 | teilweise | §7-Metrik behoben; `p1-zusatz-schritt` in `SPEC-064` offen (Architect) |
| V-6 | INFO | mehrere Merge-Basen: kein Fall; Verweis auf `git merge-base` im Vertrag genügt |
| V-7 | LOW | §7 Absatz „Mutationen (145)“ nennt noch „147 Läufe: 144 Mutationen“ |

### Verdikt Runde 2

**Closure-fähig: noch nicht, aber nur noch durch Außenstehendes:** (a) Lauf `39190523174` muss grün enden (Beleg des Closure-Triggers §5 zusammen mit dem roten Lauf `39190523004`, den es schon gibt); (b) `SPEC-064` Punkt 1 `p1-zusatz-schritt` → die Namen des Skripts (Architect); (c) V-7 und die Vorbehalte V-1 in der Closure-Notiz bzw. durch den Implementer; (d) Closure-Notiz, Lerneintrag, Register, Ausgänge der Risiken (Vorschläge in Runde 1; für Risiko 1 bis 3 nach grünem Lauf: entfallen). Kein HIGH/MEDIUM offen.
