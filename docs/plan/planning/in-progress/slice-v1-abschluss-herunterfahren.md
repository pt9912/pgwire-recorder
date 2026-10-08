# Slice slice-v1-abschluss-herunterfahren: Herunterfahren mit Frist in record und replay

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-15`](../../../../spec/lastenheft.md#lh-fa-15--ci-eignung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0003](../../adr/0003-pgwire-server-ist-driving-adapter.md)

**Berührte Spec-Stellen:** `LH-FA-13.a` · `LH-FA-13.b` · `LH-FA-02.b` · `LH-FA-03.b` · `LH-FA-17.a` · `SPEC-013` bis `SPEC-019` · `SPEC-034` · `SPEC-046`

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

**Ziel:** `record` und `replay` fahren auf `SIGINT`/`SIGTERM` kontrolliert herunter, warten dabei höchstens `--shutdown-timeout` lang auf laufende Sessions und enden mit dem Exit-Code der gemerkten Klasse; eine durch die Frist abgebrochene Interaktion ist `PGR-E4006`.

**Übernimmt:** `slice-v1-abschluss-betrieb` — dessen Teil *Herunterfahren* (Entscheidung des
Nutzers vom 2026-10-08, Schnitt nach F-345 in drei Slices; dort §7 `Gegenstand:`). Im
Einzelnen, je mit dem ursprünglichen Geber:

- **Aus `slice-extended-query-record`** (Risiko F-309/F-317, Validierung
  `docs/reviews/2026-10-05-validierung-slice-extended-query-record.md`, Frage 1): Das
  Herunterfahren bekommt eine Obergrenze. Heute wartet `record` ohne Frist auf die
  laufende Interaktion jeder Session; im Container beendet der `SIGKILL` nach der
  Stopp-Frist den Prozess, und die Aufzeichnung jeder noch wartenden Session fehlt ganz.
  Frist gezählt ab dem ersten Signal, Default 5 s, `0` heißt ohne Frist, setzbar per
  Option und Umgebungsvariable ([`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)). Nach Ablauf endet jede noch laufende Session
  zwangsweise wie bei einem Abbruch nach `LH-FA-02.b`: die laufende Interaktion wird nicht
  übernommen, die abgeschlossenen bleiben; danach wird die Aufzeichnung geschrieben.
  `PGR-E4006` der Klasse 4 (vergeben in `SPEC-034`), Exit-Code `4`; das Log nennt je
  abgebrochener Session die Session und die verworfene Interaktion. Ein zweites Signal
  lässt die Frist sofort ablaufen. Beim Beginn eine Info-Zeile mit der Zahl der Sessions,
  auf die gewartet wird.
- **Aus `slice-extended-query-replay`** (Folge-Review F-330): `replay` beantwortet beim
  Herunterfahren eine begonnene Extended-Interaktion bis zu ihrem `Sync` und wartet dabei
  ohne Frist auf einen pausierenden Client. `replay` setzt die Frist durch Schließen der
  Client-Verbindung durch (`replaySitzung`); dieses Ende behandelt `replaySitzung` heute
  als regulär, die Frist merkt es für eine unvollständige Interaktion als `PGR-E4006`.
- **Aus `slice-replay-semantik-mismatch`:** der Test, dass eine durch die Frist
  zwangsweise beendete Replay-Session mit `--fail-on-unconsumed` `PGR-E4006` vor
  `PGR-E5002` merkt und mit Exit-Code `4` endet (`LH-FA-03.b`).
- **Aus `slice-v1-abschluss-sessions`** (dort §1, Abgrenzung): die Signalbehandlung.
- **Aus `slice-harness-integration-wait`** (dort §1, Abgrenzung): die Tests mit Signal und
  Frist je Modus, geschrieben über die Helfer, die jener umbaut; sie sind die Tests in
  DoD-Punkt 1 und 2. Dazu aus dessen Verifikation (V-105) die Randform *`stderr` zur
  Laufzeit lesen* in §6.
- **Aus `slice-harness-blackbox-pgwire`** (Review F-451, Verifikation V-87; Abgrenzung in
  `slice-tests-ueberlebende-mutanten`): der Wert der Schreibfrist der Fehlerantwort beim
  Session-Ende (`meldeFrist`) und ein Test mit der Schranke als Literal; Randform in §6.
  Er gehört zu DoD-Punkt 1, weil die Frist einen Teil des Budgets von `--shutdown-timeout`
  verbraucht.
- **Aus `slice-tests-ueberlebende-mutanten-driving`** (dort §1, Abgrenzung; Verifikation
  V-95 zu `slice-lint-bestand-driving`): das Ende einer Replay-Sitzung durch die Frist und
  `PGR-E4006`. Der Test dazu entsteht hier; jener Slice schreibt später den Test zu G2 (ein
  Ende, das der Client auslöst) daneben, ohne diesen zu ändern.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Atomares Schreiben der Aufzeichnung, `--output` und `--force` — `slice-v1-abschluss-schreiben`;
  dieser Slice schreibt nach dem Zwangsende so, wie der Recording-Adapter heute schreibt.
