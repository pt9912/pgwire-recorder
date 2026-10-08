# Slice slice-harness-lint-warten-ohne-frist: Lint meldet Warten ohne Frist in Tests

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
nächsten Welle-Closure. Angelegt bei der Closure von `slice-v1-abschluss-herunterfahren`
nach Entscheidung des Nutzers vom 2026-10-08 (Register-Ausgang *geplant*); Reihenfolge
in §4 *Start*.

**Bezug:** kein Lastenheft-Bezug — die Prüfung gilt den Tests, nicht einer Zusage des
Produkts. Bindung an Entscheidungen:
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Lint-Gate mit eigenen
Prüfungen; ob die neue Prüfung eine eigene ADR braucht, entscheidet der Architect, §6).

**Berührte Spec-Stellen:** `SPEC-049` (neuer Punkt für die Prüfung, vom Architect vor
dem Code geschrieben) · `SPEC-038` (*Warten in Tests*, angewandt; die Ausnahme für
synchrone Aufrufe trägt der Architect nach Entscheidung des Nutzers vom 2026-10-08 ein) ·
`spezifikation.md` §12 (*Historie*)

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-08.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make lint` meldet in jeder Datei `*_test.go` einen Empfang aus einem Kanal und
ein Warten auf ein Prozessende, das ohne Frist wartet, als Befund
`lint: <pfad>:<zeile>: <befund>` (`SPEC-049` Punkt 9); ein `select` mit einem Fall für
einen Zeitgeber oder `time.After` gilt als Frist. Die Regel steht als neuer Punkt in
`SPEC-049`, vom Architect vor dem ersten Code-Commit geschrieben (`AGENTS.md` §3.12), und
setzt `SPEC-038` *Warten in Tests* in eine Maschine um. Danach trägt der Bestand keinen
Befund der Prüfung: `make lint` hängt an `make gates` (`SPEC-049` Punkt 10), und eine
Ausnahme, die Bestand ausblendet, gibt es nicht (Punkt 8).

**Herkunft — Geber `slice-v1-abschluss-herunterfahren`** (`AGENTS.md` §3.13).
`BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` erreicht mit dessen Closure 3×
(`slice-harness-blackbox-einstieg`: V-88, zwei `Wait` auf demselben Prozess;
`slice-lint-bestand-driving`: F-468, Warten ohne Frist auf `conn.gesetzt`;
`slice-v1-abschluss-herunterfahren`: F-485, fünf neue Tests ohne Frist, Hänger bis
`panic: test timed out after 20s`, und der erste Anlauf zu H4; Verifikation V-107). Die
Regel `SPEC-038` *Warten in Tests* stand schon, als F-485 entstand; gefunden hat ihn das
Review, kein Sensor. Entscheidung des Nutzers vom 2026-10-08: Ausgang *geplant* auf
diesen Slice, eine eigene Prüfung in `tools/harness/lint.sh` mit Gegenprobe-Fällen,
Sensor-Datei und Bereinigung des Bestands. Das Risiko *Test mit Signal hängt bis zum
Zeitlimit* des Gebers ist mit Ausgang *eingetreten* hierher gegeben.

**Gegenstand — der Bestand am Stand `bb36f2d`.** Gemessen vom Planner über den AST
(`go/parser` der Standardbibliothek, in `golang:1.27-alpine` ohne Netz, alle 37 Dateien
`*_test.go`), ohne Typinformation. Treffer der Regel im Wortlaut des Nutzers — ein
Empfang außerhalb eines `select` im Rumpf eines Tests oder einer Hilfsfunktion, nicht in
einer Goroutine und nicht in einer Methode — **13 Stellen in vier Dateien**:

- `internal/hexagon/services/record_extended_test.go` — 9: Zeilen 354, 402, 463, **529**
  (`<-up3.letzte.queryLaeuft`), 534, 714, **812** (`<-up.letzte.queryLaeuft`), 817, 842.
  Die Kern-Stellen 529 und 812 führt `slice-tests-ueberlebende-mutanten-driving` §6
  *Warten auf einen Kanal* als akzeptiertes Negativ ohne eigene Adresse; dieser Slice ist
  jetzt ihre Adresse (dort nachgezogen).
- `internal/adapters/driven/postgres/upstream_extended_test.go` — 2: Zeilen 123, 154
  (`<-empfangen`).
- `internal/adapters/driving/pgwire/server_extended_test.go` — 1: Zeile 848
  (`<-conn.gesetzt` in `TestReplayHerunterfahrenSpaeteFrist`, die Stelle aus F-468, am
  Stand `893b54e` Zeile 846).
- `internal/bootstrap/bootstrap_test.go` — 1: Zeile 353 (`exit := <-ende`).

Was die Randformen in §6 entscheiden, ist nicht mitgezählt: 17 Stellen in Methoden von
Test-Doubles (9 Empfänge, 8 `select` ohne Zeitgeber), 2 Empfänge in Funktionsliteralen,
die als Argument übergeben werden, 1 `sync.WaitGroup.Wait`, 1 `exec.Cmd.Wait` in einer
Goroutine, 8 `select` mit `default`; kein `select` nur mit `ctx.Done()`, keine Schleife
`range` über einen Kanal. Die Zahl beim Start misst der Implementer neu (DoD-Punkt 3);
die Wellen- und Sammel-Slices davor schreiben weitere Tests.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Synchrone Aufrufe des Prüflings ohne Frist, etwa in der Hilfsfunktion `session` in
  `internal/hexagon/services/record_test.go` — Bestand bleibt bewusst stehen: Nach
  Entscheidung des Nutzers vom 2026-10-08 fallen sie nicht unter `SPEC-038` *Warten in
  Tests*, die Regel gilt nur für Kanäle und Prozessenden; der Architect trägt das in
  `SPEC-038` ein. Ein Hänger wie im ersten Anlauf zu H4 bleibt damit außerhalb der
  Prüfung (Risiko in §6).
- Die Frist in `TestReplayHerunterfahrenSpaeteFrist` (Zeile 848, Review F-468) und die
  übrigen Stellen im PGWire-Adapter und im Bootstrap — liefert
  `slice-tests-ueberlebende-mutanten-driving` (DoD-Punkt 1, §6 *Warten auf einen
  Kanal*); er schließt nach §4 vor diesem Slice. Dieser Slice setzt dort keine zweite
  Frist; was die Messung beim Start in diesen Schichten noch trifft, gehört nach
  DoD-Punkt 3 zu seinem Bestand.
- Wie ein Mutations-Gate einen Mutanten zählt, der nur über die Zeitüberschreitung rot
  wird — anderer Vorgang: Randform *Laufzeit* von `slice-harness-mutation`.
- Warten ohne Frist im Produkt-Code — Schicht-Abgrenzung: `SPEC-038` regelt Tests; der
  Slice ändert Harness-Werkzeug und Testdateien, keinen Produkt-Code.
- Ein eigenes `make`-Ziel für die Prüfung — anderer Vorgang: Der Nutzer hat sie als
  eigene Prüfung in `tools/harness/lint.sh` entschieden; sie läuft mit `make lint` und
  dessen Gegenprobe an der Gate-Kette.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Die Sendung kommt von `slice-v1-abschluss-herunterfahren` (Risiko *Test mit Signal hängt
bis zum Zeitlimit*, Register-Ausgang *geplant* für
`BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`). Belege des Implementers stehen in §7,
nicht im Bericht.

- [ ] **Prüfung:** `tools/harness/lint.sh` meldet nach dem neuen Punkt von `SPEC-049`
      jeden Empfang aus einem Kanal und jedes Warten auf ein Prozessende in `*_test.go`,
      das ohne Frist wartet, mit einer Zeile `lint: <pfad>:<zeile>: <befund>`; ein
      `select` mit einem Fall für einen Zeitgeber oder `time.After` ist kein Befund; jede
      Randform aus §6 gilt so, wie der Architect sie entschieden hat. Die Prüfung läuft,
      auch wenn eine andere einen Befund hat (Punkt 9). `harness/sensors/lint.md` nennt
      sie in *Vertrag* und *Grenze* (auch: synchrone Aufrufe prüft sie nicht), der Kopf
      von `lint.sh` ihre Fälle der Gegenprobe. Beleg in §7: je Zusage Zusage · Mutation ·
      roter Fall (`AGENTS.md` §3.10).
- [ ] **Gegenprobe:** `tools/harness/lint-gegenprobe.sh` führt je Randform aus §6
      mindestens einen Fall in einer Kopie des Arbeitsbaums — rot mit genau der eigenen
      Zeile, wo der Architect einen Befund entschieden hat, ohne Befund, wo er die Form
      frei entschieden hat (mindestens: Empfang ohne `select`, `select` mit `time.After`,
      `select` mit Zeitgeber, `select` mit `default`, Test-Double, das auf die Freigabe
      des Tests wartet, `Wait` auf einen Prozess, synchroner Aufruf ohne Frist). `make
      lint-gegenprobe` ist grün; die Mutation, die die Prüfung abschaltet, macht die
      Gegenprobe rot (§7).
- [ ] **Bestand:** Vor dem ersten Code-Commit misst der Implementer die Treffer der
      Prüfung am Stand des Starts und trägt die Liste in §7 ein (Vergleich mit der Messung
      in §1). Danach meldet `make lint` keinen Befund der Prüfung. Jede bereinigte Stelle
      wartet nach `SPEC-038` *Warten in Tests*: Frist als Literal, höchstens 60 s, die
      Meldung nennt das ausgebliebene Ereignis. Die Kern-Stellen
      `record_extended_test.go:529` und `:812` werden unter einem Mutanten, der die
      einfache Anfrage vor dem Upstream anhält, binnen ihrer Frist rot statt an der
      Zeitgrenze von `go test` (Beleg in §7). Kein Produkt-Code.
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
| `spec/spezifikation.md` | update (Architect, vor dem Code) | neuer Punkt in `SPEC-049` mit den Randformen aus §6; in `SPEC-038` *Warten in Tests* die Ausnahme für synchrone Aufrufe (Entscheidung des Nutzers vom 2026-10-08), falls noch nicht eingetragen; §12 *Historie* |
| `docs/plan/adr/` | neu (nur wenn der Architect eine ADR verlangt) | [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) ist angenommen und unveränderlich (`AGENTS.md` §3.5); eine neue Prüfung, die nicht aus ihr folgt, wäre eine ergänzende ADR mit `Schärft:` auf `SPEC-049` |
| `tools/harness/lint.sh` | update | die Prüfung; Kopf *GEPRÜFT DURCH* um den neuen Punkt |
| Programm unter `tools/harness/` | neu (nur bei Erkennung über den AST, §6) | liest die Testdateien mit `go/parser` in der Stufe `lint`, ohne Netz |
| `tools/harness/lint-gegenprobe.sh` | update | Fälle je Randform, rot und frei (DoD-Punkt 2) |
| `harness/sensors/lint.md` | update | *Vertrag* und *Grenze* der Prüfung |
| `harness/README.md` | update | Zeile `make lint` der Sensors-Tabelle nennt die Prüfung |
| `internal/hexagon/services/record_extended_test.go` | update | neun Stellen nach §1, Frist nach `SPEC-038` (Negative: Mutant an `:529`, `:812`) |
| `internal/adapters/driven/postgres/upstream_extended_test.go` | update | zwei Stellen nach §1 |
| weitere `*_test.go` | update | was die Messung beim Start sonst trifft (DoD-Punkt 3) |

Die Abdeckungs-Deklarationen bleiben, wie sie sind: Der Slice ändert, wie Tests warten,
nicht, was sie abdecken; `make abdeckung-check` bestätigt es.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-tests-ueberlebende-mutanten-driving` liegt in
`done/` (WIP-Limit 1). Ist er ausdrücklich zurückgestellt, trifft die Prüfung auch dessen
Stellen im PGWire-Adapter und im Bootstrap; dann schneidet der Planner vor dem Start
neu (§1 *Ausdrücklich NICHT*). Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08
(Wellen vor Harness, M3 vor M4): die Slices von welle-v1-abschluss, darin
`slice-harness-integration-wait` direkt vor `slice-v1-abschluss-herunterfahren` und
`slice-harness-meldungskatalog-gate` direkt vor `slice-v1-abschluss-container`, dann
`slice-harness-abdeckung-gate` und `slice-harness-coverage`, dann die Slices von
welle-erster-release, danach `slice-harness-commit-struktur-id`,
`slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`,
`slice-harness-lint-warten-ohne-frist`, `slice-harness-mutation`; die Reihenfolge der
Wellen-Slices steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md).

