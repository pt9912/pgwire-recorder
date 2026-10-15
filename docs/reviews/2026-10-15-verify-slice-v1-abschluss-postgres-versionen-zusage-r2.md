# Verifikation slice-v1-abschluss-postgres-versionen-zusage, Nachlauf r2

**Rolle:** Verifier (Modul 11). **Datum:** 2026-10-15. **Spanne:** `acecaee` (Vorlauf) bis `e537ada` (HEAD). **Gegenstand:** nur das Nachgezogene zu V-1 und V-2 des [Vorlaufs](2026-10-15-verify-slice-v1-abschluss-postgres-versionen-zusage.md); entschieden in `f5ce5d6` (Architect), umgesetzt in `8ccf3c3` (Spezifikation, Lastenheft) und `e537ada` (Plan). Plan: `docs/plan/planning/in-progress/slice-v1-abschluss-postgres-versionen-zusage.md`. Belege selbst erzeugt (Lesen, `grep`, `git diff`, `make kopf-check`); kein Matrix- und kein Integrationslauf, da nur Text geändert ist. Arbeitsbaum unberührt.

**Verdikt: PASS.** V-1 geschlossen, V-2 geschlossen. 0 neue Befunde, 1 INFO.

## V-1 (MEDIUM): geschlossen

| Prüfung | Beleg | Urteil |
|---|---|---|
| `SPEC-029` Punkt 2 ohne Szenario 18 und Szenario 16 Satz 3 | `spec/spezifikation.md` Zeilen 2885 bis 2886: „Gegen die Referenz 17 läuft zusätzlich die volle Integrationssuite. Die Anmeldung mit Klartext-Passwort und mit MD5 sind nur dort geprüft.“ Diff `8ccf3c3` streicht genau die zwei Elemente. | bestätigt |
| Kein anderer Fund in `spec/`, `docs/user/`, `README.md` | `grep -rnE "Szenario 18\|Szenario 16\|Szenarios 18" spec docs/user README.md` außerhalb von `spec/lastenheft.md`: 0 Treffer. Im Lastenheft nur das Szenario selbst (Zeilen 810, 1087, 1101), keine Versionszusage. | bestätigt |
| Sendung in §1 und DoD beider Nehmer mit Kennung des Gebers | `-verbindungsende-upstream` §1 (Absatz „Übernommen aus `slice-v1-abschluss-postgres-versionen-zusage`“) und DoD-Punkt 3; `-verbindungsende-vergleich` ebenso (§1, DoD-Punkt 3). Beide Nehmer liegen in `next/`, keiner in `done/`. | bestätigt |
| Kopf führt `SPEC-029` | beide Köpfe: `SPEC-029` in *Berührte Spec-Stellen*; `make kopf-check` Exit 0 | bestätigt |
| Zählung in §8 der Nehmer | beide: Eintrag „Nachgezählt beim Eintragen der Sendung …“ mit Kennung des Gebers: Liefer-Punkte 3, Schichten 2 (upstream: Upstream-Adapter, Kern nur Port-Vertrag; vergleich: Kern mit Service, CLI-Adapter); die Sendung gehört zu DoD-Punkt 3, Spezifikation und Planung zählen nicht als Schicht (§3.13). Stimmt mit den dort vorhandenen Zählungen überein. | bestätigt |
| Ausschlüsse der Nehmer nehmen an (§3.13) | Beide Listen „Ausdrücklich NICHT“ gelesen: kein Punkt schließt `SPEC-029` oder einen Satz der Spezifikation aus. upstream schließt nur das Handbuch aus (die Sendung nennt das ausdrücklich), vergleich ebenso; die Matrix der Versionen ist in upstream als andere Slice-Aufgabe genannt, Szenario 18 läuft „hier gegen die Referenzversion 17“, deckt sich mit der Sendung. Beide Nehmer schließen nach dem Geber (Geber liegt in `in-progress/`). | bestätigt |
| Kein Zielstand-Satz im Ist-Zustand | Handbuch und README nennen die Szenarien nicht (siehe `grep`) | bestätigt |

## V-2 (LOW): geschlossen

`spec/lastenheft.md` `LH-RB-02` *Nachweis* (Zeile 937): „gegen jede unterstützte Hauptversion läuft mindestens die Teilmenge der Kompatibilitätstests, die die Spezifikation nennt.“ „Mindestens“ stimmt für 14, 15, 16 und 18 (Teilmenge laut `SPEC-029` Punkt 1) und für 17 (volle Suite, Punkt 2); der Quantor „jede … nur“ ist weg. Lastenheft: Status Draft, Version 0.1.0, im Diff `8ccf3c3` kein Versionsfeld und keine Historie-Zeile des Lastenhefts berührt; die Bump-Regel gilt ab `Accepted`.

## Plan, Größe, Gates

- **Plan §6/§7 ohne Reste:** Aussage (g) lautet „Klartext, MD5 nur 17“; Deckungstabelle „Referenz 17“ ohne beide Szenarien. Die verbleibenden Nennungen von Szenario 18 und 16 Satz 3 (§1, §3, §6 Randform, §6 Zusage-reicht-weiter, §7 Nachzug) sind ausdrücklich als nicht zugesagt und als Sendung geführt; keine nennt sie als geprüft. Das Risiko „Zusage reicht weiter als das Gate“ trägt den Ausgang *eingetreten* mit Kandidaten `slice-v1-abschluss-verbindungsende-upstream` und `-vergleich`, „offen bis Closure“; beide Kennungen sind Adressen, die annehmen (siehe V-1).
- **Größe:** `git diff --numstat 510e7ab..8ccf3c3 -- spec docs/user README.md`: README 3, Handbuch 24, Lastenheft 7, Spezifikation 22, Summe **56**, wie behauptet; Marke 400 und Meldeschwelle 300 nicht berührt.
- **`make gates`:** Nachweis `.harness/state/gates-passed.dateiliste`/`.diffsha` vom 05:09:47, nach dem letzten Commit (`e537ada`, 04:59:10); `git diff HEAD` leer. Der Plan nennt zum Nachzug einen ersten roten Lauf (`postgres-versionen-gegenprobe`, Mutation `p7-int-treiber` mit zusätzlich `p7-sigTERM-gruppe`, zeitabhängig), `make test-postgres-versionen-gegenprobe` allein danach grün, zweiter voller Lauf Exit 0. Nicht neu gefahren (Auftrag). `make kopf-check` selbst gefahren: Exit 0.

## INFO

- I-1: Der Satz in `SPEC-029` Punkt 2 „Die Anmeldung mit Klartext-Passwort und mit MD5 sind nur dort geprüft“ ist im Numerus holprig („ist“), ohne Folge für den Inhalt. Beim nächsten Berühren der Stelle (die Nehmer ergänzen dort ohnehin einen Satz) berichtigen.

## Offen für die Closure (nicht Teil dieses Nachlaufs)

DoD-Punkte 6 bis 9 (Closure-Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen) bleiben Handlung des Planners; der Ausgang des Risikos „Zusage reicht weiter als das Gate“ ist bei Closure von „Kandidat“ auf *eingetreten* mit den zwei Kennungen festzuschreiben.
