# ADR-0058: Die Zeitgrenze eines Rollenlaufs setzt ein deterministischer Wächter im Werkzeug-Aufruf durch

**Status:** Accepted

**Datum:** 2026-10-22

**Autor:** Architect; Entscheidung des Nutzers vom 2026-10-21 (wörtlich, vom Koordinator übermittelt): „Wir haben rollen (siehe .claude/agents) und commands -- max. 1 Stunde“, dazu sein „ok“ auf den Entwurf des Mechanismus (Wächter im Werkzeug-Aufruf, bash und awk, kein Hintergrundprozess, Datenquelle die vorhandene Erfassung) und seine Entscheidung vom 2026-10-22 („deine empfehlungen“), die Stunde zugleich als zweite Rückführungs-Marke zu führen. Die **Wahl der Mittel** unten (Sperrmenge, Fail-Richtung, Anker, Wanduhr, Weg über die Grenze, Umgang mit den pflichtigen Läufen) ist die **Empfehlung** des Architect innerhalb dieser Entscheidung; die Annahme dieser ADR ist die Entscheidung des Nutzers.

**Bezug:** [ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md) (Zeitbudget der Gegenprobe, 600 s; die Stunde ist eine zusätzliche Grenze darüber), [ADR-0048](0048-gegenproben-verlassen-die-gate-kette.md) (Ort der Gegenprobe), [ADR-0038](0038-gate-kette-voll-und-schnell.md) (die Hooks dieses Repos und ihr Zusammenspiel), [ADR-0056](0056-groessen-metrik-misst-den-gelieferten-diff.md) (Metrik der Größenmarke, neben der die Stunde als zweite Marke steht); kein Anspruch des Produkts, die Regel gilt für den Harness.

