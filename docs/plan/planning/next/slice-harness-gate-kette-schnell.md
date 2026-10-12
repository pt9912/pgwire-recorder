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
Nutzers vom 2026-10-12: erst die Gate-Kette, dann `slice-v1-abschluss-umgebung-anderes-kommando`,
beide vor `slice-v1-abschluss-antwortvergleich` (WIP-Limit 1); Reihenfolge in §4 *Start* und im
Drift-Log der Roadmap. Der Zuschnitt in zwei Slices (dieser und `slice-harness-gate-kette-stop-hook`)
ist **Empfehlung des Planners**, nicht entschieden (§1).

**Bezug:** kein Anspruch des Produkts — der Slice liefert ein Werkzeug des Harness.
Bindung an Entscheidungen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Gate-Kette in
voller und schneller Kette; Accepted 2026-10-12), das sich auf
[`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) stützt (die Gates
tragen die Messmethode der Prüfbarkeit).

**Berührte Spec-Stellen:** `SPEC-052` (neu: dieser Slice legt sie in `spezifikation.md`
§11 *Harness-Werkzeuge* an, der Architect schreibt sie vor dem ersten Code-Commit; die Kennung ist
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
Hashs) und die Randformen aus §6 stehen in der Spezifikation, bevor Code entsteht. Das Ziel der
schnellen Kette schreibt bei grünem Lauf ihren Stand, ersetzt den vollen Nachweis nie, und eine
Gegenprobe fängt jede Zusage der Kette (§6 *Mutationen*).

**Anlass** ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), gemessen am 2026-10-12):
`make gates` braucht rund 4,5 Minuten, davon `lint-gegenprobe` 150 s, `test-integration` 96 s und
`a-check-negativ` 11 s; die übrigen Gates höchstens 2 s. Ein Satz im Plan kostet denselben Lauf wie
eine Änderung am Code.

**Zuschnitt — Empfehlung des Planners an Architect und Nutzer, nicht entschieden.** Der Architect
trug als Liefer-Punkte der Gate-Kette vor: `SPEC-052`, schnelle Kette mit Dateiliste, Hook mit
Gegenprobe. Geprüft nach `AGENTS.md` §3.13: Das sind drei Liefer-Punkte und zwei Schichten, die
Zählung hielte. Geteilt wird trotzdem, weil der Hook das Restrisiko ist ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Konsequenzen*):
Seine Gegenprobe trägt mehr als ein Dutzend Randformen (§6), und Hook samt Gegenprobe soll in
*einer* Review-Sitzung prüfbar sein. **Schätzung des Planners, nicht gemessen:** ungeteilt 800 bis
1100 Zeilen mit Tests, dieser Slice 250 bis 450, `slice-harness-gate-kette-stop-hook` 500 bis 750;
der Architect prüft die Schätzung vor dem Code (`BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`).
Die Teilung des Architect (erst `SPEC-052`, Dateiliste und Hook, danach die schnelle Kette) ist
umgedreht, weil die Bedingung des Hooks ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Entscheidung* Punkt 3) den Stand der schnellen
Kette verlangt: Ein Hook vor der Kette fände nie einen Stand und fiele immer auf `make gates`
zurück, wäre also ein Slice ohne Lieferwert. Dieser Slice ist einzeln lieferbar: Wer ein Dokument
ändert, fährt die schnelle Kette von Hand.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Hook (`.claude/hooks/stop-require-gates.sh`), Dateiliste in `tools/harness/record-gates.sh` und
  `tools/harness/working-tree-hash.sh`, Gegenprobe des Hooks, Gate-Zeile dazu und die Zeile der
  Durchsetzungsschicht in `harness/conventions.md` — Folge-Slice `slice-harness-gate-kette-stop-hook`
  (dort §1 und DoD, im selben Commit wie dieser Plan eingetragen). Der Hook liest den Stand der
  schnellen Kette erst dort; bis dahin verlangt er nach jeder Änderung unverändert `make gates`.
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

- [ ] **Liefer-Punkt 1:** `SPEC-052` steht in `spezifikation.md` §11 (vom Architect, vor dem ersten
      Code-Commit): Gate-Kette (voll, schnell, Zugehörigkeit je Gate), Wahl der Kette (Positivliste,
      Rückfall auf voll), Form des Nachweises (Dateiliste) und jede Randform aus §6 entschieden,
      Punkte nummeriert; §12 trägt die Historie. Die Kennung löst auf, wo
      [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) sie unter `Schärft:` nennt.
- [ ] **Liefer-Punkt 2:** Die schnelle Kette ist ein `make`-Ziel (Name nach `SPEC-052`): Es fährt die
      Gates von `GATE_CHECKS` außer den Gegenproben und `test-integration`; bei grünem Lauf schreibt
      es seinen Stand, bei rotem keinen (ein älterer Stand entfällt); es schreibt und verändert den
      Nachweis des vollen Laufs nie. `make gates` ist unverändert (gleiche Gates, gleiche
      Reihenfolge, gleicher Nachweis).
- [ ] **Liefer-Punkt 3:** Die Gegenprobe der Kette (Temp-Bäume, wie `kopf-check-gegenprobe`) macht
      je Zusage aus §6 *Mutationen* einen Fall, der rot wird, und hängt mit der Kette an
      `GATE_CHECKS`; `harness/README.md` §Sensors führt je eine Zeile für die Kette und ihre
      Gegenprobe mit Bindung an [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md); `harness/conventions.md` nennt in der Zeile zu
      `grundlagen-durchsetzungsschicht.md` eine neue Träger-Datei der Kette, falls der Implementer
      eine anlegt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Die Mutationen und ihre roten Fälle stehen als Tabelle in §7 (Zusage · Mutation · roter
      Fall), vom Implementer selbst gefahren (`AGENTS.md` §3.10).
- [ ] Die gewählten Namen (Ziel der schnellen Kette, Stand-Datei, Gegenprobe) stehen in §1 und §3
      von `slice-harness-gate-kette-stop-hook`, mit Nachzählung seiner Liefer-Punkte und Schichten
      im selben Commit (Closure dieses Slice; `grep` nach ihnen im Nehmer ist nicht leer).
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
| `spec/spezifikation.md` §11 und §12 | update | `SPEC-052` und Historie; **schreibt der Architect vor dem ersten Code-Commit**, der Implementer ändert die Spezifikation nicht (`AGENTS.md` §3.12) |
| `harness/mk/` (neues Fragment) und `Makefile`, falls ein Aggregator nötig ist | neu / update | Ziel der schnellen Kette; die Menge ihrer Gates folgt aus `GATE_CHECKS` und der Regel der Spezifikation (§6, Zugehörigkeit); Namen wählt der Implementer, `SPEC-052` nennt sie |
| `tools/harness/` (neues Skript) | neu | schreibt den Stand der schnellen Kette bei grünem Lauf, in der Form von `SPEC-052`; berührt `record-gates.sh` und `working-tree-hash.sh` nicht (beide gehören `slice-harness-gate-kette-stop-hook`) |
| `tools/harness/` (neue Gegenprobe) und ihr Fragment | neu | je Zusage ein Mutant in einem Temp-Baum, nach dem Vorbild von `tools/harness/kopf-check-gegenprobe.sh`; hängt an `GATE_CHECKS` |
| `harness/README.md` §Sensors | update | Zeile der Kette und Zeile der Gegenprobe, Bindung [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) |
| `harness/conventions.md` (Träger-Tabelle, Zeile `grundlagen-durchsetzungsschicht.md`) | update, falls nötig | nennt die neue Träger-Datei der Kette; die Hook-Seite der Zeile zieht `slice-harness-gate-kette-stop-hook` nach |

- Die Menge der schnellen Kette entsteht entweder als Abzug (die langsamen Gates sind ausdrücklich
  benannt, jedes andere läuft) oder als Positivliste der schnellen Gates; welche Form gilt, ist
  Randform *Zugehörigkeit eines Gates* (§6) und gehört dem Architect.
- Der Implementer prüft vor dem Commit, ob eine weitere lebende Stelle die Gate-Kette beschreibt
  (`grep -rn "make gates"` über `AGENTS.md`, `harness/`, `.claude/commands/`, `.claude/agents/`) und
  zieht sie im selben Commit nach, soweit die Beschreibung falsch wird; die Beschreibung des
  Hooks bleibt bei `slice-harness-gate-kette-stop-hook`.
- Kein Produkt-Code, kein Handbuch, keine Tests unter `internal/` oder `test/`.
## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` ist leer (WIP-Limit 1; `slice-harness-lint-gegenprobe-last`
liegt in `done/`); dieser Slice steht direkt hinter den Slices von welle-v1-abschluss bis
`slice-v1-abschluss-zeitangaben-doku` und vor `slice-harness-gate-kette-stop-hook`,
`slice-v1-abschluss-umgebung-anderes-kommando` und `slice-v1-abschluss-antwortvergleich`
(Drift-Log der Roadmap vom 2026-10-12). Technisch hängt er an nichts. Die Reihenfolge Gate-Kette
vor `slice-v1-abschluss-umgebung-anderes-kommando` vor `slice-v1-abschluss-antwortvergleich` ist
Entscheidung des Nutzers vom 2026-10-12. **Vor dem ersten Code-Commit:** Der Architect hat `SPEC-052`
geschrieben und die Randformen aus §6 entschieden
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`, `AGENTS.md` §3.12); die Tabelle dort ist heute
*offen*. Er prüft dabei die Größenschätzung aus §1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Prüfung des Architect schätzt den Diff
  über eine Review-Sitzung (über 450 Zeilen mit Tests), oder die Gegenprobe der Kette lässt sich
  nicht neben `SPEC-052` in einer Sitzung prüfen. Schnitt dann: `SPEC-052` allein als Slice (der
  Architect schreibt sie ohnehin zuerst), die Kette mit Gegenprobe danach.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer nimmt [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md)
  zurück oder ändert die Positivliste so, dass `SPEC-052` neu zu schreiben ist; oder der Architect
  gibt eine Randform zurück, die eine Entscheidung des Nutzers braucht.
## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

`SPEC-052` mit allen Randformen aus §6 steht; die schnelle Kette und ihre Gegenprobe laufen, `make
gates` ist grün (vom Auftraggeber gefahren), die Tabelle der Mutationen steht in §7, die Namen sind in
`slice-harness-gate-kette-stop-hook` eingetragen, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12): Der Slice liefert einen neuen Vertrag. Jede Randform entscheidet
der Architect **vor dem ersten Code-Commit** in `SPEC-052` (Technik-Stratum, auch für Harness-Werkzeuge);
die Punktnummer trägt er in die Tabelle ein, wenn er die Stelle schreibt. Eine nicht genannte oder
offene Randform entscheidet der Implementer nicht, er gibt sie zurück. Die Spalte *wirkt in* sagt, wo
der Code entsteht; die Tabelle gilt für beide Slices (`slice-harness-gate-kette-stop-hook` nennt die
Hook-Randformen in seinem §6 noch einmal). Stand der Entscheidung heute: **alle offen**; was die ADR
oder der Architect schon vorgibt, steht als Vorgabe dabei.

| Randform | Wirkt in | Vorgabe / Frage | Ort |
|---|---|---|---|
| Zugehörigkeit eines Gates zu den Ketten: neues Gate ohne Eintrag; `a-check-negativ` (Gegenprobe ohne `-gegenprobe` im Namen); `build`, `test` | Kette | [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md): schnell = alles ohne Gegenproben und ohne Integrationstests. Frage: Abzug (nur die langsamen ausdrücklich benannt, ein neues Gate läuft in beiden Ketten) oder Positivliste der schnellen? Empfehlung des Planners: Abzug, weil ein Fehler Zeit kostet, kein Grün | `SPEC-052` |
| Bezugsgröße des Stands der schnellen Kette: Hash des ganzen Arbeitsbaums, Dateiliste oder Stand von `HEAD` plus Diff | Kette, Hook | offen; der Hook vergleicht ihn mit dem Arbeitsbaum | `SPEC-052` |
| Roter Lauf nach einem grünen: bleibt der alte Stand? | Kette | Empfehlung des Planners: Der Lauf entfernt den Stand zu Beginn; rot lässt keinen zurück | `SPEC-052` |
| Form der Dateiliste des vollen Laufs: Inhalt je Eintrag (Pfad, Hash), Symlink mit Ziel, gelöschte Datei, Sortierung, Trennzeichen; Pfad mit Leerzeichen, Zeilenumbruch, nicht UTF-8 | Hook (Schreiben im Folge-Slice) | [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md): Liste der Dateihashes statt eines Hashs; Einzelheiten offen | `SPEC-052` |
| Gelöschte oder umbenannte Datei | Hook | Vorgabe des Architect: beide Pfade gelten als geändert | `SPEC-052` |
| Pfad auf der Positivliste; Pfad auf der Grenze (Präfix `docs/plan/planning-x/`, `docs/reviews-alt/`); Groß- und Kleinschreibung | Hook | Vorgabe: Positivliste zu Beginn nur `docs/plan/planning/` und `docs/reviews/` ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md)); ob Verzeichnisgrenze oder Zeichenpräfix, offen | `SPEC-052` |
| Symlink: Ziel zählt mit; Link in einem inerten Pfad mit Ziel außerhalb | Hook | Vorgabe des Architect: mit Ziel; Frage: inert ist der Link nur, wenn auch sein Ziel inert ist? | `SPEC-052` |
| Leerer Diff gegen den vollen Lauf | Hook | offen (gibt `make gates` frei, weil nichts geändert ist, oder fehlt dann der Stand der schnellen Kette?) | `SPEC-052` |
| Merge-Konflikt (Konfliktmarker, Index mit Stufen) | Hook | offen; Empfehlung des Planners: voll | `SPEC-052` |
| Untracked Datei in und außerhalb der Positivliste | Hook | offen; der bisherige Hash zählt untracked Dateien mit (`git ls-files --others --exclude-standard`) | `SPEC-052` |
| Erster Lauf ohne Stand des letzten vollen grünen Laufs; **dazu der Bestand:** der Hook gibt heute bei fehlendem Stand und sauberem Baum frei (Restlücke im Kopf des Hooks, CI als Netz) | Hook | Vorgabe [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) und Architect: fail-closed, `make gates`. Widerspruch zum Bestand (sauberer Baum, kein Stand → frei): Frage, welche Lesart gilt | `SPEC-052` |
| Nicht lesbarer Vergleich (`git` scheitert, Datei nicht lesbar, Verzeichnis ohne Recht) | Hook | Vorgabe [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md): voll | `SPEC-052` |
| Datei der Dateiliste selbst gelöscht oder beschädigt (abgeschnitten, Zeile ohne Hash) | Hook | Vorgabe des Architect; Empfehlung des Planners: voll | `SPEC-052` |
| Gleichzeitige Läufe (zwei Läufe der Kette, volle und schnelle Kette zugleich, Hook liest beim Schreiben) | Kette, Hook | **Aufruf-Voraussetzung** (`AGENTS.md` §3.12): kommen die Läufe nacheinander oder gleichzeitig? Die Spezifikation bindet jede Zusage an die Antwort | `SPEC-052` |
| Zwischenstand bis `slice-harness-gate-kette-stop-hook`: `SPEC-052` sagt Verhalten des Hooks zu, das erst der Folge-Slice zeigt | Spezifikation | Frage: Wie kennzeichnet die Stelle, dass der Hook die Zusage noch nicht erfüllt (`AGENTS.md` §3.11)? | `SPEC-052` |

**Mutationen** (`AGENTS.md` §3.10): Der Implementer fährt je Zusage die Mutation selbst und trägt die
Tabelle in §7 ein. Zusagen und Mutationen, die der Plan vorgibt (die Fälle der Gegenprobe benennt der
Implementer):

| Zusage | Mutation | roter Fall |
|---|---|---|
| Die schnelle Kette fährt kein langsames Gate | `lint-gegenprobe` oder `test-integration` in die Menge der Kette | Kette fährt es → rot |
| Die schnelle Kette fährt jedes andere Gate | ein schnelles Gate (etwa `kopf-check`) aus der Menge nehmen | Gate fehlt → rot |
| Ein rotes Gate der Kette lässt keinen Stand zurück | rotes Gate in der Kette bei vorhandenem Stand | Stand bleibt → rot |
| Ein grüner Lauf schreibt den Stand, und nur der schnelle Stand | Schreiben unterdrücken; Stand an den Pfad des vollen Nachweises schreiben | kein Stand bzw. voller Nachweis entsteht → rot |
| `make gates` ist unverändert | ein Gate aus `make gates` nehmen oder die Kette in seine Aufrufkette hängen | `make gates` verändert → rot |
| Die Gegenprobe hängt an `GATE_CHECKS` und endet bei Rot mit Exit ungleich 0 | Anschluss entfernen | Fall der Gegenprobe → rot |

**Risiken:**

- Zwischen diesem Slice und `slice-harness-gate-kette-stop-hook` existiert das Ziel der schnellen
  Kette, aber der Hook liest seinen Stand nicht; ein Lauf könnte „grün“ nach der schnellen Kette
  melden, obwohl `make gates` nicht lief — **Ausgang:** entfallen: `SPEC-052` und die Zeile in
  `harness/README.md` sagen, dass die Kette den vollen Nachweis nie ersetzt, die Mutation „voller
  Nachweis entsteht“ prüft es, und der Hook verlangt bis zum Folge-Slice unverändert `make gates`.
- `SPEC-052` sagt Verhalten des Hooks zu, das nur eine Gegenprobe des Folge-Slice prüft
  (`AGENTS.md` §3.11, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) — **Ausgang:** eingetreten:
  `slice-harness-gate-kette-stop-hook` (nimmt an: DoD-Punkt 3 macht je Randform der Hook-Seite einen
  Fall der Gegenprobe).
- Die Namen von Ziel, Stand-Datei und Gegenprobe wählt der Implementer, der Folge-Slice braucht sie
  — **Ausgang:** eingetreten: `slice-harness-gate-kette-stop-hook` (die Closure dieses Slice trägt
  die Namen in dessen §1 und §3 ein, DoD-Punkt oben).
- Die Größe zeigt sich erst in der Schätzung des Architect (§1) — **Ausgang:** weiter offen:
  → `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×; geht der Slice nach der
  Prüfung zurück, ist es das dritte Auftreten und der Eintrag braucht einen Ausgang).
