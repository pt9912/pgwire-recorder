# Slice slice-harness-gate-kette-stop-hook-liste: Dateiliste als Nachweis, der Stop-Hook gibt nur bei identischer Liste frei

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
Nutzers vom 2026-10-12 (Ziel: ein Release; die Gate-Kette steht nach dem Release, hinter
`welle-erster-release` und `welle-v1-nachtrag`): direkt hinter `slice-harness-gate-kette-schnell`
(WIP-Limit 1); Reihenfolge in §4 *Start* und im Drift-Log der Roadmap. Der Zuschnitt der Hook-Arbeit in
`slice-harness-gate-kette-schnell`, diesen Slice und `slice-harness-gate-kette-stop-hook-freigabe`
(vor dem Start geteilt; der Plan hieß bis zum 2026-10-12 `slice-harness-gate-kette-stop-hook`) und die
Reihenfolge sind **Entscheidung des Nutzers vom 2026-10-12** („Deine Empfehlungen“).

**Bezug:** kein Anspruch des Produkts — der Slice liefert ein Werkzeug des Harness.
Bindung an Entscheidungen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Nachweis als
Dateiliste, Hook; Accepted 2026-10-12), das sich auf
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
`.claude/hooks/stop-require-gates.sh` gibt **nur frei, wenn die Dateiliste des Arbeitsbaums mit der
Dateiliste des letzten vollen grünen Laufs identisch ist**; jeder andere Fall verlangt `make gates`
(fail-closed). Dazu gehört **ein fehlender Stand des vollen Laufs, auch bei sauberem Baum**: Das ist eine
Änderung der Zusage des Hooks (bisher Freigabe, künftig Sperre), und ein frischer Klon oder Worktree
sperrt den ersten Stop, bis `make gates` gelaufen ist. **Entscheidung des Nutzers vom 2026-10-12**
(„Deine Empfehlungen“; [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) Punkt 3 wörtlich); der
Kopf-Kommentar des Hooks (heute „Restlücke … CI ist das Netz“) wird umgeschrieben und beschreibt das neue
Verhalten (`AGENTS.md` §3.7). Die Gegenprobe `stop-hook-gegenprobe` macht je Randform dieses Slice (§6)
einen Fall, der nur im inerten Fall freigibt. Der Hook ist das Restrisiko der Entscheidung
([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Konsequenzen*): Ohne seine Gegenprobe
schließt der Slice nicht.

**Zuschnitt (Entscheidung des Nutzers vom 2026-10-12, nach dem Schnitt des Architect aus der Prüfung vor
dem Start):** Der Plan hieß ungeteilt `slice-harness-gate-kette-stop-hook`; der Architect schätzte ihn
auf 700 bis 950 Zeilen mit Tests (Logik der Dateiliste 100 bis 150, Hook 130 bis 180, Gegenprobe 450 bis
650) und empfahl den Schnitt vor dem Start. Dieser Slice liefert die **Liste und einen Hook ohne
Positivliste**: Er ist einzeln lieferbar, denn der Hook verlangt nach jeder Änderung `make gates` — das
Verhalten, das `SPEC-052` Hook-Teil für die Zeit vor dem Gate `stop-hook-gegenprobe` nennt — und schließt
die Restlücke. `slice-harness-gate-kette-stop-hook-freigabe` liefert die Positivliste mit der Freigabe.
**Größe, Verteilung der Schätzung des Architect (Planner, nicht gemessen; der Architect bestätigt):**
dieser Slice 410 bis 590 Zeilen mit Tests (Logik der Liste 100 bis 150, Hook 60 bis 90, Gegenprobe mit
rund anderthalb Dutzend Fällen 250 bis 350); der Folge-Slice 260 bis 390.

**Voraussetzung (Geber):** `slice-harness-gate-kette-schnell` liegt in `done/`: `SPEC-052` mit allen
Randformen, das Ziel der schnellen Kette und ihr Stand. **Namen** (vom Architect am 2026-10-12 im Geber
festgelegt, hier eingetragen; kein Eintrag bei dessen Closure mehr): Ziel der schnellen Kette
`gates-schnell`, Stand der schnellen Kette `.harness/state/gates-schnell-passed.diffsha` (Hash des
Arbeitsbaums, von diesem Slice nicht gelesen), Gegenprobe der Kette `gates-schnell-gegenprobe`, Gegenprobe
des Hooks `stop-hook-gegenprobe` (dieser Slice legt sie an, der Folge-Slice erweitert sie). Die
Dateiliste des vollen Laufs liegt in einer **neuen** Datei unter `.harness/state/` (Name wählt der
Implementer dieses Slice, nicht `gates-passed.diffsha`): Ein alter Hash wird so nie als Liste gelesen.
`SPEC-052` Teil *Hook* gilt nach seinem Eingangssatz, soweit `stop-hook-gegenprobe` die Zusagen prüft;
dieser Slice liefert das Gate und ändert die Spezifikation nicht.

**Sendung an `slice-harness-gate-kette-stop-hook-freigabe`** (eingetragen im selben Commit wie dieser
Plan, §3.13): die Dateiliste mit ihrem Leser (der Hook dort liest die Liste über dieselbe Funktion), das
Gate `stop-hook-gegenprobe` mit dem Gerüst der Temp-Repos, die Prüfung nicht aufgelöster Index-Einträge
vor jeder Entscheidung nach Pfad, und der Name der Dateiliste. Den Namen wählt der Implementer; der Planner
trägt ihn bei der Closure dieses Slice in den Nehmer ein (§6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `SPEC-052`, die schnelle Kette, ihr Stand und deren Gegenprobe — `slice-harness-gate-kette-schnell`
  (liegt vor diesem Slice in `done/`): Dieser Slice setzt sie voraus und ändert die Spezifikation nicht;
  gibt der Implementer eine Randform zurück, entscheidet der Architect in `SPEC-052`.
- Positivliste (`docs/plan/planning/`, `docs/reviews/`), Grenze, Symlink-Regel mit Ziel und Spiegel,
  Freigabe bei geänderten inerten Pfaden und das Lesen des Stands der schnellen Kette —
  `slice-harness-gate-kette-stop-hook-freigabe` (dort §1 und DoD; es nimmt an): Dieser Slice kennt nur
  „identisch“ oder „voll“; die Pfad-Logik ist der zweite Lieferwert und braucht ihre eigenen Fälle.
- Aufnahme weiterer Pfade in die Positivliste (`docs/plan/adr/`, `spec/`, `docs/user/`, `README.md`) —
  Bestand bleibt bewusst stehen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Setzungen des
  Nutzers*; eine Aufnahme ist eine Änderung von `SPEC-052` mit dem Beleg, dass kein langsames Gate den
  Pfad liest. Es gibt dafür keinen Slice und keine Kennung.
- Änderung von `make gates` — [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Entscheidung*
  Punkt 1: Die volle Kette bleibt der Nachweis für alles, was nicht auf der Positivliste steht.
- Die Hooks `.claude/hooks/pretooluse-command-guard.sh` und `.claude/hooks/span-emit.sh` und den
  Eintrag in `.claude/settings.json` — anderer Vorgang: Der Slice ändert den Nachweis und die Freigabe
  im Stop-Hook, nicht den Befehls-Wächter oder die Erfassung.
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
      Arbeitsbaum; beide teilen sich die Logik, keine Dopplung zwischen Schreiber und Hook. Die
      Hash-Ausgabe für `gates-schnell` bleibt.
- [ ] **Liefer-Punkt 2:** `.claude/hooks/stop-require-gates.sh` gibt nur bei identischer Liste frei und
      verlangt in jedem anderen Fall `make gates` (fehlender oder beschädigter Stand auch bei sauberem
      Baum, nicht lesbarer Vergleich, nicht aufgelöster Index-Eintrag, Abbruch als Antwort `block`);
      der Schleifenschutz (`stop_hook_active`) bleibt; der Kopf-Kommentar beschreibt das neue Verhalten
      (`AGENTS.md` §3.7), nicht die Restlücke.
- [ ] **Liefer-Punkt 3:** Die Gegenprobe `stop-hook-gegenprobe` macht je Randform aus §6 einen Fall in
      einem Temp-Repo, der nur im inerten Fall freigibt, und hängt an `GATE_CHECKS`;
      `harness/README.md` §Sensors führt ihre Zeile mit Bindung an
      [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md); `harness/conventions.md` (Träger-Tabelle,
      Zeile `grundlagen-durchsetzungsschicht.md`) nennt die Träger des geänderten Hooks. Die Gegenprobe
      heißt `stop-hook-gegenprobe` (ein anderer Name als `make hook-gegenprobe`, die dem Commit-Träger
      `.githooks/commit-msg` gehört).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). Der Reviewer prüft
      fail-closed ausdrücklich: Jede Abweichung von „Liste identisch“ verlangt `make gates`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Die Mutationen und ihre roten Fälle stehen als Tabelle in §7 (Zusage · Mutation · roter
      Fall), vom Implementer selbst gefahren (`AGENTS.md` §3.10).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: eine weitere Datei
      `evidence/slice-harness-gate-kette-stop-hook-liste.md` in
      `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar/` (dritter Beleg: Die Prüfung des
      Architect vor dem Start fand den ungeteilten Plan über der Grenze); der Zähler folgt aus den
      Dateien (3×) und der Eintrag bekommt seinen Ausgang. Die Messung der Laufzeit liefert der
      Folge-Slice `slice-harness-gate-kette-stop-hook-freigabe`.
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
| `tools/harness/working-tree-hash.sh` | update / ersetzt | liefert dieselbe Liste für den Arbeitsbaum (Form nach `SPEC-052`: NUL-getrennte Felder Art, Wert, Pfad); bleibt für den Stand der schnellen Kette eine Hash-Ausgabe, die genau dann wechselt, wenn sich der Arbeitsbaum ändert, denn `gates-schnell` ruft sie auf; Schreiber und Hook teilen sich die Logik |
| `.claude/hooks/stop-require-gates.sh` | update | Freigabe nur bei identischer Liste, Rückfall auf voll bei jedem anderen Fall; Kopf-Kommentar neu |
| `tools/harness/` (neue Gegenprobe des Hooks) und ihr Fragment unter `harness/mk/` | neu | je Randform dieses Slice ein Fall in einem Temp-Repo; hängt an `GATE_CHECKS` |
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
dieser Slice steht direkt dahinter und vor `slice-harness-gate-kette-stop-hook-freigabe` und
`slice-harness-commit-struktur-id` (Drift-Log der Roadmap vom 2026-10-12); `slice-harness-gate-kette-schnell`
steht seinerseits nach `welle-erster-release` und `welle-v1-nachtrag`. Die Reihenfolge nach dem Release ist
Entscheidung des Nutzers vom 2026-10-12 (Ziel: ein Release). Die offenen Rückfragen an den Architect (Merge-Konflikt zu `slice-harness-gate-kette-stop-hook-liste`, Eingangssatz von `SPEC-052`, neue Mutationen, Rückführungsgrenzen) bleiben offen, bis die Slices der Gate-Kette starten; sie halten das Release nicht auf. Die Namen aus §1 stehen im Plan. Die Randformen aus §6 sind in `SPEC-052`
festgelegt (der Architect hat sie am 2026-10-12 in §6 des Gebers entschieden, dessen erster Commit
überträgt sie); dieser Slice ändert sie nicht. **Vor dem Start:** Der Architect bestätigt die Verteilung der
Randformen und Mutationen auf diesen Slice und `slice-harness-gate-kette-stop-hook-freigabe` und die
Größenschätzung aus §1 (Planner-Verteilung seiner Schätzung).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff liegt über einer Review-Sitzung
  (über 650 Zeilen mit Tests; die Schätzung liegt bei 410 bis 590). Schnitt dann nach Lieferwert:
  Dateiliste mit Schreiber und ihrer Gegenprobe (Format, beschädigte Liste) zuerst, der Vergleich im
  Hook mit seinen Fällen danach; beide mit der Gegenprobe ihrer Zusagen, denn ohne Gegenprobe schließt
  keiner.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer nimmt
  [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) zurück, oder das Format der Antwort des
  Stop-Hooks ändert sich so, dass die Gegenprobe nicht mehr gegen die echte Eingabe und Ausgabe des
  Hooks läuft.
## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (vom Auftraggeber gefahren), die Gegenprobe des Hooks rot unter
jeder Mutation der Tabelle in §6, Closure-Notiz mit Lerneintrag; der Planner trägt bei der Closure den
Namen der Dateiliste in `slice-harness-gate-kette-stop-hook-freigabe` ein (§6).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12): Der Slice liefert einen neuen Vertrag (Nachweis als Liste, Freigabe
nur bei Identität). Die Randformen sind vom Architect am 2026-10-12 **entschieden**, und zwar in der
Tabelle von §6 des Gebers `slice-harness-gate-kette-schnell` (Fundstelle `SPEC-052`; dessen erster Commit
überträgt sie in die Spezifikation, vor jedem Code). Dieser Slice trägt davon, ohne Doppelung mit
`slice-harness-gate-kette-stop-hook-freigabe`:

- Form der Dateiliste (NUL-getrennt, Art/Wert/Pfad; Arten `F`, `L`, `G`, `X`; Pfad mit Zeilenumbruch,
  Leerzeichen, Nicht-UTF-8; sortiert, ein Pfad einmal); die Liste enthält Links mit Ziel, das Auswerten
  des Ziels gehört dem Folge-Slice;
- Vergleich der Listen: ein Eintrag, der fehlt, neu ist oder abweicht (Inhalt, Link-Ziel, gelöschte,
  umbenannte, untracked Datei), verlangt `make gates`; die Pfad-Entscheidung („inert“) gehört dem
  Folge-Slice;
- leerer Diff gegen den vollen Lauf: Identität gibt frei (ohne Stand der schnellen Kette, die dieser
  Hook nicht liest); der Fall ohne Stand der schnellen Kette bei geändertem inertem Pfad ist Randform des
  Folge-Slice;
- erster Lauf ohne Stand des vollen Laufs: `make gates`, auch bei sauberem Baum (Entscheidung des Nutzers
  vom 2026-10-12, §1);
- nicht lesbarer Vergleich; Abbruch des Hooks (Antwort `block`, nie Abbruch ohne Antwort);
- Datei der Dateiliste gelöscht oder beschädigt (leer, Feldzahl nicht durch 3 teilbar, unbekannte Art,
  `F` mit Wert, der nicht 64 Hexzeichen ist, `G` oder `X` mit Wert, doppelter Pfad); ein unzulässiger
  Pfad (Segment `.`, `..` oder leer, absolut, NUL) in der Liste ist Randform des Folge-Slice, weil nur
  die Positivliste ihn liest;
- Merge-Konflikt: ein nicht aufgelöster Index-Eintrag (`git ls-files -u` nicht leer) verlangt
  `make gates`, vor jeder Entscheidung nach Pfad (hier verteilt vom Planner, siehe unten);
- Schleifenschutz (`stop_hook_active`) bleibt;
- gleichzeitige Läufe (**Aufruf-Voraussetzung**: nacheinander, `AGENTS.md` §3.12): `SPEC-052` bindet jede
  Zusage daran, der Hook sagt gleichzeitige Läufe nicht zu.

Eine nicht genannte oder offene Randform entscheidet der Implementer nicht, er gibt sie zurück.

**Mutationen** (`AGENTS.md` §3.10): Der Implementer fährt je Zusage die Mutation selbst und trägt die
Tabelle in §7 ein. Die Zeilen mit „Planner“ in der letzten Spalte sind vom Planner abgeleitet, der
Architect bestätigt sie vor dem Start:

| Zusage | Mutation | roter Fall | Herkunft |
|---|---|---|---|
| Fehlender oder unlesbarer Stand des vollen Laufs verlangt `make gates` | fehlenden Stand als „nichts geändert“ lesen | Hook gibt frei → rot | Architect |
| Fehlender Stand bei sauberem Baum verlangt `make gates` (Änderung der Zusage, Entscheidung des Nutzers) | den Fall wie bisher freigeben | Hook gibt frei → rot | Architect |
| Beschädigte Dateiliste verlangt `make gates` | Zeile ohne Hash überspringen; je Fehlform der Aufzählung oben ein Fall | Hook gibt frei → rot | Architect |
| Ein nicht lesbarer Vergleich verlangt `make gates` | Fehler von `git` ignorieren | Hook gibt frei → rot | Architect |
| Ein Abbruch des Hooks antwortet `block` | Fehlerbehandlung entfernen (Abbruch ohne Antwort) | Hook endet ohne `block` → rot | Architect |
| Der Schleifenschutz bleibt | `stop_hook_active` nicht mehr lesen | Hook blockiert endlos → rot | Architect |
| `make gates` schreibt die Dateiliste nur bei grünem Lauf | Schreiben vor die Gates ziehen | roter Lauf hinterlässt Liste → rot | Architect |
| Ein nicht aufgelöster Index-Eintrag verlangt `make gates` | `git ls-files -u` nicht lesen | Konflikt bei sonst identischer Liste gibt frei → rot | Planner (Architect: Konflikt in inertem Pfad, Fall dort nicht möglich ohne Positivliste) |
| Jede Abweichung der Liste verlangt `make gates` | nur Pfade oder nur Werte vergleichen | geänderter Inhalt, gelöschte, hinzugekommene, umbenannte Datei, geändertes Link-Ziel gibt frei → rot | Planner |
| Eine identische Liste gibt frei | Hook blockiert immer | Gutfall blockiert → rot | Planner |
| Die Liste trägt jeden Pfad (Zeilenumbruch, Leerzeichen) | zeilenweise statt NUL-getrennt lesen | Pfad mit Zeilenumbruch → rot | Planner |

**Verteilung auf den Folge-Slice** (`slice-harness-gate-kette-stop-hook-freigabe`, dort Fälle und Zeilen):
Positivliste, Grenze, Groß-/Kleinschreibung, Pfade mit `.`/`..`/leer/absolut, Symlink mit Ziel und
Spiegelregel, Stand der schnellen Kette (fehlt, rot, anderer Arbeitsbaum), leerer Diff ohne Stand der
schnellen Kette; Mutationen „nur Positivliste“, „Grenze“, „beide Pfade der Umbenennung“, „Symlink mit
Ziel“, „Spiegelregel“, „Segmente“, „Stand der schnellen Kette“, „leerer Diff“.

**Risiken:**

- Der Hook ist das Restrisiko der Entscheidung: Ein Fehler dort bleibt ungeprüft, wo die Gegenprobe
  ihn nicht fängt ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Konsequenzen*) — **Ausgang:** entfallen: Die Gegenprobe macht je
  Randform einen Fall, jede Mutation der Tabelle ist rot gesehen, und ohne sie schließt der Slice
  nicht (DoD-Punkt 3).
