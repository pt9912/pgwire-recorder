# Verifikation slice-harness-randform-stelle-und-lesart

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-15. **Spanne:** `88f9c5e` (Basis der Marke) bis `d53f830` (HEAD). **Gegenstand:** Slice `slice-harness-randform-stelle-und-lesart` (Plan unter `docs/plan/planning/in-progress/`) gegen DoD (§2), Plan §1 (Abgrenzung), §3, §4, §6, §7, §8, die Entscheidung des Nutzers vom 2026-10-15 („ja deine Vorschläge“; Punkt 1: Absatz zur Stelle einer Nachricht samt Reihenfolge mehrerer Ablehnungsgründe und Zeiger im Architect; Punkt 2: Satz zu „Lesart“ in `implement-slice.md`; kein Sensor), den Review-Report vom selben Tag (2 LOW behoben, F-3 bis F-6 INFO) und `AGENTS.md` §3.7 bis §3.14.

Eine Verifier-Vorlage führt die vendored Baseline nicht; Aufbau wie der Report von `slice-v1-abschluss-verbindungsende-format`. Der Slice liefert reine Prosa (drei Dateien, kein Code, kein Test); Belege sind deshalb Lesen des Diffs gegen §1 und §3, `git diff --numstat`, die Register-Dateien, der Stop-Hook und `make docs-check`. Nicht gefahren: voller `make gates` (Docker-Läufe nicht doppelt, der Hook trägt den Stand, siehe DoD-Punkt 4).

**Verdikt: PASS mit Befunden.** Beide Liefer-Punkte sind bestätigt, der Wortlaut trägt die DoD ohne Mutation, er geht nicht über den Beschluss hinaus und lockert keine Zusage. 0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO. Offen sind die Closure-Pflichten (Notiz, Lerneintrag, Register-Fortschreibung, Risiko-Ausgänge, Paarungen) und die Entscheidung des Planners über die Flanke `p7-int-treiber` (unten).

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1: `AGENTS.md` §3.12 trägt den Absatz *Die Stelle einer Nachricht gehört dazu* (Aufzählung der Stellen, Entscheidung je Stelle, Spezifikation trägt sie, Reihenfolge der Ablehnungsgründe), Falsch/Richtig (Anlass F-1 in `slice-v1-abschluss-verbindungsende-format`), Begründung mit `BEO-REPO/spec-randform-erst-im-review-entschieden`, Retirement-Check mit `seit slice-harness-randform-stelle-und-lesart`; `architect.md` nennt die Stelle mit Zeiger, ohne den Absatz zu wiederholen | bestätigt | `git diff 88f9c5e..HEAD -- AGENTS.md` gelesen: Absatz steht nach *Aufrufarten des Treibers*, vor der Begründung (`AGENTS.md` Zeile 339 ff.); fünf Stellen (Anfang, Mitte, direkt vor dem Ende, letzte Gruppe, nicht letzte Gruppe), „je Stelle eine Entscheidung; die Spezifikation trägt sie“, Satz zur Reihenfolge der Gründe („und die Spezifikation bindet sie“), Falsch/Richtig mit F-1. Retirement-Check um „für die Stelle einer Nachricht seit slice-harness-randform-stelle-und-lesart“ verlängert. `architect.md`: zwei Zeilen nach „Fehlt in §6 eine, … meldest du sie dem Planner“, mit Zeiger `AGENTS.md` §3.12, ohne Wiederholung der Aufzählung. Die `BEO-REPO/spec-randform-erst-im-review-entschieden`-Kennung steht in der bestehenden Begründung unverändert; der Absatz selbst führt sie nicht, der Plan verlangt sie „in der Begründung“, die bleibt dieselbe Zeile (siehe V-2). |
| 2: `implement-slice.md` (Randform-Rückgabe) trägt den Satz zu „Lesart“, „Auslegung“, „§6 nennt nicht“ in §7, Zeiger `AGENTS.md` §3.12, Anker `seit …`; kein Sensor | bestätigt | Drei Zeilen am Ende des Absatzes *Auch eine Lesart ist eine Entscheidung* (Zeile 147 ff.): „… hast du eine Randform gefunden und hältst an; ein Vermerk in §7 ersetzt die Rückgabe nicht (`AGENTS.md` §3.12, seit slice-harness-randform-stelle-und-lesart, `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`)“. Kein Skript, kein Make-Target, kein Hook im Diff (`git diff --name-only` außerhalb `docs/`: genau `AGENTS.md`, `architect.md`, `implement-slice.md`). |
| 3: Kein Wortlaut lockert oder ersetzt eine Zusage; Beleg „`git diff` ohne Löschzeile in den bestehenden Regeln, in §7“ | bestätigt in der Sache, Beleg in der Form abweichend (V-1, LOW) | `git diff` der drei Dateien hat genau eine Löschzeile: das Ende des Retirement-Checks in `AGENTS.md` (`slice-harness-gate-kette-schnell?` wird `slice-harness-gate-kette-schnell, für die Stelle einer Nachricht seit slice-harness-randform-stelle-und-lesart?`). Die Zeile wird nur verlängert, keine Bedingung wird enger, keine Zusage entfällt. §7 nennt die Zeile und begründet sie, trägt aber nicht den in der DoD verlangten Beleg „ohne Löschzeile“; die DoD-Formulierung war zu eng gefasst. |
| 4: `make gates` grün | durch den Hook getragen, nicht selbst gefahren | §7: erster voller Lauf Exit 2 (siehe Flanke), zweiter Exit 0; danach ändern nur `docs/plan/` und `docs/reviews/` den Baum (`git diff --stat 88f9c5e..HEAD` außerhalb `docs`: `AGENTS.md`, `.claude/…` ×2, alle in `95a5bc2`). `bash .claude/hooks/stop-require-gates.sh` mit `{}` gibt `{"decision":"approve"}`, Exit 0. Der Slice berührt weder Gate noch Gegenprobe; ob der volle Lauf „aus dem richtigen Grund“ grün war, ist für Prosa ohne Gate-Änderung ohne Aussage. |
| 5: Review durchgeführt, Report liegt vor | bestätigt | Report `2026-10-15-review-slice-harness-randform-stelle-und-lesart.md` (0 HIGH, 0 MEDIUM, 2 LOW, 4 INFO). F-1 und F-2 in `d53f830` am Plan behoben (nachgeprüft unten), F-3 durch §7 (Läufe stehen jetzt dort), F-4 siehe Flanke, F-5/F-6 ohne Aktion. Kein Self-Review. |
| 6 bis 9: Closure-Notiz mit Lerneintrag, Register `geplant` → `verkörpert`, Risiko-Ausgänge, Paarungen | offen (erwartet) | §7 trägt Belege des Implementers, noch keine Closure-Notiz; §2 ohne Häkchen. Handlung des Planners; Vorschläge unten. |

