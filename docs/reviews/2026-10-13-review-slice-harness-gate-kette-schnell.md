# Review-Report: slice-harness-gate-kette-schnell — 2026-10-13

**Review-Art:** Code-Review (Harness-Werkzeug und Spezifikation). Geprüft wird gegen Plan §1, §3, §4, §6 (Randform-Tabelle, Mutationen), §7 (Belege), §8, gegen [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) und gegen `AGENTS.md` §3.3, §3.6, §3.7, §3.9, §3.10, §3.11, §3.12. Die DoD prüft der Verifier.

**Gegenstand:** `d075d3f..HEAD` (drei Commits: `cb37a34` `SPEC-052` und Historie, `f633694` Fragment, zwei Skripte, Gate-Zeilen, Träger-Zeile, `a203246` Belege in §7). Geändert: `spec/spezifikation.md`, `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell.sh`, `tools/harness/gates-schnell-gegenprobe.sh`, `harness/README.md`, `harness/conventions.md`, Slice-Plan (nur §7).

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-13

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `slice-harness-gate-kette-schnell` ganz gelesen (§1 bis §8, mit der Randform-Tabelle und den Entscheidungen des Architect vom 2026-10-13 und §7 *Belege des Implementers*); den Bericht des Implementers habe ich nicht gelesen.
- [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) ganz; `SPEC-052` (Teil 1 und 2) Wort für Wort gegen die Tabelle in §6 gehalten, der Eingangssatz von Teil 2 gegen die Zeile *Zwischenstand*.
- `tools/harness/working-tree-hash.sh`, `harness/mk/enforce.mk`, `Makefile` (Aggregator), `harness/mk/kopf-check.mk` als Vorbild.
- `AGENTS.md` (Hard Rules). Vorherige Findings am selben Modul: keine Reports zur Gate-Kette.

**Ausgeführte Läufe:**

