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

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-22.a` · `SPEC-001` · `SPEC-002` · `SPEC-041` · `SPEC-043` · `ARC-008`

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

**Bindung aus `slice-v1-abschluss-verbindungsende-format`** (Closure des Gebers, Planner, 2026-10-15, `AGENTS.md` §3.13; dort §6 *Ort der Formregel einer Session*): Die Formregel einer Session (höchstens ein Ende, in der letzten Interaktion) liegt als `Session.ValidateEnd() error` im Modell, und der SQLite-Leser ruft sie je Session nach den Interaktionen, statt sie nachzubauen; `Interaction.Validate` allein lehnt ein Ende in einer Interaktion vor der letzten und ein zweites Ende nicht ab. Die Konstante `model.VersionMitEnde` (2) benutzt der Schreiber in `slice-v1-abschluss-verbindungsende-sqlite-schreiber`; der Leser vergleicht die Version mit `model.Version` und `model.VersionMitEnde`, nicht mit einem Literal. Der Aufruf geschieht über `Unmarshal` (§6 *Ort des Aufrufs*), nicht im Leser selbst.

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

- [x] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), `SPEC-001`, `SPEC-043`: Der SQLite-Leser nimmt `meta.version` 1 und 2 an, auch Version 2 ohne Ende. `connection_closed` in einer Datei der Version 1, an einer anderen Stelle als der Form, nicht nach `FATAL`/`PANIC` oder mit Feldern sind `PGR-E3003`; ein Ende in einer Interaktion vor der letzten und ein zweites Ende in einer Session sind `PGR-E3003` über `Session.ValidateEnd`; jede andere Version ist `PGR-E3002` (Test je Randform aus §6, die Dateien von Hand gebaut).
- [x] Je Zusage des neuen Vertrags eine Mutation, vom Implementer rot gesehen; die Tabelle (Zusage · Mutation · roter Test) steht in §7 (`AGENTS.md` §3.10).
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
| `internal/adapters/driven/recording/sqlite_lesen.go` | update | Leser: `liesMeta` und `pruefeMeta` liefern den Text der Version (`model.Version` oder `model.VersionMitEnde`, sonst `PGR-E3002`), `liesAufzeichnung` schreibt ihn statt fest `1` in den Baum; die Nachricht `connection_closed` und die Formregeln (`Validate`, `Session.ValidateEnd`) laufen über `Unmarshal` (Bindung aus `slice-v1-abschluss-verbindungsende-format`, Ort in §6). `sqlite_lesen_baum.go` und `sqlite_json_lesen.go` bleiben unverändert: Das Ende ist eine Nachricht wie jede andere. |
| `internal/adapters/driven/recording/sqlite_lesen_ende_test.go` (neu), `sqlite_lesen_inhalt_test.go`, `leser_test.go`, `test/integration/sqlite_lesen_e2e_test.go` | update | Lesen: Version 1 und 2, Negative je Randform, Reihenfolge der Zeilen von `interaction` und `session` mit Ende, Dateien von Hand mit SQL gebaut; die bestehenden Fälle, die `2` als unbekannte Version führten (auch `TestE2ELesenAbgelehnt`, ein Fall, keine neue Zusage), nehmen `3` |
| `docs/user/abdeckung-unit.md` | update (`make abdeckung`) | Tabelle aus den Deklarationen der Tests (`abdeckung-e2e.md` bleibt, der Slice hat keinen E2E-Test) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` liegt in `done/` (Voraussetzung; Schritt 35a: Modell, `Session.ValidateEnd`, YAML-Leser); `slice-v1-abschluss-verbindungsende-format-schreiber` (35b) liegt ebenfalls in `done/`, ist aber für diesen Slice keine Voraussetzung. Dieser Slice ist Schritt 36a der Reihenfolge in §5 von welle-v1-abschluss. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat fa8486b..<letzter Commit des Slice> -- internal docs/user`, Tests eingeschlossen, wobei die Basis `fa8486b` ist, der Stand des Hauptzweigs vor dem Anspruch (der Übergang nach `in-progress/` ändert nichts unter den genannten Pfaden); **über 700** geht der Slice zurück zur Zerlegung. **Meldeschwelle 550:** Erreicht die Summe 550, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners in dieser Metrik: 260 bis 420, Mitte 340 und damit 49 % der Marke (§8). Prüfung des Architect (2026-10-15): 220 bis 420, Mitte 320, 46 % der Marke, unter 70 %, kein Schnittvorschlag (§8).
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

**Randformen** (`AGENTS.md` §3.12), in der Spezifikation entschieden vor dem Code (Prüfung des Architect, 2026-10-15). Ort, wo nicht anders genannt: `SPEC-043` Aufzählungspunkt *Ende des Servers in `sqlite`* (Ort der Darstellung, der Stelle nach Nummern, der Version und der Rangfolge), daneben `SPEC-001` *Rangfolge der Prüfung*, `SPEC-002`, `SPEC-041`, `LH-FA-25.a`, `LH-FA-22.a` *Lesen von `sqlite`*. Fälle sind Dateien, die der Test von Hand baut (SQL-Anweisungen auf einer gültigen Datei, nicht der SQLite-Schreiber: Er schreibt bis `-sqlite-schreiber` `rec.Version` unverändert und ist kein Beleg).

*Darstellung des Endes:*
- Das Ende ist eine Zeile von `message` mit `kind` `connection_closed`, `direction` `server`, `fields` `{}`; keine neue Tabelle oder Spalte, das neutrale Schema bleibt (`SPEC-043`). Fälle: einfache Interaktion, Extended-Interaktion mit `Sync`-Gruppe, mit `Flush`-Gruppe, mit `offset_ms`; das Modell trägt das Ende als letzte Antwort ohne Felder (Mutation: der Leser lässt die Zeile aus).
- `fields` mit einem Schlüssel — `{"tag":"x"}`, `{"x":null}`, `{"":1}`, `{"type":"x"}` — ist `PGR-E3003` (Gegenstück zum Schlüssel neben `type`; Mutation: der Leser verwirft die Felder des Endes).
- `direction` `client` mit `kind` `connection_closed` ist `PGR-E3003`; `kind` `Connection_Closed` und `connection_closed ` (Leerzeichen) sind `PGR-E3003`.
- Stelle folgt den Nummern, nicht der Reihenfolge der Zeilen: die Zeile des Endes in der Datei vor der der `error_response`, mit richtigen `position`-Werten, wird angenommen (Mutation: die Sortierung entfällt).

*Version:*
- `meta.version` `1` und `2` werden angenommen, die Aufzeichnung trägt die Version der Datei (`Recording.Version`), auch `2` ohne Ende. Der Vergleich nutzt `model.Version` und `model.VersionMitEnde`, kein Literal. Mutationen: der Leser setzt den Baum fest auf `1` (heutiger Stand, `sqlite_lesen.go`; Datei der Version 2 mit Ende wird `PGR-E3003`), fest auf `2` (Datei der Version 1 mit Ende wird angenommen), nur `1` angenommen.
- Kanonische Zahlen außer `1` und `2` sind `PGR-E3002`: `0`, `3`, `-1`, `10`, `20`, `99999999999999999999` (je ein Fall; `10` und `20` fangen einen Vergleich als Text); andere Schreibweisen (`02`, `+2`, ` 2`, leer) bleiben `PGR-E3003` wie für `1` (bestehende Fälle bleiben). Die Meldung zu `PGR-E3002` nennt keine Version (`SPEC-043`).
- `connection_closed` in einer Datei der Version 1 ist `PGR-E3003` (`SPEC-001`, `SPEC-043`); `PRAGMA user_version` ändert die Version nicht: `user_version` 2 mit `meta.version` 1 und Ende ist `PGR-E3003`, `user_version` 1 mit `meta.version` 2 und Ende wird angenommen (Mutation: der Leser liest `user_version`).
- *Version 2 ohne Ende:* angenommen — `SPEC-001`, `SPEC-043`.

*Rangfolge mehrerer Ablehnungsgründe:* Zugesagt sind nur die Stufen Lesbarkeit der Datei, Formatkennung, Version, Inhalt (`SPEC-001`); `PGR-E3002` geht jedem Fehler des Inhalts vor, auch einem Ende an falscher Stelle (Fall: Version 3 mit Ende vor der `error_response`; Mutation: die Version wird nach dem Inhalt geprüft). Eine Datei der Version 1 mit Ende ist `PGR-E3003`, nicht `PGR-E3002`. Mehrere Fehler des Endes im Inhalt: derselbe Code, welcher gemeldet wird, ist nicht zugesagt. Ort: `SPEC-001`, `SPEC-043`, `SPEC-002`.

*Stellen des Endes* (`SPEC-002`, `SPEC-041`; je Stelle ein Fall mit der Datei in `sqlite`; alle `PGR-E3003`, wo nicht „angenommen“ steht): Anfang (erste Nachricht der Interaktion) · Mitte (eine Nachricht folgt) · vor der `error_response` · zwischen `error_response` und Ende eine `notice_response` · nach `ERROR` statt `FATAL` · `V` `PANIC` angenommen · `V` fehlt, `S` `FATAL`: angenommen · `V` leer, `S` `FATAL` · `Fatal` · `ready_for_query` vor und nach dem Ende (einfach) · zweites Ende in derselben Interaktion · korrekt geformtes Ende in der Interaktion 1 von 2 · Ende in zwei Interaktionen einer Session · Ende in einer Session, die nicht die letzte ist: angenommen · Extended: Ende nicht in der letzten Gruppe, nicht die letzte Server-Nachricht der letzten Gruppe, `ready_for_query` in der Gruppe mit Ende, vorletzte Gruppe mit `Sync`, letzte Gruppe ohne `Sync` und `Flush`. Die Zusagen liegen im gemeinsamen Code (`Unmarshal`, `Interaction.Validate`, `Session.ValidateEnd`), den der SQLite-Leser über den Baum benutzt (akzeptiertes Negativ: Die Mutationen dort hat `slice-v1-abschluss-verbindungsende-format` rot gesehen; dieser Slice belegt je Stelle den Weg der Übersetzung durch einen Fall und fährt diese eine Mutation selbst: der Aufruf von `Session.ValidateEnd` entfällt, die Fälle „Ende in der Interaktion 1 von 2“ und „Ende in zwei Interaktionen“ werden rot).
- *Ort des Aufrufs von `Session.ValidateEnd`* (Setzung des Architect; Bindung aus `slice-v1-abschluss-verbindungsende-format`): Der SQLite-Leser ruft sie nicht selbst. Er baut den Baum der Aufzeichnung und lässt `Unmarshal` lesen (`ladeInhalt`), und dessen Prüfung je Session ruft `ValidateEnd`; ein zweiter Aufruf im Leser wäre der Nachbau, den die Bindung ausschließt. Stelle im Ablauf: Fehler der Bibliothek und der Zeilen (Arten, Nummern, verwaiste Zeilen, JSON) vor dem Baum, danach die Prüfung des Baums; nicht zugesagt ist die Ordnung innerhalb des Inhalts (`SPEC-001`).

*Datei, Fremdschlüssel, `PRAGMA`:* Das Ende ändert nichts an Öffnen, Lesetransaktion, Fehlerklassen der Bibliothek oder verwaisten Zeilen; eine Zeile des Endes ohne Interaktion ist wie jede Nachricht ohne Interaktion `PGR-E3003` (akzeptiertes Negativ: kein eigener Fall, die Prüfung hängt nicht an `kind`). `foreign_keys` wird beim Lesen nicht erzwungen (`LH-FA-22.a`). Das neutrale Schema bleibt (`make schema-validate` unberührt).

*Verhältnis der Leser:* gleiche Verstöße geben in `yaml` und `sqlite` denselben Code, weil `sqlite` über denselben Baum und dieselbe Prüfung liest (`SPEC-043`); einen Vergleich beider Leser auf derselben ungültigen Aufzeichnung gibt es nicht (akzeptiertes Negativ: eine zweite Darstellung desselben Falls; der Roundtrip gültiger Aufzeichnungen liegt in `slice-v1-abschluss-verbindungsende-sqlite-schreiber`).

*Bindung an den SQLite-Schreiber* (`AGENTS.md` §3.13; Bestätigung des Architect an `slice-v1-abschluss-verbindungsende-sqlite-schreiber`, 2026-10-15): Die Setzung aus `-format-schreiber` — der Schreiber prüft die Form des Endes nicht — gilt für `sqlite` ebenso; `SPEC-043` sagt es jetzt, und `-sqlite-schreiber` trägt sie bereits in §1, §6 und DoD 2. Der Leser dieses Slice ist die Stelle, die eine verletzte Form ablehnt.

- *Aufruf-Voraussetzung und Aufrufarten des Treibers:* der Leser ist eine Zustandsprüfung einer Datei ohne Zusage über Reihenfolge oder Gleichzeitigkeit von Aufrufen; kein Hook und kein Gate liest ein Ergebnis dieses Slice — keine (`AGENTS.md` §3.12).

**Risiken:**

- Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`): Der Planner hat vor dem Anspruch geschnitten (Schätzung 260 bis 420, Mitte 49 % der Marke); die Rückführung in §4 benennt die Marke und die Meldeschwelle — **Ausgang: entfallen** — gemessen 387 in der Metrik von §4 (55 % der Marke 700, unter der Meldeschwelle 550, innerhalb der Schätzung 220 bis 420 des Architect); das Review lief in einer Sitzung ohne Größen-Finding.

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

