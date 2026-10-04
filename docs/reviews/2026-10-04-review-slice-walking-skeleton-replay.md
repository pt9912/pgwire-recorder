# Review-Report: slice-walking-skeleton-replay — 2026-10-04

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD ist nicht Gegenstand dieses Reviews; sie prüft der Verifier.

**Gegenstand:** `git diff c30b3f1 HEAD` ohne `docs/reviews/` (HEAD `c66cd55`): `3dbc1cf` (Verantwortlich gesetzt), `d2fb076` und `af9faa3` (reine Moves `open/` → `next/` → `in-progress/`), `c66cd55` (Port `Replayer`, Replay-Service, PGWire-Adapter für Record und Replay, Kommando `replay`, YAML-Leser mit `offset_ms` und `empty_sessions`, Unit- und E2E-Tests, Abdeckungstabellen, README, Plan).

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-04

**Eingangs-Kontext:**

- `slice-walking-skeleton-replay` (§1 bis §8), `welle-walking-skeleton.md` (Abgrenzung), `welle-replay-semantik.md`, die offenen Slices `slice-replay-semantik-mismatch`, `slice-v1-abschluss-sessions`, `slice-replay-semantik-meldungscodes` (je §1 bis §3)
- [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0003](../plan/adr/0003-pgwire-server-ist-driving-adapter.md), [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md), [ADR-0007](../plan/adr/0007-strict-replay.md), [ADR-0008](../plan/adr/0008-kein-sql-parser-in-v1.md), [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)
- `spec/lastenheft.md` LH-FA-03, LH-FA-05, LH-FA-07, LH-FA-09, LH-FA-10, LH-FA-12, LH-FA-13, LH-QA-01, LH-QA-05, Abnahmeszenarien 2 und 4
- `spec/spezifikation.md` LH-FA-03.a, LH-FA-03.b, LH-FA-05.a, LH-FA-05.b, LH-FA-09.a, LH-FA-10.a, LH-FA-12.a, LH-FA-13.b, `SPEC-011`, `SPEC-013` bis `SPEC-028`, `SPEC-034`
- `spec/architecture.md` §2 (Schichten), §2.3 (Ports), §2.4, §4 (Sequenz Replay und Mismatch), §4.1, §4.2, §4.4, §4.6, §5
- `AGENTS.md` (Hard Rules 3.1 bis 3.8), `harness/README.md`, `harness/conventions.md`
- die Vorlauf-Reports am selben Modul: `2026-10-04-review-slice-walking-skeleton-record.md` (F-230 bis F-253) und `…-record-folge.md` (F-254 bis F-272)

**Ausgeführte Läufe im Repo:** `make build` (grün), `make test` (grün), `make a-check` (0 Befunde; Hinweis wie F-248, jetzt auch für `test/integration/replay_e2e_test.go`), `make a-check-negativ` (grün), `make abdeckung-check` (grün), `make docs-check` (0 Befunde, 126 Dateien), `make test-integration` (grün, 13 Tests, darunter `TestE2EReplaySelect1`, `TestE2EReplayAbweichung`, `TestE2EReplayBeschaedigt`). Die Image-IDs von `pgwire-recorder:{build,dev,test,integration}` blieben unverändert. `make gates` lief nicht.

**Mutationen, Beobachtungen, Sonden:** Nur an einer Kopie (`git archive HEAD`) im Scratchpad. Im Runner-Skript der Kopie sind Image-Tag und Netzname umbenannt (`pgr-review-rp:*`, `pgr-rv-*`). Je Mutation liefen `docker build --target test` und, wo angegeben, der Runner gegen das gepinnte PostgreSQL-Image. Ein Kontrolllauf mit einer reinen Kommentaränderung lief grün/grün. Alle Images, Container und Netze des Reviews sind gelöscht.

