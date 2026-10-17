# Die Zahl der Anweisungsabdeckung wird zum Ziel

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Tests, die Anweisungen ausführen, ohne etwas zu prüfen, heben die Zahl der Anweisungsabdeckung, nicht die Prüfung. Die Schwelle (`SPEC-056`, Gate `make test`) ersetzt `AGENTS.md` §3.10 (je Zusage eine Mutation) nicht; ob ein Test etwas prüft, bleibt Urteil von Review und Verifikation.

## Benannt, nicht gezählt

- 2026-10-17: Risiko in §6 von `slice-harness-coverage`, bei dessen Closure mit dem Ausgang *weiter offen* hierher gestellt. Kein Auftreten beobachtet; die erste Evidence-Datei entsteht, wenn ein Review einen Test findet, der Anweisungen ausführt und nichts prüft, oder eine Anhebung der Schwelle ohne neue Zusage-Tests.