### Belege des Implementers

**Läufe.** Der erste volle `make gates` war rot (Exit 2): `TestE2ELesenAbgelehnt` führte „sqlite, Version 2“ als unbekannte Version; der Fall nimmt `3` (derselbe Fehler wie in `slice-v1-abschluss-verbindungsende-format`). Zweiter voller Lauf `make gates` grün (Exit 0, ohne Flanke `p7-int-treiber`). Zwischendurch `go test ./internal/adapters/driven/recording/` per Bind-Mount ohne Netz grün. Nacharbeit zum Review: voller `make gates` grün (Exit 0, 12 min 56 s, ohne Flanke `p7-int-treiber`); Zwischenläufe per Bind-Mount ohne Netz.

**Mutationstabelle** (`AGENTS.md` §3.10). Jede Mutation einzeln in einer Kopie des Baums (`/…/scratchpad/M`), `go test ./internal/adapters/driven/recording/ -run 'TestLesenEnde|TestLesenMeta|TestLesenMeldung|TestLesenGleiche|TestLesenUserVersion|TestLesenRoundtrip'` im Image `pgwire-recorder:test` mit `--network=none` und Bind-Mount (die Datei wird neu geschrieben, danach zurückgesetzt; der Mutant kommt an, Übersetzungsfehler wurden erkannt und der Mutant korrigiert). „Rot“ nennt die Tests, die rot wurden; weitere Fälle derselben Tests, die mitfärben, sind nicht genannt.

