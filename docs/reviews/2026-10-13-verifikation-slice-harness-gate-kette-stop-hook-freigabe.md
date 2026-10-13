# Verifikation: slice-harness-gate-kette-stop-hook-freigabe — 2026-10-13

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-13-review-slice-harness-gate-kette-stop-hook-freigabe.md`, Commit `785ed69`).

**Gegenstand:** Slice-Plan `slice-harness-gate-kette-stop-hook-freigabe` (Kopf, §1 bis §8, mit den Entscheidungen des Architect vom 2026-10-12/13/14: `23414ff`, `2f49e7e`, `d788954`, `5822f9f`, `1b7c390`, `278140c`) gegen `8349a65..278140c`. Commits: `6c6182b`, `a1fca6c`, `34a268f`, `b512b23` (Code), `2a95992`/`72111be` (Belege), `785ed69` (Review), `1b7c390` (Architect), `b6c2fbe` (F-713), `9e727fb` (F-714, F-716, F-720), `2f34a98`/`5fed425` (Belege), `278140c` (Architect: `reason`, Punkt 9).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-13

**Eingang:** Plan ganz, `SPEC-052` Teil 1 und Teil 2 (Punkte 1 bis 22, `spec/spezifikation.md` §11), [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) (Accepted), Review-Report (0 HIGH, 1 MEDIUM F-713, 3 LOW F-714 bis F-716, 7 INFO F-717 bis F-723), Diff `8349a65..HEAD -- tools harness .claude spec`. Den Bericht des Implementers habe ich nicht gelesen; ich lese seine Belege in §7 des Plans.

## 1. Selbst gefahrene Sensoren

Ergebnis vorweg: kein HIGH, kein MEDIUM, ein LOW (V-340), sonst INFO. Der Hook gibt in keiner meiner Proben öfter frei als `SPEC-052` es sagt; die Gegenprobe wird unter zwanzig eigenen Mutationen auf Freigabe-Lücken rot, an den Fällen, die der Plan nennt; F-713 ist geheilt. Der echte Arbeitsbaum blieb unberührt (`git status --short` leer; alle Proben liefen in Kopien im Scratchpad).

| Nr. | Sensor / Probe | Ergebnis |
|---|---|---|
| 1 | `make stop-hook-gegenprobe` | Exit 0, 2 min 4,2 s: „138 Fälle der Punkte 1, 9 bis 22 je Mutant rot, unveränderte Kopie grün, Anschluss an GATE_CHECKS“ (§7: 2 min 9,4 s). Die Fälle je Punkt in der Tabelle `FAELLE`: p1 3, p9 14, p10 10, p11 9, p12 23, p13 1, p14 13, p15 14, p16 3, p17 1, p18 2, p19 2, p20 29, p21 10, p22 4, Summe 138 |
| 2 | `make gates-schnell-gegenprobe` | Exit 0, 20,5 s |
| 3 | `make gates-schnell` | Exit 0, 7,9 s (Cache warm); `.harness/state/gates-schnell-passed.diffsha` gleich der Ausgabe von `tools/harness/working-tree-hash.sh` (`d1d1e2a3…e79`) |
| 4 | `make hook-gegenprobe` | Exit 0, 0,18 s (§7: 0,18 bis 0,19 s) |
| 5 | `make docs-check`, `make kopf-check` | 694 Dateien, 0 Befunde; `kopf-check` Exit 0 |
| 6 | Größe nach der Metrik (`git diff --numstat 8349a65..9e727fb -- tools harness .claude`, Einfügungen plus Löschungen) | 670 + 60 = **730** (stimmt mit §7); bis `HEAD` dieselbe Summe, weil `b6c2fbe`/`9e727fb` die letzten Commits in diesen Pfaden sind. Schätzung 550 bis 750, Marke 900. Anteil der Nacharbeit (`b512b23..9e727fb`): 71 + 42 = 113 (§7: 113). Aufteilung: Hook 206/9, Gegenprobe 419/9, `commit-msg-gegenprobe.sh` 23/28, Commands je 5/3, `record-gates.sh` 4/2, Rest je 1/1 |
| 7 | Hook im echten Arbeitsbaum | `block` („A changed path is not on the inert list“): erwartet, weil `278140c` `spec/spezifikation.md` änderte, nach dem Lauf von `make gates` am Stand `5fed425` (siehe V-341) |

### Hook von Hand (Kopie des Arbeitsbaums im Scratchpad, `git clone` von `278140c`)

Je Zeile: Kopie zurückgesetzt, Nachweis des vollen Laufs mit `tools/harness/record-gates.sh` im Normalaufruf geschrieben, Änderung angewandt, Stand der schnellen Kette (falls „mit“) als `working-tree-hash.sh`-Ausgabe in `.harness/state/gates-schnell-passed.diffsha` gelegt, `bash .claude/hooks/stop-require-gates.sh </dev/null`. Der Stand der schnellen Kette ist hier ein Stellvertreter für den Lauf von `make gates-schnell`; den Schreiber prüft `make gates-schnell-gegenprobe` (Probe 2).

| Änderung | mit Stand der schnellen Kette | ohne |
|---|---|---|
| sauberer Baum nach `record-gates` | `approve` | `approve` (Identität) |
| `docs/plan/planning/…`, `docs/reviews/…`, `docs/user/benutzerhandbuch.md`, `docs/user/abdeckung-e2e.md`, `README.md` je angehängt | `approve` (alle fünf) | `block`, `reason` nennt `make gates-schnell` (alle fünf) |
| `harness/README.md`, `docs/users/x.md`, `README.md.bak`, `docs/plan/planning-x/x.md`, `readme.md`, `docs/reviews-alt/x.md`, `docs/Plan/planning/x.md` (neu) | `block`, „not on the inert list“ (alle sieben) | — |
| `git mv` von `docs/user/` nach `internal/` / von `go.mod` nach `docs/user/` | `block` / `block` | — |
| `git mv` innerhalb `docs/user/` | `approve` | — |
| untracked Datei in `docs/user/` / in `internal/` / in `internal/` neben geändertem `README.md` | `approve` / `block` / `block` | — |
| gelöschte Datei in `docs/user/` / `go.sum` (Arbeitsbaum) / `go.sum` (`git rm`) | `approve` / `block` / `block` | — |
| Link in `docs/user/` mit Ziel `../reviews/…md` (inert) | `approve` | — |
| Link in `docs/user/` mit Ziel `/etc/passwd` · `../../go.mod` · `../reviews` (Verzeichnis der Liste) · `../../../../x` (über die Wurzel) | `block` („changed symlink is not inert“), alle vier | — |
| Art-Wechsel Datei zu Link (Ziel `../../go.mod`) | `block` | — |
| Link außerhalb der Liste (`internal/l`, im Stand vorhanden) mit Ziel `../docs/user/benutzerhandbuch.md` · `../docs/user` · `../docs/reviews/` · `../docs` · `../docs/plan` · `..` · absolut unter der Wurzel · `../docs/user/gibtsnicht.md`, danach Änderung in `README.md` bzw. `docs/plan/planning/…` | `block` („A symlink outside the inert paths points into them“), je beide Änderungen | — |
| dieselben Links mit Ziel `/etc/passwd` (absolut außerhalb) · `../../../x` (über die Wurzel) · `../go.mod` (außerhalb der Liste) | `approve` (Gegenrichtungen: Ärgernis, keine Lücke, so entschieden) | — |
| `# x` angehängt an: Hook, `dateiliste.sh`, `working-tree-hash.sh`, `spec/spezifikation.md`, `.dockerignore`, `.d-check.yml`, `Makefile`, `Dockerfile`, `tools/test/run-integration-tests.sh`, `harness/mk/gates-schnell.mk`, `go.mod` | `block`, „not on the inert list“ (alle elf) | — |

