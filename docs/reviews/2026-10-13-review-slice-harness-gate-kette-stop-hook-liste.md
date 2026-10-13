# Review-Report: slice-harness-gate-kette-stop-hook-liste — 2026-10-13

**Review-Art:** Code-Review (Harness-Werkzeug: Stop-Hook, Dateiliste, Gegenprobe). Geprüft wird gegen Plan §1, §3, §4, §6 (Randform-Tabelle, Mutationen, Entscheidungen des Architect), §7 (Belege), §8, gegen [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) und `SPEC-052` Teil 1 (Punkte 1 bis 8) und Teil 2 (Punkte 9 bis 22) und gegen `AGENTS.md` §3.3, §3.6, §3.7, §3.9, §3.10, §3.11, §3.12. Die DoD prüft der Verifier.

**Gegenstand:** `9662e13..HEAD`, Pfade `tools harness .claude spec` und der Slice-Plan. Code-Commit `b259381` (Hook, `tools/harness/dateiliste.sh`, `record-gates.sh`, `working-tree-hash.sh`, `stop-hook-gegenprobe.sh`, Fragment `harness/mk/stop-hook-gegenprobe.mk`, Ausnahmeliste, Gate-Zeile, Träger), Belege `e758865`, Entscheidungen des Architect `c005555` (`SPEC-052` Punkte 13, 20, 21, 22 und Plan §6), Nacharbeit `e92db76` (`record-gates.sh` kennt `-i`, Fälle `p21-*`, `p22-*`), Plan-Commits `32a5a47`, `78b35ab`, `9a1cc9d`.

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-13

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `slice-harness-gate-kette-stop-hook-liste` ganz gelesen (§1 bis §8, mit den Entscheidungen des Architect vom 2026-10-13, dem Urteil des Planners zur Marke und §7 *Belege des Implementers*); den Bericht des Implementers habe ich nicht gelesen.
- [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) (Entscheidung Punkt 3: Hook, fail-closed, fehlender Stand auch bei sauberem Baum); `SPEC-052` Punkte 1 bis 22 Wort für Wort gegen den Hook, `dateiliste.sh`, `record-gates.sh` und die Fälle gehalten.
- Der alte Hook und die alten Skripte (`git show 9662e13:…`) zum Vergleich der Strenge; `harness/mk/enforce.mk`, `Makefile`, `.claude/settings.json` (Aufruf des Hooks), `.claude/commands/*.md`.
- `AGENTS.md` (Hard Rules). Vorheriger Report am selben Modul: `2026-10-13-review-slice-harness-gate-kette-schnell`.

**Ausgeführte Läufe:**

- `make stop-hook-gegenprobe`: grün, 52 Fälle, 20,8 s (nicht als root, daher läuft `p20-stand-unlesbar`).
- `make gates-schnell-gegenprobe`: grün, 21,6 s. `make gates-schnell`: grün, 7,8 s; der Arbeitsbaum blieb danach sauber (`git status`).
- Metrik des Planners selbst gemessen: `git diff --numstat 9662e13..HEAD -- tools harness .claude`: 845 Einfügungen, 75 Löschungen, Summe 920, netto 770 (stimmt mit §7). Nach der Nacharbeit sind es 142 mehr als die 778, die das Urteil des Planners meint.
- Dateiliste von Hand in einem Temp-Repo mit Pfaden ` lead`, `trail `, Tab, `a]b`, `q"x`, `bs\x`, `$(touch pwn)`, `*`, `-dash`, Nicht-UTF-8, hängendem Link und Link auf ein Verzeichnis: Erzeugung, `dateiliste_pruefen` (Status 0), kein Befehl aus dem Pfad ausgeführt, Liste byte-gleich bei zweitem Lauf, `gates-passed.diffsha` = SHA-256 der Liste. Die Datei der Liste hat Modus 0600 (aus `mktemp`), ohne Wirkung.
- Hook gegen ein `SIGTERM` mitten im Lauf (hängendes `sha256sum` als Wrapper): Antwort `block`, Status 0.
- **Dreizehn Mutationen** in Kopien im Scratchpad (nicht im Arbeitsbaum), jede mit der unveränderten Gegenprobe der Kopie gefahren:

