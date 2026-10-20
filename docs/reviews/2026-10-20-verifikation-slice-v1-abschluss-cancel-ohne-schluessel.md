# Verifikationsbericht: slice-v1-abschluss-cancel-ohne-schluessel — 2026-10-20

**Rolle:** Verifier (Modul 11), frischer Kontext. **Frage:** Bauen wir es richtig, gegen Plan, DoD und Spezifikation?

**Gegenstand:** Diff `f7aeea92..impl-cancel2` (Kopf `d8a72c11`), Zweig `ver-cancel`. Der Review-Report ist im Diff enthalten.
**Maßstab:** Slice-Plan (§1, §2 DoD, §3, §6, §7), Spezifikation `LH-FA-05.e` *Länge des `CancelRequest`*, `LH-FA-13`, Review-Report vom selben Tag.

## Ergebnis in einem Satz

DoD-Liefer-Punkte 1 und 2 und das Gate sind am Endstand belegt; F-1 bis F-4 des Reviews sind behoben (F-4 mit einer benannten Grenze im `replay`). Befunde: ein MEDIUM als Rückgabe (Spezifikation sagt zur laufenden Interaktion mehr, als der Test im `replay` prüft), ein einmaliges, nicht reproduziertes Flackern, zwei LOW. Kein HIGH.

## Selbst gefahrene Sensoren (nacheinander, nie parallel zu anderen Läufen)

| Sensor | Stand | Ergebnis |
|---|---|---|
| `PGR_LAUF=pgr-it-ver PGR_TESTS_1=TestE2ERecordCancelOhneSchluessel make test-integration` | `d8a72c11` | `--- PASS: TestE2ERecordCancelOhneSchluessel (0,73 s)`. Danach rot in der Phase ohne PostgreSQL (`TestE2EOhnePostgresEnde` u. a.: `/daten/…` fehlt): Folge der Einschränkung auf einen Test, weil die erste Phase ihre Sichten nicht schrieb; kein Befund zum Slice (der volle Lauf steckt in `make gates`, grün). |
| `make gates` (ein Lauf) | `d8a72c11`, sauberer Baum | Exit 0; Abdeckung der Unit-Tests 94,41 % bei Schwelle 94 %; `test-integration`, `test-postgres-versionen` (alle Versionen), `test-container`, `lint`, `abdeckung-check`, `workflows-check` grün. Der Baum blieb danach sauber (`git status` leer). Gegenproben nicht gefahren: Der Diff ändert kein Gate und keinen Vertrag eines Werkzeugs (`AGENTS.md` §6 Schritt 6). |
| Unit-Tests des Pakets `pgwire` (`go test` im Image `pgwire-recorder:test`, `internal` als Bind-Mount einer Kopie) | `d8a72c11` | grün; unter Last (12 Busy-Loops), `-count=5` des Pakets mit `TestCancelRequest|TestProtokollrand|TestTLS`: grün (64 s); `TestCancelRequestLaenge` allein `-count=40`: 15 360 Unterfälle grün. |

Die Belege aus §7 des Plans, die die DoD verlangt (Mutationen je Zusage, Probe-Tabelle, Läufe, Metrik), stehen dort; keiner fehlt.

## Von Hand gebrochen (Mutation einer Kopie von `server.go`, Original nie verändert)

| Mutant | Wirkung | Ergebnis |
|---|---|---|
| K0: Zweig `code == codeCancelRequest && !cancelLaengeZugesagt(laenge)` aus `startkopf` entfernt (Prüfung im Kopf) | Bibliothek liest den Rest zuerst | rot: 108 Unterfälle `TestCancelRequestKopf/…/Länge_<n>_nur_Kopf,_Verbindung_offen` (je 12 für 9, 10, 11, 13, 14, 15, 269, 270, 10000), Länge 8 grün (äquivalent, siehe unten) |
| K9, K13, K269, K10000: Länge in `cancelLaengeZugesagt` zugelassen (je eine Länge) | nur diese Länge geht an die Bibliothek | rot: genau 12 Unterfälle der eigenen Länge (K13 einmal zusätzlich 1 Unterfall `TestCancelRequestLaenge/…/Länge_10000`, siehe F-V3) |
| K8: Länge 8 zugelassen | Bibliothek lehnt sofort ab | grün (`rc=0`). Der Plan führt es als äquivalent mit Invariante (Länge 8 hat keinen Rest), Ort (Nachrichtenformat, `ReceiveStartupMessage` liest 0 Bytes) und Sonde (Fall `Länge_8_nur_Kopf,_Verbindung_offen` grün unter K8). Nach `AGENTS.md` §3.10 (a) bis (c) vollständig; die Invariante ist nicht die Zusage selbst. Von mir bestätigt, nicht widerlegt. |
| W12: vor dem Zweig der Länge 12 wird Länge 12 aus dem Kopf mit `PGR-E6001` abgelehnt | Zusage „12 ist `PGR-W3001`“ | rot: 67 Unterfälle (`TestCancelRequestLaenge/…/Länge_12*`, `…Kopf/…/Länge_12,_0_und_3_Bytes_Rest`, `…ZweiteVerbindung…`, `…LaufendeInteraktion…/Länge_12`) |
| W16: 16 aus `cancelLaengeZugesagt` genommen | Bibliothek-Länge aus dem Kopf abgelehnt | rot: 58 (u. a. `TestTLSRangfolge/CancelRequest`) |
| W268: 268 genommen (Obergrenze 267) | | rot: 38 |
| K12 (mein Textersatz „12 aus `cancelLaengeZugesagt` genommen“) | | grün, weil der Zweig der Länge 12 vor dieser Prüfung steht: äquivalent, nicht der K12 des Plans; der K12 des Plans (Länge 12 sofort `PGR-E6001`) entspricht W12 und ist rot |

