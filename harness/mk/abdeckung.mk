# harness/mk/abdeckung.mk — Abdeckungstabellen je Anforderung und Pfad (ADR-0028, ADR-0033).
# `abdeckung-check` ruft das Skript mit `--check` (prueft, schreibt nichts),
# `abdeckung-gegenprobe` die Gegenprobe des Skripts; beide haengen an GATE_CHECKS
# und laufen auf dem Host (bash, awk, sort, find, cmp). Geschrieben werden die
# Tabellen nur ueber `make abdeckung` (Werkzeug, kein Gate). Weder das Rezept
# mit `--check` noch der Anschluss an GATE_CHECKS ist durch einen Fall der
# Gegenprobe gezeigt; das bleibt Review.
.PHONY: abdeckung abdeckung-check abdeckung-gegenprobe

abdeckung: ## Abdeckungstabellen aus den Deklarationen der Tests und Gegenproben schreiben (Werkzeug, kein Gate)
	@bash tools/test/abdeckung.sh

abdeckung-check: ## Abdeckungstabellen gegen die Deklarationen der Tests und Gegenproben pruefen
	@bash tools/test/abdeckung.sh --check

abdeckung-gegenprobe: ## Gegenprobe des Abdeckungs-Skripts: Fehlformen werden abgelehnt
	@bash tools/test/abdeckung-gegenprobe.sh

GATE_CHECKS += abdeckung-check abdeckung-gegenprobe