| Zusage | Mutation | roter Test |
|---|---|---|
| das Ende ist eine Zeile von `message` und wird geladen (einfach) | `interaktionBau.knoten` lässt die Zeile `connection_closed` aus | `TestLesenEndeAngenommen` (alle Fälle mit einfachem Ende); mitgefärbt `TestLesenEndeAbgelehntStelle` „ready_for_query vor/nach dem Ende“ |
| dasselbe (Extended) | `interaktionBau.gruppen` lässt die Zeile aus | `TestLesenEndeAngenommen` (Sync-, Flush-, vordere Flush-Gruppe, Zeile vor der Fehlerantwort); `TestLesenEndeAbgelehntDarstellung` „direction client, extended“, `…Stelle` „nicht in der letzten Gruppe“, „ready_for_query in der Gruppe“ (mitgefärbt, weil das Ende fehlt) |
| `fields` ohne Schlüssel (Felder des Endes werden nicht verworfen) | `nachrichtBau.knoten` verwirft alle Felder von `connection_closed` | `TestLesenEndeAbgelehntDarstellung`: „fields mit tag“, „mit Wert null“, „mit leerem Namen“, „mit type“ |
| Schlüssel `tag` ist `PGR-E3003` | verwirft nur den Schlüssel `tag` | nur „fields mit tag“ |
| Schlüssel mit Wert `null` | verwirft nur Paare mit Wert `null` | nur „fields mit Wert null“ |
| Schlüssel mit leerem Namen | verwirft nur den Schlüssel `""` | nur „fields mit leerem Namen“ |
| Schlüssel `type` | verwirft nur den Schlüssel `type` | nur „fields mit type“ |
| `direction` `client` mit `connection_closed` (einfach) | Prüfung `richtung != server` gilt nicht für `connection_closed` | nur „direction client, einfach“ |
| dasselbe (Extended) | beide Zweige der Client-Nachrichten in `gruppen` lassen `connection_closed` aus (es landet in `server`) | **grün, äquivalent:** In einer Extended-Gruppe stehen die Client-Nachrichten vor den Server-Nachrichten (`SPEC-043`), eine Client-Zeile `connection_closed` landet also vor `error_response` in `server` und ist nicht die letzte Nachricht; `Group.validate` lehnt sie mit demselben Code ab. Der Fall „direction client, extended“ belegt den Code, nicht die Prüfung der Richtung |
| `kind` genau (Großschreibung) | `nachrichtVon` setzt `kind` in Kleinbuchstaben | nur „kind in Großbuchstaben“ |
| `kind` genau (Leerzeichen) | `nachrichtVon` entfernt Leerraum von `kind` | nur „kind mit Leerzeichen“ |
| die Stelle folgt den Nummern, nicht der Reihenfolge der Zeilen | die Sortierung der Nachrichten entfällt | `TestLesenEndeAngenommen` „Zeile des Endes vor der Fehlerantwort“ (einfach und extended) |
| Interaktionen werden nach `sequence` geordnet, nicht nach der Reihenfolge der Zeilen (Datei mit Ende) | die Sortierung der Interaktionen in `liesInteraktionen` entfällt | `TestLesenEndeNummernOrdnen` „Zeilen von interaction umgekehrt“ (`PGR-E3003` „Interaktion 1 trägt die Nummer 2“ statt geladen), `TestLesenReihenfolgeDerZeilen/interaction` |
| Sessions werden nach `id` geordnet (Datei mit Ende in der zweiten Session) | die Sortierung der Sessions in `liesSessionTabelle` entfällt | `TestLesenEndeNummernOrdnen` „Zeilen von session umgekehrt“ (`PGR-E3003` „Session 1 trägt die Kennung 2“), `TestLesenReihenfolgeDerZeilen/session` |
| `Session.ValidateEnd` ruft `Unmarshal`, der Leser baut sie nicht nach | der Aufruf in `sessionFromDTO` (`yaml.go`) entfällt | `TestLesenEndeJeSession` „Interaktion 1 von 2“ und „zwei Interaktionen“; `TestLesenEndeAbgelehntStelle` dieselben zwei Fälle |
| der Leser liefert die Version der Datei (1 und 2) | der Baum trägt fest `"1"` | `TestLesenEndeAngenommen` (alle Fälle der Version 2: `PGR-E3003`) |
| dasselbe | der Baum trägt fest `"2"` | `TestLesenEndeAngenommen` „Version 1 ohne Ende“ (Version 2 statt 1), `TestLesenEndeAbgelehntDarstellung` „Version 1, einfach/extended“ (angenommen statt `PGR-E3003`) |
| Version 2 wird angenommen | nur `model.Version` angenommen | `TestLesenEndeAngenommen` (alle Fälle der Version 2), `TestLesenMeta` „version 2“ |
| Version 2 ohne Ende wird angenommen | Version 2 ohne `connection_closed` im Text ist `PGR-E3002` | nur „Version 2 ohne Ende“, „…, extended“ (`TestLesenEndeAngenommen`) und `TestLesenMeta` „version 2“ |
| die Zahl wird als Zahl verglichen (`10`) | `strings.HasPrefix(version, "1")` statt Vergleich | `TestLesenEndeVersionen` „Version 10“, über die Meldung (siehe Anmerkung: der Code bleibt `PGR-E3002`) |
| dasselbe (`20`) | `strings.HasPrefix(version, "2")` | „Version 20“, über die Meldung |
| `0` ist `PGR-E3002` | `0` angenommen | „Version 0“, über die Meldung |
| `3` ist `PGR-E3002` | `3` angenommen | „Version 3“, „Version 3, Ende vor der Fehlerantwort“ (beide Fälle) über die Meldung; `TestLesenMeta` „Version vor Inhalt/fehlender Tabelle/unbekanntem Schlüssel“ |
| `-1` ist `PGR-E3002` | `-1` angenommen | „Version -1“, über die Meldung |
| eine Zahl über 64 Bit ist `PGR-E3002` | die Prüfung der Menge in `pruefeMeta` entfällt (der Leser vergleicht Texte, es gibt keine Umwandlung; die Zahl bleibt als Text im Baum und `kopfVorpruefung` lehnt sie als nicht kanonisch ab) | „Version 99999999999999999999“ (`PGR-E3003` statt `PGR-E3002`), `TestLesenMeta` zwei Fälle; hier ist die Menge in `pruefeMeta` am Code rot; bei den kleinen Zahlen mit sauberem Inhalt nur über die Meldung, mit Inhalt der Rangfolge am Code (`TestLesenMeta` „Version vor Inhalt“) |
| `02` ist `PGR-E3003` | führende Nullen werden vor der Prüfung entfernt | „Version "02"“, `TestLesenMeta` „01“ |
| `+2` ist `PGR-E3003` | ein `+` wird vor der Prüfung entfernt | „Version "+2"“, `TestLesenMeta` „+1“ |
| ` 2` ist `PGR-E3003` | Leerzeichen davor werden entfernt | „Version " 2"“, `TestLesenMeta` „ 1“ |
| `2 ` ist `PGR-E3003` | Leerzeichen danach werden entfernt | „Version "2 "“, `TestLesenMeta` „1 “ |
| die Meldung zu `PGR-E3002` nennt keine Version | die Meldung trägt die Zahl in irgendeiner Form: `Version <Zahl>` (Mutant: die Menge in `pruefeMeta` entfällt, die Meldung stammt aus `Unmarshal`), `Fassung <Zahl>`, `version=<Zahl>` in `pruefeMeta` | `TestLesenEndeVersionen` alle E3002-Fälle (die Meldung enthält die Zahl der Datei nach Abzug von Pfad und Code in keiner Form, auch nicht ohne das Wort `Version`), `TestLesenMeldungOhneWertUndZeile` „meta version unbekannt“ |
| `PRAGMA user_version` ändert die Version nicht | der Leser setzt die Version aus `user_version` 1 oder 2 | `TestLesenEndeAngenommen` „user_version 1, meta.version 2, Ende“ und „user_version 2, meta.version 1 ohne Ende“, `TestLesenEndeAbgelehntDarstellung` „user_version 2, meta.version 1, Ende“, `TestLesenUserVersionIrrelevant` |
| `PGR-E3002` geht dem Inhalt vor | die Version wird erst nach dem Inhalt geprüft (der Baum trägt `2` für eine fremde Version, `PGR-E3002` folgt nach der Form) | `TestLesenEndeVersionen` „Version 3, Ende vor der Fehlerantwort“ (einfach und Zeilen umgekehrt; `PGR-E3003` statt `PGR-E3002`), `TestLesenMeta` „Version vor Inhalt/fehlender Tabelle/unbekanntem Schlüssel“ |

