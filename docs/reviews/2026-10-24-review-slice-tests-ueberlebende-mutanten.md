# Review-Report: slice-tests-ueberlebende-mutanten — 2026-10-24

**Review-Art:** Code — gegen Plan, Spezifikation, Hard Rules (Maintainability); nicht gegen die DoD.

**Gegenstand:** `slice-tests-ueberlebende-mutanten`, Diff `83eb9993..772c62b3` (Basis des Slice bis Kopf von main); Spec-Entscheidungen des Architect in `e3b904c0`, Implementer-Commits `73aa388f` (Tests, Kommentar) und `144c15e7` (Abdeckungszeilen), dazu die Plan-Commits `0f167533`, `58e32c7f`, `772c62b3`.

**Skill:** `.harness/skills/reviewer.md` @ 772c62b3 ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-24

> **Zitier-Form:** Kennungen statt Adressen; Baseline als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice).

**Eingangs-Kontext:**

- Slice-Plan `slice-tests-ueberlebende-mutanten` (§1, §2, §3, §6, §7 *Belege des Implementers*, §8 Zählung und Planner-Entscheidung)
- `LH-FA-02.a`, `LH-FA-18.a`, `LH-FA-22.a`, `SPEC-038`, `SPEC-049` Punkt 7 (Spezifikation, Stand `772c62b3`); `LH-FA-02`, `LH-FA-07`, `LH-FA-10`, `LH-FA-18`, `LH-FA-22` (Lastenheft)
- ADR-0034 (Black-Box-Tests über die exportierte Schnittstelle)
- `AGENTS.md` §3.9 bis §3.14 (Hard Rules)

