# Slice slice-v1-abschluss-konfigurationsdatei: Konfigurationsdatei und config show

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md), [ADR-0014](../../adr/0014-konfigurationsdatei.md), [ADR-0027](../../adr/0027-yaml-bibliothek.md), [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), [ADR-0037](../../adr/0037-yaml-bibliothek-fuer-die-anzeige-der-konfigurationsdatei.md)

**Berührte Spec-Stellen:** `LH-FA-01.a` · `LH-FA-17.a` · `LH-FA-03.b` · `LH-FA-14.a` · `SPEC-007` · `SPEC-008` · `SPEC-012` · `SPEC-014` · `SPEC-020` · `SPEC-033` · `SPEC-034` · `SPEC-046` · `ARC-005` · `ARC-013`

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

**Ziel:** Jede Option von `record` und `replay`, die am allgemeinen Leser angemeldet ist, ist auch über die Konfigurationsdatei setzbar, in der Priorität Kommandozeile vor Umgebungsvariable vor Datei vor Default; `--config`, `PGWIRE_RECORDER_CONFIG` und die Standarddatei wählen genau eine Datei, `config show` zeigt sie, die Hilfe geht jeder Prüfung vor, eine ungültige Datei ist `PGR-E2004`, und die Code-Tabelle führt `PGR-E2004` bis `PGR-E2006`.

**Übernimmt:** `slice-v1-abschluss-konfiguration` — dessen Datei-Hälfte nach dem Schnitt aus
seinem §4 (*zu groß*, eingetreten am 2026-10-08): DoD-Punkt 2 für `config show` und `--config`
(Hilfe vor jeder Prüfung, `config show`) und DoD-Punkt 3 (Datei, Verbindungen, Platzhalter,
`PGR-E2004` bis `PGR-E2006`); dort §1, Abgrenzung. Im Einzelnen, je mit dem ursprünglichen
Geber, über `slice-v1-abschluss-konfiguration` und davor `slice-v1-abschluss-betrieb`:

- **Aus `slice-replay-semantik-mismatch`:** der Schlüssel `fail_on_unconsumed` im Abschnitt
  `replay:`; die Hilfe vor jeder Prüfung von Optionen, Umgebungsvariablen und
  Konfigurationsdatei (`LH-FA-01.a`) für `config show` und `--config`, samt `--` als Ende der
  Optionen.
- **Aus `slice-replay-semantik-meldungscodes`:** der Schlüssel `log_level` auf der obersten
  Ebene der Konfigurationsdatei (`LH-FA-17.a`), mit derselben Wertemenge wie `--log-level`;
  ein ungültiger Wert ist `PGR-E2004`.
- **Aus `slice-v1-abschluss-herunterfahren`** (dort §1, Abgrenzung): der Schlüssel der Frist
  `--shutdown-timeout` in der Konfigurationsdatei.

**Abgegeben** an `slice-v1-abschluss-verbindungen-platzhalter` (dort §1, *Übernimmt*, mit der
Kennung dieses Slice), nach dem vorab benannten Schnitt aus §4 (*zu groß*, eingetreten am
2026-10-09): der frühere DoD-Punkt 3 — Grammatik der URL und `sslmode`, `${VAR}` und `$$`,
`PGR-E2005` und `PGR-E2006` als Fehler, das Auflösen eines Verbindungsnamens in `--upstream`
und der Teil des Benutzerhandbuchs dazu —, samt den Randformen dazu aus §6 (Rückgaben 5 zum
Teil, 6, 7, 8 und 10). Hier bleiben die Wahl und das Laden der Datei, die Priorität über drei
Quellen, `config show` mit der Hilfe, die Konstanten und von `connections:` die Form, die
DoD-Punkt 2 prüft: eine Abbildung, ein Name in der Form eines Werts ohne Steuerzeichen, ein
Wert als Skalar (Rückgaben 4 und 11). Bis zum Nehmer gilt ein Wert mit `$` wörtlich, und eine
URL wird nicht zerlegt.

**Schlüssel der Datei.** Die Schlüssel, die die Datei annimmt, folgen aus der Anmeldung am
allgemeinen Leser von `slice-v1-abschluss-konfiguration`; keine Option wird ein zweites Mal
angemeldet. Ein Schlüssel einer Option aus der Tabelle in `LH-FA-17.a`, die der Stand noch
nicht kennt (etwa `format`, `tls_cert` oder der Abschnitt `play:`), ist bis zu ihrem Slice ein
unbekannter Schlüssel (`PGR-E2004`). `output` und `force` sind am Leser angemeldet und damit
hier Schlüssel des Abschnitts `record:` (Rückgabe 12 in §6). Der
Test über alle angemeldeten Optionen läuft dann über drei Quellen.

**Ort des Lesens.** Die Datei liest der CLI-Adapter (`ARC-005`). Die YAML-Bibliothek ist dort
nur zum Lesen der Konfigurationsdatei zulässig (`ARC-013`,
[ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), ergänzt
[ADR-0027](../../adr/0027-yaml-bibliothek.md)); kein Typ der Bibliothek verlässt den
CLI-Adapter. Die Ausgabe derselben Datei in `config show` mit dem Kodierer der Bibliothek
trägt [ADR-0037](../../adr/0037-yaml-bibliothek-fuer-die-anzeige-der-konfigurationsdatei.md)
(Accepted, Entscheidung des Nutzers vom 2026-10-09). Die Gegenprobe des Architektur-Gates dafür liefert
`slice-v1-abschluss-konfiguration`; dass ein Import im CLI-Adapter nur der Datei dient, prüft
kein Gate, das bleibt Review.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der allgemeine Leser für Kommandozeile und Umgebungsvariable, die Hilfe von `record` und
  `replay` und die Gegenprobe des Architektur-Gates — `slice-v1-abschluss-konfiguration`;
  er liegt vor diesem Slice, und dieser setzt auf ihm auf.
- Das atomare Schreiben und das Verhalten bei vorhandenem `--output` —
  `slice-v1-abschluss-schreiben`; er folgt diesem Slice. `--output` und `--force` selbst
  meldet `slice-v1-abschluss-konfiguration` am Leser an (Entscheidung zu F-496).
- Die Optionen von `play` und der Abschnitt `play:` — `slice-v1-abschluss-einspielen` (dort
  §1); er liest sie über den allgemeinen Leser und die Datei dieses Slice.
- Grammatik der URL, `sslmode`, Platzhalter, `PGR-E2005`, `PGR-E2006` und das Auflösen eines
  Verbindungsnamens — `slice-v1-abschluss-verbindungen-platzhalter` (oben, *Abgegeben*); er
  folgt diesem Slice.
- TLS zum Upstream nach `sslmode` — `slice-v1-abschluss-verbindungen-platzhalter` (oben,
  *Abgegeben*); er gibt die Wirkung bei `play` an `slice-v1-abschluss-einspielen` weiter
  (dort §1, *Übernommen aus* jenem Slice). Nur `play` verbindet mit TLS zum Server.
- Code im Kern, im PGWire-Adapter und in den Driven-Adaptern — Schicht-Abgrenzung: Der
  Slice ändert den CLI-Adapter und den Bootstrap. **Einzige Ausnahme:** drei Konstanten
  `PGR-E2004` bis `PGR-E2006` in der Code-Tabelle `internal/hexagon/model/fehler.go`, ohne
  Logik — die Tabelle ist die eine Code-Tabelle im Quelltext
  ([ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md), Konsequenzen), und dort steht
  schon `PGR-E2001`, das nur der CLI-Adapter und der Bootstrap erzeugen (Entscheidung des
  Architect vom 2026-10-08, §6).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Die Hilfe geht jeder Prüfung von Optionen,
      Umgebungsvariablen und Konfigurationsdatei vor, auch für `config show` und `--config`,
      samt `--` als Ende der Optionen; `config show` zeigt die gewählte Datei in der Form aus
      `LH-FA-17.a`, ohne einen aufgelösten Wert (Test).
- [x] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): `--config`, `PGWIRE_RECORDER_CONFIG` und die Standarddatei wählen genau
      eine Datei; für jede am allgemeinen Leser angemeldete Option gilt die Priorität
      Kommandozeile vor Umgebungsvariable vor Konfigurationsdatei vor Default (`SPEC-007`),
      auch für den Schlüssel `fail_on_unconsumed` im Abschnitt `replay:`, den Schlüssel der
      Frist und den Schlüssel `log_level` auf der obersten Ebene (Wertemenge und Strenge von
      `--log-level`); eine ungültige Datei ist `PGR-E2004` (Test). Die Code-Tabelle führt
      dazu die Konstanten `PGR-E2005` und `PGR-E2006` ohne Erzeuger; ihre Tests liefert
      `slice-v1-abschluss-verbindungen-platzhalter`.
