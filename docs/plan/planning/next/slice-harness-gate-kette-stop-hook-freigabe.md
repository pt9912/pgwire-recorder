# Slice slice-harness-gate-kette-stop-hook-freigabe: Der Stop-Hook gibt bei geänderten inerten Pfaden frei, wenn die schnelle Kette grün ist

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
`slice-v1-abschluss-tls-client`): direkt hinter `slice-harness-gate-kette-stop-hook-liste`
(WIP-Limit 1); Reihenfolge in §4 *Start* und im Drift-Log der Roadmap. Der Schnitt der Hook-Arbeit in
`slice-harness-gate-kette-stop-hook-liste` und diesen Slice (vor dem Start, aus dem Plan
`slice-harness-gate-kette-stop-hook`) ist **Entscheidung des Nutzers vom 2026-10-12**
(„Deine Empfehlungen“).

**Bezug:** kein Anspruch des Produkts — der Slice liefert ein Werkzeug des Harness.
Bindung an Entscheidungen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Wahl der
Kette am geänderten Pfad; Accepted 2026-10-12), das sich auf
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

**Ziel:** Der Stop-Hook `.claude/hooks/stop-require-gates.sh` wählt die Kette nach `SPEC-052`: Er gibt
frei, wenn jede seit dem Stand des letzten **vollen** grünen Laufs geänderte, hinzugekommene oder
gelöschte Datei auf der Positivliste liegt (zu Beginn `docs/plan/planning/` und `docs/reviews/`) und der
Stand der schnellen Kette (`gates-schnell`) grün ist; jeder andere Fall verlangt `make gates`
(fail-closed). Die Gegenprobe `stop-hook-gegenprobe` bekommt je Randform dieses Slice einen Fall, der nur
im inerten Fall freigibt. Der Hook ist das Restrisiko der Entscheidung
([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Konsequenzen*): Ohne seine Gegenprobe
schließt der Slice nicht. Der Slice liefert die Ersparnis der Entscheidung: Ein Stop nach einem
Doku-Commit verlangt den vollen Lauf nicht mehr.

**Voraussetzung (Geber):** `slice-harness-gate-kette-stop-hook-liste` liegt in `done/`: die Dateiliste des
vollen Laufs mit ihrem Leser (dieser Slice liest sie über dieselbe Funktion, keine zweite Logik), der Hook,
der nur bei identischer Liste freigibt und sonst `make gates` verlangt, die Prüfung nicht aufgelöster
Index-Einträge **vor** jeder Entscheidung nach Pfad (dieser Slice stellt die Positivliste nach ihr, nie
davor), und das Gate `stop-hook-gegenprobe` mit dem Gerüst der Temp-Repos (dieser Slice fügt Fälle hinzu,
er legt das Gate nicht an). `slice-harness-gate-kette-schnell` liegt davor in `done/`: `SPEC-052`, das Ziel
`gates-schnell`, der Stand der schnellen Kette `.harness/state/gates-schnell-passed.diffsha` (Hash des
Arbeitsbaums aus `tools/harness/working-tree-hash.sh`). **Name der Dateiliste:** wählt der Implementer von
`slice-harness-gate-kette-stop-hook-liste`; der Planner trägt ihn bei dessen Closure hier ein (§6).
`SPEC-052` Teil *Hook* gilt nach seinem Eingangssatz, soweit `stop-hook-gegenprobe` die Zusagen prüft;
dieser Slice erweitert das Gate und ändert die Spezifikation nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Dateiliste, ihre Form, ihr Schreiber und der Vergleich „identisch oder voll“, die Randformen
  fehlender oder beschädigter Stand, nicht lesbarer Vergleich, Abbruch, nicht aufgelöster Index-Eintrag
  und Schleifenschutz — `slice-harness-gate-kette-stop-hook-liste` (liegt vor diesem Slice in `done/`):
  Dieser Slice setzt sie voraus und ändert sie nicht; ändert er sie, hat er den Plan geändert.
- `SPEC-052`, die schnelle Kette, ihr Stand und deren Gegenprobe — `slice-harness-gate-kette-schnell`
  (liegt in `done/`): Dieser Slice ändert die Spezifikation nicht; gibt der Implementer eine Randform
  zurück, entscheidet der Architect in `SPEC-052`.
- Aufnahme weiterer Pfade in die Positivliste (`docs/plan/adr/`, `spec/`, `docs/user/`, `README.md`) —
  Bestand bleibt bewusst stehen: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Setzungen des
  Nutzers*; eine Aufnahme ist eine Änderung von `SPEC-052` mit dem Beleg, dass kein langsames Gate den Pfad
  liest. Es gibt dafür keinen Slice und keine Kennung.
- Änderung von `make gates` — [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Entscheidung*
  Punkt 1: Die volle Kette bleibt der Nachweis für alles, was nicht auf der Positivliste steht.
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

- [ ] **Liefer-Punkt 1:** `.claude/hooks/stop-require-gates.sh` gibt nach `SPEC-052` frei, wenn alle
      geänderten Pfade auf der Positivliste liegen (Grenze mit Schrägstrich, bytegenau und case-sensitiv;
      gelöscht, umbenannt mit beiden Pfaden, untracked wie getrackt; Pfad mit `.`/`..`/leerem Segment,
      absolut oder mit NUL gilt als beschädigte Liste) **und** der Stand der schnellen Kette grün ist
      (Hash des aktuellen Arbeitsbaums gleich der Datei; fehlend, ungleich oder nicht lesbar heißt kein
      Stand); bei identischer Liste gibt er ohne Stand der schnellen Kette frei; jeder andere Fall
      verlangt `make gates`. Der Kopf-Kommentar sagt die Positivliste zu, soweit die Gegenprobe sie prüft
      (`AGENTS.md` §3.7, §3.11).
- [ ] **Liefer-Punkt 2:** Ein Symlink ist inert, wenn Link-Pfad und das lexikalisch gegen das Verzeichnis
      des Links aufgelöste Ziel (Eintrag `L`) auf der Positivliste liegen; Ziel außerhalb des Repos,
      absolut oder nicht auflösbar ist nicht inert. Liegt irgendwo im Arbeitsbaum ein Link außerhalb der
      Positivliste, dessen Ziel auf sie zeigt, ist jede Änderung nicht inert (Spiegelregel).
- [ ] **Liefer-Punkt 3:** Die Gegenprobe `stop-hook-gegenprobe` macht je Randform dieses Slice (§6) einen
      Fall in einem Temp-Repo, der nur im inerten Fall freigibt; ihre Zeile in `harness/README.md`
      §Sensors nennt die Positivliste und die Bindung an
      [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), soweit die Zeile aus
      `slice-harness-gate-kette-stop-hook-liste` sie bisher nicht trägt; `harness/conventions.md` (Träger-Tabelle,
      Zeile `grundlagen-durchsetzungsschicht.md`) bleibt richtig oder wird nachgezogen.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). Der Reviewer prüft
      fail-closed ausdrücklich: Jede Abweichung von „inert und Stand grün“ verlangt `make gates`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Die Mutationen und ihre roten Fälle stehen als Tabelle in §7 (Zusage · Mutation · roter
      Fall), vom Implementer selbst gefahren (`AGENTS.md` §3.10).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: eine weitere Datei
      `evidence/slice-harness-gate-kette-stop-hook-freigabe.md` in
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
| `.claude/hooks/stop-require-gates.sh` | update | Positivliste, Symlink-Regel, Stand der schnellen Kette; Rückfall auf voll bei jedem anderen Fall; Kopf-Kommentar |
| `tools/harness/` (Gegenprobe des Hooks aus `slice-harness-gate-kette-stop-hook-liste`) | update | je Randform dieses Slice ein Fall in einem Temp-Repo; das Fragment unter `harness/mk/` und die Anbindung an `GATE_CHECKS` bestehen |
| `harness/README.md` §Sensors | update, falls nötig | Zeile der Gegenprobe nennt die Positivliste |
| `harness/conventions.md` (Träger-Tabelle, Zeile `grundlagen-durchsetzungsschicht.md`) | update, falls nötig | nennt die Träger des geänderten Hooks |

- Der Hook liest die Liste über dieselbe Funktion wie in `slice-harness-gate-kette-stop-hook-liste`; die
  Positivliste ist die einzige neue Logik und steht nach der Prüfung nicht aufgelöster Index-Einträge.
- Kein Produkt-Code, kein Handbuch, keine Tests unter `internal/` oder `test/`.
## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-gate-kette-stop-hook-liste` liegt in `done/`
(WIP-Limit 1); dieser Slice steht direkt dahinter und vor `slice-v1-abschluss-tls-client` (Schritt 28) und den übrigen Release-Slices (Drift-Log der Roadmap vom
2026-10-13). Die Reihenfolge ist Entscheidung des Nutzers vom 2026-10-13. Die offenen Rückfragen an den Architect (Merge-Konflikt zu `slice-harness-gate-kette-stop-hook-liste`, Eingangssatz von `SPEC-052`, neue Mutationen, Rückführungsgrenzen) entscheidet er vor dem Start in §6. Der Name der
Dateiliste steht in §1 (Eintrag des Planners bei der Closure des Gebers; `grep` nach ihm in diesem Plan
ist nicht leer). Die Randformen aus §6 sind in `SPEC-052` festgelegt; dieser Slice ändert sie nicht.

**Größe — Verteilung der Schätzung des Architect (700 bis 950 Zeilen ungeteilt) durch den Planner, der
Architect bestätigt vor dem Start:** dieser Slice 260 bis 390 Zeilen mit Tests (Hook 60 bis 90,
Gegenprobe mit gut einem Dutzend Fällen 200 bis 300); in *einer* Review-Sitzung prüfbar.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff liegt über einer Review-Sitzung
  (über 500 Zeilen mit Tests; die Schätzung liegt bei 260 bis 390). Schnitt dann: die Symlink-Regel mit
  ihren Fällen (Liefer-Punkt 2) als eigener Slice hinter der Positivliste; beide mit der Gegenprobe ihrer
  Zusagen, denn ohne Gegenprobe schließt keiner.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer nimmt
  [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) zurück, oder das Format der Antwort des
  Stop-Hooks ändert sich so, dass die Gegenprobe nicht mehr gegen die echte Eingabe und Ausgabe des
  Hooks läuft.
## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (vom Auftraggeber gefahren), die Gegenprobe des Hooks rot unter
jeder Mutation der Tabelle in §6, die Messung als Evidence im Register, Closure-Notiz mit
Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12): Der Slice liefert einen neuen Vertrag (Wahl der Kette am Pfad).
Die Randformen sind vom Architect am 2026-10-12 **entschieden**, und zwar in der Tabelle von §6 des
Gebers `slice-harness-gate-kette-schnell` (Fundstelle `SPEC-052`; dessen erster Commit überträgt sie in
die Spezifikation, vor jedem Code). Dieser Slice trägt davon, ohne Doppelung mit
`slice-harness-gate-kette-stop-hook-liste`:

