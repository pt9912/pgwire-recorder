# Verifikationsbericht: slice-erster-release-freigabe, Runde 3 (Modul 11)

**Rolle:** Verifier · **Datum:** 2026-10-19 · **Modell:** claude-sonnet-5-5
**Gegenstand:** Diff b307b020..8c852996 (Zweig `ver-fr2`), Metrik-Basis 760c51d4.
**Maßstab:** Plan (DoD, §6, §7), SPEC-067 Punkte 11, 13, 16, 18, [ADR-0052](../plan/adr/0052-scan-des-images-vor-dem-push.md), Review-Report [Runde 2](2026-10-19-review-slice-erster-release-freigabe-r2.md).

## Verdikt

Freigabe zur Closure, ohne HIGH und ohne MEDIUM. **Summary-Zeile:** `slice-erster-release-freigabe (Verifikation Runde 3): 0 HIGH · 0 MEDIUM · 2 LOW · 2 INFO`

## Behebungen von Hand gebrochen

| Punkt | Ergebnis |
|---|---|
| F-1 Cache-Mount | `m-p16-cache-ueber-laeufe` färbt genau `p16-cache-ueber-laeufe`; eigener Mutant (Cache außerhalb des Temp-Verzeichnisses) ebenfalls |
| F-2 Marke an drei Exit-2-Stellen | `m-p16-marke-loescht-nach-entpacken`, `-nach-index`, `-nach-auswahl` färben die ineinander enthaltenen Mengen (1, 2, 3 Fälle) |
| F-3 `blobs/` | `m-p16-blobs-leer` färbt `p16-blobs-inhalt`; eigener Mutant (nur Index-Blob) ebenfalls |
| F-4 | `m-p16-index-mehr-als-zwei-abgelehnt` färbt `p16-index-attestation`, `p16-index-fremdes-os`, `p16-manifest-doppelt-variante`; Zeile in §7 vorhanden |
| F-5 Variante | `m-p16-auswahl-variante-ausgeschlossen` färbt `p16-index-variante`; `m-p16-auswahl-gleiche-variante-doppelt` färbt `p16-manifest-doppelt-variante` |
| F-6 (1) | `m-p16-entpacken-je-plattform` färbt `p16-entpacken-einmal` |
| F-6 (2) `manifests[-1]` | eigener Mutant grün; Äquivalenz nach Invariante (Punkt 5), Ort (`# P5ein-index`, `p5-layout-form`, `m-p5-ein-index-entfernt`) und Sonde tragfähig |
| F-7, F-8 | ohne Befund |

Weitere Mutationen der Gruppe färbten genau die erwarteten Fälle. Größenmetrik 723 (Marke 1200, Meldeschwelle 900).

## Findings

| ID | Kat. | Befund | Zuständig |
|---|---|---|---|
| F-V1 | LOW | §7 sagte „kein Mutant gilt als äquivalent“ neben der Äquivalenz von `manifests[-1]` | behoben im Plan (Orchestrator) |
| F-V2 | LOW | Inhalt von `oci-layout`, `index.json` und `blobs/` als Symlink nicht durch Fälle belegt (Mutanten grün); fail-closed (Trivy rot) | akzeptiertes Negativ, in §7 geführt |
| F-V3 | INFO | Sonde zu `manifests[-1]` läuft nicht durch `scan` mit Zwei-Eintrags-Layout | keiner |
| F-V4 | INFO | Beschreibung von `m-p16-cache-ueber-laeufe` in §7 nennt `$RUNNER_TEMP/release`, die Mutation `$zdir/trivy` | optional |

## Eigene Läufe

`bash -n` beider Skripte grün; `make docs-check` 0 Befunde (1248 Dateien); Mutationen der Gruppe `f` grün. Belege des Implementers, nicht wiederholt: `make release-gegenprobe` 534 Läufe grün, `make workflows-check`, `make gates` Exit 0.

## Offen bis Tag oder Trockenlauf

- Annahme (e) (`Results` für das Produkt-Image nicht leer) bis zum zweiten Trockenlauf.
- Freshness-Audit, Checkliste, Release-Commit, Gegenproben auf dem Hash, Lauf des ersten Tags, Handbuch-/README-Teil danach, Closure.
