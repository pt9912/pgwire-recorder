# Verifikation: slice-extended-query-modell — 2026-10-04

**Rolle:** Verifier (Modul 11). Geprüft wird, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules prüft der Reviewer.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-extended-query-modell.md` <!-- d-check:ignore (Pfad beim geprüften Stand) --> (Kopf, §1 bis §3, §6, §8) gegen die Commits `cd03e09` (Umsetzung) und `6b4967d` (Review-Findings) bei HEAD `6b4967d`.

**Eingang:** die DoD-Liefer-Punkte 1 bis 3 und die Commit-Message von `6b4967d` („F-283 bis F-293 behoben“). Dazu der Review-Report `2026-10-04-review-slice-extended-query-modell.md` (F-283 bis F-293, Stand `cd03e09`) und die Mutationstabelle des Implementers `2026-10-04-mutationen-slice-extended-query-modell.md` (M1 bis M11, N1 bis N10, Abschnitt „Nicht abgedeckt“). Laut Tabelle hat der Implementer die Mutationen im Arbeitsbaum angewandt und zurückgeschrieben. Der Arbeitsbaum war bei meinem Start sauber.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-04

Die Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg. Mutationen und Sonden liefen nur in einer Kopie von `git archive HEAD` im Scratchpad, nie im Repo:

- je Mutation eine eigene Kopie und `docker build --target test` mit dem Tag `pgr-verif:mut`
- die Sonden als zusätzliche Testdatei, schreibgeschützt in einen Container des Images `pgr-verif:source` gemountet, netzlos

Images und Kopien habe ich danach gelöscht. Die Images `pgwire-recorder:*` erzeugen die Make-Ziele des Repos; ihre IDs sind unverändert. `git status --short` war vor dem Anlegen dieser Datei leer.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — LH-FA-18: `LH-FA-18.a` und `SPEC-041` stimmen mit dem Domain-Modell überein; die Entscheidung steht in ADR-0012. **Bestätigt.**

Abgleich von Spezifikation und `internal/hexagon/model`:

| Spec-Aussage | Modell | Urteil |
|---|---|---|
| Client-Nachrichten `Parse` … `Sync` (7) | `clientMessages`, 7 Konstanten | gleich |
| Server-Nachrichten einschließlich `ParameterStatus` (14) | `extendedResponses`, 14 Einträge | gleich |
| Gruppe endet mit `Flush` oder `Sync`, zuerst Client-, dann Server-Nachrichten | `Group{Client, Server}`; `validate`: `flush`/`sync` genau als letzte Client-Nachricht | gleich |
| jedes `Sync` erzeugt genau ein `ReadyForQuery`; die Interaktion endet mit dem `ReadyForQuery` auf ihr `Sync` | nur die letzte Gruppe endet mit `sync`, ihre letzte Server-Nachricht ist `ready_for_query`, keine andere ist es | gleich |
| Felder je Client-Nachricht (Tabelle in `SPEC-041`), `describe`/`close` mit Zielart | `ClientMessage` mit allen Feldern; Zielart geprüft | gleich |
| `parameter_description` trägt `param_types` | `Response.ParamTypes` | gleich |

[ADR-0012](../plan/adr/0012-extended-query-gruppen.md) ist `Accepted` und unverändert. Ihre Entscheidung („Gruppen je `Flush`/`Sync`“) bildet das Modell ab. Ihr `Schärft:` nennt `LH-FA-18.a` und `SPEC-041`, also die beiden Stellen, die dieser Slice abgeglichen hat.

### Punkt 2 — Request- und Response-Typen und ihre Formregeln im Modell, je Regel ein Negativfall. **Bestätigt.**

Den Modell-Code hat `6b4967d` nicht geändert, nur die Tests. Die Mutationen des Reviews (M3, M5, M7, M9, M11, M13 rot am Stand `cd03e09`) gelten also weiter für die Regeln von `Validate`. Ich habe nachgeprüft, was das Review grün gesehen hatte, und eine Regel zusätzlich:

| # | Mutation (Kopie, Stand `6b4967d`) | Ergebnis |
|---|---|---|
| N2 | `error_response` aus `simpleResponses` | **rot:** `TestValidateGueltig` |
| N3 | `empty_query_response` aus `simpleResponses` | **rot:** `TestValidateGueltig` |
| N4 | `parameter_status` aus `simpleResponses` | **rot:** `TestValidateGueltig` |
| X2 | Pflicht „letzte Gruppe endet mit `ready_for_query`“ abgeschaltet | **rot:** `TestValidateFehler` (2 Fälle), `TestUnmarshalExtendedFehler` (2 Fälle) |

Core-Reinheit nach [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md): `extended.go` importiert nur `errors` und `fmt`, `response.go` und `recording.go` nichts; keine YAML-Tags und keine Bibliothekstypen. `make a-check` meldet 0 Befunde.

### Punkt 3 — LH-FA-07: Der Recording-Adapter schreibt und liest `type: extended` nach `SPEC-041`; eine abgeschnittene oder fremd geformte Extended-Interaktion ist beschädigt (`PGR-E3003`). **Im Kern bestätigt, mit einem Rest gegen den Wortlaut von `SPEC-041` (V-16).**

- **Roundtrip:** `TestExtendedRoundtrip` (alle Felder, NULL- und Binärparameter, `parameter_description`, Flush-Gruppe ohne Server-Nachricht, deterministisch beim zweiten Schreiben) und `TestExtendedLeereFelder` (leere Felder werden ausgeschrieben, auch `param_types: []` an `parameter_description`).
- **Beispiel aus `SPEC-041`:** `TestUnmarshalSpec041` liest den Block wörtlich, mit einer einfachen Anfrage davor.
- **Beschädigt:** `TestUnmarshalExtendedFehler` mit 31 Fällen. Jede Meldung nennt die Regel.

Die zweite Runde der Mutationstabelle (N1 bis N10) und vier eigene Mutationen habe ich in Kopien wiederholt. Jede ist **rot**, und es schlägt der Fall fehl, den die Tabelle nennt:

| # | Mutation | Rot in |
|---|---|---|
| N1 | Typprüfung des Matchers entfernt | `TestReplayExtendedAmCursor` |
| N5 | fehlendes Feld einer Client-Nachricht angenommen | 6 Fälle („parse ohne sql“ … „close ohne target“) |
| N6 | `null` als belegt gezählt | „parse mit sql: null“ |
| N7 | `param_types` an jeder Antwort angenommen | „param_types an parse_complete“, „… in einfacher Anfrage“ |
| N8 | `param_types` an `parameter_description` nicht verlangt | „parameter_description ohne param_types“ |
| N9 | leere `param_types` nicht geschrieben | `TestExtendedLeereFelder` |
| N10 | Vorprüfung auf `null` abgeschaltet | „request: null“, „responses: null“, „groups: null“ |
| X1 | Vorprüfung nur noch für `request` | „responses: null“, „groups: null“ |
| X3 | `type: query` auf Interaktionsebene als einfache Anfrage angenommen | „type: query auf Interaktionsebene“ |
| X4 | einfache Anfrage mit `groups` angenommen | „einfache Anfrage mit groups“ |
| X5 | Extended-Interaktion mit `responses` angenommen | „extended mit responses“ |

Der Rest steht in Abschnitt 3 (Sonden): `param_types: null` an einer anderen Server-Nachricht wird angenommen, obwohl `SPEC-041` ihn wörtlich beschädigt nennt (V-16). Zwei weitere Randformen sind in der Spezifikation nicht oder abweichend geregelt (V-17, V-18).

### Punkt 4 — `make gates` grün. **Bestätigt.**

Eigener Lauf von `make gates` am Stand `6b4967d`, Exit 0. Im Log:

- `abdeckung-check` und `abdeckung-gegenprobe` grün
- `a-check` 0 Befunde (mit dem bekannten Hinweis auf `test/integration/*`), `a-check-negativ` grün
- `baseline-verify` v6.13.0 OK (54 Dateien)
- `build` grün; `test` aus dem Cache des Repos
- `docs-check`: 140 Dateien, 0 Befunde; `commit-msg-gegenprobe` grün
- `test-integration` grün in beiden Phasen

Weil `test` aus dem Cache kam, habe ich die Stufe `test` in der Kopie ohne Cache neu gefahren (`--no-cache-filter test`): grün. Zusätzlich lief `go test -v` auf `TestUnmarshalExtendedFehler`, `TestReplayExtendedAmCursor`, `TestValidateGueltig`, `TestExtendedLeereFelder`, `TestUnmarshalSpec041` und `TestExtendedRoundtrip`: alle PASS.

### Prozess-Punkte der DoD (nicht Teil des Auftrags, nur vermerkt)

- Review: Der Report liegt vor, er erfasst `cd03e09`. Den Korrektur-Commit `6b4967d` prüft dieser Beleg funktional (Abschnitt 2). Ein Folge-Review gibt es nicht. Der Commit fügt aber einen neuen Lese-Mechanismus hinzu (`formSchluesselNull`, zweiter Parse über `yaml.Node`); ob dafür ein Folge-Review nötig ist, entscheidet der Planner.
- Closure-Notiz, Beobachtungs-Register, Ausgänge der Risiken aus §6 und die drei Paarungen: offen, wie bei `in-progress` zu erwarten. §6 trägt ein Risiko doppelt (V-20).

---

## 2. Stand der Review-Findings (F-283 bis F-293)

Geprüft am Code von `6b4967d`.

| ID | Stand | Beleg |
|---|---|---|
| F-283 | **getragen** | `TestReplayExtendedAmCursor` prüft eine leere und eine nicht leere einfache Anfrage gegen eine Extended-Interaktion am Cursor, dazu, dass der Cursor stehen bleibt (Warnung beim Schließen). N1 ist rot. |
| F-284 | **getragen** | `TestValidateGueltig`, Fall „query mit jedem Antworttyp“, enthält alle acht Typen. N2, N3, N4 sind rot. Die zwei übrigen Mitglieder (`row_description`, `notice_response`) waren schon im Review rot (M16, M18). |
| F-285 | **getragen** | `SPEC-041` nennt jetzt ein fehlendes Feld und ein Feld mit `null` beschädigt; der Leser verlangt jeden Schlüssel des Typs mit Wert. N5, N6 sind rot, Sonde P10 (`params: null`) wird abgelehnt. |
| F-286 | **getragen** | `SPEC-001` sagt jetzt, wo `type` je Art steht, und `SPEC-041` stellt `offset_ms` nur bei Extended neben `type`. Der Test „type: query auf Interaktionsebene“ hält die Leser-Entscheidung fest (X3 rot). Rest: V-18 (`type: ''`). |
| F-287 | **teilweise** | `param_types` an einer anderen Antwort ist abgelehnt (N7 rot), aber nicht mit dem Wert `null` (V-16). |
| F-288 | **getragen** | Der Schreiber gibt an jeder `parameter_description` `param_types` aus, auch `[]`; `SPEC-041` sagt „auch leer“. N9 ist rot. |
| F-289 | **getragen** | `formSchluesselNull` lehnt `request`, `responses` und `groups` mit `null` ab, und der Kommentar nennt die Prüfung. N10, X1 sind rot. Der Kommentar von `interactionFromDTO` sagt weiter mehr zu, als geprüft wird, nämlich bei `type` (V-18). |
| F-290 | **nicht wiederholt** | Die Message von `cd03e09` bleibt in der Historie. Die Message von `6b4967d` nennt keine `SPEC-*`- oder `ARC-*`-Kennung. |
| F-291 | **getragen** | §3 hat Zeilen für `slice-extended-query-replay.md`, `replay_test.go`, die beiden Reports und die beiden Register-Belege. |
| F-292 | **getragen** | `slice-extended-query-record` §6 übernimmt das Risiko. Im eigenen Plan steht es zweimal (V-20). |
| F-293 | **getragen** | Die Mutationstabelle liegt als Artefakt vor. Ich habe alle zehn N-Mutationen nachgefahren: jede ist rot in dem genannten Test. Die Aussage „je Regel ist die rote Mutation gesehen“ ist für die Regeln, die die Tabelle nennt, belegt. Für die Lücken unter „Nicht abgedeckt“ gilt sie nicht, und die Tabelle sagt das auch. |

---

## 3. Die Präzisierungen der Spezifikation in Leser, Schreiber und Tests

Sonden in der Kopie (`Unmarshal`/`Marshal` direkt, Ergebnis wörtlich):

| # | Eingabe | Ergebnis | Gegen die Spezifikation |
|---|---|---|---|
| P1a | `{type: parse_complete, param_types: null}` in einer Extended-Gruppe | **angenommen** | `SPEC-041`: „steht es an einer anderen Server-Nachricht, ist das Recording beschädigt“ → **Abweichung** (V-16) |
| P1b | `{type: command_complete, tag: x, param_types: ~}` in einer einfachen Anfrage | **angenommen** | dito (V-16) |
| P2 | `{type: parameter_description, param_types: null}` | `PGR-E3003`: `Schlüssel "param_types" fehlt an parameter_description` | konform |
| P3 | `client: null` | `PGR-E3003`: `Gruppe 1: ohne Client-Nachricht` (über `Validate`) | Ergebnis konform, die Meldung nennt `null` nicht |
| P8 | Gruppe ohne `client` | `PGR-E3003`: `Gruppe 1: ohne Client-Nachricht` | dito |
| P4a | Flush-Gruppe mit `server: null`, danach eine Sync-Gruppe | **angenommen**, als leere Liste | `SPEC-041` regelt `null` an `server` nicht (V-17) |
| P4b | Flush-Gruppe ohne den Schlüssel `server` | **angenommen**, als leere Liste | dito; nicht in „Nicht abgedeckt“ genannt (V-17) |
| P4c | Flush-Gruppe mit `server: []` (Gegenstück) | angenommen | konform |
| P5 | Sync-Gruppe mit `server: null` | `PGR-E3003`: `letzte Gruppe endet nicht mit ready_for_query` | konform |
| P6 | einfache Anfrage mit `type: null` an der Interaktion | **angenommen** | `SPEC-001`: „eine einfache trägt dort kein `type`“; `null` ist im Modell „fehlt“ — vertretbar, aber nicht geregelt (V-18) |
| P7 | einfache Anfrage mit `type: ''` an der Interaktion | **angenommen** | `SPEC-001`: „Ein anderer Wert von `type` an einer dieser Stellen macht das Recording zu einem beschädigten“ → **Abweichung** (V-18) |
| P9 | fremder Schlüssel an einer Gruppe | `PGR-E3003`: `field extra not found` | konform |
| P10 | `bind` mit `params: null` | `PGR-E3003`: `Schlüssel "params" fehlt an bind` | konform |
| P11 | `type: extended`, `groups: []` | `PGR-E3003`: `ohne Gruppe` | konform |
| P12 | `Marshal` eines `sync` mit `SQL` und eines `parse_complete` mit `ParamTypes` | schreibt ohne Fehler, beide Felder fallen weg; Roundtrip ungleich | Schreiber hält `SPEC-041` ein; das stille Wegfallen am `parse_complete` steht in keinem Risiko (V-21) |

**Ergebnis je Präzisierung:**

- **Fehlendes Feld / `null` = beschädigt (Client-Nachricht):** im Leser (`clientMessageDTO.UnmarshalYAML`, `belegt`), im Schreiber (alle Felder des Typs, auch leer) und in Tests (6 Fälle, N5/N6 rot). Konform.
- **`type: query` nur in `request`:** im Leser (`interactionFromDTO`, `default`-Zweig) und im Test (X3 rot). Konform. Der zweite Teil desselben Satzes in `SPEC-001` („eine einfache trägt dort kein `type`“, „ein anderer Wert … beschädigt“) ist für `type: ''` nicht umgesetzt (V-18).
- **`param_types` nur und immer an `parameter_description`:** Für „immer“ sind Schreiber und Leser konform (N8, N9 rot, P2). Für „nur“ ist der Leser konform, solange ein Wert dasteht (N7 rot), mit `null` nicht (V-16).
- **`request`/`responses`/`groups` mit `null` = beschädigt:** im Leser (`formSchluesselNull`) und in Tests (N10, X1 rot). Konform.

### Bewertung der zwei Lücken unter „Nicht abgedeckt“ gegen SPEC-041

1. **`param_types: null` an einer anderen Antwort angenommen.** Das ist keine Testlücke, sondern ein Verhalten, das dem Wortlaut von `SPEC-041` widerspricht: Der Schlüssel *steht* an einer anderen Server-Nachricht. Die Spezifikation sagt nicht „mit einem Wert“, und für Client-Nachrichten nennt derselbe Absatz `null` ausdrücklich beschädigt. Die Ursache ist die von F-289: Der Decoder liest `null` in einen Zeiger als fehlend. Im Repo liegt mit `formSchluesselNull` schon das Muster, das die Lücke schließt. Die Folge ist gering, denn das Modell erhält keinen Wert. Die DoD verlangt aber „nach `SPEC-041`“, und auch die Abdeckungszeile LH-FA-07/Negative sagt ohne Einschränkung „ebenso param_types an einer anderen Antwort als parameter_description“. **Urteil: Abweichung von `SPEC-041`, V-16 (MEDIUM).**
2. **`client: null`/`server: null` als leere Liste gelesen.** `SPEC-041` sagt „jede Gruppe hat die Client-Nachrichten (`client`) und die Server-Nachrichten (`server`)“. Seine Regel für `null` nennt nur `request`, `responses`, `groups` und die Felder einer Client-Nachricht. Für die Schlüssel einer Gruppe entscheidet die Spezifikation also nicht. Folgen im Leser:
   - `client: null` und ein fehlendes `client` enden über `Validate` doch als `PGR-E3003` (P3, P8). Nur die Meldung nennt die Ursache nicht.
   - `server: null` *und ein fehlender Schlüssel `server`* werden an einer Flush-Gruppe vor der letzten angenommen (P4a, P4b). Den fehlenden Schlüssel nennt die Tabelle nicht.
   - Eine abgeschnittene Datei fällt trotzdem immer auf. Wird sie in einer Flush-Gruppe abgeschnitten, ist diese die letzte Gruppe, und `Validate` verlangt dort `sync` mit `ready_for_query` (P5, „abgeschnitten nach der Flush-Gruppe“).

   **Urteil: keine Verletzung des Wortlauts, aber eine ungeregelte Leser-Entscheidung, dieselbe Klasse wie F-285. Das ist V-17 (LOW).**

---

## 4. Plan gegen Code-Diff

### §1 Ziel und Abgrenzung

| Plan | Code | Urteil |
|---|---|---|
| Modell trägt Extended-Typen und `Interaction.Validate` | `extended.go`, `response.go`, `recording.go` | konform |
| Spezifikation abgeglichen, Entscheidung als ADR angenommen | `LH-FA-18.a`, `SPEC-001`, `SPEC-041` geändert; [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) `Accepted` | konform |
| Recording-Adapter schreibt und liest `type: extended` | `yaml.go` | konform, Rest V-16 |
| Replay-Service im Code unverändert, ein Test hält die Typprüfung | `replay.go` nicht im Diff; `TestReplayExtendedAmCursor` | konform |
| keine PGWire-Adapter, kein SQLite, keine Abdeckung für LH-FA-18 | nicht im Diff; die Abdeckung trägt nur LH-FA-07/Negative | konform |

### §3 Dateien

Jede Datei aus `git show --stat cd03e09 6b4967d` hat eine Zeile in §3, und jede Zeile in §3 hat Änderungen im Diff. `internal/hexagon/model/doc.go` (eine Zeile) fällt unter die Zeile `internal/hexagon/model`. **Konform.**

### §6 Risiken

| Risiko | Stand im Code | Ausgang bei Closure |
|---|---|---|
| 1 — Lücken in `SPEC-041` | vier geschlossen; zwei Randformen bleiben (V-16 Code, V-17 Spec) | erst nach V-16/V-17 setzbar |
| 2 — Leser lehnt `extended` ab | Leser und Schreiber liefern `extended` | „eingetreten und gelöst“ |
| 3 — nil gegenüber leerer Liste | an `slice-extended-query-replay` übergeben | übergeben |
| 4 — `Validate` prüft keine Felder je Typ | an `slice-extended-query-record` übergeben; ohne `Response.ParamTypes` (V-21) | übergeben |
| 5 und 6 — Schreiber ohne `Validate` | **doppelt** (V-20); an `slice-extended-query-record` übergeben | übergeben |

### §8 Register

Die Belege für `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` liegen. V-18 ist ein weiterer Fall der zweiten Klasse, diesmal im Korrektur-Commit (V-19).

---

## 5. ADR-Konformität

| ADR | Zusage | Urteil |
|---|---|---|
| [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md) | Bibliothekstypen werden an der Adaptergrenze in Domain-Typen übersetzt | **konform.** Das Modell ist frei von YAML und PGWire; DTOs und die Übersetzung liegen in `recording`. `make a-check` 0 Befunde, `a-check-negativ` grün. |
| [ADR-0007](../plan/adr/0007-strict-replay.md) | exaktes sequenzielles Matching | **konform.** Der Matcher vergleicht Art und SQL. Eine Extended-Interaktion am Cursor ist für jede einfache Anfrage eine Abweichung, und der Cursor bleibt stehen (`TestReplayExtendedAmCursor`, N1 rot). |
| [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) | Gruppen je `Flush`/`Sync`, zuerst Client-, dann Server-Nachrichten | **konform.** Abschnitt 1, Punkt 1. Der Feldvergleich im Replay gehört zu `slice-extended-query-replay`. |

---

## 6. Spec-Konformität

| Stelle | Urteil |
|---|---|
| [`LH-FA-18`](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol) | unverändert. Kein Akzeptanzkriterium ist in diesem Slice belegt und keines deklariert, wie §1 es begründet. |
| [`LH-FA-07`](../../spec/lastenheft.md#lh-fa-07--persistente-recordings) Negative | **konform** für abgeschnittene Dateien und alle Fehlformen mit Test. Ausnahmen: V-16, V-18. Die Abdeckungszeile nennt mehr, als geprüft wird (V-19). |
| `LH-FA-18.a` | **konform** (Abschnitt 1). |
| `SPEC-001` | **konform** für `type: query` in `request` und `type: extended` an der Interaktion. **Abweichung** bei `type: ''` an einer einfachen Interaktion (V-18). |
| `SPEC-002` | **konform.** Die Ausgabe einfacher Anfragen ist laut Review byteidentisch zu `cd03e09~1`; den Schreibpfad für einfache Anfragen hat `6b4967d` nicht geändert, `param_types` steht dort nie. |
| `SPEC-041` | **konform** für Form, Feldtabelle, leere Felder, fehlende/`null`-Felder an Client-Nachrichten, `param_types` an `parameter_description` und `null` an den Form-Schlüsseln. **Abweichung** bei `param_types: null` an anderen Server-Nachrichten (V-16). **Ungeregelt:** `null`/fehlend an `client`/`server` (V-17). |
| `ARC-001` | **konform.** Die Zeile nennt Extended-Gruppe, Client-Nachricht und Formregeln, ohne ADR, Slice oder Welle (Hard Rule 3.4). |

---

## 7. Befunde des Verifiers

| ID | Kategorie | Befund | Beleg |
|---|---|---|---|
| V-16 | MEDIUM | `param_types: null` an einer anderen Server-Nachricht als `parameter_description` wird angenommen, in Extended-Gruppen wie in einfachen Anfragen. `SPEC-041` nennt das beschädigt („steht es an einer anderen Server-Nachricht“), und DoD-Punkt 3 verlangt „nach `SPEC-041`“. Die Mutationstabelle führt den Fall unter „Nicht abgedeckt“, als wäre er eine Testlücke. Es ist aber ein Verhalten gegen den Wortlaut, und die Abdeckungszeile LH-FA-07/Negative sagt es ohne Einschränkung zu. Ursache wie bei F-289: Der Decoder liest `null` in `*[]uint32` als fehlend. Vorschlag: die Anwesenheit des Schlüssels am `yaml.Node` prüfen, wie es `formSchluesselNull` für die Form-Schlüssel tut. Den Gegenfall gibt es auch: Die Spezifikation erlaubt `null` dort ausdrücklich, dann gehört das in `SPEC-041` mit Historienzeile. | Sonden P1a, P1b; `internal/adapters/driven/recording/yaml.go` (`responseFromDTO`, `responseDTO.ParamTypes`) |
| V-17 | LOW | Für `client`/`server` einer Gruppe regelt `SPEC-041` weder `null` noch das Fehlen. Der Leser nimmt `server: null` und den fehlenden Schlüssel `server` an einer Flush-Gruppe vor der letzten als leere Liste an (P4a, P4b). `client: null` endet nur über `Validate` als beschädigt, mit einer Meldung, die `null` nicht nennt (P3). Abgeschnittene Dateien fallen trotzdem auf. Den fehlenden Schlüssel `server` nennt „Nicht abgedeckt“ nicht. Gleiche Klasse wie F-285. Frage an den Planner: Gehören `client`/`server` in die Regel für `null` und fehlende Schlüssel, oder sagt `SPEC-041`, dass ein fehlendes oder `null`-`server` leer bedeutet? Der Leser folgt der Entscheidung, ein Test hält sie fest. | Sonden P3, P4a, P4b, P5, P8 |
| V-18 | LOW | Eine einfache Interaktion mit `type: ''` (und mit `type: null`) wird als einfache Anfrage angenommen. `SPEC-001`, in diesem Slice für F-286 präzisiert, sagt: „eine einfache trägt dort kein `type`. Ein anderer Wert von `type` an einer dieser Stellen macht das Recording zu einem beschädigten.“ `''` ist ein anderer Wert. Die Kommentare von `interactionDTO` („Ein anderer Wert von type ist beschädigt“) und `interactionFromDTO` sagen dasselbe zu, geprüft wird nur „leer oder `extended`“. Ursache: Das Feld `Type string` unterscheidet „fehlt“ nicht von `''`/`null`. Vorschlag: `type` in die Anwesenheitsprüfung am `yaml.Node` aufnehmen, mit Negativfall. | Sonden P6, P7; `yaml.go` (`interactionDTO`, `interactionFromDTO`) |
| V-19 | LOW | Mehrere Zusagen gehen weiter als die Prüfung. Der Kommentar von `TestUnmarshalExtendedFehler` und die Abdeckungszeile LH-FA-07/Negative in `docs/user/abdeckung-unit.md` sagen „einen … mit null belegten Schlüssel“ und „param_types an einer anderen Antwort“ ohne Einschränkung zu (V-16, V-17), die Kommentare zu `type` ebenso (V-18). Damit liegt die Register-Klasse `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` im Korrektur-Commit desselben Slice erneut vor, nachdem F-289 sie schon gemeldet hatte. Für den Beleg im Register und die Frage aus dem Review an den Architect (Regel „Prüfung je Zusage“). | `docs/user/abdeckung-unit.md` Zeile LH-FA-07/Negative; `yaml_test.go` Kommentar von `TestUnmarshalExtendedFehler` |
| V-20 | LOW | Plan §6 trägt das Risiko „Der Schreiber prüft nicht mit `Validate` … (Review F-292)“ zweimal wortgleich; die Doppelung kam mit `6b4967d`. Die Closure verlangt je Risiko genau einen Ausgang. Die Doppelung macht daraus zwei Einträge für ein Risiko. Eine Zeile streichen. | `git show 6b4967d -- docs/plan/planning/in-progress/slice-extended-query-modell.md` <!-- d-check:ignore (Pfad beim geprüften Stand) --> |
| V-21 | INFO | Der Schreiber lässt `ParamTypes` an jeder Antwort außer `parameter_description` still weg (P12, Roundtrip ungleich). Das ist dieselbe Grenze wie Risiko 4 für Client-Nachrichten. Risiko 4 und `slice-extended-query-record` §6 nennen aber nur die Client-Nachrichten. Hinweis an den Planner: das übergebene Risiko auf `Response.ParamTypes` erweitern. | Sonde P12; `yaml.go` (`responseToDTO`) |

---

## 8. Gegenprobe dieser Datei

`make docs-check` nach dem Anlegen dieser Datei: `d-check: 141 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

## 9. Gesamturteil

**DoD-Liefer-Punkte 1, 2 und 4 sind bestätigt. Punkt 3 ist im Kern bestätigt, ein Rest steht gegen den Wortlaut von `SPEC-041` (V-16, MEDIUM).** Grundlage sind eigene Läufe:

- `make gates` mit Exit 0
- die Stufe `test` ohne Cache in der Kopie
- 15 Mutanten (N1 bis N10, X1 bis X5), jeder rot im erwarteten Test
- 15 Sonden zu den Randformen

Die Review-Findings F-283 bis F-286, F-288, F-289, F-291 bis F-293 trägt `6b4967d`. F-287 ist nur für Werte getragen, nicht für `null` (V-16). F-290 ist nicht wiederholt.

Der Code hält [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md), [ADR-0007](../plan/adr/0007-strict-replay.md) und [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) ein. Er entspricht `LH-FA-18.a`, `SPEC-002` und `ARC-001`. `SPEC-001` und `SPEC-041` hält er mit den Ausnahmen V-16 und V-18 ein, dazu kommt die ungeregelte Randform V-17.

**Kein Befund ist HIGH.** Vor der Closure empfohlen:

1. V-16: `param_types: null` an anderen Server-Nachrichten ablehnen (Negativfall), oder `SPEC-041` ändern — Implementer.
2. V-18: `type: ''`/`null` an einer einfachen Interaktion ablehnen, oder `SPEC-001` ändern; die Kommentare an die Prüfung angleichen — Implementer.
3. V-17: `SPEC-041` für `client`/`server` entscheiden, der Leser folgt — Planner.
4. V-19: Abdeckungszeile und Test-Kommentar an die Prüfung angleichen; Beleg im Register — Implementer.
5. V-20: die doppelte Zeile in Plan §6 streichen; V-21: das übergebene Risiko in `slice-extended-query-record` erweitern — Planner.
6. Über ein Folge-Review für `6b4967d` entscheiden (neuer Lese-Mechanismus `formSchluesselNull`) — Planner.
