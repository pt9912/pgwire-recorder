# Slice slice-harness-vertraege-umzug-gate-kette: Der Vertrag der Gate-Kette nach `harness/vertraege/`, `spezifikation.md §11` entfällt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er von der nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom 2026-10-13 ([ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md)): Der Umzug der Verträge der Harness-Werkzeuge geschieht direkt hinter `slice-harness-gate-kette-stop-hook-freigabe` und vor `slice-v1-abschluss-tls-client`, geschnitten in drei Slices nach Vertrag (WIP-Limit 1, Reihenfolge: `slice-harness-vertraege-umzug-ordner`, `slice-harness-vertraege-umzug-lint`, `slice-harness-vertraege-umzug-gate-kette`); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode (4)). Bindung an Entscheidungen: [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) (Ort) mit Adaption [MR-002](../../../../harness/conventions.md#mr-002); der Vertrag bleibt von [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) geschärft. Bindung aus `slice-harness-gate-kette-stop-hook-freigabe` (V-343, benannte Spec-Lücke): siehe §1.

**Berührte Spec-Stellen:** `SPEC-052` (Vertrag von `gates-schnell` und Stop-Hook: der Ort wechselt von `spezifikation.md §11` nach `harness/vertraege/`; Punkt 22 (b) wird eindeutig gefasst) · `SPEC-049` (nur die Ortsangabe in `AGENTS.md` §3.12, vom Vorgänger-Slice gestrichen; hier nicht berührt) · `spezifikation.md §11` (*Harness-Werkzeuge*: Einleitung und letzter Vertrag gehen, die Überschrift entfällt) · `spezifikation.md §12` (*Historie*, verliert fünf Zeilen)

**Verantwortlich:** pt9912

**Autor:** pt9912 (Planner). **Datum:** 2026-10-13.

---

## 1. Ziel und Abgrenzung


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Vertrag der Gate-Kette, `SPEC-052` (262 Zeilen, 22 Punkte), steht als Datei in `harness/vertraege/`, im Text unverändert gegenüber `spezifikation.md §11` am Stand
`231089e` bis auf Punkt 22 (b), den der Architect eindeutig fasst; mit ihm verschwindet der Abschnitt *Harness-Werkzeuge* ganz aus der Spezifikation (Einleitung, Überschrift, fünf Zeilen der
*Historie*), und jeder Verweis auf den Vertrag und jeder Plan, der `spezifikation.md §11` als Ort neuer Kennungen nennt, zeigt auf den neuen Ort.

**Sendung aus `slice-harness-vertraege-umzug-ordner` (Zwischenstand mit zwei Orten).** `AGENTS.md` §3.12 nennt für `SPEC-052` (und bis `slice-harness-vertraege-umzug-lint` für `SPEC-049`) `spec/spezifikation.md` §11. Dieser Slice nimmt an: Er streicht die Ausnahme für `SPEC-052` und damit den Zwischenstand mit zwei Orten; danach nennt §3.12 nur `harness/vertraege/` als Ort, ohne Ausnahme und ohne Chronik (`AGENTS.md` §3.7).

**Bindung aus `slice-harness-gate-kette-stop-hook-freigabe` (V-343, benannte Spec-Lücke).** Punkt 22 (b) liest sich einzeln als „leere Umgebung: bei identischer Liste `approve`, bei jedem Unterschied
`block`“; erst der folgende Satz („beide Aufrufarten gelten auch für die Freigabe nach Punkt 9“) hebt das auf. Der Code ist richtig; die Gegenprobe prüft die Freigabe in beiden Aufrufarten
(`p22-freigabe-unterverzeichnis`, `p22-freigabe-leere-umgebung`, Mutanten `m-22-hash-relativ`, `m-22-hash-home`). Dieser Slice ist der Nehmer: Der Architect fasst (b) bei der Neuschrift so,
dass der Satz für sich allein gelesen nicht das Gegenteil sagt; mehr ändert er an Punkt 22 nicht. Das ist die einzige inhaltliche Änderung an einem Vertrag im Umzug (Entscheidung des Nutzers).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ordner, `README.md`, Rang, Reihe der Kennungen und die Verträge von `kopf-check`, `abdeckung` und `commit-msg`: `slice-harness-vertraege-umzug-ordner` (Voraussetzung, er liegt vor diesem Slice in `done/`). Dieser Slice trägt nur den Vertrag, seine Historie-Zeilen und das Entfernen der Einleitung nach.
- Der Vertrag des Lint-Profils und die Ausnahme für `SPEC-049` in `AGENTS.md` §3.12: `slice-harness-vertraege-umzug-lint` (Voraussetzung, liegt davor in `done/`; Sendung aus `slice-harness-vertraege-umzug-ordner` an jenen).
- Jede andere Änderung an Wortlaut oder Zusage von `SPEC-052`: Der Umzug ist wortgleich (Entscheidung des Nutzers); was der Architect darüber hinaus für nötig hält, ist eine Randform in §6 und geht, wenn es mehr als die Bindung ist, zurück an den Planner (Klasse 3: ein anderer Vorgang, Arbeit am Vertrag statt am Ort).
- Die Stop-Hook-Logik, `tools/harness/dateiliste.sh`, `tools/harness/working-tree-hash.sh`, `tools/harness/record-gates.sh`, `tools/harness/gates-schnell.sh`: nur Kommentare mit dem alten Pfad ändern sich (Klasse 4: kein Produkt-Code, keine Logik). Ausnahme nur, wenn der Architect in §6 Randform 2 entscheidet, dass die Selbständerung eines Vertrags in der Gegenprobe einen Fall bekommt.
- Pläne, die einen *bestehenden* Vertrag fortschreiben (`slice-harness-commit-struktur-id`, `slice-harness-abdeckung-gate`, `slice-harness-lint-warten-ohne-frist`): bereits in den ersten beiden Slices nachgezogen.
- Zeitdokumente: Links in Reviews und `done/` nur dort, wo sie sonst brechen; Fließtext bleibt ([ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), akzeptiertes Negativ).

