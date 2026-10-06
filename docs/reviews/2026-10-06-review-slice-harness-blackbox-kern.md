# Review-Report: slice-harness-blackbox-kern — 2026-10-06

**Review-Art:** Code (Testumbau) und Plan-Nachzug, geprüft gegen Plan, Entscheidungen (`SPEC-049` Punkt 7 und 8, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 und 5) und Hard Rules (Modul 10). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git show 839e61c`. Die Testdateien unter `internal/hexagon/model` (2) und `internal/hexagon/services` (7) sind jetzt `_test`-Pakete. Neu sind `export_test.go` in beiden Paketen. Der Plan `docs/plan/planning/in-progress/slice-harness-blackbox-kern.md` ist in §3, §6 und §7 (Mutationstabelle) nachgezogen. Schwerpunkte laut Auftrag: die Brücken nach `SPEC-049` Punkt 7, verdeckt statt behoben, Prüfung im Umbau verloren, Testliste und Abdeckungs-Deklarationen, Mutationen.

**Skill:** `.harness/skills/reviewer.md` @ `839e61c`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-harness-blackbox-kern.md` (ganz, Stand `839e61c`), dazu der Plan-Diff des Commits
- `spec/spezifikation.md` `SPEC-049` Punkt 1 bis 8; [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) (Accepted), [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)
- `AGENTS.md` §3.2, §3.3, §3.7, §3.9 bis §3.12
- Produkt-Code als Bezug der Brücken: `internal/hexagon/model/fehler.go`, `internal/hexagon/services/replay.go`, `matcher.go`, `lebendpruefung.go`, `record.go`
- Vorherige Reports: Review zu `slice-harness-lint-werkzeug` und zu `slice-harness-upgrade-v6-16`. Die höchste vergebene Nummer vor diesem Lauf war F-440. Die Klasse „Randform im Code entschieden“ war bisher dreimal LOW (Reviews zu `slice-replay-semantik-mismatch`, `slice-harness-kopf-sensor`, `slice-harness-lint-werkzeug`).

**Ausgeführte Läufe:**

