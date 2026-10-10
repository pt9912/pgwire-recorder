# Review-Report: slice-v1-abschluss-einspielen-extended — 2026-10-10

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3, §6 und §8, gegen die Spezifikation `LH-FA-20.a` (*Gruppen*, *Interaktion*, *Abbruchsignal*, Schritte 3, 6 und 7), gegen [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md), [ADR-0012](../plan/adr/0012-extended-query-gruppen.md), [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) und [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) und gegen die Hard Rules. Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** Diff `b23ff2c..80095aa`: Code und Tests `1d045fa`, Abdeckungstabellen `a9bc23a`, Belege in Plan §7 `80095aa` (+1616 −101). Geändert sind `internal/adapters/driven/postgres/einspielen.go` (Sender und Warteschlange), `internal/hexagon/ports/driven/einspielziel.go` (Operation `Gruppe`), `internal/hexagon/services/play.go` (Zählen, Fortsetzung, erwarteter Fehler), die Tests dazu, neu `internal/adapters/driven/postgres/einspielen_gruppe_test.go`, `internal/hexagon/services/play_extended_test.go` und `test/integration/play_extended_e2e_test.go`. Grundlage vor dem Code ist `b23ff2c` (Architect: *Gruppen*, Zählen je Antwort, Option A) und die früheren Randformen `f7c9c79`, `19f5512`, `e4963ee`, `85caec4`.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-extended.md`: ganz gelesen am Stand `80095aa`, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen.
- `spec/spezifikation.md` `LH-FA-20.a` ganz, `SPEC-038`; [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-18`](../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler).
- Die ADRs oben, Abschnitt *Entscheidung*; [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) für die Nebenläufigkeit je Verbindung.
- `AGENTS.md` §3.3, §3.5, §3.7 und §3.9 bis §3.13.
- Der Nehmer `slice-v1-abschluss-einspielen-extended-doku` (Kopf, §1, DoD, §8), gelesen auf die Sendung des Gebers.
- Unverändert, aber vom Diff vorausgesetzt: `toFrontendMessage` und `toResponse` in `internal/adapters/driven/postgres/upstream.go`, `Validate` in `internal/hexagon/model/extended.go`, `Play` und `session` in `internal/hexagon/services/play.go`.

**Ausgeführte Läufe:**

- **Arbeitsweise der Mutanten.** Je Mutant eine frische Kopie des Arbeitsbaums per `cp -r` ohne `-p` unter dem Scratch-Pfad (`rev-ext-m-<name>`), genau eine Ersetzung, gofmt geprüft (kein Mutant meldete eine Datei), dann `go test -count=1` des betroffenen Pakets im Image `golang:1.27` ohne Netz mit dem Modul-Cache des Images `deps`. Im Repo nichts geändert; Kopien, Container und Netz danach gelöscht.
- **Mutanten, alle rot an dem Test, den §7 nennt.** `Describe` einer Anweisung `n += 2` → `n++` (`TestPlayExtendedWarten`); kein Warten nach `Flush` (`TestPlayExtendedWarten`, `…AndereArt`, `…FehlerBrichtAb`, `…Fortsetzung`, `…Fatal`, `…ErwarteterFehler`, `…Signal`); `fehlerGesehen` wird nicht gesetzt (`…Fortsetzung`, `…ErwarteterFehler`); der Fehler beendet das Zählen nicht, `n = 0` → `n--` (`…Fortsetzung`); `fehlerGesehen` überspringt das Warten nicht (`…Fortsetzung`, `…ErwarteterFehler`); Fehler von `Gruppe` ignoriert (`…SendenLesenScheitert`); `PANIC` nicht mehr abbrechend (`TestPlayFehlerantwort`, nicht ein Test der Extended-Interaktion); Sendefehler mit anderem Code als `PGR-E4003` (`TestEinspielGruppeSendefehler`); Warteschlange LIFO (`TestEinspielGruppeNachrichten`, `…Unabhaengig`); `Null` wird zu leer (`TestEinspielGruppeNachrichten`); `TryLock` → `Lock` in `Schliesse` (fünf Tests, je 5 s Frist); Prüfung von `geschlossen` in `Gruppe` entfällt (`…SchliesseBeimSenden`); der Sender merkt auch nach `Schliesse` einen Fehler (`…SchliesseBeimSenden`); `close(s.ende)` entfällt (`…SchliesseBeendetSender`, `…VerwirftWarteschlange`); der Sender schließt bei einem Sendefehler nicht (`…Sendefehler`, `…SchliesseBeimSenden`). Ein Mutant ohne Wirkung (`Gruppe(g.Client[:len(g.Client)])`) blieb erwartungsgemäß grün und zählt nicht.
- **Integrations-Mutant.** Binary und Testbinary aus einer Kopie gebaut, in einem eigenen internen Docker-Netz gegen das gepinnte PostgreSQL-Image gefahren (`rev-ext-pg`, `rev-ext-net`, danach entfernt). Mutant: kein Warten nach `Flush` (`g.Client[len(g.Client)-1].Type != model.ClientFlush` → `true`). `TestE2EPlayExtended`, `TestE2EPlayExtendedFehler` und `TestE2EPlayExtendedGegendruck` blieben alle grün (0,43 s, 0,70 s, 1,75 s).
- **`-race`.** Nicht im Gate. Im Image `golang:1.27` (Debian, mit `gcc`) mit `CGO_ENABLED=1`, ohne Netz: `go test -race -count=3` über `internal/adapters/driven/postgres` und `internal/hexagon/services` ohne Befund; `-race -count=30 -run 'Gruppe|Einspiel'` über den Adapter ohne Befund und ohne Flackern. Nur als Hinweis; der Integrationstest lief nicht unter `-race`.
- **Gates.** `make a-check`: `gesamt: 0 Befund(e)`. `make docs-check` vor dem Commit dieses Reports: grün. `make gates` und die Läufe auf `80095aa` habe ich nicht wiederholt; das ist Sache des Verifiers.

