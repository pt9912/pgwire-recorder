# Roadmap

**Format-Regel:** Die Roadmap ist eine Reihenfolge von **Wellen**,
keine Reihenfolge von Terminen (siehe
Baseline-Regelwerk `modul-06-roadmap.md`).
Termine werden — falls überhaupt — als Konsequenz der Wellen-Schätzung
gezeigt, nicht als Treiber.

---

## Offene Wellen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte — *Offene Wellen* ist **derivativ**: Der
Zustand sind die flachen Welle-Dateien; woran gearbeitet wird, sagt das
`Welle:`-Feld der Slices in `in-progress/`. Ziel, Trigger und
Closure-Kriterien stehen in der Welle-Datei, nicht hier.

- [welle-v1-abschluss](../welle-v1-abschluss.md)
- [welle-erster-release](../welle-erster-release.md)

In Arbeit: `slice-lastenheft-pruefbarkeit` (wellenlos).



## Nächste Wellen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte, Bullet *Nächste Wellen* — die geordnete
Vorschau: je Zeile Welle, Trigger als beobachtbare Bedingung, wichtigste Slices
und geschätzter Aufwand (S/M/L, kein Termin).

Keine; alle geplanten Wellen haben eine Datei unter *Offene Wellen*.

## Meilensteine

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Welle ≠ Meilenstein ≠ Release.



| Meilenstein | Welle(n) | Trigger | Status |
|---|---|---|---|
| M1 | welle-walking-skeleton | `SELECT 1;` Record, PostgreSQL stoppen, Replay mit demselben Ergebnis | erreicht — [welle-walking-skeleton-results](../done/welle-walking-skeleton-results.md) |
| M2 | welle-extended-query | Clients mit Standardtreibern (Extended Query) funktionieren mit Host/Port-Umstellung: Abnahmeszenario 7 nachweisbar | erreicht — [welle-extended-query-results](../done/welle-extended-query-results.md) |
| M3 | welle-replay-semantik, welle-v1-abschluss | Produkt fertig: alle Anforderungen (MUSS und SOLL) umgesetzt, die Abnahmeszenarien 1 bis 10 und 12 bis 17 nachweisbar | offen |
| M4 | welle-erster-release | erster Release veröffentlicht, Abnahmeszenario 11 (Homebrew) nachgewiesen | offen |

## Abhängigkeitsgraph

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte, Bullet *Nächste Wellen* — die Abhängigkeit
steht als beobachtbare Bedingung in der `Trigger`-Spalte **und** als gerichtete
Kante hier; eine Welle, die ohne fertige Vorgängerin nicht starten kann, ist
eine Phantom-Welle.

```mermaid
flowchart LR
    W1[welle-walking-skeleton]
    W2[welle-extended-query]
    W3[welle-replay-semantik]
    W4[welle-v1-abschluss]
    W5[welle-erster-release]

    W1 --> W2
    W2 --> W3
    W3 --> W4
    W4 --> W5
```

## Abgeschlossene Wellen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

| Welle | Abschluss | Closure-Notiz |
|---|---|---|
| welle-walking-skeleton | 2026-10-04 | [welle-walking-skeleton-results](../done/welle-walking-skeleton-results.md) |
| welle-extended-query | 2026-10-05 | [welle-extended-query-results](../done/welle-extended-query-results.md) |
| welle-replay-semantik | 2026-10-06 | [welle-replay-semantik-results](../done/welle-replay-semantik-results.md) |

## Historische Trigger-Verschiebungen

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte, Bullet *Historische Trigger-Verschiebungen*
— das Drift-Log: jede Umplanung mit Datum, Änderung, Grund. Leer heißt starre
Roadmap, jede Zeile voll heißt treibende.



