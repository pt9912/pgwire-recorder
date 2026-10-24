# Review-Report: slice-tests-ueberlebende-mutanten-driving — 2026-10-24

**Review-Art:** Code — Code-Review gegen Plan und Konventionen (Hard Rules `AGENTS.md` §3.9 bis §3.14).

**Gegenstand:** `slice-tests-ueberlebende-mutanten-driving`, Diff `7ed55a23..b0aeec7a` (Zweig `impl-drv1`, von `main` aus gelesen); Commits `d24e7c98`, `020c5635`, `02b7d7c5`, `be8ca1f4`, `f841bb03`, `0752b3ec`, `b0aeec7a`; Entscheidungen des Architect in `44cb42c3`.

**Skill:** `.harness/skills/reviewer.md` @ `4155023d` ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-24

> **Zitier-Form.** Kennung statt Adresse; Baseline-Stellen als Tag und Pfad: `v6.18.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice.

**Eingangs-Kontext:**

- Slice-Plan `slice-tests-ueberlebende-mutanten-driving` (§1 bis §8, darin §7 *Belege des Implementers*)
- [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) (Black-Box-Tests über die Export-Test-Brücke)
- `LH-FA-02`, `LH-FA-03`, `LH-FA-13`, `LH-FA-18`, `LH-QA-07`; Spezifikation `LH-FA-02.a`, `LH-FA-03.a`, `LH-FA-13.a`, `SPEC-032`, `SPEC-038` (*Warten in Tests*), `SPEC-049` Punkt 7
- `AGENTS.md` (Hard Rules, §3.9 bis §3.14)

Nicht gefahren: kein Docker, keine Mutation (auf `main` läuft ein voller `make gates`). Die Mutationstabelle ist gegen den Test-Code gelesen, nicht nachgefahren; die Läufe aus §7 sind Angaben des Implementers.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Im PGWire-Adapter wartet `TestTLSNacheinanderUndGleichzeitig` im `defer` ohne Frist auf das Ende von `Serve` (`<-fertig`, Goroutine mit `s.Serve`). Die Stelle steht weder in der Tabelle des Architect noch in §7. Hängt `Serve`, endet der Test erst mit der Zeitüberschreitung des Pakets. | SPEC-038 *Warten in Tests*; Plan §2 Punkt 1 („Tests dieser beiden Schichten, die ohne Frist … warten“); `AGENTS.md` §3.12 (nicht genannte Randform) | `internal/adapters/driving/pgwire/server_tls_test.go` · `defer func() { cancel(); l.Close(); <-fertig }()` | nein — kein Gate prüft es vor `slice-harness-lint-warten-ohne-frist`; Mutation: `Serve` ignoriert das Ende von `ctx`, der Test hängt | Warten ohne Frist bleibt in der Schicht stehen |
| F-2 | LOW | Die Stelle `tls_test.go:100` (`TestRunReplayTLS`) steht in §3 und §7, aber nicht in §1, §2 und der Tabelle der Randformen in §6 (dort vier Kandidaten). Die Bestätigung des Architect, um die der Implementer bittet, steht aus; der Plan trägt die Randform erst seit `be8ca1f4`, nach dem Code-Commit `020c5635`. | `AGENTS.md` §3.12 (Randform erst im Code-Commit eingetragen), §3.9 | Plan · `tls_test.go:184` (§1 und §2 nennen vier Kandidaten) gegenüber §3 · `` `tls_test.go:100` `` | nein — Lesart des Plans; `make kopf-check` prüft es nicht | Randform erst im Code-Commit im Plan genannt |
| F-3 | INFO | `play_anmeldung_test.go` liest `<-danach` (Zeilen 171, 242) nach dem Ende von `Run`. Der Kanal wird vom Test-Server gefüllt, aber erst, wenn der Prüfling die Verbindung beendet; das liegt zwischen „Testgeschirr“ und „der Prüfling führt es herbei“ in `SPEC-038`. Die Ausnahme des Implementers ist vertretbar, aber nicht begründet entschieden. Zuständig: Architect, sonst `slice-harness-lint-warten-ohne-frist`. | SPEC-038 *Warten in Tests* | `internal/bootstrap/play_anmeldung_test.go` · `if got := <-danach; got != "X"` | nein | Warten auf Kanal des Testgeschirrs, Grenzfall |
| F-4 | INFO | Die Abdeckungs-Deklaration von `TestRunSignalVorDemOeffnen` führt nur den Pfad `Boundary`. Die Fälle `Name löst nicht auf` und `Port belegt` sind Negativpfade derselben Zusage; sie stehen in keiner Zeile von `docs/user/abdeckung-unit.md` als `Negative`. Die Tabelle sagt weniger zu, als der Test prüft; kein Verstoß gegen §3.11. | `AGENTS.md` §3.11 | `internal/bootstrap/bootstrap_test.go` · `// Abdeckung: LH-FA-13/Boundary` | ja — `make abdeckung-check` (grün laut §7) | Deklaration enger als der Test |
| F-5 | INFO | Commit `0752b3ec` ändert nur `internal/bootstrap/bootstrap_test.go` (Argumente von `bootstrap.Run` vor der Goroutine gebildet, `contextcheck`). `Run` und `internal/bootstrap/bootstrap.go` stehen nicht im Diff, `server.go` ebenso nicht: das Produktverhalten ist unverändert, §1 *Produkt-Code* ist gewahrt. | Plan §1 (*Ausdrücklich NICHT: Produkt-Code*) | `git diff 7ed55a23..HEAD --stat` · nur `*_test.go`, `spec/`, `docs/` | ja — `git diff --stat` | — |
| F-6 | INFO | §7 nennt als Summe der Größe 138 Zeilen, der Befehl aus `AGENTS.md` §3.14 ergibt am Kopf `b0aeec7a` 137 (ohne Review-Report). Marke 1500 ist weit entfernt. | `AGENTS.md` §3.14 | Plan §7 · `Summe 138 Zeilen vor dieser Closure-Notiz` | ja — Befehl aus §3.14 | Größenzahl um eins abweichend |

