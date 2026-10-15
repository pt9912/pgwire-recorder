# Slice slice-v1-abschluss-verbindungsende-format: Verbindungsende im Format: Modell und YAML-Leser (Version 2)

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

**Ziel:** Das Modell kennt das Ende des Servers, und der YAML-Leser lädt es: `connection_closed` ohne Felder als letzte Server-Nachricht der letzten Interaktion einer Session, unmittelbar nach einer `error_response` `FATAL` oder `PANIC` (`LH-FA-25.a` *Form*, `SPEC-002`, `SPEC-041`). Der Leser nimmt `version: 1` und `version: 2` an, lehnt `connection_closed` in Version 1 und jede Stelle außerhalb der Form mit `PGR-E3003` ab und jede andere Version mit `PGR-E3002` (`SPEC-001`). Nachweis ohne `record` und ohne Schreiber: handgeschriebene Aufzeichnungen als Text im Test; am Binary ändert sich nur das Laden einer solchen Aufzeichnung.

**Herkunft:** Schnitt des Planners nach der Entscheidung des Nutzers vom 2026-10-12 (Option B, [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md)): Die Arbeit ist in acht Slices der Welle geschnitten, dieser ist der erste. Modell und YAML-Leser stehen vor allen Slices, die sie benutzen (SQLite, Replay, Einspielen, Record), weil jeder `connection_closed` und die Formregeln des Modells voraussetzt.

**Rückführung und Schnitt (2026-10-15):** Dieser Plan trug Modell, Leser und Schreiber mit Roundtrip. Die Prüfung des Architect vor dem Code schätzte ihn auf 850 bis 1250 in der Metrik von §4 (Mitte etwa 1050, über der Marke 1000); der Slice ging von `in-progress/` zurück nach `next/` und ist geteilt (Entscheidung des Planners nach der Empfehlung des Architect, Schnitt nach Lieferwert: Leser | Schreiber). Diese Kennung behält Modell, Formregeln und den YAML-Leser, weil Replay, Einspielen, SQLite und Record sie als Adresse für das Laden nennen und nichts anderes von ihr erwarten; der YAML-Schreiber mit Version nach dem Inhalt und der Roundtrip sind in [`slice-v1-abschluss-verbindungsende-format-schreiber`](../next/slice-v1-abschluss-verbindungsende-format-schreiber.md) aufgegangen (Schritt 35b). Beide Teile sind einzeln lieferbar: dieser mit handgeschriebenen Aufzeichnungen, jener setzt ihn voraus.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- SQLite (Schreiber, Leser, `meta.version` in der Transaktion der ersten Session mit Ende) — `slice-v1-abschluss-verbindungsende-sqlite` (dort §1 und DoD): zweites Format mit eigenen Randformen (Transaktion); mit ihm läge dieser Plan über der Grenze einer Review-Sitzung.
- Wiedergabe im Replay — `slice-v1-abschluss-verbindungsende-replay`: Nach diesem Slice lädt `replay` eine Aufzeichnung mit Ende, gibt `connection_closed` aber noch nicht wieder (wellenintern akzeptierter Zwischenstand, ein Release liegt nicht dazwischen).
- Einspielen und Vergleich — `slice-v1-abschluss-verbindungsende-einspielen` und `slice-v1-abschluss-verbindungsende-vergleich`: Sie lesen die Aufzeichnung über die Leser dieses Slice, verhalten sich aber anders als das Format.
- YAML-Schreiber (Version 2 genau dann, wenn die Aufzeichnung ein Ende enthält, sonst Version 1), `Recording.HasEnd()` und der Roundtrip (Laden, Schreiben, Laden) — `slice-v1-abschluss-verbindungsende-format-schreiber` (dort §1 und DoD): Der Roundtrip lädt über den Leser dieses Slice und kann erst nach ihm geprüft werden, und der Plan lag mit dem Schreiber über der Marke von §4 (Rückführung vom 2026-10-15). Nach diesem Slice schreibt das Produkt weiter nur `version: 1`; eine Version 2 entsteht von Hand.
- Aufzeichnen des Endes — `slice-v1-abschluss-verbindungsende-record` (Use Case, PGWire-Adapter) und `slice-v1-abschluss-verbindungsende-upstream` (Upstream-Adapter): Der `RecordService` setzt die Version nicht; sie folgt beim Schreiben dem Inhalt (`LH-FA-25.a` *Version*, Schreiber in `slice-v1-abschluss-verbindungsende-format-schreiber`).
- Änderung am neutralen Schema `tools/schema/schema.yaml` — Entscheidung des Nutzers vom 2026-10-12 (Darstellung bestätigt): Das Ende ist ein Wert von `message.kind`, keine neue Spalte.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (dort §1 und DoD), nach dem letzten Code-Slice des Endes in derselben Welle (`AGENTS.md` §3.11, §3.13): Mit der Dokumentation als dritter Schicht läge dieser Plan über zwei Schichten; ein Ende kann ein Anwender erst mit `record` erzeugen, und das Handbuch beschreibt, was das gebaute Binary zeigt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), `SPEC-002`, `SPEC-041`: `connection_closed` ist ein Antworttyp des Modells. `Validate` nimmt ihn nur als letzte Server-Nachricht an, unmittelbar nach einer `error_response` mit dem Schweregrad `FATAL` oder `PANIC` (Feld `V`, ohne es `S`) und ohne Felder: bei einer einfachen Anfrage statt `ready_for_query`, bei einer Extended-Interaktion in der letzten Gruppe, die dann mit `Sync` oder `Flush` endet und kein `ready_for_query` trägt. Eine Session trägt höchstens ein Ende, in ihrer letzten Interaktion: `Session.ValidateEnd() error` (Test).
- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), `SPEC-001`: Der YAML-Leser nimmt `version: 1` und `version: 2` an, auch Version 2 ohne Ende. `connection_closed` in Version 1, an einer anderen Stelle als der Form, nicht nach `FATAL`/`PANIC` und mit Feldern sind `PGR-E3003`; `version: 3` und jede andere ganze Zahl sind `PGR-E3002`. Der Leser ruft `Session.ValidateEnd` nach den Interaktionen jeder Session (Test mit handgeschriebenen Aufzeichnungen als Text; die vorhandenen Tests des Lesers bleiben grün).
- [x] Je Zusage des neuen Vertrags (Modell, Formregeln, Leser) eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/model` (`response.go`, `recording.go`, `extended.go`) | update | Antworttyp `connection_closed`, Konstante der Version 2, Formregeln einer Interaktion (darunter `ohneReadyForQuery`: kein `ready_for_query` in einer einfachen Interaktion mit Ende) und einer Session (`Session.ValidateEnd`) |
| `internal/hexagon/model` (Tests: `ende_test.go`, `extended_test.go`) | update | `Validate` und `Session.ValidateEnd`: Happy/Negative je Randform aus §6; die zwei Meldungstexte der Stellung von `sync` und `flush` in `extended_test.go` folgen der neuen Meldung |
| `internal/adapters/driven/recording/yaml.go` | update | Leser: beide Versionen (`kopfVorpruefung`), in der Vorprüfung am YAML-Baum `connection_closed` nicht in Version 1 und ohne Schlüssel neben `type` (auch leer oder `null`, was das Modell am dekodierten Wert nicht sieht), Aufruf von `Session.ValidateEnd` je Session; der Schreiber bleibt unberührt (`slice-v1-abschluss-verbindungsende-format-schreiber`) |
| `internal/adapters/driven/recording` (Tests: `yaml_ende_test.go`, drei vorhandene) | update | Leser; Fixtures als Text im Test; `yaml_test.go`, `leser_test.go` und `zeitangaben_pruefen_test.go` nehmen für „unbekannte Version“ die 3 statt der 2 (der SQLite-Zweig von `leser_test.go` behält die 2); `zeitangaben_pruefen_test.go` prüft die Meldung von `TestZeitangabenAbsteigendBeschaedigt` ohne Dateipfad (`ohnePfad`, Review F-4) |
| `test/integration/sqlite_lesen_e2e_test.go` | update | der Fall „yaml, Version 2“ von `TestE2ELesenAbgelehnt` nimmt die 3 (Version 2 ist gültig); der Fall „sqlite, Version 2“ bleibt |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-gesamt.md`, `docs/user/abdeckung-vollstaendig.md` | update (`make abdeckung`) | Tabellen aus den Deklarationen der Tests; Tests, die die alte Regel prüfen (kein Ende im Modell), gehören diesem Slice |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-protokollrand-doku` liegt in `done/` (Schritt 34); dieser Slice ist Schritt 35 der Reihenfolge in §5 von welle-v1-abschluss, der erste des Blocks zum Verbindungsende (Entscheidung des Nutzers vom 2026-10-12). WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat 6d0d697..<letzter Commit des Slice> -- internal docs/user`, Tests eingeschlossen, wobei die Basis `6d0d697` der letzte Commit vor dem Übergang nach `in-progress/` ist (der Commit dieser Zeile und der Move ändern nur Dateien unter `docs/plan/`, nicht die Pfade der Metrik); **über 1000** geht der Slice zurück zur Zerlegung. **Meldeschwelle 800:** Erreicht die Summe 800, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Architect für diesen Teil 550 bis 750 (Tests Modell etwa 200, Leser etwa 250, Produktcode etwa 140, Hilfen und Abdeckungstabellen etwa 50, vor der Nacharbeit der Mutationen; der Faktor 1,3 bis 1,5 aus `slice-v1-abschluss-protokollrand` ist darin gerechnet, der Abstand der Marke zur oberen Schätzung beträgt rund ein Drittel).
- `in-progress` → `next`: Der Leser verlangt eine Änderung am Schreiber, am Port oder an einem dritten Paket außerhalb von Modell und Leser — zurück zur Zerlegung.
- `in-progress` → `open`: Der Nutzer ändert die Darstellung des Endes (`connection_closed` ohne Felder, letzte Nachricht) — Change Request an die Spezifikation.

