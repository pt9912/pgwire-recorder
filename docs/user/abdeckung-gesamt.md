# Abdeckung je Anforderung

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`) aus den
Abdeckungs-Deklarationen aller Tests. Je Anforderung zeigt die Tabelle, welche
Nachweisart welchen Pfad belegt: bei funktionalen Anforderungen die drei
Akzeptanzkriterien des Lastenhefts, bei Qualitätsanforderungen die Messung.
Vollständig belegte Anforderungen stehen zusätzlich in
`abdeckung-vollstaendig.md`, die `make doc-trace` liest.

| Anforderung | Happy | Boundary | Negative | Messung | Stand |
| --- | --- | --- | --- | --- | --- |
| [`LH-FA-02`](../../spec/lastenheft.md) | E2E, Unit | E2E, Unit | E2E, Unit | n/a | vollständig |
| [`LH-FA-05`](../../spec/lastenheft.md) | — | E2E, Unit | E2E, Unit | n/a | teilweise |
| [`LH-FA-06`](../../spec/lastenheft.md) | — | — | E2E, Unit | n/a | teilweise |
| [`LH-FA-07`](../../spec/lastenheft.md) | — | — | Unit | n/a | teilweise |
| [`LH-FA-13`](../../spec/lastenheft.md) | E2E | E2E | E2E | n/a | vollständig |
| [`LH-QA-06`](../../spec/lastenheft.md) | n/a | n/a | n/a | Unit | vollständig |
