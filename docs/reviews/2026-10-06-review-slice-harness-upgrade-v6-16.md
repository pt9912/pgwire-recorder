# Review-Report: slice-harness-upgrade-v6-16 — 2026-10-06

**Review-Art:** Harness, Pins, vendored Baseline und Planung, geprüft gegen Plan, Entscheidungen (§6 *Randformen*, entschieden in `53ea519`) und Hard Rules (Modul 10). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff 53ea519..01a45ae`, also die Commits `01f9933` (d-check `v0.82.0` gepinnt), `0fca02e` (Baseline `v6.16.0` statt `v6.13.0`, Symlinks, `MR-001`, Stand-Zeilen), `aec4aab` (Abgleich und Regelblock-Tabelle) und `01a45ae` (Skills, ADR-Index, Folge-Slice, Drift-Log). Schwerpunkte laut Auftrag: `d-check.mk`, Baseline-Tausch, Regelblock-Tabelle, Abgleich gegen das echte Delta, Folge-Slice, kein `v6.13.0` in lebenden Dateien.

**Skill:** `.harness/skills/reviewer.md` @ `01a45ae` (mit den Regeln aus `v6.16.0`: LOW nur mit Konventions-Anker, kein HIGH/MEDIUM ohne Failure-Szenario, `pfad` als Kurzzitat)
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-harness-upgrade-v6-16.md` (ganz, Stand `01a45ae`), `docs/plan/planning/open/slice-harness-gate-index-werkzeug-teil.md` (ganz)
- `AGENTS.md` §1, §3.3, §3.5, §3.8 bis §3.12; `harness/conventions.md` (ganz); `harness/README.md` §Sensors; `harness/sensors/*.md`; `.claude/agents/architect.md`; `.claude/commands/implement-slice.md` (Randform-Rückgabe)
- Baseline alt: `git archive 53ea519 .harness/baseline/v6.13.0` in einen Temp-Baum; neu: `.harness/baseline/v6.16.0/`. Vergleich mit `diff -r -I 'v6\.1[36]\.0'` (Tag-Zeilen ausgeblendet) und `git show -M --name-status 0fca02e`
- `v6.16.0` · `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt; `v6.16.0` · `regelwerk/grundlagen-referenz-richtung.md` §Spec-Straten; `v6.16.0` · `regelwerk/modul-10-review-harness.md`; `v6.16.0` · `regelwerk/modul-13-quality-gates.md`; die geänderten Vorlagen (14) des Stands
- Vorheriger Report: Review zu `slice-harness-lint-werkzeug` vom selben Tag; höchste vergebene Nummer vor diesem Lauf F-433

**Ausgeführte Läufe:** `make docs-check` (280 Dateien, 0 Befunde, Image per Digest `sha256:d28e9437…`), `make kopf-check` (Exit 0), `make baseline-verify` (`v6.16.0 OK — 54 Dateien`). `docker run --network none ghcr.io/pt9912/d-check:v0.82.0 --print-mk` gegen `d-check.mk` gehalten; `docker image inspect … RepoDigests` gelesen. In `.harness/baseline/v6.16.0/`: `sha256sum -c SHA256SUMS` OK, Pfadliste gleich `find regelwerk templates -type f | LC_ALL=C sort`, Liste sortiert. `find . -path ./.git -prune -o -xtype l -print` leer. `grep -rn "v6\.13\.0"` über die Pfade aus DoD Punkt 2. `make gates` lief nicht.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-434 | MEDIUM | `AGENTS.md` §3.12 lässt als Ort einer entschiedenen Randform weiter „in der Spezifikation (Technik-Stratum), auch für Harness-Werkzeuge …, oder einer ADR“ zu. `v6.16.0` sagt für die Festlegungen eines Werkzeugs, dass sie in der Spezifikation stehen und dass „die gesperrte ADR dieser Ort nicht“ ist. Der Abgleich in §7 führt Befund (2) als behoben, nennt §3.12 aber nicht, obwohl die DoD den Abgleich gegen §3.1 bis §3.12 verlangt. Failure-Szenario: Ein Slice der Lint-Reihe entscheidet eine Randform eines Gates nach §3.12 in dessen ADR. Nach `Accepted` ist die ADR gesperrt (§3.5), und eine später auftauchende Randform braucht eine Ersetzungs-ADR. Genau das schließt die Baseline mit „fortgeschrieben, ohne Folge-ADR“ aus. §3.8 und die Pläne der Lint-Reihe zeigen heute auf die Spezifikation und mildern das, die Regel selbst aber lässt es zu. | `AGENTS.md` §3.12; Plan §6 Risiko *Verkörperte Regel und neue Baseline widersprechen sich*; `v6.16.0` · `regelwerk/grundlagen-referenz-richtung.md` §Spec-Straten | `AGENTS.md` · „Harness-Werkzeuge (seit slice-lastenheft-pruefbarkeit), oder einer ADR“; `docs/plan/planning/in-progress/slice-harness-upgrade-v6-16.md` · „Spezifikation §11, beide Sensor-Dateien und die“ | nein (Lesen) | Verkörperte Regel nach Baseline-Sprung nicht nachgezogen |
| F-435 | LOW | Die Zeile `grundlagen-harness-dateien.md` der Regelblock-Tabelle trägt den Wert *Träger kommt mit* und daneben den Satz, dass der Bootstrap einen Teil des Regelblocks (den Teil des Gate-Index unter `harness/mk/`) nicht ablegt. Die Zelle hat damit zwei Zustände. Ohne Dauer steht der fehlende Teil als *kommt nicht mit*. Die Tabelle verlangt je Regelblock genau einen Wert, *kommt nicht mit* nur mit Grund und Dauer. Für `modul-02` und `modul-15` teilt sie den Block deshalb in zwei Zeilen je Abschnitt. | `harness/conventions.md` §Welche Regelblöcke des Baums hier einen Träger haben („Je Regelblock … genau einer von drei Werten“) | `harness/conventions.md` · „Einen Teil des Gate-Index unter `harness/mk/` legt der Bootstrap nicht ab“ | nein (Lesen) | Regelblock-Zelle trägt zwei Werte |
| F-436 | LOW | Bei der Teilersetzung in §Adoptierte Konventions-Quellen ist die öffnende Klammer „(Stand-Zeile in“ weggefallen, die schließende nach „CHANGELOG.md im Kurs-Repo)“ steht noch. Der Satz hat jetzt eine schließende Klammer ohne öffnende. | Reviewer-Skill LOW (einmaliger Tippfehler) | `harness/conventions.md` · „CHANGELOG.md im Kurs-Repo); für harte Reproduzierbarkeit“ | nein (Lesen) | Teilersetzung lässt Satzrest stehen |
| F-437 | INFO | §3.3 im Tausch-Commit nachgemessen und nicht übernommen: `git show -M --name-status 0fca02e` zeigt 55 × `R` (alle Baum-Dateien mit Rename erkannt) und 6 × `M`. `git log --follow` auf `.harness/baseline/v6.16.0/regelwerk/grundlagen-harness-dateien.md` (die Datei mit dem größten Delta) reicht bis `d027ecc`. Der Schutzzweck von §3.3 bleibt also gewahrt. Plan §6 *Alte Baseline* führt `git log --follow` als akzeptiertes Negativ; eingetreten ist es nicht. Dass die Auslegung „Ersetzen eines vendored Baums ist kein Move“ in einem Slice-Plan steht und nicht in einem Carveout oder einer Adaption, bleibt Sache des Architects. | `AGENTS.md` §3.3; Plan §6 *Alte Baseline*; Rolle Architect | `docs/plan/planning/in-progress/slice-harness-upgrade-v6-16.md` · „`git log --follow` über den Baum ist ein akzeptiertes Negativ“ | ja (`git show -M --name-status 0fca02e`) | — (Hinweis) |
| F-438 | INFO | Die Wahrheit von `MR-000` hängt an einer Lesart, die der Baseline-Text nicht ausdrücklich trägt. Der Folge-Slice sagt, dass `harness/README.md` die Targets der Werkzeug-Fragmente führen darf, solange das Werkzeug keinen eigenen Teil schreibt. `v6.16.0` · `regelwerk/grundlagen-harness-dateien.md` sagt nur, dass ein Repo *ohne* Werkzeug-Fragmente bei einer Datei bleibt; dieses Repo hat Fragmente von ai-harness-init. Die Vorlage `harness/README.template.md` sagt im Kommentar „Targets aus Werkzeug-Fragmenten gehören NICHT hierher“. Der Start des Folge-Slice hängt allein an einem externen Lauf, ohne Datum. Dass er als Slice bleibt, ist Entscheidung des Nutzers und wird hier nicht bestritten. Adressat für Lesart und Trigger: Architect, gegebenenfalls Frage an das Kurs-Repo. | `harness/conventions.md` `MR-000`; Rolle Architect | `docs/plan/planning/open/slice-harness-gate-index-werkzeug-teil.md` · „ist das die eine Datei, die die Regel für diesen Fall zulässt“ | nein | — (Hinweis) |
| F-439 | INFO | Der Skill verweist als Report-Gerüst auf `docs/reviews/review-report.template.md`. Die Datei gibt es im Repo nicht; das Gerüst liegt nur im vendored Baum, und dort hat `v6.16.0` die Spalte `Pfad` auf „Datei · wörtliches Kurzzitat“ umgestellt. Die Zeile stand schon vor dem Diff, sie ist nicht Teil der Änderung. | Reviewer-Skill §Output-Schema | `.harness/skills/reviewer.md` · „`docs/reviews/review-report.template.md`, ein Report pro Lauf“ | ja (`ls docs/reviews/`) | — (Hinweis) |
| F-440 | INFO | Der eigene Plan nennt `.harness/baseline/v6.13.0/` an sieben Stellen als Pfad in Code-Spans (§2, §3, §6, §8). Ob das unter „verweist auf einen Pfad“ fällt oder unter „den alten Stand als Wort nennen, dieser eingeschlossen“, ist eine Frage der DoD. Adressat ist der Verifier. | Plan DoD Punkt 2; Rolle Verifier | `docs/plan/planning/in-progress/slice-harness-upgrade-v6-16.md` · „duldet ein Tag-Verzeichnis zur Zeit“ | ja (`grep -rn "baseline/v6\.13\.0" docs/plan/planning/{open,next,in-progress}`) | — (Hinweis) |

## Antwort auf die Schwerpunkte

1. **`d-check.mk`.** Neu erzeugt mit genau vier Adaptionen: `diff` gegen `d-check --print-mk` von `v0.82.0` (netzlos) zeigt nur Kopfkommentar, `DCHECK_DIGEST` gepinnt, `doc-check` → `docs-check` und das Muster von `doc-help` mit `^docs?-`. Rezepte unverändert. Der Digest `sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8` ist gleich dem lokalen `RepoDigests`, und `make docs-check` läuft über `ghcr.io/pt9912/d-check@sha256:…`. Der Kopfkommentar sagt nur zu, was der Diff zeigt (`AGENTS.md` §3.11): Er nennt jetzt die `doc-help`-Adaption, und „Übrige advisory doc-*-Targets verbatim“ stimmt. Den Hinweis des erzeugten Kopfs auf die Release-Notes trägt er nicht mehr; der Plan verlangt das nicht.
2. **Baseline-Tausch.** Ein Commit (`0fca02e`): Nur `v6.16.0` liegt unter `.harness/baseline/`. `SHA256SUMS` hat 54 Zeilen in der Form `<sha256>  <pfad>`, relativ, `LC_ALL=C` sortiert, deckt genau die Dateien unter `regelwerk/` und `templates/`, und `sha256sum -c` ist OK. Beide Symlinks unter `.claude/rules/` zeigen im selben Commit auf `v6.16.0`, kein hängender Link. `MR-001` trägt den Link jetzt als Code-Span mit unverändertem Pfad samt `v6.13.0` und Anker; der Text ist gleich, wie §6 *Was am alten Pfad hängt* (2) entschieden hat. Zu §3.3 siehe F-437.
3. **Regelblock-Tabelle.** Gegen `ls .harness/baseline/v6.16.0/regelwerk/*.md` vollständig: Alle 26 Dateien stehen mindestens einmal (28 Zeilen; `modul-02` und `modul-15` je zweimal nach Abschnitt). Jeder Wert stammt aus der geschlossenen Menge (18 × *Träger kommt mit*, 8 × *liegt bei, nicht verdrahtet*, 2 × *kommt nicht mit*). Die Messzeile nennt `v6.16.0`. Zur Zelle `grundlagen-harness-dateien.md` siehe F-435.
4. **Abgleich gegen das echte Delta.** Mit ausgeblendeten Tag-Zeilen ändern sich genau `README.md`, `grundlagen-begriffe.md`, `grundlagen-harness-dateien.md`, `grundlagen-referenz-richtung.md`, `modul-03-spec.md`, `modul-10-review-harness.md`, `modul-13-quality-gates.md` und `modul-15-observability.md` inhaltlich, dazu 14 Vorlagen. Das deckt sich mit §7 *Belege*; übersehen ist kein Regelblock. Stichproben: `modul-10` (Stil-Polizist, Failure-Szenario, `pfad`) ist im Skill wortgleich mit der Vorlage. Bei `grundlagen-referenz-richtung` und `modul-03` (Werkzeug-Festlegungen im Technik-Stratum) sind ADR-Index, Spezifikation §11 und die Sensor-Dateien (verlinken `SPEC-047`/`SPEC-048`) nachgezogen, `AGENTS.md` §3.12 nicht (F-434). `grundlagen-harness-dateien`, `modul-13` und die Vorlagen von `AGENTS.md`, `harness/README.md`, `Makefile` und `.d-check.yml` (Teil des Werkzeugs) gehen an den Folge-Slice. `modul-15` wird von `harness/erfassung-feldliste.md` getragen („leer heißt unbekannt“). Die Vorlagen von Slice, Welle-Results, Review-Report und Konventionen (Form `BEO-<KUERZEL>/<slug>`, Ablage `observations/`) werden append-only behandelt; die offenen Pläne behalten `BEO-<NNN>` im Regelabsatz von §7, der neue Folge-Slice führt die neue Form.
5. **Folge-Slice.** `slice-harness-gate-index-werkzeug-teil` liegt in `open/`, aus der Vorlage `v6.16.0` (neue Register-Form), `Verantwortlich: —`, Kopf grün nach `make kopf-check`. Drift-Log und Plan §3 nennen ihn. Dass er als Slice bleibt, ist entschieden. Zur Lesart, auf der `MR-000` dabei ruht, siehe F-438.
6. **Kein `v6.13.0` in lebenden Dateien.** `grep -rn "v6\.13\.0"` über `AGENTS.md`, `README.md`, `harness/` (ohne `harness/conventions/done/`), `.claude/`, `.harness/skills/`, `tools/`, `Makefile`, `*.mk`, `harness/mk/`, `spec/`, `.d-check.yml`, `.githooks/`, `docs/plan/adr/README.md` und `docs/plan/planning/observations/` ist leer. In Plänen kommt das Wort nur im eigenen Plan, im Folge-Slice und im Drift-Log der Roadmap vor; ein Pfad unter `.harness/baseline/v6.13.0/` steht nur im eigenen Plan (F-440). Außerhalb davon nennen den alten Stand nur eingefrorene Reports unter `docs/reviews/` und `done/`, wie §1 vorsieht. `d-check v0.79.0` nennt kein lebendes Dokument mehr.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `d-check.mk` | geprüft, ohne Befund. Vier Adaptionen, Digest gepinnt, Kopf nach §3.11. |
| `.harness/baseline/v6.16.0/` | geprüft, ohne Befund. Inhalt gegen `SHA256SUMS`, Vollständigkeit, Form; gegenüber `v6.13.0` keine Datei neu, umbenannt oder entfallen. |
| `.claude/rules/` | geprüft, ohne Befund. Beide Symlinks auf `v6.16.0`, auflösend. |
| `harness/conventions/done/` | geprüft, ohne Befund. Form-Reparatur in `MR-001` nach §6, Text unverändert. |
| `harness/conventions.md` | geprüft. Befund F-435, F-436. §Baseline, Asset-URL, Stand-Zeile („Kurs-Welle 159 · 2026-10-06“ = `regelwerk/README.md`), Messzeile und Zeile `MR-001` sonst stimmig. |
| `AGENTS.md` | geprüft. Asset-URL in §1 nachgezogen. Befund F-434 (§3.12). |
| `.harness/skills/` | geprüft, ohne Befund im Diff. Beide Skills folgen den Vorlagen `v6.16.0`. Hinweis F-439 (Zeile außerhalb des Diffs). |
| `docs/plan/adr/README.md` | geprüft, ohne Befund. Ein Satz nach `adr/README.template.md`; keine angenommene ADR berührt (§3.5). |
| `docs/plan/planning/in-progress/` | geprüft. Plan §1, §3 und §6 folgen dem Diff (§3.9): Jede Datei im Diff steht in §3, kein `.go`-File, `.d-check.yml` unverändert. Drift-Log-Zeile ist eine Umplanung, keine Schließung. Hinweise F-437, F-440. |
| `docs/plan/planning/open/` | geprüft. Folge-Slice formgerecht, `make kopf-check` grün. Hinweis F-438. |
| Hard Rules §3.3, §3.5, §3.6 | geprüft, ohne Befund. Keine Lockerung an `.d-check.yml`, keine ADR geändert; §3.3 siehe F-437. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Verkörperte Regel nach Baseline-Sprung nicht nachgezogen · Regelblock-Zelle trägt zwei Werte · Teilersetzung lässt Satzrest stehen

## Verdikt

**Merge-blockierend:** ja, wegen F-434 (MEDIUM). F-435 und F-436 sind klein; F-437 bis F-440 verlangen keine Aktion des Implementers.

**Übergabe:** Die Findings gehen an den Implementer, F-437 und F-438 zusätzlich an den Architect, F-440 an den Verifier. Die Finding-Klassen gehen in §7 des Slice und von dort in den Zähler. Der Report ersetzt keine Verifikation.
