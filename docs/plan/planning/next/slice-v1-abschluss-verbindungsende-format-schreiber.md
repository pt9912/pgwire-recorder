# Slice slice-v1-abschluss-verbindungsende-format-schreiber: Verbindungsende im Format: YAML-Schreiber und Roundtrip (Version 2)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-22.a` · `SPEC-001` · `ARC-001` · `ARC-008`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der YAML-Schreiber schreibt `version: 2` genau dann, wenn die Aufzeichnung ein Ende enthält (`Recording.HasEnd()`), sonst `version: 1`, bei jedem Schreiben nach dem Inhalt und nicht nach dem Feld `version` des Modells (`LH-FA-25.a` *Version*, `SPEC-001`). Der Roundtrip (Laden, Schreiben, Laden: gleich) gilt für Aufzeichnungen mit Ende, einfach, Extended mit `Sync`-Gruppe und mit `Flush`-Gruppe, mit `offset_ms` und mit `empty_sessions`. Nachweis ohne `record`: Aufzeichnungen als Text im Test, über den Leser von `slice-v1-abschluss-verbindungsende-format`; am Binary ändert sich nichts, bis `slice-v1-abschluss-verbindungsende-record` ein Ende aufzeichnet.

**Herkunft:** Teil 2 des Schnitts vom 2026-10-15 (Rückführung von `slice-v1-abschluss-verbindungsende-format`: Prüfung des Architect vor dem Code, ungeteilt 850 bis 1250 in der Metrik von §4, über der Marke 1000; Empfehlung Leser | Schreiber, Zuschnitt des Planners). Die Sendung „Schreiber und Roundtrip“ stammt aus §1 und DoD jenes Slice; er liefert Modell, `Session.ValidateEnd` und den Leser und liegt vor diesem in `done/` (Voraussetzung: Der Roundtrip lädt über seinen Leser). Schritt 35b der Reihenfolge in §5 von welle-v1-abschluss.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Modell (`connection_closed`, Formregeln, `Session.ValidateEnd`), die Konstante der Version 2 und der YAML-Leser — `slice-v1-abschluss-verbindungsende-format` (dort §1 und DoD): Vorgänger und Voraussetzung; mit ihnen läge dieser Plan wieder über der Marke.
- SQLite (Schreiber, Leser, `meta.version`, Roundtrip beider Formate) — `slice-v1-abschluss-verbindungsende-sqlite` (dort §1 und DoD): zweites Format mit eigenen Randformen (Transaktion); es benutzt den Schreiber dieses Slice für den Roundtrip in beiden Formaten und setzt ihn voraus.
- Aufzeichnen des Endes — `slice-v1-abschluss-verbindungsende-record` (Use Case, PGWire-Adapter) und `slice-v1-abschluss-verbindungsende-upstream` (Upstream-Adapter): Der `RecordService` setzt die Version nicht; sie folgt beim Schreiben dem Inhalt, und dieser Slice liefert diese Wahl. Der Slice `-record` rechnet mit Fakes und braucht den Schreiber nicht; `-upstream` zeichnet in YAML auf und setzt ihn voraus.
- Änderung am neutralen Schema `tools/schema/schema.yaml` — Entscheidung des Nutzers vom 2026-10-12 (Darstellung bestätigt): Das Ende ist ein Wert von `message.kind`, keine neue Spalte.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (dort §1 und DoD), nach dem letzten Code-Slice des Endes in derselben Welle (`AGENTS.md` §3.11, §3.13): Mit der Dokumentation als dritter Schicht läge dieser Plan über zwei Schichten, und ein Ende kann ein Anwender erst mit `record` erzeugen.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), `SPEC-001`: `Recording.HasEnd() bool` ist wahr, wenn irgendeine Interaktion irgendeiner Session ein Ende trägt. Der YAML-Schreiber schreibt `version: 2` mit mindestens einem Ende und sonst `version: 1`, bei jedem Schreiben nach dem Inhalt, nicht nach dem Feld `version` des Modells (Feld 1 mit Ende schreibt 2, Feld 2 ohne Ende schreibt 1); eine Aufzeichnung ohne Ende, auch eine ohne Sessions oder nur mit `empty_sessions`, schreibt `version: 1` (Test; die vorhandenen Tests des Schreibers bleiben grün).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats): Der Roundtrip (Laden, Schreiben, Laden: gleich) gilt für Aufzeichnungen mit Ende — einfach, Extended mit `Sync`-Gruppe und mit `Flush`-Gruppe, mit `offset_ms`, mit `empty_sessions` —; eine Aufzeichnung ohne Ende schreibt unverändert `version: 1` (Test, Fixtures als Text im Test).
- [ ] Je Zusage des neuen Vertrags (Wahl der Version, Roundtrip) eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
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
| `internal/hexagon/model` (`recording.go`) | update | `Recording.HasEnd() bool` (Ort: Entscheidung des Architect in `slice-v1-abschluss-verbindungsende-format` §6) |
| `internal/adapters/driven/recording/yaml.go` | update | Schreiber: Version nach dem Inhalt (`HasEnd`), nicht nach dem Feld `version` |
| `internal/adapters/driven/recording` (Tests) | update | Schreiber und Roundtrip; Fixtures als Text im Test; ein Test von `HasEnd` am Modell |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` liegt in `done/` (Voraussetzung: Modell, Formregeln, Leser); dieser Slice ist Schritt 35b der Reihenfolge in §5 von welle-v1-abschluss und steht vor `slice-v1-abschluss-verbindungsende-sqlite` (dessen Roundtrip beider Formate benutzt den Schreiber). WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal docs/user`, Tests eingeschlossen, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 700** geht der Slice zurück zur Zerlegung. **Meldeschwelle 550:** Erreicht die Summe 550, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Architect für diesen Teil 300 bis 500 (Produktcode etwa 40, Tests Schreiber und Roundtrip etwa 230, Hilfen und Tabellen, vor der Nacharbeit der Mutationen; Faktor 1,3 bis 1,5 gerechnet).
- `in-progress` → `next`: Die Version muss beim Schreiben von einer dritten Schicht gesetzt werden (etwa vom `RecordService`) oder der Roundtrip verlangt eine Änderung am Leser aus `slice-v1-abschluss-verbindungsende-format` — zurück zur Zerlegung.
- `in-progress` → `open`: Der Nutzer ändert die Darstellung des Endes oder die Regel „Version 2 nur bei Ende“ — Change Request an die Spezifikation.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (`LH-FA-25.a` *Version*, `SPEC-001`; die Regel „Formatversion 2 nur bei Ende“ bestätigte der Nutzer am 2026-10-12). Der Leser, die Formregeln und ihre Randformen stehen in `slice-v1-abschluss-verbindungsende-format` §6; hier nur, was der Schreiber und der Roundtrip entscheiden:

- *Version im Modell gegen Inhalt:* das Feld `version` des Modells bestimmt nichts; der Schreiber richtet sich nach dem Inhalt (Feld 1 mit Ende schreibt 2, Feld 2 ohne Ende schreibt 1) — `SPEC-001`.
- *Ende in einer von mehreren Sessions:* ein Ende in irgendeiner Session genügt für `version: 2` (`Recording.HasEnd`); die Formregel der Session (höchstens ein Ende, in der letzten Interaktion) gilt je Session — `LH-FA-25.a` *Version*, `SPEC-001`.
- *Version 2 ohne Ende:* der Schreiber schreibt sie nicht; der Leser nimmt sie an (`slice-v1-abschluss-verbindungsende-format`) — `SPEC-001`.
- *Ohne Sessions, nur `empty_sessions`:* kein Ende, also `version: 1` — folgt aus „genau dann, wenn die Aufzeichnung ein Ende enthält“ (`LH-FA-25.a`).
- *Roundtrip-Formen:* einfach, `Sync`-Gruppe, `Flush`-Gruppe, `offset_ms`, `empty_sessions`; die Interaktion mit Ende trägt `offset_ms` wie jede andere — `SPEC-002`, `LH-FA-12.a`.
- *Aufrufer mit verletzter Formregel (Rückfrage an den Architect vor dem Code):* Ob der Schreiber eine Aufzeichnung schreibt, deren Ende nicht in der letzten Interaktion einer Session steht (heute prüft er den Inhalt nicht; die Prüfung läge bei `Session.ValidateEnd` im Aufrufer), oder sie ablehnt, ist in der Spezifikation nicht entschieden; Empfehlung des Planners: keine Prüfung im Schreiber, der Roundtrip lehnt beim Laden mit `PGR-E3003` ab. Entscheidet der Architect anders, trägt er es in `SPEC-001` ein, bevor der erste Code-Commit liegt.
- *Aufruf-Voraussetzung:* `HasEnd` und die Wahl der Version sind Zustandsprüfungen ohne Zusage über Reihenfolge oder Gleichzeitigkeit von Aufrufen; der Slice nennt keine Ordnung der Werte über mehrere Aufrufe, also gibt es keine Voraussetzung (`AGENTS.md` §3.12). Aufrufarten des Treibers: keine (kein Hook oder Gate liest ein Ergebnis dieses Schreibers).
- *Andere Schreiber:* der SQLite-Schreiber bleibt unberührt (`slice-v1-abschluss-verbindungsende-sqlite`); `model.Version` bleibt 1 für ihn.