- `make lint`: Exit 2. Unter beiden Pfaden gibt es nur Befunde im Produkt-Code (`extended.go`, `recording.go`, `replay.go`), keinen in einer `*_test.go`.
- `make test`: grün. `make abdeckung-check` und `make kopf-check`: beide Exit 0.
- Testliste: Stand `839e61c^` und `839e61c` je per `git archive` in eine frische Kopie unter dem Scratch-Pfad entpackt. Dort lief `go test -list .` für beide Pakete im Image der Stufe `deps`, ohne Netz. Die Listen sind gleich (62 Tests).
- Abdeckungs-Deklarationen: alle `Abdeckung:`-Zeilen unter beiden Pfaden vor und nach dem Commit sortiert verglichen. Sie sind gleich (43 Zeilen).
- Sechs eigene Mutationen, je in einer frischen Kopie (`cp -r` ohne `-p`) der Kopie von `839e61c`, mit `go test -count=1 -run` im Image der Stufe `deps`, ohne Netz. Der Grundlauf war grün; alle sechs Mutationen waren rot (Tabelle unten).
- Das Repo selbst blieb bis auf diese Datei unverändert.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-441 | MEDIUM | Ob eine Brücke einen Empfänger eines unexportierten Typs selbst anlegen darf, beantwortet `SPEC-049` Punkt 7 nicht. Gemeint ist `LetzteNummer`, das einen `cursor` auf der übergebenen Session baut. Der Text sagt: „Zustand setzt sie nur an einem Wert, den sie übergeben bekommt“, der `cursor` wird aber nicht übergeben. Diese offene Randform ist erst im Code-Commit `839e61c` in §6 eingetragen („Brücke mit Cursor“) und dort dem Review zugewiesen („das ist Urteil des Review“). Nach §6 und `AGENTS.md` §3.12 gehört sie dem Architect, und zwar vor dem Code. Failure-Szenario: Die Auslegung steht dann nur in diesem Report, also in einem Lauf-Beleg. Die drei weiteren Umstellungs-Slices und später der Review des Gates in `slice-harness-lint` lesen ihn nicht und entscheiden dieselbe Frage jeder für sich neu. Das ist das vierte Auftreten dieser Klasse, deshalb MEDIUM. | `AGENTS.md` §3.12; Plan §6 *Randformen* („was dort nicht steht, gibt der Implementer an den Architect zurück“) und §4 *Start*; `SPEC-049` Punkt 7; `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“ | `docs/plan/planning/in-progress/slice-harness-blackbox-kern.md` · „Ob das noch Weiterreichen ist, prüft das Werkzeug“; `internal/hexagon/services/export_test.go` · „return (&cursor{session: session}).letzteNummer()“ | nein (Lesen; `git log --follow` auf den Plan zeigt den Eintrag erst in `839e61c`) | Randform im Code entschieden |
| F-442 | INFO | Urteil des Reviewers zur Sache, für den Architect, der nach F-441 entscheidet. Die Brücke hat keine Variable, keine Funktion `Test…` und keinen Zustand auf Paketebene. Der `cursor` ist ein lokaler Wert, er lebt nur für einen Aufruf, und die übergebene Session bleibt unverändert. Das entspricht dem Zweck von Punkt 7. Eine Kopplung bleibt: Die Brücke setzt nur `session`, während der Produkt-Code seine Cursor mit `status: "I"` anlegt. Liest `letzteNummer` später ein weiteres Feld, prüft `TestReplayLetzteNummer` einen Zustand, den das Produkt nie erzeugt. Die Nummer ist auch über die exportierte Diagnose zu sehen („nach Interaktion 2 erwartet die Aufzeichnung keine weitere“, `TestReplayLebendpruefungNummerNachDemEnde`). | `SPEC-049` Punkt 7; Rolle Architect | `internal/hexagon/services/export_test.go` · „LetzteNummer ist letzteNummer eines Cursors“; `internal/hexagon/services/replay.go` · „&cursor{status: "I"}“ | nein | — (Hinweis) |
| F-443 | INFO | Zur Frage, ob die Mutationen gegen §3.10 reichen. Die Tabelle in §7 deckt jede Brücken-Funktion, die beiden Subtests mit `t.Context()`, die Literale mit Feldnamen und drei Testdaten-Funktionen. Die übrigen 14 Testdaten-Funktionen habe ich mit dem Stand davor verglichen: Sie liefern dieselben Literale wie die früheren Globalen. Weil jeder Aufruf einen frischen Wert liefert, teilen sie keine Slices und Maps mehr. Ein `reflect.DeepEqual` gegen sie ist damit strenger als vorher, nicht schwächer. Die sechs eigenen Mutationen in Bereichen, die die Tabelle nicht nennt, waren alle rot. Nach meinem Urteil reicht die Abdeckung. Ob sie die DoD trägt, prüft der Verifier. | `AGENTS.md` §3.10; Plan §6 Risiko *Prüfung geht im Umbau verloren*; Rolle Verifier | `docs/plan/planning/in-progress/slice-harness-blackbox-kern.md` · „Mutationen des Implementers“ | ja (Mutationen R1 bis R6) | — (Hinweis) |

**Eigene Mutationen** (frische Kopie je Mutation, ohne Mutation grün):

| ID | Mutation | Weg im Test | Ergebnis |
|---|---|---|---|
| R1 | `parameterStelle`: „erwartet“ und „empfangen“ bei der Anzahl vertauscht | Brücke `ParameterStelle` | rot, `TestReplayExtendedParameterStelle` |
| R2 | `abweichung`: Vergleich von `params` entfernt | Brücke `Abweichung` | rot, `TestReplayExtendedFelder` |
| R3 | `vtLeerraum`: `err == nil &&` zu `err != nil \|\|` | Brücke `VTLeerraum` | rot, `TestLebendpruefungServerversion` |
| R4 | `record.go`: Flush schließt keine Gruppe mehr ab | Testdaten-Funktionen `cParse()`, `cFlush()`, `sErr()` … | rot, `TestRecordExtendedFlushGruppen` |
| R5 | `Group.validate`: nur „sync vor dem Ende“ gemeldet, nicht „flush am Ende der letzten Gruppe“ | `model_test`, qualifizierte Literale | rot, `TestValidateFehler/letzte_Gruppe_endet_mit_flush_(abgeschnitten)` |
| R6 | `letzteNummer`: erste statt letzte Interaktion | Brücke `LetzteNummer` | rot, `TestReplayLetzteNummer` |

## Antwort auf die Schwerpunkte

1. **Brücken nach `SPEC-049` Punkt 7.**
   - Beide `export_test.go` gehören zum Paket des Codes und sind dort die einzige Testdatei. Alle übrigen Testdateien sind `model_test` bzw. `services_test`.
   - Die Brücken enthalten nur Funktionen. Es gibt keine Variable, keine Funktion `Test…`, `Benchmark…`, `Example…` oder `Fuzz…`, keinen Typ-Alias und keinen Zustand auf Paketebene.
   - `Klasse`, `IstLebendpruefung`, `VTLeerraum`, `Abweichung` und `ParameterStelle` reichen nur weiter.
   - `LetzteNummer` legt bei jedem Aufruf einen lokalen `cursor` an und setzt nichts an der übergebenen Session. Wem die Auslegung zusteht, regelt F-441; das Urteil steht in F-442.
2. **Verdeckt statt behoben.**
   - Die 17 früheren Globalen sind jetzt Funktionen, die bei jedem Aufruf ein neues Literal liefern. Sie teilen keinen Zustand, auch keine Closure über eine Variable.
   - Im Diff kommt kein `_ =` hinzu. Die fünf vorhandenen Stellen (`_ = s.CloseSession`, `id, _ =`, `s, _ =`) standen schon am Stand `839e61c^`.
   - Unter `internal/` gibt es kein `//nolint` und keine Änderung an `.golangci.yml`.
