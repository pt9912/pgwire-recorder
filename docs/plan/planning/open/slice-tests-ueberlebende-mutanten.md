# Slice slice-tests-ueberlebende-mutanten: Tests für überlebende Mutanten, die das Verhalten ändern

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein
Abnahmeszenario und kein Meilenstein hängt an ihm. Eingesammelt wird er von der
nächsten Welle-Closure. Angelegt nach Entscheidung des Nutzers vom 2026-10-07
(Review F-459 und F-461 zu `slice-lint-bestand-kern-driven`); Reihenfolge in §4
*Start*.

**Bezug:** [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-10`](../../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Tests über die öffentliche Schnittstelle). Bindung an Entscheidungen:
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Black-Box-Tests über die
Export-Test-Brücke).

**Berührte Spec-Stellen:** `LH-FA-02.a` · `LH-FA-10.a` · `LH-FA-18.a` · `SPEC-001` · `SPEC-049` (Punkt 7, nur angewandt)

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Jeder hier gesammelte Mutant, der das Verhalten ändert und heute alle Tests
überlebt, wird von einem Test über die exportierte Schnittstelle rot
([`SPEC-049`](../../../../spec/spezifikation.md#spec-049--lint-profil-lint) Punkt 7),
bevor `slice-harness-mutation` den Bestand misst. Das Mutations-Gate startet damit
nicht mit Lücken, die schon bekannt sind. Jeder Test prüft eine Zusage, die vorher in
Lastenheft oder Spezifikation steht; wo sie fehlt, entscheidet der Architect sie vor
dem Code (§6, `AGENTS.md` §3.12).

**Herkunft:** Review F-459 und F-461 zu `slice-lint-bestand-kern-driven` (Report vom
2026-10-07, Stand `a8b341e`) und dessen Plan §7 *Grüne Mutanten, eingeordnet*: Vier
Mutanten ändern das Verhalten, sind über die Schnittstelle fangbar und waren schon vor
dem Umbau grün. Die Adresse `slice-harness-coverage` nahm sie nicht an, denn sein §1
schließt neue Tests über dem gemessenen Stand aus, und die Zeilen laufen unter Tests.
`slice-harness-mutation` schließt Tests für heute überlebende Mutanten ebenfalls aus.
Der Nutzer hat am 2026-10-07 diesen Slice als Adresse entschieden, und der Architect hat
den Plan von `slice-lint-bestand-kern-driven` darauf umgestellt (`8d765ba`).

**Gegenstand — die gesammelten Mutanten.** Kennungen wie im Review; jeder ist am
Stand `a8b341e` und am Code vor dem Umbau grün.

1. **M3 — `extendedNachspielen` ohne `bestaetigt[m.Type]--`**
   (`internal/hexagon/services/replay.go`; Review M3, M3i, F-460).
   *Verhalten:* Folgt in einer Interaktion auf ein bestätigtes `parse` ein
   abgelehntes, gilt auch das abgelehnte als angenommen, und die Diagnose einer
   späteren Abweichung nennt dessen SQL.
   *Test-Idee* (Paket `services_test`, über `NewReplayService` mit einem
   `RecordingRepository` des Tests, `OpenConnection` und `ClientMessage`): Eine
   Extended-Interaktion trägt in Gruppe 1 `parse` `s1` mit `SELECT 1` und `flush`,
   Antwort `parse_complete`; in Gruppe 2 `parse` `s1` mit `SELECT 2` und `sync`,
   Antwort `ErrorResponse` (42P05, das Statement besteht schon) und `ready_for_query`.
   Die folgende Interaktion erwartet ein `bind` auf `s1`. Der Client schickt dort ein
   `bind` mit anderem Portalnamen, und die Diagnose nennt das SQL `SELECT 1`, nicht
   `SELECT 2`.
   *Mutation, die rot werden muss:* die Zeile `bestaetigt[m.Type]--` entfernt.
   *Grenze:* Geprüft wird nur ein benanntes Statement. Beim unbenannten weicht das
   Verhalten des Servers ab (Randform *Abgelehntes parse* in §6). Geprüft wird die
   Diagnose, nicht die Antwort des Replays, denn die hängt nur an der Aufzeichnung.
   Erst mit dem Test kehrt der Satz über „die ersten n Nachrichten einer Art“ in den
   Doc-Kommentar von `extendedNachspielen` zurück (Entscheidung des Architect zu
   F-460).
2. **P8 — `spalten` ohne `TableOID`** und
3. **P2 — `spalten` ohne `ColumnNumber`**
   (`internal/adapters/driven/postgres/upstream.go`; Review P8, P2, P2alt, F-461).
   *Verhalten:* Tabellen-OID beziehungsweise Spaltennummer einer Spalte in
   `RowDescription` gehen verloren. Im Record erhält der Client 0, die Aufzeichnung
   trägt 0, und das Replay liefert 0.
   *Test-Idee* (Paket `postgres_test`, über `Upstream`, `Open` und `Query` sowie
   `Send` und `Receive` gegen einen Server des Tests über TCP-Loopback, wie
   `TestOpenUndQuery` und `TestSendUndReceive`): Der Server antwortet mit einer
   `RowDescription`, deren `TableOID` und `TableAttributeNumber` ungleich 0 und
   voneinander verschieden sind. Der Test vergleicht die Spalte in allen Feldern,
   auf dem Simple- und dem Extended-Weg.
   *Mutation, die rot werden muss:* die Zeile `TableOID:` entfernt; die Zeile
   `ColumnNumber:` entfernt; je für sich.
   *Grenze:* Geprüft wird die Übersetzung im Upstream-Adapter. Den Rückweg im
   PGWire-Adapter zum Client und die Aufzeichnung prüft der Test nicht, und die
   Integrationstests vergleichen die beiden Felder nicht (P2 dort grün). Den Rückweg
   im PGWire-Adapter hat `slice-lint-bestand-driving` gemessen: Der gleiche Mutant
   überlebt dort (T3, T5, Review F-469). Ihn übernimmt
   `slice-tests-ueberlebende-mutanten-driving` als G4 (§1 dort).
4. **Y2 — `geprueftFromDTO` mit Session-Nummer 1 statt der echten**
   (`internal/adapters/driven/recording/yaml.go`; Review Y2, Y2alt, Y2alt2, F-461).
   *Verhalten:* Die Meldung eines beschädigten Recordings (`PGR-E3003`) nennt bei
   einem Fehler an einer Interaktion ab Session 2 die Session 1. Kein Test erzeugt
   einen solchen Fehler ab Session 2.
   *Test-Idee* (Paket `recording_test`, über `Unmarshal`): eine Aufzeichnung mit zwei
   Sessions, Session 1 gültig; in Session 2 je ein Fall mit falscher
   Interaktionsnummer, mit negativem `offset_ms` und mit verletzter Form. Erwartet wird
   je `PGR-E3003` mit „Session 2“ im Text.
   *Mutation, die rot werden muss:* `geprueftFromDTO(1, ii, id)` statt
   `geprueftFromDTO(sd.ID, ii, id)`; dazu je Meldung in `geprueftFromDTO` die
   Konstante 1 statt `sessionID`.
   *Grenze:* Die Kennung einer Session ist ihre Stellung, denn `sessionFromDTO`
   prüft das vorher. `si+1` statt `sd.ID` ist deshalb gleichwertig (Review Y1), und
   der Test kann beide nicht unterscheiden. Die Zusage steht heute in keinem Stratum
   (Randform *Ort in der Meldung* in §6).

**Sammelregel für weitere Funde** (Entscheidung des Nutzers vom 2026-10-07). Findet ein
Review oder eine Verifikation einen grünen Mutanten, der das Verhalten ändert und über
die exportierte Schnittstelle fangbar ist, in einem Slice, der neue Tests ausschließt,
ist die Adresse dieser Slice. Das gilt, solange er in `open/` oder `next/` liegt. Der
Plan des Fund-Slice nennt ihn mit Test-Idee und Grenze, und der Planner trägt den
Mutanten hier unter *Gegenstand* nach, mit Quelle, Test-Idee, Mutation und Grenze.
Ergäbe das einen vierten Liefer-Punkt oder eine dritte Schicht, schneidet der Planner vor
dem Start einen zweiten Slice ab (§4). Ab `in-progress/` nimmt dieser Slice nichts mehr
an, denn das wäre eine Planänderung während der Arbeit; dann legt der Planner einen
neuen Slice an.
Als Adresse nennt diesen Slice `slice-lint-bestand-driving` (§6, Risiko *Verhalten von
`replaySitzung` ändert sich unbemerkt*). Seine Funde liegen im PGWire-Adapter und im
Bootstrap, für diesen Slice eine dritte Schicht; nach dieser Regel nimmt sie der
abgeschnittene `slice-tests-ueberlebende-mutanten-driving` an (Review F-469).

**Sichtung früherer Funde** (am Stand `8d765ba`): Unter `docs/reviews/`, in den Plänen
unter `done/` samt den drei Welle-Archiven und in den offenen Plänen steht kein
weiterer grüner, verhaltensändernder Mutant ohne Adresse. Jeder frühere ist im
Folge-Review, in der Verifikation oder in der Closure als behoben belegt (etwa F-260
mit `e29851a`, F-313, F-340, F-356, F-385, F-393, F-399, V-31, V-52, V-56 bis V-60,
V-67), in einem anderen Paket gefangen (`envFailOnUnconsumed`, F-456; PGR-E3003 nach
`formVorpruefung`, F-457) oder äquivalent (siehe unten). Die Zuordnung je Fund steht in
§7 unter *Belege der Planung*.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Session aus `Open` mit dem Lesepuffer des Aufbaus (F-444, V-85,
  `BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft`) — ein Folge-Slice hat ihn
  schon: `slice-v1-abschluss-anmeldung` führt die Test-Idee in §6, weil er die
  Aufbau-Schleife in `Open` umbaut, genau an dieser Stelle. Ein Test hier liefe gegen
  Code, den jener Slice ersetzt.
- Der Wert von `meldeFrist` (F-451, V-87) — ein Folge-Slice hat ihn schon:
  `slice-v1-abschluss-herunterfahren` führt ihn als Randform in §6 (übernommen aus
  `slice-v1-abschluss-betrieb`). Einen Wert nennt keine
  Stelle der Spezifikation, und ein Test mit der Schranke als Literal folgt erst der
  Entscheidung des Architect dort.
- Äquivalente Mutanten (aus `slice-lint-bestand-kern-driven` M5, M6, Y1, Y5 und V1;
  F-449 `weiterlesen`; V9 `newSession`; VE2) — Bestand bleibt bewusst stehen: Kein Test
  kann sie fangen, und ihre Grenzen stehen in den Plänen, die sie fanden. Wie das Gate
  mit ihnen umgeht, ist die Randform *Äquivalente Mutanten* von
  `slice-harness-mutation`.
- Mutanten, die erst die Messung des Mutations-Gates findet — anderer Vorgang: Diese
  Liste gibt es erst, wenn `slice-harness-mutation` misst, und dann liegt dieser Slice
  schon in `done/`. Für sie gilt die Rückführung von jenem Slice.
- Produkt-Code — Schicht-Abgrenzung: Der Slice schreibt Tests. Im Code ändert sich nur
  der Doc-Kommentar von `extendedNachspielen` (F-460). Macht ein Test einen Fehler im
  Code sichtbar statt einer Lücke im Test, ist das ein Befund für den Architect (§4).
- Der PGWire-Adapter (`internal/adapters/driving/pgwire`) und `test/integration` —
  Schicht-Abgrenzung: Die vier Mutanten liegen im Kern (`services`) und in den
  Driven-Adaptern (`postgres`, `recording`). Mehr als diese zwei Schichten wären ein
  Schnitt (§4); die Funde im PGWire-Adapter und im Bootstrap übernimmt
  `slice-tests-ueberlebende-mutanten-driving`.
- Die Spezifikation — anderer Vorgang: Die fehlenden Zusagen (§6) schreibt der
  Architect vor dem ersten Code-Commit; der Implementer schreibt keine.
- Der Plan von `slice-lint-bestand-kern-driven` — Bestand bleibt: Er zeigt seit
  `8d765ba` hierher, und dieser Slice ändert ihn nicht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Je Liefer-Punkt ein Mutant oder ein Paar aus §1 *Gegenstand*. Der Beleg des
Implementers steht in §7 als Mutationstabelle (Zusage · Mutation · roter Test, je
Mutant dazu der Lauf ohne den neuen Test, in dem er grün bleibt), nicht im Bericht.

- [ ] **Diagnose nach abgelehntem `parse` (M3):** Ein Test in `services_test` über
      `NewReplayService` und `ClientMessage` hält die Zusage, die der Architect nach §6
      *Abgelehntes parse* in `LH-FA-10.a` oder `LH-FA-18.a` festgehalten hat. Er wird
      rot ohne `bestaetigt[m.Type]--` und ist ohne Mutation grün. Der Doc-Kommentar
      von `extendedNachspielen` nennt die Zusage über die ersten n Nachrichten wieder,
      und zwar nur so weit, wie der Test sie prüft (`AGENTS.md` §3.11).
- [ ] **Spalten-Metadaten der `RowDescription` (P8, P2):** Ein Test in `postgres_test`
      über `Open` und `Query` sowie `Send` und `Receive` hält die Zusage aus
      `LH-FA-02.a` (Schritt 4) und `LH-FA-18.a` (*Record*: unverändert weitergeleitet)
      für `TableOID` und `ColumnNumber`, nach der Lesart aus §6 *Spalten-Metadaten*. Er
      wird rot ohne `TableOID:` und, je für sich, ohne `ColumnNumber:` in `spalten`.
- [ ] **Session in der Meldung eines beschädigten Recordings (Y2):** Ein Test in
      `recording_test` über `Unmarshal` hält die Zusage, die der Architect nach §6
      *Ort in der Meldung* in `SPEC-001` oder an der Stelle festgehalten hat, die er
      dafür wählt. Für jeden Fehler an einer Interaktion in Session 2 (Nummer,
      `offset_ms`, Form) wird er rot, wenn `geprueftFromDTO` die Session 1 nennt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` | update (Architect, vor dem Code) | Die Randformen *Abgelehntes parse*, *Ort in der Meldung* und *Was ein Test festhält* aus §6, an der Stelle, die der Architect wählt; die Lesart *Spalten-Metadaten* nur, wenn er sie nicht als entschieden bestätigt |
| `internal/hexagon/services/replay_extended_test.go` | update | M3: Diagnose nach bestätigtem und danach abgelehntem `parse` desselben benannten Statements (Negative, LH-FA-10) |
| `internal/hexagon/services/replay.go` | update (nur Kommentar) | Doc-Kommentar von `extendedNachspielen`: die Zusage über die ersten n Nachrichten, so weit der Test sie prüft (F-460) |
| `internal/adapters/driven/postgres/upstream_test.go`, `upstream_extended_test.go` | update | P8, P2: `RowDescription` mit `TableOID` und `TableAttributeNumber` ungleich 0, Spalte in allen Feldern verglichen, Simple- und Extended-Weg (Happy, LH-FA-02 und LH-FA-18) |
| `internal/adapters/driven/recording/yaml_test.go` | update | Y2: Fehler an einer Interaktion in Session 2, je Nummer, `offset_ms` und Form (Negative, LH-QA-05) |

Die Abdeckungs-Deklarationen bleiben, wie sie sind: Die Tests decken neue Fälle
bestehender Anforderungen und Pfade ab. Ob eine Zeile in den Tabellen unter
`docs/user/` hinzukommt, entscheidet `make abdeckung-check`; wenn ja, gehört sie in
denselben Commit.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-coverage` und
`slice-harness-commit-struktur-id` liegen in `done/` oder sind ausdrücklich zurückgestellt
(WIP-Limit 1). Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness, M3 vor M4):
die Slices von welle-v1-abschluss, darin `slice-harness-integration-wait` direkt vor
`slice-v1-abschluss-herunterfahren` und `slice-harness-meldungskatalog-gate` direkt vor
`slice-v1-abschluss-container`, dann `slice-harness-abdeckung-gate` und
`slice-harness-coverage`, dann die Slices von welle-erster-release, danach
`slice-harness-commit-struktur-id`, `slice-tests-ueberlebende-mutanten`,
`slice-tests-ueberlebende-mutanten-driving`, `slice-harness-mutation`; die Reihenfolge
der Wellen-Slices steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md). Die Slices beider Wellen ändern
den Code, auf dem die gesammelten Mutanten liegen; vor dem Start prüft der Architect
jeden gegen den Code nach den Wellen, ob er noch besteht und noch überlebt.

*Warum nach Coverage und nicht davor.* Technisch hängt der Slice an keinem der beiden;
der Platz folgt aus drei Gründen.

1. `slice-harness-coverage` ist der letzte Slice vor M3: Mit ihm wird
   Abnahmeszenario 17 nachweisbar. Dieser Slice trägt zu keinem Abnahmeszenario bei,
   und vor Coverage verschöbe er M3 ohne Gewinn.
2. Die Messung von Coverage braucht ihn nicht. Alle vier Mutanten liegen auf
   abgedeckten Zeilen (Review F-459). Die neuen Tests heben die Anweisungs-Abdeckung
   also höchstens wenig, und ein Anstieg nach dem Einschalten lässt das Gate grün.
3. Als Adresse sammelt er umso mehr, je später er startet. Bis dahin laufen noch
   `slice-lint-bestand-driving` (sein §6 zeigt für ungefangene Zweige hierher),
   `slice-harness-lint` und die Harness-Slices. Was deren Reviews finden, schließt er
   vor der Messung des Mutations-Gates mit.

Vor `slice-harness-mutation` steht er, damit dessen Messung nicht mit Lücken startet,
die schon bekannt sind (Entscheidung des Nutzers vom 2026-10-07). Direkt nach ihm folgt
`slice-tests-ueberlebende-mutanten-driving` mit den Funden im PGWire-Adapter und im
Bootstrap; der Architect entscheidet die Randformen *Spalten-Metadaten* und *Was ein
Test festhält* hier für beide.

Erster Schritt nach dem Start, vor jedem Code-Commit: Der Architect entscheidet die
Randformen aus §6 und hält sie in der Spezifikation fest
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`). Davor, noch in `open/`, prüft
der Planner die Sammlung: Hat *Gegenstand* mehr als drei Liefer-Punkte oder mehr als
zwei Schichten, schneidet er einen zweiten Slice ab, bevor dieser nach `next/` geht.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Ein Test braucht einen
  Zugang, den die Export-Test-Brücke nach `SPEC-049` Punkt 7 nicht hergibt, und der
  Weg über die Schnittstelle führt in eine dritte Schicht, etwa für P8 und P2 nur über
  `test/integration`. Dann wird der betroffene Mutant ein eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): Ein Test macht einen Fehler im
  Produkt sichtbar statt einer Lücke im Test, etwa wenn das Replay nach §6 *Abgelehntes
  parse* anders antworten muss, als es antwortet. Dann entscheidet der Architect über
  eine Korrektur des Codes in einem eigenen Slice; der Test wartet auf sie.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; jeder Mutant aus §1 *Gegenstand* ist in §7 mit seiner Mutation rot
