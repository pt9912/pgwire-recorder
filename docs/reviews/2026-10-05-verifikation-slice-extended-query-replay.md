# Verifikation: slice-extended-query-replay — 2026-10-05

**Rolle:** Verifier (Modul 11). Geprüft wird, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules prüft der Reviewer.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-extended-query-replay.md` (Kopf, §1 bis §4, §6) gegen den Gesamt-Diff `d0e5d3b..f528c90` bei HEAD `f528c90`. Code-Commits: `b895565` (Umsetzung), `09405f0`, `c8ebf08` (F-321 bis F-329), `b697462` (F-330 bis F-338), `2d3e0f7` (F-339 bis F-344), `f528c90` (F-347, F-348). `f528c90` hat noch kein Review gesehen; er ist hier mit eigenen Mutationen und Sonden geprüft.

**Eingang:**

- die DoD-Liefer-Punkte 1 bis 3 und die Commit-Messages
- der [Review-Report](2026-10-05-review-slice-extended-query-replay.md) (F-321 bis F-329, Stand `09405f0`)
- das [Folge-Review](2026-10-05-folge-review-slice-extended-query-replay.md) (F-330 bis F-338)
- das [zweite Folge-Review](2026-10-05-folge-review-2-slice-extended-query-replay.md) (F-339 bis F-346, darin F-346 an den Verifier)
- das [dritte Folge-Review](2026-10-05-folge-review-3-slice-extended-query-replay.md) (F-347, F-348)
- der Hinweis des Implementers, dass er die Mutation M14 an Abnahmeszenario 7 nicht selbst rot gesehen hat, weil der Runner nach der ersten Phase abbricht

Der Arbeitsbaum war bei meinem Start sauber.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

Die Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg. Mutationen und Sonden liefen nur in Kopien von `git archive f528c90` im Scratchpad:

- Unit-Mutationen: je Mutation eine eigene Kopie, `go test -count=1` des betroffenen Pakets, netzlos in einem Container des Images `pgr-verif:source` (Stufe `source` des `Dockerfile`)
- E2E: Images `pgr-verif:int-<name>` (Stufe `integration`) aus der unveränderten Kopie und aus sechs Mutanten, eigenes internes Docker-Netz `pgr-verif-net`, eigenes Volume `pgr-verif-daten` und eigener PostgreSQL-Container aus dem gepinnten Image von `harness/mk/integration.mk`
- Race-Detector: Container aus `pgr-verif:source` mit `build-base` (Netz nur für `apk`), kein Gate

Images, Container, Netz, Volumes und Kopien habe ich danach gelöscht. `git status --short` war vor dem Anlegen dieser Datei leer.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), drei Behauptungen (F-346). **Bestätigt.**

**1a — Abnahmeszenario 7: mit gestoppter PostgreSQL liefert Replay dem Go-Client mit Prepared Statements dasselbe Ergebnis. Bestätigt.**

- `TestE2EOhnePostgresExtendedReplay` ist grün im eigenen `make gates`, in der zweiten Phase des Runners (`PGR_OHNE_POSTGRES=1`). Der Test prüft vorher, dass PostgreSQL nicht erreichbar ist.
- **Bewusstes Brechen in Isolation.** Der Runner bricht nach einer roten ersten Phase ab, und jeder Mutant am Replay macht schon `TestE2EReplayExtendedPgx` in der ersten Phase rot. Deshalb habe ich die Phasen getrennt:
  1. erste Phase nur mit `TestE2EVorbereitungExtendedOhnePostgres` aus dem **unveränderten** Image; die Aufzeichnung (13 795 Bytes) und die Sicht (13 Zeilen, darunter `fehler: 22012 …` und `batch: 22012 …`) liegen danach im Volume
  2. PostgreSQL gestoppt
  3. zweite Phase nur mit `^TestE2EOhnePostgresExtendedReplay$`, je Image gegen dieselbe Aufzeichnung

| Image | Mutation | Ergebnis |
|---|---|---|
| unverändert | — | grün (zweimal, vor und nach den Mutanten) |
| A1 | **Zustand ohne den Fix:** `replaySitzung` lehnt jede Extended-Nachricht mit `PGR-E6001` ab, wie vor dem Slice | **rot:** Das erste `QueryRow` scheitert, die Verbindung ist danach geschlossen, `extendedAblauf` bricht bei `Zeilen: conn closed` ab |
| U4 | Server-Nachrichten einer Gruppe nach jeder Client-Nachricht, nicht erst nach Flush/Sync | rot: `received unexpected message type *pgproto3.BindComplete` |
| U12 | Replay verschluckt die `ErrorResponse` einer Gruppe | rot: Sicht weicht ab, `fehler: timeout: context deadline exceeded` statt `22012 division by zero` |
| U1 | Matcher vergleicht keine Parameterwerte | **grün.** Erwartet: Der Ablauf sendet dieselben Werte in derselben Reihenfolge. Diese Klasse fängt `TestE2EReplayExtendedAbweichung` (unten, 2a). |

Der Test ist damit für den Zustand vor dem Slice aus dem richtigen Grund rot: Der Client erhält keine Antwort auf das Extended Query Protocol. Die Lücke aus M14 ist geschlossen. Zur Beweisführung siehe V-30.

**1b — Beim Herunterfahren beantwortet Replay eine begonnene Extended-Interaktion bis zu ihrem `Sync` und endet danach ([LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus)). Bestätigt.**

| # | Mutation | Rot in |
|---|---|---|
| A3 | **vorige Logik:** endet ctx, kehrt `replaySitzung` vor dem nächsten Lesen zurück, ohne den Use Case zu fragen | `TestReplayHerunterfahren/laufende_Interaktion`, `TestReplayHerunterfahrenSpaeteFrist` (beide `unexpected EOF`); E2E `TestE2EReplayExtendedSigtermMittenInFolge` (`Ausführung 1 nach SIGTERM: unexpected EOF`) |
| A4 | Use Case wird nur einmal gefragt; nach dem `Sync` endet die Session nicht („endet danach“) | `TestReplayHerunterfahren`, `TestReplayHerunterfahrenSpaeteFrist` (Verbindung nach dem Sync nicht geschlossen); E2E `…SigtermMittenInFolge` (`replay endet nach dem Sync nicht von selbst`, 10 s) |
| A5 | `<-geweckt` entfernt (M1 der Folge-Reviews) | `TestReplayHerunterfahrenSpaeteFrist` |
| U7 | `Shutdown` gibt das Ende immer frei | `TestReplayExtendedHerunterfahren` (`nach Nachricht 1 (parse): Ende freigegeben true`) |

Der E2E-Test prüft dazu Exit-Code 0 und das Fehlen von `PGR-W2001`; unverändert ist er grün.

**1c — Eine Aufzeichnung, deren Interaktion die Form verfehlt, ist beim Start `PGR-E3003`. Bestätigt.**

- U6 (`Validate` in `NewReplayService` wirkungslos) ist rot in `TestReplayStartformExtended`.
- Sonde auf Prozessebene: Die Aufzeichnung aus 1a habe ich so verändert, dass die letzte Gruppe der ersten Interaktion mit `flush` endet. Ergebnis: `Recording [PGR-E3003]: Session 1, Interaktion 1: Gruppe 1: sync steht genau am Ende der letzten Gruppe, flush am Ende jeder anderen`, Exit-Code 3 ([`SPEC-025`](../../spec/spezifikation.md)).
- Diese Meldung stammt vom YAML-Leser, nicht vom Service. Siehe V-29.

### Punkt 2 — [LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), vier Behauptungen. **Nicht vollständig bestätigt** (2c, V-25).

**2a — Eine Extended-Abweichung wird erkannt und nicht beantwortet. Bestätigt.**

| # | Mutation | Rot in |
|---|---|---|
| U1 | Feld `params` nicht verglichen | `TestReplayExtendedWiederholung`, `…Felder`, `…Abweichung`; E2E `TestE2EReplayExtendedAbweichung` (`erhalten "eins!", <nil>`) |
| U2 | NULL gleich leerem Wert | `TestReplayExtendedFelder` |
| U3 | nil ungleich leerer Liste | `TestReplayExtendedFelder` |
| U5 | Cursor rückt bei einer Abweichung vor | `TestReplayExtendedAbweichung`, `…DiagnoseLebensdauer` |
| U8 | einfache Anfrage, wo Extended erwartet ist, wird nicht abgelehnt | `TestReplayExtendedFalscheArt`, `TestReplayExtendedAmCursor` |
| U14 | Extended-Nachricht, wo eine einfache Anfrage erwartet ist, wird nicht abgelehnt | `TestReplayExtendedFalscheArt` (Panik) |

Der Matcher vergleicht alle elf Felder von `model.ClientMessage` (Code-Lesung gegen `internal/hexagon/model/extended.go`), wie `LH-FA-18.a` §Replay sie aufzählt.

**2b — Eine aufgezeichnete `ErrorResponse` samt Verwerfen bis zum `Sync` wird reproduziert (Abnahmeszenario 6, Extended). Bestätigt.**

- U12 ist rot in `TestReplayExtendedFehlerantwort` und, isoliert, in Abnahmeszenario 7 (1a).
- Der pgx-Ablauf enthält einen Batch, in dem nach einem Fehler die folgende Anweisung bis zum `Sync` verworfen wird. Die Sicht beim Aufzeichnen trägt `batch: 22012 …` dreimal, die Sicht im Replay ebenso.

**2c — Die Diagnose nennt Stelle, Nachrichtentypen, Feld und die SQL-Texte der Anweisungen ohne Parameterwerte (§1). Bestätigt für Abweichungen innerhalb einer Extended-Interaktion, nicht für Abweichungen der Art (V-25).**

| # | Mutation | Rot in |
|---|---|---|
| U9 | Diagnose hängt die Parameterwerte an | `TestReplayExtendedAbweichung` (wörtlicher Vergleich der Meldung) |
| U10 | Bestätigungsfilter aus `f528c90` entfernt (F-347 zurückgedreht) | `TestReplayExtendedDiagnoseLebensdauer` |
| U11 | keine Anweisungen in der Diagnose | `…Abweichung`, `…DiagnoseAnweisung`, `…DiagnoseLebensdauer` |
| U13 | Ende der Transaktion beendet die Portale nicht | `TestReplayExtendedDiagnoseLebensdauer` |

Eigene Sonden in einer Kopie (zusätzliche Testdatei, nicht im Repo):

| # | Lage | Diagnose |
|---|---|---|
| P4a | erwartet `parse "SELECT $1::text"`, empfangen einfache Anfrage `SELECT 99` | `… Gruppe 1, Nachricht 1: erwartet Client-Nachricht parse, empfangen Anfrage "SELECT 99"`: **ohne das erwartete SQL** |
| P4b | erwartet einfache Anfrage `SELECT 7`, empfangen `parse "SELECT 42"` | `… Interaktion 1: erwartet Anfrage "SELECT 7", empfangen Client-Nachricht parse`: **ohne das empfangene SQL** |
| P4c | nach der letzten Interaktion `parse "SELECT 42"` | `… keine aufgezeichnete Interaktion mehr nach 1; empfangen Client-Nachricht parse`: **ohne das empfangene SQL** |
| P5 | Flush-Gruppe `parse s1`, deren `parse_complete` spät in der folgenden Sync-Gruppe steht; danach Abweichung an `bind s1` | `Anweisung erwartet unbekannt, empfangen "SELECT 2"`: **`s1` besteht, gilt aber als unbekannt** (V-26) |

**2d — Die Klasse ist 5, und die Verbindung endet. Bestätigt.**

- A2 (nach einer Abweichung weiterlesen statt schließen) ist rot in `TestReplayExtendedAbweichungImAdapter` (`Verbindung nach der Abweichung nicht geschlossen`).
- `TestE2EReplayExtendedAbweichung` prüft Exit-Code 5 und dass `anderer-wert` weder in der Fehlerantwort noch in stderr steht; grün im Gate-Lauf.

### Punkt 3 — `make gates` grün. **Bestätigt.**

- Eigener Lauf von `make gates` an `f528c90`, **Exit 0**:
  - `abdeckung-check` und `abdeckung-gegenprobe` grün
  - `a-check` 0 Befunde, `a-check-negativ` grün
  - `baseline-verify` v6.13.0 OK, 54 Dateien
  - `docs-check` 160 Dateien, 0 Befunde
  - `commit-msg-gegenprobe` grün
  - `test-integration` in beiden Phasen grün, darunter alle `TestE2EReplayExtended*`, `TestE2EVorbereitungExtendedOhnePostgres` und in der zweiten Phase `TestE2EOhnePostgresExtendedReplay`
- Die Stufe `test` kam im Gate-Lauf aus dem Cache. In der Kopie habe ich sie darum ohne Cache gefahren (`--no-cache-filter test`): `gofmt -l` leer, `go vet -tags integration` und alle Unit-Pakete `ok`.
- Race-Detector `go test -race -count=10 -run "Replay|Herunterfahren"` über `internal/adapters/driving/pgwire` und `internal/hexagon/services`: grün.

### Prozess-Punkte der DoD (nur vermerkt)

- **Review:** Reports zu `09405f0`, `c8ebf08`, `b697462` und `2d3e0f7` liegen vor. `f528c90` ändert `objekte` (Bestätigungsfilter) und hat kein Review. Dieser Bericht prüft ihn funktional (U10, P5). Ob ein viertes Folge-Review nötig ist, entscheidet der Planner.
- **Offen, wie bei `in-progress` zu erwarten:** Closure-Notiz, Register, Ausgänge der Risiken in §6 und die Paarungen. Sechs der acht Risiken tragen „offen bis Closure“; das ist kein Ausgang aus der geschlossenen Menge und muss bei der Closure ersetzt werden.

---

## 2. Plan gegen Code

### Kopf

| Feld | Urteil |
|---|---|
| `Bezug` | konform für die LH-Kennungen. Genannt ist nur [ADR-0007](../plan/adr/0007-strict-replay.md). [ADR-0012](../plan/adr/0012-extended-query-gruppen.md), deren Wiedergabe-Hälfte der Slice liefert, fehlt; [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) steht nur in §1 (V-28). |
| `Berührte Spec-Stellen` | konform. Jede genannte Stelle ist berührt (`ARC-003`: `ports/driving/replay.go`; `ARC-006`: `server.go`). `LH-FA-11.a` fehlt, obwohl DoD-Punkt 2 Abnahmeszenario 6 beansprucht (V-28). |

### §1 Ziel und Abgrenzung

| Plan | Code | Urteil |
|---|---|---|
| `ClientMessage` am Port, Vergleich in allen Feldern, leere Liste gleich nil | `ports/driving/replay.go`, `matcher.go` | konform (U1 bis U3) |
| Freigabe der Server-Nachrichten erst nach `Flush`/`Sync` der Gruppe | `ClientMessage` | konform (U4) |
| Cursor mit Interaktion, Gruppe, Nachricht (`ARC-002`) | `cursor` | konform |
| Diagnose: Stelle, Typen, Feld, bei SQL beide Texte, bei bind/execute/describe/close das SQL der Anweisung nach dem Verlauf | `anweisungen`, `objekte` | konform innerhalb einer Extended-Interaktion; Art-Abweichung ohne SQL der Extended-Seite (V-25) |
| „nur vom Server bestätigte `parse`, `bind` und `close` wirken“ | Bestätigungsfilter je Gruppe | konform, solange die Bestätigung in der Gruppe ihrer Nachricht steht; späte Bestätigung einer Flush-Gruppe wirkt nicht (V-26) |
| Grenze: SQL-Befehle, die Objekte beenden | Kommentar an `objekte`, §6 | konform (F-348) |
| Art-Abweichung (einfach ↔ Extended) ist `PGR-E5001` | `Query`, `ClientMessage` | konform (U8, U14) |
| `replaySitzung` leitet weiter und schließt nach einer Abweichung, ohne weiterzulesen | `server.go` | konform (A1, A2) |
| Herunterfahren bis zum `Sync`, danach Ende; `Shutdown` am Port | `server.go`, `Shutdown` | konform (A3, A4, A5, U7) |
| Frist nicht verdrahtet; Schließen der Verbindung endet regulär | `verbindungsende` → `return` | konform (Code-Lesung) |
| `Validate` beim Start, `PGR-E3003` | `NewReplayService` | konform (U6) |
| Ausschlüsse | kein `--fail-on-unconsumed`, kein `PGR-E5002` im Code; kein Zwei-Richtungs-Ablauf im Replay | konform |

### §3 Dateien — **Abweichung (V-27)**

Jede Zeile in §3 hat Änderungen im Diff. Umgekehrt fehlen Zeilen für die vier Review-Reports unter `docs/reviews/`, die im Slice-Diff mitcommittet sind, und für `internal/hexagon/services/replay_extended_test.go`. Die Zeile `internal/hexagon/services` nennt nur `matcher.go` und `replay.go`.

### §4 Trigger

Keine der beiden Rückführungs-Bedingungen ist eingetreten. Die Strict-Replay-Entscheidung ist unverändert ([ADR-0007](../plan/adr/0007-strict-replay.md)), und pgx im Standardmodus braucht nichts außerhalb von `LH-FA-18`. Konform.

### §6 Risiken

Jeder genannte Beleg-Test existiert und ist grün. Die Behauptungen „mit der vorigen Logik rot“ (Risiko Herunterfahren) und „ohne `<-geweckt` rot (M1)“ habe ich mit A3 und A5 nachgefahren. Das Risiko „Herunterfahren ohne Obergrenze“ hat schon den Ausgang *eingetreten* → `slice-v1-abschluss-betrieb`; dessen §1 nennt `replay`.

---

## 3. ADR-Konformität

| ADR | Zusage | Urteil |
|---|---|---|
| [ADR-0007](../plan/adr/0007-strict-replay.md) | exaktes sequenzielles Matching, keine Normalisierung, Mismatch-Diagnose nach Spezifikation | **konform** für das Matching (2a; Namen nicht normalisiert, Wiederholung nach Position). Folgepflicht „Diagnose gemäß Spezifikation“: Lücke V-25. |
| [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) | Vergleich jeder Client-Nachricht in allen Feldern; Server-Nachrichten einer Gruppe nach deren letzter Client-Nachricht | **konform** (U1 bis U4, isoliert auch in Szenario 7). Späte Antworten einer Flush-Gruppe werden mit der folgenden Gruppe freigegeben, wie `LH-FA-18.a` es sagt (Risiko F-308, offen). |
| [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) | gilt für den Record-Pfad; „Der Replay-Pfad übernimmt diesen Zuschnitt nicht automatisch“ | **konform.** `replaySitzung` liest und antwortet abwechselnd, wie §1 es abgrenzt. Die Entscheidung, ob die Session beim Herunterfahren endet, trifft der Use Case (`Shutdown`); der Adapter führt sie aus. `TestE2EReplayExtendedGegendruck` grün. |

---

## 4. Spec-Konformität

| Stelle | Urteil |
|---|---|
| `LH-FA-18.a` §Replay, §Fehler | **konform** (1a, 2a, 2b) |
| `LH-FA-18.a` §Mismatch | **konform.** Index von Gruppe und Nachricht, erwarteter und empfangener Typ. Parameterwerte erscheinen auf keiner Log-Stufe (strenger als verlangt). |
| `LH-FA-10.a` | **teilweise.** Stelle, Interaktionsnummer, kein Sprung. „Erwartete Query“ und „tatsächlich empfangene Query“ fehlen bei der Art-Abweichung und nach dem Ende der Aufzeichnung, sobald die Extended-Seite betroffen ist (V-25). Exit-Code 5: konform. |
| `LH-FA-13.a` | **konform** für den Replay: laufende Interaktion bis `Sync`, danach keine neue (1b). Die Frist ist übergeben. |
| `LH-FA-13.b` | **konform.** `ErrorResponse` mit Meldungscode, Klasse 5 gemerkt, nur die Verbindung endet (2d). |
| `SPEC-025` | **konform** (1c, Exit 3) |
| `SPEC-033` | **konform.** SQL im Klartext, keine Parameterwerte (U9) |
| `SPEC-041` | **konform.** Jede von `record` geschriebene Aufzeichnung lädt im Replay. |
| `ARC-002` | **konform.** Matcher, Cursor und Diagnose liegen in `internal/hexagon/services`; `a-check` 0 Befunde |

**Rückwärts (einfache Anfragen):** `TestE2EReplaySelect1`, `TestE2EReplayAbweichung`, `TestE2EReplayBeschaedigt` und `TestE2EOhnePostgresReplay` sind grün im Gate-Lauf.

---

## 5. Befunde des Verifiers

| ID | Kategorie | Befund | Beleg |
|---|---|---|---|
| V-25 | MEDIUM | **Die Diagnose der Art-Abweichung nennt das SQL der Extended-Seite nicht.** Betroffen sind drei Lagen: Eine einfache Anfrage kommt, wo ein `parse` erwartet ist (P4a, erwartetes SQL fehlt). Ein `parse` kommt, wo eine einfache Anfrage erwartet ist (P4b, empfangenes SQL fehlt). Ein `parse` kommt nach dem Ende der Aufzeichnung (P4c, empfangenes SQL fehlt). `LH-FA-10.a` verlangt als Mindestinhalt „erwartete Query“ und „tatsächlich empfangene Query“. Für den einfachen Pfad hält der Code das ein: `Query` nennt das empfangene SQL auch nach dem Ende. DoD-Punkt 2 sagt ohne Einschränkung „die SQL-Texte der Anweisungen“. Die Art-Abweichung hat dieser Slice selbst eingeführt (§1). Die Erkennung und Klasse 5 sind richtig; betroffen ist nur der Inhalt der Diagnose. Das SQL ist vorhanden: bei `parse` in der Nachricht, sonst über `objekte`. `TestReplayExtendedFalscheArt` prüft nur Stelle und Typ. Dritter Fund der Klasse „Diagnose nennt die Anweisung nicht vollständig“ nach F-339 und F-347. | Sonden P4a bis P4c; `internal/hexagon/services/replay.go:105`, `:131`, `:136`; `replay_extended_test.go:227`, `:234`; Plan §1, DoD-Punkt 2 |
| V-26 | LOW | **Der Bestätigungsfilter aus `f528c90` zählt Bestätigungen je Gruppe.** `LH-FA-18.a` legt Server-Nachrichten, die nach der nächsten Client-Nachricht eintreffen, in die folgende Gruppe. Dann hat die Flush-Gruppe keine Bestätigung, und ihr `parse` gilt als nicht angenommen, obwohl der Server es bestätigt hat. Die folgende Gruppe hat dafür eine Bestätigung mehr, als sie Nachrichten dieser Art führt. Die Diagnose nennt das Statement dann „unbekannt“ (P5). Damit sagt sie das Gegenteil dessen, was Plan §1 zusagt („nur vom Server bestätigte … wirken“). Ein Zähler je Interaktion statt je Gruppe hielte die Präfix-Eigenschaft ebenfalls ein, weil das Verwerfen nach einem Fehler bis zum `Sync` reicht, also bis zum Ende der Interaktion. Die Erkennung der Abweichung ist nicht berührt. Der Commit hat kein Review. | Sonde P5; `internal/hexagon/services/replay.go:231`–`:237`; Plan §1, §6 (Risiko späte Antworten) |
| V-27 | LOW | **Plan §3 führt nicht jede geänderte Datei.** Es fehlen die vier Review-Reports unter `docs/reviews/` und `internal/hexagon/services/replay_extended_test.go`. Dieselbe Klasse wie V-23 im Vorgänger-Slice; dort wurde sie bei der Closure behoben. `BEO-REPO/plan-folgt-korrektur-nicht`, `AGENTS.md` §3.9. | `git diff --stat d0e5d3b..f528c90`; Plan §3 |
| V-28 | INFO | **Lücken im Kopf.** Der Kopf nennt [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) nicht im `Bezug`, obwohl der Slice deren Wiedergabe-Hälfte liefert und §6 sie zitiert. In `Berührte Spec-Stellen` fehlt `LH-FA-11.a`, obwohl DoD-Punkt 2 Abnahmeszenario 6 beansprucht. `LH-FA-18.a` §Fehler deckt den Inhalt ab, deshalb nur INFO. | Plan Kopf, §6, DoD-Punkt 2 |
| V-29 | INFO | **`PGR-E3003` beim Start erzeugt auf Prozessebene der YAML-Leser** (`internal/adapters/driven/recording/yaml.go:560`), bevor der Service prüft. Die Prüfung in `NewReplayService` ist eine zweite Wache für jedes andere Repository. Ihr Bruch (U6) ist nur auf Unit-Ebene sichtbar; auf Prozessebene war die Zusage aus DoD 1c schon vor dem Slice wahr. Kein Handlungsbedarf, die DoD-Aussage stimmt. | Mutation U6; Prozess-Sonde aus 1c |
| V-30 | INFO | **Die zweite Phase des Runners ist unter einer Mutation nicht beobachtbar.** Jeder Mutant, der den Replay bricht, macht schon die erste Phase rot, und der Runner endet dort. Abnahmeszenario 7 lässt sich deshalb mit `make test-integration` allein nicht bewusst brechen. Der Gate bleibt dabei richtig rot, nur die Zuordnung „welcher Test fängt was“ fehlt. Belegt ist die Aussage hier über getrennte Phasen (1a). Für spätere DoD-Behauptungen an `TestE2EOhnePostgres*` braucht es denselben Weg, etwa als Option des Runners, nur die zweite Phase gegen ein anderes Image zu fahren. | `tools/test/run-integration-tests.sh`; Abschnitt 1a |

---

## 6. Gegenprobe dieser Datei

`make docs-check` nach dem Anlegen dieser Datei: `d-check: 161 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