**Eingetreten (Grund der Rückführung `in-progress` → `next`, 2026-10-15):** Die erste Bedingung des Plans ist vor dem ersten Code-Commit eingetreten: Die Prüfung des Architect (Commit `690e353`) schätzte den ungeteilten Plan auf 850 bis 1250 in der Metrik oben, Mitte etwa 1050, über der Marke von 1000 (Schätzung des Planners 650 bis 1000). Die Zählung (drei Liefer-Punkte, zwei Schichten) hatte die Größe nicht gezeigt. Der Slice ist nach der Empfehlung des Architect geschnitten, Leser | Schreiber: dieser Plan behält Modell und Leser (550 bis 750), `slice-v1-abschluss-verbindungsende-format-schreiber` den Schreiber mit Roundtrip (300 bis 500). Gearbeitet war nichts außer dem Plan und der Prüfung des Architect.

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

- *Version 2 ohne Ende:* der Leser nimmt sie an (der Schreiber schreibt sie nicht, das ist Sache von `slice-v1-abschluss-verbindungsende-format-schreiber`) — `LH-FA-25.a` *Version*, `SPEC-001`.
- *`connection_closed` in Version 1:* `PGR-E3003` — `SPEC-001`.
- *Nicht als letzte Nachricht, nicht in der letzten Interaktion, zweimal in einer Session, mit Feldern:* `PGR-E3003` — `SPEC-002`, `SPEC-041`, `LH-FA-25.a` *Form*.
- *Nicht nach `FATAL`/`PANIC` (nach `ERROR`, ohne Fehlerantwort davor, Schweregrad nur in `S`):* `PGR-E3003` — `LH-FA-25.a` (Definition des Endes), `LH-FA-02.b` *Replay*.
- *`version: 3`, `0`, negativ, nicht ganzzahlig, fehlend:* `version: 3` und andere ganze Zahlen sind `PGR-E3002`; fehlend oder nicht kanonisch bleibt wie heute `PGR-E3003` — `SPEC-001`.
- *Extended-Interaktion mit Ende:* in der letzten Gruppe, ohne `ready_for_query`; die Gruppe endet mit `Sync` oder mit `Flush`; Gruppen davor wie bisher — `LH-FA-25.a` *Form*, `SPEC-041`.
- *Ende mit `offset_ms` und mit `empty_sessions`:* die Interaktion mit Ende trägt `offset_ms` wie jede andere; `empty_sessions` ist unabhängig vom Ende — `SPEC-002`, `LH-FA-12.a`.
- *Schweregrad und Schlüssel (Architect, 2026-10-15):* `V` entscheidet allein, wenn vorhanden, auch leer; `S` nur, wenn `V` fehlt; genauer Vergleich (`fatal` kein Ende); `V: ERROR` mit `S: FATAL` ist kein Ende; `connection_closed` mit irgendeinem Schlüssel neben `type`, auch leer oder `null`, ist `PGR-E3003`; eine weitere Nachricht (auch `notice_response`) zwischen Fehlerantwort und Ende ist `PGR-E3003` — `SPEC-002`.
- *`ready_for_query` vor dem Ende in einer einfachen Anfrage (Architect, 2026-10-15, Rückgabe aus dem Review F-1):* eine einfache Interaktion mit Ende trägt kein `ready_for_query`, an keiner Stelle; `[command_complete, ready_for_query, error_response FATAL, connection_closed]` ist `PGR-E3003` (das Ende steht statt `ready_for_query`, `LH-FA-25.a` *Form*, wie bei der Gruppe). Ohne Ende bleibt `ready_for_query` mitten in einer einfachen Interaktion unbeanstandet — akzeptiertes Negativ: die Regel gab es für Version 1 nie, sie lehnte bisher ladende Dateien ab und läge außerhalb von `LH-FA-25`; kein Widerspruch zum Lastenheft — `SPEC-002`.
- *Stelle des Endes in der Interaktion (Architect, 2026-10-15):* gültig nur als letzte Server-Nachricht der Interaktion, bei Extended der letzten Gruppe; jede andere Stelle (einfach, Extended, nicht letzte Gruppe) ist `PGR-E3003`. `Interaction.endet` und damit `Session.ValidateEnd` erkennen ein Ende an jeder Stelle, auch an falscher; ein falsch platziertes Ende in einer Interaktion vor der letzten ist `PGR-E3003`; welche Prüfung meldet, ist nicht zugesagt — `SPEC-002`.
- *Gruppen mit Ende:* letzte Gruppe mit `Sync` oder `Flush` und ohne `ready_for_query`, Gruppen davor mit `Flush`; ohne Ende bleibt die letzte Gruppe mit `Sync` und `ready_for_query` Pflicht (eine `Flush`-Gruppe als letzte ohne Ende ist `PGR-E3003`); `ready_for_query` in einer Gruppe mit Ende ist `PGR-E3003` — `SPEC-002`, `SPEC-041`.
- *Version im Modell gegen Inhalt:* der Leser setzt das Feld `version` des Modells aus der Datei (1 oder 2); die Wahl beim Schreiben nach dem Inhalt (Feld 1 mit Ende schreibt 2, Feld 2 ohne Ende schreibt 1) liegt in `slice-v1-abschluss-verbindungsende-format-schreiber` (dort §6) — `SPEC-001`.
- *Reihenfolge mehrerer Ablehnungsgründe:* Lesbarkeit, Formatkennung, Version, Inhalt (`SPEC-001`); `version: 3` mit Fehlern im Inhalt, auch mit `connection_closed`, ist `PGR-E3002`; im Inhalt sind alle Gründe `PGR-E3003`, welcher gemeldet wird, ist nicht zugesagt (`SPEC-002`). Die Meldung eines Ablehnungsgrunds im Modell nennt die erste verletzte Prüfung in der Reihenfolge des Kommentars von `Group.validate` und `Interaction.Validate`; zugesagt ist nur der Fehler.
- *Aufruf-Voraussetzung:* `Validate`, die Formprüfung der Session und die Wahl der Version sind Zustandsprüfungen ohne Zusage über Reihenfolge oder Gleichzeitigkeit von Aufrufen; der Slice nennt keine Ordnung der Werte über mehrere Aufrufe, also gibt es keine Voraussetzung (`AGENTS.md` §3.12).
- *Andere Aufrufer von `Interaction.Validate`* (Replay, Einspielen, Record-Service): `Validate` nimmt ein Ende an jeder Interaktion an, die Regel „nur in der letzten Interaktion einer Session“ ist die der Session; diese Aufrufer behalten ihr Verhalten, bis ihre Slices (`slice-v1-abschluss-verbindungsende-replay`, `-einspielen`, `-record`) es ändern. Der SQLite-Leser und -Schreiber bleiben unberührt (`slice-v1-abschluss-verbindungsende-sqlite`): Er lehnt `version: 2` weiter als `PGR-E3002` ab, bis jener Slice liegt, und `model.Version` bleibt 1 für ihn (neue Konstante für die Version 2).

