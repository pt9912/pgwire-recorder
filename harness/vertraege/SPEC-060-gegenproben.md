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
Die Punkte 13 bis 15 gelten für beide Teile und haben keinen eigenen Fall, außer wo sie einen nennen. Der
Name jedes Falls trägt die Nummer des Punktes, den er prüft. Je Zusage und je Element einer Aufzählung
gibt es einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die Gegenprobe verlangt
genau diese Menge roter Fälle je Mutation (mehr und weniger machen sie rot).

**Plattform.** Zugesagt ist ein Linux-Host oder -Container mit bash ≥ 4, GNU `make` und GNU-Werkzeugen;
die Gegenprobe der Sammelziele braucht zusätzlich `setsid`. macOS- und Windows-Hosts sind nicht zugesagt
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

4. **Zwei Ziele, ein Skript.** `make gegenproben` fährt `GEGENPROBEN_CONTAINER`, `make gegenproben-host` fährt
   `GEGENPROBEN_HOST`, beide über dasselbe Skript `tools/harness/gegenproben.sh` aus dem Fragment
   `harness/mk/gegenproben.mk`. Beide sind `.PHONY`, kein Gate, nicht in `GATE_CHECKS`, nicht in der Ausnahmeliste, hängen
   nicht voneinander ab, und `make gates` und `make gates-schnell` fahren keines von beiden. Ein Ziel fährt die Liste des
   anderen nie. Fälle: `p4-container` (nur die Einzelziele der Liste für den Container laufen), `p4-host` (nur die der
   Liste für den Host), `p4-kein-gate` (`make -n gates` nennt keines der beiden Sammelziele, und beide stehen nicht in `GATE_CHECKS`). Mutationen: das
   Ziel für den Container fährt auch die Liste für den Host; das Ziel für den Host fährt auch die Liste für den
   Container; das eine hängt an dem anderen; ein Sammelziel steht in `GATE_CHECKS`.
5. **Nacheinander, in der Reihenfolge der Liste.** Jedes Einzelziel beginnt nach dem Ende des vorigen, in der
   Reihenfolge der Liste, auch unter `-j` (eine Schleife im Skript, kein Fächer von Prerequisites). Fälle: `p5-reihenfolge`
   (die Attrappen schreiben `beginn`/`ende` und warten kurz; die Folge im Protokoll ist `beginn a, ende a, beginn b, …` in
   der Reihenfolge der Liste), `p5-j` (dasselbe unter `-j4`). Mutationen: Einzelziele im Hintergrund gestartet; die Liste
   rückwärts gefahren; die Einzelziele als Prerequisites des Sammelziels statt in der Schleife.
6. **Rotes Einzelziel, Ausgang, Zusammenfassung.** Alle Einzelziele der Liste laufen, auch nach einem roten;
   ein Einzelziel ist rot bei jedem Exit-Code ungleich 0, auch 127 (Werkzeug fehlt) und 2 (kein Ziel); das Sammelziel
   meldet das fehlende Werkzeug nicht selbst und prüft vorab nichts. Der Exit-Code des Sammelziels ist 0, wenn alle
   grün sind, sonst 1. Nach der Ausgabe des letzten Einzelziels steht auf stderr je Einzelziel in der Reihenfolge der Liste eine
   Zeile `<sammelziel>: <ziel>: grün` oder `<sammelziel>: <ziel>: rot (Exit <n>)`; sonst schreibt das
   Skript nichts. Fälle: `p6-alle-gruen` (Exit 0, eine grüne Zeile je Ziel), `p6-erstes-rot`, `p6-mittleres-rot`,
   `p6-letztes-rot` (Exit 1, die übrigen laufen und stehen in der Zusammenfassung), `p6-zwei-rot` (beide roten stehen mit
   Namen und Exit-Code dort), `p6-127` (ein Ziel endet mit 127: rot mit `Exit 127`, das nächste läuft). Mutationen: Abbruch
   beim ersten roten; Exit 0 trotz rotem; Exit-Code nur des letzten Einzelziels; Zusammenfassung nur der roten;
   Zusammenfassung ohne Exit-Code; Zusammenfassung vor der Ausgabe des letzten Ziels.
7. **Ausgabe.** Die Ausgabe jedes Einzelziels (stdout und stderr) erscheint unverändert und auf demselben Strom;
   das Skript schreibt außer der Zusammenfassung und den Meldungen der Punkte 9, 10 und 12 nichts. Zeilen von `make`
   selbst (Verzeichniswechsel, Warnungen zum Jobserver unter `-j`) sind nicht zugesagt. Fall `p7-ausgabe`: die
   Attrappe schreibt eine Zeile auf stdout und eine auf stderr; beide erscheinen auf ihrem Strom. Mutationen: stderr nach
   stdout geleitet; stdout verworfen.
