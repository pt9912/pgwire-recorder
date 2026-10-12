# Verifikation: slice-v1-abschluss-umgebung-anderes-kommando — 2026-10-12

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-12-review-slice-v1-abschluss-umgebung-anderes-kommando.md`).

**Gegenstand:** Slice-Plan `slice-v1-abschluss-umgebung-anderes-kommando` (Kopf, §1 bis §8) gegen `2b553d5..1d160ce`. Commits: `b91131c` (Code, Tests, Abdeckung, Handbuch), `f1203f8` (Lint-Nachzug), `4da61d6` (Belege), `c31c201`/`ef5665c` (Review-Report), `e1a33d1` (Nacharbeit F-1/F-2), `1d160ce` (Belege der Nacharbeit).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

**Eingang:** Plan ganz (§7 *Belege des Implementers*), Diff von `cli.go`, Handbuch und Abdeckungstabellen, `LH-FA-17.a` (*Einstellungen*, *Fehler*, Optionstabelle) in der Spezifikation, Review-Report. Den Bericht des Implementers habe ich nicht gelesen.

## 1. Selbst gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make build` | Image `dev` unverändert zwischengespeichert, Binary als Entrypoint |
| `make docs-check` | 633 Dateien, 0 Befunde |
| `make abdeckung-check` | Exit 0 |
| `grep` im Handbuch auf `ADR-`, `SPEC-`, `LH-`, `[Ss]lice`, `Welle`, `Spezifikation`, `Lastenheft`, `[Rr]eview` | kein Treffer |
| `grep "unbeacht"` in `abdeckung-*.md` | nur Treffer, die `PGWIRE_RECORDER_PASSWORD`, nicht gelesene Optionen, PEM-Text, Variablen unbenutzter Verbindungen oder `BOGUS` meinen; keine Zeile sagt „unbeachtet“ für eine Variable, die jetzt Fehler ist |
| 31 Proben am gebauten Binary (Tabelle unten) | alle wie im Handbuch und in `LH-FA-17.a` |
| 6 Mutationen per Bind-Mount in `pgwire-recorder:test` auf einer Kopie des Baums (Tabelle unten), danach Rücksetzen | alle rot; der Arbeitsbaum blieb unberührt (`git status` leer) |

### Proben am Binary

Alle ohne Datenbank (Startfehler). Umgebung nur mit den genannten Variablen.

| Aufruf | Ergebnis |
|---|---|
| `play` mit `TLS_CERT=x` | Exit 2, `PGR-E2001`, nennt `PGWIRE_RECORDER_TLS_CERT` (F-1-Fall; der Handbuchsatz behauptet dazu nichts mehr) |
| `record` mit `TLS_CERT=x` | läuft (eigenes Kommando, unbeachtet), Exit 124 nach Timeout meiner Probe |
| `play` mit `LISTEN=127.0.0.1:1`, `SHUTDOWN_TIMEOUT=5s`, `FORCE=GEHEIM`, `FAIL_ON_UNCONSUMED=false` | je Exit 2, `PGR-E2001`, Name statt Wert (kein `GEHEIM` in der Meldung) |
| `record` mit `FAIL_ON_UNCONSUMED=true` | Exit 2, `PGR-E2001`, „die record nicht kennt“ |
| `replay` mit `FAIL_ON_UNCONSUMED=false` | kein `PGR-E2001` (nächster Fehler: Aufzeichnung fehlt, `PGR-E3001`) |
| `record` mit leerer `FAIL_ON_UNCONSUMED` | Exit 2, `Pflichtoption --listen fehlt` (leer gilt als nicht gesetzt) |
| `PGWIRE_RECORDER_PASSWORD` bei `record`, `replay`, `play` | unbeachtet (`Pflichtoption … fehlt`) |
| `BOGUS=x` bei `record`, `replay`, `play`; `PGWIRE_RECORDER_=x`; `pgwire_recorder_listen=x`; `PGWIRE_RECORDER_Listen=x` bei `play` | unbeachtet |
| `config show` mit `LISTEN=x`, `FORCE=GEHEIM` | Exit 0, listet die Namen; mit `CONFIG=/gibt/es/nicht.yaml` Exit 2, `PGR-E2004` |
| `play --help`, `record -h`, `version` mit fremder Variable | Exit 0, Hilfe bzw. Version |
| `play` mit `SHUTDOWN_TIMEOUT=5s` und `LISTEN=x` | nur die in der Tabelle erste (`LISTEN`) |
| `play --fail-on-unconsumed` mit fremder `LISTEN` | Kommandozeile geht vor (`flag provided but not defined`) |
| `play` mit fremder `LISTEN` und `CONFIG=/gibt/es/nicht.yaml` | `PGR-E2001` vor der Datei |
| `play` mit `KEEP_TIMING=GEHEIM` und `--keep-timing=false` | `PGR-E2001` für die eigene Variable, unabhängig von der Priorität |
| `play` mit `INPUT` leer, `SHUTDOWN_TIMEOUT=5s` | `PGR-E2001` für `SHUTDOWN_TIMEOUT` |
| `play` mit ungültiger eigener `FORMAT=bad` (unbekannte Option) und fremder `LISTEN` | `LISTEN` (fremd steht in der Tabelle vor `FORMAT`) |
| `config show` mit `LOG_LEVEL=info` | Exit 0 |

