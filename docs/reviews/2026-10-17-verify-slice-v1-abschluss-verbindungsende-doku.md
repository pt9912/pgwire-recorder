# Verifikation `slice-v1-abschluss-verbindungsende-doku` (Schritt 42)

Rolle: Verifier (Modul 11). Diff `e9f61167..3d3d31ca`, Plan `slice-v1-abschluss-verbindungsende-doku` (in `in-progress/`), Review `docs/reviews/2026-10-17-review-slice-v1-abschluss-verbindungsende-doku.md`. Kein Handbuchtext, kein Produktcode geändert.

**Ergebnis: bestanden, 0 blockierende Befunde.** Ein LOW, zwei INFO. Der Review-Befund F-1 (DoD verlangte einen nicht belegbaren Satz) ist durch die DoD-Korrektur behoben: DoD-Punkt 1 nennt `PGR-E3002` für jede andere Version als 1 und 2, §1 trägt den Ausschluss mit Grund.

## 1. DoD gegen den Diff

| DoD-Punkt | Befund | Beleg |
|---|---|---|
| 1 Handbuch | erfüllt | Neuer Abschnitt *Das Ende der Verbindung durch die Datenbank* (Aufzeichnen mit und ohne Ende, Leerlauf, Abbruch, `ERROR`, Format `version` 1 und 2 in YAML und SQLite, `replay`, `play` mit/ohne `--allow-recorded-errors`/`--compare-responses`, `PGR-E4003`, `PGR-E5004`, Zähler); Sätze zu mehreren Verbindungen und Sitzungen angepasst; überholte Sätze zu `FATAL`/`PANIC`, Fehlerantwort und Vergleich ersetzt oder eng gefasst; Fehlercode-Tabelle `PGR-E3002`/`E3003`/`E4003`/`E5004`/`E6001` und Glossar. Optionen stehen gegen `--help` (play: `--allow-recorded-errors`, `--compare-responses`, `--continue-on-error`; record: `--format`, `--output`; replay: `--fail-on-unconsumed`; keine neue Option), Codes im Katalog `internal/hexagon/model/fehler.go`. |
| 2 README | erfüllt | Ein Absatz zum Ende (`record`, `version: 2`, `replay`, `play`, `PGR-E4003`/`PGR-E5004`) und eng gefasster Satz zur Fehlerantwort; kein Zielstand, keine Chronik. |
| 3 `make gates` | durch Beleg erfüllt, Vorbehalt LOW V-1 | §7 nennt `make gates-schnell` Exit 0, geändert sind nur Pfade der Positivliste. Selbst gefahren: `make gates-schnell` grün, Stand `.harness/state/gates-schnell-passed.diffsha` gleich dem Hash des Arbeitsbaums (`tools/harness/working-tree-hash.sh`): Stop-Hook gibt frei. |
| Review, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen (Closure) | Review liegt vor; §7 endet mit „Wird bei Closure gefüllt“. |

## 2. Entscheidungs-Konformität (`AGENTS.md` §3.11)

- Suche im Diff (Handbuch und README, nur Zugänge) nach `slice|welle|ADR|LH-|SPEC-|lastenheft|spezifikation|review|noch nicht|kommt mit|künftig|geplant|früher|seit|ältere`: kein Treffer (die Treffer „bis zu deren `Sync`“ und „bis zur Fehlerantwort“ sind Sachaussagen). Keine Chronik, kein Zielstand, kein Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen, Reviews.
- Bindung aus `-upstream-extended`: Suche nach `erweitert|Extended|Sync|Flush|herunterfahr|nur noch` in den neuen Zeilen: nichts zu Extended-Aufzeichnung mit Ende gegen PostgreSQL, nichts zu Herunterfahren mit Ende, kein „nur noch Schließen“. Eingehalten.
- Der Satz „ältere Programmversionen lehnen Version 2 ab“ fehlt, wie DoD und §1 jetzt sagen; die Zeile `PGR-E3002` sagt nur „andere Version als 1 und 2“ (am Binary: `Version 3 wird nicht unterstützt`, Exit 3).
- Spezifikation (`LH-FA-25.a`, `LH-FA-20.a`, `LH-FA-24.a`, Szenarien 16 und 18): das Handbuch beschreibt das Binary, das sich mit der Spezifikation deckt (Ende wird aufgezeichnet, wiedergegeben, beim Einspielen mit Option oder gleichem SQLSTATE im Vergleich erwartet, sonst `PGR-E4003`/`PGR-E5004`; Netzwerkabbruch bleibt `PGR-E4003`); keine Abweichung gefunden.

## 3. Stichprobe am gebauten Binary (selbst gefahren)

Image `pgwire-recorder:dev` (`facf1c729079`, wie in §7), `postgres:17-alpine` per Digest aus `harness/mk/integration.mk` (17.11), eigenes Netz, ein Lauf zur Zeit; Container und Netz danach entfernt (`docker ps -a` leer, Netz weg). Hinweis: Der erste Lauf traf ein altes Scratch-Verzeichnis (`PGR-E2002`, Datei existiert) und ist verworfen; alle unten genannten Ergebnisse stammen aus Läufen mit frischem Verzeichnis.

