# Slice slice-harness-nachweis-uebernehmen-aufrufarten: Die Gegenprobe des Nachweis-Werkzeugs trägt die Signale, die übrigen Make-Flags, den Aufruf aus einem Unterverzeichnis und die Kopplung an den Hook

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten, sechste Stelle der Reihenfolge, direkt hinter `slice-harness-nachweis-uebernehmen` und vor `slice-harness-gegenprobe-baumkopie`; wie jener Slice aus der Reihenfolge der Wellen gelöst (Entscheidung des Nutzers vom 2026-10-18, so vom Koordinator übermittelt): Start-Trigger ist `slice-harness-nachweis-uebernehmen` in `done/`, nicht `welle-harness-regeln`. Angelegt am 2026-10-18 beim Schnitt von `slice-harness-nachweis-uebernehmen` (Regel des Nutzers: schneiden, wenn die Schätzung die Marke selbst überschreitet). Der Slice liegt in `open/`: Er ist nicht priorisiert, den Anspruch trägt der Nutzer.

**Bezug:** [LH-QA-07](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur insoweit, als das Werkzeug und seine Gegenprobe dessen Messmethode (4) tragen; kein Anspruch des Produkts, der Slice ändert die Gate-Kette des Harness. Aktive ADR: [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) (Accepted, Entscheidung des Nutzers vom 2026-10-18), [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md). Bindung an das Register: `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (2×; dieser Slice fährt den Hook in der Gegenprobe, §8).

**Berührte Spec-Stellen:** [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Teil 4 und 5 und die Punkte 15 (die vier Fälle mit Signal), 16 (a), (c), (d), 17 und 18 (der Vertrag liegt vor); [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) Punkte 16 und 19 (der Hook, nur gelesen, seine Antworten sind der Gegenstand der Fälle der Kopplung).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Gegenprobe von `make nachweis-uebernehmen` (aus `slice-harness-nachweis-uebernehmen`) trägt die Fälle des Teils 2 von [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md): `SIGINT` und `SIGTERM` zwischen Prüfung und Ersetzen (Exit 130 und 143, keine Hilfsdatei, kein halbes Paar; Punkt 15, vier Fälle), `-n`, `-q`, `-t` (auch aus der Umgebung), ihre Kombination mit `-i`, `-k`, `-j` und das Langflag (Punkt 16 (a), (c), (d), zwölf Fälle), den Aufruf aus einem Unterverzeichnis (Punkt 17, zwei Fälle) und die Kopplung an den Hook (Punkt 18, vier Fälle: der unveränderte Hook antwortet nach Erfolg `approve`, vorher und nach einer Ablehnung `block`); das Skript bekommt dafür die Behandlung der Signale und, wo ein Fall es verlangt, den Wechsel in die Wurzel des Repos. Mit diesem Slice sind alle 93 Fälle von [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) geliefert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Werkzeug selbst, Quelle, Gleichheit, Ziel und `-i` (Punkte 1 bis 14, drei Fälle von Punkt 15, Punkt 16 (b)) — Klasse 1, rückwärts: `slice-harness-nachweis-uebernehmen` liefert sie und liegt in `done/`, wenn dieser Slice startet; sie ändern sich hier nur, wenn ein Fall des Teils 2 einen Fehler zeigt (Rückgabe an den Architect, wenn der Vertrag betroffen ist).
- Eine Änderung am Hook, an `dateiliste_erzeugen`, `dateiliste_pruefen` oder der Positivliste — Bestand bleibt bewusst stehen (wie in `slice-harness-nachweis-uebernehmen`): Der Hook ist der unveränderte Leser, seine Antworten sind der Gegenstand der vier Fälle.
- Die Umgebung von Claude Code, in der der Hook läuft — Schicht-Abgrenzung und akzeptiertes Negativ: Die Fälle lesen die Antwort des Hooks in einer Kopie; ob die Umgebung sie so liest, prüft keine Gegenprobe (`BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung`, der Kopf der Gegenprobe nennt es).
- Rollen-Anweisungen (wann das Werkzeug zu rufen ist) und `AGENTS.md` — Frage 4 in §6 von `slice-harness-nachweis-uebernehmen`: eine Pflicht, die eine Entscheidung des Nutzers nach der Lieferung braucht, und eine weitere Schicht (`AGENTS.md` §3.13).
- Die Übernahme des Standes der schnellen Kette und jede Lockerung der zweiten Sperre des Hooks — Entscheidung des Nutzers vom 2026-10-18 (Option (a), nichts ändern), so vom Koordinator übermittelt.

**Sendungen und Bindungen aus `slice-harness-nachweis-uebernehmen`** (dort §8 eingetragen, hier mit der Kennung des Gebers angenommen; `AGENTS.md` §3.13): (1) die Fälle und Mutationen der Punkte 15 (Signale), 16 (a), (c), (d), 17 und 18 von [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md); (2) das Rezept im Fragment `harness/mk/nachweis-uebernehmen.mk` von `slice-harness-nachweis-uebernehmen` trägt weder `$(MAKE)` noch `+` — Punkt 16 (a) prüft es, und die Mutation „das Rezept trägt `+`“ ändert das Fragment des Gebers in der Kopie; (3) das Skript arbeitet aus der Wurzel des Repos (Punkt 17); (4) Kennwörter und Exit-Codes sind die von Punkt 1, die Fälle dieses Slice lesen sie. Der Geber sagt zu diesen Punkten nichts zu (sein §1), deshalb beginnt die Zusage hier.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Signale):** Das Skript `tools/harness/nachweis-uebernehmen.sh` beendet sich bei `SIGINT` mit 130 und bei `SIGTERM` mit 143, ohne Hilfsdatei und ohne Hilfsverzeichnis, jede der beiden Dateien ist danach die alte oder die neue ganz (Punkt 15); die Gegenprobe trägt `p15-sigint-pruefen`, `p15-sigterm-pruefen` (Attrappe von `cmp` wartet), `p15-sigint-ersetzen`, `p15-sigterm-ersetzen` (Attrappe von `mv` wartet vor der ersten Ersetzung) mit den Mutationen aus Punkt 15, und jeder Fall prüft, dass die Attrappe angelaufen war (`AGENTS.md` §3.10).
- [ ] **Liefer-Punkt 2 (Aufrufarten):** Punkt 16 (a), (c), (d) und Punkt 17 von [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md): 14 Fälle (`p16-n`, `p16-q`, `p16-t`, `p16-n-env`, `p16-q-env`, `p16-t-env`, `p16-ni`, `p16-qi`, `p16-ti`, `p16-k`, `p16-j`, `p16-langflag`, `p17-relativ`, `p17-absolut`), am Skript und über `make`, je mit der Mutation, die nur die im Vertrag genannten Fälle färbt (das Rezept trägt `+`; das Skript überliest `i`, beides; das Skript liest das ganze `MAKEFLAGS`; ohne `k` und `j` lehnt es ab; eine relative Quelle gilt gegen die Wurzel; das Skript wechselt nicht in die Wurzel).
- [ ] **Liefer-Punkt 3 (Kopplung an den Hook):** Punkt 18 von [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md): `p18-erfolg`, `p18-ohne-nachweis`, `p18-untracked`, `p18-alt`, in einer Kopie des Temp-Repos mit dem Hook aus dem Arbeitsbaum; die Mutationen aus Punkt 18 (das Skript vergleicht nicht; das Skript schreibt bei jeder Ablehnung die frische Liste; das Skript schreibt keine Liste); `harness/README.md` führt die Zeile der Gegenprobe mit dem Stand der Fälle (Zahl der Fälle und Punkte nach Prüfung des Architect).
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `tools/` und `harness/` ändern); `make nachweis-uebernehmen-gegenprobe` einzeln grün mit allen 93 Fällen, die Dauer im Bericht; die Messung der Größe nach der Metrik von §4 im Bericht.
- [ ] **Probe am Endstand (Verifier):** wie in `slice-harness-nachweis-uebernehmen`, dazu die Fälle der Kopplung als Lauf am Endstand; der Verifier bricht je Punkt eine Zusage und sieht den Fall rot.
- [ ] Beobachtungs-Register: Beleg in `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (dritter Beleg, falls die Grenze erneut als offen stehen bleibt), sonst Antwort „keine Beobachtung“ in §7; bei drei Belegen legt der Planner den Folge-Slice an (§8).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review. Das Review fährt die Gegenprobe gezielt und prüft die Aufzählungen aus §6 je Element.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/nachweis-uebernehmen.sh` | update | Behandlung von `SIGINT` und `SIGTERM` (Hilfsverzeichnis und Hilfsdateien aufräumen, Exit 130 und 143), Wechsel in die Wurzel, falls `slice-harness-nachweis-uebernehmen` ihn nicht schon trägt |
| `tools/dev/harness/nachweis-uebernehmen-gegenprobe.sh` | update | 22 Fälle des Teils 2 mit Mutationen; Attrappen für `cmp` und `mv` (warten), zweite Kopie des Temp-Repos mit dem Hook aus dem Arbeitsbaum |
| `harness/README.md` | update | Zeile der Gegenprobe nennt Fälle und Punkte des Teils 2 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` setzt dieser Übergang; der Nutzer priorisiert und trägt den Anspruch. **Vorbedingungen:** (a) `slice-harness-nachweis-uebernehmen` liegt in `done/` (Werkzeug, Fragment, Gegenprobe mit 71 Fällen und der Vertrag `SPEC-062` stehen); (b) Vertrag [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) und [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) liegen vor (`AGENTS.md` §3.12); (c) **WIP:** eigener Zweig (Worktree); `Verantwortlich:` trägt beim Start Person und Zweig.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness`. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt.

**Schätzung** (Architect für Teil 2 der ungeteilten Schätzung: 27 Fälle, 200 bis 300 Zeilen; Planner nach dem Schnitt: die fünf `-i`-Fälle sind bei `slice-harness-nachweis-uebernehmen`): 22 Fälle, Skript (Signale, Wurzel) 15 bis 30, Gegenprobe 140 bis 200, `harness/README.md` 10 bis 20, zusammen **165 bis 250**; unter der Meldeschwelle 700 und der Marke 900.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900 (nach der Schätzung nicht zu erwarten) oder die Fälle sind nicht in einer Review-Sitzung prüfbar.
- `in-progress` → `open` (blockiert): Der unveränderte Hook antwortet in einem Fall der Kopplung anders, als [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 18 zusagt (Annahme des Vertrags falsch, Befund an den Planner und den Architect), oder ein Fall des Teils 2 zeigt einen Fehler im Werkzeug von `slice-harness-nachweis-uebernehmen`, den nur eine Änderung des Vertrags behebt (Rückgabe an den Architect, `.claude/commands/implement-slice.md`, Randform-Rückgabe).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; `make gates` grün auf dem Endstand; `make nachweis-uebernehmen-gegenprobe` grün mit allen 93 Fällen; die Probe am Endstand (Verifier); Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag geschrieben; die drei Paarungen getragen. Bei der Closure prüft der Planner den Zähler von `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (§8); `grep` nach `slice-harness-nachweis-uebernehmen` in §1 und §6 dieses Slice findet die vier Sendungen.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12). Die Nummern sind die des Plans von `slice-harness-nachweis-uebernehmen` (§6 dort, Numerierung unverändert); entschieden ist jede in [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted). Der Implementer entscheidet keine still (Randform-Rückgabe).