## 2. Definition of Done


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Liefer-Punkt 1 — die Datei samt V-343.** `harness/vertraege/SPEC-052-gate-kette.md` trägt den Text von `SPEC-052` aus `spezifikation.md §11` am Stand `231089e` (Zeilen 3216 bis 3477),
      unverändert bis auf die Abweichungen, die `slice-harness-vertraege-umzug-ordner` für alle Verträge festlegt, und bis auf Punkt 22 (b), der für sich allein gelesen eindeutig ist
      (Bindung aus `slice-harness-gate-kette-stop-hook-freigabe`, V-343; sein Beleg wird in §7 eigens geführt) und, nach Entscheidung des Nutzers vom 2026-10-13 (Randform 2, Variante (b)), bis auf die beiden Stellen von Punkt 12 und den Fall `p12-selbst-vertrag`. Das `README.md` des Ordners führt den Vertrag und
      die fünf Zeilen seiner Historie. Beleg in §7: der Vergleich (Befehl und Ausgabe: der einzige Unterschied außer den Abweichungen ist (b)) und der neue Wortlaut von (b).
- [x] **Liefer-Punkt 2 — die Spezifikation.** `spec/spezifikation.md` führt weder Vertrag noch Einleitung noch Überschrift von §11 noch die fünf Zeilen der Historie; die Abschnitte danach tragen ihre Nummern
      wie die Entscheidung des Architect festlegt (§6 Randform 3); die Spezifikation verweist nicht auf `harness/`; `make docs-check` grün.
- [x] **Liefer-Punkt 3 — die Verweise.** Jeder Verweis auf `SPEC-052` und jede Ortsangabe `spezifikation.md §11` in einem lebenden Plan zeigt auf den neuen Ort: `harness/README.md` (Zeilen von `make gates-schnell`,
      `make gates-schnell-gegenprobe`, `make stop-hook-gegenprobe`), die Kommentare von `.claude/hooks/stop-require-gates.sh`, `tools/harness/gates-schnell.sh`, `tools/harness/gates-schnell-gegenprobe.sh`,
      `tools/harness/stop-hook-gegenprobe.sh`, `tools/harness/record-gates.sh`, `tools/harness/dateiliste.sh`, `harness/mk/gates-schnell.mk`, `harness/mk/stop-hook-gegenprobe.mk`; in den Plänen
      `slice-harness-lh-links-pflicht`, `slice-harness-meldungskatalog-gate`, `slice-harness-coverage`, `slice-harness-mutation` der Ort neuer Kennungen (Kopf, §1, §3, §6; die nächste freie Nummer ist 053);
      `AGENTS.md` §3.12 (Sendung aus `slice-harness-vertraege-umzug-ordner`): die Ausnahme für `SPEC-052` und der Zwischenstand mit zwei Orten entfallen, der Ort ist `harness/vertraege/` ohne Ausnahme; die Rollen- und Command-Texte unter `.claude/`, falls sie `spezifikation.md §11` oder „die Spezifikation“ für Harness-Werkzeuge nennen.
      Beleg in §7: `grep -rn "spezifikation.md#spec-052\|spezifikation.md SPEC-052\|spezifikation.md §11\|spezifikation.md.*§ ?11" --exclude-dir=.git --exclude-dir=.harness .` ohne Treffer außerhalb von Zeitdokumenten und Entscheidungen (Reviews, Pläne in `done/`, ADRs, `MR-002`), und `grep -rn "§ \?12" docs/plan/planning/open docs/plan/planning/next` ohne Verweis auf die Historie der Spezifikation.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Handbuch und README unberührt: Das Produktverhalten ändert sich nicht (`AGENTS.md` §3.11).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)


Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-052-gate-kette.md` | neu | Text aus §11 (Zeilen 3216 bis 3477 am Stand `231089e`); Punkt 22 (b) neu gefasst (V-343) |
| `harness/vertraege/README.md` | update | Eintrag des Vertrags in der Tabelle der Verträge, fünf Zeilen Historie nach der Regel der Einsortierung (Ordner-Slice, Randform 5); die Einleitung bleibt (§6 Randform 4) |
| `spec/spezifikation.md` | update | §11 ganz entfernt (Überschrift, Einleitung, Vertrag), fünf Zeilen der *Historie*; die Folge-Abschnitte (*Historie*) nach Entscheidung des Architect (§6 Randform 3) |
| `harness/README.md` | update | §Sensors: Links und Nennungen in den Zeilen von `make gates-schnell`, `make gates-schnell-gegenprobe`, `make stop-hook-gegenprobe` |
| `.claude/hooks/stop-require-gates.sh`, `tools/harness/gates-schnell.sh`, `tools/harness/gates-schnell-gegenprobe.sh`, `tools/harness/stop-hook-gegenprobe.sh`, `tools/harness/record-gates.sh`, `tools/harness/dateiliste.sh`, `harness/mk/gates-schnell.mk`, `harness/mk/stop-hook-gegenprobe.mk` | update (nur Kommentare) | „Vertrag: spec/spezifikation.md SPEC-052“ nennt den neuen Pfad |
| `tools/harness/stop-hook-gegenprobe.sh` | update (Entscheidung des Nutzers vom 2026-10-13, §6 Randform 2) | ein Fall für die Selbständerung des Vertrags am neuen Pfad (`p12-selbst-vertrag`), mit Mutation; sonst unverändert |
| `AGENTS.md`, `.claude/agents/architect.md` | update, falls sie den alten Ort nennen | Ortsangabe der Spezifikation für Harness-Werkzeuge; `AGENTS.md` §3.12 verliert hier die Ausnahme für `SPEC-052` (Sendung aus `slice-harness-vertraege-umzug-ordner`, Randform 9); Stand zum Planungszeitpunkt: §3.12 ist von `slice-harness-vertraege-umzug-ordner` angepasst, `architect.md` Zeile 35 nennt „die Spezifikation des Gates“ |
| `docs/plan/planning/open/slice-harness-lh-links-pflicht.md`, `slice-harness-mutation.md`; `next/slice-harness-meldungskatalog-gate.md`, `next/slice-harness-coverage.md` | update | Ort neuer Kennungen: `harness/vertraege/` statt `spezifikation.md §11`, nächste freie Nummer 053 (`AGENTS.md` §3.9); die Ortsänderung fügt keinem dieser Pläne einen Liefer-Punkt hinzu und keine Schicht (die Festlegung liegt jetzt in der Schicht Harness, die jeder von ihnen schon berührt) |
| `docs/plan/planning/in-progress/roadmap.md` | update | Stand erst bei der Closure |
| `docs/reviews/2026-10-13-verifikation-slice-harness-gate-kette-stop-hook-freigabe.md` und weitere Reviews mit Link | update (Links) | nur, wo das Doku-Gate sonst bricht |

Ansatz: wie in `slice-harness-vertraege-umzug-ordner`: Der Text wandert per Ausschneiden, der Vergleich in §7 belegt es; Punkt 22 (b) ist die einzige Zeile, die der Architect schreibt. Gegenprobe und Code
bleiben, soweit der Architect keinen Fall für die Selbständerung verlangt (§6 Randform 2); in dem Fall zählt der Fall als Teil von Liefer-Punkt 1 (er prüft die Zusage, die der Vertrag am neuen Ort trägt) und die
Schätzung steigt um rund 60 (Faktor 2 für Gegenprobe-Code, `AGENTS.md` §3.14).

## 4. Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-vertraege-umzug-lint` liegt in `done/` (WIP-Limit 1). Die Entscheidung des Nutzers zu Randform 2 (Variante (b)) liegt seit 2026-10-13 vor; der Start wartet auf nichts weiter. Die übrigen Randformen und der Wortlaut von Punkt 22 (b) (V-343) sind in §6 entschieden; der Implementer überträgt ihn unverändert und gibt eine Form, die §6 nicht nennt, an den Architect zurück.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und Löschungen von
  `git diff --numstat <Basis>..<letzter Commit des Slice> -- spec harness AGENTS.md tools .claude docs/plan docs/reviews` liegt über 900
  (`<Basis>` ist der Commit des Übergangs `next → in-progress` auf dem Hauptzweig; Metrik nach `AGENTS.md` §3.14; Schätzung des Planners 550 bis 750: 262 Zeilen Vertragstext, doppelt gezählt, dazu Historie, Einleitung, Verweise und die Pläne; mit dem Fall aus Randform 2 rund 60 mehr, weil Gegenprobe-Code gerechnet wird).
