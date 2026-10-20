# ADR-0055: Zwei Gegenproben haben eine Obergrenze für ihre Zeit und räumen je Mutation auf

**Status:** Proposed

**Datum:** 2026-10-20

**Autor:** Architect. Anlass ist der Hinweis des Nutzers vom 2026-10-20 (so vom Koordinator übermittelt): Dauert eine Gegenprobe drei Stunden, werden keine Gegenproben mehr gefahren. Die Obergrenzen und die Wahl der Mittel unten sind die Empfehlung des Architect; die Annahme dieser ADR durch den Nutzer steht aus.

**Bezug:** [ADR-0048](0048-gegenproben-verlassen-die-gate-kette.md), [ADR-0038](0038-gate-kette-voll-und-schnell.md), [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md)

**Schärft:** [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[ADR-0048](0048-gegenproben-verlassen-die-gate-kette.md) hat die Gegenproben aus der Gate-Kette genommen, weil sie zu lange dauern, und die Pflicht zum Fahren an den Aufrufer gelegt. Die Laufzeit selbst blieb unbegrenzt (`SPEC-060` Punkt 15 (1) und (7): nicht zugesagt). Die zwei größten Gegenproben haben sie seither gesprengt: Der Umbenennungs-Slice brauchte rund zwei Stunden, der Homebrew-Slice fast drei, fast nur für Läufe von `make release-gegenprobe`. Eine Pflicht, die drei Stunden kostet, wird nicht mehr erfüllt; dann trägt der Beleg nichts, den ADR-0048 an sie gebunden hat.

Gemessen am Stand `8dbe6a14` auf dem Gegenproben-Container (20 Kerne, Last 4 bis 5), je ein Lauf mit `GP_MUT`:

| Messung | Wert |
|---|---|
| `release-gegenprobe`, eine Mutation der Gruppe `f` (112 Fälle) | 30 s (12 s user, 21 s sys), rund 50 000 Prozess-Starts |
| `release-gegenprobe`, Gruppe `b` (71 Fälle) | rund 30 s |
| `release-gegenprobe`, Gruppe `sh` (87 Fälle) | rund 6 s |
| `release-gegenprobe`, die unveränderte Kopie | 104 s |
| Mutationen von `release-gegenprobe` | 122 der Gruppe `f`, 65 `b`, 53 `sh`, 7 `p13`, 9 `gp13`, dazu rund 120 des Textes (`p1`, je Zeile zwei; billig) |
| `test-container-gegenprobe`, unveränderte Kopie | 22 s, danach rund 1 s je Mutation (135 Mutationen) |
| Arbeitsordner je Mutation von `release-gegenprobe` | rund 32 MB nach 110 Fällen, keiner wird vor dem Ende des Laufs entfernt |

Daraus folgt: Jede Mutation fährt die ganze Gruppe ihres Skripts (bis 112 Fälle), und jeder Fall startet `release.sh` mehrfach unter Attrappen; das ergibt rund 6300 s Arbeit, bei vier gleichzeitigen Läufen 25 bis 30 Minuten auf ruhigem Host und Stunden unter Last. Die Vermutung des Nutzers, verschachtelte Schleifen mit `sed` und `awk`, trifft nicht zu: `sed` und `awk` sind 2,4 Prozent der Prozess-Starts; es sind `bash`-Unterprozesse (5200), `mkdir` (2000), `rm` (1600), `jq` (1450), `sha256sum` (1440), `tr`, `cat`, `date` und `tar` der Attrappen von `docker` und der Layout-Erzeugung. Mehr als die Hälfte der Zeit ist Systemzeit der Prozess-Starts aus einer großen `bash`. `test-container-gegenprobe` läuft nacheinander mit einer Auswahl der Fälle je Mutation und liegt nach Schätzung nahe an fünf Minuten; gemessen ist nur ihr Anteil je Mutation.

## Entscheidung

Wir empfehlen **Option B: eine harte Obergrenze für den vollen Lauf von zwei Gegenproben, Anzeige und Rot bei Überschreitung, Aufräumen je Mutation, keine Streichung einer Mutation.**

1. **Obergrenzen.** `release-gegenprobe` 600 s, `test-container-gegenprobe` 300 s, Wanduhr des vollen Laufs auf ruhigem Host. Das Skript misst sich selbst, nennt `Zeit`, Obergrenze und Last in einer Zeile und endet bei Überschreitung mit Exit 1. Die übrigen Gegenproben bekommen eine Obergrenze, wenn ein eigener Slice sie gemessen hat.
2. **Aufräumen.** Der Ordner einer Mutation und der eines Falls entfallen, sobald sie fertig sind, rot wie grün. Die Spitze des Platzes liegt unter 100 MB.
3. **Die Mutationen bleiben.** Jede ist die Mutation einer Zusage oder eines Elements (`AGENTS.md` §3.10); das Budget wird über die Kosten je Mutation erreicht, nicht über deren Zahl. Der Slice senkt die Prozess-Starts je Fall (Arbeit in der `bash` statt in Unterprozessen, Layouts und Attrappen einmal je Welt statt je Aufruf, ein Ordner je Fall statt einer Kopie mit `cp -a`) und erhöht höchstens die Zahl gleichzeitiger Läufe auf 8.
4. **Reicht das nicht**, entscheidet der Nutzer zwischen zwei Lockerungen, jede mit eigener ADR: die Obergrenze erhöhen, oder je Mutation nur die Fälle prüfen, die sie berühren kann (Auswahl nach den benannten roten Fällen und ihren Nachbarn, wie `test-container-gegenprobe` es mit `nur` hält); letzteres prüft „genau die benannten Fälle rot“ nur noch unter der Auswahl.

Die Regeln im Einzelnen (Randformen, Ausgabe, Reihenfolge der Gründe, Fälle) stehen in [`SPEC-060`](../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 16.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — keine Zusage zur Zeit, nur Messung im Bericht (Bestand, `SPEC-060` Punkt 15 (7)) | nichts zu ändern | genau die drei Stunden; die Pflicht aus ADR-0048 wird nicht mehr gefahren |
| **B — Obergrenze, Rot bei Überschreitung, Kosten je Mutation senken** | belegt Ergebnisgleichheit (deklarierte rote Mengen bleiben, ein grüner Lauf prüft sie); Lauf in Minuten | die Obergrenze gilt nur auf ruhigem Host; ein Lauf unter Last ist rot ohne Befund am Code |
| C — Mutationen streichen oder zusammenlegen | schnellster Weg | senkt die Prüfung, nicht die Laufzeit (`AGENTS.md` §3.6, §3.10: je Zusage und Element eine Mutation) |
| D — nur die betroffenen Fälle je Mutation fahren | senkt die Zeit am stärksten (Faktor 10 und mehr) | „genau die benannten Fälle rot“ gilt nur unter der Auswahl; eine Mutation, die einen Fall außerhalb färbt, bleibt unbemerkt; Lockerung nach `AGENTS.md` §3.6 |
| E — Nachtlauf in der CI | kein Warten | andere Außenwirkung; getragen von den offenen Plänen für CI-Nachtläufe, nicht von dieser ADR |
| F — Obergrenze ohne Rot, nur Hinweis | kein falsches Rot unter Last | ohne Folge liest niemand den Hinweis; die Laufzeit wächst mit jeder Mutation weiter (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`) |

## Konsequenzen

- Positiv: Der volle Lauf bleibt in Minuten fahrbar; Platzbedarf und Laufzeit sind eine Zusage mit Fällen, nicht eine Beobachtung im Bericht. Wer eine Mutation hinzufügt, sieht an der Zeile, was sie kostet.
- Negativ: Unter Last kann ein grüner Lauf wegen der Zeit rot enden; die Zeile nennt die Last, wer sie liest, fährt auf ruhigem Host neu.
- Akzeptiertes Negativ: Der Vergleich `T > B` gegen `T >= B` ist nicht prüfbar (die Zeit ist nicht steuerbar, die Wirkung ist eine Sekunde).
- Akzeptiertes Negativ: Die Zahlen gelten für Linux. Grund: Die Gegenproben laufen im Gegenproben-Container (`SPEC-060` Kopf, Plattform).
- Übergang: Bis der Slice zu dieser ADR geliefert ist, fährt kein Implementer einen vollen Lauf der beiden Gegenproben; er fährt den engsten Fall mit `GP_MUT`, der Verifier den vollen Lauf einmal am Endstand. Das ist eine Anweisung der Rollen, kein Vertragstext; die Zeile für die Rollen steht im Slice-Plan.
- Folgepflicht: `tools/harness/release-gegenprobe.sh`, `tools/harness/container-gegenprobe.sh`; `harness/README.md` §Sensors (beide Zeilen nennen Obergrenze und Aufräumen erst mit der Lieferung); ein Slice für die übrigen Gegenproben (Kennung vom Planner). Lastenheft unverändert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, git, Attrappen | Obergrenze, Ausgabe, Reihenfolge der Gründe, Arbeitsordner je Mutation entfernt (Fälle nach `SPEC-060` Punkt 16 (9)) | `make release-gegenprobe` |
| bash, make, Attrappe von `docker` | dieselben Fälle für `test-container-gegenprobe` | `make test-container-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn ein voller Lauf auf ruhigem Host die Obergrenze trotz sinkender Kosten je Mutation erneut reißt; wenn eine dritte Gegenprobe mehr als zehn Minuten braucht; wenn der Nutzer einen Nachtlauf in der CI beschließt (dann ist die Obergrenze neu zu setzen).
