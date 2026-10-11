# Verifikation: slice-v1-abschluss-sqlite-format — 2026-10-11

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-11-review-slice-v1-abschluss-sqlite-format.md`, F-614 bis F-622); die Übergaben prüfe ich nach. Ich repariere nichts.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-sqlite-format.md` (Schreib-Teil A), Kopf und §1 bis §8, gegen den Gesamt-Diff `4914a4a..a4dfd15` (35 Dateien, +3257 −43). Architect: `75874ba`, `7a5a14e`. Implementer: `f052e63`, `fc59441`, `07d34c4`, `9ebab6c`, `10c10ef`, `f26ebb9`; Nacharbeit `e4ada22`, `a4dfd15`. Review: `075e63a`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingang:**

- der Plan ganz am Stand `a4dfd15`, besonders §1 (*Abgegeben*, *Ausdrücklich NICHT*), §2 (DoD), §3 (mit den Aufträgen des Architect), §6 (*Randformen*), §7 (*Belege des Implementers*, beide Mutationstabellen, Namensabgleich), §8
- `spec/spezifikation.md` `LH-FA-22.a` (*Wahl und Erkennung*, *Schreiben von `sqlite`*), `SPEC-043`, `LH-FA-07.a`
- der Nehmer `slice-v1-abschluss-sqlite-format-lesen` (§1, DoD, §3, §6) und `slice-v1-abschluss-sessions` (für `--record-empty-sessions`)
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen
- am Stand `a4dfd15`: `sqlite.go`, `sqlite_zeilen.go`, `sqlite_json.go`, `format.go`, `schema.sql`, die Diffs von `cli.go`, `bootstrap.go`, Port und `RecordService`, `.golangci.yml`, `tools/schema/schema.yaml`, `sqlite_format_e2e_test.go`
- `AGENTS.md` §3.9 bis §3.13

Beim Start war der Arbeitsbaum sauber, HEAD `a4dfd15`. Der Commit dieses Berichts enthält nur diese Datei.

**So wurden die Proben gebaut.** Alles Docker-only, Scratch-Präfix `ver-sq-`, nie im Repo.

