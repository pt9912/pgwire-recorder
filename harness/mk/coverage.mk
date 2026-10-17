# harness/mk/coverage.mk — Gegenprobe der Schwelle der Anweisungsabdeckung (ADR-0044)
# · seit slice-harness-coverage.
# Die Schwelle selbst prueft die Stufe `test` (harness/mk/build.mk, Vertrag
# harness/vertraege/SPEC-056-testabdeckung-schwelle.md). `coverage-gegenprobe` haengt an
# GATE_CHECKS und steht in GATE_SCHNELL_AUSNAHME (harness/mk/gates-schnell.mk): Sie fuehrt
# Pruefschritt und Stufe an Profilen und in Kopien des Arbeitsbaums aus und braucht bash,
# make, docker, awk und coreutils. Welche Zusagen sie prueft, nennt ihr Kopf.
# GEPRUEFT DURCH tools/harness/coverage-gegenprobe.sh: p11-gate-checks.
.PHONY: coverage-gegenprobe

coverage-gegenprobe: ## Gegenprobe der Schwelle der Anweisungsabdeckung: Profile und Kopien der Stufe test
	@DOCKER='$(DOCKER)' bash tools/harness/coverage-gegenprobe.sh

GATE_CHECKS += coverage-gegenprobe
