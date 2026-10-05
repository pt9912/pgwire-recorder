# Drittes Folge-Review: slice-extended-query-replay — 2026-10-05

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** Commit `2d3e0f7` (`git diff 3426462..2d3e0f7`). Er soll die Findings F-339 bis F-344 aus dem [zweiten Folge-Review](2026-10-05-folge-review-2-slice-extended-query-replay.md) abarbeiten. Frühere Teile des Slice prüft dieses Review nicht erneut. Schwerpunkte laut Auftrag:

- Status je Finding.
- Die Nachbildung `objekte` in `replay.go`, geprüft gegen das PostgreSQL-Protokoll: unbenannte Objekte, Status I/T/E am ReadyForQuery, einfache Anfrage, Close, Überschreiben. Dazu die Frage, ob das Verhalten bei Status `E` vertretbar ist.
- Ob `TestReplayHerunterfahrenSpaeteFrist` jetzt deterministisch ist, mit Wiederholungsprobe.
- Ob die 200-ms-Grenze die Mutation M1 trägt.

**Ablage:** Eigene Datei, Nummerierung fortlaufend ab F-347 (Reviewer-Skill §Output-Schema: „Folgeläufe als neue Datei statt Überschreibung“).

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-extended-query-replay.md` (§1, §6), `docs/plan/planning/open/slice-v1-abschluss-betrieb.md` (Kopf, §1), `docs/plan/planning/open/slice-replay-semantik-mismatch.md` (Kopf), `docs/plan/planning/welle-replay-semantik.md` (§4)
- `spec/lastenheft.md` [LH-FA-10](../../spec/lastenheft.md), [LH-FA-13](../../spec/lastenheft.md), [LH-FA-18](../../spec/lastenheft.md); `spec/spezifikation.md` LH-FA-10.a, LH-FA-13.a, `SPEC-034`, `SPEC-046`
- `internal/hexagon/model` (`Interaction`, `Validate`, Antworttypen), `internal/adapters/driving/pgwire/server.go` (`handle`, `replaySitzung`)
- `AGENTS.md` Hard Rules 3.3 bis 3.9; Register `BEO-REPO/plan-folgt-korrektur-nicht`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`
- PostgreSQL-Protokolldokumentation (Extended Query) zur Lebensdauer der Objekte. Benannte Statements leben bis zum Close oder Sitzungsende, und ein Parse auf ein bestehendes benanntes Statement scheitert. Das unbenannte Statement endet mit dem nächsten Parse auf das unbenannte Statement oder mit einer einfachen Anfrage. Portale leben bis zum Ende der Transaktion, das unbenannte zusätzlich nur bis zum nächsten Bind darauf oder bis zu einer einfachen Anfrage. Nach einem Fehler verwirft der Server alle Nachrichten bis zum Sync.
- Für Status `E`: das Verhalten des Servers nach dem PostgreSQL-Quelltext (`xact.c` `AbortCurrentTransaction`, `portalmem.c` `AtAbort_Portals`/`AtCleanup_Portals`), gelesen, nicht gegen einen Server erprobt
- Vorherige Reports am Modul: [Review](2026-10-05-review-slice-extended-query-replay.md) (F-321 bis F-329), [Folge-Review](2026-10-05-folge-review-slice-extended-query-replay.md) (F-330 bis F-338), [zweites Folge-Review](2026-10-05-folge-review-2-slice-extended-query-replay.md) (F-339 bis F-346)

**Ausgeführte Läufe im Repo:** `make docs-check` nach dem Anlegen dieser Datei. `make gates` lief nicht. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert (`git status --short` leer).

**Gegenproben und Mutationen:** Nur in Kopien (`git archive 2d3e0f7`) im Scratchpad. Die Unit-Läufe liefen in einem Container aus der Stufe `deps` des `Dockerfile` (Tag `pgr-rev3-deps`, ohne Netz). `make test-integration` lief in der unveränderten Kopie und baute dabei das Tag `pgwire-recorder:integration` aus demselben Commit neu. Image `pgr-rev3-deps` und Kopien sind gelöscht.

