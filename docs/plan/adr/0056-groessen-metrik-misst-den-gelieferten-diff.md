# ADR-0056: Die Größen-Metrik der Rückführungs-Marke misst den gelieferten Diff, nicht die Bücher des Slice

**Status:** Proposed

**Datum:** 2026-10-21

**Autor:** Architect; Auftrag des Nutzers vom 2026-10-21 (wörtlich, vom Koordinator übermittelt): „Metrik-Lücke beheben“. Die Fassung unten ist die **Empfehlung** des Architect; die Annahme und der Wortlaut in `AGENTS.md` sind die Entscheidung des Nutzers.

**Bezug:** [ADR-0032](0032-kopf-sensor-fuer-slice-plaene.md) (Form der Pläne); kein Anspruch des Produkts — die Regel gilt für den Harness. Bindung an das Register: `BEO-REPO/groessenprobe-ohne-messbaren-stand` (5×, Stand verkörpert), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`.

**Schärft:** — *(Prozess-ADR ohne Spec-Stratum; Träger der Regel ist `AGENTS.md` §3.14)*

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`AGENTS.md` §3.14 bindet jede Rückführungs-Marke in §4 eines Plans an eine Messung: die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice>` über die Pfade, die der Plan nennt. Der Zweck ist eine Zahl, die zwei Leser gleich erhalten. `BEO-REPO/groessenprobe-ohne-messbaren-stand` zählt 5 und nennt zwei Auftreten, in denen die Formel das nicht hielt: ob die Vertragsdateien des Architect unter einem genannten Pfad zählen (zwei Lesarten, 1062 gegen 783) und eine Basis, die vor einem Commit des Architect liegt (899 gegen 891).

Zwei weitere Lücken sind an der Closure des Homebrew-Nachweises gemessen (Befund des Planners, vom Architect nachgemessen; Basis `2c56dabb`, Pfade `tools harness spec docs README.md`, Stände des Closure-Zweigs):

| Stand | Formel von heute | Befehl der Entscheidung unten |
|---|---|---|
| `ba84090c` — gelieferter Stand, Ende der Arbeit | 907 | 223 |
| `3d74dd13` — Verifikationsbericht abgelegt | 1017 | 223 |
| `cd88df3e` — Berichtigungen der Verifikation | 1053 | 244 |
| `4bed49bd` — Closure-Notiz, Register | 1148 | 244 |
| `d950b370` — reiner `git mv` nach `done/` | 1148 | 244 |
| `fd140c11` — Verweise nachgezogen (Endstand) | 1148 | 244 |

1. **Der Zustandswechsel zählt als Arbeit.** Am Endstand steht die Plandatei zweimal im Diff — `464 0` unter `done/` und `0 207` unter `next/` —, weil ihre Ähnlichkeit nach den Closure-Änderungen unter die Schwelle der Umbenennungs-Erkennung fällt: 671 Zeilen für eine Datei. Ob sie fällt, hängt am Umfang der Closure: Für einen anderen geschlossenen Slice dieses Repos (Spanne `ff5c553a..5b36ae96`) paart `git` dieselbe Bewegung und zählt `124 13`. Dieselbe Mechanik ergibt damit zwei Größenordnungen, je nach Inhalt — die Zahl ist keine Eigenschaft der Arbeit.
2. **Die Formel zählt, was die Größenregel ausnimmt.** Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice nimmt Gate-Läufe, Review-Report, Closure-Notiz, Register, Risiko-Ausgänge und die drei Paarungen von der Größenzählung aus. In der Spanne oben zählt die Formel von heute Review- und Verifikationsbericht mit 179 Zeilen (69 und 110) und die Belege des Registers mit 20. Die Closure rechnete die Ausnahme von Hand heraus und schrieb dazu, die Zeilen der Notiz, die Belege des Registers und die Zeile, die die Zahl einträgt, zähle ein Plan nicht mit: Die Ausnahme war in Gebrauch, nur nicht in der Regel.

Dazu hängt die Zahl an der Konfiguration des Lesers: Über die Spanne bis `fd140c11` ergibt sie ohne Umbenennungs-Erkennung 1148, mit der Vorgabe von `git` 778. Der Rest der Konfiguration ist ohne Gewicht: Zwischen `myers`, `minimal`, `patience` und `histogram` liegen über eine Spanne von 10158 Zeilen zwei Zeilen Unterschied.

**Der Bestand ist klein.** 20 lebende Pläne nennen die Formel (`grep -l numstat` über `open/`, `next/`, `in-progress/`); 19 davon nennen Pfade außerhalb des Planungsbaums und messen heute schon, was unten entschieden wird. Einer nennt `docs` und liegt in der Closure.

## Entscheidung

Wir empfehlen **Option C: der Endstand bleibt der Messpunkt, der Befehl trägt die Ausschlüsse und die Umbenennungs-Schwelle, und die Formel steht an einer Stelle.**

