# harness/mk/range-konfiguration.mk — Gegenprobe der Blöcke `commits` und `vcs` in .d-check.yml (SPEC-064 Teil 3)
# · seit slice-harness-ci-range-konfiguration.
# `range-konfiguration-gegenprobe` steht in GEGENPROBEN_CONTAINER, nicht in GATE_CHECKS (SPEC-060): sie faehrt je Fall der Punkte 11
# und 12 einen Temp-Baum mit dem echten Image von d-check (DCHECK_REF aus d-check.mk) und je Zusage und Element eine Mutation von
# .d-check.yml. Sie braucht bash, make, git, sed, grep, Docker und die Werkzeuge von coreutils. DOCKER kommt aus build.mk.
.PHONY: range-konfiguration-gegenprobe

range-konfiguration-gegenprobe: ## Gegenprobe der Range-Pruefungen: Faelle und Mutationen von .d-check.yml je Punkt 11 und 12 von SPEC-064
	@DOCKER='$(DOCKER)' DCHECK_REF='$(DCHECK_REF)' bash tools/harness/range-konfiguration-gegenprobe.sh

GEGENPROBEN_CONTAINER += range-konfiguration-gegenprobe
