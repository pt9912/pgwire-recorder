# Verifikationsbericht: slice-harness-stop-hook-positivliste-developer — 2026-10-17

**Rolle:** Verifier (Modul 11) · **Frage:** Bauen wir es richtig, gegen Plan, DoD und Vertrag?
**Gegenstand:** Diff `97fc3339..HEAD` (Code-Stand `6d58b508`, danach nur Beleg in §7 und Review-Report); Review `docs/reviews/2026-10-17-review-slice-harness-stop-hook-positivliste-developer.md` (0 HIGH, 0 MEDIUM, 0 LOW, 2 INFO), erst nach dem eigenen Urteil gelesen.
**Eingangs-Kontext:** Slice-Plan (§1 bis §8, Abweichung `p15-ziel-verzeichnis-developer`), `SPEC-052` Punkte 12, 13, 15, `AGENTS.md` §3.9 bis §3.14, [ADR-0038](../plan/adr/0038-gate-kette-voll-und-schnell.md), Register `BEO-REPO/stop-hook-positivliste-ohne-docs-developer`. `closure-note-reviewer` nicht einschlägig (Closure steht aus).

## Eigene Läufe

| Lauf | Ergebnis |
|---|---|
| `make stop-hook-gegenprobe` | grün, 150 Fälle (Beleg nennt 150, vorher 139: 11 neue Fälle laut `FAELLE`) |
| `make docs-check` | 976 Dateien, 0 Befunde |
| `make kopf-check` | grün |
| Marke §3.14 (`git diff --numstat 97fc3339..HEAD -- .claude/hooks tools/harness harness`) | 63 + 26 = **89**, wie in §7; Schwelle 300 nicht erreicht |
| Hook am HEAD im echten Arbeitsbaum | `approve` (Nachweis und Stand der schnellen Kette passen) |
| `make gates` (voll) | von mir **nicht** gefahren; §7 nennt Exit 0 im ersten Lauf, der Stop-Hook gibt am Endstand frei (Dateiliste-Nachweis 10:40 deckt den Code-Stand, Rest liegt auf der Positivliste). Beleg des Implementers, nicht nachgefahren |

## Hook gegen das echte Verhalten (Temp-Repo, echter Hook, echte Dateiliste und Hash)

Aufbau wie in der Gegenprobe (`record-gates.sh`, Stand der schnellen Kette aus `working-tree-hash.sh`), je Fall nur eine Änderung:

| Änderung | Antwort |
|---|---|
| `docs/developer/d.md` (geändert) | `approve` |
| `docs/developer/x.md` (neu, untracked) | `approve` |
| `docs/developer-alt/x.md` | `block` |
| `docs/Developer/x.md` | `block` |
| `internal.go` (Kontrolle) | `block` |

## Eigene Mutationen (Kopie per `git clone` in den Scratchpad, nie im Arbeitsbaum)

| Mutation am Hook | Ergebnis der Gegenprobe |
|---|---|
| `docs/developer/*` → `docs/developer*` (Grenze) | rot: Kopie `ok` nicht grün mit `p12-grenze-developer-alt`, `p12-grenze-developers`, `p12-grenze-datei-developer`, `p15-ziel-verzeichnis-developer` |
| zusätzliche Zeile `docs/Developer/*) return 0` | rot: Kopie `ok` nicht grün mit `p12-case-developer` |
| `docs/developer` aus `spiegelt` entfernt | rot (`m-spiegel-ohne-developer` ändert nichts, `m-spiegel-ohne-verzeichnis` färbt nicht wie erwartet; die Gegenprobe meldet ROT) |
| Zeile `docs/developer/*` aus `inert` entfernt | rot (Gegenprobe ROT) |

Vier von vier Mutationen werden erkannt. Bei den letzten beiden habe ich nur das Schluss-Ende der Ausgabe gelesen; erster roter Fall nicht einzeln zugeordnet (das Review nennt für dieselben Mutationen `p15-spiegel-developer` bzw. `p12-developer`, `p12-developer-tief`).

