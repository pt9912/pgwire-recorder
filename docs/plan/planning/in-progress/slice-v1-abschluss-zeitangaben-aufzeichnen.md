# Slice slice-v1-abschluss-zeitangaben-aufzeichnen: Zeitangaben aufzeichnen (`--record-timing`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen)

**Berührte Spec-Stellen:** `LH-FA-21.a` · `LH-FA-17.a` · `SPEC-001` · `SPEC-002` · `SPEC-004` · `SPEC-041` · `SPEC-043` · `SPEC-034` · `ARC-001` · `ARC-002`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-11.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `record --record-timing` trägt an jeder Interaktion `offset_ms` ein — ganze Millisekunden zwischen dem Beginn der Session (`LH-FA-21.a` *Messpunkt*: der Aufruf von `OpenSession`, nach TLS und Startup des Clients) und dem Eintreffen der ersten Client-Nachricht der Interaktion, mit der monotonen Uhr des Ports —, in YAML und in SQLite; ohne die Option fehlt das Feld und die Ausgabe ist unabhängig vom Zeitverhalten. Die Aufzeichnung lässt sich in beiden Formaten laden (Roundtrip).

**Übernommen aus `slice-v1-abschluss-zeitangaben`** (dort bis zum Schnitt DoD-Punkt 1 und §3, Zeilen *PGWire-Adapter*, *Recording-Adapter* und *CLI-Adapter*, Option `--record-timing`): das Aufzeichnen der Zeitangaben.

**Übernommen aus `slice-v1-abschluss-sqlite-format-lesen`** (dort §1, Abgrenzung): die Zeitangaben in der SQLite-Datei. Das Lesen der Spalte `interaction.offset_ms` liefert jener Slice (nur Spalte mit geprüftem Typ: negative und nicht ganze Werte sind `PGR-E3003`); dieser Slice schreibt sie mit `--record-timing` und belegt den Roundtrip in beiden Formaten.

**Übernommen aus `slice-v1-abschluss-einspielen`** (dort §1, Abgrenzung): Zeitangaben; hier die Hälfte beim Aufzeichnen, die Hälfte beim Einspielen liegt in `slice-v1-abschluss-zeitangaben`.

**Übernommen aus `slice-v1-abschluss-zeitangaben-uhr`** (dort §3, Zeile `internal/bootstrap/systemuhr.go`, und §6 Risiko 3): die Verdrahtung der Systemuhr. Der Bootstrap ruft `NeueSystemuhr()` und reicht das Ergebnis an den Record-Service; bis zu diesem Slice ist die Systemuhr nirgends verdrahtet. Die Logik der Systemuhr (Warten mit Abbruch) liefert `-uhr` und zählt dort.

