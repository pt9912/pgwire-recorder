# Verifikations-Report: slice-v1-abschluss-verbindungsende-hilfetext-binary — 2026-10-16

**Rolle:** Verifier (Modul 11), frischer Kontext. Frage: Bauen wir es richtig, gegen Plan, DoD und Entscheidungen?

**Gegenstand:** `slice-v1-abschluss-verbindungsende-hilfetext-binary` (Schritt 39b, Hilfetext und Abnahmeszenario 16 Satz 3 am Binary), Diff `1de25738..a8309486`; Lieferung `cc4ca898` und `f4e9defb`, Belege `db1b0ac5` und `a8309486`, Architect vor dem Code `179ab424`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingang:** Plan (§1, §2, §3, §4, §6, §7 *Belege des Implementers*, §8), Review-Report vom selben Tag (0 HIGH, 0 MEDIUM, 1 LOW, 4 INFO; nicht neu gefahren), Spezifikation `LH-FA-20.a` (Absatz unter der Optionstabelle), `LH-FA-24.a`, Abnahmeszenario 16 Satz 3, `SPEC-029` Punkt 2, `SPEC-053` Punkt 14, `AGENTS.md` §3.9 bis §3.14.

## Ergebnis

**Der Slice erfüllt die Liefer- und Gate-Punkte der DoD; kein Befund blockiert die Closure.** Die Stichproben (Hilfetext, „gleich“, „anders“, `--allow-recorded-errors` mit Vergleich) wurden rot gesehen, jede färbt nur den Fall ihrer Zusage; die Messung 222 stimmt. Der Code-Stand seit `a8309486` ist unverändert, der Nachweis des vollen Laufs gilt. Ein LOW (V-1, entspricht F-1 des Reviews: Entscheidung unten) und zwei INFO (V-2, V-3). Offen sind die Closure-Pflichten, die erst nach der Verifikation stehen.

## 1. DoD-Punkte

| DoD-Punkt (§2) | Beleg | Urteil |
|---|---|---|
| Punkt 1: Hilfetext nennt Ende ohne Vergleich, Wirkungslosigkeit mit Vergleich, erwartetes Ende; je Satz ein Test im CLI-Adapter, Mutation färbt ihn | `cli.go` (`optionenPlay`) Wortlaut wie §6 (Diff gelesen, nur die Konstante geändert, die erste Aussage unverändert); `play_ende_hilfe_test.go`: acht Teilfälle, Block je Option; eigener Mutant `an derselben Stelle` → `an einer Stelle` färbt genau `--compare-responses: an derselben Stelle` | erfüllt (Einschränkung V-1) |
| Punkt 2: Szenario 16 Satz 3 am Binary gegen PostgreSQL, gleich: Exit 0, nächste Session, Zeile `eingespielt`/`verglichen`/`abweichend` `3/3/0`; anders: `PGR-E5004`, Exit 5, `2/2/1`; Satz in `SPEC-029` Punkt 2 im Commit des Tests; `TestE2EPlayVergleichFatal` nur Beschreibung | `play_ende_vergleich_e2e_test.go`, zwei Fixtures; `SPEC-029` Punkt 2 steht in `f4e9defb` (Commit der Tests), nicht früher; die Berichtigung des Numerus steht in `179ab424` und wird nicht doppelt geschrieben; `play_vergleich_e2e_test.go`: nur Kommentar, Code und Erwartung per Diff unberührt. Der Satz stimmt mit `SPEC-053` Punkt 14 (5) und der Teilmenge der Matrix überein | erfüllt |
| Punkt 3: je Zusage eine Mutation, Tabelle in §7 | zwölf Hilfetext-Zeilen, fünf E2E-Zeilen (M1 bis M5), drei als nicht eigens belegt benannt; meine Stichproben stimmen mit den Zeilen überein (Abschnitt 3); Beleg (a) ist falsch begründet (V-2) | erfüllt (V-2 INFO) |
| Punkt 4: `make gates` grün | `.harness/state/gates-passed.diffsha` und `.dateiliste` von 19:46 (Lauf nach `db1b0ac5`/`a8309486`, §7: Exit 0); `git diff a8309486..HEAD` berührt nur `docs/plan/planning/` und `docs/reviews/` (Positivliste), kein Code, keine Spezifikation, kein Handbuch. Selbst gefahren: `make abdeckung-check` Exit 0, `make docs-check` 0 Befunde, `make gates-schnell` Exit 0, Stand der schnellen Kette geschrieben; der Stop-Hook antwortet `approve`. Den vollen Lauf habe ich nicht wiederholt (Bedingung des Auftrags: Code-Stand gleich) | erfüllt |
| Punkt 5: Review, Report liegt vor | `2026-10-16-review-slice-v1-abschluss-verbindungsende-hilfetext-binary.md` (0 HIGH, 0 MEDIUM, 1 LOW, 4 INFO) | Report liegt vor; der Haken gehört der Closure |
| Punkte 6 bis 9: Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | §7 trägt Belege, keine Closure-Notiz, keine Ausgänge; alle Haken leer | offen, Closure-Arbeit (Vorschlag in Abschnitt 5) |

