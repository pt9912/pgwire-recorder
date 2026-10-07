# Slice slice-tests-ueberlebende-mutanten-driving: Tests für überlebende Mutanten im PGWire-Adapter und im Bootstrap

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
nächsten Welle-Closure. Angelegt nach der *Sammelregel* in §1 von
`slice-tests-ueberlebende-mutanten` (Entscheidung des Nutzers vom 2026-10-07): Die Funde
aus `slice-lint-bestand-driving` liegen im PGWire-Adapter und im Bootstrap, für jenen
Slice eine dritte Schicht. Reihenfolge in §4 *Start*.

**Bezug:** [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Tests über die öffentliche Schnittstelle). Bindung an Entscheidungen: [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Black-Box-Tests über die Export-Test-Brücke).

**Berührte Spec-Stellen:** `LH-FA-02.a` · `LH-FA-02.b` · `LH-FA-03.a` · `LH-FA-03.b` · `LH-FA-05.a` · `LH-FA-13.a` · `LH-FA-13.b` · `LH-FA-18.a` · `SPEC-034` · `SPEC-038` · `SPEC-049` (Punkt 7, nur angewandt). Welche davon eine Zusage trägt und wo eine fehlt, entscheidet der Architect vor dem Code (§6).

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

**Ziel:** Jeder hier gesammelte Mutant im PGWire-Adapter und im Bootstrap, der das
Verhalten ändert und heute alle Tests überlebt, wird von einem Test über die
exportierte Schnittstelle rot
([`SPEC-049`](../../../../spec/spezifikation.md#spec-049--lint-profil-lint) Punkt 7),
bevor `slice-harness-mutation` den Bestand misst. Jeder Test prüft eine Zusage, die
vorher in Lastenheft oder Spezifikation steht; wo sie fehlt, entscheidet der Architect
sie vor dem Code (§6, `AGENTS.md` §3.12).

**Herkunft — Geber `slice-lint-bestand-driving`** (`AGENTS.md` §3.13). Dessen §7
*Grüne Mutanten, eingeordnet* führt sieben Mutanten G1 bis G7 mit Test-Idee und Grenze;
sein §6 (Risiko *Verhalten von `replaySitzung` ändert sich unbemerkt*) gibt sie über die
*Sammelregel* von `slice-tests-ueberlebende-mutanten` an den Planner. Jener Slice nimmt
sie nicht an: Sein §1 schließt den PGWire-Adapter als dritte Schicht aus, und G1 liegt im
Bootstrap (Review F-469 zu `slice-lint-bestand-driving`, Stand `934add1`; Verifikation
V-93 bis V-95, Stand `89fba63`). Nach seiner Sammelregel schneidet der Planner darum
diesen zweiten Slice ab. Belegt sind G2, G3 und G5 zusätzlich durch temporäre Tests der
Verifikation, G1 durch einen temporären Test des Implementers (50 von 50 Läufen rot).

**Gegenstand — die gesammelten Mutanten.** Kennungen wie in §7 des Gebers; Mutation in
dessen Mutationstabelle. Alle außer G1 sind schon am Code vor dem Umbau grün (`78f28a1`),
G1 entstand mit `pgwire.Listen(ctx, …)` in `78f28a1`.

1. **G2 — Ende der Replay-Sitzung nach Verbindungsende** (R2, R4;
   `internal/adapters/driving/pgwire/server.go`, `replayLesefehler`; Review F-467).
   *Verhalten:* Unter R2 schreibt die Sitzung beim Verbindungsende eine Fehlerantwort
   und merkt `PGR-E6001` als ersten Verbindungsfehler (Exit-Code ungleich 0). Unter R4
   liest sie nach dem Verbindungsende erneut und endet nicht, solange `ctx` läuft.
   *Test-Idee* (Paket `pgwire_test`, über `pgwire.Handle` mit Replay-Fake): Nach dem
   Startup schließt der Client ohne `Terminate`; `Handle` kehrt binnen einer Frist
   zurück (fängt R4), `CloseConnection` ist einmal gerufen und `FirstErrorCode()` bleibt
   leer (fängt R2).
   *Mutation, die rot werden muss:* Bedingung `verbindungsende(err) && false` (R2);
   `return` → `continue` nach dem Lesefehler (R4); je für sich.
   *Grenze:* Die Fehlerantwort an den geschlossenen Client ist nicht lesbar, beobachtbar
   ist nur der gemerkte Code; das Ende der Sitzung nur über eine Frist und nur, solange
   `ctx` nicht endet. Gemeint ist ein Ende, das der Client auslöst. Ein Ende durch die
   Frist des Herunterfahrens (`--shutdown-timeout`) ist nicht gemeint: Dort merkt die
   Sitzung nach `slice-v1-abschluss-betrieb` für eine unvollständige Interaktion
   `PGR-E4006`, und jener Slice ergänzt den Test, statt ihn zu ändern (Verifikation
   V-95).
2. **G3 — unlesbare Nachricht im Replay** (R3; `replayLesefehler`).
   *Verhalten:* Der Mutant beendet die Sitzung ohne Fehlerantwort und ohne gemerkten
   Code.
   *Test-Idee* (`pgwire_test`, `pgwire.Handle`): Nach dem Startup ein Nachrichtenkopf
   mit einem Typ, den `pgproto3` nicht kennt (etwa `'!'`, Länge 4); der Client erhält
   eine `ErrorResponse` mit SQLSTATE `0A000` und `PGR-E6001`, `FirstErrorCode()` ist
   `PGR-E6001`.
   *Mutation, die rot werden muss:* `s.fail(…)` in `replayLesefehler` entfernt.
   *Grenze:* Welcher Lesefehler eine unlesbare Nachricht ist, legt `pgproto3` fest; der
   Test hält nur die Art fest, die er sendet.
3. **G4 — Tabellen-OID und Spaltennummer der Spaltenbeschreibung** (T3, T5;
   `rowDescription` unter `toMessage`).
   *Verhalten:* Die beiden Felder einer `RowDescription` gehen auf dem Weg zum Client
   verloren; kein Unit- und kein Integrationstest vergleicht sie.
   *Test-Idee* (`pgwire_test`, über `pgwire.ToMessage`): eine RowDescription, deren
   Spalte in jedem Feld einen Wert ungleich 0 trägt, Tabellen-OID und Spaltennummer
   voneinander verschieden; der Test vergleicht jedes Feld der `FieldDescription`.
   *Mutation, die rot werden muss:* `TableOID` → 0; `TableAttributeNumber` → 0; je für
   sich.
   *Grenze:* Die Felder sind ein Durchreichen ohne Logik; der Test fängt das Vertauschen
   oder Weglassen eines Feldes, nicht dessen Bedeutung. G4 ist der Rückweg zum Client,
   den die Grenze von P8 und P2 in `slice-tests-ueberlebende-mutanten` als ungemessen
   führt; gemessen ist er im Geber (T3, T5 grün).
4. **G5 — Merken eines gescheiterten Versands im Replay** (X4; `replayZustellen`;
   Review F-466).
   *Verhalten:* Ohne `s.sendFailed(err)` wird ein Schreibfehler an den Client nicht als
   `PGR-E4003` gemerkt, und der Lauf endet mit Exit-Code 0.
   *Test-Idee* (Verifikation V-94): Der zweite Teil von `TestReplaySent` erzwingt den
   gescheiterten Versand schon deterministisch über `net.Pipe` (der Client schließt nach
   der Anfrage). Dort genügt die zusätzliche Prüfung `FirstErrorCode() == PGR-E4003`;
   ohne Mutant ist sie in 50 von 50 Durchläufen grün, unter X4 im ersten rot.
   *Mutation, die rot werden muss:* `s.sendFailed(err)` in `replayZustellen` entfernt.
   *Grenze:* Geprüft wird der Versandfehler über `net.Pipe` nach dem Schließen durch den
   Client. Den Pfad einer nicht abbildbaren Antwort deckt G6.
5. **G6 — Ende der Replay-Sitzung nach gescheitertem Versand** (Z3; `replayZustellen`).
   *Verhalten:* Unter Z3 liest die Sitzung nach einem gescheiterten Versand weiter und
   beantwortet die nächste Anfrage.
   *Test-Idee* (`pgwire_test`, `pgwire.Handle`): Der Replay-Fake liefert auf die erste
   Anfrage einen Antworttyp ohne Abbildung, sodass `send` scheitert, bevor es schreibt;
   `Handle` kehrt zurück, ohne eine zweite Anfrage zu lesen, und `Query` ist einmal
   gerufen.
   *Mutation, die rot werden muss:* `return false` → `return true` nach dem gescheiterten
   Versand.
   *Grenze:* nur über den Pfad der nicht abbildbaren Antwort fangbar; nach einem
   Schreibfehler an der Verbindung scheitert meist schon das nächste Lesen, und die
   Sitzung endet über den Lesefehler.
6. **G7 — Ende des Wächters der Replay-Sitzung** (V13; `replayWaechter`; Verifikation
   V-93).
   *Verhalten:* Ohne `case <-fertig:` läuft der Wächter nach dem Ende der Sitzung weiter,
   bis `ctx` endet, je Replay-Verbindung eine Goroutine; danach setzt er die Lesefrist
   einer geschlossenen Verbindung.
   *Test-Idee* (`pgwire_test`, `pgwire.Handle`): Der Client endet mit `Terminate`, `ctx`
   läuft weiter; nach der Rückkehr von `Handle` kehrt die Zahl der Goroutinen
   (`runtime.NumGoroutine`) binnen einer Frist auf den Stand vor der Verbindung zurück.
   *Mutation, die rot werden muss:* `case <-fertig:` in `replayWaechter` entfernt.
   *Grenze:* über das Protokoll nicht beobachtbar, nur über die Zahl der Goroutinen des
   Prozesses; der Test läuft nicht parallel zu anderen und braucht eine Frist. Ob das
   eine Zusage ist, entscheidet der Architect (§6 *Ende des Wächters*).
7. **G1 — `Listen` ohne `WithoutCancel`** (K4; `pgwire.Listen`, gerufen aus
   `internal/bootstrap`).
   *Verhalten:* Mit schon beendetem Kontext und einem Namen als Adresse bricht die
   Namensauflösung ab: Der Lauf endet mit `PGR-E4001` und Exit-Code 4 statt regulär mit
   Exit-Code 0.
   *Test-Idee* (Paket `bootstrap_test`, über `bootstrap.Run`): schon beendeter Kontext
   und `record --listen localhost:0`; erwartet Exit-Code 0, kein `PGR-E4001`.
   *Mutation, die rot werden muss:* `context.WithoutCancel` in `pgwire.Listen` entfernt.
   *Grenze:* nur mit einem Namen als Adresse sichtbar, weil allein die Namensauflösung
   von `net` den Kontext liest; mit einer IP-Adresse bleibt der Mutant grün. Die Zusage
   „ein Signal vor dem Öffnen ist kein Startfehler“ steht heute in keinem Stratum (§6
   *Signal vor dem Öffnen*).

Die gestrichenen Sätze in den Doc-Kommentaren von `replayLesefehler`, `replayZustellen`
(Review F-465) und `replayWaechter` (V-93) kehren nur mit ihrem Test zurück und nur so
weit, wie er sie prüft (`AGENTS.md` §3.11).

**Frist im Test (Review F-468).** `TestReplayHerunterfahrenSpaeteFrist` wartet in
`server_extended_test.go:846` ohne eigene Frist auf `<-conn.gesetzt`; unter dem Mutanten
W des Gebers wird er nur über die Zeitüberschreitung des Pakets rot (10 min, ein
Goroutinen-Dump statt einer Zusicherung). Die Frist gehört als Randform zu diesem Slice
(§6 *Warten auf einen Kanal*). Der Architect entscheidet sie für alle Tests, die ohne
Frist auf einen Kanal warten, nicht nur für diesen; dieser Slice setzt sie in seinen
beiden Schichten um. Wie ein Gate einen Mutanten zählt, der nur über die
Zeitüberschreitung rot wird, ist die Randform *Laufzeit* von `slice-harness-mutation`.

**Sammelregel für weitere Funde.** Wie in §1 von `slice-tests-ueberlebende-mutanten`,
für Funde im PGWire-Adapter, im CLI-Adapter und im Bootstrap: Findet ein Review oder eine
Verifikation dort einen grünen, verhaltensändernden Mutanten, der über die exportierte
Schnittstelle fangbar ist, in einem Slice, der neue Tests ausschließt, trägt der Planner
ihn hier unter *Gegenstand* nach, solange dieser Slice in `open/` oder `next/` liegt.
Ergäbe das einen vierten Liefer-Punkt, schneidet er vor dem Start einen weiteren Slice ab.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Mutanten im Kern (`internal/hexagon/`) und in den Driven-Adaptern — übernimmt
  `slice-tests-ueberlebende-mutanten` (M3, P8, P2, Y2); er liegt vor diesem Slice.
- Das Ende einer Replay-Sitzung durch die Frist des Herunterfahrens und `PGR-E4006` —
  übernimmt `slice-v1-abschluss-betrieb` (§1 dort); G2 meint nur ein Ende, das der
  Client auslöst.
- Fristen in `test/integration` — übernimmt `slice-harness-integration-wait` für das
  Warten im Integrations-Testgeschirr; Fristen in Tests anderer Schichten setzt dieser
  Slice nicht, die Entscheidung aus §6 *Warten auf einen Kanal* nennt für sie eine
  eigene Adresse.
- Äquivalente Mutanten des Gebers („`WithoutCancel` entfernt“ an den Befunden 10 bis 13,
  Rückgabewert von `weiterlesen`, F-449) — Bestand bleibt bewusst stehen: Kein Test kann
  sie fangen; ihre Grenzen stehen in §6 des Gebers.
- Wie ein Mutations-Gate einen Mutanten zählt, der nur über eine Zeitüberschreitung rot
  wird — anderer Vorgang: Randform *Laufzeit* von `slice-harness-mutation`.
- Mutanten, die erst die Messung des Mutations-Gates findet — anderer Vorgang: Diese
  Liste entsteht erst in `slice-harness-mutation`.
- Produkt-Code — Schicht-Abgrenzung: Der Slice schreibt Tests. Im Code ändern sich nur
  die Doc-Kommentare oben. Macht ein Test einen Fehler im Code sichtbar statt einer
  Lücke im Test, ist das ein Befund für den Architect (§4).
- Die Spezifikation — anderer Vorgang: Die fehlenden Zusagen (§6) schreibt der
  Architect vor dem ersten Code-Commit; der Implementer schreibt keine.
- Der Plan des Gebers `slice-lint-bestand-driving` — Bestand bleibt: Seine §7 zeigt
  hierher, dieser Slice ändert sie nicht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Je Liefer-Punkt eine Gruppe aus §1 *Gegenstand*, nach Lieferwert geschnitten; die
Sendung kommt von `slice-lint-bestand-driving` (§7 dort, G1 bis G7). Der Beleg des
Implementers steht in §7 als Mutationstabelle (Zusage · Mutation · roter Test, je Mutant
dazu der Lauf ohne den neuen Test, in dem er grün bleibt), nicht im Bericht. Entscheidet
der Architect zu einem Mutanten, dass keine Zusage besteht (§6), entfällt dessen Test
mit Begründung in §7, und der Mutant ist mit Grenze als äquivalent geführt.

- [ ] **Lesen der Replay-Sitzung (G2, G3):** Tests in `pgwire_test` über `pgwire.Handle`
      halten die Zusagen, die der Architect nach §6 *Verbindungsende im Replay* und
      *Unlesbare Nachricht im Replay* festgehalten hat (Kandidaten `LH-FA-02.b`,
      `LH-FA-03.b`, `LH-FA-05.a`, `LH-FA-13.b`, `SPEC-034`). Sie werden rot unter R2,
      unter R4 und unter R3, je für sich, und sind ohne Mutation grün. Dazu trägt
      `TestReplayHerunterfahrenSpaeteFrist` beim Warten auf `conn.gesetzt` die Frist nach
      §6 *Warten auf einen Kanal* (Review F-468), und der Mutant W des Gebers wird dort
      rot, ohne dass das Paket die Zeit überschreitet.
- [ ] **Versand und Ende der Replay-Sitzung (G5, G6, G7):** `TestReplaySent` prüft nach
      dem gescheiterten Versand `FirstErrorCode() == PGR-E4003` (Verifikation V-94) und
      wird rot unter X4; ein Test über `pgwire.Handle` mit nicht abbildbarer Antwort wird
      rot unter Z3; für G7 ein Test nach §6 *Ende des Wächters*, rot ohne
      `case <-fertig:`, oder die Begründung in §7, dass keine Zusage besteht. Zusagen
      nach §6 *Versand im Replay* (Kandidaten `LH-FA-13.b`, `SPEC-034`).
- [ ] **Abbildung und Öffnen am Rand (G4, G1):** Ein Test in `pgwire_test` über
      `pgwire.ToMessage` hält Tabellen-OID und Spaltennummer der Spaltenbeschreibung
      nach §6 *Spalten-Metadaten* (`LH-FA-02.a`, `LH-FA-03.a`, `LH-FA-18.a`) und wird
      rot unter T3 und, je für sich, unter T5. Ein Test in `bootstrap_test` über
      `bootstrap.Run` hält die Zusage nach §6 *Signal vor dem Öffnen* (Kandidat
      `LH-FA-13.a`) und wird rot ohne `context.WithoutCancel` in `pgwire.Listen`, oder
      §7 begründet, dass keine Zusage besteht.
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
| `spec/spezifikation.md` | update (Architect, vor dem Code) | Die offenen Randformen aus §6, an der Stelle, die der Architect wählt; die Lesart *Spalten-Metadaten* nur, wenn er sie nicht mit der Entscheidung in `slice-tests-ueberlebende-mutanten` bestätigt |
| `internal/adapters/driving/pgwire/server_test.go` | update | G2, G3: Verbindungsende ohne `Terminate` und unlesbare Nachricht im Replay (Negative, LH-FA-03, LH-FA-05); G5: `TestReplaySent` um `FirstErrorCode()` erweitert, G6: nicht abbildbare Antwort, G7: Ende des Wächters (Negative, LH-FA-13) |
| `internal/adapters/driving/pgwire/server_extended_test.go` | update | Frist beim Warten auf `conn.gesetzt` (Zeile 846 am Stand `893b54e`) und an weiteren Stellen dieser Schicht nach §6 *Warten auf einen Kanal* |
| `internal/adapters/driving/pgwire/*_test.go` (Abbildung) | update | G4: `ToMessage` mit RowDescription, alle Felder verglichen (Happy, LH-FA-18) |
| `internal/bootstrap/bootstrap_test.go` | update | G1: `bootstrap.Run` mit beendetem Kontext und `--listen localhost:0` (Boundary, LH-FA-13) |
| `internal/adapters/driving/pgwire/server.go` | update (nur Kommentar) | Doc-Kommentare von `replayLesefehler`, `replayZustellen` und `replayWaechter`: die gestrichenen Sätze, so weit der jeweilige Test sie prüft |

Die Abdeckungs-Deklarationen bleiben, wie sie sind: Die Tests decken neue Fälle
bestehender Anforderungen und Pfade ab. Ob eine Zeile in den Tabellen unter
`docs/user/` hinzukommt, entscheidet `make abdeckung-check`; wenn ja, gehört sie in
denselben Commit.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-tests-ueberlebende-mutanten` liegt in
`done/` oder ist ausdrücklich zurückgestellt (WIP-Limit 1). Reihenfolge nach
Entscheidung des Nutzers vom 2026-10-05, 2026-10-06 und 2026-10-07:
`slice-harness-lint-werkzeug`, die vier Umstellungs-Slices,
`slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`, `slice-harness-lint`,
`slice-harness-commit-struktur-id`, `slice-harness-integration-wait`,
`slice-harness-meldungskatalog-gate`, `slice-harness-abdeckung-gate`,
`slice-harness-coverage`, `slice-tests-ueberlebende-mutanten`, dieser Slice,
`slice-harness-mutation`.

*Warum hier.* Technisch hängt der Slice an keinem der beiden Nachbarn. Direkt nach
`slice-tests-ueberlebende-mutanten`, weil beide nach derselben Regel sammeln und der
Architect die Randform *Spalten-Metadaten* dort für den Rückweg mitentscheidet. Vor
`slice-harness-mutation`, damit dessen Messung nicht mit Lücken startet, die schon
bekannt sind, und damit die Randform *Laufzeit* dort die Entscheidung *Warten auf einen
Kanal* von hier vorfindet. Vor `slice-v1-abschluss-betrieb` ist er nicht gebunden; liegt
jener früher in `done/`, prüft der Architect G2 gegen dessen Ende durch die Frist.

Erster Schritt nach dem Start, vor jedem Code-Commit: Der Architect entscheidet die
Randformen aus §6 und hält sie in der Spezifikation fest
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`). Davor, noch in `open/`, prüft
der Planner die Sammlung: Hat *Gegenstand* mehr als drei Liefer-Punkte oder mehr als
zwei Schichten, schneidet er einen weiteren Slice ab, bevor dieser nach `next/` geht.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Ein Test braucht einen
  Zugang, den die Export-Test-Brücke nach `SPEC-049` Punkt 7 nicht hergibt, oder die
  Frist nach §6 *Warten auf einen Kanal* trifft so viele Tests dieser Schicht, dass sie
  mit den Mutanten nicht in eine Review-Sitzung passt; dann wird der betroffene Teil
  ein eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): Ein Test macht einen Fehler im
  Produkt sichtbar statt einer Lücke im Test, etwa wenn die Replay-Sitzung bei einem
  Verbindungsende anders antworten muss, als sie antwortet. Dann entscheidet der
  Architect über eine Korrektur des Codes in einem eigenen Slice; der Test wartet auf
  sie.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; jeder Mutant aus §1 *Gegenstand* ist in §7 mit seiner Mutation rot