Weitere Grenze, mit Beleg: `internal/domain/doc.go` gibt es nicht; die Aufzählung des Auftrags („`internal/…`“) habe ich mit `internal/neu.go` (untracked) und der gelöschten `go.sum` gefahren, beide `block`.

### F-713 von Hand (Kopie, Nachweis des vollen Laufs auf `278140c`)

Alle neun Slices aus `open/` per `make slice-mv … TO=next` nach `next/` bewegt (acht über das Ziel; `slice-tests-ueberlebende-mutanten` meldete „mehrdeutig“, weil `slice-tests-ueberlebende-mutanten-driving` dasselbe Präfix trägt, der neunte per `git mv`):

| Zustand | `make hook-gegenprobe` |
|---|---|
| vor den Zügen | Exit 0 |
| nach den Zügen, Skript der Lieferung (`b6c2fbe`) | Exit 0 |
| nach den Zügen, Skript aus `1b7c390` (Mutante der Tabelle) | Exit ungleich 0 (make meldet 2; das Skript endet mit 1) |
| nach den Zügen: `make gates-schnell` | Exit 0; Hook `approve` |

Die schnelle Kette gibt damit nicht für eine Menge frei, für die der volle Lauf rot wäre. Das Skript läuft zudem in einem Verzeichnis ganz ohne `docs/plan/planning/` grün (Exit 0).

Zusagen des Commit-Trägers (`SPEC-050`) unter Mutation des Trägers in diesem Verzeichnis (Skript der Lieferung): Slice-Muster nur `done/` → rot an `slice-gegenprobe-offen`; Zeile `*-results) continue` gestrichen → rot an `welle-gegenprobe-fertig-results`; Kommentarfilter entfernt → rot an „ohne Kennung“/Kommentar; Muster der Welle auf `docs/plan/planning/*` → rot. Die Mutation ist je rot am Fall der Zusage, die sie bricht.

### Suche nach einem Pfad der Positivliste, den ein langsames Gate liest (`SPEC-052` Punkt 13)

