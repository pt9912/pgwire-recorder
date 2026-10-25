# Slice slice-play-zweites-signal-terminate-wettlauf: Das zweite Signal sendet `Terminate`, sobald die Nachricht vollständig geschrieben ist — auch bei gehaltener Bearbeitung der Anfrage

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist (`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht). Angelegt vom Planner bei der Closure von `slice-test-unter-last-ursache` (Risiko-Ausgang *eingetreten*); die Stelle in der Reihenfolge ist Empfehlung des Planners, nicht Entscheidung des Nutzers. Der Test `TestRunPlayZweitesSignal` bleibt bis zur Lieferung unter CPU-Last rot-reproduzierbar.

**Bezug:** `LH-FA-20` (Einspielen), `LH-FA-20.a` (*Abbruchsignal*: das zweite Signal sendet `Terminate`). Auslöser: [slice-test-unter-last-ursache](../done/slice-test-unter-last-ursache.md) §7 (Ursache von `TestRunPlayZweitesSignal`: Wettlauf im Code; Probe P1 5 von 5 rot ohne Last). Register: `BEO-REPO/test-faellt-unter-cpu-last-aus` (die Ursache von C steht in dessen `state.md`). Keine aktive ADR.

**Berührte Spec-Stellen:** `spezifikation.md` §`LH-FA-20.a — Einspielen`, Stichwort *Abbruchsignal* (geschärft in Commit 36b6387c; Randformen Granularität, Gesamtfrist, Gründe und Plattform vom Architect vorab in f17d9027 entschieden; dieser Slice ändert die Spezifikation nicht mehr, er bringt den Code auf ihren Wortlaut) · `LH-FA-20.a`

**Verantwortlich:** Implementer (pt9912)

**Autor:** Planner. **Datum:** 2026-10-24.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das zweite Signal von `play` sendet `Terminate`, sobald die Nachricht des Clients vollständig geschrieben ist, auch wenn der Sender seine Bearbeitung der Anfrage noch nicht abgeschlossen hat (Wortlaut von `LH-FA-20.a` *Abbruchsignal* seit 36b6387c): Im Einspiel-Adapter (`internal/adapters/driven/postgres/einspielen.go`) ersetzt ein Zustand „Nachricht teilweise geschrieben“ die Sperre `schreiben` als Stellvertreter für „der Sender schreibt“; `schliesse` fragt den Zustand und nicht die Sperre (`TryLock`). Die Wartezeit des Schreibens in `schliesse` bleibt durch `terminateFrist` (`upstream.go:107`) begrenzt, ein unbegrenztes `Lock()` ist unzulässig. `TestRunPlayZweitesSignal` (`internal/bootstrap/play_test.go`) ist der rote Test der Zusage und fällt mit der Mutation „die Sperre bleibt 300 ms nach dem Senden gehalten“ (Probe P1 von `slice-test-unter-last-ursache`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine **Änderung der Spezifikation oder des Lastenhefts** im Slice — Schicht-Abgrenzung und Rang: Die Schärfung steht seit 36b6387c in `LH-FA-20.a`, die Randformen 2 bis 5 hat der Architect vor dem Anspruch in f17d9027 ergänzt (Vorab-Änderung, nicht Teil dieses Slice); eine Randform, die sie nicht trägt, geht an den Architect zurück (`AGENTS.md` §3.12), sie entscheidet der Implementer nicht.
- Eine **Änderung an `TestE2EPlayVergleichAbweichung` oder `TestE2EPlayVergleichAlle`** — anderer Vorgang: Deren Rot ist nicht reproduziert und nicht erklärt (`BEO-REPO/test-faellt-unter-cpu-last-aus`), die Ursache ist keine Wettlaufbedingung dieses Codes.
- Eine **Änderung an Fristen anderer Tests oder an `terminateFrist`** — Bestand bleibt bewusst stehen: Die Frist von 100 ms ist die Zusage („kurze Frist“); der Slice nutzt sie, er ändert sie nicht.
- Das **Terminate in `Close` des Upstream-Adapters** (`upstream.go`, Pfad von `record`) — Bestand bleibt bewusst stehen: Der Wettlauf stammt aus `einspielSession`; `record` trägt andere Sperren.
- **Benutzerhandbuch und README** — Schicht-Abgrenzung: kein Verhalten ändert sich, das ein Anwender sieht, außer dass `Terminate` zuverlässig ankommt; das Handbuch nennt das Verhalten nicht und wird nicht beschrieben (§3.11).

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Liefer-Punkt 1 (Zustand statt Sperre):** `schliesse` (in `einspielen.go`) sendet `Terminate`, wenn keine Nachricht des Clients teilweise geschrieben ist, und sonst nicht; die Sperre `schreiben` ist nicht mehr der Stellvertreter. Die Wartezeit von `schliesse` auf das Ende eines teilweisen Schreibens ist durch `terminateFrist` begrenzt (Randformen 1 bis 3).
- [x] **Liefer-Punkt 2 (Tests mit Mutation, `AGENTS.md` §3.10):** je Zusage ein Test, der rot wird, wenn der Code sie bricht, und die Mutation vom Implementer selbst rot gesehen (Zusage · Mutation · roter Test im Bericht): (a) Terminate bei gehaltener Bearbeitung der Anfrage (Sperre 300 ms nach dem Senden gehalten): `TestRunPlayZweitesSignal` rot; (b) kein Terminate bei teilweise geschriebener Nachricht (bestehend: `TestEinspielSchliesseSchreibtNichtBeimSenden`); (c) die Frist begrenzt das Warten (ein Test, der ohne Frist, mit unbegrenztem `Lock()`, rot wird); (d) die Reihenfolge der Gründe (Randform 4), je Element ein Fall.
- [x] **Liefer-Punkt 3 (Plattform und Abdeckung):** Die Plattformentscheidung (Randform 5, in f17d9027 entschieden: geprüft nur unter Linux) ist in der Spezifikation getragen und der Abdeckungs-Deklaration der neuen Tests nachgezogen (`make abdeckung-check` grün); die Tests tragen `Abdeckung: LH-FA-20/…`.
- [x] Kein Verhalten der Kommandozeile außer dem zweiten Signal ändert sich; `TestRunPlayZweitesSignal` ist unter der künstlichen Last aus `slice-test-unter-last-ursache` §7 (K=32, K=48 Container mit Rechenschleife) in mindestens 20 Läufen je Last grün, sonst ist die Ursache nicht behoben.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben (kein Zähler wird gesetzt); keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/postgres/einspielen.go` | update | `einspielSession`: Zustand „teilweise geschrieben“ (unter `zustand` oder eigener Wert), gesetzt für die Dauer des Schreibens einer Nachricht, `schliesse` fragt ihn; Kommentar des Typs und von `Schliesse` nachziehen (§3.7) |
| `internal/adapters/driven/postgres/einspielen_ende_test.go`, `einspielen_test.go` | update | Fälle je Zusage und Mutation (Liefer-Punkt 2), Abdeckung `LH-FA-20.a` *Abbruchsignal* |
| `internal/bootstrap/play_test.go` | unverändert gefahren | `TestRunPlayZweitesSignal` ist der rote Test (a); Änderung nur, wenn der Architect es entscheidet |
| `spec/spezifikation.md` | unverändert | Randformen 1 bis 5 stehen seit 36b6387c und f17d9027 (Architect, vor dem Anspruch); Änderung im Slice nur über Rückgabe an den Architect |

Das Schreiben einer Nachricht läuft über `fe.Send` und `fe.Flush`; der Implementer prüft, ob der Zustand je Nachricht oder je `Flush` gilt (Randform 2) und nennt es dem Architect, wenn die Spezifikation es nicht entscheidet.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert (`open` → `next`, `Verantwortlich:` gesetzt), dann der reine Move auf `main`. Voraussetzung: `slice-test-unter-last-ursache` in `done/` (Ursache belegt). Kein anderer Slice ändert `einspielen.go` zur selben Zeit (`slice-v1-abschluss-zeitangaben` ist geschlossen oder liegt nicht in `in-progress/`; der Planner prüft es beim Anspruch).

**Größe.** Metrik und Befehl stehen in `AGENTS.md` §3.14. **Pfade des Plans:** `internal test` (mit den zwei Ausschlüssen des Befehls). **Marke 1500** (`AGENTS.md` §3.14). **Basis:** der Commit, der diese Datei nach `in-progress/` bewegt; sein Hash steht hier beim Anspruch: `dced44df`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik erreicht die Marke 1500; oder die Randform 2 (Zustand je Nachricht statt je `Flush`) verlangt eine Änderung an `pgproto3`-Nutzung außerhalb von `einspielSession`.
- `in-progress` → `open` (blockiert): Der Architect entscheidet eine Randform nicht, oder die Plattformentscheidung (Randform 5) verlangt einen Test unter macOS oder Windows, den das Gerüst nicht bietet.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; in §7 stehen je Zusage Mutation und roter Test, die Läufe von `TestRunPlayZweitesSignal` unter Last, der Beleg `make gates` grün; Review-Report liegt vor; Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12; Geltung: Verhalten von `play`, Record-/Replay-Pfad, Wire-Grenze). Je Randform, wo sie entschieden ist: *Spezifikation* heißt, die Stelle steht seit 36b6387c in `LH-FA-20.a` *Abbruchsignal*; *Architect* heißt, sie ist vor dem Code in der Spezifikation zu entscheiden, der Implementer entscheidet sie nicht und gibt sie zurück.