- **Binary:** `git archive HEAD` in eine frische Kopie, daraus die Stufe `runtime` des `Dockerfile` (distroless) mit eigenem Tag. Gegenstelle: ein echter PostgreSQL-17-Server (das gepinnte Image aus `harness/mk/integration.mk`, `trust`) in einem eigenen Netz; Clients: `psql` desselben Images. Gelesen habe ich die Datei mit dem `sqlite3`-Modul von Python 3.12 (Image `python:3.12-slim`, `mode=ro`) mit rohem SQL, unabhängig vom Lesecode des Adapters.
- **Unit-Mutanten:** je Mutant eine frische Kopie (`cp -r` ohne `-p` von `go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`), genau eine Ersetzung (das Skript bricht ab, wenn der Suchtext nicht genau einmal trifft), im Image der Stufe `deps`, `--network none`, `go test -count=1` über `./internal/adapters/driven/recording`, `./internal/adapters/driving/cli` und `./internal/bootstrap`. Der Grundlauf ohne Mutation war grün.
- **Integrations-Mutant und Gate-Mutanten:** je eine Kopie des Baums; `make test-integration`, `make a-check-negativ` und `make lint` liefen darin.
- **Danach:** Container, Netz, Images und Kopien entfernt (siehe Ende).

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `record --format sqlite` schreibt je Session in einer Transaktion; nur vollständige Sessions, atomar angelegt, nie als Ganzes neu geschrieben; Journal-Rest `PGR-E3001`; Wertemenge `yaml`/`sqlite`; `replay`/`play` kennen `--format` nicht. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Ende-zu-Ende über das gebaute Binary gegen einen echten Upstream | `record --format sqlite` gegen PostgreSQL 17, sieben Sessions (einfache Anfrage, `\bind` als Extended Query, `client_encoding=LATIN1`, Fehlerantworten, Datenzeile mit dem Byte `0xE9`). Exit 0 nach `SIGINT`; die Datei beginnt mit `SQLite format 3` und NUL (`xxd`) | bestätigt |
| Inhalt gegen die Liste in `SPEC-043`, rohes SQL, unabhängig gelesen | `meta`: `format=pgwire-recorder`, `version=1`, keine `empty_sessions`-Zeile; `session.id` 1 bis 7 lückenlos; `startup` und `server_parameters` als JSON-Objekte; `interaction.type` `query`/`extended`, `sql` bei `query` gesetzt und roh (`select 'caf\xe9' as c;` mit dem Byte `0xE9` unverändert), bei `extended` `NULL`, `offset_ms` `NULL`; `message`: bei `query` nur `server` mit `group_no` `NULL`, bei `extended` Client-Nachrichten `parse`, `bind`, `describe`, `execute`, `sync` mit `group_no` 1 vor den Server-Nachrichten, `position` 1 bis 11 über beide Richtungen; `fields` in jeder Zeile ein Objekt, ohne `type` (`{}` bei `sync`), kein `bytes`. `PRAGMA user_version` 0, `journal_mode` `delete`, `integrity_check` ok, `foreign_key_check` leer. Das Schema der Datei entspricht `SPEC-043` | bestätigt |
| `base64:`-Escape, Schreib-Richtung | Ein Parameter `base64:abc` steht in `bind` und `data_row` als `\base64:abc`, ein Wert `\x` als `\\x`; ein Byte `0xE9` in der Datenzeile steht als YAML-Darstellung `{"base64":"Y2Fm6Q=="}` (das ist die Darstellung des Modells, nicht der Escape). Ein Startup-Wert mit `0xE9` entstand über das Binary nicht (Docker ersetzt das Byte in der Umgebung durch U+FFFD); dafür gelten die Unit-Tests (siehe Mutanten) | bestätigt, Rest über Mutanten |
| Journal-Rest `PGR-E3001`, auch mit `--force` und ohne Zieldatei | Datei `j.sqlite-journal` allein: Exit 3, `PGR-E3001`, Pfad des Journals in der Meldung; mit `--force`: ebenso. Zieldatei vorhanden ohne `--force`: Exit 2, `PGR-E2002`, auch zusammen mit Journal (Reihenfolge Zielpfad vor Journal) | bestätigt |
| `--force` ersetzt, ergänzt nie; Zielpfad bleibt bis zur ersten Session unberührt | Datei mit sieben Sessions, `record --force`: vor der ersten Session unverändert, danach genau eine Session (`select 'neu';`) | bestätigt |
| Pfad wörtlich | Ziel `we ird/a?b#c%25.sqlite`: genau diese Datei entsteht, nichts anderes im Verzeichnis, Inhalt lesbar | bestätigt |
| Wertemenge und Quellen | `--format SQLite`, `--format ""`: `PGR-E2001`, Exit 2, „erlaubt sind yaml und sqlite“; `PGWIRE_RECORDER_FORMAT=bogus`: `PGR-E2001` mit Namen der Variable; `replay --format` und `play --format`: `PGR-E2001` „flag provided but not defined“ | bestätigt |
| Zwischenstand aus §6 | `replay` und `play` auf der Datei: `PGR-E3003` „Aufzeichnung beschädigt“, Exit 3, wie §6 es als wellenintern akzeptiert beschreibt (Ausgang *eingetreten*, `slice-v1-abschluss-sqlite-format-lesen`) | bestätigt |

**Bewusstes Brechen** (Modul 11). Jede Mutation in einer frischen Kopie, genau eine Ersetzung, `gofmt` sauber, Grundlauf grün. Alle rot aus dem benannten Grund:

