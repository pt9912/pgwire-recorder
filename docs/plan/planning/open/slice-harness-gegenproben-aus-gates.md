# Slice slice-harness-gegenproben-aus-gates: Die Gegenproben und `a-check-negativ` verlassen `make gates`; zwei Sammelziele fahren sie, eines im Gegenproben-Container, eines am Host

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten, erste Stelle der Reihenfolge, vor `slice-harness-gegenprobe-baumkopie` (Entscheidung des Nutzers vom 2026-10-18, sofort). **Start-Trigger abweichend von der Welle:** `slice-harness-devcontainer-einrichten-geruest` in `done/`, nicht `welle-harness-regeln` in `done/` (Auftrag des Nutzers „Bitte sofort beheben“, Lösung aus der Reihenfolge der Wellen vom Koordinator am 2026-10-18 angewiesen; Umplanung im Drift-Log der Roadmap); der Start-Trigger der Welle gilt für die übrigen Slices. Angelegt am 2026-10-18 nach den Dialogantworten des Nutzers desselben Tages: „Das dauert zu lange und ich verstehe auch nicht, warum das nötig ist, wenn sich nichts geändert hat“, „Der schnelle fix besteht darin es aus make gates herauszunehmen“, „Alle Gegenproben raus“, „Man kann es explizit in einem eigenen Dev-Container fahren, dann werden die anderen nicht von der eigentlichen Arbeit abgehalten“ und „Wir können Gegenproben aufteilen. Einen speziellen, der nur im Host (u-boot und devcontainer) läuft (zwei make targets)“ (Wortlaut vom Koordinator übermittelt). Der erste Zuschnitt des Planners („läuft nur bei geänderten Eingaben“) ist damit nicht der Slice für jetzt (§1, *Ausdrücklich NICHT*). Der Slice liegt in `open/`: Er ist nicht priorisiert, den Anspruch trägt der Nutzer; Vertrag und ADR des Architect liegen nicht vor (§4 *Start*); die ADR, die [ADR-0046](../../adr/0046-gegenprobe-ueberspringt-im-dev-container.md) ersetzt ([ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md), Accepted auf dem Zweig `slice-harness-devcontainer-einrichten-geruest`), liegt vor.

**Bezug:** [LH-QA-07](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur insoweit, als die Gegenproben dessen Messungen tragen (Abdeckungs-Deklaration im Kopf der Gegenproben); kein Anspruch des Produkts, der Slice ändert die Gate-Kette des Harness. Bindung an das Register: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Beleg `evidence/slice-harness-gegenproben-aus-gates.md`, mit diesem Plan angelegt). Aktive ADRs: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (wird teilweise überholt, §4 *Start*), [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) (Aufteilung in Container- und Host-Ziel; ersetzt [ADR-0046](../../adr/0046-gegenprobe-ueberspringt-im-dev-container.md)).

