# Review-Report: slice-v1-abschluss-umgebung-anderes-kommando — 2026-10-12

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3, §6 und §8, gegen die Spezifikation `LH-FA-17.a` (*Einstellungen*, Optionstabelle, *Fehler*, *Anzeige*), `LH-FA-01.a` (Hilfe), `SPEC-034` (Zeile `PGR-E2001`) und gegen `AGENTS.md` §3.2, §3.7 bis §3.13. Die DoD prüft der Verifier.

**Gegenstand:** `2b553d5..4da61d6` (drei Commits: `b91131c` Code, Tests, Abdeckung, Handbuch; `f1203f8` Lint-Nachzug; `4da61d6` Belege in §7). Geändert: `internal/adapters/driving/cli/cli.go`, neu `fremd_test.go`, `cli_test.go`, `leser_test.go`, `play_test.go`, `docs/user/abdeckung-*.md`, `docs/user/benutzerhandbuch.md`, Plan.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `slice-v1-abschluss-umgebung-anderes-kommando` ganz gelesen am Stand `4da61d6`, mit §7 *Belege des Implementers* (Mutationstabelle, Suche, Probe-Tabelle, Läufe). Den Bericht des Implementers habe ich nicht gelesen; die Mutationstabelle habe ich gegen Code und Tests gelesen, nicht nachgefahren.
- `spec/spezifikation.md` `LH-FA-17.a` (Rangfolge, Ausnahmen, *Fehler*, *Anzeige*, Optionstabelle) ganz; Zeile `PGR-E2001` des Katalogs.
- `AGENTS.md` §3 ganz.
- Vorherige Reports am CLI-Adapter: `2026-10-12-review-slice-v1-abschluss-zeitangaben.md` (Kopf gelesen).
- `make docs-check` am Stand `4da61d6`: 632 Dateien, 0 Befunde. `make gates` nicht gefahren (Auftrag).

## Zusammenfassung

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

Wiederkehrende Klasse: keine neue (die beiden LOW sind erstmalig an diesem Slice).

## Findings

### F-1 — LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.11 (Handbuch im Ist-Zustand; sagt nur zu, was ein Test prüft)
- `pfad`: `docs/user/benutzerhandbuch.md` · „und eine Variable mit dem Präfix, die zu keiner Option gehört (`PGWIRE_RECORDER_BOGUS`), bleibt bei jeder Betriebsart unbeachtet"
- `befund`: Das Handbuch führt `--tls-cert`, `--tls-key`, `--allow-plaintext`, `--session-assignment`, `--record-empty-sessions` und `--compare-responses` in keiner Tabelle. Für einen Leser des Handbuchs „gehört" `PGWIRE_RECORDER_TLS_CERT` damit zu keiner Option, das Binary lehnt sie aber bei `play` mit `PGR-E2001` ab (`zuordnung()`, `TestLeserUmgebungNichtGelesenerOptionen`). Der Implementer nennt die Lücke selbst unter *Verbleibende Risiken*; der Satz ist für diese sechs Namen dennoch nicht gedeckt.
- `verifizierbar`: ja (`play` mit `PGWIRE_RECORDER_TLS_CERT=x` gegen das gebaute Binary; die Probe-Tabelle §7 hat keine Zeile dafür)
- `klasse`: Handbuchsatz allgemeiner als die dokumentierte Optionsmenge

