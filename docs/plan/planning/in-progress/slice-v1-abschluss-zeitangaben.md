# Slice slice-v1-abschluss-zeitangaben: Zeitgetreues Einspielen (`--keep-timing`)

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

**Berührte Spec-Stellen:** `LH-FA-21.a` · `LH-FA-17.a` · `LH-FA-20` · `SPEC-034` · `ARC-002` · `ARC-005`

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

**Ziel:** `play --keep-timing` stellt den zeitlichen Abstand der aufgezeichneten Interaktionen her, ohne ihn zu verkürzen: `--timing-mode relative` (Default) oder `absolute`, bei `absolute` mit `--timing-reference connect` (Default) oder `first-request`; dazu Umgebungsvariable und Schlüssel im Abschnitt `play:`. Ungültige Verwendung ist `PGR-E2001`, eine Aufzeichnung, der bei `--keep-timing` irgendeiner Interaktion `offset_ms` fehlt, `PGR-E2003` (Exit-Code 2, Startfehler). Der Slice schließt Abnahmeszenario 13: aufgezeichnet mit `--record-timing`, eingespielt mit `--keep-timing`.

**Übernommen aus `slice-v1-abschluss-einspielen`** (dort §1, Abgrenzung): Zeitangaben und zeitgetreues Einspielen; hier die Hälfte beim Einspielen (das Ziel oben), die Hälfte beim Aufzeichnen liegt in `slice-v1-abschluss-zeitangaben-aufzeichnen`. Die Optionen `--keep-timing`, `--timing-mode` und `--timing-reference` sind bis zu diesem Slice unbekannt (`PGR-E2001`), ihre Schlüssel im Abschnitt `play:` unbekannt (`PGR-E2004`) (dort §6, *Optionen der Folge-Slices*; `slice-v1-abschluss-upstream-verbinden` ließ die Schlüssel `keep_timing` und `timing_mode` bei diesem Slice).

**Übernommen aus `slice-v1-abschluss-sqlite-format-lesen`** (dort §1, Abgrenzung): die Zeitangaben in der SQLite-Datei beim Einspielen. Das Lesen der Spalte `interaction.offset_ms` liefert jener Slice (nur Spalte mit geprüftem Typ); die Prüfung absteigender Werte liegt in `slice-v1-abschluss-zeitangaben-pruefen`, das Schreiben mit `--record-timing` in `slice-v1-abschluss-zeitangaben-aufzeichnen`. Dieser Slice spielt Aufzeichnungen mit Zeitangaben aus beiden Formaten ein.

**Abgegeben beim Schnitt vom 2026-10-11** (je mit der Kennung dieses Slice im Nehmer, `AGENTS.md` §3.13): das Aufzeichnen mit `--record-timing` an `slice-v1-abschluss-zeitangaben-aufzeichnen`, die Prüfung absteigender Werte beim Lesen (der Halbsatz „ein ungültiger Wert ist `PGR-E3003`“ aus DoD-Punkt 1) an `slice-v1-abschluss-zeitangaben-pruefen`, der Uhr-Port mit Fake-Uhr und Systemuhr an `slice-v1-abschluss-zeitangaben-uhr`. Dieser Slice benutzt den Port, wie `-uhr` ihn liefert, und die geprüften Leser. Der Port ist `Uhr` in `internal/hexagon/ports/driven/uhr.go` ([ADR-0015](../../adr/0015-uhr-port.md)) mit `Jetzt() int64` (ganze Millisekunden seit einem Ursprung, den die Uhr wählt, nie kleiner als die vorige Lesung) und `WarteBis(ctx, zeitpunkt int64) error`: kehrt beim Erreichen des Zeitpunkts mit `nil` zurück, bei einem erreichten oder vergangenen sofort, und mit dem Fehler des Kontexts, wenn er beendet ist — vor und während des Wartens und auch bei einem vergangenen Zeitpunkt; `math.MaxInt64` wartet bis zum Ende des Kontexts, `math.MinInt64` kehrt sofort zurück. In Tests die Fake-Uhr `testuhr.Neu` mit `Vorruecken` (rückt nur auf Anweisung vor, eine negative Zahl ist eine Panik) und `Wartende`; im Betrieb die Systemuhr `NeueSystemuhr` aus `internal/bootstrap`, die `slice-v1-abschluss-zeitangaben-aufzeichnen` zuerst verdrahtet und die dieser Slice für `play` an den Play-Service reicht.

