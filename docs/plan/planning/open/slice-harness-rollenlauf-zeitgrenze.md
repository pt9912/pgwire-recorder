# Slice slice-harness-rollenlauf-zeitgrenze: Ein Rollenlauf hält bei einer Stunde an und berichtet, erzwungen von einem deterministischen Wächter im Werkzeug-Aufruf

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-harness-regeln](../welle-harness-regeln.md) — **Empfehlung des Planners vom 2026-10-21** für Welle und Stelle (hinter `slice-harness-groessen-metrik`, §4 *Start*), nicht Entscheidung des Nutzers: Der Nutzer hat die **Sache** entschieden (unten), über Welle und Reihenfolge ist nichts entschieden. Angelegt in `open/`, weil der Start zwei Vorbedingungen hat, die noch fehlen (Vertrag und ADR des Architect, §4).

**Bezug:** kein Anspruch des Produkts (`LH-*`) — der Slice erzwingt eine Regel über Agentenläufe und ändert kein Verhalten des Werkzeugs. **Entscheidung des Nutzers vom 2026-10-21** (wörtlich, vom Koordinator übermittelt): „Wir haben rollen (siehe .claude/agents) und commands -- max. 1 Stunde", und „ok" auf den Entwurf des Mechanismus (Wächter im Werkzeug-Aufruf, reines bash+awk, **kein** Hintergrundprozess, Datenquelle ist der vorhandene Span-Strom; dazu die Wiederholungs-Grenze für volle Gegenproben, §1 *Ausdrücklich NICHT*). Aktive ADRs, gelesen: [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Zeitbudget 600 s je Gegenprobe — gilt auch für die Gegenprobe dieses Slice), [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (die Gegenprobe hängt an `GEGENPROBEN_CONTAINER`, nicht an `GATE_CHECKS`), [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (die Hooks dieses Repos und ihr Zusammenspiel), [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (Metrik der Größenmarke in §4). Die **ADR dieses Slice** (neues Element der Durchsetzungsschicht; Nummer vergibt der Architect, nächste freie ist heute die 58) schreibt der Architect vor dem Code, §4 und §6.

**Berührte Spec-Stellen:** `SPEC-070` (neu; Nummer vergibt der Architect, nächste freie; Vertrag in `harness/vertraege/`, vor dem ersten Code-Commit, `AGENTS.md` §3.12) · [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) (gelesen: Name und Liste der Gegenprobe, Punkte 1 bis 3; Aufnahme nach Punkt 16) · [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) (gelesen: der Stop-Hook und die Ketten — dieser Slice ändert beide nicht).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Rollenlauf, der eine Stunde läuft, **hält an und berichtet**, und das hängt nicht an seiner Disziplin: Ein Wächter im Werkzeug-Aufruf (`PreToolUse`, reines bash+awk wie `.claude/hooks/pretooluse-command-guard.sh`, kein `jq`, kein `node`) liest den Span-Strom des laufenden Agenten, rechnet die Zeit seit dem Beginn des Laufs und blockt den Aufruf mit einer Meldung, die die Handlung nennt (anhalten, berichten, Arbeit nicht zurücknehmen), sobald die Grenze erreicht ist. Die Regel steht daneben in `AGENTS.md` und in den sechs Rollen-Typen unter `.claude/agents/`, damit sie auch gelesen wird, wo der Wächter nicht greift.

**Warum ein Wächter im Aufruf und kein Beobachter:** Ein Hintergrundprozess kann nur zusehen — er hat keinen Weg, eine laufende Rolle anzuhalten. Der Wächter läuft im Werkzeug-Aufruf der Rolle selbst und kann sie anhalten. Das ist die **Entscheidung des Nutzers vom 2026-10-21** („ok" auf den Entwurf), und die ADR des Architect trägt sie mit den verworfenen Alternativen (Hintergrund-Beobachter, Regel ohne Mechanismus).

**Woran der Anlass gemessen ist:** Der Implementer-Lauf von `slice-harness-devcontainer-einrichten` lief **4,34 h** (Span-Strom: `SubagentStart` 2026-10-21T12:08:37Z, letzte Zeile 16:29:11Z, 251 Zeilen, Befehlszeit 7448 s = 48 % des Laufs), ohne dass eine Marke feuerte. Vorher gab es **kein** Zeitbudget für einen Lauf: [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) budgetiert die Gegenprobe (600 s je vollem Lauf), `.claude/commands/implement-slice.md` trägt keine Zeitregel, `AGENTS.md` §6 nur den auf zwei Gegenproben eingeschränkten Zeiger, und §4 eines Plans misst in **Zeilen** ([ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md)) — in Zeit nichts.

**Was der Slice nicht entscheidet, sondern vorlegt:** Die Randformen in §6 sind **Fragen an den Architect**, keine Antworten des Planners (`AGENTS.md` §3.12: entschieden vor dem ersten Code-Commit, in `harness/vertraege/`). Dazu gehört ausdrücklich die **Fail-Richtung**: Der Planner empfiehlt fail-open (Begründung in §6, Randform 2); der bestehende Wächter nach §3.1 failt umgekehrt, und die Asymmetrie braucht ihre Begründung im Vertrag, nicht im Plan.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die **Wiederholungs-Grenze für volle Gegenprobe-Läufe** (nach dem ersten vollen Lauf je Rollenlauf wird ein weiterer ohne `GP_NUR`/`GP_MUT` geblockt, aus `duration_ms` desselben Span-Strom gerechnet) — **vom Nutzer am 2026-10-21 freigegeben, aber nicht dieser Slice** (Klasse 3, anderer Vorgang: ein zweiter Mechanismus mit eigenen Zusagen, eigenen Randformen und eigenen Fällen im selben Träger). Anlass: 46 volle Läufe à ~175 s in einem Lauf, während der Fehler in zwei Fällen saß. Grund für den Schnitt: Mit ihr liegt die Schätzung über der Marke (§4, §8), und ihre Randformen sind andere (was „voll" heißt, woran `GP_NUR`/`GP_MUT` erkannt wird, was bei einer *anderen* Gegenprobe gilt). **Eine Adresse gibt es noch nicht** — der Planner legt den Slice an, sobald der Architect den Vertrag dieses Slice entschieden hat, weil die Wiederholungs-Grenze dessen Datenquelle, dessen Anker und dessen Fail-Richtung erbt; ihr eigener Plan hängt also am Vertrag. Bis dahin ist nichts zugewiesen (`AGENTS.md` §3.13: keine Kennung genannt, keine Sendung).
- Ein **Anker für die Hauptsitzung** (ein Command, der nicht als Subagent läuft, hat keinen `SubagentStart` und damit keinen Beginn im Strom) — Klasse 3 (anderer Vorgang: ein zweiter Anker, nicht derselbe Wächter mit einem Zusatz). Grund: Das Budget des Nutzers nennt „rollen (siehe .claude/agents) und commands"; die sechs Rollen deckt dieser Slice **exakt**, die Hauptsitzung braucht eine andere Quelle für ihren Beginn, und welche das ist, ist nicht entschieden. Der Vertrag **benennt diese Lücke**, statt sie zu verschweigen (`AGENTS.md` §3.11, §6 Randform 4a); ein Satz, der mehr zusagt, als der Wächter hält, ist der Befund, gegen den diese Regel steht.
- Eine Zusage für die Zeit **innerhalb** eines laufenden Vordergrund-Befehls — Klasse 2 (Bestand bleibt bewusst stehen) und eine **Grenze**, die der Vertrag nennt: Zwischen zwei Werkzeug-Aufrufen feuert kein Hook; läuft ein einzelner Befehl 40 Minuten, bemerkt der Wächter die Überschreitung erst danach. Grund: Ein Mechanismus, der mitten in einen Befehl greift, ist ein Hintergrundprozess mit Signalrecht — genau die verworfene Alternative, und der Nutzer hat gegen sie entschieden.
- Ein **Abbruch** des Laufs oder eine **Rücknahme** seiner Arbeit — Klasse 2 mit Begründung aus der Entscheidung selbst: „bei der Grenze hält er an und **berichtet** — kein Abbruch, keine Rücknahme von Arbeit". Grund: Ein Wächter, der einen Lauf tötet, verliert den Bericht, der die einzige verwertbare Ausgabe eines angehaltenen Laufs ist.
- Eine **Änderung der Zeitgrenze von 600 s je Gegenprobe** ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)) oder der Marken und Meldeschwellen in §4 irgendeines Plans — Klasse 3 (anderer Vorgang) und `AGENTS.md` §3.6: Die Stunde ist eine **zusätzliche** Grenze über dem Lauf, keine Senkung einer bestehenden.
- Die **Erfassungsschicht selbst** (`.claude/hooks/span-emit.sh`, `harness/mk/erfassung.mk`, `harness/erfassung-feldliste.md`) — Klasse 2: Der Wächter ist ein **Leser**; er ändert kein Feld und keine Zeile des Stroms. Grund: `harness/erfassung-feldliste.md` ist werkzeug-erzeugt und wird von einem erneuten Lauf des Werkzeugs kanonisch neu geschrieben; eine Änderung von Hand ginge verloren. Braucht der Wächter ein Feld, das es nicht gibt, ist das ein Befund an den Architect (§6, Randform 3), kein Eingriff dieses Slice.
- **Der Stop-Hook, die Gate-Ketten und ihre Nachweise** (`.claude/hooks/stop-require-gates.sh`, `tools/harness/gates-schnell.sh`, [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md)) — Schicht-Abgrenzung (Klasse 4) und `AGENTS.md` §3.6. Grund: Der Wächter urteilt über Zeit, nicht über Gates; ein Eingriff dort änderte die Freigabe-Logik des Stop-Hooks mit. Beim Review ist sofort prüfbar, dass der Diff diese Dateien nicht anfasst.
- **Produkt-Code und seine Tests** (`cmd/`, `internal/`, `test/`) und die **Nutzer-Doku** (`docs/user/`, `README.md`) — Schicht-Abgrenzung (Klasse 4). Grund: Kein Verhalten des gebauten Binaries ändert sich; ein Satz im Handbuch beschriebe etwas, das das Binary nicht tut (`AGENTS.md` §3.11).
- Die **Dateien des laufenden Slice** `slice-harness-devcontainer-einrichten` (`tools/`, `harness/`, sein Plan) in der Fassung, die dessen Implementer gerade schreibt — Klasse 3 (anderer Vorgang) und eine Vorgabe des Koordinators vom 2026-10-21. Grund: Zwei Läufe auf denselben Dateien erzeugen einen Konflikt statt einer Regel; der Start dieses Slice liegt ohnehin hinter dessen Abwicklung (§4).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die Closure-Pflichten darunter zählen nicht mit.

Liefer-Punkte (drei):

- [ ] **Liefer-Punkt 1 — der Wächter:** ein Hook unter `.claude/hooks/` (Name: Architect), reines bash+awk (kein `jq`, kein `node`, kein Netz), eingetragen in `.claude/settings.json` unter `PreToolUse` und in `harness/README.md` als **Werkzeug, kein Gate** (er prüft keinen Stand des Repos). Er erfüllt `SPEC-070`: er findet den Span-Strom des laufenden Agenten, rechnet die Zeit seit dessen Beginn, lässt unter der Grenze jeden Aufruf ohne Ausgabe durch (wie der bestehende Wächter: keine Ausgabe heißt, die normale Berechtigungs-Prüfung übernimmt), und blockt ab der Grenze mit einer Meldung, die die Handlung nennt (anhalten, berichten, nichts zurücknehmen) und die durchgelassenen Aufrufe benennt. Jede Randform aus §6 ist in `SPEC-070` entschieden, bevor die erste Zeile Code entsteht; der Befehl wird aus derselben Quelle gelesen wie beim bestehenden Wächter, nicht aus einem Span-Feld (§6, Randform 3).
- [ ] **Liefer-Punkt 2 — die Gegenprobe:** ein Ziel `make <name>-gegenprobe` (Name: Architect) mit Fragment unter `harness/mk/`, in `GEGENPROBEN_CONTAINER` ([`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 1 bis 3: bash und awk, ohne Docker), **nicht** in `GATE_CHECKS` und nicht in der Ausnahmeliste der schnellen Kette ([ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)). Sie fährt gegen **Attrappen-Ströme** (JSONL-Dateien in `mktemp -d`, kein Agentenlauf, keine echte Erfassung) und trägt **je Zusage von `SPEC-070` und je Element ihrer Aufzählungen einen Fall und eine Mutation, die nur ihn rot färbt** (`AGENTS.md` §3.10): unter der Grenze frei · über der Grenze geblockt · genau an der Grenze (die Seite, die der Vertrag nennt) · je durchgelassener Abwicklungs-Aufruf einer · fehlende Erfassung (kein Verzeichnis, kein Strom, leerer Strom, Strom ohne Anker) in der Richtung, die der Vertrag entscheidet · Strom eines **anderen** Agenten derselben Sitzung · Zeitform und Zeitzone des Ankers. Die unveränderte Kopie ist grün, jeder Fall räumt auf, und der volle Lauf liegt **unter 600 s** ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)); die gemessene Zeit steht in §7. Die Gegenprobe ist **ab ihrem Anlegen aufgenommen** im Sinn von [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16, und `harness/README.md` §Sensors führt ihre Zeile.
- [ ] **Liefer-Punkt 3 — die Regel, in einer Fassung, die der Nutzer bestätigt hat:** `AGENTS.md` trägt die Grenze als eigene harte Regel (Nummer: nächste freie, nach §3.14) mit Falsch/Richtig-Paar, Herkunfts-Anker `seit slice-harness-rollenlauf-zeitgrenze` und Retirement-Check; die **sechs** Rollen-Typen unter `.claude/agents/` (`planner`, `architect`, `implementer`, `reviewer`, `verifier`, `validator`) nennen je in einem Satz, was bei der Grenze geschieht (anhalten, berichten, nichts zurücknehmen) und dass der Wächter sie hält. Der Wortlaut wird dem Nutzer **vor** dem Commit vorgelegt (§6 der Welle: „Der Wortlaut der Regeln steht beim Nutzer"); §7 nennt Datum und Fundstelle der Bestätigung. `grep` nach der Kennung dieses Slice findet die Regel in `AGENTS.md` und je einen Satz in den sechs Rollen-Dateien.

Gates und Closure:

- [ ] `make docs-check` und `make kopf-check` grün (der Slice ändert `AGENTS.md`, `harness/README.md`, ein Dokument unter `harness/vertraege/`, eine ADR und diesen Plan).
- [ ] `make gates` grün (der Stop-Hook verlangt den vollen Lauf, weil `harness/`, `tools/` und ein Makefile-Fragment ändern).
- [ ] Die Gegenprobe aus Liefer-Punkt 2 ist am Endstand **in einem Zug** grün, mit ihrer Zeit in §7 ([`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16); während der Arbeit fährt der Implementer nur den engsten Fall mit `GP_MUT` ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md), *Übergang*).
- [ ] Der **Wächter ist einmal im Lauf gesehen**: §7 trägt den Beleg, dass er unter der Grenze nichts tut und an der Grenze blockt — gegen einen Attrappen-Strom mit vorgerückter Zeit gefahren, Eingabe und Ausgabe des Hooks wörtlich. Ein grüner Gegenprobe-Lauf allein belegt das nicht (`AGENTS.md` §3.10; der Hook läuft in der Umgebung des Agenten, die Gegenprobe im Container, `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: **neu** `BEO-REPO/rollenlauf-ohne-zeitgrenze` (kein vorhandener Eintrag trägt den Anlass, §8) mit dem Beleg dieses Slice; dazu die gemessene Zeit der neuen Gegenprobe als Beleg in `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (sie hängt nicht an `GATE_CHECKS`, die Kette wächst nicht). Keine weitere Beobachtung ist ebenfalls eine Antwort und steht in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-070-…md` | neu, **durch den Architect** vor dem Code | Zusagen des Wächters und **jede** Randform aus §6 entschieden (`AGENTS.md` §3.12); dazu die zwei Lücken benannt (langer Vordergrund-Befehl, Hauptsitzung ohne Anker) |
| `docs/plan/adr/00NN-…md` + `docs/plan/adr/README.md` | neu, **durch den Architect** vor dem Code | neues Element der Durchsetzungsschicht mit echten Alternativen (Wächter im Aufruf / Hintergrund-Beobachter / Regel ohne Mechanismus) und der Begründung der Fail-Richtung; Index-Zeile in derselben Änderung |
| `.claude/hooks/<name>.sh` (+ `tools/harness/<name>.awk`, wenn der Architect den Leser trennt) | neu | Liefer-Punkt 1: der Wächter; Form und Aufbau nach dem Muster von `.claude/hooks/pretooluse-command-guard.sh` (Kopf mit Zusage und Grenze, Extraktor in awk, fail-Richtung an einer Stelle) |
| `.claude/settings.json` | update | Liefer-Punkt 1: Eintrag unter `PreToolUse`; welcher `matcher` und welche Zeitgrenze des Hooks, entscheidet der Vertrag (§6, Randform 1) |
| `tools/harness/<name>-gegenprobe.sh` | neu | Liefer-Punkt 2: Fälle und Mutationen gegen Attrappen-Ströme, Aufräumen je Fall |
| `harness/mk/<name>.mk` | neu | Liefer-Punkt 2: das Ziel, Eintrag in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS` |
| `harness/README.md` §Sensors und §Werkzeuge | update | Zeile der Gegenprobe und Zeile des Wächters (Werkzeug, kein Gate); der Gate-Index steht **einmal** (`AGENTS.md` §4) |
| `AGENTS.md` (neue harte Regel hinter §3.14) | update | Liefer-Punkt 3: die Regel mit Paar, Anker und Retirement-Check; Wortlaut dem Nutzer vorgelegt |
| `.claude/agents/planner.md`, `architect.md`, `implementer.md`, `reviewer.md`, `verifier.md`, `validator.md` | update | Liefer-Punkt 3: je ein Satz — Grenze, Anhalten, Bericht; sechs Dateien, ein Satz, damit keine Rolle die Regel nur aus dem Block erfährt |
| `docs/plan/planning/observations/BEO-REPO/rollenlauf-ohne-zeitgrenze/` | neu bei der Closure | Register-Eintrag mit Beleg (kein Liefer-Punkt, DoD) |
| `docs/plan/planning/observations/BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/` | update bei der Closure | Beleg mit der gemessenen Zeit der neuen Gegenprobe (kein Liefer-Punkt, DoD) |

- **Reihenfolge:** (1) der Architect schreibt ADR und `SPEC-070` und entscheidet jede Randform aus §6 — **vor** dem ersten Code-Commit und vor dem Übergang nach `in-progress`; (2) der Wächter mit seinem Eintrag in `.claude/settings.json`; (3) die Gegenprobe mit ihren Fällen und Mutationen, je Zusage eine, vom Implementer selbst rot gesehen; (4) der Beleg des Laufs gegen einen Attrappen-Strom in §7; (5) der Wortlaut der Regel dem Nutzer vorgelegt, dann `AGENTS.md` und die sechs Rollen-Dateien; (6) Closure mit Register. Jeder Commit nennt die Kennung dieses Slice.
- **Reihenfolge-Grund für (5) am Ende:** Die Regel beschreibt, was der Wächter tut; steht sie vor ihm, sagt sie etwas zu, das noch kein Fall prüft (`AGENTS.md` §3.11).
- **Was der Implementer nicht entscheidet:** Trifft er eine Randform, die `SPEC-070` nicht nennt, hält er an und gibt sie dem Architect zurück (`.claude/commands/implement-slice.md`, Randform-Rückgabe; `AGENTS.md` §3.12). Den Vertragstext ändert er nicht selbst.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): WIP-Limit frei; `slice-harness-groessen-metrik` liegt in `done/` (**Empfehlung des Planners** für die Stelle: beide schreiben an `AGENTS.md` §3.14 bzw. unmittelbar dahinter, und der Metrik-Slice schreibt diesen Text neu — ein gleichzeitiger Lauf erzeugt einen Konflikt; über die Stelle hat der Nutzer nichts entschieden); `slice-harness-devcontainer-einrichten` ist abgewickelt (dessen Lauf hält `tools/` und `harness/`); **Vertrag und ADR des Architect liegen vor und sind entschieden**, jede Randform aus §6 mit einer Antwort (`AGENTS.md` §3.12); der **Wortlaut** der Regel aus Liefer-Punkt 3 ist dem Nutzer vorgelegt oder wird es in Schritt (5) — er ist kein Start-Trigger. Der Übergang `open` → `next` setzt `Verantwortlich:`.

**Größe.** Metrik nach [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (`AGENTS.md` §3.14 trägt sie nach `slice-harness-groessen-metrik`; wiederholt dieser Plan den Befehl, gilt der von dort):

```
git diff --numstat --find-renames=10% <Basis>..HEAD -- AGENTS.md .claude harness tools docs/plan/adr \
  ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' \
  | awk '{a+=$1; b+=$2} END {print a+b+0}'
```

Basis ist der Commit, der diese Datei nach `in-progress/` bewegt (sein Hash steht beim Übergang hier). Die Pfade sind `AGENTS.md`, `.claude`, `harness`, `tools` und `docs/plan/adr`; unter `docs/plan/adr` und `harness/vertraege` ist **nichts** erwartet, weil ADR und Vertrag vor der Basis liegen (Start-Trigger) — ein Nachtrag dort zählt trotzdem mit, deshalb stehen die Pfade im Befehl. Register-Belege und dieser Plan liegen unter `docs/plan/planning` und zählen nicht.

**Marke 900**, **Meldeschwelle 700**. Erreicht die Summe 700, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`); über 900 geht der Slice zurück zur Zerlegung (§Rückführungen).

**Schätzung des Planners (2026-10-21): 490 bis 900, Mitte 695** (77 % der Marke; die Mitte liegt damit über 70 %, und der Schnitt ist **vorgemerkt**, siehe Rückführung unten). Gerechnet in der Metrik oben, mit dem Faktor für Gegenprobe-Code (`AGENTS.md` §3.14: rund das Doppelte des Produktcodes):

| Teil | Schätzung |
|---|---|
| Wächter (`.claude/hooks/…`, ggf. awk-Extraktor) | 120 bis 200 |
| `.claude/settings.json` | 8 bis 12 |
| Gegenprobe (`tools/harness/…-gegenprobe.sh`) — Faktor 2 auf den Wächter | 250 bis 450 |
| Fragment `harness/mk/…` mit Liste | 15 bis 30 |
| `harness/README.md` (zwei Zeilen) | 5 bis 15 |
| `AGENTS.md` (neue harte Regel) | 30 bis 60 |
| sechs Rollen-Dateien, je ein Satz | 25 bis 50 |
| Reserve | 40 bis 80 |

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik oben liegt **über 900**, oder der Architect entscheidet eine Randform aus §6 so, dass ein vierter Liefer-Punkt entsteht (etwa ein Freigabe-Mechanismus, der die Grenze hebt, oder ein zweiter Anker für die Hauptsitzung). **Vorgemerkter Schnitt** (die Mitte der Schätzung liegt über 70 % der Marke): Liefer-Punkt 3 — die Regel in `AGENTS.md` und den sechs Rollen-Typen — verlässt den Slice und geht in einen eigenen Slice direkt dahinter, in derselben Welle; Wächter und Gegenprobe bleiben zusammen, weil eine Zusage ohne ihre Mutation nicht geliefert ist (`AGENTS.md` §3.10). Die Kennung des Nehmers vergibt der Planner beim Schnitt; hier ist keine genannt und nichts zugewiesen (`AGENTS.md` §3.13).
- `in-progress` → `open` (blockiert — Carveout?): Der Architect kann eine Randform aus §6 nicht entscheiden, ohne eine Entscheidung des Nutzers (etwa Fail-Richtung oder die Liste der durchgelassenen Aufrufe), und die liegt nicht vor; oder der Span-Strom trägt den Anker nicht, den der Wächter braucht (§6, Randform 3 und 4b), und die Erfassungsschicht ist aus §1 ausgeschlossen — dann ist der Befund ein Vorgang am Werkzeug und der Slice wartet; oder der Nutzer nimmt die Entscheidung vom 2026-10-21 zurück.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Beobachtbar: (1) `make gates` und `make docs-check` grün, und die Gegenprobe aus Liefer-Punkt 2 in einem Zug grün mit ihrer Zeit unter 600 s in §7; (2) der Beleg des Wächters im Lauf in §7 — unter der Grenze keine Ausgabe, an der Grenze die Block-Meldung, je gegen einen Attrappen-Strom, Eingabe und Ausgabe wörtlich; dazu der bestätigte Wortlaut der Regel in `AGENTS.md` und in den sechs Rollen-Dateien (`grep` findet die Kennung). Dazu die Closure-Notiz mit Lerneintrag und der Register-Eintrag aus §2.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

### Randformen — Fragen an den Architect, vor dem ersten Code-Commit zu entscheiden

Entschieden wird jede in `SPEC-070` (`AGENTS.md` §3.12: in `harness/vertraege/`, auch für Harness-Werkzeuge; eine ADR trägt nur Entscheidung und Gründe und zeigt mit `Schärft:` dorthin). Der Planner stellt die Fragen und beantwortet sie nicht.

1. **Was geblockt wird und was durchgelassen.** Die Rolle muss bei der Grenze noch committen, pushen und berichten können — sonst erzwingt der Wächter genau das Gegenteil der Regel („anhalten und berichten"). Zu entscheiden: welcher `matcher` im Hook-Eintrag (nur `Bash` wie der bestehende Wächter, oder jedes Werkzeug — bei nur `Bash` arbeitet eine Rolle mit `Read`, `Write` und `Edit` weiter); welche Befehle frei bleiben (`git commit`, `git push`, `git mv`, `git status`, `git diff`, `make slice-mv`?) und woran sie erkannt werden (Kommandoposition je Segment wie beim bestehenden Wächter); ob das Werkzeug, mit dem eine Rolle ihren Bericht abgibt, überhaupt unter den Wächter fällt. **Eine Liste, die zu eng ist, nimmt der Rolle den Bericht; eine, die zu weit ist, lässt die Arbeit weiterlaufen.**
2. **Die Fail-Richtung.** Empfehlung des Planners: **fail-open** — die Erfassung liegt gitignored, und `harness/erfassung-feldliste.md` sagt, dass ein frischer Klon das Programm nicht hat („Die Liste gilt auch dann, wenn gerade nichts erfasst wird"); alle Arbeit zu blocken, weil die Erfassung fehlt, ist schlimmer als nicht zu begrenzen. Der Wächter nach `AGENTS.md` §3.1 failt umgekehrt (**closed**, „bei Parse-Zweifel fail-closed", und der gebackene Boden ist ausdrücklich „NIE fail-open"). Die Asymmetrie zweier Hooks im selben Verzeichnis braucht ihre Begründung im Vertrag und in der ADR, nicht im Plan. Zu entscheiden je Lage getrennt: kein Verzeichnis `.harness/state/spans` · Verzeichnis ohne passenden Strom · Strom ohne Anker-Zeile · Strom mit beschädigter Zeile (Parse-Zweifel) · fehlendes `awk`.
3. **Woraus der Wächter liest, und was sein Anker ist.** Der Befehl: aus `tool_input.command` der Hook-Eingabe wie beim bestehenden Wächter — das Span-Feld `program` trägt nur das erste Token, und in der Messung führen `until`, `cd` und `cat`, weil `cd … && make …` als `cd` erfasst ist. Der **Anker** (Beginn des Laufs) ist zu entscheiden: die Zeile `event: SubagentStart` des eigenen Stroms (`ts`), die erste Zeile des Stroms, oder die Mtime der Datei; und **wie der Wächter seinen eigenen Strom findet** — der Dateiname ist `<Sitzung mit Unterstrichen>-<agent>.jsonl`, eine Sitzung trägt viele Ströme (im gemessenen Bestand über 20), und ob die Hook-Eingabe die Agenten-Kennung überhaupt trägt, ist zu prüfen, nicht anzunehmen. Dazu: **Wanduhr oder Werkzeugzeit** — der gemessene Lauf hatte 4,34 h Wanduhr und 48 % Befehlszeit; „max. 1 Stunde" ist in der einen Rechnung erreicht, in der anderen nicht. Dazu: Zeitform und Zeitzone (`ts` ist UTC mit `Z`), welches Werkzeug die Gegenwart liefert (`date -u`), und was bei einer Uhr gilt, die zurückspringt.
4. **Die Lücken, die der Vertrag benennen muss, statt sie zu verschweigen** (`AGENTS.md` §3.11): (a) Die **Hauptsitzung** hat keinen `SubagentStart` — das Budget deckt die sechs Rollen exakt, ein Command in der Hauptsitzung braucht einen anderen Anker (aus §1 ausgeschlossen, im Vertrag zu benennen). (b) **Innerhalb eines langen Vordergrund-Befehls feuert kein Hook** — die Überschreitung kann einen Befehl lang dauern, und die Zusage lautet darum nicht „nach einer Stunde", sondern „beim ersten Aufruf ab der Grenze". (c) Der Wächter wirkt nur, wo die Erfassung liegt (Folge von Randform 2). (d) Die **Plattform**: Der Hook läuft, wo Claude Code läuft — nach `LH-QA-03` sind das Linux, macOS und Windows; geprüft ist er dort, wo die Gegenprobe fährt. Je übrige Plattform eine Entscheidung: Grenze im Vertrag („geprüft nur unter Linux") oder eine Zusage mit Fall (`AGENTS.md` §3.12, Plattform-Hälfte; `BEO-REPO/verhalten-nur-unter-linux-geprueft`) — betroffen sind `date`, `stat` und der awk-Dialekt.
5. **Aufrufarten und Umgebung des Trägers** (`AGENTS.md` §3.12, Aufrufarten): Der Wächter liefert ein Urteil, das der Lauf glaubt. Zu entscheiden: Zeitgrenze des Hooks im Eintrag (der bestehende Wächter hat keine, der Emitter 5 s) und was gilt, wenn sie reißt; Verhalten bei leerer Umgebung und aus einem Unterverzeichnis (wie der Stop-Hook es zusagt); Verhalten im Dev-Container gegenüber dem Host (`BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`); Verhalten in einem **Worktree** — die Spans folgen dem Hauptbaum, nicht dem Klon (`BEO-REPO/span-zuordnung-folgt-dem-hauptbaum-nicht-dem-klon`, offen), und der Implementer arbeitet künftig im Worktree.
6. **Ob es einen Weg über die Grenze gibt.** Zu entscheiden: Darf der Nutzer die Grenze für einen Lauf heben (Umgebungsvariable, Datei, Satz im Auftrag)? Ohne Weg ist ein Lauf, der legitim länger braucht (ein Gegenprobe-Lauf von 19 min plus Arbeit), nach einer Stunde blockiert; mit Weg ist die Grenze eine Bitte. Der Planner nennt beide Hälften und entscheidet keine.

### Risiken

- **Der Wächter nimmt der Rolle den Bericht** (zu enge Liste der durchgelassenen Aufrufe, Randform 1): Der Lauf wird geblockt, kann aber nicht mehr committen oder berichten, und die Arbeit einer Stunde liegt nicht erfasst da — genau der Schaden, den die Regel verhindern soll. — **Ausgang:** <eingetreten / entfallen / weiter offen — bei der Closure>
- **Fail-open macht den Wächter dort wirkungslos, wo er gebraucht wird** (Randform 2): Fehlt die Erfassung (frischer Klon, `make span-clean` kurz davor, ein Lauf im Worktree mit anderem Pfad), läuft der Lauf unbegrenzt weiter, und niemand merkt, dass der Wächter nichts tut. Gegenmittel im Vertrag: eine Meldung des Hooks, die seine eigene Leere nennt (wie `span-report` seine Abdeckung zuerst nennt). — **Ausgang:** <bei der Closure>
- **Der Anker ist nicht eindeutig** (Randform 3): Findet der Hook seinen Strom nicht sicher, rechnet er die Zeit eines fremden Agenten derselben Sitzung — er blockt zu früh oder zu spät. Ein Fall je Lage in der Gegenprobe ist Pflicht (Liefer-Punkt 2). — **Ausgang:** <bei der Closure>
- **Latenz vor jedem Werkzeug-Aufruf:** Der Wächter läuft wie der bestehende Guard vor jedem Aufruf („der Guard laeuft vor JEDEM Bash-Call, Latenz zaehlt") und liest dabei eine Datei, die mit dem Lauf wächst (251 Zeilen nach 4,34 h). Liest er den ganzen Strom, wächst die Latenz mit der Laufzeit. — **Ausgang:** <bei der Closure>
- **Die Zeit der Gegenprobe:** Fälle mit Zeitsprüngen verführen zu `sleep`; feste Wartezeiten sind in diesem Repo ein eigener Befund (`BEO-REPO/feste-wartezeit-statt-beobachtung-im-test`), und [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) setzt 600 s. Die Zeit kommt aus dem Attrappen-Strom (ein `ts` in der Vergangenheit), nicht aus Warten. — **Ausgang:** <bei der Closure>
- **Zwei Hooks mit entgegengesetzter Fail-Richtung im selben Verzeichnis** sind eine Lesart-Falle für jeden, der später einen dritten schreibt (`BEO-REPO/harness-lesart-ohne-entscheidungsort`). Gegenmittel: die Begründung steht in der ADR und im Kopf beider Dateien, nicht nur in einem von beiden. — **Ausgang:** <bei der Closure>
- **Die Regel ohne den Mechanismus und der Mechanismus ohne die Regel:** Wird Liefer-Punkt 3 beim vorgemerkten Schnitt (§4) abgegeben, wirkt der Wächter, bevor eine Rolle seine Handlung in ihrer eigenen Datei liest. Gegenmittel: Die Block-Meldung des Wächters nennt die Handlung selbst (Liefer-Punkt 1); der Folge-Slice steht in derselben Welle direkt dahinter. — **Ausgang:** <bei der Closure>
- **Offener Punkt für den Nutzer (a):** `AGENTS.md` §3.9 nennt als Harness-Artefakte „`.claude/commands/`, `.harness/skills/`, Hooks" — `.claude/agents/` steht **nicht** in der Liste, obwohl dieser Slice sie ändert und obwohl die Rollen-Typen genauso jeden Lauf steuern. Der Planner entscheidet das nicht; zu klären ist, ob die Liste ergänzt wird und in welchem Slice (hier, oder in `slice-harness-regeln-nach-v1-abschluss`, der den Regel-Block ohnehin anfasst). — **Ausgang:** <bei der Closure>
- **Offener Punkt für den Nutzer (b):** Die Stunde ist eine **zweite Rückführungs-Marke** neben der Zeilenzahl — ein Slice, der in einer Stunde nicht umsetzbar ist, ist zu groß, und das ist eine Aussage über den **Schnitt**, nicht über den Lauf. Gehört sie in diesen Slice (dann wächst Liefer-Punkt 3) oder in `slice-harness-groessen-metrik`, der `AGENTS.md` §3.14 ohnehin neu schreibt? Der Planner legt die Frage vor und entscheidet sie nicht; `slice-harness-groessen-metrik` ist hier **nicht** als Adresse genannt und bekommt keine Sendung (`AGENTS.md` §3.13). — **Ausgang:** <bei der Closure>

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren) · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln (das Feld `liegt in` steht **nur**, wenn mit diesem Slice
wirklich etwas verkörpert wurde).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Belege des Implementers:** <Mutationstabelle je Zusage von `SPEC-070` (Zusage · Mutation · roter Fall), die Zeit des vollen Gegenprobe-Laufs, der Lauf des Wächters gegen den Attrappen-Strom mit Eingabe und Ausgabe wörtlich, Datum und Fundstelle der Bestätigung des Wortlauts durch den Nutzer>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** `BEO-REPO/rollenlauf-ohne-zeitgrenze` neu angelegt, Beleg `evidence/slice-harness-rollenlauf-zeitgrenze.md`; `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` um die gemessene Zeit ergänzt.
- **Folge-Slices:** <die Wiederholungs-Grenze (§1) — vom Planner anzulegen, sobald `SPEC-070` entschieden ist; der vorgemerkte Schnitt aus §4, falls er eintrat>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <Anker · Folge-Slice · Register, Ergebnis>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — die **zwei vorgelagerten Schritte**
(in jedem Slice-Plan) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Eine Sub-Area, `REPO` (Default des gesamten Repos, `harness/conventions.md` §Modus-Deklaration pro Sub-Area, Greenfield). Die drei Achsen sind erfüllt (eigene Artefakte unter `.claude/` und `harness/`, eigener Maßstab im Review, eigene Konventionen in `AGENTS.md`); eine feinere Ausdifferenzierung („Durchsetzungsschicht") deklariert dieses Repo nicht, und ein Slice ist kein Ort, eine zu erfinden.

**Vorgelagert — offene Beobachtungen sichten** (gemergter Stand `main`, 2026-10-21):

- **Kein vorhandener Eintrag trägt den Anlass** — einen Rollenlauf ohne Zeitgrenze. Das ist die Antwort der Sichtung, und sie ist der Grund für den neuen Eintrag `BEO-REPO/rollenlauf-ohne-zeitgrenze` in der DoD (§2). `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (14 Belege, verkörpert) misst die Laufzeit von **Gates und Gegenproben**, nicht die eines Agentenlaufs; `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` (geplant, `slice-harness-lint-warten-ohne-frist`) betrifft einen hängenden Lauf im Gate, nicht die Dauer einer Rolle.
- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` — **14 Belege, Stand verkörpert**: Dieser Slice fügt eine Gegenprobe hinzu. Sie hängt an `GEGENPROBEN_CONTAINER`, nicht an `GATE_CHECKS`, und fährt ohne Docker — die Kette wächst nicht, die Fahrt von `make gegenproben` um sie. Beleg bei der Closure (§2).
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` — **verkörpert in `AGENTS.md` §3.12** (Plattform-Hälfte): Der Wächter läuft, wo Claude Code läuft; die Plattform ist deshalb Randform 4d, zu entscheiden vor dem Code.
- `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` und `BEO-REPO/span-zuordnung-folgt-dem-hauptbaum-nicht-dem-klon` (offen) — beide treffen die Datenquelle dieses Slice: Der Hook läuft in der Umgebung des Agenten, die Gegenprobe im Container, und die Ströme folgen dem Hauptbaum, nicht dem Worktree. Beides steht als Randform 5 beim Architect; erreicht der zweite Eintrag mit diesem Slice sein drittes Auftreten, ist er keine Notiz mehr, sondern eine Lücke und braucht einen eigenen Folge-Slice (bei der Closure zu prüfen, §2).
- `BEO-REPO/meldeschwelle-ohne-artefakt`, `BEO-REPO/groessenprobe-ohne-messbaren-stand` (verkörpert) und `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` — §4 nennt Metrik, Marke, Meldeschwelle und rechnet die Gegenprobe mit dem Faktor 2.
- `BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet` — getragen: Der Kopf und §4 trennen die **Entscheidung des Nutzers** (Sache: Grenze, Mechanismus, Wiederholungs-Grenze) von der **Empfehlung des Planners** (Welle, Stelle, Schnitt, Marke, Fail-Richtung).
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` — getragen: Dieser Plan nennt **keine** Kennung als Adresse und weist nichts zu (§1 Wiederholungs-Grenze, §4 vorgemerkter Schnitt, §6 Punkt (b)); damit ist keine Annahme einzutragen und keine Nachzählung fällig.

Keiner der gesichteten Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Zählung nach `AGENTS.md` §3.13 und §3.14** (Planner, 2026-10-21, im Commit, der diesen Plan anlegt): **Liefer-Punkte 3** (Wächter mit Eintrag und Index-Zeile | Gegenprobe | Regel in `AGENTS.md` und den sechs Rollen-Typen). **Schichten 1** — der Harness (Teilung nach `slice-harness-d-check-v0-85` §8, `AGENTS.md` §3.13: Gate-Konfiguration und Gate-Werkzeuge, `Dockerfile`, Agenten-Anweisungen `AGENTS.md`, Skills, Commands und `harness/` sind **eine** Schicht; dazu zählt der Planner die Rollen-Typen unter `.claude/agents/`, weil sie Agenten-Anweisungen sind — dass `AGENTS.md` §3.9 sie in seiner Liste nicht nennt, ist der offene Punkt (a) in §6). Die Entscheidungen (`docs/plan/adr/`) sind eine zweite Schicht, liegen aber **vor** der Basis (Start-Trigger, §4) und gehören damit nicht in diesen Diff; käme ein Nachtrag an der ADR hinzu, sind es zwei. Die Planung (dieser Plan, Welle, Roadmap, Register) ist keine Schicht. **Schätzung 490 bis 900, Mitte 695** bei **Marke 900** und Meldeschwelle 700 (§4) — unter drei Liefer-Punkten und unter zwei Schichten, aber mit der Mitte bei 77 % der Marke: Der Schnitt ist vorgemerkt (§4 *Rückführungen*), nicht vollzogen. `grep` nach der Kennung dieses Slice findet ihn im Welle-Plan [welle-harness-regeln](../welle-harness-regeln.md) (§4 und §5) und im Drift-Log der Roadmap; in keinem anderen Plan, weil dieser Slice nichts zuweist.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
