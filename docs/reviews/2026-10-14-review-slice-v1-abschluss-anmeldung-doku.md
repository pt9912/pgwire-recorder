# Review slice-v1-abschluss-anmeldung-doku

Rolle: Reviewer. Datum: 2026-10-14.
Gegenstand: Diff `ea5019b..6a44529` (Move-Commit bis HEAD, fünf Commits) gegen den Plan [`slice-v1-abschluss-anmeldung-doku`](../plan/planning/done/welle-v1-abschluss/slice-v1-abschluss-anmeldung-doku.md), den Code-Slice [`slice-v1-abschluss-anmeldung`](../plan/planning/done/welle-v1-abschluss/slice-v1-abschluss-anmeldung.md), [LH-FA-05](../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), `LH-FA-05.b` (Randformen der Vermittlung, Record) in der [Spezifikation](../../spec/spezifikation.md) und `AGENTS.md` §3.7, §3.9, §3.11, §3.13, §3.14.

**Summary: 0 HIGH · 1 MEDIUM · 2 LOW · 4 INFO.**

## Findings

### F-1 MEDIUM
- `quelle`: `AGENTS.md` §3.9 (Der Slice-Plan folgt jeder Korrektur)
- `pfad`: `docs/plan/planning/done/slice-v1-abschluss-anmeldung-doku.md` · "**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle"
- `befund`: Commit `6a44529` ändert `spec/spezifikation.md` (`LH-FA-05.b`, Punkt *Aufzeichnung*: Satz zum Passwort als Startup-Parameter gestrichen). Der Kopf des Plans nennt keine Spec-Stelle, §1 sagt „Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan“, und §7 beschreibt die Spezifikation weiter mit dem gestrichenen Satz („Die Spezifikation sagt zum Passwort als Startup-Parameter, es stehe in der Aufzeichnung“; „Der Satz der Spezifikation dazu (Grenze) ist … nicht zu belegen“). Der Plan stammt aus den Commits vor `6a44529` und folgt dieser Korrektur nicht.
- `verifizierbar`: nein (`make kopf-check` prüft nur Kennungen aus §1 und §2)
- `klasse`: Plan folgt der Korrektur nicht (`BEO-REPO/plan-folgt-korrektur-nicht`)

### F-2 LOW
- `quelle`: `AGENTS.md` §3.11 und §3.7 (Zusage nur mit Prüfung; Kommentar beschreibt, was da ist)
- `pfad`: `spec/spezifikation.md` · "Ein Parameter, den der Server nicht kennt, lässt ihn"
- `befund`: Der neue Satz sagt ein Verhalten des Servers zu („lässt ihn den Aufbau ablehnen (Fehlerantwort, keine Session)“); „ihn“ ist ohne Bezug lesbar als Server oder als Recorder. Gedeckt ist es nur durch eine Probe des Implementers gegen PostgreSQL 17 in §7 des Plans, kein Test, kein Gate; „die Spezifikation sagt nichts darüber zu, ob ein Server …“ beschreibt zudem einen abwesenden Text.
- `verifizierbar`: nein
- `klasse`: Zusage im Text weiter als die Prüfung

### F-3 LOW
- `quelle`: `AGENTS.md` §3.11 (Version-Regel: Sätze über ein Gegenüber nennen die Version)
- `pfad`: `docs/user/benutzerhandbuch.md` · "Fordert die Datenbank ein anderes Verfahren an (etwa Kerberos, GSSAPI oder SSPI)"
- `befund`: Die Aussage zu Kerberos, GSSAPI und SSPI stammt von einem Stand-in-Server (Python), der nur die Nachricht `R` mit Code 2, 7 oder 9 sendet; §7 des Plans nennt das („kein echter Kerberos-Server“), das Handbuch nicht. Der Satz beschreibt, wie sich „die Datenbank“ verhält, und nennt weder Stand-in noch Version. Belegt ist die Reaktion des Werkzeugs auf diese Codes, nicht die Anforderung durch einen echten Server (zum Beispiel mit `gssencmode`); Sätze zu Verfahren und Verhalten der Datenbank ohne Version stehen außerdem an vier weiteren Stellen außerhalb des neuen Abschnitts (*Voraussetzungen* §1, *Hinweise* §4, §6 *Rollen und Rechte*, Zeile `PGR-E6001`), dort ohne Verweis auf „PostgreSQL 17“ (der Abschnitt selbst nennt ihn in *Voraussetzung*). Das ist das offene Risiko 3 aus §6 des Plans.
- `verifizierbar`: nein
- `klasse`: Belegte Umgebung nicht im Satz genannt

