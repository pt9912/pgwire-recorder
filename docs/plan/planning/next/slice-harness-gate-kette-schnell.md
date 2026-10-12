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
Nutzers vom 2026-10-12 (Ziel: ein Release): nach dem Release, hinter `welle-erster-release` und
`welle-v1-nachtrag` und vor `slice-harness-commit-struktur-id` (WIP-Limit 1); davor stehen
`slice-v1-abschluss-umgebung-anderes-kommando` und `slice-v1-abschluss-antwortvergleich`. Reihenfolge in §4
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
ungeteilt 800 bis 1100 Zeilen, dieser Slice 250 bis 450, die Hook-Slices zusammen nach dem Architect 700 bis 950, verteilt 410 bis 590 und 260 bis 390): dieser Slice 400
bis 550 Zeilen mit Tests (`SPEC-052` rund 120 bis 160, Gegenprobe der Kette 200 bis 280, Rest unter 100),
in *einer* Review-Sitzung prüfbar; die Hook-Arbeit lag ungeteilt darüber und ist deshalb geschnitten (§4 von `slice-harness-gate-kette-stop-hook-liste`). Die Zählung nach `AGENTS.md`
§3.13 hält: 3 Liefer-Punkte, 2 Schichten (§8).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Hook (`.claude/hooks/stop-require-gates.sh`), Dateiliste in `tools/harness/record-gates.sh` und
  `tools/harness/working-tree-hash.sh`, Gegenprobe des Hooks, Gate-Zeile dazu und die Zeile der
  Durchsetzungsschicht in `harness/conventions.md` — Folge-Slices `slice-harness-gate-kette-stop-hook-liste` (Dateiliste, Hook, Gegenprobe) und `slice-harness-gate-kette-stop-hook-freigabe` (Positivliste, Freigabe), dort §1 und DoD,
  im selben Commit wie dieser Plan eingetragen. Der Hook liest den Stand der
  schnellen Kette erst in `slice-harness-gate-kette-stop-hook-freigabe`; bis dahin verlangt er nach jeder Änderung unverändert `make gates`.
- Aufnahme weiterer Pfade in die Positivliste (`docs/plan/adr/`, `spec/`, `docs/user/`,
  `README.md`) — Bestand bleibt bewusst stehen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Setzungen des Nutzers*: zu Beginn nur
  `docs/plan/planning/` und `docs/reviews/`; eine Aufnahme ist eine Änderung von `SPEC-052` mit dem
  Beleg, dass kein langsames Gate den Pfad liest, das Benutzerhandbuch ist der heikle Fall (Beispiele
  laufen als Dateien gegen das gebaute Binary, `AGENTS.md` §3.11). Es gibt dafür keinen Slice und
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
      es seinen Stand, bei rotem oder abgebrochenem keinen (ein älterer Stand entfällt); es schreibt und verändert den
      Nachweis des vollen Laufs nie. `make gates` ist unverändert (gleiche Gates, gleiche
      Reihenfolge, gleicher Nachweis).
- [ ] **Liefer-Punkt 3:** Die Gegenprobe der Kette `gates-schnell-gegenprobe` (Temp-Bäume, wie `kopf-check-gegenprobe`) macht
      je Zusage aus §6 *Mutationen* einen Fall, der rot wird, und hängt mit der Kette an
      `GATE_CHECKS`; `harness/README.md` §Sensors führt je eine Zeile für die Kette und ihre
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
| `tools/harness/` (neues Skript) | neu | schreibt den Stand der schnellen Kette bei grünem Lauf, in der Form von `SPEC-052`; berührt `record-gates.sh` und `working-tree-hash.sh` nicht (beide gehören `slice-harness-gate-kette-stop-hook-liste`) |
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

