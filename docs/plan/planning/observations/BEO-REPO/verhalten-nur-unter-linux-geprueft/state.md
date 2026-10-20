**Stand:** verkörpert

Zielorte: `AGENTS.md §3.12` (die Plattform gehört zu den Randformen vor dem Code) und `.claude/agents/architect.md` — Herkunfts-Anker `seit welle-v1-abschluss` (Entscheidung des Nutzers vom 2026-10-18, Dialogantwort „Beide verkörpern“; keine Querschnittszeile in der Spezifikation, die Grenze steht an der Zusage). Die Regel sagt, dass die Grenze genannt wird; sie macht keine Plattform geprüft. Ein Sensor wäre ein Lauf auf macOS und Windows, den ein Repo mit den Host-Abhängigkeiten `bash`, `git` und `docker` nicht hat.

`LH-FA-07.a` Schritt 2 nennt die Grenze: geprüft ist das Verschieben nur unter Linux im selben Dateisystem, unter macOS und Windows nicht (Entscheidung des Nutzers, `36b9d56`; keine Schärfung je Plattform). Dieselbe Grenze nennt `LH-FA-22.a` *Schreiben von `sqlite`* unter *Grenzen* (Verschieben wie in `LH-FA-07.a` Schritt 2; Transaktionen nur unter Linux auf einem lokalen Dateisystem; Stromausfall nicht belegt). Für `sqlite` ist die Bibliothek auf allen sechs Plattformen gebaut, ihr Verhalten nur unter `linux/amd64` geprüft; der Windows-Zweig der URI-Bildung hat keine Zusage (`evidence/slice-v1-abschluss-sqlite-format.md`). Das Lesen übernimmt die Grenze in `slice-v1-abschluss-sqlite-format-lesen` §6.

Weiter offen, kein Auftreten: `slice-v1-abschluss-einspielen-tls-doku` — alle Proben des Handbuchs zu TLS laufen unter Linux im Container; das Handbuch sagt nichts über eine andere Plattform zu.

3× erreicht mit `slice-v1-abschluss-sqlite-format-lesen`. Der Ausgang ist nicht zugewiesen: Ob die Spezifikation die Zusage je Plattform schärft oder als Grenze nennt, steht aus (Architect mit dem Nutzer), und kein Slice ist Adresse für eine Regel dazu. Der Lese-Schritt der Closure von welle-v1-abschluss weist ihn zu.

Weiter offen, kein neues Auftreten: `slice-v1-abschluss-sqlite-format-doku` — alle Proben des Handbuchs zu `sqlite` laufen unter Linux im Container auf einem lokalen Dateisystem; Review und Verifikation haben dazu keinen Befund (das Handbuch macht zu Windows, macOS und Netzwerk-Dateisystemen keine Aussage). Der Zähler bleibt bei 3.

Weiter offen, kein neues Auftreten: `slice-v1-abschluss-zeitangaben-doku` — alle Proben des Handbuchs zu den Zeitangaben laufen unter Linux im Container; das Handbuch macht zu Windows, macOS und Netzwerk-Dateisystemen keine Aussage, Review und Verifikation haben dazu keinen Befund. Der Zähler bleibt bei 3.

Weiter offen, kein neues Auftreten: `slice-v1-abschluss-antwortvergleich-doku` — alle Proben des Handbuchs zum Antwortvergleich laufen unter Linux im Container; das Handbuch macht zu Windows, macOS und Netzwerk-Dateisystemen keine Aussage, Review und Verifikation haben dazu keinen Befund. Der Zähler bleibt bei 3.

Weiter offen, kein neues Auftreten: `slice-v1-abschluss-sessions-doku` — alle Proben des Handbuchs zu mehreren Verbindungen laufen unter Linux im Container; das Handbuch macht zu Windows, macOS und Netzwerk-Dateisystemen keine Aussage, Review und Verifikation haben dazu keinen Befund. Der Zähler bleibt bei 3.

Weiter offen, kein neues Auftreten: `slice-v1-abschluss-tls-client-doku` — alle Proben des Handbuchs zu TLS zum Client laufen unter Linux im Container (Implementer und Verifier); das Handbuch macht zu Windows, macOS und Netzwerk-Dateisystemen keine Aussage, Review und Verifikation haben dazu keinen Befund. Der Zähler bleibt bei 3.

Weiter offen, kein neues Auftreten: `slice-v1-abschluss-anmeldung-doku` — alle Proben des Handbuchs zur Anmeldung beim Aufzeichnen laufen unter Linux im Container (Implementer und Verifier); das Handbuch macht zu Windows und macOS keine Aussage, Review und Verifikation haben dazu keinen Befund. Der Zähler bleibt bei 3.

Weiter offen, kein neues Auftreten: `slice-v1-abschluss-protokollrand-doku` — alle Proben des Handbuchs zum Protokollrand (P1 bis P16) laufen als Prozess unter Linux; das Handbuch sagt „unter Linux geprüft; für andere Systeme nicht geprüft“ (ohne „im Container“, Review F-2) und macht zu Windows und macOS keine Aussage, Review und Verifikation haben dazu keinen Befund. Der Zähler bleibt bei 3; der Ausgang ist weiter nicht zugewiesen und gehört dem Lese-Schritt der Closure von welle-v1-abschluss.

Kein Auftreten, Gegenbeleg (`slice-v1-abschluss-cancel-ohne-schluessel`, Urteil des Planners, 2026-10-20): §6 nennt die Plattform vor dem Code (Linux im Container geprüft, macOS und Windows als Grenze an der Zusage der Zustellung der `ErrorResponse`), und die Grenze steht im Punkt der Spezifikation (Verifikation, Abschnitt „Plan-vs-Code“; Nach-Verifikation Abschnitt 5). Kein Evidence, der Zähler bleibt bei 3.
