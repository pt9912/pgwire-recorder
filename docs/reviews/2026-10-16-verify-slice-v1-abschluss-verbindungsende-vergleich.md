# Verifikations-Report: slice-v1-abschluss-verbindungsende-vergleich — 2026-10-16

**Rolle:** Verifier (Modul 11), frischer Kontext. Frage: Bauen wir es richtig, gegen Plan, DoD und Entscheidungen?

**Gegenstand:** `slice-v1-abschluss-verbindungsende-vergleich` (Schritt 39a, Ende im Vergleich, Service), Diff `48325eae..8ee8f0f2`; Lieferung `85434400`, Nacharbeit `a1904b57` und `8ee8f0f2`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingang:** Plan (§1, §2, §3, §4, §6, §7 *Belege des Implementers* und *Nacharbeit zum Review*, §8), Review-Report vom selben Tag (0 HIGH, 1 MEDIUM, 2 LOW, 3 INFO; das Review wurde nicht neu gefahren), Spezifikation `LH-FA-24.a` (*Ende der Aufzeichnung*, *Fehler*, *Meldung*), `LH-FA-20.a` (Tabelle der Fehlerregeln), `LH-FA-25.a`, `SPEC-001`, `SPEC-002` (Randformen des Endes), Architect-Commits `51145c0`, `cc563f8d`, `ccbb1073`, `AGENTS.md` §3.9 bis §3.14.

## Ergebnis

**Der Slice erfüllt die Liefer- und Gate-Punkte der DoD; kein Befund blockiert die Closure.** Alle geprüften Zusagen tragen eine von mir selbst gesehene, rote Mutation; die Nacharbeit zu F-1, F-3, F-6 trägt (F-2 und F-5 geprüft); die Messung stimmt (566). Es gibt ein LOW (V-1, Schichtzählung nach F-6 fehlt in §8) und drei INFO. Offen sind die Closure-Pflichten, die erst nach der Verifikation stehen (Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen, DoD-Haken).

## 1. DoD-Punkte

| DoD-Punkt (§2) | Beleg | Urteil |
|---|---|---|
| Punkt 1: Ende im Vergleich (erwartet / `PGR-E5004`), Extended-Flush-Fälle, drei Kommentare, Zähler, `--allow-recorded-errors` ohne Wirkung, Formen, Session endet, leeres `V` wie fehlendes | Code `play.go` (`fatal`, `fehlerantwort`, `interaktion`, `abschluss`), `play_vergleich.go` (`erste`, `gewoehnlich`), `response.go` (`beendetVerbindung`); Tests `play_vergleich_ende_test.go` (fünf Tests, Läufer `laufVergleichEnde`); alle Zusagen gegen Mutanten in Abschnitt 3 | erfüllt |
| Punkt 2: je Zusage eine Mutation, Tabelle in §7 | Tabelle steht (M01 bis M32, M36 bis M38, E1); 36 meiner Mutanten bestätigen die Zeilen, die ich gefahren habe (Abschnitt 3) | erfüllt |
| Punkt 3: `make gates` grün | `.harness/state/gates-passed.diffsha` und `.dateiliste` von 18:12 (Lauf nach `a1904b57`, laut §7 Exit 0, 12 min 46 s); die Liste unterscheidet sich vom Arbeitsbaum in genau einer Datei, dem Plan (`docs/plan/planning/`, Positivliste); `git diff a1904b57..HEAD` berührt nur den Plan; `make gates-schnell` selbst gefahren: Exit 0, Stand `gates-schnell-passed.diffsha` geschrieben und gleich dem Hash des Arbeitsbaums; der Stop-Hook antwortet `approve` | erfüllt (den vollen Lauf habe ich nicht wiederholt; Grundlage: Nachweis, leerer Code-Diff seitdem, eigener Lauf der schnellen Kette mit `make test`, `make lint`, `make a-check`, `make abdeckung-check`, `make kopf-check`) |
| Punkt 4: Review durchgeführt, Report liegt vor | `2026-10-16-review-slice-v1-abschluss-verbindungsende-vergleich.md` (0 HIGH, 1 MEDIUM, 2 LOW, 3 INFO) | Report liegt vor; der Haken gehört der Closure |
| Punkte 5 bis 8: Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | §7 trägt nur die Belege und die Nacharbeit, keine Closure-Notiz, keine Ausgänge; DoD-Haken alle leer | offen, Closure-Arbeit (Vorschlag in Abschnitt 6) |

