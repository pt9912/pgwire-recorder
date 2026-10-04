# Review-Report: slice-walking-skeleton-record, Folge-Review — 2026-10-04

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Das ist ein Folge-Lauf zu [`2026-10-04-review-slice-walking-skeleton-record.md`](2026-10-04-review-slice-walking-skeleton-record.md) (F-230 bis F-253). Die DoD ist nicht Gegenstand dieses Reviews; sie prüft der Verifier.

**Gegenstand:** `git diff b52a8e3 HEAD` ohne `docs/reviews/` (HEAD `0c99a84`): `bef7015` (Korrekturen zu F-230 bis F-253, Abdeckung je Anforderung und Pfad, [ADR-0027](../plan/adr/0027-yaml-bibliothek.md) Proposed, `harness/mk/vorgaben.mk`) und `0c99a84` ([ADR-0027](../plan/adr/0027-yaml-bibliothek.md) Accepted). Der wiederhergestellte Report `2026-10-03-review-einspielen-zeitangaben.md` ist nur auf Byte-Gleichheit mit seinem Stand vor `fb95322` geprüft.

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-04

**Eingangs-Kontext:**

- `slice-walking-skeleton-record` (§1 bis §6)
- [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0003](../plan/adr/0003-pgwire-server-ist-driving-adapter.md), [ADR-0005](../plan/adr/0005-recording-store-ist-driven-adapter.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md), [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md), [ADR-0027](../plan/adr/0027-yaml-bibliothek.md)
- `spec/lastenheft.md` LH-FA-02, LH-FA-05, LH-FA-06, LH-FA-07, LH-FA-12, LH-FA-13, LH-QA-01, LH-QA-06, LH-RB-01 (Akzeptanzkriterien Happy/Boundary/Negative)
- `spec/spezifikation.md` LH-FA-01.a (Kommandoliste), LH-FA-02.b, LH-FA-05.a, LH-FA-05.e, LH-FA-07.a, LH-FA-12.a, LH-FA-13.a, LH-FA-13.b (Prozessende, Vorrang von Exit 3), `SPEC-013` bis `SPEC-028`, `SPEC-033`, `SPEC-034`, `SPEC-041`
- `AGENTS.md` (Hard Rules 3.1 bis 3.8), `harness/README.md`, `harness/conventions.md`
- der Vorlauf-Report (F-230 bis F-253)
- Quelltext von `pgproto3` v5.11.0 (`Backend.Receive`, `translateEOFtoErrUnexpectedEOF`) aus der Stufe `test`

**Ausgeführte Läufe im Repo:** `make build` (grün; baut jetzt die Stufe `runtime`), `make test` (grün), `make a-check` (0 Befunde; Hinweis wie F-248), `make a-check-negativ` (grün), `make abdeckung-check` (grün), `make docs-check` (0 Befunde), `make test-integration` (grün, 6 Tests). `git status` danach sauber; die Image-IDs von `pgwire-recorder:{build,dev,test,integration}` blieben unverändert. `make gates` lief nicht.

**Mutationen, Beobachtungen, Sonden** — nur an Kopien (`git archive HEAD`) im Scratchpad; Testimages dort umbenannt (`pgr-review-fr-*`) und am Ende gelöscht; Integration über das Runner-Skript der Kopie gegen das gepinnte PostgreSQL-Image.

Die neun Mutationen aus F-234, an den heutigen Code angepasst, einzeln:

| # | Mutation | `docker build --target test` | Integration |
|---|---|---|---|
| M1 | `Sequence` konstant 1 | rot | rot |
| M2 | `Write` schreibt direkt per `os.WriteFile(path, …, 0o644)` | **grün** | **grün** |
| M3 | `SSLRequest` mit `S` beantwortet | rot | rot |
| M4 | Warnung des `CancelRequest` mit `PGR-E1000` statt `PGR-W3001` | **grün** | **grün** |
| M5 | `ReadyForQuery` immer `I` | rot | grün |
| M6 | `AuthenticationOk` vor der Fehlerantwort des Aufbaus (Pfad `id == 0`) | **grün** | **grün** |
| M7 | Versionsprüfung nur für Major 1 und 2 | rot | grün |
| M8 | `DataRow`-Werte umgekehrt im PGWire-Adapter | rot | grün |
| M9 | `DataRow`-Werte umgekehrt im Upstream-Adapter | grün | rot |

Mutationen am neuen Code:

