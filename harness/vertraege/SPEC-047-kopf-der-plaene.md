# SPEC-047 — Kopf der Pläne (`kopf-check`)

`make kopf-check` gleicht zwei Mengen von Kennungen ab, ohne über Inhalt zu urteilen:
Was §1 oder §2 eines lebenden Slice-Plans an Kennungen nennt, führt der Kopf des
Plans.

1. **Gegenstand.** Geprüft werden die Dateien, deren Name mit dem Präfix `slice` und
   einem Bindestrich beginnt und auf `.md` endet, flach in
   `docs/plan/planning/open/`, `docs/plan/planning/next/` und
   `docs/plan/planning/in-progress/`. Nicht geprüft werden Unterverzeichnisse dieser
   drei, `docs/plan/planning/done/` (flach und archiviert), Welle-Pläne, Roadmap,
   Register und andere Dateien.
2. **Kennungen.** Gezählt werden `LH-XX-NN`, `LH-XX-NN.x`, `SPEC-NNN` und `ARC-NNN`
   (X ein Großbuchstabe, N eine Ziffer, x ein Kleinbuchstabe) in genau dieser
   Schreibweise und als ganzes Wort; Wortzeichen sind Buchstabe, Ziffer und
   Unterstrich. Ein Punkt beendet das Wort: Folgt einer Hauptkennung nach dem Punkt
   keine Unterkennung (etwa `.ab`), zählt die Hauptkennung. Auszeichnung ist ohne
   Belang: Code-Span, Linktext, Hervorhebung und Codeblock zählen gleich. Link-Anker sind
   klein geschrieben und darum keine Nennung. Eine andere Schreibweise (klein, ohne
   führende Null, mit zusätzlicher Ziffer, ohne Bindestrich) ist keine Kennung, weder
   im Kopf noch in §1 oder §2.
3. **Bereich.** `SPEC-NNN bis SPEC-MMM`, ebenso mit `ARC`, steht im Kopf wie in §1 und
   §2 für jede Kennung von NNN bis MMM. Backticks um die beiden Enden sind erlaubt, und
   ein Zeilenumbruch innerhalb eines Absatzes zählt als Leerraum; in §1 und §2 reicht
   ein Bereich auch über Absatzgrenzen (siehe *Grenze*). Für die Kennungen
   dazwischen steht der Bereich nur, wenn beide Enden ganze Wörter derselben Klasse
   sind und NNN kleiner als MMM ist. Andere Bereichsformen (Gedankenstrich,
   Bindestrich) gibt es nicht.
4. **Gleichheit ist exakt.** Eine Unterkennung im Kopf deckt ihre Hauptkennung nicht,
   eine Hauptkennung keine Unterkennung.
5. **Kopf.** Der Kopf sind die Absätze, die mit `**Bezug:**` und
   `**Berührte Spec-Stellen:**` beginnen, vor der ersten Zeile `## `, je bis zur
   nächsten Leerzeile. Eine Leerzeile ist eine Zeile, die leer ist oder nur
   Leerzeichen und Tabs trägt. Ein Absatz beginnt nur in der ersten Zeile der Datei
   oder nach einer Leerzeile; eine Feldmarke mitten in einer Zeile oder in einer
   Folgezeile (direkt nach `# …`, nach `---` oder nach einem anderen Feld) ist kein
   Feld. Geprüft wird gegen die Vereinigung beider Felder; in welchem Feld eine
   Kennung steht, ist ohne Belang, und ein anderes Feld des Kopfs zählt nicht. `—`
   ist die leere Menge, Mehrfachnennung ist ohne Belang, eine Kopf-Kennung ohne
   Nennung in §1 oder §2 ist kein Befund.
6. **§1 und §2.** Ein Abschnitt reicht von der Zeile `## 1.` bzw. `## 2.`, die
   mitzählt, bis zur nächsten Zeile `## `; erkannt wird er an der Nummer, nicht am
   Titel. Unterüberschriften beenden ihn nicht. Ausgenommen ist der Absatz, der mit
   `Regeln dieser Sektion` beginnt, bis zur nächsten Leerzeile; dieselben Worte
   mitten in einem Absatz nehmen nichts aus. Alles übrige zählt, auch Abgrenzung,
   Herkunft und Bereits-Geliefertes; eine Markierung für eine Nennung ohne Anspruch
   gibt es nicht.
7. **Formfehler.** Fehlt einem geprüften Plan eines der beiden Kopf-Felder oder einer
   der beiden Abschnitte, ist das ein Befund.
8. **Ausgabe und Ausgang.** Je Befund eine Zeile auf stderr,
   `kopf-check: <pfad>: <abschnitt>: <befund>`, mit dem Pfad relativ zur Wurzel des
   Repos; der Abschnitt ist `Kopf`, `§1`, `§2` oder `Datei`. Der Befund lautet
   `<Kennung> fehlt im Kopf`, `Feld Bezug fehlt`, `Feld Berührte Spec-Stellen fehlt`,
   `Abschnitt ## 1. fehlt`, `Abschnitt ## 2. fehlt` oder `nicht lesbar`. Eine Kennung erscheint
   je Abschnitt einmal. Sortiert wird nach Pfad, Abschnitt und Befund in Bytefolge
   (`Kopf` vor `§1` und `§2`); stdout bleibt leer. Ein Plan, der nicht lesbar ist
   oder an dessen Auswertung die Prüfung scheitert, ist genau ein Befund
   `<pfad>: Datei: nicht lesbar`, gleich an welcher Position; die übrigen Pläne
   werden weiter gelesen. Die Prüfung liest alle Pläne und endet dann mit Ausgang 1
   bei mindestens einem Befund, mit 0 ohne Befund und ohne Ausgabe, mit 2 und einer
   Zeile, die `docs/plan/planning fehlt` nennt, wenn die Ablage
   `docs/plan/planning/` fehlt. Ein fehlendes Lifecycle-Verzeichnis enthält keinen
   Plan und ist kein Abbruch.
9. **Start ohne Stufung.** Die Prüfung ist vom ersten Lauf an voll scharf:
   `make kopf-check` hängt an der Gate-Kette von `make gates` und endet mit
   Fehlerstatus, wenn die Prüfung rot ist.

**Grenze.** Ob eine Kennung existiert, in welchem Kopf-Feld sie steht und was
außerhalb von §1 und §2 steht, prüft das Werkzeug nicht. Codeblöcke werden nicht
gesondert verfolgt: Eine Zeile `## ` in einem Codeblock beendet den Abschnitt. In §1
und §2 werden die gezählten Zeilen eines Abschnitts aneinandergehängt; ein Bereich
reicht dort auch über eine Leerzeile und über einen ausgenommenen Absatz
`Regeln dieser Sektion` hinweg. Das macht §1 und §2 nur strenger, nie still grün.
