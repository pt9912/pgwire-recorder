# Review-Report: slice-harness-devcontainer-einrichten-geruest — 2026-10-18

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-harness-devcontainer-einrichten-geruest`, Diff `8b19dc66..4f72a7a5` (Pfade `tools`, `harness`, `docs`; Zweig `rev-geruest`).

**Skill:** `.harness/skills/reviewer.md` @ 229b5efb · <!-- d-check:ignore (Skill-Pfad des Repos, Zitat des Laufs) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-18

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-harness-devcontainer-einrichten-geruest` (§1, §3, §4, §6, §7 Belege des Implementers, §8); Nehmer `slice-harness-devcontainer-einrichten-start` (§1, DoD) für die Adresse des grünen Mutanten
- [ADR-0045](../plan/adr/0045-container-start-als-werkzeug-im-repo.md) (Accepted); [SPEC-058](../../harness/vertraege/SPEC-058-container-start.md) Teil I (Punkte 2 bis 8, 13 bis 17 soweit Teil I)
- [LH-QA-03](../../spec/lastenheft.md#lh-qa-03--portabilität) (nur Host-Voraussetzungen)
- `AGENTS.md` (Hard Rules, §3.7, §3.9, §3.10, §3.11, §3.12, §3.13, §3.14)
- Eigene Läufe im Worktree: `make docs-check` (0 Befunde), `make kopf-check` (Exit 0); Sonden am Skript `tools/harness/container-start.sh` mit Eingaben aus der Umgebung (`PGR_SSH_PORT`, `PGR_WARTEN`, `PGR_DEV`); `make gates` und `make container-start-gegenprobe` nicht gefahren (Auftrag: Verifier).

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill
— Tag einsetzen, denn diese Zeile wandert in den eingefrorenen Report.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Funktion `zahl` lässt `+` und `-` in jeder Stelle außer der ersten zu: `PGR_SSH_PORT=2+2`, `22-` und `1-1` werden angenommen (Sonde: keine Meldung zu `PGR_SSH_PORT` in `eingaben`), ebenso `PGR_WARTEN`. [SPEC-058](../../harness/vertraege/SPEC-058-container-start.md) Punkt 2 verlangt Dezimalziffern; der Kommentar über `zahl` sagt „Dezimalziffern, kein Vorzeichen“. Ein solcher Port landet als `127.0.0.1:2+2:2222` in `override.json`. Kein Fall und keine Mutation berührt ein Zeichen `+`/`-` hinter der ersten Stelle (`p2-zahl-ziffern` ersetzt die Zeile mit der Klasse `[!0-9+-]` durch `:` und wird nur von `abc` gefangen). | SPEC-058 Punkt 2; `AGENTS.md` §3.10 | `tools/harness/container-start.sh` · `[[ $1 != *[!0-9+-]* ]] || return 1` | ja — `PGR_SSH_PORT=2+2 PGR_DEV=… bash tools/harness/container-start.sh` endet nicht mit Exit 2; ein Fall `p2-port-…` mit `2+2` würde es bestätigen | Eingabeform lässt Zeichen zu, die der Vertrag nicht nennt (Fall fehlt) |
| F-2 | MEDIUM | [SPEC-058](../../harness/vertraege/SPEC-058-container-start.md) Punkt 15 (e) gilt nach der Zuordnung des Vertrags schon in Teil I („für Verzeichnisse“): Die Gegenprobe bricht ab, wenn ein Name aus der Umgebung das Präfix nicht trägt. Der Fall fehlt; der Mutant `p15-schutz` (Vergleich auf das Präfix vor `rm -rf` entfällt) bleibt grün und ist in §7 als Absicherung geführt. Die Invariante („`$arbeit` stammt aus `mktemp -d -t pgr-cs-XXXXXX`“) und der Ort stehen dort, eine Sonde mit Lauf fehlt (Begründung durch Lesen, kein Aufruf mit einer Eingabe, an die der Autor nicht dachte); der Plan nennt zudem in DoD und §6 „kein Name ohne Präfix `pgr-cs-`“ als Zusage der Gegenprobe. | `AGENTS.md` §3.10; SPEC-058 Punkt 15 | `docs/plan/planning/in-progress/slice-harness-devcontainer-einrichten-geruest.md` · „`p15-schutz` (der Vergleich auf das Präfix vor `rm -rf` entfällt)“ | nein — nur ein Fall in der Gegenprobe oder die Entscheidung des Architect (Punkt 15 (e) erst mit Teil II) würde es bestätigen | Äquivalenter Mutant ohne Sonde mit Lauf; Zusage des Vertrags ohne Fall |
| F-3 | MEDIUM | Der grüne Mutant `p2-vorgabe-warten` nennt als Adresse `slice-harness-devcontainer-einrichten-start`; dessen §1 und DoD führen die Vorgabe `60` von `PGR_WARTEN` als Fall nicht, nur „Eingaben mit Vorgaben und Formen“ als vorausgesetzt (Bindung an den Geber) und die Wartezeit der Probe (A8). Der Implementer vermerkt selbst „Punkt an den Planner“; im Diff trägt der Nehmer die Sendung nicht mit der Kennung des Gebers, `grep` nach `p2-vorgabe-warten` findet sie dort nicht. | `AGENTS.md` §3.13 (Adresse nimmt an) | `docs/plan/planning/in-progress/slice-harness-devcontainer-einrichten-geruest.md` · „`p2-vorgabe-warten` (Vorgabe `61` statt `60`)“ gegen `docs/plan/planning/open/slice-harness-devcontainer-einrichten-start.md` §1 | nein — Lesen von §1 und DoD des Nehmers | Adresse nimmt nicht an |
| F-4 | LOW | Der Fall `p2-warten-abc` übergibt den Wert `ab` statt `abc` (`for p in leer: abc:ab …`); [SPEC-058](../../harness/vertraege/SPEC-058-container-start.md) Punkt 2 nennt `abc` für `PGR_WARTEN`. Die Prüfung ist dieselbe Funktion, die Abweichung ohne semantische Folge. | SPEC-058 Punkt 2 | `tools/harness/container-start-gegenprobe.sh` · `for p in leer: abc:ab 0:0 601:601` | ja — Lesen der Zeile | Fallwert weicht vom Vertrag ab |
| F-5 | INFO | Ist `$PGR_DEV/ssh` eine Datei (kein Verzeichnis, kein Paar), meldet `zustand` nichts („erzeugen“); `schluessel` scheitert erst nach `geruest` mit Exit 5, obwohl Punkt 4 bis `zustand` ohne Änderung zusagt, was ein früherer Fehler wäre. Die Randform steht weder in §6 noch in [SPEC-058](../../harness/vertraege/SPEC-058-container-start.md); Entscheidung beim Architect, kein Rückgabepunkt R1 bis R6. | `AGENTS.md` §3.12 | `tools/harness/container-start.sh` · `zs=erzeugen` | nein — ein Fall mit `ssh` als Datei würde es zeigen | Randform von Teil I nicht genannt |
| F-6 | INFO | R1 bis R6 sind in §7 mit dem heutigen Verhalten dokumentiert (R1 Gegenprobe bricht ohne `u-boot` ab, R2 Modus eines vorhandenen `ssh/`, R3 unlesbare Ausgabe von `docker version`/`u-boot --version`, R4 stdout bei Ausgabe von `u-boot`, R5 `override.json` im Vergleich „Gerüst vorhanden“, R6 Mutation „Kollision“ in Teil I nicht möglich). Der Code tut das Genannte (R3: „älter als 24“ beziehungsweise „unlesbar“; R2: `mkdir -p -m 700`); R1 trifft `make gates` auf Hosts ohne `u-boot`, `node`, `ssh-keygen`, wie das Fragment `container-start-gegenprobe.mk` es sagt. Gewertet wird nur die Dokumentation, nicht die Entscheidung. | Maintainability | `docs/plan/planning/in-progress/slice-harness-devcontainer-einrichten-geruest.md` · „Offene Punkte, an den Architect zurückgegeben“ | nein | Rückgabepunkte dokumentiert |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/container-start.sh`, Zugangsvariablen (`CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY`, `GIT_TOKEN`) | geprüft, ohne Befund; `unset` zu Beginn vor jedem Kind, das Skript liest und gibt keinen Wert aus; der Verweis `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}` bleibt Verweis; Fälle `p8-umgebung-*` und `p8-keine-werte` mit je einer Mutation je Variable laut Tabelle |
| `tools/harness/container-start.sh`, Schreibbereich (nur unter `$PGR_DEV`, kein `git`, kein `~/.ssh`, keine Zerstörung) | geprüft, ohne Befund; Schritte `eingaben`, `werkzeuge`, `zustand` schreiben nichts, `geruest`/`schluessel`/`override` arbeiten im Verzeichnis; Fälle `p1-*` und `p4-nichts-geaendert` |
| `harness/mk/container-start.mk` (Aufrufarten von `make`) | geprüft, ohne Befund; Rezept eine Zeile ohne `$(MAKE)` und ohne `+`, `.PHONY` ohne Voraussetzung; Fälle `p14-*` zu `-n`, `-q`, `-t`, `-i`, `-k`, `-j`, `MAKEFLAGS` und `-n -i` mit Mutationen (`p14-plus`, `p14-make`, `p14-phony`) |
| Punkte 3, 5, 6, 7, 8, 13 gegen die Tabelle in §7 (Fälle je Form und je Element, Mutation, die nur ihn färbt) | geprüft, ohne weiteren Befund; die Aufzählungen von Punkt 2 (13 Formen von `PGR_DEV`), Punkt 3 (sechs Werkzeuge, Versionsgrenzen) und Punkt 7 (acht Aufrufe) haben je einen Fall und eine benannte Mutation (Stichprobe der Fall- und Mutationszeilen in `tools/harness/container-start-gegenprobe.sh`); `make container-start-gegenprobe` selbst nicht gefahren (Verifier) |
| Grüne Mutanten `p2-vorgabe-warten`, `p15-int`, `p15-term` | geprüft; Invariante, Ort und Sonde stehen; Sonde von `p2-vorgabe-warten` (Zeilen 6, 36, 59, 60, 63) stimmt mit `grep -n PGR_WARTEN` überein, `PGR_WARTEN` wird im Skript nur geprüft |
| Kommentare in `container-start.sh`, `container-start-gegenprobe.sh`, den zwei Fragmenten (`AGENTS.md` §3.7) | geprüft, ohne Befund; Zusage, Kopplung, Abgrenzung oder Grenze, keine verworfene Alternative, kein abwesender Text (Ausnahme: der Kommentar über `zahl`, siehe F-1) |
| Größe (Marke 900, §3.14) | geprüft, ohne Befund; 734 ohne Vertragsdateien (1017 mit) laut §7, Aufteilung 512 + 197 + 19 + 7 ergibt 735 (Rundung), Meldeschwelle 700 überschritten, Marke nicht |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh`, `harness/README.md`, `harness/vertraege/README.md`, `docs/plan/adr/` ([ADR-0045](../plan/adr/0045-container-start-als-werkzeug-im-repo.md) und Index) | geprüft, ohne Befund; Ausnahmeliste und Fixture stimmen überein, Index nennt Ziel unter Werkzeugen und Gegenprobe unter Sensors; `make docs-check` 0 Befunde |
| Plan (§1, §3, §6 folgen dem Diff, §3.9); Kopf (`make kopf-check`) | geprüft, ohne Befund; §3 nennt die Fixture-Zeile, §7 Läufe und Zählung stimmen mit dem Diff; `make kopf-check` Exit 0 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Eingabeform lässt Zeichen zu, die der Vertrag nicht nennt (Fall fehlt) · Äquivalenter Mutant ohne Sonde mit Lauf; Zusage des Vertrags ohne Fall · Adresse nimmt nicht an · Fallwert weicht vom Vertrag ab · Randform von Teil I nicht genannt

## Verdikt

**Merge-blockierend:** ja — drei MEDIUM (F-1 Korrektheitsfehler der Eingabeprüfung samt fehlendem Fall; F-2 Zusage (e) des Vertrags ohne Fall und Mutant ohne Sonde; F-3 Adresse nimmt nicht an). Kein HIGH.

**Übergabe:** Findings gehen an den Implementer (F-3 und F-5 berühren Planner beziehungsweise Architect); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt). Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
