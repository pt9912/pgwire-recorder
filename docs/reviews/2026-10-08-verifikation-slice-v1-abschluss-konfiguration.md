# Verifikation: slice-v1-abschluss-konfiguration — 2026-10-08

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-08-review-slice-v1-abschluss-konfiguration.md`, F-496 bis F-505). F-503 bis F-505 hat er an mich übergeben.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-konfiguration.md` (der Leser-Slice nach dem Schnitt), Kopf und §1 bis §8, gegen den Gesamt-Diff `c060323..960f239`. Architect: `1101e45`, `c060323`, `096e51e`; `f19b440` betrifft den Datei-Slice und diente nur zur Abgrenzung. Schnitt: `fa4a5f1`. Implementer: `8fad493`, `a968790`, `8f7fe90`, `960f239`. Review: `93164a5`.

**Eingang:**

- der Plan ganz, besonders §1, §2 (DoD), §3, §4 (Schnitt), §6 (*Randformen*, *Risiken*) und §7 (*Belege des Implementers*)
- `spec/spezifikation.md` `LH-FA-17.a` ganz, mit Optionstabelle, *Fehler* und *Dauer*; `LH-FA-01.a` *Hilfe vor Prüfung* und *Ende der Optionen*
- [ADR-0036](../plan/adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), `.a-check.yml` (Abschnitt `tech`)
- am Stand `960f239` `internal/adapters/driving/cli/cli.go`, `leser_test.go` und `export_test.go` ganz, `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk`, die Zeile in `harness/README.md` §Sensors, `docs/user/benutzerhandbuch.md` §5
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen
- die Adressen in `slice-v1-abschluss-schreiben` und `slice-v1-abschluss-konfigurationsdatei` (je §1)

Beim Start war der Arbeitsbaum sauber, HEAD `960f239`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-08

**So wurden die Proben gebaut.** Für jede Probe eine frische Kopie per `git archive 960f239 | tar -x`, ohne `.git` und ohne `cp -p`, in einem eigenen Unterverzeichnis des Scratchpads, nie im Repo. Die Mutation setzte ein Skript, das abbricht, wenn der Suchtext nicht genau einmal trifft. Danach `chmod -R a+rX` und in der Kopie `make test` bzw. `make a-check-negativ`. Für das Binary je ein Image der Stufe `runtime` aus einer Kopie von `960f239` und von `c060323` mit eigenem Tag, aufgerufen mit `docker run --rm`. Diese Images sind danach entfernt. `make test` in den Kopien taggt dasselbe Test-Image wie im Repo, der nächste Lauf im Repo baut es neu. `make gates` lief im Arbeitsbaum vor den Mutanten. `git status` war danach sauber.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Gegenprobe des Architektur-Gates ([ADR-0036](../plan/adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md)). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| nimmt beide YAML-Modulpfade im CLI-Adapter an | Fall 8 `cli-yaml` importiert `go.yaml.in/yaml/v3` und `gopkg.in/yaml.v3` in einer Datei, Erwartung grün; Zeilen 1 und 2 der Tabelle in §7. Ungeändert sind alle zwölf Fälle grün (`make gates`, eigener Lauf V-A0: Exit 0) | bestätigt |
| lehnt jeden Modulpfad im PGWire-Adapter ab (`ARC-013`) | Fälle 9 `pgwire-yaml` und 10 `pgwire-gopkg-yaml`, je `tech-leak`; §7 Zeilen 3 und 4 | bestätigt |
| zwölf Fälle; die Ablehnung gilt je Modulpfad im Domain Model, im Postgres-Adapter und im PGWire-Adapter (F-499) | Fälle 6, 7, 11, 12 im Skript gelesen. **A1 nachgefahren** (Regel `go.yaml.in/yaml` von `driven/recording` auf `internal/adapters/driven` geweitet): **rot**, `a-check-negativ: ROT in Fall 'postgres-yamlin' — erwartet rot, a-check war gruen`. Im Review war A1 noch grün | bestätigt |
| Kopfkommentar, Fragment und Zeile in `harness/README.md` §Sensors sagen nicht mehr zu, als die Fälle prüfen | gelesen: Alle drei nennen je Modulpfad Domain Model, Postgres-Adapter und PGWire-Adapter (abgelehnt) und den CLI-Adapter (zugelassen). Genau das prüfen die Fälle 6 bis 12. Zur Schlusszeile des Skripts siehe V-115 (nicht Gegenstand der DoD) | bestätigt |
| `.a-check.yml` führt beide Modulpfade für Recording- und CLI-Adapter | Abschnitt `tech` gelesen, deckt sich mit [ADR-0036](../plan/adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) *Entscheidung* | bestätigt |

