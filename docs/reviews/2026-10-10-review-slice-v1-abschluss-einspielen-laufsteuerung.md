# Review-Report: slice-v1-abschluss-einspielen-laufsteuerung — 2026-10-10

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3, §6 und §8, gegen die Spezifikation `LH-FA-20.a` (Tabelle *Fehlerregeln*, Schritte 6 und 7, *Randformen je Schritt*) und `LH-FA-17.a`, gegen [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md), [ADR-0014](../plan/adr/0014-konfigurationsdatei.md) und [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) und gegen die Hard Rules. Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** Code-Commit `1b63ef8` auf `85caec4` (+881 −54 mit Plan §3 und Abdeckungstabellen), dazu die Belege `de00222` (Plan §7). Geändert sind `internal/adapters/driving/cli/cli.go`, `internal/hexagon/services/play.go`, der Kommentar von `internal/hexagon/ports/driving/play.go`, die Tests des Lesers und des Play-Service, neu `internal/hexagon/services/play_laufsteuerung_test.go` und `test/integration/play_laufsteuerung_e2e_test.go`. Grundlage vor dem Code ist `85caec4` (Randformen in `LH-FA-20.a` *Meldungen*, *Abbruchsignal*, *Exit-Code*); danach `ee23fad` (Architect: akzeptiertes Negativ (e), Schicht-Abgrenzung des Port-Kommentars, `ARC-003` im Kopf). Beide gelesen und gegen den Code gehalten.

**Skill:** `.harness/skills/reviewer.md` @ `6db272b`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-v1-abschluss-einspielen-laufsteuerung.md`: ganz gelesen am Stand `ee23fad`, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen.
- `spec/spezifikation.md` `LH-FA-20.a` ganz, `SPEC-038`; [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-FA-14`](../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler).
- [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) ganz; die übrigen ADRs oben, Abschnitt *Entscheidung*.
- `AGENTS.md` §3.3, §3.5, §3.7 und §3.9 bis §3.13, dort auch *Nachzählen beim Eintragen*.
- Die Adressen des Plans, §1 und DoD: `slice-v1-abschluss-einspielen-extended`, `slice-v1-abschluss-antwortvergleich`; gelesen auf Erwähnung: `slice-v1-abschluss-einspielen-anmeldung`, `slice-v1-abschluss-einspielen-tls`.
- Unverändert, aber vom Diff vorausgesetzt: `internal/bootstrap/bootstrap.go` (`play`), `model.Meldungen`, Port `Einspielziel`, `internal/adapters/driven/postgres/einspielen.go` (`Verbinde`, `aufbau`).
- Vorheriger Report am selben Modul: `slice-v1-abschluss-einspielen` (F-547 bis F-554). Die Nummern dieses Laufs beginnen bei F-564.

**Ausgeführte Läufe:**

- **Arbeitsweise der Mutanten.** Je Mutant eine frische Kopie des Arbeitsbaums per `cp -r` ohne `-p` unter dem Scratch-Pfad (`rev-lauf-<m>`), genau eine Ersetzung in `internal/hexagon/services/play.go`, gofmt-sauber (die Stufe `test` meldet keine Datei), dann `make test` und bei den grünen zusätzlich `make test-integration` in der Kopie; die Kopien danach gelöscht. Im Repo nichts geändert. Die Integrationsläufe räumen Container und Netz selbst ab; `docker ps -a` und `docker network ls` mit `pgr` danach leer.
- **Mutanten.**
  - M1: in `Play` `if abbruch.Err() != nil && !s.optionen.FinishSessionOnInterrupt {` statt `if abbruch.Err() != nil {` — `make test` grün, `make test-integration` grün.
  - M19: in `session` das Schließen bei `abbruch` (`context.AfterFunc(abbruch, us.Schliesse)`) nur ohne `FinishSessionOnInterrupt` — `make test` grün, `make test-integration` grün.
  - MA (Zeile der Tabelle in §7): `errors.Join(append(frueher, err)...)` — rot, `TestPlayFortsetzungAbbruch` alle vier Fälle.
  - MB (Zeile der Tabelle in §7): `erwartet := s.optionen.AllowRecordedErrors` — rot, `TestPlayErwarteterFehler/ohne_aufgezeichnete_error_response` und `/mit_--continue-on-error`.