## Zu den fünf Fragen des Auftrags

**(a) Mutationstabelle gegen den Test-Code.** Jede Zeile ist gegen den Code plausibel; keine Zeile ist der Test-Seite nach unerklärt rot.
- G4: `TestToMessageSpaltenMetadaten` nutzt je Spalte verschiedene Werte (16384/1, 16390/7, 0/0) und vergleicht beide Felder einzeln. T3, T5 und die `max(…,1)`-Mutationen (Spalte `n` mit 0) färben unabhängig voneinander. Ein Vertauschen beider Felder wäre ebenfalls rot (Tabellen-OID ≠ Spaltennummer). Der Test geht über `ToMessage` (Export-Brücke, `SPEC-049` Punkt 7); `send` ruft für jede Antwort dieselbe Funktion, einen zweiten Pfad für `RowDescription` gibt es nicht (`server.go`, `toMessage`).
- G1: Die Fälle `Name localhost` fangen K4 (Listen mit gekündigtem Kontext → Exit 4). `IP-Adresse` bleibt unter K4 grün und ist als Kontrolle ausgewiesen; das trägt die Zeile. `Name löst nicht auf` und `Port belegt` sind je Kommando eigene Fälle; die Mutationen „bei beendetem Kontext `return 0`“ (einzeln für `record` und `replay`, nach `*net.DNSError` getrennt) trennen sie. Die Bedingung `code != exit || Contains(PGR-E4001) != (exit==4)` prüft „kein PGR-E4001“ als eigene Bedingung neben dem Exit-Code (Zeile 468 der Tabelle). Das erfüllt `AGENTS.md` §3.10 *je Element* für die Aufzählung „Name, der nicht auflöst · Port, der nicht zu öffnen ist“. Die Meldung `Exit-Code 0, erwartet 0; stderr …` ist als Ausgabe lesbar, aber missverständlich; kein Befund.
- Frist-Tests: Jede der fünf Stellen hat eine Frist als Literal (10 s, 30 s; SPEC-038 erlaubt höchstens 60 s) und eine Meldung, die das ausgebliebene Ereignis nennt. Die Mutation `<-ctx.Done()` → `<-make(chan struct{})` in `betreiben` färbt die drei Bootstrap-Stellen über die Frist; die Mutation W färbt `TestReplayHerunterfahrenSpaeteFrist` nach 10 s statt über das Paket-Limit. `server_tls_test.go` (`paar`) ist Testgeschirr, die Mutation ändert den Code des Tests; der Plan nennt das ausdrücklich („Frist erlaubt, nicht Pflicht“).
- Eine Beobachtung zur Zeit: `.invalid` als Namen löst ohne Netz sofort auf; bei gekündigtem Kontext und aktivem `WithoutCancel` hinge der Fall an der Zeit eines Resolvers, wenn DNS erreichbar, aber stumm ist. Die Frist von 30 s fängt das; kein Befund.