| # | Mutation | Unit | Integration |
|---|---|---|---|
| N1 | Adapter meldet nach gescheitertem Senden `EndNormal` statt `EndLost` | **grün** | **grün** |
| N2 | `&& ctx.Err() == nil` entfernt | grün | rot |
| N3 | `SetReadDeadline` beim Ende von ctx entfernt | grün | rot (15 s) |
| N4 | `note` merkt sich den letzten statt den ersten Fehler | **grün** | **grün** |
| N5 | `record` endet immer mit 0 | grün | rot |
| N6 | scheiterndes `Finish` liefert die gemerkte Klasse statt 3 | **grün** | **grün** |
| N7 | unbekannter Startcode meldet `PGR-E6002` statt `PGR-E6001` | **grün** | **grün** |
| N9 | Service ignoriert `EndLost` | rot | grün |
| N10 | Datei mit `0666` statt `0644` | rot | grün |
| N11 | Service setzt `unsupported` nicht | rot | grün |
| N12 | erster Fehler wird nie gemerkt | rot | rot |

*Beobachtungstests* (zusätzliche Dateien in `test/integration` der Kopie): Client schließt nach `SELECT 1;` den Socket ohne `Terminate` → Exit-Code **4**, Log `PGR-E4003 … unexpected EOF` (F-254); dasselbe mit RST → Exit-Code 4, `connection reset by peer` (F-254); nach `SELECT 1;` eine Nachricht mit Typ `z` → keine `ErrorResponse`, Session mit `SELECT 1;` in der Aufzeichnung, Exit-Code 4, Log `PGR-E4003 … unknown message type: z` (F-255); StartupMessage 4.0 → `ErrorResponse` mit `PGR-E6001` „unbekannte Startnachricht“, Exit-Code 6 (F-262); HTTP-Anfrage an den Port → `ErrorResponse` mit `PGR-E6001`, Exit-Code 6 (F-271).

*Sonde YAML-Leser* (Test in der Kopie, Stufe `test`): `{text: a, text: b}` → gelesen als `b` ohne Fehler, ein doppelter Schlüssel auf oberster Ebene → `PGR-E3003` (F-264); `{~: true}`, `{Null: true}`, `{"null": yes}` → NULL.

*Abdeckungs-Skript* (Kopien, Host-`bash`): `--check` unter `LC_ALL` = `C`, `C.UTF-8`, `en_US.UTF-8`, `de_DE.UTF-8` am heutigen Stand grün; mit zwei Tests `TestB` und `Test_a` gleicher Deklaration unterscheidet sich die Reihenfolge zwischen `C` und `de_DE.UTF-8` (F-265). Eine Probedatei zeigte: Deklaration mit Leerzeile vor `func`, Deklaration über einer Hilfsfunktion, eingerückte Deklaration → still verworfen; Test ohne Deklaration → nicht gemeldet; ein zweiter Kommentarabsatz → in die Kurzbeschreibung übernommen; `|` im Text → Tabellenzeile mit einer Spalte zu viel; `LH-FA-99` mit allen drei Pfaden → steht in `abdeckung-vollstaendig.md`, `make docs-check` 0 Befunde, `make doc-trace` ohne Zeile dafür; eine leere Testfunktion mit `LH-FA-05/Happy` → LH-FA-05 „vollständig“, `make doc-trace` zeigt `Tests | ok` (F-257, F-258, F-259). Mutationen am Skript, danach `--check` gegen die versionierten Tabellen: Pfad-Prüfung entfernt → grün; Abbruch bei `FEHLER` entfernt → grün; `--check` schreibt statt zu prüfen → grün; `chmod 0644` entfernt → grün; Vollständigkeit ohne Negative → rot (F-259). `--check` lässt die Tabellen bei veralteter Deklaration unverändert (md5 vorher/nachher gleich). `make e2e-abdeckung` in der Kopie überschreibt `docs/user/e2e-abdeckung.md`, danach `make abdeckung-check` rot (F-267).

---

## Stand der Findings aus dem Vorlauf

