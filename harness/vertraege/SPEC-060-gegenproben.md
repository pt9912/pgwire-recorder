# SPEC-060 — Gegenproben außerhalb der Gate-Kette und ihre Sammelziele (`gegenproben`, `gegenproben-host`)

Die Gegenproben der Gates und `a-check-negativ` sind kein Teil von `make gates` und nicht von
`make gates-schnell`. Sie stehen in genau einer von zwei Listen, und je ein Sammelziel fährt eine Liste:
`make gegenproben` die Gegenproben für den Dev-Container, `make gegenproben-host` die, die `u-boot` oder
`devcontainer` brauchen und deshalb nur am Host laufen. Beide Sammelziele sind kein Gate; kein Hook und kein
Stand glaubt ihr Ergebnis, geglaubt wird allein der Exit-Code, den der Aufrufer liest. Die Kennung legt
fest, wann ein Ziel in welche Liste gehört, wie die Gegenprobe der Kette das hält, was die Sammelziele
zusagen, woran sie scheitern und was die Pflicht zum Fahren ist. Der Grund der Aufteilung steht in
[ADR-0048](../../docs/plan/adr/0048-gegenproben-verlassen-die-gate-kette.md); diese Kennung wiederholt
ihn nicht. Die Namen der Listen und Sammelziele und die Entscheidungen, dass `make gates-schnell`
vorerst bleibt und weder ein Stand der Gegenproben noch ein Hook, der ihn liest, entsteht, sind
Entscheidung des Nutzers vom 2026-10-18; ändert er sie, ändert sich nur diese Kennung mit den Fragmenten.

Der Vertrag hat zwei Teile an einer Naht: **Teil 1** (Punkte 1 bis 3) legt Listen und Namensregel fest und
wird von `make gates-schnell-gegenprobe` geprüft; **Teil 2** (Punkte 4 bis 12) legt die Sammelziele fest und
wird von `make gegenproben-gegenprobe` geprüft, einer Gegenprobe mit Attrappen als Einzelziele, ohne Docker.
Die Punkte 13 bis 15 gelten für beide Teile und haben keinen eigenen Fall, außer wo sie einen nennen. **Punkt 16** (Zeit und Platz einer Gegenprobe) gilt für einzelne Gegenproben, nicht für die Sammelziele; seine Fälle stehen in der Gegenprobe, die er nennt. **Punkt 17** (Start mit ignoriertem `SIGINT`) gilt für jede Gegenprobe mit einem Signalfall und für das Skript der Sammelziele und hat keinen Fall. Der
Name jedes Falls trägt die Nummer des Punktes, den er prüft. Je Zusage und je Element einer Aufzählung
gibt es einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die Gegenprobe verlangt
genau diese Menge roter Fälle je Mutation (mehr und weniger machen sie rot).

**Plattform.** Zugesagt ist ein Linux-Host oder -Container mit bash ≥ 4, GNU `make` und GNU-Werkzeugen;
das Skript der Sammelziele braucht zusätzlich `pgrep` (procps), die Gegenprobe der Sammelziele `setsid`, `pgrep` und `timeout`. macOS- und Windows-Hosts sind nicht zugesagt
(akzeptiertes Negativ; Grund: kein solcher Host steht zur Messung bereit, und die Einzelziele brauchen
Docker mit Linux-Semantik; die Skripte lehnen sie nicht ab). Die Sammelziele fahren Einzelziele, deren
Plattform ihr eigener Vertrag nennt.

**Randform → Punkt.**