**Berührte Spec-Stellen:** neuer Vertrag der Sammelziele und der Namens- und Zuordnungsregel in `harness/vertraege/` (Nummer vergibt der Architect; nächste freie Nummer nach `SPEC-059`, also `SPEC-060`; `SPEC-058` liegt auf dem Zweig von `slice-harness-devcontainer-einrichten-geruest`, noch nicht auf main); [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) Punkte 1, 2, 3, 4, 9 und 13; [`SPEC-053`](../../../../harness/vertraege/SPEC-053-postgres-versionen.md) Punkte 12 und 13; [`SPEC-054`](../../../../harness/vertraege/SPEC-054-meldungskatalog.md) Punkt 11; [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 12; [`SPEC-056`](../../../../harness/vertraege/SPEC-056-testabdeckung-schwelle.md) Punkt 11; [`SPEC-057`](../../../../harness/vertraege/SPEC-057-integration-aufraeumen.md) Punkt 10; [`SPEC-047`](../../../../harness/vertraege/SPEC-047-kopf-der-plaene.md) Punkt 9 und [`SPEC-049`](../../../../harness/vertraege/SPEC-049-lint-profil.md) Punkt 10 nur in dem, was ihre Gegenprobe zum Anschluss prüft; `SPEC-058` (Stand des Zweigs) Punkte zur Ausnahmeliste und zur Erkennung des Dev-Containers.

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

**Ziel:** `make gates` fährt keine Gegenprobe mehr: Jedes Gate, dessen Name auf `-gegenprobe` endet, und `a-check-negativ` verlassen `GATE_CHECKS`. Sie bleiben als Ziele bestehen, hängen an einer von zwei neuen Listen, und je ein Sammelziel fährt eine Liste nacheinander: `make gegenproben` (vorläufiger Name) fährt die Gegenproben, die im Dev-Container laufen, also alle ohne `u-boot` und `devcontainer`; `make gegenproben-host` (vorläufiger Name) fährt die Gegenproben, die diese beiden Werkzeuge brauchen, heute `container-start-gegenprobe`, später ihren Teil II. Die Namen schlägt der Architect vor. Beide Sammelziele sind kein Teil von `make gates`.

**Betroffen (Zählung am Stand `16e5cc7d`, Planner):** Aus `GATE_CHECKS` gehen zwölf Ziele in zwölf Fragmenten unter `harness/mk/`: elf Gegenproben (`abdeckung-gegenprobe`, `test-container-gegenprobe`, `coverage-gegenprobe`, `gates-schnell-gegenprobe`, `hook-gegenprobe`, `integration-volumes-gegenprobe`, `kopf-check-gegenprobe`, `lint-gegenprobe`, `meldungskatalog-check-gegenprobe`, `stop-hook-gegenprobe`, `test-postgres-versionen-gegenprobe`) und `a-check-negativ` (`arch-negativ.mk`); dazu `container-start-gegenprobe`, sobald `slice-harness-devcontainer-einrichten-geruest` (Zweig, noch nicht auf main) liegt (§4 *Start*). Im Gate bleiben `abdeckung-check`, `a-check`, `baseline-verify`, `build`, `docs-check`, `kopf-check`, `lint`, `meldungskatalog-check`, `test`, `test-container`, `test-integration` und `test-postgres-versionen`. Die Ausnahmeliste der schnellen Kette schrumpft damit auf `test-container`, `test-integration` und `test-postgres-versionen`. Der Stop-Hook, die Dateiliste und der Nachweis des vollen Laufs bleiben unverändert im Code.

**Ohne Kette bleibt eine Regression unbemerkt, bis jemand die Gegenprobe fährt.** Das ist eine Lockerung der Kette (`AGENTS.md` §3.6) und der Preis dieses Slice (§6, erstes Risiko). Er wird nicht durch einen Stand oder eine Eingabeliste gemildert (das war der erste Zuschnitt, unten), sondern durch eine Pflicht, wann die Gegenprobe zu fahren ist: bei einer Änderung am geprüften Werkzeug oder Gate, an der Gegenprobe selbst oder an deren Vertrag; am Ende jeder Welle; vor dem Release. Review und Verifikation fahren die betroffene Gegenprobe gezielt (Regel des engsten Sensors). Diesen Slice trägt davon der Satz in `AGENTS.md` §6 Schritt 6 und die Gate-Zeilen in `harness/README.md`; die Rollen-Anweisungen und die Fahrstelle liefert `slice-harness-gegenproben-aus-gates-fahrstelle` (Abgrenzung unten).

**Die Pflicht-Fahrstelle für `make gegenproben` ist ein eigener Dev-Container** (Entscheidung des Nutzers vom 2026-10-18, Name vorläufig `pgr-gp`): Er läuft im Hintergrund auf dem eingelieferten Stand eines Zweigs, während Implementer, Planner und Hook am Host und im Arbeitscontainer ungestört arbeiten. `make gegenproben-host` läuft am Host (`u-boot`, `devcontainer`). Der Container wird von Hand oder von einem Agenten gestartet und fährt nur die Gegenproben (in ihm wird kein Slice abgearbeitet); ein `post-receive`-Hook im Bare-Repo des Containers startet bei jeder Einlieferung eines Zweigs `make gegenproben` auf diesem Stand, die Closure von Welle und Release liest das Ergebnis. `make gegenproben-host` fährt, wer es braucht, am Host. Was sie voraussetzen und wie, steht in `slice-harness-gegenproben-aus-gates-fahrstelle`.

**Übernimmt:** nichts. Die Erkennung des Dev-Containers und das Überspringen der Fälle mit `u-boot` und `devcontainer` entfallen mit der Aufteilung (Entscheidung des Nutzers vom 2026-10-18); [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) (Accepted) ersetzt [ADR-0046](../../adr/0046-gegenprobe-ueberspringt-im-dev-container.md) und macht `container-start-gegenprobe` zu einem Ziel für den Host, das weder an `GATE_CHECKS` hängt noch in der Ausnahmeliste der schnellen Kette steht. `slice-harness-devcontainer-einrichten-geruest` (SPEC-058 und der Code der Gegenprobe) hat den Vertrag und den Code umgestellt (Zweig, Stand `8fa3dd24`); dieser Slice nimmt nichts davon mit, er hängt die Gegenprobe nur an die Host-Liste (unten, *Ausdrücklich NICHT*).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- „Die Gegenprobe läuft nur, wenn sich ihre Eingaben geändert haben“ (Hash der Eingaben, Stand unter `.harness/state/`, ungültig bei Änderung) — der erste Zuschnitt dieses Auftrags, vom Nutzer verworfen („Alle Gegenproben raus“). Er bräuchte je Gegenprobe eine vollständige Eingabeliste, und eine unvollständige ist ein stilles Grün ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), Befund 2); der Nutzer will die Kette lockern, nicht bedingen. **Möglicher späterer Slice**, falls die Pflicht nicht trägt; vermerkt als Beleg im Register (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `evidence/slice-harness-gegenproben-aus-gates.md`). Eine Adresse gibt es nicht; ein Slice legt der Planner an, wenn der Nutzer es entscheidet.
- Rollen-Anweisungen (`implement-slice`, `close-welle`, Reviewer- und Verifier-Skill, Implementer-Typ), der Abschnitt in `docs/developer/devcontainer.md` über den Gegenproben-Container und die Sendung an `slice-erster-release-freigabe` — `slice-harness-gegenproben-aus-gates-fahrstelle`, direkt hinter diesem Slice. Begründung: Schicht-Abgrenzung (Rollen-Anweisungen und Wartungs-Doku sind zwei weitere Schichten, `AGENTS.md` §3.13) und Größe; der Satz in `AGENTS.md` §6 hält die Pflicht in der Zwischenzeit.
- Ein Stand oder Nachweis der Gegenproben (analog `.harness/state/gates-passed.diffsha`) und ein Hook, der ihn liest — anderer Vorgang und eine Entscheidung des Nutzers: Der Hook glaubt heute den Nachweis von `make gates`, und ein zweiter Nachweis brächte die Laufzeit zurück, die der Nutzer streichen will. Nicht angelegt; Frage in §6 (Risiko 1).
- Der Code und der Vertrag der Gegenprobe des Container-Starts (SPEC-058, Erkennung des Dev-Containers entfernt) — `slice-harness-devcontainer-einrichten-geruest` hat sie nach [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) umgestellt; sie liegt dort, ihr Anschluss an die Host-Liste ist dieser Slice.
- Ob `make gates-schnell` überflüssig wird — nur genannt, nicht entschieden: Nach dem Umbau fährt `make gates` alles außer den Gegenproben, `make gates-schnell` davon alles außer `test-container`, `test-integration` und `test-postgres-versionen`; die zwei Ketten unterscheiden sich weiter durch diese drei Gates. Entscheidung des Nutzers (§6, Frage 2).
- Kleinere Kopien, `git ls-files`-Helfer, Reihenfolge der Gegenproben nach Laufzeit — `slice-harness-gegenprobe-baumkopie` (Kopie nach `git ls-files`); der Slice ändert, wo die Gegenproben hängen, nicht, was sie tun.
- Neue Fälle in den Gegenproben außer denen für Anschluss, Namens- und Zuordnungsregel und Sammelziele — andere Vorgänge; Bestand bleibt bewusst stehen.
- Produkt-Code, Handbuch, `README.md` — Schicht-Abgrenzung: kein Verhalten des gebauten Binary.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Kette):** Kein Gate mit Namen auf `-gegenprobe` und nicht `a-check-negativ` hängt an `GATE_CHECKS`; jedes hängt an genau einer der beiden Listen (Container-Liste, Host-Liste; Namen und Form im Vertrag), die Fragmente unter `harness/mk/` tragen die Zeile dafür; `GATE_SCHNELL_AUSNAHME` nennt nur noch `test-container`, `test-integration` und `test-postgres-versionen`. Die Namens- und Zuordnungsregel steht im Vertrag vor dem ersten Code-Commit und hat Fälle und Mutationen in `gates-schnell-gegenprobe`: je Element der Aufzählung (Gegenprobe-Name, `a-check-negativ`, Container-Liste, Host-Liste) ein Fall, der nur es bricht; die Fälle der Gegenproben, die heute den Anschluss an `GATE_CHECKS` prüfen (`container-gegenprobe`, `coverage-gegenprobe`, `kopf-check-gegenprobe`, `lint-gegenprobe`, `meldungskatalog-check-gegenprobe`, `postgres-versionen-gegenprobe`, `stop-hook-gegenprobe`, `gates-schnell-gegenprobe`), prüfen den Anschluss an die Liste.
- [ ] **Liefer-Punkt 2 (Sammelziele):** `make gegenproben` und `make gegenproben-host` (Namen laut Vertrag) fahren je ihre Liste nacheinander, alle auch nach einem roten Einzelziel, enden mit Exit ungleich 0 bei mindestens einem roten und nennen die roten beim Namen; die Entscheidung je Randform aus §6 (Gruppe G) steht im Vertrag, und eine Gegenprobe des Sammelskripts (Fragment unter `harness/mk/`, an `GATE_CHECKS` ausdrücklich nicht, sondern an der Container-Liste, Zeile im Gate-Index) trägt je Randform und je Element ihrer Aufzählungen einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10).
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `harness/` und `tools/` ändern), und `make gegenproben` grün auf dem Endstand, **gefahren im Gegenproben-Container, falls er vorliegt, sonst am Host** (der Bericht nennt, wo); `make gegenproben-host` grün am Host, sobald seine Liste nicht leer ist (§4 *Start*). Der Bericht nennt je Lauf die Dauer; `make gates` ohne Gegenproben ist damit gemessen (das Soll der Auftragsstellung, kein Vertragspunkt).
- [ ] Doku-Update: `harness/README.md` §Sensors (Zeile der Sammelziele, Bindung und Spalte „hängt an“ der zwölf Ziele, Zeile von `make gates-schnell` und `make gates`), `AGENTS.md` §6 Schritt 6 und §4 (ein Satz zur Pflicht; Wortlaut nach der Entscheidung des Nutzers, Empfehlung in §6, Frage 3).
- [ ] Beobachtungs-Register: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` mit dem Beleg dieses Slice geführt (liegt an, mit diesem Plan); das Risiko 1 aus §6 geht als neue Beobachtung ins Register (`BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`, 1×, Ausgang *weiter offen*); kein Zähler wird gesetzt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review. Das Review fährt die Gegenprobe der Kette und des Sammelskripts gezielt.
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
| `harness/vertraege/SPEC-<NNN>-gegenproben.md` und die Punkte aus dem Kopf in `SPEC-052`, `-053`, `-054`, `-055`, `-056`, `-057`, ggf. `-047`, `-049` | neu und update (Architect, vor dem Code) | Vertrag der Sammelziele, Namens- und Zuordnungsregel, Entscheidung je Randform aus §6; Nachzug der Aussagen „hängt an `GATE_CHECKS`“ und der Ausnahmeliste; Zählung im Kopf: zwölf bis vierzehn Punkte in acht Verträgen (Stand `16e5cc7d`, ohne `SPEC-058`) |
| neue ADR (Architect, Proposed) | neu | Lockerung der Kette (`AGENTS.md` §3.6); ergänzt oder ersetzt Teile von [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (Alternative D dort), nennt die ADRs, die „hängt an `GATE_CHECKS`“ tragen (Stand `16e5cc7d`: 0025, 0028, 0029, 0032, 0033, 0034, 0038, 0042, 0043, 0044, auf dem Zweig 0045 und 0046); [ADR-0046](../../adr/0046-gegenprobe-ueberspringt-im-dev-container.md) ersetzt [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) (liegt vor), die neue ADR führt [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) als Bezug |
| zwölf Fragmente unter `harness/mk/` (`arch-negativ`, `abdeckung`, `container`, `coverage`, `gates-schnell`, `hook-gegenprobe`, `integration`, `kopf-check`, `lint`, `meldungskatalog-check`, `postgres-versionen`, `stop-hook-gegenprobe`) | update | `GATE_CHECKS +=` wird zur Zeile der Liste; Kommentarköpfe, die „hängt an GATE_CHECKS“ sagen, folgen (`AGENTS.md` §3.7) |
| `harness/mk/gegenproben.mk`, `tools/harness/gegenproben.sh` | neu | die zwei Listen, beide Sammelziele, das Skript, das eine Liste nacheinander fährt, sammelt und nach `MAKEFLAGS` entscheidet |
| `tools/harness/gegenproben-gegenprobe.sh`, `harness/mk/gegenproben-gegenprobe.mk` | neu | Gegenprobe des Sammelskripts in Temp-Bäumen mit Attrappen als Ziele (kein Docker): je Randform ein Fall und eine Mutation |
| `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh` | update | Ausnahmeliste auf drei Einträge; Deckung: jedes Ziel mit Namen auf `-gegenprobe` und `a-check-negativ` steht in genau einer Liste und in keiner Hälfte von `GATE_CHECKS`; Fixture und Mutationen je Element |
| `tools/harness/container-gegenprobe.sh`, `coverage-gegenprobe.sh`, `kopf-check-gegenprobe.sh`, `lint-gegenprobe.sh`, `meldungskatalog-check-gegenprobe.sh`, `stop-hook-gegenprobe.sh`, `tools/test/postgres-versionen-gegenprobe.sh` | update | die Fälle und Mutationen zum Anschluss (heute 70 Zeilen mit `GATE_CHECKS` oder `GATE_SCHNELL_AUSNAHME`) prüfen den Anschluss an die Liste |
| `harness/README.md`, `AGENTS.md` | update | Gate-Index, Pflicht-Satz (Doku-Update der DoD) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. **Vorbedingungen:** (a) Der Architect hat Vertrag und ADR geschrieben, mit der Entscheidung je Randform aus §6 (`AGENTS.md` §3.12), den Namen der Sammelziele und der Listen und der Entscheidung, welche Teile von [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) ergänzt und welche ersetzt werden; ADR als Proposed, Annahme durch den Nutzer; Vertrag und ADR sind Commits des Architect vor dem Move nach `in-progress/`. (b) **Start-Trigger:** `slice-harness-devcontainer-einrichten-geruest` liegt in `done/` (die Host-Liste braucht `container-start-gegenprobe`, sonst ist sie leer, §6, G7; nach [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) hängt die Gegenprobe dort nicht an `GATE_CHECKS`). (c) **Reihenfolge:** Der Slice ist aus der Reihenfolge der Wellen gelöst (Auftrag „sofort“; Drift-Log der Roadmap): Er startet, sobald (a) und (b) gelten, und wartet nicht auf `welle-harness-regeln`; in der Welle bleibt er der erste Slice (§5 der Welle), der Start-Trigger der Welle gilt für die übrigen Slices. (d) **Wortlaut der Pflicht in `AGENTS.md` §6:** Der Nutzer sagt beim Start ausdrücklich Ja zum Wortlaut (Frage 3); ohne dieses Ja ändert der Slice `AGENTS.md` nicht. (e) **WIP:** Der Slice läuft auf eigenem Zweig; `Verantwortlich:` trägt beim Start Person und Zweig (`<Person> (Zweig <branch-name>)`, `modul-05-planning-harness.md` §Lifecycle als State Machine).

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness AGENTS.md`. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; Vertrag und ADR (Commits des Architect vor dem Start) zählen nicht.

