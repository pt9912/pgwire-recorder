# Die Laufzeit von `make gates` wächst mit jeder Gegenprobe

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein neues Gate bringt nach `AGENTS.md` §3.10 seine Gegenprobe mit, und eine Gegenprobe,
die das Werkzeug je Fall in einer Kopie laufen lässt, kostet ein Vielfaches des Gates
selbst. Jede weitere hängt an dieselbe Kette `make gates`. Wird der Lauf so lang, dass er
vor dem Handoff übersprungen wird, ist das Gate ein stilles Rot. Die Beobachtung ist
eine Laufzeit mit Zahl, kein Grund für eine Lockerung (`AGENTS.md` §3.6).
