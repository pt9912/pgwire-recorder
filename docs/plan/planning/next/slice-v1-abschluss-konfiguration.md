# Slice slice-v1-abschluss-konfiguration: Konfiguration über Kommandozeile, Umgebung und Datei

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0014](../../adr/0014-konfigurationsdatei.md)

**Berührte Spec-Stellen:** `LH-FA-01.a` · `LH-FA-17.a` · `LH-FA-03.b` · `SPEC-007` · `SPEC-008` · `SPEC-012` · `SPEC-014` · `SPEC-034` · `SPEC-046`

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

**Ziel:** Jede Option von `record` und `replay` ist per Kommandozeile, Umgebungsvariable und Konfigurationsdatei setzbar, in der Priorität Kommandozeile vor Umgebungsvariable vor Datei vor Default; die Hilfe geht jeder Prüfung vor, und ein Fehler in jeder Quelle ist `PGR-E2001` oder `PGR-E2004` bis `PGR-E2006`.

**Übernimmt:** `slice-v1-abschluss-betrieb` — dessen Teil *Konfiguration* (DoD-Punkt 2 ohne
`--output` und `--force`, DoD-Punkt 3; Entscheidung des Nutzers vom 2026-10-08, Schnitt nach
F-345; dort §7 `Gegenstand:`). Im Einzelnen, je mit dem ursprünglichen Geber:

- **Aus `slice-replay-semantik-mismatch`:** der Schlüssel `fail_on_unconsumed` im Abschnitt
  `replay:`; die Prüfung jeder gesetzten Umgebungsvariable einer Option des Kommandos, auch
  wenn die Kommandozeile vorgeht (`PGR-E2001`, `LH-FA-17.a`), im allgemeinen Leser für alle
  Optionen; die Hilfe vor jeder Prüfung von Optionen, Umgebungsvariablen und
  Konfigurationsdatei (`LH-FA-01.a`), auch für `config show` und `--config`, samt `--` als
  Ende der Optionen.
- **Aus `slice-replay-semantik-meldungscodes`:** der Schlüssel `log_level` auf der obersten
  Ebene der Konfigurationsdatei (`LH-FA-17.a`), mit derselben Wertemenge wie `--log-level`;
  ein ungültiger Wert ist `PGR-E2004`.
- **Aus `slice-v1-abschluss-herunterfahren`** (dort §1, Abgrenzung): der Schlüssel der Frist
  `--shutdown-timeout` in der Konfigurationsdatei, und dass der allgemeine Leser auch
  `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT` liest.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `--output`, `--force` und die strengen Werte für `--force` — `slice-v1-abschluss-schreiben`;
  er folgt diesem Slice und liest `--force` über den allgemeinen Leser von hier.
- Die Optionen von `play` — `slice-v1-abschluss-einspielen` (dort §1: `--fail-on-unconsumed`
  und `--log-level` bei `play`); er setzt auf dem allgemeinen Leser dieses Slice auf.
- Signale und die Frist selbst — `slice-v1-abschluss-herunterfahren`; hier nur ihre
  Quellen.
- Code im Kern, im PGWire-Adapter und in den Driven-Adaptern — Schicht-Abgrenzung: Der
  Slice ändert den CLI-Adapter und den Bootstrap.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Für jede Option von `record` und `replay`, auch `--shutdown-timeout`, gilt
      die Priorität Kommandozeile vor Umgebungsvariable vor Konfigurationsdatei vor Default
      (`SPEC-007`); ein allgemeiner Leser prüft jede gesetzte Umgebungsvariable einer Option
      des Kommandos, auch wenn die Kommandozeile vorgeht (`PGR-E2001`), und ersetzt die
      Einzel-Leser von `--fail-on-unconsumed`, `--log-level` und `--shutdown-timeout` (Test).
- [ ] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Die Hilfe geht jeder Prüfung von Optionen,
      Umgebungsvariablen und Konfigurationsdatei vor, auch für `config show` und
      `--config`, samt `--` als Ende der Optionen; `config show` zeigt die gewählte Datei,
      ohne einen aufgelösten Wert (Test).
