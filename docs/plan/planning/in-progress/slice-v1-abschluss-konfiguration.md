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

**Ziel:** Jede Option von `record` und `replay` ist per Kommandozeile und Umgebungsvariable setzbar, über einen allgemeinen Leser in der Priorität Kommandozeile vor Umgebungsvariable vor Default, der jede gesetzte Umgebungsvariable prüft; die Hilfe geht jeder dieser Prüfungen vor, und das Architektur-Gate lässt die YAML-Bibliothek im CLI-Adapter zu und hält sie aus dem PGWire-Adapter.

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
- **Aus `slice-v1-abschluss-schreiben`** (dort §1, *Abgegeben*; Entscheidung des Architect vom
  2026-10-08 zu F-496): `--output` und `--force` am allgemeinen Leser, damit
  `PGWIRE_RECORDER_OUTPUT` und `PGWIRE_RECORDER_FORCE` gelesen und geprüft werden
  (`LH-FA-17.a`), und die strengen Werte für `--force` (`true`, `false`), die
  `slice-v1-abschluss-schreiben` aus `slice-replay-semantik-mismatch` übernommen hatte.

**Optionen, die der Stand kennt.** Der allgemeine Leser ist die einzige Stelle, an der eine
Option von `record` und `replay` angemeldet wird; Kommandozeile und Umgebungsvariable folgen
aus dieser einen Anmeldung, und `slice-v1-abschluss-konfigurationsdatei` leitet die Schlüssel
der Datei aus ihr ab. Ein Test läuft über alle angemeldeten Optionen (zwei Quellen,
Priorität, leerer Wert, Default gegen die Optionstabelle), sodass eine später angemeldete
Option ihn mitnimmt; er bricht ab, bis ihr Default aus der Tabelle im Test steht. Ein
zweiter Test hält die *einzige Stelle*: Jede Option der Tabelle, die nicht am Leser
angemeldet ist, ist auf der Kommandozeile unbekannt.
Auch `--output` (Pflicht) und `--force` (Wahrheitswert, Default `false`) sind dort
angemeldet, in der Reihenfolge der Tabelle in `LH-FA-17.a` nach `--upstream`; `parseRecord`
meldet keine Option mehr am FlagSet an.

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
  die Datei-Hälfte ist auf 900 bis 1300 Zeilen geschätzt und hatte zwölf offene Randformen
  (dort §6, vor dem Code entschieden).
- Das atomare Schreiben und das Verhalten bei vorhandenem `--output` —
  `slice-v1-abschluss-schreiben`; er folgt der Konfigurationsdatei und nutzt `--output` und
  `--force`, wie der allgemeine Leser von hier sie liefert.
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

- [x] [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md): Die Gegenprobe des Architektur-Gates (`make a-check-negativ`) nimmt beide
      YAML-Modulpfade im CLI-Adapter an und lehnt jeden im PGWire-Adapter ab (`ARC-013`);
      Kopfkommentar, Fragment und Zeile in `harness/README.md` §Sensors sagen nicht mehr zu,
      als die Fälle prüfen (Gegenprobe).
- [x] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Für jede Option von `record` und `replay`, auch `--output`,
      `--force` und `--shutdown-timeout`, gilt die Priorität Kommandozeile vor
      Umgebungsvariable vor Default (`SPEC-007`, `SPEC-008`); ein allgemeiner Leser prüft
      jede gesetzte Umgebungsvariable einer Option des Kommandos, auch wenn die Kommandozeile
      vorgeht (`PGR-E2001`), und ersetzt die Einzel-Leser von `--fail-on-unconsumed`,
      `--log-level` und `--shutdown-timeout`, deren Tests unverändert grün bleiben (Test).