| # | Mutation | Unit | Integration |
|---|---|---|---|
| R1 | Cursor rückt auch bei Mismatch vor | rot | — |
| R2 | Vergleich mit Whitespace-Normalisierung (`strings.Fields`) | rot | — |
| R3 | Handshake immer aus der letzten statt aus der nächsten freien Session | **grün** | — |
| R4 | Anfrage nach Ende der Session erhält die Antworten der letzten Interaktion | rot | — |
| R5 | `CloseConnection` warnt nie | rot | — |
| R6 | `Unassigned` meldet nie | rot | — |
| R7 | Ohne freie Session wird die letzte erneut zugeordnet statt `PGR-E5003` | rot | — |
| R8 | Adapter verschluckt die Warnung von `CloseConnection` | rot | grün |
| R10 | `replay` endet immer mit 0 | grün | rot |
| R11 | Bootstrap loggt die Warnung für nie zugeordnete Sessions nicht | **grün** | **grün** |
| R12 | Mismatch-Text ohne Interaktionsnummer, erwartete und empfangene Query (`"Session %d: Abweichung"`) | **grün** | — |
| R13 | `PGR-E3004` entfällt | rot | — |
| R14a | Leser: Prüfung auf negatives `offset_ms` entfernt | **grün** | — |
| R14b | Leser: `empty_sessions` wird ignoriert, eine leere Session bleibt `PGR-E3003` | **grün** | — |
| R14c | Leser: `offset_ms` wird nicht ins Modell übernommen | **grün** | — |
| R14d | Leser: `empty_sessions` wird nicht ins Modell übernommen | **grün** | — |
| R15 | `ParameterStatus` ohne Sortierung (Map-Reihenfolge) | rot in 1 von 3 Läufen | — |
| R16 | Session wird schon beim Verbindungsaufbau zugeordnet (Zählung `connection` statt `first-request`) | rot | — |

*Sonde psql/libpq:* Ein Container mit `postgres:17-alpine` lieferte PostgreSQL und `psql`. `record` zeichnete `SELECT 1;` und `SELECT 'a' AS t;` auf. Danach wurde der PostgreSQL-Container entfernt und `replay` mit derselben Datei gestartet. Ergebnisse:

- `psql` (sslmode `disable` und `prefer`) baut die Session auf und erhält dieselben Ergebnisse. Der Handshake besteht aus `AuthenticationOk`, 14 `ParameterStatus` und `ReadyForQuery`, ohne `BackendKeyData`.
- Eine zweite Verbindung mit Anfrage erhält eine FATAL-Antwort mit `PGR-E5003`.
- Innerhalb der Session erhält `SELECT 'b' AS t;` statt `SELECT 'a' AS t;` die FATAL-Antwort `Replay [PGR-E5001]: Session 1, Interaktion 2: erwartet "SELECT 'a' AS t;", empfangen "SELECT 'b' AS t;"`. Danach folgt die Warnung `code=PGR-W2001` „Session 1: 1 von 2 Interaktionen nicht verbraucht“ ohne Code im Text.
- Exit-Code nach SIGTERM: 5.

*Sonde Startfehler* (Binary, `--input` auf vorbereitete Dateien):