Gelesen und mit `grep` nach `docs/`, `README`, `planning`, `reviews`, `*.md`, `find .`, `git ls-files` durchsucht: `tools/test/run-integration-tests.sh`, `Dockerfile`, `.dockerignore`, `tools/harness/lint-gegenprobe.sh`, `tools/arch/a-check-negativ.sh`, `tools/harness/kopf-check-gegenprobe.sh`, `tools/test/abdeckung-gegenprobe.sh`, `tools/harness/gates-schnell-gegenprobe.sh`, `tools/harness/stop-hook-gegenprobe.sh`, `tools/hook/commit-msg-gegenprobe.sh`, `test/integration/*`, `internal/**` (Go-Dateien), alle `harness/mk/*.mk`. Treffer außerhalb der Gates der schnellen Kette: keiner, der eine Datei der Positivliste aus dem Arbeitsbaum liest.

- Die Build-Kontexte der Docker-Gates sind die Allowlist der `.dockerignore` (`go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`, `.golangci.yml`, `tools/harness/lint.sh`); `docs/`, `README.md` gehören nicht dazu. `test/integration/record_e2e_test.go` nennt `docs/user/abdeckung-e2e.md` nur im Kommentar.
- `abdeckung-gegenprobe` und `kopf-check-gegenprobe` bauen ihre Bäume mit `mktemp -d`; `gates-schnell-gegenprobe` und `stop-hook-gegenprobe` legen eigene Fälle an (Pfade wie `docs/user/u.md` sind erfundene Dateien); `a-check-negativ` kopiert den Baum, `a-check` liest nach `.a-check.yml` nur Go-Dateien; `lint-gegenprobe` baut über denselben Kontext. `tools/harness/e2e-abdeckung.sh` liest `docs/user`, ist aber kein Gate (kein Eintrag in `GATE_CHECKS`, nicht von `run-integration-tests.sh` gerufen).
- `hook-gegenprobe` war der einzige Leser (F-713); er ist nach `b6c2fbe` ein eigenes Temp-Repo (F-713 von Hand oben).

### Mutationen (eigene Kopien im Scratchpad, je Mutant eine frische Kopie von `278140c`, im Hook mutiert, danach die unveränderte `stop-hook-gegenprobe.sh` der Kopie gefahren)

Jeder Lauf endete mit Exit 1. Maßgeblich ist die erste Zeile `Kopie 'ok' ist nicht grün: …`; die weiteren Zeilen („Mutant … ändert nichts“, „färbt auch … rot“) kommen daher, dass die Gegenprobe ihre eigenen Mutanten nochmals auf dem schon mutierten Hook anwendet.

| Nr. | Zusage | Mutation | rote Fälle der Kopie `ok` |
|---|---|---|---|
| V1 | Grenze mit Schrägstrich | `docs/plan/planning*)` | `p12-grenze-planning-x` |
| V2 | nur die Wurzel-`README.md` | `README.md \| */README.md` | `p12-harness-readme` |
| V3 | case-sensitiv | `shopt -s nocasematch` | `p12-case-planning`, `p12-case-readme` |
| V4 | gelöschte und umbenannte Dateien zählen mit beiden Pfaden | zweite Schleife in `diff_bilden` (Pfade nur des Standes) entfernt | `p14-geloescht-index-nicht-inert`, `p14-umbenennung-nach-inert`, `p12-unzulaessig-*` (fünf) |
| V5 | untracked zählt nach seinem Pfad | Pfade, die nicht im Stand stehen, überspringen (pauschal inert) | `p18-untracked-nicht-inert`, `p14-neu-nicht-inert`, `p14-umbenennung-von-inert` und die Grenzfälle |
| V6 | Spiegelregel | Aufruf `spiegel` entfernt | alle fünf `p15-spiegel-*` |
| V7 | Stand gehört zum aktuellen Arbeitsbaum | nur auf Vorhandensein prüfen (`[ -n … ]`) | `p10-fremd` |
| V8 | `harness/README.md` nicht inert | in `inert` aufgenommen | `p12-harness-readme` |
| V9 | Pfad mit `..` nie inert | Prüfung auf `..` aus `zulaessig` entfernt | `p12-unzulaessig-punktpunkt` |
| V10 | nachgestelltes leeres Segment (Pfad auf `/`) | `\| */` aus `zulaessig` entfernt | `p12-unzulaessig-leer-ende` |
| V11 | scheiterndes Hash-Skript sperrt | `\|\| true` am Aufruf | `p10-skript-scheitert` |
| V12 | Pfadregel vor dem Stand der schnellen Kette (beobachtbar am `reason`) | `schnell_gruen` vor `pfadregel` | `p9-stand-fehlt-nicht-inert` |
| V13 | Art `X` nie inert | Zeile entfernt | `p12-art-x` |
| V14 | Symlink zählt mit Ziel | Zeile `link_inert` entfernt | acht Fälle von Punkt 15 (`p15-ziel-aussen`, `-ueber-wurzel`, `-absolut`, `-absolut-aussen`, `-ausserhalb-liste`, `-verzeichnis`, `p15-art-datei-zu-link`, `-link-zu-datei`) |
| V15 | Spiegelregel je Zieltyp (Vorfahr) | `docs \| docs/plan` entfernt | `p15-spiegel-vorfahr` (allein) |
| V16 | Grenze von `docs/user/` | `docs/user*)` | `p12-grenze-users`, `p12-grenze-datei-user`, `p15-ziel-verzeichnis` |
| V17 | leere Ausgabe des Hash-Skripts sperrt | Prüfung entfernt | `p10-skript-leer` |
| V18 | Konflikt steht vor der Identität | Prüfung `git ls-files -u` entfernt | `p17-konflikt-identisch`, `p9-konflikt`, `p20-git-u-scheitert`, `p20-schleifenschutz-konflikt`, `-git-u` |
| V19 | `spec/spezifikation.md` nicht inert | in `inert` aufgenommen | `p12-selbst-spec` |
| V20 | `.dockerignore` nicht inert (Bedingung (b) von Punkt 13) | in `inert` aufgenommen | `p9-nicht-inert-dockerignore` |

