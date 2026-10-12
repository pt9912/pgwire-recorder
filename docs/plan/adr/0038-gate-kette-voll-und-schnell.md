# ADR-0038: Gate-Kette in voller und schneller Kette, Wahl am geänderten Pfad

**Status:** Proposed

**Datum:** 2026-10-12

**Autor:** Architect (Empfehlung; Entscheidung beim Nutzer)

**Bezug:** [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0034](0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md), [ADR-0001](0001-hexagonale-architektur.md)

**Schärft:** — (Prozess-ADR; die Spezifikation der Kettenwahl entsteht mit dem Folge-Slice als eigene Stelle im Technik-Stratum, `SPEC-049` bleibt unverändert, siehe Konsequenzen)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Stop-Hook `.claude/hooks/stop-require-gates.sh` gibt einen Lauf nur frei, wenn der Arbeitsbaum-Hash (`tools/harness/working-tree-hash.sh`) dem Hash des letzten grünen `make gates` entspricht (`tools/harness/record-gates.sh`). Jede Änderung, auch ein reiner Doku-Commit (Plan, Report, Handbuch), verlangt damit den vollen Lauf.

Gemessen am 2026-10-12, je Gate einzeln: `lint-gegenprobe` 150 s (Load im Mittel 57, Maximum 75 auf 20 Kernen), `test-integration` 96 s, `a-check-negativ` 11 s, alle übrigen höchstens 2 s; `make gates` insgesamt etwa 4,5 Minuten. Die Gegenproben prüfen die Gates, nicht das Produkt.

Zwei Befunde aus dem Bestand bestimmen die Wahl:

1. **Die Eingabe einer Gegenprobe ist der ganze Arbeitsbaum.** `lint-gegenprobe` fährt ihre Mutanten in Kopien des Arbeitsbaums und verlangt, dass die unveränderte Kopie grün ist; sie prüft damit auch den aktuellen Go-Code mit dem aktuellen Profil. Eine Liste „Eingaben der Gegenprobe“ (Profil, Skripte, Dockerfile-Stufe, Fragmente) wäre unvollständig, sobald ein Go-Quelltext ein Lint-Befund wird. Dasselbe gilt für `a-check-negativ`.
2. **Eine falsche Eingabe-Erkennung ist ein stilles Grün.** Wer eine Gegenprobe überspringt, obwohl sich eine ihrer Eingaben geändert hat, bekommt einen Nachweis „alle Gates grün“, der nicht stimmt; nichts im Lauf zeigt es an (`grundlagen-durchsetzungsschicht.md` §Vier Design-Eigenschaften: fail-closed, Nachweis über Inhalt). Der Fehler läge in Code, den kein Gate prüft.

`make gates` ist heute die Zusage „voller Nachweis“. `AGENTS.md` §3.6 gilt für Schwellen; eine Bedingung an einem Gate ändert diese Zusage und gehört deshalb in eine ADR, obwohl keine Schwelle sinkt.

## Entscheidung (Empfehlung)

Wir wählen **Option B, eng gefasst: zwei Ketten, die Wahl trifft eine Positivliste mit Rückfall auf die volle Kette.**

1. `make gates` bleibt unverändert die volle Kette und der einzige Nachweis für Code, Gates, Profile, Skripte, `Dockerfile`, Fragmente, `go.mod` und alles, was nicht ausdrücklich freigegeben ist.
2. Es kommt eine **schnelle Kette** hinzu, die nur die Gates ohne Gegenproben und ohne Integrationstests fährt. Sie ersetzt den vollen Nachweis nie, sie ersetzt nur den Lauf für Änderungen, die keine Eingabe eines langsamen Gates sind.
3. Der Stop-Hook wählt nicht über eine Liste der Eingaben je Gegenprobe, sondern **umgekehrt über eine Positivliste inerter Pfade** (Planung, Reports, Review-Berichte; weitere Pfade nur, wenn belegt ist, dass kein langsames Gate sie liest). Verglichen wird der Arbeitsbaum mit dem Stand des letzten **vollen** grünen Laufs; liegt jede geänderte, hinzugekommene oder gelöschte Datei auf der Positivliste und ist der Stand der schnellen Kette grün, gibt der Hook frei. Jeder andere Fall, ein fehlender Stand des vollen Laufs, ein nicht lesbarer Vergleich, ein unbekannter Pfad, verlangt `make gates` (fail-closed).
4. Die Positivliste, die Form des Nachweises (Liste der Dateihashes des vollen Laufs statt eines einzelnen Hashs) und die Randformen (Umbenennung, Löschung, Symlink, untracked Datei, Pfad auf der Grenze) stehen in der Spezifikation, nicht hier (`AGENTS.md` §3.8). Gegenproben des Hooks sind Teil der Lieferung (`AGENTS.md` §3.10).

