# Verifikation: slice-harness-blackbox-kern — 2026-10-06

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft; F-441 prüfe ich hier auf Auftrag als erneutes Review mit.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-blackbox-kern.md` (Kopf, §1 bis §8) gegen den Gesamt-Diff `a12f16f..ebbf5b8`. Darin: `839e61c` (Umbau, Mutationstabelle in §7), `2bea0f0` ([Review](2026-10-06-review-slice-harness-blackbox-kern.md), F-441 bis F-443), `b0b2286` (Architect, [`SPEC-049`](../../spec/spezifikation.md#spec-049--lint-profil-lint) Punkt 7 und *Grenze*) und `ebbf5b8` (Nacharbeit F-441).

**Eingang:**

- die DoD-Liefer-Punkte, §3, §6 *Randformen* und *Risiken*, die Mutationstabelle in §7 und die Commit-Messages
- `spec/spezifikation.md` `SPEC-049` Punkt 7, 8 und *Grenze* am Stand `ebbf5b8`; [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 2, 4 und 5; [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md); [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- Produkt-Code als Bezug: `internal/hexagon/services/replay.go` (`NewReplayService`, `OpenConnection`, `zuordnen`, `letzteNummer`, `Query`, `ClientMessage`), `lebendpruefung.go`, `record.go`, `internal/hexagon/model/fehler.go`
- der Review-Report (F-441 bis F-443, Stand `839e61c`)
- die DoD-Bestätigung des Implementers nur in der Zusammenfassung des Auftrags (62 Tests, Deklarationen gleich, `make lint` 0 in den Testdateien und 72 modulweit, kein `_ =`, kein `//nolint`). Ein Bericht mit den Lint-Zeilen vor und nach lag mir nicht vor (V-82); alle Belege unten habe ich in diesem Kontext selbst erzeugt.

Bei meinem Start war der Arbeitsbaum sauber, HEAD `ebbf5b8`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

**Wie die Proben gebaut wurden.** Je Stand (`a12f16f` vor dem Umbau, `839e61c`, `ebbf5b8`) eine Kopie aus `git archive` im Scratchpad. Testlisten und Mutationen liefen im Image der Stufe `deps` (aus `ebbf5b8` gebaut) mit Bind-Mount der Kopie und `--network=none`. Jede Mutation lief in einer frischen Kopie (`cp -r` ohne `-p`) von `ebbf5b8`; das Ersetzungsskript bricht ab, wenn der Suchtext nicht genau einmal trifft. Lint-Läufe auf Kopien mit `docker build --progress=plain --target lint`. Das Repo selbst blieb unverändert.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — `_test`-Pakete, Testliste, Abdeckung, Brücke. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Testdateien gehören zu `_test`-Paketen | Kopfzeilen: 2× `package model_test`, 7× `package services_test`; im Paket des Codes nur je `export_test.go` | bestätigt |
| Testliste vor und nach gleich | `go test -list .` für beide Pakete an `a12f16f` und `ebbf5b8`: je 62 Zeilen, `diff` leer | bestätigt (62 Tests) |
| Abdeckungs-Deklarationen unverändert | alle `Abdeckung:`-Zeilen unter beiden Pfaden sortiert: je 43, `diff` leer; zusätzlich die Zuordnung Test → Deklaration (Kommentarblock über `func Test…`) an beiden Ständen gleich | bestätigt |
| `make abdeckung-check` grün | Exit 0, auch im Gate-Lauf | bestätigt |
| Unexportiertes nur über die Brücke nach Punkt 7 (neue Fassung) | siehe Abschnitt 2 | bestätigt |
| keine Prüfung entfallen | Testliste gleich; der Fall „ohne Interaktion 0“ ist akzeptiertes Negativ (Abschnitt 3); Mutationen Abschnitt 4 | bestätigt |