- `in-progress` → `open` (blockiert — Carveout?): Der Wortlaut aus §6 Randform 1 stellt sich beim Vergleich als Änderung der Zusage heraus (dann entscheidet der Architect, ob V-343 ein eigener Slice wird), oder
  `slice-harness-vertraege-umzug-lint` schließt nicht.

## 5. Closure-Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (inklusive `make stop-hook-gegenprobe`; die Dateien unter `harness/vertraege/` liegen nicht auf der Positivliste des Stop-Hooks), Closure-Notiz mit Lerneintrag. Der Abschnitt *Harness-Werkzeuge* ist aus `spec/spezifikation.md` verschwunden: `MR-002` bleibt aktiv, sein Zustand (Umzug unvollendet) wird in der Closure auf den Stand gebracht (Entscheidung des Planners, `MR-002` ist Adaption, nicht Zeitdokument).

## 6. Risiken und offene Punkte


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Punkt 22 (b) wird neu gefasst, und die Neufassung ändert die Zusage (aus „eindeutig“ wird „anders“) — **Ausgang: entfallen** — Grund: Einleitungssatz des Punktes und die Fälle `p22-freigabe-unterverzeichnis`, `p22-freigabe-leere-umgebung` mit den Mutanten `m-22-hash-relativ`, `m-22-hash-home` blieben unverändert und grün (139 Fälle); der Verifier las (b) für sich allein als Freigabe in beiden Aufrufarten, das Review (F-4, INFO) fand die Zusage unverändert.
- Nach dem Wegfall von §11 bleiben Pläne, die den Abschnitt noch nennen (die vier aus der DoD sind bekannt, ein fünfter ist möglich); `grep` findet nur die Schreibweisen, die er kennt (`BEO-REPO/plan-folgt-korrektur-nicht`) — **Ausgang: entfallen** — Grund: Der `grep` der DoD, die Suche des Reviews und die des Verifiers (je mit eigenen Schreibweisen: `spezifikation.md#spec-052`, `spezifikation.md SPEC-052`, `§11`, `§ 12`) fanden nur die vier bekannten Pläne, `MR-002` (in der Closure berichtigt) und `slice-harness-lint-warten-ohne-frist` (dort ist §11 die umnummerierte *Historie*, richtig); ein fünfter Plan mit der Ortsangabe fand sich nicht.
- Die Kette der Stop-Hooks: Der Vertrag liegt nicht mehr in `spec/`; Punkt 12 sagt „diese Spezifikation“, und der Fall `p12-selbst-spec` prüft `spec/spezifikation.md`, nicht den Vertrag — **Ausgang: entfallen** — Grund: Variante (b) der Randform 2 (Entscheidung des Nutzers) ist geliefert: Punkt 12 nennt den Vertrag, der Fall `p12-selbst-vertrag` und die Mutation „Positivliste um `harness/vertraege/` erweitert“ färben genau diesen Fall rot (Implementer und Verifier selbst gefahren).
- Größe: Der Vertrag ist der größte der fünf; liegt die Summe über der Marke, schneidet der Planner die Verweise ab (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, 2×) — **Ausgang: entfallen** — Grund: Summe 766 im Stand der Verifikation (`8b1ede6..491ef31`), Marke 900 nicht erreicht; kein Schnitt nötig, kein Auftreten der Beobachtung (§7).

