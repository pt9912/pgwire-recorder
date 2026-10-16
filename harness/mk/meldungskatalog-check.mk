# harness/mk/meldungskatalog-check.mk — Abgleich von Code-Tabelle und Meldungskatalog (ADR-0011)
# · seit slice-harness-meldungskatalog-gate.
# `meldungskatalog-check` und `meldungskatalog-check-gegenprobe` haengen an GATE_CHECKS: das eine
# gleicht die Codes aus internal/hexagon/model/fehler.go, aus SPEC-034 und aus den Tabellen des
# Benutzerhandbuchs ab (Vertrag harness/vertraege/SPEC-054-meldungskatalog.md), das andere haelt das
# Skript an Temp-Baeumen je Punkt des Vertrags. Beide laufen auf dem Host: meldungskatalog-check
# braucht bash, awk, sort und find; die Gegenprobe dazu make, mktemp, grep, GNU sed (sed -i), cksum
# und die Werkzeuge von coreutils. Die Pruefung steht in der schnellen Kette, die Gegenprobe in
# GATE_SCHNELL_AUSNAHME (harness/mk/gates-schnell.mk).
.PHONY: meldungskatalog-check meldungskatalog-check-gegenprobe

meldungskatalog-check: ## Codes aus fehler.go, SPEC-034 und Benutzerhandbuch gegeneinander abgleichen
	@bash tools/harness/meldungskatalog-check.sh

meldungskatalog-check-gegenprobe: ## Gegenprobe des Abgleichs: Fehlformen je Punkt von SPEC-054
	@bash tools/harness/meldungskatalog-check-gegenprobe.sh

GATE_CHECKS += meldungskatalog-check meldungskatalog-check-gegenprobe