- `make gates-schnell-gegenprobe`: grün, 4,1 s.
- `make gates-schnell` (Cache warm): grün, 7,3 s; Stand `.harness/state/gates-schnell-passed.diffsha` gleich der Ausgabe von `working-tree-hash.sh`; der Nachweis `gates-passed.diffsha` blieb unberührt. Die Laufzeitangabe in §7 (7,2 bis 7,5 s) ist damit belegt, nicht geschätzt.
- `make docs-check`: 673 Dateien, 0 Befunde.
- `make -n gates-schnell` (Trockenlauf) selbst gefahren: siehe F-694.
- Drei Mutationen der Gegenprobe selbst in Kopien im Scratchpad (A: Überschuss `docs-check` in der Ausnahmeliste; B: Hash vor dem Lauf statt danach; C: `record-gates.sh` löscht den Stand der schnellen Kette), dazu zwei weitere (D: neues Gate samt Eintrag in der Liste; E: `gate_variable` liest eine falsche Variable). Ergebnis: A, D, E rot; B und C **grün** (F-695).
- **Nicht gefahren:** `make gates`, `make test-integration`, `make lint-gegenprobe`, `make hook-gegenprobe` (Auftrag); die Mutationstabelle des Implementers in §7 habe ich nicht Zeile für Zeile wiederholt, nur die oben genannten fünf; ein Lauf der schnellen Kette mit geändertem Go-Quelltext (102 s in §7) und der rote Lauf in §7 sind nicht nachgefahren. Der Stop-Hook, `record-gates.sh` und `working-tree-hash.sh` sind im Diff unberührt, ihr Verhalten habe ich nicht gefahren.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-694 | MEDIUM | Ein Trockenlauf `make -n gates-schnell` führt das Rezept aus, weil die Zeile `$(MAKE)` enthält: Das Skript entfernt den Stand, fährt `make -n <Ziele>` (druckt, fährt nichts, Exit 0) und schreibt dann den Stand. Gefahren: Exit 0, `gates-schnell-passed.diffsha` entsteht mit dem Hash des Arbeitsbaums, ohne dass ein Gate lief. Das widerspricht Punkt 5 („Bei einem grünen Lauf … schreibt sie den Hash“) und der Zusage „Stand nur bei grünem Lauf“; der Hook der Folge-Slices liest genau diese Datei als „grün“. `make -n gates` hat das Problem nicht (Rezept ohne `$(MAKE)`), und die Gegenprobe fährt nur `-n gates`, nie `-n gates-schnell`. | `SPEC-052` Punkt 5 und 6; `AGENTS.md` §3.10 | `harness/mk/gates-schnell.mk` · `@bash tools/harness/gates-schnell.sh '$(MAKE)'` | ja — `make -n gates-schnell; ls .harness/state` | Stand ohne Lauf (Trockenlauf schreibt) |
| F-695 | MEDIUM | Drei Zusagen von `SPEC-052` haben keine fangende Mutation. (a) Punkt 5 „Hash des Arbeitsbaums **zum Ende** des Laufs“: Mutant B (Hash vor `"$make_bin" "$@"` berechnet) bleibt grün, weil die Attrappen den Baum nicht ändern. (b) Punkt 7 „der volle Lauf berührt den Stand der schnellen Kette nicht“: Mutant C (`record-gates.sh` löscht den Stand) bleibt grün, die Gegenprobe fährt `make gates` nie, nur `-n gates`. (c) Punkt 5 „Nur die schnelle Kette schreibt diese Datei“ hat denselben Befund wie (b). Beide Mutanten gefahren: Ausgabe `gruen — Punkte 1 bis 7 …`, Exit 0. Die Tabelle in §7 nennt keine Zeile dazu, §6 *Mutationen* auch nicht; die Zusage steht trotzdem in der Spezifikation. | `AGENTS.md` §3.10, §3.11 | `tools/harness/gates-schnell-gegenprobe.sh` · `# Punkt 7: ohne vollen Nachweis entsteht keiner.`; `spec/spezifikation.md` · `der volle Lauf berührt den Stand der\n   schnellen Kette nicht.` | ja — Mutanten B und C in einer Kopie | Zusage ohne fangende Mutation |
| F-696 | MEDIUM | `SPEC-052` sagt „Die Punkte 1 bis 8 prüft `make gates-schnell-gegenprobe`“, und der Kopf der Gegenprobe sagt „gegen die Punkte 1 bis 8“. Punkt 8 (aufruf-Voraussetzung) hat nach demselben Kopf („hat keinen Fall“), nach der Zeile in `harness/README.md` („Punkt 1 bis 7“) und nach der Schlusszeile der Gegenprobe („Punkte 1 bis 7“) keinen. Der Satz in der Spezifikation reicht weiter als die Prüfung. Die Wendung stammt nicht aus §6; der Implementer hat sie ergänzt. | `AGENTS.md` §3.11 | `spec/spezifikation.md` · `Die Punkte 1 bis 8 prüft `make gates-schnell-gegenprobe`.`; `tools/harness/gates-schnell-gegenprobe.sh` · `gegen die Punkte 1 bis 8 von SPEC-052` | ja — Vergleich mit der Schlusszeile | Zusage im Satz weiter als die Prüfung |
| F-697 | LOW | Die Ausnahmeliste hat sieben Einträge (mit `gates-schnell-gegenprobe`), Punkt 13 von `SPEC-052` und die Zeile in §6 zählen sechs auf („Kein Gate der Ausnahmeliste (`lint-gegenprobe` … `test-integration`) liest …“). Punkt 2 und 4 verlangen jedes Gate mit `-gegenprobe` im Namen in der Liste, der Name ist in §1 entschieden: Der Implementer hat die einzig mögliche Lesart genommen, aber die Aufzählung in Punkt 13 stimmt damit nicht mehr mit der Liste überein, und §6 (Plan folgt der Korrektur, §3.9) ist nicht nachgezogen (§7 sagt ausdrücklich „§6 ist nicht geändert“). Entscheidung des Architect, siehe unten. | `AGENTS.md` §3.9, §3.11 | `spec/spezifikation.md` · `(`lint-gegenprobe`, `kopf-check-gegenprobe`,\n    `abdeckung-gegenprobe`, `hook-gegenprobe`, `a-check-negativ`, `test-integration`)`; `harness/mk/gates-schnell.mk` · `GATE_SCHNELL_AUSNAHME :=` | ja — Zählung | Aufzählung weicht von der Liste ab |
| F-698 | LOW | Der Kopf von `gates-schnell.mk` sagt „Beide laufen auf dem Host mit bash, make, git, sha256sum, sort und coreutils“. Die Kette fährt `build`, `test`, `lint`, `a-check`, `docs-check` über Docker (gelaufen: `docker build`, `docker run`); nur die Gegenprobe läuft ohne Docker. `harness/README.md` sagt es richtig („Docker für die Gates selbst“). Die Werkzeugliste der Gegenprobe (`setsid`, GNU `sed -i`, `mktemp`) hat keinen Test. | `AGENTS.md` §3.11 | `harness/mk/gates-schnell.mk` · `Beide laufen auf dem Host mit bash, make, git, sha256sum,` | nein | Kommentar weiter als die Prüfung |
| F-699 | INFO | `SPEC-052` Punkt 3: „ein Fehler der Liste kostet Zeit, nie ein Gate“. Das gilt für ein Gate, das fehlt; ein Überschuss in der Liste (ein schnelles Gate irrtümlich eingetragen) nimmt es aus der schnellen Kette. Gefahren (Mutanten A und D): beide werden rot, aber nur über die Kopie (`p2-schnell` bei den acht bekannten Gate-Namen, `p4-eintrag` bei einem neuen Namen, den die Attrappen nicht führen); die Prüfung im Arbeitsbaum (`deckung`) ist einseitig. Das Ergebnis trägt, die Meldung nennt nicht die Ursache. Hinweis an den Architect, keine Aktion erwartet. | Maintainability | `tools/harness/gates-schnell-gegenprobe.sh` · `ist="$(deckung "$wurzel" | sort -u | tr '\n' ' ')"` | ja | Überschuss in der Liste nur indirekt gefangen |
| F-700 | INFO | Die Prüfung der Deckung wählt Namen mit `*-gegenprobe`; die Spezifikation sagt „`-gegenprobe` im Namen“. Ein Name wie `x-gegenprobe-y` fiele durch. Bestand und die künftige Kennung `stop-hook-gegenprobe` passen. | Maintainability | `tools/harness/gates-schnell-gegenprobe.sh` · `*-gegenprobe | a-check-negativ | test-integration)` | ja | Muster enger als die Wortlaut-Regel |
| F-701 | INFO | §7 nennt einen roten Lauf der Kette (Kommentar ohne Leerzeile am Ende von `cmd/pgwire-recorder/main.go`, Exit 2) als Beleg für „kein Stand nach Rot“ und sagt, den Grund habe der Implementer nicht untersucht. Als Beleg trägt er nur den Exit und das Fehlen des Stands; die Gegenprobe deckt den Fall mit einer Attrappe (`p6-rot`). | Maintainability | `docs/plan/planning/in-progress/slice-harness-gate-kette-schnell.md` · `den Grund des Rots habe ich nicht untersucht` | nein | Beleg ohne Ursache |
| F-702 | INFO | `gates-schnell-gegenprobe` hängt an `GATE_CHECKS`; `make gates` fährt damit ein Gate mehr (rund 4 s). Das ist die Zusage von Punkt 1 und 4 und keine Lockerung nach `AGENTS.md` §3.6; `record-gates.sh`, Hook und Nachweis sind unverändert. | `AGENTS.md` §3.6 | `harness/mk/gates-schnell.mk` · `GATE_CHECKS += gates-schnell-gegenprobe` | ja | — |

