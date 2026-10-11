# Slice slice-v1-abschluss-sqlite-format: SQLite als Aufzeichnungsformat

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [ADR-0005](../../adr/0005-recording-store-ist-driven-adapter.md), [ADR-0013](../../adr/0013-sqlite-recording-backend.md)

**Berührte Spec-Stellen:** `LH-FA-22.a` · `LH-FA-07.a` · `SPEC-001` · `SPEC-043` · `SPEC-044` · `ARC-008` · `ARC-014`

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `record --format sqlite` speichert die Aufzeichnung je Session in einer Transaktion in einer SQLite-Datei, und Replay und Einspielen erkennen das Format der Datei selbst.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-sqlite-format-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- Umwandlung zwischen den Formaten und weitere Formate — Out-of-Scope von LH-FA-22.
- Das Standardformat YAML — bleibt unverändert aus `slice-walking-skeleton-record`.

**Übernommen von `slice-v1-abschluss-schreiben`:** Was von `LH-FA-07.a` *Temporäre Datei*
(Name, Rechte, Fehlschlag, übrig gebliebene Datei) und von der Probedatei aus *Zielpfad beim
Start* für `sqlite` gilt, nennt §6 dieses Slice als Randformen und entscheidet der Architect
in `LH-FA-22.a` *Schreiben von `sqlite`*, vor dem ersten Code-Commit (`AGENTS.md` §3.12).
Dort gilt heute nur: Zielpfad wie in `LH-FA-07.a`, temporäre Datei im Verzeichnis der
Zieldatei, atomar verschoben.
Dazu, mit der Kennung `slice-v1-abschluss-schreiben` (Entscheidung des Nutzers vom
2026-10-09): die Grenze der Prüfung für `sqlite` — ob für das atomare Verschieben in
`LH-FA-22.a` *Schreiben von `sqlite`* dieselbe Grenze gilt wie in `LH-FA-07.a` Schritt 2 (geprüft
nur unter Linux im selben Dateisystem, macOS und Windows nicht) und ob sie für das Ergänzen
einer Session in einer Transaktion eine eigene braucht. §6 nennt sie als Randform, der Architect
dieses Slice entscheidet sie in `LH-FA-22.a`, vor dem ersten Code-Commit (`AGENTS.md` §3.12).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat): Eine Aufzeichnung mit `--format sqlite` liefert im Replay und beim Einspielen dasselbe Verhalten wie dieselbe Aufzeichnung als YAML (Roundtrip-Gleichheit, Abnahmeszenario 14); die Datei enthält zu jedem Zeitpunkt nur vollständige Sessions.
- [ ] [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats): Eine Datei mit unbekannter Version (`PGR-E3002`) oder ohne gültige Aufzeichnung (`PGR-E3003`) wird in beiden Formaten erkannt (Test).
- [ ] Die SQLite-Bibliothek ist im Architektur-Gate (`.a-check.yml`, `tech`-Regel) auf den Recording-Adapter begrenzt; eine absichtliche Verletzung lässt `make a-check` fehlschlagen (Test).
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
| `internal/adapters/driven/recording` | update | zweiter Adapter hinter `RecordingRepository`: SQLite, Formaterkennung, Transaktion je Session |
| `internal/adapters/driving/cli` | update | Option `--format` |
| `tools/schema/schema.yaml`, `harness/mk/schema.mk` | vorhanden | neutrales Schema der Tabellenform (d-migrate); das SQL für SQLite wird daraus erzeugt und im Adapter eingebettet |
| `.a-check.yml` (`tech`-Regel) | prüfen | Die Regel für `modernc.org/sqlite` liegt vor; der Slice belegt sie mit einer absichtlichen Verletzung |
| `test/integration` | update | Roundtrip-Gleichheit beider Formate |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-walking-skeleton-record` ist `done` (YAML-Format liegt vor), und [ADR-0013](../../adr/0013-sqlite-recording-backend.md) ist `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: die SQLite-Bibliothek verlangt native Abhängigkeiten für die Zielplattformen — zurück zur Zerlegung.
- `in-progress` → `open`: Beide Formate tragen das Modell nicht gleich — Entscheidung klären.

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

- **Randform** (aus `slice-v1-abschluss-schreiben`, §1 *Übernommen*): Grenze der Prüfung des
  atomaren Verschiebens und der Transaktion je Session unter Linux, macOS und Windows — offen,
  entscheidet der Architect dieses Slice in `LH-FA-22.a` vor dem ersten Code-Commit —
  **Ausgang:** offen bis zu seiner Prüfung.
- Die gewählte Bibliothek (reines Go, siehe [ADR-0024](../../adr/0024-sqlite-bibliothek.md)) ist für die sechs Zielplattformen erst durch den Build im Slice belegt — **Ausgang:** offen bis Closure.

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

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-sqlite-format-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 3. Schichten: zwei Schichten (Recording-Adapter, CLI-Adapter; das Schema und die Regel in `.a-check.yml` liegen vor und werden geprüft).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