- Pfad auf der Positivliste; Pfad auf der Grenze (`docs/plan/planning-x/`, `docs/reviews-alt/`, eine Datei
  `docs/reviews`); Groß- und Kleinschreibung (`docs/Plan/planning/`);
- Pfad mit Segment `.`, `..` oder leerem Segment, absoluter Pfad, Pfad mit NUL: beschädigte Liste, nie
  inert;
- gelöschte oder umbenannte Datei (beide Pfade gelten als geändert und müssen inert sein);
- untracked Datei (zählt wie eine getrackte; eine durch `.gitignore` ausgeblendete ist nicht Teil der
  Liste, Bestand);
- Symlink (Link-Pfad und Ziel inert; Ziel lexikalisch gegen das Verzeichnis des Links aufgelöst; Ziel
  außerhalb des Repos, absolut oder nicht auflösbar nicht inert) und gespiegelter Link (Link außerhalb der
  Positivliste mit Ziel darin: jede Änderung nicht inert; der Bestand hat Links nur unter `.claude/rules/`
  mit Zielen außerhalb);
- leerer Diff gegen den vollen Lauf ohne Stand der schnellen Kette (frei, Teil des Gutfalls der
  Identität aus `slice-harness-gate-kette-stop-hook-liste`; hier der Fall, dass ein vorhandener roter
  Stand ihn nicht sperrt);
