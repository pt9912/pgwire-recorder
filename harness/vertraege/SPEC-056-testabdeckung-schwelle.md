# SPEC-056 — Schwelle der Anweisungsabdeckung der Unit-Tests (`test`, `coverage-gegenprobe`)

Die Stufe `test` des `Dockerfile` (Ziel `make test`) misst die Anweisungsabdeckung der
Unit-Tests und endet rot, wenn sie unter der Schwelle liegt; `make coverage-gegenprobe`
prüft das. Das ist Teil 2 der Messmethode von
[`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes). Die Abdeckung je
Anforderung und Pfad ist ein anderer Gegenstand ([`SPEC-048`](SPEC-048-abdeckung.md)). Diese
Kennung legt fest, was gemessen wird, woran die Stufe rot wird und was sie nicht zusagt. Die Regel, nach der die
Zahl der Messung folgt, und den Quellstand der Messung nennt die Entscheidung
([ADR-0044](../../docs/plan/adr/0044-schwelle-der-anweisungsabdeckung.md)).

Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes, den er prüft. Die Gegenprobe
prüft den Prüfschritt `tools/harness/coverage.sh` an Profilen, die sie schreibt (ohne Docker),
und die Stufe in Kopien des Arbeitsbaums unter eigenem Temp-Pfad (mit Docker, über `DOCKER`).
Die unveränderte Kopie ist grün, und je Fall bricht eine Mutation genau ihn. Unter Punkt 4, 6 und
10 zählt jedes Element der Aufzählung einzeln.

1. **Schwelle.** Die Schwelle steht an einer Stelle, in dieser Datei, in der Zeile
**Schwelle:** `94` Prozent der Anweisungen.
   (Die Zeile steht ohne Einrückung am Zeilenanfang; was davor steht, ist gleich; sie
   beginnt mit `**Schwelle:** `, der Wert steht in Backticks. Eine eingerückte Zeile zählt nicht.) Einheit: Prozent der Anweisungen, ganze
   Zahl. Grund: gemessen 94,45 % (4201 von 4448 Anweisungen) bei Quellstand `69aaf3f7`, auf die ganze
   Zahl abgerundet; der Abstand von 19 Anweisungen ist der Spielraum für Umbauten. Die Regel, nach der
   die Zahl folgt, steht in der Entscheidung. `make test` liest den Wert aus dieser Zeile und gibt
   ihn der Stufe als Build-Argument ohne Vorgabewert; das `Dockerfile` führt keine Zahl. Eine Senkung
   ist eine neue Entscheidung (`AGENTS.md` §3.6), eine Anhebung ein Commit dieser Zeile.
2. **Messung.** Gemessen werden die Unit-Tests der Stufe `test`: ein Lauf von `go test` über alle
   Pakete mit `-covermode=set` und ohne `-coverpkg`. Die Zahl meint also den Anteil der Anweisungen,
   die ein Test des eigenen Pakets mindestens einmal ausführt. Die Integrationstests zählen nicht mit
   (Grund: ein Profil aus dem Prozess des Binaries bräuchte einen Weg aus dem Container, den die
   Stufen nicht haben, und die Stufe schreibt nichts in den Arbeitsbaum).
3. **Gesamtwert.** Der Wert ist die Summe der ausgeführten Anweisungen durch die Summe aller
   Anweisungen im Profil, über alle Pakete, nicht das Mittel der Pakete. Eine Schwelle je Paket gibt
   es nicht; die Ausgabe nennt jedes Paket (Punkt 7).
4. **Pakete.** Jedes Paket mit mindestens einer Anweisung im Profil zählt, ohne Ausnahmeliste:
   (a) ein Paket `main` ohne Testdatei (`cmd/`), (b) ein anderes Paket ohne Testdatei, auch ein
   Hilfspaket für Tests, steht mit 0 % im Profil und zählt mit seinen Anweisungen; (c) ein Block mit
   der Anzahl 0 zählt nichts, und ein Paket nur aus solchen Blöcken hat keine Zeile in der Ausgabe;
   (d) ein Paket ohne Anweisung (Ports, `doc.go`) hat keinen Block und steht nicht in der Ausgabe.
   Einen Test eines anderen Pakets zählt ein Paket nicht (Punkt 2). Für generierten Code gilt keine
   Regel: Der Prüfschritt liest nur das Profil, und die Anweisungen zählen wie alle.
5. **Vergleich.** Grün ist der Stand, wenn `ausgeführt * 100 >= Schwelle * gesamt` in ganzen
   Zahlen; ein Stand genau an der Schwelle ist grün, einer eine Anweisung darunter rot. Es wird
   vor dem Vergleich nicht gerundet. Der angezeigte Wert ist auf zwei Nachkommastellen abgerundet,
   sodass ein angezeigter Wert unter der Schwelle immer rot bedeutet.
6. **Eingabe des Prüfschritts.** `tools/harness/coverage.sh <Profil> <Schwelle>` endet mit
   Exit-Code ungleich 0 und nennt die Form, nie grün, bei: (a) Profil fehlt oder ist nicht lesbar,
   (b) Profil ist leer, (c) erste Zeile fehlt oder ist nicht `mode: set`, (d) eine Zeile ohne die
   Form `<Datei>:<Bereich> <Anzahl> <Zähler>` mit Anzahl und Zähler als ganzen Zahlen, (e) ein
   Block (Datei und Bereich) kommt zweimal vor, (f) die Summe der Anweisungen ist 0, (g) Schwelle
   fehlt oder ist leer, (h) Schwelle ist keine ganze Zahl von 1 bis 100 in Dezimalschreibweise ohne
   führende Null (`94.5`, `094`, `0`, `101`, `94 94` und `-1` sind abgelehnt). Zwei Zeilen
   `**Schwelle:**` in dieser Datei ergeben `94 94` und damit (h). Hat eine Eingabe mehrere
   Gründe, gilt die Reihenfolge (a), (b), (c), (g), (h), dann (d), (e), (f): Profil und Schwelle
   werden vor den Zeilen geprüft; die Zeilen werden der Reihe nach gelesen, und die erste Zeile,
   die (d) oder (e) verletzt, bestimmt die Meldung; (f) folgt erst, wenn alle Zeilen
   angenommen sind.
7. **Ausgabe.** Der Prüfschritt schreibt, bei grün und bei rot gleich, je Paket eine Zeile mit
   Importpfad, Wert und Zahl der Anweisungen, aufsteigend nach Importpfad, dann als letzte Zeile
   den Gesamtwert mit Schwelle. Bei rot steht in der Gesamtzeile zusätzlich die Differenz (Schwelle
   minus angezeigter Wert, in Prozentpunkten), und der Exit-Code ist ungleich 0. Die Ausgabe
   erscheint nur, wenn die Schicht läuft; bei einem Treffer im Build-Cache gibt BuildKit sie nicht
   noch einmal aus (ein roter Stand wird nie zwischengespeichert).
8. **Stufe `test`.** Die Stufe fährt in dieser Reihenfolge `gofmt`, `go vet` über allen Code
   einschließlich der Integrationstests, die Unit-Tests mit Profil und den Prüfschritt. Jeder
   Schritt bricht die Stufe, sodass ein roter Test auch bei erfüllter Schwelle rot bleibt und ein
   Formatfehler auch bei erfüllter Schwelle. Der Prüfschritt liegt in der Allowlist von
   `.dockerignore` und braucht in der Stufe nur das Image der Stufe; die Stufe bleibt netzlos
   und schreibt nichts in den Arbeitsbaum. Die Stufe bleibt in der schnellen Kette
   ([`SPEC-052`](SPEC-052-gate-kette.md) Punkt 2): Sie liest nichts unter `docs/` und nicht die
   `README.md`; die Schwelle liest das Ziel auf dem Host aus dieser Datei, nicht der Build-Kontext.
9. **Aufrufarten des Treibers.** Das Rezept von `test` bleibt eine Zeile ohne `$(MAKE)` und ohne
   `+`; die Schwelle liest ein `$(shell)` beim Einlesen des Makefiles, das nur liest. Damit fährt
   `make` unter `-n`, `-q` und `-t` kein Docker, und `MAKEFLAGS` aus der Umgebung gilt wie das Flag;
   je Flag und für `MAKEFLAGS` ein Fall mit einer Attrappe von `docker` (über `DOCKER`), die
   Mutation ist ein Rezept mit `$(MAKE)`. `-i` verhält sich wie bei jedem Gate (`SPEC-052`
   Punkt 21) und hat keinen eigenen Fall.
10. **Gegenprobe.** Fälle ohne Docker (Profile, die sie schreibt): Punkt 3 (Gewicht nach
    Anweisungen, nicht Mittel der Pakete), Punkt 4 (a) bis (c), Punkt 5 (genau an der Schwelle grün,
    eine Anweisung darunter rot, angezeigter Wert abgerundet), Punkt 6 (a) bis (h) je Form und die Reihenfolge der Gründe je benachbartem Paar der Folge (a)/(b), (b)/(c), (c)/(g), (g)/(h), (h)/(d), (d)/(e), (e)/(f) sowie (d) und (e) in vertauschter Zeilenfolge, Punkt 1 (eine eingerückte Zeile zählt nicht, eine Zeile mitten im Text zählt nicht), Punkt 7
    (Liste bei grün und bei rot, Reihenfolge, Gesamtzeile zuletzt, Differenz nur bei rot). Fälle
    mit Docker, je einer: Schwelle 100 in dieser Datei der Kopie wird rot (die Schwelle kommt in der
    Stufe an; die Mutation setzt einen Vorgabewert im `Dockerfile`), Zeile `**Schwelle:**` fehlt
    wird rot (Punkt 1), ein neues Paket ohne Test mit vielen Anweisungen, das nur ein Test eines
    anderen Pakets ausführt, wird rot (Punkt 2 und 4; die Mutation ist `-coverpkg=./...`), ein
    Formatfehler und ein roter Test bleiben bei erfüllter Schwelle rot (Punkt 8). Die Gegenprobe
    nennt in `harness/README.md` §Sensors ihre Bindung und trägt die Abdeckungs-Deklaration von
    Teil 2 von `LH-QA-07` ([`SPEC-048`](SPEC-048-abdeckung.md)).
11. **Anschluss der Gegenprobe.** `coverage-gegenprobe` steht in der Liste für den Container
    ([`SPEC-060`](SPEC-060-gegenproben.md) Punkte 1 und 2), nicht in `GATE_CHECKS`; sie fährt `docker build`.
    Sie liest `docs/user/` und die `README.md` des Arbeitsbaums nicht (Regel an den Schreiber,
    `SPEC-052` Punkt 13): Sie baut aus einer Kopie unter
    eigenem Temp-Pfad, und der Build-Kontext schließt beides aus.
12. **Grenzen.** (1) Die Zahl sagt, dass Anweisungen laufen, nicht, dass etwas geprüft wird; die
    Prüfung bleibt `AGENTS.md` §3.10. (2) Ein Profil mit `atomic` oder `count` (etwa mit `-race`)
    ist abgelehnt (Punkt 6); ein Umstieg ist eine neue Entscheidung. (3) Das Gate vergleicht den
    Stand nicht mit dem vorigen und hebt die Schwelle nicht an. (4) Die Zeilen der Ausgabe
    erscheinen nur bei laufender Schicht (Punkt 7). (5) Gleichzeitige Läufe sind nicht zugesagt.
