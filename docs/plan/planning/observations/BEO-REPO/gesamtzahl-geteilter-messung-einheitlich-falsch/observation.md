# Eine einheitlich falsche Gesamtzahl einer geteilten Messung bleibt unbemerkt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Gesamtzahl n und Zuordnung eines Teils zur Messmethode einer Anforderung (Pfad
`Messung-<i>-von-<n>`) sind Urteil dessen, der deklariert. `make abdeckung-check` fängt
verschiedene n für dieselbe Anforderung, nicht ein n, das in allen Deklarationen gleich und
falsch ist; eine Anforderung stünde dann als vollständig da, obwohl ein Teil fehlt oder
keiner ist.