belegt, und die Verifikation hat das nachgefahren. `make gates` ist grün, die
Closure-Notiz mit Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12). Der Slice liefert keinen neuen Vertrag, aber jeder
Test macht eine Zusage prüfbar. Wo die Zusage in keinem Stratum steht, entscheidet der
Architect vor dem ersten Code-Commit und hält sie in der Spezifikation fest. Was dort
nicht steht, entscheidet der Implementer nicht; er gibt es zurück.

- **Abgelehntes parse** (M3) — **offen**. `LH-FA-10.a` verlangt die „erwartete Query“,
  `LH-FA-18.a` *Mismatch* ergänzt Gruppe, Nachricht und Typ. Was die erwartete Query
  bei `bind`, `execute`, `describe` und `close` ist, sagt keine Stelle. Der Code nennt
  das SQL des Statements, das der Server bis dahin angenommen hat. Ein abgelehntes oder
  bis zum `Sync` verworfenes `parse` legt nichts an. Die Bestätigungen zählt er je Art
  in der ganzen Interaktion, und die ersten n Nachrichten einer Art mit n Bestätigungen
  gelten als angenommen. Offen ist, ob das die Zusage ist und wo sie steht. Offen ist
  außerdem das unbenannte Statement: Nach dem Quelltext von PostgreSQL
  (`exec_parse_message`) gibt der Server ein bestehendes unbenanntes Statement frei,
  bevor er ein neues `parse` darauf prüft, also auch dann, wenn dieses `parse`
  scheitert. Der Code behält das alte. Ob das eine Grenze von v1 ist oder eine Korrektur
  des Codes verlangt (dann eigener Slice, §4), entscheidet der Architect nach einer
  Prüfung gegen den Server. Bis dahin prüft der Test nur ein benanntes Statement.