- [x] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung): Die Hilfe von `record` und `replay` geht jeder Prüfung von Optionen
      und Umgebungsvariablen am allgemeinen Leser vor, auch einer ungültigen
      Umgebungsvariable, und nennt die Umgebungsvariablen und ihre Priorität; `--` beendet die Optionen (Test). Beleg in §7
      für alle drei Punkte: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/architecture.md`, `spec/spezifikation.md`, `.a-check.yml` | erledigt (Architect, vor dem Code, 2026-10-08) | Sicht (`ARC-013`, §2, §6) und `tech`-Regel: YAML-Bibliothek auch im CLI-Adapter ([ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md)); `LH-FA-17.a` mit den Entscheidungen des Nutzers aus §6 |
| `tools/arch/a-check-negativ.sh` | update (geliefert, 8fad493) | drei Fälle: CLI-Adapter importiert `go.yaml.in/yaml/v3` und `gopkg.in/yaml.v3` → a-check meldet nichts; PGWire-Adapter (`internal/adapters/driving/pgwire`) importiert `go.yaml.in/yaml/v3` → `tech-leak`, und ebenso `gopkg.in/yaml.v3` (je Modulpfad eine Regel, je Regel ein Fall). Die beiden letzten halten die Erlaubnis auf den CLI-Adapter statt auf alle Driving-Adapter. Nach dem Review (F-499) zwei weitere: Domain Model importiert `gopkg.in/yaml.v3`, Postgres-Adapter importiert `go.yaml.in/yaml/v3`, sodass jeder Modulpfad an jedem der drei genannten Orte abgelehnt wird. Kopfkommentar (Zahl der Fälle, *YAML nur im Recording-Adapter*) und Schlusszeile auf *Recording- und CLI-Adapter*; Beschreibung von `a-check-negativ` in `harness/mk/arch-negativ.mk` und Zeile in `harness/README.md` §Sensors nachziehen (`AGENTS.md` §3.11). Mutation: eine Regel auf `internal/adapters/driving` weiten → Fall ihres Modulpfads im PGWire-Adapter rot; CLI-Adapter aus einer Regel nehmen → Fall des CLI-Adapters rot |
| `internal/adapters/driving/cli` | update (geliefert, a968790; Nacharbeit 960f239) | allgemeiner Leser, an dem jede Option einmal angemeldet wird, auch `--output` und `--force` mit den strengen Werten für `--force` (geliefert in 960f239, Entscheidung zu F-496), für Kommandozeile und Umgebungsvariable; er prüft jeden gesetzten Wert und ersetzt die Einzel-Leser von `--fail-on-unconsumed`, `--log-level` und `--shutdown-timeout`; Priorität; Hilfe von `record` und `replay` nennt die Umgebungsvariablen. Die Hilfe vor jeder Prüfung trägt `Parse` seit `slice-replay-semantik-mismatch`: Es prüft die Hilfe-Angabe vor jedem Kommando, also auch vor dem allgemeinen Leser |
| `internal/adapters/driving/cli` (Unit-Tests) | update (geliefert, a968790; Nacharbeit 960f239) | `leser_test.go` über alle angemeldeten Optionen, Reihenfolge, fremde Umgebung, Hilfe nach `LH-FA-17.a` und `LH-FA-01.a`; mit 960f239 dazu `TestLeserNurAngemeldete` (F-498), `TestLeserLeererWert` (F-497) und der Default aus der Optionstabelle in `tabelle()` (F-505); die vorhandenen Tests der drei Einzel-Leser in `cli_test.go`, `frist_test.go` und `internal/bootstrap` bleiben unverändert und grün (Risiko in §6) |
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
  die *Datei* ist `slice-v1-abschluss-konfigurationsdatei` in `next/`. Der Architect bestätigt
  am 2026-10-08 (F-503), dass das trägt: Die Pflichten der Rückführung (Bedingung vorab, Grund
  nachgetragen) stehen hier, der Schnitt geschah im selben Zug, und der Zustand stimmt für das,
  was die Datei jetzt beschreibt. Der Weg über `next/` und zurück ergäbe zwei reine
  `git mv` mit demselben Endstand. Akzeptiertes Negativ: Den Schnitt zeigt nicht die
  Verzeichnis-Historie, sondern dieser Absatz und `fa4a5f1`. Die Nummern in der Bedingung
  oben zählen die DoD vor dem Schnitt. Nach dem Schnitt ist DoD-Punkt 1 die Gegenprobe des
  Architektur-Gates, Punkt 2 der allgemeine Leser (der frühere Punkt 1 ohne Datei) und
  Punkt 3 die Hilfe von `record` und `replay` (aus dem früheren Punkt 2); die Hilfe für
  `config show` und `--config` ging mit der Datei, weil es beide ohne Datei nicht gibt.
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
- **`PGWIRE_RECORDER_OUTPUT` und `PGWIRE_RECORDER_FORCE`** — gelesen und geprüft wie jede
  Umgebungsvariable einer Option (`LH-FA-17.a`, Optionstabelle); `PGWIRE_RECORDER_FORCE=ja`
  ist `PGR-E2001`. Entscheidung des Architect vom 2026-10-08 zu F-496: Die Spezifikation galt,
  der Plan hatte eine Ausnahme behauptet, die sie nicht trägt; beide Optionen kommen an den
  Leser (§1).

*Kommandozeile*

- **Leerer Wert auf der Kommandozeile** (`--listen=`), auch neben gesetzter Umgebungsvariable
  derselben Option — gesetzt und ungültig, `PGR-E2001`, geprüft mit der Kommandozeile, also
  vor den Umgebungsvariablen; nur die leere Umgebungsvariable gilt als nicht gesetzt.
  `LH-FA-17.a`, Entscheidung des Architect vom 2026-10-08 zu F-497 (dieselbe Strenge wie
  beim leeren Wert boolescher Optionen und der Dauer). Testfall: `record --listen=
  --upstream h:1 --output r.yaml` mit `PGWIRE_RECORDER_LISTEN=127.0.0.1:1` und
  `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT=x` endet mit `PGR-E2001`, die Meldung nennt `listen`,
  weder `Pflichtoption` noch die Variable; dazu je Option der leere Wert im Test über alle
  angemeldeten Optionen. Mutant L6 (leerer Wert gilt als nicht gesetzt) wird damit rot.

*Hilfe und Fehler*

- **Hilfe nach einer ungültigen Option oder neben einer ungültigen Umgebungsvariable** — die
  Hilfe geht vor; `LH-FA-01.a`.
- **Wert in der Meldung** — eine Meldung zu einer Umgebungsvariable nennt ihren Namen, nie
  ihren Wert; eine Meldung zur Kommandozeile darf den Wert nennen
  (`invalid value "trace" for flag -log-level`), weil der Aufruf ihn selbst enthält und die
  Kommandozeile kein Passwort trägt. `LH-FA-17.a` *Fehler*, Entscheidung des Architect vom
  2026-10-08 zu V-116; das Lastenheft (`LH-FA-17`) verlangt nur, dass die Anzeige der
  Konfiguration keine Geheimnisse zeigt. Kein Code-Nachtrag.
- **Reihenfolge bei mehreren Fehlern** — Abbruch beim ersten Fehler; Reihenfolge
  Kommandozeile, Umgebungsvariablen nach der Tabelle, zuletzt Pflichtoptionen; Entscheidung
  des Nutzers vom 2026-10-08, `LH-FA-17.a`. Die Stellen der Datei, der Kombinationen, von
  `--upstream` und `PGR-E2005` in dieser Reihenfolge liefert
  `slice-v1-abschluss-konfigurationsdatei`.

**Risiken:**

- Der allgemeine Leser ersetzt drei Einzel-Leser; eine Prüfung eines Einzel-Lesers kann
  dabei wegfallen, statt in den allgemeinen überzugehen
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). Gegenmittel: Die vorhandenen Tests der
  drei Einzel-Leser bleiben unverändert und grün (§3) — **Ausgang:** entfallen. Die Tests
  der drei Einzel-Leser in `cli_test.go`, `frist_test.go` und `internal/bootstrap` sind im
  Diff `c060323..960f239` unverändert und grün; der Vergleich alter gegen neuer Leser fand
  keine weggefallene Prüfung (Review F-504: 43 von 46 Fällen gleich, die Abweichungen sind
  die Reihenfolge der Fehler nach `LH-FA-17.a` und `PGWIRE_RECORDER_LISTEN`; Verifikation
  Abschnitt 3: 26 von 29 Fällen gleich, abweichend ein Zufallsname und die strengen Werte
  für `--force` nach F-496; alle Abweichungen gewollt).
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

- **Belege zur DoD (Implementer):** siehe *Belege des Implementers* unten. DoD-Punkt 4
  (`make gates` grün, V-113): Lauf der Verifikation am Stand `960f239`, Exit 0
  (`docs/reviews/2026-10-08-verifikation-slice-v1-abschluss-konfiguration.md`, Abschnitt 1,
  Punkt 4); Lauf der Closure am Stand `647f0ff` (nach 960f239 nur Plan und Spezifikation
  geändert), Exit 0, darin `a-check-negativ: gruen`, `d-check: 371 Datei(en) geprüft, 0
  Befund(e)`, `baseline-verify: v6.16.0 OK`, `run-integration-tests: gruen`,
  `kopf-check-gegenprobe: gruen`, `abdeckung-gegenprobe: gruen`, `lint-gegenprobe: gruen`.
  Die Nacharbeit zum Review ist am Diff geprüft: F-496 bis F-502 und F-505 in `960f239`,
  nachgefahren von der Verifikation (Mutanten L6, L3, L9, F1, H1, A1 rot); V-114 (§3 nach
  `960f239`) und V-113 (diese Zeile) im Closure-Commit; V-116 vom Architect in `647f0ff`.
- **Was hat funktioniert:** Die Prüfung von §6 vor dem Code (`1101e45`, `c060323`) trug:
  Acht offene Randformen entschied der Nutzer vor dem ersten Code-Commit, und keine
  davon entschied der Code anders. Die in §4 vorab benannte Rückführung *zu groß* trat
  ein und schnitt entlang der vorab benannten Linie (*Leser* und *Datei*); der
  Implementer hielt nach zwei Liefer-Commits an und gab die zwölf Randformen der Datei
  zurück, statt sie im Code zu entscheiden (§3.12 hielt für die Datei-Hälfte). Das
  Verhalten der drei ersetzten Einzel-Leser blieb gleich (Review F-504, Verifikation
  Abschnitt 3), die Gegenprobe des Architektur-Gates fängt Weiten und Wegnehmen je
  Modulpfad (A1 bis A4 rot), und nach dem Schnitt passte der Diff in eine Review-Sitzung.
- **Was ging anders als geplant:**
  1. Der Slice war zu groß. Mit vier Übernahmen (aus `slice-v1-abschluss-betrieb`, darin
     aus `slice-replay-semantik-mismatch`, `slice-replay-semantik-meldungscodes` und
     `slice-v1-abschluss-herunterfahren`) schätzte der Implementer die Datei-Hälfte auf 900
     bis 1300 Zeilen; Schnitt in `fa4a5f1`, die Datei ging an
     `slice-v1-abschluss-konfigurationsdatei`. Der Schnitt blieb ohne `git mv` in
     `in-progress/` (F-503, V-117; Auftrag des Nutzers laut Drift-Log der Roadmap, vom
     Architect in §4 bestätigt): Die State Machine kennt für *zu groß* nur
     `in_progress → next`, und `MR-000` erklärt für die Lifecycle-Regeln keine Adaption.
  2. Randformen außerhalb der Spezifikation entschieden: Dass `record` die Variablen von
     `--output` und `--force` nicht liest, stand im Code-Commit `a968790` nur im Plan,
     gegen die Optionstabelle von `LH-FA-17.a` (F-496); den leeren Wert auf der
     Kommandozeile neben gesetzter Variable entschied der Code still (F-497). Beide
     entschied der Architect in `096e51e`, umgesetzt in `960f239`. Die Lesart von „nie einen
     Wert“ für Meldungen zur Kommandozeile fand erst die Verifikation (V-116, Architect
     `647f0ff`).
  3. Zusagen ohne fangenden Test: „einzige Stelle der Anmeldung“ (F-498, L3 grün),
     *unerwartetes Argument* vor der Umgebung (F-500, L8 grün), der Default gegen die
     eigene Anmeldung statt gegen die Tabelle (F-505, L9 im Allgemein-Test grün). Der
     Sensor-Text der Gegenprobe sagte mehr zu als die Fälle (F-499, A1 grün); die
     Schlusszeile von `tools/arch/a-check-negativ.sh` tut es weiter (V-115, siehe
     *Folge-Slices*).
  4. Der Plan folgte der Nacharbeit nicht ganz: §4 zählte die DoD vor dem Schnitt (F-502),
     §3 nannte `--output` und `--force` nach `960f239` noch „offen“ (V-114), §7 trug keinen
     Beleg für `make gates` (V-113). Der Kommentar an `wahrheitswert` beschrieb eine
     entfallene Rolle (F-501).
  - **Summary-Zeilen:** Review
    `docs/reviews/2026-10-08-review-slice-v1-abschluss-konfiguration.md`: „0 HIGH · 3
    MEDIUM · 4 LOW · 3 INFO (F-496 Ausnahme `PGWIRE_RECORDER_OUTPUT`/`_FORCE` nur im Plan
    entschieden, im Code-Commit eingeführt, gegen die Optionstabelle von `LH-FA-17.a`;
    F-497 leerer Wert auf der Kommandozeile neben gesetzter Variable vom Code entschieden,
    Mutant grün; F-498 „einzige Stelle der Anmeldung“ ungeprüft, Option am FlagSet vorbei
    fällt nicht auf; F-499 Sensor-Text sagt YAML-Ablehnung im Domain Model und
    Postgres-Adapter weiter zu als die Fälle; F-500 Reihenfolge *unerwartetes Argument* vor
    Umgebung ungeprüft; F-501 Kommentar an `wahrheitswert` beschreibt entfallene Rolle;
    F-502 §4 zählt die DoD vor dem Schnitt; F-503 bis F-505 Lifecycle ohne `git mv`,
    gleiches Verhalten belegt, Standardwert im Allgemein-Test selbstbezüglich).
    Wiederkehrende Klassen: `BEO-REPO/spec-randform-erst-im-review-entschieden` (F-496,
    F-497), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-498, F-499),
    `BEO-REPO/plan-folgt-korrektur-nicht` (F-502).“ Verifikation
    `docs/reviews/2026-10-08-verifikation-slice-v1-abschluss-konfiguration.md`: „0 HIGH · 0
    MEDIUM · 2 LOW · 3 INFO (V-113 Gate-Beleg fehlt in §7; V-114 §3 nennt
    `--output`/`--force` „offen“; V-115 Schlusszeile der Gegenprobe weiter als die Fälle;
    V-116 Wert in der Meldung der Kommandozeile, Lesart von `LH-FA-17.a` *Fehler* offen;
    V-117 Verbleib in `in-progress/` nach dem Schnitt ohne `MR`).“ Verdikt: DoD-Liefer-Punkte
    1 bis 3 und `make gates` bestätigt, kein blockierender Befund. Ausgänge: V-113 und
    V-114 in diesem Commit, V-116 in `647f0ff`, V-115 und V-117 unter
    *Beobachtungs-Register*.
- **Steering-Loop-Eintrag:** Benannte Spec-Lücken, beide in `LH-FA-17.a` geschlossen:
  (1) Ein leerer Wert auf der Kommandozeile neben gesetzter Umgebungsvariable derselben
  Option war nicht entschieden (F-497); jetzt gesetzt und `PGR-E2001`, Architect `096e51e`.
  (2) „nie einen Wert“ in *Fehler* ließ offen, ob es auch für Meldungen zur Kommandozeile
  gilt (V-116); jetzt nur für Datei und Umgebungsvariablen, Architect `647f0ff`. Dazu eine
  neue Beobachtung statt einer Adaption: `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel`
  (V-117, 1×) — eine eingetretene Rückführung *zu groß*, deren Schnitt im selben Zug
  geschieht, bleibt in `in-progress/`. Ob daraus eine Adaption (`MR-<NNN>` in
  `harness/conventions.md`) wird, entscheidet der Nutzer; das Register zählt bis dahin.
  Kein Feld `liegt in`: Mit diesem Slice ist nichts verkörpert; die Spec-Stellen tragen
  keinen Herkunfts-Anker.

  Retirement-Checks: `AGENTS.md` §3.9 (seit welle-walking-skeleton) ist wieder
  aufgetreten (F-502, V-114, V-113); die Regel bleibt. §3.10 (seit welle-extended-query)
  ist wieder aufgetreten (F-500, F-505); die Regel bleibt, der Sensor ist mit
  `slice-harness-mutation` geplant. §3.11 (seit welle-extended-query) ist wieder
  aufgetreten (F-498, F-499, V-115); die Regel bleibt. §3.12 (seit
  slice-harness-randformen-vor-code) ist wieder aufgetreten (F-496 im Plan statt in der
  Spezifikation, F-497 still im Code, V-116 eine Lesart der Spezifikation selbst); die
  Regel bleibt, für die Datei-Hälfte hat sie getragen. §3.13 (seit
  slice-lint-bestand-kern-driven) ist nicht wieder aufgetreten: Jede Übernahme in §1 nennt
  ihren Geber, und die Nehmer der Abgrenzungen nennen diesen Slice (Review, Negativbefund
  *Plan*; Paarungen unten).
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `647f0ff` (Zähler =
  Dateien unter `evidence/`).
  - `BEO-REPO/spec-randform-erst-im-review-entschieden`: **Beleg**, 14× → 15× (F-496,
    F-497; dazu V-116), Stand verkörpert, bleibt.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: **Beleg**, 21× → 22× (F-498,
    F-499; dazu V-115), Stand verkörpert, bleibt.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: **Beleg**, 15× → 16× (F-500, F-505:
    Zusagen des allgemeinen Lesers, deren Mutation der Implementer nicht gefahren hatte),
    Stand verkörpert, bleibt.
  - `BEO-REPO/plan-folgt-korrektur-nicht`: **Beleg**, 17× → 18× (F-502, V-114), Stand
    verkörpert, bleibt.
  - `BEO-REPO/slice-waechst-durch-uebernahmen`: **Beleg**, 1× → 2× (der Schnitt nach §4,
    `fa4a5f1`), offen.
  - `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel`: **neu**, 1× (F-503, V-117), offen.
  - Ohne Beleg: `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (bleibt 4×,
    verkörpert; F-496 und F-497 gab der Implementer nicht zurück, das Review fand sie —
    sie zählen oben; die zwölf Randformen der Datei gab er vor dem Code zurück),
    `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (bleibt 1×; Risiko entfallen, keine
    Prüfung weggefallen), `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (bleibt 3×,
    verkörpert), `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht` (bleibt 3×,
    verkörpert; die Belege in §7 erreichten beide Prüfer, nur der Gate-Lauf fehlte dort),
    `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (bleibt 2×; neuer Vertrag, die
    grünen Mutanten zählen unter `negativtests-…` und `zusage-…`).

  Einmalig und nicht eingetragen: F-501 (Kommentar an einem Typ nach dem Umbau, `AGENTS.md`
  §3.7), V-113 (ein fehlender Beleg-Satz in §7), F-504 (Negativbefund). Mit diesem Slice
  erreicht kein Eintrag die Schwelle 3× neu.
- **Folge-Slices:** keiner aus Risiko oder Register. V-115: Die Schlusszeile
  „YAML nur im Recording- und CLI-Adapter“ stand schon vor dem Slice so („nur im
  Recording-Adapter“); kein offener Harness-Slice nimmt die Gegenprobe des
  Architektur-Gates an, ohne seinen Schnitt zu verlassen (`slice-harness-meldungskatalog-gate`,
  `slice-harness-gate-index-werkzeug-teil`, `slice-harness-mutation` haben je einen anderen
  Gegenstand). Ausgang: Beleg in `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; die
  Zeile bleibt als Bestand stehen, bis der Nutzer eine Adresse entscheidet.
  Aus §1 *Ausdrücklich NICHT*: `slice-v1-abschluss-konfigurationsdatei` (Datei, `config
  show`, `PGR-E2004` bis `PGR-E2006`; als nächster in der Reihe),
  `slice-v1-abschluss-schreiben` (atomares Schreiben, vorhandenes `--output`),
  `slice-v1-abschluss-einspielen` (Optionen von `play` über den allgemeinen Leser); alle
  drei in `next/`. Signale und Frist lieferte `slice-v1-abschluss-herunterfahren` vor diesem
  Slice (`done/`); er ist Geber, kein Folge-Slice.
