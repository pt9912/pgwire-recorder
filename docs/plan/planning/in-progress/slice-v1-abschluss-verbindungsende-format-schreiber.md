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

**Berührte Spec-Stellen:** `LH-FA-25.a` · `LH-FA-22.a` · `SPEC-001` · `SPEC-002` · `ARC-001` · `ARC-008`

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

**Ziel:** Der YAML-Schreiber schreibt `version: 2` genau dann, wenn die Aufzeichnung ein Ende enthält (`Recording.HasEnd()`), sonst `version: 1`, bei jedem Schreiben nach dem Inhalt und nicht nach dem Feld `version` des Modells (`LH-FA-25.a` *Version*, `SPEC-001`). Der Roundtrip (Laden, Schreiben, Laden: gleich) gilt für Aufzeichnungen mit Ende, einfach, Extended mit `Sync`-Gruppe und mit `Flush`-Gruppe, mit `offset_ms` und mit `empty_sessions`. Nachweis ohne `record`: Aufzeichnungen als Modell im Test (Text nur für die Darstellung und für den Fall ohne Ende, dessen Bytes bleiben), über den Leser von `slice-v1-abschluss-verbindungsende-format`; am Binary ändert sich nichts, bis `slice-v1-abschluss-verbindungsende-record` ein Ende aufzeichnet.

