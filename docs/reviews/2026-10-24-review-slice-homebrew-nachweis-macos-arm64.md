# Review-Report: slice-homebrew-nachweis-macos-arm64 — 2026-10-24

**Review-Art:** Code — Code-Review (Doku-Diff) gegen Plan und Konventionen, ohne DoD-Prüfung.

**Gegenstand:** Commit `66468b30` (Zweig `impl-mac1`, Basis `00634a58`): `README.md`, `docs/user/benutzerhandbuch.md`, `docs/maintainer/releasing.md`.

**Skill:** `.harness/skills/reviewer.md` @ `00634a58` ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-24

> **Zitier-Form.** Kennung, nicht Adresse: Slices als `slice-<Kennung>`, Baseline-Stellen als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice).

**Eingangs-Kontext:**

- `slice-homebrew-nachweis-macos-arm64` (Plan: §1 Ziel und Abgrenzung, §2 DoD, §3 Plan, §6 Risiken)
- `LH-FA-19.a` (Absatz *Bericht je Kombination*) und `SPEC-067` Punkt 20 (6)
- keine ADR im Commit genannt; [ADR-0053](../plan/adr/0053-homebrew-tap-in-der-release-pipeline.md) nur mittelbar (Homebrew)
- `LH-FA-19`
- `AGENTS.md` (Hard Rules, besonders §3.9, §3.11)

---

## Findings

