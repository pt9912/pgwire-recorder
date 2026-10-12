# Verifikation: slice-v1-abschluss-antwortvergleich — 2026-10-12

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-12-review-slice-v1-abschluss-antwortvergleich.md`).

**Gegenstand:** Slice-Plan `slice-v1-abschluss-antwortvergleich` (Kopf, DoD, §1, §3, §4, §6, §7) gegen `08787ed..c03f502`. Commits: `6ef3cf6`, `79a380c`, `b0122fd`, `25ccf77`, `8e8b60d`, `d99f24e` (Implementer), `d6d7214` (Review-Report), `f15305e`, `7da285c`, `c03f502` (Nacharbeit F-665 bis F-668); Architect-Entscheidungen `d8b8611`, `35f1ec0`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

**Eingang:** Plan ganz (§7 *Belege des Implementers* und *Belege der Nacharbeit*), Diff von `internal/`, `test/`, `docs/user`, `README.md`, `LH-FA-24.a` (mit *Randformen des Vergleichs*), `LH-FA-21.a`, `LH-FA-20.a` (Fehlertabelle), Review-Report. Den Bericht des Implementers habe ich nicht gelesen.

## 1. Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make docs-check` | 637 Dateien, 0 Befunde |
| `make kopf-check`, `make abdeckung-check` | Exit 0 |
| `make build` | Image `pgwire-recorder:dev` gebaut |
| `make test-integration` (volle Suite, Stand `c03f502`, Arbeitsbaum sauber) | `EXIT 0`, 107 `--- PASS`, kein `FAIL`; darin die zehn Tests `TestE2EPlayVergleich*` (`Gleich`, `Abweichung` mit vier Unterfällen, `Fortsetzung`, `Fatal`, `Quellen`, `Extended`, `ExtendedFehler`, `Hinweise`) grün gegen die reale Datenbank |
| Proben am gebauten Binary | alle wie spezifiziert (Tabelle unten) |
| 11 Mutationen per Bind-Mount in `pgwire-recorder:test` auf einer Kopie des Baums im Scratchpad | 10 rot, 1 grün (V5b, Befund V-1); danach `diff -r internal` gegen das Repo: identisch, `git status` leer |
| `grep -in compare` in `docs/user/benutzerhandbuch.md` und `README.md` | kein Treffer (Ist-Zustand, Doku-Folge-Slice trägt es) |

Nicht gefahren: `make gates` (volle Kette), `make test`, `make lint`, `make a-check`, `lint-gegenprobe`, `a-check-negativ`; die übrigen 60 Mutanten der Tabelle in §7 (siehe Abschnitt 6).

### Proben am Binary (ohne Datenbank, Startfehler)

| Aufruf | Ergebnis |
|---|---|
| `play --help` | nennt `--compare-responses[=true\|false]`, `PGR-E5004 (Exit-Code 5)`, Standard `false`, `PGWIRE_RECORDER_COMPARE_RESPONSES` |
| `record` / `replay` mit `PGWIRE_RECORDER_COMPARE_RESPONSES=true` | Exit 2, `PGR-E2001`, „gehört zu einer Option, die record/replay nicht kennt“ |
| `play` mit der Variable `true` | wird gelesen (nächster Fehler: `--upstream` fehlt) |
| `play` mit Variable `banane`, `1`, `True` | `PGR-E2001`, „erlaubt sind true und false“ |
| `play` mit leerer Variable | gilt als nicht gesetzt (nächster Fehler: `--upstream` fehlt) |
| `play --compare-responses=banane` | `PGR-E2001` |
| Konfigurationsdatei `play: compare_responses: true` | wird gelesen; `maybe` ist `PGR-E2004` (`play.compare_responses: erlaubt sind true und false`) |
| Konfigurationsdatei `record: compare_responses: true` | `PGR-E2004`, unbekannter Schlüssel |
| `config show` | listet `compare_responses` und `PGWIRE_RECORDER_COMPARE_RESPONSES` |
| `play --compare-responses` mit leerer Aufzeichnung | `PGR-E3003` (Aufzeichnung wird vor der Verbindung geprüft) |

