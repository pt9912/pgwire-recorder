# Slice slice-v1-abschluss-konfiguration: Konfiguration über Kommandozeile und Umgebung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0027](../../adr/0027-yaml-bibliothek.md), [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md)

**Berührte Spec-Stellen:** `LH-FA-01.a` · `LH-FA-17.a` · `LH-FA-03.b` · `LH-FA-14.a` · `SPEC-007` · `SPEC-008` · `SPEC-012` · `SPEC-014` · `SPEC-020` · `SPEC-034` · `SPEC-046` · `ARC-005` · `ARC-013`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Jede Option von `record` und `replay` außer `--output` und `--force` ist per Kommandozeile und Umgebungsvariable setzbar, über einen allgemeinen Leser in der Priorität Kommandozeile vor Umgebungsvariable vor Default, der jede gesetzte Umgebungsvariable prüft; die Hilfe geht jeder dieser Prüfungen vor, und das Architektur-Gate lässt die YAML-Bibliothek im CLI-Adapter zu und hält sie aus dem PGWire-Adapter.

**Übernimmt:** `slice-v1-abschluss-betrieb` — dessen Teil *Konfiguration* (DoD-Punkt 2 ohne
`--output` und `--force`, DoD-Punkt 3; Entscheidung des Nutzers vom 2026-10-08, Schnitt nach
F-345; dort §7 `Gegenstand:`). Davon gab dieser Slice am 2026-10-08 nach dem Schnitt aus §4
(*zu groß*) die Konfigurationsdatei an `slice-v1-abschluss-konfigurationsdatei` ab (dort §1,
*Übernimmt*, mit der Kennung dieses Slice); hier bleibt, was ohne Datei lieferbar ist. Im
Einzelnen, je mit dem ursprünglichen Geber:

- **Aus `slice-replay-semantik-mismatch`:** die Prüfung jeder gesetzten Umgebungsvariable einer
  Option des Kommandos, auch wenn die Kommandozeile vorgeht (`PGR-E2001`, `LH-FA-17.a`), im
  allgemeinen Leser für alle Optionen; die Hilfe vor jeder Prüfung von Optionen und
  Umgebungsvariablen bei `record` und `replay` (`LH-FA-01.a`), samt `--` als Ende der
  Optionen. Abgegeben: der Schlüssel `fail_on_unconsumed` und die Hilfe für `config show` und
  `--config`.
- **Aus `slice-replay-semantik-meldungscodes`:** abgegeben, der Schlüssel `log_level` gehört
  zur Datei.
- **Aus `slice-v1-abschluss-herunterfahren`** (dort §1, Abgrenzung): dass der allgemeine Leser
  auch `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT` liest. Abgegeben: der Schlüssel der Frist in der
  Konfigurationsdatei.

**Optionen, die der Stand kennt.** Der allgemeine Leser ist die einzige Stelle, an der eine
Option von `record` und `replay` angemeldet wird; Kommandozeile und Umgebungsvariable folgen
aus dieser einen Anmeldung, und `slice-v1-abschluss-konfigurationsdatei` leitet die Schlüssel
der Datei aus ihr ab. Ein Test läuft über alle angemeldeten Optionen (zwei Quellen,
Priorität), sodass eine später angemeldete Option ihn ohne eigenen Plan-Punkt mitnimmt.
`--output` und `--force` meldet dieser Slice nicht am Leser an: Sie liest bis
`slice-v1-abschluss-schreiben` nur die Kommandozeile, wie bisher; der Nehmer meldet sie an
(dort §1, *Übernimmt*, mit der Kennung dieses Slice). Die Hilfe von `record` nennt die
Ausnahme.

**Ort des Lesens.** Die Datei liest der CLI-Adapter (`ARC-005`) in
`slice-v1-abschluss-konfigurationsdatei`. Dieser Slice liefert dafür nur die Gegenprobe des
Architektur-Gates: Die YAML-Bibliothek ist im CLI-Adapter zugelassen (`ARC-013`,
[ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), ergänzt
[ADR-0027](../../adr/0027-yaml-bibliothek.md)) und im PGWire-Adapter abgelehnt. Sicht und
`tech`-Regel in `.a-check.yml` sind nachgezogen; dass ein Import im CLI-Adapter nur der Datei
dient, prüft kein Gate, das bleibt Review.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Konfigurationsdatei: `--config`, `PGWIRE_RECORDER_CONFIG`, die Standarddatei, die
  Schlüssel `fail_on_unconsumed`, `log_level` und der Frist, benannte Verbindungen
  (`connections`, `--upstream <Name>`, `sslmode`), `${VAR}` und `$${VAR}`, `config show`, die
  Hilfe für `config show` und `--config`, `PGR-E2004` bis `PGR-E2006` —
  `slice-v1-abschluss-konfigurationsdatei` (dort §1, *Übernimmt*). Grund: Schnitt nach §4,
  die Datei-Hälfte ist auf 900 bis 1300 Zeilen geschätzt und hat zwölf offene Randformen
  (dort §6).
