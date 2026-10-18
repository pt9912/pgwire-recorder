# SPEC-061 — Hook für die Einlieferung in den Gegenproben-Container (`gegenproben-nach-einlieferung`, `gegenproben-nach-einlieferung-gegenprobe`)

Der Gegenproben-Container (Name vorläufig `pgr-gp`) fährt nur Gegenproben; in ihm wird kein Slice abgearbeitet. Sein
Bare-Repo trägt einen `post-receive`-Hook, dessen Inhalt das Skript `tools/harness/gegenproben-nach-einlieferung.sh`
ist: Jede Einlieferung eines Zweigs startet `make gegenproben` (die Liste für den Dev-Container, [`SPEC-060`](SPEC-060-gegenproben.md))
im Hintergrund auf genau dem eingelieferten Stand und legt das Ergebnis des beendeten Laufs ab. Das Skript ist kein Gate: Es
steht im Abschnitt *Werkzeuge* von `harness/README.md`, nicht in `GATE_CHECKS`, und kein Hook und kein Stand am Host glaubt
sein Ergebnis (Entscheidung des Nutzers vom 2026-10-18, [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 14). Der Grund für den
Container und seine Pflicht steht in [ADR-0048](../../docs/plan/adr/0048-gegenproben-verlassen-die-gate-kette.md); der Hook ist
die kleinste Form, die Einlieferung und Fahrt verbindet. Eine eigene ADR entsteht nicht (Entscheidung des Architect): Die
Entscheidung „der automatische Lauf liegt nur im eigenen Container, auf dem eingelieferten Stand“ trägt [ADR-0048](../../docs/plan/adr/0048-gegenproben-verlassen-die-gate-kette.md); der Hook ist
ihre Umsetzung, wirkt nur im Container und ändert weder Gate noch Schwelle (`AGENTS.md` §3.6).

Diese Kennung legt fest, was das Skript zusagt, woran es nichts startet und was die Gegenprobe belegt. Sie sagt nichts über
das, was `make gegenproben` fährt ([`SPEC-060`](SPEC-060-gegenproben.md)), nichts über den Aufbau des Containers
([`SPEC-058`](SPEC-058-container-start.md)) und nichts über den Wortlaut der Anweisungen der Rollen (Punkt 15 nennt nur die Regel an den Leser).
Der Name `gegenproben-nach-einlieferung-gegenprobe` folgt der Namensregel ([`SPEC-060`](SPEC-060-gegenproben.md) Punkt 2): Die Gegenprobe
steht in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS` und nicht in der Ausnahmeliste der schnellen Kette; ihr Fragment ist
`harness/mk/gegenproben-nach-einlieferung-gegenprobe.mk`, ihr Skript `tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh`.
Sie arbeitet in Temp-Verzeichnissen mit einem Bare-Repo, einer Kopie des Skripts als Hook und einer Attrappe als `make` (im `PATH` vorangestellt),
ruft Docker nicht und braucht `git`, `flock`, `setsid`, `pgrep` und `timeout`. **Voraussetzung der Gegenprobe: ein Benutzer ohne root** (Entscheidung des Architect; Grund: der Fall `p13-ablage-weg` macht die Ablage mit `chmod` unschreibbar, und root schließt `chmod` nicht aus). Als root endet sie mit Exit 2 und einer Meldung, die den Fall nennt, ohne Überspringen (wie [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 3); der Gegenproben-Container läuft als `vscode`. Dass die Prüfung selbst als root greift, belegt kein Fall (akzeptiertes Negativ; Grund: ein Fall bräuchte root).

Der Name jedes Falls trägt die Nummer des Punktes, den er prüft. Je Zusage und je Element einer Aufzählung gibt es einen Fall
und eine Mutation, die nur ihn (oder die genannten Fälle) rot färbt (`AGENTS.md` §3.10); die Gegenprobe verlangt genau diese Menge roter Fälle.

**Plattform.** Zugesagt ist ein Linux-Container oder -Host mit bash ≥ 4, `git` (gemessen im Image des Dev-Containers: 2.55.0), `flock` und
`setsid` (util-linux), `pgrep` (procps), GNU `make` und GNU-Coreutils. macOS- und Windows-Hosts sind nicht zugesagt (akzeptiertes Negativ; Grund: der Hook
läuft im Gegenproben-Container, und kein anderer Host steht zur Messung bereit); das Skript lehnt sie nicht ab, es meldet ein fehlendes
Werkzeug (Punkt 13).

**Randform → Punkt.**

| Randform des Plans | Punkt |
|---|---|
| H1 der Push wartet nicht | 5 |
| H2 welche Refs, mehrere Zweige | 3 |
| H3 zwei Einlieferungen kurz hintereinander | 6 |
| H4 gelöschter Zweig | 4 |
| H5 derselbe Hash erneut | 7 |
| H6 Name des Zweigs im Dateinamen | 8 |
| H7 Ergebnis nur bei beendetem Lauf | 9 |
| H8 nur die Container-Liste | 10 |
| H9 Marke des Containers | 2 |
| H10 wohin der Hook-Code gehört | 1 |
| H11 Berechtigung, Zustand, Reihenfolge der Gründe | 13 |
| H12 welcher Hash der Closure zählt | 15 |
| Aufrufarten des Treibers (`MAKEFLAGS`) | 11 |
| Aufruf-Voraussetzung | 14 |
| Plattform | Kopf dieser Kennung, 16 |
| neu vom Architect: Umgebung des Hooks (`GIT_*`) | 11 |
| neu vom Architect: Arbeitsbaum je Lauf | 12 |
| neu vom Architect: Zeilenformat und Form des Hashes | 4 |
| neu vom Architect: Länge des Namens | 8 |
| neu vom Architect: Ort der Ablage | 1, 2 |
| neu vom Architect: hängender Lauf, Zeitlimit | 16 |

**Begriffe.** *Ablage* ist das Verzeichnis, das die Marke (Punkt 2) nennt. *Name* ist der Zweigname ohne `refs/heads/`, kodiert nach Punkt 8.
*Hash* ist der neue Wert der Zeile. Je eingeliefertem Zweig entstehen in der Ablage vier Pfade, `<N>` für `gegenproben-<Name>-<Hash>`:
das Ergebnis `<N>`, die Marke `<N>.laeuft`, das Log `<N>.log` und, während des Laufs, ein Arbeitsbaum `arbeitsbaum.<Zufall>`; dazu die Sperre `gegenproben.lock`.

1. **Gegenstand, Ort, Kopie, Ablage.** Das Skript ist eine einzelne Datei im Repo und braucht zur Laufzeit keine andere Datei des Repos: Es
   wird als **Kopie** mit `install -m 755` nach `<Bare-Repo>/hooks/post-receive` gelegt, kein Symlink. Grund: Der Lauf checkt den eingelieferten Hash
   in einen eigenen Arbeitsbaum aus, ein Symlink in einen Klon zeigte auf eine Datei, die ein Wechsel des Klons ersetzt, und `bash` liest ein
   laufendes Skript fortlaufend. Das Skript liest nichts als die Zeilen von `git` auf stdin, die Marke (Punkt 2) und das Bare-Repo; es
   braucht keinen Klon. Die **Ablage** liegt deshalb nicht im Klon, sondern dort, wo die Marke sie nennt (Entscheidung des Architect; Abweichung vom Vorschlag des
   Planners „Pfad im Klon“; Grund: ein Klon ist keine Voraussetzung der Fahrt, und ein fehlender Klon wäre eine Randform mehr). Dass die Kopie
   dem Skript im Repo entspricht, prüft nichts (akzeptiertes Negativ; Grund: die Kopie liegt außerhalb des Repos, und das Aufsetzen wiederholt man nach
   jeder Änderung des Skripts; `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt` führt es). Fälle: `p1-kopie` (das Skript läuft aus einer Kopie in einem
   Verzeichnis ohne Repo und ohne den Pfad des Repos in der Umgebung und lädt nichts nach), `p1-echter-push` (ein echter `git push` in das Bare-Repo mit
   dem Skript als `hooks/post-receive` löst einen Lauf aus, der Push endet mit Exit 0 vor dem Ende des Laufs). Mutationen: das Skript lädt eine Datei relativ zu `$0`
   nach (färbt `p1-kopie`); das Skript liest stdin nicht (färbt `p1-echter-push`).
2. **Marke des Containers (H9).** Das Skript tut nur etwas, wenn in `<Bare-Repo>/gegenproben-ablage` eine einzige Zeile steht: der absolute Pfad der Ablage, eines
   vorhandenen Verzeichnisses (ob es schreibbar ist, entscheidet Punkt 13). Die Zeile zum Aufsetzen legt die Marke und die Kopie des Hooks zusammen an; das Skript lehnt damit einen Hook
   ab, der allein in das Bare-Repo eines Arbeits-Containers kopiert wurde. Wer die ganze Zeile in einem Arbeits-Container ausführt, hat den Hook dort aufgesetzt
   (akzeptiertes Negativ; Grund: ein Name des Containers als Marke wäre eine Liste, die dieselbe Hand pflegt wie das Aufsetzen; das Dokument nennt den Container in der Zeile).
   Fehlt die Marke, ist sie leer, nicht absolut, hat mehr als eine Zeile oder nennt sie kein vorhandenes Verzeichnis: Meldung nach Punkt 13, Exit 0, **nichts** wird angelegt, kein Prozess
   gestartet, das Bare-Repo nicht verändert. Fälle: `p2-fehlt`, `p2-leer`, `p2-relativ`, `p2-zwei-zeilen`, `p2-kein-verzeichnis` (je einer; bei jedem: Exit 0, Meldung mit
   Schlüsselwort, die Ablage und das Bare-Repo unverändert, die Attrappe von `make` nicht gerufen). Mutationen: je Form die Prüfung gestrichen.
3. **Welche Refs (H2).** Ein Lauf entsteht nur für eine Zeile, deren Ref mit `refs/heads/` beginnt. Tags (`refs/tags/…`), Notizen (`refs/notes/…`), `refs/headsx/…` und
   jeder andere Ref lösen keinen Lauf aus, und zwar ohne Meldung. Liefert ein Push mehrere Zweige ein, entsteht je Zweig ein Lauf, **nacheinander in der Reihenfolge der Zeilen** (ein
   Hintergrundprozess je Push, der die Zeilen abarbeitet). Fälle: `p3-tag`, `p3-notes`, `p3-headsx` (kein Lauf, keine Ausgabe), `p3-mehrere`
   (drei Zweige, die Attrappe protokolliert Beginn und Ende; die Folge ist Zeile 1, 2, 3 ohne Überlappung, je Zweig ein Ergebnis). Mutationen: das Muster `refs/*` statt `refs/heads/*`
   (färbt `p3-tag`, `p3-notes`); das Muster `refs/heads` ohne den Schrägstrich (färbt `p3-headsx`); die Zweige laufen gleichzeitig; die Zeilen laufen rückwärts; nur die erste Zeile läuft (die drei färben `p3-mehrere`).
4. **Zeilen: gelöscht, Form (H4).** Zeilen haben die Form `<alt> <neu> <ref>`. Ist `<neu>` nur Nullen (40 oder 64 Stellen), wurde der Zweig gelöscht: kein Lauf, keine Meldung, ein
   vorhandenes Ergebnis, Log und Marke bleiben. Ist `<neu>` weder 40 noch 64 Stellen aus `0-9a-f`, oder fehlt ein Feld der Zeile, endet nur diese Zeile mit einer Meldung (Schlüsselwort `zeile`),
   die übrigen laufen. Leere Eingabe ist kein Fehler: kein Lauf, keine Meldung; ebenso eine **leere oder nur aus Leerraum bestehende Zeile** mitten in der Eingabe: schweigend übersprungen (Grund: `git` sendet sie nie, und eine Meldung nützte dem Pusher nicht). Eine Zeile mit **mehr als drei Feldern** ist ungültig (ein Ref enthält nach `git check-ref-format` kein Leerzeichen): Meldung `zeile`, die Zeile läuft nicht. Die **letzte Zeile ohne Zeilenumbruch** wird verarbeitet wie jede andere. Fälle: `p4-null40`, `p4-null64`, `p4-ergebnis-bleibt` (ein Ergebnis eines früheren Laufs desselben Zweigs steht nach der Löschung unverändert da),
   `p4-hash-kurz` (39 Stellen), `p4-hash-gross` (Großbuchstaben), `p4-hash-41` (41 Stellen), `p4-hash-64` (64 Stellen aus `0-9a-f`, angenommen), `p4-feld-fehlt`, `p4-stdin-leer`, `p4-leerzeile` (leere Zeile zwischen zwei guten: keine Meldung, beide laufen), `p4-leerzeichen-zeile`, `p4-vier-felder` (Meldung, kein Lauf), `p4-ohne-umbruch` (die letzte Zeile ohne Umbruch läuft), `p4-zeile-kaputt-dann-gut` (eine
   kaputte und eine gute Zeile: die gute läuft). Mutationen: Null-Hashes laufen; Nullen nur mit 40 Stellen erkannt (färbt `p4-null64`); keine Prüfung der Form; nur 40 Stellen angenommen (färbt `p4-hash-64`);
   Großbuchstaben angenommen; 41 Stellen angenommen; eine kaputte Zeile beendet das Skript (färbt `p4-zeile-kaputt-dann-gut`); eine leere Zeile meldet (färbt `p4-leerzeile`, `p4-leerzeichen-zeile`); mehr als drei Felder gelten als Ref mit Leerzeichen (färbt `p4-vier-felder`); die letzte Zeile ohne Umbruch entfällt (färbt `p4-ohne-umbruch`); das Ergebnis wird bei Löschung entfernt.
5. **Der Push wartet nicht (H1).** Der Hook kehrt zurück, bevor `make` endet, mit Exit 0 in jedem Fall dieser Kennung (auch bei den Meldungen aus Punkt 13); er schreibt nichts auf stdout; auf stderr stehen nur
   die Meldungen nach Punkt 13 und, je gestartetem Lauf, eine Zeile `gegenproben-nach-einlieferung: gestartet: <Zweig> <Hash>` (Anfang zugesagt, der Rest ist Text). Der Lauf hat weder
   stdout noch stderr noch stdin des Hooks: Er schreibt in das Log und liest `/dev/null`, sonst hielte er die Leitung des Pushs offen (`git` wartet auf das Ende der Leitung). Er überlebt das Ende der Sitzung und ein
   Signal an die Prozessgruppe des Hooks (eigene Sitzung). Fälle: `p5-vor-ende` (die Attrappe wartet; der Hook ist zurück, bevor sie endet; das Ergebnis fehlt zu diesem Zeitpunkt), `p5-leitung` (Aufruf mit stdout und stderr an eine Leitung, deren Leser bis zum Ende der
   Leitung wartet: der Leser endet vor dem Ende des Laufs), `p5-gruppe` (ein `SIGTERM` an die Prozessgruppe des Hook-Aufrufs beendet den Lauf nicht), `p5-stdout-leer`, `p5-gestartet-zeile`.
   Mutationen: der Hook wartet auf den Lauf (färbt `p5-vor-ende`); der Lauf erbt stdout und stderr (färbt `p5-leitung`); der Lauf bleibt in der Sitzung des Hooks (färbt `p5-gruppe`); der Hook gibt Exit 1 zurück, wenn der Lauf rot endet;
   eine Zeile auf stdout; die Zeile `gestartet` fehlt.
6. **Zwei Einlieferungen kurz hintereinander (H3).** Ein Lauf zur Zeit: Die Läufe zweier Einlieferungen **überlappen nicht**, der laufende Lauf wird nicht abgebrochen, der zweite wartet auf die Sperre
   (`gegenproben.lock` in der Ablage) und läuft danach; alle Läufe laufen nacheinander und jeder legt sein Ergebnis ab. Die Reihenfolge der Läufe **verschiedener** Einlieferungen ist nicht zugesagt (akzeptiertes Negativ;
   Grund: jedes Ergebnis gilt je Hash, die Reihenfolge ändert es nicht, und eine Warteschlange mit Ordnung wäre ein Dienst). Fälle: `p6-nacheinander` (zwei Hook-Aufrufe im Abstand von
   einer Sekunde, die Attrappe wartet; Beginn und Ende der Protokolle überlappen nicht), `p6-kein-abbruch` (der erste Lauf endet mit seinem eigenen Ergebnis `exit=0`, der zweite hat auch eines). Mutationen: ohne Sperre (färbt `p6-nacheinander`); die Sperre wird nicht
   erwartet, sondern der zweite Lauf entfällt, wenn sie belegt ist (färbt `p6-kein-abbruch`); der zweite Aufruf beendet den ersten (färbt `p6-kein-abbruch`).
7. **Derselbe Hash erneut (H5).** Ein erneuter Push desselben Zweigs mit demselben Hash, ein Force-Push auf einen früheren Stand und ein anderer Zweig mit demselben Hash lösen einen neuen Lauf aus. Ein vorhandenes Ergebnis desselben Namens wird
   **beim Aufruf des Hooks entfernt** (vor der Rückkehr) und nach dem Lauf neu geschrieben; es wird nie als Stand wiederverwendet. Zwei Läufe desselben Namens zugleich in der Warteschlange: die Marke endet mit dem ersten, das Ergebnis des letzten gewinnt (akzeptiertes Negativ;
   Grund: die Closure liest erst nach dem Ende aller Läufe, und ein Zähler je Marke wäre ein Dienst). Fälle: `p7-neu` (erster Lauf Exit 1, zweiter Exit 0: das Ergebnis ist `exit=0`), `p7-altes-weg` (zwischen Rückkehr des Hooks und Ende des Laufs existiert kein Ergebnis dieses Namens), `p7-anderer-zweig` (zwei Zweige,
   ein Hash: zwei Ergebnisse). Mutationen: ein vorhandenes Ergebnis lässt den Lauf entfallen (färbt `p7-neu`); das alte Ergebnis wird nicht entfernt (färbt `p7-altes-weg`); der Name enthält den Zweig nicht (färbt `p7-anderer-zweig`).
8. **Name im Dateinamen (H6).** Der Name entsteht aus dem Zweignamen, indem jedes Byte außer `A-Za-z0-9._-` als `%` und zwei Großbuchstaben-Hexziffern geschrieben wird (`thema/teil` wird `thema%2Fteil`, `ü` wird `%C3%BC`, `%` wird `%25`); verschiedene
   Zweige ergeben verschiedene Namen, und jede Form, die `git check-ref-format` zulässt, landet in genau einer Datei der Ablage. Ein Punkt am Anfang bleibt, weil jeder Pfad mit `gegenproben-` beginnt. Ist der kodierte Name länger als 120 Zeichen,
   endet der Zweig mit einer Meldung (Schlüsselwort `name`) und ohne Lauf (Grenze der Dateinamen; Entscheidung des Architect). Fälle: `p8-schraeg`, `p8-prozent`, `p8-nichtascii`, `p8-sonderzeichen` (`+`), `p8-unveraendert` (`a.b_c-d` und `.x` bleiben), `p8-eindeutig` (`a/b` und `a%2Fb` ergeben zwei Ergebnisse), `p8-lang`
   (121 Zeichen: Meldung, kein Lauf), `p8-grenze` (120 Zeichen: Lauf). Mutationen: keine Kodierung (der Schrägstrich legt ein Unterverzeichnis an, färbt `p8-schraeg`); `%` bleibt unkodiert (färbt `p8-prozent`, `p8-eindeutig`); nur der Schrägstrich wird kodiert (färbt `p8-nichtascii`, `p8-sonderzeichen`);
   Kleinbuchstaben-Hex (färbt `p8-schraeg`); die Grenze bei 121 (färbt `p8-lang`) und bei 119 (färbt `p8-grenze`).
9. **Ergebnis, Marke, Log (H7).** Beim Aufruf des Hooks legt er je Lauf die Marke `<N>.laeuft` an. Der Lauf wartet auf die Sperre, legt den Arbeitsbaum an (Punkt 12), ruft `make` (Punkt 10) mit Ausgabe in `<N>.log` (das Log wird beim Lauf neu beschrieben, nicht fortgesetzt)
   und schreibt **nur nach dem Ende von `make`** das Ergebnis: zuerst in eine Datei der Ablage, die mit `.` beginnt, dann mit `mv` auf `<N>`; danach entfernt er die Marke und den Arbeitsbaum. Das Ergebnis hat genau diese Zeilen in dieser Reihenfolge: `zweig=<Zweig ohne refs/heads/>`, `hash=<Hash>`, `exit=<Exit von make>`,
   `ergebnis=gruen` (genau bei Exit 0) oder `ergebnis=rot`, `beginn=<UTC>` (nach der Sperre, unmittelbar vor `make`), `ende=<UTC>` (nach `make`), `rot=<Ziele>`; Zeit im Format `%Y-%m-%dT%H:%M:%SZ`, `rot` die Namen der Ziele aus den Zeilen
   `gegenproben: <ziel>: rot (Exit <n>)` des Logs, durch je ein Leerzeichen getrennt, in der Reihenfolge des Logs, leer wenn keine. Ein Ergebnis steht nur nach einem **beendeten** `make`: Ein Abbruch (`SIGTERM`, Punkt 16 für `SIGKILL`), ein
   Neustart des Containers, eine Sperre, ein Arbeitsbaum oder ein Lauf, der nicht zustande kam, hinterlassen **kein** Ergebnis und eine **verwaiste Marke**; die Marke ohne Ergebnis bedeutet „läuft oder abgebrochen“ und ist nie ein Ergebnis. `ergebnis=gruen` und eine Marke
   zugleich liest der Leser nicht als grün (Punkt 15). Ein `SIGTERM` an den Hintergrundprozess, während er wartet oder `make` läuft, beendet `make` samt Nachkommen, räumt den Arbeitsbaum, lässt die Marke stehen, legt kein Ergebnis an und endet mit Exit 143 (`SIGINT` und
   `SIGHUP` sind nicht zugesagt: ein Hintergrundprozess einer nicht interaktiven Shell ignoriert `SIGINT`; akzeptiertes Negativ, Grund: der Lauf gehört keiner Sitzung). Fälle: `p9-gruen` (Zeilen, Reihenfolge, `exit=0`, `ergebnis=gruen`, leeres `rot=`, Zeitformat, `beginn` ≤ `ende`), `p9-rot`
   (Exit 2: `ergebnis=rot`, `rot=` nennt zwei Ziele in der Reihenfolge des Logs), `p9-marke` (die Marke steht beim Aufruf des Hooks da und nach dem Ende des Laufs nicht), `p9-nicht-vorher` (während die Attrappe wartet, steht kein Pfad `<N>` in der Ablage; Marke, Log und eine Datei, deren Name mit `.` beginnt, sind dort erlaubt), `p9-sigterm-lauf`, `p9-sigterm-wartet` (kein Ergebnis, Marke bleibt, `make` und sein Hintergrundkind beendet, Exit 143), `p9-log` (die Ausgabe von stdout und stderr der Attrappe steht im Log; beim zweiten Lauf nicht doppelt), `p9-beginn-nach-sperre` (der zweite Lauf trägt `beginn` nach dem `ende` des ersten).
   Mutationen: `ergebnis=gruen` trotz Exit 2; `ergebnis` aus dem Exit des letzten Ziels; `rot=` nur das erste Ziel; Zeilen vertauscht; Zeitformat mit lokaler Zeit; Marke bleibt nach dem Ende; Marke fehlt; das Ergebnis wird mit `>` unter dem Endnamen angelegt (färbt `p9-nicht-vorher`); Ergebnis bei `SIGTERM`
   trotzdem (färbt beide `p9-sigterm-*`); die Nachkommen von `make` bleiben am Leben; `SIGTERM` in der Wartezeit ohne Behandlung (färbt `p9-sigterm-wartet`); Exit 0 statt 143; das Log wird angehängt statt neu beschrieben; `beginn` vor der Sperre (färbt `p9-beginn-nach-sperre`).
10. **Nur die Container-Liste (H8; Entscheidung des Nutzers vom 2026-10-18).** Der Lauf ruft genau `make gegenproben`: ein Argument, keine Option, kein Variablen-Argument. Er ruft nie `gegenproben-host`, kein Einzelziel und kein `-j`, `-k`, `-i`, `-n`. Die Attrappe protokolliert ihr `argv`. Fall `p10-argv` (genau ein
    Argument, `gegenproben`). Mutationen: `gegenproben-host`; zusätzlich `-k`; zusätzlich `-j4`; das Einzelziel statt des Sammelziels.
11. **Umgebung des Laufs (MAKEFLAGS, `GIT_*`).** Der Lauf ruft `make` ohne die Variablen `MAKEFLAGS`, `MFLAGS`, `MAKELEVEL`, `MAKEOVERRIDES` und ohne jede Variable, deren Name mit `GIT_` beginnt, auch wenn der Aufruf des Hooks sie trägt (`git` setzt `GIT_DIR` und, bei Quarantäne, `GIT_QUARANTINE_PATH`, und
    ein geerbtes `-i` ließe `make gegenproben` mit Exit 2 enden, [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 10 b). Fälle, je Variable einer: `p11-makeflags` (der Aufruf trägt `MAKEFLAGS=-i`), `p11-mflags`, `p11-makelevel`, `p11-makeoverrides`, `p11-git-dir`, `p11-git-quarantaene`; die Attrappe schreibt ihre Umgebung. Mutationen: je Variable wird
    sie nicht zurückgesetzt; nur `MAKEFLAGS` und `MFLAGS` (färbt die übrigen); nur `GIT_DIR` (färbt `p11-git-quarantaene`).
12. **Arbeitsbaum je Lauf.** Der Lauf checkt genau den Hash der Zeile in einen **eigenen** Arbeitsbaum in der Ablage aus (Name `arbeitsbaum.<Zufall>`, je Lauf neu), ohne Zweig (losgelöst), und ruft `make` dort. Er entfernt ihn nach dem Ende von `make`, grün wie rot, und bei `SIGTERM`; er verändert
    keinen Ref und keine Konfiguration des Bare-Repos. Ein Rest eines abgebrochenen Laufs (`SIGKILL`) wird nicht wiederverwendet und nicht gelöscht (Punkt 16). Fälle: `p12-stand` (die Attrappe schreibt das Arbeitsverzeichnis, `git rev-parse HEAD`, `git status --porcelain` und `git symbolic-ref -q HEAD`: das Verzeichnis liegt in der Ablage, `HEAD` ist der Hash der Zeile, der
    Baum sauber, kein Zweig), `p12-entfernt-gruen`, `p12-entfernt-rot`, `p12-bare-unveraendert` (Refs und `config` des Bare-Repos vor und nach dem Lauf gleich), `p12-rest` (ein früherer Rest liegt da, der Lauf verwendet ihn nicht und lässt ihn stehen). Mutationen: der Lauf fährt `make` im Verzeichnis des Bare-Repos
    oder in einem Klon; der Spitzenstand des Zweigs statt des Hashes (färbt `p12-stand` bei einem Hash, der nicht die Spitze ist); der Arbeitsbaum bleibt nach grün; nach rot; ein fester Name; der Lauf löscht Reste.
13. **Voraussetzungen und Reihenfolge der Gründe (H11).** Das Skript prüft in dieser Reihenfolge und endet beim ersten Treffer **für den ganzen Aufruf** (Schritte 1 bis 3) oder **für die Zeile** (Schritt 4): (1) die Marke (Punkt 2), Schlüsselwort `marke`; (2) die Werkzeuge `git`, `flock`, `setsid`, `make` und `pgrep` (es findet beim Abbruch die Nachkommen von `make`) im `PATH`, Schlüsselwort
    `werkzeug` und der Name des ersten fehlenden; (3) die Ablage, Schlüsselwort `ablage`: schreibbar zum Zeitpunkt des Aufrufs (das Anlegen von `gegenproben.lock` oder einer Datei in ihr gelingt; scheitert es, endet der Aufruf mit dieser Meldung); (4) je Zeile: Ref (Punkt 3, schweigt), Löschung (schweigt), Form des Hashes (Punkt 4), Länge des Namens (Punkt 8).
    Danach, im Hintergrund und nur im Log mit Zeile `gegenproben-nach-einlieferung: <sperre|arbeitsbaum>: …`: (5) die Sperre, (6) der Arbeitsbaum (der Hash fehlt im Bare-Repo, `git worktree add` scheitert), (7) der Lauf. Jede Meldung beginnt mit `gegenproben-nach-einlieferung: ` und dem Schlüsselwort und einem Doppelpunkt; zugesagt und geprüft sind nur
    Präfix und Schlüsselwort. In jedem Fall endet der Hook mit Exit 0, schreibt **kein** Ergebnis (nie `gruen`) und startet für den betroffenen Teil nichts. Ein Lauf, in dem `make` selbst endet (auch mit 127), ist ein beendeter Lauf mit Ergebnis `rot` (Punkt 9); ebenso `make`, das ein Signal beendet (Exit über 128, etwa 137 nach `SIGKILL` von außen): Ergebnis `rot` mit diesem `exit=` (Fall `p13-make-signal`, die Attrappe beendet sich mit `SIGKILL`; Mutation: bei Exit über 128 entfällt das Ergebnis). Das ist nicht der Abbruch des Laufs durch `SIGTERM` (Punkt 9). Ist die Datei nicht ausführbar, ignoriert `git` den Hook
    mit einer Warnung, der Push gelingt und es läuft nichts; das Skript sieht es nicht (Grenze; das Dokument nennt die Probe: Testeinlieferung und Ergebnis). Fälle: `p13-werkzeug-git`, `p13-werkzeug-flock`, `p13-werkzeug-setsid`, `p13-werkzeug-make`, `p13-werkzeug-pgrep` (der `PATH` enthält die übrigen Werkzeuge, das eine fehlt), `p13-ablage-weg` (Marke nennt ein Verzeichnis, das beim Aufruf nicht mehr
    schreibbar ist), `p13-sperre` (die Sperre ist ein Verzeichnis: Log-Zeile `sperre`, kein Ergebnis, Marke bleibt), `p13-hash-fehlt` (der Hash steht nicht im Bare-Repo: Log-Zeile `arbeitsbaum`, kein Ergebnis, Marke bleibt, Exit 0), `p13-make-127` (die Attrappe endet mit 127: Ergebnis `exit=127`, `ergebnis=rot`), Paare der Reihenfolge `p13-marke-vor-werkzeug`
    (keine Marke und kein `git`: Schlüsselwort `marke`), `p13-werkzeug-vor-ablage`, `p13-ablage-vor-zeile` (Ablage nicht schreibbar und kaputte Zeile: Schlüsselwort `ablage`). Mutationen: je Werkzeug (auch `pgrep`) die Prüfung gestrichen; das Paar umgekehrt (je eine); Exit 1 statt 0 bei jedem Grund (färbt jeden Fall dieses Punktes); ein Ergebnis `gruen` ohne beendetes `make` (färbt `p13-sperre`,
    `p13-hash-fehlt`); die Marke entfällt bei Fehler in Sperre oder Arbeitsbaum.
14. **Aufruf-Voraussetzung (§3.12).** Zugesagt ist die Folge der Läufe über die Sperre (Punkt 6). Zwei Aufrufe des Hooks **gleichzeitig** mit gleichem Namen (Zweig und Hash) haben für die Schritte vor der Sperre (Marke anlegen, altes Ergebnis entfernen) keine Zusage (akzeptiertes Negativ; Grund: ein
    Push desselben Zweigs mit demselben Hash im selben Augenblick gibt es nicht ohne Absicht, und das Ergebnis des letzten Laufs gilt, Punkt 7); die Zusagen von `git` über gleichzeitige Pushes in dasselbe Bare-Repo gelten nicht für die Ablage. Das Skript ruft `make` nicht aus einem Rezept, deshalb gibt es kein `-n`, `-q`, `-t` für den Hook; `MAKEFLAGS` aus der Umgebung ist Punkt 11.
15. **Lesen des Ergebnisses; welcher Hash zählt (H12).** Regel an den Leser, kein Fall (Grund: der Leser ist die Closure von Welle und Release und das Review, ein Mensch oder ein Agent; ein Sensor dafür wäre ein Stand der Gegenproben, den [`SPEC-060`](SPEC-060-gegenproben.md) Punkt 14 ausschließt). Die Closure liest die Ablage des Gegenproben-Containers
    über `ssh` mit dem Hash als Muster `gegenproben-*-<Hash>` (voller Hash, `git rev-parse`); ein Hash gilt als grün, wenn mindestens eine Datei vorliegt, alle Dateien `ergebnis=gruen` tragen und keine Marke `gegenproben-*-<Hash>.laeuft` steht. Eine Datei `rot`, eine Marke desselben Hashes oder kein Ergebnis heißt:
    nicht grün, die Closure fährt `make gegenproben` selbst. Der Hash ist der des Commits, der den **letzten nicht inerten Stand** trägt; nicht inert heißt ausgenommen die Positivliste von [`SPEC-052`](SPEC-052-gate-kette.md) Punkt 12 (Planung, Reviews, Dokumentation). Ein Commit, der nur inerte Pfade ändert, braucht keinen neuen Lauf; die Closure-Notiz nennt beide Hashes.
    Das Ergebnis ist ein Nachweis zum Nachlesen, kein Stand: Es liegt im Container, nicht im Arbeitsbaum des Hosts, und der Stop-Hook liest es nicht.
16. **Grenzen.** (1) Ein Zeitlimit des Laufs gibt es nicht; ein hängender Lauf hält die Sperre, bis er endet oder `SIGTERM` bekommt, und die folgenden Läufe warten (akzeptiertes Negativ, Entscheidung des Architect; Grund: die Dauer einer Fahrt ist nicht gemessen, und ein Limit wäre eine ungemessene Zahl; der Implementer
    misst die Dauer in der ersten Fahrt, `SPEC-060` Punkt 15 (7)). (2) `SIGKILL`, ein Neustart des Containers und ein Ausfall des Docker-Dienstes im Container hinterlassen eine verwaiste Marke, ein Log ohne Ergebnis, einen Rest des Arbeitsbaums und möglicherweise Container und Netze der Einzelziele; aufgeräumt wird von Hand (Dokument). (3) Ergebnisse, Logs und Reste wachsen ohne Grenze; kein
    Aufräumen durch das Skript. (4) Mehrere Zweige in einem Push fahren nacheinander, die Dauer ist die Summe. (5) Die Uhr ist die des Containers, `UTC`. (6) Ob `make gegenproben` im frischen Container von selbst grün ist, belegt die erste Fahrt der Lieferung, kein Fall ([`SPEC-060`](SPEC-060-gegenproben.md) Punkt 3). (7) Dass die Kopie im Bare-Repo dem Skript im Repo entspricht, prüft nichts (Punkt 1).
    (8) Die Locale des Containers bestimmt die Sprache von `make`; das Skript liest die Zeilen der Zusammenfassung von `gegenproben.sh`, die in keiner Locale übersetzt sind.
