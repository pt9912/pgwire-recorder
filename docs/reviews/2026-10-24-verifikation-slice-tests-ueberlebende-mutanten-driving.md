# Verifikation slice-tests-ueberlebende-mutanten-driving

**Rolle:** Verifier. **Stand:** main `8b621273`, Basis `7ed55a23`, Architect-Entscheidungen `44cb42c3`.
**Urteil: bestanden.** Kein HIGH und kein MEDIUM. Zwei LOW, zwei INFO; ein Punkt geht an den Architect (nicht blockierend, siehe V-1).

## Vorgehen

Alle Mutationen in einer frischen Kopie des Arbeitsbaums (ohne `.git`, eigener Pfad unter dem Scratchpad), genau eine Textersetzung je Mutant mit Trefferzahl 1 geprüft, danach `touch`. Lauf: `go test -count=1` im Image `golang@sha256:b8c9b82a…` (das Image der Stufe `test`), Bind-Mount, `--network=none`, Modulcache einmalig mit Netz geladen und schreibgeschützt eingebunden. Kein Build-Kontext, kein `make gates`, keine Gegenprobe.

## Mutationen am Endstand

| Zusage | Mutation | mit dem neuen Test | ohne ihn (`-skip`) |
|---|---|---|---|
| Spalten-Metadaten (`LH-FA-02.a`, `LH-FA-03.a`) | T3: `TableOID: c.TableOID` → `0` in `ToMessage` | rot: `TestToMessageSpaltenMetadaten` „Tabellen-OID 0, erwartet 16384“ (id und name) | grün (`ok … 21,05 s`) |
| dieselbe | T5: `TableAttributeNumber` → `0` | rot: „Spaltennummer 0, erwartet 1“ und „… 7“ | grün (`ok … 21,04 s`) |
| Startphase (`LH-FA-13.a`) | K4: `WithoutCancel` aus `Listen` entfernt | rot: `…/record_Name_localhost`, `…/replay_Name_localhost`; `IP-Adresse`, `Name löst nicht auf` und `Port belegt` bleiben grün | grün (`ok … 6,52 s`) |
| Startfehler bleibt: schluckt jeden | `if ctx.Err() != nil { return 0 }` nach `Listen`, je `record` und `replay` | rot: genau `Name löst nicht auf` und `Port belegt` des betroffenen Kommandos | — |
| Name, der nicht auflöst | bei beendetem Kontext und `*net.DNSError` `return 0`, je Kommando | rot: genau `…Name_löst_nicht_auf` des Kommandos | — |
| Port belegt | bei beendetem Kontext und Fehler ohne `DNSError` `return 0`, je Kommando | rot: genau `…Port_belegt` des Kommandos | — |
| Frist `conn.gesetzt` | W: `SetReadDeadline` im Wächter entfernt | rot nach 10,01 s: „der Prüfling setzt die Lesefrist nicht binnen 10 s“ (`server_extended_test.go:851`), kein Paket-Timeout | — |
| Frist Ende von `Run` | `<-ctx.Done()` in `betreiben` → `<-make(chan struct{})` | rot nach 30 s, nicht durch den Paket-Timeout (300 s gesetzt): `TestRunRecordSchreibfehlerJeStufe` (30,05 s), `TestRunReplayTLS` (`tls_test.go:106`, drei Fälle, 90,3 s), `TestRunTLSLogOhneSchluessel` (`tls_test.go:121`, zwei Fälle, 60,6 s) | — |
| Frist Ende von `Serve` | `Serve` kehrt nach `s.wg.Wait()` nie zurück | rot nach 10,01 s: „Serve endet nicht binnen 10 s nach dem Schließen des Listeners“ (`server_tls_test.go:438`) | — |

Jeder rote Lauf nennt die Zeile des Tests und den erwarteten Grund, nicht einen Aufbaufehler.

**G2, G3, G5, G6 am Endstand**, je engster Test, alle rot:

| Mutant | Test | Ergebnis |
|---|---|---|
| G2/R2 `&& false` in `verbindungsende(err)` | `TestProtokollrandClientNachrichten` | rot (Fälle „Abbruch mitten in einer Query“, „im Kopf einer Query“) |
| G2/R4 `return false` → `continue` | `TestReplayEndeTerminateMittenInDerGruppe` | rot nach 5,00 s: „Verbindung endet nicht“ |
| G3/R3 `s.fail(…)` entfernt | `TestProtokollrandClientNachrichten` | rot (Fälle `Parse ohne Nullbyte`, `Nachricht p ohne Nullbyte`) |
| G5/X4 `s.sendFailed(err)` entfernt | `TestReplayEndeSchreibfehlerBleibtUnverbraucht` | rot: „erster Fehler "", erwartet PGR-E4003“ |
| G6/Z3 `return folgeEnde` → `return folgeWeiter` | `TestReplayEndeNurAlsLetzteAntwort` | rot nach 5,00 s: „read pipe: i/o timeout“ |

