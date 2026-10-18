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

**Ziel:** `make gegenproben` (Namen laut `SPEC-060` Punkt 4) fährt die Gegenproben der Container-Liste (`GEGENPROBEN_CONTAINER`, alle ohne `u-boot` und `devcontainer`), `make gegenproben-host` die der Host-Liste (`GEGENPROBEN_HOST`, heute `container-start-gegenprobe`), beide über ein Skript (`tools/harness/gegenproben.sh`, Fragment `harness/mk/gegenproben.mk`): nacheinander in der Reihenfolge der Liste, alle auch nach einem roten Einzelziel, Exit 1 des Skripts bei mindestens einem roten (über `make` ist nur Exit ungleich 0 zugesagt: `make` gibt einem scheiternden Rezept Exit 2, bei `-i` Exit 0; Absatz vor Punkt 4 von [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md)), eine Zusammenfassung mit den Namen, je Einzelziel eine Zeile mit Name und Ausgang; der Standardeingang jedes Einzelziels ist `/dev/null`. Beide Sammelziele sind kein Gate und kein Teil von `make gates`; kein Hook und kein Stand glaubt ihr Ergebnis. Eine Gegenprobe der Sammelziele (`make gegenproben-gegenprobe`, an der Container-Liste, mit Attrappen als Einzelziele, ohne Docker) trägt die Fälle und Mutationen aus `SPEC-060` Teil 2: 44 Fälle und 48 Mutationen (der Architect rechnete vor dem Code mit 37 Fällen bei 16 Randformen der Gruppe G; die Nacharbeit zum Review fügte Fälle zu Punkt 4, 8 und 11 und Mutationen je Signal und je Element hinzu). Die Exit-Codes prüft die Gegenprobe am Skript, über `make` nur Exit ungleich 0.

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

