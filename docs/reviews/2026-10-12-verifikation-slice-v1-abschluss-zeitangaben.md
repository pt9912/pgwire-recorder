# Verifikation: slice-v1-abschluss-zeitangaben — 2026-10-12

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-12-review-slice-v1-abschluss-zeitangaben.md`, F-653 bis F-659); die Nacharbeit und die Entscheidung des Architect prüfe ich nach. Ich repariere nichts.

**Gegenstand:** Slice-Plan `slice-v1-abschluss-zeitangaben` (Kopf, §1 bis §8) gegen den Gesamt-Diff `d44154c..ed55527` (20 Dateien, 1802 Zeilen hinzu, 23 weg). Implementer: `61a63d5`, `2edde5c`, `ed55527`. Review: `803c0ac`. Entscheidung des Architect: `5dac153`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

**Eingang:** der Plan ganz am Stand `ed55527` (§7 *Belege des Implementers* nur als Beleg-Teile; §7 ist bis zur Closure unvollständig, das ist kein Befund); [`LH-FA-21`](../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen); [`LH-FA-21.a`](../../spec/spezifikation.md#lh-fa-21a--zeitangaben) (*Einspielen*, *Uhr*, *Randformen des Einspielens* mit *Prüfung der Werte*); [`LH-FA-17.a`](../../spec/spezifikation.md#lh-fa-17a--konfiguration); [`LH-FA-20.a`](../../spec/spezifikation.md#lh-fa-20a--einspielen); [ADR-0015](../plan/adr/0015-uhr-port.md); der Review-Report; der Code `internal/hexagon/services/play_zeit.go`, `play.go`, `internal/adapters/driving/cli/zeitangaben.go`, `cli.go`, `internal/bootstrap/bootstrap.go`, `internal/hexagon/model/fehler.go`. Den Bericht des Implementers habe ich nicht gelesen. Beim Start war der Arbeitsbaum sauber, HEAD `ed55527`.

**So wurden die Proben gebaut.** Alles Docker-only, Scratch im Scratchpad der Sitzung, nie im Repo.

- *Mutanten:* je Mutant eine frische Kopie des Baums (`cp -r` ohne `-p`, ohne `.git`), genau eine Ersetzung (das Skript bricht ab, wenn sie nichts oder mehr als eine Stelle trifft), `gofmt -l` im selben Lauf (nie ein Treffer), `go test -count=1` im Go-Image des `Dockerfile` (Stufe `deps`) mit `--network none` und dem Bau-Cache als Bind-Mount. Die Mutanten sind unabhängig vom Implementer neu gebaut, im Wortlaut von [`LH-FA-21.a`](../../spec/spezifikation.md#lh-fa-21a--zeitangaben), nicht aus §7 übernommen.
- *Ende zu Ende:* das aus dem Stand gebaute Binary (Stufe `integration` des `Dockerfile`) in einem eigenen internen Docker-Netz gegen das gepinnte PostgreSQL-Image (Digest aus `harness/mk/integration.mk`), Aufzeichnen mit `record --record-timing` über einen Recorder-Container, Client `psql` im Server-Container mit Pausen (`sleep`) zwischen den Anweisungen, die Zeit am Server gemessen mit `clock_timestamp()` und `backend_start`. Aufgeräumt per Name: `docker rm -fv` der Container, das Netz, die zwei Images dieser Sitzung.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `--keep-timing` hält die Pausen ein (Abnahmeszenario 13). **Bestätigt.**

Die Tests existieren: Jeder Name aus §7 und der DoD (36) hat genau eine Fundstelle (`grep "^func Name("`). Ohne den Fix aus dem richtigen Grund rot (Service, `go test ./internal/hexagon/services`; jeder Mutant nur eine Zusage):

| Zusage (Wortlaut `LH-FA-21.a`) | Mutant | rot, mit Test |
|---|---|---|
| `relative`: seit dem Beginn der vorigen Interaktion | Bezug ist die Lesung beim Warten | `TestPlayZeitRelative`, `TestPlayZeitNachFehlerantwort` |
| `relative`: die Verspätung wird nicht aufgeholt | Bezug ist der Beginn der ersten Interaktion | `TestPlayZeitRelativeVerspaetung`, `TestPlayZeitNachFehlerantwort` |
| `relative`: vor der ersten Interaktion nie | die erste wartet bis zu ihrem Wert | `TestPlayZeitRelativeVerspaetung`, `TestPlayZeitSessionRelative`, drei weitere |
| `absolute` mit `connect`: vor der ersten bis zum `offset_ms` | die erste wartet nie | `TestPlayZeitAbsoluteConnect`, `TestPlayZeitSessionAbsolute`, `TestPlayZeitGrosserWert` |
| `absolute` mit `first-request`: vor der ersten nie | die erste wartet | `TestPlayZeitAbsoluteErsteAnfrage` |
| `first-request` ist ein eigener Bezugspunkt | im Service gilt immer `connect` | `TestPlayZeitAbsoluteErsteAnfrage` |
| `absolute` ist ein eigener Modus (Verspätung aufgeholt) | der Modus ist nie `absolute` | `TestPlayZeitAbsoluteConnect`, `TestPlayZeitSessionAbsolute`, `TestPlayZeitGrosserWert` |
| `connect` ist der Beginn des Aufbaus, vor TLS und Anmeldung | die Lesung steht nach `Verbinde` | `TestPlayZeitAbsoluteConnect` |
| `connect` wird am Beginn jedes Aufbaus gelesen | eine Lesung für alle Sessions (Paketvariable) | `TestPlayZeitSessionAbsolute` |
| Beginn einer Interaktion: nach Ablauf der Wartezeit | der Beginn wird vor dem Warten festgehalten | neun Tests, darunter `TestPlayZeitRelative`, `TestPlayZeitNachFehlerantwort` |
| Scheitert der Aufbau, gibt es keine Wartezeit | die erste Interaktion wartet vor `Verbinde` | `TestPlayZeitAufbauScheitert`, `TestPlayZeitAbsoluteConnect`, `TestPlayZeitAbbruchMitFinish` |
| Nach dem ersten Signal wird die Uhr nicht mehr aufgerufen | `warte` prüft `ctx` nicht vor dem Aufruf | `TestPlayZeitAbbruchMitFinish` |
| Das erste Signal beendet jede Wartezeit sofort | die Wartezeit hängt nicht an `ctx` | `TestPlayZeitAbbruchWartezeit`, `TestPlayZeitAbbruchMitFinish`, `TestPlayZeitGrosserWert` |
| Das zweite Signal beendet sie ebenso | die Wartezeit hängt nicht an `abbruch` | `TestPlayZeitZweitesSignalInWartezeit` |
| ohne die Option wird die unterbrochene Interaktion nicht gesendet | sie wird gesendet | `TestPlayZeitAbbruchWartezeit` |
| mit `--finish-session-on-interrupt` wird sie gesendet | sie wird nicht gesendet | `TestPlayZeitAbbruchMitFinish` |
| `Play` kehrt erst nach dem Schließen zurück | `<-geschlossen` entfernt | `TestPlayZeitZweitesSignalInWartezeit` |
| Überlauf: begrenzt, nie negativ | keine Begrenzung | `TestPlayZeitGrosserWert` |
| Ohne `--keep-timing` wartet `play` nie und liest die Uhr nie | der Takt gilt mit Uhr auch ohne die Option | `TestPlayZeitOhneKeepTiming` |
| Das zweite Signal in der Wartezeit sendet die Interaktion nicht | ein Fehler des Wartens zählt als Erfolg, solange `ctx` lebt | `TestPlayZeitZweitesSignalInWartezeit` |
| Der Bootstrap verdrahtet die Systemuhr | `services.Wartezeiten(NeueSystemuhr())` entfernt | `TestRunPlayKeepTiming` (Nil-Zugriff) |

Dazu am Binary (`make test-integration`-Image, Tests dieses Slice, eigene Mutanten, jeweils ein neues Integrationsimage): der Takt nie aktiv macht `TestE2EPlayKeepTiming` und `TestE2EPlayKeepTimingQuellen` rot; `PGR-E2003` nie geprüft macht `TestE2EPlayKeepTimingOhneOffset` rot. Ohne Mutation liefen alle sechs Integrationstests des Slice (`TestE2EPlayKeepTiming*`, `TestE2EReplayIgnoriertZeitangaben`) grün.

**Ende zu Ende am gebauten Binary** (Aufzeichnung mit `record --record-timing`, drei Interaktionen, davon eine Extended-Interaktion über `\bind`; `offset_ms` 4, 552, 954 in `yaml`, 4, 555, 958 in `sqlite`; Zeit am Server gemessen):

| Probe | Ergebnis |
|---|---|
| `yaml` und `sqlite`, `relative`, `absolute` (`connect`), `absolute` mit `first-request` | Exit 0; Abstände der Marken 551 und 402 ms, 552 und 953 ms (zum Backend-Beginn), 552 und 954 ms — die Untergrenzen der Aufzeichnung gehalten (kleinste Unterschreitung 1 ms, siehe V-196) |
| ohne `--keep-timing` | Abstände 13 und 3 ms: kein Warten |
| `relative`, erste Interaktion bei `offset_ms` 3000 (zweite 3100) | die erste läuft nach 6 ms; die zweite 97 ms später: die erste wartet nie |
| `absolute` (`connect`), dieselbe Datei | die erste läuft nach 3004 ms (3000 aufgezeichnet), die zweite 95 ms später |
| `absolute` mit `first-request`, dieselbe Datei | die erste nach 5 ms, die zweite 98 ms später: die erste wartet nie |
| `relative`, erste Interaktion dauert 1,5 s (`pg_sleep`), zweite bei 500, dritte bei 1000 | zweite 1506 ms nach der ersten (Verspätung nicht aufgeholt), dritte 500 ms nach der zweiten (wartet) |
| `absolute` (`connect`) und `first-request`, dieselbe Datei | zweite 1515 ms nach der ersten, dritte 3 ms nach der zweiten (Verspätung aufgeholt) |
| erstes Signal (SIGINT) in der Wartezeit (nächste Interaktion bei 60 s), ohne `--finish-session-on-interrupt` | Exit 0 nach 188 ms; nur die erste Anweisung auf dem Server |
| dasselbe mit `--finish-session-on-interrupt` | Exit 0 nach 206 ms; die übrigen Interaktionen laufen sofort zu Ende (Marken 2 und 3 auf dem Server) |
| `replay` auf einer Aufzeichnung mit `offset_ms` 60 000 | Client fertig nach 58 ms, Exit 0: Zeitangaben ignoriert |

### Punkt 2: Optionen, Variablen, Schlüssel, `PGR-E2001`. **Bestätigt.**

Mutanten im CLI-Adapter (`go test ./internal/adapters/driving/cli`), alle rot, jeder nur eine Zusage:

| Zusage | Mutant | rot, mit Test |
|---|---|---|
| `--timing-mode` ohne `--keep-timing` ist `PGR-E2001` | der Modus aus der Prüfliste | `TestPlayZeitVerwendung` |
| `--timing-reference` ohne `--keep-timing` ist `PGR-E2001` | der Bezug aus der Prüfliste | `TestPlayZeitVerwendung` |
| `--timing-reference` mit `relative` ist `PGR-E2001` | die Prüfung entfällt | `TestPlayZeitVerwendung` |
| „gesetzt“ heißt: irgendeine Quelle nennt den Wert, auch den Default | „gesetzt“ heißt: vom Default verschieden | `TestPlayZeitVerwendung` |
| nur Kleinbuchstaben: `Absolute` | zugelassen | `TestPlayZeitWerteUngueltig` |
| `first_request` | zugelassen | `TestPlayZeitWerteUngueltig` |
| der leere Wert (Bezug, Modus je ein Mutant) | zugelassen | `TestPlayZeitWerteUngueltig` |
| Meldung nennt den Wert von Variable und Datei nicht | der Text hängt den Wert an | `TestPlayZeitWerteUngueltig` |
| Default `relative` / `connect` | `absolute` / `first-request` | `TestPlayZeitWerte`, `TestPlayZeitQuellen`; der Modus auch `TestPlayZeitVerwendung` |
| Option, Variable und Schlüssel gelten | Modus und Bezug aus der Quelle ignoriert (je ein Mutant) | `TestPlayZeitWerte`, `TestPlayZeitQuellen` |
| bei `record` und `replay` unbekannt | `--keep-timing` bei `record`, `--timing-mode` bei `replay` angemeldet | `TestDateiUngueltig`, `TestConfigShowFehler`, `TestLeserOptionen`, `TestLeserFremdeUmgebung` (und weitere) |
| `--help` nennt `PGR-E2003`, „Standard false“, den Beginn der vorigen und der ersten Interaktion | je ein Textmutant | `TestPlayZeitHilfe` (vier Mutanten, alle rot) |

**Am Binary** (gebautes Binary, `play` mit unerreichbarem Upstream, damit nur die Prüfung der Verwendung spricht): `--timing-mode` oder `--timing-reference` ohne `--keep-timing` (Option, Variable, Datei), `--timing-reference connect` mit `relative` (Default und ausdrücklich), `--keep-timing=false` mit `--timing-mode absolute` in beiden Reihenfolgen, `Absolute`, `RELATIVE`, `first_request`, der leere Wert und die Kurzform `-k` sind `PGR-E2001`, Exit 2, ohne Verbindung. Der Wert aus Variable und Datei steht nie in der Meldung (`GEHEIM` als Probe). Die Schlüssel `keep_timing`, `timing_mode`, `timing_reference` mit ungültigem Wert (`GEHEIM`, `Absolute`, `first_request`, `1`) sind `PGR-E2004`; `record.keep_timing`, `replay.timing_mode`, `record.timing_mode` und `timing_reference` oben sind `PGR-E2004` unbekannt; `record --keep-timing`, `record --timing-mode`, `replay --timing-reference` und `play --record-timing` sind `PGR-E2001`. Die Umgebungsvariable `PGWIRE_RECORDER_KEEP_TIMING=GEHEIM` stört `record` nicht (fremde Variable), wohl `play`. `--keep-timing` auf der Kommandozeile geht `PGWIRE_RECORDER_KEEP_TIMING=false` vor (der Lauf erreicht den Verbindungsaufbau). `play --help` nennt die drei Optionen, Standards, Variablen und `PGR-E2001`/`PGR-E2003`; `config show` nennt die drei Schlüssel der Datei und den Namen der gesetzten Variablen ohne Wert.

### Punkt 3: `PGR-E2003`. **Bestätigt.**

| Zusage | Mutant | rot, mit Test |
|---|---|---|
| jede Interaktion wird geprüft | nur die erste jeder Session | `TestPlayZeitFehltOffset` |
| jede Session mit Interaktionen | nur die erste Session | `TestPlayZeitFehltOffset` |
| die Meldung nennt die erste ohne Wert, nicht die letzte | die letzte | `TestPlayZeitFehltOffset` |
| Session und Interaktion in der richtigen Reihenfolge | vertauscht | `TestPlayZeitFehltOffset` |
| Klasse Konfiguration, Exit-Code 2 | Code `PGR-E3009` | `TestPlayZeitFehltOffset`, `TestRunPlayKeepTimingOhneOffset` |
| ohne `--keep-timing` entfällt die Prüfung | die Prüfung läuft immer | `TestPlayZeitFehltOffset`, `TestPlayZeitOhneKeepTiming` (und der Bestand) |
| beschädigte Interaktion (`PGR-E3003`) geht `PGR-E2003` vor | die Prüfung auf `offset_ms` steht vor `Validate` | `TestPlayZeitFehltOffset`, `TestPlayZeitFehltOffsetRand` |

**Am Binary:** Fehlt der zweiten Interaktion `offset_ms`, ist es `Konfiguration [PGR-E2003]: Session 1, Interaktion 2 …` auf einer Zeile, Exit 2, ohne Log-Zeile und ohne Wert, vor jeder Verbindung (Upstream unerreichbar, sonst wäre es `PGR-E4002`); mit zwei Sessions nennt die Meldung die erste Session mit der Lücke (`Session 2, Interaktion 1`); dieselbe Meldung mit `absolute` und `first-request`; ohne `--keep-timing` spielt dieselbe Datei ein (Exit 0, drei Marken). Die Ordnung nach der Entscheidung des Architect: eine Aufzeichnung mit beschädigter dritter Interaktion (endet nicht mit `ready_for_query`) und fehlendem `offset_ms` an der zweiten ist mit `--keep-timing` `Recording [PGR-E3003]: Session 1, Interaktion 3 …`, Exit 3, ohne die Option dieselbe Meldung; ein negativer `offset_ms` ist `PGR-E3003` (Laden). Eine nicht lesbare Datei bleibt `PGR-E3001`.

### Punkt 4 bis 9: Gate, Review, Closure, Register, Risiken, Paarungen. **Nicht Gegenstand dieses Laufs** (Closure-Pflichten; §7 vor der Closure unvollständig, kein Befund). Das Review-Häkchen: Der Report liegt vor; seine Nacharbeit prüfe ich unten. `make gates`: siehe Verdikt.

---

## 2. Nacharbeit zum Review und Entscheidung des Architect

- **F-653 (Ordnung `PGR-E3003` vor `PGR-E2003`):** Die Entscheidung steht in [`LH-FA-21.a`](../../spec/spezifikation.md#lh-fa-21a--zeitangaben) *Prüfung der Werte* (Spec-Commit `5dac153`, Technik-Stratum, kein Lastenheft-Bump nötig) und in Plan §6. Der Code folgt ihr (`Validate`-Schleife vor `pruefeZeitangaben`); sie ist doppelt gefangen (Mutant oben, Binary-Probe). Behoben.
- **F-654:** Der Kommentar von `PlayOptions` sagt zu, der Service erwarte genau die gültigen Werte, und sagt nichts mehr über andere zu. Die zwei Mutanten bleiben grün und sind für jeden zugelassenen Wert äquivalent (beide Zweige stimmen für `relative`/`absolute` und `connect`/`first-request` überein). Der CLI-Adapter lässt den Bezug bei `relative` leer (`TestPlayZeitWerte`). Behoben, die Begründung trägt.
- **F-655:** `TestPlayZeitHilfe` prüft „Standard false“ und die Beschreibung der Modi und Bezugspunkte; vier Textmutanten rot (oben). Behoben.
- **F-658, F-657, F-659:** akzeptierte Negative beziehungsweise Umbruch; im Plan geführt. Der Kommentar von `Wartezeiten` nennt die Pflicht des Aufrufers; der Mutant „Verdrahtung entfernt“ ist rot. Kein Befund hier.
- **Eingriff `<-geschlossen` im `defer` von `session`:** Im Play-Pfad ohne Zeitangaben läuft der Bestand grün (`go test` über alle drei Pakete), und am Binary: zwei Signale hintereinander in einer laufenden Anfrage (`pg_sleep(30)`, ohne `--keep-timing`) beenden `play` mit Exit 0 binnen 220 ms, ohne Hang. Kein Deadlock; der Mutant, der das Warten entfernt, wird rot (`TestPlayZeitZweitesSignalInWartezeit`), das Warten ist also getragen.

## 3. Spezifikation gegen Code

[`LH-FA-21.a`](../../spec/spezifikation.md#lh-fa-21a--zeitangaben) Satz für Satz gegen `play_zeit.go`, `play.go`, `zeitangaben.go`, `cli.go`:

| Zusage | Stelle | Befund |
|---|---|---|
| Uhr nur über den Port; der Kern bindet keine Zeitbibliothek | `play_zeit.go` importiert `context`, `math`, Modell, Ports | erfüllt (`make a-check` Teil des Gate-Laufs) |
| `relative`, `absolute`, `connect`, `first-request` | `takt.ziel`, `neuerTakt` | erfüllt (Mutanten, Binary-Proben) |
| Jede Session eigene Uhr, Sessions nacheinander | `neuerTakt` je `session` | erfüllt |
| Beginn der Interaktion = Senden der ersten Client-Nachricht nach der Wartezeit | `beginnt` nach `warte` | erfüllt |
| `PGR-E2003` nach dem Laden, nach `Validate`, vor jeder Verbindung, ohne Wert | `NewPlayService` | erfüllt |
| Sehr großer Wert: keine Obergrenze, Überlauf begrenzt, abbrechbar | `plus`, `warte` | erfüllt |
| Verbindungsaufbau scheitert: keine Wartezeit | `warte` steht nach `Verbinde` | erfüllt |
| ungültige Verwendung `PGR-E2001`; `PGR-E2004` für Schlüssel | `zeitPlay`, `zeitmodus`, `zeitbezug` | erfüllt |
| Port `Player` unverändert | `git diff` über `internal/hexagon/ports` leer | erfüllt |

Die Sicht (`spec/architecture.md`) ist unberührt; nichts darin verlangt eine Änderung.

## 4. Plan gegen Code, Hard Rules, Namensabgleich

- §1/§3/§6/§8 folgen dem Code bis auf V-193: drei Liefer-Punkte (Wartezeiten mit `PGR-E2003`, Optionen, Verwendung), zwei Schichten (Kern, CLI-Adapter); der Bootstrap trägt eine Zeile ohne Verzweigung und zählt als Verdrahtung nicht. Jeder Dateiname und jeder Bezeichner aus §3 und §7 steht im Code (`zeitPlay`, `neuerTakt`, `Wartezeiten`, `CodeTimingMissing`, `NeueSystemuhr`, `Vorruecken`, `Wartende`, `toleranz`); jeder Testname hat genau eine Fundstelle; die Einträge des Beobachtungs-Registers aus §8 und §6 liegen unter `docs/plan/planning/observations/BEO-REPO/`.
- Kein `//nolint` in `cmd/`, `internal/`, `test/`. Keine Kennung eines Spezifikations- oder Sichtstratums in den Commit-Messages von `d44154c..HEAD` ([ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md)); jede nennt `slice-v1-abschluss-zeitangaben` und `LH-FA-21`.
- Handbuch und `README.md` unberührt (`git diff d44154c HEAD -- docs/user README.md` leer); der Slice liefert sie über `slice-v1-abschluss-zeitangaben-doku`. Dieser Nehmer nennt `slice-v1-abschluss-zeitangaben` in Ziel, Abgrenzung und DoD (nimmt an, `AGENTS.md` §3.13).
- Roadmap: „In Arbeit“ nennt den Slice, kein Chronik-Text.
- `make abdeckung-check`, `make kopf-check`, `make docs-check`, `make baseline-verify`: Teil des Gate-Laufs am Ende.

