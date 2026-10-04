# Vollständig abgedeckte Anforderungen

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`). Hier steht eine
Anforderung erst, wenn Tests alle drei Pfade des Lastenhefts (Happy, Boundary,
Negative) belegen; die Spalte nennt die beteiligten Nachweisarten. Diese Datei
liest `make doc-trace` (`trace.coverage`); Teilabdeckung steht in
`abdeckung-gesamt.md`.

| Anforderung | Nachweisarten |
| --- | --- |
| [`LH-FA-02`](../../spec/lastenheft.md) | E2E, Unit |
| [`LH-FA-07`](../../spec/lastenheft.md) | E2E, Unit |
| [`LH-FA-13`](../../spec/lastenheft.md) | E2E |
