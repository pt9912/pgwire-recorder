# SPEC-054 — Abgleich von Code-Tabelle und Meldungskatalog (`meldungskatalog-check`)

`make meldungskatalog-check` gleicht die Meldungscodes `PGR-…` dreier Quellen ab, ohne über
Inhalt zu urteilen: die Code-Tabelle im Quelltext (C), den Katalog der Spezifikation (K, Tabelle
in `SPEC-034`) und die Fehlercode- und Warnungstabellen des Benutzerhandbuchs (H). Die Punkte 1
bis 10 prüft `make meldungskatalog-check-gegenprobe`; Punkt 11 hat dort den Fall mit Makefile.

1. **Quellen.** Gelesen werden drei Dateien und ein Baum, alle relativ zur Wurzel des Repos.
   - C: jede Zeile von `internal/hexagon/model/fehler.go`, die nach beliebigem Einzug mit
     `Code` und Wortzeichen beginnt, dann Leerraum, `=`, Leerraum und eine Zeichenkette in
     doppelten Anführungszeichen trägt. Der Code ist der Text zwischen den Anführungszeichen.
     Eine Zeile mit anderem Anfang (Kommentar, Klassen-Tabelle, Funktion) ist keine Konstante.
   - K: die Zeilen der Tabelle in `spec/spezifikation.md` von der Zeile `### SPEC-034 ` bis zur
     nächsten Zeile, die mit `#` beginnt.
   - H: dieselbe Zeilenform in `docs/user/benutzerhandbuch.md` in den Abschnitten, die mit den
     Zeilen `### Fehlercodes` und `### Warnungen` beginnen und je bis zur nächsten Zeile enden,
     die mit `#` beginnt; beide Abschnitte zusammen sind H.
   - Eine Tabellenzeile ist eine Zeile, deren erste Zelle ein Text in Backticks ist
     (`` | `PGR-E2001` | … ``); der Code ist dieser Text. Kopf- und Trennzeile sind keine
     Zeilen. Eine Zelle gilt als ein Code, auch wenn sie mehrere nennt (dann Punkt 2).
   - Der Baum für Punkt 6: die Go-Dateien unter `cmd/` und `internal/`, deren Name nicht auf
     `_test.go` endet.
2. **Form.** Ein Code in C, K oder H passt auf `^PGR-[EWI][0-9]{4}$`, sonst ist das ein Befund
   `Form` mit dem Text als Code. Er nimmt an den Regeln der Punkte 3 bis 5 nicht teil. `PGR-E123`,
   `PGR-X0001`, ein Platzhalter und eine Zelle mit zwei Codes sind Form-Befunde.
3. **Doppelt.** Steht ein Code in einer Quelle mehrfach, ist das ein Befund `doppelt`, je Code und
   Quelle einmal, an der zweiten Fundstelle; der Code zählt einmal.
4. **Mengen.** Je Regel ein Befund je Code:
   - a. Ein Code in C, der nicht in K steht: `nicht im Katalog`, an der Zeile der Konstante.
   - b. Ein Code in K, der nicht in C steht und nicht auf `000` endet: `nicht in der
     Code-Tabelle`, an der Zeile des Katalogs. Ein Code in K, der auf `000` endet und nicht in
     C steht, ist zugelassen (ein Rückfall, den der Code nie erzeugt).
   - c. Ein Code in C, der nicht in H steht: `nicht im Handbuch`, an der Zeile der Konstante.
   - d. Ein Code in H, der nicht in C steht: `nicht in der Code-Tabelle`, an der Zeile des
     Handbuchs, auch wenn K ihn führt. Das Handbuch führt damit genau die Codes, die der Code
     erzeugt; einen nicht erzeugten Rückfall führt es nicht.
5. **Rückfall je Klasse.** Für jede Klasse 1 bis 6 steht `PGR-E<k>000` in K, sonst ein Befund
   `Rückfall fehlt` mit diesem Code, am Pfad von K ohne Zeile.
6. **Code außerhalb der Code-Tabelle.** Eine Zeile einer Datei des Baums aus Punkt 1, nicht
   `fehler.go`, mit `"PGR-` (Anführungszeichen unmittelbar davor) ist ein Befund `Code
   außerhalb der Code-Tabelle`; der Code ist der Text ab `PGR-` aus Wortzeichen. Eine Zeichenkette
   in Backticks, eine Konstante in `fehler.go` unter anderem Namen als `Code…` und ein
   Kommentar ohne Anführungszeichen zählen nicht (Grenze).
