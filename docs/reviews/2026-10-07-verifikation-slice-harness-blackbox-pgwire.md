# Verifikation: slice-harness-blackbox-pgwire — 2026-10-07

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft. F-450, F-451 und F-453 hat er an mich übergeben, F-448 und F-449 an den Implementer.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-blackbox-pgwire.md` (Kopf, §1 bis §8) gegen den Gesamt-Diff `6ff1b3b..1b6fa29`. Darin liegen `b37d6e7` (Umbau, Belege in §7), `7a951bd` ([Review](2026-10-07-review-slice-harness-blackbox-pgwire.md), F-448 bis F-453) und `1b6fa29` (Nacharbeit F-448 und F-449, nur Plan).

**Eingang:**

- DoD-Liefer-Punkte, §1, §3, §6 *Randformen* und *Risiken*, §7 *Belege des Implementers*, Commit-Messages
- `spec/spezifikation.md` `SPEC-049` Punkt 7 und 8, der Absatz *Abbruch* im Record-Teil (Session-Ende, blockiertes Schreiben); [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 und 5; [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md); [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- `.claude/commands/implement-slice.md` Schritt 19 (Einordnung grüner Mutanten)
- Produkt-Code `internal/adapters/driving/pgwire/server.go` (`handle`, `recordSitzung`, `richtungen`, `wecke`, `weiterlesen`, `beende`, `clientRichtung`, `fail`, `note`, `zielart`, `toClientMessage`, `toMessage`)
- der Review-Report und die Pläne `slice-lint-bestand-driving`, `slice-v1-abschluss-betrieb`, `slice-harness-mutation` (für die Adresse zu F-451)

Bei meinem Start war der Arbeitsbaum sauber, HEAD `1b6fa29`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

**Wie die Proben gebaut wurden.** `6ff1b3b` (vor dem Umbau, `b37d6e7^`) und `1b6fa29` habe ich per `git archive` in je eine Kopie im Scratchpad entpackt. Testlisten, Testläufe, Sonde und Mutationen liefen im Image der Stufe `deps` (eigener Tag, aus dem Arbeitsbaum gebaut), mit Bind-Mount der Kopie und `--network=none`. Jede Mutation lief in einer frischen Kopie (`cp -r` ohne `-p`). Das Ersetzungsskript bricht ab, wenn der Suchtext nicht genau einmal trifft, und prüft per Prüfsumme, dass sich die Datei geändert hat. Die Lint-Läufe auf den Kopien liefen mit `docker build --progress=plain --target lint`. Das Repo selbst blieb unverändert. `make lint` und `make gates` liefen im Arbeitsbaum.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `_test`-Pakete, Testliste, Abdeckung, Brücke. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Testdateien gehören zu `_test`-Paketen | An `1b6fa29` sind `server_test.go`, `server_extended_test.go` und `server_meldung_test.go` `package pgwire_test`. Im Paket des Codes liegt nur `export_test.go`. | bestätigt |
| Testliste vorher und nachher gleich | `go test -list .` an `6ff1b3b` und `1b6fa29`: je 40 Tests, `diff` leer. `go test -count=1 -v` liefert an beiden Ständen dieselben 53 Zeilen `--- PASS` (sortiert, ohne Zeitangabe), kein `FAIL`. | bestätigt |
| keine Prüfung entfallen | bis auf den Rückgabewert von `weiterlesen` ja, siehe Abschnitt 3 und 4 | bestätigt mit eingeordneter Ausnahme |
| Abdeckungs-Deklarationen unverändert | An beiden Ständen 25 Deklarationen; die vollständigen Kommentarblöcke ab `// Abdeckung:` (92 Zeilen, je Datei sortiert) sind gleich, ebenso die Zuordnung Deklaration → nächstes `func Test…`. | bestätigt |
| `make abdeckung-check` grün | Exit 0 einzeln und im Gate-Lauf (Abschnitt 6) | bestätigt |
| Unexportiertes nur über die Brücke | Der Compiler erzwingt es für `pgwire_test`. Die Brücke selbst: siehe Abschnitt 2. | bestätigt |