| Nr. | Mutation | Ergebnis |
|---|---|---|
| M1 | Prüfung des Konflikts hinter die Identität, vor `frei` verschoben | grün (äquivalent, siehe F-708) |
| M2 | fehlender Stand bei sauberem Baum freigeben | rot: `p19-ohne-stand-sauber` |
| M3 | Fehlerantwort `block` in `sicherung` entfernt | rot: `p20-git-u-scheitert`, `p20-git-liste-scheitert`, `p20-datei-nicht-lesbar`, `p20-abbruch`, `p20-lib-fehlt`, `p20-kein-repo`, `p20-ohne-git` |
| M5 | `-i` nicht erkannt (`*i*` unwirksam) | rot: alle sechs `p21-i-*` |
| M6 | `IFS=` im `read` von `dateiliste_erzeugen` entfernt (Pfade verlieren führende und abschließende Leerzeichen) | **grün** (F-703) |
| M7 | `read` ohne `-r` in `dateiliste_erzeugen` (Backslash im Pfad) | **grün** (F-703) |
| M8 | `sort -zu` ohne `LC_ALL=C` | rot: `p22-leere-umgebung` (die leere Umgebung sortiert in `C`, der Schreiber in der Locale) |
| M9 | keine Sortierung (`cat` statt `sort -zu`) | rot: `p11-form` |
| M10 | Schleifenschutz hinter `git rev-parse` und `source`, vor die Prüfung des Konflikts gesetzt | **grün** (F-704) |
| M11 | `mktemp` des Hooks im Arbeitsbaum | rot: `p16-identisch`, `p22-*` (Schreiben in den Baum macht die Liste ungleich) |
| M13 | Zweig für Symlinks entfernt | rot: `p11-form` |
| — | `chmod 000` des Standes entfernt (Simulation von root) | **Gate rot** (F-705) |