- [x] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Die Vorgaben des Architect vom 2026-10-09 (§6, *Rückgaben vom
      2026-10-09*) sind als Test belegt, je mit ihrer Mutation rot gesehen: R1 —
      `output: "null"` und `output: '~'` setzen den Pfad `null` bzw. `~`, `Null` und `NULL`
      ohne Anführungszeichen sind `PGR-E2004`, rot über den grünen Mutanten aus §7 und über
      die Prüfung auf `!!null` als Textvergleich; R2 — die Standarddatei als Link auf ein
      fehlendes Ziel wählt keine Datei, und `config show` sagt das in der ersten Zeile, rot
      über `os.Lstat` statt `os.Stat`; Lesart 3 — `--config` mit relativem Pfad, die erste
      Zeile von `config show` ist genau dieser Pfad, rot über `filepath.Abs`; Lesart 2 — ein
      unbekanntes Unterkommando als `show` behandelt, rot über `{"config", "zeige"}` (Test).
      Beleg in §7 für alle drei Punkte: je Zusage Zusage · Mutation · roter Test
      (`AGENTS.md` §3.10).
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
| `spec/spezifikation.md` (`LH-FA-17.a`) | update (Architect, vor dem Code, erledigt am 2026-10-08 und am 2026-10-09) | Entscheidung der zwölf Rückgaben aus §6 (`AGENTS.md` §3.12); R1, R2 und fünf Lesarten nach dem Code von DoD-Punkt 1 und 2 (`cfc4d6b`) |
| `internal/adapters/driving/cli` | update | Wahl und Laden der Datei (`--config`, `PGWIRE_RECORDER_CONFIG`, Standarddatei; nur eine reguläre Datei, Kodierung und Zeilenenden nach den Befunden des Reviews, §6); Schlüssel aus der Anmeldung am allgemeinen Leser als dritte Quelle der Priorität; `fail_on_unconsumed`, `log_level`, Frist; Form von `connections:` (Abbildung, Name, Skalar); Kommando `config show` und seine Hilfe — geliefert in `7a80393` |
| `internal/hexagon/model/fehler.go` | update | nur die drei Konstanten `PGR-E2004` bis `PGR-E2006` in der Code-Tabelle (§1, Ausnahme; §6) |
| `internal/bootstrap` | update | `config show` ausführen (Ausgabe auf `stdout`); die zusammengeführten Optionen an die Use Cases geben |
| `internal/adapters/driving/cli` (Unit-Tests), `internal/bootstrap` (Tests), `test/integration` | update | Happy/Boundary/Negative nach `LH-FA-17.a` und `LH-FA-01.a`; der Test über alle angemeldeten Optionen läuft über drei Quellen. Die Tests der Vorgaben vom 2026-10-09 (DoD-Punkt 3) — R1 in `TestDateiGueltig` und `TestDateiUngueltig`, R2 und Lesart 3 in `TestConfigShowOhneDatei`, die Mutation zu Lesart 2 über `TestConfigShowFehler` — geliefert in `9332356`; die Tests zu den Befunden des Reviews vom 2026-10-09 (§6) in `TestDateiTag`, `TestDateiKodierung`, `TestDateiKeineRegulaere`, `TestConfigShowOhneBOM`, `TestDateiUngueltig`, `TestConfigShow` und `TestE2EConfigShowStdin` |
| `docs/user/benutzerhandbuch.md` | update, falls abweichend | §4 *Die gewählte Konfigurationsdatei anzeigen* und §5 *Konfigurationsdatei* ohne die Absätze zu Verbindungen und Platzhaltern (die zieht `slice-v1-abschluss-verbindungen-platzhalter` nach) beschreiben den Zielstand; nachgezogen wird, was die Entscheidungen aus §6 ändern |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-konfiguration` liegt in `done/` (der
allgemeine Leser, aus dessen Anmeldung die Schlüssel folgen). Schritt 4 der Reihenfolge in §5
von [welle-v1-abschluss](../welle-v1-abschluss.md) (Entscheidung des Nutzers vom 2026-10-08).
Die zwölf Randformen und den Ort der Konstanten aus §6 entschied der Architect am 2026-10-08
vor dem Code (`AGENTS.md` §3.12).

**Erneuter Start** (`next` → `in-progress`) nach der Rückführung unten: dieser Plan im
geschnittenen Zuschnitt liegt auf dem Hauptzweig. Die Arbeit ist DoD-Punkt 3 (die Tests der
Vorgaben vom 2026-10-09); DoD-Punkt 1 und 2 sind geliefert (`7a80393`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar, oder eine Änderung im Kern über die Konstanten hinaus wird nötig.
  Schnitt dann: *Datei* (DoD-Punkt 1 und 2) und *Verbindungen und Platzhalter* (DoD-Punkt 3);
  der zweite setzt den ersten voraus. **Eingetreten am 2026-10-09** (Grund unten). Für den
  geschnittenen Zuschnitt: Ein Test einer Vorgabe vom 2026-10-09 verlangt eine Änderung am
  Code, die den Diff über eine Review-Sitzung hebt; dann zurück an den Architect, weil er dort
  *Code bleibt* entschied.
- `in-progress` → `open` (blockiert — Carveout?): Die Form der Datei nach [ADR-0014](../../adr/0014-konfigurationsdatei.md) trägt
  einen Schlüssel nicht, den eine Option braucht, ohne eine neue Entscheidung; dann zuerst
  die Entscheidung.

**Grund der Rückführung `in-progress` → `next`, eingetreten am 2026-10-09:** Der Implementer
hielt nach DoD-Punkt 1 und 2 an; deren Diff allein umfasst rund 1250 Zeilen, davon rund 600
Tests (`7a80393`, Beleg `a063f2f` in §7), mit dem damaligen DoD-Punkt 3 läge er über einer
Review-Sitzung. Das ist die erste Bedingung oben. Geschnitten wird nach dem vorab benannten
Schnitt: *Datei* bleibt hier, *Verbindungen und Platzhalter* geht an
`slice-v1-abschluss-verbindungen-platzhalter` (§1, *Abgegeben*). Der Architect entschied
vorher R1, R2 und fünf Lesarten zu DoD-Punkt 1 und 2 (`cfc4d6b`); deren Tests fehlen und sind
jetzt DoD-Punkt 3. Der Übergang läuft formal über `next/` (Entscheidung des Nutzers vom
2026-10-09 zu V-117), als reiner `git mv` nach diesem Commit, dann `next` → `in-progress`.

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

**Randformen** (`AGENTS.md` §3.12) — je Randform, wo sie entschieden ist. Die entschiedenen
zogen mit dem Schnitt aus §6 von `slice-v1-abschluss-konfiguration` hierher; der Architect
prüfte sie dort am 2026-10-08 vor dem Code. Die zwölf Rückgaben des Implementers und den Ort
der Konstanten entschied er am 2026-10-08 hier, vor dem Code (unten). Zwei Rückgaben und fünf
Lesarten aus dem Lauf an DoD-Punkt 1 und 2 entschied er am 2026-10-09 (unten, *Rückgaben vom
2026-10-09*). Drei Befunde des Reviews vom 2026-10-09 (F-506, F-507, F-517) entschied er am
selben Tag (unten, *Befunde des Reviews vom 2026-10-09*), zwei Befunde der Verifikation
vom 2026-10-09 (V-118, V-119) ebenso (unten, *Befunde der Verifikation vom 2026-10-09*); offen
ist keine. [ADR-0037](../../adr/0037-yaml-bibliothek-fuer-die-anzeige-der-konfigurationsdatei.md) ist
angenommen (Entscheidung des Nutzers vom 2026-10-09). Die Randformen zu Verbindungen und Platzhaltern und die
Rückgaben 6, 7, 8 und 10 sowie der Teil der Verbindungen aus Rückgabe 5 zogen mit dem Schnitt
vom 2026-10-09 nach §6 von `slice-v1-abschluss-verbindungen-platzhalter`; hier steht je eine
Zeile mit dem Ort.

*Umgebung, soweit sie die Datei betrifft*

- **Gesetzte, aber leere `PGWIRE_RECORDER_CONFIG`** — gilt als nicht gesetzt; `LH-FA-17.a`.
- **Umgebungsvariable mit Präfix ohne passende Option** — unbeachtet; `config show` nennt ihren
  Namen; `LH-FA-17.a`.

*Datei: Wahl und Form*

- **Ort der Standarddatei je Betriebssystem** — nur `.pgwire-recorder.yaml` im aktuellen
  Verzeichnis, auf jedem Betriebssystem, kein Benutzer- oder Systemverzeichnis; `LH-FA-17.a`.
- **Standarddatei fehlt / Datei aus `--config` oder `PGWIRE_RECORDER_CONFIG` fehlt** — keine
  Datei gelesen / `PGR-E2004`; `LH-FA-17.a`. Standarddatei vorhanden, aber nicht lesbar:
  `PGR-E2004` (Tabelle *Fehler*).
- **`PGWIRE_RECORDER_CONFIG` nennt eine fehlende Datei, `--config` geht vor** — gelesen und
  geprüft wird nur die gewählte Datei; `LH-FA-17.a` (*genau eine Datei*).
- **Groß- und Kleinschreibung von Schlüsseln** — Schlüssel heißen wie die Optionen;
  `Log_Level` ist ein unbekannter Schlüssel (`PGR-E2004`); `LH-FA-17.a`.
- **Unbekannter Schlüssel, Schlüssel im falschen Abschnitt** — `PGR-E2004`; `LH-FA-17.a`.
  Schlüssel einer Option, die der Stand noch nicht kennt: unbekannt (§1).
- **Doppelter Schlüssel** — ungültiges YAML, `PGR-E2004`; `LH-FA-17.a`. Hinweis an den
  Implementer: Das Lesen über die Node-Schnittstelle der Bibliothek lehnt Doppelte nicht
  selbst ab.
- **Schlüssel `config` in einem Abschnitt** — gibt es nicht, unbekannter Schlüssel
  (`PGR-E2004`); Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`.
- **Wertform in der Datei** (YAML `true` gegen `"true"`, `True`, `yes`, `1`; `shutdown_timeout: 0`
  als Zahl; leerer Wert, `null`, Liste oder Abbildung als Wert; Anker und Aliase) — Text des
  Skalars mit der Wertemenge der Option; leerer Wert, `null`, Liste, Abbildung, Anker, Aliase
  und Merge-Schlüssel `PGR-E2004`; Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`.
- **Form der Dauer** — dieselbe in Option, Umgebungsvariable und Datei; ungültig `PGR-E2001`
  bzw. `PGR-E2004`; `LH-FA-17.a` (*Dauer*). `0` ohne Anführungszeichen ist gültig (Wertform
  oben).
- **`log_level` ungültig** — `PGR-E2004`; `LH-FA-17.a`.
- **Relative Pfade in der Datei** (`input`, `output`, später Zertifikate) — relativ zum
  aktuellen Verzeichnis; Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`.

*Verbindungen und Platzhalter* — Name bei `--upstream`, `sslmode=require` bei `record`,
Klartext-Passwort, nicht gesetzte oder leere Variable, Platzhalter in den von `record`
ignorierten Teilen, `$${VAR}`, Rekursion, Form des Namens und Einsetzen in die URL: gezogen
nach `slice-v1-abschluss-verbindungen-platzhalter` (dort §6), unverändert entschieden in
`LH-FA-17.a`.

*Anzeige und Fehler*

- **`config show` und Geheimnisse** — keine Maskierung, und keine nötig: Ein Passwort steht
  in keiner gültigen Datei (`PGR-E2006`, erzeugt erst von
  `slice-v1-abschluss-verbindungen-platzhalter`; Risiko unten), Platzhalter erscheinen
  unaufgelöst, Umgebungsvariablen nur mit Namen; `LH-FA-17` (Lastenheft: *die Anzeige zeigt sie nicht*) und
  `LH-FA-17.a` (*Anzeige*). `LH-RB-01` gilt der Aufzeichnung (keine Maskierungsfunktion für
  deren Inhalt, `SPEC-033`), nicht der Anzeige der Konfiguration; kein Widerspruch.
- **`config show`: aktive Umgebungsvariablen** — gesetzt und nicht leer, auch ohne passende
  Option, geprüft nur `PGWIRE_RECORDER_CONFIG`; `LH-FA-17.a`. Dass `config show` kein
  `PGR-E2005` meldet, zog nach `slice-v1-abschluss-verbindungen-platzhalter` (dort §6).
- **Form der Ausgabe von `config show`** — erste Zeile Pfad der Datei oder dass keine
  gefunden wurde, danach der Inhalt als YAML (zwei Leerzeichen, Reihenfolge der Datei, ohne
  Kommentare), danach die Namen der aktiven Umgebungsvariablen sortiert, alles auf `stdout`;
  Entscheidung des Nutzers vom 2026-10-08, `LH-FA-17.a`.
- **Hilfe nach einer ungültigen Option, auch bei `config show` und `--config`** — die Hilfe
  geht vor; `LH-FA-01.a`.
- **Reihenfolge bei mehreren Fehlern** — Abbruch beim ersten Fehler; Reihenfolge
  Kommandozeile, Umgebungsvariablen nach der Tabelle, Datei in ihrer Reihenfolge, zuletzt
  Pflichtoptionen, Kombinationen, `--upstream` und `PGR-E2005`; Entscheidung des Nutzers vom
  2026-10-08, `LH-FA-17.a`. Den Teil der Verbindungen (innerhalb einer URL, `--upstream`,
  `PGR-E2005`, Port nach dem Einsetzen) liefert `slice-v1-abschluss-verbindungen-platzhalter`
  (dort §6).


*Rückgaben des Implementers vom 2026-10-08* (aus dem abgebrochenen Lauf an
`slice-v1-abschluss-konfiguration`), entschieden vom Architect am 2026-10-08 vor dem Code,
alle in `LH-FA-17.a`, abgeleitet aus den Grundsätzen dort: Text des Skalars zählt, Laden ohne
Kommando (`config show` prüft die Datei ohne Kommando, also prüft jedes Kommando sie ganz),
geschlossene Fehlertabelle, URL vor dem Einsetzen zerlegt, Abbruch beim ersten Fehler.

1. **Leere Datei oder nur Kommentare** — setzt nichts, kein Fehler; keine Ursache der
   Fehlertabelle trifft zu.