## DoD-Konformität

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| Hook und Gegenprobe: `inert`, `spiegelt`, Kopf-Kommentar nennen `docs/developer/`; Fälle und Mutationen aus §3; Gegenprobe grün | erfüllt | Diff des Hooks (3 Stellen); alle 11 Fälle und 7 Mutationen aus §3 im Diff (`m-ohne-developer`, `m-developer-erste-ebene`, drei `aufnahme`, `m-spiegel-ohne-developer`, erweiterte `m-verzeichnis-inert`, `m-case-egal`); bestehende Mutationen um die neuen Fälle ergänzt (`m-art-x-wie-datei`, `m-link-ohne-ziel`, `m-spiegel-ohne-verzeichnis`, `SPIEGEL_ALLE`) |
| Vertrag: `SPEC-052` Punkte 12, 13, 15; Sätze in `SPEC-053`/`-055`/`-057` | erfüllt | Diff gelesen; Hook-Liste, Kopf-Kommentar und Vertrag nennen dieselben fünf Einträge; Beleg Punkt 13 trägt (`.dockerignore` ohne `docs/`, nur `docs-check` liest den Pfad, akzeptiertes Negativ fortgeschrieben); Link-Bestand „keiner auf `docs/developer`“ stimmt mit dem Index überein |
| Gate-Index: Zeile `make stop-hook-gegenprobe` in `harness/README.md` | erfüllt | Diff: Pfad in der Positivliste der Zeile |
| `make kopf-check`, `make docs-check` grün; `make gates` einmal grün | erfüllt, `gates` nur als Beleg | siehe oben |
| Review-Report liegt vor | erfüllt | Datei vorhanden, 0 HIGH/MEDIUM/LOW |
| Je Zusage Zusage · Mutation · roter Test (§3.10, je Element) | erfüllt | Tabelle in §7 deckt Tiefe, drei Grenzformen, Groß-/Kleinschreibung, Art X, Link-Ziel, Verzeichnis-Ziel, Spiegelregel; die „nur sie“-Prüfung der Gegenprobe erzwingt die Mengen |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen für den Planner | siehe unten |

## Plan-vs-Code-Diff (§3.9)

- §3 (fünf Zeilen der Tabelle) entspricht dem Diff Datei für Datei; nichts außerhalb von §3 im Diff; Produktcode, `docs/developer/`, Commands, Agenten und `conventions.md` unberührt (§1-Abgrenzung gehalten); `grep` nach `docs/developer` in `.claude/commands` ist leer.
- Abweichung `p15-ziel-verzeichnis-developer` (hängt an `link_inert`, Mutation `m-verzeichnis-inert`): in §3 und §7 genannt, deckt sich mit dem Code.
- Kopf (`Bezug`, `Berührte Spec-Stellen`) führt die Punkte 12, 13, 15 und die drei Nehmer; `kopf-check` grün.
- Metrik und Marke: 89 nach der Metrik von §3.14, unter der Schätzung 120 bis 220 (Überschätzung, harmlos).
- Zählung §8: drei Liefer-Punkte, zwei Schichten; der Diff berührt weder Doku noch Produktcode. Kein Plan nennt die Kennung (`grep` nur die eigene Datei): keine Sendung eingegangen.

## Findings

| ID | Schwere | Befund |
|---|---|---|
| V-1 | INFO | Bestätigt Review F-2: `p15-ziel-verzeichnis-developer` setzt nur `../developer/`; §3 und Punkt 15 nennen auch die Form ohne Schrägstrich. Gefangen wird sie über andere Fälle (meine Mutation `docs/developer*` färbt den Fall mit), einen eigenen Fall hat sie nicht. Nicht blockierend; der Bestand für `docs/user` hat die umgekehrte Lücke. |
| V-2 | INFO | Voller `make gates` nicht von mir nachgefahren; Exit 0 beruht auf dem Beleg des Implementers in §7 (plausibel: Hook gibt am HEAD frei). |

Keine DoD-Verletzung, keine Zusage ohne Prüfung, kein Mismatch zwischen Plan und Code gefunden.

## Voraussetzungen der Closure (Planner schreibt sie)

- Register `BEO-REPO/stop-hook-positivliste-ohne-docs-developer`: `observation.md` und `state.md` stehen bei 1× und „offen“; `evidence/slice-harness-stop-hook-positivliste-developer.md` **fehlt noch** (DoD-Punkt „Register fortgeschrieben“, kein Zähler setzen; §8 spricht von 2×, `state.md` ist bei der Closure aus den Dateien zu lesen). Außerdem nennt `observation.md` die Positivliste noch ohne `docs/developer/`: Der Eintrag ist bei der Closure nachzuziehen (Stand `offen` bleibt unter 3×; ob die Regel damit schon verkörpert ist, entscheidet der Planner).
- Risiko-Ausgänge in §6: (1) künftiger Leser von `docs/developer/` → *weiter offen*, ins Register (Voraussetzung: Evidenz-Datei, siehe oben); (2) Sätze der Nehmer-Verträge → *entfallen*, begründet, durch den Diff getragen; (3) Zählung → *entfallen*, begründet, `grep` bestätigt. Alle drei tragen.
- Lerneintrag und Finding-Klassen aus dem Review (Historie-Zeile bei Vertragsänderung fehlt; Aufzählungselement ohne eigenen Fall) gehören in §7 und ins Register; DoD-Häkchen und Closure-Notiz stehen aus, §7 „Closure-Notiz“ ist noch leer.

## Verdikt

**Bestätigt** für die Liefer-Punkte (Hook, Gegenprobe, Vertrag, Gate-Index). Offen sind allein die Pflichten der Closure (Notiz, Häkchen, Register-Evidenz, Lerneintrag); keine davon ist eine Verletzung der Lieferung.