**Herkunft:** Teil des Schnitts von `slice-v1-abschluss-zeitangaben` vom 2026-10-11 (Entscheidung des Nutzers vom 2026-10-11: vor dem Anspruch zu schneiden, weil der Plan vier Schichten berührte — Kern mit Uhr-Port, PGWire-Adapter, Recording-Adapter, CLI-Adapter — und damit über der Grenze von zwei lag, `AGENTS.md` §3.13; den Zuschnitt in vier Teile nach Lieferwert entschied der Planner mit seiner Empfehlung). Die vier Teile: `slice-v1-abschluss-zeitangaben-pruefen` (Lesen prüft die Werte), `slice-v1-abschluss-zeitangaben-uhr` (Uhr-Port und Systemuhr), `slice-v1-abschluss-zeitangaben-aufzeichnen` (`--record-timing`) und `slice-v1-abschluss-zeitangaben` (`--keep-timing`, schließt Abnahmeszenario 13); Handbuch und README zu allen vier liegen gemeinsam im Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` hinter dem letzten Teil (wie `slice-v1-abschluss-sqlite-format-doku` hinter dem Schreiben und dem Lesen: Zwischen den Code-Slices und dem Doku-Slice hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit allen).


**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Prüfung absteigender Werte beim Lesen — `slice-v1-abschluss-zeitangaben-pruefen` (dort §1 und DoD, mit der Kennung dieses Slice): Dieser Slice schreibt aufsteigende Werte und lädt sie über die geprüften Leser; die Prüfung selbst liegt dort und steht in der Reihenfolge vor diesem Slice.
- Der Uhr-Port, die Fake-Uhr und die Systemuhr — `slice-v1-abschluss-zeitangaben-uhr` (dort §1 und DoD, [ADR-0015](../../adr/0015-uhr-port.md)): Dieser Slice benutzt den Port, wie jener ihn liefert. Der Port ist `Uhr` in `internal/hexagon/ports/driven/uhr.go` mit `Jetzt() int64` (ganze Millisekunden seit einem Ursprung, den die Uhr wählt, nie kleiner als die vorige Lesung) und `WarteBis(ctx, zeitpunkt int64) error`; in Tests die Fake-Uhr `testuhr.Neu` mit `Vorruecken` und `Wartende`, im Betrieb die Systemuhr `NeueSystemuhr` aus `internal/bootstrap`. Dieser Slice ruft nur `Jetzt` (Messen wartet nicht). Ändert er eine Signatur, ist das ein Plan-Wechsel jenes Slice, nicht eine Ergänzung.
- `--keep-timing`, `--timing-mode`, `--timing-reference` und das Warten beim Einspielen — `slice-v1-abschluss-zeitangaben` (dort §1 und DoD): Abnahmeszenario 13 läuft erst, wenn beide Hälften da sind; es schließt in jenem Slice, der hinter diesem steht.
- Eine neue Version des Recording-Formats — jedes neue Feld, auch ein optionales, erhöht `version` (`SPEC-001`); `offset_ms` steht schon in Version 1 (`SPEC-002`, `SPEC-043`), dieser Slice fügt kein Feld hinzu (Rückführung in §4).
- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` (dort §1 und DoD), hinter dem letzten der vier Teile in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über der Grenze der Schichten.
- Änderungen am PGWire-Adapter — Schicht-Abgrenzung: Der Use Case misst, nicht der Adapter (Entscheidung des Architect vom 2026-10-11, Randform *Messpunkte* in §6); der Port `Recorder` ändert seine Signatur nicht. Zeigt sich im Code, dass der Adapter doch die Zeit mitgeben muss, ist das die dritte Schicht und die Rückführung in §4.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen): `record --record-timing` trägt `offset_ms` an jeder Interaktion ein, einfach wie Extended, gemessen mit `Jetzt` der Uhr des Ports (`driven.Uhr`, Systemuhr `NeueSystemuhr` vom Bootstrap gereicht); die Aufzeichnung lädt in YAML und in SQLite per Roundtrip mit den Werten (E2E am Binary mit beiden Formaten, darin eine Pause des Clients zwischen zwei Interaktionen, deren Untergrenze der Unterschied der `offset_ms` einhält — ohne verdrahtete Systemuhr wäre er 0 —, und Test mit Fake-Uhr `testuhr.Neu` und `Vorruecken` im Service); ohne die Option fehlt das Feld und die Ausgabe ist unabhängig vom Zeitverhalten des Laufs, in YAML byte-gleich zwischen zwei Läufen mit verschiedenen Pausen (Test, Boundary). Belege: `TestRecordZeitEinfach`, `TestRecordZeitNullSteht`, `TestRecordZeitNichtAbsteigend`, `TestRecordZeitExtendedErsteNachricht` und `TestRecordZeitJedeSessionEigenerBeginn` im Service mit Fake-Uhr (`internal/hexagon/services/record_zeit_test.go`), `TestRecordZeitRoundtrip` mit den echten Schreibern (`internal/bootstrap/record_zeit_test.go`), `TestE2ERecordTiming` am Binary mit der Pause des Clients, `TestE2ERecordOhneTiming` ohne die Option (`test/integration/zeitangaben_aufzeichnen_e2e_test.go`); jeder Test ist in der Tabelle von §7 mit seiner Mutation rot gesehen und von der Verifikation neu gebaut.
- [x] Option `--record-timing` für `record` mit Umgebungsvariable und Schlüssel der Konfigurationsdatei nach den Regeln der Optionstabelle (`LH-FA-17.a`, Default `false`, Priorität Option vor Variable vor Datei, ungültiger Wert `PGR-E2001` bzw. `PGR-E2004`); bei `replay` und `play` unbekannt (`PGR-E2001`); `config show` und `--help` nennen sie (Test am Binary). Belege: `TestRecordTimingQuellen`, `TestRecordTimingHilfe`, `TestLeserAlleOptionen`, `TestLeserOptionen`, `TestLeserFremdeUmgebung`, `TestParseRecord`, `TestDateiUngueltig`, `TestConfigShowFehler` (`internal/adapters/driving/cli`), am Binary `TestE2ERecordTimingQuellen`, `TestE2ERecordTimingUngueltig`, `TestE2ERecordTimingHilfeUndConfigShow`.
- [x] `make gates` grün (Planner, sauberer Baum, Stand `3727254`; nach dem letzten Commit der Closure bestätigt der Bericht den Lauf erneut).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Belege: §7, `evidence/slice-v1-abschluss-zeitangaben-aufzeichnen.md` in fünf Einträgen.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services` (Record-Service, `record.go`) | update | `offset_ms` je Interaktion mit der Uhr des Ports, eingestellt über die Option `services.Zeitangaben(uhr, aufzeichnen)` von `NewRecordService`; Basis je Session = Lesung beim Aufruf von `OpenSession`, Ende = Lesung beim Übernehmen der ersten Nachricht der Interaktion (`Query`, erste `ClientMessage` nach dem Ende der vorigen Interaktion), vor jedem Warten; der Port `Recorder` ändert seine Signatur nicht, sein Kommentar nennt die Voraussetzung nacheinander kommender Aufrufe der Client-Richtung (`internal/hexagon/ports/driving/record.go`, Entscheidung des Architect zu F-647), ebenso die Kommentare von `Query` und `ClientMessage` |
| `internal/hexagon/model` (`recording.go`) | update | Kommentar von `OffsetMS` nennt „Beginn der Session“ in der Lesart von `LH-FA-21.a` *Messpunkt*; kein neues Feld |
| `internal/adapters/driving/cli` | update | Option `--record-timing`, Umgebungsvariable, Schlüssel im Abschnitt `record:`, `config show`, Hilfe |
| `internal/bootstrap` | update | `NeueSystemuhr()` aufrufen und die Uhr mit `o.RecordTiming` an den Record-Service reichen, `services.Zeitangaben(NeueSystemuhr(), o.RecordTiming)` ohne Verzweigung (reine Verdrahtung, keine Schicht, `AGENTS.md` §3.13; die Logik der Systemuhr liegt in `slice-v1-abschluss-zeitangaben-uhr`) |
| Tests, `test/integration` | neu/update | Service mit Fake-Uhr (`internal/hexagon/services/record_zeit_test.go`), Roundtrip beider Formate mit den echten Schreibern (`internal/bootstrap/record_zeit_test.go`) und am Binary, Optionstabelle; die Testhilfe `internal/bootstrap/offsetlese` lädt die Aufzeichnung für die Integrationstests (unter `test/` darf der Recording-Adapter nicht importiert werden), `TestBinaryOhneTesthilfen` hält sie aus dem Binary; `docs/user/abdeckung-*.md` von `make abdeckung` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-zeitangaben-pruefen` und `slice-v1-abschluss-zeitangaben-uhr` liegen in `done/`, dazu `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` (der Roundtrip in SQLite; Schritt 21 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md); WIP-Limit 1). Vor dem ersten Code-Commit entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Messpunkt liegt im PGWire-Adapter — das wären Kern, PGWire-Adapter und CLI-Adapter, drei Schichten —, oder die Verdrahtung der Systemuhr braucht im Bootstrap Logik (mehr als `NeueSystemuhr()` aufrufen und das Ergebnis reichen: eine Datei lesen, einen Fehler einstufen, verzweigen) und der Bootstrap zählt als dritte Schicht (`AGENTS.md` §3.13), oder der Diff ist nicht in einer Review-Sitzung prüfbar; Schnitt dann: Messpunkt im Adapter mit der Signatur des Ports als eigener Slice vor der Option.
- `in-progress` → `next`: Das Zeitmodell verlangt eine Änderung am Recording-Format über das in Version 1 enthaltene Feld `offset_ms` hinaus; jedes neue Feld, auch ein optionales, erhöht `version` (`SPEC-001`) — zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Die Fake-Uhr bildet die monotone Uhr nicht ab oder der Port aus `slice-v1-abschluss-zeitangaben-uhr` trägt den Messpunkt nicht — Carveout, dann zuerst der Plan jenes Slice.

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

