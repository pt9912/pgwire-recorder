# Slice slice-harness-gegenproben-aus-gates-sammelziele: Zwei Sammelziele fahren die Gegenproben, `make gegenproben` im Gegenproben-Container, `make gegenproben-host` am Host

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten, zweite Stelle der Reihenfolge, direkt hinter `slice-harness-gegenproben-aus-gates` und vor `slice-harness-gegenproben-aus-gates-fahrstelle` (Entscheidung des Nutzers vom 2026-10-18, sofort). Wie jener Slice aus der Reihenfolge der Wellen gelöst: Start-Trigger ist `slice-harness-gegenproben-aus-gates` in `done/`, nicht `welle-harness-regeln` in `done/`; der Start-Trigger der Welle gilt für die übrigen Slices (Drift-Log der Roadmap, 2026-10-18). Angelegt am 2026-10-18 durch den Schnitt von `slice-harness-gegenproben-aus-gates`: Der Architect schätzte den ungeteilten Slice auf 610 bis 1100 Zeilen (Metrik von `AGENTS.md` §3.14, Marke 900) und empfahl den Schnitt an der Naht des Vertrags (`SPEC-060` Teil 1 gegen Teil 2); nach der Regel des Nutzers (schneiden, wenn die Schätzung die Marke selbst überschreitet) schneidet der Planner so. Dieser Slice behält Teil 2: die Sammelziele mit ihrer Gegenprobe.