1. *Zustand statt Sperre:* „teilweise geschrieben“ ist der Zustand; die Sperre `schreiben` ist kein Stellvertreter. Terminate bleibt aus, solange eine Nachricht teilweise geschrieben ist (sie wird nicht zerrissen). Entschieden: Spezifikation.
2. *Granularität des Zustands:* je Nachricht oder je `Flush` einer Gruppe, und was gilt zwischen zwei vollständig geschriebenen Nachrichten einer Gruppe. Entschieden: Spezifikation (f17d9027): je Schreibvorgang (`Flush`), nicht je Nachricht; zwischen den Nachrichten einer Gruppe gibt es keine Stelle für `Terminate`.
3. *Frist:* Das Schreiben von `Terminate` und das Warten auf ein teilweises Schreiben sind durch `terminateFrist` (`upstream.go:107`, 100 ms) begrenzt; ein unbegrenztes `Lock()` ist unzulässig; reißt die Frist, bleibt `Terminate` aus und die Verbindung wird geschlossen, ohne Meldung. Entschieden: Spezifikation („binnen einer kurzen Frist“); die Zahl ist die Konstante; seit f17d9027 eine Gesamtfrist von 100 ms ab dem zweiten Signal für Warten und Schreiben zusammen, eingereihte ungeschriebene Gruppen werden danach nicht mehr geschrieben.
4. *Reihenfolge der Gründe, aus denen `Terminate` ausbleibt* (`AGENTS.md` §3.12 *Stelle einer Nachricht*): bereits geschlossen (kein zweites Schließen) → `SchliesseOhneTerminate` → Nachricht teilweise geschrieben → Frist gerissen; je Grund ein Fall. Stellen der Nachricht: vor dem ersten Byte einer Anfrage, mitten in der Nachricht, nach dem letzten Byte der Nachricht bei noch laufender Bearbeitung (der Fall dieses Slice), nach Rückkehr von `Anfrage`; je Stelle eine Entscheidung. Entschieden: Spezifikation (f17d9027): `Terminate` bleibt aus bei erstem Schließen schon geschehen, `SchliesseOhneTerminate`, gerissener Gesamtfrist; „teilweise geschrieben“ ist kein eigener Grund, es wird bis zur Frist gewartet, danach ist es der Grund „Frist gerissen“. Warten und Schreiben als Grund der Frist sind nicht unterscheidbar (akzeptiertes Negativ, ohne Folge); ein Fall je Grund, die Reihenfolge der letzten zwei ist nicht beobachtbar. Stellen: vor dem ersten Byte, nach dem letzten Byte bei laufender Bearbeitung und nach Rückkehr senden `Terminate`; mitten im Schreibvorgang wartet es bis zur Frist.
5. *Plattform* (`LH-QA-03`): Die Eigenschaft hängt vom Schreiben auf eine TCP-Verbindung und von der Zeitplanung ab. Geprüft: Linux (Container). macOS und Windows: Grenze in der Spezifikation an der Zusage („geprüft nur unter Linux“) oder eine Zusage mit Test; ein Handbuchsatz nennt nur die geprüfte Plattform. Entschieden: Spezifikation (f17d9027): geprüft nur unter Linux, Grenze für macOS und Windows an der Zusage; ein Handbuchsatz entfällt (§1).
6. *Aufruf-Voraussetzung:* `Anfrage`, `Gruppe` und `Schliesse` kommen aus verschiedenen Goroutinen gleichzeitig; die Zusage gilt dafür. Aufrufarten des Treibers: nicht berührt.
7. *Mutation* (§3.10): die Sperre bleibt 300 ms nach dem Senden gehalten, in einer Kopie und im Test als Wartezeit des Schreibens nachgebaut (Probe P1 von `slice-test-unter-last-ursache`); der rote Test ist `TestRunPlayZweitesSignal`, ein deterministischer Test im Adapter ergänzt ihn, weil der Wettlauf dort ohne Netz und Zeitgrenze von 30 s fällt.