- Die Gegenprobe fährt den Hook in Temp-Repos mit dessen Eingabe, nicht in der Umgebung von Claude
  Code; ob die Umgebung die Antwort (`decision`) so liest, prüft sie nicht — **Ausgang:** weiter offen:
  → `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (neu bei der Closure angelegt, 1×).
- Zwei Läufe gleichzeitig (der Hook liest, während `make gates` die Liste schreibt) — **Ausgang:**
  entfallen: `SPEC-052` bindet die Zusagen an Läufe nacheinander (Aufruf-Voraussetzung, entschieden
  am 2026-10-12, `AGENTS.md` §3.12), und der Hook sagt gleichzeitige Läufe nicht zu; eine halb
  geschriebene Liste endet bei voll (beschädigt oder „kein Stand“), ein atomares Schreiben bringt keinen
  neuen Schutz.
- Der Hook sperrt nach der Änderung der Zusage (fehlender Stand bei sauberem Baum) den ersten Stop in
  einem frischen Klon oder Worktree — **Ausgang:** eingetreten ist es mit dem Slice selbst (gewollt,
  Entscheidung des Nutzers vom 2026-10-12, [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md)
  Punkt 3); der Kopf-Kommentar des Hooks sagt es.
- Der Name der Dateiliste wählt der Implementer, der Folge-Slice liest sie über dieselbe Funktion —
  **Ausgang:** eingetreten: `slice-harness-gate-kette-stop-hook-freigabe` (nimmt an: §1 und DoD nennen die
  Liste, die Funktion und `stop-hook-gegenprobe`; der Planner trägt den Namen bei der Closure dieses Slice
  ein, mit Nachzählung des Nehmers im selben Commit; bis `grep` nach dem Namen im Nehmer trifft, ist das
  Risiko ohne Ausgang).
- Ein unzulässiger Pfad in der Liste (Segment `..`, absolut) wird in diesem Slice nicht gelesen, weil
  der Hook ohne Positivliste jeden Unterschied sperrt — **Ausgang:** eingetreten:
  `slice-harness-gate-kette-stop-hook-freigabe` (nimmt an: Randform und Mutation „Segmente“ stehen dort in
  §6 und DoD).
- Die Größe zeigte sich erst in der Prüfung des Architect (ungeteilt 700 bis 950 statt 500 bis 750) —
  **Ausgang:** weiter offen: → `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×; dies ist
  der dritte Beleg, die Evidence legt diese Closure an, DoD).
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

- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (2×, offen) — Gegenstand der Hook-Arbeit; die Messung
  (Freigabe eines Stops nach einem Doku-Commit) setzt die Positivliste voraus und liegt deshalb bei
  `slice-harness-gate-kette-stop-hook-freigabe`, das den dritten Beleg und den Ausgang liefert. Dieser
  Slice liefert keine Messung der Ersparnis.
- `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (1×, offen) — berührt: nicht lesbarer Eintrag und
  beschädigte Dateiliste sind Randformen dieses Slice (§6); übergeht der Hook einen Eintrag still, ist es
  das zweite Auftreten.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — berührt vor allem der Folge-Slice
  (Positivliste); hier nur der Kommentar des Hooks, der nur zusagt, was die Gegenprobe prüft.
- `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) — berührt: §1 Schnitt vor dem
  Start; dieser Slice legt bei der Closure den dritten Beleg an (DoD).
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×, offen) — der Slice ist wellenlos; die
  Randformen entscheidet der Architect vor dem Code des Gebers, §4 *Start* nennt es.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) und
  `BEO-REPO/slice-waechst-durch-uebernahmen` (verkörpert, `AGENTS.md` §3.13) — Mutationen in §6, kein
  Kommentar sagt mehr zu, als die Gegenprobe prüft, die Sendungen sind nachgezählt.