## Trägt der Wortlaut die DoD ohne Mutation (§3.10 für Prosa)?

Für Prosa gibt es keinen Test, der rot wird; die Prüfung ist, ob jede Zusage des Beschlusses als eigenes Satzglied im Wortlaut steht, sodass das Streichen eines Gliedes im Lesen auffällt. Gegengelesen Glied für Glied (Streichung eines Gliedes gedacht, Lücke benannt):

| Zusage des Beschlusses | Satzglied im Wortlaut | Lücke bei Streichung |
|---|---|---|
| Vertrag mit Nachricht, die eine Folge beendet oder ersetzt | „(Ende, `ready_for_query`, `sync`)“ | Auslöser verschwände; Regel fiele auf jeden Vertrag |
| §6 nennt je Nachrichtenart die zulässigen Stellen | „je Nachrichtenart die zulässigen Stellen in der Folge“ | Tabelle Nachricht × Stelle entfiele |
| fünf Stellen | Klammer „Anfang, Mitte, direkt vor dem Ende, letzte Gruppe, nicht letzte Gruppe“ | je Element eine Stelle; der Falsch-Fall (F-1) fiele durch „Mitte“ und „direkt vor dem Ende“ |
| Entscheidung je Stelle, Spezifikation trägt sie | „und je Stelle eine Entscheidung; die Spezifikation trägt sie“ | Entscheidung bliebe im Code (der Fehler aus F-1) |
| Reihenfolge mehrerer Ablehnungsgründe, Spezifikation bindet sie | zweiter Satz des Absatzes | Anlass `slice-v1-abschluss-protokollrand`, V-1 (`_pq_.` vor `PGR-E6003`) bliebe offen |
| Zeiger im Architect | Satz in `architect.md` mit `AGENTS.md` §3.12 | Prüfung vor dem Code ohne Zeiger |
| „Lesart“/„Auslegung“/„§6 nennt nicht“ in §7 = Randform, anhalten | Satz in `implement-slice.md` | Vermerk in §7 bliebe zulässiger Ersatz der Rückgabe (Anlass F-1 desselben Slice: Implementer nannte die Lesarten in §7, gab sie nicht zurück) |