### Punkt 2 — `make lint` ohne Befund in den Testdateien. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| 31 Befunde vorher in den Testdateien | `a12f16f` (Kopie): 103 Befunde gesamt, davon 31 unter `internal/hexagon/{model,services}/*_test.go`: `testpackage` 9 (je eine Zeile in den neun Testdateien), `contextcheck` 5 (`services/record_extended_test.go:792`, `:793`, `:815`, `:816`, `:817`), `gochecknoglobals` 17 (`services/record_extended_test.go:17`–`:29` zwölf, `services/replay_extended_test.go:24`–`:27` vier, `services/replay_unverbraucht_test.go:21` eine) | bestätigt, gleich der Aufteilung der ADR |
| 0 nachher | `make lint` im Arbeitsbaum an `ebbf5b8`: Exit 2, 72 Befunde; unter beiden Pfaden 12 Zeilen, alle im Produkt-Code (`model/extended.go`, `model/recording.go`, `services/replay.go`), keine in einer `*_test.go` | bestätigt |
| kein neuer Befund | `comm -13` der sortierten Befundzeilen `a12f16f` gegen `ebbf5b8`: leer; 103 − 31 = 72 | bestätigt |
| ohne `//nolint`, ohne Änderung an `.golangci.yml` | `grep -rn nolint internal/` leer; `git diff a12f16f ebbf5b8 -- .golangci.yml` leer | bestätigt |
| kein `_ =` vor einem Fehler, den `errcheck` meldet (§6) | Blank-Zuweisungen in den Testdateien an beiden Ständen verglichen (Paketpräfix entfernt): gleich bis auf ein zusätzliches `id, _ := s.OpenConnection(ctx)` in `TestReplayLetzteNummer`. `OpenConnection` liefert `(SessionID, []Response)`, keinen Fehler. Die zwei `_ = s.CloseSession(…)` in `record_extended_test.go` sind Bestand und waren nie ein `errcheck`-Befund | bestätigt |

### Punkt 3 — `make gates` grün. **Bestätigt.**

`make gates` im Arbeitsbaum an `ebbf5b8`: Exit 0. Im Log u. a. `abdeckung-gegenprobe`, `a-check-negativ`, `baseline-verify` (54 Dateien), `commit-msg-gegenprobe`, `run-integration-tests`, `kopf-check-gegenprobe` je grün.

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| Review mit Report | liegt vor (F-441 bis F-443); die Nacharbeit `ebbf5b8` hat kein eigenes Review, Abschnitt 2 prüft sie |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen: §7 trägt Platzhalter, §6 *Risiken* drei Ausgänge „— (bei Closure)“. Sache der Closure |

---

## 2. F-441 als erneutes Review — hält die Umsetzung `SPEC-049` Punkt 7 in der neuen Fassung?

**Ja.** Die neue Fassung (`b0b2286`) sagt: Einen Wert eines unexportierten Typs legt nur der Produkt-Code an, weder Brücke noch Test, auch nicht lokal, als Literal oder über einen Typ-Alias.

| Prüfung | Beleg | Urteil |
|---|---|---|
| Brücke `model` | `internal/hexagon/model/export_test.go`: nur `Klasse(ziffer int) string { return klassen[ziffer] }` — Index in die Tabelle des Pakets, kein neuer Wert | hält |
| Brücke `services` | `internal/hexagon/services/export_test.go`: `IstLebendpruefung`, `VTLeerraum`, `Abweichung`, `ParameterStelle` rufen je die Funktion des Pakets; Parameter und Rückgaben sind `string`, `bool`, `map[string]string`, `model.ClientMessage`, `[]model.Value`. `LetzteNummer` mit `&cursor{…}` ist entfernt | hält |
| kein Alias, keine Variable, keine Konstante, kein `new(`, kein Literal eines kleingeschriebenen Typs in beiden Brücken | `grep -nE '\btype\b\|\bnew\(\|&?[a-z]\w*\{\|\bvar\b\|\bconst\b'` auf beide Dateien: leer | hält |
| Tests | Die Tests sind `_test`-Pakete; einen unexportierten Typ erreichen sie nur über einen Alias der Brücke, und es gibt keinen. `grep -nE '(services\|model)\.[a-z]'` in allen Testdateien: leer | hält |
| `TestReplayLetzteNummer` | prüft über `NewReplayService`, `OpenConnection` und `Query` (Session mit den Interaktionen 2 und 4, dritte Anfrage `PGR-E5001` „nach Interaktion 4 …“), wie §6 *Brücke mit Cursor* es festlegt | hält |
| Plan nachgezogen | §3 (Zeilen `export_test.go`, `replay_lebendpruefung_test.go`), §6 (*Export-Test-Brücke*, *White-Box-Zugriffe*, *Brücke mit Cursor*), §7 (zwei Mutationszeilen über `ReplayService`) | hält |

