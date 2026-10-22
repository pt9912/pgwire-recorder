# harness/mk/arch-negativ.mk — Gegenprobe des Architektur-Gates (ADR-0001, ADR-0010,
# ADR-0027, ADR-0036, ADR-0024, ADR-0015)
# · seit slice-walking-skeleton-build-gates.
# Sie haengt an GEGENPROBEN_CONTAINER (SPEC-060). Sie belegt fuer pgproto3 und crypto/tls, dass a-check
# einen Import im Domain Model und im Recording-Adapter ablehnt und in beiden
# PGWire-Adaptern zulaesst, und fuer jeden der beiden YAML-Modulpfade, dass
# a-check ihn im Domain Model, im Postgres-Adapter und im PGWire-Adapter ablehnt
# und im CLI-Adapter zulaesst, und fuer modernc.org/sqlite, dass a-check den Import im
# Domain Model, im Postgres-Adapter, im PGWire-Adapter und im CLI-Adapter ablehnt und
# im Recording-Adapter zulaesst, und fuer time (ADR-0015), dass a-check den Import im
# Domain Model, in den Driven Ports und in den Services ablehnt, im PGWire-Adapter zulaesst, im Composition Root nicht
# ablehnt (von der Regel ausgenommen) und runtime im Domain Model zulaesst; andere Regeln
# prueft sie nicht.
# Image und Runtime stammen aus a-check.mk.
.PHONY: a-check-negativ

a-check-negativ: ## Architektur-Gate-Gegenprobe: pgproto3 und crypto/tls nur in den PGWire-Adaptern, YAML nicht im PGWire-Adapter, aber im CLI-Adapter, SQLite nur im Recording-Adapter, time nicht im Kern
	@A_CHECK_IMAGE='$(A_CHECK_IMAGE)' DOCKER='$(DOCKER)' bash tools/dev/arch/a-check-negativ.sh

GEGENPROBEN_CONTAINER += a-check-negativ
