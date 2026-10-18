**Vorgang:** slice-harness-gegenproben-aus-gates
**Größe:** Sprache von `make` (Locale `de_DE.UTF-8`).
**Fund:** `make gates-schnell-gegenprobe` war am Host unter `LANG=de_DE.UTF-8` rot, im Dev-Container grün: die Namensregel las `make -p` im Text einer Sprache. Behoben mit `LC_ALL=C` im Aufruf von `make` und dem Fall `p2-locale`. Beleg am Ort der Ursache: `BEO-REPO/make-ausgabe-locale-abhaengig-gelesen/evidence/slice-harness-gegenproben-aus-gates.md`. Nachgetragen vom Planner (Closure von `slice-harness-gegenproben-aus-gates-fahrstelle`, 2026-10-18).
