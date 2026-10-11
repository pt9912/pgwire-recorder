# Review-Report: slice-v1-abschluss-sqlite-format-doku — 2026-10-11

**Review-Art:** Doku-Review. Geprüft wird gegen Plan §1, §3, §6 und §8, gegen `AGENTS.md` §3.9, §3.11, §3.13, gegen [LH-FA-22](../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) und [ADR-0013](../plan/adr/0013-sqlite-recording-backend.md), nicht gegen die DoD.

**Gegenstand:** `221ee96..3234506` (zwei Commits, +220 −26 über vier Dateien: `docs/user/benutzerhandbuch.md`, `README.md`, der Slice-Plan, die Roadmap-Zeile). Kein Produkt-Code, kein Test, kein Gate im Diff.

**Skill:** `.harness/skills/reviewer.md` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

**Eingangs-Kontext:** Plan ganz gelesen mit §7 *Belege des Implementers* (der Bericht des Implementers liegt nicht vor); `AGENTS.md` §3; Report `2026-10-11-review-slice-v1-abschluss-sqlite-format.md` und `…-lesen.md` zum selben Modul (Nummern dieses Laufs ab F-631).

**Ausgeführte Läufe** (Docker-only; `make build` am Stand `3234506`, Image `d221c756b469`; Scratch außerhalb des Repos; eigenes Netz und ein PostgreSQL-17-Container aus dem gepinnten Image von `harness/mk/integration.mk`, danach per Name entfernt):