- `sessions: []` → `PGR-E3004`, Exit 3
- nur leere Session mit `empty_sessions: true` → `PGR-E3004`, Exit 3
- leere Session ohne Kennzeichnung → `PGR-E3003`, Exit 3
- `version: 9` → `PGR-E3002`, Exit 3
- `offset_ms: -5` → `PGR-E3003`, Exit 3
- fehlende Datei → `PGR-E3001`, Exit 3
- `offset_ms: 12` → startet

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-273 | MEDIUM | Drei Zusagen des Replay-Pfads hält kein Test; die Mutationen R3, R11 und R12 bleiben in Unit- und Integrationstests grün. (1) Der Handshake einer Verbindung nimmt die Serverparameter der nächsten noch nicht zugeordneten Session (LH-FA-12.a). Kein Test prüft das mit mehr als einer Session. (2) Die Mismatch-Diagnose nennt Session, erwartete Interaktionsnummer, erwartete und empfangene Query (LH-FA-10.a). Die Tests prüfen nur den Code `PGR-E5001`. (3) Am Laufende wird die Warnung `PGR-W2001` für nie zugeordnete Sessions geloggt (LH-FA-03.b, LH-FA-13.b §Prozessende). Den Aufruf in `replay` prüft kein Test. Die Sortierung der `ParameterStatus` fängt `TestReplayStrictSequential` nur zufällig (R15: 1 von 3 Läufen rot), weil die Reihenfolge bei zwei Parametern an der Map-Iteration hängt. Drittes Auftreten der Klasse (F-234, F-260): Pflege-Schwelle des Skills. | LH-FA-10.a; LH-FA-12.a; LH-FA-03.b; LH-FA-13.b; Reviewer-Skill MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag“ | `internal/hexagon/services/replay_test.go` · `TestReplayMismatch`, `TestReplaySessionZuordnung`, `TestReplayStrictSequential`; `internal/bootstrap/bootstrap.go` · `replay` (Aufruf von `Unassigned`); `test/integration/replay_e2e_test.go` · `TestE2EReplayAbweichung` | ja — Mutationen R3, R11, R12, R15 | Protokollrand ohne Negativtest |
| F-274 | MEDIUM | Der Plan beschreibt einen kleineren Umfang als der Diff. §1 nimmt „Mismatch-Diagnose und Exit-Code `5`“ (welle-replay-semantik) und „Mehrere aufgezeichnete Sessions“ (welle-v1-abschluss) aus. §3 desselben Plans nennt aber „Abweichung mit Exit-Code 5“ und „Zuordnung `first-request`“. Der Code liefert dazu `PGR-E5001` mit Diagnose, `PGR-E5003`, die Zuordnung mehrerer Sessions und `PGR-W2001`. `Bezug` und `Berührte Spec-Stellen` nennen weder LH-FA-10, LH-FA-12 und LH-FA-13 noch LH-FA-03.b, `SPEC-018`, `SPEC-027` und `SPEC-034`. Die Commit-Message nennt LH-FA-10. `docs/user/abdeckung-gesamt.md` führt LH-FA-10 als vollständig belegt. `welle-walking-skeleton.md` grenzt dieselben Punkte aus. Die offenen Slices `slice-replay-semantik-mismatch` und `slice-v1-abschluss-sessions` planen den gelieferten Teil weiter. §3 wurde nachgezogen, §1 nicht. Das ist das Muster von F-241, das vierte Auftreten der Klasse nach F-215, F-241 und F-266, die beiden letzten davon LOW. Damit gilt die MEDIUM-Regel des Skills; die Pflege-Schwelle wurde schon bei F-266 erreicht. | Maintainability; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“; `v6.13.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice | Slice-Plan §1, Kopf (`Bezug`, `Berührte Spec-Stellen`); `docs/plan/planning/welle-walking-skeleton.md` §Abgrenzung; `docs/plan/planning/open/slice-replay-semantik-mismatch.md` §1/§2; `docs/plan/planning/open/slice-v1-abschluss-sessions.md` §1/§2 | ja (Lesen gegen Diff) | Plan-Zeile nach Korrektur nicht nachgezogen |
| F-275 | MEDIUM | Der Leser nimmt zwei Felder des Formats Version 1 neu an: `offset_ms` mit Ablehnung negativer Werte und `empty_sessions`, das leere Sessions zulässt. Es gibt keinen Test dafür; `yaml_test.go` ist im Diff unverändert. Die vier Mutationen R14a bis R14d bleiben grün: negative Werte angenommen, Kennzeichnung ignoriert, `OffsetMS` oder `EmptySessions` nicht ins Modell übernommen. Plan §6 nennt genau diese Felder als Risiko. | LH-FA-07 Negative; LH-FA-12.a; LH-FA-21.a; `SPEC-041`; Reviewer-Skill MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag“ | `internal/adapters/driven/recording/yaml.go` · `recordingDTO`, `interactionDTO`, `fromDTO`, `toDTO`; `internal/adapters/driven/recording/yaml_test.go` | ja — Mutationen R14a bis R14d | Formatfeld ohne Test |
| F-276 | LOW | Laut Kommentar am Typ `Server` ist „genau einer gesetzt“ (`Recorder` oder `Replayer`). Der Code prüft das nicht. Sind beide gesetzt, gewinnt `Replayer` still. Ist keiner gesetzt, führt der erste Verbindungsaufbau zu einem nil-Zugriff. Die Zusage ist weiter als der Code. Zweites Auftreten der Klasse (F-262). | `AGENTS.md` §3.7; Maintainability | `internal/adapters/driving/pgwire/server.go` · Typ `Server` (Kommentar), `open`, `query`, `close` | ja (Lesen) | Zusage weiter als Code |
| F-277 | LOW | Die Warnung `PGR-W2001` reist als `*model.Error` durch Port und Bootstrap (`CloseConnection`, `Unassigned`). Für diesen Wert liefert `ExitCode()` 2, denn die erste Ziffer ist bei Warnungen der Bereich. `Error()` ergäbe „Konfiguration [PGR-W2001]: …“. `SPEC-034` sagt: „Eine Warnung trägt keine Klasse und ändert den Exit Code nicht.“ Heute ist das nicht sichtbar, weil nur `Msg` und `Code` geloggt werden. Der Typ-Kommentar sagt „klassifizierter Fehler“. | `SPEC-034`; Maintainability | `internal/hexagon/ports/driving/replay.go` · `CloseConnection`; `internal/hexagon/services/replay.go` · `CloseConnection`, `Unassigned`; `internal/hexagon/model/fehler.go` · `Error`, `ExitCode` | ja (Lesen; `(&model.Error{Code: "PGR-W2001"}).ExitCode()` = 2) | Warnung im Fehlertyp |
| F-278 | LOW | Der Port `Replayer` verwendet `model.SessionID` als Nummer der Client-Verbindung. Dieselbe Typbezeichnung trägt sonst die Kennung einer aufgezeichneten Session. Die Diagnosen zählen in zwei Räumen mit demselben Typ: `PGR-E5003` „Verbindung 2“, `PGR-E5001` „Session 1“. Wer die Stelle ändert, sieht am Typ nicht, welcher Raum gemeint ist. | Maintainability; `spec/architecture.md` §2.3 (Ports in fachlichen Begriffen) | `internal/hexagon/ports/driving/replay.go`; `internal/hexagon/services/replay.go` · `verbindung`, `next` | ja (Lesen) | Typ trägt zwei Bedeutungen |
| F-279 | INFO | Plan §6, Risiko 1 bleibt offen. Der Replay-Handshake besteht aus `AuthenticationOk`, aufgezeichneten `ParameterStatus` und `ReadyForQuery`, ohne `BackendKeyData`, unabhängig von `user`/`database` des Clients. Beobachtet sind `pgconn` (E2E) und `psql` 17/libpq (Sonde). Andere typische Treiber (JDBC, Npgsql, asyncpg) sind nicht geprüft. Record reicht `BackendKeyData` ebenfalls nicht weiter (`SPEC-004`). Ob das für LH-FA-05.b „typische PostgreSQL-Treiber“ genügt, beurteilt der Validator. | LH-FA-05.b; `spec/architecture.md` §4.4; Verweis Validator | `internal/hexagon/services/replay.go` · `OpenConnection`; `internal/adapters/driving/pgwire/server.go` · `handle` | ja — Sonde je Treiber | Treiber-Kompatibilität nur an zwei Clients belegt |
| F-280 | INFO | `TestE2EReplaySelect1` deklariert LH-FA-03/Happy („kein erreichbarer PostgreSQL-Server“), aber PostgreSQL läuft im selben Netz weiter. Die Bedingung gilt nur strukturell, weil `replay` keine Upstream-Option hat. Die Sonde mit entferntem PostgreSQL-Container bestätigt das Verhalten. Ob die DoD-Zeile „mit gestoppter PostgreSQL“ damit belegt ist, prüft der Verifier. LH-FA-03/Boundary ist zu Recht nicht deklariert. | LH-FA-03; [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md); Verweis Verifier | `test/integration/replay_e2e_test.go` · `TestE2EReplaySelect1`; `docs/user/abdeckung-e2e.md` | ja — derselbe Test mit beendetem PostgreSQL-Container | Vorbedingung des Kriteriums nur strukturell hergestellt |
| F-281 | INFO | Den ersten Teil der Bedingung aus §4.1 (gleicher Anfrage-Typ) sichert der Leser, nicht der Matcher. `fromDTO` lehnt jeden Typ außer `query` ab, auch das für Version 1 spezifizierte `extended`, mit `PGR-E3003` „Anfrage-Typ "extended" unbekannt“. Der Service vergleicht nur den SQL-Text. Plan §6, Risiko 2 („lesen oder gezielt ablehnen“): Die Ablehnung erfolgt über den allgemeinen Pfad für beschädigte Aufzeichnungen. Hinweis an Implementer (Ausgang des Risikos) und an `slice-extended-query-replay`. | `spec/architecture.md` §4.1; [ADR-0007](../plan/adr/0007-strict-replay.md); LH-FA-18 | `internal/adapters/driven/recording/yaml.go` · `fromDTO`; `internal/hexagon/services/replay.go` · `Query` | ja (Lesen) | Matcher-Bedingung an anderer Schicht verankert |
| F-282 | INFO | `--fail-on-unconsumed` (LH-FA-03.a, LH-FA-03.b, `PGR-E5002`) plant kein Slice unter `docs/plan/planning/open/`. Zwei Slices verweisen den Punkt an welle-v1-abschluss, dort nennt ihn keiner. Dieser Slice liefert die Warnung `PGR-W2001` ohne die Option. Hinweis an den Planner. | LH-FA-03.b; `SPEC-012` | `docs/plan/planning/open/` (kein Treffer für `fail-on-unconsumed`) | ja (Lesen) | Spec-Option ohne planenden Slice |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/model` — Core-Reinheit | geprüft. Einziger Import `fmt`; `EmptySessions`, `OffsetMS` als Modellfelder ohne Bibliothekstyp; neue Codes stimmen mit der Tabelle in `SPEC-034` überein. Befund F-277. |
| `internal/hexagon/ports/driving` — §2.3 Ports | geprüft. `Replayer` mit `context` und Model-Typen; keine Connection, kein `pgproto3`. Befund F-277, F-278. |
| `internal/hexagon/services` — Strict Matching (`SPEC-011`, LH-FA-09.a, §4.1, §4.2) | geprüft. Importe `context`, `sort`, `sync`, `model`, `ports/driven`; kein Driving Port, kein Logging, kein Prozessende. Vergleich `erwartet.Request.SQL != sql` byte-genau ohne Normalisierung (R2 rot); Cursor rückt nur bei Gleichheit vor (R1 rot); identische Queries nach Position (`TestReplayStrictSequential`); Anfrage nach Ende → `PGR-E5001` (R4 rot); Cursor-Zustand je Verbindung im Service, der Adapter kennt ihn nicht ([ADR-0003](../plan/adr/0003-pgwire-server-ist-driving-adapter.md)). Kein SQL-Parser ([ADR-0008](../plan/adr/0008-kein-sql-parser-in-v1.md)). Befund F-273, F-281. |
| `internal/hexagon/services` — Zuordnung und Nebenläufigkeit (LH-FA-12.a) | geprüft. `first-request`: Zuordnung erst bei der ersten Anfrage aus `frei` in Reihenfolge der Kennung, Sessions ohne Interaktion übersprungen, Verbindung ohne Anfrage zählt nicht (R16 rot); keine freie Session → `PGR-E5003` (R7 rot); Handshake aus `frei[0]`, sonst `letzte`; alle Zugriffe unter `mu`; die zugeordnete Session ist eine Kopie, ihre Interaktionen werden nur gelesen. `PGR-W2001` bei Verbindungsende mit Rest (auch nach einem Mismatch, LH-FA-03.b wörtlich) und am Laufende für nie zugeordnete Sessions. `PGR-E3004` bei Aufzeichnung ohne zuordenbare Session (R13 rot, Sonde). Befund F-273. |
| `internal/adapters/driving/pgwire` — Replay-Modus (LH-FA-05.b, §4.4) | geprüft. Handshake ohne Upstream: `AuthenticationOk` vom Adapter, `ParameterStatus` und `ReadyForQuery` aus dem Use Case; Fehler des Use Case → FATAL-`ErrorResponse` mit Code im Text, SQLSTATE `XX000`, erster Fehler gemerkt, Verbindung endet (LH-FA-13.b); Warnung als Attribut `code`, Text ohne Code (`SPEC-034`). `pgproto3` bleibt im Adapter (`make a-check-negativ`). Befund F-276, F-279. |
| `internal/adapters/driven/recording` | geprüft. `KnownFields` bleibt; `empty_sessions` auf oberster Ebene, `offset_ms` je Interaktion, `omitempty` beim Schreiben; Sonden wie oben. Befund F-275, F-281. |
| `internal/adapters/driving/cli` | geprüft, ohne Befund. `replay` mit Pflichtoptionen `--listen`, `--input` (LH-FA-03.a); fehlende Option, Restargument → `PGR-E2001`; Hilfe in der Usage. |
| `internal/bootstrap` — Exit-Codes (LH-FA-13.b, `SPEC-016`, `SPEC-018`) | geprüft. Laden vor `Listen`, Startfehler über `fail` mit der Klasse des Codes (Sonden: E3001/E3002/E3003/E3004 → 3); Mismatch und `PGR-E5003` → 5 über `FirstErrorCode` (R10 rot, Sonde); Warnung nie zugeordneter Sessions nach `Serve`. Die Erfüllung von `Replayer` hält die Composition Root fest (§4.6). Befund F-273. |
| `test/integration`, Unit-Tests | geprüft. Record → zehn Replay-Läufe mit Vergleich von Spalten, OIDs, Zeilen und Befehlsabschluss (LH-QA-01-Messung); Abweichung mit `PGR-E5001` und Exit 5; beschädigte Aufzeichnung mit Exit 3. Befund F-273, F-280. |
| Abdeckungs-Deklarationen ([ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)) | geprüft gegen die Kriterien des Lastenhefts. LH-FA-09 Happy/Boundary/Negative, LH-FA-10 Happy/Boundary (Whitespace)/Negative, LH-FA-07/Happy (Neustart → Replay), LH-FA-13/Negative (Klasse 5), LH-QA-01/Messung (zehn Läufe), LH-FA-03/Negative treffen ihr Kriterium; LH-FA-03/Boundary und LH-FA-12 sind nicht deklariert, das ist ehrlich. Befund F-280; Plan-Bezug F-274. |
| Hard Rule 3.1 — Docker-only | geprüft, ohne Befund. Build, Test, Integration, Gates in Docker. |
| Hard Rule 3.2 — Suppression | geprüft, ohne Befund. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. `d2fb076` und `af9faa3` sind reine Renames (0 Zeilen), die Inhaltsänderungen stehen in `3dbc1cf` und `c66cd55`. |
| Hard Rule 3.5, 3.8 — ADRs | geprüft, ohne Befund. Keine ADR im Diff. |
| Hard Rule 3.6 — Gates nicht lockern | geprüft, ohne Befund. Kein Gate, keine Schwelle, keine Regel in `.a-check.yml`/`.d-check.yml` geändert. |
| Hard Rule 3.7 — Kommentare | geprüft. Neue Kommentare im Indikativ mit Zusage, Kopplung oder Abgrenzung; keine verworfene Alternative, kein abwesender Text, kein abgebrochener Satz. Befund F-276. |
| ADR-Konformität | geprüft. [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) (Schichten, `make a-check` 0 Befunde), [ADR-0003](../plan/adr/0003-pgwire-server-ist-driving-adapter.md) (Adapter ohne Fachlogik), [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md) (Model-Typen an der Portgrenze), [ADR-0007](../plan/adr/0007-strict-replay.md) und [ADR-0008](../plan/adr/0008-kein-sql-parser-in-v1.md) ohne Befund; [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) siehe F-280. |
| `README.md` | geprüft, ohne Befund. Stand deckt sich mit dem Code. |
| Commit-Message `c66cd55` | geprüft. Trägt `slice-walking-skeleton-replay` und `LH-*`, keine Struktur-IDs; Abweichung zum Plan-Bezug siehe F-274. |
| `spec/` | geprüft, ohne Befund. Nicht im Diff. |