F-1 und F-2 sind damit behoben und gepinnt: die Reihenfolge Länge vor Rest (Zeitpunkt der Ablehnung) färbt K0 rot; die Grenze „kein Antwort ohne Rest bei 12, 16, 17, 268 vor dem Ende der Verbindung“ färbt W12 rot. F-3: Plan §3 nennt `startkopf`, `cancelOhneSchluessel` und den Grund (Komplexitätsgrenze) — behoben.

## Plan-vs-Code, Spezifikation

- Code entspricht §3 und `LH-FA-05.e`: Reihenfolge der Gründe `startkopf`: Länge außerhalb 8 bis 10000 (`PGR-W3003`) vor Zweig 12 vor Prüfung 8 bis 11, 13 bis 15, 269 bis 10000 (`PGR-E6001` mit `ErrorResponse`, ohne Rest) vor Bibliothek (16 bis 268); fehlender Rest bei 12 liest `cancelOhneSchluessel` (`PGR-E6001`). TLS-Pflicht ohne `PGR-E6003`: Fall „unverschlüsselt bei TLS-Pflicht“ vorhanden (M4 laut §7).
- Spezifikation trägt die Nacharbeit (Zeitpunkt der Ablehnung aus dem Kopf, Grenze ohne Frist, Änderungsprotokoll-Zeile). Plan §1, §3, §6 und §8 folgen (`AGENTS.md` §3.9). Randformen in §6 stehen vor dem Code in der Spezifikation (Commit `9255aa9d` vor dem Slice-Start); die Nacharbeit ergänzte sie nach Entscheidung des Architect, nicht durch den Implementer.
- Handbuch und README gegen die Probe-Tabelle in §7 gelesen (Quantoren „jede“, „keine“, „andere Verbindungen laufen unberührt weiter“, „ohne Fehler und ohne den Exit-Code zu ändern“, „jede andere Länge von 8 bis 10000“, „bedient es nicht“): jede Aussage hat eine Probe oder einen Unit-Fall mit roter Mutation (§7). Ist-Zustand, keine Chronik, im Handbuch keine Verweise auf Spezifikation, Slice oder Review. Kein Beispiel als Datei im Diff.
- Abdeckungs-Deklarationen: `abdeckung-check` in `make gates` grün; die Texte der Deklarationen stimmen mit den Testkommentaren überein. `TestCancelRequestZweiteVerbindungUnberuehrt` sagt nur noch „sendet nach dem CancelRequest ihre Anfrage“ (F-4-Korrektur), `TestCancelRequestLaufendeInteraktionUnberuehrt` sagt „in record“.
- Größenmetrik (`AGENTS.md` §3.14): `git diff --numstat f7aeea92..impl-cancel2 -- internal test docs/user README.md` = 463 Einfügungen + 14 Löschungen = **477** (stimmt mit §7 überein; die Basis `f7aeea92` ist der Übergangs-Commit des Slice), Marke 700 nicht erreicht, Schätzung 350 um 36 % überschritten. Zwei Liefer-Punkte, zwei Schichten (Adapter, Dokumentation) wie in §8.

## Befunde

