# Slice slice-harness-gate-kette-stop-hook: Der Stop-Hook wählt die Kette am geänderten Pfad

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
Nutzers vom 2026-10-12 (Gate-Kette vor `slice-v1-abschluss-umgebung-anderes-kommando` vor
`slice-v1-abschluss-antwortvergleich`): direkt hinter `slice-harness-gate-kette-schnell`
(WIP-Limit 1); Reihenfolge in §4 *Start* und im Drift-Log der Roadmap. Der Zuschnitt in zwei Slices
ist **Empfehlung des Planners**, nicht entschieden.

**Bezug:** kein Anspruch des Produkts — der Slice liefert ein Werkzeug des Harness.
Bindung an Entscheidungen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Wahl der
Kette am geänderten Pfad, Nachweis als Dateiliste; Accepted 2026-10-12), das sich auf
[`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) stützt.

**Berührte Spec-Stellen:** `SPEC-052` (von `slice-harness-gate-kette-schnell` in
`spezifikation.md` §11 angelegt, mit den Randformen vom Architect entschieden; dieser Slice ändert
sie nicht)

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

**Ziel:** Der Nachweis des vollen Laufs ist eine Dateiliste statt eines Hashs
(`tools/harness/record-gates.sh`, `tools/harness/working-tree-hash.sh`), und der Stop-Hook
`.claude/hooks/stop-require-gates.sh` wählt die Kette nach `SPEC-052`: Er gibt frei, wenn jede seit
dem Stand des letzten **vollen** grünen Laufs geänderte, hinzugekommene oder gelöschte Datei auf der
Positivliste liegt (zu Beginn `docs/plan/planning/` und `docs/reviews/`) und der Stand der schnellen
Kette grün ist; jeder andere Fall verlangt `make gates` (fail-closed). Eine Gegenprobe des Hooks macht
je Randform einen Fall, der nur im inerten Fall freigibt. Der Hook ist das Restrisiko der Entscheidung
([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Konsequenzen*): Ohne seine Gegenprobe
schließt der Slice nicht.

**Voraussetzung (Geber):** `slice-harness-gate-kette-schnell` liegt in `done/`: `SPEC-052` mit allen
Randformen, das Ziel der schnellen Kette und ihr Stand. **Namen, die der Implementer des Gebers
wählt** (Ziel der schnellen Kette, Datei ihres Stands, Gegenprobe der Kette): trägt die Closure des
Gebers hier ein, mit Nachzählung im selben Commit (`AGENTS.md` §3.13); bis dahin steht hier
*offen* — Ziel: ?, Stand-Datei: ?, Gegenprobe der Kette: ?. Der Implementer beginnt nicht, solange
`grep` nach den Namen in diesem Plan leer bleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `SPEC-052`, die schnelle Kette, ihr Stand und deren Gegenprobe — `slice-harness-gate-kette-schnell`
  (liegt vor diesem Slice in `done/`): Dieser Slice setzt sie voraus und ändert die Spezifikation nicht;
  gibt der Implementer eine Randform zurück, entscheidet der Architect in `SPEC-052`.
- Aufnahme weiterer Pfade in die Positivliste (`docs/plan/adr/`, `spec/`, `docs/user/`, `README.md`) —
  Bestand bleibt bewusst stehen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Setzungen des Nutzers*; eine Aufnahme ist eine Änderung von
  `SPEC-052` mit dem Beleg, dass kein langsames Gate den Pfad liest. Es gibt dafür keinen Slice und
  keine Kennung.
- Änderung von `make gates` — [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Entscheidung* Punkt 1: Die volle Kette bleibt der Nachweis
  für alles, was nicht auf der Positivliste steht.
- Die Hooks `.claude/hooks/pretooluse-command-guard.sh` und `.claude/hooks/span-emit.sh` und den
  Eintrag in `.claude/settings.json` — anderer Vorgang: Der Slice ändert die Wahl der Kette im
  Stop-Hook, nicht den Befehls-Wächter oder die Erfassung.
- Eine Prüfung gegen die echte Hook-Umgebung von Claude Code — die Gegenprobe fährt den Hook in
  Temp-Repos mit dessen Eingabe und liest dessen Ausgabe; ob die Umgebung die Antwort so liest, prüft
  sie nicht (Risiko in §6).
- Produkt-Code, Benutzerhandbuch, `README.md` — Schicht-Abgrenzung.
## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1:** Der Nachweis des vollen Laufs ist die Dateiliste nach `SPEC-052`:
      `tools/harness/record-gates.sh` schreibt sie nur nach grünem `make gates`,
      `tools/harness/working-tree-hash.sh` (oder sein Nachfolger) liefert dieselbe Liste für den
      Arbeitsbaum; beide teilen sich die Logik, keine Dopplung zwischen Schreiber und Hook.
- [ ] **Liefer-Punkt 2:** `.claude/hooks/stop-require-gates.sh` wählt die Kette nach `SPEC-052`
      (Positivliste, Stand des vollen und Stand der schnellen Kette, Rückfall auf voll bei jedem
      anderen Fall); der Schleifenschutz (`stop_hook_active`) bleibt.
- [ ] **Liefer-Punkt 3:** Die Gegenprobe des Hooks macht je Randform aus §6 einen Fall in einem
      Temp-Repo, der nur im inerten Fall freigibt, und hängt an `GATE_CHECKS`;
      `harness/README.md` §Sensors führt ihre Zeile mit Bindung an [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md);
      `harness/conventions.md` (Träger-Tabelle, Zeile `grundlagen-durchsetzungsschicht.md`) nennt die
      Träger des geänderten Hooks. Der Name der Gegenprobe ist ein anderer als der von
      `make hook-gegenprobe` (die gehört dem Commit-Träger `.githooks/commit-msg`).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). Der Reviewer prüft
      fail-closed ausdrücklich: Jede Abweichung von „inert und Stand grün“ verlangt `make gates`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Die Mutationen und ihre roten Fälle stehen als Tabelle in §7 (Zusage · Mutation · roter
      Fall), vom Implementer selbst gefahren (`AGENTS.md` §3.10).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: eine weitere Datei
      `evidence/slice-harness-gate-kette-stop-hook.md` in
      `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/` mit der Messung dieses Slice (Dauer von
      `make gates` und der schnellen Kette; Freigabe eines Stops nach einem Doku-Commit; Last und
      Laufzeit dessen, was der Hook spart); der Zähler folgt aus den Dateien (3×) und der Eintrag
      bekommt seinen Ausgang. Die Messung ist Beleg, kein Liefer-Punkt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/record-gates.sh` | update | schreibt die Dateiliste statt eines Hashs, nur nach grünem `make gates` |
| `tools/harness/working-tree-hash.sh` | update / ersetzt | liefert dieselbe Liste für den Arbeitsbaum (getrackt und untracked, Symlinks mit Ziel, gelöschte Dateien); Name und Form nach `SPEC-052` |
| `.claude/hooks/stop-require-gates.sh` | update | Wahl der Kette: Positivliste, Stand des vollen und der schnellen Kette, Rückfall auf voll |
| `tools/harness/` (neue Gegenprobe des Hooks) und ihr Fragment unter `harness/mk/` | neu | je Randform ein Fall in einem Temp-Repo; hängt an `GATE_CHECKS` |
| `harness/README.md` §Sensors | update | Zeile der Gegenprobe, Bindung [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) |
| `harness/conventions.md` (Träger-Tabelle, Zeile `grundlagen-durchsetzungsschicht.md`) | update | nennt die Träger des geänderten Hooks (Dateiliste, Hook, gegebenenfalls ersetzte Datei) |

- Der Hook liest den Nachweis über dieselbe Funktion, die ihn schreibt; ein zweiter Leser der Liste
  wäre die Logik-Dopplung, die der Kopf von `working-tree-hash.sh` heute ausschließt.
- Der Implementer prüft vor dem Commit, ob eine weitere lebende Stelle den Hash als Nachweis nennt
  (`grep -rn "gates-passed\|working-tree-hash\|record-gates"` über `AGENTS.md`, `harness/`,
  `.claude/`, `tools/`, `docs/`) und zieht sie im selben Commit nach.
- Kein Produkt-Code, kein Handbuch, keine Tests unter `internal/` oder `test/`.
## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-gate-kette-schnell` liegt in `done/` (WIP-Limit 1);
dieser Slice steht direkt dahinter und vor `slice-v1-abschluss-umgebung-anderes-kommando` und
`slice-v1-abschluss-antwortvergleich` (Drift-Log der Roadmap vom 2026-10-12). Die Namen aus §1 stehen
im Plan. Die Randformen aus §6 sind in `SPEC-052` entschieden (der Architect hat sie vor dem Code des
Gebers geschrieben); dieser Slice ändert sie nicht.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Prüfung des Architect schätzt den Diff
  über eine Review-Sitzung (über 750 Zeilen mit Tests). Schnitt dann nach Lieferwert: Dateiliste
  mit dem Hook im Stand „alles verlangt voll“ als einer, die Positivliste mit der Freigabe als der
  nächste; beide mit der Gegenprobe ihrer Zusagen, denn ohne Gegenprobe schließt keiner.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer nimmt
  [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) zurück, oder das Format der Antwort des
  Stop-Hooks ändert sich so, dass die Gegenprobe nicht mehr gegen die echte Eingabe und Ausgabe des
  Hooks läuft.
## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (vom Auftraggeber gefahren), die Gegenprobe des Hooks rot unter
jeder Mutation der Tabelle in §7, die Messung als Evidence im Register, Closure-Notiz mit
Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12): Der Slice liefert einen neuen Vertrag (Wahl der Kette am Pfad).
Die Randformen sind in `SPEC-052` entschieden, vom Architect vor dem ersten Code-Commit von
`slice-harness-gate-kette-schnell` (dessen §6 führt die Tabelle mit den Fragen; die Punktnummer in
`SPEC-052` trägt die Closure des Gebers hier ein). Eine nicht genannte oder offene Randform
entscheidet der Implementer nicht, er gibt sie zurück. Die Hook-Seite, je mit Ort `SPEC-052`:

- gelöschte oder umbenannte Datei (beide Pfade gelten als geändert);
- Pfad auf der Positivliste; Pfad auf der Grenze (`docs/plan/planning-x/`);
- Symlink (mit Ziel); untracked Datei; Merge-Konflikt; leerer Diff gegen den vollen Lauf;
- erster Lauf ohne Stand des letzten vollen grünen Laufs (fail-closed), dazu der Bestand:
  fehlender Stand bei sauberem Baum gibt heute frei;
- nicht lesbarer Vergleich; Datei der Dateiliste selbst gelöscht oder beschädigt;
- Stand der schnellen Kette fehlt, ist rot oder gehört zu einem anderen Arbeitsbaum;
- gleichzeitige Läufe (**Aufruf-Voraussetzung**: nacheinander oder gleichzeitig, `AGENTS.md` §3.12).

**Mutationen** (`AGENTS.md` §3.10): Der Implementer fährt je Zusage die Mutation selbst und trägt die
Tabelle in §7 ein. Vorgegeben:

| Zusage | Mutation | roter Fall |
|---|---|---|
| Nur ein Pfad der Positivliste gibt frei | `spec/` in die Positivliste aufnehmen | Änderung in `spec/` gibt frei → rot |
| Ein Pfad auf der Grenze zählt nicht | Präfix ohne Verzeichnisgrenze vergleichen | `docs/plan/planning-x/` gibt frei → rot |
| Gelöschte und umbenannte Dateien zählen mit beiden Pfaden | nur den neuen Pfad prüfen | Umbenennung aus einem nicht inerten Pfad in einen inerten gibt frei → rot |
| Ein Symlink zählt mit seinem Ziel | Link ohne Ziel vergleichen | Link in inertem Pfad mit Ziel außerhalb gibt frei → rot |
| Fehlender oder unlesbarer Stand des vollen Laufs verlangt `make gates` | fehlenden Stand als „nichts geändert“ lesen | Hook gibt frei → rot |
| Beschädigte Dateiliste verlangt `make gates` | Zeile ohne Hash überspringen | Hook gibt frei → rot |
| Ohne grünen Stand der schnellen Kette gibt der Hook nicht frei | Prüfung des Stands der schnellen Kette entfernen | rote oder fehlende Kette gibt frei → rot |
| Ein nicht lesbarer Vergleich verlangt `make gates` | Fehler von `git` ignorieren | Hook gibt frei → rot |
| Der Schleifenschutz bleibt | `stop_hook_active` nicht mehr lesen | Hook blockiert endlos → rot |
| `make gates` schreibt die Dateiliste nur bei grünem Lauf | Schreiben vor die Gates ziehen | roter Lauf hinterlässt Liste → rot |