**Risiken:**

- Der Zustand „teilweise geschrieben“ ist je `Flush` zu grob: Eine Gruppe mit mehreren Nachrichten hält ihn, und `Terminate` käme bei großen Gruppen später als nach der Spezifikation — **Ausgang:** offen bis Closure (Randform 2, Architect).
- Der Fix verlagert den Wettlauf, ohne ihn zu beseitigen (`TestRunPlayZweitesSignal` bleibt unter Last in einzelnen Läufen rot) — **Ausgang:** offen bis Closure; Beleg: mindestens 20 Läufe je Last (DoD).
- Ein Test mit Zeitgrenze gerät selbst unter Last ins Rot (`BEO-REPO/test-faellt-unter-cpu-last-aus`) — **Ausgang:** weiter offen → `BEO-REPO/test-faellt-unter-cpu-last-aus`.
- Der Slice erreicht die Marke von §4 (1500, Metrik `AGENTS.md` §3.14) — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Die Closure-Notiz (Lerneintrag, Risiko-Ausgänge, Register) schreibt der Planner bei Closure, vor dem `git mv` nach `done/`.

### Belege des Implementers

**Umsetzung** (Commit `b0a7e10a`, Basis `dced44df`): `einspielSession` führt `schreibEnde chan struct{}` unter `zustand`; ungleich `nil` heißt „ein Schreibvorgang (ein `Flush` von `Anfrage` oder einer Gruppe) läuft“. `beginneSchreiben` prüft `geschlossen` und setzt den Zustand in einem Schritt, `endeSchreiben` löscht ihn und weckt `schliesse`. `schliesse` setzt `geschlossen` und liest `schreibEnde` im selben Schritt, wartet höchstens bis zur Gesamtfrist `terminateFrist` ab dem Aufruf (Warten und `SetWriteDeadline(frist)` teilen sie) und sendet danach `Terminate`; reißt die Frist, bleibt es aus. Die Sperre `schreiben` hält nur noch `Anfrage` und den Sender auseinander. `Anfrage` nach dem Schließen liefert `PGR-E4003`, ohne zu schreiben (vorher ein Write auf die geschlossene Verbindung, Fehler gleich). Keine Randform außerhalb von §6 berührt; Randformen 1 bis 5 stehen in `LH-FA-20.a` *Abbruchsignal* (f17d9027), die Spezifikation ist unverändert.

