# ADR-0059: Die Obergrenze für den vollen Lauf einer Gegenprobe steigt von zehn auf zwanzig Minuten

**Status:** Accepted

**Datum:** 2026-10-25

**Autor:** Architect; angenommen durch Entscheidung des Nutzers vom 2026-10-25 (Rückfrage des Koordinators, Wahl „Obergrenze per ADR auf 1200 s“): Die Obergrenze der Zeit einer Gegenprobe wird von 600 s auf 1200 s angehoben; jeder Fall und jede Mutation bleiben.

**Bezug:** [ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md) (Punkt 1 wird ersetzt, ADR-0055 selbst bleibt unverändert), [ADR-0048](0048-gegenproben-verlassen-die-gate-kette.md), [ADR-0058](0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md)

**Ersetzt:** Punkt 1 („Obergrenze“) von [ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md), nur die Zahl 600 s und das Wort „zehn Minuten“; die Punkte 2 bis 5 von ADR-0055 gelten unverändert.

**Schärft:** [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0055](0055-gegenproben-haben-ein-zeitbudget.md) Punkt 5 sah vor: Reicht das Senken der Kosten je Mutation nicht, entscheidet der Nutzer zwischen einer höheren Obergrenze und einer Auswahl der Fälle je Mutation, jede mit eigener ADR; der Re-Evaluierungs-Trigger nennt den Fall, dass der Nutzer die zehn Minuten ändert. `slice-harness-gegenproben-schneller` hat die Kosten gesenkt und gemessen (Plan §7): Die Stichprobe fiel von 73 s auf 29 s (Faktor 2,5). Die Gruppe `h` (Mutationen `m-p19-*`, rund 100) lässt sich kaum beschleunigen (510 s auf 500 s bei `GP_PARALLEL=8`), weil ihre Fälle echtes `git` gegen ein Probe-Tap fahren. Die Hochrechnung für den vollen Lauf von `release-gegenprobe` ergibt rund 985 s bei `GP_PARALLEL=8`; unter 600 s liegt er nicht, unter 1200 s mit Abstand. <!-- d-check:status-provenance -->

## Entscheidung

Wir wählen **Option A: Die Obergrenze für den vollen Lauf jeder Gegenprobe ist 1200 s (zwanzig Minuten); jeder Fall und jede Mutation bleiben.**

Sie ersetzt Punkt 1 von ADR-0055 für jede Gegenprobe gleich („für jede gleich“ bleibt); alles andere an ADR-0055 gilt weiter: Anzeige und Rot bei Überschreitung, Aufnahme, Aufräumen je Mutation, keine Streichung einer Mutation. Die Regeln im Einzelnen (Vorgabe von `GP_BUDGET`, Ausgabe, Reihenfolge der Gründe, Fälle) stehen in [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — Obergrenze für alle auf 1200 s** | trägt die gemessene Hochrechnung (985 s) mit Abstand; Ergebnisgleichheit bleibt (jede Mutation, jeder Fall, dieselben roten Mengen); eine Zahl für jede Gegenprobe, keine Sonderregel | ein voller Lauf kostet bis zu 20 Minuten; Rollenläufe, die einen vollen Lauf fahren, haben weniger Rest in der Stunde von ADR-0058 (Folge unten) |
| B — Gruppe `h` mit eigener Grenze | der Rest bliebe bei 600 s | zwei Zahlen in einem Lauf, ein Teil-Budget je Gruppe, das niemand gefordert hat; die Obergrenze gilt je Gegenprobe, nicht je Gruppe (ADR-0055 Punkt 1), und die Messung der Gruppen einzeln wäre ein eigener Mechanismus |
| C — je Mutation nur Fälle auswählen, die sie berühren kann | größter Zeitgewinn (Faktor 10 und mehr) | „genau die benannten Fälle rot“ gilt nur noch unter der Auswahl; Lockerung nach `AGENTS.md` §3.6 mit Verlust an Prüfung; der Nutzer hat sie nicht gewählt |
| D — Slice mit Teilziel schließen (ein Lauf über 600 s bleibt, Tabelle `nicht gemessen`) | keine Lockerung | der Zustand `nicht gemessen` ruht auf einem Carveout (`CO-001`, `CO-003`); die Pflicht am Wellenende und vor dem Release bliebe unerfüllbar fahrbar, was ADR-0055 gerade beenden sollte |
| E — Obergrenze nicht ändern, Slice zurück in `next` (Rückführung des Plans) | hält die Zahl | die Messung zeigt, dass weiteres Beschleunigen vor allem Gruppe `h` trifft, die an echtem `git` hängt; Aufwand ohne belegte Wirkung |

## Konsequenzen

- Positiv: Der volle Lauf von `release-gegenprobe` und `test-container-gegenprobe` kann mit allen Fällen und Mutationen aufgenommen werden; keine Prüfung wird gesenkt.
- Negativ: Die Zahl gilt für jede Gegenprobe, auch für schnelle; eine Gegenprobe, die bisher unter 600 s lag, kann bis 1200 s wachsen, ohne rot zu werden. Das Wachstum zeigt die Zeile `Zeit <T> s von höchstens <B> s` (`SPEC-060` Punkt 16 (2)).
- Folge für die Stunde eines Rollenlaufs: [ADR-0058](0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) und `SPEC-070` Punkt 10 rechneten mit 600 s für den vollen Lauf einer aufgenommenen Gegenprobe neben den Gates. Mit 1200 s bleiben im ungünstigen Fall (Gates 600 s, Gegenprobe 1200 s) rund 30 Minuten für die Arbeit; `SPEC-070` Punkt 10 trägt die neue Rechnung. Die Zeitgrenze von einer Stunde selbst bleibt.
- Akzeptiertes Negativ: Die Zahlen gelten für Linux, wie in ADR-0055. Der Vergleich `T > B` gegen `T >= B` bleibt nicht prüfbar (ADR-0055).
- Folgepflicht: `SPEC-060` Punkt 16, `SPEC-055` Punkt 12, `SPEC-069` Punkt 17 und `SPEC-070` zeigen auf diese ADR; die Skripte der aufgenommenen Gegenproben und ihre Fälle für die Vorgabe tragen die neue Zahl (Slice, der sie liefert). Lastenheft unverändert.
- Der Re-Evaluierungs-Trigger von ADR-0055 „wenn der Nutzer die zehn Minuten ändert“ ist mit dieser Entscheidung eingetreten und erledigt.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, git, Attrappen | Vorgabe der Obergrenze, Zeile mit Zeit, Rot bei Überschreitung (Fälle nach `SPEC-060` Punkt 16 (9)) | `make release-gegenprobe` |
| bash, make, Attrappe von `docker` | dieselben Fälle | `make test-container-gegenprobe` |
| bash, git, make | dieselben Fälle | `make container-einrichten-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn ein voller Lauf auf ruhigem Host 1200 s erneut reißt; wenn der Nutzer die zwanzig Minuten ändert; wenn die Zeitgrenze eines Rollenlaufs (ADR-0058) geändert wird, weil der volle Lauf neben den Gates nicht mehr trägt.
