# Verifikation: slice-walking-skeleton-record — 2026-10-04

**Rolle:** Verifier (Modul 11). Geprüft wird, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules prüft der Reviewer.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-walking-skeleton-record.md` <!-- d-check:ignore (Pfad beim geprüften Stand) --> (§1, §3, §6) gegen `git diff 41f7462 HEAD` bei HEAD `fea1182`. Das sind 13 Commits, 48 Dateien, +4080/−38.

**Eingang:** DoD-Liefer-Punkte 1 bis 3 und die Sensor-Angaben des Implementers. Dazu die Review-Reports `2026-10-04-review-slice-walking-skeleton-record.md` (F-230 bis F-253) und `2026-10-04-review-slice-walking-skeleton-record-folge.md` (F-254 bis F-272, Stand `0c99a84`).

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-04

Alle Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg. Mutationen und Sonden liefen nur in Kopien von `git archive HEAD` im Scratchpad, nie im Repo. Das Image dafür (`pgr-verifier:test`, Stufe `test` des Dockerfiles, ohne Cache gebaut) habe ich danach gelöscht. Die Images `pgwire-recorder:*` erzeugen die Make-Ziele des Repos; sie lagen schon vorher vor, ihre IDs sind unverändert. `git status --short` war nach allen Läufen leer.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — LH-FA-02, LH-FA-06: `SELECT 1;` über den Recorder, Ergebnis der realen Instanz, Interaktion geordnet im Recording (Integrationstest). **Bestätigt.**

Beleg aus dem eigenen Lauf von `make gates` (Teil `test-integration`). Das gepinnte Image `postgres:17-alpine@sha256:b0f9560a…` lief in einem eigenen, internen Docker-Netz:

| Test | Ergebnis |
|---|---|
| `TestE2ERecordSelect1` | PASS |
| neun weitere `TestE2ERecord*` (Upstream nicht erreichbar, mehrere Interaktionen, SSL abgelehnt, nicht unterstützt, Beenden mit offener Verbindung, Verbindung ohne Anfrage, Ende ohne Terminate, Schreibfehler am Ende, fremde Anfrage) | PASS |
| Runner | `run-integration-tests: gruen` |

`TestE2ERecordSelect1` arbeitet Black-Box. Es startet das Binary mit `record`, fragt mit `pgconn` `SELECT 1;` ab und prüft, dass das Ergebnis `1` ist. Nach SIGTERM erwartet es Exit-Code 0. Danach liest es die Datei und prüft drei Dinge: Formatkennung und Version stehen darin, genau eine Session mit einer Interaktion ist aufgezeichnet, und die Antworten stehen in der Reihenfolge `row_description` → `data_row` → `command_complete` → `ready_for_query`. Adressen und Pfade des Laufs dürfen nicht vorkommen. Der Test importiert nichts aus `internal/` (geprüft per `grep`).

**Eigene Mutationsprobe** (in der Kopie: Binary und Testbinary gebaut, gegen ein eigenes PostgreSQL gelaufen, nur `TestE2ERecordSelect1`):

| # | Mutation | Ergebnis |
|---|---|---|
| E1 | `toMessage` schickt dem Client statt des Wertes `9` (`server.go`) | rot: `Ergebnis über den Recorder …` |
| E2 | Service übernimmt die Interaktion nicht (`record.go`, `append` verworfen) | rot: `ready_for_query` fehlt, keine Interaktion |
| E3 | YAML-Adapter schreibt die Antworten in umgekehrter Reihenfolge | rot: `data_row … steht nicht in der Reihenfolge` |
| E4 | Service schreibt nie (weder in `CloseSession` noch in `Finish`) | rot: Aufzeichnung fehlt |
| E5 | Upstream-Adapter nimmt die Upstream-Adresse als `ParameterStatus` auf | rot: `rechnerspezifische Angabe "…-pg:5432"` |

Alle fünf Mutanten werden gefangen. Die unmutierte Kopie läuft grün. Der Punkt trägt.

### Punkt 2 — LH-FA-07: Recording persistent, mit Formatkennung und Version, per Roundtrip ladbar (Adapter-Contract-Test). **Bestätigt**, mit einer Reichweitengrenze.

Den Test-Teil habe ich ohne Build-Cache neu gefahren. Grund: In `make gates` war die Stufe `test` `CACHED` und hat ihre Ausgabe nicht gezeigt. Der Neulauf (`go vet -tags integration ./...` und `go test ./...`) ergab für alle fünf Packages mit Tests `ok`. Mit `-v` liefen `TestRoundtrip`, `TestMarshalDeterministisch`, `TestPrepareVorhandeneDatei`, `TestWriteRechteUndAtomar`, `TestUnmarshalNullUngequotet` und `TestUnmarshalFehler` (8 Unterfälle), alle PASS.

`TestRoundtrip` läuft über die Port-Methoden `Prepare` → `Write` → `Load` von `recording.YAML` gegen eine echte Datei. Es vergleicht per `reflect.DeepEqual` und prüft, dass `format: pgwire-recorder`, `version: 1`, `base64: AP8=` und `"null": true` im Text stehen und keine temporären Reste bleiben.

**Eigene Mutationsprobe** (`go test -run 'TestRoundtrip$'` in der Kopie):

| # | Mutation | Ergebnis |
|---|---|---|
| U1 | Binärwerte immer als `text` statt `base64` | rot (`base64: AP8=` fehlt) |
| U2 | `Write` benennt die temporäre Datei nicht um (nicht persistent) | rot (`Load: PGR-E3001 … no such file`) |
| U3 | `version` wird nicht geschrieben | rot (`Load: PGR-E3002 Version 0`) |
| U4 | `server_parameters` werden nicht geschrieben | rot (`Roundtrip weicht ab`) |
| U5 | Formatkennung im Modell auf `pgwire-rec` geändert (Schreiber und Leser gleich) | rot (`format: pgwire-recorder` fehlt) |
| U6 | NULL als leerer Text geschrieben | rot (`Roundtrip weicht ab`) |

Alle sechs Mutanten werden gefangen. Der Punkt trägt in der Form, in der ihn die DoD formuliert.

**Grenze:** Ein Port-übergreifender Contract-Test (eine Testsuite gegen `driven.RecordingRepository`, die jeder Adapter durchlaufen muss) existiert nicht. Der Test prüft den einzigen Adapter direkt. Für diesen Slice genügt das. Mit dem SQLite-Adapter (LH-FA-22) wird es eine Frage an den Planner. Außerdem lädt kein Test die Datei, die das Binary selbst geschrieben hat. Die E2E-Prüfung sucht nur nach Zeichenfolgen, der Roundtrip läuft über ein im Test gebautes Modell. Mehr dazu in §5, Befund V-3.

### Punkt 3 — `make gates` grün. **Bestätigt.**

| Kommando | Ausgabe (gekürzt) | Exit |
|---|---|---|
| `make gates` (HEAD `fea1182`, sauberer Baum, vor dem Anlegen dieser Datei) | `abdeckung-gegenprobe: gruen …` · a-check `gesamt: 0 Befund(e)` · `a-check-negativ: gruen …` · `baseline-verify: v6.13.0 OK — 54 Dateien` · `build` (Stufe `runtime`) · `test` · `d-check: 113 Datei(en) geprüft, 0 Befund(e)` · `commit-msg-gegenprobe: gruen …` · 10 × E2E PASS, `run-integration-tests: gruen` | 0 |
| `make abdeckung-check` (einzeln) | keine Ausgabe | 0 |
| `make doc-trace` | RTM: LH-FA-02, LH-FA-05, LH-FA-07, LH-FA-13, LH-QA-06 mit `Tests \| ok` | 0 |
| `docker run --rm --network none pgwire-recorder:dev version` | `pgwire-recorder dev`; Image-User `nonroot`, Entrypoint `/pgwire-recorder` | 0 |

Hinweis: Bei a-check bleibt `test/integration/record_e2e_test.go` „in keiner Schicht, ungeprüft“. Dass der Test nichts aus `internal/` importiert, hält also kein Gate; ich habe es per `grep` geprüft.

### Prozess-Punkte der DoD (nicht Teil des Auftrags, nur vermerkt)

- **Review:** Beide Reports liegen vor. Der Folge-Review deckt den Stand bis `0c99a84` ab. Die Commits `e29851a` (Korrekturen zu F-254 bis F-272, [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)), `bebac0c` ([ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) Accepted) und `fea1182` (`PGR-W3003`, Spec-Änderung LH-FA-05.e, `SPEC-045`) hat kein Review erfasst.
- Closure-Notiz, Register, Risiko-Ausgänge und Paarungen sind noch offen (§7 leer). Für `in-progress` ist das erwartbar.

## 2. Plan gegen Code-Diff

### §1 Ziel und Abgrenzung

| Plan §1 | Code | Befund |
|---|---|---|
| Ziel: `record` vermittelt `SELECT 1;` und schreibt ein versioniertes YAML-Recording | Punkt 1 und 2 | konform |
| Kein Replay; Roundtrip-Test genügt | kein `replay`-Kommando, `TestRoundtrip` | konform |
| Fehlerreplay ausgenommen; Fehlerantworten nur vermittelt und aufgezeichnet | `ErrorResponse`/`NoticeResponse` werden abgebildet (`toResponse`, `toMessage`) | konform |
| Parallele Sessions werden vermittelt und in der Reihenfolge ihres Endes übernommen; LH-FA-12 gilt nicht als belegt | `CloseSession` vergibt `id = len+1`; keine Deklaration für LH-FA-12 (F-261 behoben) | konform |
| Signalbehandlung: Ende nach der laufenden Interaktion, ein zweites Signal beendet sofort | `main.go`: `stop()` nach dem ersten Signal; `SetReadDeadline` beim Ende von ctx; `OpenSession`/`close` mit `WithoutCancel` | konform |
| Nur Upstream ohne Passwort-Anmeldung | `Upstream.Open`: Cleartext/MD5/SASL → `PGR-E6001` | konform (aber siehe V-2) |

### §3 Dateien

| Plan §3 | Diff | Befund |
|---|---|---|
| `internal/hexagon/model`, `services`, `ports/driving`, `ports/driven` — neu | `fehler.go`, `recording.go`, `response.go`, `record.go`, Ports `recording.go`, `upstream.go`, `record.go` | konform |
| `internal/adapters/driving/pgwire`, `…/cli` — Startup samt Protokollrand, `Query`, `Terminate`, Verbindungsende; `record`, `version` | `server.go`, `errorfields.go`, `cli.go` | konform |
| `internal/adapters/driven/postgres`, `…/recording` | `upstream.go`, `yaml.go` | konform |
| `internal/bootstrap` | `bootstrap.go` | konform |
| `test/integration` — neun Fälle | zehn `TestE2ERecord*`; der zehnte (`FremdeAnfrage`, `PGR-W3003`) kam mit `fea1182` | **kleine Abweichung:** §3 nennt den Fall nicht, §6 schon |
| Runner und `integration.mk` | `tools/test/run-integration-tests.sh`, `harness/mk/integration.mk` | konform; schreibt nichts in den Baum (bestätigt) |
| `Dockerfile`, `.dockerignore`, `build.mk`, `go.mod`, `go.sum` | vorhanden; Stufen `deps` (Netz), `test`/`build`/`integration` mit `--network=none`, `runtime` distroless nonroot | konform |
| `.a-check.yml` (`tech`) | eine Zeile `go.yaml.in/yaml` | konform |
| Abdeckungs-Skripte, `abdeckung.mk`, `docs/user/abdeckung-*.md`, `.d-check.yml` (`trace.coverage`) | vorhanden | konform |
| Unit-Tests | fünf `_test.go` | konform |
| `harness/mk/vorgaben.mk` | vorhanden | konform |
| — | `cmd/pgwire-recorder/main.go` (Signale, Exit-Code) | **nicht in §3.** Sachlich gedeckt durch §1 (Signalbehandlung) |
| — | `spec/spezifikation.md`: LH-FA-05.e (vier Zeilen Protokollrand), `SPEC-041` (Leseregel `null`), `SPEC-045` (neu), `PGR-W3003` | **nicht in §3**, und `SPEC-041` und `SPEC-045` fehlen in „Berührte Spec-Stellen“. §6 nennt die Spec-Festlegung (Risiko 2) |
| — | `docs/user/benutzerhandbuch.md` (`PGR-W3003`), `README.md`, `harness/README.md` §Sensors, [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) bis [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) samt Index | nicht in §3. Folgeänderungen nach Schritt 7 des Workflows oder Entscheidungsartefakte, keine Code-Abweichung |

### §6 Risiken

| Risiko | Stand im Code | Befund |
|---|---|---|
| Anmeldeverfahren des Upstreams | jedes Verfahren → `PGR-E6001` (`TestNichtVermittelbar`) | wie beschrieben; Ausgang „offen bis Closure“ ist zutreffend |
| fremde erste Nachricht | `PGR-W3003`, Exit-Code unverändert (`TestFremdeErsteNachricht`, `TestE2ERecordFremdeAnfrage`); TCP-Probe ohne Bytes nur `Debug` | „entfallen“ ist zutreffend |
| SIGKILL verliert laufende Sessions | Schreiben nach jedem Session-Ende und in `Finish` | wie beschrieben |
| netzlose Quelle der Abhängigkeiten | durch [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) gelöst (Stufe `deps`, sonst `--network=none`) | der Text sagt noch „offen bis Closure“ und „sonst scheitert der Build“. Sachlich ist das Risiko eingetreten und behandelt; bei der Closure sollte das der Ausgang sein |

## 3. ADR-Konformität

| ADR | Entscheidung | Befund |
|---|---|---|
| [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) | Ports & Adapters | **konform.** `make a-check` meldet 0 Befunde. Die Services importieren nur `model` und `ports/driven`. Die Verdrahtung liegt in `internal/bootstrap` (`composition_root`). |
| [ADR-0003](../plan/adr/0003-pgwire-server-ist-driving-adapter.md) | PGWire-Server ist Driving Adapter | **konform.** `internal/adapters/driving/pgwire` ruft `driving.Recorder` auf; `RecordService` erfüllt den Port (Prüfung zur Übersetzungszeit in `bootstrap.go`). |
| [ADR-0004](../plan/adr/0004-postgresql-upstream-ist-driven-adapter.md) | Upstream über Outbound Port | **konform.** `driven.Upstream`/`UpstreamSession`, Implementierung in `internal/adapters/driven/postgres`, eine Upstream-Verbindung je Session (LH-FA-02.a Schritt 3). |
| [ADR-0005](../plan/adr/0005-recording-store-ist-driven-adapter.md) | `RecordingRepository`-Port, YAML-Adapter | **konform.** Port mit `Prepare`/`Write`/`Load`, `recording.YAML` implementiert ihn. |
| [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md) | Bibliothekstypen an der Grenze übersetzt | **konform.** `toResponse` (Upstream) und `toMessage` (Server) übersetzen zwischen `pgproto3` und `model.Response`. Das Modell importiert nur `fmt`. |
| [ADR-0009](../plan/adr/0009-implementierungssprache-go.md) | Go | **konform.** Go 1.27, `CGO_ENABLED=0`. |
| [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) | `pgproto3` nur in den beiden PGWire-Adaptern | **konform.** Nur `server.go` und `upstream.go` importieren es. `make a-check-negativ` ist grün. Der Integrationstest nutzt `pgconn` als Client außerhalb jeder Schicht; das ist Testcode und kein Produkt. |
| [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) | Multistage-Dockerfile, netzlose Stufen, gepinnte Images, Integration in `make gates` | **konform.** `golang`, `distroless` und `postgres` sind per Digest gepinnt. `test-integration` steht in `GATE_CHECKS`. Die in der ADR genannte Folge („der Build-Cache überspringt einen unveränderten Testlauf“) habe ich beobachtet und durch einen ungecachten Lauf ausgeglichen. |
| [ADR-0027](../plan/adr/0027-yaml-bibliothek.md) | `go.yaml.in/yaml/v3`, beide Modulpfade nur im Recording-Adapter | **konform.** Eigene Sonde: `go.yaml.in/yaml/v3` im Postgres-Adapter → `tech-leak … außerhalb internal/adapters/driven/recording`, Exit 1. Mit `gopkg.in/yaml.v3` ebenso. Im Service → `app-impurity`, Exit 1. Grenze: Eine Gegenprobe in `make gates` hält diese Zeile nicht fest (`a-check-negativ` deckt nur `pgproto3` und `crypto/tls` ab). Siehe V-6. |
| [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) | Deklaration je Pfad; die RTM liest nur vollständig belegte Anforderungen; Gate auf Aktualität und Gegenprobe | **mechanisch konform.** `abdeckung-check` und `abdeckung-gegenprobe` stehen in `GATE_CHECKS`, ein `TestE2E*` ohne Deklaration ist `FEHLER`, `LC_ALL=C` ist gesetzt, `trace.coverage` liest nur `abdeckung-vollstaendig.md`. **Inhaltlich** treffen zwei Deklarationen ihr Kriterium nicht (V-1, V-2). Das ist genau die Lücke, die die ADR unter „Negativ“ dem Review zuweist. |

Keine Accepted-ADR wurde im Diff inhaltlich geändert (Hard Rule 3.5). [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) bis [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) sind neu und stehen im Index.

## 4. Spec-Konformität der Aufzeichnung

| Stelle | Zusage | Befund |
|---|---|---|
| `SPEC-001` | `format: pgwire-recorder`, `version: 1`; unbekannte Version → `PGR-E3002`; fehlende oder fremde Kennung und unbekannter `type` → `PGR-E3003` | **konform** (`Unmarshal`, `TestUnmarshalFehler`, U3, U5). Hinweis: `type: extended` gehört laut `SPEC-001` zu Version 1, dieser Leser kennt es nicht und meldet es als beschädigt. Das entspricht dem Wortlaut („den ein Leser nicht kennt“). Siehe V-5. |
| `SPEC-002` | logisches Modell `sessions[].id/startup/interactions[].sequence/request/responses` | **konform.** Zusätzlich gibt es das optionale Feld `server_parameters`. |
| `SPEC-003` | nicht verlustfrei als String darstellbare Bytes eindeutig gekennzeichnet | **konform.** Ungültiges UTF-8 → `base64` (U1 rot). |
| `SPEC-004` | deterministisch; keine rechnerspezifischen Angaben | **konform.** `TestMarshalDeterministisch`; E2E prüft Listen-, Upstream- und Ausgabepfad (E5 rot). `BackendKeyData` wird nicht aufgezeichnet. |
| LH-FA-02.a | `--listen`/`--upstream`/`--output` Pflicht, eine Upstream-Verbindung je Session; Listen-Port → `PGR-E4001`/4; Upstream → `PGR-E4002` | **konform.** Eigene Sonde mit dem Binary: Listen nicht nutzbar → `Netzwerk [PGR-E4001]`, Exit 4. Fehlende Pflichtoption → `PGR-E2001`, Exit 2. |
| LH-FA-02.b | Ende nach `ReadyForQuery` mit oder ohne `Terminate` ist regulär; Abbruch vor `ReadyForQuery` → `PGR-E4003`, die unvollständige Interaktion entfällt | **konform.** `verbindungsende` umfasst `ErrUnexpectedEOF` (F-254 behoben, `TestE2ERecordEndeOhneTerminate`). Kann eine Antwort nicht zugestellt werden, gilt `EndLost` und die letzte Interaktion entfällt (`TestAntwortNichtZugestellt`). Upstream-Abbruch → `PGR-E4003`, die Interaktion wird nicht angehängt. |
| LH-FA-05.c | `SSLRequest` → `N`, danach unverschlüsselt | **konform** (`TestSSLUndGSSMitN`, `TestE2ERecordSSLAbgelehnt`). |
| LH-FA-05.e | 3.0; andere Version → `PGR-E6002`; Sonderanfrage 1234.x → `PGR-E6001`; `CancelRequest` → `PGR-W3001`; Verbindung ohne Nachricht still; Länge < 8 oder > `SPEC-045` → `PGR-W3003` ohne Antwort und ohne Exit-Wirkung | **konform.** `startup` liest den Startcode per `Peek` vor `pgproto3` (F-262 behoben). |
| LH-FA-07.a | vorhandenes `--output` ohne `--force` → `PGR-E2002`/2 vor dem Annehmen; Schreiben nach jedem Session-Ende und am Ende über eine temporäre Datei mit Umbenennen; ein Lauf ohne Session schreibt ein gültiges Recording | **konform.** Eigene Sonde: `Konfiguration [PGR-E2002] … existiert bereits`, Exit 2. `Prepare` läuft vor `Listen`. `TestWriteRechteUndAtomar`; `sessions: []` in den E2E-Fällen. |
| LH-FA-13.b | Verbindungsfehler beenden nur die Verbindung, der erste wird gemerkt; Exit-Code der gemerkten Klasse, sonst 0; ein Schreibfehler am Ende → 3 mit Vorrang | **konform.** `note` merkt den ersten Fehler (`TestErsterFehlerZaehlt`); `Finish`-Fehler → 3 (`TestE2ERecordSchreibfehlerAmEnde`). Hinweis: Eine `FATAL`-Antwort des Servers im Verbindungsaufbau (Pfad `id == 0`, etwa eine fehlende Datenbank) geht unverändert an den Client und zählt nicht als Verbindungsfehler. Für `record` legt die Spezifikation das nicht fest (die Tabelle mit `PGR-E4002` für „`FATAL` beim Aufbau“ gilt für das Einspielen); siehe V-7. |
| `SPEC-034` | Codes `PGR-[EWI][0-9]{4}`, erste Ziffer = Klasse; Kopf `<klasse> [<code>]: <Ursache>`; Warnungen ändern den Exit-Code nicht | **konform.** `model.Error.Error()` liefert den Kopf; Log-Zeilen tragen `code=`. Warnungen laufen nicht über `note`. Ein `Accept`-Fehler trägt jetzt den Rückfall `PGR-E4000` (F-263). |

## 5. Ehrlichkeit der Abdeckung (RTM, F-257)

Die Behebung von F-257 habe ich gegen das Lastenheft gelesen. LH-FA-02/Boundary steht jetzt über `TestE2ERecordVerbindungOhneAnfrage` und `TestRecordSessionOhneAnfrage` und trifft damit das Kriterium. LH-FA-06/Boundary und LH-QA-01 sind entfernt, LH-FA-06 steht ehrlich als „teilweise“. LH-QA-06/Messung trifft die Messmethode. LH-FA-02 und LH-FA-13 halte ich für ehrlich belegt (zu LH-FA-13 siehe V-4). **Zwei Anforderungen weist die RTM als `Tests | ok` aus, obwohl ein Akzeptanzkriterium nicht erfüllt ist:**

| # | Klasse | Befund | Beleg |
|---|---|---|---|
| V-1 | MEDIUM | **LH-FA-07/Negative** („unvollständige oder beschädigte Datei → als ungültiges Recording erkannt“) ist über `TestUnmarshalFehler` deklariert. Der Leser erkennt aber eine **unvollständige** Datei nicht, sofern sie an einer Zeilengrenze abgeschnitten ist. Damit gilt LH-FA-07 zu Unrecht als vollständig. Das atomare Schreiben verhindert, dass der Recorder selbst solche Dateien erzeugt. Das Kriterium spricht aber vom Laden, etwa nach Kopie oder Transport. | Sonde in der Kopie (`Unmarshal`): Abschnitt nach `version: 1` → `err=<nil>` (gültig, 0 Sessions). Abschnitt nach `data_row` ohne `command_complete`/`ready_for_query` → `err=<nil>`. Abschnitt nach `request` ohne `responses` → `err=<nil>`. Nur die leere Datei → `PGR-E3003`. |
| V-2 | MEDIUM | **LH-FA-05** steht als vollständig. Das Lastenheft fasst den Umfang aber als „Verbindungsaufbau (Start-up einschließlich Authentifizierung)“, und Happy lautet „wird aufgezeichnet **beziehungsweise wiedergegeben**“. Plan §1 und §6 nehmen die Authentifizierung aus: Jedes Anmeldeverfahren ergibt `PGR-E6001`, das Risiko ist „offen bis Closure“. Wiedergabe gibt es nicht. Die Deklaration `LH-FA-05/Happy` über `TestE2ERecordSelect1` belegt damit nur die Record-Hälfte ohne Anmeldung. Die RTM sagt mehr, als der Slice selbst behauptet. | `make doc-trace`: `LH-FA-05 … Tests \| ok`; `upstream.go` (`AuthenticationCleartextPassword`, `MD5`, `SASL` → `CodeUnsupported`) |
| V-3 | LOW | **LH-FA-07/Happy** („nach Neustart für ein Replay verwendbar“) ist nur über den Roundtrip im Prozess vertreten. **LH-FA-07/Boundary** belegt nur die Hälfte „keine rechnerspezifischen Angaben“, nicht „Replay funktioniert unverändert“. Kein Test lädt die vom Binary geschriebene Datei mit `Load`. Für einen Record-Slice ist der Ersatz vertretbar. Bliebe V-1 offen, hätte LH-FA-07 nach dem Replay-Slice drei Pfade, aber keinen davon ganz. | `yaml_test.go`, `record_e2e_test.go` |
| V-4 | INFO | **LH-FA-13**: Boundary nennt auch den Windows-Konsolenabbruch; der lässt sich in diesem Lauf nicht belegen. Negative nennt „Betriebs- oder Konfigurationsfehler“. Die Klasse 2 (Konfiguration) belegt nur `TestParseFehler` (Unit, ohne Deklaration), am Prozess belegt sie kein Test. Meine Sonde zeigt das Binary korrekt (`PGR-E2001`/`PGR-E2002` → 2). | Sonde mit dem Binary |

Zur Frage aus F-272, ob `Tests | ok` einen Lauf belegt: Die RTM weist eine **Deklaration** aus, keinen Lauf. Die Tabellen sagen das selbst („kein Lauf-Beleg; dass die Tests grün laufen, sichert `make gates`“). Für diesen Stand habe ich beides gefahren: Die deklarierten Tests laufen grün (§1), und die Tabellen sind aktuell (`make abdeckung-check` Exit 0). Ehrlich ist die RTM damit in ihrer Mechanik, aber noch nicht in ihrem Inhalt (V-1, V-2).

## 6. Weitere Befunde

| # | Klasse | Befund | Ort |
|---|---|---|---|
| V-5 | INFO | Der Leser setzt `KnownFields(true)` und lehnt damit auch Felder ab, die die Spezifikation schon für Version 1 vorsieht (`offset_ms` je Interaktion, `empty_sessions`; `SPEC-002`, `SPEC-041`), ebenso `type: extended`. Gemessen an diesem Slice ist das richtig. Spätere Slices, die diese Felder schreiben, müssen den Leser erweitern, sonst sind ihre Aufzeichnungen für diesen Stand „beschädigt“. Hinweis an den Planner für `slice-walking-skeleton-replay`. | `internal/adapters/driven/recording/yaml.go` · `Unmarshal` |
| V-6 | INFO | Die `tech`-Zeile für `go.yaml.in/yaml` ([ADR-0027](../plan/adr/0027-yaml-bibliothek.md)) wirkt (eigene Sonde), aber keine Gegenprobe in `make gates` hält sie fest. Wird sie gelöscht, bleibt alles grün. `a-check-negativ` deckt nur `pgproto3` und `crypto/tls` ab. | `tools/arch/a-check-negativ.sh`, `.a-check.yml` |
| V-7 | INFO | `FATAL` des Upstreams im Verbindungsaufbau (zum Beispiel eine fehlende Datenbank) wird an den Client durchgereicht und nicht als Verbindungsfehler gezählt; der Lauf endet mit 0. Die Spezifikation regelt das für `record` nicht. Hinweis an den Planner. | `internal/adapters/driving/pgwire/server.go` · `handle` (Pfad `id == 0`) |
| V-8 | LOW | Der Paketkommentar des Integrationstests nennt `docs/user/e2e-abdeckung.md` als Ziel von `make abdeckung`. Das Skript schreibt aber `docs/user/abdeckung-e2e.md` (`AGENTS.md` §3.7, Zusage weiter als Code). | `test/integration/record_e2e_test.go` · Paketkommentar |

## 7. Gegenprobe dieser Datei

`make docs-check` nach dem Anlegen dieser Datei: `d-check: 114 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