**Herkunft:** Teil 2 des Schnitts vom 2026-10-15 (Rückführung von `slice-v1-abschluss-verbindungsende-format`: Prüfung des Architect vor dem Code, ungeteilt 850 bis 1250 in der Metrik von §4, über der Marke 1000; Empfehlung Leser | Schreiber, Zuschnitt des Planners). Die Sendung „Schreiber und Roundtrip“ stammt aus §1 und DoD jenes Slice; er liefert Modell, `Session.ValidateEnd` und den Leser und liegt vor diesem in `done/` (Voraussetzung: Der Roundtrip lädt über seinen Leser). Schritt 35b der Reihenfolge in §5 von welle-v1-abschluss.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Modell (`connection_closed`, Formregeln, `Session.ValidateEnd`), die Konstante der Version 2 (`model.VersionMitEnde`; dieser Slice schreibt sie, definiert sie nicht neu) und der YAML-Leser — `slice-v1-abschluss-verbindungsende-format` (dort §1 und DoD): Vorgänger und Voraussetzung; mit ihnen läge dieser Plan wieder über der Marke.
- SQLite (Schreiber, Leser, `meta.version`, Roundtrip beider Formate) — `slice-v1-abschluss-verbindungsende-sqlite` (dort §1 und DoD): zweites Format mit eigenen Randformen (Transaktion); es benutzt den Schreiber dieses Slice für den Roundtrip in beiden Formaten und setzt ihn voraus.
- Aufzeichnen des Endes — `slice-v1-abschluss-verbindungsende-record` (Use Case, PGWire-Adapter) und `slice-v1-abschluss-verbindungsende-upstream` (Upstream-Adapter): Der `RecordService` setzt die Version nicht; sie folgt beim Schreiben dem Inhalt, und dieser Slice liefert diese Wahl. Der Slice `-record` rechnet mit Fakes und braucht den Schreiber nicht; `-upstream` zeichnet in YAML auf und setzt ihn voraus.
- Änderung am neutralen Schema `tools/schema/schema.yaml` — Entscheidung des Nutzers vom 2026-10-12 (Darstellung bestätigt): Das Ende ist ein Wert von `message.kind`, keine neue Spalte.
- Benutzerhandbuch und `README.md` — `slice-v1-abschluss-verbindungsende-doku` (dort §1 und DoD), nach dem letzten Code-Slice des Endes in derselben Welle (`AGENTS.md` §3.11, §3.13): Mit der Dokumentation als dritter Schicht läge dieser Plan über zwei Schichten, und ein Ende kann ein Anwender erst mit `record` erzeugen.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), `SPEC-001`: `Recording.HasEnd() bool` ist wahr, wenn irgendeine Interaktion irgendeiner Session ein Ende trägt. Der YAML-Schreiber schreibt `version: 2` mit mindestens einem Ende und sonst `version: 1`, bei jedem Schreiben nach dem Inhalt, nicht nach dem Feld `version` des Modells (Feld 1 mit Ende schreibt 2, Feld 2 ohne Ende schreibt 1); eine Aufzeichnung ohne Ende, auch eine ohne Sessions oder nur mit `empty_sessions`, schreibt `version: 1`; der Schreiber prüft die Form des Endes nicht, auch ein Ende an falscher Stelle schreibt 2 (Test je Stelle aus §6; `SPEC-001`; die vorhandenen Tests des Schreibers bleiben grün).
- [ ] [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats): Der Roundtrip (Laden, Schreiben, Laden: gleich) gilt für Aufzeichnungen mit Ende — einfach, Extended mit `Sync`-Gruppe und mit `Flush`-Gruppe, mit `offset_ms`, mit `empty_sessions` —; eine Aufzeichnung ohne Ende schreibt unverändert `version: 1` (Test, Aufzeichnungen als Modell im Test, der Leser lädt die geschriebene Datei; Text nur für die Darstellung und den Fall ohne Ende).
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
| `internal/hexagon/model` (`recording.go`, `hasend_test.go`) | update, neu | `Recording.HasEnd() bool` (Ort: Entscheidung des Architect in `slice-v1-abschluss-verbindungsende-format` §6); `hasend_test.go`: ein Fall je Ort und Session (`TestHasEnd`) |
| `internal/adapters/driven/recording/yaml.go` | update | Schreiber: `formatVersion` wählt `model.VersionMitEnde` mit Ende, sonst `model.Version`, nach dem Inhalt (`HasEnd`), nicht nach dem Feld `version`; `toDTO` ruft sie |
| `internal/adapters/driven/recording/yaml_schreiben_ende_test.go` | neu | Tests des Schreibers und des Roundtrips: `TestSchreibenVersionNachInhalt`, `TestSchreibenNimmtFormverletzungAn`, `TestSchreibenEndeDarstellung`, `TestRoundtripMitEnde`; Aufzeichnungen als Modell im Test, der Leser lädt die geschriebene Datei |
| `spec/spezifikation.md` | update (eng gefasst) | `SPEC-001` *Schreiber in `yaml`*: die Zusage „liest die Aufzeichnung einmal je Schreiben“ entfällt, kein Test prüft sie; keine Entscheidung des Architect ändert sich |
| `docs/user/abdeckung-unit.md` | update (`make abdeckung`) | Tabelle aus den Deklarationen der Tests; `abdeckung-e2e.md` bleibt (kein E2E-Test) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format` liegt in `done/` (Voraussetzung: Modell, Formregeln, Leser); dieser Slice ist Schritt 35b der Reihenfolge in §5 von welle-v1-abschluss und steht vor `slice-v1-abschluss-verbindungsende-sqlite` (dessen Roundtrip beider Formate benutzt den Schreiber). WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal docs/user`, Tests eingeschlossen, wobei die Basis `6c6dcd9` ist, der Stand des Hauptzweigs vor dem Anspruch (der Übergang nach `in-progress/` ändert nichts unter den genannten Pfaden); **über 700** geht der Slice zurück zur Zerlegung. **Meldeschwelle 550:** Erreicht die Summe 550, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Architect für diesen Teil 300 bis 500 (Produktcode etwa 40, Tests Schreiber und Roundtrip etwa 230, Hilfen und Tabellen, vor der Nacharbeit der Mutationen; Faktor 1,3 bis 1,5 gerechnet).
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
- *Aufrufer mit verletzter Form (Architect, 2026-10-15, vor dem Code):* Der Schreiber prüft die Form des Endes nicht — Stelle, Anzahl, `ready_for_query` daneben, Felder neben `type`; er schreibt, was die Aufzeichnung trägt, wie er heute jede andere Form ungeprüft schreibt. Gründe: Eine Prüfung allein für das Ende wäre eine Ausnahme im Schreiber und brächte einen neuen Fehlerweg in `Write` (Meldung, Katalog, Aufräumen der temporären Datei); die Form sichern der Aufrufer und der Leser. Folge: eine Aufzeichnung mit verletzter Form ergibt eine Datei, die der Leser mit `PGR-E3003` ablehnt (Test: schreiben nimmt an, laden lehnt ab). Das Gewicht liegt beim Aufrufer: `slice-v1-abschluss-verbindungsende-record` trägt die Bindung (dort §1 und §6, Eintrag im Commit dieser Prüfung) — `SPEC-001` *Schreiber in `yaml`*.
- *Stelle des Endes für die Wahl der Version (die Nachricht beendet die Folge, `AGENTS.md` §3.12):* `HasEnd` und die Version fragen nur, ob irgendwo `connection_closed` steht, nicht wo; die Stellen sind darum alle gleich entschieden — Version 2, kein Fehler im Schreiber: Anfang der Folge · Mitte · direkt vor dem Ende, also an der richtigen Stelle (letzte Nachricht) · letzte Gruppe · nicht letzte Gruppe · letzte Interaktion einer Session · nicht letzte Interaktion · nicht letzte Session. Die Zulässigkeit der Stelle entscheidet der Leser (`SPEC-002`, `PGR-E3003`); dort steht auch `ready_for_query` vor, nach oder statt dem Ende. Reihenfolge mehrerer Ablehnungsgründe: der Schreiber lehnt nichts ab, also gibt es keine; beim Laden gilt `SPEC-001` *Rangfolge der Prüfung*.
- *Modellfeld `version` außerhalb von 1 und 2* (0, 3, negativ): ändert nichts, der Schreiber schreibt 1 oder 2 nach dem Inhalt (`SPEC-001`); `Format` schreibt er unverändert.
- *`connection_closed` mit Feldern im Modell* (`Tag`, `Fields`, `Notice` …): der Schreiber schreibt sie, der Leser lehnt sie ab (`SPEC-001`, `SPEC-002`); er verwirft sie nicht still.
- *Darstellung des Endes:* Blockform, ein Listenelement mit dem einzigen Schlüssel `type`, in `responses` und in `server` einer Gruppe, eingerückt wie die übrigen Nachrichten (`SPEC-002`, `SPEC-004`); es gibt keine Flow-Form und keine Reihenfolge neben `type`, weil kein weiterer Schlüssel steht.
- *Version 1 byte-gleich:* eine Aufzeichnung ohne Ende ergibt dieselben Bytes wie vor der Änderung; Ausnahme ist nur das Modellfeld 2 ohne Ende (schreibt 1). Der Roundtrip ist gleich bis auf das Feld `version` einer Datei der Version 2 ohne Ende (`SPEC-001`).
- *Aufruf-Voraussetzung:* `HasEnd` und die Wahl der Version sind Zustandsprüfungen ohne Zusage über Reihenfolge oder Gleichzeitigkeit von Aufrufen; der Slice nennt keine Ordnung der Werte über mehrere Aufrufe, also gibt es keine Voraussetzung (`AGENTS.md` §3.12). Aufrufarten des Treibers: keine (kein Hook oder Gate liest ein Ergebnis dieses Schreibers). Bestätigt (Architect, 2026-10-15, auf Bitte des Verifiers `V-2`): Das Streichen von „liest die Aufzeichnung einmal je Schreiben“ in `SPEC-001` (`827763c`) ändert weder Entscheidung (b) noch die Aufruf-Voraussetzung „keine“; der Halbsatz war eine ungeprüfte Zusage (`AGENTS.md` §3.11), die Zusage über Reihenfolge und Gleichzeitigkeit steht unverändert.
- *Andere Schreiber:* der SQLite-Schreiber bleibt unberührt (`slice-v1-abschluss-verbindungsende-sqlite`); `model.Version` bleibt 1 für ihn.
- *Konstante der Version 2* (Bindung aus `slice-v1-abschluss-verbindungsende-format`, geliefert im Modell): Der Schreiber schreibt `version: 2` mit `model.VersionMitEnde` und `version: 1` mit `model.Version`; ein Literal ist es nicht — `SPEC-001`.

