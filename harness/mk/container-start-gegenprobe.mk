# harness/mk/container-start-gegenprobe.mk — Gegenprobe des Werkzeugs container-start (ADR-0045, ADR-0047;
# Vertrag harness/vertraege/SPEC-058-container-start.md, Punkte 15 und 16). EIN HOST-ZIEL, KEIN GATE: sie hängt
# nicht an GATE_CHECKS und steht nicht in GATE_SCHNELL_AUSNAHME. Sie braucht bash, make, node, ssh-keygen und
# u-boot im Pfad und die GNU-Werkzeuge von coreutils; sie ruft Docker nicht und braucht kein Netz.
.PHONY: container-start-gegenprobe

container-start-gegenprobe: ## Gegenprobe von container-start (nur Host): Fälle und Mutationen je Punkt von SPEC-058 (Teil I) — KEIN Gate
	@bash tools/harness/container-start-gegenprobe.sh
