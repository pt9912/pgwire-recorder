# Verifikationsbericht: slice-play-feldvergleich-doku — 2026-10-23

**Rolle:** Verifier (Frage: Bauen wir es richtig? gegen Plan und DoD). **Gegenstand:** Zweig `impl-diff-all-doku`, Endstand `a502e40d`, Basis `3d3832d9`.

**Ergebnis: bestanden** (kein Befund, der die Closure sperrt; zwei Hinweise unten).

## DoD gegen Beleg

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| Liefer-Punkt 1 Handbuch | bestätigt | Optionstabelle, Zeilen `PGR-E2001`/`PGR-E2004`/`PGR-E5004`, Abschnitt *Alle Abweichungen einer Interaktion melden*, Satz zur ersten Abweichung eng gefasst; keine Verweise auf Spec/ADR/Slice/Welle, keine Chronik (Diff gelesen); Beispiele stehen als Zeilen im Test `play_vergleich_handbuch_e2e_test.go` und laufen gegen das gebaute Binary |
| Liefer-Punkt 2 README | bestätigt | eine Stelle (Meldung `PGR-E5004`) ergänzt, übrige begründet in §7 unverändert |
| `make gates` grün | Beleg liegt vor (Plan §7) | selbst gefahren: `make docs-check` (1419 Dateien, 0 Befunde), `make kopf-check`, `make abdeckung-check` ohne Fehler; volles `make gates` nicht wiederholt |
| Review-Report | liegt vor, nach dem eigenen Urteil gelesen | 0 HIGH, 0 MEDIUM, 4 LOW, 1 INFO; F-1 bis F-5 laut §7 nachgearbeitet (Commit `a502e40d`) |
| Closure, Register, Risiko-Ausgänge, Paarungen | offen | gehören dem Planner; §6 trägt noch keine Ausgänge (erwartet) |

## Eigene Läufe

- Engster Fall: `make test-integration` mit `PGR_TESTS_1` auf die vier Tests (`…Einzelheiten`, `…Obergrenzen`, `…Quellen`, `…Fortsetzung`): alle 4 Tests und alle Unterfälle PASS. Die Phase ohne PostgreSQL ist bei diesem Teillauf rot (fehlende Sicht-Dateien der vollen Phase 1), wie §7 vermerkt, und zählt nicht.
- **Bewusstes Brechen (Modul 11):**
  - `maxPositionen` 20 -> 21 in `internal/hexagon/services/play_vergleich.go`: `21_Abweichungen` und `23_Abweichungen` rot, übrige grün. Belegt „höchstens 20 Positionen, `und N weitere`".
  - `maxEinzelne` 20 -> 21: `21_Spalten` und `23_Spalten` rot, `20_Spalten` grün. Belegt „höchstens 20 Spalten, Grenze zählt Spalten".
  - Beide Mutanten zurückgesetzt (`git status` sauber). Die Mutation der Kombinationsprüfung (`PGR-E2001`) und die Obergrenzen-Mutanten stehen laut §7 gesehen; die Kombinationsprüfung habe ich nicht erneut gefahren.
- Marke 1500 (§3.14, Befehl aus `AGENTS.md`, Pfade `docs/user README.md test`): 341 Zeilen; Plan nennt 340 (Zählung vor der Nacharbeit bzw. Rundung, unerheblich, weit unter der Marke).
- Plan-vs-Code: Diff berührt nur `README.md`, `docs/user/` (Handbuch, `abdeckung-e2e.md`), den neuen Test, Plan und Review-Report; kein Produkt-Code, keine Spec-Änderung; Abgrenzung §1 gehalten; Liefer-Punkte 2, Schichten 1.

## §3.11 und Entscheidungs-Konformität

Jeder Satz zu `--diff-all` im Handbuch hat einen Fall: Quellen (Option, Umgebung, Schlüssel), `PGR-E2001` ohne Vergleich (auch `=false`, auch aus Umgebung/Datei), `record`/`replay` (Option, Variable, Schlüssel), Einzelheiten je Art, 20/21/23 für Positionen und Spalten, Fortsetzung mit und ohne `--continue-on-error`, Exit 5. Die Nacharbeit zu F-1 bis F-4 ist im Text und im Test sichtbar. Keine ADR oder Hard Rule verletzt.

## Hinweise (nicht sperrend)

1. Handbuchsatz „je abweichender Spalte … Name und Typ-OID beider Seiten, soweit sie abweichen" hat keinen Fall für eine Abweichung nur des Namens und keinen für „Spaltenzahl verschieden und Spalte abweichend" zugleich (Review F-5 b, vom Implementer bewusst nicht ergänzt). Die Aussage ist ein „soweit"-Satz und wird durch Typ-OID und Name zusammen abgedeckt; ein Mutant „Name nur bei gleicher OID" bliebe grün.
2. Die Tests lesen das Handbuch nicht; die Kopplung Satz zu Fall hält allein §7 des Plans. Ändert später jemand einen Beispielwert im Handbuch, wird kein Test rot. Bekannte Eigenschaft von §3.11, hier nur genannt.

## Verbleibende Closure-Pflichten (Planner)

§6 Risiken 1 bis 3 (Ausgänge: entfallen mit Begründung, durch DoD 1 und die Läufe belegt); Closure-Notiz mit Lerneintrag und Register (Finding-Klassen aus dem Review); reiner `git mv` nach `done/`.