**Anmerkung zu den Zahlen.** `Unmarshal` prüft die Version noch einmal (`kopfVorpruefung`) und lehnt eine kanonische Zahl außer 1 und 2 mit `PGR-E3002` ab, wenn `pruefeMeta` sie durchlässt. Eine Mutation an der Menge in `pruefeMeta` (`0`, `3`, `-1`, `10`, `20`) ändert deshalb den Code nicht, nur die Meldung (aus `Unmarshal` statt der ohne Zahl). Grenze: Mit sauberem Inhalt ist die Menge in `pruefeMeta` für diese Zahlen am Code nicht zu isolieren (äquivalent im Meldungscode) und wird nur über den Text der Meldung rot; als Ganzes ist sie am Code rot (`TestLesenMeta` „Version vor Inhalt“, „Version vor fehlender Tabelle“, „Version vor unbekanntem Schlüssel“ für die Rangfolge Version vor Inhalt, und die Zahl über 64 Bit), und die Zusage bleibt die der Spezifikation (Code `PGR-E3002`, keine Version in der Meldung), nicht „die Menge liegt in `pruefeMeta`“. Der Test prüft den Text unabhängig vom Wortlaut: Die Meldung (ohne Pfad und ohne `PGR-E3002`) enthält die Zahl der Datei nicht als Teilstring, auch nicht ohne das Wort `Version`. Die Zahl über 64 Bit ist die Ausnahme: Hier hält allein `pruefeMeta` den Code `PGR-E3002`, ein Wegfall der Menge macht daraus `PGR-E3003`, und der Test ist am Code rot (`kopfVorpruefung` nimmt die Zahl nicht an, ein Satz „`Unmarshal` hält die Menge zusätzlich“ gilt für sie nicht).