- `--output`, `--force` und die strengen Werte für `--force` — `slice-v1-abschluss-schreiben`;
  er folgt der Konfigurationsdatei und liest `--force` über den allgemeinen Leser von hier.
- Die Optionen von `play` — `slice-v1-abschluss-einspielen` (dort §1: `--fail-on-unconsumed`
  und `--log-level` bei `play`); er setzt auf dem allgemeinen Leser dieses Slice auf.
- Signale und die Frist selbst — `slice-v1-abschluss-herunterfahren`; hier nur ihre
  Quellen Kommandozeile und Umgebungsvariable.
- Code im Kern, im PGWire-Adapter, in den Driven-Adaptern und im Bootstrap —
  Schicht-Abgrenzung: Der Slice ändert den CLI-Adapter; dazu die Gegenprobe der `tech`-Regel
  des Architektur-Gates (Gate-Werkzeug, keine Schicht).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md): Die Gegenprobe des Architektur-Gates (`make a-check-negativ`) nimmt beide
      YAML-Modulpfade im CLI-Adapter an und lehnt jeden im PGWire-Adapter ab (`ARC-013`);
      Kopfkommentar, Fragment und Zeile in `harness/README.md` §Sensors sagen nicht mehr zu,
      als die Fälle prüfen (Gegenprobe).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Für jede Option von `record` und `replay` außer `--output` und
      `--force`, auch `--shutdown-timeout`, gilt die Priorität Kommandozeile vor
      Umgebungsvariable vor Default (`SPEC-007`, `SPEC-008`); ein allgemeiner Leser prüft
      jede gesetzte Umgebungsvariable einer Option des Kommandos, auch wenn die Kommandozeile
      vorgeht (`PGR-E2001`), und ersetzt die Einzel-Leser von `--fail-on-unconsumed`,
      `--log-level` und `--shutdown-timeout`, deren Tests unverändert grün bleiben (Test).
- [ ] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung): Die Hilfe von `record` und `replay` geht jeder Prüfung von Optionen
      und Umgebungsvariablen am allgemeinen Leser vor, auch einer ungültigen
      Umgebungsvariable, und nennt die Umgebungsvariablen und ihre Priorität, bei `record`
      mit der Ausnahme `--output` und `--force`; `--` beendet die Optionen (Test). Beleg in §7
      für alle drei Punkte: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/architecture.md`, `spec/spezifikation.md`, `.a-check.yml` | erledigt (Architect, vor dem Code, 2026-10-08) | Sicht (`ARC-013`, §2, §6) und `tech`-Regel: YAML-Bibliothek auch im CLI-Adapter ([ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md)); `LH-FA-17.a` mit den Entscheidungen des Nutzers aus §6 |
| `tools/arch/a-check-negativ.sh` | update (geliefert, 8fad493) | drei Fälle: CLI-Adapter importiert `go.yaml.in/yaml/v3` und `gopkg.in/yaml.v3` → a-check meldet nichts; PGWire-Adapter (`internal/adapters/driving/pgwire`) importiert `go.yaml.in/yaml/v3` → `tech-leak`, und ebenso `gopkg.in/yaml.v3` (je Modulpfad eine Regel, je Regel ein Fall). Die beiden letzten halten die Erlaubnis auf den CLI-Adapter statt auf alle Driving-Adapter. Kopfkommentar (Zahl der Fälle, *YAML nur im Recording-Adapter*) und Schlusszeile auf *Recording- und CLI-Adapter*; Beschreibung von `a-check-negativ` in `harness/mk/arch-negativ.mk` und Zeile in `harness/README.md` §Sensors nachziehen (`AGENTS.md` §3.11). Mutation: eine Regel auf `internal/adapters/driving` weiten → Fall ihres Modulpfads im PGWire-Adapter rot; CLI-Adapter aus einer Regel nehmen → Fall des CLI-Adapters rot |
| `internal/adapters/driving/cli` | update (geliefert, a968790) | allgemeiner Leser, an dem jede Option außer `--output` und `--force` einmal angemeldet wird, für Kommandozeile und Umgebungsvariable; er prüft jeden gesetzten Wert und ersetzt die Einzel-Leser von `--fail-on-unconsumed`, `--log-level` und `--shutdown-timeout`; Priorität; Hilfe von `record` und `replay` nennt die Umgebungsvariablen. Die Hilfe vor jeder Prüfung trägt `Parse` seit `slice-replay-semantik-mismatch`: Es prüft die Hilfe-Angabe vor jedem Kommando, also auch vor dem allgemeinen Leser |
| `internal/adapters/driving/cli` (Unit-Tests) | update (geliefert, a968790) | `leser_test.go` über alle angemeldeten Optionen, Reihenfolge, fremde Umgebung, Hilfe nach `LH-FA-17.a` und `LH-FA-01.a`; die vorhandenen Tests der drei Einzel-Leser in `cli_test.go`, `frist_test.go` und `internal/bootstrap` bleiben unverändert und grün (Risiko in §6) |
| `docs/user/benutzerhandbuch.md` | kein update | §5 *Einstellungen* beschreibt Optionen, Umgebungsvariablen und Priorität schon im Zielstand; die Datei beschreibt `slice-v1-abschluss-konfigurationsdatei` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-herunterfahren` liegt in `done/` (der allgemeine Leser
übernimmt dessen Umgebungsvariable). Schritt 3 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)
(Entscheidung des Nutzers vom 2026-10-08). Vor dem ersten Code-Commit prüft der Architect
die Randformen aus §6 und entscheidet die offenen (`AGENTS.md` §3.12).

