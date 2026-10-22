# harness/mk/container-einrichten-gegenprobe.mk — Gegenprobe des Werkzeugs container-einrichten (ADR-0057,
# ADR-0055; Vertrag harness/vertraege/SPEC-069-container-einrichten.md Punkt 17). KEIN GATE: sie steht in
# GEGENPROBEN_CONTAINER (SPEC-060 Punkte 2 und 3 — sie ruft weder u-boot noch devcontainer), nicht in
# GATE_CHECKS und nicht in GATE_SCHNELL_AUSNAHME. Sie faehrt ohne Container, ohne Docker und ohne Netz und
# braucht bash, git, make, sed, awk, grep, find und die GNU-Werkzeuge von coreutils. Obergrenze des vollen Laufs:
# 600 s (SPEC-060 Punkt 16, ADR-0055).
.PHONY: container-einrichten-gegenprobe

container-einrichten-gegenprobe: ## Gegenprobe von container-einrichten: Faelle und Mutationen zu den Punkten 1 bis 17 von SPEC-069 — KEIN Gate
	@bash tools/dev/harness/container-einrichten-gegenprobe.sh

GEGENPROBEN_CONTAINER += container-einrichten-gegenprobe