- **Spalten-Metadaten** (P8, P2) — **entschieden, Lesart zu bestätigen**.
  `LH-FA-02.a` Schritt 4: „Der Recorder verändert fachliche Inhalte nicht“;
  `LH-FA-18.a` *Record*: Der Recorder „leitet alle Nachrichten unverändert … weiter“;
  das Lastenheft schließt bei `LH-FA-02` die „Veränderung von Serverantworten“ aus.
  Gelesen wird: Tabellen-OID und Spaltennummer einer `RowDescription` sind Teil der
  Serverantwort und gehen unverändert an den Client und in die Aufzeichnung. Der
  Architect bestätigt diese Lesart für den Simple-Weg, auf dem „fachliche Inhalte“
  steht; verneint er sie, prüft der Test nur den Extended-Weg.
- **Ort in der Meldung** (Y2) — **offen**. Keine Stelle in Lastenheft oder
  Spezifikation sagt, dass die Meldung eines beschädigten Recordings (`PGR-E3003`)
  Session und Interaktion nennt; das tun heute nur der Code und
  `TestVorpruefungNenntOrt`. Offen sind die Zusage selbst, ihr Ort und die Zählung
  (Kennung oder Stellung in der Datei; beide sind gleich, weil `sessionFromDTO` das
  vorher prüft). Ohne Zusage gibt es keinen Test, und Y2 wird zum äquivalenten Fall
  mit ausgewiesener Grenze; der Liefer-Punkt entfällt dann mit Begründung in §7.
