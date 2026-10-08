# harness/mk/arch-negativ.mk — Gegenprobe des Architektur-Gates (ADR-0001, ADR-0010,
# ADR-0027, ADR-0036)
# · seit slice-walking-skeleton-build-gates.
# Sie haengt an GATE_CHECKS. Sie belegt fuer pgproto3 und crypto/tls, dass a-check
# einen Import im Domain Model und im Recording-Adapter ablehnt und in beiden
# PGWire-Adaptern zulaesst, und fuer die YAML-Modulpfade, dass a-check sie im
# Domain Model, im Postgres-Adapter und je Modulpfad im PGWire-Adapter ablehnt
# und beide im CLI-Adapter zulaesst; andere Regeln prueft sie nicht.
# Image und Runtime stammen aus a-check.mk.
.PHONY: a-check-negativ

a-check-negativ: ## Architektur-Gate-Gegenprobe: pgproto3 und crypto/tls nur in den PGWire-Adaptern, YAML nicht im PGWire-Adapter, aber im CLI-Adapter
	@A_CHECK_IMAGE='$(A_CHECK_IMAGE)' DOCKER='$(DOCKER)' bash tools/arch/a-check-negativ.sh

GATE_CHECKS += a-check-negativ
