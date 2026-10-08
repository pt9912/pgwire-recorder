# Review-Report: slice-v1-abschluss-herunterfahren — 2026-10-08

**Review-Art:** Code-Review gegen Plan §1, §3 und §6, die Spezifikation (`LH-FA-13.a` mit *Randfälle des Herunterfahrens* und *Weitere Signale*, `LH-FA-13.b`, `LH-FA-14.a` `sessions`, `LH-FA-17.a` *Dauer*, `SPEC-051`, `LH-FA-18.a` *Abbruch*, `SPEC-038`), [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md), [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) und die Hard Rules (Modul 10). Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff d1bb8e4..c644fc3`, fünf Commits:

- `17e4925`: `record` mit Frist (Composition Root, CLI, PGWire-Adapter, Kern `CloseSession` mit `EndForced`, Tests in allen Schichten, E2E).
- `64f48b3`: `replay` mit Frist, neue Port-Methode `Forced` am `Replayer`.
- `61ad648`: Exit-Code je Klasse, Benutzerhandbuch.
- `5e20a87`: Integrationstest ohne `pgproto3` nach dem Befund von `make a-check`.
- `c644fc3`: Plan §7 *Belege des Implementers*.

28 Dateien, 2437 Zeilen hinzu und 136 entfernt; davon etwa 550 Zeilen Produkt-Code, der Rest Tests, Handbuch, Abdeckungstabellen und Plan.

**Skill:** `.harness/skills/reviewer.md` am Stand `c644fc3`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-08

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-v1-abschluss-herunterfahren.md`: ganz gelesen am Stand `c644fc3`, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen.
- `spec/lastenheft.md` `LH-FA-13`; `spec/spezifikation.md` `LH-FA-13.a`, `LH-FA-13.b`, `LH-FA-14.a`, `LH-FA-17.a` (*Dauer* und Optionstabelle), `LH-FA-03.b` (*Zeitpunkte*, *Fehlerebene*), `SPEC-013` bis `SPEC-019`, `SPEC-046`, `SPEC-051`, `SPEC-038` *Warten in Tests*.
- [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) (Entscheidung E: Adapter meldet Ereignisse, Service entscheidet), `.a-check.yml`.
- `AGENTS.md` §3.3, §3.7, §3.9 bis §3.12.
- Bestand außerhalb des Diffs, soweit der Diff ihn ruft: `internal/adapters/driven/postgres/upstream.go` (`Open`, `Close`), `internal/hexagon/services/record.go` (`Query`, `AwaitServer`, `CloseSession`), `internal/hexagon/model/fehler.go` (`exitCode`).
- Register: `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/spec-randform-erst-im-review-entschieden` (Namen aus Plan §8 und `AGENTS.md`).
- Vorheriger Report: `slice-harness-integration-wait` (bis F-484). Die Nummern dieses Laufs beginnen bei F-485.

**Ausgeführte Läufe:**

- `make test`, `make a-check` (0 Befunde), `make lint` (Exit 0) am Stand `c644fc3`.
- **Mutanten**, je Mutant eine frische Kopie per `cp -r` ohne `-p` unter dem Scratch-Pfad der Sitzung, genau eine Ersetzung per Skript mit geprüfter Trefferzahl 1, eigenes Image-Tag der Stufe `source`, `go test` in einem Container mit `--network none`; danach Image und Kopie entfernt.
  - Stichproben aus §7, alle rot wie angegeben: R13 (`TestRecordZwangsendeBeimSendenDesAufbaus`), P13 (`TestReplayZwangsendeWaehrendShutdown`), P15 (`TestReplayZwangsendeClientLiestNicht`), R20 (`TestRunRecordZweitesSignalOhneFrist`, `…LangeFrist`).
  - Eigene: M1 `return "nicht geschrieben"` → Text des Falls `n == 0`: **grün** (F-487). M2 Zweig `else { r.s.note(err) }` in `richtungen.zwangsende` entfernt: **grün** über alle Pakete unter `internal/` (F-486). M11 `dauerForm` ohne führende Nullen: **grün** (F-490). M29 Schreibfrist der Replay-Sitzung beim Zwangsende `3 * meldeFrist`: rot (`TestReplayZwangsendeClientLiestNicht`). M30 dasselbe im Record: rot (`TestRecordZwangsendeClientLiestNicht`, `TestZwangsendeNebeneinander`). M-hang `handle` kehrt nach der Startnachricht zurück, `-timeout 20s`, nur `TestRecordZwangsendeImAufbau`: `panic: test timed out after 20s`, Goroutine in `server_zwangsende_test.go:229` (F-485).
