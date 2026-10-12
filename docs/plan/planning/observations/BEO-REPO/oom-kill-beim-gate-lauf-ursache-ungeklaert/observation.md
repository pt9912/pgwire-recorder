# Ein OOM-Kill beim Gate-Lauf hat eine ungeklärte Ursache

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Beim Gate-Lauf auf dem Entwicklungsrechner (20 Kerne) tötete der Kernel einen Prozess wegen Speichermangels. Welcher Prozess zu welchem Gate gehörte, ist nicht belegt, und eine Senkung der Parallelität einer Gegenprobe behebt ihn möglicherweise nicht. Die Beobachtung ist eine Ursache ohne Befund, keine Zusage und kein Grund für eine Lockerung (`AGENTS.md` §3.6).
