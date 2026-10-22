# Slice slice-harness-regeln-rueckbau: Die eigenen harten Regeln von `AGENTS.md` kosten weniger je Produkt-Slice: Raten-Kriterium statt Retirement-Check, eine Marke von 1500 Zeilen ohne Schätzung und Meldeschwelle, Review-Stufen mit Failure-Szenario, Geltungsbereich von §3.10 und §3.12

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Geprüft und nicht gewählt: `welle-harness-regeln` — sie ist pausiert (Entscheidung im Kopf), und dieser Slice ist der Rückbau, den ihre Pause verlangt, nicht ihr Gegenstand; ihr Closure-Trigger (alle Regeln eingetragen) widerspräche seinem Ziel.

**Bezug:** Kein Anspruch des Produkts (`LH-*`) — Harness-Arbeit; Anlass ist die Frage des Nutzers, warum am Produkt kaum Code entsteht (Entscheidung im Kopf).

**Berührte Spec-Stellen:** —

**Entscheidung:** Entscheidung des Nutzers vom 2026-10-22, wörtlich „ja“, so vom Koordinator übermittelt, auf den Bericht des Planners mit fünf Entscheidungen (der Koordinator hatte „ja für alle außer 5“ angeboten und zu 5 „nein, zurückstellen“ empfohlen). **Angenommen sind 1 bis 4:** (1) der Retirement-Check wird einmal ausgesetzt, künftig gilt das Raten-Kriterium (tritt eine Klasse in über der Hälfte der Slices auf: Sensor oder streichen); (2) `welle-harness-regeln` pausiert, weitere Regel-Slices zu `AGENTS.md` §3.10 bis §3.13 pausieren, zuerst dieser Rückbau-Slice, dann `slice-version-option` und `slice-v1-abschluss-verbindungsende-upstream-extended`; (3) Marke und Review-Stufen: die Marke ist eine Zahl mit dem Standard 1500, Faktor, 70-%-Schnitt, Meldeschwelle und Messzeile in §7 entfallen, und in `.harness/skills/reviewer.md` ist ein Verstoß gegen §3.10, §3.12 oder §3.13 ohne Failure-Szenario LOW, nicht „mindestens MEDIUM“; (4) der Geltungsbereich von §3.10 und §3.12 wird begrenzt (Record-/Replay-Pfad, Format, Protokollrand; „je Element“ nur für Sätze der Spezifikation; Randformen nur für eine neue Wire-, Datei- oder Formatgrenze, bei einer Option oder einem Hilfetext ein Verweis auf die Geschwister-Regel). **Nicht angenommen** ist 5 (`slice-harness-mutation` vorziehen); es entsteht **kein neuer Sensor** (Auftrag des Koordinators vom 2026-10-22, begründet mit dem Wunsch des Nutzers nach schnellerer Produktentwicklung statt neuer Harness-Maschinerie: der Vorschlag `make groesse` entfällt, die Messung bleibt der Befehl in §3.14). Den Start in `next/` mit `Verantwortlich: pt9912` hat der Nutzer entschieden (Koordinator, 2026-10-22); die Stelle **nach** `slice-harness-gegenproben-umzug-dev` und vor `slice-version-option` ist Auftrag des Koordinators (der Nutzer hat „zuerst dieser Slice, dann `slice-version-option`“ entschieden, Punkt 2). **Die sechs Fragen des Berichts des Planners** (§6, Randformen 1, 2, 3, 6 und 7, dazu die Pause der Wächter der Stunde in §5 der Welle) **hat der Nutzer am 2026-10-22 entschieden** (wörtlich „1. ja“, so vom Koordinator übermittelt: alle sechs Empfehlungen des Planners gelten). **Der Wortlaut in §3 ist Entwurf des Planners**; der Nutzer bestätigt den Diff vor dem Merge nach `main` (Muster „Kein Halt, Diff vor Merge“ wie bei `slice-harness-rollenlauf-zeitgrenze-regel`).

**Verantwortlich:** pt9912

**Autor:** Planner. **Datum:** 2026-10-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Regeln von `AGENTS.md` §3.9 bis §3.14 kosten einen kleinen Produkt-Slice weniger, ohne die Fehlerklassen zurückzuholen, die sie verhindern, und die Träger, die sie lesen, folgen. Drei Liefer-Punkte, je einer pro angenommener Entscheidung: (1) das **Raten-Kriterium** als neue Regel `AGENTS.md` §3.15 an die Stelle der Retirement-Check-Sätze in §3.9 bis §3.14; (2) die **Marke** in §3.14 (eine Zahl, Standard 1500; Faktor, 70-%-Schnitt, Meldeschwelle, Schätzung und Messzeile in §7 entfallen); (3) der **Geltungsbereich** von §3.10 und §3.12 mit den **Review-Stufen** im Reviewer-Skill. Jeder Punkt trägt `AGENTS.md` und die Anweisungen, die denselben Satz lesen (§3).