| ID | Kategorie | Befund | Ort | Zuständig |
|---|---|---|---|---|
| F-V1 | MEDIUM | Die Spezifikation sagt zum `CancelRequest` für `record` und `replay` gleich zu: „eine laufende Interaktion dieser zweiten Verbindung wird nicht abgebrochen“. Der Test `TestCancelRequestLaufendeInteraktionUnberuehrt` prüft nur `record`; im `replay` bleibt die Zusage ungeprüft (§7 nennt es eine Grenze, der Fake hat keine haltende Anfrage). Die Deklaration sagt „in record“, die Spezifikation nicht. Das ist eine Randform-Frage (Zusage der Spezifikation enger fassen oder einen Fall im `replay` verlangen); ich entscheide sie nicht. **Rückgabe an den Architect.** | `spec/spezifikation.md` LH-FA-05.e (Absatz *Länge des `CancelRequest`*, „und eine laufende Interaktion …“); `server_cancel_test.go` | Architect |
| F-V2 | LOW | Laufzeit der neuen Tests: `TestCancelRequestKopf` hat rund 100 Unterfälle, die je bis zu 2 s auf „keine Antwort“ warten, und läuft mit `t.Parallel()`. Die Gesamtzeit hängt an `GOMAXPROCS`: Mit 20 Kernen etwa 20 s je Lauf des Pakets, mit `GOMAXPROCS=2` gemessen 385 s für `-count=4` (rund 96 s je Lauf). Auf Runnern mit wenigen Kernen (`make test`, CI) wächst die Zeit des Pakets entsprechend; kein Fehler, aber keine Zeitgrenze des Gates ist dafür belegt. | `server_cancel_test.go` `TestCancelRequestKopf` | Implementer / Planner (Urteil) |
| F-V3 | LOW (nicht reproduziert) | Beim ersten Lauf des Mutanten K13 färbte zusätzlich ein Unterfall `TestCancelRequestLaenge/…/Länge_10000` (Fall mit vollem Rest von 9992 Bytes, nicht parallel) rot. Drei Wiederholungen von K13 und 15 360 Unterfälle des Originals unter Last (12 Busy-Loops) blieben ohne diesen Befund. Die Meldung des einen Laufs habe ich nicht gesichert. Hypothese, nicht belegt: Der Recorder lehnt 269 bis 10000 jetzt aus dem Kopf ab und schließt mit ungelesenen Bytes im Empfangspuffer; das kann auf Loopback ein `RST` statt eines geordneten Endes ergeben, und `cancelSenden` liest mit `io.ReadAll` (ein Fehler `connection reset` ist kein Timeout und läuft in den Vergleich, die `ErrorResponse` kann dann fehlen). Wenn das auftritt, ist es ein Flackern der Fälle 269, 270, 10000 mit Rest, nicht ein Fehler des Codes. | `server_cancel_test.go` `cancelSenden` | Beobachtung für §7 / Implementer |
| F-V4 | LOW | Plan §3 *Größe* nennt als gemessenen Stand weiter „354 (Basis `f7aeea92`, letzter Commit `f3ba1835`)“; die Nacharbeit hat 477 gemessen (nur in §7). Nach `AGENTS.md` §3.9 folgt §3 jeder Korrektur; die Zahl in §3 ist veraltet, die Entscheidung (Marke 700) bleibt unberührt. | Plan §3 *Größe* | Implementer |

Keine DoD-Verletzung, keine Entscheidung (ADR) verletzt, keine Randform vom Implementer entschieden.

## DoD-Stand

| Punkt | Stand |
|---|---|
| 1: Länge 12 schließt mit `PGR-W3001` in `record` und `replay` (Unit, 12 Bytes); Handbuch/README Ist-Zustand | erfüllt; mit Mutanten M1/W12 belegt (rot für den Fall gebrochen) |
| 2: `LH-FA-13` Exit-Code 0 nach pgx-Kontext-Abbruch (Integrationstest) | erfüllt; selbst gefahren, grün (0,73 s); Mutant M1 im Binary laut §7 rot (Exit 6), von mir nicht erneut gefahren |
| `make gates` grün | erfüllt am Endstand `d8a72c11`, ein Lauf |
| Review-Report liegt vor | erfüllt (`docs/reviews/2026-10-20-review-…`, ohne Nach-Review der Nacharbeit; dort suchte ich selbst: F-V1 bis F-V4) |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | noch offen (§6 Risiken tragen „offen bis Closure“); gehören zur Closure, kein Befund |

## Rückgaben

1. An den Architect: F-V1 (Zusage „laufende Interaktion“ im `replay`: Spezifikation enger fassen oder Fall verlangen).
2. An den Implementer: F-V4 (§3 der Größe nachziehen); F-V2/F-V3 in §7 als Beobachtung notieren oder entscheiden lassen.

Dieser Bericht ist ein Lauf-Beleg; er repariert nichts.
