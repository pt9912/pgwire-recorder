# Verifikations-Report: slice-v1-abschluss-verbindungsende-record — 2026-10-16

**Rolle:** Verifier (Modul 11), frischer Kontext. Frage: Bauen wir es richtig, gegen Plan, DoD und Entscheidungen?

**Gegenstand:** `slice-v1-abschluss-verbindungsende-record` (Schritt 40, Ende des Servers im Record: Use Case und Zustellung), Diff `9c19eba6..96b2f8ab`; Lieferung `a4daf4da`, Tests `a16c27f0`, Nacharbeit zum Review `96b2f8ab`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

**Eingang:** Plan (§1, §2, §3, §4, §6, §7 *Belege des Implementers* und *Nacharbeit zum Review*, §8), Review-Report vom selben Tag (0 HIGH, 2 MEDIUM, 1 LOW, 4 INFO; nicht neu gefahren), Spezifikation `LH-FA-25.a` (*Meldung des Endes beim Aufzeichnen*, *Stelle von `connection_closed` beim Aufzeichnen*, *Reihenfolge der Gründe beim Aufzeichnen*, *Herunterfahren beim Aufzeichnen*), `LH-FA-18.a`, `SPEC-002`, `SPEC-043`, Architect-Commits `23374c20`, `d5c42925`, `ccbb1073`, `AGENTS.md` §3.9 bis §3.14.

## Ergebnis

**Die Liefer- und Gate-Punkte der DoD sind erfüllt bis auf einen Befund, der vor der Closure behoben werden sollte:** Die Zusage „nach dem Ende nimmt der Use Case nichts mehr an“ hat für den **einfachen Weg** keine Mutation (V-1, MEDIUM): Der Mutant, der `serverBeendet()` in `Query` auslässt, bleibt grün, und §7 sagt „Keine weiteren grünen Mutanten“. Alle übrigen geprüften Zusagen fangen eine von mir selbst gesehene, rote Mutation; die Nacharbeit zu F-1 und F-3 trägt; die Messung stimmt (860). `make abdeckung-check` und `make gates-schnell` selbst gefahren: Exit 0; der Stop-Hook antwortet `approve`. Dazu V-2 (LOW, offene Sendung F-2 des Reviews, Planner) und zwei INFO. Offen sind die Closure-Pflichten (Notiz, Register, Ausgänge, Haken).

## 1. DoD-Punkte

| DoD-Punkt (§2) | Beleg | Urteil |
|---|---|---|
| Punkt 1: einfache Interaktion mit Ende im Use Case (`FATAL`/`PANIC` + `connection_closed`), Antworten zurück, Interaktion abgeschlossen, kein `PGR-E4003`/`PGR-E6001`, nichts danach, `offset_ms`, `--record-empty-sessions`, `wirdVerworfen` bleibt die eine Stelle | `record.go` (`bisEnde`, `endetMitEnde`, `serverBeendet`, Zweige in `Query`); Tests `TestRecordEndeSchweregrad`, `Stelle`, `NichtsDahinter`, `Enden`, `ZeitUndLeereSessions`, `EinzigeInteraktion`; Mutanten unten | erfüllt, **mit V-1**: „nichts danach“ ist für `Query` als Weg des Endes ohne Mutation |
| Punkt 2: Extended-Interaktion, Ende in der letzten Gruppe mit `Sync` oder `Flush`, Pipeline gehört nicht zur Aufzeichnung, offene Gruppe `PGR-E6001` samt Session verworfen, `ERROR` ohne Ende bleibt `PGR-E6001`, Kommentare „Abschluss“ | `TestRecordEndeGruppe`, `DanachNichts` (`Pipeline`), `Enden` (mit Pipeline); die bleibenden Tests `TestRecordExtendedFehlerantwortVorDemAbbruch`, `…SendefehlerNachFehlerantwort` unverändert grün; `grep -n "Sync"` über beide Dateien: übrig sind nur Kommentare, die das `Sync` selbst meinen (`record.go` 73, 81, 207, 279, 282; `ports/driving/record.go` 49), Zeile 207 beschreibt `Query` bei gesendetem `Sync` und ist an das Feld `sync` gebunden | erfüllt |
| Punkt 3: Adapter schreibt bis einschließlich der Fehlerantwort, schließt beide Verbindungen, ohne `ReadyForQuery`/Meldung/Exit-Code; Schreibfehler davor `EndWriteFailed`; `send` lässt `connection_closed` aus mit Fall und Mutation im Record-Pfad | `server.go` (`schreibe`), `TestRecordEndeSendetDavorUndSchliesst` (vier Läufe), `TestRecordEndeSchreibfehler`; Mutanten A1 bis A3, A7 | erfüllt |
| Punkt 4: je Zusage eine Mutation, Tabelle in §7 | Tabelle steht; 12 meiner Mutanten bestätigen die Zeilen (Abschnitt 3); eine Zusage ohne Mutation (V-1) | **teilweise** (V-1) |
| Punkt 5: `make gates` grün | Nachweis `.harness/state/gates-passed.diffsha`/`.dateiliste` von 21:45, nach `96b2f8ab` (21:31); Stop-Hook antwortet `approve` (Liste gleich dem Arbeitsbaum oder Unterschied nur auf der Positivliste); `git diff 6c3b2484..96b2f8ab -- internal` berührt nur `record_ende_test.go` (+33); Arbeitsbaum sauber; `make gates-schnell` selbst gefahren: Exit 0 (`make test`, `lint`, `a-check`, `abdeckung-check`, `kopf-check` darin); den vollen Lauf habe ich nicht wiederholt | erfüllt |
| Punkt 6: Review-Report liegt vor | `2026-10-16-review-slice-v1-abschluss-verbindungsende-record.md` (0 HIGH, 2 MEDIUM, 1 LOW, 4 INFO); F-1 und F-3 behoben (geprüft), F-2 offen (V-2) | Report liegt vor; Haken gehört der Closure |
| Punkte 7 bis 10: Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | §7 trägt Belege und Nacharbeit, keine Closure-Notiz, keine Ausgänge; Haken leer | offen, Closure-Arbeit (Abschnitt 6) |