| Zusage | Mutation | roter Test |
|---|---|---|
| eine Transaktion je Session, `Commit` | `Commit` durch `Rollback` | `TestSQLiteEinstellungen`, `TestSQLiteSonderzeichenImPfad`, `TestSQLiteWerteAlsParameter` und weitere (Datei bleibt leer) |
| nie als Ganzes neu geschrieben | Datei bei jedem Aufruf neu angelegt (`if !s.angelegt` zu `if true`) | `TestSQLiteStandJeAdapter`, `TestSQLiteUebrigGebliebeneDateien`, `TestSQLiteRechte` und weitere |
| `<Zielpfad>-journal` | Name `.journal` | `TestSQLitePrepareJournal`, `TestSQLitePrepareJournalReihenfolge` |
| Verknüpfung nach ihrem Ziel | `Lstat` statt `Stat` | `TestSQLitePrepareJournalReihenfolge` |
| Nummern vorab (`PGR-E1000`) | Prüfung entfernt | `TestSQLiteNummernNeueDatei`, `TestSQLiteNummernBestehendeDatei` |
| … über alle Sessions, nicht nur die offenen | Prüfung nur für Sessions ab `s.ergaenzt` | `TestSQLiteNummernBestehendeDatei` allein |
| … vor Anlegen und Öffnen | Prüfung hinter `legeAn` verschoben (Umstellung in zwei Ersetzungen) | `TestSQLiteNummernNeueDatei` allein (die Datei entsteht) |
| weniger Sessions als die Datei: `PGR-E1000` | Prüfung entfernt | `TestSQLiteWenigerSessionsAlsDieDatei` |
| Datei wird vor jedem Ergänzen geprüft | Prüfung der Nummern (`id != n`) entfernt; Prüfung der Zahl (`n != erwartet`) entfernt | `TestSQLiteDateiGeprueft` (beide einzeln) |
| `mode=rw`: fehlt die Datei, wird sie nicht angelegt | `mode=rw` entfernt | `TestSQLiteDateiGeprueft` |
| jeder Fehler beim Öffnen ist `PGR-E3001` | Code `PGR-E1000` | `TestSQLiteDateiGeprueft` |
| Rollback-Journal statt WAL, auch bei fremder Datei im WAL-Modus (F-614) | `journal_mode(DELETE)` aus der URI | `TestSQLiteFremdeWALDatei` allein |
| `foreign_keys=ON` | Pragma entfernt | `TestSQLiteEinstellungen`, `TestSQLiteFremdschluessel` |
| `BEGIN IMMEDIATE` | `_txlock=immediate` entfernt | `TestSQLiteTransaktionImmediate` |
| Sperrfrist fünf Sekunden | Konstante 1 | `TestSQLiteEinstellungen`, `TestSQLiteSperrfrist` |
| Abbruchsignal unterbricht keine Transaktion | `ctx` statt `context.WithoutCancel(ctx)` | `TestSQLiteAbbruchUnterbrichtNicht` |
| Schließen nach der Transaktion: `PGR-E3001` | Fehler des Schließens verworfen | `TestSQLiteSchliessenScheitert` |
| Zähler der ergänzten Sessions | Zähler nicht erhöht | `TestSQLiteEinstellungen`, `TestSQLiteSchliessenScheitert` und weitere |
| `angelegt` erst nach Erfolg | `angelegt` nie gesetzt | dieselbe Gruppe |
| `base64:`: gültiger String mit `\` oder `base64:` bekommt ein weiteres `\` | nur das Fluchtzeichen nicht geprüft | `TestSQLiteFluchtzeichen`, `TestSQLiteUmschreibungKehrtZurueck` |
| Bytes ohne UTF-8 stehen als `base64:` | Präfix beim Schreiben weggelassen | `TestSQLiteStringOhneUTF8`, `TestSQLiteUmschreibungKehrtZurueck` |
| `<`, `>`, `&` nicht maskiert | `SetEscapeHTML(true)` | `TestSQLiteJSONTypen` |
| `interaction.sql` unverändert | Text durch `maskiere` geschickt | `TestSQLiteSQLUnveraendert`, `TestSQLiteFluchtzeichen` |
| `fields` ohne `type` | `type` bleibt | `TestSQLiteInhalt`, `TestSQLiteFieldsGleichYAML`, `TestSQLiteLeereWerte` und weitere |
| `group_no` ab 1 | 0-basiert | `TestSQLiteInhalt` |
| `position` über beide Richtungen | Client-Zähler um 1 verschoben | `TestSQLiteInhalt`, `TestSQLiteFieldsGleichYAML` |
| `sql` nur bei `query` | `sql` bei `query` nicht gesetzt | `TestSQLiteInhalt`, `TestSQLiteSQLUnveraendert` und weitere |
| `empty_sessions` nur mit leeren Sessions | Zeile nie geschrieben | `TestSQLiteEmptySessions` |
| Werte als Parameter | Startup, Feld oder SQL in den Anweisungstext (die Mutanten Z16 bis Z18 des Implementers sind nicht von mir gefahren) | `TestSQLiteWerteAlsParameter` (von Z5 oben mitgefangen) |
| `Load` Platzhalter `PGR-E1000` | liefert `nil` | `TestSQLiteLoadPlatzhalter` |
| Fabrik und Wertemenge | `sqlite` aus `New` gestrichen (liefert YAML); Groß-/Kleinschreibung zugelassen | `TestNewFormat`, `TestRunRecordFormat`; `TestParseRecordFormat`, `TestLeserAlleOptionen` |
| Bootstrap reicht `o.Format` durch | `recording.New("yaml")` | `TestRunRecordFormat` (Unit); am Binary `TestE2ERecordFormatSQLite`, `TestE2ERecordFormatQuellen`, `TestE2ERecordFormatSQLiteJournal` (`make test-integration` in einer Kopie, drei Tests rot) |

Grün und äquivalent, wie §7 es führt, von mir nachgefahren: `synchronous(FULL)` aus der URI (der Standard der Bibliothek ist `FULL`, `TestSQLiteEinstellungen` liest `2`) und `db.SetMaxOpenConns(1)`. Beide ändern kein beobachtbares Verhalten; der Kommentar an `oeffne` sagt die Zahl nicht mehr zu.

### Punkt 2: Schema in `tools/schema/schema.yaml`, eingebettetes SQL daraus erzeugt, ein Test legt es in einer In-Memory-Datenbank an und prüft gegen `SPEC-043`. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `schema.yaml`: `message.bytes` entfällt, `fields` `required`, `session.server_parameters` kommt hinzu | `git diff 4914a4a..HEAD -- tools/schema/schema.yaml` zeigt genau diese drei Änderungen; `make schema-validate`: `Validation passed: 0 warning(s)` | bestätigt |
| eingebettetes SQL ist daraus erzeugt | `make schema-generate` gegen `internal/adapters/driven/recording/schema.sql`: gleich bis auf die Kopfzeilen (Zeitstempel) | bestätigt |
| Test gegen `SPEC-043` | `TestSQLiteSchema` besteht; Mutanten: `fields` nullable (`"fields" TEXT`), rot an `TestSQLiteSchema`; `server_parameters` entfernt, rot an vielen Tests, darunter die der Spalten. Die Mutanten S1 bis S9 des Implementers habe ich bis auf diese zwei nicht gefahren | bestätigt |

Der Test prüft `schema.sql` gegen `SPEC-043`, nicht gegen `schema.yaml` (Entscheidung des Nutzers, §3); das ist beschrieben und kein Befund.

### Punkt 3: Bibliothek im Architektur-Gate auf den Recording-Adapter begrenzt (Test), Erlaubnisliste des Lint-Gates führt sie. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `make a-check-negativ` | `a-check-negativ: gruen — alle siebzehn Faelle` am Arbeitsbaum | bestätigt |
| die Gegenprobe fängt die Verletzung | `.a-check.yml` in einer Kopie: `modernc.org/sqlite` zusätzlich im CLI-Adapter zugelassen; `make a-check-negativ` in der Kopie: `ROT in Fall 'cli-sqlite' — erwartet rot, a-check war gruen`, Exit 2 | bestätigt |
| Erlaubnisliste des Lint-Gates | `gomodguard_v2.allowed` führt `modernc.org/sqlite` (`.golangci.yml`, eine Zeile im Diff); Zeile in einer Kopie entfernt, `make lint`: `sqlite.go:18:2: import of package modernc.org/sqlite is blocked … (gomodguard_v2)` | bestätigt |
| keine Lockerung | `.golangci.yml` ändert nur die Erlaubnisliste (Nachzug nach ADR-0034 Punkt 3, `AGENTS.md` §3.6); kein `//nolint` im Diff | bestätigt |