**Schätzung (Planner, 2026-10-18; Metrik wie oben; Gegenprobe-Code mit dem Doppelten des Produktcodes für Neues, Bestand mit gezählten Zeilen):**

| Teil | Zeilen |
|---|---|
| zwölf Fragmente (eine Zeile ersetzt, Kommentarköpfe nachgezogen, je 4 bis 8 Zeilen in und aus) | 50 bis 90 |
| `gegenproben.mk` und `gegenproben.sh` (zwei Listen, zwei Ziele, Skript mit `MAKEFLAGS`, Schleife, Sammlung) | 70 bis 110 |
| Nachzug der sieben Gegenproben zum Anschluss (heute 70 Zeilen betroffen; Fälle und Mutationen bleiben, das Ziel wechselt) | 100 bis 200 |
| `gates-schnell-gegenprobe` (Ausnahmeliste, Deckung beider Richtungen, Fixture, 11 Mutationen der Ausnahmeliste entfallen, Mutationen je Element der neuen Regel) | 60 bis 120 |
| `gegenproben-gegenprobe` (rund 14 Randformen aus Gruppe G, das Doppelte des Produktcodes und mehr, weil die Attrappen Zeilen kosten) | 200 bis 320 |
| `harness/README.md` (Zeilen, die sich ändern, zählen je einmal ein und aus) und `AGENTS.md` | 40 bis 70 |
| **Zusammen** | **520 bis 910** |