**Erster Code-Commit:** Die Bedingungen sind seit 2026-10-08 erfüllt —
[ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) ist `Accepted`, Sicht und `tech`-Regel sind nachgezogen, und jede
Randform in §6 ist in `LH-FA-17.a` entschieden. Der erste Commit des Implementers ist die
Gegenprobe des Architektur-Gates (§3).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar, oder eine Änderung im Upstream-Adapter oder im Kern wird nötig.
  Schnitt dann: *Leser* (DoD-Punkt 1 und die Hilfe aus DoD-Punkt 2, ohne Datei) und
  *Datei* (DoD-Punkt 3 und `config show`); der zweite setzt den ersten voraus.
  **Eingetreten am 2026-10-08:** Der Implementer hielt nach der Gegenprobe (8fad493) und dem
  allgemeinen Leser (a968790) an; er schätzte die Datei-Hälfte auf 900 bis 1300 Zeilen, nicht
  in einer Review-Sitzung prüfbar, und gab zwölf Randformen der Datei und die Frage nach dem
  Ort von `PGR-E2004` bis `PGR-E2006` an den Architect zurück. Der Schnitt wurde ohne
  `git mv` umgesetzt: Dieser Slice bleibt in `in-progress/` als *Leser* mit dem Gelieferten,
  die *Datei* ist `slice-v1-abschluss-konfigurationsdatei` in `next/`. Die Hilfe aus
  DoD-Punkt 2 blieb, soweit sie `record` und `replay` betrifft; die für `config show` und
  `--config` ging mit der Datei, weil es beide ohne Datei nicht gibt.
- `in-progress` → `open` (blockiert — Carveout?): keine Bedingung mehr. Die bisherige — die
  Form der Datei nach [ADR-0014](../../adr/0014-konfigurationsdatei.md) trägt einen Schlüssel nicht, den eine Option braucht — betrifft
  nur die Datei und steht in §4 von `slice-v1-abschluss-konfigurationsdatei`.

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

**Randformen** (`AGENTS.md` §3.12) — je Randform, wo sie entschieden ist; vom Architect am
2026-10-08 vor dem Code geprüft. Keine ist offen. Die Randformen der Datei (Wahl und Form,
Verbindungen und Platzhalter, `config show`) gingen mit dem Schnitt (§4) an
`slice-v1-abschluss-konfigurationsdatei` (dort §6).

*Umgebung*

- **Gesetzte, aber leere Umgebungsvariable einer Option** — gilt als nicht gesetzt;
  `LH-FA-17.a`.