### F-2 — LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.9 (Plan folgt dem Diff)
- `pfad`: `docs/plan/planning/in-progress/slice-v1-abschluss-umgebung-anderes-kommando.md` · „`docs/user/abdeckung-unit.md` | erzeugt (`make abdeckung`) | folgt den Deklarationen"
- `befund`: §3 und die DoD nennen nur `abdeckung-unit.md`. Der Diff ändert zusätzlich `abdeckung-gesamt.md` (`LH-FA-01` von „teilweise" auf „vollständig", Boundary „Unit") und `abdeckung-vollstaendig.md` (neue Zeile `LH-FA-01`). §7 nennt die Dateien in *Läufe*, §3 nicht; die Statusänderung von `LH-FA-01` ist durch die neue Deklaration `LH-FA-01/Boundary` in `TestLeserUmgebungConfigShowUndHilfe` entstanden, nicht vorab geplant.
- `verifizierbar`: ja (`make abdeckung-check`, `git diff --stat`)
- `klasse`: Plan §3 nennt erzeugte Folgedateien unvollständig

### F-3 — INFO

- `kategorie`: INFO
- `quelle`: Maintainability / `AGENTS.md` §3.7
- `pfad`: `internal/adapters/driving/cli/cli.go` · „auch die, die der Stand noch nicht liest"
- `befund`: Das Wort „noch nicht" in `zuordnung()` und `liesUmgebung` ist zeitbezogen; der Kommentar bleibt Zusage und Kopplung (jede Option der Tabelle, Unbeachtetes bei eigenen Kommandos), und `TestLeserUmgebungNichtGelesenerOptionen` prüft, was er sagt. Keine Aktion erwartet; Hinweis für die Slices, die den Leser liefern: dort ändert sich der Satz mit.
- `verifizierbar`: nein
- `klasse`: zeitbezogene Wortwahl in Kommentar

### F-4 — INFO

- `kategorie`: INFO
- `quelle`: `LH-FA-17.a` *Fehler* / Plan §6 *Optionen, die der Stand noch nicht liest*
- `pfad`: `internal/adapters/driving/cli/cli.go` · „bleibt die Variable unbeachtet"
- `befund`: Die Spezifikation sagt die Prüfung des Werts jeder Option des Kommandos zu, der Code lässt die Variable der sechs ungelesenen Optionen beim eigenen Kommando unbeachtet. Das ist vom Architect am 2026-10-12 entschieden (§6) und durch Test gehalten; hier nur als bekannte Abweichung vom Wortlaut der Spezifikation festgehalten, Auflösung durch die drei Slices, die den Leser liefern. Zuständig: Architect/Planner.
- `verifizierbar`: ja (`TestLeserUmgebungNichtGelesenerOptionen`)
- `klasse`: bewusste Abweichung Code gegen Spezifikationswortlaut

## Prüfpunkte des Auftrags

1. **Verhalten gegen die Spezifikation.** Rangfolge Kommandozeile (`liesKommandozeile`) → ein Durchlauf in Tabellenreihenfolge (`liesUmgebung` über `zuordnung()`) → Datei (`ladeGewaehlte`) stimmt mit *Fehler* überein; nur der erste Fehler (Rückgabe beim ersten Treffer); leer gilt als nicht gesetzt (`v == ""`); Meldung nennt Namen, nicht Wert; `PGWIRE_RECORDER_PASSWORD` und `PGWIRE_RECORDER_CONFIG` sind keine Zeilen von `zuordnung()`; Präfix ohne Option und andere Schreibweise unbeachtet (`os.Getenv` je Tabellenname, exakt); `config show` ruft `lies` nicht (`parseConfig`); Hilfe vor `lies` (`TestParseHilfeVorPruefung`); sechs ungelesene Optionen: `index` kennt sie nicht, eigenes Kommando `continue`, anderes Kommando Fehler. Kein Befund.
2. **`zuordnung()` gegen die Spalte *Kommando*.** Zeile für Zeile verglichen (26 Zeilen, `--config` ausgenommen, `--log-level` an letzter Stelle): Namen, Kommandos und Reihenfolge stimmen mit der Tabelle in `LH-FA-17.a` überein. Kein Befund. (Der Test hält dieselbe Tabelle hart kodiert, akzeptiertes Negativ nach §6; mit diesem Lauf als Drift-Prüfung gegen die Spezifikation erledigt.)
3. **Mutationen.** Jede Zusage der vorgegebenen Tabelle in §6 hat eine Zeile in §7 mit rotem Test; die Tests sind nach Lesart fähig, die Mutation zu fangen (Werte `true`/`5s`/`GEHEIM`, leere Variable, Namen-Prüfung, beide Reihenfolgen, `EqualFold`-Fall mit `pgwire_recorder_listen`). Zusagen ohne roten Test: keine gefunden. Nicht nachgefahren.
4. **§3.11.** Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (der Anker `#5-einstellungen` ist ein interner Link, `make docs-check` grün); jeder Satz hat eine Zeile der Probe-Tabelle bis auf den Fall aus F-1. Hilfetext „Jede Option ist auch über ihre Umgebungsvariable setzbar" bleibt wahr (er beschreibt die Optionen des jeweiligen Kommandos; `TestLeserHilfe` hält ihn unverändert). Kommentare tragen Zusage/Kopplung/Rang-Zeiger; keiner beschreibt eine verworfene Alternative oder abwesenden Text.
5. **§3.9 und §3.13.** Plan folgt dem Diff bis auf F-2; §8 zählt 2 Liefer-Punkte und 2 Schichten (CLI-Adapter, Nutzer-Doku), Bootstrap ohne Logik (siehe Punkt 7), keine Sendung in einen Nehmer, damit kein Nehmer nachzuzählen. Kein weiterer Befund.
6. **Lint-Profil.** Kein `//nolint` in `internal/`, `cmd/`, `test/` (Suche leer); `.golangci.yml` im Diff unverändert; der Lint-Nachzug `f1203f8` macht die Tabellen zu Funktionen statt Ausnahmen zu schreiben.
7. **Bootstrap.** `internal/bootstrap` ist im Diff nicht berührt; keine Logik dort.

## Negativbefund-Zeilen (geprüft, ohne Befund)

- `internal/adapters/driving/cli/` (`cli.go`, `fremd_test.go`, `cli_test.go`, `leser_test.go`, `play_test.go`): geprüft, ohne Befund außer F-3 und F-4.
- `internal/bootstrap/`: unverändert, geprüft, ohne Befund.
- `internal/hexagon/**`: nicht berührt, Core-Reinheit unverändert, ohne Befund.
- `spec/`: nicht berührt (Abgrenzung §1 eingehalten), ohne Befund.
- `docs/user/benutzerhandbuch.md`: geprüft, Befund F-1.
- `docs/user/abdeckung-*.md`: geprüft, Befund F-2 (Plan), Tabellen selbst ohne Befund.
- `.golangci.yml`, `.a-check.yml`, `harness/`: nicht berührt, ohne Befund.
- ADRs: keine berührt, keine verletzt (`ADR-0001`, `ADR-0034` gelesen), ohne Befund.
- `docs/plan/planning/in-progress/slice-…`: Randformen in §6 vor dem Code entschieden (`AGENTS.md` §3.12), keine im Code-Commit neu eingetragen, ohne Befund.
