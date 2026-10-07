# Verifikation: slice-harness-blackbox-driven — 2026-10-07

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft. F-444 und F-447 hat er an mich übergeben.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-blackbox-driven.md` (Kopf, §1 bis §8) gegen den Gesamt-Diff `aff75ee..874b29d`. Darin liegen `58e863c` (Architect: [`SPEC-049`](../../spec/spezifikation.md#spec-049--lint-profil-lint) Punkt 7, `newSession`, Plan §1, §3 und §6), `1c30e43` (Umbau, Belege in §7) und `874b29d` ([Review](2026-10-07-review-slice-harness-blackbox-driven.md), F-444 bis F-447). Während meines Laufs kam `f0ea7be` hinzu (Architect, Wortlaut von Punkt 7 nach F-445, nur `spec/spezifikation.md`). Die Datei habe ich nicht angefasst.

**Eingang:**

- DoD-Liefer-Punkte, §1, §3, §6 *Randformen* und *Risiken*, §7 *Belege des Implementers*, Commit-Messages
- `spec/spezifikation.md` `SPEC-049` Punkt 7 in der Fassung `f0ea7be` und Punkt 8; [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 und 5; [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md); [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- Produkt-Code: `internal/adapters/driven/postgres/upstream.go` an `aff75ee` und `874b29d`, `internal/adapters/driven/recording/yaml.go`
- der Review-Report (F-444 bis F-447)
- die DoD-Bestätigung des Implementers aus §7 und die Zusammenfassung des Auftrags (14 und 13 Tests, Deklarationen gleich, `make lint` 0 in den Testdateien und 67 modulweit, kein `_ =`, kein `//nolint`, Produkt-Code nur `upstream.go`)

Bei meinem Start war der Arbeitsbaum sauber, HEAD `874b29d`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

**Wie die Proben gebaut wurden.** Ich habe `aff75ee` (vor dem Umbau) und `874b29d` per `git archive` in je eine Kopie im Scratchpad entpackt. Testlisten, Sonde und Mutationen liefen im Image der Stufe `deps` (eigener Tag, aus dem Arbeitsbaum gebaut), mit Bind-Mount der Kopie und `--network=none`. Jede Mutation lief in einer frischen Kopie (`cp -r` ohne `-p`) von `874b29d`. Das Ersetzungsskript bricht ab, wenn der Suchtext nicht genau einmal trifft. Die Lint-Läufe auf den Kopien liefen mit `docker build --progress=plain --target lint`. Das Repo selbst blieb unverändert. `make lint`, `make abdeckung-check`, `make kopf-check` und `make gates` liefen im Arbeitsbaum.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `_test`-Pakete, Testliste, Abdeckung, Brücke. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Testdateien gehören zu `_test`-Paketen | Kopfzeilen an `874b29d`: `upstream_test.go` und `upstream_extended_test.go` sind `package postgres_test`, `yaml_test.go` ist `package recording_test`. Im Paket des Codes steht nur `postgres/export_test.go`. | bestätigt |
| Testliste vorher und nachher gleich | `go test -list .` je Paket an `aff75ee` und `874b29d` ergibt `postgres` 14 und `recording` 13, der `diff` ist leer. `go test -v` liefert an beiden Ständen dieselben 108 Zeilen `--- PASS` (sortiert, ohne Zeitangabe) und kein `FAIL`. | bestätigt |
| Abdeckungs-Deklarationen unverändert | An beiden Ständen gibt es 11 `Abdeckung:`-Zeilen unter beiden Pfaden, der `diff` ist leer. Auch die Zuordnung Test → Deklaration (Kommentarblock über `func Test…`) ist an beiden Ständen gleich. | bestätigt |
| `make abdeckung-check` grün | Exit 0 einzeln und im Gate-Lauf (`GATE_CHECKS` in `harness/mk/abdeckung.mk`). §7 nennt den Lauf nur über `make gates`. | bestätigt |
| Unexportiertes nur über die Brücke nach Punkt 7 | siehe Abschnitt 3 | bestätigt |