**Schärft:** [`SPEC-070`](../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Ein Rollenlauf hatte keine Zeitgrenze. [ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md) budgetiert die Gegenprobe, die Rückführungs-Marke eines Plans misst Zeilen ([ADR-0056](0056-groessen-metrik-misst-den-gelieferten-diff.md)), und `.claude/commands/implement-slice.md` trug keine Zeitregel. Der Implementer-Lauf eines Slice dauerte 4,34 Stunden (Wanduhr, 48 % davon Befehlszeit), ohne dass eine Marke feuerte; Review und Verifikation sahen erst den fertigen Stand.

Gemessen an der Erfassung des Hauptbaums (Stand 2026-10-22, 854 Ströme, 39 569 Zeilen; die Zahlen im Einzelnen stehen in der Spezifikation): Ein Agent kann fortgesetzt werden, und die Erfassung schreibt dann einen neuen `SubagentStart` in denselben Strom; gerechnet ab dem **letzten** Start dauern 10 von 849 Läufen länger als eine Stunde, ab dem **ersten** 100 Ströme. Die Hauptsitzung hat keinen `SubagentStart`. Der längste Bash-Aufruf des Bestands dauerte 601 s; zwischen zwei Werkzeug-Aufrufen feuert kein Hook.

Das Budget des Nutzers verlangt dreierlei: Der Lauf hält bei der Grenze an und berichtet; er wird dabei nicht abgebrochen und nichts wird zurückgenommen; und das Mittel ist deterministisch, keine Disziplin. Offen ist, wie weit der Wächter gehen darf, ohne der Rolle den Bericht zu nehmen, und was er tut, wo er nichts sieht.

## Entscheidung

Wir wählen **Option C: ein `PreToolUse`-Hook, der den Strom der Erfassung des laufenden Agenten liest und ab der Grenze die Aufrufe sperrt, die den Lauf verlängern, und die Abwicklung durchlässt.** Im Einzelnen (Regeln, Listen, Fälle: [`SPEC-070`](../../../harness/vertraege/SPEC-070-rollenlauf-zeitgrenze.md)):

1. **Im Aufruf, nicht daneben.** Nur der Hook im Werkzeug-Aufruf der Rolle kann sie anhalten; ein Hintergrundprozess kann zusehen und nichts weiter, und mit Signalrecht wäre er der Abbruch, den der Nutzer ausgeschlossen hat. Er liest nur, schreibt keine Zeile und keine Datei und ist kein Gate.
2. **Beginn und Zeit.** Der Beginn des Laufs ist der **letzte** `SubagentStart` im Strom des Agenten, gerechnet in Wanduhr. Die Fortsetzung eines Agenten durch den Koordinator beginnt damit eine neue Stunde; Werkzeugzeit zählt nicht, weil die Stunde des Nutzers die ist, die der Koordinator wartet.
3. **Kleine Sperrmenge.** Gesperrt sind die Werkzeuge, die Arbeit schreiben oder Arbeit starten, und die Befehle, die erkennbar Zeit kosten; alles andere bleibt frei, namentlich Lesen, Berichten und `git`. Eine zu enge Freigabe nähme der Rolle den Bericht, den die Entscheidung ausdrücklich verlangt; eine zu weite lässt Arbeit weiterlaufen, und das ist der kleinere Fehler. Die Prüfung ist ein Stolperdraht, keine Sandbox.
4. **Fail-open, mit Selbstauskunft.** Jede Unsicherheit des Hooks (keine Erfassung, kein Strom, kein Anker, unlesbare Eingabe, fehlendes Werkzeug) ist eine Freigabe mit einer Zeile auf stderr, die den Grund nennt. Das ist die **umgekehrte** Richtung des Wächters nach `AGENTS.md` §3.1, mit Grund: Jener braucht nur den Text des Befehls, der immer da ist, und ein verfehlter Block ist ein installierter Paketmanager, der nicht zurückzunehmen ist. Dieser braucht gitignorierten Zustand, der in einem frischen Klon, nach einem Aufräumen und in einem Worktree fehlen kann, und ein Block aus Zweifel hielte dort jede Arbeit jedes Subagenten an, still. Ein verfehlter Block heißt nur, dass diese Grenze für diesen Lauf nicht greift.
5. **Kein Hebel im Hook.** Es gibt keine Variable, keine Datei und keinen Satz im Auftrag, die die Grenze heben; sonst wäre sie eine Bitte. Der Weg über die Grenze ist die Entscheidung des Koordinators nach dem Bericht (Fortsetzung oder neuer Lauf, Punkt 2).
6. **Die pflichtigen Läufe passen in die Stunde, bis auf zwei Fälle.** `make gates` und der volle Lauf einer aufgenommenen Gegenprobe zusammen belegen höchstens rund 20 Minuten. Nicht in die Stunde neben der Arbeit passt der volle Lauf einer noch nicht aufgenommenen Gegenprobe (25 Minuten und mehr) und der Lauf der Gegenproben am Ende einer Welle. Für den ersten bleibt die Ausnahme von `AGENTS.md` §6 (Implementer: engster Fall, Verifier: einmal der volle Lauf), und der Verifier startet den Lauf zuerst und berichtet ihn, wenn die Stunde endet, ohne dass er endete, als **Offenes**, nicht als grün; für den zweiten ist es die Einlieferung bzw. der Koordinator. Die Pflicht zum Fahren aus [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 14 ruht dadurch nicht.
7. **Was nicht gedeckt ist, steht in der Spezifikation:** die Hauptsitzung (ein Command, das in ihr läuft), die Zeit innerhalb eines laufenden Vordergrund-Befehls, ein Befehl im Hintergrund, die Denkzeit und ein Agent, der die Grenze umgehen will. Plattform: geprüft ist Linux; macOS und Windows sind Grenze.
8. **Die Stunde ist zugleich eine zweite Rückführungs-Marke** neben der Zeilenmarke (Entscheidung des Nutzers vom 2026-10-22): Ein Slice, der in einer Stunde nicht umsetzbar ist, ist zu groß, und das ist eine Aussage über den Schnitt, nicht über den Lauf. Der Wortlaut steht in der Regel von `AGENTS.md`; der Wächter schreibt kein Ereignis, die Marke stützt sich auf den Bericht der Rolle.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nur die Regel an den Leser (Text in `AGENTS.md` und den Rollen) | nichts zu bauen, nichts zu pflegen | genau der Zustand, der den Anlass gab: kein Wächter, die Disziplin trägt; ein Lauf, der eine Stunde nicht bemerkt, liest die Regel nicht neu |
| B — Hintergrund-Beobachter über dem Strom | sieht auch die Zeit innerhalb eines Befehls | kann nur zusehen; ein Prozess mit Signalrecht wäre der Abbruch, den der Nutzer ausschließt, und ein eigener Prozess mit eigener Pflege und eigenen Ausfällen |
| **C — Hook im Aufruf, kleine Sperrmenge, fail-open** | hält die Rolle an, ohne den Bericht zu nehmen; deterministisch; ein Leser ohne Zustand | wirkt nur zwischen zwei Aufrufen und nur, wo die Erfassung liegt; eine Umgehung durch Wortwahl ist nicht erkannt |
| D — Hook, der den Lauf beendet | die Grenze ist hart | verliert den Bericht, die einzige verwertbare Ausgabe eines angehaltenen Laufs; vom Nutzer ausgeschlossen |
| E — Hook mit Hebel (Variable, Datei, Satz im Auftrag) | ein legitim längerer Lauf geht weiter | die Grenze wird zur Bitte; wer sie heben will, hebt sie |
| F — Hook, der bei Zweifel sperrt (fail-closed) | keine Lücke bei fehlender Erfassung | hält in einem frischen Klon, einem Worktree und nach einem Aufräumen die Arbeit aller Subagenten an, ohne dass jemand die Ursache sähe |
| G — Beginn aus dem ersten Start, Werkzeugzeit statt Wanduhr | einfacher zu lesen (erste Zeile) | ein fortgesetzter Agent bliebe gesperrt (100 Ströme überschreiten die Stunde ab dem ersten Start, 10 ab dem letzten); die Werkzeugzeit ist ein Anteil der Wanduhr (48 % im Anlass) |
| H — alles außer einer kleinen Freigabeliste sperren | die Arbeit hört sicher auf | jede Commit-Message mit einem Wort, jeder `ssh`-Aufruf in den Arbeitsklon eines Dev-Containers nähme der Rolle den Bericht |

## Konsequenzen

- Positiv: Ein Lauf hält an und berichtet, und das hängt nicht an seiner Disziplin; die Fortsetzung ist eine Entscheidung des Koordinators, sichtbar als zweiter Start im Strom. Der Hook schreibt nichts, hat keinen Zustand und ändert die Erfassung nicht.
- Negativ: Die Grenze wirkt erst beim ersten Aufruf nach der Stunde (höchstens ein Vordergrund-Befehl später, gemessen 601 s); in der Hauptsitzung wirkt sie gar nicht. Ein Wort in einer Commit-Message (`make`, `docker`, `sleep`) sperrt den Aufruf nach der Grenze; die Meldung sagt es.
- Akzeptiertes Negativ: Wo die Erfassung nichts schreibt, schweigt der Hook und nennt es nur auf stderr; ob Claude Code das anzeigt, ist nicht zugesagt. Grund: ein Block aus Zweifel wäre schlimmer (Punkt 4), und die Gegenmaßnahme ist die Selbstauskunft.
- Akzeptiertes Negativ: Die Rolle kann den Strom ändern und damit ihren Anker; die Datei ist gitignored und nicht zugriffsbeschränkt. Grund: Die Grenze ist für einen Agenten gedacht, der sie befolgen soll, nicht für einen, der sie umgehen will.
- Akzeptiertes Negativ: Die Plattformen macOS und Windows sind Grenze. Grund: kein solcher Host steht zur Messung bereit; der Hook läuft, wo Claude Code läuft, und nutzt nur Mittel, die auf allen drei Plattformen vorhanden sind.
- Folgepflicht: Hook, Leser der Eingabefelder, Eintrag in `.claude/settings.json`, Gegenprobe mit Fragment und die Zeilen in `harness/README.md` (Werkzeug, Gegenprobe); die Zeile der Gegenprobe in der Tabelle von [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16. Die Regel in `AGENTS.md` (mit der Stunde als zweiter Rückführungs-Marke, dem Satz zum Lauf des Verifiers im Hintergrund und `.claude/agents/` in der Liste der Harness-Artefakte von §3.9) und je ein Satz in den sechs Rollen-Typen; der Wortlaut steht beim Nutzer. Kennungen der Slices vergibt der Planner.
- Folgepflicht, eigene Entscheidung: Die Wiederholungs-Grenze für volle Gegenprobe-Läufe (Entscheidung des Nutzers vom 2026-10-21, ein eigener Slice) erbt Datenquelle, Anker und Fail-Richtung dieser ADR; ein Anker für die Hauptsitzung ist eine eigene Entscheidung, nicht ein Zusatz dieses Hooks.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, awk, Attrappen-Ströme | je Zusage und je Element der Aufzählungen von `SPEC-070` ein Fall und eine Mutation, die nur ihn färbt; Zeitbudget nach `SPEC-060` Punkt 16 | `make zeitgrenze-gegenprobe` (steht in `GEGENPROBEN_CONTAINER`) |
| Agentenlauf | der Hook sperrt und lässt durch gegen den eigenen Strom des Implementers (Handlung, kein Sensor, `SPEC-070` Punkt 13) | — |

## Re-Evaluierungs-Trigger

Wenn die Selbstauskunft des Hooks in Läufen auftaucht, in denen die Erfassung da sein sollte (dann trägt die Fail-Richtung nicht); wenn `BEO-REPO/rollenlauf-ohne-zeitgrenze` nach der Lieferung ein weiteres Auftreten zählt (dann ist die Hauptsitzung oder die Umgehung die Lücke); wenn Claude Code die Eingabe von `PreToolUse` oder die Form des Blocks ändert; wenn die Erfassung ein Ereignis für das Ende eines Laufs schreibt; wenn der Nutzer die Stunde ändert.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-22 | Proposed (Empfehlung des Architect; Entscheidungen des Nutzers vom 2026-10-21 und 2026-10-22) | `docs/plan/planning/open/slice-harness-rollenlauf-zeitgrenze.md` |
| 2026-10-22 | Accepted (Entscheidung des Nutzers, Wortlaut „adr accept“) | Vorlage der Koordination am 2026-10-22 |
| 2026-10-22 | Geltungsbereich bestätigt (Entscheidung des Nutzers, Wortlaut „2. ja“ auf die Frage der Koordination, ob der Wächter nur für Rollenläufe gilt und Commands und Hauptsitzung im Regeltext als nicht erzwungen stehen) | Frage 2 der Koordination am 2026-10-22 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen entstehen als neue ADR mit `Supersedes ADR-0058`.
