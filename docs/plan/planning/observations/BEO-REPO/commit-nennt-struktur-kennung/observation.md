# Eine Commit-Message nennt eine Struktur-Kennung der Spezifikation

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Commit-Messages nennen `SPEC-*`- oder `ARC-*`-Kennungen, obwohl `AGENTS.md` §5 Regel 1
Struktur-IDs aus der Commit-Message ausschließt. Der Commit-Hook prüft nur, ob eine
zugelassene Kennung vorhanden ist, nicht, ob eine ausgeschlossene dasteht; das Review
liest die Messages und sieht es nicht.