## Summary

| Kategorie | Anzahl (neu) |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:**

- Protokollrand ohne Negativtest (3. Auftreten, vgl. F-234, F-260 — Pflege-Schwelle)
- Plan-Zeile nach Korrektur nicht nachgezogen (4. Auftreten, vgl. F-215, F-241, F-266)
- Formatfeld ohne Test
- Zusage weiter als Code (2. Auftreten, vgl. F-262)
- Warnung im Fehlertyp
- Typ trägt zwei Bedeutungen
- Treiber-Kompatibilität nur an zwei Clients belegt
- Vorbedingung des Kriteriums nur strukturell hergestellt
- Matcher-Bedingung an anderer Schicht verankert
- Spec-Option ohne planenden Slice

## Verdikt

**Merge-blockierend:** nein — kein HIGH. Der kritische Pfad hält. Der Vergleich ist byte-genau, ein Mismatch rückt den Cursor nicht vor, identische Queries werden nach Position verbraucht, und nach dem Ende der Session gibt es keine geratene Antwort. Die Zuordnung `first-request` entspricht LH-FA-12.a; Exit-Codes 3 und 5 entsprechen LH-FA-13.b. Belegt ist das durch Mutationen und Sonden, auch mit gestoppter PostgreSQL und `psql`.

Vor der Closure verdienen die drei MEDIUM Aufmerksamkeit. Bei F-273 und F-275 tragen die Diagnose nach LH-FA-10.a, die Handshake-Vorlage und die beiden neuen Formatfelder keinen Test. F-274 betrifft die Planung: Der Slice liefert Gegenstand zweier späterer Wellen, ohne dass §1, Bezug und die betroffenen Slices das sagen.

**Übergabe:**

- F-273, F-275, F-276, F-277, F-278 und F-281 an den Implementer.
- F-274 und F-282 an den Planner.
- F-279 an den Validator, F-280 an den Verifier.
- Steering-Loop: F-273 erreicht mit der Klasse „Protokollrand ohne Negativtest“ die Pflege-Schwelle (drittes Auftreten). F-274 wiederholt eine Klasse, die schon bei F-266 an der Schwelle stand. Für beide Klassen geht an den Architect die Frage, ob ein Gate (Mutationslauf, Abgleich Plan §1/§3 gegen den Diff) oder ein Eintrag in `AGENTS.md` sie verhindert hätte.
- Widerspricht der Implementer einer MEDIUM-Einstufung, genügt Annahme oder Begründung; der Konflikt-Pfad über den Architect gilt ab HIGH oder ab dem dritten gleichen Konflikttyp.
- Die Finding-Klassen gehen in die Slice-Closure §7. Die DoD-Konformität prüft der Verifier separat.
