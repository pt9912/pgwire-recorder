# SPEC-052 — Gate-Kette und Wahl der Kette (`gates-schnell`)

`make gates` ist die volle Kette der Gates, `make gates-schnell` ihre schnelle Hälfte, und
`make gates-schnell-gegenprobe` prüft die schnelle Kette. Die Kennung legt fest, welche Gates
die beiden Ketten fahren, welchen Stand die schnelle Kette hinterlässt und, im zweiten Teil,
nach welcher Regel der Stop-Hook `.claude/hooks/stop-require-gates.sh` zwischen den Ketten
wählt. Alle Zusagen dieser Kennung gelten unter der Aufruf-Voraussetzung von Punkt 8.

**Teil 1: die Ketten.** Die Punkte 1 bis 7 prüft `make gates-schnell-gegenprobe`; Punkt 8 ist
eine Bindung der Zusagen und hat keinen Fall.

1. **Volle Kette.** Die volle Kette ist `GATE_CHECKS` selbst, ohne zweite Liste: Jedes Gate,
   das ein Fragment unter `harness/mk/` an `GATE_CHECKS` hängt, fährt `make gates`. Ein Gate
   kann aus ihr nicht still fallen, weil kein Fragment sie aufzählt. `make gates` schreibt
   den Nachweis des vollen Laufs `.harness/state/gates-passed.diffsha`, dazu die Dateiliste
   (Punkt 11); bei den Make-Flags `-i`, `-n`, `-q`, `-t` gilt Punkt 21.
2. **Schnelle Kette.** Die schnelle Kette ist `GATE_CHECKS` ohne die Ausnahmeliste. Die
   Ausnahmeliste steht im Fragment der Kette, `harness/mk/gates-schnell.mk`, und nennt jedes
   Gate, dessen Name auf `-gegenprobe` endet, dazu `a-check-negativ`, `test-container`, `test-integration` und `test-postgres-versionen`. `build`,
   `test` und alle übrigen Gates fährt die schnelle Kette. Sie fährt kein Gate der
   Ausnahmeliste.
3. **Neues Gate.** Ein neues Gate läuft in beiden Ketten, bis es ausdrücklich in die
   Ausnahmeliste aufgenommen wird; ein Fehler der Liste kostet Zeit, nie ein Gate der vollen
   Kette: Ein fehlender Eintrag lässt ein langsames Gate in der schnellen Kette laufen, ein
   Eintrag zu viel nimmt ein Gate aus der schnellen Kette, die volle fährt es weiter.
4. **Deckung der Ausnahmeliste.** Jeder Eintrag der Ausnahmeliste steht in `GATE_CHECKS`,
   und jedes Gate in `GATE_CHECKS`, dessen Name auf `-gegenprobe` endet, sowie `a-check-negativ`,
   `test-container`, `test-integration` und `test-postgres-versionen` steht in der Ausnahmeliste. Ein neues Gate dieser Namen macht die
   Gegenprobe der Kette rot, bis jemand entscheidet.
5. **Stand der schnellen Kette.** Bei einem grünen Lauf der schnellen Kette schreibt sie den
   Hash des ganzen Arbeitsbaums zum Ende des Laufs, die unveränderte Ausgabe von
   `tools/harness/working-tree-hash.sh`, in `.harness/state/gates-schnell-passed.diffsha`.
   Ändert ein Gate des Laufs den Arbeitsbaum, gilt der Hash nach dem Lauf.
6. **Roter oder abgebrochener Lauf.** Nach einem roten oder abgebrochenen Lauf der schnellen
   Kette existiert kein Stand; ein älterer entfällt. Zugesagt ist das Ergebnis, nicht der
   Mechanismus. Ein Stand entsteht nur nach einem Lauf, in dem die Gates tatsächlich liefen:
   Ein Aufruf mit dem Make-Flag `-n`, `-q` oder `-t` (das Rezept fährt kein Gate oder keines
   wirklich) schreibt keinen Stand und entfernt keinen; ein Aufruf mit `-i` (Fehler werden
   überlesen) schreibt keinen, ein älterer entfällt. Maßgeblich sind die Flags in `MAKEFLAGS`,
   auch aus der Umgebung. Stehen `-n`, `-q` oder `-t` zusammen mit `-i`, gilt `-n`, `-q`,
   `-t`: der Lauf fährt kein Gate, ein vorhandener Stand bleibt unberührt.
