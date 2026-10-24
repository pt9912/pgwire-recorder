# Slice slice-tests-ueberlebende-mutanten-driving: Tests für überlebende Mutanten im PGWire-Adapter und im Bootstrap

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten (Entscheidung des Nutzers vom 2026-10-18; bis dahin ohne Welle). Angelegt nach der *Sammelregel* in §1 von
`slice-tests-ueberlebende-mutanten` (Entscheidung des Nutzers vom 2026-10-07): Die Funde
aus `slice-lint-bestand-driving` liegen im PGWire-Adapter und im Bootstrap, für jenen
Slice eine dritte Schicht. Reihenfolge in §4 *Start*.

**Bezug:** [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-18`](../../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Tests über die öffentliche Schnittstelle). Bindung an Entscheidungen: [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Black-Box-Tests über die Export-Test-Brücke).

**Berührte Spec-Stellen:** `LH-FA-02.a` · `LH-FA-03.a` · `LH-FA-13.a` · `LH-FA-18.a` · `SPEC-032` · `SPEC-038` · `SPEC-049` (Punkt 7, nur angewandt). Welche davon eine Zusage trägt und wo eine fehlt, entscheidet der Architect vor dem Code (§6).

**Verantwortlich:** pt9912 (Zweig impl-drv1)

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
Verhalten ändert und am Stand `f4ff9387` alle Tests überlebt (G1, G4; G7 ist kein Test, siehe
unten), wird von einem Test über die exportierte Schnittstelle rot
([`SPEC-049`](../../../../harness/vertraege/SPEC-049-lint-profil.md#spec-049--lint-profil-lint) Punkt 7),
bevor `slice-harness-mutation` den Bestand misst. Dazu warten die Tests dieser beiden
Schichten, die ohne Frist auf einen Kanal oder das Ende eines Laufs warten, künftig mit
der Frist nach `SPEC-038` *Warten in Tests*. Jeder Test prüft eine Zusage, die vorher in
Lastenheft oder Spezifikation steht; wo sie fehlt, entscheidet der Architect sie vor dem
Code (§6, `AGENTS.md` §3.12).

**Herkunft — Geber `slice-lint-bestand-driving`** (`AGENTS.md` §3.13). Dessen §7
*Grüne Mutanten, eingeordnet* führt sieben Mutanten G1 bis G7 mit Test-Idee und Grenze;
sein §6 (Risiko *Verhalten von `replaySitzung` ändert sich unbemerkt*) gibt sie über die
*Sammelregel* von `slice-tests-ueberlebende-mutanten` an den Planner. Jener Slice nimmt
sie nicht an: Sein §1 schloss den PGWire-Adapter als dritte Schicht aus, und G1 liegt im
Bootstrap (Review F-469 zu `slice-lint-bestand-driving`; Verifikation V-93 bis V-95). Nach
seiner Sammelregel schnitt der Planner darum diesen zweiten Slice ab.

**Aktualität** (Planner, 2026-10-24, Stand `f4ff9387`; **gemessen**, nicht am Code
gesichtet — Lehre Y2 des Gebers `slice-tests-ueberlebende-mutanten`): Je Mutant eine
Kopie des Arbeitsbaums unter eigenem Pfad, `go test -count=1 -timeout 100s` im Image
`pgwire-recorder:test` ohne Netz, ein Lauf je Mutant (kein Lauf von `make gates`, keine
50 Wiederholungen). Der Code ist seit dem Geber umgebaut: `replayZustellen` meldet über
den Typ `folge` statt über `bool`, das Zwangsende und das Herunterfahren stehen in
`replayLauf`/`replayHerunterfahren`; die Mutation Z3 heißt darum
`return folgeEnde` → `return folgeWeiter` nach `s.sendFailed(err)`.

| Mutant | Mutation | Ergebnis am Stand | gefangen von |
|---|---|---|---|
| G2 / R2 | `&& false` in der Bedingung von `verbindungsende(err)` in `replayLesefehler` | rot | `TestProtokollrandStartnachricht`, `TestProtokollrandClientNachrichten`, `TestProtokollrandAnfragenFolgen` |
| G2 / R4 | `return false` → `continue` nach `replayLesefehler` in `replayLauf` | rot, nach der Frist von 5 s in `warte`; danach läuft das Paket bis zur Zeitüberschreitung (die Schleife endet nicht) | `TestReplayEndeTerminateMittenInDerGruppe` |
| G3 / R3 | `s.fail(…)` in `replayLesefehler` entfernt | rot | `TestProtokollrandClientNachrichten` |
| G5 / X4 | `s.sendFailed(err)` in `replayZustellen` entfernt | rot | `TestReplayEndeSchreibfehlerBleibtUnverbraucht`, `TestReplayEndeNurAlsLetzteAntwort` |
| G6 / Z3 | `return folgeEnde` → `return folgeWeiter` nach `s.sendFailed(err)` | rot, nach der Frist von 5 s | `TestReplayEndeNurAlsLetzteAntwort` |
| G4 / T3 | `TableOID: c.TableOID` → `0` in `rowDescription` | **grün** | — |
| G4 / T5 | `TableAttributeNumber: c.ColumnNumber` → `0` | **grün** | — |
| G7 / V13 | `case <-fertig:` in `replayWaechter` entfernt | **grün** | — |
| G1 / K4 | `context.WithoutCancel` in `pgwire.Listen` entfernt | **grün** (`bootstrap_test` und `pgwire_test`) | — |

Pfade und Funktionen aus §3 bestehen; `TestReplayHerunterfahrenSpaeteFrist` wartet seit
dem Geber in `server_extended_test.go:848` (damals 846) weiter ohne Frist auf
`<-conn.gesetzt`. Kein Test in `bootstrap_test` ruft `record` oder `replay` mit einem
Namen als Adresse (alle `127.0.0.1:0`). Seit der Anlage hat kein Review und keine
Verifikation einen weiteren Mutanten an diese Adresse gesendet (Suche nach der Kennung
unter `docs/`: nur Pläne und das Register des Gebers).

**Gegenstand — die gesammelten Mutanten.** Kennungen wie in §7 des Gebers; Mutation in
dessen Mutationstabelle, am heutigen Code wie in der Tabelle oben. Zwei bekommen einen Test
(G1, G4), G7 entfällt als Test (keine Zusage), dazu die Frist:

1. **Frist beim Warten auf einen Kanal** (Review F-468). `TestReplayHerunterfahrenSpaeteFrist`
   wartet in `server_extended_test.go:848` ohne eigene Frist auf `<-conn.gesetzt`; unter dem
   Mutanten W des Gebers wird er nur über die Zeitüberschreitung des Pakets rot (10 min, ein
   Goroutinen-Dump statt einer Zusicherung). Die Frist gehört als Randform zu diesem Slice
   (§6 *Warten auf einen Kanal*). Der Architect entscheidet sie für alle Tests dieser beiden
   Schichten, die ohne Frist auf einen Kanal oder das Ende eines Laufs warten; die Sichtung
   des Planners nennt als weitere Kandidaten `server_tls_test.go:54` (`server = <-angenommen`),
   `bootstrap_test.go:353` (`exit := <-ende`) und `tls_test.go:184` (`<-ende`, Paket
   `bootstrap`). Wie ein Gate einen Mutanten zählt, der nur über die Zeitüberschreitung rot
   wird, ist die Randform *Laufzeit* von `slice-harness-mutation`. Die Prüfung auf Warten
   ohne Frist in `make lint` kommt nach diesem Slice mit `slice-harness-lint-warten-ohne-frist`
   (Abgrenzung dort, §1): Sie setzt an Zeile 848 und in den beiden Schichten dieses Slice
   keine zweite Frist, sie findet, was hier übrig bleibt.
2. **G4 — Tabellen-OID und Spaltennummer der Spaltenbeschreibung** (T3, T5;
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
   führt; gemessen ist er im Geber (T3, T5 grün) und am Stand `f4ff9387` (siehe oben). Die
   Spezifikation sagt die Kette zum Client **und** zur Aufzeichnung zu (`LH-FA-02.a`
   *Spalten-Metadaten*); die Übersetzung im Upstream-Adapter hält
   `slice-tests-ueberlebende-mutanten` (P8, P2, geschlossen), den Rückweg zum Client hält G4
   hier. Die Aufzeichnung der zwei Felder (Schreiben in `yaml` und `sqlite`) gehört nicht
   zu G4 (siehe *Ausdrücklich NICHT*).
3. **G7 — Ende des Wächters der Replay-Sitzung** (V13; `replayWaechter`; Verifikation
   V-93).
   *Verhalten:* Ohne `case <-fertig:` läuft der Wächter nach dem Ende der Sitzung weiter,
   bis `ctx` endet, je Replay-Verbindung eine Goroutine; danach setzt er die Lesefrist
   einer geschlossenen Verbindung.
   *Test-Idee* (`pgwire_test`, `pgwire.Handle`): Der Client endet mit `Terminate`, `ctx`
   läuft weiter; nach der Rückkehr von `Handle` kehrt die Zahl der Goroutinen
   (`runtime.NumGoroutine`) binnen einer Frist auf den Stand vor der Verbindung zurück.
   *Mutation, die rot werden muss:* `case <-fertig:` in `replayWaechter` entfernt.
   *Grenze:* über das Protokoll nicht beobachtbar, nur über die Zahl der Goroutinen des
   Prozesses; der Test läuft nicht parallel zu anderen und braucht eine Frist. Der
   Architect hat entschieden: keine Zusage (§6 *Ende des Wächters*, `SPEC-032`); der Test
   entfällt.
4. **G1 — `Listen` ohne `WithoutCancel`** (K4; `pgwire.Listen`, gerufen aus
   `internal/bootstrap`).
   *Verhalten:* Mit schon beendetem Kontext und einem Namen als Adresse bricht die
   Namensauflösung ab: Der Lauf endet mit `PGR-E4001` und Exit-Code 4 statt regulär mit
   Exit-Code 0.
   *Test-Idee* (Paket `bootstrap_test`, über `bootstrap.Run`): schon beendeter Kontext
   und `record --listen localhost:0`; erwartet Exit-Code 0, kein `PGR-E4001`. Der Test
   führt auch `replay`, die IP-Adresse als Kontrolle und die beiden Startfehler, die
   `LH-FA-13.a` *Startphase* im selben Satz behält (Name, der nicht auflöst, und belegter
   Port: Exit-Code 4, `PGR-E4001`); ohne sie bliebe ein Mutant grün, der das Signal vor dem
   Öffnen jeden Startfehler schlucken lässt.
   *Mutation, die rot werden muss:* `context.WithoutCancel` in `pgwire.Listen` entfernt.
   *Grenze:* nur mit einem Namen als Adresse sichtbar, weil allein die Namensauflösung
   von `net` den Kontext liest; mit einer IP-Adresse bleibt der Mutant grün. Die Zusage
   „ein Signal vor dem Öffnen ist kein Startfehler“ steht in `LH-FA-13.a` *Startphase* (§6
   *Signal vor dem Öffnen*, vom Architect für die Namensauflösung ergänzt); geprüft nur unter Linux.

**Entfallen als Gegenstand (Planner, 2026-10-24, nach der Messung oben):** G2, G3, G5 und
G6 sind am Stand `f4ff9387` gefangen; ihre Tests entstanden nach der Anlage dieses Plans
(2026-10-07): `TestProtokollrand…` mit `slice-v1-abschluss-protokollrand` (2026-10-15),
`TestReplayEnde…` mit `slice-v1-abschluss-verbindungsende-replay` und
`slice-v1-abschluss-verbindungsende-replay-extended` (2026-10-16, nach `git log -S`). Kein
neuer Test: Der Bestand fängt sie, und die Aussage „überlebt alle Tests“ des Gebers
(Review F-466, F-467, Verifikation V-94, V-95) gilt für sie nicht mehr. Der Implementer fährt die vier
Mutationen am Endstand einmal nach und trägt die Läufe in §7 ein; wird eine von ihnen dann
grün (ein anderer Slice änderte einen Test), ist sie Gegenstand wie G1, G4 und G7. Die
Sätze in den Doc-Kommentaren von `replayLesefehler` und `replayZustellen` (Review F-465)
kehren nur zurück, soweit ein benannter Test sie prüft, und der Implementer nennt ihn
in §7 (`AGENTS.md` §3.11); der von `replayWaechter` (V-93) kehrt nur mit dem Test aus G7
zurück.

**Sammelregel für weitere Funde.** Wie in §1 von `slice-tests-ueberlebende-mutanten`,
für Funde im PGWire-Adapter, im CLI-Adapter und im Bootstrap: Findet ein Review oder eine
Verifikation dort einen grünen, verhaltensändernden Mutanten, der über die exportierte
Schnittstelle fangbar ist, in einem Slice, der neue Tests ausschließt, trägt der Planner
ihn hier unter *Gegenstand* nach, solange dieser Slice in `open/` oder `next/` liegt.
Ergäbe das einen vierten Liefer-Punkt, schneidet er vor dem Start einen weiteren Slice ab.
Seit dem Anspruch (`in-progress/`) trägt er nichts mehr nach; ein neuer Fund bekommt einen
eigenen Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Mutanten im Kern (`internal/hexagon/`) und in den Driven-Adaptern — übernimmt
  `slice-tests-ueberlebende-mutanten` (M3, P8, P2, Y2); er liegt seit seiner Closure in `done/`.
- Die Aufzeichnung von Tabellen-OID und Spaltennummer (Verifikation F-2 zu
  `slice-tests-ueberlebende-mutanten`) — Schicht-Abgrenzung: Der Test läge im Driven-Adapter
  `recording`, mit PGWire-Adapter und Bootstrap die dritte Schicht. Es ist keine Adresse: Ein Mutant dort
  ist nicht gemessen, und `slice-harness-mutation` schließt Tests für heute überlebende Mutanten aus.
  Der Stand steht im Register (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`); ob ein eigener Slice
  entsteht, entscheidet der Nutzer.