- **Zugang der Tests** — **entschieden** in `SPEC-049` Punkt 7: Die Tests liegen in den
  Paketen `<name>_test` und gehen über die exportierte Schnittstelle. Ein
  `RecordingRepository` und einen Server über TCP-Loopback stellt der Test, die Brücke
  wächst nicht.
- **Was ein Test festhält** — **offen**. Ob ein Test dieses Slice von einer Meldung
  nur Code, Exit-Code und die zugesagten Bestandteile des Textes vergleicht oder den
  ganzen Text, entscheidet keine Stelle der Spezifikation und keine ADR. Die
  Entscheidung des Architect zu F-463 in `slice-lint-bestand-kern-driven` (§6 dort)
  entscheidet es auch nicht: Sie sagt, dass die Reihenfolge zweier gleichzeitiger
  Fehler kein Vertrag ist, dass nach außen nur der Text geht (`PGR-E3003`, Exit-Code
  3) und dass die Charakterisierungstests dort den Bestand festhalten, ohne etwas
  zuzusagen; diese vergleichen bewusst den ganzen Text. Wie ein Test vergleicht, der
  eine Zusage prüft, sagt sie nicht (Verifikation V-90 zu jenem Slice). Der Architect
  entscheidet das vor dem ersten Code-Commit für die Tests aus Liefer-Punkt 1 und 3
  und hält es fest; bis dahin entscheidet der Implementer es nicht (`AGENTS.md`
  §3.12).

