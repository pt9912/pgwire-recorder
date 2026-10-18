# harness/mk/gegenproben.mk — die zwei Sammelziele der Gegenproben (ADR-0048; Vertrag
# harness/vertraege/SPEC-060-gegenproben.md, Punkte 4 bis 15) · seit slice-harness-gegenproben-aus-gates-sammelziele.
# `gegenproben` fährt GEGENPROBEN_CONTAINER, `gegenproben-host` fährt GEGENPROBEN_HOST, beide über
# tools/harness/gegenproben.sh, nacheinander in der Reihenfolge der Liste. KEIN GATE: beide stehen weder in
# GATE_CHECKS noch in GATE_SCHNELL_AUSNAHME, hängen nicht voneinander ab, und kein Hook und kein Stand liest ihr Ergebnis.
# Das Rezept übergibt $(MAKE) an das Skript; make führt es deshalb auch unter -n, -q und -t aus, und das Skript
# entscheidet an MAKEFLAGS (Punkt 10). Der Arbeitsbaum bleibt unberührt (Punkt 11).
# Wo die Einzelziele laufen, sagt die Dokumentation (Punkt 14); `gegenproben-host` braucht u-boot, devcontainer, node,
# ssh-keygen und Docker am Host.
# GEPRÜFT DURCH tools/harness/gegenproben-gegenprobe.sh (harness/mk/gegenproben-gegenprobe.mk).
.PHONY: gegenproben gegenproben-host

gegenproben: ## Fährt die Gegenproben der Liste für den Dev-Container nacheinander, alle auch nach einem roten (KEIN Gate)
	@bash tools/harness/gegenproben.sh '$(MAKE)' gegenproben $(GEGENPROBEN_CONTAINER)

gegenproben-host: ## Fährt die Gegenproben der Liste für den Host nacheinander, alle auch nach einem roten (KEIN Gate)
	@bash tools/harness/gegenproben.sh '$(MAKE)' gegenproben-host $(GEGENPROBEN_HOST)
