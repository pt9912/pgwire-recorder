# Verifikationsbericht: slice-harness-gegenproben-schneller — 2026-10-25

**Rolle:** Verifier (Modul 11), Frage „Bauen wir es richtig?“ gegen Plan, DoD und Entscheidungen.
**Gegenstand:** `slice-harness-gegenproben-schneller`, Zweig `impl-gps1`, Diff `b60a6b9d..6cb151da`.
**Maßstab:** DoD im Plan, `SPEC-060` Punkt 16 mit den Randformen (g) bis (j) (Architect `36bab2f2`, `e1223cfe`), [ADR-0059](../plan/adr/0059-obergrenze-der-gegenprobe-zwanzig-minuten.md), [ADR-0055](../plan/adr/0055-gegenproben-haben-ein-zeitbudget.md), Reviews r1 und r2.
**Eingang:** Belege des Implementers in §7; die Logs des vollen Laufs am Endstand `ee68abcd` (vom Koordinator gefahren) habe ich gelesen, nicht neu gefahren.

## Ergebnis: bestanden mit Vorbehalten

Alle Liefer-Punkte (Tempo und Ergebnisgleichheit, Zeitbudget und Ordner je Mutation, Fälle zu Punkt 16 (9)) sind am Code bestätigt, und die Belege tragen. Offen sind die Closure-Punkte der DoD (Tabelle in `SPEC-060`, Platz-Spitze im Beleg, Fahrt unter fremder Umgebung) mit Adressen unten. Kein Befund bricht eine Zusage.

## Eigene Sensoren (nur engste Fälle, kein voller Lauf, kein Docker parallel)

| Sensor | Ergebnis |
|---|---|
| `make docs-check` | 0 Befunde (1488 Dateien) |
| `make kopf-check`, `make abdeckung-check` | Exit 0 |
| Stichprobe `GP_PARALLEL=1`, drei Mutationen (`m-p16-pin-ohne-digest`, `m-p2-trockenlauf-image-ruft-skopeo`, `m-p2-trockenlauf-ruft-skopeo`) | grün, 29 s, Last 2,5 (Ziel höchstens 37 s; vorher 73 s). Nur `m-p16-pin-ohne-digest`: 13 s |
| Dasselbe unter `LC_ALL=C TZ=Asia/Tokyo` (nur `m-p16-pin-ohne-digest`) | grün, 13 s |
| `GP_MUT='^m-p14-' GP_ENTDECKE=1 GP_PARALLEL=4` | `grün (30 Läufe)`, 164 s; jede der 30 Mutationen färbt genau die deklarierten Fälle |
| `GP_MUT='^m-p13-gegenprobe'` (Nachkommen nach Abbruch) | `grün (7 Läufe)`, 39 s |
| `GP_MUT='^ok$'` (unveränderte Kopie mit allen Fällen, auch den neuen `+5`, ` 5`, `5 `) | grün, 156 s |
| `m-p3-3-version-datei-entfernt`, `m-p5-summe-name-os` (Snapshot b/f) | grün; `GP_ENTDECKE` zeigt 65 bzw. 56 rote Fälle: Ein roter gemeinsamer Schritt färbt jeden abhängigen Fall, der Stand je Mutation maskiert nichts |
| `test-container-gegenprobe`, `GP_MUT='^m-p12-(budget-gruende-vertauscht\|ordner\|budget-ungueltig-ohne-exit)'` | grün, 68 s |
| Platz-Spitze (`du -sk` je Sekunde auf `TMPDIR`): `m-p3-` und `m-p5-` (107 Mutationen, `GP_PARALLEL=4`) und `m-p19-push-` | 6,5 MB bzw. 5,2 MB; Ordner danach leer |
| Größe nach `AGENTS.md` §3.14 (Befehl dort, Pfade `tools harness/README.md harness/vertraege`) | 581 Zeilen, Marke 1500 nicht erreicht |

Logs des Koordinators (`gps-voll-release.log`, `gps-voll-container.log`, `gps-voll.txt`) gelesen: `release-gegenprobe: grün (732 Läufe)`, Zeit 982 s von höchstens 1200 s, Last 8,06, Exit 0; `test-container-gegenprobe` 533 s, Last 3,09, Exit 0.

## DoD gegen Code