- Neue Tests für G2, G3, G5 und G6 — Bestand bleibt bewusst stehen: Die Messung oben zeigt sie
  gefangen; ein zweiter Test für denselben Mutanten prüfte keine weitere Zusage.
- Das Ende einer Replay-Sitzung durch die Frist des Herunterfahrens und `PGR-E4006` —
  übernimmt `slice-v1-abschluss-herunterfahren` (§1 dort, `done/`); dieser Slice berührt es nicht.
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
  hierher, dieser Slice ändert sie nicht; die Berichtigung zu G2, G3, G5 und G6 steht hier.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Je Liefer-Punkt eine Gruppe aus §1 *Gegenstand*, nach Lieferwert geschnitten; die
Sendung kommt von `slice-lint-bestand-driving` (§7 dort, G1, G4 und G7; G2, G3, G5 und G6
entfallen nach §1 *Aktualität*). Der Beleg des Implementers steht in §7 als
Mutationstabelle (Zusage · Mutation · roter Test, je Mutant dazu der Lauf ohne den neuen
Test, in dem er grün bleibt), nicht im Bericht; vor dem Code misst er G1, G4 und G7 gegen den
Stand und nennt, ob sie überleben. Entscheidet der Architect zu einem Mutanten, dass keine
Zusage besteht (§6; entschieden ist das nur für G7), entfällt dessen Test mit Begründung in §7; der Mutant
ist ein akzeptiertes Negativ mit diesem Grund, nicht „äquivalent“: Er ändert das Verhalten, und
Invariante, Ort und Sonde (`AGENTS.md` §3.10) gibt es nicht.

