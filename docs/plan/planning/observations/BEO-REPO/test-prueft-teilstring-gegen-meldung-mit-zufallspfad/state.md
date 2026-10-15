**Stand:** offen

Erstes Auftreten, unter 3×. Behoben im Slice (`ohnePfad` streicht jeden Pfad unter `os.TempDir()` aus der Meldung, bevor gegen Werte geprüft wird); die Verifikation reproduzierte die Ursache mit `TMPDIR=/tmp/t4709x` am Test der Basis (rot) und am geänderten Test (grün). Ob andere Tests dieselbe Form tragen (Teilstring gegen eine Meldung mit Pfad), ist nicht gesucht; ohne Adresse, Vorschlag an den Nutzer: ein Satz in `.claude/commands/implement-slice.md` Schritt 19 (Meldung vor dem Teilstring-Vergleich von Pfaden befreien), erst ab 3×.
