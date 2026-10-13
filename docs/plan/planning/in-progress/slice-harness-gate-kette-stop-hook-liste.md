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
Nutzers vom 2026-10-13 (die Gate-Kette steht vor den übrigen Release-Slices, vor
`slice-v1-abschluss-tls-client`): direkt hinter `slice-harness-gate-kette-schnell`
(WIP-Limit 1); Reihenfolge in §4 *Start* und im Drift-Log der Roadmap. Der Zuschnitt der Hook-Arbeit in
`slice-harness-gate-kette-schnell`, diesen Slice und `slice-harness-gate-kette-stop-hook-freigabe`
(vor dem Start geteilt; der Plan hieß bis zum 2026-10-12 `slice-harness-gate-kette-stop-hook`) und die
Reihenfolge sind **Entscheidung des Nutzers vom 2026-10-12** („Deine Empfehlungen“).

**Bezug:** kein Anspruch des Produkts — der Slice liefert ein Werkzeug des Harness.
Bindung an Entscheidungen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Nachweis als
Dateiliste, Hook; Accepted 2026-10-12), das sich auf
[`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) stützt.

**Berührte Spec-Stellen:** `SPEC-052` (von `slice-harness-gate-kette-schnell` in
`spezifikation.md` §11 angelegt, mit den Randformen vom Architect entschieden; die Punkte 13 und 20 bis 22
hat der Architect am 2026-10-13 aus den Rückgaben dieses Slice ergänzt; dieser Slice ändert sie nicht)

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
einen Fall, der nur im inerten Fall freigibt. Zum Nachweis gehört das Make-Flag `-i` (`SPEC-052` Punkt 21: ein
Lauf mit überlesenen Fehlern schreibt keinen Nachweis, ein älterer entfällt) und zum Hook die Aufrufart
(Punkt 22: Arbeitsverzeichnis im Unterverzeichnis, leere Umgebung). Der Hook ist das Restrisiko der Entscheidung
([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Konsequenzen*): Ohne seine Gegenprobe
schließt der Slice nicht.

**Zuschnitt (Entscheidung des Nutzers vom 2026-10-12, nach dem Schnitt des Architect aus der Prüfung vor
dem Start):** Der Plan hieß ungeteilt `slice-harness-gate-kette-stop-hook`; der Architect schätzte ihn
auf 700 bis 950 Zeilen mit Tests (Logik der Dateiliste 100 bis 150, Hook 130 bis 180, Gegenprobe 450 bis
650) und empfahl den Schnitt vor dem Start. Dieser Slice liefert die **Liste und einen Hook ohne
Positivliste**: Er ist einzeln lieferbar, denn der Hook verlangt nach jeder Änderung `make gates` — das
Verhalten, das der Hook heute für jede Änderung zeigt — und schließt
die Restlücke. `slice-harness-gate-kette-stop-hook-freigabe` liefert die Positivliste mit der Freigabe.
**Größe, Verteilung der Schätzung des Architect (Planner, nicht gemessen; der Architect bestätigt):**
dieser Slice 410 bis 590 Zeilen mit Tests (Logik der Liste 100 bis 150, Hook 60 bis 90, Gegenprobe mit
rund anderthalb Dutzend Fällen 250 bis 350); der Folge-Slice 260 bis 390. **Architect, 2026-10-13:** Die Verteilung stimmt; durch die geteilten Mutationen (§6, Abweichung der Liste in vier Fällen, Leerzeichen, Konflikt bei identischer Liste) und die sechs Fehlformen der beschädigten Liste liegt die Gegenprobe bei rund 25 Fällen, tabellengetrieben 280 bis 400; dieser Slice **440 bis 620** Zeilen mit Tests, der Folge-Slice 280 bis 420.

**Voraussetzung (Geber):** `slice-harness-gate-kette-schnell` liegt in `done/`: `SPEC-052` mit allen
Randformen, das Ziel der schnellen Kette und ihr Stand. **Namen** (vom Architect am 2026-10-12 im Geber
festgelegt, hier eingetragen; kein Eintrag bei dessen Closure mehr): Ziel der schnellen Kette
`gates-schnell`, Stand der schnellen Kette `.harness/state/gates-schnell-passed.diffsha` (Hash des
Arbeitsbaums, von diesem Slice nicht gelesen), Gegenprobe der Kette `gates-schnell-gegenprobe`, Gegenprobe
des Hooks `stop-hook-gegenprobe` (dieser Slice legt sie an, der Folge-Slice erweitert sie). Die
Dateiliste des vollen Laufs liegt in einer **neuen** Datei unter `.harness/state/`, neben
`gates-passed.diffsha`: `.harness/state/gates-passed.dateiliste` (Name vom Implementer gewählt): Ein alter Hash
wird so nie als Liste gelesen.
`SPEC-052` Teil *Hook* gilt nach seinem Eingangssatz (Wortlaut: `slice-harness-gate-kette-schnell` §6, Zeile *Zwischenstand*), soweit ein Fall von `stop-hook-gegenprobe` die Zusage nennt;
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
- Positivliste (`docs/plan/planning/`, `docs/reviews/`, `docs/user/`, `README.md`; Architect, 2026-10-13, Beleg in
  `slice-harness-gate-kette-schnell` §6), Grenze, Symlink-Regel mit Ziel und Spiegel,
  Freigabe bei geänderten inerten Pfaden und das Lesen des Stands der schnellen Kette —
  `slice-harness-gate-kette-stop-hook-freigabe` (dort §1 und DoD; es nimmt an): Dieser Slice kennt nur
  „identisch“ oder „voll“; die Pfad-Logik ist der zweite Lieferwert und braucht ihre eigenen Fälle.
- Aufnahme weiterer Pfade in die Positivliste (`docs/plan/adr/`, `spec/`) —
  Bestand bleibt bewusst stehen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Setzungen des
  Nutzers*; eine Aufnahme ist eine Änderung von `SPEC-052` mit dem Beleg, dass kein langsames Gate den
  Pfad liest. `docs/user/` und `README.md` sind seit der Entscheidung des Nutzers vom 2026-10-13 mit Beleg
  aufgenommen (`slice-harness-gate-kette-schnell` §6); die Positivliste liefert
  `slice-harness-gate-kette-stop-hook-freigabe`, nicht dieser Slice. Es gibt dafür keinen Slice und keine Kennung.
- Änderung von `make gates` — [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Entscheidung*
  Punkt 1: Die volle Kette bleibt der Nachweis für alles, was nicht auf der Positivliste steht.
- Die Hooks `.claude/hooks/pretooluse-command-guard.sh` und `.claude/hooks/span-emit.sh` und den
  Eintrag in `.claude/settings.json` — anderer Vorgang: Der Slice ändert den Nachweis und die Freigabe
  im Stop-Hook, nicht den Befehls-Wächter oder die Erfassung.
- `bash -e` und eine fehlende Variable `CLAUDE_PROJECT_DIR` als Aufrufart des Hooks — nicht Gegenstand
  (`SPEC-052` Punkt 22): der Hook setzt seine Optionen selbst, und der Aufruf des Treibers startet ihn ohne die
  Variable nicht; kein Code des Hooks läuft, er kann nicht antworten.
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
      Hash-Ausgabe für `gates-schnell` bleibt. Ein Lauf mit `-i` in `MAKEFLAGS` schreibt weder Liste noch
      Hash, ein älterer Nachweis beider entfällt; `-n`, `-q`, `-t` lassen ihn stehen (`SPEC-052` Punkt 21).
- [ ] **Liefer-Punkt 2:** `.claude/hooks/stop-require-gates.sh` gibt nur bei identischer Liste frei und
      verlangt in jedem anderen Fall `make gates` (fehlender oder beschädigter Stand auch bei sauberem
      Baum, nicht lesbarer Vergleich, nicht aufgelöster Index-Eintrag, Abbruch als Antwort `block`);
      der Schleifenschutz (`stop_hook_active`) bleibt; der Kopf-Kommentar beschreibt das neue Verhalten
      (`AGENTS.md` §3.7), nicht die Restlücke, und sagt die Sperre bei fehlendem Stand zu, auch bei sauberem
      Baum (Entscheidung des Nutzers vom 2026-10-12); der Fall der Gegenprobe dazu steht in Liefer-Punkt 3.
      Die Prüfung nicht aufgelöster Index-Einträge ist die erste Prüfung nach dem Schleifenschutz, auch vor
      „Liste identisch“. Der Hook urteilt über das Repo seines Arbeitsverzeichnisses und liest keine
      Umgebungsvariable (`SPEC-052` Punkt 22).
- [ ] **Liefer-Punkt 3:** Die Gegenprobe `stop-hook-gegenprobe` macht je Randform aus §6 einen Fall in
      einem Temp-Repo, der nur im inerten Fall freigibt (dazu `p21-*` für die Make-Flags des Nachweises und
      `p22-*` für die Aufrufarten), und hängt an `GATE_CHECKS`;
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
| `tools/harness/record-gates.sh` | update | schreibt die Dateiliste statt eines Hashs, nur nach grünem `make gates`; liest `MAKEFLAGS` wie `gates-schnell.sh`: bei `i` keine Schreibung, ein älterer Nachweis entfällt (`SPEC-052` Punkt 21) |
| `tools/harness/working-tree-hash.sh` | update / ersetzt | liefert dieselbe Liste für den Arbeitsbaum (Form nach `SPEC-052`: NUL-getrennte Felder Art, Wert, Pfad); bleibt für den Stand der schnellen Kette eine Hash-Ausgabe, die genau dann wechselt, wenn sich der Arbeitsbaum ändert, denn `gates-schnell` ruft sie auf; Schreiber und Hook teilen sich die Logik |
| `.claude/hooks/stop-require-gates.sh` | update | Freigabe nur bei identischer Liste, Rückfall auf voll bei jedem anderen Fall; Kopf-Kommentar neu |
| `tools/harness/dateiliste.sh` | neu | die Funktionen der Liste (`dateiliste_erzeugen`, `dateiliste_pruefen`) und der Pfad `.harness/state/gates-passed.dateiliste`; `record-gates.sh`, `working-tree-hash.sh` und der Hook laden sie mit `source` (eine Logik, keine Dopplung) |
| `tools/harness/stop-hook-gegenprobe.sh` und `harness/mk/stop-hook-gegenprobe.mk` | neu | je Randform dieses Slice ein Fall in einem Temp-Repo, darunter `p21-*` (Make-Flags des Nachweises) und `p22-*` (Aufrufarten); hängt an `GATE_CHECKS` |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh` | update | `stop-hook-gegenprobe` endet auf `-gegenprobe` und steht damit in der Ausnahmeliste (`SPEC-052` Punkt 2 und 4); die Kopien der Gegenprobe der Kette tragen `dateiliste.sh` und führen den Namen in ihrer Attrappen-Liste |
| `harness/README.md` §Sensors | update | Zeile der Gegenprobe, Bindung [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) |
| `harness/conventions.md` (Träger-Tabelle, Zeile `grundlagen-durchsetzungsschicht.md`) | update | nennt die Träger des geänderten Hooks (Dateiliste, Hook, gegebenenfalls ersetzte Datei) |

- Jeder Fall der Gegenprobe trägt die Nummer des Punktes aus `SPEC-052` in seinem Namen (Architect, 2026-10-13; sonst ist „ein Fall nennt den Punkt“ im Eingangssatz nicht prüfbar).
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
dieser Slice steht direkt dahinter und vor `slice-harness-gate-kette-stop-hook-freigabe`,
`slice-v1-abschluss-tls-client` (Schritt 28) und den übrigen Release-Slices (Drift-Log der Roadmap vom 2026-10-13);
`slice-harness-gate-kette-schnell` steht seinerseits vor `slice-v1-abschluss-tls-client`. Die Reihenfolge ist
Entscheidung des Nutzers vom 2026-10-13. Die Rückfragen des Planners hat der Architect am 2026-10-13 entschieden (§1 Größe, §4 Rückführung, §6 Mutationen, Wortlaut des Eingangssatzes im Geber). Die Namen aus §1 stehen im Plan. Die Randformen aus §6 sind in `SPEC-052`
festgelegt (der Architect hat sie am 2026-10-12 in §6 des Gebers entschieden, dessen erster Commit
überträgt sie); dieser Slice ändert sie nicht. **Vor dem Start:** Der Architect hat die Verteilung der
Randformen und Mutationen auf diesen Slice und `slice-harness-gate-kette-stop-hook-freigabe` und die
Größenschätzung aus §1 am 2026-10-13 bestätigt (mit den Änderungen in §1 und §6).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff liegt über einer Review-Sitzung
  (über 700 Zeilen mit Tests; die Schätzung des Architect liegt bei 440 bis 620). Schnitt dann nach Lieferwert:
  Dateiliste mit Schreiber und ihrer Gegenprobe (Format, beschädigte Liste) zuerst, der Vergleich im
  Hook mit seinen Fällen danach; beide mit der Gegenprobe ihrer Zusagen, denn ohne Gegenprobe schließt
  keiner.
- **Urteil zur Marke (Planner, 2026-10-13, Größe des Code-Commits `b259381`).** *Metrik* (§4 nannte keine; hier
  festgelegt, wie bei `slice-v1-abschluss-sessions`): Summe aus Einfügungen und Löschungen von
  `git diff --numstat 9662e13..HEAD -- tools harness .claude` (Skripte, Fragmente, Hook, Tests und
  Gegenprobe eingeschlossen; `docs/` und `spec/` ausgeschlossen, hier unberührt). *Messung:* 704 Einfügungen,
  74 Löschungen, Summe 778, netto 630 (`stop-hook-gegenprobe.sh` 537, Hook 65/45, `dateiliste.sh` 56,
  `record-gates.sh` 16/6, `working-tree-hash.sh` 12/20, übrige 20). Die Marke von 700 ist nach der Metrik
  **eingetreten** (78 darüber); die Schätzung des Architect (440 bis 620) liegt 158 unter der Messung. Nach
  Einfügungen allein (704) ist sie knapp erreicht, nach netto (630) nicht: die Probe hing an der Lesart.
  *Test-Code zählt mit* („mit Tests“ in dieser Marke); ohne die Gegenprobe läge der Rest bei 241.
  *Prüfbarkeit in einer Sitzung:* ja. Prüfgegenstand sind der Hook (110 Zeilen geändert), die Liste (56) und
  eine tabellengetriebene Gegenprobe, die der Reviewer gegen die Mutationstabelle in §6 liest, nicht Zeile
  für Zeile als Logik; 41 Mutanten sind rot gesehen. *Urteil:* **keine Rückführung, Fortsetzung.** Der Schnitt
  des Vorab-Satzes (Liste zuerst, Vergleich danach) trennte fertigen, grünen Code, der Hook ohne seine
  Fälle wäre ein Gate ohne Gegenprobe, und zwei Reviews desselben Hooks kosten mehr als eine Sitzung. Der
  Reviewer meldet in seiner Summary-Zeile, ob der Diff in einer Sitzung prüfbar war; meldet er nein,
  entscheidet der Planner neu. Belege im Register: §8.
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
- Make-Flags beim Schreiben des Nachweises (`make -i gates` mit rotem Gate, `-n`, `-q`, `-t`, `MAKEFLAGS` aus der
  Umgebung, Kombinationen): entschieden in `SPEC-052` Punkt 21 (Architect, 2026-10-13, siehe unten);
- Aufrufarten des Hooks (`bash -e`, leere Umgebung, Arbeitsverzeichnis in einem Unterverzeichnis oder in einem
  anderen Repo, fehlende `CLAUDE_PROJECT_DIR`): entschieden in `SPEC-052` Punkt 22 (Architect, 2026-10-13, siehe unten).

Eine nicht genannte oder offene Randform entscheidet der Implementer nicht, er gibt sie zurück.

**Mutationen** (`AGENTS.md` §3.10): Der Implementer fährt je Zusage die Mutation selbst und trägt die
Tabelle in §7 ein. Die Zeilen mit „Planner“ in der letzten Spalte sind vom Planner abgeleitet; der
Architect hat sie am 2026-10-13 bestätigt, geschärft oder geteilt (Spalte *Herkunft*). Je Zeile bricht die Mutation nur ihre Zusage:

| Zusage | Mutation | roter Fall | Herkunft |
|---|---|---|---|
| Fehlender oder unlesbarer Stand des vollen Laufs verlangt `make gates` | fehlenden Stand als „nichts geändert“ lesen | Hook gibt frei → rot | Architect |
| Fehlender Stand bei sauberem Baum verlangt `make gates` (Änderung der Zusage, Entscheidung des Nutzers) | den Fall wie bisher freigeben | Hook gibt frei → rot | Architect |
| Beschädigte Dateiliste verlangt `make gates` | Zeile ohne Hash überspringen; je Fehlform der Aufzählung oben ein Fall | Hook gibt frei → rot | Architect |
| Ein nicht lesbarer Vergleich verlangt `make gates` | Fehler von `git` ignorieren | Hook gibt frei → rot | Architect |
| Ein Abbruch des Hooks antwortet `block` | Fehlerbehandlung entfernen (Abbruch ohne Antwort) | Hook endet ohne `block` → rot | Architect |
| Der Schleifenschutz bleibt | `stop_hook_active` nicht mehr lesen | Hook blockiert endlos → rot | Architect |
| `make gates` schreibt die Dateiliste nur bei grünem Lauf | Schreiben vor die Gates ziehen | roter Lauf hinterlässt Liste → rot | Architect |
| Ein nicht aufgelöster Index-Eintrag verlangt `make gates`, vor „identisch“ | `git ls-files -u` nicht lesen, oder erst die Identität prüfen | Konflikt bei sonst identischer Liste gibt frei → rot | Planner, vom Architect bestätigt am 2026-10-13 (der frühere Vermerk „in inertem Pfad nicht möglich“ entfällt: der Fall liegt an der Identität, nicht am Pfad) |
| Ein geänderter Inhalt oder ein geändertes Link-Ziel verlangt `make gates` | nur Pfade vergleichen | geänderter Inhalt; geändertes Link-Ziel gibt frei → rot | Planner, geteilt (Architect, 2026-10-13) |
| Eine Umbenennung verlangt `make gates` | nur Werte (Hashes) vergleichen | Umbenennung bei gleichem Inhalt gibt frei → rot | Planner, geteilt |
| Eine gelöschte Datei verlangt `make gates` | nur Pfade des Arbeitsbaums gegen den Stand prüfen (gelöschte ignorieren) | gelöschte Datei gibt frei → rot | Planner, geteilt |
| Eine untracked Datei verlangt `make gates` | `git ls-files` ohne `--others` | neue untracked Datei gibt frei → rot | Planner, geteilt |
| Eine identische Liste gibt frei | Hook blockiert immer | Gutfall blockiert → rot | Planner, bestätigt |
| Die Liste trägt jeden Pfad (Zeilenumbruch) | zeilenweise statt NUL-getrennt lesen | Pfad mit Zeilenumbruch → rot | Planner, bestätigt |
| Die Liste trägt jeden Pfad (Leerzeichen) | Pfad ungequotet expandieren (Wortteilung) | Pfad mit Leerzeichen → rot | Architect, 2026-10-13 |
| `make -i gates` schreibt keinen Nachweis, ein älterer entfällt (`SPEC-052` Punkt 21) | die `-i`-Prüfung in `record-gates.sh` entfernen | `p21-i-rot` (Attrappe `exit 1`) und `p21-i-gruen` schreiben Liste und Hash → rot | Architect, 2026-10-13 |
| `-i` gilt in jeder Schreibweise von `MAKEFLAGS` | nur das literale `-i` im Wortlaut von `MAKEFLAGS` suchen; nur das erste Wort genau `i` gelten lassen (statt der Buchstaben des ersten Worts) | Gemessen (§7): `make` schreibt `-i` als `i` in `MAKEFLAGS`, ein literales `-i` kommt nie vor; der erste Mutant färbt alle sechs `p21-i-*` rot (auch `p21-i-rot`, anders als der Plan des Architect annahm), der zweite nur `p21-i-ki` (`ik`) | Architect, 2026-10-13; Fassung nach Messung |
| Ein älterer Nachweis entfällt unter `-i` | das Entfernen von Liste und Hash weglassen (nur nicht schreiben) | `p21-i-rot` und `p21-i-gruen` mit vorherigem grünen Lauf: Liste oder Hash bleibt → rot | Architect, 2026-10-13 |
| `-n`, `-q`, `-t` schreiben nichts und entfernen nichts, auch mit `-i` (Vorrang, Punkt 21) | das Rezept in `harness/mk/enforce.mk` mit `+` führen (läuft dann unter `-n`); dazu bei `-q` zusätzlich die Voraussetzungen von `record-gates` im `Makefile` entfernen | Gemessen (§7): mit `+` allein `p21-n`, `p21-t`, `p21-ni` (älterer Nachweis geändert) → rot, `p21-q` bleibt grün (`make -q` hält an der Voraussetzung, das Rezept läuft nicht); mit `+` und ohne Voraussetzungen alle vier rot | Architect, 2026-10-13; Fassung nach Messung |
| Unter `-i` entfällt jede der beiden Dateien | nur den Hash entfernen (Liste bleibt); nur die Liste entfernen (Hash bleibt) | je alle sechs `p21-i-*` rot | Implementer (jede Bedingung der Zusage „Liste und Hash“ ist eine Zeile) |
| Der Hook gilt für das Repo seines Arbeitsverzeichnisses, auch aus einem Unterverzeichnis (Punkt 22 a) | `git rev-parse --show-toplevel` durch `pwd` ersetzen | `p22-unterverzeichnis` (identisch `approve`, geändert `block`) → rot | Architect, 2026-10-13 |
| Der Hook liest keine Umgebung (Punkt 22 b) | `: "${HOME:?}"` vor die erste Prüfung setzen | `p22-leere-umgebung` (`env -i`, identisch → `approve`) → rot | Architect, 2026-10-13 |

Sortierung und Einmaligkeit der Liste haben keine eigene Mutation (**akzeptiertes Negativ**, Architect, 2026-10-13; Grund: Schreiber und Leser erzeugen die Liste mit derselben Funktion, eine abweichende Sortierung macht aus „identisch“ höchstens „voll“, nie „frei“; ein doppelter Pfad in der gelesenen Liste ist beschädigt und hat seinen Fall oben).

**Entschieden vom Architect, 2026-10-13** (Rückgaben des Implementers aus §7; Fundstelle `SPEC-052`):

1. **`make -i gates`.** Für den Nachweis des vollen Laufs gilt dieselbe Regel wie für `gates-schnell`
   (`SPEC-052` Punkt 21): `-i` in `MAKEFLAGS` (jede Schreibweise, auch aus der Umgebung, mit `-k`/`-j`) schreibt
   weder Liste noch Hash, ein älterer Nachweis beider entfällt, auch bei grünem Lauf; `-n`/`-q`/`-t` schreiben
   nichts und entfernen nichts, auch zusammen mit `-i`. Ort: `tools/harness/record-gates.sh` (der Nachweis ist der
   letzte Schritt der Kette, `GATE_CHECKS` bleibt unberührt). `-n`/`-q`/`-t` hält `make` selbst (das Rezept in
   `harness/mk/enforce.mk` trägt weder `$(MAKE)` noch `+` und läuft nicht); ein Skript-Wächter dafür entfällt
   (er wäre unerreichbar und nicht zugesagt), der Fall hält die Eigenschaft. Der Implementer prüft die Form von
   `MAKEFLAGS` im Rezept an `make` selbst (`-j4 -i`, `-ki`, `--ignore-errors`, Umgebung) und liest sie wie
   `tools/harness/gates-schnell.sh`. Fälle in `stop-hook-gegenprobe` (Temp-Repo, Attrappen, `FAIL=fake-b`, vorher
   ein grüner Lauf als älterer Nachweis): `p21-i-rot`, `p21-i-gruen`, `p21-i-env`, `p21-i-ki`, `p21-i-lang`,
   `p21-i-j`, `p21-n`, `p21-q`, `p21-t`, `p21-ni` (`-ni` und `MAKEFLAGS=ni`: Nachweis byte-gleich); Mutationen in der
   Tabelle oben. Wirkung auf `-liste`: Liste und Hash sind nur nach einem Lauf ohne überlesene Fehler da; der
   Hook ändert sich nicht. Wirkung auf `slice-harness-gate-kette-stop-hook-freigabe`: keine, die Fälle dort setzen
   den Stand mit `record-gates.sh` im Normalaufruf. Akzeptiertes Negativ: ein roter Lauf ohne `-i` lässt einen
   älteren Nachweis stehen (Grund in Punkt 21).
2. **Aufrufarten des Hooks** (`SPEC-052` Punkt 22): Zusage mit Fall (a) Arbeitsverzeichnis im Unterverzeichnis
   (`p22-unterverzeichnis`) und (b) leere Umgebung (`p22-leere-umgebung`: `env -i "$BASH" Hook`); Nicht-Gegenstand
   `bash -e` (der Hook setzt die Optionen selbst) und fehlende `CLAUDE_PROJECT_DIR` (der Hook startet nicht, kein
   Code läuft); akzeptiertes Negativ: Arbeitsverzeichnis in einem anderen Repo (Sache des Treibers, Punkt 8; der
   Hook urteilt über das Repo, in dem er läuft; ohne Stand, ohne `dateiliste.sh` oder ohne Repo endet er in `block`,
   gedeckt durch `p19-ohne-stand-*`, `p20-lib-fehlt`, `p20-kein-repo`). Unter `block` fallen damit: git fehlt oder
   scheitert, kein Repo, `dateiliste.sh` fehlt, Stand fehlt oder beschädigt (alle schon Fälle); neu keiner. Der
   Implementer zieht den Kopf-Kommentar von `tools/harness/stop-hook-gegenprobe.sh` nach (dort steht „nicht
   Gegenstand“ für `bash -e`, leere Umgebung, `CLAUDE_PROJECT_DIR`) und lässt `hook()` eine leere Umgebung
   zu. Weicht die Antwort in `env -i` ab, weil Git-Konfiguration der Umgebung die Liste ändert (nur Richtung `block`
   möglich), gibt er den Fall zurück.
3. **Pfad mit `..` in der beschädigten Liste:** bestätigt, er bleibt beim Folge-Slice (Mutation dort). In diesem
   Slice sperrt ein unzulässiger Pfad im Stand über die Nicht-Identität; `SPEC-052` Punkt 20 sagt das, ein eigener
   Fall und eine Mutation entfallen hier (kein Code unterscheidet ihn von jedem anderen Unterschied).
4. **`SPEC-052` Punkt 13:** die Aufzählung ist durch eine Regel ersetzt (jedes Gate der Ausnahmeliste mit
   `-gegenprobe` im Namen außer `lint-gegenprobe` fährt in selbst angelegten Temp-Bäumen). Bestätigt: Die Temp-Repos
   von `stop-hook-gegenprobe` tragen nur Makefile, Fragment, Skripte und Hook (`fixture()`), sie liest `docs/user/`
   und `README.md` nicht; der Beleg gilt.
5. **Zählung (§3.13):** unverändert 3 Liefer-Punkte, 1 Schicht (Make-Flag-Prüfung gehört zu Liefer-Punkt 1
   *Dateiliste*, die Fälle zu 3 *Gegenprobe*). Größe: Code in `record-gates.sh` etwa 10 bis 15 Zeilen, Fälle und
   Mutanten in der Gegenprobe etwa 50 bis 80 Zeilen (zwölf Fälle, sechs Mutanten), Kommentar-Nachzug wenige
   Zeilen; zusammen etwa 60 bis 95 Zeilen über dem Stand von 711 Einfügungen. Die Größe entscheidet der Planner.

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

### Belege des Implementers

**Name der Dateiliste:** `.harness/state/gates-passed.dateiliste` (neben `gates-passed.diffsha`, das
bleibt und ist der Hash dieser Liste). Funktionen: `dateiliste_erzeugen` und `dateiliste_pruefen` in
`tools/harness/dateiliste.sh`; Schreiber (`record-gates.sh`), Hash (`working-tree-hash.sh`) und Hook laden sie mit
`source`. Der Planner trägt den Namen bei der Closure in `slice-harness-gate-kette-stop-hook-freigabe` ein (§6).

**Läufe** (Code-Commit `b259381`, Baum sauber): `make gates` Exit 0, 5 min 45 s; darin `stop-hook-gegenprobe: gruen — 40 Fälle …`
(12 s), `gates-schnell-gegenprobe: gruen`, `d-check: 680 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.0 OK`. Danach der Hook
im Arbeitsbaum: `{"decision":"approve"}`; vor dem Lauf, ohne Liste: `block` („file list missing“). Der Hook ist in einer
Kopie (Scratchpad) entwickelt und erst nach grüner Gegenprobe eingespielt worden; der alte Stand: `git show 9662e13:.claude/hooks/stop-require-gates.sh`
(zurück: `git checkout 9662e13 -- .claude/hooks/stop-require-gates.sh`).

**Lauf der Nacharbeit** (Commit `78b35ab`, Code-Commit `e92db76`, Baum sauber): `make gates` Exit 0, 5 min 39 s; darin
`stop-hook-gegenprobe: gruen — 52 Fälle der Punkte 1, 10, 11, 14, 16 bis 22 …` (22 s), `gates-schnell-gegenprobe: gruen`,
`d-check: 682 Datei(en) geprüft, 0 Befund(e)`, `baseline-verify: v6.18.0 OK`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`. Danach der Hook im
Arbeitsbaum: `{"decision":"approve"}`. Der Hook selbst ist in dieser Nacharbeit unverändert; geändert sind `record-gates.sh` und die Gegenprobe.
**Größe der Nacharbeit** (Metrik des Planners aus §4, `git diff --numstat 9662e13..HEAD -- tools harness .claude`): 845 Einfügungen, 75 Löschungen,
Summe 920, netto 770 (`stop-hook-gegenprobe.sh` 659, `record-gates.sh` 34/6, Hook 65/45, `dateiliste.sh` 56); gegenüber dem Stand 778 sind es +142
(`stop-hook-gegenprobe.sh` +122, `record-gates.sh` +18/0, Fragment-Kommentar und Zeile in `harness/README.md` je 1). Die Erwartung des Architect (+60 bis 95) ist
überschritten: die Gegenprobe wächst um zwölf Fälle mit Hilfsfunktionen und zehn Mutanten, `record-gates.sh` um 18 Zeilen mit Kommentar.

**Befund der Gegenprobe an der unveränderten Kopie (vorher/nachher):** Erster Lauf: `Kopie 'ok' ist nicht grün: p20-git-liste-scheitert
p20-datei-nicht-lesbar` — der Hook antwortete leer, weil die Antwort des Abbruch-Handlers in der Umleitung `> "$aktuell"` des gerade
scheiternden Befehls landete. Behoben mit Dateikennung 3 für jede Antwort (`exec 3>&1`, `>&3`); danach beide Fälle grün.

**Mutationen** (`bash tools/harness/stop-hook-gegenprobe.sh`; Mutant = `sed`/`awk` auf der Kopie von Hook oder Skript im Temp-Repo, vom Skript selbst
gefahren; „rot“ heißt: der Fall meldet sich, „grün bleibt“: der Fall bleibt mit dem Mutanten grün). Die Skripte laufen auf dem Host, kein
Docker-Build-Kontext: der Mutant liegt in der Kopie, die bash liest.

| Zusage | Mutation | roter Fall (grün bleibt) |
|---|---|---|
| Fehlender Stand verlangt `make gates` | `[ -f "$stand" ] \|\| frei` | `p19-ohne-stand-geaendert`, `p19-ohne-stand-sauber` |
| Fehlender Stand bei sauberem Baum verlangt `make gates` | bei sauberem Baum freigeben, sonst `block` | `p19-ohne-stand-sauber` (grün: `p19-ohne-stand-geaendert`) |
| Nicht lesbarer Stand verlangt `make gates` | `[ -r "$stand" ] \|\| frei` | `p20-stand-unlesbar` (läuft nicht als root; hier als Nutzer 1000 gelaufen) |
| Beschädigt: leere Datei | `[ -s ]` entfernt | `p20-leer` |
| Beschädigt: Feldzahl nicht durch 3 | beide `read … \|\| return 1` für Wert und Pfad entfernt; nur das für den Pfad entfernt; Prüfung des unterminierten letzten Felds entfernt | `p20-feldzahl-4`, `p20-feldzahl-5`; `p20-feldzahl-5` (grün: `-4`); `p20-feldzahl-unterminiert` |
| Beschädigt: unbekannte Art | `*) return 1` → `*) ;;` | `p20-art` |
| Beschädigt: `F` mit Wert ohne 64 Hexzeichen | `{64}` → `+`; `[0-9a-fA-F]` → `.` | `p20-hash-kurz` (grün: `-nichthex`); `p20-hash-nichthex` (grün: `-kurz`) |
| Beschädigt: `G` oder `X` mit Wert | nur `G` ohne Prüfung; nur `X` ohne Prüfung | `p20-wert-g` (grün: `-x`); `p20-wert-x` (grün: `-g`) |
| Beschädigt: doppelter Pfad | Prüfung des Pfads entfernt | `p20-doppelt` |
| Nicht lesbarer Vergleich: `git ls-files -u` scheitert | `\|\| true` | `p20-git-u-scheitert` |
| Nicht lesbarer Vergleich: Liste des Arbeitsbaums | `\|\| :` hinter dem Erzeuger im Hook; `pipefail` im Erzeuger aus | `p20-git-liste-scheitert`, `p20-datei-nicht-lesbar`; `p20-git-liste-scheitert` (grün: `p20-datei-nicht-lesbar`) |
| Nicht lesbarer Vergleich: eine Datei nicht lesbar | `\|\| true` hinter `sha256sum` im Erzeuger | `p20-datei-nicht-lesbar` (grün: `p20-git-liste-scheitert`) |
| Ein Abbruch antwortet `block` | `trap sicherung EXIT` entfernt | `p20-abbruch`, `p20-lib-fehlt`, `p20-kein-repo`, `p20-ohne-git` |
| Der Schleifenschutz bleibt | `if false`; Wert `true` nicht gelesen | `p20-schleifenschutz` (grün: `-false`); `p20-schleifenschutz-false` (grün: `p20-schleifenschutz`) |
| Die Liste entsteht nur bei grünem Lauf | `record-gates` vor die Gates gezogen | `p1-liste-nicht-bei-rot` (grün: `p1-liste-bei-gruen`) |
| Ein grüner Lauf schreibt Liste und Hash | `mv` der Liste entfernt; Schreiben des Hashs entfernt | `p1-liste-bei-gruen`; `p1-hash-neben-liste` |
| Merge-Konflikt verlangt `make gates`, vor „identisch“ | `git ls-files -u` nicht gelesen; Prüfung hinter die Identität verschoben | `p17-konflikt-identisch` (beide) |
| Geänderter Inhalt, geändertes Link-Ziel | nur Pfade vergleichen; Werte von `F` ignoriert; Werte von `L` ignoriert | `p14-inhalt`, `p14-link-ziel`; `p14-inhalt` (grün: Link-Ziel, Umbenennung); `p14-link-ziel` (grün: Inhalt, Umbenennung) |
| Umbenennung | nur Werte (sortiert) vergleichen | `p14-umbenennung` (grün: Inhalt, Link-Ziel) |
| Gelöschte Datei | nur neue und abweichende Einträge, `G` ausgenommen | `p14-geloescht`, `p14-geloescht-index` (grün: Inhalt, Umbenennung) |
| Neue Datei (getrackt, untracked, am Eintragsende abgeschnittene Liste) | nur Einträge des Standes gegen den Arbeitsbaum | `p14-neu-getrackt`, `p18-untracked`, `p20-abgeschnitten-eintragsende` (grün: Inhalt, gelöschte Datei) |
| Untracked Datei | `--others` entfernt | `p18-untracked` (grün: `p14-neu-getrackt`) |
| Pfad mit Zeilenumbruch | zeilenweise statt NUL-getrennt | `p11-zeilenumbruch`, `p11-form` (grün: `p11-leerzeichen`) |
| Pfad mit Leerzeichen | `$pfad` ungequotet im Erzeuger | `p11-leerzeichen`, `p11-zeilenumbruch`, `p11-form` (die Wortteilung trifft den Zeilenumbruch mit; die Trennung der beiden leistet der Mutant davor) |
| Form (Art, Wert, Pfad; vier Arten) | Pfad vor Wert; `X` als `G` | `p11-form` |
| Eine identische Liste gibt frei | Hook blockiert immer | `p16-identisch` |
| Der Hash wechselt genau mit dem Arbeitsbaum | Hash konstant; Hash mit Zeitstempel | `p10-hash-wechselt` (grün: `-stabil`); `p10-hash-stabil` (grün: `-wechselt`) |
| Das Gate hängt an `GATE_CHECKS`, sein Exit gilt | Zeile `GATE_CHECKS +=` gestrichen; `-@bash` | `anschluss-gate-checks`; `anschluss-exit` |

**Nacharbeit (Code-Commit `e92db76`, `SPEC-052` Punkte 21 und 22):** Der Nachweis liest `MAKEFLAGS` in `record-gates.sh`;
die Gegenprobe führt 12 Fälle mehr (52 Fälle) und 10 Mutanten mehr. Die Form von `MAKEFLAGS` im Rezept ist an `make` 4.3 gemessen
(Scratchpad, `printf 'x:\n\t@echo "[$$MAKEFLAGS]"'`): `-i` → `i`; `-j4 -i` → `i -j4 --jobserver-auth=3,4`; `-ki` → `ik`;
`--ignore-errors` → `i`; `MAKEFLAGS=i`, `=-i`, `=--ignore-errors` in der Umgebung → `i`; `-ni` und `-i -n` führen das Rezept
nicht aus. Die Fälle `p21-*` fahren `make gates` ohne `-s`, weil `-s` den Buchstaben `s` ins erste Wort schreibt (`is`).
Jeder Fall legt vorher einen grünen Lauf als älteren Nachweis an und ändert danach `a.txt`, damit ein Schreiben den
Nachweis ändert (`p21-n`, `p21-q`, `p21-t`, `p21-ni`: byte-gleich zu dem des grünen Laufs).

| Zusage | Mutation | roter Fall (grün bleibt) |
|---|---|---|
| `-i` schreibt keinen Nachweis, ein älterer entfällt (Punkt 21) | `case`-Zweig `*i*` unwirksam gemacht (`*nie*`) | `p21-i-rot`, `p21-i-gruen`, `p21-i-env`, `p21-i-ki`, `p21-i-lang`, `p21-i-j` (grün: `p21-n`, `p21-q`, `p21-t`, `p21-ni`) |
| `-i` gilt in jeder Schreibweise | `MAKEFLAGS` auf literales `*-i*` prüfen | alle sechs `p21-i-*` (auch `p21-i-rot`: `make` schreibt `i`, nicht `-i`; die Erwartung „`p21-i-rot` bleibt grün“ im Plan stimmte nicht) |
| `-i` gilt in jeder Schreibweise (Buchstaben im ersten Wort) | erstes Wort muss genau `i` sein | `p21-i-ki` (`ik`) (grün: `p21-i-rot`, `p21-i-gruen`, `p21-i-env`, `p21-i-lang`, `p21-i-j`, `p21-n`, `p21-q`, `p21-t`, `p21-ni`) |
| Unter `-i` entfällt der Hash | nur die Liste entfernen | alle sechs `p21-i-*` |
| Unter `-i` entfällt die Liste | nur den Hash entfernen | alle sechs `p21-i-*` |
| Unter `-i` entfällt ein älterer Nachweis | das Entfernen weglassen | alle sechs `p21-i-*` (mit dem vorherigen grünen Lauf) |
| `-n`, `-t`, `-n` mit `-i` schreiben nichts und entfernen nichts | Rezept in `enforce.mk` mit `+` | `p21-n`, `p21-t`, `p21-ni` (grün: `p21-q`, alle `p21-i-*`) |
| `-q` schreibt nichts und entfernt nichts | Rezept mit `+` und `record-gates` ohne Voraussetzungen (`Makefile`, Gates vor `record-gates` gezogen) | `p21-n`, `p21-q`, `p21-t`, `p21-ni` (grün: alle `p21-i-*`); mit `+` allein bleibt `p21-q` grün: `make -q` hält an der Voraussetzung, das Rezept läuft nicht |
| Der Hook urteilt über das Repo seines Arbeitsverzeichnisses, auch aus dem Unterverzeichnis (Punkt 22 a) | `git rev-parse --show-toplevel` durch `pwd` | `p22-unterverzeichnis` (grün: `p22-leere-umgebung`) |
| Der Hook liest keine Umgebungsvariable (Punkt 22 b) | `: "${HOME:?}"` vor der Eingabe | `p22-leere-umgebung` (grün: `p22-unterverzeichnis`) |

`env -i` gibt bei identischer Liste `approve`, bei Änderung `block`; die Git-Konfiguration der Umgebung ändert die Antwort hier nicht
(keine Rückgabe nach Auftrag 2). Die zehn Mutanten dieser Nacharbeit sind alle rot an den genannten Fällen, die genannten grünen Fälle bleiben grün.

Die 41 Mutanten der ersten Tabelle: alle rot an den genannten Fällen, kein grüner Mutant. Zwei Einzel-Mutanten waren äquivalent und sind deshalb in der Form der Tabelle gefahren: Entfällt nur
die Prüfung des Werts beim Lesen, fängt die des Pfads dieselbe Liste (Feldzahl); `F` mit leerem Wert fängt zusätzlich die Hex-Prüfung.
Sortierung und Einmaligkeit haben keinen Fall (akzeptiertes Negativ, §6); Punkt 8 von `SPEC-052` (nacheinander) hat keinen Fall, der Kopf der
Gegenprobe und des Hooks nennen die Bindung.

**Größe:** Code-Commit `b259381` gegenüber `9662e13`: 711 Einfügungen, 78 Löschungen in 11 Dateien (davon `stop-hook-gegenprobe.sh` 537, Hook 110 Zeilen
geändert, `dateiliste.sh` 56, Plan und Doku 9). Über der Marke von 700 Zeilen aus §4 (Schätzung des Architect: 440 bis 620); Rückführung `in-progress → next` ist Entscheidung des Planners.

**Rückgaben an den Architect (entschieden am 2026-10-13, Commit `c005555`; die Entscheidung steht in §6 und `SPEC-052` Punkte 13, 20 bis 22, umgesetzt im Code-Commit `e92db76`):**

1. `make -i gates` mit rotem Gate schreibt die Liste und endet mit 0 (Probe im Scratchpad: Attrappe mit `exit 1`, `make -s -i gates`). `SPEC-052` Punkt 6 nennt `-i` nur für die schnelle Kette; für den Nachweis des vollen Laufs ist die Randform nicht genannt.
2. Aufrufarten des Treibers (`AGENTS.md` §3.12): §6 nennt sie nicht als Aufrufarten; als Fälle von Punkt 20 (git scheitert) sind „kein Repo“ und „ohne git“ gefahren. Offen ohne Fall: `bash -e` (der Hook setzt eigene Optionen), leere Umgebung (`env -i`), Arbeitsverzeichnis in einem anderen gültigen Repo (der Hook liest dessen Baum), `CLAUDE_PROJECT_DIR` fehlt (der Hook liest sie nicht; der Aufruf in `.claude/settings.json` ruft ihn dann gar nicht).
3. Aufgabe des Auftrags nennt „Pfad mit `..`“ unter den beschädigten Listen; §6 weist unzulässige Pfade dem Folge-Slice zu (Segmente). Nicht umgesetzt; in diesem Slice sperrt jeder Unterschied, ein Pfad mit `..` im Stand ist nie identisch mit der Liste des Arbeitsbaums.

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

- `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (3×, *geplant* mit der Kennung von
  `slice-harness-gate-kette-stop-hook-freigabe`, seit `slice-harness-gate-kette-schnell`) — Gegenstand der
  Hook-Arbeit; die Messung (Freigabe eines Stops nach einem Doku-Commit) setzt die Positivliste voraus und
  liegt deshalb bei `slice-harness-gate-kette-stop-hook-freigabe`, das den vierten Beleg und den Ausgang
  *verkörpert* liefert. Dieser
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

