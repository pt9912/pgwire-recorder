**Stand:** offen

Ohne Adresse, unter 3×. Ursache nicht untersucht; ob die Gegenprobe die Disposition selbst
setzen soll oder der Aufruf von `make gates` dokumentiert wird, ist offen.

Zweites Auftreten (Urteil des Planners, `slice-harness-devcontainer-einrichten`, 2026-10-21; Beleg `evidence/slice-harness-devcontainer-einrichten.md`): dieselbe Form bei einer anderen Gegenprobe (Fall `p17-budget-sigint`). Die offene Frage oben (setzt die Gegenprobe die Disposition selbst, oder wird der Aufruf dokumentiert) ist für die Gegenproben entschieden: Der Architect legte in `SPEC-060` Punkt 17 die Regel an den Aufrufer fest, mit gemessenem Rezept; dort steht auch, dass der Hook für die Fahrt der Sammelziele `SIGINT` und `SIGQUIT` zurücksetzt (`BEO-REPO/hook-erbt-deskriptor-und-ignorierte-signale`). Der Zähler steht bei 2, unter der Schwelle, `offen`. Nicht getragen: keine Rollen-Anweisung (`AGENTS.md`, `.claude/commands/`) nennt den Aufruf; ob eine das soll, entscheidet der Nutzer, der Planner trägt keine Sendung ein.
