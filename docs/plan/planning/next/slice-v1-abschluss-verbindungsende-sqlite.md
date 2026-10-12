# Slice slice-v1-abschluss-verbindungsende-sqlite: Verbindungsende in SQLite: `meta.version` in der Transaktion der ersten Session mit Ende

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0013](../../adr/0013-sqlite-recording-backend.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-22.a` · `SPEC-001` · `SPEC-043` · `ARC-008`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der SQLite-Schreiber und der SQLite-Leser tragen das Ende des Servers wie das YAML-Format (`slice-v1-abschluss-verbindungsende-format`): Das Ende steht als Nachricht mit `kind` `connection_closed` und `fields` `{}` am Ende der letzten Interaktion der Session. Der erste Schreibvorgang legt `meta.version` = 1 an, und die Transaktion der ersten Session mit Ende setzt sie in derselben Transaktion auf 2; scheitert die Transaktion, bleibt die Version 1 und die Datei trägt nur die zuvor ergänzten Sessions (`LH-FA-25.a` *Version*, `SPEC-043`). Der Leser nimmt beide Versionen an. Dieselbe Aufzeichnung mit Ende liest in beiden Formaten gleich (Roundtrip).

**Übernommen aus `slice-v1-abschluss-verbindungsende-format`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): SQLite als zweites Format des Endes. Das Modell mit seinen Formregeln und der YAML-Teil liegen in jenem Slice; dieser Slice benutzt `Validate` und die Fehlerklassen unverändert und setzt ihn voraus.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Modell, seine Formregeln und das YAML-Format — `slice-v1-abschluss-verbindungsende-format`: Vorgänger in der Reihenfolge; hier nur der Schreiber, der Leser und das Schema-Verhalten von `sqlite`.
- Ein Ende, das `record` aus einem echten Upstream aufzeichnet, in einer SQLite-Datei — `slice-v1-abschluss-verbindungsende-upstream` (dort §1 und DoD, E2E von Abnahmeszenario 18): Hier schreibt der Test das Modell mit Ende von Hand; der Upstream-Adapter liefert es erst dort.
- Wiedergabe, Einspielen und Vergleich — `-replay`, `-einspielen`, `-vergleich`: andere Schichten und andere Verträge.
- Änderung am neutralen Schema `tools/schema/schema.yaml` und am Schema-Text — Entscheidung des Nutzers vom 2026-10-12 (bestätigt): `connection_closed` ist ein Wert der Spalte `message.kind`; `make schema-validate` bleibt unberührt.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), `SPEC-043`: Der SQLite-Schreiber legt `meta.version` = 1 an und setzt sie in der Transaktion der ersten Session mit Ende auf 2; die Transaktion einer Session ohne Ende, die danach kommt, und jede spätere ändert sie nicht, eine Datei ohne Ende bleibt bei 1. Das Ende steht als Nachricht `connection_closed` mit `fields` `{}` in der Reihenfolge der Nachrichten der Interaktion. Scheitert die Transaktion der ersten Session mit Ende, bleibt die Version 1 und die Datei trägt nur die zuvor ergänzten Sessions; der nächste Schreibvorgang ergänzt, was fehlt (Test mit eingespeistem Fehler).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), `SPEC-043`: Der SQLite-Leser nimmt `meta.version` 1 und 2 an, auch Version 2 ohne Ende. `connection_closed` in einer Datei der Version 1, an einer anderen Stelle als der Form, nicht nach `FATAL`/`PANIC` oder mit Feldern sind `PGR-E3003`; jede andere Version ist `PGR-E3002` (Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat): Roundtrip in beiden Formaten: dieselbe Aufzeichnung mit Ende — einfach, Extended mit `Sync`- und mit `Flush`-Gruppe, mit `offset_ms` — wird als YAML und als SQLite geschrieben und liest gleich; ein späterer Schreibvorgang an eine Datei der Version 2 ergänzt Sessions, ohne die Version zu ändern (Test).
- [ ] Je Zusage des neuen Vertrags eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/recording/sqlite.go`, `sqlite_zeilen.go`, `sqlite_json.go` | update | Schreiber: Nachricht `connection_closed`, `meta.version` in der Transaktion |
| `internal/adapters/driven/recording/sqlite_lesen.go`, `sqlite_lesen_baum.go`, `sqlite_json_lesen.go` | update | Leser: beide Versionen, Nachricht `connection_closed`, Formregeln über `Validate` |
| `internal/adapters/driven/recording` (Tests) | update | Schreiben (Version, Transaktion mit Fehler), Lesen (Negative je Randform), Roundtrip beider Formate |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` liegt in `done/` (Voraussetzung; Schritt 35); dieser Slice ist Schritt 36 der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Prüfung des Architect schätzt den Diff mit Tests auf über 1000 Zeilen (Schätzung des Planners 400 bis 750) — Schnitt vor dem Anspruch: Schreiber | Leser, Zuschnitt: Planner (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`).
- `in-progress` → `next`: Das Setzen der Version in der Transaktion verlangt eine Änderung am Port `RecordingRepository` oder am Modell — zurück zur Zerlegung (dritte Schicht).
- `in-progress` → `open`: Das neutrale Schema muss sich ändern (der Nutzer hat es am 2026-10-12 ausgeschlossen) — Change Request.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Setzungen des Architect). Ort: `LH-FA-25.a` *Version*, `SPEC-001`, `SPEC-043`, `LH-FA-22.a`.

- *Erster Schreibvorgang:* legt `meta.version` = 1 an, auch wenn die erste Session ein Ende trägt; die Transaktion dieser Session setzt sie auf 2 — `LH-FA-25.a` *Version*.
- *Transaktion der ersten Session mit Ende scheitert:* zurückgerollt, Version bleibt 1, der nächste Schreibvorgang ergänzt — `LH-FA-25.a` *Version*, `SPEC-043` (Schreibvorgang).
- *Zweite und spätere Sessions mit Ende, Session ohne Ende nach der Version 2:* ändern die Version nicht — `SPEC-043`.
- *Version 2 ohne Ende:* der Leser nimmt sie an — `SPEC-001`.
- *`connection_closed` in einer Datei der Version 1, nicht als letzte Nachricht, nicht nach `FATAL`/`PANIC`, mit Feldern ungleich `{}`:* `PGR-E3003` — `SPEC-043`, `SPEC-001`.
- *Version 3 und andere:* `PGR-E3002` — `SPEC-001`.
- *Ende mit `offset_ms`:* wie jede Interaktion — `SPEC-043`, `LH-FA-21.a`.
- *Jeder Schreibvorgang öffnet die Datei neu und prüft sie:* die Prüfung nimmt Version 2 an — `SPEC-043` (Schreibvorgang).

**Risiken:**

- Ein Fehler in der Transaktion ist ohne Einspeisung nicht erreichbar; ob die vorhandenen Tests des Schreibers ein Mittel dafür haben (etwa in `sqlite_einstellungen_test.go`), prüft der Architect vor dem Code; fehlt es, ist das Mittel Teil dieses Slice — **Ausgang:** offen bis Closure.
- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×): die Rückführung in §4 benennt den Schnitt — **Ausgang:** offen bis Closure.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×; Belege `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen`, die beiden Vorgänger dieses Formats) — der Plan nennt seine Schätzung mit Tests und den Schnitt vorab; wird der Eintrag mit diesem Slice zum dritten Mal belegt, ist er eine Lücke. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13): Liefer-Punkte: 3 (Schreiber, Leser, Roundtrip). Schichten: eine (Recording-Adapter, `ARC-008`). Schätzung des Planners: 400 bis 750 Zeilen mit Tests.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
