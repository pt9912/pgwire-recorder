# ADR-0024: SQLite-Bibliothek ohne native Abhängigkeit

**Status:** Proposed

**Datum:** 2026-10-04

**Autor:** pt9912

**Bezug:** [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [ADR-0009](0009-implementierungssprache-go.md), [ADR-0013](0013-sqlite-recording-backend.md)

**Schärft:** [`LH-FA-22.a`](../../../spec/spezifikation.md#lh-fa-22a--aufzeichnungsformat), [`SPEC-035`](../../../spec/spezifikation.md#spec-035--zielplattformen), [`ARC-014`](../../../spec/architecture.md#3-externe-abhängigkeiten)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0013](0013-sqlite-recording-backend.md) führt SQLite als zweites Aufzeichnungsformat ein und lässt die Bibliothek offen, bevorzugt aber eine ohne native Abhängigkeit. Das Binary entsteht für sechs Plattformen (Linux, macOS, Windows je `amd64` und `arm64`), und alle Builds laufen in Containern ohne Host-Toolchain ([ADR-0009](0009-implementierungssprache-go.md)). Eine Bibliothek, die einen C-Compiler je Zielplattform verlangt, macht diese Builds schwerer und das Image größer.

## Entscheidung

Wir wählen `modernc.org/sqlite`, eine SQLite-Bibliothek in reinem Go ohne cgo. Sie wird nur im Recording-Adapter verwendet. Geprüft am 2026-10-04: aktuelle Version v1.60.1 (veröffentlicht am 2026-09-29), Lizenz BSD-3-Clause (SQLite selbst gemeinfrei), kein cgo, Windows, macOS und Linux jeweils für `amd64` und `arm64` unterstützt, das Repository wird aktiv gepflegt (über 800 Commits).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — `github.com/mattn/go-sqlite3` (cgo) | verbreitet, schnell | verlangt einen C-Compiler je Zielplattform; Cross-Builds für Windows und `arm64` aufwendig; widerspricht „ohne native Abhängigkeit" |
| B — `github.com/ncruces/go-sqlite3` (WebAssembly) | ohne cgo | zusätzliche Laufzeit für WebAssembly; weniger verbreitet |
| C — SQLite nicht verwenden | keine Zusatzabhängigkeit | widerspricht [`LH-FA-22`](../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat) |
| **D — `modernc.org/sqlite` (reines Go)** | ein Build ohne C-Toolchain für alle Plattformen; passt zur Docker-only-Regel | langsamer als die C-Bibliothek; größeres Binary |

## Konsequenzen

- Positiv: Cross-Builds für alle Zielplattformen bleiben einfach; das Image braucht keine C-Laufzeit.
- Negativ: Die Bibliothek ist nach eigener Angabe etwa 1,3- bis 2-mal langsamer als native SQLite; die Binärgröße wächst; sie hängt von `modernc.org/libc` in genau der Version ab, die ihre `go.mod` festlegt (Abhängigkeit exakt pinnen).
- Folgepflicht: Eine `tech`-Regel in `.a-check.yml` beschränkt die Bibliothek auf den Recording-Adapter; der Roundtrip-Test beider Formate (siehe [ADR-0013](0013-sqlite-recording-backend.md)) läuft mit ihr.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check | `modernc.org/sqlite` nur im Recording-Adapter (`tech`-Regel, mit der Annahme) | `make a-check` |

## Re-Evaluierungs-Trigger

Wenn die Geschwindigkeit bei großen Aufzeichnungen nicht genügt, die Bibliothek nicht mehr gepflegt wird oder ihre Lizenz sich ändert.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