## 2. Plan gegen Code, Entscheidungs-Konformität

- **Plan gegen Diff (§3):** Der Diff berührt `services/record.go`, `ports/driving/record.go` (nur Kommentare), `model/recording.go` (nur `EndServer`, ohne Logik), `pgwire/server.go`, die Tests (`record_ende_test.go`, `server_record_ende_test.go`, Fakes in `record_test.go`, `server_test.go`), `abdeckung-unit.md`, Roadmap, Plan, Review, `spec/spezifikation.md` — genau die Zeilen der Tabelle in §3. `ports/driven/upstream.go` bleibt unberührt wie zugesagt. Kein `nolint`, keine neuen Importe im Produktcode. Zwei Schichten (Kern, PGWire-Adapter), drei Liefer-Punkte; die Zählung von `docs/user` bleibt Planner-Entscheidung (Review F-7).
- **`LH-FA-25.a` gegen Code:** *Meldung*: `bisEnde` kürzt hinter der ersten `connection_closed`, `Query` und `AwaitServer` liefern bis einschließlich; der Adapter sendet davor, meldet `Delivered`, beendet mit `EndServer`, liest nichts mehr (Wache `r.beendet`), der Upstream wird im Use Case geschlossen. Schreibfehler: `EndWriteFailed`, `interactions[:len-unzugestellt]`, `PGR-E4003`. *Stelle*: Prüfung der Form des Modells (`pruefeEnde`, `ohneReadyForQuery`) plus der Zweig „`connection_closed` ohne laufende Interaktion“ in `zuordnen` (`PGR-E6001`), jede andere Nachricht `PGR-E1000`; die Texte von `d5c42925` stehen so im Code. *Form*: `serverEnde` und `l.offen = nil` schließen alles Weitere aus. *Reihenfolge der Gründe*: Fehler des Upstream-Ports vor der Form (`TestRecordEndePortfehlerVorForm`, M16); „dann die Form der Session“ wird nicht gerufen (siehe nächster Punkt). *Herunterfahren*: `Shutdown` mitten in der Interaktion gibt nicht frei, `Delivered` nach dem Ende schon; Zwangsende nach gelesenem Ende: keine Meldung, Session geschrieben (auch vor `Delivered`, F-1); Zwangsende vor dem Lesen: `PGR-E4006`.
- **`ValidateEnd` wird nicht gerufen:** `grep -n ValidateEnd internal/hexagon/services/*.go` ohne Treffer im Produktcode; die Entscheidung des Architect (§6: Reihenfolge erfüllt die Form von selbst) trägt über `serverEnde`, aber nur auf dem Weg über `AwaitServer` (V-1); die Tests rufen `sess.ValidateEnd()` auf der geschriebenen Session (`endeAngenommen`).
- **`ErrSessionEnded` vor `ErrShutdown`:** im Test festgelegt (`Anfrage, herunterfahren true`), im Port-Kommentar genannt (Review F-4, Architect, keine Aktion erwartet).
- **`SPEC-002`/`SPEC-043`, `V` und `S`:** `beendetVerbindung` liest `V`, bei leerem oder fehlendem `V` `S`, genau verglichen; Mutant „nur `S`/`V` leer“ ist rot.
- **Plan folgt Code (`AGENTS.md` §3.9):** §1, §3, §6, §7, Kopf stimmen mit dem Diff. Beobachtung (3) in §7 nennt `d5c42925`, die Spezifikation trägt die Entscheidung dort.
- **Kommentare (`AGENTS.md` §3.7, §3.11):** Die Kommentare an `serverEnde`, `serverBeendet`, `Query`, `AwaitServer`, `ClientMessage`, `schreibe`, `EndServer` und im Port sind Indikativ; zu jeder Zusage gibt es einen Test, mit einer Ausnahme in der Stärke: „nach dem Ende nimmt es nichts mehr an“ (`Query`, `serverEnde`) ist für den Weg des Endes über `Query` ungeprüft (V-1).

