# harness/mk/hook-gegenprobe.mk — Gegenprobe des Commit-Traegers (ADR-0025, ADR-0029).
# Sie haengt an GEGENPROBEN_CONTAINER (SPEC-060) und prueft .githooks/commit-msg mit Message-Dateien:
# vorhandene Slice- und Welle-Kennungen werden angenommen, erfundene Kennungen,
# Namen von Closure-Notizen und eine fehlende Kennung abgelehnt. Sie braucht weder
# Docker noch Netz.
.PHONY: hook-gegenprobe

hook-gegenprobe: ## Gegenprobe des Commit-Traegers: benannte Slice- und Welle-Kennungen
	@bash tools/dev/hook/commit-msg-gegenprobe.sh

GEGENPROBEN_CONTAINER += hook-gegenprobe