Kein Mutant blieb grün; V12 zeigt, dass die Ordnung „Pfadregel, zuletzt Stand“ (F-719) jetzt beobachtbar ist. V1, V2, V3, V4, V5, V6, V7, V8, V9, V10, V11 und V12 sind die zwölf Mutationen auf die im Auftrag genannten Lücken (Verzeichnisgrenze ohne Schrägstrich, `README.md` in jeder Tiefe, Groß-/Kleinschreibung, nur neuer Pfad, untracked pauschal inert, Spiegelregel, Stand nur auf Vorhandensein, `harness/README.md` inert, `..` im Stand, trailing slash, Fehlschlagen des Hash-Skripts, Reihenfolge Stand-vor-Pfadregel).

**Zur Frage des `reason` (V-340):** In einer weiteren Kopie (R1) habe ich `schnell_gruen` aufgeteilt: ein *fehlender* Stand nennt `make gates-schnell`, ein *ungleicher* nur noch `make gates`. Die Kopie `ok` blieb mit allen 138 Fällen grün; kein Fall prüft den `reason` bei fremdem oder leerem Stand.

## 2. DoD Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| Liefer-Punkt 1: Hook gibt nach `SPEC-052` frei, wenn alle geänderten Pfade auf der Positivliste liegen (Grenze mit Schrägstrich, bytegenau, case-sensitiv; gelöscht, umbenannt mit beiden Pfaden, untracked wie getrackt; `.`/`..`/leeres Segment/absolut als beschädigte Liste an den Pfaden der Differenz; Art `X` nie inert) **und** der Stand der schnellen Kette der Hash des Arbeitsbaums ist (aus `working-tree-hash.sh` gerufen; fehlend, ungleich, leer, nicht lesbar heißt kein Stand; fehlt das Skript oder scheitert es, `block`, beim Schleifenschutz `approve`); identische Liste ohne Stand der schnellen Kette frei; sonst `make gates`; Kopf-Kommentar sagt zu, soweit die Gegenprobe prüft | erfüllt, mit einer Einschränkung (V-340) | Hook von Hand (alle Zeilen oben), Probe 1, V1 bis V20. Positivliste nur in `inert()` (eine Stelle); Reihenfolge im Code: Schleifenschutz, `git ls-files -u`, Stand fehlt, `dateiliste_pruefen`, Identität (`cmp`), `pfadregel`, `schnell_gruen`; der Hook ruft kein `make` und baut den Hash nicht nach. Der Kopf des Hooks sagt den Wortlaut des `reason` weiter zu als Fall und Spezifikation (V-340). Die Bedingung der Freigabe „kein langsames Gate liest die Positivliste“ ist mit der Heilung von F-713 (`b6c2fbe`) erfüllt (Abschnitt 1) |
| Liefer-Punkt 2: Symlink inert, wenn Link-Pfad und lexikalisch aufgelöstes Ziel auf der Liste liegen; Ziel außerhalb, absolut, nicht auflösbar, Verzeichnis der Liste nicht inert; Spiegelregel (Ziel auf der Liste, Verzeichnis der Liste, Vorfahr, Wurzel, absolut unter der Wurzel) | erfüllt | Hook von Hand (Link-Zeilen, 8 Mirror-Ziele, 3 Gegenrichtungen), V6, V14, V15; `p15-*` 14 Fälle |
| Liefer-Punkt 3: Gegenprobe macht je Randform einen Fall in einem Temp-Repo, nur im inerten Fall freigebend; Zeile in `harness/README.md` §Sensors nennt Positivliste und Bindung an [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md); `harness/conventions.md` Träger-Zeile richtig | erfüllt | Probe 1 (138 Fälle, 70 neu gegenüber dem Geber, zwei weitere durch die Nacharbeit); `harness/README.md` Zeile 79: „Punkt 1, 9 bis 22“ (stimmt mit der Meldung der Gegenprobe), Positivliste mit allen vier Einträgen, `harness/README.md` ausdrücklich nicht, Symlinks mit Ziel und Spiegelregel, Bindung [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) und „seit slice-harness-gate-kette-stop-hook-liste“. `harness/conventions.md` Zeile `grundlagen-durchsetzungsschicht.md` nennt den Hook mit Positivliste und schneller Kette; `harness/mk/stop-hook-gegenprobe.mk` Kopf „Punkte 1, 9 bis 22“; `grep -rn "nur bei identischer\|gibt nur frei" tools harness .claude` ohne Treffer in lebenden Stellen (Treffer nur in `done/`-Plänen und in `SPEC-052` Punkt 10 als Satz über einen fremden Stand) |
| `make gates` grün | laut Implementer erfüllt am Stand `5fed425`; nicht von mir gefahren; am `HEAD` nicht belegt (V-341) | §7: Exit 0, 7 min 22,9 s, Hook danach `approve`; `278140c` änderte danach `spec/spezifikation.md` (nicht inert), der Hook im echten Arbeitsbaum antwortet `block` (Probe 7) |
| Review durchgeführt, Report unter `docs/reviews/` | erfüllt | Report `785ed69`: 0 HIGH, 1 MEDIUM, 3 LOW, 7 INFO; fail-closed ausdrücklich geprüft (23 Mutationen). Die Nacharbeit (`b6c2fbe`, `9e727fb`) hat keinen zweiten Report; ich habe sie in Abschnitt 1 nachgefahren (F-713 von Hand, `p10-fehlt`, `p9-stand-fehlt-nicht-inert`, `p12-unzulaessig-leer-ende`, V10, V12) |
| Mutationen als Tabelle in §7, vom Implementer selbst gefahren (§3.10, je Element einer Aufzählung) | erfüllt | §7 trägt drei Tabellen (Mutationen des Slice, Zusagen des Trägers, `reason`) mit Mutant, rotem und grünem Fall; jede Zeile der Tabelle in §6 hat eine Entsprechung (ich habe die Zeilen der Spalten Positivliste, Grenze, Symlink, Spiegel, Stand, Hash-Skript, Reihenfolge, Selbständerung, `reason` und F-713 gegengelesen). Meine zwanzig Mutationen bestätigen die roten Fälle, die ich nachgefahren habe, ohne Abweichung |
| Closure-Notiz mit Lerneintrag | offen (Closure) | §7 trägt nur *Belege des Implementers* und *Belege der Nacharbeit* |
| Beobachtungs-Register: Evidence `slice-harness-gate-kette-stop-hook-freigabe.md` in `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/` (4×, Ausgang *verkörpert* oder *offen*) | offen (Closure) | Datei nicht angelegt. Messwerte stehen in §7 (7 min 18 s bis 7 min 22,9 s voller Lauf; schnelle Kette 7,9 s bis 8 s, meine Probe 3: 7,9 s; Hook 4,5 s je Aufruf) |
| Jedes Risiko aus §6 trägt einen Ausgang | offen (Closure) | §6 trägt je Risiko einen vorgeschlagenen Ausgang; die Closure trägt sie ein |
| Drei Paarungen | offen (Closure) | |