## Zu den vorgelegten Auslegungen und Prüfpunkten

1. **`SPEC-052` gegen §6.** Die Zeilen *Zugehörigkeit*, *Bezugsgröße*, *Roter Lauf*, *Form der Dateiliste*, *gelöscht/umbenannt*, *Positivliste/Grenze*, *Beleg*, *Symlink*, *Leerer Diff*, *Merge-Konflikt*, *Untracked*, *Erster Lauf*, *nicht lesbar*, *beschädigt*, *gleichzeitig* sind inhaltlich übertragen; bei Punkt 13 sind Attribuierungen („Entscheidung des Nutzers“, „Beleg des Architect“, „Mutation 2“) entfernt, „Ehrlich:“ gestrichen, „Mutation 2“ durch „Punkt 4“ ersetzt: Übertragung, keine neue Entscheidung. Der Eingangssatz von Teil 2 steht wörtlich. Teil 2 ist durch den Eingangssatz nicht zugesagt (`make help` führt kein Ziel `stop-hook-gegenprobe`). Keine Nennung von ADR, Slice, Welle oder Commit im neuen Text. kein Widerspruch zu [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) gefunden; die Positivliste mit `docs/user/` und `README.md` ist die Setzung des Nutzers vom 2026-10-13 (ADR sagt „zu Beginn nur …“, die Aufnahme ist dort als Änderung der Spezifikation mit Beleg vorgesehen). Ergänzungen des Implementers: die Rückbezüge, der Absatz „Der Name jedes Falls … Punkte 9 bis 20 gelten nach dieser Regel“ (aus der Zeile *Zwischenstand*, Vorgabe an die Hook-Slices) und der Satz „Die Punkte 1 bis 8 prüft …“ (F-696). Der Satz in Punkt 9 „Der Hook vergleicht …“ steht im nicht zugesagten Teil.
2. **Auslegung 1 (`gates-schnell` nicht an `GATE_CHECKS`).** Zwingend: Hinge es dran, enthielte `GATE_CHECKS` ohne die Ausnahmeliste das Ziel selbst (Rekursion), und `make gates` fuhre die Kette (Mutation „die Kette in seine Aufrufkette hängen“ in §6). §6 und `SPEC-052` sagen nur „die Gegenprobe hängt an `GATE_CHECKS`“; der Wortlaut der DoD („hängt mit der Kette an `GATE_CHECKS`“) ist die ungenaue Stelle, nicht der Vertrag. Kein stillschweigender Vertragswechsel in der Spezifikation; die DoD gehört dem Architect/Planner zur Klarstellung.
3. **Auslegung 2 (sieben Einträge).** Zwingend aus Punkt 2 und 4 plus dem fest vergebenen Namen; die Aufzählung in Punkt 13 und §6 ist dadurch veraltet (F-697). Kein neuer Vertrag, aber eine Spezifikationsstelle, die der Architect nachzieht.
4. **Verhalten Wort für Wort.** Kette = `GATE_CHECKS` ohne Liste: ja (`$(filter-out …)` über `GATE_CHECKS`, Rezept spät ausgewertet). Neues Gate läuft in beiden Ketten: ja (`neu-x`). Stand = unveränderte Ausgabe von `working-tree-hash.sh` (Skript unverändert): ja. Stand nur bei grünem Lauf: bis auf den Trockenlauf (F-694) ja; `rm -f` vor dem ersten Ziel deckt Rot und Abbruch (Mutant `m-entfernen-erst-bei-rot` fängt die Spätvariante). Ersetzt den vollen Nachweis nie: das Skript berührt `gates-passed.diffsha` nicht. `make gates` fährt die Kette nicht und ist unverändert (`Makefile`, `record-gates.sh`, `enforce.mk` im Diff unberührt).
5. **§3.10.** Fangende Mutationen vorhanden für Punkt 1 (zwei), 2 (drei), 3, 4 (beide Richtungen), 5 (Schreiben, Pfad), 6 (Rot, Abbruch, Spätentfernen), 7 (Pfad, neu entstehend), Anschluss und Exit (drei). Ohne: F-694, F-695. Die Gegenprobe baut Kopien mit Attrappen als Gates, `setsid -w` mit `kill -KILL 0` als Abbruch, ohne Docker; ein stilles Grün durch falsche Ausnahmeliste fängt sie (Mutanten A, D, E).
6. **§3.9.** §1, §3, §6, §8 und der Kopf des Plans widersprechen dem Diff nicht; Ausnahme: die Aufzählung in §6 (F-697). Größe: 434 Zeilen laut `git diff --shortstat` ohne Plan, unter der Rückführungsgrenze von 650 und innerhalb der Schätzung (420 bis 580). Eine Review-Sitzung hat getragen.
7. **§3.6.** Keine Lockerung: `make gates` fährt ein Gate mehr, keines weniger.
8. **Robustheit.** `set -euo pipefail`, Quoting und Pfade mit Leerzeichen in beiden Skripten ohne Befund; `'$(MAKE)'` in Einzelanführungszeichen ist für den Pfad `make` ausreichend. Kein `//nolint`, `.golangci.yml` unberührt. Aufruf-Voraussetzung (nacheinander) steht in Punkt 8, ohne Fall (akzeptiert nach §6).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/spezifikation.md` §11 `SPEC-052`, §12 Historie | geprüft; Befunde F-696, F-697; sonst ohne Befund (kein ADR-/Slice-Bezug, Eingangssatz wörtlich, Teil 2 nicht zugesagt, keine Randform neu entschieden) |
| `harness/mk/gates-schnell.mk` | geprüft; F-694, F-698; Ausnahmeliste, Anschluss, Kommentar zu `gates-schnell` nicht an `GATE_CHECKS` gedeckt |
| `tools/harness/gates-schnell.sh` | geprüft; F-694, F-695; sonst ohne Befund (Quoting, `set -e`, Reihenfolge Entfernen → Lauf → Hash → Schreiben) |
| `tools/harness/gates-schnell-gegenprobe.sh` | geprüft; F-695, F-696, F-699, F-700; sonst ohne Befund (Arbeitsbaum bleibt unberührt, Attrappen, `setsid`, Auswertung `rot`/`gruen` mit `exakt`) |
| `harness/README.md` (zwei Gate-Zeilen), `harness/conventions.md` (Träger-Zeile) | geprüft, ohne Befund bis auf die Zahl „Punkt 1 bis 7“ (stimmt mit der Gegenprobe, nicht mit der Spezifikation, F-696); Träger-Zeile nennt Fragment und Skript, die Gegenprobe-Datei ist über das Fragment gedeckt |
| `Makefile`, `harness/mk/enforce.mk`, `tools/harness/record-gates.sh`, `tools/harness/working-tree-hash.sh`, `.claude/hooks/` | im Diff unberührt, geprüft, ohne Befund |
| Slice-Plan §1, §3, §4, §8 und Kopf | geprüft, ohne Befund (außer §6, F-697); §7 enthält Belege, keine Plan-Entscheidung |
| Commits | geprüft, ohne Befund: reine Inhalts-Commits, kein Move (§3.3), `SPEC-052` vor dem ersten Code-Commit (§3.12); Messages nennen Slice und [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md) |
| Lint-Profil, Produkt-Code, Tests unter `internal/` und `test/` | im Diff unberührt, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Stand ohne Lauf (Trockenlauf schreibt) · Zusage ohne fangende Mutation · Zusage im Satz weiter als die Prüfung · Aufzählung weicht von der Liste ab · Kommentar weiter als die Prüfung · Überschuss in der Liste nur indirekt gefangen · Muster enger als die Wortlaut-Regel · Beleg ohne Ursache

## Verdikt

**Merge-blockierend:** nein im Sinn von HIGH; die drei MEDIUM (F-694 bis F-696) sind vor der Closure zu schließen (§3.10: Zusage ohne Mutation; F-694 ist ein reales Fehlverhalten der Kette, kein Wortlautfehler).

**Entscheidungen, die der Architect treffen muss:**

1. F-697: Punkt 13 von `SPEC-052` und die Zeile in §6 um den siebten Eintrag `gates-schnell-gegenprobe` nachziehen (oder die Aufzählung durch „jedes Gate der Ausnahmeliste“ ersetzen); damit auch bestätigen, dass der Beleg für `docs/user/` und `README.md` für diese Gegenprobe gilt.
2. DoD Liefer-Punkt 3 („hängt mit der Kette an `GATE_CHECKS`“): Klarstellung, dass die Gegenprobe gemeint ist (Auslegung 1 ist die einzig mögliche, bestätigt).
3. F-696: ob Punkt 8 einen Fall bekommt oder der Satz „Die Punkte 1 bis 8 prüft …“ auf 1 bis 7 gekürzt wird; F-695: ob Punkt 5 „zum Ende des Laufs“ und Punkt 7 zweite Hälfte eine Mutation erhalten oder enger gefasst werden (Zusage nach §3.11 enger fassen oder prüfen, nicht beides offen lassen).
4. F-699: ob Punkt 3 („nie ein Gate“) um den Überschuss in der Liste ergänzt wird.

**Übergabe:** Findings gehen an den Implementer; F-697 und die Entscheidungen oben an den Architect; die Finding-Klassen zusätzlich in die Slice-Closure §7. Gate-Ergebnis: `make gates` nicht gefahren (Auftrag des Aufrufers); `make docs-check` grün (0 Befunde).
