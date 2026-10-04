# Abdeckung je Anforderung

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`) aus den
Abdeckungs-Deklarationen aller Tests. Je Anforderung zeigt die Tabelle, welche
Nachweisart welchen der drei Pfade des Lastenhefts belegt. Vollständig ist eine
Anforderung, wenn alle drei Pfade belegt sind; nur diese stehen in
`abdeckung-vollstaendig.md`, die `make doc-trace` liest.

| Anforderung | Happy | Boundary | Negative | Stand |
| --- | --- | --- | --- | --- |
| [`LH-FA-02`](../../spec/lastenheft.md) | E2E, Unit | Unit | E2E, Unit | vollständig |
| [`LH-FA-05`](../../spec/lastenheft.md) | — | E2E, Unit | E2E, Unit | teilweise |
| [`LH-FA-06`](../../spec/lastenheft.md) | E2E, Unit | E2E | — | teilweise |
| [`LH-FA-07`](../../spec/lastenheft.md) | E2E, Unit | Unit | Unit | vollständig |
| [`LH-FA-12`](../../spec/lastenheft.md) | — | Unit | — | teilweise |
| [`LH-FA-13`](../../spec/lastenheft.md) | E2E | E2E | E2E | vollständig |
| [`LH-QA-01`](../../spec/lastenheft.md) | — | Unit | — | teilweise |
| [`LH-QA-06`](../../spec/lastenheft.md) | Unit | — | Unit | teilweise |