*Warum nach den beiden Sammel-Slices.* `make lint` ist ein Gate, und `SPEC-049` Punkt 8
lässt keine Ausnahme für Bestand zu: Die Prüfung kommt erst, wenn der Bestand ohne Befund
ist. `slice-tests-ueberlebende-mutanten-driving` setzt die Frist in
`TestReplayHerunterfahrenSpaeteFrist` schon (DoD-Punkt 1) und wendet `SPEC-038` im
PGWire-Adapter und im Bootstrap an; davor müsste dieser Slice ihm diesen Teil abnehmen und
seine DoD neu schneiden. Danach misst er nur, was übrig bleibt, im Kern und in den
Driven-Adaptern. Der Preis: Die neuen Tests der beiden Sammel-Slices (G2, G7 und die
Frist aus F-468 beobachten über Fristen) entstehen ohne den Sensor; ihr Bestand fällt
beim Start unter DoD-Punkt 3. *Warum vor `slice-harness-mutation`.* Dessen Randform
*Laufzeit* entscheidet, wie ein Mutant zählt, der nur über die Zeitüberschreitung rot
wird; mit dieser Prüfung wartet kein Test mehr ohne Frist auf einen Kanal oder ein
Prozessende, und ein solcher Mutant wird binnen der Frist rot.

