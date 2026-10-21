**Stand:** verkörpert

Zielort: `AGENTS.md §3.14` (jede Rückführungs-Marke in §4 nennt ihre Metrik: Summe aus Einfügungen und Löschungen über die genannten Pfade, Tests eingeschlossen) — Herkunfts-Anker `seit slice-harness-gate-kette-stop-hook-liste` (Beleg `evidence/slice-harness-gate-kette-stop-hook-liste.md`; Entscheidung des Nutzers vom 2026-10-13 „Regel in AGENTS.md aufnehmen“).

Weiter offen, kein Auftreten (Urteil des Planners, `slice-harness-gate-kette-schnell`): §4 nannte die Grenze „über 650 Zeilen mit Tests“ ohne Metrik. Unter jeder Lesart von Code und Spezifikation (406 ohne `spec/`, 563 mit) liegt der Diff unter der Grenze; nur mit Plänen und Reviews läge er darüber, und die zählt niemand als „mit Tests“. Die Probe hing nicht an der Lesart. Der Planner legte die Metrik bei der Closure in §4 fest.

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-harness-devcontainer-einrichten-geruest` — die Metrik nennt die Pfade, nicht ob die Vertragsdateien des Architect unter einem genannten Pfad zählen; 1062 gegen 783 (Beleg `evidence/slice-harness-devcontainer-einrichten-geruest.md`). Der Planner legte die Lesart in §4 des Plans fest; die Regel in `AGENTS.md` §3.14 ist nicht geschärft, die Lücke ist benannt.

Weiteres Auftreten, Regel bleibt (Zielort wie oben): `slice-harness-ci-workflows` — die Basis der Metrik liegt vor einem Commit des Architect, der den Stand der Entscheidung setzt; zwei Lesarten (899, 891), beide unter der Marke (Beleg `evidence/slice-harness-ci-workflows.md`, V-5). Die Regel in `AGENTS.md` §3.14 ist nicht geschärft, die Lücke ist benannt. Zähler 5.

Kein Auftreten, Gegenbeleg (`slice-harness-devcontainer-einrichten`, Urteil des Planners, 2026-10-21): §4 nannte Metrik als Kommando, Basis als Hash (`28bb1f01`) und Pfade; Reviewer (1164 an `3ca372dc`), Verifier (1341 an `d179df23`), Nachverifikation und Planner (je 1407 an `2c73122a`) fuhren dasselbe Kommando und lasen jede Zahl am genannten Stand. Eine Lesart wurde nicht gebraucht. Kein Evidence, der Zähler bleibt bei 5.