### Punkt 2: allgemeiner Leser, Priorität Kommandozeile > Umgebung > Default (`LH-FA-17`). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| jede Option von `record` und `replay`, auch `--output`, `--force` und `--shutdown-timeout`, ist am Leser angemeldet | `optionen()` in `cli.go`: `record` mit listen, upstream, output, force, shutdown-timeout, log-level; `replay` mit listen, input, fail-on-unconsumed, shutdown-timeout, log-level. Das ist die Reihenfolge der Optionstabelle in `LH-FA-17.a`, ohne die noch nicht gelieferten Optionen. `parseRecord`/`parseReplay` rufen nur `lies`, keine Option geht am FlagSet vorbei an | bestätigt |
| Priorität Kommandozeile vor Umgebung vor Default | `TestLeserAlleOptionen` je Option; §7 Zeile 1. Binary: `PGWIRE_RECORDER_LISTEN=127.0.0.1:5999` neben `--listen=bad-addr` endet mit `PGR-E4001 … Adresse bad-addr`, die Kommandozeile gilt. Allein aus der Umgebung (listen, upstream, output, log-level) startet `record` | bestätigt |
| Default aus der Optionstabelle, nicht aus der eigenen Anmeldung (F-505) | `tabelle()` in `leser_test.go` hält die Tabelle, der Test vergleicht „ohne Quelle“ mit `--<option>=<Default der Tabelle>`. **L9 nachgefahren** (Default von `--log-level` bei `record` `warn`): **rot**, jetzt auch `TestLeserAlleOptionen`, dazu `TestParseLogLevel`, `TestParseRecord`, `TestRunRecordLogLevel` und fünf weitere | bestätigt, F-505 erledigt |
| jede gesetzte Umgebungsvariable wird geprüft, auch wenn die Kommandozeile vorgeht (`PGR-E2001`) | §7 Zeile 3. Binary: `PGWIRE_RECORDER_FORCE=ja` neben `--force` → `PGR-E2001: Umgebungsvariable PGWIRE_RECORDER_FORCE: erlaubt sind true und false`, Exit 2 | bestätigt |
| `--output` und `--force` am Leser (F-496), `PGWIRE_RECORDER_OUTPUT` genügt, strenge Werte für `--force` | Binary: `record --listen 127.0.0.1:5999 --upstream h:1` mit `PGWIRE_RECORDER_OUTPUT=/tmp/r.yaml` startet und endet erst durch den Abbruch der Probe (kein `PGR-E2001`). `PGWIRE_RECORDER_FORCE=ja` → `PGR-E2001`; `--force=1` → `PGR-E2001`. **F1 nachgefahren** (Wertemenge von `--force` nach `strconv.ParseBool`): **rot** in `TestLeserAlleOptionen` (`PGWIRE_RECORDER_FORCE=1 neben --force=false: … erhalten <nil>`) und `TestLeserHilfe` (`--force=1: <nil>`) | bestätigt |
| leere Umgebungsvariable gilt als nicht gesetzt | Binary: `PGWIRE_RECORDER_LISTEN=` → `Pflichtoption --listen fehlt`; `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT=` und `PGWIRE_RECORDER_LOG_LEVEL=` ohne Wirkung. §7 Zeile 4 | bestätigt |
| leerer Wert auf der Kommandozeile ist `PGR-E2001`, auch neben gesetzter Variable (F-497) | Binary: `record --listen= --upstream h:1 --output r.yaml` mit `PGWIRE_RECORDER_LISTEN=127.0.0.1:1` und `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT=x` → `PGR-E2001 … invalid value "" for flag -listen: ein leerer Wert ist ungültig`. Die Meldung nennt die Option, weder die Variable noch `Pflichtoption`, wie §6 es als Testfall fordert. **L6 nachgefahren** (leerer Wert in `kommandozeile.Set` übergangen): **rot**, `TestLeserAlleOptionen`, `TestLeserLeererWert`, die drei `TestParse…Werte`, `TestRunStartfehlerJeStufe` | bestätigt |
| einzige Stelle der Anmeldung (F-498) | **L3 nachgefahren** (`fs.String("session-assignment", …)` bei `replay` am Leser vorbei): **rot**, `TestLeserNurAngemeldete` (`replay --session-assignment=wert: erwartet unbekannte Option, erhalten <nil>`) | bestätigt |
| Reihenfolge der Fehler: Kommandozeile, Umgebung nach der Tabelle, zuletzt Pflichtoptionen | Binary mit `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT=5` und `PGWIRE_RECORDER_LOG_LEVEL=INFO`: `--log-level=trace` meldet die Kommandozeile, `zusatz` meldet `unerwartetes Argument "zusatz"`, ohne Fehler auf der Kommandozeile kommt `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT` vor `…_LOG_LEVEL`. `…_FAIL_ON_UNCONSUMED=1` ohne Optionen wird vor `Pflichtoption` gemeldet. `record` mit `…_FORCE=x` und `…_LOG_LEVEL=x` meldet `FORCE`. Ohne alles kommt `Pflichtoption --listen fehlt`. §7 Zeilen 7 bis 9 | bestätigt |
| Dauer-Form (`LH-FA-17.a` *Dauer*) | Binary über `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT`: `0`, `05s`, `00s`, `1ms` angenommen (es folgt `PGR-E3001` der fehlenden Aufzeichnung). `00`, `1m30s`, `1.5s`, `-1s`, `5S`, ` 5s`, `5` → `PGR-E2001`. `9999999999999m` → `PGR-E2001 … länger als die längste darstellbare Dauer` | bestätigt |
| Umgebungsvariable einer fremden Option oder mit Präfix ohne Option unbeachtet | `TestLeserFremdeUmgebung` (fünf Namen); Binary: `record` mit `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=x` läuft weiter bis `PGR-E3001` | bestätigt |
| ersetzt die drei Einzel-Leser, deren Tests unverändert grün bleiben | `git diff --stat c060323 960f239 -- internal/ cmd/`: nur `cli.go`, `export_test.go` (nur ergänzt) und `leser_test.go` (neu). `cli_test.go`, `frist_test.go` und `internal/bootstrap` sind unverändert und im Gate grün. Gleiches Verhalten siehe Abschnitt 3 | bestätigt |