- `make docs-check` nach dem Schreiben dieses Reports (siehe Verdikt).
- **Nicht gefahren:** `make gates`, `make test-integration`, `make lint-gegenprobe`, `make hook-gegenprobe` (Auftrag); die 51 Mutanten der Gegenprobe habe ich nicht einzeln nachgezählt, die Gegenprobe fährt sie selbst (49 `fixture m-…` und zwei Anschluss-Mutanten gezählt, §7 nennt 41 und 10); den Hook nicht in der Umgebung von Claude Code (Grenze der Gegenprobe, Plan §6); ein echtes `make -i gates` gegen das echte `Makefile` mit allen Gates (nur gegen die Attrappen der Gegenprobe); Läufe unter bash 3.2 oder als root (Verhalten von mir nur abgeleitet beziehungsweise simuliert, F-705, F-710).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-703 | MEDIUM | Pfade mit führendem oder abschließendem Leerzeichen, mit Backslash, Tabulator oder Nicht-UTF-8 haben keinen Fall. Zwei Mutationen von `dateiliste_erzeugen` (M6 `IFS=` entfernt, M7 `read` ohne `-r`) lassen die Gegenprobe grün; die Fälle `p11-*` nennen nur ein Leerzeichen im Pfad. Failure-Szenario gefahren: Mit M6 steht ` lead` in der Liste als `G` (Pfad zu `lead` verkürzt, Datei „fehlt“), eine Änderung an der Datei ändert die Liste nicht, der Hook antwortet `approve` für einen Baum, den kein `make gates` gesehen hat. Die Spezifikation sagt „Zeilenumbruch, Leerzeichen, Tabulator und Nicht-UTF-8 im Pfad sind damit ohne Sonderfall“; der Plan §6 nennt Nicht-UTF-8 als Randform, ein Fall oder eine Mutation dazu steht nicht in §6 und §7. | `AGENTS.md` §3.10, §3.11; `SPEC-052` Punkt 11 | `tools/harness/dateiliste.sh` · `while IFS= read -r -d '' pfad; do`; `spec/spezifikation.md` · `Zeilenumbruch, Leerzeichen, Tabulator und Nicht-UTF-8 im Pfad sind damit ohne Sonderfall` | ja — M6 und M7 in einer Kopie | Zusage ohne fangende Mutation (Pfadform) |
| F-704 | MEDIUM | Die Zusage „der Schleifenschutz steht davor“ (`SPEC-052` Punkte 17 und 20: erste Prüfung nach dem Schleifenschutz; Schleifenschutz steht vor dem Abbruch-Pfad) hat keine fangende Mutation. M10 setzt den Schleifenschutz hinter `git rev-parse`, `cd` und `source`, unmittelbar vor `git ls-files -u`: Gegenprobe grün. Die zwei Fälle `p20-schleifenschutz` und `p20-schleifenschutz-false` haben einen gültigen Stand, keinen Konflikt und kein scheiterndes `git`; sie unterscheiden die Stelle des Schleifenschutzes nicht. Failure-Szenario: Ein Hook, der bei Konflikt, fehlendem Stand oder scheiterndem `git` mit `stop_hook_active: true` erneut `block` antwortet, erzeugt die Endlosschleife, die der Schleifenschutz verhindern soll; die Gegenprobe merkt es nicht. Das absichtliche fail-open selbst ist benannt (Kopf des Hooks: „antwortet `approve`, damit ein dauerhaft rotes Gate keine Endlosschleife erzeugt“; Kommentar am Code; zwei Fälle belegen den Grundfall), in der Spezifikation steht es nur als „steht davor und bleibt“ ohne zu sagen, dass es Konflikt und fehlenden Stand überstimmt, und der Kopf der Gegenprobe führt es nicht als Grenze auf. | `AGENTS.md` §3.10, §3.11; `SPEC-052` Punkte 17, 20 | `.claude/hooks/stop-require-gates.sh` · `if [[ "$input" =~ \"stop_hook_active\"[[:space:]]*:[[:space:]]*true ]]; then`; `tools/harness/stop-hook-gegenprobe.sh` · `f_p20_schleifenschutz() {` | ja — M10 in einer Kopie | Zusage ohne fangende Mutation (Ort des Schleifenschutzes, mehrteilig) |
| F-705 | LOW | Als root endet die Gegenprobe rot, obwohl ihr Kopf zusagt, der Fall `p20-stand-unlesbar` werde „laut übersprungen“. Der Fall überspringt sich (gibt nichts aus), aber der Mutant `m-unlesbar-als-unveraendert` verlangt von ihm eine Ausgabe und meldet „färbt 'p20-stand-unlesbar' nicht rot“; Exit 1 (Simulation: `chmod 000` im Skript einer Kopie entfernt, `[ -r ]` ist dann wahr). Das Gate ist als root nicht grün zu bekommen; der Kommentar beschreibt das Überspringen, nicht den Ausgang. | `AGENTS.md` §3.11 | `tools/harness/stop-hook-gegenprobe.sh` · `Als root bleibt eine Datei mit Rechten 000 lesbar; dann ist p20-stand-unlesbar nicht` | ja — Simulation wie beschrieben | Kommentar weiter als die Prüfung |
| F-706 | LOW | Lebende Stellen sagen den alten Nachweis: `.claude/commands/implement-slice.md` („`record-gates`, das einen Content-Hash des Working Tree stempelt“), `plan-welle.md` und `close-welle.md` („macht den Stempel ungültig“); der Plan §3 verlangt das `grep` über `.claude/` und das Nachziehen im selben Commit, die Befehle sind im Diff unberührt. Keine Aussage ist falsch im Ergebnis, aber sie nennen weder die Dateiliste noch die neue Sperre (fehlender Stand auch bei sauberem Baum, `make -i gates` schreibt keinen Nachweis), die jeder Implementer-Lauf jetzt trifft. Ebenso §3 des Plans (Zeile `record-gates.sh`: „schreibt die Dateiliste statt eines Hashs“): `gates-passed.diffsha` bleibt und ist der Hash der Liste. | `AGENTS.md` §3.9; Plan §3 (Prüfung vor dem Commit) | `.claude/commands/implement-slice.md` · `einen Content-Hash des`; `docs/plan/planning/in-progress/slice-harness-gate-kette-stop-hook-liste.md` · `schreibt die Dateiliste statt eines Hashs` | ja — `grep -rn "Content-Hash" .claude` | Plan-Auftrag „Stellen nachziehen“ nicht erfüllt |
| F-707 | LOW | Zwei Zusagen reichen weiter als ihr Fall. (a) Der Kommentar an `sicherung` nennt „Signal“ unter den Enden, die `block` antworten; gefahren (`SIGTERM`) stimmt es, aber kein Fall prüft es (`p20-abbruch` ist ein `exit 3`). (b) `SPEC-052` Punkt 22: „liest keine Umgebungsvariable“; der Fall `p22-leere-umgebung` prüft die leere Umgebung (und die Mutation `HOME`), der Hook liest aber `PATH` (Programmsuche), `TMPDIR` (`mktemp`) und `GIT_*` (über `git`). Der Satz ist enger als er klingt, der Fall belegt „leere Umgebung gibt dieselbe Antwort“. | `AGENTS.md` §3.11 | `.claude/hooks/stop-require-gates.sh` · `# Jedes Ende ohne Antwort (Abbruch durch set -e, fehlende Datei, Signal)`; `spec/spezifikation.md` · `und liest keine Umgebungsvariable` | ja — `kill -TERM`; Lesen von `mktemp` | Zusage im Satz weiter als die Prüfung |
| F-708 | INFO | Zu Prüfpunkt (1) Reihenfolge: Konflikt und Identität zu vertauschen ist äquivalent, solange beide Prüfungen vor `frei` stehen (M1: grün, kein stilles Grün: die Antwort ist in beiden Ordnungen `block`). Wirksam ist nur die Mutation, die den Konflikt hinter `frei` schiebt (`m-identitaet-zuerst`, rot an `p17-konflikt-identisch`). Die Ordnung „Konflikt vor fehlendem Stand“ ist nur am Text von `reason` sichtbar, nicht an der Antwort; die Gegenprobe prüft `reason` nicht, die Spezifikation sagt nur „erste Prüfung nach dem Schleifenschutz“. Keine Aktion erwartet. | Maintainability | `tools/harness/stop-hook-gegenprobe.sh` · `fixture m-identitaet-zuerst` | ja | Ordnung nur im Begründungstext sichtbar |
| F-709 | INFO | Sortierung und Einmaligkeit haben doch Fälle, obwohl Kopf der Gegenprobe und Punkt 11 „keinen eigenen Fall“ sagen: `p11-form` vergleicht mit einer bytegenau sortierten, von Hand geschriebenen Liste (M9 rot), und `p22-leere-umgebung` fängt `sort` ohne `LC_ALL=C` (M8 rot). Der Kommentar ist enger als die Prüfung; das ist die sichere Richtung, kein Befund. | Maintainability | `tools/harness/stop-hook-gegenprobe.sh` · `Sortierung und Einmaligkeit der Liste haben keinen eigenen Fall` | ja | Kommentar enger als die Prüfung |
| F-710 | INFO | Umgebungsabhängigkeiten des Hooks, die weder sein Kopf noch ein Fall nennt: `local -A` in `dateiliste_pruefen` (bash ab 4) und `readlink -z`, `sort -z` (GNU). Unter bash 3.2 oder BSD-Werkzeugen bricht der Hook ab und antwortet `block` (fail-closed, keine Freigabe), also dauerhaft. Ebenso `TMPDIR` im Arbeitsbaum: die temporäre Datei steht dann in der Liste, der Hook antwortet `block`. `mktemp` ohne Pfad schreibt sonst nach `/tmp`, nicht in den Baum; der Hook schreibt im Baum nichts (M11 rot). Das Repo setzt „Docker und GNU make“ voraus, nicht bash 4; INFO an den Architect, ob die Voraussetzung in den Kopf gehört. | Maintainability | `tools/harness/dateiliste.sh` · `local -A gesehen=()` | ja — Lauf unter bash 3.2 (nicht gefahren) | Werkzeug-Voraussetzung ungenannt |
| F-711 | INFO | Zu Prüfpunkt (10) Größe: Summe 920, netto 770 (`stop-hook-gegenprobe.sh` 659 von 920 Zeilen); der Planner maß 778 vor der Nacharbeit und entschied „Fortsetzung“ unter der Bedingung dieses Satzes. Die Nacharbeit liegt mit +142 über der Erwartung des Architect (+60 bis 95), und das Urteil des Planners bezieht sich auf den Stand 778. **Der Diff war in einer Sitzung prüfbar** (ich habe Hook, Liste, Schreiber, die Gegenprobe als Tabelle, die Spezifikation und den Plan gelesen und dreizehn Mutationen gefahren), aber nur so, dass die Gegenprobe gegen die Mutationstabelle gelesen und nicht Zeile für Zeile als Logik; ohne die Tabelle in §7 wäre es knapp. Der Planner entscheidet nicht neu, wenn er diese Lesart bestätigt. | Maintainability; Plan §4 | `docs/plan/planning/in-progress/slice-harness-gate-kette-stop-hook-liste.md` · `Der Reviewer meldet in seiner Summary-Zeile` | ja — `git diff --numstat` | Größe über der Schätzung |
| F-712 | INFO | Zu Prüfpunkt (7), `AGENTS.md` §3.6: Der Hook sperrt jetzt öfter (fehlender Stand auch bei sauberem Baum; Konflikt; unlesbare oder beschädigte Liste) und gibt nirgends mehr frei als der alte: der alte Hash umfasste dieselbe Menge (`--cached --others --exclude-standard`, Link mit Ziel, `GONE`), die neue Liste unterscheidet zusätzlich `X` und trägt Modus wie der Hash nicht (akzeptiertes Negativ, Punkt 11). Eine Verschärfung nach [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) Punkt 3 und Entscheidung des Nutzers vom 2026-10-12, keine Lockerung; kein ADR nötig. Ignorierte Dateien sind nicht Teil der Liste (Punkt 18, akzeptiertes Negativ): im Bestand ist nur `.harness/state/` ignoriert (`git status --ignored`), ein Gate liest daraus nichts. | `AGENTS.md` §3.6 | `.claude/hooks/stop-require-gates.sh` · `# - Ein fehlender, nicht lesbarer oder beschädigter Stand sperrt, auch bei sauberem Baum` | ja | — |