### F-4 INFO
- `quelle`: `AGENTS.md` §3.9
- `pfad`: `docs/reviews/2026-10-14-verifikation-slice-v1-abschluss-anmeldung.md` · "liegt in `next/`"
- `befund`: Commit `f149fda` (`make slice-mv`) hat im Verifikations-Report des Code-Slice nur den Link auf `in-progress/` gezogen; der Text daneben sagt weiter „liegt in `next/`“. Mechanik des Werkzeugs, kein Befund gegen den Implementer; der Report ist Lauf-Beleg.
- `klasse`: Move zieht Link, nicht Prosa nach

### F-5 INFO
- `quelle`: `AGENTS.md` §3.9
- `pfad`: Plan · "Fordert die Datenbank" nicht; DoD · "§7 Zeile `PGR-E6001` (Ursache „beim Aufzeichnen verlangt die Datenbank ein Passwort“)"
- `befund`: Die DoD zählt fünf ersetzte Stellen auf; der Diff ändert zusätzlich *Mehrere Verbindungen* (keine Sitzung), *Zeitangaben aufzeichnen*, `PGR-E4002` und `PGR-E4003`. §3 des Plans nennt alle; nur der Wortlaut der DoD ist enger. Kein Befund gegen den Diff.
- `klasse`: DoD-Wortlaut hinter dem Bestand

### F-6 INFO
- `quelle`: Maintainability
- `pfad`: `docs/user/benutzerhandbuch.md` · "Ebenso keine Sitzung ist eine Verbindung, deren"
- `befund`: Der eingefügte Satz steht mitten in einem Absatz, der danach mit „Das gilt auch für eine Verbindung, deren“ weitergeht, und die Zeile ist ohne Umbruch lang. Der Bezug von „Das gilt auch“ bleibt lesbar („keine Sitzung“). Ohne Konventions-Anker keine Aktion.
- `klasse`: Zeilenumbruch uneinheitlich

### F-7 INFO
- `quelle`: `LH-FA-05.b`
- `pfad`: `docs/user/benutzerhandbuch.md` · "Sendet die Datenbank in der Anmeldung eine Nachricht, die der Ablauf nicht vorsieht, ist"
- `befund`: Die Spezifikation trennt: eine nicht lesbare Nachricht des Servers im Aufbau ist `PGR-E6001`, eine lesbare außerhalb der Reihenfolge `PGR-E4002`. Das Handbuch nennt nur `PGR-E4002`; für den Fall „nicht lesbar“ und für „Anwendung sendet in der Anmeldung etwas anderes als Passwort oder `Terminate`“ (`PGR-E6001`) sagt es nichts. Das ist eine Lücke, keine falsche Aussage.
- `klasse`: Handbuch nennt nicht alle Randformen des Katalogs

## Eigene Prüfung (nicht auf Zusage)

**Aufbau:** `pgwire-recorder:dev` (Image 06:23, nach dem letzten Code-Commit `755d786` 06:22; seither nur Dokumente), `postgres:17-alpine` (17.11, `psql` 17.11) im eigenen Docker-Netz `rv-net`, Rollen `pw_scram`, `pw_md5`, `pw_klar` (Klartext in `pg_hba.conf`, `password`), `pw_reject`, `app`; ein Lauf zur Zeit, kein `pkill`. Container und Netz danach entfernt (`docker ps -a` und `docker network ls` ohne Treffer).