Erster Schritt nach dem Start, vor jedem Code-Commit: Der Architect entscheidet die
Randformen aus §6 und schreibt den neuen Punkt in `SPEC-049`
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`); der Implementer misst den
Bestand neu.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Messung beim Start trifft
  so viele Stellen oder so viele Testpakete, dass Prüfung, Gegenprobe und Bereinigung
  nicht in eine Review-Sitzung passen. Dann wird die Bereinigung ein eigener Slice vor
  diesem — nicht danach, denn ohne bereinigten Bestand ist das Gate rot.
- `in-progress` → `open` (blockiert — Carveout?): Die Erkennung, die der Architect wählt,
  braucht in der Stufe `lint` etwas, das sie ohne Netz nicht hat (etwa Typinformation
  über Module außerhalb von `deps`), oder eine Randform lässt sich ohne Markierung im
  Testcode nicht unterscheiden (Test-Double gegen Warten auf den Prüfling) und braucht
  eine Entscheidung, die `SPEC-038` ändert.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make lint` ist mit der Prüfung grün, die Gegenprobe macht jeden
entschiedenen Fall sichtbar, und die Verifikation hat Prüfung, Gegenprobe und den
Mutanten an `:529`/`:812` nachgefahren. Die Closure-Notiz mit Lerneintrag ist
geschrieben, und `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` wechselt von *geplant*
auf *verkörpert* mit Zielort und Herkunfts-Anker `seit slice-harness-lint-warten-ohne-frist`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12). Die Prüfung ist ein neuer Vertrag eines
Harness-Werkzeugs; entschieden wird in der Spezifikation (`SPEC-049`, neuer Punkt), vor
dem ersten Code-Commit, vom Architect. Gesetzt vom Nutzer am 2026-10-08 sind nur die
ersten drei; alle übrigen sind **offen**. Die Zahlen sind die Messung des Planners am
Stand `bb36f2d` (§1).

