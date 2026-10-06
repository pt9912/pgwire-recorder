# Slice slice-replay-semantik-meldungscodes: Meldungscodes, Fehlertext und Logging

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-replay-semantik.

**Bezug:** [`LH-FA-14`](../../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md)

**Berührte Spec-Stellen:** `LH-FA-14.a` · `SPEC-034` · `SPEC-005` · `SPEC-006` · `SPEC-013` bis `SPEC-028` · `LH-FA-13.b` · `LH-FA-17.a` · `LH-FA-01.a` · `LH-FA-10.a` · `LH-FA-18.a` · `SPEC-033` · `LH-FA-11.a` · `ARC-001` · `ARC-005` · `ARC-006` · `ARC-009`

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Jeder Fehler und jede Warnung trägt einen `PGR-…`-Meldungscode; der Fehlertext beginnt mit `<klasse> [<code>]: `, Logs gehen nach `stderr` mit einstellbarem Level (`--log-level`). Die Randformen dieser Verträge sind vor dem Code entschieden (§6, Architect 2026-10-06) und stehen in `LH-FA-14.a` und `SPEC-034` §Ausgabe.

**Schon geliefert, mit Test** (Bestand der Wellen davor; dieser Slice baut darauf, statt ihn neu zu liefern):

- Code-Tabelle im Quelltext (`internal/hexagon/model/fehler.go`, nur die Codes, die der Code erzeugt), Kopf `<klasse> [<code>]: ` in `model.Error` und Exit-Code aus der ersten Ziffer des Codes — geprüft für `PGR-E2001` mit Exit-Code 2 (`cli_test.go`), `PGR-E2001`/`PGR-E3001` als Startfehler auf `stderr` (`bootstrap_test.go`), `PGR-E3003` mit Exit-Code 3 und `PGR-E5001` mit Exit-Code 5 (`replay_e2e_test.go`), den Kopf `Replay [PGR-E5002]: ` (`replay_unverbraucht_test.go`).
- `ErrorResponse` mit `FATAL`, dem Code im Meldungstext und SQLSTATE `0A000` beziehungsweise `08006` (`server_test.go`, `server_extended_test.go`).
- Warnungen `PGR-W2001`, `PGR-W3001`, `PGR-W3003` mit Attribut `code` (`server_test.go`, `record_e2e_test.go`, `unverbraucht_e2e_test.go`).
- Hilfe und `version` auf `stdout`, nichts auf `stderr` (`bootstrap_test.go`).
- Parameterwerte erscheinen nicht in der Diagnose einer Extended-Abweichung; sie nennt nur das Feld `params` (`replay_extended_test.go`).

**Neu in diesem Slice:**

- Fehlertext einzeilig, Kopf genau einmal mit dem Code des äußersten Fehlers (heute trägt eine Kette zwei Köpfe, etwa `PGR-E3003` um `PGR-E3003` beim Lesen eines Werts), fremde Fehler mit Kopf (heute ohne: der Annahmefehler im PGWire-Adapter und jeder nicht eingeordnete Fehler in `note` und `fail`), Klassenname nach `SPEC-034` (`SPEC-034` §Ausgabe).
- `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` bei `record` und `replay` mit den vier Stufen, der Strenge des Werts und dem Inhalt der Stufen; heute ist die Stufe fest `info` (`LH-FA-14.a`).
- Test je Fehlerklasse (Exit 1 bis 6) für Kopf, Klassenname und Exit-Code, dazu ein Test, dass jeder Code der Tabelle im Quelltext die Form aus `SPEC-034` hat und seine Klasse ergibt.
- Gleichrangige Fehler (`errors.Join`, etwa aus `CloseSession`) als je eine Meldung mit eigenem Kopf, die erste gemerkt, nicht klassifizierte als Ursache der ersten (`SPEC-034` §Ausgabe *Gleichrangige Fehler*); die Zeile beim Prozessende folgt derselben Regel. Eine nicht annehmbare Verbindung ist der gemerkte Verbindungsfehler `PGR-E4000` (`LH-FA-13.b`); die Debug-Zeile einer Verbindung ohne Startnachricht nennt den Bibliothekstext unter `grund` (`LH-FA-14.a` §Zeilenform).

