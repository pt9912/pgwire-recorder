# harness/mk/archiv-bezug.mk — Anforderungen im Stub archivierter Slices (SPEC-068) · seit slice-harness-doc-trace-archiv.
# Werkzeug, KEIN Gate: `archiv-bezug` steht weder in GATE_CHECKS noch in einer Kette, und kein Hook und kein Stand liest sein Ergebnis.
# Das Rezept ruft das Skript genau einmal und enthaelt weder $(MAKE) noch ein Plus; `make -n`, `-q` und `-t` schreiben darum nichts (SPEC-068 Punkt 10).
# `make archiv-bezug` faehrt alle Archive unter docs/plan/planning/done/*/archiv.zip, `make archiv-bezug WELLE=<welle-id>` eines.
# GEPRUEFT DURCH tools/dev/harness/archiv-bezug-gegenprobe.sh (harness/mk/archiv-bezug-gegenprobe.mk).
.PHONY: archiv-bezug

archiv-bezug: ## Schreibt die Anforderungen der archivierten Slices als Zeile Bezug in ihre Stubs (WELLE=<welle-id> fuer eine Welle)
	@bash tools/harness/archiv-bezug.sh "$(CURDIR)" "$(WELLE)"