## 5. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-193 | LOW | §3 des Plans folgt dem Code nicht in den Testdateien: Er nennt `play_zeit_test.go` (Service, CLI, Bootstrap), den E2E-Test und `TestDateiUngueltig`, nicht aber die geänderten Tests `TestLeserOptionen`, `TestLeserAlleOptionen`, die Tabelle `tabelle()` und `TestLeserFremdeUmgebung` in `leser_test.go`, `TestParsePlayFremdeUmgebung` in `play_test.go` des CLI-Adapters und die Testhilfe `beiSchliesse` in `internal/hexagon/services/play_test.go`. Die Änderungen sind folgerichtig (kein Test geschwächt: die Platzhalter für die drei Optionen sind durch gleichwertige Fälle ersetzt) und im Diff des Reviews gesichtet; es fehlt nur die Zeile in §3 (`AGENTS.md` §3.9). | `docs/plan/planning/in-progress/slice-v1-abschluss-zeitangaben.md` §3, Zeile *Tests, `test/integration`* | Implementer (§3 ergänzen) vor der Closure |
| V-194 | INFO | Der Gate-Absatz in §7 nennt den Stand `61a63d5`. Der Stand `ed55527` unterscheidet sich im Code (Kommentare in `play.go`, ein Testfall in `play_zeit_test.go` des CLI-Adapters) und in der Spezifikation; der Absatz sagt „Spätere Commits ändern nur diesen Absatz“, ohne den späteren Stand zu nennen. Der Gate-Lauf am Stand `ed55527` ist belegt durch den Auftrag und durch meinen Lauf nach dem Commit dieses Berichts (siehe Verdikt); die Closure-Notiz nennt den Stand, auf dem sie sich beruft. | `docs/plan/planning/in-progress/slice-v1-abschluss-zeitangaben.md` §7 *Gate-Läufe* | Implementer, Closure |
| V-195 | INFO | Die Binary-Tests trennen weder `relative` von `absolute` noch `first-request` von `connect`: Die Binary-Mutanten „Modus nie `absolute`“ und „Bezug von `relative` ist die Lesung beim Warten“ blieben grün (nur Untergrenzen, Toleranz 50 ms, und der erste `offset_ms` der Probe liegt bei 4 ms). §7 führt unter *Grüne Mutanten* nur die Trennung `first-request`/`connect`. Die Unterscheidung tragen die Service-Tests (Mutanten oben alle rot), und die DoD verlangt am Binary die Untergrenze (Punkt 1); der Befund ist also keine Lücke der DoD, die Aussage in §7 ist nur enger zu fassen. Die Ende-zu-Ende-Proben oben (Pausen länger als die Antwort, erste Interaktion bei 3000 ms) trennen die Modi am Binary. | `test/integration/zeitangaben_einspielen_e2e_test.go` · `TestE2EPlayKeepTiming` | Implementer (eine Zeile unter *Grüne Mutanten* in §7), kein Test erwartet |
| V-196 | INFO | Die am Server gemessenen Abstände können die aufgezeichneten um wenige Millisekunden unterschreiten (97 statt 100 ms, 953 statt 954 ms), weil `clock_timestamp()` die Ausführung misst, nicht das Senden. Die Zusage gilt für das Senden durch `play` (gemessen mit der Uhr des Services), die Toleranz von 50 ms im E2E-Test deckt das. Keine Verletzung; ein Test mit Toleranz 0 am Server würde flackern. Bestätigt F-657. | `test/integration/zeitangaben_einspielen_e2e_test.go` · `const toleranz = 50.0` | ohne Blockade |
| V-197 | INFO | In der Spezifikation ist der Umbruch der geänderten Stelle *Prüfung der Werte* ungleichmäßig (der Satz „Die Prüfung sieht jede Session“ steht allein auf einer kurzen Zeile). Inhalt und Sätze sind vollständig. | `spec/spezifikation.md` · `Die Prüfung sieht jede Session` | Architect, kosmetisch |