- Stand der schnellen Kette fehlt, ist rot oder gehört zu einem anderen Arbeitsbaum (der Hash des
  aktuellen Arbeitsbaums weicht von der Datei ab oder die Datei fehlt): voll;
- gleichzeitige Läufe (**Aufruf-Voraussetzung**: nacheinander, `AGENTS.md` §3.12): `SPEC-052` bindet jede
  Zusage daran, der Hook sagt gleichzeitige Läufe nicht zu.

Eine nicht genannte oder offene Randform entscheidet der Implementer nicht, er gibt sie zurück.

**Mutationen** (`AGENTS.md` §3.10): Der Implementer fährt je Zusage die Mutation selbst und trägt die
Tabelle in §7 ein. Die Zeilen mit „Planner“ in der letzten Spalte sind vom Planner abgeleitet, der
Architect bestätigt sie vor dem Start:

| Zusage | Mutation | roter Fall | Herkunft |
|---|---|---|---|
| Nur ein Pfad der Positivliste gibt frei | `spec/` in die Positivliste aufnehmen | Änderung in `spec/` gibt frei → rot | Architect |
| Ein Pfad auf der Grenze zählt nicht | Präfix ohne Verzeichnisgrenze vergleichen | `docs/plan/planning-x/` gibt frei → rot | Architect |
| Gelöschte und umbenannte Dateien zählen mit beiden Pfaden | nur den neuen Pfad prüfen | Umbenennung aus einem nicht inerten Pfad in einen inerten gibt frei → rot | Architect |
| Ein Symlink zählt mit seinem Ziel | Link ohne Ziel vergleichen | Link in inertem Pfad mit Ziel außerhalb gibt frei → rot | Architect |
| Ohne grünen Stand der schnellen Kette gibt der Hook nicht frei | Prüfung des Stands der schnellen Kette entfernen | rote oder fehlende Kette gibt frei → rot | Architect |
| Ein Link außerhalb der Positivliste mit Ziel darin macht jede Änderung nicht inert | die Spiegelregel entfernen | Änderung in `docs/plan/planning/` gibt frei → rot | Architect |
| Ein Pfad mit `..`, absolut oder leerem Segment gilt als beschädigte Liste | Segmente nicht prüfen | `docs/plan/planning/../../internal/x.go` gibt frei → rot | Architect |
| Ein leerer Diff gibt ohne Stand der schnellen Kette frei | Stand der schnellen Kette auch dann verlangen | Hook blockiert bei identischem Baum → rot | Architect |
| Der Stand der schnellen Kette gehört zum aktuellen Arbeitsbaum | den Hash nicht mit dem aktuellen vergleichen | Stand eines anderen Arbeitsbaums gibt frei → rot | Planner |
| Die Positivliste vergleicht case-sensitiv | Groß-/Kleinschreibung nicht unterscheiden | `docs/Plan/planning/` gibt frei → rot | Planner |
| Eine untracked Datei zählt wie eine getrackte | untracked Dateien aus der Entscheidung nehmen | untracked Datei in `internal/` neben geändertem inertem Pfad gibt frei → rot | Planner |