**Prüfung des Architect vor dem Code** (2026-10-13, Stand `795f1c0`): Bezug zu [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), [ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) und [MR-002](../../../../harness/conventions.md#mr-002) bestätigt. Nachgemessen am Stand `231089e`: `SPEC-052` 3216 bis 3477 (262 Zeilen, 22 Punkte), Einleitung und Überschrift von §11 2871 bis 2879, fünf Historie-Zeilen 3576 bis 3580. Größe in der Metrik von `AGENTS.md` §3.14: 262 Zeilen doppelt 524, Einleitung neun Zeilen, fünf gelöschte und fünf neue Historie-Zeilen, `harness/README.md` und Kommentare rund 22, Pläne, `AGENTS.md`, `MR-002` und Links rund 50: **680 bis 780** mit dem Fall aus Randform 2 (ohne ihn 620 bis 700) (Fall, Mutation und Wortlaut, Gegenprobe-Code mit rund dem Doppelten gerechnet). Die Marke 900 bleibt; der Abstand ist knapp genug, dass der Implementer bei 800 mitzählt und dem Planner meldet (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`). Zählung bestätigt: 3 Liefer-Punkte (V-343 und die Selbständerung samt Fall gehören zu Punkt 1: Es sind dieselbe Datei und dieselbe Zusage am neuen Ort, kein vierter Gegenstand) und 2 Schichten (Spezifikation; Harness mit `tools/`, `.claude/`, `AGENTS.md`). Kein Schnitt, keine Folge-ADR.

**Randformen** (`AGENTS.md` §3.12): Es gelten die Entscheidungen 1 bis 11 aus §6 von `slice-harness-vertraege-umzug-ordner`. Zusätzlich, vom Architect am 2026-10-13 vor dem ersten Code-Commit, Ort: dieser Plan (der Vertragstext entsteht erst mit Liefer-Punkt 1; der Implementer überträgt den Wortlaut unten unverändert):

1. *Punkt 22 (b) — entschieden (Nutzer hat der Klarstellung zugestimmt, Wortlaut vom Architect).* Der Teil des Punktes ab „(b)“ lautet künftig: „(b) eine leere Umgebung (`env -i`; `PATH` nicht gesetzt, `bash` mit absolutem Pfad gestartet) gibt bei identischer Liste `approve`, bei einem Unterschied, der nach Punkt 9 freigibt (jede geänderte, hinzugekommene oder gelöschte Datei auf der Positivliste und grüner Stand der schnellen Kette), ebenfalls `approve`, bei jedem anderen Unterschied `block`; die Freigabe nach Punkt 9 gilt damit in beiden Aufrufarten, (a) und (b).“ Der Satz „beide Aufrufarten gelten auch für die Freigabe nach Punkt 9 (inerter Unterschied und grüner Stand der schnellen Kette geben `approve`)“ entfällt, sein Inhalt steht in (b). Die Zusage ändert sich nicht: Der Einleitungssatz des Punktes („gibt in einer leeren Umgebung dieselbe Antwort“) und die Fälle `p22-freigabe-unterverzeichnis`, `p22-freigabe-leere-umgebung` mit den Mutanten `m-22-hash-relativ`, `m-22-hash-home` tragen sie schon. Der Rest von Punkt 22 bleibt Zeichen für Zeichen.
2. *Selbständerung (Punkt 12) am neuen Ort — entschieden: (b) (Entscheidung des Nutzers vom 2026-10-13).* Punkt 12 sagt zweimal „diese Spezifikation“ („und diese Spezifikation nennt sie“; „oder dieser Spezifikation ist nicht inert und verlangt den vollen Lauf“), und `p12-selbst-spec` prüft `spec/spezifikation.md`. Am neuen Ort ist „diese Spezifikation“ der Vertrag selbst, `harness/vertraege/SPEC-052-gate-kette.md` steht nicht auf der Positivliste, die Aussage bleibt wahr, aber kein Fall prüft sie für ihren Ort (`AGENTS.md` §3.10, §3.11). *(a)* wortgleich lassen: die Aussage steht ohne Fall am neuen Ort. *(b)* die beiden Stellen auf den Ort fassen, „dieser Vertrag nennt sie“ und „oder diesem Vertrag (`harness/vertraege/SPEC-052-gate-kette.md`) ist nicht inert“, dazu den Fall `p12-selbst-vertrag` (`pfad_fall` mit `block` für diesen Pfad) und die Mutation „Positivliste im Hook um `harness/vertraege/` erweitert“ (Eintrag `aufnahme` wie bei `spec/*`); `p12-selbst-spec` bleibt als Fall der Pfadklasse `spec/`. Der Nutzer hat (b) gewählt, mit `p12-selbst-spec` als bleibendem Fall; die Abweichung von „wortgleich“ ist damit seine Entscheidung. Die Zusage („eine Änderung an der Datei des Vertrags verlangt den vollen Lauf“) bleibt dieselbe, der Wortlaut folgt dem Ort. Die beiden Historie-Zeilen mit dem Datum 2026-10-14 (Entscheidung des Nutzers vom 2026-10-13) ziehen wortgleich um, mit dem Datum.
3. *Nummer der Historie nach dem Wegfall von §11 — entschieden: nachnummerieren.* `## 12. Historie` heißt `## 11. Historie`. Belegt am 2026-10-13: kein Link auf `#12-historie` oder `#11-` im Repo (`grep`), kein Verweis auf §12 in `spec/`; Nennungen „`spezifikation.md` §12“ stehen in `slice-harness-lint-warten-ohne-frist` (Kopf) und im Fließtext von `docs/reviews/2026-10-03-spec-erstfassung.md` (Zeitdokument, bleibt). Der Implementer zieht den Kopf von `slice-harness-lint-warten-ohne-frist` nach und `grep`t vorher `§ \?12` in `docs/plan/planning/open` und `next`.
4. *Dopplung der Einleitung.* Entschieden: Die Einleitung und die Überschrift von §11 entfallen in der Spezifikation; das `README.md` des Ordners bleibt unverändert (Randform 4 des Ordner-Slice).
5. *Kennung für Verträge nach diesem Slice.* Die nächste freie Nummer ist 053; sie steht in den vier Plänen aus Liefer-Punkt 3. Das `README.md` des Ordners nennt die Regel, nicht die Zahl (Randform 6 des Ordner-Slice).
6. *`MR-002`.* Der Zustandsabsatz („Der Umzug geschieht erst im Umzug-Slice …“) lautet danach als Zustand: aktiv, jeder Vertrag eines Harness-Werkzeugs steht in `harness/vertraege/`; keine Chronik (`AGENTS.md` §3.7). Adaption und Auflösungs-Trigger bleiben; der Absatz im Geltungsbereich „bisher `spec/spezifikation.md` §11, künftig …“ wird auf „`harness/vertraege/`“ gekürzt.
7. *Rollen-Datei `.claude/agents/architect.md`.* Entschieden: unverändert. „Die Spezifikation des Gates“ bleibt als Aussage wahr (die Festlegung eines Gates); das Technik-Stratum bleibt (`MR-002` Adaption 2). Die Command-Texte unter `.claude/commands/` nennen §11 nicht (geprüft am 2026-10-13); der `grep` aus Liefer-Punkt 3 belegt es erneut.
8. *Nachzug der vier Pläne.* Der Ort neuer Kennungen heißt dort `harness/vertraege/` und die nächste freie Nummer 053; ein Plan, der „Abschnitt für Harness-Werkzeuge der Spezifikation“ im Fließtext sagt, wird im selben Zug berichtigt. `slice-harness-abdeckung-gate` hat der Ordner-Slice schon nachgezogen.
9. *Ortsangabe in `AGENTS.md` §3.12 (Sendung aus `slice-harness-vertraege-umzug-ordner`, `AGENTS.md` §3.13).* Zusage: Nach diesem Slice nennt §3.12 `harness/vertraege/` als Ort und keine Ausnahme mehr (`spec/spezifikation.md` §11 kommt dort nicht mehr vor, auch kein „bis zu deren Umzug“). Beleg: `grep -n "spezifikation.md. §11\|§11" AGENTS.md` zeigt keine Zeile, die §11 als Ort für Harness-Werkzeuge nennt. Zeile in einer bereits berührten Datei und Schicht (Harness), derselbe Liefer-Punkt 3.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Die Belege aus §2 (Vergleich gegen den Stand `231089e`, Wortlaut von Punkt 22 (b), `grep` der alten Pfade und Orte) stehen hier, nicht im Bericht.

### Belege des Implementers

Basis der Größenprobe: `8b1ede6` (Move `next → in-progress`). Commits: `90c1b37` (Vertrag, Spezifikation, Historie, Fall `p12-selbst-vertrag`), `ebccbeb` (Verweise, Kommentare, `AGENTS.md` §3.12, vier Pläne, Kopf von `slice-harness-lint-warten-ohne-frist`).

**Läufe.** `make gates` auf `ebccbeb`: Exit 0, 7 min 27 s, `lint-gegenprobe: gruen — … jeder Fall je Punkt von SPEC-049 wie erwartet`, `stop-hook-gegenprobe: gruen — 139 Fälle der Punkte 1, 9 bis 22 je Mutant rot, unveränderte Kopie grün, Anschluss an GATE_CHECKS` (138 vorher, ein Fall mehr: `p12-selbst-vertrag`). `make docs-check` auf demselben Stand: `718 Datei(en) geprüft, 0 Befund(e)`. Die volle `stop-hook-gegenprobe` lief damit im Gate und einzeln (siehe Mutation 2).

**Liefer-Punkt 1 — wortgleich bis auf die drei benannten Stellen.**

```
git show 231089e:spec/spezifikation.md | sed -n '3216,3477p' | sed '1s/^### /# /;$d' | diff - harness/vertraege/SPEC-052-gate-kette.md
```

Ausgabe: genau zwei Hunks, `120,122c120,122` (Punkt 12: „und dieser Vertrag nennt sie“, „oder diesem Vertrag (`harness/vertraege/SPEC-052-gate-kette.md`) ist nicht inert“) und `247,249c247,250` (Punkt 22 (b)); sonst leer (Abweichung nur `#` statt `###` in der Überschrift, keine Leerzeile am Ende; 262 Zeilen). Punkt 22 (b), neuer Wortlaut: „(b) eine leere Umgebung (`env -i`; `PATH` nicht gesetzt, `bash` mit absolutem Pfad gestartet) gibt bei identischer Liste `approve`, bei einem Unterschied, der nach Punkt 9 freigibt (jede geänderte, hinzugekommene oder gelöschte Datei auf der Positivliste und grüner Stand der schnellen Kette), ebenfalls `approve`, bei jedem anderen Unterschied `block`; die Freigabe nach Punkt 9 gilt damit in beiden Aufrufarten, (a) und (b).“ Der Satz „beide Aufrufarten gelten auch für die Freigabe nach Punkt 9 …“ entfällt; der Einleitungssatz des Punktes und der Rest bleiben Zeichen für Zeichen (Hunk auf drei Zeilen begrenzt). Umbruch: Die neuen Zeilen 122 und 250 sind länger als die Nachbarzeilen, weil der Wortlaut des Architect unverändert übernommen und nur an den Stellen der Hunks umbrochen ist. Die fünf Historie-Zeilen: `diff <(tail -5 harness/vertraege/README.md) <(sed -n '3576,3580p' <Stand 231089e>)` ohne Ausgabe (wortgleich, am Ende der Tabelle, nach der Zeile vom 2026-10-08; die zwei Zeilen vom 2026-10-14 mit Datum). Der Vertragstext nennt keinen Pfad der Spezifikation: `grep -n "spec/\|spezifikation\|](" harness/vertraege/SPEC-052-gate-kette.md` ohne Treffer.

**Liefer-Punkt 2 — Spezifikation.** `spec/spezifikation.md`: 3233 Zeilen am Stand `8b1ede6`, 2957 jetzt; `git diff 8b1ede6 HEAD --stat -- spec`: eine Einfügung (die Überschrift `## 11. Historie`), 277 Löschungen (Überschrift, Einleitung und Vertrag von §11 samt Leerzeile, fünf Historie-Zeilen). `grep -c "SPEC-052\|^## 11\. Harness\|^## 12\. Historie"` ergibt vorher 8, nachher 0. `grep -n "§ \?12\|Harness-Werkzeug\|SPEC-052\|harness/" spec/*.md` ohne Treffer (Einleitung und Überschrift entfallen, kein Verweis auf `harness/`). `## 11. Historie` ist die letzte Überschrift; kein Link auf `#12-historie` im Repo.

**Liefer-Punkt 3 — Verweise.** Der `grep` aus der DoD (`spezifikation.md#spec-052`, `spezifikation.md SPEC-052`, `spezifikation.md §11`, `spezifikation.md.*§ ?11`, ohne `.git` und `.harness`) trifft außerhalb von Zeitdokumenten (Reviews, `done/`, ADRs, Register) nur diesen Plan (Zitate), `MR-002` (Zustandsabsatz und Geltungsbereich; Randform 6, vom Planner in der Closure auf den Stand gebracht, §5) und `open/slice-harness-lint-warten-ohne-frist.md:26,168` (dort ist `spezifikation.md` §11 die umnummerierte *Historie*, nach Randform 3 richtig). `grep -rn "§ \?12" docs/plan/planning/open docs/plan/planning/next`: ohne Treffer (vor dem Nachzug zwei Treffer, beide `slice-harness-lint-warten-ohne-frist`, jetzt §11). `grep -n "§11\|§ 11" AGENTS.md`: keine Zeile. Nachgezogen: `harness/README.md` Zeile von `gates-schnell` (Link auf `vertraege/SPEC-052-gate-kette.md`; die Zeilen von `gates-schnell-gegenprobe` und `stop-hook-gegenprobe` nennen `SPEC-052` ohne Link und bleiben); acht Dateien mit Ortsangabe im Kommentar (`.claude/hooks/stop-require-gates.sh`, `tools/harness/gates-schnell.sh`, `gates-schnell-gegenprobe.sh`, `stop-hook-gegenprobe.sh`, `record-gates.sh`, `dateiliste.sh`, `harness/mk/gates-schnell.mk`, `harness/mk/stop-hook-gegenprobe.mk`; die Kommentare mit dem Pfad sind `grep`-belegt, kein Gate fängt sie); `AGENTS.md` §3.12 (Ausnahme und Zwischenstand gestrichen, „in `harness/vertraege/` (Ort seit slice-harness-vertraege-umzug-gate-kette)“); `slice-harness-mutation`, `slice-harness-coverage`, `slice-harness-lh-links-pflicht`, `slice-harness-meldungskatalog-gate` (Kopf, §1, §3, §6; nächste freie Nummer 053); Kopf und §3 von `slice-harness-lint-warten-ohne-frist` (§12 → §11). Keine Datei in `docs/reviews/` und keine in `done/` brauchte einen Nachzug (`make docs-check` 0 Befunde; Links dorthin auf `SPEC-052` gab es nicht). `.claude/agents/architect.md` unverändert (Randform 7), die Command-Texte nennen §11 nicht (`grep -rn "§11\|§ 11\|SPEC-052" .claude` trifft nur den Kommentar des Hooks).

**Zusage · Mutation · roter Test** (einzeln gefahren, danach `git checkout`; kein Docker-Build-Kontext betroffen, daher kein `touch` nötig außer beim Hook):

| # | Zusage | Mutation | roter Lauf |
|---|---|---|---|
| 1 | der Text ist wortgleich | ein `X` vor Zeile 3 von `SPEC-052-gate-kette.md` | der Vergleich oben meldet `3c3` mit beiden Zeilen |
| 2 | Punkt 12 (neu): eine Änderung an der Datei des Vertrags verlangt den vollen Lauf (Fall `p12-selbst-vertrag`, Pfad `harness/vertraege/SPEC-052-gate-kette.md`) | Positivliste im Hook um `harness/vertraege/*) return 0 ;;` erweitert (`touch` auf den Hook) | `make stop-hook-gegenprobe`: `ROT — Kopie 'ok' ist nicht grün: p12-selbst-vertrag`, nur dieser Fall; zusätzlich im grünen Lauf `aufnahme 'harness/vertraege/*' p12-selbst-vertrag` (Mutant färbt den Fall rot, sonst meldet die Gegenprobe „färbt nicht rot“) |
| 3 | `p12-selbst-spec` bleibt als Fall der Pfadklasse `spec/` | unverändert; `aufnahme 'spec/*' p12-selbst-spec` steht weiter in der Gegenprobe, 139 Fälle grün | grüner Lauf `ebccbeb` (die Gegenprobe verlangt, dass der Mutant den Fall rot färbt) |
| 4 | Punkt 22 (b) sagt für sich allein gelesen nicht das Gegenteil (Freigabe in beiden Aufrufarten) | der alte Satz „beide Aufrufarten gelten auch für die Freigabe nach Punkt 9“ wieder eingesetzt | `grep -c "beide Aufrufarten gelten auch" harness/vertraege/SPEC-052-gate-kette.md` meldet 1 statt 0 (ein `grep`, kein Gate; die Zusage selbst tragen `p22-freigabe-unterverzeichnis`, `p22-freigabe-leere-umgebung` mit `m-22-hash-relativ`, `m-22-hash-home`, unverändert grün) |
| 5 | die Spezifikation führt keinen Abschnitt *Harness-Werkzeuge* und keine Kennung `SPEC-052` mehr | Zeile `## 11. Harness-Werkzeuge` am Ende | `grep -c "^## 11. Harness-Werkzeuge" spec/spezifikation.md` meldet 1 statt 0 (ein `grep`) |
| 6 | die Spezifikation verweist nicht auf `harness/` | Link auf `../harness/vertraege/README.md` am Ende | `make docs-check`: `matrix-forbidden — Referenz spec-straten → aussen ist nicht erlaubt`, 1 Befund |
| 7 | die Historie der Spezifikation heißt `## 11. Historie` | `## 12. Historie` | `grep -c "^## 12. Historie" spec/spezifikation.md` meldet 1 statt 0 (ein `grep`) |
| 8 | jeder Link auf den Vertrag zeigt auf den neuen Ort | `harness/README.md:78` zurück auf `../spec/spezifikation.md#spec-052…` | `make docs-check`: 1 Befund, `anchor-missing` |
| 9 | `AGENTS.md` §3.12 nennt keine Ausnahme für `SPEC-052` | „, für den Vertrag `SPEC-052` in `spec/spezifikation.md` §11“ wieder eingesetzt | `grep -n "spezifikation.md. §11" AGENTS.md` meldet Zeile 305 (ein `grep`) |
| 10 | die fünf Historie-Zeilen stehen im `README.md` des Ordners | die letzte Zeile gelöscht | `grep -c "Harness-Werkzeuge:" harness/vertraege/README.md` meldet 7 statt 8 |

Rang- und Ortsaussagen in Prosa haben kein Gate; die `grep`-Zeilen sind der Beleg, keine Sperre. Bei Mutation 2 färbt der Mutant genau einen Fall, nicht den Aufbau.

**Größe** (Metrik `AGENTS.md` §3.14, `git diff --numstat 8b1ede6..ebccbeb -- spec harness AGENTS.md tools .claude docs/plan docs/reviews`, Summe aus Einfügungen und Löschungen): **636** (+314, −322) vor diesem Beleg-Commit; Marke 900 und die Meldeschwelle 800 nicht erreicht; Schätzung des Architect 680 bis 780 mit Fall: darunter, weil die Verweise weniger Dateien betrafen als angesetzt und die Verweise weniger Dateien betrafen. Eine Review-Runde steht noch aus (Schätzung wie bei den Vorgängern rund 20 bis 25 % darauf).

**Rückgaben an den Architect:** keine; keine Randform außerhalb von §6 aufgetreten. **Hinweis an den Planner (kein Befund):** (1) `MR-002` (Zustandsabsatz, Geltungsbereich) ist nach §5 Closure-Arbeit und nicht angefasst. (2) In `slice-harness-lh-links-pflicht` trug die Zählung „zwei Schichten (Spezifikation: neue Stelle in §11 und Links; Harness: …)“ eine Schicht *Spezifikation*, die ihre Stelle nun in `harness/vertraege/` hat; ich habe das Etikett in „Vertrag“ geändert und die Zahl nicht angefasst. Ob der Plan danach weiter auf zwei Schichten zählt (die neue Stelle liegt in der Schicht Harness), entscheidet der Planner. (3) Die zwei Zeilen von `harness/README.md` zu `gates-schnell-gegenprobe` und `stop-hook-gegenprobe` nennen `SPEC-052` ohne Link und waren schon vorher ohne Link; unverändert.

**Beobachtung:** keine Evidence-Datei vom Implementer vorgeschlagen (Urteil der Closure).

### Closure-Notiz

**Geliefert:** `SPEC-052` steht in `harness/vertraege/SPEC-052-gate-kette.md` (262 Zeilen, 22 Punkte), wortgleich gegenüber `spezifikation.md §11` am Stand `231089e` bis auf die drei benannten Stellen (Punkt 12 an zwei Stellen, Punkt 22 (b) nach V-343); das `README.md` des Ordners führt den Vertrag und die fünf Historie-Zeilen. Die Spezifikation führt weder Vertrag noch Einleitung noch Überschrift von §11 noch die fünf Zeilen (`## 11. Historie` heißt jetzt so, kein Verweis auf `harness/`). Verweise: `harness/README.md`, acht Kommentardateien, `AGENTS.md` §3.12 (Ausnahme und Zwischenstand gestrichen), vier Pläne in `open/` und `next/` und der Kopf von `slice-harness-lint-warten-ohne-frist` zeigen auf den neuen Ort. Der Fall `p12-selbst-vertrag` (139 Fälle) prüft die Selbständerung des Vertrags am neuen Pfad. Verifikation: [Bericht](../../../reviews/2026-10-13-verifikation-slice-harness-vertraege-umzug-gate-kette.md) bestätigt die drei Liefer-Punkte, `make gates` (Exit 0, volle Kette auf `491ef31`), `make docs-check`, `make kopf-check`, die Unberührtheit von Handbuch und README (kein Pfad unter `docs/user/`, Wurzel-`README.md`, `internal/`, `cmd/`, `test/` im Diff) und acht selbst gefahrene Mutationen; keine HIGH/MEDIUM-Befunde.

**Größe** (Metrik `AGENTS.md` §3.14, `git diff --numstat 8b1ede6..491ef31 -- spec harness AGENTS.md tools .claude docs/plan docs/reviews`): **766** (+442, −324, Stand der Verifikation, mit Review-Report, ohne Verifikationsbericht); mit dem Verifikationsbericht bis `e578c54` 821 (nur Zeitdokument über der Meldeschwelle 800); Marke 900 nicht erreicht. Die Schätzung des Architect (680 bis 780 mit Fall) lag am oberen Rand und wurde im Stand der Verifikation getroffen; die Review-Runde stand in ihr (Lerneintrag der Vorgänger, rund 20 %).

**Ausgang je Risiko aus §6:** alle vier *entfallen*, mit Grund in §6 (Punkt 22 (b): Zusage unverändert, `p22-*` grün · fünfter Plan mit `§11`: drei unabhängige Suchen fanden nur die bekannten Treffer · Punkt 12: Fall und Mutation geliefert · Größe 766 unter Marke).

**Ausgang F-3 (Review LOW, Verifikation L-1; Nehmer `slice-harness-lh-links-pflicht`, Zeile „zwei Schichten (Spezifikation; …)“ gegen „(Vertrag; …)“):** *entfallen*, mit Grund: Der Nehmer zählt zwei Liefer-Punkte und zwei Schichten in jeder Lesart (die neue Stelle liegt in der Schicht Harness, die Schicht Spezifikation bleibt wegen der Links in `spec/`); die Etiketten meinten dieselbe Zählung. Die Nachzählung steht jetzt in §8 des Nehmers, im selben Commit wie diese Closure, mit der Kennung dieses Gebers; ebenso in §8 von `slice-harness-mutation`, `slice-harness-coverage` und `slice-harness-meldungskatalog-gate` (dort drei Liefer-Punkte, höchstens zwei Schichten; beim letzten drei mit der bedingten Zeile des Bestands-Befunds, was schon vor der Sendung so war und beim Architect liegt). Das Etikett in der datierten Zeile bleibt, weil eine Änderung die Chronik der Zählung verfälschte. Dass die Nachzählung erst hier steht, ist das Auftreten im Register (unten). F-1 (Kommentare der Gegenprobe) und F-2 (261/262, 3476/3477) sind behoben (`c85257c`, `491ef31`); F-4 bis F-6 (INFO) verlangten keine Aktion.

**`MR-002`:** [MR-002](../../../../harness/conventions/MR-002-harness-vertraege-ausserhalb-der-spezifikation.md) bleibt **aktiv**, nicht nach `conventions/done/`. Der Auflösungs-Trigger (das Kurs-Repo nennt einen Ort außerhalb der Produktspezifikation für Werkzeugverträge, oder das Repo führt wieder einen Abschnitt in `spec/spezifikation.md` ein; für Adaption 5 ein Vertrag, der sich auf eine ADR oder einen Slice stützt) ist nicht eingetreten; der Umzug vollzieht die Adaption, löst sie nicht auf, denn Ort, Rang und Referenz-Richtung von `harness/vertraege/` weichen weiter von der Baseline-Regel ab. Nachgezogen: der Zustandsabsatz nennt den Zustand (jeder Vertrag eines Harness-Werkzeugs steht in `harness/vertraege/`), im Geltungsbereich entfällt „bisher … künftig“, und die zwei Sätze der Adaption 5, die den Umzug als künftig beschrieben, nennen den Zustand (§3.7). `harness/conventions.md` bleibt unverändert: die Tabelle der aktiven Adaptionen nennt den Geltungsbereich ohne Chronik und ist mit der Datei im Einklang.

**Roadmap:** Die Zeile „Als Nächstes“ nennt `slice-v1-abschluss-tls-client` zuerst; der Umzug ist erledigt und steht dort nicht mehr. Das Drift-Log trägt keine Schließung (`AGENTS.md` §3.7); es ist keine Umplanung eingetreten.

**Beobachtungs-Register:**
- `BEO-REPO/slice-waechst-durch-uebernahmen`: Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.13 *Nachzählen beim Eintragen*): vier Nehmer erhielten die Sendung ohne Nachzählung in §8 (Review F-3); Evidence-Datei `evidence/slice-harness-vertraege-umzug-gate-kette.md`. Zweites Auftreten in Folge bei einer Sendung, die nur den Ort ändert.
- `BEO-REPO/plan-folgt-korrektur-nicht`: Wiederaufgetreten, Regel bleibt (`AGENTS.md` §3.9): §7 nannte 261 statt 262 Zeilen und zwei verschiedene Zeilenzahlen im Vergleichsbefehl (Review F-2); Evidence-Datei `evidence/slice-harness-vertraege-umzug-gate-kette.md`.
- `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung`: **kein drittes Auftreten, 2× bleibt, offen.** Die Beobachtung hält fest, dass die Gegenprobe nicht prüft, wie die Umgebung von Claude Code die Antwort des Hooks liest (und was der Treiber bei Überschreiten einer Dauer tut). Dieser Slice ändert den Ort des Vertrags und fügt der Gegenprobe den Fall `p12-selbst-vertrag` hinzu, der eine Pfadklasse des Hooks in einem Temp-Repo prüft; an der Grenze zum Treiber ändert er nichts, und weder Review noch Verifikation benennen sie als neu offen. Das „Treffer“ in §8 vor dem Code las die Beobachtung als berührt, nicht als aufgetreten (Randform 2 betraf die Selbständerung, nicht die Umgebung). Keine Evidence-Datei, kein Ausgang; die Schwelle 3× ist nicht erreicht.
- `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`: kein neues Auftreten, 2× bleibt; dieser Slice war die größte Hälfte des Schnitts (766 gegen Marke 900), die Marke trat nicht ein, und die Meldeschwelle 800 überschritt nur der Verifikationsbericht (821).
- `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum`: keine Evidence-Datei; der Ort der Festlegung wurde nicht übersehen. Mit diesem Slice ist der Umzug aller fünf Verträge vollzogen; der Stand bleibt `offen` (1×).
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: kein Auftreten; F-1 (Kommentare der Gegenprobe enger als die Prüfung) ist die Gegenrichtung, kein Zusage-Überschuss, und ist behoben.