## 2. Plan gegen Code, Entscheidungs-Konformität

- **Plan gegen Diff (§3):** Der Diff berührt `cli.go`, `play_ende_hilfe_test.go`, `play_ende_vergleich_e2e_test.go`, zwei Fixtures, `play_vergleich_e2e_test.go`, `spec/spezifikation.md`, die zwei Abdeckungstabellen, Roadmap, Plan und `slice-v1-abschluss-verbindungsende-upstream` (Plan-Nachzug) — die Pfade der Tabelle in §3 und nicht mehr. Kein Produktcode des Service, des Ports, des Einspiel-Adapters oder des Bootstrap; die zweite Rückführung von §4 ist nicht eingetreten. Kein `nolint`.
- **`LH-FA-20.a`, Absatz unter der Optionstabelle (`179ab424`):** viererlei zu `--allow-recorded-errors` und eine Aussage zu `--compare-responses`; `optionenPlay` trägt jede, der Test je Aussage (acht Needles) deckt sie. „Ein Test je Aussage prüft sie“ ist im Endstand wahr (Review F-5: im Zwischenstand vor `cc4ca898` nicht; Architect-Commit vor dem Code, Absicht, keine Aktion).
- **`LH-FA-24.a` *Ende der Aufzeichnung* und Abnahmeszenario 16 Satz 3:** gleicher SQLSTATE an derselben Stelle ist erwartet (Exit 0, nächste Session läuft, `3/3/0`), anderer SQLSTATE ist `PGR-E5004` der Art Fehler mit beiden SQLSTATE in der Meldung, Exit 5, Abbruch, `2/2/1`, auch mit `--allow-recorded-errors`. Der Fall „anders“ prüft die volle Meldung, genau eine `PGR-E5004`, kein `PGR-E4003` und die Datenbank (`to_regclass`). Das Lastenheft nennt im Szenario nur „gleiche Stelle, gleicher Fehlercode, nächste Session“ und impliziert den Abweichungsfall durch Satz 2; der Test deckt beides.
- **`SPEC-029` Punkt 2:** Der Satz „Das erwartete Verbindungsende beim Einspielen (Abnahmeszenario 16, dritter Satz) ist nur dort geprüft“ nennt nur Szenario 16 Satz 3 und die Referenz 17; die Matrix fährt es nicht (`SPEC-053` Punkt 14 (5)). Er steht erst im Commit des Tests (`AGENTS.md` §3.11). Der Numerus „die mit MD5“ steht in `179ab424`.
- **Handbuch und README:** nicht berührt (Abgrenzung §1, Folge-Slice `-doku`); `AGENTS.md` §3.11 verlangt hier keine Änderung.
- **Fixtures (`AGENTS.md` §3.9, §3.12):** Die Abweichung vom Wortlaut des Architect (`CREATE TABLE`, DO-Block statt `SELECT 1`/`SELECT pg_terminate_backend`) ist in §6 nachgezogen, im selben Commit wie die Tests, und entscheidet keine Randform (Review F-3, geteilt).

## 3. Bewusstes Brechen (Modul 11), eigene Stichproben

Weg: Hilfetext in einer Kopie (`git archive HEAD`) im Scratchpad, Image `deps` mit `--network=none`; E2E in einer Kopie des ganzen Baums mit `make test-integration` (Parameter `PGR_ROLLEN=scram`, `PGR_TESTS_1` mit den zwei neuen Tests, nur Phase 1 des Runners), ein Lauf zur Zeit; der Arbeitsbaum blieb unberührt (`git status` sauber).