belegt oder mit Begründung als äquivalent geführt, und die Verifikation hat das
nachgefahren. `make gates` ist grün, die Closure-Notiz mit Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12). Der Slice liefert keinen neuen Vertrag, aber jeder
Test macht eine Zusage prüfbar. Für das Replay sagt die Spezifikation an keiner Stelle
ausdrücklich, was ein Lesefehler oder ein Versandfehler an der Client-Verbindung ist;
`LH-FA-02.b` regelt Verbindungsende und nicht lesbare Nachricht nur für `record`. Wo die
Zusage in keinem Stratum steht, entscheidet der Architect vor dem ersten Code-Commit und
hält sie in der Spezifikation fest. Was dort nicht steht, entscheidet der Implementer
nicht; er gibt es zurück.

- **Verbindungsende im Replay** (G2) — **offen**. Schließt der Client die Verbindung
  ohne `Terminate`, ist das im Replay regulär, ohne gemerkten Fehler, und die Sitzung
  endet. Kandidaten: `LH-FA-02.b` *Verbindungsende* (für `record`), `LH-FA-03.b`
  *Zeitpunkte* („der Client schließt sie“ als ein Ende der Verbindung). Offen ist, ob
  das auch mitten in einer Interaktion gilt oder dort `PGR-E4003` ist; der Test aus G2
  schließt nach dem Startup, vor jeder Anfrage. Abzugrenzen vom Ende durch die Frist
  (`slice-v1-abschluss-betrieb`, Verifikation V-95).