**Ort der Formregel einer Session (Entscheidung des Architect, 2026-10-15):** Methode `Session.ValidateEnd() error` im Modell (`internal/hexagon/model`), aufgerufen vom YAML-Leser nach den Interaktionen jeder Session; dazu `Recording.HasEnd() bool` für die Wahl der Version beim Schreiben (geliefert von `slice-v1-abschluss-verbindungsende-format-schreiber`, das den Aufrufer hat). Grund: Die Regel braucht nur die Interaktionen der Session, also gehört sie an `Session`; `Recording` trägt keine Regel über Sessions hinweg. Das Modell bleibt Domäne, der Adapter importiert es schon (`adapters-driven → model`), es entsteht keine neue Kante und kein neues Paket, `make a-check` bleibt unberührt. Die Methode liegt im Kern, damit `slice-v1-abschluss-verbindungsende-record` und der SQLite-Slice dieselbe Regel aufrufen können statt sie nachzubauen.

**Zählung der Mutationstabelle (Architect, 2026-10-15; `AGENTS.md` §3.10):** ein Fall und eine Mutation, die nur ihn bricht, je Element von: den vier Verletzungen der Form (nicht letzte Nachricht · nicht letzte Interaktion · zweites Ende in einer Session · Felder), je für einfache Anfrage und Extended, wo die Prüfung eine andere Stelle des Codes ist; den Schritten der Ortsangabe der Extended-Form (letzte Server-Nachricht · letzte Gruppe · keine `ready_for_query` in der Gruppe); den Schweregraden (`FATAL` und `PANIC` angenommen; `ERROR`, keine Fehlerantwort, nur `S`, `V` gegen `S`, leeres `V`, `fatal`); „unmittelbar“ (eine Nachricht dazwischen); den Versionen (1 und 2 angenommen; 3, 0, negativ, zu groß: `PGR-E3002`). Die Wahl der Version beim Schreiben und die Formen des Roundtrips (einfach · `Sync`-Gruppe · `Flush`-Gruppe · `offset_ms` · `empty_sessions`) zählen im Plan von `slice-v1-abschluss-verbindungsende-format-schreiber`. Mutationen werden vom Implementer vor der Übergabe gefahren; der Testanteil am Diff kann dadurch auf das Mehrfache des Produktcodes steigen, das ist in der Schätzung in §8 gerechnet.

**Risiken:**

