# ADR-0027: YAML-Bibliothek go.yaml.in/yaml/v3

**Status:** Accepted

**Datum:** 2026-10-04

**Autor:** pt9912

**Bezug:** [`LH-FA-07`](../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-QA-06`](../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), [ADR-0005](0005-recording-store-ist-driven-adapter.md), [ADR-0009](0009-implementierungssprache-go.md)

**Schärft:** [`SPEC-001`](../../../spec/spezifikation.md#spec-001--recording-serialisierung-und-versionierung), [`ARC-013`](../../../spec/architecture.md#3-externe-abhängigkeiten)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Format `yaml` der Aufzeichnung braucht eine YAML-Bibliothek im Recording-Adapter. Die verbreitete Bibliothek `gopkg.in/yaml.v3` wird nicht mehr gepflegt; ihre Pflege hat das YAML-Projekt unter einem neuen Modulpfad übernommen. Die Bibliothek muss deterministisch schreiben (`SPEC-004`), unbekannte Schlüssel ablehnen können und ohne cgo bauen.

## Entscheidung

Wir wählen `go.yaml.in/yaml/v3`, den gepflegten Nachfolger von `gopkg.in/yaml.v3` mit derselben Schnittstelle. Sie steht nur im Recording-Adapter; das Architektur-Gate beschränkt beide Modulpfade auf ihn.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — `gopkg.in/yaml.v3` | am weitesten verbreitet | nicht mehr gepflegt; Sicherheitskorrekturen kommen nicht mehr |
| B — `sigs.k8s.io/yaml` | verbreitet im Kubernetes-Umfeld | Umweg über JSON; Reihenfolge und Typen der YAML-Darstellung schwerer zu steuern |
| C — eigenes Format statt YAML | keine Abhängigkeit | widerspricht `SPEC-001` (YAML als Standardformat) |
| **D — `go.yaml.in/yaml/v3`** | gepflegt; gleiche Schnittstelle wie A; Node-API für das Lesen ungequoteter `null`-Schlüssel | jüngerer Modulpfad, weniger verbreitet |

## Konsequenzen

- Positiv: Die Abhängigkeit wird gepflegt; der Wechsel von A ist eine Änderung des Importpfads.
- Negativ: Sie schreibt YAML-1.1-Schlüsselwörter wie `null` oder `y` in Anführungszeichen; Leser und Beispiele müssen beide Formen tragen.
- Folgepflicht: Die `tech`-Regel in `.a-check.yml` führt beide Modulpfade, damit auch der alte Pfad außerhalb des Recording-Adapters auffällt.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | `go.yaml.in/yaml` und `gopkg.in/yaml` nur im Recording-Adapter (`tech`-Regel) | `make a-check` |

## Re-Evaluierungs-Trigger

Wenn die Bibliothek nicht mehr gepflegt wird oder eine YAML-Version gefordert wird, die sie nicht trägt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Proposed | — |
| 2026-10-04 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