Die Mutation „nicht aufgelöster Index-Eintrag“ der Tabelle des Architect steht in
`slice-harness-gate-kette-stop-hook-liste` (dort als Konflikt bei sonst identischer Liste; Verteilung des
Planners, damit die Prüfung vor der Pfad-Entscheidung liegt und in der Zwischenstufe nicht fehlt).

**Risiken:**

- Der Hook ist das Restrisiko der Entscheidung: Ein Fehler dort bleibt ungeprüft, wo die Gegenprobe
  ihn nicht fängt ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *Konsequenzen*) — **Ausgang:** entfallen: Die Gegenprobe macht je
  Randform einen Fall, jede Mutation der Tabelle ist rot gesehen, und ohne sie schließt der Slice
  nicht (DoD-Punkt 3).
- Die Gegenprobe fährt den Hook in Temp-Repos mit dessen Eingabe, nicht in der Umgebung von Claude
  Code; ob die Umgebung die Antwort (`decision`) so liest, prüft sie nicht — **Ausgang:** weiter offen:
  → `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (von `slice-harness-gate-kette-stop-hook-liste`
  angelegt, 1×; mit diesem Slice 2×, Evidence bei der Closure).
- Ein Pfad der Positivliste wird später von einem langsamen Gate gelesen ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), *akzeptiertes
  Negativ*) — **Ausgang:** entfallen: Ein späterer Leser ändert Code oder Gates und löst damit den
  vollen Lauf aus; die ADR hat das Negativ akzeptiert.
- Zwei Läufe gleichzeitig (der Hook liest, während die Kette ihren Stand schreibt) — **Ausgang:**
  entfallen: `SPEC-052` bindet die Zusagen an Läufe nacheinander (Aufruf-Voraussetzung, entschieden
  am 2026-10-12, `AGENTS.md` §3.12), und der Hook sagt gleichzeitige Läufe nicht zu; ein halb
  geschriebener Stand endet bei voll („kein Stand“).
- Die Positivliste liest sich als Zusage, die die Gegenprobe nur an Pfaden auf und neben der Grenze
  prüft (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, 2×) — **Ausgang:** entfallen: Der
  Kopf-Kommentar sagt nur zu, was ein Fall der Gegenprobe prüft (`AGENTS.md` §3.11), und jede Mutation
  der Tabelle ist rot gesehen; tritt es dennoch auf, legt das Review den dritten Beleg an.
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
  Hook und Kette, Herkunfts-Anker `seit slice-harness-gate-kette-stop-hook-freigabe`, sofern der Hook die
  Laufzeit eines Doku-Commits zeigbar senkt, sonst *offen* mit der Begründung.
- `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (1×, offen) — berührt: Symlink und nicht lesbarer
  Eintrag sind Randformen (§6); übergeht der Hook einen Eintrag still, ist es das zweite Auftreten.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — berührt: Die Positivliste liest
  sich als Zusage; die Gegenprobe prüft Pfade auf und neben der Grenze. Tritt es wieder auf, steht der
  Eintrag bei 3×.
