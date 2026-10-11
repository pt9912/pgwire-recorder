# Review-Report: slice-v1-abschluss-zeitangaben-aufzeichnen — 2026-10-12

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3, §6 und §8, gegen die Spezifikation `LH-FA-21.a` (*Aufzeichnen*, *Uhr*, *Messpunkt*) und `LH-FA-17.a` (Optionstabelle), gegen [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0015](../plan/adr/0015-uhr-port.md), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md) und gegen `AGENTS.md` §3.3 bis §3.13. Die DoD prüft der Verifier.

**Gegenstand:** `f4b7a97..833960b` (vier Commits: `9ed8d68` Option und Messpunkt, `0bfbfe4` Testhilfe `offsetlese`, `2e23f50` Fälle von `record_timing`, `833960b` Plan und Roadmap). Geändert: `internal/hexagon/services/record.go`, `internal/hexagon/model/recording.go`, `internal/adapters/driving/cli/cli.go`, `internal/bootstrap/bootstrap.go`, `internal/bootstrap/offsetlese`, `internal/bootstrap/testhilfen_import_test.go`, Tests in vier Paketen und `test/integration`, `docs/user/abdeckung-*.md`, Plan, Roadmap.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `slice-v1-abschluss-zeitangaben-aufzeichnen` ganz gelesen am Stand `833960b`, mit §7 *Belege des Implementers* (Mutationstabelle, grüne Mutanten, Gegenbeispiel, Grenzen). Den Bericht des Implementers habe ich nicht gelesen; seine Aussagen zu den Mutanten und zum Gegenbeispiel habe ich nachgefahren und gegen den Code gelesen.
- `spec/spezifikation.md` `LH-FA-21.a` (*Aufzeichnen*, *Uhr*, *Messpunkt*, *Prüfen*), `LH-FA-17.a` (Optionstabelle, Meldungen), Lastenheft [`LH-FA-21`](../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen).
- Nehmer `slice-v1-abschluss-zeitangaben` und `slice-v1-abschluss-zeitangaben-doku`; der Diff ändert sie nicht. `slice-v1-abschluss-sessions` (liefert `record_empty_sessions`).
- Vorherige Reports am selben Modul: `slice-v1-abschluss-zeitangaben-uhr` und `-pruefen` (Nummern dieses Laufs ab F-647).
- [ADR-0015](../plan/adr/0015-uhr-port.md), [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), `AGENTS.md` (Hard Rules).

**Ausgeführte Läufe** (Docker-only, Kopien unter dem Sitzungsverzeichnis ohne `.git` und ohne Zeitstempel, nichts im Repo geändert; Unit-Mutanten mit `go test -run` im Image `pgwire-recorder:test` ohne Netz, gofmt je Mutant geprüft; Binary-Mutanten mit `make test-integration` in der Kopie, Container, Netz und Volume vom Runner entfernt):

