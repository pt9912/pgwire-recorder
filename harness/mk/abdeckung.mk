# harness/mk/abdeckung.mk — Abdeckungstabellen je Anforderung und Pfad (ADR-0028, ADR-0033).
# `abdeckung-check` und `abdeckung-gegenprobe` haengen an GATE_CHECKS: das eine
# prueft, dass die Tabellen unter docs/user/abdeckung-*.md den Deklarationen der
# Tests und der Gegenproben entsprechen, das andere, dass das Skript Fehlformen ablehnt. Beide laufen
# auf dem Host (bash, awk, sort, find, cmp, stat). Geschrieben werden die
# Tabellen nur ueber `make abdeckung` (Werkzeug, kein Gate).
.PHONY: abdeckung abdeckung-check abdeckung-gegenprobe

abdeckung: ## Abdeckungstabellen aus den Deklarationen der Tests und Gegenproben schreiben (Werkzeug, kein Gate)
	@bash tools/test/abdeckung.sh

abdeckung-check: ## Abdeckungstabellen gegen die Deklarationen der Tests und Gegenproben pruefen
	@bash tools/test/abdeckung.sh --check

abdeckung-gegenprobe: ## Gegenprobe des Abdeckungs-Skripts: Fehlformen werden abgelehnt
	@bash tools/test/abdeckung-gegenprobe.sh

GATE_CHECKS += abdeckung-check abdeckung-gegenprobe
