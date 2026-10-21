# Die enge Auswahl nach einem neuen Fall einer Gegenprobe übersieht die Mutation, die ihn sieht

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Eine Gegenprobe führt je Mutation die Menge der Fälle, die rot werden müssen, und jede Mutation wählt die Fälle, die sie sehen kann, über einen Filter. Wer der Gegenprobe einen Fall hinzufügt, färbt damit jede Mutation rot, deren Filter ihn wählt und die sein Verhalten bricht. Fährt der Implementer danach nur die Mutationen, die er selbst angefasst hat (die enge Auswahl nach `AGENTS.md` §6, wie sie das Zeitbudget der Gegenprobe verlangt), bleibt die Mutation unbemerkt, deren Erwartung den neuen Fall nicht nennt; der volle Lauf des Verifiers endet rot (Nachverifikation N-1, HIGH).

Abgrenzung zu `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: Dort fehlt für eine Zusage der Test oder die Mutation; hier sind beide da, und die Erwartung einer *anderen* Mutation ist durch den neuen Fall falsch geworden. Abgrenzung zu `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig`: Dort ist das zusätzliche Rot zeitabhängig und folgt keiner Änderung der Gegenprobe; hier ist es deterministisch und folgt einer Änderung an ihr.
