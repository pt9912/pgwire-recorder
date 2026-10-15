# Slice slice-v1-abschluss-verbindungsende-format: Verbindungsende im Format: Modell und YAML (Version 2)

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

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-02.b` · `LH-FA-22.a` · `SPEC-001` · `SPEC-002` · `SPEC-041` · `ARC-001` · `ARC-008`

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

**Ziel:** Das Modell kennt das Ende des Servers, und das YAML-Format trägt es: `connection_closed` ohne Felder als letzte Server-Nachricht der letzten Interaktion einer Session, unmittelbar nach einer `error_response` `FATAL` oder `PANIC` (`LH-FA-25.a` *Form*, `SPEC-002`, `SPEC-041`). Der Schreiber schreibt `version: 2` genau dann, wenn die Aufzeichnung ein Ende enthält, sonst `version: 1`; der Leser nimmt beide Versionen an, lehnt `connection_closed` in Version 1 und jede Stelle außerhalb der Form mit `PGR-E3003` ab und jede andere Version mit `PGR-E3002` (`SPEC-001`). Nachweis ohne `record`: handgeschriebene Aufzeichnungen und der Roundtrip; am Binary ändert sich nur das Laden einer solchen Aufzeichnung.

**Herkunft:** Schnitt des Planners nach der Entscheidung des Nutzers vom 2026-10-12 (Option B, [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md)): Die Arbeit ist in acht Slices der Welle geschnitten, dieser ist der erste. Modell und YAML stehen vor allen Slices, die sie benutzen (SQLite, Replay, Einspielen, Record), weil jeder `connection_closed` und die Formregeln des Modells voraussetzt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- SQLite (Schreiber, Leser, `meta.version` in der Transaktion der ersten Session mit Ende) — `slice-v1-abschluss-verbindungsende-sqlite` (dort §1 und DoD): zweites Format mit eigenen Randformen (Transaktion); mit ihm läge dieser Plan über der Grenze einer Review-Sitzung.
- Wiedergabe im Replay — `slice-v1-abschluss-verbindungsende-replay`: Nach diesem Slice lädt `replay` eine Aufzeichnung mit Ende, gibt `connection_closed` aber noch nicht wieder (wellenintern akzeptierter Zwischenstand, ein Release liegt nicht dazwischen).
- Einspielen und Vergleich — `slice-v1-abschluss-verbindungsende-einspielen` und `slice-v1-abschluss-verbindungsende-vergleich`: Sie lesen die Aufzeichnung über die Leser dieses Slice, verhalten sich aber anders als das Format.
- Aufzeichnen des Endes — `slice-v1-abschluss-verbindungsende-record` (Use Case, PGWire-Adapter) und `slice-v1-abschluss-verbindungsende-upstream` (Upstream-Adapter): Der `RecordService` setzt die Version nicht; sie folgt beim Schreiben dem Inhalt (`LH-FA-25.a` *Version*).
- Änderung am neutralen Schema `tools/schema/schema.yaml` — Entscheidung des Nutzers vom 2026-10-12 (Darstellung bestätigt): Das Ende ist ein Wert von `message.kind`, keine neue Spalte.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (dort §1 und DoD), nach dem letzten Code-Slice des Endes in derselben Welle (`AGENTS.md` §3.11, §3.13): Mit der Dokumentation als dritter Schicht läge dieser Plan über zwei Schichten; ein Ende kann ein Anwender erst mit `record` erzeugen, und das Handbuch beschreibt, was das gebaute Binary zeigt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), `SPEC-002`, `SPEC-041`: `connection_closed` ist ein Antworttyp des Modells. `Validate` nimmt ihn nur als letzte Server-Nachricht an, unmittelbar nach einer `error_response` mit dem Schweregrad `FATAL` oder `PANIC` (Feld `V`, ohne es `S`) und ohne Felder: bei einer einfachen Anfrage statt `ready_for_query`, bei einer Extended-Interaktion in der letzten Gruppe, die dann mit `Sync` oder `Flush` endet und kein `ready_for_query` trägt. Eine Session trägt höchstens ein Ende, in ihrer letzten Interaktion (Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), `SPEC-001`: Der YAML-Leser nimmt `version: 1` und `version: 2` an, auch Version 2 ohne Ende. `connection_closed` in Version 1, an einer anderen Stelle als der Form, nicht nach `FATAL`/`PANIC` und mit Feldern sind `PGR-E3003`; `version: 3` und jede andere ganze Zahl sind `PGR-E3002` (Test).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), `SPEC-001`: Der YAML-Schreiber schreibt `version: 2` mit mindestens einem Ende und sonst `version: 1`, bei jedem Schreiben nach dem Inhalt. Der Roundtrip (Laden, Schreiben, Laden: gleich) gilt für Aufzeichnungen mit Ende — einfach, Extended mit `Sync`-Gruppe und mit `Flush`-Gruppe, mit `offset_ms`, mit `empty_sessions` —; eine Aufzeichnung ohne Ende schreibt unverändert `version: 1` (Test; die vorhandenen Tests des Schreibers bleiben grün).
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
| `internal/hexagon/model` (`response.go`, `recording.go`, `extended.go`) | update | Antworttyp `connection_closed`, Version 2, Formregeln einer Interaktion und einer Session |
| `internal/hexagon/model` (Tests) | update | `Validate`: Happy/Boundary/Negative je Randform aus §6 |
| `internal/adapters/driven/recording/yaml.go` | update | Leser: beide Versionen, `PGR-E3002`/`PGR-E3003`; Schreiber: Version nach dem Inhalt |
| `internal/adapters/driven/recording` (Tests) | update | Leser, Schreiber, Roundtrip; Fixtures als Text im Test |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests; Tests, die die alte Regel prüfen (kein Ende im Modell), gehören diesem Slice |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-protokollrand-doku` liegt in `done/` (Schritt 34); dieser Slice ist Schritt 35 der Reihenfolge in §5 von welle-v1-abschluss, der erste des Blocks zum Verbindungsende (Entscheidung des Nutzers vom 2026-10-12). WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Prüfung des Architect vor dem Code schätzt den Diff auf über 1000 Zeilen, gemessen als Summe aus Einfügungen und Löschungen von `git diff --numstat 63c497c..<letzter Commit des Slice> -- internal docs/user`, Tests eingeschlossen (`AGENTS.md` §3.14; Schätzung des Planners 650 bis 1000, Tests eingeschlossen, Verhältnis Test zu Code etwa 3 zu 1 wie bei `slice-v1-abschluss-antwortvergleich`, der Architect rechnet für Tabellen- und Gegenprobe-Code mit dem Doppelten des Produktcodes; `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`) — Schnitt vor dem Anspruch, nicht danach (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`): Leser | Schreiber, Zuschnitt: Planner.
- `in-progress` → `next`: Die Version muss beim Schreiben von einer dritten Schicht gesetzt werden (etwa vom `RecordService`) — zurück zur Zerlegung.
- `in-progress` → `open`: Der Nutzer ändert die Darstellung des Endes (`connection_closed` ohne Felder, letzte Nachricht) — Change Request an die Spezifikation.

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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Setzungen des Architect; die Darstellung und die Regel „Formatversion 2 nur bei Ende“ bestätigte der Nutzer am 2026-10-12). Ort, wenn nicht anders genannt: `LH-FA-25.a`, `SPEC-001`, `SPEC-002`, `SPEC-041`.

- *Version 2 ohne Ende:* der Leser nimmt sie an, der Schreiber schreibt sie nicht — `LH-FA-25.a` *Version*, `SPEC-001`.
- *`connection_closed` in Version 1:* `PGR-E3003` — `SPEC-001`.
- *Nicht als letzte Nachricht, nicht in der letzten Interaktion, zweimal in einer Session, mit Feldern:* `PGR-E3003` — `SPEC-002`, `SPEC-041`, `LH-FA-25.a` *Form*.
- *Nicht nach `FATAL`/`PANIC` (nach `ERROR`, ohne Fehlerantwort davor, Schweregrad nur in `S`):* `PGR-E3003` — `LH-FA-25.a` (Definition des Endes), `LH-FA-02.b` *Replay*.
- *`version: 3`, `0`, negativ, nicht ganzzahlig, fehlend:* `version: 3` und andere ganze Zahlen sind `PGR-E3002`; fehlend oder nicht kanonisch bleibt wie heute `PGR-E3003` — `SPEC-001`.
- *Extended-Interaktion mit Ende:* in der letzten Gruppe, ohne `ready_for_query`; die Gruppe endet mit `Sync` oder mit `Flush`; Gruppen davor wie bisher — `LH-FA-25.a` *Form*, `SPEC-041`.
- *Ende mit `offset_ms` und mit `empty_sessions`:* die Interaktion mit Ende trägt `offset_ms` wie jede andere; `empty_sessions` ist unabhängig vom Ende — `SPEC-002`, `LH-FA-12.a`.
- *Schweregrad und Schlüssel (Architect, 2026-10-15):* `V` entscheidet allein, wenn vorhanden, auch leer; `S` nur, wenn `V` fehlt; genauer Vergleich (`fatal` kein Ende); `V: ERROR` mit `S: FATAL` ist kein Ende; `connection_closed` mit irgendeinem Schlüssel neben `type`, auch leer oder `null`, ist `PGR-E3003`; eine weitere Nachricht (auch `notice_response`) zwischen Fehlerantwort und Ende ist `PGR-E3003` — `SPEC-002`.
- *Gruppen mit Ende:* letzte Gruppe mit `Sync` oder `Flush` und ohne `ready_for_query`, Gruppen davor mit `Flush`; ohne Ende bleibt die letzte Gruppe mit `Sync` und `ready_for_query` Pflicht (eine `Flush`-Gruppe als letzte ohne Ende ist `PGR-E3003`); `ready_for_query` in einer Gruppe mit Ende ist `PGR-E3003` — `SPEC-002`, `SPEC-041`.
- *Version im Modell gegen Inhalt:* das Feld `version` des Modells bestimmt nichts; der Schreiber richtet sich nach dem Inhalt (Feld 1 mit Ende schreibt 2, Feld 2 ohne Ende schreibt 1) — `SPEC-001`.
- *Reihenfolge mehrerer Ablehnungsgründe:* Lesbarkeit, Formatkennung, Version, Inhalt (`SPEC-001`); `version: 3` mit Fehlern im Inhalt, auch mit `connection_closed`, ist `PGR-E3002`; im Inhalt sind alle Gründe `PGR-E3003`, welcher gemeldet wird, ist nicht zugesagt (`SPEC-002`). Die Meldung eines Ablehnungsgrunds im Modell nennt die erste verletzte Prüfung in der Reihenfolge des Kommentars von `Group.validate` und `Interaction.Validate`; zugesagt ist nur der Fehler.
- *Aufruf-Voraussetzung:* `Validate`, die Formprüfung der Session und die Wahl der Version sind Zustandsprüfungen ohne Zusage über Reihenfolge oder Gleichzeitigkeit von Aufrufen; der Slice nennt keine Ordnung der Werte über mehrere Aufrufe, also gibt es keine Voraussetzung (`AGENTS.md` §3.12).
- *Andere Aufrufer von `Interaction.Validate`* (Replay, Einspielen, Record-Service): `Validate` nimmt ein Ende an jeder Interaktion an, die Regel „nur in der letzten Interaktion einer Session“ ist die der Session; diese Aufrufer behalten ihr Verhalten, bis ihre Slices (`slice-v1-abschluss-verbindungsende-replay`, `-einspielen`, `-record`) es ändern. Der SQLite-Leser und -Schreiber bleiben unberührt (`slice-v1-abschluss-verbindungsende-sqlite`): Er lehnt `version: 2` weiter als `PGR-E3002` ab, bis jener Slice liegt, und `model.Version` bleibt 1 für ihn (neue Konstante für die Version 2).

**Ort der Formregel einer Session (Entscheidung des Architect, 2026-10-15):** Methode `Session.ValidateEnd() error` im Modell (`internal/hexagon/model`), aufgerufen vom YAML-Leser nach den Interaktionen jeder Session; dazu `Recording.HasEnd() bool` für die Wahl der Version beim Schreiben. Grund: Die Regel braucht nur die Interaktionen der Session, also gehört sie an `Session`; `Recording` trägt keine Regel über Sessions hinweg. Das Modell bleibt Domäne, der Adapter importiert es schon (`adapters-driven → model`), es entsteht keine neue Kante und kein neues Paket, `make a-check` bleibt unberührt. Die Methode liegt im Kern, damit `slice-v1-abschluss-verbindungsende-record` und der SQLite-Slice dieselbe Regel aufrufen können statt sie nachzubauen.

**Zählung der Mutationstabelle (Architect, 2026-10-15; `AGENTS.md` §3.10):** ein Fall und eine Mutation, die nur ihn bricht, je Element von: den vier Verletzungen der Form (nicht letzte Nachricht · nicht letzte Interaktion · zweites Ende in einer Session · Felder), je für einfache Anfrage und Extended, wo die Prüfung eine andere Stelle des Codes ist; den Schritten der Ortsangabe der Extended-Form (letzte Server-Nachricht · letzte Gruppe · keine `ready_for_query` in der Gruppe); den Schweregraden (`FATAL` und `PANIC` angenommen; `ERROR`, keine Fehlerantwort, nur `S`, `V` gegen `S`, leeres `V`, `fatal`); „unmittelbar“ (eine Nachricht dazwischen); den Versionen (1 und 2 angenommen; 3, 0, negativ, zu groß: `PGR-E3002`); der Wahl der Version beim Schreiben (mit Ende 2, ohne 1, Feld gegen Inhalt in beide Richtungen); den Formen des Roundtrips (einfach · `Sync`-Gruppe · `Flush`-Gruppe · `offset_ms` · `empty_sessions`). Mutationen werden vom Implementer vor der Übergabe gefahren; der Testanteil am Diff kann dadurch auf das Mehrfache des Produktcodes steigen, das ist in der Schätzung in §8 gerechnet.

**Risiken:**

- Zwischenstand: Nach diesem Slice lädt `replay` eine Aufzeichnung mit Ende, der PGWire-Adapter kennt `connection_closed` aber noch nicht — **Ausgang:** offen bis Closure (erwartet: eingetreten, `slice-v1-abschluss-verbindungsende-replay`, direkt dahinter in der Reihenfolge nach `slice-v1-abschluss-verbindungsende-sqlite`).
- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 2×; `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`, 2×): Die Prüfung des Architect vor dem Code (2026-10-15) schätzt 850 bis 1250 (§8) und empfiehlt den Schnitt Leser | Schreiber; die Rückführung nach §4 ist eingetreten — **Ausgang:** offen bis Closure (erwartet: eingetreten, Zuschnitt und Kennungen beim Planner).
- Die Formregel einer Session („höchstens ein Ende, in der letzten Interaktion“) braucht einen Ort im Modell, den die Spezifikation nicht vorgibt (Methode an `Session` oder `Recording`, Aufruf aus den Lesern); die Wahl hat der Architect getroffen (oben, `Session.ValidateEnd`) — **Ausgang:** entfallen: entschieden vor dem Code, kein offener Punkt.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (2×) und `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×) — der Plan nennt seine Schätzung mit Tests und die Rückführung vorab (§4, §6); schätzt der Architect den Plan trotzdem über die Grenze, ist das das dritte Auftreten und beide Einträge sind Lücken mit eigenem Slice. `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` ist verkörpert (`AGENTS.md` §3.10): die Mutationstabelle steht in DoD und §7. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13): Liefer-Punkte: 3 (Modell mit Formregeln, YAML-Leser, YAML-Schreiber mit Roundtrip). Schichten: zwei (Kern mit dem Modell, `ARC-001`; Recording-Adapter, `ARC-008`); der Handbuch- und README-Teil liegt in `slice-v1-abschluss-verbindungsende-doku`. Schätzung des Planners: 650 bis 1000 Zeilen mit Tests.