- **Unlesbare Nachricht im Replay** (G3) — **offen**. Kandidaten: `LH-FA-05.a` *Nicht
  unterstützte Protokollnachrichten*, `SPEC-034` *Zustellung an den Client* (SQLSTATE
  `0A000`), `LH-FA-13.b`. Offen ist, ob eine Nachricht, deren Typ die Bibliothek nicht
  kennt, im Replay wie in `record` eine nicht unterstützte Interaktion ist.
- **Versand im Replay** (G5, G6) — **offen**. Dass ein gescheiterter Versand an den
  Client `PGR-E4003` ist und die Sitzung beendet, steht für das Replay in keinem
  Stratum; für eine nicht abbildbare Antwort ist offen, welcher Code gemerkt wird.
- **Ende des Wächters** (G7) — **offen**. Ob das Ende einer Goroutine mit ihrer Sitzung
  eine Zusage ist (etwa unter `SPEC-032`) oder eine Eigenschaft des Codes ohne Vertrag.
  Ohne Zusage bleibt G7 mit Grenze als äquivalent für Tests über die Schnittstelle
  geführt, und der Satz im Doc-Kommentar kehrt nicht zurück.
- **Spalten-Metadaten** (G4) — **zu bestätigen**. Dieselbe Lesart wie in
  `slice-tests-ueberlebende-mutanten` §6: Tabellen-OID und Spaltennummer einer
  `RowDescription` sind Teil der Serverantwort und gehen unverändert an den Client
  (`LH-FA-02.a` Schritt 4, `LH-FA-18.a` *Record*) und im Replay so, wie sie aufgezeichnet
  sind (`LH-FA-03.a` *Ausgabe*). Der Architect entscheidet sie einmal für beide Slices.
