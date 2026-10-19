# Review-Report: slice-erster-release-freigabe — 2026-10-19

**Review-Art:** Code, gegen Plan, Entscheidungen und Hard Rules.
**Gegenstand:** Diff e2372b18..impl-fr (Kopf 9f5425cb, Basis der Metrik 760c51d4). Geprüft: release.sh, release.yml, release-gegenprobe.sh, docs/maintainer/releasing.md, harness/README.md, Slice-Plan §7.
**Skill:** .harness/skills/reviewer.md @ 229b5efb · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-19
**Eingangs-Kontext:** Slice-Plan §1, §3, §4, §6, §7; ADR-0052 und ADR-0051; SPEC-067 Punkte 1, 11, 13, 16, 17, 18; LH-FA-16; AGENTS.md §3.7 und §3.9 bis §3.14.

**Eigene Läufe (nur engste Sensoren):** `bash -n` auf beide Skripte grün; `make docs-check` grün (1245 Dateien, 0 Befunde); zehn Mutationen m-p16-* per `GP_MUT` mit `GP_ENTDECKE=1` stimmen mit der Tabelle in §7 überein; Größenmetrik 527 stimmt; der Pin `aquasec/trivy@sha256:af6acf9a…855daa` löst als Tag 0.75.0 auf (vom Host gemessen); sechs eigene Sonden-Mutationen in einer Wegwerfkopie, drei blieben grün (F-2).

## Findings

| ID | Kat. | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Code trägt vorläufige Lesarten für neun Randformen, die SPEC-067 Punkt 16 nicht entscheidet. Der Implementer hat weitergearbeitet statt anzuhalten. Die Gegenprobe bindet die Lesarten mit Fällen und Mutationen (p16-user, p16-cache-dir, p16-cache-kein-zwischenspeicher, p16-ohne-image-mit-layout, p16-alte-marke, p16-eingaben-loescht-marke, p16-image-loescht-marke). Entscheidet der Architect anders, sind Code, Fälle und Mutationen neu zu schreiben. | AGENTS.md §3.12 | Plan §7 · "vorläufige Lesart"; release.sh · "db_alter() {" | nein, Entscheidung des Architect | Randform im Code vor der Entscheidung des Architect |
| F-2 | MEDIUM | §7 sagt "Grüne Mutanten: keiner". Drei Wächter haben aber keinen Fall und bleiben bei Entfernen grün: `\|\| scan_fehler pruefmenge …` hinter `# P16auswertung`, dasselbe hinter `# P16zaehlung`, und `(type == "array") and`. Ist `Results` kein Array oder `Vulnerabilities` nicht iterierbar, wird leer ausgewertet und der Scan endet grün. Punkt 16 sagt "Gescheitert … nie grün". | AGENTS.md §3.10; SPEC-067 Punkt 16 | release.sh · `# P16auswertung`, `# P16zaehlung`, `(type == "array") and (length > 0)` | ja, Mutation per GP_MUT färbt keinen Fall | Fail-closed-Wächter ohne Fall; grüner Mutant nicht berichtet |
| F-3 | LOW | Die Werkzeug-Vorprüfung der Gegenprobe nennt weder `jq` noch GNU `date`, obwohl Kopf und Gate-Index sie verlangen. Fehlt `jq`, werden viele p16-Fälle rot ohne Hinweis auf die Ursache, statt Exit 2 mit "jq fehlt". | Maintainability | release-gegenprobe.sh · `for t in git sed grep awk cmp tar mktemp … sha256sum; do` | ja, PATH ohne jq | Werkzeugliste der Gegenprobe nicht nachgezogen |
| F-4 | LOW | Die Zeile `scan` der Fehlertabelle rät bei "nicht gelaufen (Exit 2: Netz, Datenbank, Plattform)" zu "wiederholen". Die Stichwörter `plattform` und `pruefmenge` zeigen eine falsche Annahme über Trivy oder das Skript, dafür verlangt Punkt 18 (b) einen neuen, höheren Tag. | SPEC-067 Punkte 11 und 18 (b) | releasing.md · "nicht gelaufen: wiederholen; Befund: den Befund beheben" | nein | Wiederanlauf-Rat widerspricht Punkt 18 (b) |
| F-5 | INFO | Ein `UpdatedAt` in der Zukunft ergibt ein negatives Alter, und `[ "$alter" -le "$DB_GRENZE" ]` nimmt es an. ADR-0052 nennt nur "älter als 72 Stunden oder nicht lesbar". Rückgabe an den Architect. | ADR-0052 | release.sh · `[ "$alter" -le "$DB_GRENZE" ]` | ja, Fixture mit Zukunftsdatum | Randform Zeitpunkt in der Zukunft nicht entschieden |
| F-6 | INFO | releasing.md §10 Punkt 2 nennt "Ablage von `pgr-gp`". SPEC-061 nennt den Namen "vorläufig", SPEC-067 Punkt 17 nennt keinen Namen. | SPEC-061 | releasing.md · "der Ablage von `pgr-gp` per `ssh`" | nein | Vorläufiger Name in der Checkliste |
| F-7 | INFO | (a) SPEC-067 Punkt 18 endet mit "(Frage an den Nutzer)" zu v0.1.1, obwohl ADR-0052 die Entscheidung festhält und releasing.md sie als entschieden führt (Architect). (b) §1 und DoD führen für den Scan einen Handbuch- und README-Teil; der Diff liefert ihn nicht ("Entwurf … nicht geliefert"), und §1 und DoD sind nicht nachgezogen (Verifier, AGENTS.md §3.9). | SPEC-067 Punkt 18; AGENTS.md §3.9 | SPEC-067 · "(Frage an den Nutzer)"; Plan · "Entwurf … (nicht geliefert" | ja, `grep` | Vertrag und Plan folgen der Entscheidung nicht |

