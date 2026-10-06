# Verifikation: slice-replay-semantik-meldungscodes — 2026-10-06

**Rolle:** Verifier (Modul 11). Ich prüfe, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules prüft der Reviewer.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-replay-semantik-meldungscodes.md` (Kopf, §1 bis §6, §8) gegen den Gesamt-Diff `852c6d0..7dc26b7`. Spezifikation: `ac12870`, `f2c428a`, `ba2de96`. Code: `1348a9a`, `a9f4b15`, `7dc26b7`. Review-Report: `90c9079`.

**Eingang:**

- die DoD-Liefer-Punkte und die Commit-Messages
- der [Review-Report](2026-10-06-review-slice-replay-semantik-meldungscodes.md) (F-399 bis F-405, Stand `1348a9a`). Die Nacharbeit `a9f4b15`, `ba2de96` und `7dc26b7` hat kein eigenes Review gesehen. Ich habe sie mit eigenen Mutationen geprüft (Abschnitt 2).
- [`LH-FA-01`](../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-10`](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-11`](../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers), [`LH-FA-13`](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-14`](../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-FA-18`](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) mit den Abnahmeszenarien 4 und 6. In der Spezifikation: [`LH-FA-01.a`](../../spec/spezifikation.md#lh-fa-01a--kommandos-und-hilfe), [`LH-FA-10.a`](../../spec/spezifikation.md#lh-fa-10a--mismatch), [`LH-FA-13.b`](../../spec/spezifikation.md#lh-fa-13b--fehlerebenen-und-prozessstatus), [`LH-FA-14.a`](../../spec/spezifikation.md#lh-fa-14a--logging-und-diagnose), [`LH-FA-17.a`](../../spec/spezifikation.md#lh-fa-17a--konfiguration), [`LH-FA-18.a`](../../spec/spezifikation.md#lh-fa-18a--extended-query-ablauf-aufzeichnung-matching) §Mismatch, [`SPEC-033`](../../spec/spezifikation.md#spec-033--sicherheit), [`SPEC-034`](../../spec/spezifikation.md#spec-034--meldungscodes) §Ausgabe
- [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md)
- der Welle-Plan `docs/plan/planning/welle-replay-semantik.md` §3, `.claude/commands/close-welle.md` Schritt 1 und die Folge-Slices `docs/plan/planning/open/slice-v1-abschluss-betrieb.md` und `docs/plan/planning/open/slice-v1-abschluss-einspielen.md`

Bei meinem Start war der Arbeitsbaum sauber, HEAD `7dc26b7`. Dieser Commit enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

Die Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg.

**Wie die Mutanten gebaut wurden.** Alles lief unter `verify-mc/` im Scratchpad. Für jeden Mutanten habe ich eine frische Kopie aus `git archive 7dc26b7` angelegt. Die Ersetzung erfolgte über ein Skript, das abbricht, wenn die Stelle nicht genau einmal vorkommt. Den Docker-Build-Kontext habe ich umgangen (implement-slice Schritt 19, Bind-Mount). Ein Image der Stufe `deps` des `Dockerfile` (`vmc-deps:local`) trägt nur die Module und keinen Quellstand.

- **Unit-Tests:** Die Kopie ist eingebunden, darin läuft netzlos `go test -count=1`.
- **E2E:** Binary und Integrationstests entstehen im Container aus der schreibgeschützt eingebundenen Kopie. Sie laufen gegen `postgres:17-alpine` mit dem Digest aus `harness/mk/integration.mk`, in einem eigenen internen Docker-Netz, mit `-test.v`.

Ohne Mutation ist die Kopie auf beiden Wegen grün, dazu `gofmt -l` leer und `go vet -tags integration ./...` ohne Befund. `TestE2EReplayAbweichung`, `TestE2EReplayExtendedAbweichung`, `TestE2EReplayNichtVerbrauchtRangfolge` und `TestE2EReplaySelect1` liefen dreimal hintereinander grün (`-test.count=3`).

Eigene Container, das Netz, das Cache-Volume, das Image `vmc-deps:local` und die Kopien habe ich danach gelöscht. Die Images, die `make gates` unter den Namen des Repos baut (`pgwire-recorder:dev`, `:test`, `:integration`), sind nicht meine und bleiben stehen.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — [LH-FA-14](../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Kopf, Klasse und Exit-Code je Fehlerklasse. **Bestätigt**, mit zwei Befunden an Regeln aus `ba2de96` (V-56, V-57).

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| je Klasse 1 bis 6 Kopf `<klasse> [<code>]: ` mit Klassenname nach [`SPEC-034`](../../spec/spezifikation.md#spec-034--meldungscodes) §Ausgabe | `TestFehlerKopfJeKlasse` je Klasse; VE3 (Name der Klasse 1) rot | bestätigt |
| … als Zeile beim Prozessende, Exit-Code der Klasse | `TestFailJeKlasse` je Klasse 1 bis 6; VE4 (`PGR-E1000` mit Exit-Code 0) und VB3 (Rohtext statt Fehlertext) rot | bestätigt |
| … als Attribut `error` und als Meldungstext der `ErrorResponse`, SQLSTATE nach Klasse (`0A000`, `08006`, sonst `XX000`), `FATAL`, keine weiteren Felder | `TestFehlerantwortJeKlasse` vergleicht die ganze Nachricht (`DeepEqual`) mit dem Text des Attributs `error`; VN2 und VW2 rot | bestätigt |
| Fehlertext ist eine Zeile (*Eine Zeile*): LF, CR LF, CR mit Leerraum danach, VT, FF, U+0085, U+2028, U+2029 bleiben | `TestFehlerEineZeile`; VU1 bis VU4 rot | bestätigt |
| eine Kette trägt einen Kopf mit dem Code des äußersten Fehlers (*Fehlerkette*) | `TestFehlerKette`; VK1 (innerer Kopf), VK2 (Kopf durch eine fremde Hülle), VK3 (Code des innersten) rot | bestätigt |
| … Hülle mit eigenem Text um mehrere Ursachen ist eine Kette (`ba2de96`, F-402) | `TestFehlerKette` S1 bis S3 und „Hülle innen“; VG1 (jede `Unwrap() []error` gilt als gleichrangig) rot | bestätigt. Welcher Code bei Ursachen in verschiedener Tiefe gilt, prüft kein Test (V-57). |
| gleichrangige Fehler je eine Meldung, erste gemerkt, nicht klassifizierte als Ursache der ersten (*Gleichrangige Fehler*) | `TestFehlerGleichrangig`, `TestNoteGleichrangig`, `TestFailJeKlasse` (Join); VG2, VG3, VG4, VN1, VB1, VB2 rot | bestätigt |
| nicht eingeordneter Fehler ist `PGR-E1000` mit Kopf (*Fremde Fehler*) | `TestFehlerFremd`, `TestFehlerantwortJeKlasse`, `TestFailJeKlasse`; VE1 rot | bestätigt |
| `ErrorResponse` bei mehreren Meldungen: genau eine, die der ersten (*Zustellung an den Client*, `ba2de96`, F-401) | `TestFailGleichrangig`. VF1 (letzte statt erste) und VF3 (Text der letzten) rot. **VF2 (je Meldung eine `ErrorResponse`) grün.** | „die erste“ bestätigt, „genau eine“ nicht (V-56) |
| nicht annehmbare Verbindung ist der gemerkte Verbindungsfehler `PGR-E4000` ([`LH-FA-13.b`](../../spec/spezifikation.md#lh-fa-13b--fehlerebenen-und-prozessstatus)), Serve nimmt weiter an | `TestAnnahmefehler`; VA1, VA2, VA3 rot | bestätigt |
| jeder Code der Tabelle im Quelltext hat die Form aus `SPEC-034` und ergibt seine Klasse | `TestCodeTabelle` liest die Konstanten aus dem Quelltext | bestätigt (Reviewer-Mutation M18 ist dieselbe Klasse wie VE3; eigene Mutation an der Tabelle nicht gefahren) |
| Abbildung beim Herunterfahren | laut DoD Gegenstand von `slice-v1-abschluss-betrieb` | nicht geprüft, ausgeschlossen |

### Punkt 2 — [LH-FA-14](../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [LH-FA-17](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): `--log-level`. **Bestätigt mit Einschränkung** (V-58).

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Option und `PGWIRE_RECORDER_LOG_LEVEL` bei `record` und `replay` | `TestParseLogLevel`, `TestParseLogLevelUmgebung`; VL6 (`replay` ohne Option) rot | bestätigt |
| vier Stufen, jede zeigt ihre und die strengeren, Standard `info` | `TestLoggerSchwelle`, `TestRunLogLevel`, `TestRunRecordLogLevel`; VS1, VS2 rot | bestätigt |
| jeder andere Wert ist `PGR-E2001`, auch leer, großgeschrieben, anderer Name | `TestParseLogLevelWerte`; VL3, VL4, VL7 rot | bestätigt |
| Reihenfolge: Standard, Umgebung, Kommandozeile, letzte Angabe gilt ([`LH-FA-17.a`](../../spec/spezifikation.md#lh-fa-17a--konfiguration)) | VL1 (Umgebung geht vor) und VL5 (erste Angabe gilt) rot | bestätigt |
| ungültige Umgebungsvariable ist `PGR-E2001`, auch neben gültiger Option; leere gilt als nicht gesetzt | `TestParseLogLevelUmgebungUngueltig`; VL2 rot | bestätigt |
| Log-Zeile trägt `level`, Fehler `code` und `error`, Warnungen `code` | VW1 (Warnung am Prozessende ohne `code`), VW2 (`error` umbenannt), VW3 (Warnung im Adapter ohne `code`) rot | bestätigt |
| Startfehler auf jeder Stufe als Fehlertext, vorher nichts anderes | `TestRunStartfehlerJeStufe` für `replay`; VS3 (Startfehler als Log-Zeile) und VS4 (Zeile vor dem Parsen) rot. **Für `record` grün:** VS5 (Startfehler bei `record` auf `error` unterdrückt), VS6 (Schreibfehler am Ende als Log-Zeile außer auf `info`), auch in allen `TestE2ERecord*` | für `replay` bestätigt, für `record` nicht durch einen Test (V-58) |
| auf `info` keine Zeile je Verbindung oder Interaktion (F-399) | `TestInfoOhneZeileJeVerbindung`; VI1 (Zeile je Verbindung) rot bei `record`, VI2 (Zeile je Anfrage) rot bei `replay` | bestätigt |
| Debug-Zeile ohne Startnachricht mit `grund` statt `error` | `TestOhneStartnachrichtGrund`; VD1 rot | bestätigt |
| Hilfe geht Wert und Umgebungsvariable vor ([`LH-FA-01.a`](../../spec/spezifikation.md#lh-fa-01a--kommandos-und-hilfe)) | `TestParseLogLevelHilfe`; VH1 (Umgebung vor der Hilfe geprüft) rot | bestätigt |
| `version` liest keine Umgebungsvariable (`LH-FA-01.a`, §6 „wirkt auf keines davon“) | Code: `Parse` ruft `logLevelOption` für `version` nicht. VV1 (`version` prüft die Variable) grün | am Code gelesen, kein Test (V-59) |

### Punkt 3 — [LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Parameter-Nummer, keine Werte, V-52. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Nummer des ersten abweichenden Parameters, ab 1 ([`LH-FA-18.a`](../../spec/spezifikation.md#lh-fa-18a--extended-query-ablauf-aufzeichnung-matching) §Mismatch) | `TestReplayExtendedParameterStelle`, `TestReplayExtendedAbweichung`; VP1 (ab 0) und VP4 (letzter statt erster) rot | bestätigt |
| bei abweichender Zahl beide Anzahlen, erwartet vor empfangen | VP2 (vertauscht) und VP5 (keine Anzahlen) rot | bestätigt |
| kein Wert in der Diagnose ([`SPEC-033`](../../spec/spezifikation.md#spec-033--sicherheit)) | VP3 (Wert neben der Nummer) rot | bestätigt |
| kein Wert im Log, auch bei `--log-level debug`, und in der `ErrorResponse` (E2E) | `TestE2EReplayExtendedAbweichung` mit `--log-level debug`; EP1 (Adapter schreibt jeden Parameterwert als Debug-Zeile) rot. Die einzige Debug-Zeile im Code ist die ohne Startnachricht. | bestätigt |
| V-52, Unit: `Query` liefert bei Abweichung keine Antworten, auch nach dem Ende der Aufzeichnung | `TestReplayMismatch`; VQ1 und VQ2 rot | bestätigt |
| V-52, E2E: `TestE2EReplayAbweichung` verlangt kein Ergebnis von `ReadAll` | EV1 (aufgezeichnete Antworten vor `PGR-E5001`) rot mit „1 Ergebnisse, nil“. EV2 (ohne `ReadyForQuery`) rot mit „1 Ergebnisse“ neben `PGR-E5001`. Gegenprobe: EV2 gegen den Test aus `852c6d0` grün, das ist die Lücke aus V-52. | bestätigt |
| Abdeckung `LH-FA-10/Negative` (E2E und Unit) | `docs/user/abdeckung-e2e.md`, `docs/user/abdeckung-unit.md`; `make abdeckung-check` grün | bestätigt |

### Punkt 4 — `make gates` grün. **Bestätigt.**

`make gates` habe ich selbst im Repo auf `7dc26b7` gefahren, Arbeitsbaum sauber. **Exit-Code 0.**

- `a-check`: 0 Befunde; `a-check-negativ`, `abdeckung-gegenprobe`, `commit-msg-gegenprobe`, `kopf-check-gegenprobe`: grün
- `baseline-verify`: v6.13.0 OK, 54 Dateien
- d-check: 235 Dateien, 0 Befunde
- Stufe `test` aus dem Build-Cache. Den Unit-Lauf habe ich deshalb zusätzlich in der Kopie mit `-count=1` gefahren, dazu `gofmt -l` und `go vet`: alle Pakete `ok`, ohne Befund.
- `run-integration-tests`: grün, beide Phasen, 46 E2E-Tests `PASS`, keiner `FAIL`
- `make kopf-check` und `make abdeckung-check` einzeln: Exit-Code 0

### Prozess-Punkte der DoD (nur vermerkt, kein Häkchen gesetzt)

Der Review-Report liegt vor. Closure-Notiz, Register, Risiko-Ausgänge und Paarungen stehen aus. §7 ist leer, und vor der Closure ist das richtig. Die drei Risiken in §6 tragen „offen bis Closure“. Für das erste gilt V-57: Der Ausgang „entfällt, wenn der Kettentest aus DoD 1 grün ist“ trägt nur für Ketten mit einer Ursache.

---

## 2. Bewusstes Brechen

Je Mutant eine frische Kopie. 64 Mutationen, davon 56 rot. Grün blieben VE2 (äquivalent), VF2, VK4, VS5, VS6, VV1 sowie VB4 und EV3, die jeweils der andere Weg rot färbt (EB4 und VQ2).

**Fehlertext und Meldungen** (`internal/hexagon/model/fehler.go`)

| # | Mutation | Rot in | Grund (Testmeldung) |
|---|---|---|---|
| VK1 | **ein Kopf je Kette:** innerer klassifizierter Fehler mit `Error()` statt `ursache` | `TestFehlerKette`, `TestFehlerGleichrangig` | „Recording [PGR-E3001]: Session 1, Interaktion 2: Recording [PGR-E3003]: …“ |
| VK2 | Hülle ersetzt den inneren Text nicht durch seine Ursache | `TestFehlerKette` | „Kette durch einen nicht klassifizierten Fehler“: innerer Kopf `[PGR-E4002]` steht |
| VK3 | Kopf mit dem Code des innersten Fehlers | `TestFehlerKette` | `PGR-E3003` statt `PGR-E3001` |
| VK4 | erster klassifizierter Fehler in Breitensuche statt Tiefensuche | **grün** | V-57 |
| VG1 | **gleichrangig gegen Hülle:** jede `Unwrap() []error` ist eine Zusammenfassung (F-402 zurückgedreht) | `TestFehlerKette` S1, S1 fremd zuerst, S1 ohne Klasse, S2, Hülle innen | zwei Meldungen statt einer, „Kontext“ fehlt |
| VG2 | keine Zusammenfassung erkannt | `TestFehlerKette` S3, `TestFehlerGleichrangig` | eine Meldung statt drei, Trenner `" "` statt `"; "` |
| VG3 | nicht klassifizierter wird eigene Meldung `PGR-E1000` | `TestFehlerGleichrangig` | drei Meldungen statt zwei |
| VG4 | Reihenfolge der Meldungen umgekehrt | `TestFehlerGleichrangig` | erste Meldung `PGR-E3001` statt `PGR-E4003` |
| VE1 | **`PGR-E1000`** ohne Kopf | `TestFehlerKette`, `TestFehlerFremd`, `TestFehlerEineZeile`, `TestFehlerGleichrangig` | „connection reset by peer“ ohne „sonstiger Fehler [PGR-E1000]: “ |
| VE2 | Rückfall-Exit-Code 0 für einen Code ohne Ziffer 1 bis 6 | grün | äquivalent: `TestCodeTabelle` schließt solche Codes aus |
| VE3 | Klassenname der Klasse 1 `intern` | `TestFehlerKopfJeKlasse` u. a. | „intern [PGR-E1000]“ |
| VE4 | `PGR-E1000` mit Exit-Code 0 | `TestFehlerKopfJeKlasse`, `TestFailJeKlasse` | „Exit-Code 0, erwartet 1“ |
| VU1 | **Umbrüche:** durch nichts statt ein Leerzeichen ersetzt | `TestFehlerEineZeile` | „x: ab“ |
| VU2 | Leerraum nach dem Umbruch bleibt | `TestFehlerEineZeile` | `"a  \t    b"` |
| VU3 | FF als Umbruch | `TestFehlerEineZeile` | `\f` ersetzt |
| VU4 | `Error()` ohne `einzeilig` | `TestFehlerEineZeile` | „Umbruch in der eigenen Meldung“ |

**PGWire-Adapter** (`internal/adapters/driving/pgwire/server.go`)

| # | Mutation | Rot in | Grund (Testmeldung) |
|---|---|---|---|
| VF1 | **`fail` stellt die letzte Meldung zu** | `TestFailGleichrangig` | `08006`, „Netzwerk [PGR-E4003]: weg“ statt `0A000` und `PGR-E6001` |
| VF2 | `fail` stellt je Meldung eine `ErrorResponse` zu | **grün** | V-56 |
| VF3 | SQLSTATE der ersten, Text der letzten Meldung | `TestFailGleichrangig` | Text `PGR-E4003` |
| VA1 | **Annahmefehler `PGR-E4000`** nur geloggt, nicht gemerkt | `TestAnnahmefehler` | „gemerkt \"\"“ |
| VA2 | Annahmefehler als `PGR-E1000` | `TestAnnahmefehler` | `code=PGR-E1000` |
| VA3 | Annahmefehler beendet `Serve` | `TestAnnahmefehler` | „Serve nach dem Fehler nicht weiter“ |
| VN1 | `note` protokolliert nur die erste Meldung | `TestNoteGleichrangig`, `TestFailGleichrangig` | eine Log-Zeile statt zwei |
| VN2 | Attribut `error` mit `err.Error()` | `TestFehlerantwortJeKlasse`, `TestNoteGleichrangig` | `"nicht\neingeordnet"`, Join mit Umbrüchen |
| VW2 | Schlüssel `error` umbenannt | `TestFehlerantwortJeKlasse`, `TestNoteGleichrangig`, `TestAnnahmefehler` | „Fehlertext \"\"“ |
| VW3 | Warnung im Adapter ohne `code` | `TestReplayModus` | „Warnung PGR-W2001 fehlt“ |
| VD1 | Debug-Zeile ohne Startnachricht unter `error` | `TestOhneStartnachrichtGrund` | `error=EOF` an `level=DEBUG` |
| VI1 | **keine Zeile je Verbindung auf `info`:** Info-Zeile nach dem Öffnen der Session | `TestInfoOhneZeileJeVerbindung` | „record: Zeile der Stufe INFO je Verbindung oder Anfrage“ |
| VI2 | Info-Zeile je beantworteter einfacher Anfrage im Replay | `TestInfoOhneZeileJeVerbindung` | „replay: Zeile der Stufe INFO …“ |

**Bootstrap** (`internal/bootstrap/bootstrap.go`)

| # | Mutation | Rot in | Grund (Testmeldung) |
|---|---|---|---|
| VB1 | Zeile beim Prozessende nur für die erste Meldung | `TestFailJeKlasse` | zweite Zeile `PGR-E4000` fehlt |
| VB2 | Exit-Code der letzten Meldung | `TestFailJeKlasse` | „Exit-Code 4 … erwartet 3“ |
| VB3 | Zeile beim Prozessende als `err.Error()` | `TestFailJeKlasse` | „intern\nkaputt“ |
| VB4 | Meldungen der nie zugeordneten Sessions nur ohne gemerkten Code geloggt | Unit grün; **E2E rot** (EB4) | — |
| VS1 | **Strenge:** `debug` als `info` | `TestLoggerSchwelle` | „debug: [INFO WARN ERROR]“ |
| VS2 | `error` zeigt auch `warn` | `TestRunLogLevel`, `TestLoggerSchwelle` | Zeile `WARN PGR-W2001` auf `error` |
| VS3 | **Startfehler** bei `replay` als Log-Zeile statt Fehlertext | `TestRunStartfehlerJeStufe` | `time=… level=ERROR …` auf `error`, `warn`, `debug` |
| VS4 | Log-Zeile vor dem Parsen | `TestRunLogLevel`, `TestRunStartfehlerJeStufe` | `level=INFO msg=start` vor dem Fehlertext |
| VS5 | Startfehler bei `record` (Listen-Port) auf `error` unterdrückt | **grün** (Unit und alle `TestE2ERecord*`) | V-58 |
| VS6 | Schreibfehler am Ende bei `record` außer auf `info` als Log-Zeile | **grün** (Unit und alle `TestE2ERecord*`) | V-58 |
| VW1 | Warnung am Ende von `replay` ohne `code` | `TestRunLogLevel` | „[WARN] erwartet [WARN PGR-W2001]“ |

**CLI** (`internal/adapters/driving/cli/cli.go`)

| # | Mutation | Rot in | Grund (Testmeldung) |
|---|---|---|---|
| VL1 | **Reihenfolge:** Umgebungsvariable nach `fs.Parse`, geht der Kommandozeile vor | `TestParseLogLevelUmgebung` | Stufe aus der Umgebung statt der Option |
| VL2 | **Strenge:** ungültige Umgebungsvariable ignoriert | `TestParseLogLevelUmgebungUngueltig` | kein `PGR-E2001` |
| VL3 | Groß- und Kleinschreibung egal | `TestParseLogLevelWerte`, `…UmgebungUngueltig` | `INFO`, `Info` angenommen |
| VL4 | leerer Wert gilt als Standard | `TestParseLogLevelWerte` | `--log-level=` angenommen |
| VL5 | erste statt letzter Angabe gilt | `TestParseLogLevel` | „\"debug\" … erwartet \"error\"“ |
| VL6 | `replay` ohne `--log-level` | `TestParseLogLevel`, `…Umgebung`, `…UmgebungUngueltig` | „flag provided but not defined: -log-level“ |
| VL7 | Alias `warning` | `TestParseLogLevelWerte`, `…UmgebungUngueltig` | `warning` angenommen |
| VH1 | Umgebungsvariable vor der Hilfe geprüft | `TestParseLogLevelHilfe` | `PGR-E2001` statt Hilfe |
| VV1 | `version` prüft `PGWIRE_RECORDER_LOG_LEVEL` | **grün** | V-59 |

**Use Case** (`internal/hexagon/services/matcher.go`, `replay.go`)

| # | Mutation | Rot in | Grund (Testmeldung) |
|---|---|---|---|
| VP1 | **Parameter-Nummer** ab 0 | `TestReplayExtendedAbweichung`, `…ParameterStelle` | `(Parameter $0)` |
| VP2 | Anzahlen vertauscht | `TestReplayExtendedParameterStelle` | „Anzahl erwartet 1, empfangen 3“ |
| VP3 | **nie Werte:** Wert neben der Nummer | beide Tests | `(Parameter $1 = "geheim-wert")` |
| VP4 | letzter statt erster abweichender Parameter | `TestReplayExtendedParameterStelle` | `$3` statt `$2` |
| VP5 | keine Anzahlen bei abweichender Zahl | `TestReplayExtendedParameterStelle` | `""` |
| VQ1 | **V-52, Unit:** `Query` liefert bei abweichender Anfrage die aufgezeichneten Antworten | `TestReplayMismatch` | „erwartet PGR-E5001 ohne Antworten, erhalten [command_complete A, ready_for_query]“ |
| VQ2 | … nach dem Ende der Aufzeichnung die der letzten Interaktion | `TestReplayMismatch` | „Anfrage über die Aufzeichnung hinaus: [command_complete A …]“ |

**E2E** (`test/integration`)

| # | Mutation | Rot in | Grund (Testmeldung) |
|---|---|---|---|
| EV1 | **V-52, E2E:** Use Case liefert die aufgezeichneten Antworten, der Adapter sendet sie vor `PGR-E5001` | `TestE2EReplayAbweichung` | „1 Ergebnisse, nil“ |
| EV2 | wie EV1, ohne `ReadyForQuery` | `TestE2EReplayAbweichung` | „1 Ergebnisse“ neben „Replay [PGR-E5001]“; gegen den Test aus `852c6d0` grün |
| EV3 | wie EV2 für eine Anfrage nach dem Ende der Aufzeichnung | **grün** in allen `TestE2EReplay*`, `TestE2EFehlerreplay*`, `TestE2EErgebnisarten*` (19 Tests); Unit rot (VQ2) | V-60 |
| EP1 | **Parameterwerte bei `debug`:** Adapter schreibt jeden Wert als Debug-Zeile | `TestE2EReplayExtendedAbweichung` | „Diagnose nennt den Parameterwert“ (die Nummer bleibt unverändert, rot macht der Wert) |
| EB4 | wie VB4 | `TestE2EReplayNichtVerbrauchtRangfolge` | „PGR-E5002 fehlt oder steht nicht in Reihe“ |
| ES1 | **Szenario 6:** SQLSTATE jeder Fehlerantwort `XX000` (Record und Replay) | `TestE2EFehlerreplayEinfach` | „Sicht beim Aufzeichnen“ |
| ES2 | Meldung jeder Fehlerantwort mit angehängtem `[PGR-E5000]` (Record und Replay) | `TestE2EFehlerreplayEinfach` | „Sicht beim Aufzeichnen“ |
| ES3 | dasselbe nur im Replay (Use Case) | `TestE2EFehlerreplayEinfach` | „Sicht im Replay“ |

**Sonden** (Kopie, nur gelesen):

- `TestFailGleichrangig` unter VF2: Nach dem ersten `Receive` ist `netz.Len()` 0, ein zweites `Receive` liefert die zweite `ErrorResponse` (`08006`, „Netzwerk [PGR-E4003]: weg“). Ohne Mutation liefert es `unexpected EOF`. Grund und Folge siehe V-56.
- `Meldungen(fmt.Errorf("K %w / %w", fmt.Errorf("i %w", a), b))` mit `a` = `PGR-E6001`, `b` = `PGR-E3001` ergibt eine Meldung mit `PGR-E6001` (Tiefensuche). Siehe V-57.
- `Meldungen(fmt.Errorf("%v", a))` ergibt „sonstiger Fehler [PGR-E1000]: nicht unterstützt [PGR-E6001]: a“, also zwei Köpfe. Kein Code im Repo bettet einen Fehler über `%v` oder `%s` ein (`grep`). Siehe Grenzen.

---

## 3. Closure-Trigger der Welle: Abnahmeszenario 4 und 6 für Simple Query

Der Trigger in `welle-replay-semantik` §3 verlangt beide Szenarien „für Simple Query end-to-end automatisiert nachgewiesen“. close-welle Schritt 1 verlangt je Zusage einen Test, der unter ihrer Mutation rot wird, auch für die verneinende.

- **Szenario 4** („erkennt die Abweichung und liefert … einen eindeutigen Fehler, statt eine unpassende aufgezeichnete Antwort zu verwenden“, [LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage)) steckt in `TestE2EReplayAbweichung` (`pgconn.Exec`, also Simple Query).
  - *Bejahende Zusage:* `PGR-E5001` beim Client und Exit-Code 5 prüft der Test. Die Felder der Diagnose nach [LH-FA-10.a](../../spec/spezifikation.md#lh-fa-10a--mismatch) prüft `TestReplayMismatch`.
  - *Verneinende Zusage:* Für die abweichende Anfrage sind EV1 und EV2 E2E rot und VQ1 im Unit-Test. Die Gegenprobe mit dem Test aus `852c6d0` bleibt grün. Damit ist V-52 geschlossen.
  - *Randform:* Eine Anfrage nach dem Ende der Aufzeichnung ist nach [LH-FA-10.a](../../spec/spezifikation.md#lh-fa-10a--mismatch) ebenfalls ein Mismatch. Für sie färbt nur der Unit-Test die Mutation rot (VQ2), E2E bleibt sie grün (EV3, V-60).
- **Szenario 6** ([LH-FA-11](../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers)) steckt in `TestE2EFehlerreplayEinfach`. Dieser Slice ändert den Test nicht. Auf `7dc26b7` sind ES1, ES2 und ES3 rot, ES3 bei „Sicht im Replay“. Die neue Zustellung des Recorders (`fail`) berührt die wiedergegebene `ErrorResponse` des Servers nicht. ES2 zeigt das: Ein angehängter Code würde rot. Das ist die Zusage aus `SPEC-034` *Zustellung an den Client*, dass eine wiedergegebene `ErrorResponse` keinen Meldungscode trägt.

**Urteil:** Den Abnahme-Teil des Closure-Triggers sehe ich erfüllt: Szenario 4 und Szenario 6 sind für Simple Query E2E nachgewiesen, auch die verneinende Zusage von Szenario 4 hat eine rote Mutation. Das Szenario beschreibt die abweichende Anfrage, und die ist E2E geprüft. Die Randform „nach dem Ende der Aufzeichnung“ trägt der Unit-Test. Wer das für den Trigger E2E verlangt, muss V-60 vor der Welle-Closure schließen. Offen bleiben die übrigen Bedingungen aus §3: Alle Slices liegen in `done/` (dieser noch nicht), `make gates` beim Schließen, Closure-Notiz der Welle.

---

## 4. Plan gegen Code (`AGENTS.md` §3.9)

| Stelle | Ergebnis |
|---|---|
| Kopf `Bezug` | nennt [LH-FA-14](../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [LH-FA-17](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [LH-QA-05](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md). Stimmig. |
| Kopf `Berührte Spec-Stellen` | deckt jede Stelle, die der Diff in der Spezifikation ändert: `LH-FA-13.b`, `LH-FA-14.a`, `LH-FA-18.a`, `SPEC-033`, `SPEC-034`. Nennt die Sicht-Stellen `ARC-001`, `ARC-005`, `ARC-006` und `ARC-009` (F-403 behoben). `make kopf-check` grün. |
| §1 | Ziel, „Schon geliefert“, „Neu“, die zwei Übernahmen und die Ausschlüsse stimmen mit dem Diff überein. Kein Produktionscode außerhalb von §3 im Diff. |
| §2 | Die drei Liefer-Punkte entsprechen dem gelieferten Umfang. Die Regeln aus `ba2de96` (eine `ErrorResponse`, Hülle um mehrere Ursachen) nennt die DoD nicht eigens. Sie stehen in §3 und §6, und §3.9 verlangt die DoD nicht. |
| §3 | Jede Zeile hat ihren Gegenstand im Diff. Die Zeile zu `fehler.go` nennt die Erkennung der reinen Zusammenfassung am Text (`7dc26b7`), die zu `fehler_test.go` S1 bis S3. Die Zeile zu `server.go` nennt „`fail` stellt genau eine `ErrorResponse` … zu“ und die Reihenfolge im Join aus `CloseSession` (F-405 behoben). `record.go` liegt wie zugesagt nicht im Diff. |
| §6 | Die Zeilen F-401 und F-402 sind ergänzt (`ba2de96`). Der Code folgt beiden, mit den Lücken aus V-56 und V-57. Risiko 1 hängt an V-57. |
| Folge-Slice `slice-v1-abschluss-betrieb` | §1 (Übernahme), DoD 3 (`log_level`, ungültig `PGR-E2004`) und die §3-Zeile für `internal/adapters/driving/cli` (allgemeiner Leser ersetzt beide Einzel-Leser, Schlüssel `log_level`) sind nachgezogen (F-400 behoben). Die Info-Zeile beim Beginn des Herunterfahrens steht dort in §1, und ihre Zuordnung zur DoD ist ausdrücklich dem Planner überlassen. Das stand schon vor diesem Slice so. |
| Folge-Slice `slice-v1-abschluss-einspielen` | §1, DoD 1 (Wertemenge, Strenge, Schwelle, Umgebungsvariable neben gültiger Option) und die §3-Zeile für die CLI sind nachgezogen (F-400 behoben). |

---

## 5. Status der Findings aus dem Review

| Finding | Status | Beleg |
|---|---|---|
| F-399 | behoben | `TestInfoOhneZeileJeVerbindung` für `record` (zwei einfache Anfragen) und `replay` (einfach und Extended); VI1, VI2 rot |
| F-400 | behoben | DoD und §3 beider Folge-Slices (Abschnitt 4) |
| F-401 | **teilweise behoben** | Entschieden in [`SPEC-034`](../../spec/spezifikation.md#spec-034--meldungscodes) *Zustellung an den Client* und §6 (`ba2de96`). Der Code stellt die erste Meldung zu, VF1 und VF3 sind rot. „Genau eine“ prüft `TestFailGleichrangig` nicht, VF2 bleibt grün (V-56). |
| F-402 | behoben | Entschieden in *Fehlerkette* und *Gleichrangige Fehler* (`ba2de96`), Erkennung am Text (`7dc26b7`), S1 bis S3; VG1 rot. Offen bleibt die Rangfolge bei Ursachen verschiedener Tiefe (V-57). |
| F-403 | behoben | Kopf nennt `ARC-001`, `ARC-005`, `ARC-006`, `ARC-009` |
| F-404 | behoben | Deklaration an `TestLoggerOrtszeit` entfernt, Abdeckung neu geschrieben, `make abdeckung-check` grün |
| F-405 | erledigt (INFO) | §3 nennt die Reihenfolge im Join aus `CloseSession` und warum sie heute ohne Folge ist |

---

## 6. Befunde

| ID | Befund | Ort | Empfehlung an den Planner |
|---|---|---|---|
| V-56 | **„Genau eine `ErrorResponse`“ ist nicht geprüft.** `TestFailGleichrangig` will mit `netz.Len() > 0` eine zweite Nachricht erkennen. Das `pgproto3.Frontend` liest den `bytes.Buffer` aber beim ersten `Receive` ganz in seinen Puffer, deshalb ist `netz.Len()` danach immer 0. Die Sonde unter VF2 zeigt: Ein zweites `Receive` liefert die zweite `ErrorResponse`, und der Test bleibt grün. Der Kommentar am Test und `fail` sagen „genau eine“ zu, `SPEC-034` *Zustellung an den Client* ebenso. Heute erreicht kein Pfad `fail` mit einem Join (Review F-401), die Wirkung ist also latent. F-401 ist damit nur zur Hälfte behoben. Klasse `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. | `internal/adapters/driving/pgwire/server_meldung_test.go:263`–`:294` (`TestFailGleichrangig`); `internal/adapters/driving/pgwire/server.go:688`–`:697` | Den Test eine zweite Nachricht lesen lassen und `io.ErrUnexpectedEOF` beziehungsweise `io.EOF` verlangen, statt `netz.Len()` zu prüfen. Danach VF2 als Gegenprobe. Am Code ist nichts zu ändern. |
| V-57 | **Kettencode bei Ursachen in verschiedener Tiefe ist ungeprüft, und der Wortlaut lässt zwei Lesarten zu.** `SPEC-034` *Fehlerkette* (`ba2de96`): „Code und Klasse sind die des ersten klassifizierten Fehlers unter ihm, in der Reihenfolge seiner Ursachen und von außen nach innen“. Der Code nimmt `errors.As`, also eine Tiefensuche: Bei `K %w / %w` mit `i %w a` und `b` gilt `a` (Sonde). Eine Breitensuche, die `b` als äußeren nimmt, bleibt in allen Tests grün (VK4). Sie wäre mit „äußerster klassifizierter Fehler“ aus demselben Absatz ebenso vereinbar. Kein Pfad im Code erzeugt heute mehrere `%w`. Risiko 1 in §6 nennt als Ausgang „entfällt, wenn der Kettentest aus DoD 1 grün ist“. Grün ist er, für diese Form aber kein Beleg. | `spec/spezifikation.md` `SPEC-034` §Ausgabe *Fehlerkette*; `internal/hexagon/model/fehler.go:165`–`:198` (`Meldungen`); `internal/hexagon/model/fehler_test.go:68`–`:93` | Vor der Closure vom Architect den Wortlaut auf eine Lesart festlegen lassen, etwa „Ursachen in ihrer Reihenfolge, jede von außen nach innen“ für die Tiefensuche des Codes. Dazu ein Fall in `TestFehlerKette`, den VK4 rot färbt. Der Ausgang von Risiko 1 hängt an diesem Fall. |
| V-58 | **Startfehler und Schreibfehler am Ende auf jeder Stufe sind für `record` nicht durch einen Test belegt.** DoD 2: „bei `record` und `replay` … ein Startfehler erscheint auf jeder Stufe als Fehlertext“. [LH-FA-14.a](../../spec/spezifikation.md#lh-fa-14a--logging-und-diagnose) sagt dasselbe für den Fehler beim Schreiben des Recordings am Ende zu. `TestRunStartfehlerJeStufe` ruft nur `replay`. VS5 (Startfehler bei `record` auf `error` unterdrückt) und VS6 (Schreibfehler am Ende außer auf `info` als Log-Zeile) bleiben in Unit und E2E grün. `TestE2ERecordSchreibfehlerAmEnde` läuft auf der Standardstufe. Am Code gelesen: `record` meldet beide über `fail`, und `fail` hängt nicht am Logger. Das Verhalten ist also richtig, nur der Nachweis fehlt für `record`. Die Mutanten sind gezielt stufenabhängig gebaut. | `internal/bootstrap/bootstrap_test.go:188`–`:218`; `internal/bootstrap/bootstrap.go:49`–`:76` | `TestRunStartfehlerJeStufe` um `record` erweitern, etwa mit vorhandenem `--output` ohne `--force` (`PGR-E2002`), dazu ein Fall für den Schreibfehler am Ende mit `--log-level=error`. Danach VS5 und VS6 als Gegenprobe. |
| V-59 | **`version` mit ungültiger `PGWIRE_RECORDER_LOG_LEVEL` ist ungeprüft.** [LH-FA-01.a](../../spec/spezifikation.md#lh-fa-01a--kommandos-und-hilfe): „`version` liest weder Umgebungsvariablen …“. §6 nennt für Hilfe und `version`: „`--log-level` wirkt auf keines davon“. `TestRunVersion` setzt nur `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED`. VV1 bleibt grün. Der Code ruft `logLevelOption` für `version` nicht, die Wirkung ist heute also keine. Der Reviewer hat das gesehen und nicht als Befund geführt. Ich führe es, weil der Testkommentar „Es liest keine Umgebungsvariable“ weiter reicht als die Prüfung. | `internal/bootstrap/bootstrap_test.go:32`–`:41` | In `TestRunVersion` auch `PGWIRE_RECORDER_LOG_LEVEL=INFO` setzen. Klein, kann mit V-58 zusammen gehen. |
| V-60 | **Szenario 4 nach dem Ende der Aufzeichnung ist nur im Unit-Test belegt.** Liefern Use Case und Adapter für eine Anfrage nach der letzten Interaktion deren Antworten vor `PGR-E5001`, bleiben alle E2E-Tests grün (EV3, 19 Tests). Rot wird nur `TestReplayMismatch` (VQ2). Den Closure-Trigger halte ich trotzdem für erfüllt (Abschnitt 3). | `test/integration/replay_e2e_test.go:99`–`:120`; `internal/hexagon/services/replay_test.go:99`–`:124` | Nur vermerkt. Wer den Trigger für diese Randform E2E verlangt, zieht vor der Welle-Closure einen zweiten Fall in `TestE2EReplayAbweichung` nach: `SELECT 1;` beantwortet, dann `SELECT 1;` erneut, kein Ergebnis. |
| V-61 | **Struktur-IDs in Commit-Messages.** `ac12870` nennt `SPEC-033` und `SPEC-034`, `f2c428a` und `ba2de96` nennen `SPEC-034`. `AGENTS.md` §5 Regel 1: „Struktur-IDs … gehören nicht in die Commit-Message.“ Der Review-Report führt die Commit-Messages `ac12870` bis `1348a9a` als „geprüft, ohne Befund … keine nennt `SPEC-*` oder `ARC-*`“. Für `ac12870` und `f2c428a` trifft das nicht zu. Dasselbe Muster steht in 15 der letzten 300 Commits. Kein Sensor hält die Regel, der Commit-Hook prüft nur, ob eine Kennung vorhanden ist. | Commit-Messages `ac12870`, `f2c428a`, `ba2de96`; `AGENTS.md` §5 | Kein Umschreiben der Historie. Bei der Closure entscheiden, ob die Regel einen Träger im Commit-Hook bekommt oder gelockert wird. Beides ist eine Harness-Entscheidung, nicht Gegenstand dieses Slice. |

**Grenzen dieser Verifikation:**

- `make gates` lief im Repo selbst, nicht in einem frischen Klon. Der Arbeitsbaum war vor und nach dem Lauf sauber (`git status --short` leer). Die Stufe `test` kam aus dem Build-Cache. Den Unit-Lauf habe ich deshalb in der Kopie mit `-count=1` wiederholt.
- Ein Fehler, den Code über `%v` oder `%s` in einen anderen Text einbettet, ergäbe zwei Köpfe (Sonde). `SPEC-034` *Fehlerkette* schließt das aus, ein Test oder Sensor nicht. Heute tut kein Code das (`grep` auf `Errorf`, `Sprintf`, `errors.New` mit `%v`/`%s` und Fehlerargument leer). Ich führe das nicht als Befund.
- Für `record` belegt `TestInfoOhneZeileJeVerbindung` nur einfache Anfragen, keine Extended-Interaktion. Die DoD verlangt das nicht je Protokoll.
- Den Race-Detector habe ich nicht gefahren. Die Stufen bauen mit `CGO_ENABLED=0`. Dass `firstCode` nur unter `s.mu` gelesen und geschrieben wird, habe ich am Code gelesen (`note`, `FirstErrorCode`).

---

## 7. Gesamturteil

Liefer-Punkt 1 und 3 sowie `make gates` sind **bestätigt**. Liefer-Punkt 2 ist **bestätigt mit Einschränkung**: Für `replay` hat jede Teilbehauptung einen Test, der unter ihrer Mutation aus dem richtigen Grund rot wird. „Startfehler auf jeder Stufe“ ist für `record` nur am Code gelesen (V-58). An den Regeln aus der Nacharbeit (`ba2de96`, `7dc26b7`) ist „die erste Meldung“ belegt, „genau eine `ErrorResponse`“ nicht (V-56). Die Rangfolge bei Ursachen verschiedener Tiefe ist weder eindeutig formuliert noch geprüft (V-57). F-399, F-400, F-402 bis F-405 sind behoben, F-401 zur Hälfte. Der Plan folgt dem Code, die Folge-Slices sind nachgezogen, `make kopf-check` ist grün. Den Abnahme-Teil des Closure-Triggers der Welle sehe ich erfüllt: Szenario 4, mit der verneinenden Zusage, und Szenario 6 sind für Simple Query E2E nachgewiesen. V-52 ist geschlossen.

**Summary:** 2/3 Liefer-Punkte bestätigt, 1 mit Einschränkung (DoD 2, V-58) · `make gates` grün auf `7dc26b7` · 56 von 64 Mutationen rot. Grün blieben VF2 (V-56), VK4 (V-57), VS5 und VS6 (V-58), VV1 (V-59), EV3 (V-60, Unit rot), VB4 (E2E rot) und VE2 (äquivalent) · Befunde V-56 bis V-58 vor der Closure, V-59 klein, V-60 und V-61 vermerkt.