### Punkt 2: `make lint` ohne Befund in den Testdateien. **Bestätigt, mit falscher Gesamtzahl im Beleg (V-84).**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| 12 Befunde vorher in den Testdateien | Kopie `aff75ee`: 16 Zeilen unter beiden Pfaden, davon 12 in `*_test.go`. Das sind `testpackage` 3, `errcheck` 5 (`upstream_extended_test.go:98`, `:178`, `:206`, `:220`, `:314`), `gochecknoglobals` 3 (`upstream_test.go:65`, `:184`, `:186`) und `unused` 1 (`yaml_test.go:434`). Die Zeilen stimmen Wort für Wort mit dem Block *Vorher* in §7 überein. | bestätigt |
| 0 nachher | `make lint` im Arbeitsbaum (Exit 2) und Kopie `874b29d`: Unter beiden Pfaden stehen nur 4 Zeilen, alle im Produkt-Code (`upstream.go:100`, zweimal `upstream.go:226`, `yaml.go:538`), wie im Block *Nachher* in §7. Eine `lint:`-Zeile der eigenen Prüfungen gibt es nicht. | bestätigt |
| kein neuer Befund | Die Befundzeilen ohne Zeile und Spalte habe ich an beiden Ständen sortiert. `comm -13` (nur nachher) ist leer, `comm -23` (nur vorher) hat genau 12 Zeilen, alle unter den beiden Pfaden. | bestätigt |
| Gesamtzahl modulweit „79 → 67“ (§7, Auftrag) | golangci-lint meldet **72** Befunde an `aff75ee` und **60** an `874b29d` (`72 issues:` und `60 issues:`; ebenso viele Zeilen `<datei>:<zeile>:<spalte>:`). `make lint` im Arbeitsbaum meldet ebenfalls 60. Die Differenz 12 stimmt. Die beiden Zahlen in §7 liegen je um 7 zu hoch. Die 72 passt zum Endstand von `slice-harness-blackbox-kern` (103 → 72); zwischen dessen Ende und `aff75ee` ändert sich kein Go-Code. | **abweichend**, V-84 |
| ohne `//nolint`, ohne Änderung an `.golangci.yml` | `grep -rn nolint internal/` ist an beiden Ständen leer, `diff` von `.golangci.yml` ist leer. | bestätigt |
| kein `_ =` vor einem gemeldeten Fehler | Alle Blank-Zuweisungen in den drei Testdateien habe ich an beiden Ständen verglichen (Paketpräfix entfernt): 94 vorher, 93 nachher. Die Unterschiede sind nur `bereit` → `bereit()`, `beendet` → `beendet()`, `toFrontendMessage` → `ToFrontendMessage` und das entfernte `_ = ss.conn.SetReadDeadline(…)` (jetzt mit `t.Fatal` geprüft). Kein `_ =` kommt hinzu. Die drei `_ = s.Close()` in `upstream_test.go` und das in `TestCloseMitSchreibfrist` sind Bestand und waren nie ein Befund von `errcheck` (F-446). | bestätigt |

### Punkt 3: `make gates` grün. **Bestätigt.**

`make gates` im Arbeitsbaum an `f0ea7be` (Go-Code gleich `874b29d`): Exit 0. Im Log sind unter anderem `abdeckung-gegenprobe`, `a-check-negativ`, `commit-msg-gegenprobe`, `run-integration-tests` und `kopf-check-gegenprobe` grün. Für diesen Punkt trägt §7 nur die Behauptung („`make gates` grün am Stand dieses Commits“), keine Ausgabe. Diesen Lauf habe ich deshalb selbst gefahren.

### Produkt-Code: nur `upstream.go`, ohne Verhaltensänderung und ohne Export. **Bestätigt.**

- `git diff --stat aff75ee 874b29d`: Die einzige Go-Datei ohne `_test` ist `internal/adapters/driven/postgres/upstream.go` (10 Zeilen).
- Exportierte Bezeichner auf Paketebene (`type`, `func`, `var`, `const` mit Großbuchstaben, Methoden) sind an beiden Ständen dieselben: `Upstream`, `Open` und die Methoden von `session`. `newSession` ist unexportiert.
- Verhalten: Vorher legte `Open` `fe := pgproto3.NewFrontend(conn, conn)` nach dem Dial an und gab an ReadyForQuery `&session{conn: conn, fe: fe}` zurück. Jetzt legt es `s := newSession(conn)` an derselben Stelle an, liest über `fe := s.fe` und gibt `s` zurück. `newSession` setzt genau `conn` und `fe` aus `NewFrontend(conn, conn)`. Frontend-Objekt, Reihenfolge, Fehlercodes und die Fehlerpfade (`conn.Close()`, `nil`-Session) sind gleich. Neu ist nur, dass das Struct schon vor dem Startup existiert. Keine Verhaltensänderung.

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| Review mit Report | liegt vor (F-444 bis F-447) |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen: §7 trägt Platzhalter, §6 *Risiken* hat drei Ausgänge „— (bei Closure)“. Das ist Sache der Closure. |