### Punkt 2: `make lint` ohne Befund in den Testdateien. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| 6 Befunde vorher in den Testdateien | Kopie `6ff1b3b`: `60 issues:`, 20 Zeilen unter `internal/adapters/driving/pgwire`, davon 6 in `*_test.go` (`testpackage` 3, `gochecknoglobals` 1 für `sendeMu`, `revive` `context-as-argument` 2). Wort für Wort wie der Block *Vorher* in §7. | bestätigt |
| 0 nachher | Kopie `1b6fa29` und `make lint` im Arbeitsbaum (Exit 2): `54 issues:`, unter dem Pfad 14 Zeilen, alle in `server.go`. Sortiert identisch mit dem Block *Nachher* in §7. Keine `lint:`-Zeile der eigenen Prüfungen, auch nicht an `export_test.go`. | bestätigt |
| modulweit 60 → 54, kein neuer Befund | Befundzeilen ohne Zeile und Spalte, sortiert: `comm -13` (nur nachher) leer, `comm -23` (nur vorher) genau die 6 Zeilen der Testdateien. Die 14 Zeilen in `server.go` sind an beiden Ständen mit Zeile und Spalte gleich. | bestätigt |
| Zeilen vorher und nachher im Bericht | §7 nennt vorher die 6 Testdatei-Zeilen und nachher die 14 Produkt-Zeilen mit dem Satz, sie seien dieselben wie vorher. 6 + 14 = 20, die Gleichheit habe ich gemessen. | bestätigt |
| ohne `//nolint` | `grep -rn nolint` über `internal/`, `cmd/`, `test/` an `1b6fa29` leer | bestätigt |
| ohne Änderung an `.golangci.yml` | `git diff 6ff1b3b 1b6fa29 -- .golangci.yml tools/harness/lint.sh` leer | bestätigt |
| kein `_ =` vor einem gemeldeten Fehler | Alle Blank-Zuweisungen in den drei Testdateien an beiden Ständen verglichen (Paketpräfix entfernt, sortiert): je 44. Unterschiede nur in der Parameterreihenfolge (`ctx` zuerst) und `toClientMessage` → `ToClientMessage`. `_ = fe.Flush()` steht vorher in `sende`, nachher in `nebenher`. Der neue `TestWeiterlesenNachWecken` prüft den Fehler von `SetDeadline`. | bestätigt |

### Punkt 3: `make gates` grün. **Bestätigt.**

`make gates` im Arbeitsbaum an `1b6fa29` (Go-Code gleich `b37d6e7`): siehe Abschnitt 6. §7 trägt für diesen Punkt nur die Zeile „`make gates` grün am Stand dieses Commits“, ohne Ausgabe; den Lauf habe ich deshalb selbst gefahren.

### Kein Produkt-Code geändert. **Bestätigt.**

`git diff --stat 6ff1b3b 1b6fa29`: geändert sind die drei Testdateien, neu ist `export_test.go`, dazu Plan und Review-Report. `server.go`, `errorfields.go` und `doc.go` sind unverändert.

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| Review mit Report | liegt vor (F-448 bis F-453) |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen: §7 trägt Platzhalter, §6 *Risiken* hat drei Ausgänge „— (bei Closure)“. Das ist Sache der Closure; Vorschläge in Abschnitt 7. |

---

## 2. Brücke und White-Box-Zugriffe

**Brücke nach `SPEC-049` Punkt 7.** `export_test.go` ist `package pgwire` und die einzige Testdatei dort. Sie hat fünf Konstanten (`CodeCancelRequest`, `CodeSSLRequest`, `CodeGSSEncRequest`, `MajorSpezial`, `MeldeFrist`) und fünf Funktionen (`Handle`, `Fail`, `Note`, `ToClientMessage`, `ToMessage`), die je genau eine unexportierte Funktion oder Methode aufrufen und die Argumente unverändert weiterreichen. Keine Variable, kein Typ, kein Alias, keine Funktion `Test…`, kein Wert eines unexportierten Typs. `handle` ruft auch das Produkt (`Serve`), `richtungen` legt nur `recordSitzung` an. Hält.