## 2. Plan gegen Code, Entscheidungs-Konformität

- **Plan gegen Diff (§3):** Der Diff berührt `play.go`, `play_vergleich.go`, `play_vergleich_ende_test.go` (neu), `play_vergleich_test.go`, `response.go`, `ende_test.go`, `yaml_ende_test.go`, `sqlite_lesen_ende_test.go`, `docs/user/abdeckung-unit.md`, die Roadmap und den Plan — genau die Pfade der Tabelle in §3. `git diff --stat 48325eae..8ee8f0f2` über `internal/hexagon/ports`, `internal/adapters/driving`, `internal/bootstrap`, `cmd`, `test`, `harness`, `tools`, `Makefile`, `.golangci.yml`, `.a-check.yml` ist leer; **Port, Einspiel-Adapter, CLI, Hilfetext, `TestE2EPlayVergleichFatal` und `SPEC-029` sind unberührt (bestätigt)**. Die beiden Leser (`internal/adapters/driven/recording`) sind nur in Tests berührt, ihr Produktcode nicht. Kein `nolint` im Diff, keine neuen Importe im Hexagon (die neue Testdatei importiert `model` und `services`).
- **Spezifikation (`LH-FA-24.a`, `LH-FA-20.a`, `LH-FA-25.a`) gegen Code:** Die drei Punkte von *Ende der Aufzeichnung* und der vierte (letzte `Flush`-Gruppe) stehen im Code: `fatal` hängt `connection_closed` an, zählt `verglichen`, ruft `erste(ende)`; `erste` liefert nil nur, wenn `gewoehnlich` nil liefert (gleiche Folge, damit gleicher SQLSTATE und aufgezeichnetes Ende an der Stelle), sonst die Art Fehler an der Position der Fehlerantwort, wenn die erste verschiedene Position dort oder dahinter liegt, sonst die frühere. Der aufgezeichnete SQLSTATE ist leer, wenn dort keine `error_response` steht (*Fehler*, ergänzt in `ccbb1073`). `--continue-on-error` und `--allow-recorded-errors` wie die Tabelle der Fehlerregeln. Verbindungsverlust geht vor jedem Vergleich (`PGR-E4003`, nicht gezählt).
- **Architect-Commits:** `ccbb1073` (leeres `V`, SQLSTATE nur an `error_response`) steht im Code und in der Spezifikation (Änderungshistorie, *Fehler*, Randformen des Endes); `51145c0` und `cc563f8d` (gewöhnlicher Vergleich bei normaler Antwort, Ausnahme der letzten `Flush`-Gruppe) sind durch die Fälle `Gewoehnlich` und M19 bis M21 getragen.
- **Randformen (`AGENTS.md` §3.12):** Der Diff entscheidet keine Randform still; die Stelle von `connection_closed`, die Reihenfolge der Gründe und die Schweregrade stehen in §6 und der Spezifikation.
- **Kommentare (`AGENTS.md` §3.7, §3.11):** Die geänderten Kommentare von `interaktion`, `fehlerantwort`, `fatal`, `abschluss`, `Zaehler`, `erste` und `beendetVerbindung` sagen zu, was Tests prüfen (Zuordnung in Abschnitt 3); der Schutz `ende <= len(v.soll)` in `erste` ist als Grenze ohne Test gekennzeichnet (Review F-5), und seine Begründung stimmt: Ist `u.position >= ende`, sind die Positionen davor gleich, und eine Aufzeichnung, die auf `ready_for_query` oder `connection_closed` endet, reicht dann mindestens bis `ende`.
- **Altes Wort:** Die Suche über `internal/` (und `test`, `harness`, `tools`, `docs/user`, `README.md`) nach „stets eine Abweichung“, „sofort die Abweichung“, „vorhanden, auch leer“, „entscheidet es allein“, `erste(true)`, `erste(false)` findet nichts Veraltetes; „bis zum ReadyForQuery“ steht nur dort, wo wirklich bis zum `ReadyForQuery` gelesen wird (`anfrage`, `extended` für `Sync`, `bisBereit`, Aufnahme, die Tests der `Sync`-Fälle); in Spezifikation und Handbuch lautet die Stelle `LH-FA-20.a` Zeile 1948 weiter richtig („erwartetes Ende oder sofort die Abweichung“).