- [ ] **Frist beim Warten (F-468):** Die Tests dieser beiden Schichten, die ohne Frist auf
      einen Kanal oder das Ende eines Laufs warten, tragen die Frist nach §6 *Warten auf
      einen Kanal* (Kandidaten: `server_extended_test.go:848`, `server_tls_test.go:54`,
      `bootstrap_test.go:353`, `tls_test.go:184`; welche der Architect festlegt). Der Mutant W
      des Gebers wird in `TestReplayHerunterfahrenSpaeteFrist` rot, ohne dass das Paket die
      Zeit überschreitet.
- [ ] **Abbildung (G4) und Ende des Wächters (G7):** Ein Test in `pgwire_test` über
      `pgwire.ToMessage` hält Tabellen-OID und Spaltennummer der Spaltenbeschreibung
      nach §6 *Spalten-Metadaten* (`LH-FA-02.a`, `LH-FA-03.a`, `LH-FA-18.a`) und wird
      rot unter T3 und, je für sich, unter T5. Für G7 kein Test: §6 *Ende des Wächters*
      entscheidet „keine Zusage“; §7 führt den Mutanten als akzeptiertes Negativ (nicht äquivalent) mit diesem Grund.
- [ ] **Öffnen am Rand (G1):** Ein Test in `bootstrap_test` über `bootstrap.Run` hält die
      Zusage nach §6 *Signal vor dem Öffnen* (`LH-FA-13.a` *Startphase*) und wird rot ohne
      `context.WithoutCancel` in `pgwire.Listen`; geprüft nur unter Linux (Grenze für macOS und Windows steht in `LH-FA-13.a`).
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
| `spec/spezifikation.md` | update (Architect, vor dem Code) | Die offenen Randformen aus §6 (Ende des Wächters, Signal vor dem Öffnen, Spalten-Metadaten im Replay, Stellen der Frist), an der Stelle, die der Architect wählt; die Lesart *Spalten-Metadaten* für den Rückweg steht seit `slice-tests-ueberlebende-mutanten` in `LH-FA-02.a` |
| `internal/adapters/driving/pgwire/server_extended_test.go` | update | Frist beim Warten auf `conn.gesetzt` (Zeile 848 am Stand `f4ff9387`) und an weiteren Stellen dieser Schicht nach §6 *Warten auf einen Kanal* (Kandidat `server_tls_test.go:54`) |
| `internal/adapters/driving/pgwire/*_test.go` (Abbildung, Wächter) | update | G4: `ToMessage` mit RowDescription, alle Felder verglichen (Happy, LH-FA-18); G7: kein Test (keine Zusage, §6) |
| `internal/bootstrap/bootstrap_test.go` | update | G1: `bootstrap.Run` mit beendetem Kontext und `--listen localhost:0` (Boundary, LH-FA-13), dazu die zwei Startfehler derselben Zusage; Frist beim Warten auf das Ende des Laufs (`bootstrap_test.go:353`, `tls_test.go:184` und, nach derselben Entscheidung des Architect für den Lauf der Kommandozeile, `tls_test.go:100`) |
| `internal/adapters/driving/pgwire/server.go` | update (nur Kommentar) | Doc-Kommentar von `replayWaechter` (der gestrichene Satz, nur mit dem Test aus G7) und, wo ein benannter Test sie prüft, von `replayLesefehler` und `replayZustellen` |

