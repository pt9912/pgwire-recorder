# Slice slice-harness-rollenlauf-zeitgrenze: Ein Rollenlauf wird nach einer Stunde angehalten: Ein Wächter im Werkzeug-Aufruf liest den Strom der Erfassung, sperrt ab der Grenze, was den Lauf verlängert, und lässt die Abwicklung durch — Stufe 1, die Stunde wirkt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-harness-regeln](../welle-harness-regeln.md), **hinter** `slice-harness-rollenlauf-zeitgrenze-regel` (§4 und §5 der Welle) — **Empfehlung des Planners** für Welle und Stelle (angelegt am 2026-10-21, am 2026-10-22 nachgezogen), nicht Entscheidung des Nutzers: Der Nutzer hat die **Sache** entschieden (Bezug), über Welle und Reihenfolge ist nichts entschieden. Angelegt in `open/`; **geschnitten am 2026-10-22** in vier Slices (§1), dieser ist die erste Stufe des Wächters.

**Bezug:** kein Anspruch des Produkts (`LH-*`) — der Slice erzwingt eine Regel über Agentenläufe und ändert kein Verhalten des Werkzeugs. **Entscheidungen des Nutzers** (jeweils vom Koordinator übermittelt): 2026-10-21, wörtlich „Wir haben rollen (siehe .claude/agents) und commands -- max. 1 Stunde“, und „ok“ auf den Entwurf des Mechanismus (Wächter im Werkzeug-Aufruf, bash und awk, **kein** Hintergrundprozess, Datenquelle die vorhandene Erfassung); 2026-10-22, „deine empfehlungen“: die Stunde ist zugleich zweite Rückführungs-Marke und steht in der Regel (`slice-harness-rollenlauf-zeitgrenze-regel`), `AGENTS.md` §3.9 führt `.claude/agents/` auf; 2026-10-22, wörtlich „2. ja“: Der Wächter deckt nur Rollenläufe, und der Regeltext nennt Commands und Hauptsitzung als nicht erzwungen; 2026-10-22: [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) ist **vom Nutzer angenommen**, die Statuszeile (`Accepted`) setzt der Architect, bis zum Merge steht sie auf `Proposed`. Aktive ADRs, gelesen: [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Zeitbudget 600 s je Gegenprobe — gilt auch für die Gegenprobe dieses Slice), [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (die Gegenprobe hängt an `GEGENPROBEN_CONTAINER`, nicht an `GATE_CHECKS`), [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (die Hooks dieses Repos und ihr Zusammenspiel), [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (Metrik der Größenmarke in §4).

**Berührte Spec-Stellen:** [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) (liegt vor, vom Architect; der Vertrag des Wächters mit den Randformen der sechs Fragen dieses Plans; **die Zuordnung seiner Punkte zu den zwei Stufen ist offen**, §6 Randform 7) · [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) (gelesen: Name und Liste der Gegenprobe, Punkte 1 bis 3; Aufnahme nach Punkt 16, deren Zeile der Architect nachträgt) · [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) (gelesen: der Stop-Hook und die Ketten — dieser Slice ändert beide nicht) · [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 (6) (gelesen: Ort des Belegs eines Nachweises, dem der Beleg im Agentenlauf folgt).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-21 (Schnitt und Nachzug nach `SPEC-070`: 2026-10-22).

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Rollenlauf (ein Subagent unter `.claude/agents/`) wird nach einer Stunde angehalten, und das hängt nicht an seiner Disziplin: Ein Wächter im Werkzeug-Aufruf (`PreToolUse`, bash und awk wie `.claude/hooks/pretooluse-command-guard.sh`, kein `jq`, kein `node`) liest den Strom der Erfassung des laufenden Agenten, rechnet die Wanduhr seit dem letzten `SubagentStart` und sperrt ab 3600 s die Werkzeuge `Write`, `Edit`, `NotebookEdit` und `Agent` und die `Bash`-Aufrufe, die `make`, `docker` oder `sleep` nennen. `Read`, `Grep`, `Glob`, `SubagentHandback` und `git commit` und `git push` bleiben frei, damit die Rolle committen, übergeben und berichten kann; die Block-Meldung nennt die Handlung selbst. Das ist **Stufe 1 von zweien: die Stunde wirkt**. Die Härtung und die Selbstauskunft des Hooks sind Stufe 2 (`slice-harness-rollenlauf-zeitgrenze-haertung`). Der Vertrag ist [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) mit der Entscheidung [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md); dieser Plan wiederholt ihn nicht.

**Geschnitten am 2026-10-22 — warum so.** Die Schätzung des Architect für alles zusammen liegt bei 800 bis 1400 (Mitte 1100, Marke 900 in der Metrik von [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md)); die Gegenprobe (550 bis 950 für rund 125 Fälle und 180 Mutationen) ist der Treiber. Erwogen und verworfen:

1. *Regel | Wächter mit Gegenprobe* (zwei Teile): Die Regel ist allein lieferbar und klein (`slice-harness-rollenlauf-zeitgrenze-regel`), aber Hook und Gegenprobe bleiben bei 703 bis 1212 (Mitte 958, 106 % der Marke): `AGENTS.md` §3.10 verlangt je Zusage und je Element einer Aufzählung einen Fall und eine Mutation, und die Zusagen eines Hooks sind eine Lieferung. Als einziger Schnitt reicht das nicht.
2. *Hook | Gegenprobe* in zwei Slices: Ein Hook ohne die Mutationen seiner Zusagen verletzt §3.10; eine Gegenprobe ohne Hook hat nichts zu prüfen.
3. *Leser | Zeit | Sperre* (nach Bauteilen): ein Schicht-Schnitt; jeder Teil ist einzeln nutzlos (Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice).
4. *Sperrmenge in zwei Stufen* (nur `Bash` | dazu die Werkzeuge): trennt rund 7 % der Fälle.
5. **Gewählt: nach Wirkung und Härtung.** Stufe 1 (dieser Slice) liefert, was die Stunde wirken lässt: wer begrenzt ist, Strom, Anker und Zeit, die ganze Sperrmenge, Ausgabe und Exit, kein Weg über die Grenze, den Nachweis im Lauf, die Gegenprobe samt Aufnahme. Stufe 2 liefert, was den Hook an den Rändern hält und seine Lage nennt: strenge Zeitform des Ankers, Reihenfolge der Gründe, Gründe der Freigabe mit der Zeile auf stderr, Umgebung. Geschätzt in der Metrik: Stufe 1 440 bis 759 (Mitte 600, 67 % der Marke 900), Stufe 2 273 bis 493 (Mitte 383, 55 % der Marke 700); beide unter 70 %. **Voraussetzung:** Der Architect ordnet die Punkte von `SPEC-070` den zwei Lieferungen zu (Randform 7 in §6, mit dem Vorschlag des Planners); ohne diese Zuordnung steht der Vertrag als Ganzes, und eine Stufe lieferte Zusagen ohne Fall.

Fehlt die Zuordnung, bleibt die Alternative eine Entscheidung des Nutzers über Marke und Meldeschwelle (Präzedenz: `slice-erster-release-veroeffentlichung`, Marke 1300); der Planner hebt sie nicht.

**Zusagen, die diese Stufe nicht gibt, und warum das trägt:** Der Hook gibt bei jeder Unsicherheit **still** frei (keine Zeile auf stderr, keine Aufzählung der Gründe) und prüft die Zeitform des Ankers nur in der Gestalt, nicht im Wertebereich; beides liefert Stufe 2. Die Richtung ist die sichere (ein verfehlter Block heißt nur, dass die Grenze für diesen Lauf nicht greift, und die Regel an den Leser bleibt, `slice-harness-rollenlauf-zeitgrenze-regel`).

**Sendungen** (eingetragen vom Planner im Commit des Schnitts, `AGENTS.md` §3.13): an `slice-harness-rollenlauf-zeitgrenze-regel` (die Regel in `AGENTS.md` und den sechs Rollen), an `slice-harness-rollenlauf-zeitgrenze-haertung` (die Punkte der zweiten Stufe) und an `slice-harness-rollenlauf-wiederholungsgrenze` (die Wiederholungs-Grenze); je in deren §1 und DoD angenommen, mit Nachzählung in deren §8.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die **Regel** in `AGENTS.md` (die Grenze als harte Regel, die Stunde als zweite Rückführungs-Marke, der Satz zum Verifier, die Zeile in §3.9) und der **Satz in den sechs Rollen-Typen** — Klasse 1: `slice-harness-rollenlauf-zeitgrenze-regel`. Grund: andere Dateien, anderer Vorgang; der Wächter ist ein Werkzeug, die Regel ein Text an den Leser. Die Block-Meldung des Hooks nennt die Handlung selbst (Punkt 6 von [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md)), damit sie nicht an der Regel hängt.
- **Stufe 2** — Härtung und Selbstauskunft — Klasse 1: `slice-harness-rollenlauf-zeitgrenze-haertung`. Grund: Aufteilung des Vertrags nach Wirkung und Härtung, damit jede Lieferung unter 70 % der Marke bleibt (oben).
- Die **Wiederholungs-Grenze für volle Gegenprobe-Läufe** (nach dem ersten vollen Lauf je Rollenlauf wird ein weiterer ohne `GP_NUR` oder `GP_MUT` gesperrt, aus `duration_ms` desselben Stroms) — Klasse 1: `slice-harness-rollenlauf-wiederholungsgrenze` (vom Nutzer am 2026-10-21 freigegeben). Grund: ein zweiter Mechanismus mit eigenen Zusagen und Randformen im selben Träger; er erbt Datenquelle, Anker und Fail-Richtung dieses Wächters ([ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md), Folgepflicht).
- Ein **Anker für die Hauptsitzung** und eine Grenze für **Commands** in ihr — Klasse 3 (anderer Vorgang) und Grenze 12 (1) von [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md); eine Adresse gibt es nicht. Grund: Das Budget des Nutzers nennt „rollen und commands“; der Wächter deckt die Rollenläufe (Entscheidung des Nutzers vom 2026-10-22, „2. ja“), die Hauptsitzung hat keinen `SubagentStart`, und welche Quelle ihren Beginn liefert, ist nicht entschieden.
- Eine Zusage für die Zeit **innerhalb** eines laufenden Vordergrund-Befehls — Klasse 2 und eine **Grenze** des Vertrags (Punkt 12 (2)): Zwischen zwei Werkzeug-Aufrufen feuert kein Hook. Grund: Ein Mechanismus, der in einen Befehl greift, ist ein Hintergrundprozess mit Signalrecht, die verworfene Alternative.
- Ein **Abbruch** des Laufs oder eine **Rücknahme** seiner Arbeit — Klasse 2, mit der Begründung aus der Entscheidung selbst: Ein Wächter, der einen Lauf tötet, verliert den Bericht.
- Eine **Änderung der Zeitgrenze von 600 s je Gegenprobe** ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)) oder der Marken in §4 irgendeines Plans — Klasse 3 und `AGENTS.md` §3.6: Die Stunde ist eine zusätzliche Grenze.
- Die **Erfassungsschicht** (`.claude/hooks/span-emit.sh`, `harness/mk/erfassung.mk`, `harness/erfassung-feldliste.md`) — Klasse 2: Der Wächter ist ein Leser; die Feldliste wird vom Werkzeug kanonisch neu geschrieben. Braucht der Wächter ein Feld, das es nicht gibt, ist das ein Befund an den Architect.
- **Der Stop-Hook, die Gate-Ketten und ihre Nachweise** ([`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md)) — Schicht-Abgrenzung (Klasse 4) und `AGENTS.md` §3.6: Der Wächter urteilt über Zeit, nicht über Gates.
- **Vertragstext**, auch die Zeile der Gegenprobe in der Tabelle von [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16 — Der Implementer ändert keinen Vertrag; die Zeile trägt der Architect bei seiner Prüfung vor der Closure (DoD).
- **Produkt-Code und seine Tests** (`cmd/`, `internal/`, `test/`) und die **Nutzer-Doku** (`docs/user/`, `README.md`) — Schicht-Abgrenzung (Klasse 4): Kein Verhalten des gebauten Binary ändert sich.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die Closure-Pflichten darunter zählen nicht mit.

Liefer-Punkte (zwei):

- [ ] **Liefer-Punkt 1 — der Wächter, Stufe 1:** `.claude/hooks/pretooluse-zeitgrenze.sh` mit dem Leser `tools/harness/hook-felder.awk` (neu), reines bash und awk (kein `jq`, kein `node`, kein Netz), eingetragen in `.claude/settings.json` unter `PreToolUse` in einem eigenen Eintrag (`matcher` `""`, `timeout` 5; der Eintrag des Wächters nach `AGENTS.md` §3.1 bleibt unberührt) und mit den Zeilen in `harness/README.md` als **Werkzeug, kein Gate**. Er erfüllt die Punkte von [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md), die die Zuordnung des Architect der ersten Stufe gibt (Vorschlag des Planners: §6, Randform 7): wer begrenzt ist (Punkt 2), Strom, Anker und Zeit (Punkt 3 ohne die Bereichsformen von (e)), die Sperrmenge (Punkt 5), Ausgabe und Exit (Punkt 6 ohne die Zeile auf stderr), die allgemeine Freigabe bei Unsicherheit (Punkt 7 ohne Gründe und Auskunft), Umgebung (Punkt 8 (1) und (4)), kein Weg über die Grenze (Punkt 9). Jede Randform dieser Punkte ist in `SPEC-070` entschieden, bevor die erste Zeile Code entsteht.
- [ ] **Liefer-Punkt 2 — die Gegenprobe, Stufe 1:** ein Ziel `make zeitgrenze-gegenprobe` mit Fragment `harness/mk/zeitgrenze.mk` und Skript `tools/harness/zeitgrenze-gegenprobe.sh` (Namen: [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) Punkt 14), in `GEGENPROBEN_CONTAINER` ([`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 1 bis 3: bash und awk, ohne Docker), **nicht** in `GATE_CHECKS` und nicht in der Ausnahmeliste der schnellen Kette ([ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)). Sie fährt gegen **Attrappen-Ströme** (JSONL in `mktemp -d`, kein Agentenlauf, keine echte Erfassung) und eine Attrappe von `date` und trägt **je Zusage der ersten Stufe und je Element ihrer Aufzählungen einen Fall und eine Mutation, die nur ihn rot färbt** (`AGENTS.md` §3.10; die Namen der Fälle tragen die Nummer des Punktes). Die unveränderte Kopie ist grün, jeder Fall räumt auf, **kein Fall wartet** (die Zeit kommt aus dem Anker, nicht aus `sleep`), und der volle Lauf liegt **unter 600 s** ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)); die gemessene Zeit steht in §7. Der Code von `container-einrichten-gegenprobe.sh` darf übernommen werden. `harness/README.md` §Sensors führt die Zeile der Gegenprobe.

Gates und Closure:

- [ ] `make docs-check` und `make kopf-check` grün (der Slice ändert `harness/README.md` und diesen Plan; Vertrag und Entscheidung liegen vor der Basis).
- [ ] `make gates` grün (der Stop-Hook verlangt den vollen Lauf, weil `.claude/`, `harness/` und `tools/` ändern).
- [ ] Die Gegenprobe aus Liefer-Punkt 2 ist am Endstand **in einem Zug** grün, mit ihrer Zeit in §7 ([`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16); während der Arbeit fährt der Implementer nur den engsten Fall mit `GP_MUT` ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md), *Übergang*).
- [ ] **Der Wächter ist einmal im Agentenlauf gesehen** ([`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) Punkt 13, Handlung, kein Sensor): Der Implementer legt in seinem eigenen Lauf eine Marke an und liest sie, sichert seinen Strom mit `cp` und nennt die Summe, setzt `ts` der letzten Startzeile um zwei Stunden zurück, ruft `Bash` mit `echo make` auf (gesperrt, die Meldung im Ergebnis), dann `Write` auf eine Datei im Temp-Ordner (gesperrt), dann `Bash` mit `git status` (durchgelassen), und stellt den Strom aus der Sicherung wieder her (gleiche Summe). Der **Beleg ist eine eigene Datei unter `docs/reviews/`** (Namensform und Inhalt nach dem Ort des Belegs in [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20: Kopf mit Datum, Stand, Plattform und Fassung von Claude Code, Befehle und Ausgabe **ungekürzt**, die Datei wird einmal abgelegt und nicht geändert); §7 nennt Pfad, Summe nach `sha256sum`, die Zeile mit dem Exit-Code und den Kopf, nicht die Ausgabe. Kein Attrappen-Strom und kein grüner Lauf der Gegenprobe ersetzt ihn (`BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`). Zeigt der Beleg, dass die Eingabe von `PreToolUse` **kein** `agent_id` trägt oder dass die Form `decision: block` nicht greift, hält der Implementer an und gibt es dem Architect zurück (Randform-Rückgabe); die Wahl des Ankers und die Form des Blocks fallen dann.
- [ ] **Die Zeile der Gegenprobe in der Tabelle von [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16** (`zeitgrenze-gegenprobe`, Zustand `aufgenommen`, Fall `p14-aufnahme` im Vertrag) steht: Sie ist Vertragstext, den die Rolle des Implementers nicht ändert; sie trägt der **Architect** bei seiner Prüfung vor der Closure (Empfehlung des Architect), und die Closure prüft, dass sie steht.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: ein **weiterer Beleg** in `BEO-REPO/rollenlauf-ohne-zeitgrenze` (der Eintrag besteht seit der Closure von `slice-harness-devcontainer-einrichten`; eine Datei in `evidence/`, kein neues Verzeichnis); dazu die gemessene Zeit der neuen Gegenprobe als Beleg in `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (sie hängt nicht an `GATE_CHECKS`, die Kette wächst nicht). Keine weitere Beobachtung ist ebenfalls eine Antwort und steht in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen; die Folge-Slice-Paarung prüft Existenz **und** Annahme der Sendungen an `slice-harness-rollenlauf-zeitgrenze-regel`, `slice-harness-rollenlauf-zeitgrenze-haertung` und `slice-harness-rollenlauf-wiederholungsgrenze`: `grep` nach der Kennung dieses Slice findet sie in deren §1 und DoD (`AGENTS.md` §3.13).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md`, `docs/plan/adr/0058-…md` | liegen vor (Architect); die **Zuordnung der Punkte zu den zwei Stufen** ergänzt der Architect vor dem Code (Randform 7) | Zusagen und Randformen entschieden (`AGENTS.md` §3.12) |
| `.claude/hooks/pretooluse-zeitgrenze.sh` | neu | Liefer-Punkt 1: der Wächter; Form und Aufbau nach `.claude/hooks/pretooluse-command-guard.sh` |
| `tools/harness/hook-felder.awk` | neu | Liefer-Punkt 1: der Leser der Felder der Eingabe (`session_id`, `agent_id`, `tool_name`); den Befehl liest der Hook mit `tools/harness/extract-command.awk` |
| `.claude/settings.json` | update | Liefer-Punkt 1: eigener Eintrag unter `PreToolUse` (`matcher` `""`, `timeout` 5) |
| `harness/README.md` §Werkzeuge und §Sensors | update | Zeile des Wächters (Werkzeug, kein Gate) und Zeile der Gegenprobe; der Gate-Index steht **einmal** (`AGENTS.md` §4) |
| `harness/mk/zeitgrenze.mk` | neu | Liefer-Punkt 2: das Ziel, Eintrag in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS` |
| `tools/harness/zeitgrenze-gegenprobe.sh` | neu | Liefer-Punkt 2: Fälle und Mutationen gegen Attrappen-Ströme und eine Attrappe von `date`, Aufräumen je Fall |
| `docs/reviews/<Datum>-nachweis-slice-harness-rollenlauf-zeitgrenze.md` | neu (Implementer, Lauf-Beleg) | der Nachweis im Agentenlauf; von der Größenzählung ausgenommen ([ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md)) |
| `harness/vertraege/SPEC-060-gegenproben.md` Punkt 16 | update **durch den Architect** | die Zeile der Gegenprobe mit `aufgenommen` |
| `docs/plan/planning/observations/BEO-REPO/rollenlauf-ohne-zeitgrenze/`, `…/gate-laufzeit-waechst-mit-gegenprobe/` | update bei der Closure | Beleg und gemessene Zeit (kein Liefer-Punkt, DoD) |

- **Reihenfolge:** (1) der Architect ordnet die Punkte von `SPEC-070` den zwei Stufen zu (**vor** dem ersten Code-Commit und vor dem Übergang nach `in-progress`); (2) der Wächter mit seinem Eintrag in `.claude/settings.json` — der **Nachweis im Agentenlauf** folgt unmittelbar, weil er die Annahme belegt, auf der der Vertrag ruht (Eingabe mit `agent_id`, Form des Blocks); (3) die Gegenprobe mit ihren Fällen und Mutationen, je Zusage eine, vom Implementer selbst rot gesehen; (4) Closure mit Register und der Zeile in `SPEC-060`. Jeder Commit nennt die Kennung dieses Slice.
- **Was der Implementer nicht entscheidet:** Trifft er eine Randform, die `SPEC-070` nicht nennt, hält er an und gibt sie dem Architect zurück (`.claude/commands/implement-slice.md`, Randform-Rückgabe; `AGENTS.md` §3.12). Den Vertragstext ändert er nicht, auch nicht die Zeile in `SPEC-060`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): WIP-Limit frei; **Vertrag und ADR liegen vor** — `SPEC-070` vom Architect, [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) vom Nutzer angenommen (die Statuszeile `Accepted` mit dem Merge) — **und der Architect hat die Punkte der zwei Stufen zugeordnet**, jede Randform der ersten Stufe mit einer Antwort (`AGENTS.md` §3.12). Reihenfolge: **hinter** `slice-harness-rollenlauf-zeitgrenze-regel` (**Empfehlung des Planners**, kein Start-Trigger: Die Regel an den Leser gilt ab ihrer Lieferung und wartet nicht auf den Wächter, und der Wächter nicht auf sie); davor liegt `slice-harness-groessen-metrik`. Der Übergang `open` → `next` setzt `Verantwortlich:`.

**Größe.** Metrik nach [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (der Befehl steht in `AGENTS.md` §3.14 nach `slice-harness-groessen-metrik`; wiederholt dieser Plan ihn, gilt der von dort):

```
git diff --numstat --find-renames=10% <Basis>..HEAD -- AGENTS.md .claude harness tools \
  ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' \
  | awk '{a+=$1; b+=$2} END {print a+b+0}'
```

Basis ist der Commit, der diese Datei nach `in-progress/` bewegt (sein Hash steht beim Übergang hier). Unter `harness/vertraege` ist nichts erwartet (Vertrag vor der Basis; ein Nachtrag zählt mit, deshalb steht `harness` im Befehl); der Beleg des Nachweises liegt unter `docs/reviews` und zählt nicht.

**Marke 900, Meldeschwelle 700.** Erreicht die Summe 700, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`); über 900 geht der Slice zurück zur Zerlegung.

**Schätzung des Planners (2026-10-22), in der Metrik oben, mit dem Faktor 3 bis 4,6 für Gegenprobe-Code gegenüber Produktcode** (Messwerte in `slice-harness-schaetzung-planner-faktor` §1: 2,9, 3,2 und 4,6): **440 bis 759, Mitte 600** (67 % der Marke; unter 70 %). Die Gegenprobe ist mit 72 der 125 Fälle der Schätzung des Architect (550 bis 950) angesetzt:

| Teil | Schätzung |
|---|---|
| Wächter (`.claude/hooks/pretooluse-zeitgrenze.sh`) und Leser (`tools/harness/hook-felder.awk`) | 90 bis 150 |
| `.claude/settings.json` | 8 bis 12 |
| Fragment `harness/mk/zeitgrenze.mk` | 15 bis 30 |
| `harness/README.md` (zwei Zeilen) | 10 bis 20 |
| Gegenprobe (`tools/harness/zeitgrenze-gegenprobe.sh`), Faktor 3,5 bis 3,6 auf den Wächter | 317 bis 547 |

Der Beleg des Nachweises im Agentenlauf und die Zeile in `SPEC-060` zählen nicht (Pfad ausgenommen bzw. Vertragstext des Architect, höchstens 5 Zeilen).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt **über 900**, oder der Architect entscheidet eine Randform so, dass ein dritter Liefer-Punkt entsteht (etwa ein Freigabe-Mechanismus, der die Grenze hebt). Der Schnitt nach Lieferwert ist **vollzogen**, nicht vorgemerkt (§1); ein weiterer Schnitt wäre eine neue Zuordnung der Punkte durch den Architect.
- `in-progress` → `open` (blockiert): Der Nachweis im Agentenlauf zeigt, dass die Eingabe von `PreToolUse` kein `agent_id` trägt oder die Form des Blocks nicht greift (Befund an den Architect: Anker und Form fallen); oder die Zuordnung der Punkte zu den Stufen steht nicht im Vertrag; oder der Nutzer nimmt die Entscheidung vom 2026-10-21 zurück.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Beobachtbar: (1) `make gates` und `make docs-check` grün, und die Gegenprobe aus Liefer-Punkt 2 in einem Zug grün mit ihrer Zeit unter 600 s in §7; (2) der Beleg des Wächters im Agentenlauf als Datei unter `docs/reviews/`, mit Pfad, Summe, Exit-Zeile und Kopf in §7 (Block-Meldung bei `echo make`, `Write` gesperrt, `git status` durchgelassen, Strom wiederhergestellt); dazu die Zeile `zeitgrenze-gegenprobe` mit `aufgenommen` in `SPEC-060` Punkt 16, die Closure-Notiz mit Lerneintrag und der Beleg im Register aus §2.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

### Randformen — entschieden in `SPEC-070`, eine offen für den Architect

Die sechs Fragen des früheren Plans (Sperrmenge und Abwicklung, Fail-Richtung, Befehl und Anker, die Lücken, Aufrufarten und Umgebung, Weg über die Grenze) und die Reihenfolge der Gründe, die Stelle der Aufrufe, die Aufruf-Voraussetzung, die pflichtigen Läufe und die Plattform sind in [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) entschieden; die Tabelle *Randform → Punkt* steht dort, und dieser Plan wiederholt sie nicht (`AGENTS.md` §3.12). **Offen ist eine**, vor dem ersten Code-Commit zu entscheiden:

7. **Die Zuordnung der Punkte zu zwei Lieferungen (Stufung).** Vorschlag des Planners, Empfehlung und keine Entscheidung; der Vertrag trägt sie erst, wenn der Architect sie festhält. Jede Zusage steht in genau einer Stufe, mit ihren Fällen und Mutationen; zwischen den Stufen gibt es keine Zusage ohne Prüfung, und Stufe 1 ist allein lieferbar (der Hook wirkt, still und sicher):

| Punkt von `SPEC-070` | Stufe 1 (dieser Slice) | Stufe 2 (`slice-harness-rollenlauf-zeitgrenze-haertung`) |
|---|---|---|
| 1 Eintrag, Leser, Befehl | alle Fälle | — |
| 2 wer begrenzt ist | alle Fälle | — |
| 3 Anker, Strom, Zeit | (a) bis (d), (f), (g), Schaltjahr, Monatsende, Zeitzonen; von (e) die Gestalt `YYYY-MM-DDTHH:MM:SSZ` (`p3-form-versatz`, `-bruchteil`, `-klein-z`, `-leerzeichen`, `-leer`, `-fehlt`, `-text-davor`) | von (e) die Wertebereiche (`p3-form-monat-13`, `-tag-32`, `-februar-30`, `-stunde-24`, `-minute-60`, `-sekunde-60`) |
| 4 Reihenfolge der Gründe | — | alle Fälle |
| 5 Sperrmenge | alle Fälle (Werkzeuge und Wörter, Wortgrenze, Durchgelassenes) | — |
| 6 Ausgabe und Exit | `p6-block-text`, `p6-frei-leer`, Exit 0 in jedem Zweig | `p6-stderr-eine-zeile` |
| 7 Freigabe bei Unsicherheit | die allgemeine Zusage (Freigabe, still, Exit 0), ein Fall und eine Mutation | die Aufzählung der Gründe, je Grund die Zeile auf stderr |
| 8 Umgebung und Aufrufe | (1) Wurzel, (4) Worktree und Container | (2) leere Umgebung, (3) Zeitgrenze des Hooks, (5) gleichzeitige Aufrufe |
| 9 kein Weg über die Grenze | alle Fälle | — |
| 10 bis 12 pflichtige Läufe, Plattform, Grenzen | Prosa im Vertrag, kein Fall | — |
| 13 Nachweis im Agentenlauf | Handlung dieses Slice | nur bei geänderter Eingabe oder Form |
| 14 Gegenprobe, Zeit, Aufnahme | Gerüst, Aufnahme, Zeile in `SPEC-060` | Erweiterung, Zeit weiter unter 600 s |

**Risiken:**

- **Der Wächter nimmt der Rolle den Bericht** (zu enge Liste der durchgelassenen Aufrufe): Der Lauf wird gesperrt, kann aber nicht mehr committen oder berichten — genau der Schaden, den die Regel verhindern soll. — **Ausgang (vorgesehen):** entfallen, wenn die Fälle `p5-frei-*` und der Nachweis im Agentenlauf (`git status` durchgelassen) ihn zeigen; sonst eingetreten → Rückführung (§4).
- **Die Annahme des Vertrags fällt** (die Eingabe von `PreToolUse` trägt kein `agent_id`, die Form `decision: block` greift nicht): Der Hook wirkte in keinem Lauf, und die Gegenprobe zeigte es nicht. — **Ausgang (vorgesehen):** eingetreten → Randform-Rückgabe an den Architect, Slice nach `open` (§4); entfallen, wenn der Nachweis im Agentenlauf sie bestätigt.
- **Der Hook schweigt in dieser Stufe, wo der Strom fehlt** (kein Grund auf stderr; die Auskunft liefert Stufe 2): Ein Lauf läuft unbegrenzt, und niemand sieht, dass der Wächter nichts tut. — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/rollenlauf-ohne-zeitgrenze`; Gegenmittel: Stufe 2 steht direkt dahinter, in derselben Welle.
- **Latenz vor jedem Werkzeug-Aufruf** (der Hook liest eine Datei, die mit dem Lauf wächst): gemessen sind 21 ms für einen Strom von 5 MB (`SPEC-070` Punkt 8 (3)). — **Ausgang (vorgesehen):** entfallen, mit der Zahl des Vertrags.
- **Die Zeit der Gegenprobe:** Fälle mit Zeitsprüngen verführen zu `sleep` (`BEO-REPO/feste-wartezeit-statt-beobachtung-im-test`), und [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) setzt 600 s. — **Ausgang (vorgesehen):** entfallen, wenn kein Fall wartet (die Zeit kommt aus dem Anker) und der volle Lauf unter 600 s liegt (§7).
- **Zwei Hooks mit entgegengesetzter Fail-Richtung im selben Verzeichnis** sind eine Lesart-Falle für den, der einen dritten schreibt (`BEO-REPO/harness-lesart-ohne-entscheidungsort`). — **Ausgang (vorgesehen):** entfallen, wenn die Begründung in der ADR und im Kopf beider Dateien steht, nicht nur in einer.
- **Die Zuordnung der Stufen steht nicht im Vertrag** (Randform 7): Der Slice lieferte Zusagen ohne Fall. — **Ausgang (vorgesehen):** eingetreten → Rückführung nach `open` (§4); die Alternative ist eine Entscheidung des Nutzers über Marke und Meldeschwelle (§1).
- **Der Slice passt nicht in einen Rollenlauf von einer Stunde** (Hook, Gegenprobe mit 72 Fällen, Nachweis): Er läuft in mehreren Läufen, und jeder hält nach der Regel an und berichtet. — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/rollenlauf-ohne-zeitgrenze` (Beleg: die Zahl der Läufe); die Reihenfolge der Punkte lässt jeden Lauf an einer Naht enden.
- **Ein Wort in einer Commit-Message (`make`, `docker`, `sleep`) sperrt den Aufruf nach der Grenze:** akzeptiertes Negativ im Vertrag (Stolperdraht, keine Sandbox); die Meldung sagt es. — **Ausgang (vorgesehen):** entfallen, weil im Vertrag entschieden (Punkt 5, Grenze 12 (5)).
- **Die Zeile in `SPEC-060` fehlt zur Closure:** Sie trägt der Architect, nicht die Rolle des Implementers. — **Ausgang (vorgesehen):** entfallen, wenn die Closure sie findet; sonst eingetreten → die Closure hält an und gibt es dem Architect zurück.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren) · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln (das Feld `liegt in` steht **nur**, wenn mit diesem Slice
wirklich etwas verkörpert wurde).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Hier stehen:

- **Was hat funktioniert / was ging anders als geplant:** —
- **Belege des Implementers:** die Mutationstabelle je Zusage der ersten Stufe von `SPEC-070` (Zusage · Mutation · roter Fall), die Zeit des vollen Gegenprobe-Laufs, und der **Nachweis im Agentenlauf** — Pfad der Datei unter `docs/reviews/`, ihre Summe nach `sha256sum`, die Zeile mit dem Exit-Code und der Kopf (Datum, Stand, Plattform, Fassung von Claude Code), nicht die Ausgabe.
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** ein weiterer Beleg in `BEO-REPO/rollenlauf-ohne-zeitgrenze` (`evidence/slice-harness-rollenlauf-zeitgrenze.md`); `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` um die gemessene Zeit ergänzt.
- **Folge-Slices:** `slice-harness-rollenlauf-zeitgrenze-regel`, `slice-harness-rollenlauf-zeitgrenze-haertung`, `slice-harness-rollenlauf-wiederholungsgrenze` (Dateien in `open/` oder `done/`).
- **Risiken aus §6:** jedes mit genau einem Ausgang — siehe §6.
- **Drei Paarungen:** Anker · Folge-Slice · Register, Ergebnis.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — die **zwei vorgelagerten Schritte**
(in jedem Slice-Plan) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine Sub-Area, `REPO` (Default des gesamten Repos, `harness/conventions.md` §Modus-Deklaration pro Sub-Area, Greenfield). Die drei Achsen sind erfüllt (eigene Artefakte unter `.claude/` und `harness/`, eigener Maßstab im Review, eigene Konventionen in `AGENTS.md`); eine feinere Ausdifferenzierung deklariert dieses Repo nicht, und ein Slice ist kein Ort, eine zu erfinden.

**Vorgelagert — offene Beobachtungen sichten** (gemergter Stand `main`, 2026-10-22):

- `BEO-REPO/rollenlauf-ohne-zeitgrenze` (1×, unter der Schwelle; `Stand` offen) — der Auslöser: Ein Rollenlauf ohne Zeitgrenze. Der Eintrag besteht seit der Closure von `slice-harness-devcontainer-einrichten`; dieser Slice ist die Adresse des Wächters, und sein Beleg ist eine weitere Datei in `evidence/`.
- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (14 Belege, `verkörpert`): Dieser Slice fügt eine Gegenprobe hinzu. Sie hängt an `GEGENPROBEN_CONTAINER`, nicht an `GATE_CHECKS`, und fährt ohne Docker — die Kette wächst nicht, die Fahrt von `make gegenproben` um sie. Beleg bei der Closure (§2).
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (`verkörpert` in `AGENTS.md` §3.12): Der Wächter läuft, wo Claude Code läuft; die Plattform ist in `SPEC-070` entschieden (Linux geprüft, macOS und Windows Grenze).
- `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (10×, `geplant` mit `slice-harness-randform-umgebung-der-gegenprobe`) und `BEO-REPO/span-zuordnung-folgt-dem-hauptbaum-nicht-dem-klon` (offen): beide treffen die Datenquelle; der Vertrag nennt die Umgebung (Punkt 8) und den Nachweis im Agentenlauf (Punkt 13). Erreicht der zweite Eintrag mit diesem Slice sein drittes Auftreten, braucht er einen eigenen Folge-Slice (bei der Closure zu prüfen).
- `BEO-REPO/meldeschwelle-ohne-artefakt`, `BEO-REPO/groessenprobe-ohne-messbaren-stand` (`verkörpert`) und `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`: §4 nennt Metrik, Marke, Meldeschwelle und rechnet die Gegenprobe mit dem Faktor 3 bis 4,6.
- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (3×, offen): Der Schnitt dieses Plans vermeidet das Muster (beide Stufen unter 70 % der Marke, kein zweiter Schnitt vorab benannt); Voraussetzung ist die Zuordnung des Architect (Randform 7).
- `BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet` und `BEO-REPO/folge-slice-adresse-nimmt-nicht-an`: getragen — der Kopf trennt Entscheidung des Nutzers (mit Datum und Wortlaut) und Empfehlung des Planners; die drei Adressen existieren und nehmen die Sendungen in §1 und DoD an.

Keiner der gesichteten Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Zählung nach `AGENTS.md` §3.13 und §3.14** (Planner, 2026-10-22, im Commit, der diesen Plan schneidet und die Sendungen einträgt): **Liefer-Punkte 2** (Wächter mit Leser, Eintrag und Index-Zeilen | Gegenprobe). **Schichten 1** — der Harness (Hook, Leser, Fragment, Skript, `harness/README.md`, `.claude/settings.json`; Teilung nach `slice-harness-d-check-v0-85` §8: Gate-Konfiguration, Gate-Werkzeuge und Agenten-Anweisungen sind eine Schicht); die Entscheidung und der Vertrag liegen **vor** der Basis und gehören nicht in diesen Diff, ein Nachtrag macht zwei; die Planung ist keine Schicht. **Schätzung 440 bis 759, Mitte 600** bei **Marke 900** und Meldeschwelle 700 (§4), mit dem Faktor 3 bis 4,6 für Gegenprobe-Code — unter drei Liefer-Punkten, unter zwei Schichten und unter 70 % der Marke. **Herkunft der Zählung:** Der ungeteilte Plan (angelegt am 2026-10-21) hatte 3 Liefer-Punkte (Wächter | Gegenprobe | Regel), 1 Schicht, die Schätzung des Planners 490 bis 900 mit dem Faktor 2; die Schätzung des Architect nach `SPEC-070` (800 bis 1400) und der Faktor 3 bis 4,6 hoben sie über die Marke; der Schnitt gibt die Regel (Liefer-Punkt 3) an `slice-harness-rollenlauf-zeitgrenze-regel` und die zweite Stufe an `slice-harness-rollenlauf-zeitgrenze-haertung` ab. Die Sendungen an die drei Nehmer legen in diesem Slice keinen Liefer-Punkt und keine Schicht an.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
