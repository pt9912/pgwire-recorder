# Verifikation: slice-lint-bestand-kern-driven — 2026-10-07

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft. F-464 hat er an mich übergeben. F-459 bis F-463 hat der Architect in `8d765ba` entschieden, F-460 und F-463 hat der Implementer in `b215e70` umgesetzt.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-lint-bestand-kern-driven.md` (Kopf, §1 bis §8) gegen den Gesamt-Diff `7ac016b..253d0d1`. Darin liegen `cedd891` (Charakterisierungstests), `d840888` (Umbau), `ac2e662` (Belege in §7), `a8b341e` ([Review](2026-10-07-review-slice-lint-bestand-kern-driven.md), F-459 bis F-464), `8d765ba` (Architect), `b215e70` (Kommentare und §7) und `253d0d1` (Planner legt `slice-tests-ueberlebende-mutanten` an). Der Go-Code von `253d0d1` unterscheidet sich von `d840888` nur in Kommentaren (`b215e70`).

**Eingang:**

- Plan §1, §2 (DoD), §3, §4, §6 *Randformen* und *Risiken*, §7 *Belege des Implementers*, §8; die Commit-Messages der acht Commits
- `spec/spezifikation.md` `SPEC-049` Punkt 4, 8 und 9 (`git diff 7ac016b 253d0d1 -- spec/` ist leer) und der Abschnitt zu `LH-FA-18.a` (*Fehler*); [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 5; [ADR-0001](../plan/adr/0001-hexagonale-architektur.md); [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- `AGENTS.md` §3.6, §3.9 bis §3.12
- Produkt-Code an `7ac016b` und `253d0d1`: `internal/hexagon/model/extended.go`, `recording.go`, `internal/hexagon/services/replay.go` (auch `NewReplayService`, `OpenConnection`, `zuordnen`), `internal/adapters/driven/postgres/upstream.go`, `internal/adapters/driven/recording/yaml.go`
- `docs/plan/planning/open/slice-tests-ueberlebende-mutanten.md` ganz; die Hunks von `253d0d1` in `slice-harness-mutation`, `slice-harness-coverage`, `slice-lint-bestand-driving` und der Roadmap

Bei meinem Start war der Arbeitsbaum sauber, HEAD `253d0d1`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

**Wie die Proben gebaut wurden.** `7ac016b`, `cedd891`, `d840888` und `253d0d1` habe ich per `git archive` und `tar -x --touch` in je eine Kopie im Scratchpad entpackt, ohne `cp -p`. Lint lief je Kopie mit `docker build --progress=plain --target lint` aus dem eigenen Kontext. Testlisten, Unit-Tests, Differenzproben und Unit-Mutationen liefen im Image der Stufe `deps` unter eigenem Tag, per Bind-Mount und mit `--network=none`. Jede Mutation lief in einer frischen Kopie aus `git archive 253d0d1`. Das Ersetzungsskript bricht ab, wenn der Suchtext nicht genau einmal trifft, und prüft per Prüfsumme, dass sich die Datei geändert hat. Für den Integrationslauf habe ich in der Kopie den Image-Tag des Runners auf einen eigenen umgeschrieben. Der Runner legt je Lauf ein eigenes internes Netz, eigene Container und ein eigenes Volume an (Namen mit seiner PID) und räumt sie danach ab. Danach habe ich meine beiden Images und die Kopien entfernt. Hostweit habe ich nichts aufgeräumt, und kein Lauf hat in den Arbeitsbaum geschrieben. `make gates` lief im Arbeitsbaum.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Komplexität unter den Schwellen, ohne Verhaltensänderung. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| die vier Funktionen unter den Schwellen aus `SPEC-049` Punkt 4 (`gocognit` > 20, `gocyclo` > 15, `cyclop` > 15) | Kopie `253d0d1`, in `.golangci.yml` nur die drei Schwellen auf 1 gesetzt: `(Group).validate` 6/7/7, `validateClient` 11/10/10, `validateServer` 7/6/6, `(*cursor).objekte` 12/10/10, `extendedNachspielen` 12/7/7, `toResponse` 2/10/10, `ohneFelder` –/7/7, `spalten` –/2/2, `werte` 4/3/3, `fromDTO` 4/4/4, `sessionFromDTO` 6/6/6, `geprueftFromDTO` 5/6/6 (`gocognit`/`gocyclo`/`cyclop`). Jeder Wert stimmt mit der Tabelle in §7 überein. Der höchste liegt bei 12 von 20 bzw. 10 von 15, keiner knapp unter einer Schwelle. | bestätigt |
| Testliste vorher und nachher gleich | `go test -list .` je Paket: `model` 9 / 10 / 10 / 10, `services` 53 / 53 / 53 / 53, `postgres` 14 / 14 / 14 / 14, `recording` 13 / 14 / 14 / 14 (an `7ac016b` / `cedd891` / `d840888` / `253d0d1`). `diff` `cedd891` gegen `d840888` und `d840888` gegen `253d0d1` ist leer. Gegen `7ac016b` kommen nur `TestValidateFehlerReihenfolge` und `TestUnmarshalFehlerReihenfolge` hinzu. | bestätigt |
| keine Erwartung geändert | `d840888` berührt keine Testdatei. `b215e70` ändert in den beiden Testdateien nur Kommentarzeilen. `cedd891` fügt nur an. `go test -count=1 -v` der vier Pakete liefert an `cedd891` und `253d0d1` dieselben 230 `--- PASS`-Zeilen (sortiert, ohne Zeitangabe), kein `FAIL`. | bestätigt |
| `make test` und `make test-integration` grün | im Gate-Lauf, Abschnitt 5 | bestätigt |

**Eigene Differenzproben gegen den Code vor dem Umbau (`cedd891` = Code von `7ac016b`).** Beide Proben liegen nur in der Kopie, sind deterministisch geseedet und hashen die Ausgabe je Fall.

- `Validate`: 300 000 zufällige Extended-Interaktionen mit 1 bis 3 Gruppen, 0 bis 3 Client-Nachrichten (die sieben Arten und eine unbekannte, drei Zielarten einschließlich leer) und 0 bis 3 Server-Nachrichten (acht Arten einschließlich einer unbekannten). Es entstehen 37 verschiedene Fehlertexte, und der Hash über alle Ergebnisse ist an beiden Ständen gleich. Damit ist die Reihenfolge der Fehler auch bei drei und mehr gleichzeitigen Verletzungen unverändert, nicht nur in den 13 Fällen des Charakterisierungstests.
- `Unmarshal`: 40 000 zufällige Aufzeichnungen mit 1 bis 3 Sessions, falschen Kennungen und Nummern, negativem und gültigem `offset_ms`, `empty_sessions` und fünf Formen (gültig einfach, abgeschnitten, gültig extended, flush am Ende, unbekannter Typ). Verglichen werden der Fehlertext und bei Erfolg die ganze Aufzeichnung (`OffsetMS` dereferenziert). Es entstehen 1700 verschiedene Ergebnisse, und der Hash ist an beiden Ständen gleich.

Für `toResponse` und `objekte` habe ich keine Differenzprobe gebaut. Gelesen: Der Type-Switch in `toResponse` und `ohneFelder` läuft über disjunkte Zeigertypen, und `%T` an `msg` statt an `m` nennt denselben dynamischen Typ. In `objekte` steht die Haltebedingung Zeichen für Zeichen unverändert in `halt`, und `extendedNachspielen` gibt `false` genau dort zurück, wo vorher `return o` stand. Das bestätigt F-464.

### Punkt 2: `make lint` ohne Befund unter den vier Pfaden. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| vorher 16 unter den vier Pfaden, 32 modulweit | Kopie `7ac016b`: 32 Befundzeilen, 16 davon unter `internal/hexagon/model`, `internal/hexagon/services`, `internal/adapters/driven/postgres` und `internal/adapters/driven/recording`, keine in einer `*_test.go`, keine `lint:`-Zeile der eigenen Prüfungen. Die 16 Zeilen sind sortiert byte-gleich mit dem Block *Vorher* in §7. `cedd891`: dieselben 32 Zeilen. | bestätigt |
| nachher 0 unter den vier Pfaden, modulweit 16 | Kopien `d840888` und `253d0d1`: je 16 Befundzeilen, gleich. Keine unter den vier Pfaden. Verteilung: `internal/adapters/driving/pgwire/server.go` 14, `internal/adapters/driving/cli/cli.go` 1, `internal/bootstrap/bootstrap.go` 1. | bestätigt |
| kein neuer Befund | sortiert mit Zeile und Spalte: `comm -13` (nur nachher) ist leer, `comm -23` (nur vorher) sind genau die 16 Zeilen aus §7. | bestätigt |
| ohne `//nolint` | Das Muster aus `tools/harness/lint.sh` über `internal`, `cmd` und `test` an `253d0d1` trifft nichts, und die Lint-Läufe zeigen keine `lint:`-Zeile. | bestätigt |
| ohne neue Ausnahme in `.golangci.yml` | `git diff 7ac016b 253d0d1 -- .golangci.yml tools/harness/lint.sh Dockerfile .dockerignore` ist leer. | bestätigt |
| Zeilen vorher und nachher im Bericht | §7 nennt die 16 Zeilen vorher und für nachher „keine Zeile“. Das habe ich gemessen. | bestätigt |

