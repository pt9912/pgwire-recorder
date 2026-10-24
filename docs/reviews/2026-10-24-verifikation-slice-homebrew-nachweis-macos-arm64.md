# Verifikationsbericht: slice-homebrew-nachweis-macos-arm64 — 2026-10-24

**Rolle:** Verifier (DoD- und Plan-Konformität, kein Review, keine Reparatur) · **Modell:** claude-sonnet-5-5
**Gegenstand:** Zweig `impl-mac1`, Stand `f675d8e7`, Basis `00634a58`, Basis-Hash des Slice (§4) `34a8cf11`.
**Maßstab:** DoD §2 des Slice-Plans, `LH-FA-19.a` (Absatz *Bericht je Kombination*), `SPEC-067` Punkt 20 (6), `AGENTS.md` §3.9, §3.11, §3.14.

## DoD-Punkte

| DoD-Punkt | Ergebnis | Beleg |
|---|---|---|
| 1 Doku: Handbuch, README, `releasing.md` §1, §7, §10, Historie nennen macOS/arm64 `0.2.0` als berichtet; macOS/amd64 und Linux/arm64 Grenze; Linux/amd64 allein gefahren | erfüllt | Diff `66468b30` plus `5dae30d9`; Handbuch (Voraussetzungen, Abschnitt zur Installation mit Homebrew, Abschnitt zum Binary vom GitHub-Release), README (Homebrew-Zeile, Binaries-Satz), `releasing.md` §1, §7, §10, Zeile 0.18 und Kopf-Version 0.18 |
| 2 Zitat und Merkmale in §1 | erfüllt | §1 zitiert wortgetreu: „ich habe es eben auf einem macos installiert und die version abgefragt“, „Darwin arm64 -- Homebrew 7.0.10“, Befehl `brew install pt9912/pgwire-recorder/pgwire-recorder`, „pgwire-recorder 0.2.0“, Datum 2026-10-24; keine Beleg-Datei, keine Summe (Diff berührt keine) |
| 3 `docs-check`, `kopf-check` grün | erfüllt, selbst gefahren | `make docs-check`: „1467 Datei(en) geprüft, 0 Befund(e)“, Exit 0; `make kopf-check`: Exit 0. `make gates` nicht gefahren (Auftrag) |
| 4 Review-Report | erfüllt | `docs/reviews/2026-10-24-review-slice-homebrew-nachweis-macos-arm64.md` (0 HIGH, 0 MEDIUM, 3 LOW, 3 INFO) |
| 5 bis 8 Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen, erwartet | §7 „(entsteht bei Closure)“; Closure-Pflichten, nicht Gegenstand dieser Verifikation; §6 trägt je Risiko die Ausgangsform im Plan |

## Prüfungen laut Auftrag

- **Kein „gefahren/geprüft/belegt“ für macOS/arm64:** `grep` über README, Handbuch, `releasing.md` findet „gefahren“ nur für Linux/amd64 und als Verneinung („Bericht, nicht gefahren“); „belegt“ nur für Linux/amd64, Bau/Dateikopf/Summe und in „kein Beleg für `replay`, `brew test` und Exit-Codes“. Kein Satz belegt macOS/arm64 als geprüft. Ohne Befund.
- **Handbuch (§3.11):** keine Verweise auf Spec, Lastenheft, ADR, Slice, Welle, Review; keine Chronik; der interne Anker zum Homebrew-Abschnitt löst auf (docs-check grün). Ohne Befund.
- **Linux/amd64 allein gefahren:** README, Handbuch, `releasing.md` §1 und §7 ja. macOS/amd64 und Linux/arm64 bleiben „ausgeführt ist die Installation dort nicht“ bzw. Grenze. Ohne Befund.
- **Nichts zu `brew test`, `replay`, Exit-Codes, Ausgabe der Installation auf macOS:** nur in der Verneinung (`releasing.md` §7). Ohne Befund.
- **Plan-vs-Diff (§3.9):** geänderte Dateien: `README.md`, `docs/user/benutzerhandbuch.md`, `docs/maintainer/releasing.md` (alle in §3 und DoD), Plan (nachgezogen mit §10), Review-Report. `releasing.md` §10 steht seit `5dae30d9` in §3 und DoD (F-2 behoben). Spec, Vertrag, Freigabe `0.2.0`, Abdeckung, Skripte unberührt.
- **Review-Findings:** F-1 (`releasing.md` §7: „das Handbuch nennt die geprüfte Kombination und den Bericht“) behoben; F-2 (Plan §3 und DoD nennen §10) behoben; F-3 („nicht von uns ausgeführt“ entfernt, Handbuch schreibt „gebaut, aber nicht ausgeführt“) behoben. INFO F-4 bis F-6 sind keine Pflicht.
- **Marke 1500 (§3.14, Messbefehl aus `AGENTS.md`):** Basis `34a8cf11`, Pfade `docs/user docs/maintainer README.md`, beide Ausschlüsse: Summe **41** (mit Basis `00634a58` ebenfalls 41). Weit unter 1500; Marke in §4 genannt.

## Anmerkungen ohne Befund-Rang (INFO)

- DoD 1 nennt „Handbuch §2 *Installation*“; die Änderung der Voraussetzungen liegt in §1 *Einleitung* (Handbuch Zeilen 66 bis 80). §3 führt die Datei, die Änderung ist inhaltlich gedeckt; die Klammer der DoD ist ungenau.
- `releasing.md` §2 „Abnahme-Umgebungen“ sagt weiter, `darwin` sei „nicht ausgeführt“; das bezieht sich auf die Prüfungen dieses Repos und widerspricht dem Bericht nicht.
- §7 des Plans trägt noch keine Belege des Implementers zum vollen `make gates`-Lauf; er entsteht mit der Closure.

## Verdikt

**Konform.** Alle Liefer-Punkte der DoD sind belegt, `docs-check` und `kopf-check` grün, Plan und Diff stimmen überein, Review-Findings F-1 bis F-3 behoben, Marke nicht berührt. Offen sind nur die Closure-Pflichten (Notiz, Lerneintrag, Register, Paarungen), die dem Planner folgen. Nicht gefahren: Gegenproben, voller `make gates`-Lauf (Auftrag).