| DoD-Punkt | Befund |
|---|---|
| 1 `release-gegenprobe` höchstens 1200 s, Ergebnis gleich, Stichprobe höchstens 37 s | erfüllt. 982 s (Log), Stichprobe von mir 29 s. Deklarationszeilen: `git diff -U0 b60a6b9d..HEAD` entfernt keine Zeile `mutr`, `mf`, `mh`, `mutv`, `mut1`, `mut13g`; geändert sind zwei Zeilen, die nur erweitert werden (die Liste `ok` um `FAELLE_GP14`, `FAELLE_GP14S`; `m-p14-budget-ungueltig-ohne-exit` um die neuen Fälle). Mengen-Vergleich `GP_ENTDECKE=1` vor und nach: 568 Mutationen der Gruppen `p1`, `sh`, `b`, `f` und 116 `m-p19-` ohne Unterschied (§7, F-9 und Lauf 3). |
| 2 Fälle und Mutationen von Punkt 16 (9) | erfüllt. Ich sah 30 Mutationen `m-p14-` rot auf genau ihre Fälle, je Zahlform (`08`, `09`, `00`, `007`, sieben und zwanzig Stellen, `+5`, Leerraum vorn und hinten, `999999`, `10`, `0`), `Zeit`-Zeile, Last, Reihenfolge (c) vor (d), Ordner je Mutation und Fall, `GP_FALL` (Protokoll und `stderr`), Signale. |
| 3 `test-container-gegenprobe` höchstens 1200 s | erfüllt, 533 s (Log). |
| 4 Beleg der vollen Läufe | erfüllt mit Abweichung, siehe V-1. |
| 5 Tabelle in `SPEC-060` Punkt 16 | **offen**, siehe V-2. |
| 6 Ordner je Mutation und Fall entfernt, Spitze unter 100 MB, voller Lauf in einem Zug | Entfernen erfüllt (Fälle `p14-ordner-*` und `p12-ordner-mutation` mit Mutationen rot gesehen). Spitze: von mir an Stichproben gemessen (6,5 MB), §7 trägt keine Messung am vollen Lauf, siehe V-3. Der volle Lauf in einem Zug liegt vor (732 Läufe, `ee68abcd`). |
| 7 `make gates` | erfüllt am Stand `838ec960` (284 s, §7); `6cb151da` ändert danach nur den Plan (`git diff --stat`). Nicht neu gefahren, wie beauftragt. |
| 8 Vorschlag zum Wortlaut | erfüllt (§7). `AGENTS.md` und `.claude/` im Diff `b60a6b9d..HEAD`: leer, der Übergangssatz in `AGENTS.md` §6 Punkt 6 und `implement-slice.md` Schritt 15 ist **nicht** geändert. Die Entscheidung des Nutzers steht aus. |
| 9 bis 13 (Review, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) | gehören zur Closure, noch leer. Review r1 und r2 liegen vor (r2: 0 HIGH, 0 MEDIUM). |

## Plan-vs-Code-Diff

- **Randformen (g) bis (j):** im Code gehalten. `GP_BUDGET=0` trägt `B` = 0; die Zahlform ist `^(0|[1-9][0-9]{0,5})$`; `GP_FALL` und `GP_ENTDECKE` nach (h); die Reihenfolge der Gründe gilt in beiden Skripten; `>=` ist als akzeptiertes Negativ geführt, nicht als äquivalent (`AGENTS.md` §3.10, kein Verstoß).
- **Obergrenze 1200 s:** in `release-`, `container-` und `container-einrichten-gegenprobe` und im `.mk`-Kommentar nachgezogen. Der Rest-Treffer `600` in `release-gegenprobe.sh:435` ist ein Zeitversatz der Attrappe, kein Obergrenzen-Wert. In `SPEC-060`, `SPEC-069` und `SPEC-070` stehen 600 s nur als Herkunft („vorher“).
- **Umbau der Fälle `b` und `f`** (Stand je Mutation unter `$run/snap`, Pfad `$run/t.w` bleibt): von Punkt 16 (9) gedeckt (Randform-neutral, Fehler eines gemeinsamen Schrittes färbt jeden abhängigen Fall; belegt durch die 65 roten Fälle oben).
- **§3.9 (Plan folgt Korrektur):** §1, §3, §6 und der Kopf nennen 1200 s, die Randformen und die Ausnahme zu Rückgabe 5. Kleinigkeit: In §7 steht noch „Review r2 (Commit folgt)“ (Zeile 237); der Commit `9d9035d6` existiert.
- **Größe:** 581 Zeilen, ein Liefer-Bereich (Harness-Werkzeug), keine Rückführung ausgelöst; die Rückführung nach Faktor 2 (Stichprobe) ist nicht eingetreten (29 s gegen 37 s).

