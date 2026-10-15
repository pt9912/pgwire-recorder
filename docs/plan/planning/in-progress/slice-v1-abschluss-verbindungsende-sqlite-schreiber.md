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

**Verantwortlich:** pt9912

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

**Bindung aus `slice-v1-abschluss-verbindungsende-format-schreiber`** (dort §6 *Aufrufer mit verletzter Form*, Setzung des Architect, 2026-10-15; Eintrag des Planners, 2026-10-15, `AGENTS.md` §3.13): Der Schreiber prüft die Form des Endes nicht — Stelle, Anzahl, `ready_for_query` daneben, Felder neben `type`; er schreibt, was die Aufzeichnung trägt, und eine verletzte Form ergibt eine Datei, die der Leser (`slice-v1-abschluss-verbindungsende-sqlite`) mit `PGR-E3003` ablehnt (Test: schreiben nimmt an, laden lehnt ab). Der SQLite-Schreiber trägt dieselbe Setzung (Architect, 2026-10-15, vor dem Code, in §6 und in der Spezifikation; Setzung, keine Entscheidung des Nutzers); das Gewicht liegt beim Aufrufer `slice-v1-abschluss-verbindungsende-record`, das die Bindung in §1 trägt.

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

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal docs/user`, Tests eingeschlossen, wobei die Basis der Stand des Hauptzweigs vor dem Anspruch ist (**Basis: `b5933e85de09b686e4e204d17d6bb85d0a092b37`**, eingetragen vom Planner beim Anspruch am 2026-10-15; der Übergang nach `in-progress/` ändert nichts unter den genannten Pfaden); **über 1000** geht der Slice zurück zur Zerlegung. **Meldeschwelle 800:** Erreicht die Summe 800, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik: 500 bis 800, Mitte 650 und damit 65 % der Marke (§8), unter der Grenze von 70 %, am 2026-10-15 beim Anspruch bestätigt; Schätzung des Architect nach der Prüfung (Produktcode 40 bis 60, Tests mit Faktor 5 bis 10 samt Tabellen 400 bis 650, `docs/user` 20 bis 50): 480 bis 760, Mitte 620 und damit 62 %, unter 70 %, kein Schnitt (Faktor Tests zu Produktcode mindestens 5, gemessene Datenpunkte 8 und 11); schneidet die Prüfung des Architect höher, ist der Schnitt Version in der Transaktion | Roundtrip (Zuschnitt: Planner).
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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Setzungen des Architect, 2026-10-15). Ort: `LH-FA-22.a` *Ende und Version beim Schreiben von `sqlite`*, `SPEC-001` *Schreiber in `sqlite`*, `SPEC-043` (Ende des Servers). Keine Abweichung vom Lastenheft (`LH-FA-25` nennt nur die neuere Formatversion; die Wahl in der Transaktion steht in [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) Punkt 4).

- *Erster Schreibvorgang:* legt `meta.version` = 1 an, auch wenn die erste Session ein Ende trägt, und gleich welchen Wert das Modellfeld `version` hat (0, 1, 2, 3, negativ). Bisher schreibt `legeSchemaAn` das Feld des Modells; das ändert dieser Slice (`model.Version`).
- *Welche Sessions mit Ende die Version setzen:* eine Session trägt ein Ende, wenn irgendeine Antwort irgendeiner ihrer Interaktionen `connection_closed` ist (`responses` und `server` einer Gruppe), an jeder Stelle, auch einer falschen, wie `Recording.HasEnd` im YAML. Der Adapter wendet `HasEnd` auf eine Aufzeichnung mit genau dieser Session an; das Modell bekommt keine Methode (sonst dritte Schicht, §4).
- *Wer setzt die Version:* nur die Transaktion der ersten Session mit Ende, die die Datei noch nicht trägt (Stand im Schreiber: gesetzt, sobald die Transaktion abgeschlossen ist, vor dem Schließen der Datenbank). Sie ändert genau die Zeile `version` mit dem Wert 1; findet sie keine solche (Zeile fehlt, Wert ist 2), ist das ein Fehler der Transaktion, `PGR-E3001`, die Session wird nicht ergänzt. Jede andere Transaktion lässt `meta` unberührt; die Reihenfolge der Anweisungen in der Transaktion ist nicht zugesagt.
- *Transaktion der ersten Session mit Ende scheitert* (auch beim Setzen der Version): zurückgerollt, keine Zeile der Session, Version 1, der nächste Schreibvorgang ergänzt und setzt 2 — `LH-FA-22.a`. Die Datei ist in jedem Fall eine gültige Aufzeichnung: Die temporäre Datei ist vor der ersten Transaktion verschoben, das Aufräumen von `legeAn` ist für die Version nicht berührt (kein neuer Weg dort).
- *Schließen der Datenbank scheitert nach der Transaktion* (Ende-Session): Version 2 steht, der nächste Schreibvorgang ändert sie nicht und ergänzt die Session nicht doppelt — `LH-FA-22.a`.
- *Abbruch des Kontexts:* unterbricht die Transaktion mit dem Ende nicht (Bestand `TestSQLiteAbbruchUnterbrichtNicht`) — `LH-FA-22.a`.
- *`PRAGMA user_version`:* wird nicht geschrieben (Stand 0 bleibt) — `SPEC-043`.
- *Zweite und spätere Sessions mit Ende, Session ohne Ende nach der Version 2, Ende nur in einer späteren oder nicht letzten Session, ohne Sessions, nur `empty_sessions`:* Version nach dem Inhalt wie in `SPEC-001` *Schreiber in `sqlite`*; ohne Sessions und mit nur leeren Sessions bleibt 1.
- *Darstellung des Endes:* Zeile von `message` mit `kind` `connection_closed`, `direction` `server`, `fields` genau der Text `{}`; `sequence` der Interaktion; `position` in der Zählung der übrigen Nachrichten; `group_no` `NULL` bei einer einfachen, die Nummer der Gruppe bei einer Extended-Interaktion (die Zeile steht unter den Server-Nachrichten der Gruppe, nach den Client-Nachrichten). Sie liegt an der Stelle des Modells, nicht an einer vom Schreiber gewählten — `SPEC-043`. Reihenfolge der Zeilen in der Datei: nicht zugesagt, die Tests ordnen nach den Nummern. Eine Aufzeichnung ohne Ende ergibt dieselben Zeilen wie bisher, `meta.version` = 1.
- *Aufrufer mit verletzter Form* (Bindung aus `slice-v1-abschluss-verbindungsende-format-schreiber`, bestätigt für SQLite): der Schreiber prüft die Form des Endes nicht; er schreibt, was das Modell trägt, und der Leser lehnt mit `PGR-E3003` ab — `SPEC-043`, `SPEC-001`. Das Gewicht liegt beim Aufrufer `slice-v1-abschluss-verbindungsende-record` (Bindung dort in §1 und §6 um `sqlite` ergänzt, Commit dieser Prüfung).
- *`connection_closed` mit Feldern im Modell* (`Tag`, `Fields`, `Notice` …): der Schreiber schreibt sie in `fields`, wie der YAML-Schreiber, der Leser lehnt sie ab (`PGR-E3003`); er verwirft sie nicht still — `SPEC-043`.
- *Stelle des Endes (die Nachricht beendet die Folge, `AGENTS.md` §3.12):* der Schreiber lehnt keine Stelle ab und setzt für jede die Version 2; der Leser entscheidet. Je Stelle ein Fall „schreiben nimmt an (Version 2, keine Meldung) · laden“: Anfang der Folge (`PGR-E3003`) · Mitte (`PGR-E3003`) · direkt nach der `error_response`, letzte Nachricht (lädt) · letzte Gruppe (lädt) · nicht letzte Gruppe (`PGR-E3003`) · letzte Interaktion einer Session (lädt) · nicht letzte Interaktion (`PGR-E3003`) · nicht letzte Session (lädt) · zwei Enden in einer Interaktion (`PGR-E3003`) · zwei Enden in zwei Interaktionen einer Session (`PGR-E3003`) · `ready_for_query` vor dem Ende (`PGR-E3003`) · `ready_for_query` nach dem Ende (`PGR-E3003`) · Feld neben `type` (`PGR-E3003`). Die Fälle sind dieselben wie in `TestSchreibenNimmtFormverletzungAn` (YAML); eine geteilte Tabelle, die beide Formate durchläuft, spart Zeilen.
- *Reihenfolge mehrerer Ablehnungsgründe beim Schreiben:* das Ende fügt keinen Ablehnungsgrund hinzu; es gilt der Bestand: zuerst die Nummern der Sessions (`PGR-E1000`, Datei unberührt, eine noch nicht angelegte wird nicht angelegt), dann Anlegen und Öffnen (`PGR-E3001`), dann die Transaktion (`PGR-E3001`). Ein Fall: Session mit Ende und falscher Nummer im ersten Schreibvorgang ist `PGR-E1000`, keine Datei. Beim Laden gilt `SPEC-001` *Rangfolge der Prüfung* (Leser, nicht dieser Slice).
- *Roundtrip:* einfach · `Sync`-Gruppe · `Flush`-Gruppe · `offset_ms` (auch an der Interaktion mit Ende) · `empty_sessions` · Version 2 ohne Ende (Datei laden, schreiben, laden: gleich bis auf `version`, das Schreiben ergibt 1) · dieselbe Aufzeichnung als YAML und als SQLite liest gleich, auch in `version`; ein späterer Schreibvorgang an eine Datei der Version 2 ergänzt Sessions, ohne die Version zu ändern — `SPEC-001`.
- *Konstante der Version 2* (Bindung aus `slice-v1-abschluss-verbindungsende-format`): `model.VersionMitEnde` und `model.Version`, kein Literal — `SPEC-001`.
- *Aufruf-Voraussetzung:* der Schreiber vergibt keine Werte, die über mehrere Aufrufe geordnet oder eindeutig sein müssten; die Zusagen sind Zustandsprüfungen ohne Zusage über Reihenfolge oder Gleichzeitigkeit (`AGENTS.md` §3.12) — keine Voraussetzung. Aufrufarten des Treibers: keine; kein Hook oder Gate liest ein Ergebnis dieses Schreibers. Vorhandene Zieldatei und Journal prüft `Prepare`, nicht dieser Slice.

**Mittel für den Fehler in der Transaktion** (Antwort auf das erste Risiko, Architect): Kein neuer Eingriff im Produktcode und keiner in `export_test.go`. Zwei Wege reichen. (1) Ein Auslöser in der Datei, den der Test mit `rohSchreibe` anlegt: `CREATE TRIGGER … BEFORE UPDATE ON "meta" BEGIN SELECT RAISE(ABORT, '…'); END;`; der Schreiber öffnet die Datei und stößt auf ihn beim Setzen der Version, nachdem die Zeilen der Session eingefügt sind. Die Folge prüft Zeilen und Werte: `PGR-E3001` mit der Ursache im Text (Anwesenheit der Ursache), keine Session-Zeile, `meta.version` 1; nach `DROP TRIGGER` ergänzt der nächste Schreibvorgang und setzt 2. Der Leser ignoriert Auslöser (`LH-FA-22.a`). (2) Eine Interaktion mit doppelter `sequence` in der Ende-Session scheitert beim Einfügen (Primärschlüssel); im ersten Schreibvorgang belegt das, dass die Datei mit Version 1 und ohne Session bleibt (erzeugte der Schreiber beim Anlegen schon 2, bliebe der Fall rot) und dass die Version nicht in einer eigenen Transaktion davor steht.

**Zählung der Mutationstabelle** (`AGENTS.md` §3.10, ein Fall und eine Mutation, die nur ihn bricht, je Element; Fixtures so, dass kein anderes Element den Fall mitfärbt):

- *Wahl der Version:* mit Ende 2 · ohne Ende 1 · Modellfeld 1 mit Ende schreibt 2 · Modellfeld 2 ohne Ende schreibt 1 · Modellfeld 0, 3 und -1 (je einer) · ohne Sessions · nur `empty_sessions` · Ende nur in einer späteren Session (nach dem ersten Schreibvorgang noch 1) · Ende nur in einer nicht letzten Session · Ende in einer Extended-Gruppe, ohne dass irgendwo eine einfache Interaktion ein Ende trägt · Ende in einer einfachen Interaktion, ohne Gruppen mit Ende (Mutation: `HasEnd` sieht nur `Responses` oder nur `Groups`) · erster Schreibvorgang mit Ende in der ersten Session (Anlegen 1, Transaktion 2).
- *Wer setzt:* zweite Session mit Ende und Session ohne Ende nach der Version 2 ändern `meta` nicht (Auslöser auf `meta` nach der Version 2: der Schreibvorgang gelingt; Mutation: Setzen bei jeder Ende-Session) · Zeile `version` fehlt · Zeile `version` steht auf 2 (je `PGR-E3001`, Session nicht ergänzt; Mutationen: Zeilenzahl nicht geprüft · ohne `value = '1'`) · `user_version` bleibt 0.
- *Fehlerweg:* Auslöser beim Setzen · doppelte `sequence` im ersten Schreibvorgang · Schließen scheitert nach der Ende-Transaktion (nächster Schreibvorgang gelingt, Version 2; Mutation: Stand erst nach dem Schließen gesetzt) · Abbruch des Kontexts mit Ende · der nächste Schreibvorgang nach einem Fehler ergänzt und setzt 2.
- *Darstellung:* Zeile des Endes in einer einfachen Interaktion (`position`, `group_no` `NULL`, `fields` genau `{}`, Form des Werts: Text mit zwei Zeichen) · in einer `Sync`-Gruppe · in einer `Flush`-Gruppe (`group_no`, `position` über beide Richtungen) · mit `offset_ms` an der Interaktion · Aufzeichnung ohne Ende: alle Zeilen der vier Tabellen wie bisher · `connection_closed` mit Feld im Modell: Feld in `fields`, Laden `PGR-E3003`.
- *Stellen* (oben): je Stelle Version 2 und kein Fehler beim Schreiben, Laden nach dem Muster; Mutation je Stelle: der Schreiber lehnt sie ab oder schreibt 1 · Session mit Ende und falscher Nummer `PGR-E1000`, keine Datei.
- *Roundtrip:* die sechs Formen oben, je Form SQLite für sich und gleich dem YAML.

Fixtures: `sEnde`, `sSession`, `sRec`, `sEinfach`, `sExtended`, `sGruppe` aus `yaml_schreiben_ende_test.go` und `sessionsMit`, `rohSchreibe`, `roh`, `zeilen` aus `sqlite_hilfen_test.go` sind im selben Paket vorhanden. Eine Abwesenheitszusage (kein Ende, keine Ablehnung, keine Version) prüft die Form des Werts (`err == nil`, Zeilenzahl, Wert von `meta.version`), nicht den Wortlaut einer Meldung (`BEO-REPO/abwesenheitszusage-an-wortlaut-der-meldung-geprueft`).

**Risiken:**

- Ein Fehler in der Transaktion ist ohne Einspeisung nicht erreichbar. Geprüft (Architect, 2026-10-15): Das Mittel ist ein Auslöser in der Datei (`rohSchreibe`) und eine doppelte `sequence`, beides ohne Änderung am Produktcode oder an `export_test.go` (oben *Mittel für den Fehler in der Transaktion*) — **Ausgang:** offen bis Closure; zu belegen ist, dass beide Fälle rot werden, wenn die Version außerhalb der Transaktion gesetzt wird.
- Der Diff liegt über der Marke (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`: Tests mit Mutationen je Element, Roundtrip je Form; Schätzung des Architect 480 bis 760, Mitte 620 (62 %), Meldeschwelle 800) — **Ausgang:** offen bis Closure.

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

Nachgezählt bei der Prüfung des Architect vor dem Code (2026-10-15, `AGENTS.md` §3.13): Liefer-Punkte 2 (Schreiber mit Version in der Transaktion; Roundtrip). Schichten **zwei**: Recording-Adapter (`ARC-008`) und `docs/user` (die Abdeckungstabellen, die `make abdeckung` aus den Deklarationen der Tests schreibt; `docs/user` zählt als Schicht, §3.13 *Schichten zählen alle Pläne gleich*) — die Zählung „eine“ vom Anlegen war um diese Tabellen zu niedrig, die Grenze von zwei hält. Das Modell bleibt unberührt (`Recording.HasEnd` auf eine Aufzeichnung mit einer Session, keine neue Methode), der Port ebenso. Schätzung in der Metrik von §4: 480 bis 760, Mitte 620, 62 % der Marke 1000, unter 70 %.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
