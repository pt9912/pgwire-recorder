# Slice slice-harness-gate-kette-schnell: Die schnelle Gate-Kette und `SPEC-052`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt
wird er von der nächsten Welle-Closure. Eingeordnet am 2026-10-12 nach Entscheidung des
Nutzers vom 2026-10-13: vor den übrigen Release-Slices, direkt hinter `slice-v1-abschluss-sessions-doku` und
vor `slice-v1-abschluss-tls-client` (Schritt 28 von welle-v1-abschluss; WIP-Limit 1). Grund: Der Stop-Hook verlangt nach
jeder Änderung einen vollen `make gates` von etwa 5 Minuten. Reihenfolge in §4
*Start* und im Drift-Log der Roadmap. Der Zuschnitt in drei Slices (dieser, `slice-harness-gate-kette-stop-hook-liste` und `slice-harness-gate-kette-stop-hook-freigabe`; die Hook-Arbeit hieß bis zum Schnitt vor dem Start `slice-harness-gate-kette-stop-hook`) und die Reihenfolge dieser zuerst sind **Entscheidung des Nutzers vom 2026-10-12** („Deine
Empfehlungen“, §1).

**Bezug:** kein Anspruch des Produkts — der Slice liefert ein Werkzeug des Harness.
Bindung an Entscheidungen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Gate-Kette in
voller und schneller Kette; Accepted 2026-10-12), das sich auf
[`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) stützt (die Gates
tragen die Messmethode der Prüfbarkeit).

**Berührte Spec-Stellen:** `SPEC-052` (neu: dieser Slice legt sie in `spezifikation.md`
§11 *Harness-Werkzeuge* an, als erster Commit des Slice und vor jedem Code; die Entscheidungen
stehen in §6 dieses Plans, der Implementer überträgt sie und entscheidet nichts; die Kennung ist
in [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) unter `Schärft:` vergeben und steht
bis dahin ohne Anker) · `spezifikation.md` §12 (*Historie*)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `SPEC-052` legt die Gate-Kette fest, und ein `make`-Ziel fährt ihre schnelle Hälfte.
Die volle Kette (`make gates`, unverändert), die schnelle Kette ohne Gegenproben und ohne
Integrationstests, die Positivliste inerter Pfade, die Form des Nachweises (Dateiliste statt eines
Hashs) und die Randformen aus §6 stehen in der Spezifikation, bevor Code entsteht. Das Ziel
`make gates-schnell` schreibt bei grünem Lauf seinen Stand, ersetzt den vollen Nachweis nie, und
`make gates-schnell-gegenprobe` fängt jede Zusage der Kette (§6 *Mutationen*). Die Namen sind
entschieden (Architect, 2026-10-12): `gates-schnell`, `gates-schnell-gegenprobe`, Stand-Datei
`.harness/state/gates-schnell-passed.diffsha`, und für die Folge-Slices `stop-hook-gegenprobe`.

**Anlass** ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), gemessen am 2026-10-12):
`make gates` braucht rund 4,5 Minuten, davon `lint-gegenprobe` 150 s, `test-integration` 96 s und
`a-check-negativ` 11 s; die übrigen Gates höchstens 2 s. Ein Satz im Plan kostet denselben Lauf wie
eine Änderung am Code.

**Zuschnitt — Entscheidung des Nutzers vom 2026-10-12 („Deine Empfehlungen“):** Der Slice
`slice-harness-gate-kette-schnell` liefert `SPEC-052`, die schnelle Kette und die Gegenprobe der
Kette; `slice-harness-gate-kette-stop-hook-liste` liefert Dateiliste, Hook (frei nur bei identischer Liste) und Hook-Gegenprobe, `slice-harness-gate-kette-stop-hook-freigabe` Positivliste und Freigabe. Dieser
Slice kommt zuerst, weil die Bedingung des Hooks ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md),
*Entscheidung* Punkt 3) den Stand der schnellen Kette verlangt: Ein Hook vor der Kette fände nie einen
Stand und fiele immer auf `make gates` zurück, wäre also ein Slice ohne Lieferwert. Dieser Slice ist
einzeln lieferbar: Wer ein Dokument ändert, fährt die schnelle Kette von Hand. Geteilt wird, weil der
Hook das Restrisiko ist (*Konsequenzen* der ADR) und seine Gegenprobe mehr als zwei Dutzend Randformen
trägt. **Größe, geprüft vom Architect am 2026-10-12** (Schätzung, nicht gemessen; der Planner schätzte
ungeteilt 800 bis 1100 Zeilen, dieser Slice 250 bis 450, die Hook-Slices zusammen nach dem Architect 700 bis 950, verteilt 410 bis 590 und 260 bis 390; am 2026-10-13 korrigiert auf 440 bis 620 und 280 bis 420): dieser Slice 400
bis 550 Zeilen mit Tests (`SPEC-052` rund 120 bis 160, Gegenprobe der Kette 200 bis 280, Rest unter 100); **Architect, 2026-10-13:** 420 bis 580, weil `SPEC-052` mit der Positivliste um `docs/user/` und `README.md` samt Beleg und mit der Reihenfolge der Konfliktprüfung um 20 bis 30 Zeilen wächst,
in *einer* Review-Sitzung prüfbar; die Hook-Arbeit lag ungeteilt darüber und ist deshalb geschnitten (§4 von `slice-harness-gate-kette-stop-hook-liste`). **Architect, 2026-10-13 (nach dem Review):** Die Nacharbeit (Prüfung der Make-Flags im Skript, vier bis fünf Fälle, ein Fall für den Stand am Ende des Laufs, Spezifikation netto kürzer) wächst den Diff um 40 bis 80 Zeilen auf 475 bis 515 gegenüber 434; die Rückführungsgrenze von 650 bleibt. Die Zählung nach `AGENTS.md`
§3.13 hält: 3 Liefer-Punkte, 2 Schichten (§8); die Nacharbeit trägt keine Sendung an einen anderen Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Hook (`.claude/hooks/stop-require-gates.sh`), Dateiliste in `tools/harness/record-gates.sh` und
  `tools/harness/working-tree-hash.sh`, Gegenprobe des Hooks, Gate-Zeile dazu und die Zeile der
  Durchsetzungsschicht in `harness/conventions.md` — Folge-Slices `slice-harness-gate-kette-stop-hook-liste` (Dateiliste, Hook, Gegenprobe) und `slice-harness-gate-kette-stop-hook-freigabe` (Positivliste, Freigabe), dort §1 und DoD,
  im selben Commit wie dieser Plan eingetragen. Der Hook liest den Stand der
  schnellen Kette erst in `slice-harness-gate-kette-stop-hook-freigabe`; bis dahin verlangt er nach jeder Änderung unverändert `make gates`.
- Aufnahme weiterer Pfade in die Positivliste über `docs/plan/planning/`, `docs/reviews/`, `docs/user/`
  und `README.md` hinaus (`docs/plan/adr/`, `spec/`) — Bestand bleibt bewusst stehen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md),
  *Setzungen des Nutzers*; eine Aufnahme ist eine Änderung von `SPEC-052` mit dem Beleg, dass kein
  langsames Gate den Pfad liest. Für `docs/user/` und `README.md` ist dieser Beleg geführt (§6, Zeile
  *Beleg für `docs/user/` und `README.md`*; Entscheidung des Nutzers vom 2026-10-13, Beleg des Architect
  vom 2026-10-13); `docs/plan/adr/` und `spec/` haben keinen Beleg. Es gibt dafür keinen Slice und
  keine Kennung.
- Änderung von `make gates` (Zusammensetzung, Reihenfolge, Zusagen) — [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Entscheidung*
  Punkt 1: Die volle Kette bleibt unverändert der einzige Nachweis für Code, Gates, Profile,
  Skripte, `Dockerfile`, Fragmente und `go.mod`; ein Gate zu entfernen wäre eine Lockerung nach
  `AGENTS.md` §3.6.
- Schnellere Gegenproben (weniger Parallelität, Aufteilung, bedingte Läufe) — anderer Vorgang:
  `slice-harness-lint-gegenprobe-last` hat den Standard der Last gesenkt; [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) wählt die
  Wahl der Kette am Pfad, nicht die Beschleunigung einer Gegenprobe.
- Produkt-Code, Benutzerhandbuch, `README.md` — Schicht-Abgrenzung: Der Slice ändert Werkzeuge des
  Harness und die Spezifikation, keine Zusage des Produkts.
## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1:** `SPEC-052` steht in `spezifikation.md` §11, als erster Commit des Slice und vor
      jedem Code, übertragen aus den Entscheidungen von §6 (der Implementer entscheidet nichts; Abweichung
      vom Wortlaut von §6 ist ein Befund): Gate-Kette (voll, schnell, Zugehörigkeit je Gate), Wahl der
      Kette (Positivliste, Rückfall auf voll), Form des Nachweises (Dateiliste) und jede Randform aus §6,
      Punkte nummeriert; der Teil über den Hook trägt die Kennzeichnung aus §6 (letzte Zeile der
      Tabelle); §12 trägt die Historie. Die Kennung löst auf, wo
      [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) sie unter `Schärft:` nennt.
- [ ] **Liefer-Punkt 2:** Die schnelle Kette ist das `make`-Ziel `gates-schnell`: Es fährt `GATE_CHECKS`
      ohne die Ausnahmeliste (Gegenproben und `test-integration`, §6); bei grünem Lauf schreibt
      es seinen Stand, bei rotem oder abgebrochenem keinen (ein älterer Stand entfällt); bei einem Aufruf mit `-n`, `-q`, `-t` oder `-i` in `MAKEFLAGS` schreibt es keinen (§6, Zeile *Make-Flags*); es schreibt und verändert den
      Nachweis des vollen Laufs nie. `make gates` ist unverändert (gleiche Gates, gleiche
      Reihenfolge, gleicher Nachweis).
- [ ] **Liefer-Punkt 3:** Die Gegenprobe der Kette `gates-schnell-gegenprobe` (Temp-Bäume, wie `kopf-check-gegenprobe`) macht
      je Zusage aus §6 *Mutationen* einen Fall, der rot wird, und **die Gegenprobe** hängt an
      `GATE_CHECKS`; `gates-schnell` selbst hängt nicht an `GATE_CHECKS` (Architect, 2026-10-13: Auslegung 1 des Implementers bestätigt, die einzig mögliche Lesart); `harness/README.md` §Sensors führt je eine Zeile für die Kette und ihre
      Gegenprobe mit Bindung an [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md); `harness/conventions.md` nennt in der Zeile zu
      `grundlagen-durchsetzungsschicht.md` eine neue Träger-Datei der Kette, falls der Implementer
      eine anlegt. Die Zeile von `stop-hook-gegenprobe` trägt der Hook-Slice.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Die Mutationen und ihre roten Fälle stehen als Tabelle in §7 (Zusage · Mutation · roter
      Fall), vom Implementer selbst gefahren (`AGENTS.md` §3.10).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, falls eine Beobachtung anfiel; kein
      Zähler wird gesetzt. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/spezifikation.md` §11 und §12 | update | `SPEC-052` und Historie; **erster Commit des Slice, vor jedem Code** (`AGENTS.md` §3.12): der Implementer überträgt den Wortlaut aus §6, ohne zu entscheiden; der Reviewer prüft die Übereinstimmung mit §6 |
| `harness/mk/` (neues Fragment) und `Makefile`, falls ein Aggregator nötig ist | neu / update | Ziel `gates-schnell`; die Menge ihrer Gates ist `GATE_CHECKS` ohne die Ausnahmeliste im Fragment (§6, Zugehörigkeit); Ziel-, Gegenprobe- und Stand-Name stehen fest (§1), Skriptnamen wählt der Implementer |
| `tools/harness/` (neues Skript) | neu | schreibt den Stand der schnellen Kette bei grünem Lauf, in der Form von `SPEC-052`, und prüft vorher die Make-Flags in `MAKEFLAGS` (`SPEC-052` Punkt 6: `-n`, `-q`, `-t`, `-i`); berührt `record-gates.sh` und `working-tree-hash.sh` nicht (beide gehören `slice-harness-gate-kette-stop-hook-liste`) |
| `tools/harness/` (neue Gegenprobe `gates-schnell-gegenprobe`) und ihr Fragment | neu | je Zusage ein Mutant in einem Temp-Baum, nach dem Vorbild von `tools/harness/kopf-check-gegenprobe.sh`; hängt an `GATE_CHECKS` |
| `harness/README.md` §Sensors | update | Zeile der Kette und Zeile der Gegenprobe, Bindung [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) |
| `harness/conventions.md` (Träger-Tabelle, Zeile `grundlagen-durchsetzungsschicht.md`) | update, falls nötig | nennt die neue Träger-Datei der Kette; die Hook-Seite der Zeile ziehen `slice-harness-gate-kette-stop-hook-liste` und `slice-harness-gate-kette-stop-hook-freigabe` nach |

- Die Menge der schnellen Kette ist ein Abzug (Architect, 2026-10-12, §6): die langsamen Gates stehen
  ausdrücklich in der Ausnahmeliste, jedes andere läuft.
- Der Implementer prüft vor dem Commit, ob eine weitere lebende Stelle die Gate-Kette beschreibt
  (`grep -rn "make gates"` über `AGENTS.md`, `harness/`, `.claude/commands/`, `.claude/agents/`) und
  zieht sie im selben Commit nach, soweit die Beschreibung falsch wird; die Beschreibung des
  Hooks bleibt bei `slice-harness-gate-kette-stop-hook-liste` und `slice-harness-gate-kette-stop-hook-freigabe`.
- Kein Produkt-Code, kein Handbuch, keine Tests unter `internal/` oder `test/`.
## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` ist leer (WIP-Limit 1). Dieser Slice steht vor
`slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`,
`slice-v1-abschluss-tls-client` (Schritt 28) und den übrigen Release-Slices (Drift-Log der Roadmap vom 2026-10-13). Technisch hängt er an nichts;
kein Slice von welle-v1-abschluss, welle-erster-release oder welle-v1-nachtrag ist seine Voraussetzung. Die
Reihenfolge ist Entscheidung des Nutzers vom 2026-10-13. Die Rückfragen des Planners (Merge-Konflikt, Eingangssatz von `SPEC-052`, abgeleitete Mutationen, Rückführungsgrenzen) hat der Architect am 2026-10-13 entschieden (§6, Zeilen *Merge-Konflikt* und *Zwischenstand*; Größe §1 und unten). **Vor dem ersten Code-Commit:** Der Architect hat die Randformen aus §6
am 2026-10-12 entschieden und die Größe geprüft (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`,
`AGENTS.md` §3.12); die Stelle `SPEC-052` selbst entsteht als erster Commit des Slice (Übertragung von §6,
ohne Code), und erst danach folgt der erste Code-Commit. Ergibt sich beim Übertragen eine Randform, die
§6 nicht nennt, hält der Implementer an und gibt sie dem Architect zurück.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff liegt über einer Review-Sitzung
  (über 650 Zeilen mit Tests; die Schätzung des Architect liegt bei 420 bis 580; die Grenze steigt mit der Schätzung, weil das Mehr Spezifikationstext ist, den der Reviewer gegen diese Tabelle liest), oder die Gegenprobe der Kette lässt sich
  nicht neben `SPEC-052` in einer Sitzung prüfen. Schnitt dann: `SPEC-052` allein als Slice (sie entsteht
  ohnehin zuerst), die Kette mit Gegenprobe danach.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer nimmt [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md)
  zurück oder ändert die Positivliste so, dass `SPEC-052` neu zu schreiben ist; oder der Implementer
  gibt eine Randform zurück, die eine Entscheidung des Nutzers braucht.
## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

`SPEC-052` mit allen Randformen aus §6 steht; die schnelle Kette und ihre Gegenprobe laufen, `make
gates` ist grün (vom Auftraggeber gefahren), die Tabelle der Mutationen steht in §7, Closure-Notiz mit
Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12): Der Slice liefert einen neuen Vertrag. Jede Randform ist hier vom
Architect am 2026-10-12 **vor dem ersten Code-Commit entschieden**; Ort der Festlegung ist `SPEC-052`
(Technik-Stratum, auch für Harness-Werkzeuge), und der erste Commit des Slice überträgt die Zeilen
dieser Tabelle dorthin, ohne Code (§4 *Start*). Bis dahin ist **diese Tabelle** die Fundstelle; die
Spec-Stelle entsteht mit diesem Slice, ein Eintrag in der Spezifikation vorab wäre eine Stelle ohne
Slice, der sie liefert. Eine nicht genannte oder offene Randform entscheidet der Implementer nicht, er
gibt sie zurück. Die Tabelle gilt für alle drei Slices (`slice-harness-gate-kette-stop-hook-liste` und `slice-harness-gate-kette-stop-hook-freigabe` nennen die Hook-Randformen, verteilt vom Planner und vom Architect zu bestätigen, in ihrem §6). Die Spalte *wirkt in* sagt, wo der Code entsteht. Die Zeile *Erster Lauf ohne Stand* ist seit dem 2026-10-12 Entscheidung des Nutzers. **Architect, 2026-10-13:** Alle Randformen dieser Tabelle sind entschieden, keine ist offen; die Aufruf-Voraussetzung (nacheinander) ist in der Zeile *Gleichzeitige Läufe* gebunden. **Architect, 2026-10-13 (nach dem Review, F-694 bis F-700):** Die Zeilen *Zugehörigkeit*, *Roter Lauf nach einem grünen* und *Beleg …* sind nachgezogen, die Zeile *Make-Flags* und die akzeptierten Negative am Ende der Tabelle sind neu; `SPEC-052` ist entsprechend geändert (vor dem Code übertragen, deshalb änderbar). Auslegung 2 des Implementers ist bestätigt: die Ausnahmeliste hat **sieben** Einträge, auch `gates-schnell-gegenprobe`, und der Beleg unten gilt für sie (ihre Kopien tragen nur Makefile, Fragmente und Skripte der Kette). **Architect, 2026-10-13:** Die Positivliste umfasst von Beginn an auch `docs/user/` und `README.md` (Entscheidung des Nutzers vom 2026-10-13); die Zeilen *Pfad auf der Positivliste …* und *Beleg für `docs/user/` und `README.md` …* tragen sie, `SPEC-052` überträgt sie.

| Randform | Wirkt in | Entscheidung (Architect, 2026-10-12) | Ort |
|---|---|---|---|
| Zugehörigkeit eines Gates zu den Ketten: neues Gate ohne Eintrag; `a-check-negativ`; `build`, `test` | Kette | **Abzug.** Die volle Kette ist `GATE_CHECKS` selbst, ohne zweite Liste (`make gates` unverändert); ein Gate kann aus ihr nicht still fallen, weil kein Fragment sie aufzählt. Die schnelle Kette ist `GATE_CHECKS` ohne die **Ausnahmeliste** im Fragment der Kette: jedes Gate, dessen Name auf `-gegenprobe` endet (Wortlaut an das Muster `*-gegenprobe` der Gegenprobe angeglichen, F-700), dazu `a-check-negativ` und `test-integration`. Ein neues Gate läuft deshalb in beiden Ketten, bis es ausdrücklich in die Liste aufgenommen wird; ein Fehler der Liste kostet Zeit, nie ein Gate der vollen Kette (ein Eintrag zu viel nimmt ein Gate aus der schnellen Kette, die volle fährt es weiter; die Gegenprobe fängt ihn, F-699). `build`, `test` und alle übrigen laufen schnell. Die Gegenprobe der Kette verlangt: jeder Eintrag der Liste steht in `GATE_CHECKS`, und jedes Gate in `GATE_CHECKS`, dessen Name auf `-gegenprobe` endet, sowie `a-check-negativ` und `test-integration` steht in der Liste (ein neues Gate dieser Namen macht sie rot, bis jemand entscheidet) | `SPEC-052` |
| Bezugsgröße des Stands der schnellen Kette | Kette, Hook | Der Stand der schnellen Kette ist der **Hash des ganzen Arbeitsbaums** zum Ende eines grünen Laufs (Ausgabe von `tools/harness/working-tree-hash.sh`, in diesem Slice unverändert; ändert ein Gate des Laufs den Baum, gilt der Hash nach dem Lauf, Fall der Gegenprobe), in `.harness/state/gates-schnell-passed.diffsha`, geschrieben von der schnellen Kette; dass **nur** sie schreibt, ist nicht zugesagt (akzeptiertes Negativ, unten). Der Hook vergleicht ihn mit dem Hash des aktuellen Arbeitsbaums; ungleich, fehlend oder nicht lesbar heißt: kein Stand. Der Stand des vollen Laufs ist davon getrennt und eine Dateiliste (Zeile unten); `working-tree-hash.sh` bleibt für den Hook-Slice eine Ausgabe, die genau dann wechselt, wenn sich der Arbeitsbaum ändert | `SPEC-052` |
| Roter Lauf nach einem grünen | Kette | Nach einem roten **oder abgebrochenen** Lauf der schnellen Kette existiert kein Stand; ein älterer entfällt. Zugesagt ist das Ergebnis, nicht der Mechanismus (Entfernen zu Beginn des Laufs ist der Weg, der auch den Abbruch deckt). Die schnelle Kette berührt den Nachweis des vollen Laufs nie; ob der volle Lauf den Stand der schnellen Kette berührt, ist nicht zugesagt (akzeptiertes Negativ, unten) | `SPEC-052` |
| Form der Dateiliste des vollen Laufs | Hook (Schreiben in `slice-harness-gate-kette-stop-hook-liste`) | Eine Datei aus **NUL-getrennten Feldern, drei je Eintrag**: Art, Wert, Pfad. Art `F` (reguläre Datei; Wert = SHA-256 des Inhalts, 64 Hexzeichen), `L` (Symlink; Wert = `readlink`-Ziel unverändert), `G` (im Index, im Arbeitsbaum fehlend; Wert leer), `X` (anderer Typ: Verzeichnis, Gitlink; Wert leer). Pfad wie `git ls-files -z` ihn liefert (getrackt und untracked, ohne ignorierte), Sortierung bytegenau, ein Pfad höchstens einmal. Zeilenumbruch, Leerzeichen, Tabulator und Nicht-UTF-8 im Pfad sind damit ohne Sonderfall. Sortierung und Einmaligkeit sind Form, keine Zusage mit eigener Mutation (**akzeptiertes Negativ, Architect, 2026-10-13**; Grund: Schreiber und Leser erzeugen die Liste mit derselben Funktion, eine abweichende Sortierung macht aus „identisch“ höchstens „voll“, nie „frei“). Der Ausführungs-Modus gehört nicht zum Eintrag (**akzeptiertes Negativ**: wie der Hash heute; eine Modusänderung ohne Inhaltsänderung bleibt ungeprüft, CI ist das Netz) | `SPEC-052` |
| Gelöschte oder umbenannte Datei | Hook | Verglichen wird je Pfad die Menge beider Listen (Stand des vollen Laufs, aktueller Arbeitsbaum): ein Pfad, dessen Eintrag fehlt, neu ist oder abweicht, gilt als geändert. Eine Umbenennung ist damit ein gelöschter und ein hinzugekommener Pfad, **beide** müssen inert sein | `SPEC-052` |
| Pfad auf der Positivliste; Pfad auf der Grenze; Groß- und Kleinschreibung | Hook | Inert ist ein Pfad, der mit `docs/plan/planning/`, `docs/reviews/` oder `docs/user/` beginnt, **einschließlich des Schrägstrichs**, oder genau `README.md` ist (Wurzel), bytegenau und case-sensitiv. `docs/plan/planning-x/…`, `docs/reviews-alt/…`, `docs/users/…`, `docs/Plan/planning/…`, eine Datei namens `docs/reviews` oder `docs/user`, `readme.md`, `README.md.bak` und `harness/README.md` sind nicht inert (`harness/README.md` bleibt draußen: der Beleg unten gilt nur für die Wurzel-`README.md`; wer sie aufnehmen will, belegt sie). Ein Pfad mit Segment `.`, `..` oder leerem Segment, ein absoluter Pfad oder ein Pfad mit NUL gilt als beschädigte Liste (nächste Zeilen), nie als inert | `SPEC-052` |
| Beleg für `docs/user/` und `README.md` (Aufnahme in die Positivliste); Abdeckungstabellen `docs/user/abdeckung-*.md` | Hook, Spezifikation | **Aufgenommen, ohne Ausnahme** (Entscheidung des Nutzers vom 2026-10-13; Beleg und Abwägung des Architect vom 2026-10-13). Kein langsames Gate (Ausnahmeliste: `lint-gegenprobe`, `kopf-check-gegenprobe`, `abdeckung-gegenprobe`, `hook-gegenprobe`, `gates-schnell-gegenprobe`, `a-check-negativ`, `test-integration`) liest eine Datei unter `docs/user/` oder die `README.md`: die Docker-Gates (`test-integration`, `lint-gegenprobe` über `docker build` einer Kopie des Arbeitsbaums, ebenso `build`, `test`, `lint`) sehen nur den Build-Kontext, und der ist eine Allowlist, die `docs/` und `README.md` nicht freigibt (`.dockerignore`: `*`, dann nur `go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`, `.golangci.yml`, `tools/harness/lint.sh`); ein Go-Test, der dort läse, wäre im Container rot. `a-check-negativ` kopiert den Baum, aber `a-check` liest nach `.a-check.yml` nur `**/*.go`. `kopf-check-gegenprobe`, `abdeckung-gegenprobe`, `hook-gegenprobe` und `gates-schnell-gegenprobe` fahren in Temp-Bäumen (`abdeckung-gegenprobe` baut ihr eigenes `docs/user/`; die Kopien von `gates-schnell-gegenprobe` tragen nur Makefile, Fragmente und Skripte). Die Tabellen `docs/user/abdeckung-*.md` liest nur `abdeckung-check` (`tools/test/abdeckung.sh --check`), und das läuft in der schnellen Kette (kein Name auf `-gegenprobe`): Eine Änderung an ihnen ohne Änderung an Tests (die Deklarationen stehen in `*_test.go`, nicht inert) ist von Hand möglich und wird von `abdeckung-check` rot gemeldet, sodass die schnelle Kette sie fängt; deshalb **nicht** ausgenommen, kein Rückfall auf voll. Das Handbuch und die `README.md` liest von den Gates der schnellen Kette `docs-check` (Verweise, Anker; `.d-check.yml`, Wurzel `.`). Die Aufnahme gilt, solange (a) `docs-check` und `abdeckung-check` in der schnellen Kette stehen (die Gegenprobe der Kette hält die Ausnahmeliste, Mutation 2) und (b) der Build-Kontext `docs/` und `README.md` nicht freigibt; beides sind Änderungen an Gates oder `.dockerignore`, also nicht inert und im vollen Lauf. **Ehrlich:** Die schnelle Kette fährt für eine Handbuch-Änderung kein Gate weniger *über das Handbuch* als die volle (kein langsames Gate liest es), aber sie liefert nicht den Nachweis „`make gates` grün über den ganzen Baum“, und ob ein Beispiel im Handbuch gegen das gebaute Binary läuft (`AGENTS.md` §3.11), prüft **kein** Gate, sondern die Probe von Implementer und Verifier; sie bleibt Pflicht der Rollen und wird durch die schnelle Kette nicht ersetzt. Akzeptiertes Negativ (Grund: gleiche Folge wie in [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *akzeptiertes Negativ*): Ein künftiger Leser von `docs/user/` oder `README.md` in einem Host-Skript eines langsamen Gates ändert Skript oder Gate und löst den vollen Lauf aus; der Beleg ist dann beim Review dieser Änderung neu zu führen, eine eigene Sperre dafür entfällt | `SPEC-052` |
| Symlink | Hook | Inert ist ein geänderter Link nur, wenn **Link-Pfad und Ziel** inert sind; das Ziel wird aus dem Eintrag `L` lexikalisch gegen das Verzeichnis des Links aufgelöst (ohne das Dateisystem zu befragen) und muss danach innerhalb des Repos auf der Positivliste liegen; ein Ziel außerhalb des Repos, ein absolutes Ziel oder ein nicht lexikalisch auflösbares ist nicht inert. Gespiegelt: Liegt irgendwo im Arbeitsbaum ein Link **außerhalb** der Positivliste, dessen Ziel auf sie zeigt, liest ein Gate über ihn inerte Pfade; dann gilt **jede** Änderung als nicht inert. Der Bestand hat Links nur unter `.claude/rules/` mit Zielen außerhalb der Positivliste | `SPEC-052` |
| Leerer Diff gegen den vollen Lauf | Hook | Ist die Liste des Arbeitsbaums mit dem Stand des vollen Laufs identisch **und** gibt es keinen nicht aufgelösten Index-Eintrag (Zeile *Merge-Konflikt*), gibt der Hook frei, **ohne** Stand der schnellen Kette: der volle Nachweis deckt den Baum, die schnelle Kette ist eine Teilmenge. Das ist das heutige Verhalten bei gleichem Hash | `SPEC-052` |
| Merge-Konflikt | Hook | Ein nicht aufgelöster Index-Eintrag (`git ls-files -u` nicht leer) ist nie inert: voll. **Die Prüfung steht vor jeder Entscheidung, auch vor „identisch“** (Architect, 2026-10-13): Die Liste trägt die Index-Stufen nicht, ein Konflikt kann bei identischer Liste bestehen, und ein Hook, der die Identität zuerst prüfte, gäbe ihn frei. Sie ist damit die erste Prüfung des Hooks nach dem Schleifenschutz und gehört zur Liefermenge von `slice-harness-gate-kette-stop-hook-liste`. Konfliktmarker im Inhalt einer inerten Datei ändern nichts, ihr Pfad entscheidet | `SPEC-052` |
| Untracked Datei | Hook | Zählt wie eine getrackte (`git ls-files --others --exclude-standard`): in der Positivliste inert, sonst nicht. Eine durch `.gitignore` ausgeblendete Datei ist nicht Teil der Liste, wie heute (Bestand, bewusst stehen gelassen: ein Gate, das sie liest, hat sie auch heute nicht im Hash) | `SPEC-052` |
| Erster Lauf ohne Stand des letzten vollen grünen Laufs; **Widerspruch zum Bestand** | Hook | **Es gilt [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) *Entscheidung* Punkt 3 wörtlich: ein fehlender Stand des vollen Laufs verlangt `make gates`, auch bei sauberem Baum.** Der Grund: Die ADR ist die jüngere und höher rangierende Quelle, nennt keine Ausnahme, und der Bestand (Kopf von `.claude/hooks/stop-require-gates.sh`, „Restlücke … CI ist das Netz“) ist ein Kommentar des Hooks, keine Entscheidung. Das ist eine **Änderung der Zusage des Hooks** (bisher Freigabe, künftig Sperre) und damit des Kommentars im Hook (`AGENTS.md` §3.7); sie fällt `slice-harness-gate-kette-stop-hook-liste` zu (dort §6, DoD, Fall der Gegenprobe). **Der Kopf-Kommentar des Hooks sagt die Sperre zu** (Zusage im Sinn von `AGENTS.md` §3.7, ersetzt „Restlücke … CI ist das Netz“), und der Fall „fehlender Stand bei sauberem Baum“ der Gegenprobe prüft sie (Architect, 2026-10-13). Nebenwirkung: Ein frischer Klon oder Worktree (`.harness/state/` ist ignoriert) sperrt den ersten Stop, bis `make gates` lief. **Entscheidung des Nutzers vom 2026-10-12** („Deine Empfehlungen“, Empfehlung des Architect) | `SPEC-052` |
| Nicht lesbarer Vergleich; Abbruch des Hooks | Hook | Scheitert `git`, ist eine Datei nicht lesbar oder bricht der Hook mitten im Lauf ab, gilt: **voll**, und zwar als Antwort `block`, nicht als Abbruch ohne Antwort (ein Ausgang ungleich 0 ohne Antwort könnte als Freigabe gelesen werden; ob die Umgebung ihn so liest, ist unbekannt). Der Schleifenschutz (`stop_hook_active`) steht davor und bleibt | `SPEC-052` |
| Dateiliste selbst gelöscht oder beschädigt | Hook | Voll. **Beschädigt** ist: leere Datei, Feldzahl nicht durch 3 teilbar, unbekannte Art, `F` mit Wert, der nicht 64 Hexzeichen ist, `G` oder `X` mit Wert, doppelter Pfad, Pfad nach der Zeile oben unzulässig. Eine am Eintragsende abgeschnittene Liste ist nicht beschädigt, aber unschädlich: die fehlenden Pfade gelten als hinzugekommen und entscheiden nach ihrem Pfad | `SPEC-052` |
| Gleichzeitige Läufe (**Aufruf-Voraussetzung**, `AGENTS.md` §3.12) | Kette, Hook | Zugesagt ist **nacheinander**: Änderung des Arbeitsbaums, Lauf der vollen Kette, Lauf der schnellen Kette und Hook kommen nicht gleichzeitig; `SPEC-052` bindet jede Zusage daran und sagt gleichzeitige Läufe nicht zu. Ein Hook, der eine halb geschriebene Datei liest, endet in „kein Stand“ oder „beschädigt“ und damit bei voll; mehr wird nicht zugesagt, ein atomares Schreiben und ein Fall dafür entfallen (Grund: gleiche Folge, kein neuer Fehlerpfad) | `SPEC-052` |
| Zwischenstand bis `slice-harness-gate-kette-stop-hook-liste` und `slice-harness-gate-kette-stop-hook-freigabe` (`AGENTS.md` §3.11) | Spezifikation | **Architect, 2026-10-13:** Das „soweit“ trägt, eine Zwischenstufe in `SPEC-052` entfällt. `SPEC-052` trennt den Teil über die **Kette** (geprüft von `gates-schnell-gegenprobe`) vom Teil über den **Hook**. Der Hook-Teil beginnt mit dem Satz (Wortlaut für den Implementer): *Die Zusagen dieses Teils gelten, soweit das Gate `stop-hook-gegenprobe` sie prüft, und nur so weit: Ein Punkt, den kein Fall der Gegenprobe nennt, ist nicht zugesagt. Führt `make help` kein Ziel dieses Namens, ist der Hook-Teil insgesamt nicht zugesagt.* Die Bedingung ist der Name des Gates und der Fall, kein Datum; der Satz nennt weder `slice-harness-gate-kette-stop-hook-liste` noch `slice-harness-gate-kette-stop-hook-freigabe` und bleibt nach beiden unverändert wahr, `slice-harness-gate-kette-stop-hook-freigabe` ändert die Spezifikation nicht. Damit „nennt“ prüfbar ist, trägt jeder Fall der Gegenprobe die Nummer des Punktes aus `SPEC-052` in seinem Namen (Vorgabe an beide Hook-Slices, §3 dort). Bis `slice-harness-gate-kette-stop-hook-freigabe` fehlen der Gegenprobe die Fälle der Positivliste; die Freigabe bei geänderten Pfaden ist bis dahin nicht zugesagt, und der Satz sagt es, ohne die Stufe zu benennen. Eine Aussage über das Verhalten des Hooks vor dem Gate (heute „nach jeder Änderung `make gates`“) steht nicht in `SPEC-052`: keine Prüfung deckt sie, sie wäre eine Zusage ohne Test | `SPEC-052` |
| Make-Flags, die kein Gate fahren oder Fehler überlesen (F-694); **Wirkt in:** Kette | Kette, Hook | **Architect, 2026-10-13.** Ein Stand entsteht nur nach einem Lauf, in dem die Gates tatsächlich liefen. `-n`, `-q`, `-t` (das Rezept enthält `$(MAKE)` und läuft deshalb auch mit `-n`, ohne ein Gate zu fahren): kein Stand geschrieben, ein vorhandener bleibt unberührt, es lief ja nichts; `-i` (Fehler überlesen, der Lauf endet mit Exit 0 trotz rotem Gate): kein Stand geschrieben, ein älterer entfällt wie nach Rot. **Zusammen (`-n`, `-q` oder `-t` mit `-i`): `-n`, `-q`, `-t` haben Vorrang, der Stand bleibt unberührt** (Architect, 2026-10-13: Wo nichts lief, gibt es nichts zu entwerten, und ein Trockenlauf ändert den Zustand nie; der Hook liest weiter einen Stand, der zum Baum passt oder am Hash scheitert). Maßgeblich sind die Flags in `MAKEFLAGS` (die Buchstaben im ersten Wort, wenn es nicht mit `-` beginnt), auch aus der Umgebung; die Prüfung steht im Skript `tools/harness/gates-schnell.sh` vor dem Entfernen des Standes. Andere Flags (`-k`, `-j`, `-s`) ändern nichts. Der Hook der Folge-Slices liest die Datei unverändert; er braucht dafür nichts | `SPEC-052` Punkt 6 |
| Akzeptierte Negative der Kette (F-695) | Kette | **Architect, 2026-10-13.** (a) Dass **nur** die schnelle Kette `gates-schnell-passed.diffsha` schreibt, und (b) dass der volle Lauf sie nicht löscht, sind **nicht zugesagt** und stehen nicht in `SPEC-052`. Grund: Löscht der volle Lauf den Stand, verlangt der Hook `make gates` (sichere Richtung); schreibt er ihn, hat er alle Gates der schnellen Kette gefahren (Obermenge). Beide Pfade liegen in `record-gates.sh`, das dieser Slice nicht berührt; eine Gegenprobe müsste `make gates` fahren. Kein Fall, keine Bindung an die Hook-Slices (§3.13): Ändert ein Folge-Slice `record-gates.sh`, bleibt die Aussage wahr | `SPEC-052` Punkte 5 und 7 (nicht zugesagt) |

**Mutationen** (`AGENTS.md` §3.10): Der Implementer fährt je Zusage die Mutation selbst und trägt die
Tabelle in §7 ein. Zusagen und Mutationen, die der Plan vorgibt (die Fälle der Gegenprobe benennt der
Implementer):

| Zusage | Mutation | roter Fall |
|---|---|---|
| Die schnelle Kette fährt kein langsames Gate | `lint-gegenprobe` oder `test-integration` aus der Ausnahmeliste nehmen | Kette fährt es → rot |
| Die schnelle Kette fährt jedes andere Gate | ein schnelles Gate (etwa `kopf-check`) in die Ausnahmeliste nehmen | Gate fehlt → rot |
| Ein neues Gate läuft in beiden Ketten (Abzug, nicht Positivliste) | die schnelle Kette aus einer Positivliste statt aus `GATE_CHECKS` bilden; im Temp-Baum ein zusätzliches Gate ohne Eintrag | Gate fehlt in der schnellen Kette → rot |
| Die Ausnahmeliste deckt sich mit `GATE_CHECKS` | im Temp-Baum ein Eintrag der Liste ohne Gate; ein Gate `x-gegenprobe` ohne Eintrag | beide rot |
| Ein rotes Gate der Kette lässt keinen Stand zurück | rotes Gate in der Kette bei vorhandenem Stand; Lauf mitten im Gate abbrechen | Stand bleibt → rot |
| Ein grüner Lauf schreibt den Stand, und nur den schnellen Stand | Schreiben unterdrücken; Stand an den Pfad des vollen Nachweises schreiben | kein Stand bzw. voller Nachweis entsteht → rot |
| `make gates` ist unverändert und enthält jedes Gate von `GATE_CHECKS` | ein Gate aus `make gates` nehmen oder die Kette in seine Aufrufkette hängen | `make gates` verändert → rot |
| Der Hash wird zum Ende des Laufs gebildet (Punkt 5) | den Hash vor dem Lauf bilden; Fall: ein Gate der Kopie ändert den Arbeitsbaum | Stand gleich Hash vor dem Lauf → rot |
| Ein Lauf mit `-n`, `-q`, `-t` in `MAKEFLAGS` schreibt und entfernt keinen Stand, mit `-i` schreibt er keinen und ein älterer entfällt (Punkt 6) | die Prüfung der Flags im Skript streichen; Fälle in der Kopie: `make -n gates-schnell`, `MAKEFLAGS=n`, `-q`, `-t` je mit und ohne vorhandenen Stand, `-i` mit rotem Gate bei vorhandenem Stand | Stand entsteht, ist weg oder bleibt → rot |
| Die Gegenprobe hängt an `GATE_CHECKS` und endet bei Rot mit Exit ungleich 0 | Anschluss entfernen | Fall der Gegenprobe → rot |

**Risiken:**

- Zwischen diesem Slice und `slice-harness-gate-kette-stop-hook-liste` existiert das Ziel der schnellen
  Kette, aber der Hook liest seinen Stand nicht; ein Lauf könnte „grün“ nach der schnellen Kette
  melden, obwohl `make gates` nicht lief — **Ausgang:** entfallen: `SPEC-052` und die Zeile in
  `harness/README.md` sagen, dass die Kette den vollen Nachweis nie ersetzt, die Mutation „voller
  Nachweis entsteht“ prüft es, und der Hook verlangt bis `slice-harness-gate-kette-stop-hook-freigabe` nach jeder Änderung `make gates`.
- `SPEC-052` sagt Verhalten des Hooks zu, das nur eine Gegenprobe der Folge-Slices prüft
  (`AGENTS.md` §3.11, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) — **Ausgang:** eingetreten:
  `slice-harness-gate-kette-stop-hook-liste` und `slice-harness-gate-kette-stop-hook-freigabe` (nehmen an: DoD-Punkt 3 macht je Randform ihrer Hälfte der Hook-Seite einen
  Fall der Gegenprobe; die Kennzeichnung in `SPEC-052` bindet die Zusage an das Gate `stop-hook-gegenprobe`).
- Die Größe zeigte sich erst in der Prüfung des Architect (§1) — **Ausgang:** weiter offen:
  → `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`; die Prüfung vom 2026-10-12 fand
  die Größe der ungeteilten Hook-Arbeit über der Grenze (heute §4 von `slice-harness-gate-kette-stop-hook-liste`), der Beleg wird bei der Closure von `slice-harness-gate-kette-stop-hook-liste` angelegt.
- Der Nutzen schrumpft, wenn die Gegenproben schneller werden ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Re-Evaluierungs-Trigger*)
  — **Ausgang:** entfallen: Der Slice sagt keine Zeitersparnis zu; die Messung liefert
  `slice-harness-gate-kette-stop-hook-freigabe`, der Re-Evaluierungs-Trigger steht in der ADR.
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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

### Belege des Implementers

**Commits** (Basis `d075d3f`): `cb37a34` (`SPEC-052` und Historie, ohne Code), `f633694` (Fragment
`harness/mk/gates-schnell.mk`, Skripte `tools/harness/gates-schnell.sh` und
`tools/harness/gates-schnell-gegenprobe.sh`, zwei Gate-Zeilen in `harness/README.md`, Träger-Zeile in
`harness/conventions.md`), `a203246` (Belege), `b5cb7b3` (Entscheidungen des Architect nach dem Review: `SPEC-052`, §6),
`0e8fa93` (Nacharbeit: Make-Flags, Fälle der Gegenprobe, Kopf des Fragments, Gate-Zeile), `6e22c4b` (Entscheidung des Architect zu `-n`/`-q`/`-t` mit `-i`),
`68bfa16` (Fall `p6-ni` und Mutant `m-i-vor-n` in der Gegenprobe), dieser Beleg-Commit.
**Größe** (`git diff --shortstat d075d3f HEAD` ohne `docs/`): 563 Zeilen eingefügt, 1 entfernt (`SPEC-052` 157, Gegenprobe 338,
Fragment 23, Skript 42, Doku 3); das sind 129 mehr als die 434 des ersten Stands und mehr als die 40 bis 80 der Schätzung
des Architect (475 bis 515); unter der Rückführungsgrenze von 650. Das Mehr liegt in der Gegenprobe (+99: ein Fall für das
Ende des Laufs, die Make-Flags je mit und ohne Stand, `-n` mit `-i`, dreizehn Mutanten); die Rückführungsgrenze von 650 ist nicht erreicht (87 darunter).

**`SPEC-052`** ist die Übertragung der Tabelle aus §6, nummeriert: Teil 1 (Punkte 1 bis 8) trägt die Zeilen
*Zugehörigkeit*, *Bezugsgröße* (Stand), *Roter Lauf nach einem grünen* und *Gleichzeitige Läufe*; Teil 2
(Punkte 9 bis 20) beginnt mit dem Eingangssatz aus der Zeile *Zwischenstand* im Wortlaut und trägt die
übrigen Zeilen, *Beleg für `docs/user/` und `README.md`* eingeschlossen. Ergänzt habe ich nur die
Rückbezüge zwischen den Punkten (etwa „Punkt 12“ statt „Zeile oben“) und, in Punkt 9, den Satz aus
[ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) *Entscheidung* Punkt 3; Verweise auf ADR und Slice
stehen nicht in der Spezifikation. `make docs-check` nach dem Eintrag und nach den Gate-Zeilen: 0 Befunde.

**Mutationen** (`AGENTS.md` §3.10), am echten Arbeitsbaum gefahren (Mutation per `sed` auf die Datei, dann
`bash tools/harness/gates-schnell-gegenprobe.sh`, danach `git checkout` der Datei); „roter Fall“ ist die
erste Zeile der Ausgabe, die Kopie `ok` ist die unveränderte Kopie mit Attrappen. Mutanten kommen ohne
Docker-Build an: Gegenprobe und Kette laufen auf dem Host und lesen die Dateien direkt.

| Zusage | Mutation (am Arbeitsbaum) | roter Fall |
|---|---|---|
| Die schnelle Kette fährt kein langsames Gate | `lint-gegenprobe` aus `GATE_SCHNELL_AUSNAHME` genommen | Kopie `ok`: `p2-langsam`; Arbeitsbaum: `p4-name` |
| ebenso | `test-integration` aus der Liste genommen | Kopie `ok`: `p2-langsam`; Arbeitsbaum: `p4-name` |
| Die schnelle Kette fährt jedes andere Gate | `kopf-check` in die Liste genommen | Kopie `ok`: `p2-schnell` (und `p6-rot`, weil das rote Gate nicht läuft) |
| Ein neues Gate läuft in beiden Ketten | Kette aus fester Liste statt `GATE_CHECKS` gebildet | Kopie `ok`: `p3-neu` (die Attrappe `neu-x` ohne Eintrag fehlt im Log) |
| Die Ausnahmeliste deckt sich mit `GATE_CHECKS` (Eintrag ohne Gate) | `phantom` in die Liste | Kopie `ok` und Arbeitsbaum: `p4-eintrag` |
| ebenso (Gate ohne Eintrag) | Fragment `harness/mk/zz.mk` mit `GATE_CHECKS += x-gegenprobe`, ohne Eintrag in der Liste | Arbeitsbaum: `p4-name` |
| Ein rotes Gate der Kette lässt keinen Stand zurück | `rm -f "$stand"` im Skript gestrichen | Kopie `ok`: `p6-rot` (und `p6-abbruch`) |
| Ein abgebrochener Lauf lässt keinen Stand zurück | `rm -f` zu Beginn durch `trap … ERR` ersetzt (Entfernen erst im Fehlerfall) | Kopie `ok`: `p6-abbruch`, **nicht** `p6-rot` (der Abbruch beendet die Prozessgruppe, der Trap läuft nicht) |
| Ein grüner Lauf schreibt den Stand | die Zeile, die ihn schreibt, gestrichen | Kopie `ok`: `p5-stand` |
| Nur den schnellen Stand, der volle Nachweis bleibt | Stand an den Pfad `gates-passed.diffsha` geschrieben | Kopie `ok`: `p7-nachweis` (der vorhandene volle Nachweis wird überschrieben) und `p7-nachweis-neu` (er entsteht ohne `make gates`) |
| `make gates` ist unverändert und enthält jedes Gate | `record-gates: $(filter-out kopf-check,$(GATE_CHECKS))` im `Makefile` | Kopie `ok`: `p1-voll` |
| ebenso | `record-gates: $(GATE_CHECKS) gates-schnell` im `Makefile` | Kopie `ok`: `p1-voll` |
| Die Gegenprobe hängt an `GATE_CHECKS` | die Zeile `GATE_CHECKS += gates-schnell-gegenprobe` gestrichen | Kopie `ok`: `anschluss-gegenprobe` (und `p4-eintrag`) |
| Die Gegenprobe endet bei Rot mit Exit ungleich 0 | `-@` vor dem Aufruf der Gegenprobe im Rezept | Kopie `ok`: `anschluss-exit-gegenprobe` |
| Die Kette endet bei Rot mit Exit ungleich 0 | `-@` vor dem Aufruf von `gates-schnell.sh` im Rezept | Kopie `ok`: `anschluss-exit-schnell` (und `p6-rot`) |

**Nacharbeit nach dem Review** (F-694, F-695; derselbe Weg: Mutation per `sed` auf `tools/harness/gates-schnell.sh`, danach
`bash tools/harness/gates-schnell-gegenprobe.sh`, danach die Datei aus einer Sicherung zurück, `touch`; die Skripte laufen auf dem Host,
der Mutant kommt ohne Docker-Build an). Die erste Zeile der Ausgabe ist der rote Fall der Kopie `ok`:

| Zusage (`SPEC-052`) | Mutation (am Arbeitsbaum) | roter Fall |
|---|---|---|
| Punkt 5: der Hash gilt zum Ende des Laufs | Hash vor dem Aufruf von `make` gebildet (ein Gate der Kopie hängt eine Zeile an eine Datei im Baum) | `p5-ende` (nur dieser) |
| Punkt 6: `-n` schreibt keinen Stand (ohne vorhandenen) | die Zeile für `n` gestrichen | `p6-n-schreibt` und `p6-nenv-schreibt` (Aufruf `-n` und `MAKEFLAGS=n` aus der Umgebung) |
| Punkt 6: `-n` entfernt keinen vorhandenen Stand | das Entfernen unbedingt (`if true`) | `p6-n-entfernt`, `p6-nenv-entfernt`, `p6-q-entfernt`, `p6-t-entfernt` |
| Punkt 6: `-n`, `-t` schreiben keinen Stand (trotz Lauf ohne Gate) | die Bedingung „trocken“ vor dem Schreiben gestrichen | `p6-n-schreibt`, `p6-nenv-schreibt`, `p6-t-schreibt` |
| Punkt 6: `-q` entfernt keinen vorhandenen Stand | die Zeile für `q` gestrichen | `p6-q-entfernt` |
| Punkt 6: `-t` schreibt keinen Stand | die Zeile für `t` gestrichen | `p6-t-schreibt` |
| Punkt 6: `-i` schreibt keinen Stand | die Zeile für `i` gestrichen; die Bedingung „überlesen“ vor dem Schreiben gestrichen | `p6-i-entfernt`, `p6-i-schreibt` |
| Punkt 6: bei `-i` entfällt ein älterer Stand | `-i` wie `-n` behandelt (Entfernen übersprungen) | `p6-i-entfernt` (nur dieser; rotes Gate und grüner Lauf bei vorhandenem Stand) |
| Punkt 6: andere Flags (`-k`, `-j`, `-s`) ändern nichts | `k` wie `n` behandelt; ganze `MAKEFLAGS` statt des ersten Worts durchsucht (`--no-print-directory` enthält `n`) | `p6-andere`; bei der zweiten Mutation `p5-stand` und weitere |
| Punkt 6: alle Flags | die Prüfung der Flags ganz gestrichen | `p6-n-schreibt`, `p6-nenv-schreibt`, `p6-q-entfernt`, `p6-t-schreibt`, `p6-i-entfernt`, `p6-i-schreibt` |
| Punkt 7: die Kette löscht den vollen Nachweis nie | `rm -f` des Skripts nennt auch `gates-passed.diffsha` | `p7-nachweis` |

Der Fall `-q` mit Stand hat keinen Mutanten, der nur „schreibt“ bricht (äquivalent): Ohne die Prüfung endet der innere
`make -q` mit Exit 1, bevor das Skript schreibt; die Hälfte „schreibt keinen“ trägt dort `make` selbst, die Hälfte „entfernt
keinen“ trägt `p6-q-entfernt`. Für die Zusage „auch aus der Umgebung“ gibt es einen Fall (`MAKEFLAGS=n`), aber keinen Mutanten,
der nur ihn bricht: make reicht Aufruf-Flags und Umgebung gleich als Buchstaben im ersten Wort durch, das Skript kennt den
Unterschied nicht (gemessen: `make -n` und `MAKEFLAGS=n make` liefern beide `n`; `-n -j2` liefert `n -j2`, `-j2` liefert ` -j2 …`).
Die Zusagen „nur die schnelle Kette schreibt `gates-schnell-passed.diffsha`“ und „der volle Lauf berührt sie nicht“ stehen nicht
mehr in `SPEC-052` (Entscheidung des Architect, §6 *Akzeptierte Negative der Kette*); die Gegenprobe fährt `make gates` nie, nur
`make -n gates`, und ihr Kopf sagt es.

Zusätzlich trägt die Gegenprobe je Zeile der Tabelle einen Mutanten in der Kopie (`m-…`), der die
benannte Prüfung rot färben muss, bei `m-eintrag-ohne-gate` und `m-entfernen-erst-bei-rot` nur sie;
sie sind der Grund, warum die Prüfung nicht aus Versehen leer läuft. Die unveränderte Kopie und der
Arbeitsbaum sind grün (4 s). Kein Mutant blieb grün.

**Läufe der schnellen Kette** (`make gates-schnell`, Docker-Cache warm, Go-Quelltext unverändert):
`docs-check` 1,4 s, `kopf-check` 0,2 s, `baseline-verify` 0,0 s, `abdeckung-check` 1,6 s, `lint` 0,5 s,
`test` 0,5 s, `build` 0,6 s, `a-check` 0,4 s; die Kette zusammen 7,2 bis 7,5 s, Stand
`.harness/state/gates-schnell-passed.diffsha` gleich der Ausgabe von `working-tree-hash.sh`. Mit einer
geänderten Zeile in `cmd/pgwire-recorder/main.go` (Kommentar, danach zurückgesetzt) 102 s, weil `test`,
`lint`, `build` und `a-check` neu laufen; `make gates` braucht rund 270 s ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), gemessen am 2026-10-12).
Ein roter Lauf (ein Kommentar ohne Leerzeile am Ende derselben Datei) beendete die Kette mit Exit 2 und ließ
keinen Stand zurück. Ursache (F-701, nachgefahren: `// probe` an `cmd/pgwire-recorder/main.go` angehängt, `make gates-schnell`,
Datei zurückgesetzt): reproduzierbar, `make test` meldet in der Stufe `test` „nicht gofmt-formatiert: cmd/pgwire-recorder/main.go“
(`gofmt -l` nicht leer); der Kommentar am Dateiende ohne Leerzeile davor ist nicht gofmt-konform. Das Rot kommt vom Gate `test`, nicht
von der Kette, und der Beleg „kein Stand nach Rot“ trägt jetzt auch eine untersuchte Ursache. Die Zeit ist keine Zusage der Spezifikation.
Nach der Nacharbeit: `make -n gates-schnell` und `MAKEFLAGS=n make gates-schnell` lassen den Stand unberührt (gefahren am Arbeitsbaum, Exit 0,
kein Stand entstanden); `make gates-schnell` 7,8 s, Stand gleich der Ausgabe von `working-tree-hash.sh`.
`make gates-schnell-gegenprobe` braucht jetzt rund 17 s (vorher 4 s), weil jede der rund vierzig Kopien nun rund 30 Läufe der Attrappen fährt.
**F-702:** `make gates` fährt die Gegenprobe der Kette mit (Punkt 1 und 4 von `SPEC-052`, `GATE_CHECKS += gates-schnell-gegenprobe`); das ist keine
Lockerung nach `AGENTS.md` §3.6 (ein Gate mehr, keines weniger), die Laufzeit steigt entsprechend (rund 17 s statt 4 s), kein Handlungsbedarf.

**Auslegung, die ich dem Architect zur Bestätigung vorlege** (keine Randform; die Spezifikation und die
Tabelle in §6 lassen je nur diese Lesart zu, §6 ist nicht geändert):

Beide Auslegungen hat der Architect am 2026-10-13 bestätigt (§6, DoD Liefer-Punkt 3; Zeilen *Zugehörigkeit* und *Beleg*; Ausnahmeliste: sieben Einträge, Code, `harness/README.md`, `SPEC-052` Punkt 13 und Gegenprobe zählen gleich).

1. `gates-schnell` hängt **nicht** an `GATE_CHECKS`. Die Zeile *Die Gegenprobe hängt an `GATE_CHECKS` …*
   in §6 und die Zusage *`make gates` ist unverändert* (Mutation „die Kette in seine Aufrufkette hängen“)
   lassen es nicht zu: `gates-schnell` in `GATE_CHECKS` hieße, `make gates` fährt die schnelle Kette ein
   zweites Mal, und `GATE_CHECKS` ohne die Ausnahmeliste enthielte das Ziel selbst. DoD Liefer-Punkt 3
   („hängt mit der Kette an `GATE_CHECKS`“) lese ich als: die Gegenprobe hängt an `GATE_CHECKS`.
2. `gates-schnell-gegenprobe` steht in der Ausnahmeliste (sieben Einträge): Punkt 4 von `SPEC-052` verlangt jedes Gate mit `-gegenprobe`
   im Namen in der Liste, sonst wäre die Gegenprobe der Kette bei ihrem eigenen Namen rot. Sie liest
   weder `docs/user/` noch `README.md` (Kopien tragen nur `Makefile`, Fragmente und Skripte), der Beleg
   gilt damit auch für sie.

**Randformen:** keine zurückgegeben; Teil 2 der Spezifikation (Hook) ist durch den Eingangssatz bis zu
`stop-hook-gegenprobe` nicht zugesagt, `make help` führt kein Ziel dieses Namens.

**Nacharbeit, Klasse der Befunde** (Schritt 19, „als Klasse beheben“): F-694 war „ein Aufruf, in dem kein Gate wirklich lief, schreibt den Stand“; das Merkmal ist
der Make-Flag, der das Rezept leerlaufen lässt oder Fehler überliest. Gefahren ist je Ausprägung ein Fall und eine Mutation (`n`, `q`, `t`, `i`; dazu `-k`, `-j`, `-s`
als Flags, die nichts ändern); die Kombination `-n -i` hat den Fall `p6-ni` (siehe Rückgabe unten). F-695 war „Zusage ohne fangende Mutation“; behoben für
Punkt 5 (Ende des Laufs) und Punkt 7 (löscht nie), die übrigen beiden Zusagen hat der Architect aus `SPEC-052` genommen. F-696: der Satz der Spezifikation nennt 1 bis 7,
Kopf der Gegenprobe, Schlusszeile und Zeile in `harness/README.md` ebenso. F-698: der Kopf von `gates-schnell.mk` sagt, was die Gates brauchen, und dass die
Gegenprobe Docker nicht ruft (kein Fall, `gates-schnell-gegenprobe.sh` ruft es nicht auf). F-699 und F-700: nach der Entscheidung des Architect Wortlaut `-gegenprobe`
im Namen (`SPEC-052` Punkt 2 und 4) und Muster `*-gegenprobe` stehen jetzt gleich; der Überschuss in der Liste wird über `p2-schnell` und `p4-eintrag` gefangen.

**Rückgabe an den Architect, entschieden** (Architect, 2026-10-13, Commit `6e22c4b`; Zeile *Make-Flags* in §6, `SPEC-052` Punkt 6): **`MAKEFLAGS` mit `-n` (oder `-q`, `-t`) und
`-i` zugleich.** Es gilt `-n`/`-q`/`-t`: Der Lauf schreibt keinen Stand, ein vorhandener bleibt unberührt. Das ist die Reihenfolge, die `gates-schnell.sh` schon trägt; der Code ist unverändert.
Fall der Gegenprobe: `p6-ni` (Kopie mit Attrappen, `MAKEFLAGS=ni`, einmal ohne und einmal mit vorhandenem Stand; geprüft werden `p6-ni-schreibt` und `p6-ni-entfernt`,
nach dem Muster der übrigen Flag-Fälle `p6-<name>-…`; der Name im Auftrag lautete `p6-flags-ni`). Mutation (`AGENTS.md` §3.10), selbst gefahren:

| Zusage (`SPEC-052` Punkt 6) | Mutation (Kopie `m-i-vor-n`, im Skript) | roter Fall |
|---|---|---|
| `-n`, `-q`, `-t` haben Vorrang vor `-i`: der Stand bleibt unberührt | `*i*) ueberlesen=1 ;;` zu `*i*) ueberlesen=1; trocken=0 ;;` (`-i` gewinnt, der Stand entfällt) | `p6-ni-entfernt` (nur dieser, `rot … exakt`; kein anderer Fall ist rot, weil `-i` allein und `n`, `q`, `t` allein `trocken` schon anders setzen) |

Die Gegenprobe ist mit dem Mutanten und unverändert grün (`make gates-schnell-gegenprobe`, rund 17 s).

**Stand der Gates:** `make gates` Exit 0 am Baum des Commits dieser Belege (Hash `23d0f621…62baae94`, davor `make docs-check` 0 Befunde); diese Zeile selbst ändert danach nur §7, `make gates` lief mit ihr ein zweites Mal.

**Kommentare** (`AGENTS.md` §3.7, §3.11): neu geschrieben sind der Kopf von `gates-schnell.mk`, von
`gates-schnell.sh` und von `gates-schnell-gegenprobe.sh`; jede Zusage darin hat eine Prüfung der Tabelle
(Ausnahmeliste, Stand, kein Stand nach Rot und Abbruch, voller Nachweis, Anschluss an `GATE_CHECKS`).
Die Grenze der Gegenprobe (Attrappen als Gates, Aufruf-Voraussetzung ohne Fall) steht in ihrem Kopf. Die
Zeile `GEPRÜFT DURCH` im Fragment nennt die Gegenprobe, ohne Punkte zu zählen.

**Weitere Stellen der Gate-Kette** (Plan §3): `grep -rn "make gates"` über `AGENTS.md`, `harness/`,
`.claude/commands/` und `.claude/agents/` fand keine Beschreibung, die durch die schnelle Kette falsch wird
(`make gates` ist unverändert, der Hook bleibt bei `make gates`). Handbuch und `README.md` sind nicht berührt.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/`
gesichtet (Zähler = Dateien unter `evidence/`, gelesen am Stand `f1f4415`, 2026-10-12). Treffer:

- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (2×, offen) — Gegenstand dieses Slice und seiner
  Folge-Slices. Dieser Slice liefert keine Messung der Ersparnis (der Hook liest den Stand erst im
  Folge-Slices); den dritten Beleg und den Ausgang liefert `slice-harness-gate-kette-stop-hook-freigabe`.
- `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) — berührt: §1 trägt eine
  Schätzung in Zeilen, die Prüfung des Architect vor dem Code ist der Prüfpunkt. Geht der Slice
  zurück, erreicht der Eintrag 3× und braucht einen Ausgang.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — berührt: Die Menge der
  schnellen Kette liest sich als Zusage; die Gegenprobe prüft sie je Gate, damit kein Gate still
  fehlt. Tritt es wieder auf, steht der Eintrag bei 3×.
- `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (1×, offen) — berührt die Hook-Seite (Symlink, nicht
  lesbarer Eintrag, beschädigte Dateiliste), die `slice-harness-gate-kette-stop-hook-liste` und `slice-harness-gate-kette-stop-hook-freigabe` liefern; hier
  nur als Randform in §6.
- `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum` (1×, offen) — der Slice folgt der Regel:
  `SPEC-052` entsteht vor dem Code, die ADR zeigt mit `Schärft:` auf sie.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×, offen) — der Slice ist wellenlos; die
  Randformen sind noch nicht entschieden, §4 *Start* verlangt es vor dem ersten Code-Commit.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10),
  `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) und
  `BEO-REPO/slice-waechst-durch-uebernahmen` (verkörpert, `AGENTS.md` §3.13) — neuer Vertrag mit
  Mutationen und Randformen in §6; die Sendungen an `slice-harness-gate-kette-stop-hook-liste` und `slice-harness-gate-kette-stop-hook-freigabe` sind dort
  nachgezählt.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Zählung nach `AGENTS.md` §3.13:** Liefer-Punkte: 3 (`SPEC-052`; die Kette mit Stand; die Gegenprobe
mit den Gate-Zeilen). Gate-Läufe, Review, Closure-Notiz und Register zählen nicht. Schichten: 2 —
Harness-Werkzeug (`tools/harness`, `harness/mk`, `harness/README.md`, `harness/conventions.md`) und
Spezifikation (`spezifikation.md` §11, §12; die Zählregel nimmt sie nicht aus, deshalb gezählt); keine
Nutzer- oder Wartungs-Doku, kein Produkt-Code, kein Bootstrap. Planung (Pläne, Roadmap, Register) zählt
nicht. Nachgezählt am 2026-10-13 bei der Umordnung (vor die Release-Slices): Die Sendungen an die beiden Nehmer sind
unverändert, die Zählung bleibt 3 Liefer-Punkte und 2 Schichten. Vom Architect am 2026-10-12 nachgeprüft (mit den Entscheidungen aus §6), am 2026-10-13 mit der erweiterten Positivliste erneut: 3 Liefer-Punkte, 2 Schichten; geschätzt 420 bis 580
Zeilen mit Tests, in einer Review-Sitzung prüfbar.
**Sendungen in `slice-harness-gate-kette-stop-hook-liste`** (Dateiliste, Hook ohne Positivliste, Gegenprobe des
Hooks, die Namen aus §1, die Hook-Randformen aus §6 nach der Verteilung dort) **und
`slice-harness-gate-kette-stop-hook-freigabe`** (Positivliste, Freigabe, Symlink-Regel, Stand der schnellen
Kette, die übrigen Randformen): In beiden zählt §8 mit den Sendungen 3 Liefer-Punkte und 1 Schicht; die Namen
stehen seit dem 2026-10-12 in allen drei Plänen und sind kein vierter Liefer-Punkt. Die Größe der Nehmer prüft
deren §4 (440 bis 620 und 280 bis 420 Zeilen, Verteilung der Schätzung des Architect durch den Planner).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