- **Empfang aus einem Kanal außerhalb eines `select`** (`<-ch`, `v := <-ch`,
  `v, ok := <-ch`) im Rumpf eines Tests oder einer Hilfsfunktion — **Befund**
  (Entscheidung des Nutzers). 13 Stellen.
- **`select` mit einem Fall für einen Zeitgeber oder `time.After`** — **kein Befund**
  (Entscheidung des Nutzers). Offen ist, woran die Prüfung den Fall erkennt:
  `<-time.After(…)`, `<-time.Tick(…)`, `<-t.C` eines `time.NewTimer` oder
  `time.NewTicker` (Name beliebig), ein Kanal, den ein Helfer mit Zeitgeber liefert.
- **Synchroner Aufruf des Prüflings ohne Frist** (Hilfsfunktion `session` in
  `internal/hexagon/services/record_test.go`) — **kein Befund**; nach Entscheidung des
  Nutzers fällt er nicht unter `SPEC-038` *Warten in Tests*, der Architect trägt das dort
  ein.
- **`select` mit `ctx.Done()` ohne Zeitgeber** — offen. Eine Frist ist es nur, wenn der
  Kontext eine trägt (`context.WithTimeout`, `WithDeadline`, `t.Context()` unter
  `-timeout` nicht); im Testkörper heute 0 Stellen, `<-ctx.Done()` allein 1 in einer
  Methode eines Doubles (`server_test.go:94`).
