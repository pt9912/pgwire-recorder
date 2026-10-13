# harness/mk/gates-schnell.mk — die schnelle Gate-Kette (ADR-0038; Vertrag SPEC-052)
# · seit slice-harness-gate-kette-schnell.
# `gates-schnell` fährt GATE_CHECKS ohne GATE_SCHNELL_AUSNAHME und schreibt bei
# grünem Lauf den Stand .harness/state/gates-schnell-passed.diffsha; den Nachweis
# des vollen Laufs (`make gates`) schreibt und ersetzt es nie. Die Ausnahmeliste
# nennt jedes Gate mit `-gegenprobe` im Namen, dazu a-check-negativ und
# test-integration; `gates-schnell-gegenprobe` hält sie gegen GATE_CHECKS. Die
# Gegenprobe hängt an GATE_CHECKS, `gates-schnell` selbst nicht (es wäre sonst
# eigener Bestandteil). Beide laufen auf dem Host mit bash, make, git, sha256sum,
# sort und coreutils; die Gegenprobe braucht dazu setsid, grep, sed und awk.
# GEPRÜFT DURCH tools/harness/gates-schnell-gegenprobe.sh.
GATE_SCHNELL_AUSNAHME := abdeckung-gegenprobe a-check-negativ gates-schnell-gegenprobe hook-gegenprobe kopf-check-gegenprobe lint-gegenprobe test-integration

.PHONY: gates-schnell gates-schnell-gegenprobe

gates-schnell: ## Schnelle Gate-Kette (ohne Gegenproben und Integrationstests); schreibt den Stand der schnellen Kette
	@bash tools/harness/gates-schnell.sh '$(MAKE)' $(filter-out $(GATE_SCHNELL_AUSNAHME),$(GATE_CHECKS))

gates-schnell-gegenprobe: ## Gegenprobe der schnellen Kette: Ausnahmeliste, Stand, voller Nachweis, Anschluss an GATE_CHECKS
	@bash tools/harness/gates-schnell-gegenprobe.sh

GATE_CHECKS += gates-schnell-gegenprobe