### Punkt 3: Kommentar über `(*cursor).letzteNummer` (V-81). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| „ohne erwartete Interaktion 0“ gestrichen | Der Satz fehlt an `253d0d1`. | bestätigt |
| Invariante mit Grund | „Die Session eines Cursors hat mindestens eine erwartete Interaktion: zuordnen vergibt nur Sessions aus frei, und NewReplayService nimmt dort nur solche auf.“ Das stimmt mit dem Code überein. `NewReplayService` hängt eine Session nur mit `len(sess.Interactions) > 0` (nach `ohneLebendpruefungen`) an `frei`. `zuordnen` setzt `c.session` nur aus `s.frei[0]`. Die beiden Aufrufer von `letzteNummer` laufen erst nach `zuordnen`. `OpenConnection` legt den Cursor ohne Session an und ruft `letzteNummer` nicht. Den Grund hat das Review mit L1 rot belegt. | bestätigt |
| kein Wert für den leeren Zweig zugesagt | „Die Prüfung auf die leere Liste schützt allein den Index.“ nennt keinen Wert. | bestätigt |
| Zweig bleibt, Testliste gleich | `if len(c.session.Interactions) == 0 { return 0 }` steht unverändert. Testliste: siehe Punkt 1. | bestätigt |

### Punkt 4: `make gates` grün. **Bestätigt.**

