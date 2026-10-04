# harness/mk/arch-negativ.mk — Gegenprobe des Architektur-Gates (ADR-0001, ADR-0010,
# ADR-0027)
# · seit slice-walking-skeleton-build-gates.
# Sie haengt an GATE_CHECKS. Sie belegt fuer pgproto3 und crypto/tls, dass a-check
# einen Import im Domain Model und im Recording-Adapter ablehnt und in beiden
# PGWire-Adaptern zulaesst, und fuer beide YAML-Modulpfade, dass a-check sie
# ausserhalb des Recording-Adapters ablehnt; andere Regeln prueft sie nicht.
# Image und Runtime stammen aus a-check.mk.
.PHONY: a-check-negativ

a-check-negativ: ## Architektur-Gate-Gegenprobe: pgproto3 und crypto/tls nur in den PGWire-Adaptern
	@A_CHECK_IMAGE='$(A_CHECK_IMAGE)' DOCKER='$(DOCKER)' bash tools/arch/a-check-negativ.sh

GATE_CHECKS += a-check-negativ
