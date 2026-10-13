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

**Berührte Spec-Stellen:** `SPEC-052` (Vertrag von `gates-schnell` und Stop-Hook: der Ort wechselt von `spezifikation.md §11` nach `harness/vertraege/`; Punkt 22 (b) wird eindeutig gefasst) · `spezifikation.md §11` (*Harness-Werkzeuge*: Einleitung und letzter Vertrag gehen, die Überschrift entfällt) · `spezifikation.md §12` (*Historie*, verliert fünf Zeilen)

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

**Bindung aus `slice-harness-gate-kette-stop-hook-freigabe` (V-343, benannte Spec-Lücke).** Punkt 22 (b) liest sich einzeln als „leere Umgebung: bei identischer Liste `approve`, bei jedem Unterschied
`block`“; erst der folgende Satz („beide Aufrufarten gelten auch für die Freigabe nach Punkt 9“) hebt das auf. Der Code ist richtig; die Gegenprobe prüft die Freigabe in beiden Aufrufarten
(`p22-freigabe-unterverzeichnis`, `p22-freigabe-leere-umgebung`, Mutanten `m-22-hash-relativ`, `m-22-hash-home`). Dieser Slice ist der Nehmer: Der Architect fasst (b) bei der Neuschrift so,
dass der Satz für sich allein gelesen nicht das Gegenteil sagt; mehr ändert er an Punkt 22 nicht. Das ist die einzige inhaltliche Änderung an einem Vertrag im Umzug (Entscheidung des Nutzers).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ordner, `README.md`, Rang, Reihe der Kennungen und die Verträge von `kopf-check`, `abdeckung` und `commit-msg`: `slice-harness-vertraege-umzug-ordner` (Voraussetzung, er liegt vor diesem Slice in `done/`). Dieser Slice trägt nur den Vertrag, seine Historie-Zeilen und das Entfernen der Einleitung nach.
- Der Vertrag des Lint-Profils: `slice-harness-vertraege-umzug-lint` (Voraussetzung, liegt davor in `done/`).
- Jede andere Änderung an Wortlaut oder Zusage von `SPEC-052`: Der Umzug ist wortgleich (Entscheidung des Nutzers); was der Architect darüber hinaus für nötig hält, ist eine Randform in §6 und geht, wenn es mehr als die Bindung ist, zurück an den Planner (Klasse 3: ein anderer Vorgang, Arbeit am Vertrag statt am Ort).
- Die Stop-Hook-Logik, `tools/harness/dateiliste.sh`, `tools/harness/working-tree-hash.sh`, `tools/harness/record-gates.sh`, `tools/harness/gates-schnell.sh`: nur Kommentare mit dem alten Pfad ändern sich (Klasse 4: kein Produkt-Code, keine Logik). Ausnahme nur, wenn der Architect in §6 Randform 2 entscheidet, dass die Selbständerung eines Vertrags in der Gegenprobe einen Fall bekommt.
- Pläne, die einen *bestehenden* Vertrag fortschreiben (`slice-harness-commit-struktur-id`, `slice-harness-abdeckung-gate`, `slice-harness-lint-warten-ohne-frist`): bereits in den ersten beiden Slices nachgezogen.
- Zeitdokumente: Links in Reviews und `done/` nur dort, wo sie sonst brechen; Fließtext bleibt ([ADR-0040](../../adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), akzeptiertes Negativ).

## 2. Definition of Done


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 — die Datei samt V-343.** `harness/vertraege/SPEC-052-gate-kette.md` trägt den Text von `SPEC-052` aus `spezifikation.md §11` am Stand `231089e` (Zeilen 3216 bis 3477),
      unverändert bis auf die Abweichungen, die `slice-harness-vertraege-umzug-ordner` für alle Verträge festlegt, und bis auf Punkt 22 (b), der für sich allein gelesen eindeutig ist
      (Bindung aus `slice-harness-gate-kette-stop-hook-freigabe`, V-343; sein Beleg wird in §7 eigens geführt). Das `README.md` des Ordners führt den Vertrag und
      die fünf Zeilen seiner Historie. Beleg in §7: der Vergleich (Befehl und Ausgabe: der einzige Unterschied außer den Abweichungen ist (b)) und der neue Wortlaut von (b).
