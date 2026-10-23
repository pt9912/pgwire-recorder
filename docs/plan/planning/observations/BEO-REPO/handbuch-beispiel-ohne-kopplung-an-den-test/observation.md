# Ein Handbuchbeispiel und sein Test hängen nur an der Tabelle im Plan zusammen

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Beispiel oder Satz des Handbuchs ist durch einen E2E-Test gegen das gebaute Binary belegt
(`AGENTS.md` §3.11), aber der Test liest das Handbuch nicht: Welcher Satz zu welchem Fall
gehört, steht allein in der Tabelle in §7 des Slice-Plans. Ändert später jemand einen Wert im
Handbuch, wird kein Test rot; die Zeile in §7 ist Archiv, kein Sensor.

Abgrenzung zu `BEO-REPO/plan-nennt-test-der-nicht-existiert`: Dort fehlt der genannte Test;
hier existiert er, aber nichts hält Beispielzeile und erwartete Zeichenkette gleich.
Kleinere Lösung vor einem Sensor: der Test liest die Beispielzeile aus dem Handbuch und
vergleicht sie mit dem Lauf.