**Randformen** (`AGENTS.md` §3.12) — vom Architect vor dem Code am 2026-10-11 entschieden, Ort: [`LH-FA-21.a`](../../../../spec/spezifikation.md#lh-fa-21a--zeitangaben) *Aufzeichnen*, *Uhr*, *Messpunkt*. Entscheidungen des Architect (Empfehlung, bis der Nutzer anderes sagt; der Nutzer hat am 2026-10-11 verfügt, ohne Nachfrage weiterzumachen).

- *Messpunkte:* Der Use Case misst. Beginn der Session ist der Aufruf von `OpenSession` (vor dem Aufbau der Upstream-Session; TLS und Startup des Clients davor zählen nicht), Ankunft der Aufruf, mit dem der Use Case die erste Nachricht der Interaktion übernimmt, gelesen vor jedem Warten — entschieden ebenda. Ein Client, der Anfragen vorausschickt, wird um die Wartezeit auf die vorige Antwort später gemessen (Grenze, kein Fehler; der Wert ist nie kleiner als die Zeit bis zur Ankunft am Socket). Der PGWire-Adapter bleibt unberührt; der Slice hat zwei Schichten.
- *Uhr des Ports:* Eine Uhr je Lauf von `record`, damit Beginn und Ankunft von derselben Uhr stammen; ein Wert hat nur als Unterschied zweier Lesungen derselben Uhr eine Bedeutung, und keine Lesung selbst wird aufgezeichnet — entschieden: `LH-FA-21.a` *Uhr* ([ADR-0015](../../adr/0015-uhr-port.md)). Dieser Slice ruft `WarteBis` nicht auf; die Zusagen des Wartens (Fehler des Kontexts, ferner Zeitpunkt) gelten ihm nicht.
- *Rundung:* Es wird nicht gerundet: Der Wert ist die Differenz zweier abgeschnittener Lesungen der Uhr (ganze Millisekunden) — entschieden ebenda. Die Empfehlung des Planners (abschneiden) gilt für die Uhr, nicht für die Differenz.
- *Erste Interaktion:* kann 0 sein; `offset_ms` ist nie negativ — entschieden (ganze Zahl ≥ 0, die Uhr ist monoton).
- *Nicht absteigend:* Die Interaktionen werden in der Reihenfolge ihrer Ankunft nummeriert, die Lesungen werden nicht kleiner, also steigen die Werte nicht ab — entschieden ebenda. Ergibt ein Test mit der Fake-Uhr ein Gegenbeispiel (etwa zwei Gruppen einer Extended-Session, deren Nummer in anderer Reihenfolge fällt als ihre erste Nachricht ankommt), hält der Implementer an und gibt es zurück. Entschieden am 2026-10-12 (F-647, Empfehlung des Architect, der Nutzer hat verfügt, ohne Nachfrage weiterzumachen; Ort: `LH-FA-21.a` *Messpunkt*, Sicht: Zeile *Record-Use-Case*): Die Zusage gilt unter der Bedingung, dass die Aufrufe der Client-Richtung einer Session (`Query`, `ClientMessage`) nacheinander kommen; der Use Case ordnet gleichzeitige Aufrufe nicht um und sagt für sie keine Ordnung zu. Gründe gegen die Absicherung im Use Case: Die Nummer vergibt `Query` erst nach der Antwort des Upstream, die Lesung liegt bei der Ankunft; beides unter einer Sperre zu halten hieße, die Client-Richtung für die Dauer des Upstream-Aufrufs zu sperren oder nach Ankunft umzunummerieren — ein Umbau des Nebenläufigkeitsmodells für einen Fall, den kein Aufrufer erzeugt. *Akzeptiertes Negativ:* Dass der PGWire-Adapter die Bedingung einhält (eine Schleife, eine wartende `Query` hält das Lesen an), prüft kein Test; ein Adapter-Test wäre eine dritte Schicht dieses Slice (`AGENTS.md` §3.13), und die Eigenschaft ändert sich nur mit einem Umbau des Adapters, den dessen Review sieht. Die Spezifikation sagt es nicht als Zusage des Adapters (`AGENTS.md` §3.11). Der Roundtrip lädt über die geprüften Leser.
- *Extended-Interaktion:* das Feld steht auf der Ebene der Interaktion; Messpunkt ist die erste Client-Nachricht der ersten Gruppe, gleich welche Art — entschieden ebenda. Überlappende Gruppen gibt es nicht: Eine Interaktion endet mit `Sync`, die nächste beginnt danach.
- *Lebendprüfung, Fehlerantwort, Abbruch und Herunterfahren:* Die Lebendprüfung trägt das Feld wie jede Anfrage; eine Interaktion, die unvollständig endet, wird nicht übernommen und hat keinen Wert; jede aufgezeichnete Interaktion trägt mit der Option das Feld — entschieden ebenda.
- *Mehrere Sessions und `--record-empty-sessions`:* Basis je Session ihr eigener Beginn; eine Session ohne Interaktion trägt kein Feld — entschieden.
- *Ohne die Option:* Das Feld fehlt, und die Ausgabe ist unabhängig vom Zeitverhalten des Laufs, ausgenommen `Flush`-Gruppen ohne wartenden Client (`SPEC-004`) — entschieden.
- *Option:* Name, Umgebungsvariable `PGWIRE_RECORDER_RECORD_TIMING`, Schlüssel `record_timing` im Abschnitt `record:`, Default und Wertemenge stehen in der Optionstabelle und in `LH-FA-17.a` — entschieden; bei `replay` und `play` ist die Option unbekannt (`PGR-E2001`).
- *Roundtrip:* Die Schreiber (`yaml`, `sqlite`) tragen `offset_ms` schon; dieser Slice schreibt sie nicht neu, er belegt sie mit dem Roundtrip je Format (aufzeichnen, laden, Werte gleich) — entschieden.

**Risiken:**

- Zeitabstände auf einem ausgelasteten Rechner schwanken; der E2E-Test prüft nur Ordnung und Untergrenze (nicht absteigend, nicht negativ), nicht eine bestimmte Dauer — **Ausgang:** entfallen, Begründung in §7.
- Der Messpunkt braucht entgegen der Entscheidung doch den PGWire-Adapter, und der Slice berührt drei Schichten (`BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`, 2×, offen): Der Architect hat den Messpunkt vor dem Code in den Use Case gelegt, das Risiko ist damit klein — **Ausgang:** entfallen, Begründung in §7; die Rückführung in §4 ist nicht eingetreten.

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

Geschrieben bei der Closure (Planner, 2026-10-12, Stand `3727254`).

- **Belege zur DoD:** die Testverweise stehen in der DoD; *Belege des Implementers* unten tragen Zusage, Mutation und roten Test. Die Verifikation (`docs/reviews/2026-10-12-verifikation-slice-v1-abschluss-zeitangaben-aufzeichnen.md`) bestätigt beide Liefer-Punkte: 25 Mutanten (14 im Service, 5 im CLI-Adapter, 2 an den Schreibern, 3 an der Verdrahtung, 1 an der Verbotsliste) sind unabhängig neu gebaut und rot aus dem richtigen Grund; am gebauten Binary stimmen die Quellen der Option, ihre Rangfolge und die Abstände der Werte (302, 202, 303 bei Pausen von 300, 200, 300 ms), und zwei Läufe ohne die Option mit verschiedenen Pausen sind byte-gleich. Namensabgleich des Planners: jeder Testname aus §7 und der DoD hat genau eine Zeile `^func <Name>(` unter `internal/` oder `test/`. `make gates` lief grün am Stand `3727254` (Planner, auf sauberem Baum vor den weiteren Closure-Commits).
- **Was hat funktioniert:** Der Messpunkt im Use Case hielt den Slice bei zwei Schichten; der PGWire-Adapter blieb unberührt, die Rückführungen aus §4 traten nicht ein. Die Entscheidung des Architect vor dem Code (Messpunkt, Rundung, Uhr je Lauf) fing alle Randformen bis auf die Aufruf-Voraussetzung; die Fake-Uhr aus `-uhr` trug die Tests ohne Änderung ihrer Signatur. Die Mutationstabelle fand die Zusagen im Wortlaut der Spezifikation und ließ die Verifikation sie neu bauen.
- **Was ging anders als geplant:**
  1. Das Gegenbeispiel zu „Werte steigen nicht ab“ (gleichzeitige Aufrufe von `Query` und `ClientMessage`) stand als Satz in §7 statt als Rückgabe an den Architect (F-647, MEDIUM); entschieden nach dem Review (`fb88966`, `8ad1f3d`).
  2. Der Beginn der Messung stand in `LH-FA-21.a` *Aufzeichnen* („Annahme der Verbindung“) anders als in *Messpunkt* (V-189, LOW); angeglichen in `29c8628`.
  3. Der Platzhalter in `TestDateiUngueltig` ruhte darauf, dass `slice-v1-abschluss-sessions` `record_empty_sessions` nicht liefert, ohne dass der Nehmer es wusste (F-650); mit ihm das akzeptierte Negativ zu F-647. Beides trug der Planner bei der Closure in den Nehmer ein.
  4. §3 folgte den Code-Commits erst im vierten (F-649); die Zeile der Verbotsliste fehlte in der Tabelle (F-648).
- **Summary-Zeilen:** Review `docs/reviews/2026-10-12-review-slice-v1-abschluss-zeitangaben-aufzeichnen.md`: 0 HIGH, 1 MEDIUM, 2 LOW, 2 INFO; Klassen: Zusage hängt an nicht genannter Aufruf-Voraussetzung (F-647) · Mutation einer Zusage nicht in §7 geführt (F-648) · Plan folgt der Korrektur nicht im selben Commit (F-649) · Test ruht auf der Nicht-Lieferung eines Folge-Slice (F-650) · Stand `vollständig` weiter als die Prüfung (F-651) · Zusage „byte-gleich“ gegen Vorstand nicht direkt geprüft (F-652). Verifikation: „0 HIGH · 0 MEDIUM · 1 LOW (V-189) · 3 INFO (V-190, V-191, V-192)“. Ausgänge: F-647 durch `fb88966` und `8ad1f3d`; F-648 und F-649 durch `833960b` und `8ad1f3d`; F-650 durch die Eintragung in `slice-v1-abschluss-sessions` (`3727254`); F-651 und V-190 ohne Aktion (mechanische Folge von `SPEC-048`, die Zeile stimmt, sobald `slice-v1-abschluss-zeitangaben` die Anforderung schließt); F-652 durch den engeren Wortlaut der DoD; V-189 durch `29c8628`; V-192 ohne Aktion.
- **Entscheidungen des Planners zur Verifikation:**
  - *V-191, Lebendprüfung am Binary:* kein Binary-Fall. Der Code kennt keinen Sonderpfad für sie (`record.go` enthält kein `istLebendpruefung`); ein Fall am Binary braucht einen Client, der Text aus Kommentar oder Leerraum schickt, und `psql` tut das nicht. `TestRecordZeitLebendpruefung` trägt die Zusage mit einer Fixture, die den Pfad trifft, und sein Mutant ist rot. Ein eigener Wire-Client im Integrationstest wäre eine neue Testhilfe für einen Pfad ohne eigenen Code.
  - *V-192:* ohne Aktion. `psql` 17 kennt `\parse`, `\bind_named` und `\close_prepared` nicht; die Arten `Close` und `Flush` und mehrere Gruppen belegen die Unit-Tests je Art mit einem Mutanten. Eine einzelne Probe-Grenze des Werkzeugs, kein Muster; wiederholt sie sich, bekommt sie einen Eintrag.
- **Steering-Loop-Eintrag:** Geschärfte Regel: Sagt ein Vertrag eine Ordnung der Werte über mehrere Aufrufe einer Session zu, nennt §6 als Randform, ob die Aufrufe nacheinander oder gleichzeitig kommen, und die Spezifikation bindet die Zusage an diese Voraussetzung; ein Gegenbeispiel im Bericht des Implementers ist keine Rückgabe. Anlass: F-647 (und die Bindung F-650, die der Planner bei der Closure eintrug). Liegt in `AGENTS.md §3.12` (seit slice-v1-abschluss-zeitangaben-aufzeichnen); gezählt unter `BEO-REPO/spec-randform-erst-im-review-entschieden`.

  Retirement-Checks: §3.9 (seit welle-walking-skeleton) ist wieder aufgetreten (F-649); die Regel bleibt. §3.10 (seit welle-extended-query) ist nicht wieder aufgetreten: Jede Zusage hatte eine Mutation und einen roten Test, nur die Zeile der Verbotsliste fehlte in der Tabelle (F-648, einmalig, nicht eingetragen). §3.11 ist in der Hälfte *Plan-Zeile* wieder aufgetreten (F-652, F-651); die Regel bleibt. §3.12 ist wieder aufgetreten (F-647) und ist geschärft. §3.13 *Auch eine Bindung ist eine Sendung* (seit slice-v1-abschluss-einspielen-laufsteuerung) ist wieder aufgetreten (F-650); die Regel bleibt.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `3727254` (Zähler = Dateien unter `evidence/`); jeder Beleg trägt `evidence/slice-v1-abschluss-zeitangaben-aufzeichnen.md`.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden`: **Beleg**, 23× → 24× (F-647), Stand verkörpert, mit der Regelzeile in `AGENTS.md §3.12`.
  - `BEO-REPO/plan-folgt-korrektur-nicht`: **Beleg**, 31× → 32× (F-649), Stand verkörpert, bleibt.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: **Beleg**, 40× → 41× (F-652, F-651 in einer Datei), Stand verkörpert, Sensor `slice-harness-mutation` geplant, bleibt.
  - `BEO-REPO/folge-slice-adresse-nimmt-nicht-an`: **Beleg**, 7× → 8× (F-650 und die Bindung zu F-647), Stand verkörpert, bleibt; die Zeile zum Port der Uhr trägt die Feststellung, dass dieser Slice nur `Jetzt` benutzte.
  - Neu: `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` (1×, V-189), Stand offen. Keine Adresse: Die Ursache liegt in der Fortschreibung der Spezifikation von Abschnitt zu Abschnitt; ein Slice dafür entstünde unter 3× ohne Gegenstand.
  - Gesichtet ohne Beleg: `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (bleibt 2×) und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (bleibt 2×): Die Rückführungen aus §4 traten nicht ein, zwei Schichten, zwei Liefer-Punkte. `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (bleibt): F-647 gehört nicht dorthin, der Code entschied nichts. `BEO-REPO/liste-gruener-mutanten-unvollstaendig` und `BEO-REPO/gruener-mutant-faelschlich-aequivalent` (je 2×): Die zwei grünen Mutanten stimmen, die Verifikation baute alle neu. `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (bleibt 1×): V-192 gilt dem Client-Werkzeug, nicht dem Server. `BEO-REPO/test-faellt-unter-cpu-last-aus` (bleibt 1×): kein Auftreten.
  - Einmalig und nicht eingetragen: F-648 (Zeile fehlte, Test trug die Zusage, behoben), V-191, V-192 (Proben-Grenzen, ohne Aktion).
- **Folge-Slices:** keiner neu. Der nächste Schritt nach §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) ist `slice-v1-abschluss-zeitangaben` (Schritt 22). Angenommen haben die Sendungen: `slice-v1-abschluss-sessions` (Platzhalter in `TestDateiUngueltig` und Aufrufe einer Session nacheinander, §1, §3, §4, §6, §8) und `slice-v1-abschluss-zeitangaben-doku` (Lesart von *Messpunkt* für den Beginn, §1, aus `29c8628`).
- **Risiken aus §6:** zwei. *Zeitabstände auf einem ausgelasteten Rechner schwanken, der E2E-Test prüft nur Ordnung und Untergrenze*: **entfallen** — die Untergrenze kann Last nicht verletzen: Pausen werden unter Last länger, nie kürzer, und der Unterschied zweier abgeschnittener Lesungen ist nie kleiner als die Pause (`floor(x + 300) = floor(x) + 300`); eine Obergrenze prüft kein Test, also kein Flackern nach oben. Dass die Untergrenze die Verdrahtung trägt, zeigt der Mutant mit konstanter Uhr (`[0 0 0]`, `TestE2ERecordTiming` rot), und die Verifikation maß 302, 202 und 303 am Binary. Was die Untergrenze nicht fängt (eine Systemuhr mit falscher Skala), tragen die Tests von `slice-v1-abschluss-zeitangaben-uhr`. *Der Messpunkt braucht doch den PGWire-Adapter*: **entfallen** — der Messpunkt liegt im Use Case, `git diff --stat` der Verifikation nennt keine Datei unter `internal/adapters/driving/pgwire`, zwei Schichten bestätigt, die Rückführung aus §4 trat nicht ein; der Eintrag `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` bleibt bei 2×. Eingetreten ist keines. Die Bedingung der Aufrufe nacheinander (F-647) war kein Risiko aus §6, sondern eine Randform; sie ist entschieden und als Bindung in `slice-v1-abschluss-sessions` eingetragen.
- **Drei Paarungen:** Anker: Der Steering-Loop-Eintrag trägt das Feld `liegt in` mit `AGENTS.md §3.12`; der Absatz steht dort mit dem Anker `seit slice-v1-abschluss-zeitangaben-aufzeichnen`, und `BEO-REPO/spec-randform-erst-im-review-entschieden/state.md` nennt Zielort und Anker. Folge-Slice: `slice-v1-abschluss-sessions`, `slice-v1-abschluss-zeitangaben` und `slice-v1-abschluss-zeitangaben-doku` liegen in `next/`; `grep -n "slice-v1-abschluss-zeitangaben-aufzeichnen"` findet die Kennung des Gebers in §1 jedes der drei, und kein Ausschluss trifft die Sendungen. Register: Die vier Einträge mit Beleg und der neue tragen `evidence/slice-v1-abschluss-zeitangaben-aufzeichnen.md`; die übrigen genannten bestehen als Verzeichnisse mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft die Paarungen erneut.

