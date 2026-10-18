# harness/mk/container-start.mk — das Werkzeug `container-start` (ADR-0045; Vertrag
# harness/vertraege/SPEC-058-container-start.md). KEIN GATE: das Ziel steht nicht in GATE_CHECKS, und
# kein Hook und kein Stand glaubt sein Ergebnis. Das Rezept ist eine Zeile ohne $(MAKE) und ohne `+`,
# das Ziel .PHONY ohne Voraussetzungen: unter -n, -q und -t läuft nichts (SPEC-058 Punkt 14).
# Es braucht docker, u-boot, devcontainer, node, ssh und ssh-keygen am Host (SPEC-058 Punkt 3).
.PHONY: container-start

container-start: ## Dev-Container-Gerüst unter $PGR_DEV anlegen — KEIN Gate
	@bash tools/harness/container-start.sh
