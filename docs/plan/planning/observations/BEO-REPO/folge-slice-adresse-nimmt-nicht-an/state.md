**Stand:** verkörpert

Zielort: `AGENTS.md §3.13`, `.claude/commands/plan-welle.md` Schritt 6, `.claude/commands/implement-slice.md` Schritte 19 und 24, `.harness/skills/reviewer.md` §Klassifikation (MEDIUM), `.claude/commands/close-welle.md` Schritt 3 Paarung (b) und `.claude/agents/planner.md` — Herkunfts-Anker `seit slice-lint-bestand-kern-driven`.

Geschärft: `AGENTS.md §3.13` *Auch eine Bindung ist eine Sendung* — Herkunfts-Anker `seit slice-v1-abschluss-einspielen-laufsteuerung`.

Weiter offen, kein Auftreten (keine Datei in `evidence/`, der Zähler bleibt): Risiko aus §6 von `slice-v1-abschluss-zeitangaben-uhr` — der Port steht vor seinem ersten Benutzer, und `slice-v1-abschluss-zeitangaben-aufzeichnen` oder `slice-v1-abschluss-zeitangaben` braucht eine andere Signatur als `Jetzt` und `WarteBis`. Das Gegenmittel steht (beide Nehmer nennen die Namen in §1 und DoD, nachgezählt); gelesen wird die Zeile, wenn beide Nehmer geschlossen sind: Ein Plan-Wechsel des Ports ist dann ein Beleg, ohne ihn endet die Zeile mit dieser Feststellung.

Weiter offen, kein neues Auftreten (Risiko aus §6 von `slice-v1-abschluss-zeitangaben-uhr`, Port vor seinem ersten Benutzer): `slice-v1-abschluss-zeitangaben-aufzeichnen` ist geschlossen und benutzte nur `Jetzt` in der gelieferten Signatur, kein Plan-Wechsel des Ports. Gelesen wird die Zeile weiter, wenn `slice-v1-abschluss-zeitangaben` (ruft `WarteBis`) geschlossen ist.

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.13): `slice-v1-abschluss-sessions` (V-308) — eine Nahtstelle, die der Architect im Geber festlegt, stand in den Nehmern nur teilweise (Beleg `evidence/slice-v1-abschluss-sessions.md`).

Geschärft: der Satz zur Nahtstelle des Architect in `.claude/agents/planner.md` (jeder Punkt in §1 und §6 des bindenden Nehmers, mit Nachzählung, bevor der Geber nach `done/` geht; `grep` nach der Kennung des Gebers in jedem Nehmer der Reihe) — Herkunfts-Anker `seit slice-v1-abschluss-sessions`.

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.13): `slice-harness-vertraege-umzug-ordner` (Review F-1) — eine bewusst stehen gelassene Ortsangabe (Zwischenstand) ist eine Sendung an die Nehmer des Umzugs (Beleg `evidence/slice-harness-vertraege-umzug-ordner.md`).

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.13): `slice-v1-abschluss-verbindungsende-einspielen` (F-3, V-2) — eine Bindung des Nehmers an einen Test des Gebers, den der Geber nicht führte; danach nannte der Nehmer für ihn die falsche Stelle (Beleg `evidence/slice-v1-abschluss-verbindungsende-einspielen.md`). Abgegrenzt gegen `BEO-REPO/plan-nennt-test-der-nicht-existiert`: Der Name stand in einer Sendung, nicht im Beleg einer Zusage.

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.13): `slice-v1-abschluss-verbindungsende-einspielen-extended` (F-1, V-1) — eine Bindung an den Nehmer stand nur in §6 und §8, nicht in §1 und DoD; die Kommentare, die die Verifikation dem Nehmer zuwies, standen dort nicht (Beleg `evidence/slice-v1-abschluss-verbindungsende-einspielen-extended.md`). Die Nachzählung in §8 des Nehmers fehlte nicht (sie nannte nur eine nicht vorhandene DoD-Zeile), deshalb kein Zuwachs bei `BEO-REPO/slice-waechst-durch-uebernahmen`.

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.13): `slice-v1-abschluss-verbindungsende-record` (F-2, V-2) — eine Test-Idee nannte einen Nehmer als Adresse, der sie nicht führte und schon an der Grenze von 3 Liefer-Punkten und 2 Schichten lag; zurückgenommen, die Idee gestrichen (Beleg `evidence/slice-v1-abschluss-verbindungsende-record.md`).

Kein Auftreten, Gegenbeleg (`slice-erster-release-freigabe`, Urteil des Planners, 2026-10-19): Die Bindung von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` stand bei der Eröffnung in §1, §6 und §8 des Nehmers; der Befund zum Hook-Lauf geht mit der Closure in einem Commit an den Geber dieser Bindung zurück, mit Nachzählung in dessen §8. Kein Evidence, der Zähler bleibt bei 13.

Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.13): `slice-release-v0-2-0` (F-3 LOW) und `slice-go-toolchain-1-27-3` (F-4 LOW) — in zwei Plänen derselben Sitzung stand ein Punkt bei einer Adresse, die ihn nicht führt (ein geschlossener Slice; ein offener Slice, der ihn ausschließt); beide nachgezogen als „nicht übernommen“, ohne Adresse (ein Beleg für beide Pläne in `evidence/slice-release-v0-2-0.md`). Keine Zeile in `AGENTS.md` oder `.claude/` von der Closure. Zähler 14.