### Mutationen (selbst gefahren)

Weg: Kopie von `cmd`, `internal`, `go.mod`, `go.sum` im Scratchpad, Bind-Mount in `pgwire-recorder:test`, `go test ./internal/adapters/driving/cli/`. Basislauf grün.

| Zusage | Mutation | roter Test |
|---|---|---|
| fremde, gesetzte, nicht leere Variable ist `PGR-E2001` | Prüfung entfernt (`if false && …`) | `TestLeserUmgebungAndererKommandos`, `TestLeserUmgebungNichtGelesenerOptionen`, `TestLeserUmgebungReihenfolge`, `TestParseFailOnUnconsumedRecord`, `TestParsePlayFailOnUnconsumed`, `TestParsePlayFremdeUmgebung` |
| Option → Kommando richtig | `output` auch `replay` zugeordnet | `TestLeserUmgebungAndererKommandos` |
| leer gilt als nicht gesetzt | `LookupEnv` nur im fremden Zweig | `TestLeserUmgebungAndererKommandos` (zusätzlich `TestLeserFremdeUmgebung`, `TestLeserUmgebungNebenKommandozeile`, `TestLeserAlleOptionen` u. a.) |
| leer gilt als nicht gesetzt (allgemein) | `LookupEnv` für alle | 8 und mehr Tests rot (`TestDatei*`, Aufbau bricht) |
| Meldung nennt nie den Wert | `NAME=Wert` in der Meldung | sechs Tests, darunter `TestLeserUmgebungAndererKommandos` |
| jeder Wert, gültig oder nicht | fremde Variable nur bei `v != "true"` ablehnen | `TestLeserUmgebungAndererKommandos`, `TestLeserUmgebungReihenfolge`, `TestParseFailOnUnconsumedRecord`, `TestParsePlayFailOnUnconsumed` |

Erster Versuch der Mutation „leer gilt als gesetzt“ baute nicht (von mir falsch geschrieben, Variablenname `ok` doppelt); mit anderem Namen wiederholt. Keine grüne Mutation. Rücksetzen: `diff` gegen die Vorlage leer, `git status` im Repo leer.