**Übernommen aus `slice-extended-query-replay`** (Review F-328, Validierung `docs/reviews/2026-10-05-validierung-slice-extended-query-replay.md`, Befund 3): Die Regel zu Parameterwerten stand in `LH-FA-18.a` §Mismatch mit der Bedingung „wenn der Log-Level nicht `debug` ist“, in `SPEC-033` ohne Bedingung; die Diagnose nannte den Index des abweichenden Parameters nicht. **Vom Nutzer entschieden am 2026-10-06** und in die Spezifikation geschrieben: Parameterwerte erscheinen nie in der Diagnose, weder im Log noch in der `ErrorResponse`, auf keinem Log-Level (`LH-FA-18.a` §Mismatch an `SPEC-033` angeglichen); die Diagnose nennt die Nummer des ersten abweichenden Parameters ohne Wert, bei abweichender Zahl beide Anzahlen (`LH-FA-18.a` §Mismatch). Die Umsetzung liefert DoD-Punkt 3.

**Übernommen aus `slice-replay-semantik-fehlerreplay`** (Verifikation `docs/reviews/2026-10-06-verifikation-slice-replay-semantik-fehlerreplay.md`, V-52): Abnahmeszenario 4 ist für Simple Query nur zur Hälfte nachgewiesen. `TestE2EReplayAbweichung` prüft `PGR-E5001` und Exit-Code 5, deklariert auch „keine Antwort“, verwirft aber die Ergebnisse von `ReadAll`; `TestReplayMismatch` deklariert `LH-FA-10/Negative` und prüft nur Code und Cursor. Die Mutation „bei abweichender Anfrage liefert das Replay die Antworten der aufgezeichneten Anfrage, der Driving-Adapter sendet sie vor `PGR-E5001`“ bleibt grün (VF24). Der Code ist richtig, es fehlt nur der Nachweis für `LH-FA-10` Negative („statt eine unpassende aufgezeichnete Antwort zu verwenden“). Ohne ihn schließt welle-replay-semantik nicht (Closure-Trigger, Abnahmeszenario 4). Der Nachweis gehört in DoD-Punkt 3, weil er dieselbe Anforderung und dieselbe Diagnose einer Abweichung betrifft; er braucht keinen Produktionscode.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Gate, das die Code-Tabelle im Quelltext mit dem Katalog des Handbuchs und mit `SPEC-034` abgleicht — Folge-Slice der welle-v1-abschluss (dort unter *Ausdrücklich nicht* geführt), nach dem Katalog in der Betriebsdokumentation (`slice-v1-abschluss-container`); hier prüft ein Test nur Form und Klasse der Codes im Quelltext.
- Codes, die erst ein späterer Auslöser erzeugt (`PGR-E2003` bis `PGR-E2007`, `PGR-E4004` bis `PGR-E4006`, `PGR-E5004`, `PGR-E6003`, `PGR-W3002`), und deren Warnungen — kommen mit dem Slice, der den Auslöser liefert, in die Tabelle; sie folgen `SPEC-034`.
- Schlüssel `log_level` der Konfigurationsdatei — `slice-v1-abschluss-betrieb` (Konfigurationsdatei für alle Optionen); hier nur Option und Umgebungsvariable.
- `--log-level` bei `play` — `slice-v1-abschluss-einspielen` (Optionen von `play`); `play` gibt es noch nicht.
- Info-Zeile beim Beginn des Herunterfahrens und Exit-Code nach dem Herunterfahren je Klasse — `slice-v1-abschluss-betrieb`.
- Allgemeiner Leser für Umgebungsvariablen — `slice-v1-abschluss-betrieb`; `--log-level` liest seine Variable wie `--fail-on-unconsumed` (Risiko in §6).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-14`](../../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Je Fehlerklasse (`SPEC-013` bis `SPEC-019`, Exit 1 bis 6) beginnt der Fehlertext mit dem Kopf aus `SPEC-034` §Ausgabe (Klassenname, Code), als Zeile beim Prozessende, als Attribut `error` und als Meldungstext der `ErrorResponse` mit SQLSTATE nach Klasse, und der Exit-Code ist der der Klasse; der Fehlertext ist eine Zeile, eine Fehlerkette trägt einen Kopf mit dem Code des äußersten Fehlers, ein nicht eingeordneter Fehler ist `PGR-E1000` mit Kopf; jeder Code der Tabelle im Quelltext hat die Form aus `SPEC-034` und ergibt seine Klasse (Test je Klasse und je Regel; die Abbildung beim Herunterfahren prüft `slice-v1-abschluss-betrieb`).
- [ ] [`LH-FA-14`](../../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): `--log-level` und `PGWIRE_RECORDER_LOG_LEVEL` wirken bei `record` und `replay` wie `LH-FA-14.a` sagt: vier Stufen, jede zeigt nur ihre und die strengeren Zeilen, Standard `info`; jeder andere Wert, auch leer und großgeschrieben, ist `PGR-E2001`, eine ungültige Umgebungsvariable auch neben einer gültigen Option; eine Log-Zeile trägt `level`, Fehler `code` und `error`, Warnungen `code`; ein Startfehler erscheint auf jeder Stufe als Fehlertext (Test je Stufe und je Randform aus §6).
- [ ] [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol): Die Diagnose einer Parameter-Abweichung nennt die Nummer des ersten abweichenden Parameters (ab 1), bei abweichender Zahl beide Anzahlen, ohne einen Wert; kein Parameterwert erscheint in Log oder `ErrorResponse`, auch bei `--log-level debug` (`LH-FA-18.a` §Mismatch, `SPEC-033`; Test). Bei einer abweichenden einfachen Anfrage erhält der Client keine aufgezeichnete Antwort, nur `PGR-E5001` (Abnahmeszenario 4, `LH-FA-10` Negative, übernommen aus `slice-replay-semantik-fehlerreplay`): `TestE2EReplayAbweichung` verlangt, dass `ReadAll` kein Ergebnis liefert, `TestReplayMismatch`, dass `Query` bei Abweichung keine Antworten liefert; beide werden rot, wenn das Replay die Antworten der aufgezeichneten Anfrage vor `PGR-E5001` liefert (E2E und Unit, Abdeckung `LH-FA-10/Negative`).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/model/fehler.go` | update | Code-Tabelle bleibt hier (kein eigenes Paket); `Error()` einzeilig, ein Kopf je Kette; `Meldungen` zerlegt nur eine reine Zusammenfassung (`errors.Join`, erkannt daran, dass ihr Text genau der ihrer Teile ist) in gleichrangige Meldungen, eine Hülle mit eigenem Text um mehrere Ursachen bleibt eine Kette, fremde werden `PGR-E1000` (`SPEC-034` §Ausgabe) |
| `internal/hexagon/model/fehler_test.go` | neu | DoD 1: Kopf, Klassenname und Exit-Code je Klasse, Kette samt Hülle um mehrere Ursachen (S1 bis S3), fremde Fehler, Zeilenumbruch, gleichrangige Fehler, Form jedes Codes der Tabelle (aus dem Quelltext gelesen) |
| `internal/hexagon/services/replay.go`, `matcher.go` | update | DoD 3: Nummer des ersten abweichenden Parameters, beide Anzahlen |
| `internal/hexagon/services/replay_test.go`, `replay_extended_test.go` | update | DoD 3: `TestReplayMismatch` prüft, dass `Query` bei Abweichung keine Antworten liefert; Nummer und Anzahlen ohne Wert |
| `internal/adapters/driving/cli` | update | DoD 2: `--log-level` und Umgebungsvariable bei `record` und `replay`, Hilfetext |
| `internal/bootstrap` | update | DoD 1 und 2: Stufe des Loggers aus der Option; Zeile beim Prozessende mit Kopf auch für einen nicht eingeordneten Fehler |
| `internal/adapters/driving/pgwire/server.go` | update | DoD 1: Annahmefehler über `note` als Verbindungsfehler `PGR-E4000` (gemerkt), nicht eingeordneter Fehler in `note` mit Kopf `PGR-E1000`; `note` zerlegt einen `errors.Join` in eigene Meldungen, merkt die erste, hängt nicht klassifizierte als Ursache an; `fail` stellt genau eine `ErrorResponse` mit der ersten Meldung zu; Debug-Zeile ohne Startnachricht mit `grund` statt `error`. `record.go` (`CloseSession`) bleibt unverändert: der `errors.Join` ordnet Verbindungsende, Schreiben des Recordings, Schließen zum Upstream; das Schließen entsteht vor dem Schreiben, steht aber dahinter, und das ist ohne Folge, solange `Upstream.Close` nur nicht klassifizierte Fehler liefert, die keine eigene Meldung werden |
| Tests in `cli`, `bootstrap`, `pgwire` (dort neu `server_meldung_test.go`) | update | DoD 1 und 2: je Klasse (`ErrorResponse`, Zeile beim Prozessende), gleichrangige Fehler, Annahmefehler, `grund`, je Stufe bei `record` und `replay`, Zeilenform, Startfehler je Stufe bei `record` und `replay`, Schreibfehler am Ende je Stufe bei `record`, `version` ohne Prüfung der Umgebungsvariable, Werte und Umgebungsvariable, Hilfe |
| `test/integration/replay_e2e_test.go` | update | DoD 3, Abnahmeszenario 4: `TestE2EReplayAbweichung` prüft die Ergebnisse von `ReadAll` (keine), Deklaration `LH-FA-10/Negative` ergänzt; `TestE2EReplayNachDemEnde`: eine Anfrage nach dem Ende der Aufzeichnung erhält `PGR-E5001` und kein Ergebnis (`LH-FA-10.a`) |
| `test/integration/extended_replay_e2e_test.go` | update | DoD 3: `TestE2EReplayExtendedAbweichung` läuft mit `--log-level debug`; `ErrorResponse` und Log nennen `(Parameter $1)`, keinen Parameterwert |
| `docs/user/benutzerhandbuch.md` | update | Werte von `--log-level` und ihre Strenge, Zeilenform, Nummer des abweichenden Parameters |
| `docs/user/abdeckung-*.md` | update | von `make abdeckung` aus den geänderten Deklarationen geschrieben |