### Übrige DoD-Punkte (konstant je Slice)

- **`make gates` grün:** Der Beleg in §7 gilt für `9ebab6c`; seither änderten `e4ada22` und `a4dfd15` Code und Tests (V-172). Ich habe `make gates` auf dem Endstand nach dem Commit dieses Berichts gefahren; das Ergebnis steht in der Übergabe-Nachricht (nach diesem Lauf folgt kein Commit).
- **Review:** Der Report liegt vor (`075e63a`), aus einem anderen Lauf. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen, Lerneintrag:** noch offen; §7 trägt noch keine Closure-Notiz, die drei Risiken aus §6 tragen die Ausgänge *entfallen* (Grenze der Prüfung, in `LH-FA-22.a` *Grenzen* entschieden), *offen bis Closure* (Bibliothek für sechs Plattformen) und *eingetreten* (`slice-v1-abschluss-sqlite-format-lesen`). Das ist die Reihenfolge, kein Befund. Für die Closure aus meiner Sicht: Risiko 2 ist für den Build belegt (§7: sechs Plattformen gebaut), das Verhalten ist nur unter Linux belegt (`BEO-REPO/verhalten-nur-unter-linux-geprueft`); das Bauen habe ich nicht wiederholt, `linux/amd64` im Image `runtime` habe ich ausgeführt.