**Messung der White-Box-Zugriffe gegengeprüft, auf anderem Weg als der Implementer.** In einer Kopie von `6ff1b3b` habe ich die drei Testdateien auf `package pgwire_test` mit Punkt-Import des Pakets umgestellt und mit `go test -gcflags=-e` übersetzt. Jede Fehlermeldung ist ein Zugriff auf Unexportiertes:

| Meldung | Anzahl | §7 |
|---|---|---|
| `s.handle undefined` | 10 | 10 |
| `undefined: toMessage`, `undefined: toClientMessage` | 4, 4 | 4, 4 |
| `unexported field log in struct literal` | 4 | 4 |
| `s.note undefined`, `s.fail undefined` | 3, 2 | 3, 2 |
| `undefined: meldeFrist` | 2 | 2 |
| `undefined: codeSSLRequest`, `codeGSSEncRequest`, `codeCancelRequest`, `majorSpezial` | je 1 | je 1 |
| `undefined: richtungen` | 1 | 1 (Felder und Methoden folgen daraus und werden vom Compiler nicht einzeln gemeldet) |

Die Zahlen stimmen mit der Tabelle in §7 überein.

---

## 3. F-448 und F-449: Nacharbeit im Plan

**F-448: §6 ist jetzt widerspruchsfrei und folgt dem Code.** Der Punkt *White-Box-Zugriffe im Bestand* nennt die Messung am Stand `6ff1b3b`, verweist für die Zuordnung auf §7 und sagt: `meldeFrist` reicht die Brücke als Konstante weiter, gelesen, nicht gesetzt. Das stimmt mit dem Punkt *Export-Test-Brücke* („Zustand nur an einem übergebenen Wert, nie auf Paketebene“) und mit dem Code überein: `MeldeFrist` ist eine Konstante, kein Test setzt eine Frist, `TestFehlerantwortMitFrist` liest sie nur. Dass `richtungen` nur das Produkt anlegt und `wecke`/`weiterlesen` über `Handle` an einer vom Test gestellten Verbindung geprüft werden, stimmt mit `TestWeiterlesenNachWecken` überein. Eine offene Randform nennt §6 nicht mehr. Die übrigen Punkte von §6 (Fakes der Ports, Übrige Befunde, Eigener Schnitt) widersprechen dem Code nicht. Ein Rest steht in V-86.

**F-449: Die Einordnung nach Schritt 19 trägt.** Selbst nachgefahren:

| ID | Mutation | an `6ff1b3b` | an `1b6fa29` (alle Tests des Pakets) |
|---|---|---|---|
| G1 | `weiterlesen`: nach dem Wecken `return true` → `return false` | rot, `TestWeiterlesenNachWecken` „Wecken nicht gemeldet“ | grün |
| G2 | `weiterlesen`: ohne Signal `return false` → `return true` | rot, `TestWeiterlesenNachWecken` „Wecken ohne Signal gemeldet“ | grün |

Beide Mutanten sind im Produkt äquivalent: `weiterlesen` hat genau einen Aufrufer, `clientRichtung` (`server.go` · „r.weiterlesen()“), der den Wert als Anweisung verwirft. Lesefrist und Lesen hängen nur am Signal in `weck`, nicht an der Rückgabe. Nach Schritt 19 fängt einen äquivalenten Mutanten kein Test, und die Grenze in §6 trägt ihn; die steht jetzt im Punkt *White-Box-Zugriffe im Bestand* („Liest das Produkt den Wert künftig, braucht er einen Test über die Schnittstelle“). Nach `SPEC-049` Punkt 7, letzter Fall, ist das, was nur an einem selbst angelegten `richtungen`-Wert zu sehen war, über die Schnittstelle zu prüfen; der Wegfall ist zulässig. Die Grenze hat einen absehbaren Auslöser, siehe V-86.