## 7. Gesamturteil

**DoD-Punkt 1 und 3 sind bestätigt; Punkt 2 ist bis auf die Teilbehauptung 2c bestätigt.** Grundlage sind eigene Läufe:

- `make gates` mit Exit 0 an `f528c90`
- die Stufe `test` ohne Cache
- der Race-Detector 10-fach auf den Replay-Tests
- 20 Unit-Mutationen (U1 bis U14, A1 bis A6), alle rot im erwarteten Test
- 7 E2E-Läufe mit Mutanten, alle rot. Ausnahme ist U1 in Szenario 7, was erwartet ist; U1 ist im Abweichungs-Test rot.
- Abnahmeszenario 7 in getrennten Phasen: grün unverändert, rot ohne den Fix (A1)
- 5 eigene Sonden (Prozess-Start mit beschädigter Aufzeichnung, P4a bis P4c, P5)

Der Code hält [ADR-0007](../plan/adr/0007-strict-replay.md) im Matching, [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) und [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) ein. Er entspricht `LH-FA-18.a`, `LH-FA-13.a`/`.b`, `SPEC-025`, `SPEC-033`, `SPEC-041` und `ARC-002`. Bei `LH-FA-10.a` bleibt die Lücke V-25.

**Kein Befund ist HIGH.** Vor der Closure empfohlen:

1. V-25: Diagnose der Art-Abweichung und des Endes um das SQL der Extended-Seite ergänzen, mit Testfall; oder DoD-Punkt 2 und §1 begründet einschränken — Implementer/Planner.
2. V-26: Bestätigungen je Interaktion statt je Gruppe zählen, mit einem Fall für die späte Bestätigung; oder die Grenze in §1 und im Kommentar nennen — Implementer.
3. V-27, V-28: Plan §3 und den Kopf nachziehen — Implementer/Planner.
4. Über ein Review für `f528c90` entscheiden — Planner.
5. V-30 als Kandidat für einen Lerneintrag (neuer Sensor oder Runner-Option) in die Closure-Notiz nehmen — Planner.

---

## Nachverifikation zu 5a1da90

**Gegenstand:** Die Commits `1ab9db2` (V-25 bis V-28) und `5a1da90` (Test Q1 zu F-349) bei HEAD `5a1da90`. Geprüft wird nur dreierlei: DoD-Punkt 2 gegen `LH-FA-10`/`LH-FA-10.a` und `LH-FA-18.a` §Mismatch, mit Schwerpunkt auf der Diagnose bei Art-Abweichung und nach dem Ende der Aufzeichnung (P4a bis P4c); der Stand von V-26 bis V-28; ob DoD 1 und 3 schlechter geworden sind.