Ein Satz ohne Prüfung bleibt: dass der Implementer den `implement-slice.md`-Satz befolgt, prüft nichts (Entscheidung des Nutzers: kein Sensor, Plan §1 und §6). Das ist zugesagt als Regel, nicht als Gate; der Wortlaut sagt nichts zu, was ein Test prüfen müsste (`AGENTS.md` §3.11): „hältst an“ ist Anweisung an den Schreiber. Ergebnis: trägt.

**Wortlaut über den Beschluss hinaus (Zusatz) oder darunter (Lockerung)?** Nein. Der Absatz enthält genau die Glieder des Beschlusses plus das Falsch/Richtig-Paar, das die DoD verlangt und das Teil des bestätigten Vorschlags war (Plan §3 *Gelieferter Wortlaut*). Die sprachlichen Abweichungen vom Entwurf (Plural „Nachrichten“, „und die Spezifikation bindet sie“, im Architect kein „fehlt sie, meldest du sie dem Planner“, im Command „als eigene Feststellung“) sind in §7 vollständig genannt und ändern keine Aussage; der gestrichene Halbsatz im Architect steht im Satz davor. „Als eigene Feststellung“ ist eine Einengung gegenüber dem bloßen Beschlusswort („Lesart“ in §7), aber aus §6 des Plans (Zitat ist keine Randform) begründet und im Review (Negativbefund) als Umsetzung bewertet; sie schwächt die Regel nicht, weil ein Zitat nie eine Feststellung des Implementers war.

## Plan gegen Diff

| Prüfung | Ergebnis |
|---|---|
| §1 Abgrenzung | Diff außerhalb `docs/`: nur `AGENTS.md`, `.claude/agents/architect.md`, `.claude/commands/implement-slice.md`. `plan-welle.md`, `planner.md`, `implementer.md`, `.harness/skills/` unberührt; kein Sensor, kein Produkt-Code, kein Gate, kein Handbuch, §3.14 unberührt (Hunk in §3.12). Eingehalten. |
| §3 (Plan) | Nach `d53f830` heißt der Block „Gelieferter Wortlaut“ und gibt die drei Fassungen wörtlich wie im Diff wieder (verglichen mit `git diff 95a5bc2`). F-2 des Reviews damit behoben. Restbestand „Entwurf“ in §6 (Risiko 4, „drei Dateien, Entwurf mit Falsch/Richtig“) und §1/§7 als Gegenüberstellung: harmlos (V-3). |
| §6 Randformen | F-1 des Reviews (Randform „nicht berührt“): §6 trägt jetzt „der Plan entscheidet es nicht über den Wortlaut hinaus (offene Frage an den Orchestrator, Review F-5)“, also keine Entscheidung ohne Träger. Der Wortlaut in `AGENTS.md` sagt weiter nichts dazu; das ist die offene Frage (V-2). |
| §7 gegen Diff | Zahlen stimmen (14 + 1, 2 + 0, 3 + 0, Summe 20), die eine Löschzeile ist benannt, die Läufe stehen jetzt in §7 (F-3 behoben). |
| Größe (Metrik §4) | `git diff --numstat 88f9c5e..HEAD -- AGENTS.md .claude`: 14/1, 2/0, 3/0, Summe **20**, selbst gerechnet. Unter Meldeschwelle 80 und Marke 100; unter der Schätzung von 30 bis 50 (die Schätzung war ungenau nach oben, im Sinn von `AGENTS.md` §3.14 kein Verstoß). Zwei Liefer-Punkte, zwei Schichten (Briefing, Rollen-/Command-Artefakte). |
| `AGENTS.md` §3.7 | indikativ, Herkunft als `seit …`-Anker, keine Chronik. Ohne Befund. |
| `AGENTS.md` §3.8, §3.5 | keine ADR berührt. |
| `AGENTS.md` §3.9 | Plan §1, §3, §6 folgen dem Diff nach `d53f830`; Kopf unverändert richtig (`Berührte Spec-Stellen: —`). Ohne Befund. |
| `AGENTS.md` §3.13 | Register-Zuweisung beider Beobachtungen: beide `state.md` tragen „**Stand:** geplant — Kennung: `slice-harness-randform-stelle-und-lesart`“. Der Nehmer (dieser Slice) liegt in `in-progress/`, schließt die Sendung in §1 nicht aus, die DoD führt sie (Register-Punkt). Die Adresse nimmt an. Keine Sendung an einen anderen Nehmer, daher keine Nachzählung bei Dritten; §8 nennt das. |
| `AGENTS.md` §3.3 | Move `next/` → `in-progress/` in eigenem Commit `1d3e6bf`; Inhalt danach. Ohne Befund. |
| Verweise | `make docs-check` siehe *Gates und Zustand*. |