7. **Der volle Nachweis bleibt getrennt.** Die schnelle Kette schreibt, verändert und löscht
   den Nachweis des vollen Laufs nie und ersetzt ihn nie.
8. **Aufruf-Voraussetzung.** Zugesagt ist nacheinander: Änderung des Arbeitsbaums, Lauf der
   vollen Kette, Lauf der schnellen Kette und Hook kommen nicht gleichzeitig. Gleichzeitige
   Läufe sind nicht zugesagt. Ein Hook, der eine halb geschriebene Datei liest, endet in „kein
   Stand“ oder „beschädigt“ und damit bei der vollen Kette; mehr wird nicht zugesagt, ein
   atomares Schreiben und ein Fall dafür entfallen.

**Teil 2: der Hook.** Die Zusagen dieses Teils gelten, soweit das Gate
`stop-hook-gegenprobe` sie prüft, und nur so weit: Ein Punkt, den kein Fall der Gegenprobe
nennt, ist nicht zugesagt. Führt `make help` kein Ziel dieses Namens, ist der Hook-Teil
insgesamt nicht zugesagt.

Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes, den er prüft. Die Punkte
9 bis 22 gelten nach dieser Regel und für den Hook, nicht für die Gegenprobe der Ketten.

9. **Wahl der Kette.** `make gates` ist der einzige Nachweis für Code, Gates, Profile,
   Skripte, `Dockerfile`, Fragmente, `go.mod` und alles, was Punkt 12 nicht freigibt. Der
   Hook vergleicht den Arbeitsbaum mit dem Stand des letzten vollen grünen Laufs. Liegt jede
   geänderte, hinzugekommene oder gelöschte Datei auf der Positivliste (Punkt 12) und ist der
   Stand der schnellen Kette grün, gibt er frei. Jeder andere Fall verlangt `make gates`.
   Fehlt der Stand der schnellen Kette, obwohl jede geänderte Datei auf der Positivliste
   liegt, nennt der `reason` `make gates-schnell` (es genügt; `make gates` genügt ebenso).
   Fehlt er bei einer nicht inerten Änderung, nennt der `reason` nur `make gates`, nie
   `make gates-schnell` (damit ist auch die Reihenfolge „Pfadregel, zuletzt Stand“
   beobachtbar). Die Gegenprobe prüft den `reason` in diesen beiden Fällen. Für alle übrigen
   Sperren, auch für einen Stand eines anderen Arbeitsbaums, einen leeren und einen nicht
   lesbaren Stand, ist der Wortlaut des `reason` nicht zugesagt.
    Reihenfolge der Prüfungen: Schleifenschutz (Punkt 20), nicht aufgelöster Index-Eintrag
    (Punkt 17), fehlender oder beschädigter Stand des vollen Laufs (Punkte 19, 20),
    Identität der Listen (Punkt 16), danach die Pfadregel (Punkte 12, 14, 15, 18), zuletzt der
    Stand der schnellen Kette (Punkt 10). Die Positivliste ersetzt keine der früheren Prüfungen.
10. **Bezugsgröße des Stands der schnellen Kette.** Der Hook vergleicht den Stand der
    schnellen Kette (Punkt 5) mit dem Hash des aktuellen Arbeitsbaums; ungleich, fehlend oder
    nicht lesbar heißt: kein Stand. `tools/harness/working-tree-hash.sh` bleibt eine Ausgabe,
    die genau dann wechselt, wenn sich der Arbeitsbaum ändert. Der Hook ruft dieses Skript und
    baut den Hash nicht nach; fehlt das Skript oder scheitert es, gilt Punkt 20 (`block`,
    beim Schleifenschutz `approve`). Verglichen wird der Hash, nicht die Bytes der Datei:
    ein angehängter Zeilenumbruch ändert nichts (akzeptiertes Negativ; Grund: Zugesagt ist nur
    der Schreiber der schnellen Kette (Punkt 5), und eine fremde Datei mit dem passenden Hash
    gibt nur frei, wenn der Hash des Arbeitsbaums stimmt, nie öfter). Der Hook ruft kein
    `make`; die Make-Flags `-n`, `-q`, `-t`, `-i` und `MAKEFLAGS` wirken auf den Stand nur
    beim Schreiben (Punkt 6), und der Hook stützt sich darauf, ohne den Ursprung der Datei zu
    prüfen: Ein Lauf mit `-i` schreibt keinen Stand und löscht einen älteren, also kann kein
    Stand eines solchen Laufs vorliegen; den Fall dafür führt die Gegenprobe der Ketten
    (Punkt 6), nicht die des Hooks.