- **`select` mit `default`** — offen; wartet nicht, ist also kein Warten ohne Frist. 8
  Stellen (davon 6 im Testkörper).
- **`select` ohne Zeitgeber und ohne `default`** — offen, ob er wie ein Empfang zählt. 8
  Stellen, alle in Methoden von Doubles.
- **Schleife `range` über einen Kanal** — offen; sie wartet bis zum Schließen. Ohne
  Typinformation ist sie von einer Schleife über eine Slice nicht zu unterscheiden. 0
  Stellen über einen Kanal (alle Schleifen im Bestand laufen über Slices oder Maps).
- **`sync.WaitGroup.Wait`** — offen; die Entscheidung des Nutzers nennt Kanäle und
  Prozessenden, nicht Wartegruppen. 1 Stelle (`record_test.go:301`, wartet auf
  Goroutinen, die den Prüfling synchron rufen — nah am synchronen Aufruf).
- **Warten auf ein Prozessende** (`exec.Cmd.Wait`, `os.Process.Wait`) — **Befund** ohne
  Frist (Entscheidung des Nutzers); offen, wie es zu `SPEC-038` passt, wonach je Prozess
  genau eine Stelle wartet und jeder weitere Leser ihr Ergebnis liest: Die eine Stelle
  `r.waitErr = cmd.Wait()` liegt in einer Goroutine (`record_e2e_test.go:381`), die Leser
  warten mit Frist.
- **Empfang oder `Wait` in einer Goroutine des Tests** (`go func() { … <-ch … }()`) —
  offen; die Goroutine hält den Test nicht an, wer ihr Ergebnis liest, schon. 0 Empfänge,
  1 `Wait` (oben).
- **Test-Doubles, die auf eine Freigabe durch den Test warten** — nach `SPEC-038` kein
  Befund; offen, woran die Prüfung sie erkennt: Methode mit Empfänger in einer
  `*_test.go`, Funktionsliteral, das der Test einem Double oder Testserver übergibt
  (`upstream_extended_test.go:352`, `<-stumm`), Name des Kanals, Markierung. 17 Stellen
  in Methoden, 1 im Funktionsliteral. Eine Methode, die auf den Prüfling wartet statt auf
  den Test, wäre unter jeder Regel nach Ort frei.
- **Hilfsfunktionen** — offen: (a) ein Helfer mit Frist, dem der Test ein
  Funktionsliteral mit einem Empfang übergibt (`fertigBinnen(t, 10*time.Second, "Send",
  func() { … <-sendeFehler … })`, `upstream_extended_test.go:340`) — der Empfang ist
  durch den Helfer begrenzt, steht aber syntaktisch frei; (b) ein Helfer, der selbst mit
  Zeitgeber wartet (`ereignisBinnen`, `zurueckBinnen`, `endetBinnen`, `warteEnde`) — frei
  nach der Regel zum `select`; (c) ein Helfer ohne Frist, den Tests rufen — Befund an der
  Stelle im Helfer, nicht an jedem Aufruf (heute 0).
- **Erkennung über den AST statt über Regex** — offen. Die Prüfungen nach Punkt 6 und 7
  sind Regex in `bash`; ein Regex sieht `<-` auch in Strings und Kommentaren und kann den
  Fall eines `select` nicht seinem Block zuordnen. Ein Programm mit `go/parser` der
  Standardbibliothek läuft in der Stufe `lint` (Go im Image, ohne Netz) und unterscheidet
  `select`-Fälle, Goroutinen, Methoden und Funktionsliterale; `range` über einen Kanal
  braucht darüber hinaus Typinformation (`go/types` mit den Modulen aus `deps`). Das
  Programm fällt selbst unter `make lint` und braucht seine Gegenprobe. Die Messung in §1
  lief über `go/parser`, ohne Typen.
