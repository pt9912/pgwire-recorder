# Review-Report: slice-harness-gegenproben-aus-gates-sammelziele — 2026-10-18

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; das ist die Frage des Verifiers).

**Gegenstand:** `slice-harness-gegenproben-aus-gates-sammelziele`, Zweig `rev-sz` (Implementer-Stand 16477376), Diff `09d35ed8..HEAD` über `tools harness AGENTS.md docs`.

**Skill:** `.harness/skills/reviewer.md` @ 229b5efb ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

> **Zitier-Form.** Dieser Report friert ein; Slices werden mit Kennung zitiert. Die ADR- und SPEC-Kennungen sind auf Wunsch des Auftrags auflösende Links.

**Eingangs-Kontext:**

- `slice-harness-gegenproben-aus-gates-sammelziele` (Plan in `in-progress/`, §7 *Belege des Implementers*: Metrik, Läufe, Mutationstabelle, Rückgaben 1 bis 5)
- [ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md) (Accepted), [ADR-0047](../plan/adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md)
- [SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) Teil 2 (Punkte 4 bis 12) und Punkte 13 bis 15, [SPEC-052](../../harness/vertraege/SPEC-052-gate-kette.md) Punkte 6 und 8 (gelesen)
- [LH-QA-07](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur mittelbar (Plan-Kopf)
- `AGENTS.md` (Hard Rules §3.7, §3.9, §3.10, §3.11, §3.12, §3.14; §6 Schritt 6)

**Gefahrene Sensoren (Auftrag des Nutzers: immer der engste):** `bash -n` über beide neuen Skripte, `make kopf-check` (Exit 0), `make docs-check` (0 Befunde), `bash tools/harness/gegenproben-regel.sh` (Exit 0), `make gegenproben-gegenprobe` einmal in der Sprache `de_DE.UTF-8` mit `LANGUAGE=de` (grün, 27 s), `git diff --numstat`. Dazu 14 zusätzliche Mutationen in einer Kopie unter dem Scratchpad (nicht im Repo), gefahren mit `GP_MUT` und `GP_ENTDECKE=1` (11 s), und ein Einzelversuch des Abbruchs mit einem Einzelziel, das ein Kind startet. Keine volle Gegenprobe, kein `make gates`, kein `make gegenproben`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Abbruch (Punkt 8) beendet die Nachkommen des Einzelziels rekursiv mit `pgrep -P`, aber alle vier Fälle fahren eine Attrappe, die ihren Prozess mit `exec sleep` ersetzt und keine Kinder hat. Die Mutation, die `beende` auf den Kindprozess beschränkt, färbt keinen Fall rot (grüner Mutant, nicht im §7-Satz „Grüne Mutanten: keine“ enthalten). Ein Einzelversuch mit einem Einzelziel, das `sleep 31 &` startet, zeigt: mit der Rekursion ist das Kind nach `SIGTERM` an das Skript beendet, ohne sie lebt es weiter. Das echte Einzelziel (`make <ziel>` mit Rezept und Docker-Aufrufen) hat Nachkommen. | `AGENTS.md` §3.10 (Zusage „beendet das laufende Einzelziel“, Punkt 8) | `tools/harness/gegenproben.sh` · `for k in $(pgrep -P "$1" 2> /dev/null); do beende "$k" "$2"; done;` | ja — Mutation „Schleife über `pgrep` gestrichen“ in `gegenproben-gegenprobe.sh`; zurzeit grün | Zusage mit ungeprüftem Zweig: Nachkommen des Einzelziels beim Abbruch |
| F-2 | MEDIUM | Punkt 11 zählt „grün, rot, abgebrochen oder unter `-n`, `-q`, `-t`“ auf; der Fall `p11-kein-stand` fährt grün, rot und `-n`. Für „abgebrochen“, `-q` und `-t` gibt es keinen Fall, und die einzige Mutation (Datei nach grünem Lauf) bricht nur ein Element. Eine Mutation, die beim Abbruch oder nur unter `-q` eine Datei unter `.harness/state/` schreibt, bleibt grün; eine, die bei rotem Lauf schreibt, wird zwar rot, hat aber keine eigene Mutation. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `tools/harness/gegenproben-gegenprobe.sh` · `mut m-p11-schreibt` und `# Punkt 11` | ja — Mutationen „Datei im Zweig `abbruch`“ und „Datei im Zweig `*q*`“ grün | Aufzählung in der Zusage: Fall und Mutation nur für einen Teil der Elemente |
| F-3 | MEDIUM | Punkt 8 zählt zwei Signale und zwei Adressaten auf, die Mutation „Exit-Codes vertauscht“ bricht beide Signale zugleich, und es gibt keine Mutation, die nur ein Signal oder nur einen Adressaten bricht. Die Gegenprobe kann sie unterscheiden: „INT nach 143“ färbt nur die beiden `sigint`-Fälle, „TERM nach 130“ nur die `sigterm`-Fälle, „Trap für INT gestrichen“ nur `p8-sigint-skript`, „Trap für TERM gestrichen“ nur `p8-sigterm-skript` (alle vier in der Kopie gefahren, Ergebnis wie genannt). | `AGENTS.md` §3.10 (Aufzählung zählt je Element) · [SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) Kopf („Je Zusage und je Element einer Aufzählung … eine Mutation, die nur ihn rot färbt“) | `tools/harness/gegenproben-gegenprobe.sh` · `mut m-p8-vertauscht` | ja — die vier Mutationen sind in der Kopie nachvollziehbar | Aufzählung in der Zusage: Mutation bricht mehrere Elemente zugleich |
| F-4 | LOW | Das Gate-Index-Zeile von `make gegenproben-gegenprobe` sagt, die Flags `-n`, `-q`, `-t`, `-i`, `-k`, `-j` würden „einzeln, aus der Umgebung und kombiniert“ geprüft. `-k` und `-j` haben nur Fälle auf der Kommandozeile (`p10-k`, `p10-j`; mit `-i` kombiniert `p10-i-k`, `p10-i-j`), `-n`, `-q`, `-t` sind nicht untereinander kombiniert. Dieselbe Zeile und der Kommentar in `gegenproben-gegenprobe.mk` nennen als Voraussetzung „bash, make, setsid, pgrep“; der Kopf des Skripts nennt zusätzlich `timeout` und „bash ≥ 5“, der Vertrag (Plattform) „bash ≥ 4“ und nur `setsid`. | `AGENTS.md` §3.11 (ein Satz sagt nur zu, was ein Test prüft) | `harness/README.md` · `die Flags `-n`, `-q`, `-t`, `-i`, `-k`, `-j` einzeln, aus der Umgebung und kombiniert`; `tools/harness/gegenproben-gegenprobe.sh` · `Braucht bash ≥ 5, make, setsid, pgrep, timeout.` | nein — Prosa gegen Fallliste, kein Gate-Lauf | Zusage im Gate-Index weiter als die Fälle |
| F-5 | LOW | Punkt 4 und der Kommentar in `gegenproben.mk` sagen, die Sammelziele stünden nicht in der Ausnahmeliste und `make gates-schnell` fahre keines; die Fälle (`p4-kein-gate`: `make -n gates` und `GATE_CHECKS`) prüfen beides nicht, und der Mutationssatz setzt nur `gegenproben` (nicht `gegenproben-host`) in `GATE_CHECKS`. Ein Eintrag in `GATE_SCHNELL_AUSNAHME` bleibt hier grün (in der Kopie gefahren); das Ziel ohne Eintrag in `GATE_CHECKS` fängt nach dem Vertrag eine andere Gegenprobe ([SPEC-052](../../harness/vertraege/SPEC-052-gate-kette.md) Punkt 4), in dieser steht dafür kein Fall. Die Mutation für `gegenproben-host` in `GATE_CHECKS` färbt `p4-kein-gate` rot (gefahren), hat aber keine eigene Zeile. | `AGENTS.md` §3.11 und §3.10 | `harness/mk/gegenproben.mk` · `beide stehen weder in GATE_CHECKS noch in GATE_SCHNELL_AUSNAHME` | ja — Mutation `GATE_SCHNELL_AUSNAHME += gegenproben` grün | Zusage im Kommentar ohne eigenen Fall |
| F-6 | INFO | Der Wortlaut von `AGENTS.md` §6 Schritt 6 entspricht nach Entfernen der Backticks Zeichen für Zeichen dem Satz des Nutzers (`grep -F` ergibt 1 Treffer); geändert ist genau diese eine Zeile (ein Einfügen, ein Löschen), der Satz vor dem Semikolon stammt aus 8c649c4e. In `harness/README.md` §Sensors steht derselbe zweite Halbsatz zusätzlich; das ist eine Wiederholung der Regel an zweiter Stelle, keine Erweiterung. | `AGENTS.md` §6 (Entscheidung des Nutzers vom 2026-10-18) | `AGENTS.md` · `am Ende jeder Welle und vor dem Release fahren` | ja — `grep -F` | — |
| F-7 | INFO | Metrik: `git diff --numstat 09d35ed8..HEAD -- tools harness AGENTS.md` ergibt 426 Einfügungen und 2 Löschungen, Summe 428, wie in §7 wörtlich genannt; der Letzte Commit (16477376) berührt nur den Plan. Unter Meldeschwelle 700 und Marke 900. Fälle (37) und Mutationen (34) stimmen mit der Zählung in `harness/README.md`, dem Fußsatz der Gegenprobe und §7 überein. | `AGENTS.md` §3.14 | `docs/plan/planning/in-progress/slice-harness-gegenproben-aus-gates-sammelziele.md` · `Summe **428**` | ja | — |
| F-8 | INFO | Locale: Weder `gegenproben.sh` (liest nur `MAKEFLAGS`-Buchstaben und Exit-Codes) noch die Gegenprobe (filtert Zeilen mit `^make`, vergleicht Exit-Codes und eigene Ausgaben, sortiert mit `LC_ALL=C`) liest Text von `make`; `make gegenproben-gegenprobe` unter `de_DE.UTF-8` und `LANGUAGE=de` ist grün. Es gibt dafür keinen Fall (der Fall `p2-locale` gehört der Regel aus Teil 1). | Maintainability (Fund „Ausgabe von make ist locale-abhängig“) | `tools/harness/gegenproben-gegenprobe.sh` · `grep -v '^make' "$st/err.raw"` | ja — Lauf in der Locale | — |
| F-9 | INFO | Rückgaben des Implementers (gehen an den Architect; hier nur die Dokumentation): (1) Der Gate-Index sagt für `make gegenproben` „Exit-Code ungleich 0“ und für die Gegenprobe „über `make` nur Exit ungleich 0“, das stimmt mit dem Code überein; nicht dokumentiert ist, dass `make -i gegenproben` mit Exit 0 endet, obwohl das Skript ablehnt (Exit 2 wird von `-i` überlesen), kein Einzelziel läuft und nur die Meldung auf stderr bleibt. (2) Die Zeile `rot (Exit n)` trägt den Exit von `make <ziel>` (2), nicht 127; der Fall `p6-127` fährt 127 nur an der Attrappe. (3) `pgrep` und `timeout` stehen nicht im Vertrag (siehe F-4). (4) Der Standardeingang der Einzelziele ist `/dev/null`; dokumentiert nirgends, und `x-stdin` (Eingabe statt `/dev/null`) bleibt grün. (5) Die Mutation „Prüfung liest die Liste des anderen“ ist als vertauschte Rezepte umgesetzt und färbt neun Fälle statt nur die von Punkt 9. Plan §1, §3 und §6 stimmen mit dem unveränderten Vertrag überein (Exit 1 und 2 stehen dort als Vertragsaussage), §3.9 ist damit bis zur Entscheidung des Architect nicht verletzt; folgt der Vertrag den Rückgaben, ziehen G3, G4 und G7 in §6 nach. | `AGENTS.md` §3.9, §3.12 | `docs/plan/planning/in-progress/slice-harness-gegenproben-aus-gates-sammelziele.md` · `Rückgaben an den Architect` | nein — Entscheidung des Architect | Rückgabe an den Architect, Vertrag und Code weichen im Exit-Code über `make` ab |
| F-10 | INFO | Der Lauf von `make gegenproben` im Gegenproben-Container folgte auf `make gates` im selben Container (§7 nennt das selbst); das Risiko 1 aus §6 (grün ohne vorheriges `make gates`) und die Regel an die Lieferung in [SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 3 sind damit nicht belegt. Der Ausgang liegt beim Verifier und Orchestrator, nicht bei dieser Prüfung. | `AGENTS.md` §3.12 · Plan §6 Risiken | `docs/plan/planning/in-progress/slice-harness-gegenproben-aus-gates-sammelziele.md` · `kein frischer Container` | nein — Fahrt im frischen Container | — |
| F-11 | INFO | Der Zweig `trap - INT TERM` im Hintergrund-Unterprozess des Einzelziels (Rückgabe 3 „zurückgesetztes `SIGINT`/`SIGTERM`“) hat keinen Fall: entfernt bleibt die Gegenprobe grün. Ob die Zeile etwas bewirkt, ist offen; weder Kommentar noch §7 führen eine Invariante dafür an. | `AGENTS.md` §3.10 | `tools/harness/gegenproben.sh` · `( trap - INT TERM; exec "$make_bin" "$z" ) &` | ja — Mutation „`trap - INT TERM` gestrichen“ grün | Zusage mit ungeprüftem Zweig: Zeile ohne Fall |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `AGENTS.md` §6 Schritt 6 (nur die erlaubte Änderung, beide Hälften = Wortlaut des Nutzers) | geprüft, ohne Befund (F-6) |
| Fälle und Mutationen gegen Punkte 4 bis 12 von [SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) (37 Fälle, 34 Mutationen, Namen und Mengen wie im Vertrag; Aufrufarten `-n -q -t -i` je Form, `MAKEFLAGS`, `--ignore-errors`, Reihenfolge der Gründe) | geprüft; Befunde nur zu Punkt 8 (F-1, F-3) und 11 (F-2), sonst ohne Befund |
| Kommentare in `harness/mk/gegenproben.mk`, `gegenproben-gegenprobe.mk`, `tools/harness/gegenproben.sh` (§3.7: Klassen, kein Konjunktiv, keine Chronik) | geprüft; Zusagen siehe F-4, F-5, sonst ohne Befund |
| `harness/README.md` §Sensors: Zeilen der Sammelziele und der Gegenprobe, Werkzeug-Tabelle | geprüft; Befund F-4 (Prosa), sonst ohne Befund |
| Plan (§1, §3, §6, §8) gegen Diff, Kopf (`kopf-check` Exit 0) | geprüft, ohne Befund (F-9 zur Rückgabe) |
| Metrik (428, Marke 900) | geprüft, ohne Befund (F-7) |
| Sprache von `make` in Sammelskript und Gegenprobe | geprüft, ohne Befund (F-8) |
| `docs/` (`make docs-check`, 0 Befunde) | geprüft, ohne Befund |
| `internal/`, `cmd/`, `test/`, `spec/` | nicht berührt vom Diff |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 6 |

**Finding-Klassen dieses Laufs:** Zusage mit ungeprüftem Zweig · Aufzählung in der Zusage: Fall und Mutation nicht je Element · Zusage im Gate-Index weiter als die Fälle

## Verdikt

**Merge-blockierend:** ja — F-1 bis F-3 sind MEDIUM (Mutationen, die nach §3.10 fehlen; jede ist in der Gegenprobe mit einer Zeile `mut` und ohne Eingriff in das Skript zu ergänzen). Der Code des Skripts selbst zeigte in keinem Versuch ein falsches Verhalten; offen ist die Prüfung seiner Zusagen. Die Rückgaben 1 bis 5 gehen an den Architect, nicht an den Implementer.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). Dieser Report ersetzt keine Verifikation.
