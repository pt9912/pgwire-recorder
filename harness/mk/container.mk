# harness/mk/container.mk — Binaries je Zielplattform, Container-Image und Smoke
# (ADR-0043; Vertrag harness/vertraege/SPEC-055-container-und-binaries.md)
# · seit slice-v1-abschluss-container.
# Die Listen der Plattformen stehen hier und nur hier (SPEC-055 Punkt 1); sie sind keine
# Eingabe des Aufrufers. `binaries` schreibt die Binaries mit SHA256SUMS nach DIST und ist
# ein Werkzeug; `test-container` baut sie zweimal und vergleicht die Summen, prüft Kopf,
# Image und Smoke; die Gegenprobe (GEGENPROBEN_CONTAINER, SPEC-060) fährt es mit einer Attrappe von docker. Das Rezept jedes
# Ziels ist eine Zeile ohne $(MAKE).
# GEPRÜFT DURCH tools/harness/container-gegenprobe.sh.
VERSION ?= dev
DIST ?= dist
BINARY_PLATTFORMEN := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64
IMAGE_PLATTFORMEN := linux/amd64 linux/arm64

# Die Werte gehen in die Umgebung der Rezepte, nicht in die Befehlszeile: Ein Wert mit
# Anführungszeichen oder Leerzeichen erreicht das Skript unverändert.
export VERSION DIST BINARY_PLATTFORMEN IMAGE_PLATTFORMEN

.PHONY: binaries test-container test-container-gegenprobe

binaries: ## Binaries aller Zielplattformen mit SHA256SUMS nach DIST bauen (Werkzeug, kein Gate)
	@DOCKER='$(DOCKER)' bash tools/harness/binaries.sh

test-container: VERSION := 0.0.0-test
test-container: ## zwei Builds mit gleichen Summen, Kopf der Binaries, Image je Plattform, Smoke im Container
	@POSTGRES_IMAGE='$(POSTGRES_IMAGE)' DOCKER='$(DOCKER)' bash tools/harness/run-container.sh

test-container-gegenprobe: ## Gegenprobe von binaries und test-container: Eingabe, Summen, Kopf, Image, Smoke, Aufräumen, Aufrufarten, mit Attrappe von docker
	@bash tools/harness/container-gegenprobe.sh

GATE_CHECKS += test-container
GEGENPROBEN_CONTAINER += test-container-gegenprobe
