# harness/mk/archiv-bezug-gegenprobe.mk — Gegenprobe von archiv-bezug (SPEC-068) · seit slice-harness-doc-trace-archiv.
# `archiv-bezug-gegenprobe` steht in GEGENPROBEN_CONTAINER, nicht in GATE_CHECKS (SPEC-060): sie faehrt je Punkt 1 bis 12 von SPEC-068
# Faelle und Mutationen in Temp-Baeumen mit Archiven, die sie mit zip baut. Sie braucht bash, make, zip, unzip, GNU sed, grep, awk, sort, cmp,
# stat, find, md5sum und die Werkzeuge von coreutils; sie ruft kein Docker.
.PHONY: archiv-bezug-gegenprobe

archiv-bezug-gegenprobe: ## Gegenprobe von archiv-bezug: Faelle und Mutationen je Punkt 1 bis 12 von SPEC-068
	@bash tools/dev/harness/archiv-bezug-gegenprobe.sh

GEGENPROBEN_CONTAINER += archiv-bezug-gegenprobe