## Findings

### V-1 LOW

- `quelle`: Plan §2 DoD-Punkt 3, `AGENTS.md` §3.9
- `pfad`: `docs/plan/planning/in-progress/slice-harness-randform-stelle-und-lesart.md` · §2 „Beleg: `git diff` der drei Dateien ohne Löschzeile in den bestehenden Regeln, in §7“
- `befund`: Der Diff hat eine Löschzeile (Ende des Retirement-Checks in `AGENTS.md`). Der in der DoD genannte Beleg („ohne Löschzeile“) ist wörtlich nicht erbringbar. In der Sache ist die DoD erfüllt (die Zeile ist eine Verlängerung, jede Bedingung bleibt), und §7 nennt die Zeile. Die Formulierung der DoD war falsch gefasst, weil jede Ergänzung des Retirement-Checks die Zeile umbricht.
- `verifizierbar`: ja (`git diff 88f9c5e..HEAD -- AGENTS.md .claude | grep '^-[^-]'`)
- `klasse`: DoD-Beleg enger als der Vorgang
- Vorschlag: In der Closure-Notiz die Belegform auf „keine Löschzeile außer dem verlängerten Ende des Retirement-Checks; keine Zusage entfällt“ stellen; die DoD-Zeile nicht still ändern, sondern in §7 benennen.

### V-2 INFO

- `quelle`: `AGENTS.md` §3.12, Plan §6
- `pfad`: `AGENTS.md` · „nennt §6 je Nachrichtenart die zulässigen Stellen in der Folge“
- `befund`: Der Wortlaut sagt nicht, ob eine nicht berührte Stelle in §6 verneint werden muss („nicht berührt“) oder wegfällt (Review F-1/F-5, im Plan als offene Frage an den Orchestrator geführt). Ebenso gibt es kein Falsch/Richtig-Paar für die Reihenfolge der Gründe, und die Begründungs-Zeile des Absatzes benennt `BEO-REPO/spec-randform-erst-im-review-entschieden` nicht neu (sie steht in der bestehenden Begründung, die für alle Absätze gilt). Beides liegt im Rahmen des Beschlusses; ein Zusatz wäre eine Änderung des Wortlauts und braucht die Antwort des Nutzers.
- `verifizierbar`: nein (Urteil)
- `klasse`: Regelabsatz mit offener Frage, vom Plan geführt
- Vorschlag: In der Closure als *weiter offen* (Risiko 1, Register `BEO-REPO/spec-randform-erst-im-review-entschieden`) führen; tritt die Frage in einem Vertrag auf, legt der Planner dem Nutzer einen Satz vor.

### V-3 INFO

- `quelle`: `AGENTS.md` §3.9
- `pfad`: `docs/plan/planning/in-progress/slice-harness-randform-stelle-und-lesart.md` · §6 Risiko 4 („Entwurf mit Falsch/Richtig“)
- `befund`: Ein Wort „Entwurf“ steht noch in §6, obwohl §3 den gelieferten Wortlaut führt. Der Satz ist als Risiko-Beschreibung weiter lesbar (er meint den Planungsstand vor der Lieferung); kein falscher Sachverhalt, aber der letzte Rest der alten Lesart.
- `verifizierbar`: ja (`grep -n Entwurf`)
- `klasse`: Plan folgt dem Diff nicht ganz (`BEO-REPO/plan-folgt-korrektur-nicht`, kleinste Form)
- Vorschlag: In der Closure-Notiz mit dem Ausgang von Risiko 4 überschreiben („gemessen 20, unter Marke“); keine eigene Aktion.

### V-4 INFO