1. **Der Befehl.** Gemessen wird mit

   ```
   git diff --numstat --find-renames=10% <Basis>..HEAD -- <Pfade des Plans> \
     ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' \
     | awk '{a+=$1; b+=$2} END {print a+b+0}'
   ```

2. **Basis** ist der Commit, der die Plandatei nach `in-progress/` bewegt; seinen Hash trägt §4 beim Übergang ein. Was davor liegt — Vertrag und Entscheidung des Architect vor dem Start —, zählt nicht. Ob die Basis dieser Commit oder sein Elternteil ist, ändert die Zahl nicht: Beide berühren nur ausgeschlossene Pfade (gemessen: 244 für die drei Kandidaten der Spanne oben).
3. **Messpunkt** ist der committete Stand des Zweigs, auf dem der Slice läuft, zu jeder Zeit der Arbeit — nicht erst der letzte Commit. Eine Marke, die erst nach der Arbeit messbar ist, löst keine Rückführung aus, und die Meldeschwelle braucht die Zahl vorher. Nicht gezählt ist, was noch nicht committet ist.
4. **Pfade.** Die Pfade bleiben die des Plans; die zwei Ausschlüsse stehen in jedem Fall im Befehl, damit eine weite Pfadangabe (`docs`, `.`) den Planungsbaum und die Reports nicht wieder einfängt. Nicht gezählt sind damit Plan, Closure-Notiz, Roadmap, Welle-Plan, Beobachtungs-Register, Review- und Verifikationsbericht. Gezählt bleiben Entscheidungen (`docs/plan/adr/`, `docs/plan/carveouts/`) sowie Nutzer- und Wartungs-Doku: Sie wachsen mit dem Umfang der Arbeit.
5. **Die Umbenennungs-Schwelle 10 % ist ein Pin, keine Heilung.** Sie nimmt die Zahl aus der Konfiguration des Lesers heraus und senkt die Klippe so weit, dass eine bewegte Datei mit ihrer Bearbeitung zählt statt mit ihrer Größe. Die Doppelzählung der Plandatei löst der Ausschluss, nicht dieses Flag.
6. **Träger ist `AGENTS.md` §3.14, einmal.** §4 eines Plans nennt Pfade, Marke und Meldeschwelle und zeigt auf diese Stelle; wiederholt ein Plan den Befehl, gilt der von dort. Ein Vertrag unter `harness/vertraege/` ist heute nicht der Ort: Dieser Ordner legt Werkzeuge fest, und die Schnittstelle eines Werkzeugs ist sein `make`-Target; für die Marke gibt es keines. Mit einem Sensor wandert die Formel in dessen Vertrag, und §3.14 zeigt dann dorthin.
7. **Der Bestand wird nicht nachgezogen.** Der berichtigte Befehl liefert für jeden lebenden Plan eine Zahl, die nicht über der alten liegt: Ausschlüsse nehmen nur weg, und die niedrigere Schwelle paart nur mehr. Keine Marke greift dadurch früher, also ist kein Plan rückwirkend zu groß und keiner neu zu schneiden. Ein Plan, der die alte Zeile trägt, wird mit dem Befehl aus §3.14 gelesen (Punkt 6); ein Plan, dessen Pfade in den Planungsbaum reichen, wird beim nächsten Anfassen nachgezogen. Die Richtung „nicht über der alten“ gilt nicht für eine Datei, die aus einem ausgeschlossenen Pfad in einen genannten wandert (akzeptiertes Negativ unten); kein lebender Plan liefert das.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (Bestand) | nichts zu ändern | dieselbe Arbeit ergibt 244 oder 1148; ein Slice, dessen Pfade den Planungsbaum fassen, reißt die Marke durch seinen Zustandswechsel |
| B — Messpunkt auf den letzten Inhalts-Commit vor dem Move legen | löst die Doppelzählung ohne Ausschluss | die Bücher zählen weiter (907 gegen 223 am selben Stand); „letzter Inhalts-Commit“ ist selbst eine Lesart; während der Arbeit nicht messbar |
| **C — Endstand, Ausschlüsse und Pin im Befehl, Formel an einer Stelle** | eine Zahl je Arbeit, unabhängig von Basis-Wahl, Zustandswechsel und Konfiguration (223/244 über alle sechs Stände); schreibt die Ausnahme auf, die die Closure schon von Hand anwandte | der Befehl ist länger und trägt zwei Magie-Pfadangaben, die ein Leser kennen muss |
| D — `--no-renames` statt der Schwelle | keine Heuristik, nichts zu paaren | eine bewegte Datei zählt mit ihrer Größe; ein reiner Umzug von Vertragstext reißt Marken ohne Prüfaufwand |
| E — Formel in einen Vertrag unter `harness/vertraege/`, mit Sensor | prüfbar statt Prosa | ein Vertrag ohne Werkzeug und `make`-Target widerspricht dem Ordner; Skript, Fragment, Gegenprobe und Index sind ein eigener Slice, und ein Sensor auf einer falschen Formel ist schlechter als keiner |
| F — Netto-Zeilen oder nur Einfügungen | eine Zahl, kein Flag | misst den Prüfaufwand schlechter (eine ersetzte Zeile ist zwei Zeilen Review); die Marken des Bestands sind in der Summe gesetzt |