- [ ] **Liefer-Punkt 2 — die Spezifikation.** `spec/spezifikation.md` führt weder Vertrag noch Einleitung noch Überschrift von §11 noch die fünf Zeilen der Historie; die Abschnitte danach tragen ihre Nummern
      wie die Entscheidung des Architect festlegt (§6 Randform 3); die Spezifikation verweist nicht auf `harness/`; `make docs-check` grün.
- [ ] **Liefer-Punkt 3 — die Verweise.** Jeder Verweis auf `SPEC-052` und jede Ortsangabe `spezifikation.md §11` in einem lebenden Plan zeigt auf den neuen Ort: `harness/README.md` (Zeilen von `make gates-schnell`,
      `make gates-schnell-gegenprobe`, `make stop-hook-gegenprobe`), die Kommentare von `.claude/hooks/stop-require-gates.sh`, `tools/harness/gates-schnell.sh`, `tools/harness/gates-schnell-gegenprobe.sh`,
      `tools/harness/stop-hook-gegenprobe.sh`, `tools/harness/record-gates.sh`, `tools/harness/dateiliste.sh`, `harness/mk/gates-schnell.mk`, `harness/mk/stop-hook-gegenprobe.mk`; in den Plänen
      `slice-harness-lh-links-pflicht`, `slice-harness-meldungskatalog-gate`, `slice-harness-coverage`, `slice-harness-mutation` der Ort neuer Kennungen (Kopf, §1, §3, §6; die nächste freie Nummer ist 053);
      `AGENTS.md` §3.12 und die Rollen- und Command-Texte unter `.claude/`, falls sie `spezifikation.md §11` oder „die Spezifikation“ für Harness-Werkzeuge nennen.
      Beleg in §7: `grep -rn "spezifikation.md#spec-052\|spezifikation.md SPEC-052\|spezifikation.md §11\|spezifikation.md.*§ ?11" --exclude-dir=.git --exclude-dir=.harness .` ohne Treffer außerhalb von Zeitdokumenten.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Handbuch und README unberührt: Das Produktverhalten ändert sich nicht (`AGENTS.md` §3.11).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)


Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-052-gate-kette.md` | neu | Text aus §11 (Zeilen 3216 bis 3477 am Stand `231089e`); Punkt 22 (b) neu gefasst (V-343) |
| `harness/vertraege/README.md` | update | Eintrag des Vertrags, fünf Zeilen Historie, ggf. Wegfall der Dopplung der Einleitung (§6 Randform 4) |
| `spec/spezifikation.md` | update | §11 ganz entfernt (Überschrift, Einleitung, Vertrag), fünf Zeilen der *Historie*; die Folge-Abschnitte (*Historie*) nach Entscheidung des Architect (§6 Randform 3) |
| `harness/README.md` | update | §Sensors: Links und Nennungen in den Zeilen von `make gates-schnell`, `make gates-schnell-gegenprobe`, `make stop-hook-gegenprobe` |
| `.claude/hooks/stop-require-gates.sh`, `tools/harness/gates-schnell.sh`, `tools/harness/gates-schnell-gegenprobe.sh`, `tools/harness/stop-hook-gegenprobe.sh`, `tools/harness/record-gates.sh`, `tools/harness/dateiliste.sh`, `harness/mk/gates-schnell.mk`, `harness/mk/stop-hook-gegenprobe.mk` | update (nur Kommentare) | „Vertrag: spec/spezifikation.md SPEC-052“ nennt den neuen Pfad |
| `tools/harness/stop-hook-gegenprobe.sh` | update, nur nach Entscheidung des Architect (§6 Randform 2) | ein Fall für die Selbständerung des Vertrags am neuen Pfad (`p12-selbst-vertrag`), mit Mutation; sonst unverändert |
| `AGENTS.md`, `.claude/agents/architect.md` | update, falls sie den alten Ort nennen | Ortsangabe der Spezifikation für Harness-Werkzeuge; Stand zum Planungszeitpunkt: `AGENTS.md` §3.12 ist von `slice-harness-vertraege-umzug-ordner` angepasst, `architect.md` Zeile 35 nennt „die Spezifikation des Gates“ |
| `docs/plan/planning/open/slice-harness-lh-links-pflicht.md`, `slice-harness-mutation.md`; `next/slice-harness-meldungskatalog-gate.md`, `next/slice-harness-coverage.md` | update | Ort neuer Kennungen: `harness/vertraege/` statt `spezifikation.md §11`, nächste freie Nummer 053 (`AGENTS.md` §3.9); die Ortsänderung fügt keinem dieser Pläne einen Liefer-Punkt hinzu und keine Schicht (die Festlegung liegt jetzt in der Schicht Harness, die jeder von ihnen schon berührt) |
| `docs/plan/planning/in-progress/roadmap.md` | update | Stand erst bei der Closure |
| `docs/reviews/2026-10-13-verifikation-slice-harness-gate-kette-stop-hook-freigabe.md` und weitere Reviews mit Link | update (Links) | nur, wo das Doku-Gate sonst bricht |

Ansatz: wie in `slice-harness-vertraege-umzug-ordner`: Der Text wandert per Ausschneiden, der Vergleich in §7 belegt es; Punkt 22 (b) ist die einzige Zeile, die der Architect schreibt. Gegenprobe und Code
bleiben, soweit der Architect keinen Fall für die Selbständerung verlangt (§6 Randform 2); in dem Fall zählt der Fall als Teil von Liefer-Punkt 1 (er prüft die Zusage, die der Vertrag am neuen Ort trägt) und die
Schätzung steigt um rund 60 (Faktor 2 für Gegenprobe-Code, `AGENTS.md` §3.14).

## 4. Trigger


Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-vertraege-umzug-lint` liegt in `done/` (WIP-Limit 1). Erster Schritt nach dem Start, vor jedem Commit an Spezifikation, `harness/` oder
`tools/`: Der Architect entscheidet die Randformen aus §6 und schreibt den Wortlaut von Punkt 22 (b) (V-343) und gibt ihn dem Implementer als Text.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und Löschungen von
  `git diff --numstat <Basis>..<letzter Commit des Slice> -- spec harness AGENTS.md tools .claude docs/plan docs/reviews` liegt über 900
  (`<Basis>` ist der Commit des Übergangs `next → in-progress` auf dem Hauptzweig; Metrik nach `AGENTS.md` §3.14; Schätzung des Planners 550 bis 750: 262 Zeilen Vertragstext, doppelt gezählt, dazu Historie, Einleitung, Verweise und die Pläne; mit dem Fall aus Randform 2 rund 60 mehr, weil Gegenprobe-Code gerechnet wird).