Die Abdeckungs-Deklarationen bleiben, wie sie sind: Die Tests decken neue Fälle
bestehender Anforderungen und Pfade ab. Ob eine Zeile in den Tabellen unter
`docs/user/` hinzukommt, entscheidet `make abdeckung-check`; wenn ja, gehört sie in
denselben Commit.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-tests-ueberlebende-mutanten` liegt in
`done/` oder ist ausdrücklich zurückgestellt (WIP-Limit 1). Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness, M3 vor M4):
die drei Slices der Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`; Entscheidung des Nutzers vom 2026-10-13, vor die Release-Slices gezogen), dann die Slices von welle-v1-abschluss, darin `slice-harness-integration-wait` direkt vor
`slice-v1-abschluss-herunterfahren` und `slice-harness-meldungskatalog-gate` direkt vor
`slice-v1-abschluss-container`, dann `slice-harness-abdeckung-gate` und
`slice-harness-coverage`, dann die Slices von welle-erster-release und welle-v1-nachtrag, danach
`slice-harness-commit-struktur-id`, `slice-tests-ueberlebende-mutanten`,
`slice-tests-ueberlebende-mutanten-driving`, `slice-harness-lint-warten-ohne-frist`,
`slice-harness-mutation`; die Reihenfolge
der Wellen-Slices steht in §5 von [welle-v1-abschluss](../done/welle-v1-abschluss/welle-v1-abschluss.md).