2. **Mehrere YAML-Dokumente** — ein zweites Dokument ist ungültiges YAML (`PGR-E2004`); es
   still zu übergehen, ließe Eingabe unbeachtet ([ADR-0014](../../adr/0014-konfigurationsdatei.md): jede unbekannte Eingabe ist ein
   Startfehler). Ein einzelnes `---` vor dem einen Dokument ist zulässig.
3. **Oberste Ebene keine Abbildung** — `PGR-E2004`.
4. **Leerer Abschnitt, leeres `connections:`** — `{}` setzt nichts; ohne Inhalt oder `null`
   ist `PGR-E2004` wie ein leerer Wert. Akzeptiertes Negativ: Ein Abschnitt, dessen Inhalt
   ganz auskommentiert ist, wird damit ungültig; das folgt der Strenge bei `null`
   (Entscheidung des Nutzers), und die Meldung nennt den Abschnitt.
5. **Abschnitt eines anderen Kommandos** — mitgeprüft. Der Teil der Verbindungen (jede nicht
   benutzte Verbindung mitgeprüft; vom Kommando hängen nur `sslmode=require` bei `record` und
   die Variablen der Platzhalter ab) zog nach `slice-v1-abschluss-verbindungen-platzhalter`
   (dort §6, Rückgabe 5).
6. **URL ohne Port** — gezogen nach `slice-v1-abschluss-verbindungen-platzhalter` (dort §6,
   Rückgabe 6).
7. **Platzhalter-Syntax in ignorierten Teilen** — gezogen nach
   `slice-v1-abschluss-verbindungen-platzhalter` (dort §6, Rückgabe 7).
8. **Eingesetzter Wert macht seinen Teil ungültig** — gezogen nach
   `slice-v1-abschluss-verbindungen-platzhalter` (dort §6, Rückgabe 8).
9. **YAML-Tags** — jeder ausdrücklich geschriebene Tag ist `PGR-E2004`, wie Anker und Aliase;
   ein Tag behauptet eine Bedeutung jenseits des Texts.
10. **URL-Sonderformen** — gezogen nach `slice-v1-abschluss-verbindungen-platzhalter` (dort
    §6, Rückgabe 10).
11. **Name einer Verbindung** — Form eines Werts (Text des Skalars; leer, `null`, Tag
    ungültig), dazu kein Steuerzeichen, weil die Meldung den Namen in einer Zeile nennt
    (`SPEC-034`); sonst jeder Name, auch mit Leerraum.
12. **`output` und `force` in der Datei** — Schlüssel des Abschnitts `record:` mit diesem
    Slice, weil `slice-v1-abschluss-konfiguration` beide am Leser anmeldet (Entscheidung des
    Architect vom 2026-10-08 zu F-496; sie ersetzt die frühere Antwort *unbekannt bis
    `slice-v1-abschluss-schreiben`*). `LH-FA-17.a` führt beide in der Tabelle, kein
    Sonderfall.

*Rückgaben und Lesarten des Implementers vom 2026-10-09* (aus §7, *Belege des Implementers*),
entschieden vom Architect am 2026-10-09 in `LH-FA-17.a`, nach dem Code von DoD-Punkt 1 und 2
(`7a80393`). In keinem Fall ändert sich der Code; je Fall eine Vorgabe für Test und Mutation,
die vor der Übergabe ans Review rot gesehen wird (`AGENTS.md` §3.10).

- **R1 — `"null"` oder `"~"` in Anführungszeichen** — Text des Skalars, gültig, wo die
  Wertemenge der Option ihn annimmt; `null` meint die YAML-Null, ohne Anführungszeichen `~`,
  `null`, `Null`, `NULL`. Grund: Dieselbe Wertemenge in allen drei Quellen trägt, und
  `--output=null` ist auf der Kommandozeile ein Pfad; Anführungszeichen sind in YAML der Weg,
  diesen Text zu schreiben. Code bleibt. *Test:* `record:` mit `output: "null"` und mit
  `output: '~'` setzt den Pfad `null` bzw. `~` (`TestDateiGueltig`); `Null` und `NULL` ohne
  Anführungszeichen sind `PGR-E2004` (`TestDateiUngueltig`). *Mutation:* der grüne Mutant aus
  §7 (Text in Anführungszeichen als `null`) wird rot; dazu die Prüfung auf `!!null` durch einen
  Textvergleich mit `null` und `~` ersetzt, rot über `Null`/`NULL`.
- **R2 — Standarddatei als Link ins Leere** — keine Datei; die Standarddatei existiert, wenn
  ihr Pfad nach Auflösung der Links auflöst; lässt sich das nicht feststellen (Schleife), ist
  sie vorhanden, aber nicht lesbar (`PGR-E2004`). Grund: „sofern sie existiert“ ist die
  Antwort des Betriebssystems auf den Pfad, wie bei `test -e`, und `config show` zeigt in der
  ersten Zeile, dass keine Datei gilt; die Abweichung wird damit nicht still. Akzeptiertes
  Negativ: Wer die Standarddatei als Link auf eine gelöschte Datei hält, startet ohne Datei;
  `--config` und `PGWIRE_RECORDER_CONFIG` melden denselben Fall als fehlende Datei
  (`PGR-E2004`). Code bleibt (`os.Stat`). *Test:* Standarddatei als Link auf ein fehlendes
  Ziel wählt keine Datei, `config show` sagt das in der ersten Zeile (`TestDateiWahl` oder
  `TestConfigShowOhneDatei`). *Mutation:* `os.Lstat` statt `os.Stat`, rot über diesen Test;
  die Schleife hält `TestDateiNichtLesbar` schon.
- **Lesart 1 — `config --help`** — bestätigt: Hilfe von `config show` (`LH-FA-01.a`: das erste
  Argument ist ein bekanntes Kommando). Getestet (`TestDateiHilfeVorPruefung`), keine Vorgabe.
- **Lesart 2 — `config` allein, anderes Unterkommando, Argument nach `config show`** —
  bestätigt: `PGR-E2001`, wie das unerwartete Argument bei `record` und `replay`. Getestet
  (`TestConfigShowFehler`, alle drei Fälle). *Mutation*, falls nicht gefahren: ein
  unbekanntes Unterkommando als `show` behandelt, rot über `{"config", "zeige"}`.
- **Lesart 3 — erste Zeile von `config show`** — bestätigt: der Pfad, wie er gewählt wurde,
  ein relativer bleibt relativ, die Standarddatei heißt `.pgwire-recorder.yaml`. *Test:*
  `--config` mit relativem Pfad, erste Zeile genau dieser Pfad. *Mutation:* `filepath.Abs`
  auf den Pfad, rot über diesen Test.
- **Lesart 4 — Datei nur mit `---`** — bestätigt: ein Dokument ohne Inhalt, oberste Ebene keine
  Abbildung, `PGR-E2004`; das folgt der Strenge bei `null` (Rückgabe 4 oben). Getestet
  (`TestDateiUngueltig`, *nur ---*), keine Vorgabe.
- **Lesart 5 — doppelte Schlüssel** — bestätigt: nach ihrem Text verglichen, `1` und `"1"`
  sind derselbe Schlüssel; das folgt „Text des Skalars zählt“. Getestet
  (`TestDateiUngueltig`, *doppelt in Anführungszeichen*), keine Vorgabe.

*Befunde des Reviews vom 2026-10-09*, entschieden vom Architect am 2026-10-09 in
`LH-FA-17.a` bzw. in [ADR-0037](../../adr/0037-yaml-bibliothek-fuer-die-anzeige-der-konfigurationsdatei.md).
Der Code folgt nach dieser Entscheidung; je Fall Test und Mutation, vor der Übergabe ans Review
rot gesehen (`AGENTS.md` §3.10). Die Fälle stehen in `TestDateiGueltig`, `TestDateiUngueltig`
bzw. `TestDateiNichtLesbar`, wo nichts anderes steht.

- **F-506 — Tag** — der nicht spezifische Tag `!` ist ein Tag (`LH-FA-17.a`). Die Prüfung über
  `TaggedStyle` kommt zurück und **ergänzt** die Prüfung am Text: `TaggedStyle` fängt jeden Tag
  mit Namen unabhängig von Zeile und Spalte; nur `!` allein setzt die Bibliothek nicht als
  `TaggedStyle` (Sonde des Architect: `log_level: ! info` ergibt `!!str`, ohne `TaggedStyle`),
  ihn fängt allein die Prüfung am Text. Diese ist verlässlich, weil die Zählung von Zeile und
  Spalte der Bibliothek folgt: führendes BOM vor dem Zählen entfernt, Zeilenenden `\n`, `\r\n`
  und `\r`, und was die Bibliothek anders zählt, ist vorher abgelehnt (F-507). Die Kommentare an
  `form` und `zeilen` sagen genau das zu und nicht mehr (§3.11).
  *Tests:* (1) `form` direkt mit einem Knoten `Style: TaggedStyle` und Zeile außerhalb des Texts
  (`z` leer) meldet „Tag ist ungültig“; *Mutation:* Prüfung über `TaggedStyle` entfernt, rot.
  (2) `log_level: ! info` und `log_level: !!str info` je mit führendem BOM, mit Zeilenende nur
  `\r` (Tag in Zeile 2) und mit `\r\n` sind `PGR-E2004`; *Mutation:* BOM vor dem Zählen nicht
  entfernt bzw. nur nach `\n` gezählt, rot über den Fall `!` allein. (3) Gültig bleibt
  `record:\r  force: true\nlog_level: info\n#          !\n` (das `!` steht in der Spalte, die eine Zählung
  nur nach `\n` dem Wert von `log_level` zuordnete); *Mutation:* nur nach `\n` gezählt, rot.
- **F-507 (a) — Kodierung** — nur UTF-8; ein BOM als erstes Zeichen wird **hingenommen** (Entscheidung des Nutzers vom
  2026-10-09) und
  übergangen (YAML 1.2 erlaubt ihn am Anfang des Streams, Werkzeuge unter Windows schreiben
  ihn); ein BOM an anderer Stelle, eine ungültige UTF-8-Folge und UTF-16 oder UTF-32, auch mit
  BOM, sind ungültiges YAML, `PGR-E2004` mit der Zeile (`LH-FA-17.a`). Die Bibliothek nimmt
  UTF-16 mit BOM und ein BOM in einem Wert in Anführungszeichen an (Sonde des Architect); die
  Prüfung liegt deshalb vor dem YAML.
  *Tests:* gültige Datei mit führendem BOM setzt dieselben Werte wie ohne; BOM am Anfang von
  Zeile 2 und BOM in einem Wert in Anführungszeichen sind `PGR-E2004` mit „Zeile 2“; dieselbe
  Datei als UTF-16LE und als UTF-16BE mit BOM und als UTF-16LE ohne BOM ist `PGR-E2004`; ein
  Byte `0xFF` in einem Kommentar ist `PGR-E2004`, die Meldung nennt kein Byte der Datei.
  *Mutation:* Prüfung auf BOM nach dem ersten Zeichen entfernt, rot über den Wert in
  Anführungszeichen; Prüfung auf UTF-8 entfernt, rot über UTF-16 mit BOM.
- **F-507 (b) — Zeilenende** — `\n`, `\r\n` und `\r` allein nach YAML 1.2; `U+0085`, `U+2028`
  und `U+2029` sind an jeder Stelle ungültiges YAML, `PGR-E2004` mit der Zeile, auch in
  Anführungszeichen und in einem Kommentar (`LH-FA-17.a`). Grund: Die Bibliothek liest sie
  nach YAML 1.1 als Zeilenende, YAML 1.2 nicht; jede andere Lesart widerspräche entweder YAML 1.2
  oder dem Wert, den die Bibliothek liefert.
  *Tests:* Datei nur mit `\r` und Datei mit `\r\n` setzen dieselben Werte wie mit `\n`; je eines
  der drei Zeichen in einem Wert in Anführungszeichen und `U+2028` in einem Kommentar sind
  `PGR-E2004` mit der Zeile, in der das Zeichen steht. *Mutation:* Prüfung der drei Zeichen
  entfernt, rot über den Wert in Anführungszeichen.
