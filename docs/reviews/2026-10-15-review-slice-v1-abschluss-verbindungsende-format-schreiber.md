# Review-Report: slice-v1-abschluss-verbindungsende-format-schreiber — 2026-10-15

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten)

**Gegenstand:** `slice-v1-abschluss-verbindungsende-format-schreiber` (Schritt 35b), Diff `6c6dcd9..HEAD` (Entscheidung des Architect `4940cce`, Umsetzung `48d50a3`)

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec0 · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-v1-abschluss-verbindungsende-format-schreiber` (§1, §3, §4, §6 mit Zählung der Mutationstabelle, §7 Belege des Implementers, §8); Nehmer-Pläne `slice-v1-abschluss-verbindungsende-record`, `-sqlite`, `-upstream`, `-doku` (Bindungen, Nachzählung)
- [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md) (Aufzeichnung kennt das Verbindungsende)
- `LH-FA-25` (.a), `LH-FA-22`, `LH-QA-06`; `SPEC-001` (*Schreiber in `yaml`*), `SPEC-002`, `SPEC-004`
- `AGENTS.md` (Hard Rules, insbesondere §3.7 bis §3.14)
- vorherige Reports am Modul: `2026-10-15-review-slice-v1-abschluss-verbindungsende-format`

**Eigene Läufe des Reviewers** (Kopie des Arbeitsbaums `git archive HEAD` im Scratchpad, Image `pgwire-recorder:test`, `--network=none`, Bind-Mount, kein `make gates`): Ausgangslauf `go test -run 'Schreiben|Roundtrip|HasEnd'` grün; 31 Mutanten, jeweils die Originaldatei zurückkopiert:

- Version: `return 2` statt `model.VersionMitEnde` **grün** (äquivalent, siehe F-4); Feld 1 / 0 / 2 / 3 mit und ohne Ende je einzeln; immer 1; immer 2; Feld übernehmen — alle rot wie in §7 behauptet, bis auf die als nicht trennbar gemeldete Zeile „Feld 1 mit Ende" (F-1).
- `HasEnd`: nur `Groups`, nur `Responses`, nur erste / nur letzte Interaktion, nur erste / nur letzte Session, wahr ohne Sessions — je rot an genau den Fällen der Tabelle in §7.
- Stellen: ein Mutant, der „nicht letzte Interaktion" ablehnt, färbt zusätzlich „zwei Enden"; einer, der „zwei Enden" ablehnt, färbt nur diesen Fall; einer für „Feld neben `type`" färbt dieses und `TestSchreibenEndeDarstellung/mit_Feld_neben_type`.
- Darstellung: Flow-Form, Einrückung eine Stufe tiefer, Zeilenumbruch mehr — rot (`TestSchreibenEndeDarstellung`, die Einrückung zusätzlich die „richtige Stelle" und den Roundtrip).
- `VersionMitEnde = 3` im Modell: rot (`TestVersionMitEnde` und die Fälle mit Ende).
- Hilfsdatei / Umgebungsvariable des Implementers: nicht im Diff (`git diff --stat` ohne weitere Dateien, `git status` sauber, kein `Getenv` im Paket `recording`). `//nolint`: keines unter `cmd`, `internal`, `test`.
- `make docs-check`: 0 Befunde.

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt**; bei Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | In `TestSchreibenVersionNachInhalt` tragen alle Fälle mit Ende das Modellfeld 1; der Mutant „Feld 1 schreibt 1" färbt darum fünf Fälle statt nur „Feld 1 mit Ende" (selbst reproduziert; §7 nennt es offen als „nicht zu trennen"). Die Zusage „Feld 1 mit Ende schreibt 2" ist ein eigenes Element der Aufzählung in §6 und hat keinen Fall, den nur ihr Mutant bricht. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `internal/adapters/driven/recording/yaml_schreiben_ende_test.go` · `{"Ende nur in der letzten Session", sRec(1, sSession(1, "u", ohne), sSession(2, "u", mitEnde)), "2"}` | ja — Mutant `rec.HasEnd() && rec.Version != 1` in `formatVersion` | Element einer Aufzählung ohne isolierenden Fall |
| F-2 | LOW | §1 *Ziel* („Aufzeichnungen als Text im Test") und DoD-Punkt 2 („Fixtures als Text im Test") nennen Text-Fixtures; die Tests von Version, Stellen und Roundtrip bauen die Aufzeichnung als Modell, und §3 sagt inzwischen „als Modell im Test". Nur die Fälle „ohne Ende, Bytes bleiben" und die Darstellung tragen Text. | `AGENTS.md` §3.9 (Plan folgt jeder Korrektur) | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-format-schreiber.md` · `Nachweis ohne `record`: Aufzeichnungen als Text im Test` | nein — kein Gate liest Plan gegen Diff | Plan §1/DoD nicht auf §3 nachgezogen |
| F-3 | LOW | Die Deklaration von `TestRoundtripMitEnde` (und DoD-Punkt 2, `SPEC-001`) sagt „Laden, Schreiben und erneutes Laden"; der Test geht vom handgebauten Modell aus (Schreiben, Laden, Vergleich), er lädt vorher nichts. Der Weg Laden, Schreiben ist nur für eine Aufzeichnung ohne Ende (`ohne Ende, Bytes bleiben`) geprüft. | `AGENTS.md` §3.11 (Zusage nur, was ein Test prüft) | `internal/adapters/driven/recording/yaml_schreiben_ende_test.go` · `// Abdeckung: LH-FA-25/Happy, LH-FA-22/Happy, LH-QA-06/Messung — Laden, Schreiben und erneutes Laden` | ja — Mutant, der nur eine geladene (nicht handgebaute) Form mit Ende verändert, bliebe grün | Zusage weiter als der Test |
| F-4 | INFO | Der Mutant `return 2` statt `model.VersionMitEnde` bleibt grün (selbst reproduziert). Die Begründung „äquivalent, Verhalten gleich" trägt nach §3.10: Die Schreibweise ist nicht beobachtbar; erkennbar ist sie nur zusammen mit `VersionMitEnde = 3` (rot im Schreiber-Test und in `TestVersionMitEnde`). Der Satz „ein Literal ist es nicht" in §6 steht ohne eigenen Test. | `AGENTS.md` §3.10 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-format-schreiber.md` · `**Grüner Mutant:** `return 2` statt `model.VersionMitEnde`` | ja — Mutation `VersionMitEnde = 3` | Äquivalenter Mutant, Begründung genügt |
| F-5 | INFO | §6 nennt für das Modellfeld `version` „0, 3, negativ"; die Fälle decken 0 und 3, keinen negativen Wert. Der Code liest das Feld nicht, ein Mutant, der nur negative Werte anders behandelt, bliebe grün. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `internal/adapters/driven/recording/yaml_schreiben_ende_test.go` · `{"Feld 3 mit Ende", sRec(3, sSession(1, "u", mitEnde)), "2"}` | ja — Mutant `rec.Version < 0` in `formatVersion` | Element einer Aufzählung ohne Fall |
| F-6 | INFO | §8 sagt „kein Nehmer eingetragen"; der Commit `4940cce` trug die Bindung „Schreiber prüft die Form nicht" in `slice-v1-abschluss-verbindungsende-record` (§1, §6, §8 mit Zählung: Liefer-Punkte 3, Schichten zwei) ein, die Messzeile des Implementers nennt das nicht. In `-sqlite` und `-upstream` steht die Bindung nicht; sie liegt dort indirekt über den Use Case in `-record`. | `AGENTS.md` §3.13 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-format-schreiber.md` · `kein Nehmer eingetragen` | nein | Messzeile im Plan unvollständig gegenüber dem Stand der Nehmer |
| F-7 | INFO | `SPEC-001` sagt zu, der Schreiber lese die Aufzeichnung „einmal je Schreiben"; kein Test prüft das (die Zählung ruft `HasEnd` aus `toDTO` einmal). Es ist Technik-Stratum, kein Kommentar, und die Aussage ist als Zustandsprüfung ohne Folge für den Aufrufer formuliert. | Maintainability | `spec/spezifikation.md` · `Der Schreiber liest die Aufzeichnung einmal je Schreiben` | nein | Zusage in der Spezifikation ohne Test |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driven/recording/yaml.go` (`formatVersion`, `toDTO`; Version nach dem Inhalt vor `toDTO`, `model.VersionMitEnde` und `model.Version`, kein Literal, kein Lesen von `rec.Version`) | geprüft, ohne Befund |
| `internal/hexagon/model/recording.go` (`HasEnd` ohne Stellenprüfung, nutzt `endet()` für Responses und Gruppen, Core-Reinheit) | geprüft, ohne Befund |
| `internal/adapters/driven/recording/yaml_schreiben_ende_test.go` (Entscheidung (a): Schreiben nimmt an, Laden `PGR-E3003`, richtige Stelle lädt; je Stelle ein Fall mit eigenem isolierbaren Mutanten, außer „nicht letzte Interaktion" färbt prinzipbedingt „zwei Enden" mit) | geprüft, ohne Befund (außer F-1, F-3, F-5) |
| `internal/hexagon/model/hasend_test.go` (ein Fall je Ort und Session, 13 Fälle; alle sieben Mutanten an genau den genannten Fällen rot) | geprüft, ohne Befund |
| Darstellung (Blockform in `responses` und `server`, Einrückung, Version 1 byte-gleich ohne Ende) | geprüft, ohne Befund; die Blockform ist die allgemeine Ausgabe des Schreibers, ein Sonderpfad für das Ende existiert nicht |
| Kommentare (`formatVersion`, `HasEnd`, Testhilfen) | geprüft, ohne Befund (§3.7/§3.11: Zusage mit Test, kein Konjunktiv, keine Chronik) |
| `//nolint`, Gate-Lockerung, Hilfsdatei oder Umgebungsvariable des Implementers im Diff | geprüft, ohne Befund (keine) |
| `docs/user/abdeckung-unit.md` | geprüft, ohne Befund (Zeilen entsprechen den Deklarationen) |
| `spec/spezifikation.md` (`SPEC-001` *Schreiber in `yaml`*, `SPEC-002`-Blockform) | geprüft, ohne Befund außer F-7 (Architect-Commit vor dem Code; Stratum nennt keine ADR, keinen Slice) |
| Randformen §6 (§3.12): im Code-Commit `48d50a3` unverändert; keine neu eingetragen | geprüft, ohne Befund |
| Nehmer `-record` (§1, §6, §8 mit Zählung), `-sqlite`, `-upstream`, `-doku` (§3.13, Adresse nimmt an, Nachzählung) | geprüft, ohne Befund außer F-6 |
| Größe (§3.14): Summe 311 (Einfügungen 310 + Löschungen 1, Pfade `internal docs/user`, nachgemessen), Meldeschwelle 550 und Marke 700 nicht erreicht; Liefer-Punkte 2, Schichten zwei | geprüft, ohne Befund |
| `docs/plan/planning/` (Roadmap „In Arbeit", `done/…-format` Verweis nach `in-progress/`, Move getrennt vom Inhalt) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Element einer Aufzählung ohne isolierenden Fall · Plan §1/DoD nicht auf §3 nachgezogen · Zusage weiter als der Test · Äquivalenter Mutant mit tragender Begründung · Messzeile im Plan unvollständig gegenüber dem Stand der Nehmer · Zusage in der Spezifikation ohne Test

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Die drei LOW sind Nacharbeit am Plan und an Testfixtures, keine Fehler im Verhalten. Die selbst gefahrenen Mutationen bestätigen die Tabelle in §7; die als nicht trennbar oder äquivalent genannten Zeilen („Feld 1 mit Ende", `return 2`, „nicht letzte Interaktion" gegen „zwei Enden") sind ehrlich benannt, die ersten beiden tragen nach §3.10 mit den Anmerkungen F-1 und F-4, die dritte ist prinzipbedingt (zwei Enden in einer Session enthalten ein Ende in einer nicht letzten Interaktion; die Gegenrichtung ist isoliert).

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** — er wird über Läufe hinweg nicht wieder gelesen. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
