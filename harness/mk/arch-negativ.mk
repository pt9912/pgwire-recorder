# harness/mk/arch-negativ.mk — Gegenprobe des Architektur-Gates (ADR-0001).
# Sie haengt an GATE_CHECKS: Ein gruenes `make a-check` belegt nur, dass nichts
# bricht; die Gegenprobe belegt, dass das Gate einen verbotenen Import meldet.
# Image und Runtime stammen aus a-check.mk.
.PHONY: a-check-negativ

a-check-negativ: ## Architektur-Gate-Gegenprobe: pgproto3-Import im Domain Model muss scheitern
	@A_CHECK_IMAGE='$(A_CHECK_IMAGE)' DOCKER='$(DOCKER)' bash tools/arch/a-check-negativ.sh

GATE_CHECKS += a-check-negativ