14. **Signale.** `SIGINT` und `SIGTERM` beenden das Skript mit 130 und 143, ohne Hilfsdatei und ohne Hilfsverzeichnis; jede der beiden Dateien ist alt oder neu ganz, nie halb. Ob nach einem Signal zwischen den beiden Ersetzungen das Paar zusammengehört, ist nicht zugesagt (akzeptiertes Negativ, Grund: das Fenster sind zwei `mv`; Hash neu und Liste alt lehnt Punkt 7 ab). Vier Fälle: je Signal und Phase (prüfen, ersetzen). Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 15.
17. **Make-Flags außer `-i`.** `-n`, `-q`, `-t` (Kommandozeile und `MAKEFLAGS` der Umgebung): `make` führt das Rezept nicht aus, nichts ändert sich, auch dann nicht, wenn die Übernahme gelänge; der Exit-Code unter `-q` ist nicht zugesagt. `-i` zusammen mit `-n`, `-q` oder `-t`: `-n`, `-q`, `-t` gelten, keine Meldung des Skripts. `-k`, `-j` und Langflags: ohne Wirkung (die Auswahl liest nur das erste Wort von `MAKEFLAGS`). Zwölf Fälle: `p16-n`, `p16-q`, `p16-t`, die drei `-env`, `p16-ni`, `p16-qi`, `p16-ti`, `p16-k`, `p16-j`, `p16-langflag`. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 16 (a), (c), (d).
18. **Aufruf aus einem Unterverzeichnis.** Das Skript urteilt über den Baum der Wurzel (`git rev-parse --show-toplevel`); eine relative Quelle gilt gegen das Arbeitsverzeichnis des Aufrufs; ein Aufruf über `make -C` ist Sache von `make` (akzeptiertes Negativ). Zwei Fälle: `p17-relativ`, `p17-absolut`. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 17.
19. **Aufruf-Voraussetzung** (nacheinander, ohne Fall): Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 19; der Fall der Kopplung nimmt sie als gegeben.
20. **Der unveränderte Hook liest das Ergebnis.** Nach Erfolg `approve` ohne Stand der schnellen Kette, vorher bei fehlendem Nachweis `block`, nach Ablehnung wegen einer untracked Datei `block`, nach Ablehnung mit älterem gültigem Nachweis eines anderen Baums dessen Befund (`block`). Vier Fälle; der Fall `p18-untracked` legt die untracked Datei außerhalb der Positivliste an. Der Hook ist ein zweites Netz: Bildete das Skript die Liste ohne `--others`, träfe der Hook die untracked Datei weiter und sperrte; diese Mutation färbt `p9-untracked` (bei `slice-harness-nachweis-uebernehmen`), nicht `p18-untracked`. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 18.