## 3. Entscheidungs-Konformität

| Entscheidung | Urteil | Beleg |
|---|---|---|
| [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) Entscheidung 1: `make gates` bleibt der Nachweis für alles, was nicht auf der Positivliste steht | erfüllt | Hook von Hand (elf Pfade, alle `block`); `make gates` unverändert (`harness/mk/*.mk` und `Makefile` nicht im Diff) |
| [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) Entscheidung 2: schnelle Kette ohne Gegenproben und ohne Integrationstests | erfüllt | Ausnahmeliste unverändert (acht Einträge, `gates-schnell.mk` nicht im Diff); `gates-schnell-gegenprobe` grün (Probe 2) |
| [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) Entscheidung 3: Positivliste nur mit Beleg, dass kein langsames Gate den Pfad liest | erfüllt nach Heilung | Suche in Abschnitt 1; F-713 behoben (`b6c2fbe`), Beleg in `SPEC-052` Punkt 13 für `docs/plan/planning/` und `docs/reviews/` geführt |
| Nutzer 2026-10-13: Positivliste `docs/plan/planning/`, `docs/reviews/`, `docs/user/` (einschließlich Abdeckungstabellen), Wurzel-`README.md` | erfüllt | Hook von Hand, V2, V8, V16 |
| Architect 2026-10-14 (`5822f9f`): Reihenfolge, Hash-Aufruf, Zulässigkeit an der Differenz, Art `X`, Symlink-Einzelheiten, Selbständerung, Unterverzeichnis und leere Umgebung | erfüllt | Reihenfolge: V12, V18 und drei Mutanten der Gegenprobe; Selbständerung: elf Pfade von Hand, V19/V20; Unterverzeichnis/`env -i`: `p22-freigabe-*` in Probe 1 |
| Architect 2026-10-14 (b) (`1b7c390`): F-713 Entscheidung (c), F-716 Fall, F-720 `reason`, F-717 Satz in Punkt 21, F-718 Ersetzung der Mutanten, F-721 Laufzeit, F-722 `LC_ALL` | erfüllt | F-713 von Hand; F-716: `p12-unzulaessig-leer-ende` rot unter V10, `m-ohne-leeres-segment-ende` in der Gegenprobe; F-720: `reason` mit `make gates-schnell` nur bei fehlendem Stand und inerter Änderung (Hook von Hand), `p10-fehlt`/`p9-stand-fehlt-nicht-inert` rot unter V12; F-717: Punkt 21 letzter Satz nennt die Positivliste; F-718: Zeilen in §6 tragen die Ersetzung; F-721: Punkt 20 nennt keine Dauer, `settings.json` im Diff unberührt; F-722: `export LC_ALL=C` bleibt ohne Zusage im Kopf |
| Architect `278140c`: `reason` nennt `make gates-schnell` in `p10-fehlt`, nicht in `p9-stand-fehlt`; Punkt 9 auf „fehlender Stand“ gefasst | erfüllt | `SPEC-052` Punkt 9 sagt: fehlender Stand bei nur inerten Änderungen nennt `make gates-schnell`, bei nicht inerter Änderung nur `make gates`; Fälle `p10-fehlt` und `p9-stand-fehlt-nicht-inert` prüfen genau das; Wortlaut für fremd/leer/nicht lesbar „nicht zugesagt“ |
| Keine Lockerung (§3.6) | erfüllt | Der Hook sperrt nirgends weniger als nötig: jeder Fall außerhalb von „identisch“ und „jeder Unterschied inert und Hash der schnellen Kette gleich“ ergibt `block` (Hook von Hand, V1 bis V20). Die schnelle Kette gibt nicht öfter frei als der volle Lauf für dieselbe Menge: kein Gate der Ausnahmeliste liest die Positivliste (Abschnitt 1, F-713 von Hand). Die Freigabe über die Positivliste ist die von [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) entschiedene Wahl, keine Senkung einer Schwelle |
| Aufruf-Voraussetzung nacheinander (Punkt 8, §3.12) | erfüllt | Kopf des Hooks und `SPEC-052`; kein Fall (so vorgesehen) |
| Akzeptierte Negative (Hook entscheidet über die eigene Änderung; Seite des Arbeitsbaums bei unzulässigen Pfaden; Moduswechsel; Ursprung der Stand-Datei; ignorierte Dateien) | als akzeptiert geführt | `SPEC-052` Punkte 10, 11, 12, 18; ich habe keine Probe gefunden, die über sie hinausginge |

