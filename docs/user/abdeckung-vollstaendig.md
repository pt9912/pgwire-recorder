# Vollständig abgedeckte Anforderungen

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`). Hier steht eine
Anforderung erst, wenn Tests alle ihre Pfade belegen: bei funktionalen
Anforderungen Happy, Boundary und Negative, bei Qualitätsanforderungen die
Messung. Die Spalte nennt die beteiligten Nachweisarten. Diese Datei liest
`make doc-trace` (`trace.coverage`); Teilabdeckung steht in
`abdeckung-gesamt.md`.

| Anforderung | Nachweisarten |
| --- | --- |
| [`LH-FA-02`](../../spec/lastenheft.md) | E2E, Unit |
| [`LH-FA-05`](../../spec/lastenheft.md) | E2E, Unit |
| [`LH-FA-07`](../../spec/lastenheft.md) | E2E, Unit |
| [`LH-FA-13`](../../spec/lastenheft.md) | E2E |
| [`LH-QA-06`](../../spec/lastenheft.md) | Unit |