**Zählung der Mutationstabelle (Planner nach der Rechnung des Architect aus `slice-v1-abschluss-verbindungsende-format` §6, ergänzt vom Architect 2026-10-15; `AGENTS.md` §3.10):** ein Fall und eine Mutation, die nur ihn bricht, je Element von:

- der Wahl der Version: mit Ende 2 · ohne Ende 1 · Feld 1 mit Ende schreibt 2 · Feld 2 ohne Ende schreibt 1 · Ende nur in einer späteren Session · Ende nur in einer nicht letzten Session · ohne Sessions · nur `empty_sessions` · Feld `version` 0 (oder 3) · Ende in einer Extended-Gruppe · Ende in einer einfachen Interaktion (`HasEnd` sieht beide Orte; Mutation: nur `Responses` oder nur `Groups` prüfen);
- den Stellen nach dem Muster aus §6 (Anfang · Mitte · richtig · letzte Gruppe · nicht letzte Gruppe · nicht letzte Interaktion · nicht letzte Session): je Stelle Version 2 und kein Fehler beim Schreiben, Laden lehnt mit `PGR-E3003` ab (Ausnahme: die richtige Stelle lädt); Mutation je Stelle: der Schreiber lehnt sie ab oder schreibt 1;
- den Formen des Roundtrips: einfach · `Sync`-Gruppe · `Flush`-Gruppe · `offset_ms` · `empty_sessions` · Version 2 ohne Ende (gleich bis auf `version`);
- der Darstellung: ein Text mit Bytes für ein Ende in `responses` und eines in `server` (Blockform, Einrückung) · eine Aufzeichnung ohne Ende, deren Text Byte für Byte erhalten bleibt (Laden, Schreiben) · `connection_closed` mit Feld im Modell (Feld wird geschrieben, Laden lehnt ab).