### Punkt 3: Hilfe vor jeder Prüfung am Leser (`LH-FA-01`). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Hilfe vor jeder Prüfung, auch neben ungültiger Umgebungsvariable | Binary: `record --help` mit `…_SHUTDOWN_TIMEOUT=x`, `…_FORCE=ja`, `…_LOG_LEVEL=INFO` → Hilfe auf stdout, stderr leer, Exit 0. `replay --listen= --bogus -h` mit ungültiger Variable → Hilfe, Exit 0. `record --force=1 --help=false` mit `…_FORCE=ja` → Hilfe, Exit 0. `replay --input -help` → Hilfe, Exit 0. **Eigener Mutant H1** (Umgebungsvariablen von `record` am Anfang von `Parse` vor der Hilfe geprüft): **rot**, `TestParseLogLevelHilfe`, `TestParseShutdownTimeoutHilfe` (`["record" "--help"]: erwartet Hilfe, erhalten … Umgebungsvariable …`), dazu `TestLeserReihenfolge`, `TestLeserLeererWert`, `TestRunVersionLogLevel` | bestätigt |
| nennt die Umgebungsvariablen und ihre Priorität, ohne Ausnahme | Hilfetext des Binary: „Jede Option ist auch über ihre Umgebungsvariable setzbar: PGWIRE_RECORDER_ und der Name in Großbuchstaben mit _ statt -, etwa PGWIRE_RECORDER_LISTEN; die Kommandozeile geht ihr vor. Ein leerer Wert auf der Kommandozeile ist ungültig.“ Ohne „außer“. §7 Zeilen *Hilfe nennt …* | bestätigt |
| `--` beendet die Optionen | Binary: `-- --help` → `unerwartetes Argument "--help"`, Exit 2. `--input --` → `flag needs an argument: -input`, `PGR-E2001`. `--input=--` nimmt `--` als Wert. §7 Zeilen 4 und 5 der Hilfe-Tabelle | bestätigt |
| `version` liest keine Umgebungsvariable | Binary: `version` mit `…_SHUTDOWN_TIMEOUT=x` → Exit 0 | bestätigt |
| Beleg in §7 je Zusage: Zusage · Mutation · roter Test | Tabellen zu Punkt 1 (6 Zeilen), Punkt 2 (20 Zeilen), Punkt 3 (7 Zeilen) | vorhanden |

### Punkt 4: `make gates` grün. **Bestätigt durch eigenen Lauf, Beleg in §7 fehlt (V-113).**