**Nachtrag (Planner, 2026-10-13, Marke eingetreten, §4 *Urteil zur Marke*):** Zwei Einträge erhalten einen Beleg und
erreichen 3× ohne Regel: `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Schätzung 440 bis 620, gemessen 778,
das Mehr in der gemeinsamen Gegenprobe) und `BEO-REPO/groessenprobe-ohne-messbaren-stand` (Marke ohne Metrik, 630
netto, 704 eingefügt, 778 als Summe). `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` bleibt bei 2×
(Closure, DoD). Die Zuweisung (*verkörpert* oder *geplant*) steht beim Nutzer.

Kein Eintrag erreicht mit diesem Plan die Schwelle 3× neu; der dritte Beleg zu
`BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` entsteht bei der Closure.

**Zählung nach `AGENTS.md` §3.13:** Liefer-Punkte: 3 (Dateiliste; Hook mit Freigabe nur bei identischer
Liste; Gegenprobe mit den Gate-Zeilen). Messung, Gate-Läufe, Review, Closure-Notiz und Register zählen
nicht. Schichten: 1 — Harness-Werkzeug (`tools/harness`, `.claude/hooks`, `harness/mk`,
`harness/README.md`, `harness/conventions.md`); keine Spezifikation (sie ändert
`slice-harness-gate-kette-schnell`), keine Nutzer- oder Wartungs-Doku, kein Produkt-Code, kein Bootstrap.
Nachgezählt am 2026-10-13 bei der Umordnung (vor die Release-Slices): Die Sendungen von und an die Nachbarn sind
unverändert, die Zählung bleibt 3 Liefer-Punkte und 1 Schicht.
**Sendung des Gebers `slice-harness-gate-kette-schnell`** (die Namen aus §1, die Randformen der Hook-Seite
aus §6, soweit hier verteilt): nachgezählt mit der Sendung 3 Liefer-Punkte und 1 Schicht; die Namen sind
Teil der Liefer-Punkte 1 bis 3, kein vierter. **Sendung an `slice-harness-gate-kette-stop-hook-freigabe`**
(§1): dort nachgezählt, im selben Commit wie dieser Plan. Geschätzt 440 bis 620 Zeilen mit Tests (§1,
Verteilung der Schätzung des Architect, am 2026-10-13 korrigiert); in einer Review-Sitzung prüfbar. Vom Architect am 2026-10-13 nachgezählt: 3 Liefer-Punkte, 1 Schicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
