# Die Gegenprobe des Stop-Hooks prüft nicht die Umgebung, in der Claude Code ihn liest

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`make stop-hook-gegenprobe` fährt den Hook in Temp-Repos mit dessen Eingabe und liest seine Ausgabe. Ob die
Umgebung von Claude Code die Antwort (`decision`) so liest, prüft sie nicht; ein Hook, der in der
Gegenprobe richtig antwortet und vom Treiber anders gelesen wird, bleibt ungeprüft. Der Hook ist das
Restrisiko von [ADR-0038](../../../../adr/0038-gate-kette-voll-und-schnell.md); die Beobachtung hält fest,
wie oft dieselbe Grenze in einem Plan als offen stehen bleibt.