---

## 2. Übergaben des Reviews

| Befund | Prüfung | Urteil |
|---|---|---|
| F-614 Pragmas ohne Mutant | `journal_mode(DELETE)` entfernt: `TestSQLiteFremdeWALDatei` rot (allein). `synchronous(FULL)` entfernt: grün, §7 führt es als äquivalent mit Grund | behoben |
| F-615 Kommentare weiter als der Test | `oeffne` sagt die Zahl der Verbindungen nicht mehr, `dsnVon` sagt `%`, `?` und `#` und nennt den Test; `SetMaxOpenConns`-Mutant grün und in §7 als äquivalent geführt | behoben |
| F-616 Plan folgt dem Code | §1 *Übernommen* nennt die Randformen als entschieden (`LH-FA-22.a`), §8 sagt nicht mehr „kein Produktionscode“ | behoben |
| F-617 doppelte Nummer im Aufruf | `schreibe` prüft `rec.Sessions[i].ID == i+1` für alle Sessions als erste Prüfung (Mutanten oben); `LH-FA-22.a` und Plan sagen es; der Satz „lässt die Datei auf dem Stand vor dem Aufruf“ gilt nun für jeden Aufruf mit falschen Nummern | behoben |
| F-618 Windows-Zweig ohne Test | `dsnVon` hat den Zweig weiter; §7 führt ihn als ungeprüfte Plattform-Grenze (W10) | zur Kenntnis (Planner) |
| F-619 Größe | §7 nennt 2794 Zeilen; `git diff --numstat 4914a4a..9ebab6c` zählt 2794 hinzu, 676 Go-Code, 1975 Go-Tests, die Zahlen stimmen. Das Maß gehört dem Planner | zur Kenntnis |
| F-620 N1b rot aus anderem Grund | §7 führt N1b und N1c nach der Entscheidung zu F-617 als „Rückhalt, nicht erzwungen“. Der Primärschlüssel von `session` ist nur noch Rückhalt, die Prüfung vorab fängt doppelte Nummern; die Zeile trägt den Anspruch nicht mehr, weil sie ihn nicht mehr erhebt | behoben |
| F-621 Zufallsteil des Journals | Zusage der Spezifikation durch Wahrscheinlichkeit; Entscheidung des Architect, §7 führt es als akzeptiertes Negativ | zur Kenntnis (Architect) |
| F-622 Schichtzählung des Bootstrap | Architect hat entschieden: Fehlerweitergabe ohne Unterscheidung ist Verdrahtung, zwei Schichten. `bootstrap.go` ändert `+5 −1`, der Zweig unterscheidet nichts. Bestätigt | behoben |

