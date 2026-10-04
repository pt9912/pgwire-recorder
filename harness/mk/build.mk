# harness/mk/build.mk — Build und Test über das Multistage-Dockerfile (ADR-0026).
# `build` und `test` haengen an GATE_CHECKS; der Root-Aggregator faehrt sie via
# make gates. Netz braucht nur die Stufe deps des Dockerfiles, und nur bei leerem
# Build-Cache oder geaenderten Abhaengigkeiten.
#
# `go-mod-tidy` ist ein Werkzeug, kein Gate: es aktualisiert go.mod und go.sum mit
# Netz im gepinnten Go-Image und schreibt beide in den Arbeitsbaum.
GO_IMAGE ?= golang:1.27-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414

DOCKER ?= docker
DOCKER_BUILD ?= $(DOCKER) build --progress=plain

.PHONY: build test go-mod-tidy

build: ## Binary und Produkt-Image bauen (Dockerfile, Stufe runtime; netzlos ausser deps)
	$(DOCKER_BUILD) --target runtime -t pgwire-recorder:dev .

test: ## vet und Unit-Tests (Dockerfile, Stufe test; netzlos ausser deps)
	$(DOCKER_BUILD) --target test -t pgwire-recorder:test .

go-mod-tidy: ## go.mod und go.sum mit Netz aktualisieren (Werkzeug, kein Gate)
	$(DOCKER) run --rm -u "$$(id -u):$$(id -g)" -v "$(CURDIR)":/src -w /src \
		-e GOTOOLCHAIN=local -e GOCACHE=/tmp/go-cache -e GOPATH=/tmp/go -e HOME=/tmp \
		$(GO_IMAGE) go mod tidy

GATE_CHECKS += build test