- **Risiken aus §6:** zwei. *Einzel-Leser* (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`):
  **entfallen**, Begründung in §6. *[ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) nicht angenommen*: **entfallen**, die ADR ist
  `Accepted`. Die Randformen in §6 sind Entscheidungen, keine Risiken.
- **Drei Paarungen:** Anker: Der Steering-Loop-Eintrag trägt kein Feld `liegt in`; nichts
  zu prüfen. Folge-Slice: `grep -n "slice-v1-abschluss-konfiguration"` findet die Kennung
  in §1 jedes Nehmers — `slice-v1-abschluss-konfigurationsdatei` (`next/`, §1
  *Übernimmt*), `slice-v1-abschluss-schreiben` (`next/`, §1 *Abgegeben* und
  *Ausdrücklich NICHT*: Leser hier, atomares Schreiben und vorhandenes `--output` dort im
  Ziel), `slice-v1-abschluss-einspielen` (`next/`, §1 *Übernommen aus
  `slice-v1-abschluss-konfiguration`*); keiner schließt die Sendung in seinem §1
  *Ausdrücklich NICHT* aus, keiner liegt in `done/`. Register:
  `BEO-REPO/spec-randform-erst-im-review-entschieden`,
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`,
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/plan-folgt-korrektur-nicht`,
  `BEO-REPO/slice-waechst-durch-uebernahmen` und
  `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel` tragen
  `evidence/slice-v1-abschluss-konfiguration.md`; die übrigen genannten Einträge bestehen
  als Verzeichnis, jedes mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft
  erneut.

### Belege des Implementers

Stand nach der Nacharbeit zum Review (F-496 bis F-502, F-505). Weg der Mutanten: je Mutant
eine frische Kopie des Arbeitsbaums (`mktemp -d`, `tar` ohne `.git`, ein neuer Pfad je
Mutant, darum keine gleiche mtime am selben Pfad), die Änderung mit einem Skript, das die
Trefferzahl 1 prüft; dann `make test` bzw. `tools/arch/a-check-negativ.sh` in der Kopie. Der
Arbeitsbaum blieb unberührt. Jede Zeile ist am Stand der Nacharbeit gefahren.

**Gegenprobe des Architektur-Gates** (DoD-Punkt 1; `make a-check-negativ`, zwölf Fälle).
Mutationen an `.a-check.yml`; ungeändert sind alle zwölf Fälle grün, Exit 0. Die Fälle 11
(`gopkg.in/yaml.v3` im Domain Model) und 12 (`go.yaml.in/yaml/v3` im Postgres-Adapter) sind
für F-499 neu; damit gilt die Ablehnung je Modulpfad an allen drei Orten, die README,
Fragment und Kopf nennen.

| Zusage | Mutation | roter Fall |
|---|---|---|
| `go.yaml.in/yaml` im CLI-Adapter zugelassen | CLI-Adapter aus dieser Regel genommen | `cli-yaml` (a-check: `tech-leak`) |
| `gopkg.in/yaml` im CLI-Adapter zugelassen | CLI-Adapter aus dieser Regel genommen | `cli-yaml` |
| `go.yaml.in/yaml` im PGWire-Adapter abgelehnt | Regel auf `internal/adapters/driving` geweitet | `pgwire-yaml` |
| `gopkg.in/yaml` im PGWire-Adapter abgelehnt | Regel auf `internal/adapters/driving` geweitet | `pgwire-gopkg-yaml` |
| `go.yaml.in/yaml` im Postgres-Adapter abgelehnt (A1 aus dem Review) | Regel von `driven/recording` auf `internal/adapters/driven` geweitet | `postgres-yamlin` |
| `gopkg.in/yaml` im Postgres-Adapter abgelehnt | Regel von `driven/recording` auf `internal/adapters/driven` geweitet | `postgres-yaml` |

Grün, eingeordnet: Das Domain Model in die `adapter`-Liste einer der beiden Regeln
aufgenommen (je Modulpfad ein Mutant) — **äquivalent**: a-check lehnt den Import im Domain
Model über dessen Rolle `domain` ab, nicht über die `tech`-Regel; die Fälle 6 und 11 bleiben
rot wie erwartet, das Verhalten ist gleich.

**Allgemeiner Leser: Kommandozeile, Umgebung, Priorität** (DoD-Punkt 2). Die Tests der drei
Einzel-Leser (`cli_test.go`, `frist_test.go`, `internal/bootstrap`) sind unverändert und grün;
`export_test.go` ist nur ergänzt.

| Zusage | Mutation in `cli.go` | rote Tests |
|---|---|---|
| Kommandozeile vor Umgebungsvariable | Umgebung vor Kommandozeile übernommen | `TestLeserAlleOptionen`, `TestParseFailOnUnconsumedUmgebung`, `TestParseLogLevelUmgebung`, `TestParseShutdownTimeoutUmgebung`, `TestRunLogLevel` |
| Umgebungsvariable setzt jede Option | Umgebung nicht gelesen | `TestLeserAlleOptionen`, `TestLeserHilfe`, `TestLeserReihenfolge`, die sechs `…Umgebung…`-Tests der Einzel-Leser, `TestRunLogLevel`, `TestRunStartfehlerJeStufe` |
| gesetzte Umgebungsvariable geprüft, auch wenn die Kommandozeile vorgeht | Prüfung nur ohne Kommandozeile | `TestLeserAlleOptionen`, `TestLeserHilfe`, `TestParseFailOnUnconsumedUmgebungNebenOption`, `TestParseLogLevelUmgebungUngueltig`, `TestParseShutdownTimeoutUmgebungUngueltig`, `TestRunStartfehlerJeStufe` |
| leere Umgebungsvariable gilt als nicht gesetzt | leere Variable geprüft und übernommen | 30 Tests, darunter `TestLeserAlleOptionen`, `TestLeserFremdeUmgebung`, `TestParseRecord` |
| leerer Wert auf der Kommandozeile ist gesetzt und `PGR-E2001` (F-497, L6) | leerer Wert in `kommandozeile.Set` als nicht gesetzt übergangen | `TestLeserAlleOptionen`, `TestLeserLeererWert`, `TestParseFailOnUnconsumedWerte`, `TestParseLogLevelWerte`, `TestParseShutdownTimeoutWerte`, `TestRunStartfehlerJeStufe` |
| leerer Wert auch bei einer Option der Art `text` ungültig | Prüfung des leeren Werts in `artText` entfernt | `TestLeserAlleOptionen`, `TestLeserLeererWert` |
| Umgebungsvariablen in der Reihenfolge der Tabelle geprüft | Schleife rückwärts | `TestLeserReihenfolge` |
| Kommandozeile vor den Umgebungsvariablen geprüft | Umgebung vor `fs.Parse` geprüft | `TestLeserReihenfolge`, `TestLeserLeererWert` |
| unerwartetes Argument vor den Umgebungsvariablen (F-500, L8) | Prüfung *unerwartetes Argument* hinter die Umgebung | `TestLeserReihenfolge` |
| einzige Stelle der Anmeldung (F-498, L3) | `parseReplay` meldet `--session-assignment` am FlagSet vorbei an | `TestLeserNurAngemeldete` |
| Pflichtoption ohne Quelle ist `PGR-E2001` | Pflicht-Prüfung entfernt | `TestLeserAlleOptionen`, `TestLeserReihenfolge`, `TestParseFehler`, `TestParseReplay` |
| ohne Quelle gilt der Default der Optionstabelle (F-505, L9) | Standard von `--log-level` bei `record` `warn` | `TestLeserAlleOptionen` und acht Einzeltests, darunter `TestParseLogLevel`, `TestRunRecordLogLevel` |
| ohne Quelle gilt der Standardwert | Standardwert leer | `TestLeserAlleOptionen`, `TestParseLogLevel`, `TestParseRecord`, `TestParseReplay`, `TestParseShutdownTimeout`, `TestParseLogLevelUmgebung`, `TestParseShutdownTimeoutUmgebung` |
| Name der Umgebungsvariable mit `_` statt `-` | `-` bleibt stehen | `TestLeserOptionen`, `TestLeserReihenfolge`, die sechs `…Umgebung…`-Tests, `TestRunLogLevel`, `TestRunStartfehlerJeStufe` |
| boolesche Option ohne Wert ist `true` | `schalter` aus | `TestLeserHilfe`, `TestParseFailOnUnconsumed`, `TestParseFailOnUnconsumedUmgebung`, `TestParseFailOnUnconsumedUmgebungNebenOption`, `TestParseRecord`, `TestRunLogLevel` |
| `--output` am Leser: Umgebungsvariable allein genügt (F-496) | `output` aus `optionen("record")` genommen | `TestLeserHilfe`, `TestLeserOptionen` und 19 Einzeltests |
| `--force` am Leser (F-496) | `force` aus `optionen("record")` genommen | `TestLeserHilfe`, `TestLeserOptionen`, `TestParseRecord` |
| `--force` nur `true` oder `false`, `PGWIRE_RECORDER_FORCE=ja` und `--force=1` sind `PGR-E2001` (F-496) | Wertemenge von `--force` nach `strconv.ParseBool` | `TestLeserAlleOptionen`, `TestLeserHilfe` |

Eine Umgebungsvariable mit Präfix ohne passende Option bleibt unbeachtet: Der Leser fragt nur
die Namen seiner Optionen ab; `TestLeserFremdeUmgebung` belegt es für fünf Namen. Ein Mutant,
der sie liest, müsste eine neue Abfrage einführen; keiner gefahren.

F-501: `wahrheitswert`, `dauer` und `stufe` sind an keinem FlagSet mehr angemeldet, sie
prüfen nur noch (`art`) und setzen (`setzeDauer`); `IsBoolFlag` und die drei `String`-Methoden
sind entfernt, der Kommentar an `wahrheitswert` nennt die Rolle als Prüfer, und den Schalter
ohne Wert trägt `art.schalter` (Zeile *boolesche Option ohne Wert* oben). Mutant L2 aus dem
Review gibt es damit nicht mehr.

**Hilfe vor den Prüfungen am Leser** (DoD-Punkt 3). Jede Bedingung der Zusage hat eine
eigene Zeile; keine Mutation blieb grün, keine Lücke im Code.

| Zusage | Mutation in `Parse` / `hilfeVerlangt` / `lies` / Hilfetext | rote Tests |
|---|---|---|
| Hilfe vor der Prüfung der Umgebungsvariablen am Leser | gesetzte Umgebungsvariablen der Optionen vor der Hilfe geprüft | `TestParseHilfeVorPruefung`, `TestParseLogLevelHilfe`, `TestParseShutdownTimeoutHilfe`, `TestRunHilfe`, `TestLeserReihenfolge`, `TestLeserLeererWert` |
| Hilfe vor der Prüfung der Werte auf der Kommandozeile am Leser | `--<option>=<wert>` vor der Hilfe geprüft | `TestParseHilfeVorPruefung`, `TestParseLogLevelHilfe`, `TestLeserAlleOptionen`, `TestLeserLeererWert` |
| Hilfe auch vor `--force` | `--force=<wert>` vor der Hilfe geprüft | `TestParseHilfeVorPruefung`, `TestLeserAlleOptionen`, `TestLeserNurAngemeldete` |
| nach `--` ist eine Hilfe-Angabe ein gewöhnliches Argument | `hilfeVerlangt` ohne Halt an `--` | `TestParseKeineHilfe`, `TestParseEndeDerOptionen` |
| `--` beendet die Optionen am Leser | `lies` ohne `endeDerOptionen` | `TestParseEndeDerOptionen`, `TestRunEndeDerOptionen` |
| Hilfe nennt die Umgebungsvariablen, ohne Ausnahme (F-496) | Ausnahme `außer --output` in den Satz gesetzt | `TestLeserHilfe` |
| Hilfe nennt, dass ein leerer Wert auf der Kommandozeile ungültig ist (F-497) | Satz entfernt | `TestLeserHilfe` |

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