- `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (2×, offen) — den dritten Beleg legt
  `slice-harness-gate-kette-stop-hook-liste` an; hier berührt: §4 Größe und Rückführung.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×, offen) — der Slice ist wellenlos; die
  Randformen entscheidet der Architect vor dem Code des Gebers, §4 *Start* nennt es.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) und
  `BEO-REPO/slice-waechst-durch-uebernahmen` (verkörpert, `AGENTS.md` §3.13) — Mutationen in §6, kein
  Kommentar sagt mehr zu, als die Gegenprobe prüft, die Sendungen sind nachgezählt.

Dieser Slice führt den dritten Beleg für `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` und ist mit
`slice-harness-gate-kette-schnell` selbst die Antwort darauf; sonst erreicht keiner der Einträge mit
diesem Plan die Schwelle 3× neu.

**Zählung nach `AGENTS.md` §3.13:** Liefer-Punkte: 3 (Positivliste mit Freigabe und Stand der schnellen
Kette; Symlink-Regel; Fälle der Gegenprobe mit der Gate-Zeile). Messung als Beleg, Gate-Läufe, Review und
Closure-Notiz zählen nicht. Schichten: 1 — Harness-Werkzeug (`.claude/hooks`, `tools/harness`,
`harness/README.md`, `harness/conventions.md`); keine Spezifikation (sie ändert
`slice-harness-gate-kette-schnell`), keine Nutzer- oder Wartungs-Doku, kein Produkt-Code, kein Bootstrap.
**Sendungen** des Gebers `slice-harness-gate-kette-stop-hook-liste` (Liste mit Leser, Gerüst von
`stop-hook-gegenprobe`, Prüfung nicht aufgelöster Index-Einträge vor der Pfad-Entscheidung, Name der
Dateiliste) und der Randformen und Mutationen der Hook-Seite aus `slice-harness-gate-kette-schnell`, soweit
verteilt: nachgezählt mit den Sendungen 3 Liefer-Punkte und 1 Schicht; der Name der Dateiliste ist Teil
der Liefer-Punkte 1 bis 3, kein vierter. Geschätzt 260 bis 390 Zeilen mit Tests (§4). Nachgezählt am
2026-10-13 bei der Umordnung (vor die Release-Slices): Die Sendungen sind unverändert, die Zählung bleibt 3
Liefer-Punkte und 1 Schicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