Meldeschwelle 700, **Marke 900**. Die Mitte (715) liegt knapp über der Meldeschwelle, die Obergrenze (910) über der Marke: Der Slice ist am Rand des Schnitts. Der Planner hat deshalb bereits geschnitten (Rollen-Anweisungen, Doku-Abschnitt und Sendung an den Release in `slice-harness-gegenproben-aus-gates-fahrstelle`); ein weiterer Schnitt unterbleibt. Die Spanne hängt am Nachzug der Gegenproben und an der Zahl der Randformen der Gegenprobe; zählt der Architect mehr als 14 Randformen in Gruppe G oder mehr als zwölf Punkte im Nachzug der Verträge, rechnet er neu.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder der Vertrag verlangt mehr als drei Liefer-Punkte oder mehr als zwei Schichten, oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner (Vorschlag: Kettenumbau mit Nachzug der Gegenproben getrennt von den Sammelzielen mit ihrer Gegenprobe).
- `in-progress` → `open` (blockiert): Eine Gegenprobe wird außerhalb von `make gates` nur grün, wenn ein anderes Gate vorher lief (versteckte Abhängigkeit zwischen Gates; Befund an den Planner).

## 5. Closure-Trigger

DoD vollständig; `make gates` grün auf dem Endstand; `make gegenproben` grün im Gegenproben-Container oder, wenn er nicht vorliegt, am Host (Beleg im Bericht mit Dauer); `make gegenproben-host` grün am Host; Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag geschrieben; die drei Paarungen getragen.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des neuen Vertrags (`AGENTS.md` §3.12).** Je Randform steht die Empfehlung des Planners; entschieden ist keine, der Architect legt sie vor dem Code in den Vertrag (`harness/vertraege/`), und der Implementer entscheidet keine still (Randform-Rückgabe). Was aus dem Bestand entscheidbar war, steht als „Bestand“.