## 4. Plan gegen Diff (`AGENTS.md` §3.9) und Kommentare (§3.11)

- **§3 (Plan):** nennt Hook, Gegenprobe, `harness/README.md`, `harness/conventions.md`, Fragment-Kopf, die drei Commands, `record-gates.sh` (Kopf) und `tools/hook/commit-msg-gegenprobe.sh`; der Diff in `tools harness .claude` berührt genau diese elf Dateien (`git diff --stat`), dazu `spec/spezifikation.md` (Architect). F-715 ist behoben. Kein Produkt-Code, kein Handbuch, keine Tests unter `internal/` oder `test/`.
- **§1, §6, §7, §8:** folgen dem Diff; §6 nennt die Entscheidungen von (b) mit Datum, §7 die Belege der Nacharbeit mit „zurück“-Befehl, §8 die Zählung (3 Liefer-Punkte, 1 Schicht; die Heilung von F-713 steht unter Liefer-Punkt 1). Der Kopf nennt als geschärfte Punkte nur 9, 10, 12, 15, 22 (V-342).
- **Commands:** `implement-slice.md`, `plan-welle.md`, `close-welle.md` nennen Positivliste (`inert` als Quelle) und `make gates-schnell` und stimmen mit meinen Proben überein (identisch oder jeder Unterschied inert und schnelle Kette grün frei; fehlender Nachweis sperrt auch bei sauberem Baum).
- **Kommentare, die nur zusagen, was ein Fall prüft (§3.11):** Hook (Konflikt, fehlender/beschädigter Stand, Eintragsarten, Positivliste, Art `X`, Symlink, Stand der schnellen Kette bis auf V-340, Schleifenschutz, Läufe nacheinander), `stop-hook-gegenprobe.sh` (Kopf nennt Grenzen: Attrappen, root, Schleifenschutz, drei grüne Mutanten, Seite des Arbeitsbaums), `stop-hook-gegenprobe.mk`, `record-gates.sh` (F-714 behoben), `commit-msg-gegenprobe.sh` (Tabelle der fünfzehn Fälle, alle in der Gegenprobe vorhanden und unter Mutation rot) sagen nicht mehr zu, als ein Fall prüft. Ausnahme: V-340.
- **Commits:** reine Inhalts-Commits, kein `git mv`; Slice-Commits nennen Slice-Kennung und [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md).

## 5. Befunde

Kein HIGH, kein MEDIUM.

