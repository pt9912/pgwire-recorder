# Review-Report: slice-v1-abschluss-einspielen-anmeldung-doku — 2026-10-10

**Review-Art:** Code-Review (Dokumentation). Geprüft wird gegen Plan §1, §3, §6 und §7, gegen `AGENTS.md` §3.9, §3.11 (inkl. *Handbuch und README beschreiben den Ist-Zustand*), §3.12, §3.13 und gegen `LH-FA-20`. Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** Diff `e8fd0b2..1898740`: Randformen und Probe-Matrix `36480b8` (Architect), Handbuch, README und §3 `e70ebac`, Belege `8cc0a48`, Gate-Zeile `1898740` (Implementer). Geändert sind `docs/user/benutzerhandbuch.md`, `README.md` und der Slice-Plan.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-v1-abschluss-einspielen-anmeldung-doku` (§1, §3, §6 mit den Randformen des Architect, §7 *Belege des Implementers*)
- `slice-doku-ist-stand` §6 und §7 (Randformen und Befehlsfolge für Doku-Slices)
- `LH-FA-20`, `LH-FA-17`, `LH-QA-05`, `LH-RB-01`
- `AGENTS.md` (Hard Rules)

**Selbst gefahren** (gebautes Binary `sha256:c2871ccf3ceb`, gleich dem Image im Beleg des Implementers; gepinntes PostgreSQL 17 aus `harness/mk/integration.mk` mit `pg_hba.conf` je Benutzer auf `scram-sha-256`, `md5`, `password`, dazu ein SCRAM-Benutzer mit U+00A0 im Passwort und die Standardzeile `trust`; Aufzeichnung von `SELECT 1` mit `record` und `psql`; Scratch-Präfix `rev-anmdoku-`, Container, Netz, Image und Dateien danach entfernt):

- Passwort je Verfahren richtig (SCRAM mit Leerzeichen, `$`, `%`; MD5; Klartext) meldet an, falsch endet mit `PGR-E4005`, SQLSTATE `28P01`, Exit 4; unbekannter Benutzer `28000`; die Meldung nennt `Session 1`.
- Variable leer oder ungesetzt gegen SCRAM: `PGR-E4005` „der Server verlangt ein Passwort, play hat keines“. `trust` mit gesetzter Variable gelingt.
- Vorrang: Verbindung mit Passwort-Platzhalter und falscher Variable meldet an; Platzhalter falsch und Variable richtig endet mit `PGR-E4005`; Verbindung ohne Passwortteil nimmt die Variable; `host:port` nimmt die Variable. Platzhalter-Variable leer oder ungesetzt: `PGR-E2005`, Exit 2; `user:@host` und `?password=`: `PGR-E2006`.
- `config show` nennt bei gesetzter Variable am Ende `PGWIRE_RECORDER_PASSWORD` ohne Wert, bei leerer Variable nichts; Beispiel §5 als Datei (Zeilen 572 bis 583) unverändert ausgegeben; mit der Zeile `test` auf den Server gerichtet meldet `--upstream test` ohne `--user` als Benutzer der URL an (Log-Stufe `debug`, `GEHEIM` in keiner Ausgabe); `sslmode=require` bei `play`: `PGR-E2004`.
- SASLprep-Grenze: U+00A0 im Passwort endet gegen `scram-sha-256` mit `PGR-E4005`, `28P01`.
- `record` gegen einen SCRAM-Benutzer: Client-Fehler und Log mit `PGR-E6001`, Exit 6.
- Befehlsfolge aus `slice-doku-ist-stand` §7 am Handbuch: Optionen nur `--h`, `--help`, `--version`, `--name`, `--rm`; Codes ohne Ausgabe; Variablen genau `PGWIRE_RECORDER_PASSWORD`.
- `grep` über Handbuch und README auf Spezifikation, Lastenheft, ADR, slice, welle, review und Zielstand-Wörter: nur die Verweise in `README.md` §Auditierbarkeit (Links ohne Standaussage) und Alltagswörter („später“, „aktuellen“, „noch nicht“ bei einer Datei).
- Namensabgleich §7: `TestAnmeldungNichtUnterstuetzt` (Fälle `nur PLUS`, `Kerberos`, `GSS`, `SSPI`, `SCM`), `einspielen_anmeldung_test.go`, `envPassword` in `cli.go` und `TestE2EPlayAnmeldungFehler` (prüft `PGR-E4005`, Exit 4 und `kein`) existieren. `TestAnmeldungNichtUnterstuetzt` prüft den Code `PGR-E4005` und dass nichts gesendet wird, nicht den Meldungstext; das Handbuch sagt zu dieser Gruppe nur den Code.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-597 | MEDIUM | §3 *Ansatz* sagt, die Befehlsfolge aus `slice-doku-ist-stand` §7 laufe vor der Übergabe grün. Sie meldet am Handbuch aber bei jedem Lauf genau `PGWIRE_RECORDER_PASSWORD` (Variable ohne Option); die „bekannte Ausnahme“ steht nur in den Belegen von §7, nicht in §3, und die Mutante „Zeile `PGWIRE_RECORDER_PASSWORD` zurück“ der Befehlsfolge ist nicht mehr trennscharf. Jeder spätere Lauf der Folge (Verifikation, Closure, weitere Doku-Slices) liest die eine Ausgabezeile als Rot oder als Normalfall und kann sie nicht von einer echten verwaisten Variable unterscheiden. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-anmeldung-doku.md` · „Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe grün.“ | ja — die Befehlsfolge am Handbuch (Variablen gegen Optionen) | Plan folgt der Korrektur nicht (Befehlsfolge) |
| F-598 | LOW | Der Satz in §6 *Rollen und Rechte* ist inhaltlich richtig (für `record` und `replay` wie vorher, für `play` mit Benutzer, Datenbank und Passwort; belegt durch die Proben), steht aber nur in der §3-Zeile des Implementers und nicht unter den Randformen in §6 des Plans, die der Architect vor dem Code entscheidet. Die Randform-Zeile „`PGR-E4005` ohne Passwort-Aussage“ führt als Orte nur §1, §4, §5 und §7; der Implementer meldet die Lücke selbst als Fund 2. | `AGENTS.md` §3.12 | `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-anmeldung-doku.md` · „§6 *Rollen und Rechte* (Nachbarsatz zur Anmeldung, der für `play` nicht mehr gilt)“ | nein — Urteil | Nachbarsatz-Randform erst im Code-Commit genannt |
| F-599 | INFO | §1 *Voraussetzungen* und §4 *Voraussetzung* verweisen für das Passwort auf den Abschnitt *Konfigurationsdatei*; die Variable `PGWIRE_RECORDER_PASSWORD` und ihr Vorrang stehen im Absatz hinter der Optionstabelle von §5, vor diesem Anker. §4 nennt die Variable selbst, §1 nicht. Kein Satz ist falsch. Für die zuständige Rolle: Doku-Autor. | Maintainability | `docs/user/benutzerhandbuch.md` · „nennen Sie es `play` (siehe Konfigurationsdatei)“ (Link auf den Abschnitt *Konfigurationsdatei*) | nein | Verweisziel enthält die Zusage nicht |
| F-600 | INFO | Die SASLprep-Grenze („kann bei SCRAM scheitern“) hat als Beleg nur die Probe in §7 und die Grenze der Spezifikation, keinen Test oder Gate; der Satz ist als Grenze mit „kann“ gefasst und sagt nicht mehr zu, als die Probe zeigt. Gebrochen wird `AGENTS.md` §3.11 damit nicht. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „ein Passwort mit\nZeichen, die SASLprep ändert, kann bei SCRAM scheitern“ | nein | Grenze nur durch Probe belegt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Handbuch ohne Chronik, Zielstand, Verweise auf Spezifikation, Lastenheft, ADR, Slice, Welle, Review | geprüft, ohne Befund |
| `README.md` ohne Stand-Aussagen und Chronik (Satz zu `play` und `record`; die Links in §Auditierbarkeit nennen keinen Stand) | geprüft, ohne Befund |
| Handbuch §1 und §4 *Voraussetzung*, *Hinweise* (Passwortquellen, Verfahren, `PGR-E4005`, Klartext-Satz) gegen Binary und PostgreSQL 17 | geprüft, ohne Befund |
| Handbuch §5 Absatz *Passwort*, Vorrang, leere Variable, `trust`, `config show`, Ausgaben ohne `GEHEIM`, SASLprep | geprüft, ohne Befund |
| Handbuch §5 Beispiel `test` als Datei, `config show`, `play --upstream test`, Platzhalter-Absatz (`PGR-E2005`, `PGR-E2006`) | geprüft, ohne Befund |
| Handbuch §4 Beispiel `play` | geprüft, ohne Befund |
| Handbuch §6 *Rollen und Rechte* für `record`, `replay`, `play` | geprüft, ohne Befund (Plan-Lücke siehe F-598) |
| Handbuch §7 Zeilen `PGR-E4005` (Verfahren, die `play` nicht kann: `nur PLUS`, Kerberos, GSS, SSPI durch Test), `PGR-E6001`, `PGR-E2005`, `PGR-E2006` und alle weiteren „ohne Passwort“-Sätze (jeder betrifft `record`) | geprüft, ohne Befund |
| Plan §1, §3, §6 gegen den Diff (§3.9), Adressen (§3.13: keine Zuweisung im Diff, §6 *Zuweisungen* sagt „keine“) | geprüft, ohne Befund außer F-597 |
| §7 Belege: genannte Tests, Dateien und Symbole existieren | geprüft, ohne Befund |
| Mutations- und Randform-Pflicht (§3.10, §3.12): Slice liefert keinen Vertrag, kein Code im Diff | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Plan folgt der Korrektur nicht (Befehlsfolge) · Nachbarsatz-Randform erst im Code-Commit genannt · Verweisziel enthält die Zusage nicht · Grenze nur durch Probe belegt

## Verdikt

**Merge-blockierend:** ja — wegen F-597 (MEDIUM): §3 verlangt einen grünen Lauf der Befehlsfolge, den es mit der benannten Ausnahme nicht gibt. F-598 bis F-600 blockieren nicht.

**Übergabe:** F-597 und F-598 gehen an den Implementer, F-597 mit Widerspruchsrecht an den Planner (die Befehlsfolge selbst steht in `done/slice-doku-ist-stand.md` §7, der Implementer hat sie als Fund 1 gemeldet). Ein Widerspruch senkt die Einstufung nicht; ab HIGH mit Rollen-Widerspruch läuft der Konflikt-Pfad über den Architect. Die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg; er ersetzt keine Verifikation.