Der Nutzer entscheidet, ob er diese Empfehlung annimmt und welche Pfade zur Positivliste gehören; bis dahin ist die ADR ein Vorschlag.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Gegenproben nur fahren, wenn sich ihre Eingaben ändern | spart den Lauf auch bei Änderungen an Docs und unbeteiligtem Code; `make gates` bleibt ein Aufruf | Der Eingabesatz einer Gegenprobe ist der ganze Arbeitsbaum (Befund 1); eine Eingabeliste pflegt dieselbe Hand, die die Gegenprobe ändert, und ein Fehler ist ein stilles Grün (Befund 2); die Erkennung bräuchte selbst eine Gegenprobe, die den Fehler der Liste nicht findet, wenn sie dieselbe Liste benutzt |
| **B — zwei Ketten, Positivliste inerter Pfade, Rückfall auf voll** | Fehler der Liste kosten Zeit, kein Grün; die Positivliste ist klein und je Eintrag belegbar; der volle Nachweis für Code bleibt wörtlich `make gates` | Der Nachweis wird eine Dateiliste statt eines Hashs; der Hook wird komplexer; die Schnellkette ist ein zweites Ziel, das jemand pflegt; Doku-Commits sparen die Zeit, Code-Commits nicht |
| C — unverändert lassen | keine neue Mechanik, kein neues Fehlerbild | 4,5 Minuten auch für einen Satz im Plan; der Druck, Läufe zu bündeln, kostet Rückmeldung; die Last (Load 57 auf 20 Kernen) hält andere Läufe auf demselben Host auf |
| D — Gegenproben aus `make gates` herausnehmen, nur in einem eigenen Ziel | schnellster Standardlauf | Lockerung von `make gates` (`AGENTS.md` §3.6); ein Gate, das niemand fährt, ist eine Behauptung (`modul-13-quality-gates.md`, Harness-Lüge) |

Zu C bleibt festzuhalten: Die Messung ist eine Last-, keine Korrektheitsfrage. Wäre `lint-gegenprobe` bei geringerer Parallelität schneller (offene Folge des Plans zur Gegenprobe-Parallelität), schrumpft der Nutzen von B; B lohnt sich auch dann nur für Änderungen, die gar keine Eingabe sind.

## Konsequenzen

- Positiv: Doku- und Planungs-Commits brauchen Sekunden statt Minuten; Code, Gates und Konfiguration behalten den vollen Nachweis; ein Fehler in der Pfadwahl führt zum vollen Lauf, nicht zu einem Grün.
- Negativ: `tools/harness/record-gates.sh`, `tools/harness/working-tree-hash.sh` und der Hook ändern sich; der Hook ist die Stelle, die den Nachweis prüft, ein Fehler dort bleibt ungeprüft, solange keine Gegenprobe ihn fängt. Das ist das Restrisiko dieser Option und der Grund für Punkt 4 der Entscheidung.
- Akzeptiertes Negativ: Ein Commit, der nur inerte Pfade ändert, aber eine Datei einführt, die später ein Gate liest, wird beim ersten Lauf der vollen Kette gefunden, nicht beim Doku-Commit. Der Grund: Die Positivliste nennt nur Pfade, die kein langsames Gate liest; ein späterer Leser ändert Code oder Gates und löst damit den vollen Lauf aus.
- Folgepflicht: ein Folge-Slice (Kennung vergibt der Planner nach der Entscheidung des Nutzers) liefert die Spezifikations-Stelle, die schnelle Kette, die Nachweisform, die Hook-Änderung samt Gegenprobe und trägt die Gate-Zeilen in `harness/README.md` §Sensors nach. Die Kennung wird hier nicht erfunden. `harness/conventions.md` §Welche Regelblöcke des Baums hier einen Träger haben nennt `grundlagen-durchsetzungsschicht.md` und seine Träger; sie ziehen mit.

## Fitness Function (falls maschinell prüfbar)

Entsteht mit dem Folge-Slice: eine Gegenprobe des Hooks, die je Pfadklasse (inert, nicht inert, unbekannt, gelöscht, umbenannt) und bei fehlendem oder unlesbarem Stand des vollen Laufs prüft, dass nur der inerte Fall freigibt.

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, git | Positivliste und Rückfall auf die volle Kette | wird mit dem Folge-Slice benannt |

## Re-Evaluierungs-Trigger

Wenn die volle Kette unter eine Minute fällt; wenn die Positivliste mehr Fehlentscheidungen als Fälle zählt (ein inerter Pfad wird von einem Gate gelesen); wenn eine CI die volle Kette je Änderung fährt und der Hook nur noch lokal schützt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-12 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
