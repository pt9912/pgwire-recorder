# Verifikation: slice-v1-abschluss-sqlite-format-lesen — 2026-10-11

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-11-review-slice-v1-abschluss-sqlite-format-lesen.md`, F-623 bis F-630); die Übergaben prüfe ich nach. Ich repariere nichts.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-sqlite-format-lesen.md` (Lese-Hälfte des Formats `sqlite`), Kopf und §1 bis §8, gegen den Gesamt-Diff `97f3ca6..e0b56b1` (21 Dateien, 2702 Zeilen hinzu, 46 weg; Go: 886 Zeilen Code, 1444 Zeilen Tests). Implementer: `45af81d`, `60a448f`, `18695d4`, `aee68b3`, `415b725`, `a6c5b56`, `182355c`, `e0b56b1`. Architect: `ea1ef40`, `a6bcc45`, `39cc927`, `c43b775`. Review: `cbb922c`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingang:**

- der Plan ganz am Stand `e0b56b1`, besonders §1 (*Übernimmt*, *Ausdrücklich NICHT*), §2 (DoD), §3 (mit den Aufträgen), §6 (*Randformen*), §7 (*Belege des Implementers* und die drei Nacharbeiten, Namensabgleich), §8
- `spec/spezifikation.md` `LH-FA-22.a` (*Wahl und Erkennung*, *Lesen von `sqlite`*, *Grenzen*), `SPEC-043`, der Diff `97f3ca6..HEAD` daran
- die Nehmer `slice-v1-abschluss-sqlite-format-doku` (§1, DoD) und `slice-v1-abschluss-zeitangaben` (§1, DoD, §4)
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen
- am Stand `e0b56b1`: `leser.go`, `sqlite_lesen.go`, `sqlite_lesen_baum.go`, `sqlite_json_lesen.go`, die Diffs von `sqlite.go` und `bootstrap.go`, die Tests `TestLesenAbgeschnitten`, `TestLesenFehlendeTabelleOderSpalte`, `TestLesenJSONZahlen` und `sqlite_lesen_e2e_test.go`
- `AGENTS.md` §3.9 bis §3.13

Beim Start war der Arbeitsbaum sauber, HEAD `e0b56b1`. Der Commit dieses Berichts enthält nur diese Datei.

**So wurden die Proben gebaut.** Alles Docker-only, Scratch-Präfix `ver-sl-`, nie im Repo.