- [ ] **Liefer-Punkt 1 (Sammelziele):** `make gegenproben` und `make gegenproben-host` (Namen laut `SPEC-060`) fahren je ihre Liste nacheinander, alle auch nach einem roten Einzelziel, enden mit Exit ungleich 0 bei mindestens einem roten (das Skript mit 1, über `make` nur ungleich 0), schreiben die Zusammenfassung mit den Namen, brechen bei `SIGINT` und `SIGTERM` ohne Zusammenfassung ab (das Skript beendet das laufende Einzelziel mit seinen Nachkommen), lehnen eine leere Liste und `-i` ab (das Skript mit Exit 2) und fahren unter `-n`, `-q`, `-t` nichts; jede Entscheidung je Randform aus §6 (Gruppe G) steht im Vertrag (`SPEC-060` Punkte 4 bis 12), vor dem ersten Code-Commit. Fragment `harness/mk/gegenproben.mk` und Skript `tools/harness/gegenproben.sh`; beide Ziele stehen nicht an `GATE_CHECKS` und nicht in der Ausnahmeliste.
- [ ] **Liefer-Punkt 2 (Gegenprobe der Sammelziele):** `make gegenproben-gegenprobe` (Fragment unter `harness/mk/`, an `GEGENPROBEN_CONTAINER` und ausdrücklich nicht an `GATE_CHECKS`, Zeile im Gate-Index) trägt je Randform und je Element ihrer Aufzählungen aus §6 einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10): 44 Fälle und 48 Mutationen, darunter je Signal und je Adressat des Abbruchs, die Nachkommen des Einzelziels, sechs Fälle zu Punkt 11 (grün, rot, abgebrochen, `-n`, `-q`, `-t`), der Standardeingang und die Sammelziele weder in `GATE_CHECKS` noch in der Ausnahmeliste; die Gegenprobe verlangt je Mutation genau die Menge roter Fälle, die sie nennt (Namen der Fälle tragen die Nummer des Punktes). Attrappen als Einzelziele, kein Docker.
- [ ] **Liefer-Punkt 3 (Gate-Index und Satz):** `harness/README.md` §Sensors führt die Zeilen der zwei Sammelziele (als Werkzeuge, kein Gate) und der Gegenprobe der Sammelziele (nur, was Fälle prüfen: Flags einzeln auf der Kommandozeile und aus der Umgebung, `-k` und `-j` nur auf der Kommandozeile, Exit-Codes am Skript, `make -i gegenproben` mit Exit 0); `AGENTS.md` §6 Schritt 6 trägt den Rest des Satzes (Ende jeder Welle und Release, `make gegenproben` im Gegenproben-Container und `make gegenproben-host` am Host) in dem Wortlaut, den der Nutzer am 2026-10-18 entschieden hat.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `harness/` und `tools/` ändern) und `make gates-schnell-gegenprobe` grün (die Deckung zählt die neue Gegenprobe); `make gegenproben` grün auf dem Endstand, **gefahren im Gegenproben-Container, falls er vorliegt, sonst am Host** (der Bericht nennt, wo und die Dauer); `make gegenproben-host` grün am Host. Gemessen (§7): `make gates` 267 s im Container `pgr-dev2` (Stand `f5b4abad`), `make gegenproben` 1095 s im Container `pgr-dev2` (Stand `0fe98aec`; danach änderten sich nur die Zeile `trap - INT TERM` in `gegenproben.sh` und die Gegenprobe der Sammelziele, die `make gegenproben` selbst fährt, 42 s), `make gegenproben-host` 255 s am Host (Stand `f5b4abad`). Dieser Container-Lauf folgte auf `make gates` im selben Container; der Beleg für das Risiko 1 aus §6 fehlt (Ausgang dort).
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
| `harness/mk/gegenproben.mk`, `tools/harness/gegenproben.sh` | neu | die zwei Sammelziele, das Skript, das eine Liste nacheinander fährt (Einzelziele im Hintergrund mit Standardeingang `/dev/null`, damit ein Signal wirkt; beim Abbruch endet das Einzelziel mit seinen Nachkommen), sammelt und nach `MAKEFLAGS` entscheidet; braucht `pgrep`; fängt `INT` und `TERM` selbst (ein zurückgesetztes `trap - INT TERM` im Unterprozess bewirkt nichts und steht nicht darin) |
| `tools/harness/gegenproben-gegenprobe.sh`, `harness/mk/gegenproben-gegenprobe.mk` | neu | Gegenprobe des Sammelskripts in Temp-Bäumen mit Attrappen als Ziele (kein Docker; bash ≥ 5, make, setsid, pgrep, timeout): 44 Fälle und 48 Mutationen, je Randform und je Element ein Fall und eine Mutation; hängt an `GEGENPROBEN_CONTAINER`, nicht an `GATE_CHECKS` |
| `harness/README.md`, `AGENTS.md` | update | Gate-Index (Zeilen der Sammelziele und ihrer Gegenprobe; im Absatz über den Sensors der Satz zum Ende jeder Welle und zum Release), Rest des Satzes in Schritt 6 |

**Arbeitsort.** Der Implementer arbeitet im Container `pgr-dev2` (Port 2223), Klon `/workspaces/pgr-dev2/pgwire-recorder`; dort laufen `make gates`, `make gegenproben-gegenprobe` und der Container-Teil, `make gegenproben`. Das Host-Ziel `make gegenproben-host` läuft am Host und wird in einem Wegwerf-Worktree des eingelieferten Commits geprüft, nicht im Arbeitsbaum des Containers.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` setzt dieser Übergang; der Nutzer priorisiert und trägt den Anspruch. **Vorbedingungen (beim Anspruch am 2026-10-18 am Stand `998738e8` geprüft: (a) bis (d) erfüllt, (e) mit diesem Übergang):** (a) Vertrag (`SPEC-060` Teil 2) und ADR liegen vor, mit der Entscheidung je Randform aus §6 (`AGENTS.md` §3.12), den Namen der Sammelziele und Listen; erfüllt, Accepted am 2026-10-18 (Zweig `arch-gates`, Commits des Architect vor dem Move nach `in-progress/`). (b) **Start-Trigger:** `slice-harness-gegenproben-aus-gates` liegt in `done/` (die Listen und die Namensregel stehen; die Host-Liste enthält `container-start-gegenprobe`, sonst ist sie leer, §6, G7). (c) **Reihenfolge:** Der Slice ist aus der Reihenfolge der Wellen gelöst und wartet nicht auf `welle-harness-regeln`; in der Welle bleibt er an zweiter Stelle. (d) **Wortlaut der Pflicht in `AGENTS.md` §6:** als erteilt eingetragen, Entscheidung des Nutzers vom 2026-10-18 (Wortlaut in §1); der Kettenumbau (`done/`) trägt die erste Hälfte bis zum Semikolon, dieser Slice den Rest. (e) **WIP:** Der Slice läuft auf eigenem Zweig; `Verantwortlich:` trägt beim Start Person und Zweig.

**Basis der Größenmetrik: 09d35ed8.** **Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness AGENTS.md`. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; Vertrag und ADR (Commits des Architect vor dem Start) zählen nicht.