**Lerneintrag (Steering Loop):** Form: Wiederaufgetreten, Regel bleibt (`BEO-REPO/slice-waechst-durch-uebernahmen`, Herkunfts-Anker `seit slice-v1-abschluss-einspielen`); geändert wird `AGENTS.md` nicht. Die Regel *Nachzählen beim Eintragen* fehlte zum zweiten Mal in Folge bei einer Sendung, die nur einen Ort ändert; hier trug der Implementer die Sendung in vier Nehmer ein und meldete die uneinheitliche Zählung in einem von ihnen als Hinweis an den Planner, statt sie im selben Commit in §8 zu führen. Dass keiner der vier über der Grenze lag, ist das Glück dieses Falls, nicht die Wirkung der Regel. Was funktionierte: Die Randform 2 (Selbstbezug „diese Spezifikation“ im Vertragstext) hat der Architect vor dem Code gefunden, der Nutzer hat sie entschieden, und der Fall `p12-selbst-vertrag` fing die Mutation rot; die Wortgleichheit des Umzugs ist durch einen Befehl belegt (`diff` gegen `231089e`, zwei Hunks), nicht durch Lesen. Vorschläge an den Nutzer, nicht als Eintrag (derselbe erste Vorschlag stand schon in der Closure von `slice-harness-vertraege-umzug-lint`): (1) `AGENTS.md` §3.13 *Nachzählen beim Eintragen* um „auch bei einer Ortsangabe, je Nehmer eine Zählung in dessen §8, und liegt der Nehmer schon darüber, nennt die Zählung das und den Schnitt“ ergänzen; (2) `.claude/commands/implement-slice.md` Schritt 19: ein Hinweis an den Planner ersetzt die Nachzählung nicht, sie steht im Commit der Sendung. Retirement-Check bestanden (Regel weder entfernt noch gelockert).

