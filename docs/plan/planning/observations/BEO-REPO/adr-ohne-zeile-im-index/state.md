**Stand:** offen

Zähler 0; der einzige bekannte Fund steht in `observation.md` unter *Benannt, nicht gezählt*
(2026-10-21, Commit `4d9fd4ae`). Kein Ausgang: Unterhalb der Schwelle ist `offen` der
Normalzustand, und ein Ausgang wird ab 3× zugewiesen.

Angelegt vom Planner am 2026-10-21 beim Herrichten von `welle-dev-container` (Auftrag des
Nutzers vom selben Tag, den Befund als Register-Eintrag zu führen).

## Prüfung: gestrichen (Entscheidung des Nutzers vom 2026-10-21)

Der Nutzer hat am 2026-10-21 entschieden (wörtlich, vom Koordinator übermittelt: „ok deine
empfehlugnen“ auf eine Vorlage mit sechs Punkten, deren vierter diese Frage war): Eine Prüfung
entsteht **nur**, wenn `.d-check.yml` sie mit einer Regel tragen kann, sonst keine. Der Planner
hat das am 2026-10-21 am gepinnten Image gemessen (Digest aus `d-check.mk`, `--manual`,
`--print-config` und Läufe gegen Temp-Bäume); sie kann es **nicht**:

- Das Modul `mentions` (`DC-FA-MENT-001.a`) ist die nächste Form: Es prüft je Mitglied einer
  Artefakt-Menge, dass es in irgendeinem Dokument einer Dokument-Menge wörtlich genannt ist.
  Mit `artifacts: ["docs/plan/adr/[0-9]*.md"]`, `documents: ["docs/plan/adr/README.md"]` und
  `match: basename` läuft die Regel: auf dem heutigen Bestand grün (57 von 57 Artefakten
  erwähnt), und nach dem Entfernen der Index-Zeilen einer ADR rot mit `artifact-unmentioned`.
- Sie trägt den Befund dieses Eintrags aber nicht. Gemessen auf genau dem Stand, der ihn trug
  (`4d9fd4ae^`, Haupttabelle ohne die ADRs 0051 bis 0055): **grün, 55 von 55**. Die fehlenden
  Zeilen der Haupttabelle fielen nicht auf, weil die Tabelle der Beziehungen dieselben Dateien
  nennt. `mentions.documents` ist eine Datei-Menge; einen Abschnitt oder eine Tabelle darin kann
  die Regel nicht adressieren.
- Das Modul `structure` adressiert einen Abschnitt und dessen Tabelle, misst aber nur Zellen
  **vorhandener** Zeilen (Länge je Spaltenname, Ordnung der Schlüsselspalte) — eine fehlende
  Zeile hat keine Zelle. Die zweite Hälfte des Befundes (eine Zeile, die die Tabelle bricht)
  fängt `structure` über `table.column` mit `cell-min-chars`; die erste Hälfte nicht.
- Kein anderes Modul des Images bindet eine Artefakt-Menge an einen Abschnitt (`links`,
  `anchors`, `ids`, `matrix`, `citations`, `reviews`, `file`, `targets`, `planning`, `tracked`,
  `commits`, `vcs`, `workflows`).

**Die Lücke bleibt damit ungeprüft:** `make docs-check` meldet weiterhin nicht, wenn eine
ADR-Datei keine Zeile in der Haupttabelle des Index hat; die Regel dazu bleibt Prosa
(`AGENTS.md` §5 Zeile 3, `harness/README.md` §Traceability rules).

**Gestrichen ist der Sensor, nicht die Beobachtung.** Der Ausgang `gestrichen` der drei Ausgänge
(`../../README.md` §Die drei Ausgänge) verlangt die Begründung, warum die Beobachtung **nicht
mehr auftreten kann**; sie kann weiter auftreten. Der Eintrag bleibt deshalb `offen` und zählt
weiter. Erreicht er 3×, bekommt er einen Ausgang, und die Mittel sind dann ein eigenes Skript
mit Vertrag und Gegenprobe oder eine Regel einer neueren Fassung von d-check, die einen
Abschnitt adressiert. Ein Slice dafür ist **nicht** angelegt: Das ist Entscheidung des Nutzers,
und ein `geplant` mit einer Kennung, die es nicht gibt, wäre keine Adresse (`AGENTS.md` §3.13).

Die Zeile für [ADR-0056](../../../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) steht seit dem Proposed-Commit `85454f8e` in beiden Tabellen des Index; der Architect hat Status und Datum am 2026-10-22 gesetzt (`71255295`). `slice-harness-groessen-metrik` trägt dazu nichts mehr (sein §6: Risiko entfallen); ein Gegenbeleg aus dem Slice entfällt, der Zähler bleibt 0.

Erstes gezähltes Auftreten: `slice-harness-gegenproben-schneller` (Review r1 F-6, LOW) — die Hauptzeile von [ADR-0059](../../../../adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md) fehlte in der Haupttabelle (Beleg `evidence/slice-harness-gegenproben-schneller.md`); nachgetragen. Zähler 1, unter der Schwelle, `offen`.
