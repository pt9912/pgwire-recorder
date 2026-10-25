# Review-Report: slice-harness-gegenproben-schneller — 2026-10-25

**Review-Art:** Code — Code-Review gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; das ist die Frage des Verifiers).

**Gegenstand:** `slice-harness-gegenproben-schneller`, Zweig `impl-gps1`, Diff `b60a6b9d..d57a1b3a` (14 Dateien, +549/−99; gemessen nach `AGENTS.md` §3.14: 494 Zeilen, Marke 1500 nicht erreicht).

**Skill:** `.harness/skills/reviewer.md` @ 4155023d · <!-- d-check:ignore (Adopter-spezifischer Skill-Pfad) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-25

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

- Plan `slice-harness-gegenproben-schneller` (§1 bis §8, §7 Belege des Implementers, Läufe 1 bis 4)
- `SPEC-060` Punkt 16 (inkl. Randformen (g), (h), (i) und Entscheidung des Architect in `36bab2f2`), `SPEC-055` Punkt 12, `SPEC-069` Punkt 17, `SPEC-070` Punkt 10
- [ADR-0055](../plan/adr/0055-gegenproben-haben-ein-zeitbudget.md), [ADR-0059](../plan/adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md), [ADR-0058](../plan/adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) (Bezug)
- keine `LH-*`-Kennung (Harness-Werkzeug)
- `AGENTS.md` (Hard Rules, insbesondere §3.2, §3.5, §3.7, §3.8, §3.9, §3.10, §3.11, §3.12, §3.13, §3.14)

**Eigene Sonden dieses Laufs** (nur lesen und engste Fälle, kein Docker, keine vollen Läufe, kein `make gates`; Kopien der Skripte im Scratchpad, am Arbeitsbaum nichts geändert):