**Paarungen:** (a) *Anker* — der Lerneintrag trägt kein Feld `liegt in`, keine Verkörperung; nichts zu prüfen. (b) *Folge-Slice* — die Sendung aus `slice-harness-vertraege-umzug-ordner` (Ausnahme für `SPEC-052` in `AGENTS.md` §3.12) ist getragen: `grep -n "SPEC-052\|§11" AGENTS.md` zeigt keinen Ort in `spec/`, §3.12 nennt `harness/vertraege/` mit der Kennung dieses Slice. Die Sendung dieses Slice an die vier Nehmer (Ort neuer Kennungen, Nummer 053) ist getragen: `grep -c "slice-harness-vertraege-umzug-gate-kette"` findet die Kennung in `slice-harness-lh-links-pflicht` (1), `slice-harness-mutation` (2), `slice-harness-coverage` (2) und `slice-harness-meldungskatalog-gate` (1), jeder Nehmer liegt in `open/` oder `next/` und schließt den Ort nicht aus. `slice-harness-lint-warten-ohne-frist` trägt nur die umnummerierte *Historie* (§12 zu §11), keine Sendung. (c) *Register* — `BEO-REPO/slice-waechst-durch-uebernahmen` und `BEO-REPO/plan-folgt-korrektur-nicht` existieren mit neuer Evidence-Datei, die fünf weiteren genannten Kennungen existieren als Verzeichnis.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel `REPO`, Greenfield); der Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` am Stand `231089e` gesichtet. Treffer: `werkzeug-festlegung-ausserhalb-technik-stratum` (1×, offen: Ort der Festlegung; mit diesem Slice ist der Umzug vollständig) ·
`plan-folgt-korrektur-nicht` (41×, verkörpert: die vier Pläne mit `§11`) · `zusage-im-kommentar-weiter-als-pruefung` (51×, verkörpert; geplanter Sensor in `slice-harness-mutation`: Punkt 12 und der Fall aus Randform 2) ·
`stop-hook-gegenprobe-prueft-nicht-die-umgebung` (2×, offen: Variante (b) der Randform 2 ist entschieden, der Slice trifft die Beobachtung) · `schnitt-laesst-haelfte-an-der-grenze` (2×, offen: dieser Slice ist die größte Hälfte, Marke und Schätzung in der Metrik von `AGENTS.md` §3.14) ·
`groessenprobe-ohne-messbaren-stand` und `schaetzung-in-zeilen-unterschaetzt-tests` (je 3×, verkörpert in §3.14) · `slice-waechst-durch-uebernahmen` (3×, verkörpert: Nachzählung unten).

**Nachgezählt** (`AGENTS.md` §3.13): Liefer-Punkte 3 (Datei samt V-343, Spezifikation, Verweise). Schichten 2: Spezifikation (`spec/`) und Harness (`harness/`, `tools/`, `.claude/`, `AGENTS.md`); Planung ist keine Schicht, Produkt-Code unberührt. **Mit der Sendung aus `slice-harness-vertraege-umzug-ordner` nachgezählt:** `AGENTS.md` §3.12 steht schon in Liefer-Punkt 3 und in der Schicht Harness; es bleiben 3 Liefer-Punkte und 2 Schichten, die Schätzung 550 bis 750 (mit Fall 680 bis 780) wächst um eine Zeile, Marke 900 unberührt. Die Bindung aus `slice-harness-gate-kette-stop-hook-freigabe` fügt keinen Liefer-Punkt hinzu: Sie gehört zu Punkt 1 (dieselbe Datei), ihr Beleg wird in §7 eigens geführt; der Planner urteilt so, der Architect bestätigt vor dem Code.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