- `record` mit `--format sqlite` und yaml, Ablauf einfach und erweitert (`\bind`); `replay` aus beiden Dateien liefert die Antworten der Datenbank. Tabellen und Spalten der Datei stimmen mit dem Abschnitt *Eine SQLite-Aufzeichnung abfragen* überein (`interaction.sql` bei `extended` leer, `message.kind = 'parse'` trägt das SQL); beide `sqlite3`-Aufrufe liefen wörtlich.
- `--format`: `SQLite`, `sqlite3`, leer als Option ergeben `PGR-E2001`; Variable `SQLite` mit gültiger Option `PGR-E2001`; `format: SQLite` in der Datei `PGR-E2004`; leere Variable gilt als nicht gesetzt; `replay` und `play` mit `--format`: `PGR-E2001`.
- `PGR-E2002` bei sqlite auf vorhandene sqlite- und yaml-Datei und umgekehrt; leeres `<Ziel>-journal`: `PGR-E3001` auch mit `--force` und ohne Zieldatei, mit Zieldatei ohne `--force` zuerst `PGR-E2002`; bei yaml kein Befund.
- `--force` mit beiden Formaten: Datei vor der ersten beendeten Verbindung byte-gleich, danach ersetzt, nie ergänzt.
- Sperre mit `begin exclusive` für 14 s: `replay` endet nach 5,5 s mit `PGR-E3001` (`SQLITE_BUSY`). `replay` und `play` auf `:ro`: Prüfsumme gleich, kein Journal. Lesen während des Aufzeichnens: `replay` und `sqlite3` sehen die beendete Verbindung. Aufzeichnung ohne Verbindung: gültige Datei ohne Sitzung, `replay` `PGR-E3004`, `play` Exit 0 ohne Meldung. Datei leer und 15 Byte `SQLite format 3`: `PGR-E3003`.
- Journal-Anleitung: ein 0-Byte-Journal bleibt nach einer Abfrage mit `sqlite3` liegen (F-631).
- Commit-Messages: keine `SPEC-`/`ARC-`-Kennung. `grep` im Handbuch auf ADR, SPEC, LH, Slice, Welle, Spezifikation, Lastenheft, Review: kein Treffer; Chronik-Wörter („noch nicht“, „bisher“, „kommt mit“, „geplant“): kein neuer Treffer.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-631 | MEDIUM | Der Hinweis zum Journal führt jedes `<Zieldatei>-journal` auf „Rollback-Journal eines früheren Laufs“ zurück und verspricht: Eine Abfrage mit `sqlite3` macht den Abbruch rückgängig und entfernt die Datei, danach startet `record` wieder. Ein Journal der Länge 0 (Prozess hart beendet zwischen Anlegen und erstem Schreiben des Journals, oder eine fremde Datei dieses Namens) bleibt nach der Abfrage liegen; `record` meldet weiter `PGR-E3001`, und der Handbuchtext nennt keinen Ausweg. Die Probe des Implementers (§7) fährt nur das Journal mit Inhalt (9728 Byte). | AGENTS.md §3.11 (Kettenaussage „macht … rückgängig und entfernt“, nur in einer Richtung geprobt) | `docs/user/benutzerhandbuch.md` · „dann macht SQLite den Abbruch rückgängig und entfernt die Datei; danach startet `record`“ (Hinweis zum Journal im Abschnitt *Das Format der Aufzeichnung festlegen*) | ja (Kopie der Datei, 0-Byte-Journal, `sqlite3 … "select count(*) from session"`, `ls`) | Kettenaussage nur in einer Richtung geprobt |
| F-632 | LOW | Die Zeile `PGR-E3001` in §7 sagt zu, die Meldung nenne bei einer gesperrten Datei `SQLITE_BUSY`. `TestLesenSperrfrist` prüft nur, dass die Meldung „locked“ enthält; dieselbe Tabelle sagt oben, nur Klasse und Code seien verlässlich, der Text danach könne sich ändern. Der Wortlaut folgt der Bibliothek und ist über keinen Test zugesagt. | AGENTS.md §3.11 (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`); Wiederholung von F-615 | `docs/user/benutzerhandbuch.md` · „(die Meldung nennt `SQLITE_BUSY`)“ | ja (Test auf den Wortlaut oder Satz enger fassen) | Zusage im Handbuch weiter als der Test |
| F-633 | LOW | „Eine Datei, die nicht mit der Kennung einer SQLite-Datei beginnt, liest es als YAML, auch eine leere; sie ist dann beschädigt (`PGR-E3003`).“ „Sie“ ist nicht eindeutig: Gelesen auf die ganze Datei sagt der Satz, jede Datei ohne Kennung (auch eine gültige YAML-Aufzeichnung) sei beschädigt. Gemeint ist nur die, die als YAML nicht lesbar ist. | Maintainability (Verständlichkeit, AGENTS.md §3.11 Ist-Zustand) | `docs/user/benutzerhandbuch.md` · „sie ist dann beschädigt (`PGR-E3003`)“ | nein | Satzbezug mehrdeutig |
| F-634 | INFO | Hinweis (a) des Implementers: Das Handbuch sagt nichts Falsches. `PGR-E3004` nennt ausdrücklich nur `replay` („meldet den Fehler beim Start“), und der Abschnitt zu `play` sagt über eine Aufzeichnung ohne Sitzung nichts. Probe: `play` Exit 0, nur „play gestartet“ und „play beendet“. Ob das Register einen Eintrag braucht, entscheidet der Planner. | Maintainability | `docs/user/benutzerhandbuch.md` · „`replay` meldet den Fehler beim Start“ | ja | Verhalten ohne Meldung bei leerer Aufzeichnung |
| F-635 | INFO | Hinweis (b): Das Handbuch sagt zum Inhalt der Meldungen bei einer kaputten Datei nichts zu (nur Code und Ursache-Klasse), außer der Stelle in F-632. Die 15-Byte-Datei meldet `PGR-E3003` mit einem Stück des Dateiinhalts (`cannot unmarshal !!str …`); das Handbuch sagt dazu nichts. | Maintainability | `docs/user/benutzerhandbuch.md` · Zeile `PGR-E3003` | ja | Meldung nennt Dateiinhalt (bekannte Lücke, kein Befund dieses Diffs) |
| F-636 | INFO | Der Satz „Neben der Zieldatei liegt während des Ergänzens eine Datei `<Zieldatei>-journal`“ hat in der Probe-Matrix (§7) keine Zeile; die Matrix prüft nur, dass nach dem Lesen kein Journal bleibt. Der Abschnitt führt zwölf Hinweise, die Nutzung und Randfälle mischen; die Lesbarkeit hängt am Aufbau, nicht an einer Regel. | AGENTS.md §3.11 | `docs/user/benutzerhandbuch.md` · „liegt während des Ergänzens eine Datei“ | nein | Satz ohne eigene Probezeile |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/user/benutzerhandbuch.md` Tabelle §5 (`--format`, Variable, Standard, Werte) gegen `record --help` und Binary | geprüft, ohne Befund |
| Fehlercode-Tabelle §7 (`PGR-E2001`, `PGR-E2004`, `PGR-E3002`, `PGR-E3003`) | geprüft, ohne Befund außer F-632 (Zeile `PGR-E3001`) |
| Priorität Option > Variable > Datei, leere Variable, ungültige Variable mit gültiger Option | geprüft, ohne Befund (Probe in beide Richtungen) |
| `--force` ersetzt nie ergänzt, Datei bis zur ersten Aktualisierung unverändert, Formatwechsel | geprüft, ohne Befund |
| `PGR-E2002` vor `PGR-E3001`, `PGR-E3001` mit `--force` und ohne Zieldatei | geprüft, ohne Befund |
| Sperre fünf Sekunden, lesend auf `:ro`, Lesbarkeit während des Aufzeichnens, Aufzeichnung ohne Verbindung | geprüft, ohne Befund |
| Tabellen und Spalten, beide `sqlite3`-Beispiele (laufen wörtlich), `interaction.sql` bei `extended` | geprüft, ohne Befund |
| Grenze „unter Linux im Container auf einem lokalen Dateisystem“: keine Aussage zu Windows, macOS, großen Dateien | geprüft, ohne Befund |
| Chronik, Zielstand, Verweise auf Spezifikation, Lastenheft, ADR, Slice, Welle, Review im Handbuch | geprüft, ohne Befund |
| `README.md` (ein Satz; Verweise auf `spec/` unverändert) | geprüft, ohne Befund |
| FAQ und Glossar (Format, Binärdatei, beide Formate) | geprüft, ohne Befund |
| Plan §3, §7, §8, DoD-Zählung (zwei Liefer-Punkte, eine Schicht), Roadmap-Zeile, `AGENTS.md` §3.9; Namensabgleich genannter Dateien | geprüft, ohne Befund |
| §3.13: kein Nehmer neu eingetragen, keine Adresse im Diff | geprüft, ohne Befund |
| Commit-Messages (`SPEC-`/`ARC-`-Kennung), §3.3 (kein Move im Diff) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Kettenaussage nur in einer Richtung geprobt · Zusage im Handbuch weiter als der Test · Satzbezug mehrdeutig · Verhalten ohne Meldung bei leerer Aufzeichnung · Satz ohne eigene Probezeile

## Verdikt

**Merge-blockierend:** nein. Keine HIGH-Findings. F-631 (MEDIUM) ist ein Handbuchsatz, der in einem Randfall (Journal der Länge 0) in eine Sackgasse führt.

**Übergabe:** F-631 bis F-633 an den **Implementer** (Satz enger fassen oder Fall aufnehmen, `AGENTS.md` §3.11); F-634 an den **Planner**; F-635 und F-636 ohne Blockade. Die Finding-Klassen gehen in die Slice-Closure §7. DoD-Konformität prüft der Verifier.