- `in-progress` → `open` (blockiert — Carveout?): Der Architect findet für (b) keinen Wortlaut, der den Vertrag nicht inhaltlich ändert (dann entscheidet er, ob V-343 ein eigener Slice wird), oder
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

- Punkt 22 (b) wird neu gefasst, und die Neufassung ändert die Zusage (aus „eindeutig“ wird „anders“) — **Ausgang:** offen bis Closure.
- Nach dem Wegfall von §11 bleiben Pläne, die den Abschnitt noch nennen (die vier aus der DoD sind bekannt, ein fünfter ist möglich); `grep` findet nur die Schreibweisen, die er kennt (`BEO-REPO/plan-folgt-korrektur-nicht`) — **Ausgang:** offen bis Closure.
- Die Kette der Stop-Hooks: Der Vertrag liegt nicht mehr in `spec/`; Punkt 12 sagt „diese Spezifikation“, und der Fall `p12-selbst-spec` prüft `spec/spezifikation.md`, nicht den Vertrag — **Ausgang:** offen bis Closure (Randform 2).
- Größe: Der Vertrag ist der größte der fünf; liegt die Summe über der Marke, schneidet der Planner die Verweise ab (`BEO-REPO/schnitt-laesst-haelfte-an-der-grenze`, 2×) — **Ausgang:** offen bis Closure.

**Randformen** (`AGENTS.md` §3.12): Es gelten die Entscheidungen aus §6 von `slice-harness-vertraege-umzug-ordner`. Zusätzlich, vom Architect vor dem ersten Commit, Ort: `harness/vertraege/README.md` und der Vertragstext:

1. *Punkt 22 (b).* Neuer Wortlaut, den der Architect schreibt; Maßstab: allein gelesen sagt (b), dass bei identischer Liste `approve` und bei jedem Unterschied `block` gilt, *und dass die Freigabe nach Punkt 9 in beiden Aufrufarten gilt*; die Gegenprobe prüft beides schon.
2. *Selbständerung (Punkt 12) am neuen Ort.* Der Vertrag sagt „eine Änderung an … dieser Spezifikation ist nicht inert und verlangt den vollen Lauf“; gemeint war `spec/spezifikation.md`, und `p12-selbst-spec` prüft diesen Pfad. Nach dem Umzug liegen die Verträge unter `harness/vertraege/`, das nicht auf der Positivliste steht (jede Änderung verlangt den vollen Lauf, wie bisher). Entscheidung des Architect: (a) Wortlaut wortgleich lassen und „diese Spezifikation“ bedeutet jetzt die Datei des Vertrags (Gegenprobe prüft den Pfad `spec/spezifikation.md`, nicht den Vertrag: die Zusage trüge keinen Fall für ihren neuen Ort), oder (b) den Wortlaut auf Hook, Skripte und „diese Datei und `spec/spezifikation.md`“ fassen und einen Fall `p12-selbst-vertrag` mit Mutation (Positivliste um `harness/vertraege/` erweitert) ergänzen. Empfehlung des Planners: (b), weil `AGENTS.md` §3.11 und §3.10 sonst eine Zusage ohne Fall zulassen.
3. *Nummern der Abschnitte nach dem Wegfall von §11.* Die Spezifikation führt *Historie* als §12; entweder bleibt die Nummer (Lücke bei §11) oder sie wird §11. Verweise auf `spezifikation.md §12` stehen im Bestand in `slice-harness-lint-warten-ohne-frist` und `slice-harness-commit-struktur-id` (Kopf, §3); ob Anker der Form `#12-historie` irgendwo verlinkt sind, prüft der Architect mit `grep`. Vorschlag: nachnummerieren und diese Verweise nachziehen, weil eine Lücke bei §11 als Fehler gelesen wird; bleibt die Nummer, entfällt der Nachzug.
4. *Dopplung der Einleitung.* Die Einleitung von §11 steht im `README.md` seit `slice-harness-vertraege-umzug-ordner`; hier entfällt der Text in der Spezifikation (Löschen, keine Änderung am `README.md`).
5. *Kennung für Verträge nach diesem Slice.* Die nächste freie Nummer ist 053; sie steht im `README.md` des Ordners und in den vier Plänen aus Liefer-Punkt 3.
6. *`MR-002`.* Der Zustandsabsatz („Der Umzug geschieht erst im Umzug-Slice …“) wird auf den Stand gebracht; Adaption und Auflösungs-Trigger bleiben (Entscheidung des Planners, Urteil beim Review).

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
`stop-hook-gegenprobe-prueft-nicht-die-umgebung` (Treffer nur, wenn (b) der Randform 2 gewählt wird) · `schnitt-laesst-haelfte-an-der-grenze` (2×, offen: dieser Slice ist die größte Hälfte, Marke und Schätzung in der Metrik von `AGENTS.md` §3.14) ·
`groessenprobe-ohne-messbaren-stand` und `schaetzung-in-zeilen-unterschaetzt-tests` (je 3×, verkörpert in §3.14) · `slice-waechst-durch-uebernahmen` (3×, verkörpert: Nachzählung unten).

**Nachgezählt** (`AGENTS.md` §3.13): Liefer-Punkte 3 (Datei samt V-343, Spezifikation, Verweise). Schichten 2: Spezifikation (`spec/`) und Harness (`harness/`, `tools/`, `.claude/`, `AGENTS.md`); Planung ist keine Schicht, Produkt-Code unberührt. Die Bindung aus `slice-harness-gate-kette-stop-hook-freigabe` fügt keinen Liefer-Punkt hinzu: Sie gehört zu Punkt 1 (dieselbe Datei), ihr Beleg wird in §7 eigens geführt; der Planner urteilt so, der Architect bestätigt vor dem Code.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
