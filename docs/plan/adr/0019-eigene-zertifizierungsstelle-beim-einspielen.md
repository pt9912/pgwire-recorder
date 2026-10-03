# ADR-0019: Eigene Zertifizierungsstelle beim Einspielen

**Status:** Proposed

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [ADR-0016](0016-einspielen-anmeldung-und-tls.md)

**Schärft:** [`LH-FA-20.a`](../../../spec/spezifikation.md#lh-fa-20a--einspielen), [`LH-FA-17.a`](../../../spec/spezifikation.md#lh-fa-17a--konfiguration)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0016](0016-einspielen-anmeldung-und-tls.md) prüft das Serverzertifikat gegen den Zertifikatsspeicher des Systems und schließt das Überspringen aus. Dessen Re-Evaluierungs-Trigger („eigene Zertifizierungsstelle pro Aufruf") ist eingetreten: Testserver und Adapter tragen Zertifikate einer eigenen Zertifizierungsstelle, die nicht im Systemspeicher liegt, und im Container lässt sich der Systemspeicher nicht bequem ändern. Die Entscheidung von 0016 bleibt gültig; diese ADR ergänzt sie.

## Entscheidung

Wir wählen: `play` kennt die Option `--upstream-ca` (PEM-Datei mit einem oder mehreren Zertifikaten). Die Zertifikate **ergänzen** den Systemspeicher, ersetzen ihn nicht. Die Prüfung der Kette und des Host-Namens bleibt unverändert; ein Überspringen gibt es weiterhin nicht. Die Option gilt nur mit TLS zum Server (`PGR-E2001` sonst), eine unlesbare oder ungültige Datei ist `PGR-E2007`.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nichts tun, Zertifizierungsstelle im Systemspeicher ablegen | keine neue Option | im Container und in CI umständlich; ändert das Image oder den Host |
| B — Prüfung abschaltbar machen | einfachste Bedienung | schwächt die Sicherheit; widerspricht [ADR-0016](0016-einspielen-anmeldung-und-tls.md) |
| C — Eigene Zertifizierungsstelle ersetzt den Systemspeicher | strengere Vertrauensbasis | ein Server mit öffentlichem Zertifikat wäre mit der Option nicht mehr prüfbar |
| **D — Eigene Zertifizierungsstelle ergänzt den Systemspeicher** | Prüfung bleibt vollständig; eine Option für beide Fälle | das Vertrauen wächst um die genannten Zertifikate |

## Konsequenzen

- Positiv: Einspielen gegen Testserver mit eigenem Zertifikat ohne Änderung am System.
- Negativ: Der Anwender entscheidet über zusätzliches Vertrauen und trägt die Verantwortung für die Datei.
- Folgepflicht: Handbuch beschreibt die Option; Tests decken gültige, ungültige und fehlende Datei ab.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Tests | `play` kennt keine Option, die die Zertifikatsprüfung abschaltet | `make gates` |

## Re-Evaluierungs-Trigger

Wenn Client-Zertifikate für `play` gefordert werden oder die Zertifizierungsstelle je Verbindung statt je Aufruf gelten muss.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