## Die neun Rückgaben an den Architect (nicht vom Reviewer entschieden)

| # | Randform | Vertrag muss vorab entscheiden? |
|---|---|---|
| 1 | Quelle des Alters der Datenbank, Cache-Mount, `--cache-dir`, `--user` | ja |
| 2 | Marke `image.ok` gegen `image/digest` und `layout.tar` | ja |
| 3 | Form der Meldung bei Exit 2 | ja |
| 4 | Reihenfolge mehrerer Gründe | ja (§3.12) |
| 5 | `scan.ok` und Zustand früherer Läufe (Punkte 11 und 13) | ja |
| 6 | Summenzeile bei Grün | ja |
| 7 | Zeile `scan` in der Tabelle von Punkt 11 | ja |
| 8 | Image ohne Paket, `Results` leer | Vertrag bestätigt oder ändert; gemessen wird im Trockenlauf |
| 9 | Nicht gemessene Trivy-Annahmen | Verfahren nach falscher Annahme; gemessen wird im Trockenlauf |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| release.yml (Scan ohne `if` und `env`, zwischen Image und Veröffentlichen, gleich Punkt 1b) | geprüft, ohne Befund |
| release.sh Scan (Aufruf E1 bis E10, Schwelle, Behebung fehlend/null/leer, Marke, Exit 1/2, Pin, kein Token, keine Ausnahmedatei) | geprüft, ohne Befund außer F-1, F-2, F-5 |
| release-gegenprobe.sh (Fallnamen, Tabelle §7 gegen zehn gemessene Mutationen, nachgezogene Erwartungen von A und B) | geprüft, ohne Befund außer F-2, F-3 |
| releasing.md §5 und §10 gegen Punkte 17 und 18 | geprüft, ohne Befund außer F-4, F-6 |
| harness/README.md (Zeile release-gegenprobe) | geprüft, ohne Befund |
| Größenmetrik 527 | geprüft, ohne Befund |
| Core-Reinheit, Suppression, ADR-Immutabilität, Move gegen Inhalt | nicht berührt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 3 |

**Summary-Zeile:** `slice-erster-release-freigabe: 0 HIGH · 2 MEDIUM · 2 LOW · 3 INFO`

**Finding-Klassen:** Randform im Code vor der Entscheidung des Architect · Fail-closed-Wächter ohne Fall; grüner Mutant nicht berichtet · Werkzeugliste der Gegenprobe nicht nachgezogen · Wiederanlauf-Rat widerspricht Punkt 18 (b) · Randform Zeitpunkt in der Zukunft nicht entschieden · Vorläufiger Name in der Checkliste · Vertrag und Plan folgen der Entscheidung nicht

## Verdikt

**Merge-blockierend:** ja, wegen der zwei MEDIUM. F-1 ist durch die Entscheidungen des Architect (Commit dddf7cac auf arch-fr) aufzulösen; F-2 verlangt je einen Fall und eine Mutation für die drei Wächter.

**Übergabe:** Findings an den Implementer; F-5 und F-7 (a) an den Architect; F-7 (b) an den Verifier; Finding-Klassen in die Closure §7.
