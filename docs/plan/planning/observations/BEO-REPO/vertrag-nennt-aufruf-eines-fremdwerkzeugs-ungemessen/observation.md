# Ein Vertrag sagt den Aufruf eines Fremdwerkzeugs zu, den erst der Trockenlauf misst

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der Vertrag eines Harness-Werkzeugs (hier `SPEC-067` Punkt 16, Scan mit Trivy) legt Aufruf und Wirkung eines Fremdwerkzeugs fest (Layout als Tar, `--platform`, Ort der Datenbank-Metadaten) und führt die Annahmen darüber als „nicht gemessen“; der Code ist so geschrieben, dass eine Abweichung rot endet, nie grün. Eine der Annahmen ist falsch, und erst der erste Lauf auf GitHub zeigt es: Architect, Implementer, Review und Verifikation haben bis dahin Vertrag, Code, Fälle und Mutationen auf ihr gebaut, danach folgt eine Behebungsrunde. Die Messung am Host mit dem gepinnten Image war in Minuten möglich und lief erst nach dem roten Lauf.

Abgrenzung zu `BEO-REPO/plan-annahme-geruest-nicht-gemessen`: Dort nennt ein Plan Aufrufe nach dem Auftrag und ein Nehmer übernimmt sie; hier trägt der Vertrag des Architect die Annahme und der Code folgt ihm. Abgrenzung zu `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`: Dort entscheidet der Implementer eine Randform; hier ist sie entschieden, aber auf einer ungemessenen Annahme.
