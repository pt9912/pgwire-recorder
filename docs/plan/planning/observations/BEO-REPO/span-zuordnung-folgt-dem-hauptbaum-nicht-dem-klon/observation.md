# Die Spans ordnen Arbeit im Klon eines Dev-Containers nach dem Zustand des Hauptbaums zu

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die Erfassung (`harness/mk/erfassung.mk`, Hook `.claude/hooks/span-emit.sh`) schreibt die Felder `slice`, `requirement`, `adr`, `branch` und `commit` aus dem Repo, in dem der Hook läuft: dem Hauptbaum des Hosts. Arbeitet ein Agent im Klon eines Dev-Containers (Aufrufe über `ssh`), steht dort der Zustand des Hauptbaums (`main`, alle Slices in `in-progress/`), nicht der des Klons. Zwei Implementer an zwei Slices tragen dieselbe Slice-Liste und denselben Zweig; die Token-Zuordnung je Slice ([Modul 15](../../../../../../.harness/baseline/v6.18.0/regelwerk/modul-15-observability.md)) ist damit nur auf die Menge aller angespruchten Slices genau. Außerdem sieht der Hook nur das erste Wort jedes Befehls (`cd`, `$S`), nicht, was im Container läuft.

Abgrenzung: Das ist keine Lücke der Gate-Kette, sondern der Erfassung. Sie war für Arbeit im Hauptbaum ausgelegt; die Arbeit im Klon (`docs/developer/devcontainer.md` §6) ist später hinzugekommen.