---

## 3. Hexagon und Architektur-Gate

- [ADR-0001](../plan/adr/0001-hexagonale-architektur.md): `git diff 4914a4a..HEAD -- internal/hexagon` ändert nur die Kommentare an `RecordingRepository.Write` und `RecordService`; kein Import der Bibliothek im Core. `.a-check.yml` ist unverändert (die Regel für `modernc.org/sqlite` lag vor).
- [ADR-0013](../plan/adr/0013-sqlite-recording-backend.md) und [ADR-0024](../plan/adr/0024-sqlite-bibliothek.md): `go.mod` führt `modernc.org/sqlite` v1.60.1 und `modernc.org/libc` v1.77.1 exakt wie §3 es nennt, `golang.org/x/sync` steigt indirekt auf v0.23.0; sonst nur indirekte Zeilen.
- `AGENTS.md` §3.2: kein `//nolint` im Diff. §3.3: kein Move im Diff. §3.6: nur die Erlaubnisliste, keine Senkung.
- Port `Write`-Signatur, Handbuch und `README.md` sind unverändert (`git diff --stat 4914a4a..HEAD` nennt weder `docs/user/benutzerhandbuch.md` noch `README.md`).
- Abdeckungstabellen: `make abdeckung-check` Exit 0.
- Commit-Messages `4914a4a..a4dfd15` (elf Commits): kein `SPEC-` und kein `ARC-` (`grep` auf Betreff und Rumpf ohne Treffer, ADR-0035); jede nennt die Slice-Kennung und `LH-FA-22`.

---

## 4. Plan gegen Code

- **§1:** Geliefert sind das Schreiben von `sqlite`, die Option `--format`, die Fabrik, das Schema mit Ersatztest, die Gate-Regel mit Erlaubnisliste. Die Abgrenzungen halten: Lesen, Erkennung, `PGR-E3002` beim Lesen, Roundtrip stehen nicht im Diff; Handbuch und `README.md` nicht im Diff; das Standardformat YAML unverändert (`yaml.go` nicht im Diff).
- **§3.13, Nehmer `slice-v1-abschluss-sqlite-format-lesen`:** Jede Sendung aus §1 *Abgegeben* steht dort in §1 *Übernimmt* oder in der DoD: Formaterkennung und Lesetransaktion, Umkehr des `base64:`-Escapes und `{}`, `PGR-E3002`/`PGR-E3003`/`meta.version`/`meta.empty_sessions`/Waisen und Einstufung, Roundtrip und gleiche Ablehnung beider Formate, der Platzhalter `Load` mit der URI-Bildung („kein zweiter Weg“). Der Nehmer liegt in `next/` und ist nicht geschlossen; Liefer-Punkte drei, Schichten zwei (Adapter und Integrationstest, Bootstrap als Verdrahtung). Der Nehmer nimmt an.
- **§2 Zählung:** drei Liefer-Punkte (DoD 1 bis 3), zwei Schichten (Recording-Adapter, CLI-Adapter); der Bootstrap reicht nur `o.Format` und den Fehler der Fabrik (Verdrahtung nach der Entscheidung des Architect), die Kommentare sind keine Schicht. Das stimmt mit dem Diff (`bootstrap.go` +5 −1).
- **§3:** Jede geänderte Datei steht dort, auch `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk`, `harness/README.md`, `tools/harness/lint-gegenprobe.sh`. Die Aufträge zur Nummern-Vorabprüfung (a) bis (d) sind in `sqlite_schreiben_test.go` als `TestSQLiteNummernNeueDatei`, `TestSQLiteNummernBestehendeDatei`, `TestSQLiteLueckeInDerDatei` umgesetzt (Namen per `grep` vorhanden).
- **§6:** Jede Randform, die ich am Binary fahren konnte, folgt ihrem Ort: Journal (auch `--force`, auch ohne Zieldatei, Reihenfolge Zielpfad vor Journal), Pfad mit Sonderzeichen, Wertemenge und Quellen von `--format`, `replay`/`play` ohne die Option, `fields` immer als Objekt, `interaction.sql` roh, Zwischenstand. Nicht am Binary erreichbar und nur über Tests belegt: Nachholen nach fehlgeschlagener Transaktion, verschwundene oder ersetzte Datei, Fehler beim Schließen, Nummern des Aufrufs (der Dienst vergibt sie lückenlos), Journal der temporären Datei (V-174 zu `--record-empty-sessions`).
- **§7:** Namensabgleich: aus §3, §6 und §7 zog ich jeden Namen `Test…` (56) und jeden Dateinamen in Backticks (37) und suchte `^func <Name>(` in `internal/**` und `test/**` sowie die Dateien im Baum: alles gefunden, bis auf Prosa-Treffer (`Test`, `Testnamen`, `Tests`, `TestSQLite…` als Gattungsname), Baseline-Dateien (`grundlagen-bootstrap.md`, `grundlagen-traceability.md`, `modul-06-roadmap.md`) und `interaction.sql` (eine Spalte). Alle genannten Commits existieren. Die Zahlen zur Größe stimmen. Falsch bzw. veraltet sind die Zusammenfassung der grünen Mutanten und der Gate-Beleg (V-172).
- **§8:** Sub-Area `REPO`, Greenfield, wie in `harness/conventions.md`; Zählung wie oben; Register gesichtet.