11. **Form des Stands des vollen Laufs.** Der Stand des vollen Laufs ist eine Dateiliste: eine
    Datei aus NUL-getrennten Feldern, drei je Eintrag: Art, Wert, Pfad. Art `F` (reguläre
    Datei; Wert = SHA-256 des Inhalts, 64 Hexzeichen), `L` (Symlink; Wert = `readlink`-Ziel
    unverändert), `G` (im Index, im Arbeitsbaum fehlend; Wert leer), `X` (anderer Typ:
    Verzeichnis, Gitlink; Wert leer). Der Pfad ist wie `git ls-files -z` ihn liefert (getrackt
    und untracked, ohne ignorierte), die Sortierung bytegenau, ein Pfad höchstens einmal.
    Zeilenumbruch, Leerzeichen, Tabulator und Nicht-UTF-8 im Pfad sind damit ohne Sonderfall.
    Die Erzeugung setzt bash ab Version 4 und GNU-Werkzeuge (`readlink -z`, `sort -z`) voraus,
    wie das Repo GNU `make` voraussetzt; ein Fall mit einer anderen Umgebung gibt es nicht,
    die Voraussetzung ist nicht zugesagt. Fehlt sie, bricht der Lauf ab, und Punkt 20 gilt.
    Sortierung und Einmaligkeit sind Form, keine Zusage mit eigenem Fall (akzeptiertes
    Negativ; Grund: Schreiber und Leser erzeugen die Liste mit derselben Funktion, eine
    abweichende Sortierung macht aus „identisch“ höchstens „voll“, nie „frei“). Der
    Ausführungs-Modus gehört nicht zum Eintrag (akzeptiertes Negativ: wie der Hash heute; eine
    Modusänderung ohne Inhaltsänderung bleibt ungeprüft).
12. **Positivliste, Grenze der Pfade, Groß- und Kleinschreibung.** Inert ist ein Pfad, der mit
    `docs/plan/planning/`, `docs/reviews/`, `docs/user/` oder `docs/developer/` beginnt, einschließlich des
    Schrägstrichs, oder genau `README.md` ist (Wurzel), bytegenau und case-sensitiv.
    `docs/plan/planning-x/…`, `docs/reviews-alt/…`, `docs/users/…`, `docs/developer-alt/…`,
    `docs/developers/…`, `docs/Developer/…`, `docs/Plan/planning/…`,
    eine Datei namens `docs/reviews`, `docs/user` oder `docs/developer`, `readme.md`, `README.md.bak` und
    `harness/README.md` sind nicht inert (`harness/README.md` bleibt draußen: der Beleg von
    Punkt 13 gilt nur für die Wurzel-`README.md`; wer sie aufnehmen will, belegt sie). Ein
    Pfad mit Segment `.`, `..` oder leerem Segment (auch am Ende: der Pfad endet auf `/`), ein
    absoluter Pfad oder ein Pfad mit NUL
    gilt als beschädigte Liste (Punkt 20), nie als inert. Geprüft werden die Pfade, in denen
    sich Stand des vollen Laufs und Arbeitsbaum unterscheiden, beider Listen; ein solcher
    Pfad im Stand gibt nie frei. Für die Seite des Arbeitsbaums gibt es keinen herstellbaren
    Fall, weil `git ls-files` ihn nicht liefert (akzeptiertes Negativ). Ein Pfad mit NUL ist
    in der Form (Punkt 11) nicht darstellbar und endet als Feldzahl-Fehler (Punkt 20), ohne
    eigenen Fall dieses Punktes. Ein geänderter Eintrag der Art `X` (Verzeichnis, Gitlink)
    ist nie inert, auch unter einem Pfad der Positivliste: sein Inhalt steht nicht in der
    Liste. Die Positivliste steht an einer Stelle, im Hook
    `.claude/hooks/stop-require-gates.sh`, und dieser Vertrag nennt sie; eine Änderung an
    Hook, `tools/harness/dateiliste.sh`, `tools/harness/working-tree-hash.sh` oder diesem
    Vertrag (`harness/vertraege/SPEC-052-gate-kette.md`) ist nicht inert und verlangt den vollen Lauf. Akzeptiertes Negativ: Der
    geänderte Hook entscheidet über seine eigene Änderung (er verlangt den vollen Lauf, wenn
    er sich nicht selbst freigibt); eine Abweichung zwischen der Liste im Hook und dieser
    Stelle fängt nur der Fall je Pfadklasse der Gegenprobe, kein Vergleich.
