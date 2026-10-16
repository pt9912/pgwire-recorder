# Verifikations-Report: slice-v1-abschluss-verbindungsende-einspielen-grenzen — 2026-10-16

**Rolle:** Verifier (Modul 11), frischer Kontext. Frage: Bauen wir es richtig, gegen Plan, DoD und Entscheidungen?

**Gegenstand:** `slice-v1-abschluss-verbindungsende-einspielen-grenzen` (Schritt 38c), Diff `fa149c68..63dbbdb3`, Reports bis `45ac8530`; Lieferung `670ec850`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingang:** Plan (§1, §2, §3, §4, §6, §7 *Belege des Implementers*, §8), Review-Report vom selben Tag (0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO), Spezifikation `LH-FA-20.a` *Interaktion* und *Meldungen*, `LH-FA-25.a`, `SPEC-034` (`PGR-E4003`), Architect-Commits `526a50f`, `5c8c592`, `e92825a`, `51145c0`, `AGENTS.md` §3.9 bis §3.14.

## Ergebnis

**Der Slice erfüllt seine Liefer-Punkte und die Gate-Punkte der DoD; kein Befund blockiert die Closure.** Produktcode 0 Zeilen (bestätigt), alle fünf Zusagen tragen eine selbst gesehene, isolierte Mutation, die Messung stimmt (92). Offen sind die Closure-Pflichten, die erst nach der Verifikation stehen (Closure-Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen), und eine Lücke zwischen Lastenheft und Spezifikation, die der Slice nicht verursacht (V-1, Architect).

## 1. DoD-Punkte

| DoD-Punkt (§2) | Beleg | Urteil |
|---|---|---|
| Zwei Fälle am Binary gegen PostgreSQL, eine Aufzeichnung je Fall | `test/integration/play_ende_grenzen_e2e_test.go` (`TestE2EPlayEndeGrenzeOhneOption`, `TestE2EPlayEndeGrenzeOhneEnde`), `testdata/play-ende-grenze-ohne-option.yaml` und `play-ende-grenze-ohne-ende.yaml`; keine geteilte Fixture mit `play-ende-erwartet.yaml`, getrennte Datenbanken; beide Fälle `PASS` in jedem meiner Läufe ohne Mutant | erfüllt |
| Fall 1: Ende aufgezeichnet, ohne Option: `PGR-E4003`, Exit 4 | Fixture trägt `FATAL 57P01` und `connection_closed` als letzte Antwort; Aufruf ohne `--allow-recorded-errors` | erfüllt, am Grund rot gesehen (Mutant 1) |
| Fall 2: `ERROR` aufgezeichnet, Server sendet `FATAL`, mit Option: `PGR-E4003`, Exit 4 | Fixture trägt `ERROR 57P01` und `ready_for_query`, kein Ende | erfüllt, am Grund rot gesehen (Mutanten 2 und 3) |
| Je Zusage eine Mutation, Tabelle in §7 | Tabelle steht in §7; alle fünf Zeilen von mir nachgefahren (siehe 3) | erfüllt |
| `make gates` grün | Nachweis `gates-passed.diffsha` vom 16:06 für den Stand `670ec850`; `git diff fa149c68..HEAD -- internal cmd` leer; der Hash des Arbeitsbaums weicht seitdem nur um Dateien der Positivliste ab (`docs/reviews/`, `docs/plan/planning/`); `make gates-schnell` von mir gefahren: Exit 0, Stand `gates-schnell-passed.diffsha` geschrieben; der Stop-Hook gibt frei | erfüllt (den vollen Lauf habe ich nicht wiederholt; Grundlage ist der Nachweis und der leere Code-Diff) |
| Review durchgeführt, Report liegt vor | `2026-10-16-review-slice-v1-abschluss-verbindungsende-einspielen-grenzen.md`; DoD-Haken im Plan noch leer, richtig für diesen Stand | Report liegt vor, Haken setzt die Closure |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | noch nicht geschrieben (§7 trägt nur die Belege) | offen, Closure-Arbeit (siehe 6) |

## 2. Produktcode 0 Zeilen

`git diff --stat fa149c68..HEAD -- internal cmd` ist leer (selbst gefahren). Der Diff `fa149c68..63dbbdb3` berührt sechs Pfade: Roadmap (Zeile „In Arbeit“), den Plan (Move von `next/` nach `in-progress/` mit Inhalt), `docs/user/abdeckung-e2e.md`, den Test und zwei Fixtures. Kein `nolint`, kein Import im Hexagon geändert.

## 3. Bewusstes Brechen (Modul 11), eigene Mutanten

Mutanten im Arbeitsbaum, **nicht** committed (der Plan lässt den Commit vor dem Lauf zu, weil der Nachweis der Gates am Hash hängt; hier ist der Nachweis nicht Gegenstand, das Image baut aus dem Arbeitsbaum), `touch` vor dem Lauf, je ein voller `make test-integration`, ein Lauf zur Zeit, danach `git checkout internal`; `git status` am Ende sauber.