**G7** (`case <-fertig:` in `replayWaechter` entfernt, ganzes Paket `pgwire`): grün (`ok … 20,97 s`). Das ist das akzeptierte Negativ.

**Zeittests, 20 Läufe ohne Mutant** (`-count=20`): `pgwire` (`TestReplayHerunterfahrenSpaeteFrist`, `TestToMessageSpaltenMetadaten`, alle `TestTLS…`) `ok … 50,5 s`; `bootstrap` (`TestRunSignalVorDemOeffnen`, `TestRunRecordSchreibfehlerJeStufe`, `TestRunReplayTLS`, `TestRunTLSLogOhneSchluessel`) `ok … 17,7 s`. Beide Pakete liefen gleichzeitig, also unter Last, ohne Flattern. Die 50 Läufe des Implementers (§7) habe ich nicht wiederholt; 20 sind die geforderte Zahl.

## DoD

- Frist beim Warten: bestätigt (W rot nach 10 s ohne Paket-Timeout; alle Stellen der Tabelle in §6 tragen eine Frist).
- G4: bestätigt (T3, T5 je für sich rot, ohne Test grün).
- G1: bestätigt (K4 rot nur auf den beiden Namensfällen; die drei Startfehler-Mutanten färben genau ihre Fälle).
- G7 kein Test: bestätigt. §7 nennt den Mutanten „akzeptiertes Negativ, nicht äquivalent“, Grund „keine Zusage in `SPEC-032`“, ohne Invariante, Ort und Sonde. Das Wort „äquivalent“ steht nur in der Negation. `replayWaechter` trägt keinen Satz zum Ende der Sitzung (`server.go` im Diff unverändert); `SPEC-032` führt die Nicht-Zusage (Absatz „Nicht zugesichert ist, was eine beendete Verbindung im Prozess zurücklässt“).
- `make abdeckung-check` Exit 0, `make kopf-check` Exit 0, `make docs-check` 1455 Dateien, 0 Befunde (alle im Worktree am Stand von main). `make gates` nicht gefahren, wie beauftragt; `git diff 7ed55a23..HEAD` zeigt seit der Gate-Lauf-Angabe keine Änderung unter `internal/`, `cmd/`, `test/`, `harness/` oder `tools/` außer den Testdateien und `spec/`, die in den Lauf des Implementers eingingen.
- Offene DoD-Punkte (Review, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) sind Sache von Review und Closure; die Felder der Closure-Notiz in §7 stehen noch auf `<…>`.

## Entscheidungs-Konformität (Architect, `44cb42c3`)

- `SPEC-032` (G7), `LH-FA-03.a` und `LH-FA-02.a` (Spalten-Metadaten), `LH-FA-13.a` Startphase (Namensauflösung gehört zum Öffnen; Plattformgrenze genannt), `SPEC-038` (Lauf der Kommandozeile, Kanal des Testgeschirrs): Der Diff setzt sie wie entschieden um. Der Doc-Kommentar zum Test der Spalten-Metadaten sagt nur zu, was der Test prüft (kein Satz „Weg ist für Record und Replay derselbe“); der Wortlaut der Spezifikation sagt den Weg für den Client zu, der Test prüft `ToMessage`, das ist dieselbe Abbildung für beide Modi.
- Zugang: `TestToMessageSpaltenMetadaten` steht in `package pgwire_test` und ruft `pgwire.ToMessage`, `TestRunSignalVorDemOeffnen` in `package bootstrap_test` und ruft `bootstrap.Run` (`SPEC-049` Punkt 7); die Export-Test-Brücke wächst nicht.
- `0752b3ec` ändert nur `internal/bootstrap/bootstrap_test.go` (+2 −3): bestätigt, Produkt-Code unverändert.
- Abdeckung: `docs/user/abdeckung-unit.md` trägt `TestToMessageSpaltenMetadaten` unter `LH-FA-02` und `LH-FA-03` (Happy) und `TestRunSignalVorDemOeffnen` unter `LH-FA-13` (Boundary), gleich den Deklarationen im Test; `abdeckung-gesamt.md` ändert nur die Zelle `LH-FA-03` (Unit kommt hinzu).
- Marke 1500 (`AGENTS.md` §3.14), Befehl wörtlich, Basis `7ed55a23`, Pfade `internal/adapters/driving/pgwire internal/bootstrap spec`: **147**. Marke nicht erreicht; die Zahl in §7 stimmt.
- Plan gegen Diff: §3 nennt `server.go` „nur Kommentar“, der Diff ändert es nicht (Kommentar nicht nötig, keine Abweichung). Alle anderen Dateien aus §3 stimmen mit dem Diff überein; §1 und §6 nennen `tls_test.go:100` und das Ende von `Serve` (Nacharbeit F-1, F-2).