Eigener Lauf im Arbeitsbaum am Stand `960f239`: Exit 0. Darunter `a-check-negativ: gruen`, `d-check: 370 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: … OK`, `run-integration-tests: gruen`, `kopf-check-gegenprobe: gruen`, `abdeckung-gegenprobe: gruen`.

### Übrige DoD-Punkte (konstant je Slice)

- Review: liegt vor (`93164a5`). F-496 bis F-502 und F-505 haben einen Ausgang (Abschnitt 2).
- Closure-Notiz, Register, Risiko-Ausgänge, Paarungen: In §7 stehen noch Platzhalter, und das Risiko *Einzel-Leser* steht in §6 mit „offen bis Closure“. Das ist Closure-Arbeit nach dieser Verifikation, kein Befund. Für das Risiko trägt Abschnitt 3 den Ausgang *entfallen*.

---

## 2. Review-Befunde F-496 bis F-505: nachgefahren

| Befund | Ausgang | Nachgefahren | Urteil |
|---|---|---|---|
| F-496 `PGWIRE_RECORDER_OUTPUT`/`_FORCE` | Architect `096e51e`, umgesetzt in `960f239`, §1 und §6 | Binary (Punkt 2), Mutant F1 rot | umgesetzt |
| F-497 leerer Wert auf der Kommandozeile | `LH-FA-17.a` (`096e51e`), Handbuch §5, Code | Binary, Mutant L6 rot | umgesetzt |
| F-498 einzige Stelle | `TestLeserNurAngemeldete` | Mutant L3 rot | umgesetzt |
| F-499 Sensor-Text weiter als die Fälle | Fälle 11 und 12 | Mutant A1 rot | umgesetzt |
| F-500 unerwartetes Argument vor Umgebung | `TestLeserReihenfolge` (zwei Fälle) | Binary (Punkt 2); L8 nicht selbst gefahren, die Testzeilen 223 bis 230 gelesen | umgesetzt |
| F-501 Kommentar an `wahrheitswert` | `IsBoolFlag` und `String` der Prüfer entfernt, Kommentar nennt die Rolle als Prüfer | gelesen, `cli.go` Z. 388 bis 392 | umgesetzt |
| F-502 Zählung in §4 | §4 nennt die Zählung vor und nach dem Schnitt | gelesen | umgesetzt |
| F-503 Lifecycle ohne `git mv` | Architect bestätigt in §4 | siehe V-117 | an den Planner |
| F-504 gleiches Verhalten | — | Abschnitt 3 | bestätigt |
| F-505 Default selbstbezüglich | `tabelle()` im Test | Mutant L9 rot in `TestLeserAlleOptionen` | umgesetzt |

---

## 3. Gleiches Verhalten der ersetzten Einzel-Leser (Stichprobe gegen `c060323`)

Ich habe dieselben Aufrufe mit dem Binary von `c060323` und dem von `960f239` gefahren, 29 Fälle zu `--fail-on-unconsumed`, `--log-level` und `--shutdown-timeout`. Verglichen wurden Ausgabe und Exit-Code. Darunter waren Schalter ohne Wert, `=`, leerer Wert, `1`, `INFO`, `1m30s`, `00`, `05s`, Mehrfachangabe in beiden Reihenfolgen, ungültige Variable allein und neben allen drei Optionen sowie leere Variable. **26 Fälle sind gleich**, Code und Text eingeschlossen. Abweichend sind drei:

- `record` mit `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=x`: Beide melden dasselbe `PGR-E3001`, nur der Zufallsname der Probe-Datei ist anders. Gleiches Verhalten.
- `record --force=x` und `--force=`: Beide sind `PGR-E2001`, der Grund lautet jetzt „erlaubt sind true und false“ statt „parse error“. Das ist die gewollte Änderung aus F-496. `--force=1` ist jetzt ebenfalls `PGR-E2001`, wie `LH-FA-17.a` es verlangt.