**Beleg des Nutzens (Messung des Planners, 2026-10-22).** Quellen: `git log`, `docs/reviews/`, `evidence/` des Registers, die Erfassung unter `.harness/state/spans` (gitignoriert, daher stehen die Zahlen hier). Gemessen sind: (a) fünf Produkt-Slices (`slice-v1-abschluss-cancel-ohne-schluessel`, `-verbindungsende-upstream`, `-record`, `-format`, `-hilfetext-binary`) haben einen Plan von 258 bis 343 Zeilen, davon §7 35 bis 60 %, 15 bis 30 Commits, davon 2 bis 6 am Code, im Mittel 2 h 15 vom Anspruch bis `done/`, 0 bis 3 Nacharbeitsrunden mit je einem vollen `make gates` (10 min 55 s und 13 min 8 s in zwei §7); (b) Rollenminuten der 16 Slices zum Verbindungsende und zu `CancelRequest` (Wanduhr zwischen erstem und letztem Ereignis eines Laufs, Spans): Implementer 58 %, Verifier 19, Architect 9, Reviewer 8, Planner 6, davon 30 bis 55 % der Implementer-Zeit Warten auf `make`; (c) seit dem 2026-10-09 nennt die Quelle von 67 der 73 MEDIUM und HIGH in 60 Produkt-Reviews eine Regel aus §3.9 bis §3.14 (nicht alle sind Formalie: ein grüner Mutant als äquivalent und eine Testlücke waren echt); der letzte HIGH mit Code-Defekt war am 2026-10-05; (d) die acht Pläne zum Verbindungsende des Servers vom 2026-10-12 sind fünfzehn Slices geworden, Schnitte nach der Marke; Reviews von 560 bis 941 Zeilen dauerten 5 bis 23 Minuten; (e) Rate je Klasse in den 73 produktseitig geschlossenen Slices: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` 59 (81 %), `BEO-REPO/plan-folgt-korrektur-nicht` 51 (70 %), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` 42 (58 %), `BEO-REPO/spec-randform-erst-im-review-entschieden` 36 (49 %), `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` 19 (26 %). **Geschätzt, nicht gemessen:** die Minuten je Regel. **Grenzen:** Das Register beginnt am 2026-10-04 und die Regeln §3.9 bis §3.12 stammen vom 04./05.10., eine Messung vor ihnen gibt es nicht; nur die Schichtzählung (§3.13, fünf Auftreten bis zum 12.10., danach keines) zeigt eine Wirkung; ein Auftreten im Register ist ein Fund im Review, kein Fehler im Produkt; ein fortgesetzter Agent zählt in den Spans einmal.