**Nacharbeit zum Review (Findings F-1 bis F-3).** F-1: `pruefeEndeFaelle` prüft bei `PGR-E3002` nach Abzug von Pfad (`ohnePfad` aus `zeitangaben_pruefen_test.go`) und Code, dass die Zahl der Datei nicht als Teilstring in der Meldung steht; mit sauberem Inhalt wird die Menge in `pruefeMeta` für die kleinen Zahlen nur über die Meldung rot; als Ganzes ist sie am Code rot (Rangfolge Version vor Inhalt, Zahl über 64 Bit; Anmerkung zu den Zahlen). Der Fall *Version vor Inhalt* trägt nur die Zahl `3`. F-2: die Zeile zur Zahl über 64 Bit und die Anmerkung sind an den Code angeglichen (Textvergleich, keine Umwandlung). F-3: „Richtung `client` in Extended“ bleibt äquivalent (Zeile oben); die Sortierung der Interaktionen und Sessions hat einen Fall mit Ende (`TestLesenEndeNummernOrdnen`, `endeSQLGeordnet`; die Session-Zeilen stehen umgekehrt, weil `session.id` der Schlüssel der Tabelle ist, über `CREATE TABLE … ORDER BY id DESC`). Grenze: Beide Mutanten werden `PGR-E3003` über die Prüfung der Nummern („trägt die Nummer 2“, „trägt die Kennung 2“), bevor `ValidateEnd` läuft; der Fall belegt den Weg des Endes durch die Ordnung, keine eigene Zusage des Endes. Die Sortierung der Nachrichten hat ihren Fall mit Ende schon (Zeile oben).