## Zu den Prüfpunkten des Auftrags

1. **Hook gegen Punkte 9 bis 22.** Reihenfolge: Schleifenschutz → `git rev-parse` und `source` (Aufbau, keine Entscheidung) → Konflikt (`git ls-files -u`) → fehlender Stand → `dateiliste_pruefen` → Liste erzeugen → `cmp -s` (byte-gleich) → `frei`; jedes andere Ende `voll`. Entspricht den Punkten 16, 17, 19, 20. Antwort über Dateikennung 3 (der Befund der Gegenprobe vor/nachher in §7 ist belegt: ohne `exec 3>&1` hielte die Umleitung die Antwort fest); `sicherung` antwortet bei jedem Ende ohne Antwort `block` mit Status 0, auch bei `SIGTERM` (gefahren). Temp-Datei über `mktemp` in `$TMPDIR` oder `/tmp`, vom Trap entfernt; keine Schreibzugriffe auf den Baum (M11 rot). Das absichtliche fail-open (`stop_hook_active` → `approve`): Der Wortlaut des Hooks beschreibt, was da ist (§3.7); der Test ist nicht durch Text im JSON zu täuschen (ein Wert in einer Zeichenkette steht als `\"…\"`, das Muster verlangt `"stop_hook_active"` ohne Gegenschrägstrich davor). Belegt ist der Grundfall; Ort und Vorrang gegenüber Konflikt, fehlendem Stand und `git`-Fehler sind es nicht (F-704).
2. **Dateiliste.** Format, Arten `F`/`L`/`G`/`X`, Wert, Pfad, NUL-getrennt, bytegenau sortiert (`LC_ALL=C sort -zu`), Symlink mit Ziel (`readlink -z`), gelöscht (`G`), umbenannt (zwei Einträge), untracked, `.harness/state/` ignoriert und `.git` nie in `git ls-files`: ohne Befund. `dateiliste_pruefen` (leer, Feldzahl, Hash, doppelt) mit Sonderzeichen im Pfad gefahren (Status 0, kein Befehl ausgeführt). `gates-passed.diffsha` = SHA-256 der Liste (`p1-hash-neben-liste`). Befund: F-703.
3. **`record-gates.sh`.** Grüner Lauf schreibt Liste und Hash, `-i` (jede Schreibweise des ersten Worts von `MAKEFLAGS`, auch `ik`, `i -j4 …`, Umgebung) entfernt beide, `-n`/`-q`/`-t` hält `make` selbst (das Rezept in `enforce.mk` trägt weder `$(MAKE)` noch `+`; im Diff unberührt), kombiniert mit `-i` ebenso: `p21-*` und die Mutanten `m-21-*` decken jede Bedingung (Entfernen nur der Liste, nur des Hashs, nichts, `+`, `+` ohne Voraussetzung); M5 rot an allen sechs `p21-i-*`. Ohne Befund bis auf den Hinweis, dass der Nachweis bei rotem Lauf ohne `-i` stehen bleibt (akzeptiertes Negativ nach Punkt 21).
4. **Gegenprobe.** Tabellengetrieben (Fälle in `FAELLE`, Mutanten als Zeilen mit Rot- und Grünlisten), Attrappen als Gates, Temp-Repos, ohne Docker, jeder Fall mit Punktnummer; `rot` verlangt auch, dass die genannten grünen Fälle grün bleiben. Stilles Grün: F-703 und F-704 sind die zwei Stellen, an denen eine Mutation grün bleibt; F-705 ist ein stilles Rot als root. Verneinende Zusagen („kein Nachweis bei rotem Lauf“, „`-n`/`-q`/`-t` entfernen nichts“, „der Hook schreibt nicht“): die ersten beiden haben Fälle und Mutanten, die dritte ist nicht zugesagt (weder im Hook noch in `SPEC-052`) und wirkt nur indirekt über `p16-identisch` (M11).
5. **Kommentare, Gate-Zeile, Träger.** Kopf von Hook, `dateiliste.sh`, `record-gates.sh` und Gegenprobe beschreiben den Zustand, nicht eine Alternative oder einen abwesenden Text; Befunde nur F-705, F-707. Die Zeile in `harness/README.md` und die Träger in `harness/conventions.md` stimmen mit dem Code überein (Fälle der Punkte 1, 10, 11, 14, 16 bis 22; Fragment; `dateiliste.sh`). `SPEC-052` Punkt 13: Die Temp-Repos von `stop-hook-gegenprobe` tragen nur Makefile, `enforce.mk`, Hook, Skripte und Probedateien (`fixture()`), die Gegenprobe liest aus dem Arbeitsbaum nur die sieben genannten Pfade, weder `docs/user/` noch `README.md`; der Beleg gilt.
6. **Ausnahmeliste.** `stop-hook-gegenprobe` steht in `GATE_SCHNELL_AUSNAHME` und in `LANGSAM` der Gegenprobe der Kette, die Kopien dort tragen `dateiliste.sh`; `make gates-schnell-gegenprobe` grün, `make gates-schnell` fährt das Gate nicht.
7. **§3.6.** F-712.
8. **§3.9.** §1, §6, Kopf, §8 folgen dem Diff (Zählung 3 Liefer-Punkte, 1 Schicht; Namen; Entscheidungen des Architect mit Commit); Befund nur F-706.
9. **Läufe.** Siehe oben.
10. **Größe.** F-711.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.claude/hooks/stop-require-gates.sh` | geprüft; F-704, F-707, F-710; sonst ohne Befund (Reihenfolge, Dateikennung 3, Trap, Status 0, Quoting, kein Schreiben in den Baum) |
| `tools/harness/dateiliste.sh` | geprüft; F-703, F-710; sonst ohne Befund (Format, Sortierung `C`, Links, `G`/`X`, `pipefail`, Sonderzeichen im Pfad) |
| `tools/harness/record-gates.sh`, `harness/mk/enforce.mk` (unberührt) | geprüft, ohne Befund (`-i` in jeder Form, `-n`/`-q`/`-t` vom Rezept gehalten, atomares Schreiben nicht zugesagt) |
| `tools/harness/working-tree-hash.sh` | geprüft, ohne Befund (Hash = SHA-256 der Liste, Fehler des Erzeugers schlägt durch, `pipefail`) |
| `tools/harness/stop-hook-gegenprobe.sh`, `harness/mk/stop-hook-gegenprobe.mk` | geprüft; F-703, F-704, F-705, F-707, F-708, F-709; sonst ohne Befund (Arbeitsbaum bleibt unberührt, Anschluss an `GATE_CHECKS`, Exit gilt, Attrappen, `GIT_CEILING_DIRECTORIES`) |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh` | geprüft, ohne Befund (Eintrag und Kopie mit `dateiliste.sh`, Gegenprobe grün) |
| `harness/README.md` (Gate-Zeile), `harness/conventions.md` (Träger-Zeile) | geprüft, ohne Befund |
| `spec/spezifikation.md` `SPEC-052` Punkte 13, 20 bis 22, Historie | geprüft; F-707 (b); sonst ohne Befund (kein ADR-/Slice-Bezug im Text, Punkte 21 und 22 mit Fall, Nicht-Gegenstand benannt) |
| Slice-Plan §1, §3, §4, §6, §7, §8, Kopf | geprüft; F-706, F-711; sonst ohne Befund (Rückgaben stehen als Rückgabe, nicht im Code-Commit entschieden: die Randformen 21 und 22 hat der Architect in `c005555` vor `e92db76` entschieden, `AGENTS.md` §3.12) |
| `.claude/commands/` | geprüft; F-706 |
| Commits | geprüft, ohne Befund: kein Move, Inhalt und Belege getrennt (§3.3), Messages nennen Slice und [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) |
| Lint-Profil, Produkt-Code, Tests unter `internal/` und `test/` | im Diff unberührt, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 5 |