Der Schreiber prüft nicht (oben); ein Test, der das sichert, ist der Fall „Stelle falsch: schreiben nimmt an“ je Stelle — Mutation: Prüfung im Schreiber.

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

### Belege des Implementers

**Läufe:** `make gates` grün (Exit 0, ohne Flanke in der Versionsmatrix-Gegenprobe); davor `make lint` einzeln (erster Lauf rot: `gochecknoglobals` für vier Hilfswerte im Test, danach als Funktionen). Stand der Nacharbeit nach dem Review (`F-1` bis `F-7`): Läufe und Mutationen unten mit „Nacharbeit“ markiert. Der Stand der Mutationen: Bind-Mount des Arbeitsbaums in das Image `pgwire-recorder:test` mit `--network=none`, `go test -run 'Schreiben|Roundtrip|HasEnd'` über `./internal/adapters/driven/recording/` und `./internal/hexagon/model/`; nach jedem Lauf die Originaldateien zurückkopiert (kein Build-Kontext, also kein mtime-Problem). Stellen und Darstellung als Mutanten über eine Hilfsdatei im Paket, gesteuert über eine Umgebungsvariable; sie ist nicht Teil des Diffs.

**Metrik** (`AGENTS.md` §3.14): Summe aus Einfügungen und Löschungen von `git diff --numstat 6c6dcd9..HEAD -- internal docs/user`: 311 vor der Nacharbeit, 320 nach ihr (Einfügungen 319, Löschungen 1; Produktcode 26, Tests 287, Tabelle 7); Meldeschwelle 550 und Marke 700 nicht erreicht.

**Vorhandene Tests:** kein vorhandener Test mit `Recording{}` oder Version 0 ändert sich; alle Tests des Repos bleiben ohne Änderung grün.

**Zusage · Mutation · roter Test** (jede Zeile: der Mutant färbt genau die genannten Fälle rot, wo nichts anderes steht):