| # | Lauf oder Mutation | Ergebnis |
|---|---|---|
| K1 | Kopie ohne Änderung: `gofmt -l .`, `go vet -tags integration ./...`, Unit-Tests von `internal/hexagon/services` und `internal/adapters/driving/pgwire`, `tools/test/abdeckung.sh --check` | grün, `gofmt` leer |
| K2 | Kopie ohne Änderung: `make test-integration` | grün, darunter `TestE2EReplayExtendedPipeline` und `TestE2EReplayExtendedSigtermMittenInFolge` |
| K3 | `TestReplayHerunterfahrenSpaeteFrist` unverändert, `-count=300`, einmal mit Standard-`GOMAXPROCS`, einmal mit `GOMAXPROCS=1` | 600 von 600 grün |
| M1 | `<-geweckt` entfernt | rot, 40 von 40 (20 Standard, 20 mit `GOMAXPROCS=1`), `Nachricht 1: unexpected EOF` |
| M1b | wie M1, zusätzlich 250 ms Schlaf vor dem Zurücksetzen der Frist | grün, 5 von 5 (Grenze der 200 ms, siehe Antwort 3) |
| T1 | Zeitprobe aus dem zweiten Folge-Review: `replaySitzung` schläft 20 ms nach einer Nachricht ohne Antwort | grün, 100 von 100 (50 Standard, 50 mit `GOMAXPROCS=1`). F-341 behoben |
| D1 | Portal merkt den Namen des Statements und löst ihn erst beim Fehler auf | rot (`D1 Portal behält das Statement des bind`) |
| D2 | `objekte` nimmt die Nachricht am Cursor mit (`ni > c.nachricht`) | rot (P1, P2, P3 und weitere) |
| L1 | Ende der Transaktion nach einer Extended-Interaktion nicht nachgespielt | rot (`P3 Portal endet mit der Transaktion`) |
| L2 | einfache Anfrage entfernt das unbenannte Statement nicht | rot (`P2`) |
| L3 | jedes ReadyForQuery beendet die Portale, auch Status `T` | rot (`P3 Portal in offener Transaktion`) |
| L4 | close auf ein Statement entfernt es nicht | rot (`P1`) |
| S1 | Sonde: `parse s1 "SELEC 1"` scheitert (ErrorResponse, ReadyForQuery `I`), dann `bind s1`, wo Parse erwartet ist | Diagnose nennt `empfangen "SELEC 1"`; das Statement besteht auf dem Draht nicht (F-347) |
| S2 | Sonde: `s1` besteht mit `SELECT 1`, ein zweites `parse s1 "SELECT 2"` scheitert, dann `bind s1`, wo Parse erwartet ist | Diagnose nennt `empfangen "SELECT 2"`; auf dem Draht trägt `s1` weiter `SELECT 1` (F-347) |
| S3 | Sonde: `s1` besteht, die einfache Anfrage `DEALLOCATE s1`, dann `bind s1` | Diagnose nennt `empfangen "SELECT 1"` (F-348) |
| S4 | Sonde: Portal in Transaktion (`T`), dann Fehler (`E`), dann `execute` auf das Portal, wo Parse erwartet ist | Diagnose nennt das SQL des Portals (Antwort 2) |

---

## Status der Findings F-339 bis F-344

