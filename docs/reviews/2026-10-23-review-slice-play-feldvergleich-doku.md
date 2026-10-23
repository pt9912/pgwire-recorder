# Review-Report: slice-play-feldvergleich-doku — 2026-10-23

**Review-Art:** Code — Code-Review gegen Plan, Spezifikation und Hard Rules (`AGENTS.md`).

**Gegenstand:** Zweig `impl-diff-all-doku`, Diff `main..HEAD` (`5f35239f` Handbuch, README und Test, `81c4768d` Plan §7).

**Skill:** `.harness/skills/reviewer.md` @ `81c4768d`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-23

> **Zitier-Form.** Kennung, nicht Adresse: Slices als `slice-<Kennung>`, Sensoren als `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code.

**Eingangs-Kontext:**

- Slice-Plan `slice-play-feldvergleich-doku` (§1 bis §4, §6, §7 Belege des Implementers, §8)
- Spezifikation `LH-FA-24.a` (Absatz *Alle Abweichungen (`--diff-all`)*), `LH-FA-17.a`
- `LH-FA-24`; keine ADR im Diff genannt
- `AGENTS.md` §3.9, §3.11, §3.13
- `v6.18.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice

Gelesen, nicht gefahren: kein `make gates`, kein Integrationslauf (laut Implementer grün am Endstand); jede Handbuchaussage zu `--diff-all` gegen die vier Tests in `test/integration/play_vergleich_handbuch_e2e_test.go` und den Text der Spezifikation gelesen. Beide Beispielzeilen des Handbuchs stimmen im Meldungstext und im Zähler mit den erwarteten Zeichenketten der Tests überein.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Das Handbuch sagt „Ohne `--compare-responses` ist `--diff-all` ein ungültiger Aufruf“; nach der Spezifikation gilt das für wirksam `true`, `--diff-all=false` ist auch ohne Vergleich gültig. Der Test belegt nur die Form ohne Wert (`--diff-all`, Umgebungsvariable `true`, Schlüssel `true`). | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „Ohne `--compare-responses` ist `--diff-all` ein ungültiger Aufruf“ | ja — Fall `--diff-all=false` ohne Vergleich im Test | Handbuchsatz weiter als der Test |
| F-2 | LOW | Der Kommentar über `TestE2EPlayHandbuchDiffAllEinzelheiten` sagt, die Meldung habe „ohne die Option nicht“ die Form mit der Gesamtzahl, und das Handbuch sagt „ohne `--diff-all` hat sie die Form aus Antworten beim Einspielen vergleichen“. Jeder Lauf dieser Datei setzt `--diff-all` (`handbuchLauf`, `vergleichLauf` mit der Option); der Teil ohne Option wird hier nicht gefahren. | `AGENTS.md` §3.11 | `test/integration/play_vergleich_handbuch_e2e_test.go` · „ohne die Option nicht“ | ja — Lauf ohne Option in dieser Datei | Kommentar weiter als der Test |
| F-3 | LOW | Das Handbuch nennt für `record` und `replay` die Option und `PGWIRE_RECORDER_DIFF_ALL` als `PGR-E2001`. Der Test fährt die Umgebungsvariable nur bei `replay`; bei `record` fehlt der Fall. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „ebenso `PGWIRE_RECORDER_DIFF_ALL`“ | ja — Fall `record Umgebungsvariable` | Handbuchsatz weiter als der Test |
| F-4 | LOW | Das Handbuch schreibt im Abschnitt zum Ende der Verbindung „nennt die Meldung ohne `--diff-all` diese Abweichung“; das legt nahe, dass sie mit der Option etwas anderes nennt. Kein Test fährt FATAL oder PANIC mit `--diff-all` und einer davor abweichenden Antwort; der Implementer führt den Fall in §7 als nicht belegt. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „nennt die Meldung ohne `--diff-all` diese Abweichung“ | ja — Lauf mit FATAL und abweichender Antwort davor, mit Option | Handbuchsatz deutet Verhalten ohne Beleg an |
| F-5 | INFO | Ohne erwartete Aktion, Stellen am Rand der Belege: (a) Die Zeile `PGR-E2004` zählt die Schlüssel mit Wertemenge `true`/`false` auf; `diff_all` fehlt, der Wert ist im Handbuch nicht genannt (laut §7 bewusst). (b) Eine Abweichung nur des Namens einer Spalte und die Kombination aus verschiedener Spaltenzahl und abweichender Spalte hat kein Test. (c) Die Zustände `I` und `T` im Beispiel zu `Transaktionsstatus` erklärt das Handbuch nicht. (d) `docs/user/abdeckung-e2e.md` ändert sich im Diff (erzeugt von `make abdeckung`), steht aber nicht in §3 des Plans. | Maintainability | `docs/user/benutzerhandbuch.md` · „`compare_responses` und `allow_plaintext` (erlaubt sind `true` und `false`)“ | nein | Stellen am Rand der Belege |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Handbuch, Abschnitt *Alle Abweichungen einer Interaktion melden*: Obergrenzen (20 Positionen, 20 Spalten, „und N weitere“, Spalten statt Angaben, Gesamtzahl), Einzelheiten je Art, Fortsetzung mit und ohne `--continue-on-error`, Exit-Code 5, `PGR-E2001`/`PGR-E2004` je Quelle und Betriebsart | geprüft, ohne Befund (bis auf F-1, F-3, F-4) |
| Zählung „Spalten statt Felder“ im Handbuch (Satz und Beispiel `Name … Typ-OID` als ein Eintrag) | geprüft, ohne Befund; `…Obergrenzen` belegt 20 Spalten mit 40 Angaben, die 21. und 23. |
| Beispielzeilen des Handbuchs (Typ-OID, COMMIT gegen BEGIN) samt Zähler | geprüft, ohne Befund |
| Stellen „erste Abweichung“ im Handbuch (Meldung, Fehlercode `PGR-E5004`, FATAL) | geprüft, ohne Befund; eng gefasst („Ohne `--diff-all`“) |
| Optionstabelle und Zeile `PGR-E2001` | geprüft, ohne Befund |
| `README.md` | geprüft, ohne Befund; nur die erste Stelle beschreibt eine einzelne Abweichung, sie nennt die Option; die übrigen Stellen bleiben aus dem Grund in §7 |
| Verweise auf Spezifikation, Lastenheft, ADR, Slice, Welle, Review im Handbuch und README | geprüft, ohne Befund; keine im Diff |
| Chronik („noch nicht“, „kommt mit“) im Handbuch und README | geprüft, ohne Befund |
| Kopf des Handbuchs (Version 0.12, Stand) | geprüft, ohne Befund |
| Plan §1 (Abgrenzung), §3, §4 (Marke 1500, Basis `3d3832d9`), §7, §8 gegen den Diff; kein Produkt-Code im Diff | geprüft, ohne Befund (bis auf F-5 d) |
| `AGENTS.md` §3.13 (Adresse, Nachzählen): Liefer-Punkte 2, Schichten 1 | geprüft, ohne Befund |
| Commit-Regel §3.3 (kein Move im Diff) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Handbuchsatz weiter als der Test · Kommentar weiter als der Test · Handbuchsatz deutet Verhalten ohne Beleg an · Stellen am Rand der Belege

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM; die vier LOW sind Aussagen, die weiter reichen als ein Test (`AGENTS.md` §3.11), ohne dass das Binary sie bricht.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation gegen die DoD.