## Vorbehalte (mit Adresse)

- **V-1 Wer den vollen Lauf fuhr, und Umgebung.** Die Entscheidung des Nutzers (Rückgabe 5) nennt den Implementer; den Lauf am Endstand `ee68abcd` fuhr der Koordinator, der Implementer die Läufe an `7b8387a1` und `76456d59` (beide mit Befund, behoben, §7). Die Last beim Lauf (8,06) liegt über einem Fünftel der Kerne (Punkt 16 (1): „ruhig“ heißt unter 4 bei 20 Kernen); die 982 s gelten damit formal nicht als Messung auf ruhigem Host. Ein ruhiger Host ist eher schneller, die Marge zu 1200 s beträgt 218 s. Nach `ee68abcd` änderten sich 20 Zeilen in drei Gegenproben (Fälle `+5`, ` 5`, `5 `); sie sind nur an engsten Fällen gefahren (ich fuhr die unveränderte Kopie von `release-gegenprobe` und die `p14`-Mutationen grün). Adresse: Planner, in §7 der Closure als akzeptiert mit Last nennen; Risiko 3 aus §6 (`BEO-REPO/lastabhaengige-zaehlung`) bleibt *weiter offen*.
- **V-2 Tabelle in `SPEC-060` Punkt 16 und `CO-001`/`CO-003`.** Die Zeile `release-gegenprobe`, `test-container-gegenprobe` steht noch auf `nicht gemessen (Fälle von (9) fehlen noch)`; die Fälle liegen vor, der volle Lauf ist gemessen. DoD 5 ist offen. Die Spalte ist Vertragstext (Architect), die Stand-Zeile zu `CO-001` und `CO-003` gehört in die Closure (Planner). Adresse: Architect für die Zeile (`aufgenommen`) vor dem `git mv`, Planner für die Closure-Zeile. Das ist keine Lücke der Lieferung, sondern ein noch nicht gesetzter Haken; ohne ihn bleibt der Slice nicht schließbar.
- **V-3 Platz-Spitze im Beleg.** §7 trägt keine `du`-Messung am vollen Lauf (DoD 6 verlangt „im Bericht“). Meine Messung an drei Teilmengen (6,5 MB bei 107 Mutationen mit `GP_PARALLEL=4`, 5,2 MB bei `m-p19-push-`) und der leere Ordner danach belegen die Größenordnung weit unter 100 MB, nicht die Spitze des vollen Laufs mit `GP_PARALLEL=8`. Adresse: Planner, in der Closure als „an Stichproben gemessen“ nennen oder den Lauf des Koordinators mit `du` wiederholen (Aufwand 16 Minuten, nicht erforderlich).
- **V-4 Fahrt unter fremdem Benutzer, fremder Locale und Last** (`SPEC-060` Punkt 16 (8), Randform (f)): vom Implementer nicht geliefert. Ich fuhr die Stichprobe unter `LC_ALL=C TZ=Asia/Tokyo` (grün, 13 s) als Benutzer ohne root (uid 1000); als root und unter Last habe ich nicht gefahren. Der Vertrag verlangt die Zeit „genannt, ohne Geltung der Grenze“; ein roter Befund ist nicht zu erwarten, eine Zusage ruht aber auf einer nicht gefahrenen Fahrt. Adresse: Planner, entweder als Closure-Beleg nachholen oder mit der Beobachtung `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` als *weiter offen* führen (Nutzer entscheidet).
- **V-5 Regel-Entscheidung des Nutzers.** Der Übergangssatz in `AGENTS.md` §6 Punkt 6 und `implement-slice.md` Schritt 15 steht unverändert; der Vorschlag liegt in §7. Adresse: Nutzer (Entscheidung steht aus).

## Nicht geprüft

Voller Lauf beider Gegenproben (nicht gefahren, nur Logs), `make gates` (nach `838ec960` nur Plan-Änderungen), Fahrt als root und unter Last, voller Mengen-Vergleich der Gruppe `h` am Endstand (Lauf 3 vergleicht `m-p19-` vor und nach dem Umbau; der volle Lauf am Endstand belegt die deklarierten Mengen).