---

## Findings

Die Nummern dieses Laufs beginnen bei F-580. Vorheriger Report am selben Modul: `slice-v1-abschluss-einspielen-laufsteuerung` (F-564 bis F-569).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-580 | MEDIUM | Der Kommentar von `TestE2EPlayExtendedGegendruck` und die gleichlautende Zeile in `docs/user/abdeckung-e2e.md` (`LH-FA-18`, Boundary) sagen zu, eine Gruppe mit `Flush`, deren Ausgabe verzögert eintrifft, werde „abgewartet, bevor die nächste Gruppe geht“. Der Test prüft nur den Endzustand der Tabelle; der Server arbeitet die Nachrichten ohnehin in Reihenfolge ab, ob `play` wartet oder nicht. *Failure-Szenario:* Der Mutant „kein Warten nach `Flush`“ ließ alle drei Extended-E2E-Tests grün (selbst gefahren); ein Umbau, der das Warten bricht, wird nur von den Unit-Tests mit Fake gefangen, die Tabelle nennt aber den Lauf gegen die reale Instanz als Beleg. Das Muster war in diesem Modul schon in F-564 (MEDIUM) und F-565 (LOW). | `AGENTS.md` §3.11 (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) | `test/integration/play_extended_e2e_test.go` · „wird abgewartet, bevor“; `docs/user/abdeckung-e2e.md` · „wird abgewartet, bevor die nächste Gruppe geht“ | ja (Mutant „kein Warten nach `Flush`“ in `play.go`, `make test-integration`) | Zusage im Kommentar weiter als die Prüfung |
| F-581 | LOW | Die Abdeckungs-Deklaration von `TestPlayExtendedErwarteterFehler` sagt, ein Fehler sei erwartet, „gleich in welcher Gruppe der Server den Fehler sendet“. Der Fake lässt den Server den Fehler in allen vier Fällen in der ersten Gruppe senden (`gruppen: {{fehler("42P01", "a")}, nil, {antwort(model.ResponseReadyForQuery)}}`); nur die Aufzeichnung ändert sich. Der Satz über die Stelle des Serverfehlers hat keinen Test. | `AGENTS.md` §3.11 | `internal/hexagon/services/play_extended_test.go` · „gleich in welcher Gruppe der Server den Fehler sendet“; `docs/user/abdeckung-unit.md` (dieselbe Zeile) | ja (Satz gegen `ziel := &fakeZiel{sessions: …}` gelesen; kein Gate prüft den Wortlaut) | Zusage im Kommentar weiter als die Prüfung |
| F-582 | LOW | `toFrontendMessage` liefert für eine Zielart außer `statement` und `portal` und für einen unbekannten Typ `PGR-E1000`, und `Gruppe` sendet dann nichts von der Gruppe. Das ist eine Randform ohne Eintrag in §6 und ohne Stelle in `LH-FA-20.a`; §7 meldet sie nachträglich als „angefallen“. Die Abdeckungszeile `TestEinspielGruppeNichtAbbildbar` zitiert dafür *Interaktion*, dort steht davon nichts. `Validate` lässt die Nachrichten beim Laden nie durch, der Pfad ist im Zielstand nicht erreichbar, und die Abbildung ist die des Upstream-Adapters, die dort schon `PGR-E1000` lieferte. | `AGENTS.md` §3.12, §3.11 | `internal/adapters/driven/postgres/upstream.go` · „ohne Abbildung“; `docs/user/abdeckung-unit.md` · „liefert Gruppe als Fehler PGR-E1000“ | nein (kein Gate prüft, ob §6 die Randform nennt) | Randform im Code entschieden und nachträglich gemeldet |
| F-583 | INFO | Zu den drei „beim Bauen angefallenen“ Punkten aus §7, gegen §6 gelesen: (1) nicht gedeckt, siehe F-582. (2) `ReadyForQuery` in einer `Flush`-Gruppe zählt wie jede Antwort: gedeckt, §6 *Akzeptierte Negative* („wartet wie auf jede Antwort ohne Frist“). (3) Die Warteschlange ist unbegrenzt: §6 nennt es nicht; sie bleibt durch den Ablauf klein, weil `extended` nach jeder `Flush`-Gruppe liest, ehe es die nächste einreiht, und nur nach einer Fehlerantwort die übrigen Gruppen der einen Interaktion einreiht. Kein Failure-Szenario. Ebenso: §6 sagt, der Sendefehler komme aus `Naechste` „oder der nächsten Operation“, der Port-Kommentar sagt „der nächsten Gruppe“; `Anfrage` nach einem Sendefehler scheitert an der geschlossenen Verbindung ohnehin mit `PGR-E4003`, und `extended` endet immer in einem Lesen. | `AGENTS.md` §3.12 | Plan §7 · „Randformen. Keine neu entschieden.“; `internal/adapters/driven/postgres/einspielen.go` · `warteschlange [][]pgproto3.FrontendMessage` | nein | Hinweis ohne Aktion |
| F-584 | INFO | Nebenläufigkeit des Senders, durch Lesen und `-race` geprüft, ohne Datenrennen und ohne Fund bei Goroutine-Lecks oder Verklemmung: Das Wecken über `weck` (Puffer 1) verliert kein Signal, weil die Gruppe vor dem Signal eingereiht und der Sender die Warteschlange vor dem Warten prüft; `Schliesse` setzt `geschlossen` unter `zustand`, und der Sender prüft es nach dem Erwerb von `schreiben`; ein blockiertes `Flush` beendet `conn.Close`. Eine schmale Lücke bleibt: Hält der Sender `schreiben` noch zwischen der Rückkehr von `Flush` und dem `Unlock`, während `Naechste` die letzte Antwort schon geliefert hat, scheitert das `TryLock` in `Schliesse`, und ein Abbruch endet ohne `Terminate`. Das deckt *soweit möglich* in `LH-FA-20.a` Schritt 6, ein Test hält es nicht. | [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md); `LH-FA-20.a` Schritt 6 | `internal/adapters/driven/postgres/einspielen.go` · „if s.schreiben.TryLock() {“ | nein | Hinweis ohne Aktion |
| F-585 | INFO | Der Gegendruck-Test verklemmt den Mutanten „`Gruppe` sendet selbst“ nur, wenn die Puffer der Verbindung (Senden plus Empfangen) kleiner sind als die 16 MB Parameter und die 32 MB Ausgabe. Auf einem Host mit sehr großen `tcp_rmem`/`tcp_wmem` bliebe der Mutant grün. Der Lauf selbst flackerte nicht (rund 2 s gegen die Frist von 60 s). Außerdem zählen `TestEinspielGruppeSchliesseBeendetSender` und `warteAufKeinenSender` die Goroutinen des Senders im ganzen Prozess und hängen damit an jedem früheren Test des Pakets, der `Schliesse` nicht erreicht; die Meldung nennt das. | `AGENTS.md` §3.10 | `test/integration/play_extended_e2e_test.go` · „parameterMB“; `internal/adapters/driven/postgres/einspielen_gruppe_test.go` · „func senderZahl()“ | nein | Hinweis ohne Aktion |
| F-586 | INFO | Der Kommentar vor `sendeGruppe` bricht nach „… oder weil das“ mitten in der Zeile um und ist die erste Zeile länger als die übrigen; der Satz ist vollständig, der Kommentar von `einspielSession` nennt `senderLaeuft`, `weck` und `ende` nicht unter dem, was `zustand` schützt. Keine Zusage darin ist falsch. | Maintainability | `internal/adapters/driven/postgres/einspielen.go` · „Es liefert falsch, wenn der Sender enden soll“ | nein | Hinweis ohne Aktion |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services/play.go`, Zählen je Antwort | geprüft, ohne Befund. `Describe` einer Anweisung zählt zwei, `Flush` und `Sync` keine, `DataRow`, `NoticeResponse` und `ParameterStatus` zählen nicht, jede andere Antwort zählt eine; `NotificationResponse` verwirft der Adapter. Gegen `LH-FA-20.a` *Gruppen* Zeile für Zeile verglichen; die Mutanten oben sind rot. |
| `play.go`, Fehler, Fortsetzung, Abbruch | geprüft, ohne Befund. Nach `PGR-E4004` ohne Optionen kein weiteres Lesen und keine weitere Gruppe; mit `--continue-on-error` und bei erwartetem Fehler zählt `play` nicht mehr, sendet den Rest ohne Warten und liest bis zum `ReadyForQuery`; `FATAL`/`PANIC` vor der Prüfung auf erwartet; `mitFehlerantwort` liest `Responses` und `Groups[].Server`; das erste Signal wird in `extended` nicht geprüft, das zweite endet über `abbruch.Err()` in `Play` ohne Fehler. |
| `internal/adapters/driven/postgres/einspielen.go` und `ports/driven/einspielziel.go` | geprüft, ohne Befund außer F-584 bis F-586. Option A umgesetzt wie in §6: `Gruppe` kehrt zurück, die Reihenfolge der Aufrufe bleibt, Sendefehler ist `PGR-E4003` aus `Naechste` und der nächsten Gruppe, `Schliesse` beendet den Sender, ein Fehler danach bleibt ohne Folge. Der Port-Kommentar sagt nur zu, was `TestEinspielGruppe*` hält. |
| `toFrontendMessage` | geprüft, ohne Befund außer F-582. `Null` ist SQL-NULL, ein leerer Wert ein leerer Wert (Bytes `nil` und leer), Formate, Namen, `MaxRows` und Zielart wie aufgezeichnet; die Mutanten sind rot. |
| Hexagon, `make a-check` | geprüft, ohne Befund. `gesamt: 0 Befund(e)`; Port und Service importieren nur `model` und `driven`, `pgproto3` bleibt im Adapter. Die 14 Integrationsdateien liegen in keiner Schicht (bekannter Hinweis, kein Befund). |
| Suppression (`//nolint`) und `.golangci.yml` | geprüft, ohne Befund. Der Diff enthält keine `nolint`-Direktive und ändert das Profil nicht. |
| `SPEC-038` *Warten in Tests* | geprüft, ohne Befund. Die Warteschleifen und Kanalwartestellen im Adapter-Test tragen 5 s als Literal und nennen das ausgebliebene Ereignis; die E2E-Läufe warten über `starteBisFrist` mit 30 s und 60 s als Literal; die Service-Tests rufen synchron auf. |
| Plan §1, §3, §6 gegen den Diff (`AGENTS.md` §3.9) | geprüft, ohne Befund. §3 nennt alle geänderten Dateien einschließlich `TestPlayStart`, `TestRunPlayStartfehler` und `TestE2EPlayZwischenstand`; §1 *Ausdrücklich NICHT* (kein Code in CLI, Bootstrap, PGWire-Adapter) stimmt, im Bootstrap ändert sich nur ein Test. Der Kopf führt die Kennungen, die §1 nennt. |
| `AGENTS.md` §3.13, Nehmer `slice-v1-abschluss-einspielen-extended-doku` | geprüft, ohne Befund. Der Nehmer nennt die Sendung mit der Kennung des Gebers in §1 und DoD, liegt in `next/`, und §8 zählt zwei Liefer-Punkte und eine Schicht; der Geber zählt drei Liefer-Punkte und zwei Schichten (Play-Service mit Port, Upstream-Adapter). Der Diff trägt keine neue Sendung ein. |
| Hard Rules 3.3, 3.5, 3.7; Commit-Messages | geprüft, ohne Befund. Kein Move, keine ADR geändert, keine Spec-Datei im Diff. Die Kommentare im Code sind Zusagen oder Kopplungen. Die Messages nennen `slice-v1-abschluss-einspielen-extended` und [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung). |
| §3.10, Mutationstabelle in §7 | geprüft, ohne Befund außer F-580 und F-581. Die gefahrenen Stichproben entsprechen den Zeilen der Tabelle; kein Mutant der Stichprobe blieb grün. |
| Größe | geprüft, ohne Befund. Rund 1200 der 1600 eingefügten Zeilen sind Tests; die Rückführungen aus §4 sind aus Sicht des Reviews nicht eingetreten. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 4 |

- F-580: Der E2E-Gegendruck-Test hält das Warten nach `Flush` nicht, Kommentar und Tabelle sagen es zu.
- F-581: Abdeckung von `TestPlayExtendedErwarteterFehler` nennt die Stelle des Serverfehlers, der Test variiert sie nicht.
- F-582: Unabbildbare Client-Nachricht ist `PGR-E1000`, nicht in §6, Zitat der Spezifikation trägt es nicht.
- F-583: §7-Punkte gegen §6 gelesen; Warteschlange durch den Ablauf klein.
- F-584: Sender ohne Datenrennen und ohne Leck; schmale Lücke bei `Terminate`.
- F-585: Gegendruck hängt an der Puffergröße des Hosts; Goroutinen-Zählung koppelt Tests.
- F-586: Kommentar vor `sendeGruppe`.

Wiederkehrende Klassen: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-580, F-581; zuvor F-564, F-565), `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (F-582).

**Finding-Klassen dieses Laufs:** Zusage im Kommentar weiter als die Prüfung · Randform im Code entschieden und nachträglich gemeldet

## Verdikt

**Merge-blockierend:** ja, bis F-580 bearbeitet ist. Das gelieferte Verhalten ist an den Prüfpunkten des Auftrags richtig: Zählen je Antwort, Gegendruck und Unabhängigkeit von Senden und Lesen, Schließen mitten im Senden, Sendefehler, Abbruch ohne weitere Gruppe, erwarteter Fehler in jeder Gruppe einschließlich `Groups[].Server`, Null gegen leer und die Ablösung des Zwischenstands; die Stichproben-Mutanten sind rot. Offen ist die Weite einer Zusage im E2E-Test und in der Tabelle (F-580), die den Lauf gegen die Instanz als Beleg für ein Warten nennt, das dieser Lauf nicht hält. F-581 und F-582 blockieren nicht; F-583 bis F-586 sind Hinweise.

**Übergabe:** F-580 bis F-582 gehen an den Implementer. Widerspricht er F-580, läuft der Konflikt über den Architect. F-582 geht zusätzlich an den Architect zur Kenntnis, ob die Randform in `LH-FA-20.a` oder als akzeptiertes Negativ in §6 stehen soll. F-583 bis F-586 gehen an den Verifier zur Kenntnis; F-584 und F-585 sind Hinweise für eine künftige Welle, kein Auftrag an diesen Slice. Die Finding-Klassen gehen in §7 der Closure. Dieser Report ersetzt keine Verifikation.