- Der Nutzen schrumpft, wenn die Gegenproben schneller werden ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Re-Evaluierungs-Trigger*)
  — **Ausgang:** entfallen: Der Slice sagt keine Zeitersparnis zu; die Messung liefert
  `slice-harness-gate-kette-stop-hook`, der Re-Evaluierungs-Trigger steht in der ADR.
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

- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (2×, offen) — Gegenstand dieses Slice und seines
  Folge-Slice. Dieser Slice liefert keine Messung der Ersparnis (der Hook liest den Stand erst im
  Folge-Slice); den dritten Beleg und den Ausgang liefert `slice-harness-gate-kette-stop-hook`.
- `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) — berührt: §1 trägt eine
  Schätzung in Zeilen, die Prüfung des Architect vor dem Code ist der Prüfpunkt. Geht der Slice
  zurück, erreicht der Eintrag 3× und braucht einen Ausgang.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — berührt: Die Menge der
  schnellen Kette liest sich als Zusage; die Gegenprobe prüft sie je Gate, damit kein Gate still
  fehlt. Tritt es wieder auf, steht der Eintrag bei 3×.
- `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (1×, offen) — berührt die Hook-Seite (Symlink, nicht
  lesbarer Eintrag, beschädigte Dateiliste), die `slice-harness-gate-kette-stop-hook` liefert; hier
  nur als Randform in §6.
