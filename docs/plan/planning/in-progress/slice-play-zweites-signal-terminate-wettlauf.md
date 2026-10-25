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

**Berührte Spec-Stellen:** `spezifikation.md` §`LH-FA-20.a — Einspielen`, Stichwort *Abbruchsignal* (geschärft in Commit 36b6387c; Randformen Granularität, Gesamtfrist, Gründe und Plattform vom Architect vorab in f17d9027 entschieden, Randform 8 auf Rückgabe in fa8a8479; dieser Slice ändert die Spezifikation nicht mehr, er bringt den Code auf ihren Wortlaut) · `LH-FA-20.a`

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

**Ziel:** Das zweite Signal von `play` sendet `Terminate`, sobald die Nachricht des Clients vollständig geschrieben ist, auch wenn der Sender seine Bearbeitung der Anfrage noch nicht abgeschlossen hat (Wortlaut von `LH-FA-20.a` *Abbruchsignal* seit 36b6387c): Im Einspiel-Adapter (`internal/adapters/driven/postgres/einspielen.go`) ersetzt ein Zustand „Nachricht teilweise geschrieben“ die Sperre `schreiben` als Stellvertreter für „der Sender schreibt“; `schliesse` fragt den Zustand und nicht die Sperre (`TryLock`). Die Wartezeit des Schreibens in `schliesse` bleibt durch `terminateFrist` (`upstream.go:107`) begrenzt, ein unbegrenztes `Lock()` ist unzulässig. `TestRunPlayZweitesSignal` (`internal/bootstrap/play_test.go`) ist der rote Test der Zusage und fällt mit der Mutation „die Sperre bleibt 300 ms nach dem Senden gehalten“ (Probe P1 von `slice-test-unter-last-ursache`). Bei der Closure zieht der Planner den Satz in `.claude/commands/implement-slice.md` Schritt 16 („Rot von `TestRunPlayZweitesSignal` ist ein Befund … bis ein Folge-Slice es behebt“) auf den Ist-Zustand nach: Der Folge-Slice ist dieser, die Bedingung des Satzes ist erfüllt (`AGENTS.md` §3.9, Harness-Artefakt); keine weitere Änderung am Command.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine **Änderung der Spezifikation oder des Lastenhefts** im Slice — Schicht-Abgrenzung und Rang: Die Schärfung steht seit 36b6387c in `LH-FA-20.a`, die Randformen 2 bis 5 hat der Architect vor dem Anspruch in f17d9027 ergänzt (Vorab-Änderung, nicht Teil dieses Slice), die Randform 8 auf Rückgabe des Implementers in fa8a8479; eine Randform, die sie nicht trägt, geht an den Architect zurück (`AGENTS.md` §3.12), sie entscheidet der Implementer nicht.
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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben (kein Zähler wird gesetzt); keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen (Lage-Prüfung nach dem Move steht aus, §7).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/postgres/einspielen.go` | update | `einspielSession`: Zustand `schreibEnde` (ein Kanal unter `zustand`, ungleich `nil` heißt: ein Schreibvorgang läuft), gesetzt für die Dauer eines `Flush` von `Anfrage` oder einer Gruppe (je Schreibvorgang, nicht je Nachricht, Randform 2); `schliesse` fragt ihn und wartet bis zur Gesamtfrist; Kommentar des Typs und von `Schliesse` nachziehen (§3.7) |
| `internal/adapters/driven/postgres/einspielen_terminate_test.go` | neu | Fälle je Zusage und Mutation (Liefer-Punkt 2), Abdeckung `LH-FA-20.a` *Abbruchsignal*; die bestehenden Tests in `einspielen_ende_test.go`, `einspielen_test.go` und `einspielen_gruppe_test.go` bleiben unverändert und fangen mehrere der Mutationen |
| `internal/adapters/driven/postgres/export_test.go` | update | Brücke `NeueEinspielSessionMitFrist` und `WartePunktlich` (reicht nur weiter): Die Tests der Reihenfolge fahren mit einer Frist von einer Sekunde statt 100 ms, damit Last sie nicht reißt |
| `internal/hexagon/ports/driven/einspielziel.go` | update | Kommentar von `Schliesse` nennt das Warten auf den Schreibvorgang und die Gesamtfrist statt „sofort“ |
| `internal/bootstrap/play_test.go` | unverändert gefahren | `TestRunPlayZweitesSignal` ist der rote Test (a); Änderung nur, wenn der Architect es entscheidet |
| `spec/spezifikation.md` | vom Architect geändert, nicht vom Implementer | Randformen 1 bis 5 stehen seit 36b6387c und f17d9027 (Architect, vor dem Anspruch), die Randform 8 seit fa8a8479 (Architect, auf Rückgabe im Slice, im Diff des Zweigs); weitere Änderung nur über Rückgabe an den Architect |
| `docs/user/abdeckung-unit.md` | erzeugt | `make abdeckung` schreibt die Zeilen der neuen Tests aus ihren Deklarationen; nicht von Hand |
| `.claude/commands/implement-slice.md` | update (Closure, Planner) | Schritt 16: der Satz zu `TestRunPlayZweitesSignal` nennt den Ist-Zustand statt „bis ein Folge-Slice es behebt“ (§1) |

Das Schreiben einer Nachricht läuft über `fe.Send` und `fe.Flush`; der Zustand gilt je `Flush`, nicht je Nachricht (Randform 2, in der Spezifikation entschieden).

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
8. *`Anfrage` nach dem Schließen* (`AGENTS.md` §3.12, *Stelle einer Nachricht*): `Anfrage` ist ein Schreibvorgang wie eine Gruppe; beginnt sie nach dem zweiten Signal, darf sie nicht neben das `Terminate` schreiben (zwei Schreiber auf einem `Frontend`). Der Code liefert `PGR-E4003` ohne zu schreiben (vorher ein Write auf die geschlossene Verbindung, derselbe Code). Die Spezifikation nennt die Stelle nicht (sie nennt nur eingereihte Gruppen). Entschieden: Spezifikation (`LH-FA-20.a` *Abbruchsignal*, Absatz zu eingereihten Gruppen): kein Schreiben, `PGR-E4003`, der Adapter gibt nur den Fehler zurück (Teilsatz „ohne Meldung“ entfallen, V-2); der Code bleibt unverändert, `TestEinspielAnfrageNachSchliessen` (beide Arten, Fehlercode und kein Byte nach dem Schließen) trägt §3.10 (entschieden in fa8a8479).
7. *Mutation* (§3.10): die Sperre bleibt 300 ms nach dem Senden gehalten, in einer Kopie und im Test als Wartezeit des Schreibens nachgebaut (Probe P1 von `slice-test-unter-last-ursache`); der rote Test ist `TestRunPlayZweitesSignal`, ein deterministischer Test im Adapter ergänzt ihn, weil der Wettlauf dort ohne Netz und Zeitgrenze von 30 s fällt.

**Risiken:**

- Der Zustand „teilweise geschrieben“ ist je `Flush` zu grob: Eine Gruppe mit mehreren Nachrichten hält ihn, und `Terminate` käme bei großen Gruppen später als nach der Spezifikation — **Ausgang:** offen bis Closure (Randform 2, Architect).
- Der Fix verlagert den Wettlauf, ohne ihn zu beseitigen (`TestRunPlayZweitesSignal` bleibt unter Last in einzelnen Läufen rot) — **Ausgang:** offen bis Closure; Beleg: mindestens 20 Läufe je Last (DoD).
- Ein Test mit Zeitgrenze gerät selbst unter Last ins Rot (`BEO-REPO/test-faellt-unter-cpu-last-aus`) — **Ausgang:** weiter offen → `BEO-REPO/test-faellt-unter-cpu-last-aus`.
- Der Slice erreicht die Marke von §4 (1500, Metrik `AGENTS.md` §3.14) — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Die Closure-Notiz (Lerneintrag, Risiko-Ausgänge, Register) schreibt der Planner bei Closure, vor dem `git mv` nach `done/`.

### Belege des Implementers

**Umsetzung** (Commit `b0a7e10a`, Basis `dced44df`): `einspielSession` führt `schreibEnde chan struct{}` unter `zustand`; ungleich `nil` heißt „ein Schreibvorgang (ein `Flush` von `Anfrage` oder einer Gruppe) läuft“. `beginneSchreiben` prüft `geschlossen` und setzt den Zustand in einem Schritt, `endeSchreiben` löscht ihn und weckt `schliesse`. `schliesse` setzt `geschlossen` und liest `schreibEnde` im selben Schritt, wartet höchstens bis zur Gesamtfrist `terminateFrist` ab dem Aufruf (Warten und `SetWriteDeadline(frist)` teilen sie) und sendet danach `Terminate`; reißt die Frist, bleibt es aus. Die Sperre `schreiben` hält nur noch `Anfrage` und den Sender auseinander. `Anfrage` nach dem Schließen liefert `PGR-E4003`, ohne zu schreiben (vorher ein Write auf die geschlossene Verbindung, Fehler gleich). Randformen 1 bis 5 stehen in `LH-FA-20.a` *Abbruchsignal* (f17d9027); der Slice selbst ändert die Spezifikation nicht. **Eine Randform war nicht genannt und ist im Code entschieden** (Review F-4): `Anfrage` nach dem Schließen (§6 Punkt 8). Rückgabe an den Architect: *Randform:* Stelle „`Anfrage` beginnt nach dem zweiten Signal“ · *Frage:* soll die Spezifikation sie neben „eingereihte Gruppen werden nicht mehr geschrieben“ entscheiden (kein Schreiben, `PGR-E4003`, ohne Meldung)? Der Code schreibt nicht, weil ein Schreiben neben `Terminate` zwei Schreiber auf einem `Frontend` hieße; die Entscheidung des Architect kann ihn bestätigen oder ändern. Entschieden: angenommen in fa8a8479, der Code bleibt unverändert.

**Mutationen** (jeweils eine Änderung an `einspielen.go`, Lauf `go test ./internal/adapters/driven/postgres/ -run Einspiel` im Image `pgwire-recorder:test` mit dem Arbeitsbaum per Bind-Mount, danach zurückgesetzt; Mutant kommt so ohne Cache an):

| Zusage (Wortlaut `LH-FA-20.a`) | Mutation | roter Test (aus dem Grund der Zusage) |
|---|---|---|
| Terminate bei vollständig geschriebener Anfrage, Bearbeitung nicht abgeschlossen | `endeSchreiben` in `Anfrage` 3 s nach dem `Flush` (die Gesamtfrist der Tests ist 1 s) | `TestEinspielSchliesseSendetNachVollstaendigerAnfrage` (kein Terminate hinter der Anfrage); dazu `TestRunPlayZweitesSignal` mit der Verzögerung von 300 ms 30,3 s rot (Lauf vor der Nacharbeit) („play endet binnen 30 s nach dem zweiten Signal nicht“ bzw. kein `X`), ohne Mutation grün |
| dasselbe für eine Gruppe | `endeSchreiben` im Sender 3 s nach dem `Flush` | `TestEinspielSchliesseSendetNachVollstaendigerGruppe` (daneben `TestEinspielGruppeNachrichten`, `TestEinspielGruppeUnabhaengig`) |
| kein Terminate in einen teilweise geschriebenen Schreibvorgang | `schliesse` wartet nicht (`wartePunktlich(nil, …)`) | `TestEinspielGruppeSchliesseBeimSenden` (Panic des Fakes: zweites `Write` während des ersten, `sendeSperre` und `steuerConn`), `TestEinspielSchliesseSchreibtNichtBeimSenden` |
| Frist begrenzt das Warten | Timer `time.Hour` statt `time.Until(frist)` | `TestEinspielWartePunktlich`, `TestEinspielSchliesseBeimSenden`, `TestEinspielGruppeSchliesseBeimSenden` (je 5 s rot, Lauf nach der Nacharbeit); im ersten Lauf zusätzlich `TestEinspielSchliesseSchreibtNichtBeimSenden/Schliesse`, `…SchliesseBeendetSender`, `…VerwirftWarteschlange` |
| gerissene Frist: kein Terminate | `wartePunktlich` liefert nach Ablauf `true` | `TestEinspielGruppeSchliesseBeimSenden` (Panic: zweites `Write`) |
| Gesamtfrist für Warten und Schreiben zusammen | `SetWriteDeadline(time.Now().Add(s.frist))` nach dem Warten | `TestEinspielSchliesseGesamtfrist` (Frist nicht nach Ende des Vorgangs plus Gesamtfrist) |
| Ort „Ende vor der Frist“ (bei abgelaufener Frist und beendetem Vorgang wird gesendet) | Vorab-Prüfung `select … default` in `wartePunktlich` entfernt | `TestEinspielWartePunktlich` (200 Wiederholungen; ohne die Prüfung wählt das `select` zufällig) |
| Grund `SchliesseOhneTerminate` | `SchliesseOhneTerminate` ruft `schliesse(true)` | `TestEinspielSchliesseOhneTerminate`, `TestEinspielSchliesseNurEinmal/SchliesseOhneTerminate,…` (2 Fälle) |
| Grund „erstes Schließen schon geschehen“ | `erstes := true` | `TestEinspielSchliesseNurEinmal/Schliesse,_Schliesse` |
| eingereihte Gruppen nach dem Signal nicht mehr geschrieben | `beginneSchreiben` prüft `geschlossen` nicht | `TestEinspielGruppeSchliesseVerwirftWarteschlange` (beide Arten) |
| `Anfrage` nach dem Schließen endet mit `PGR-E4003` (Randform 8) | `Anfrage` liefert `CodeUnsupported` statt `CodeConnectionLost` | `TestEinspielAnfrageNachSchliessen` (beide Arten), daneben `TestEinspielNaechsteFehler/Senden_scheitert` |
| nach dem Schließen schreibt `Anfrage` nichts | `Anfrage` setzt „schreibt“ ohne Prüfung von `geschlossen` | `TestEinspielAnfrageNachSchliessen` (beide Arten) |

„Ohne Meldung“ (Randform 8): Der Adapter hat keinen Ausgabeweg und gibt nur den Fehler zurück; ob der Aufrufer ihn meldet, belegt dieser Slice nicht (Grenze: kein Test dieses Teilsatzes, Adresse: Planner bei der Closure). Die Reihenfolge der letzten zwei Gründe (Teilweise geschrieben und Frist) ist laut Spezifikation nicht beobachtbar und hat keinen eigenen Fall. Zwei der Fälle (Warten entfällt, Terminate trotz Frist) werden durch den Panic des Fakes (`close of closed channel` in `einspielen_gruppe_test.go`) rot, nicht durch eine Meldung; Grund ist in beiden Fällen ein zweites gleichzeitiges `Write`. Eine Mutation, die `beginneSchreiben` ignoriert und dann `endeSchreiben` ruft (Panic `close of nil channel` in einem anderen Test), zählt nicht als Beleg für die Zeile *Anfrage nach dem Schließen*; die Zeile trägt die zweite Mutation.

**Last** (`TestRunPlayZweitesSignal` als vorab gebautes Testbinary, je Lauf ein Container `--network none`, Last wie `slice-test-unter-last-ursache` §7 (a): K Container mit `while :; do :; done` im Image `postgres:17-alpine@sha256:b0f9560a…`, 5 s Vorlauf; Host Linux, 20 Kerne, kein anderer Container): K=32: 20 Läufe, 0 rot (Last 1 min vorher 10,7, nachher 19,8); K=48: 20 Läufe, 0 rot (vorher 29,9, nachher 39,1). Zum Vergleich der alte Stand laut `slice-test-unter-last-ursache` §7: 8 von 20 und 9 von 20 rot bei 16 Kernen. Grenze: Die Hosts unterscheiden sich (16 gegen 20 Kerne), die Lastzahl blieb bei K=32 um die Kernzahl; ein Lauf des alten Codes auf diesem Host steht nicht daneben. Das Skript liegt im Scratchpad der Sitzung, nicht im Repo.

**Nacharbeit nach dem Review (F-1 bis F-5, Klasse „Zeitgrenze ohne Spielraum“):** Die Reihenfolge-Tests (`…SendetNachVollstaendigerAnfrage`, `…Gruppe`, `…Gesamtfrist`) fahren mit einer Gesamtfrist von einer Sekunde (`NeueEinspielSessionMitFrist`, `langeFrist`); die Verzögerung der Mutanten ist 3 s. `TestEinspielSchliesseGesamtfrist` wartet nicht mehr auf `time.Sleep(80 ms)` ab dem Aufruf, sondern auf den Beginn des Schließens (eine Gruppe wird abgelehnt) plus 200 ms und lässt `Terminate` erst dann zu; Spielraum nun rund 800 ms statt 20 ms. Die Spezifikation bleibt bei 100 ms (`terminateFrist`), die Zahl der Session ist nur die Frist des Tests. Die Abdeckungszeilen sind enger gefasst (F-2: „sobald der Schreibvorgang endet“, „200 ms nach Beginn des Schließens“). §3 nennt die Dateien des Diffs (F-3). Die Vorab-Prüfung in `wartePunktlich` hat Fall und Mutation (F-5). `endeSchreiben` läuft per `defer` in `schreibeGruppe` (INFO), der Kommentar von `Schliesse` im Port nennt Warten und Gesamtfrist.

**Last der neuen Tests** (Testbinary `go test -c ./internal/adapters/driven/postgres`, Muster `SendetNachVollstaendig|Gesamtfrist|WartePunktlich|AnfrageNachSchliessen|Schliesse`, je Lauf ein Container `--network none`, Last wie oben, 20 Kerne): K=32: 20 Läufe, 0 rot (Last 1 min vorher 10,0, nachher 24,2); K=48: 20 Läufe, 0 rot (vorher 30,4, nachher 44,7). Grenze: 0 rot bei 20 Läufen belegt keinen Spielraum; ihn benennt die Konstruktion (1 s Frist, 200 ms Wartezeit).

**Gates:** `make gates` am Endstand (Commit `b0a7e10a`) und erneut nach der Nacharbeit (Code-Commit `9295e31f`), beide Exit 0; **Endstand nach Review r2** (Commit `0321dd35`, Spezifikation mit `fa8a8479`): `make gates` Exit 0, `workflows-check: kein Befund`; die folgenden Angaben zum ersten Lauf: Exit 0, `run-integration-tests: gruen`, `test-postgres-versionen: grün`, `workflows-check: kein Befund`. `make abdeckung-check` lief vorher rot („nicht aktuell“, die vier neuen Tests), `make abdeckung` schrieb `docs/user/abdeckung-unit.md` neu, danach grün im Gate. Keine Gegenprobe gefahren: kein Gate, keine Gegenprobe und kein Vertrag unter `harness/` berührt.

**Größe (§3.14):** 252 Zeilen nach dem Befehl der Regel (Basis `dced44df`, Pfade `internal test`), Marke 1500 nicht nahe.

**Kommentare (§3.7/§3.11):** neue und geänderte Kommentare in `einspielen.go` (Typ `einspielSession`, `beginneSchreiben`, `endeSchreiben`, `Schliesse`, `schliesse`, `wartePunktlich`) beschreiben den Ist-Zustand; jede Zusage darin hat die Tests der Tabelle. Handbuch und README unverändert (§1). Der Kommentar des Ports `Schliesse` nennt das Warten auf den Schreibvorgang und die Gesamtfrist (statt „sofort“).

**Beobachtungen für die Closure:** Die Prüfung von `TestEinspielSchliesseGesamtfrist` hängt seit der Nacharbeit an einem Spielraum von rund 800 ms (Frist 1 s, Wartezeit 200 ms nach Beginn des Schließens); unter extremer Last könnte sie rot werden (`BEO-REPO/test-faellt-unter-cpu-last-aus`, Risiko 3 in §6).

### Closure (Planner, 2026-10-25)

**Gelesen:** Review-Reports r1 (F-1 MEDIUM, F-2 bis F-5 LOW, F-6 bis F-8 INFO) und r2 (F-1 bis F-3 LOW, F-4, F-5 INFO), Verifikationsbericht (bestanden mit Vorbehalten V-1, V-2), alle unter `docs/reviews/2026-10-25-*`. r1 ist mit der Nacharbeit behoben (r2 (a), (b)). Offen aus r2 war der Nachzug von §3 (F-2, V-1) und die Mutation zum Code von Randform 8 (F-3, in §7 nachgetragen); beides erledigt.

**Zwei Befunde der Closure:**

1. *V-1 (§3 folgte dem Code nicht):* behoben im Commit `9ba623b9`: Die Zeile `einspielen.go` nennt den Zustand `schreibEnde` je `Flush`, die Zeile zu `spec/spezifikation.md` nennt die Änderung durch den Architect (Randform 8), §1 und Kopf nennen fa8a8479, und §3 führt die Zeile `implement-slice.md`. §4 verweist auf Randform 2 („je Nachricht statt je `Flush`“) und stimmt.
2. *V-2 („ohne Meldung“ ohne Test):* **benannte Spec-Lücke, entschieden in @HASH@ (Architect, 2026-10-25): die engere Fassung, kein Folge-Slice; der Satz in `LH-FA-20.a` lautet jetzt „sie endet mit `PGR-E4003`, der Adapter gibt nur den Fehler zurück“.** Der Satz steht in `LH-FA-20.a` *Abbruchsignal* („sie endet mit `PGR-E4003`, ohne Meldung“, fa8a8479). Der Adapter hat keinen Ausgabeweg; die Tests des Slice zeigen nur, dass `Anfrage` nach dem Schließen nichts schreibt und `PGR-E4003` liefert. Ob der Aufrufer (Dienst `play`, Bootstrap) den Fehler dieser Stelle meldet, belegt kein Test (`TestRunPlayZweitesSignal` prüft „kein `level=ERROR`“ nur für den Fall, dass `Terminate` ankommt). *Randform · Frage an den Architect:* Anfrage nach dem zweiten Signal · soll der Teilsatz „ohne Meldung“ auf das verengt werden, was geprüft ist („sie endet mit `PGR-E4003`; der Adapter gibt nur den Fehler zurück“, die Meldung regelt die Zeile *Abbruchsignal*: die unterbrochene Interaktion ist kein Fehler), oder soll ein Test auf Dienst- oder Bootstrap-Ebene den Teilsatz tragen (`play` mit einer Anfrage nach dem zweiten Signal, stderr ohne Meldung, Exit 0)? Empfehlung des Planners, nicht Entscheidung: die engere Fassung, weil sie keinen neuen Test an einer Stelle verlangt, die der Adapter nicht kennt (`AGENTS.md` §3.11). Bis zur Entscheidung steht der Satz ungeprüft in der Spezifikation; ein Folge-Slice entsteht nur, wenn der Architect den Test wählt.

**Lerneintrag (Steering Loop).**

- **Geschärfte Regel:** `implement-slice.md` Schritt 16 nennt für `TestRunPlayZweitesSignal` den Ist-Zustand: Rot ist ein Befund, keine Wiederholung; der Wettlauf im Code ist behoben (Anker `seit slice-play-zweites-signal-terminate-wettlauf`, Herkunft `BEO-REPO/test-faellt-unter-cpu-last-aus`). Der Wortlaut der Zeile ist am 2026-10-24 vom Nutzer für den Kern bestätigt; die Berichtigung des Satzes auf den Ist-Zustand ist eine Änderung, die der Planner auf Anweisung des Koordinators vornahm und die der Nutzer noch nicht bestätigt hat (zu bestätigen oder streichen lassen; die Zusätze (i) und (ii) von `slice-test-unter-last-ursache` Vorbehalt 5 stehen weiter offen).
- **Was gemessen wurde:** Das Rot war ein Wettlauf im Produktcode (`TryLock` als Stellvertreter für „der Sender schreibt“). Die Zwischenregel „einmal wiederholen“ aus `slice-test-unter-last-ursache` hätte ihn verdeckt; sie nimmt `TestRunPlayZweitesSignal` deshalb aus. Nach der Behebung: 0 von 20 und 0 von 20 rot bei K=32 und K=48 (vorher 8 und 9 von 20 auf einem anderen Host), Verifikation 6 von 6 bei K=32.
- **Folge, Empfehlung und nicht Entscheidung des Nutzers:** Ein Test mit Zeitgrenze bekommt den Spielraum durch Konstruktion (Ereignis statt Schlafen, Test-Frist ein Vielfaches der Wartezeit), nicht durch Läufe; „0 rot in 40 Läufen“ belegte hier nichts, das Review r1 F-1 fand 20 ms Spielraum, wo der Plan 80 ms nannte, ohne dass ein Lauf rot war. Der Vorschlag steht im `state.md` von `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test` (dritte Evidence-Datei, Schwelle 3×); die Zuweisung trifft der Nutzer.
- **Benannte Spec-Lücke:** V-2 oben.
- **Kein neuer Sensor.**

**Risiken aus §6 und ihr Ausgang** (je genau einer):

1. *Zustand je `Flush` zu grob, `Terminate` käme bei großen Gruppen später* — **entfallen:** Der Architect hat die Granularität in f17d9027 als je Schreibvorgang entschieden und die Gesamtfrist von 100 ms gesetzt; der Code folgt (`TestEinspielSchliesseSendetNachVollstaendigerGruppe`, Mutation in §7), eine Abweichung von der Spezifikation entstand nicht.
2. *Der Fix verlagert den Wettlauf* — **entfallen:** 0 von 40 Läufen rot unter Last (Implementer), 6 von 6 (Verifikation); die Konstruktion lässt `Terminate` nach dem Ende des Schreibvorgangs ausbleiben nur über die Gesamtfrist, die die Spezifikation zusagt. Das Restfenster (Goroutine länger als 100 ms angehalten) steht unter Risiko 3.
3. *Ein Test mit Zeitgrenze gerät unter Last ins Rot* — **weiter offen** → `BEO-REPO/test-faellt-unter-cpu-last-aus` (Stand im `state.md`: Spielraum der Tests rund 800 ms aus der Konstruktion, `TestRunPlayZweitesSignal` selbst mit der Frist von 100 ms der Spezifikation); kein Auftreten, keine neue Evidence-Datei, Zähler 7.
4. *Die Marke 1500 wird erreicht* — **entfallen:** 341 Zeilen nach dem Befehl von `AGENTS.md` §3.14 (Basis `dced44df`, Pfade `internal test`), gemessen am Stand `c190c61d` vom Verifier und von der Closure (`9ba623b9` ändert keinen gezählten Pfad).

**Beobachtungs-Register** (Schritt 25, Belege `evidence/slice-play-zweites-signal-terminate-wettlauf.md`, Commit `65ed52ec`): `BEO-REPO/plan-folgt-korrektur-nicht` (86), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (99), `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (30), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (66), `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test` (3, **Schwelle erreicht**, Zuordnung der Klasse „Test mit Zeitgrenze ohne Spielraum unter Last“ zu diesem Eintrag ist Urteil des Planners; sie lag dort näher als bei `BEO-REPO/test-faellt-unter-cpu-last-aus`, weil der Test grün war und das Review die feste Zeit fand); `BEO-REPO/test-faellt-unter-cpu-last-aus` nur im `state.md` (Risiko 3, kein Auftreten). Der Zähler wird nicht gesetzt.