| Datum | Was wurde geändert? | Warum? |
|---|---|---|
| 2026-10-03 | welle-erster-release angelegt; Abnahmeszenario 11 (Homebrew) wandert von M3 nach M4 | Die Homebrew-Formel entsteht erst bei einem veröffentlichten, stabilen Release; M3 und der erste Release bildeten einen Zirkel |
| 2026-10-03 | welle-extended-query vor welle-replay-semantik gezogen; Start-Trigger von welle-replay-semantik von „welle-walking-skeleton done“ auf „welle-extended-query done“ geändert; Meilensteine M2 und M3 neu zugeschnitten | Extended Query gehört zu v1: verbreitete Treiber nutzen es standardmäßig, ohne es genügt die Umstellung von Host und Port nicht |
| 2026-10-03 | welle-v1-abschluss um `slice-v1-abschluss-tls-client` und `slice-v1-abschluss-antwortvergleich` erweitert; M3 umfasst die Abnahmeszenarien 12 bis 16 | Neue SOLL-Anforderungen LH-FA-23 (TLS zum Client) und LH-FA-24 (Antwortvergleich) aus dem Bedarf eines E2E-Einsatzes |
| 2026-10-04 | welle-v1-abschluss um `slice-v1-abschluss-anmeldung` erweitert | Der Record-Modus vermittelt im Walking Skeleton keine Passwort-Anmeldung; LH-FA-05.b verlangt sie |
| 2026-10-05 | welle-extended-query um `slice-extended-query-lebendpruefung` erweitert; M2 schließt erst danach | Validierung von `slice-extended-query-replay`: Die Lebendprüfung von Connection-Pools (`-- ping` nach Leerlauf) macht das strenge Replay zeitabhängig, und die Zusage „nur Host und Port“ aus LH-FA-18 hält für die häufigste Einsatzform nicht; der Nutzer hat entschieden, solche Anfragen zu tolerieren |
| 2026-10-05 | welle-v1-abschluss um `slice-v1-abschluss-postgres-versionen` erweitert | Die Validierung von `slice-extended-query-replay` und [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md) machen die Serverversion relevant: Das Replay entscheidet selbst, was PostgreSQL als leere Anfrage behandelt, und das hängt am Scanner des Servers; die Spezifikation nennt keine Version, getestet wird nur gegen PostgreSQL 17. Der Nutzer hat entschieden, die Hauptversionen der letzten fünf Jahre zu unterstützen |
| 2026-10-05 | `slice-v1-abschluss-postgres-versionen` startet vor dem Welle-Trigger von welle-v1-abschluss, sobald `slice-extended-query-lebendpruefung` done ist | Entscheidung des Nutzers: Die Lebendprüfung erkennt Leerraum jetzt abhängig von der Serverversion, und die Zusage über die unterstützten Versionen soll nicht bis nach welle-replay-semantik ungeprüft bleiben |
| 2026-10-05 | `slice-harness-lint` und `slice-harness-coverage` (wellenlos) angelegt; sie starten nach `slice-replay-semantik-fehlerreplay` und vor dem nächsten großen Slice, in dieser Reihenfolge (WIP-Limit 1) | Entscheidung des Nutzers: SOLID-naher Lint und eine Schwelle für die Testabdeckung vor dem nächsten großen Slice |
| 2026-10-06 | `slice-lastenheft-pruefbarkeit` und vier Umstellungs-Slices (`slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`, `slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`) wellenlos angelegt; Reihenfolge vor dem nächsten großen Slice: Lastenheft, `slice-harness-lint`, die vier Umstellungs-Slices, `slice-harness-coverage` | Entscheidung des Nutzers vom 2026-10-05: Tests werden auf Black-Box-Pakete umgestellt (bis dahin `testpackage` gestuft), die Profil-Schwellen des Vorbilds gelten, und das Lastenheft bekommt eine Qualitätsanforderung an die Prüfbarkeit des Quellcodes |
| 2026-10-06 | M3 um Abnahmeszenario 17 (Prüfbarkeit des Quellcodes) erweitert; es entsteht mit `slice-lastenheft-pruefbarkeit` und wird mit `slice-harness-coverage` nachweisbar, sodass die wellenlose Reihe von Lastenheft bis Coverage vor M3 liegt | Entscheidung des Nutzers vom 2026-10-06: die neue Qualitätsanforderung hat Priorität MUSS; das Team hat ein eigenes Abnahmeszenario entschieden (alle drei Prüfungen grün, jede nachweislich rot) |
| 2026-10-06 | `slice-harness-mutation` (wellenlos) angelegt; die wellenlose Reihe endet nach `slice-harness-coverage` mit ihm | Entscheidung des Nutzers vom 2026-10-06 bei der Closure von welle-replay-semantik: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` und `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` sind als Prosa ausgeschöpft und bekommen ein Mutations-Gate als Sensor; M3 hängt nicht an ihm |
| 2026-10-06 | `slice-harness-abdeckung-gate` (wellenlos) angelegt; Reihenfolge der wellenlosen Reihe nach dem Lastenheft: `slice-harness-lint`, die vier Umstellungs-Slices, `slice-harness-abdeckung-gate`, `slice-harness-coverage`, `slice-harness-mutation` | Entscheidung des Nutzers vom 2026-10-06: Der Nachweis von LH-QA-07 in den Abdeckungstabellen nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) braucht eine Erweiterung des Abdeckungs-Skripts in einem eigenen Slice; er steht nach der Black-Box-Umstellung, weil Teil 3 erst gilt, wenn `testpackage` überall scharf ist, und vor Coverage, das Teil 2 deklariert |
| 2026-10-06 | `slice-harness-vertraege-spezifikation` (wellenlos) angelegt; er steht in der wellenlosen Reihe nach `slice-lastenheft-pruefbarkeit` und vor `slice-harness-lint` | Antwort des Kurs-Repos ai-harness-course auf den Change Request „Spezifikations-Ort für Verträge der Harness-Werkzeuge“: Der Ort ist die Spezifikation; der Slice legt den Abschnitt an und verlegt die Verträge von `kopf-check` und `abdeckung` dorthin, bevor Lint, Umstellung, Abdeckungs-Erweiterung, Coverage und Mutation ihre Randformen dort entscheiden |