| ID | Vorlauf | Stand | Beleg |
|---|---|---|---|
| F-230 | HIGH | **behoben**, Restpfad siehe F-255 | `CloseSession` trägt `model.SessionEnd`; der Service verwirft bei `EndUnsupported` oder einer nicht unterstützten Serverantwort und vergibt dann keine `id`; der Adapter meldet `Parse` und andere bekannte, nicht unterstützte Nachrichten mit `EndUnsupported`. `TestE2ERecordNichtUnterstuetzt` (COPY nach `SELECT 1;`) und `TestRecordSessionNichtUnterstuetzt` halten es; N11 rot. Services importieren nur `model` und `ports/driven`. |
| F-231 | MEDIUM | **behoben**, neue Folge siehe F-254 | `FirstErrorCode` und `exitCode`; der Test erwartet 4; N5 und N12 rot. Ein scheiterndes `Finish` endet mit 3 vor der gemerkten Klasse (LH-FA-13.b), ungetestet (N6, F-260). |
| F-232 | MEDIUM | **behoben** | `SetReadDeadline` beim Ende von ctx, `Query` des Upstreams beachtet ctx nicht, die laufende Interaktion läuft also zu Ende; `close` mit `context.WithoutCancel`. `main` gibt nach dem ersten Signal die Signale frei. `TestE2ERecordBeendenMitOffenerVerbindung` hält es (N3 rot). Das zweite Signal ist nur gelesen, nicht gemessen. Restfenster beim Aufbau siehe F-269. |
| F-233 | MEDIUM | **teilweise** | Startcode vor `pgproto3` gelesen; 2.0, 3.1, 3.2 → `PGR-E6002` (`TestAndereProtokollversion`, M7 rot). Major 0 und ab 4 → `PGR-E6001`, siehe F-262. |
| F-234 | MEDIUM | **teilweise** | Unit-Tests für PGWire- und Upstream-Adapter, Service und Leser, sechs E2E-Fälle. Sechs der neun Mutationen werden gefangen; M2, M4, M6 bleiben grün, dazu N1, N4, N6, N7 (F-260). |
| F-235 | MEDIUM | **behoben**, Nebenwirkung siehe F-264 | `valueDTO.UnmarshalYAML` liest `!!null`-Schlüssel als `null`; `TestUnmarshalNullUngequotet`. |
| F-236 | MEDIUM | **behoben** | Anfrage- und Antwort-Typ geprüft, mehr als eine Wertform ist `PGR-E3003`; drei neue Fälle in `TestUnmarshalFehler`. |
| F-237 | MEDIUM | **behoben** | `tech` führt `go.yaml.in/yaml` und `gopkg.in/yaml`; die Bibliothekswahl trägt [ADR-0027](../plan/adr/0027-yaml-bibliothek.md). |
| F-238 | LOW | **behoben**, Code siehe F-263 | `Serve` loggt vorübergehende `Accept`-Fehler und nimmt weiter an; `net.ErrClosed` beendet. |
| F-239 | LOW | **behoben**, Codes siehe F-263 | TCP-Probe auf `Debug`; Verbindungsverlust mit `PGR-E4003`; scheitert das Senden, endet die Session mit `EndLost`. |
| F-240 | LOW | **behoben**, Einordnung siehe F-263 | `toMessage` liefert einen Fehler statt `I` oder `nil`; `TestToMessageErfindetNichts`, M5 rot. |
| F-241 | LOW | **teilweise** | §1, §6 und `Berührte Spec-Stellen` nachgezogen; §3 nicht vollständig, siehe F-266. |
| F-242 | LOW | **behoben**, Nebenwirkung siehe F-256 | Probedatei im Zielverzeichnis bei `Prepare`; `TestPrepareVorhandeneDatei` mit fehlendem Verzeichnis. |
| F-243 | LOW | **behoben** | Testcontainer benannt, `aufraeumen` entfernt ihn mit PostgreSQL und Netz. |
| F-244 | LOW | **behoben**, neue Kollision siehe F-267 | Der Runner schreibt nichts mehr; die Tabellen schreibt nur das Werkzeug `make abdeckung`, das Gate `make abdeckung-check` liest nur; die Spalte `Ort` trägt keine Zeilennummern. |
| F-245 | LOW | **behoben** | `go vet -tags integration ./...` in der Stufe `test`. |
| F-246 | LOW | **behoben** | `harness/mk/vorgaben.mk` nimmt `docs/reviews` vom Verweis-Nachzug aus; der Report ist byte-gleich mit seinem Stand vor `fb95322`. |
| F-247 | LOW | **behoben** | Kommando `version` (LH-FA-01.a führt es); Plan-Zeile siehe F-266. |
| F-248 | INFO | **offen**, bewusst nicht umgesetzt | `make a-check` meldet weiterhin `test/integration/record_e2e_test.go` außerhalb jeder Schicht. |
| F-249 | INFO | **teilweise** | `make build` baut die Stufe `runtime`; gestartet wird das Produkt-Image nicht, die Fitness Function von [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) prüft weiter kein Gate; Sensors-Zeile siehe F-266. |
| F-250 | INFO | **offen**, bewusst nicht umgesetzt; erweitert, siehe F-270 | — |
| F-251 | INFO | **offen**, bewusst nicht umgesetzt | — |
| F-252 | INFO | **teilweise** | `TestRecordGleichzeitigeSessions` (20 Sessions); ohne `-race`, weil `CGO_ENABLED=0`. |
| F-253 | INFO | **teilweise** | SQLSTATE je Klasse (`0A000`, `08006`, `XX000`); eine Festlegung in der Spezifikation gibt es weiter nicht. Hinweis bleibt beim Planner. |