- S1: `bash`-Arithmetik der Prüfung `(( T > budget ))` mit `budget` ∈ {`08`, `09`, `00`, `0010`, `99999999999999999999`} (alle bestehen `^[0-9]+$`).
- S2: Kopie von `release-gegenprobe.sh` mit `T=5` statt der Messung und `GP_BUDGET=5`, `GP_MUT='^fall-ohne-funktion$'`, einmal Original, einmal mit `>=` statt `>`.
- S3: `GP_MUT='^m-p13-gegenprobe-(int|term)-kinder'` an drei Ständen: HEAD, HEAD ohne die Zeile `[ -d "$run" ] || mkdir -p "$run"` in `faelle`, und Stand `5ebede1a` (vor der Behebung).
- S4: 3000 Zufallseingaben, `${s#*"$a"}` gegen `${s/"$a"/<Marke>}` als Eindeutigkeitsprüfung.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | `GP_BUDGET` mit führender Null oder über 64 Bit besteht die Prüfung `^[0-9]+$`, aber die Arithmetik versagt still: `08` und `09` melden `((: Wert zu groß für die Basis` und lassen die Prüfung aus (Exit 0 trotz Überschreitung), `99999999999999999999` ebenso, `00` und `0010` werden als 0 bzw. 8 gelesen, und `00` ist rot bei jeder Zeit über 0 (`00 != 0`). Dieselbe Prüfung steht in allen drei Skripten; keine Randform in §6, kein Fall, die Zusage „ungültiger Wert Exit 2“ nennt nur Leerstring, `-1`, `1.5`, `abc`. Szenario: Ein Aufrufer setzt `GP_BUDGET=0900`, die Überschreitung bleibt ohne Rot, die Zeile `Zeit … von höchstens 0900 s` behauptet eine Prüfung. | `SPEC-060` Punkt 16 (4); `AGENTS.md` §3.12 (Randform nicht vor dem Code entschieden), §3.10 | `tools/dev/harness/release-gegenprobe.sh` · `[[ $GP_BUDGET =~ ^[0-9]+$ ]]` und `(( T > budget ))`; ebenso `container-gegenprobe.sh`, `container-einrichten-gegenprobe.sh` | ja — Sonde S1; ein Fall `GP_BUDGET=08` in `p14-*`/`p12-*` würde es zeigen | Randform von GP_BUDGET (Zahlform) nicht entschieden, Arithmetik statt Textvergleich |
| F-2 | MEDIUM | §7 führt den grünen Mutanten `>=` statt `>` als „äquivalent in der Praxis, akzeptiertes Negativ“ mit dem Grund, die Zeit sei nicht steuerbar; Invariante, Ort und Sonde fehlen. Die Zeit ist steuerbar: Die Ersetzung `T=5` in der Kopie (Mechanismus `ZE`/`zk`, den die Fälle `p14-*` ohnehin nutzen) mit `GP_BUDGET=5` trennt Original (Exit 0, `Zeit 5 s von höchstens 5 s`) und Mutant (Exit 1, `ROT — Zeitbudget: 5 s über 5 s`). Die Zusage „höchstens `B` s“ (gleich ist grün) bricht, ohne dass ein Test rot wird; der Mutant gilt in allen drei Gegenproben. Das Negativ ist in `SPEC-060` Punkt 16 (9) vom Architect als akzeptiert festgehalten; die Sonde widerlegt dessen Begründung, daher geht der Punkt an den Architect. | `AGENTS.md` §3.10 (grüner Mutant gilt als äquivalent erst mit Invariante, Ort und Sonde) | `harness/vertraege/SPEC-060-gegenproben.md` · „Die Mutation zur Nutzung von `>=` statt `>` bleibt grün“; Plan §7 · „`>=` statt `>` bei der Prüfung: äquivalent in der Praxis“ | ja — Sonde S2 | Grüner Mutant ohne Invariante/Ort/Sonde als äquivalent geführt |
| F-3 | MEDIUM | Die Behebung `7b8387a1` (`[ -d "$run" ] \|\| mkdir -p "$run"` je Fall in `faelle`) trägt keine Mutation verlässlich, und der Grund steht nirgends im Skript: Die Fälle `p13-gegenprobe-sigint` und `-sigterm` erkennen überlebende Nachkommen daran, dass diese den Ordner neu anlegen (Kopplung zwischen `faelle` und `abbruchfall`/`reste_da`). Sonde S3: am Stand `5ebede1a` färbt keine der beiden Mutationen `m-p13-gegenprobe-{int,term}-kinder-bleiben` ihren Fall (2 von 2 Läufen „nichts“); mit HEAD ohne die Zeile war es in 3 von 4 Läufen rot, im vierten grün (zeitabhängig); HEAD 3 von 3 grün. Szenario: Wer die Zeile als überflüssig entfernt (der Bestand hatte sie schon einmal je Fall), bekommt nur mit Glück Rot; die Zusage „Nachkommen werden beendet“ bricht dann ohne Test. | `AGENTS.md` §3.10 (Mutation je Zusage), §3.7 (Kopplung ohne Kommentar) | `tools/dev/harness/release-gegenprobe.sh` · `[ -d "$run" ] || mkdir -p "$run"` in `faelle`; `reste_da() { [ -n "$(ls -A "$t"` | ja — Sonde S3 | Fall erkennt Nachkommen über Nebenwirkung statt Beobachtung; Behebung ohne tragende Mutation |
| F-4 | LOW | Die Zusage „Ordner jedes Falls entfällt, sobald der Fall endet“ ist nur für `t.<Fall>` der Gruppe `p1` geprüft (`p14-ordner-fall` fährt `FAELLE_P1_KLEIN`). Die Elemente `t.<Fall>.*`, `t.w`, `t.w.*` und `tmp.*` derselben Entfernzeile (der einzige Ordner der Gruppen `b`, `f`, `h`) hat keine Mutation, die nur sie bricht; gelesen, nicht gefahren. Die Zeile in `harness/README.md` („Ordner je Mutation und je Fall entfernt“) sagt über die Gruppen `b`, `f`, `h` mehr zu, als ein Fall prüft. | `AGENTS.md` §3.10 (je Element), §3.11 | `tools/dev/harness/release-gegenprobe.sh` · `for d in "$run/t.$n" "$run/t.$n".* "$run/t.w"` und `Z_P1_KLEIN`; `harness/README.md` · „Ordner je Mutation und je Fall entfernt“ | nein — Lesen; ein Mutant, der nur `t.w` fallen lässt, würde es zeigen | Entfernen je Element nicht je Element gefangen |
| F-5 | LOW | In `container-gegenprobe.sh` fehlt die Mutation zur Reihenfolge der Gründe (c) vor (d); der Fall `p12-budget-gruende` steht, Plan §7 begründet das Fehlen damit, dass die roten Zeilen im Lauf entstehen. `SPEC-060` Punkt 16 (9) verlangt je Element einen Fall und eine Mutation, die nur ihn bricht; `release-gegenprobe.sh` hat `m-p14-budget-gruende-vertauscht`. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `tools/dev/harness/container-gegenprobe.sh` · `ist p12-budget-gruende`; Plan §7 · „ohne Mutation zur Reihenfolge der Gründe“ | nein — Lesen | Fehlende Mutation zu einem Element der Zusage |
| F-6 | LOW | Die Hauptzeile von [ADR-0059](../plan/adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md) fehlt in der Tabelle `ID · Titel · Status · Bezug` des ADR-Index; nur die Relationstabelle („ergänzende und teilweise ersetzende ADRs“) führt sie. Die Zeilen 0055 bis 0058 stehen in beiden. | `AGENTS.md` §5 Regel 3 (neue ADRs im Index), ADR-Index „Bei `Accepted`: diesen Index aktualisieren“ | `docs/plan/adr/README.md` · Tabelle `\| ID \| Titel \| Status \| Bezug \|` endet bei `[0058]` | ja — Lesen des Index; ein Sensor dafür liegt nicht vor | ADR-Index unvollständig geführt |
| F-7 | LOW | Der Plan folgt zwei Korrekturen nicht: (1) §1 *Übergang für die Rollen* sagt „Die Zeile je Gegenprobe in `harness/README.md` §Sensors bleibt bis zur Lieferung unverändert“, die Zeilen sind seit `5ebede1a` geändert. (2) §3 (Tabelle) trägt „`AGENTS.md` §6 Punkt 6 und `implement-slice.md` Schritt 15: Übergangssatz mit der Lieferung herausnehmen“, der Diff berührt beide Dateien nicht, und die DoD führt den Punkt nicht; ebenso nennt §3 die Tabelle von `SPEC-060` Punkt 16, die DoD hat dafür eine offene Zeile. Der Punkt in (2) kann beim Closure verloren gehen. | `AGENTS.md` §3.9 | Plan · „bleibt bis zur Lieferung unverändert (§3.11: sie sagt nur zu“; Plan · „nimmt den Übergangssatz zu `GP_MUT` und dem vollen Lauf des Verifiers heraus“ | nein — Lesen | Plan folgt Korrektur nicht |
| F-8 | LOW | Ein Kommentar im Code beschreibt die verworfene Alternative: „`${neu#*"$alt"}` ist bei 260 KB quadratisch“ steht an der Eindeutigkeitsprüfung in `mutr`; „(und ist rot, wie zuvor)“ in der Kommentarzeile zu `snap_hole` verweist auf abwesenes Verhalten. Kein Failure-Szenario (ein Lauf handelt nicht danach). | `AGENTS.md` §3.7 | `tools/dev/harness/release-gegenprobe.sh` · `# eindeutig: nach dem Ersatz der ersten Stelle trägt der Rest keine zweite`; `# … jeder Fall fährt ihn selbst (und ist rot, wie zuvor)` | nein | Kommentar trägt verworfene Alternative |
| F-9 | INFO | Ergebnisgleichheit (Auftrag a): Gelesen und kein Befund. Der Stand je Mutation (`snap_hole`/`snap_sichere`) wird nur bei `bad=0` des ersten Falls abgelegt; ein roter gemeinsamer Schritt legt nichts ab, und jeder weitere Fall fährt ihn selbst und ist rot. Cache `gp-layout` (je Welt, Schlüssel Variante), `hsumme` (reine Funktion von Plattform und Version), `tap_leer` (der einzige Zweitaufruf löscht vorher), `git init --template=` und `#!bash` in den Attrappen lockern keinen Fall. Rest: Nach `snap_hole` setzt kein Schritt `code`, `out`, `err` neu; ein Fall, der nach `bild` ohne `lauf` einen davon liest, sähe den Wert des Vorgängers — der volle Lauf (727 grün) zeigt keinen. Die Gleichheit ruht auf den unveränderten Deklarationszeilen und dem grünen Lauf; ein Vorher-Mengenvergleich liegt nur für `m-p19-` vor (Plan §7 Lauf 3), für `sh`, `b`, `f`, `p1` nur nachher. | `SPEC-060` Punkt 16 (9) | Plan §7 · „Ein Vorher-Lauf aller Gruppen fehlt“ | nein | — |
| F-10 | INFO | Eindeutigkeitsprüfung (Auftrag c): gleichwertig. `${x/"$a"/<Marke>}` ersetzt das erste Vorkommen; ein weiteres Vorkommen im Rest wird gefunden, überlappende Treffer trennt die Marke wie `${x#*"$a"}` den Rest abschneidet; Vorbedingung (`a` enthält die Marke nicht) gilt. Sonde S4: 0 Abweichungen in 3000 Eingaben. Die Prüfung selbst (Meldung „nicht eindeutig“) hat weiter keinen eigenen Fall (Bestand). | Maintainability | `tools/dev/harness/release-gegenprobe.sh` · `rest="${neu/"$alt"/$'\001'}"` | ja — Sonde S4 | — |
| F-11 | INFO | Test-Haken (Auftrag d): `GP_INNEN`, `GP_FALL`-Semantik (Protokoll und `stderr` der unveränderten Kopie vor dem Entfernen, Auswahl unverändert, keine Vorgabe), Vorgabe nur ohne `GP_MUT` und `GP_FALL`, Geltung bei `GP_ENTDECKE`, Reihenfolge der Gründe in `release-gegenprobe.sh` stimmen mit `SPEC-060` Punkt 16 (4), (5), (7), (9) und den Randformen (g) bis (i) überein; Mutationen je Element sind vorhanden (Ausnahmen: F-2, F-5). `zk_start` in `container-gegenprobe.sh` prüft die Textstellen nicht wie `zk` (ein weggefallener Anker ließe den Fall mit Signal ohne Verzögerung laufen; die Mutationen `m-p12-budget-sig*-zeile` fangen den Bruch der Zusage trotzdem). `zs` wartet mit festen `sleep 0.5` nach der Beobachtung der Attrappe. | `SPEC-060` Punkt 16 | `tools/dev/harness/container-gegenprobe.sh` · `zk_start() {` | nein | — |
| F-12 | INFO | [ADR-0059](../plan/adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md) formal (Auftrag f): `Ersetzt` nennt nur Punkt 1 von [ADR-0055](../plan/adr/0055-gegenproben-haben-ein-zeitbudget.md), ADR-0055 selbst ist im Diff unverändert (§3.5 gewahrt); Regeldetails stehen mit Verweis in `SPEC-060` (§3.8 gewahrt). Der Marker `d-check:status-provenance` in der Kontext-Zeile ist der Ausweg, den `.d-check.yml` für eine Accepted-ADR mit Slice-Kennung als Provenance nennt: angemessen. Er kam durch den Implementer in einem späteren Commit in eine angenommene ADR (nur Kommentar, Text unverändert). „Hochrechnung 985 s trägt mit Abstand“: die vollen Läufe lagen bei 964 s und 1044 s (80 und 87 Prozent von 1200 s) auf einem Host, der nicht ruhig im Sinne von `SPEC-060` Punkt 16 (1) war; der Re-Evaluierungs-Trigger der ADR deckt das. [ADR-0057](../plan/adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md) und [ADR-0058](../plan/adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) nennen weiter 600 s (angenommen, unveränderlich; ADR-0059 `Ersetzt` führt sie nicht). | `AGENTS.md` §3.5, §3.8 | `docs/plan/adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md` · `<!-- d-check:status-provenance -->` | nein | — |
| F-13 | INFO | Nachzüge der Zahl 1200 (Auftrag g): vollständig in `SPEC-060` Punkt 16, `SPEC-055` Punkt 12, `SPEC-069` Punkt 17, `SPEC-070` (Kopf, Punkt 10, 14), `harness/README.md`, `harness/mk/container-einrichten-gegenprobe.mk` und den drei Skripten (`grep` nach `budget=600`/`höchstens 600` im Code leer). Rechnung in `SPEC-070` Punkt 10 (b): „im Regelfall unter 25 min (rund 300 s und bis zu 1200 s)“ ergibt genau 25 min, nicht darunter; Sache des Architect. **Noch 600 s oder „zehn Minuten“ (Sache des Planner, nur gemeldet):** `slice-harness-gegenproben-grenze-rest` (Titel, Ziel, DoD „höchstens 600 s“, Risiken), `slice-harness-rollenlauf-zeitgrenze`, `slice-harness-rollenlauf-zeitgrenze-haertung`, `slice-harness-rollenlauf-wiederholungsgrenze` (DoD „unter 600 s“), `slice-harness-container-start-u-boot-0-13`, `CO-001` und `CO-003` (Titel „unter zehn Minuten“, Trigger „600 s“), das Register `BEO-REPO/rollenlauf-ohne-zeitgrenze`, die Roadmap-Zeile vom 2026-10-21; dazu in `SPEC-060` Punkt 16 der Satz „Stand 2026-10-21 noch nicht angelegt“ zur Kennung des Nehmers. | `AGENTS.md` §3.9 | `harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md` · „Beide zusammen belegen im Regelfall unter 25 min“ | ja — `grep` nach `600 s` | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/plan/adr/0059-*` Inhalt (§3.5, §3.8), `docs/plan/adr/0055-*` unverändert | geprüft, ohne Befund (Hauptzeile im Index: F-6) |
| `harness/vertraege/SPEC-055`, `SPEC-069`, `SPEC-070` (Zahl 1200, Verweise auf ADR-0059) | geprüft, ohne Befund außer dem Rechenhinweis in F-13 |
| `harness/vertraege/SPEC-060` (Punkt 16 (1), (4), (5), (7), (9) Ergänzungen, keine Zahl 600 mehr) | geprüft, ohne Befund |
| `harness/README.md` (Zeilen `test-container-gegenprobe`, `release-gegenprobe`): Zusagen gegen Fälle | geprüft; eine Zusage weiter als ihr Test (F-4), sonst ohne Befund |
| `harness/mk/container-einrichten-gegenprobe.mk`, `tools/dev/harness/container-einrichten-gegenprobe.sh` (Vorgabe 1200, Fall und Mutation) | geprüft, ohne Befund (F-1 gilt auch hier) |
| `tools/dev/harness/release-gegenprobe.sh`: Snapshots, Caches, Attrappen, Zeitbudget, Entfernen je Mutation | geprüft; F-1 bis F-5, F-8, F-9 |
| `tools/dev/harness/container-gegenprobe.sh`: Zeitbudget, Baum je Mutation, Fälle `p12-*` | geprüft; F-1, F-2, F-5, F-11 |
| Rückführungs-Marke (`AGENTS.md` §3.14) | geprüft, ohne Befund: 494 von 1500 Zeilen |
| `docs/plan/planning/done/welle-v1-nachtrag/` (zwei nachgezogene Verweise) | geprüft, ohne Befund (reine Zeigerkorrektur auf `in-progress/`) |
| Produkt-Code (`cmd/`, `internal/`), `Dockerfile`, `.claude/`, `AGENTS.md` | nicht im Diff, ohne Befund (siehe F-7 zur ausstehenden Änderung an `AGENTS.md`) |
| Suppression (`//nolint`, `#noqa`) | geprüft, ohne Befund; keine neue Suppression im Diff |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 5 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Randform von GP_BUDGET (Zahlform) nicht entschieden, Arithmetik statt Textvergleich · Grüner Mutant ohne Invariante/Ort/Sonde als äquivalent geführt · Fall erkennt Nachkommen über Nebenwirkung statt Beobachtung; Behebung ohne tragende Mutation · Entfernen je Element nicht je Element gefangen · Fehlende Mutation zu einem Element der Zusage · ADR-Index unvollständig geführt · Plan folgt Korrektur nicht · Kommentar trägt verworfene Alternative

## Verdikt

**Merge-blockierend:** ja — drei MEDIUM (F-1, F-2, F-3), je mit reproduzierter Sonde. F-2 betrifft ein vom Architect in `SPEC-060` Punkt 16 (9) als akzeptiert festgehaltenes Negativ; die Sonde widerlegt dessen Begründung (Zeit steuerbar), der Implementer ändert die Zusage nicht selbst, der Punkt geht an den Architect (Fall und Mutation oder neue Begründung mit Invariante, Ort und Sonde). F-1 verlangt eine Entscheidung der Randform (Zahlform von `GP_BUDGET`) im Vertrag vor dem Code; sie gehört dem Architect.

Die Umbauten für das Tempo (Auftrag a, c, d) tragen nach Lesen und Sonden: kein gelockerter Fall gefunden (F-9, F-10, F-11). Kein HIGH.

**Übergabe:** Findings gehen an den Implementer (F-1 und F-2 über den Architect, weil sie Vertragstext berühren); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und muss es nicht. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext). Die nicht geprüften DoD-Punkte (Platz-Spitze unter 100 MB, Fahrt unter fremdem Benutzer, fremder Locale und Last, Tabelle in `SPEC-060`) liegen beim Verifier.
