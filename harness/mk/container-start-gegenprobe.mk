# harness/mk/container-start-gegenprobe.mk — Gegenprobe des Werkzeugs container-start (ADR-0045, ADR-0047;
# Vertrag harness/vertraege/SPEC-058-container-start.md, Punkte 15 und 16). EIN HOST-ZIEL, KEIN GATE: sie steht in
# GEGENPROBEN_HOST (SPEC-060), nicht in GATE_CHECKS und nicht in GATE_SCHNELL_AUSNAHME. Sie braucht bash, make, node, ssh, ssh-keygen,
# u-boot, docker und devcontainer im Pfad und die GNU-Werkzeuge von coreutils; die Gruppe Docker startet einen Container des Präfixes pgr-cs- und baut ein Image.
.PHONY: container-start-gegenprobe

container-start-gegenprobe: ## Gegenprobe von container-start (nur Host): Fälle und Mutationen je Punkt von SPEC-058 — KEIN Gate
	@bash tools/harness/container-start-gegenprobe.sh

GEGENPROBEN_HOST += container-start-gegenprobe