3. **Prüfung im Umbau verloren.**
   - Stichproben mit `git show 839e61c^:<datei>` geprüft: `TestRecordEndeNachErfolgreichemUpstream`, `TestFehlerKette`, `TestFehlerGleichrangig`, `TestCodeTabelle`, `TestReplayLetzteNummer`, `TestReplayExtendedFelder` und die Testdaten in `record_extended_test.go` und `replay_extended_test.go`.
   - Die Subtests mit `t.Context()`: Der Kontext wird erst nach dem Ende des Subtests abgebrochen. Beide lesen das Ergebnis ihrer Goroutine noch im Subtest, das Verhalten ist also gleich.
   - Die Literale mit Feldnamen: `model.Meldung` hat genau die Felder `Code` und `Text`. Der `DeepEqual` vergleicht dasselbe wie vorher.
   - `TestCodeTabelle` liest weiter nur `fehler.go`. Die Brücke ändert die Menge der Codes nicht.
   - Ich habe keinen Test gefunden, der weniger prüft als vorher.
4. **Testliste und Abdeckungs-Deklarationen.** Beide sind unverändert, siehe *Ausgeführte Läufe*. `make abdeckung-check` ist grün.
5. **Mutationen.** Siehe F-443 und die Tabelle oben.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/model/*_test.go` | geprüft, ohne Befund. Paket `model_test`, Literale mit Feldnamen, keine Globale, kein Befund von `make lint`. |
| `internal/hexagon/services/*_test.go` (sieben Dateien) | geprüft, ohne Befund. Paket `services_test`, Fakes implementieren die exportierten Ports, keine Globale, kein Befund von `make lint`. |
| `internal/hexagon/model/export_test.go` | geprüft, ohne Befund. |
| `internal/hexagon/services/export_test.go` | geprüft. Befund F-441, Hinweis F-442. |
| Produkt-Code unter `internal/hexagon/` | geprüft, ohne Befund. Der Diff ändert ihn nicht, es gibt keinen neuen Export nur für Tests (Plan §1). |
| Hard Rule 3.2, 3.6 (Suppression, Gate-Lockerung) | geprüft, ohne Befund. Kein `//nolint`, `.golangci.yml` ist unverändert. |
| Hard Rule 3.3 | geprüft, ohne Befund. Der Commit enthält keinen Move. Die Moves `350ef18` und `bd5132f` davor sind rein. |
| Hard Rule 3.7, 3.11 (Kommentare) | geprüft, ohne Befund. Die neuen Kommentare (`// Testdaten…`, Kopf der Brücken, `// streng sind die Optionen mit FailOnUnconsumed.`) beschreiben, was da ist. |
| Hard Rule 3.9 (Plan folgt Korrektur) | geprüft, ohne Befund. §3 und §6 sind im selben Commit nachgezogen, §1 stimmt mit dem Diff überein, der Kopf mit `SPEC-049` (`make kopf-check` grün). |
| Hard Rule 3.12 | geprüft. Befund F-441. |
| Commit-Message | geprüft, ohne Befund. Sie nennt `slice-harness-blackbox-kern`, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) und `LH-QA-07`, aber keine Struktur-ID (`AGENTS.md` §5 Regel 1). |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 2 |

**Summary:** 0 HIGH · 1 MEDIUM (F-441: Brücke mit Cursor als Randform erst im Code-Commit benannt und dem Review statt dem Architect zugewiesen) · 0 LOW · 2 INFO (F-442 Urteil zur Brücke für den Architect, F-443 Mutationen reichen). Wiederkehrende Klasse: „Randform im Code entschieden“ (`BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`, viertes Auftreten in Reviews).

**Finding-Klassen dieses Laufs:** Randform im Code entschieden

## Verdikt

**Merge-blockierend:** ja, wegen F-441. Am Code selbst blockiert nichts: Die Tests prüfen nicht weniger als vorher, und nach dem Urteil in F-442 hält die Brücke den Zweck von Punkt 7 ein. Blockiert ist, dass die Auslegung von Punkt 7 bisher keinen Entscheider und keinen dauerhaften Ort hat.

**Übergabe:** F-441 und F-442 gehen an den Architect, der über die Auslegung von `SPEC-049` Punkt 7 für einen selbst angelegten Empfänger entscheidet. Danach geht F-441 an den Implementer zum Nachziehen. F-443 geht an den Verifier. Die Finding-Klasse geht in §7 des Slice und von dort in den Zähler. Dieser Report ersetzt keine Verifikation.