## Neue Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-254 | MEDIUM | Eine Client-Verbindung, die nach einem `ReadyForQuery` ohne `Terminate` endet, zählt als Verbindungsfehler `PGR-E4003`, und der Lauf endet mit Exit-Code 4 (beobachtet für FIN und RST). `pgproto3` liefert bei sauberem Ende `io.ErrUnexpectedEOF`, nicht `io.EOF`; LH-FA-02.b nennt das Ende „mit oder ohne `Terminate`“ regulär, `SPEC-013` verlangt dann 0. Das trifft auch eine Pool-Probe, die sich anmeldet und den Socket schließt (LH-FA-07.a). Kein Test schließt ohne `Terminate`. | LH-FA-02.b; LH-FA-13.b; `SPEC-013` | `internal/adapters/driving/pgwire/server.go` · `handle` (Zweig nach `be.Receive`) | ja — Integrationstest: `SELECT 1;`, Socket ohne `Terminate` schließen, Exit-Code 0 erwartet | Reguläres Verbindungsende als Fehler gezählt |
| F-255 | MEDIUM | Eine Client-Nachricht, die `pgproto3` nicht lesen kann (unbekannter Typ, nicht dekodierbarer Inhalt), endet ohne `ErrorResponse`; die Session wird mit `EndNormal` übernommen und der Fehler als `PGR-E4003` gezählt (beobachtet mit Typ `z`). LH-FA-05.e verlangt für eine nicht unterstützte Interaktion eine `ErrorResponse`, LH-FA-12.a verwirft die Session. Kein regulärer Client erzeugt solche Nachrichten; die Klasse des Vorlaufs (F-230) bleibt auf diesem Pfad offen. | LH-FA-05.a; LH-FA-05.e; LH-FA-12.a | `internal/adapters/driving/pgwire/server.go` · `handle` (Fehlerzweig von `be.Receive`) | ja — Integrationstest mit einer Nachricht unbekannten Typs nach `SELECT 1;` | Verworfene Session wird übernommen |
| F-256 | MEDIUM | `Write` setzt die Aufzeichnung fest auf `0644`, unabhängig von der umask und auch dann, wenn `--force` eine Datei mit engeren Rechten ersetzt. `SPEC-033` legt die Dateirechte in die Verantwortung des Anwenders, LH-RB-01 nennt sensible Inhalte; F-242 bemängelte das Abweichen von der umask, `TestRoundtrip` sichert jetzt `0644` zu. | `SPEC-033`; LH-RB-01; LH-FA-07.a | `internal/adapters/driven/recording/yaml.go` · `Write`; `internal/adapters/driven/recording/yaml_test.go` · `TestRoundtrip` | ja — Lauf mit `umask 077`, Rechte der Aufzeichnung | Dateirechte fest statt aus der Umgebung |
| F-257 | MEDIUM | Die Pfad-Deklarationen treffen die Akzeptanzkriterien des Lastenhefts nicht durchgehend, und die RTM übernimmt sie: LH-FA-02/Boundary steht über `TestRecordSessionVerbindungsende` (Kriterium: Client ohne Anfrage), LH-FA-07/Boundary über `TestRecordSessionOhneAnfrage` (Kriterium: Aufzeichnung auf anderem Rechner ohne rechnerspezifische Angaben, ohne Test), LH-FA-06/Boundary über mehrere Interaktionen (Kriterium: mehrere Ergebnismengen oder keine Zeilen). LH-QA-01 und LH-QA-06 haben keine Happy/Boundary/Negative-Kriterien, der Skript-Kopf nennt die Pfade dennoch „die drei Akzeptanzkriterien des Lastenhefts“. `make doc-trace` zeigt LH-FA-02 und LH-FA-07 als `Tests | ok`. | LH-FA-02, LH-FA-06, LH-FA-07, LH-QA-01, LH-QA-06; Maintainability | `internal/hexagon/services/record_test.go` · Deklarationen über `TestRecordSessionVerbindungsende`, `TestRecordSessionOhneAnfrage`; `test/integration/record_e2e_test.go` · über `TestE2ERecordMehrereInteraktionen`; `tools/test/abdeckung.sh` · Kopfkommentar; `docs/user/abdeckung-vollstaendig.md` | ja (Lesen gegen `spec/lastenheft.md`) | Abdeckungs-Deklaration trifft das Akzeptanzkriterium nicht |
| F-258 | MEDIUM | Die Prüfung „jeder `TestE2E*` trägt eine Abdeckungs-Zeile“ (bisher Exit 1 von `make test-integration`) ist entfallen; das neue Skript meldet einen Test ohne Deklaration nicht. Es kamen Formprüfungen hinzu, die entfallene Regel wurde nicht übernommen. Ob das eine Lockerung nach Hard Rule 3.6 ist, entscheidet der Architect; zweites Auftreten der Klasse von F-237. | `AGENTS.md` §3.6 | `tools/test/run-integration-tests.sh` (entfernter Abschnitt); `tools/test/abdeckung.sh` | ja — Probedatei mit Test ohne Deklaration, `make abdeckung-check` grün | Gate-Regel ersetzt statt ergänzt |
| F-259 | MEDIUM | Das neue Gate-Skript hat keine Gegenprobe, und Fehlformen fallen still durch: eine Deklaration mit Leerzeile vor `func`, über einer Hilfsfunktion oder eingerückt wird verworfen; ein Folgeabsatz geht in den Text; `\|` sprengt die Tabelle; erfundene Kennungen (`LH-FA-99`) stehen in `abdeckung-vollstaendig.md` ohne Befund von `make docs-check`. Vier Mutationen am Skript (Pfad-Prüfung, Abbruch bei `FEHLER`, `--check` schreibt, `chmod`) bleiben unter `make abdeckung-check` grün. Die übrigen Gates führen eine Gegenprobe (`make a-check-negativ`, `make hook-gegenprobe`). | Reviewer-Skill MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag“; Maintainability | `tools/test/abdeckung.sh` · awk-Teil, Schritt 3; `harness/mk/abdeckung.mk` | ja — die Mutationen und die Probedatei | Gate ohne Gegenprobe |
| F-260 | MEDIUM | Nach den Korrekturen sind diese Zusagen ohne fangenden Test: nicht atomares Schreiben (M2), Code der `CancelRequest`-Warnung (M4), `AuthenticationOk` vor der Fehlerantwort des Aufbaus (M6), `EndLost` im Adapter (N1), erster statt letzter Fehler (N4), Vorrang von Exit 3 bei scheiterndem `Finish` (N6), Code für unbekannte Startnachrichten (N7) sowie das Ende ohne `Terminate` (F-254). | LH-FA-07.a; LH-FA-05.e; LH-FA-02.b; LH-FA-13.b; Reviewer-Skill MEDIUM „fehlende Negativtests“ | `internal/adapters/driving/pgwire/server_test.go`; `internal/adapters/driven/recording/yaml_test.go`; `test/integration/record_e2e_test.go` | ja — die Mutationen M2, M4, M6, N1, N4, N6, N7 | Protokollrand ohne Negativtest |
| F-261 | LOW | `TestRecordGleichzeitigeSessions` und LH-FA-12/Boundary zeigen gleichzeitige Sessions als belegt, Plan §1 nimmt „Zuordnung und Kennungen paralleler Sessions nach LH-FA-12.a“ aus. Plan und Deklaration sagen über denselben Umfang Verschiedenes. | Maintainability; Plan §1 | Slice-Plan §1; `internal/hexagon/services/record_test.go` · Deklaration über `TestRecordGleichzeitigeSessions` | ja (Lesen) | Plan und Abdeckung widersprechen sich |
| F-262 | LOW | Eine StartupMessage mit Major 0 oder ab 4 (beobachtet: 4.0) erhält `PGR-E6001` „unbekannte Startnachricht“ statt `PGR-E6002`. LH-FA-05.e verlangt `PGR-E6002` für jede andere Protokollversion; der Kommentar von `startup` sagt „eine andere Protokollversion als 3.0 ist PGR-E6002“ zu. | LH-FA-05.e; `AGENTS.md` §3.7 | `internal/adapters/driving/pgwire/server.go` · `startup` (Kommentar und Fallunterscheidung nach dem Startcode) | ja — rohe StartupMessage 4.0 | Zusage weiter als Code |
| F-263 | LOW | Mehrere Codes passen nicht zur Ursache: ein `Accept`-Fehler trägt `PGR-E4001` („Listen-Port nicht zu öffnen“), `SPEC-034` sieht für eine Ursache ohne eigenen Code den Rückfall der Klasse vor; ein Abbildungsfehler in `toMessage` wird als `PGR-E4003` mit `EndLost` gezählt; ein gescheitertes Zwischenschreiben (`PGR-E3001`) wird als „Verbindungsfehler“ geloggt und als erster Verbindungsfehler gemerkt; im Pfad `id == 0` steht `PGR-E4003` an einer Warnung und zählt nicht. | `SPEC-034`; LH-FA-13.b; LH-FA-14.a | `internal/adapters/driving/pgwire/server.go` · `Serve`, `send`, `close`, `note`, `handle` (Pfad `id == 0`) | ja (Lesen) | Meldungscode passt nicht zur Ursache |
| F-264 | LOW | `valueDTO.UnmarshalYAML` umgeht die Erkennung doppelter Schlüssel: `{text: a, text: b}` wird als `b` gelesen, während ein doppelter Schlüssel an anderer Stelle `PGR-E3003` ist. | `SPEC-001`; LH-FA-07 Negative | `internal/adapters/driven/recording/yaml.go` · `valueDTO.UnmarshalYAML` | ja — `Unmarshal` mit doppeltem Schlüssel in einem Wert | Leser nimmt widersprüchliche Werte an |
| F-265 | LOW | Das Gate hängt an der Locale des Hosts: `sort` und `find \| sort` laufen ohne feste Locale, zwei Tests gleicher Anforderung und gleichen Pfads (`TestB`, `Test_a`) erscheinen unter `C` und `de_DE.UTF-8` in verschiedener Reihenfolge. Derselbe Baum ist dann je nach Rechner grün oder rot. | Maintainability; LH-QA-03 | `tools/test/abdeckung.sh` · `sort`-Aufrufe | ja — `make abdeckung-check` unter zwei Locales mit der Probedatei | Ausgabe hängt an der Umgebung des Hosts |
| F-266 | LOW | Dokumente sind dem Diff nicht gefolgt: Plan §3 nennt für `…/cli` nur `record` (neu: `version`), für `test/integration` nur zwei Fälle (jetzt sechs) und für den Runner „schreibt die E2E-Abdeckungstabelle“ (er schreibt nichts mehr); die Sensors-Zeile von `make build` nennt die Stufe `build`, das Ziel baut `runtime`. Drittes Auftreten der Klasse (F-215, F-241) — Pflege-Schwelle des Skills. | Maintainability; `v6.13.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice | Slice-Plan §3; `harness/README.md` §Sensors · Zeile `make build` | ja (Lesen gegen `make help` und Diff) | Plan-Zeile nach Korrektur nicht nachgezogen |
| F-267 | LOW | `make e2e-abdeckung` (Bootstrap-Werkzeug, in `make help`) schreibt `docs/user/e2e-abdeckung.md` aus `tools/harness/selbstpruefung.sh`; danach ist `make abdeckung-check` rot. Der Kopf von `abdeckung.mk` sagt zu, die Tabellen würden „ausschliesslich“ über `make abdeckung` geschrieben. | Maintainability; `AGENTS.md` §3.7 | `harness/mk/abdeckung.mk` · Kopfkommentar; `harness/mk/e2e-abdeckung.mk` · `E2E_ABDECKUNG_ZIEL` | ja — `make e2e-abdeckung`, dann `make abdeckung-check` | Zwei Erzeuger schreiben dieselbe Datei |
| F-268 | LOW | Die Sensors-Zeile bindet `make abdeckung-check` an [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md); die ADR entscheidet über Abdeckung, Deklarationsform und Vollständigkeitsregel nichts. Die Regel „nur vollständig belegte Anforderungen gehen in `trace.coverage`“ steht nur in Skript, Plan und `.d-check.yml`. | Maintainability; `harness/README.md` §Sensors (Spalte Bindung) | `harness/README.md` §Sensors · Zeile `make abdeckung-check` | ja (Lesen) | Bindung zeigt auf eine ADR, die es nicht entscheidet |
| F-269 | LOW | Endet ctx, nachdem die StartupMessage gelesen ist, aber bevor `OpenSession` wählt, scheitert `DialContext` sofort; der Client erhält `PGR-E4002`, und der sonst fehlerfreie Lauf endet mit 4. Der Kommentar von `Run` sagt zu, jede Verbindung ende „nach ihrer laufenden Interaktion“. Nur gelesen, nicht gemessen. | LH-FA-13.a; LH-FA-13.b; `AGENTS.md` §3.7 | `internal/adapters/driving/pgwire/server.go` · `handle` (`OpenSession` mit ctx); `internal/adapters/driven/postgres/upstream.go` · `Open` | nein (Zeitfenster) | Herunterfahren bricht Verbindungsaufbau als Fehler ab |
| F-270 | INFO | `make abdeckung-check` ist ein Gate, das ganz auf dem Host läuft (`bash`, `awk`, `sort -V`, `paste`, `sed`, `cmp`, `find`, `mktemp`); die Sensors-Zeile nennt das nicht, anders als bei `make hook-gegenprobe` („bash, ohne Docker“). Erweitert F-250; Hinweis an den Architect. | `AGENTS.md` §3.1 | `harness/mk/abdeckung.mk`; `tools/test/abdeckung.sh` | ja (Lesen) | Host-Abhängigkeit über die deklarierte Menge hinaus |
| F-271 | INFO | Jede Verbindung, deren erste acht Bytes kein bekannter Startcode sind (HTTP-Gesundheitsprobe, Port-Scan), zählt als `PGR-E6001` und färbt den Lauf auf Exit-Code 6. Ob eine fremde Protokollanfrage eine „nicht unterstützte Interaktion“ nach LH-FA-13.b ist, legt die Spezifikation nicht fest. Hinweis an den Planner. | LH-FA-05.e; LH-FA-13.b | `internal/adapters/driving/pgwire/server.go` · `startup` (Zweig `default`) | ja — HTTP-Anfrage an den Port | Meldungscode ohne Spec-Zuordnung gewählt |
| F-272 | INFO | [ADR-0027](../plan/adr/0027-yaml-bibliothek.md) hält unter Konsequenzen fest, Leser und Beispiele müssten `null` gequotet und ungequotet tragen. `SPEC-041` zeigt nur `null: true`; die Leseregel für beide Formen steht nur in der ADR. Als Folge der Wahl ist der Satz mit Hard Rule 3.8 verträglich; Hinweis an Architect und Planner. Dass die RTM mit `Tests \| ok` eine Deklaration und keinen Lauf ausweist, bewertet der Verifier. | `AGENTS.md` §3.8; `SPEC-041`; Verweis Verifier | `docs/plan/adr/0027-yaml-bibliothek.md` · Konsequenzen; `docs/user/abdeckung-vollstaendig.md` | ja (Lesen) | Norm nur in der ADR |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/model` — Core-Reinheit | geprüft, ohne Befund. `SessionEnd` ist ein Model-Typ ohne Import; die Kommentare nennen Folge und Kriterium. |
| `internal/hexagon/ports/driving` | geprüft, ohne Befund. `CloseSession` trägt `model.SessionEnd`; Signaturen nur mit Model-Typen und `context`. |
| `internal/hexagon/services` — Schichtregeln | geprüft, ohne Befund. Importiert `errors`, `sync`, `model`, `ports/driven`, keinen Driving Port; die Erfüllung von `Recorder` hält weiter die Composition Root fest. `EndLost` kürzt die letzte Interaktion, `EndUnsupported` oder `unsupported` verwerfen ohne `id`; `id` lückenlos in der Reihenfolge der Sessionenden (LH-FA-12.a). Der Session-Zustand wird nur von der eigenen Goroutine berührt. |
| `internal/adapters/driving/pgwire` | geprüft. Startcode per `Peek` vor `pgproto3`; SSL/GSS → `N` und erneute Prüfung; `CancelRequest` → `PGR-W3001`; SQLSTATE nach Klasse; `fail` merkt den ersten Fehler; `close` mit `context.WithoutCancel`; Lesefrist nur beim Ende von ctx, Schreiben unberührt; `errorfields.go` ist reine Verschiebung samt `strings.Contains`. Befund F-254, F-255, F-262, F-263, F-269, F-271. |
| `internal/adapters/driven/postgres` | geprüft, ohne Befund im Diff (nur Tests). Der Test-Server belegt `BackendKeyData` außerhalb der Antworten (`SPEC-004`), Anmeldeverlangen, Fehlerantwort im Aufbau und `CopyOutResponse`. |
| `internal/adapters/driven/recording` | geprüft. Probedatei mit `CreateTemp` und Entfernen, Fehler als `PGR-E3001`; Typprüfung, Wertform, `!!null`. Befund F-256, F-264. |
| `internal/adapters/driving/cli`, `internal/bootstrap`, `cmd/pgwire-recorder` | geprüft. `version` ohne Argumente, sonst `PGR-E2001`; Exit-Code aus dem ersten Verbindungsfehler, `Finish`-Fehler mit 3 vorrangig; zweites Signal nach `stop()` mit dem Standardverhalten (Lesart der Standardbibliothek). Befund F-269. |
| `test/integration`, Unit-Tests | geprüft. Sechs E2E-Fälle, Prüfung der Reihenfolge in der Aufzeichnung, Exit-Codes 0/4/6. Befund F-257, F-260, F-261. |
| `tools/test/abdeckung.sh` | geprüft. `set -euo pipefail`, Arbeitsverzeichnis per `mktemp -d` mit `trap`; Folgezeilen werden angehängt, Kennungen über Zeilen hinweg geteilt; `--check` schreibt nicht (md5 gleich), Schreiben nur bei Abweichung mit `0644`; Vollständigkeit verlangt alle drei Pfade, Nachweisarten je Pfad zusammengeführt; `trace.coverage` liest nur die vollständige Liste. Befund F-257, F-258, F-259, F-265, F-267, F-270. |
| `harness/mk/abdeckung.mk`, `harness/mk/integration.mk`, `tools/test/run-integration-tests.sh` | geprüft. Kein Gate schreibt mehr eine Datei, die ein anderes liest; Runner benennt den Testcontainer und räumt ihn auf. Befund F-267, F-268. |
| `harness/mk/vorgaben.mk` | geprüft, ohne Befund. `=` gewinnt über das `?=` in `slice-mv.mk` unabhängig von der Include-Reihenfolge; die Vorgabe-Pfade bleiben enthalten. |
| `Dockerfile`, `harness/mk/build.mk` | geprüft. `vet -tags integration`, Stufe `runtime` gebaut; Befund F-266 (Sensors-Zeile). |
| `.a-check.yml`, `.d-check.yml` | geprüft, ohne Befund. Beide YAML-Modulpfade auf den Recording-Adapter beschränkt ([ADR-0027](../plan/adr/0027-yaml-bibliothek.md) Fitness Function); `trace.coverage` liest `abdeckung-vollstaendig.md`. |
| Hard Rule 3.1 — Docker-only | geprüft. Build, Test, Integration, a-check, d-check in Docker; Befund F-270. |
| Hard Rule 3.2 — Suppression | geprüft, ohne Befund. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. Im Bereich kein Move. |
| Hard Rule 3.5 — ADR-Immutabilität | geprüft, ohne Befund. `0c99a84` ändert an [ADR-0027](../plan/adr/0027-yaml-bibliothek.md) nur Status und Geschichte; [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) unverändert. |
| Hard Rule 3.6 — Gates nicht lockern | geprüft. Ein Gate hinzugekommen, `tech` ergänzt; Befund F-258. |
| Hard Rule 3.7 — Kommentare | geprüft. Indikativ, Zusage/Kopplung/Grenze; keine verworfene Alternative, kein abwesender Text, kein abgebrochener Satz. Zusagen weiter als Code siehe F-262, F-267, F-269. |
| Hard Rule 3.8 — ADR ohne Regeldetails | geprüft. [ADR-0027](../plan/adr/0027-yaml-bibliothek.md) nennt Entscheidung, vier Alternativen, Gründe, Fitness Function; keine Fallunterscheidung der Spezifikation nacherzählt. Hinweis F-272. |
| ADR-Index, `harness/README.md` §Sensors, `README.md` | geprüft. Index-Zeile deckt sich mit dem Kopf von [ADR-0027](../plan/adr/0027-yaml-bibliothek.md); `abdeckung`, `abdeckung-check` existieren in `make help`. Befund F-266, F-268. |
| Commit-Messages `bef7015`, `0c99a84` | geprüft, ohne Befund. Beide tragen `slice-walking-skeleton-record` und `LH-*`; `bef7015` nennt [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) im Betreff und [ADR-0027](../plan/adr/0027-yaml-bibliothek.md) im Text; keine Struktur-IDs. |
| `spec/` | geprüft, ohne Befund. Nicht im Diff. |