Für diesen Punkt trägt §7 nur die Zeile „`make gates` grün an `d840888`, an `ac2e662` und am Stand des Commits …“, ohne Ausgabe. Den Lauf habe ich deshalb an `253d0d1` selbst gefahren: Exit 0 (Abschnitt 5).

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| Review mit Report | liegt vor (F-459 bis F-464) |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen: §7 trägt Platzhalter, und zwei Risiken in §6 tragen „— (bei Closure)“. Das ist Sache der Closure; Vorschläge stehen in Abschnitt 6. |

### Weitere Schwerpunkte des Auftrags

- **Charakterisierungstests: nur die zwei neu.** Bestätigt (Punkt 1). An `cedd891`, also am alten Code, sind sie grün, denn die Grundläufe dort sind grün. Die Fälle zählen 13 und 5, wie §7 angibt.
- **Abdeckungs-Deklarationen gleich.** An `7ac016b` und `253d0d1` stehen je 137 Zeilen `Abdeckung:` in den Testdateien unter `internal` und `test`. Ohne Zeilennummer und sortiert sind beide Mengen gleich. Der Diff `7ac016b..253d0d1` enthält keine Zeile mit `Abdeckung`, und `docs/user/` ist unverändert. `make abdeckung-check` ist im Gate-Lauf grün.

---

## 2. Mutationen, Stichprobe gegen §7

Grundlauf an `253d0d1` ohne Mutation: die Unit-Tests von `services` sind grün (W0). Integration ohne Mutation: grün im Gate-Lauf.

