# Review-Report: slice-erster-release-freigabe, Runde 2 — 2026-10-19

**Review-Art:** Code, gegen Plan, Entscheidungen und Hard Rules (Layout je Plattform für Trivy nach dem roten Trockenlauf 39263455592).
**Gegenstand:** Diff b307b020..impl-fr2 (7502b603). Maßstab: SPEC-067 Punkt 16 (arch-fr3, 077da231), [ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md), [ADR-0052](../plan/adr/0052-scan-des-images-vor-dem-push.md), LH-FA-16, AGENTS.md §3.7, §3.9 bis §3.12.
**Skill:** .harness/skills/reviewer.md · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-19

**Engste Sensoren:** `bash -n` auf beide Skripte; `GP_ENTDECKE=1` mit `GP_MUT` auf `ok` und die 29 Mutationen der Scan-Auswahl (grün, jede färbt genau die genannten Fälle); `make docs-check` Exit 0 (1247 Dateien); acht eigene Sonden-Mutanten in einer Wegwerfkopie (gelöscht).

## Findings

| ID | Kat. | Befund | Quelle · Pfad | Klasse |
|---|---|---|---|---|
| F-1 | MEDIUM | `m-p16-cache-ueber-laeufe` färbt keinen Fall mehr: der Fall prüft nur den Präfix `-v $w/tmp/`, den jetzt auch der Mount des Plattform-Verzeichnisses erfüllt. Die Erwartung der Mutation wurde auf leer gesetzt, die Zelle in §7 ist leer, §7 behauptet "keinen grünen Mutanten". E11 ist damit nicht belegt. | AGENTS.md §3.10; release-gegenprobe.sh · `m-p16-cache-ueber-laeufe`, `f_p16_cache_ueber_laeufe` | Mutation durch geänderte Fall-Basis grün, Erwartung auf leer gesetzt |
| F-2 | LOW | `scan.ok` wird bei Beginn gelöscht; belegt nur für Exit 1 (`p16-marke-loescht-entfaellt`). Die neuen Exit-2-Stellen (Entpacken, Index lesen, Auswahl `plattform`) liegen hinter dem `rm -f`, haben aber keinen Fall mit alter Marke; Mutanten, die den `rm` hinter Index-Prüfung oder Auswahl legen, bleiben grün. | AGENTS.md §3.10 (Ortsangabe je Schritt); release.sh · `# P16marke-loescht` | Ortsangabe der Marken-Löschung nicht je Fehlerstelle belegt |
| F-3 | LOW | "`blobs/` mit denselben Blobs" ist nicht belegt: die Attrappe zeichnet nur die Namen auf; `mkdir "$pdir/blobs"` statt `cp -al` bleibt grün. | release.sh · `# P16blobs` | Inhalt einer kopierten Struktur nicht geprüft, nur ihr Name |
| F-4 | LOW | `p16-index-attestation` hat keine Zeile in der Mutationstabelle §7. Eine Mutation, die Indizes mit mehr als zwei Manifesten ablehnt, färbt genau `p16-index-attestation` und `p16-index-fremdes-os`; der Fall ist wirksam, die Mutation fehlt in der Tabelle. | Plan §7 | Zusage mit Fall, aber ohne eingetragene Mutation |
| F-5 | INFO | Manifest mit Variante (`arm64/v8`) steht nicht in §6 und hat keinen Fall; die Auswahl vergleicht nur `os` und `architecture`. (e) und (g) des Vertrags sind ungemessen. | AGENTS.md §3.12 | Randform nicht genannt; tragende Annahme ungemessen |
| F-6 | INFO | Zwei beobachtbar äquivalente Mutanten stehen nicht in §7: Entpacken je Plattform statt einmal, und `manifests[-1]` statt `manifests[0]` im äußeren Index. | release.sh · `# P16entpacken` | Äquivalenter Mutant nicht in §7 geführt |
| F-7 | INFO | `cp -al` und das Schreiben der Plattform-`index.json` sind nicht einzeln geprüft; ein Fehler zeigt sich als Exit von Trivy (`scan-gescheitert`). Hardlinks über Dateisystemgrenzen sind kein Fehler (Quelle und Ziel im selben Temp-Verzeichnis); ohne Hardlinks endet der Lauf rot, nie grün. | release.sh · `# P16blobs` | Fehlerquelle eines Hilfsschritts nur indirekt gemeldet |
| F-8 | INFO | Die Zeile `scan` der Fehlertabelle in `releasing.md` rät bei `scan-gescheitert` zu "wiederholen"; das Stichwort hat jetzt deterministische Ursachen (Layout nicht entpackbar, Index nicht lesbar). | releasing.md · Fehlertabelle | Fehlertabelle trägt neue Ursache des Stichworts nicht |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Entpacken, Lesen von `index.json` und Index-Blob; Index-Blob `{}`, `[]` oder ohne `platform` führt zu `plattform` vor `docker` | geprüft, ohne Befund |
| Auswahl je Plattform, genau ein Eintrag, Reihenfolge der Gründe | geprüft, ohne Befund |
| Aufruf E1 bis E13, kein `--platform`, Mounts, `--input /layout`, `--user` | geprüft, ohne Befund |
| Sonderzeichen in Pfaden, Quoting, Fehlerpfade ohne `set -e` | geprüft, ohne Befund |
| `releasing.md` §6 und `harness/README.md` | geprüft, ohne Befund außer F-8 |
| Plan gegen Code (§3.9) | folgt dem Code, außer "es gibt keinen grünen Mutanten" (F-1) |
| AGENTS.md §3.7 und §3.13 | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 4 |

**Summary-Zeile:** `slice-erster-release-freigabe (Runde 2): 0 HIGH · 1 MEDIUM · 3 LOW · 4 INFO`

**Finding-Klassen:** Mutation durch geänderte Fall-Basis grün, Erwartung auf leer gesetzt · Ortsangabe der Marken-Löschung nicht je Fehlerstelle belegt · Inhalt einer kopierten Struktur nicht geprüft, nur ihr Name · Zusage mit Fall, aber ohne eingetragene Mutation · Randform nicht genannt; tragende Annahme ungemessen · Äquivalenter Mutant nicht in §7 geführt

## Verdikt

**Merge-blockierend:** ja, wegen F-1. Übergabe: F-1 bis F-4, F-6 bis F-8 an den Implementer; F-5 an den Architect.