- **Binary:** `git archive HEAD` in eine frische Kopie, daraus die Stufe `runtime` des `Dockerfile` (distroless, nonroot) mit eigenem Tag. Gegenstelle: ein echter PostgreSQL-17-Server (das gepinnte Image aus `harness/mk/integration.mk`, `trust`) in einem eigenen Netz; Clients: `psql` desselben Images. Beschädigte, fremde und angepasste Dateien baute ich mit dem `sqlite3`-Modul von Python 3.12 im Image `python:3.12-slim` (rohes SQL, unabhängig vom Leser des Adapters) aus der Aufzeichnung des Binaries.
- **Unit-Mutanten:** je Mutant eine frische Kopie (`cp -r` ohne `-p` von `go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`), genau eine Ersetzung (das Skript bricht ab, wenn der Suchtext nicht genau einmal trifft), `gofmt -l` der Datei leer, im Image der Stufe `deps`, `--network none`, `go test -count=1 -run <Tests>` über `./internal/adapters/driven/recording`. Der Grundlauf der Pakete `recording` und `bootstrap` war grün.
- **Integrations-Mutant:** eine Kopie des ganzen Baums, `make test-integration` darin.
- **Danach:** Container, Netz, Images und Kopien entfernt (siehe Ende).

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Roundtrip-Gleichheit beider Formate über das Binary, Formaterkennung nach Inhalt, Datei nur lesend und nie verändert. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Ende-zu-Ende gegen einen echten Upstream | `record --format sqlite` und `record --format yaml` mit denselben Anfragen gegen PostgreSQL 17 (einfache Anfragen, Datenzeile mit `base64:abc`, `\x41`, leerem Text, `NULL`, `café 日本`, temporäre Tabelle, Fehlerantwort, Extended Query über `\bind` mit `NULL`, Latin-1-Byte `0xE9` in einer Datenzeile). Exit 0 nach `SIGINT`, die Datei beginnt mit `SQLite format 3` und NUL | bestätigt |
| `replay` liest beide Formate gleich | Dieselben `psql`-Skripte gegen `replay --input rec.sqlite` und `replay --input rec.yaml`: die Ausgaben beider Formate sind byte-gleich (`diff` leer) | bestätigt |
| `play` spielt beide Formate gleich ein | Eine Aufzeichnung mit `CREATE TABLE` und `INSERT` (Werte `NULL`, `café`, `base64:abc`, `\x41`) in beiden Formaten, eingespielt in zwei leere Datenbanken: gleicher Inhalt der Tabelle. Mit der Aufzeichnung der Fehlerantwort enden beide Formate mit derselben Meldung `PGR-E4004`, Exit 4 | bestätigt |
| Erkennung nach Inhalt, nicht nach Endung | Dateien mit der Endung `.file` und Inhalt wie eine leere Datei, 15 Bytes und `SQLite format 3` ohne NUL lesen als YAML und werden als beschädigt gemeldet (`PGR-E3003`); die Gegenrichtung (SQLite-Inhalt, andere Endung) belegen `TestLeserErkennungOhneEndung` und die Mutanten unten | bestätigt |
| Datei nur lesend, nie verändert | `sha256sum` von `rec.sqlite` vor und nach drei `replay`-Läufen gleich; `replay` läuft auch auf einem schreibgeschützten Mount (`-v …:ro`); das Verzeichnis trägt keine Nebendatei, die der Leser angelegt hätte | bestätigt |
| Eingabepfad wörtlich | Dateien `we ird?a#b%25.sqlite`, `q?x.sqlite`, `h#x.sqlite`, `p%41.sqlite` (die Datei `pA.sqlite` gibt es nicht) lesen `replay` (läuft) und `play` (verbindet sich, `PGR-E4002` am toten Ziel); führendes `//` (`//c/sub/rec.sqlite`) liest beide | bestätigt |
| Journal-Rest neben der Datei | `jr.sqlite` mit einer Datei `jr.sqlite-journal`: `PGR-E3001` „attempt to write a readonly database“; nicht zugesagt (Grenze in `LH-FA-22.a`), die Beobachtung deckt sich mit §6 | bestätigt (Grenze) |
| `offset_ms` | `-1` und `1.5`: `PGR-E3003`; `250` und `9223372036854775807`: gelesen | bestätigt |

**Bewusstes Brechen** (Modul 11). Jede Mutation in einer frischen Kopie, genau eine Ersetzung, `gofmt` sauber, Grundlauf grün. Alle rot, die Fehlermeldung des Tests nennt den benannten Grund:

| Zusage | Mutation | roter Test |
|---|---|---|
| nur lesend geöffnet | `?mode=ro&` zu `?mode=rw&` | `TestLesenNurLesend` (Meldung: „Schreibversuch in der Lesetransaktion: <nil>“) |
| Erkennung an den ersten 16 Bytes | immer `yaml` (Bedingung `&& false`) | `TestLeserErkennungOhneEndung`, `TestLeserErkennungAnDenErsten16Bytes`, `TestLesenGleicheAblehnungBeiderFormate` |
| … | immer `sqlite` (`true \|\|`) | dieselben drei |
| eine einzige Lesetransaktion | `SetMaxOpenConns(2)` | grün, äquivalent: die Transaktion hält ihre eine Verbindung; die Anfragen außerhalb der Transaktion (O6) fing der Implementer mit einem anderen Mutanten |
| Eingabepfad wörtlich, auch mit führendem `//` | `filepath.Abs` entfällt | `TestLesenPfadMitFuehrendemDoppelstrich` |
| `replay` erkennt selbst (K1) und `play` erkennt selbst (K2) | `recording.YAML{}` statt `recording.Leser{}` an beiden Stellen in `bootstrap.go`, `make test-integration` in einer Kopie des ganzen Baums | rot: `TestE2EReplaySQLite`, `TestE2EReplayExtendedSQLite`, `TestE2EPlaySQLite`, `TestE2ELesenAbgelehnt` |