---

## 2. F-444: Ruft `Open` `newSession` für die zurückgegebene Session, und trägt die Grenze?

**Ja, `Open` ruft es für genau diese Session** (Abschnitt 1, Produkt-Code). Zwei Mutanten sind dabei zu unterscheiden, denn sie fallen verschieden aus:

| Mutant | Was er bricht | Tests von `postgres` | Sonde (unten) | `make lint` |
|---|---|---|---|---|
| R1 (Review): an ReadyForQuery `return newSession(conn), …` | Verhalten: Die Session hat ein neues Frontend, der Lesepuffer des Aufbaus geht verloren. Punkt 7 ist nicht verletzt, denn `newSession` erzeugt den zurückgegebenen Wert. | grün | **rot** | — |
| V9 (eigen): `Open` legt wieder `fe := pgproto3.NewFrontend(conn, conn)` und `&session{conn: conn, fe: fe}` an, `newSession` ruft nur noch die Brücke | Punkt 7 („eine Funktion, die nur die Brücke ruft, gibt es nicht“). Das Verhalten ist gleich. | grün | grün | 60, kein neuer Befund, auch nicht von `unused` |

**Verletzt das die DoD oder §6? Nein.**

- Die DoD verlangt keinen Test für diese Kopplung.
- §6 *Session über `net.Pipe`* weist sie ausdrücklich als Urteil des Review aus (Grenze von `SPEC-049`). Das Review hat geurteilt (F-444), und meine Lektüre bestätigt das Urteil.
- V9 ist ein äquivalenter Mutant. Kein Verhaltenstest kann ihn fangen, und `unused` sieht `newSession` über die Brücke als benutzt. Für V9 ist die ausgewiesene Grenze deshalb die einzig mögliche Form.
- R1 war schon an `aff75ee` ungeprüft, denn das Literal `&session{conn: conn, fe: fe}` hielt kein Test. Mit dem Umbau ist also keine Prüfung verloren gegangen (Risiko *Prüfung geht im Umbau verloren*). Die Grenze darf stehen bleiben.

**Ist ein Test über die Schnittstelle möglich? Ja, für R1.** Die Sonde `TestSondeAufbauPuffer` hat nur in den Kopien gelegen, nicht im Repo. Sie gehört zu `package postgres_test` und nutzt nur `Upstream.Open`, `Receive` und `Close`:

- Der Server nimmt die StartupMessage an.
- Er schickt `AuthenticationOk`, `ReadyForQuery` und direkt dahinter eine `NoticeResponse` mit **einem** `Flush`, also einem Schreibvorgang.
- Danach wartet er 300 ms und schließt die Verbindung.
- Der Test ruft nach `Open` einmal `Receive` und erwartet genau diese Notice.

Ergebnisse:

- An `874b29d` ist die Sonde in 50 Läufen (`-count=50`) grün.
- Mit R1 ist sie rot, aus dem richtigen Grund. Die Notice lag im Puffer des verworfenen Frontends, und `Receive` endet mit `PGR-E4003 … vor ReadyForQuery beendet: unexpected EOF`, als der Server schließt.
- Mit V9 bleibt sie grün, wie erwartet.

**Grenze der Sonde.** Dass beide Nachrichten in einem Lesevorgang ankommen, garantiert TCP nicht. Auf Loopback ist es bei einem einzigen Schreibvorgang in der Praxis stabil, aber nicht zugesagt. Über `net.Pipe` wäre es deterministisch, doch `Upstream.Dialer` ist ein `net.Dialer` (Struct, nur TCP), und `Open` ist darüber nicht zu erreichen. Die Sonde wäre ein neuer Fall. §1 schließt das für diesen Slice aus („Neue Fälle oder geänderte Erwartungen — ein anderer Vorgang“).