## 3. Bewusstes Brechen (Modul 11), eigene Mutanten

Weg: Kopie von `go.mod`, `go.sum`, `internal` in ein frisches Verzeichnis je Mutant (Scratchpad), genau eine Ersetzung mit erzwungenem Treffer, `docker run --rm --network=none` im Image der Stufe `deps`, `go test -count=1 -run TestRecordEnde` über `services` beziehungsweise `pgwire`; ein Lauf zur Zeit, unveränderte Kopie grün, Arbeitsbaum unberührt.

| # | Zusage / Variante | Mutant | Ergebnis |
|---|---|---|---|
| F1a | Zwangsende nach gelesenem Ende vor der Zustellung: keine Meldung | `CloseSession`: Zwangsende zählt als laufend, wenn `serverEnde && unzugestellt > 0` | rot: `ZwangsendeVorZustellung` (`einfach`, `erweitert`) |
| F1b | dasselbe in `laeuft()` | `laeuft()` mit `\|\| serverEnde && unzugestellt > 0` | rot: dieselben zwei Läufe |
| M17 | Anfrage nach dem Ende nicht an den Upstream (auch vor `Delivered`) | Prüfung am Anfang von `Query` entfernt | rot: `Anfrage, herunterfahren true`, `Anfrage vor der Zustellung` (nach 2 s) |
| M15 | Hinter `connection_closed` übernimmt der Use Case nichts | `bisEnde` kürzt nicht | rot: `NichtsDahinter` (`einfach`, `erweitert`) |
| – | dasselbe, nur `AwaitServer` | `rs = bisEnde(rs)` entfernt | rot: `NichtsDahinter/erweitert` allein |
| – | dasselbe, nur `Query` | `responses = bisEnde(responses)` entfernt | rot: `NichtsDahinter/einfach` allein |
| V | `V` geht vor `S` | `beendetVerbindung` liest nur `S` | rot: `Schweregrad` (`V FATAL, S ERROR`, `V ERROR, S FATAL`) |
| A1 | `send` lässt `connection_closed` aus | Auslassung entfernt | rot: `SendetDavorUndSchliesst` (alle vier), `Schreibfehler` |
| A2 | Adapter schließt nach dem Ende | `beende` bei Ende entfernt | rot: `SendetDavorUndSchliesst` (alle vier, je 5 s) |
| A3 | Ereignis `EndServer` | `EndWriteFailed` statt `EndServer` | rot: alle vier Läufe |
| A3b | `Delivered` vor dem Schließen | `Delivered` nach `beende` | rot: alle vier Läufe |
| A7a | `schreibe` liefert nach dem Ende wahr (Wache bleibt) | `return true` statt `false` | **grün**, äquivalent, siehe Abschnitt 4 |
| **V1m** | Nach dem Ende nimmt der Use Case nichts mehr an, **einfacher Weg** | `Query`: `if endetMitEnde(responses) { l.serverBeendet() }` entfernt | **grün, kein Äquivalent** (V-1) |