**Prüfbarkeit in einer Sitzung: ja.** Der Diff (Summe 920, netto 770 Zeilen, davon 659 in der tabellengetriebenen Gegenprobe) war in einer Sitzung prüfbar, wenn die Gegenprobe gegen die Mutationstabelle in §7 gelesen wird; der Planner muss nicht neu entscheiden. Die Nacharbeit (+142) liegt über der Erwartung des Architect.

**Finding-Klassen dieses Laufs:** Zusage ohne fangende Mutation (Pfadform) · Zusage ohne fangende Mutation (Ort des Schleifenschutzes, mehrteilig) · Kommentar weiter als die Prüfung · Plan-Auftrag „Stellen nachziehen“ nicht erfüllt · Zusage im Satz weiter als die Prüfung · Ordnung nur im Begründungstext sichtbar · Werkzeug-Voraussetzung ungenannt

## Verdikt

**Merge-blockierend:** nein im Sinn von HIGH. Der Hook ist fail-closed in allen gefahrenen Fehlerformen (Abbruch, `SIGTERM`, `git`-Fehler, fehlende Dateien, beschädigte Liste: `block`, nie keine Antwort) und sperrt in keinem Fall weniger als der alte. F-703 und F-704 sind vor der Closure zu schließen (§3.10: Zusage ohne Mutation; F-703 hat ein gefahrenes Fehlverhalten, `approve` für eine geänderte Datei mit Leerzeichen am Pfadrand, nur unter einer Mutation des Codes, nicht im Code von heute).