**Beispiele des neuen Abschnitts als Datei:** Beide Blöcke per `awk` aus dem Handbuch gezogen, ausgeführt mit `pgwire-recorder` als Funktion (Container), Host `postgres` → `rv-pg`, `./recordings/` → `/recordings/`, `localhost` → `rv-rec`. Ergebnis: Zeile `1`, Exit 0, Aufzeichnung trägt `user: app`, `database: app`, `sql: select 1`; `grep -c -E 'geheim|SCRAM-SHA|n,,n=|biws'` ergibt 0. Ein zweiter Lauf des ersten Blocks gegen dieselbe Datei endet mit `PGR-E2002` (Datei existiert); das Handbuch beschreibt diesen Fall an anderer Stelle, kein Befund.

**MD5, Klartext und Ablehnung:** `record` mit `pw_md5` (`md5`) und `pw_klar` (Klartext): beide liefern ihren Benutzer, zwei Sitzungen in der Datei, 0 Treffer auf `GEHEIM` und `md5[0-9a-f]{32}`. Falsches Passwort (`pw_md5`) und Regel `reject` (`pw_reject`): Text der Fehlerantwort wie direkt am Server (`password authentication failed …`, `pg_hba.conf rejects connection …`), im Log keine Zeile, Exit 0, keine Sitzung in der Datei.

**Fehlercodes mit Stand-in-Server (Python, nur zur Auslösung):** Code 7 → `PGR-E6001` „Anmeldeverfahren (Code 7) … nicht vermittelt“, Exit 6, `sessions: []` auch mit `--record-empty-sessions`. Code 11 als erste Nachricht → `PGR-E4002`, Exit 4. Server schließt nach Code 3 und Antwort → `PGR-E4003`, Exit 4. Server schließt direkt nach der Startnachricht → ebenfalls `PGR-E4003`, Exit 4. In allen vier Fällen erhält `psql` die Meldung als `FATAL`, `sessions: []`.

**Satz zu `PGR-E4003` gegen `LH-FA-05.b` (*Ende des Servers*) und Probe P6b:** Die Spezifikation: Senden einer Antwort scheitert oder die Verbindung des Servers endet in der Anmeldung → `PGR-E4003`; Senden der Startnachricht bleibt `PGR-E4002`. Das Handbuch („beendet sie die Verbindung in der Anmeldung, ist es `PGR-E4003`“, Tabellenzeile „oder beim Aufzeichnen hat die Datenbank sie in der Anmeldung beendet“) deckt sich mit Spezifikation und Läufen (zwei Auslöser, ein Code, Exit 4). Kein Befund.

**Channel Binding:** Ohne TLS: `channel_binding=require` → „channel binding required, but SSL not in use“, `prefer` → `1`. Mit `--tls-cert`/`--tls-key` (selbst signiert, EC P-256) und `sslmode=require`: `require` → „server did not offer an authentication method that supports channel binding“, `prefer` → `1`, `pw_md5` → `1`. Entspricht dem Handbuch.

**Zeitangaben:** Rohclient mit `--record-timing`: 1,5 s Wartezeit zwischen `R` (Code 3) und Passwort → `offset_ms: 1509`, ohne Wartezeit `8`. Entspricht dem geänderten Satz.

**Weitere Stichproben:** `record --help` nennt keine Option für Benutzer, Passwort oder Datenbank. `GSSENCRequest` an den Recorder → `N`, danach läuft die Anmeldung weiter (Info zu F-3). Katalog `internal/hexagon/model/fehler.go` führt `PGR-E4002`, `PGR-E4003`, `PGR-E6001`; alle drei stehen in der Tabelle des Handbuchs mit der neuen Ursache.