**Nicht erbracht:** Die Mutationen habe ich nicht nachgefahren. Auf dem Host liegt kein Go (Docker-only, `AGENTS.md` §3.1), und auf main lief ein voller `make gates`; Docker-Läufe waren ausgeschlossen. Die Mutationstabelle aus §7 habe ich stattdessen gegen den Test-Code und den Produkt-Code gelesen und je Zeile nachgerechnet (siehe Negativbefunde). Das ist Lesen, kein Lauf; das Nachfahren bleibt der Verifikation (§5 des Plans).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Kopf des Plans nennt unter `Bezug` weder `LH-FA-22` (Y2 hält `LH-FA-22.a`, die Spezifikation dort ist geändert) noch `LH-FA-07` (der Test `TestUnmarshalFehlerNenntSession` ist dagegen deklariert). `Berührte Spec-Stellen` führt `LH-FA-22.a`; `make kopf-check` prüft nur §1 und §2, nicht `Bezug`. | `AGENTS.md` §3.9 (Kopf: `Bezug`, `Berührte Spec-Stellen`) | `docs/plan/planning/in-progress/slice-tests-ueberlebende-mutanten.md` · „**Bezug:** [`LH-FA-02`]" | ja — `make doc-trace` zeigt Y2 nicht unter `LH-FA-22` und `LH-FA-07` | Kopf des Plans führt Bezug zum Diff nicht nach |
| F-2 | LOW | Der Doc-Kommentar von `extendedNachspielen` sagt „Sind n Bestätigungen einer Art aufgezeichnet, sind die ersten n Nachrichten dieser Art angenommen“ für jede Art (`parse`, `bind`, `close`). Der neue Test belegt es für `parse` auf ein benanntes Statement mit n = 1; der Plan (DoD 1) verlangt „nur so weit, wie der Test sie prüft“. Ein Szenario, in dem eine Zusage bricht, ohne dass ein Test rot wird, finde ich nicht: Die Zeile `bestaetigt[m.Type]--` ist für alle Arten dieselbe, und „erste n statt letzte n“ trennt der Test (erstes `parse` angenommen, zweites abgelehnt). | `AGENTS.md` §3.11 (Ein Kommentar sagt nur zu, was ein Test prüft) | `internal/hexagon/services/replay.go` · „Sind n Bestätigungen einer Art aufgezeichnet, sind die ersten n Nachrichten dieser Art angenommen;“ | nein — kein Gate prüft Kommentarreichweite | Kommentar reicht weiter als der Test |
| F-3 | LOW | §8 des Plans zählt die fünf von `make abdeckung` erzeugten Zeilen in `docs/user/abdeckung-unit.md` nicht als Schicht (Entscheidung des Planners vom 2026-10-24) und räumt selbst ein, dass der Wortlaut von §3.13 die Gegenlesart ebenso trägt („formal 3“). `AGENTS.md` §3.13 nennt `docs/user/` als Schicht und nimmt Erzeugtes nicht aus; die Entscheidung liegt nicht beim Nutzer, der sonst Schichtbegriffe festlegt (so die Einträge in §3.13). Der Plan hat die Abweichung offen gelegt und dem Nutzer einen Vorschlag zur Klärung gemacht. Ein Szenario, in dem der Slice dadurch nicht in einer Sitzung prüfbar wäre, besteht nicht (191 Zeilen, Marke 1500). | `AGENTS.md` §3.13 *Schichten zählen alle Pläne gleich* | `docs/plan/planning/in-progress/slice-tests-ueberlebende-mutanten.md` · „**Sie zählt nicht.** Entscheidung des Planners über den" | nein — Urteil über Wortlaut und Zweck | Schichtzählung weicht per Planner-Entscheidung vom Wortlaut ab |
| F-4 | INFO | §3 des Plans nennt vier Abdeckungs-Deklarationen (`LH-FA-10`/Negative, `LH-FA-02`/Happy, `LH-FA-18`/Happy, `LH-FA-07`/Negative) und „fünf Tabellenzeilen“. Der M3-Test deklariert zusätzlich `LH-FA-18`/Negative; das ergibt die fünfte Zeile, steht aber in §3 nicht. Die Zeilen in `abdeckung-unit.md` stimmen mit den Deklarationen der Tests überein. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-tests-ueberlebende-mutanten.md` · „(`LH-FA-10`/Negative, `LH-FA-02`/Happy, `LH-FA-18`/Happy, `LH-FA-07`/Negative)" | ja — `make abdeckung-check` (grün, kennt die Plan-Zeile aber nicht) | Plan zählt Deklarationen und Zeilen verschieden |
| F-5 | INFO | Der Test M3 prüft neben `SELECT 1` im Text auch, dass `SELECT 2` fehlt. Das ist nach `SPEC-038` kein ganzer Text, aber `LH-FA-18.a` sagt nur, dass die Diagnose das bestehende SQL nennt; „nennt das abgelehnte nicht“ folgt aus „wie es nach den Nachrichten besteht“, steht aber nicht wörtlich dort. Der Plan (§6 *Was ein Test festhält*) verlangt genau diesen Vergleich, und der Mutant wird schon von der ersten Bedingung rot. Keine Aktion erwartet. | `SPEC-038` *Meldungstext in Tests* | `internal/hexagon/services/replay_extended_test.go` · `strings.Contains(err.Error(), "SELECT 1") \|\| strings.Contains(err.Error(), "SELECT 2")` | nein | — |
| F-6 | INFO | `TestUnmarshalFehlerNenntSession` vergleicht nur `Session 2`; `LH-FA-22.a` sagt Session **und** Interaktion. Die Interaktionsnummer in beiden Meldungen halten Bestandstests (`yaml_test.go` „Session 1: Interaktion 2 trägt die Nummer 3“; `sqlite_lesen_inhalt_test.go` „Session 2, Interaktion 1“), ein Mutant nur an dieser Zahl wird dort rot. Der Plan (DoD 3) verlangt nur die Session. Der Mutant `si+1` statt `sd.ID` ist gleichwertig (`sessionFromDTO` prüft `sd.ID != si+1` vor der Schleife); §7 führt ihn nicht als Mutant mit Sonde und Lauf, der Plan nennt ihn als Grenze, daher kein Befund nach `AGENTS.md` §3.10. | `AGENTS.md` §3.10 (Aufzählung je Element) | `internal/adapters/driven/recording/yaml_test.go` · `strings.Contains(err.Error(), "Session 2")` | nein | — |
| F-7 | INFO | Ich habe die Mutationen nicht gefahren (kein Go am Host, kein Docker-Lauf neben dem Gate-Lauf); Beleg der Rotfärbung ist damit nur §7 des Plans, gelesen gegen den Code. Die Verifikation fährt sie nach (§5 des Plans). Zuständig: Verifier. | Maintainability | Report, Abschnitt „Nicht erbracht“ | ja — Nachfahren je Zeile der Tabelle in §7, frische Kopie, `go test` im Image | Rotfärbung nur gelesen, nicht gefahren |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services` (`replay_extended_test.go`, `replay.go` nur Kommentar) | geprüft; Test M3 gegen `extendedNachspielen` gelesen: Ohne `bestaetigt[m.Type]--` bleibt die Zahl 1 für beide `parse`, das zweite setzt `s1` auf `SELECT 2`, die Diagnose nennt es für das `bind` mit anderem Portal, und `Contains("SELECT 1")` schlägt fehl. Mit der Zeile steht `SELECT 1`. Die Reihenfolge der Gruppen und der Cursor (Interaktion 2, Nachricht 1) stimmen. Der Test vergleicht `PGR-E5001` und Teilstrings, nicht den ganzen Text (`SPEC-038`); die Grenze zum unbenannten Statement ist nicht behauptet. Der Kommentar-Satz „ein weiteres parse auf ein benanntes Statement ändert dieses nicht“ ist vom Test getragen. Befund nur F-2, F-5 |
| `internal/adapters/driven/postgres` (`upstream_test.go`, `upstream_extended_test.go`) | geprüft, ohne Befund: `spalten` ist die einzige Übersetzung (`upstream.go`, Aufruf für Simple und Extended); beide Tests setzen `TableOID` 16384/16999 und `TableAttributeNumber` 3/7, verschieden und ungleich 0, und vergleichen die Spalte mit `reflect.DeepEqual` in allen Feldern. Jede der Mutationen `TableOID:` und `ColumnNumber:` färbt beide Tests rot (0 statt Wert); die Spalte `c` ohne Tabelle belegt die Zusage „0“ der Spezifikation. Pakete `postgres_test`, nur exportierte Schnittstelle (SPEC-049 Punkt 7), kein `//nolint` |
| `internal/adapters/driven/recording` (`yaml_test.go`) | geprüft, ohne Befund; die Kopf-Zeichenkette hängt die zweite Interaktion in Session 2 an (Nummer 3 bzw. Form ohne `ready_for_query`), Session 1 ist gültig. Der Fall „unbekannte Art“ in `TestLesenMeldungOhneWertUndZeile` liegt in `sqlite_lesen_inhalt_test.go` Zeile 409 und geht durch `geprueftFromDTO` mit der Meldung `Session %d, Interaktion %d`; die Berichtigung in §1 stimmt damit (färbt `geprueftFromDTO(1, …)` und die Konstante in der Meldung der Form, nicht die Konstante in „trägt die Nummer“). Die zwei Fälle des Tests sind genau die zwei Meldungen von `geprueftFromDTO`; „negatives `offset_ms`“ ist zu Recht entfallen (`formVorpruefung`, `sessionFromDTO`) |
| `spec/spezifikation.md` (Diff des Architect `e3b904c0`) | geprüft, ohne Befund: keine ADR-, Slice-, Wellen- oder Commit-Kennung in den neuen Absätzen (`AGENTS.md` §3.4); `SPEC-029` und die Verweise „siehe unten“ lösen auf; die Regeln der Meldungen stimmen mit `yaml.go` überein; Historie-Zeile ergänzt |
| `docs/user/abdeckung-unit.md` | geprüft, ohne Befund: genau fünf neue Zeilen, Text, Pfad und Testname entsprechen den Deklarationen in den Testdateien; nur diese Datei ist geändert (kein Handbuch, keine README); `make abdeckung-check` und `make docs-check` laut §7 grün |
| `docs/plan/planning/` (Plan, Verweis in `slice-lint-bestand-kern-driven`) | geprüft; Befunde F-1, F-3, F-4. Y2-Berichtigung in §1 (Gegenstand 4, Aktualität, DoD) stimmt mit Code und Tests überein; §7 enthält keine Zahl, die der Diff widerlegt (Marke: mein Befehl nach `AGENTS.md` §3.14 ergibt 191). Kein Move und keine Inhaltsänderung in einem Commit (`AGENTS.md` §3.3). Closure-Abschnitte (Lerneintrag, Risiko-Ausgänge) stehen noch leer; sie sind vor `done/` fällig, nicht Gegenstand dieses Diffs |
| Suppression, Core-Reinheit, ADRs | geprüft, ohne Befund: keine `//nolint`, kein neuer Import in `internal/hexagon/**`, keine ADR berührt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Kopf des Plans führt Bezug zum Diff nicht nach · Kommentar reicht weiter als der Test · Schichtzählung weicht per Planner-Entscheidung vom Wortlaut ab · Plan zählt Deklarationen und Zeilen verschieden · Rotfärbung nur gelesen, nicht gefahren

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. F-1 und F-2 sind klein im Plan beziehungsweise im Kommentar zu beheben; F-3 entscheidet der Nutzer bei der Closure (der Plan legt den Vorschlag zu `AGENTS.md` §3.13 vor).

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation; die Mutationen nachzufahren ist Sache des Verifiers.