### Mutationen (selbst gefahren)

Weg: Kopie von `cmd`, `internal`, `go.mod`, `go.sum` im Scratchpad, Bind-Mount nach `/src` in `pgwire-recorder:test`, `go test` über die Pakete des Mutanten, Datei danach zurückgesetzt. Basislauf der Pakete `services` und `bootstrap` grün.

| # | Zusage | Mutation | Ergebnis (roter Test) |
|---|---|---|---|
| V1 | verneinend: nach `FATAL` liest `play` nichts mehr (M35) | `fatal()` gerufen, Fehler verworfen, Lesen geht weiter | rot: `TestVergleichFatal`, `…FatalRand`, `…FatalFortsetzung`, `…FatalExtendedGruppen` |
| V2 | verneinend: keine zweite Zeile `play beendet` (C9) | zweite `log.Info("play beendet")` davor | rot: `TestRunPlay`, `…Fehlerantwort`, `…Vergleich`, `…SignalImStart` |
| V3 | mehrteilig: Fehler vergleicht SQLSTATE, nicht den Schweregrad (M18) | `|| V` verschieden | rot: `TestVergleichStruktur` |
| V4 | mehrteilig: `FATAL` und `PANIC` (M43) | nur `FATAL` | rot: `TestPlayFehlerantwort`, `TestVergleichFatal` |
| V5 | nur Wörter aus Ziffern sind Zahlen (M7) | Wort mit führender Ziffer gilt als Zahl | **grün** (mein Mutant ist für `T1` äquivalent; siehe V5b) |
| V5b | dieselbe Zusage | `strings.ContainsAny(w, "0123456789")` statt „nur Ziffern“ | **grün**, siehe V-1 |
| V6 | Abbruch vor der nächsten Interaktion ohne `--continue-on-error` (M27) | Abbruch-Bedingung `if false` | rot: `TestVergleichAbbruchNachInteraktion` |
| V7 | Position ab 1 (M2) | `u.position = i` | rot: zwölf Tests |
| Z1 | Vergleich ändert die Wartezeit nicht (F-665) | `warte` entfällt mit `CompareResponses` | rot: `TestPlayZeitMitVergleich` |
| E2 | Zeile zum Ende erscheint auch ohne Interaktion (F-666) | Zeile nur bei `Eingespielt != 0` | rot: `TestRunPlayVergleich`, `TestRunPlaySignalImStart` |

Die Tabelle in §7 führt 70 + 8 Mutanten; die Zeile „M7 … `TrimRight` über Ziffern und Leerzeichen“ ist rot (laut Implementer), die Zusage „nur Wörter aus Ziffern“ ist damit nur gegen diese eine Mutation geprüft.