| ID | Zeile in §7 / Auftrag | Mutation (wie ausgeführt) | Ergebnis |
|---|---|---|---|
| W1 | F-460: „eine Art ohne Bestätigung in der Interaktion wird nicht nachgespielt“ | `bestaetigt[m.Type] > 0` → `>= 0` | **rot**, `TestReplayExtendedDiagnoseLebensdauer`, aus dem richtigen Grund: Die Fälle „S2 abgelehntes parse überschreibt nicht“ (erwartet `"SELECT 3"`, empfangen `"SELECT 2"`) und „verworfenes bind legt kein Portal an“ scheitern, weil eine unbestätigte Nachricht nachgespielt wird. |
| W2 | Bestätigungen je Interaktion gezählt | Zählschleife nur über `in.Groups[:1]` | rot, `TestReplayExtendedAbweichung`, `TestReplayExtendedDiagnoseSpaeteBestaetigung` |
| W3 | grüner Mutant M3 (Adresse `slice-tests-ueberlebende-mutanten`) | Zeile `bestaetigt[m.Type]--` entfernt | grün, wie §7 und das Review angeben |
| W4 | Nachricht am Cursor wird nicht nachgespielt | `ni >= c.nachricht` → `ni > c.nachricht` | rot, `TestReplayExtendedDiagnoseLebensdauer` |
| W5 | nach dem Halt kein Transaktionsende | Ergebnis von `extendedNachspielen` ignoriert | rot, `TestReplayExtendedAbweichung`, `TestReplayExtendedDiagnoseAnweisung`, `TestReplayExtendedDiagnoseLebensdauer` |
| X1 | sync/flush am Ende vor den Server-Nachrichten | `validateServer` vor die sync/flush-Prüfung | rot, `TestValidateFehlerReihenfolge/sync_am_Ende_einer_vorderen_Gruppe_vor_einer_unbekannten_Server-Nachricht` |
| X2 | `QF1001`, Bedingung *letzte Stelle* | `(!letzte \|\| si != len(g.Server)-1)` → `(!letzte)` | rot, `TestValidateFehler/ready_for_query_vor_dem_Ende_der_letzten_Gruppe` und zwei Fälle von `TestValidateFehlerReihenfolge` |
| X3 | Session ohne Interaktion nur mit `empty_sessions` | Prüfung in `sessionFromDTO` entfernt | rot, `TestUnmarshalFehler/Session_ohne_Interaktion`, `TestFelderDerVersion1` |
| X4 | grüner Mutant Y2 | `geprueftFromDTO(1, ii, id)` | grün, wie §7 angibt |
| X5 | `ohneFelder`: no_data | `NoData` → `ResponseEmptyQueryResponse` | rot, `TestSendUndReceive` |
| X6 | eine andere Antwort ist nicht unterstützt | `ohneFelder` liefert am Ende `ResponseNoData, true` | rot, `TestReceiveFehler/COPY`; `TestNichtVermittelbar/COPY` hängt dabei bis zum Zeitlimit von 10 min |
| X7, X8 | grüne Mutanten P8 und P2 | Zeile `TableOID:` bzw. `ColumnNumber:` in `spalten` entfernt | beide grün in den Unit-Tests von `postgres`, wie §7 angibt (Integration hat das Review gefahren, P2) |
| I1 | `ohneFelder`: empty_query_response | `EmptyQueryResponse` → `ResponseNoData`, Integrationslauf | rot, `TestE2EReplayLebendpruefungPgxpool`, `…DatabaseSQL`, `…WiePostgres`, wie §7 angibt |

Damit ist der F-460-Satz über die Art ohne Bestätigung von einer Mutation gefangen, wie §6 es für sein Bleiben verlangt. Die vier grünen Mutanten aus §6 sind tatsächlich grün (W3, X4, X7, X8). Zehn Zeilen der Tabelle in §7 habe ich nachgefahren (W1, W2, W4, W5, X1 bis X3, X5, X6, I1), und jede ist aus dem genannten Grund rot.

---

## 3. F-460, F-463, F-464 und der gestrichene Satz über den Server

