# ADR-0035: Struktur-Kennungen im Commit-Hook abgelehnt

**Status:** Accepted

**Datum:** 2026-10-08

**Autor:** pt9912

**Bezug:** [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md), [ADR-0029](0029-benannte-welle-kennungen-im-commit-hook.md)

**Schärft:** [`SPEC-050`](../../../spec/spezifikation.md#spec-050--struktur-kennungen-im-commit-träger-commit-msg)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`AGENTS.md` §5 Regel 1 sagt, dass Struktur-Kennungen der Spezifikation und der Sicht nicht in die Commit-Message gehören. Die Regel hat keinen Träger; das Review hat das Muster dreimal erst nach dem Commit gefunden (`BEO-REPO/commit-nennt-struktur-kennung`). Der Träger `.githooks/commit-msg` gehört dem Repo und bleibt bei einem erneuten Bootstrap stehen, die mitgelieferte Prüfung wird dabei neu geschrieben ([ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md)). [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md) und [ADR-0029](0029-benannte-welle-kennungen-im-commit-hook.md) entscheiden, dass der Träger eine Message mit benannter Slice- oder Welle-Kennung durchlässt und jede andere unverändert weiterreicht; eine Ablehnung ändert beides für Messages mit Struktur-Kennung.

**Bestand** (Stand `25325be`, 489 Commits): 27 Messages nennen eine Struktur-Kennung, alle in Großschreibung mit drei Ziffern im Fließtext oder im Betreff. Keine nennt eine kleingeschriebene Form, einen Anker oder einen Platzhalter, und keine ist eine Merge- oder Revert-Message. Fünf Betreffs tragen eine Struktur-Kennung; ein Revert eines dieser Commits zitierte sie.

## Entscheidung

Wir wählen **F**: Der Träger `.githooks/commit-msg` lehnt eine Message ab, die eine Struktur-Kennung nennt, vor jeder Annahme nach [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md) und [ADR-0029](0029-benannte-welle-kennungen-im-commit-hook.md). Gezählt wird nur die genaue Schreibweise als ganzes Wort, im selben Lesebereich wie die Annahme; Merge- und Revert-Messages sind ausgenommen. Die Einzelregeln stehen in [`SPEC-050`](../../../spec/spezifikation.md#spec-050--struktur-kennungen-im-commit-träger-commit-msg).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — So lassen, das Review findet es | keine Änderung | dreimal erst nach dem Commit gefunden; ein gepushter Commit bleibt |
| B — Ablehnung in der mitgelieferten Prüfung | eine Prüfung statt zwei Stellen | ein erneuter Bootstrap schreibt sie neu, die Ablehnung ginge still verloren |
| C — eigenes Skript, das der Träger ruft | Träger bleibt kurz | eine Datei mehr für eine Regel, deren Ort dem Repo schon gehört |
| D — jede Schreibweise zählt (auch `spec-049`, Anker) | fängt mehr | sperrt Links mit Anker auf die Spezifikation, die eine Message legitim tragen darf; im Bestand nie vorgekommen |
| E — ohne Ausnahme für Merge und Revert | keine Lücke | ein Revert eines der fünf Commits mit Struktur-Kennung im Betreff scheiterte an dem Betreff, den git zitiert |
| **F — im Träger, genaue Schreibweise, vor der Annahme, Merge und Revert ausgenommen** | überlebt den Bootstrap; fängt jede Form aus dem Bestand; Reverts bleiben möglich | eine selbst geschriebene Message mit Betreff `Merge ` oder `Revert ` bleibt ungeprüft |

## Konsequenzen

- Positiv: Eine Struktur-Kennung in der Message fällt im aktivierten Klon beim Commit auf, nicht im Review.
- Negativ (akzeptiert): Ein Klon ohne Aktivierung, `--no-verify`, ein anderes Kommentarzeichen und eine selbst geschriebene Message mit Betreff `Merge ` oder `Revert ` bleiben ungeprüft; das Review bleibt zweiter Leser. Bestehende Commits bleiben, wie sie sind.
- Folgepflicht: Die Gegenprobe des Trägers hält je Punkt von [`SPEC-050`](../../../spec/spezifikation.md#spec-050--struktur-kennungen-im-commit-träger-commit-msg) einen Fall; ihre vorhandenen Fälle bleiben mit ihrer Erwartung.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Gegenprobe des Trägers | Message mit Struktur-Kennung abgelehnt, auch neben einer angenommenen Kennung; andere Schreibweisen, Kommentarzeilen, Diff unter der Scissors-Zeile und Merge/Revert angenommen bzw. wie ohne diese Regel | `make gates` |

## Re-Evaluierungs-Trigger

Wenn `AGENTS.md` §5 Regel 1 Struktur-Kennungen zulässt, wenn die Spezifikation oder die Sicht ein anderes Kennungs-Schema einführt, oder wenn eine Merge- oder Revert-Message mit selbst geschriebener Struktur-Kennung im Review gefunden wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-08 | Proposed | — |
| 2026-10-08 | Accepted (Entscheidung des Nutzers) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
