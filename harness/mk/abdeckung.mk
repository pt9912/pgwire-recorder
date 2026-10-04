# harness/mk/abdeckung.mk — Abdeckungstabellen je Anforderung und Pfad.
# `abdeckung-check` haengt an GATE_CHECKS und prueft nur: Die Tabellen unter
# docs/user/ entsprechen den Abdeckungs-Deklarationen der Tests. Geschrieben
# werden sie ausschliesslich ueber `make abdeckung` (Werkzeug, kein Gate), damit
# kein Gate eine Datei schreibt, die ein anderes Gate liest.
.PHONY: abdeckung abdeckung-check

abdeckung: ## Abdeckungstabellen aus den Test-Deklarationen schreiben (Werkzeug, kein Gate)
	@bash tools/test/abdeckung.sh

abdeckung-check: ## Abdeckungstabellen gegen die Test-Deklarationen pruefen
	@bash tools/test/abdeckung.sh --check

GATE_CHECKS += abdeckung-check