### Punkt 2: unbekannte Version `PGR-E3002`, ungültige Aufzeichnung `PGR-E3003`, in beiden Formaten gleich. **Bestätigt, mit einer Lücke (V-175).**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `version`: jede kanonische Zahl außer 1 ist `PGR-E3002`, alles andere an `version` `PGR-E3003` | Am Binary: `2` und `99999999999999999999` `PGR-E3002`; `01`, `+1`, `eins` `PGR-E3003`; Exit 3, `replay` wie `play` | bestätigt |
| Reihenfolge Formatkennung, Version, Inhalt | `format` falsch bei `version` 2 meldet `PGR-E3003` (`TestLesenMeta`, Mutanten rot) | bestätigt |
| Ablehnungen am Binary, `replay` und `play` gleich (Exit 3, kein Start) | beschädigter Seitenkopf (0xFF je Seite), Kopf ohne Inhalt, Kopf und Müll, fremdes Schema, fehlende Tabelle `meta`, fehlende Spalten `startup`, `server_parameters`, `fields`, `kind`, `direction`, `type`, `group_no`, `offset_ms`, `position`, `session_id`, `sequence`, `id`, `key` und `value`, doppelter Schlüssel in `meta`, unbekannter Schlüssel in `meta`, verwaiste Interaktion und Nachricht, `type` und `direction` unbekannt, Text in `sequence`, `sql` `NULL` bei einfacher Anfrage, JSON ungültig, Zahl `1.5e0`, Base64 ungültig, `null` als Wert, doppelter Schlüssel, `<<` als Schlüssel, einzelnes Surrogat `\ud800` | bestätigt |
| gleiche Ablehnung beider Formate | `TestLesenGleicheAblehnungBeiderFormate` (unbekannter Schlüssel, `null`, doppelter Schlüssel, Version); Mutanten oben (immer `yaml`, immer `sqlite`) rot; Mutant „`PGR-E3003` statt `PGR-E3002` für Version ≠ 1“ rot an `TestLesenMeta` und `TestLesenGleicheAblehnungBeiderFormate` | bestätigt |
| Schlüssel `<<` (F-630) | Am Binary lehnen `replay` und `play` eine Datei mit `{"<<":"…"}` in `startup` mit `PGR-E3003` ab; Mutant `c.Value == "<<" \|\|` entfällt in `ohneVerweise`: rot an `TestLesenMergeKey` | bestätigt |
| **fehlende Spalte** `interaction.sql` | **Nicht abgelehnt**, siehe V-175 | **Lücke** |

### Punkt 3: Einstufung der Fehler nach primärem Ergebniscode, abgeschnittene Datei `PGR-E3003`. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Datei fehlt, Verzeichnis, nicht lesbar (Rechte 0) | `PGR-E3001` mit Pfad und Ursache, Exit 3, `replay` wie `play` | bestätigt |
| gesperrte Datei | Eine zweite Verbindung hält `BEGIN EXCLUSIVE`: `replay` wartet fünf Sekunden (gemessen 5 s) und meldet `PGR-E3001` „database is locked (5) (SQLITE_BUSY)“ | bestätigt |
| abgeschnitten an der Seitengrenze, mitten in einer Seite, auf 100 und 4096 Bytes, beschädigte Seiten | `PGR-E3003` („database disk image is malformed“, „file is not a database“), Exit 3 | bestätigt |
| zehn primäre Codes `PGR-E3001`, jeder andere `PGR-E3003` | `SQLITE_FULL` aus der Liste entfernt: rot an `TestLesenStufeNachErgebniscode`; ohne die Maske `& 0xff`: rot; ohne den Sonderfall `IOERR_SHORT_READ`: rot | bestätigt |
| jeder Fehler der Bibliothek `PGR-E3001` (statt Stufung) | rot an `TestLesenFehlendeTabelleOderSpalte` und `TestLesenAbgeschnitten` | bestätigt |
| Fehler ohne Typ der Bibliothek (Abbruch von `ctx`) `PGR-E3001` | Rückfall `PGR-E3003`: rot an `TestLesenNichtLesbarUndKopfOhneInhalt` und `TestLesenAbbruchWaehrendDesLesens` | bestätigt |

`IOERR_SHORT_READ` liefert keine echte Datei (die Bibliothek meldet `NOTADB` und `CORRUPT`); die Zuordnung trägt der Test mit der Zahl, das hält §7 selbst fest.

### Zusagen der Meldung und der Randformen (Nacharbeiten des Reviews)

