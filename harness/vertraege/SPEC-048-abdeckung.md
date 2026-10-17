# SPEC-048 — Abdeckung je Anforderung und Pfad (`abdeckung`)

`make abdeckung` schreibt die Abdeckungstabellen `docs/user/abdeckung-*.md` aus den
Abdeckungs-Deklarationen der Tests und, für Qualitätsanforderungen und Randbedingungen, der
Gegenproben der Gates; `make abdeckung-check` prüft, dass sie dem erzeugten Stand entsprechen, und
schreibt nichts. Beide wenden dieselben Regeln an. Die Nachweisart *Gate* und die geteilte
Messung (Punkte 7 und 8) stehen in [ADR-0033](../../docs/plan/adr/0033-gate-nachweise-in-der-abdeckung.md).

1. **Deklaration.** Gelesen werden die Testdateien (Name auf `_test.go`) im Repo außer
   unter den Verzeichnissen `.git/` und `.harness/` an der Wurzel des Repos; ein
   gleichnamiges Verzeichnis tiefer im Baum wird gelesen. In einer gelesenen
   Testdatei steht direkt über `func Test…`, ohne Leerzeile und am Zeilenanfang:

   ```text
   // Abdeckung: <Anforderung>/<Pfad>, … — <Kurzbeschreibung>
   ```

   Folgezeilen mit `//` setzen sie fort und werden mit einem Leerzeichen angefügt.
   Fehlformen sind: eine eingerückte Deklaration (Leerzeichen oder Tab), eine Deklaration ohne „ — “ und eine
   Deklaration, auf die nicht unmittelbar ein Test folgt (Leerzeile, andere Funktion,
   weitere Deklaration, Dateiende).
2. **Anforderung und Pfad.** Der Pfad ist bei funktionalen Anforderungen (`LH-FA`)
   eines der Akzeptanzkriterien `Happy`, `Boundary` oder `Negative`, bei
   Qualitätsanforderungen (`LH-QA`) und Randbedingungen (`LH-RB`) `Messung` oder
   `Messung-<i>-von-<n>` (Punkt 8); jede andere Paarung ist eine Fehlform. Eine Anforderung, die das Lastenheft nicht als
   Überschrift dritter Ebene führt (`### <Anforderung> …`), ist eine Fehlform.
3. **Nachweisart nach Ort.** Ein Test in einer Datei unter `test/integration/` an der Wurzel
   des Repos ist ein E2E-Nachweis, jeder andere ein Unit-Nachweis; ein gleichnamiges
   Verzeichnis tiefer im Baum macht den Test zum Unit-Nachweis. Eine Deklaration im Kopf einer
   Gegenprobe (Punkt 7) ist ein Gate-Nachweis. Jeder gelesene Test, dessen Name mit
   `TestE2E` beginnt, trägt eine Deklaration, gleich in welchem Verzeichnis er liegt;
   fehlt sie, ist das eine Fehlform.