- **Unit-Mutanten des Services** (Pakete `internal/hexagon/services`, je eine Ersetzung), alle rot und im Wortlaut der Zusage von `LH-FA-21.a`: Operanden vertauscht (`TestRecordZeitEinfach`, `[-40 -45]`); Beginn nach `upstream.Open` gelesen (`TestRecordZeitBeginnVorAufbau`, `[0]` statt `[30]`); Lesung von `Query` nach der Warteschleife (`TestRecordZeitAnkunftVorWarten`, `Query liest die Uhr nicht, bevor sie auf die Extended-Interaktion wartet`); Extended nur bei `Parse` (sechs Arten rot); Extended erst beim `Sync` (sechs Arten rot, `TestRecordZeitExtendedMehrereGruppen` rot); Lesung bei der ersten Nachricht jeder Gruppe (`[20]` statt `[10]`); `nil` bei Wert 0 (`TestRecordZeitNullSteht`); `nil` für alle `Query` und `nil` für alle Extended; `CloseSession` hängt offene Interaktionen an (`TestRecordZeitUnvollstaendigOhneWert`, `[10 20]` statt `[10]`); `offset` liefert ohne Option einen Wert (`TestRecordOhneZeitangaben`); `lesung` ruft `WarteBis` (`TestRecordZeitWartetNie`).
- **Mutanten der Verdrahtung am Binary** (`make test-integration`): `Zeitangaben(NeueSystemuhr(), true)` fest rot (`TestE2ERecordOhneTiming`: `yaml: offset_ms [3 4 4], erwartet ohne Feld`; `TestE2ERecordTimingQuellen`); `Zeitangaben(NeueSystemuhr(), false)` fest rot (`TestE2ERecordTiming`, `TestE2ERecordTimingQuellen`); konstante Uhr `Jetzt() == 0` rot (`TestE2ERecordTiming`: `[0 0 0]` in beiden Formaten).
- **Testhilfe.** `bootstrap.go` importiert `offsetlese` (leerer Import): `TestBinaryOhneTesthilfen` rot (`das Binary linkt …/internal/bootstrap/offsetlese`). Der Eintrag in der Verbotsliste trägt, `go list -deps ./cmd/...` läuft ohne Tests, die Testhilfe erreicht das Binary nicht.
- **Kommandozeile.** `record --record-timing=GEHEIM` und `record --force=GEHEIM` am gebauten Binary: dieselbe Meldung der Bibliothek (`invalid boolean value "GEHEIM" for -…: erlaubt sind true und false`, `PGR-E2001`, Exit-Code 2); Umgebung und Datei nennen den Wert nicht (`TestE2ERecordTimingUngueltig` prüft das).
- **Gegenbeispiel „Werte steigen nicht ab“.** `internal/adapters/driving/pgwire/server.go`: `Query` und `ClientMessage` werden nur aus `nachricht` der einen Schleife `clientRichtung` gerufen, `Query` blockiert die Schleife bis zur Antwort.
- **Commit-Messages** `f4b7a97..HEAD`: keine `SPEC-`/`ARC-`-Kennung. `//nolint`: keines im Diff. Handbuch, `README.md` und `spec/` unberührt. Kein Move im Diff (§3.3).
- **Gates.** `make gates` auf dem sauberen Baum nach dem Commit dieses Reports: Ergebnis im Verdikt.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-647 | MEDIUM | Die Zusage „Die Werte einer Session steigen nicht ab“ gilt nur, wenn die Aufrufe der Client-Richtung einer Session nacheinander kommen. Ruft ein Aufrufer `ClientMessage` (neue Extended-Interaktion, Lesung `t2`) während einer wartenden `Query` (Lesung `t1 < t2`) auf, bekommt die `Query` die frühere Lesung und die höhere Nummer, und der Recorder schreibt eine Aufzeichnung, die er selbst als `PGR-E3003` ablehnt. Über den PGWire-Adapter tritt das nicht ein (eine Schleife, `Query` blockiert sie, von mir gelesen), aber weder der Port `Recorder` noch ein Test nennen die Voraussetzung: der Port sagt nur, die beiden Richtungen dürften gleichzeitig laufen. §6 verlangt für ein Gegenbeispiel, dass der Implementer anhält und es zurückgibt; es steht als Satz in §7 (*Gegenbeispiel*), ohne Entscheidung des Architect. | `AGENTS.md` §3.12, `AGENTS.md` §3.11, `LH-FA-21.a` *Messpunkt* | `internal/hexagon/services/record.go` · `ankunft := s.lesung()` (Zeilen 186 und 247); `internal/hexagon/ports/driving/record.go` · `Aufrufe der beiden Richtungen dürfen für dieselbe Session`; Plan §6 *Nicht absteigend* und §7 *Gegenbeispiel zu* | nein — ein Test müsste `Query` und `ClientMessage` gleichzeitig rufen, wie es kein Aufrufer des Adapters tut | Zusage hängt an nicht genannter Aufruf-Voraussetzung |
| F-648 | LOW | Die Zusage „Kein Paket unter `./cmd/...` importiert die Testhilfe“ (Kommentar von `doc.go`) hat in §7 keine Zeile *Zusage · Mutation · roter Test*. Der Test trägt sie (der leere Import im Bootstrap macht `TestBinaryOhneTesthilfen` rot, von mir gefahren), die Tabelle führt die Mutation nicht. | `AGENTS.md` §3.10 | `docs/plan/planning/in-progress/slice-v1-abschluss-zeitangaben-aufzeichnen.md` · Tabelle *Zusage · Mutation · roter Test* (keine Zeile zu `offsetlese`); `internal/bootstrap/offsetlese/doc.go` · `Kein Paket unter` | ja — `make test` mit dem Import im Bootstrap | Mutation einer Zusage nicht in §7 geführt |
| F-649 | LOW | Die drei Code-Commits `9ed8d68`, `0bfbfe4`, `2e23f50` lassen §3 des Plans stehen, wie es vor der Option `services.Zeitangaben` und der Testhilfe `offsetlese` war; erst `833960b` zieht §3, §7 und §8 nach. Der Endstand stimmt mit dem Code überein. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-v1-abschluss-zeitangaben-aufzeichnen.md` · Zeile `Tests, test/integration` in §3 (Diff `833960b`) | nein | Plan folgt der Korrektur nicht im selben Commit |
| F-650 | INFO | `TestDateiUngueltig` ist nicht geschwächt: Der Fall „Option, die der Stand nicht kennt“ steht weiter da, mit `record_empty_sessions` statt `record_timing`, dazu vier neue Fälle für `record_timing` (in `replay:`, `play:`, oben, ungültiger Wert; die Fälle der Ebene tragen den Wert `true`, damit nur die Ebene den Fehler macht). Der Platzhalter ruht darauf, dass der Stand `record_empty_sessions` nicht kennt; `slice-v1-abschluss-sessions` liefert diese Option und nennt den Platzhalter in §1 und DoD nicht. Der Test wird dort laut rot, nicht still. Ohne Aktion erwartet, Hinweis an den Planner. | `AGENTS.md` §3.13 | `internal/adapters/driving/cli/datei_test.go` · `{"Option, die der Stand nicht kennt"` | ja — `make test` | Test ruht auf der Nicht-Lieferung eines Folge-Slice |
| F-651 | INFO | `docs/user/abdeckung-gesamt.md` und `abdeckung-vollstaendig.md` führen `LH-FA-21` jetzt als `vollständig`, während `--keep-timing` (Einspielen, `slice-v1-abschluss-zeitangaben`) noch aussteht. Die Zeile ist die mechanische Folge der Deklarationen (Happy, Boundary und Negative belegt, `SPEC-048`); das Werkzeug prüft nicht, ob die Anforderung als ganze belegt ist (*Grenze* dort). Kein Fehler des Diffs, der Wortlaut von `vollständig` ist enger als er klingt. | `AGENTS.md` §3.11 | `docs/user/abdeckung-gesamt.md` · `LH-FA-21` mit Stand `vollständig` | nein | Stand `vollständig` weiter als die Prüfung |
| F-652 | INFO | Zu (e): „Ohne Option byte-gleich mit der Fassung vor dem Slice“ stimmt durch den Code (ohne Option ist `OffsetMS` immer `nil`, die Schreiber sind unverändert) und durch die schon vorhandenen Aufzeichnungstests (`yaml_test.go`, `record_e2e_test.go` führen Aufzeichnungen ohne das Feld, nicht Zeile für Zeile gegen den Vorstand verglichen). `TestE2ERecordOhneTiming` selbst vergleicht zwei Läufe mit verschiedenen Pausen miteinander und prüft, dass `offset_ms` fehlt; einen Vergleich gegen den Stand vor dem Slice führt kein Test, und der Satz „byte-gleich der heutigen“ in der DoD ist so gelesen enger zu fassen als geprüft. Verifier. | `AGENTS.md` §3.11 | `test/integration/zeitangaben_aufzeichnen_e2e_test.go` · `yaml unterscheidet sich zwischen den Läufen` | ja — Mutant `true` fest ist rot | Zusage „byte-gleich“ gegen Vorstand nicht direkt geprüft |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services/record.go`: Beginn vor `upstream.Open`, Lesung als Erstes in `Query` und `ClientMessage`, Extended mit der ersten Nachricht der ersten Gruppe (alle Arten), keine Lesung ohne Option, kein `WarteBis`, unvollständige Interaktion ohne Wert, Lebendprüfung ohne Sonderpfad (kein `istLebendpruefung` in `record.go`); 12 Mutanten rot, die zwei in §7 als grün geführten Fälle stimmen | geprüft, ohne Befund |
| `internal/hexagon/model/recording.go`: Kommentar von `OffsetMS` nennt Zusage in Indikativ, deckt sich mit *Messpunkt*; kein neues Feld | geprüft, ohne Befund |
| `internal/adapters/driving/cli`: `--record-timing[=true\|false]` in der Optionstabelle (Default `false`, nach `--format`), Umgebungsvariable, Schlüssel nur unter `record:`, bei `replay`/`play` unbekannt (`PGR-E2001`, `PGR-E2004`), Hilfetext sagt nur, was `Messpunkt` und Tests tragen; Meldung zur Kommandozeile nennt den Wert wie `--force` (zulässig nach `LH-FA-17.a`), zu Umgebung und Datei nie | geprüft, ohne Befund |
| `internal/bootstrap/bootstrap.go`: eine Zeile `services.Zeitangaben(NeueSystemuhr(), o.RecordTiming)`, keine Verzweigung, damit Verdrahtung und keine Schicht (`AGENTS.md` §3.13); drei Mutanten am Binary rot | geprüft, ohne Befund |
| `internal/bootstrap/offsetlese` und `testhilfen_import_test.go`: Eingriff notwendig (Integrationstests unter `test/` dürfen den Recording-Adapter nicht importieren, `.a-check.yml`), liegt im Bootstrap, wo der Adapter importiert werden darf; Verbotsliste hält die Testhilfe aus dem Binary, Mutant rot (F-648 betrifft nur die Zeile in §7) | geprüft, ohne Befund |
| `internal/adapters/driving/pgwire`: unberührt (Messpunkt im Use Case, zwei Schichten, Rückführung aus §4 nicht eingetreten); Einzel-Schleife der Client-Richtung gelesen für F-647 | geprüft, ohne Befund |
| Zu (d), Flackern von `TestE2ERecordTiming`: geprüft werden nur Untergrenzen (`got[1]-got[0] >= 300`, `got[2]-got[1] >= 200`) und `>= 0`; der Unterschied zweier abgeschnittener Lesungen ist nie kleiner als die Pause (`floor(x+300) = floor(x)+300`), eine belastete CPU macht Pausen länger, nie kürzer, also kein Flackern in diese Richtung; keine Obergrenze, damit keine Zeitabhängigkeit nach oben (Risiko 1 aus §6, Ausgang bleibt für die Closure) | geprüft, ohne Befund |
| `TestE2ERecordTimingQuellen`, `TestE2ERecordTimingUngueltig`, `TestE2ERecordTimingHilfeUndConfigShow`: Quellen-Rangfolge Option vor Umgebung vor Datei, `config show` nennt Schlüssel und Namen der Variable ohne Wert | geprüft, ohne Befund |
| `docs/user/abdeckung-*.md`: Deklarationen sagen zu, was ihr Test prüft (Stichprobe der neuen Zeilen gegen Testkörper); `make abdeckung-check` Teil des Gate-Laufs | geprüft, ohne Befund (F-651 betrifft nur die Wortbedeutung) |
| Plan §1, §3, §6, §7, §8 gegen den Code: zwei Liefer-Punkte, zwei Schichten (Kern, CLI-Adapter), Bootstrap als Verdrahtung; Nehmer `slice-v1-abschluss-zeitangaben` und `-doku` nehmen an (Sendungen stehen dort mit den Kennungen der Geber, der Diff ändert sie nicht, `AGENTS.md` §3.13); keine neue Randform im Code-Commit (§3.12) | geprüft, ohne Befund (F-649 betrifft nur den Zeitpunkt der Nachführung) |
| Hard Rules: kein `//nolint`, keine `SPEC-`/`ARC-`-Kennung in Commit-Messages, Kommentare im Indikativ ohne verworfene Alternative (§3.7), Handbuch und `README.md` unberührt (§3.11), Roadmap-Zeile nennt Zustand ohne Chronik | geprüft, ohne Befund |
| Core-Reinheit: `internal/hexagon/**` importiert nur das Domain Model und den Port `driven.Uhr`, keine Zeitbibliothek, kein Adapter | geprüft, ohne Befund |
| Suppression und Gate-Lockerung: nichts in `.golangci.yml`, `.a-check.yml`, `.d-check.yml` geändert | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusage hängt an nicht genannter Aufruf-Voraussetzung · Mutation einer Zusage nicht in §7 geführt · Plan folgt der Korrektur nicht im selben Commit · Test ruht auf der Nicht-Lieferung eines Folge-Slice · Stand `vollständig` weiter als die Prüfung · Zusage „byte-gleich“ gegen Vorstand nicht direkt geprüft

## Verdikt

**Merge-blockierend:** ja, für F-647 bis zur Entscheidung des Architect; sonst nein (kein HIGH). Die Zusagen von `LH-FA-21.a` *Aufzeichnen* und *Messpunkt* sind in Service, Option und Verdrahtung gefangen: zwölf Unit-Mutanten und drei Mutanten am Binary sind rot, die Testhilfe erreicht das Binary nicht.

**Übergabe:**

- **F-647** an den **Architect**: Er entscheidet, ob die Voraussetzung (Aufrufe der Client-Richtung einer Session nacheinander) eine Zeile im Vertrag des Ports `Recorder` und der Spezifikation (*Messpunkt*) bekommt oder die Zusage enger gefasst wird. Das Gegenbeispiel steht in §7 des Plans; anhalten und zurückgeben verlangt §6.
- **F-648** und **F-649** an den **Implementer**, ohne Blockade: Zeile zur Verbotsliste in §7 nachtragen.
- **F-650** an den **Planner**: Der Platzhalter-Fall von `TestDateiUngueltig` ruht darauf, dass `slice-v1-abschluss-sessions` die Option noch nicht liefert; ob dieser Slice ihn in §1 nennt, ist dessen Planung.
- **F-651** und **F-652** an den **Verifier**, ohne Aktion.
- Die **Finding-Klassen** gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg; DoD- und Spec-Konformität prüft der Verifier.