7. **Zurückgezogen, reserviert.** Der Katalog hat keine Form für einen zurückgezogenen Code, und
   das Werkzeug kennt keine; die Schwere `I` und der Warnbereich 9 sind Codes wie die übrigen.
   Ein zurückgezogener Code ohne Konstante und ohne Endung `000` ist ein Befund nach 4b. Der
   Slice, der einen Code zurückzieht, ändert diese Kennung in derselben Änderung (Stand
   2026-10-17: kein Code ist zurückgezogen).
8. **Sperren.** Nichts wird abgeglichen und es gibt keine Befund-Zeile, wenn eine dieser
   Bedingungen gilt; geprüft wird in dieser Reihenfolge, die erste gilt, und die Zeile auf
   stderr lautet `meldungskatalog-check: <grund>`: eine der drei Dateien fehlt oder ist
   nicht lesbar (`<pfad> fehlt`, `<pfad> nicht lesbar`), `cmd/` oder `internal/` fehlt
   (`<verzeichnis> fehlt`), ein Abschnitt fehlt (`Abschnitt <überschrift> fehlt`), eine Quelle
   hat keine Zeile (C: `<pfad> ohne Code`; K und H: `<pfad> ohne Codes`; H zählt beide
   Abschnitte zusammen). Eine Sperre geht jedem Befund vor.
9. **Ausgabe und Ausgang.** Je Befund eine Zeile auf stderr,
   `meldungskatalog-check: <pfad>[:<zeile>]: <code>: <befund>`, der Pfad relativ zur Wurzel.
   Sortiert wird nach Pfad, Code und Befund in Bytefolge, danach nach der Zeile als Zahl;
   stdout bleibt leer. Ausgang 0 ohne Befund und ohne Ausgabe, 1 bei mindestens einem Befund
   (alle sind gelesen und ausgegeben), 2 bei einer Sperre.
10. **Aufrufarten.** Das Werkzeug nimmt höchstens die Wurzel als Argument (ohne: das aktuelle
    Verzeichnis), schreibt nichts und hinterlässt keinen Stand; was `make` mit `-n`, `-q`, `-t`,
    `-i` oder `MAKEFLAGS` aus der Umgebung tut, ändert seine Aussage nicht, weil kein Hook und
    kein Gate mehr als den Ausgang liest. Dass `-i` den roten Ausgang überliest, fängt der
    Nachweis des vollen Laufs (`SPEC-052`, Punkt 21), nicht dieses Werkzeug.
11. **Anschluss ohne Stufung.** `meldungskatalog-check` und `meldungskatalog-check-gegenprobe`
    hängen über das Fragment `harness/mk/meldungskatalog-check.mk` an `GATE_CHECKS` und enden bei
    rotem Lauf mit Ausgang ungleich 0. Die Prüfung ist vom ersten Lauf an scharf und am Stand
    2026-10-17 grün (C 30 Codes, K 34, H 30; die vier Codes von K ohne Konstante sind die
    Rückfälle `PGR-E2000`, `PGR-E3000`, `PGR-E5000`, `PGR-E6000`). Die Prüfung fährt in der
    schnellen Kette und muss dort bleiben, weil die Positivliste von `SPEC-052` Änderungen am
    Handbuch ohne vollen Lauf zulässt; die Gegenprobe steht in der Ausnahmeliste. Die Gegenprobe
    baut ihre Bäume selbst und liest weder Handbuch noch Spezifikation des Arbeitsbaums; für
    diesen Punkt kopiert sie nur Makefile und Fragment.

**Grenze.** Nicht geprüft sind die Spalten *Klasse / Bereich*, *Bedeutung* und *Maßnahme*
(passt die erste Ziffer zur Klasse, stimmt der Text), Codes in Kommentaren, Prosa, Tests,
`test/`, Plänen und Lastenheft, ein Code in Backticks im Quelltext, der Exit Code und der Text
eines Codes im Verhalten, und ob das Handbuch eine Zeile sinnvoll beschreibt. Ob K eine
vergebene Zeile ohne Konstante nur als Rückfall führt, regelt Punkt 4b; ein Code, der vergeben
und noch nicht erzeugt ist, hat in diesem Stand keine Form. Die Abschnitts-Erkennung folgt den
Zeilen `#`: Eine Zeile `#` in einem Codeblock beendet einen Abschnitt. Die Form ist ein
Shell-Werkzeug mit `awk` auf dem Host und nicht ein Go-Test, weil der Build-Kontext des Stage
`test` weder `spec/` noch `docs/` enthält.