**Start** (`next` → `in-progress`): `in-progress/` ist leer (WIP-Limit 1); `welle-erster-release` und
`welle-v1-nachtrag` sind `done`. Dieser Slice steht direkt hinter diesen Wellen und vor
`slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe` und
`slice-harness-commit-struktur-id` (Drift-Log der Roadmap vom 2026-10-12). Technisch hängt er an nichts. Die
Reihenfolge nach dem Release ist Entscheidung des Nutzers vom 2026-10-12 (Ziel: ein Release). Die offenen Rückfragen an den Architect (Merge-Konflikt zu `slice-harness-gate-kette-stop-hook-liste`, Eingangssatz von `SPEC-052`, neue Mutationen, Rückführungsgrenzen) bleiben offen, bis die Slices der Gate-Kette starten; sie halten das Release nicht auf. **Vor dem ersten Code-Commit:** Der Architect hat die Randformen aus §6
am 2026-10-12 entschieden und die Größe geprüft (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`,
`AGENTS.md` §3.12); die Stelle `SPEC-052` selbst entsteht als erster Commit des Slice (Übertragung von §6,
ohne Code), und erst danach folgt der erste Code-Commit. Ergibt sich beim Übertragen eine Randform, die
§6 nicht nennt, hält der Implementer an und gibt sie dem Architect zurück.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff liegt über einer Review-Sitzung
  (über 600 Zeilen mit Tests; die Schätzung des Architect liegt bei 400 bis 550), oder die Gegenprobe der Kette lässt sich
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
gibt sie zurück. Die Tabelle gilt für alle drei Slices (`slice-harness-gate-kette-stop-hook-liste` und `slice-harness-gate-kette-stop-hook-freigabe` nennen die Hook-Randformen, verteilt vom Planner und vom Architect zu bestätigen, in ihrem §6). Die Spalte *wirkt in* sagt, wo der Code entsteht. Die Zeile *Erster Lauf ohne Stand* ist seit dem 2026-10-12 Entscheidung des Nutzers.

| Randform | Wirkt in | Entscheidung (Architect, 2026-10-12) | Ort |
|---|---|---|---|
| Zugehörigkeit eines Gates zu den Ketten: neues Gate ohne Eintrag; `a-check-negativ`; `build`, `test` | Kette | **Abzug.** Die volle Kette ist `GATE_CHECKS` selbst, ohne zweite Liste (`make gates` unverändert); ein Gate kann aus ihr nicht still fallen, weil kein Fragment sie aufzählt. Die schnelle Kette ist `GATE_CHECKS` ohne die **Ausnahmeliste** im Fragment der Kette: jedes Gate mit `-gegenprobe` im Namen, dazu `a-check-negativ` und `test-integration`. Ein neues Gate läuft deshalb in beiden Ketten, bis es ausdrücklich in die Liste aufgenommen wird; ein Fehler der Liste kostet Zeit, nie ein Gate. `build`, `test` und alle übrigen laufen schnell. Die Gegenprobe der Kette verlangt: jeder Eintrag der Liste steht in `GATE_CHECKS`, und jedes Gate in `GATE_CHECKS` mit `-gegenprobe` im Namen sowie `a-check-negativ` und `test-integration` steht in der Liste (ein neues Gate dieser Namen macht sie rot, bis jemand entscheidet) | `SPEC-052` |
| Bezugsgröße des Stands der schnellen Kette | Kette, Hook | Der Stand der schnellen Kette ist der **Hash des ganzen Arbeitsbaums** zum Ende eines grünen Laufs (Ausgabe von `tools/harness/working-tree-hash.sh`, in diesem Slice unverändert), in `.harness/state/gates-schnell-passed.diffsha`, geschrieben nur von der schnellen Kette. Der Hook vergleicht ihn mit dem Hash des aktuellen Arbeitsbaums; ungleich, fehlend oder nicht lesbar heißt: kein Stand. Der Stand des vollen Laufs ist davon getrennt und eine Dateiliste (Zeile unten); `working-tree-hash.sh` bleibt für den Hook-Slice eine Ausgabe, die genau dann wechselt, wenn sich der Arbeitsbaum ändert | `SPEC-052` |
| Roter Lauf nach einem grünen | Kette | Nach einem roten **oder abgebrochenen** Lauf der schnellen Kette existiert kein Stand; ein älterer entfällt. Zugesagt ist das Ergebnis, nicht der Mechanismus (Entfernen zu Beginn des Laufs ist der Weg, der auch den Abbruch deckt). Der volle Lauf berührt diesen Stand nicht, die schnelle Kette den des vollen Laufs nie | `SPEC-052` |
| Form der Dateiliste des vollen Laufs | Hook (Schreiben in `slice-harness-gate-kette-stop-hook-liste`) | Eine Datei aus **NUL-getrennten Feldern, drei je Eintrag**: Art, Wert, Pfad. Art `F` (reguläre Datei; Wert = SHA-256 des Inhalts, 64 Hexzeichen), `L` (Symlink; Wert = `readlink`-Ziel unverändert), `G` (im Index, im Arbeitsbaum fehlend; Wert leer), `X` (anderer Typ: Verzeichnis, Gitlink; Wert leer). Pfad wie `git ls-files -z` ihn liefert (getrackt und untracked, ohne ignorierte), Sortierung bytegenau, ein Pfad höchstens einmal. Zeilenumbruch, Leerzeichen, Tabulator und Nicht-UTF-8 im Pfad sind damit ohne Sonderfall. Der Ausführungs-Modus gehört nicht zum Eintrag (**akzeptiertes Negativ**: wie der Hash heute; eine Modusänderung ohne Inhaltsänderung bleibt ungeprüft, CI ist das Netz) | `SPEC-052` |
| Gelöschte oder umbenannte Datei | Hook | Verglichen wird je Pfad die Menge beider Listen (Stand des vollen Laufs, aktueller Arbeitsbaum): ein Pfad, dessen Eintrag fehlt, neu ist oder abweicht, gilt als geändert. Eine Umbenennung ist damit ein gelöschter und ein hinzugekommener Pfad, **beide** müssen inert sein | `SPEC-052` |
| Pfad auf der Positivliste; Pfad auf der Grenze; Groß- und Kleinschreibung | Hook | Inert ist ein Pfad, der mit `docs/plan/planning/` oder `docs/reviews/` beginnt, **einschließlich des Schrägstrichs**, bytegenau und case-sensitiv. `docs/plan/planning-x/…`, `docs/reviews-alt/…`, `docs/Plan/planning/…` und eine Datei namens `docs/reviews` sind nicht inert. Ein Pfad mit Segment `.`, `..` oder leerem Segment, ein absoluter Pfad oder ein Pfad mit NUL gilt als beschädigte Liste (nächste Zeilen), nie als inert | `SPEC-052` |
| Symlink | Hook | Inert ist ein geänderter Link nur, wenn **Link-Pfad und Ziel** inert sind; das Ziel wird aus dem Eintrag `L` lexikalisch gegen das Verzeichnis des Links aufgelöst (ohne das Dateisystem zu befragen) und muss danach innerhalb des Repos auf der Positivliste liegen; ein Ziel außerhalb des Repos, ein absolutes Ziel oder ein nicht lexikalisch auflösbares ist nicht inert. Gespiegelt: Liegt irgendwo im Arbeitsbaum ein Link **außerhalb** der Positivliste, dessen Ziel auf sie zeigt, liest ein Gate über ihn inerte Pfade; dann gilt **jede** Änderung als nicht inert. Der Bestand hat Links nur unter `.claude/rules/` mit Zielen außerhalb der Positivliste | `SPEC-052` |
| Leerer Diff gegen den vollen Lauf | Hook | Ist die Liste des Arbeitsbaums mit dem Stand des vollen Laufs identisch, gibt der Hook frei, **ohne** Stand der schnellen Kette: der volle Nachweis deckt den Baum, die schnelle Kette ist eine Teilmenge. Das ist das heutige Verhalten bei gleichem Hash | `SPEC-052` |
| Merge-Konflikt | Hook | Ein nicht aufgelöster Index-Eintrag (`git ls-files -u` nicht leer) ist nie inert: voll. Konfliktmarker im Inhalt einer inerten Datei ändern nichts, ihr Pfad entscheidet | `SPEC-052` |
| Untracked Datei | Hook | Zählt wie eine getrackte (`git ls-files --others --exclude-standard`): in der Positivliste inert, sonst nicht. Eine durch `.gitignore` ausgeblendete Datei ist nicht Teil der Liste, wie heute (Bestand, bewusst stehen gelassen: ein Gate, das sie liest, hat sie auch heute nicht im Hash) | `SPEC-052` |
| Erster Lauf ohne Stand des letzten vollen grünen Laufs; **Widerspruch zum Bestand** | Hook | **Es gilt [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) *Entscheidung* Punkt 3 wörtlich: ein fehlender Stand des vollen Laufs verlangt `make gates`, auch bei sauberem Baum.** Der Grund: Die ADR ist die jüngere und höher rangierende Quelle, nennt keine Ausnahme, und der Bestand (Kopf von `.claude/hooks/stop-require-gates.sh`, „Restlücke … CI ist das Netz“) ist ein Kommentar des Hooks, keine Entscheidung. Das ist eine **Änderung der Zusage des Hooks** (bisher Freigabe, künftig Sperre) und damit des Kommentars im Hook (`AGENTS.md` §3.7); sie fällt `slice-harness-gate-kette-stop-hook-liste` zu (dort §6, DoD, Fall der Gegenprobe). Nebenwirkung: Ein frischer Klon oder Worktree (`.harness/state/` ist ignoriert) sperrt den ersten Stop, bis `make gates` lief. **Entscheidung des Nutzers vom 2026-10-12** („Deine Empfehlungen“, Empfehlung des Architect) | `SPEC-052` |
| Nicht lesbarer Vergleich; Abbruch des Hooks | Hook | Scheitert `git`, ist eine Datei nicht lesbar oder bricht der Hook mitten im Lauf ab, gilt: **voll**, und zwar als Antwort `block`, nicht als Abbruch ohne Antwort (ein Ausgang ungleich 0 ohne Antwort könnte als Freigabe gelesen werden; ob die Umgebung ihn so liest, ist unbekannt). Der Schleifenschutz (`stop_hook_active`) steht davor und bleibt | `SPEC-052` |
| Dateiliste selbst gelöscht oder beschädigt | Hook | Voll. **Beschädigt** ist: leere Datei, Feldzahl nicht durch 3 teilbar, unbekannte Art, `F` mit Wert, der nicht 64 Hexzeichen ist, `G` oder `X` mit Wert, doppelter Pfad, Pfad nach der Zeile oben unzulässig. Eine am Eintragsende abgeschnittene Liste ist nicht beschädigt, aber unschädlich: die fehlenden Pfade gelten als hinzugekommen und entscheiden nach ihrem Pfad | `SPEC-052` |
| Gleichzeitige Läufe (**Aufruf-Voraussetzung**, `AGENTS.md` §3.12) | Kette, Hook | Zugesagt ist **nacheinander**: Änderung des Arbeitsbaums, Lauf der vollen Kette, Lauf der schnellen Kette und Hook kommen nicht gleichzeitig; `SPEC-052` bindet jede Zusage daran und sagt gleichzeitige Läufe nicht zu. Ein Hook, der eine halb geschriebene Datei liest, endet in „kein Stand“ oder „beschädigt“ und damit bei voll; mehr wird nicht zugesagt, ein atomares Schreiben und ein Fall dafür entfallen (Grund: gleiche Folge, kein neuer Fehlerpfad) | `SPEC-052` |
| Zwischenstand bis `slice-harness-gate-kette-stop-hook-liste` und `slice-harness-gate-kette-stop-hook-freigabe` (`AGENTS.md` §3.11) | Spezifikation | `SPEC-052` trennt den Teil über die **Kette** (geprüft von `gates-schnell-gegenprobe`) vom Teil über den **Hook**. Der Hook-Teil beginnt mit dem Satz: *Die Zusagen dieses Teils gelten, soweit das Gate `stop-hook-gegenprobe` sie prüft; solange `make help` kein Ziel dieses Namens führt, ist der Hook-Teil nicht zugesagt, und der Hook verlangt nach jeder Änderung `make gates`.* Die Bedingung ist der Name des Gates, kein Datum; `slice-harness-gate-kette-stop-hook-liste` ändert die Stelle nicht, das Ziel entsteht mit ihr; der Satz bindet jede Zusage an das, was das Gate prüft: bis `slice-harness-gate-kette-stop-hook-freigabe` prüft es nur die Zusagen von `slice-harness-gate-kette-stop-hook-liste`, die Positivliste bleibt bis dahin nicht zugesagt (Planner, vom Architect zu bestätigen) | `SPEC-052` |

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
nicht. Nachgezählt am 2026-10-12 bei der Umordnung (nach dem Release): Die Sendungen an die beiden Nehmer sind
unverändert, die Zählung bleibt 3 Liefer-Punkte und 2 Schichten. Vom Architect am 2026-10-12 nachgeprüft (mit den Entscheidungen aus §6); geschätzt 400 bis 550
Zeilen mit Tests, in einer Review-Sitzung prüfbar.
**Sendungen in `slice-harness-gate-kette-stop-hook-liste`** (Dateiliste, Hook ohne Positivliste, Gegenprobe des
Hooks, die Namen aus §1, die Hook-Randformen aus §6 nach der Verteilung dort) **und
`slice-harness-gate-kette-stop-hook-freigabe`** (Positivliste, Freigabe, Symlink-Regel, Stand der schnellen
Kette, die übrigen Randformen): In beiden zählt §8 mit den Sendungen 3 Liefer-Punkte und 1 Schicht; die Namen
stehen seit dem 2026-10-12 in allen drei Plänen und sind kein vierter Liefer-Punkt. Die Größe der Nehmer prüft
deren §4 (410 bis 590 und 260 bis 390 Zeilen, Verteilung der Schätzung des Architect durch den Planner).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