**Mutationen** (jeweils eine Änderung an `einspielen.go`, Lauf `go test ./internal/adapters/driven/postgres/ -run Einspiel` im Image `pgwire-recorder:test` mit dem Arbeitsbaum per Bind-Mount, danach zurückgesetzt; Mutant kommt so ohne Cache an):

| Zusage (Wortlaut `LH-FA-20.a`) | Mutation | roter Test (aus dem Grund der Zusage) |
|---|---|---|
| Terminate bei vollständig geschriebener Anfrage, Bearbeitung nicht abgeschlossen | `endeSchreiben` in `Anfrage` 300 ms nach dem `Flush` | `TestEinspielSchliesseSendetNachVollstaendigerAnfrage` (kein Terminate hinter der Anfrage); dazu `TestRunPlayZweitesSignal` 30,3 s rot („play endet binnen 30 s nach dem zweiten Signal nicht“ bzw. kein `X`), ohne Mutation grün |
| dasselbe für eine Gruppe | `endeSchreiben` im Sender 300 ms nach dem `Flush` | `TestEinspielSchliesseSendetNachVollstaendigerGruppe` (daneben `TestEinspielGruppeNachrichten`, `TestEinspielGruppeUnabhaengig`) |
| kein Terminate in einen teilweise geschriebenen Schreibvorgang | `schliesse` wartet nicht (`wartePunktlich(nil, …)`) | `TestEinspielGruppeSchliesseBeimSenden` (Panic des Fakes: zweites `Write` während des ersten, `sendeSperre` und `steuerConn`), `TestEinspielSchliesseSchreibtNichtBeimSenden` |
| Frist begrenzt das Warten | Timer `time.Hour` statt `time.Until(frist)` | `TestEinspielSchliesseBeimSenden`, `TestEinspielSchliesseSchreibtNichtBeimSenden/Schliesse`, `TestEinspielGruppeSchliesseBeimSenden`, `…SchliesseBeendetSender`, `…VerwirftWarteschlange` (je 5 s: Schliesse endet nicht) |
| gerissene Frist: kein Terminate | `wartePunktlich` liefert nach Ablauf `true` | `TestEinspielGruppeSchliesseBeimSenden` (Panic: zweites `Write`) |
| Gesamtfrist für Warten und Schreiben zusammen | `SetWriteDeadline(time.Now().Add(terminateFrist))` nach dem Warten | `TestEinspielSchliesseGesamtfrist` (Frist nach dem Ende des Vorgangs plus 100 ms) |
| Grund `SchliesseOhneTerminate` | `SchliesseOhneTerminate` ruft `schliesse(true)` | `TestEinspielSchliesseOhneTerminate`, `TestEinspielSchliesseNurEinmal/SchliesseOhneTerminate,…` (2 Fälle) |
| Grund „erstes Schließen schon geschehen“ | `erstes := true` | `TestEinspielSchliesseNurEinmal/Schliesse,_Schliesse` |
| eingereihte Gruppen nach dem Signal nicht mehr geschrieben | `beginneSchreiben` prüft `geschlossen` nicht | `TestEinspielGruppeSchliesseVerwirftWarteschlange` (beide Arten) |
| nach dem Schließen schreibt `Anfrage` nichts | `Anfrage` setzt „schreibt“ ohne Prüfung von `geschlossen` | `TestEinspielAnfrageNachSchliessen` (beide Arten) |

