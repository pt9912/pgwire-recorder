# ADR-0039: Die Aufzeichnung kennt das Verbindungsende des Servers

**Status:** Proposed

**Datum:** 2026-10-12

**Autor:** Architect; Entscheidung des Nutzers vom 2026-10-12 für die Richtung (Option B: „Die Aufzeichnung soll ein Verbindungsende kennen“); Umfang und Form sind Empfehlung, bis der Nutzer sie freigibt

**Bezug:** [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-07`](../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-11`](../../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers), [ADR-0013](0013-sqlite-recording-backend.md), [ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md), [ADR-0023](0023-antwortvergleich-entscheidung.md), [ADR-0030](0030-full-duplex-im-record-pfad.md)

**Schärft:** [`LH-FA-02.b`](../../../spec/spezifikation.md#lh-fa-02b--abschluss-einer-interaktion), [`SPEC-001`](../../../spec/spezifikation.md#spec-001--recording-serialisierung-und-versionierung), [`SPEC-043`](../../../spec/spezifikation.md#spec-043--recording-sqlite-format), [`LH-FA-20.a`](../../../spec/spezifikation.md#lh-fa-20a--einspielen), [`LH-FA-24.a`](../../../spec/spezifikation.md#lh-fa-24a--vergleich-der-antworten); die Stellen tragen die Regeln erst nach der Freigabe des Lastenhefts (Rangfolge: Lastenheft, dann Spezifikation, dann Code)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Das Lastenheft sagt in [`LH-FA-24`](../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen) (Boundary) und im Abnahmeszenario 16, ein Verbindungsende, das die Aufzeichnung zeigt und das der Server genauso verursacht, gelte beim Vergleich als erwartet. Die Aufzeichnung kann heute kein Verbindungsende zeigen: Jede Interaktion endet mit `ReadyForQuery`, und eine Interaktion, in der der Upstream nach einer `ErrorResponse` die Verbindung beendet, macht die Session zu einer nicht unterstützten (`PGR-E6001`, nicht aufgezeichnet). Beim Einspielen ist eine `ErrorResponse` mit `FATAL` oder `PANIC` deshalb stets eine Abweichung. Der Satz im Lastenheft hat keinen Fall, auf den er zutrifft.

Anwender, die das Verhalten ihrer Anwendung gegen einen vom Server beendeten Verbindungsaufbau testen wollen (Abbruch durch einen Administrator, Ablauf eines Leerlauf-Zeitlimits, Herunterfahren des Servers), können solche Läufe heute weder aufzeichnen noch wiedergeben noch vergleichen.

Annahmen: Ein Server beendet eine Verbindung mit einer Fehlerantwort `FATAL` oder `PANIC` und schließt danach den Datenstrom; ein Datenstrom, der ohne solche Antwort endet (Netzwerkabbruch, Reset), ist ein Ereignis des Netzes und kein Verhalten des Servers, das sich wiederholen ließe. Es gibt noch keine veröffentlichte Version des Produkts, also keine installierten Binaries, die eine neue Formatversion ablehnen müssten.

## Entscheidung

Wir wählen **Option B: Die Aufzeichnung kennt das Verbindungsende des Servers.** Die Richtung entschied der Nutzer am 2026-10-12 ausdrücklich, auch wenn sie vor dem ersten Release liegt und ihn verschiebt. Die folgenden Punkte sind der Umfang, den der Architect empfiehlt; sie gelten erst mit der Freigabe des Nutzers.

1. **Ein Verbindungsende ist Teil der Aufzeichnung**, wenn der Server es in einer laufenden Interaktion mit einer Fehlerantwort herbeiführt. Es ist das letzte Ereignis seiner Session und steht an der Stelle der Interaktion, an der es eintrat. Ein Ende ohne Fehlerantwort bleibt ein unerwartetes Verbindungsende und wird nicht aufgezeichnet. Der Fall „Fehlerantwort, dann Ende“ (heute nicht unterstützt) wird damit unterstützt, soweit er ein Ende des Servers ist; die Einzelfälle legt die Spezifikation fest.
2. **Das Format erhält eine neue Version.** Die Leser lesen die bisherige und die neue Version, in `yaml` wie in `sqlite`; eine ältere Version des Produkts lehnt eine Aufzeichnung mit Ende als unbekannte Version ab. Welche Version der Schreiber wählt, legt die Spezifikation fest; empfohlen ist die neue nur für Aufzeichnungen, die ein Ende enthalten.
3. **Das Wiedergeben folgt der Aufzeichnung:** Replay beendet die Verbindung an der aufgezeichneten Stelle. **Der Vergleich beim Einspielen** behandelt ein Ende, das der Server an derselben Stelle und mit derselben Fehlerklasse herbeiführt wie aufgezeichnet, als erwartet; jedes andere Ende bleibt eine Abweichung oder ein Verbindungsfehler wie bisher. Das Einspielen löst das Ende nicht selbst aus: Es sendet die aufgezeichneten Anfragen, das Ende bewirkt der Server.
4. **Reihenfolge der Änderung:** zuerst das Lastenheft (Entwurf zur Freigabe durch den Nutzer), dann die Spezifikation, dann Code. Der Code beginnt nicht vor der Entscheidung jeder Randform nach `AGENTS.md` §3.12.

Keine angenommene ADR wird ersetzt (siehe Konsequenzen). Die Regeln im Einzelnen (Fälle, Codes, Darstellung in `yaml` und `sqlite`, Versionswahl) stehen nicht hier, sondern in der Spezifikation (`AGENTS.md` §3.8).

## Verglichene Alternativen

Die Buchstaben der Optionen A und C sind die dieser ADR; B ist die Wahl des Nutzers.

| Option | Pro | Contra |
|---|---|---|
| A — Lastenheft zurücknehmen: Satz 3 des Szenarios 16 und die Boundary-Klausel streichen, die Aufzeichnung kennt weiter kein Ende | kein Formatwechsel, kein Mehraufwand, Release unverzögert | Läufe, die der Server beendet, bleiben nicht aufzeichenbar und nicht vergleichbar; jede `FATAL` ist beim Einspielen eine Abweichung; der Bedarf, der den Satz ins Lastenheft brachte, bleibt unerfüllt |
| **B — die Aufzeichnung kennt das Ende (Formatversion 2)** | Aufzeichnen, Wiedergeben und Vergleichen schließen sich zu einem Kreis; das Lastenheft sagt dann, was das Produkt tut; vor dem ersten Release kostet der Formatwechsel keine installierte Basis | Änderung an Modell, beiden Formaten, Record, Replay und Einspielen; neue Randformen; der Release verschiebt sich; ältere Binaries lehnen Aufzeichnungen mit Ende ab |
| C — kein Formatwechsel: das Ende nur beim Einspielen erwarten lassen, etwa über eine Option, die ein Ende zulässt | kein neues Format; wenig Code | die Aufzeichnung bleibt blind für das Ende: Replay kann es nicht wiedergeben, der Vergleich kennt die Stelle nicht und die Option wäre eine Behauptung des Anwenders statt Teil der Aufzeichnung |
| D — das Ende als Nachricht im bestehenden Format ohne Versionswechsel darstellen | Leser und Schreiber der Version 1 bleiben gültig | ein Leser der Version 1 hielte eine solche Interaktion für beschädigt (`PGR-E3003`) statt für unbekannt (`PGR-E3002`); die Regel „ein neues Feld erhöht die Version“ ([`SPEC-001`](../../../spec/spezifikation.md#spec-001--recording-serialisierung-und-versionierung)) würde umgangen |

## Konsequenzen

- Positiv: Das Verbindungsende ist reproduzierbar und prüfbar; Szenario 16 hat in allen Sätzen einen Fall; das Lastenheft und die Aufzeichnung sagen dasselbe.
- Negativ: Der Release verschiebt sich um die Slices dieser Änderung. Das Format bekommt eine zweite Version, die jeder Leser dauerhaft trägt. Aufzeichnungen mit Ende lesen nur Binaries mit dieser Version.
- Akzeptiertes Negativ: Ein Ende ohne Fehlerantwort (Netzwerkabbruch) und ein Ende im Aufbau vor dem ersten `ReadyForQuery` bleiben nicht aufzeichenbar. Der Grund: Sie sind keine Zusage des Servers, die ein zweiter Lauf wiederholen könnte, und das Einspielen kann sie nicht auslösen.
- Beziehung zu angenommenen ADRs, ohne `Supersedes`: [ADR-0017](0017-einspielen-sequenziell-und-fehlersemantik.md) hält, dass ein Verbindungsende beim Einspielen ein Verbindungsfehler ist; B nimmt davon das aufgezeichnete, gleich herbeigeführte Ende aus, das die Spezifikation als Regel festlegt, wie [ADR-0023](0023-antwortvergleich-entscheidung.md) es für Fehlerantworten tat. [ADR-0013](0013-sqlite-recording-backend.md) und die Versionsregel in [`SPEC-001`](../../../spec/spezifikation.md#spec-001--recording-serialisierung-und-versionierung) bleiben gültig; B ergänzt eine Version. [ADR-0030](0030-full-duplex-im-record-pfad.md) weist die Einstufung eines Verbindungsendes dem Record-Service zu; das gilt weiter. Die Regel „Fehlerantwort vor dem Abbruch ist nicht unterstützt“ steht nur in der Spezifikation und ändert sich dort.
- Folgepflicht: Lastenheft-Entwurf zur Freigabe durch den Nutzer; danach Spezifikation (die Stellen oben und `LH-FA-11.a`, `LH-FA-03.b`, `LH-FA-12.a`); danach ein Schnitt in Slices durch den Planner. Jede Zusage des neuen Vertrags bekommt eine Mutation (`AGENTS.md` §3.10), jede Randform ist vor dem Code entschieden (§3.12). Beim Annehmen zieht der Index mit.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | nicht maschinell prüfbar; die Zusagen des Vertrags prüfen die Tests der Slices | — |

## Re-Evaluierungs-Trigger

Wenn Anwender ein Ende ohne Fehlerantwort oder ein Ende im Aufbau aufzeichnen wollen, oder wenn ein Server eine Verbindung ohne Fehlerantwort beendet, die das Produkt reproduzierbar auslösen kann.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-12 | Proposed; Richtung (Option B) vom Nutzer entschieden, Umfang Empfehlung | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