| ID | Stand | Beleg |
|---|---|---|
| F-339 | behoben, mit Rest | `objekte` (`replay.go:212`–`:264`) spielt den Verlauf nach. Close entfernt sein Ziel, eine einfache Anfrage entfernt das unbenannte Statement und das unbenannte Portal, und ein ReadyForQuery mit Status `I` beendet alle Portale. Die Sonden P1 bis P3 stehen als Fälle in `TestReplayExtendedDiagnoseLebensdauer`; L1, L2 und L4 sind rot. Plan §1 (`:34`) nennt die drei Fälle. Rest: Objekte, deren Anlage der Server abgelehnt hat, gelten als angelegt (F-347). |
| F-340 | behoben | Neue Fälle für close auf Statement und Portal, describe auf ein Portal, ein nach dem Bind überschriebenes unbenanntes Statement und die Grenze am Cursor. D1 und D2 sind jetzt rot. Die neue Abdeckungszeile (`docs/user/abdeckung-unit.md:35`) entspricht der Deklaration, und `abdeckung.sh --check` ist grün. |
| F-341 | behoben | `cancel()` folgt erst, wenn der Use Case das Parse hat und die Sitzung wieder im `Read` steht (`server_extended_test.go:762`–`:766`). Freigegeben wird nach dem Bind, und `Execute` geht erst nach `gesetzt` hinaus. K3 600 von 600, T1 jetzt grün, M1 weiter rot. Siehe Antwort 3. |
| F-342 | behoben | `zeilenText` schreibt NULL als `NULL` und jeden anderen Wert gequotet (`extended_replay_e2e_test.go:172`–`:191`). Die Sicht beim Aufzeichnen muss `ausführung 2: [NULL]` enthalten (`:202`). Damit ist belegt, dass der Test einen NULL-Wert vergleicht. K2 grün. |
| F-343 | behoben | Der Kopf von `slice-v1-abschluss-betrieb` nennt `SPEC-046` (`:16`). |
| F-344 | behoben | (1) Der Bezug von `slice-replay-semantik-mismatch` (`:14`) nennt jetzt wie die Welle-Zeile LH-FA-03 und LH-FA-13 und dazu die ADR. Der Titel in Zeile 1 stimmt mit dem der Welle-Zeile überein; die Welle-Zeile ergänzt nur in Klammern den Gegenstand. (2) `slice-v1-abschluss-betrieb` §1 `:38` nennt den vergebenen `PGR-E4006` und `SPEC-034` als Vergabestelle. |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-347 | LOW | `nachspielen` legt bei jedem aufgezeichneten Parse ein Statement an und bei jedem Bind ein Portal, unabhängig von der Antwort des Servers. Hat der Server die Nachricht abgelehnt, oder hat er sie nach einem früheren Fehler bis zum Sync verworfen, gibt es das Objekt auf dem Draht nicht. Ein zweites Parse auf ein bestehendes benanntes Statement scheitert, und das alte bleibt bestehen. Die Diagnose nennt dann als „empfangen“ das SQL eines Objekts, das es nicht gibt (S1), oder das falsche SQL (S2). Plan §1 sagt, das SQL werde genannt, „soweit es nach dem Verlauf der Aufzeichnung vor dem Cursor besteht“; der Kommentar an `anweisungen` sagt Entsprechendes. Die Fehlerantworten gehören zu diesem Verlauf. Die Abweichung wird weiter richtig erkannt, und Parameterwerte erscheinen nicht; betroffen ist nur der Inhalt der Diagnose. Der Fall setzt eine Fehlerantwort vor dem Cursor voraus. Er fehlt in den drei Fällen, die §1 aufzählt. | LH-FA-10.a („tatsächlich empfangene Query“); LH-FA-18.a | `internal/hexagon/services/replay.go:184`–`:188`, `:239`–`:248`; `docs/plan/planning/in-progress/slice-extended-query-replay.md:34` | ja (Sonden S1, S2) | Diagnose leitet den Bezug ohne Lebensdauer der Protokollobjekte her |
| F-348 | INFO | Für den Implementer, ohne erwartete Aktion. Die Nachbildung sieht nur Protokollnachrichten. SQL-Befehle, die Objekte beenden, sieht sie nicht. `DEALLOCATE` und `DISCARD ALL` entfernen Statements (S3), `ROLLBACK TO SAVEPOINT` entfernt die Portale aus dem Untertransaktionsteil, und `CLOSE` entfernt einen Cursor. Der Status am ReadyForQuery zeigt keinen dieser Fälle an. Ohne SQL zu deuten, kann der Replay das nicht erkennen. „Verlauf der Aufzeichnung“ in §1 deckt den Fall nur, wenn man darunter die Protokollebene versteht. | LH-FA-10.a | `internal/hexagon/services/replay.go:212`–`:216` | — | — (Hinweis) |

## Antwort auf die Prüffragen

