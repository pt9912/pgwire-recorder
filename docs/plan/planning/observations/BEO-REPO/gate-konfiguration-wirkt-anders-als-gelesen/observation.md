# Gate-Konfiguration wirkt anders, als sie gelesen wird

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Regel der Gate-Konfiguration liest sich als Zusage, das Werkzeug wertet sie
aber anders aus: zwei `tech`-Einträge mit demselben Muster in `.a-check.yml` wertet
a-check nur einmal aus. Ein grünes Gate belegte die Zusage nicht; erst eine
Gegenprobe mit erlaubtem und verbotenem Fall machte es sichtbar.
