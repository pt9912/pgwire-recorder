# harness/mk/hook-gegenprobe.mk — Gegenprobe des Commit-Traegers (ADR-0025, ADR-0029).
# Sie haengt an GATE_CHECKS und prueft .githooks/commit-msg mit Message-Dateien:
# vorhandene Slice- und Welle-Kennungen werden angenommen, erfundene Kennungen,
# Namen von Closure-Notizen und eine fehlende Kennung abgelehnt. Sie braucht weder
# Docker noch Netz.
.PHONY: hook-gegenprobe

hook-gegenprobe: ## Gegenprobe des Commit-Traegers: benannte Slice- und Welle-Kennungen
	@bash tools/hook/commit-msg-gegenprobe.sh

GATE_CHECKS += hook-gegenprobe
