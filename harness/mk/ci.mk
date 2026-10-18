# harness/mk/ci.mk — Gegenprobe von .github/workflows/ci.yml und tools/harness/ci-range.sh (SPEC-064 Teile 1 und 2)
# · seit slice-harness-ci.
# `ci-gegenprobe` steht in GEGENPROBEN_CONTAINER, nicht in GATE_CHECKS (SPEC-060): sie haelt den Text von ci.yml (Punkt 1), die Berechnung der
# Range in kleinen Git-Repositories (Punkte 6 bis 9) und den Abbruch (Punkt 18), je Zusage mit Mutationen. Sie braucht bash, git, GNU sed, grep, pgrep,
# setsid, timeout und die Werkzeuge von coreutils, weder Docker noch make zum Pruefen; die CI selbst ist kein Gate dieser Kette.
.PHONY: ci-gegenprobe

ci-gegenprobe: ## Gegenprobe der CI: Text von ci.yml (SPEC-064 Punkt 1), Range von ci-range.sh (Punkte 6 bis 9), Abbruch (Punkt 18)
	@bash tools/harness/ci-gegenprobe.sh

GEGENPROBEN_CONTAINER += ci-gegenprobe
