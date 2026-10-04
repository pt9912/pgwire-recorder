# harness/mk/vorgaben.mk — Belegungen dieses Repos fuer Werkzeug-Fragmente, die
# ein Bootstrap kanonisch neu schreibt; dieses Fragment schreibt er nicht.
#
# Der Verweis-Nachzug von `make slice-mv` laesst die Review-Reports aus: Sie sind
# Lauf-Belege und nennen den Pfad, den ein Plan zur Zeit des Reports hatte.
SLICE_MV_AUSGENOMMENE_PFADE = :!.harness/baseline :!docs/plan/adr :!docs/reviews
