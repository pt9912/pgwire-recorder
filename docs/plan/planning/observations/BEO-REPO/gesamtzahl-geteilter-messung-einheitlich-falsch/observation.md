# Eine einheitlich falsche Gesamtzahl einer geteilten Messung bleibt unbemerkt

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Gesamtzahl n und Zuordnung eines Teils zur Messmethode einer Anforderung (Pfad
`Messung-<i>-von-<n>`) sind Urteil dessen, der deklariert. `make abdeckung-check` fängt
verschiedene n für dieselbe Anforderung, nicht ein n, das in allen Deklarationen gleich und
falsch ist; eine Anforderung stünde dann als vollständig da, obwohl ein Teil fehlt oder
keiner ist.
**Stand:** offen

Risiko aus `slice-harness-abdeckung-gate` (§6), Ausgang weiter offen; Grenze in `SPEC-048` und [ADR-0033](../../../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Re-Evaluierungs-Trigger). Kein Auftreten: Die einzige geteilte Messung (`LH-QA-07`, n = 3) hat Teil 1 und 3 deklariert, Teil 2 fehlt bewusst und die Anforderung steht als *teilweise* da. Mit `slice-harness-coverage` kommt Teil 2.