Mutationen der Nacharbeit (Kopie des Baums, Bind-Mount, `--network=none`, Image `pgwire-recorder:test`, `go test -run 'TestLesenEnde|TestLesenMeta|TestLesenMeldung|TestLesenReihenfolge'`, jede einzeln, danach zurückgesetzt):

| Zusage | Mutation | roter Test |
|---|---|---|
| Meldung ohne die Zahl (F-1, Ausprägung 1) | die Menge in `pruefeMeta` entfällt (Meldung aus `Unmarshal`: `Version <Zahl>`) | `TestLesenEndeVersionen` (Meldung nennt die Zahl), `TestLesenMeta`, `TestLesenMeldungOhneWertUndZeile` |
| dasselbe (Ausprägung 2) | `pruefeMeta` hängt `version=<Zahl>` an | `TestLesenEndeVersionen` alle E3002-Fälle (alle acht), `TestLesenMeldungOhneWertUndZeile` |
| dasselbe (Ausprägung 3) | `pruefeMeta` meldet `Fassung <Zahl> wird nicht unterstützt` (der alte Test mit `Version <Zahl>` bliebe grün) | `TestLesenEndeVersionen` alle E3002-Fälle, `TestLesenMeldungOhneWertUndZeile` |
| dasselbe (Meldung in `kopfVorpruefung`) | Menge in `pruefeMeta` entfällt und die Meldung in `kopfVorpruefung` heißt `Fassung <Zahl>` | `TestLesenEndeVersionen`, `TestLesenMeta`, `TestLesenMeldungOhneWertUndZeile` |
| Meldung in `kopfVorpruefung` allein anders | nur `Fassung <Zahl>` in `kopfVorpruefung` | **grün, äquivalent:** `pruefeMeta` lehnt jede fremde Zahl vor `Unmarshal` ab, die Meldung von `kopfVorpruefung` ist für den SQLite-Leser nicht erreichbar |
| Interaktionen nach `sequence` | die Sortierung der Interaktionen entfällt | `TestLesenEndeNummernOrdnen` „Zeilen von interaction umgekehrt“, `TestLesenReihenfolgeDerZeilen/interaction` |
| Sessions nach `id` | die Sortierung der Sessions entfällt | `TestLesenEndeNummernOrdnen` „Zeilen von session umgekehrt“, `TestLesenReihenfolgeDerZeilen/session` |