| Zeile §7 | Probe | Beobachtet | Handbuch |
|---|---|---|---|
| 1, 4, 5, 6, 12 Aufzeichnen mit Ende | `record`, `pg_terminate_backend`, `select 2`, `docker stop` | `psql` 1: `FATAL:  terminating connection due to administrator command` + vier Zeilen, Exit 2; `psql` 2: Zeile `2`, Exit 0; `record` Exit 0, Log ohne `PGR-`; Datei `version: 2`, ein `connection_closed` (Zeile 51) | stimmt |
| 14, 15 Wiedergeben | `replay --fail-on-unconsumed`, dieselben zwei Aufrufe | gleiche Meldung, Exit 2; dann `2`; `replay` Exit 0, kein `PGR-W2001` | stimmt |
| 16 `PGR-E4003` | `play` ohne Option | `Session 1, Interaktion 1: Fehlerantwort des Servers 57P01 …, die Verbindung endet`, `eingespielt=1`, Exit 4 | beide Zeilen wörtlich bis auf Zeitstempel |
| 17 mit Option | `play --allow-recorded-errors` | `eingespielt=2 verglichen=0 abweichend=0`, Exit 0 | stimmt |
| 19, 18 `PGR-E5004` | Aufnahme mit `idle_session_timeout` (Datei: `C: 57P05`, `connection_closed`, `version: 2`), `play --compare-responses` | `Position 1, Art Nachrichtenart`, `eingespielt=2 verglichen=2 abweichend=1`, Exit 5; ohne Vergleich Exit 0 | stimmt |
| 20 `PGR-E3002`/`E3003` | Kopien von `ende.yaml` mit Version 3 und 1 | `Version 3 wird nicht unterstützt`, Exit 3; `connection_closed in Version 1`, Exit 3 | stimmt |
| 9 Abbruch | `docker kill` des Servers in `pg_sleep(30)` | `PGR-E4003`, `record` Exit 4, `psql`: `FATAL:  Netzwerk [PGR-E4003]` | stimmt |

Marke (§3.14): `git diff --numstat e9f61167 HEAD -- docs/user README.md test/integration` = 9+4 und 244+30, Summe 287 bei Marke 600, Meldeschwelle 480 nicht erreicht. Stimmt mit §7.

## 4. Quantoren und Zahlwörter

Stichprobe der Wörter „jede“, „nie“, „nur“, „auch“, „gleich welcher SQLSTATE“, „für `FATAL` wie `PANIC` und für beide Formate“, „mehrere Sitzungen mit einem Ende“, „andere Verbindungen laufen weiter“, „auch `--continue-on-error` ändert das nicht“, „keine Wirkung“: je eine Zeile in §7 (4, 6, 7, 14, 16 bis 19), im Bestand durch die Tabelle *Durchsicht des Bestands* mit Entscheidung je Treffer. `PANIC` und `ERROR` mit Abbruch sind im Handbuch ausdrücklich als nur mit Ersatz-Server geprüft gekennzeichnet; PostgreSQL 17.11 als einzige Serverversion ist genannt (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`). Kein Quantor ohne Beleg gefunden.

## Findings

| ID | Klasse | Befund | Zuständigkeit |
|---|---|---|---|
| V-1 | LOW | Der Sensor-Beleg für DoD-Punkt 3 ist `make gates-schnell`, nicht `make gates`; das ist nach der Positivliste des Stop-Hooks zulässig (nur `docs/user/`, Wurzel-`README.md` und `docs/plan/planning/` geändert, Hash bestätigt), die DoD-Zeile sagt aber `make gates`. Kein Befund am Stand, nur die Lesart in der Closure-Notiz benennen. | Planner (Closure-Notiz) |
| V-2 | INFO | Das Kopfdatum des Handbuchs `Stand: 16.10.2026` liegt vor dem Tag des Abschlusses (Review und DoD-Korrektur 2026-10-17); der Inhalt änderte sich seit dem 16.10. nicht, kein sachlicher Fehler. | Implementer nach Ermessen |
| V-3 | INFO | Die Beispiele des neuen Abschnitts liegen nicht als Datei in einem Test (Plan §3 entscheidet das so, Review F-2); belegt durch den Handlauf des Implementers und diese Stichprobe. Eine spätere Binary-Änderung färbt kein Gate. | zur Kenntnis (§1 schließt einen Sensor aus) |

## Vorschlag Risiko-Ausgänge (für die Closure)

1. Extended-Aufzeichnung mit Ende gegen PostgreSQL im Handbuch zugesagt: *entfallen*, Suche im Diff ohne Treffer, Bindung in §1 eingehalten.
2. Zusage über PostgreSQL 17 hinaus: *entfallen*, das Handbuch nennt 17.11 und `psql` 17.11 als geprüft, `PANIC` und `ERROR` mit Abbruch als nur Ersatz-Server; zusätzlich *weiter offen* für `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (Zähler 2×, falls der Planner die Wiederholung zählt).
3. Handbuch-Lücke zu `FATAL`/`PANIC` bei `--compare-responses`: *eingetreten* und mit diesem Slice ersetzt (Sätze am Binary bestätigt), Zähler `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`.

## Gates

`make docs-check`: 0 Befunde (vor dem Commit des Reports); `make kopf-check`: Exit 0; `make gates-schnell`: Exit 0, Stand gleich Hash des Arbeitsbaums. Arbeitsbaum vor dem Report sauber.