**Empfehlung: kein eigener Folge-Slice, sondern ein Register-Eintrag bei der Closure mit einer benannten Adresse.** Die Begründung:

- Es geht um einen Test und keine Produktänderung. Ein eigener Slice dafür wäre kleiner als sein Rahmen.
- `slice-harness-mutation` schließt „Neue Tests für Mutanten, die heute überleben“ ausdrücklich aus und nimmt die Sendung deshalb nicht an.
- `slice-v1-abschluss-anmeldung` baut die Schleife des Aufbaus in `Open` für den Anmeldeaustausch ohnehin um (§3: `internal/adapters/driven/postgres`, Anmeldeaustausch). Dort wird die Frage, welches Frontend der Aufbau liest und welches die Session trägt, neu entschieden, und dort gehört der Test hin.
- Vorschlag an den Planner: In der Closure dieses Slice einen Eintrag (etwa `BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft`) mit Ausgang *geplant* auf `slice-v1-abschluss-anmeldung`, dazu ein Satz in dessen §6, dass die Sonde oben als Fall mitkommt. Dieser Slice wächst dadurch nicht.

---

## 3. Brücken nach `SPEC-049` Punkt 7 (Fassung `f0ea7be`)

| Prüfung | Beleg | Urteil |
|---|---|---|
| einzige Testdatei im Paket des Codes | `postgres`: nur `export_test.go` ist `package postgres`. `recording` hat keine Brücke und braucht keine (kein `recording.<klein>` in `yaml_test.go`). | hält |
| nur Funktionen, die weiterreichen; kein Alias, keine Variable, keine Konstante, kein `new(`, kein Literal, keine `Test…`-Funktion | Die Suche nach `type`, `var`, `const`, `new(`, `&session`, `session{` in `export_test.go` ist leer. Es gibt drei Funktionen: `ToFrontendMessage` (reicht an `toFrontendMessage`), `NeueSession` (reicht an `newSession`, liefert `driven.UpstreamSession`), `Verbindung` (liest das Feld `conn` einer übergebenen Session, setzt nichts). | hält |
| die Funktion, an die die Brücke zum Erzeugen weiterreicht, ruft auch der Produkt-Code | `newSession` ruft `Open` für die zurückgegebene Session (Abschnitt 2). Dass dies so bleibt, prüft nichts (V9), siehe Grenze. | hält, Grenze ausgewiesen |
| Werte exportierter Typen stellt der Test | `TestCloseMitSchreibfrist` stellt `net.Pipe()` und übergibt das eine Ende an `NeueSession`. | hält |
| kein selbst angelegter Wert eines unexportierten Typs in Test oder Brücke | Die Suche nach `postgres.<klein>` und `recording.<klein>` in den Testdateien ist leer. Es gibt kein `session`-Literal mehr; die Messung per AST in §7 („0 angelegte Werte“) stimmt damit überein. | hält |
| `Verbindung` gegen den neuen Wortlaut | Der Satz „Was nur an einem Wert eines unexportierten Typs zu sehen ist, den die Brücke weder übergeben bekommt noch über eine solche Funktion erzeugen lässt …“ trifft `Verbindung` nicht. Die Funktion bekommt die Session übergeben, und die Session stammt aus `Open` oder `NeueSession`. Die Mehrdeutigkeit aus F-445 ist damit aufgelöst. | hält |

---

## 4. Mutationen nach `AGENTS.md` §3.10, Stichprobe gegen §7

Der Grundlauf beider Pakete an `874b29d` ist grün. Jede Mutation lief in einer eigenen frischen Kopie mit `go test -count=1 -run` auf den genannten Test.