- [ ] [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Benannte Verbindungen (`connections`, `--upstream <Name>`, `sslmode`),
      `--config`, `PGWIRE_RECORDER_CONFIG`, die Standarddatei, der Schlüssel
      `fail_on_unconsumed` im Abschnitt `replay:`, der Schlüssel `log_level` auf der obersten
      Ebene (Wertemenge und Strenge von `--log-level`, ungültig `PGR-E2004`) und `$${VAR}`
      verhalten sich wie spezifiziert; eine ungültige Datei, eine nicht gesetzte Variable und
      ein Klartext-Passwort sind `PGR-E2004` bis `PGR-E2006` (Test). Beleg in §7 für alle
      drei Punkte: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
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
| `internal/adapters/driving/cli` | update | allgemeiner Leser für Umgebungsvariablen, der jeden gesetzten Wert prüft und die Einzel-Leser von `--fail-on-unconsumed`, `--log-level` und `--shutdown-timeout` ersetzt; Priorität; Hilfe vor jeder Prüfung; `config show`; Konfigurationsdatei mit `connections`, `fail_on_unconsumed`, `log_level`, `$${VAR}` |
| `internal/bootstrap` | update | benannte Verbindung (`--upstream <Name>`, `sslmode`) an den Upstream übergeben |
| `internal/adapters/driving/cli` (Unit-Tests), `test/integration` | update | Happy/Boundary/Negative nach `LH-FA-17.a` und `LH-FA-01.a` |
| `docs/user/benutzerhandbuch.md` | update | Konfiguration über drei Quellen und ihre Priorität |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-herunterfahren` liegt in `done/` (der allgemeine Leser
übernimmt dessen Umgebungsvariable). Schritt 3 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)
(Entscheidung des Nutzers vom 2026-10-08). Vor dem ersten Code-Commit prüft der Architect
die Randformen aus §6 und entscheidet die offenen (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Benannte Verbindungen mit
  `sslmode` verlangen eine Änderung am Upstream-Adapter, also eine dritte Schicht; dann
  gehen die benannten Verbindungen als eigener Slice der Welle heraus.
- `in-progress` → `open` (blockiert — Carveout?): Die Form der Datei nach [ADR-0014](../../adr/0014-konfigurationsdatei.md) trägt
  einen Schlüssel nicht, den eine Option braucht, ohne eine neue Entscheidung; dann zuerst
  die Entscheidung.

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

- **Gesetzte, aber leere Umgebungsvariable** — offen: Der Architect bestätigt vor dem Code
  in `LH-FA-17.a`, ob sie als gesetzt gilt.
- **Unbekannter Schlüssel in der Datei** — `PGR-E2004` als ungültige Datei; der Architect
  bestätigt vor dem Code, dass `LH-FA-17.a` das deckt.
- **`--config` nennt eine fehlende Datei, die Standarddatei fehlt** — der erste Fall ist
  ein Fehler, der zweite nicht; der Architect bestätigt vor dem Code die Stelle in
  `LH-FA-17.a`.
- **`$${VAR}` nicht gesetzt, Maskierung von `$`** — entschieden in `LH-FA-17.a`
  (`PGR-E2005` bzw. die Maskierung dort).
- **Klartext-Passwort in der Datei** — `PGR-E2006`; entschieden in `LH-FA-17.a`.
- **Boolesche Werte in der Datei** (YAML `true` gegen Zeichenkette `"true"`) — offen: Der
  Architect entscheidet vor dem Code in der Spezifikation.
- **Hilfe nach einer ungültigen Option** — die Hilfe geht vor; entschieden in `LH-FA-01.a`.
- **`log_level` ungültig** — `PGR-E2004`; entschieden in `LH-FA-17.a`.

**Risiken:**

- Der allgemeine Leser ersetzt drei Einzel-Leser; eine Prüfung eines Einzel-Lesers kann
  dabei wegfallen, statt in den allgemeinen überzugehen
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×) — **Ausgang:** offen bis Closure.

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
