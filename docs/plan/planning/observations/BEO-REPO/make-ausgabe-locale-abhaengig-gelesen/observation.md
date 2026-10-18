# Ein Werkzeug liest die Ausgabe von `make` im Text einer Sprache

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Skript, das die Ausgabe eines Programms mit Mustern liest, erkennt sie nur in der Sprache, in der es sie gesehen hat. GNU make 4.3 schreibt in `make -p` unter `de_DE.UTF-8` „Auszuführendes Rezept (aus „…““ statt „recipe to execute (from '…'“. Die Namensregel der Gegenproben (`tools/harness/gegenproben-regel.sh`) fand dann kein Ziel mit Rezept, und `make gates-schnell-gegenprobe` war am Host rot, im Dev-Container (Locale C) grün; die Sonde des Implementers deckte nur eine Locale. Behoben durch `LC_ALL=C` im Aufruf von `make` und den Fall `p2-locale` (Mutation `p2-locale-nicht-c`).

Abgrenzung zu `BEO-REPO/verhalten-nur-unter-linux-geprueft`: Dort fehlt eine Plattform; hier hängt dasselbe Verhalten an der Spracheinstellung derselben Plattform.