**Eingang:** dieser Bericht (V-25 bis V-30), das [vierte Folge-Review](2026-10-05-folge-review-4-slice-extended-query-replay.md) (F-349) und die Commit-Messages. Der Arbeitsbaum war bei meinem Start sauber.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Selbst gefahren.** Mutationen und Sonden liefen nur in Kopien von `git archive 5a1da90` im Scratchpad: je Kopie `go test -count=1` netzlos in einem Container aus der Stufe `source` des `Dockerfile` (Hilfs-Tag `pgr-verif2:source`). Kopien und Hilfs-Images (`pgr-verif2:source`, `pgr-verif2:test`) sind gelöscht.

### N1. DoD-Punkt 2 — **bestätigt**

**Bewusstes Brechen: der Zustand ohne den Fix.** Mutante NA setzt `replay.go` auf den Stand `f528c90` zurück, die Tests bleiben auf `5a1da90`. Ergebnis: `TestReplayExtendedDiagnoseArt` ist rot, und zwar aus dem richtigen Grund. Jede Meldung trägt Stelle, Typ und Klasse `PGR-E5001`, aber kein SQL der Extended-Seite. Das ist genau die Lücke aus V-25:

- P4a: `erwartet Client-Nachricht parse, empfangen Anfrage "SELECT 99"`
- P4b: `erwartet Anfrage "SELECT 7", empfangen Client-Nachricht parse`
- P4c: `keine aufgezeichnete Interaktion mehr nach 1; empfangen Client-Nachricht parse`

