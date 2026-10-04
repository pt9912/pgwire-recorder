# ADR-0029: Benannte Welle-Kennungen im Commit-Hook

**Status:** Proposed

**Datum:** 2026-10-04

**Autor:** pt9912

**Bezug:** [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md)

**Schärft:** —

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Der Commit-Hook verlangt in jeder Commit-Message eine Kennung. Seit [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md) nimmt der Träger `.githooks/commit-msg` die Namen vorhandener Slices an; Wellen-Namen nimmt er nicht an. Das Archivierungs-Werkzeug (Make-Ziel `archive-welle`) verschiebt die Zeitdokumente einer geschlossenen Welle und committet mit dem Welle-Namen als einziger Kennung. Bei der Closure der ersten Welle scheiterte dieser Commit am Hook; das Werkzeug ist im Repo nicht benutzbar. Die mitgelieferte Prüfung wird bei jedem Bootstrap neu geschrieben und ist kein Ort für eine Repo-Regel.

## Entscheidung

Wir wählen **D**: Der Träger `.githooks/commit-msg` nimmt zusätzlich zu den Slice-Namen nach [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md) den Namen einer vorhandenen Welle an, also den Namen einer Welle-Datei im Index unter `docs/plan/planning/`, flach, in `done/` oder im Archiv-Ordner der Welle unter `done/`. Die Closure-Notiz einer Welle (Endung `-results.md`) ist keine Welle-Datei. Alle anderen Fälle gehen unverändert an die mitgelieferte Prüfung.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — So lassen, Archivierung von Hand | keine Änderung am Hook | das Werkzeug ist unbenutzbar; die Verweis-Nachführung des Werkzeugs entfällt |
| B — Commit des Werkzeugs mit `--no-verify` | keine Änderung am Hook | umgeht den Wächter an genau der Stelle, die er schützen soll |
| C — Jedes Wort mit dem Welle-Präfix durchlassen | einfachste Änderung | Tippfehler und erfundene Namen gälten als Kennung |
| **D — Nur Namen vorhandener Wellen durchlassen** | das Werkzeug funktioniert; ein erfundener Name fällt weiter durch; gleiche Form wie [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md) | der Hook liest auch den Archiv-Ordner unter `done/` |

## Konsequenzen

- Positiv: `make archive-welle` läuft mit Verweis-Nachführung durch; der Commit-Hook bleibt für alle anderen Messages unverändert.
- Negativ: Ein Welle-Name ist eine schwächere Kennung als eine Anforderung; ein Commit, der nur ihn trägt, sagt nicht, welche Zusage berührt ist. Das gilt nur für Moves und Archivierung.
- Folgepflicht: Der Träger trägt die Regel; die Gegenprobe zeigt, dass ein erfundener Welle-Name und der Name einer Closure-Notiz abgelehnt und eine vorhandene Welle angenommen wird.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Gegenprobe des Trägers | erfundener Welle-Name und Name einer Closure-Notiz abgelehnt, vorhandene Welle (flach, in `done/`, archiviert) angenommen | `make gates` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
Wenn die mitgelieferte Prüfung benannte Welle-Kennungen selbst kennt oder das Archivierungs-Werkzeug mit einer anderen Kennung committet.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