---

## 4. Mutationen nach `AGENTS.md` §3.10, Stichprobe gegen §7

Der Grundlauf des Pakets an `1b6fa29` ist grün. Jede Mutation lief mit `go test -count=1 -run` auf die genannten Tests.

| ID | Zeile in §7 | Mutation (wie ausgeführt) | Ergebnis | Grund in der Ausgabe |
|---|---|---|---|---|
| V1 | Frist bleibt nach dem Wecken stehen | `weiterlesen` setzt nach dem Signal die Frist zurück | rot, `TestWeiterlesenNachWecken` | `"frist abgelaufen zurück"` |
| V2 | Signal wird nicht verbraucht | `select` in `weiterlesen` entfernt, setzt immer zurück | rot, `TestWeiterlesenNachWecken` | `"frist abgelaufen zurück"` |
| V3 | Wecken hinterlegt ein Signal | in `wecke` das Senden in `weck` entfernt | rot, `TestWeiterlesenNachWecken` | `"frist abgelaufen zurück"` |
| V4 | ohne Wecken wird zurückgesetzt | `weiterlesen` setzt ohne Signal nicht zurück | rot, `TestWeiterlesenNachWecken` | „Lesefrist nicht zurückgesetzt“ nach 2 s |
| V5 | Parameterwerte kopiert | `Bind`: `append([]byte{}, p...)` → `p` | rot, `TestToClientMessage` | „Parameterwert nicht kopiert“ |
| V6 | fremde Zielart ist ein Fehler | `zielart` `default` liefert `TargetStatement` | rot, `TestToClientMessage` („Close mit Zielart X angenommen“) und `TestExtendedNichtUnterstuetzt/Zielart` (Antwort bleibt aus, Timeout) | wie §7 |
| V7 | unbekannter Antworttyp nicht erfunden | `toMessage` `default` liefert `NoData` | rot, `TestToMessageErfindetNichts` | „unbekannter Typ angenommen“ |
| V8 | ParameterDescription trägt die Typen | `ParameterOIDs: nil` | rot, `TestToMessageExtended` | `ParameterOIDs:[]uint32(nil)` |
| V9 | SQLSTATE nach Klasse | Klasse 6 liefert `XX000` | rot, `TestFehlerantwortJeKlasse`, `TestFailGleichrangig` | `Code:XX000`, erwartet `0A000` |
| V10 | genau eine ErrorResponse | `fail` sendet sie zweimal | rot, `TestFailGleichrangig` | „2 Nachrichten statt genau einer ErrorResponse“ |
| V11 | der erste Fehler zählt | `note` überschreibt `firstCode` immer | rot, `TestErsterFehlerZaehlt` | „erster Fehler: "PGR-E6001"“ |
| V12 | je Meldung eine Log-Zeile | `note` loggt nur die erste | rot, `TestNoteGleichrangig` | Log nur mit `PGR-E4003` |
| V13 | Fehlerantwort höchstens `MeldeFrist` lang | `SetWriteDeadline` in `beende` entfernt | rot, `TestFehlerantwortMitFrist` | „Sitzung kehrt nicht binnen 2s zurück“ |
| V14 | SSLRequest erkannt | `codeSSLRequest = 80877199` | rot, `TestSSLUndGSSMitN` | „Code 80877199: Antwort "E"“ |
| V15 | Felder an den Use Case | `Parse` trägt als `Statement` die Anfrage | rot, `TestExtendedSyncGruppe` | „Client-Nachrichten: …“ |
| G1, G2 | *Ohne roten Test* | Abschnitt 3 | wie §7 | äquivalent |
| G3 | *Ohne roten Test* | `meldeFrist = 3 * time.Second` | grün an `1b6fa29` **und an `6ff1b3b`** | Abschnitt 5 |