Falsifikationsversuche, die nichts fanden: Ordnung der Prüfungen mit beschädigter Interaktion und fehlendem Wert (Binary); Session ohne Interaktion mit `--keep-timing`; zwei Signale in laufender Anfrage ohne Zeitangaben; `--keep-timing=false` gegen Variable und Datei; Variable fremder Kommandos; Überlauf des Zeitpunkts (`TestPlayZeitGrosserWert`, Mutant rot); Wartezeit nach dem ersten Signal mit und ohne `--finish-session-on-interrupt` am Binary; `replay` mit einer Stunde Pause in den Zeitangaben; Aufzeichnung in `sqlite` und `yaml` mit denselben Ergebnissen.

---

## Verdikt

**Befunde:** 0 HIGH, 0 MEDIUM, 1 LOW (V-193), 4 INFO (V-194 bis V-197).

Die drei Liefer-Punkte der DoD sind **bestätigt**: Die Tests existieren (36 Namen, je eine Fundstelle), werden ohne den Fix aus dem richtigen Grund rot (52 neu gebaute Mutanten im Wortlaut von [`LH-FA-21.a`](../../spec/spezifikation.md#lh-fa-21a--zeitangaben): 28 im Service, 19 im CLI-Adapter, davon vier im Hilfetext, 1 in der Verdrahtung und 4 am Binary, davon 2 rot und 2 grün aus dem in V-195 genannten Grund; alle anderen rot, keiner aus falschem Grund), und das gebaute Binary hält die Pausen in beiden Formaten und allen Betriebsarten ein, wartet nie vor der ersten Interaktion von `relative` und `first-request`, holt bei `absolute` die Verspätung auf und bei `relative` nicht, und beendet die Wartezeit mit dem ersten Signal sofort. Die Entscheidung des Architect zur Ordnung von `PGR-E3003` und `PGR-E2003` steht in der Spezifikation, im Plan und im Code und ist am Binary bestätigt. Der Eingriff `<-geschlossen` verursacht im Play-Pfad ohne Zeitangaben keinen Deadlock.

**Closure-blockierend:** nichts im Code. Vor der Closure: V-193 (§3 nachziehen), die Stand-Angabe des Gate-Absatzes (V-194) und die Closure-Pflichten aus §2 des Plans (Notiz mit Lerneintrag, Register, Ausgänge der zwei Risiken aus §6: das erste, Schwankung der Zeitabstände, ist am Binary nur als Untergrenze geprüft und im Mess-Verhalten bestätigt (V-196); das zweite hängt an `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients`).

**Finding-Klasse dieses Laufs:** „Plan folgt dem Code nicht“ (V-193, `BEO-REPO/plan-folgt-korrektur-nicht`, nur der Teil Testdateien); die Klasse „Zusage weiter als der Test“ (V-175, V-182) trat nicht auf.

**Aufgeräumt:** Container (`docker rm -fv`), Netz und die zwei Images dieser Sitzung per Name entfernt; kein `prune`; nur benannte Verzeichnisse im Scratchpad (`w/`, `gocache/`, `e2e/daten/`). Der Arbeitsbaum des Repos war bis zu diesem Bericht unverändert.

**Gate:** `make gates` auf dem sauberen Baum nach dem Commit dieses Berichts: Ergebnis in der Übergabe-Nachricht (ein Gate-Lauf belegt nur den Stand, auf dem er lief).
