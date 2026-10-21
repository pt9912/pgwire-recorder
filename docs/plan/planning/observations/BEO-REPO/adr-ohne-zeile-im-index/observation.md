# Eine neue ADR hat keine Zeile in der Haupttabelle des Index

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine ADR liegt als Datei unter `docs/plan/adr/`, aber der Index `docs/plan/adr/README.md`
führt sie nicht in seiner Haupttabelle (oder nur in der Tabelle der Beziehungen, oder mit einer
Zeile, deren Spalten die Tabelle brechen). `make docs-check` meldet das nicht: Es prüft, dass
eine Kennung `ADR-NNNN` als Link auf `docs/plan/adr/` auflöst, nicht dass jede ADR-Datei eine
Zeile im Index hat. Die Regel dazu steht als Prosa (`AGENTS.md` §5 Zeile 3,
`harness/README.md` §Traceability rules), ohne Sensor; wer den Index liest, hält ihn für
vollständig.

Abgrenzung zu `BEO-REPO/historie-zeile-nennt-neuen-satz-der-spezifikation-nicht`: Dort fehlt
der **Historie** eines Vertrags die Zeile zu einer Änderung. Hier fehlt dem **Index** die Zeile
zu einem ganzen Artefakt.

## Benannt, nicht gezählt

- 2026-10-21, Befund des Nutzers: Die Zeilen für die ADRs 0051 bis 0055 fehlten in der
  Haupttabelle, und zwei Zeilen der Tabelle der Beziehungen waren als Tabelle zerbrochen (eine
  ohne abschließendes `|`, eine mit einer Spalte zu viel). Behoben im Commit `4d9fd4ae`; ein
  voller Lauf der Gate-Kette war davor grün. Kein abgeschlossener Vorgang trägt den Fund — er
  fiel zwischen zwei Vorgängen auf —, deshalb keine Datei in `evidence/` und kein
  Zählerzuwachs.
