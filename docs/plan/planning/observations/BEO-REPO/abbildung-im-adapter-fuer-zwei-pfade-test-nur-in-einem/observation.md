# Eine Abbildung im Adapter gilt für zwei Pfade und ist nur für einen getestet

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die Abbildung von `connection_closed` im PGWire-Adapter (`send` lässt die letzte Antwort aus,
der Aufrufer schließt) ist gemeinsamer Code für den Replay-Pfad und den Record-Pfad. Geliefert und
getestet wird sie in einem Slice für den Replay-Pfad; der Record-Pfad hat bis zu seinem eigenen Slice
keinen Test dafür, und eine Änderung der Aufteilung träfe beide.

Abgrenzung zu `BEO-REPO/fall-tabelle-je-format-doppelt`: Dort steht dieselbe Fallmenge je Format in
zwei Tabellen. Hier steht eine Funktion, die zwei Pfade benutzen, und ein Pfad prüft sie.