**Was der Rückbau an der Baseline ändert (am Text geprüft, `.harness/baseline/v6.18.0/regelwerk/`).** Er braucht **keinen MR**: (a) Der **Retirement-Check** ist Baseline (`grundlagen-traceability.md` §Sensor 2), aber er verlangt, dass die Herkunft einer Regel vor dem Entfernen oder Lockern **konsultiert und das Ergebnis dokumentiert** wird; er verbietet das Lockern nicht, wenn die Beobachtung wieder auftrat. Die Sätze in `AGENTS.md` lasen sich als Sperre. §3.15 behält die Konsultation (Zähler und Rate im Plan des auslösenden Slice) und ersetzt die Sperr-Lesart durch das Raten-Kriterium; dieser Slice dokumentiert seine eigene Konsultation hier und in §7. Wollte der Nutzer auch die Konsultation fallen lassen, wäre das eine Abweichung von dieser Baseline-Regel und ein MR (`Ersetzt-Baseline-Regel: grundlagen-traceability.md` §Sensor 2); dieser Slice liefert das nicht. (b) **Marke, Faktor, Schätzung, Meldeschwelle** sind eigene Regel; die Baseline verlangt nur, dass §4 die Bedingung der Rückführung nennt und dass ein Slice höchstens drei Liefer-Punkte und zwei Schichten hat (`modul-05-planning-harness.md` §Trigger je Lifecycle-Übergang, §Ziel-Form: Slice). (c) **Review-Stufen:** `modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill sagt „Kein HIGH- oder MEDIUM-Finding ohne Failure-Szenario“; die Änderung stellt den Skill dieser Regel wieder gleich. (d) **Geltungsbereich §3.10 und §3.12:** Randformen-Liste und Mutation je Zusage stehen nirgends in der Baseline; nur `modul-11-verification.md` §Bewusstes Brechen verlangt vom Verifier den Rot-Beleg für korrektheitskritische DoD-Punkte, und der Record-/Replay-Pfad bleibt im Geltungsbereich.

**Bestand bleibt bewusst stehen** (Auftrag des Koordinators vom 2026-10-22: die neue Regel gilt für Pläne, die nach ihr entstehen): 23 Pläne in `open/`, `next/` und `in-progress/` führen Marke, Meldeschwelle und Schätzung in §4; sie behalten sie. Die Regel in §3.14 sagt es. Für die zwei nächsten Produkt-Pläne (`slice-version-option`, `slice-v1-abschluss-verbindungsende-upstream-extended`) ist offen, ob der Planner §4 beim Anspruch nachzieht (§6, Randform 7).

**Beim Anlegen erledigt (Planner, kein Liefer-Punkt):** `welle-harness-regeln` ist pausiert (Pause und Wiederaufnahme in §5 der Welle), das Drift-Log der Roadmap trägt die Umplanung, die zwei Produkt-Pläne nennen die Reihenfolge in ihrem Start, und `slice-harness-rollenlauf-zeitgrenze-regel` trägt die Bindung an die Nummer §3.15 und den Wortlaut ohne Retirement-Check (Sendung nach `AGENTS.md` §3.13).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Geltungsbereich der Dev-Seite** (`tools/dev/`, Positivliste des Stop-Hooks, der MR zur Dev-Seite, die Rollen-Sequenz dort) — Klasse 3, anderer Vorgang: Er gehört zur ADR über die Trennung von Entwicklung und Nutzen der Harness-Werkzeuge (Entwurf auf dem Zweig `adr-dev-nutzen`, auf `main` noch ohne Kennung) und kommt als eigener Slice **hinter** diesem; er fasst `AGENTS.md` ebenfalls an, der spätere rechnet den Wortlaut des früheren ein. Dieser Slice hat keinen Eintrag für ihn und legt ihn nicht an.
- **Gates und Werkzeuge unter `harness/vertraege/`** — Klasse 2, Bestand: §3.10 und §3.12 gelten für sie unverändert (der Wortlaut sagt es); der Nutzer hat den Geltungsbereich für das Produkt entschieden (Record-/Replay-Pfad, Format, Protokollrand). Offen: Randform 3.
- **Ein Sensor jeder Art** (`make groesse`, das Mutations-Gate `slice-harness-mutation`, ein Gate für die Paarung der Folge-Slices) — Klasse 3: Entscheidung des Nutzers vom 2026-10-22 (kein neuer Sensor; 5 zurückgestellt). `slice-harness-mutation` bleibt, was er ist.
- **Inhalt von §3.8, §3.9, §3.11 und §3.13** — Klasse 2, Bestand: Es entfallen allein die Retirement-Check-Sätze; kein Absatz wird verschlankt (der Nutzer hat 1, 3 und 4 angenommen, nicht die Vorschläge zu §3.9, §3.11 und §3.13 aus dem Bericht). Der Vorschlag, den Volltext der Regeln aus `AGENTS.md` auszulagern, ist nicht angenommen.
- **Die Pläne des Bestands** (Marke, Meldeschwelle, Schätzung in §4 von 23 Plänen) — Klasse 2, Bestand (oben).
- **Die Slices der pausierten Welle** und `slice-harness-regeln-gate-konfiguration` — Klasse 1/2: Sie warten (§5 der Welle, Drift-Log); beim Wiederaufnehmen rechnet jeder den neuen Wortlaut ein. Dieser Slice ändert keinen von ihnen außer der Bindung in `slice-harness-rollenlauf-zeitgrenze-regel`.
- **Produkt-Code, Gates, Hooks, Spezifikation, Benutzerhandbuch, Register-Zähler** — Schicht-Abgrenzung: **kein Code**, kein Verhalten des gebauten Binary. Die Klassen im Register, deren Regel entfällt, bekommen ihren Ausgang bei der Closure (Planner), nicht als Liefer-Punkt.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Raten-Kriterium):** `AGENTS.md` trägt die Regel §3.15 mit Falsch/Richtig-Paar, Begründung und Herkunfts-Anker `seit slice-harness-regeln-rueckbau`; die 13 Retirement-Check-Stellen in §3.9 bis §3.13 sind entfallen, ohne dass ein Herkunfts-Anker `seit …` verloren geht (Beleg: `grep -n 'Retirement-Check' AGENTS.md` findet nur Zeilen in §3.15; `git diff` zeigt in §3.9 bis §3.13 nur Löschzeilen der Stellen aus §3). `.claude/commands/close-welle.md` (Schritt 3, Lese-Schritt) und `.claude/commands/implement-slice.md` (Schritt 25, wellenlos) nennen die Rate je fortgeschriebener Klasse mit dem Zeiger auf §3.15.
- [ ] **Liefer-Punkt 2 (Marke):** `AGENTS.md` §3.14 trägt die Marke 1500 als eine Zahl, den Messbefehl unverändert, und keinen Faktor, keinen 70-%-Schnitt, keine Meldeschwelle, keine Schätzung des Planners oder des Architect und keine Messzeile in §7; das Falsch/Richtig-Paar zum Faktor ist durch ein Paar zur Marke ersetzt. `.claude/commands/plan-welle.md` Schritt 6 nennt die Marke statt Schätzung, Faktor und Meldeschwelle.
- [ ] **Liefer-Punkt 3 (Geltungsbereich und Review-Stufen):** `AGENTS.md` §3.10 und §3.12 tragen den Geltungsbereich aus §3 (Record-/Replay-Pfad, Format, Protokollrand; „je Element“ nur für Sätze der Spezifikation; Randformen nur für eine neue Wire-, Datei- oder Formatgrenze, bei Option und Hilfetext der Verweis auf die Geschwister-Regel; Gates und Werkzeuge ohne Einschränkung); `.claude/commands/implement-slice.md` (Schritte 12, 13, 19 und Randform-Rückgabe) und `.claude/agents/architect.md` folgen; `.harness/skills/reviewer.md` stuft einen Verstoß gegen §3.10, §3.12 oder §3.13 ohne Failure-Szenario als LOW (die vier Einträge unter MEDIUM nennen „MEDIUM nur mit Failure-Szenario, sonst LOW“; die zwei Kommentar-Einträge unter HIGH tragen denselben Zusatz, Randform 6).
- [ ] Der Wortlaut folgt §3 (Treue zum Entwurf; jede Abweichung steht mit Grund in §7); der Nutzer bestätigt den Diff, **bevor** er nach `main` kommt, nicht vor dem Commit auf dem Zweig des Implementers; Datum und Fundstelle der Bestätigung stehen bei der Closure in §7.
- [ ] Kein Wortlaut außerhalb der Stellen in §3 ändert eine Zusage von `AGENTS.md` (Beleg: `git diff` von `AGENTS.md` berührt nur die Stellen aus §3; in §7).
- [ ] `make docs-check` und `make kopf-check` grün.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `.claude/` ändert).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Herkunfts-Anker `seit slice-harness-regeln-rueckbau` in `AGENTS.md` (`liegt in`).
- [ ] Beobachtungs-Register (`../observations/`): die Klassen, deren Regel mit diesem Slice entfällt (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, `BEO-REPO/meldeschwelle-ohne-artefakt`, `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`, `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`), bekommen in ihrer `state.md` den Ausgang *gestrichen* mit Begründung; die Konsultation der Herkunft (Zähler und Rate) für die gelockerten Regeln steht in §7; kein Zähler wird gesetzt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit); die Folge-Slice-Paarung prüft Existenz **und** Annahme (`AGENTS.md` §3.13): `slice-harness-rollenlauf-zeitgrenze-regel` nennt diesen Slice in seinem §1 (`grep` nach der Kennung dieses Slice).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3.15 (neu, hinter §3.14) | update (ergänzen) | das Raten-Kriterium mit Falsch/Richtig, Begründung, Konsultation der Herkunft (Liefer-Punkt 1) |
| `AGENTS.md` §3.9 bis §3.13 | update (löschen) | die 13 Retirement-Check-Stellen (Tabelle unten); sonst unverändert (Liefer-Punkt 1) |
| `AGENTS.md` §3.14 | update (ersetzen) | Marke, Messbefehl, Falsch/Richtig, Begründung (Wortlaut unten; Liefer-Punkt 2) |
| `AGENTS.md` §3.10, §3.12 | update (ersetzen, je ein Absatz und ein Halbsatz) | Geltungsbereich (Liefer-Punkt 3) |
| `.harness/skills/reviewer.md` | update | Stufen der vier Einträge unter MEDIUM, ein Satz in der Klassifikation, ein Eintrag unter LOW (Liefer-Punkt 3); `.claude/agents/reviewer.md` zeigt auf den Skill und bleibt |
| `.claude/commands/implement-slice.md` | update | Schritte 12, 13, 19 und Randform-Rückgabe (Geltungsbereich), Schritt 25 (Rate) |
| `.claude/commands/plan-welle.md` Schritt 6 | update | Marke statt Schätzung, Faktor, Meldeschwelle; Geltungsbereich der Randformen |
| `.claude/agents/architect.md` | update | Geltungsbereich der Randformen, Prüfung des Verweises |
| `.claude/commands/close-welle.md` Schritt 3 | update | die Rate im Lese-Schritt |
| `docs/plan/planning/observations/BEO-REPO/…/state.md` | update bei Closure | Ausgang *gestrichen* für die vier Klassen der entfallenden Regel (kein Liefer-Punkt) |

Tests und Gates: keine neuen, kein Sensor (§1). Der Slice liefert Prosa; die Wirkung zeigt das Register (die Rate je Klasse in den zehn Slices hinter diesem) und die Rollenminuten der folgenden Produkt-Slices.

**Entwurf des Wortlauts** (Entwurf des Planners; die Richtung ist Entscheidung des Nutzers vom 2026-10-22, §1; Zeilennummern gelten für `AGENTS.md` am Stand `56ff5194`; der Nutzer bestätigt den Diff vor dem Merge nach `main`).

**(A) Neue Regel, hinter §3.14 (Liefer-Punkt 1).** Die Nummer §3.15 nimmt dieser Slice; `slice-harness-rollenlauf-zeitgrenze-regel` trägt dort §3.16 (Bindung in dessen §1).

```text
### 3.15 Eine eigene harte Regel, die in über der Hälfte der Slices auftritt, wird Sensor oder fällt (seit slice-harness-regeln-rueckbau)

