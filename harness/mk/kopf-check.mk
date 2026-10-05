# harness/mk/kopf-check.mk — Kopf-Sensor fuer Slice-Plaene (ADR-0032).
# `kopf-check` und `kopf-check-gegenprobe` haengen an GATE_CHECKS: das eine
# prueft, dass der Kopf jedes Slice-Plans in open/, next/ und in-progress/ die
# Kennungen aus §1 und §2 fuehrt, das andere haelt das Skript an Temp-Baeumen je
# Nummer der ADR. Beide laufen auf dem Host (bash, awk, sort, find).
.PHONY: kopf-check kopf-check-gegenprobe

kopf-check: ## Kopf der lebenden Slice-Plaene gegen die Kennungen aus §1 und §2 pruefen
	@bash tools/harness/kopf-check.sh

kopf-check-gegenprobe: ## Gegenprobe des Kopf-Sensors: Fehlformen je Nummer von ADR-0032
	@bash tools/harness/kopf-check-gegenprobe.sh

GATE_CHECKS += kopf-check kopf-check-gegenprobe
