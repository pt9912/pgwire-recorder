# Verifikation: slice-play-zweites-signal-terminate-wettlauf

**Rolle:** Verifier (Modul 11). **Datum:** 2026-10-25. **Zweig:** `impl-wettlauf1`, Diff `dced44df..6eca756c`.
**Gegenstand:** `LH-FA-20.a` *Abbruchsignal* (Spezifikation f17d9027 und fa8a8479).
**Urteil: bestanden mit Vorbehalten** (zwei LOW, kein Befund gegen die Lieferung).

## 1. Belege des Implementers (§7)

Vorhanden und vollständig für die DoD: Mutationstabelle (Zusage, Mutation, roter Test), Last K=32 und K=48 je 20 Läufe, Gates, Größe. Der Beleg der Zeile „Last der neuen Tests“ liegt ebenfalls vor. Die Größe nach Befehl `AGENTS.md` §3.14 ist am Endstand 341 (Marke 1500; §7 nennt 252 vom ersten Stand, die Nacharbeit kam hinzu) und verletzt nichts.

## 2. Eigene Sensoren

| Sensor | Ergebnis |
|---|---|
| `go test -run Einspiel` (Adapter) und `TestRunPlayZweitesSignal` im Image `pgwire-recorder:test`, `--network none` | grün |
| `make abdeckung-check`, `make kopf-check` | Exit 0 |
| Last K=32 (`while`-Schleife, `postgres:17-alpine`, `--network none`, Last 11 bis 18): `TestRunPlayZweitesSignal` und die neuen Zeittests, 6 Läufe | 6 von 6 grün |

## 3. DoD gegen Code

- **Liefer-Punkt 1 (Zustand statt Sperre):** bestätigt. `schliesse` liest `schreibEnde` unter `zustand` im selben Schritt, in dem es `geschlossen` setzt; kein `TryLock` mehr. Das Warten ist durch `frist` (Gesamtfrist ab Aufruf, geteilt mit `SetWriteDeadline`) begrenzt. Entspricht Randform 1 bis 3.
- **Liefer-Punkt 2 (Tests mit Mutation):** bestätigt durch eigene Mutanten (Kopie des Zweigstands, Image wie oben, Original danach wiederhergestellt):

| Zusage | Mutation | Ergebnis |
|---|---|---|
| kein Warten (TryLock-ähnlich: Schreibzustand ignoriert) | `wartePunktlich(nil, …)` | rot: `panic: close of closed channel` (zweites gleichzeitiges `Write` im Fake), wie in §7 beschrieben |
| Frist begrenzt das Warten | Timer `time.Hour` | rot: 6 Tests, u. a. `TestEinspielWartePunktlich`, `…SchliesseBeimSenden` |
| Gesamtfrist | `SetWriteDeadline(time.Now().Add(s.frist))` nach dem Warten | rot: `TestEinspielSchliesseGesamtfrist` |
| Frist gerissen: kein Terminate | `case <-t.C: return true` | rot (Panic, s. o.) |
| `Anfrage` nach dem Schließen schreibt nicht | `beginneSchreiben`-Prüfung in `Anfrage` entfernt | rot: `TestEinspielAnfrageNachSchliessen` |
| Terminate nach vollständiger Anfrage bei gehaltener Bearbeitung | `endeSchreiben` 3 s nach dem `Flush` | rot: `…SendetNachVollstaendigerAnfrage`, `…Gesamtfrist`; `TestRunPlayZweitesSignal` rot nach 33 s |

  Nicht gefahren (Stichprobe): Reihenfolge der Gründe `SchliesseOhneTerminate` und „erstes Schließen“, `…VerwirftWarteschlange`, Vorab-Prüfung in `wartePunktlich`, Code `PGR-E4003`. Sie stützen sich auf §7 und die Reviews.
- **Liefer-Punkt 3 (Plattform, Abdeckung):** Spezifikation trägt „geprüft nur unter Linux“ (Randform 5); `abdeckung-check` grün.
- **Kein anderes Verhalten, Last-Punkt:** Kein Verhalten außerhalb des Abbruchsignals im Diff; Last-Beleg von mir in kleinerem Umfang bestätigt (6 Läufe K=32, nicht 20, nicht K=48). Die 20 Läufe je Last stammen allein vom Implementer.
- **`make gates`:** nicht von mir gefahren (Auftrag: läuft auf `main`); §7 belegt Exit 0 am Endstand.
- Die Punkte Review, Closure-Notiz, Register, Risiko-Ausgänge und Paarungen sind offen und gehören dem Planner bei der Closure.

## 4. Plan gegen Code

- Dateien des Diffs entsprechen §3 (Adapter, Test, `export_test.go`, Port-Kommentar, Spezifikation durch den Architect, `docs/user/abdeckung-unit.md` durch `make abdeckung`). `play_test.go` unverändert gefahren.
- Die Spezifikationsänderung (fa8a8479, Randform 8) ist dem Plan zufolge Architect-Entscheidung auf Rückgabe, nicht Implementer-Entscheidung; §1 und §6 Punkt 8 nennen sie.

## 5. Vorbehalte

- **V-1 (LOW):** §3, Zeile `einspielen.go` führt noch „Zustand ‚teilweise geschrieben‘ … gesetzt für die Dauer des Schreibens einer Nachricht“; der Code und Randform 2 führen „schreibt“ je `Flush`. Dasselbe Finding F-2 aus Review r2 blieb im Plan stehen (`AGENTS.md` §3.9, `BEO-REPO/plan-folgt-korrektur-nicht`). Der Planner zieht es bei der Closure nach.
- **V-2 (LOW):** Der Teilsatz „ohne Meldung“ (Randform 8) hat keinen Test; §7 nennt ihn als Grenze mit Adresse Planner. Er sagt in der Spezifikation mehr zu, als ein Test prüft (`AGENTS.md` §3.11); der Planner entscheidet bei der Closure über engere Fassung oder Folge.
- Zeittests tragen Spielraum von rund 800 ms (Frist 1 s); 0 rot unter Last belegt ihn nicht, er folgt aus der Konstruktion (`BEO-REPO/test-faellt-unter-cpu-last-aus`, weiter offen).

## 6. Grenzen dieser Verifikation

Keine Gegenprobe (kein Gate berührt), kein `make gates`, kein K=48, keine 20 Läufe; macOS und Windows nicht geprüft (Grenze in der Spezifikation).