- `quelle`: `harness/README.md` §Traceability rules, Register
- `pfad`: `docs/plan/planning/observations/BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig/state.md`
- `befund`: Das Register steht bei „2, offen, ohne Adresse“. Der erste volle `make gates` dieses Slice färbte erneut `p7-int-treiber` und zusätzlich `p7-sigTERM-skript` (Beleg: §7 des Plans, „Läufe“; Review F-4; die Flanke steht im Review-Report beschrieben). Der Slice berührt weder Gate noch Gegenprobe. Das ist das dritte Auftreten; siehe *Flanke* unten. Ich habe die Läufe nicht wiederholt (Docker-Läufe nicht parallel und nicht doppelt); die Aussage ruht auf §7 und dem Review.
- `verifizierbar`: ja (`make gates` wiederholt; nur ein Rot ohne Änderung an Gate oder Gegenprobe bestätigt es)
- `klasse`: Flanke der Gegenprobe färbt zusätzlichen Fall (Register `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig`)
- Vorschlag: siehe unten.

## Risiko-Ausgänge (Vorschläge für die Closure)

| Risiko (§6) | Vorschlag | Begründung |
|---|---|---|
| Aufzählung der Stellen zu eng oder zu weit | *weiter offen* → `BEO-REPO/spec-randform-erst-im-review-entschieden` | Die Aufzählung ist der Beschluss; ob ein Vertrag eine Stelle hat, die sie nicht nennt, zeigt erst ein Auftreten. Dazu die offene Frage zu „nicht berührt“ (V-2) in der Register-Zeile. |
| Satz in `implement-slice.md` gelesen, aber nicht befolgt (kein Sensor) | *weiter offen* → `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` | Entscheidung des Nutzers: kein Sensor; ein Sensor wird erst begründet, wenn das Register ein Auftreten trotz der Regel zählt (Plan §1). |
| Auftreten nach dem Slice schwächt oder bestätigt die Regel (Urteil des Planners) | *entfallen*, mit Begründung | Bis zur Closure steht kein Auftreten nach `95a5bc2` im Register (beide `state.md` führen als letztes `slice-v1-abschluss-verbindungsende-format`, vor dem Slice); der Retirement-Check wird erst bei einem späteren Auftreten gestellt. |
| Slice wächst über Meldeschwelle oder Marke | *entfallen*, mit Begründung | Gemessen 20 in der Metrik von §4, unter Meldeschwelle 80 und Marke 100. |

## Lerneintrag, Kandidaten (Urteil des Planners)

1. **Geschärfte Regel (Form des Lerneintrags):** `AGENTS.md` §3.12 (Stelle einer Nachricht, Reihenfolge der Gründe) und `implement-slice.md` (Vermerk in §7 ersetzt die Rückgabe nicht), `liegt in`: `AGENTS.md` §3.12 und `.claude/commands/implement-slice.md` Randform-Rückgabe, Anker `seit slice-harness-randform-stelle-und-lesart`. Register: beide Beobachtungen von `geplant` auf `verkörpert` (Zielort und Anker), Zähler folgen den Dateien in `evidence/`, kein Auftreten durch die Lieferung.
2. **Sicherheit der Aussage:** Die Regel wirkt erst, wenn das Register nach `95a5bc2` ein Auftreten zählt oder ausbleibt; das ist die Beobachtung, nicht der Slice.
3. **DoD-Formulierung** (V-1): „ohne Löschzeile“ für einen Slice, der einen Retirement-Check verlängert, ist nicht erbringbar. Kandidat für die Planer-Praxis, nicht für eine Regel: Beleg als „Löschzeilen außer Verlängerungen, die nichts lockern“ fassen. Kein eigener Register-Eintrag, solange es ein Fall ist.
4. **Wiederkehrende Finding-Klasse aus dem Review:** `BEO-REPO/plan-folgt-korrektur-nicht` (F-2, Plan §3 folgt dem gelieferten Wortlaut nicht), ein Auftreten; Beleg unter `evidence/` des Eintrags.

## Flanke `p7-int-treiber`: was der Planner bei der Closure eintragen muss

Stand im Register (`BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig`): „Zähler 2, unter der Schwelle“, `evidence/` führt `slice-v1-abschluss-postgres-versionen-zusage.md` (`p7-sigTERM-gruppe`) und `slice-v1-abschluss-protokollrand.md` (`p7-sigTERM-skript`, V-8). Dieser Slice ist das **dritte Auftreten**: im ersten vollen `make gates` färbt `p7-int-treiber` zusätzlich `p7-sigTERM-skript`, Einzellauf und zweiter voller Lauf grün, der Diff berührt Gate und Gegenprobe nicht (Prosa unter `AGENTS.md` und `.claude/`). Zusammen mit der Beschreibung der Muster im Review-Report (zeitabhängige Kopplung zwischen dem Abbruch des Treibers per `SIGKILL` unter INT und einem nachlaufenden TERM-Fall der Gruppe `p7-sig`; Beschreibung, keine Ursache) ist die Schwelle **erreicht**: nach dem Wortlaut des Eintrags ist es keine Notiz mehr, sondern eine Lücke, und sie braucht einen eigenen Slice.

