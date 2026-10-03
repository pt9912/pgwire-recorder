# ADR-0016: Einspielen: Anmeldung und TLS als Client

**Status:** Accepted

**Datum:** 2026-10-03

**Autor:** pt9912

**Bezug:** [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0004](0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](0010-verwendung-von-pgproto3.md), [ADR-0014](0014-konfigurationsdatei.md)

**Schärft:** [`LH-FA-20.a`](../../../spec/spezifikation.md#lh-fa-20a--einspielen), [`LH-FA-17.a`](../../../spec/spezifikation.md#lh-fa-17a--konfiguration)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Im Record-Modus vermittelt der Recorder die Anmeldung des Clients und kennt kein TLS zum Upstream. Beim Einspielen ist der Recorder selbst Client: er muss sich anmelden und darf eine TLS-Verbindung verlangen. Die Aufzeichnung enthält die Anmeldung nicht als Geheimnis, und ein Passwort gehört weder in die Aufzeichnung noch in eine Option, die in der Prozessliste steht.

## Entscheidung

Wir wählen: `play` meldet sich mit Klartext-Passwort, MD5 oder SCRAM-SHA-256 an. Das Passwort kommt aus dem Platzhalter der benutzten benannten Verbindung ([ADR-0014](0014-konfigurationsdatei.md)), sonst aus `PGWIRE_RECORDER_PASSWORD`; eine Option dafür gibt es nicht. Mit TLS (`--upstream-tls` oder `sslmode=require`) wird das Serverzertifikat gegen den Zertifikatsspeicher des Systems geprüft; ein Überspringen der Prüfung gibt es nicht. Eine fehlgeschlagene Anmeldung, ein nicht unterstütztes Verfahren, abgelehntes TLS und ein ungültiges Zertifikat sind `PGR-E4005`, alle anderen Fehler beim Aufbau `PGR-E4002`.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Nur Verbindungen ohne Anmeldung und ohne TLS | einfachste Umsetzung | taugt nicht für Testsysteme mit Anmeldung |
| B — Passwort als Option und Prüfung des Zertifikats abschaltbar | bequem im Test | Passwort in der Prozessliste; stille Abschwächung der Sicherheit |
| C — Passwort aus der Aufzeichnung | kein Zusatzaufwand | Geheimnisse in der Aufzeichnung widersprechen deren Weitergabe |
| **D — Passwort nur aus Umgebung oder Platzhalter, strenge Zertifikatsprüfung** | keine Geheimnisse in Aufruf und Datei; ein einheitlicher Sicherheitsstand | Testserver mit selbst signiertem Zertifikat brauchen eine Zertifizierungsstelle im Systemspeicher |

## Konsequenzen

- Positiv: Das Einspielen gegen reale Testsysteme ist möglich, ohne Geheimnisse preiszugeben.
- Negativ: Selbst signierte Zertifikate lassen sich nur über den Zertifikatsspeicher des Systems nutzen.
- Folgepflicht: Handbuch nennt die Wege zum Passwort und zum Zertifikatsspeicher; Tests decken je Fehlerursache einen Code ab.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Tests | `play` kennt keine Option für Passwort und keine für das Überspringen der Zertifikatsprüfung | `make gates` |

## Re-Evaluierungs-Trigger

Wenn ein Anwender eine eigene Zertifizierungsstelle pro Aufruf angeben muss oder weitere Anmeldeverfahren gefordert werden.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-03 | Proposed | — |
| 2026-10-03 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