- **Ort und Text des Befunds** — offen: Zeile des Empfangs bzw. des `select` bzw. des
  `Wait`; Text in der Form von Punkt 9.
- **Gegenstand** — Dateien `*_test.go` unter `cmd/`, `internal/` und `test/`, auch mit
  dem Build-Tag `integration` (Punkt 1); die Brücke `export_test.go` gehört dazu.
- **Ausnahme** — keine (Punkt 8: keine Regel, die Bestand ausblendet; Punkt 6: kein
  Kommentar-Marker).
- **ADR** — offen, ob die Prüfung aus
  [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) folgt oder eine
  ergänzende ADR braucht (`AGENTS.md` §3.5, §3.8).

**Risiken:**

- **Die Prüfung meldet zu viel** — eine Erkennung nach Ort (Methode, Funktionsliteral)
  trennt Doubles nicht sicher vom Warten auf den Prüfling; die Bereinigung ändert dann
  Doubles, die `SPEC-038` ausnimmt, oder lässt ein echtes Warten frei. Gegenprobe-Fälle
  je Seite (DoD-Punkt 2). — **Ausgang:** — (bei Closure)
- **Ein Hänger, den die Prüfung nicht sieht** — ein synchroner Aufruf des Prüflings (wie
  im ersten Anlauf zu H4 in `slice-v1-abschluss-herunterfahren`) oder eine Randform, die
  der Architect frei entscheidet. Tritt einer auf, zählt ihn
  `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` weiter. — **Ausgang:** — (bei Closure)
- **Der Bestand wächst bis zum Start** — die Slices der Wellen und die beiden
  Sammel-Slices schreiben Tests mit Kanälen; die Messung beim Start (DoD-Punkt 3)
  entscheidet über die Rückführung aus §4. — **Ausgang:** — (bei Closure)
- **Laufzeit der Gegenprobe** — jeder neue Fall läuft die Stufe `lint` in einer Kopie
  (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, 1×); Laufzeit vorher und nachher in
  §7. — **Ausgang:** — (bei Closure)
- **Mutant kommt im Build-Kontext nicht an** — eine Mutation mit unveränderter Größe und
  mtime überträgt BuildKit nicht (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`,
  1×); die Mutanten laufen in frischen Kopien. — **Ausgang:** — (bei Closure)

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
`docs/plan/planning/observations/BEO-REPO/` am Stand der Closure von
`slice-v1-abschluss-herunterfahren` gesichtet (Zähler = Dateien unter `evidence/`).
Treffer:

- `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` (3×, *geplant* auf diesen Slice) — der
  Gegenstand; mit der Closure wird er *verkörpert* (§5).
- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (1×) — die Gegenprobe wächst um die
  Fälle aus DoD-Punkt 2; Risiko in §6.
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — die Prüfung tritt neben die
  Prüfungen nach Punkt 6 bis 8, nicht an ihre Stelle; die Gegenprobe-Fälle der übrigen
  Punkte bleiben grün.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×) — liest sich eine Regel
  als Zusage, die das Werkzeug anders auswertet, zeigt es erst die Gegenprobe mit freiem
  und rotem Fall; DoD-Punkt 2 verlangt beide.
- `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — Risiko in §6.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — wellenlos; §4 *Start* nennt
  den Architect-Schritt vor dem Code.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (14×, verkörpert in `AGENTS.md`
  §3.12) und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (4×, verkörpert)
  — darum stehen die offenen Randformen vor dem Code in §6.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (15×, `AGENTS.md` §3.10) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (21×, §3.11) — je Zusage eine
  Mutation (DoD-Punkt 1); die Sensor-Datei sagt nur zu, was die Gegenprobe prüft.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — §1 nennt den Geber; der
  Hinweis in `slice-tests-ueberlebende-mutanten-driving` §6 zeigt hierher.

Keiner der Einträge unter der Schwelle erreicht mit diesem Plan 3×; vor dem Code entsteht
keine neue Lücke.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