**Risiken:**

- **Die Sammlung wächst über den Schnitt** — weitere Funde (§1 *Sammelregel*) heben
  die Zahl der Liefer-Punkte über drei oder ziehen eine dritte Schicht hinein. Der
  Planner schneidet vor `open` → `next` (§4). — **Ausgang:** — (bei Closure)
- **Der Rückweg im PGWire-Adapter ist ungemessen** — ob ein Mutant, der `TableOID`
  oder `ColumnNumber` auf dem Weg zum Client fallen lässt, dort überlebt. Der Test aus
  Liefer-Punkt 2 fängt ihn nicht. Gemessen hat ihn `slice-lint-bestand-driving` (T3, T5
  grün, §7 dort G4); er ist ein Fund nach der Sammelregel in einer dritten Schicht und
  liegt bei `slice-tests-ueberlebende-mutanten-driving`. Der Implementer fährt die
  Mutation hier nicht mehr. — **Ausgang:** — (bei Closure)
- **Der Test hält den Bestand statt der Zusage fest** — ein Vergleich des ganzen
  Meldungstextes machte jede Umformulierung rot, ohne dass eine Zusage bricht
  (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, in Gegenrichtung). Je Test nur
  das, was der Architect nach §6 *Was ein Test festhält* entscheidet. — **Ausgang:**
  — (bei Closure)