**Risiken:**

- Der Hook ist das Restrisiko der Entscheidung: Ein Fehler dort bleibt ungeprüft, wo die Gegenprobe
  ihn nicht fängt ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Konsequenzen*) — **Ausgang:** entfallen: Die Gegenprobe macht je
  Randform einen Fall, jede Mutation der Tabelle ist rot gesehen, und ohne sie schließt der Slice
  nicht (DoD-Punkt 3).
- Die Gegenprobe fährt den Hook in Temp-Repos mit dessen Eingabe, nicht in der Umgebung von Claude
  Code; ob die Umgebung die Antwort (`decision`) so liest, prüft sie nicht — **Ausgang:** weiter offen:
  → `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (neu bei der Closure angelegt, 1×).
- Ein Pfad der Positivliste wird später von einem langsamen Gate gelesen ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *akzeptiertes
  Negativ*) — **Ausgang:** entfallen: Ein späterer Leser ändert Code oder Gates und löst damit den
  vollen Lauf aus; die ADR hat das Negativ akzeptiert.
- Zwei Läufe gleichzeitig (der Hook liest, während `make gates` die Liste schreibt) — **Ausgang:**
  entfallen: `SPEC-052` bindet die Zusagen an Läufe nacheinander (Aufruf-Voraussetzung,
  `AGENTS.md` §3.12), und der Hook sagt gleichzeitige Läufe nicht zu; entscheidet der Architect die
  Voraussetzung anders, wählt die Closure den Ausgang neu (atomare Schreibweise und ein Fall der
  Gegenprobe, oder weiter offen im Register).
- Die Größe zeigt sich erst in der Schätzung des Architect (`BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`,
  2×) — **Ausgang:** weiter offen: → dieser Eintrag, falls der Slice zurückgeht.
- Der Nutzen schrumpft, wenn die Gegenproben schneller werden ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Re-Evaluierungs-Trigger*)
  — **Ausgang:** entfallen: Der Slice sagt keine Zeitersparnis zu; die Messung steht als Beleg im
  Register.
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

- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (2×, offen) — Gegenstand: Mit der Messung dieses
  Slice (DoD) steht der Eintrag bei 3× und braucht seinen Ausgang; er lautet *verkörpert*, Zielort
  Hook und Kette, Herkunfts-Anker `seit slice-harness-gate-kette-stop-hook`, sofern der Hook die
  Laufzeit eines Doku-Commits zeigbar senkt, sonst *offen* mit der Begründung.
- `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (1×, offen) — berührt: Symlink, nicht lesbarer Eintrag
  und beschädigte Dateiliste sind Randformen des Hooks (§6); übergeht der Hook einen Eintrag still,
  ist es das zweite Auftreten.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — berührt: Die Positivliste liest
  sich als Zusage; die Gegenprobe prüft Pfade auf und neben der Grenze. Tritt es wieder auf, steht der
  Eintrag bei 3×.