**Gruppe K — Kette, Namens- und Zuordnungsregel**

1. **Namensregel.** Ein Gate, dessen Name auf `-gegenprobe` endet, und `a-check-negativ` hängen an genau einer Liste, nie an `GATE_CHECKS`. Ein Gate dieses Namens in `GATE_CHECKS` macht die Gegenprobe der Kette rot (Ersatz von `SPEC-052` Punkt 4). Ein Gate, das beiden Listen angehört, ebenso. Ein Gate mit Namen auf `-gegenprobe`, das an keiner Liste hängt, ebenso (es fiele still aus). Eine Gegenprobe mit anderem Namen (`-negativ`, `-mutation`) erkennt die Regel nicht: akzeptiertes Negativ, Grund: eine Regel am Namen ist so weit wie der Name. Entschieden: offen.
2. **Zuordnung zu Container oder Host (Regel an den Schreiber).** Eine neue Gegenprobe kommt in die Host-Liste, wenn sie `u-boot` oder `devcontainer` aufruft, sonst in die Container-Liste. Vorschlag: die Zuordnung ist Entscheidung des Schreibers im Fragment (eine Zeile); die Gegenprobe der Kette prüft, dass jedes Ziel in genau einer Liste steht, und kann die Wahl nicht prüfen, weil sie nicht weiß, was ein Skript aufruft (akzeptiertes Negativ). Die Wahl wird im Gegenproben-Container geprüft: Eine Host-Gegenprobe in der Container-Liste endet dort rot wegen des fehlenden Werkzeugs, eine Container-Gegenprobe in der Host-Liste läuft am Host unbemerkt weiter (der falsche Fehler ist der erste, nicht der zweite). Vorschlag zusätzlich: eine Kopfzeile in jeder Gegenprobe, die ihre Liste nennt. Entschieden: offen.
3. **Ausnahmeliste der schnellen Kette** — drei Einträge; ein Eintrag einer Gegenprobe darin ist ein Fehler der Deckung (Ersatz von `SPEC-052` Punkt 4). Entschieden: offen (Bestand: Punkt 3 sagt „ein Fehler der Liste kostet Zeit, nie ein Gate“; hier kostet er nichts mehr, die Aussage entfällt).
4. **Nachweis des vollen Laufs und Hook.** `make gates` schreibt Dateiliste und Hash wie bisher (Bestand, `SPEC-052` Punkte 9 bis 11, 21); der Satz „`make gates` ist der einzige Nachweis für Code, Gates, Profile, Skripte …“ (Punkt 9) wird enger: Der Nachweis schließt die Gegenproben nicht mehr ein. Der Hook-Vertrag Teil 2 (Punkte 9 bis 22) und der Code des Hooks bleiben, Entscheidung: Der Hook verlangt keinen Stand der Gegenproben. Entschieden: offen (Frage 4).
5. **Beleg für die Positivliste (Punkt 13).** Er stützt sich auf die Gates der Ausnahmeliste, die `docs/` nicht lesen; die Gegenproben stehen nach dem Umbau außerhalb, der Beleg gilt für sie weiter, wenn eine Änderung unter `docs/plan/planning/` nie eine Gegenprobe rot macht. Entscheidung: Punkt 13 nennt die Gegenproben als Leser, die nie aus dem Arbeitsbaum lesen, und die Pflicht gilt auch für sie. Entschieden: offen.
6. **Aufrufarten von `make` für `make gates`.** `-n`, `-q`, `-t`, `-i` und `MAKEFLAGS` wirken auf Stand und Nachweis wie in `SPEC-052` Punkt 6 und 21; die Kette ist kleiner, die Zusagen bleiben (Bestand). Fälle in `gates-schnell-gegenprobe` bleiben, die Attrappen fahren die verkleinerte Kette. Entschieden: Bestand.
7. **`make gates-schnell` und `make gates`.** Beide Ketten unterscheiden sich nur durch `test-container`, `test-integration`, `test-postgres-versionen`. Ob `gates-schnell` überflüssig wird, ist nicht Gegenstand (Frage 2); der Vertrag lässt beide stehen. Entschieden: offen.