| # | Zusage / Variante | Mutant (Stelle in `internal/`) | Ergebnis |
|---|---|---|---|
| 1 | Ende immer erwartet (Option nicht geprüft) | `services/play.go` Z. 261 `l.ende = traegtEnde(in)` | nur `TestE2EPlayEndeGrenzeOhneOption` rot (Exit-Code 0, erwartet 4); `…OhneEnde` grün |
| 2 | Option allein genügt (Aufzeichnung nicht geprüft) | Z. 261 `l.ende = s.optionen.AllowRecordedErrors` | nur `TestE2EPlayEndeGrenzeOhneEnde` rot (Exit-Code 0); `…OhneOption` grün |
| 3 | eigene Variante: `FATAL` gilt als erwartet, sobald die Aufzeichnung eine Fehlerantwort trägt (`erwartet` statt `ende`) | Z. 414 `if l.ende \|\| l.erwartet {` | nur `TestE2EPlayEndeGrenzeOhneEnde` rot (Exit-Code 0); `…OhneOption` grün. Ein anderer Mutant als Nr. 2, am selben Fall gefangen |
| 4 | Stelle der Interaktion (Session und Interaktion) | Z. 203 Meldung nur `Session %d` | beide Grenzfälle rot, dazu sechs Bestandstests (`TestE2EPlayFehlerantwort`, `…Fortsetzung`, `…FortsetzungAbbruch`, `…ErwarteterFehler`, `…VergleichAbweichung`, `…VergleichFatal`) |
| 5 | Code `PGR-E4003` | Z. 417 `CodeServerError` statt `CodeConnectionLost` | beide Grenzfälle rot, dazu `TestE2EPlayFortsetzungAbbruch` |
| 6 | Exit-Code 4 | `model/fehler.go` `exitCode`: `PGR-E4003` liefert 5 | beide Grenzfälle rot (Exit-Code 5, erwartet 4), dazu `TestE2ERecordExtendedAbbruch`, `TestE2ERecordExtendedSigtermNachBlockade`, `TestE2EPlayFortsetzungAbbruch` |
| 7 | eigene Variante: SQLSTATE und Meldung der Fehlerantwort aus der Meldung gestrichen | Z. 417 Text ohne `%s „%s“` | **alle Tests grün (Exit 0)**, siehe V-2 |

Die Tabelle in §7 stimmt mit meinen Läufen überein (die Zeilennummern 203, 261, 417 stimmen mit `play.go`). Die Isolation von 1 und 2 bleibt, auch gegen Variante 3: Jeder Mutant färbt nur den Fall seiner Zusage.

## 4. Entscheidungs-Konformität

- `LH-FA-20.a` *Interaktion*: ohne `--allow-recorded-errors` ist `FATAL` sofort `PGR-E4003`, auch wenn die Aufzeichnung ein Ende trägt (Fall 1); mit der Option ist das Ende erwartet nur, wenn die Aufzeichnung der Interaktion `connection_closed` trägt (Fall 2). Code (`play.go` Z. 261, 414 bis 417) und Tests entsprechen der Spezifikation und den Architect-Commits (ein Grund, `PGR-E4003`, keine Reihenfolge; *Stelle* ist die Interaktion).
- `LH-FA-25.a` *Form*: Das Ende steht als letzte Antwort der einfachen Anfrage; Fixture 1 hält das ein.
- `SPEC-034` / `PGR-E4003` Exit 4: von den Tests am Binary gefangen (Mutanten 5 und 6).
- Randformen (§3.12): Der Slice fügt keinen Vertrag hinzu; keine Randform wird im Diff still entschieden. Die Stelle des `connection_closed` (a) und der eine Grund (d) stehen in §6 und der Spezifikation.
- Adressen (§3.13): Der Diff weist keinem Nehmer etwas zu. Die Übernahme aus `slice-v1-abschluss-verbindungsende-einspielen-extended` steht in §1 des Nehmers dieses Slice.
- Hard Rules (§3.2, §3.9, §3.10, §3.11): ohne Befund (der Plan folgt dem Diff; §3 nennt die Dateinamen, die der Diff trägt; jede Zusage hat eine Mutation; die Abdeckungs-Texte sagen zu, was die Tests prüfen, siehe V-3).

## 5. Plan-vs-Code-Diff, Abdeckung, Messung

- **Plan gegen Diff:** §3 nennt `play_ende_grenzen_e2e_test.go`, zwei Fixtures und `docs/user/abdeckung-e2e.md`; der Diff liefert genau diese. Der Test hat die zugesagten Namen. §1 bis §8 folgen dem Diff.
- **Abdeckung:** `make abdeckung-check` Exit 0 (selbst gefahren); die Tabelle trägt je Test zwei Zeilen (`LH-FA-20` und `LH-FA-25`), insgesamt vier, wie §7 meldet.
- **Fixtures je Fall:** eine Datei je Fall, beide per `go:embed`, keine andere Datei bindet sie.
- **Messung nachgerechnet:** `git diff --numstat fa149c68..63dbbdb3 -- internal test docs/user` ergibt 52 + 18 + 18 + 4 = **92**; bis `HEAD` ebenfalls 92. Marke 400, Meldeschwelle 300, beide nicht erreicht; die Schätzung 140 bis 190 lag darüber (46 je Fall gegen 65 bis 70 im Geber, weil die Hilfe `spieleEndeNichtErwartet` geteilt ist). Die Basis `fa149c68` ist der Stand beim Anspruch.