**Schätzung (Architect, Prüfung vom 2026-10-18; Metrik wie oben; Gegenprobe-Code mit dem Doppelten des Produktcodes; Doku-Anteil Planner):**

| Teil | Zeilen |
|---|---|
| `gegenproben.mk` und `gegenproben.sh` (zwei Ziele, Skript mit `MAKEFLAGS`, Schleife, Sammlung, Signale) | 80 bis 140 |
| `gegenproben-gegenprobe` (37 Fälle mit Mutationen, 16 Randformen; Attrappen kosten Zeilen) | 250 bis 420 |
| `harness/README.md` (Zeilen der Sammelziele und der Gegenprobe) und `AGENTS.md` (Rest des Satzes) | 20 bis 40 |
| **Zusammen** | **350 bis 600** |

Meldeschwelle 700, **Marke 900**; die Obergrenze liegt unter der Meldeschwelle. Zählt der Architect später mehr als 16 Randformen in Gruppe G oder mehr als 37 Fälle, rechnet er neu.

**Gemessen (Closure, 2026-10-18):** `git diff --numstat 09d35ed8..076db09b -- tools harness AGENTS.md` = Einfügungen 484, Löschungen 6, Summe **490** (wörtlich); ohne die Änderung des Vertrags durch den Architect nach der Basis (`harness/vertraege`, 21 Zeilen) **469**. Unter der Schätzung (350 bis 600), unter der Meldeschwelle 700 und der Marke 900; die Rückführung *zu groß* trat nicht ein. Die Gegenprobe trägt mit 381 Zeilen gut vier Fünftel (381 von 469) und liegt, mit 44 statt 37 Fällen, im oberen Teil ihrer Spanne (250 bis 420).

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

**Randformen des neuen Vertrags (`AGENTS.md` §3.12).** Je Randform steht die Empfehlung des Planners und, wo entschieden, der Ort der Entscheidung (Zeile „Entschieden“): Der Architect legte sie vor dem Code in den Vertrag (`harness/vertraege/`, `SPEC-060` Teil 2), und der Implementer entschied keine still (Randform-Rückgabe; fünf Rückgaben, vom Architect am 2026-10-18 im Vertrag entschieden). Was aus dem Bestand entscheidbar war, steht als „Bestand“. Die Gruppe K (Kette) gehört `slice-harness-gegenproben-aus-gates`.

**Gruppe G — Sammelziele**