### Belege des Implementers

**Läufe:** Die Mutanten liefen in einer frischen Kopie des Baums ohne `.git` (`tar`, nicht im Repo), die Unit-Mutanten mit `go test` im Image `pgwire-recorder:test` ohne Netz und die Kopie als Bind-Mount, die Mutanten des Binaries mit `make test-integration` in der Kopie. Die Datei der Mutation wurde je Mutant aus dem Repo zurückgeschrieben. `make test`, `make lint`, `make a-check` und `make a-check-negativ` liefen grün nach `0bfbfe4`; `make gates` auf sauberem Baum nennt der Bericht (Stand nach dem letzten Commit).

**Zusage · Mutation · roter Test** (Wortlaut `LH-FA-21.a` *Aufzeichnen*, *Uhr*, *Messpunkt*, Optionstabelle von `LH-FA-17.a`):

| Zusage | Mutation (nur sie) | roter Test |
|---|---|---|
| Der Wert ist der Unterschied zweier Lesungen (Ankunft minus Beginn) | Operanden vertauscht (`l.beginn - ankunft`) | `TestRecordZeitEinfach` (`[-40 -45]`) |
| Beginn der Session ist der Aufruf von `OpenSession` | Beginn konstant 0 | `TestRecordZeitEinfach` (`[140 145]`), `TestRecordZeitJedeSessionEigenerBeginn` |
| … vor dem Aufbau der Upstream-Session | Beginn nach `upstream.Open` gelesen | `TestRecordZeitBeginnVorAufbau` (`[0]` statt `[30]`) |
| Es wird nicht gerundet | Differenz auf 10 ms aufgerundet | `TestRecordZeitEinfach` (`[40 50]`) |
| Jede Session misst von ihrem eigenen Beginn | Beginn der ersten Session für alle | `TestRecordZeitJedeSessionEigenerBeginn` (`Session b: [105]` statt `[5]`) |
| Ankunft einer einfachen Anfrage ist der Aufruf von `Query` (nicht das Ende) | Lesung beim Übernehmen der Antwort | `TestRecordZeitAnkunftVorUpstream` (`[510]` statt `[10]`) |
| … gelesen, bevor der Aufruf auf etwas wartet | Lesung nach der Warteschleife von `Query` | `TestRecordZeitAnkunftVorWarten` (`Query liest die Uhr nicht, bevor sie auf die Extended-Interaktion wartet`) |
| Ankunft einer Extended-Interaktion ist die erste Client-Nachricht … | Lesung beim Abschluss der Interaktion | `TestRecordZeitExtendedErsteNachricht` (sieben Arten, `[130]` statt `[20]`), `TestRecordZeitExtendedMehrereGruppen`, `TestRecordZeitNichtAbsteigend` |
| … gleich welche Art (a: Parse, Bind, Describe, Execute, Close, Flush, Sync) | Lesung nur bei `Parse` (sonst 0) | `TestRecordZeitExtendedErsteNachricht` (sechs Arten rot, Parse grün) |
| … gleich welche Art (b) | Lesung beim `Sync` der Gruppe | `TestRecordZeitExtendedErsteNachricht` (sechs Arten rot, Sync grün), `TestRecordZeitExtendedMehrereGruppen` |
| … der ersten Gruppe (nicht jeder) | Lesung bei der ersten Nachricht jeder Gruppe | `TestRecordZeitExtendedMehrereGruppen` (`[20]` statt `[10]`), `TestRecordZeitExtendedErsteNachricht` (Art Flush) |
| Die erste Interaktion kann 0 tragen, das Feld steht dann da | `nil` bei Wert 0 | `TestRecordZeitNullSteht` (`[-1]` statt `[0]`) |
| Die Werte einer Session steigen nicht ab (Nummer = Reihenfolge der Ankunft; einfach, Extended, Pipelining) | Wert der Extended-Interaktion um 100 ms je abgeschlossener Interaktion verschoben (`[10 -80 -170 240]`) | `TestRecordZeitNichtAbsteigend` (prüft die genauen Werte `[10 20 30 240]`; eine eigene Prüfung auf Absteigen gäbe es neben ihnen nicht her) |
| Die Lebendprüfung trägt das Feld wie jede Anfrage | `nil` für Text, den `istLebendpruefung` annimmt | `TestRecordZeitLebendpruefung` (`[-1 -1]` statt `[3 7]`) |
| Eine unvollständig endende Interaktion wird nicht übernommen, hat keinen Wert | `CloseSession` hängt die offenen Interaktionen an | `TestRecordZeitUnvollstaendigOhneWert` (`[10 20]` statt `[10]`) |
| Ohne die Option fehlt das Feld | `offset` liefert immer einen Wert | `TestRecordOhneZeitangaben` (`[0 0]` statt ohne Feld) |
| Ohne die Option liest der Service die Uhr nie (Kommentar von `Zeitangaben`) | `lesung` liest immer | `TestRecordOhneZeitangaben` (`4 Lesungen ohne Zeitaufzeichnung`) |
| Dieser Slice ruft `WarteBis` nicht auf (§1) | `lesung` ruft nach `Jetzt` `WarteBis` | `TestRecordZeitWartetNie` |
| Roundtrip yaml: aufzeichnen, laden, Werte gleich | Schreiber lässt `offset_ms` weg / schreibt Wert + 1 | `TestRecordZeitRoundtrip` (`yaml … ohne offset_ms` / `13, erwartet 12`) |
| Roundtrip sqlite: aufzeichnen, laden, Werte gleich | Schreiber lässt `offset_ms` weg / schreibt Wert + 1 | `TestRecordZeitRoundtrip` (`sqlite … ohne offset_ms` / `13, erwartet 12`) |
| Ohne die Option lädt die Aufzeichnung beider Formate ohne Feld | (derselbe Test, Zweig `mit=false`) | `TestRecordZeitRoundtrip`: kein eigener Mutant, der nur diesen Zweig bricht; er ist durch `TestRecordOhneZeitangaben` und `TestE2ERecordOhneTiming` gebunden |
| Der Bootstrap verdrahtet die Systemuhr (ohne sie wäre der Unterschied 0) | `services.Zeitangaben(konstUhr{}, …)` mit Uhr, die immer 0 liest | `TestE2ERecordTiming` (`rec.yaml: [0 0 0]`, `rec.sqlite: [0 0 0]`) |
| Die Pause des Clients zwischen zwei Interaktionen ist die Untergrenze des Unterschieds (300 ms, 200 ms) | dieselbe Mutation | `TestE2ERecordTiming` (Abstände 0) |
| Der Bootstrap reicht die Option durch | `services.Zeitangaben(NeueSystemuhr(), false)` | `TestE2ERecordTiming` (`[-1 -1 -1]`), `TestE2ERecordTimingQuellen` (Umgebung und Datei) |
| Ohne die Option ist die Ausgabe unabhängig vom Zeitverhalten (YAML byte-gleich zwischen zwei Läufen mit den Pausen 0 und 300 ms; ein Vergleich gegen den Stand vor dem Slice ist nicht zugesagt, F-652) | `services.Zeitangaben(NeueSystemuhr(), true)` | `TestE2ERecordOhneTiming` (`yaml nennt offset_ms ohne die Option`, `yaml unterscheidet sich zwischen den Läufen`, sqlite `[8 10 11]`), `TestE2ERecordTimingQuellen` (`false` auf der Kommandozeile) |
| Option: Wert `--record-timing` setzt `RecordTiming` | `RecordTiming = false` | `TestRecordTimingQuellen` (Option, Wert, Datei, Umgebung), `TestLeserAlleOptionen` |
| Option: Default `false` | Default `true` | `TestLeserAlleOptionen` (`erwartet Default "false"`), `TestRecordTimingQuellen`, `TestParseRecord` |
| Option: Reihenfolge der Optionstabelle (nach `--format`) | Option nach `--log-level` | `TestLeserOptionen` |
| Option: bei `replay` und `play` unbekannt (`PGR-E2001`) | Option bei `replay` mit angemeldet | `TestLeserOptionen`, `TestLeserAlleOptionen`, `TestLeserFremdeUmgebung`; am Binary `TestE2ERecordTimingUngueltig` (`replay: Exit-Code 3` statt 2) |
| Option: Schlüssel `record_timing` im Abschnitt `record:` | Option aus `optionen("record")` entfernt | `TestLeserOptionen`, `TestRecordTimingQuellen` (`record.record_timing: unbekannter Schlüssel`) |
| Option: Schlüssel bei `replay:`, `play:` und oben unbekannt (`PGR-E2004`) | Option bei `replay` mit angemeldet | `TestDateiUngueltig`, `TestConfigShowFehler` (Fall `record_timing in replay`, Wert `true`, damit nur die Ebene den Fehler macht) |
| Option: nur `true` und `false` (Wertemenge Wahrheitswert), sonst `PGR-E2001` bzw. `PGR-E2004` | Wertemenge `text` statt Wahrheitswert | `TestLeserAlleOptionen`, `TestDateiUngueltig` (Fall `record_timing ungültig`), `TestRecordTimingQuellen` |
| Kein Paket unter `./cmd/...` importiert die Testhilfe `offsetlese` (Kommentar von `doc.go`) | leerer Import der Testhilfe in `internal/bootstrap/bootstrap.go` (frische Kopie) | `TestBinaryOhneTesthilfen` (`das Binary linkt …/internal/bootstrap/offsetlese`) |
| `--help` nennt die Option und ihre Umgebungsvariable | Hilfetext ohne den Absatz / ohne den Namen der Variable | `TestRecordTimingHilfe` (beide rot), am Binary `TestE2ERecordTimingHilfeUndConfigShow` (Variable fehlt) |

