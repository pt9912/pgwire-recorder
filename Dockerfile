# Dockerfile — Build, Test, Integration und Produkt-Image (ADR-0026).
# Eingangs-Images sind per Digest gepinnt; eine Anhebung ist ein bewusster
# Commit. Netz braucht nur die Stufe deps; alle Stufen danach laufen mit
# `RUN --network=none`.

# --- deps: Module aus go.mod und go.sum, geprüft gegen go.sum. Die Stufe und
# alle von ihr abgeleiteten laufen auf der Plattform des Bau-Hosts; `build`
# kompiliert per GOOS/GOARCH auf die Zielplattform.
FROM --platform=$BUILDPLATFORM golang:1.27-alpine@sha256:8a5910f31396cd4d89662f56c68b3ae31d374308270a1c3bd96672ee5ed43414 AS deps
WORKDIR /src
ENV CGO_ENABLED=0 GOTOOLCHAIN=local GOFLAGS=-mod=readonly
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# --- source: der Quellstand; .dockerignore lässt nur go.mod, go.sum, cmd/,
# internal/ und test/ in den Kontext.
FROM deps AS source
COPY . .

# --- test: vet über allen Code einschließlich der Integrationstests, dazu die
# Unit-Tests.
FROM source AS test
RUN --network=none go vet -tags integration ./... && go test -trimpath -buildvcs=false ./...

# --- build: das Binary für die Zielplattform, reproduzierbar (-trimpath,
# ohne VCS-Stempel).
FROM source AS build
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
RUN --network=none GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -buildvcs=false \
      -ldflags="-s -w -X main.version=$VERSION" \
      -o /out/pgwire-recorder ./cmd/pgwire-recorder

# --- integration: Testbinary der Integrationstests und das Binary des
# Bau-Hosts in einem Image; `make test-integration` startet es im Netz einer
# PostgreSQL-Instanz.
FROM source AS integration
RUN --network=none go build -trimpath -buildvcs=false -o /out/pgwire-recorder ./cmd/pgwire-recorder \
 && go test -c -tags integration -trimpath -buildvcs=false -o /out/integration.test ./test/integration
ENV PGR_BINARY=/out/pgwire-recorder PGR_FIXTURES=/src/test/integration/testdata
ENTRYPOINT ["/out/integration.test", "-test.v", "-test.count=1"]

# --- runtime: distroless, nonroot, nur das Binary.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS runtime
COPY --from=build /out/pgwire-recorder /pgwire-recorder
USER nonroot
ENTRYPOINT ["/pgwire-recorder"]