| ID | Schwere | Befund | Beleg |
|---|---|---|---|
| V-340 | LOW | Der Kopf des Hooks sagt weiter zu, als ein Fall prüft und als `SPEC-052` Punkt 9 festlegt. Er schreibt: „Fehlt der Stand oder gehört er zu einem anderen Arbeitsbaum, obwohl jede Änderung inert ist, nennt der `reason` `make gates-schnell`; jede andere Sperre nennt nur `make gates`.“ Die Spezifikation sagt den Wortlaut nur für den *fehlenden* Stand bei inerter Änderung und für die nicht inerte Änderung zu und nennt „Stand eines anderen Arbeitsbaums, leer, nicht lesbar“ ausdrücklich nicht zugesagt; die Gegenprobe ruft `grund`/`kein_grund` nur in `p10-fehlt` und `p9-stand-fehlt-nicht-inert`. Mutation R1 (der Hinweis nennt `make gates-schnell` nur bei fehlendem, nicht bei ungleichem Stand) lässt alle 138 Fälle grün. Kein Fehlverhalten (der Wortlaut ist ein Hinweis, die Sperre bleibt), aber `AGENTS.md` §3.11 verlangt, dass ein Kommentar nur zusagt, was ein Test prüft. Heilung wäre im Kopf: „Fehlt der Stand …“ auf den fehlenden Stand fassen und den Satz „jede andere Sperre nennt nur `make gates`“ auf die nicht inerte Änderung beschränken. Ein Satz im Kopf, kein Code, keine Spezifikation. | `.claude/hooks/stop-require-gates.sh` Zeilen 27 bis 29; `tools/harness/stop-hook-gegenprobe.sh` Zeilen 586, 610 (einzige Aufrufe von `grund`/`kein_grund`); `spec/spezifikation.md` Punkt 9; Mutation R1 |
| V-341 | LOW | Der DoD-Punkt „`make gates` grün“ ist am Stand `5fed425` belegt (§7: Exit 0, 7 min 22,9 s, Hook `approve`), nicht am letzten Stand. `278140c` änderte `spec/spezifikation.md` (Punkt 9), und diese Datei ist nach Punkt 12 nie inert. Der Hook im echten Arbeitsbaum antwortet deshalb `block` (Probe 7). Für die Closure ist ein voller Lauf nach dem letzten nicht inerten Commit nötig (Plan §5: „vom Auftraggeber gefahren“); diesen Lauf habe ich nicht gefahren, und er ist im Plan noch nicht eingetragen. Mein eigener Report liegt unter `docs/reviews/` (inert) und ändert das nicht. | Plan §7 Zeile „Stand der Gates der Nacharbeit“; `git diff 5fed425..278140c --stat`; Probe 7 |
| V-342 | INFO | Der Kopf des Plans nennt unter *Berührte Spec-Stellen* nur die Schärfung vom 2026-10-14 vor dem Code (Punkte 9, 10, 12, 15, 22). Die Entscheidungen nach dem Review (b) und von `278140c` änderten Punkte 9, 12, 13, 20, 21; sie stehen in §6, nicht im Kopf. `AGENTS.md` §3.9 verlangt, den Kopf nachzuziehen. Kein Sensor (`make kopf-check` prüft Kennungen aus §1 und §2, nicht diesen Satz). | Plan Kopf Zeilen 26 bis 32; §6 *Entscheidungen des Architect nach dem Review* |
| V-343 | INFO | `SPEC-052` Punkt 22 (b) liest sich einzeln als „leere Umgebung: bei identischer Liste `approve`, bei jedem Unterschied `block`“; der folgende Satz („beide Aufrufarten gelten auch für die Freigabe nach Punkt 9“) hebt das auf. Die Gegenprobe prüft die Freigabe in beiden Aufrufarten (`p22-freigabe-unterverzeichnis`, `-leere-umgebung`; Mutanten `m-22-hash-relativ`, `m-22-hash-home`). Der Wortlaut der Spezifikation ist beim Lesen von (b) allein missverständlich; eine Entscheidung des Architect, kein Code. | `spec/spezifikation.md` Punkt 22 |
| V-344 | INFO | Die Nacharbeit nach dem Review (b) (`b6c2fbe`, `9e727fb`: 113 Zeilen, Hook, Gegenprobe, Skript des Commit-Trägers) hat keinen zweiten Review-Report; die Rollenfolge verlangt ihn nicht. Ich habe die Änderungen gelesen und die vier Teile gefahren: F-713 (Abschnitt 1), `reason` (Hook von Hand, V12), F-716 (V10), F-714 (`grep`). Kein Befund daraus. | Abschnitt 1 |
| V-345 | INFO | Die Probe, ob ein Beispiel im Handbuch gegen das gebaute Binary läuft (`SPEC-052` Punkt 13, F-723), bleibt Rollenpflicht von Implementer und Verifier. Der Hook gibt eine Änderung in `docs/user/` nach der schnellen Kette frei, ohne sie. Dieser Slice ändert kein Handbuch; ich habe sie nicht gefahren, weil sie hier keinen Gegenstand hat. Für künftige Slices mit Handbuch-Änderung bleibt sie ein Punkt der Verifikation und nicht des Hooks. | `SPEC-052` Punkt 13 |

## 6. Zu den Prüfpunkten des Auftrags

