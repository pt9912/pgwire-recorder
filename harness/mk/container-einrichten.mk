# harness/mk/container-einrichten.mk — das Werkzeug `container-einrichten` (ADR-0057; Vertrag
# harness/vertraege/SPEC-069-container-einrichten.md). KEIN GATE: das Ziel steht nicht in GATE_CHECKS, und
# kein Hook und kein Stand glaubt sein Ergebnis. Das Rezept ist eine Zeile ohne $(MAKE) und ohne `+`, das
# Ziel .PHONY ohne Voraussetzungen: unter -n, -q und -t laeuft nichts (SPEC-069 Punkt 16).
# Es braucht ssh und git am Host (SPEC-069 Punkt 3) und einen laufenden Dev-Container, den
# `make container-start` liefert (SPEC-058 Punkte 2 und 11).
# GEPRUEFT DURCH tools/harness/container-einrichten-gegenprobe.sh.
.PHONY: container-einrichten

container-einrichten: ## Bare-Repo, main, Klon, Identitaet, Commit-Traeger und Zweig im laufenden Dev-Container einrichten — KEIN Gate
	@bash tools/harness/container-einrichten.sh