- **F-507 (c) — keine reguläre Datei** — die gewählte Datei muss, Links gefolgt, eine reguläre
  Datei sein; Verzeichnis, FIFO, Gerät und Socket sind vorhanden, aber nicht lesbar
  (`PGR-E2004`, „nicht lesbar“ mit Quelle, ohne Pfad), geprüft vor dem Öffnen, weil schon
  das Öffnen einer FIFO blockiert. Das gilt für alle drei Quellen, die Standarddatei
  eingeschlossen. `--config /dev/stdin` liest eine umgeleitete reguläre Datei und lehnt Pipe
  und Terminal ab (Entscheidung des Nutzers vom 2026-10-09); Lesen von `stdin` gibt es nicht, eine Datei lässt sich direkt nennen.
  Akzeptiertes Negativ: Wird der Pfad zwischen Prüfung und Öffnen durch eine FIFO ersetzt,
  blockiert der Start; das verlangt, dass jemand die Datei während des Starts austauscht, und
  hat keine andere Wirkung als ein hängender Start.
  *Tests:* `--config` auf ein Verzeichnis (Bestand) und auf `/dev/null` ist `PGR-E2004`
  „nicht lesbar“; eine FIFO (`syscall.Mkfifo`, Datei mit `//go:build unix`) ist `PGR-E2004`,
  `ladeDatei` läuft dabei in einer Goroutine, und der Test wartet auf ihr Ergebnis mit einer
  Frist als Literal (`SPEC-038`); die Standarddatei als Verzeichnis ist `PGR-E2004`; `--config`
  als Link auf eine reguläre Datei wird gelesen. E2E in `konfiguration_e2e_test.go`: das Binary
  mit `config show --config /dev/stdin` und `stdin` aus einer Pipe endet mit Exit `2` und
  `PGR-E2004` innerhalb der Frist; mit `stdin` aus einer regulären Datei zeigt es deren Inhalt.
  *Mutation:* Prüfung entfernt, rot über `/dev/null` (leere Datei angenommen) und über die
  Frist bei FIFO und Pipe; `os.Lstat` statt `os.Stat`, rot über den Link.
- **F-507 (d) — leerer Schlüssel, Schlüssel kein Skalar** (F-511) — unbekannt, die Meldung
  nennt als Stelle die Abbildung, in der er steht (`oberste Ebene`, Abschnitt); `LH-FA-17.a`.
  *Test:* `"": 1` und `? [a]` auf der obersten Ebene und in `record:` nennen `oberste Ebene`
  bzw. `record`. *Mutation:* Stelle wieder aus dem Text des Schlüssels, rot.
- **F-508 — Meldung zu einem doppelten Schlüssel** — wie jede Meldung zu ungültigem YAML die
  Zeile, nicht der Text des Schlüssels; die Phase *YAML* weiß noch nicht, ob er an der Stelle
  eines Werts steht (`LH-FA-17.a` *Fehler*). Doppelte in jeder Tiefe, auch in einer Liste,
  sind ungültiges YAML (F-510). *Tests:* der Fall der Sonde unter `connections.a` nennt die
  Zeile und weder `GEHEIM` noch den Schlüssel; Doppelte in dritter Ebene und in einer Liste an
  einer Wertstelle sind „ungültiges YAML“, nicht „erwartet ein Skalar“. *Mutation:* Text des
  Schlüssels in der Meldung, rot; M1 und M2 aus dem Review, rot.
- **Ausgabe von `config show`** — ohne BOM und mit `\n` als Zeilenende, gleich wie die Datei
  geschrieben ist (`LH-FA-17.a` *Anzeige*). *Test:* Datei mit BOM und `\r` zeigt kein `\r` und
  kein BOM. Keine Mutation: Die Zusage hält der Kodierer; der Test hält sie gegen eine Ausgabe
  des Texts der Datei.
- **F-517 — Kodierer in `config show`** — Wiederausgabe ist kein Lesen im Sinn von
  [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md); die Nutzung ist
  legitim, aber nicht dokumentiert. Verdikt: nachziehen durch die ergänzende
  [ADR-0037](../../adr/0037-yaml-bibliothek-fuer-die-anzeige-der-konfigurationsdatei.md)
  (Accepted, Entscheidung des Nutzers vom 2026-10-09), keine Supersession, weil
  [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) sonst unverändert
  gilt. Der Code bleibt; keine Vorgabe an den Implementer. Die Sicht nennt in `ARC-013` und §2
  *Zusätzliche Einschränkungen* Lesen und Anzeigen der Konfigurationsdatei (Architect,
  2026-10-09). `.a-check.yml` bleibt unberührt: Die `tech`-Regel erlaubt den Import im
  CLI-Adapter, nicht einen Zweck, und der Import ist derselbe.

*Befunde der Verifikation vom 2026-10-09*, entschieden vom Architect am 2026-10-09 in
`LH-FA-17.a`. Der Code folgt nach dieser Entscheidung; je Fall Test mit genauer Zahl und
Mutation, vor der Übergabe rot gesehen (`AGENTS.md` §3.10). V-121 betrifft den Nehmer und steht
als offene Randform in §6 von `slice-v1-abschluss-verbindungen-platzhalter`.

- **V-118 — Zeile bei einem Konstrukt über mehrere Zeilen** — die Zeile, in der das Konstrukt
  (Folge, Abbildung, Text in Anführungszeichen) beginnt, sonst die Zeile, in der der Fehler
  erkannt wird, gezählt ab 1. Grenze: Beginnt das Konstrukt in Zeile 1, ist es die Zeile der
  Erkennung. Grund: Die Bibliothek gibt die Zeile der Erkennung nur ohne umgebendes Konstrukt
  heraus, sonst dessen Anfang; mehr lässt sich ohne eigenen Parser nicht bestimmen, und der
  Anfang des Konstrukts führt zu der Stelle, an der zu suchen ist. Sonde des Architect
  (Bibliothek v3.0.5): Fehler des Parsers nennen die Zeile ab 0, Fehler des Scanners ab 1;
  liegen beide Marken in der ersten Zeile, nennt der Text keine. *Hinweis an den
  Implementer:* Fehler des Parsers sind an ihren elf Texten in `parserc.go` erkennbar, im
  genauen Vergleich, nicht am Anfang, weil Texte des Scanners ebenso mit `did not find
  expected` beginnen: `did not find expected` gefolgt von `<stream-start>`,
  `<document start>`, `node content`, `'-' indicator`, `key`, `',' or ']'` und `',' or '}'`,
  dazu `found undefined tag handle`, `found duplicate %YAML directive`,
  `found incompatible YAML document` und `found duplicate %TAG directive`; ihre Zahl plus 1. Ein
  Fehler ohne Zahl, der kein Fehler aus dem Bezug zwischen Knoten ist, liegt in Zeile 1. Der
  Kommentar daran ist eine Kopplung an die Version der Bibliothek; die Tests unten halten sie.
  *Tests* (`TestDateiUngueltig`, genaue Meldung):
  `log_level: info\nrecord:\n  force: true\nd: [x\n` nennt „Zeile 4“;
  `log_level: info\nrecord:\n  force: true\n  output: x\n  - y\n` nennt „Zeile 3“ (Abbildung
  unter `record:`); `a: 1\nb: 2\n]\n` nennt „Zeile 3“ (Grenze); `a: 1\nb: c: d\n` nennt
  „Zeile 2“ (Scanner); der Fall *Syntax* `replay:\n  listen: "GEHEIM\n` nennt „Zeile 2“;
  `]\n` nennt „Zeile 1“. *Mutation:* plus 1 entfernt, rot über `d: [x` und `]` in Zeile 3;
  plus 1 auch für den Scanner, rot über `b: c: d`; fehlende Zahl nicht als Zeile 1, rot über
  `]\n`.
- **V-119 — Meldung ohne Zeile** — jede Meldung zu ungültigem YAML nennt eine Zeile, außer
  einem Fehler aus dem Bezug zwischen Knoten: Alias ohne Anker und Anker, der sich selbst
  enthält (Grenze, die Bibliothek nennt keine Stelle; die Meldung nennt die Ursache ohne den
  Namen des Ankers). Zeichen, die YAML 1.2 nicht als druckbar zulässt (C0 außer Tabulator und
  Zeilenende, auch NUL, `U+007F`, C1, `U+FFFE`, `U+FFFF`), lehnt `kodierung` vor dem YAML mit
  der Zeile ab; damit trägt auch UTF-16LE ohne BOM eine Zeile. *Tests* (`TestDateiKodierung`):
  NUL in einem Wert in Zeile 2 nennt „Zeile 2“; `\x01` in einem Wert in Anführungszeichen in
  Zeile 3 nennt „Zeile 3“; `U+007F` und `U+0080` je mit Zeile; ein Tabulator in einem Wert in
  Anführungszeichen bleibt gültig; UTF-16LE ohne BOM nennt „Zeile 1“; `a: *x\n` nennt
  „ungültiges YAML“ mit der Ursache, ohne Zeile und ohne `x`. *Mutation:* Prüfung der nicht
  druckbaren Zeichen entfernt, rot über NUL und UTF-16LE ohne BOM (Meldung ohne Zeile); der
  Tabulator mit abgelehnt, rot. Der Kopfkommentar von `ladeDatei` sagt die Grenze zu (§3.11).

*Ort der Konstanten* — im Model, `internal/hexagon/model/fehler.go`, als drei Konstanten ohne
Logik (§1, Ausnahme). Grund: [ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) führt eine Code-Tabelle im Quelltext, die mit dem
Katalog gleich bleibt, und der Bestand legt dort schon Codes ab, die nur ein Adapter erzeugt
(`PGR-E2001` im CLI-Adapter, `PGR-E4001` im PGWire-Adapter, `PGR-E3001` im Recording-Adapter); eine zweite Tabelle im
CLI-Adapter wäre eine zweite Quelle. Keine ADR nötig: Die Entscheidung folgt aus
[ADR-0011](../../adr/0011-meldungscodes-praefix-pgr.md) und dem Bestand.

Aus `slice-v1-abschluss-konfiguration` zieht kein Risiko mit: Dessen Risiko der Einzel-Leser
bleibt dort, das zu [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md) ist entfallen (die ADR ist `Accepted`).

**Risiken:**

- Die Datei-Hälfte ist auf 900 bis 1300 Zeilen geschätzt (Bericht des Implementers vom
  2026-10-08); das liegt an der Grenze einer Review-Sitzung. Gegenmittel: der vorab benannte
  Schnitt in §4, ausgeführt am 2026-10-09 (Rest an
  `slice-v1-abschluss-verbindungen-platzhalter`) — **Ausgang:** eingetreten, Folge-Slice
  `slice-v1-abschluss-verbindungen-platzhalter`. Der Diff von DoD-Punkt 1 und 2 allein lag bei
  rund 1250 Zeilen (`7a80393`); der Schnitt nach §4 (`fdf1cd8`) gab Verbindungen und
  Platzhalter dorthin (dort §1, *Übernimmt*, DoD-Punkt 1 bis 3).