- **Signal vor dem Öffnen** (G1) — **offen**. Keine Stelle sagt, wie ein Lauf endet,
  dessen Kontext vor dem Öffnen des Listeners schon beendet ist: regulär mit Exit-Code 0
  (heute) oder als Startfehler `PGR-E4001` (`LH-FA-13.b` *Startfehler*,
  `LH-FA-13.a`). Ohne Zusage gibt es keinen Test, und G1 bleibt mit Grenze als
  äquivalent geführt; der Doc-Kommentar von `Listen` sagt dazu weiter nichts zu.
- **Warten auf einen Kanal** (Review F-468) — **offen**. Ob ein Test, der auf einen
  Kanal wartet, eine eigene Frist trägt, wie lang sie ist und wo die Regel steht (etwa
  `SPEC-038`). Der Architect entscheidet für alle Tests, die ohne Frist auf einen Kanal
  warten, nicht nur für `server_extended_test.go:846`; dieser Slice setzt die Regel im
  PGWire-Adapter und im Bootstrap um, für andere Schichten nennt die Entscheidung eine
  Adresse. Die Zählung eines Mutanten, der nur über die Zeitüberschreitung rot wird,
  bleibt die Randform *Laufzeit* von `slice-harness-mutation`.
- **Zugang der Tests** — **entschieden** in `SPEC-049` Punkt 7: Die Tests liegen in
  `pgwire_test` und `bootstrap_test` und gehen über die exportierte Schnittstelle; die
  Brücke wächst nicht.