Kein Eintrag erreicht mit diesem Plan die Schwelle 3× neu; der dritte Beleg zu
`BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` entsteht bei der Closure.

**Zählung nach `AGENTS.md` §3.13:** Liefer-Punkte: 3 (Dateiliste; Hook mit Freigabe nur bei identischer
Liste; Gegenprobe mit den Gate-Zeilen). Messung, Gate-Läufe, Review, Closure-Notiz und Register zählen
nicht. Schichten: 1 — Harness-Werkzeug (`tools/harness`, `.claude/hooks`, `harness/mk`,
`harness/README.md`, `harness/conventions.md`); keine Spezifikation (sie ändert
`slice-harness-gate-kette-schnell`), keine Nutzer- oder Wartungs-Doku, kein Produkt-Code, kein Bootstrap.
Nachgezählt am 2026-10-12 bei der Umordnung (nach dem Release): Die Sendungen von und an die Nachbarn sind
unverändert, die Zählung bleibt 3 Liefer-Punkte und 1 Schicht.
**Sendung des Gebers `slice-harness-gate-kette-schnell`** (die Namen aus §1, die Randformen der Hook-Seite
aus §6, soweit hier verteilt): nachgezählt mit der Sendung 3 Liefer-Punkte und 1 Schicht; die Namen sind
Teil der Liefer-Punkte 1 bis 3, kein vierter. **Sendung an `slice-harness-gate-kette-stop-hook-freigabe`**
(§1): dort nachgezählt, im selben Commit wie dieser Plan. Geschätzt 410 bis 590 Zeilen mit Tests (§1,
Verteilung der Schätzung des Architect durch den Planner); in einer Review-Sitzung prüfbar.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
