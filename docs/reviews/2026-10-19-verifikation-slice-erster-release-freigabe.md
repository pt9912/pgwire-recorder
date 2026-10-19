# Verifikationsbericht: slice-erster-release-freigabe (Modul 11)

**Rolle:** Verifier · **Datum:** 2026-10-19 · **Modell:** claude-sonnet-5-5
**Gegenstand:** Diff e2372b18..impl-fr (Kopf 5a828f74), Metrik-Basis 760c51d4.
**Maßstab:** Plan (DoD, §7), SPEC-067 Punkte 11, 13, 16, 17, 18, [ADR-0052](../plan/adr/0052-scan-des-images-vor-dem-push.md), Review-Report [2026-10-19-review-slice-erster-release-freigabe.md](2026-10-19-review-slice-erster-release-freigabe.md).

## Verdikt

Der Diff erfüllt die Zusagen, die gegen ihn prüfbar sind. Kein HIGH, kein MEDIUM; der Merge ist nicht blockiert.

**Summary-Zeile:** `slice-erster-release-freigabe (Verifikation): 0 HIGH · 0 MEDIUM · 2 LOW · 3 INFO`

## Eigene Läufe (engste Sensoren, nacheinander)

| Lauf | Ergebnis |
|---|---|
| `bash -n` auf `release.sh` und `release-gegenprobe.sh` | grün |
| `make docs-check` | Exit 0, 1246 Dateien, 0 Befunde |
| `make kopf-check` | Exit 0 |
| `make workflows-check` | Exit 0, kein Befund |
| Unveränderte Kopie, Fall `ok` | grün |
| Alle 107 Mutationen `m-p16-*` (`GP_ENTDECKE=1`) | grün (107 Läufe), je Mutation exakt die erwartete Menge roter Fälle |
| Größenmetrik (Summe Einfügungen und Löschungen, ohne `docs/plan`, `docs/reviews`) | 617 (Marke 1200, Meldeschwelle 900) |

Nicht gefahren: volle Gegenprobe (511 Läufe) und `make gates`; Belege des Implementers in Plan §7, nicht widerlegt.

## Behebungen des Reviews, von Hand gebrochen

| Finding | Ergebnis | Beleg (roter Fall) |
|---|---|---|
| F-1 Randformen | behoben (Punkte 11, 13, 16 entscheiden die neun Rückgaben) | Zuordnung Code zu Vertrag stichprobenweise |
| F-2 drei Wächter | behoben (ein `jq`-Aufruf der Auswertung, Wächter je mit Fall) | `p16-ergebnis-kein-array`, `p16-vulnerabilities-nicht-iterierbar` |
| F-3 Werkzeugliste | behoben | `p16-werkzeug-jq`, `p16-werkzeug-date` |
| F-4 Wiederholungs-Rat | behoben in `releasing.md` §11 (Punkt 18 b) | Text gelesen |
| F-5 Zukunft von `UpdatedAt` | behoben (unter −300 s = `datenbank`) | `p16-db-zukunft`, `p16-db-zukunft-toleranz` |
| F-6 Name `pgr-gp` | behoben | `grep` |
| F-7 (a) und (b) | behoben | Vertragstext, Plan §1 und DoD 1 |

Die sieben Fälle `p16-reihenfolge-*` stehen alle in der Gegenprobe, jeder mit einer Mutation, die genau seinen Fall färbt. Summenzeile bei Grün: `m-p16-gruen-zeile-fehlt` und `-auf-stderr` färben je acht Fälle; eigene Mutationen (`>&2` an der Grünzeile, `.FixedVersion != null`) färbten die erwarteten Fälle. Grün blieb keine Mutation.

## DoD-Konformität

| DoD-Punkt | Stand |
|---|---|
| 1 Scan (LH-FA-16) | erfüllt, soweit gegen den Diff prüfbar; Handbuch und README tragen nichts dazu, wie §1 es sagt |
| 2 Freshness-Audit | §5 beschreibt es; Ausführung offen bis Tag |
| 3 Checkliste §10 und erster Tag | Checkliste erfüllt; Tag offen |
| `make gates` grün | Beleg des Implementers: Exit 0 |
| Review-Report | erfüllt |
| Closure, Lerneintrag, Register, Risiko-Ausgänge, Paarungen | offen (Closure) |

## Findings

| ID | Kat. | Befund | Zuständig |
|---|---|---|---|
| V-1 | LOW | Tabelle in SPEC-067 Punkt 11 (Zeile `scan`: "wiederholen") und `releasing.md` §11 decken sich nicht; Punkt 18 (b) folgt die Doku | Architect (nachgezogen, d5e0b0f3) |
| V-2 | LOW | Plan §7 trägt bis zur Closure Platzhalter, Risiken in §6 "offen bis Closure"; Frage zum Release-Commit unbeantwortet | Closure |
| V-3 | INFO | Einleitungssatz der Randformen in §6 veraltet | Implementer bei Closure (§3.9) |
| V-4 | INFO | SPEC-067 Punkt 16 nennt Mutationen als Fälle (`p16-db-70-stunden`, `-74-stunden`) | Architect (nachgezogen, d5e0b0f3) |
| V-5 | INFO | Die Fälle prüfen den Scan gegen eine Attrappe von `docker`; die Annahmen (a) bis (f) über Trivy sind nicht belegt; Abweichung endet rot, nie grün | Trockenlauf |

## Offen bis Tag oder Trockenlauf (kein Mangel gegen diesen Diff)

- Trockenlauf des Scans auf GitHub: Messung der Annahmen (a) bis (f) und Lauf des Pins `aquasec/trivy@sha256:af6acf9a…855daa` (0.75.0).
- Freshness-Audit mit ausgefüllter Tabelle, Abfrage höchstens sieben Tage vor dem Tag.
- Das Ja des Nutzers zum Tag, der Release-Commit (`VERSION`, Zeile der Änderungshistorie), Gegenproben und `make gates` auf dem getaggten Hash.
- Der Lauf des Tags (Wortlaute von `skopeo` und `gh`, Schreibrecht, Sichtbarkeit, Smoke aus beiden Registries, Anhänge, `latest`).
- Der Handbuch- und README-Teil in einem eigenen Commit nach dem Tag.
- Closure, Lerneintrag, Register, Risiko-Ausgänge, Paarungen.
