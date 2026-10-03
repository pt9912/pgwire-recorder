# ADR-0014: Konfigurationsdatei mit benannten Verbindungen

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0011](0011-meldungscodes-praefix-pgr.md)

**Schärft:** [`LH-FA-17.a`](../../../spec/spezifikation.md#lh-fa-17a--konfiguration), [`SPEC-034`](../../../spec/spezifikation.md#spec-034--meldungscodes)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Optionen von `record`, `replay` und `play` sind zahlreich; wiederkehrende Läufe (CI, lokale Entwicklung) brauchen eine Datei, die sie festhält und benannte Verbindungen zu PostgreSQL trägt. Verbindungen enthalten Passwörter. Eine Datei, die im Repository liegt, darf keine Geheimnisse tragen. Fehler beim Laden müssen sich unterscheiden lassen, weil CI sie unterschiedlich behandelt.

## Entscheidung

Wir wählen **eine YAML-Datei** `.pgwire-recorder.yaml`, gewählt aus `--config`, `PGWIRE_RECORDER_CONFIG` oder dem aktuellen Verzeichnis (in dieser Reihenfolge, genau eine Datei, keine Zusammenführung). Die Schlüssel heißen wie die Optionen; `log_level` und `connections:` stehen oben, alle anderen in einem Abschnitt je Kommando. Verbindungen sind URLs mit dem einzigen Parameter `sslmode` (`disable`, `require`). Geheimnisse stehen nie in der Datei: ein Passwort ist nur als Platzhalter `${VAR}` erlaubt, der aus der Umgebung aufgelöst wird. Ein Klartext-Passwort und jede unbekannte Eingabe sind Startfehler mit eigenem Meldungscode (`PGR-E2004` bis `PGR-E2006`). Der Rang bei Mehrfachangabe lautet CLI-Argument, Umgebungsvariable, Datei, Standardwert. `config show` zeigt die Datei mit unaufgelösten Platzhaltern.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Keine Datei, nur CLI und Umgebung | kein Dateiformat, keine Ladefehler | lange Aufrufe in CI-Skripten; benannte Verbindungen nicht möglich |
| B — Mehrere Dateien zusammenführen (Benutzer, Projekt) | gemeinsame Vorgaben, projektweise Überschreibung | Herkunft eines Werts schwer zu erkennen; mehr Fehlerfälle |
| C — TOML oder JSON statt YAML | strengere Syntax | anderes Format als das der Aufzeichnung (YAML); Platzhalter und Kommentare in JSON nicht möglich |
| **D — Eine YAML-Datei, Platzhalter für Geheimnisse** | ein Format im ganzen Werkzeug; Herkunft eindeutig; Geheimnisse außerhalb des Repositorys | Platzhalter-Syntax muss gelernt werden; eine Datei erlaubt keine Schichtung |

## Konsequenzen

- Positiv: Läufe sind aus einer versionierbaren Datei wiederholbar; ein Passwort gelangt nicht in die Datei.
- Negativ: Umgebungsvariablen für Platzhalter müssen gesetzt sein; Fehler im Laden beenden den Start.
- Folgepflicht: Handbuch beschreibt Dateiformat und Fehlercodes; Slice für die Konfiguration trägt die Tests je Fehlercode.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Tests | jede Fehlerursache der Spezifikation liefert ihren Code, die Meldung enthält nie einen Wert | `make gates` |

## Re-Evaluierungs-Trigger

Wenn Einstellungen aus mehreren Quellen zusammengeführt werden müssen (Benutzer- und Projektdatei) oder ein weiterer Parameter einer Verbindung gefordert wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