**Ist-Zustand und Verweise (`grep`):** `grep -n -i -E "slice|welle|adr-|LH-|SPEC-|ARC-|spec/|lastenheft|spezifikation|review|noch nicht|kommt mit|geplant|später|bisher|künftig"` über das Handbuch: Treffer nur Zeilen 14, 593, 675 (`später` als Beschreibung des Produkts bzw. der Reihenfolge), 167, 1277, 1659 (`noch nicht gültig`, `existiert noch nicht` als Zustand von Zertifikat und Datei); keiner im Diff, keiner ist Chronik, Zielstand oder Verweis. README: nur Verweise nach oben (`spec/…`, `make doc-trace`), im Diff keine Aussage über einen Stand. Die Reste von „ohne Passwort“ im Handbuch betreffen den neuen Abschnitt (Anwendung ohne Passwort), die Schlüsseldatei und die Support-Zeile; der alte Satz „Datenbank muss den Benutzer ohne Passwort anmelden“ steht nirgends mehr.

**Quantoren:** Neu: „weder in der Datei noch im Log, auch nicht mit `--log-level debug`“ (Datei und Log in der Probe 0 Treffer, Debug-Log trägt im Testbestand `TestE2ERecordAnmeldung`), „jedes Passwort“ bei `replay` (vom Implementer in P3 belegt, hier nicht neu gefahren), „unverändert“ (Text gleich; byte-gleich nur für Klartext belegt, die Spezifikation sagt „mit allen Feldern“). Im Bestand: die Sätze zu `PGWIRE_RECORDER_PASSWORD` und zur Konfigurationsdatei (§5) hält die Durchsicht des Plans für unverändert; die Stelle „Keine Sitzung ist, auch mit der Option …“ ist ergänzt. Eine Lücke habe ich nicht gefunden.

**Version-Regel:** Der neue Abschnitt nennt PostgreSQL 17 und libpq 17, nennt Windows, macOS und andere Treiber nicht und sagt „andere Versionen … nicht geprüft“; README: „geprüft mit PostgreSQL 17“. Rest siehe F-3.

**Versionszeile:** `Version: 0.8` (vorher 0.7), `Stand: 14.10.2026`; es gibt keinen Änderungsverlauf im Handbuch, der nachzuziehen wäre.

**Plan-Nachzug §3.9:** §3 des Plans nennt jede im Diff geänderte Handbuch-Stelle (neuer Abschnitt, §1, §4 mit zwei Hinweisen und Zeitangaben, §6, §7 mit vier Zeilen, Versionszeile) und README. Offen: die Spezifikationsänderung (F-1).

**Größe §3.14:** `git diff --numstat ea5019b..HEAD -- docs/user README.md`: Handbuch 81 + 15 = 96, README 4 + 3 = 7, Summe 103; Marke 300 und Meldeschwelle 220 nicht erreicht, Zahl in §7 stimmt. Der Plan ist nicht in der Metrik.

**Adressen (§3.13):** Der Diff weist keinem anderen Slice etwas zu; `slice-v1-abschluss-postgres-versionen` ist im Risiko 3 genannt (Empfehlung des Planners, schon in `done/`-Plan des Code-Slice) und ändert sich in diesem Diff nicht.

## Geprüft, ohne Befund

- `docs/user/benutzerhandbuch.md`, neuer Abschnitt *Eine Anwendung mit Passwort-Anmeldung aufzeichnen*: Beispiele und Verhaltensaussagen gegen das Binary, siehe oben.
- `docs/user/benutzerhandbuch.md`, Fehlercode-Tabelle und *Die Anwendung kann sich nicht verbinden*: Codes, Exit-Codes und Ursachen gegen Katalog und Läufe.
- `README.md`: Satz zu `record` und `play`; nur Verweise nach oben.
- `internal`, `cmd`, `test`, `harness`, `.claude`, `.harness`, `Makefile`: im Diff nicht berührt.
- Move-Commit `ea5019b`: rein; Inhaltsänderungen stehen in eigenen Commits (`AGENTS.md` §3.3).
- Suppression, Core-Reinheit, ADRs: nicht berührt.