Alle 15 Zeilen der Tabelle in §7 habe ich nachgefahren; jede ist rot aus dem genannten Grund. Die drei Einträge unter *Ohne roten Test* habe ich für G1 bis G3 nachgestellt; den Mutanten `nebenher` ohne Sperre hat das Review grün nachgestellt (R6), ich nicht. Mit R2 bis R4 aus dem Review sehe ich keine Prüfung, die im Umbau verloren ging, außer dem eingeordneten Rückgabewert. F-453 ist bestätigt.

---

## 5. F-450 und F-451

**F-450: `TestWeiterlesenNachWecken` ist deterministisch genug.** Lesend:

- Nach `cancel()` ruft der Wächter in `recordSitzung` `wecke` (Signal, `frist`). Die Client-Richtung scheitert am Lesen (`abgelaufen`), `weiterlesen` verbraucht das Signal und lässt die Frist stehen. `Shutdown` des Fakes liefert ohne Einstellung `false`, also folgt ein zweites Lesen, das wieder scheitert (`abgelaufen`). Danach setzt `weiterlesen` zurück (`zurück`), und das Lesen blockiert. Steht die Client-Richtung bei `cancel()` noch nicht im Lesen, ergibt sich dieselbe Folge, weil die Frist schon in der Vergangenheit liegt.
- Andere Quellen von Ereignissen gibt es nicht: Der Test schickt keine Server-Antworten, also weckt `schreibe` nicht.
- Die einzige verdeckte Abhängigkeit ist die, die das Review benennt: Der Startphasen-Wächter in `handle` muss `fertig` gesehen haben, bevor `ctx` endet. `close(fertig)` macht ihn lauffähig; dazwischen liegen `AuthenticationOk`, die Aufbau-Antworten, der Startup des Clients, das `Sync` und das Warten auf `c:sync`, jeweils mit blockierenden Pipe-Rundläufen. Bei `GOMAXPROCS=1` kommt der freigegebene Goroutine-Wächter am nächsten Blockierpunkt dran.

Gedrosselt in einer eigenen Kopie gefahren:

| Lauf | Ergebnis |
|---|---|
| `--cpus=0.2`, `GOMAXPROCS=1`, `-count=500`, nur dieser Test | grün (500 von 500, 31 s) |
| `--cpus=0.5`, `GOMAXPROCS=8`, `-count=300 -cpu 1,2,8`, nur dieser Test | grün (900 von 900, 56 s) |
| `--cpus=0.3`, `-count=10 -cpu 1,4`, ganzes Paket | grün |

Zusammen mit den 600 Läufen des Reviews sind das 2000 grüne Läufe dieses Tests unter Drossel. Kein Befund; den Eingriff R7 des Reviews habe ich nicht wiederholt.

**F-451: Die Lücke bestand vor diesem Slice.** Belege:

- G3 (`meldeFrist = 3 * time.Second`) ist auch an `6ff1b3b` grün. Der alte Test las die Konstante ebenso, die Prüfung hing nie an einem festen Wert.
- `spec/spezifikation.md` nennt den Wert nicht. Der Absatz *Abbruch* im Record-Teil sagt nur: „Endet eine Session, schließt der Recorder die Client-Verbindung; ein blockiertes Schreiben an den Client endet damit.“ Die Default-Tabelle führt `--shutdown-timeout` (`SPEC-046`), keine Schreibfrist für die Fehlerantwort. Das Lastenheft nennt sie ebenfalls nicht.
- Eingeführt hat die Frist `09cc0e9` in `slice-extended-query-record`. Die Verifikation dort hat die fehlende Prüfung als V-22 gemeldet; daraufhin kam `TestFehlerantwortMitFrist`, der die Schranke prüft, nicht den Wert.

Der Slice selbst wächst nicht: §1 schließt geänderte Erwartungen aus, und ein Literal im Test würde einen nicht spezifizierten Wert festschreiben. Zur Adresse siehe V-87.