*Warum hier.* Technisch hängt der Slice an keinem der beiden Nachbarn. Direkt nach
`slice-tests-ueberlebende-mutanten`, weil beide nach derselben Regel sammeln und der
Architect die Randform *Spalten-Metadaten* dort für den Rückweg mitentscheidet. Vor
`slice-harness-mutation`, damit dessen Messung nicht mit Lücken startet, die schon
bekannt sind, und damit die Randform *Laufzeit* dort die Entscheidung *Warten auf einen
Kanal* von hier vorfindet. Vor `slice-v1-abschluss-herunterfahren` ist er nicht gebunden; nach der
Reihenfolge vom 2026-10-08 liegt jener früher in `done/`, und der Architect prüft G2 gegen
dessen Ende durch die Frist und jeden gesammelten Mutanten gegen den Code nach den Wellen.

**Start-Trigger und Ausnahme (Anspruch am 2026-10-24):** Der Nehmer `slice-tests-ueberlebende-mutanten`
liegt in `done/`, die Bedingung oben ist erfüllt. Der Start-Trigger von
welle-gegenproben-und-mutanten (`welle-harness-regeln` in `done/`) ist es nicht; der Nutzer hat am
2026-10-24 mit „ok deine Empfehlung“ (in dieser Sitzung gesprochen) den Anspruch vor der Welle als
Ausnahme angenommen, wie beim Geber (Drift-Log der Roadmap, 2026-10-24). Der Slice hängt technisch an
keinem Harness-Slice. Die Messung der Mutanten (§1 *Aktualität*) steht am Stand `f4ff9387`; vor dem
Code misst der Implementer G1, G4 und G7 erneut und nennt in §7, ob sie überleben (Y2).

**Rückführungs-Marke** (`AGENTS.md` §3.14): Marke 1500, Metrik und Befehl dort. Pfade des Plans:
`internal/adapters/driving/pgwire`, `internal/bootstrap`, `spec`. Basis-Hash des Moves nach
`in-progress/`: `7ed55a23`.