1. **F-339 bis F-344.** Siehe Statustabelle. Alle sechs sind behoben, F-339 mit einem Rest bei abgelehnten Nachrichten (F-347).
2. **`objekte` gegen das PostgreSQL-Protokoll.**
   - **Unbenannte Objekte und Überschreiben:** Parse auf `""` ersetzt das unbenannte Statement, und Bind auf `""` ersetzt das unbenannte Portal. Ein Portal trägt das SQL seines Statements zum Zeitpunkt des Bind, und ein späteres Parse ändert es nicht (D1 rot). Close auf ein Statement lässt Portale stehen, die daraus gebunden wurden. Auch das entspricht dem Server, weil ein Portal den Plan selbst hält. Bei benannten Statements weicht das Überschreiben vom Server ab (F-347, S2).
   - **Einfache Anfrage:** Sie entfernt das unbenannte Statement und das unbenannte Portal. Danach bewertet der Status ihres ReadyForQuery das Transaktionsende. Das stimmt mit der Protokolldokumentation überein. `Validate` garantiert das abschließende ReadyForQuery (`validateQuery`), und eine Extended-Interaktion hat mindestens eine Gruppe. `in.Groups[len(in.Groups)-1]` kann daher nicht überlaufen.
   - **Status `I` und `T`:** Bei `I` enden alle Portale, benannte und unbenannte. Das ist das Ende der Transaktion, auch bei einer impliziten Transaktion, die durch einen Fehler endet. Bei `T` bleiben die Portale bestehen. L3 ist rot.
   - **Status `E`:** Dass die Portale bestehen bleiben, ist vertretbar. Ein Fehler in einem Transaktionsblock bricht die Transaktion ab (`AbortTransaction`, Blockzustand `TBLOCK_ABORT`). `CleanupTransaction` und damit `AtCleanup_Portals`/`PortalDrop` laufen aber erst mit dem ROLLBACK. Bis dahin sind die Portale unter ihrem Namen vorhanden und nur als fehlgeschlagen markiert. Ein Execute darauf scheitert an „current transaction is aborted“, nicht an „portal does not exist“. Die Diagnose bezieht sich also auf ein Objekt, das es dem Namen nach gibt (S4). „unbekannt“ wäre ebenso vertretbar. Ein Fehlbefund ist keine der beiden Lesarten. Beendet ein ROLLBACK den Block, kommt ReadyForQuery mit `I`, und die Portale enden. Belegt ist das aus dem Quelltext, nicht mit einem Server; der Implementer nennt es ungeprüft, und kein Test hält es.
   - **Grenzen:** Abgelehnte und verworfene Nachrichten (F-347). SQL-Befehle, die Objekte beenden (F-348).
3. **`TestReplayHerunterfahrenSpaeteFrist`.**
   - **Deterministisch für richtigen Code: ja.** `cancel()` folgt erst, wenn zwei Dinge gelten: Der Use Case hat das Parse (`len == 1`), und `liest` meldet ein laufendes `Read`. Weil `Read` beim Verlassen `liest` zurücksetzt, bevor der Use Case die Nachricht bekommt, bezeichnet die Und-Bedingung genau das nächste `Read`. Der Wächter hängt danach an `freigabe`. Das Bind beendet das `Read`, und die Sitzung wartet an `<-geweckt`. Erst nach `gesetzt` geht `Execute` hinaus, und zu diesem Zeitpunkt hat die Sitzung die Frist entweder zurückgesetzt oder wird es vor ihrem nächsten `Read` tun. Kein Schritt hängt an einer Schlafzeit; das `bis(200 ms, …)` läuft im richtigen Code immer ab und verzögert nur. K3 600 von 600, T1 (die Lage, die F-341 rot machte) 100 von 100 grün.
   - **Rest:** Eine fehlerhafte Variante, in der der Wächter die Frist nie setzt, ließe den Test an `<-conn.gesetzt` hängen, bis `go test` abbricht, statt nach Sekunden rot zu werden. Das betrifft nur die Laufzeit im Fehlerfall, nicht das Ergebnis.