**Herkunft:** Teil des Schnitts von `slice-v1-abschluss-zeitangaben` vom 2026-10-11 (Entscheidung des Nutzers vom 2026-10-11: vor dem Anspruch zu schneiden, weil der Plan vier Schichten berührte — Kern mit Uhr-Port, PGWire-Adapter, Recording-Adapter, CLI-Adapter — und damit über der Grenze von zwei lag, `AGENTS.md` §3.13; den Zuschnitt in vier Teile nach Lieferwert entschied der Planner mit seiner Empfehlung). Die vier Teile: `slice-v1-abschluss-zeitangaben-pruefen` (Lesen prüft die Werte), `slice-v1-abschluss-zeitangaben-uhr` (Uhr-Port und Systemuhr), `slice-v1-abschluss-zeitangaben-aufzeichnen` (`--record-timing`) und `slice-v1-abschluss-zeitangaben` (`--keep-timing`, schließt Abnahmeszenario 13); Handbuch und README zu allen vier liegen gemeinsam im Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` hinter dem letzten Teil (wie `slice-v1-abschluss-sqlite-format-doku` hinter dem Schreiben und dem Lesen: Zwischen den Code-Slices und dem Doku-Slice hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit allen).


**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `--record-timing` und das Schreiben von `offset_ms` — `slice-v1-abschluss-zeitangaben-aufzeichnen` (dort §1 und DoD): Ein anderer Weg, der Schreiben, Messpunkt und eine Option des Befehls `record` trägt; die Hälfte beim Einspielen braucht ihn nicht, weil Tests Aufzeichnungen mit Zeitangaben von Hand bauen.
- Die Prüfung absteigender Werte beim Lesen — `slice-v1-abschluss-zeitangaben-pruefen` (dort §1 und DoD): Sie gilt auch ohne diesen Slice (`replay`, `play` ohne Option) und liegt in einer anderen Schicht, dem Recording-Adapter.
- Der Uhr-Port, die Fake-Uhr und die Systemuhr — `slice-v1-abschluss-zeitangaben-uhr` (dort §1 und DoD): Dieser Slice ändert die Signaturen `Jetzt` und `WarteBis` nicht; verlangt das Warten eine andere, ist das ein Plan-Wechsel jenes Slice.
- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-zeitangaben-doku` (dort §1 und DoD), hinter dem letzten der vier Teile in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über der Grenze der Schichten.
- Zeitgetreue Antworten im Replay-Modus — Out-of-Scope von LH-FA-21.
- Paralleles Einspielen — Out-of-Scope von LH-FA-20 und LH-FA-21.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-21`](../../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen): `--keep-timing` hält die Pausen in beiden Modi (`relative`, `absolute`) und beiden Bezugspunkten (`connect`, `first-request`) ein, ohne sie zu verkürzen (Abnahmeszenario 13: Test mit Fake-Uhr `testuhr.Neu` im Service, in dem `Vorruecken` die Wartenden (`Wartende`) von `WarteBis` weckt, dazu ein E2E am Binary, aufgezeichnet mit `--record-timing`, eingespielt mit `--keep-timing`, geprüft auf die Untergrenze der Pausen); ohne `--keep-timing` werden vorhandene Zeitangaben ignoriert, und `replay` ignoriert sie immer (Test, Boundary).
- [ ] Optionen `--keep-timing`, `--timing-mode` und `--timing-reference` für `play` mit Umgebungsvariablen und Schlüsseln im Abschnitt `play:` nach der Optionstabelle (`LH-FA-17.a`; Default `false`, `relative`, `connect`; Priorität Option vor Variable vor Datei); bei `record` und `replay` unbekannt; ungültige Verwendung — `--timing-mode` oder `--timing-reference` ohne `--keep-timing`, `--timing-reference` zusammen mit `relative`, ein Wert außerhalb der Aufzählungen — ist `PGR-E2001`, Exit-Code 2, Startfehler; `config show` und `--help` nennen sie (Test am Binary, Negativ je Fall).
- [ ] Eine Aufzeichnung, der bei `--keep-timing` irgendeiner Interaktion `offset_ms` fehlt, wird als ungültige Verwendung abgelehnt: `PGR-E2003`, Exit-Code 2, Startfehler, vor der ersten Anfrage (Test, Negativ).
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
| `internal/hexagon/services` (Play-Service: `play.go`, neu `play_zeit.go`) | update, neu | `NewPlayService` nimmt die Option `Wartezeiten(uhr)` und in `PlayOptions` die Felder `KeepTiming`, `TimingMode`, `TimingReference`; vor jeder Interaktion wartet `warte` über die Uhr des Ports: `relative` und `absolute`, Bezugspunkt `connect` (Lesung am Beginn des Verbindungsaufbaus, vor TLS und Anmeldung) oder `first-request`; Prüfung auf `offset_ms` aller Interaktionen im Start, nach dem Laden (`PGR-E2003`); das erste Abbruchsignal beendet die Wartezeit, das zweite ebenso, und `Play` kehrt erst nach dem Schließen der Verbindung zurück |
| `internal/hexagon/model` (`fehler.go`) | update | der Meldungscode `PGR-E2003` (`CodeTimingMissing`), Klasse Konfiguration, Exit-Code 2 |
| `internal/hexagon/ports/driving` (Player) | keine Änderung | der Port bleibt `Play(ctx, ablauf)`; die Optionen reisen in `PlayOptions`, die Uhr in der Option `Wartezeiten` des Konstruktors |
| `internal/adapters/driving/cli` (`cli.go`, neu `zeitangaben.go`) | update, neu | Optionen `--keep-timing`, `--timing-mode`, `--timing-reference`, Umgebungsvariablen, Schlüssel im Abschnitt `play:`, Wertemengen in Kleinbuchstaben, Verwendungsfehler `PGR-E2001` aus der Zusammenführung (`zeitPlay`: „gesetzt“ ist, was irgendeine Quelle nennt), Hilfe (`config show` zeigt Schlüssel und Variablen ohne eigenen Code) |
| `internal/bootstrap` | update | `services.Wartezeiten(NeueSystemuhr())` an `NewPlayService` (reine Verdrahtung, keine Schicht, `AGENTS.md` §3.13; die Logik der Systemuhr liegt in `slice-v1-abschluss-zeitangaben-uhr`); `ctx` (erstes Signal) und `ablauf` (zweites) gibt der Bootstrap dem Service schon heute |
| Tests, `test/integration` | neu/update | Service mit Fake-Uhr nach `LH-FA-21` (`play_zeit_test.go`: Happy je Modus und Bezugspunkt, Boundary, Negativ), Optionen (`play_zeit_test.go` im CLI-Adapter), Verdrahtung (`internal/bootstrap/play_zeit_test.go`), Binary (`zeitangaben_einspielen_e2e_test.go`: aufzeichnen → einspielen, Quellen, Verwendung, `PGR-E2003`, Hilfe, `replay`); der Platzhalter `keep_timing` in `TestDateiUngueltig` (Schlüssel unbekannt) entfällt, `keep_timing` in `record:`, `timing_mode` in `replay:` und `timing_reference` auf der obersten Ebene sind dort Fälle, die für immer ungültig bleiben, dazu `shutdown_timeout` in `play:`; `docs/user/abdeckung-*.md` von `make abdeckung` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-zeitangaben-pruefen`, `slice-v1-abschluss-zeitangaben-uhr` und `slice-v1-abschluss-zeitangaben-aufzeichnen` liegen in `done/`, dazu `slice-v1-abschluss-einspielen`, `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` (Schritt 22 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter `-aufzeichnen`; WIP-Limit 1). Vor dem ersten Code-Commit entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar, oder die Wartelogik braucht eine dritte Schicht (PGWire-Adapter, Recording-Adapter, Bootstrap mit Logik); Schnitt dann: `relative` mit `connect` hier, `absolute` mit seinen Bezugspunkten in einem eigenen Slice.
- `in-progress` → `open` (blockiert — Carveout?): Die Fake-Uhr bildet die monotone Uhr nicht ab — Carveout.

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