Ebenfalls rot ist `TestReplayExtendedDiagnoseSpaeteBestaetigung` (P5: `Anweisung erwartet unbekannt`). Unverändert ist das Paket grün.

**Einzelne Mutationen**, je eine Kopie:

| # | Mutation | Ergebnis |
|---|---|---|
| NB | `Query` ohne `nachrichtText` für die erwartete Nachricht (P4a) | rot nur in `…DiagnoseArt` P4a |
| NC | Ende der Aufzeichnung ohne `nachrichtText` (P4c) | rot nur in `…DiagnoseArt` P4c, P4c bind |
| ND | Art-Abweichung in `ClientMessage` ohne `nachrichtText` (P4b) | rot nur in `…DiagnoseArt` P4b, P4b bind, Q1 |
| NE | Bestätigungen je Gruppe statt je Interaktion (V-26 zurückgedreht) | rot nur in `…DiagnoseSpaeteBestaetigung` |
| NF | Grenze an der einfachen Anfrage am Cursor entfernt (F-349) | rot nur in `…DiagnoseArt` Q1: `bind (Anweisung unbekannt)` statt `"SELECT 1"` |
| NG | `nachrichtText` hängt immer die Parameterwerte an | rot in `…DiagnoseArt`, aber nur zufällig: in P4a reicht der geprüfte Teilstring über das angehängte `[]` hinaus |
| NG2 | `nachrichtText` nennt die Parameterwerte, sobald die Nachricht welche trägt | **grün** im ganzen Baum `./internal/...` (V-31) |