Der Planner trägt bei der Closure ein:

1. **Evidenz-Datei** `evidence/slice-harness-randform-stelle-und-lesart.md` unter dem Eintrag (Zähler folgt den Dateien): Auftreten im ersten vollen `make gates` des Slice, Mutation `p7-int-treiber`, zusätzlich rot `p7-sigTERM-skript`, Einzellauf und zweiter voller Lauf grün, Slice berührt Gate und Gegenprobe nicht; Beleg: §7 des Plans, Review F-4, dieser Report V-4. Die Beschreibung der Flanke und das Muster aus dem Review-Report nicht kopieren, sondern zitieren.
2. **`state.md`:** `Zähler 3`, Schwelle erreicht, **Stand** von „offen“ auf „geplant — Kennung: `slice-harness-gegenprobe-matrix-p7-flanke`“ (Kennung ist Vorschlag; sie muss dem tatsächlichen Slice entsprechen), nicht auf „verkörpert“, weil nichts verkörpert ist. Chronik nicht anhängen (`AGENTS.md` §3.7, Zustandsfelder): der Satz „Zähler 2, unter der Schwelle: …“ wird durch den Stand und den Beleg ersetzt, nicht ergänzt.
3. **Eigener Slice** unter `docs/plan/planning/open/` (Gegenstand: Ursache der Kopplung zwischen `p7-int-treiber` und den TERM-Fällen der Gruppe `p7-sig` in `tools/test/postgres-versionen-gegenprobe.sh` und `SPEC-053` Punkt 7 klären und eine Mutation/ein Fall, der nur die Gruppe `p7-sig` färbt, zusagen; mit §6 für Randformen, `AGENTS.md` §3.10 für die Zusage, und die Sichtung durch Slice-Planung). Die Adresse muss annehmen (§3.13): Der Nehmer führt die Sendung in §1 und DoD, mit der Kennung dieses Gebers (`slice-harness-randform-stelle-und-lesart`) und der beiden früheren; kein vorhandener Slice in `open/`/`next/` führt die Gegenprobe der Matrix (`slice-v1-abschluss-postgres-versionen-matrix` liegt in `done/`, ist also keine Adresse).
4. **Roadmap:** Einordnung des neuen Slice (Wellen-Zugehörigkeit oder wellenlos mit Reihenfolge); eine Zeile im Drift-Log nur für die Umplanung, keine Schließung.
5. **Kein Risiko-Ausgang dieses Slice:** Die Flanke ist keine Folge der Lieferung (Slice berührt Gate und Gegenprobe nicht) und steht deshalb nicht unter §6; sie gehört ins Register, nicht zu den vier Risiken. Die Closure-Notiz nennt sie nur als Beobachtung.
6. **Hook/Gate-Bewertung:** Der Stop-Hook gab nach dem zweiten Lauf `approve`; die Flanke blockiert die Closure nicht (grüner Stand, zweiter Lauf Exit 0), solange kein Carveout nötig ist. Wiederholt sich das Rot **vor** dem eigenen Slice im Lauf eines späteren Slice, bleibt es ein Auftreten im Zähler dieses Slice (4×), kein neuer Eintrag.

## Gates und Zustand

- Eigene Läufe: `git diff --numstat 88f9c5e..HEAD -- AGENTS.md .claude` (20); `git diff 88f9c5e..HEAD -- AGENTS.md .claude` gelesen, Löschzeilen gefiltert (eine); `bash .claude/hooks/stop-require-gates.sh` (`approve`, Exit 0); Register-Dateien beider Beobachtungen und der Flanke gelesen; `make docs-check` (Ergebnis in der Übergabe).
- Voller `make gates` nicht gefahren; der Stop-Hook trägt den Stand, der Diff nach dem Lauf berührt nur `docs/plan/` und `docs/reviews/`.
- Aufräumen: keine Container, keine Kopien, kein Netz; `git status` zeigt nur diesen Report (nicht committiert).
