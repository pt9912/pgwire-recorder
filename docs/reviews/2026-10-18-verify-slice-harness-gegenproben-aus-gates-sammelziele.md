# Verifikation: slice-harness-gegenproben-aus-gates-sammelziele

Verifier, 2026-10-18. Stand: Zweig `rev-sz`, Implementer-Stand `16477376`, Basis `09d35ed8`. Plan: [slice-harness-gegenproben-aus-gates-sammelziele](../plan/planning/in-progress/slice-harness-gegenproben-aus-gates-sammelziele.md) (§2 DoD, §7 Belege). Vertrag: [SPEC-060](../../harness/vertraege/SPEC-060-gegenproben.md) Teil 2. Entscheidung: [ADR-0048](../plan/adr/0048-gegenproben-verlassen-die-gate-kette.md). Gefahren am Host, nur die engsten Sensoren; die Läufe des Implementers (`make gates`, `make gegenproben` im Container, `make gegenproben-host`) wurden nicht wiederholt.

## Ergebnis je Punkt

| # | Prüfung | Ergebnis |
|---|---|---|
| 1 | `make gegenproben-gegenprobe` ganz, `LANG=LC_ALL=de_DE.UTF-8` | grün, Exit 0, 27 s, 37 Fälle, 34 Mutationen (`grep -c '^mut '`), Arbeitsbaum danach unverändert |
| 2 | 10 Mutanten einzeln (`GP_MUT`, `GP_ENTDECKE=1`) | alle rot genau wie in der Tabelle aus §7, siehe unten |
| 3 | `make -pn gates` und Auswertung der Listen | `gegenproben` und `gegenproben-host` weder in `GATE_CHECKS` noch im Rezept von `make -n gates`, nicht in `harness/mk/gates-schnell.mk`; `gegenproben-gegenprobe` in `GEGENPROBEN_CONTAINER`; `GEGENPROBEN_HOST` = `container-start-gegenprobe` (Risiko 2 aus §6: entfallen, belegt) |
| 4 | `make gates-schnell-gegenprobe` | grün, 60 s, die Namensregel akzeptiert die neue Gegenprobe |
| 5 | `make docs-check`, `make kopf-check`, `make gates-schnell` | alle grün (0 Befunde; 1106 Dateien) |
| 6 | Probelauf mit rotem Einzelziel (Temp-Verzeichnis mit echtem Skript und Fragment, Ziele `a-gruen`, `b-rot` (Exit 3), `c-gruen`) | alle drei liefen; Zusammenfassung zuletzt auf stderr (`grün`, `rot (Exit 2)`, `grün`); `-n`: nichts gefahren, Exit 0; `-i`: Meldung, kein Ziel; leere Host-Liste: `Liste leer`, Exit 2 |
| 7 | Metrik `git diff --numstat 09d35ed8..HEAD -- tools harness AGENTS.md` | 426 + 2 = **428**, wie beim Implementer; unter Meldeschwelle 700 und Marke 900 (Risiko 3: entfallen) |
| 8 | `AGENTS.md` §6 Schritt 6 | `grep -cF` auf den ganzen Satz des Nutzers = 1: Wortlaut Zeichen für Zeichen (erste Hälfte vom Kettenumbau, Rest von diesem Slice); Plan-vs-Code unten |

Mutanten (Name · rote Fälle): `m-p4-in-gate-checks` · `p4-kein-gate`; `m-p4-host-haengt-an-container` · `p4-host`; `m-p5-prerequisites` · `p5-j, p9-container, p10-i, p10-i-env, p10-i-lang, p10-i-k, p10-i-j, p10-j, p10-k`; `m-p6-exit0` · `p6-erstes-rot, -mittleres-rot, -letztes-rot, -zwei-rot, p6-127, p10-j, p10-k`; `m-p7-stdout-verworfen` · `p7-ausgabe`; `m-p8-ziel-bleibt` · `p8-sigint-skript, p8-sigterm-skript`; `m-p9-exit0` · `p9-container, p9-host`; `m-p10-i-vor-n` · `p10-in, p10-iq, p10-it`; `m-p11-schreibt` · `p11-kein-stand`; `m-p12-liste-vor-i` · `p12-i-leer`. Damit mindestens ein Mutant je Punkt 4 bis 12; jede Menge entspricht §7. Die Gegenprobe selbst verlangt je Mutation genau diese Menge (Ausgabe „gruen").

## Plan-vs-Code

- Dateien aus §3 liegen alle vor (Skript, Fragment, Gegenprobe samt Fragment, `harness/README.md`, `AGENTS.md`); es ist nichts Zusätzliches im Diff (7 Dateien, davon der Plan).
- DoD 1 und 2: Sammelziele und Gegenprobe vorhanden, nicht an `GATE_CHECKS`, nicht in der Ausnahmeliste. DoD 3: Zeilen im Gate-Index vorhanden (zwei Werkzeuge, eine Gegenprobe), Satz in §6 vollständig.
- Die DoD-Häkchen stehen noch offen (Closure folgt dem Review); §7 trägt nur die Belege, Notiz, Risiko-Ausgänge und Paarungen sind „—" (zu erwarten vor der Closure).

## Befunde

1. **Risiko 1 aus §6 ist nicht belegt (offen für die Closure).** Der Container-Lauf von `make gegenproben` folgte auf `make gates` im selben Container (so auch §7); der vorgesehene Ausgang „entfallen" verlangt einen frischen Gegenproben-Container ohne vorheriges `make gates`. Der Implementer nennt das selbst. Der Planner muss entscheiden: Fahrt im frischen Container nachholen, oder Risiko „weiter offen" ins Register (`BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`).
2. **Vertragsabweichungen im Verhalten über `make` (bekannt, Rückgaben 1 und 2 in §7, Vertrag unverändert, Entscheidung beim Architect).** Beobachtet im Probelauf: über `make` endet ein rotes Sammelziel mit Exit 2 (Vertrag Punkt 6: 1), `-i` über `make` mit Exit 0 und Meldung (Punkt 10 (b): 2), leere Liste über `make` mit Exit 2 (stimmt). Das Skript selbst hält 1, 2, 130 und 143; die Gegenprobe prüft über `make` nur Exit ungleich 0. Die Zeile `rot (Exit n)` trägt den Exit von `make <ziel>` (2), nicht den der Gegenprobe. Der Vertrag sagt es noch nicht, bis der Architect ihn nachzieht, stimmt Punkt 6 wörtlich nicht mit dem Verhalten über `make` überein.
3. Kein Befund zu Mutanten, Metrik, Satzwortlaut oder Anschluss an die Ketten. Nicht gefangen und im Vertrag als Grenze benannt: Signal nur an `make`, `SIGKILL` (Punkt 15).

Nicht geprüft (bewusst, vom Implementer belegt): `make gates` im Container, `make gegenproben` im Container (1095 s), `make gegenproben-host` (260 s).

Nichts repariert, nichts committet außer diesem Report.