## 2. DoD Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| Liefer-Punkt 1: Prüfung im CLI-Adapter nach §1 und `LH-FA-17.a` (gesetzt, nicht leer → `PGR-E2001` mit Namen ohne Wert; leer = nicht gesetzt; `config show`, Hilfe, `version`; Ausnahmen) | erfüllt | `cli.go`: `zuordnung()` (26 Zeilen gegen die Tabelle, `--config` und Passwort keine Zeilen), `liesUmgebung` in einem Durchlauf in Tabellenreihenfolge, danach Datei; alle Zusagen am Binary geprobt (oben) |
| …die vier Tests aus §1 folgen dem Vertrag | erfüllt | `TestParseFailOnUnconsumedRecord`, `TestLeserFremdeUmgebung`, `TestParsePlayFailOnUnconsumed`, `TestParsePlayFremdeUmgebung` im Diff, im Basislauf grün, in den Mutationen rot |
| …jede Zusage hat einen Test, der unter ihrer Mutation rot wird | erfüllt, soweit ich fuhr | 6 Mutationen selbst rot; die übrigen Zeilen der Tabelle in §7 (Rangfolge, Hilfe, `config show`, `EqualFold`, Passwort, `BOGUS`, Gleichheitstest) habe ich nicht nachgefahren, die Testnamen existieren und lesen sich passend |
| …Abdeckungstabellen folgen den Deklarationen | erfüllt | `make abdeckung-check` Exit 0; Zeilen zu den vier Tests neu gefasst; `TestPlayPasswort` unverändert |
| Liefer-Punkt 2: Handbuch im Ist-Zustand, Abschnitt *Einstellungen* und Zeile `PGR-E2001` | erfüllt | Absatz in §5 und Katalogzeile (Anker `#5-einstellungen` löst auf, `make docs-check` grün); kein Verweis nach oben, keine Chronik; jeder Satz hat eine Probe (Abschnitt 3) |
| `make gates` grün | erfüllt laut Auftraggeber (Exit 0 an `1d160ce`) und §7 (Stände `f1203f8`, `e1a33d1`); nicht von mir gefahren | |
| Review, Report unter `docs/reviews/` | erfüllt | 0 HIGH, 0 MEDIUM, 2 LOW (F-1, F-2), 2 INFO; F-1 und F-2 behoben (Abschnitt 4) |
| Mutationstabelle und Suche in §7 | erfüllt | Tabelle 19 Zeilen; die Suche nennt 29 Dateien; ich wiederholte sie nicht (die Tests sind grün, `make gates` laut Auftraggeber grün) |
| Closure-Notiz mit Lerneintrag | offen (Closure) | |
| Beobachtungs-Register | offen (Closure) | |
| Risiken mit Ausgang (vier in §6) | offen (Closure) | die Ausgänge in §6 stehen als Vorschlag da; die Closure bestätigt sie |
| Paarungen | offen (Closure) | |

## 3. Verhaltenssätze des Handbuchs gegen die Probe-Tabelle in §7 (AGENTS.md §3.11)

| Satz im Handbuch | Zeile in §7 | eigene Probe |
|---|---|---|
| gesetzt und nicht leer ist `PGR-E2001`, gleich welchen Wert (`true`, `5s`, ungültig) | ja (drei Werte) | bestätigt |
| nennt Namen, nie den Wert | ja (`GEHEIM`) | bestätigt |
| `FAIL_ON_UNCONSUMED` bei `record` und `play`, `LISTEN` bei `play` | ja | bestätigt |
| Betriebsarten gemeinsam in einer Umgebung: Variablen der anderen Betriebsart nicht setzen | folgt aus den ersten Zeilen | bestätigt |
| `config show` prüft nur `PGWIRE_RECORDER_CONFIG` | ja | bestätigt |
| unbeachtet: `PASSWORD` bei `record`/`replay`; `BOGUS` bei allen drei; `pgwire_recorder_listen` | ja (BOGUS bei allen drei seit `1d160ce`) | bestätigt |
| Zeile `PGR-E2001` des Katalogs | ja | bestätigt |
| (unverändert) leer gilt als nicht gesetzt | ja | bestätigt |

Kein Satz ohne Zeile gefunden. Kein Satz im Handbuch behauptet etwas über die sechs noch nicht gelesenen Optionen; der unveränderte Satz zur Hilfe und die ändernden Sätze lesen sich für `PGWIRE_RECORDER_TLS_CERT` bei `play` nicht falsch mit (das Binary lehnt ab, das Handbuch sagt dazu nichts, denn die Tabelle führt die Option nicht). Der Satz „Die Option gibt es nur bei `replay`“ zu `--fail-on-unconsumed` bleibt wahr.