| Zusage | Mutant | Ergebnis |
|---|---|---|
| Hilfetext `--compare-responses`: „an derselben Stelle“ | `cli.go`: `an derselben Stelle mit` → `an einer Stelle mit` | rot nur `TestPlayEndeHilfe/--compare-responses:_an_derselben_Stelle` |
| „gleich“: gleiche Stelle, gleicher SQLSTATE ist erwartet | `play.go` `fatal`: `if u == nil && false` (Ende nie erwartet) | rot nur `TestE2EPlayEndeVergleichGleich`; „anders“ (beide Läufe) grün |
| „anders“: anderer SQLSTATE ist `PGR-E5004`, Exit 5 | `fatal`: `if u == nil \|\| true` (jedes Ende erwartet) | rot beide Läufe von `TestE2EPlayEndeVergleichAnders`; „gleich“ grün |
| `--allow-recorded-errors` ohne Wirkung mit Vergleich | `fatal`: `if u == nil \|\| s.optionen.AllowRecordedErrors` | rot nur `TestE2EPlayEndeVergleichAnders/play_ende_vergleich_anders_allow`; der Lauf ohne Option und „gleich“ grün |

Dazu der Mutant zu F-1: In `--compare-responses` zusätzlich „mit --compare-responses ohne Wirkung;“ vor „ein Ende des Servers“ — `TestPlayEndeHilfe` bleibt **grün** (bestätigt, siehe V-1).

## 4. Beurteilung der Befunde des Reviews

**F-1 (Hilfetest prüft Anwesenheit im Block, nicht Stelle und nicht Abwesenheit im anderen Block).** Mutant gefahren, grün: bestätigt. Ist das eine Zusage ohne Test im Sinne von `AGENTS.md` §3.10? **Nicht im strengen Sinn.** Die Zusagen, die Spezifikation, Plan und Kommentar machen, sind die Anwesenheit je Aussage im Block der Option („die Aussagen stehen im Block ihrer Option“, „viererlei“); jede davon fängt die Mutation „Aussage oder Glied entfernt“ (zwölf Zeilen, eine selbst gefahren). Weder die Spezifikation noch der Testkommentar sagen die Abwesenheit im Block der anderen Option oder die Reihenfolge innerhalb des Blocks zu; die Reihenfolge der Sätze ist wie der Umbruch keine Zusage (Plan §6). Aber: Der Plan §7 trägt nur „nicht gefahren“, keinen Grund, und das Verschieben in den falschen Block ist eine Fehlform, die ein Leser des Hilfetextes sieht („mit --compare-responses ohne Wirkung“ im Block von `--compare-responses` wäre sinnlos). Die Abwesenheit ist billig prüfbar (Needle der Gegenoption darf im Block nicht stehen). **Urteil: LOW (V-1), tragbar als akzeptiertes Negativ mit Grund oder durch eine Abwesenheitsprüfung zu schließen; Zuständigkeit Implementer (Nacharbeit zum Review).** Die Stelle innerhalb des Blocks bleibt unzugesagt und braucht keinen Test.

**F-2 (Beleg „kein `PGR-E4003` nicht trennbar“ ist falsch).** Aus der Struktur des Tests nachvollzogen: Die Prüfung „kein `PGR-E4003`“ ist eine eigene Klausel der Bedingung (`strings.Contains(stderr, "PGR-E4003")`) neben Exit-Code, genau einer `code=PGR-E5004`-Meldung und Zähler; eine zusätzliche `PGR-E4003`-Meldung bei sonst gleichem Ergebnis färbt allein sie. Den Mutanten X2 des Reviews habe ich nicht neu gefahren; die Klausel ist unabhängig lesbar, die Begründung in §7 (a) ist unrichtig, die Prüfung trägt. Hinweis ohne Folge für die Zusage: V-2 (INFO, Implementer, Beleg in §7).

**F-4 (Historie-Zeile).** Bestätigt: Die Zeile vom 2026-10-16 in `spec/spezifikation.md` sagt „Berichtigung des Numerus, keine Aussage geändert“ für `SPEC-029` Punkt 2; `f4e9defb` fügt dort den Satz zum erwarteten Verbindungsende hinzu, ohne Historie-Zeile. V-3 (INFO, Implementer).

## 5. Messung, Closure-Vorbereitung