**Gruppe G — Sammelziele**

1. **Zwei Ziele, zwei Listen, ein Skript.** Name der Ziele und der Listen: Vorschlag des Architect. Beide Ziele sind kein Gate, hängen nicht an `GATE_CHECKS`, hängen nicht aneinander; das Container-Ziel fährt die Host-Liste nicht. Entschieden: offen.
2. **Reihenfolge.** Nacheinander in der Reihenfolge der Liste (Reihenfolge der Includes, alphabetisch nach Fragmentdatei), nie parallel, auch nicht unter `-j` (eine Schleife im Skript statt eines Fächers von Prerequisites): Die Last (Docker, Load 57 auf 20 Kernen, gemessen in [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md)) ist der Grund für das Herausnehmen. Entschieden: offen.
3. **Rotes Einzelziel.** Alle Ziele der Liste fahren, auch nach einem roten (wie `-k`); die Zusammenfassung nennt die roten beim Namen, Exit 1. Ein Abbruch durch `SIGINT` oder `SIGTERM` beendet das laufende Ziel und startet kein weiteres (Exit 130 und 143). Je ein Fall. Entschieden: offen.
4. **Aufrufarten.** `-n`, `-q`, `-t`: Das Rezept trägt weder `$(MAKE)` noch `+`, make führt es nicht aus (wie `record-gates`); `-i`: Die Einzelziele erben `-i` und gäben ein rotes Skript als Exit 0 zurück, deshalb lehnt das Skript `-i` in `MAKEFLAGS` ab (Exit 2, Meldung), statt grün zu melden; `-k`, `-j` ohne Wirkung. `-i` zusammen mit `-n`, `-q` oder `-t`: `-n`, `-q`, `-t` gelten, nichts läuft. Kein Stand wird geschrieben. Je Flag ein Fall, auch aus der Umgebung und kombiniert. Entschieden: offen.
5. **Fehlendes Werkzeug.** Fehlt `docker`, `bash` oder `setsid`, ist das betroffene Einzelziel rot mit seiner eigenen Meldung, das Sammelziel fährt weiter; kein Vorab-Test der Werkzeuge. Entschieden: offen.
6. **Ausgabe.** Die Ausgabe der Einzelziele erscheint unverändert; die Zusammenfassung steht zuletzt und nennt je Ziel Name und Ausgang. Ein Überspringen innerhalb eines Ziels ([ADR-0046](../../adr/0046-gegenprobe-ueberspringt-im-dev-container.md), ersetzt durch [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md)) gibt es nach der Aufteilung nicht mehr. Entschieden: offen.
7. **Leere Liste.** Die Host-Liste ist ohne `container-start-gegenprobe` leer (§4 *Start*); eine leere Liste ist ein Fehler (Exit 2, Meldung), kein stilles Grün. Entschieden: offen.
8. **Eine Gegenprobe steht in beiden Zielen oder in keinem.** Siehe K1; das Sammelskript prüft die Zuordnung nicht, die Gegenprobe der Kette tut es.
9. **Aufruf.** Aus einem Unterverzeichnis des Repos wie von der Wurzel; ohne Git-Arbeitsbaum Exit 2. Entschieden: offen.
10. **Gleichzeitige Läufe und Worktrees.** Zwei Fahrten im selben Arbeitsbaum oder zusammen mit `make gates` sind nicht zugesagt (`SPEC-052` Punkt 8, Aufruf-Voraussetzung nacheinander); zwei Worktrees am selben Docker-Daemon teilen Namen von Containern und Netzen (Präfixe `pgr-ivg-`, `pgr-probe-`), die Gegenprobe der Container schützt nur durch ihre Präfixe. Der Vertrag bindet die Zusage an „nacheinander“; die Fahrstelle im eigenen Container (`slice-harness-gegenproben-aus-gates-fahrstelle`) hat einen eigenen Daemon. **Aufruf-Voraussetzung (§3.12):** Das Sammelskript sagt keine Ordnung über mehrere Aufrufe zu. Entschieden: offen.
11. **Dev-Container-Erkennung entfällt.** Keine Erkennung, kein Überspringen; die Zuordnung (K2) ersetzt sie. Der Vertrag der Container-Start-Gegenprobe (`SPEC-058`) hat den Punkt zur Erkennung nach [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) verloren (umgestellt im Zweig des Geruest-Slice). Entschieden: Nutzer, 2026-10-18; [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) Accepted.

