# Verifikation: slice-v1-abschluss-zeitangaben-uhr — 2026-10-12

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-11-review-slice-v1-abschluss-zeitangaben-uhr.md`, F-640 bis F-646); die Nacharbeit prüfe ich nach. Ich repariere nichts.

**Gegenstand:** Slice-Plan `slice-v1-abschluss-zeitangaben-uhr` (Kopf, §1 bis §8) gegen den Gesamt-Diff `aa07a01..7127433` (21 Dateien, 906 Zeilen hinzu, 36 weg). Implementer: `0ca3161`, `685bf54`, `cf66288`, `a866965`, `29caf1e`, `7127433`. Review: `cf69828`, `368e857`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

**Eingang:** der Plan ganz am Stand `7127433` (§7 *Belege des Implementers* nur als Beleg-Teile; §7 ist bis zur Closure unvollständig, das ist kein Befund); [`LH-FA-21.a`](../../spec/spezifikation.md#lh-fa-21a--zeitangaben) *Uhr*; [ADR-0015](../plan/adr/0015-uhr-port.md); der Review-Report; der Code `internal/hexagon/ports/driven/uhr.go`, `internal/testuhr`, `internal/testzeit`, `internal/bootstrap/systemuhr.go`, `.a-check.yml`, `tools/arch/a-check-negativ.sh`. Den Bericht des Implementers habe ich nicht gelesen. Beim Start war der Arbeitsbaum sauber, HEAD `7127433`.

**So wurden die Proben gebaut.** Alles Docker-only, Scratch im Scratchpad der Sitzung, nie im Repo. Mutanten: je Mutant eine frische Kopie des Baums ohne `.git` (`tar`, entpackt mit `-m`, keine Zeitstempel übernommen), genau eine `sed`-Ersetzung (das Skript bricht ab, wenn sie nichts ändert), `gofmt -l` im selben Lauf, `go test -count=1 -run …` im Image `pgwire-recorder:test` mit `--network none`. Mutanten der Regel: Kopie, dann `tools/arch/a-check-negativ.sh` bzw. das gepinnte a-check-Image gegen die Kopie.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Port, Fake-Uhr, Tests des Vertrags. **Bestätigt.**

Alle Tests, die der Plan nennt, existieren (`grep "func Test…("` je Name, genau eine Fundstelle). Ohne den Fix aus dem richtigen Grund rot (Wortlaut der Zusage in [`LH-FA-21.a`](../../spec/spezifikation.md#lh-fa-21a--zeitangaben) *Uhr*, jeder Mutant nur eine Zeile):

| Zusage | Mutant | rot, mit Meldung |
|---|---|---|
| Fake: weckt jeden Wartenden, dessen Zeitpunkt erreicht ist (gleich) | `<` statt `<=` im Wecken | `TestUhrWartetBisZumZeitpunkt`, `TestUhrWecktAlleErreichten`, `TestUhrBeendeterKontextGehtVorBeimWecken` (kehren nicht zurück) |
| Warten: erreichter Zeitpunkt sofort | `<` statt `<=` im Eintritt | `TestUhrErreichterZeitpunktSofort` (kehrt nicht zurück), dazu die drei obigen |
| Warten: ein beendeter Kontext geht vor, auch bei vergangenem Zeitpunkt | Vorab-Prüfung entfernt | `TestUhrBeendeterKontextVorDemWarten` (`WarteBis(-1)`: nil statt canceled) |
| Fake: Kontext beendet, wenn der Zeitpunkt zugleich erreicht wird | nach dem Wecken `nil` | `TestUhrBeendeterKontextGehtVorBeimWecken` (`nil` statt canceled) |
| Fake: negative Zahl ist Programmierfehler | Panik entfernt | `TestUhrVorrueckenNegativPanik` |
| Fake: null ist erlaubt (Grenze) | Panik schon bei `ms <= 0` | `TestUhrRueckenNurAufAnweisung` (Panik, siehe V-188) |
| Fake: rückt nur auf Anweisung vor | `Jetzt` addiert die Realzeit | `TestUhrRueckenNurAufAnweisung` (`Jetzt nach Wartezeit ohne Vorruecken: 20`) |
| Fake: rückt nur auf Anweisung vor | `WarteBis` setzt die Uhr auf den Zeitpunkt | `TestUhrWartetBisZumZeitpunkt` (kehrt bei Jetzt 4 zurück) |
| Fake: um die genannte Zahl | `jetzt = ms` statt `+=` | `TestUhrRueckenNurAufAnweisung` (3 statt 10) |
| Fake: weckt jeden | nur den ersten | `TestUhrWecktAlleErreichten` (zweiter Wartender kehrt nicht zurück) |
| Fake: weckt nur die erreichten | alle | `TestUhrWartetBisZumZeitpunkt`, `TestUhrWecktAlleErreichten` (Zeitpunkt 10 bei Jetzt 7) |
| Warten: ein Abgebrochener wartet nicht mehr | `delete` im Abbruch entfernt | `TestUhrAbbruchWaehrendDesWartens` (`Wartende nach Abbruch: 1`) |

Die Mutanten sind unabhängig vom Implementer neu gebaut, nicht aus §7 übernommen.

### Punkt 2: Systemuhr im Composition Root. **Bestätigt.**

| Zusage | Mutant | rot, mit Meldung |
|---|---|---|
| erreichter Zeitpunkt sofort (`<=`) | `<` | `TestSystemuhrErreichterZeitpunktOhneWarten` (400 Aufrufe brauchten 400 ms; F-642 behoben) |
| beendeter Kontext geht vor | Prüfung am Schleifenkopf entfernt | `TestSystemuhrBeendeterKontextVorDemWarten` (`WarteBis(-1)`: nil) |
| Abbruch beendet das Warten | `ctx.Done` ignoriert | `TestSystemuhrFernerZeitpunktDrehtNicht`, `TestSystemuhrAbbruchWaehrendDesWartens` (kehren nicht binnen 5 s zurück) |
| Fehler des Kontexts beim Abbruch | `nil` | beide obigen (nil statt Fehler) |
| abgeschnitten | gerundet | `TestSystemuhrLesungAbgeschnitten` (Lesung 5 über Obergrenze 4) |
| Ursprung bei Erzeugung / nicht Wanduhr | `UnixNano`; Ursprung eine Stunde zurück | `TestSystemuhrUrsprungBeiErzeugung`, `TestSystemuhrLesungAbgeschnitten` |
| aus der monotonen Zeit | `Round(0)` am Ursprung | `TestSystemuhrUrsprungTraegtMonotonenMesswert` |
| Warten nicht vor dem Zeitpunkt | halbe Wartezeit, `return nil` nach dem Timer | `TestSystemuhrWartenUntergrenze` (`WarteBis(1) kehrte bei Lesung 0 zurück`) |
| vergangener Zeitpunkt sofort, auch `MinInt64` | Abstand vor dem Vergleich (`zeitpunkt-jetzt <= 0`, Überlauf) | `TestSystemuhrVergangenerZeitpunktSofort` (`WarteBis(-9223372036854775808)` hängt) |
| ferner Zeitpunkt wartet, ohne zu drehen | Begrenzung der Stufe entfernt | `TestSystemuhrFernerZeitpunktDrehtNicht` (345 ms CPU-Zeit gegen Grenze 150 ms) |

Die Fake-Uhr erfüllt den Port (`var _ driven.Uhr = (*testuhr.Uhr)(nil)`, vom Compiler geprüft). `NeueSystemuhr` hat keinen Aufrufer außer den Tests; Plan §3 sagt „nicht in einen Service verdrahtet“ ausdrücklich (F-645, kein Befund).

### Punkt 3: Architekturregel, Gegenprobe, Migration der Tests. **Bestätigt.**

- `.a-check.yml`: `{pattern: "^time$", match: regex, adapter: [internal/adapters, test/, internal/testuhr, internal/testzeit, internal/testpki]}`. `make a-check` grün (0 Befunde). Der Composition Root steht nicht in der Liste; `composition_root` nimmt ihn von der Regel aus (der Kommentar in der Datei sagt das so).
- `make a-check-negativ` grün, **23 Fälle** (`grep -c '^fall '` = 23; Kopf des Skripts „Dreiundzwanzig“, Schlusszeile „dreiundzwanzig“, `harness/README.md` und `harness/mk/arch-negativ.mk` nennen die Regel in derselben Fassung). Jeder der Fälle 18 bis 23 trägt seine Beschreibung: 18 bis 20 lehnen Domain Model, Driven Ports, Services ab und a-check nennt die Datei; 21 PGWire-Adapter grün; 22 Composition Root grün (Beschreibung „von der Regel ausgenommen, kein Listeneintrag“ trifft: der Fall belegt den Zustand, nicht einen Eintrag); 23 `runtime` im Domain Model grün.
- Mutanten der Regel, mit dem Skript der Gegenprobe in Kopien: Eintrag entfernt → `model-time`, `ports-time`, `services-time` rot-erwartet, a-check grün; Muster `time` statt `^time$` → `model-runtime` rot; `match: regex` entfernt → dieselben drei Fälle rot; `internal/adapters` entfernt → alle grün-erwarteten Fälle rot, darunter `pgwire-time`.
- Jeder Eintrag der Liste trägt (der Kommentar in `.a-check.yml` sagt „Jeder Eintrag steht für Dateien, die `time` einbinden“): je einzeln entfernt meldet a-check 17 `tech-leak` unter `test/integration` (`test/`), 1 in `internal/testuhr/testuhr_test.go`, 2 in `internal/testzeit/testzeit.go`, 1 in `internal/testpki/testpki.go`, mehrere unter `internal/adapters`. Keine Gate-Lockerung: `.golangci.yml`, `.d-check.yml` unverändert, die Regel ist neu und streng ([`AGENTS.md`](../../AGENTS.md) §3.6).
- Die drei Testdateien in `internal/hexagon/services` binden `time` nicht mehr ein (`grep '"time"' internal/hexagon` leer); `_ "time"` in `play_test.go` einer Kopie lässt a-check rot werden (`play_test.go:4: app-impurity: Application importiert time`). Der Diff dieser Dateien tauscht nur die Aufrufe; jeder Wert ist gleich geblieben (30000 ms, 2000 ms, 100 ms, 50 ms, 1 ms).
- `TestWartenInMillisekunden` und `TestFristLaeuftNachMillisekundenAb` existieren; der Mutant Mikrosekunden macht den ersten rot (`Schlafe(30) dauerte 1.13ms`).

### Punkt 4 bis 9: Gate, Review, Closure, Register, Risiken, Paarungen. **Nicht Gegenstand dieses Laufs** (Closure-Pflichten; §7 vor der Closure unvollständig, kein Befund). `make gates`: siehe Verdikt.

---

## 2. Spezifikation und Sicht gegen Code

`LH-FA-21.a` *Uhr* gegen den Code, Satz für Satz:

| Zusage | Stelle | Befund |
|---|---|---|
| ganze Millisekunden, `int64`, Ursprung wählt die Uhr, kein Typ der Zeitbibliothek | `uhr.go` (`Jetzt() int64`, `WarteBis(ctx, int64)`), Datei bindet nur `context` | erfüllt |
| abgeschnitten, nie kleiner als die vorige | Systemuhr `time.Since / time.Millisecond` (monotone Zeit, ganzzahlige Division); Fake nur `+=` mit `ms >= 0` | erfüllt |
| beendeter Kontext geht vor, auch bei vergangenem Zeitpunkt; während des Wartens; sonst `nil`, erreicht/vergangen sofort | Systemuhr: Kontext am Schleifenkopf vor dem Vergleich, nach jedem Timer erneut; Fake: Vorab-Prüfung, nach dem Wecken `ctx.Err()` | erfüllt |
| Systemuhr: Ursprung bei Erzeugung, monoton | `NeueSystemuhr`: `time.Now()` mit Messwert der monotonen Zeit | erfüllt |
| Fake: rückt nur auf Anweisung, nicht negativ (Panik), weckt alle erreichten Wartenden | `testuhr.go` | erfüllt |
| Kern nutzt keine Zeitfunktion der Standardbibliothek | `.a-check.yml`, Gegenprobe, `grep` | erfüllt am Import (Grenze im Kommentar der Regel genannt, §6 Randform *Regel*) |

Die Sicht (`spec/architecture.md`) nennt die Uhr als Driven Port und die Systemuhr als Baustein des Composition Root mit der Fake-Uhr in Tests; der Code folgt dem. [ADR-0015](../plan/adr/0015-uhr-port.md) ist umgesetzt, die Folgepflicht (Architekturregel) ist der dritte Liefer-Punkt.

## 3. Flackern (Prüfpunkt 4)

Die Tests des Plans (`internal/bootstrap` `TestSystemuhr*`, `internal/testuhr`, `internal/testzeit`) liefen 25 Mal hintereinander bei 60 gleichzeitigen Rechenschleifen auf 20 Kernen: **0 Fehlschläge** in 75 Läufen, darunter `TestSystemuhrFernerZeitpunktDrehtNicht`, `TestSystemuhrAbbruchWaehrendDesWartens` (Untergrenze 30 ms) und `TestFristLaeuftNachMillisekundenAb`. Im ersten, breiteren Lauf (alle Pakete einschließlich `internal/hexagon/services`, 15 Durchgänge bei 40 Schleifen) fiel nur `TestRunPlayZweitesSignal` aus, siehe V-187.

## 4. Plan gegen Code, Hard Rules, Namensabgleich

- §1/§3/§6/§8 folgen dem Code: drei Liefer-Punkte (Port mit Fake, Systemuhr, Regel mit Migration), zwei Schichten (Kern, Bootstrap mit Logik). §3 nennt jetzt auch `export_test.go` und beide `doc.go` (F-645 (4) nachgezogen). Jeder Dateiname und jeder Name in §3 (`Neu`, `Jetzt`, `WarteBis`, `Vorruecken`, `Wartende`, `Nach`, `Schlafe`, `NeueFrist`, `Abgelaufen`, `NeueSystemuhr`, `UrsprungMonoton`) steht im Code; jeder Testname in §7 hat genau eine Fundstelle.
- Kein `//nolint` in `cmd/`, `internal/`, `test/`. Keine Kennung eines Spezifikations- oder Sichtstratums in den Commit-Messages von `aa07a01..HEAD` ([ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md) laut Hook), Kennung `slice-v1-abschluss-zeitangaben-uhr` und `LH-FA-21` in jeder.
- Handbuch und `README.md` unberührt (`git diff aa07a01 HEAD -- docs/user README.md` leer); der Slice ändert kein Verhalten des Binaries.
- `make abdeckung-check`, `make kopf-check`, `make docs-check`, `make baseline-verify`: grün.
- Nehmer: `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` nennen den Port, die Fake-Uhr und die Systemuhr mit der Kennung dieses Slice in §1 (nehmen an), aber weder `Jetzt` noch `WarteBis` (V-186, F-643 beim Planner).