**Eigene Sonden** (zusätzliche Testdatei in der Kopie):

| # | Lage | Diagnose |
|---|---|---|
| P6 | `parse s1` in Interaktion 1; Interaktion 2 erwartet `bind s1`; empfangen wird eine einfache Anfrage | `erwartet Client-Nachricht bind (Anweisung "SELECT 1"), empfangen Anfrage "SELECT 99"`, richtig |
| P7 | mitten in einer Interaktion nach der Flush-Gruppe; erwartet `bind s1`; empfangen wird eine einfache Anfrage | `… Gruppe 2, Nachricht 1: erwartet Client-Nachricht bind (Anweisung "SELECT $1::text") …`, richtig |
| P8 | erwartet `execute` auf das unbenannte Portal in offener Transaktion; empfangen wird eine einfache Anfrage | `erwartet Client-Nachricht execute (Anweisung "SELECT 5")`, richtig |
| P9 | `bind` mit dem Wert `GEHEIM`, einmal bei der Art-Abweichung und einmal nach dem Ende | `… bind (Anweisung unbekannt)` bzw. `… bind (Anweisung "SELECT 1")`; der Wert erscheint nicht |

**Gegen die Spec.** `LH-FA-10.a` verlangt erwartete und tatsächlich empfangene Query. Beide stehen jetzt in allen drei Lagen der Art-Abweichung und nach dem Ende, soweit die Nachricht sich auf eine Anweisung bezieht. Wo die Aufzeichnung das Objekt nicht kennt, steht „unbekannt“. Die Session und die Interaktionsnummer stehen ebenfalls drin. `LH-FA-18.a` §Mismatch verlangt Gruppe, Nachricht und die Nachrichtentypen: Bei P4a, P6 bis P8 steht die Stelle mit Gruppe und Nachricht. Bei P4b fehlt sie zu Recht, denn erwartet war eine einfache Anfrage. Parameterwerte erscheinen nach Code-Lesung und P9 nicht (`SPEC-033`); dass kein Test diese Zusage auf den neuen Pfaden hält, ist V-31. Damit sind alle vier Teilbehauptungen von DoD-Punkt 2 bestätigt (2a, 2b, 2d unverändert aus dem Hauptteil, 2c jetzt vollständig). **V-25 ist behoben.**