Erster Schritt nach dem Start, vor jedem Code-Commit: Der Architect entscheidet die
Randformen aus §6 und hält sie in der Spezifikation fest
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`). Davor, noch in `open/`, hat
der Planner die Sammlung geprüft (2026-10-24): drei Liefer-Punkte, zwei Schichten (§8
*Nachzählung*), kein weiterer Schnitt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Ein Test braucht einen
  Zugang, den die Export-Test-Brücke nach `SPEC-049` Punkt 7 nicht hergibt, oder die
  Frist nach §6 *Warten auf einen Kanal* trifft so viele Tests dieser Schicht, dass sie
  mit den Mutanten nicht in eine Review-Sitzung passt; dann wird der betroffene Teil
  ein eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): Ein Test macht einen Fehler im
  Produkt sichtbar statt einer Lücke im Test, etwa wenn der Lauf mit beendetem Kontext
  anders enden muss, als er endet (G1). Dann entscheidet der
  Architect über eine Korrektur des Codes in einem eigenen Slice; der Test wartet auf
  sie.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; jeder Mutant aus §1 *Gegenstand* ist in §7 mit seiner Mutation rot
belegt oder (G7) als akzeptiertes Negativ mit Begründung geführt, die vier entfallenen (G2, G3, G5, G6)
stehen mit ihrem Lauf am Endstand in §7, und die Verifikation hat das nachgefahren. `make gates` ist grün, die Closure-Notiz mit Lerneintrag ist geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12). Der Slice liefert keinen neuen Vertrag, aber jeder
Test macht eine Zusage prüfbar. Wo die Zusage in keinem Stratum steht, entscheidet der
Architect vor dem ersten Code-Commit und hält sie in der Spezifikation fest. Was dort
nicht steht, entscheidet der Implementer nicht; er gibt es zurück. Die Randformen zum
Verbindungsende, zur unlesbaren Nachricht und zum Versand im Replay (G2, G3, G5, G6)
entfallen mit ihren Mutanten (§1 *Aktualität*); die Bestandstests, die sie fangen, tragen
ihre Zusagen.

Stand der Randformen (Architect, 2026-10-24, vor dem ersten Code-Commit; Randform ·
Entscheidung · Ort):

- **Ende des Wächters** (G7) — **entschieden: keine Zusage** (`SPEC-032`, Absatz *Nicht
  zugesichert ist, was eine beendete Verbindung im Prozess zurücklässt*). Der Wächter darf
  nach dem Ende der Sitzung bis zum Ende des Prozesses weiterlaufen; eine Goroutine je
  Replay-Verbindung kostet in einem Testwerkzeug nichts Messbares, und ein Test über
  `runtime.NumGoroutine` wäre flatterhaft (kein Parallellauf, Frist) und belegte keine
  Zusage. Es entfällt der Test aus G7. **Der Mutant ist nicht äquivalent**: Er ändert das
  Verhalten (die Goroutine bleibt), `AGENTS.md` §3.10 verlangt Invariante, Ort und Sonde
  nur für „äquivalent“, und die gibt es hier nicht. §7 führt ihn als **akzeptiertes
  Negativ** mit dem Grund „keine Zusage in `SPEC-032`“, ohne Sonde. Der Satz im
  Doc-Kommentar von `replayWaechter` kehrt nicht zurück.
- **Spalten-Metadaten** (G4) — **entschieden**: Record in `LH-FA-02.a` (Absatz
  *Spalten-Metadaten*), Replay in `LH-FA-03.a` (Absatz *Spalten-Metadaten*, neu): Tabellen-OID
  und Spaltennummer einer aufgezeichneten `RowDescription` gehen unverändert an den Client,
  einfach wie Extended. Der Test über `pgwire.ToMessage` trägt beide Zusagen (der Weg ist
  für Record und Replay derselbe). Die Aufzeichnung der Felder bleibt ausgeschlossen (§1).
- **Signal vor dem Öffnen** (G1) — **entschieden, eine Zusage bestand schon**: `LH-FA-13.a`
  *Randfälle des Herunterfahrens*, *Startphase* („Ein Signal, bevor der Prozess lauscht,
  bricht den Start nicht ab“; der Planner las „keine Stelle sagt es“, die Stelle sagt
  Exit-Code `0` mit der Ausnahme, dass ein Startfehler seinen Code behält). Der Architect
  hat dort ergänzt, dass das Öffnen des Listeners samt Namensauflösung zu den
  Startprüfungen gehört: Exit-Code `0`, kein `PGR-E4001`; ein Name, der nicht auflöst, und
  ein Port, der nicht zu öffnen ist, bleiben Startfehler. Der Test in `bootstrap_test` ist
  Pflicht, der Ausgang „keine Zusage“ in §2 entfällt. Eine Änderung des Verhaltens ist
  nicht vorgeschlagen und keine Entscheidung des Nutzers nötig. **Plattform** (`LH-QA-03`:
  Linux, macOS, Windows je `amd64` und `arm64`): geprüft nur unter Linux, in der Stufe
  `test` löst die Hosts-Datei `localhost` auf; macOS und Windows sind Grenze, im Text von
  `LH-FA-13.a` an der Zusage genannt, kein Test dort (die Gates laufen nur im
  Linux-Container). Der Implementer fährt die Mutation in einer Kopie und prüft, dass der
  Test ohne `context.WithoutCancel` rot wird; löst `localhost` in der Stufe `test` nicht
  auf, hält er an und gibt es zurück (kein anderer Name, keine Annahme).
- **Warten auf einen Kanal** (Review F-468) — **entschieden** in `SPEC-038` *Warten in
  Tests* (Regel, Absatz zum Lauf der Kommandozeile und zum Kanal des Testgeschirrs, neu).
  Je Kandidat:

  | Stelle | Entscheidung |
  |---|---|
  | `server_extended_test.go:848` (`<-conn.gesetzt`) | **Frist Pflicht.** Den Wert auf `gesetzt` führt der Prüfling herbei (Lesefrist durch Wächter oder Zwangsende); die Freigabe `close(conn.freigabe)` davor ist das Double und fällt nicht darunter. |
  | `server_tls_test.go:54` (`server = <-angenommen`) | **Keine Pflicht, Frist erlaubt und zu setzen.** Das `Accept` auf dem Listener des Tests nach dessen eigenem `Dial` ist Testgeschirr, der Prüfling führt es nicht herbei. Der Implementer setzt dennoch eine Frist: `slice-harness-lint-warten-ohne-frist` meldet jeden Empfang im Rumpf einer Hilfsfunktion (§1 dort, Entscheidung des Nutzers) und kann das Geschirr nicht unterscheiden; zwei Zeilen jetzt sparen dort eine Ausnahme. |
  | `bootstrap_test.go:353` (`exit := <-ende`) | **Frist Pflicht.** Der Lauf der Kommandozeile (`bootstrap.Run`) ist der Einstiegspunkt des Prozesses; der Test startet ihn in einer Goroutine und liest sein Ende über einen Kanal. Das ist ein gestarteter Prozess im Sinn der Regel und zugleich ein Wert, den der Prüfling auf einem Kanal liefert; beide Lesarten führen zur Frist. |
  | `tls_test.go:184` (`<-ende`, Paket `bootstrap`) | **Frist Pflicht**, Gründe wie `bootstrap_test.go:353`. |

  Für die zwei Stellen am Ende von `bootstrap.Run` gilt dazu „je Prozess genau eine Stelle
  wartet, jeder weitere Leser liest ihr Ergebnis“: Liest ein Test das Ende desselben Laufs
  an mehr als einer Stelle, wartet eine mit Frist und die andere liest ihr Ergebnis; findet
  der Implementer eine solche Stelle, trägt er es in §7 ein. Die Entscheidung gilt für
  beide Schichten dieses Slice. Im Kern (`record_extended_test.go:529`, `:812`) und in
  `test/integration` setzt dieser Slice keine Frist (Adressen wie oben:
  `slice-harness-lint-warten-ohne-frist`, `slice-harness-integration-wait`). Die Zählung
  eines Mutanten, der nur über die Zeitüberschreitung rot wird, bleibt die Randform
  *Laufzeit* von `slice-harness-mutation`.
- **Zugang der Tests** — **entschieden** in `SPEC-049` Punkt 7: Die Tests liegen in
  `pgwire_test` und `bootstrap_test` und gehen über die exportierte Schnittstelle; die
  Brücke wächst nicht.
- **Was ein Test festhält** — **entschieden** (Architect, 2026-10-24, bei
  `slice-tests-ueberlebende-mutanten`) in `SPEC-038`, Absatz *Meldungstext in Tests*: Code und
  zugesagte Textteile, nicht der ganze Text; sie gilt für beide Slices.

**Risiken:**

- **Die Sammlung wächst über den Schnitt** — weitere Funde (§1 *Sammelregel*), etwa aus
  `slice-harness-lint`, heben die Liefer-Punkte über drei. Seit dem Anspruch bekommt
  ein neuer Fund einen eigenen Slice (§1 *Sammelregel*). — **Ausgang:** — (bei Closure)
- **Ein Test hängt am Zeitverhalten** — die Frist aus F-468 (G7 entfällt, §6) beobachtet über
  eine Frist; zu knapp gesetzt macht sie den Test flatterhaft, zu weit lässt sie einen
  roten Lauf lange hängen (`BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`, 2×;
  `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients`, 1×, verwandt). Je Test mindestens 50 Läufe ohne Mutant grün, in §7. — **Ausgang:** — (bei
  Closure)
- **Der Test hält den Bestand statt der Zusage fest** — ein Vergleich des ganzen
  Meldungstextes machte jede Umformulierung rot, ohne dass eine Zusage bricht. Je Test
  nur das, was der Architect nach §6 *Was ein Test festhält* entscheidet. —
  **Ausgang:** — (bei Closure)
- **Die Aktualität steht auf einem Lauf je Mutant** — §1 *Aktualität* misst G1, G4 und G7
  grün und G2, G3, G5, G6 rot mit je einem Lauf; ein flatternder Test könnte eine Zeile
  verfälschen (Y2 des Gebers: gesichtet statt gemessen). Der Implementer misst G1, G4, G7
  vor dem Code und die vier entfallenen am Endstand erneut, in §7. — **Ausgang:** — (bei
  Closure)
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
`docs/plan/planning/observations/BEO-REPO/` am Stand `f4ff9387` gesichtet (Zähler =
Dateien unter `evidence/`, neu gelesen beim Anspruch 2026-10-24). Treffer:

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (97×, verkörpert in `AGENTS.md`
  §3.11; Sensor geplant mit `slice-harness-mutation`) — die nach F-465 und V-93
  gestrichenen Sätze kehren nur mit einem benannten Test zurück (§1), und je Test gilt das
  Risiko *Der Test hält den Bestand statt der Zusage fest* (§6).
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (65×, verkörpert in `AGENTS.md`
  §3.10) — G1, G4 und G7 sind Lücken dieser Klasse an Verträgen, die vor §3.10
  entstanden; der Slice schließt sie, er schärft die Regel nicht.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (51×, verkörpert in `AGENTS.md`
  §3.12) — die offenen Randformen stehen deshalb vor dem Code in §6.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (14×, verkörpert in `AGENTS.md` §3.13) —
  dieser Slice ist die Adresse; §1 nennt den Geber, und die Sendung steht in §1 und DoD;
  die Berichtigung zu G2, G3, G5 und G6 steht ebenfalls dort.
- `BEO-REPO/gruener-mutant-faelschlich-aequivalent` (11×, verkörpert in `AGENTS.md` §3.10) —
  trifft G7: Der Architect hält keine Zusage fest (§6), §7 führt G7 als akzeptiertes
  Negativ und nicht als äquivalent, weil Invariante, Ort und Sonde fehlen; G1 hat eine Zusage und einen Test.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — der Slice ist wellenlos;
  §4 *Start* nennt den Architect-Schritt vor dem Code.
- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (1×) — G7 und die Frist aus
  F-468 beobachten über eine Frist; Risiko in §6.
- `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — ein Risiko in §6.
- `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` (3×) — die Randform *Warten auf einen
  Kanal* in §6 ist seine Adresse; wartet ein neuer Test ohne Frist, trifft er die Klasse
  ein viertes Mal.