- Die Entscheidungen der zwölf Rückgaben ändern `LH-FA-17.a`; das Benutzerhandbuch beschreibt
  die Datei schon im Zielstand und kann abweichen (§3); den Teil zu Verbindungen und
  Platzhaltern trägt `slice-v1-abschluss-verbindungen-platzhalter` — **Ausgang:** entfallen.
  Begründung: Die eine Abweichung im Teil dieses Slice (§4 des Handbuchs nannte `PGR-E2005`
  für `config show`, V-120) ist in `36d8d49` behoben; §5 *Konfigurationsdatei* stimmt nach
  der Verifikation (Abschnitt 5) mit dem Binary überein, soweit dieser Slice liefert. Den Teil
  zu Verbindungen und Platzhaltern führt der Nehmer als Liefer-Punkt (dort DoD-Punkt 3), nicht
  als Risiko hier.
- Zwischen der Closure dieses Slice und der von `slice-v1-abschluss-verbindungen-platzhalter`
  nimmt der Stand eine Datei mit Klartext-Passwort an, und `config show` zeigt es; das
  verletzt die Zusage aus `LH-FA-17` (*die Anzeige zeigt sie nicht*). Gegenmittel: der Nehmer
  ist der nächste Schritt der Reihenfolge, kein Release liegt dazwischen — **Ausgang:**
  eingetreten, Folge-Slice `slice-v1-abschluss-verbindungen-platzhalter`. Mit dieser Closure
  liegt der Stand auf dem Hauptzweig; der Nehmer führt `PGR-E2006` in jeder Verbindung, auch
  bei `config show` (dort DoD-Punkt 3), und ist nach §5 von
  [welle-v1-abschluss](../welle-v1-abschluss.md) der nächste Schritt. Kein Tag und kein
  Release liegen dazwischen (Review F-516, Verifikation Abschnitt 2).

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

- **Belege zur DoD (Implementer):** siehe *Belege des Implementers* unten. DoD-Punkt 1 bis 3
  und `make gates` bestätigt die Verifikation am Stand `c5b7d96`
  (`docs/reviews/2026-10-09-verifikation-slice-v1-abschluss-konfigurationsdatei.md`,
  Abschnitt 1; 22 Mutanten rot, kein grüner). Die Nacharbeit zur Verifikation ist am Diff von
  `36d8d49` geprüft: V-118 (`ungueltigesYAML`: Zahl des Parsers plus 1, ohne Zahl Zeile 1,
  Fälle mit genauer Meldung in `TestDateiUngueltig`), V-119 (`nichtDruckbar` in `kodierung`,
  Fälle in `TestDateiKodierung`, Alias ohne Anker mit genauer Meldung) und V-120 (Handbuch
  §4: `PGR-E2004` oder `PGR-E2006`, nie `PGR-E2005`) entsprechen den Vorgaben in §6. Drei der
  sechs Mutanten aus §7 hat der Planner am Stand `36d8d49` in je einer frischen Kopie
  nachgefahren (`git archive`, Ersetzung mit Trefferzahl 1, `go test -count=1` im Image der
  Stufe `deps`, ohne Netz), alle rot aus dem genannten Grund: plus 1 entfernt
  (`TestDateiUngueltig`, *Parser, Folge* „Zeile 3“ statt „Zeile 4“, ebenso *Abbildung unter
  record* und *Grenze*), Prüfung der nicht druckbaren Zeichen entfernt
  (`TestDateiKodierung`, NUL bis `U+FFFF` „Zeile 1“ statt der Zeile des Zeichens), Fehler ohne
  Zahl ohne Zeile (`TestDateiUngueltig`, *ohne Zahl*). `make gates` am Baum von `36d8d49`
  vor dem Commit: Exit 0 (*Belege des Implementers*, letzter Absatz zur Verifikation); nach
  `36d8d49` änderten sich nur Plan und Register.
- **Was hat funktioniert:** Der vorab benannte Schnitt aus §4 griff ein zweites Mal entlang
  der benannten Linie (*Datei* gegen *Verbindungen und Platzhalter*), der Implementer hielt
  nach DoD-Punkt 2 an, und danach passte der Diff in eine Review-Sitzung (Review,
  Negativbefund *Größe*). Die zwölf Rückgaben vom 2026-10-08, vor dem Code entschieden,
  hielten: Keine davon öffnete Review oder Verifikation wieder. Der Architect entschied die
  Befunde mit Sonden an der Bibliothek (`TaggedStyle` bei `!`, UTF-16 mit BOM, Zeilen des
  Parsers ab 0), und die Nacharbeit folgte den Vorgaben ohne weitere Runde: 14 Mutanten zur
  Nacharbeit am Review und sechs zur Verifikation, alle rot. Keine Meldung zu Datei oder
  Umgebung nennt einen Wert (Verifikation Abschnitt 3), das Hexagon blieb bei drei Konstanten
  im Model.
- **Was ging anders als geplant:**
  1. Der Slice war ein zweites Mal zu groß: Schon die Schätzung beim Schnitt vom 2026-10-08
     lag an der Grenze einer Review-Sitzung (§6, Risiko 1), der zweite Schnitt war nur vorab
     benannt, nicht ausgeführt; er trat nach DoD-Punkt 1 und 2 ein (`fdf1cd8`), der Slice
     lief über `next/` zurück.
  2. Randformen im Code entschieden und erst danach zurückgegeben: R1, R2 und fünf Lesarten
     stehen im Code-Commit `7a80393`, zurückgegeben in §7, bestätigt in `cfc4d6b` (F-509).
  3. Randformen der Schicht unter dem YAML entschied erst das Review: Kodierung, Zeilenende,
     Dateiart, leerer Schlüssel (F-507); die Verifikation fand die Zeile bei Fehlern des
     Parsers (V-118) und Meldungen ohne Zeile (V-119).
  4. Ein grüner Mutant hieß äquivalent und seine Prüfung wurde entfernt, obwohl Eingaben ihn
     unterscheiden (F-506, `TaggedStyle`); ebenso bei F-511 für die Meldung.
  5. Eine Zusage der Spezifikation brach der Code ohne fangenden Test: die Meldung zu
     Doppelten nannte Text an einer Wertstelle (F-508); M3 und M4 blieben grün (F-512,
     F-513). Kommentare sagten mehr zu als die Tests (F-506, F-510, F-511), und §3 folgte dem
     Liefer-Commit `9332356` nicht (F-514).
  6. Nach der Verifikation fand die Sonde des Planners einen Zweig ohne Test: Der Kommentar
     an `ungueltigesYAML` und der Kopfkommentar von `ladeDatei` sagen für einen Anker, der
     sich selbst enthält, eine Meldung „Anker enthält sich selbst“ ohne Zeile zu. Über
     `leseYAML` ist der Zweig nicht erreichbar (`a: &x [*x]`, `a: &x {b: *x}` und die
     Blockform lesen ohne Fehler; die Datei ist danach als Anker ungültig, mit Stelle), und
     der Mutant, der seine Meldung ändert, bleibt grün. Die Spezifikation sagt nur, dass ohne
     Zeile allein diese Fälle bleiben; ihr widerspricht der Stand nicht. **Ausgang: offen,
     Entscheidung des Nutzers vor dem `git mv`** (Nacharbeit hier: Zweig und Zusage enger,
     oder Randform an `slice-v1-abschluss-verbindungen-platzhalter`, der `datei.go` ohnehin
     ändert).
  - **Summary-Zeilen:** Review
    `docs/reviews/2026-10-09-review-slice-v1-abschluss-konfigurationsdatei.md`: „0 HIGH · 4
    MEDIUM · 6 LOW · 3 INFO (F-506 die Tag-Prüfung über die Textstelle lässt Tags bei BOM,
    UTF-16, `\r` und `U+2028` durch, die entfernte `TaggedStyle`-Prüfung war nicht
    äquivalent; F-507 Kodierung, Zeilenende und nicht reguläre Dateien sind nirgends
    entschieden; F-508 die Meldung zu Doppelten nennt Text an Wertstellen; F-509 R1, R2 und
    die Lesarten wurden im Code entschieden und erst danach zurückgegeben; F-510 bis F-513
    Zusagen weiter als die Prüfung bzw. Mutanten grün; F-514 §3 nennt gelieferte Tests
    „Offen“; F-515 TLS-Adresse mit verschiedenem Geber; F-516 bis F-518 Zwischenrisiko
    trägt, Encoder unter [ADR-0036](../../adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), Priorität unabhängig getestet). Wiederkehrende Klassen:
    `BEO-REPO/spec-randform-erst-im-review-entschieden` (F-507),
    `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (F-509),
    `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (F-508, F-512, F-513),
    `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-506, F-510, F-511),
    `BEO-REPO/plan-folgt-korrektur-nicht` (F-514).“ Verifikation
    `docs/reviews/2026-10-09-verifikation-slice-v1-abschluss-konfigurationsdatei.md`: „0 HIGH ·
    1 MEDIUM · 2 LOW · 1 INFO (V-118 die Zeile zu ungültigem YAML ist bei Parser-Fehlern die
    Zeile vor dem Konstrukt; V-119 C0-Steuerzeichen, UTF-16LE ohne BOM und ein Alias ohne
    Anker ergeben eine Meldung ohne Zeile; V-120 Handbuch §4 nennt `PGR-E2005` für `config
    show`; V-121 ein Verbindungsname in der Form einer URL trägt ein Passwort in Meldung und
    Anzeige, Randform für den Nehmer).“ Verdikt: DoD-Liefer-Punkte 1 bis 3 und `make gates`
    bestätigt. Ausgänge: F-506 bis F-508, F-510 bis F-514 in `c5b7d96`, F-517 als
    [ADR-0037](../../adr/0037-yaml-bibliothek-fuer-die-anzeige-der-konfigurationsdatei.md)
    (`5cbfd1b`, `8d320cc`), F-515 in §1, F-509 unter *Beobachtungs-Register*; V-118 und V-119
    vom Architect in `f10a397`, mit V-120 umgesetzt in `36d8d49`; V-121 als offene Randform in
    §6 von `slice-v1-abschluss-verbindungen-platzhalter` (`f10a397`).
- **Steering-Loop-Eintrag:** Benannte Spec-Lücke, in `LH-FA-17.a` geschlossen: Die
  Spezifikation regelte die Datei auf der Ebene des YAML (Text des Skalars, Tags, Anker,
  Doppelte), nicht die Schicht darunter, auf der die Bibliothek still entscheidet:
  Kodierung und BOM, Zeilenenden und nicht druckbare Zeichen, die Dateiart vor dem Öffnen,
  und welche Zeile eine Meldung nennt, wenn die Bibliothek anders oder gar nicht zählt.
  Daraus kamen F-506, F-507, V-118 und V-119. Der Architect schloss die Lücke vor der
  Nacharbeit (`5cbfd1b`: Kodierung, Zeilenende, reguläre Datei, leerer Schlüssel; `f10a397`:
  Zeile bei einem Konstrukt über mehrere Zeilen, nicht druckbare Zeichen, Meldung ohne
  Zeile). Dazu eine neue Beobachtung: `BEO-REPO/gruener-mutant-faelschlich-aequivalent`
  (F-506, F-511, 1×). Kein Feld `liegt in`: Mit diesem Slice ist nichts verkörpert; die
  Spec-Stellen tragen keinen Herkunfts-Anker.

  Retirement-Checks: `AGENTS.md` §3.9 (seit welle-walking-skeleton) ist wieder aufgetreten
  (F-514); die Regel bleibt. §3.10 (seit welle-extended-query) ist wieder aufgetreten
  (F-508, F-512, F-513); die Regel bleibt, der Sensor ist mit `slice-harness-mutation`
  geplant. §3.11 (seit welle-extended-query) ist wieder aufgetreten (F-506, F-510, F-511 und
  der Zweig aus Punkt 6 oben); die Regel bleibt. §3.12 (seit slice-harness-randformen-vor-code)
  ist wieder aufgetreten, beide Hälften: Randformen, die §6 nicht nannte, fand das Review
  (F-507) und die Verifikation (V-118, V-119), und R1, R2 und die Lesarten entschied der
  Code vor der Rückgabe (F-509), obwohl `implement-slice` die Randform-Rückgabe vor dem Code
  verlangt; die Regel bleibt. §3.13 (seit slice-lint-bestand-kern-driven) ist nicht wieder
  aufgetreten: Jede Adresse in §1 nimmt an (F-515 ist LOW, die Adresse nimmt an, nur der
  Geber hieß verschieden; behoben in `c5b7d96`). `implement-slice` Schritt 19
  (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`): nicht wieder aufgetreten, jeder Mutant
  in einer frischen Kopie per Bind-Mount.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `36d8d49` (Zähler =
  Dateien unter `evidence/`).
  - `BEO-REPO/spec-randform-erst-im-review-entschieden`: **Beleg**, 15× → 16× (F-507, V-118,
    V-119), Stand verkörpert, bleibt.
  - `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`: **Beleg**, 4× → 5× (F-509),
    Stand verkörpert, bleibt.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: **Beleg**, 16× → 17× (F-508, F-512,
    F-513), Stand verkörpert, bleibt. F-508 zählt hier, nicht unter `zusage-…`: Die Zusage
    „nie einen Wert“ steht in der Spezifikation, nicht in einem Kommentar, Hilfetext oder
    Plan; der Code brach sie, und kein Test fing es. Eine eigene Klasse *Meldung nennt einen
    Wert* trägt der Bestand nicht: Derselbe Inhalt bei V-121 ist vor dem Code als Randform
    beim Nehmer, nicht als Befund.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: **Beleg**, 22× → 23× (F-506, F-510,
    F-511; dazu der Zweig aus Punkt 6 oben), Stand verkörpert, bleibt.
  - `BEO-REPO/plan-folgt-korrektur-nicht`: **Beleg**, 18× → 19× (F-514), Stand verkörpert,
    bleibt.
  - `BEO-REPO/gruener-mutant-faelschlich-aequivalent`: **neu**, 1× (F-506, F-511), offen.
    Nicht unter `BEO-REPO/liste-gruener-mutanten-unvollstaendig`: Dort sind die grünen
    Mutanten eines Umbaus im Bestand richtig eingeordnet, aber nicht alle gefunden; hier ist
    ein gefundener Mutant im neuen Code falsch eingeordnet.
  - `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`: **neu**, 1× (der zweite Schnitt nach §4,
    `fdf1cd8`), offen.
  - `BEO-REPO/slice-waechst-durch-uebernahmen`: **kein Beleg**, bleibt 2× (offen). Dieser
    Slice nahm nach seiner Anlage keine Übernahme auf; die drei aus §1 brachte der Schnitt
    von `slice-v1-abschluss-konfiguration` mit, und dessen Beleg zählt sie schon. Zu groß
    machte ihn der Kern der Datei selbst (DoD-Punkt 1 und 2 allein rund 1250 Zeilen) bei
    einer Größe, die beim Schnitt schon geschätzt war; das zählt unter
    `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`.
  - Ohne Beleg: `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel` (bleibt 1×, die
    Rückführung lief über `next/`), `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (bleibt
    3×, verkörpert; alle Nehmer nehmen an), `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht`
    (bleibt 3×, verkörpert; die Belege in §7 erreichten beide Prüfer),
    `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (bleibt 2×, siehe oben),
    `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (bleibt 1×; kein Gate geändert),
    `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (bleibt 1×).

  Einmalig und nicht eingetragen: F-515 (Geber der Sendung verschieden genannt, Adresse
  nimmt an), F-516 bis F-518 (Hinweise und Negativbefund), V-120 (Handbuch, Risiko 2 aus §6),
  V-121 (Randform beim Nehmer). Mit diesem Slice erreicht kein Eintrag die Schwelle 3× neu.
  Über der Schwelle stehen nur Einträge mit Ausgang (`commit-nennt-struktur-kennung` 3×
  geplant; `folge-slice-adresse-nimmt-nicht-an` 3×,
  `implementer-bericht-erreicht-pruefer-nicht` 3×,
  `randform-im-code-entschieden-dann-zurueckgegeben` 5×,
  `spec-randform-erst-im-review-entschieden` 16×, `negativtests-fehlen-bei-neuem-vertrag`
  17×, `plan-folgt-korrektur-nicht` 19×, `zusage-im-kommentar-weiter-als-pruefung` 23×,
  verkörpert; `roter-lauf-haengt-bis-zum-zeitlimit` 3× geplant;
  `white-box-liste-vor-code-nur-namenssuche` 3× gestrichen).