## Summary

| Kategorie | Anzahl (neu) |
|---|---|
| HIGH | 0 |
| MEDIUM | 7 |
| LOW | 9 |
| INFO | 3 |

Vorlauf: 15 behoben (davon neun mit Verweis auf ein neues Finding), 6 teilweise, 3 offen und bewusst nicht umgesetzt (F-248, F-250, F-251).

**Finding-Klassen dieses Laufs:** Reguläres Verbindungsende als Fehler gezählt · Verworfene Session wird übernommen (2. Auftreten, vgl. F-230) · Dateirechte fest statt aus der Umgebung · Abdeckungs-Deklaration trifft das Akzeptanzkriterium nicht · Gate-Regel ersetzt statt ergänzt (2. Auftreten, vgl. F-237) · Gate ohne Gegenprobe · Protokollrand ohne Negativtest (2. Auftreten, vgl. F-234) · Plan und Abdeckung widersprechen sich · Zusage weiter als Code · Meldungscode passt nicht zur Ursache · Leser nimmt widersprüchliche Werte an · Ausgabe hängt an der Umgebung des Hosts · Plan-Zeile nach Korrektur nicht nachgezogen (3. Auftreten, vgl. F-215, F-241 — Pflege-Schwelle) · Zwei Erzeuger schreiben dieselbe Datei · Bindung zeigt auf eine ADR, die es nicht entscheidet · Herunterfahren bricht Verbindungsaufbau als Fehler ab · Host-Abhängigkeit über die deklarierte Menge hinaus (2. Auftreten, vgl. F-250) · Meldungscode ohne Spec-Zuordnung gewählt (2. Auftreten, vgl. F-253) · Norm nur in der ADR

