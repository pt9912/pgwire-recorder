# harness/mk/hook-gegenprobe.mk — Gegenprobe des Commit-Traegers (ADR-0025).
# Sie haengt an GATE_CHECKS und prueft .githooks/commit-msg mit Message-Dateien:
# eine vorhandene Slice-Kennung wird angenommen, eine erfundene und eine fehlende
# Kennung abgelehnt. Sie braucht weder Docker noch Netz.
.PHONY: hook-gegenprobe

hook-gegenprobe: ## Gegenprobe des Commit-Traegers: benannte Slice-Kennungen
	@bash tools/hook/commit-msg-gegenprobe.sh

GATE_CHECKS += hook-gegenprobe
