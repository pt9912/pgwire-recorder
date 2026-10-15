# Slice slice-v1-abschluss-verbindungsende-sqlite-schreiber: Verbindungsende in SQLite: Schreiber (`meta.version` in der Transaktion der ersten Session mit Ende) und Roundtrip beider Formate

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0013](../../adr/0013-sqlite-recording-backend.md)

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-22.a` · `SPEC-001` · `SPEC-002` · `SPEC-043` · `ARC-008`

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-15.

**Herkunft:** Teil 2 des Schnitts von `slice-v1-abschluss-verbindungsende-sqlite` vom 2026-10-15 (Planner, vor dem Anspruch: Schätzung des Ganzen 760 bis 1220 in der Metrik von §4, Mitte 990, 99 % der Marke 1000, über der Grenze von 70 %). Die Sendung „Schreiber und Roundtrip“ stammt aus §1 und DoD jenes Slice; er liefert den Leser und liegt vor diesem in `done/` (Voraussetzung: Der Roundtrip lädt über seinen Leser). Schritt 36b der Reihenfolge in §5 von welle-v1-abschluss.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der SQLite-Schreiber legt beim ersten Schreibvorgang `meta.version` = 1 an, und die Transaktion der ersten Session mit Ende setzt sie in derselben Transaktion auf 2 (`model.VersionMitEnde`); scheitert die Transaktion, bleibt die Version 1 und die Datei trägt nur die zuvor ergänzten Sessions (`LH-FA-25.a` *Version*, `SPEC-043`). Das Ende steht als Nachricht mit `kind` `connection_closed` und `fields` `{}` in der Reihenfolge der Nachrichten der Interaktion. Dieselbe Aufzeichnung mit Ende liest in beiden Formaten gleich (Roundtrip).

**Übernommen aus `slice-v1-abschluss-verbindungsende-format`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): SQLite als zweites Format des Endes. Das Modell mit seinen Formregeln und der YAML-Teil liegen in jenem Slice; dieser Slice benutzt `Validate` und die Fehlerklassen unverändert und setzt ihn voraus.

**Bindung aus `slice-v1-abschluss-verbindungsende-format`** (Closure des Gebers, Planner, 2026-10-15, `AGENTS.md` §3.13; dort §6 *Ort der Formregel einer Session*): Der Schreiber setzt `meta.version` auf `model.VersionMitEnde` (2); `model.Version` bleibt 1. Die Formregel einer Session (`Session.ValidateEnd`) ruft der Leser, nicht der Schreiber.

**Bindung aus `slice-v1-abschluss-verbindungsende-format-schreiber`** (dort §6 *Aufrufer mit verletzter Form*, Setzung des Architect, 2026-10-15; Eintrag des Planners, 2026-10-15, `AGENTS.md` §3.13): Der Schreiber prüft die Form des Endes nicht — Stelle, Anzahl, `ready_for_query` daneben, Felder neben `type`; er schreibt, was die Aufzeichnung trägt, und eine verletzte Form ergibt eine Datei, die der Leser (`slice-v1-abschluss-verbindungsende-sqlite`) mit `PGR-E3003` ablehnt (Test: schreiben nimmt an, laden lehnt ab). Ob der SQLite-Schreiber dieselbe Setzung trägt, bestätigt der Architect vor dem Code (Setzung des Planners, keine Entscheidung des Nutzers); das Gewicht liegt beim Aufrufer `slice-v1-abschluss-verbindungsende-record`, das die Bindung in §1 trägt.

**Übernommen aus `slice-v1-abschluss-verbindungsende-format-schreiber`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-15, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-format`): der YAML-Schreiber mit Version nach dem Inhalt. Der Roundtrip in beiden Formaten (DoD-Punkt 2) schreibt dieselbe Aufzeichnung als YAML und als SQLite; den YAML-Schreiber liefert jener Slice, dieser benutzt ihn und setzt ihn voraus.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der SQLite-Leser (beide Versionen, Formregeln, `PGR-E3002`/`PGR-E3003`) — `slice-v1-abschluss-verbindungsende-sqlite`: Vorgänger in der Reihenfolge; dieser Slice lädt über ihn und ändert ihn nicht.
- Das Modell, seine Formregeln und der YAML-Leser — `slice-v1-abschluss-verbindungsende-format`; der YAML-Schreiber — `slice-v1-abschluss-verbindungsende-format-schreiber`: Vorgänger; hier nur der SQLite-Schreiber und der Roundtrip.
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
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat): Roundtrip in beiden Formaten: dieselbe Aufzeichnung mit Ende — einfach, Extended mit `Sync`- und mit `Flush`-Gruppe, mit `offset_ms` — wird als YAML und als SQLite geschrieben und liest gleich; ein späterer Schreibvorgang an eine Datei der Version 2 ergänzt Sessions, ohne die Version zu ändern; eine verletzte Form des Endes nimmt der Schreiber an, der Leser lehnt sie mit `PGR-E3003` ab (Test).
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
| `internal/adapters/driven/recording` (Tests) | update | Schreiben (Version, Transaktion mit Fehler), Roundtrip beider Formate |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-sqlite` (Leser, Schritt 36a), `slice-v1-abschluss-verbindungsende-format` (35a) und `slice-v1-abschluss-verbindungsende-format-schreiber` (35b, YAML-Schreiber für den Roundtrip) liegen in `done/` (Voraussetzung). Dieser Slice ist Schritt 36b der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (der Planner trägt den SHA beim Anspruch ein; der Übergang nach `in-progress/` ändert nichts unter den genannten Pfaden); **über 1000** geht der Slice zurück zur Zerlegung. **Meldeschwelle 800:** Erreicht die Summe 800, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik: 500 bis 800, Mitte 650 und damit 65 % der Marke (§8); schneidet die Prüfung des Architect höher, ist der Schnitt Version in der Transaktion | Roundtrip (Zuschnitt: Planner).
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
- *Ende mit `offset_ms`:* wie jede Interaktion — `SPEC-043`, `LH-FA-21.a`.
- *Jeder Schreibvorgang öffnet die Datei neu und prüft sie:* die Prüfung nimmt Version 2 an — `SPEC-043` (Schreibvorgang).
- *Aufrufer mit verletzter Form* (Bindung aus `slice-v1-abschluss-verbindungsende-format-schreiber`): der Schreiber prüft die Form des Endes nicht; Stelle, Anzahl und Felder entscheidet der Leser — Architect bestätigt vor dem Code, ob die Setzung für SQLite gilt (`SPEC-001`, `SPEC-002`).
- *Konstante der Version 2* (Bindung aus `slice-v1-abschluss-verbindungsende-format`): `model.VersionMitEnde`, ein Literal ist es nicht — `SPEC-001`.
- *Aufruf-Voraussetzung:* der Schreiber vergibt keine Werte, die über mehrere Aufrufe einer Session geordnet oder eindeutig sein müssten; ein Schreibvorgang an dieselbe Datei läuft unter der Sperre der Datei (`LH-FA-22.a`); keine Ordnungszusage — keine Voraussetzung (`AGENTS.md` §3.12). Aufrufarten des Treibers: keine (kein Hook oder Gate liest ein Ergebnis dieses Schreibers).

**Risiken:**

- Ein Fehler in der Transaktion ist ohne Einspeisung nicht erreichbar; ob die vorhandenen Tests des Schreibers ein Mittel dafür haben (etwa in `sqlite_einstellungen_test.go`), prüft der Architect vor dem Code; fehlt es, ist das Mittel Teil dieses Slice — **Ausgang:** offen bis Closure.
- Der Diff liegt über der Marke (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`: Tests mit Mutationen je Element, Roundtrip je Form; Schätzung 500 bis 800 bei Meldeschwelle 800) — **Ausgang:** offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand). Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (beide `verkörpert`, `AGENTS.md` §3.14) — dieser Plan nennt Metrik, Meldeschwelle und Marke in §4 und ist vor dem Anspruch geschnitten. `BEO-REPO/meldeschwelle-ohne-artefakt`: §4 nennt die Meldeschwelle. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-15, `AGENTS.md` §3.13): Liefer-Punkte: 2 (Schreiber mit Version in der Transaktion, Roundtrip). Schichten: eine (Recording-Adapter, `ARC-008`). Schätzung des Planners in der Metrik von §4: Produktcode 80 bis 110 (Version in der Transaktion, Nachricht `connection_closed`, Einspeisung des Fehlers), Tests mit dem Faktor Tests zu Produktcode von mindestens 5 (gemessene Datenpunkte 8 und 11) 400 bis 550, Roundtrip-Fälle und Tabellen 20 bis 140: **500 bis 800, Mitte 650, 65 % der Marke 1000**, oberes Ende 80 %; kein Schnitt, die Meldeschwelle 800 trägt.

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-verbindungsende-format-schreiber` und der Bindungen aus `slice-v1-abschluss-verbindungsende-format` und `-format-schreiber` (Planner, 2026-10-15): Der YAML-Schreiber ist Voraussetzung und keine Lieferung; die Konstante `model.VersionMitEnde` gehört zu DoD-Punkt 1, die Bindung „Schreiber prüft nicht“ zu DoD-Punkt 2 (ein Fall „schreiben nimmt an, laden lehnt ab“ je Stelle aus `format-schreiber` §6). Liefer-Punkte weiter 2, Schichten weiter eine.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