- **Folge-Slices:** `slice-v1-abschluss-verbindungen-platzhalter` (früherer DoD-Punkt 3,
  Risiko 1 und 3 aus §6, V-121 als offene Randform; der nächste Schritt nach §5 von
  [welle-v1-abschluss](../welle-v1-abschluss.md)). Aus §1 *Ausdrücklich NICHT*:
  `slice-v1-abschluss-schreiben` (atomares Schreiben, vorhandenes `--output`),
  `slice-v1-abschluss-einspielen` (Optionen von `play`, Abschnitt `play:`). Der allgemeine
  Leser liegt in `slice-v1-abschluss-konfiguration` (`done/`); er ist Geber, kein
  Folge-Slice.
- **Risiken aus §6:** drei. *Größe der Datei-Hälfte*: **eingetreten**, Folge-Slice
  `slice-v1-abschluss-verbindungen-platzhalter`. *Handbuch weicht ab*: **entfallen**,
  Begründung in §6. *Klartext-Passwort zwischen den Closures*: **eingetreten**, Folge-Slice
  `slice-v1-abschluss-verbindungen-platzhalter`. Die Randformen in §6 sind Entscheidungen,
  keine Risiken.
- **Drei Paarungen:** Anker: Der Steering-Loop-Eintrag trägt kein Feld `liegt in`; nichts
  zu prüfen. Folge-Slice: `grep -n "slice-v1-abschluss-konfigurationsdatei"` findet die
  Kennung in §1 jedes Nehmers — `slice-v1-abschluss-verbindungen-platzhalter` (`next/`, §1
  *Übernimmt*: Grammatik der URL, `sslmode`, Platzhalter, `PGR-E2005` und `PGR-E2006`;
  DoD-Punkt 3 trägt `PGR-E2006` auch bei `config show`; §6 die Randform zu V-121),
  `slice-v1-abschluss-schreiben` (`next/`, §1 *Abgegeben* und *Ausdrücklich NICHT*: Datei
  hier, atomares Schreiben und vorhandenes `--output` dort im Ziel),
  `slice-v1-abschluss-einspielen` (`next/`, §1 *Übernommen aus
  `slice-v1-abschluss-konfigurationsdatei`*: Abschnitt `play:`); keiner schließt die Sendung
  in seinem §1 *Ausdrücklich NICHT* aus, keiner liegt in `done/`. Register:
  `BEO-REPO/spec-randform-erst-im-review-entschieden`,
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`,
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`,
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/plan-folgt-korrektur-nicht`,
  `BEO-REPO/gruener-mutant-faelschlich-aequivalent` und
  `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` tragen
  `evidence/slice-v1-abschluss-konfigurationsdatei.md`; die übrigen genannten Einträge
  bestehen als Verzeichnis, jedes mit nicht leerem `evidence/`. Die nächste Welle-Closure
  prüft erneut.

### Belege des Implementers

**Stand.** Geliefert sind DoD-Punkt 1 und 2 (*Datei*). DoD-Punkt 3 (*Verbindungen und
Platzhalter*) ist nicht begonnen: Der Diff der beiden ersten Punkte umfasst rund 1250
Zeilen (davon rund 600 Tests), mit DoD-Punkt 3 läge er über einer Review-Sitzung. Damit
ist die Bedingung der Rückführung *zu groß* aus §4 eingetreten, mit dem dort vorab
benannten Schnitt; der Implementer hält nach DoD-Punkt 2 an, den Übergang entscheidet der
Planner. Bis DoD-Punkt 3 gilt im gelieferten Stand: `connections:` ist eine Abbildung, ein
Name hat die Form eines Werts ohne Steuerzeichen, ein Wert ist ein Skalar; die Grammatik
der URL, `sslmode`, `${VAR}` und `$$` (in jedem Wert der Datei), das Klartext-Passwort
(`PGR-E2006`), `PGR-E2005` und der Name bei `--upstream` sind nicht umgesetzt, ein Wert
mit `$` gilt wörtlich. Die Konstanten `PGR-E2005` und `PGR-E2006` stehen schon in
`internal/hexagon/model/fehler.go`, ohne Erzeuger.

*Nachtrag des Planners vom 2026-10-09:* Mit dem Schnitt nach §4 heißt DoD-Punkt 3 dieses Plans
die Tests der Vorgaben vom 2026-10-09; der frühere DoD-Punkt 3, von dem dieser Absatz spricht,
ging an `slice-v1-abschluss-verbindungen-platzhalter` (§1, *Abgegeben*). Die Belege unten
gelten DoD-Punkt 1 und 2.

**Weg der Mutanten.** Je Mutant eine frische Kopie des Arbeitsbaums in einem neuen
Verzeichnis (`tar` ohne `.git`, kein `cp -p`), die Änderung mit einem Skript, das die
Trefferzahl 1 prüft, dann `go test -count=1` für `internal/adapters/driving/cli` und
`internal/bootstrap` im Image der Stufe `test`, die Kopie per Bind-Mount statt
Build-Kontext; danach wird nur das Verzeichnis des Mutanten gelöscht. Der Arbeitsbaum
blieb unberührt. Gefahren am Stand `7a80393`: 66 Mutanten, 65 rot, einer grün; dieser ist unter
DoD-Punkt 3 rot gesehen.

**DoD-Punkt 1 — Hilfe vor jeder Prüfung, `config show`.**

| Zusage | Mutation | rote Tests |
|---|---|---|
| Hilfe vor der Datei, auch bei `--config` und `config show` | `Parse` lädt die gewählte Datei vor der Hilfe-Suche | `TestDateiHilfeVorPruefung`, `TestDateiNichtLesbar`, `TestDateiWahl` |
| `config --help` ist die Hilfe von `config show` | Hilfe von `config` entfernt (globale Hilfe) | `TestDateiHilfeVorPruefung` |
| globale Hilfe nennt `config show` | Zeile und Abschnitt entfernt | `TestDateiHilfeVorPruefung` |
| Hilfe von `record` und `replay` nennt `--config` und die Datei | `optionConfig` aus der Hilfe von `record` genommen | `TestDateiHilfeVorPruefung` |
| `--` beendet die Optionen auch bei `config show` | `config show -- x` überspringt `--` | `TestDateiHilfeVorPruefung` |
| erste Zeile: Pfad der gewählten Datei | feste Zeile statt Pfad | `TestConfigShow`, `TestConfigShowOhneDatei`, `TestRunConfigShow` |
| ohne Datei sagt es die erste Zeile | Zeile weggelassen | `TestConfigShowOhneDatei` |
| Inhalt ohne Kommentare | Kommentare nicht entfernt | `TestConfigShow`, `TestRunConfigShow` |
| zwei Leerzeichen Einzug | Einzug 4 | `TestConfigShow` |
| Umgebungsvariablen nach Namen sortiert | ohne Sortierung | `TestConfigShow` |
| Umgebungsvariablen ohne Wert | `NAME=wert` ausgegeben | `TestConfigShow`, `TestConfigShowOhneDatei` |
| nur gesetzte und nicht leere | leere mitgezählt | `TestConfigShow`, `TestConfigShowOhneDatei`, `TestRunConfigShow` |
| auch ohne passende Option | nur Variablen mit Option | `TestConfigShow` |
| ungültige Datei: Fehlercode des Ladens, keine Anzeige | Ladefehler übergangen, Anzeige ohne Datei | `TestConfigShowFehler`, `TestRunConfigShow` |
| von den Umgebungsvariablen prüft `config show` nur `PGWIRE_RECORDER_CONFIG` | Umgebungsvariablen von `replay` mitgeprüft | `TestConfigShowOhneDatei` |
| `config show` kennt nur `--config` | Optionen von `replay` angemeldet | `TestConfigShowFehler` |
| `config` ohne `show` ist `PGR-E2001` | `config` allein als `config show` | `TestConfigShowFehler` |
| Ausgabe auf `stdout`, nichts auf `stderr` | Anzeige auf `stderr` | `TestRunConfigShow` |

**DoD-Punkt 2 — Wahl der Datei, Priorität, Schlüssel, ungültige Datei.**

| Zusage | Mutation | rote Tests |
|---|---|---|
| `--config` vor `PGWIRE_RECORDER_CONFIG` | Reihenfolge getauscht | `TestDateiWahl` |
| `PGWIRE_RECORDER_CONFIG` wählt die Datei | Variable nie gelesen | `TestConfigShowFehler`, `TestConfigShowOhneDatei`, `TestDateiNichtLesbar`, `TestDateiWahl` |
| Standarddatei nur, wenn sie existiert | fehlende Standarddatei gewählt | `TestDateiWahl`, `TestConfigShowOhneDatei`, `TestLeserAlleOptionen` und weitere Tests des Lesers |
| leere `PGWIRE_RECORDER_CONFIG` gilt als nicht gesetzt | `LookupEnv` statt nicht leer | `TestDateiWahl`, `TestDateiNichtLesbar`, `TestConfigShowOhneDatei`, `TestRunConfigShow`, `TestRunReplayAusDatei` |
| genannte Datei fehlt: `PGR-E2004` | fehlende Datei setzt nichts | `TestDateiNichtLesbar`, `TestConfigShowFehler` |
| Standarddatei vorhanden, nicht lesbar: `PGR-E2004` | Fehler von `Stat` als *keine Datei* | `TestDateiNichtLesbar` |
| genau eine Datei: neben `--config` wird die aus der Umgebung nicht gelesen | fehlende Datei aus der Umgebung neben `--config` gemeldet | `TestDateiWahl` |
| Meldung nennt die Quelle, nicht den Pfad | Grund mit Pfad des Betriebssystems | `TestDateiNichtLesbar` |
| Umgebungsvariable vor Datei | Datei überschreibt die Umgebung | `TestLeserAlleOptionen` |
| Kommandozeile vor Datei | Datei überschreibt die Kommandozeile | `TestLeserAlleOptionen` |
| Datei vor Standardwert | Wert der Datei nie übernommen | `TestLeserAlleOptionen`, `TestDateiGueltig`, `TestDateiWahl`, `TestRunReplayAusDatei` |
| Pflichtoption auch aus der Datei | Pflicht nur aus Kommandozeile oder Umgebung | `TestLeserAlleOptionen`, `TestDateiGueltig`, `TestRunReplayAusDatei` |
| Wert der Datei mit der Wertemenge der Option geprüft | Prüfung entfernt | `TestLeserAlleOptionen`, `TestDateiUngueltig`, `TestDateiReihenfolge`, `TestConfigShowFehler`, `TestRunConfigShow` |
| ungültiger Wert der Datei ist `PGR-E2004`, nicht `PGR-E2001` | Code `PGR-E2001` | dieselben fünf |
| geprüft, auch wenn die Kommandozeile die Option setzt | Ladefehler übergangen, wenn die Kommandozeile setzt | `TestLeserAlleOptionen`, `TestDateiUngueltig`, `TestDateiNichtLesbar` |
| Meldung nennt den Schlüssel, nie den Wert | Wert in die Meldung | `TestLeserAlleOptionen`, `TestDateiUngueltig` |
| ungültiges YAML ohne Text der Bibliothek | Fehler der Bibliothek angehängt | `TestDateiUngueltig` (Fall *Alias ohne Anker*) |
| Code ist `PGR-E2004` | Konstante `PGR-E2003` | `TestRunConfigShow`, `TestRunReplayAusDatei` |
| `log_level` auf der obersten Ebene, auch bei `record` | `oben` an `log-level` von `record` entfernt | `TestLeserAlleOptionen` |
| `log_level` in einem Abschnitt ungültig | `log_level` im Abschnitt angenommen | `TestLeserAlleOptionen`, `TestDateiUngueltig`, `TestConfigShowFehler` |
| Schlüssel mit `_` statt `-` | Name unverändert | 7 Tests, darunter `TestLeserAlleOptionen`, `TestConfigShow` |
| unbekannter Schlüssel oben | angenommen | `TestLeserAlleOptionen`, `TestDateiUngueltig`, `TestConfigShowFehler` |
| unbekannter Schlüssel im Abschnitt | übergangen | `TestLeserAlleOptionen`, `TestDateiUngueltig`, `TestConfigShowFehler`, `TestRunReplayAusDatei` |
| Abschnitt nur mit Optionen seines Kommandos | Optionen beider Kommandos | `TestDateiUngueltig`, `TestConfigShowFehler` |
| kein Schlüssel `config` | `config` im Abschnitt angenommen | `TestDateiUngueltig`, `TestConfigShowFehler` |
| Abschnitt `play:` unbekannt (§1) | `play` als Abschnitt angenommen | `TestDateiUngueltig` (Fall `play: {}`), `TestConfigShowFehler` |
| Groß- und Kleinschreibung zählt | Schlüssel klein verglichen | `TestDateiUngueltig`, `TestConfigShowFehler` |
| doppelter Schlüssel ist ungültiges YAML | nicht erkannt | `TestDateiUngueltig`, `TestConfigShowFehler` |
| doppelter Schlüssel vor den Schlüsseln geprüft (Phase *YAML*) | Prüfung nach den Schlüsseln | `TestDateiUngueltig` (Fall *doppelt vor unbekanntem Schlüssel*) |
| doppelter Schlüssel in jeder Tiefe | nur oberste Ebene | `TestDateiUngueltig`, `TestConfigShowFehler` |
| höchstens ein Dokument | zweites übergangen | `TestDateiUngueltig`, `TestConfigShowFehler` |
| oberste Ebene ist eine Abbildung | jede Form angenommen | `TestDateiUngueltig`, `TestConfigShowFehler` |
| Abschnitt ohne Inhalt oder `null` ungültig | `null` als Abschnitt angenommen | `TestDateiUngueltig`, `TestConfigShowFehler` |
| `connections:` ist eine Abbildung | jede Form angenommen | `TestDateiUngueltig`, `TestConfigShowFehler` |
| leere Abbildung `{}` gültig | leere Abbildung abgelehnt | `TestDateiGueltig` |
| leere Datei und nur Kommentare setzen nichts | leere Datei als Fehler | `TestDateiGueltig`, `TestConfigShowOhneDatei` |
| jeder Tag ungültig, auch `!` | nur Tags der Bibliothek (`TaggedStyle`) | `TestDateiUngueltig`, `TestConfigShowFehler` |
| Tag ungültig | Prüfung des Tags entfernt | `TestDateiUngueltig`, `TestConfigShowFehler` |
| Anker ungültig | Prüfung entfernt | `TestDateiUngueltig`, `TestConfigShowFehler` |
| `null`, `~` und leerer Wert ungültig | nur leerer Wert | `TestDateiUngueltig`, `TestConfigShowFehler` |
| Liste oder Abbildung als Wert ungültig | Prüfung entfernt | `TestDateiUngueltig`, `TestConfigShowFehler` |
| Name einer Verbindung nicht leer | Prüfung entfernt | `TestDateiUngueltig`, `TestConfigShowFehler` |
| Name einer Verbindung ohne Steuerzeichen | Prüfung entfernt | `TestDateiUngueltig`, `TestConfigShowFehler` |
| Wert einer Verbindung ist ein Skalar | nur Anker und Tag geprüft | `TestDateiUngueltig`, `TestConfigShowFehler` |
| Kommandozeile vor Datei geprüft | Ladefehler vor dem Fehler der Kommandozeile | `TestDateiReihenfolge`, `TestDateiNichtLesbar` |
| Umgebungsvariablen vor Datei geprüft | Datei vor den Umgebungsvariablen geladen | `TestDateiReihenfolge` |
| Datei in ihrer Reihenfolge, oberste Ebene | rückwärts | `TestDateiReihenfolge` |
| Datei in ihrer Reihenfolge, im Abschnitt | rückwärts | `TestDateiReihenfolge`, `TestDateiUngueltig` |
| relativer Pfad relativ zum aktuellen Verzeichnis | Pfad umgeschrieben | `TestDateiGueltig` |

**DoD-Punkt 3 — Vorgaben des Architect vom 2026-10-09** (§6, `cfc4d6b`). Gleicher Weg
der Mutanten, gefahren am Stand des Commits, der diese Tests liefert: fünf Mutanten,
alle rot.

| Zusage | Mutation | rote Tests |
|---|---|---|
| R1: `"null"` und `'~'` in Anführungszeichen sind Text, `output` setzt den Pfad `null` bzw. `~` | Text `null` oder `~` auch in Anführungszeichen als `null` (der zuvor grüne Mutant) | `TestDateiGueltig` |
| R1: `null` ist die YAML-Null ohne Anführungszeichen, auch `Null` und `NULL` | Prüfung auf `!!null` durch Textvergleich mit leer, `null` und `~` ersetzt | `TestDateiUngueltig`, `TestConfigShowFehler` (Fälle `Null`, `NULL`) |
| R2: Standarddatei als Link auf ein fehlendes Ziel wählt keine Datei, `config show` sagt das in der ersten Zeile | `os.Lstat` statt `os.Stat` | `TestConfigShowOhneDatei` |
| Lesart 3: erste Zeile von `config show` ist der Pfad aus `--config`, auch relativ | Pfad mit `filepath.Abs` | `TestConfigShowOhneDatei` |
| Lesart 2: ein unbekanntes Unterkommando von `config` ist `PGR-E2001` | jedes Unterkommando als `show` | `TestConfigShowFehler` (Fall `{"config", "zeige"}`) |

**Grüne Mutanten, eingeordnet.**

- *Schlüssel ist kein Skalar* (eine eigene Meldung für einen Schlüssel, der Liste oder
  Abbildung ist): **äquivalent im Code** — ein solcher Schlüssel hat den leeren Text, keine
  Option und kein Abschnitt heißt so, er ist unbekannt (`PGR-E2004`). Die Meldung nannte
  dabei keine Stelle (F-511); seit der Nacharbeit nennt sie die Abbildung (`stelleVon`,
  Zeile *F-511* unten).
- *Alias ungültig* (eigene Prüfung): **äquivalent** — ein Alias verweist auf einen Anker
  davor in der Datei, und den lehnt die Prüfung in der Reihenfolge der Datei zuerst ab;
  ein Alias ohne Anker ist ungültiges YAML. Die Prüfung ist entfernt, der Fall *Alias* in
  `TestDateiUngueltig` bleibt.
- *Tag über `TaggedStyle`* neben dem ersten Zeichen `!`: **nicht äquivalent**, die
  Einstufung als äquivalent war falsch (F-506). Die Prüfung am Text hängt an Zeile und
  Spalte der Bibliothek, und die zählte bei BOM, `\r`, `U+2028`, `U+0085` und UTF-16 anders
  als `zeilen`. Die Prüfung über `TaggedStyle` ist zurück und ergänzt die am Text; der
  Mutant ist rot (Zeile *F-506* unten).

**Befunde des Reviews vom 2026-10-09** (§6, Vorgaben 1 bis 8, und F-512, F-513). Gleicher Weg
der Mutanten, gefahren am Stand des Commits, der die Nacharbeit liefert: 14 Mutanten, alle
rot; der E2E-Mutant (Prüfung der regulären Datei entfernt) per `make test-integration` in
einer frischen Kopie.

| Zusage | Mutation | rote Tests |
|---|---|---|
| F-506: ein Tag, den die Bibliothek als `TaggedStyle` markiert, ist ungültig, auch ohne Text an Zeile und Spalte | Prüfung über `TaggedStyle` entfernt | `TestDateiTag` (`form` direkt) |
| F-506: `!` allein mit führendem BOM ist ungültig | BOM vor dem Zählen nicht entfernt | `TestDateiTag` |
| F-506: Zeilenende `\r` allein zählt; `!` in einem Kommentar macht die Datei nicht ungültig | nur nach `\n` gezählt | `TestDateiTag` |
| F-507 (a): BOM nach dem ersten Zeichen ist `PGR-E2004` mit der Zeile | Prüfung entfernt | `TestDateiKodierung` |
| F-507 (a): nur UTF-8, UTF-16 mit BOM ist `PGR-E2004` | Prüfung auf UTF-8 entfernt | `TestDateiKodierung` |
| F-507 (b): `U+0085`, `U+2028`, `U+2029` sind `PGR-E2004` mit der Zeile | Prüfung der drei Zeichen entfernt | `TestDateiKodierung` |
| F-507 (c): nur eine reguläre Datei, geprüft vor dem Öffnen | Prüfung nur auf Verzeichnis | `TestDateiKeineRegulaere` (`/dev/null`, FIFO über die Frist), `TestE2EConfigShowStdin` (Pipe über die Frist von 15 s) |
| F-507 (c): Links gefolgt | `os.Lstat` statt `os.Stat` | `TestDateiKeineRegulaere` (Link auf eine reguläre Datei) |
| F-508: Meldung zu einem doppelten Schlüssel nennt die Zeile, nie den Schlüssel | Text des Schlüssels in der Meldung | `TestDateiUngueltig` (Wertstelle mit `GEHEIM`) |
| F-510: Doppelte auch in einer Liste (M1) | kein Abstieg in Listen | `TestDateiUngueltig` |
| F-510: Doppelte in jeder Tiefe (M2) | nur bis zur zweiten Ebene | `TestDateiUngueltig` |
| F-511: leerer oder nicht skalarer Schlüssel nennt die Abbildung als Stelle | Stelle aus dem Text des Schlüssels | `TestDateiUngueltig` |
| F-512: Name einer Verbindung ohne DEL und C1 (M3) | nur Steuerzeichen unter `0x20` | `TestDateiUngueltig`, `TestConfigShowFehler` |
| F-513: `config show` listet nur Variablen mit dem Präfix am Anfang (M4) | `strings.Contains` statt `strings.HasPrefix` | `TestConfigShow` |

Ohne Mutation, wie §6 sagt: die Ausgabe von `config show` ohne BOM und mit `\n`
(`TestConfigShowOhneBOM`); die Zusage hält der Kodierer. Die Kommentare an `form`,
`zeilen`, `kodierung` und `doppelte` sagen nur zu, was diese Tests prüfen (Vorgabe 8).

Läufe zur Nacharbeit, am Baum des Liefer-Commits vor dem Commit: `make test` grün,
`make lint` Exit 0, `make abdeckung` geschrieben, `make gates` Exit 0 (darin
`d-check: 382 Datei(en) geprüft, 0 Befund(e)`, `run-integration-tests: gruen` mit
`TestE2EConfigShowStdin`, `a-check-negativ: gruen`, `abdeckung-gegenprobe: gruen`,
`kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `baseline-verify: v6.16.0 OK`).