- **Messung nachgerechnet:** `git diff --numstat 1de25738..a8309486 -- internal test docs/user` ergibt Summe **222** (`cli.go` 12, Hilfetest 60, E2E-Datei 78, Fixtures 28 + 28, `play_vergleich_e2e_test.go` 9, Tabellen 5 + 2). Marke 400 nicht erreicht (56 %), Meldeschwelle 300 nicht erreicht, Schätzung 145 bis 355 getroffen. Die Pfade ändern sich nach `f4e9defb` nicht mehr.
- **Vorschlag der Risiko-Ausgänge zur Closure (§6, Urteil Planner):**
  1. *Diff über einer Review-Sitzung* — **entfallen**: Messung 222 unter der Meldeschwelle, Review in einer Sitzung.
  2. *Serververhalten nur gegen PostgreSQL 17* — **weiter offen** → `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (fünftes Auftreten; die Grenze steht in `SPEC-029` Punkt 2 und `SPEC-053` Punkt 14 (5)).
  3. *Handgeschriebene Aufzeichnungen* — **weiter offen** → `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` (Auftreten: die Fixture brauchte die Form der Antworten, der erste Lauf war rot, Beobachtung 1 in §7).
  4. *Fall am Binary rot wegen Service* — **entfallen**: im Diff kein Produktcode des Service, alle Läufe grün.
  5. *Satz in `SPEC-029` sagt mehr zu als der Test prüft* — **entfallen**: Der Satz nennt nur Szenario 16 Satz 3 und die Referenz 17; die Matrix fährt den Fall nicht (`SPEC-053` Punkt 14 (5)), der Test läuft in der vollen Suite.
- **Finding-Klassen für die Closure-Notiz:** Hilfetext-Test prüft Anwesenheit, nicht Zuordnung · Beleg-Aussage enger als die Prüfung · Historie nennt neuen Satz der Spezifikation nicht (Review) und, aus der Verifikation, nichts darüber hinaus.

## Findings

| ID | Kategorie | Befund | Zuständig | Pfad |
|---|---|---|---|---|
| V-1 | LOW | `TestPlayEndeHilfe` prüft die Anwesenheit je Aussage im Block der Option; die Abwesenheit im Block der anderen Option ist nicht geprüft (Mutant „mit --compare-responses ohne Wirkung“ zusätzlich im Block von `--compare-responses` bleibt grün, gefahren), und §7 nennt das Verschieben nur „nicht gefahren“, ohne Grund. Keine Zusage ohne Test im Sinne von `AGENTS.md` §3.10 (Spezifikation und Kommentar sagen die Abwesenheit nicht zu), aber eine ungenannte Lücke. Entweder eine Abwesenheitsprüfung im Test (die Needles der Gegenoption fehlen im Block) oder in §7 ein akzeptiertes Negativ mit Grund (Reihenfolge und Stelle im Block sind keine Zusage). | Implementer (Nacharbeit zum Review, vor der Closure) | `internal/adapters/driving/cli/play_ende_hilfe_test.go`; Plan §7 |
| V-2 | INFO | §7 „Nicht mit eigenem Mutanten belegt (a)“: „kein `PGR-E4003`“ sei nicht trennbar; die Prüfung ist eine eigene Klausel und trennbar. Beleg-Satz berichtigen. | Implementer | Plan §7 |
| V-3 | INFO | Historie-Zeile der Spezifikation vom 2026-10-16 nennt für `SPEC-029` Punkt 2 „keine Aussage geändert“; der Satz zu Szenario 16 Satz 3 aus `f4e9defb` hat keine Historie-Zeile. Zeile ergänzen oder die Zeile berichtigen. | Implementer | `spec/spezifikation.md` Änderungshistorie |

**Merge-blockierend:** nein. V-1 bis V-3 gehen mit der Nacharbeit zum Review vor der Closure; eine Entscheidung gegen V-1 gehört in §7.

## Übergabe

An den Planner: DoD-Punkte 1 bis 4 erfüllt, Punkt 5 mit diesem Report und dem Review belegt, Punkte 6 bis 9 sind Closure-Arbeit (Notiz mit Lerneintrag, Register, Ausgänge nach Abschnitt 5, Paarungen). Der Nehmer `slice-v1-abschluss-verbindungsende-upstream` (Schritt 41) prüft die Berichtigung des Numerus, sie steht in `179ab424`.