Entscheidung des Nutzers vom 2026-10-22. Gilt für §3.8 bis §3.14 und für jede weitere eigene
Regel dieser Liste, deren Begründung auf eine Klasse des Beobachtungs-Registers zeigt
(`BEO-REPO/<slug>`). Für jede Klasse dieser Art, zu der eine Closure eine Datei unter `evidence/`
anlegt, rechnet der Planner dabei die Rate: die Zahl der zehn zuletzt geschlossenen Slices, zu
denen die Klasse eine Datei unter `evidence/` führt, geteilt durch zehn; gezählt wird ab dem Slice
hinter `slice-harness-regeln-rueckbau`. Liegt die Rate über 0,5, legt der Planner dem Nutzer vor,
die Regel zu streichen oder durch einen Sensor zu ersetzen (ein Gate prüft sie); er nennt die
kleinere Lösung zuerst. Bis zur Entscheidung gilt die Regel.

Wer eine dieser Regeln entfernt oder lockert, nennt im Plan des Slice den Zähler und die Rate ihrer
Klasse; das ist die Konsultation der Herkunft, die der Retirement-Check der Baseline verlangt
(`grundlagen-traceability.md` §Sensor 2). Dass die Beobachtung wieder aufgetreten ist, sperrt das
Lockern nicht.