---

## 6. Gate-Ergebnis

- `make gates` an `1b6fa29` im Arbeitsbaum: **grün** (Exit 0). Im Log grün unter anderem `abdeckung-gegenprobe`, `a-check-negativ`, `commit-msg-gegenprobe`, `run-integration-tests` und `kopf-check-gegenprobe`, kein `--- FAIL`.
- `make abdeckung-check` und `make kopf-check` einzeln: Exit 0.
- `make lint` im Arbeitsbaum: Exit 2, `54 issues:`, keiner in einer Testdatei unter `internal/adapters/driving/pgwire`, keiner neu gegenüber `6ff1b3b`.
- `gofmt -l` des Pakets leer, `go vet` grün (Kopie `1b6fa29`).

---

## 7. Plan gegen Code

- **Kopf:** `Bezug` (`LH-QA-07`, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)) und `Berührte Spec-Stellen` (`SPEC-049` Punkt 7 und 8, unverändert) stimmen; der Slice ändert keine Spec-Stelle. `make kopf-check` ist grün.
- **§1:** Das Ziel ist erreicht (Abschnitt 1). Die Zerlegung „`testpackage` 3, dazu `gochecknoglobals` und `revive`, zusammen 3“ stimmt (1 + 2). Die Abgrenzung hält: kein Produkt-Code, kein Export, `.golangci.yml` und Lint-Gegenprobe unverändert, keine neue Testfunktion.
- **§3:** Jede Zeile hat ihre Datei im Diff und umgekehrt. `export_test.go` gibt es einmal, nur in `pgwire`.
- **§4:** Keine Rückführung ausgelöst: Umgeschrieben ist ein Test, der Rest ist umgestellt; die Brücke reichte ohne Entscheidung des Architect.
- **§5:** `make gates` grün und `make lint` ohne Befund in den Testdateien sind gemessen.
- **§6 *Randformen*:** nach `1b6fa29` ohne offene Frage (Abschnitt 3). `b37d6e7` entscheidet keine Randform, die §6 nicht nennt.
- **§6 *Risiken*:** *Prüfung geht im Umbau verloren* ist bis auf den eingeordneten, äquivalenten Rückgabewert nicht eingetreten (Abschnitt 3 und 4). *Befund verdeckt statt behoben* ist nicht eingetreten: kein neues `_ =`, `nebenher` legt je Aufruf eine eigene Sperre an (21 Aufrufe, je einer pro Test oder Teiltest und Frontend, kein `t.Parallel`). *Abdeckung sinkt* bleibt bei `slice-harness-coverage`. Die Ausgänge setzt die Closure.
- **§7 *Belege*:** vollständig für das, was die DoD verlangt (Testliste, Deklarationen, `make abdeckung-check`, Lint-Zeilen vorher und nachher, Mutationen); `make gates` nur als Zeile mit Stand, von mir selbst gefahren. Die Zahlen stimmen.
- **§8:** Die Sichtung nennt den Stand `175fd2d`; die Treffer sind plausibel, keine Schwelle erreicht.

---