- **Hilfe am Binary.** `make build`, dann `pgwire-recorder:dev play --help` netzlos: die drei Optionen mit Umgebungsvariablen wie im Code. `PGWIRE_RECORDER_FINISH_SESSION_ON_INTERRUPT=yes`: `PGR-E2001` mit dem Namen der Variable.
- **Gates.** `make a-check`: `gesamt: 0 Befund(e)`. `make test` am Stand `ee23fad`: grün. `make abdeckung-check`: grün. `make docs-check` vor dem Commit dieses Reports: grün.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-564 | MEDIUM | Dass das zweite Signal auch mit `--finish-session-on-interrupt` die Verbindung sofort schließt und die Unterbrechung kein Fehler ist, hält keine Mutation: M1 und M19 bleiben in `make test` und `make test-integration` grün. `TestPlayZweitesSignalNachFehlerantwort` und die Tests des Kerns fahren das zweite Signal nur ohne die Option, `TestE2EPlayFinishSession` sendet nur ein Signal. *Failure-Szenario:* Ein Umbau wie M19 lässt `play` mit der Option nach dem zweiten `Strg+C` weiterlaufen, bis die Session zu Ende ist, auch während einer langen Anfrage gegen eine Datenbank, deren Inhalt sich ändert. Ein Umbau wie M1 meldet die Unterbrechung als `PGR-E4003` mit Exit-Code 4 statt 0. Kein Gate wird rot; die Kommentare an Port und Service sagen beides zu. | `LH-FA-20.a` *Abbruchsignal* („Das zweite Signal … schließt sie … und beendet den Prozess; die unterbrochene Interaktion … kein Fehler“); Plan §6 *Fehlerantwort vor dem zweiten Signal*; `AGENTS.md` §3.10, §3.11 | `internal/hexagon/services/play.go` · `stop := context.AfterFunc(abbruch, us.Schliesse)`; · `if abbruch.Err() != nil {`; `internal/hexagon/ports/driving/play.go` · „wird ablauf geschlossen (zweites Signal), sofort.“ | ja (M1, M19; `make test`, `make test-integration`) | Zusage nur für einen Teil ihrer Fälle von einer Mutation gehalten |
| F-565 | LOW | Der Hilfetext von `play` sagt mehr zu, als die Tests halten. `--continue-on-error` „endet mit Exit-Code 4“: Nach einem abbrechenden `PGR-E6001` endet es mit 6 (`TestPlayFortsetzungAbbruch/PGR-E6001_in_Session_2`, `TestE2EPlayFortsetzungAbbruch/PGR-E6001`). `--allow-recorded-errors` „eine Fehlerantwort ist erwartet“: Eine mit `FATAL` ist `PGR-E4003` (`TestPlayErwarteterFehler/FATAL`). Plan §6, akzeptiertes Negativ (d), und §7 *Handbuch und Hilfe* sagen, die Hilfe sage nur zu, was ein Test prüft. | `AGENTS.md` §3.11; Plan §6 *Akzeptierte Negative* (d) | `internal/adapters/driving/cli/cli.go` · „der nächsten Anfrage weiter und endet mit Exit-Code 4;“; · „eine Fehlerantwort ist erwartet und kein Fehler, wenn die“ | ja (die genannten Tests gegen den Wortlaut gelesen; kein Gate prüft den Wortlaut) | Zusage im Hilfetext weiter als die Prüfung |
| F-566 | INFO | Zum Handbuch §4, das der Diff nicht ändert und von dem §7 sagt, jede Aussage zu den drei Optionen prüfe ein E2E-Test: „Mit `--continue-on-error` läuft es weiter und endet am Ende mit Exit-Code 4“ und „Mit `--allow-recorded-errors` gilt ein Fehler nicht, wenn auch die aufgezeichnete Anfrage mit einem Fehler beantwortet wurde“ tragen dieselbe Weite wie F-565; `TestE2EPlayFortsetzungAbbruch` zeigt Exit-Code 6. | Plan §7 *Handbuch und Hilfe*; `LH-FA-20.a` *Exit-Code* | `docs/user/benutzerhandbuch.md` · „Mit `--continue-on-error` läuft es weiter und endet am Ende“ | nein | — (Hinweis an den Verifier) |
| F-567 | LOW | Der Code-Commit `1b63ef8` ändert den Kommentar des Driving Ports `Player` (Schicht *Ports*, `ARC-003`) und nennt ihn in §3. Kopf (`Berührte Spec-Stellen`) und Schicht-Abgrenzung in §1 ziehen erst `ee23fad` nach. Am Stand `ee23fad` sind Kopf, §1, §3 und §8 stimmig. | `AGENTS.md` §3.9 („im selben Commit … dazu den Kopf“) | Plan, Kopf · „`ARC-002` · `ARC-003` · `ARC-005`“ | ja (`git show 1b63ef8 -- docs/plan` ändert nur §3; `make kopf-check` fängt es nicht, weil §1 die Kennung damals nicht nannte) | Plan folgt der Korrektur erst im Folge-Commit |
| F-568 | INFO | Das akzeptierte Negativ (e) schließt den abbrechenden Fehler ohne Meldungscode über den Vertrag der Ports aus, „an den auch `slice-v1-abschluss-einspielen-anmeldung` und `slice-v1-abschluss-einspielen-tls` mit neuen Fehlern im Aufbau gebunden sind“. Keiner der beiden nennt diesen Slice oder die Bindung. Getragen ist sie vom Kommentar an `Verbinde` (Einstufung `PGR-E4002`, `PGR-E4005`). Dass jeder Fehler des Adapters einen Code trägt, prüft kein Test. Die Aussage am heutigen Code habe ich an `Verbinde` und `aufbau` nachgelesen; sie hält. | Plan §6 *Akzeptierte Negative* (e); `AGENTS.md` §3.13 (gelesen, keine der vier Adressen-Klassen) | Plan §6 · „`EinspielSession` (je Fehler sein Code), an den auch“ | nein | — (Hinweis an Architect und Planner) |
| F-569 | INFO | Der Kommentar von `play` im Bootstrap sagt „Der Exit-Code ist der des Fehlers“. Seit diesem Slice ist der Fehler eine Zusammenfassung, und der Exit-Code ist der der ersten Meldung. Die Kopplung, dass der Service den abbrechenden Fehler zuerst stellt und der Bootstrap den Code der ersten Meldung nimmt, steht nur am Service und am Port. Der Bootstrap liegt außerhalb dieses Slice (§1, Schicht-Abgrenzung). | `AGENTS.md` §3.7 (Kopplung); `LH-FA-20.a` *Exit-Code* | `internal/bootstrap/bootstrap.go` · „Der Exit-Code ist der des Fehlers, ohne Fehler 0.“ | nein | — (Hinweis an den Architect für den nächsten Slice mit Bootstrap) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Hexagon, [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) | geprüft, ohne Befund. `services/play.go` importiert weiter nur `context`, `errors`, `fmt`, `model` und `ports/driven`. Der Port ändert nur den Kommentar. `make a-check`: 0 Befunde. Der Hinweis auf 13 Dateien unter `test/integration/` ohne Schicht ist Bestand; die neue E2E-Datei kommt dazu. |
| [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0014](../plan/adr/0014-konfigurationsdatei.md) | geprüft, ohne Befund. Die Sessions laufen weiter nacheinander, die Antworten werden außer `ErrorResponse` verworfen. Die beiden Lockerungen wirken wie entschieden, Verbindungsfehler brechen auch mit `--continue-on-error` ab. Die drei Optionen hängen am allgemeinen Leser mit Kommandozeile, Umgebung und `play:`, `v == "true"` wie `--force` und `--fail-on-unconsumed`. |
| Reihenfolge der Fehler im `errors.Join` | geprüft, ohne Befund. `Play` liefert bei Abbruch `errors.Join(append([]error{err}, frueher...)...)`, sonst `errors.Join(frueher...)`. `model.Meldungen` zerlegt eine reine Zusammenfassung in Reihenfolge (`nebeneinander`). Der Bootstrap schreibt je Meldung eine Zeile und nimmt den Code der ersten; das entspricht *Meldungen* und *Exit-Code*. MA ist rot. |
| Exit-Code nach früherem `PGR-E4004` | geprüft, ohne Befund. Ein abbrechender `PGR-E6001` steht vor den früheren, Exit-Code 6; `PGR-E4003` und `PGR-E4002` ergeben 4. Unit- und E2E-Test halten beide Fälle. Ein abbrechender Fehler ohne Code ist akzeptiertes Negativ (e), siehe F-568. |
| Erstes und zweites Signal mit `--finish-session-on-interrupt` | geprüft, Befund F-564. Erstes Signal: Ohne Option endet die Session vor der nächsten Interaktion, mit Option laufen alle; der Aufbau läuft mit `abbruch`, nicht mit `ctx`; vor der nächsten Session bricht die Schleife ab, auch mit Option. Fehler nach dem ersten Signal zählen wie ohne Signal. Zweites Signal: Der Code schließt die Verbindung über `abbruch` und verwirft den Fehler der Unterbrechung, frühere `PGR-E4004` bleiben. Das Verhalten ist richtig, mit der Option aber ungeprüft. |
| Erwarteter Fehler aus der Aufzeichnung | geprüft, ohne Befund. `mitFehlerantwort` sucht an jeder Stelle der aufgezeichneten Antworten, nicht nach SQLSTATE. Ein erwarteter Fehler ist keine Meldung und liest bis `ReadyForQuery` weiter, auch ohne `--continue-on-error`. `FATAL` und `PANIC` werden vor der Prüfung auf „erwartet“ ausgewertet. MB ist rot. |
| `AGENTS.md` §3.12, je Operation gegen §6 | geprüft, ohne Befund. Verglichen habe ich die Optionen und ihre Quellen, den erwarteten Fehler, die Fortsetzung bis `ReadyForQuery`, mehrere Fehlerantworten in einer Interaktion, `FATAL` bei erwartetem Fehler, das Verbindungsende beim Weiterlesen, die Reihenfolge der Zeilen und den Exit-Code. Dazu das erste Signal im Aufbau, in und nach der letzten Interaktion, zwischen Sessions und nach einer Fehlerantwort, das zweite Signal nach einer Fehlerantwort und die Kombinationen der Optionen. Jede Entscheidung des Codes steht in §6 mit Stelle in `LH-FA-20.a`. Der Code-Commit ändert §6 nicht. |
| `AGENTS.md` §3.9 | geprüft, Befund F-567. §3 folgt den Dateien des Diffs, §1 und §6 dem Gelieferten. |
| `AGENTS.md` §3.13 mit *Nachzählen beim Eintragen* | geprüft, ohne Befund außer F-568 (Hinweis). Neu genannte Adresse im Diff ist `slice-v1-abschluss-antwortvergleich` für die Kombinationen mit `--compare-responses`. Dessen DoD-Punkt 3 führt sie, und sein §1 nennt diesen Slice als Voraussetzung. `slice-v1-abschluss-einspielen-extended` führt die Teile [L·E] mit der Kennung dieses Slice. §8 zählt nach jeder Eintragung drei Liefer-Punkte und zwei Schichten. |
| `AGENTS.md` §3.10, Tabelle in §7 | geprüft, Befund F-564. Die Stichproben MA und MB werden mit dem genannten Test rot, die Zeilen der Tabelle dazu halten. |
| `AGENTS.md` §3.11, Kommentare, Hilfe, Abdeckungs-Deklarationen, Handbuch | geprüft, Befunde F-564, F-565 und F-566. Die Deklarationen der neuen Unit- und E2E-Tests sagen nur zu, was der Test prüft: Ablauf, Meldungstexte, Exit-Code und Wirkung in der Tabelle. `make abdeckung-check` ist grün. Die Hilfe am Binary entspricht dem Code. |
| `SPEC-038` *Warten in Tests* | geprüft, ohne Befund. `TestE2EPlayFinishSession` wartet an genau einer Stelle auf das Ende des Prozesses, mit 30 s als Literal und einer Meldung, die das Ereignis nennt. `warteAufSchlaf` hat eine eigene Frist von 30 s. Die Unit-Tests rufen `Play` synchron auf, außer `TestPlayZweitesSignalNachFehlerantwort` über `warteAufPlay` und `bis` mit 30 s. |
| Spezifikation `85caec4`, Strata | geprüft, ohne Befund. Die neuen Absätze nennen weder ADR noch Slice, Welle oder Commit. Die Zeile der Änderungshistorie nennt nur Lastenheft-Kennungen. |
| Hard Rules 3.3, 3.5, 3.7; Commit-Messages | geprüft, ohne Befund. Kein Move, keine ADR geändert. Die neuen Kommentare sind Zusagen oder Kopplungen. Die Messages nennen `slice-v1-abschluss-einspielen-laufsteuerung` und [`LH-FA-20`](../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung). Der Code-Commit nennt zusätzlich [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration); keiner nennt eine `SPEC-` oder `ARC-`Kennung. |
| Größe | geprüft, ohne Befund. Der Diff war in einer Sitzung prüfbar. `internal/bootstrap`, die Driven-Adapter, der PGWire-Adapter und `model` sind nicht im Diff; die Rückführungen aus §4 sind aus Sicht des Reviews nicht eingetreten. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

