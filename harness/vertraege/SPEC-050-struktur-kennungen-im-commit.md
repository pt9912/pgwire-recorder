# SPEC-050 — Struktur-Kennungen im Commit-Träger (`commit-msg`)

`.githooks/commit-msg` ist der Träger, den git vor jedem Commit mit dem Pfad der Datei
der vorgeschlagenen Message aufruft; Exit ungleich 0 bricht den Commit ab. Aktiviert
wird er mit `make hooks-install`, seine Gegenprobe ist `make hook-gegenprobe`. Diese
Kennung legt fest, wann er eine Message wegen einer Struktur-Kennung ablehnt. Was er
sonst annimmt und was er an die mitgelieferte Prüfung weiterreicht, legt sie nicht
fest.

1. **Ort.** Die Ablehnung liegt im Träger `.githooks/commit-msg` selbst, nicht in der
   mitgelieferten Prüfung `tools/harness/commit-msg-traceability.sh`.
2. **Lesebereich.** Gelesen wird die Message-Datei bis vor die erste Zeile, die genau
   `# ------------------------ >8 ------------------------` lautet (je 24 Bindestriche),
   die Scissors-Zeile von `git commit -v`; eine Zeile in anderer Form beendet den
   Bereich nicht. Nicht gelesen werden Kommentarzeilen: Zeilen, deren erstes Zeichen
   nach führenden Leerzeichen und Tabs `#` ist.
3. **Struktur-Kennung.** `SPEC-NNN` und `ARC-NNN` mit genau drei Ziffern, in genau
   dieser Schreibweise und als ganzes Wort; Wortzeichen sind Buchstabe, Ziffer und
   Unterstrich, jedes andere Zeichen grenzt ab, auch `-`, `/`, `#`, `.`, Backtick,
   Anführungszeichen, `>` und Wagenrücklauf. Die Kennung zählt darum auch in Backticks,
   in einem Zitat, in Klammern, als Linktext und als Teil eines Pfads oder
   Dateinamens. Keine Struktur-Kennung ist eine andere Schreibweise (`spec-049`,
   `Spec-049`, `arc-003`), eine andere Ziffernzahl (`SPEC-49`, `SPEC-0491`) oder ein
   Platzhalter (`SPEC-<NNN>`, `SPEC-NNN`, `ARC-*`); ein Link-Anker wie
   `spezifikation.md#spec-049--lint-profil-lint` ist klein geschrieben und darum keine
   Nennung. Ob die Kennung in Spezifikation oder Sicht existiert, ist ohne Belang. In
   `SPEC-035 bis SPEC-040` sind die beiden Enden je eine Nennung.
4. **Merge und Revert.** Beginnt der Betreff mit `Merge ` oder `Revert `, sucht der
   Träger keine Struktur-Kennung. Betreff ist die erste Zeile des Lesebereichs, die
   weder leer ist noch nur Leerraum trägt noch Kommentarzeile ist, ohne führende
   Leerzeichen und Tabs.
5. **Vorrang.** Die Suche nach Struktur-Kennungen läuft vor jeder Annahme: Eine
   Message mit Struktur-Kennung wird abgelehnt, auch wenn sie daneben eine Kennung
   trägt, die der Träger oder die mitgelieferte Prüfung annimmt. Eine Message ohne
   Struktur-Kennung im Lesebereich, eine nach Punkt 4 ausgenommene und ein Aufruf ohne
   lesbare Message-Datei gehen den Weg, den der Träger ohne diese Kennung geht.
6. **Ausgabe und Ausgang.** Je Nennung eine Zeile auf stderr,
   `commit-msg: Zeile <n>: Struktur-Kennung <Kennung>`, mit `<n>` der Zeilennummer in
   der Message-Datei ab 1, Kommentarzeilen mitgezählt. Eine Kennung, die in derselben
   Zeile mehrmals steht, erscheint für diese Zeile einmal; geordnet wird nach Zeile,
   innerhalb einer Zeile nach erstem Auftreten. Danach folgt genau eine Zeile
   `commit-msg: Struktur-Kennungen (SPEC-NNN, ARC-NNN) gehoeren nicht in die Commit-Message`.
   stdout bleibt leer, der Ausgang ist 1.

**Grenze.** Ein Klon ohne `make hooks-install` und ein Commit mit
`git commit --no-verify` bleiben ungeprüft. Ein anderes Kommentarzeichen als `#`
(`core.commentChar`) kennt der Träger nicht: Eine solche Zeile wird gelesen. Eine selbst
geschriebene Message, deren Betreff mit `Merge ` oder `Revert ` beginnt, geht nach
Punkt 4 ungeprüft durch. Eine Kennung, die ein Zeilenumbruch teilt, ist keine Nennung.
Gelesen wird nur die Message-Datei; was der Commit ändert, liest der Träger nicht.