**F-460: umgesetzt, mit einem Rest (V-89).** Der zweite Absatz des Doc-Kommentars von `extendedNachspielen` sagt nur noch zu, was ein Test fängt: Die Zählung über die ganze Interaktion fängt W2, die Art ohne Bestätigung fängt W1. Der Satz über „die ersten n“ ist weg. Im **ersten** Absatz, den `d840888` neu geschrieben hat, steht aber weiter „spielt die **angenommenen** Client-Nachrichten einer Extended-Interaktion in Sendereihenfolge nach“. Unter W3 (ohne Dekrement) spielt die Funktion ein abgelehntes `parse` nach einem angenommenen derselben Art nach, und alle Tests bleiben grün. Das Wort „angenommenen“ sagt also genau das zu, was mit dem gestrichenen Satz entfallen sollte. Siehe V-89.

**Der gestrichene Satz über den Server: in Ordnung.** `b215e70` streicht zusätzlich „Der Server bestätigt jedes angenommene parse, bind und close; nach einer Ablehnung verwirft er bis zum Sync.“ Der Satz sagte nichts über den Code zu, er beschrieb das Verhalten des Servers als Grund. Seine zweite Hälfte steht normativ in der Spezifikation (Abschnitt `LH-FA-18.a`, *Fehler*: „Nach einer `ErrorResponse` verwirft ein Server Nachrichten bis zum nächsten `Sync`“). Die erste Hälfte ist eine Tatsache des Protokolls, die ein Unit-Test nicht prüfen kann. Die Entscheidung des Architect in §6 lautet „An seiner Stelle steht nur, was ein Test fängt“, und die Streichung bleibt in diesem Rahmen. Verloren geht nur die Erklärung, warum eine Bestätigung als Annahme zählt. Sie kommt mit dem Test aus `slice-tests-ueberlebende-mutanten` zurück (dort DoD Punkt 1). Kein Befund.

**F-463: umgesetzt.** Beide Kopfkommentare sagen jetzt „Der Test hält fest, welche der Bestand meldet; eine Zusage ist das nicht, und wer die Reihenfolge bewusst ändert, passt ihn an.“ Das entspricht §6 dem Sinn nach und fast wörtlich. „und welche, steht fest“ ist an beiden Stellen weg. Die Doc-Kommentare von `Group.validate`, `sessionFromDTO` und `geprueftFromDTO` beschreiben die Reihenfolge, und X1 bis X3 sowie die Mutationen in §7 fangen sie. Dass `b215e70` die beiden Kommentarblöcke umbricht, ohne die Zeilen neu zu füllen, ist Form (V-92).

**F-464: bestätigt.** Siehe Abschnitt 1, Punkt 1: Differenzproben für `Validate` und `Unmarshal`, Lesen für `toResponse` und `objekte`.

---

## 4. Adresse `slice-tests-ueberlebende-mutanten`

Das Risiko *Verhalten ändert sich unbemerkt* in §6 nennt drei Bedingungen: Der Slice ist bei der Closure angelegt, er ist nicht geschlossen, und er schließt die vier Mutanten nicht aus.