## 4. Review-Findings

| Finding | Behoben? | Beleg |
|---|---|---|
| F-1 Handbuchsatz zu Variablen mit Präfix allgemeiner als die Optionsmenge | ja | Satz nennt nur `BOGUS`, `pgwire_recorder_listen` und `PASSWORD` bei `record`/`replay`; `play` mit `TLS_CERT=x` am Binary ist `PGR-E2001`, und das Handbuch behauptet das Gegenteil nicht mehr; Probe `BOGUS` bei allen drei Betriebsarten ergänzt und von mir bestätigt |
| F-2 Plan §3 nennt `abdeckung-gesamt.md` und `-vollstaendig.md` nicht | ja | §3 führt beide und die Deklaration `LH-FA-01/Boundary`; `make abdeckung-check` Exit 0 |
| F-3, F-4 (INFO) | entschieden, keine Findings | |

## 5. Befunde

**HIGH / MEDIUM:** keine.

**LOW:**

- V-1 (Plan, Zeilennummern): Plan §1, §3 und DoD Liefer-Punkt 1 nennen die Zeilen 22, 172, 176, 219, 224 von `docs/user/abdeckung-unit.md`. Nach dem Diff stehen die Zeilen bei 173, 177, 180 und weiteren; die Nummern gelten nur für den Stand vor `b91131c`. Kein inhaltlicher Fehler; die Closure nennt Test-Namen statt Zeilen, damit die Angabe nicht veraltet. `verifizierbar`: ja (`grep -n` auf die vier Test-Namen in `abdeckung-unit.md`).

**INFO:**

- V-2: Die Mutationstabelle in §7 führt bei „leer gilt als nicht gesetzt (allgemein)“ „68 Tests“; ich sah beim Nachfahren der Variante (`LookupEnv` für alle) acht und mehr `TestDatei*`-Fälle rot, die Zahl habe ich nicht gezählt. Die Zusage trägt (§7 sagt selbst, die Zeile belege sie nur zusammen mit der Zeile darüber).
- V-3: `record` mit `PGWIRE_RECORDER_TLS_CERT` läuft an und ignoriert die Variable (Stand des Codes, vom Architect entschieden); für den Nutzer des Handbuchs ist das nicht beschrieben, weil die Option dort fehlt. Kein Befund, nur Hinweis für die Slices, die den Leser liefern.

## 6. Was ich nicht selbst gefahren habe

- `make gates` (volle Kette) und `make test-integration`; grün laut Auftraggeber an `1d160ce`.
- Die Mutationen der Zeilen Rangfolge (Variable vor Datei, Durchlauf-Reihenfolge, nur erster Fehler), Hilfe vor Prüfung, `config show` prüft nur `CONFIG`, `EqualFold`, Passwort wie Option, jede Variable mit Präfix fremd, `optionen("play")` um eine Option erweitert, „Name weggelassen“. Ich lese die Tests passend; am Binary habe ich die zugehörigen Verhalten geprobt (Abschnitt 1).
- Die Suche nach Tests, die Variablen setzen (29 Dateien), und die `grep`-Suche nach „unbeacht“ im Bestand außerhalb der Abdeckungstabellen.
- Probe mit laufender Datenbank (nicht nötig; alle Zusagen sind Startfehler).
- Closure-Punkte (Notiz, Register, Risiko-Ausgänge, Paarungen): offen, führt die Closure.

## 7. Urteil

DoD-Liefer-Punkte 1 und 2 erfüllt, Entscheidungen des Architect (§6) im Code und am Binary wiederzufinden, Plan folgt dem Diff (V-1 ausgenommen). Keine DoD-Verletzung, die Tests oder Review nicht sehen würden. Offen bleibt die Closure.