1. **Zwei Ziele, zwei Listen, ein Skript.** Name der Ziele und der Listen: Vorschlag des Architect. Beide Ziele sind kein Gate, hängen nicht an `GATE_CHECKS`, hängen nicht aneinander; das Container-Ziel fährt die Host-Liste nicht. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 1 und 4: `gegenproben` und `gegenproben-host`, Listen `GEGENPROBEN_CONTAINER` und `GEGENPROBEN_HOST`, Skript `tools/harness/gegenproben.sh`, Fragment `harness/mk/gegenproben.mk`.
2. **Reihenfolge.** Nacheinander in der Reihenfolge der Liste (Reihenfolge der Includes, alphabetisch nach Fragmentdatei), nie parallel, auch nicht unter `-j` (eine Schleife im Skript statt eines Fächers von Prerequisites): Die Last (Docker, Load 57 auf 20 Kernen, gemessen in [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md)) ist der Grund für das Herausnehmen. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 5. Abweichung: Zugesagt ist die Reihenfolge, in der `make` die Liste ausgibt, nicht „alphabetisch nach Fragmentdatei“ (Grund: die Reihenfolge der Dateien beim `include` mit Muster ist Sache von `make`).
3. **Rotes Einzelziel.** Alle Ziele der Liste fahren, auch nach einem roten (wie `-k`); die Zusammenfassung nennt die roten beim Namen, Exit 1. Ein Abbruch durch `SIGINT` oder `SIGTERM` beendet das laufende Ziel und startet kein weiteres (Exit 130 und 143). Je ein Fall. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 6 und 8; ein Abbruch schreibt keine Zusammenfassung. Exit 1, 130 und 143 sind die des Skripts (Absatz vor Punkt 4, Entscheidung des Architect nach der Rückgabe 1 des Implementers); über `make` endet ein rotes Sammelziel mit Exit 2, zugesagt ist nur „ungleich 0“. Das Abbruchsignal gilt der Prozessgruppe und dem Skript allein; das laufende Einzelziel endet mit seinen Nachkommen (`pgrep -P`).
4. **Aufrufarten.** `-n`, `-q`, `-t`: Das Rezept trägt weder `$(MAKE)` noch `+`, make führt es nicht aus (wie `record-gates`); `-i`: Die Einzelziele erben `-i` und gäben ein rotes Skript als Exit 0 zurück, deshalb lehnt das Skript `-i` in `MAKEFLAGS` ab (Exit 2, Meldung), statt grün zu melden; `-k`, `-j` ohne Wirkung. `-i` zusammen mit `-n`, `-q` oder `-t`: `-n`, `-q`, `-t` gelten, nichts läuft. Kein Stand wird geschrieben. Je Flag ein Fall, auch aus der Umgebung und kombiniert. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 10, Reihenfolge der Gründe Punkt 12. Exit 2 bei `-i` ist der des Skripts; über `make` endet `make -i gegenproben` mit Exit 0 und der Meldung (Absatz vor Punkt 4). Abweichung: Das Rezept übergibt `$(MAKE)` an das Skript wie `gates-schnell`, `make` führt es also auch unter `-n`, `-q`, `-t` aus, und das Skript entscheidet an `MAKEFLAGS` (Grund: Das Skript braucht das `make` des Aufrufers, und nur so hat es es; zugesagt ist das Ergebnis: kein Einzelziel läuft).
5. **Fehlendes Werkzeug.** Fehlt `docker`, `bash` oder `setsid`, ist das betroffene Einzelziel rot mit seiner eigenen Meldung, das Sammelziel fährt weiter; kein Vorab-Test der Werkzeuge. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 6 (Fall `p6-127`).
6. **Ausgabe.** Die Ausgabe der Einzelziele erscheint unverändert; die Zusammenfassung steht zuletzt und nennt je Ziel Name und Ausgang. Ein Überspringen innerhalb eines Ziels (das Überspringen der ersetzten ADR, ersetzt durch [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md)) gibt es nach der Aufteilung nicht mehr. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 6 und 7; die Zeilen der Zusammenfassung sind in Punkt 6 festgelegt, `<n>` ist der Exit von `<make> <ziel>` (unter echtem `make` 2). Der Standardeingang jedes Einzelziels ist `/dev/null` (Punkt 7, akzeptiertes Negativ ohne Fall).
7. **Leere Liste.** Die Host-Liste ist ohne `container-start-gegenprobe` leer (§4 *Start*); eine leere Liste ist ein Fehler (Exit 2, Meldung), kein stilles Grün. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 9 (Exit 2 des Skripts; über `make` ungleich 0).
8. **Eine Gegenprobe steht in beiden Zielen oder in keinem.** Das Sammelskript prüft die Zuordnung nicht, die Gegenprobe der Kette tut es (Gruppe K, K1 in `slice-harness-gegenproben-aus-gates`). Entschieden: [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 2.
9. **Aufruf.** Aus einem Unterverzeichnis des Repos wie von der Wurzel; ohne Git-Arbeitsbaum Exit 2. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): Abweichung von der Empfehlung: [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 13 sagt weder einen Aufruf aus einem Unterverzeichnis noch Exit 2 ohne Git-Arbeitsbaum zu (Grund: Die Schnittstelle ist das Make-Ziel, das Rezept läuft im Verzeichnis des `Makefile`, und das Skript ruft kein `git`).
10. **Gleichzeitige Läufe und Worktrees.** Zwei Fahrten im selben Arbeitsbaum oder zusammen mit `make gates` sind nicht zugesagt (`SPEC-052` Punkt 8, Aufruf-Voraussetzung nacheinander); zwei Worktrees am selben Docker-Daemon teilen Namen von Containern und Netzen (Präfixe `pgr-ivg-`, `pgr-probe-`), die Gegenprobe der Container schützt nur durch ihre Präfixe. Der Vertrag bindet die Zusage an „nacheinander“; die Fahrstelle im eigenen Container (`slice-harness-gegenproben-aus-gates-fahrstelle`) hat einen eigenen Daemon. **Aufruf-Voraussetzung (§3.12):** Das Sammelskript sagt keine Ordnung über mehrere Aufrufe zu. Entschieden (Architect, 2026-10-18; Empfehlung bis zur Annahme von [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 13.
11. **Dev-Container-Erkennung entfällt.** Keine Erkennung, kein Überspringen; die Zuordnung (K2 in `slice-harness-gegenproben-aus-gates`) ersetzt sie. Entschieden: Nutzer, 2026-10-18; [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) Accepted. Im Vertrag: [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 3.
12. **Plattform des Skripts.** Das Skript läuft ab bash 4 und braucht `pgrep` (procps); die Gegenprobe der Sammelziele braucht zusätzlich `setsid` und `timeout` und bash ≥ 5. Entschieden (Architect, 2026-10-18): [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Kopf, *Plattform*.

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

**Belege des Implementers** (Stand 2026-10-18; Code-Commit `0fe98aec`, Nacharbeit zum Review in `f5b4abad`):

- **Metrik** (`AGENTS.md` §3.14): `git diff --numstat 09d35ed8..f5b4abad -- tools harness AGENTS.md` = Einfügungen 484, Löschungen 6, Summe **490** (zuvor 428 am Commit `0fe98aec`). Ohne die Änderung des Vertrags durch den Architect nach der Basis (`harness/vertraege`: 17 und 4) ist sie **469**: `gegenproben-gegenprobe.sh` 381, `gegenproben.sh` 53, `gegenproben.mk` 17, `gegenproben-gegenprobe.mk` 11, `harness/README.md` 5, `AGENTS.md` 2. Unter der Meldeschwelle 700 und der Marke 900.
- **Läufe am Endstand `f5b4abad`:** `make gates` grün, 267 s (Container `pgr-dev2`); `make gegenproben-gegenprobe` grün, 42 s (vier Mutationen gleichzeitig); `make gegenproben-host` am Host in einem Wegwerf-Worktree des eingelieferten Commits grün, 255 s (vorher keine Gegenprobe in `pgr-dev3`; Worktree entfernt, kein Container und kein Volume mit Präfix `pgr-cs-` geblieben); `make docs-check` 0 Befunde. `make gegenproben` im Container lief zuletzt am Commit `0fe98aec` grün (1095 s, alle 13 Einträge; danach änderten sich `gegenproben.sh` um eine Zeile und die Gegenprobe der Sammelziele, die sie selbst fährt, `make gates-schnell-gegenprobe` unverändert). **Dieser Lauf folgte auf `make gates` im selben Container und war kein frischer:** das Risiko 1 aus §6 bleibt unbelegt (Befund F-10 des Review, geht an den Planner).
- **Zusage · Mutation · roter Test** (Textersatz in der Kopie, `mut` in `tools/harness/gegenproben-gegenprobe.sh`, 48 Mutationen, 44 Fälle; die Gegenprobe verlangt je Mutation genau die genannten Fälle rot, mit `GP_ENTDECKE=1` einmal gesehen; „+" = Folgefälle derselben Ausgabe):

| Punkt | Zusage | Mutation | rote Fälle |
|---|---|---|---|
| 4 | Container-Ziel fährt nur seine Liste | fährt auch die Host-Liste; hängt an `gegenproben-host` | `p4-container` (+ `p9-container`, `p10-j`, `p10-k` bei der ersten) |
| 4 | Host-Ziel fährt nur seine Liste | fährt auch die Container-Liste; hängt an `gegenproben` | `p4-host` (+ `p9-host`) |
| 4 | kein Gate | `GATE_CHECKS += gegenproben`; `+= gegenproben-host` | `p4-kein-gate` |
| 4 | nicht in der Ausnahmeliste | `GATE_SCHNELL_AUSNAHME += gegenproben`; `+= gegenproben-host` | `p4-kein-ausnahme` |
| 5 | nacheinander in Reihenfolge, auch unter `-j` | Hintergrund; rückwärts; Prerequisites | `p5-reihenfolge`, `p5-j` (+ Zusammenfassung, Abbruch); Prerequisites `p5-j` (+ `p9-container`, `p10-i*`, `p10-j`, `p10-k`) |
| 6 | alle laufen nach einem roten | Abbruch beim ersten; ohne `k` | `p6-erstes-rot`, `p6-mittleres-rot`, `p6-zwei-rot`, `p6-127`, `p10-j`, `p10-k` (ohne `k`: nicht `p10-k`) |
| 6 | Exit 1 bei rotem | Exit 0; Exit des letzten Ziels | `p6-*-rot`, `p6-127`, `p10-j`, `p10-k`; letzter: nur wo ein früheres rot, das letzte grün ist |
| 6 | Zusammenfassung vollständig, mit Exit-Code, zuletzt | nur die roten; ohne Exit-Code; vor der Ausgabe des letzten Ziels | `p6-alle-gruen` (erste und dritte), `p6-*-rot`, `p6-127`, `p10-j`, `p10-k` |
| 7 | Ausgabe unverändert, gleicher Strom | stderr nach stdout; stdout verworfen | `p7-ausgabe` |
| 7 | Standardeingang `/dev/null` | `< /dev/zero` | `p7-eingabe` |
| 8 | Abbruch: kein weiteres Ziel, keine Zusammenfassung | nächstes startet; Zusammenfassung erscheint | alle vier `p8-*` |
| 8 | 130 für `SIGINT`, 143 für `SIGTERM` | beide vertauscht; nur `INT` nach 143; nur `TERM` nach 130 | alle vier; `p8-sigint-gruppe`, `p8-sigint-skript`; `p8-sigterm-gruppe`, `p8-sigterm-skript` |
| 8 | Signal nur an das Skript wirkt | Trap für `INT` gestrichen; Trap für `TERM` gestrichen | nur `p8-sigint-skript`; nur `p8-sigterm-skript` |
| 8 | das laufende Einzelziel endet | Einzelziel bleibt am Leben | `p8-sigint-skript`, `p8-sigterm-skript` |
| 8 | seine Nachkommen enden (die Attrappe startet ein Hintergrundkind) | Schleife über `pgrep -P` gestrichen | `p8-sigint-skript`, `p8-sigterm-skript` |
| 9 | leere Liste: Exit 2, Meldung, kein Ziel | Exit 0; Listen vertauscht | `p9-container`, `p9-host` (vertauscht: + `p4-*`, `p5-*`, `p6-erstes-rot`, `p10-j`, `p10-k`) |
| 10 | `-n`, `-q`, `-t` fahren nichts | je Buchstabe übersehen | `p10-n` (+ `-env`, `p10-in`, `p12-n-leer`); `p10-q` (+ `-env`, `p10-iq`); `p10-t` (+ `-env`, `p10-it`) |
| 10 | `-i` lehnt ab (Exit 2, Meldung) | `i` fällt aus; `i` wie `n` | die fünf `p10-i*` und `p12-i-leer` |
| 10 | `n`, `q`, `t` gelten vor `i` | `i` vor `n` | `p10-in`, `p10-iq`, `p10-it` |
| 10 | nur das erste Wort von `MAKEFLAGS` | ganzes `MAKEFLAGS` | `p5-j`, `p10-i-j`, `p10-j` |
| 10 | `-j` ändert nichts | `j` gleichzeitig | `p5-j`, `p10-j` |
| 11 | kein Stand, je Element | schreibt nach grünem Lauf; bei rotem; beim Abbruch; nur unter `-n`; nur unter `-q`; nur unter `-t` | `p11-gruen`; `p11-rot`; `p11-abbruch`; `p11-n`; `p11-q`; `p11-t` (je nur der eine) |
| 12 | Reihenfolge der Gründe | Liste vor `-n`; Liste vor `-i` | `p12-n-leer`, `p12-i-leer`; `p12-i-leer` |

- **Nachgeführt zum Review:** F-1 Attrappe mit Hintergrundkind in allen vier Abbruch-Fällen, Mutation `m-p8-nachkommen-bleiben`; F-2 sechs Fälle und Mutationen zu Punkt 11 (statt eines Falls); F-3 vier Mutationen je Signal (`m-p8-int-nach-143`, `m-p8-term-nach-130`, `m-p8-trap-int-gestrichen`, `m-p8-trap-term-gestrichen`); F-4, F-9 Gate-Index (Flags, Voraussetzungen, `make -i` Exit 0, Standardeingang); F-5 Fall `p4-kein-ausnahme` und Mutation für `gegenproben-host` in `GATE_CHECKS`; F-11 siehe nächster Punkt.
- **F-11, `trap - INT TERM` im Unterprozess:** Die Mutation „Zeile gestrichen" färbte keinen Fall. Grund (gemessen): Das Skript trapt `INT` und `TERM` vor dem Start, und gefangene Signale setzt `exec` auf den Standard zurück; die Zeile bewirkte nichts. Ich habe sie aus `tools/harness/gegenproben.sh` entfernt (damit auch ein Host-Lauf am geänderten Skript). Unterstützt ist, dass das Skript mit nicht ignoriertem `SIGINT` startet (ein im Hintergrund einer nicht interaktiven Shell gestartetes Skript kann `INT` nicht fangen; die Gegenprobe setzt es dafür zurück).
- **Invariante ohne Mutation:** Der Fall `p10-i*` verlangt über `make -i` Exit 0 (Verhalten von `make`, kein Code dieses Slice); eine Mutation, die ihn bricht, gibt es nicht.
- **Grüne Mutanten:** keine. Ohne Fall bleiben ein Signal nur an `make` und `SIGKILL` (Punkt 15).
- **Rückgaben:** Die fünf Rückgaben aus dem ersten Lauf hat der Architect am 2026-10-18 im Vertrag entschieden (`cff2d59a`); offen bleibt F-10 (Risiko 1, frischer Container) beim Planner.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Hintergrund; der Slice liefert die Ziele, die die Gegenproben fahren), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (DoD 2 trägt je Element einen Fall und eine Mutation), `BEO-REPO/spec-randform-erst-im-review-entschieden` (§6 nennt die Randformen vor dem Code), `BEO-REPO/slice-waechst-durch-uebernahmen` (der Slice übernimmt nichts; der Schnitt steht in der Welle-Zeile), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groessenprobe-ohne-messbaren-stand` (§4 rechnet in der Metrik mit dem Doppelten für Neues). Keiner der Treffer erreicht mit diesem Slice 3× neu.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Liefer-Punkte 3 (Sammelziele mit Skript | Gegenprobe der Sammelziele | Gate-Index und Satz); der Slice liegt an der Grenze, ein vierter Punkt hat keinen Platz. Schichten 1: Harness-Werkzeug mit Gegenprobe und Gate-Index; `AGENTS.md` §6 ist ein Satz im selben Harness-Artefakt, Planung, Vertrag und ADR sind keine Schicht; Rollen-Anweisungen und Wartungs-Doku liegen in `slice-harness-gegenproben-aus-gates-fahrstelle`. Schätzung in der Metrik von §4: 350 bis 600 (Marke 900, Meldeschwelle 700). Der Plan trägt Sendungen ein an `slice-harness-gegenproben-aus-gates` (Rest des Satzes in `AGENTS.md` §6 und die Sammelziele sind abgegeben; Nachzählung dort in §8), `slice-harness-gegenproben-aus-gates-fahrstelle` (setzt diesen Slice voraus; Nachzählung dort in §8) und `slice-harness-gegenprobe-baumkopie` (`make gegenproben` seit diesem Slice; Nachzählung dort in §8).
