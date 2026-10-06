# Slice slice-v1-abschluss-betrieb: Signale, Schreiben und Konfiguration

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-08`](../../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-15`](../../../../spec/lastenheft.md#lh-fa-15--ci-eignung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)

**Berührte Spec-Stellen:** `LH-FA-01.a` · `LH-FA-07.a` · `LH-FA-13.a` · `LH-FA-13.b` · `LH-FA-17.a` · `LH-FA-02.b` · `LH-FA-03.b` · `SPEC-007` · `SPEC-008` · `SPEC-012` · `SPEC-013` bis `SPEC-019` · `SPEC-034` · `SPEC-046`

**Verantwortlich:** —
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Prozess fährt auf `SIGINT`/`SIGTERM` kontrolliert herunter, das Recording wird atomar geschrieben, `--output`/`--force` verhalten sich wie spezifiziert, und alle Optionen sind per CLI und Umgebungsvariable setzbar.

**Übernommen aus `slice-extended-query-record`** (Risiko F-309/F-317, Validierung `docs/reviews/2026-10-05-validierung-slice-extended-query-record.md`, Frage 1): Das Herunterfahren bekommt eine Obergrenze, in `record` und in `replay` (Optionstabelle, `LH-FA-13.a`). Heute wartet `record` ohne Frist auf die laufende Interaktion jeder Session; im Container beendet der `SIGKILL` nach der Stopp-Frist den Prozess, und die Aufzeichnung jeder noch wartenden Session fehlt ganz. **Übernommen aus `slice-extended-query-replay`** (Folge-Review F-330): `replay` beantwortet beim Herunterfahren eine begonnene Extended-Interaktion bis zu ihrem `Sync` und wartet dabei ebenso ohne Frist auf einen pausierenden Client. `replay` setzt die Frist durch Schließen der Client-Verbindung durch (`replaySitzung`); dieses Ende behandelt `replaySitzung` heute als regulär, die Frist muss es für eine unvollständige Interaktion als `PGR-E4006` merken.

- Frist, gezählt ab dem Signal, Default 5 s, setzbar per Option und Umgebungsvariable ([`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)); `0` heißt ohne Frist.
- Nach Ablauf endet jede noch laufende Session zwangsweise wie bei einem Abbruch nach `LH-FA-02.b`: die laufende Interaktion wird nicht übernommen, die abgeschlossenen bleiben; danach wird die Aufzeichnung geschrieben. In `replay` endet die Session durch Schließen der Client-Verbindung, und eine unvollständige Interaktion ist `PGR-E4006`.
- Meldungscode `PGR-E4006` der Klasse 4 (vergeben in `SPEC-034`, nicht `PGR-E4003`), Exit-Code `4`; das Log nennt je abgebrochener Session die Session und die verworfene Interaktion.
- Ein zweites Signal lässt die Frist sofort ablaufen, statt den Prozess hart zu beenden.
- Beim Beginn des Herunterfahrens eine Info-Zeile mit der Zahl der Sessions, auf die gewartet wird.

Die Ergänzung von [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) Boundary ist bestätigt und steht im Lastenheft; `LH-FA-13.a` (Frist, Zwangsende, `PGR-E4006`, weiteres Signal, Info-Zeile), die Optionstabelle (`--shutdown-timeout`) und `SPEC-046` sind spezifiziert. Ob der Gegenstand in DoD-Punkt 1 passt oder einen eigenen Slice der `welle-v1-abschluss` braucht, entscheidet der Planner beim `open` → `next` dieses Slice (zweites Folge-Review zu `slice-extended-query-replay`, F-345: DoD-Punkt 1 bündelt atomares Schreiben, Exit-Code aller Klassen und die Frist in zwei Modi; §3 führt vier Komponenten). Zeigt sich erst in der Arbeit, dass DoD-Punkt 1 nicht in eine Review-Sitzung passt, gilt die Rückführung in §4.

**Übernommen aus `slice-replay-semantik-mismatch`:** der Schlüssel `fail_on_unconsumed` im Abschnitt `replay:` (mit der Konfigurationsdatei für alle Optionen); die Werte boolescher Optionen nach `LH-FA-17.a` auch für `--force` (heute nimmt es `1` und `t` an); die Prüfung jeder gesetzten Umgebungsvariable einer Option des Kommandos, auch wenn die Kommandozeile vorgeht (`PGR-E2001`, `LH-FA-17.a`), im allgemeinen Leser für alle Optionen; die Hilfe vor jeder Prüfung von Optionen, Umgebungsvariablen und Konfigurationsdatei (`LH-FA-01.a`), auch für `config show` und `--config`, samt `--` als Ende der Optionen; und der Test, dass eine durch die Frist zwangsweise beendete Replay-Session mit `--fail-on-unconsumed` `PGR-E4006` vor `PGR-E5002` merkt und mit Exit-Code `4` endet (`LH-FA-03.b`).

**Übernommen aus `slice-replay-semantik-meldungscodes`:** der Schlüssel `log_level` auf der obersten Ebene der Konfigurationsdatei (`LH-FA-17.a`), mit derselben Wertemenge wie `--log-level`; ein ungültiger Wert ist `PGR-E2004`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Container-Image — `slice-v1-abschluss-container`.
- Meldungscodes der neuen Fehler — gehören mit zu `PGR-E2002` und folgen `SPEC-034`; die Tabelle selbst liefert `slice-replay-semantik-meldungscodes`.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Herunterfahren schreibt das Recording atomar und liefert den Exit-Code der gemerkten Klasse, für alle Klassen aus `SPEC-013` bis `SPEC-019`; die Frist `--shutdown-timeout` begrenzt das Warten in `record` und `replay`, und eine dabei unvollständige Interaktion ist `PGR-E4006`; eine so beendete Replay-Session mit `--fail-on-unconsumed` merkt `PGR-E4006` vor `PGR-E5002` und endet mit Exit-Code 4 (Test mit Signal, je Modus).
- [ ] [`LH-FA-08`](../../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Vorhandenes `--output` wird ohne `--force` abgelehnt; `--force` und jede boolesche Option nehmen nur `true` oder `false`; Priorität CLI vor Umgebungsvariable vor Konfigurationsdatei vor Default, und jede gesetzte Umgebungsvariable einer Option des Kommandos wird geprüft, auch wenn die Kommandozeile vorgeht (`PGR-E2001`); die Hilfe geht jeder Prüfung von Optionen, Umgebungsvariablen und Konfigurationsdatei vor, auch für `config show` und `--config`; `config show` zeigt die gewählte Datei, ohne einen aufgelösten Wert (Test).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Benannte Verbindungen (`connections`, `--upstream <Name>`, `sslmode`), `--config`, `PGWIRE_RECORDER_CONFIG`, die Standarddatei, der Schlüssel `fail_on_unconsumed` im Abschnitt `replay:`, der Schlüssel `log_level` auf der obersten Ebene (Wertemenge und Strenge von `--log-level`, ungültig `PGR-E2004`) und `$${VAR}` verhalten sich wie spezifiziert; eine ungültige Datei, eine nicht gesetzte Variable und ein Klartext-Passwort sind `PGR-E2004` bis `PGR-E2006` (Test).
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
| `internal/adapters/driving/cli` | update | Optionstabelle, Priorität, Signale; allgemeiner Leser für Umgebungsvariablen, der jeden gesetzten Wert prüft und die zwei Einzel-Leser von `--fail-on-unconsumed` und `--log-level` ersetzt; strenge boolesche Werte auch für `--force`; Hilfe vor jeder Prüfung; Schlüssel `fail_on_unconsumed` (übernommen aus `slice-replay-semantik-mismatch`); Schlüssel `log_level` auf der obersten Ebene, ungültig `PGR-E2004` (übernommen aus `slice-replay-semantik-meldungscodes`) |
| `internal/adapters/driving/pgwire`, `internal/bootstrap` | update | Frist in `recordSitzung` und `replaySitzung`; `replay` schließt bei Ablauf die Client-Verbindung und merkt `PGR-E4006` |
| `internal/adapters/driven/recording` | update | temporäre Datei, atomares Verschieben |
| `test/integration` | update | Happy/Boundary/Negative |
| `docs/user/benutzerhandbuch.md` | update | Herunterfahren mit Frist in `record` und `replay`; der Abschnitt nennt auch `replay` im Container (heute `SIGKILL` nach der Stopp-Frist, Exit-Code `137`, keine Log-Zeile; Validierung zu `slice-extended-query-replay`, Frage 5) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-replay-semantik` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Atomares Schreiben verlangt plattformspezifische Lösungen — zurück zur Zerlegung.
- `in-progress` → `next`: DoD-Punkt 1 passt nicht in eine Review-Sitzung (F-345) — die Frist `--shutdown-timeout` in beiden Modi geht als eigener Slice der `welle-v1-abschluss` heraus.
- `in-progress` → `open`: Signalbehandlung ist im Container nicht prüfbar — Carveout.


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

- Atomarität des Verschiebens ist plattformabhängig ("bestmöglich atomar") — **Ausgang:** offen bis Closure.

- Herunterfahren ohne Obergrenze (aus `slice-extended-query-record`, Review F-309, Folge-Review F-317, Validierung Frage 1): Eine Session mit laufender Interaktion — einfache Anfrage oder Extended-Interaktion ohne `Sync` — hält das Herunterfahren beliebig lange; der Container-Stopp verliert dann ihre Aufzeichnung ganz, ohne Log-Zeile und mit Exit-Code `137`. Gegenstand siehe §1; Lastenheft und Spezifikation sind ergänzt — **Ausgang:** offen bis Closure.

- Herunterfahren ohne Obergrenze im Replay (aus `slice-extended-query-replay`, Folge-Review F-330): Ein Client, der mitten in einer Extended-Interaktion pausiert, hält `replay` nach `SIGTERM` beliebig lange; endet die Session durch Schließen der Verbindung, behandelt `replaySitzung` das heute als reguläres Ende (Exit-Code 0 statt 4). Im Container endet ein solches Warten nach der Stopp-Frist mit `SIGKILL`, Exit-Code `137` und ohne Log-Zeile nach dem Start; verloren geht nichts, weil `replay` nichts schreibt (Validierung zu `slice-extended-query-replay`, Frage 5, Sonde S10) — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