## 3. Bewusstes Brechen (Modul 11), eigene Mutanten

Weg: Kopie von `go.mod`, `go.sum`, `internal`, `cmd` in ein frisches Verzeichnis je Mutant (Scratchpad, eigener Unterordner), genau eine Ersetzung mit erzwungenem Treffer, `docker run --rm --network=none` im Image der Stufe `deps`, `go test -count=1` über `./internal/hexagon/services/` (bei M37 zusätzlich Modell und beide Leser); der Arbeitsbaum blieb unberührt (`git status` sauber), ein Lauf zur Zeit, unveränderte Kopie grün (alle drei Pakete). Gezählt sind rote Zeilen `--- FAIL` (Ausgabe auf 40 Zeilen begrenzt, daher „≥ 40“).

| # | Zusage / Variante | Mutant | Ergebnis |
|---|---|---|---|
| F1 | nach dem Ende endet die Session (F-1) | `session`: `continue` statt `return nil` nach `errSessionEnde` | nur `TestVergleichEndeSessionEndet`, 17 Zeilen: `erwartetes Ende` e S F (12), `Abweichung` e S mit `--continue-on-error` (4); **beide Fälle rot**, sonst nichts |
| E1 | SQLSTATE zählt nur an `error_response` (F-3) | Typ-Check in `erste` entfällt | nur `TestVergleichEndeAbweichung/Server früher, Feld C an command_complete`, 13 Zeilen |
| M37 | leeres `V` wie fehlendes, Modell (F-6) | `beendetVerbindung` zurück auf „vorhanden, auch leer“ | `TestValidateEndeGueltig`, `TestUnmarshalEnde`, `TestLesenEndeAngenommen` rot (4 Zeilen) |
| M38 | dasselbe, Service | `fehlerantwort`: Rückfall auf `S` nur bei fehlendem Feld | 32 Zeilen: `TestPlayEndeExtendedErwartet`, `TestPlayEndeErwartet`, `TestVergleichEndeErwartet`, `…Abweichung` (jeweils Fall `V leer, S FATAL`) |
| M14, M14b | Server `PANIC` ist Ende; Server `FATAL` ist Ende | `fatal := schwere == "FATAL"` beziehungsweise `== "PANIC"` | ≥ 40, `TestVergleichEndeErwartet` und Bestandstests |
| M15 | Schweregrad aus `S`, wenn `V` fehlt | kein Rückfall auf `S` | ≥ 40, `…Erwartet/Server S FATAL ohne V` und `V leer, S FATAL` |
| M15b | dasselbe nur für `PANIC` | Rückfall auf `S` ignoriert `PANIC` | 34, `S PANIC ohne V` |
| M16 | `V` geht `S` vor | `S` zuerst | ≥ 40 (`V ERROR bei S FATAL` rot) |
| M13 | `FATAL` und `PANIC` beidseitig dasselbe Ende | Vergleich der Fehlerantworten nimmt `V` hinzu | ≥ 40, nur `TestVergleichEndeErwartet` (aufgezeichnet `PANIC`, Server `V PANIC`) |
| M02 | Glied „SQLSTATE gleich“ | `erste` gibt bei Art Fehler an der Fehlerantwort nil | 13, nur `Abweichung/anderer SQLSTATE` (e S F) |
| M03 | Glied „aufgezeichnetes Ende an der Stelle“ | `erste` gibt hinter der Fehlerantwort nil | ≥ 40, nur `TestVergleichEndeAbweichung` (die fünf Fälle `aufgezeichnet ERROR ohne Ende`) |
| M04 | frühere Abweichung geht vor | Umwandlung auch für frühere Positionen | 17, `Server später`, `frühere Abweichung` |
| M05 | Aufzeichnung ohne Fehlerantwort an der Stelle | Umwandlung nur mit `error_response` | 33 |
| M27, M28, M29 | Art Fehler; Position der Fehlerantwort; SQLSTATE des Servers in der Meldung | Art Nachrichtenart; `position: u.position`; `sqlstateIs: ""` | je ≥ 40, nur `TestVergleichEndeAbweichung` |
| M01 | `play` hängt `connection_closed` an | Anfügen entfällt | ≥ 40, nur `TestVergleichEndeErwartet` (Review F-4 bestätigt: nur dort beobachtbar) |
| M10, M11, M10b | „nichts danach“ / erwartet ohne Meldung | `fatal` gibt nil statt `errSessionEnde`; `PGR-E4003` beim erwarteten Ende; `fatal(l)` aufgerufen, aber weiter gelesen | je ≥ 40, `TestVergleichEndeErwartet` (M10b: Ablauf mit weiteren `Naechste`) |
| M24 | Schließen ohne `Terminate` mit Vergleich | `SchliesseOhneTerminate` nur ohne Vergleich | ≥ 40 |
| M09a, M09b | `--continue-on-error` | Fehler immer zurück; immer weiter gereicht | 21 und 14 Zeilen (Bestandstests und `Abweichung`) |
| M12 | `--allow-recorded-errors` ohne Wirkung | `fatal` meldet bei der Option das Ende als erwartet | ≥ 40, `TestVergleichEndeAbweichung` |
| MA | dasselbe, zweite Form | `l.ende` auch mit Vergleich gesetzt und in `fehlerantwort` vor `l.v` geprüft (zwei Ersetzungen) | ≥ 40, `TestVergleichEndeErwartet`; die einzelne Ersetzung der Zuweisung ist äquivalent (beide Felder werden erst nach `l.v != nil` gelesen) |
| M06, M07, M08, M30 | Zähler `verglichen`, `abweichend` einmal, erwartetes Ende nicht abweichend, `eingespielt` | `Verglichen++` entfällt; `Abweichend += 2`; `Abweichend++` vor der Prüfung auf nil; `Eingespielt++` entfällt | je ≥ 40 |
| M18 | gewöhnliche Regeln bei aufgezeichnetem Ende und normaler Antwort | `interaktion` vergleicht bei einfacher Anfrage mit Ende nicht | 13, `TestVergleichEndeGewoehnlich` |
| M19, M20, M21, M21b | letzte `Flush`-Gruppe: Rückkehr vor `bisBereit`; Lesen endet mit der Fehlerantwort; Vergleich; `verglichen` | Rückkehr nur ohne Vergleich; `n = 0` entfällt; `interaktion` vergleicht die letzte `Flush`-Gruppe nicht; `Verglichen++` in `abschluss` entfällt | M19: 13, nur Form Flush; M20: 24; M21: 13; M21b: ≥ 40 |
| M22 | Verbindungsverlust geht vor jedem Vergleich | `interaktion` vergleicht auch bei Fehler des Lesens | ≥ 40; `TestVergleichEndeVerbindungsverlust` rot in allen Formen (gezielt gefahren) |

