# Ein Beleg trägt die Zusage aus einem anderen Grund als dem benannten

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Test oder ein Mutant wird als Beleg für eine Zusage geführt, wird aber aus einem anderen Grund rot
oder bleibt aus einem anderen Grund grün: Der Mutant macht Dutzende Tests rot, weil eine Voraussetzung des
Aufbaus bricht, oder ein Fall in einer Tabelle bleibt grün, weil die Fixture eine andere Prüfung auslöst.
Namen und Zeile in §7 nennen die Zusage; ob sie der Grund ist, zeigt erst eine Sonde, die nur die Zusage
bricht.

Abgrenzung zu `BEO-REPO/gruener-mutant-faelschlich-aequivalent` (ein grüner Mutant wird falsch eingeordnet)
und `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (der Test fehlt): Hier ist der Test da und der Mutant
fällt oder bleibt, aber nicht wegen der Zusage.
