# Review-Report: slice-v1-abschluss-zeitangaben-pruefen — 2026-10-11

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3, §6 und §8, gegen die Spezifikation `LH-FA-21.a` (*Prüfen*), `SPEC-043` und `SPEC-001`, gegen [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md) und gegen `AGENTS.md` §3.

**Gegenstand:** `8e430ce..2a316b1` (vier Commits: Code und Tests `78a6b4c`, Plan `c3ba98f`, Fixture-Datei `341b569`, Belege `2a316b1`). Geändert: `internal/adapters/driven/recording/yaml.go`, `yaml_test.go`, `docs/user/abdeckung-*.md`, Plan; neu `zeitangaben_pruefen_test.go`, `test/integration/zeitangaben_pruefen_e2e_test.go`, `test/integration/testdata/zeitangaben-absteigend.sqlite`.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-v1-abschluss-zeitangaben-pruefen.md` ganz gelesen am Stand `2a316b1`, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen.
- `spec/spezifikation.md` `LH-FA-21.a` (*Prüfen*), `SPEC-001` (Version und Formatkennung), `SPEC-043`.
- Nehmer: `slice-v1-abschluss-zeitangaben` (§1, DoD, Start-Bedingung) und die Geschwister im Schnitt; der Diff ändert sie nicht.
- `AGENTS.md` §3.3 bis §3.13.
- Vorherige Reports am selben Modul: `slice-v1-abschluss-sqlite-format` und `-lesen` (Nummern dieses Laufs ab F-637).

**Ausgeführte Läufe** (Docker-only, Scratch unter dem Sitzungsverzeichnis, nichts im Repo geändert):

- **Arbeitsweise der Mutanten.** Frische Kopie (`cp -r` ohne `-p`: `go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`), genau eine Ersetzung in `yaml.go`, `go test` des Pakets `recording` im Image `pgwire-recorder:test`.
  - **Rot:** M9 `Style` ungeprüft (`TestZeitangabenFormYAML`, `TestZeitangabenVersionVorForm`); M13 `!!null` angenommen (`TestZeitangabenFormYAML`); M16 Session und Interaktion vertauscht (`TestZeitangabenFormYAML`, `TestZeitangabenFormSQLite`); führende Null erlaubt (`s == "0"`-Zweig entfernt: `TestZeitangabenGueltig`, `TestZeitangabenFormYAML`); Bereich von `int64` ungeprüft (`TestZeitangabenFormYAML`); Vorgänger nur der erste Wert (`TestZeitangabenAbsteigendBeschaedigt`); Fehler der Bibliothek zuerst (`TestZeitangabenFormYAML`, `TestZeitangabenVersionVorForm`); Befund nur für die erste Session (`TestZeitangabenFormYAML`); Befund als Ursache umhüllt (`TestZeitangabenFormSQLite`); Befund vor Formatkennung und Version (`TestZeitangabenVersionVorForm`).
  - **Grün:** keiner.
- **Fixture.** `zeitangaben-absteigend.sqlite` per `sqlite3` gelesen: eine Session, zwei Interaktionen, `offset_ms` 4711 und 4709 (Typ `integer`). Der Mutant „Vergleich entfernt“ lässt `replay` hängen (E1 in §7); die Datei ist sonst gültig und wird aus dem richtigen Grund abgelehnt, nicht wegen einer anderen Prüfung.
- **Reihenfolge der Fehler.** Testprogramm in der Kopie: Datei mit `version: abc` und `offset_ms: 1.5` liefert `PGR-E3002 Version 0 wird nicht unterstützt`; mit `bogus: 1` statt `version: abc` `PGR-E3003` an der Form (F-637).
- **Commit-Messages** `8e430ce..HEAD`: keine `SPEC-`/`ARC-`-Kennung. `//nolint`: keines. Handbuch und `README.md` unangetastet.
- **Gates.** `make gates` auf dem sauberen Baum nach dem Commit dieses Reports: Ergebnis im Verdikt.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-637 | MEDIUM | Der Diff entscheidet im Code die Rangfolge mehrerer Fehler zwischen Version, Form und Bibliothek (Vorprüfung vor dem Dekodieren, Befund erst nach Formatkennung und Version, Fehler der Bibliothek tritt zurück). §6 nennt sie nicht; §6 sagt nur „jeder ist `PGR-E3003`“, die Entscheidung steht als „Hinweis an den Architect (keine Randform)“ in §7 des Code-Commits. Die Folge reicht über das Genannte hinaus: Eine Datei mit `version: abc` und einem unzulässigen `offset_ms` meldet jetzt `PGR-E3002` mit „Version 0“ (einen Wert, der nicht in der Datei steht), ohne den zweiten Fehler `PGR-E3003`; §7 behauptet „beide `PGR-E3003`“. Die Spezifikation lässt die Rangfolge offen („nicht zugesagt“), die Tests `TestZeitangabenVersionVorForm` und die Zeilen M19/M20 in §7 führen sie als Zusage. | `AGENTS.md` §3.12; `LH-FA-21.a` *Prüfen* | `internal/adapters/driven/recording/yaml.go` · `err != nil && vor == nil`; Plan · `beide `PGR-E3003``; `internal/adapters/driven/recording/zeitangaben_pruefen_test.go` · `func TestZeitangabenVersionVorForm` | ja — Testprogramm: `version: abc` mit `offset_ms: 1.5` ergibt `PGR-E3002` | Randform im Code entschieden, §6 nennt sie nicht |
| F-638 | LOW | §3 des Plans nennt als Änderungen in `yaml.go` nur `formVorpruefung` und `sessionFromDTO`; geändert sind auch `Unmarshal` (Reihenfolge, die Quelle von F-637), `geprueftFromDTO` (die Prüfung `offset_ms negativ` entfällt) und neu `offsetForm`, `kanonischeDezimalzahl`. §1 lässt „Typ und Vorzeichen … liefern die Leser schon“ stehen, obwohl das Vorzeichen jetzt die Form in `formVorpruefung` prüft. | `AGENTS.md` §3.9; `BEO-REPO/plan-folgt-korrektur-nicht` | Plan · `Form des Werts am YAML-Baum, Vergleich der `offset_ms` je Session in `sessionFromDTO``; Plan · `negativ in YAML) liefern die Leser schon` | nein — Urteil, `make kopf-check` fängt nur den Kopf | Plan folgt der Korrektur nicht |
| F-639 | INFO | Die Datei `zeitangaben-absteigend.sqlite` (32 KB, binär) ist eingebettet, ohne dass im Repo steht, wie sie entsteht (kein Generator, keine Zeile in Plan oder Test). Wird `schema.sql` später geändert, lehnt der Test sie aus einem anderen Grund ab oder nicht mehr; angezeigt wird das durch die Prüfung von „Session 1, Interaktion 2“ nur teilweise. Die Begründung (Import von SQLite unter `test/integration` verbietet `make a-check`) trägt; zuständig: Architect. | Maintainability | `test/integration/testdata/zeitangaben-absteigend.sqlite`; `test/integration/zeitangaben_pruefen_e2e_test.go` · `//go:embed testdata/zeitangaben-absteigend.sqlite` | nein | Binäre Fixture ohne Herkunft |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `yaml.go` `offsetForm`, `kanonischeDezimalzahl`: `0` oder Ziffern ohne führende Null, kein Vorzeichen, Bereich von `int64`, `Style != 0` deckt Anführungszeichen und ausdrücklichen Typ; die entfernte Prüfung `Kind != ScalarNode` ist äquivalent (ein Nicht-Skalar hat den Wert `""`, den `kanonischeDezimalzahl` ablehnt; Alias und Anker sind vorher abgelehnt) | geprüft, ohne Befund; Mutanten (Style, führende Null, Bereich, `!!null`) rot |
| `yaml.go` `sessionFromDTO`: Vergleich gegen den vorigen vorhandenen Wert, je Session, `<` (gleich gültig), fehlender Wert lässt den Vorgänger stehen; Meldung mit Session und Interaktion, ohne Wert | geprüft, ohne Befund; Mutanten (Vorgänger nur der erste) rot, der Rest der Tabelle in §7 plausibel gegen Tests gelesen |
| Meldungen der Form und des Vergleichs nennen Session und Interaktion, nie einen Wert und keine Zeile, in `yaml` und `sqlite` (SQLite läuft über `Unmarshal`) | geprüft, ohne Befund; Test je Schreibweise und je Typ der Spalte |
| `TestUnmarshalFehlerReihenfolge`: Anpassung zulässig (der Test nennt sich selbst keine Zusage und verlangt die Anpassung bei bewusster Änderung); neue Fälle „offset_ms vor der Nummer“, „Nummer vor Absteigen“, „Form vor Absteigen“ decken die neue Reihenfolge; der entfernte Fall `offset_ms negativ` ist durch die Form ersetzt | geprüft, ohne Befund (Reihenfolge selbst: F-637) |
| `zeitangaben_pruefen_test.go`, `zeitangaben_pruefen_e2e_test.go`: Kommentare und Abdeckungs-Deklarationen sagen nur zu, was geprüft wird; `--keep-timing` nicht als geprüft geführt (Tests rufen `replay` und `play` ohne die Option, Plan §7 sagt es) | geprüft, ohne Befund |
| Fixture: Werte 4711 und 4709, sonst gültig, aus dem richtigen Grund abgelehnt; kein `docs-check`- oder Größenproblem im Diff erkennbar (Gates im Verdikt) | geprüft, ohne Befund außer F-639 |
| §3.13: Nehmer `slice-v1-abschluss-zeitangaben` nimmt die Prüfung ausdrücklich nicht an (§1 *Ausdrücklich NICHT*) und setzt `-pruefen` in der Start-Bedingung voraus; der Diff trägt keine neue Sendung ein; Nachzählung in §8: zwei Liefer-Punkte, eine Schicht | geprüft, ohne Befund |
| Core-Reinheit und a-check: Änderung nur im Recording-Adapter; `strconv` und `errors` aus der Standardbibliothek; `internal/hexagon/**` unberührt | geprüft, ohne Befund |
| Hard Rules: kein `//nolint`, keine `SPEC-`/`ARC-`-Kennung in den Commit-Messages, Handbuch und `README.md` unberührt, Abdeckungstabellen mit den Deklarationen erzeugt, kein Move im Diff (§3.3) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Randform im Code entschieden, §6 nennt sie nicht · Plan folgt der Korrektur nicht · Binäre Fixture ohne Herkunft

## Verdikt

**Merge-blockierend:** nein (kein HIGH). Die Prüfung der Form und des Vergleichs ist durch die nachgefahrenen Mutanten gefangen, in beiden Formaten, ohne grünen Mutanten.

**Übergabe:**

- **F-637** an den **Architect**: die Rangfolge von Version, Form und Bibliothek entscheiden und in §6 und in der Spezifikation festhalten oder die Tests enger fassen; danach der **Implementer** (`AGENTS.md` §3.12, §3.11).
- **F-638** an den **Implementer** (Plan §3 und §1 nachziehen, `AGENTS.md` §3.9).
- **F-639** ohne Blockade an den Architect.
- Die **Finding-Klassen** gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg; DoD- und Spec-Konformität prüft der Verifier.