Zu V1m: Ich habe in einer Kopie den Test `TestVerifyEinfachDanach` ergänzt (`beginne` mit den Antworten `FATAL` + `connection_closed`, `Query`, `Delivered`, dann `Query` und `ClientMessage`): im Original grün, mit dem Mutanten rot („Query nach dem Ende: <nil>“). Der Fall gehört also zu einer Zusage der Spezifikation (*Form*: „nimmt danach keine Anfrage und keine Client-Nachricht mehr an“), kein Äquivalent. Im Betrieb fängt der Adapter ihn über `r.beendet` (A7-Invariante), aber der Use Case sagt es selbst zu (Kommentar an `serverEnde`, `Query` und Port).

## 4. Register `gruener-mutant-faelschlich-aequivalent` (3×) — die in §7 als äquivalent geführten Mutanten

| Mutant | Invariante und Ort in §7 | Meine Prüfung | Urteil |
|---|---|---|---|
| A7 (`schreibe` liefert wahr nach dem Ende, Wache bleibt) | `beende` setzt `r.beendet` vor `CloseSession`, `clientRichtung` prüft `r.beendet` vor jedem Lesen; Ort `clientRichtung` in `server.go` | `beende` (`einmal.Do`) setzt `beendet.Store(true)` als Erstes (`server.go` 594); `clientRichtung` prüft vor jedem `Receive` (711), `serverRichtung` ruft `AwaitServer`, das nach dem Ende `ErrSessionEnded` liefert und ohne Meldung endet; A7a selbst gefahren: grün; A7b (Wache ebenfalls entfernt) laut §7 rot | äquivalent, Invariante und Ort tragen |
| Doppelabsicherung M9 (`pruefeEnde` und `ohneReadyForQuery` lehnen `[Fehlerantwort, ready_for_query, connection_closed]` beide ab) | jede allein bleibt grün; beide aufgehoben ist rot | Quelltext gelesen (`response.go` 136 bis 171): die erste Prüfung trifft „connection_closed nur als letzte Nachricht“ und „unmittelbar nach `FATAL`“, die zweite jedes `ready_for_query` neben dem Ende; für diese Folge lehnen beide ab; M8 (nur `ohneReadyForQuery` entfernt) fängt „`ready_for_query` davor“ | äquivalent im Sinne der doppelten Absicherung, mit Ort und Invariante belegt |
| M17 ohne `Shutdown` | in der Nacharbeit zurückgenommen | M17 rot (siehe Abschnitt 3) | korrigiert, trägt |
| **nicht geführt:** V1m | – | grün, kein Äquivalent | **fehlt in §7** (V-1): „Keine weiteren grünen Mutanten“ ist falsch |

