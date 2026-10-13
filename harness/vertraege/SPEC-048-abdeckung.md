# SPEC-048 — Abdeckung je Anforderung und Pfad (`abdeckung`)

`make abdeckung` schreibt die Abdeckungstabellen `docs/user/abdeckung-*.md` aus den
Abdeckungs-Deklarationen der Tests; `make abdeckung-check` prüft, dass sie dem
erzeugten Stand entsprechen, und schreibt nichts. Beide wenden dieselben Regeln an.

1. **Deklaration.** Gelesen werden die Testdateien (Name auf `_test.go`) im Repo außer
   unter den Verzeichnissen `.git/` und `.harness/` an der Wurzel des Repos; ein
   gleichnamiges Verzeichnis tiefer im Baum wird gelesen. In einer gelesenen
   Testdatei steht direkt über `func Test…`, ohne Leerzeile und am Zeilenanfang:

   ```text
   // Abdeckung: <Anforderung>/<Pfad>, … — <Kurzbeschreibung>
   ```

   Folgezeilen mit `//` setzen sie fort und werden mit einem Leerzeichen angefügt.
   Fehlformen sind: eine eingerückte Deklaration, eine Deklaration ohne „ — “ und eine
   Deklaration, auf die nicht unmittelbar ein Test folgt (Leerzeile, andere Funktion,
   weitere Deklaration, Dateiende).
2. **Anforderung und Pfad.** Der Pfad ist bei funktionalen Anforderungen (`LH-FA`)
   eines der Akzeptanzkriterien `Happy`, `Boundary` oder `Negative`, bei
   Qualitätsanforderungen (`LH-QA`) und Randbedingungen (`LH-RB`) `Messung`; jede
   andere Paarung ist eine Fehlform. Eine Anforderung, die das Lastenheft nicht als
   Überschrift dritter Ebene führt (`### <Anforderung> …`), ist eine Fehlform.
3. **Nachweisart nach Ort.** Ein Test in einer Datei unter `test/integration/` ist ein
   E2E-Nachweis, jeder andere ein Unit-Nachweis. Jeder gelesene Test, dessen Name mit
   `TestE2E` beginnt, trägt eine Deklaration, gleich in welchem Verzeichnis er liegt;
   fehlt sie, ist das eine Fehlform.
4. **Tabellen.** Unter `docs/user/` entstehen vier Tabellen:
   - `docs/user/abdeckung-e2e.md` und `docs/user/abdeckung-unit.md`: je Deklaration
     und Paar aus Anforderung und Pfad eine Zeile mit Anforderung (Link auf das
     Lastenheft), Pfad, Kurzbeschreibung, Test und Datei, sortiert nach Anforderung,
     Pfad und Test; ein `|` in der Kurzbeschreibung ist maskiert.
   - `docs/user/abdeckung-gesamt.md`: je Anforderung mit mindestens einer
     Deklaration eine Zeile mit den Spalten `Happy`, `Boundary`, `Negative` und
     `Messung` — je die Nachweisarten, die den Pfad belegen, durch Komma getrennt,
     `—` für unbelegt, `n/a` für nicht anwendbar — und dem Stand `vollständig` oder
     `teilweise`.
   - `docs/user/abdeckung-vollstaendig.md`: nur die vollständig belegten
     Anforderungen, je mit den beteiligten Nachweisarten. Vollständig ist eine
     `LH-FA` mit `Happy`, `Boundary` und `Negative`, gleich aus welcher Nachweisart,
     eine `LH-QA` oder `LH-RB` mit `Messung`; Teilabdeckung steht hier nicht.
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
   Tabelle. Fehlt das Lastenheft, endet sie mit 2 und einer Meldung auf stderr.

**Grenze.** Ob eine Deklaration ihren Pfad tatsächlich belegt, prüft das Werkzeug
nicht: Es liest Deklarationen, es führt keine Tests aus. Führt das Lastenheft keine
Überschrift `### LH-…`, endet die Prüfung mit 1 ohne eine Zeile auf stderr.