| Zusage | Beleg | Urteil |
|---|---|---|
| Meldung nennt Pfad und Art, nie Wert, nie Zeile (F-623) | 18 präparierte Dateien, jede mit dem Wert `geheim-wert` an einer anderen Stelle (Formatkennung, `version`, Schlüssel in `meta`, `startup`, Schlüssel und Feld in `fields`, doppelter Schlüssel, `type`, `direction`, `kind`, JSON-Text, Base64, Zahl, `sequence`, `<<`, `null`, Waise): `grep -c geheim` über alle Ausgaben von `replay` und `play` ist 0, ebenso `Zeile`/`line` in Meldungen zu `sqlite`. Mutanten: Pfad entfällt bei eigenen Fehlern (MP1) und Wert in eigener Meldung (ME2): rot an `TestLesenMeldungOhneWertUndZeile`, die Fehlertexte nennen „Meldung ohne den Pfad der Datei“ | bestätigt |
| einzelnes Surrogat `PGR-E3003` (Entscheidung des Architect) | Am Binary `\ud800` in `startup`: `PGR-E3003`; Mutant „Prüfung entfällt“: rot an `TestLesenSurrogate` | bestätigt |
| Umkehr von `\` und `base64:` | Mutant „RawStdEncoding statt StdEncoding“: rot an `TestLesenStringOhneUTF8`, `TestLesenUmkehrDerUmschreibung` | bestätigt |
| Pfad mit führendem `//` auch beim Schreiben (F-624) | Mutant oben rot | bestätigt |
| Entscheidung zu F-623 (Art auf dem Weg über `Unmarshal` enger gefasst) | `LH-FA-22.a` und §6 sagen „Aufzeichnung beschädigt“ mit Nummern; das Binary meldet genau das (`unbekannter Schlüssel`, `null`, `<<`, doppelter Schlüssel: `Inhalt der SQLite-Datei beschädigt: Aufzeichnung beschädigt`; Kopf einer Interaktion: `Session 1, Interaktion 1`) | bestätigt |
| `yaml`-Meldungen vorerst unverändert (Entscheidung des Nutzers) | `yaml.go` ist im Diff nicht berührt. Folge, die der Nutzer mit der Entscheidung in Kauf nahm: Die Meldung für eine 15-Byte-Datei im Binary nennt Zeile und Anfang des Inhalts (Meldung `yaml: unmarshal errors: line 1: cannot unmarshal !!str` und dem Anfang des Inhalts); siehe V-178 | bestätigt, Beobachtung für das Register |

### Übrige DoD-Punkte (konstant je Slice)

- **`make gates` grün:** Der Beleg in §7 gilt für `45af81d`; seither änderten `18695d4`, `415b725` und `182355c` Code und Tests (V-176). Ich habe `make gates` auf dem Endstand nach dem Commit dieses Berichts gefahren; das Ergebnis steht in der Übergabe-Nachricht (nach diesem Lauf folgt kein Commit).
- **Review:** Der Report liegt vor (`cbb922c`), aus einem anderen Lauf. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen, Lerneintrag:** noch offen; §7 trägt noch keine Closure-Notiz, die drei Risiken aus §6 tragen noch keinen Ausgang. Das ist die Reihenfolge (die Closure macht der Planner), kein Befund. Aus meiner Sicht für die Closure: Risiko 1 (Zwischenstand) ist *entfallen*, belegt durch Punkt 1; Risiko 2 (abgeschnittene Datei an beiden Stellen) ist belegt durch Punkt 3, mit dem Vorbehalt zu `IOERR_SHORT_READ` oben; Risiko 3 (nur unter Linux geprüft) bleibt *weiter offen* ins Register.

---

## 2. Übergaben des Reviews

| Befund | Prüfung | Urteil |
|---|---|---|
| F-623 Meldung ohne Pfad oder Art | Pfad in jeder Meldung zu einem Fehler des Inhalts (18 Dateien am Binary); die Art auf dem Weg über `Unmarshal` ist laut Entscheidung des Architect enger gefasst und deckt sich mit `LH-FA-22.a` und §6; MP1 bis MP6 in §7 | behoben |
| F-624 O5 nicht äquivalent | Test `TestLesenPfadMitFuehrendemDoppelstrich` vorhanden, Mutant rot; §7 berichtigt | behoben |
| F-625 `rows.Err()` ohne Mutant | §7 führt U7 als grün, mit Begründung (Zweig als Absicherung der Bibliothek, durch keinen Test belegt); die Spezifikation sagt „während des Lesens“, der Test bricht zwischen den Tabellen ab | zur Kenntnis (Grenze, in §7 benannt) |
| F-626 Kommentare | `Load` nennt, was der Leser prüft und was `Unmarshal`; `JSONKnoten` in `export_test.go` nennt `!!map` und `!!seq` | behoben |
| F-627 Plan folgt dem Code | §8 nennt 2330 Zeilen Go (886 und 1444), das stimmt mit `git diff --numstat 97f3ca6..182355c -- '*.go'`; Rest zu V-176 | behoben |
| F-628 Größe | 2330 Zeilen Go gegen „rund 2000“ aus §4; die Entscheidung des Planner steht nicht im Plan, V-177 | zur Kenntnis (Planner) |
| F-629 Plattformzweig | `uriVon` trägt den Zweig weiter; §6 führt das Risiko, Ausgang bei der Closure | zur Kenntnis (Planner) |
| F-630 Schlüssel `<<` | Entscheidung des Architect in §7 und `LH-FA-22.a`; Test `TestLesenMergeKey` vorhanden, Mutant rot, kein Produktcode im Diff dazu | behoben |

