# harness/mk/build.mk — Build und Test des Go-Moduls, nur in Docker (AGENTS.md §3.1).
# Beide Ziele haengen an GATE_CHECKS; der Root-Aggregator faehrt sie via make gates.
#
# GO_IMAGE ist per Digest gepinnt (Manifestliste fuer amd64 und arm64); eine Anhebung
# ist ein bewusster Commit. Der Lauf ist netzlos: das Modul hat keine externen
# Abhaengigkeiten, GOTOOLCHAIN=local verhindert einen Toolchain-Download.
# `go build ./...` schreibt kein Binary in den Arbeitsbaum.
GO_IMAGE ?= golang:1.27-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414

DOCKER ?= docker

GO_RUN = $(DOCKER) run --rm --network none \
	-u "$$(id -u):$$(id -g)" \
	-v "$(CURDIR)":/src:ro -w /src \
	-e GOTOOLCHAIN=local -e GOFLAGS=-mod=readonly -e CGO_ENABLED=0 \
	-e GOCACHE=/tmp/go-cache -e GOPATH=/tmp/go -e HOME=/tmp \
	$(GO_IMAGE)

.PHONY: build test

build: ## Go-Modul bauen (Docker, netzlos, schreibgeschuetzt)
	$(GO_RUN) go build -trimpath -buildvcs=false ./...

test: ## Go-Tests (Docker, netzlos, schreibgeschuetzt)
	$(GO_RUN) go test -trimpath -buildvcs=false ./...

GATE_CHECKS += build test
