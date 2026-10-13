# harness/mk/stop-hook-gegenprobe.mk — Gegenprobe des Stop-Hooks und der Dateiliste des vollen
# Laufs (ADR-0038; Vertrag SPEC-052, Punkte 1, 10, 11, 14, 16 bis 20)
# · seit slice-harness-gate-kette-stop-hook-liste.
# Sie hängt an GATE_CHECKS und prüft .claude/hooks/stop-require-gates.sh mit
# tools/harness/dateiliste.sh, working-tree-hash.sh und record-gates.sh in Temp-Repos. Sie
# braucht kein Docker und kein Netz; sie braucht bash, git, make, mktemp, awk, sed, cmp, comm
# und die Werkzeuge von coreutils (cp, mv, rm, ln, chmod, readlink, sha256sum, sort, tr).
# Sie steht in der Ausnahmeliste der schnellen Kette (gates-schnell.mk).
.PHONY: stop-hook-gegenprobe

stop-hook-gegenprobe: ## Gegenprobe des Stop-Hooks und der Dateiliste: Fälle je Punkt von SPEC-052, je Mutant rot
	@bash tools/harness/stop-hook-gegenprobe.sh

GATE_CHECKS += stop-hook-gegenprobe
