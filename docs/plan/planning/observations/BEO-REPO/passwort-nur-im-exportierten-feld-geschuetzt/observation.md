# Ein Passwort ist nur im exportierten Feld einer Struktur gegen Formatierung geschützt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Wert vom Typ `Passwort` gibt sich bei keiner Formatierung preis, weil der Typ `Format`
trägt. `fmt` ruft die Methode eines Felds aber nur, wenn das Feld exportiert ist: Hält eine
Struktur das Passwort in einem unexportierten Feld (`passwort Passwort` oder `passwort
string`), gibt `%+v` auf der Struktur es aus. Latent, solange nichts die Struktur formatiert.
Maßnahme je Fall: ein eigenes `Format` an der Struktur oder ein exportiertes Feld.