- `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) — berührt: §4
  Rückführung, Schätzung in `slice-harness-gate-kette-schnell` §1.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×, offen) — der Slice ist wellenlos; die
  Randformen entscheidet der Architect vor dem Code des Gebers, §4 *Start* nennt es.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) und
  `BEO-REPO/slice-waechst-durch-uebernahmen` (verkörpert, `AGENTS.md` §3.13) — Mutationen in §6, kein
  Kommentar sagt mehr zu, als die Gegenprobe prüft, die Sendung des Gebers ist nachgezählt.

Dieser Slice führt den dritten Beleg für `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` und ist mit
`slice-harness-gate-kette-schnell` selbst die Antwort darauf; sonst erreicht keiner der Einträge mit
diesem Plan die Schwelle 3× neu.

**Zählung nach `AGENTS.md` §3.13:** Liefer-Punkte: 3 (Dateiliste; Hook mit Wahl der Kette; Gegenprobe
mit den Gate-Zeilen). Messung als Beleg, Gate-Läufe, Review und Closure-Notiz zählen nicht. Schichten:
1 — Harness-Werkzeug (`tools/harness`, `.claude/hooks`, `harness/mk`, `harness/README.md`,
`harness/conventions.md`); keine Spezifikation (sie ändert `slice-harness-gate-kette-schnell`), keine
Nutzer- oder Wartungs-Doku, kein Produkt-Code, kein Bootstrap. **Sendung des Gebers
`slice-harness-gate-kette-schnell`** (die Namen aus §1, die Randformen der Hook-Seite aus §6):
nachgezählt mit der Sendung 3 Liefer-Punkte und 1 Schicht; die Namen sind Teil der Liefer-Punkte 1 bis
3, kein vierter.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
