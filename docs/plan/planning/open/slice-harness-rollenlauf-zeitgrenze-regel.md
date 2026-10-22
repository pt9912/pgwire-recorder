# Slice slice-harness-rollenlauf-zeitgrenze-regel: Die Regel der Stunde steht im Briefing und in den sechs Rollen: Ein Rollenlauf hält bei einer Stunde an und berichtet, die Stunde ist zweite Rückführungs-Marke, und `.claude/agents/` steht in der Liste von `AGENTS.md` §3.9; dazu der Wortlaut von `AGENTS.md` §3.1, der die Host-Ziele nennt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-harness-regeln, **hinter** `slice-harness-groessen-metrik` und **vor** `slice-harness-rollenlauf-zeitgrenze` (§4 und §5 der Welle). **Zuordnung und Stelle sind Empfehlung des Planners vom 2026-10-22**, nicht Entscheidung des Nutzers: Beide Slices schreiben an `AGENTS.md` bei §3.14 (der Metrik-Slice schreibt den Absatz neu, die Stunden-Marke steht daneben), und die Regel an den Leser gilt ab ihrer Lieferung, ohne auf den Wächter zu warten. Entschieden hat der Nutzer die **Sache**: die Grenze von einer Stunde je Rollenlauf (2026-10-21, wörtlich: „Wir haben rollen (siehe .claude/agents) und commands -- max. 1 Stunde“) und am 2026-10-22 („deine empfehlungen“, so vom Koordinator übermittelt), dass (a) `AGENTS.md` §3.9 `.claude/agents/` aufführt und (b) die Stunden-Marke in der Regel dieses Slice steht, nicht in `slice-harness-groessen-metrik`, und am selben Tag (wörtlich „2. ja“, so vom Koordinator übermittelt auf dessen Frage), (c) der Geltungsbereich: Der Wächter deckt nur Rollenläufe, und der Regeltext nennt Commands und Hauptsitzung als nicht erzwungen. **Mitgenommen am 2026-10-22** (Auftrag des Koordinators, §3.13): der Wortlaut von `AGENTS.md` §3.1 (Docker-only, Host-Ziele mit ihren Werkzeugen), dessen **Sache** der Nutzer am 2026-10-22 entschieden hat (vom Koordinator übermittelt) und für den der Architect einen Wortlaut-Entwurf schrieb; **Wortlaut und Träger sind nicht entschieden**, die Zuordnung zu diesem Slice ist Empfehlung des Planners (§1). **Geschnitten am 2026-10-22** aus `slice-harness-rollenlauf-zeitgrenze` (dessen früherer Liefer-Punkt 3); angelegt in `open/`, weil er nicht priorisiert ist.