13. **Beleg für die Positivliste (`docs/plan/planning/`, `docs/reviews/`, `docs/user/`,
    `docs/developer/`, `README.md`); Abdeckungstabellen.** Alle fünf sind ohne Ausnahme aufgenommen. Kein Gate der
    Ausnahmeliste (Punkt 2) liest eine Datei unter `docs/plan/planning/`, `docs/reviews/`,
    `docs/user/` oder `docs/developer/` oder die `README.md`: die Docker-Gates (`test-integration`,
    `test-postgres-versionen` (dieselbe Stufe des `Dockerfile`), `test-container` (baut mit dem
    Build-Kontext und liest die Fragmente `harness/mk/container.mk` und `harness/mk/integration.mk`),
    `lint-gegenprobe` über `docker build` einer Kopie des Arbeitsbaums, ebenso `build`, `test`,
    `lint`) sehen nur den Build-Kontext, und der ist eine Allowlist, die `docs/` und
    `README.md` nicht freigibt (`.dockerignore`: `*`, dann nur `go.mod`, `go.sum`, `cmd/`,
    `internal/`, `test/`, `.golangci.yml`, `tools/harness/lint.sh`); ein Go-Test, der dort
    läse, wäre im Container rot. `a-check-negativ` kopiert den Baum, aber `a-check` liest nach
    `.a-check.yml` nur `**/*.go`. Jedes Gate der Ausnahmeliste, dessen Name auf `-gegenprobe` endet, außer
    `lint-gegenprobe` (oben), fährt in Temp-Bäumen, die es selbst anlegt, und liest
    `docs/user/`, `docs/developer/` und die `README.md` des Arbeitsbaums nicht (`abdeckung-gegenprobe` baut ihr
    eigenes `docs/user/`; die Temp-Repos von `gates-schnell-gegenprobe` und
    `stop-hook-gegenprobe` tragen nur Makefile, Fragmente, Skripte und Hook; `integration-volumes-gegenprobe` legt Kopien von
    Dateien unter `tools/test/` in einem Temp-Verzeichnis an und spricht Docker an, [`SPEC-057`](SPEC-057-integration-aufraeumen.md) Punkt 10). Den Planungsbaum
    liest von den Gates der Ausnahmeliste keines aus dem Arbeitsbaum: `kopf-check-gegenprobe`
    baut ihren eigenen, und `hook-gegenprobe` legt die Slices und Wellen ihrer Fälle in einem
    eigenen Temp-Repo an (Regel an den Schreiber dieses Gates: ein Fall bezieht Slice und Welle
    nie aus dem Index des Arbeitsbaums; sonst kann eine Lifecycle-Bewegung unter
    `docs/plan/planning/` das Gate rot machen, ohne dass die schnelle Kette es fährt).
    `docs/plan/planning/` lesen von der schnellen Kette `docs-check` und `kopf-check`;
    `docs/reviews/` und `docs/developer/` liest `docs-check` (`.d-check.yml`: Wurzel `.`, `docs/developer/` steht nicht
    in `scan.ignore`; Stand 2026-10-17: kein Gate der Ausnahmeliste nennt den Pfad, ein Gate der schnellen Kette außer
    `docs-check` auch nicht; ob Schritte des Dokuments gegen den Container laufen, prüft kein Gate). Ein künftiges Gate
    dieser Art, das beide liest, ist eine Änderung an einem Gate: nicht inert, voller Lauf. Die Tabellen `docs/user/abdeckung-*.md` liest nur `abdeckung-check`
    (`tools/test/abdeckung.sh --check`), und das läuft in der schnellen Kette (kein Name auf
    `-gegenprobe`): Eine Änderung an ihnen ohne Änderung an Tests (die Deklarationen stehen in
    `*_test.go`, nicht inert) ist von Hand möglich und wird von `abdeckung-check` rot gemeldet,
    sodass die schnelle Kette sie fängt; deshalb nicht ausgenommen, kein Rückfall auf die
    volle Kette. Das Handbuch und die `README.md` liest von den Gates der schnellen Kette
    `docs-check` (Verweise, Anker; `.d-check.yml`, Wurzel `.`) und `meldungskatalog-check` (die Codes der
    Fehlercode- und Warnungstabellen des Handbuchs, `SPEC-054`). Die Aufnahme gilt, solange
    (a) `docs-check`, `abdeckung-check` und `meldungskatalog-check` in der schnellen Kette stehen (die Gegenprobe der
    Kette hält die Ausnahmeliste, Punkt 4) und (b) der Build-Kontext `docs/` und `README.md`
    nicht freigibt; beides sind Änderungen an Gates oder `.dockerignore`, also nicht inert und
    im vollen Lauf. Die schnelle Kette fährt für eine Handbuch-Änderung kein Gate weniger über
    das Handbuch als die volle (kein langsames Gate liest es), aber sie liefert nicht den
    Nachweis „`make gates` grün über den ganzen Baum“, und ob ein Beispiel im Handbuch gegen
    das gebaute Binary läuft, prüft kein Gate, sondern die Probe von Implementer und Verifier;
    sie bleibt Pflicht der Rollen und wird durch die schnelle Kette nicht ersetzt. Akzeptiertes
    Negativ (Grund: gleiche Folge wie beim akzeptierten Negativ der Positivliste): Ein
    künftiger Leser von `docs/plan/planning/`, `docs/reviews/`, `docs/user/`, `docs/developer/` oder `README.md`
    in einem Host-Skript eines langsamen
    Gates ändert Skript oder Gate und löst den vollen Lauf aus; der Beleg ist dann beim Review
    dieser Änderung neu zu führen, eine eigene Sperre dafür entfällt.