---

## 3. Hexagon und Architektur-Gate

- [ADR-0001](../plan/adr/0001-hexagonale-architektur.md): `git diff --stat 97f3ca6..HEAD -- internal/hexagon` ist leer; die Bibliothek wird nur im Recording-Adapter eingebunden (`.a-check.yml` unverändert, `make a-check` läuft in `make gates`).
- [ADR-0013](../plan/adr/0013-sqlite-recording-backend.md) und [ADR-0024](../plan/adr/0024-sqlite-bibliothek.md): `go.mod` und `go.sum` sind unverändert; `modernc.org/sqlite/lib` gehört zum Modul des Treibers.
- `AGENTS.md` §3.2: `grep -rn nolint cmd internal test` ohne Treffer. §3.3: kein Move im Diff (`git diff --name-status` zeigt nur `M` und `A`). §3.6: `.golangci.yml` unberührt.
- Handbuch und `README.md`: `git diff --stat 97f3ca6..HEAD -- docs/user/benutzerhandbuch.md README.md` leer.
- Abdeckungstabellen: `make abdeckung-check` Exit 0; `make kopf-check` Exit 0.
- Commit-Messages `97f3ca6..e0b56b1` (13 Commits): kein `SPEC-` und kein `ARC-` (`grep` über Betreff und Rumpf ohne Treffer, [ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md)); jede nennt die Slice-Kennung und `LH-FA-22`.

---

## 4. Plan gegen Code

- **§1:** Geliefert sind die Erkennung an den ersten 16 Bytes, der Leser in einer Lesetransaktion mit `mode=ro` und fünf Sekunden Sperrfrist, die Umkehr von `base64:`, `{}` als fehlendes Feld, die Fehlerklassen, `meta.version` und `meta.empty_sessions`, der Ersatz des Platzhalters `Load` und der Roundtrip über das Binary. Die Abgrenzungen halten: Schreiben, Fabrik und Schema stehen nicht im Diff (nur `uriVon` aus `dsnVon` herausgelöst); Handbuch und `README.md` nicht im Diff; `yaml.go` unberührt; die Zeitangaben sind nur Spalte mit geprüftem Typ.
- **§3.13, Nehmer:** `slice-v1-abschluss-sqlite-format-doku` (in `next/`, nicht geschlossen) führt den Lese-Slice in Kopf, Ziel, Herkunft, §1, DoD, §3 und §4; sein Auftrag „Grenzen dieses Slice“ trägt die Grenzen des Lesens (Journal, Speicher, `<<`). `slice-v1-abschluss-zeitangaben` (in `next/`) nennt den Lese-Slice im Start (§4) und trägt in der DoD „Recording in YAML und SQLite per Roundtrip laden, ungültiger Wert `PGR-E3003`“; das Lesen von `offset_ms` ist geliefert (Punkt 1). Beide nehmen an. Nachzählung im Geber: drei Liefer-Punkte, eine Schicht (Recording-Adapter); der Bootstrap tauscht nur einen Wert (Verdrahtung, `bootstrap.go` `+2 −2`).
- **§3:** Jede geänderte Datei steht in der Tabelle; die Aufzeichner `aufnehmenExtended` und `nimmAufEinspielAblaufAls` existieren (`grep`); `tools/test/run-integration-tests.sh` ist unverändert.
- **§6:** Jede Randform, die ich am Binary fahren konnte, folgt ihrem Ort: Typ-Prüfung der Spalten, `fields` fehlt oder kein Objekt, Erkennung (leer, kürzer, ohne NUL), nur lesend auch auf schreibgeschütztem Mount, `meta`-Regeln, zusätzliche Tabellen und Spalten ignoriert (`extras.sqlite` liest), verwaiste Zeilen, Fehlerklassen, Journal-Rest (nicht zugesagt), Pfad wörtlich, Zahl, Surrogat, Meldung, Abbruch, doppelter Schlüssel in `meta`, Gruppen, `<<`. Nur über Tests belegt: Schließen nach dem Lesen, `ctx`-Abbruch, `IOERR_SHORT_READ`, Speicher des Lesers (ohne Obergrenze, keine Zusage). Die Entscheidungen des Nutzers vom 2026-10-11 stehen im Plan (§6 *Speicher*; §7 *Hinweise an den Architect*, Punkte 2 und 7) und in `LH-FA-22.a` *Grenzen*.
- **§7:** Namensabgleich: aus §2, §3, §6 und §7 zog ich jeden Namen `Test…` (36 verschiedene) und jeden Dateinamen in Backticks und suchte `^func <Name>(` in `internal/**` und `test/**` sowie die Dateien im Baum: alles gefunden, bis auf Prosa-Treffer (`Test`, `Testnamen`, `Tests`) und `interaction.sql` (eine Spalte). Alle genannten Commits existieren. Die Zahlen zur Größe stimmen (V-176 nennt die eine Unschärfe).
- **§8:** Sub-Area `REPO`, Greenfield, wie in `harness/conventions.md`; Zählung wie oben; Register gesichtet.