**Bezug:** [LH-QA-07](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur insoweit, als die Gegenproben dessen Messungen tragen; kein Anspruch des Produkts, der Slice ändert die Gate-Kette des Harness. Bindung an das Register: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Beleg des Kettenumbaus: `evidence/slice-harness-gegenproben-aus-gates.md`) und das mit `slice-harness-gegenproben-aus-gates` entstehende `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`. Aktive ADRs: [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (Accepted am 2026-10-18, Entscheidung des Nutzers; schärft `SPEC-060`), [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) (die Host-Liste, die ein Sammelziel fährt).

**Berührte Spec-Stellen:** `SPEC-060` Teil 2 (Punkte 4 bis 12, die Sammelziele) und die Punkte 13 bis 15, soweit sie die Sammelziele betreffen (Vertrag in `harness/vertraege/`, vom Architect am 2026-10-18 geschrieben); [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) Punkt 8 (Aufruf-Voraussetzung nacheinander), nur gelesen.

**Verantwortlich:** pt9912 (Zweig slice-harness-gegenproben-aus-gates-sammelziele)

**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make gegenproben` (Namen laut `SPEC-060` Punkt 4) fährt die Gegenproben der Container-Liste (`GEGENPROBEN_CONTAINER`, alle ohne `u-boot` und `devcontainer`), `make gegenproben-host` die der Host-Liste (`GEGENPROBEN_HOST`, heute `container-start-gegenprobe`), beide über ein Skript (`tools/harness/gegenproben.sh`, Fragment `harness/mk/gegenproben.mk`): nacheinander in der Reihenfolge der Liste, alle auch nach einem roten Einzelziel, Exit 1 bei mindestens einem roten, eine Zusammenfassung mit den Namen. Beide Sammelziele sind kein Gate und kein Teil von `make gates`; kein Hook und kein Stand glaubt ihr Ergebnis. Eine Gegenprobe der Sammelziele (`make gegenproben-gegenprobe`, an der Container-Liste, mit Attrappen als Einzelziele, ohne Docker) trägt die Fälle und Mutationen aus `SPEC-060` Teil 2 (nach dem Architect 37 Fälle bei 16 Randformen der Gruppe G).

**Voraussetzung, kein Teil dieses Slice:** Die Listen, die Namensregel und der Umbau der Kette liefert `slice-harness-gegenproben-aus-gates` (`done/` vor dem Start); die Sammelziele lesen die Listen, die er anlegt.

**Pflicht zum Fahren, zweite Hälfte.** Der Nutzer hat am 2026-10-18 den Wortlaut für `AGENTS.md` §6 Schritt 6 entschieden: „Wer ein Gate, eine Gegenprobe oder einen Vertrag ändert, fährt die betroffene Gegenprobe (`make <name>`, Zuordnung in `harness/README.md` §Sensors); am Ende jeder Welle und vor dem Release fahren `make gegenproben` im Gegenproben-Container und `make gegenproben-host` am Host.“ Der Satz nennt Ziele, die erst dieser Slice liefert (`AGENTS.md` §3.11: eine Zusage nennt nur, was existiert). Der Planner teilt ihn deshalb (Vorschlag des Planners, nicht Teil der Entscheidung des Nutzers): `slice-harness-gegenproben-aus-gates` trägt die Hälfte bis zum Semikolon, dieser Slice den Rest, sodass nach seiner Closure der Wortlaut des Nutzers vollständig in Schritt 6 steht. Die Teilung ist am Stand `998738e8` geprüft: Schritt 6 endet heute mit „… (`make <name>`, Zuordnung in `harness/README.md` §Sensors).“ (erste Hälfte, vom Kettenumbau); dieser Slice ersetzt den Schlusspunkt durch das Semikolon und setzt an: „; am Ende jeder Welle und vor dem Release fahren `make gegenproben` im Gegenproben-Container und `make gegenproben-host` am Host.“ Beide Hälften zusammen ergeben den Wortlaut des Nutzers Zeichen für Zeichen; DoD-Punkt 3 prüft das mit `grep` auf den ganzen Satz. Die Rollen-Anweisungen und die Fahrstelle liefert `slice-harness-gegenproben-aus-gates-fahrstelle` (unten).

**Übernimmt:** nichts. Teil 2 des Vertrags ist durch den Schnitt Gegenstand dieses Slice, nicht durch eine Übernahme von einem Slice, der geschlossen wird.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Listen `GEGENPROBEN_CONTAINER` und `GEGENPROBEN_HOST`, die Namens- und Zuordnungsregel (`SPEC-060` Teil 1), der Anschluss der zwölf Fragmente, die Ausnahmeliste der schnellen Kette und die Gegenprobe der Kette — `slice-harness-gegenproben-aus-gates`, dessen `done/` dieser Slice voraussetzt. Begründung: Naht des Vertrags und Größe (Schätzung des ungeteilten Slice über der Marke); der Kettenumbau ist die Beschleunigung, die der Nutzer will, und steht vorn.
- Rollen-Anweisungen (`implement-slice`, `close-welle`, Reviewer- und Verifier-Skill, Implementer-Typ), der Abschnitt in `docs/developer/devcontainer.md` über den Gegenproben-Container, der Hook für die Einlieferung und die Sendung an `slice-erster-release-freigabe` — `slice-harness-gegenproben-aus-gates-fahrstelle`, direkt hinter diesem Slice. Begründung: Schicht-Abgrenzung (Rollen-Anweisungen und Wartungs-Doku sind weitere Schichten, `AGENTS.md` §3.13) und Größe.
- Ein Stand oder Nachweis der Gegenproben (analog `.harness/state/gates-passed.diffsha`) und ein Hook, der ihn liest — anderer Vorgang und Entscheidung des Nutzers (Frage 4 in `slice-harness-gegenproben-aus-gates`, Empfehlung des Planners: nein): Ein zweiter Nachweis brächte die Laufzeit zurück, die der Nutzer streichen will. Das Sammelskript schreibt nichts in den Arbeitsbaum (`SPEC-060` Punkt 11).
- „Die Gegenprobe läuft nur, wenn sich ihre Eingaben geändert haben“ — vom Nutzer verworfen (Abgrenzung in `slice-harness-gegenproben-aus-gates` §1, dort mit Grund und Zeiger auf das Register).
- Kleinere Kopien und `git ls-files`-Helfer der Gegenproben — `slice-harness-gegenprobe-baumkopie`; die Sammelziele ändern nicht, was eine Gegenprobe tut, nur, wie die Liste gefahren wird.
- Neue Fälle in den Gegenproben der Gates — andere Vorgänge; Bestand bleibt bewusst stehen.
- Produkt-Code, Handbuch, `README.md` — Schicht-Abgrenzung: kein Verhalten des gebauten Binary.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Sammelziele):** `make gegenproben` und `make gegenproben-host` (Namen laut `SPEC-060`) fahren je ihre Liste nacheinander, alle auch nach einem roten Einzelziel, enden mit Exit ungleich 0 bei mindestens einem roten, schreiben die Zusammenfassung mit den Namen, brechen bei `SIGINT` und `SIGTERM` ohne Zusammenfassung ab, lehnen eine leere Liste und `-i` ab und fahren unter `-n`, `-q`, `-t` nichts; jede Entscheidung je Randform aus §6 (Gruppe G) steht im Vertrag (`SPEC-060` Punkte 4 bis 12), vor dem ersten Code-Commit. Fragment `harness/mk/gegenproben.mk` und Skript `tools/harness/gegenproben.sh`; beide Ziele stehen nicht an `GATE_CHECKS` und nicht in der Ausnahmeliste.
- [ ] **Liefer-Punkt 2 (Gegenprobe der Sammelziele):** `make gegenproben-gegenprobe` (Fragment unter `harness/mk/`, an `GEGENPROBEN_CONTAINER` und ausdrücklich nicht an `GATE_CHECKS`, Zeile im Gate-Index) trägt je Randform und je Element ihrer Aufzählungen aus §6 einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die Gegenprobe verlangt je Mutation genau die Menge roter Fälle, die der Vertrag nennt (Namen der Fälle tragen die Nummer des Punktes). Attrappen als Einzelziele, kein Docker.
- [ ] **Liefer-Punkt 3 (Gate-Index und Satz):** `harness/README.md` §Sensors führt die Zeilen der zwei Sammelziele (als Werkzeuge, kein Gate) und der Gegenprobe der Sammelziele; `AGENTS.md` §6 Schritt 6 trägt den Rest des Satzes (Ende jeder Welle und Release, `make gegenproben` im Gegenproben-Container und `make gegenproben-host` am Host) in dem Wortlaut, den der Nutzer am 2026-10-18 entschieden hat.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `harness/` und `tools/` ändern) und `make gates-schnell-gegenprobe` grün (die Deckung zählt die neue Gegenprobe); `make gegenproben` grün auf dem Endstand, **gefahren im Gegenproben-Container, falls er vorliegt, sonst am Host** (der Bericht nennt, wo und die Dauer); `make gegenproben-host` grün am Host.
- [ ] Beobachtungs-Register: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` mit der Messung der Fahrten ergänzt (Beleg `evidence/slice-harness-gegenproben-aus-gates-sammelziele.md`, angelegt bei der Closure); kein Zähler wird gesetzt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review. Das Review fährt die Gegenprobe der Sammelziele gezielt.
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
| `harness/vertraege/SPEC-060-gegenproben.md` Teil 2 | Vertrag (Architect, liegt vor dem Code) | Sammelziele, Entscheidung je Randform aus Gruppe G |
| `harness/mk/gegenproben.mk`, `tools/harness/gegenproben.sh` | neu | die zwei Sammelziele, das Skript, das eine Liste nacheinander fährt, sammelt und nach `MAKEFLAGS` entscheidet |
| `tools/harness/gegenproben-gegenprobe.sh`, `harness/mk/gegenproben-gegenprobe.mk` | neu | Gegenprobe des Sammelskripts in Temp-Bäumen mit Attrappen als Ziele (kein Docker): je Randform ein Fall und eine Mutation; hängt an `GEGENPROBEN_CONTAINER` |
| `harness/README.md`, `AGENTS.md` | update | Gate-Index (Zeilen der Sammelziele und ihrer Gegenprobe), Rest des Satzes in Schritt 6 |

**Arbeitsort.** Der Implementer arbeitet im Container `pgr-dev2` (Port 2223), Klon `/workspaces/pgr-dev2/pgwire-recorder`; dort laufen `make gates`, `make gegenproben-gegenprobe` und der Container-Teil, `make gegenproben`. Das Host-Ziel `make gegenproben-host` läuft am Host und wird in einem Wegwerf-Worktree des eingelieferten Commits geprüft, nicht im Arbeitsbaum des Containers.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` setzt dieser Übergang; der Nutzer priorisiert und trägt den Anspruch. **Vorbedingungen (beim Anspruch am 2026-10-18 am Stand `998738e8` geprüft: (a) bis (d) erfüllt, (e) mit diesem Übergang):** (a) Vertrag (`SPEC-060` Teil 2) und ADR liegen vor, mit der Entscheidung je Randform aus §6 (`AGENTS.md` §3.12), den Namen der Sammelziele und Listen; erfüllt, Accepted am 2026-10-18 (Zweig `arch-gates`, Commits des Architect vor dem Move nach `in-progress/`). (b) **Start-Trigger:** `slice-harness-gegenproben-aus-gates` liegt in `done/` (die Listen und die Namensregel stehen; die Host-Liste enthält `container-start-gegenprobe`, sonst ist sie leer, §6, G7). (c) **Reihenfolge:** Der Slice ist aus der Reihenfolge der Wellen gelöst und wartet nicht auf `welle-harness-regeln`; in der Welle bleibt er an zweiter Stelle. (d) **Wortlaut der Pflicht in `AGENTS.md` §6:** als erteilt eingetragen, Entscheidung des Nutzers vom 2026-10-18 (Wortlaut in §1); der Kettenumbau (`done/`) trägt die erste Hälfte bis zum Semikolon, dieser Slice den Rest. (e) **WIP:** Der Slice läuft auf eigenem Zweig; `Verantwortlich:` trägt beim Start Person und Zweig.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness AGENTS.md`. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; Vertrag und ADR (Commits des Architect vor dem Start) zählen nicht.

**Schätzung (Architect, Prüfung vom 2026-10-18; Metrik wie oben; Gegenprobe-Code mit dem Doppelten des Produktcodes; Doku-Anteil Planner):**

| Teil | Zeilen |
|---|---|
| `gegenproben.mk` und `gegenproben.sh` (zwei Ziele, Skript mit `MAKEFLAGS`, Schleife, Sammlung, Signale) | 80 bis 140 |
| `gegenproben-gegenprobe` (37 Fälle mit Mutationen, 16 Randformen; Attrappen kosten Zeilen) | 250 bis 420 |
| `harness/README.md` (Zeilen der Sammelziele und der Gegenprobe) und `AGENTS.md` (Rest des Satzes) | 20 bis 40 |
| **Zusammen** | **350 bis 600** |

Meldeschwelle 700, **Marke 900**; die Obergrenze liegt unter der Meldeschwelle. Zählt der Architect später mehr als 16 Randformen in Gruppe G oder mehr als 37 Fälle, rechnet er neu.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner (Vorschlag: Skript mit Gegenprobe der Punkte 4 bis 9 getrennt von den Punkten 10 bis 12, den Aufrufarten).
- `in-progress` → `open` (blockiert): Eine Gegenprobe der Container-Liste wird im Gegenproben-Container nur grün, wenn ein anderes Gate vorher lief (versteckte Abhängigkeit zwischen Gates; Befund an den Planner, Rückgabe an `slice-harness-gegenproben-aus-gates`).

## 5. Closure-Trigger

DoD vollständig; `make gates` grün auf dem Endstand; `make gegenproben` grün im Gegenproben-Container oder, wenn er nicht vorliegt, am Host (Beleg im Bericht mit Dauer); `make gegenproben-host` grün am Host; Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag geschrieben; die drei Paarungen getragen.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des neuen Vertrags (`AGENTS.md` §3.12).** Je Randform steht die Empfehlung des Planners; entschieden ist keine, der Architect legt sie vor dem Code in den Vertrag (`harness/vertraege/`, `SPEC-060` Teil 2), und der Implementer entscheidet keine still (Randform-Rückgabe). Was aus dem Bestand entscheidbar war, steht als „Bestand“. Die Gruppe K (Kette) gehört `slice-harness-gegenproben-aus-gates`.

**Gruppe G — Sammelziele**

1. **Zwei Ziele, zwei Listen, ein Skript.** Name der Ziele und der Listen: Vorschlag des Architect. Beide Ziele sind kein Gate, hängen nicht an `GATE_CHECKS`, hängen nicht aneinander; das Container-Ziel fährt die Host-Liste nicht. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 1 und 4: `gegenproben` und `gegenproben-host`, Listen `GEGENPROBEN_CONTAINER` und `GEGENPROBEN_HOST`, Skript `tools/harness/gegenproben.sh`, Fragment `harness/mk/gegenproben.mk`.
2. **Reihenfolge.** Nacheinander in der Reihenfolge der Liste (Reihenfolge der Includes, alphabetisch nach Fragmentdatei), nie parallel, auch nicht unter `-j` (eine Schleife im Skript statt eines Fächers von Prerequisites): Die Last (Docker, Load 57 auf 20 Kernen, gemessen in [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md)) ist der Grund für das Herausnehmen. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 5. Abweichung: Zugesagt ist die Reihenfolge, in der `make` die Liste ausgibt, nicht „alphabetisch nach Fragmentdatei“ (Grund: die Reihenfolge der Dateien beim `include` mit Muster ist Sache von `make`).
3. **Rotes Einzelziel.** Alle Ziele der Liste fahren, auch nach einem roten (wie `-k`); die Zusammenfassung nennt die roten beim Namen, Exit 1. Ein Abbruch durch `SIGINT` oder `SIGTERM` beendet das laufende Ziel und startet kein weiteres (Exit 130 und 143). Je ein Fall. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 6 und 8; ein Abbruch schreibt keine Zusammenfassung.
4. **Aufrufarten.** `-n`, `-q`, `-t`: Das Rezept trägt weder `$(MAKE)` noch `+`, make führt es nicht aus (wie `record-gates`); `-i`: Die Einzelziele erben `-i` und gäben ein rotes Skript als Exit 0 zurück, deshalb lehnt das Skript `-i` in `MAKEFLAGS` ab (Exit 2, Meldung), statt grün zu melden; `-k`, `-j` ohne Wirkung. `-i` zusammen mit `-n`, `-q` oder `-t`: `-n`, `-q`, `-t` gelten, nichts läuft. Kein Stand wird geschrieben. Je Flag ein Fall, auch aus der Umgebung und kombiniert. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 10, Reihenfolge der Gründe Punkt 12. Abweichung: Das Rezept übergibt `$(MAKE)` an das Skript wie `gates-schnell`, `make` führt es also auch unter `-n`, `-q`, `-t` aus, und das Skript entscheidet an `MAKEFLAGS` (Grund: Das Skript braucht das `make` des Aufrufers, und nur so hat es es; zugesagt ist das Ergebnis: kein Einzelziel läuft).
5. **Fehlendes Werkzeug.** Fehlt `docker`, `bash` oder `setsid`, ist das betroffene Einzelziel rot mit seiner eigenen Meldung, das Sammelziel fährt weiter; kein Vorab-Test der Werkzeuge. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 6 (Fall `p6-127`).
6. **Ausgabe.** Die Ausgabe der Einzelziele erscheint unverändert; die Zusammenfassung steht zuletzt und nennt je Ziel Name und Ausgang. Ein Überspringen innerhalb eines Ziels (das Überspringen der ersetzten ADR, ersetzt durch [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md)) gibt es nach der Aufteilung nicht mehr. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 6 und 7; die Zeilen der Zusammenfassung sind in Punkt 6 festgelegt.
7. **Leere Liste.** Die Host-Liste ist ohne `container-start-gegenprobe` leer (§4 *Start*); eine leere Liste ist ein Fehler (Exit 2, Meldung), kein stilles Grün. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 9.
8. **Eine Gegenprobe steht in beiden Zielen oder in keinem.** Das Sammelskript prüft die Zuordnung nicht, die Gegenprobe der Kette tut es (Gruppe K, K1 in `slice-harness-gegenproben-aus-gates`). Entschieden: [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 2.
9. **Aufruf.** Aus einem Unterverzeichnis des Repos wie von der Wurzel; ohne Git-Arbeitsbaum Exit 2. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): Abweichung von der Empfehlung: [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 13 sagt weder einen Aufruf aus einem Unterverzeichnis noch Exit 2 ohne Git-Arbeitsbaum zu (Grund: Die Schnittstelle ist das Make-Ziel, das Rezept läuft im Verzeichnis des `Makefile`, und das Skript ruft kein `git`).
10. **Gleichzeitige Läufe und Worktrees.** Zwei Fahrten im selben Arbeitsbaum oder zusammen mit `make gates` sind nicht zugesagt (`SPEC-052` Punkt 8, Aufruf-Voraussetzung nacheinander); zwei Worktrees am selben Docker-Daemon teilen Namen von Containern und Netzen (Präfixe `pgr-ivg-`, `pgr-probe-`), die Gegenprobe der Container schützt nur durch ihre Präfixe. Der Vertrag bindet die Zusage an „nacheinander“; die Fahrstelle im eigenen Container (`slice-harness-gegenproben-aus-gates-fahrstelle`) hat einen eigenen Daemon. **Aufruf-Voraussetzung (§3.12):** Das Sammelskript sagt keine Ordnung über mehrere Aufrufe zu. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 13.
11. **Dev-Container-Erkennung entfällt.** Keine Erkennung, kein Überspringen; die Zuordnung (K2 in `slice-harness-gegenproben-aus-gates`) ersetzt sie. Entschieden: Nutzer, 2026-10-18; [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) Accepted. Im Vertrag: [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 3.

**Aufrufarten des Treibers (§3.12):** Siehe G4; das Sammelskript liefert keinen Stand und keinen Nachweis, den ein Hook glaubt (die Aufrufarten für `make gates` stehen als K6 in `slice-harness-gegenproben-aus-gates`).

**Fahrstelle, Stand und Ressourcen (für `slice-harness-gegenproben-aus-gates-fahrstelle`, hier genannt, damit der Vertrag sie bindet).** Der Gegenproben-Container sieht nur den eingelieferten, committeten Stand eines Zweigs, nicht den Arbeitsbaum des Hosts; was nicht eingeliefert ist, ist nicht gefahren. Zwei Container desselben Klons sind nicht zugesagt. Ressourcen laut `docs/developer/devcontainer.md` §10: je Container 8 CPUs und 12 GiB (Vorgabe), Host 20 CPUs und 31 GiB; Arbeitscontainer und Gegenproben-Container zusammen 16 CPUs und 24 GiB. Der Vertrag sagt dazu nichts zu; die Zusage steht in der Dokumentation, soweit gemessen.

### Risiken

- **Eine Gegenprobe ist ohne Kette nur grün, wenn sie von selbst hält:** Sie lief im Verbund mit `make gates` und kann von einem Zustand leben, den ein anderes Gate herstellt (Image, Cache, Datei im Arbeitsbaum). — **Ausgang (vorgesehen):** entfallen, wenn `make gegenproben` im frischen Gegenproben-Container grün endet, ohne dass vorher `make gates` lief (Beleg im Bericht); sonst eingetreten → Rückführung nach `open` (§4).
- **Die Host-Liste ist zur Fahrt leer**, weil `container-start-gegenprobe` doch nicht in ihr steht (Kettenumbau `slice-harness-gegenproben-aus-gates`, Risiko dort; Beleg: `grep`). — **Ausgang (vorgesehen):** entfallen, wenn beim Start `grep container-start-gegenprobe harness/mk/*.mk` sie in der Host-Liste zeigt; sonst eingetreten → Rückführung nach `open` (§4).
- **Die Größe liegt unter der Marke, die Spanne ist breit** (Schätzung 350 bis 600, Marke 900); sie hängt an der Zahl der Fälle der Gegenprobe. — **Ausgang (vorgesehen):** entfallen, wenn die Messung nach der Metrik unter 900 bleibt; darüber meldet der Implementer (§4 Rückführung *zu groß*).
- **Die Pflicht zum Fahren wird nicht gelesen oder nicht befolgt** (Risiko 1 und 2 von `slice-harness-gegenproben-aus-gates`; dieser Slice liefert die Ziele, die sie nennt). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`.

**Fragen an den Nutzer:** Keine neue. Der Wortlaut der Pflicht ist entschieden (§1); ob `make gates-schnell` bleibt (Frage 2) und ob ein Stand der Gegenproben entsteht (Frage 4) stehen in `slice-harness-gegenproben-aus-gates` §6.

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
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-harness-gegenproben-aus-gates-sammelziele.md` in `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/` bei der Closure (Messung der Fahrten)
- **Folge-Slices:** `slice-harness-gegenproben-aus-gates-fahrstelle` (Rollen-Anweisungen, Gegenproben-Container, Release-Pflicht — ist eine Datei in `open/`)
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Hintergrund; der Slice liefert die Ziele, die die Gegenproben fahren), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (DoD 2 trägt je Element einen Fall und eine Mutation), `BEO-REPO/spec-randform-erst-im-review-entschieden` (§6 nennt die Randformen vor dem Code), `BEO-REPO/slice-waechst-durch-uebernahmen` (der Slice übernimmt nichts; der Schnitt steht in der Welle-Zeile), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groessenprobe-ohne-messbaren-stand` (§4 rechnet in der Metrik mit dem Doppelten für Neues). Keiner der Treffer erreicht mit diesem Slice 3× neu.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Liefer-Punkte 3 (Sammelziele mit Skript | Gegenprobe der Sammelziele | Gate-Index und Satz); der Slice liegt an der Grenze, ein vierter Punkt hat keinen Platz. Schichten 1: Harness-Werkzeug mit Gegenprobe und Gate-Index; `AGENTS.md` §6 ist ein Satz im selben Harness-Artefakt, Planung, Vertrag und ADR sind keine Schicht; Rollen-Anweisungen und Wartungs-Doku liegen in `slice-harness-gegenproben-aus-gates-fahrstelle`. Schätzung in der Metrik von §4: 350 bis 600 (Marke 900, Meldeschwelle 700). Der Plan trägt Sendungen ein an `slice-harness-gegenproben-aus-gates` (Rest des Satzes in `AGENTS.md` §6 und die Sammelziele sind abgegeben; Nachzählung dort in §8), `slice-harness-gegenproben-aus-gates-fahrstelle` (setzt diesen Slice voraus; Nachzählung dort in §8) und `slice-harness-gegenprobe-baumkopie` (`make gegenproben` seit diesem Slice; Nachzählung dort in §8).