| Randform des Plans | Punkt |
|---|---|
| K1 Namensregel | 2 |
| K2 Zuordnung zu Container oder Host | 3 |
| K3 Ausnahmeliste der schnellen Kette | 2 (letzter Satz), [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 4 |
| K4 Nachweis des vollen Laufs und Hook | 14, [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 9 |
| K5 Beleg für die Positivliste | [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 13 |
| K6 Aufrufarten von `make` für `make gates` | Bestand, [`SPEC-052`](SPEC-052-gate-kette.md) Punkte 6 und 21 |
| K7 `make gates-schnell` neben `make gates` | [`SPEC-052`](SPEC-052-gate-kette.md) Punkte 1 und 2 (beide Ketten bleiben); ob sie bleiben, entscheidet der Nutzer |
| G1 zwei Ziele, zwei Listen, ein Skript | 1, 4 |
| G2 Reihenfolge | 5 |
| G3 rotes Einzelziel, Abbruch | 6, 8 |
| G4 Aufrufarten | 10 |
| G5 fehlendes Werkzeug | 6 |
| G6 Ausgabe | 6, 7 |
| G7 leere Liste | 9 |
| G8 Ziel in beiden Listen oder in keiner | 2 |
| G9 Aufruf aus einem Unterverzeichnis | 13 |
| G10 gleichzeitige Läufe, Aufruf-Voraussetzung | 13 |
| G11 Erkennung des Dev-Containers entfällt | 3 |
| Eintrag ohne Ziel, doppelter Eintrag | 2 |
| Reihenfolge mehrerer Gründe | 12 |
| Kein Stand | 11 |
| Plattform | Kopf dieser Kennung, 15 |
| Z1 Obergrenze der Zeit eines vollen Laufs, Überschreitung | 16 (1) bis (4) |
| Z2 Arbeitsordner je Mutation und je Fall, rote Mutation, Abbruch, Parallelität | 16 (5) |
| Z3 Aufrufarten (`GP_MUT`, `GP_FALL`, `GP_ENTDECKE`, `GP_PARALLEL`, `GP_BUDGET`, `make`), Reihenfolge der Gründe | 16 (3), (4), (7) |
| Z4 Plattform und Umgebung der Fahrt (Last, Platz, root, Locale, Zeitzone) | 16 (6), (8) |
| Z5 Ergebnisgleichheit nach einer Beschleunigung | 16 (9) |
| Z6 Start mit ignoriertem `SIGINT`, Rezept für den Hintergrundstart, Plattform | 17 (für das Skript der Sammelziele auch 8) |

**Teil 1: Listen und Namensregel.**

1. **Zwei Listen.** `GEGENPROBEN_CONTAINER` und `GEGENPROBEN_HOST` sind Variablen von `make`. Jedes Fragment
   unter `harness/mk/`, das eine Gegenprobe definiert, hängt deren Namen mit `+=` an genau eine von beiden; ein
   Eintrag ist der Name eines Make-Ziels, ohne Leerraum und ohne führendes `-`. Die Reihenfolge einer Liste
   ist die, in der `make` sie ausgibt. Die Sammelziele selbst stehen in keiner Liste. Die Liste für den
   Dev-Container nimmt jede Gegenprobe, die weder `u-boot` noch `devcontainer` aufruft; die Liste für den Host
   nimmt jede, die eines von beiden aufruft (Punkt 3). **Keine Initialisierung** (Entscheidung des Architect): Beide Listen
   werden nirgends mit `:=` gesetzt, weil das Makefile generiert ist und ein Fragment, das setzt, je nach Reihenfolge der
   Fragmente Einträge verwirft. `+=` auf eine ungesetzte Variable ergibt eine rekursive Variable; zugesagt ist deshalb nur der
   **ausgewertete** Wert, den `make` im Rezept und in einer Regel mit `$(…)` liefert, nicht der Text in `make -p`. Wer die Liste liest
   (die Sammelziele, die Regel von Punkt 2), liest ihn ausgewertet; die Regel liest die Listen und `GATE_CHECKS` über `make --eval`
   mit einem Ziel, das den Wert ausgibt; ein Fragment darf einen Eintrag aus einer Variable bilden. Fall
   `p1-liste-aus-variable`: ein Eintrag, der aus einer Variable gebildet wird, ist in der Liste und läuft im Sammelziel; Mutation: die
   Regel oder das Sammelziel liest den unausgewerteten Text (`make -p`).
2. **Namensregel und Eindeutigkeit.** Als Gegenprobe im Sinn dieser Kennung gilt jedes Ziel mit Rezept in einem
   Fragment unter `harness/mk/`, dessen Name auf `-gegenprobe` endet (Element N1), und `a-check-negativ`
   (Element N2). Jedes solche Ziel steht in genau einer der beiden Listen, in dieser genau einmal, und in
   keiner Hälfte von `GATE_CHECKS` (weder in der vollen Kette noch in der schnellen). Ein Eintrag einer Liste hat ein Ziel
   mit Rezept. Eine Verletzung macht `make gates-schnell-gegenprobe` rot, im Arbeitsbaum selbst wie in
   der Kopie. Der Name `gegenproben-gegenprobe` (die Gegenprobe der Sammelziele) steht in der Liste für den
   Dev-Container. Fälle, je einer für N1 und N2: (V1) das Ziel steht in `GATE_CHECKS`; (V2) es steht in beiden Listen;
   (V3) es steht in keiner Liste; (V4) es steht zweimal in derselben Liste (Fälle `p2-v4-n1` und `p2-v4-n2` in der Liste für den Container, `p2-v4-host` in der für den Host); dazu je Liste ein Fall (V6) mit einem Eintrag
   ohne Ziel; dazu ein Fall `p2-andere-namen`: `gegenproben`, `gegenproben-host`, ein Ziel `x-negativ`, ein Ziel
   `x-gegenprobe-alt` und ein Ziel `gegenprobe` stehen in `GATE_CHECKS` und in keiner Liste ohne Befund (die Regel
   erkennt nur den Namen; akzeptiertes Negativ, Grund: eine Regel am Namen ist so weit wie der Name). Mutationen: je
   Fall die Prüfung, die ihn fängt, gestrichen; für `p2-andere-namen` eine Prüfung, die den Teilstring `gegenprobe`
   statt des Endes `-gegenprobe` sucht, und je Name der Aufzählung (`gegenproben`, `gegenproben-host`, `x-negativ`, `x-gegenprobe-alt`,
   `gegenprobe`) eine Mutation, die nur ihn als Gegenprobe liest (`p2-andere-<Name>`). Fall `p2-rezept-ausserhalb-fragment`: ein Ziel mit Namen auf `-gegenprobe`, dessen Rezept im
   `Makefile` und nicht in einem Fragment unter `harness/mk/` steht, ohne Befund (die Regel gilt dem Fragment; akzeptiertes Negativ,
   Grund: das `Makefile` ist generiert und trägt kein Ziel mit Rezept außer den Hilfsmitteln); Mutation: die Regel sucht auch außerhalb
   des Verzeichnisses (`p2-rezept-ohne-herkunft`). Die grüne Kopie hat beide Listen mit Einträgen. Die Ausnahmeliste der
   schnellen Kette nennt kein Ziel dieser Kennung: ein Eintrag, der nicht in `GATE_CHECKS` steht, macht
   [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 4 rot, ein Eintrag, der dort steht, Punkt 2 dieser Kennung (V1); eine eigene Prüfung gibt es dafür nicht.
   **Sprache von `make`** (Voraussetzung): Die Regel ruft `make` immer mit `LC_ALL=C` auf; die Ausgabe von `make -p` und der Meldungen, die sie
   liest, hängt sonst von der Sprache des Hosts ab (GNU make 4.3 in `de_DE.UTF-8` meldet „Rezept (aus …“ statt „recipe to execute (from …“).
   Fall `p2-locale`: `LC_ALL=en_US.utf8` mit `LANGUAGE=de` erzwingt die deutsche Ausgabe auch ohne installierte Locale `de_DE`; der Fall
   prüft vorab mit `grep` auf „Rezept (aus“, dass sie so erzwungen ist, und die Regel bleibt grün; ist keine Locale installiert, die das
   erzwingt, ist er rot und nennt es (kein Überspringen). Mutation `p2-locale-nicht-c`: `LC_ALL=C` entfällt. **Arbeitsbaum:** Fall
   `p2-arbeitsbaum`: eine Kopie der Gegenprobe als eigener Arbeitsbaum, in dem `lint-gegenprobe` in `GATE_CHECKS` steht, endet mit Exit
   ungleich 0 (die Zusage „im Arbeitsbaum selbst wie in der Kopie“); Mutation `p2-arbeitsbaum-gestrichen`: die Prüfung des Arbeitsbaums
   ist durch `ist=""` ersetzt, und die Kopie bleibt grün.
3. **Zuordnung ist Entscheidung des Schreibers.** Wer eine Gegenprobe anlegt, wählt die Liste mit einer Zeile im
   Fragment: Host, wenn sie `u-boot` oder `devcontainer` aufruft, sonst Container. Die Gegenprobe der Kette prüft
   nicht die Wahl, nur dass es eine einzige gibt (Punkt 2), weil sie nicht weiß, was ein Skript aufruft
   (akzeptiertes Negativ; Grund: ein Parser der Skripte wäre eine Liste, die dieselbe Hand pflegt wie die
   Gegenprobe). Die Folgen eines Irrtums: Eine Host-Gegenprobe in der Liste für den Container endet bei der Fahrt
   im Gegenproben-Container rot mit der Meldung des fehlenden Werkzeugs; eine Container-Gegenprobe in der Liste
   für den Host läuft am Host unbemerkt weiter (akzeptiertes Negativ; Grund: der falsche Fehler ist der erste,
   und die Fahrt am Host ist ohnehin Pflicht, Punkt 14). Es gibt keine Erkennung des Dev-Containers und kein
   Überspringen eines Falls; fehlt einer Gegenprobe ein Werkzeug, ist sie rot und nennt es
   ([ADR-0047](../../docs/plan/adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md)). Im Bestand
   dieser Entscheidung steht `container-start-gegenprobe` in der Liste für den Host und kein anderes Ziel.
   Jede Gegenprobe der Liste für den Container ist ohne vorherigen Lauf von `make gates` grün; die Lieferung
   belegt es durch eine Fahrt von `make gegenproben` in einem frischen Gegenproben-Container (Regel an die
   Lieferung, kein Fall; Grund: ein Fall dafür wäre ein Gate, das diese Entscheidung verlässt).

**Teil 2: die Sammelziele.**

**Exit-Codes und Naht der Gegenprobe** (Entscheidung des Architect): Die Exit-Codes der Punkte 6, 8, 9, 10 und 12 sind die des **Skripts**
`tools/harness/gegenproben.sh`. `make` selbst gibt jedem scheiternden Rezept Exit 2 (bei `-i` Exit 0); über `make` ist deshalb nur
„Exit ungleich 0“, das Ziel, die Ausgabe und die Reihenfolge zugesagt. Die Gegenprobe prüft die Codes am Skript, das sie mit einer Attrappe
als `make` ruft (Aufruf `gegenproben.sh <make> <sammelziel> <ziel>…`; Punkt 13 gilt dem Aufrufer, nicht der Gegenprobe), und über `make`
das Rezept und die Ausgabe. Bei `-i` bleibt über `make` nur die Meldung.

4. **Zwei Ziele, ein Skript.** `make gegenproben` fährt `GEGENPROBEN_CONTAINER`, `make gegenproben-host` fährt
   `GEGENPROBEN_HOST`, beide über dasselbe Skript `tools/harness/gegenproben.sh` aus dem Fragment
   `harness/mk/gegenproben.mk`. Beide sind `.PHONY`, kein Gate, nicht in `GATE_CHECKS`, nicht in der Ausnahmeliste, hängen
   nicht voneinander ab, und `make gates` und `make gates-schnell` fahren keines von beiden. Ein Ziel fährt die Liste des
   anderen nie. Fälle: `p4-container` (nur die Einzelziele der Liste für den Container laufen), `p4-host` (nur die der
   Liste für den Host), `p4-kein-gate` (`make -n gates` nennt keines der beiden Sammelziele, und beide stehen nicht in `GATE_CHECKS`), `p4-kein-ausnahme` (keines
   steht in `GATE_SCHNELL_AUSNAHME`). Mutationen: das
   Ziel für den Container fährt auch die Liste für den Host; das Ziel für den Host fährt auch die Liste für den
   Container; das eine hängt an dem anderen (je Richtung); je Sammelziel einzeln: es steht in `GATE_CHECKS` (`p4-kein-gate`), es steht in
   `GATE_SCHNELL_AUSNAHME` (`p4-kein-ausnahme`).
5. **Nacheinander, in der Reihenfolge der Liste.** Jedes Einzelziel beginnt nach dem Ende des vorigen, in der
   Reihenfolge der Liste, auch unter `-j` (eine Schleife im Skript, kein Fächer von Prerequisites). Fälle: `p5-reihenfolge`
   (die Attrappen schreiben `beginn`/`ende` und warten kurz; die Folge im Protokoll ist `beginn a, ende a, beginn b, …` in
   der Reihenfolge der Liste), `p5-j` (dasselbe unter `-j4`). Mutationen: Einzelziele im Hintergrund gestartet; die Liste
   rückwärts gefahren; die Einzelziele als Prerequisites des Sammelziels statt in der Schleife.
6. **Rotes Einzelziel, Ausgang, Zusammenfassung.** Alle Einzelziele der Liste laufen, auch nach einem roten;
   ein Einzelziel ist rot bei jedem Exit-Code ungleich 0, auch 127 (Werkzeug fehlt) und 2 (kein Ziel); `<n>` ist der Exit von `<make> <ziel>`:
   bei einem scheiternden Rezept unter echtem `make` immer 2, 127 nur, wenn `make` selbst fehlt (der Fall `p6-127` läuft an der Attrappe); das Sammelziel
   meldet das fehlende Werkzeug nicht selbst und prüft vorab nichts. Der Exit-Code des Sammelziels ist 0, wenn alle
   grün sind, sonst 1. Nach der Ausgabe des letzten Einzelziels steht auf stderr je Einzelziel in der Reihenfolge der Liste eine
   Zeile `<sammelziel>: <ziel>: grün` oder `<sammelziel>: <ziel>: rot (Exit <n>)`; sonst schreibt das
   Skript nichts. Fälle: `p6-alle-gruen` (Exit 0, eine grüne Zeile je Ziel), `p6-erstes-rot`, `p6-mittleres-rot`,
   `p6-letztes-rot` (Exit 1, die übrigen laufen und stehen in der Zusammenfassung), `p6-zwei-rot` (beide roten stehen mit
   Namen und Exit-Code dort), `p6-127` (ein Ziel endet mit 127: rot mit `Exit 127`, das nächste läuft). Mutationen: Abbruch
   beim ersten roten; Exit 0 trotz rotem; Exit-Code nur des letzten Einzelziels; Zusammenfassung nur der roten;
   Zusammenfassung ohne Exit-Code; Zusammenfassung vor der Ausgabe des letzten Ziels.
7. **Ausgabe.** Die Ausgabe jedes Einzelziels (stdout und stderr) erscheint unverändert und auf demselben Strom;
   das Skript schreibt außer der Zusammenfassung und den Meldungen der Punkte 9, 10 und 12 nichts. Vom Wortlaut der Meldung bei `-i`
   (Punkt 10 b) ist nur der Anfang `<sammelziel>: -i ` zugesagt und geprüft. Der Standardeingang jedes Einzelziels ist `/dev/null`
   (Folge des Hintergrundstarts, den die Behandlung der Signale braucht, Punkt 8); ein Einzelziel, das liest, bekommt Ende der Eingabe
   (akzeptiertes Negativ, ohne Fall; Grund: ein Terminal als Eingang wäre mit dem Abbruch nicht zu vereinbaren). Zeilen von `make`
   selbst (Verzeichniswechsel, Warnungen zum Jobserver unter `-j`) sind nicht zugesagt. Fall `p7-ausgabe`: die
   Attrappe schreibt eine Zeile auf stdout und eine auf stderr; beide erscheinen auf ihrem Strom. Fall `p7-eingabe`: der Standardeingang jedes Einzelziels ist `/dev/null`
   (Eingabe leer, Ende sofort). Mutationen: stderr nach stdout geleitet; stdout verworfen; der Eingang ist `/dev/zero` statt `/dev/null`.
8. **Abbruch.** Ein `SIGINT` oder `SIGTERM` beendet das laufende Einzelziel und startet kein weiteres; das Skript endet
   mit 130 (`SIGINT`) oder 143 (`SIGTERM`) und schreibt keine Zusammenfassung. Die Zusage gilt für ein Signal an die
   Prozessgruppe und für ein Signal nur an das Skript. Fälle: je Signal und Adressat einer, vier insgesamt
   (`p8-sigint-gruppe`, `p8-sigterm-gruppe`, `p8-sigint-skript`, `p8-sigterm-skript`); die Attrappe des mittleren von drei Zielen
   wartet, das dritte schreibt bei Beginn in das Protokoll; die Attrappe startet außerdem ein Hintergrundkind (`sleep 31`); geprüft werden Exit-Code, dass die Attrappe **und ihre Nachkommen** beendet sind, dass das
   dritte Ziel nicht begann und dass keine Zusammenfassung erscheint. Das Skript beendet das laufende Einzelziel samt seinen Nachkommen
   (Voraussetzung `pgrep`, Kopf). **Voraussetzung:** Das Skript wird mit nicht ignoriertem `SIGINT` gestartet; ein Aufrufer, der es mit
   ignoriertem `SIGINT` startet, ist nicht zugesagt (akzeptiertes Negativ, allgemein Punkt 17; Grund: eine ignorierte Disposition erbt jedes Kind, und die
   Gegenprobe setzt sie für ihren Hintergrundstart zurück). Mutationen: das nächste Ziel startet nach dem Signal; die Exit-Codes
   130 und 143 vertauscht; je Signal einer: `SIGINT` endet mit 143 (färbt die beiden Fälle mit `SIGINT`), `SIGTERM` endet mit 130 (die beiden mit `SIGTERM`);
   die Behandlung von `SIGINT` gestrichen (färbt `p8-sigint-skript`); die von `SIGTERM` gestrichen (färbt `p8-sigterm-skript`); das laufende Einzelziel bleibt am Leben (sichtbar nur im Fall mit dem Signal an das Skript);
   die Nachkommen des Einzelziels bleiben am Leben (ebenso); eine Zusammenfassung erscheint.
   `SIGKILL` und ein Ausfall von Docker sind nicht zugesagt (Punkt 15).
9. **Leere Liste.** Eine leere Liste ist ein Fehler, kein stilles Grün: Exit 2, eine Zeile `<sammelziel>: Liste leer` auf stderr, kein
   Einzelziel. Fälle `p9-container` und `p9-host`, je Sammelziel einer (das Fragment der Kopie hängt nichts an die Liste). Mutationen:
   Exit 0;
   das Fragment übergibt dem Skript die Liste des anderen Sammelziels (Rezepte mit vertauschten Listen; das Skript kennt nur die ihm
   übergebene Liste). Diese Mutation färbt außer `p9-*` auch Fälle der Punkte 4, 5, 6 und 10 (`p4-*`, `p5-*`, `p6-erstes-rot`, `p10-j`,
   `p10-k`); die Gegenprobe verlangt genau diese Menge.
10. **Aufrufarten von `make`.** Zugesagt ist das Ergebnis, nicht der Mechanismus; das Rezept übergibt `$(MAKE)` an das Skript, `make`
    führt es deshalb auch unter `-n`, `-q` und `-t` aus, und das Skript entscheidet an den Buchstaben im ersten Wort von
    `MAKEFLAGS`, wenn es nicht mit `-` beginnt (wie [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 6).
    (a) `-n`, `-q`, `-t`: kein Einzelziel läuft, Exit 0, keine Zusammenfassung und keine Meldung, auch bei leerer Liste.
    (b) `-i`: Exit 2, eine Zeile `<sammelziel>: -i …` auf stderr, kein Einzelziel; Grund: ein Einzelziel erbte `-i`, ein rotes Skript
    gäbe Exit 0 zurück, und das Sammelziel meldete Grün. `-i` gilt in jeder Schreibweise (`-i`, `--ignore-errors`), aus der Umgebung
    und zusammen mit anderen Flags (`-k`, `-j`). (c) `-i` zusammen mit `-n`, `-q` oder `-t`: (a) gilt. (d) `-k` und `-j` ändern nichts.
    Fälle, je Flag einer: `p10-n`, `p10-q`, `p10-t` auf der Kommandozeile und `p10-n-env`, `p10-q-env`, `p10-t-env` als `MAKEFLAGS`
    der Umgebung; `p10-i`, `p10-i-env`, `p10-i-lang` (`--ignore-errors`), `p10-i-k`, `p10-i-j`; `p10-in`, `p10-iq`, `p10-it`
    (kein Einzelziel, Exit 0); `p10-k` (ein rotes Ziel, Exit 1, alle Ziele liefen, Zusammenfassung wie ohne `-k`) und `p10-j` (wie
    Punkt 5, unter `-j4`). Mutationen: das Skript überliest `n`; `q`; `t`; `i` (`i` fällt aus der Auswahl); `i` wird als `n` gelesen
    (Exit 0 statt 2); `i` gilt vor `n` (Exit 2 statt 0 bei `p10-in`); das Skript liest das ganze `MAKEFLAGS` statt des ersten Worts; ohne `k` bricht
    das Skript beim ersten roten ab; `j` startet die Einzelziele gleichzeitig.
11. **Kein Stand, kein Nachweis.** Ein Lauf eines Sammelziels, ob grün, rot, abgebrochen oder unter `-n`, `-q`, `-t`, legt keine Datei im
    Arbeitsbaum an und ändert keine, auch nichts unter `.harness/state/`; kein Hook und kein anderes Gate liest sein Ergebnis. Fälle
    `p11-gruen`, `p11-rot`, `p11-abbruch`, `p11-n`, `p11-q`, `p11-t`: die Dateiliste des Arbeitsbaums der Kopie (Pfade und Inhalte) ist nach
    einem grünen Lauf, einem roten, einem abgebrochenen und unter `-n`, `-q`, `-t` unverändert. Mutationen, je Fall eine, die nur ihren
    Fall färbt: das Skript schreibt eine Datei unter `.harness/state/` im grünen Lauf, im roten, beim Abbruch, unter `-n`, `-q`, `-t`.
12. **Reihenfolge der Gründe.** Das Skript prüft in dieser Reihenfolge und endet beim ersten Treffer: (1) `-n`, `-q`, `-t`: Exit 0, nichts; (2) `-i`:
    Exit 2; (3) leere Liste: Exit 2; (4) die Einzelziele. Fälle: `p12-n-leer` (`-n` bei leerer Liste: Exit 0, keine Meldung), `p12-i-leer` (`-i` bei
    leerer Liste: die Meldung nennt `-i`, nicht die Liste); das Paar (1)/(2) deckt `p10-in`. Mutationen: die Prüfung der Liste vor `-n`; die Prüfung der Liste vor
    `-i`.

**Für beide Teile.**

13. **Schnittstelle und Aufruf-Voraussetzung.** Die Schnittstelle sind die beiden Make-Ziele; das Skript wird nicht unmittelbar gerufen, und ein
    Aufruf aus einem Unterverzeichnis (etwa mit `make -C`) ist Sache von `make` (akzeptiertes Negativ; Grund: das Rezept läuft im Verzeichnis des
    `Makefile`). Zugesagt ist nacheinander: zwei Fahrten im selben Arbeitsbaum, eine Fahrt zusammen mit `make gates` und zwei Fahrten in
    zwei Arbeitsbäumen am selben Docker-Daemon sind nicht zugesagt ([`SPEC-052`](SPEC-052-gate-kette.md) Punkt 8); die Einzelziele tragen
    die Präfixe ihrer Namen für Container und Netze, und mehr als diese Trennung zusagt kein Vertrag. Das Skript sagt keine Ordnung über mehrere Aufrufe zu.
14. **Pflicht zum Fahren (Entscheidung des Nutzers vom 2026-10-18).** Regel an den Aufrufer, kein Fall (Grund: ein Gate, das sie erzwänge, wäre das
    Gate, das diese Entscheidung verlässt): Wer ein Gate, eine Gegenprobe oder einen Vertrag eines Gates ändert, fährt die betroffene Gegenprobe;
    am Ende jeder Welle und vor dem Release laufen `make gegenproben` im Gegenproben-Container und `make gegenproben-host` am Host. Welche
    Rolle wann fährt und wo das Ergebnis liegt, sagen die Anweisungen der Rollen, nicht dieser Vertrag. Für eine Gegenprobe im Zustand `nicht gemessen` (Punkt 16) ist der engste Fall mit `GP_MUT=<Muster>` die Fahrt (Entscheidung des Nutzers vom 2026-10-22, Wortlaut „ok“). Der Stop-Hook verlangt keinen Stand der
    Gegenproben, und `make gates` grün sagt nichts über sie (akzeptiertes Negativ; Grund: ein Stand der Gegenproben brächte die Laufzeit zurück,
    die die Entscheidung streicht; die Folge einer nicht gefahrenen Gegenprobe ist ein unbemerkter Bruch des geprüften Werkzeugs bis zur
    nächsten Fahrt, `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`).
15. **Grenzen.** (00) Die Sprache von `make` ist Voraussetzung des Hosts, kein Wert, den ein Aufrufer setzt: Die Regel erzwingt `LC_ALL=C` (Punkt 2); Aufrufe von `make` außerhalb der Regel behalten die Sprache des Hosts (kein Fall; Grund: nur die Regel liest die Ausgabe). (0) Der Kommentar in Zeile 2 des generierten `Makefile` („jedes Fragment haengt seine Checks an `GATE_CHECKS`“) gilt für die Gate-Fragmente; die Fragmente, die eine Gegenprobe definieren, hängen sie an eine der beiden Listen. Die Abweichung bleibt stehen (Grund: ein Lauf von `ai-harness-init` schreibt die Datei neu, und eine Änderung hier ginge verloren); kein Fall. (1) Ein Zeitlimit mitten im Lauf gibt es nicht (die Obergrenze von Punkt 16 wird nach dem Lauf gelesen); `SIGKILL` und der Ausfall des Docker-Dienstes sind nicht zugesagt. (2) Ein
    Signal nur an `make` (nicht an das Skript und nicht an die Prozessgruppe) ist Verhalten von `make` und nicht zugesagt; ein Fall auf den
    Zeitpunkt nach dem Ende von `make` wäre zeitabhängig. (3) Die Warnung von `make` zum Jobserver unter `-j` ist weder zugesagt noch
    ausgeschlossen (Punkt 7). (4) Ob eine Gegenprobe das prüft, was ihr Gate zusagt, ist Gegenstand des Vertrags ihres Gates, nicht dieser
    Kennung. (5) Ob die Zuordnung zu Container oder Host stimmt, prüft nichts (Punkt 3). (6) Namen mit Leerraum oder führendem `-` in einer
    Liste sind nicht zugesagt. (7) Die Zeit der Fahrt ist für aufgenommene Gegenproben nach Punkt 16 zugesagt; für die im Zustand `nicht gemessen` ist sie nicht zugesagt, sie wird gemessen und im Bericht genannt, soweit gelesen, und der Nehmer der Tabelle trägt die Aufnahme.

16. **Zeit und Platz einer Gegenprobe (Entscheidung des Nutzers vom 2026-10-21, wörtlich: „die gegenproben werden zeitlich begrenzt. höchstens 10 minuten“; [ADR-0055](../../docs/plan/adr/0055-gegenproben-haben-ein-zeitbudget.md)).**
    Der Punkt gilt für jede Gegenprobe der Listen `GEGENPROBEN_CONTAINER` und `GEGENPROBEN_HOST` einschließlich `a-check-negativ`, jede nennt in ihrem Kopf die Obergrenze und diesen Punkt. **Aufnahme:** Die Grenze gilt für eine Gegenprobe ab ihrer Aufnahme; aufgenommen ist sie, wenn ein Slice ihren vollen Lauf gemessen, die Fälle von (9) in ihr eingeführt und sie in der Tabelle unten auf `aufgenommen` gesetzt hat. Bis dahin steht sie in der Tabelle im benannten Zustand `nicht gemessen` mit der Kennung des Slice, der sie aufnimmt; ein Eintrag ohne Kennung eines vorhandenen Slice in `open/`, `next/` oder `in-progress/` ist ein Fehler dieses Vertrags, kein Zustand (Grund: ein Zustand ohne Nehmer ist eine Ausnahme ohne Frist). Der Vorschlag des Architect für die Nehmer steht in der Tabelle; die Kennung legt der Planner an (Stand 2026-10-21 noch nicht angelegt).
    | Gegenprobe | Zustand | Aufnahme durch |
    |---|---|---|
    | `release-gegenprobe`, `test-container-gegenprobe` | nicht gemessen (Fälle von (9) fehlen noch) | `slice-harness-gegenproben-schneller` |
    | `lint-gegenprobe`, `stop-hook-gegenprobe`, `gegenproben-nach-einlieferung-gegenprobe`, `ci-gegenprobe`, `coverage-gegenprobe`, `gates-schnell-gegenprobe`, `gegenproben-gegenprobe`, `workflows-check-gegenprobe`, `range-konfiguration-gegenprobe`, `integration-volumes-gegenprobe`, `test-postgres-versionen-gegenprobe`, `archiv-bezug-gegenprobe`, `abdeckung-gegenprobe`, `kopf-check-gegenprobe`, `meldungskatalog-check-gegenprobe`, `hook-gegenprobe`, `a-check-negativ` | nicht gemessen | `slice-harness-gegenproben-grenze-rest` (Vorschlag des Architect, vom Planner anzulegen: misst jede, nimmt die unter der Grenze auf, schneidet die übrigen) |
    | `container-start-gegenprobe` (Liste für den Host) | nicht gemessen | ebenfalls `slice-harness-gegenproben-grenze-rest` (Messung am Host) |
    | `container-einrichten-gegenprobe` | aufgenommen (trägt die Fälle von (9); voller Lauf gemessen bei ihrer Lieferung) | `slice-harness-devcontainer-einrichten` (angelegt mit ihr) |
    Diese Tabelle trägt jeder Slice mit, der eine Gegenprobe aufnimmt oder anlegt; eine neue Gegenprobe steht ab ihrem Anlegen als `aufgenommen`, wenn sie die Fälle von (9) trägt, sonst als `nicht gemessen` mit Nehmer.
    (1) **Obergrenze.** Ein voller Lauf (ohne `GP_MUT`, mit allen Fällen und Mutationen) einer aufgenommenen Gegenprobe dauert höchstens 600 s (zehn Minuten, für jede gleich), gemessen als Wanduhr vom Start des Skripts bis nach der Auswertung aller Mutationen, auf dem Gegenproben-Container
    (Linux) bei ruhigem Host. Ruhig heißt: Die Last des Hosts (`/proc/loadavg`, erster Wert) liegt unter einem Fünftel der Kernzahl; die Zahl
    gilt nicht für einen Host unter Last (`BEO-REPO/lastabhaengige-zaehlung`). (2) **Ausgabe.** Der Lauf endet mit einer Zeile auf stderr
    `<name>: Zeit <T> s von höchstens <B> s, Last <L>` (`T` ganze Sekunden, `L` der erste Wert von `/proc/loadavg`, `-`, wenn die Datei fehlt oder
    unlesbar ist); sie steht bei grünem und rotem Ergebnis und fehlt bei Abbruch durch ein Signal und bei Exit 2. (3) **Überschreitung.** Ist `T`
    größer als `B`, steht davor die Zeile `<name>: ROT — Zeitbudget: <T> s über <B> s`, und der Exit-Code ist 1, auch wenn jede Mutation wie
    erwartet endete; gleich `B` ist grün (nicht zugesagt, siehe (9) letzter Satz). Der Lauf wird bei Erreichen von `B` nicht abgebrochen, denn ein
    Abbruch kürzte das Ergebnis; er endet, wenn er endet. Eine Überschreitung unter Last ist kein Befund über den Code: Wer sie
    sieht, liest `L` und fährt auf einem ruhigen Host neu (Grund: ein Lauf, der die Last erkennt und die Prüfung auslässt, wäre selbst
    ein Gate mit eigener Zusage). Ein Lauf, der auf ruhigem Host die Obergrenze reißt, ist rot, bis eine Mutation schneller wird oder die
    Obergrenze per ADR geändert ist (Punkt 16 (9)). (4) **`GP_BUDGET=<s>`.** Setzt `B` in Sekunden; `0` schaltet die Prüfung ab; ohne
    die Variable ist `B` die Obergrenze aus (1), und nur dann nur bei einem Lauf ohne `GP_MUT` (mit `GP_MUT` gibt es keine Vorgabe, die Zeile
    nach (2) steht ohne `B` als `Zeit <T> s`). Ein Wert, der keine ganze Zahl ohne Vorzeichen ist, endet vor jedem Fall mit Exit 2 und
    einer Zeile `<name>: GP_BUDGET: keine ganze Zahl`. **Reihenfolge der Gründe:** (a) `GP_BUDGET` ungültig (Exit 2), (b) Abbruch durch
    `SIGINT` (130) oder `SIGTERM` (143) ohne Zeile nach (2), (c) rote Mutation oder Fall (Exit 1, Zeilen wie bisher), (d) Zeitbudget (Exit 1);
    (c) und (d) stehen beide in der Ausgabe, wenn beide zutreffen.
    (5) **Arbeitsordner.** Der Ordner einer Mutation (bei `release-gegenprobe` `$arbeit/<name>`, bei `test-container-gegenprobe` `baum-<n>`)
    wird entfernt, sobald ihre Fälle bewertet sind, die rote wie die grüne, und der Ordner jedes Falls, sobald der Fall endet (Grund: ein
    Ordner je Mutation und Fall summierte sich auf über 5 GB je Lauf). Jede Mutation hat einen eigenen Ordner und kein gemeinsames
    Verzeichnis mit einer anderen; mehrere Mutationen zugleich (`GP_PARALLEL`) entfernen je ihren. Die Spitze des Platzbedarfs des
    Arbeitsordners liegt bei der Vorgabe von `GP_PARALLEL` unter 100 MB. Bei `SIGINT`, `SIGTERM` und Ende entfernt das Skript den ganzen
    Arbeitsordner, den der laufenden Mutationen eingeschlossen, nachdem es seine Nachkommen beendet hat (Bestand, Fälle `p13-*`). Ein
    Ordner, der nach der Mutation für die Diagnose stehen bliebe, gibt es nicht: Der Arbeitsordner war schon bisher mit dem Lauf weg.
    (6) **Plattform.** Zeit und Platz sind auf dem Gegenproben-Container (Linux) gemessen; für macOS und Windows gilt die Grenze aus dem Kopf
    dieser Kennung (nicht zugesagt). (7) **Aufrufarten.** `GP_MUT`, `GP_ENTDECKE`, `GP_FALL` und `GP_PARALLEL` behalten ihre Bedeutung; unter
    `make` gelten `-n`, `-q`, `-t`, `-i` und `MAKEFLAGS` wie bisher ohne eigenen Fall ([`SPEC-067`](SPEC-067-release.md) Kopf von Punkt 14);
    `GP_BUDGET` aus der Umgebung von `make` erreicht das Skript. Die Vorgabe von `GP_PARALLEL` darf von 4 auf höchstens 8 steigen, wenn zwei
    volle Läufe in Folge grün sind (Messung im Bericht); Grund für die Grenze: Fälle mit Wartezeiten und `pgrep` fallen unter Last
    aus (`BEO-REPO/test-faellt-unter-cpu-last-aus`). (8) **Umgebung der Fahrt.** Die Messung nach (1) verlangt einen ruhigen Host. Platz,
    Benutzer (root und ohne root), Sprache und Zeitzone ändern die Zahlen nach (1) nicht; die Fahrt am Endstand läuft einmal unter einem
    fremden Benutzer, einer fremden Locale und unter Last und nennt die Zeit dabei, ohne dass die Obergrenze für sie gilt.
    (9) **Ergebnisgleichheit und Lockerung.** Eine Beschleunigung ändert weder die Mutationen noch ihre erwarteten roten Fälle noch die Liste
    der Fälle. Beleg: der Unterschied der Deklarationszeilen (Zeilen, die eine Mutation, ihre roten Fälle oder die Fall-Listen deklarieren) zwischen
    Basis und Endstand ist leer, und ein voller Lauf am Endstand ist grün (die Gegenprobe vergleicht je Mutation die tatsächlichen mit den
    deklarierten roten Fällen). Eine Mutation zu streichen oder zusammenzulegen, einen Fall zu streichen oder die Obergrenze zu erhöhen ist
    eine Lockerung (`AGENTS.md` §3.6) und braucht eine ADR; die Option, die roten Fälle nur unter einer Auswahl der Fälle zu prüfen
    (wie [`SPEC-055`](SPEC-055-container-und-binaries.md) mit `nur`), ist es auch, wo sie die Prüfung der übrigen Fälle je Mutation auslässt.
    **Fälle und Mutationen** (in der jeweiligen Gegenprobe, Name mit der Nummer ihres Vertrags; je Element eines): Überschreitung rot (Kopie mit
    `GP_BUDGET=1` und einer billigen Mutation endet mit Exit 1 und der Zeile nach (3)); Grün unter der Grenze (`GP_BUDGET=3600`, Zeile
    nach (2) ohne `ROT`); `GP_BUDGET=0` ohne Prüfung; ungültiger Wert Exit 2 vor jedem Fall (Leerstring, `-1`, `1.5`, `abc`, je einer);
    Vorgabe nur ohne `GP_MUT` (Lauf mit `GP_MUT` und ohne `GP_BUDGET` hat keine Zeile `ROT — Zeitbudget`); Zeile nach (2) mit `Last` und
    `-` bei fehlender Datei; Reihenfolge der Gründe (c) vor (d); keine Zeile nach (2) bei `SIGINT` und `SIGTERM`; Arbeitsordner einer Mutation
    nach ihrem Ende entfernt (eine Kopie ohne `rm` färbt den Fall rot), je Fall entfernt, und nach Abbruch weg (Bestand). Die Mutation
    zur Nutzung von `>=` statt `>` bleibt grün (äquivalent in der Praxis; akzeptiertes Negativ: die Zeit ist nicht steuerbar, die Wirkung beträgt eine
    Sekunde).

17. **Start mit ignoriertem `SIGINT` (Entscheidung des Architect vom 2026-10-21, nach F-V10 von `slice-harness-devcontainer-einrichten`; Regel an den Aufrufer, kein Fall).**
    Der Punkt gilt für jede Gegenprobe der Listen `GEGENPROBEN_CONTAINER` und `GEGENPROBEN_HOST` einschließlich `a-check-negativ`, die in einem Fall
    `SIGINT` an einen Prozess sendet, den sie selbst gestartet hat, und für das Skript der Sammelziele (Punkt 8). **Voraussetzung:** `SIGINT` ist beim
    Start nicht ignoriert. **Nicht zugesagt** ist der Start mit ignoriertem `SIGINT` (ein Start mit `&` aus einer Shell ohne Jobsteuerung, `nohup`, ein
    Aufrufer, der es ignoriert): Ergebnis, Meldung, Exit-Code und Dauer der Fälle mit `SIGINT` sind dann nicht zugesagt, und ein Rot dieser Fälle ist
    kein Befund über den Code. Akzeptiertes Negativ; Grund: `bash` kann eine beim Start ignorierte Disposition weder fangen noch zurücksetzen (ein
    `trap` auf `INT` bleibt ohne Wirkung, gemessen unten), und jeder Nachkomme erbt sie; ein Zurücksetzen oder Ablehnen am Start wäre je Gegenprobe Code
    mit Fall und Mutation (`AGENTS.md` §3.10) für einen Aufruf, den die Fahrt der Sammelziele nicht erzeugt, denn der Hook setzt `SIGINT` und `SIGQUIT`
    für den Lauf von `make gegenproben` zurück ([`SPEC-061`](SPEC-061-gegenproben-nach-einlieferung.md) Punkt 11, Fälle `p11-erbt-int` und
    `p11-erbt-quit`). Eine Gegenprobe, die den Start mit ignoriertem `SIGINT` dennoch besteht (etwa durch `env --default-signal` an den Startstellen
    ihrer Signalfälle), sagt es nicht zu, solange ihr Vertrag es nicht mit einem Fall nennt.
    **Regel an den Aufrufer:** im Vordergrund starten, oder im Hintergrund mit zurückgesetzter Disposition: `setsid -f env --default-signal=INT,QUIT
    <Befehl> > <Datei> 2>&1`, der Befehl ist `make <Ziel>` oder `bash tools/dev/<bereich>/<name>.sh` (`<bereich>` ist `harness`, `test`, `hook` oder `arch`; `env --default-signal` braucht GNU coreutils ab 8.31,
    wie [`SPEC-061`](SPEC-061-gegenproben-nach-einlieferung.md) Punkt 11). **Gemessen** am 2026-10-21 (bash 5.2.37, coreutils 9.7) mit einem Aufrufer,
    der `SIGINT` mit `trap '' INT` ignoriert: Bei direktem Start, bei `setsid -f` allein und bei `nohup` bleibt `SIGINT` im Skript ignoriert (`trap -p INT`
    zeigt `trap -- '' SIGINT`, die Maske `SigIgn` in `/proc/<pid>/status` trägt das Bit von `SIGINT`, `0x2`, und ein `kill -INT` an das Skript ändert
    nichts); bei `setsid -f env --default-signal=INT,QUIT bash …` ist es nicht ignoriert, und der `trap` greift. `setsid -f` allein ist deshalb **kein**
    Rezept. Wer den Verdacht hat, liest die Maske des gestarteten Prozesses. **Plattform** (`AGENTS.md` §3.12): geprüft ist Linux (Dev-Container und
    Gegenproben-Container); macOS und Windows sind nicht zugesagt (Kopf dieser Kennung), das Rezept nennt eine Option von GNU coreutils, die dort nicht
    gemessen ist. Kein Fall (Grund: Der Punkt sagt nichts zu, was ein Test prüfen müsste; die Messung steht hier, die Zusage des Hooks prüft
    [`SPEC-061`](SPEC-061-gegenproben-nach-einlieferung.md)).