**Randformen** (`AGENTS.md` §3.12) — vom Architect vor dem Code am 2026-10-11 entschieden, Ort: [`LH-FA-21.a`](../../../../spec/spezifikation.md#lh-fa-21a--zeitangaben) *Einspielen* und *Randformen des Einspielens*. Entscheidungen des Architect (Empfehlung, bis der Nutzer anderes sagt; der Nutzer hat am 2026-10-11 verfügt, ohne Nachfrage weiterzumachen).

- *Entschieden in `LH-FA-21.a` *Einspielen*:* `relative` wartet vor jeder Interaktion außer der ersten einer Session, bis seit dem Beginn der vorigen mindestens der Unterschied ihrer beiden `offset_ms` vergangen ist; eine Verspätung wird nicht aufgeholt. `absolute` wartet bis zum Bezugspunkt plus aufgezeichnetem Abstand; eine Verspätung wird aufgeholt, der Abstand zweier Interaktionen kann kürzer sein. `connect` ist der Beginn des Verbindungsaufbaus zum Server, vor TLS und Anmeldung (aufgezeichnet gilt `offset_ms` = 0, Aufbau und Anmeldung zählen als Verspätung); `first-request` ist der Beginn der ersten Interaktion der Session (aufgezeichnet gilt deren `offset_ms`). Jede Session beginnt mit eigener Uhr, die Sessions laufen nacheinander.
- *Entschieden ebenda, ungültige Verwendung:* `--timing-mode` oder `--timing-reference` ohne `--keep-timing`, `--timing-reference` zusammen mit `relative` und ein Wert außerhalb der Aufzählungen sind `PGR-E2001`; fehlt irgendeiner Interaktion der Aufzeichnung `offset_ms`, ist die Verwendung `PGR-E2003` (Exit-Code 2, Startfehler). Ohne `--keep-timing` werden Zeitangaben ignoriert, der Replay-Modus ignoriert sie immer.
- *Zusagen des Ports am Warten:* Der Fehler eines beendeten Kontexts geht vor, auch bei einem vergangenen Zeitpunkt; darum ruft der Service nach dem ersten Abbruchsignal `WarteBis` nicht mehr auf (jeder Aufruf lieferte den Fehler), und `--finish-session-on-interrupt` spielt ohne Wartezeit zu Ende. Ein ferner Zeitpunkt (`math.MaxInt64` nach der Begrenzung) wartet bis zum Abbruch — entschieden: `LH-FA-21.a` *Uhr* ([ADR-0015](../../adr/0015-uhr-port.md)).
- *Abbruch einer Wartezeit:* Das erste Abbruchsignal beendet jede Wartezeit sofort, auch mit `--finish-session-on-interrupt`, danach wartet `play` nicht mehr; ohne die Option sendet es die folgende Interaktion nicht, mit ihr spielt es die übrigen Interaktionen der Session ohne Wartezeit zu Ende — entschieden, `LH-FA-21.a` *Randformen des Einspielens*. Die Empfehlung des Planners (Wartezeit endet, Session endet wie ohne Wartezeit) gilt ohne die Option; mit ihr hält das Wort „Session beenden“ aus `LH-FA-20.a` Schritt 7 die Session zu Ende.
- *Wartezeit nach einem Fehler mit `--continue-on-error`:* Der Bezug bleibt der Beginn der vorigen aufgezeichneten Interaktion, auch wenn sie mit einer Fehlerantwort endete — entschieden ebenda.
- *Beginn einer Interaktion:* der Zeitpunkt, zu dem `play` ihre erste Client-Nachricht sendet, nach Ablauf der Wartezeit — entschieden ebenda.
- *Sehr großer Wert:* keine Obergrenze und keine Warnung, die Wartezeit ist durch das Signal abbrechbar (akzeptiertes Negativ: Die Aufzeichnung und `--keep-timing` sind der Wunsch des Anwenders); der Zeitpunkt wird bei Überlauf von `int64` auf den größten Wert begrenzt (Mutation: Überlauf bleibt negativ → keine Wartezeit) — entschieden ebenda.
- *Sessions ohne Interaktion:* keine Prüfung auf `offset_ms`, keine Wartezeit — entschieden.
- *Verbindungsaufbau scheitert:* kein Warten, der Fehler gilt wie ohne Option; `connect` wird am Beginn jedes Aufbaus gelesen — entschieden ebenda.
- *Wartezeit vor der ersten Interaktion:* `relative` wartet nie, `absolute` mit `connect` bis zum `offset_ms` der ersten Interaktion, mit `first-request` nie — entschieden ebenda.
- *Gesetzt:* Eine Quelle (Option, Variable, Datei), die den Wert des Defaults nennt, zählt als gesetzt; `--timing-reference connect` mit `relative` und `--keep-timing=false` mit `--timing-mode absolute` sind `PGR-E2001`; die Prüfung liegt bei den Kombinationen nach der Zusammenführung, die Reihenfolge der Optionen ändert nichts — entschieden ebenda. Der leere Wert, `Absolute` und `first_request` sind als Option und Variable `PGR-E2001`, als Schlüssel der Datei `PGR-E2004`; Kurzformen gibt es nicht.
- *Ordnung von `PGR-E3003` und `PGR-E2003`:* Eine beschädigte Interaktion (`Validate`) ist `PGR-E3003`, Exit-Code 3, und geht `PGR-E2003` vor; gehört zum Laden (`LH-FA-20.a` Schritt 1) — entschieden vom Architect (2026-10-12, F-653), `LH-FA-21.a` *Prüfung der Werte*. `PGR-E6001` hat in `play` heute keine Prüfung der Art, die Ordnung gilt für jede künftige.
- *`PGR-E2003`:* nach dem Laden, vor `PGR-E6001` und vor jeder Verbindung; die Meldung nennt Session und Interaktion der ersten ohne `offset_ms`, keinen Wert; ohne `--keep-timing` entfällt die Prüfung — entschieden ebenda.
- *Mit `--compare-responses`:* ändert die Wartezeit nicht — entschieden ebenda.
- *Option, Variable, Schlüssel:* Namen, Defaults und Wertemenge stehen in der Optionstabelle und in `LH-FA-17.a` — entschieden; Groß- und Kleinschreibung der Werte wie bei den anderen Aufzählungen (nur Kleinbuchstaben) — entschieden ebenda.

**Risiken:**

- Zeitabstände auf einem ausgelasteten Rechner schwanken; der Test prüft nur die Untergrenze (Pausen nie kürzer), das E2E nur diese — **Ausgang:** offen bis Closure.
- Das Warten hängt am Zeitverhalten des Clients oder Servers im Replay (`BEO-REPO/replay-haengt-am-zeitverhalten-des-clients`, 1×, offen): Der E2E-Test wartet auf Zeit, die nur die Fake-Uhr des Service-Tests deterministisch macht — **Ausgang:** offen bis Closure.

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

**Läufe:** Die Mutanten liefen in frischen Kopien des Baums (`cp -r` ohne `-p`, ohne `.git`, im Scratchpad, nicht im Repo; je Mutant eine eigene Kopie, danach entfernt), die Unit-Mutanten mit `go test` im Go-Image des `Dockerfile` ohne Zwischenstand und die Kopie als Bind-Mount, die Mutanten des Binaries mit `make test-integration` in der Kopie (der `ENTRYPOINT` der Kopie auf die Tests dieses Slice beschränkt). Dies sind 58 Mutanten (54 auf Unit-Ebene: 31 im Service, 22 im CLI-Adapter, 1 im Bootstrap; 4 am Binary); die Spalte „roter Test“ nennt die Tests, die zur Zusage gehören, nicht jeden roten Test des Laufs (beim Mutanten „Prüfung läuft immer“, bei „Fehlerantwort bricht nicht ab“ und bei „Umgebung vor Option“ sind es mehr).

**Gate-Läufe:** `make gates` lief grün am Stand `61a63d5` auf sauberem Baum, darin `make test`, `make lint`, `make a-check`, `make a-check-negativ`, `make abdeckung-check` und `make test-integration` mit den Tests dieses Slice (`TestE2EPlayKeepTiming`, `TestE2EPlayKeepTimingQuellen`, `TestE2EPlayKeepTimingUngueltig`, `TestE2EPlayKeepTimingOhneOffset`, `TestE2EPlayKeepTimingHilfeUndConfigShow`, `TestE2EReplayIgnoriertZeitangaben`); `make kopf-check` und `make docs-check` ebenso. Spätere Commits ändern nur diesen Absatz.

**Zusage · Mutation · roter Test** (Wortlaut `LH-FA-21.a` *Einspielen* und *Randformen des Einspielens*, Optionstabelle von `LH-FA-17.a`):

| Zusage | Mutation (nur sie) | roter Test |
|---|---|---|
| `relative`: wartet, bis seit dem **Beginn** der vorigen Interaktion mindestens der Unterschied ihrer beiden `offset_ms` vergangen ist (Beginn = Senden ihrer ersten Client-Nachricht, nicht das Ende ihrer Antwort) | Bezug ist die Lesung beim Warten, also das Ende der Antwort der vorigen | `TestPlayZeitRelative` (B wartet bis 1007 statt 1000), `TestPlayZeitNachFehlerantwort` |
| `relative`: eine Verspätung wird nicht aufgeholt | Bezug ist der Beginn der ersten Interaktion, `relative` holt auf | `TestPlayZeitRelativeVerspaetung` (C wartet nicht), `TestPlayZeitNachFehlerantwort` |
| `relative`: vor der ersten Interaktion einer Session wartet sie nie | die erste Interaktion wartet bis zu ihrem `offset_ms` | `TestPlayZeitRelativeVerspaetung` (A, `offset_ms` 5000, nicht gesendet), `TestPlayZeitSessionRelative` |
| `relative`: gleicher `offset_ms` hintereinander wartet nicht | (kein Mutant, der nur sie bricht: Die Rechnung `Beginn + 0` kehrt mit jeder Uhr sofort zurück; der Fall D in `TestPlayZeitRelativeVerspaetung` zeigt, dass der Lauf endet) | grün, äquivalent: nicht beobachtbar außer durch das Ende des Laufs |
| `relative`: auch vor einer Extended-Interaktion | keine Wartezeit vor `RequestExtended` | `TestPlayZeitRelative` (die dritte Interaktion ist `gruppe sync`) |
| `absolute`: wartet bis zum Bezugspunkt plus aufgezeichnetem Abstand, eine Verspätung wird aufgeholt (der Abstand kann kürzer sein) | der Modus ist `relative` | `TestPlayZeitAbsoluteConnect` (C folgt B nach 10 statt 50 ms), `TestPlayZeitSessionAbsolute`, `TestPlayZeitGrosserWert/absolute,_connect` |
| `connect` ist der Beginn des Verbindungsaufbaus, vor TLS und Anmeldung | `neuerTakt` nach `Verbinde` gelesen | `TestPlayZeitAbsoluteConnect` (A bei 150 statt 100; der Aufbau dauert 50 ms) |
| `connect` wird am Beginn jedes Aufbaus gelesen (jede Session eigene Uhr) | einmal je Lauf gelesen | `TestPlayZeitSessionAbsolute` (C ohne Wartezeit) |
| `relative`: jede Session eigene Uhr | der Takt wird über die Sessions geteilt | `TestPlayZeitSessionRelative`, `TestPlayZeitSessionAbsolute` |
| `absolute` mit `connect` wartet vor der ersten Interaktion bis zu ihrem `offset_ms` | die erste Interaktion wartet nie | `TestPlayZeitAbsoluteConnect` (A wird vor 100 gesendet), `TestPlayZeitSessionAbsolute`, `TestPlayZeitGrosserWert/absolute,_connect` |
| `absolute` mit `first-request` wartet vor der ersten Interaktion nie | die erste wartet bis `connect` + `offset_ms` | `TestPlayZeitAbsoluteErsteAnfrage` (A nicht ohne Vorrücken gesendet) |
| `first-request` ist der Beginn der ersten Interaktion der Session (aufgezeichnet deren `offset_ms`) | der Bezug ist die Lesung am Aufbau statt des Beginns der ersten | `TestPlayZeitAbsoluteErsteAnfrage` (B kommt bei 249 statt 250) |
| `--timing-reference first-request` wählt den Bezugspunkt | im Service gilt immer `connect` | `TestPlayZeitAbsoluteErsteAnfrage` |
| der Bezug von `relative` ist der Beginn der vorigen aufgezeichneten Interaktion, auch nach einer Fehlerantwort mit `--continue-on-error` | nach einer weitergereichten Fehlerantwort gilt der Beginn der Interaktion davor | `TestPlayZeitNachFehlerantwort` (C sofort statt bei 210); der Mutant „Bezug = Ende“ ebenso |
| Scheitert der Aufbau, gibt es keine Wartezeit | die Wartezeit der ersten Interaktion liegt vor `Verbinde` | `TestPlayZeitAufbauScheitert` (ein Aufruf von `WarteBis` statt keiner), `TestPlayZeitAbsoluteConnect`, `TestPlayZeitAbbruchMitFinish` |
| Bricht das Einspielen mit einem Fehler ab, gibt es keine weitere Wartezeit | eine Fehlerantwort bricht nicht ab | `TestPlayZeitKeineWartezeitNachAbbruch` (zwei Aufrufe von `WarteBis`; rot auch `TestPlayFehlerantwort` und weitere des Bestands) |
| Abbruch: das erste Signal beendet jede Wartezeit sofort | die Wartezeit hängt nicht an `ctx` | `TestPlayZeitAbbruchWartezeit`, `TestPlayZeitAbbruchMitFinish` (Lauf endet nicht binnen 10 s) |
| Abbruch ohne `--finish-session-on-interrupt`: die folgende Interaktion wird nicht gesendet | die unterbrochene Interaktion wird gesendet | `TestPlayZeitAbbruchWartezeit` (`anfrage B` gesendet) |
| Abbruch mit `--finish-session-on-interrupt`: die übrigen Interaktionen der Session ohne Wartezeit zu Ende (die unterbrochene eingeschlossen) | die unterbrochene Interaktion wird nicht gesendet | `TestPlayZeitAbbruchMitFinish` |
| nach dem ersten Signal wird nicht mehr gewartet (der Service ruft `WarteBis` nicht mehr auf) | `warte` prüft `ctx` nicht vor dem Aufruf | `TestPlayZeitAbbruchMitFinish` (`zaehlUhr` zählt einen Aufruf nach dem Signal) |
| das zweite Signal beendet jede Wartezeit sofort | die Wartezeit hängt nicht an `abbruch` | `TestPlayZeitZweitesSignalInWartezeit` (Lauf endet nicht binnen 10 s) |
| die Verbindung ist geschlossen, wenn `Play` zurückkehrt (Terminate) | `session` wartet nicht auf das Schließen durch das zweite Signal | `TestPlayZeitZweitesSignalInWartezeit` (`schliesse` fehlt, das Schließen dauert 100 ms) |
| große Werte: keine Obergrenze, der Zeitpunkt wird bei Überlauf auf den größten Wert begrenzt, nie negativ | Überlauf bleibt negativ | `TestPlayZeitGrosserWert` (alle drei Betriebsarten: `anfrage B` beziehungsweise `anfrage A` gesendet) |
| Ohne `--keep-timing` wartet `play` nie und liest die Uhr nie; die Uhr ist dann nicht nötig | der Takt gilt mit Uhr auch ohne `KeepTiming` | `TestPlayZeitOhneKeepTiming` (Lauf endet nicht binnen 10 s, Lesungen 1) |
| Mit `--keep-timing` hält `play` die Pausen am Binary ein (Untergrenze, beide Formate, jede Betriebsart) | der Takt ist nie aktiv | `TestE2EPlayKeepTiming`, `TestE2EPlayKeepTimingQuellen` (Abstände wenige ms statt mindestens 350) |
| Der Bootstrap verdrahtet die Systemuhr im Play-Pfad | `services.Wartezeiten(NeueSystemuhr())` entfernt | `TestRunPlayKeepTiming` (Panik der fehlenden Uhr, rot) |
| `PGR-E2003`: jede Interaktion wird geprüft | nur die erste jeder Session | `TestPlayZeitFehltOffset` (spätere Interaktion, spätere Session, zwei fehlen) |
| `PGR-E2003`: jede Session mit Interaktionen | nur die erste Session | `TestPlayZeitFehltOffset` (spätere Session) |
| `PGR-E2003`: die Meldung nennt die erste ohne Wert in der Reihenfolge der Sessions und Interaktionen | die letzte statt der ersten | `TestPlayZeitFehltOffset` (zwei fehlen) |
| `PGR-E2003`: die Meldung nennt keinen Wert | die Meldung nennt den Wert der vorigen Interaktion | `TestPlayZeitFehltOffset` (`4242` in der Meldung) |
| `PGR-E2003`: ohne `--keep-timing` entfällt die Prüfung | die Prüfung läuft immer | `TestPlayZeitFehltOffset` (ohne `--keep-timing`), `TestPlayZeitOhneKeepTiming` |
| `PGR-E2003`: eine Session ohne Interaktion verlangt nichts | eine Session ohne Interaktion ist `PGR-E2003` | `TestPlayZeitFehltOffsetRand` |
| `PGR-E2003`: der Code und die Klasse Konfiguration (Exit-Code 2) | Code `PGR-E3009` | `TestPlayZeitFehltOffset` (Kopf `Konfiguration [PGR-E2003]`) |
| `PGR-E2003` am Binary: Startfehler vor der ersten Verbindung, eine Zeile ohne Log-Zeile | die Prüfung gibt immer `nil` | `TestE2EPlayKeepTimingOhneOffset`, am Bootstrap `TestRunPlayKeepTimingOhneOffset` |
| Ein Fehler beim Laden (`PGR-E3003`) geht `PGR-E2003` vor (Schritt 1 vor der Prüfung der Werte) | die Prüfung steht vor der der Interaktionen | `TestPlayZeitFehltOffsetRand` |
| `replay` ignoriert Zeitangaben immer | `ReplayService.Query` wartet auf `offset_ms` | `TestE2EReplayIgnoriertZeitangaben` (die Aufzeichnung steht auf eine Stunde, der Lauf endet nicht in 30 s) |
| Option: `--timing-mode` ohne `--keep-timing` ist `PGR-E2001` | die Prüfung entfällt für den Modus | `TestPlayZeitVerwendung` (fünf Fälle, Option, Umgebung, Datei) |
| Option: `--timing-reference` ohne `--keep-timing` ist `PGR-E2001` | die Prüfung entfällt für den Bezug | `TestPlayZeitVerwendung` (vier Fälle) |
| Option: `--timing-reference` zusammen mit `relative` ist `PGR-E2001`, auch mit dem Standardmodus | die Prüfung entfällt | `TestPlayZeitVerwendung` (vier Fälle) |
| „Gesetzt“ ist, was irgendeine Quelle nennt, auch den Wert des Standards (Modus) | „gesetzt“ heißt: vom Standard verschieden | `TestPlayZeitVerwendung` (`--timing-mode=relative` ohne `--keep-timing`, Umgebung, Datei) |
| „Gesetzt“ (Bezug) mit dem Standard `connect` | dasselbe für den Bezug | `TestPlayZeitVerwendung` (`--timing-reference connect` mit `relative`) |
| `--keep-timing=false` mit `--timing-mode absolute` ist `PGR-E2001` | dieselbe Mutation wie in der Zeile zu `--timing-mode` ohne `--keep-timing` | `TestPlayZeitVerwendung` (`--keep-timing=false` mit Modus, umgekehrte Reihenfolge) |
| Die Reihenfolge der Optionen ändert nichts | (äquivalent: Die Prüfung liegt in `zuletzt`, nach der Zusammenführung aller Quellen; ein Mutant, der die Reihenfolge der Kommandozeile einbezieht, hat in dieser Struktur keine Stelle) | `TestPlayZeitVerwendung` (Fall „Reihenfolge umgekehrt“), `TestPlayZeitWerte` (Fall „Reihenfolge der Optionen“) grün ohne Mutant |
| Die Meldung der Verwendung nennt Optionen, keinen Wert aus Umgebung oder Datei | die Meldung hängt den Wert des Modus an | `TestPlayZeitVerwendung` (Meldung genau) |
| Werte nur in Kleinbuchstaben: `Absolute` ist `PGR-E2001` (Option, Variable), `PGR-E2004` (Schlüssel) | Gleichheit ohne Beachtung der Schreibweise | `TestPlayZeitWerteUngueltig` |
| … `first_request` | der Unterstrich ist erlaubt | `TestPlayZeitWerteUngueltig` |
| … der leere Wert (Modus, Bezug) | leer ist erlaubt | `TestPlayZeitWerteUngueltig` |
| Die Meldung zu Umgebungsvariable und Schlüssel nennt den Wert nie | der Text der Wertemenge hängt den Wert an (Modus, Bezug, je ein Mutant) | `TestPlayZeitWerteUngueltig` (beide rot); am Binary `TestE2EPlayKeepTimingUngueltig` ohne eigenen Mutanten |
| Default `relative` und `connect` | Standardmodus `absolute`, Standardbezug `first-request` | `TestPlayZeitWerte`, `TestPlayZeitQuellen`; `TestLeserAlleOptionen` für `--keep-timing` |
| Mit `--keep-timing` allein gilt `relative` ohne Bezug; bei `absolute` gilt `connect`; ohne `--keep-timing` sind die Felder leer | der Bezug wird auch bei `relative` gesetzt, der Modus ohne `--keep-timing` nicht geleert, `first-request` wird `connect` | `TestPlayZeitWerte`, `TestParsePlay` und weitere des Bestands |
| Quellen: Option, Umgebungsvariable, Schlüssel im Abschnitt `play:` | der Schlüssel `keep_timing` wirkt nicht | `TestPlayZeitQuellen` (alle drei Optionen), `TestLeserAlleOptionen` |
| Bei `record` und `replay` sind die Optionen, Variablen und Schlüssel unbekannt | `--keep-timing` bei `record` mit angemeldet | `TestLeserOptionen`, `TestLeserAlleOptionen`, `TestDateiUngueltig`, `TestConfigShowFehler`; am Binary `TestE2EPlayKeepTimingUngueltig` |
| `--help` nennt die Optionen, Werte, Standards, Variablen und `PGR-E2001`/`PGR-E2003` | der Name `--timing-reference` im Hilfetext verändert; der Absatz zu `PGR-E2003` entfernt; „Standard false“ zu „Standard true“ bei `--keep-timing`; „Beginn der ersten Interaktion“ zu „Ende der ersten Interaktion“ bei `--timing-reference`; „Beginn der vorigen Interaktion“ zu „Ende“ bei `--timing-mode` (je ein Mutant, nach dem Review F-655) | `TestPlayZeitHilfe` (beide rot); am Binary `TestE2EPlayKeepTimingHilfeUndConfigShow` |
| `config show` nennt die Schlüssel der Datei und die Namen gesetzter Variablen ohne Wert | (kein Code dieses Slice: `anzeige` zeigt die Datei ohne Kenntnis einzelner Optionen; der Test belegt das Verhalten für die neuen Schlüssel) | `TestE2EPlayKeepTimingHilfeUndConfigShow`, grün ohne Mutant |
| Keine Kurzformen der Optionen | (nichts geliefert: Der Slice meldet keine Kurzform an) | kein Test |
| Priorität Option vor Variable vor Datei (`--keep-timing=false` geht Variable und Datei vor) | die Umgebungsvariable geht der Option vor | `TestPlayZeitQuellen`, `TestLeserAlleOptionen` und weitere des Bestands; am Binary `TestE2EPlayKeepTimingQuellen` (`Option false vor Datei`, `Option false vor Umgebung`), dafür kein eigener Mutant am Binary |

**Nachtrag nach dem Review (F-654 bis F-659):**

- *F-654, Werte außerhalb der Aufzählung im Service:* Der Kommentar von `PlayOptions` sagt nichts mehr über unbekannte Werte zu. Die Spezifikation macht sie am Start zu `PGR-E2001`, der Service sieht nur gültige Werte (`TimingRelative`/`TimingAbsolute`, `ReferenceConnect`/`ReferenceFirstRequest`; der CLI-Adapter lässt den Bezug bei `relative` leer). Die Mutanten `TimingMode != TimingRelative` statt `== TimingAbsolute` und `TimingReference != ReferenceConnect` statt `== ReferenceFirstRequest` bleiben grün und sind **äquivalent** für jeden gültigen Wert (beide Zweige stimmen für alle zugelassenen Werte überein); ein Test für einen ungültigen Wert würde eine Zusage prüfen, die die Spezifikation nicht macht.
- *F-655, Hilfetext:* `TestPlayZeitHilfe` prüft jetzt „Standard false“ bei `--keep-timing`, die Beschreibung der beiden Modi und die der beiden Bezugspunkte im Wortlaut; die drei Textmutanten oben sind rot.
- *F-658, `KeepTiming` ohne `Wartezeiten`:* akzeptiertes Negativ. Es gibt genau einen Aufrufer, den Bootstrap, und `TestRunPlayKeepTiming` wird rot, wenn er die Option nicht reicht; ein Fehlerpfad im Konstruktor wäre Code ohne Aufrufer, der ihn auslöst. Der Kommentar von `Wartezeiten` nennt die Pflicht des Aufrufers und dass der Lauf sonst in einem Nil-Zugriff endet (nicht zugesagt).
- *F-657, Flackern des Binary-Tests:* Die einzige Richtung, in der `TestE2EPlayKeepTiming` flackern kann, ist `absolute` mit `connect`: Der Anker ist `backend_start`, der nach dem Beginn des Aufbaus liegt, und die 50 ms Toleranz decken diesen Abstand; auf einem stark belasteten Runner kann er sie überschreiten (nicht beobachtet). Akzeptiertes Negativ, die Untergrenze selbst kann Last nicht verletzen.
- *F-659:* der Kommentar von `session` ist gleichmäßig umbrochen.

**Grüne Mutanten:**

- *gleicher `offset_ms` hintereinander* und *Reihenfolge der Optionen*: äquivalent, siehe die beiden Zeilen. Ein Test ändert daran nichts, die Grenze trägt der Aufbau (die Prüfung liegt nach der Zusammenführung).
- *`config show`* und *Kurzformen*: kein Mutant, weil dieser Slice dafür keinen Code schreibt; die Zeilen stehen, damit die Aussage der DoD nicht als geprüft gilt, wo sie es durch einen Mutanten nicht ist.
- Die E2E-Untergrenze trennt `first-request` nicht von `connect` (Mutant „der Bezug ist immer `connect`“ blieb am Binary grün, weil die Untergrenze zum Backend-Beginn stärker ist als die des Beginns der ersten Interaktion); die Zusage tragen `TestPlayZeitAbsoluteErsteAnfrage`, `TestPlayZeitWerte` und `TestPlayZeitQuellen` (rot für „`first-request` wird `connect`“ und „Bezug aus der Lesung am Aufbau“).

**Gegenbeispiele, die der Code nicht verhindert (akzeptiertes Negativ, zur Closure zu führen):** Die Untergrenze am Binary ist auf der Server-Uhr gemessen (`clock_timestamp()`) und trägt eine Toleranz von 50 ms (`toleranz` in `zeitangaben_einspielen_e2e_test.go`), weil Senden und Ausführung streuen; eine Obergrenze prüft kein Test. Bei `absolute` mit `connect` liegt der Beginn des Backends (der Anker der Messung) nach der Lesung am Beginn des Aufbaus; die Toleranz deckt den Abstand.

**Ordnung der Prüfungen (entschieden vom Architect, 2026-10-12, F-653/F-656):** `PGR-E3003` einer beschädigten Interaktion geht `PGR-E2003` vor, weil `Validate` zum Laden gehört (`LH-FA-20.a` Schritt 1); `LH-FA-21.a` *Prüfung der Werte* sagt es jetzt. Der Satz „vor `PGR-E6001`“ bleibt: die Prüfung der Art ist in `play` ohne Fall, die Ordnung gilt für jede künftige Art. Code und `TestPlayZeitFehltOffsetRand` bleiben unverändert (Mutation „die Prüfung steht vor der der Interaktionen“, Zeile in §7 oben).

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
- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (1×, offen) — berührt, weil das Einspielen auf Zeit wartet (Risiko in §6).
- `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (2×, offen) — Dieser Plan ist der Rest des Schnitts vom 2026-10-11 mit drei Liefer-Punkten und zwei Schichten; die Größe misst der Architect vor dem Code, und die Rückführung in §4 nennt die Bedingung vorab.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 3. Schichten: zwei (Kern mit dem Play-Service, CLI-Adapter); der Bootstrap reicht nur die Uhr weiter und zählt als reine Verdrahtung nicht (`AGENTS.md` §3.13). Nachgezählt vom Implementer am Diff (2026-10-12): Liefer-Punkte 3 und Schichten 2 wie geplant; der Bootstrap trägt eine Zeile (`services.Wartezeiten(NeueSystemuhr())`) ohne Verzweigung, der Meldungscode `PGR-E2003` liegt im Kern. Nachgezählt beim Eintragen der Signatur des Ports (Planner, 2026-10-12, Closure von `slice-v1-abschluss-zeitangaben-uhr`): Liefer-Punkte 3 unverändert (die Namen stehen in DoD-Punkt 1, die Zusagen des Ports am Warten in §6), Schichten 2; `NeueSystemuhr()` aufrufen und reichen ist Verdrahtung, `ctx` und `ablauf` reicht der Bootstrap schon heute. Nachgezählt beim Schnitt und vom Architect vor dem Code (2026-10-11): Der Plan stand vorher bei drei Liefer-Punkten und vier Schichten (Kern mit Uhr-Port, PGWire-Adapter, Recording-Adapter, CLI-Adapter; die Dokumentation ausgenommen); nach dem Schnitt hat jeder der vier Teile höchstens drei Liefer-Punkte und höchstens zwei Schichten (`-pruefen` 2 und 1, mit Helfer im Kern 2; `-uhr` 3 und 2; `-aufzeichnen` 2 und 2; dieser Slice 3 und 2). Die Nachzählungen früherer Sendungen — Schichtzählung aus `slice-doku-ist-stand` (2026-10-10) und Sendung aus `slice-v1-abschluss-sqlite-format-lesen` (2026-10-11) — galten dem ungeschnittenen Plan; die Sendung aus `slice-v1-abschluss-sqlite-format-lesen` steht jetzt in diesem Plan (§1, Einspielen aus beiden Formaten) und in `-aufzeichnen` (Schreiben, Roundtrip) und `-pruefen`.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