**Entscheidungen, die der Architect treffen muss:**

1. F-703: Fälle und Mutationen für Pfade mit führendem oder abschließendem Leerzeichen, Backslash, Tabulator und Nicht-UTF-8 (je Zusage eine Mutation), oder den Satz „Zeilenumbruch, Leerzeichen, Tabulator und Nicht-UTF-8 im Pfad sind damit ohne Sonderfall“ in Punkt 11 auf das fassen, was die Fälle prüfen (§3.11, nicht beides offen lassen).
2. F-704: ein Fall, der `stop_hook_active: true` mit Konflikt, mit fehlendem Stand und mit scheiterndem `git` zeigt (jeweils `approve`), oder die Formulierung „steht davor“ in den Punkten 17 und 20 enger fassen; dazu, ob das fail-open des Schleifenschutzes (überstimmt Konflikt und fehlenden Stand) in `SPEC-052` und im Kopf der Gegenprobe ausdrücklich als akzeptierte Grenze steht. Es ist Bestand des alten Hooks, keine neue Lockerung.
3. F-707 (b): „liest keine Umgebungsvariable“ in Punkt 22 auf „leere Umgebung gibt dieselbe Antwort“ fassen, oder einen Fall je Variable (`PATH`, `TMPDIR`, `GIT_*`) liefern.
4. F-705: ob die Gegenprobe als root rot sein darf (dann sagt der Kommentar es) oder der Mutant für den übersprungenen Fall ebenfalls entfällt.
5. F-710: ob die Voraussetzung bash ab 4 und GNU `readlink`/`sort` in den Kopf des Hooks und der Spezifikation gehört.

**Übergabe:** Findings gehen an den Implementer (F-703 bis F-707 mit der Entscheidung des Architect zu 1 bis 5; F-706 allein beim Implementer); die Finding-Klassen zusätzlich in die Slice-Closure §7. Gate-Ergebnis: `make gates` nicht gefahren (Auftrag des Aufrufers); `make stop-hook-gegenprobe`, `make gates-schnell-gegenprobe` und `make gates-schnell` grün; `make docs-check` siehe Commit.