### N2. Stand V-26 bis V-28

| ID | Stand | Beleg |
|---|---|---|
| V-25 | behoben | N1; NA bis ND rot |
| V-26 | behoben | Zählung je Interaktion (`replay.go`, `objekte`); NE rot in `TestReplayExtendedDiagnoseSpaeteBestaetigung`, und dieser Test ist die frühere Sonde P5 |
| V-27 | behoben, mit Rest (V-32) | §3 nennt `replay_extended_test.go` und die Berichte unter `docs/reviews/` |
| V-28 | behoben | Der Kopf nennt [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) im `Bezug` und `LH-FA-11.a` in den berührten Spec-Stellen |
| F-349 (Review) | behoben | NF rot in Q1 |

### N3. DoD 1 und 3 — **unverändert bestätigt**

- **`make gates` an `5a1da90`, eigener Lauf, Exit 0:**
  - `abdeckung-gegenprobe` grün, `make abdeckung-check` gesondert Exit 0
  - `a-check` 0 Befunde, `a-check-negativ` grün
  - `baseline-verify` v6.13.0 OK, 54 Dateien
  - `d-check` 162 Dateien, 0 Befunde
  - `commit-msg-gegenprobe` grün
  - `run-integration-tests: gruen`, in beiden Phasen; darunter alle `TestE2EReplayExtended*` und in der zweiten Phase `TestE2EOhnePostgresExtendedReplay`, also Abnahmeszenario 7
