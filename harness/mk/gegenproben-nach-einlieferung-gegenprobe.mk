# harness/mk/gegenproben-nach-einlieferung-gegenprobe.mk — Gegenprobe des Hooks für die Einlieferung in den Gegenproben-Container (ADR-0048;
# Vertrag harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md, Punkte 1 bis 13) · seit slice-harness-gegenproben-aus-gates-fahrstelle.
# KEIN GATE: sie steht in GEGENPROBEN_CONTAINER (Namensregel, SPEC-060 Punkt 2), nicht in GATE_CHECKS. Sie fährt
# tools/harness/gegenproben-nach-einlieferung.sh als hooks/post-receive eines Temp-Bare-Repos mit einer Attrappe als make und ruft Docker nicht;
# sie braucht bash ≥ 5, git, flock, setsid, pgrep, timeout und einen Benutzer ohne root.
.PHONY: gegenproben-nach-einlieferung-gegenprobe

gegenproben-nach-einlieferung-gegenprobe: ## Gegenprobe des Hooks für die Einlieferung: Fälle und Mutationen je Punkt 1 bis 13 von SPEC-061 (KEIN Gate)
	@bash tools/dev/harness/gegenproben-nach-einlieferung-gegenprobe.sh

GEGENPROBEN_CONTAINER += gegenproben-nach-einlieferung-gegenprobe