4. **Die 200-ms-Grenze bei M1.** Sie trägt. Ohne `<-geweckt` setzt die Sitzung die Frist sofort zurück und steht Mikrosekunden nach dem Bind wieder im `Read`. Der Test sieht `liest` und gibt den Wächter frei. Dessen Frist bricht das `Read` ab, und die Sitzung meldet einen Fehler (M1 rot, 40 von 40, auch mit `GOMAXPROCS=1`). Erreicht die Sitzung das `Read` nicht binnen 200 ms, wird ebenfalls freigegeben. Solange sie die Frist noch nicht zurückgesetzt hat, scheitert ihr `Read` an der Frist des Wächters, und M1 bleibt rot. Grün bleibt M1 nur, wenn die Sitzung länger als 200 ms zwischen Bind und Zurücksetzen steht und erst nach dem Wächter zurücksetzt (M1b). Das setzt einen Stillstand von mehr als 200 ms voraus, den der Test nicht ausschließt; in K3 und M1 trat er nicht auf.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services` | geprüft. `objekte` läuft unter `mu`, nur im Mismatch-Zweig, ohne Zustandsänderung am Cursor; die Grenze am Cursor ist richtig (D2 rot). Parameterwerte werden nicht gelesen. Befund F-347, F-348. |
| Core-Reinheit | geprüft, ohne Befund. `replay.go` importiert keine neuen Pakete; `clear` ist eingebaut. `make a-check` lief nicht. |
| `internal/adapters/driving/pgwire` | geprüft, ohne Befund. Nur der Test ändert sich; der Wrapper hält nur Fristen an, die nicht null sind, wie zuvor, und `close(c.gesetzt)` läuft einmal, weil allein der Wächter von `replaySitzung` nach dem Verbindungsaufbau eine solche Frist setzt. Der Wächter in `handle` ist zu diesem Zeitpunkt über `fertig` beendet. |
| `test/integration` | geprüft, ohne Befund. K2 grün. |
| `docs/user/` | geprüft, ohne Befund. Die neue Abdeckungszeile entspricht der Deklaration des Tests, und ihre Aussagen sind durch Fälle gedeckt (D1, D2, L1 bis L4 rot). |
| `docs/plan/planning/` — Hard Rule 3.9 | geprüft. Plan §1 dieses Slice folgt dem Code; die Köpfe der Folge-Slices sind nachgezogen. Befund F-347 betrifft die Reichweite der Zusage in §1. |
| `spec/` — Strata, Hard Rule 3.4 | geprüft, ohne Befund. Keine Spec-Datei im Diff. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. `2d3e0f7` enthält keinen Move. |
| Hard Rules 3.5, 3.6, 3.8 | geprüft, ohne Befund. Keine ADR und keine Gate-Konfiguration im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft, ohne Befund. Die neuen Kommentare stehen im Indikativ. In `server_extended_test.go:770`–`:773` steht „ohne das Warten stünde sie schon wieder im Read“. Der Satz beschreibt keine verworfene Entwurfsalternative. Er nennt den Fehler, den der Test erkennt, und gehört damit zu dessen Zusage; die Form ist eine Stilfrage. |
| Commit-Message `2d3e0f7` | geprüft, ohne Befund. Sie nennt `slice-extended-query-replay`, [LH-FA-10](../../spec/lastenheft.md), [LH-FA-13](../../spec/lastenheft.md) und [LH-FA-18](../../spec/lastenheft.md), dazu keine Struktur-ID. |
| Register `BEO-REPO/*` | geprüft. Kein Finding dieses Laufs gehört zu `BEO-REPO/plan-folgt-korrektur-nicht` oder `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:**

- Diagnose leitet den Bezug ohne Lebensdauer der Protokollobjekte her (F-347; zweites Mal nach F-339, jetzt auf der Seite der Anlage)

## Verdikt

**Merge-blockierend nach dem Skill:** nein; kein HIGH. F-339 bis F-344 sind abgearbeitet. `TestReplayHerunterfahrenSpaeteFrist` ist für richtigen Code deterministisch, und M1 bleibt rot. Der verbleibende Befund betrifft den Inhalt der Diagnose nach einer Fehlerantwort in der Aufzeichnung.

**Übergabe:**

- F-347 an den Implementer. Widerspricht er der Einstufung, genügt eine Begründung im Plan oder in der Closure-Notiz.
- F-348 ist ein Hinweis ohne erwartete Aktion.
- Die Finding-Klasse geht in die Closure-Notiz §7. Die DoD-Konformität prüft der Verifier separat; F-345 und F-346 aus dem zweiten Folge-Review bleiben dort adressiert.