**Befund vorher und nachher.** Vorher: `liesAufzeichnung` schrieb fest `Value: "1"`, `pruefeMeta` nahm nur `"1"` an (`grep -n 'Value: "1"\|version != "1"' internal/adapters/driven/recording/sqlite_lesen.go`: zwei Treffer). Nachher: kein Treffer; `grep -n 'strconv.Itoa(model.Version' internal/adapters/driven/recording/sqlite_lesen.go` nennt `model.Version` und `model.VersionMitEnde`.

**Größe.** Messung in der Metrik von §4 am Arbeitsbaum vor dem Commit (`git diff --numstat fa8486b -- internal docs/user`, die neue Datei mit `git add -N`): 330 (Einfügungen plus Löschungen: Produktcode 39, Tests 282, Tabelle 9, `test/integration` liegt außerhalb der Pfade der Metrik). Nach der Nacharbeit (`git diff --numstat fa8486b..HEAD -- internal docs/user` nach dem Commit): 387 (359 Einfügungen plus 28 Löschungen: Produktcode 39, Tests 338, Tabelle 10). Meldeschwelle 550 und Marke 700 nicht erreicht.

**Randformen.** Keine nicht genannte oder offene Randform gefunden; die Meldung zu einem Fehler von `Session.ValidateEnd` im SQLite-Leser nennt Session und nicht die verletzte Regel (`ladeInhalt` gibt `me.Msg` aus, Ursache bleibt in `Unwrap`); zugesagt ist dort nur der Code, der Test prüft nur ihn.

### Closure (Planner, 2026-10-15)

**Review und Verifikation.** Review (0 HIGH, 0 MEDIUM, 1 LOW F-1, 2 INFO F-2, F-3): alle bearbeitet in `654ad5a`. Verifikation: PASS, 0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO; 24 Mutanten nachgefahren, alle wie in §7 behauptet, beide äquivalenten Mutanten tatsächlich äquivalent. V-1 ist in §7 berichtigt (Anmerkung zu den Zahlen, Nacharbeit, Zeile zur Zahl über 64 Bit), V-3 im Kopf (`SPEC-041`, `make kopf-check` Exit 0). V-2 („in keiner Form“ weiter als die Prüfung der Ziffernfolge) und V-4 (`TestLesenEndeNummernOrdnen` belegt den Weg des Endes durch die Ordnung, keine eigene Zusage des Endes) bleiben als **Grenze** stehen; Code, Test-Kommentar und Abdeckungszeile sind nicht geändert.

**Gate-Nachweis.** Der volle `make gates` lief nach der Nacharbeit grün (Exit 0, §7 *Läufe*, ohne Flanke `p7-int-treiber`). Die Closure ändert nur Dateien auf der Positivliste (`docs/plan/planning/`).

**Gegenstand:** geliefert (kein Fall von „übernommen“ oder „entfallen“); der Schreiber ist vor dem Anspruch in `slice-v1-abschluss-verbindungsende-sqlite-schreiber` aufgegangen.

**Größe.** 387 in der Metrik von §4 (330 vor der Nacharbeit; Produktcode 39, Tests 338, Tabelle 10, Faktor Tests zu Produktcode 8,7), Schätzung des Architect 220 bis 420; Datenpunkt im Register bei `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Vermerk, kein Zählerzuwachs).

**Risiken aus §6 — Ausgang je Risiko** (geschlossene Menge):

| Risiko | Ausgang | Folge-Slice bzw. Begründung |
|---|---|---|
| Der Diff liegt über einer Review-Sitzung | entfallen | gemessen 387, 55 % der Marke 700, unter der Meldeschwelle 550; Review in einer Sitzung ohne Größen-Finding; Liefer-Punkte 1, Schichten zwei |

**Register** (Zähler = Dateien unter `evidence/`):

| Eintrag | Zähler | Fund | Stand |
|---|---|---|---|
| `BEO-REPO/abwesenheitszusage-an-wortlaut-der-meldung-geprueft` | 1 (neu) | F-1: Abwesenheitszusage über eine Meldung gegen einen Wortlaut geprüft | offen; nicht dieselbe Beobachtung wie `BEO-REPO/test-prueft-teilstring-gegen-meldung-mit-zufallspfad` (dort entscheidet ein Zufallspfad über einen Fehlalarm, hier ist der Test zu eng und lässt eine Verletzung durch), dessen Zähler deshalb bleibt |
| `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | 61 | F-2, V-1, V-2: Beleg und Anmerkung weiter als die Probe | verkörpert, Regel bleibt |
| `BEO-REPO/beleg-traegt-zusage-aus-anderem-grund` | 3 | kein Auftreten: V-4 ist eine im Plan benannte Grenze | verkörpert, Vermerk |
| `BEO-REPO/gruener-mutant-faelschlich-aequivalent` | 2 | kein Auftreten, Gegenbeleg: beide Äquivalenzen halten, begründet über eine Invariante der Struktur | offen, vermerkt |
| `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` | 5 | kein Auftreten: 387 in 220 bis 420, Faktor 8,7 | geplant (`slice-harness-schaetzung-planner-faktor`), Vermerk |

