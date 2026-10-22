# harness/mk/release.mk — Gegenprobe von .github/workflows/release.yml und tools/harness/release.sh (SPEC-067)
# · seit slice-erster-release-veroeffentlichung.
# `release-gegenprobe` steht in GEGENPROBEN_CONTAINER, nicht in GATE_CHECKS (SPEC-060): sie haelt den Text von release.yml (Punkt 1), die Schritte von release.sh
# mit Attrappen fuer make, gh, docker (buildx) und skopeo und fuer tools/harness/container-smoke.sh (Punkte 2 bis 11) und den Abbruch (Punkt 13), je Zusage
# mit Mutationen. Sie braucht bash, git, GNU sed, grep, awk, cmp, tar, pgrep, setsid, timeout, sha256sum und die Werkzeuge von coreutils; sie ruft weder Docker noch make
# auf, und die Pipeline selbst ist kein Gate dieser Kette.
.PHONY: release-gegenprobe

release-gegenprobe: ## Gegenprobe der Veroeffentlichung: Text von release.yml (SPEC-067 Punkt 1), Schritte von release.sh (Punkte 2 bis 11), Abbruch (Punkt 13)
	@bash tools/dev/harness/release-gegenprobe.sh

GEGENPROBEN_CONTAINER += release-gegenprobe
