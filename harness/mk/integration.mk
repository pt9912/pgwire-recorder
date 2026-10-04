# harness/mk/integration.mk — Integrationstests gegen eine reale PostgreSQL-Instanz
# (ADR-0026). Haengt an GATE_CHECKS. Der Runner startet das gepinnte
# PostgreSQL-Image und das Testimage (Dockerfile, Stufe integration) in einem
# eigenen Docker-Netz; er schreibt nichts in den Arbeitsbaum. Eine zweite Phase
# stoppt PostgreSQL und spielt eine Aufzeichnung der ersten ab
# (· seit slice-walking-skeleton-replay).
POSTGRES_IMAGE ?= postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24

.PHONY: test-integration

test-integration: ## Integrationstests gegen PostgreSQL in Docker
	@POSTGRES_IMAGE='$(POSTGRES_IMAGE)' DOCKER='$(DOCKER)' bash tools/test/run-integration-tests.sh

GATE_CHECKS += test-integration
