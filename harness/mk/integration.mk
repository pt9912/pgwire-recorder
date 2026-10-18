# harness/mk/integration.mk — Integrationstests gegen eine reale PostgreSQL-Instanz
# (ADR-0026). `test-integration` haengt an GATE_CHECKS, `integration-volumes-gegenprobe` an
# GEGENPROBEN_CONTAINER (SPEC-060). Der Runner startet das gepinnte
# PostgreSQL-Image und das Testimage (Dockerfile, Stufe integration) in einem
# eigenen Docker-Netz; er schreibt nichts in den Arbeitsbaum. Eine zweite Phase
# stoppt PostgreSQL und spielt eine Aufzeichnung der ersten ab
# (· seit slice-walking-skeleton-replay). Der Runner entfernt seine Container mit
# ihren anonymen Volumes und prüft nach Namen, dass keines bleibt (SPEC-057);
# `integration-volumes-gegenprobe` prüft das gegen reales Docker.
# GEPRÜFT DURCH tools/test/integration-volumes-gegenprobe.sh.
POSTGRES_IMAGE ?= postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24

.PHONY: test-integration integration-volumes-gegenprobe

test-integration: ## Integrationstests gegen PostgreSQL in Docker
	@POSTGRES_IMAGE='$(POSTGRES_IMAGE)' DOCKER='$(DOCKER)' bash tools/test/run-integration-tests.sh

integration-volumes-gegenprobe: ## Gegenprobe des Aufräumens der Integrationstests: Volumes der eigenen Container, mit realem Docker
	@POSTGRES_IMAGE='$(POSTGRES_IMAGE)' DOCKER='$(DOCKER)' bash tools/test/integration-volumes-gegenprobe.sh

GATE_CHECKS += test-integration
GEGENPROBEN_CONTAINER += integration-volumes-gegenprobe