---

## 5. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-172 | LOW | §7 trägt drei Belege, die dem Endstand nicht mehr entsprechen (`AGENTS.md` §3.9). (a) Die Zeile „alle rot bis auf drei äquivalente (nächster Absatz)“ steht vor dem Nachtrag; §7 führt danach sechs grüne Mutanten (B2, W10, W14, `synchronous(FULL)`, `SetMaxOpenConns(1)`, N1c), der „nächste Absatz“ ist der Nachtrag. (b) Der Lauf `make gates` in der Läufe-Tabelle gilt für `9ebab6c`; danach änderten `e4ada22` und `a4dfd15` `sqlite.go`, die Tests und den Plan, und der Beleg des Implementers für `make gates` auf dem Endstand fehlt in §7 (die DoD verlangt ihn; ich habe ihn selbst gefahren, siehe Übergabe). (c) Der Absatz *Stand* nennt „vier Commits seit `4914a4a`“, es sind elf. | `docs/plan/planning/in-progress/slice-v1-abschluss-sqlite-format.md` · „alle rot bis auf drei äquivalente“; · „`make gates` | auf `9ebab6c`“; · „in vier Commits seit `4914a4a`“ | Implementer: die drei Stellen auf den Endstand ziehen oder den Stand ausdrücklich nennen; der Planner trägt den Gate-Beleg vor der Closure nach (Lauf `a4dfd15` oder Endstand) |
| V-173 | INFO | DoD-Punkt 1 verlangt „Test über das Binary; Inhalt der Datei gegen die Liste in `SPEC-043` mit rohem SQL, ohne `sqlite3`“. Die Tests über das Binary (`sqlite_format_e2e_test.go`) lesen die Datei nur als Bytes (Kopfzeile, Teilstrings); der Inhalt gegen `SPEC-043` wird allein im Paket des Adapters (`TestSQLiteInhalt`) mit rohem SQL geprüft, nicht an einer Datei, die das Binary schrieb. Zwischen beiden steht der Bootstrap-Test `TestRunRecordFormat` und der Mutant E1 (rot an den drei Tests über das Binary). Ich habe die Kette einmal geschlossen: Binary, echter Upstream, unabhängig gelesen, Inhalt nach `SPEC-043` (Punkt 1 oben), ohne Abweichung. Ob die DoD-Klammer beide Teile an einem Test fordert, ist Lesart des Planners. | `test/integration/sqlite_format_e2e_test.go` · `TestE2ERecordFormatSQLite`; `internal/adapters/driven/recording/sqlite_test.go` · `TestSQLiteInhalt` | Planner: Lesart; Validator: kein Befund am Verhalten |
| V-174 | INFO | `--record-empty-sessions` gibt es im gebauten Binary noch nicht (`PGR-E2001`); die Zeile `meta.empty_sessions` des Schreibens (§6 und `SPEC-043`) ist deshalb nur über den Adapter-Test (`TestSQLiteEmptySessions`, Mutant rot) erreichbar, nicht am Binary. Die Option liefert `slice-v1-abschluss-sessions` (`next/`); dessen §1 und DoD nennen `sqlite` nicht. Das Modell trägt das Feld (`model.Recording.EmptySessions`), YAML und SQLite lesen es gleich, es ist keine Sendung im Sinn von `AGENTS.md` §3.13. | `docs/plan/planning/next/slice-v1-abschluss-sessions.md` · „`--record-empty-sessions`“ | Planner: ob `slice-v1-abschluss-sessions` die Wirkung mit `--format sqlite` mitprüft (ein Satz in dessen DoD), oder ob der Test im Adapter genügt |