- **Umgebungsvariable mit Präfix ohne passende Option** und die einer fremden Option —
  unbeachtet; `LH-FA-17.a`.
- **Ungültige Umgebungsvariable neben gesetzter Kommandozeile** — `PGR-E2001`; `LH-FA-17.a`.
- **Groß- und Kleinschreibung im Namen einer Umgebungsvariable** — folgt dem Betriebssystem;
  unter Windows unterscheidet es nicht. Akzeptiertes Negativ: keine Regel, kein Test, weil
  das Produkt die Umgebung nicht selbst liest, sondern über das Betriebssystem.
- **`PGWIRE_RECORDER_OUTPUT` und `PGWIRE_RECORDER_FORCE`** — bis
  `slice-v1-abschluss-schreiben` nicht gelesen, auch mit ungültigem Wert (§1); die Hilfe von
  `record` nennt die Ausnahme.

*Hilfe und Fehler*

- **Hilfe nach einer ungültigen Option oder neben einer ungültigen Umgebungsvariable** — die
  Hilfe geht vor; `LH-FA-01.a`.
- **Reihenfolge bei mehreren Fehlern** — Abbruch beim ersten Fehler; Reihenfolge
  Kommandozeile, Umgebungsvariablen nach der Tabelle, zuletzt Pflichtoptionen; Entscheidung
  des Nutzers vom 2026-10-08, `LH-FA-17.a`. Die Stellen der Datei, der Kombinationen, von
  `--upstream` und `PGR-E2005` in dieser Reihenfolge liefert
  `slice-v1-abschluss-konfigurationsdatei`.

**Risiken:**

- Der allgemeine Leser ersetzt drei Einzel-Leser; eine Prüfung eines Einzel-Lesers kann
  dabei wegfallen, statt in den allgemeinen überzugehen
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). Gegenmittel: Die vorhandenen Tests der
  drei Einzel-Leser bleiben unverändert und grün (§3) — **Ausgang:** offen bis Closure.
- [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) wird nicht angenommen; dann liest der CLI-Adapter keine
  YAML-Datei, und der Slice geht nach `open` zurück (§4, blockiert), bis eine andere
  Entscheidung den Ort des Lesens trägt — **Ausgang:** entfallen, die ADR ist seit
  2026-10-08 `Accepted`.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Belege zur DoD (Implementer):** siehe *Belege des Implementers* unten.

### Belege des Implementers

Weg der Mutanten: Gegenprobe des Architektur-Gates je Mutant in einer frischen Kopie des
Arbeitsbaums (`mktemp -d`, `tar` ohne `.git`), dort `.a-check.yml` mit `sed` geändert und
`tools/arch/a-check-negativ.sh` aus der Kopie gefahren; der Arbeitsbaum blieb unberührt.

**Gegenprobe des Architektur-Gates** (`make a-check-negativ`, zehn Fälle):

| Zusage | Mutation an `.a-check.yml` | roter Fall |
|---|---|---|
| `go.yaml.in/yaml` im CLI-Adapter zugelassen | CLI-Adapter aus der Regel `go.yaml.in/yaml` genommen | `cli-yaml` (a-check: `tech-leak`) |
| `gopkg.in/yaml` im CLI-Adapter zugelassen | CLI-Adapter aus der Regel `gopkg.in/yaml` genommen | `cli-yaml` |
| beide zugleich | CLI-Adapter aus beiden Regeln genommen | `cli-yaml` |
| `go.yaml.in/yaml` im PGWire-Adapter abgelehnt | Regel `go.yaml.in/yaml` auf `internal/adapters/driving` geweitet | `pgwire-yaml` |
| `gopkg.in/yaml` im PGWire-Adapter abgelehnt | Regel `gopkg.in/yaml` auf `internal/adapters/driving` geweitet | `pgwire-gopkg-yaml` |
| beide zugleich | beide Regeln geweitet | `pgwire-yaml`, `pgwire-gopkg-yaml` |

Ungeändert: alle zehn Fälle grün, Exit 0.

**Allgemeiner Leser: Kommandozeile, Umgebung, Priorität** (DoD-Punkt 1, ohne Datei). Weg der
Mutanten: je Mutant eine frische Kopie des Arbeitsbaums (`mktemp -d`, `tar` ohne `.git`), die
Änderung an `internal/adapters/driving/cli/cli.go` mit einem Skript, dann `make test` in der
Kopie; ein neuer Pfad je Mutant, darum keine gleiche mtime am selben Pfad. Die vorhandenen
Tests der drei Einzel-Leser (`cli_test.go`, `frist_test.go`, die Tests in
`internal/bootstrap`) sind unverändert und grün; `export_test.go` ist nur ergänzt.