Die Reihenfolge der letzten zwei Gründe (Teilweise geschrieben und Frist) ist laut Spezifikation nicht beobachtbar und hat keinen eigenen Fall. Zwei der Fälle (Warten entfällt, Terminate trotz Frist) werden durch den Panic des Fakes (`close of closed channel` in `einspielen_gruppe_test.go`) rot, nicht durch eine Meldung; Grund ist in beiden Fällen ein zweites gleichzeitiges `Write`. Eine Mutation, die `beginneSchreiben` ignoriert und dann `endeSchreiben` ruft (Panic `close of nil channel` in einem anderen Test), zählt nicht als Beleg für die Zeile *Anfrage nach dem Schließen*; die Zeile trägt die zweite Mutation.

**Last** (`TestRunPlayZweitesSignal` als vorab gebautes Testbinary, je Lauf ein Container `--network none`, Last wie `slice-test-unter-last-ursache` §7 (a): K Container mit `while :; do :; done` im Image `postgres:17-alpine@sha256:b0f9560a…`, 5 s Vorlauf; Host Linux, 20 Kerne, kein anderer Container): K=32: 20 Läufe, 0 rot (Last 1 min vorher 10,7, nachher 19,8); K=48: 20 Läufe, 0 rot (vorher 29,9, nachher 39,1). Zum Vergleich der alte Stand laut `slice-test-unter-last-ursache` §7: 8 von 20 und 9 von 20 rot bei 16 Kernen. Grenze: Die Hosts unterscheiden sich (16 gegen 20 Kerne), die Lastzahl blieb bei K=32 um die Kernzahl; ein Lauf des alten Codes auf diesem Host steht nicht daneben. Das Skript liegt im Scratchpad der Sitzung, nicht im Repo.

