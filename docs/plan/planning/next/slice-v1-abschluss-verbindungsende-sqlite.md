# Slice slice-v1-abschluss-verbindungsende-sqlite: Verbindungsende in SQLite: Leser (`meta.version` 1 und 2, `connection_closed`)

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

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-22.a` · `SPEC-001` · `SPEC-002` · `SPEC-043` · `ARC-008`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12. Geschnitten am 2026-10-15 (Planner, vor dem Anspruch): Der Schreiber mit `meta.version` in der Transaktion und der Roundtrip beider Formate sind in `slice-v1-abschluss-verbindungsende-sqlite-schreiber` aufgegangen.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der SQLite-Leser trägt das Ende des Servers wie der YAML-Leser (`slice-v1-abschluss-verbindungsende-format`): Das Ende steht als Nachricht mit `kind` `connection_closed` und `fields` `{}` am Ende der letzten Interaktion der Session. Der Leser nimmt `meta.version` 1 und 2 an (auch 2 ohne Ende); `connection_closed` in einer Datei der Version 1, an einer anderen Stelle als der Form, nicht nach `FATAL`/`PANIC` oder mit Feldern ist `PGR-E3003`, jede andere Version `PGR-E3002` (`LH-FA-25.a` *Version*, `SPEC-001`, `SPEC-043`). Die Dateien dafür baut der Test von Hand.

**Übernommen aus `slice-v1-abschluss-verbindungsende-format`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): SQLite als zweites Format des Endes. Das Modell mit seinen Formregeln und der YAML-Teil liegen in jenem Slice; dieser Slice benutzt `Validate` und die Fehlerklassen unverändert und setzt ihn voraus.

**Bindung aus `slice-v1-abschluss-verbindungsende-format`** (Closure des Gebers, Planner, 2026-10-15, `AGENTS.md` §3.13; dort §6 *Ort der Formregel einer Session*): Die Formregel einer Session (höchstens ein Ende, in der letzten Interaktion) liegt als `Session.ValidateEnd() error` im Modell, und der SQLite-Leser ruft sie je Session nach den Interaktionen, statt sie nachzubauen; `Interaction.Validate` allein lehnt ein Ende in einer Interaktion vor der letzten und ein zweites Ende nicht ab. Die Konstante `model.VersionMitEnde` (2) benutzt der Schreiber in `slice-v1-abschluss-verbindungsende-sqlite-schreiber`; der Leser vergleicht die Version mit `model.Version` und `model.VersionMitEnde`, nicht mit einem Literal.

**Abgegeben an `slice-v1-abschluss-verbindungsende-sqlite-schreiber`** (Schnitt des Planners, 2026-10-15, `AGENTS.md` §3.13): der Schreiber (`meta.version` = 1 beim ersten Schreibvorgang, 2 in der Transaktion der ersten Session mit Ende, Fehler in der Transaktion) und der Roundtrip beider Formate. Der Schreiber liegt hinter diesem Slice, weil sein Roundtrip über den Leser lädt; dieser Slice braucht den YAML-Schreiber (`slice-v1-abschluss-verbindungsende-format-schreiber`) nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Modell, seine Formregeln und der YAML-Leser — `slice-v1-abschluss-verbindungsende-format`; der YAML-Schreiber — `slice-v1-abschluss-verbindungsende-format-schreiber`: Vorgänger in der Reihenfolge; hier nur der SQLite-Leser.
- Der SQLite-Schreiber, `meta.version` in der Transaktion und der Roundtrip beider Formate — `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (dort §1 und DoD): anderer Vertrag (Transaktion, Einspeisung eines Fehlers), und mit ihm läge dieser Plan über der Marke; der Test dieses Slice baut die Dateien von Hand.
- Ein Ende, das `record` aus einem echten Upstream aufzeichnet, in einer SQLite-Datei — `slice-v1-abschluss-verbindungsende-upstream` (dort §1 und DoD, E2E von Abnahmeszenario 18).
- Wiedergabe, Einspielen und Vergleich — `-replay`, `-einspielen`, `-vergleich`: andere Schichten und andere Verträge.
- Änderung am neutralen Schema `tools/schema/schema.yaml` und am Schema-Text — Entscheidung des Nutzers vom 2026-10-12 (bestätigt): `connection_closed` ist ein Wert der Spalte `message.kind`; `make schema-validate` bleibt unberührt.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (`AGENTS.md` §3.11, §3.13).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), `SPEC-001`, `SPEC-043`: Der SQLite-Leser nimmt `meta.version` 1 und 2 an, auch Version 2 ohne Ende. `connection_closed` in einer Datei der Version 1, an einer anderen Stelle als der Form, nicht nach `FATAL`/`PANIC` oder mit Feldern sind `PGR-E3003`; ein Ende in einer Interaktion vor der letzten und ein zweites Ende in einer Session sind `PGR-E3003` über `Session.ValidateEnd`; jede andere Version ist `PGR-E3002` (Test je Randform aus §6, die Dateien von Hand gebaut).
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
| `internal/adapters/driven/recording/sqlite_lesen.go`, `sqlite_lesen_baum.go`, `sqlite_json_lesen.go` | update | Leser: beide Versionen, Nachricht `connection_closed`, Formregeln über `Validate` und `Session.ValidateEnd` (Bindung aus `slice-v1-abschluss-verbindungsende-format`) |
| `internal/adapters/driven/recording` (Tests) | update | Lesen: Version 1 und 2, Negative je Randform, Dateien von Hand gebaut |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` liegt in `done/` (Voraussetzung; Schritt 35a: Modell, `Session.ValidateEnd`, YAML-Leser); `slice-v1-abschluss-verbindungsende-format-schreiber` (35b) liegt ebenfalls in `done/`, ist aber für diesen Slice keine Voraussetzung. Dieser Slice ist Schritt 36a der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat fa8486b..<letzter Commit des Slice> -- internal docs/user`, Tests eingeschlossen, wobei die Basis `fa8486b` ist, der Stand des Hauptzweigs vor dem Anspruch (der Übergang nach `in-progress/` ändert nichts unter den genannten Pfaden); **über 700** geht der Slice zurück zur Zerlegung. **Meldeschwelle 550:** Erreicht die Summe 550, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik: 260 bis 420, Mitte 340 und damit 49 % der Marke (§8).
- `in-progress` → `next`: Die Version oder die Formregel verlangt eine Änderung am Port `RecordingRepository` oder am Modell — zurück zur Zerlegung (dritte Schicht).
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