**Aufrufarten des Treibers (§3.12):** Siehe G4 und K6; das Sammelskript liefert keinen Stand und keinen Nachweis, den ein Hook glaubt.

**Fahrstelle, Stand und Ressourcen (für `slice-harness-gegenproben-aus-gates-fahrstelle`, hier genannt, damit der Vertrag sie bindet).** Der Gegenproben-Container sieht nur den eingelieferten, committeten Stand eines Zweigs, nicht den Arbeitsbaum des Hosts; was nicht eingeliefert ist, ist nicht gefahren. Zwei Container desselben Klons sind nicht zugesagt. Ressourcen laut `docs/developer/devcontainer.md` §10: je Container 8 CPUs und 12 GiB (Vorgabe), Host 20 CPUs und 31 GiB; Arbeitscontainer und Gegenproben-Container zusammen 16 CPUs und 24 GiB. Der Vertrag dieses Slice sagt dazu nichts zu; die Zusage steht in der Dokumentation, soweit gemessen.

### Risiken

- **Ohne Kette bleibt eine Regression eines Gates unbemerkt, bis jemand die Gegenprobe fährt** (`make gates` grün sagt nichts mehr über die Gegenproben; der Stop-Hook gibt auf `make gates` frei; ein Gate, das niemand prüft, ist eine Behauptung, `modul-13-quality-gates.md`). — **Ausgang (vorgesehen):** weiter offen → neues `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt` im Register. Mitigation: die Pflicht (Änderung am Werkzeug, an der Gegenprobe oder am Vertrag; Ende jeder Welle; vor dem Release), Review und Verifikation fahren die betroffene Gegenprobe gezielt (`slice-harness-gegenproben-aus-gates-fahrstelle`). Tritt eine unbemerkte Regression auf, ist der Ausgang *eingetreten* → Slice für einen Stand der Gegenproben (Frage 4) oder für den verworfenen Zuschnitt „nur bei geänderten Eingaben“.
- **Die Pflicht trägt nur, wenn jemand sie liest, und ein Sensor für „geprüftes Werkzeug geändert“ fehlt** (eine Zuordnung Werkzeug zu Gegenprobe ist die Eingabeliste, die der erste Zuschnitt brauchte). — **Ausgang (vorgesehen):** weiter offen → dasselbe `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`.
- **Eine Gegenprobe ist ohne Kette nur grün, wenn sie von selbst hält:** Sie lief im Verbund mit `make gates` und kann von einem Zustand leben, den ein anderes Gate herstellt (Image, Cache, Datei im Arbeitsbaum). — **Ausgang (vorgesehen):** entfallen, wenn `make gegenproben` im frischen Gegenproben-Container grün endet, ohne dass vorher `make gates` lief (Beleg im Bericht); sonst eingetreten → Rückführung nach `open` (§4).
- **Die Host-Liste hat keine Gegenprobe, weil `slice-harness-devcontainer-einrichten-geruest` nicht in `done/` liegt oder `container-start-gegenprobe` doch an `GATE_CHECKS` oder in der Ausnahmeliste hängt** ([ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) sagt das Gegenteil; der Zweig trägt es, main noch nicht). — **Ausgang (vorgesehen):** entfallen, wenn beim Start `grep container-start-gegenprobe harness/mk/*.mk` sie weder an `GATE_CHECKS` noch in der Ausnahmeliste zeigt und der Slice in `done/` liegt (Beleg im Bericht); sonst eingetreten → Rückführung nach `open` (§4).
- **Die Größe liegt am Rand** (Schätzung 520 bis 910, Marke 900). — **Ausgang (vorgesehen):** entfallen, wenn die Messung nach der Metrik unter 900 bleibt; darüber meldet der Implementer (§4 Rückführung *zu groß*).
- **Der Hook-Vertrag bleibt formal intakt, sein Nachweis verliert Gehalt** (siehe K4 und K5). — **Ausgang (vorgesehen):** entfallen, wenn Punkt 9 und 13 von `SPEC-052` die engere Aussage tragen (Review prüft den Wortlaut) und Frage 4 vom Nutzer entschieden ist; sonst weiter offen → dasselbe Register-Objekt.