**Gegenprobe: Regel aus Punkt 7 gebrochen** (`p7-bruch`, Kopie von `ebbf5b8`). In `services/export_test.go` ergänzt: `type Cursor = cursor` und `LetzteNummer` mit `(&cursor{session: session}).letzteNummer()`; dazu `bruch_test.go` in `services_test` mit `services.Cursor{}` und `services.LetzteNummer(&model.Session{ID: 1})` (der frühere Fall „ohne Interaktion 0“). Ergebnis: `go vet` still, `go test -run '^TestBruchPunkt7$'` **PASS**, `make lint`-Stufe 72 Befunde, `comm -3` gegen den Grundlauf leer. Der Bruch ist also für Test, `vet` und Lint unsichtbar — genau die *Grenze* von `SPEC-049`. Halt ist allein Review und Verifikation; das bestätigt die Zeile in `docs/plan/planning/next/slice-harness-lint.md` („kein Rot-Fall, Urteil des Review“).

**F-441:** behoben. Die Auslegung hat einen Entscheider (Architect, `b0b2286`) und einen dauerhaften Ort (`SPEC-049` Punkt 7 und *Grenze*); die drei übrigen Umstellungs-Slices und `slice-harness-lint` tragen den Hinweis in §6. **F-442:** erledigt — die Kopplung an `status: "I"` entfällt mit der Brücke. **F-443:** bestätigt, Abschnitt 4.

---

## 3. Akzeptiertes Negativ „ohne Interaktion 0“ — ist der Zweig über die Schnittstelle unerreichbar?

**Ja.** Der Zweig `len(c.session.Interactions) == 0` in `(*cursor).letzteNummer` setzt einen Cursor mit einer Session ohne Interaktion voraus.

1. `c.session` wird nur in `zuordnen` gesetzt (`replay.go:254`–`:261`), und zwar aus `s.frei[0]`; `OpenConnection` legt Cursor mit `session == nil` an (`replay.go:113`). Andere Schreibstellen auf `.session` gibt es im Paket nicht (Suche nach `session:`, `.session =`).
2. `s.frei` füllt nur `NewReplayService`, und nur mit Sessions, deren Interaktionen nach `ohneLebendpruefungen` nicht leer sind (`replay.go:81`–`:84`). Ist keine solche da, endet der Aufbau mit `PGR-E3004`.
3. `Interactions` einer zugeordneten Session schreibt danach nichts mehr: Die einzigen Zuweisungen auf `.Interactions` im Paket stehen in `record.go` (Aufzeichnung) und an `replay.go:81` (vor der Aufnahme in `s.frei`). `ohneLebendpruefungen` liefert eine neue Liste, ein späteres Ändern des geladenen Werts durch das Repository erreicht sie nicht.
4. `letzteNummer` wird nur in `Query` (`:153`) und `ClientMessage` (`:211`) gerufen, beide nach `zuordnen`.

Damit ist der Fall über `ReplayService` nicht herzustellen; der frühere Test prüfte einen Zustand, den das Produkt nicht erzeugt. Die Begründung in §6 trägt. Eine Folge davon steht in V-81.

---

## 4. Mutationen nach `AGENTS.md` §3.10 — Stichprobe gegen §7

Grundlauf beider Pakete an `ebbf5b8`: grün. Jede Mutation in eigener frischer Kopie, `go test -count=1 -run` auf den genannten Test.

| Zeile in §7 | Mutation (wie ausgeführt) | Ergebnis | Grund in der Ausgabe |
|---|---|---|---|
| letzte Nummer ist die aufgezeichnete | `letzteNummer`: `return len(c.session.Interactions)` | rot, `TestReplayLetzteNummer` | „nach Interaktion 2 erwartet …“ statt 4 |
| letzte Nummer ist die der letzten Interaktion | `letzteNummer`: `Interactions[0].Sequence` | rot, `TestReplayLetzteNummer` | „nach Interaktion 2 erwartet …“ statt 4 |
| jeder Fehlercode ergibt eine Klasse (Brücke `Klasse`) | Zeile `6:` aus `klassen` gestrichen | rot, `TestCodeTabelle` | „CodeUnsupported = PGR-E6001 ergibt keine Klasse“ |
| eine Kette trägt den Code des klassifizierten Fehlers | `Meldungen`: `Code: CodeInternal` statt `me.Code` | rot, `TestFehlerKette` | `Code:"PGR-E1000"` statt `PGR-E4002` |
| `;` macht keine Lebendprüfung (Brücke `IstLebendpruefung`) | `leerraum` um `;` erweitert | rot, `TestLebendpruefungErkennung` | `istLebendpruefung(";", false) = true` |
| `\v` ab Hauptversion 17 (Brücke `VTLeerraum`) | `>=` zu `>` | rot, `TestLebendpruefungServerversion` | `server_version "17.0" … = false, erwartet true` |
| Query nach dem Ende liefert `ErrSessionEnded` | Prüfung `beendet` nach dem Upstream-Aufruf in `Query` entfernt | rot, `TestRecordEndeNachErfolgreichemUpstream/Query` | `Query: <nil>` |
| mit `FailOnUnconsumed` ist Unverbrauchtes ein Fehler | `if false && s.streng` | rot, `TestReplayNichtVerbrauchtMeldung` | `PGR-W2001 … erwartet PGR-E5002` |

