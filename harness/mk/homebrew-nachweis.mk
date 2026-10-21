# harness/mk/homebrew-nachweis.mk — Nachweis der Installation aus dem Tap (Abnahmeszenario 11,
# harness/vertraege/SPEC-067-release.md Punkt 20) · seit slice-erster-release-homebrew-nachweis.
# Werkzeug, KEIN Gate: `homebrew-nachweis` haengt nicht an GATE_CHECKS, steht in keiner der beiden
# Gegenproben-Listen (SPEC-060), und kein Hook und kein Stand liest sein Ergebnis. Es urteilt nicht;
# sein Ergebnis ist seine Ausgabe, die ein Mensch in §7 des Slice-Plans eintraegt (Punkt 20 (4), (7)).
# Es braucht NETZ (Tap, Anhaenge des GitHub-Release, Homebrew im Image), anders als jedes Gate dieser
# Kette; fehlt das Netz oder ist der Tap nicht erreichbar, endet es mit dem Exit-Code des gescheiterten
# Schrittes — kein Ueberspringen, kein Gruen ohne Lauf, kein Rueckfall auf eine lokale Formel.
# `-n`, `-q`, `-t`, `-i` und MAKEFLAGS haben keinen Fall (Punkt 20 (5), Grund wie Punkt 13).
#
# Die letzte Zeile des Rezepts nennt den Exit-Code des LAUFS (`docker run`), nicht ohne weiteres den
# des Skripts: 125, 126 und 127 meldet die Docker-CLI selbst (Image nicht ziehbar, Daemon nicht
# erreichbar, Befehl im Container nicht ausfuehrbar oder nicht gefunden) — genau der Weg von Punkt 20
# (4) —, und dann lief das Skript nicht. Wessen Code es ist, sagt die Zeile `Skript endet mit
# Exit-Code <n>`, die das Skript selbst als letzte schreibt: fehlt sie, kam es nicht zum Ende. Den
# Code 7 waehlt allein das Skript, fuer die abgelaufene Grenze des Wartens auf den Start des
# Replay-Laufs (Punkt 20 (3)); er liegt ausserhalb der Codes 0 bis 6 des Produkts. Der
# Grund fuer die Zeile des Rezepts ueberhaupt: ueber `make` sieht der Aufrufer bei einem roten Lauf
# nur die 2 von make, und der Beleg in §7 soll den Code des Laufs nennen (Punkt 20 (2), (6)).
HOMEBREW_IMAGE ?= homebrew/brew:4.6.20@sha256:b0072bfdebf5934ae24b93b44a1928a88057399b3283ffa0177bb86084fdedfd # 4.6.20, Digest des Index, gegen die Registry aufgeloest am 2026-10-21

DOCKER ?= docker

.PHONY: homebrew-nachweis

homebrew-nachweis: ## Nachweis aus dem echten Tap im gepinnten homebrew/brew-Container (braucht Netz, kein Gate)
	@echo "homebrew-nachweis: Image $(strip $(HOMEBREW_IMAGE))"
	@s=0; $(DOCKER) run --rm \
		-e HOMEBREW_NO_ANALYTICS=1 -e HOMEBREW_NO_AUTO_UPDATE=1 -e HOMEBREW_NO_ENV_HINTS=1 \
		-v "$(CURDIR)/tools/harness/homebrew-nachweis.sh":/tmp/homebrew-nachweis.sh:ro \
		$(strip $(HOMEBREW_IMAGE)) bash /tmp/homebrew-nachweis.sh --schritte < /dev/null || s=$$?; \
	case "$$s" in \
	0) echo "homebrew-nachweis: der Lauf endet mit Exit-Code 0";; \
	7) echo "homebrew-nachweis: der Lauf endet mit Exit-Code 7; den waehlt das Skript fuer die abgelaufene Grenze des Wartens auf den Start des Replay-Laufs (SPEC-067 Punkt 20 (3)) — kein Code der Docker-CLI und keiner der Codes 0 bis 6 des Produkts; der Nachweis ist nicht erbracht";; \
	125|126|127) echo "homebrew-nachweis: der Lauf endet mit Exit-Code $$s; mit 125, 126 und 127 meldet die Docker-CLI, dass sie selbst, der Daemon, das Image oder der Aufruf im Container scheiterte — steht darueber keine Zeile 'Skript endet mit Exit-Code', kam das Skript nicht zum Ende und der Nachweis ist nicht erbracht";; \
	*) echo "homebrew-nachweis: der Lauf endet mit Exit-Code $$s; wessen Code das ist, sagt die Zeile 'Skript endet mit Exit-Code' darueber — fehlt sie, stammt er nicht vom Skript";; \
	esac; \
	exit $$s