- `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (4×) — trifft diesen Slice nur, wenn
  Review oder Verifikation hier einen grünen Mutanten finden, den §7 nicht führt.

Keiner der Einträge unter der Schwelle erreicht mit diesem Slice allein 3×; vor dem Code
entsteht keine neue Lücke.

**Nachzählung** (Planner, 2026-10-24, Anspruch, Stand `f4ff9387`): Liefer-Punkte 3
(Frist · G4 · G1; G7 ist kein Liefer-Punkt, es gibt keinen Test), Schichten 2 (PGWire-Adapter, Bootstrap; Tests zählen als die
Schicht ihres Codes; die Spezifikation schreibt der Architect vor dem Code). Gegenüber der
Fassung vom 2026-10-07 entfallen vier Mutanten (§1 *Aktualität*), die Frist wird zu einem
eigenen Punkt, der Plan wird kleiner, nicht größer. Fügt `make abdeckung` Zeilen unter
`docs/user/abdeckung-*.md` hinzu, wären sie formal eine dritte Schicht; der Geber
`slice-tests-ueberlebende-mutanten` entschied das für sich als Einzelfall (§8 dort,
Planner), der Nutzer hat den Satz für `AGENTS.md` §3.13 nicht entschieden, und hier gilt
dasselbe wie dort: Die Zeilen kommen nur, wenn ein neuer Test eine Deklaration trägt, die
die Tabelle noch nicht führt, und zählen nicht als Schicht. Die Aufzeichnung der
Spalten-Metadaten bleibt draußen (§1), sie wäre die dritte Schicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
