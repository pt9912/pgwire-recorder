# Ein Mutant kommt im Docker-Build-Kontext nicht an

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Mutations- oder Testlauf über den Docker-Build-Kontext sieht eine geänderte Datei
nicht, weil BuildKit eine Datei nicht neu überträgt, deren Größe und mtime dem zuletzt
übertragenen Stand desselben Pfads gleichen; `--no-cache` ändert daran nichts. Ein
Mutant bleibt dann zu Unrecht grün, oder ein Zurücksetzen kommt nicht an.
