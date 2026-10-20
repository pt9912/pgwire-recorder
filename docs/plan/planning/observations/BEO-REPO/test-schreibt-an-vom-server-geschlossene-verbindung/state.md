**Stand:** offen

Erstes Auftreten: `slice-v1-abschluss-cancel-ohne-schluessel` (Verifikation F-V3, LOW; Beleg `evidence/slice-v1-abschluss-cancel-ohne-schluessel.md`). Im Slice behoben: `cancelSenden` toleriert `ECONNRESET` und `EPIPE` beim Schreiben, fünfmal 100 Läufe der 36 Unterfälle ohne Befund. Unter der Schwelle; `offen`, weil derselbe Fehler jeden Test trifft, der mehr sendet, als der Server vor dem Schließen liest.