Die Isolation der Konjunktionsglieder trägt: M02 färbt nur den Fall des anderen SQLSTATE, M03 nur die Fälle ohne aufgezeichnetes Ende, M04 und M05 nur ihre Fälle. Die Mutation M01 ist über die Schnittstelle nur im erwarteten Ende beobachtbar, wie §7 und Review F-4 sagen.

## 4. Nacharbeit zum Review (Fokus der Anfrage)

- **F-1 (`TestVergleichEndeSessionEndet`):** bestätigt. Der Mutant `continue` statt `return nil` färbt den Fall `erwartetes Ende` in allen drei Formen und den Fall `Abweichung` in beiden Formen mit `--continue-on-error` (ohne die Option bricht der Lauf ab, die Interaktion D ist dort nicht erreichbar). Die Aufzeichnung des Falls `erwartetes Ende` hat das Ende nicht in der letzten Interaktion der Session und wäre als Datei `PGR-E3003`; der Kommentar des Tests nennt das und die Grenze (der Service prüft die Stelle nicht, der Leser tut es). Das ist eine ehrliche Konstruktion für eine Service-Zusage, kein Verstoß.
- **F-6 (leeres `V` wie fehlendes):** In der Spezifikation entschieden (Absatz *Randformen des Endes in der Datei* in `SPEC-002`, Änderungshistorie vom 2026-10-16), vom Architect gesetzt (`ccbb1073`) und von Tests getragen: Modell (`TestValidateEndeGueltig`), YAML-Leser (`TestUnmarshalEnde`), SQLite-Leser (`TestLesenEndeAngenommen`), Service (M38; der Fall `V leer, S FATAL` läuft in jedem Test, der `fatalFaelle()` nutzt). Die drei alten Ablehnungszeilen sind entfernt, ihre Kommentare angepasst; kein Zeuge der Ablehnung bleibt. **Auswirkung auf das Laden:** nur Lockerung. `beendetVerbindung` hat genau einen Aufrufer (die Stelle von `connection_closed` im Modell), Dateien, die bisher luden, laden weiter; Dateien mit `V: ''` und `S: FATAL`/`PANIC` vor `connection_closed` waren `PGR-E3003` und laden jetzt. `LH-FA-22.a` (SQLite) verweist für Stelle und Form des Endes auf `SPEC-002` und `SPEC-041` und zieht mit; `SPEC-043` und das Lastenheft sagen nichts zum leeren `V` (gesucht). Der Aufzeichner schreibt kein leeres Feld (`set` überspringt `""` in `upstream.go`), also entsteht der Fall nur in handgeschriebenen oder umgeschriebenen Dateien. Schreiber-Roundtrip: Ich habe in einer Kopie eine Probe gefahren (Aufzeichnung mit `V: ""`, `S: FATAL`, Ende; YAML schreiben, laden, gleich; SQLite schreiben, laden, gleich): grün, der Schreiber schreibt `V: ""` und der Leser nimmt es an. Die Probe ist **kein Test im Repo**, keine Zusage verlangt ihn (V-2).
- **F-3 (Feld `C` an `command_complete`):** bestätigt, E1 rot und isoliert (13 Zeilen, nur dieser Fall).
- **F-2 (Formen je Fall):** §3 und DoD nennen die Formen; nachgezählt: Von den elf Fällen von `TestVergleichEndeAbweichung` fahren drei (`anderer SQLSTATE`, `Server früher`, `Server früher, Feld C`) in drei Formen, acht in zwei (fünf `aufgezeichnet ERROR ohne Ende` und `ohne Fehlerantwort der Aufzeichnung` haben aufgezeichnet kein Ende und damit keine `Flush`-Gruppe; `Server später` und `frühere Abweichung` erreichen die Fehlerantwort in der Form Flush nicht, weil das Lesen mit der Antwortzahl endet — von mir nachgerechnet am Fake); der Kommentar von `ohneFlush` nennt den Grund. Der Plan stimmt mit dem Diff.
- **F-5 (Kommentar am Schutz):** Kommentar vorhanden und inhaltlich richtig (siehe Abschnitt 2), als Grenze ohne Test gekennzeichnet.

