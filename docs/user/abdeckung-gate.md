# Gate-Abdeckung je Anforderung und Pfad

Erzeugt von `make abdeckung` (`tools/test/abdeckung.sh`) aus den
Abdeckungs-Deklarationen im Kopf der Gegenproben unter `tools/`. Die Datei ist eine
Abdeckungs-Deklaration, kein Lauf-Beleg; dass die Gegenproben laufen, sichert
`make gates`. Die Gesamtsicht steht in `abdeckung-gesamt.md`.

| Anforderung | Pfad | Kurzbeschreibung | Nachweis |
| --- | --- | --- | --- |
| [`LH-QA-07`](../../spec/lastenheft.md) | Messung-1-von-3 | Teil 1, statische Analyse nach dem Regelprofil: der Grundlauf ohne Befund (p0-grundlauf), die Linter und ihre Schwellen (p3-*, p4-*, p5-*), kein `//nolint` (p6-*) und die Ausnahmen mit Begründung im Profil (p8-*) | `tools/harness/lint-gegenprobe.sh` |
| [`LH-QA-07`](../../spec/lastenheft.md) | Messung-3-von-3 | Teil 3, Lage der Unit-Tests: ein Test in der Einheit wird rot (p7-whitebox-*), ebenso ein Zugriff auf Internes an der Export-Test-Brücke vorbei (p7-internal-test, p7-endet-auf-export-test, p7-test-in-bruecke-*) | `tools/harness/lint-gegenprobe.sh` |
