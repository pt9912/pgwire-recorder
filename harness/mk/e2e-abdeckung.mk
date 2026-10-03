# harness/mk/e2e-abdeckung.mk — Fragment des Erzeugers der E2E-Abdeckungs-Sicht,
# emittiert von ai-harness-init. EIN KOMMANDO, KEIN GATE.
#
# WARUM ES AN KEINER GATE-KETTE HAENGT: der Erzeuger urteilt ueber den QUELLTEXT des
# E2E-Skripts, nicht ueber den Zustand des Baums. Ein Gate an dieser Stelle faerbte rot,
# weil jemand eine Deklaration noch nicht geschrieben hat — nicht, weil etwas kaputt ist.
# Es wird ausdruecklich gerufen und schreibt seine Sicht nur bei Abweichung.
#
# DIESE DATEI IST KONVERGENT (ADR-0007 Festlegung 3): jeder Lauf des Werkzeugs schreibt
# sie kanonisch neu, und ein Edit an den Belegungen unten ist danach still weg. Sie sind
# darum die VORGABE, nicht der Setz-Ort. Dasselbe gilt fuer das Root-Makefile: es ist der
# generierte Aggregator und wird ebenso neu geschrieben. Gesetzt wird am Aufruf
#   make e2e-abdeckung E2E_ABDECKUNG_QUELLE=tools/harness/mein-e2e.sh
# oder dauerhaft in einem EIGENEN Fragment unter harness/mk/ mit einem Namen, den dieses
# Werkzeug nicht schreibt — etwa harness/mk/vorgaben.mk:
#   E2E_ABDECKUNG_QUELLE = tools/harness/mein-e2e.sh
# Der Aggregator bindet es ueber `include harness/mk/*.mk` mit ein, und kein Lauf
# entfernt, was er nicht selbst angelegt hat. Das einfache `=` gewinnt unabhaengig von
# der Include-Reihenfolge: steht es vorher, greift das `?=` unten nicht mehr; steht es
# nachher, ueberschreibt es dessen Belegung.
#
# ABHAENGIGKEIT. bash und coreutils. Kein git, kein Docker, kein Netz.
.PHONY: e2e-abdeckung

# Das E2E-Skript, dessen Stufen gelesen werden. Die Vorgabe ist die mitgelieferte
# Selbstpruefung — das eine E2E, das jedes gebootstrappte Repo fuehrt.
E2E_ABDECKUNG_QUELLE ?= tools/harness/selbstpruefung.sh
# Das Wort, mit dem eine Stufen-Kopfzeile dieses Skripts beginnt.
E2E_ABDECKUNG_PRAEFIX ?= selbstpruefung
# Die Spec-Datei, aus deren `### `-Ueberschriften die Anker der Kennungsspalte stammen.
# Eine Kennung, die dort nicht aufloest, steht als Code-Span ohne Link.
E2E_ABDECKUNG_SPEC ?= spec/lastenheft.md
# Die zu schreibende Sicht.
E2E_ABDECKUNG_ZIEL ?= docs/user/e2e-abdeckung.md

# Fehlt das Skript, sagt dieses Ziel das und bricht ab, statt auf ein fehlendes
# Programm zu zeigen.
e2e-abdeckung: ## die E2E-Abdeckungs-Sicht aus den Stufen-Deklarationen erzeugen — KEIN Gate
	@test -x tools/harness/e2e-abdeckung.sh || { echo "e2e-abdeckung: tools/harness/e2e-abdeckung.sh liegt nicht ausfuehrbar — das Werkzeug legt sie bei jedem Lauf kanonisch ab; ein erneuter Lauf stellt sie her."; exit 2; }
	@E2E_ABDECKUNG_QUELLE="$(E2E_ABDECKUNG_QUELLE)" \
	 E2E_ABDECKUNG_PRAEFIX="$(E2E_ABDECKUNG_PRAEFIX)" \
	 E2E_ABDECKUNG_SPEC="$(E2E_ABDECKUNG_SPEC)" \
	 E2E_ABDECKUNG_ZIEL="$(E2E_ABDECKUNG_ZIEL)" \
	 bash tools/harness/e2e-abdeckung.sh
