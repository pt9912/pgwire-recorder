# Slice slice-v1-abschluss-anmeldung: Anmeldung des Clients im Record-Modus vermitteln

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md), [ADR-0041](../../adr/0041-anmeldung-im-record-als-rueckruf-des-upstream-ports.md) (Accepted, Entscheidung des Nutzers vom 2026-10-14)

**Berührte Spec-Stellen:** `LH-FA-05.b` · `LH-FA-02.a` · `SPEC-034` · `ARC-006` · `ARC-007` · `architecture.md §4.4`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `record` vermittelt die Anmeldung des Clients beim Upstream transparent (Klartext-Passwort, MD5, SCRAM-SHA-256), sodass Anwendungen gegen Server mit Passwort-Anmeldung aufgezeichnet werden können.

**Übernommen aus `slice-harness-blackbox-driven`** (Review F-444, Verifikation V-85; Abgrenzung in `slice-tests-ueberlebende-mutanten`): ein Test, dass die Session aus `Open` den Lesepuffer des Aufbaus trägt, mit Test-Idee und Grenze in §6; er gehört zu DoD-Punkt 1, weil der Umbau der Aufbau-Schleife in `Open` genau diese Stelle berührt.

**Übernommen aus `slice-doku-ist-stand`** (Verifikation V-152, Review F-573; bekannter Teilstand, keine neue Abweichung): Der Meldungstext des Binaries für ein Anmeldeverfahren, das `record` nicht vermittelt („der Upstream verlangt ein Anmeldeverfahren (…), das dieser Stand nicht vermittelt“, `internal/adapters/driven/postgres/upstream.go`), wird beim Liefern von SCRAM-SHA-256 so gefasst, dass er nur noch die Verfahren nennt, die `record` weiterhin nicht vermittelt; er gehört zu DoD-Punkt 2.