4. **Tabellen.** Der Wortlaut der Kopftexte ist nicht zugesagt. Unter `docs/user/` entstehen
   fünf Tabellen:
   - `docs/user/abdeckung-e2e.md` und `docs/user/abdeckung-unit.md`: je Deklaration
     und Paar aus Anforderung und Pfad eine Zeile mit Anforderung (Link auf das
     Lastenheft), Pfad, Kurzbeschreibung, Test und Datei, sortiert nach Anforderung,
     Pfad und Test; ein `|` in der Kurzbeschreibung ist maskiert.
   - `docs/user/abdeckung-gate.md`: je Deklaration und Paar aus Anforderung und Pfad eine
     Zeile mit Anforderung (Link auf das Lastenheft), Pfad, Kurzbeschreibung und Nachweis (der
     Pfad der Gegenprobe), sortiert nach Anforderung, Pfad (ein Teil nach seiner Nummer) und
     Nachweis; dieselbe Maskierung. Ohne Gate-Deklaration trägt sie nur Kopf und Tabellenkopf.
   - `docs/user/abdeckung-gesamt.md`: je Anforderung mit mindestens einer
     Deklaration eine Zeile mit den Spalten `Happy`, `Boundary`, `Negative` und
     `Messung` — je die Nachweisarten, die den Pfad belegen, durch Komma getrennt,
     `—` für unbelegt, `n/a` für nicht anwendbar — und dem Stand `vollständig` oder
     `teilweise`. Die Nachweisarten stehen in der Folge `E2E`, `Gate`, `Unit`. Bei einer
     geteilten Messung (Punkt 8) zeigt die Spalte `Messung` alle Teile von 1 bis n, aufsteigend,
     durch „; “ getrennt, je Teil `<i>/<n> <Nachweisarten>` oder `<i>/<n> —` für einen
     unbelegten Teil.
   - `docs/user/abdeckung-vollstaendig.md`: nur die vollständig belegten
     Anforderungen, je mit den beteiligten Nachweisarten. Vollständig ist eine
     `LH-FA` mit `Happy`, `Boundary` und `Negative`, gleich aus welcher Nachweisart,
     eine `LH-QA` oder `LH-RB` mit `Messung` oder mit allen Teilen 1 bis n einer geteilten
     Messung, gleich aus welcher Nachweisart; Teilabdeckung steht hier nicht. Ein Teil in
     zwei Nachweisarten zählt einmal.
5. **Schreiben und Prüfen.** `make abdeckung` schreibt nur die Tabellen, die vom
   erzeugten Stand abweichen oder fehlen, setzt deren Rechte auf 0644 und meldet je
   geschriebener Tabelle eine Zeile `abdeckung: <pfad> geschrieben` auf stdout; eine
   aktuelle Tabelle bleibt unberührt, auch in ihren Rechten. `make abdeckung-check`
   vergleicht nur und meldet je veralteter Tabelle eine Zeile mit ihrem Pfad auf
   stderr.
6. **Ausgang.** Die Prüfung endet mit 0 ohne Fehlform und, beim Prüfen, ohne
   veraltete Tabelle. Sie endet mit 1 bei mindestens einer Fehlform, je mit einer
   Zeile auf stderr, die Datei und Zeile nennt (bei einem Test `TestE2E…` ohne
   Deklaration die Zeile des Tests), oder beim Prüfen bei mindestens einer veralteten
   Tabelle. Sie endet mit 2 und einer Meldung auf stderr, die nennt, was fehlt, und
   prüft dann nichts, wenn (in dieser Reihenfolge) das Lastenheft
   fehlt, das Lastenheft keine Zeile führt, die mit `### LH-` beginnt, oder das Verzeichnis `docs/user/`
   fehlt; `make abdeckung` legt es nicht an. Ob eine solche Überschrift die Form einer Anforderung
   hat, prüft dieser Ausgang nicht: Ein Lastenheft nur mit `### LH-XX-01` besteht ihn, und jede
   Deklaration scheitert dann an (4) mit Exit 1 (akzeptiertes Negativ: ohne Deklaration bleiben die
   Tabellen leer, und ein Lastenheft ohne Anforderung fällt beim ersten Test auf). **Reihenfolge der Gründe** je Deklaration und
   Element, ein Element nennt höchstens einen: (1) Form der Deklaration (Punkte 1 und 7),
   (2) Schreibweise von Anforderung und Pfad einschließlich der Paarung und der Zahlform
   (Punkt 2, 8), (3) Wertebereich von i und n, (4) die Anforderung steht nicht im Lastenheft,
   (5) Widerspruch zwischen Deklarationen (Punkt 8). Die Prüfung (5) betrachtet nur Elemente
   ohne einen Grund aus (1) bis (4). Ein Fehler aus (5) steht an der späteren Fundstelle,
   geordnet nach Dateipfad (bytegenau) und Zeile, und nennt die frühere: bei verschiedenen n die
   erste Deklaration (in dieser Ordnung), die das n der Anforderung festlegt, bei `Messung` neben
   einer geteilten Messung die erste Deklaration des jeweils anderen Pfads; ein Element, das
   widerspricht, legt nichts fest, und jedes weitere Element mit dem widersprechenden n steht
   ebenfalls im Fehler. Die Zeile ist die Zeile der Deklaration, bei einer Deklaration im Kopf
   einer Gegenprobe die Zeile ihres Anfangs `# Abdeckung: …`, bei einer eingerückten Deklaration
   und einer Deklaration nach dem Kopf die Zeile selbst. Der Wortlaut der Meldung ist nicht
   zugesagt, nur Datei, Zeile und der Grund nach der Reihenfolge.
