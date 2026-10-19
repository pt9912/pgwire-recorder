# Der Archiv-Stub einer Welle verliert die Anforderungs-Kennung des Slice

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`archive-welle` (`ai-harness-init`) kürzt den Volltext eines Slice beim Archivieren auf einen Stub aus Identität, Archiv-Zeiger, Zustand und den Kennungen, die den Slice überlebt haben; die Kennungen der Anforderungen (`LH-FA-*`, `LH-QA-*`, `LH-RB-*`) aus dem Kopf des Volltexts gehören nicht dazu. `make doc-trace` liest die Slices unter `done/` rekursiv und zählt eine Anforderung nur, wenn ihre Kennung im Text einer Datei steht: Mit jeder Archivierung kehren Waisen in die Matrix zurück, und kein Gate meldet es, weil `make doc-trace` keines ist.

Die Handlung bis zur Aufnahme upstream ist `make archiv-bezug` nach jeder Archivierung (`SPEC-068`); sie ist eine Handlung, kein Sensor.