**Grüne Mutanten:**

- *Lebendprüfung:* Der Code kennt keinen Sonderpfad für sie; der Mutant der Tabelle ist künstlich (er legt einen an) und rot. Dass es keinen gibt, belegt der Test, nicht der Code.
- *Ohne-Zweig des Roundtrip* (Zeile „Ohne die Option lädt die Aufzeichnung …“ der Tabelle): kein Mutant, der nur ihn bricht; die Zusage trägt `TestRecordOhneZeitangaben` (Dienst) und `TestE2ERecordOhneTiming` (Binary, beide Formate).

**Gegenbeispiel zu „Werte steigen nicht ab“:** Kein Test fand eines. Beim Lesen des Codes fiel ein Fall auf, den der Vertrag der Benutzung ausschließt: Riefe ein Aufrufer `ClientMessage` (beginnt eine neue Extended-Interaktion) *während* einer wartenden `Query` auf, bekäme die `Query` die frühere Lesung, aber die höhere Nummer. Der PGWire-Adapter ruft beide aus derselben Schleife der Client-Richtung (`internal/adapters/driving/pgwire/server.go`, `nachricht`), `Query` blockiert sie bis zur Antwort; der Fall tritt über den Adapter nicht ein, und ein Test dafür ruft die Methoden gleichzeitig auf, wie es kein Aufrufer tut. Entscheidung des Architect (2026-10-12, F-647): die Zusage ist an die Bedingung gebunden, siehe §6 *Nicht absteigend*; Auftrag an den Implementer: der Satz im Vertrag des Ports (`internal/hexagon/ports/driving/record.go`).