8. **Abbruch.** Ein `SIGINT` oder `SIGTERM` beendet das laufende Einzelziel und startet kein weiteres; das Skript endet
   mit 130 (`SIGINT`) oder 143 (`SIGTERM`) und schreibt keine Zusammenfassung. Die Zusage gilt für ein Signal an die
   Prozessgruppe und für ein Signal nur an das Skript. Fälle: je Signal und Adressat einer, vier insgesamt
   (`p8-sigint-gruppe`, `p8-sigterm-gruppe`, `p8-sigint-skript`, `p8-sigterm-skript`); die Attrappe des mittleren von drei Zielen
   wartet, das dritte schreibt bei Beginn in das Protokoll; geprüft werden Exit-Code, dass die Attrappe beendet ist, dass das
   dritte Ziel nicht begann und dass keine Zusammenfassung erscheint. Mutationen: das nächste Ziel startet nach dem Signal;
   die Exit-Codes 130 und 143 vertauscht; das laufende Einzelziel bleibt am Leben (sichtbar nur im Fall mit dem Signal an das Skript); eine Zusammenfassung erscheint.
   `SIGKILL` und ein Ausfall von Docker sind nicht zugesagt (Punkt 15).
9. **Leere Liste.** Eine leere Liste ist ein Fehler, kein stilles Grün: Exit 2, eine Zeile `<sammelziel>: Liste leer` auf stderr, kein
   Einzelziel. Fälle `p9-container` und `p9-host`, je Sammelziel einer (das Fragment der Kopie hängt nichts an die Liste). Mutationen:
   Exit 0; die Prüfung liest die Liste des anderen Sammelziels.
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
    Arbeitsbaum an und ändert keine, auch nichts unter `.harness/state/`; kein Hook und kein anderes Gate liest sein Ergebnis. Fall
    `p11-kein-stand`: die Dateiliste des Arbeitsbaums der Kopie (Pfade und Inhalte) ist nach einem grünen, einem roten Lauf und unter
    `-n` unverändert. Mutation: das Skript schreibt nach grünem Lauf eine Datei unter `.harness/state/`.
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
    Rolle wann fährt und wo das Ergebnis liegt, sagen die Anweisungen der Rollen, nicht dieser Vertrag. Der Stop-Hook verlangt keinen Stand der
    Gegenproben, und `make gates` grün sagt nichts über sie (akzeptiertes Negativ; Grund: ein Stand der Gegenproben brächte die Laufzeit zurück,
    die die Entscheidung streicht; die Folge einer nicht gefahrenen Gegenprobe ist ein unbemerkter Bruch des geprüften Werkzeugs bis zur
    nächsten Fahrt, `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`).
15. **Grenzen.** (00) Die Sprache von `make` ist Voraussetzung des Hosts, kein Wert, den ein Aufrufer setzt: Die Regel erzwingt `LC_ALL=C` (Punkt 2); Aufrufe von `make` außerhalb der Regel behalten die Sprache des Hosts (kein Fall; Grund: nur die Regel liest die Ausgabe). (0) Der Kommentar in Zeile 2 des generierten `Makefile` („jedes Fragment haengt seine Checks an `GATE_CHECKS`“) gilt für die Gate-Fragmente; die Fragmente, die eine Gegenprobe definieren, hängen sie an eine der beiden Listen. Die Abweichung bleibt stehen (Grund: ein Lauf von `ai-harness-init` schreibt die Datei neu, und eine Änderung hier ginge verloren); kein Fall. (1) Ein Zeitlimit des ganzen Laufs gibt es nicht; `SIGKILL` und der Ausfall des Docker-Dienstes sind nicht zugesagt. (2) Ein
    Signal nur an `make` (nicht an das Skript und nicht an die Prozessgruppe) ist Verhalten von `make` und nicht zugesagt; ein Fall auf den
    Zeitpunkt nach dem Ende von `make` wäre zeitabhängig. (3) Die Warnung von `make` zum Jobserver unter `-j` ist weder zugesagt noch
    ausgeschlossen (Punkt 7). (4) Ob eine Gegenprobe das prüft, was ihr Gate zusagt, ist Gegenstand des Vertrags ihres Gates, nicht dieser
    Kennung. (5) Ob die Zuordnung zu Container oder Host stimmt, prüft nichts (Punkt 3). (6) Namen mit Leerraum oder führendem `-` in einer
    Liste sind nicht zugesagt. (7) Die Zeit der Fahrt ist nicht zugesagt; sie wird gemessen und im Bericht genannt, soweit gelesen.