**Befunde der Verifikation vom 2026-10-09** (§6, V-118, V-119; V-120 Handbuch §4). Gleicher
Weg der Mutanten, gefahren am Stand des Commits, der die Nacharbeit liefert: sechs Mutanten,
alle rot; die 14 Mutanten der Nacharbeit zum Review dort erneut gefahren, alle rot.

| Zusage | Mutation | rote Tests |
|---|---|---|
| V-118: Fehler des Parsers (elf Texte aus `parserc.go`) nennen ihre Zahl plus 1 | plus 1 entfernt | `TestDateiUngueltig` (Fälle `d: [x` → „Zeile 4“, `]` in Zeile 3) |
| V-118: Fehler des Scanners nennen ihre Zahl | plus 1 auch beim Scanner | `TestDateiUngueltig` (Fall `b: c: d` → „Zeile 2“) |
| V-118: ein Fehler ohne Zahl liegt in Zeile 1 | ohne Zahl keine Zeile | `TestDateiUngueltig` (Fall `]` → „Zeile 1“) |
| V-119: nicht druckbare Zeichen (NUL, `U+0001`, `U+007F`, `U+0080`, `U+FFFE`, `U+FFFF`) mit Zeile, auch UTF-16LE ohne BOM in Zeile 1 | Prüfung entfernt | `TestDateiKodierung` |
| V-119: ein Tabulator bleibt gültig | Tabulator mit abgelehnt | `TestDateiKodierung` |
| V-119: Alias ohne Anker nennt die Ursache ohne Zeile und ohne Namen | Name in die Meldung | `TestDateiKodierung`, `TestDateiUngueltig` |