7. **Deklaration im Kopf einer Gegenprobe (Nachweisart Gate).** Gelesen werden die regulären Dateien
   (ein Symlink nicht), deren Name auf `-gegenprobe.sh` endet, in `tools/` an der Wurzel des Repos in jeder
   Tiefe; eine Zeile `# Abdeckung: …` in jeder anderen Datei wird nicht gelesen (eine
   Deklaration dort fehlt in den Tabellen, die Anforderung steht dann höchstens als
   `teilweise`; akzeptiertes Negativ). Der Kopf einer Datei sind ihre ersten zusammenhängenden
   Zeilen, die mit `#` beginnen (die Shebang-Zeile eingeschlossen), bis zur ersten Zeile, die
   nicht mit `#` beginnt. Im Kopf steht am Zeilenanfang:

   ```text
   # Abdeckung: <Anforderung>/<Pfad>, … — <Kurzbeschreibung>
   ```

   Folgezeilen, die mit `# ` und Text beginnen, setzen sie fort (mit einem Leerzeichen
   angefügt); sie endet an einer Zeile aus `#` und höchstens Leerzeichen, an jeder anderen Zeile des
   Kopfs, die nicht mit `# ` beginnt, an der nächsten Deklaration und am Ende des Kopfs. Mehrere Deklarationen in einem Kopf sind zulässig. Fehlformen sind: eine
   eingerückte Deklaration (Leerzeichen oder Tab vor `# Abdeckung:`, an jeder Stelle der Datei, nicht
   nur im Kopf; anderer Leerraum wird nicht geprüft), eine Deklaration ohne „ — “, eine Zeile
   `# Abdeckung: ` am Zeilenanfang nach dem Kopf und eine Anforderung `LH-FA` (eine
   Gate-Deklaration belegt kein Produktverhalten). Eine Gegenprobe ohne Deklaration erzeugt
   weder Fehler noch Zeile.
8. **Geteilte Messung.** Der Pfad `Messung-<i>-von-<n>` nennt den Teil i einer Messung mit n
   Teilen, bei `LH-QA` und `LH-RB`, in der Deklaration eines Tests wie in der einer Gegenprobe.
   i und n sind einstellige Dezimalzahlen ohne Vorzeichen, Leerraum und führende Null
   (`[1-9]`); es gilt i ≤ n und 2 ≤ n. Fehlformen sind: eine andere Schreibweise (mehrstellig,
   `01`, Vorzeichen, Leerraum), i > n, n = 1 (verletzt ein Element beide, i > n und n = 1, nennt die Meldung einen der beiden Gründe; welchen,
   ist nicht zugesagt), verschiedene n für dieselbe Anforderung und
   `Messung` neben einer geteilten Messung derselben Anforderung, gleich in welcher Datei.
   Derselbe Teil in zwei Deklarationen ist zulässig und ergibt zwei Zeilen wie zwei Tests am
   selben Pfad. Vollständig ist eine geteilte Messung, wenn jeder Teil 1 bis n in mindestens
   einer Deklaration steht (Punkt 4).

**Grenze.** Ob eine Deklaration ihren Pfad tatsächlich belegt, prüft das Werkzeug
nicht: Es liest Deklarationen, es führt keine Tests und keine Gegenproben aus. Es prüft auch
nicht, ob eine Gegenprobe an der Kette der Gates hängt, ob die Gesamtzahl n der Teile
stimmt (nur, dass sie für eine Anforderung einheitlich ist) und ob ein Teil zur
richtigen Messmethode der Anforderung gehört; das bleibt Review.