| Zusage | Mutation | roter Test |
|---|---|---|
| `HasEnd`: Ende in `Responses` | nur `Groups` prüfen | `TestHasEnd/Ende_in_Responses` und vier weitere Fälle mit Ende in `Responses` |
| `HasEnd`: Ende in `Server` einer Gruppe | nur `Responses` prüfen | `TestHasEnd/Ende_in_Server_einer_Gruppe`, `…_einer_nicht_letzten_Gruppe` |
| `HasEnd`: an jeder Stelle (nicht letzte Nachricht) | nur die letzte Antwort einer Interaktion prüfen / nur die letzte Server-Nachricht einer Gruppe prüfen | `TestHasEnd/Ende_in_Responses,_nicht_letzte_Nachricht` / `…_in_Server_einer_Gruppe,_nicht_letzte_Nachricht` |
| `HasEnd`: auch in einer nicht letzten Gruppe | nur die letzte Gruppe prüfen | `TestHasEnd/Ende_in_Server_einer_nicht_letzten_Gruppe` |
| `HasEnd`: jede Interaktion | nur die erste / nur die letzte Interaktion | `TestHasEnd/Ende_in_der_letzten_Interaktion` / `…_in_einer_nicht_letzten_Interaktion` |
| `HasEnd`: jede Session | nur die erste / nur die letzte Session | `TestHasEnd/Ende_nur_in_der_letzten_Session`, `…_in_der_mittleren_von_drei_Sessions` / `…_nur_in_der_ersten_Session` |
| `HasEnd` falsch ohne Sessions / ohne Interaktion | wahr bei `len(Sessions)==0` / wahr bei einer Session ohne Interaktion | `TestHasEnd/ohne_Sessions` / `…Session_ohne_Interaktion` |
| Version 2 mit Ende | immer `model.Version` | alle Fälle mit Ende in `TestSchreibenVersionNachInhalt`, `TestRoundtripMitEnde` u. a. |
| Version 1 ohne Ende | immer `model.VersionMitEnde` | `…/Feld_0,_1,_2,_3_ohne_Ende`, `…/ohne_Sessions`, `…/nur_empty_sessions`, `TestSchreibenEndeDarstellung/ohne_Ende,_Bytes_bleiben`, `TestRoundtrip` (vorhanden) |
| Feld 1 mit Ende schreibt 2 (Nacharbeit `F-1`) | Feld 1 mit Ende schreibt 1 (`rec.HasEnd() && rec.Version != 1`) | genau `TestSchreibenVersionNachInhalt/Feld_1_mit_Ende`; alle anderen Aufzeichnungen mit Ende tragen jetzt das Feld 2 (0 und 3 in ihren eigenen Fällen), nur dieser Fall trägt 1 |
| Feld 2 ohne Ende schreibt 1 | Feld 2 schreibt 2 | `TestSchreibenVersionNachInhalt/Feld_2_ohne_Ende`, `TestRoundtripMitEnde/Version_2_ohne_Ende` (nur diese) |
| Feld negativ, mit und ohne Ende (Nacharbeit `F-5`) | negatives Feld mit Ende schreibt 1 (`rec.HasEnd() && rec.Version >= 0`) · negatives Feld ohne Ende schreibt 2 (`rec.HasEnd() \|\| rec.Version < 0`) | je genau `…/negatives_Feld_mit_Ende`, `…/negatives_Feld_ohne_Ende` |
| Feld 0 / Feld 3, je mit und ohne Ende (in der Nacharbeit je einzeln neu gefahren: `!= 0`, `!= 3`, `\|\| == 0`, `\|\| == 3`) | Feld 0 mit Ende schreibt 1 · Feld 0 ohne Ende schreibt 2 · Feld 3 mit Ende schreibt 1 · Feld 3 ohne Ende schreibt 2 | je genau `…/Feld_0_mit_Ende`, `…/Feld_0_ohne_Ende`, `…/Feld_3_mit_Ende`, `…/Feld_3_ohne_Ende` |
| Feld des Modells bestimmt nichts | Feld 0 bis 3 übernehmen | alle Fälle mit Ende, `Feld_0_ohne_Ende`, `Feld_2_ohne_Ende`, `Feld_3_ohne_Ende` u. a. |
| Ohne Sessions 1 / nur `empty_sessions` 1 | `len(Sessions)==0` schreibt 2 / `EmptySessions` schreibt 2 | `TestSchreibenVersionNachInhalt/ohne_Sessions` / `…/nur_empty_sessions` (je dazu der Leser-Roundtrip der Bestandstests) |
| Konstante `model.VersionMitEnde` | `VersionMitEnde = 3` | alle Fälle mit Ende (der Schreiber folgt der Konstante) |
| Schreiber lehnt nichts ab, je Stelle: Anfang · Mitte · nicht letzte Interaktion · nicht letzte Gruppe · mitten in der letzten Gruppe · zwei Enden · `ready_for_query` daneben · Feld neben `type` · nicht letzte Session | der Schreiber lehnt die Stelle ab (je ein Mutant) | `TestSchreibenNimmtFormverletzungAn/` mit dem Namen der Stelle; „nicht letzte Interaktion“ färbt zusätzlich „zwei Enden“ (zwei Enden in einer Session heißt ein Ende in einer nicht letzten Interaktion) |
| Das Ende an der richtigen Stelle lädt | der Schreiber lehnt die richtige Stelle ab (einfach / letzte Gruppe) | `…/richtige_Stelle,_einfach` u. a. / `…/richtige_Stelle,_letzte_Gruppe` u. a. |
| Laden lehnt die Formverletzung ab | (im Leser, `slice-v1-abschluss-verbindungsende-format`) | `TestSchreibenNimmtFormverletzungAn` erwartet `PGR-E3003` je Stelle |
| Roundtrip einfach · Sync-Gruppe · Flush-Gruppe · `offset_ms` · `empty_sessions` | je ein Mutant, der das Ende (einfach, Sync, Flush), `offset_ms` oder `empty_sessions` nur der Aufzeichnung mit dem Kennzeichen der Form verliert | je genau `TestRoundtripMitEnde/<Form>` |
| Laden, Schreiben, erneutes Laden (Nacharbeit `F-3`): von der geladenen Aufzeichnung aus schreiben und laden ist dieselbe Aufzeichnung | kein Mutant des Schreibers, den nur dieser zweite Durchgang färbt: der Schreiber ist eine Funktion des Modells, und die geladene Aufzeichnung unterscheidet sich von der gebauten nur im Feld `version`, das jeder Mutant schon im ersten Durchgang färbt (Grenze, kein Mutant zugesagt) | `TestRoundtripMitEnde/<Form>`, zweiter Durchgang (grün, belegt den Pfad; rot würde ihn ein Schreiber, der eine geladene Aufzeichnung anders schriebe als eine gebaute) |
| Roundtrip Version 2 ohne Ende (gleich bis auf `version`) | Feld 2 schreibt 2 | `TestRoundtripMitEnde/Version_2_ohne_Ende` |
| Darstellung: Blockform in `responses` / in `server` | Flow-Form `{type: connection_closed}` an der einen / der anderen Stelle | `TestSchreibenEndeDarstellung/in_responses` (und `mit_Feld_neben_type`) / `…/in_server_einer_Gruppe` |
| Darstellung: Einrückung | eine Stufe tiefer in `responses` / in `server` | `…/in_responses` u. a. / `…/in_server_einer_Gruppe` u. a. |
| `connection_closed` mit Feld schreibt das Feld | das Feld verwerfen | `TestSchreibenEndeDarstellung/mit_Feld_neben_type`, `TestSchreibenNimmtFormverletzungAn/Feld_neben_type` |
| Aufzeichnung ohne Ende, Bytes bleiben | ein Zeilenumbruch mehr am Ende ohne Ende | `TestSchreibenEndeDarstellung/ohne_Ende,_Bytes_bleiben` (nur dieser) |

