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

**Bezug:** kein Anspruch des Produkts (`LH-*`) — der Slice erzwingt eine Regel über Agentenläufe und ändert kein Verhalten des Werkzeugs. **Entscheidungen des Nutzers** (jeweils vom Koordinator übermittelt): 2026-10-21, wörtlich „Wir haben rollen (siehe .claude/agents) und commands -- max. 1 Stunde“, und „ok“ auf den Entwurf des Mechanismus (Wächter im Werkzeug-Aufruf, bash und awk, **kein** Hintergrundprozess, Datenquelle die vorhandene Erfassung); 2026-10-22, „deine empfehlungen“: die Stunde ist zugleich zweite Rückführungs-Marke und steht in der Regel (`slice-harness-rollenlauf-zeitgrenze-regel`), `AGENTS.md` §3.9 führt `.claude/agents/` auf; 2026-10-22, wörtlich „2. ja“: Der Wächter deckt nur Rollenläufe, und der Regeltext nennt Commands und Hauptsitzung als nicht erzwungen; 2026-10-22: [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) ist **vom Nutzer angenommen** und steht auf `Accepted`. Aktive ADRs, gelesen: [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Zeitbudget 600 s je Gegenprobe — gilt auch für die Gegenprobe dieses Slice), [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (die Gegenprobe hängt an `GEGENPROBEN_CONTAINER`, nicht an `GATE_CHECKS`), [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (die Hooks dieses Repos und ihr Zusammenspiel), [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (Metrik der Größenmarke in §4).

**Berührte Spec-Stellen:** [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) (liegt vor, vom Architect; der Vertrag des Wächters mit den Randformen der sechs Fragen dieses Plans; **die Zuordnung seiner Punkte und Fälle zu den zwei Stufen steht im Vertrag**, Abschnitt *Stufen*) · [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) (gelesen: Name und Liste der Gegenprobe, Punkte 1 bis 3; Aufnahme nach Punkt 16, deren Zeile der Architect nachträgt) · [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) (gelesen: der Stop-Hook und die Ketten — dieser Slice ändert beide nicht) · [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 (6) (gelesen: Ort des Belegs eines Nachweises, dem der Beleg im Agentenlauf folgt).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-21 (Schnitt und Nachzug nach `SPEC-070`: 2026-10-22; Nachzug nach der Zuordnung der Stufen im Vertrag: 2026-10-22).

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
5. **Gewählt: nach Wirkung und Härtung; der Architect hat den Schnitt bestätigt und zugeordnet** ([`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md), Abschnitt *Stufen*, Entscheidung vom 2026-10-22; [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) bleibt unverändert, kein `Supersedes`). Stufe 1 (dieser Slice) liefert, was die Stunde wirken lässt: wer begrenzt ist, Strom, Anker (in der Gestalt) und Zeit, die ganze Sperrmenge, Ausgabe und Exit, kein Weg über die Grenze, die Wurzel, die **Freigabe in den Lagen, die ein Fall zeigt**, den Nachweis im Lauf, die Gegenprobe samt Aufnahme. Stufe 2 liefert, was den Hook an den Rändern hält und seine Lage nennt (Zeile auf stderr, Reihenfolge der Gründe, Bereich des Ankers, Umgebung). Der Vertrag zählt rund **78 Fälle** für Stufe 1 und rund 29 für Stufe 2 (zusammen 107). Geschätzt in der Metrik (§4, nachgerechnet am 2026-10-22): Stufe 1 **476 bis 815, Mitte 645, das sind 72 % der Marke 900** (der Architect schätzt 470 bis 800, Mitte 630, 70 %), Stufe 2 200 bis 405 (Mitte 302, 60 % der Marke 500). **Die Lage der ersten Stufe zur 70 %-Linie und zur Meldeschwelle steht in §4.**

**Was nach dieser Stufe allein zugesagt ist, und was nicht** (der Vertrag sagt es im Abschnitt *Stufen*; der Plan zeigt dorthin und wiederholt es nicht): Zugesagt ist die **Freigabe mit Exit 0 und ohne Ausgabe auf stdout** in den Lagen, die ein Fall zeigt: Hauptsitzung, kein Verzeichnis, kein Strom, Strom ohne Start oder leer, Anker in falscher Gestalt, Eingabe leer, Eingabe kein JSON, `awk` fehlt, ein Befehl, den der Extraktor nicht liest, ein `agent_id` unterhalb der obersten Ebene, ein Unicode-Escape in einem anderen Schlüssel. **Nicht zugesagt** ist jede Zeile auf stderr (nur die Hauptsitzung bleibt dort leer), die Freigabe bei jeder anderen Unsicherheit, ein Anker mit einem Wert außerhalb des Bereichs (weder Freigabe noch Sperre), die Reihenfolge der Gründe, die leere Umgebung, die Eingabe bis 1 MB und gleichzeitige Aufrufe. Die Sätze „Jede Unsicherheit des Hooks ist eine Freigabe“ und „Exit 0 in jedem Zweig“ gelten **uneingeschränkt erst nach Stufe 2**. Die Richtung der zugesagten Lagen ist die sichere (ein verfehlter Block heißt nur, dass die Grenze für diesen Lauf nicht greift, und die Regel an den Leser bleibt, `slice-harness-rollenlauf-zeitgrenze-regel`). **Folge für jeden Satz über den Hook:** `harness/README.md`, der Kommentar im Kopf des Hooks und des Fragments, Handbuch und Regeltext nennen nach Stufe 1 nur diese Zusagen, **keine Selbstauskunft auf stderr und keine allgemeine Freigabe** (Vertrag, *Zwischen den Stufen* (d); DoD und Randform 8).

**Sendungen** (eingetragen vom Planner im Commit des Schnitts, `AGENTS.md` §3.13): an `slice-harness-rollenlauf-zeitgrenze-regel` (die Regel in `AGENTS.md` und den sechs Rollen), an `slice-harness-rollenlauf-zeitgrenze-haertung` (die Punkte der zweiten Stufe) und an `slice-harness-rollenlauf-wiederholungsgrenze` (die Wiederholungs-Grenze); je in deren §1 und DoD angenommen, mit Nachzählung in deren §8.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die **Regel** in `AGENTS.md` (die Grenze als harte Regel, die Stunde als zweite Rückführungs-Marke, der Satz zum Verifier, die Zeile in §3.9) und der **Satz in den sechs Rollen-Typen** — Klasse 1: `slice-harness-rollenlauf-zeitgrenze-regel`. Grund: andere Dateien, anderer Vorgang; der Wächter ist ein Werkzeug, die Regel ein Text an den Leser. Die Block-Meldung des Hooks nennt die Handlung selbst (Punkt 6 von [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md)), damit sie nicht an der Regel hängt.
- **Stufe 2** — Härtung und Selbstauskunft — Klasse 1: `slice-harness-rollenlauf-zeitgrenze-haertung`. Grund: Aufteilung des Vertrags nach Wirkung und Härtung (Abschnitt *Stufen* des Vertrags), damit jede Lieferung unter ihrer Marke bleibt; die Zeile auf stderr, die Reihenfolge der Gründe, der Bereich des Ankers und die Fälle der Umgebung (Punkt 8 (2), (3), (5)) gehören dorthin.
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

- [ ] **Liefer-Punkt 1 — der Wächter, Stufe 1:** `.claude/hooks/pretooluse-zeitgrenze.sh` mit dem Leser `tools/harness/hook-felder.awk` (neu), reines bash und awk (kein `jq`, kein `node`, kein Netz), eingetragen in `.claude/settings.json` unter `PreToolUse` in einem eigenen Eintrag (`matcher` `""`, `timeout` 5; der Eintrag des Wächters nach `AGENTS.md` §3.1 bleibt unberührt) und mit den Zeilen in `harness/README.md` als **Werkzeug, kein Gate**. Er erfüllt die Punkte, die der Abschnitt *Stufen* von [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) der ersten Stufe gibt: Eintrag, Leser, Befehl (Punkt 1), wer begrenzt ist (Punkt 2), Strom, Anker und Zeit mit der Gestalt des Ankers (Punkt 3 ohne den Bereich von (e)), die Sperrmenge (Punkt 5), Ausgabe und Exit ohne die Zeile auf stderr (Punkt 6), die Freigabe in den Lagen des Vertrags (Punkt 7, sieben Fälle), Wurzel und Worktree (Punkt 8 (1) und (4)), kein Weg über die Grenze (Punkt 9). Jede Randform dieser Punkte hat der Vertrag entschieden. **Jeder Satz über den Hook nennt nur diese Zusagen: keine Selbstauskunft auf stderr, keine Freigabe bei jeder Unsicherheit** — die Zeile des Wächters in `harness/README.md`, die dieser Liefer-Punkt schreibt, und der Kommentar im Kopf des Hooks (Vertrag, *Zwischen den Stufen* (d); `AGENTS.md` §3.11; Randform 8); die Zeile der Gegenprobe in `harness/README.md` schreibt Liefer-Punkt 2.
- [ ] **Liefer-Punkt 2 — die Gegenprobe, Stufe 1:** ein Ziel `make zeitgrenze-gegenprobe` mit Fragment `harness/mk/zeitgrenze.mk` und Skript `tools/dev/harness/zeitgrenze-gegenprobe.sh` (Namen: [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) Punkt 14), in `GEGENPROBEN_CONTAINER` ([`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 1 bis 3: bash und awk, ohne Docker), **nicht** in `GATE_CHECKS` und nicht in der Ausnahmeliste der schnellen Kette ([ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)). Sie fährt gegen **Attrappen-Ströme** (JSONL in `mktemp -d`, kein Agentenlauf, keine echte Erfassung) und eine Attrappe von `date` und trägt **je Zusage der ersten Stufe und je Element ihrer Aufzählungen einen Fall und eine Mutation, die nur ihn rot färbt** (`AGENTS.md` §3.10; die Namen der Fälle tragen die Nummer des Punktes). **Die Fälle der ersten Stufe** nach dem Abschnitt *Stufen* von `SPEC-070` (78): Punkt 1 (`p1-eintrag`, `p1-schreibt-nichts`, `p1-befehl-nicht-aus-program`); Punkt 2 (sieben, darunter `p2-agent-in-tool-input`); Punkt 3 (23: alle `p3-*` ohne die sechs Bereichsformen, dazu die sieben der Gestalt `p3-form-versatz`, `-bruchteil`, `-klein-z`, `-leerzeichen`, `-leer`, `-fehlt`, `-text-davor`); Punkt 5 (28, darunter `p5-unicode-in-tool-input` und `p5-bash-unlesbar`); Punkt 6 (`p6-block-text`, `p6-frei-leer`, `p6-exit-0-block`, `p6-exit-0-frei`, `p6-exit-0-fehler`); Punkt 7 (sieben: `p7-kein-verzeichnis`, `p7-kein-strom`, `p7-strom-leer`, `p7-kein-start`, `p7-eingabe-leer`, `p7-eingabe-kein-json`, `p7-kein-awk`, je Exit 0, stdout leer und eine Mutation „Block statt Freigabe“); Punkt 8 (`p8-wurzel-projektdir`, `p8-wurzel-ort-der-datei`, `p8-unterverzeichnis`); Punkt 9 (`p9-fortsetzung`); Punkt 14 (`p14-aufnahme` und die Fälle von [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16 (9)). **Die Erwartung jedes Falls dieser Stufe ist Exit 0 und stdout wie im Vertrag; über stderr ist nichts erwartet, außer dass es bei der Hauptsitzung leer bleibt.** Fälle mit Grund und Zeile (`p5-bash-unlesbar`, die `p3-form-*` und die `p7`-Fälle) prüfen hier die Freigabe, nicht die Zeile; Stufe 2 ergänzt die Erwartung und ändert weder Namen noch erwartete Mengen roter Fälle. Die unveränderte Kopie ist grün, jeder Fall räumt auf, **kein Fall wartet** (die Zeit kommt aus dem Anker, nicht aus `sleep`), und der volle Lauf liegt **unter 600 s** ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)); die gemessene Zeit steht in §7. Der Code von `container-einrichten-gegenprobe.sh` darf übernommen werden. `harness/README.md` §Sensors führt die Zeile der Gegenprobe.

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
| `harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md`, `docs/plan/adr/0058-…md` | liegen vor (Architect); die **Zuordnung der Punkte und Fälle zu den zwei Stufen** steht im Vertrag (Abschnitt *Stufen*) | Zusagen und Randformen entschieden (`AGENTS.md` §3.12) |
| `.claude/hooks/pretooluse-zeitgrenze.sh` | neu | Liefer-Punkt 1: der Wächter; Form und Aufbau nach `.claude/hooks/pretooluse-command-guard.sh` |
| `tools/harness/hook-felder.awk` | neu | Liefer-Punkt 1: der Leser der Felder der Eingabe (`session_id`, `agent_id`, `tool_name`); den Befehl liest der Hook mit `tools/harness/extract-command.awk` |
| `.claude/settings.json` | update | Liefer-Punkt 1: eigener Eintrag unter `PreToolUse` (`matcher` `""`, `timeout` 5) |
| `harness/README.md` §Werkzeuge und §Sensors | update | Zeile des Wächters (Werkzeug, kein Gate) und Zeile der Gegenprobe, **nach Stufe 1 ohne Selbstauskunft und ohne Freigabe bei jeder Unsicherheit** (Randform 8); der Gate-Index steht **einmal** (`AGENTS.md` §4) |
| `harness/mk/zeitgrenze.mk` | neu | Liefer-Punkt 2: das Ziel, Eintrag in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS` |
| `tools/dev/harness/zeitgrenze-gegenprobe.sh` | neu | Liefer-Punkt 2: Fälle und Mutationen gegen Attrappen-Ströme und eine Attrappe von `date`, Aufräumen je Fall |
| `docs/reviews/<Datum>-nachweis-slice-harness-rollenlauf-zeitgrenze.md` | neu (Implementer, Lauf-Beleg) | der Nachweis im Agentenlauf; von der Größenzählung ausgenommen ([ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md)) |
| `harness/vertraege/SPEC-060-gegenproben.md` Punkt 16 | update **durch den Architect** | die Zeile der Gegenprobe mit `aufgenommen` |
| `docs/plan/planning/observations/BEO-REPO/rollenlauf-ohne-zeitgrenze/`, `…/gate-laufzeit-waechst-mit-gegenprobe/` | update bei der Closure | Beleg und gemessene Zeit (kein Liefer-Punkt, DoD) |

- **Reihenfolge:** (1) die Zuordnung der Punkte von `SPEC-070` zu den zwei Stufen steht im Vertrag (erfüllt, vor dem Übergang nach `in-progress`); (2) der Wächter mit seinem Eintrag in `.claude/settings.json` — der **Nachweis im Agentenlauf** folgt unmittelbar, weil er die Annahme belegt, auf der der Vertrag ruht (Eingabe mit `agent_id`, Form des Blocks); (3) die Gegenprobe mit ihren Fällen und Mutationen, je Zusage eine, vom Implementer selbst rot gesehen; (4) Closure mit Register und der Zeile in `SPEC-060`. Jeder Commit nennt die Kennung dieses Slice.
- **Was der Implementer nicht entscheidet:** Trifft er eine Randform, die `SPEC-070` nicht nennt, hält er an und gibt sie dem Architect zurück (`.claude/commands/implement-slice.md`, Randform-Rückgabe; `AGENTS.md` §3.12). Den Vertragstext ändert er nicht, auch nicht die Zeile in `SPEC-060`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): WIP-Limit frei; **Vertrag und ADR liegen vor** — `SPEC-070` vom Architect, [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) vom Nutzer angenommen (`Accepted`) — **und die Zuordnung der Punkte zu den zwei Stufen steht im Vertrag** (Abschnitt *Stufen*), jede Randform der ersten Stufe mit einer Antwort (`AGENTS.md` §3.12). Reihenfolge: **hinter** `slice-harness-rollenlauf-zeitgrenze-regel` (**Empfehlung des Planners**, kein Start-Trigger: Die Regel an den Leser gilt ab ihrer Lieferung und wartet nicht auf den Wächter, und der Wächter nicht auf sie); davor liegt `slice-harness-groessen-metrik`. Der Übergang `open` → `next` setzt `Verantwortlich:`. **Größe gegen die Regel von `AGENTS.md` §3.14** (Entscheidung des Nutzers vom 2026-10-22, wörtlich „Beim Start je Plan entscheiden (Empfohlen)“): Die Mitte der Schätzung liegt bei 71,7 % der Marke 900, über 70 % (Mitte 645,5; die Stufe 1 ist die Hälfte nach einem Schnitt, für die dieselbe Linie gilt); es wird nicht rückwirkend geschnitten. Beim Übergang nach `next/`, spätestens vor dem Anspruch, rechnet der Planner die Mitte erneut gegen 70 % der Marke und schneidet nach der Regel oder legt dem Nutzer die datierte Entscheidung „ohne weiteren Schnitt“ vor.

**Größe.** Metrik nach [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (der Befehl steht in `AGENTS.md` §3.14 nach `slice-harness-groessen-metrik`; wiederholt dieser Plan ihn, gilt der von dort):

```
git diff --numstat --find-renames=10% <Basis>..HEAD -- AGENTS.md .claude harness tools \
  ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' \
  | awk '{a+=$1; b+=$2} END {print a+b+0}'
```

Basis ist der Commit, der diese Datei nach `in-progress/` bewegt (sein Hash steht beim Übergang hier). Unter `harness/vertraege` ist nichts erwartet (Vertrag vor der Basis; ein Nachtrag zählt mit, deshalb steht `harness` im Befehl); der Beleg des Nachweises liegt unter `docs/reviews` und zählt nicht.

**Marke 900, Meldeschwelle 700.** Erreicht die Summe 700, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`); über 900 geht der Slice zurück zur Zerlegung.

**Schätzung des Planners (2026-10-22, nachgerechnet nach der Zuordnung der Stufen), in der Metrik oben, mit dem Faktor 3 bis 4,6 für Gegenprobe-Code gegenüber Produktcode** (Messwerte in `slice-harness-schaetzung-planner-faktor` §1: 2,9, 3,2 und 4,6): **476 bis 815, Mitte 645** (72 % der Marke). Die Gegenprobe ist mit 78 Fällen angesetzt, je Fall 4,4 bis 7,6 Zeilen mit seinen Mutationen (aus der Schätzung des Architect von 550 bis 950 für 125 Fälle):

| Teil | Schätzung |
|---|---|
| Wächter (`.claude/hooks/pretooluse-zeitgrenze.sh`) und Leser (`tools/harness/hook-felder.awk`) | 100 bis 160 |
| `.claude/settings.json` | 8 bis 12 |
| Fragment `harness/mk/zeitgrenze.mk` | 15 bis 30 |
| `harness/README.md` (zwei Zeilen) | 10 bis 20 |
| Gegenprobe (`tools/dev/harness/zeitgrenze-gegenprobe.sh`), Faktor 3,4 bis 3,7 auf den Wächter | 343 bis 593 |

**Lage zur 70 %-Linie und zur Meldeschwelle.** Die Mitte (645) liegt bei 72 % der Marke, an der Linie und knapp darüber (der Architect: 630, 70 %); die Obergrenze (815) liegt über der Meldeschwelle 700 und bei 91 % der Marke. **Die Marke 900 hält auf der Schätzung**, mit 9 % Abstand an der Obergrenze, und die Schätzungen dieses Repos lagen zuletzt unter den Messungen (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`); der Planner hebt sie nicht. **Geprüft, ob etwas in Stufe 2 gehört:** nichts, ohne eine Zusage der ersten Stufe zu streichen. Die sieben Fälle der Freigabe, die beiden Fälle der Zeitzone (`p3-zeitzone-*`, echte Uhr), `p6-exit-0-fehler` (Dateirechte, entfällt als root) und die Fälle der Wurzel schützen davor, dass der Hook einen Lauf aus Zweifel oder falscher Zeit sperrt, und die Sperrmenge ist als Ganzes zugesagt. Ein weiterer Schnitt wäre eine neue Zuordnung des Architect und ist **nicht** vorab benannt (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`); die Rückführung unten nennt die Bedingung. Die Entscheidung des Nutzers vom 2026-10-16 (nicht schneiden, solange die Schätzung die Marke selbst nicht überschreitet) gilt für diese Schätzung: Die Obergrenze liegt unter der Marke.

Der Beleg des Nachweises im Agentenlauf und die Zeile in `SPEC-060` zählen nicht (Pfad ausgenommen bzw. Vertragstext des Architect, höchstens 5 Zeilen).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt **über 900**, oder der Architect entscheidet eine Randform so, dass ein dritter Liefer-Punkt entsteht (etwa ein Freigabe-Mechanismus, der die Grenze hebt). Der Schnitt nach Lieferwert ist **vollzogen**, nicht vorgemerkt (§1); ein weiterer Schnitt wäre eine neue Zuordnung der Punkte durch den Architect. Erreicht die Summe die Meldeschwelle 700, hält der Implementer an und meldet dem Planner; die Entscheidung („weiter“ oder „schneiden“, mit dem gemessenen Stand) steht als datierte Zeile in §7.
- `in-progress` → `open` (blockiert): Der Nachweis im Agentenlauf zeigt, dass die Eingabe von `PreToolUse` kein `agent_id` trägt oder die Form des Blocks nicht greift (Befund an den Architect: Anker und Form fallen); oder der Nutzer nimmt die Entscheidung vom 2026-10-21 zurück.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Beobachtbar: (1) `make gates` und `make docs-check` grün, und die Gegenprobe aus Liefer-Punkt 2 (die 78 Fälle der ersten Stufe) in einem Zug grün mit ihrer Zeit unter 600 s in §7; (2) der Beleg des Wächters im Agentenlauf als Datei unter `docs/reviews/`, mit Pfad, Summe, Exit-Zeile und Kopf in §7 (Block-Meldung bei `echo make`, `Write` gesperrt, `git status` durchgelassen, Strom wiederhergestellt); dazu die Zeile `zeitgrenze-gegenprobe` mit `aufgenommen` in `SPEC-060` Punkt 16, die Closure-Notiz mit Lerneintrag und der Beleg im Register aus §2.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

### Randformen — entschieden in `SPEC-070`

Die sechs Fragen des früheren Plans (Sperrmenge und Abwicklung, Fail-Richtung, Befehl und Anker, die Lücken, Aufrufarten und Umgebung, Weg über die Grenze) und die Reihenfolge der Gründe, die Stelle der Aufrufe, die Aufruf-Voraussetzung, die pflichtigen Läufe und die Plattform sind in [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) entschieden; die Tabelle *Randform → Punkt* steht dort, und dieser Plan wiederholt sie nicht (`AGENTS.md` §3.12). Auch die **Stufung** hat der Architect entschieden (Randform 7).

7. **Die Zuordnung der Punkte zu zwei Lieferungen (Stufung) — entschieden.** Der Architect hat sie am 2026-10-22 im Vertrag festgehalten (Abschnitt *Stufen* vor *Gemessen*, dazu die Stufe an jeder Punkt-Überschrift und die Tabelle *Randform → Punkt*); **die Tabelle des Vertrags gilt**, der frühere Vorschlag des Planners an dieser Stelle entfällt. Gegenüber dem Vorschlag hat der Architect geändert: Die Freigabe ist in Stufe 1 nur in den Lagen zugesagt, die ein Fall zeigt (sieben neue `p7`-Fälle), die allgemeine Zusage gilt erst nach Stufe 2; `p2-agent-in-tool-input` und `p5-unicode-in-tool-input` zogen von Stufe 2 nach Stufe 1; Fälle mit Grund und stderr-Zeile haben getrennte Erwartungen je Stufe; der Bereich des Ankers ist reine Härtung.
8. **Sätze über den Hook nach Stufe 1.** `harness/README.md` (die Zeile des Wächters und der Gegenprobe), der Kommentar im Kopf des Hooks und des Fragments, Handbuch und Regeltext nennen nach Stufe 1 **nur die Zusagen der ersten Stufe** — keine Selbstauskunft auf stderr, keine allgemeine Freigabe (Vertrag, *Zwischen den Stufen* (d); `AGENTS.md` §3.11). Der Satz steht im DoD von Liefer-Punkt 1, weil dieser Liefer-Punkt die Zeile des Wächters und den Kommentar des Hooks schreibt; die Zeile der Gegenprobe schreibt Liefer-Punkt 2 mit derselben Beschränkung. Die Regel in `slice-harness-rollenlauf-zeitgrenze-regel` sagt zur Erzwingung nichts zu, was Stufe 1 nicht hält.
9. **Getrennte Erwartung je Fall.** Ein Fall mit Grund und Zeile prüft in Stufe 1 die Freigabe (Exit 0, stdout), nicht die Zeile; Stufe 2 ergänzt die Erwartung an stderr, lockert den Fall nie und ändert weder Namen noch erwartete Mengen roter Fälle (Vertrag, *Zwischen den Stufen* (b)). Eine Randform, die ein Fall der Stufe 2 aufwirft und die Stufe 1 berührt, geben beide Stufen an den Architect zurück ((f) dort).

**Risiken:**

- **Der Wächter nimmt der Rolle den Bericht** (zu enge Liste der durchgelassenen Aufrufe): Der Lauf wird gesperrt, kann aber nicht mehr committen oder berichten — genau der Schaden, den die Regel verhindern soll. — **Ausgang (vorgesehen):** entfallen, wenn die Fälle `p5-frei-*` und der Nachweis im Agentenlauf (`git status` durchgelassen) ihn zeigen; sonst eingetreten → Rückführung (§4).
- **Die Annahme des Vertrags fällt** (die Eingabe von `PreToolUse` trägt kein `agent_id`, die Form `decision: block` greift nicht): Der Hook wirkte in keinem Lauf, und die Gegenprobe zeigte es nicht. — **Ausgang (vorgesehen):** eingetreten → Randform-Rückgabe an den Architect, Slice nach `open` (§4); entfallen, wenn der Nachweis im Agentenlauf sie bestätigt.
- **Der Hook schweigt in dieser Stufe, wo der Strom fehlt** (über stderr ist nach Stufe 1 nichts zugesagt, die Auskunft liefert Stufe 2): Ein Lauf läuft unbegrenzt, und niemand sieht, dass der Wächter nichts tut. — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/rollenlauf-ohne-zeitgrenze`; Gegenmittel: Stufe 2 steht direkt dahinter, in derselben Welle.
- **Latenz vor jedem Werkzeug-Aufruf** (der Hook liest eine Datei, die mit dem Lauf wächst): gemessen sind 21 ms für einen Strom von 5 MB (`SPEC-070` Punkt 8 (3)). — **Ausgang (vorgesehen):** entfallen, mit der Zahl des Vertrags.
- **Die Zeit der Gegenprobe:** Fälle mit Zeitsprüngen verführen zu `sleep` (`BEO-REPO/feste-wartezeit-statt-beobachtung-im-test`), und [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) setzt 600 s. — **Ausgang (vorgesehen):** entfallen, wenn kein Fall wartet (die Zeit kommt aus dem Anker) und der volle Lauf unter 600 s liegt (§7).
- **Zwei Hooks mit entgegengesetzter Fail-Richtung im selben Verzeichnis** sind eine Lesart-Falle für den, der einen dritten schreibt (`BEO-REPO/harness-lesart-ohne-entscheidungsort`). — **Ausgang (vorgesehen):** entfallen, wenn die Begründung in der ADR und im Kopf beider Dateien steht, nicht nur in einer.
- **Ein Satz über den Hook sagt nach Stufe 1 mehr zu, als Stufe 1 zusagt** (die Selbstauskunft auf stderr oder „jede Unsicherheit ist eine Freigabe“ in `harness/README.md` oder im Kopf des Hooks): Die Zeile zusagte, was kein Fall der ersten Stufe zeigt (`AGENTS.md` §3.11). — **Ausgang (vorgesehen):** entfallen, wenn das Review die Sätze gegen *Nach Stufe 1 allein ist zugesagt* im Vertrag liest; ein Satz darüber ist *eingetreten* → Befund an den Planner.
- **Die Schätzung liegt unter der Messung:** Die Obergrenze (815) liegt bei 91 % der Marke 900 und über der Meldeschwelle 700; die Beobachtung `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` zählt 19 Belege. — **Ausgang (vorgesehen):** offen bis Closure; bei Erreichen von 700 meldet der Implementer (§4), über 900 geht der Slice zurück zur Zerlegung; liegt die Messung über der Obergrenze der Spanne, ist es ein Auftreten im Register.
- **Der Slice passt nicht in einen Rollenlauf von einer Stunde** (Hook, Gegenprobe mit 78 Fällen, Nachweis): Er läuft in mehreren Läufen, und jeder hält nach der Regel an und berichtet. — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/rollenlauf-ohne-zeitgrenze` (Beleg: die Zahl der Läufe); die Reihenfolge der Punkte lässt jeden Lauf an einer Naht enden.
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
- **Folge-Slices:** `slice-harness-rollenlauf-zeitgrenze-regel`, `slice-harness-rollenlauf-zeitgrenze-haertung`, `slice-harness-rollenlauf-wiederholungsgrenze` (Dateien in `open/`, `next/` oder `done/`).
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
- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` (3×, `geplant` mit `slice-harness-schaetzung-planner-faktor`): Dieser Plan hält das Muster nicht ganz ein: Stufe 1 liegt mit der Mitte bei 72 % der Marke (der Architect: 70 %), also an der Linie und nicht darunter; kein zweiter Schnitt ist vorab benannt, die Marke liegt über der Obergrenze der Schätzung (§4, *Lage zur 70 %-Linie*).
- `BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet` und `BEO-REPO/folge-slice-adresse-nimmt-nicht-an`: getragen — der Kopf trennt Entscheidung des Nutzers (mit Datum und Wortlaut) und Empfehlung des Planners; die drei Adressen existieren und nehmen die Sendungen in §1 und DoD an.

Keiner der gesichteten Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Zählung nach `AGENTS.md` §3.13 und §3.14** (Planner, 2026-10-22, im Commit, der diesen Plan schneidet und die Sendungen einträgt; nachgezählt am 2026-10-22 nach der Zuordnung der Stufen im Vertrag): **Liefer-Punkte 2** (Wächter mit Leser, Eintrag und Index-Zeilen | Gegenprobe). **Schichten 1** — der Harness (Hook, Leser, Fragment, Skript, `harness/README.md`, `.claude/settings.json`; Teilung nach `slice-harness-d-check-v0-85` §8: Gate-Konfiguration, Gate-Werkzeuge und Agenten-Anweisungen sind eine Schicht); die Entscheidung und der Vertrag liegen **vor** der Basis und gehören nicht in diesen Diff, ein Nachtrag macht zwei; die Planung ist keine Schicht. **Schätzung 476 bis 815, Mitte 645** bei **Marke 900** und Meldeschwelle 700 (§4), mit dem Faktor 3 bis 4,6 für Gegenprobe-Code, **78 Fälle** nach dem Vertrag — unter drei Liefer-Punkten und unter zwei Schichten, **an der 70 %-Linie** (72 %) und mit der Obergrenze unter der Marke. **Herkunft der Zählung:** Der ungeteilte Plan (angelegt am 2026-10-21) hatte 3 Liefer-Punkte (Wächter | Gegenprobe | Regel), 1 Schicht, die Schätzung des Planners 490 bis 900 mit dem Faktor 2; die Schätzung des Architect nach `SPEC-070` (800 bis 1400) und der Faktor 3 bis 4,6 hoben sie über die Marke; der Schnitt gab die Regel (Liefer-Punkt 3) an `slice-harness-rollenlauf-zeitgrenze-regel` und die zweite Stufe an `slice-harness-rollenlauf-zeitgrenze-haertung` ab; die Zuordnung des Architect hob Stufe 1 von 72 auf 78 Fälle (sieben `p7`-Fälle der Freigabe, zwei umbenannte Fälle, ohne den Bereich des Ankers) und die Schätzung von 440 bis 759 auf 476 bis 815. Die Sendungen an die drei Nehmer legen in diesem Slice keinen Liefer-Punkt und keine Schicht an.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der die Ortsangabe aus `slice-harness-gegenproben-umzug-dev` einträgt): Die Gegenproben liegen seit jenem Slice unter `tools/dev/<bereich>/`, die Pfade der Gegenproben in diesem Plan sind nachgezogen (`tools/dev/` ist wie `tools/harness/` Harness-Werkzeug, keine neue Schicht). Ortsangabe der Gegenprobe geändert, Liefer-Punkte (2) und Schichten (1) unverändert.