## 5. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-186 | LOW | F-643 ist offen und wird hier bestätigt: §6 Risiko 3 nennt als Gegenmittel, dass die beiden Nehmer die Signatur in §1 nennen. `grep -c "Jetzt\|WarteBis"` in `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` liefert 0. Die Nehmer nehmen an (kein Ausschluss trifft), die Namen fehlen. Wer sie einträgt, zählt nach (`AGENTS.md` §3.13). | `docs/plan/planning/next/slice-v1-abschluss-zeitangaben-aufzeichnen.md` §1; `…/slice-v1-abschluss-zeitangaben.md` §1 | Planner (bereits so im Review); Risiko 3 trägt bis dahin keinen Ausgang |
| V-187 | INFO | `TestRunPlayZweitesSignal` (`internal/bootstrap`, nicht Teil dieses Diffs) fällt unter hoher CPU-Last aus („binnen 30 s keine Client-Nachricht 'X'“): 3 von 15 Durchgängen bei 40 Rechenschleifen, und ebenso am Stand `aa07a01` (2 von 12). Nicht durch diesen Slice verursacht; kein Befund gegen ihn. | `internal/bootstrap/play_test.go:274` | Beobachtung für das Register, kein Slice-Befund |
| V-188 | INFO | Der Mutant „null ist erlaubt“ (`ms <= 0` löst Panik aus) wird rot, weil `Vorruecken(0)` im Test selbst eine Panik wirft und das Testbinary abbricht; die übrigen Tests des Pakets laufen dann nicht. Rot aus dem richtigen Grund, aber ohne Meldung; `TestUhrRueckenNurAufAnweisung` fängt die Panik nicht. | `internal/testuhr/testuhr_test.go` · `u.Vorruecken(0)` | ohne Blockade |