14. **Gelöschte oder umbenannte Datei.** Verglichen wird je Pfad die Menge beider Listen
    (Stand des vollen Laufs, aktueller Arbeitsbaum): ein Pfad, dessen Eintrag fehlt, neu ist
    oder abweicht, gilt als geändert. Eine Umbenennung ist damit ein gelöschter und ein
    hinzugekommener Pfad, beide müssen inert sein.
15. **Symlink.** Inert ist ein geänderter Link nur, wenn Link-Pfad und Ziel inert sind; das
    Ziel wird aus dem Eintrag `L` lexikalisch gegen das Verzeichnis des Links aufgelöst (ohne
    das Dateisystem zu befragen) und muss danach innerhalb des Repos auf der Positivliste
    liegen; ein Ziel außerhalb des Repos, ein absolutes Ziel oder ein nicht lexikalisch
    auflösbares ist nicht inert. Gespiegelt: Liegt irgendwo im Arbeitsbaum ein Link außerhalb
    der Positivliste, dessen Ziel auf sie zeigt, liest ein Gate über ihn inerte Pfade; dann
    gilt jede Änderung als nicht inert. Der Bestand hat Links nur unter `.claude/rules/` mit
    Zielen außerhalb der Positivliste, keinen auf `docs/developer` (Stand 2026-10-17). Einzelheiten: Geprüft wird je geändertem Link jeder
    Eintrag `L`, der ihn trägt (Stand und Arbeitsbaum, auch bei einem Wechsel der Art, etwa
    Datei zu Link). Ziel „auf der Positivliste“ heißt für die Inertheit: ein Pfad unter einem
    Verzeichnis-Eintrag mit Schrägstrich oder genau `README.md`; ein Ziel, das ein Verzeichnis
    der Liste selbst benennt (`docs/user`, `docs/user/`, `docs/developer`, `docs/developer/`), ist nicht inert. Für die Spiegelregel
    zählt weiter: Ein Ziel, das auf der Liste liegt, ein Verzeichnis der Liste benennt oder
    ein Vorfahr davon ist (`docs`, `docs/plan`, die Wurzel `.`), spiegelt; maßgeblich ist die
    Liste des Arbeitsbaums. Ein absolutes Ziel unter der Wurzel des Repos wird gegen sie
    aufgelöst und zählt dann für die Spiegelregel wie ein relatives; für die Inertheit ist es
    nicht inert. Ein Ziel, das über die Wurzel hinaus auflöst, ist außerhalb und spiegelt
    nicht.