**Grüner Mutant:** `return 2` statt `model.VersionMitEnde` im Schreiber bleibt grün, weil die Konstante 2 ist (äquivalent: das Verhalten ist gleich). „Kein Literal“ ist eine Zusage über die Schreibweise; das Verhalten prüft `TestVersionMitEnde` des Modells (`VersionMitEnde == 2`) und die Mutation `VersionMitEnde = 3` oben. Sonst blieb kein Mutant grün.

**Nacharbeit (Review `F-1` bis `F-7`):** `F-1` isoliert (oben, Mutation rot gesehen, nur der eine Fall); `F-2` §1 und DoD Punkt 2 auf „Modell im Test“ gezogen; `F-3` Test ergänzt (zweiter Durchgang, oben) und die Deklaration auf „Schreiben und Laden einer Aufzeichnung und, von der geladenen aus, Schreiben und erneutes Laden“ gefasst; `F-7` in `SPEC-001` nur gestrichen („liest die Aufzeichnung einmal je Schreiben“), die Zusage über Reihenfolge und Gleichzeitigkeit bleibt, keine Entscheidung des Architect (`4940cce`) ändert sich; `F-5` Fälle „negatives Feld mit/ohne Ende“ mit Mutationen (oben, je genau ein Fall rot); `F-6` §8 angepasst; `F-4` kein Handeln, die Äquivalenz von `return 2` bleibt oben dokumentiert. Klasse „Element der Aufzählung ohne isolierenden Fall“: die Fixtures mit Ende tragen nicht mehr das Feld der Aufzählung; Mutationen der Aufzählung (1 mit Ende, 2 ohne Ende, 0 und 3 je mit und ohne Ende, negativ je mit und ohne Ende) einzeln gefahren, je genau ein Fall rot (Ausnahme `Feld_2_ohne_Ende`, dazu `Roundtrip/Version_2_ohne_Ende`, der dasselbe Element ist). Lauf der Nacharbeit: `make gates` grün (Exit 0, ohne Flanke), Stand: Arbeitsbaum dieses Commits.