Falsifikationsversuche, die nichts fanden: `math.MinInt64`/`math.MaxInt64` an der Systemuhr (Überlauf der Differenz, Stufen der größten Dauer); beendeter Kontext mit `MaxInt64`; Wecken und Kontextende gleichzeitig am Fake; Wartende nach Abbruch; `runtime` gegen das Muster; Regel mit jedem Listeneintrag einzeln entfernt.

---

## Verdikt

**Befunde:** 0 HIGH, 0 MEDIUM, 1 LOW (V-186), 2 INFO (V-187, V-188).

Die drei Liefer-Punkte der DoD sind **bestätigt**: Die Tests existieren, werden ohne den Fix aus dem richtigen Grund rot (neu gebaute Mutanten je Zusage im Wortlaut von [`LH-FA-21.a`](../../spec/spezifikation.md#lh-fa-21a--zeitangaben), 25 Mutanten am Code und 10 an der Regel, alle rot, keiner aus falschem Grund), und Code, Gegenprobe (23 Fälle, Beschreibungen wahr) und Dokumentation stimmen mit Plan und Spezifikation überein. Kein Flackern der neuen Tests in 75 Läufen unter Last.

**Closure-blockierend:** nichts im Code. Vor der Closure: V-186 (Planner trägt die Namen in die Nehmer ein oder gibt dem Risiko 3 einen anderen Ausgang), die Closure-Pflichten aus §2 des Plans (Review-Häkchen, Notiz mit Lerneintrag, Register, Ausgänge der drei Risiken: Risiko 1 erwartet *entfallen* mit dem Grund aus §6, Risiko 2 *entfallen* (dieselben Werte, 75 Läufe ohne Flackern), Risiko 3 bis V-186 offen).

**Finding-Klasse dieses Laufs:** keine neue; die Klasse „Zusage weiter als der Test“ (V-175, V-182) trat nicht auf.

**Aufgeräumt:** nur benannte Verzeichnisse unter `w/` im Scratchpad der Sitzung; keine Container, Netze oder Images angelegt. Der Arbeitsbaum des Repos war bis zu diesem Bericht unverändert.

**Gate:** `make gates` auf dem sauberen Baum nach dem Commit dieses Berichts: Ergebnis in der Übergabe-Nachricht (ein Gate-Lauf belegt nur den Stand, auf dem er lief).
