# Eine Gegenprobe kopiert den Arbeitsbaum mit allem, was darin liegt, und bricht, wenn eine Datei sich währenddessen ändert

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Gegenproben, die einen Arbeitsbaum mutieren, legen vorher eine Kopie an: `tools/arch/a-check-negativ.sh` mit `tar --exclude=./.git -cf - . | tar -xf -`, `tools/harness/lint-gegenprobe.sh` und `tools/harness/coverage-gegenprobe.sh` mit `cp -R` über `*` und `.[!.]*` ohne `.git`. Alle drei nehmen jede Datei des Baums mit, auch die ignorierten Laufzeitdateien unter `.harness/state/` (Spans, Stände der Gates), in die die Hooks laufender Agenten schreiben. `tar` bricht dann ab (`tar: ./.harness/state/spans: Datei hat sich beim Lesen geändert`), der volle Lauf von `make gates` endet rot ohne Befund am Code; `cp -R` bricht nicht, nimmt aber einen halb geschriebenen Stand mit.

Abgrenzung zu `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`: Dort blockiert der Hook andere Sitzungen; hier stört die Arbeit der Sitzungen das Gate. Abgrenzung zu `BEO-REPO/mutanten-skript-checkt-kopien-in-den-arbeitsbaum-ein`: Dort landet eine Kopie im Arbeitsbaum; hier ist der Arbeitsbaum die Quelle der Kopie.