- Der Schlüssel der Frist in der Konfigurationsdatei und der allgemeine Leser der
  Umgebungsvariablen — `slice-v1-abschluss-konfiguration`; hier wird
  `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT` gelesen wie die Umgebungsvariablen der übrigen
  Optionen heute, und jener Slice ersetzt den Leser.
- Ein Ende der Replay-Sitzung, das der Client auslöst (G2) —
  `slice-tests-ueberlebende-mutanten-driving`.
- Container-Image — `slice-v1-abschluss-container`.
- Neue Meldungscodes über `PGR-E4006` hinaus — Bestand bleibt stehen: `SPEC-034` vergibt
  `PGR-E4006`, die Tabelle liegt seit `slice-replay-semantik-meldungscodes`.
- Code im Kern (`internal/hexagon/`) und in den Driven-Adaptern — Schicht-Abgrenzung: Der
  Slice ändert die Driving-Adapter (`pgwire`, `cli`) und den Bootstrap.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-15`](../../../../spec/lastenheft.md#lh-fa-15--ci-eignung): Nach `SIGINT` oder `SIGTERM` wartet `record`
      höchstens `--shutdown-timeout` (ab dem ersten Signal, `0` ohne Frist) auf laufende
      Sessions; danach endet jede noch laufende Session wie bei einem Abbruch nach
      `LH-FA-02.b`, die Aufzeichnung wird geschrieben, das Log nennt je Session die
      verworfene Interaktion, der Prozess merkt `PGR-E4006` und endet mit Exit-Code 4; ein
      zweites Signal lässt die Frist sofort ablaufen; beim Beginn steht eine Info-Zeile mit
      der Zahl der Sessions; `meldeFrist` hat den Wert der Spezifikation (Test mit Signal
      über die Helfer aus `slice-harness-integration-wait`; Test mit der Schranke als
      Literal).
- [ ] [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Die Frist begrenzt das Warten in `replay`; bei Ablauf
      schließt `replay` die Client-Verbindung, eine unvollständige Interaktion ist
      `PGR-E4006`, Exit-Code 4; mit `--fail-on-unconsumed` merkt die Session `PGR-E4006`
      vor `PGR-E5002` und der Prozess endet mit Exit-Code 4 (Test mit Signal).
- [ ] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Das kontrollierte Herunterfahren liefert den Exit-Code der gemerkten
      Klasse, für alle Klassen aus `SPEC-013` bis `SPEC-019` (0 ohne zuvor aufgetretenen
      Verbindungsfehler); das Benutzerhandbuch beschreibt das Herunterfahren mit Frist in
      `record` und `replay`, auch `replay` im Container (Test je Klasse). Beleg in §7 für
      alle drei Punkte: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
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
| `internal/bootstrap` | update | Signale, Frist ab dem ersten Signal, zweites Signal, Info-Zeile, Exit-Code der gemerkten Klasse |
| `internal/adapters/driving/pgwire` | update | Frist in `recordSitzung` und `replaySitzung`; `replay` schließt bei Ablauf die Client-Verbindung und merkt `PGR-E4006`; `meldeFrist` nach der Spezifikation |
| `internal/adapters/driving/cli` | update | Option `--shutdown-timeout` und `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT` für `record` und `replay` (Optionstabelle, `SPEC-046`) |
| `internal/adapters/driving/pgwire` (Unit-Tests) | update | `meldeFrist` mit der Schranke als Literal, nach der Entscheidung in §6 |
| `test/integration` | update | Signal und Frist je Modus, Happy/Boundary/Negative nach `LH-FA-13.a`, über die Helfer aus `slice-harness-integration-wait` |
| `docs/user/benutzerhandbuch.md` | update | Herunterfahren mit Frist in `record` und `replay`; `replay` im Container (heute `SIGKILL` nach der Stopp-Frist, Exit-Code `137`, keine Log-Zeile; Validierung zu `slice-extended-query-replay`, Frage 5) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-integration-wait` liegt in `done/` (die Tests mit Signal und
Frist entstehen über dessen Helfer). Schritt 2 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)
(Entscheidung des Nutzers vom 2026-10-08). Vor dem ersten Code-Commit prüft der Architect
die Randformen aus §6 und entscheidet die offenen (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Frist in `replay` verlangt
  einen Umbau im Kern (`internal/hexagon/`), also eine dritte Schicht; dann geht die
  Frist in `replay` (DoD-Punkt 2) als eigener Slice der Welle heraus.
- `in-progress` → `open` (blockiert — Carveout?): Signalbehandlung ist im Container nicht
  prüfbar — Carveout.

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

**Randformen** (`AGENTS.md` §3.12) — je Randform, wo sie entschieden ist:

- **Frist `0`** — ohne Frist; entschieden in `LH-FA-13.a` und der Optionstabelle (`SPEC-046`).
- **Form des Werts** (Einheit `ms`, `s` oder `m`; ohne Einheit, negativ, leer) — die Form
  steht in der Optionstabelle; welcher Code einen ungültigen Wert meldet, bestätigt der
  Architect vor dem Code in der Spezifikation.
- **Zweites Signal** — die Frist läuft sofort ab; entschieden in `LH-FA-13.a`.
- **Signal ohne laufende Session** (auch vor der ersten Verbindung) — offen: Der
  Architect bestätigt vor dem Code, dass der Prozess ohne Warten endet und welche
  Info-Zeile dann steht.
- **Ablauf der Frist während des Schreibens der Aufzeichnung** — offen: Der Architect
  entscheidet vor dem Code, ob das Schreiben noch zur Frist zählt.
- **Ende durch die Frist mitten in einer Extended-Interaktion im Replay** — `PGR-E4006`;
  entschieden in `LH-FA-13.a`; mit `--fail-on-unconsumed` `PGR-E4006` vor `PGR-E5002`,
  entschieden in `LH-FA-03.b`.
- **Wert der Schreibfrist der Fehlerantwort beim Session-Ende** (seit slice-harness-blackbox-pgwire,
  Review F-451, Verifikation V-87; übernommen aus `slice-v1-abschluss-betrieb`) — der
  Recorder schreibt die Fehlerantwort beim Ende einer Session höchstens `meldeFrist` lang an
  einen nicht lesenden Client (`internal/adapters/driving/pgwire/server.go`, im Code eine
  Sekunde); den Wert nennen weder Spezifikation noch Lastenheft, der Absatz *Abbruch* im
  Record-Teil sagt nur, dass das Schließen der Client-Verbindung ein blockiertes Schreiben
  beendet. Die Frist verbraucht einen Teil des Budgets von `--shutdown-timeout`. Offen: Der
  Architect entscheidet den Wert vor dem Code in der Spezifikation; danach folgt ein Test
  mit der Schranke als Literal, rot bei `meldeFrist = 3 * time.Second`. Heute liest
  `TestFehlerantwortMitFrist` die Konstante über die Brücke und fängt nur das Entfernen der
  Frist.
- **`stderr` zur Laufzeit lesen** (seit slice-harness-integration-wait, Verifikation V-105)
  — offen: Der Helfer aus `slice-harness-integration-wait` liest `stderr` ohne Wettlauf erst
  nach dem Ende des Prozesses (dort §6 *`stderr` erst nach dem Ende lesen*), nicht solange
  der Prozess läuft. Wartet ein Test vor dem zweiten Signal auf die Info-Zeile mit der Zahl
  der Sessions oder auf die Log-Zeile einer abgebrochenen Session, entscheidet der Architect
  vor dem Code, wie: die Log-Zeilen erst nach dem Ende prüfen, oder ein synchronisierter
  Puffer im Geschirr.

**Risiken:**

- Herunterfahren ohne Obergrenze (aus `slice-extended-query-record`, Review F-309,
  Folge-Review F-317, Validierung Frage 1; übernommen aus `slice-v1-abschluss-betrieb`):
  Eine Session mit laufender Interaktion — einfache Anfrage oder Extended-Interaktion ohne
  `Sync` — hält das Herunterfahren beliebig lange; der Container-Stopp verliert dann ihre
  Aufzeichnung ganz, ohne Log-Zeile und mit Exit-Code `137` — **Ausgang:** offen bis
  Closure.
- Herunterfahren ohne Obergrenze im Replay (aus `slice-extended-query-replay`, Folge-Review
  F-330; übernommen aus `slice-v1-abschluss-betrieb`): Ein Client, der mitten in einer
  Extended-Interaktion pausiert, hält `replay` nach `SIGTERM` beliebig lange; endet die
  Session durch Schließen der Verbindung, behandelt `replaySitzung` das heute als
  reguläres Ende (Exit-Code 0 statt 4) — **Ausgang:** offen bis Closure.
- Ein Test mit Signal, dessen Prozess unter einer Mutation nicht endet, hängt bis zum
  Zeitlimit von `go test` (`BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`, 2×); jeder
  neue Test wartet darum über die Helfer aus `slice-harness-integration-wait` mit eigener
  Frist — **Ausgang:** offen bis Closure.

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

- **Belege zur DoD (Implementer):** <je Zusage: Zusage · Mutation · roter Test (`AGENTS.md` §3.10)>
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

- `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` (2×) — die Tests mit Signal; ein Risiko in
  §6. Hinge ein Test dieses Slice bis zum Zeitlimit, erreichte der Eintrag 3×.
- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (1×) — die Frist im Replay hängt an
  einem pausierenden Client; die Randformen in §6 nennen die Fälle.
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