Die Fälle zu ungültigem YAML in `TestDateiUngueltig` prüfen die Meldung genau. Der Kommentar
an `parserFehler` ist eine Kopplung an die Version v3.0.5 der Bibliothek, der Kopfkommentar
von `ladeDatei` nennt die Grenze (Alias ohne Anker, Anker, der sich selbst enthält). Handbuch
§4 (V-120): `config show` meldet `PGR-E2004` oder `PGR-E2006`, nie `PGR-E2005`. Läufe am Baum
des Liefer-Commits vor dem Commit: `make test` grün, `make lint` Exit 0, `make abdeckung`
geschrieben, `make gates` Exit 0 (darin `d-check: 383 Datei(en) geprüft, 0 Befund(e)`,
`run-integration-tests: gruen`, `a-check-negativ: gruen`, `abdeckung-gegenprobe: gruen`,
`kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `baseline-verify: v6.16.0 OK`).

**Rückgaben und Lesarten** (an den Architect, ohne Eintrag in §6; entschieden in
`cfc4d6b`, §6 *Rückgaben vom 2026-10-09*, belegt unter DoD-Punkt 3 oben):

- *R1 — Wert `"null"` oder `"~"` in Anführungszeichen:* Text des Skalars (gültig, wo die
  Wertemenge ihn annimmt) oder `null` (`PGR-E2004`)? Entschieden: Text.
- *R2 — Standarddatei als symbolischer Link ins Leere:* Der Code folgt `os.Stat`: Ein
  Link ins Leere ist *keine Datei*, eine Schleife *vorhanden, aber nicht lesbar*
  (`PGR-E2004`, `TestDateiNichtLesbar`). Entschieden: ein Link ins Leere ist keine Datei.
- *Lesarten aus `LH-FA-01.a` und `LH-FA-17.a`*, im Code so umgesetzt und getestet:
  `config` ist das bekannte Kommando der Hilfe, `config --help` gibt die Hilfe von
  `config show`; `config` ohne `show`, ein anderes Unterkommando und ein Argument nach
  `config show` sind `PGR-E2001` (`TestConfigShowFehler`); die erste Zeile von
  `config show` ist der Pfad, wie er gewählt wurde (relativ bleibt relativ); eine Datei
  nur mit `---` ist weder leer noch nur Kommentar, ihre oberste Ebene ist keine
  Abbildung (`PGR-E2004`); doppelte Schlüssel werden nach ihrem Text verglichen (`1` und
  `"1"` sind derselbe Schlüssel).

**Läufe:** `make gates` am Inhalt von `7a80393` (Lauf unmittelbar vor dem Commit, Baum
gleich), Exit 0, darin `d-check: 379 Datei(en) geprüft, 0 Befund(e)`,
`baseline-verify: v6.16.0 OK`, `a-check-negativ: gruen`, `run-integration-tests: gruen`
(darin `TestE2EConfigShow`, `TestE2EReplayKonfigurationsdateiUngueltig`),
`abdeckung-gegenprobe: gruen`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`.
Dazu `make lint` Exit 0 (`0 issues.`), `make abdeckung` geschrieben,
`make abdeckung-check` grün im Gate-Lauf.

Für DoD-Punkt 3: `make test` grün, `make lint` Exit 0, `make abdeckung` geschrieben und
`make gates` Exit 0 am Baum des Commits, der die Tests liefert, vor dieser Zeile (darin
`d-check: 380 Datei(en) geprüft, 0 Befund(e)`, `run-integration-tests: gruen`,
`a-check-negativ: gruen`, `abdeckung-gegenprobe: gruen`, `kopf-check-gegenprobe: gruen`,
`lint-gegenprobe: gruen`, `baseline-verify: v6.16.0 OK`).

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `a968790` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/slice-waechst-durch-uebernahmen` (1×) — `slice-v1-abschluss-konfiguration` wuchs
  durch Übernahmen bei drei Liefer-Punkten, bis die Datei-Hälfte nicht mehr in eine
  Review-Sitzung passte; dieser Slice entsteht aus dem Schnitt. Der Beleg gehört in die
  Closure jenes Slice (dann 2×); hier ein Risiko in §6.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (14×, verkörpert in `AGENTS.md` §3.12)
  und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (4×, verkörpert) — der
  Implementer hielt an und gab zwölf Randformen zurück, statt sie im Code zu entscheiden;
  der Architect entschied sie vor dem Code (§6).
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (15×, `AGENTS.md` §3.10) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (21×, §3.11) — je Zusage eine
  Mutation, Beleg in §7.
- `BEO-REPO/plan-folgt-korrektur-nicht` (17×, §3.9) — dieser Plan entsteht aus einem
  Schnitt; jede Korrektur zieht §1, §3 und §6 im selben Commit nach.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — die Übernahme oben steht mit
  der Kennung des Gebers; Geber und die Nehmer `slice-v1-abschluss-schreiben` und
  `slice-v1-abschluss-einspielen` zeigen im selben Commit hierher. Beim Schnitt vom 2026-10-09
  nennt der Nehmer `slice-v1-abschluss-verbindungen-platzhalter` diesen Slice unter
  *Übernimmt*, und `slice-v1-abschluss-einspielen` nennt beide unter *Übernommen aus*.

*Nachsichtung beim Schnitt vom 2026-10-09* (Stand `cfc4d6b`):

- `BEO-REPO/slice-waechst-durch-uebernahmen` (2×) — dieser Slice wuchs durch die Übernahmen
  oben über eine Review-Sitzung, bei drei Liefer-Punkten; das ist das Muster des Eintrags.
  Trägt seine Closure den Beleg ein, steht der Eintrag bei 3× und braucht einen eigenen
  Folge-Slice; die Entscheidung gehört in die Closure.
- `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel` (1×) — diesmal läuft die Rückführung über
  `next/` (Entscheidung des Nutzers vom 2026-10-09 zu V-117); kein weiterer Beleg.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu; keine neue Lücke vor
dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