### Risiken

- **Der Hook wird zum dritten Mal berührt, ohne dass die Gegenprobe seine Umgebung prüft** (`BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung`, 2×; die Fälle der Kopplung fahren den Hook und lesen seine Antwort, ob die Umgebung von Claude Code sie so liest, prüft keiner; der Kopf der Gegenprobe nennt es). — **Ausgang (vorgesehen):** weiter offen → dieselbe Beobachtung; legt der Planner bei der Closure einen dritten Beleg an, braucht sie nach Modul 5 einen eigenen Folge-Slice (Frage 6 im Plan von `slice-harness-nachweis-uebernehmen`).
- **Die Mutation „das Rezept trägt `+`“ ändert das Fragment von `slice-harness-nachweis-uebernehmen`** und färbt Fälle, die der Geber nicht kennt; ändert `slice-harness-nachweis-uebernehmen` bis zum Start Namen oder Aufbau des Fragments, trifft die Mutation ins Leere. — **Ausgang (vorgesehen):** entfallen, wenn `grep` im Fragment von `slice-harness-nachweis-uebernehmen` beim Start das Rezept ohne `$(MAKE)` und `+` zeigt und die Mutation in der Gegenprobe rot färbt (Fall `p16-n` rot).
- **Der Geber `slice-harness-nachweis-uebernehmen` zieht die Sendungen nicht nach**, wenn er Vertrag oder Namen ändert (`AGENTS.md` §3.13). — **Ausgang (vorgesehen):** entfallen mit `grep` nach `slice-harness-nachweis-uebernehmen` in §1 und §6 (Closure-Trigger).
- **Ein Fall der Signale ist zeitabhängig** (die Attrappe wartet, das Signal geht an das Skript). — **Ausgang (vorgesehen):** entfallen, wenn jeder Fall prüft, dass die Attrappe angelaufen war, und die Gegenprobe dreimal in Folge grün ist (Bericht).

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

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** Beleg oder Antwort „keine Beobachtung“ zu `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (bei der Closure)
- **Folge-Slices:** — (ein dritter Beleg im Register würde einen anlegen)
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen in jedem Slice-Plan.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (**2×**; die vier Fälle der Kopplung fahren den Hook in einer Kopie, die Umgebung von Claude Code prüft keiner; erhält die Beobachtung bei der Closure einen dritten Beleg, ist sie keine Notiz mehr, sondern eine Lücke und braucht einen eigenen Folge-Slice), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (DoD trägt je Fall eine Mutation und den Nachweis, dass die Attrappe anlief), `BEO-REPO/spec-randform-erst-im-review-entschieden` (§6 nennt die Randformen, entschieden in [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md)), `BEO-REPO/slice-waechst-durch-uebernahmen` (dieser Slice übernimmt vier Sendungen von `slice-harness-nachweis-uebernehmen`; Zählung unten), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groessenprobe-ohne-messbaren-stand` (§4 rechnet in der Metrik). Keiner der übrigen Treffer erreicht mit diesem Slice 3× neu.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-18, im Commit des Schnitts, der die Sendungen von `slice-harness-nachweis-uebernehmen` einträgt):** Liefer-Punkte 3 (Signale | Aufrufarten | Kopplung an den Hook); die Sendungen sind diese drei Punkte, ein vierter hat keinen Platz. Schichten 1: Harness-Werkzeug mit Gegenprobe und Gate-Index; Planung, Vertrag und ADR sind keine Schicht, Dokumentation unter `docs/` wird hier nicht geändert. Schätzung in der Metrik von §4: 165 bis 250 (Marke 900). Der Geber `slice-harness-nachweis-uebernehmen` zählt in seinem §8 mit: Liefer-Punkte 3, Schichten 2, 665 bis 940.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der die Ortsangabe aus `slice-harness-gegenproben-umzug-dev` einträgt): Die Gegenproben liegen seit jenem Slice unter `tools/dev/<bereich>/`, die Pfade der Gegenproben in diesem Plan sind nachgezogen (`tools/dev/` ist wie `tools/harness/` Harness-Werkzeug, keine neue Schicht). Ortsangabe der Gegenprobe geändert, Liefer-Punkte (3) und Schichten (1) unverändert.