## 8. Gesamturteil

**Die DoD-Liefer-Punkte 1 bis 3 sind bestätigt.** Grundlage sind eigene Läufe (`make gates`, ein ungecachter Unit-Testlauf, `make doc-trace`, `make abdeckung-check`) und eine eigene Mutationsprobe in Kopien: 5 von 5 E2E-Mutanten und 6 von 6 Roundtrip-Mutanten sind rot. Der Code hält [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0003](../plan/adr/0003-pgwire-server-ist-driving-adapter.md) bis [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md), [ADR-0009](../plan/adr/0009-implementierungssprache-go.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md), [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) und [ADR-0027](../plan/adr/0027-yaml-bibliothek.md) ein. Die Aufzeichnung entspricht `SPEC-001` bis `SPEC-004`, LH-FA-02.a/.b, LH-FA-05.c/.e, LH-FA-07.a, LH-FA-13.b und `SPEC-034`. Plan und Code stimmen bis auf nicht nachgezogene Plan-Zeilen überein (§2).

**Nicht bestätigt ist die Ehrlichkeit der RTM** im Sinne von [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md): LH-FA-07 (V-1) und LH-FA-05 (V-2) stehen als `Tests | ok`, obwohl je ein Akzeptanzkriterium nicht erfüllt ist. Das ist keine DoD-Verletzung, denn die DoD verlangt die RTM nicht. Es ist aber eine Zusage des Slice-Plans (§3: „`trace.coverage` liest nur die vollständig belegten Anforderungen“), die inhaltlich nicht trägt.

**Übergabe an den Planner:**

1. V-1 und V-2: entweder die Deklarationen zurücknehmen (`LH-FA-07/Negative` nur für „beschädigt“, `LH-FA-05/Happy` erst mit Anmeldung und Replay), oder den Leser um eine Vollständigkeitsprüfung erweitern (zum Beispiel: jede Interaktion endet mit `ready_for_query`, `sessions` ist vorhanden). Danach `make abdeckung`.
2. Die Commits `e29851a`, `bebac0c` und `fea1182` hat kein Review erfasst, darunter die Spec-Änderung LH-FA-05.e/`SPEC-045`. Vor der Closure ist zu entscheiden, ob ein Folge-Review nötig ist.
3. Plan §3 und „Berührte Spec-Stellen“ nachziehen (`cmd/pgwire-recorder/main.go`, `spec/spezifikation.md` mit `SPEC-041` und `SPEC-045`, zehnter E2E-Fall). Beim Risiko 4 aus §6 bei der Closure den Ausgang nach [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) setzen.
4. V-5 bis V-8 als Hinweise für die Folge-Slices (Replay-Leser, Gegenprobe für die YAML-Zeile, `FATAL` beim Aufbau in `record`, veralteter Kommentar).
