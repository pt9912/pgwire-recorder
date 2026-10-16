# Review-Report: slice-v1-abschluss-verbindungsende-hilfetext-binary — 2026-10-16

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-hilfetext-binary` (Schritt 39b), Diff `1de25738..a8309486` (Basis laut Plan §4; die Commits `179ab424` und `06f959d9` des Architect liegen dazwischen; Planungs-Commits nach `a8309486` sind nicht Gegenstand).

**Skill:** `.harness/skills/reviewer.md` @ `8f90c592` · <!-- d-check:ignore (Skill-Pfad des Repos, kein Link-Ziel des Reports) -->
**Modell:** `claude-sonnet-5-5` · **Datum:** 2026-10-16

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

- Slice-Plan `slice-v1-abschluss-verbindungsende-hilfetext-binary` (§1, §3, §4, §6, §7, §8), Stand `a8309486`
- [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md), [ADR-0042](../plan/adr/0042-versionsmatrix-der-integrationstests.md) (Bezug des Plans)
- `LH-FA-20.a` (neuer Absatz unter der Optionstabelle, `179ab424`), `LH-FA-24.a`, `LH-FA-25.a`, Abnahmeszenario 16 Satz 3, `SPEC-029` Punkt 2
- Geber `slice-v1-abschluss-verbindungsende-vergleich`, `slice-v1-abschluss-verbindungsende-einspielen`, `-einspielen-extended`, `-einspielen-grenzen`
- `AGENTS.md` (Hard Rules §3.7 bis §3.14)

**Eigene Stichproben** (Zusagen des Implementers nicht übernommen, sondern gefahren):

- Hilfetext, 12 Mutanten in Kopien (Image `deps`, `--network=none`, ein Lauf zur Zeit) gegen `TestPlayEndeHilfe`: ERROR-Satz, „ebenso ein“, FATAL entfernt und PANIC entfernt je Option, „dann läuft die nächste Session“, „ohne Wirkung“, „erwartet“, „an derselben Stelle“, „mit demselben SQLSTATE“, „sonst eine Abweichung“. Jeder färbt genau seinen Teilfall, die Tabelle in §7 stimmt. Dazu zwei eigene Mutanten, die grün blieben (F-1).
- E2E, `make test-integration` im Arbeitsbaum nach Wegwerf-Commit, danach `git reset --hard HEAD~1`: M1 rot nur `TestE2EPlayEndeVergleichGleich`; M2 rot beide Läufe von `TestE2EPlayEndeVergleichAnders` und `TestE2EPlayVergleichFatal`; M3 rot nur `…Anders/play_ende_vergleich_anders_allow`; M4 rot beide Läufe von „Anders“ und drei Bestandstests; M5 rot „Gleich“, beide Läufe von „Anders“ und `TestE2EPlayVergleichFatal`. Alle fünf wie in §7 behauptet.
- Eigene E2E-Mutanten: X1 (zusätzliche `PGR-E4003`-Meldung beim erwarteten Ende) rot nur `TestE2EPlayEndeVergleichGleich`; X2 (zusätzliche `PGR-E4003`-Meldung neben der einen `PGR-E5004`, Exit-Code 5, Zähler 2/2/1 unverändert) rot beide Läufe von „Anders“ allein über die Prüfung „kein `PGR-E4003`“ (F-2).
- Messung nachgerechnet: `cli.go` 12, Hilfetest 60, E2E-Datei 78, Fixtures je 28, `play_vergleich_e2e_test.go` 9, Tabellen 5 + 2 = 7; Summe 222 (Marke 400, Meldeschwelle 300). Die Pfade der Metrik ändern sich nach `f4e9defb` nicht mehr.

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | `TestPlayEndeHilfe` prüft, dass ein Satz im Block seiner Option steht, nicht seine Stelle darin und nicht, dass er im Block der anderen fehlt. Der Mutant „dann läuft die nächste Session“ vor „ebenso ein Ende des Servers“ (der Satz hängt dann am ERROR-Teil) und der Mutant „mit --compare-responses ohne Wirkung“ zusätzlich im Block von `--compare-responses` bleiben grün; §7 nennt das Verschieben in den Block der anderen Option als nicht gefahren. Zuständig: Implementer. | `AGENTS.md` §3.10 (Ortsangabe zählt je Schritt) | `internal/adapters/driving/cli/play_ende_hilfe_test.go` · `strings.Contains(block, f.aussage)` | ja — Kopie mit dem Mutanten, `go test -run TestPlayEndeHilfe` bleibt grün | Hilfetext-Test prüft Anwesenheit, nicht Zuordnung |
| F-2 | INFO | §7 „Nicht mit eigenem Mutanten belegt (a)“ nennt die Prüfung „kein `PGR-E4003`“ nicht trennbar von Exit-Code 5 und der einen `PGR-E5004`-Meldung. Sie ist trennbar: eine zusätzliche `PGR-E4003`-Meldung bei sonst gleichem Ergebnis (X2) färbt beide Läufe von „Anders“ allein über diese Prüfung. Die Prüfung trägt, die Begründung im Beleg ist unrichtig. Zuständig: Implementer (Beleg in §7). | `AGENTS.md` §3.10 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-hilfetext-binary.md` · „ist von der auf Exit-Code 5 und die eine Meldung `PGR-E5004` nicht zu trennen“ | ja — Mutant X2 mit `make test-integration` | Beleg-Aussage enger als die Prüfung |
| F-3 | INFO | Die Fixtures tragen entgegen dem Wortlaut des Architect (`SELECT 1`, `SELECT pg_terminate_backend(...)`) `CREATE TABLE ende_vor` und den DO-Block; §6 ist im Commit `f4e9defb` nachgezogen, dem Commit der Tests. Das entscheidet keine Randform eines Vertrags (der Vergleich der Struktur steht in `LH-FA-24.a`), `AGENTS.md` §3.9 ist erfüllt, §3.12 nicht berührt. Die Fixture bleibt eine handgeschriebene Aufzeichnung; die Form entspricht der aufgezeichneten Fehlerantwort in `TestE2EPlayVergleichFatal`, und der Server prüft die Antwortseite im Lauf, nicht die Aufzeichnungsseite: Das Risiko „handgeschriebene Aufzeichnung“ bleibt mit dem Ausgang „weiter offen“ tragfähig. Keine Aktion. | `AGENTS.md` §3.12, `BEO-REPO/handgeschriebene-aufzeichnung-ohne-beleg-gegen-den-server` | `test/integration/testdata/play-ende-vergleich-gleich.yaml` · `DO $$ BEGIN PERFORM pg_terminate_backend` | nein — Urteil über Plan und Register | Fixture-Form erst im Lauf gegen den Server entschieden |
| F-4 | INFO | Die Historie-Zeile vom 2026-10-16 nennt `SPEC-029` Punkt 2 „Berichtigung des Numerus, keine Aussage geändert“; `f4e9defb` fügt dort die Aussage zum erwarteten Verbindungsende beim Einspielen hinzu, die keine Historie-Zeile trägt. Zuständig: Implementer. | Maintainability | `spec/spezifikation.md` · „keine Aussage geändert“ | nein — kein Gate liest die Historie | Historie nennt neuen Satz der Spezifikation nicht |
| F-5 | INFO | Der Absatz in `LH-FA-20.a` („Ein Test je Aussage prüft sie“) steht seit `179ab424` (Architect, vor dem Code), `TestPlayEndeHilfe` erst in `cc4ca898`. Im Endstand ist der Satz wahr (acht Teilfälle, je Aussage einer); im Zwischenstand sagte die Spezifikation einen Test zu, den es nicht gab. Keine Aktion, wenn der Architect-Commit vor dem Code Absicht war. | `AGENTS.md` §3.11 | `spec/spezifikation.md` · „Ein Test je Aussage prüft sie.“ | nein | Zusage der Spezifikation vor ihrem Test |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driving/cli/cli.go` | geprüft, ohne Befund: nur die Konstante `optionenPlay` geändert (Diff), Wortlaut Zeichen für Zeichen der des Architect in §6, die erste Aussage von `--allow-recorded-errors` unverändert, höchstens 78 Zeichen je Zeile |
| `internal/adapters/driving/cli/play_ende_hilfe_test.go` | geprüft, ohne Befund außer F-1: acht Needles wie in §6, Block je Option, Kommentar trägt nur Zusage und Abgrenzung, die zwölf Mutanten aus §7 nachgefahren und rot |
| `test/integration/play_ende_vergleich_e2e_test.go`, Fixtures | geprüft, ohne Befund außer F-2 und F-3: Zähler `3/3/0` und `2/2/1`, Meldung mit Stelle, Position, Art und beiden SQLSTATE (volle Meldung, kein Teilstring eines Zufallspfads), `to_regclass` als zweite Sicht; M1 bis M5 wie in §7 |
| `test/integration/play_vergleich_e2e_test.go` | geprüft, ohne Befund: nur die Abdeckungs-Deklaration von `TestE2EPlayVergleichFatal` enger gefasst, Code und Erwartung per Diff unberührt |
| Produktcode außerhalb des Hilfetextes (`internal/hexagon`, Adapter, Bootstrap) | geprüft, ohne Befund: im Diff keine Änderung |
| `docs/user/abdeckung-e2e.md`, `docs/user/abdeckung-unit.md` | geprüft, ohne Befund: Deklarationen sagen nur zu, was die Tests prüfen; „handgeschriebene Aufzeichnung“ genannt, Szenario 18 nicht behauptet |
| `spec/spezifikation.md` | geprüft, ohne Befund außer F-4 und F-5: der Satz in `SPEC-029` Punkt 2 steht erst im Commit der E2E-Tests, die Berichtigung des Numerus nicht doppelt, kein Artefakt unterhalb genannt |
| Plan (`§1`, `§3`, `§6`, `§7`, `§8`) und Nehmer `slice-v1-abschluss-verbindungsende-upstream` | geprüft, ohne Befund außer F-3: Plan folgt dem Diff, Messung 222 stimmt, der Nehmer führt die Prüfung der Berichtigung ohne neue Last, keine neue Sendung |
| Gates `make gates` | nicht erneut gefahren (Gegenstand des Verifiers); die Stichproben liefen mit `make test-integration` im Wegwerf-Commit, der Arbeitsbaum ist wieder auf `8f90c592` |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Hilfetext-Test prüft Anwesenheit, nicht Zuordnung · Beleg-Aussage enger als die Prüfung · Fixture-Form erst im Lauf gegen den Server entschieden · Historie nennt neuen Satz der Spezifikation nicht · Zusage der Spezifikation vor ihrem Test

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. F-1 ist LOW mit Anker `AGENTS.md` §3.10 und lässt sich im Review-Nachgang oder mit dem Doku-Folge-Slice schließen; eine Entscheidung dagegen gehört in §7.

**Übergabe:** Findings gehen an den Implementer (F-1, F-2, F-4), F-3 und F-5 sind Hinweise ohne Aktion; die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