Berührte Schichten: Core (`model`, `services`) und Driving-Adapter (`cli`, `pgwire`, mit `bootstrap` als Verdrahtung) — zwei.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-replay-semantik-mismatch` ist `done`, und die Fassung von `LH-FA-18.a` §Mismatch und `SPEC-033` zu Parameterwerten und Parameter-Index ist vom Nutzer bestätigt (§1, 2026-10-06); die Randformen aus §6 sind vor dem ersten Code-Commit entschieden.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Der Diff passt nicht in eine Review-Sitzung, oder `--log-level` verlangt den allgemeinen Optionsleser aus `slice-v1-abschluss-betrieb` — zurück zur Zerlegung; DoD 2 geht dann als eigener Slice der Welle ab (Vorschlag: `slice-replay-semantik-log-level`).
- `in-progress` → `open`: Eine Klasse lässt sich nicht eindeutig zuordnen — Carveout.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen der neuen Verträge** (`AGENTS.md` §3.12), je mit dem Ort ihrer Entscheidung; geprüft und entschieden vom Architect am 2026-10-06, die Nutzerentscheidungen vom selben Tag eingeschlossen:

| Randform | Entscheidung | Ort |
|---|---|---|
| Mehrzeilige Ursache, Kontinuationszeilen | Fehlertext ist eine Zeile; ein Zeilenumbruch wird mit dem Leerraum danach durch ein Leerzeichen ersetzt, auch im Text einer Bibliothek; keine Kontinuationszeilen | `SPEC-034` §Ausgabe *Eine Zeile* |
| Fehler vor dem Lesen von `--log-level`, ungültiger Wert von `--log-level` | Startfehler erscheint auf jeder Stufe als Zeile beim Prozessende; vorher schreibt der Prozess nichts anderes | `LH-FA-14.a` §Log-Level |
| Werte von `--log-level`, Strenge | genau `error`, `warn`, `info`, `debug`; kleingeschrieben; leer, großgeschrieben, anderer Name (`warning`, `trace`, `off`) sind `PGR-E2001`; Umgebungsvariable: leer gilt als nicht gesetzt, ungültig ist `PGR-E2001` auch neben gültiger Option; Mehrfachangabe: letzte gilt; nur nach dem Kommando | `LH-FA-14.a` §Log-Level, `LH-FA-17.a` |
| Stufen und ihr Inhalt | `error` Verbindungsfehler, `warn` Warnungen (jede mit Code), `info` Start, Ende und Beginn des Herunterfahrens, keine Zeile je Verbindung, `debug` je Verbindung, nicht Vertrag; Schwelle zeigt die strengeren Stufen mit | `LH-FA-14.a` §Log-Level |
| Code für Warnungen | jede Zeile der Stufe `warn` trägt `PGR-W…` als Attribut `code`; ein Hinweis ohne Maßnahme geht nach `info` oder `debug` | `LH-FA-14.a`, `SPEC-034` §Ausgabe |
| `stdout` gegen `stderr`, auch bei Hilfe und `version` | Logs und Fehlertext auf `stderr`; auf `stdout` nur Hilfe, `version`, `config show`, die nichts auf `stderr` schreiben; `--log-level` wirkt auf keines davon | `LH-FA-14.a`, `LH-FA-01.a` |
| Zeitstempel und Zeilenform | `logfmt`: `time` (RFC 3339, Millisekunden, Zonenversatz), `level` (`DEBUG` bis `ERROR`), `msg`, dann Attribute; Vertrag nur `level`, `code`, `error`; Zeile beim Prozessende ist genau der Fehlertext | `LH-FA-14.a` §Zeilenform |
| `ErrorResponse` gegen Log | derselbe Fehlertext wie `error`, unabhängig vom Log-Level; `FATAL`; SQLSTATE `0A000`, `08006`, sonst `XX000`; keine weiteren Felder | `SPEC-034` §Ausgabe *Zustellung an den Client* |
| Interne Fehler ohne Code | `PGR-E1000` mit Kopf, Exit-Code 1; ein Fehlertext ohne Kopf kommt nicht vor | `SPEC-034` §Fehler und §Ausgabe *Fremde Fehler* |
| Fehler aus Bibliotheken | Code der Stelle, die ihn einordnet, Bibliothekstext als Ursache (nicht Vertrag); eine weitergeleitete oder wiedergegebene `ErrorResponse` des Servers trägt keinen Code | `SPEC-034` §Ausgabe *Fremde Fehler*, `LH-FA-11.a` |
| Fehlerkette mit mehreren Codes | ein Kopf, Code und Klasse des äußersten klassifizierten Fehlers | `SPEC-034` §Ausgabe *Fehlerkette* |
| Parameterwerte, auch bei `debug` | nie in Log oder `ErrorResponse` (Nutzer, 2026-10-06) | `LH-FA-18.a` §Mismatch, `SPEC-033` |
| Parameter-Index, abweichende Anzahl | Nummer des ersten abweichenden Parameters ab 1 wie `$1`; bei abweichender Zahl beide Anzahlen statt einer Nummer (Nutzer, 2026-10-06) | `LH-FA-18.a` §Mismatch |
| Gleichrangige Fehler (`errors.Join` in `CloseSession`), Fehler ohne Klasse daneben, Annahmefehler (Rückgabe des Implementers 1) | je klassifiziertem Fehler eine Meldung mit eigenem Kopf in Reihenfolge des Entstehens, der erste gemerkt; nicht klassifizierter wird Ursache der ersten, nur nicht klassifizierte ergeben eine `PGR-E1000`; Vorrang des Schreibfehlers gilt nur beim Prozessende; nicht annehmbare Verbindung ist Verbindungsfehler `PGR-E4000` und zählt für den Exit-Code | `SPEC-034` §Ausgabe *Gleichrangige Fehler*, `LH-FA-13.b` |
| Mehrere Meldungen in `fail` (F-401) | genau eine `ErrorResponse` mit Text und SQLSTATE der ersten, gemerkten Meldung; die übrigen nur im Log | `SPEC-034` §Ausgabe *Zustellung an den Client* |
| Hülle mit eigenem Text um mehrere Ursachen, etwa mehrere `%w` (F-402) | Kette: eine Meldung mit ganzem Text, Code und Klasse des ersten klassifizierten Fehlers darunter (Reihenfolge der Ursachen, außen nach innen), kein innerer Kopf, sonst `PGR-E1000`; gleichrangig ist nur eine reine Zusammenfassung (`errors.Join`) | `SPEC-034` §Ausgabe *Fehlerkette*, *Gleichrangige Fehler* |
| Suchreihenfolge in einer Hülle mit mehreren Ursachen verschiedener Tiefe (V-57) | Tiefensuche wie `errors.As`: jede Ursache ganz, bis innen, vor der nächsten; ein tiefer klassifizierter Fehler in der ersten Ursache geht einem flachen in der zweiten vor | `SPEC-034` §Ausgabe *Fehlerkette* |
| Was ein Zeilenumbruch ist (Rückgabe 2) | LF, CR LF, einzelnes CR, mit folgenden Leerzeichen, Tabs und Umbrüchen; VT, FF, U+0085, U+2028, U+2029 nicht | `SPEC-034` §Ausgabe *Eine Zeile* |
| Attribut `error` an einer Debug-Zeile (Rückgabe 3) | `error` nur an `error`-Zeilen und immer mit Kopf; Bibliothekstext an anderen Stufen unter `grund` | `LH-FA-14.a` §Zeilenform |
| Abgleich Code-Tabelle mit Code und Katalog | hier ein Test auf Form und Klasse der Codes im Quelltext; das Gate gegen Katalog und `SPEC-034` ist Folge-Slice (§1) | §1, `SPEC-034` §Stabilität |

**Akzeptierte Negative** (bewusst offen gelassen, mit Grund — die nächste Runde liest sie als entschieden):

- Abweichung in `param_types`, `param_formats` oder `result_formats` nennt weiter nur das Feld, ohne Nummer — es sind keine Werte, und keine Anforderung verlangt die Nummer; die Diagnose bleibt eindeutig.
- VT, FF, U+0085, U+2028 und U+2029 bleiben in der Zeile beim Prozessende roh stehen — sie brechen in üblichen Terminals und Log-Sammlern keine Zeile, und in der Log-Zeile escapt `logfmt` sie ohnehin.
- `time` der Log-Zeile ist nicht deterministisch — Tests prüfen `level`, `code` und `error`, nicht `time`.
- Was auf `debug` steht, ist nicht Vertrag — geprüft wird nur die Schwelle und dass kein Parameterwert erscheint.

**Risiken:**

- Vorrangfolge bei Fehlerketten mit mehreren Codes — entschieden in `SPEC-034` §Ausgabe *Fehlerkette*; der Re-Evaluierungs-Trigger von [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) ist geprüft: die Entscheidung bleibt, die Rangregel ist Spezifikationsdetail (`AGENTS.md` §3.8), keine Folge-ADR — **Ausgang:** offen bis Closure (entfällt, wenn der Kettentest aus DoD 1 grün ist).
- Parameterwerte bei `debug` (aus `slice-extended-query-replay`, Review F-328) — vom Nutzer am 2026-10-06 entschieden, `LH-FA-18.a` §Mismatch an `SPEC-033` angeglichen — **Ausgang:** offen bis Closure (entfällt mit dem Test aus DoD 3).
- `--log-level` liest seine Umgebungsvariable mit einem zweiten Einzel-Leser neben dem von `--fail-on-unconsumed`; der allgemeine Leser kommt mit `slice-v1-abschluss-betrieb` — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `observations/BEO-REPO/` am 2026-10-06 durchgegangen (17 Einträge, gemergter Stand). Treffer für diesen Slice:

- `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` — 2× (`slice-harness-kopf-sensor`, `slice-replay-semantik-mismatch`), offen. Ein dritter Beleg wäre die Schwelle; deshalb stehen die Randformen in §6 vor dem Code entschieden, und eine weitere hält der Implementer an, statt sie im Code zu setzen.
- `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` — 1×, offen, Regel in `.claude/commands/implement-slice.md` Schritt 19; betrifft die Mutationen nach `AGENTS.md` §3.10 dieses Slice.
- `BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar` — 1×, offen; betrifft die E2E-Mutation aus DoD 3, deshalb trägt `TestReplayMismatch` den Nachweis zusätzlich als Unit-Test.
- Verkörpert und mit diesem Slice als Retirement-Check berührt: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (§3.10), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (§3.11), `BEO-REPO/plan-folgt-korrektur-nicht` (§3.9), `BEO-REPO/spec-randform-erst-im-review-entschieden` (§3.12).
- Keine Treffer in den übrigen Einträgen; `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` trifft nicht zu, der Slice läuft in einer Welle und war vor dem Code beim Architect.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Sub-Area `*` ist Greenfield, `harness/conventions.md`: Spezifikation führt, Code folgt).
