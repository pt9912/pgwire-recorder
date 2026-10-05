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
| [`LH-FA-03`](../../spec/lastenheft.md) | E2E | — | E2E, Unit | n/a | teilweise |
| [`LH-FA-04`](../../spec/lastenheft.md) | E2E | — | — | n/a | teilweise |
| [`LH-FA-05`](../../spec/lastenheft.md) | — | E2E, Unit | E2E, Unit | n/a | teilweise |
| [`LH-FA-06`](../../spec/lastenheft.md) | — | — | E2E, Unit | n/a | teilweise |
| [`LH-FA-07`](../../spec/lastenheft.md) | E2E | — | Unit | n/a | teilweise |
| [`LH-FA-09`](../../spec/lastenheft.md) | E2E | Unit | E2E, Unit | n/a | vollständig |
| [`LH-FA-10`](../../spec/lastenheft.md) | E2E, Unit | Unit | Unit | n/a | vollständig |
| [`LH-FA-11`](../../spec/lastenheft.md) | E2E, Unit | — | — | n/a | teilweise |
| [`LH-FA-13`](../../spec/lastenheft.md) | E2E | E2E, Unit | E2E, Unit | n/a | vollständig |
| [`LH-FA-17`](../../spec/lastenheft.md) | — | E2E, Unit | Unit | n/a | teilweise |
| [`LH-FA-18`](../../spec/lastenheft.md) | E2E, Unit | E2E, Unit | E2E, Unit | n/a | vollständig |
| [`LH-QA-02`](../../spec/lastenheft.md) | n/a | n/a | n/a | E2E | vollständig |
| [`LH-QA-06`](../../spec/lastenheft.md) | n/a | n/a | n/a | Unit | vollständig |