## 5. Messung

`git diff --numstat 48325eae..8ee8f0f2 -- internal test docs/user` ergibt Einfügungen 407, Löschungen 159, **Summe 566**; bis `a1904b57` dieselbe Zahl (die Commits danach ändern nur Plan). Aufschlüsselung wie §7: `play_vergleich_ende_test.go` 315, `play_vergleich_test.go` 109, `play_vergleich.go` 55, `play.go` 39, `abdeckung-unit.md` 22, `response.go` 8, drei Testdateien je 6. Marke 800 zu 71 % (70,75 %), Meldeschwelle 640 nicht erreicht, Basis `48325eae` ist der Stand beim Anspruch (Plan-Commit `aa822df2`). Die Schätzung des Architect (305 bis 640, Mitte 470) ist eingehalten, der Wert liegt 96 über der Mitte; Faktor Tests zu Produktcode (442 zu 102 mit Kommentaren) 4,3.

## Findings

| ID | Schwere | Befund | Zuständig |
|---|---|---|---|
| V-1 | LOW | Die Aufnahme von F-6 in den Slice (Modell-Zeile, Tests der beiden Leser) steht in §3, DoD und §6, aber nicht in §8: Dort steht als letzte Zählung „Schichten: eine (Kern mit dem Service)“, keine Nachzählung nach `ccbb1073` und `a1904b57`. Die Änderung berührt das Modell (Logik, zählt nach `AGENTS.md` §3.13) und Tests des Lese-Adapters; mit `docs/user` (konservativ als Schicht geführt) wären es mehr als die zwei Schichten der Größenregel, und die Rückführung „Adapter muss sich ändern“ in §4 ist für den Produktcode des Adapters nicht eingetreten, für seine Tests nicht benannt. §1 *Ziel* nennt die Änderung des Ladens ebenfalls nicht (`AGENTS.md` §3.9). Kein Code-Mangel; ein Plan, der dem Diff nicht ganz folgt. | Implementer, im Closure-Commit vor dem `git mv`: Zählung in §8 mit Kennung des Gebers und Satz in §1; Planner entscheidet, ob die Zählung die Größenregel berührt |
| V-2 | INFO | Das Laden von Dateien mit leerem `V` vor einem Ende ändert sich (bisher `PGR-E3003`); es ist entschieden und getestet, aber kein Test im Repo fährt Schreiben und Laden einer Aufzeichnung mit leerem `V` (meine Probe in einer Kopie war grün). `LH-FA-20.a` (Zeile zum Schweregrad) und `SPEC-002` (Satz zur Stelle von `connection_closed`) sagen weiter „Feld `V`, ohne es `S`“; der Absatz *Randformen des Endes in der Datei* löst es, ohne dass die beiden Stellen auf ihn zeigen. | Architect (bei Gelegenheit: Zeiger oder Wortlaut „fehlt oder leer“); kein Eingriff vor der Closure |
| V-3 | INFO | Der Slice liegt bei 71 % der Marke 800 (Schätzung: Mitte 59 %), die Nacharbeit trug 81 von 566 Zeilen. Die Marke und die Meldeschwelle sind nicht erreicht; die Beobachtung gehört in die Closure-Notiz (Register: `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; ob sie zählt, entscheidet der Planner). | Planner |
| V-4 | INFO | Die Handbuch-Lücke (§1, akzeptiert, Entscheidung des Nutzers vom 2026-10-12) besteht mit diesem Slice: `docs/user/benutzerhandbuch.md` beschreibt `FATAL`/`PANIC` mit Vergleich noch nach der alten Regel; die Zusage liefert `slice-v1-abschluss-verbindungsende-doku`. Keine Aktion hier. | Planner (zur Kenntnis) |

Kein HIGH, kein MEDIUM.

## 6. Risiko-Ausgänge: Vorschlag zur Closure (nichts eingetragen)

Der Plan führt in §6 zwei Risiken, beide „offen bis Closure“; die Ausgänge setzt der Implementer in §7 im Commit vor dem `git mv`:

1. *Diff über einer Review-Sitzung* (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`): **entfallen** mit Begründung „vor dem Anspruch geschnitten; Marke 800 und Meldeschwelle 640 nicht erreicht (gemessen 566, 71 %), Schätzung 305 bis 640 eingehalten“. Dazu die Beobachtung V-3 in die Closure-Notiz und, wenn der Planner es zählt, ins Register unter `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (Zähler dort 2×, Stand offen).
2. *Szenario 16 Satz 3 nur gegen PostgreSQL 17* (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`): **entfallen** für diesen Slice (kein Server gefahren), mit dem Zeiger, dass das Risiko in §6 von `slice-v1-abschluss-verbindungsende-hilfetext-binary` steht (geprüft: dort §6, Zeile mit „4× bei Beginn“, und das akzeptierte Negativ in §1); die Adresse nimmt an.

