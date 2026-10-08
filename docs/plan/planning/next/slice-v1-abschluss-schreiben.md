# Slice slice-v1-abschluss-schreiben: Atomares Schreiben, --output und --force

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-08`](../../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0005](../../adr/0005-recording-store-ist-driven-adapter.md)

**Berührte Spec-Stellen:** `LH-FA-07.a` · `LH-FA-17.a` · `SPEC-014` · `SPEC-016` · `SPEC-034`

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

**Ziel:** `record` schreibt die Aufzeichnung atomar (temporäre Datei, dann Verschieben), lehnt ein vorhandenes `--output` ohne `--force` ab, und `--force` nimmt nur `true` oder `false`.

**Übernimmt:** `slice-v1-abschluss-betrieb` — dessen Teil *Schreiben* (atomares Schreiben aus
DoD-Punkt 1, `--output` und `--force` aus DoD-Punkt 2; Entscheidung des Nutzers vom
2026-10-08, Schnitt nach F-345; dort §7 `Gegenstand:`). Im Einzelnen, je mit dem
ursprünglichen Geber:

- **Aus `slice-replay-semantik-mismatch`:** die Werte boolescher Optionen nach
  `LH-FA-17.a` auch für `--force` (heute nimmt es `1` und `t` an).
- **Aus `slice-v1-abschluss-herunterfahren`** (dort §1, Abgrenzung): atomares Schreiben,
  `--output` und `--force`; das atomare Schreiben gilt auch für die Aufzeichnung nach
  dessen Zwangsende (DoD-Punkt 1).
- **Aus `slice-v1-abschluss-konfiguration`** (dort §1, Abgrenzung): `--output`, `--force`
  und die strengen Werte für `--force`, gelesen über dessen allgemeinen Leser.
- **Aus `slice-v1-abschluss-konfigurationsdatei`** (dort §1, Abgrenzung): die Schlüssel
  `output` und `force` der Konfigurationsdatei; sie folgen aus der Anmeldung von `--output`
  und `--force` am allgemeinen Leser und sind bis dahin unbekannt (`PGR-E2004`). Das ist der
  Stand des Plans zu Rückgabe 12 in §6 jenes Slice; entscheidet der Architect anders, zieht er
  diese Zeile im selben Commit nach.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Herunterfahren und die Frist — `slice-v1-abschluss-herunterfahren`; er liegt vor diesem
  Slice, und das atomare Schreiben gilt auch für die Aufzeichnung nach seinem Zwangsende.
- Der allgemeine Leser — `slice-v1-abschluss-konfiguration`; dieser Slice meldet `--output`
  und `--force` an dessen Leser an.
- Die Konfigurationsdatei — `slice-v1-abschluss-konfigurationsdatei`; dieser Slice liest die
  Schlüssel `output` und `force` über deren Laden, ohne eigenen Code für die Datei.
- Das Schreiben im SQLite-Format — anderer Vorgang: `slice-v1-abschluss-sqlite-format`
  schreibt seine Datei je Session in einer Transaktion; dieser Slice schreibt die Datei
  des Standardformats YAML.
- Code im Kern, im PGWire-Adapter und im Bootstrap — Schicht-Abgrenzung: Der Slice ändert
  den Recording-Adapter und den CLI-Adapter.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings): `record` schreibt die Aufzeichnung in eine temporäre Datei und verschiebt
      sie danach an `--output`; ein Abbruch während des Schreibens, auch nach dem
      Zwangsende des Herunterfahrens, lässt unter `--output` keine teilweise Datei zurück
      (Test).
- [ ] [`LH-FA-08`](../../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration): Ein vorhandenes `--output` wird ohne `--force` abgelehnt
      (Exit-Code 2, `SPEC-014`); `--force` nimmt wie jede boolesche Option nur `true` oder
      `false`, über Kommandozeile, Umgebungsvariable und Konfigurationsdatei (Test). Beleg
      in §7 für beide Punkte: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
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
| `internal/adapters/driven/recording` | update | temporäre Datei, atomares Verschieben |
| `internal/adapters/driving/cli` | update | `--output` ohne `--force` abgelehnt; strenge Werte für `--force` über den allgemeinen Leser |
| `internal/adapters/driven/recording` (Unit-Tests), `test/integration` | update | Happy/Boundary/Negative nach `LH-FA-07.a` und `LH-FA-17.a` |
| `docs/user/benutzerhandbuch.md` | update | atomares Schreiben, `--output` und `--force` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-konfiguration` und
`slice-v1-abschluss-konfigurationsdatei` liegen in `done/` (`--force` wird über den allgemeinen
Leser gelesen, auch aus der Datei). Schritt 5 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md)
(Entscheidung des Nutzers vom 2026-10-08). Vor dem ersten Code-Commit prüft der Architect
die Randformen aus §6 und entscheidet die offenen (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Atomares Schreiben verlangt
  plattformspezifische Lösungen — zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Das Verschieben ist auf einer
  Zielplattform nicht atomar zu haben, und die Spezifikation sagt nicht, was dann gilt;
  dann zuerst die Entscheidung, gegebenenfalls ein Carveout.

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

- **Ort und Name der temporären Datei** (im Zielverzeichnis, damit das Verschieben auf
  demselben Dateisystem bleibt) — offen: Der Architect entscheidet vor dem Code in der
  Spezifikation.
- **Übrig gebliebene temporäre Datei eines früheren Laufs** — offen: Der Architect
  entscheidet vor dem Code, ob sie überschrieben oder gemeldet wird.
- **`--output` ist ein Verzeichnis oder liegt in einem fehlenden Verzeichnis** — offen: Der
  Architect bestätigt vor dem Code den Code (`PGR-E3…` nach `SPEC-016` oder `PGR-E2…` nach
  `SPEC-014`).
- **Vorhandenes `--output` ohne `--force`** — Exit-Code 2; entschieden in `SPEC-014`.
- **Werte von `--force`** — nur `true` oder `false`; entschieden in `LH-FA-17.a`.
- **Fehlschlag des Verschiebens** — offen: Der Architect entscheidet vor dem Code Code und
  Verbleib der temporären Datei.

**Risiken:**

- Atomarität des Verschiebens ist plattformabhängig („bestmöglich atomar“; übernommen aus
  `slice-v1-abschluss-betrieb`) — **Ausgang:** offen bis Closure.

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