- **Integrations-Mutant** R38 (Zustellung in `richtungen.zwangsende` durch `note` ersetzt), Stufe `integration` der Kopie mit eigenem Tag, eigenem `--internal`-Netz, eigenem PostgreSQL-Container (Image und Digest aus `harness/mk/integration.mk`) und eigenem Volume, `-test.run '^TestE2ERecordFristLaeuftAb$'`: rot, `Client erhält "" (<nil>) statt der Fehlerantwort PGR-E4006`. Danach sind keine Container, Netze, Volumes oder Images mit den Lauf-Namen übrig.
- **Probe Startphase** (F-489) in einer Kopie: Ein Testfall ruft `bootstrap.Run` 400-mal mit beendetem Kontext, während ein Client in einer engen Schleife auf den Port verbindet; gezählt wurde, ob eine Verbindung angenommen wurde (`sessions=1` oder Debug-Zeile *ohne Startnachricht*). Ergebnis 0 von 400.
- `make docs-check` vor dem Commit.
- Nicht gefahren: `make gates`, `make test-integration` über den ganzen Satz, die übrigen 64 Mutanten aus §7, Lauf mit `-race` (im Image nicht möglich).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-485 | MEDIUM | Fünf neue Tests warten ohne eigene Frist auf ein Ereignis, das der Prüfling herbeiführt: Der Fake schließt den Kanal erst, wenn der Adapter bzw. der Use Case ihn ruft. *Failure-Szenario:* Eine Regression ruft `OpenSession` nicht mehr (Mutant M-hang); `TestRecordZwangsendeImAufbau` hängt dann in Zeile 229 bis zur Zeitgrenze von `go test` (`panic: test timed out after 20s`), statt mit einer Meldung rot zu werden. Dasselbe gilt für `TestRecordZwangsendeNachBemerktemEnde` (146), `TestRecordZwangsendeBrichtAufbauAb` (365), `TestReplayZwangsendeWaehrendShutdown` (579) und `TestRecordZwangsendeEinfacheAnfrage` (105). §6 sagt für dieses Risiko zu, jeder neue Test warte mit eigener Frist. | `SPEC-038` *Warten in Tests*; Plan §6 Risiko *Test mit Signal hängt*; `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` | `internal/adapters/driving/pgwire/server_zwangsende_test.go` · `<-rec.openLaeuft` (Z. 229, 365); · `<-rec.closeLaeuft`; · `<-rep.shutdownLaeuft`; `internal/hexagon/services/record_zwangsende_test.go` · `<-up.letzte.queryLaeuft` | ja (Mutant M-hang oben, `go test -timeout`) | Test wartet ohne eigene Frist auf den Prüfling |
| F-486 | MEDIUM | Der Kommentar an `richtungen.zwangsende` sagt zu: „einen anderen Fehler merkt es nur“. Der Zweig ist ungeprüft; Mutant M2 (Zweig entfernt) bleibt in allen Paketen grün, §7 führt keinen Mutanten dazu. *Failure-Szenario:* Beim Zwangsende einer Session ohne laufende Interaktion (etwa *abgeschlossen, nicht zugestellt*) scheitert das Schreiben der Aufzeichnung in `CloseSession` vorübergehend. Ohne den Zweig steht der Fehler nicht im Log und wird nicht gemerkt; gelingt `Finish` danach, endet der Lauf mit 0 statt 3 (`LH-FA-13.b`). | `AGENTS.md` §3.10, §3.11; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `internal/adapters/driving/pgwire/server.go` · „einen anderen Fehler merkt“; · `r.s.note(err)` (Z. 593) | ja (Mutant M2, `make test`) | Zusage im Kommentar ohne Mutation |
| F-487 | LOW | `CloseSession` mit `EndForced` meldet für eine nicht übernehmbare Session mit abgeschlossenen Interaktionen „nicht geschrieben“ ohne Grund. Das ist eine dritte Form neben den beiden, die `LH-FA-13.a` *Meldung* nennt (Kennung, oder „ohne abgeschlossene Interaktion nicht geschrieben“), und §6 führt sie nicht. Erreichbar ist der Zweig nur, wenn das Zwangsende zwischen dem Merken von `PGR-E6001` und dem Ende der Session durch den Adapter eintrifft. Kein Test hält ihn (Mutant M1 grün). | `AGENTS.md` §3.12, §3.10 | `internal/hexagon/services/record.go` · `return "nicht geschrieben"` | ja (Mutant M1, `make test`) | Randform vom Code entschieden, die §6 nicht nennt |
| F-488 | LOW | Das Handbuch sagt an zwei Stellen mehr zu, als geprüft ist. (a) Der Hinweis unter *Eine Anwendung aufzeichnen* sagt ohne Bedingung, dass beim Ablauf der Frist „Ihre Anwendung … `PGR-E4006`“ erhält und das Werkzeug mit Exit-Code 4 endet. Eine Verbindung im Aufbau endet nach `LH-FA-13.a` ohne Meldung, und `TestRunRecordFristLaeuftAb` prüft für diesen Fall Exit-Code 0. Der Abschnitt *Herunterfahren mit Frist* nennt die Bedingung. (b) Unter *Wiedergeben* steht, `replay` liefere nach `docker stop` „einen Exit-Code der Tabelle … statt 137“. Kein Test fährt `replay` in einem Container; die Aussage folgt aus der Rechnung 5 s + `SPEC-051` < 10 s. Das Container-Image liefert erst `slice-v1-abschluss-container`. | `AGENTS.md` §3.11; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `docs/user/benutzerhandbuch.md` · „noch laufende Verbindung zwangsweise. Die abgeschlossenen Anfragen einer“ (Z. 179); · „liefert statt 137“ (Z. 264) | nein (Lesen) | Zusage im Handbuch weiter als Prüfung |
| F-489 | INFO | Zur Einordnung von G01 in §7 („über die Schnittstelle nicht erreichbar“). `betreiben` startet `Serve` in einer Goroutine und schließt den Listener danach. Ist der Kontext schon beendet (Signal in der Startphase), kann `Accept` eine Verbindung, die zwischen `Listen` und `Close` im Backlog liegt, noch annehmen. `LH-FA-13.a` *Startphase* sagt: „Sonst nimmt der Prozess keine Verbindung an“. Meine Probe fand das Fenster in 400 Läufen nicht (0 von 400); erreichbar ist es über die Schnittstelle, beobachtet ist es nicht. Für Architect und Verifier. | `LH-FA-13.a` *Startphase*; `AGENTS.md` §3.11 | `internal/bootstrap/bootstrap.go` · `server.Serve(ctx, l)`; `docs/plan/planning/in-progress/slice-v1-abschluss-herunterfahren.md` · „über die Schnittstelle nicht“ | nein (Probe 0 von 400, kein `-race`) | Einordnung eines grünen Mutanten stärker als belegt |
| F-490 | INFO | Rückgabe an den Architect, *Lesart von `00`/`05s`*. Der Code liest die Ausnahme „ohne Einheit außer `0`“ wörtlich (`00` abgelehnt) und die Zahl numerisch (`05s`, `00s` angenommen). Mutant M11 (führende Nullen abgelehnt) bleibt grün: Keine der beiden Lesarten ist durch einen Test gehalten. Wie der Architect auch entscheidet, die Zusage braucht danach einen Test (§3.10). | `LH-FA-17.a` *Dauer*; `AGENTS.md` §3.10, §3.12 | `internal/adapters/driving/cli/cli.go` · `var dauerForm = regexp.MustCompile(` | ja (Mutant M11) | — (Rückgabe an den Architect) |
| F-491 | INFO | Rückgabe an den Architect, *Upstream-Verbindung im Aufbau*. `oeffne` bricht den Kontext des Aufbaus ab. `postgres.Upstream.Open` beachtet ihn nur in `DialContext`; den Handshake danach liest es ohne Kontext und ohne Frist. Hängt der Server dort, bleibt die Goroutine aus `oeffne` mit der Verbindung zum Server offen, bis `Open` zurückkehrt oder der Prozess endet. `LH-FA-13.a` sagt: „die Verbindungen zu Client und Server werden geschlossen“. Im Binary schließt das Prozessende sie gleich nach dem Schreiben der Aufzeichnung. Das Schließen beim Zwangsende selbst verlangt eine Änderung im Driven-Adapter, die §1 ausschließt und die §4 als Rückführungs-Bedingung nennt. | `LH-FA-13.a` *Zwangsende*; Plan §1, §4 | `internal/adapters/driving/pgwire/server.go` · `octx, abbrechen := context.WithCancel(context.WithoutCancel(ctx))`; `internal/adapters/driven/postgres/upstream.go` · `conn, err := u.Dialer.DialContext(ctx, "tcp", u.Address)` | nein (Lesen) | — (Rückgabe an den Architect) |
| F-492 | INFO | Rückgabe an den Architect, *Mutant X03*. Gelesen: `model.exitCode` liefert für jede Ziffer außerhalb des Bereichs den Rückfall 1. Nimmt der Mutant die Ziffer 1 aus dem Bereich, liefert der Rückfall für `PGR-E1…` ebenfalls 1. Kein Test kann den Mutanten unterscheiden; er ist äquivalent. Die Funktion ist Bestand und im Diff unverändert. | `SPEC-015`; Plan §7 | `internal/hexagon/model/fehler.go` · `code[5] >= '1' && code[5] <= '6'` | ja (Lesen) | — (Rückgabe an den Architect) |
| F-493 | INFO | Negativbefund zum Hexagon, für den Verifier. `PGR-E4006` entsteht im Use Case: Record in `CloseSession` aus `EndForced` und `laeuft()`, Replay in `Forced` aus `ungesendet` und `mitten()`. Der Adapter meldet nur das Ereignis. Er liest den Code des Ergebnisses einmal (`== model.CodeShutdownTimeout`), und zwar nur für die Wahl zwischen Zustellen und Merken, nicht für die Einstufung. `Forced` ändert nach Kommentar und Code keinen Zustand; `CloseConnection` folgt in `replaySitzung` danach. Record und Replay melden das Ereignis über verschiedene Port-Formen (`SessionEnd` gegen eigene Methode); das folgt dem Bestand der beiden Ports. Keine Imports gegen `.a-check.yml` (`make a-check` 0 Befunde). | [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md), [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) | `internal/hexagon/ports/driving/replay.go` · `Forced(ctx context.Context, id model.SessionID) error`; `internal/adapters/driving/pgwire/server.go` · `if model.Meldungen(err)[0].Code == model.CodeShutdownTimeout {` | ja (`make a-check`) | — (Negativbefund) |
| F-494 | INFO | Negativbefund zur Nebenläufigkeit, gelesen ohne `-race`. **Signale:** `signal.Notify` steht genau einmal in `main`; die Goroutine liest das erste Signal, beendet `ctx`, liest das zweite und schließt `ablauf`. Danach liest niemand mehr; das dritte Signal bleibt ungelesen im Puffer (Größe 1), jedes weitere verwirft die Laufzeit am vollen Kanal, keines wirkt. **Zwangsende:** `zmu` ordnet `zwang` gegen `anmelden` und `umstellen`, ein spät angemeldetes Ziel führt seine Wirkung selbst aus. Eine Wirkung kann doppelt laufen (Schnappschuss und `umstellen`); im Record fängt `einmal` das, im Replay sind die Fristen idempotent. Kein Lock-Zyklus: `schreibe` gibt `schreiben` frei, bevor es `beende` ruft, und `CloseSession` hält `l.mu` nicht über `upstream.Close`. **Goroutinen:** Einzig die Goroutine aus `oeffne` überlebt `Serve`; sie läuft höchstens bis `Open` zurückkehrt und ist über `s.mu` gegen `Finish` geordnet (zu ihrer Verbindung F-491). | `LH-FA-13.a` *Zweites Signal*, *Weitere Signale* | `cmd/pgwire-recorder/main.go` · `signale := make(chan os.Signal, 1)`; `internal/adapters/driving/pgwire/server.go` · `func (s *Server) anmelden(ende func()) *zwangsziel {` | nein (kein Lauf mit `-race`) | — (Negativbefund) |
| F-495 | INFO | Zur Größe. Der Diff war in einer Sitzung prüfbar, liegt aber an der Grenze: drei Liefer-Punkte, zwei Schichten, rund 550 Zeilen Produkt-Code in sieben Dateien, dazu rund 1500 Zeilen Tests. Die Rückführungs-Bedingung aus §4 ist nicht eingetreten. Wäre sie eingetreten, gäbe §4 den Schnitt schon vor: DoD-Punkt 2 (`replay`, `Forced`, P-Mutanten) als eigener Slice. | Plan §4; `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice | `docs/plan/planning/in-progress/slice-v1-abschluss-herunterfahren.md` · „dann geht die“ | nein | — (Negativbefund) |

---

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `cmd/pgwire-recorder` | geprüft, ohne Befund. Signalbehandlung siehe F-494. |
| `internal/bootstrap` | geprüft. `betreiben`: Frist ab `ctx.Done()`, `0` ohne Zeitgeber, `ablauf` auch bei `0`, ohne offene Verbindung sofort zurück; Startphase siehe F-489. `exitCode` und `TestExitCodeJeKlasse` decken `SPEC-013` bis `SPEC-019`. Stichprobe R20 rot. |
| `internal/adapters/driving/cli` | geprüft. Form der Dauer wie `LH-FA-17.a` *Dauer*, Überlaufgrenze exakt (Literal `9223372036854ms` angenommen, `…855ms` abgelehnt), Umgebung vor Standard, Option vor Umgebung, ungültige Umgebung auch neben der Option `PGR-E2001`. Lesart führender Nullen siehe F-490. |
| `internal/adapters/driving/pgwire` | geprüft. Befund F-486; Hexagon F-493, Nebenläufigkeit F-494, Aufbau F-491. `meldeFrist` = 1 s als Literal geprüft (R17) und an beiden Zustellstellen gehalten (M29, M30 rot). Stichproben R13, P13, P15 rot. |
| `internal/hexagon/model`, `ports/driving`, `services` | geprüft. Befund F-487. Core-Reinheit: keine neuen Imports. Record: verworfene Nummer = abgeschlossene + 1, abgeschlossen-nicht-zugestellt bleibt. Replay: `Forced` ohne Zustandsänderung. |
| Tests (`*_zwangsende_test.go`, `frist_test.go`, `exitcode_test.go`, `test/integration/herunterfahren_e2e_test.go`) | geprüft. Befund F-485. Die Integrationstests lesen `stderr` erst nach `warteEnde` und warten vor dem zweiten Signal über `warteAbgelehnt` mit eigener Frist (§6 *`stderr` zur Laufzeit lesen*). R38 im Integrationslauf rot. |
| `docs/user/benutzerhandbuch.md` | geprüft. Befund F-488. Optionstabelle, Log-Stufen, `sessions` als Vertrag, Exit-Code 3 mit Vorrang und die Form der Dauer stimmen mit der Spezifikation und den Tests überein. |
| `docs/user/abdeckung-*.md` | geprüft, ohne Befund. Die Zeilen folgen den Deklarationen der neuen Tests; `make abdeckung-check` meldet §7 als grün (nicht selbst gefahren). |
| Plan `docs/plan/planning/in-progress/slice-v1-abschluss-herunterfahren.md` | geprüft. §6 ist in keinem Code-Commit geändert (nur `c644fc3`, §7); MEDIUM *Randform im Code-Commit* ohne Befund. §3 folgt dem Code (§3.9). Zur Einordnung von G01 siehe F-489. |
| Hard Rule 3.3 | geprüft, ohne Befund. Im Diff gibt es keinen Move. |
| Hard Rule 3.7 (Kommentar-Klassen) | geprüft, ohne Befund. Die neuen Kommentare beschreiben den Ist-Zustand; der Kommentar an `replaySitzung`, der vorher die fehlende Frist nannte, ist ersetzt. Zusagen ohne Prüfung siehe F-486. |
| Hard Rule 3.13 / MEDIUM *Adresse nimmt nicht an* | geprüft, ohne Befund. Der Diff nennt keinen Slice neu als Adresse; `slice-v1-abschluss-container` in F-488 nennt dieser Report, nicht der Diff. |
| Commit-Messages | geprüft, ohne Befund. Alle fünf nennen `slice-v1-abschluss-herunterfahren` und `LH-FA-13`, keine Struktur-ID. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 7 |

**Summary:** 0 HIGH · 2 MEDIUM · 2 LOW · 7 INFO (F-485 fünf neue Tests warten ohne eigene Frist auf den Prüfling, Hänger bis zur Zeitgrenze nachgewiesen; F-486 Zweig „anderen Fehler nur merken“ beim Zwangsende ohne Mutation, Mutant grün; F-487 dritte Meldungsform „nicht geschrieben“ nicht in §6, ungeprüft; F-488 Handbuch sagt Zustellung und Exit-Code 4 ohne Bedingung und den Container-Exit-Code ohne Test zu; F-489 bis F-495 Startphasen-Fenster, drei Rückgaben an den Architect, Hexagon, Nebenläufigkeit, Größe). Wiederkehrende Klassen: `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` (F-485), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-486, F-488), `BEO-REPO/spec-randform-erst-im-review-entschieden` (F-487).

**Finding-Klassen dieses Laufs:** Test wartet ohne eigene Frist auf den Prüfling · Zusage im Kommentar ohne Mutation · Randform vom Code entschieden, die §6 nicht nennt · Zusage im Handbuch weiter als Prüfung · Einordnung eines grünen Mutanten stärker als belegt

## Verdikt

**Merge-blockierend:** nein. F-485 und F-486 brauchen vor der Closure einen Ausgang. Die Einstufung von `PGR-E4006` liegt im Use Case, die Port-Methode `Forced` ist sauber, und die Signalbehandlung hat weder doppelte Handler noch einen Lock-Zyklus. Die Stichproben aus §7 sind so rot wie angegeben. Offen sind ein Testgeschirr, das bei bestimmten Regressionen hängt statt zu melden, gerade das Risiko, das §6 ausschließen wollte, und eine Kommentar-Zusage ohne Mutation.

**Übergabe:** F-485, F-486 und F-488 gehen an den Implementer. F-487 geht nach `AGENTS.md` §3.12 an den Architect (Randform), danach an den Implementer. F-489 geht an Architect und Verifier. F-490 bis F-492 sind Beiträge zu den laufenden Entscheidungen des Architect, F-493 bis F-495 gehen an den Verifier. Dieser Report ersetzt keine Verifikation.