---

## 5. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-175 | MEDIUM | `LH-FA-22.a` *Lesen von `sqlite`* sagt zu: Eine Datei mit fehlender Tabelle oder Spalte ist `PGR-E3003`, und DoD-Punkt 3 bindet die Einstufung daran. Fehlt die Spalte `interaction.sql`, liest der Leser eine Datei, die nur einfache Anfragen trägt, **ohne Fehler**: Die Anfrage `SELECT "sql" FROM …` mit der Spalte in doppelten Anführungszeichen gilt in SQLite, wenn es die Spalte nicht gibt, als Zeichenkette `'sql'`, und der Leser nimmt sie als SQL-Text der Interaktion. Belege: (a) Am Binary: aus einer mit `record --format sqlite` geschriebenen Datei mit zwei einfachen Anfragen `ALTER TABLE interaction DROP COLUMN sql`; `replay --input` startet (kein `PGR-E3003`), und ein Client mit der Anfrage `sql` bekommt die aufgezeichnete Antwort `1`. (b) Im Paket: `sqliteDatei(t, sessionsMit(1))`, `lockere`, dann dasselbe `DROP COLUMN`; `Leser{}.Load` liefert `err=<nil>` und `Request.SQL == "sql"`. Für jede andere Spalte lehnt der Leser ab (ich habe jede der 17 Spalten einzeln entfernt), weil der Ersatzwert dort an der Art, am Wert oder an der Form scheitert; bei `interaction.sql` scheitert er nur für Extended-Interaktionen (`sql` ist dort `NULL` vorgeschrieben). Der Test `TestLesenFehlendeTabelleOderSpalte` führt den Fall `interaction.sql`, bleibt aber aus einem anderen Grund grün: Die Fixture `lesefixture` trägt eine Extended-Interaktion, deren „`sql`“ dann ungleich `NULL` ist. Zusage der Spezifikation weiter als der Test; Mutation und Test dazu fehlen (`AGENTS.md` §3.10). | `internal/adapters/driven/recording/sqlite_lesen.go` · `SELECT "session_id", "sequence", "type", "offset_ms", "sql" FROM "interaction"`; `internal/adapters/driven/recording/sqlite_lesen_test.go` · `"interaction.sql":  \`ALTER TABLE interaction DROP COLUMN sql\`` | Implementer: den Fall erkennen (die Spalten vorab prüfen oder die Anfrage ohne Rückfall auf eine Zeichenkette stellen), Test mit einer Datei nur aus einfachen Anfragen und Mutation dazu, §7 und §6 nachziehen; Architect nur, falls `LH-FA-22.a` anders gelesen werden soll |
| V-176 | LOW | §7 trägt Belege, die dem Endstand nicht mehr entsprechen (`AGENTS.md` §3.9). (a) Die Läufe-Tabelle nennt `make gates` und `make lint` für `45af81d`; seither änderten `18695d4` (Surrogat, Meldungen, `ctx`), `415b725` (Pfad, Art, `//`, Abbruch) und `182355c` (Test zu `<<`) Code und Tests, und weder die drei Nacharbeiten noch der Stand-Absatz nennen einen Lauf von `make lint` oder `make gates` auf dem Endstand (die DoD verlangt `make gates` grün; ich habe es selbst gefahren, siehe Übergabe). (b) Der Absatz *Stand* sagt „in einem Code-Commit, `45af81d` seit `97f3ca6`“; es sind vier Commits mit Code oder Tests. (c) §8 führt „2330 Zeilen Go … (`git diff --numstat 97f3ca6..415b725 -- '*.go'`)“; dieser Befehl zählt 2311, erst `97f3ca6..182355c` (der zweite Befehl im Satz) zählt 2330 mit 886 und 1444. | `docs/plan/planning/in-progress/slice-v1-abschluss-sqlite-format-lesen.md` · „`make gates` \| auf `45af81d`“; „in einem Code-Commit“; „Gemessen nach der Nacharbeit zum Review“ | Implementer: die drei Stellen auf den Endstand ziehen; der Planner trägt den Gate-Beleg vor der Closure nach |
| V-177 | INFO | Die Rückführungsschwelle aus §4 („rund 2000 Zeilen“) ist überschritten (2330 Zeilen Go, davon 1444 Tests, überwiegend Tabellen mit rohem SQL). §7 hält nur die Empfehlung des Architect fest („hinnehmen“); die Entscheidung des Planner steht nirgends (F-628). Die Prüfbarkeit in einer Sitzung hat der Review festgestellt; ich konnte sie an dieser Größe nicht gegenprüfen. | `docs/plan/planning/in-progress/slice-v1-abschluss-sqlite-format-lesen.md` · „Hinweise an den Architect“, Punkt 8 | Planner: Entscheidung zu §4 in die Closure-Notiz |
| V-178 | INFO | Für die Closure: (a) Die Entscheidung des Nutzers, die Regel „Meldung ohne Wert und Zeile“ nur für `sqlite` gelten zu lassen, hat eine am Binary sichtbare Folge: `replay` auf einer 15-Byte-Datei (Anfang `SQLite format 3`, ohne Kopf) meldet `yaml: unmarshal errors: line 1: cannot unmarshal !!str` mit Zeile und dem Anfang des Inhalts. §7 trägt „geht als Beobachtung ins Register (der Planner trägt sie bei der Closure ein)“; der Eintrag fehlt noch. (b) Risiko 2 aus §6 (abgeschnittene Datei) ist an beiden Stellen belegt, aber `IOERR_SHORT_READ` erreicht keine echte Datei; der Ausgang *entfallen* braucht diese Begründung. (c) Ein Rollback-Journal neben der Datei liest das Binary unter Linux als `PGR-E3001` („attempt to write a readonly database“); nicht zugesagt, kein Befund. | `docs/plan/planning/in-progress/slice-v1-abschluss-sqlite-format-lesen.md` · „Hinweise an den Architect“, Punkt 2 | Planner |

