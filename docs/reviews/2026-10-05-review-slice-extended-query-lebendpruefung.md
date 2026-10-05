# Review-Report: slice-extended-query-lebendpruefung — 2026-10-05

**Review-Art:** Code und Spezifikation, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff 90f6fb1..0841891`, also die Commits `f906c04` ([ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) als Proposed, Index, Plan), `4abad8e` (ADR angenommen), `eb61782` (Lastenheft), `7ffd029` (Spezifikation) und `0841891` (Code, Tests, Handbuch, Abdeckung, Plan) für `slice-extended-query-lebendpruefung`. Schwerpunkte laut Auftrag: Erkennung gegen den Scanner von PostgreSQL, Transaktionsstatus der Antwort außerhalb der Reihe, Überspringen aufgezeichneter Lebendprüfungen, Spec-Änderungen über den Auftrag hinaus, Stabilität der Integrationstests, Register-Lehren, Commit-Hook.

**Skill:** `.harness/skills/reviewer.md` @ `0841891`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-extended-query-lebendpruefung.md` (Kopf, §1 bis §6, §8), `docs/plan/planning/open/slice-v1-abschluss-sessions.md` (Kopf, §1, §2, §6)
- `spec/lastenheft.md` [LH-FA-03](../../spec/lastenheft.md#lh-fa-03--replay-modus), [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [LH-FA-18](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [LH-QA-02](../../spec/lastenheft.md#lh-qa-02--geringe-eingriffe-in-die-anwendung)
- `spec/spezifikation.md`: LH-FA-03.a, LH-FA-03.b, LH-FA-09.a, LH-FA-10.a, LH-FA-12.a, LH-FA-18.a (Interaktion, Gruppen, Replay), `SPEC-011`, Meldungscodes
- [ADR-0007](../plan/adr/0007-strict-replay.md), [ADR-0008](../plan/adr/0008-kein-sql-parser-in-v1.md), [ADR-0012](../plan/adr/0012-extended-query-gruppen.md), [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md), ADR-Index
- `AGENTS.md` Hard Rules 3.3 bis 3.9; Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/spec-randform-erst-im-review-entschieden`, `BEO-REPO/plan-folgt-korrektur-nicht`, `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients`
- `docs/user/benutzerhandbuch.md` (Wiedergabe, Mit einem Datenbanktreiber arbeiten, Meldungscodes, Fragen), `docs/user/abdeckung-*.md`
- Vorheriger Report am Modul: [Review zu `slice-extended-query-replay`](2026-10-05-review-slice-extended-query-replay.md) (bis F-329); höchste vergebene Nummer im Repo vor diesem Lauf F-349

**Ausgeführte Läufe im Repo:** `make docs-check` und `make abdeckung-check` nach dem Anlegen dieser Datei. `make gates` lief nicht. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert (`git status --short` leer).

**Mutationen und Sonden:** Nur in Kopien (`git archive 0841891`) im Scratchpad. Unit-Tests des Pakets `internal/hexagon/services` in einem Container aus der Stufe `deps` des `Dockerfile` (ohne Netz). Integrationstests mit `-test.run Lebendpruefung` aus einem Image der Stufe `integration` gegen `postgres:17-alpine` (Digest aus `harness/mk/`) in einem internen Docker-Netz. Sonde S1 mit `psql -c` in `postgres:16-alpine` und `postgres:17-alpine` ohne Netz. Images, Container, Netz, Cache-Volume und Kopien sind gelöscht.

| # | Lauf oder Mutation | Ergebnis |
|---|---|---|
| K0 | Kopie ohne Änderung: `go vet`, Unit-Tests `internal/hexagon/services` | grün |
| M1 | Antwort außerhalb der Reihe immer mit Status `I` | rot (`TestReplayLebendpruefungAusserDerReihe`, `…Extended`) |
| M2 | `ClientMessage` merkt den Status nicht | rot (`TestReplayLebendpruefungExtended`) |
| M3 | Status nach dem Handshake leer | rot (drei Tests) |
| M4 | Lebendprüfung auch mitten in einer Extended-Interaktion beantwortet | rot (`TestReplayLebendpruefungExtended`) |
| M5 | `mitten` nur über `nachricht` (nach einer Flush-Gruppe nicht mitten) | rot (`TestReplayLebendpruefungExtended`) |
| M6 | aufgezeichnete Lebendprüfungen nicht herausgenommen | rot (`TestReplayLebendpruefungAufgezeichnet`) |
| M7 | Filter auch auf Extended-Interaktionen (leerer `Parse`) | rot (fünf Extended-Tests) |
| M8 | Meldung nach der letzten Interaktion in `Query` nennt `c.pos` statt der aufgezeichneten Nummer | rot (`TestReplayLebendpruefungAufgezeichnet`) |
| M9 | dasselbe in `ClientMessage` | **grün** (F-356) |
| M10 | Lebendprüfung ordnet eine Session zu | rot (`TestReplayLebendpruefungAusserDerReihe`) |
| M11 | `merke` nimmt das erste statt das letzte `ReadyForQuery` | grün; äquivalent, weil jede Antwortliste nach `Validate` höchstens ein `ReadyForQuery` trägt |
| M12 | `Query` merkt den Status nicht | rot (`TestReplayLebendpruefungAusserDerReihe`) |
| L1 | `\v` nicht als Leerraum | rot (`TestLebendpruefungErkennung`) |
| L2 | Zeilenkommentar endet nur an `\n` | rot |
| L3 | Blockkommentare nicht verschachtelt | rot |
| L4 | offener Blockkommentar zählt | rot (zwei Tests) |
| L5 | `--` im Blockkommentar beginnt einen Zeilenkommentar | rot (drei Tests) |
| L6 | U+00A0 als Leerraum | rot |
| L7 | `;` als Leerraum | rot (drei Tests) |
| L8 | Zeilenkommentar überspringt ein Zeichen zu viel | rot |
| I0 | Integrationstests `Lebendpruefung` ohne Änderung, dreimal | 3× grün; je Lauf 20 bis 21 s, davon je 9,7 s für `Pgxpool` und `DatabaseSQL`, 0,35 s für `WiePostgres` |
| I1 | Integration: Antwort außerhalb der Reihe abgeschaltet | rot (alle drei Tests, Exit-Code 5) |
| I2 | Integration: aufgezeichnete Lebendprüfungen nicht herausgenommen | rot (`Pgxpool`, `DatabaseSQL` in den Lagen mit Pausen beim Aufzeichnen) |
| I3 | Integration: Status immer `I` | rot (`WiePostgres`) |
| S1 | `psql -c` mit `"\v-- ping\v"`, `"/* a /* b */ c */"`, `"-- a\rSELECT 1"` gegen PostgreSQL 16 und 17 | 17: leere Antwort, leere Antwort, Ergebnis `1`. **16: `ERROR: syntax error at or near ""`** für den Text mit `\v`; die anderen beiden wie 17 (F-350) |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-350 | HIGH | Der vertikale Tabulator ist nur in PostgreSQL 17 Leerraum. PostgreSQL 16 beantwortet eine Anfrage aus `\v` und einem Kommentar mit einem Syntaxfehler (Sonde S1). Die Spezifikation zählt `\v` ohne Serverversion zum Leerraum. Der Code folgt ihr und filtert beim Laden ohne Blick auf die aufgezeichnete Antwort. Eine gegen PostgreSQL 16 aufgezeichnete Interaktion mit `\v` und `ErrorResponse` wird deshalb übersprungen. Kommt dieselbe Anfrage in der Wiedergabe, erhält der Client `EmptyQueryResponse` statt des aufgezeichneten Fehlers. In einer Transaktion bleibt der Status dabei `T`, wo der Server `E` meldete. Das ist eine Antwort, die nicht aus der Aufzeichnung stammt und von der des aufgezeichneten Servers abweicht. Das Lastenheft sagt „beantwortet … wie PostgreSQL“, die ADR schließt aus, „was PostgreSQL mit einem Fehler beantwortete“. Das Lastenheft schränkt die Serverversion nicht ein. Der Kommentar an `leerraum` sagt ohne Version zu, es seien „die Zeichen, die der Scanner von PostgreSQL als Leerraum liest“. Das Risiko steht in Plan §6 als offen und ist mit S1 eingetreten. Der Auslöser ist selten, kein beobachteter Treiber sendet `\v`. Der Skill führt eine unpassende Antwort im Replay-Pfad aber ohne Häufigkeitsschwelle als HIGH. | [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) Negative; [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) §Entscheidung „Erkennung“; LH-FA-09.a; Plan §6 Risiko 5; Reviewer-Skill HIGH „Korrektheitsfehler im kritischen Pfad (… unpassende Antwort)“; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `spec/spezifikation.md:384`–`:395`; `internal/hexagon/services/lebendpruefung.go:5`–`:7`; `internal/hexagon/services/replay.go:72`–`:84`; `docs/plan/planning/in-progress/slice-extended-query-lebendpruefung.md:121` | ja (Sonde S1) | Spec-Regel an eine Serverversion gebunden, ohne sie zu nennen |
| F-351 | MEDIUM | Der Folge-Slice `slice-v1-abschluss-sessions` ist nicht nachgezogen. Dieser Slice hat LH-FA-12.a um einen Absatz zu Lebendprüfungen ergänzt; er erbt diese Sätze. Sein DoD-Punkt 1 sagt weiter „eine Anfrage darüber hinaus ist ein Mismatch“. Für eine Lebendprüfung gilt das jetzt nicht mehr. Kopf, §1 und §6 nennen weder LH-FA-09.a noch [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md). Unerwähnt bleibt auch die Kopplung an `--session-assignment connection`, die jener Slice liefert. `NewReplayService` nimmt Sessions nur aus Lebendprüfungen schon beim Laden aus `frei`. LH-FA-12.a sagt für `connection` aber, die n-te Verbindung erhalte die Session mit der `id` n, und nennt den Fall einer Session nur aus Lebendprüfungen nur für `first-request`. Plan §1 verweist zwar für Sonde S8 auf jenen Slice, die neuen Sätze reisen aber nicht mit. Hard Rule 3.9 verlangt, betroffene Folge-Slices im selben Commit nachzuziehen. Dasselbe Muster war viermal LOW (F-291, F-297, F-307, F-315) und zuletzt MEDIUM (F-323). | `AGENTS.md` §3.9; `BEO-REPO/plan-folgt-korrektur-nicht` (verkörpert); Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“ | `docs/plan/planning/open/slice-v1-abschluss-sessions.md:14`–`:16`, `:49`; `spec/spezifikation.md:466`–`:484`; `internal/hexagon/services/replay.go:60`–`:63` | ja (Lesen) | Folge-Slice nach Spec-Ergänzung nicht nachgezogen |
| F-352 | LOW | Das `Schärft:`-Feld von [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) nennt LH-FA-09.a, LH-FA-18.a, LH-FA-12.a, LH-FA-03.b und `SPEC-011`. LH-FA-03.a (Schritt 4) und LH-FA-10.a fehlen, obwohl `7ffd029` beide nach dieser ADR ändert und die Historie der Spezifikation sie unter demselben Eintrag nennt. Inhaltlich tragen die Änderungen: Sie folgen aus „löst keine Session-Zuordnung aus“ und „auch vor der ersten und nach der letzten“ der ADR, ohne diese zu erweitern. Die Aufwärts-Deklaration ist aber unvollständig, und die angenommene ADR ist unveränderlich. Wer von LH-FA-03.a oder LH-FA-10.a nach der verbindlich machenden ADR sucht (`make doc-trace`), findet sie nicht. Die Folgepflicht der ADR nennt die beiden Stellen ebenfalls nicht. | `AGENTS.md` §3.4, §3.5; ADR-Index §Konventionen (`Schärft:` aufwärts); Maintainability | `docs/plan/adr/0031-lebendpruefungen-im-replay.md:11`, `:50`; `spec/spezifikation.md:146`–`:150`, `:425`–`:428`, `:1593` | ja (Lesen, `make doc-trace`) | Schärft-Feld deckt die geänderten Spec-Stellen nicht |
| F-353 | LOW | Die Geschichte von [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) führt nur die Zeile „Proposed“, obwohl der Status `Accepted` ist. Jede andere angenommene ADR im Verzeichnis trägt eine Zeile „Accepted“. `4abad8e` ändert nur Status und Index. | ADR-Vorlage (Geschichte); Maintainability | `docs/plan/adr/0031-lebendpruefungen-im-replay.md:3`, `:63`–`:67` | ja (Lesen) | Zustandswechsel der ADR ohne Geschichtszeile |
| F-354 | LOW | Der Satz „Das Einspielen (`play`) führt Sessions ohne Interaktion nicht aus“ steht jetzt am Ende des neuen Absatzes über Lebendprüfungen in LH-FA-12.a. Direkt davor heißt es: „Eine Session, deren Interaktionen alle Lebendprüfungen sind, ist eine Session ohne Interaktion“. Zusammen gelesen hieße das, `play` führe Sessions nur aus Lebendprüfungen nicht aus. LH-FA-09.a sagt dagegen: „Record und Einspielen behandeln Lebendprüfungen wie jede andere Anfrage“, ebenso Plan §1. Die genauere Stelle löst den Widerspruch, der Satz in LH-FA-12.a lädt aber den späteren Slice zum Einspielen zur falschen Lesart ein. | LH-FA-12.a; LH-FA-09.a; Plan §1 (Record und `play` unverändert); Maintainability | `spec/spezifikation.md:480`–`:484`, `:411`–`:412` | ja (Lesen) | Satz durch Umstellung in neuen Bezug gerückt |
| F-355 | LOW | Randform 2: Eine Aufzeichnung nur aus Lebendprüfungen endet beim Start mit `PGR-E3004`. LH-FA-09.a sagt das jetzt ausdrücklich. Zwei Texte zu `PGR-E3004` passen dazu nicht. LH-FA-03.a Schritt 2 sagt „ein Recording ohne Session ist nicht verwendbar“; die Aufzeichnung hat aber Sessions. Das Handbuch erklärt den Code mit „Beim Aufzeichnen hat keine Verbindung eine Anfrage gestellt“; die Verbindungen haben aber Lebendprüfungen gestellt. Die Meldungstabelle der Spezifikation („ohne Session, die sich im Replay zuordnen lässt“) passt. Zudem verlangt die ADR, vorhandene Aufzeichnungen blieben „ohne Umwandlung verwendbar“. Eine solche Aufzeichnung startete vorher und startet jetzt nicht mehr. Bewertung der Randform unter Schwerpunkt 3. | LH-FA-03.a; LH-FA-09.a; [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) §Konsequenzen; Maintainability | `spec/spezifikation.md:142`–`:143`, `:409`–`:410`; `docs/user/benutzerhandbuch.md:612`; `internal/hexagon/services/replay.go:65`–`:67` | ja (Lesen; `TestReplayLebendpruefungAufgezeichnet`) | Meldungstext deckt den neuen Auslöser nicht |
| F-356 | LOW | Plan §3 sagt zu: „die Meldung nach der letzten Interaktion nennt deren aufgezeichnete Nummer“. `0841891` ändert das an zwei Stellen, in `Query` und in `ClientMessage`. Geprüft ist nur `Query` (M8 rot). In `ClientMessage` bleibt die Mutation zurück auf `c.pos` grün (M9). Zudem sagt die Meldung „keine aufgezeichnete Interaktion mehr nach 4“, wenn die Aufzeichnung danach noch eine Lebendprüfung mit Nummer 5 enthält (Testfall in `TestReplayLebendpruefungAufgezeichnet`). Wörtlich ist das falsch; gemeint ist „keine zu erwartende“. | Plan §3; LH-FA-10.a; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` | `internal/hexagon/services/replay.go:131`–`:132`, `:184`–`:185`; `internal/hexagon/services/replay_lebendpruefung_test.go:99`–`:102`; `docs/plan/planning/in-progress/slice-extended-query-lebendpruefung.md:84` | ja (M9) | Zusage an zwei Stellen umgesetzt, an einer geprüft |
| F-357 | LOW | Der neue Satz in der Negative von [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) nimmt nur die eingehende Lebendprüfung aus („beantwortet das Replay sie wie PostgreSQL, ohne die Aufzeichnung zu verbrauchen“). Für die andere Richtung, die aufgezeichnete, die ausbleibt, sagt er nichts. Streng gelesen ist eine fachliche Anfrage an der Stelle einer aufgezeichneten Lebendprüfung „eine Anfrage ohne passende Aufzeichnung“, für die LH-FA-10 greift. Plan §6 Risiko 4 begründet den Lastenheft-Satz genau damit, dass eine ADR das Lastenheft nicht schärfen darf. Für das Überspringen trägt diese Begründung gleichermaßen, gelöst ist sie nur über die Lesart von „passend“ in LH-FA-09.a. Das kann tragen, ist aber nicht ausdrücklich bestätigt. Hinweis an den Spec-Verantwortlichen, ohne Entscheidung. | [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [LH-FA-10](../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage); Plan §6 Risiko 4; Source Precedence (ADR schärft nicht das Lastenheft) | `spec/lastenheft.md:332`–`:336`; `spec/spezifikation.md:404`–`:412` | ja (Lesen) | Lastenheft-Satz deckt nur eine Richtung der Ausnahme |
| F-358 | INFO | Randform 4, Stabilität der Integrationstests: dreimal grün mit Laufzeiten von 20 bis 21 s für die drei Tests (I0); `make test-integration` wird dadurch um etwa 20 s länger. Die Pausen von 1,5 s liegen 0,5 s über der Schwelle von pgx. Da `time.Sleep` nicht kürzer schläft, sind sie nach unten sicher. Die Bedingung `nMit > nOhne` trägt die Behauptung „Pausen lösen Lebendprüfungen aus“. Bei `database/sql` hängt sie daran, dass `stdlib` die erste Wiederverwendung prüft und den Lauf ohne Pausen nicht länger als 1 s ruhen lässt. Unter starker Last könnte der Lauf ohne Pausen eine zweite Prüfung erhalten und die Bedingung rot färben. Das ist ein Fehlalarm, kein falsches Grün. Nicht beobachtet. I1 bis I3 zeigen, dass die Tests die drei Kernmutationen fangen. | `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients`; [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) §Fitness Function | `test/integration/lebendpruefung_e2e_test.go:22`, `:78`–`:102` | ja (I0 bis I3) | — (kein Fehlermuster) |
| F-359 | INFO | Die Commit-Message von `0841891` nennt die Slice-Kennung doch. `slice-extended-query-lebendpruefung` steht als eigene Zeile am Ende des Rumpfs, dazu `LH-FA-09`, `LH-FA-18`, `LH-QA-02` und [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) im Titel. Die Hook-Regel ist damit doppelt erfüllt. `4abad8e`, `eb61782` und `7ffd029` nennen `LH-*`-Kennungen, `f906c04` die Slice-Kennung. Keine Message nennt eine `SPEC-*`- oder `ARC-*`-Kennung. Die Prämisse im Auftrag, die Kennung fehle, trifft nicht zu. | `harness/README.md` §Traceability rules; [ADR-0025](../plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md); `AGENTS.md` §5 Regel 1 | `git show -s 0841891` | ja (`make hook-gegenprobe`, Lesen) | — (kein Fehlermuster) |

## Antwort auf die Schwerpunkte

1. **Erkennung gegen den Scanner von PostgreSQL.** Für PostgreSQL 17 deckt sich `istLebendpruefung` mit `scan.l`. Leerraum `[ \t\n\r\f\v]`. Zeilenkommentar `--` bis `[\n\r]` oder Textende. Blockkommentare verschachtelt: Jedes `/*` öffnet eine Ebene, auch wenn es mit einem `*/` überlappt, etwa `/*/`, und `\*+\/` schließt, auch `**/`. `--` hat im Blockkommentar keine Bedeutung. Ein offener Kommentar ist keine Lebendprüfung, `;` ebenso wenig. Der leere Text ist eine. Nicht-ASCII ist byteweise ausgeschlossen; PostgreSQL liest solche Bytes als Bezeichnerbeginn, also als Anweisung. Jede dieser Regeln fängt ein Tabellenfall (L1 bis L8), und der Integrationstest prüft die Tabelle gegen den Server. **Randform 1** (keine Serverversion in der Spezifikation) trägt nicht. Sonde S1 zeigt, dass PostgreSQL 16 `\v` ablehnt. Die übrigen Regeln verhalten sich in 16 gleich (S1, Stichprobe). Ob eine Versionsangabe, das Streichen von `\v` oder ein Blick auf die aufgezeichnete Antwort folgt, entscheidet der Spec-Verantwortliche (F-350).
2. **Transaktionsstatus.** Der Status kommt aus dem letzten `ReadyForQuery`, das die Verbindung erhalten hat: nach `Query` aus der Antwort der Interaktion, nach `ClientMessage` aus jeder Gruppe, Flush-Gruppen ohne `ReadyForQuery` lassen ihn stehen. Nach dem Handshake ist er `I`. Fehler in einer Transaktion ergeben `E` aus der Aufzeichnung (M1 bis M3, M12, I3 rot). Der Zustand liegt je Verbindung im Cursor. Eine neue Verbindung beginnt mit `I`, auch wenn eine andere in einer Transaktion steht. Eine Verbindung, die nur Lebendprüfungen stellt, bleibt ohne Session und erhält immer `I`, auch wenn keine Session mehr frei ist (M10 rot). Nach einem Mismatch ändert sich der Status nicht, die Verbindung endet ohnehin (`SPEC-018`).
3. **Überspringen aufgezeichneter Lebendprüfungen.** Der Cursor läuft über die gefilterte Liste. `first-request` überspringt Sessions nur aus Lebendprüfungen, und weder `CloseConnection` noch `Unassigned` zählt sie als nicht verbraucht, was `PGR-W2001` und damit `PGR-E5002` unter `--fail-on-unconsumed` erbt (M6 rot). Diagnosen nennen die aufgezeichnete `sequence` (M8 rot; Lücke M9 in F-356). **Randform 2** (`PGR-E3004` bei einer Aufzeichnung nur aus Lebendprüfungen) ist eine schlüssige Ableitung, aber eine Ableitung. Die ADR zählt als Wirkung des Lesens „als stünden sie nicht darin“ nur Cursor, Session-Zuordnung und nicht verbrauchte Interaktionen auf. Die Startprüfung in LH-FA-03.a gehört nicht zu den dreien. Mit der Ableitung spricht: Ohne sie startete ein Replay, das nichts wiedergeben kann; das ist für einen Testlauf eher ein stilles Grün als eine Hilfe. Gegen sie spricht der Satz der ADR, vorhandene Aufzeichnungen blieben ohne Umwandlung verwendbar. Die Texte zu `PGR-E3004` sind nicht nachgezogen (F-355). Ob die Ableitung trägt, entscheidet der Architect beziehungsweise der Spec-Verantwortliche.
4. **Spec-Änderungen über den Auftrag hinaus.** LH-FA-03.a Schritt 4 und LH-FA-10.a sind nötig. Ohne sie widersprächen „Anfrage ohne freie Session ist `PGR-E5003`“ und „Anfrage nach der letzten Interaktion ist ein Mismatch“ der Antwort vor der ersten und nach der letzten Interaktion. Inhaltlich decken sie sich mit der ADR, nur das `Schärft:`-Feld nennt sie nicht (F-352). Die Spezifikation folgt dem neuen Lastenheft-Satz und der ADR, die Schärft-Richtung stimmt: Die ADR zeigt aufwärts, die Spezifikation nennt keine ADR. Die ADR trägt Entscheidung und Gründe und verweist für Zeichen und Zählung auf die Spezifikation (§3.8). Abweichungen: F-354 (`play`-Satz in LH-FA-12.a), F-355 (`PGR-E3004`), F-357 (Lastenheft nur eine Richtung).
5. **Stabilität der Integrationstests.** Stabil in drei Läufen, etwa 20 s Zusatzlaufzeit. **Randform 4** (`nMit > nOhne` für `database/sql`) trägt. Der Kommentar nennt den Grund, die Prüfung bei der ersten Wiederverwendung, zutreffend. Ein Fehlalarm unter Last ist möglich, aber nicht beobachtet (F-358).
6. **Register-Lehren.** Die Zusagen aus LH-FA-09.a und LH-FA-18.a zu Lebendprüfungen fängt je ein Test mit Mutation (M1 bis M8, M10, M12, L1 bis L8, I1 bis I3). Ausnahme: die Nummer in der `ClientMessage`-Meldung (F-356). Weiter zusagt, als geprüft wird, nur ein Kommentar: der an `leerraum` (F-350). Plan-Kopf, §1, §3 und §6 folgen dem Code. Die eingetretene Lage von Risiko 5 steht noch als „offen“, was bis zur Closure zulässig ist. Nicht nachgezogen ist der Folge-Slice `slice-v1-abschluss-sessions` (F-351).
7. **Commit-Hook.** Erfüllt; die Kennung steht in `0841891` (F-359).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services` | geprüft. Erkennung als reine Funktion ohne SQL-Parser ([ADR-0008](../plan/adr/0008-kein-sql-parser-in-v1.md) unberührt); Antwort außerhalb der Reihe vor der Zuordnung; Sperre über den ganzen Aufruf; Filter nur für `Query`-Interaktionen, ein leerer `Parse` bleibt verglichen (M7 rot). Befund F-350, F-356. |
| `internal/hexagon/ports/driving` | geprüft, ohne Befund. Nur der Kommentar an `Query`; er sagt zu, was `TestReplayLebendpruefungAusserDerReihe` prüft. |
| Core-Reinheit | geprüft, ohne Befund. `lebendpruefung.go` importiert nur `strings`; kein Adapter, kein `pgproto3`. |
| `internal/adapters/driving/pgwire` | geprüft, ohne Befund. Nicht im Diff; `EmptyQueryResponse` übersetzt der Adapter schon, wie Plan §3 sagt (I0 grün). |
| `test/integration` | geprüft. Abdeckungs-Deklarationen sagen zu, was die Tests prüfen. Befund F-358. |
| `go.mod` | geprüft, ohne Befund. `puddle/v2` und `golang.org/x/sync` als indirekte Abhängigkeiten von `pgxpool`, nur im Integrationstest genutzt. |
| `docs/user/` | geprüft. Abdeckungstabellen entsprechen den Deklarationen (`make abdeckung-check`). Handbuch beschreibt Antwort, Überspringen, Grenze und `ShouldPing`. Befund F-355. |
| `spec/lastenheft.md` | geprüft. Kein Verweis auf ADR, Slice oder Spezifikation; Draft, kein Versionssprung nötig (`harness/conventions.md` §Versionierung). Befund F-357. |
| `spec/spezifikation.md` — Strata | geprüft. Keine ADR-, Slice- oder Commit-Nennung im Diff. Befund F-350, F-354, F-355. |
| `spec/architecture.md` — Hard Rule 3.4 | geprüft, ohne Befund. Nicht im Diff; `ARC-002` im Kopf des Plans ist nur Bezug. |
| `docs/plan/adr/` — Hard Rules 3.5, 3.8 | geprüft. `4abad8e` ändert nach der Annahme nichts am Inhalt; nach `4abad8e` kein Commit an der ADR. Die ADR trägt Entscheidung, Alternativen und Gründe, die Einzelregeln stehen in der Spezifikation. Index und Konventionstabelle nachgezogen. Befund F-352, F-353. |
| `docs/plan/planning/` — Hard Rule 3.9 | geprüft. Kopf, §1, §3 und §6 folgen dem Code. Befund F-351. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. Kein Move im Diff. |
| Hard Rule 3.6 — Gates nicht lockern | geprüft, ohne Befund. Keine Gate-Konfiguration im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft. Neue Kommentare sind Zusagen, Kopplungen oder Grenzen im Indikativ. Die Grenze an `objekte` (das Entfernen des unbenannten Statements durch eine einfache Anfrage) trifft zu. Befund F-350 (Kommentar an `leerraum`). |
| Commit-Messages `f906c04` bis `0841891` | geprüft. Befund F-359 (INFO, erfüllt). |
| Register `BEO-REPO/*` | geprüft. F-351 ist ein weiterer Fall von `BEO-REPO/plan-folgt-korrektur-nicht` (verkörpert, §3.9; Retirement-Check betroffen). F-350 und F-356 betreffen `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, F-356 auch `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`. Randform 2 (F-355) ist ein Kandidat für `BEO-REPO/spec-randform-erst-im-review-entschieden`, der Implementer hat sie aber selbst gemeldet. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 6 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:**

- Spec-Regel an eine Serverversion gebunden, ohne sie zu nennen (F-350; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Folge-Slice nach Spec-Ergänzung nicht nachgezogen (F-351; Register `BEO-REPO/plan-folgt-korrektur-nicht`)
- Schärft-Feld deckt die geänderten Spec-Stellen nicht (F-352)
- Zustandswechsel der ADR ohne Geschichtszeile (F-353)
- Satz durch Umstellung in neuen Bezug gerückt (F-354)
- Meldungstext deckt den neuen Auslöser nicht (F-355)
- Zusage an zwei Stellen umgesetzt, an einer geprüft (F-356; Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`)
- Lastenheft-Satz deckt nur eine Richtung der Ausnahme (F-357)

## Verdikt

**Merge-blockierend:** ja, wegen F-350. Der Auslöser ist selten: eine Anfrage mit `\v`, gegen PostgreSQL 16 oder älter aufgezeichnet. Für diesen Fall liefert das Replay aber eine Antwort, die weder aus der Aufzeichnung stammt noch der des aufgezeichneten Servers entspricht. Davon abgesehen trägt der Kern: Erkennung, Antwort außerhalb der Reihe, Status und Überspringen fangen Mutationen auf Unit- und Integrationsebene.

**Übergabe:**

- F-350 an den Spec-Verantwortlichen und den Architect, mit diesem Report und Sonde S1 als Artefakt. Die Frage ist eine Spec-Entscheidung, keine Implementierungsfrage, und berührt die Aussage der ADR zu Fehlerantworten. Widerspricht der Implementer der Einstufung HIGH, läuft der Konflikt-Pfad über den Architect als Sequenz mit Übergabe-Artefakten (Modul 8).
- F-351 an Implementer und Planner, weil `slice-v1-abschluss-sessions` mitbetroffen ist.
- F-352 und F-353 an den Architect: Die ADR ist angenommen, ob die Nachträge als Geschichtszeile zulässig sind oder eine Folge-ADR verlangen, entscheidet er.
- F-354, F-355 und F-357 an den Spec-Verantwortlichen; F-355 zusammen mit der Bewertung von Randform 2.
- F-356 an den Implementer. F-358 und F-359 ohne erwartete Aktion.
- Bei MEDIUM und darunter genügt eine Begründung im Plan oder in der Closure-Notiz.
- Die Finding-Klassen gehen in die Closure-Notiz §7. Die DoD-Konformität prüft der Verifier separat.
