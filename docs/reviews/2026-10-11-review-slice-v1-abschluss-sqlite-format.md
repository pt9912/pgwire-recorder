# Review-Report: slice-v1-abschluss-sqlite-format — 2026-10-11

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3, §6 und §8, gegen die Spezifikation `LH-FA-22.a` (*Wahl und Erkennung*, *Schreiben von `sqlite`*), `SPEC-043`, `SPEC-044`, `SPEC-049` und `LH-FA-07.a`, gegen [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0013](../plan/adr/0013-sqlite-recording-backend.md), [ADR-0024](../plan/adr/0024-sqlite-bibliothek.md), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) und gegen `AGENTS.md` §3.

**Gegenstand:** `4914a4a..f26ebb9` (sieben Commits, +3072 −32 über 34 Dateien; Code und Tests `f052e63`, `fc59441`, `07d34c4`, `9ebab6c`, Belege des Implementers und Abdeckung `10c10ef`, Randformen des Architect `75874ba`, Tests dazu `f26ebb9`). Neu: `internal/adapters/driven/recording/` (`sqlite.go`, `sqlite_zeilen.go`, `sqlite_json.go`, `format.go`, `schema.sql`, fünf Testdateien), `internal/bootstrap/format_test.go`, `test/integration/sqlite_format_e2e_test.go`. Geändert: `export_test.go`, `cli.go` und Tests, `bootstrap.go`, Kommentare an Port und `RecordService`, `go.mod`, `go.sum`, `.golangci.yml`, `tools/harness/lint-gegenprobe.sh`, `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk`, `harness/README.md`, `tools/schema/schema.yaml`, `spec/spezifikation.md`, `docs/user/abdeckung-*.md`, Plan.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-v1-abschluss-sqlite-format.md` ganz gelesen am Stand `f26ebb9`, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen.
- `spec/spezifikation.md` `LH-FA-22.a` (*Schreiben von `sqlite`*, einschließlich der Änderungen in `75874ba`), `SPEC-043`; `LH-FA-07.a`.
- Nehmer: `slice-v1-abschluss-sqlite-format-lesen` (gelesen auf die Sendungen, §1, DoD, §3, §6; der Diff ändert ihn nicht).
- `AGENTS.md` §3.3 bis §3.13.
- Vorherige Reports am selben Modul: keiner zum Recording-Adapter; die Nummern dieses Laufs beginnen bei F-614.

**Ausgeführte Läufe** (alles Docker-only, Scratch unter dem Sitzungsverzeichnis, nichts im Repo geändert):

- **Arbeitsweise der Mutanten.** Frische Kopie (`cp -r` ohne `-p`: `go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`), genau eine Ersetzung, `go test -count=1` des Pakets `recording` im gepinnten Go-Image (`golang:1.27-alpine@sha256:8a5910…`).
  - **Rot:** `foreign_keys(OFF)` und `foreign_keys` entfernt (`TestSQLiteEinstellungen`, `TestSQLiteFremdschluessel`); Sessions in umgekehrter Reihenfolge ergänzt (`TestSQLiteHalbeSessionNie`, `TestSQLiteTransaktionJeSession`); Prüfung „weniger Sessions als die Datei“ entfernt (W33, `TestSQLiteWenigerSessionsAlsDieDatei`); `INSERT OR REPLACE` statt `INSERT` bei `session` (nur `TestSQLiteDoppelteSessionNummer`).
  - **Grün:** `journal_mode(DELETE)` aus der URI entfernt, `synchronous(FULL)` entfernt (beide F-614); `db.SetMaxOpenConns(1)` entfernt (F-615); Leerzeichen roh statt `%20` (W14, wie von §7 geführt, F-615).
  - **Besonderheit N1b:** Primärschlüssel von `session` aus `schema.sql` entfernt macht 33 Tests rot, nicht nur den benannten (SQLite lehnt den Fremdschlüssel auf eine Spalte ohne Eindeutigkeit ab); F-620.
- **Randform „doppelte Nummer in einem Aufruf“.** Testprogramm in der Kopie: `Write` mit den Nummern 1, 2, 2 auf eine neue Datei liefert `PGR-E3001` (`UNIQUE constraint failed: session.id`), die Datei trägt 1 und 2; derselbe Aufruf wiederholt liefert wieder `PGR-E3001`, die Datei bleibt bei 1 und 2 (F-617).
- **Zählung.** `git diff --numstat 4914a4a..9ebab6c`: 2794 Zeilen hinzu, 676 Go-Code, 1975 Go-Tests (stimmt mit §7); Commit-Messages: keine `SPEC-`/`ARC-`-Kennung.
- **Gates.** `make gates` auf dem sauberen Baum nach dem Commit dieses Reports: Ergebnis im Verdikt.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-614 | LOW | Die Zusagen `journal_mode=DELETE` und `synchronous=FULL` stehen als Pragmas in der URI; entfernt man sie, bleibt jeder Test grün, weil der Standard der Bibliothek auf einer frischen Datei derselbe Wert ist (`TestSQLiteEinstellungen` liest über dieselbe `Oeffne`-URI). §7 führt als grüne Mutanten nur B2, W10, W14; die beiden Entfernungen fehlen dort. Szenario: Ein anderer Prozess ersetzt die Datei durch eine mit passenden Nummern im WAL-Modus (der Modus steht in der Datei); `Write` prüft nur die Nummern, schreibt im WAL-Modus und lässt `-wal`/`-shm` liegen, entgegen „keine WAL-Nebendateien“. Kein Test mit einer Fremddatei im WAL-Modus. | `AGENTS.md` §3.10; `LH-FA-22.a` *Schreiben von `sqlite`* | `internal/adapters/driven/recording/sqlite.go` · `_pragma=journal_mode(DELETE)&_pragma=synchronous(FULL)` | ja — Kopie ohne beide Pragmas: `go test` des Pakets grün | Mutation zu einer Zusage grün, nicht als äquivalent geführt |
| F-615 | LOW | Zwei Kommentare sagen mehr zu, als ein Test prüft. `oeffne`: „Es hält höchstens eine Verbindung“ — der Mutant ohne `db.SetMaxOpenConns(1)` bleibt grün. `dsnVon`: „Jedes Byte außer Buchstaben, Ziffern und -._~/ steht als %XX“ — W14 (Leerzeichen roh) bleibt grün; `DSN` ist für den Test exportiert, geprüft wird der Text der URI aber nicht. | `AGENTS.md` §3.11; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `internal/adapters/driven/recording/sqlite.go` · `Es hält höchstens eine Verbindung und prüft`; `internal/adapters/driven/recording/sqlite.go` · `Jedes Byte außer Buchstaben, Ziffern und -._~/ steht als %XX` | ja — beide Mutanten in frischer Kopie grün | Zusage im Kommentar weiter als Prüfung |
| F-616 | LOW | Der Plan folgt dem Code an drei Stellen nicht. §1 *Übernommen* beschreibt die Randformen des Schreibens noch als offen („entscheidet der Architect … vor dem ersten Code-Commit“, „Dort gilt heute nur: Zielpfad wie in `LH-FA-07.a` …“), obwohl `LH-FA-22.a` sie seit `bb67108`/`75874ba` entscheidet. §8 sagt „das Repo enthält noch keinen Produktionscode“, was dem Repo und `harness/conventions.md` widerspricht (der Satz stand schon vor dem Diff, der Diff lässt ihn als Kontext stehen). | `AGENTS.md` §3.9; `BEO-REPO/plan-folgt-korrektur-nicht` | Plan · `Dort gilt heute nur: Zielpfad wie in`; Plan · `das Repo enthält noch keinen Produktionscode` | nein — Urteil, `make kopf-check` fängt nur den Kopf | Plan folgt der Entscheidung nicht |
| F-617 | LOW | Die Zusage „doppelte Nummer ist `PGR-E3001`, die Datei bleibt auf dem Stand vor dem Aufruf“ (Plan §3 und §7) gilt nur für eine Nummer, die die Datei schon trägt (der Test fährt genau das). Trägt ein einziger Aufruf zwei gleiche neue Nummern, ist die erste committet (Transaktion je Session) und die zweite scheitert; die Datei steht danach nicht auf dem Stand vor dem Aufruf, und jeder weitere `Write` scheitert erneut an derselben Session. Das ist mit `LH-FA-22.a` vereinbar (jede Session eine Transaktion, „die Datei trägt nur die zuvor ergänzten Sessions“), steht aber weder in §6 noch hat es einen Test; der Satz im Plan ist weiter gefasst als der Test. | `AGENTS.md` §3.11, §3.12 (Unterfall einer genannten Randform) | Plan · `lässt die Datei auf dem Stand vor dem Aufruf`; `internal/adapters/driven/recording/sqlite_schreiben_test.go` · `rec.Sessions[1].ID = 1` | ja — Testprogramm: `Write` mit 1, 2, 2 trägt danach 1 und 2 | Zusage im Plan weiter als der Test |
| F-618 | INFO | Der Zweig `runtime.GOOS == "windows"` in `dsnVon` (`"/" + abs`, danach `C%3A`) hat keinen Test und keine Zusage; die Form der URI für Laufwerke und UNC-Pfade ist ungeprüft. §7 führt W10 und die Grenze Plattform (`BEO-REPO/verhalten-nur-unter-linux-geprueft`). Kein Kommentar und keine Spezifikationsstelle sagt Windows zu; Zuständig: Verifier/Planner für die Beobachtung. | Maintainability | `internal/adapters/driven/recording/sqlite.go` · `abs = "/" + abs` | nein | Plattformzweig ohne Prüfung |
| F-619 | INFO | Größe: 2794 Zeilen statt rund 1800 geschätzt (Code 676, Tests 1975, überwiegend Tabellen mit rohem SQL). Die Rückführungsbedingung *nicht in einer Sitzung prüfbar* ist formal nicht ausgelöst; in dieser Sitzung prüfbar war der Diff. Das Maß gehört dem Planner. | Plan §4; `v6.18.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice | Plan · `gemessen 2794 Zeilen hinzu` | nein | Größe über der Schätzung |
| F-620 | INFO | N1b (Primärschlüssel von `session` entfernt) belegt `TestSQLiteDoppelteSessionNummer` nicht einzeln: 33 Tests werden rot, weil das Schema ohne Primärschlüssel den Fremdschlüssel ablehnt. Ein Mutant, der nur die Eindeutigkeit aufhebt (`INSERT OR REPLACE` bei `session`), macht genau diesen Test rot; der Test trägt also, der Beleg ist schwächer als die Zeile in §7. Ob die Zeile den Anspruch „aus dem richtigen Grund“ trägt, entscheidet der Verifier. | `AGENTS.md` §3.10 | Plan · `Primärschlüssel von `session` entfernt (N1b)` | ja — Kopie mit `INSERT OR REPLACE` | Mutant rot aus anderem Grund |
| F-621 | INFO | `LH-FA-22.a` sagt, das Journal einer übrig gebliebenen temporären Datei trage „einen anderen Zufallsteil als jede spätere temporäre Datei“; `neueTempDatei` prüft nur die `.tmp`-Datei auf Belegung, nicht das `-journal`. Die Zusage trägt die Zufallswahl (16 Hexziffern), nicht eine Prüfung; der Test nutzt einen festen Namen. Bei Gleichheit würde SQLite das Journal in die neue Datei zurückrollen. | `AGENTS.md` §3.11 | `spec/spezifikation.md` · `es trägt einen anderen Zufallsteil als jede spätere temporäre Datei` | nein | Zusage durch Wahrscheinlichkeit, nicht Prüfung |
| F-622 | INFO | Die Schichtzählung nennt den Bootstrap Verdrahtung; er trägt jetzt einen Zweig (`if err != nil { return fail(...) }`), der laut §7 (B2) durch kein Programm erreichbar ist. `AGENTS.md` §3.13 nennt „verzweigen“ als Logik; das Weiterreichen eines Fehlers ist kein Einstufen. Grenzfall, Urteil beim Architect; zählt der Bootstrap, liegt der Slice bei drei Schichten. | `AGENTS.md` §3.13 | `internal/bootstrap/bootstrap.go` · `repo, err := recording.New(o.Format)` | nein | Schichtzählung im Grenzfall |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driven/recording/sqlite.go`, `sqlite_zeilen.go`: Transaktion je Session mit `BEGIN IMMEDIATE` (`_txlock`), `mode=rw`, Sperrfrist, `foreign_keys`, Prüfung der Nummern vor dem Ergänzen, Zähler, Schließen nach der Transaktion, `context.WithoutCancel`, `PGR-E3001` für jeden Fehler der Bibliothek, `PGR-E1000` bei weniger Sessions; alle Werte als Parameter | geprüft, ohne Befund außer F-614, F-615, F-617, F-621; Mutanten Reihenfolge, `foreign_keys`, W33 rot |
| `sqlite.go` `Prepare`/`pruefeSQLite`: Pfadprüfung vor Journal, `Stat` für das Journal, nicht prüfbar `PGR-E3001`, auch mit `--force`; Meldungen nennen Pfad der Datei, nie einen Wert aus Konfiguration oder Umgebung | geprüft, ohne Befund |
| `sqlite_json.go`: ein Codec, der Baum aus `yaml.Node`, `base64:` für Werte und Schlüssel, `\`-Escape, keine HTML-Maskierung, `interaction.sql` roh (Zeilen gehen nicht durch den Codec); `Load` als Platzhalter `PGR-E1000` ohne Pfad | geprüft, ohne Befund (Kommentare nennen den Ist-Zustand) |
| `schema.sql`, `tools/schema/schema.yaml`, `meta`, `empty_sessions`, `server_parameters`, `PRAGMA user_version` 0 | geprüft, ohne Befund; Ersatztest `TestSQLiteSchema` statt Abgleich ist Entscheidung des Nutzers (Plan §3) |
| `format.go`, `cli.go`, `bootstrap.go`: Wertemenge `yaml`/`sqlite`, `Formate()` und gemeinsame Zusage der Fabrik, Hilfetext, `replay`/`play` kennen die Option nicht; Meldungen der Option nennen nie den Wert (generische Leser-Tests, `werteJeArt`, Datei-Fall mit `GEHEIM`) | geprüft, ohne Befund außer F-622 |
| `internal/hexagon/**`: nur Kommentare an Port und `RecordService`; kein Import der Bibliothek (Core-Reinheit) | geprüft, ohne Befund |
| `.a-check.yml`, `tools/arch/a-check-negativ.sh` (siebzehn Fälle), `harness/mk/arch-negativ.mk`, `harness/README.md`: Skript, Kopfkommentar, Gate-Zeile und Makefile-Beschreibung sagen dasselbe | geprüft, ohne Befund |
| `.golangci.yml` (`gomodguard_v2.allowed` führt die Bibliothek), `tools/harness/lint-gegenprobe.sh` (Fall `p5-gomodguard-frei-sqlite`); `//nolint`: keine; keine Lockerung | geprüft, ohne Befund |
| `go.mod`, `go.sum`: `modernc.org/sqlite` v1.60.1, `modernc.org/libc` v1.77.1 exakt, `golang.org/x/sync` v0.23.0 indirekt, Rest indirekt | geprüft, ohne Befund |
| `spec/spezifikation.md` (Änderungen in `75874ba`): keine ADR-, Slice- oder Wellen-Kennung, Drift-Log nennt nur Umplanung; Entscheidungen decken sich mit dem Code | geprüft, ohne Befund |
| `docs/user/abdeckung-*.md` (erzeugt); Handbuch und `README.md` unangetastet (Doku-Slice) | geprüft, ohne Befund |
| §3.13: Nehmer `slice-v1-abschluss-sqlite-format-lesen` nimmt die Sendungen mit der Kennung des Gebers in §1 und DoD an (Platzhalter `Load`, URI-Bildung „kein zweiter Weg“, Roundtrip, Lesen); der Diff ändert ihn nicht. Nachzählung im Plan §8 des Gebers: drei Liefer-Punkte, zwei Schichten | geprüft, ohne Befund |
| §3.12: jede Operation des Schreibens hat eine Zeile in §6 oder die Entscheidungen von `75874ba` vor `f26ebb9`; die vier Rückgaben des Implementers kamen vor den Tests, nicht nach Code | geprüft, ohne Befund außer F-617 (Unterfall) |
| §3.3: kein Move im Diff; Commit-Messages ohne `SPEC-`/`ARC-`-Kennung | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Mutation zu einer Zusage grün, nicht als äquivalent geführt · Zusage im Kommentar weiter als Prüfung · Plan folgt der Entscheidung nicht · Zusage im Plan weiter als der Test · Plattformzweig ohne Prüfung · Größe über der Schätzung · Mutant rot aus anderem Grund · Zusage durch Wahrscheinlichkeit, nicht Prüfung · Schichtzählung im Grenzfall

## Verdikt

**Merge-blockierend:** nein. Keine HIGH- und keine MEDIUM-Findings. Der Schreibpfad (Transaktion je Session, Einstellungen, Journal-Rest, Escape, Pfad wörtlich, Nummern, Gate und Erlaubnisliste) ist durch die nachgefahrenen Mutanten gefangen, bis auf die in F-614 und F-615 benannten.

**Übergabe:**

- **F-614, F-615, F-617** an den **Implementer**: Mutant oder Test je Zusage, oder die Zusage enger fassen (`AGENTS.md` §3.11); F-617 zusätzlich Satz im Plan.
- **F-616** an den **Implementer** (Plan §1 und §8 nachziehen, `AGENTS.md` §3.9).
- **F-618 bis F-622** ohne Blockade: F-618 und F-619 an den Planner, F-620 an den Verifier, F-621 und F-622 an den Architect.
- Die **Finding-Klassen** gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg; DoD- und Spec-Konformität prüft der Verifier.