| ID | Zeile in §7 | Mutation (wie ausgeführt) | Ergebnis | Grund in der Ausgabe |
|---|---|---|---|---|
| V1 | Pflichtmutation, `Close` mit Schreibfrist | Zeile `_ = s.conn.SetWriteDeadline(…)` in `Close` entfernt | rot, `TestCloseMitSchreibfrist` | „Close kehrt nicht binnen 2s zurück“ (Terminate blockiert auf `net.Pipe`) |
| V2 | Send nach ErrorResponse ist PGR-E6001 | Flush-Fehler in `Send` liefert `model.Errorf(model.CodeConnectionLost, err, "x")` | rot, `TestSendNachFehlerantwort/ErrorResponse_gelesen` | „erwartet PGR-E6001, erhalten … [PGR-E4003]: x: … broken pipe“ |
| V3 | Describe ohne gültige Zielart ist PGR-E1000 | `if false && !ok` | rot, `TestReceiveFehler/Zielart` | „erwartet PGR-E1000, erhalten <nil>“ |
| V4 | unbekannter Nachrichtentyp ist PGR-E1000 | Zweig `default` liefert `&pgproto3.Sync{}, nil` | rot, `TestReceiveFehler/Zielart` | „erwartet PGR-E1000, erhalten <nil>“ |
| V5 | Fehler von `Close` ist Testfehler (`schliesse`) | `Close` liefert nach `conn.Close()` immer `net.ErrClosed` | rot, `TestSendUndReceive`, `TestReceiveEndetMitReadyForQuery`, `TestReceiveFehler/COPY`, `…/Verbindungsende` (erste Fehlzeilen) | „use of closed network connection“ aus `schliesse` |
| V6 | Aufbau nur ParameterStatus und ReadyForQuery (`bereit()`) | im Zweig `BackendKeyData` eine `notice_response` anhängen | rot, `TestOpenUndQuery` | „unerwartete Aufbau-Antwort: … notice_response“ |
| V7 | fremde Formatkennung ist PGR-E3003 | `if false && d.Format != model.Format` | rot, `TestUnmarshalFehler/fremdes_Format` | „erwartet PGR-E3003, erhalten <nil>“ |
| V8 | die letzte ErrorResponse zählt (`beendet()`, `ergebnis()`) | `merke` überschreibt nur, wenn noch nichts gemerkt ist oder `nil` kommt | rot, `TestFehlerantwortVorDemAbbruch/zwei_Fehlerantworten,_die_letzte_zählt` | „… nach der Fehlerantwort XX001 „erste““ statt 57P01 |
| V9 | (ohne Zeile in §7) Punkt 7: `Open` ohne `newSession` | Abschnitt 2 | grün in Tests, Sonde und `make lint` | äquivalenter Mutant, ausgewiesene Grenze |

Die acht Zeilen der Tabelle in §7 habe ich alle nachgefahren. Jede ist rot aus dem richtigen Grund, die Pflichtmutation eingeschlossen. Mit den sieben Mutationen des Reviews (R1 bis R7) sehe ich keine Prüfung, die im Umbau verloren ging. F-447 ist bestätigt.

---

## 5. Plan gegen Code

- **Kopf:** `Bezug` stimmt. `Berührte Spec-Stellen` nennt `SPEC-049` Punkt 7 und 8 sowie die Änderung von Punkt 7 unter der Kennung dieses Slice (`58e863c`). `f0ea7be` ändert Punkt 7 erneut, nur im Wortlaut; der Kopf deckt das inhaltlich. `make kopf-check` ist grün.
- **§1:** Das Ziel ist erreicht (Abschnitt 1). Die Abgrenzung hält: Produkt-Code nur an einer Stelle von `upstream.go`, kein Export, `.golangci.yml` unverändert, keine neue Testfunktion, Testliste gleich. Die Zerlegung „testpackage 3, dazu … zusammen 9“ stimmt (5 + 3 + 1).
- **§3:** Jede Zeile hat ihre Datei im Diff, und jede Datei des Diffs hat ihre Zeile. Die Zeile zu `export_test.go` nennt dieselben drei Weiterreichungen, die die Datei enthält.
- **§5:** „`make gates` grün“ und „`make lint` ohne Befund in den Testdateien dieser Pfade“ sind gemessen (Abschnitt 1).
- **§6 *Randformen*:** Jeder Punkt hat einen Entscheidungsort (`SPEC-049` Punkt 7 und 8, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 und 5, Architect 2026-10-07), und `1c30e43` entscheidet keine Randform, die §6 nicht nennt. Der Satz „das ist Urteil des Review“ zu `newSession` ist eingelöst (F-444, Abschnitt 2).
- **§6 *Risiken*:** *Prüfung geht im Umbau verloren* ist nach Abschnitt 4 nicht eingetreten. *Befund verdeckt statt behoben* ebenfalls nicht: `schliesse` prüft (V5), die Testdaten-Funktionen liefern frische Werte (V6, V8), kein neues `_ =`. *Abdeckung sinkt* bleibt bei `slice-harness-coverage`. Die Ausgänge setzt die Closure.
- **§7 *Belege*:** Die Belege sind vollständig für das, was die DoD verlangt (Testliste, Deklarationen, Lint-Zeilen vorher und nachher, Mutationen). `make gates` ist dort nur behauptet; ich habe es selbst gefahren. Die Gesamtzahlen sind falsch (V-84).

