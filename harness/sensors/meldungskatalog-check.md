# `make meldungskatalog-check` — gleicht die Meldungscodes aus `fehler.go`, `SPEC-034` und dem Benutzerhandbuch ab

---

## Vertrag

Rot heißt: Die Codes `PGR-…` der Code-Tabelle (`internal/hexagon/model/fehler.go`), des Katalogs
(Tabelle in `SPEC-034`) und der Tabellen `### Fehlercodes` und `### Warnungen` des Benutzerhandbuchs
verletzen [`SPEC-054`](../vertraege/SPEC-054-meldungskatalog.md#spec-054--abgleich-von-code-tabelle-und-meldungskatalog-meldungskatalog-check):
ein Code hat die falsche Form, steht in einer Quelle doppelt, fehlt in einer der anderen, ein Rückfall
`PGR-E<k>000` fehlt im Katalog, oder eine Go-Datei außer der Code-Tabelle trägt ein Literal `"PGR-`.
Was als Quelle, Code und Befund gilt, sagen die Punkte 1 bis 7 des Vertrags; diese Datei entscheidet
nichts davon.

## Grenze — was das Grün nicht abdeckt

Die Grenze von `SPEC-054`; jede ist ein akzeptiertes Negativ und permanent, solange sie gilt.

1. **Spalten** — *Klasse / Bereich*, *Bedeutung* und *Maßnahme* sind nicht gelesen; ob die erste Ziffer
   zur Klasse passt oder der Text stimmt, sagt der Vergleich nicht.
2. **Andere Fundorte** — Codes in Kommentaren, Prosa, Tests, `test/`, Plänen und Lastenheft sind nicht
   gelesen, ebenso ein Code in Backticks im Quelltext.
3. **Verhalten** — der Exit Code und der Text eines Codes im Lauf des Binaries; das prüfen die Tests
   des Produkts.
4. **Handbuch** — ob eine Zeile sinnvoll beschreibt, was der Code meldet.
5. **Codeblöcke** — eine Zeile `#` in einem Codeblock beendet einen Abschnitt; was danach in ihm steht,
   ist nicht gelesen.

## Ausgabe und Ausgänge

Je Befund eine Zeile auf stderr, `meldungskatalog-check: <pfad>[:<zeile>]: <code>: <befund>`, sortiert
nach Pfad, Code und Befund, danach nach der Zeile (Punkt 9 des Vertrags); stdout bleibt leer.

| Exit | Bedeutung |
|---|---|
| 0 | kein Befund; keine Ausgabe |
| 1 | mindestens ein Befund; alle Quellen sind gelesen, jede Befund-Zeile steht da |
| 2 | eine Sperre; nichts abgeglichen, eine Zeile `meldungskatalog-check: <grund>` |

`make` meldet jeden Ausgang ungleich 0 als Fehler. Auseinander hält die Ausgänge die Zeile auf stderr:
eine Befund-Zeile nennt einen Pfad und einen Code, eine Sperre nur den Grund. Aus dem Rot führt:

- `nicht im Katalog` oder `nicht im Handbuch` → den Code in `SPEC-034` bzw. in die Tabelle des
  Handbuchs eintragen (Handbuch nur, wenn der Code erzeugt wird);
- `nicht in der Code-Tabelle` → die Konstante in `fehler.go` ergänzen oder die Zeile in `SPEC-034`
  bzw. im Handbuch streichen; ein nicht erzeugter Rückfall `…000` steht nur im Katalog;
- `Code außerhalb der Code-Tabelle` → das Literal durch die Konstante aus `fehler.go` ersetzen;
- `Rückfall fehlt` → `PGR-E<k>000` für die Klasse in `SPEC-034` eintragen;
- `Form` und `doppelt` → die Zeile berichtigen.

## Sperren

| Grund | Abhilfe |
|---|---|
| `<pfad> fehlt`, `<pfad> nicht lesbar` | `make meldungskatalog-check` in der Wurzel des Repos starten; Datei lesbar machen |
| `cmd/ fehlt` | in der Wurzel des Repos starten |
| `<pfad> nicht lesbar` für ein Verzeichnis oder eine Go-Datei unter `cmd/` und `internal/` | Rechte herstellen; gemeldet wird der erste Pfad in Bytefolge |
| `<pfad> ist ein Symlink` für einen Symlink unter oder an `cmd/` und `internal/` | den Symlink durch die Datei oder das Verzeichnis ersetzen; gemeldet wird der erste Pfad in Bytefolge, nach „nicht lesbar“ |
| `Abschnitt <überschrift> fehlt` | die Überschrift `### SPEC-034 `, `### Fehlercodes` bzw. `### Warnungen` wiederherstellen |
| `<pfad> ohne Code`, `<pfad> ohne Codes` | die Quelle trägt keine Tabellenzeile bzw. keine Konstante in der Form des Vertrags |

Die Reihenfolge der Sperren nennt Punkt 8 des Vertrags; die erste gilt.

Grün heißt auch bei `make meldungskatalog-check-gegenprobe`: jeder Fall, der lief. Als root sind die Fälle mit nicht lesbaren Pfaden, ohne die Locale `de_DE.UTF-8` der Fall `nr9-bytefolge` nicht herstellbar; sie werden auf stderr als übersprungen gemeldet, und der Lauf endet dennoch mit Exit 0.

## Bindung

[ADR-0011](../../docs/plan/adr/0011-meldungscodes-praefix-pgr.md) · Vertrag
[`SPEC-054`](../vertraege/SPEC-054-meldungskatalog.md#spec-054--abgleich-von-code-tabelle-und-meldungskatalog-meldungskatalog-check)

`make meldungskatalog-check` fährt in der schnellen Kette (`make gates-schnell`), weil die Positivliste
von [`SPEC-052`](../vertraege/SPEC-052-gate-kette.md) Änderungen am Handbuch ohne vollen Lauf zulässt.
Seine Gegenprobe `make meldungskatalog-check-gegenprobe` hängt nicht an der Gate-Kette von `make gates`;
sie steht in `GEGENPROBEN_CONTAINER`, prüft, ob das Werkzeug richtig ist, und ihr Vertrag
steht in ihrer Zeile in [`../README.md` §Sensors](../README.md#sensors-feedback-gates).
