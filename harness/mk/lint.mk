# harness/mk/lint.mk — Lint ueber die Stufe `lint` des Dockerfile (ADR-0034).
# `lint` ist ein Werkzeug, kein Gate: das Fragment haengt nichts an GATE_CHECKS.
# Der Lauf endet mit Fehlerstatus bei einem Befund (SPEC-049 Punkt 9) und
# schreibt nichts in den Arbeitsbaum. Netz braucht nur die Stufe deps.
.PHONY: lint

lint: ## Go-Code nach .golangci.yml und eigenen Pruefungen pruefen (Werkzeug, kein Gate; netzlos ausser deps)
	$(DOCKER_BUILD) --target lint .