Acht von sechzehn Zeilen, darunter beide `letzteNummer`-Mutationen: alle rot aus dem richtigen Grund. Zusammen mit den sechs eigenen Mutationen des Reviews (R1 bis R5 weiter gültig; R6 lief noch über die Brücke und ist durch die erste Zeile oben ersetzt) sehe ich keine Prüfung, die im Umbau verloren ging.

---

## 5. §6 — hat jeder Punkt eine Entscheidung mit Ort?

| Punkt | Entscheidung · Ort | Urteil |
|---|---|---|
| Export-Test-Brücke | `SPEC-049` Punkt 7 (neue Fassung `b0b2286`), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 | ja |
| Übrige Befunde der Testdateien | §6 selbst, vor dem Code (Stand `a12f16f`, Architect-Prüfung nach §4); zitiert „Entscheidung 5“. Die ADR trägt `_ =`-Verbot und „keine Ausnahme nur für Bestand“, nicht aber `t.Context()`; sie nennt für neue Kontexte `context.WithoutCancel` (V-83) | ja, Ort ungenau |
| White-Box-Zugriffe im Bestand | je Zugriff entschieden: vier über die Brücke, `cursor`/`letzteNummer` über die Schnittstelle (*Brücke mit Cursor*) | ja |
| Fakes der Ports | Ports exportiert; Fakes greifen auf kein unexportiertes Feld (Kompilat der `_test`-Pakete) | ja, ohne offenen Fall |
| Mutationstests auf unexportierte Teile | `AGENTS.md` §3.10; Tabelle in §7 | ja |
| Brücke mit Cursor | `SPEC-049` Punkt 7, Architect 2026-10-06 nach F-441/F-442; akzeptiertes Negativ mit Grund | ja |
| Risiken (drei) | Ausgang „— (bei Closure)“ | offen bis zur Closure, zulässig |

---

## 6. Plan gegen Code

- **Kopf:** `Bezug` stimmt. `Berührte Spec-Stellen` nennt `SPEC-049` Punkt 7 und 8 mit dem Zusatz „der Slice ändert die Stelle nicht“ — `b0b2286` hat Punkt 7 und *Grenze* unter der Kennung dieses Slice geändert (V-80). `make kopf-check` ist grün; den Zusatz prüft er nicht.
- **§1:** Ziel stimmt. Die Abgrenzung „Produkt-Verhalten, Spezifikation, Lastenheft — Der Slice ändert Testdateien“ ist durch `b0b2286` überholt (V-80). Kein Produkt-Code im Diff, kein neuer Export, `.golangci.yml` unverändert.
- **§3:** Jede Zeile hat ihre Datei im Diff; die Zeilen zu `export_test.go` und `TestReplayLetzteNummer` folgen `ebbf5b8`.
- **§5:** „`make gates` grün“ und „`make lint` ohne Befund in den Testdateien“ sind gemessen (Abschnitt 1).
- **§6, §7:** siehe Abschnitte 2 bis 5.
- **Folge-Slices:** `slice-harness-blackbox-driven`, `-pgwire`, `-einstieg` und `slice-harness-lint` tragen den Hinweis zu Punkt 7 (`b0b2286`).

---