**Randformen:** keine zusätzliche gefunden; §6 und `SPEC-001` entscheiden alle Fälle, die der Code berührt. Keine Rückgabe an den Architect.

**Beobachtungen:** aus der Umsetzung keine eigene; die Eintragung ins Register entscheidet die Closure.


## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand). Treffer in der Sub-Area `REPO`: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (beide `verkörpert`, `AGENTS.md` §3.14; mit der Rückführung von `slice-v1-abschluss-verbindungsende-format` ein Auftreten trotz der Regel, Urteil im Register) — dieser Plan nennt Metrik, Meldeschwelle und Marke in §4 und liegt nach der Schätzung des Architect unter der Marke. `BEO-REPO/meldeschwelle-ohne-artefakt` (1×): §4 nennt die Meldeschwelle. `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` ist verkörpert (`AGENTS.md` §3.10): die Mutationstabelle steht in DoD und §7. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-15, `AGENTS.md` §3.13): Liefer-Punkte: 2 (Schreiber mit `HasEnd` und Version nach dem Inhalt, Roundtrip). Schichten: zwei (Kern mit dem Modell, `ARC-001`; Recording-Adapter, `ARC-008`); der Handbuch- und README-Teil liegt in `slice-v1-abschluss-verbindungsende-doku`. Schätzung des Architect (2026-10-15, aus dem ungeteilten Plan): 300 bis 500 Zeilen mit Tests in der Metrik von §4; nach der Prüfung vor dem Code 340 bis 500 (Produktcode etwa 40, Tabellen und Fälle der Mutationen etwa das Siebenfache, vorhandene Hilfen in `yaml_ende_test.go` genutzt), Mitte 420 und damit 60 % der Marke, oberes Ende 71 %; kein Schnitt, die Meldeschwelle 550 trägt (Schnitt, falls sie erreicht wird: Version und Stellen | Roundtrip und Darstellung). Bindung aus `slice-v1-abschluss-verbindungsende-format` (Closure des Gebers, 2026-10-15): die Namen `model.VersionMitEnde` und `model.Version` für die Wahl der Version; sie gehört zu DoD-Punkt 1, ein Aufruf vorhandener Konstanten, ohne Liefer-Punkt und Schicht (Liefer-Punkte weiter 2, Schichten weiter zwei). Sendungen anderer Slices an diesen: `slice-v1-abschluss-verbindungsende-sqlite` und `-upstream` nennen ihn als Voraussetzung (dort §1, Start und §8 nachgezählt).

Messung (Implementer, nach dem Code und nach der Nacharbeit zum Review): Summe aus Einfügungen und Löschungen von `git diff --numstat 6c6dcd9..HEAD -- internal docs/user` 311, nach der Nacharbeit 320 (Schätzung 340 bis 500, also darunter; Meldeschwelle 550 nicht erreicht). Liefer-Punkte und Schichten wie beim Anlegen: 2 und zwei. Nehmer: Der Architect trug die Bindung „Schreiber prüft die Form nicht“ in `slice-v1-abschluss-verbindungsende-record` ein (dort §1, §6 und §8, nachgezählt: Liefer-Punkte 3, Schichten zwei, Commit `4940cce`); `slice-v1-abschluss-verbindungsende-sqlite` und `-upstream` tragen sie nicht, sie liegt dort indirekt über den Use Case aus `-record`. Der Implementer trug in diesem Lauf nichts in einen Nehmer ein.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).