## 2. DoD Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1: gleiche Antwortstruktur Erfolg, abweichende `PGR-E5004`/Exit 5; Serverfehler ohne Vorbild `PGR-E5004` statt `PGR-E4004`; `FATAL`/`PANIC` `PGR-E5004` auch bei gleichem SQLSTATE; Netzwerkabbruch `PGR-E4003`; nach `FATAL`/`PANIC`/Verbindungsverlust nichts mehr senden, sonst bis `ReadyForQuery` vergleichen und vor der nächsten Interaktion abbrechen; Extended bis `Sync` (Abnahmeszenario 16) | erfüllt | Code: `fehlerantwort`/`fatal`/`abschluss` in `play.go`; E2E `TestE2EPlayVergleichGleich`, `…Abweichung/Serverfehler_ohne_Vorbild`, `…Fatal` (`pg_terminate_backend`), `…ExtendedFehler` selbst gefahren grün; Mutanten V1, V4, V6 rot. Dritter Satz von Szenario 16: siehe Abschnitt 5, offen beim Nutzer |
| 2: jede Art der Abweichung erkannt, gemeldet mit Session, Sequenznummer, Position (ab 1, erste ohne Gegenstück), Art, SQLSTATE beider Seiten nie Text; Unterschiede in Zeilenwerten, Zeilenzahlen, Hinweisen, abschließenden Zahlen nicht (einfach, Extended, `Flush`-Gruppen) | erfüllt, mit V-1 | `play_vergleich.go` liest sich Zeile für Zeile gegen die Tabelle in `LH-FA-24.a`; V3, V7 rot; `TestE2EPlayVergleichHinweise`, `…Extended` grün; eine Teil-Zusage ist nur gegen eine Mutation belegt (V-1) |
| 3: `--continue-on-error`, `--allow-recorded-errors` ohne Wirkung mit Vergleich, Exit 5 am Ende, Verbindungsfehler 4, Signale, Zeile `play beendet` mit `eingespielt`, `verglichen`, `abweichend` und keiner zweiten Zeile; ohne Option kein Vergleich; Leser (Kommandozeile, Variable, Schlüssel `compare_responses`), `nichtGelesen()` kürzt die Option | erfüllt | Proben am Binary (Abschnitt 1); `fremd_test.go`: `compare-responses` aus `nichtGelesen()` entfernt, `…ZuordnungGleichTabelle` auf „fünf“; Mutanten V2, V6, E2 rot; `TestE2EPlayVergleichQuellen`, `…Fortsetzung` grün |
| `make gates` grün | laut Implementer Exit 0 an `7da285c`, Auftraggeber meldet Exit 0 an `c03f502`; von mir nicht gefahren; `docs-check`, `kopf-check`, `abdeckung-check`, `test-integration` selbst grün | |
| Review durchgeführt, Report liegt vor | erfüllt | `docs/reviews/2026-10-12-review-slice-v1-abschluss-antwortvergleich.md` (0 HIGH, 1 MEDIUM, 3 LOW, 4 INFO) |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen (Closure) | die Ausgänge in §6 stehen als „offen bis Closure“ und „entfallen (Architect)“ |

## 3. Entscheidungen des Architect gegen den Code

