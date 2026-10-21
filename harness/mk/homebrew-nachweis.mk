# harness/mk/homebrew-nachweis.mk — Nachweis der Installation aus dem Tap (Abnahmeszenario 11,
# harness/vertraege/SPEC-067-release.md Punkt 20) · seit slice-erster-release-homebrew-nachweis.
# Werkzeug, KEIN Gate: `homebrew-nachweis` haengt nicht an GATE_CHECKS, steht in keiner der beiden
# Gegenproben-Listen (SPEC-060), und kein Hook und kein Stand liest sein Ergebnis. Es urteilt nicht;
# sein Ergebnis ist seine Ausgabe, die ein Mensch in §7 des Slice-Plans eintraegt (Punkt 20 (4), (7)).
# Es braucht NETZ (Tap, Anhaenge des GitHub-Release, Homebrew im Image), anders als jedes Gate dieser
# Kette; fehlt das Netz oder ist der Tap nicht erreichbar, endet es mit dem Exit-Code des gescheiterten
# Schrittes — kein Ueberspringen, kein Gruen ohne Lauf, kein Rueckfall auf eine lokale Formel.
# `-n`, `-q`, `-t`, `-i` und MAKEFLAGS haben keinen Fall (Punkt 20 (5), Grund wie Punkt 13).
# Die letzte Zeile des Rezepts nennt den Exit-Code des Skripts im Container: ueber `make` sieht der
# Aufrufer bei einem roten Lauf sonst nur die 2 von make, und der Beleg nennt den Code des Skripts
# (Punkt 20 (2), (6)).
HOMEBREW_IMAGE ?= homebrew/brew:4.6.20@sha256:b0072bfdebf5934ae24b93b44a1928a88057399b3283ffa0177bb86084fdedfd # 4.6.20, Digest des Index, gegen die Registry aufgeloest am 2026-10-21

DOCKER ?= docker

.PHONY: homebrew-nachweis

homebrew-nachweis: ## Nachweis aus dem echten Tap im gepinnten homebrew/brew-Container (braucht Netz, kein Gate)
	@echo "homebrew-nachweis: Image $(strip $(HOMEBREW_IMAGE))"
	@s=0; $(DOCKER) run --rm \
		-e HOMEBREW_NO_ANALYTICS=1 -e HOMEBREW_NO_AUTO_UPDATE=1 -e HOMEBREW_NO_ENV_HINTS=1 \
		-v "$(CURDIR)/tools/harness/homebrew-nachweis.sh":/tmp/homebrew-nachweis.sh:ro \
		$(strip $(HOMEBREW_IMAGE)) bash /tmp/homebrew-nachweis.sh --schritte < /dev/null || s=$$?; \
	echo "homebrew-nachweis: Exit-Code $$s (Code des Skripts im Container, nicht der von make)"; \
	exit $$s