Die Falsifikationsversuche, die nichts fanden: Journal ohne Zieldatei, mit `--force`, zusammen mit vorhandener Zieldatei ohne `--force`; Zielpfad mit `?`, `#`, `%25` und Leerzeichen; `--force` über eine Datei mit sieben Sessions; ungültige Werte aus Option, Umgebung und leerem Wert; Verzeichnis des Ziels fehlt (`PGR-E3001` vor der ersten Verbindung); `PRAGMA integrity_check` und `foreign_key_check` der geschriebenen Datei; zwei Escape-Fälle (`base64:`, `\`) am Binary; `replay`/`play` auf der SQLite-Datei (Zwischenstand wie in §6 beschrieben); das Image `runtime` (distroless, nonroot) schreibt die Datei mit den Rechten `0644` unter umask 022.

Befunde: 0 HIGH, 0 MEDIUM, 1 LOW, 2 INFO. Alle drei DoD-Verhaltenspunkte sind bestätigt, Punkt 1 auch gegen eine Gegenstelle und einen Leser, die nicht der eigene Code sind; rund 40 Mutanten der Unit-Tests, ein Mutant am Binary (`make test-integration`), einer des Architektur-Gates und einer des Lint-Gates wurden von mir gefahren, alle rot aus dem benannten Grund, außer den zwei als äquivalent geführten.

**Finding-Klassen dieses Laufs:** Belege im Plan vom Endstand überholt (Summe, Gate-Lauf, Commit-Zahl) · DoD-Klammer „Binary und Inhalt“ auf zwei Tests verteilt · Wirkung einer Option am Binary erst durch einen späteren Slice erreichbar

**Closure-blockierend:** nein. V-172 ist ein Nachzug an §7 (drei Stellen), keine Änderung am Code; V-173 und V-174 sind Lesarten des Planners.

**Übergabe an den Planner:** dieser Bericht mit V-172 (an den Implementer, Gate-Beleg an den Planner), V-173 und V-174 (an den Planner); dazu für die Closure die Ausgänge der drei Risiken aus Abschnitt 1 (*Übrige DoD-Punkte*) und die Klassen oben. Die Entscheidungen des Architect zu F-617 und F-622 tragen; der Nehmer `slice-v1-abschluss-sqlite-format-lesen` nimmt alle Sendungen an.

**Aufgeräumt:** Container `ver-sq-pg`, `ver-sq-rec`, `ver-sq-rec2`, `ver-sq-rec3` und die Läufe `pgr-it-*`; Netz `ver-sq-net`; Images `ver-sq-rt`, `ver-sq-deps`; die Kopien und Verzeichnisse mit dem Präfix `ver-` aus dieser Verifikation im Scratchpad. `docker ps -a`, `docker network ls`, `docker volume ls` und `docker images` mit `ver-sq` und `pgr-it` danach leer; die Images des Arbeitsbaums (`pgwire-recorder:*`) blieben stehen.