**Zählung der Mutationstabelle (Planner nach der Rechnung des Architect aus `slice-v1-abschluss-verbindungsende-format` §6, Prüfung durch den Architect vor dem Code; `AGENTS.md` §3.10):** ein Fall und eine Mutation, die nur ihn bricht, je Element von: der Wahl der Version beim Schreiben (mit Ende 2, ohne 1, Feld gegen Inhalt in beide Richtungen, Ende nur in einer späteren Session, ohne Sessions, nur `empty_sessions`); den Formen des Roundtrips (einfach · `Sync`-Gruppe · `Flush`-Gruppe · `offset_ms` · `empty_sessions`).

**Risiken:**

- Zwischenstand: Nach diesem Slice kann der Schreiber `version: 2` schreiben, aber kein Produktpfad erzeugt ein Ende (`slice-v1-abschluss-verbindungsende-record`, `-upstream`) — **Ausgang:** offen bis Closure (erwartet: eingetreten, `slice-v1-abschluss-verbindungsende-record` und `-upstream`, die das Ende aufzeichnen).
- Die Wahl der Version beim Schreiben ist Voraussetzung des Roundtrips beider Formate in `slice-v1-abschluss-verbindungsende-sqlite` und des Abnahmeszenarios 18 in YAML in `slice-v1-abschluss-verbindungsende-upstream` — **Ausgang:** offen bis Closure (erwartet: eingetreten, `slice-v1-abschluss-verbindungsende-sqlite` und `-upstream`, beide nennen diesen Slice in §1 und im Start).
- Der Diff liegt über der Marke (Schätzung 300 bis 500, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`: Tests mit Mutationen je Element) — **Ausgang:** offen bis Closure; Meldeschwelle 550 und Marke 700 stehen in §4.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand). Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (beide `verkörpert`, `AGENTS.md` §3.14; mit der Rückführung von `slice-v1-abschluss-verbindungsende-format` ein Auftreten trotz der Regel, Urteil im Register) — dieser Plan nennt Metrik, Meldeschwelle und Marke in §4 und liegt nach der Schätzung des Architect unter der Marke. `BEO-REPO/meldeschwelle-ohne-artefakt` (1×): §4 nennt die Meldeschwelle. `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` ist verkörpert (`AGENTS.md` §3.10): die Mutationstabelle steht in DoD und §7. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-15, `AGENTS.md` §3.13): Liefer-Punkte: 2 (Schreiber mit `HasEnd` und Version nach dem Inhalt, Roundtrip). Schichten: zwei (Kern mit dem Modell, `ARC-001`; Recording-Adapter, `ARC-008`); der Handbuch- und README-Teil liegt in `slice-v1-abschluss-verbindungsende-doku`. Schätzung des Architect (2026-10-15, aus dem ungeteilten Plan): 300 bis 500 Zeilen mit Tests in der Metrik von §4. Sendungen anderer Slices an diesen: `slice-v1-abschluss-verbindungsende-sqlite` und `-upstream` nennen ihn als Voraussetzung (dort §1, Start und §8 nachgezählt).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