## Konsequenzen

- Positiv: Zwei Leser erhalten dieselbe Zahl, und sie steht während der Arbeit zur Verfügung statt erst danach. Über Review, Verifikation, Nacharbeit, Closure und Move bleibt sie stehen (223 bis 244), statt um 241 zu wachsen.
- Positiv: Die Schätzung wird vergleichbar. Die Schätzung des Architect nannte den Plan gesondert mit 60 bis 110 Zeilen; ohne ihn sind es 240 bis 390, und die Messung liegt mit 244 darin. Die 1148 hätten ein Versagen der Schätzung behauptet, das nicht stattfand.
- Negativ: Der Befehl ist nicht mehr auswendig zu schreiben; wer ihn abtippt, lässt leicht einen Ausschluss weg und erhält eine zu große Zahl. Gegenmittel ist Punkt 6: Der Befehl steht an einer Stelle und wird von dort kopiert.
- Akzeptiertes Negativ: Ein Slice, dessen Liefer-Punkte selbst im Planungsbaum liegen (ein Plan, die Roadmap), misst 0. Grund: Der Ausschluss ist eine Pfadangabe und nimmt auch eine ausdrücklich genannte Datei darunter weg. Solche Arbeit begrenzen Liefer-Punkte und Schichten; §4 eines solchen Plans sagt das, statt eine Zahl zu nennen.
- Akzeptiertes Negativ: Eine Datei, die aus einem ausgeschlossenen Pfad in einen genannten wandert, zählt als vollständige Einfügung — die Pfadangabe wirkt vor der Umbenennungs-Erkennung. Grund: selten, und die Richtung ist die vorsichtige.
- Akzeptiertes Negativ: Eine Binärdatei trägt `-` statt Zahlen und zählt 0. Grund: Zeilen sind dort keine Größe; ein Archiv im Pfad ist kein Prüfaufwand in Zeilen.
- Akzeptiertes Negativ: Die Messung setzt voraus, dass die Basis ein Vorfahr des gemessenen Stands ist und der Zweig keine fremde Arbeit aufgenommen hat; wer Arbeit des Hauptzweigs in den Slice-Zweig holt, misst sie mit. Grund: Kein Flag trennt das; der Zweig entsteht am Übergang und endet mit dem Slice.
- Akzeptiertes Negativ: Der Diff-Algorithmus wird nicht festgelegt. Grund, gemessen: zwei Zeilen Unterschied auf 10158 — unter jeder Auflösung, mit der eine Marke gesetzt wird.
- Folgepflicht: Ein Folge-Slice (Kennung vergibt der Planner) trägt den Wortlaut in `AGENTS.md` §3.14 ein, nachdem der Nutzer ihn bestätigt hat, und setzt Status und Index dieser ADR nach. Im selben Commit trägt er die Bindung in den Plan ein, der die Schätz-Hälfte von §3.14 ergänzt (`AGENTS.md` §3.13): Dessen Definition of Done verlangt einen Diff ohne Löschzeile in der Regel selbst und nimmt diese Berichtigung darum nicht an. Kein Sensor; der Zähler bleibt `BEO-REPO/groessenprobe-ohne-messbaren-stand`.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | keine; die Marke prüft heute niemand. Urteilsfrei prüfbar wäre die **Form** — ein Plan in `open/`, `next/` oder `in-progress/`, dessen §4 eine Marke „zu groß“ nennt, nennt Pfade, Marke und Meldeschwelle und zeigt auf §3.14 —, nicht ihr Inhalt. Dieser Sensor entsteht, wenn der Zähler trotz der Berichtigung wieder steigt. | — |

## Re-Evaluierungs-Trigger

Wenn `BEO-REPO/groessenprobe-ohne-messbaren-stand` nach dieser Berichtigung ein weiteres Auftreten zählt (dann ist der Sensor fällig, und die Formel wandert in seinen Vertrag); wenn ein Slice geplant wird, dessen Liefer-Punkte im Planungsbaum liegen (dann braucht der Ausschluss eine benannte Ausnahme); wenn der Nutzer Marke oder Meldeschwelle als Größen ändert.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-21 | Proposed (Empfehlung des Architect; Auftrag des Nutzers „Metrik-Lücke beheben“, Messung an der Closure von `docs/plan/planning/in-progress/slice-erster-release-homebrew-nachweis.md`) | Zweig `arch-metrik` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen entstehen als neue ADR mit `Supersedes ADR-0056`.