**Schätzung des Architect (2026-10-15)**, in der Metrik von §4 (Summe aus Einfügungen und Löschungen, `internal` und `docs/user`, Tests eingeschlossen): Produktcode Modell etwa 70, YAML-Leser etwa 70, YAML-Schreiber etwa 40 (zusammen etwa 180); Tests Modell etwa 200, Leser etwa 250, Schreiber und Roundtrip etwa 230, Hilfen und Abdeckungstabellen etwa 50 (zusammen etwa 730). Das sind 850 bis 1000 vor der Nacharbeit der Mutationen; `slice-v1-abschluss-protokollrand` lag nach derselben Rechnung bei 800 bis 1000 und maß 1234 (Faktor 1,3 bis 1,5). Die realistische Spanne ist **850 bis 1250, Mitte etwa 1050, über der Marke 1000**. Empfehlung: Schnitt vor dem Anspruch nach §4, **Leser | Schreiber**. Erster Teil: Modell mit Formregeln (`Session.ValidateEnd`), YAML-Leser, Mutationen dazu (etwa 550 bis 750; das Laden einer Aufzeichnung mit Ende, was `-replay`, `-einspielen` und `-vergleich` voraussetzen). Zweiter Teil: YAML-Schreiber mit Version nach dem Inhalt und Roundtrip, Mutationen dazu (etwa 300 bis 500; Voraussetzung nur für `-record` und `-upstream`). Der erste ist einzeln lieferbar (handgeschriebene Aufzeichnungen), der zweite setzt ihn voraus. Der Planner schneidet, vergibt die Kennungen und stellt die Adressen um (`-replay` nennt als Start `…-format` in `done/`; `AGENTS.md` §3.13). Dies ist das dritte Auftreten von `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` im Sinne von Abschnitt „offene Beobachtungen sichten“: beide sind keine Notiz mehr, sondern Lücken, und brauchen einen eigenen Folge-Slice (Zuschnitt: Planner).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