| Bedingung | Befund | Urteil |
|---|---|---|
| angelegt | `docs/plan/planning/open/slice-tests-ueberlebende-mutanten.md` liegt seit `253d0d1` in `open/`. | erfüllt |
| führt die vier | §1 *Gegenstand* nennt alle vier mit Quelle, Test-Idee, Mutation und Grenze. Erstens M3, `extendedNachspielen` ohne `bestaetigt[m.Type]--`. Zweitens und drittens P8 und P2, `spalten` ohne `TableOID` bzw. `ColumnNumber`, je für sich. Viertens Y2, `geprueftFromDTO` mit Session 1, zusammen mit den Konstanten je Meldung. Je ein DoD-Liefer-Punkt verlangt die rote Mutation. Die Test-Idee zu M3 fängt W3: Ohne Dekrement überschreibt das zweite, abgelehnte `parse` `s1` in Gruppe 2 das erste, und die Diagnose nennt `SELECT 2`. | erfüllt |
| schließt sie nicht aus | Kein Punkt unter *Ausdrücklich NICHT* trifft einen der vier. Die Schicht-Abgrenzung (PGWire-Adapter, `test/integration`) lässt P8 und P2 im Paket `postgres_test` stehen. „Produkt-Code“ lässt den Kommentar von `extendedNachspielen` ausdrücklich zu. | erfüllt, mit einer Bedingung (unten) |
| bleibt nach diesem Slice offen | Start erst, wenn `slice-harness-coverage` in `done/` liegt. Die Reihe in §4 lautet: dieser Slice, `slice-lint-bestand-driving`, `slice-harness-lint`, `slice-harness-commit-struktur-id`, `slice-harness-integration-wait`, `slice-harness-abdeckung-gate`, `slice-harness-coverage`, dann der Nehmer, danach `slice-harness-mutation`. Bei der Closure dieses Slice liegt er also sicher in `open/` oder `next/`. `slice-harness-mutation` §4 startet erst nach ihm, und sein §1 verweist auf ihn statt die Tests auszuschließen. | erfüllt |

**Die eine Bedingung:** Y2 ist nur bedingt angenommen. Nach §6 des Nehmers (Randform *Ort in der Meldung*, offen) wird Y2 „zum äquivalenten Fall mit ausgewiesener Grenze“, wenn der Architect keine Zusage festhält, und der Liefer-Punkt entfällt dann mit Begründung. Das ist kein Ausschluss, denn der Fall hat dann einen begründeten Ausgang statt eines stillen. Nur die Formulierung des Risiko-Ausgangs hier sollte das wissen (Vorschlag in Abschnitt 6). Für M3 gilt Ähnliches in kleinerem Umfang: Geprüft wird nur ein benanntes Statement. Das ist eine Grenze des Tests, keine Ablehnung des Mutanten.

Ein Hinweis zum Nehmer, außerhalb dieser DoD: V-90.

---

## 5. Gate-Ergebnis

- `make gates` an `253d0d1` im Arbeitsbaum: **grün**, Exit 0. Im Log sind unter anderem grün: `a-check`, `a-check-negativ`, `baseline-verify`, `abdeckung-gegenprobe`, `commit-msg-gegenprobe`, `run-integration-tests` (kein `--- FAIL`, 2 `--- SKIP`), `kopf-check-gegenprobe` und `d-check` mit 0 Befunden.
- Lint (Stufe `lint`) an `7ac016b`, `cedd891`, `d840888` und `253d0d1`: Exit ungleich 0 mit 32, 32, 16 bzw. 16 Befunden (Abschnitt 1, Punkt 2). `make lint` ist ein Werkzeug, kein Gate.
- Commit-Messages: Alle acht nennen `slice-lint-bestand-kern-driven` und [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), keine nennt eine `SPEC-` oder `ARC-`Kennung.

---

## 6. Plan gegen Code