## Findings

| ID | Schwere | Befund | Zuständig |
|---|---|---|---|
| V-1 | INFO (Spec-Lücke, kein Befund an diesem Slice) | `LH-FA-25.a` *Negative* verlangt für ein Ende, das die Aufzeichnung „mit einem anderen Fehlercode zeigt“, einen Verbindungsfehler (beim Einspielen ohne Vergleich); `LH-FA-20.a` *Interaktion* (Spezifikation) lässt das Ende mit `--allow-recorded-errors` erwartet sein, „gleich welcher SQLSTATE und an welcher Stelle“. Das Lastenheft rangiert höher. Beide Fixtures dieses Slice tragen `57P01`; der Slice entscheidet es nicht und beansprucht es nicht (bestätigt Review F-4). | Architect (Spezifikation, bei Gelegenheit, gegebenenfalls Change Request an das Lastenheft); Planner nimmt den Punkt als Beobachtung oder Folge-Slice-Adresse |
| V-2 | INFO | SQLSTATE und Meldung der Fehlerantwort in der Meldung von `PGR-E4003` sind am Binary nicht belegt: Mutant 7 bleibt in `make test-integration` grün (vertieft Review F-1). Eine Zusage dieses Slice ist es nicht (§6 *Zähler und Nachricht*); `LH-FA-20.a` *Meldungen* sagt sie allgemein zu, und die Tests des Service (`play_ende_test.go` Konstante `e4003B`, `play_extended_test.go`) tragen den Wortlaut mit SQLSTATE und Meldung. Ich habe den Mutanten dort nicht gefahren, der Test ist gelesen. Keine Änderung an diesem Slice nötig; wer die Zusage am Binary haben will, nimmt sie in einen eigenen Fall. | Implementer zur Kenntnis (Closure-Notiz: akzeptiertes Negativ, Beleg liegt im Service); keine Aktion |
| V-3 | INFO | Die Abdeckungs-Zeile von `TestE2EPlayEndeGrenzeOhneEnde` beginnt mit „wie TestE2EPlayEndeGrenzeOhneOption, mit …“ und steht in `docs/user/abdeckung-e2e.md` alphabetisch vor ihrer Bezugszeile (Review F-3). `make abdeckung-check` ist grün, die Zeile sagt zu, was der Test prüft; sie ist ohne die Nachbarzeile nicht eigenständig lesbar. | Implementer zur Kenntnis; kein Eingriff vor der Closure |
| V-4 | INFO | Der Plan kennt in §7 die Schätzung 140 bis 190 gegen gemessen 92 als Beobachtung („Gegenrichtung“ zu `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`). Das gehört in die Closure-Notiz des Slice; ob es zählt, entscheidet der Planner. | Planner |

Kein HIGH, kein MEDIUM, kein LOW.

## 6. Risiko-Ausgänge: Vorschlag zur Closure (nichts eingetragen)

Der Plan führt vier Risiken in §6, alle mit „offen bis Closure“. Die Ausgänge setzt der Implementer in §7 im Commit vor dem `git mv`; mein Vorschlag, den ich für tragend halte:

1. *Grüner Test ohne Mutation:* **entfallen** — Mutanten 1 bis 6 rot gesehen (§7 und diese Verifikation), mit Begründung „Isolation der beiden Fälle bestätigt“.
2. *Handgeschriebene Aufzeichnungen ohne Beleg gegen den Server* (`BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server`, bisher 3×): **weiter offen** ins Register; dieser Slice berührt die Beobachtung erneut (4×) — der Zähler steht über der Schwelle 3, die Adresse ist `slice-v1-abschluss-verbindungsende-upstream-extended`; der Planner entscheidet, ob die Adresse trägt.
3. *Verhalten des Servers nur gegen Version 17* (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×): **weiter offen** ins Register (2× mit diesem Slice); `make test-postgres-versionen` fährt die Fälle dieses Slice nicht (er hängt nur an `make test-integration`).
4. *Befund verlangt Produktcode:* **entfallen** — `git diff fa149c68..HEAD -- internal cmd` ist leer, keine Rückführung.

Dazu in die Closure-Notiz: die Finding-Klassen des Reviews (vier INFO) und V-4.

## Eigene Läufe

- `git diff --numstat` / `--stat` (Messung, Produktcode)
- `make abdeckung-check` → Exit 0
- `make gates-schnell` → Exit 0, Stand geschrieben; Stop-Hook gibt frei
- sieben `make test-integration` mit Mutant (Nr. 1 bis 7), ein Lauf zur Zeit; Arbeitsbaum danach sauber
- `make docs-check` vor dem Commit dieses Reports (siehe Commit-Nachricht)

## Verdikt

**Verifikation bestanden** für die Liefer- und Gate-Punkte der DoD; keine Rückführung, keine Änderung an Produktcode nötig. Offen: Closure-Pflichten (Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen) beim Implementer, V-1 beim Architect.