**Fragen an den Nutzer** (Empfehlungen des Planners, nicht entschieden):

1. **Wann?** Entschieden: Der Slice startet unabhängig von `welle-harness-regeln` (Auftrag „sofort“, Lösung aus der Reihenfolge vom Koordinator am 2026-10-18 angewiesen; §4 *Start*). Keine offene Frage mehr.
2. **Bleibt `make gates-schnell`** neben dem verkleinerten `make gates`? Nur genannt, nicht entschieden; offene Frage für den Nutzer nach der Messung der Dauer von `make gates` (DoD).
3. **Wortlaut der Pflicht** in `AGENTS.md` §6 Schritt 6: „Wer ein Gate, eine Gegenprobe oder einen Vertrag ändert, fährt die betroffene Gegenprobe (`make <name>`, Zuordnung in `harness/README.md` §Sensors); am Ende jeder Welle und vor dem Release fahren `make gegenproben` im Gegenproben-Container und `make gegenproben-host` am Host.“ Empfehlung: so, als Satz in Schritt 6, nicht als neue Regel in §3. Die Änderung an `AGENTS.md` braucht das ausdrückliche Ja des Nutzers beim Start des Slice (§4 *Start* (d)).
4. **Soll ein Stand der Gegenproben** (grün je Fahrt, Hash des Arbeitsbaums) entstehen und der Hook ihn für nicht inerte Änderungen an Gates und Werkzeugen verlangen? Das brächte die Laufzeit zurück; Empfehlung des Planners: nein (kein Stand der Gegenproben, kein Hook), erst wenn die Pflicht nicht trägt (Risiko 1).

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe/evidence/slice-harness-gegenproben-aus-gates.md` mit dem Plan angelegt (Planner, 2026-10-18); die Closure ergänzt die Messung; `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt/` entsteht bei der Closure (Risiko 1)
- **Folge-Slices:** `slice-harness-gegenproben-aus-gates-fahrstelle` (Rollen-Anweisungen, Gegenproben-Container, Release-Pflicht — ist eine Datei in `open/`)
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Hauptgegenstand; Verhältnis: der Slice beantwortet die Laufzeit durch Herausnehmen, nicht durch Kürzen; Beleg liegt an), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (DoD 1 und 2 tragen je Element einen Fall und eine Mutation), `BEO-REPO/spec-randform-erst-im-review-entschieden` (§6 nennt die Randformen vor dem Code), `BEO-REPO/slice-waechst-durch-uebernahmen` (der Slice übernimmt nichts; die Schnitte stehen in §1), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groessenprobe-ohne-messbaren-stand` (§4 rechnet in der Metrik mit dem Doppelten für Neues), `BEO-REPO/gate-kopie-bricht-bei-laufenden-agenten` (Nachbar: `slice-harness-gegenprobe-baumkopie`, der dahinter steht). Keiner der Treffer erreicht mit diesem Slice 3× neu, `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` ausgenommen, der bereits `verkörpert` ist und einen Beleg erhält.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag und die ADR des Architect stehen vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Liefer-Punkte 2 (Kette mit Namens- und Zuordnungsregel und Nachzug der Gegenproben | Sammelziele mit Gegenprobe), Platz für einen dritten bleibt. Schichten 1: Harness-Werkzeug mit Gegenprobe und Gate-Index; `AGENTS.md` §6 ist ein Satz im selben Harness-Artefakt, Planung, Vertrag und ADR sind keine Schicht; Rollen-Anweisungen und Wartungs-Doku liegen in `slice-harness-gegenproben-aus-gates-fahrstelle`. Schätzung in der Metrik von §4: 520 bis 910 (Marke 900, Meldeschwelle 700). Der Plan trägt zwei Sendungen ein: an `slice-harness-gegenproben-aus-gates-fahrstelle` (Folge-Slice, neu, nimmt an: dessen §1) und an `slice-harness-gegenprobe-baumkopie` und `slice-harness-lint-warten-ohne-frist` (Nachzug der DoD-Zeile `make gates`, kein Liefer-Punkt, keine Schicht; Nachzählung dort in §8, im selben Commit).
