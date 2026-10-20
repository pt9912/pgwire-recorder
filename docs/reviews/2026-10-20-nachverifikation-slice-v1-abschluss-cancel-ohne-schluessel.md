# Nach-Verifikation slice-v1-abschluss-cancel-ohne-schluessel

**Rolle:** Verifier (Modul 11), frischer Kontext. **Zweig:** `ver-cancel2` ab `origin/arch-cancel5` (708d94c3).
**Gegenstand:** `git diff d8a72c11..origin/arch-cancel5` gegen die Befunde F-V1 bis F-V4 der ersten Verifikation
(`2026-10-20-verifikation-slice-v1-abschluss-cancel-ohne-schluessel.md`, Zweig `ver-cancel`).
**Maßstab:** Slice-Plan §1, §2, §3, §6, §7, §8; Spezifikation `LH-FA-05.e`; `AGENTS.md` §3.9 bis §3.14.
Dieser Bericht ist ein Lauf-Beleg und repariert nichts; Randformen entscheidet er nicht.

## 1. Ergebnis in einem Satz

F-V1, F-V3 und F-V4 sind behoben (Mutanten selbst rot gesehen), F-V2 bleibt als LOW-Beobachtung beim Architect oder Planner;
es bleiben drei neue LOW-Befunde und kein HIGH oder MEDIUM. Empfehlung: Merge nach `main` nach Closure möglich (Abschnitt 9).

## 2. F-V1 (laufende Interaktion im replay): behoben, mit einer Enge

- Die Spezifikation (8fb96cdc, 708d94c3) bestimmt "laufend" je Modus. Der Test `TestCancelRequestLaufendeInteraktionUnberuehrtReplay`
  (Länge 12 und 16) sendet Parse, Bind, Execute ohne Sync, wartet auf drei Nachrichten am Use Case, schickt den CancelRequest
  auf einer anderen Verbindung, prüft `geöffnet 1, geschlossen 0, erzwungen 0, Shutdown 0`, sendet Sync und Terminate und verlangt am
  Ende der Antworten `1 2 D C Z` ohne `E` und ohne `PGR-`, dazu genau Parse, Bind, Execute, Sync beim Use Case. Das trägt "vollständige
  aufgezeichnete Antworten" und "Extended-Folge vor Sync".
- Mutanten selbst an einer Kopie gebrochen (Bind-Mount einer Kopie von `internal/`, Original danach unberührt), alle rot:
  - MV1 (`cancelGeschlossen` ruft `Zwangsende()`): "geöffnet 1, geschlossen 1, erzwungen 1 ..." im Replay-Test (2 Fälle), zusätzlich im Record-Test.
  - MV2 (`CloseConnection` im replay): "geschlossen 1" (2 Fälle).
  - MV3 (`OpenConnection` im replay): "geöffnet 2" (2 Fälle).