- **Was ein Test festhält** — wie in `slice-tests-ueberlebende-mutanten` §6, dort
  **offen**; die Entscheidung gilt für beide Slices.

**Risiken:**

- **Die Sammlung wächst über den Schnitt** — weitere Funde (§1 *Sammelregel*), etwa aus
  `slice-harness-lint`, heben die Liefer-Punkte über drei. Der Planner schneidet vor
  `open` → `next` (§4). — **Ausgang:** — (bei Closure)
- **Ein Test hängt am Zeitverhalten** — G2, G7 und die Frist aus F-468 beobachten über
  eine Frist; zu knapp gesetzt macht sie den Test flatterhaft, zu weit lässt sie einen
  roten Lauf lange hängen (`BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`, 2×;
  `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients`, 1×, verwandt). Je Test mindestens 50 Läufe ohne Mutant grün, in §7. — **Ausgang:** — (bei
  Closure)
- **Der Test hält den Bestand statt der Zusage fest** — ein Vergleich des ganzen
  Meldungstextes machte jede Umformulierung rot, ohne dass eine Zusage bricht. Je Test
  nur das, was der Architect nach §6 *Was ein Test festhält* entscheidet. —
  **Ausgang:** — (bei Closure)
- **Mutant kommt im Build-Kontext nicht an** — eine Mutation mit unveränderter Größe und
  mtime überträgt BuildKit nicht (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`,
  1×). Die Mutationen laufen in einer frischen Kopie unter eigenem Pfad. —
  **Ausgang:** — (bei Closure)

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
- **Drei Paarungen:** <…>

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `893b54e` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (17×, verkörpert in `AGENTS.md`
  §3.11; Sensor geplant mit `slice-harness-mutation`) — die nach F-465 und V-93
  gestrichenen Sätze kehren nur mit ihrem Test zurück (§1), und je Test gilt das Risiko
  *Der Test hält den Bestand statt der Zusage fest* (§6).
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (13×, verkörpert in `AGENTS.md`
  §3.10) — die sieben Mutanten sind Lücken dieser Klasse an Verträgen, die vor §3.10
  entstanden; der Slice schließt sie, er schärft die Regel nicht.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (11×, verkörpert in `AGENTS.md`
  §3.12) — die offenen Randformen stehen deshalb vor dem Code in §6.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, verkörpert in `AGENTS.md` §3.13) —
  dieser Slice ist die Adresse; §1 nennt den Geber, und die Sendung steht in §1 und DoD.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — der Slice ist wellenlos;
  §4 *Start* nennt den Architect-Schritt vor dem Code.
- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (1×) — G2, G7 und die Frist aus
  F-468 beobachten über eine Frist; Risiko in §6.
- `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — ein Risiko in §6.

Nachtrag bei der Closure von `slice-lint-bestand-driving` (Stand `7cb298c`), zwei neue
Einträge mit je zwei Belegen, beide aus dem Geber:

- `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` (2×) — die Randform *Warten auf einen
  Kanal* in §6 ist seine Adresse; berührt dieser Slice ihn noch einmal, etwa weil ein
  neuer Test ohne Frist wartet, erreicht er 3×.
- `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (2×) — trifft diesen Slice nur, wenn
  Review oder Verifikation hier einen grünen Mutanten finden, den §7 nicht führt.

Keiner der Einträge unter der Schwelle erreicht mit diesem Slice allein 3×; vor dem Code
entsteht keine neue Lücke.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