## 7. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-80 | Plan folgt Korrektur nicht (`AGENTS.md` §3.9, `BEO-REPO/plan-folgt-korrektur-nicht`) | `b0b2286` ändert `SPEC-049` Punkt 7 und *Grenze* und trägt die Kennung dieses Slice. Kopf und §1 sagen weiter das Gegenteil: „der Slice ändert die Stelle nicht“ und „Spezifikation … Der Slice ändert Testdateien“. §3.9 verlangt, dass wer im Slice die Spezifikation ändert, im selben Commit Kopf und §1 nachzieht; §6 nennt die Entscheidung, Kopf und §1 nicht. Vor der Closure nachziehen, etwa: Punkt 7 in der Fassung des Architects nach F-441/F-442, und §1 nimmt die Spezifikation bis auf diese Stelle aus. Zugleich Retirement-Check der Hard Rule: die Beobachtung ist wieder aufgetreten | Plan Kopf `Berührte Spec-Stellen`; §1 letzter Ausschluss | `git show b0b2286 --stat` (u. a. `spec/spezifikation.md`, 12 Zeilen); Kopf am Stand `ebbf5b8` |
| V-81 | Zusage im Kommentar weiter als Prüfung (`AGENTS.md` §3.11, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) | Der Kommentar über `(*cursor).letzteNummer` sagt „ohne erwartete Interaktion 0“ zu. Seit `ebbf5b8` prüft das kein Test mehr, und nach Abschnitt 3 kann es keiner über die Schnittstelle. Das akzeptierte Negativ ist richtig; die Zusage im Produkt-Code bleibt aber stehen, und §6 gibt ihr keine Adresse („kein eigener Folge-Slice, die Abdeckung misst `slice-harness-coverage`“ — eine Zahl, keine Entscheidung über Zweig und Kommentar). Kein Produkt-Code in diesem Slice (§1); Adresse vorschlagen, etwa `slice-lint-bestand-kern-driven` (Produkt-Code dieser Pakete): Zweig und Satz streichen oder als Invariante benennen | `internal/hexagon/services/replay.go:185`–`:193`; Plan §6 *Brücke mit Cursor* | Abschnitt 3; `p7-bruch` zeigt, dass nur ein selbst angelegter Cursor den Zweig erreicht |
| V-82 | Beleg fehlt im Eingang | Liefer-Punkt 2 verlangt, dass „die Zeilen der Ausgabe unter diesen Pfaden vor und nach dem Umbau im Bericht stehen“. Der Bericht des Implementers lag mir nur als Zusammenfassung (Zahlen) vor, nicht mit Zeilen. Abschnitt 1 Punkt 2 ersetzt ihn durch eigene Läufe mit Datei und Zeile; für die Closure kann §7 auf diesen Bericht zeigen | DoD Punkt 2 | — |
| V-83 | Hinweis (Ort einer Entscheidung ungenau) | §6 *Übrige Befunde* begründet `t.Context()` mit „Entscheidung 5“. Die ADR trägt dort für neue Kontexte `context.WithoutCancel` als erwartete Bereinigung und schweigt zu Tests. Für einen Test ist `t.Context()` sachlich richtig (beide Subtests lesen ihr Ergebnis vor dem Ende des Subtests, Review Punkt 3), und der Slice liefert keinen neuen Vertrag, §3.12 greift also nicht. Die drei übrigen Umstellungs-Slices treffen dieselbe Frage; ein Satz in `SPEC-049` oder in deren §6 („in Testdateien `t.Context()`“) erspart ihnen die eigene Auslegung | Plan §6 *Übrige Befunde der Testdateien* | [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 5, letzter Satz |

---

## 8. Gate-Ergebnis und Urteil

`make gates` an `ebbf5b8` im Arbeitsbaum: **grün** (Exit 0). `make lint`: Exit 2, 72 Befunde, keiner in einer Testdatei unter `internal/hexagon/model` oder `internal/hexagon/services`, keiner neu gegenüber `a12f16f`.

**Urteil:** Die drei Liefer-Punkte sind erfüllt und selbst belegt: 62 Tests gleich, 43 Abdeckungs-Deklarationen gleich und gleich zugeordnet, 31 → 0 Befunde in den Testdateien, modulweit 103 → 72 ohne neuen Befund, kein `//nolint`, kein neues `_ =`. F-441 ist behoben; Brücken und Tests halten `SPEC-049` Punkt 7 in der neuen Fassung, ohne selbst angelegten Wert eines unexportierten Typs und ohne Alias. Das akzeptierte Negativ trägt. Die Stichprobe von acht Mutationen ist rot aus dem richtigen Grund.

**Vor der Closure:** V-80 nachziehen (Kopf und §1), V-81 eine Adresse geben. V-82 und V-83 sind Hinweise für die Closure-Notiz und die übrigen Umstellungs-Slices. Wiederkehrende Klassen dieses Laufs: „Plan folgt Korrektur nicht“ (V-80), „Zusage im Kommentar weiter als Prüfung“ (V-81).