**Lese-Schritt und Rate (`AGENTS.md` §3.15).** Seit `slice-harness-regeln-rueckbau` sind elf Slices geschlossen: `slice-v1-abschluss-verbindungsende-upstream-extended`, `slice-version-option`, `slice-play-feldvergleich`, `slice-play-feldvergleich-doku`, `slice-go-toolchain-1-27-3`, `slice-release-v0-2-0`, `slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`, `slice-homebrew-nachweis-macos-arm64`, `slice-test-unter-last-ursache` und dieser (die Closure von `slice-test-unter-last-ursache` nannte neun und ließ `slice-tests-ueberlebende-mutanten` aus; das ändert dort nichts, weil unter zehn). Die zehn zuletzt geschlossenen sind die ab `slice-version-option`. Rate je Klasse, zu der eine Datei dazukam und deren Regel in `AGENTS.md` §3 auf sie zeigt: `plan-folgt-korrektur-nicht` (§3.9) 5 von 10 = **0,5**; `zusage-im-kommentar-weiter-als-pruefung` (§3.11) 5 von 10 = **0,5**; `randform-im-code-entschieden-dann-zurueckgegeben` 3 von 10 = 0,3 (Regel steht in `implement-slice.md`, nicht in §3); `negativtests-fehlen-bei-neuem-vertrag` (§3.10) 1 von 10 = 0,1. **Keine Rate über 0,5, kein Vorschlag zum Streichen oder Ersetzen.** Zwei Klassen stehen genau an der Grenze; ein weiteres Auftreten in einem der nächsten Slices legt dem Nutzer die Frage vor. `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test` hat 3× und keine Regel in §3; sein Vorschlag steht im `state.md`.