Kein Eintrag erreicht mit diesem Slice neu 3× ohne Regel.

**Lerneintrag — geschärfte Regel als Vorschlag, benannte Lücke im Beleg.** Geändert wird keine Regel, kein `liegt in`. (1) Eine Zusage der Form „die Meldung nennt X nicht“ wurde gegen einen Wortlaut geprüft, und Mutanten mit anderem Wortlaut blieben grün (F-1); der Test prüft jetzt die Form des Werts. (2) Ein Beleg in §7 beschrieb eine Mutation, wie sie der Verfasser aus dem Code ableitete, nicht wie er sie fuhr (F-2, V-1); gefahren, nicht abgeleitet, trug die Tabelle. Was funktionierte: Der Schnitt vor dem Anspruch hielt den Diff bei 55 % der Marke, die Randformen standen vor dem Code, und zwei äquivalente Mutanten wurden über eine Invariante der Struktur belegt und hielten der Nachfahrt stand. Vorschläge an den Nutzer, nicht von dieser Closure in `AGENTS.md` oder `.claude/` geändert: (a) `.claude/commands/implement-slice.md` Schritt 19: „Eine Abwesenheitszusage über eine Meldung prüft die Form des Werts, nicht den Wortlaut, und der Beleg fährt mindestens zwei andere Wortlaute“ (Wirkung: F-1 träte nicht auf; Gegenargument: erst 1× aufgetreten, Schwelle 3× nicht erreicht, Empfehlung daher: abwarten). (b) Derselbe Schritt: Eine Anmerkung in §7 über die Isolierbarkeit einer Mutation nennt den Test und den Grund, der gefahren wurde (Wirkung: V-1 träte nicht auf; die Klasse hat 61 Belege). Beides ist Empfehlung, keine Entscheidung.

**Paarungen.** (a) Anker: Der Lerneintrag trägt kein Feld `liegt in`, keine Verkörperung; nichts zu prüfen. (b) Folge-Slice: Genannt sind `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Bindung, §6) und `slice-harness-schaetzung-planner-faktor` (Datenpunkt); beide existieren in `next/`, und `grep -n "slice-v1-abschluss-verbindungsende-sqlite\b"` findet die Kennung des Gebers in `-sqlite-schreiber` (Herkunft Zeile 22, Zeile 47, Start Zeile 92); der Datenpunkt steht im Register. (c) Register: Alle genannten Kennungen `BEO-REPO/…` bestehen als Verzeichnis, jedes mit nicht leerem `evidence/`.

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

Nachgezählt bei der Prüfung des Architect vor dem Code (2026-10-15, `AGENTS.md` §3.13): Liefer-Punkte 1. Schichten zwei: der Recording-Adapter und `docs/user` (die Abdeckungstabellen zählen als Nutzer-Doku-Schicht, anders als die Zeile des Planners vom selben Tag „eine“); die Spezifikation und der Plan sind keine Schicht dieser Lieferung. Schätzung in der Metrik von §4 mit dem Faktor Tests zu Produktcode von mindestens 5: Produktcode 20 bis 40 (Version aus `meta` in den Baum, Vergleich mit zwei Konstanten, Kommentare), Tests 190 bis 370 (rund 60 Zeilen Hilfe zum Bauen der Dateien, rund 45 Fälle zu 2 bis 4 Zeilen, Tabellen und Aufrufe), Tabellen etwa 10: **220 bis 420, Mitte 320, 46 % der Marke 700**. Die Meldeschwelle 550 gilt unverändert.

Nachgezählt beim Umsetzen (Implementer, `AGENTS.md` §3.13): Liefer-Punkte 1. Schichten zwei: der Recording-Adapter und `docs/user` (Abdeckungstabelle); `test/integration` ändert einen Fall und ist keine zusätzliche Schicht der Lieferung. Gemessen 330 in der Metrik von §4 (Schätzung des Architect 220 bis 420), 47 % der Marke 700, unter der Meldeschwelle 550.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
