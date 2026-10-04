# ADR-0025: Benannte Slice-Kennungen im Commit-Hook

**Status:** Accepted

**Datum:** 2026-10-04

**Autor:** pt9912

**Bezug:** [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit)

**Schärft:** —

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Commit-Hook verlangt in jeder Commit-Message eine Kennung. Die mitgelieferte Prüfung kennt für Slices nur die Kennung mit Ziffern. Dieses Repo vergibt Slice-Kennungen als Namen aus Welle und Aspekt, und `harness/README.md` sagt jede Slice-Kennung als zulässig zu. Das Lifecycle-Werkzeug (Make-Ziel für den Wechsel eines Slice) committet mit dem Slice-Namen als einziger Kennung und scheitert deshalb am Hook; Moves müssen von Hand laufen. Die mitgelieferte Prüfung wird bei jedem Bootstrap neu geschrieben und ist kein Ort für eine Repo-Regel; der Träger `.githooks/commit-msg` gehört dem Repo.

## Entscheidung

Wir wählen: Der Träger `.githooks/commit-msg` lässt eine Commit-Message zusätzlich durch, wenn sie die Kennung eines vorhandenen Slice trägt, also den Namen einer Slice-Datei unter `docs/plan/planning/`. Alle anderen Fälle reicht er unverändert an die mitgelieferte Prüfung weiter.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — So lassen, Moves von Hand | keine Änderung am Hook | das Lifecycle-Werkzeug ist unbenutzbar; die Zusage in `harness/README.md` stimmt nicht |
| B — Slices künftig mit Ziffern benennen | Hook bleibt unverändert | alle bestehenden Slices, Wellen und Verweise wären umzubenennen |
| C — Jedes Wort mit dem Slice-Präfix durchlassen | einfachste Änderung | auch Tippfehler und der Name des Lifecycle-Werkzeugs gälten als Kennung |
| **D — Nur Namen vorhandener Slices durchlassen** | das Lifecycle-Werkzeug funktioniert; ein erfundener Name fällt weiter durch | der Hook liest das Planungsverzeichnis |

## Konsequenzen

- Positiv: Lifecycle-Wechsel laufen über das Werkzeug mit Verweis-Nachzug; die Zusage in `harness/README.md` gilt.
- Negativ: Ein Commit, der einen Slice löscht oder umbenennt, nennt ihn nicht mehr als vorhandene Datei; er braucht eine andere Kennung.
- Folgepflicht: Der Träger trägt die Regel; eine Gegenprobe zeigt, dass ein erfundener Slice-Name abgelehnt und ein vorhandener angenommen wird.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Gegenprobe des Trägers | erfundener Slice-Name wird abgelehnt, vorhandener angenommen, Message ohne Kennung abgelehnt | `make gates` |

## Re-Evaluierungs-Trigger

Wenn die mitgelieferte Prüfung benannte Slice-Kennungen selbst kennt oder das Repo auf nummerierte Slices wechselt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Proposed | — |
| 2026-10-04 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
