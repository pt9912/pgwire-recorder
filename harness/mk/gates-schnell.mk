# harness/mk/gates-schnell.mk — die schnelle Gate-Kette (ADR-0038; Vertrag harness/vertraege/SPEC-052-gate-kette.md)
# · seit slice-harness-gate-kette-schnell.
# `gates-schnell` fährt GATE_CHECKS ohne GATE_SCHNELL_AUSNAHME und schreibt bei
# grünem Lauf den Stand .harness/state/gates-schnell-passed.diffsha; den Nachweis
# des vollen Laufs (`make gates`) schreibt und ersetzt es nie. Die Ausnahmeliste
# nennt test-container, test-integration und test-postgres-versionen;
# `gates-schnell-gegenprobe` hält sie gegen GATE_CHECKS und hält die Namensregel der
# Gegenproben (SPEC-060 Punkte 1 bis 3): Jedes Ziel mit Namen auf `-gegenprobe` und
# `a-check-negativ` steht in genau einer der Listen GEGENPROBEN_CONTAINER und
# GEGENPROBEN_HOST und in keiner Hälfte von GATE_CHECKS. Die Gegenprobe steht in
# GEGENPROBEN_CONTAINER. Die Gates der Kette brauchen, was ihre Ziele brauchen (`build`,
# `test`, `lint`, `a-check`, `docs-check` fahren über Docker); die Gegenprobe fährt
# Attrappen als Gates und ruft Docker nicht.
# GEPRÜFT DURCH tools/dev/harness/gates-schnell-gegenprobe.sh.
GATE_SCHNELL_AUSNAHME := test-container test-integration test-postgres-versionen

.PHONY: gates-schnell gates-schnell-gegenprobe

gates-schnell: ## Schnelle Gate-Kette (ohne test-container, test-integration und test-postgres-versionen); schreibt den Stand der schnellen Kette
	@bash tools/harness/gates-schnell.sh '$(MAKE)' $(filter-out $(GATE_SCHNELL_AUSNAHME),$(GATE_CHECKS))

gates-schnell-gegenprobe: ## Gegenprobe der schnellen Kette: Ausnahmeliste, Stand, voller Nachweis, Namensregel der Gegenproben
	@bash tools/dev/harness/gates-schnell-gegenprobe.sh

GEGENPROBEN_CONTAINER += gates-schnell-gegenprobe