- F-564: Zweites Signal mit `--finish-session-on-interrupt` hält keine Mutation (M1, M19 grün in Unit und E2E).
- F-565: Hilfe von `play` sagt Exit-Code 4 und „erwartet“ weiter zu als geprüft.
- F-566: Handbuch §4 mit derselben Weite; §7 nennt es geprüft.
- F-567: Kopf und §1 folgen dem Port-Kommentar erst im Commit des Architect.
- F-568: Die Bindung von Anmeldung und TLS in (e) steht in keinem der beiden Slices.
- F-569: Der Kommentar im Bootstrap nennt die neue Kopplung nicht.

Wiederkehrende Klassen: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (F-564), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-564, F-565), `BEO-REPO/plan-folgt-korrektur-nicht` (F-567). Die Klasse von F-564 trat im vorigen Report dreimal auf (F-547 bis F-549); mit F-564 ist sie zum vierten Mal da, damit greift §Pflege des Reviewer-Skills (Hinweis für die Closure).

**Finding-Klassen dieses Laufs:** Zusage nur für einen Teil ihrer Fälle von einer Mutation gehalten · Zusage im Hilfetext weiter als die Prüfung · Plan folgt der Korrektur erst im Folge-Commit

## Verdikt

**Merge-blockierend:** ja, bis F-564 bearbeitet ist. Das gelieferte Verhalten ist an allen vier Prüfpunkten des Auftrags richtig: Reihenfolge im `errors.Join`, Exit-Code nach früherem `PGR-E4004`, erstes und zweites Signal mit der Option und erwarteter Fehler aus der Aufzeichnung. Zum zweiten Signal mit der Option fehlt der Test, der das Verhalten hält (F-564). F-565 und F-567 blockieren nicht; F-567 ist am Stand `ee23fad` behoben.

**Übergabe:** F-564 und F-565 gehen an den Implementer. Widerspricht er F-564, läuft der Konflikt über den Architect. F-567 geht an den Implementer zur Kenntnis, eine Aktion am Plan ist nicht offen. F-566 geht an den Verifier, F-568 an Architect und Planner, F-569 an den Architect. Die Finding-Klassen gehen in §7 der Closure. Dieser Report ersetzt keine Verifikation.