**Randformen** sind vor dem Code entschieden (Architect, 2026-10-14, `AGENTS.md` §3.12): [`LH-FA-05.b`](../../../../spec/spezifikation.md#lh-fa-05b--authentifizierung) *Randformen der Vermittlung (Record)*. Der Slice liefert genau sie; jede Aufzählung und Ortsangabe dort zählt je Element (§6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-anmeldung-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- Anmeldung beim Einspielen (`play`) — `slice-v1-abschluss-einspielen-anmeldung` (dort §1, *Übernimmt*; seit dem Schnitt von `slice-v1-abschluss-einspielen` vom 2026-10-09); `play` rechnet den Austausch selbst, `record` vermittelt ihn.
- Client-Zertifikate und GSSAPI/SSPI — außerhalb des Funktionsumfangs von v1 (`LH-FA-05.b`); sie enden als nicht vermitteltes Verfahren (`PGR-E6001`, DoD-Punkt 2).
- Andere SASL-Mechanismen als SCRAM-SHA-256 und Channel Binding — nicht zugesagt (Grenze in `LH-FA-05.b`): Der Recorder liest die Liste der Mechanismen nicht und terminiert TLS nur zum Client.
- Eine Option, Umgebungsvariable oder ein Schlüssel für die Anmeldung — `record` vermittelt, die Zugangsdaten gehören dem Client (`LH-FA-05.b`); `zuordnung()` und `nichtGelesen()` im CLI-Adapter bleiben unberührt, dort entsteht kein Code.
- Anmeldung gegen jede unterstützte PostgreSQL-Version — `slice-v1-abschluss-postgres-versionen` (dort §1 *Anmeldung je Version* und DoD-Punkt 2; Entscheidung des Nutzers vom 2026-10-07): Dessen Matrix fährt den Fall SCRAM-SHA-256 des Anmeldetests aus diesem Slice gegen jede Version außer 17 und startet erst, wenn dieser Slice in `done/` liegt; hier läuft der Anmeldetest gegen die Referenzversion 17. Der Test muss deshalb als Fall SCRAM-SHA-256 einzeln wählbar sein.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-05`](../../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol): Ein Client meldet sich über `record` an einem Server mit Klartext-Passwort, MD5 und SCRAM-SHA-256 an und führt eine Anfrage aus; Passwort und Anmeldenachrichten stehen nicht in der Aufzeichnung (Integrationstest).
- [x] Eine fehlgeschlagene Anmeldung geht als Fehlerantwort des Servers unverändert an den Client; ein nicht vermitteltes Verfahren endet mit eindeutigem Meldungscode, und sein Meldungstext sagt nicht mehr „das dieser Stand nicht vermittelt“, sobald SCRAM-SHA-256 vermittelt wird (aus `slice-doku-ist-stand`, V-152; Test).
- [x] Der Port zwischen PGWire-Adapter, Record-Service und Upstream-Adapter trägt den Anmeldeaustausch, ohne dass der Core PGWire-Typen kennt (`make a-check`).
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
| `internal/hexagon/model`, `internal/hexagon/ports/driving`, `…/driven`, `internal/hexagon/services` | update | Rückruf der Anmeldung (`model.Anmeldung`, `Anmeldeanforderung`, `ErrClientEnde`): Anforderung und Antwort als für den Kern undurchsichtige Werte zwischen Client und Upstream ([ADR-0041](../../adr/0041-anmeldung-im-record-als-rueckruf-des-upstream-ports.md)) |
| `internal/adapters/driving/pgwire`, `internal/adapters/driven/postgres` | update | Upstream-Adapter: Reihenfolge des Servers und Weitergabe der Anforderungen (`vermittlung.go`); Meldungstext zu einem nicht vermittelten Verfahren (aus `slice-doku-ist-stand`, V-152). PGWire-Adapter: Austausch mit dem Client (`anmeldung.go`; die Nachricht `p` liest `SetAuthType(AuthTypeSASLContinue)` als `SASLResponse`, die jeden Rumpf unverändert nimmt); im PGWire-Adapter endet der Wächter mit der Startnachricht, damit die Lesefrist des ersten Signals für die Anmeldung nicht gilt (`LH-FA-05.b` *Abbruch und Herunterfahren*; Grenze: eine im Rennfenster schon gesetzte Frist ist ein akzeptiertes Negativ der Spezifikation) |
| `test/integration`, `tools/test/run-integration-tests.sh`, `docs/user/abdeckung-*.md` | update | Server mit `password`, `md5` und `scram-sha-256` in `pg_hba.conf` (die Benutzer legt der Runner schon an), dazu die Zeile `reject` für `play_reject`; Abdeckungstabellen über `make abdeckung` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-walking-skeleton-record` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Der Anmeldeaustausch verlangt einen anderen Zuschnitt des Record-Ports als den Rückruf in `Open` ([ADR-0041](../../adr/0041-anmeldung-im-record-als-rueckruf-des-upstream-ports.md)) — zurück zur Zerlegung. **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- internal test docs/user`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 2500** geht der Slice zurück (Architect, 2026-10-14: die Marke des Vorgängers `slice-v1-abschluss-tls-client`, den ein Review in einer Sitzung trug; gemessen dort 2351). **Meldeschwelle** neben der Marke: Erreicht die Summe **2000**, meldet der Implementer den Stand dem Architect, bevor er weiterarbeitet; über Schnitt oder Anhebung der Marke entscheidet der Planner. **Schätzung** (Architect, 2026-10-14, in dieser Metrik): Produktcode rund 400 Zeilen (Port und Modell 50, Record-Service 50, PGWire-Adapter 140, Upstream-Adapter 160); Tests nach dem Faktor 2 rund 800 (ein Fall und eine Mutation je Element der Aufzählungen in `LH-FA-05.b`, §6) und der Integrationstest samt Server-Konfiguration rund 290; zusammen rund 1500, mit dem Aufschlag der letzten Code-Slices (20 bis 30 %) und der Review-Runde (20 bis 30 %) **1800 bis 2400**. Schnitt dann nach Lieferwert: Klartext und MD5 mit Port und Fehlerpfaden, danach SCRAM-SHA-256 (`SetAuthType` für SASL, drei Schritte).
- `in-progress` → `open`: Die Bibliothek ([ADR-0010](../../adr/0010-verwendung-von-pgproto3.md)) liest die Antworten der Anmeldung (`p` als `SASLInitialResponse` und `SASLResponse`) oder die Fortsetzungen (Code 11, 12) nicht, ohne dass der Slice sie selbst dekodiert — Entscheidung des Nutzers über die Wahl der Bibliothek, bevor weitergearbeitet wird. (Die frühere Bedingung „Kanalbindung“ entfällt: `LH-FA-05.b` sagt Channel Binding nicht zu.)

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

**Randformen** (`AGENTS.md` §3.12) — entschieden in [`LH-FA-05.b`](../../../../spec/spezifikation.md#lh-fa-05b--authentifizierung) *Randformen der Vermittlung (Record)*, Architect 2026-10-14: fehlender, leerer und unbekannter Wert (leere Antwort `p`, unbekannter Code), Alias (keiner), mehrfache Angabe (zweite Anforderung, Antwort vor der Anforderung), Aufrufart (nacheinander, nie gleichzeitig; keine Option, `zuordnung()` und `nichtGelesen()` unberührt), Abbruch (Ende durch den Client; erstes Signal, Zwangsende), Verhalten je Serverversion, Schutz der Zugangsdaten, Reihenfolge am Start (nichts Neues). **Je Element ein Fall und eine Mutation, die nur es bricht** (`AGENTS.md` §3.10): die Codes 2, 6, 7, 9 und ein unbekannter (`PGR-E6001`); die Reihenfolge-Fälle des Servers (`PGR-E4002`); `p`, `Terminate`, Ende des Datenstroms, Reset und eine andere Nachricht des Clients in der Anmeldung; Klartext, MD5, SCRAM-SHA-256; Fehlerantwort am Anfang und mitten im SCRAM-Austausch; jeder Ort ohne Zugangsdaten (Aufzeichnung, jede Log-Stufe, Meldung); erstes Signal und Zwangsende je beim Warten auf den Client in der Anforderung und in der Antwort. Die Meldung zu einem nicht vermittelten Verfahren nennt den Code, nicht den Typ der Bibliothek. Nachentschieden (Architect, 2026-10-14, auf die Randform-Rückgabe des Implementers): ein Code ohne Verfahren (2, 6, 7, 9, unbekannt) nach `AuthenticationOk` ist `PGR-E4002` (Stelle vor Code; ein Fall je Code); ein zweites Code 11 im laufenden SASL-Austausch ist `PGR-E4002` (Code und Fall); ein Schreibfehler an den Client bei einer Anforderung ist Ende durch den Client ohne Meldung (Fall und Mutation, die ihn als `PGR-E4003` einstuft); ein Schreibfehler an den Server bei einer Antwort ist `PGR-E4003` (Fall: der Server setzt die Verbindung zurück, der Rückruf des Tests wartet darauf, bevor er antwortet; ist das nicht stabil zu erreichen, steht es hier als akzeptiertes Negativ mit diesem Grund). Verhalten: Der Recorder nimmt eine schon gesetzte Lesefrist zurück; akzeptiertes Negativ (`LH-FA-05.b` *Abbruch und Herunterfahren*): Dass sie im Rennfenster zwischen Startnachricht und Ende des Wächters zurückgenommen wird, prüft kein Test, weil das Fenster nicht deterministisch zu treffen ist. Mehraufwand in der Metrik rund 130 bis 170 Zeilen (gemessen 1701, unter der Meldeschwelle 2000).

- **Größe, Schichten (Entscheidung des Nutzers vom 2026-10-14):** Der Plan berührt drei Schichten (Kern mit Ports, PGWire-Adapter, Upstream-Adapter) und liegt über der Grenze von zwei. Ein Schnitt nach Lieferwert (Klartext und MD5, danach SCRAM-SHA-256) senkt die Zahl nicht: Jeder Teil braucht Port, Client-Seite und Server-Seite des Austauschs; ein Schnitt nach Schicht erzeugt Zombie-Slices. Entscheidung des Nutzers vom 2026-10-14: Der Plan läuft ungeteilt mit drei Schichten (Ausnahme von der Zwei-Schichten-Regel), mit Marke und Meldeschwelle aus §4; der Schnitt nach Lieferwert bleibt die Rückführungs-Option bei Überschreiten der Marke 2500. Keine Log-Zeile bei fehlgeschlagener Anmeldung, kein neuer Katalog-Code (Entscheidung des Nutzers vom 2026-10-14, `LH-FA-05.b`) — **Ausgang bei Closure: entfallen.** Gemessen 1701 gegen die Marke 2500, die Meldeschwelle 2000 nicht erreicht, die Rückführung ist nicht eingetreten; keine Log-Zeile bei fehlgeschlagener Anmeldung (W21 rot) und kein neuer Katalog-Code (`internal/hexagon/model/fehler.go` steht nicht im Diff) halten (§7, Risiken aus §6).
- Die Vermittlung muss die Anmeldenachrichten aus der Aufzeichnung heraushalten, sonst stehen Passwort-Hashes darin (LH-RB-01) — **Ausgang bei Closure: entfallen.** Weder Aufzeichnung noch Log tragen Anforderung oder Antwort: S2, E4 und W18/W19 rot, `TestE2ERecordAnmeldung` prüft Aufzeichnung und Debug-Log auf die Spuren aller drei Verfahren (§7, Risiken aus §6).
- Dass die Session aus `Open` den Lesepuffer des Aufbaus trägt, prüft kein Test (`BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft`, seit slice-harness-blackbox-driven); der Umbau der Aufbau-Schleife berührt genau diese Stelle. Test-Idee: Der Server schickt AuthenticationOk, ReadyForQuery und eine NoticeResponse in einem Flush, der Test ruft nach `Open` einmal `Receive` und erwartet die Notice; rot, wenn `Open` eine Session mit neuem Frontend zurückgibt. Grenze: Dass beide Nachrichten in einem Lesevorgang ankommen, ist über TCP-Loopback stabil, aber nicht zugesagt; `net.Pipe` erreicht `Open` nicht, weil `Upstream.Dialer` ein `net.Dialer` ist — **Ausgang bei Closure: entfallen.** `TestOpenSessionTraegtLesepufferDesAufbaus` fängt U17 (§7, Risiken aus §6).

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

### Belege des Implementers

**Stand.** Geliefert sind DoD-Punkt 1 bis 3 und die beiden Übernahmen aus §1 (Lesepuffer des
Aufbaus, Meldungstext) in diesen Commits: `ce47692` (Rückruf in Modell, Ports und Service;
`vermittlung.go` im Upstream-Adapter; `anmeldung.go` und Wächter im PGWire-Adapter; Unit- und
Integrationstests; `reject` im Runner), `1de6493` und `a504fa4` (Testhilfe des Skript-Servers),
`6d74175`, `fddadc2`, `297c6fb`, `fcedfeb`, `4c38fb7` (weitere Fälle, Abdeckungstabellen über
`make abdeckung`), `962543d` (§3 dieses Plans). Neu sind `internal/hexagon/model/anmeldung.go`,
`internal/adapters/driven/postgres/vermittlung.go`, `internal/adapters/driving/pgwire/anmeldung.go`,
`internal/hexagon/services/record_anmeldung_test.go`,
`internal/adapters/driven/postgres/upstream_anmeldung_test.go`,
`internal/adapters/driving/pgwire/server_anmeldung_test.go` und
`test/integration/record_anmeldung_e2e_test.go`; geändert `upstream.go` (Aufbau in `aufbauen`),
`server.go` (Wächter in `startupBisSignal`, `oeffne` mit `be`, `ErrClientEnde`, Rücksetzen von
`SetAuthType`), die Ports, `record.go`, `tools/test/run-integration-tests.sh` (Zeile `reject` für
`play_reject`) und die Abdeckungstabellen unter `docs/user/`. Der CLI-Adapter, `zuordnung()`,
`nichtGelesen()` und `internal/bootstrap/bootstrap.go` sind unverändert. Handbuch und `README.md`
liefert der Doku-Folge-Slice.

**Randformen und Rückgabe.** Keine Randform aus §6 oder der Spezifikation entschieden, kein
Code-Commit ändert §6. Die vier Rückgaben des Implementers hat der Architect entschieden
(`468997b`, `LH-FA-05.b`); der Code folgt ihnen, jeweils mit Fall und Mutation:

| Entscheidung | Fall | Mutation | roter Test |
|---|---|---|---|
| Code ohne Verfahren (2, 6, 7, 9, unbekannt) nach `AuthenticationOk` ist `PGR-E4002`, vor ihm `PGR-E6001` | je Code ein Fall | Y1: Klasse vertauscht (nach Ok `PGR-E6001`) | `TestVermittlungCodeOhneVerfahrenNachOk` (2, 6, 7, 9, 4242) |
| zweites Code 11 im laufenden SASL-Austausch ist `PGR-E4002` | Reihenfolge-Zeile `[10 11 11]` | Y2: zweites 11 zugelassen | `TestVermittlungReihenfolge/[10_11_11]` |
| Schreibfehler an den Client bei einer Anforderung ist das Ende des Clients, ohne Meldung | `…EndeDesClients` (vor dem Lesen der Anforderung), `…Zwangsende` | W5 und Y3: als `PGR-E4003` eingestuft | beide Zeilen |
| Schreibfehler an den Server bei der Antwort `p` ist `PGR-E4003` | Server setzt die Verbindung zurück (Linger 0), der Rückruf wartet 300 ms | Y4: als `PGR-E4002` eingestuft | `TestVermittlungSchreibfehlerAmServer` |

**Größe.** Metrik aus §4: `git diff --numstat fa5a64c..<letzter Commit des Slice> -- internal test docs/user` ergibt
am Ende 1612 Einfügungen und 89 Löschungen, Summe **1701** (`4c38fb7`: 1623; nach der Entscheidung des Architect 1665; bei der Verifikation 1697; mit `01573f5` 1701). Die Meldeschwelle 2000 und die Marke 2500 sind nicht
erreicht; die Schätzung war 1800 bis 2400 (§4). Nach dem Review und der Verifikation sind nur Spezifikation, Plan und `01573f5` dazugekommen.

**Läufe.**

| Lauf | Stand | Ergebnis |
|---|---|---|
| `make gates` (voll, ein Lauf, kein pkill; deckt die Spezifikation `e1dc57b` mit ab) | `3e8a057`, letzter Stand nach der Verifikation (nur §7 danach), Metrik 1701 | Exit 0, kein FAIL |
| `make gates` (voll, ein Lauf, kein pkill) | `84ad1de`, Stand nach dem Review (nur §7 danach) | Exit 0, kein FAIL |
| `make gates` (voll) | `4c38fb7` | Exit 0; darin `run-integration-tests: gruen`, `d-check: 737 Datei(en) geprüft, 0 Befund(e)`, `0 issues.`, `lint-gegenprobe: gruen`, `a-check-negativ: gruen`, `abdeckung-gegenprobe: gruen`, `stop-hook-gegenprobe: gruen` |
| `make gates` | `962543d` | Exit 2: `TestE2ERecordAnmeldung/md5` rot, weil „md5“ im Benutzernamen `play_md5` der Aufzeichnung steht; die Spur ist seit `4c38fb7` `md5[0-9a-f]{32}`, der Lauf darüber grün |
| `make lint` | vor `ce47692` (nach dem Auszug von `startupBisSignal`) | Exit 0, `0 issues.`; davor 2 Befunde (`cyclop`, `gocyclo`: `handle` mit 16) |
| `make a-check` | `4c38fb7` | `gesamt: 0 Befund(e)`; Mutanten unten |
| `make abdeckung-check` | `fcedfeb` | Exit 0 |
| `make kopf-check`, `make docs-check` | nach dem Eintrag in §3 | Exit 0, 0 Befunde |

**Weg der Mutanten.** Je Mutant ein Skript (Scratch-Pfad der Sitzung), das genau ein Vorkommen im
Arbeitsbaum ersetzt (bricht bei fehlendem Vorkommen ab), die Tests des Pakets mit `-run` im
`deps`-Image des Dockerfile mit Bind-Mount des Arbeitsbaums fährt (kein Build-Kontext, die mtime ist
ohne Belang) und die Datei danach mit `git checkout --` zurücksetzt; ein Mutant der Integration läuft
über eine Kopie von `tools/test/run-integration-tests.sh` mit `-test.run '^TestE2ERecordAnmeldung'`
und `docker build --target integration`. Gefahren: 140 Mutanten am Quelltext (S 3, U 18, N 9, R 22, F 19, X 5, W 25 mit W3a und W3c, G 19, Z 5, Y 5, E 4, M 5, W18real 1) und 2 Architekturmutanten (eine
zusätzliche Datei im Kern, eine im Upstream-Adapter; `make a-check`), zusammen 142; 139 rot. Drei sind grün und äquivalent (W3, W3a, Y5, unten); Y1 und Y2 sind seit der Entscheidung rot.

**DoD-Punkt 1 — die Anmeldung wird vermittelt, Passwort und Anmeldenachrichten stehen nicht in der Aufzeichnung.**

| Zusage | Mutation | rote Tests |
|---|---|---|
| Klartext-Passwort (Code 3) wird vermittelt | U1: Code 3 in `folge` nicht zugelassen | `TestVermittlungVerfahren` (Klartext-Passwort, leere Antwort) und 8 weitere Zeilen mit Code 3 (6 von `…Reihenfolge`, Ende bei Klartext, `…FehlerantwortDesServers` nach Klartext) |
| MD5 (Code 5) wird vermittelt | U2: Code 5 nicht zugelassen | `TestVermittlungVerfahren` (MD5, Antwort mit NUL-Bytes), Reihenfolge `[5 11]`, `[5 3]`, `…FehlerDesRueckrufs` (PGR-E6001 des Clients) |
| SASL (Code 10) wird vermittelt | U3: Code 10 nicht zugelassen | `TestVermittlungVerfahren` (SCRAM-SHA-256) |
| SASLContinue (Code 11) wird vermittelt | U4: Code 11 nicht zugelassen | `TestVermittlungVerfahren` (SCRAM-SHA-256), 5 Zeilen Reihenfolge, 2 Zeilen `…FehlerDesRueckrufs`, `…FehlerantwortDesServers` (mitten im SCRAM) |
| SASLFinal (Code 12) wird vermittelt | U5: Code 12 nicht zugelassen; N9: Code 12 nicht erkannt | `TestVermittlungVerfahren` (SCRAM-SHA-256), Reihenfolge `[10 11 12 11]`, `[10 11 12 3]`, `…FehlerDesRueckrufs` (Ende bei Code 12) |
| AuthenticationOk ohne Anforderung kommt an | U6 | `TestOpenUndQuery`, Reihenfolge `[0 3]`, `[0 0]`, `[0 11]` |
| AuthenticationOk nach Klartext/MD5 kommt an | U7 | `TestVermittlungVerfahren` (Klartext-Passwort, MD5, leere Antwort, Antwort mit NUL-Bytes) |
| AuthenticationOk nach Code 12 kommt an | U8 | `TestVermittlungVerfahren` (SCRAM-SHA-256) |
| AuthenticationOk geht nicht an den Rückruf und nicht in die Antworten des Aufbaus | U13: Ok an den Rückruf | `TestVermittlungVerfahren` (alle 5 Zeilen), Reihenfolge `[0 3]`, `[0 0]`, `[0 11]` |
| Antwort des Clients unverändert als `p` an den Server, auch leer und mit beliebigen Bytes | U9: als `PasswordMessage` gesendet | `TestVermittlungVerfahren` (alle 5 Zeilen) |
| Code 12 geht an den Rückruf | U10: Rückruf für Code 12 übersprungen | `TestVermittlungVerfahren` (SCRAM-SHA-256), `…FehlerDesRueckrufs` (Ende bei Code 12) |
| auf Code 12 geht keine Antwort an den Server | U11: Antwort auch auf Code 12 | `TestVermittlungVerfahren` (SCRAM-SHA-256), Reihenfolge `[10 11 12 12]`, `[10 11 12 11]`, `[10 11 12 3]` |
| ein Fehler des Rückrufs bei Code 12 endet den Aufbau | U12: Fehler bei Code 12 übergangen | `TestVermittlungFehlerDesRueckrufs` (Ende bei Code 12) |
| ein Fehler des Rückrufs endet Open mit genau diesem Fehler | U15: Fehler verschluckt | 30 Zeilen: `…FehlerDesRueckrufs` (alle 5), Reihenfolge, `Code 7` |
| der Rückruf kommt vom Aufrufer des Service bis zum Upstream | S1: `nil` statt des Rückrufs | `TestRecordAnmeldungRueckruf` |
| Service: keine Session, wenn der Client in der Anmeldung endet (auch mit `LeereSessions`) | S3: Session trotz `ErrClientEnde` | `TestRecordAnmeldungRueckruf` |
| Anforderung unverändert (Code, Rest) an den Client | Z1: Rest der Nachricht weggelassen; Z2: Code immer 3 | `TestAnmeldungVermittelt` (MD5, SASL, SASL-Fortsetzung, Antwort von 70000 Bytes, SASL-Abschluss), je Mutant dieselben 5 Zeilen |
| Antwort des Clients unverändert an den Use Case | Z3: um ein Byte gekürzt; Z4: `nil` statt leerer Folge; W6: `p` als Passwort gelesen | Z3: Klartext, MD5, SASL, SASL-Fortsetzung, 70000 Bytes; Z4: leere Antwort; W6: 6 Zeilen von `TestAnmeldungVermittelt`, beide von `…ErstesSignal`, `…AndereNachricht`, `…FehlerantwortDesServers` |
| keine Höchstlänge der Antwort | W7: `SetMaxBodyLen(10000)` | `TestAnmeldungVermittelt` (Antwort von 70000 Bytes) |
| auf Code 12 erwartet der Adapter keine Antwort | W8: Antwort auch auf Code 12 erwartet | `TestAnmeldungVermittelt` (SASL-Abschluss ohne Antwort) |
| AuthenticationOk vor den übrigen Antworten des Aufbaus | W22: dahinter | `TestAnmeldungVermittelt` (7 Zeilen), beide von `…ErstesSignal`, `…AndereNachricht` |
| weder Passwort noch Anforderung stehen in der Aufzeichnung (Service) | S2: Antwort im `Startup` der Session | `TestRecordAnmeldungNichtInDerAufzeichnung` |
| … (Integration, Klartext, MD5, SCRAM-SHA-256) | E4: Antwort im `Startup` der Session | `TestE2ERecordAnmeldung` (scram-sha-256, md5, password); die Spuren `md5[0-9a-f]{32}`, `SCRAM`, `n,,n=` und `biws` fangen Antworten, die kein GEHEIM tragen |
| ein Client meldet sich über `record` an und führt eine Anfrage aus (alle drei Verfahren) | E1: Rückruf des PGWire-Adapters endet sofort | `TestE2ERecordAnmeldung` (alle 3), `…Fehlerantwort` (3), `…ClientBeendet` |
| kein Log einer Stufe nennt Anforderung oder Antwort (Integration: Spuren `GEHEIM`, `md5[0-9a-f]{32}`, `SCRAM`, `n,,n=`, `biws` auch im Debug-Log, alle drei Verfahren; W18real: Debug-Zeile mit der Antwort im Adapter, rot in scram-sha-256, md5 und password von `TestE2ERecordAnmeldung`) | W18: `Debug` mit der Antwort; W19: `Debug` mit der Anforderung | W18: 9 Zeilen (Klartext, MD5, SASL, SASL-Fortsetzung, Terminate, Ende des Datenstroms nach der Anforderung, Reset, Warten auf die Antwort, Fehlerantwort); W19: 10 Zeilen |
| eine Meldung nennt den Code, nie den Inhalt | X3: Reihenfolge-Meldung mit dem Inhalt in Hex; X4: Client-Nachricht mit `%+v` statt `%T`; N7: Text der Bibliothek in der Meldung | X3: 15 Zeilen Reihenfolge; X4: `Anfrage`, `Extended-Nachricht`; N7: Code 2, 6, 9, unbekannter Code |
| auf den Client wartet die Anmeldung ohne eigene Frist | X5: Lesefrist von 1 s nach der Startnachricht | `TestAnmeldungWartetOhneFrist` |
| erstes Signal beim Schreiben der Anforderung und beim Warten auf die Antwort unterbricht nicht | W1: Wächter läuft weiter und setzt die Lesefrist; W2: dasselbe mit der Schreibfrist | beide Mutanten rot in beiden Zeilen von `TestAnmeldungErstesSignal` (die Frist des Wächters wirkt auch auf das spätere Lesen bzw. Schreiben; ein Mutant, der nur eine der beiden Zeilen bricht, fand sich nicht) |
| Zwangsende beendet die Anmeldung beim Schreiben der Anforderung und beim Warten auf die Antwort | W3c: `z.ende` leer **und** `oeffne` kehrt beim Zwangsende nicht zurück; W5: Schreibfehler der Anforderung ist ein Verbindungsfehler; W4: Lesefehler nach dem Zwangsende nicht auf das Ende des Clients abgebildet | W3c: beide Zeilen von `TestAnmeldungZwangsende`; W5: `…Zwangsende` (Schreiben) und `…EndeDesClients` (vor dem Lesen); W4: `…Zwangsende` (Warten) |
| die Session aus `Open` trägt den Lesepuffer des Aufbaus (aus `slice-harness-blackbox-driven`) | U17: `Open` liefert eine Session mit neuem Frontend | `TestOpenSessionTraegtLesepufferDesAufbaus` |

**DoD-Punkt 2 — Fehlerantwort des Servers unverändert, nicht vermitteltes Verfahren mit eindeutigem Code.**

| Zusage | Mutation | rote Tests |
|---|---|---|
| Fehlerantwort am Anfang, nach Klartext, nach Code 10 und mitten im SCRAM geht zurück | U16: Fehlerantwort wird PGR-E6001 | `TestVermittlungFehlerantwortDesServers` (alle 4 Zeilen) |
| … gleich welcher SQLSTATE | X1: nur 28P01; X2: nur Klasse 28 | X1: 28000, 3D000, 53300, 57P03, XX000; X2: 3D000, 53300, 57P03, XX000 (`TestVermittlungFehlerantwortJederSQLSTATE`) |
| … mit jedem Feld im Upstream-Adapter | F1 bis F18 (je Feld von `noticeFields` entfernt: S, V, C, M, D, H, P, p, q, W, s, t, c, d, n, F, L, R), F99 (unbekannte Felder) | je Mutant alle 4 Zeilen von `TestVermittlungFehlerantwortDesServers` |
| … mit jedem Feld im PGWire-Adapter | G1 bis G18 (je Feld von `errorResponse` entfernt), G99 (unbekannte Felder) | je Mutant `TestAnmeldungFehlerantwortDesServers` |
| … wie bei einer Verbindung unmittelbar zum Server | E3: Feld `Routine` entfernt | `TestE2ERecordAnmeldungFehlerantwort` (scram-sha-256, md5, password, reject) |
| Fehlerantwort geht ohne AuthenticationOk an den Client | Z5: AuthenticationOk davor | `TestAnmeldungFehlerantwortDesServers` |
| Fehlerantwort geht überhaupt an den Client | W23: Versand entfernt | `TestAnmeldungFehlerantwortDesServers` |
| keine Log-Zeile bei fehlgeschlagener Anmeldung | W21: Log-Zeile im Pfad ohne Session | `TestAnmeldungFehlerantwortDesServers` |
| der Recorder sendet dem Server nach dem Abbruch weder eine Nachricht noch Terminate | U14: `Close` mit Terminate statt `conn.Close` | 40 Zeilen: alle von `…NichtVermittelt`, `…Reihenfolge`, `…FehlerDesRueckrufs`, `…FehlerantwortDesServers` |
| Code 2 ist PGR-E6001 | N1 | `…NichtVermittelt/Code 2` |
| Code 6 ist PGR-E6001 | N2 | `…/Code 6` |
| Code 7 ist PGR-E6001 | N4 (als Code in der Reihenfolge); N5 (Typ der Bibliothek in der Meldung) | `…/Code 7` |
| Code 9 ist PGR-E6001 | N3 | `…/Code 9` |
| ein unbekannter Code ist PGR-E6001 | N6 | `…/unbekannter Code` |
| eine nicht lesbare Nachricht des Servers im Aufbau ist PGR-E6001 | N8: PGR-E4003 | `…/nicht lesbare Nachricht`, `…/Nachricht R ohne Code` |
| der Meldungstext sagt nicht mehr „das dieser Stand nicht vermittelt“ (V-152) | U18: alter Text | `…NichtVermittelt` (Code 2, 6, 7, 9, unbekannter Code) |
| Reihenfolge des Servers: jede Zelle der Tabelle ist PGR-E4002 | R1 bis R22 und M7 bis M10 (je (Stand, Code) zugelassen; M7 bis M10 sind die Zellen (Code 11 beantwortet, 10), (… , 5), (Code 12 gelesen, 10), (… , 5) mit den Zeilen `[10 11 10]`, `[10 11 5]`, `[10 11 12 10]`, `[10 11 12 5]`); Zellen mit Code 3 nach 11/12 nur über `[10 11 3]` und `[10 11 12 3]` | je Mutant genau seine Zeile(n) von `TestVermittlungReihenfolge`: `[8]`, `[11]`, `[12]`, `[3 11]`+`[5 11]`, `[3 12]`, `[3 8]`, `[3 3]`+`[5 3]`, `[3 5]`, `[3 10]`, `[10 12]`, `[10 0]`, `[10 3]`, `[10 5]`, `[10 10]`, `[10 11 0]`, `[10 11 3]`, `[10 11 12 12]`, `[10 11 12 11]`, `[10 11 12 3]`, `[0 3]`, `[0 0]`, `[0 11]` |
| jede andere Nachricht des Clients als `p` und Terminate ist PGR-E6001 | W12: Anfrage gilt als Antwort; W13: Parse; W14: Sync; W15: FunctionCall; W16: nicht lesbare Nachricht beendet still | je Mutant genau seine Zeile von `TestAnmeldungAndereNachricht` (Anfrage, Extended-Nachricht, Sync, Funktion, nicht lesbar) |
| eine Nachricht `p` nach der Anmeldung ist wie zuvor nicht unterstützt | W17: `SetAuthType` nicht zurückgesetzt | `TestAnmeldungAndereNachricht` (Nachricht `p` nach dem Aufbau, mit Anmeldung) |
| Terminate in der Anmeldung ist kein Verbindungsfehler | W9 | `TestAnmeldungEndeDesClients/Terminate` |
| Ende des Datenstroms ist kein Verbindungsfehler | W10 | `…/Ende des Datenstroms nach der Anforderung`, `…/Reset` |
| Reset ist kein Verbindungsfehler | W11: `ECONNRESET` nicht als Ende der Verbindung | `…/Reset` |
| das Ende des Clients erzeugt weder Meldung noch Session | W20 (Handler behandelt `ErrClientEnde` wie einen Fehler); E2 (Integration) | W20: 4 Zeilen von `TestAnmeldungEndeDesClients`; E2: `TestE2ERecordAnmeldungClientBeendet` |

**DoD-Punkt 3 — der Port trägt den Austausch, der Kern kennt keine PGWire-Typen.**

| Zusage | Mutation | Befund |
|---|---|---|
| der Kern importiert keine PGWire-Bibliothek | Datei `internal/hexagon/model/mutant_import.go` mit Import von `pgproto3` | `make a-check`: `core-impurity: Kern importiert github.com/jackc/pgx/v5/pgproto3`, `gesamt: 1 Befund(e)`, Exit 2 |
| der Driven Adapter importiert keinen Driving Adapter | Datei im Upstream-Adapter mit Import des PGWire-Adapters | `make a-check`: `lateral-adapter`, `gesamt: 1 Befund(e)`, Exit 2 |
| ohne die Dateien | Stand `4c38fb7` | `gesamt: 0 Befund(e)` |

**Grüne Mutanten und äquivalente Fälle.**

- **Lesen am Server** (`LH-FA-05.b` *Ende des Servers*): `TestVermittlungLesefehlerAmServer` (Ende des Datenstroms, Reset); Mutant M5 (`istVerbindungsende` in `lesefehler` falsch) rot in beiden Zeilen.
- **Wächter zurücknehmen** (`SetReadDeadline(time.Time{})` in `startupBisSignal`, Mutant M14 des Reviews) bleibt grün. Akzeptiertes Negativ, seit `e1dc57b` als Grenze in `LH-FA-05.b` und in §6: Das Fenster (Signal nach der Startnachricht, vor dem Ende des Wächters) ist nicht stabil zu treffen; der Kommentar sagt nur zu, was `TestAnmeldungErstesSignal` und `TestAnmeldungWartetOhneFrist` prüfen.

- **Y5** (Fehler beim Senden der Antwort `p` übergangen) ist äquivalent: Auf der zurückgesetzten Verbindung
  scheitert das folgende Lesen mit demselben `PGR-E4003`; kein Fall unterscheidet die beiden Stellen.
- **W3** (`oeffne` kehrt beim Zwangsende nicht zurück) und **W3a** (`z.ende` schließt die Verbindung im Aufbau
  nicht) sind je für sich äquivalent: Beide Mechanismen beenden die Anmeldung beim Zwangsende, jeder reicht
  allein (`handle` schließt die Verbindung beim Rückkehren, `z.ende` schließt sie ohnehin). Zusammen (W3c) rot.

**Kommentare.** Jeder neue oder geänderte Kommentar ist gegen `AGENTS.md` §3.7 und §3.11 geprüft. Zusagen in
Kommentaren und ihre Tests: `model.Anmeldung`/`ErrClientEnde` („Ende des Datenstroms, Reset, Terminate, Zwangsende“)
→ `TestAnmeldungEndeDesClients`, `TestAnmeldungZwangsende`; `RecordService.OpenSession` („liest und hält nicht“)
→ `TestRecordAnmeldungNichtInDerAufzeichnung`; `startupBisSignal` → `TestAnmeldungErstesSignal`;
`Upstream.Open` („Kein Terminate“) → die Zeilen von U14; `rohAnforderung`, `antwort` („jeden Rumpf unverändert,
auch einen leeren“) → `TestAnmeldungVermittelt` (leere Antwort, 70000 Bytes). Die Aufruf-Voraussetzung
(„Aufrufe nacheinander“) ist eine Voraussetzung des Aufrufers (`LH-FA-05.b`), keine Zusage des Adapters; er
ruft den Rückruf in der Schleife von `aufbauen` nacheinander, und kein Mutant dazu ist gefahren.

### Closure

**Review und Verifikation.** Review `fe9cec1` (0 HIGH, 2 MEDIUM F-1 und F-2, 1 LOW F-3, 4 INFO F-4 bis F-7); Nacharbeit zu F-1 und F-2 in `755d786` (Lesefehler am Server, vier Zellen der Reihenfolge; M5 und M7 bis M10 rot in genau ihrer Zeile), zu F-4 und F-7 in `84ad1de` (Lauf auf dem letzten Stand, Zahl in §6). F-3 (M14) ist ein akzeptiertes Negativ, seit `e1dc57b` und `3fdf74e` in `LH-FA-05.b` *Abbruch und Herunterfahren* und in §6 entschieden. F-5 (Wortlaut der ADR) und F-6 (Text der Bibliothek, von `TestVermittlungNichtVermittelt` gefangen): keine Aktion. Verifikation `b27b137` (DoD-Punkte 1 bis 4 bestätigt gegen das gebaute Binary und PostgreSQL 17, `make gates` voll auf `adc3aee` grün, rund 60 Mutationen aus dem richtigen Grund rot, Handbuch und README unberührt): V-1 bis V-6.

**Entscheidungen zu den Übergaben (Planner, 2026-10-14).**

- **V-1** (LOW, Integration prüft das Log der Stufe `debug` nur auf `GEHEIM`): behoben in `01573f5`, die Spuren `md5[0-9a-f]{32}`, `SCRAM`, `n,,n=` und `biws` gelten auch für `rec.stderr`. Den Mutanten W18real (Debug-Zeile mit der Antwort, Integrationslauf) hat der Implementer nach der Nacharbeit gefahren (Zeile in der Tabelle oben); die Verifikation hat ihn danach nicht erneut gefahren.
- **V-2** (INFO, Ort des akzeptierten Negativs M14): erledigt, die Grenze und die Lesart „die Rücknahme der Lesefrist ist das Verhalten“ stehen in der Spezifikation (`e1dc57b`, `3fdf74e`), nicht nur in §7 und im Kommentar. M14 bleibt grün (akzeptiert).
- **V-3** (INFO, Mutation S3 im Service nicht gleichwertig nachzubauen): kein Befund am Slice. „Keine Session bei `ErrClientEnde`“ ist am Handler (W20, vier Zeilen) und in der Integration (`TestE2ERecordAnmeldungClientBeendet`) belegt, im Service über S3b.
- **V-4** (INFO, Rechnung der Mutanten): berichtigt (140 am Quelltext und 2 Architekturmutanten, 142, davon 139 rot und drei grün; M14 des Reviews steht gesondert). **V-5** (`CACHED`): kein Befund am Slice, die Belege kamen aus Läufen mit `--no-cache-filter`. **V-6**: mit dieser Closure erledigt.
- **Größe (§3.14).** Gemessen **1701** (1612 Einfügungen, 89 Löschungen; `git diff --numstat fa5a64c..3fdf74e -- internal test docs/user`; davon Produktcode 422, Tests 1249, Abdeckungstabellen 30) gegen die Marke 2500: **nicht erreicht**, die Rückführung `in-progress` → `next` ist nicht eingetreten. Die Meldeschwelle 2000 stand in §4 und ist nicht erreicht. Abgleich mit der Schätzung des Architect (1800 bis 2400): die Messung liegt 99 Zeilen (5 %) unter der Untergrenze. Seine Basis (Produktcode 400, Tests 800, Integration 290, zusammen 1500) lag 12 % unter der Messung; Produktcode traf (422), die Tests lagen mit dem Faktor 2,96 statt 2 über dem Ansatz, und der Aufschlag für Nacharbeiten deckte das ab.
- **Randformen im Code vor der Spezifikation.** Die vier Rückgaben des Implementers (Y1 bis Y4, Tabelle oben) kamen nach dem Code (`ce47692` bis `4c38fb7`, entschieden in `468997b`): Y1 und Y2 waren vor der Entscheidung grüne Mutanten, der Code hatte also eine Lesart (zweites Code 11 zugelassen). Der Architect entschied (`468997b`) vor dem Code-Commit der Nacharbeit (`e7e6953`: zweites Code 11, Code ohne Verfahren nach `AuthenticationOk`, Schreibfehler am Server), der den Code an die Entscheidung anpasste; kein Code-Commit änderte §6.

**Risiken aus §6 — Ausgang je Risiko** (geschlossene Menge):

| Risiko | Ausgang | Begründung |
|---|---|---|
| Größe und Schichten: drei Schichten, ungeteilt (Entscheidung des Nutzers vom 2026-10-14); keine Log-Zeile bei fehlgeschlagener Anmeldung, kein neuer Katalog-Code | entfallen | Die Rückführung über die Marke 2500 ist nicht eingetreten (gemessen 1701, Meldeschwelle 2000 nicht erreicht, Review in einer Sitzung); W21 (Log-Zeile im Pfad ohne Session) rot in `TestAnmeldungFehlerantwortDesServers`, `internal/hexagon/model/fehler.go` steht nicht im Diff (`git diff fa5a64c..3fdf74e --name-only`). |
| Anmeldenachrichten in der Aufzeichnung (Passwort-Hashes, `LH-RB-01`) | entfallen | S2 und E4 (je Antwort eigener Schlüssel, Verifikation: E4c) rot in `TestRecordAnmeldungNichtInDerAufzeichnung` und `TestE2ERecordAnmeldung` (alle drei Verfahren); die Verifikation las Aufzeichnung und Debug-Log eines eigenen Laufs; W18 und W19 (Debug-Zeile mit Antwort oder Anforderung) rot in den Tests des Adapters, W18real in der Integration (nach `01573f5`). |
| Der Lesepuffer des Aufbaus in der Session aus `Open` ist ungeprüft (`BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft`) | entfallen | `TestOpenSessionTraegtLesepufferDesAufbaus` ist rot unter U17 (`Open` liefert eine Session mit neuem Frontend; Implementer und Verifikation). Die Grenze (beide Nachrichten in einem Lesevorgang über TCP-Loopback, nicht zugesagt) steht in §6 und bleibt dort; Register: gestrichen, unten. |

Andere Risiken führt §6 nicht. Das akzeptierte Negativ M14 (Lesefrist im Rennfenster) und der Mutant Y5 sowie W3 und W3a (äquivalent) sind keine Risiken mit Ausgang und stehen mit ihrer Begründung in §6 und oben.

**Register** (Zähler = Dateien unter `evidence/`, gelesen am Stand dieser Closure; Beleg je `evidence/slice-v1-abschluss-anmeldung.md`):

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` — F-1 (Hälfte einer Zusage: Lesen am Server), F-2 (Zellen der Reihenfolge-Tabelle), V-1 (Zusage „keine Spur im Log“ in der Integration nur für das Passwort): 33 → 34. Verkörpert (`AGENTS.md` §3.10); Sensor geplant (`slice-harness-mutation`).
- `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` — Y1 bis Y4, zurückgegeben nach dem Code: 14 → 15 (Urteil des Planners, wie `slice-v1-abschluss-tls-client` F-9). Verkörpert (`.claude/commands/implement-slice.md` Randform-Rückgabe).
- `BEO-REPO/spec-randform-erst-im-review-entschieden` — F-3 und V-2: das Rennfenster der Rücknahme der Lesefrist war in §6 und Spezifikation nicht genannt, der Code und sein Kommentar sagten „gilt nicht“; entschieden erst nach dem Review: 28 → 29 (Urteil des Planners). Verkörpert (`AGENTS.md` §3.12).
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — F-3 (Kommentar von `startupBisSignal` sagte „gilt nicht“, M14 grün): 53 → 54. Verkörpert (`AGENTS.md` §3.11); Sensor geplant (`slice-harness-mutation`).
- `BEO-REPO/plan-folgt-korrektur-nicht` — F-4, F-7, V-4: der Beleg in §7 stand auf einem älteren Code-Stand, §6 nannte die überholte Zahl 1623, die Rechnung der Mutanten ging nicht auf: 43 → 44. Verkörpert (`AGENTS.md` §3.9).
- `BEO-REPO/session-traegt-puffer-des-aufbaus-ungeprueft` — Stand `offen` → `gestrichen`: Der Test existiert (Begründung im `state.md`).
- Kein Auftreten: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Schätzung 1800 bis 2400, gemessen 1701, unter der Untergrenze; die Regel von `AGENTS.md` §3.14 hat gegriffen, Zählerstand 4), `BEO-REPO/meldeschwelle-ohne-artefakt` (die Schwelle 2000 stand diesmal in §4 des Plans und ist nicht erreicht: Gegenbeleg, Zähler bleibt 1, kein Ausgang), `BEO-REPO/groessenprobe-ohne-messbaren-stand` (die Marke nannte ihre Metrik), `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (beide Nehmer nehmen an, unten), `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (Handbuch und README unberührt, der Nehmer führt den Teil; der Vermerk „Teil geliefert“ kommt mit der Closure von `slice-v1-abschluss-anmeldung-doku`), `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (nicht gezählt: die Referenzversion 17 ist in §1 benannt, die Matrix liefert `slice-v1-abschluss-postgres-versionen`).
- Kein nicht verkörperter Eintrag erreicht mit diesem Slice 3×; die verkörperten Regeln bleiben, `AGENTS.md` wird nicht geändert. Keine neue Beobachtung angefallen (F-6, die Abhängigkeit vom Fehlertext der Bibliothek, fängt ein Test).

**Lerneintrag — Wiederaufgetreten, Regel bleibt** (`AGENTS.md` §3.10, §3.12, §3.9); geändert wird keine Regel, kein `liegt in`. Beide MEDIUM-Funde des Reviews waren Zusagen, die §7 als vollständig führte und die ein Element ausließen, obwohl die Schärfung vom 2026-10-13 („je Element einer Aufzählung ein Fall und eine Mutation“) schon galt: F-2 (die Tabelle Stand × Code der Reihenfolge: 22 Zeilen von Hand gewählt, vier Zellen fehlten, §7 schrieb „jede Zelle“) und F-1 (zwei Hälften einer Zusage, Lesen und Schreiben am Server: nur das Schreiben hatte einen Fall). V-1 hat dieselbe Gestalt (Aufzeichnung und Log, nur eines geprüft). Die gemeinsame Ursache ist die Auswahl der Fälle von Hand aus dem Wortlaut; eine Tabelle des Codes (Stand, Code) erzeugt ihre Zellen selbst. Was funktionierte: Die Rückgaben des Implementers kamen an, der Architect entschied vor der Nacharbeit, jede Nacharbeit bekam eine Zeile Zusage · Mutation · roter Test, und die Verifikation fuhr rund 60 Mutationen selbst und kennzeichnete zwei eigene Fehlversuche (E4 mit festem Schlüssel, X3). Vorschläge an den Nutzer, nicht als Eintrag und nicht von dieser Closure geändert: (1) `.claude/commands/implement-slice.md` Schritt 19 (oder `AGENTS.md` §3.10): Eine Zusage über eine Zustandstabelle (Stand × Eingabe) bekommt einen Test, der alle Zellen der Tabelle des Codes durchläuft, statt einer von Hand gewählten Liste; die Mutation je Zelle ist dann die Mutation je Eintrag. (2) Die Meldeschwelle neben der Marke in §4 (Zahl und Metrik) hat in diesem Plan getragen; der Vorschlag aus `slice-v1-abschluss-tls-client` §7 (in `AGENTS.md` §3.14 oder `architect.md`) bleibt offen. (3) Offen aus `slice-v1-abschluss-tls-client-doku`: Handbuchsätze zu einem Gegenüber nennen dessen Version; hier betrifft es `slice-v1-abschluss-anmeldung-doku` (belegt ist PostgreSQL 17), der Planner hat es dort als Risiko in §6 eingetragen, ohne Regel. Retirement-Check bestanden (keine Regel entfernt oder gelockert).

**Paarungen.** (a) Anker: Der Lerneintrag trägt kein Feld `liegt in`, keine Verkörperung; nichts zu prüfen. (b) Folge-Slice: Genannt sind `slice-v1-abschluss-anmeldung-doku` (Handbuch und README), `slice-v1-abschluss-postgres-versionen` (Fall SCRAM-SHA-256 des Anmeldetests je Version) und `slice-tests-ueberlebende-mutanten` (Adresse für den Test zum Lesepuffer, hier geliefert). Die beiden ersten existieren in `next/` und nehmen an, `slice-tests-ueberlebende-mutanten` (`open/`) ist keine Adresse einer Sendung dieser Closure, sondern verweist auf diesen Slice: `grep -n "slice-v1-abschluss-anmeldung"` findet die Kennung des Gebers im §1 und in der DoD des Doku-Slice (Ziel, Herkunft, DoD-Punkt 1 und 2), im §1 *Anmeldung je Version*, in der DoD (Punkt 1) und im §3 von `slice-v1-abschluss-postgres-versionen`, in der Abgrenzung und der Tabelle von `slice-tests-ueberlebende-mutanten`. Der Fall SCRAM-SHA-256 ist einzeln wählbar (`TestE2ERecordAnmeldung/scram-sha-256`). Die Sendung dieser Closure (Handbuchsätze zu einem Gegenüber nennen dessen Version, belegt ist PostgreSQL 17) steht im selben Commit in §6 und §8 von `slice-v1-abschluss-anmeldung-doku`; Nachzählung: Liefer-Punkte 2, Schichten eine (Dokumentation), innerhalb der Grenze. Die Zeile zu R1 an `Open` in `slice-tests-ueberlebende-mutanten` zeigt auf das gelieferte Ergebnis. (c) Register: Alle genannten Kennungen `BEO-REPO/…` bestehen als Verzeichnis, jedes mit nicht leerem `evidence/`; `session-traegt-puffer-des-aufbaus-ungeprueft` trägt seine Datei aus `slice-harness-blackbox-driven`.

**Stand der Gates.** Die Läufe stehen oben; die Verifikation fuhr `make gates` voll auf `adc3aee`. Seitdem ändern `01573f5` einen Test, `e1dc57b` und `3fdf74e` die Spezifikation, die übrigen Commits nur Plan und Berichte. Den Lauf auf dem Stand der Closure nennt der Bericht des Planners, nicht dieser Absatz.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-anmeldung-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 3. Schichten: drei Schichten (Kern mit Ports, PGWire-Adapter, Upstream-Adapter) und damit über zwei, auch ohne die Dokumentation; das war vor der Abgabe so und hängt nicht an ihr: Der Architect hat geprüft (2026-10-14): Ein Schnitt senkt die Schichtzahl nicht (§6); Entscheidung des Nutzers vom 2026-10-14: ungeteilt mit drei Schichten, Marke 2500, Meldeschwelle 2000 (§4). Der Meldungstext zu einem nicht vermittelten Verfahren (aus `slice-doku-ist-stand`, V-152) liegt in `internal/adapters/driven/postgres/upstream.go` und in DoD-Punkt 2.

Gesichtet bei der Closure (Planner, 2026-10-14; die Aussage „es trägt nur seine `README.md`“ oben galt zum Datum des Plans): Das Register trägt inzwischen `BEO-REPO/`. Treffer und Zählung stehen in §7, *Register*; keiner der nicht verkörperten Einträge erreicht mit diesem Slice 3×. Nachgezählt bei der Closure (`AGENTS.md` §3.13): Liefer-Punkte 3, drei Schichten (Kern mit Ports, PGWire-Adapter, Upstream-Adapter; Ausnahme des Nutzers vom 2026-10-14), unverändert; der Bootstrap blieb im Produktcode unverändert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
