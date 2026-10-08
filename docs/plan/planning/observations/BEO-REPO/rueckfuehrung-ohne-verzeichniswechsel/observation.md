# Eine Rückführung „zu groß“ wird geschnitten, ohne dass der Slice das Verzeichnis wechselt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die vorab benannte Rückführung *zu groß* tritt ein, und der Schnitt geschieht im selben
Zug: Der gelieferte Teil bleibt als geschnittener Slice in `in-progress/`, der Rest geht
an einen neuen Slice in `next/`. Den Übergang `in_progress → next`, den die State Machine
der Baseline für *zu groß* vorsieht, zeigt die Verzeichnis-Historie nicht; Bedingung und
Grund stehen in §4 des Slice-Plans. `MR-000` in `harness/conventions.md` erklärt für die
Lifecycle-Regeln keine Adaption.