**Falsch:** Eine Klasse steht in acht von zehn Slices im Register, und der nächste Slice fügt der
Regel einen Absatz hinzu, weil das Review sie wieder fand
(`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: 59 von 73 Produkt-Slices).
**Richtig:** Der Planner legt die Rate vor, und der Nutzer entscheidet über Sensor oder Streichen,
bevor ein weiterer Absatz entsteht.

**Begründung:** Ein Auftreten im Register sagt, dass das Review die Verletzung findet, nicht, dass
die Regel sie verhindert, und nicht, was sie kostet. Die Sätze „Entfernen oder Lockern setzt den
Retirement-Check voraus: Ist die Beobachtung wieder aufgetreten?“ beantwortete jede Regel, die das
Review prüft, mit ja und lasen sich als Sperre (Messung in `slice-harness-regeln-rueckbau` §1).
```

**(B) Die 13 Stellen in §3.9 bis §3.13 (Liefer-Punkt 1): der Satz entfällt, der Rest des Absatzes bleibt.** Kein Herkunfts-Anker `seit …` geht verloren: Jeder steht in der Überschrift oder am Anfang des Absatzes, den er nennt (geprüft für §3.9 bis §3.13 am Stand `56ff5194`).

| Zeilen | Stelle | Satz, der entfällt (Anfang) |
|---|---|---|
| 194–196 | §3.9, erster Absatz | „Entfernen des Zusatzes setzt den Retirement-Check voraus …“ |
| 211–213 | §3.9, Begründung | „Entfernen oder Lockern dieser Regel setzt den Retirement-Check voraus …“ |
| 227–229 | §3.10, Begründung | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ |
| 245–246 | §3.10, Aufzählung und Ortsangabe | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ |
| 265–266 | §3.10, grüner Mutant | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ |
| 312–316 | §3.11, Begründung | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ (samt der Aufzählung der drei Anker) |
| 384–385 | §3.12, Plattform | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ |
| 390–394 | §3.12, Begründung | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ (samt der Aufzählung der vier Anker) |
| 411–413 | §3.13, Begründung | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ |
| 426–427 | §3.13, Bindung | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ |
| 446–447 | §3.13, Ortsangabe | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ |
| 477–478 | §3.13, Schichten | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ und „Für die Zählung des Bootstrap gilt derselbe Check …“ (zwei Sätze, eine Stelle) |
| 489–490 | §3.13, Nachzählen | „Entfernen oder Lockern setzt den Retirement-Check voraus …“ |
| (545–568) | §3.14 | entfallen mit dem Absatz unter (C) |

**(C) `AGENTS.md` §3.14, ersetzt ab Zeile 492 bis 568 (Liefer-Punkt 2).** Unverändert bleiben der Befehlsblock, die Absätze zu Basis, Messpunkt und Pfaden und das erste Falsch/Richtig-Paar mit seiner Begründung (`BEO-REPO/groessenprobe-ohne-messbaren-stand`, ohne den Retirement-Check-Satz); der Rest wird so:

```text
### 3.14 Eine Rückführungs-Marke nennt ihre Metrik (seit slice-harness-gate-kette-stop-hook-liste)

Entscheidung des Nutzers vom 2026-10-13. Die Metrik ist am 2026-10-21 nach [ADR-0056](docs/plan/adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) berichtigt (Entscheidung des Nutzers vom selben Tag). Jede Rückführungs-Marke in §4 eines Plans („zu groß, über N Zeilen mit Tests“) nennt ihre Metrik. **Die Marke ist 1500 Zeilen** (Entscheidung des Nutzers vom 2026-10-22); eine andere Zahl nennt §4 mit ihrem Grund. Gemessen wird mit

    [Befehlsblock unverändert]

[Absätze „Basis“, „Messpunkt“, „Pfade“ unverändert.] **Der Befehl steht hier, einmal:** §4 eines Plans nennt Pfade und Marke und zeigt auf diese Stelle; wiederholt ein Plan den Befehl, gilt der von hier.

Erreicht die Summe die Marke, hält der Implementer an und meldet es dem Planner; der Planner schneidet (`in-progress → next`) oder entscheidet mit dem Nutzer, ohne Schnitt weiterzuarbeiten, und trägt die Entscheidung mit Datum in §7 ein. Gemessen wird, wenn der Diff die Marke zu erreichen scheint oder ein Review sie in Frage stellt. Ob ein Plan zu groß ist, sagen vor dem Anspruch die Liefer-Punkte und die Schichten (Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice), nicht eine Schätzung in Zeilen: Der Planner nennt keinen Faktor und keine Meldeschwelle, der Architect gibt keine Schätzung ab, und §7 trägt keine Messzeile. Ein Plan, der vor `slice-harness-regeln-rueckbau` angelegt wurde, behält Marke, Meldeschwelle und Schätzung seines §4, bis der Planner ihn beim Anspruch nachzieht.

[Erstes Falsch/Richtig-Paar (Summe 1148 gegen 244) und seine Begründung unverändert, ohne den Retirement-Check-Satz.]

**Falsch:** §4 nennt Marke 700 und die Schätzung 376 bis 844 mit dem Faktor 10; die Mitte liegt bei 87 % der Marke, der Planner schneidet vor dem Anspruch in zwei Hälften mit je eigener Schätzung, Marke und Meldeschwelle (`slice-v1-abschluss-verbindungsende-replay-extended`), und die acht Slices zum Verbindungsende des Servers werden fünfzehn.
**Richtig:** §4 nennt „Marke 1500“ und die Pfade; der Slice wird nach Liefer-Punkten und Schichten geschnitten und an der Marke gemessen, wenn der Diff sie erreicht.

**Begründung:** Reviews von 560 bis 941 Zeilen dauerten 5 bis 23 Minuten, und jeder Schnitt kostete einen Plan von rund 300 Zeilen, einen Review, eine Verifikation und einen vollen Gate-Lauf (Messung in `slice-harness-regeln-rueckbau` §1). Die Schätzung in Zeilen bediente nur den Schnitt und lag in 19 Slices unter der Messung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`); die Liefer-Punkte und die Schichten schneiden, wo die Baseline es verlangt.
```

**(D) `AGENTS.md` §3.10 (Liefer-Punkt 3).** Der erste Absatz (Zeilen 217–220) wird:

```text
Wer am Record-/Replay-Pfad (Nachrichten und Bytes zwischen Client, Aufzeichnung und Server), an
einem Format der Aufzeichnung oder am Protokollrand eine Zusage der Spezifikation neu liefert oder
ändert, schreibt zu jeder dieser Zusagen einen Test, der rot wird, wenn der Code die Zusage bricht,
und fährt diese Mutation vor der Übergabe selbst. Der Bericht nennt je Zusage: Zusage · Mutation ·
roter Test. Für eine Option oder einen Hilfetext der Kommandozeile, für Testdaten, Kommentare und
Doku genügt ein Test je Zusage der Spezifikation, ohne Mutation (§3.11). Für den Vertrag eines Gates
oder Werkzeugs unter `harness/vertraege/` gilt die Regel ohne diese Einschränkung.
```

Im Absatz *Aufzählung und Ortsangabe zählen je Element* (Zeile 232) wird „Je Element einer Aufzählung in einer Zusage und je Schritt einer Ortsangabe“ zu „Je Element einer Aufzählung in einer Zusage der Spezifikation und je Schritt einer Ortsangabe darin“. Der Absatz zum grünen Mutanten, die Falsch/Richtig-Paare und die Begründungen bleiben (ohne die Sätze aus (B)).

**(E) `AGENTS.md` §3.12 (Liefer-Punkt 3).** Der erste Satz (Zeilen 320–323) wird, der Rest des Absatzes (ab „Entschieden ist jede vor dem ersten Code-Commit“) bleibt:

```text
Ein Slice, der eine neue Wire-, Datei- oder Formatgrenze liefert (eine Nachrichtenfolge am
Protokollrand, ein Aufzeichnungsformat, eine Datei, die ein Leser einliest), nennt in §6 des
Slice-Plans deren Randformen, etwa fehlender, leerer oder unbekannter Wert, Alias, Abbruch und
Herunterfahren, was ein Diagnose-Feld je Nachrichtenart heißt, Verhalten je Serverversion, und
je Randform, wo sie entschieden ist.
```

Dahinter, nach dem Absatz (nach „… (seit slice-harness-blackbox-kern).“), ein Absatz:

```text
Bei einer Option oder einem Hilfetext der Kommandozeile nennt §6 die Spezifikationsstelle der
Geschwister-Option, deren Regeln die neue übernimmt (zum Beispiel die Hilfe-Angabe in `LH-FA-01.a`),
und was von ihr abweicht; was der Verweis deckt, ist entschieden, eine Abweichung ist eine Randform
wie oben. Für den Vertrag eines Gates oder Werkzeugs unter `harness/vertraege/` gilt die Regel ohne
diese Einschränkung.
```

**(F) Die Träger (Liefer-Punkte 1 bis 3).** `.claude/agents/reviewer.md` führt keine Stufen und zeigt auf den Skill; **der Skill `.harness/skills/reviewer.md` trägt sie** (Zeilen 60–83 unter MEDIUM, Zeilen 85–87 unter LOW).

- *`.harness/skills/reviewer.md`, vier Einträge unter MEDIUM* (Randform im Code-Commit, Adresse nimmt nicht an, äquivalenter Mutant ohne Invariante, Ort oder Sonde, Nehmer nicht nachgezählt): „mindestens MEDIUM, gegen …“ wird je zu „MEDIUM nur mit Failure-Szenario, sonst LOW, gegen …“ mit dem Szenario je Eintrag: der Code entscheidet die Randform anders, als die Spezifikation sie verlangt · die Sendung erreicht ihren Nehmer nicht (ein Ausschluss trifft sie, der Nehmer liegt in `done/`) · die Sonde des Reviewers trennt Original und Mutant · der Nehmer liegt mit der Sendung über drei Liefer-Punkten oder zwei Schichten. *Unter LOW, neuer Eintrag:* „Verstoß gegen `AGENTS.md` §3.10, §3.12 oder §3.13 ohne Failure-Szenario (Entscheidung des Nutzers vom 2026-10-22; seit slice-harness-regeln-rueckbau)“. Die zwei Kommentar-Einträge unter HIGH („Kommentar trägt keine der Kommentar-Klassen“, „Zustandsfeld trägt Chronik“) tragen den Zusatz „nur mit Failure-Szenario, sonst LOW“ (Randform 6, Fassung 6a); die übrigen HIGH-Einträge bleiben.
- *`implement-slice.md`:* Schritt 12 „Liefert der Slice einen neuen Vertrag …“ wird „Liefert der Slice eine neue Wire-, Datei- oder Formatgrenze (`AGENTS.md` §3.12) …“; in der Randform-Rückgabe ein Satz „Bei einer Option oder einem Hilfetext deckt der Verweis auf die Geschwister-Regel in §6 die Randform (`AGENTS.md` §3.12).“; Schritt 19 „Für einen neuen Vertrag gilt das für jede Zusage“ wird „Am Record-/Replay-Pfad, an einem Format und am Protokollrand gilt das für jede Zusage der Spezifikation (`AGENTS.md` §3.10, Geltungsbereich); für eine Option oder einen Hilfetext genügt ein Test je Zusage“ und der Schritt nennt den Geltungsbereich für die Absätze dahinter; Schritt 25 ein Satz „Ohne Welle rechnet die Slice-Closure die Rate je fortgeschriebener Klasse nach `AGENTS.md` §3.15.“.
- *`plan-welle.md` Schritt 6:* „Schätze den Umfang … schneide vor dem Anspruch (`AGENTS.md` §3.14, seit slice-harness-schaetzung-planner-faktor)“ wird „Trage die Marke von `AGENTS.md` §3.14 in §4 ein; geschnitten wird nach Liefer-Punkten und Schichten, nicht nach einer Schätzung in Zeilen“; „Liefert der Slice einen neuen Vertrag, nennt §6 dessen Randformen“ wird „Liefert der Slice eine neue Wire-, Datei- oder Formatgrenze, nennt §6 deren Randformen (bei Option und Hilfetext den Verweis auf die Geschwister-Regel)“.
- *`architect.md`:* „Randformen eines neuen Vertrags“ wird „Randformen einer neuen Wire-, Datei- oder Formatgrenze“; ein Satz „Bei einer Option oder einem Hilfetext prüfst du, dass der Verweis in §6 die Randform deckt.“
- *`close-welle.md` Schritt 3 (Lese-Schritt):* ein Satz „Zugleich rechnet der Planner die Rate je fortgeschriebener Klasse, die eine eigene Regel von `AGENTS.md` trägt, nach §3.15 und legt dem Nutzer vor, was über 0,5 liegt.“

**Reihenfolge:** `AGENTS.md` zuerst (die Quelle), dann die Träger, alles auf dem Zweig des Implementers, ohne Halt vor dem Commit; der Plan folgt in demselben Commit (`AGENTS.md` §3.9). Der Nutzer bestätigt den Diff vor dem Merge nach `main` (Kopf dieses Plans).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Start-Trigger ist eingetreten: `slice-harness-gegenproben-umzug-dev` liegt in `done/` und ist auf `main` gemerged (Stand `edd59b89`); `AGENTS.md` §6 Punkt 6 steht dort geändert (eine Zeile, die Zeilenzahl von `AGENTS.md` blieb, die Zeilennummern der Stellen in §3 gelten unverändert, geprüft am 2026-10-22). Reihenfolge: dieser Slice kommt nach dem Umzug der Gegenproben und vor `slice-version-option` (Auftrag des Koordinators vom 2026-10-22). Der Nutzer hat den Start am 2026-10-22 entschieden („1. ja“, so vom Koordinator übermittelt, auf „Mit ‚ja‘ starte ich den Rückbau-Slice“) und alle sechs Empfehlungen des Planners angenommen (§6). `Verantwortlich:` ist gesetzt; WIP-Limit 1. Die Bewegung nach `in-progress/` und die Basis trägt der Koordinator ein. Hinter diesem Slice kommt der Slice zur Dev-Seite (§1, ohne Kennung).

**Größe** (Metrik nach `AGENTS.md` §3.14 am Stand vor diesem Slice, [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md); der Befehl steht dort): Summe aus Einfügungen und Löschungen über die Pfade `AGENTS.md`, `.claude`, `.harness/skills`, mit den Ausschlüssen des Befehls. Basis ist der Commit, der diese Datei nach `in-progress/` bewegt; **Basis: noch nicht gesetzt** (sie wird beim Übergang eingetragen). **Marke 600, Meldeschwelle 480:** Erreicht die Summe 480, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet; die Entscheidung des Planners steht als datierte Zeile in §7. **Schätzung des Planners vor dem Anspruch** (diese Regel gilt für diesen Plan noch, weil er vor ihr entsteht): `AGENTS.md` 170 bis 230 (13 Stellen mit 2 bis 5 Zeilen löschen, §3.14 etwa 45 Zeilen löschen und 30 schreiben, §3.15 etwa 25 schreiben, §3.10 und §3.12 etwa 25), `.harness/skills/reviewer.md` 25 bis 40, `implement-slice.md` 20 bis 35, `plan-welle.md` 10 bis 20, `architect.md` 6 bis 12, `close-welle.md` 4 bis 8 — zusammen **235 bis 345, Mitte 290, 48 % der Marke 600**, die Obergrenze unter der Meldeschwelle. Faktor Tests zu Produktcode: **keine Messung** (Prosa, kein Testcode). Kein Schnitt: Die Mitte liegt unter 70 % der Marke, der Plan hat drei Liefer-Punkte und zwei Schichten (§8).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt über 600, oder der Wortlaut verlangt mehr als die Stellen aus §3 (ein Sensor, ein MR, ein Absatz in §3.8, §3.9, §3.11 oder §3.13 über die entfallenden Sätze hinaus) — zurück zur Zerlegung (§1).
- `in-progress` → `open` (blockiert): Der Nutzer bestätigt den Wortlaut am Diff vor dem Merge nicht und hat keine Fassung, die er bestätigen würde, oder er will den Retirement-Check ohne Konsultation und Dokumentation aufheben (das ist eine Abweichung von der Baseline und ein MR, §1). Der Slice wartet auf die Entscheidung, ohne dass ein Wortlaut auf `main` in `AGENTS.md` steht; was der Nutzer teilweise bestätigt, nimmt der Merge auf, den Rest nicht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, `make gates` grün, Bestätigung des Diffs durch den Nutzer vor dem Merge nach `main` in §7 mit Datum und Fundstelle, `grep -n 'Retirement-Check' AGENTS.md` findet nur Zeilen in §3.15, Closure-Notiz mit Lerneintrag geschrieben (geschärfte Regel: `AGENTS.md` §3.15, `liegt in`).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen der Regel** (`AGENTS.md` §3.12; der Slice liefert keinen Vertrag von Code, nur Regeln — genannt ist, was offen ist). Entscheidungsort ist der Wortlaut, am Diff vor dem Merge nach `main` bestätigt:

*Randformen 1, 2, 3, 6 und 7 hat der Nutzer am 2026-10-22 entschieden* (wörtlich „1. ja“, so vom Koordinator übermittelt, auf dessen Nachricht „Mit ‚ja‘ starte ich den Rückbau-Slice. Wenn dir eine der sechs Antworten nicht passt, nenne die Nummer.“; er nannte keine Nummer: **alle sechs Empfehlungen des Planners gelten**, der Slice startet). Die Entwürfe in §3 sind damit in diesen Punkten Richtung des Nutzers; der Diff vor dem Merge bleibt zu bestätigen.

1. *Fenster des Raten-Kriteriums* — **entschieden** (Empfehlung des Planners angenommen, „1. ja“): zehn zuletzt geschlossene Slices, ab dem Slice hinter diesem, Schwelle über 0,5, gerechnet nur für Klassen, die die laufende Closure fortschreibt. **Folge für §3:** (A) gilt wie entworfen; der Wortlaut ändert sich nicht.
2. *Die Schätzung des Planners entfällt mit* — **entschieden** (Empfehlung angenommen, „1. ja“): Faktor und 70-%-Schnitt waren ihr einziger Verbraucher; sie entfällt ganz, eine freiwillige Größenordnung gibt es nicht. **Folge für §3:** (C) und `plan-welle.md` Schritt 6 gelten wie entworfen; sie nennen keine Schätzung.
3. *Gates und Werkzeuge im Geltungsbereich* — **entschieden** (Empfehlung angenommen, „1. ja“): §3.10 und §3.12 gelten für den Vertrag eines Gates oder Werkzeugs unter `harness/vertraege/` unverändert, bis die Dev-Seite ihren Regelsatz hat. **Folge für §3:** (D) und (E) gelten wie entworfen (der Satz „ohne diese Einschränkung“ bleibt).
4. *Wer ist „Record-/Replay-Pfad“* — **entschieden im Entwurf:** „Nachrichten und Bytes zwischen Client, Aufzeichnung und Server“ erfasst `record`, `replay` und `play`; die Formulierung folgt der HIGH-Liste des Reviewer-Skills (Korrektheit im kritischen Pfad).
5. *Verweis auf die Geschwister-Regel* — **entschieden im Entwurf:** §6 nennt die Spezifikationsstelle und die Abweichung; der Architect prüft, dass der Verweis die Randform deckt; eine nicht gedeckte ist eine Randform wie bisher.
6. *Kommentarklassen als HIGH im Reviewer-Skill* — **entschieden** (Empfehlung angenommen, „1. ja“): Die zwei Einträge „Kommentar trägt keine der Kommentar-Klassen“ und „Zustandsfeld trägt Chronik“ werden gestuft wie die übrigen: MEDIUM mit Failure-Szenario, sonst LOW (Quelle-Spalte der Produkt-Reviews: §3.11 1 HIGH, 13 MEDIUM, 51 LOW, 27 INFO; §3.7 2 MEDIUM, 35 LOW, 4 INFO). **Folge für §3:** (F) trägt für die zwei Einträge im Skill den Zusatz „nur mit Failure-Szenario, sonst LOW“, Liefer-Punkt 3; ob sie unter HIGH bleiben oder unter MEDIUM stehen, hängt am Baseline-Befund (6a). **Baseline-Befund, am Text geprüft und im Vorschlag falsch gesagt („ändert nur den Skill“):** `grundlagen-harness-dateien.md` §Was ein Kommentar trägt (Feedback-Hälfte) nennt „der HIGH-Eintrag *Kommentar trägt keine der fünf Klassen*“ und „der HIGH-Eintrag *Zustandsfeld trägt Chronik*“ als Träger. Die Stufung weicht davon ab und braucht einen MR (`Ersetzt-Baseline-Regel` diese Stelle), oder die Einträge bleiben unter HIGH mit dem Zusatz „nur mit Failure-Szenario, sonst LOW“ (so sagt es `modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill allgemein: kein HIGH ohne Failure-Szenario). **Offen (6a):** welche der beiden Fassungen der Nutzer will; der Entwurf in (F) nimmt die zweite (Einträge bleiben unter HIGH, kein MR), weil sie ohne Abweichung von der Baseline dieselbe Wirkung für Verstöße ohne Failure-Szenario hat; der Nutzer bestätigt oder wählt den MR am Diff.
7. *Marke der zwei nächsten Produkt-Pläne* — **entschieden** (Empfehlung angenommen, „1. ja“): `slice-version-option` (Marke 400, Meldeschwelle 320, Schätzung 165 bis 315) und `slice-v1-abschluss-verbindungsende-upstream-extended` (Marke 700, Meldeschwelle 550) sind Bestand; beim Anspruch trägt der Planner „Marke 1500“ in §4 ein und streicht Meldeschwelle und Schätzung; die DoD von `slice-version-option` (§3.10 je Zusage mit Mutation) richtet sich dann nach dem neuen Geltungsbereich, eine Option verlangt dort keine Mutation. **Folge für §3:** In (C) endet der Satz zum Bestand mit „, bis der Planner ihn beim Anspruch nachzieht“ (Folgerung des Planners für jeden Bestandsplan, am Diff zu bestätigen); die zwei Pläne nennen die Entscheidung in ihrem Start.
8. *Nummer §3.15* — **entschieden:** dieser Slice vergibt sie; `slice-harness-rollenlauf-zeitgrenze-regel` trägt §3.16 und den Zeiger auf §3.15 statt eines Retirement-Check-Satzes (Bindung in dessen §1, beim Anlegen eingetragen).
9. *Aufruf-Voraussetzung, Aufrufarten des Treibers, Plattform, Stelle einer Nachricht* — nicht berührt (`AGENTS.md` §3.12): Der Slice ändert nur den Geltungsbereich ihres ersten Absatzes.

