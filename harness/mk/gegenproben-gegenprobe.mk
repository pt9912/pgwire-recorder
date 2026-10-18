# harness/mk/gegenproben-gegenprobe.mk — Gegenprobe der Sammelziele (ADR-0048; Vertrag
# harness/vertraege/SPEC-060-gegenproben.md, Teil 2, Punkte 4 bis 12) · seit slice-harness-gegenproben-aus-gates-sammelziele.
# KEIN GATE: sie steht in GEGENPROBEN_CONTAINER (Namensregel, Punkt 2), nicht in GATE_CHECKS. Sie fährt
# tools/harness/gegenproben.sh in Kopien eines Temp-Verzeichnisses mit Attrappen als Einzelziele und ruft Docker nicht;
# sie braucht bash ≥ 5, make, setsid, pgrep und timeout.
.PHONY: gegenproben-gegenprobe

gegenproben-gegenprobe: ## Gegenprobe der Sammelziele: Fälle und Mutationen je Punkt 4 bis 12 von SPEC-060 (KEIN Gate)
	@bash tools/harness/gegenproben-gegenprobe.sh

GEGENPROBEN_CONTAINER += gegenproben-gegenprobe