## Befunde

**V-1 (LOW) — zu beurteilen: Bestätigung des Architect zu `tls_test.go:100` und zu `<-fertig` in `server_tls_test.go`.** Die Entscheidung `SPEC-038` trägt beide, ohne Rückfrage zur Konformität: `tls_test.go:100` ist dieselbe Konstellation wie `:184` (Lauf der Kommandozeile in einer Goroutine, Ergebnis über `ende`), und der Absatz „Der Lauf der Kommandozeile …“ nennt sie ausdrücklich. `<-fertig` wartet auf das Ende von `Serve`, das der Prüfling herbeiführt (Ende eines gestarteten Prozesses nach dem Wortlaut der Regel). Selbst wer es als Testgeschirr liest, findet im selben Absatz, dass eine Frist dort „nicht verboten“ ist. Kein Befund gegen den Code. **Offen bleibt nur die Form:** §1, §2 und die Tabelle in §6 tragen für beide Zeilen „Bestätigung des Architect steht aus“; der Text soll vor der Closure entweder bestätigt oder gestrichen werden, sonst steht in `done/` ein ausstehender Satz. Das ist eine Zeile des Architect (Einzelbestätigung), keine Rückfrage mit offenem Ausgang.

**V-2 (LOW) — Architect zurückfragen: `play_anmeldung_test.go` `<-danach` ohne Frist (Zeilen 171 und 242; Zeile 192 trägt `time.After`).** `SPEC-038` entscheidet es nicht eindeutig. Der Satz „Ein Test-Double, das auf eine Freigabe durch den Test wartet, fällt nicht darunter; der Test, der danach auf den Prüfling wartet, schon“ trifft hier zu: `danach` füllt der Test-Server (Double), aber mit dem, was der Prüfling gesendet hat; der Test wartet also auf ein Ereignis, das der Prüfling herbeiführt. Der Absatz zum „Kanal, den nur das Testgeschirr füllt“ passt dagegen auf ein `Accept`, nicht auf einen Wert, der aus den Bytes des Prüflings entsteht. Der Implementer hat die Stelle in §7 vermerkt und nicht angefasst (Plan: Entscheidung beim Architect, sonst `slice-harness-lint-warten-ohne-frist`). Das ist nach dem Text kein Mangel dieses Slice, weil die Stelle nicht auf der Liste des Architect stand; die Antwort legt aber fest, ob die Liste „alle Tests dieser beiden Schichten“ (§1) vollständig ist. **Empfehlung:** Architect fragen, in einem Satz (Frist Pflicht oder Geschirr); ist es Pflicht, ist es eine dritte Stelle in diesem Slice oder eine Sendung an `slice-harness-lint-warten-ohne-frist` mit Eintrag dort (§3.13). Nicht blockierend für die Verifikation: Die DoD fordert die Frist für die Stellen „welche der Architect festlegt“.

**V-3 (INFO) — Abdeckungs-Deklaration von `TestRunSignalVorDemOeffnen` führt nur `Boundary`** (Review F-4). Der Text der Zeile nennt auch die Startfehler; `make abdeckung-check` ist grün. Die Deklaration sagt weniger zu, als der Test prüft, kein Verstoß gegen §3.11. Wer will, ergänzt `LH-FA-13/Negative` mit eigener Zeile.

**V-4 (INFO) — Closure-Felder in §7 stehen auf Platzhaltern** (`Was hat funktioniert` bis `Drei Paarungen`), die DoD-Punkte für Review, Closure-Notiz, Register, Risiko-Ausgänge und Paarungen sind offen. Erwartet vor Closure, nicht Teil dieser Verifikation. Die vier Risiken in §6 tragen noch kein Ausgang; die Belege dafür (Läufe ohne Mutant, Aktualität, Kopie unter eigenem Pfad) liegen in §7 und in diesem Bericht vor.

## Nicht geprüft

Kein `make gates`, keine Gegenprobe, kein Lauf auf macOS oder Windows (die Grenze steht in `LH-FA-13.a`). Die 50 Läufe des Implementers nicht wiederholt (20 Läufe gefahren). Der Review-Report ist nicht gelesen.