- *Version 2 ohne Ende:* der Leser nimmt sie an — `SPEC-001`.
- *`connection_closed` in einer Datei der Version 1, nicht als letzte Nachricht, nicht nach `FATAL`/`PANIC`, mit Feldern ungleich `{}`:* `PGR-E3003` — `SPEC-043`, `SPEC-001`.
- *Version 3 und andere:* `PGR-E3002` — `SPEC-001`.
- *Ende mit `offset_ms`:* wie jede Interaktion — `SPEC-043`, `LH-FA-21.a`.
- *Formregel der Session und Konstante der Version 2* (Bindung aus `slice-v1-abschluss-verbindungsende-format`): Der Leser ruft `Session.ValidateEnd` je Session; ein Ende in einer Interaktion vor der letzten und ein zweites Ende in einer Session sind `PGR-E3003`, welche Prüfung meldet, ist nicht zugesagt (`SPEC-002`); `meta.version` 2 heißt `model.VersionMitEnde`, ein Literal ist es nicht — `SPEC-001`, `SPEC-002`.
- *Aufruf-Voraussetzung und Aufrufarten des Treibers:* der Leser ist eine Zustandsprüfung einer Datei ohne Zusage über Reihenfolge oder Gleichzeitigkeit von Aufrufen; kein Hook und kein Gate liest ein Ergebnis dieses Slice — keine (`AGENTS.md` §3.12).

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`): Der Planner hat vor dem Anspruch geschnitten (Schätzung 260 bis 420, Mitte 49 % der Marke); die Rückführung in §4 benennt die Marke und die Meldeschwelle — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand). Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (beide `verkörpert`, `AGENTS.md` §3.14) — dieser Plan nennt Metrik, Meldeschwelle und Marke in §4 und ist vor dem Anspruch geschnitten, weil die Schätzung des Ganzen über 70 % der Marke lag. Keine weiteren Treffer.

Nachgezählt beim Schnitt (Planner, 2026-10-15, `AGENTS.md` §3.13): Liefer-Punkte: 1 (Leser mit beiden Versionen und den Formregeln des Endes). Schichten: eine (Recording-Adapter, `ARC-008`). Schätzung des Planners in der Metrik von §4: Produktcode 40 bis 60 (Version annehmen, Nachricht `connection_closed`, Aufruf von `Session.ValidateEnd`), Tests mit dem Faktor Tests zu Produktcode von mindestens 5 (gemessene Datenpunkte 8 und 11) 200 bis 360, Tabellen etwa 10: **260 bis 420, Mitte 340, 49 % der Marke 700**. Das Ganze vor dem Schnitt lag bei 760 bis 1220 gegen die Marke 1000 (Mitte 990, 99 %).

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-verbindungsende-format-schreiber` und der Bindung aus `slice-v1-abschluss-verbindungsende-format` (Planner, 2026-10-15): Der YAML-Schreiber ist hier keine Voraussetzung mehr (er gehört zum Roundtrip in `-sqlite-schreiber`); der Aufruf von `Session.ValidateEnd` gehört zum Leser-Punkt, `model.VersionMitEnde` zum Vergleich der Version. Liefer-Punkte weiter 1, Schichten weiter eine.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