## 8. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-86 | Hinweis (Closure, Grenze mit Auslöser) | §6 *Mutationstests auf unexportierte Teile* sagt weiter „muss auch nach dem Umbau fangen; sonst fehlt eine Prüfung“, ohne die Ausnahme. G1 und G2 fangen nicht mehr; getragen ist das durch die Grenze im Punkt *White-Box-Zugriffe im Bestand* und durch §7. Das ist kein Widerspruch zum Code, aber die Closure muss beim Risiko *Prüfung geht im Umbau verloren* den Rückgabewert ausdrücklich nennen, sonst liest der Ausgang wie „nichts verloren“. Vorschlag: Ausgang *entfallen*, Begründung „einzige weggefallene Prüfung ist die Rückgabe von `weiterlesen`, die das Produkt nicht liest (G1, G2 äquivalent)“. Die Grenze hat einen absehbaren Auslöser: `slice-lint-bestand-driving` zerlegt `(*richtungen).clientRichtung` (`gocognit`, `cyclop`) und zieht den Kontext aus `richtungen`; ein Umbau dort, der die Rückgabe auswertet, macht G1/G2 zu Verhaltensmutanten ohne Test. Vorschlag: In der Closure unter *Folge-Slices* einen Satz an `slice-lint-bestand-driving` §6 geben (die Grenze und G1/G2 als Mutationen nach dem Umbau). | Plan §6 *Mutationstests auf unexportierte Teile*; `internal/adapters/driving/pgwire/server.go` · „r.weiterlesen()“ | Abschnitt 3, G1/G2 |
| V-87 | Hinweis (benannte Spec-Lücke, Adresse) | F-451: Der Wert der Schreibfrist für die Fehlerantwort beim Session-Ende ist eine Lücke in der Spezifikation, die vor diesem Slice bestand (Abschnitt 5). Adresse: **nicht** `slice-lint-bestand-driving`. Sein §1 schließt Verhaltensänderung, geänderte Erwartungen in Tests und Änderungen an der Spezifikation aus; er nähme die Sendung nicht an. **Nicht** `slice-harness-mutation`. Sein §1 schließt „Neue Tests für Mutanten, die heute überleben“ aus. Vorschlag: `slice-v1-abschluss-betrieb`. Er spezifiziert und testet ohnehin das Ende von Sessions unter einer Frist (`--shutdown-timeout`, Zwangsende, `LH-FA-13.a`), und die Schreibfrist der Fehlerantwort verbraucht einen Teil genau dieses Budgets. Der Architect entscheidet den Wert vorher in der Spezifikation (Randform nach `AGENTS.md` §3.12); die Test-Idee „Schranke als Literal“ kommt dann mit. Dazu in der Closure dieses Slice die Zeile unter *Folge-Slices* und ein Satz „Übernommen aus `slice-harness-blackbox-pgwire`“ in §1 von `slice-v1-abschluss-betrieb`, durch den Planner. Lerneintrag in der Form *benannte Spec-Lücke*. Keinen weiteren Beleg in `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` oder `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: Dort steht dieser Fund schon aus `slice-extended-query-record` (V-22); ein zweiter Beleg für denselben Fund würde doppelt zählen. Der heutige Kommentar an `meldeFrist` nennt keinen Wert und sagt nur zu, was `TestFehlerantwortMitFrist` prüft. | `internal/adapters/driving/pgwire/server.go` · „const meldeFrist = time.Second“; `spec/spezifikation.md` Absatz *Abbruch* | Abschnitt 5, G3 an beiden Ständen |

---

## 9. Urteil

Die drei Liefer-Punkte sind erfüllt und selbst belegt:

- Die Tests laufen als `pgwire_test`, die Testliste ist gleich (40 Tests, 53 PASS-Zeilen an beiden Ständen).
- Die 25 Abdeckungs-Deklarationen sind gleich, Wort für Wort und in der Zuordnung.
- `make lint`: in den Testdateien 6 → 0, modulweit 60 → 54, kein neuer Befund, kein `//nolint`, kein neues `_ =`, `.golangci.yml` unverändert.
- Kein Produkt-Code geändert. Die Brücke hält `SPEC-049` Punkt 7; die White-Box-Zugriffe aus §7 sind unabhängig nachgemessen.
- Die 15 Mutationen aus §7 sind rot aus dem genannten Grund; die grünen Mutanten sind nach Schritt 19 tragfähig eingeordnet.
- F-448 ist behoben, F-449 trägt, F-450 ohne Befund, F-451 bestätigt als Lücke aus dem Bestand.

Kein Befund blockiert die Closure. V-86 und V-87 sind Vorschläge für Risiko-Ausgang, Folge-Slices und Lerneintrag.