- **Mutant kommt im Build-Kontext nicht an** — eine Mutation mit unveränderter Größe
  und mtime überträgt BuildKit nicht (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`,
  1×). Die Mutationen laufen in einer frischen Kopie unter eigenem Pfad oder mit
  `go test` im Image der Stufe `deps`. — **Ausgang:** — (bei Closure)
- **Das Replay ist am unbenannten Statement falsch** — bestätigt die Prüfung aus §6
  *Abgelehntes parse* das Verhalten des Servers, nennt die Diagnose nach einem
  gescheiterten `parse` auf das unbenannte Statement ein Statement, das der Server
  nicht mehr kennt. Ausgang nach der Entscheidung des Architect: Grenze von v1 oder
  eigener Slice. — **Ausgang:** — (bei Closure)

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>

*Belege der Planung* (Planner, 2026-10-07, Stand `8d765ba`). Gesichtet wurden die
Reports unter `docs/reviews/`, die Pläne unter `done/` samt den Archiven der drei
geschlossenen Wellen und die Pläne in `open/`, `next/` und `in-progress/`, je nach den
Wörtern „grün“, „Mutant“, „fangbar“ und „überlebt“. Je früherem grünem Mutanten, der
das Verhalten ändert, sein Verbleib:

| Fund | Quelle | Verbleib |
|---|---|---|
| M2, M4, M6, N1, N4, N6, N7 (F-234, F-260) | Folge-Review `slice-walking-skeleton-record` | behoben in `e29851a` („Folge-Review-Findings F-254 bis F-272“), Verifikation an `fea1182` ohne grünen Mutanten |
| MC (V-9) | Verifikation `slice-walking-skeleton-replay` | behoben mit der zweiten Phase des Integrations-Runners (Closure, Archiv der Welle) |
| C, D, G (F-313) | Folge-Review `slice-extended-query-record` | behoben, Verifikation: C, D, G rot |
| X1 (V-22) | Verifikation `slice-extended-query-record` | behoben, `…-mutationen-slice-extended-query-record.md`: vor dem Test grün |
| D1, D2 (F-340) | Folge-Review 2 `slice-extended-query-replay` | behoben, Folge-Review 3 |
| NG2 (V-31) | Verifikation `slice-extended-query-replay` | behoben, Fälle P4b und P4c in `replay_extended_test.go` |
| M8 an `ClientMessage` (F-356) | Review `slice-extended-query-lebendpruefung` | behoben, Folge-Review |
| M6 (F-385) | Review `slice-replay-semantik-mismatch` | behoben, VM21 rot |
| M1 (F-393) | Review `slice-replay-semantik-fehlerreplay` | behoben, VF03 rot |
| M44b, M44c (F-399) | Review `slice-replay-semantik-meldungscodes` | behoben, VI1 und VI2 rot |
| VF2 (V-56), VK4 (V-57), VS5, VS6 (V-58), VV1 (V-59) | Verifikation `slice-replay-semantik-meldungscodes` | vor der Closure behoben: `TestFailGleichrangig` liest bis zum Ende des Puffers, V-57 entschieden in `a749370` mit Test in `596aa7f`, V-58 in `596aa7f`, `TestRunVersionLogLevel` |
| EV3 (V-60) | dieselbe | E2E geschlossen mit `TestE2EReplayNachDemEnde` (Closure von welle-replay-semantik) |
| VF24 (V-52) | Verifikation `slice-replay-semantik-fehlerreplay` | geschlossen in `slice-replay-semantik-meldungscodes` |
| `envFailOnUnconsumed` (F-456), PGR-E3003 nach `formVorpruefung` (F-457) | Review `slice-harness-blackbox-einstieg` | in einem anderen Paket rot (`TestE2EReplayNichtVerbraucht`; `TestUnmarshalExtendedFehler`, `TestVorpruefungNenntOrt`) |
| R1 an `Open` (F-444, V-85) | Review `slice-harness-blackbox-driven` | Adresse `slice-v1-abschluss-anmeldung` §6 |
| `meldeFrist` (F-451, V-87) | Review `slice-harness-blackbox-pgwire` | Adresse `slice-v1-abschluss-herunterfahren` §6 (übernommen aus `slice-v1-abschluss-betrieb`) |
| R1 an `weiterlesen` (F-449), V9, VE2, M14 | Reviews und Verifikationen | äquivalent |

Grüne Mutanten an Harness-Skripten (etwa V-67 an `abdeckung`, V-72 an der
Why-Prüfung) betreffen Werkzeuge, nicht das Produkt. Beide sind in ihrem Slice behoben
und gehören nicht zu dieser Sammlung.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `8d765ba` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (16×, verkörpert in `AGENTS.md`
  §3.11) — F-460 gehört zu dieser Klasse. Der Satz über die ersten n Nachrichten kehrt
  deshalb erst mit seinem Test zurück (DoD, Liefer-Punkt 1), und je Test gilt das
  Risiko *Der Test hält den Bestand statt der Zusage fest* (§6).
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (13×, verkörpert in `AGENTS.md`
  §3.10) — die vier Mutanten sind Lücken dieser Klasse an Verträgen, die vor §3.10
  entstanden. Der Slice schließt sie, er schärft die Regel nicht.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (11×, verkörpert in `AGENTS.md`
  §3.12) — die zwei offenen Randformen stehen deshalb vor dem Code in §6.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — der Slice ist
  wellenlos; §4 *Start* nennt den Architect-Schritt vor dem Code.
- `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — ein Risiko in §6.
- `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×) — die Prüfung des
  unbenannten Statements gegen den Server (§6 *Abgelehntes parse*) betrifft
  Serververhalten. Der Architect nennt, gegen welche Hauptversionen und welchen
  Quelltext er prüft.
- `BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft` (1×) — bleibt bei
  `slice-v1-abschluss-anmeldung` (§1, Abgrenzung).

Für „Folge-Slice-Adresse nimmt die Sendung nicht an“ (F-459) führt das Register seit
der Closure von `slice-lint-bestand-kern-driven` den Eintrag
`BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×: F-330, F-431, F-459; offen, der
Ausgang liegt beim Nutzer). Dieser Slice ist selbst eine Adresse: Nach der
*Sammelregel* in §1 trägt der Planner jeden neuen Fund hier unter *Gegenstand* nach,
die Annahme steht damit im Nehmer. Keiner der Einträge unter der Schwelle erreicht mit diesem Slice allein
3×; vor dem Code entsteht keine neue Lücke.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
