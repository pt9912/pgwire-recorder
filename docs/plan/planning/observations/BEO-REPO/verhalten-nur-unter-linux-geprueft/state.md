**Stand:** offen

`LH-FA-07.a` Schritt 2 nennt die Grenze: geprüft ist das Verschieben nur unter Linux im selben Dateisystem, unter macOS und Windows nicht (Entscheidung des Nutzers, `36b9d56`; keine Schärfung je Plattform). Dieselbe Grenze nennt `LH-FA-22.a` *Schreiben von `sqlite`* unter *Grenzen* (Verschieben wie in `LH-FA-07.a` Schritt 2; Transaktionen nur unter Linux auf einem lokalen Dateisystem; Stromausfall nicht belegt). Für `sqlite` ist die Bibliothek auf allen sechs Plattformen gebaut, ihr Verhalten nur unter `linux/amd64` geprüft; der Windows-Zweig der URI-Bildung hat keine Zusage (`evidence/slice-v1-abschluss-sqlite-format.md`). Das Lesen übernimmt die Grenze in `slice-v1-abschluss-sqlite-format-lesen` §6.

Weiter offen, kein Auftreten: `slice-v1-abschluss-einspielen-tls-doku` — alle Proben des Handbuchs zu TLS laufen unter Linux im Container; das Handbuch sagt nichts über eine andere Plattform zu.

3× erreicht mit `slice-v1-abschluss-sqlite-format-lesen`. Der Ausgang ist nicht zugewiesen: Ob die Spezifikation die Zusage je Plattform schärft oder als Grenze nennt, steht aus (Architect mit dem Nutzer), und kein Slice ist Adresse für eine Regel dazu. Der Lese-Schritt der Closure von welle-v1-abschluss weist ihn zu.

Weiter offen, kein neues Auftreten: `slice-v1-abschluss-sqlite-format-doku` — alle Proben des Handbuchs zu `sqlite` laufen unter Linux im Container auf einem lokalen Dateisystem; Review und Verifikation haben dazu keinen Befund (das Handbuch macht zu Windows, macOS und Netzwerk-Dateisystemen keine Aussage). Der Zähler bleibt bei 3.

Weiter offen, kein neues Auftreten: `slice-v1-abschluss-zeitangaben-doku` — alle Proben des Handbuchs zu den Zeitangaben laufen unter Linux im Container; das Handbuch macht zu Windows, macOS und Netzwerk-Dateisystemen keine Aussage, Review und Verifikation haben dazu keinen Befund. Der Zähler bleibt bei 3.

Weiter offen, kein neues Auftreten: `slice-v1-abschluss-antwortvergleich-doku` — alle Proben des Handbuchs zum Antwortvergleich laufen unter Linux im Container; das Handbuch macht zu Windows, macOS und Netzwerk-Dateisystemen keine Aussage, Review und Verifikation haben dazu keinen Befund. Der Zähler bleibt bei 3.
