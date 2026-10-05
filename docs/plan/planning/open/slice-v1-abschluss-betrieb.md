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

**Bezug:** [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-08`](../../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-15`](../../../../spec/lastenheft.md#lh-fa-15--ci-eignung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)

**Berührte Spec-Stellen:** `LH-FA-07.a` · `LH-FA-13.a` · `LH-FA-13.b` · `LH-FA-17.a` · `SPEC-007` · `SPEC-008` · `SPEC-012`

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

**Übernommen aus `slice-extended-query-record`** (Risiko F-309/F-317, Validierung `docs/reviews/2026-10-05-validierung-slice-extended-query-record.md`, Frage 1): Das Herunterfahren bekommt eine Obergrenze. Heute wartet `record` ohne Frist auf die laufende Interaktion jeder Session; im Container beendet der `SIGKILL` nach der Stopp-Frist den Prozess, und die Aufzeichnung jeder noch wartenden Session fehlt ganz.

- Frist, gezählt ab dem Signal, Default 5 s, setzbar per Option und Umgebungsvariable ([`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)); `0` heißt ohne Frist.
- Nach Ablauf endet jede noch laufende Session zwangsweise wie bei einem Abbruch nach `LH-FA-02.b`: die laufende Interaktion wird nicht übernommen, die abgeschlossenen bleiben; danach wird die Aufzeichnung geschrieben.
- Eigener Meldungscode der Klasse 4 (Vergabe nach `SPEC-034`, nicht `PGR-E4003`), Exit-Code `4`; das Log nennt je abgebrochener Session die Session und die verworfene Interaktion.
- Ein zweites Signal lässt die Frist sofort ablaufen, statt den Prozess hart zu beenden.
- Beim Beginn des Herunterfahrens eine Info-Zeile mit der Zahl der Sessions, auf die gewartet wird.

Die Ergänzung von [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) Boundary (begrenzte, einstellbare Zeit für das kontrollierte Beenden) wartet auf die Bestätigung des Nutzers; erst danach folgen `LH-FA-13.a`, `LH-FA-13.b` und die Optionstabelle der Spezifikation. Ob der Gegenstand in DoD-Punkt 1 passt oder einen eigenen Slice der `welle-v1-abschluss` braucht, entscheidet der Schnitt nach dieser Bestätigung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Container-Image — `slice-v1-abschluss-container`.
- Meldungscodes der neuen Fehler — gehören mit zu `PGR-E2002` und folgen `SPEC-034`; die Tabelle selbst liefert `slice-replay-semantik-meldungscodes`.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Herunterfahren schreibt das Recording atomar und liefert den Exit-Code der gemerkten Klasse, für alle Klassen aus `SPEC-013` bis `SPEC-019` (Test mit Signal).
- [ ] [`LH-FA-08`](../../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Vorhandenes `--output` wird ohne `--force` abgelehnt; Priorität CLI vor Umgebungsvariable vor Konfigurationsdatei vor Default; `config show` zeigt die gewählte Datei, ohne einen aufgelösten Wert (Test).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Benannte Verbindungen (`connections`, `--upstream <Name>`, `sslmode`), `--config`, `PGWIRE_RECORDER_CONFIG`, die Standarddatei und `$${VAR}` verhalten sich wie spezifiziert; eine ungültige Datei, eine nicht gesetzte Variable und ein Klartext-Passwort sind `PGR-E2004` bis `PGR-E2006` (Test).
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
| `internal/adapters/driving/cli` | update | Optionstabelle, Priorität, Signale |
| `internal/adapters/driven/recording` | update | temporäre Datei, atomares Verschieben |
| `test/integration` | update | Happy/Boundary/Negative |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-replay-semantik` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Atomares Schreiben verlangt plattformspezifische Lösungen — zurück zur Zerlegung.
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

- Herunterfahren ohne Obergrenze (aus `slice-extended-query-record`, Review F-309, Folge-Review F-317, Validierung Frage 1): Eine Session mit laufender Interaktion — einfache Anfrage oder Extended-Interaktion ohne `Sync` — hält das Herunterfahren beliebig lange; der Container-Stopp verliert dann ihre Aufzeichnung ganz, ohne Log-Zeile und mit Exit-Code `137`. Gegenstand siehe §1; die Lastenheft-Änderung wartet auf die Bestätigung des Nutzers — **Ausgang:** offen bis Closure.

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
