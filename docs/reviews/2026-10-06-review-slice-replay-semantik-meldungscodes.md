# Review-Report: slice-replay-semantik-meldungscodes — 2026-10-06

**Review-Art:** Code und Spezifikation, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff 852c6d0..1348a9a`, also die Commits `ac12870` (Randformen in der Spezifikation, Plan §1/§3/§6, Folge-Slices), `f2c428a` (Entscheidung der drei Rückgaben des Implementers, Spezifikation, Plan §6) und `1348a9a` (Code, Tests, Handbuch, Abdeckung, Plan, `slice-v1-abschluss-betrieb`) für `slice-replay-semantik-meldungscodes`. Schwerpunkte laut Auftrag: `model.Meldungen` und `Error()`, `note`/`fail` in Server und Bootstrap, `--log-level`, Parameter-Nummer, V-52, §3.10 bis §3.12, Test-Hygiene um `time.Local`, Plan und Folge-Slice.

**Skill:** `.harness/skills/reviewer.md` @ `1348a9a`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-replay-semantik-meldungscodes.md` (Kopf, §1 bis §6, §8), `docs/plan/planning/open/slice-v1-abschluss-betrieb.md` (Kopf, §1, §2, §3), `docs/plan/planning/open/slice-v1-abschluss-einspielen.md` (§1, §2, §3)
- `spec/lastenheft.md` [LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [LH-FA-14](../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [LH-FA-17](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol)
- `spec/spezifikation.md`: LH-FA-01.a, LH-FA-13.b, LH-FA-14.a (§Log-Level, §Zeilenform), LH-FA-17.a, LH-FA-18.a §Mismatch, `SPEC-033`, `SPEC-034` §Ausgabe (*Eine Zeile*, *Fehlerkette*, *Gleichrangige Fehler*, *Fremde Fehler*, *Zustellung an den Client*)
- `spec/architecture.md` `ARC-001`, `ARC-005`, `ARC-006`, `ARC-009`, §5 (Fehlermodelle, Observability)
- [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md), [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)
- `AGENTS.md` §3.3 bis §3.12; Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/plan-folgt-korrektur-nicht`, `BEO-REPO/spec-randform-erst-im-review-entschieden`, `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`
- `docs/user/benutzerhandbuch.md` (Log-Ausgaben, Fehlerbehebung, Extended-Abweichung), `docs/user/abdeckung-*.md`
- Vorherige Reports am Modul: [Review zu `slice-replay-semantik-mismatch`](2026-10-05-review-slice-replay-semantik-mismatch.md) (F-386, F-388), [Review zu `slice-replay-semantik-fehlerreplay`](2026-10-05-review-slice-replay-semantik-fehlerreplay.md) (F-395), [Verifikation zu `slice-replay-semantik-fehlerreplay`](2026-10-06-verifikation-slice-replay-semantik-fehlerreplay.md) (V-52); höchste vergebene Nummer vor diesem Lauf F-398

**Ausgeführte Läufe im Repo:** `make a-check` (0 Befunde) und `bash tools/test/abdeckung.sh --check` (Exit 0) vor dem Anlegen dieser Datei, `make docs-check` danach. `make gates` lief nicht. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert (`git status --short` leer).

**Mutationen und Sonden:** Nur in einer Kopie (`git archive 1348a9a`) im Scratchpad. Unit-Tests in einem Container aus der Stufe `deps` des `Dockerfile` ohne Netz, dazu `gofmt -l` und `go vet -tags integration ./...` (beide ohne Befund). Integrationstests mit `-test.run` aus einem Image der Stufe `integration` gegen `postgres:17-alpine` (Digest aus `harness/mk/integration.mk`) in einem internen Docker-Netz. Kopie, Images, Container, Netz und Volume sind gelöscht.

| # | Lauf oder Mutation | Ergebnis |
|---|---|---|
| K0 | Kopie ohne Änderung: Unit-Tests `./internal/...` | grün |
| I0 | Kopie ohne Änderung: `TestE2EReplayAbweichung`, `TestE2EReplayExtendedAbweichung` | grün |
| M01 | `note` überschreibt den gemerkten Code bei jedem Fehler | rot (`TestErsterFehlerZaehlt`, `TestReplayNichtVerbrauchtFehler`) |
| M02 | `note` protokolliert nur die erste Meldung | rot (`TestNoteGleichrangig`) |
| M03 | `ErrorResponse` mit `err.Error()` statt Fehlertext | rot (`TestFehlerantwortJeKlasse`, Fall `PGR-E1000`) |
| M04 | Annahmefehler wie vorher nur protokolliert | rot (`TestAnnahmefehler`) |
| M05 | Debug-Zeile ohne Startnachricht unter `error` statt `grund` | rot (`TestOhneStartnachrichtGrund`) |
| M06 | `note` merkt den Code der letzten Meldung | rot (`TestNoteGleichrangig`) |
| M07 | `fail` stellt die letzte statt der ersten Meldung zu | **grün** (F-401) |
| M08 | einzelnes CR kein Umbruch | rot (`TestFehlerEineZeile`) |
| M09 | VT als Umbruch | rot (`TestFehlerEineZeile`) |
| M10 | U+2028 als Umbruch | rot (`TestFehlerEineZeile`) |
| M11 | nicht klassifizierter Fehler als Ursache der letzten statt der ersten Meldung | rot (`TestFehlerGleichrangig`) |
| M12 | verschachtelter Join nicht zerlegt | rot (`TestFehlerGleichrangig`) |
| M13 | innerer klassifizierter Fehler mit eigenem Kopf | rot (`TestFehlerKette`, `TestFehlerGleichrangig`) |
| M14 | Kopf eines umhüllten klassifizierten Fehlers nicht entfernt | rot (`TestFehlerKette`) |
| M15 | `Meldungen` ohne Normalisierung der Umbrüche | rot (`TestFehlerEineZeile`) |
| M16, M17 | Trenner `, ` statt `; ` (nur fremde Fehler, Join in der Kette) | rot (`TestFehlerGleichrangig`) |
| M18 | Klassenname `Konfig` | rot (`TestFehlerKopfJeKlasse`) |
| M19 | Tabulator nicht als Leerraum nach einem Umbruch | rot (`TestFehlerEineZeile`) |
| M20 | Zeile beim Prozessende nur für die erste Meldung | rot (`TestFailJeKlasse`) |
| M21 | Exit-Code der letzten Meldung | rot (`TestFailJeKlasse`) |
| M22 | Stufe ignoriert, immer `info` | rot (`TestRunLogLevel`) |
| M23 | Log-Zeile vor dem Parsen | rot (`TestRunLogLevel`) |
| M24 | Fehler der nie zugeordneten Sessions überschreibt den Code eines Verbindungsfehlers | Unit grün; abgedeckt durch `TestE2EReplayNichtVerbrauchtRangfolge` (Exit-Code 6), dort im vorigen Review als B1 rot |
| M25 | `time` in UTC statt Ortszeit | rot (`TestLoggerOrtszeit`) |
| M26 | `debug` als `info` | rot (`TestLoggerSchwelle`) |
| M27, M28 | Ende von `replay` auf `debug`, Ende von `record` entfällt | rot (`TestRunLogLevel`, `TestRunRecordLogLevel`) |
| M29 | Großschreibung angenommen | rot (`TestParseLogLevelWerte`) |
| M30 | leere Umgebungsvariable gilt als gesetzt | rot (`TestParseLogLevel`) |
| M31 | ungültige Umgebungsvariable ignoriert | rot (`TestParseLogLevelUmgebungUngueltig`) |
| M32 | Umgebungsvariable nach `fs.Parse` ausgewertet (geht der Kommandozeile vor) | rot (`TestParseLogLevelUmgebung`) |
| M33 | Standard `warn` | rot (`TestParseRecord`, `TestParseReplay`) |
| M34 | erste statt letzter Angabe gilt | rot (`TestParseLogLevel`, `TestParseLogLevelUmgebung`) |
| M35 | `warning` als Alias | rot (`TestParseLogLevelWerte`, `TestParseLogLevelUmgebungUngueltig`) |
| M36 | Hilfe ohne Umgebungsvariable | rot (`TestParseLogLevelHilfe`) |
| M37 | `replay` liest die Umgebungsvariable nicht | rot (`TestParseLogLevelUmgebung`) |
| M38 | Parameter-Nummer ab 0 | rot (`TestReplayExtendedAbweichung`, `TestReplayExtendedParameterStelle`) |
| M39 | Anzahlen vertauscht | rot (`TestReplayExtendedParameterStelle`) |
| M40 | Parameterwert in der Diagnose (zwei Formen) | rot (beide Tests) |
| M41 | letzter statt erster abweichender Parameter | rot (`TestReplayExtendedParameterStelle`) |
| M42 | V-52, Unit: `Query` liefert bei Abweichung die aufgezeichneten Antworten | rot (`TestReplayMismatch`) |
| M43 | V-52, Unit: Antworten bei Anfrage über die Aufzeichnung hinaus | rot (`TestReplayMismatch`) |
| E1 | V-52, E2E: Service liefert die aufgezeichneten Antworten, Adapter sendet sie vor `PGR-E5001` | rot (`TestE2EReplayAbweichung`: 1 Ergebnis, `err` nil) |
| E2 | wie E1, Antworten ohne `ReadyForQuery`, also mit `PGR-E5001` beim Client | rot mit dem neuen Test (1 Ergebnis neben `PGR-E5001`); **grün** mit dem Test aus `852c6d0` (Gegenprobe zu V-52) |
| E3 | Adapter protokolliert Parameterwerte auf `debug` | rot (`TestE2EReplayExtendedAbweichung`) |
| M44 | `info`-Zeile am Anfang von `handle` | rot, nur nebenbei: `TestOhneStartnachrichtGrund` zählt die Zeilen der TCP-Probe |
| M44b, M44c | `info`-Zeile je Session (nach `open`) oder je Anfrage | **grün** (F-399) |
| M45 | Warnung `PGR-W3001` ohne `code` | rot (`TestCancelRequest`) |
| M46 | Zeile ohne Startnachricht auf `info` | rot (`TestOhneStartnachrichtGrund`) |
| S1 | `fmt.Errorf("Kontext %w und %w", a, b)` mit `a`, `b` klassifiziert | zwei Meldungen `… a`, `… b`; „Kontext“ fehlt (F-402) |
| S2 | `Errorf(PGR-E4002, fmt.Errorf("K %w / %w", a, x), "aussen")` | `Netzwerk [PGR-E4002]: aussen: a; x`; „K“ fehlt (F-402) |
| S3 | `fmt.Errorf("Kontext: %w", errors.Join(a, b))` | eine Meldung `Recording [PGR-E3001]: Kontext: a; b`, ein Kopf, kein Text fehlt |
| S4 | fremde Hülle mit Umbruch um einen klassifizierten Fehler | `Netzwerk [PGR-E4000]: k: m: x y`, ein Kopf, eine Zeile |

Insgesamt 53 Mutationen (E1 bis E3 eingeschlossen) und vier Sonden. 50 Mutationen sind rot. M07 und M44b/M44c bleiben grün, M24 deckt ein E2E-Test.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-399 | MEDIUM | LH-FA-14.a sagt für `info` zu: „keine Zeile je Verbindung oder Interaktion“. Kein Test prüft das. Eine `info`-Zeile je Session (nach `open`) oder je Anfrage bleibt grün (M44b, M44c). Rot wird nur eine Zeile auf dem Pfad der TCP-Probe (M44), und das nur, weil `TestOhneStartnachrichtGrund` die Zeilen zählt. `TestRunLogLevel` und `TestRunRecordLogLevel` laufen mit beendetem Kontext, also ohne Verbindung. Der Implementer hat die Lücke selbst gemeldet. Nach §3.10 bleibt sie ein Befund: Die Zusage gehört zum neuen Vertrag des Slice und steht in §6 („Stufen und ihr Inhalt“). Dieselbe Klasse war zuletzt MEDIUM (F-384, F-390 bis F-392). | `AGENTS.md` §3.10; LH-FA-14.a §Log-Level; `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`; Reviewer-Skill MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag“ | `spec/spezifikation.md:741`–`:742`; `internal/adapters/driving/pgwire/server.go:102`, `:127`; `internal/bootstrap/bootstrap_test.go:129`; `docs/plan/planning/in-progress/slice-replay-semantik-meldungscodes.md:138` | ja (M44b, M44c) | Zusage ohne Test und Mutation |
| F-400 | MEDIUM | Die Folge-Slices führen die Übernahmen nur in §1. `slice-v1-abschluss-betrieb` nennt den Schlüssel `log_level` (mit `PGR-E2004`) im Absatz „Übernommen aus `slice-replay-semantik-meldungscodes`“. DoD-Punkt 3 nennt aber nur `fail_on_unconsumed`, und auch die §3-Zeile für `internal/adapters/driving/cli` nennt den Schlüssel nicht. `slice-v1-abschluss-einspielen` übernimmt `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` bei `play`. Weder seine DoD noch seine §3-Zeile für die CLI nennen sie. Der Verifier jener Slices misst gegen die DoD. Was nur in §1 steht, kann dort ungeprüft durchgehen. Das ist nach F-388 und F-395 (beide LOW) das dritte Auftreten dieser Klasse, daher MEDIUM. | `AGENTS.md` §3.9 (Folge-Slice); `BEO-REPO/plan-folgt-korrektur-nicht`; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“ | `docs/plan/planning/open/slice-v1-abschluss-betrieb.md:46`, `:63`, `:81`; `docs/plan/planning/open/slice-v1-abschluss-einspielen.md:36`, `:54`–`:56`, `:75`; `docs/plan/planning/in-progress/slice-replay-semantik-meldungscodes.md:57`–`:58` | ja (Lesen) | Übernahme nur in §1/§6, nicht in DoD und Kopf |
| F-401 | LOW | Der Kommentar an `fail` sagt zu, dass die *erste* Meldung des Fehlers als `ErrorResponse` zugestellt wird. Eine Mutation, die die letzte zustellt, bleibt grün (M07). Kein Test gibt einen `errors.Join` in `fail`. Heute erreicht kein Pfad `fail` mit einem Join: `OpenSession`, `Query` und `ClientMessage` liefern Einzelfehler, und der Join aus `CloseSession` geht nur an `note`. Welche Meldung bei mehreren an den Client geht, nennt weder §6 noch `SPEC-034` *Zustellung an den Client*. Der Code legt es ohne Beobachtung fest. Es ist LOW, nicht HIGH nach §3.12, weil der Pfad heute unerreichbar ist und die Wahl aus „gemerkt wird der erste“ folgt. | `AGENTS.md` §3.11, §3.12; `SPEC-034` §Ausgabe *Zustellung an den Client*; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `internal/adapters/driving/pgwire/server.go:688`–`:696`, `:711`; `internal/adapters/driving/pgwire/server_meldung_test.go:69` | ja (M07) | Kommentar-Zusage ohne Test, Pfad unerreichbar |
| F-402 | LOW | `nebeneinander` und `ursache` erkennen gleichrangige Fehler an der Schnittstelle `Unwrap() []error`. Die trägt auch `fmt.Errorf` mit mehreren `%w`. Ein solcher Hüllfehler wird als gleichrangig zerlegt, und sein eigener Text fällt weg: aus `fmt.Errorf("Kontext %w und %w", a, b)` werden zwei Meldungen ohne „Kontext“ (S1), in einer Kette fällt „K“ weg (S2). Die Kommentare nennen nur `errors.Join`, und §6 nennt die Form nicht. Heute erzeugt der Code keinen Fehler mit mehreren `%w` (`grep '%w'`), die Wirkung ist also latent. Ein einfaches `%w` um einen Join behandelt der Code wie eine Kette (S3), ohne Textverlust. | `SPEC-034` §Ausgabe *Fehlerkette*, *Gleichrangige Fehler*; `AGENTS.md` §3.11; Maintainability | `internal/hexagon/model/fehler.go:80`–`:84`, `:92`–`:99`, `:162`–`:177` | ja (S1, S2) | Typ-Erkennung weiter als die Zusage im Kommentar |
| F-403 | LOW | `Berührte Spec-Stellen` nennt keine Stelle der Sicht. Der Diff ändert aber den CLI-Adapter (neue Option und Umgebungsvariable, `ARC-005`), die Composition Root (Logger nach dem Parsen, Zeile beim Prozessende, `ARC-009`), den PGWire-Adapter (`note`, `fail`, Annahmefehler, `ARC-006`) und das Domain Model (`Meldung`, `Meldungen`, `ARC-001`). Die Aussagen der Sicht bleiben unverändert. Das ist das zweite Auftreten nach F-386. | `AGENTS.md` §3.9 (Kopf); Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice; Maintainability | `docs/plan/planning/in-progress/slice-replay-semantik-meldungscodes.md:16` | ja (Lesen) | Kopf nennt die berührten Sicht-Stellen nicht vollständig |
| F-404 | LOW | `TestLoggerOrtszeit` deklariert `LH-FA-14/Boundary`. Das Kriterium lautet „bleibt die Ausgabe auf der konfigurierten Detailstufe; die Detailstufe ist einstellbar“. Der Test prüft die Zeilenform von `time` (Ortszeit, Millisekunden, Zonenversatz), nicht die Detailstufe. Die Zeile in `abdeckung-unit.md` belegt damit für die Boundary etwas, das sie nicht prüft. Die Boundary tragen ohnehin `TestRunLogLevel`, `TestLoggerSchwelle` und `TestParseLogLevel`, die Wirkung ist klein. Die Klasse ist dieselbe wie bei F-383. | [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) §Konsequenzen; `AGENTS.md` §3.11; LH-FA-14 Akzeptanzkriterien | `internal/bootstrap/bootstrap_test.go:270`; `docs/user/abdeckung-unit.md:80` | ja (Lesen) | Abdeckungs-Deklaration verfehlt das Kriterium des Lastenhefts |
| F-405 | INFO | Plan §3 begründet „`record.go` (`CloseSession`) bleibt unverändert: der `errors.Join` ist die Reihenfolge des Entstehens“. In `CloseSession` entsteht `upErr` (Schließen zum Upstream) aber vor dem Schreiben des Recordings und steht im Join dahinter. Das hat heute keine beobachtbare Folge: `upErr` kommt aus `net.Conn.Close`, ist nicht klassifiziert und wird nie eine eigene Meldung. Die Reihenfolge würde zählen, sobald `Upstream.Close` einen klassifizierten Fehler liefert. | `SPEC-034` §Ausgabe *Gleichrangige Fehler*; `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-replay-semantik-meldungscodes.md:95`; `internal/hexagon/services/record.go:397`, `:407`; `internal/adapters/driven/postgres/upstream.go:210`–`:218` | ja (Lesen) | Plan-Begründung trifft den Code nicht genau |

## Antwort auf die Schwerpunkte

1. **`model.Meldungen` und `Error()`.**
   - *Ein Kopf je Kette:* `Error()` und `Meldungen` setzen den Kopf einmal. Innere klassifizierte Fehler tragen nur `ursache` bei, auch durch eine fremde Hülle hindurch (M13, M14 rot).
   - *Gleichrangig:* Ein Join wird rekursiv zerlegt (M12), je klassifiziertem Fehler eine Meldung in Join-Reihenfolge. Fremde Fehler hängen an der ersten (M11), nur fremde ergeben `PGR-E1000` mit `; ` (M16).
   - *Doppelt oder gar nicht:* Ein Join in einem Join ist nicht doppelt und fehlt nicht. Ein `%w` um einen Join ist eine Meldung mit einem Kopf: Code des ersten klassifizierten, alle Texte als Ursache (S3). Der zweite Code erscheint dann nicht als eigene Meldung, das entspricht *Fehlerkette*. Ein `Errorf` um einen Join verhält sich ebenso (`TestFehlerGleichrangig`, letzter Fall). Text geht nur bei mehreren `%w` verloren (F-402).
   - *Umbrüche:* LF, CR LF und CR samt Leerraum werden ersetzt, VT, FF, U+0085, U+2028 und U+2029 bleiben (M08 bis M10, M15, M19 rot). Klassennamen und Exit-Codes je Klasse sind geprüft (M18), ebenso die Code-Tabelle aus dem Quelltext (`TestCodeTabelle`).
   - [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md) ist eingehalten: `000` als Rückfall, Kopf `<klasse> [<code>]: `. Der Re-Evaluierungs-Trigger ist in §6 behandelt, die Rangregel steht in `SPEC-034`, wie `AGENTS.md` §3.8 es verlangt.
2. **`note`/`fail`.**
   - *Merken:* Der erste Code wird unter der Sperre über Verbindungen hinweg gemerkt (M01 rot über `TestErsterFehlerZaehlt`), bei einem Join der der ersten Meldung (M06).
   - *`ErrorResponse`:* derselbe Text wie `error` (M03), SQLSTATE nach Klasse (`TestFehlerantwortJeKlasse`, `DeepEqual` schließt weitere Felder aus). Lücke bei mehreren Meldungen: F-401.
   - *Annahmefehler:* `PGR-E4000` mit Kopf, gemerkt, Serve nimmt weiter an (M04).
   - *Bootstrap:* `fail` schreibt je Meldung eine Zeile, Exit-Code der ersten (M20, M21). Die Zeile beim Prozessende hängt nicht am Logger.
3. **`--log-level`.**
   - *Reihenfolge:* Standard `info`, dann die Umgebungsvariable vor `fs.Parse`, dann die Kommandozeile, letzte Angabe gilt (M30, M32 bis M34). Ungültige Variable auch neben gültiger Option ist `PGR-E2001` (M31).
   - *Logger und Startfehler:* Der Logger entsteht erst nach `cli.Parse` (M23 rot). Startfehler auf `error`, `warn`, `debug` und mit Variable `error` sind genau eine Zeile ohne `time=` (`TestRunStartfehlerJeStufe`). Für `info` belegt das der ältere Bestand.
   - *Hilfe:* Sie geht Variable und Wert vor (`TestParseLogLevelHilfe`, Variable `INFO`). `version` liest die Variable nicht, weil `Parse` für `version` `logLevelOption` nicht ruft. Ein Test mit gesetzter ungültiger Variable fehlt für `version`. Das ist kein Befund, weil der Pfad die Variable gar nicht berührt.
   - *Vor dem Kommando:* `--log-level` vor dem Kommando ist `PGR-E2001` (`TestParseLogLevelWerte`).
4. **Parameter.** Nummer ab 1, erster abweichender, beide Anzahlen in der Reihenfolge erwartet/empfangen (M38, M39, M41). Kein Wert in der Diagnose (M40), auch nicht auf `debug` im Log (E3). Die einzige Stelle, die `Params` in einen Text bringt, ist `parameterStelle`, und die nennt nur Zahlen. `param_types`, `param_formats` und `result_formats` nennen weiter nur das Feld, wie in §6 als akzeptiertes Negativ notiert.
5. **V-52.** Beide Tests werden aus dem richtigen Grund rot.
   - *Unit:* `TestReplayMismatch` wird rot, wenn `Query` bei Abweichung oder nach dem Ende Antworten liefert (M42, M43).
   - *E2E:* `TestE2EReplayAbweichung` wird rot, wenn der Adapter die aufgezeichneten Antworten vor `PGR-E5001` sendet (E1, E2).
   - *Gegenprobe:* Der Test aus `852c6d0` bleibt in der Form E2 grün. Das ist genau die Lücke aus V-52, und der neue Test schließt sie.
6. **§3.10, Stichprobe.** 53 eigene Mutationen, 50 rot. Grün blieben M07 (F-401) und M44b/M44c (F-399). Die gemeldete Lücke „auf `info` keine Zeile je Verbindung“ ist bestätigt (F-399). Eine weitere Zusage ohne Mutation ist die Wahl der ersten Meldung in `fail` (F-401).
7. **§3.11.**
   - *Hilfetext:* „Standard info; eine Stufe zeigt auch die strengeren; Umgebungsvariable …“ ist geprüft (M22, M26, M36).
   - *Handbuch:* Stufen-Tabelle, Strenge des Werts, Zeilenform, Startfehler, Parameter-Nummer, Fehlertext einzeilig und SQLSTATE sind jeweils durch einen Test belegt. Die Zeile zu `info` sagt nur „Start und Ende“ zu, nicht „keine Zeile je Verbindung“, und reicht deshalb nicht weiter als die Prüfung.
   - *Kommentare und Deklarationen:* Zu weit reichen der Kommentar an `fail` (F-401), die Kommentare zu `errors.Join` (F-402) und eine Abdeckungs-Deklaration (F-404).
8. **`time.Local`.** Kein Test im Repo ruft `t.Parallel` (`grep -rn t.Parallel internal test` leer). `Run` wartet auf das Ende von `Serve` (`<-done`), bevor es zurückkehrt, also loggt keine Goroutine eines früheren Tests weiter, während `TestLoggerOrtszeit` `time.Local` umsetzt. Pakete laufen in getrennten Prozessen. Kein Flackern zu erwarten.
9. **Plan und Folge-Slice.**
   - *Plan:* §1, §3 und §6 folgen dem Code. Jede Datei des Diffs steht in §3, die drei Rückgaben sind in §6 vor dem Code entschieden (`f2c428a` vor `1348a9a`).
   - *Folge-Slices:* `slice-v1-abschluss-betrieb` ist in §1 und §3 nachgezogen (allgemeiner Leser ersetzt beide Einzel-Leser), seine DoD nicht (F-400). Für `slice-v1-abschluss-einspielen` gilt dasselbe (F-400).
   - *Kopf:* F-403. *Begründung in §3 zu `CloseSession`:* F-405.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/model` | geprüft. Kopf, Kette, gleichrangige und fremde Fehler, Umbrüche und Code-Tabelle folgen `SPEC-034` (M08 bis M19 rot). Befund F-402. |
| `internal/hexagon/services` | geprüft, ohne Befund. `parameterStelle` nennt keine Werte, und der V-52-Nachweis trägt (M38 bis M43 rot). |
| `internal/adapters/driving/pgwire` | geprüft. `note` und Annahmefehler folgen *Gleichrangige Fehler* und LH-FA-13.b, `grund` an der Debug-Zeile (M01 bis M06, M45, M46 rot). Befund F-399, F-401. |
| `internal/adapters/driving/cli` | geprüft, ohne Befund. Wertemenge, Strenge, Umgebungsvariable, Mehrfachangabe und Hilfe sind belegt (M29 bis M37 rot). |
| `internal/bootstrap` | geprüft. Logger nach dem Parsen, Zeile beim Prozessende je Meldung, Ortszeit (M20 bis M28 rot). Befund F-399 (kein Test mit Verbindung), F-404. |
| Core-Reinheit, [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) | geprüft, ohne Befund. `make a-check` 0 Befunde. `internal/hexagon/model` importiert neu nur `regexp` und `strings`. |
| [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md) | geprüft, ohne Befund. Form, Rückfall `PGR-E1000`, Kopf; ADR nicht verändert. |
| `test/integration` | geprüft, ohne Befund. Stabil (I0), V-52 belegt (E1, E2), Parameterwerte auf `debug` belegt (E3). |
| `docs/user/` | geprüft. Abdeckungstabellen entsprechen den Deklarationen (`abdeckung.sh --check`). Handbuch ohne zu weite Zusage. Befund F-404. |
| `spec/spezifikation.md` — Strata | geprüft, ohne Befund. Keine ADR-, Slice- oder Commit-Nennung im Diff, Historie nachgetragen. |
| `spec/lastenheft.md`, `spec/architecture.md` | geprüft, ohne Befund. Nicht im Diff. |
| `docs/plan/planning/` — Hard Rule 3.9 | geprüft. §1, §3 und §6 folgen dem Code. Befund F-400, F-403, F-405. |
| Hard Rule 3.3 | geprüft, ohne Befund. Kein Move im Diff. |
| Hard Rule 3.5, 3.6 | geprüft, ohne Befund. Keine ADR und keine Gate-Konfiguration im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft, ohne Befund. Neue Kommentare sind Zusagen oder Kopplungen im Indikativ. Zur Reichweite siehe F-401 und F-402. |
| Hard Rule 3.8 | geprüft, ohne Befund. Die Rangregel der Fehlerkette steht in `SPEC-034`, keine neue ADR. |
| Hard Rule 3.10 | geprüft. 50 von 53 Mutationen rot. Befund F-399, F-401. |
| Hard Rule 3.11 | geprüft. Befund F-401, F-402, F-404. |
| Hard Rule 3.12 | geprüft. Die Randformen aus §6 sind vor dem Code entschieden (`ac12870`, `f2c428a` vor `1348a9a`), auch die drei Rückgaben. Der Code folgt ihnen. Nicht genannt und im Code festgelegt sind nur Formen auf heute unerreichbaren Pfaden (F-401, F-402). |
| Commit-Messages `ac12870` bis `1348a9a` | geprüft, ohne Befund. Jede nennt `slice-replay-semantik-meldungscodes` und `LH-*`-Kennungen, keine nennt `SPEC-*` oder `ARC-*`. |
| Register `BEO-REPO/*` | geprüft. F-399 betrifft `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`. F-401, F-402 und F-404 betreffen `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. F-400 und F-403 betreffen `BEO-REPO/plan-folgt-korrektur-nicht` (verkörpert, §3.9; Retirement-Check betroffen). `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` bekommt keinen dritten Beleg: Die Rückgaben kamen vor dem Code. |

**Summary:** 0 HIGH · 2 MEDIUM (F-399, F-400) · 4 LOW (F-401 bis F-404) · 1 INFO (F-405). Wiederkehrende Klassen: „Zusage ohne Test und Mutation“ (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`), „Übernahme nur in §1/§6, nicht in DoD und Kopf“ (`BEO-REPO/plan-folgt-korrektur-nicht`, drittes Auftreten).