**Nachführung von §3:** Die Code-Commits `9ed8d68`, `0bfbfe4` und `2e23f50` ließen §3 stehen; die Zeilen zu `services.Zeitangaben` und `offsetlese` kamen erst mit `833960b` (F-649). Die Kommentarzeilen zur Voraussetzung nacheinander kommender Aufrufe (F-647) stehen mit §3 im selben Commit.

**Grenzen:** (a) Die Untergrenze im E2E (300 ms, 200 ms) gilt für den Unterschied der Werte zweier Lesungen, die ganze Millisekunden abschneiden; der Unterschied ist nie kleiner als die Pause (die Differenz abgeschnittener Lesungen ist größer als die Dauer minus 1 und ganz). Eine obere Grenze prüft kein Test (§6, Risiko 1). (b) Das Kommandozeilen-Verhalten bei einem ungültigen Wert (die Meldung der Bibliothek nennt den Wert der Kommandozeile) ist das aller Wahrheitswert-Optionen und nicht Gegenstand dieses Slice (`LH-FA-17.a`: die Meldung zur Kommandozeile darf den Wert nennen). (c) In `TestDateiUngueltig` ist der Fall „Option, die der Stand nicht kennt“ von `record_timing` auf `record_empty_sessions` umgestellt (die Option kennt der Stand weiter nicht); `record_timing` hat eigene Fälle.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`, gelesen am Stand des Schnitts, 2026-10-11). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12) — neuer Vertrag, Randformen in §6.
- `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×, offen) — Der Messpunkt ist der Grund, dass die Größe erst die Prüfung des Architect zeigen kann; die Rückführung in §4 nennt die Bedingung vorab. Ein Eintreten wäre ein weiterer Beleg.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) — Hilfetext und Abdeckungs-Deklaration sagen nur zu, was ein Test prüft.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten: zwei (Kern mit dem Record-Service, CLI-Adapter); der Bootstrap reicht die Uhr weiter und zählt als reine Verdrahtung nicht; Messpunkt im PGWire-Adapter wäre die dritte Schicht (§4). Nachgezählt beim Schnitt und vom Architect vor dem Code (2026-10-11, `AGENTS.md` §3.13; Messpunkt im Use Case, zwei Schichten bestätigt); nachgezählt beim Eintragen der Signatur des Ports und der Verdrahtung der Systemuhr (Planner, 2026-10-12, Closure von `slice-v1-abschluss-zeitangaben-uhr`): Liefer-Punkte 2 unverändert (die Namen stehen in DoD-Punkt 1, die Verdrahtung gehört zu ihm), Schichten 2 (Kern, CLI-Adapter); der Bootstrap ruft einen Konstruktor und reicht das Ergebnis, das ist Verdrahtung und zählt nicht, die Logik der Systemuhr zählt in `-uhr`; die Sendungen aus `slice-v1-abschluss-sqlite-format-lesen` und `slice-v1-abschluss-einspielen` stehen in §1 mit den Kennungen der Geber, in DoD-Punkt 1. Nachgezählt vom Implementer nach dem Code (2026-10-12): Liefer-Punkte 2, Schichten 2 (Kern mit dem Record-Service, CLI-Adapter); `internal/bootstrap` ruft `services.Zeitangaben(NeueSystemuhr(), o.RecordTiming)` ohne Verzweigung, das ist Verdrahtung und zählt nicht (`AGENTS.md` §3.13); die Testhilfe `internal/bootstrap/offsetlese` ist kein Produkt-Code im Binary (`TestBinaryOhneTesthilfen`).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