Die Falsifikationsversuche, die nichts fanden: Pfade mit `?`, `#`, `%`, Leerzeichen und führendem `//` für beide Kommandos; Datei auf schreibgeschütztem Mount; Datei im Write-Ahead-Log-Modus (liest; nicht zugesagt) und Datei mit zusätzlicher Tabelle, Spalte und Index (liest, wie zugesagt); Sperre länger als fünf Sekunden (`PGR-E3001` nach 5 s); `offset_ms` negativ, gebrochen, sehr groß; Wert `geheim-wert` an 18 Stellen ohne Auftreten in einer Meldung; `replay` und `play` verhalten sich bei jeder Ablehnung gleich (gleicher Code, Exit 3, `play` verbindet sich nicht).

Befunde: 0 HIGH, 1 MEDIUM, 1 LOW, 2 INFO. Die drei DoD-Verhaltenspunkte sind bestätigt, Punkt 2 mit der Lücke aus V-175; 16 Mutanten der Unit-Tests und eine Doppelmutation am Binary (`make test-integration`) wurden von mir gefahren, alle rot aus dem benannten Grund, außer dem einen als äquivalent geführten (`SetMaxOpenConns(2)`).

**Finding-Klassen dieses Laufs:** Fehlende Spalte durch den Rückfall auf eine Zeichenkette gelesen, Test grün aus anderem Grund (Zusage der Spezifikation weiter als der Test) · Belege im Plan vom Endstand überholt (Gate-Lauf, Commit-Zahl, Messbefehl)

**Closure-blockierend:** V-175 ja, bis der Implementer den Fall erkennt oder der Architect `LH-FA-22.a` anders fasst; V-176 ist ein Nachzug an §7 (drei Stellen), V-177 und V-178 sind Entscheidungen und Einträge des Planner.

