# harness/mk/lint.mk — Lint-Gate ueber die Stufe `lint` des Dockerfile (ADR-0034)
# · seit slice-harness-lint.
# `lint` und `lint-gegenprobe` haengen an GATE_CHECKS. `lint` prueft den Go-Code
# nach .golangci.yml und den eigenen Pruefungen aus tools/harness/lint.sh und endet
# bei einem Befund mit Fehlerstatus (Vertrag SPEC-049, Lesart
# harness/sensors/lint.md). Die Gegenprobe faehrt Mutanten in Kopien des
# Arbeitsbaums und prueft je Fall Ausgang und Zeile; sie laeuft auf dem Host mit
# bash, make, docker, grep, sed, awk und coreutils. Was sie nicht prueft, nennt ihr
# Kopf unter OFFEN.
# GEPRUEFT DURCH tools/harness/lint-gegenprobe.sh: p10-gate-checks, p10-gates-rot,
# p10-gegenprobe-scharf.
.PHONY: lint lint-gegenprobe

lint: ## Go-Code nach .golangci.yml und eigenen Pruefungen pruefen (Gate; netzlos ausser deps)
	$(DOCKER_BUILD) --target lint .

lint-gegenprobe: ## Gegenprobe des Lint-Gates: Mutanten je Punkt von SPEC-049 in Kopien des Arbeitsbaums
	@DOCKER_BUILD='$(DOCKER_BUILD)' bash tools/harness/lint-gegenprobe.sh

GATE_CHECKS += lint lint-gegenprobe