Gelesen wurde jede geänderte Zeile einzeln gegen: „berichtet“ mit Kombination und Version, nie „gefahren/geprüft/belegt“; nichts zu `brew test`, `replay`, Exit-Codes, Ausgabe der Installation; Linux/amd64 allein „gefahren“; macOS/amd64 und Linux/arm64 Grenze; Ist-Zustand ohne Chronik und ohne Verweise auf Spec, Slice, Review.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Satz sagt „das Handbuch nennt nur die geprüfte Kombination“ und im selben Klammerausdruck, für macOS auf `arm64` liege ein Bericht vor. Das Handbuch nennt nach dem Diff auch den Bericht; der erste Halbsatz ist damit nicht mehr wahr. | `AGENTS.md` §3.11 | `docs/maintainer/releasing.md` · „das Handbuch nennt nur die geprüfte Kombination; für macOS auf `arm64` liegt ein Bericht“ (§7) | nein — kein Gate liest den Satz | Satz in Wartungs-Doku widerspricht der eigenen Ergänzung |
| F-2 | LOW | Der Diff ändert §10 (Satz „`darwin` und `windows` sind gebaut, nicht ausgeführt (§2); für `darwin` auf `arm64` nennt das Handbuch zusätzlich den Bericht (§7)“). Plan §3 und DoD nennen für `releasing.md` nur §1, §7 und die Änderungshistorie. | `AGENTS.md` §3.9 | `docs/maintainer/releasing.md` · „für `darwin` auf `arm64` nennt das Handbuch zusätzlich den Bericht (§7), sonst nur das.“ (§10) · Plan `slice-homebrew-nachweis-macos-arm64` §3 | nein | Plan folgt der Korrektur nicht (§10 fehlt in §3 und DoD) |
| F-3 | LOW | Das Handbuch schreibt im Abschnitt „Binary vom GitHub-Release“ neu „nicht von uns ausgeführt“. Im Handbuch steht sonst „nicht ausgeführt“ (Voraussetzungen, Bau), das README schreibt „von den Prüfungen nicht ausgeführt“; der Handbuch-Text hat kein „wir“ als Subjekt, und die Stellen sagen Verschiedenes. | Maintainability | `docs/user/benutzerhandbuch.md` · „gebaut, aber nicht von uns ausgeführt“ | nein | uneinheitliche Formulierung derselben Aussage in Handbuch und README |
| F-4 | INFO | Die Kopfzeile des Handbuchs bleibt `Version: 0.13`, `Stand: 24.10.2026`, obwohl der Diff drei Aussagen ergänzt. Eine Zeile in der Änderungshistorie braucht es nicht: ihre Zeilen tragen Software-Versionen und sind für die dritte Zelle die Quelle der Release-Notes (`SPEC-067` Punkt 3 (5), Punkt 9); eine Zeile ohne Release stünde dort fremd. Ein Bump der Kopf-Version ist uneinheitlich geübt (`47b1e489` bumpt bei neuem Abschnitt, `3bc2f931` bei einer Verengung nicht); `Stand` stimmt. | Maintainability | `docs/user/benutzerhandbuch.md` · „Version: 0.13“ | nein | Kopf-Version ohne feste Regel bei Doku-Änderung |
| F-5 | INFO | Die Abschnitte zum Bau (`docs/user/benutzerhandbuch.md` · „Die Binaries für macOS und Windows sind gebaut, aber nicht ausgeführt“) bleiben unberührt. Sie meinen das lokal mit `make binaries` gebaute Binary (`dev`), nicht das Release-Binary des Berichts; ein Widerspruch zum Bericht besteht nicht, ein Leser könnte beide aber nebeneinander lesen. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „Die Ausgabe ist `pgwire-recorder dev`. Die Binaries für macOS und Windows sind gebaut, aber“ | nein | Aussage zum Bau neben Bericht zum Release-Binary ohne Abgrenzung |
| F-6 | INFO | Handbuch und README schreiben „berichtet ein Nutzer“ ohne Datum. „Nutzer“ meint nach `LH-FA-19.a` den Nutzer des Repos (Auftraggeber), im Handbuch liest es sich als beliebiger Anwender. `LH-FA-19.a` verlangt für Handbuch und README Kombination und Version, kein Datum; das Datum steht in `releasing.md`. | `LH-FA-19.a` | `README.md` · „berichtet ein Nutzer“ | nein | Quelle des Berichts im Handbuch mehrdeutig benannt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `README.md` Zeile zu Homebrew (Linux/amd64 „gefahren“, macOS/arm64 „berichtet“ mit Version 0.2.0 und Homebrew 7.0.10, macOS/amd64 und Linux/arm64 „nennt die Formel die Dateien, ausgeführt ist die Installation dort nicht“) | geprüft, ohne Befund |
| `README.md` Satz zu den Binaries („von den Prüfungen nicht ausgeführt (für macOS auf `arm64` steht oben ein Bericht)“) | geprüft, ohne Befund |
| `docs/user/benutzerhandbuch.md` Voraussetzungen (Homebrew-Satz und Binary-Satz; „geprüft, dass sie gebaut sind“ meint den Bau und bleibt Bau-bezogen) | geprüft, ohne Befund |
| `docs/user/benutzerhandbuch.md` Abschnitt „Mit Homebrew installieren“ (kein `brew test`, kein `replay`, kein Exit-Code, „nicht gefahren“ nur als Verneinung; Anker `#mit-homebrew-installieren` löst auf) | geprüft, ohne Befund |
| Handbuch und README: kein Verweis auf Spec, Lastenheft, ADR, Slice, Welle, Review; keine Chronik („noch nicht“, „kommt mit“) | geprüft, ohne Befund |
| Handbuch Plattformgrenze und Binary-Satz (Zeilen 69 bis 75, 148) sowie Image-Sätze („Gestartet ist nur das Image für `linux/amd64`“) | geprüft, ohne Befund (Image betrifft den Bericht nicht; Linux/amd64 bleibt allein „gefahren“) |
| `docs/maintainer/releasing.md` §1 und Änderungshistorie (Zeile 0.18 mit Kopf-Version 0.18, Datum, Inhalt konsistent mit §1 und §7; Version 0.2.0, Homebrew 7.0.10, Datum 2026-10-24 gleich dem Plan §1) | geprüft, ohne Befund |
| `docs/maintainer/releasing.md` §7 „Die drei übrigen“ (Zählung bleibt drei; Grenze macOS/amd64 und Linux/arm64 korrekt) | geprüft, ohne Befund (Wortlaut siehe F-1) |
| Wortlaut gegen `LH-FA-19.a` Absatz *Bericht je Kombination* und `SPEC-067` Punkt 20 (6) (Kombination, Version, „berichtet“, Schritte (1) und (2) getragen, Schritt (3), `brew test`, Ausgabe, Exit-Codes nicht getragen, je Kombination einzeln) | geprüft, ohne Befund |
| Commit-Regel (kein Move im Commit, Kennungen `slice-…` und `LH-FA-19` in der Message) | geprüft, ohne Befund |
| Dateien außerhalb des Diffs (`spec/`, `harness/vertraege/`, Freigabe `0.2.0`, Abdeckungstabellen, Skripte) | geprüft, ohne Befund (unberührt, wie Plan §1) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Satz in Wartungs-Doku widerspricht der eigenen Ergänzung · Plan folgt der Korrektur nicht (§10 fehlt in §3 und DoD) · uneinheitliche Formulierung derselben Aussage in Handbuch und README · Kopf-Version ohne feste Regel bei Doku-Änderung · Aussage zum Bau neben Bericht zum Release-Binary ohne Abgrenzung · Quelle des Berichts im Handbuch mehrdeutig benannt

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM; die drei LOW sind Wortlaut und Planpflege.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation — DoD-Konformität prüft der Verifier separat.