## Verdikt

**Merge-blockierend:** nein — kein HIGH. Der Kern von F-230 ist behoben: Port, Service und Adapter tragen das Session-Ende mit Grund, und eine Session mit nicht unterstützter Interaktion erhält keine `id`. F-255 lässt einen Restpfad offen, den nur Nachrichten außerhalb des Protokolls erreichen.

Vor der Closure verdienen zwei MEDIUM Aufmerksamkeit. F-254 kehrt F-231 um: Jetzt endet ein Lauf mit 4, wenn ein Client die Verbindung nach einem `ReadyForQuery` ohne `Terminate` schließt, gegen den ausdrücklichen Wortlaut von LH-FA-02.b. F-257 betrifft die Abdeckung: Die RTM weist LH-FA-02 und LH-FA-07 als vollständig belegt aus, gestützt auf Deklarationen, die andere Kriterien prüfen als das Lastenheft nennt.

**Übergabe:** Findings an den Implementer. F-258 an den Architect (Hard Rule 3.6, Klasse von F-237), dazu F-268, F-270 und F-272. F-261, F-266 und F-271 an den Planner. F-257 und F-272 (RTM als Deklaration) zusätzlich an den Verifier. F-266 erreicht die Pflege-Schwelle des Skills (drittes Auftreten): Steering-Loop-Frage, ob ein Gate Plan §3 gegen die berührten Pfade des Diffs hält. Widerspricht der Implementer einer MEDIUM-Einstufung, genügt hier Annahme oder Begründung; der Konflikt-Pfad über den Architect gilt ab HIGH oder ab dem dritten gleichen Konflikttyp. Die Finding-Klassen gehen in die Slice-Closure §7. Die DoD-Konformität prüft der Verifier separat.
