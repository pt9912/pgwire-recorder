# harness/mk/container-start-gegenprobe.mk — Gegenprobe des Werkzeugs container-start (ADR-0045; Vertrag
# harness/vertraege/SPEC-058-container-start.md, Punkt 15 und 16). Sie hängt an GATE_CHECKS und steht in
# GATE_SCHNELL_AUSNAHME (harness/mk/gates-schnell.mk). Sie braucht bash, make, node, ssh-keygen und u-boot
# im Pfad und die GNU-Werkzeuge von coreutils; sie ruft Docker nicht und braucht kein Netz.
.PHONY: container-start-gegenprobe

container-start-gegenprobe: ## Gegenprobe von container-start: Fälle und Mutationen je Punkt von SPEC-058 (Teil I)
	@bash tools/harness/container-start-gegenprobe.sh

GATE_CHECKS += container-start-gegenprobe