16. **Leerer Diff gegen den vollen Lauf.** Ist die Liste des Arbeitsbaums mit dem Stand des
    vollen Laufs identisch und gibt es keinen nicht aufgelösten Index-Eintrag (Punkt 17), gibt
    der Hook frei, ohne Stand der schnellen Kette: der volle Nachweis deckt den Baum, die
    schnelle Kette ist eine Teilmenge.
17. **Merge-Konflikt.** Ein nicht aufgelöster Index-Eintrag (`git ls-files -u` nicht leer) ist
    nie inert: voll. Die Prüfung steht vor jeder Entscheidung, auch vor „identisch“: Die Liste
    trägt die Index-Stufen nicht, ein Konflikt kann bei identischer Liste bestehen, und ein
    Hook, der die Identität zuerst prüfte, gäbe ihn frei. Sie ist die erste Prüfung des Hooks
    nach dem Schleifenschutz (`stop_hook_active`, Punkt 20). Konfliktmarker im Inhalt einer inerten Datei
    ändern nichts, ihr Pfad entscheidet.
18. **Untracked Datei.** Eine untracked Datei zählt wie eine getrackte
    (`git ls-files --others --exclude-standard`): in der Positivliste inert, sonst nicht. Eine
    durch `.gitignore` ausgeblendete Datei ist nicht Teil der Liste (Bestand, bewusst stehen
    gelassen: ein Gate, das sie liest, hat sie auch heute nicht im Hash).
19. **Erster Lauf ohne Stand des vollen Laufs.** Ein fehlender Stand des vollen Laufs
    verlangt `make gates`, auch bei sauberem Baum. Ein frischer Klon oder Worktree
    (`.harness/state/` ist ignoriert) sperrt damit den ersten Stop, bis `make gates` lief.
    Der Kopf-Kommentar des Hooks sagt diese Sperre zu.
20. **Nicht lesbarer Vergleich; Abbruch des Hooks; beschädigte Dateiliste.** Scheitert `git`,
    ist eine Datei nicht lesbar oder bricht der Hook mitten im Lauf ab, gilt: voll, und zwar
    als Antwort `block`, nicht als Abbruch ohne Antwort. Davon ausgenommen ist der
    Schleifenschutz: Ist `stop_hook_active` wahr, antwortet der Hook `approve`, auch bei einem
    Konflikt (Punkt 17), bei fehlendem Stand (Punkt 19), wenn `git` scheitert und wenn
    `tools/harness/dateiliste.sh` fehlt; bei `false` und bei fehlendem Feld gilt `block`.
    Akzeptierte Grenze: Das Fail-open ist gewollt, damit ein dauerhaft rotes Gate keine
    Endlosschleife erzeugt, und es überstimmt auch die Sperre durch Konflikt und fehlenden
    Stand. Laufzeit: Der Hook sagt keine Dauer zu (im Bestand einige Sekunden, die Liste wird
    mehrfach berechnet). Der Eintrag des Treibers setzt kein `timeout`; was der Treiber bei
    Überschreiten seiner eigenen Frist tut, ist nicht zugesagt und hat keinen Fall (akzeptiertes
    Negativ; Grund: wie bei der Umgebung des Treibers prüft die Gegenprobe nur den Hook, nicht
    den Treiber). Ist die Dateiliste des vollen Laufs gelöscht
    oder beschädigt, gilt ebenfalls: voll. Beschädigt ist: leere Datei, Feldzahl nicht durch 3
    teilbar, unbekannte Art, `F` mit Wert, der nicht 64 Hexzeichen ist, `G` oder `X` mit Wert,
    doppelter Pfad, Pfad nach Punkt 12 unzulässig. Ein unzulässiger Pfad im Stand sperrt in
    jedem Fall, weil er nie mit der Liste des Arbeitsbaums identisch ist; ob ein Hook mit
    Positivliste ihn als Beschädigung erkennt, bevor er ihn liest, prüft ein Fall dieses
    Hooks, nicht der Hook ohne Positivliste. Eine am Eintragsende abgeschnittene Liste
    ist nicht beschädigt, aber unschädlich: die fehlenden Pfade gelten als hinzugekommen und
    entscheiden nach ihrem Pfad.