**Übergabe an den Planner:** dieser Bericht mit V-175 (an den Implementer), V-176 (an den Implementer, Gate-Beleg an den Planner), V-177 und V-178 (an den Planner); dazu für die Closure die Ausgänge der drei Risiken (Abschnitt 1, *Übrige DoD-Punkte*) und die Klassen oben. Die Entscheidungen des Architect zu F-623 und F-630 tragen; die Nehmer `slice-v1-abschluss-sqlite-format-doku` und `slice-v1-abschluss-zeitangaben` nehmen an.

**Aufgeräumt:** Container `ver-sl-pg`, `ver-sl-rec`, `ver-sl-rep`, `ver-sl-lock` und die Läufe `pgr-it-*`; Netz `ver-sl-net`; Images `ver-sl-rt`, `ver-sl-deps`; die Kopien und Verzeichnisse mit dem Präfix `ver-sl-` im Scratchpad. `docker ps -a`, `docker network ls`, `docker volume ls` und `docker images` mit `ver-sl` und `pgr-it` danach leer; die Images des Arbeitsbaums (`pgwire-recorder:*`) blieben stehen.

---

## 6. Nachtrag: Nachprüfung der Nacharbeit zu V-175 und V-176 (Stand `1d787cb`)

Geprüft in frischer Kopie (`git archive HEAD`), Docker-only, Scratch-Präfix `ver-sn-`. Die Nacharbeit besteht aus `79a821a` (Code und Test) und `1d787cb` (nur Plan).

| Befund | Prüfung | Urteil |
|---|---|---|
| V-175 fehlende Spalte `interaction.sql` | Aus einer Datei von `record --format sqlite` (zwei einfache Anfragen, echter PostgreSQL 17) entfernte ich jede der 17 Spalten einzeln (bei Schlüsselspalten durch Neuaufbau der Tabelle). `replay` und `play`: alle 34 Läufe Exit 3, `PGR-E3003`, Pfad der Datei in der Meldung, kein Wert aus der Datei; bei `interaction.sql` „no such column: sql“ statt des Textes `sql`. Der Leser nennt die Bezeichner jetzt in eckigen Klammern (`[sql]`). | behoben |
| Test aus dem richtigen Grund rot | Neuer Test `TestLesenJedeFehlendeSpalteUndTabelle` (Datei nur aus einfachen Anfragen, 17 Spalten, 4 Tabellen, Zähler 17). Mutant Q-sql (nur `sql` wieder in doppelten Anführungszeichen) und Mutant Q-interaction (die ganze Abfrage von `interaction` zurück): beide rot mit „Spalte interaction.sql fehlt: <nil> (2 Sessions)“, `gofmt` sauber, frische Kopie. | bestätigt |
| Schreibpfad, gleiches Muster | `grep` über `sqlite.go` und `sqlite_zeilen.go`: eine Abfrage mit doppelten Anführungszeichen, `SELECT "id" FROM "session"` (`pruefeSessions`); fehlt die Spalte, scheitert `Scan` der Zeichenkette in eine Zahl und es folgt `PGR-E3001`, was die Spezifikation für eine Datei mit anderen Nummern verlangt. Die `INSERT`-Anweisungen nennen Spalten in einer Liste, dort ist ein fehlender Bezeichner ein Fehler. Aus dem Quelltext gelesen, den Zweig habe ich nicht gefahren. | kein Fehlverhalten gefunden |
| V-176 Belege im Plan | Gate-Beleg steht auf `79a821a` (letzter Code-Commit; `1d787cb` ändert nur den Plan); Stand-Absatz nennt „mit `45af81d` und den Nacharbeiten“; §8 nennt `97f3ca6..79a821a`: 2369 Zeilen Go, 888 Code, 1481 Tests, das stimmt mit `git diff --numstat` überein. | behoben |
| Namensabgleich | Alle `Test…`-Namen aus dem Plan gefunden (außer Prosa-Treffer); der neue Test steht in `docs/user/abdeckung-unit.md`. Commit-Messages ab `400e9ec` ohne `SPEC-`/`ARC-`. | bestätigt |

Offen bleiben die Einträge des Planner aus V-177 und V-178 (Entscheidung zu §4, Register, Ausgänge der Risiken). V-175 ist damit nicht mehr closure-blockierend. Ob `make gates` auf dem Stand dieses Nachtrags grün ist, steht in der Übergabe-Nachricht.

**Aufgeräumt:** Container `ver-sn-pg`, `ver-sn-rec`, Netz `ver-sn-net`, Images `ver-sn-rt`, `ver-sn-deps` und die Kopien im Scratchpad; `docker ps -a` und `docker images` mit `ver-sn` und `pgr-it` danach leer.