- **Jeder DoD-Punkt einzeln:** Abschnitt 2.
- **Hook von Hand (Freigabe-Fälle):** Abschnitt 1 (alle Zeilen des Auftrags: sauber, Änderung in den vier Pfaden der Liste mit und ohne Stand, `harness/README.md`, `docs/users/`, `readme.md`, `README.md.bak`, `docs/plan/planning-x/`, Umbenennung in beide Richtungen, untracked, gelöscht, Symlink mit Ziel außerhalb und innerhalb samt Spiegelregel, Änderung an Hook, `dateiliste.sh`, `working-tree-hash.sh`, `spec/spezifikation.md`, `.dockerignore`, `.d-check.yml`, `Makefile`, `Dockerfile`, `tools/…`, `harness/mk/…`, `internal/…`).
- **F-713 selbst:** behoben und von Hand bestätigt (Abschnitt 1).
- **Pfad der Positivliste, den ein langsames Gate liest:** keiner gefunden (Abschnitt 1).
- **Mutationen:** §7 nennt je Zusage und je Element eine Mutation mit rotem Fall; ich habe zwanzig (V1 bis V20) plus R1 gefahren, alle rot an den benannten Fällen, R1 bewusst grün (V-340).
- **Folgt `SPEC-052` dem Code (Punkte 9 bis 22, `reason` nach `278140c`):** ja. Punkt 9 (Reihenfolge, `reason` nur bei fehlendem Stand und inerter Änderung), 10 (Hash, Skript fehlt/scheitert), 12 (Positivliste, Grenze, unzulässige Pfade, Art `X`, Selbständerung), 13 (Beleg, jetzt auch für `docs/plan/planning/` und `docs/reviews/`), 14, 15 (Symlink, Spiegelregel), 16 bis 22 stimmen mit meinen Proben überein. V-343 betrifft nur einen Wortlaut.
- **F-713 bis F-723:** F-713 behoben (`b6c2fbe`, von Hand); F-714 behoben (`record-gates.sh` Kopf, `grep`); F-715 behoben (Plan §3); F-716 behoben (Fall und Mutation, V10); F-717 behoben (Punkt 21); F-718 bestätigt; F-719 gelöst durch `reason`-Fälle (V12); F-720 behoben; F-721 und F-722 als INFO entschieden; F-723 INFO (V-345).
- **Plan folgt dem Diff (§3.9):** ja, mit V-342.
- **Kommentare (§3.11):** ja, mit V-340.
- **`harness/README.md`, `harness/conventions.md`, Commands:** stimmen.
- **Keine Lockerung (§3.6):** bestätigt (Abschnitt 3).
- **Größe und Prüfbarkeit:** 730 in der Metrik von `AGENTS.md` §3.14 (Schätzung 550 bis 750, Marke 900; Basis `8349a65`, letzter Code-Commit `9e727fb`). Prüfbar in einer Sitzung: ja. Hook (206 Zeilen) und die Tabellenform der Gegenprobe (Fall = eine Zeile mit Punktnummer, Mutant = eine `rot`-Zeile mit Rot- und Grünliste) tragen die Prüfung; ich habe Hook und Spezifikation vollständig gelesen, die Gegenprobe über Fälle, Kopf und Mutanten und alle Fälle durch Mutation gefahren.

## 7. Offen und nicht gefahren

- **Offen, weil die Closure sie führt:** Closure-Notiz mit Lerneintrag, Beobachtungs-Register (Evidence und Ausgang von `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`), Ausgang jedes Risikos aus §6, die drei Paarungen, `git mv` nach `done/`.
- **Nicht gefahren:** `make gates` (auf Anweisung; V-341), `make test-integration`, `make lint-gegenprobe`, `make a-check-negativ`, `make abdeckung-check`; der Hook in der Umgebung von Claude Code (Grenze der Gegenprobe, Plan §6, `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung`); die Gegenprobe als root und unter bash 3; ein echter Doku-Commit mit `git commit` und Stop in Claude Code. Der Stand der schnellen Kette war in den Hook-Proben ein Stellvertreter (die Ausgabe von `working-tree-hash.sh` in die Datei gelegt), außer in der F-713-Probe, die `make gates-schnell` in der Kopie fuhr. Nicht jede Zeile der Mutationstabellen in §7 habe ich nachgefahren (zwanzig eigene plus R1; die übrigen fährt die Gegenprobe selbst unter ihren Mutanten, grün in Probe 1).

## Verdikt

Die Liefer-Punkte 1 bis 3, Review und Mutationen sind erfüllt; die Entscheidungs-Konformität ist gegeben; der Hook gibt in keiner meiner Proben öfter frei als `SPEC-052` zusagt, und die schnelle Kette gibt nicht für eine Menge frei, für die der volle Lauf rot wäre (F-713 geheilt, kein weiterer Leser der Positivliste unter den langsamen Gates). `make gates` ist laut Implementer am Stand `5fed425` grün, am `HEAD` nicht belegt (V-341); nicht von mir gefahren. Es bleiben die Closure-Pflichten und drei kleine Nachzüge für den Implementer bzw. Architect: V-340 (Satz im Kopf des Hooks enger fassen), V-341 (voller Lauf nach dem letzten nicht inerten Commit vor der Closure), V-342 (Kopf des Plans). Kein Rückgabegrund an den Planner.