Einordnung: V1m ist nicht „fälschlich äquivalent“ (er wird nicht so geführt), sondern fehlt in der Liste: näher an `BEO-REPO/liste-gruener-mutanten-unvollstaendig` beziehungsweise, da der Slice einen neuen Vertrag liefert, an `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (Zusage ohne Test); die Zuordnung entscheidet der Planner bei der Closure.

## 5. Messung

`git diff --numstat 9c19eba6..96b2f8ab -- internal test docs/user`: **835 Einfügungen, 25 Löschungen, Summe 860**; Produktcode (`internal` ohne `*_test.go`) 130; unter der Meldeschwelle 1000 und der Marke 1400; stimmt mit §7.

## 6. Findings

| ID | Kategorie | Befund | Zuständig |
|---|---|---|---|
| V-1 | MEDIUM | Die Zusage „nach dem Ende nimmt der Use Case nichts mehr an“ (`LH-FA-25.a` *Form*, Kommentar an `serverEnde`/`Query`, Port) hat für den Weg über `Query` keinen Test: `TestRecordEndeDanachNichts` erreicht das Ende nur über `AwaitServer` (`nachEnde`), `Enden`/„einfache Anfrage“ ruft danach keine weitere Anfrage. Der Mutant, der `serverBeendet()` in `Query` auslässt, bleibt grün (`go test -run TestRecordEnde`, `services`), mit einem Fall „einfache Anfrage mit Ende, dann `Query` und `ClientMessage`“ ist er rot. §7 führt den Satz „Keine weiteren grünen Mutanten“. Behebung: ein Fall (rund 15 Zeilen) in `TestRecordEndeDanachNichts` mit `beginne(t, ende)`, `einfach()`, `Delivered`, dann `Query` und `ClientMessage`, Zeile und Mutante in die Tabelle von §7, Messung neu. | Implementer |
| V-2 | LOW | Review F-2 ist offen: Die Grenze (1) in §7 nennt den Test des Schreibfehlers in `serverRichtung` (Extended) „ohne Adresse (Ausgang bei F-2 unten)“; `slice-v1-abschluss-verbindungsende-upstream` führt ihn nicht (`grep -n "Schreibfehler\|serverRichtung"` im Plan: kein Treffer) und ist nicht annehmend. Der Implementer hat die falsche Adresse zurückgenommen und die Entscheidung an den Planner gegeben. Closure-Pflicht: Sendung samt Nachzählung in einem Nehmer (`AGENTS.md` §3.13) oder die Idee streichen mit Begründung; ohne eines von beiden bleibt eine ungeprüfte Grenze ohne Adresse. | Planner |
| V-3 | INFO | Die Spezifikation `LH-FA-25.a` *Stelle von `connection_closed` beim Aufzeichnen* (`d5c42925`) hat einen Absatz, in dem der Satz „Die Zeile ‚`ready_for_query` davor‘ entscheidet nur die einfache Anfrage.“ und der Anfang des folgenden Satzes („Das Ende liegt in der Gruppe …“) auf einer Zeile stehen (Umbruch fehlt, Lesbarkeit; inhaltlich kein Befund). | Architect (ohne erwartete Aktion) |
| V-4 | INFO | „Dann die Form der Session (höchstens ein Ende)“ in *Reihenfolge der Gründe beim Aufzeichnen* hat im Record keinen Test und keinen Pfad: Der Use Case ruft `ValidateEnd` nicht (Entscheidung des Architect), die Zusage wird von `serverEnde` getragen. Mit V-1 ist der Weg über `Query` gedeckt. Kein eigener Fall nötig; für die Closure als „von der Reihenfolge erfüllt“ vermerken. | Planner |

## 7. Vorschlag für die Risiko-Ausgänge (zur Closure; nicht eingetragen)

- *Diff an der Grenze einer Review-Sitzung:* **entfallen** — gemessen 860 (Marke 1400, Meldeschwelle 1000), das Review lief in einer Sitzung, der vorgemerkte Schnitt wurde nicht gebraucht. Die drei Einträge des Registers aus §8 (`schaetzung-in-zeilen-unterschaetzt-tests`, `groesse-erst-in-der-pruefung-des-architect-sichtbar`, `schnitt-laesst-haelfte-an-der-grenze`) sind damit nicht zum dritten Mal belegt; die Schätzung lag über der Messung (Mitte 1050 gegen 860).
- *`slice-v1-abschluss-sessions` und `slice-v1-abschluss-anmeldung` ändern `record.go`/Ports vor diesem Slice:* **entfallen** — beide liegen in `done/`, der Slice lief auf deren Stand (Basis `9c19eba6`), die Tests von `CloseSession`, `TestRecordZwangsendeLeereSession` und Nachbarn sind grün.
- *Fakes bilden nur ab, was die Spezifikation sagt; bis `-upstream` liefert kein echter Upstream ein Ende:* **eingetreten** → `slice-v1-abschluss-verbindungsende-upstream` (E2E, Abnahmeszenario 18); das ist ein Folge-Slice mit Kennung.
- Hinzu kommen die Ausgänge der Review-Findings aus §7 (F-2 mit V-2; F-7 Schichtzählung `docs/user` beim Planner) und die Finding-Klassen für den Zähler (Zusage ohne Mutation in zwei Läufen: Zustand zwischen Lesen und Zustellen im Review, Weg über `Query` hier).

## 8. Eigene Läufe

- `make abdeckung-check`: Exit 0.
- `make gates-schnell`: Exit 0 (Arbeitsbaum danach sauber).
- Stop-Hook `.claude/hooks/stop-require-gates.sh`: `approve`.
- Mutanten: 13 (Abschnitt 3), eine Kopie mit Zusatztest für V1m; ein Lauf zur Zeit, kein `pkill`, kein paralleler Docker-Lauf.
- Nicht gefahren: voller `make gates` (Nachweis von 21:45 nach `96b2f8ab`, Code-Stand gleich), der Review-Lauf, die Mutanten, die ich nicht in Abschnitt 3 nenne.