- Die Stufe `test` kam im Gate-Lauf aus dem Cache. In der Kopie lief sie deshalb mit `--no-cache-filter test`: `gofmt -l` leer, `go vet -tags integration` und alle sechs Unit-Pakete `ok`.
- **DoD 1:** Die beiden Commits ändern nur den Mismatch-Zweig des Use Case (`nachrichtText`, `objekte`), den Plan und die Abdeckungstabelle. Matching, Freigabe der Gruppen, Herunterfahren und `Validate` beim Start sind nicht berührt, und die Belege aus 1a bis 1c sind im Gate-Lauf grün. `objekte` läuft weiter unter `mu`. Nach dem Ende der Aufzeichnung kann es den Bereich nicht überschreiten, weil die Schleife über alle Interaktionen läuft (P4c, P9). Die getrennten Phasen aus 1a habe ich nicht wiederholt, weil der Replay-Pfad außerhalb der Diagnose unverändert ist.

### N4. Neue Befunde

| ID | Kategorie | Befund | Beleg |
|---|---|---|---|
| V-31 | LOW | **Auf den neuen Pfaden hält kein Test die Zusage „nie Parameterwerte“.** `nachrichtText` liest nur Typ, Namen und SQL; der Code ist richtig (P9). Fügt man aber Parameterwerte in die Diagnose ein, sobald die Nachricht welche trägt, bleibt der ganze Baum grün (NG2). Grund: `TestReplayExtendedDiagnoseArt` sendet `bind` stets ohne Werte und vergleicht nur Teilstrings. U9 im Hauptteil deckt den Pfad innerhalb einer Interaktion (`anweisungen`), nicht die Art-Abweichung und nicht das Ende. DoD-Punkt 2 und §1 sagen „ohne Parameterwerte“ beziehungsweise „nie Parameterwerte“ ohne Einschränkung (`SPEC-033`). Es ist dieselbe Klasse wie F-349: eine Zusage der Diagnose ohne Test, der sie bricht. Die DoD-Aussage selbst stimmt. | Mutation NG2; Sonde P9; `internal/hexagon/services/replay.go` (`nachrichtText`); `replay_extended_test.go` (`TestReplayExtendedDiagnoseArt`) |
| V-32 | INFO | **Plan §3 zählt die Berichte nicht mehr richtig.** Die Zeile `docs/reviews/` nennt „drei Folge-Reviews“; seit `e017ac2` liegen vier vor. Die Verzeichnis-Zeile deckt die Datei ab, nur die Zählung ist veraltet (`AGENTS.md` §3.9). | `git diff --stat d0e5d3b..5a1da90`; Plan §3 |

### N5. Gesamturteil der Nachverifikation

| DoD-Punkt | Urteil |
|---|---|
| 1 ([LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol)) | **bestätigt**, unverändert |
| 2 ([LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage)) | **bestätigt**. 2c ist jetzt vollständig; ohne den Fix ist `TestReplayExtendedDiagnoseArt` aus dem richtigen Grund rot (NA). |
| 3 (`make gates`) | **bestätigt**, Exit 0 an `5a1da90` |

**Kein Befund ist HIGH oder MEDIUM.** Vor der Closure empfohlen:

1. V-31: in `TestReplayExtendedDiagnoseArt` ein `bind` mit Wert senden und prüfen, dass der Wert fehlt, bei der Art-Abweichung und nach dem Ende — Implementer.
2. V-32: Zählung in §3 nachziehen — Implementer.
3. Die Klasse „Zusage der Diagnose ohne Test, der sie bricht“ ist jetzt der vierte Fund (F-339, F-347, F-349, V-31) und gehört als Kandidat in die Closure-Notiz — Planner.
