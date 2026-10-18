# harness/mk/workflows.mk — Gate der Workflow-Dateien unter .github/workflows (ADR-0050)
# · seit slice-harness-ci-workflows.
# `workflows-check` haengt an GATE_CHECKS, `workflows-check-gegenprobe` an GEGENPROBEN_CONTAINER (SPEC-060):
# das eine prueft jede Workflow-Datei auf Syntax und Ausdruecke (actionlint im Image unten), Rechte und
# Zeitgrenzen und darauf, dass .d-check.yml die Pin-Pruefung bindet (Vertrag
# harness/vertraege/SPEC-063-workflows-check.md); die Pins selbst prueft allein `docs-check`. Die Gegenprobe haelt
# das Skript an Temp-Baeumen und faehrt drei Faelle gegen die echten Images. Das Gate braucht bash, awk,
# find, sort, readlink, mktemp und Docker; die Gegenprobe dazu make, mktemp, sed, cksum und die Werkzeuge von coreutils.
# DOCKER kommt aus build.mk; DCHECK_REF aus d-check.mk.
.PHONY: workflows-check workflows-check-gegenprobe

# Das Image von actionlint: Repository, Digest des Index; der Tag steht hinter dem Digest.
WORKFLOWS_ACTIONLINT := rhysd/actionlint@sha256:887a259a5a534f3c4f36cb02dca341673c6089431057242cdc931e9f133147e9# 1.7.7

export DOCKER
export WORKFLOWS_ACTIONLINT
workflows-check: ## Workflow-Dateien pruefen: Syntax, Rechte, Zeitgrenzen und Bindung der Pin-Pruefung (Gate)
	@bash tools/harness/workflows-check.sh

workflows-check-gegenprobe: ## Gegenprobe des Workflow-Gates: Faelle und Mutationen je Punkt von SPEC-063
	@DOCKER='$(DOCKER)' WORKFLOWS_ACTIONLINT='$(WORKFLOWS_ACTIONLINT)' DCHECK_REF='$(DCHECK_REF)' bash tools/harness/workflows-check-gegenprobe.sh

GATE_CHECKS += workflows-check
GEGENPROBEN_CONTAINER += workflows-check-gegenprobe