**Beleg `make gates`:** am Endstand nach r2 Exit 0 (oben, Code-Stand `0321dd35`); die Commits der Closure ändern nur Plan, Register und eine Zeile von `implement-slice.md`, keine Gates gefahren (Auftrag: Integration macht der Koordinator).

**Drei Paarungen.** Anker: der Satz in `.claude/commands/implement-slice.md` Schritt 16 (`liegt in`, Anker `seit slice-play-zweites-signal-terminate-wettlauf`). Folge-Slice: keiner genannt, keiner nötig; die Frage V-2 geht an den Architect, ein Slice entsteht nur auf dessen Entscheidung. Register: die Pfade oben liegen in `observations/BEO-REPO/` und führen je ein nicht leeres `evidence/`. Die Lage-Prüfung nach dem Move (Modul 6) steht aus.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-24). Treffer in `REPO`: `BEO-REPO/test-faellt-unter-cpu-last-aus` (7×; mit `slice-test-unter-last-ursache` ist die Ursache von `TestRunPlayZweitesSignal` benannt: Wettlauf im Code; dieser Slice behebt ihn); `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (§3.10 gilt für diesen Slice: je Zusage eine Mutation). Eine Schwelle wird mit diesem Slice nicht neu erreicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-24, im Commit, der diesen Plan anlegt):** Liefer-Punkte 3 (Zustand statt Sperre | Tests mit Mutation | Plattform und Abdeckung). Schichten 2: Driven Adapter `postgres` (Einspiel-Session) mit seinen Tests und Spezifikation (Technik-Stratum, Architect; zählt als Dokumentation der Spezifikation, keine Nutzerdoku-Schicht); der Test im Bootstrap (`TestRunPlayZweitesSignal`) wird gefahren und nicht geändert. Planung zählt keine Schicht; Meldungscode-Konstanten entstehen nicht. Der Slice liegt nicht über der Grenze; Marke 1500.