**(b) `tls_test.go:100`.** Vereinbar mit der Entscheidung des Architect, keine Rückgabe nötig. Der neue Absatz in `SPEC-038` („Lauf der Kommandozeile … in einer Goroutine … gestarteter Prozess“) und die Begründung der Zeile `:184` treffen `TestRunReplayTLS` wörtlich; der Satz „Die Entscheidung gilt für beide Schichten dieses Slice“ deckt ihn. Die Tabelle des Architect ist eine Aufzählung der Kandidaten des Planners, nicht abschließend („welche der Architect festlegt“). Offen bleibt die Form (F-2): Der Plan nennt die Stelle in §1/§2/§6 nicht, die Bestätigung steht aus. `play_anmeldung_test.go` und `server_extended_test.go:790` (`<-c.freigabe`, Test-Double) sind gegen `SPEC-038` vertretbar ohne Frist (Double bzw. Testgeschirr; zu `danach` siehe F-3). Das Gegenstück, das tatsächlich nicht gedeckt ist, ist `server_tls_test.go` `defer … <-fertig` (F-1).

**(c) Zusatzumfang im G1-Test.** Der Plan folgt (§1 G1 *Test-Idee*, §3 Zeile Bootstrap und §7 nennen `replay`, IP-Kontrolle und die zwei Startfehler; `be8ca1f4`, Plan folgt Korrektur, §3.9). Die Startfehler-Fälle sind keine Erweiterung der Zusagen, sondern die Aufzählung in `LH-FA-13.a` *Startphase* (Ergänzung des Architect in `44cb42c3`), und tragen je eine Mutation. Zählung: drei Liefer-Punkte, zwei Schichten; der Test bleibt in `bootstrap_test`. Kein Befund.

**(d) `0752b3ec`.** Siehe F-5: nur Test-Code, Produkt-Code unverändert, kein Widerspruch zu §1. Die Annahme, `Run` sei geändert worden, trifft nicht zu; geändert ist die Stelle im Test, die `Run` aufruft.

**(e) G7.** Als akzeptiertes Negativ nach `AGENTS.md` §3.10 korrekt geführt: Grund „keine Zusage in `SPEC-032`“ (der Absatz *Nicht zugesichert ist, was eine beendete Verbindung im Prozess zurücklässt* steht in `spec/spezifikation.md`), ausdrücklich nicht „äquivalent“, ohne behauptete Invariante, Ort oder Sonde, mit Lauf am Endstand (grün). Der Doc-Kommentar von `replayWaechter` trägt keinen Satz zum Ende der Sitzung (§3.11 gewahrt). Kein Befund.

**(f) Abdeckungszeilen.** `make abdeckung` schrieb zwei Zeilen für `TestToMessageSpaltenMetadaten` (`LH-FA-02` und `LH-FA-03`, je `Happy`) und eine für `TestRunSignalVorDemOeffnen` (`LH-FA-13`, `Boundary`); `abdeckung-gesamt.md` führt für `LH-FA-03` nun `E2E, Unit` in der ersten Spalte. Die Zeilen sagen nur zu, was `ToMessage` bzw. `Run` prüfen; die Plattformgrenze steht im Text. Hinweis zu den Negativpfaden: F-4. Dass die Zeilen die Tabellen bei `make abdeckung-check` erfüllen, ist nur aus §7 übernommen und nicht von mir gefahren.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driving/pgwire` (`server_test.go`, `server_extended_test.go`) | geprüft, ohne Befund außer F-1 (`server_tls_test.go`) |
| `internal/bootstrap` (`bootstrap_test.go`, `tls_test.go`) | geprüft, ohne Befund außer F-2 und F-3 |
| `spec/spezifikation.md` (`LH-FA-03.a`, `LH-FA-13.a`, `SPEC-032`, `SPEC-038`, Änderungs-Log) | geprüft, ohne Befund; Ergänzungen tragen Zusage und Grenze, keine ADR-, Slice- oder Wellen-Bezüge |
| `docs/user/abdeckung-*.md` | geprüft, ohne Befund außer F-4 (INFO) |
| Produkt-Code (`server.go`, `bootstrap.go` und übrige Nicht-Test-Dateien unter `internal/`) | geprüft (nicht im Diff), ohne Befund |
| Hard Rules: Core-Reinheit, `//nolint`, ADR-Immutabilität, Move + Inhalt in einem Commit | geprüft, ohne Befund (keine Moves im Diff, kein `nolint`, keine ADR berührt) |
| Kommentare der neuen Tests (§3.11) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Warten ohne Frist bleibt in der Schicht stehen · Randform erst im Code-Commit im Plan genannt · Warten auf Kanal des Testgeschirrs, Grenzfall · Deklaration enger als der Test · Größenzahl um eins abweichend

## Verdikt

**Merge-blockierend:** nein — keine HIGH- oder MEDIUM-Findings. F-1 und F-2 sind LOW ohne Failure-Szenario außerhalb einer hängenden Zeitüberschreitung; ob F-1 gegen den DoD-Punkt „Frist beim Warten“ trägt, entscheidet der Verifier. Die Bestätigung von `tls_test.go:100` (F-2) bleibt beim Architect.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen zusätzlich in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation (DoD, Läufe und Mutationen am Endstand prüft der Verifier).
