# Review-Report: slice-harness-gegenproben-schneller (r2) — 2026-10-25

**Review-Art:** Code — Folgereview (r2) gegen Plan, Entscheidungen und Hard Rules; geprüft wird die Nacharbeit zu den Findings F-1 bis F-8 des ersten Reviews.

**Gegenstand:** `slice-harness-gegenproben-schneller`, Zweig `impl-gps1`, Diff `b60a6b9d..ee68abcd`, Schwerpunkt `e1223cfe..ee68abcd` (Nacharbeit); erster Review: `docs/reviews/2026-10-25-review-slice-harness-gegenproben-schneller.md` (`bccd9a65`)

**Skill:** `.harness/skills/reviewer.md` @ `ee68abcd`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-25

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein Ausfüll-Hinweis)*. Dieser
> Report friert ein; was er zitiert, bewegt sich weiter. Deshalb: **Kennung, nicht Adresse** —
> `slice-<Kennung>` statt seines Lifecycle-Pfads, `make <target>` statt eines Links auf die
> Sensor-Datei, eine Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-harness-gegenproben-schneller` (§1, §3, §6 Zeilen (g) bis (j), §7 Belege des Implementers)
- [ADR-0059](../plan/adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md), [ADR-0055](../plan/adr/0055-gegenproben-haben-ein-zeitbudget.md) (Punkt 1 ersetzt durch ADR-0059)
- `SPEC-060` Punkt 16 (4), (5), (9) und Punkt 17 (Stand des Architect `e1223cfe`), `SPEC-070` Punkt 10
- Erster Review-Report (F-1 bis F-13)
- `AGENTS.md` (Hard Rules, insbesondere §3.7, §3.9, §3.10, §3.11)

---

## Findings

Aufträge (a) bis (f) des Orchestrators sind unten in den Negativbefunden einzeln beantwortet. Die vier Regex-Mutationen der Zahlform, die Mutationen zu den Nachkommen und die geänderten Erwartungslisten habe ich mit `GP_MUT` gefahren (siehe Negativbefunde).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | `SPEC-060` Punkt 16 (4) zählt als ungültige Formen auch „Vorzeichen“ und „Leerraum“ auf; die Fälle decken vom Vorzeichen nur `-1` (kein `+5`) und von Leerraum keine Form (` 5`, `5 `). Eine Mutation, die genau eine davon zulässt (etwa `^[+ ]?(0\|[1-9][0-9]{0,5})$`), bliebe in allen drei Skripten grün; ohne Failure-Szenario (die Form hat keine Folge außer der abweichenden Meldung). | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `tools/dev/harness/release-gegenprobe.sh` · `f_p14_budget_ungueltig_abc() { budget_ungueltig 'abc'; }` (gleiche Lage in `container-gegenprobe.sh` `'abc:abc'` und `container-einrichten-gegenprobe.sh` `abc:abc`) | ja — Mutation `^[+ ]?…` in einer Kopie, `make release-gegenprobe` mit `GP_MUT` | Aufzählung der Zusage ohne Fall je Element |
| F-2 | INFO | Die neue Prüfung in `abbruchfall` schließt Überlebende per `kill -0` aus. Ein Zombie, den ein Init ohne Reaping (Container ohne `--init`) nicht einsammelt, würde als lebend gelesen; die frühere Prüfung (`pgrep -f` auf den Pfad) sah Zombies nicht. Die Prüfung nimmt außerdem nicht an, dass die Menge zur Zeit des Signals mehr als den Sitzungsführer enthält; ein leerer Fang färbte nichts. Im Lauf mit beiden Mutationen sind beide Fälle rot, also trägt der Fall hier. | Maintainability | `tools/dev/harness/release-gegenprobe.sh` · `lebende="$(pgrep -s "$pid" 2> /dev/null \| tr '\n' ' ')"` | nein — hängt am Init des Gegenproben-Containers | Prüfung auf Überlebende hängt an Prozess-Umgebung |
| F-3 | INFO | §7 des Plans trägt für den Endstand `ee68abcd` keine Zeile zum vollen Lauf (Befehl, Hash, Dauer, Last, Ausgang): die letzten Einträge nennen `76456d59` mit drei roten Befunden und „nicht wiederholt“. Die grünen Läufe (`release-gegenprobe` 732 Läufe in 982 s, Last 8,06; `test-container-gegenprobe` 533 s, Last 3,09) liegen nur als Logs im Scratchpad. Ebenso ist `make gates` nicht am Endstand belegt (312 s am Stand `76456d59`, danach änderten sich zwei Skripte). | `AGENTS.md` §6 Punkt 6; Plan §1 *Übergang für die Rollen* („belegt ihn in §7“) | `docs/plan/planning/in-progress/slice-harness-gegenproben-schneller.md` · `*Voller Lauf am Stand `76456d59`*` | nein — Einlage des Implementers in §7 | Beleg des Endstands nicht in §7 |
| F-4 | INFO | Der Wortlaut-Vorschlag für `AGENTS.md` §6 Punkt 6 und `implement-slice.md` Schritt 15 steht noch nicht in §7; der DoD-Punkt *Vorschlag liegt vor* ist offen und §1 verweist auf „Wortlaut in §7“. Das ist kein Mangel der Nacharbeit (F-7 verlangte den Punkt, nicht den Text), nur der noch offene Teil. | Maintainability | `docs/plan/planning/in-progress/slice-harness-gegenproben-schneller.md` · `- [ ] Der Vorschlag zum Wortlaut für `AGENTS.md` §6 Punkt 6` | nein | Offener DoD-Punkt (an den Verifier) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| (a) F-1 erster Review: `GP_BUDGET`-Zahlform `^(0\|[1-9][0-9]{0,5})$` in `release-`, `container-` und `container-einrichten-gegenprobe` | geprüft, ohne Befund außer F-1 oben: Regex in allen drei Skripten gleich; Fälle `08`, `09`, `00`, `007`, `1000000`, `99999999999999999999` (je Exit 2, genau eine Zeile) und grün `999999`, `10`; vier Mutationen der Punkt-16-Liste (`^[0-9]{1,6}$`, `^(0\|[1-9][0-9]*)$`, `^[1-9][0-9]{0,5}$`, `^(0\|[1-9][0-9]{0,4})$`) und dazu `10` abgelehnt, `ohne-exit` und `vor-faellen` mit den zehn ungültigen Fällen. Gefahren mit `GP_MUT`: `m-p12-budget-fuehrende-null-erlaubt` und `-gruende-vertauscht` und die Erwartungs-Mutationen (`container-gegenprobe`), `m-p14-budget-stellen-unbegrenzt` (`release-gegenprobe`): grün, je Mutation genau die benannten Fälle. Reihenfolge der Gründe (a) vor (b) vor (c) vor (d) bleibt wie in `SPEC-060` Punkt 16 (4). |
| (b) F-2 erster Review: `>=` | geprüft, ohne Befund: `SPEC-060` Punkt 16 (9) führt `>=` als nicht äquivalent und als akzeptiertes Negativ mit Grund (Zeit nicht steuerbar, eine Sekunde); §7 des Plans folgt ihm; im Code steht kein Satz „äquivalent“ zu `>=`. |
| (c) F-3 erster Review: `abbruchfall` | geprüft, ohne Befund: `pgrep -s` vor dem Signal fängt die Sitzung der Kopie (der Fall läuft, sobald `fall_laeuft` steht); die Kopie wartet in `beende_kinder` per `kill -0` auf ihre Nachkommen, danach prüft `abbruchfall` dieselben Prozesse; die Mutationen `m-p13-gegenprobe-int-kinder-bleiben` und `-term-kinder-bleiben` lassen die Kinder leben und färben genau ihren Fall (gefahren: grün, 2 Läufe, 16 s). Das trägt über die Prozesse selbst, nicht über den Ordner. Der Kommentar vor `[ -d "$run" ] \|\| mkdir -p "$run"` nennt die Kopplung/Abgrenzung (Klasse Abgrenzung, §3.7) im Indikativ. Grenzen: F-2. |
| (d) F-4, F-5, F-7, F-8 erster Review | geprüft, ohne Befund: README-Zeilen nennen nur noch, was ein Fall prüft (Ordner je Mutation, `t.<Fall>` der Fälle `p1` je Fall, Zahlform; `GP_INNEN` nicht mehr genannt, §3.11); `m-p12-budget-gruende-vertauscht` färbt nur `p12-budget-gruende`; Plan §1 und §3 nachgezogen, DoD-Punkt ergänzt (offen, F-4); die zwei Kommentare der Skripte beschreiben den Ist-Zustand; Hauptzeile von ADR-0059 im ADR-Index vorhanden (F-6 erster Review); „bis zu 25 min“ in `SPEC-070` Punkt 10 berichtigt (F-13 erster Review). |
| (e) Erwartungslisten `m-p14-`/`m-p12-budget-zeile-mit-obergrenze`, `-last-null`, `-last-nie-gelesen` und `m-p17-…` | geprüft, ohne Befund: die Fälle `…-zulaessig-999999` und `…-10` prüfen dieselbe Zeile `Zeit <T> s von höchstens <B> s, Last <L>` wie `…-unter`. Fällt die Zeile aus, bleibt `Last` bei 0 oder ungelesen, sind sie notwendig rot; die Erwartung wurde nach der Ursache ergänzt, nicht um Grün zu erzwingen. Die Gegenprobe der Mutationen (`GP_MUT`, gefahren) nennt je Mutation genau die Fälle, mit den neuen Namen grün. Die Fälle von `…-null`, `-ueber` und `-last-fehlt` sind nicht über-erwartet. |
| (f) Belege | geprüft: Logs `gps-voll-release.log` (732 Läufe, 982 s, Last 8,06, grün) und `gps-voll-container.log` (533 s, Last 3,09, grün) liegen im Scratchpad vor und stimmen mit den Angaben des Auftrags überein; die Eintragung in §7 und `make gates` am Endstand fehlen (F-3). Der Lauf liegt über der Last-Grenze von `SPEC-060` Punkt 16 (1) (Host nicht ruhig); der Plan nennt das. |
| `harness/vertraege/SPEC-060`, `SPEC-070` Punkt 10, `harness/README.md` (Nacharbeit) | geprüft, ohne Befund |
| `docs/plan/adr/` (ADR-0059 Index, ADR-0055 unverändert) | geprüft, ohne Befund |
| Produkt-Code (`cmd/`, `internal/`), `Dockerfile`, `.claude/`, `AGENTS.md` | nicht im Diff, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Aufzählung der Zusage ohne Fall je Element · Prüfung auf Überlebende hängt an Prozess-Umgebung · Beleg des Endstands nicht in §7 · Offener DoD-Punkt

Alle Findings F-1 bis F-8 des ersten Reviews sind behoben (F-1 bis F-3 MEDIUM tragen, F-4 bis F-8 LOW erledigt); neu sind ein LOW und drei INFO, keines blockierend.

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. F-1 (LOW) darf der Implementer annehmen oder mit Grund ablehnen; F-3 und F-4 (INFO) sind Aufgaben vor der Closure, keine Änderung des Diffs.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg** und wird über Läufe hinweg nicht wieder gelesen. Er ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