- Zwischenstand: Nach diesem Slice lädt `replay` eine Aufzeichnung mit Ende, der PGWire-Adapter kennt `connection_closed` aber noch nicht — **Ausgang:** eingetreten: `slice-v1-abschluss-verbindungsende-replay` (Schritt 37, gibt `connection_closed` wieder; nennt in §1 „Übernommen aus `slice-v1-abschluss-verbindungsende-format`“ und in §4 diesen Slice als Start).
- Der ungeteilte Plan lag über einer Review-Sitzung (Prüfung des Architect vom 2026-10-15: 850 bis 1250, §8; `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`) — **Ausgang:** eingetreten: `slice-v1-abschluss-verbindungsende-format-schreiber` nimmt den Schreiber mit Roundtrip auf (Schnitt vom 2026-10-15, §1 und §4).
- Dieser Teil überschreitet trotz des Schnitts Meldeschwelle oder Marke (Schätzung 550 bis 750; `slice-v1-abschluss-protokollrand` maß 23 % über der Obergrenze, das wären etwa 920; `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, Tests mit Tabelle von Randformen und Mutationen je Element) — **Ausgang:** entfallen: gemessen 655 in der Metrik von §4 (Meldeschwelle 800 und Marke 1000 nicht erreicht, Schätzung 550 bis 750 getroffen); `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` trat nicht auf (Vermerk im Register).
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

### Belege des Implementers

**Größe** (Metrik von §4, `git diff --numstat 6d0d697..<letzter Commit> -- internal docs/user`): erste Fassung (Stand des Commits „Lint-Befunde behoben“) 590; nach der Nacharbeit aus dem Review (Stand des Arbeitsbaums vor dem Commit der Nacharbeit) 627 Einfügungen und 28 Löschungen, Summe **655**. Meldeschwelle 800 und Marke 1000 nicht erreicht; Schätzung 550 bis 750. Außerhalb der Metrik: 2 Zeilen in `test/integration/sqlite_lesen_e2e_test.go`.

**Läufe.** `make gates` rot (Exit 2) beim ersten Lauf: `TestE2ELesenAbgelehnt` (Integrationstest) nahm „yaml, Version 2“ als unbekannte Version; behoben (Fall nimmt die 3). Zweiter Lauf rot: `make lint` (`cyclop`/`gocyclo`/`gocognit` für `Group.validate` und `interaktionVorpruefung`, `gochecknoglobals` für die Test-Variable `fatal`); behoben durch `validateAbschluss`, `antwortVorpruefung`, `fatalFelder()`, danach `make lint` allein grün (0 issues). Der letzte volle Lauf steht unten (Eintrag *Letzter Lauf*). Zwischendurch lief die Suite per Bind-Mount (`go test ./internal/...`, ohne Netz) grün. Nacharbeit aus dem Review (F-1 bis F-4, 2026-10-15): Bind-Mount-Läufe mit dem Image `pgwire-recorder:test` (`--network=none`, Quelle als Volume), Mutanten in einer frischen Kopie des Baums, je Lauf `./internal/hexagon/...` und `./internal/adapters/driven/recording/ -run 'Unmarshal|Ende'`.

**Mutationstabelle** (`AGENTS.md` §3.10). Jede Mutation einzeln angewandt, `go test ./internal/...` (alle Unit-Tests, Bind-Mount, `--network=none`; der Mutant kommt im Lauf an, weil die Datei neu geschrieben wird, danach zurückgesetzt), die roten Tests sind die der ganzen Suite. „Genannte Tests“ ist die Liste in der Spalte; die Namen kürzen die Unterfälle (`TestUnmarshalEnde/…`) auf den Test. Ein Mutant, der die Zusage bricht und bei dem die Spalte nur Tests der Zusage nennt, ist „nur diese“.

| Zusage | Mutation | roter Test (Suite `./internal/...`) |
|---|---|---|
| Schweregrad `FATAL` angenommen | `beendetVerbindung` ohne `FATAL` | `TestValidateEndeGueltig`, `TestUnmarshalEnde` (alle FATAL-Fälle), `TestUnmarshalEndeAbgelehnt/nicht_letzte_Interaktion…`, `TestUnmarshalEndeJeSession` (benutzen FATAL) |
| `PANIC` angenommen | ohne `PANIC` | nur `TestValidateEndeGueltig`, `TestUnmarshalEnde/einfach,_PANIC` und `…PANIC_nur_in_S` |
| `S` zählt, wenn `V` fehlt | Rückfall auf `S` abgeschaltet | nur `TestValidateEndeGueltig`, `TestUnmarshalEnde/einfach,_FATAL_nur_in_S` und `…PANIC_nur_in_S` |
| `V` entscheidet vor `S` | immer `S` gelesen | `TestValidateEndeFehler/query_V_ERROR_gegen_S_FATAL`, `…leeres_V_gegen_S_FATAL`, `TestValidateEndeGueltig` (V FATAL gegen S ERROR), `TestUnmarshalEndeAbgelehnt/V_ERROR_gegen_S_FATAL`, `…leeres_V…`, `TestUnmarshalEnde/einfach,_PANIC`, `…V_FATAL_gegen_S_ERROR` |
| leeres `V` entscheidet allein | Rückfall auch bei leerem `V` | nur `TestValidateEndeFehler/query_leeres_V_gegen_S_FATAL`, `TestUnmarshalEndeAbgelehnt/leeres_V_gegen_S_FATAL` |
| genauer Vergleich (`fatal`, `panic`) | `strings.ToUpper` | nur `TestValidateEndeFehler/query_fatal_klein_in_V`, `…_in_S`, `…panic_klein_in_V`, `TestUnmarshalEndeAbgelehnt/fatal_klein` |
| genauer Vergleich (Leerraum) | `strings.TrimSpace` | nur `TestValidateEndeFehler/query_FATAL_mit_Leerzeichen` |
| nur eine `error_response` zählt | Typprüfung entfernt | nur `TestValidateEndeFehler/query_Hinweis_mit_FATAL_statt_Fehlerantwort` |
| unmittelbar (nichts dazwischen) | irgendeine FATAL-Antwort davor genügt | nur `TestValidateEndeFehler/query_Hinweis_zwischen_…`, `…extended_Hinweis_zwischen_…`, `TestUnmarshalEndeAbgelehnt/Hinweis_zwischen_Fehler_und_Ende` |
| eine Fehlerantwort steht davor (Anfang) | `idx == 0` angenommen | nur `TestValidateEndeFehler/query_nur_connection_closed`, `TestUnmarshalEndeAbgelehnt/nur_Ende` |
| eine Fehlerantwort FATAL/PANIC steht davor (nicht nur nicht am Anfang) | zweiter Operand entfernt | `TestValidateEndeFehler` (ERROR, ohne Fehlerantwort, Schweregradfälle, Hinweis), `TestUnmarshalEndeAbgelehnt` (dieselben Fälle) |
| ohne Felder: Spalten · Werte · Tag · Felder · Name · Value · TxStatus · Parametertypen | je Feld die Klausel in `ohneFelder` durch `true` ersetzt (acht Mutationen) | je Feld nur der eigene Fall `TestValidateEndeFehler/query_mit_<Feld>` (bei `Tag` zusätzlich `extended_mit_Feld`, der `Tag` benutzt) |
| ohne Felder (alle) | `ohneFelder`-Prüfung abgeschaltet | die acht Feldfälle und `extended_mit_Feld` |
| Schlüssel neben `type` in der Datei: mit Wert · leer · `null` | Prüfung `len(r.Content) != 2` entfernt, und drei Mutanten, die je eine Art ausblenden | `TestUnmarshalEndeAbgelehnt/Schlüssel_…`: ohne Prüfung alle fünf; ohne die Art „Wert“ nur `Schlüssel_mit_Wert` und `Schlüssel_unbekannt`; ohne „leer“ nur `Schlüssel_leer` und `Schlüssel_leer,_extended`; ohne „null“ nur `Schlüssel_null` |
| nur als letzte Nachricht (einfach und Extended) | Stellenprüfung in `pruefeEnde` entfernt | `TestValidateEndeFehler/query_nicht_letzte_Nachricht…`, `extended_nicht_letzte_Nachricht`, `extended_mit_ready_for_query_nach_dem_Ende`, `TestUnmarshalEndeAbgelehnt/nicht_letzte_Nachricht,_einfach` und `…_extended`; einfach und Extended teilen die Stelle im Code, eine Mutation färbt beide |
| einfache Anfrage: Ende statt `ready_for_query` | Endetyp als letzte Antwort nicht angenommen | `TestValidateEndeGueltig` und die Ende-Fälle von `TestUnmarshalEnde` (alle einfachen Fälle) |
| Typ bekannt in einfacher Anfrage / in Extended | `connection_closed` aus `simpleResponses` / `extendedResponses` | die Gutfälle und die Fälle zu `connection_closed` der einfachen Anfrage bzw. der Extended-Interaktion |
| einfache Anfrage: Regeln des Endes angewandt | `pruefeEnde` in `validateQuery` entfernt | die einfachen Fälle von `TestValidateEndeFehler` und `TestUnmarshalEndeAbgelehnt` |
| Extended: letzte Gruppe mit Ende, Ende erkannt | `ende := false` | `TestValidateEndeGueltig`, `TestUnmarshalEnde/extended,_*`, `TestUnmarshalEndeAbgelehnt/nicht_letzte_Interaktion,_extended` |
| Extended: letzte Gruppe mit Ende endet mit `flush` | `flush` nur ohne Ende abgelehnt, mit Ende wie ohne | nur `TestValidateEndeGueltig` (extended mit flush), `TestUnmarshalEnde/extended,_flush` |
| Extended: letzte Gruppe mit Ende endet mit `sync` | `sync` bei Ende abgelehnt | `TestUnmarshalEnde/extended,_sync`, `…vordere_Flush-Gruppe`, die Extended-Fälle mit `sync` |
| Extended: Gruppen davor enden mit `flush` | Prüfung `sync && !letzte` entfernt | `TestValidateEndeFehler/extended_vordere_Gruppe_mit_sync…`, `TestUnmarshalEndeAbgelehnt/vordere_Gruppe_mit_sync_vor_Ende`, `TestValidateFehler/vordere_Gruppe_endet_mit_sync`, `TestValidateFehlerReihenfolge/sync_am_Ende…` |
| Extended: ohne Ende endet die letzte Gruppe mit `sync` | Prüfung für die letzte Gruppe mit `ende` vertauscht | `TestValidateEndeFehler/extended_ohne_Ende,_letzte_Gruppe_mit_flush`, `TestUnmarshalEndeAbgelehnt/letzte_Gruppe_mit_flush_ohne_Ende`, `TestValidateFehler/letzte_Gruppe_endet_mit_flush…`, `TestUnmarshalExtendedFehler/abgeschnitten_nach_der_Flush-Gruppe`, `TestRecordExtendedNichtDarstellbar/ReadyForQuery_ohne_Sync` |
| Extended: Ende nur in der letzten Gruppe | Prüfung `!letzte && g.endet()` entfernt | nur `TestValidateEndeFehler/extended_Ende_in_einer_vorderen_Gruppe`, `…in_der_vorderen_Gruppe,_letzte_Gruppe_regulär`, `TestUnmarshalEndeAbgelehnt/Ende_in_einer_vorderen_Gruppe` |
| Extended: Regeln des Endes je Gruppe angewandt | `pruefeEnde(g.Server)` entfernt | die Extended-Fälle von `TestValidateEndeFehler` und `TestUnmarshalEndeAbgelehnt` |
| Extended: ohne Ende bleibt `ready_for_query` am Ende Pflicht | Prüfung entfernt | `TestValidateEndeFehler/extended_ohne_Ende_und_ohne_ready_for_query`, `TestValidateFehler/letzte_Gruppe_ohne_…`, `TestUnmarshalExtendedFehler/abgeschnitten_in_den_Server-Nachrichten` |
| Extended: mit Ende kein `ready_for_query` verlangt | Pflicht auch mit Ende | `TestValidateEndeGueltig`, `TestUnmarshalEnde/extended,_*` |
| Extended: kein `ready_for_query` in der Gruppe mit Ende | Klausel `si != len-1` in `validateServer` entfernt | `TestValidateEndeFehler/extended_mit_ready_for_query_in_der_Gruppe_mit_Ende`, `TestUnmarshalEndeAbgelehnt/ready_for_query_in_der_Gruppe_mit_Ende`, dazu die vorhandenen `TestValidateFehler/ready_for_query_vor_dem_Ende…` und `TestValidateFehlerReihenfolge/…` (die Regel ist die vorhandene, ein eigener Code dafür fehlt; „nur diese“ gilt nicht) |
| Session: Ende in der letzten Interaktion erlaubt | `ii != letzte` entfernt | `TestSessionValidateEnd/Ende_in_der_einzigen_Interaktion`, `…_letzten_von_zwei`, `Extended-Ende_in_der_letzten_von_drei`, die Ende-Fälle von `TestUnmarshalEnde` |
| Session: Interaktion ohne Ende nicht beanstandet | `i.endet()` ersetzt durch `true` | jeder Fall mit einer Interaktion ohne Ende, in der Suite breit (u. a. `TestExtendedRoundtrip`, `TestLesen…`); die Zusage „nur die Interaktion mit Ende zählt“ trägt `TestSessionValidateEnd/ohne_Ende` |
| Session: nicht letzte Interaktion, einfach | Responses-Zweig von `Interaction.endet` entfernt | `TestSessionValidateEnd/einfaches_Ende_in_der_ersten_von_zwei`, `…_mittleren_von_dreien`, `zwei_Enden_in_aufeinanderfolgenden_Interaktionen`, `TestUnmarshalEndeAbgelehnt/nicht_letzte_Interaktion,_einfach`, `zwei_Enden_in_einer_Session` |
| Session: nicht letzte Interaktion, Extended | Groups-Zweig entfernt | `TestSessionValidateEnd/Extended-Ende_in_der_ersten_von_zwei`, `…_mittleren_von_dreien`, `zwei_Enden,_vorletztes_Extended`, `TestUnmarshalEndeAbgelehnt/nicht_letzte_Interaktion,_extended` |
| Session: zweites Ende | `ii < letzte-1` statt `ii != letzte` | die Fälle mit Ende in der ersten und mittleren Interaktion, dazu `zwei_Enden…`; eine Mutation, die nur das zweite Ende bricht, gibt es nicht (siehe Grenzen) |
| Leser ruft `Session.ValidateEnd` | Aufruf entfernt | nur `TestUnmarshalEndeAbgelehnt/nicht_letzte_Interaktion,_einfach`, `…_extended`, `zwei_Enden_in_einer_Session` |
| `VersionMitEnde` ist 2, `Version` bleibt 1 | `VersionMitEnde = 3` / `Version = 2` | `TestVersionMitEnde` und die Ende-Fälle des Lesers / `TestVersionMitEnde` und die SQLite-Tests (`TestLesen…`, `TestRecordZeitRoundtrip`, …) |
| Version 2 angenommen | nur `model.Version` angenommen | `TestUnmarshalEnde` (alle Version-2-Fälle), `TestUnmarshalVersionen/Version_2` |
| Version 3 / 0 / negativ / über 64 Bit: `PGR-E3002` | je ein Mutant, der genau diese Zahl annimmt (vier Mutationen) | `TestUnmarshalVersionen/Version_3` (+`TestUnmarshalFehler/unbekannte_Version`, `TestZeitangabenVersionVorForm/Version_3…`, `TestLesenGleicheAblehnungBeiderFormate`) / `…Version_0` (+`TestZeitangabenVersionVorForm/Version_0`) / `…Version_-1` (+`TestZeitangabenVersionVorForm/Version_-1`) / nur `TestUnmarshalVersionen/Version_99999999999999999999` |
| Version 2 ohne Ende angenommen | Leser lehnt Version 2 ohne Ende ab | nur `TestUnmarshalEnde/Version_2_ohne_Ende`, `TestUnmarshalVersionen/Version_2` |
| `connection_closed` in Version 1 ist `PGR-E3003` | Prüfung entfernt / auf Version 2 vertauscht | nur `TestUnmarshalEndeAbgelehnt/Version_1,_einfach` und `…_extended` / die Ende-Fälle von `TestUnmarshalEnde` und beide Fälle zu Version 1 |
| Extended-Antworten der Gruppen werden vorgeprüft | Anhängen der Gruppen-Antworten entfernt | `TestUnmarshalEndeAbgelehnt/Schlüssel_leer,_extended`, `Version_1,_extended`, dazu die vorhandenen `param_types`-Fälle |
| Das Feld `version` des Modells trägt die Version der Datei | `Version: model.Version` im Leser | `TestUnmarshalEnde` (Version-2-Fälle) |
| Ende mit `offset_ms` / mit `empty_sessions` geladen | Vorprüfung lehnt `offset_ms` neben dem Ende ab / `emptySessions` mit Ende | nur `TestUnmarshalEnde/mit_offset_ms` / nur `TestUnmarshalEnde/mit_empty_sessions` |
| Version geht dem Inhalt voran (auch bei Schlüssel neben `type`) | `formVorpruefung` vor `kopfVorpruefung` | `TestUnmarshalVersionen/Version_3_mit_Schlüssel_neben_type`, dazu die vorhandenen `TestZeitangabenVersionVorForm/…` |
| einfach mit Ende: kein `ready_for_query` am Anfang (F-1 a) | `ohneReadyForQuery` überspringt Index 0 | nur `TestValidateEndeFehler/query_ready_for_query_am_Anfang,…`, `TestUnmarshalEndeAbgelehnt/ready_for_query_am_Anfang,…` |
| einfach mit Ende: kein `ready_for_query` in der Mitte | überspringt die Indizes zwischen Anfang und der Stelle drei vor dem Ende | nur `…ready_for_query_in_der_Mitte,…` (Validate und Leser) |
| einfach mit Ende: kein `ready_for_query` direkt vor der Fehlerantwort | überspringt den Index drei vor dem Ende | nur `…ready_for_query_direkt_vor_der_Fehlerantwort,…` (Validate und Leser) |
| einfach mit Ende: kein `ready_for_query` zwischen Fehlerantwort und Ende | `pruefeEnde` lässt ein `ready_for_query` vor `connection_closed` als Vorgänger durch (M4a), dazu `ohneReadyForQuery` überspringt den Index zwei vor dem Ende (M4) | M4a und M4: nur `…ready_for_query_zwischen_Fehlerantwort_und_Ende` (Validate und Leser); M4b, die Überspringung allein, bleibt grün (siehe Grenzen) |
| ohne Ende bleibt `ready_for_query` mitten in der einfachen Interaktion unbeanstandet (Gegenfall) | ohne Ende wird jedes `ready_for_query` außer dem letzten abgelehnt | `TestValidateEndeGueltig` (Fall „query ohne Ende, ready_for_query mitten“), `TestUnmarshalEnde/Version_2_ohne_Ende,_ready_for_query_mitten…`, dazu `TestVergleichStruktur/Aufzeichnung_länger,_Anfang_gleich/extended=false` (Services, Bestand ohne Ende) |
| Ende mitten in der Interaktion (F-1 b): einfach, Extended letzte Gruppe | Stellenprüfung in `pruefeEnde` entfernt | zusätzlich zu den Fällen oben `TestValidateEndeFehler/query_Ende_in_der_Mitte,_Nachricht_ohne_ready_for_query_danach`, `…extended_Ende_mitten_in_der_letzten_Gruppe`, `TestUnmarshalEndeAbgelehnt/Ende_mitten_in_der_letzten_Interaktion,_einfach` und `…letzten_Gruppe,_extended` |
| Ende in einer nicht letzten Gruppe | Prüfung `!letzte && g.endet()` abgeschaltet | nur `TestValidateEndeFehler/extended_Ende_in_einer_vorderen_Gruppe`, `…in_der_vorderen_Gruppe,_letzte_Gruppe_regulär`, `TestUnmarshalEndeAbgelehnt/Ende_in_einer_vorderen_Gruppe` |
| `ValidateEnd`: falsch platziertes Ende in einer Interaktion vor der letzten, einfach (Mutation 1) | Responses-Zweig von `Interaction.endet` erkennt nur die letzte Nachricht | nur `TestSessionValidateEnd/falsch_platziertes_Ende_in_der_ersten_von_zwei,_einfach` und `…am_Anfang_der_ersten_von_zwei,_einfach` |
| `ValidateEnd`: falsch platziertes Ende, Extended (Mutation 2) | `Group.endet` erkennt nur die letzte Server-Nachricht | nur `TestSessionValidateEnd/falsch_platziertes_Ende_in_der_letzten_Gruppe_der_ersten_von_zwei,_Extended` und `…in_einer_vorderen_Gruppe_der_ersten_von_zwei,_Extended` |
| Stelle des Endes in `pruefeEnde` (Mutation 3, mit den Fällen aus F-1 b) | Stellenprüfung entfernt | sechs Fälle von `TestValidateEndeFehler` und vier von `TestUnmarshalEndeAbgelehnt` (Liste: die „nicht letzte Nachricht“-Fälle, „Ende mitten …“, „ready_for_query nach dem Ende“) |
| F-4: Meldung nennt keinen Wert (Test pfadunabhängig) | `offset_ms kleiner …` hängt den Wert des kleineren Eintrags an / den vorigen Wert | `TestZeitangabenAbsteigendBeschaedigt` (12 Unterfälle je Mutant, ohne Pfad geprüft) |

**Grüne Mutanten und Grenzen, ehrlich.**

- *Zweites Ende in einer Session* hat keine eigene Mutation: Zwei Enden heißt, dass mindestens eines nicht in der letzten Interaktion steht; die Positionsprüfung in `ValidateEnd` trägt beides, und ein getrennter Zähler wäre toter Code (jeder Mutant daran wäre äquivalent). Die Fälle `zwei_Enden_in_aufeinanderfolgenden_Interaktionen` und `zwei_Enden,_vorletztes_Extended` stehen im Test, die Mutation `ii < letzte-1` färbt sie zusammen mit den Fällen für „nicht in der letzten Interaktion“.
- *`ready_for_query` in der Gruppe mit Ende* hat keinen eigenen Code; die vorhandene Klausel in `validateServer` trägt es, die Mutation daran färbt die neuen und vorhandene Tests.
- *Einfache und Extended teilen den Code* der Stelle, der Felder und des Schweregrads (`pruefeEnde`), eine Mutation färbt beide Arten; die Zählung des Architect „je einfach und Extended, wo die Prüfung eine andere Stelle ist“ trifft nur auf die Typlisten, den Abschluss der Gruppe und die Einbindung zu, und jede hat ihre Zeile.
- *Flake `TestZeitangabenAbsteigendBeschaedigt/…sqlite` (Review F-4):* Ursache ist der Pfad, nicht Last. Die Meldung des SQLite-Lesers nennt den Dateipfad, `t.TempDir()` trägt eine Zufallszahl, und der Test lehnte jede Meldung mit einem der Werte (`4709`, `4710`, `4300`, `4900`, `4800`) als Teilstring ab. Reproduziert: Der Test der Basis (Stand `HEAD` vor der Nacharbeit) mit `TMPDIR=/tmp/t4709x` ist rot (`Meldung nennt den Wert 4709: … /tmp/t4709x/TestZeitangabenAbsteigendBeschaedigt…/rec.sqlite: …`); derselbe Aufruf mit dem geänderten Test (`ohnePfad` streicht jeden Pfad unter `os.TempDir()` aus der Meldung) ist grün. Dass der Test den Wert weiter fängt, zeigen die Mutanten in der Tabelle (Wert des kleineren und des vorigen Eintrags in der Meldung, je rot).
- *`ready_for_query` zwischen Fehlerantwort und Ende* ist doppelt gedeckt: `pruefeEnde` („unmittelbar“) lehnt es ab, bevor `ohneReadyForQuery` läuft. Die Überspringung dieser Stelle in `ohneReadyForQuery` allein (M4b) bleibt grün und ist äquivalent, weil die Meldung ohnehin von `pruefeEnde` kommt. Die Zusage bricht erst M4 (beide Stellen zugleich: `pruefeEnde` lässt ein `ready_for_query` als Vorgänger durch und `ohneReadyForQuery` überspringt die Stelle); M4 färbt nur den Fall zur Stelle rot (Validate und Leser), die Verifikation fuhr ihn nach. M4a allein (nur `pruefeEnde` gelockert) färbt den Fall ebenfalls, aber nur über den Meldungstext: `ohneReadyForQuery` meldet weiter einen Fehler, die Zusage bleibt gewahrt (die Meldung ist nicht zugesagt, `SPEC-002`; Verifikation V-2).
- *Mutation 1 und 2 färben nur Tests des Modells:* Im Leser fängt `Interaction.Validate` ein falsch platziertes Ende, bevor `Session.ValidateEnd` läuft; die Leser-Fälle „falsch platziertes Ende in einer Interaktion vor der letzten“ prüfen nur den Fehler (`PGR-E3003`), und die Meldung ist nicht zugesagt (`SPEC-002`). `Session.ValidateEnd` prüft für sich `TestSessionValidateEnd`.
- *Nicht mutiert:* die Meldungstexte der neuen Ablehnungen (die Tests prüfen Teilstrings der Regel, nicht den ganzen Text); die Blockform `type: connection_closed` im Leser (Fall `Ende in Blockform`, kein eigener Mutant); die Längenprüfung `len(g.Server) > 0` bei der Erkennung des Endes in `Group.validate`.
- *SQLite:* Eine SQLite-Datei mit der Art `connection_closed` geht über `Unmarshal` mit Version 1 und ist `PGR-E3003`; ein Test dafür steht in diesem Slice nicht (die Regel des SQLite-Lesers gehört `slice-v1-abschluss-verbindungsende-sqlite`). `version2` in SQLite bleibt `PGR-E3002` (`TestE2ELesenAbgelehnt`, `TestLesenGleicheAblehnungBeiderFormate`).
- *Allgemein:* `make abdeckung` hat `LH-FA-25` in `abdeckung-gesamt.md` als „vollständig“ (Unit) eingetragen, obwohl Record, Replay, Einspielen und SQLite fehlen; die Tabelle ist Ausgabe des Werkzeugs und zählt die Deklarationen, nicht den Umfang der Anforderung. Für das Review.
- *Randformen:* Die zwei Lesarten der ersten Fassung (`ready_for_query` mitten in einer einfachen Anfrage vor dem Ende; `Interaction.endet` zählt `connection_closed` an jeder Stelle) hatte der Code entschieden, ohne dass §6 sie nannte (Review F-1, `AGENTS.md` §3.12). Der Architect hat beide entschieden und in `SPEC-002` und §6 eingetragen (Commit `aed8a49`): Eintrag *`ready_for_query` vor dem Ende in einer einfachen Anfrage* (kein `ready_for_query` an keiner Stelle bei einem Ende; ohne Ende unbeanstandet, akzeptiertes Negativ) und Eintrag *Stelle des Endes in der Interaktion* (nur letzte Server-Nachricht, bei Extended der letzten Gruppe; `Interaction.endet` und `ValidateEnd` erkennen ein Ende an jeder Stelle; zugesagt ist nur der Fehler). Die Nacharbeit setzt beide um, ohne `endet` oder `ValidateEnd` zu ändern: Die Fälle verlangten es nicht.
- *Nacharbeit als Klasse (F-1):* Merkmal ist die Stelle einer Nachricht in der Nachrichtenfolge. Für `ready_for_query` vor dem Ende sind die Ausprägungen Anfang · Mitte · direkt vor der Fehlerantwort · zwischen Fehlerantwort und Ende je ein Fall (Validate und Leser) mit eigener Mutation; für `connection_closed` mitten in der Folge: einfach, Extended in der letzten Gruppe, Extended in einer nicht letzten Gruppe, und in `ValidateEnd` (Interaktion vor der letzten) einfach, Extended in der letzten und in einer vorderen Gruppe. Offen gelassen: Ausprägungen, die der Architect nicht nannte (zum Beispiel `ready_for_query` nach dem Ende in einer einfachen Anfrage: gedeckt durch den bestehenden Fall „query nicht letzte Nachricht“, der vor `ohneReadyForQuery` anschlägt).
- *Review F-3:* nicht geändert; die Abdeckungstabelle ist Ausgabe des Werkzeugs, die Lesart von „vollständig“ bei teilweise gelieferter Anforderung liegt beim Planner. *F-5:* keine Aktion.

**Letzter Lauf.** `make gates` auf dem Arbeitsbaum der Nacharbeit aus dem Review (2026-10-15, vor ihrem Commit): Exit 0, alle Gates grün, darunter `lint`, `a-check`, `test`, `test-integration`, `test-postgres-versionen` und die Gegenproben. Davor ein roter Lauf der Nacharbeit: `docs-check` meldete eine Kennung ohne Link in der Eingangs-Kontext-Zeile des Review-Reports (`id-unlinked`); die Zeile trägt jetzt den Link, sonst nichts am Report. Keine Datei außerhalb von `docs/plan/` hat sich seit dem grünen Lauf geändert.

### Closure (Planner, 2026-10-15)

**Review und Verifikation.** Review `3d14ef7` (0 HIGH, 1 MEDIUM F-1, 1 LOW F-2, 3 INFO F-3 bis F-5); F-1 durch den Architect (`aed8a49`, `SPEC-002` und Plan §6) und die Nacharbeit `91a5228` behoben, F-2 und F-4 behoben, F-3 siehe Register, F-5 keine Aktion. Verifikation `357b872`: **PASS mit Befunden**, 0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO (V-1 bis V-4); alle DoD-Punkte bestätigt, 18 Mutanten nachgefahren, jeder färbt die in §7 genannten Tests (M4b grün, äquivalent), die Abgrenzung ist eingehalten, die Messung 655 stimmt. V-1 ist in `8ad8446` berichtigt, V-2 ist Kenntnisnahme (die Tests binden die Reihenfolge der Prüfungen mehr als der Vertrag; bei einem Refactor der Prüfreihenfolge die Texte nachziehen), V-3 siehe Register, V-4 (Änderung eines Review-Reports durch den Implementer, formal, im Commit-Text nicht genannt) ist kein Rücksprung; künftig nennt der Commit-Text die Änderung eines Reports.

**Gate-Nachweis.** Der volle `make gates` lief auf dem Arbeitsbaum der Nacharbeit grün (Exit 0, §7 *Letzter Lauf*); danach änderten sich nur Dateien auf der Positivliste (`docs/plan/planning/`, `docs/reviews/`); der Stop-Hook gab `approve`. Die Verifikation fuhr `make gates` nicht erneut, die Läufe der beiden Unit-Pakete und der Mutanten selbst.

**Größe.** 655 in der Metrik aus §4 (erste Fassung 590, Nacharbeit +65), Meldeschwelle 800 und Marke 1000 nicht erreicht, Schätzung 550 bis 750 getroffen; Liefer-Punkte 2, Schichten 2 (§8).

**Risiken aus §6 — Ausgang je Risiko** (geschlossene Menge):

| Risiko | Ausgang | Begründung |
|---|---|---|
| Zwischenstand: `replay` lädt eine Aufzeichnung mit Ende, der PGWire-Adapter kennt `connection_closed` noch nicht | eingetreten | Folge-Slice `slice-v1-abschluss-verbindungsende-replay`; er nimmt an (§1 „Übernommen aus … format“, §4 Start) |
| Der ungeteilte Plan lag über einer Review-Sitzung | eingetreten | Folge-Slice `slice-v1-abschluss-verbindungsende-format-schreiber` nimmt den Schreiber mit Roundtrip auf (§1 und DoD); Register: `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (Beleg dieses Slice), Lücke bei `slice-harness-schaetzung-planner-faktor` |
| Dieser Teil überschreitet Meldeschwelle oder Marke | entfallen | gemessen 655 unter Meldeschwelle 800 und Marke 1000, in der Schätzung; die Randform-Rückgabe an den Architect trug den Zuwachs von 590 auf 655 |
| Ort der Formregel einer Session | entfallen | entschieden vor dem Code (`Session.ValidateEnd`), geliefert und nachgewiesen (`TestSessionValidateEnd`, der Leser ruft sie) |

**Register** (Zähler = Dateien unter `evidence/`, Stand dieser Closure; Beleg je `evidence/slice-v1-abschluss-verbindungsende-format.md`):

| Eintrag | Zähler | Fund | Stand |
|---|---|---|---|
| `BEO-REPO/spec-randform-erst-im-review-entschieden` | 31 | F-1: Stelle von `ready_for_query` und `connection_closed` in der Nachrichtenfolge | verkörpert, Regel bleibt |
| `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` | 18 | F-1: bemerkt, in §7 vermerkt, nicht zurückgegeben | verkörpert, Regel bleibt |
| `BEO-REPO/plan-folgt-korrektur-nicht` | 48 | F-2: §3 nennt die falsche Abdeckungstabelle | verkörpert, Regel bleibt |
| `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | 59 | V-1: „gibt es nicht“ in §7 weiter als die Probe | verkörpert, Regel bleibt |
| `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` | 5 | kein Auftreten: 655 in 550 bis 750 | geplant (`slice-harness-schaetzung-planner-faktor`), Vermerk |
| `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` | 4 | Beleg vom Schnitt, unverändert | geplant (`slice-harness-schaetzung-planner-faktor`) |
| `BEO-REPO/test-prueft-teilstring-gegen-meldung-mit-zufallspfad` | 1 | F-4: Teilstring gegen Meldung mit `t.TempDir()`, nicht Last | neu, offen, unter der Schwelle |
| `BEO-REPO/abdeckung-vollstaendig-vor-lieferung-der-anforderung` | 1 | F-3, V-3: `vollständig` nach dem ersten von acht Slices | neu, offen, ohne Adresse (`slice-harness-abdeckung-gate` ist keine, siehe `state.md`) |

Kein Eintrag erreicht mit diesem Slice neu 3× ohne Regel; die vier Wiederauftreten betreffen Einträge, die schon `verkörpert` sind, und lösen den Retirement-Check ihrer Regeln aus (keine Regel entfernt oder gelockert). F-3 ist damit beim Planner entschieden: Register-Eintrag statt Änderung der Tabelle, Lesart `vollständig` = drei Pfade gedeckt (`SPEC-048`), keine Aussage über die Lieferung.

**Lerneintrag — benannte Spec-Lücke, Wiederaufgetreten, Regel bleibt** (`AGENTS.md` §3.9, §3.10, §3.11, §3.12); geändert wird keine Regel, kein `liegt in`. Die Spec-Lücke ist **die Stelle einer Nachricht in der Nachrichtenfolge als Randform**: `ready_for_query` vor dem Ende einer einfachen Anfrage und `connection_closed` mitten in der Folge standen weder in §6 noch in `SPEC-002`; der Implementer entschied sie im Code und vermerkte sie nur in §7 (F-1), der Architect trug sie nach (`aed8a49`). Die Nacharbeit zählte die Ausprägungen der Klasse (Anfang · Mitte · direkt davor · zwischen Fehlerantwort und Ende; einfach · Extended letzte Gruppe · nicht letzte Gruppe) und gab jeder einen Fall mit eigener Mutation; die Verifikation fuhr sie nach. Was funktionierte: Der Architect entschied vor dem Code (`690e353`) mehr Randformen als je zuvor in einem Format-Slice, die Mutationstabelle zählte je Element, und die Doppelung zweier Verteidigungen (M4b äquivalent, M4 bricht die Zusage) wurde ehrlich ausgewiesen. Neu: ein Test mit Teilstring gegen einen Zufallspfad (`t.TempDir()`) lieferte einen seltenen Fehlalarm, den der Implementer der Last zuschrieb; Review und Verifikation fanden die Ursache durch Reproduktion mit `TMPDIR` (Register, 1×). Vorschläge an den Nutzer, nicht als Eintrag und nicht von dieser Closure in `AGENTS.md` oder `.claude/` geändert: (1) `AGENTS.md` §3.12 oder `architect.md`: Die Stelle einer Nachricht in der Nachrichtenfolge ist eine Randform (Tabelle Nachricht × Stelle: Anfang, Mitte, direkt davor, direkt danach, letzte Gruppe, nicht letzte Gruppe) und steht in §6 vor dem Code. (2) `.claude/commands/implement-slice.md` Randform-Rückgabe: Eine Lesart, die §7 „Lesart“ nennt, ist eine Rückgabe und hält den Implementer an (der Satz „ein Satz im Bericht ist keine Rückgabe“ steht in §3.12 nur bei der Aufruf-Voraussetzung). Retirement-Check bestanden (keine Regel entfernt oder gelockert).

**Paarungen.** (a) Anker: Der Lerneintrag trägt kein Feld `liegt in`, keine Verkörperung; nichts zu prüfen. (b) Folge-Slice: Genannt sind `slice-v1-abschluss-verbindungsende-replay`, `slice-v1-abschluss-verbindungsende-format-schreiber`, `slice-v1-abschluss-verbindungsende-sqlite`, `slice-v1-abschluss-verbindungsende-record` und `slice-harness-schaetzung-planner-faktor`; alle existieren in `next/` und nehmen an: `grep -n "slice-v1-abschluss-verbindungsende-format"` findet die Kennung des Gebers in §1 von `-replay` (Zeile 35), `-format-schreiber` (Zeilen 33 bis 39), `-sqlite` (Zeilen 33 bis 39) und `-record` (Zeilen 35 und 37), `slice-harness-schaetzung-planner-faktor` trägt die Lücke aus `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (Register, `state.md`). Die Bindung dieses Gebers (`Session.ValidateEnd`, `model.VersionMitEnde`) steht seit `ac45571` in §1, §6 und §8 (Nachzählung) von `-sqlite`, `-record` und `-format-schreiber`. (c) Register: Alle genannten Kennungen `BEO-REPO/…` bestehen als Verzeichnis, jedes mit nicht leerem `evidence/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand). Treffer in der Sub-Area `REPO`: `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` und `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` — beide stehen im Register als `verkörpert` (`AGENTS.md` §3.14, seit slice-harness-gate-kette-stop-hook-liste), sind also keine Lücke ohne Zuweisung; dass die Prüfung des Architect den Plan trotzdem über die Marke schätzte, ist ein Auftreten trotz der Regel und löst deren Retirement-Check aus (Urteil des Planners vom 2026-10-15 im Register, `state.md`; Zuweisung eines Folge-Slice: Entscheidung des Nutzers). `BEO-REPO/meldeschwelle-ohne-artefakt` (1×): §4 nennt die Meldeschwelle neben der Marke. `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` ist verkörpert (`AGENTS.md` §3.10): die Mutationstabelle steht in DoD und §7. Keine weiteren Treffer.

Nachgezählt nach dem Schnitt (Planner, 2026-10-15, `AGENTS.md` §3.13): Liefer-Punkte: 2 (Modell mit Formregeln, YAML-Leser); der Plan trug vorher 3 (dazu Schreiber mit Roundtrip, jetzt in `slice-v1-abschluss-verbindungsende-format-schreiber`). Schichten: zwei (Kern mit dem Modell, `ARC-001`; Recording-Adapter, `ARC-008`); der Handbuch- und README-Teil liegt in `slice-v1-abschluss-verbindungsende-doku`. Schätzung: 550 bis 750 Zeilen mit Tests in der Metrik von §4 (Architect, aus der Rechnung unten; der ungeteilte Plan lag bei 850 bis 1250).

**Schätzung des Architect (2026-10-15)** für den ungeteilten Plan, in der Metrik von §4 (Summe aus Einfügungen und Löschungen, `internal` und `docs/user`, Tests eingeschlossen): Produktcode Modell etwa 70, YAML-Leser etwa 70, YAML-Schreiber etwa 40 (zusammen etwa 180); Tests Modell etwa 200, Leser etwa 250, Schreiber und Roundtrip etwa 230, Hilfen und Abdeckungstabellen etwa 50 (zusammen etwa 730). Das sind 850 bis 1000 vor der Nacharbeit der Mutationen; `slice-v1-abschluss-protokollrand` lag nach derselben Rechnung bei 800 bis 1000 und maß 1234 (Faktor 1,3 bis 1,5). Die realistische Spanne war **850 bis 1250, Mitte etwa 1050, über der Marke 1000**. Empfehlung: Schnitt vor dem Anspruch nach §4, **Leser | Schreiber**; erster Teil (dieser Slice): Modell mit Formregeln, YAML-Leser, Mutationen dazu, 550 bis 750; zweiter Teil: YAML-Schreiber mit Version nach dem Inhalt und Roundtrip, 300 bis 500, setzt den ersten voraus. Der Planner hat danach geschnitten (§1).

**Gemessen (Implementer, 2026-10-15):** erste Fassung Summe 590, nach der Nacharbeit aus dem Review Summe 655 (`git diff --numstat 6d0d697 -- internal docs/user`, Arbeitsbaum vor dem Commit der Nacharbeit) in der Metrik von §4, Liefer-Punkte 2, Schichten 2 (Kern mit dem Modell, Recording-Adapter); die Schätzung lag bei 550 bis 750, die Meldeschwelle 800 wurde nicht erreicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
