# Dockerfile — Build, Test, Integration, Lint und Produkt-Image (ADR-0026).
# Eingangs-Images sind per Digest gepinnt; eine Anhebung ist ein bewusster
# Commit. Netz braucht nur die Stufe deps; alle Stufen danach laufen mit
# `RUN --network=none`.

# --- deps: Module aus go.mod und go.sum, geprüft gegen go.sum. Die Stufe und
# alle von ihr abgeleiteten laufen auf der Plattform des Bau-Hosts; `build`
# kompiliert per GOOS/GOARCH auf die Zielplattform.
FROM --platform=$BUILDPLATFORM golang:1.27-alpine@sha256:f92b6ef800e499660581efdabdf25d9d817a9d124eaf900924f0504e7e27e12d AS deps
WORKDIR /src
ENV CGO_ENABLED=0 GOTOOLCHAIN=local GOFLAGS=-mod=readonly
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# --- source: der Quellstand; .dockerignore lässt nur go.mod, go.sum, cmd/,
# internal/, test/, .golangci.yml, tools/harness/lint.sh und tools/harness/coverage.sh in den
# Kontext.
FROM deps AS source
COPY . .

# --- test: Formatierung (gofmt; nennt die nicht formatierten Dateien), vet
# über allen Code einschließlich der Integrationstests, dazu die Unit-Tests mit
# Profil (`-covermode=set`, ohne `-coverpkg`). Die zweite Schicht vergleicht das
# Profil mit der Schwelle (tools/harness/coverage.sh; Vertrag SPEC-056). Die Schwelle
# kommt als Build-Argument ohne Vorgabewert von `make test`; sie steht nur in SPEC-056.
# Formatierung · seit slice-extended-query-record.
FROM source AS test
RUN --network=none f="$(gofmt -l .)" && { [ -z "$f" ] || { echo "nicht gofmt-formatiert:"; echo "$f"; exit 1; }; } \
 && go vet -tags integration ./... \
 && go test -trimpath -buildvcs=false -covermode=set -coverprofile=/tmp/coverage.out ./...
ARG COVERAGE_SCHWELLE
RUN --network=none sh tools/harness/coverage.sh /tmp/coverage.out "$COVERAGE_SCHWELLE"

# --- build: das Binary für die Zielplattform, reproduzierbar (-trimpath,
# ohne VCS-Stempel).
FROM source AS build
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
RUN --network=none GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -buildvcs=false \
      -ldflags="-s -w -X main.version=$VERSION" \
      -o /out/pgwire-recorder ./cmd/pgwire-recorder

# --- dist: nur das Binary der Zielplattform; `make binaries` und `make test-container`
# lesen es per `--output type=local` aus (Vertrag SPEC-055).
FROM scratch AS dist
COPY --from=build /out/pgwire-recorder /pgwire-recorder

# --- integration: Testbinary der Integrationstests und das Binary des
# Bau-Hosts in einem Image; `make test-integration` startet es im Netz einer
# PostgreSQL-Instanz.
FROM source AS integration
RUN --network=none go build -trimpath -buildvcs=false -o /out/pgwire-recorder ./cmd/pgwire-recorder \
 && go test -c -tags integration -trimpath -buildvcs=false -o /out/integration.test ./test/integration
ENV PGR_BINARY=/out/pgwire-recorder
ENTRYPOINT ["/out/integration.test", "-test.v", "-test.count=1"]

# --- lint: golangci-lint nach dem Profil .golangci.yml und die eigenen
# Prüfungen aus tools/harness/lint.sh, auf der Plattform des Bau-Hosts, mit den
# Modulen aus deps, ohne Netz. Teil der Gate-Kette über `make lint`; welche
# Zusagen tools/harness/lint-gegenprobe.sh prüft und welche offen sind, nennt ihr Kopf.
FROM --platform=$BUILDPLATFORM golangci/golangci-lint:v2.14.0@sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f AS lint
WORKDIR /src
ENV CGO_ENABLED=0 GOTOOLCHAIN=local GOFLAGS=-mod=readonly
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .
RUN --network=none bash tools/harness/lint.sh

# --- runtime: distroless, nonroot, nur das Binary.
FROM gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab AS runtime
COPY --from=build /out/pgwire-recorder /pgwire-recorder
USER nonroot
ENTRYPOINT ["/pgwire-recorder"]