**Risiken:**

- Der Rückbau holt eine Fehlerklasse zurück: Ohne Mutation je Zusage für Optionen und Hilfetexte, ohne MEDIUM für Regelverstöße ohne Failure-Szenario und mit einer Marke von 1500 findet das Review weniger, und ein Fehler am Record-/Replay-Pfad bleibt unentdeckt — **Ausgang:** weiter offen → `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (ein Auftreten am Record-/Replay-Pfad, ein HIGH mit Code-Defekt oder eine Marke, die ein Review reißt, in den zehn Slices hinter diesem zählt; der Geltungsbereich lässt den Pfad voll stehen).
- Das Raten-Kriterium löst auf dem Bestand aus und erzeugt Vorschläge, die niemand gewollt hat (neue Harness-Arbeit statt weniger) — **Ausgang:** offen bis Closure (erwartet: entfallen, die Zählung beginnt hinter diesem Slice und rechnet nur fortgeschriebene Klassen; Randform 1).
- Der Nutzer ändert am Diff den Wortlaut, streicht einen Teil oder bestätigt ihn nicht — **Ausgang:** offen bis Closure (erwartet: entfallen, der bestätigte Wortlaut steht nach dem Merge in `AGENTS.md`; was er nicht bestätigt, nimmt der Merge nicht auf, und §7 nennt es).
- Der Implementer schreibt ohne Halt vor dem Commit einen Wortlaut, den der Nutzer am Diff ablehnt, und die Arbeit auf dem Zweig fällt weg — **Ausgang:** offen bis Closure (erwartet: entfallen, die Richtung ist entschieden, der Implementer folgt dem Entwurf aus §3, jede Abweichung steht in §7 mit Grund).
- Die pausierten Regel-Slices lesen den alten Wortlaut (Retirement-Check, Schätzung, Meldeschwelle, §3.15) — **Ausgang:** offen bis Closure (erwartet: entfallen, die Pause in §5 der Welle nennt die Wiederaufnahme: Der Planner prüft jeden Plan gegen den neuen Wortlaut; `slice-harness-rollenlauf-zeitgrenze-regel` trägt die Bindung schon).
- Zwei Läufe schreiben gleichzeitig an `AGENTS.md` (`slice-harness-gegenproben-umzug-dev` ändert §6 Punkt 6; der Slice zur Dev-Seite folgt) — **Ausgang:** entfallen: Der Umzug liegt in `done/` und ist gemerged (Stand `edd59b89`), sein Wortlaut in §6 Punkt 6 ist im Ist-Stand, den dieser Slice liest; der Slice zur Dev-Seite folgt später und rechnet den Wortlaut dieses Slice ein (§1).
- Der Slice wächst über Marke oder Meldeschwelle von §4 (Schätzung 235 bis 345 gegen Marke 600) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Hier stehen: die Bestätigung des Diffs durch den Nutzer vor dem Merge nach `main` (Datum, Fundstelle, je Teil), die Abweichungen des Wortlauts vom Entwurf in §3 mit Grund, die Konsultation der Herkunft nach `grundlagen-traceability.md` §Sensor 2 (Zähler und Rate je gelockerte Klasse am Tag des Merge), die Größe in der Metrik von §4 (Basis, letzter Commit, Summe) gegen die Schätzung, die Ausgänge der Register-Klassen aus der DoD und die Ausgänge der Risiken.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-22). Treffer in der Sub-Area `REPO`: die Klassen der Regeln §3.9 bis §3.14 mit den Zählern aus §1 (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` 93×, `BEO-REPO/plan-folgt-korrektur-nicht` 79×, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` 65×, `BEO-REPO/spec-randform-erst-im-review-entschieden` 50×, `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` 27×, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` 19×, `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` 13×, `BEO-REPO/gruener-mutant-faelschlich-aequivalent` 11×), alle `verkörpert` und über der Schwelle; der Slice berührt sie, weil er ihre Regeln lockert, und erhöht keinen Zähler. `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (14×) und `BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet` bleiben unberührt (Kopf und §1 trennen Entscheidung und Empfehlung).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-22, im Commit, der diesen Plan anlegt):** Liefer-Punkte **3** (Raten-Kriterium | Marke | Geltungsbereich und Review-Stufen). Schichten **2**, konservativ gezählt: Briefing (`AGENTS.md`) und Anweisungen der Rollen und Commands (`.claude/`, `.harness/skills/`); nach der Teilung in `slice-harness-d-check-v0-85` §8 wäre es eine (Agenten-Anweisungen). Planung zählt keine Schicht; kein Produkt-Code. Schätzung in der Metrik von §4 (Marke 600, Meldeschwelle 480): **235 bis 345, Mitte 290, 48 % der Marke**; Faktor keine Messung. Der Slice trägt eine Sendung in `slice-harness-rollenlauf-zeitgrenze-regel` ein (die Bindung an §3.15); die Nachzählung steht dort. Kein Ausschluss von §1 trifft die Sendung.