21. **Make-Flags beim Schreiben des Nachweises des vollen Laufs.** Der Hook glaubt den Nachweis,
    also schreibt ihn `make gates` nur nach einem Lauf, in dem die Gates tatsächlich liefen und
    kein Fehler überlesen wurde. Ein Aufruf mit dem Make-Flag `-i` (`--ignore-errors`, in
    `MAKEFLAGS` gleich welcher Schreibweise, auch aus der Umgebung und zusammen mit anderen
    Flags wie `-k` oder `-j`) schreibt weder die Dateiliste noch `gates-passed.diffsha`; ein
    älterer Nachweis beider Dateien entfällt, auch bei grünem Lauf. Ein Aufruf mit `-n`, `-q`
    oder `-t` schreibt keinen Nachweis und entfernt keinen; das gilt auch zusammen mit `-i`
    (Vorrang wie in Punkt 6). Zugesagt ist das Ergebnis, nicht der Mechanismus: `-n`, `-q` und
    `-t` hält `make` selbst, weil das Rezept des Nachweises weder `$(MAKE)` noch `+` trägt und
    deshalb nicht läuft; `-i` prüft das Skript, das den Nachweis schreibt. Ein roter Lauf ohne
    `-i` lässt einen älteren Nachweis stehen, weil das Rezept des Nachweises nicht läuft
    (akzeptiertes Negativ; Grund: der ältere Nachweis beschreibt den Baum seines grünen Laufs;
    der Hook gibt auf ihn nur frei, wenn die Liste identisch ist oder jeder Unterschied auf der
    Positivliste liegt und der Stand der schnellen Kette grün ist, Punkt 9).
22. **Aufrufarten des Hooks.** Der Hook entscheidet über den Baum des Repos, in dem sein
    Arbeitsverzeichnis liegt (`git rev-parse --show-toplevel`), nicht über das Repo seiner Datei,
    und gibt in einer leeren Umgebung dieselbe Antwort. Zugesagt mit Fall: (a) ein Arbeitsverzeichnis in einem
    Unterverzeichnis des Repos gibt dieselbe Antwort wie die Wurzel; (b) eine leere Umgebung
    (`env -i`; `PATH` nicht gesetzt, `bash` mit absolutem Pfad gestartet) gibt bei identischer
    Liste `approve`, bei einem Unterschied, der nach Punkt 9 freigibt (jede geänderte,
    hinzugekommene oder gelöschte Datei auf der Positivliste und grüner Stand der schnellen
    Kette), ebenfalls `approve`, bei jedem anderen Unterschied `block`; die Freigabe nach
    Punkt 9 gilt damit in beiden Aufrufarten, (a) und (b). Der Hook ruft dafür nur `bash`, `git`, die bisher genutzten Programme und
    `tools/harness/working-tree-hash.sh`. Nicht zugesagt ist die Wirkung einzelner
    gesetzter Variablen (`PATH` für die Programmsuche, `TMPDIR` für die temporäre Datei, `GIT_*`
    über `git`); akzeptiertes Negativ: der Hook liest keine davon selbst, die Aufrufart des
    Treibers setzt sie nicht, und ein `TMPDIR` im Arbeitsbaum sperrt (`block`), gibt nie frei.
    Nicht Gegenstand: `bash -e` (der Hook setzt
    `set -euo pipefail` selbst, die Option ändert nichts) und eine fehlende Variable
    `CLAUDE_PROJECT_DIR` (der Aufruf des Treibers startet den Hook dann nicht; kein Code des
    Hooks läuft, er kann nicht antworten). Akzeptiertes Negativ: Ein Arbeitsverzeichnis in einem
    anderen Repo ist keine Randform dieses Hooks, sondern des Treibers (Aufruf-Voraussetzung,
    Punkt 8): der Hook urteilt über das Repo, in dem er läuft, nach dessen eigenem Stand; ein
    Repo ohne `tools/harness/dateiliste.sh`, ohne Stand oder ohne Git-Repo endet in `block`
    (Fälle von Punkt 19 und 20).