**Gates:** `make gates` am Endstand (Commit `b0a7e10a`): Exit 0, `run-integration-tests: gruen`, `test-postgres-versionen: grün`, `workflows-check: kein Befund`. `make abdeckung-check` lief vorher rot („nicht aktuell“, die vier neuen Tests), `make abdeckung` schrieb `docs/user/abdeckung-unit.md` neu, danach grün im Gate. Keine Gegenprobe gefahren: kein Gate, keine Gegenprobe und kein Vertrag unter `harness/` berührt.

**Größe (§3.14):** 252 Zeilen nach dem Befehl der Regel (Basis `dced44df`, Pfade `internal test`), Marke 1500 nicht nahe.

**Kommentare (§3.7/§3.11):** neue und geänderte Kommentare in `einspielen.go` (Typ `einspielSession`, `beginneSchreiben`, `endeSchreiben`, `Schliesse`, `schliesse`, `wartePunktlich`) beschreiben den Ist-Zustand; jede Zusage darin hat die Tests der Tabelle. Handbuch und README unverändert (§1). Der Kommentar des Ports `Schliesse` („soweit die Verbindung es sofort annimmt“) bleibt; er deckt die Gesamtfrist nicht ausdrücklich, die Spezifikation tut es.

**Beobachtungen für die Closure:** Die Prüfung von `TestEinspielSchliesseGesamtfrist` hängt an einem Spielraum von 80 ms (Schliesse muss binnen 80 ms nach dem Aufruf seine Frist gesetzt haben); unter extremer Last könnte sie rot werden (`BEO-REPO/test-faellt-unter-cpu-last-aus`, Risiko 3 in §6).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-24). Treffer in `REPO`: `BEO-REPO/test-faellt-unter-cpu-last-aus` (7×; mit `slice-test-unter-last-ursache` ist die Ursache von `TestRunPlayZweitesSignal` benannt: Wettlauf im Code; dieser Slice behebt ihn); `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (§3.10 gilt für diesen Slice: je Zusage eine Mutation). Eine Schwelle wird mit diesem Slice nicht neu erreicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-24, im Commit, der diesen Plan anlegt):** Liefer-Punkte 3 (Zustand statt Sperre | Tests mit Mutation | Plattform und Abdeckung). Schichten 2: Driven Adapter `postgres` (Einspiel-Session) mit seinen Tests und Spezifikation (Technik-Stratum, Architect; zählt als Dokumentation der Spezifikation, keine Nutzerdoku-Schicht); der Test im Bootstrap (`TestRunPlayZweitesSignal`) wird gefahren und nicht geändert. Planung zählt keine Schicht; Meldungscode-Konstanten entstehen nicht. Der Slice liegt nicht über der Grenze; Marke 1500.