| Entscheidung (§6, `LH-FA-24.a`) | Code / Beleg | Urteil |
|---|---|---|
| Normalisierung: `data_row`, `notice_response`, `parameter_status` entfallen | `nimmNormalisiert`; V3-nahe Tests, `…Hinweise` E2E | konform |
| Befehl: Zahlen am Ende, Ziffern-Tag leer | `befehl()`/`nurZiffern()` | konform; Teil-Zusage nur einseitig geprüft (V-1) |
| Fehler ohne `C`: leerer Text, zwei ohne `C` gleich; Schweregrad nie, außer `FATAL`/`PANIC` | `nachrichtenUnterschied`, `fehlerantwort` | konform |
| `FATAL`/`PANIC` ohne `ErrorResponse` an der Stelle: Art Fehler, aufgezeichneter SQLSTATE leer, frühere Abweichung geht vor (`35f1ec0`) | `erste(true)`: `grenze = min(len(ist)-1, len(soll))`, danach Art Fehler | konform |
| Position: erste verschiedene Position, Länge plus 1 nur bei gleichem Anfang (`35f1ec0`) | `erste()` mit Schleife vor der Längenprüfung | konform; E3 (laut Implementer) rot |
| Meldung: SQLSTATE beider Seiten, Text nie | `unterschied.text()` | konform |
| `--keep-timing` mit Vergleich: Wartezeit davor (`LH-FA-21.a`) | `warte` vor `interaktion`; `TestPlayZeitMitVergleich`; Z1 von mir rot | konform, F-665 behoben |
| Verbindungsverlust ohne Fehlerantwort `PGR-E4003`, nach Abweichung 4 | `Naechste`-Fehler unverändert durchgereicht; `TestVergleichVerbindungsverlust` | konform |
| Zähler und Zeile `play beendet` (`info`, eine Zeile, auch ohne Interaktion) | `bootstrap.go`; V2, E2 rot | konform |
| Variable beim eigenen und anderen Kommando (Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando`) | Proben am Binary; `nichtGelesen()` gekürzt; `liesUmgebung`/`zuordnung()` im Diff unverändert | Bindung erfüllt |

## 4. Review-Findings F-665 bis F-668

| Finding | Behoben? | Beleg |
|---|---|---|
| F-665 MEDIUM: Zusage `LH-FA-21.a` für `--keep-timing` mit Vergleich ohne Test | ja | `play_vergleich_zeit_test.go` (`TestPlayZeitMitVergleich`, Fake-Uhr: 999 ms keine Sendung, 1000 ms Sendung, beide Interaktionen verglichen, zweite `PGR-E5004`); Abdeckungszeile in `abdeckung-unit.md` Zeile 403; Mutant Z1 von mir rot (Z2, Z3 laut Implementer); §6 zeigt auf die Randform in `LH-FA-21.a` |
| F-666 LOW: vier Zusagen ohne Zeile in §7 | ja | Zeilen E1 bis E4b in *Belege der Nacharbeit*; E2 von mir rot |
| F-667 LOW: Deklaration `LH-FA-17/Boundary` an `TestE2EPlayVergleichGleich` | ja | `abdeckung-e2e.md`: Zeile 170 trägt `LH-FA-24` Boundary, `LH-FA-17` nicht mehr; `make abdeckung-check` Exit 0 |
| F-668 LOW: Plan folgt dem Diff nicht (§1 „bindet an `liesUmgebung`“, §3 „Abbildung auf Exit-Code 5“) | ja | §1 sagt „Eintrag in `optionen("play")`, `liesUmgebung` bleibt unverändert“, §3 nennt den Exit-Code aus der Klasse des Katalogeintrags; entspricht dem Diff (`cli.go`: ein Eintrag in `optionen("play")`, kein Eingriff in `liesUmgebung`) |

## 5. Befunde

**HIGH:** keine.

**MEDIUM:**

- V-1 (§3.10, Zusage ohne fangende Mutation): *Befehl* — „nur Wörter, die ganz aus Ziffern bestehen, sind abschließende Zahlen“. Der Test `TestVergleichStruktur/Wort_mit_Ziffer_ist_keine_Zahl` vergleicht `T1` (Aufzeichnung) mit `T` (Server) und erwartet eine Abweichung. Bei der Mutation `nurZiffern` → `strings.ContainsAny(w, "0123456789")` (jedes Wort mit einer Ziffer gilt als Zahl) wird `T1` zu leer und `T` bleibt `T`: auch dann verschieden, der Test bleibt grün (`ok  internal/hexagon/services`). Die Zeile M7 in §7 nennt nur die Mutation `TrimRight` über Ziffern und Leerzeichen. Beleg: `$ mut.sh … 'return strings.ContainsAny(w, "0123456789")'` → `ok`. Behebung: ein Fall mit zwei verschiedenen Wörtern gleicher Ziffernform (`T1` gegen `T2`, erwartet Abweichung), samt Mutationszeile. Wirkung im Betrieb gering (PostgreSQL-Tags tragen keine Wörter mit eingemischten Ziffern), aber die Regel verlangt die Mutation je Zusage.

**LOW:**

- V-2 (Größenprobe §4): §7 berichtet „Die Probe von 1400 Zeilen ist damit nicht eingetreten“, gestützt auf einen geschätzten Abzug („etwa 500 Zeilen“) für DoD-Punkt 3. Gemessen: `git diff --shortstat 08787ed..8e8b60d -- . ':!docs' ':!spec'` = 1441 Zeilen hinzugefügt, 54 entfernt; bereits `6ef3cf6` allein 1350 + 48; Stand `c03f502` 1494 + 54. Einen Commit, der nur DoD-Punkte 1 und 2 erfüllt, gibt es nicht, §4 bindet die Probe aber an diesen Stand. Der Implementer legt dies offen und begründet die Lesart; die Zahlen sind ehrlich genannt. Offen bleibt, ob der Planner die geschätzte Lesart gelten lässt oder den Wert 1441 gegen 1400 als Eintritt wertet (dann hätte der Implementer anhalten müssen). Für die Closure: Zählung nennen, ob der Slice als „dritter Beleg“ von `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` zählt (§8 sagt: tritt die Probe ein, ist es der dritte). Nicht höher gewertet: Das Gesamtmaß liegt unter der Schätzung des Architect in §8 (1800 bis 2300 Zeilen mit Tabellen).

**INFO:**

- V-3: Zusage-Prosa in `Zaehler()`/`PlayOptions` („mit ihm hat `AllowRecordedErrors` keine Wirkung“) ist durch `TestVergleichErwarteterFehler` gedeckt (Mutant M32 laut Implementer rot); nicht von mir nachgefahren.
- V-4: **Offen beim Nutzer, kein Befund gegen den Slice.** Abnahmeszenario 16, dritter Satz (aufgezeichnetes Verbindungsende, das der Server genauso verursacht, gilt als erwartet) ist nach der Bereinigung der Spezifikation unerreichbar: Die Aufzeichnung kennt kein Verbindungsende, `FATAL`/`PANIC` ist stets `PGR-E5004`, ein Verbindungsverlust ohne Fehlerantwort `PGR-E4003`. Kein Test belegt ihn, keiner kann es; das Lastenheft ist nicht geändert.
- V-5: INFO F-669 bis F-672 des Reviews sind bekannt und nicht erneut geprüft.

## 6. Was ich nicht selbst gefahren habe

- `make gates` (volle Kette), `make test`, `make lint`, `make a-check`, `a-check-negativ`, `lint-gegenprobe`, `hook-gegenprobe`; grün laut Implementer und Auftraggeber (Exit 0 an `c03f502`). Meine Läufe: `docs-check`, `kopf-check`, `abdeckung-check`, `test-integration`, `go test` der Pakete `services` und `bootstrap` (im Basislauf der Mutationen).
- Die übrigen Mutanten der Tabellen in §7 (60 von 78): Ich lese die Testnamen als passend, bestätigt sind nur die elf oben (V1 bis V7, Z1, E2). Das Ergebnis V-1 zeigt, dass „Mutation gelistet“ und „Zusage gefangen“ nicht dasselbe sind; die Auswahl der Mutanten pro Zusage ist daher der Punkt, an dem eine weitere Stichprobe lohnt.
- Integrationstests der Signale (`TestVergleichErstesSignal`, `…ZweitesSignal`) nur als Unit-Tests vorhanden, nicht E2E; ich habe sie nicht mutiert.
- Vergleiche mit einer zweiten Serverversion (akzeptiertes Negativ (b) in §6, gesetzt vom Architect).
- Closure-Punkte (Notiz, Register, Risiko-Ausgänge, Paarungen): offen, führt die Closure.

## 7. Urteil

DoD-Liefer-Punkte 1 bis 3 erfüllt, die Entscheidungen des Architect (`LH-FA-24.a` mit *Randformen des Vergleichs*, `LH-FA-21.a`, `LH-FA-20.a`) sind im Code, am gebauten Binary und in den Integrationstests wiederzufinden, die Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando` ist erfüllt, F-665 bis F-668 sind behoben, der Plan folgt dem Diff, Handbuch und README sind unverändert. Eine Verifier-only-Lücke: V-1 (Zusage „nur Wörter aus Ziffern“ gegen eine Mutation geschützt, `ContainsAny` bleibt grün); sie ist vor der Closure mit einem Testfall zu schließen. V-2 (Größenprobe als Schätzung) braucht eine Entscheidung des Planners oder eine Zeile in der Closure. Szenario 16, dritter Satz, liegt beim Nutzer.
