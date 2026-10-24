# harness/mk/build.mk — Build und Test über das Multistage-Dockerfile (ADR-0026).
# `build` und `test` haengen an GATE_CHECKS; der Root-Aggregator faehrt sie via
# make gates. Netz braucht nur die Stufe deps des Dockerfiles, und nur bei leerem
# Build-Cache oder geaenderten Abhaengigkeiten.
#
# `test` liest die Schwelle der Anweisungsabdeckung beim Einlesen des Makefiles aus
# harness/vertraege/SPEC-056-testabdeckung-schwelle.md (Punkt 1: die Zeile mit
# `**Schwelle:**` am Zeilenanfang, Wert in Backticks; eine eingerueckte Zeile zaehlt
# nicht) und reicht sie der Stufe als Build-Argument; das Rezept ist eine Zeile ohne $(MAKE)
# (SPEC-056 Punkt 9). Zwei Zeilen dort ergeben "94 94", eine fehlende einen leeren Wert;
# beides lehnt der Pruefschritt ab.
#
# `go-mod-tidy` ist ein Werkzeug, kein Gate: es aktualisiert go.mod und go.sum mit
# Netz im gepinnten Go-Image und schreibt beide in den Arbeitsbaum.
GO_IMAGE ?= golang:1.27.3-alpine@sha256:b8c9b82affbbb82a020bbcf998dad2ba9516944ac7e883524642a481a9505e9c

COVERAGE_SCHWELLE := $(shell sed -n 's/^\*\*Schwelle:\*\* `\([^`]*\)`.*/\1/p' harness/vertraege/SPEC-056-testabdeckung-schwelle.md)

DOCKER ?= docker
DOCKER_BUILD ?= $(DOCKER) build --progress=plain

.PHONY: build test go-mod-tidy

build: ## Binary und Produkt-Image bauen (Dockerfile, Stufe runtime; netzlos ausser deps)
	$(DOCKER_BUILD) --target runtime -t pgwire-recorder:dev .

test: ## gofmt, vet, Unit-Tests und Schwelle der Anweisungsabdeckung (Dockerfile, Stufe test; netzlos ausser deps)
	$(DOCKER_BUILD) --build-arg COVERAGE_SCHWELLE="$(COVERAGE_SCHWELLE)" --target test -t pgwire-recorder:test .

go-mod-tidy: ## go.mod und go.sum mit Netz aktualisieren (Werkzeug, kein Gate)
	$(DOCKER) run --rm -u "$$(id -u):$$(id -g)" -v "$(CURDIR)":/src -w /src \
		-e GOTOOLCHAIN=local -e GOCACHE=/tmp/go-cache -e GOPATH=/tmp/go -e HOME=/tmp \
		$(GO_IMAGE) go mod tidy

GATE_CHECKS += build test