**Bezug:** Kein Anspruch des Produkts (`LH-*`) — der Slice schreibt eine Regel über Agentenläufe. Bindung an das Register: `BEO-REPO/rollenlauf-ohne-zeitgrenze` (1×). Regelstellen: `AGENTS.md` §3.9, §6 und eine neue Regel hinter §3.14; die sechs Rollen-Typen unter `.claude/agents/`. Aktive ADRs: [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) (Proposed; die Stunde als Marke steht dort unter Punkt 8), [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (Metrik der Marke, neben der die Stunde steht), [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Zeitbudget der Gegenprobe).

**Berührte Spec-Stellen:** [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) Punkte 10 und 12 (gelesen: die pflichtigen Läufe in der Stunde und die Grenzen; die Regel zeigt dorthin und wiederholt sie nicht). Der Slice ändert keinen Vertrag.

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Grenze steht dort, wo die Rollen sie lesen: (1) `AGENTS.md` trägt als eigene harte Regel (Nummer: nächste freie, nach §3.14) die Grenze von einer Stunde für einen **Rollenlauf** (ein Subagent unter `.claude/agents/`): Der Lauf hält bei der Grenze an, committet, übergibt seinen Stand und berichtet Stand und Offenes; er bricht nicht ab und nimmt nichts zurück; der Koordinator entscheidet danach über einen neuen Lauf. Dazu die Stunde als **zweite Rückführungs-Marke** neben der Zeilenzahl (Entscheidung des Nutzers vom 2026-10-22): Ein Slice, der in einer Stunde nicht umsetzbar ist, ist zu groß, und das ist eine Aussage über den Schnitt, nicht über den Lauf. Dazu der **Satz zum Verifier** nach `AGENTS.md` §6: Der volle Lauf einer Gegenprobe, die [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md) Punkt 10 (c) als nicht aufgenommen führt, ist seine erste Handlung und läuft im Hintergrund; ist seine Stunde erreicht, ohne dass der Lauf endete, berichtet er ihn als Offenes. Dazu die Zeile `.claude/agents/` in der Liste der Harness-Artefakte von §3.9 (Entscheidung des Nutzers vom 2026-10-22). (2) Die **sechs** Rollen-Typen unter `.claude/agents/` nennen je in einem Satz, was bei der Grenze geschieht (anhalten, committen, übergeben, berichten, nichts zurücknehmen); `verifier.md` nennt zusätzlich die Handlung zum vollen Lauf. (3) `AGENTS.md` §3.1 trägt den Wortlaut, der die Platzhalter der Vorlage ersetzt und die Host-Ziele mit ihren Werkzeugen nennt (Entwurf des Architect, §3).

**Mitgenommen: `AGENTS.md` §3.1** (Eintrag des Planners im Commit der Mitnahme, 2026-10-22, `AGENTS.md` §3.13; Quelle: der Wortlaut-Entwurf des Architect, den die Koordination übermittelte). Die Adresse ist dieser Slice, weil er `AGENTS.md` ohnehin anfasst und den Weg über den Nutzer für den Wortlaut schon trägt; `slice-harness-devcontainer-einrichten-doku` nimmt ihn **nicht** an (sein §1 und §4 schließen eine zweite Schicht aus; dass er allein `docs/developer/devcontainer.md` ändert, ist der Zweck seines Schnitts). Der Wortlaut ist **Empfehlung des Architect**; ob er so in `AGENTS.md` steht, entscheidet der Nutzer, nicht der Implementer. Der Herkunfts-Anker (Planner): `seit slice-harness-rollenlauf-zeitgrenze-regel`. Die Frage des Architect zum letzten Satz der Ausnahme beantwortet §6, Randform 8.

**Die Regel ist eine Regel an den Leser und sagt nicht mehr zu.** Dass ein Wächter die Grenze für Rollenläufe erzwingt, sagt dieser Slice nicht: Den Wächter liefert `slice-harness-rollenlauf-zeitgrenze` (Vertrag [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md)), und eine Zusage ohne Fall ist der Befund, gegen den `AGENTS.md` §3.11 steht. Der Wortlaut zeigt auf `SPEC-070` als Rang-Zeiger und nennt Commands und Hauptsitzung: Die Grenze gilt für sie als Regel, wird aber nicht erzwungen (`SPEC-070` Punkt 12 (1); Entscheidung des Nutzers vom 2026-10-22, wörtlich „2. ja“).

**Entwurf des Wortlauts und Entscheidung des Nutzers.** Regeländerungen in `AGENTS.md` und `.claude/agents/` sind eine Entscheidung des Nutzers. Entschieden hat er die Sache und (a) bis (c) (Kopf); **den Wortlaut hat er nicht entschieden**. Der Implementer legt den Wortlaut aus §3 (oder seine Überarbeitung) dem Nutzer vor und committet erst nach der Bestätigung; bis dahin ist er eine Empfehlung. Was der Nutzer nicht bestätigt, nimmt der Implementer nicht auf und nennt es in §7.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der **Wächter** (Hook, Leser, Eintrag in `.claude/settings.json`, Gegenprobe) — Klasse 1: `slice-harness-rollenlauf-zeitgrenze` und `slice-harness-rollenlauf-zeitgrenze-haertung` (beide nennen die Regel als ihren Ausschluss und diesen Slice als Adresse). Grund: anderer Vorgang und andere Dateien; die Regel an den Leser ist die Stufe, die der Nutzer am 2026-10-21 beschrieb.
- Die **Wiederholungs-Grenze** für volle Gegenprobe-Läufe — Klasse 1: `slice-harness-rollenlauf-wiederholungsgrenze`. Grund: ein zweiter Mechanismus mit eigenen Zusagen; hier steht kein Satz dazu.
- Ein **Anker für die Hauptsitzung** und eine Grenze für **Commands** in ihr — Klasse 3 und Grenze von `SPEC-070` Punkt 12 (1); eine Adresse gibt es nicht. Die Regel nennt beide als gemeint und nicht erzwungen.
- Die **Beleg-Orte** in den Rollen-Anweisungen (der Beleg eines Nachweises als Datei unter `docs/reviews/`) — Klasse 1: `slice-harness-regeln-implementer-uebergabe` (Sendung von `slice-harness-devcontainer-einrichten`). Grund: dieselben Dateien `implementer.md` und `verifier.md`, aber ein anderer Satz mit anderer Herkunft; der später schließende rechnet den Wortlaut des ersten ein.
- Die **Rückführungs-Marken in §4 bestehender Pläne** — Bestand bleibt bewusst stehen: Die Stunde als Marke gilt für Pläne, die nach der Regel entstehen oder geschnitten werden; ein Plan wird beim nächsten Anfassen nachgezogen.
- Änderungen an `AGENTS.md` §3.14 selbst (Metrik, Schätzung, Meldeschwelle) — `slice-harness-groessen-metrik` und `slice-harness-schaetzung-planner-faktor`; die Stunden-Marke steht **neben** §3.14 in der neuen Regel, nicht in dessen Absatz.
- Produkt-Code, Gates, Hooks, Spezifikation, Benutzerhandbuch — Schicht-Abgrenzung: **kein Code**, kein Verhalten des gebauten Binary.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Briefing):** `AGENTS.md` trägt (a) die Regel der Stunde für den Rollenlauf mit Falsch/Richtig-Paar, Herkunfts-Anker `seit slice-harness-rollenlauf-zeitgrenze-regel` und Retirement-Check, mit Commands und Hauptsitzung als gemeint und nicht erzwungen und dem Zeiger auf `SPEC-070`; (b) die Stunde als zweite Rückführungs-Marke, mit dem Satz, dass sie eine Aussage über den Schnitt ist; (c) in §6 den Satz zum Verifier; (d) in §3.9 die Zeile `.claude/agents/` in der Liste der Harness-Artefakte. Der Wortlaut ist vom Nutzer bestätigt, **bevor** der Commit entsteht; Datum und Fundstelle stehen in §7.
- [ ] **Liefer-Punkt 2 (sechs Rollen):** `planner.md`, `architect.md`, `implementer.md`, `reviewer.md`, `verifier.md` und `validator.md` unter `.claude/agents/` nennen je einen Satz (Grenze, anhalten, committen, übergeben, berichten, nichts zurücknehmen) mit dem Zeiger auf die Regel in `AGENTS.md`, ohne sie zu wiederholen; `verifier.md` zusätzlich die Handlung des vollen Laufs im Hintergrund. `grep` nach der Kennung dieses Slice findet die Regel in `AGENTS.md` und je einen Satz in den sechs Dateien.
- [ ] **Liefer-Punkt 3 (`AGENTS.md` §3.1):** §3.1 trägt statt der Platzhalter der Vorlage den Wortlaut aus §3 (Entwurf des Architect; Docker-only, die Ausnahme der Host-Ziele mit ihren Werkzeugen, Falsch/Richtig, Begründung) mit dem Herkunfts-Anker `seit slice-harness-rollenlauf-zeitgrenze-regel`; jede Aussage über ein Ziel stimmt mit `harness/README.md` §Sensors und `make help` überein (`AGENTS.md` §3.11). Der Wortlaut ist vom Nutzer bestätigt, **bevor** der Commit entsteht; Datum und Fundstelle stehen in §7.
- [ ] Kein Wortlaut lockert oder ersetzt eine Zusage von `AGENTS.md` §3.9 oder §6; der Diff ergänzt (Beleg: `git diff` ohne Löschzeile in den Regeln selbst, in §7). Ausgenommen ist §3.1: Dort ersetzt der Wortlaut den Text der Vorlage (`<venv/SDK/Toolchain-Install>`, `<z.B. …>`), keine Zusage.
- [ ] `make docs-check` und `make kopf-check` grün.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `.claude/` ändert).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Herkunfts-Anker `seit slice-harness-rollenlauf-zeitgrenze-regel` in `AGENTS.md` (`liegt in`), soweit aufgenommen.
- [ ] Beobachtungs-Register (`../observations/`): ein weiterer Beleg in `BEO-REPO/rollenlauf-ohne-zeitgrenze` (eine Datei in `evidence/`), falls ein Rollenlauf trotz der Regel die Stunde überschritt; sonst „keine Beobachtung angefallen“ in §7; kein Zähler wird gesetzt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der nächsten Welle-Closure; die Folge-Slice-Paarung prüft Existenz **und** Annahme: `slice-harness-rollenlauf-zeitgrenze` nennt diesen Slice in seinem §1 und seiner DoD (`grep` nach der Kennung dieses Slice, `AGENTS.md` §3.13).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` (neue harte Regel hinter §3.14) | update (ergänzen) | die Regel der Stunde mit Paar Falsch/Richtig, Begründung, Retirement-Check; die Stunde als zweite Marke; Zeiger auf `SPEC-070` |
| `AGENTS.md` §6 | update (ergänzen) | der Satz zum Verifier nach dem Satz zu den Gegenproben (Schritt 6) |
| `AGENTS.md` §3.9 | update (ergänzen) | die Zeile `.claude/agents/` in der Liste der Harness-Artefakte (Entscheidung des Nutzers vom 2026-10-22) |
| `AGENTS.md` §3.1 | update (ersetzt den Text der Vorlage) | Docker-only mit der Ausnahme der Host-Ziele und ihrer Werkzeuge; Entwurf des Architect (Wortlaut unten) |
| `.claude/agents/planner.md`, `architect.md`, `implementer.md`, `reviewer.md`, `verifier.md`, `validator.md` | update | je ein Satz; `verifier.md` zusätzlich die Handlung zum vollen Lauf |
| `docs/plan/planning/observations/BEO-REPO/rollenlauf-ohne-zeitgrenze/` | update bei Closure | Beleg oder „keine Beobachtung“ (kein Liefer-Punkt) |

Tests und Gates: keine neuen. Der Slice liefert Prosa (`AGENTS.md` §3.10 verlangt Mutationen für Verträge von Code und Gates); die Wirkung zeigt das Register (ein Lauf über der Stunde trotz der Regel).

**Entwurf des Wortlauts** (Empfehlung des Planners, keine Entscheidung des Nutzers; Nutzerentscheidung über den Wortlaut vor der Umsetzung):

> *`AGENTS.md`, neue Regel hinter §3.14:* **Ein Rollenlauf hat höchstens eine Stunde.** Ein Lauf eines Subagenten (eine Rolle unter `.claude/agents/`) hält bei einer Stunde Wanduhr seit seinem Beginn an: Er committet den Stand, übergibt ihn und berichtet Stand und Offenes; er bricht nicht ab und nimmt keine Arbeit zurück. Der Koordinator entscheidet nach dem Bericht über einen neuen Lauf; eine Fortsetzung beginnt eine neue Stunde. Für ein Command in der Hauptsitzung gilt die Grenze ebenso, erzwungen wird sie dort nicht (`SPEC-070` Punkt 12 (1)). Die Stunde ist zugleich eine zweite Rückführungs-Marke neben der Zeilenmarke von §3.14: Ein Slice, der in einer Stunde nicht umsetzbar ist, ist zu groß; das ist eine Aussage über den Schnitt, nicht über den Lauf. **Falsch:** Der Implementer arbeitet vier Stunden, weil keine Zeilenmarke feuert, und Review und Verifikation sehen den fertigen Stand. **Richtig:** Nach einer Stunde committet er, übergibt und berichtet; der Planner schneidet, wenn der Rest nicht in einen Lauf passt. Dazu Begründung und Retirement-Check `seit slice-harness-rollenlauf-zeitgrenze-regel`.
>
> *§6, nach dem Satz zu den Gegenproben (Schritt 6):* Der Verifier startet den vollen Lauf einer Gegenprobe, die `SPEC-070` Punkt 10 (c) als nicht aufgenommen führt, als erste Handlung seines Laufs im Hintergrund; ist seine Stunde erreicht, ohne dass der Lauf endete, berichtet er ihn als Offenes (Befehl, Basis-Hash, Ablage der Ausgabe, Startzeit), nicht als grün und nicht als rot.
>
> *§3.9, Liste der Harness-Artefakte:* `.claude/agents/` ergänzen.
>
> *Sechs Rollen-Dateien, je ein Satz:* Bei einer Stunde Wanduhr seit dem Beginn deines Laufs hältst du an: committen, übergeben, Stand und Offenes berichten, nichts zurücknehmen (`AGENTS.md`, Regel zum Rollenlauf). *`verifier.md`, zusätzlich:* der Satz zum vollen Lauf aus §6.

> *`AGENTS.md` §3.1 — Entwurf des Architect, Empfehlung, vom Planner unverändert übernommen:*
>
> ```text
> ### 3.1 Docker-only
>
> Kein Toolchain-Install am Host (Compiler, Interpreter, Paketmanager). Was gebaut, getestet oder geprüft wird,
> läuft über `make` in Docker; der Host braucht dafür `bash`, GNU `make`, `git` und Docker. Das gilt für jedes
> Gate der Kette (`GATE_CHECKS`) und für `make gegenproben`.
>
> **Ausnahme, abschließend:** `make container-start`, `make container-einrichten`, `make homebrew-nachweis` und
> die Gegenproben der Liste `GEGENPROBEN_HOST` (heute `container-start-gegenprobe`, gefahren von
> `make gegenproben-host`) sind Werkzeuge des Hosts: Sie legen einen Container an oder richten ihn ein oder
> prüfen gegen einen echten Dienst und brauchen am Host mehr. `container-start` und die Gegenprobe brauchen
> zusätzlich `u-boot`, `devcontainer`, `node`, `ssh` und `ssh-keygen`, `container-einrichten` braucht `ssh`,
> `homebrew-nachweis` braucht Netz; je Ziel steht es in `harness/README.md` §Sensors. Keines steht in
> `GATE_CHECKS`, und kein Gate ruft eines von ihnen auf. Ein weiteres Host-Ziel braucht eine ADR, die es hier
> nennt. Programme, die der Bootstrap unter `.harness/state/bin/` ablegt, sind kein Install.
>
> **Falsch:** `pip install …` oder `go build` am Host; ein Gate in `GATE_CHECKS`, das `u-boot` oder `devcontainer` ruft.
> **Richtig:** `make test`, `make gates`; `make container-start`, wo die Tabelle es als Host-Ziel führt.
>
> **Begründung:** Toolchain-Reproduzierbarkeit und Supply-Chain-Defense; ein Gate mit einem Host-Werkzeug liefe im
> Container nicht und prüfte dort nichts.
> ```
>
> Mit Anker `seit slice-harness-rollenlauf-zeitgrenze-regel` (Planner). Zwei Befunde des Planners zum Entwurf stehen in §6, Randform 8 und 9; der Wortlaut bleibt unverändert, bis der Nutzer entscheidet.

**Reihenfolge:** `AGENTS.md` zuerst (die Quelle), dann die sechs Rollen-Dateien, alles in einem Commit nach der Bestätigung; der Plan des Slice folgt in demselben Commit (`AGENTS.md` §3.9).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. `slice-harness-groessen-metrik` liegt in `done/` (**Empfehlung des Planners** für die Stelle: Der Metrik-Slice schreibt `AGENTS.md` §3.14 neu, und die Regel entsteht unmittelbar dahinter; ein gleichzeitiger Lauf erzeugt einen Konflikt statt einer Regel). Kein Start-Trigger ist [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) in `Accepted` oder der Wächter: Die Regel steht als Regel an den Leser ohne ihn (§1). WIP-Limit 1. Der Implementer legt den Wortlaut dem Nutzer vor und committet erst nach der Bestätigung; bis dahin hält er an und gibt an den Aufrufer zurück.

**Größe.** Metrik nach [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (der Befehl steht in `AGENTS.md` §3.14 nach `slice-harness-groessen-metrik`; wiederholt dieser Plan ihn, gilt der von dort): Summe aus Einfügungen und Löschungen über die Pfade `AGENTS.md` und `.claude`, mit den Ausschlüssen des Befehls; Basis ist der Commit, der diese Datei nach `in-progress/` bewegt. **Marke 200, Meldeschwelle 160:** Erreicht die Summe 160, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet; die Entscheidung steht als datierte Zeile in §7. Schätzung des Planners (2026-10-22, in dieser Metrik): `AGENTS.md` (die Regel mit Paar, Begründung und Retirement-Check 25 bis 45, Marken-Satz und Satz in §6 und die Zeile in §3.9 zusammen 8 bis 15), sechs Rollen-Sätze 15 bis 30, der zusätzliche Satz in `verifier.md` 4 bis 8, dazu `AGENTS.md` §3.1 (der Wortlaut des Architect ersetzt den Text der Vorlage) 22 bis 35, zusammen **74 bis 133, das sind 37 bis 67 % der Marke**. Der Faktor Tests zu Produktcode entfällt: Es entsteht kein Code und keine Gegenprobe.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt über 200, oder die Regel verlangt eine siebte Datei oder einen Sensor — zurück zur Zerlegung (§1).
- `in-progress` → `open` (blockiert): Der Nutzer bestätigt keinen Wortlaut und hat keine Fassung, die er bestätigen würde — der Slice wartet auf die Entscheidung, ohne dass ein Wortlaut in `AGENTS.md` steht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, `make gates` grün, Bestätigung des Wortlauts durch den Nutzer in §7 mit Datum, `grep` nach der Kennung dieses Slice findet die Regel in `AGENTS.md` und je einen Satz in den sechs Rollen-Dateien, Closure-Notiz mit Lerneintrag geschrieben (geschärfte Regel: `AGENTS.md`, `liegt in`, soweit aufgenommen).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen der Regel** (`AGENTS.md` §3.12; der Slice liefert keinen Vertrag von Code, aber eine Regel, deren Ränder der Planner sonst im Einzelfall entscheidet). Entscheidungsort ist der Wortlaut, vor dem Commit bestätigt; die Zusagen des Wächters stehen in [`SPEC-070`](../../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md):

1. *Wer:* Der Rollenlauf ist der Lauf eines Subagenten (eine Rolle unter `.claude/agents/`); ein Command in der Hauptsitzung ist gemeint und nicht erzwungen (`SPEC-070` Punkt 12 (1)); der Nutzer hat das am 2026-10-22 bestätigt (wörtlich „2. ja“, so vom Koordinator übermittelt).
2. *Beginn:* ab dem Beginn des Laufs; setzt der Koordinator den Agenten fort, beginnt eine neue Stunde (`SPEC-070` Punkt 9). Der Wortlaut sagt es.
3. *Was „anhalten“ heißt:* committen, übergeben, berichten; kein Abbruch, keine Rücknahme von Arbeit.
4. *Übergabe je Rolle:* Der Vorschlag des Planners im Register und `SPEC-070` Punkt 6 sagen „in das Austausch-Repo pushen“; das ist der Weg des Implementers im Dev-Container. Planner und Architect pushen nach den Aufträgen des Koordinators nichts. Der Entwurf sagt deshalb „übergibt ihn“; ob der Wortlaut je Rolle den Weg nennt (Push, Commit auf `main`, Bericht), ist **Entscheidung des Nutzers beim Wortlaut.**
5. *Pflichtige Läufe:* Die Stunde trägt `make gates` und den vollen Lauf einer aufgenommenen Gegenprobe; nicht den vollen Lauf einer nicht aufgenommenen (der Satz zum Verifier) und nicht `make gegenproben` am Ende einer Welle (`SPEC-070` Punkt 10). Der Wortlaut zeigt auf `AGENTS.md` §6 und `SPEC-070`, statt sie zu wiederholen.
6. *Stunden-Marke:* „in einer Stunde nicht umsetzbar“ ist ein Urteil nach dem Bericht des Implementers; der Planner entscheidet über den Schnitt. Es gibt keinen Sensor (die Marke stützt sich auf den Bericht, [ADR-0058](../../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) Punkt 8).
7. *Aufruf-Voraussetzung, Aufrufarten des Treibers, Plattform, Stelle einer Nachricht:* nicht berührt (`AGENTS.md` §3.12); die Regel liefert keinen Vertrag, die Plattform des Wächters steht in `SPEC-070`.
8. *`AGENTS.md` §3.1, letzter Satz der Ausnahme* („Programme, die der Bootstrap unter `.harness/state/bin/` ablegt, sind kein Install“) — die Frage des Architect, ob der Satz gilt: **Empfehlung des Planners: bestätigen, mit einer Verengung des Wortlauts.** Gelesen: `.harness/state/bin/ai-harness-init` ist ein statisch gelinktes Go-Binary, gitignoriert, vom Bootstrap abgelegt, nicht vom Agenten installiert; es läuft am Host ohne Docker und wird vom Emitter `.claude/hooks/span-emit.sh`, von `make span-report` und von `make archive-welle` gerufen (`harness/mk/erfassung.mk`, `harness/mk/archivierung.mk`). Ohne den Satz wäre die Liste der Ausnahme, die „abschließend“ heißt, falsch: Diese Ziele laufen am Host und stehen dort nicht. Es ist weder Compiler noch Interpreter noch Paketmanager. Die Verengung: Der Satz nennt das Programm und seine Aufrufer („heute `ai-harness-init`, gerufen von …“), damit er nicht jedes Programm unter `.harness/state/bin/` deckt. **Entscheidung des Nutzers beim Wortlaut.**
9. *`AGENTS.md` §3.1, erster Satz des Entwurfs* („Was gebaut, getestet oder geprüft wird, läuft über `make` in Docker; der Host braucht dafür `bash`, GNU `make`, `git` und Docker“) — **Befund des Planners an den Nutzer:** Gates der Kette laufen ohne Docker am Host (`make kopf-check`, `make abdeckung-check`, `make meldungskatalog-check`: „bash und awk, ohne Docker“ in `harness/README.md` §Sensors), und die Gegenproben der Liste `GEGENPROBEN_CONTAINER` brauchen zusätzlich `awk`, `sed`, `grep`, `setsid`, `timeout`. Der Satz stimmt für Gates nicht ausnahmslos und nennt die Standardwerkzeuge nicht; Empfehlung: „läuft über `make`“ und die Standardwerkzeuge von POSIX und GNU nennen oder „in Docker“ auf Build und Test beschränken. Der Wortlaut bleibt Entwurf des Architect, bis der Nutzer entscheidet.

**Risiken:**

- Der Nutzer ändert den Wortlaut, streicht einen Teil oder bestätigt keinen — **Ausgang:** offen bis Closure (erwartet: entfallen, der bestätigte Wortlaut steht in `AGENTS.md`; bei Ablehnung geht der Slice zurück nach `open/`, §4).
- Die Regel ist Prosa, bis `slice-harness-rollenlauf-zeitgrenze` den Wächter liefert: Ein Lauf, der die Stunde nicht bemerkt, bleibt ungebremst — **Ausgang:** weiter offen → `BEO-REPO/rollenlauf-ohne-zeitgrenze` (ein Auftreten trotz der Regel zählt; der Wächter ist die Antwort, Adresse `slice-harness-rollenlauf-zeitgrenze`).
- Der Wortlaut verlangt von Planner und Architect einen Push, den ihre Aufträge ausschließen (Randform 4) — **Ausgang:** entfallen, wenn der bestätigte Wortlaut den Weg je Rolle nennt oder „übergibt“ sagt.
- Der Satz zum Verifier steht neben dem Satz zu den Gegenproben in `AGENTS.md` §6 (Schritt 6) und widerspricht ihm — **Ausgang:** offen bis Closure (erwartet: entfallen, der Satz ergänzt den Satz zum vollen Lauf des Verifiers und ersetzt ihn nicht).
- Die Änderungen treffen auf `slice-harness-groessen-metrik` und `slice-harness-schaetzung-planner-faktor` (`AGENTS.md` §3.14) und auf `slice-harness-regeln-implementer-uebergabe` (`implementer.md`, `verifier.md`) — **Ausgang:** offen bis Closure (erwartet: entfallen, andere Absätze und Sätze; der Slice, der später schließt, rechnet den Wortlaut des ersten ein).
- Der Wortlaut von §3.1 ist Entwurf des Architect, und seine zwei Befunde (Randformen 8 und 9) sind offen — **Ausgang:** offen bis Closure (erwartet: entfallen, der Nutzer bestätigt oder ändert den Wortlaut; ohne Bestätigung entfällt Liefer-Punkt 3 und §7 nennt „nicht aufgenommen“).
- Der Slice nimmt mit §3.1 eine fremde Regel auf (`BEO-REPO/slice-waechst-durch-uebernahmen`) — **Ausgang:** entfallen, weil die Zählung in §8 steht (3 Liefer-Punkte, 2 Schichten, 37 bis 67 % der Marke) und kein Ausschluss von §1 sie trifft; über der Marke geht der Punkt zurück an den Planner.
- Der Slice wächst über Marke oder Meldeschwelle von §4 (Schätzung 74 bis 133 gegen Marke 200) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Hier stehen: die Bestätigung des Wortlauts durch den Nutzer (Datum, Fundstelle, je Teil), die Größe in der Metrik von §4 (Basis, letzter Commit, Summe) mit dem Ergebnis der Schätzung gegen die Messung, der Beleg des Registers oder „keine Beobachtung angefallen“, und die Ausgänge der Risiken.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-22). Treffer in der Sub-Area `REPO`: `BEO-REPO/rollenlauf-ohne-zeitgrenze` (1×, unter der Schwelle) — der Auslöser; `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (14×, `verkörpert`) — die Stunde ist eine zusätzliche Grenze über den 600 s der Gegenprobe; `BEO-REPO/groessenprobe-ohne-messbaren-stand` (5×, `verkörpert`) — die Stunde als Marke stützt sich auf den Bericht, nicht auf einen Diff; `BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet` (1×) — Kopf und §1 trennen Entscheidung und Empfehlung. Eine Schwelle wird mit diesem Slice nicht neu erreicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-22, im Commit, der diesen Plan anlegt):** Liefer-Punkte 2 (Briefing: Regel der Stunde | sechs Rollen), **nachgezählt am 2026-10-22 mit der Mitnahme von §3.1: 3** (dazu `AGENTS.md` §3.1). Schichten 2, konservativ gezählt (die Mitnahme legt keine an: §3.1 steht in `AGENTS.md`): Briefing (`AGENTS.md`) und Agenten-Rollen (`.claude/agents/`); nach der Teilung in `slice-harness-d-check-v0-85` §8 wäre es eine (Agenten-Anweisungen). Planung zählt keine Schicht; kein Produkt-Code. Schätzung 52 bis 98, mit §3.1 **74 bis 133** in der Metrik von §4 (Marke 200, Meldeschwelle 160); der Faktor Tests zu Produktcode entfällt. Der Slice trägt keine Sendung in einen anderen Nehmer ein; `slice-harness-rollenlauf-zeitgrenze` nennt ihn als Adresse der Regel (dort §1 und DoD, mit Nachzählung in dessen §8).