- **Kopf:** `Bezug` ([`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0001](../plan/adr/0001-hexagonale-architektur.md)) und `Berührte Spec-Stellen` (`SPEC-049`, unverändert) stimmen. `kopf-check` ist im Gate-Lauf grün.
- **§1:** Das Ziel ist erreicht. Die Aufteilung 7 + 8 + 1 = 16 stimmt mit dem Lauf überein. Die Abgrenzung hält: Kein Befund unter `driving/`, `bootstrap` oder `cmd/` ist berührt, die 16 übrigen sind unverändert. Keine Testdatei hatte vorher einen Befund. Keine Ausnahme, kein `//nolint`. Keine exportierte Signatur ist geändert, `Query` behält seine Signatur, und `ctx` heißt `_`. Die Ausnahme für Charakterisierungstests hat der Architect bestätigt, und die Tests halten ihre Grenze ein: zwei Funktionen, im eigenen Commit vor dem Umbau, am alten Code grün, ohne Deklaration.
- **§3:** Jede Zeile hat ihre Änderung im Diff und umgekehrt. `export_test.go` ist nicht berührt.
- **§4:** Keine Rückführung ausgelöst. Der Umbau war in einer Review-Sitzung prüfbar, und kein Befund verlangte eine Verhaltensänderung.
- **§5:** `make gates` grün und Lint ohne Befund unter den vier Pfaden sind gemessen.
- **§6 *Randformen*:** Alle fünf entschiedenen hält der Code ein. Er entscheidet keine Randform, die §6 nicht nennt.
- **§6 *Risiken*, Vorschläge für die Ausgänge** (setzen muss sie die Closure):
  - *Verhalten ändert sich unbemerkt:* **eingetreten** → Folge-Slice `slice-tests-ueberlebende-mutanten`. Begründung: Vier grüne, verhaltensändernde Mutanten (M3, P8, P2, Y2), alle schon am Code vor dem Umbau grün. Der Nehmer führt sie in §1 und nimmt an (Abschnitt 4). Im Ausgang zu nennen: Y2 entfällt dort mit Begründung, wenn der Architect die Zusage zum Ort in der Meldung nicht festhält.
  - *Befund verschoben statt behoben:* **entfallen**. Begründung: Die höchste Komplexität einer herausgelösten Funktion liegt bei 12 von 20 bzw. 10 von 15 (Abschnitt 1). `ctx` heißt `_`, weil `driven.UpstreamSession` den Parameter verlangt (§6).
  - *Aufzählung weicht vom Lauf ab:* trägt schon seinen Ausgang (entfallen). Das bestätige ich, denn der Lauf an `7ac016b` deckt sich mit §1.
- **§7 *Belege*:** Vollständig für das, was die DoD verlangt: Lint-Zeilen vorher und nachher, Komplexität, Testliste, Charakterisierungstests, Mutationen, grüne Mutanten. `make gates` steht nur als Zeile, den Lauf habe ich selbst gefahren. Die Zahlen stimmen. Der Kopf der Belege nennt einen veralteten Stand (V-91).
- **§8:** Die Sichtung nennt den Stand `0d464b4`, die Treffer sind plausibel. Für die Closure zu prüfen: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` ist mit F-460 und V-89 in diesem Slice getroffen.

---

## 7. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-89 | LOW (Kommentar-Zusage ohne Test, Rest von F-460) | Der erste Absatz des Doc-Kommentars von `extendedNachspielen` sagt „spielt die angenommenen Client-Nachrichten … nach“. Ohne `bestaetigt[m.Type]--` (W3) spielt die Funktion ein abgelehntes `parse` nach einem angenommenen derselben Art nach, und alle Tests bleiben grün. „angenommenen“ sagt damit dasselbe zu wie der Satz über „die ersten n“, der nach der Entscheidung des Architect zu F-460 entfallen sollte, bis der Test aus `slice-tests-ueberlebende-mutanten` kommt. `AGENTS.md` §3.11 verlangt eine engere Fassung oder einen Test. Vorschlag für die enge Fassung: „spielt die Client-Nachrichten einer Extended-Interaktion in Sendereihenfolge nach, deren Art in der Interaktion bestätigt ist, bis …“. Das prüfen W1 und W2. Adresse: Implementer, nur ein Kommentar, vor der Closure. Oder der Architect nimmt das Wort ausdrücklich in die Rückkehr der Zusage mit dem Test des Nehmers auf (dort DoD Punkt 1). | `internal/hexagon/services/replay.go` · „extendedNachspielen spielt die angenommenen Client-Nachrichten einer“ | W3 grün; W1, W2 rot |
| V-90 | Hinweis (Planner, Nehmer) | §6 von `slice-tests-ueberlebende-mutanten` führt *Was ein Test festhält* als „entschieden … (Architect zu F-463 in `slice-lint-bestand-kern-driven`)“: „Nach außen gelten Code, Exit-Code und die zugesagten Bestandteile des Textes. Ein Test vergleicht nur diese, nicht den ganzen Text.“ Die Entscheidung zu F-463 in §6 dieses Slice sagt das nicht. Sie sagt, dass die Reihenfolge kein Vertrag ist und dass nach außen nur der Text geht, mit `PGR-E3003` und Exit-Code 3. Wie ein Test vergleichen soll, sagt sie nicht. Die Regel für Tests ist also eine neue Aussage mit fremder Herkunft. Sie berührt diesen Slice nicht, denn seine Charakterisierungstests vergleichen bewusst den ganzen Text und sagen nichts zu. Vorschlag: im Nehmer als offene Randform an den Architect, oder die Herkunft richtig angeben. | `docs/plan/planning/open/slice-tests-ueberlebende-mutanten.md` · „Ein Test vergleicht nur diese, nicht den“ | Lesen: `git show 8d765ba` |
| V-91 | Hinweis (Form der Belege) | Der Kopf der Belege in §7 sagt „Stand `d840888`“. Die Belege enthalten aber Ergebnisse späterer Stände: die Mutation `>= 0` „aus `git archive 8d765ba`“, den Gate-Lauf „am Stand des Commits, der die Kommentare … fasst“ (`b215e70`) und die neue Adresse. Für die Closure sollte der Kopf `b215e70` nennen oder je Zeile den Stand. | `docs/plan/planning/in-progress/slice-lint-bestand-kern-driven.md` · „Belege des Implementers** (Stand `d840888`“ | Lesen |
| V-92 | Hinweis (Form) | `b215e70` hat die Kopfkommentare der beiden Charakterisierungstests umgeschrieben, ohne die Zeilen neu umzubrechen. Je eine Zeile ist rund 130 bzw. 105 Zeichen lang, die übrigen Zeilen desselben Blocks höchstens rund 85. `gofmt` und `make lint` melden das nicht, es ist nur Form. | `internal/hexagon/model/extended_test.go` · „Client-Nachrichten der Reihe nach vor, je Nachricht erst der Typ, dann die Stellung von flush und sync, dann die Zielart; danach“; `internal/adapters/driven/recording/yaml_test.go` · „die Sessions der Reihe nach vor, je Session erst ihre Kennung, dann ob sie Interaktionen trägt, dann ihre“ | Lesen |

---

## 8. Urteil

Die drei Liefer-Punkte sind erfüllt, und ich habe sie selbst belegt:

- **Lint:** Modulweit sinken die Befunde von 32 auf 16, unter den vier Pfaden von 16 auf 0. Es gibt keinen neuen Befund, kein `//nolint`, und `.golangci.yml`, `lint.sh` und `Dockerfile` sind unverändert. Die 16 übrigen Befunde liegen unter `driving/`, `bootstrap` und `cmd/` und gehören zu `slice-lint-bestand-driving`.
- **Komplexität:** Alle vier Funktionen und ihre acht Hilfsfunktionen liegen deutlich unter den Schwellen, und die Werte stimmen mit §7 überein.
- **Ohne Verhaltensänderung:** Die Testliste ist gleich, abgesehen von den zwei Charakterisierungstests. An `cedd891` und `253d0d1` stehen dieselben 230 PASS-Zeilen. Die Differenzproben über 300 000 bzw. 40 000 Fälle ergeben für `Validate` und `Unmarshal` an beiden Ständen dieselben Ergebnisse, auch in der Reihenfolge der Fehler. F-464 ist bestätigt.
- **V-81:** Der Kommentar entspricht §6. Die Invariante stimmt mit dem Code überein, und der Zweig ist unverändert.
- Die Abdeckungs-Deklarationen sind gleich, und `make gates` ist grün.
- F-460 ist bis auf ein Wort umgesetzt (V-89), F-463 ist umgesetzt. Die Mutation `> 0` → `>= 0` ist rot, und der gestrichene Satz über den Server ist in Ordnung.
- Die Adresse `slice-tests-ueberlebende-mutanten` nimmt alle vier Mutanten an und bleibt nach diesem Slice offen. Y2 nimmt sie unter einer offenen Randform an.

Kein Befund blockiert die DoD. V-89 sollte vor der Closure einen Ausgang bekommen: eine engere Fassung des Kommentars durch den Implementer oder eine ausdrückliche Übernahme in den Nehmer. V-90 geht an den Planner, V-91 und V-92 an die Closure. Abschnitt 6 enthält Vorschläge für die Risiko-Ausgänge.