---

## 6. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-84 | Beleg mit falscher Zahl | §7 nennt für `make lint` „79 Befunde im Repo“ vorher und „67“ nachher, der Auftrag übernimmt die 67. Gemessen sind 72 an `aff75ee` und 60 an `874b29d` und im Arbeitsbaum (golangci-lint `72 issues:` und `60 issues:`, keine `lint:`-Zeile). Die Differenz 12, die Zeilen unter beiden Pfaden und „kein neuer Befund“ stimmen; der DoD-Punkt ist erfüllt. Die Zahlen gehen aber in die Closure-Notiz und in die Folge-Slices (`slice-lint-bestand-kern-driven`, `slice-harness-lint`), die gegen sie messen. Vor der Closure in §7 auf 72 → 60 berichtigen. Die Herkunft der +7 kann ich nicht rekonstruieren. | Plan §7 *`make lint`, Zeilen unter beiden Pfaden* | Abschnitt 1, Punkt 2; Endstand 72 in `slice-harness-blackbox-kern` |
| V-85 | Hinweis (Grenze trägt, Lücke adressieren) | F-444: Die Grenze in §6 trägt. Die Kopplung nach Punkt 7 (V9) ist äquivalent und nur durch Urteil prüfbar. Der Verhaltensmutant R1 (neues Frontend für die Session) lässt sich dagegen über die Schnittstelle fangen, die Sonde aus Abschnitt 2 ist rot unter R1 und grün in 50 Läufen ohne Mutation. Kein Bruch von DoD oder §6, die Lücke bestand schon vor dem Umbau. Vorschlag: Register-Eintrag bei der Closure mit Ausgang *geplant* auf `slice-v1-abschluss-anmeldung`, der `Open` ohnehin umbaut, und dort ein Satz in §6. Kein eigener Folge-Slice, dieser Slice wächst nicht. | `internal/adapters/driven/postgres/upstream.go` · „return s, responses, nil“; Plan §6 *Session über `net.Pipe`* | Abschnitt 2 (R1, V9, Sonde) |

---

## 7. Gate-Ergebnis und Urteil

- `make gates` an `f0ea7be` im Arbeitsbaum: **grün** (Exit 0).
- `make lint`: Exit 2 mit 60 Befunden. Keiner davon steht in einer Testdatei unter `internal/adapters/driven/postgres` oder `internal/adapters/driven/recording`, und keiner ist neu gegenüber `aff75ee`.
- `make abdeckung-check` und `make kopf-check`: grün.

**Urteil:** Die drei Liefer-Punkte sind erfüllt und selbst belegt:

- Die Tests laufen als `_test`-Pakete, die Testliste ist gleich (`postgres` 14, `recording` 13, 108 PASS-Zeilen).
- Die 11 Abdeckungs-Deklarationen sind gleich und gleich zugeordnet.
- In den Testdateien gehen die Befunde von 12 auf 0, modulweit von 72 auf 60 ohne neuen Befund. Es gibt kein `//nolint` und kein neues `_ =`.
- Der Produkt-Code ändert sich nur in `upstream.go`, ohne Export und ohne Verhaltensänderung.
- Die Brücke hält `SPEC-049` Punkt 7 in der Fassung `f0ea7be`.
- Die acht Mutationen aus §7 sind rot aus dem richtigen Grund, `SetWriteDeadline` eingeschlossen.

Vor der Closure ist V-84 zu berichtigen. V-85 ist ein Vorschlag für den Lerneintrag und das Register und blockiert nicht.