- "Cursor unverändert" ist nur teilweise gepinnt (L-1): Der Fake hat keinen Cursor; geprüft sind Open, Close, Forced, Shutdown und die
  `ClientMessage`-Folge. Zwei Mutanten, die ich selbst setzte (`cancelGeschlossen` ruft im replay `Query(ctx, 1, "SELECT 1")` bzw.
  `Sent(ctx, 1)`), bleiben **grün**. Beides sind Aufrufe, die den Cursor eines echten Use Case bewegen können; der Test zählt sie nicht
  (`sentAufrufe()` gibt es im Fake). Der Satz in §7 "der Use Case sieht keinen Aufruf außer den vier Nachrichten und dem Aufbau" ist
  damit weiter als die Prüfung (§3.11); die Deklaration in `abdeckung-unit.md` ist enger ("beendet, bricht ab oder öffnet nichts, genau die vier
  Nachrichten") und trägt. Die Spezifikationsformel "ohne Veränderung ihres Cursors" hat als Unit-Nachweis diese Lücke.

## 3. F-V3 (Länge 269, 270, 10000 und 9 bis 15): behoben

- 269, 270, 10000 mit vollem Rest verlangen nur Ende der Verbindung, `PGR-E6001` im Log und als gemerkter Fehler, keine Warnung `PGR-W3001`,
  keine Session (`restUngelesen`); `cancelSenden` toleriert `ECONNRESET` und `EPIPE` beim Schreiben. 3 x 100 Läufe der 36 Unterfälle
  (6 Stellen, 2 Modi, 3 Längen) ohne Befund, dazu 40 Läufe der Längen 9, 10, 11, 13, 14, 15 mit Rest ohne Befund, 20 Läufe des Replay-Tests ohne Befund.
- MV4 (Verbindung bleibt 30 s offen): rot, 3 Fälle je 10 s "Verbindung nicht geschlossen". MV5 (ab 269 weder Antwort noch Eintrag noch gemerkter
  Fehler): rot in `TestCancelRequestLaenge` (viele Fälle). MV6 (Warnung statt Fehler ab 269): rot. MV7 (nur gemerkt, keine ErrorResponse, alle
  Längen der Gruppe): rot in `TestCancelRequestKopf` und `TestCancelRequestLaenge`.
- Länge 9 bis 15 mit Rest (Zusage aus 708d94c3): Die Tests senden Kopf und Rest in **einem** `Write` (`cancelSenden`: `append(st.vor, f.daten...)`,
  `cancelVoll(l)` = Kopf plus `l-8` Bytes, also 1 bis 7 Bytes) und verlangen die ErrorResponse (`fehler: true`; `restUngelesen` nur für 269, 270, 10000).
  Die Behauptung des Architect stimmt im Code: `tls.go:59` `bufio.NewReader(leser{conn})` (4096 Bytes), `server.go:868` `k.br.Peek(8)`; ein Write von
  8 plus höchstens 7 Bytes kommt auf Loopback in einem `Read` an und liegt im Lesepuffer; über TLS ist es ein Record, und nach dem TLS-Schluss
  entsteht ein neuer `bufio.Reader` auf der TLS-Verbindung (`tls.go:110`). Die Zusage hängt an der Pufferung, nicht an Code, der den Rest liest;
  die Spezifikation bindet sie an "in einer Sendung" und nennt die Grenze, das ist konsistent.
- Quantoren: "Länge 9 bis 15 mit mehr als 7 Bytes Rest unzugesagt" widerspricht keinem Test, denn kein Test verlangt dort die ErrorResponse. Die
  Formulierung ist aber mehrdeutig (L-2): Bei Länge 9 bis 15 hat der deklarierte Rest höchstens 7 Bytes, "mehr als 7 Bytes Rest" kann nur gesendeten
  Überschuss hinter der Nachricht meinen; der Wortlaut sagt nicht, ob "Rest" der deklarierte oder der gesendete ist. Einen Fall mit Überschuss gibt es nicht.

## 4. F-V4 (Größe): behoben

`git diff --numstat f7aeea92..HEAD -- internal test docs/user README.md` (HEAD = 708d94c3; die Commits nach 0ee53a97 ändern nur Plan und Spezifikation):
README 3+1, abdeckung-e2e 2+0, abdeckung-unit 8+0, Handbuch 12+3, `server.go` 50+8, `server_cancel_test.go` 402+0, `server_protokollrand_test.go` 3+2,
`cancel_e2e_test.go` 66+0; zusammen **560** (546 Einfügungen, 14 Löschungen). Plan §3 nennt 560 mit Basis, Stand `0ee53a97` und Pfaden; Marke 700 nicht erreicht.
Metrik nach §3.14 vollständig.

## 5. Plan folgt Code und Spezifikation (§3.9), Handbuch (§3.11), Mutanten (§3.10)

- §1 trägt die Abgrenzung "Lesen oder Verwerfen des ungelesenen Rests" mit der Ausnahme 9 bis 15; §3 und §7 (Mutantentabelle MV1 bis MV7, Läufe, Messung, F-V2)
  stimmen mit meinen Läufen überein. MV1 bis MV7 sind in §7 keine "äquivalenten" Mutanten; die Pflicht (a) bis (c) greift nicht.
- L-3: §6 Randform 2 lautet weiter "Länge 8 bis 11, 13 bis 15, 269 bis zur Höchstlänge: `PGR-E6001` mit `ErrorResponse` (Exit-Code 6), Punkt"; für 269 bis 10000
  mit Rest schränkt die spätere Randform (Ablehnung aus dem Kopf bei ungelesenem Rest) das ein, die erste Zeile bleibt ohne Verweis stehen (§3.9).
- Handbuch und README sagen für 8 bis 11, 13 bis 15, 269 bis 10000 nur `PGR-E6001` und Exit-Code 6, nicht die ErrorResponse. Das geht nicht über die
  Spezifikation hinaus; für 269 bis 10000 mit Rest belegt es der gemerkte Fehler im Unit-Test (kein E2E-Lauf mit Rest; die Probe sendet Rest 0).
- Die Plattformgrenze (macOS, Windows) steht an der Zusage in der Spezifikation und im Plan §6; das Handbuch sagt nichts über die Zustellung zu.

## 6. Gates

`git diff c30b830a..HEAD --stat`: nur der Slice-Plan (7 Zeilen) und `spec/spezifikation.md` (3 Zeilen); Code, Tests, Handbuch und Deklarationen sind seit dem
grünen `make gates` des Implementers (Exit 0, 10 min 55 s auf `c30b830a`) byte-gleich. Ein neuer Volllauf ist nicht sinnvoll: Er prüfte dieselben Bytes des
Codes; nur die Spezifikation könnte `docs-check`, `kopf-check`, `meldungskatalog-check` und `abdeckung-check` berühren. Dafür habe ich `make gates-schnell`
am Endstand gefahren: Exit 0 (Cache-Treffer, 12 s), `baseline-verify` OK, `docs-check` "1270 Datei(en) geprüft, 0 Befund(e)", gofmt, vet und Unit-Tests mit
Schwelle 94 grün, `workflows-check` ohne Befund. Der Stop-Hook liest den Nachweis des vollen Laufs; `spec/` liegt nicht auf seiner Positivliste, das ist beim
Merge zu beachten (Vorgang, nicht Gegenstand dieses Berichts).

## 7. Befunde

| Nr. | Stufe | Befund | Zuständig |
|---|---|---|---|
| L-1 | LOW | "Cursor unverändert" im replay-Test nur über Open, Close, Forced, Shutdown und die `ClientMessage`-Folge gepinnt; Mutanten mit `Query` oder `Sent` in `cancelGeschlossen` bleiben grün. Entweder `Query` und `Sent` im Test zählen, oder den Satz in §7 und die Formel der Spezifikation auf das Geprüfte eng fassen. | Implementer (Test, §7), Architect (Formel, falls eng gefasst) |
| L-2 | LOW | Spezifikation 708d94c3: "Länge 9 bis 15 mit mehr als 7 Bytes Rest" ist mehrdeutig (deklarierter oder gesendeter Rest); kein Fall mit gesendetem Überschuss. Nicht von mir entschieden. | Architect |
| L-3 | LOW | Plan §6 Randform "8 bis 11, 13 bis 15, 269 bis zur Höchstlänge: PGR-E6001 mit ErrorResponse" steht ohne Einschränkung neben der späteren Randform (§3.9). | Implementer (Plan) |
| F-V2 | LOW (unverändert) | Laufzeit von `TestCancelRequestKopf` hängt an `GOMAXPROCS` (rund 100 parallele Fälle mit bis zu 2 s); in §7 als Beobachtung geführt, nicht entschieden. | Architect oder Planner |

Kein HIGH, kein MEDIUM.

## 8. DoD-Stand (Plan §2)

- Liefer-Punkt 1 (Länge 12: Unit-Test, Handbuch, README): erfüllt; Mutanten rot gesehen.
- Liefer-Punkt 2 (pgx mit Kontext-Abbruch, Exit-Code 0, Integrationstest): unverändert seit der ersten Verifikation; Code und Test byte-gleich.
- `make gates` grün: Nachweis des Implementers auf `c30b830a`, am Endstand ergänzt um `gates-schnell` (Abschnitt 6).
- Offen und nicht Aufgabe dieses Laufs: Closure-Notiz mit Lerneintrag, Register, Ausgänge der zwei Risiken in §6 (stehen noch "offen bis Closure"), drei Paarungen.

## 9. Empfehlung

Merge nach `main` ist vertretbar, sobald der Implementer die Closure erledigt (Notiz, Lerneintrag, Register, Risiko-Ausgänge) und L-3 im Plan nachzieht;
L-1 und L-2 sind LOW und dürfen als Beobachtung oder Folgepunkt ins Register statt als Blocker. Ein weiterer Volllauf von `make gates` ist nur nötig, wenn der
Merge Code ändert.