Die Reihenfolge-Abweichung aus F-504 (a) habe ich nicht erneut gezogen; mein Lauf bestätigt die Reihenfolge nach `LH-FA-17.a` *Fehler* (Punkt 2). **Risiko *Einzel-Leser* (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`):** Ich schlage als Ausgang *entfallen* vor. Begründung: Die Tests der Einzel-Leser sind unverändert grün, und in Review und Verifikation ist keine Prüfung weggefallen.

---

## 4. Abgrenzung gegen `slice-v1-abschluss-konfigurationsdatei`

- **Nichts aus dem Datei-Slice liegt im Code.** Ein `grep` nach `yaml` und `config` in `internal/adapters/driving/cli/*.go` ohne Tests findet nichts. `go.mod` und `go.sum` sind im Diff unverändert. `config show` und `--config` gibt es nicht: `"config"` steht nur in der Liste `alleOptionen()` von `leser_test.go`, und dort muss es eine unbekannte Option sein. `PGR-E2004` bis `PGR-E2006` kommen im Diff nicht vor.
- **Es fehlt nichts, was der Leser-Slice zusagt.** Alle Zusagen aus §1 sind belegt: die Übernahmen aus `slice-replay-semantik-mismatch`, `slice-v1-abschluss-herunterfahren` und `slice-v1-abschluss-schreiben`, die einzige Stelle der Anmeldung, `parseRecord` ohne eigene Anmeldung und die Gegenprobe. Die Adressen nehmen an: `slice-v1-abschluss-schreiben` §1 nennt *Abgegeben* an diesen Slice, `slice-v1-abschluss-konfigurationsdatei` §1 *Übernimmt* mit seiner Kennung. Beide liegen in `next/`.
- **Schicht-Abgrenzung:** Der Diff berührt nur den CLI-Adapter und die Gegenprobe. Kern, PGWire-Adapter, Driven-Adapter und Bootstrap sind unberührt.

---

## 5. Handbuch §5 gegen das Verhalten

Was dieser Slice liefert, stimmt mit dem Binary überein: Priorität, Umgebungsvariable je Option der Tabelle, leerer Wert auf der Kommandozeile `PGR-E2001` auch neben gesetzter Variable, leere Variable nicht gesetzt, Wahrheitswerte `true`/`false`, Hilfe vor jeder Prüfung auch bei ungültiger Umgebungsvariable, `--` als Ende der Optionen, `--input=--`, keine Option `--version`. §5 beschreibt daneben den Zielstand (Konfigurationsdatei, `--format`, `play`, TLS-Optionen), den das Binary noch nicht hat. Plan §3 sagt das ausdrücklich („schon im Zielstand“). Darum ist das hier kein Befund.

---

## 6. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-113 | LOW | DoD-Punkt 4 (`make gates` grün) hat in §7 *Belege des Implementers* keinen Beleg. §7 nennt nur `make test` und `tools/arch/a-check-negativ.sh` in Kopien. Mein eigener Lauf am Stand `960f239` ist grün (Abschnitt 1, Punkt 4). Die Bestätigung hängt damit an diesem Bericht, nicht am Beleg des Implementers. | `docs/plan/planning/in-progress/slice-v1-abschluss-konfiguration.md` §7 | Implementer: Zeile mit Lauf und Stand nachtragen, oder auf diesen Bericht zeigen |
| V-114 | LOW | §3 folgt der Nacharbeit nicht ganz (`AGENTS.md` §3.9). Die Zeile `internal/adapters/driving/cli` sagt „update (geliefert, a968790)“ und „(`--output` und `--force` offen, Entscheidung zu F-496)“. Beide Optionen sind aber seit `960f239` am Leser geliefert. Die Test-Zeile nennt ebenfalls nur `a968790`, obwohl `leser_test.go` in `960f239` wuchs (`TestLeserNurAngemeldete`, `TestLeserLeererWert`, `tabelle()`). „offen“ liest sich als nicht geliefert. | `docs/plan/planning/in-progress/slice-v1-abschluss-konfiguration.md` §3 · „`--output` und `--force` offen“ | Implementer |
| V-115 | INFO | Die Schlusszeile von `tools/arch/a-check-negativ.sh` lautet „YAML nur im Recording- und CLI-Adapter“. Das ist weiter als die Fälle: Geprüft wird die Ablehnung im Domain Model, im Postgres- und im PGWire-Adapter. Services, Ports und die Composition Root (dort ist jeder Import zulässig) prüft die Gegenprobe nicht. Die DoD nennt Kopfkommentar, Fragment und README-Zeile, und diese drei sind genau. Die Formulierung „nur …“ stand schon vor dem Slice so („nur im Recording-Adapter“). | `tools/arch/a-check-negativ.sh` · „YAML nur im Recording- und CLI-Adapter“ | Planner: eigener Vorgang oder Closure-Beobachtung zu `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` |
| V-116 | INFO | Ein ungültiger Wert auf der Kommandozeile erscheint in der Meldung, etwa `invalid value "trace" for flag -log-level`. `LH-FA-17.a` *Fehler* sagt: „Jede Ursache … nennt in der Meldung die Stelle …, nie einen Wert.“ Ob der Satz auch für `PGR-E2001` aus der Kommandozeile gilt oder nur für die Ursachen der Datei, sagt die Spezifikation nicht eindeutig. Das Verhalten ist Bestand: Für die drei früheren Optionen ist es identisch mit `c060323` (Abschnitt 3). Neu ist es für Text-Optionen, bei denen aber nur der leere Wert scheitert (`""`). Meldungen zu Umgebungsvariablen nennen keinen Wert. | `internal/adapters/driving/cli/cli.go` · `"ungültige Verwendung von %s"` | Architect: Lesart klären, bevor `slice-v1-abschluss-konfigurationsdatei` Meldungen mit Werten aus der Datei baut |
| V-117 | INFO | F-503: Nach dem Schnitt blieb der Slice in `in-progress/`, ohne den Übergang `in_progress → next`, den die State Machine für *zu groß* vorsieht (`v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Lifecycle als State Machine). Die Abweichung ist nicht still: §4 nennt Bedingung und Grund, der Architect bestätigt, und `fa4a5f1` trägt den Schnitt. Ein `MR` deklariert diese Lesart aber nicht. Ich bestätige, dass die Pflichten der Rückführung (Bedingung vorab, Grund nachgetragen) in §4 stehen. Ob der Verbleib als Lesart der Baseline trägt, entscheide ich nicht. | `docs/plan/planning/in-progress/slice-v1-abschluss-konfiguration.md` §4 · „Der Schnitt wurde ohne“ | Planner: Lerneintrag oder `MR` erwägen, falls sich das wiederholt |

---

## 7. Gefahrene Mutanten

| Mutant | Änderung | Ziel | Ergebnis |
|---|---|---|---|
| V0 | keine | `make test` | grün |
| L6 | `kommandozeile.Set`: leerer Wert gilt als nicht gesetzt | `make test` | rot: `TestLeserAlleOptionen`, `TestLeserLeererWert`, `TestParseFailOnUnconsumedWerte`, `TestParseLogLevelWerte`, `TestParseShutdownTimeoutWerte`, `TestRunStartfehlerJeStufe` |
| L3 | `lies`: `--session-assignment` bei `replay` direkt am FlagSet | `make test` | rot: `TestLeserNurAngemeldete` |
| L9 | Default von `--log-level` bei `record` `warn` | `make test` | rot: `TestLeserAlleOptionen` und acht Einzeltests |
| F1 | Wertemenge von `--force` nach `strconv.ParseBool` | `make test` | rot: `TestLeserAlleOptionen`, `TestLeserHilfe` |
| H1 | `Parse` prüft zuerst die Umgebungsvariablen von `record`, dann die Hilfe | `make test` | rot: `TestParseLogLevelHilfe`, `TestParseShutdownTimeoutHilfe`, `TestLeserReihenfolge`, `TestLeserLeererWert`, `TestRunVersionLogLevel` |
| A0 | keine | `make a-check-negativ` | grün |
| A1 | Regel `go.yaml.in/yaml` von `driven/recording` auf `internal/adapters/driven` geweitet | `make a-check-negativ` | rot: Fall `postgres-yamlin` |

Jeder rote Mutant scheitert aus dem richtigen Grund: Die erste Fehlerzeile nennt die mutierte Zusage. Kein Mutant blieb grün.

---

## Verdikt

**DoD-Liefer-Punkte 1 bis 3 und `make gates`: bestätigt.** Das Verhalten des Binary entspricht `LH-FA-17.a` und `LH-FA-01.a`, soweit der Slice sie zusagt, und der Slice hält die Abgrenzung zum Datei-Slice ein. Es gibt keinen blockierenden Befund. Vor der Closure: V-113 und V-114 an den Implementer (Plan-Nachtrag ohne Code). V-116 an den Architect. V-115 und V-117 an den Planner. Danach stehen nur noch die Closure-Pflichten in §7 aus, mit dem Ausgang *entfallen* für das Risiko *Einzel-Leser* (Abschnitt 3).

**Summary:** 0 HIGH · 0 MEDIUM · 2 LOW · 3 INFO (V-113 Gate-Beleg fehlt in §7; V-114 §3 nennt `--output`/`--force` „offen“; V-115 Schlusszeile der Gegenprobe weiter als die Fälle; V-116 Wert in der Meldung der Kommandozeile, Lesart von `LH-FA-17.a` *Fehler* offen; V-117 Verbleib in `in-progress/` nach dem Schnitt ohne `MR`).