- `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum` (1×, offen) — der Slice folgt der Regel:
  `SPEC-052` entsteht vor dem Code, die ADR zeigt mit `Schärft:` auf sie.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×, offen) — der Slice ist wellenlos; die
  Randformen sind noch nicht entschieden, §4 *Start* verlangt es vor dem ersten Code-Commit.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10),
  `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) und
  `BEO-REPO/slice-waechst-durch-uebernahmen` (verkörpert, `AGENTS.md` §3.13) — neuer Vertrag mit
  Mutationen und Randformen in §6; die Sendung an `slice-harness-gate-kette-stop-hook` ist dort
  nachgezählt.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Zählung nach `AGENTS.md` §3.13:** Liefer-Punkte: 3 (`SPEC-052`; die Kette mit Stand; die Gegenprobe
mit den Gate-Zeilen). Gate-Läufe, Review, Closure-Notiz, Register und die Namen-Eintragung zählen
nicht. Schichten: 2 — Harness-Werkzeug (`tools/harness`, `harness/mk`, `harness/README.md`,
`harness/conventions.md`) und Spezifikation (`spezifikation.md` §11, §12; die Zählregel nimmt sie
nicht aus, deshalb gezählt); keine Nutzer- oder Wartungs-Doku, kein Produkt-Code, kein Bootstrap.
Planung (Pläne, Roadmap, Register) zählt nicht.
**Sendung in `slice-harness-gate-kette-stop-hook`** (Hook, Dateiliste, Gegenprobe des Hooks, die
Namen aus DoD-Punkt „Namen“): Dort zählt §8 mit der Sendung 3 Liefer-Punkte und 1 Schicht; die Namen
gehören zu seinen Liefer-Punkten 1 und 2, sie sind kein vierter.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
