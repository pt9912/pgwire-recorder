# Setzungsregister — Entscheidungsvorlage zu F-77

**Gegenstand:** Alle Setzungen S-1 bis S-46 aus den Review-Reports unter `docs/reviews/`. Eine Setzung ist eine Festlegung der Spezifikation oder der Planung, die das Lastenheft nicht trägt. Dieses Register ist ein Lauf-Beleg und eine Entscheidungsvorlage, keine Verifikation. Die Spalte „Entscheidung" füllt der Auftraggeber.

**Stand:** 2026-10-04, HEAD `35e7af3`.

## Lesehilfe

| Empfehlung | Bedeutung |
|---|---|
| **heben** | ins Lastenheft aufnehmen (Lastenheft-Commit vor den abhängigen Dokumenten; das Lastenheft ist `Draft`, ein Change Request entfällt) |
| **übernehmen** | bleibt Präzisierung in der Spezifikation, das Lastenheft engt nicht ein |
| **erledigt** | durch spätere Änderungen getragen |

## A — Überschreiten das Lastenheft: Entscheidung nötig

| Nr. | Setzung | Fundort | Empfehlung | Entscheidung |
|---|---|---|---|---|
| S-1 | Plattformen Linux, macOS, Windows je `amd64` und `arm64`; Image Linux `amd64` und `arm64` | SPEC-035, LH-FA-16.a, Handbuch, Releasing | **heben** in `LH-QA-03`: verpflichtet zu Windows-Build und -Test | |
| S-2 | Veröffentlichung in `ghcr.io/pt9912/pgwire-recorder` und `docker.io/pt9912/pgwire-recorder` | LH-FA-16.a, SPEC-031, Releasing | **heben** in `LH-FA-16`: das Abnahmeszenario 9 braucht einen Fundort | |
| S-8 | Windows-Konsolenabbruch entspricht `SIGINT` | LH-FA-13.a | **übernehmen**, sofern S-1 gehoben wird (Folge der Plattformwahl) | |
| S-9 | Reproduzierbarer Build für Binary und Image als Closure-Bedingung | Releasing, Welle, Slice | **heben** in `LH-QA-01` oder streichen; ohne Lastenheft-Anker ist es eine Prozessforderung | |
| S-11 | `PGR-E4002` und `PGR-E4003` beenden das Einspielen immer | LH-FA-20.a | **erledigt** (`e0ea6a5`): das Lastenheft sagt jetzt, dass nicht erreichbarer Server, fehlgeschlagene Anmeldung und abgebrochene Verbindung immer abbrechen; Weiterlaufen gilt nur für Fehlerantworten auf Anfragen | |
| S-12 | Passwort nur über `PGWIRE_RECORDER_PASSWORD` oder Platzhalter; `--user` und `--database` für alle Sessions | LH-FA-20.a, LH-FA-17.a | **heben** in `LH-FA-20`: „Zugangsdaten nicht als Option" | |
| S-13 | Zeitmodell: `offset_ms` ab Sessionbeginn, je Session eigene Uhr, kein Zeitbezug zwischen Sessions | LH-FA-21.a | **übernehmen** mit einem Satz im Lastenheft („Sessions laufen ohne gegenseitigen Zeitbezug") | |
| S-15 | Abnahmeszenario 11 (Homebrew) ist keine Bedingung von M3; Nachweis mit dem ersten stabilen Release | Roadmap, Welle, Releasing | **übernehmen** als Planungsentscheidung; das Lastenheft nennt Szenario 11 weiter als Abnahmekriterium von v1 | |
| S-20 | `PGWIRE_RECORDER_PASSWORD` ist Rückfall hinter dem Passwort der Verbindung | LH-FA-17.a, LH-FA-20.a | **übernehmen**, mit S-12 heben | |
| S-21 | `config show` zeigt den Dateiinhalt, die aktiven Umgebungsvariablen nur mit Namen | LH-FA-17.a | **übernehmen**: `LH-FA-17` sagt „Die gewählte Konfiguration lässt sich anzeigen, ohne Geheimnisse preiszugeben"; die frühere Lesart „wirksame Konfiguration" stand nicht im Lastenheft. Kein Widerspruch | |
| S-22 | Anmeldung am Server mit Klartext, MD5 und SCRAM-SHA-256; Zertifikat gegen den Systemspeicher, kein Überspringen | LH-FA-20.a | **heben** in `LH-FA-20` (Authentifizierung und Prüfung stehen dort schon allgemein); das Problem der selbst signierten Zertifikate ist durch `--upstream-ca` gelöst | |
| S-28 | TLS-Mindestversion 1.2 | LH-FA-23.a | **übernehmen**, in `LH-FA-23` ein Satz „zeitgemäße Verfahren" ohne Zahl | |
| S-35 | Interaktion ohne `ready_for_query` wird nicht verglichen | LH-FA-24.a | **erledigt**: das Lastenheft sagt jetzt „verglichen, soweit sie reicht" | |
| S-44 | Zusammenfassung am Ende des Laufs (Log-Zeile) | LH-FA-24.a | **heben** in `LH-FA-24` oder streichen; das Lastenheft verlangt sie nicht | |

## B — Teilweise vom Lastenheft getragen: übernehmen

| Nr. | Setzung | Fundort |
|---|---|---|
| S-3 | Homebrew über eigenen Tap, Formel mit SHA-256-Prüfung, kein Quellbau, nur stabile Releases | LH-FA-19.a, SPEC-042 |
| S-10 | Serverfehler beim Einspielen: Klasse „Netzwerk" (`PGR-E4004`, Exit-Code 4) | LH-FA-20.a, SPEC-017, SPEC-034 |
| S-16 | Replay ordnet mit der ersten Anfrage zu, Record nach Verbindungsannahme | LH-FA-12.a |
| S-18 | Konfigurationsdatei als YAML, genau eine Datei, Standardname im aktuellen Verzeichnis | LH-FA-17.a |
| S-19 | Klartext-Passwort in der Datei ist ein Konfigurationsfehler | LH-FA-17.a |
| S-23 | Ein Signal beendet das Einspielen mit Exit-Code 0, wenn kein Fehler auftrat | LH-FA-20.a |
| S-24 | Modus `absolute` holt Verspätung auf; Bezugspunkte `connect` und `first-request` | LH-FA-21.a |
| S-25 | Mit `--record-empty-sessions` zählt die Annahme der Verbindung, sonst die erste Anfrage | LH-FA-12.a |
| S-26 | SQLite: Formaterkennung am Dateikopf, Transaktion je Session, vier Tabellen | LH-FA-22.a, SPEC-043 |
| S-29 | Eigenes Zertifikat: Lesbarkeit, PEM, Zusammengehörigkeit, Ablauf; Namen nicht | LH-FA-23.a |
| S-30 | Abgewiesene Klartext-Verbindung ist `PGR-E6003` (Exit-Code 6); gescheiterter Handshake nur Warnung | LH-FA-23.a |
| S-38 | Verschlüsselter Schlüssel ist `PGR-E2007`; Prüfung beim Start | LH-FA-23.a |
| S-40 | Klartext-`CancelRequest` bei TLS bleibt `PGR-W3001` | LH-FA-23.a, LH-FA-05.e |
| S-42 | Ablauf eines Server- oder CA-Zertifikats ist `PGR-E4005` beim Verbinden | LH-FA-20.a |
| S-45 | Erwartetes Verbindungsende: gleicher SQLSTATE, Session endet, nächste läuft | LH-FA-20.a, LH-FA-24.a |

## C — Das Lastenheft schweigt, ohne zu widersprechen: übernehmen

S-4, S-5, S-6 (siehe F-70, erledigt), S-7, S-14, S-17 (mit [ADR-0015](../plan/adr/0015-uhr-port.md) getragen), S-27, S-33, S-34, S-37, S-39, S-46.

## D — Erledigt

S-31, S-32, S-36, S-41, S-43, S-35 (siehe A).

## Entscheidungsbedarf

Die Widersprüche S-11 und S-21 sind aufgelöst: S-11 durch eine Lastenheft-Änderung, S-21 war ein Lesefehler. Die Zeilen mit **heben** machen die Spezifikation tragfähig, ohne dass sich das Verhalten ändert. Nach der Entscheidung geht jede gehobene Setzung als Lastenheft-Commit voraus; die abhängigen Dokumente ziehen danach nach.

## Stand nach der Entscheidung (2026-10-04)

Der Auftraggeber hat die Empfehlungen der Gruppe A übernommen; S-9 abgeschwächt.

| Nr. | Ergebnis |
|---|---|
| S-1 | gehoben in `LH-QA-03` (`b725e50`) |
| S-2 | gehoben in `LH-FA-16` (`b725e50`) |
| S-8 | gehoben in `LH-FA-13` Boundary (`b725e50`) |
| S-9 | abgeschwächt gehoben in `LH-QA-01`: nur das Binary ist reproduzierbar (gleiche Prüfsumme); Welle, Releasing und Slices nachgezogen |
| S-11 | erledigt durch Lastenheft-Änderung (`e0ea6a5`) |
| S-12, S-22 | gehoben in `LH-FA-20` (`b725e50`) |
| S-13 | gehoben in `LH-FA-21` (`b725e50`) |
| S-15, S-20 | übernommen, ohne Lastenheft-Änderung |
| S-21 | übernommen, kein Widerspruch |
| S-28 | gehoben als „zeitgemäße TLS-Version" in `LH-FA-23`; die Mindestversion 1.2 bleibt Präzisierung der Spezifikation |
| S-44 | gehoben in `LH-FA-24` (`b725e50`) |
| Gruppen B und C | übernommen |