| Zusage | Mutation in `lies` / `optionen` | rote Tests |
|---|---|---|
| Kommandozeile vor Umgebungsvariable | Umgebung vor Kommandozeile übernommen | `TestLeserAlleOptionen`, `TestParseFailOnUnconsumedUmgebung`, `TestParseLogLevelUmgebung`, `TestParseShutdownTimeoutUmgebung`, `TestRunLogLevel` |
| Umgebungsvariable setzt jede Option | Umgebung nicht gelesen | `TestLeserAlleOptionen`, `TestLeserReihenfolge`, die sechs `…Umgebung…`-Tests der Einzel-Leser, `TestRunLogLevel`, `TestRunStartfehlerJeStufe` |
| gesetzte Umgebungsvariable geprüft, auch wenn die Kommandozeile vorgeht | Prüfung nur ohne Kommandozeile | `TestLeserAlleOptionen`, `TestParseFailOnUnconsumedUmgebungNebenOption`, `TestParseLogLevelUmgebungUngueltig`, `TestParseShutdownTimeoutUmgebungUngueltig`, `TestRunStartfehlerJeStufe` |
| leere Umgebungsvariable gilt als nicht gesetzt | leere Variable geprüft und übernommen | 30 Tests, darunter `TestLeserAlleOptionen`, `TestLeserFremdeUmgebung`, `TestParseRecord` |
| Umgebungsvariablen in der Reihenfolge der Tabelle geprüft | Schleife rückwärts | `TestLeserReihenfolge` |
| Kommandozeile vor den Umgebungsvariablen geprüft | Umgebung vor `fs.Parse` geprüft | `TestLeserReihenfolge` |
| Pflichtoption ohne Wert ist `PGR-E2001` | Pflicht-Prüfung entfernt | `TestLeserAlleOptionen`, `TestLeserReihenfolge`, `TestParseReplay` |
| ohne Quelle gilt der Standardwert | Standardwert leer | `TestLeserAlleOptionen`, `TestParseLogLevel`, `TestParseRecord`, `TestParseReplay`, `TestParseShutdownTimeout`, `TestParseLogLevelUmgebung`, `TestParseShutdownTimeoutUmgebung` |
| Name der Umgebungsvariable mit `_` statt `-` | `-` bleibt stehen | `TestLeserOptionen`, `TestLeserReihenfolge`, die sechs `…Umgebung…`-Tests, `TestRunLogLevel`, `TestRunStartfehlerJeStufe` |
| boolesche Option ohne Wert ist `true` | `schalter` aus | `TestParseFailOnUnconsumed`, `TestParseFailOnUnconsumedUmgebung`, `TestParseFailOnUnconsumedUmgebungNebenOption`, `TestRunLogLevel` |
| Hilfe von `replay` nennt die Umgebungsvariablen | Satz geändert | `TestLeserHilfe` |
| `--output` liest nur die Kommandozeile (Hilfe von `record`) | `--output` aus `PGWIRE_RECORDER_OUTPUT` | `TestLeserHilfe` |

Eine Umgebungsvariable mit Präfix ohne passende Option bleibt unbeachtet: Der Leser fragt nur
die Namen seiner Optionen ab; `TestLeserFremdeUmgebung` belegt es für fünf Namen. Ein Mutant,
der sie liest, müsste eine neue Abfrage einführen; keiner gefahren.
- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `3864e44` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — der allgemeine Leser ersetzt
  Einzel-Leser; ein Risiko in §6.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (12×, verkörpert in `AGENTS.md`
  §3.12) — darum nennt §6 die Randformen, und der Architect prüft sie vor dem Code.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (14×, `AGENTS.md` §3.10) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (19×, §3.11) — je Zusage eine
  Mutation, Beleg in §7.
- `BEO-REPO/plan-folgt-korrektur-nicht` (17×, §3.9) — dieser Plan entsteht aus einem
  Schnitt; jede Korrektur zieht §1, §3 und §6 im selben Commit nach.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — jede Übernahme oben steht
  mit der Kennung des Gebers; die Geber zeigen im selben Commit hierher.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu; keine neue Lücke vor
dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