Dazu in die Closure-Notiz: die Finding-Klassen des Reviews (Zusage ohne Mutation: Session endet nach Ende; Plan nennt mehr Formen als der Diff trägt; grüner Mutant als äquivalent eingeordnet; Schweregrad-Lesart verschieden) und V-1. Die Lerneintrag-Form ist Sache des Implementers; ich schlage vor: eine Zusage der Art „endet nach dem Ereignis“ braucht einen Fall mit einem Nachfolger (hier D hinter B), den eine Tabelle mit B als letzter Interaktion nicht liefert.

## Eigene Läufe

- `git diff --numstat` / `--stat` (Messung, Pfade, `nolint`, Importe)
- `make gates-schnell` → Exit 0, Stand geschrieben; Stop-Hook antwortet `approve`
- Vergleich der Dateiliste des vollen Laufs mit dem Arbeitsbaum (eine Abweichung: der Plan)
- 36 Unit-Mutanten in Kopien im Image der Stufe `deps` (`--network=none`), ein Lauf zur Zeit, Arbeitsbaum danach sauber; eine Probe zum Roundtrip von leerem `V` in einer Kopie
- `make docs-check` vor dem Commit dieses Reports (siehe Commit-Nachricht)

## Verdikt

**Verifikation bestanden** für die Liefer- und Gate-Punkte der DoD; keine Rückführung, kein Eingriff in Produktcode nötig. Offen: Closure-Pflichten (Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen, DoD-Haken) beim Implementer, V-1 beim Implementer und Planner, V-2 beim Architect.
