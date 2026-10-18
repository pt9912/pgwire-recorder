# Slice slice-harness-gegenproben-aus-gates-fahrstelle: Hook für die Einlieferung in den Gegenproben-Container: `tools/harness/gegenproben-nach-einlieferung.sh` mit Gegenprobe

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten, dritte Stelle der Reihenfolge, direkt hinter `slice-harness-gegenproben-aus-gates-sammelziele` und vor `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (Entscheidung des Nutzers vom 2026-10-18; wie jener Slice aus der Reihenfolge der Wellen gelöst: Start-Trigger ist `slice-harness-gegenproben-aus-gates-sammelziele` in `done/` (das seinerseits `slice-harness-gegenproben-aus-gates` voraussetzt), nicht `welle-harness-regeln`). Angelegt am 2026-10-18 als Folge-Slice von `slice-harness-gegenproben-aus-gates`. **Geschnitten am 2026-10-18** (Entscheidung des Nutzers: Schnitt ja; Regel des Nutzers: schneiden, wenn die Schätzung die Marke selbst überschreitet): Der Architect schätzte den ungeteilten Slice nach `SPEC-061` auf 740 bis 1345 Zeilen (Marke 900). Dieser Slice behält den Hook mit seiner Gegenprobe (Kennung bleibt); die Pflicht in den Rollen-Anweisungen, der Abschnitt im Dokument über den Gegenproben-Container und die erste Fahrt sind in `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` aufgegangen. Der Slice liegt in `open/`: Er ist nicht priorisiert, den Anspruch trägt der Nutzer.

**Bezug:** [LH-QA-07](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur insoweit, als die Gegenproben dessen Messungen tragen; kein Anspruch des Produkts, der Slice liefert ein Harness-Werkzeug. Bindung an das Register: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` und das mit `slice-harness-gegenproben-aus-gates` entstehende `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`. Keine aktive ADR berührt (der Hook setzt [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) um, eine eigene ADR entsteht nicht, Entscheidung des Architect).

**Berührte Spec-Stellen:** der Vertrag des Hooks für die Einlieferung [`SPEC-061`](../../../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md) Punkte 1 bis 14 und 16 (vom Architect geschrieben, vor dem Code); der Vertrag der Sammelziele [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkt 2 (Namensregel, Liste `GEGENPROBEN_CONTAINER`) und Punkt 10 b (`MAKEFLAGS`), nur gelesen. Punkt 15 von `SPEC-061` (Lesen des Ergebnisses) ist Regel an den Leser und gehört zu `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`.

**Verantwortlich:** pt9912

**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Skript im Repo, `tools/harness/gegenproben-nach-einlieferung.sh`, ist der Inhalt des `post-receive`-Hooks im Bare-Repo des Gegenproben-Containers (Name vorläufig `pgr-gp`): Jede Einlieferung eines Zweigs startet dort `make gegenproben` (nur die Container-Liste) im Hintergrund auf genau dem eingelieferten Hash und legt bei beendetem Lauf eine Ergebnisdatei ab; der Hook endet sofort mit Exit 0. Eine Gegenprobe des Skripts (`gegenproben-nach-einlieferung-gegenprobe`, Fragment `harness/mk/gegenproben-nach-einlieferung-gegenprobe.mk`, Eintrag in `GEGENPROBEN_CONTAINER` nach `SPEC-060` Punkt 2, nicht in `GATE_CHECKS`) arbeitet in Temp-Repos mit einem echten Bare-Repo, einer Kopie des Skripts als Hook und einer Attrappe als `make`, ohne Docker; sie trägt je Zusage und je Element einer Aufzählung aus `SPEC-061` einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10). `harness/README.md` führt die Zeile der Gegenprobe (Sensors) und die Zeile des Skripts (Werkzeuge).

**`pgr-gp` ist ein spezieller Container, der nur die Gegenproben fährt** (Ergänzung des Nutzers vom 2026-10-18); in ihm wird kein Slice abgearbeitet. Das Skript läuft nur dort: Es tut nur etwas, wenn die Marke des Containers im Bare-Repo vorliegt (`SPEC-061` Punkt 2).

**Sendungen an `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`** (eingetragen vom Planner im Commit, der den Schnitt vornimmt, `AGENTS.md` §3.13): (1) *Bindung:* Der Abschnitt im Dokument, die Zeile zum Aufsetzen und die Regel an den Leser ruhen auf den Zusagen dieses Slice: Name und Ort des Skripts, Kopie mit `install -m 755`, Marke `<Bare-Repo>/gegenproben-ablage`, Pfade `gegenproben-<Name>-<Hash>` in der Ablage (`SPEC-061` Punkte 1, 2, 8, 9 und 15). Ändert dieser Slice eine davon, zieht er sie im Nehmer nach oder meldet es dem Planner. (2) *Die Einlieferung, die den Hook zum ersten Mal an echtem `make gegenproben` zeigt*, ist die erste Fahrt dort; sie und die Messung der Dauer (`SPEC-061` Punkt 16 (1)) gehören dem Nehmer. Beide Sendungen stehen in dessen §1, DoD und §6 mit der Kennung dieses Slice.

**Sendungen, die der Schnitt vom Ursprung dieses Slice in `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` gelegt hat** (Nachweis der Adresse: `grep` nach der Kennung dieses Slice im Nehmer): die Sendung von `slice-harness-gegenproben-aus-gates-sammelziele` (erste Fahrt im frischen Container, Risiko 1 dort), die Sendung an `slice-erster-release-freigabe` (Punkt der Freigabe-Checkliste), die Sendung an die Welle-Closure und die Bindung an `docs/developer/devcontainer.md` §6 und §9 aus `slice-harness-devcontainer-einrichten` und `slice-harness-devcontainer-einrichten-start`. Dieser Slice trägt keine davon; die Nehmer nennen den Geber `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (Nachzählung dort jeweils in §8, im Commit des Schnitts).

**Übernimmt:** nichts.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Pflicht in den Rollen-Anweisungen (`implement-slice`, `close-welle`, Typen von Implementer, Reviewer, Verifier, Reviewer-Skill, `AGENTS.md` §6), der Abschnitt über den Gegenproben-Container in `docs/developer/devcontainer.md` und die erste Fahrt — `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`: Anweisung und Beschreibung nennen das Skript, das hier entsteht; die Fahrt braucht den Hook. Der Schnitt trennt das Werkzeug (Schicht Harness-Werkzeug) von Betrieb und Beschreibung, weil die Schätzung des ungeteilten Slice über der Marke lag.
- Der Kettenumbau, die Listen und die Namensregel, die erste Hälfte des Satzes in `AGENTS.md` §6 Schritt 6 — `slice-harness-gegenproben-aus-gates`; die Sammelziele, ihre Namen und Gegenprobe und die zweite Hälfte des Satzes — `slice-harness-gegenproben-aus-gates-sammelziele`; dieser Slice setzt beider `done/` voraus (die Ziele müssen existieren, bevor das Skript eines ruft).
- Ein Scheduler, eine CI-Pipeline, ein Warteschlangen-Dienst oder eine Benachrichtigung — Entscheidung des Nutzers: Der Container wird von Hand oder von einem Agenten gestartet; der `post-receive`-Hook im Bare-Repo von `pgr-gp` ist die kleinste Form, die Einlieferung und Fahrt verbindet, und nicht mehr. Ein Zeitlimit des Laufs ebenso: Empfehlung des Architect, nach Empfehlung entschieden (nicht ausdrücklich vom Nutzer bestätigt), `SPEC-061` Punkt 16 (1), bis die erste Fahrt die Dauer misst.
- Der Hook in den Bare-Repos der Arbeits-Container (`pgr-dev`, `pgr-dev2`, …) — Schicht-Abgrenzung und Last: Der Arbeits-Container soll nicht ausgebremst werden; das Skript lehnt den Lauf ab, wo die Marke fehlt (`SPEC-061` Punkt 2). Dass jemand die ganze Zeile zum Aufsetzen dort ausführt, ist ein akzeptiertes Negativ (`SPEC-061`).
- Ein Ziel, das den Gegenproben-Container anlegt, einliefert und startet (`make container-einrichten` und ein Ziel für die Fahrt) — `slice-harness-devcontainer-einrichten` automatisiert Bare-Repo, Einlieferung, Klon und Zweig; ein Ziel für das Aufsetzen des Hooks wäre ein anderer Vorgang (Werkzeug am Container statt Hook); eine Adresse gibt es nicht, ein Slice legt der Planner an, wenn der Nutzer es entscheidet.
- `u-boot` und `devcontainer` im Gegenproben-Container — nicht gebraucht: Die Gegenproben, die sie brauchen, laufen am Host (`make gegenproben-host`, Entscheidung des Nutzers vom 2026-10-18).
- Ein Stand oder Nachweis der Gegenproben, den der Stop-Hook am Host liest — `slice-harness-gegenproben-aus-gates` §1 und §6 (Frage 4 dort, Empfehlung des Planners: nein); `SPEC-060` Punkt 14. Das Ergebnis des Hooks liegt im Container, nicht im Arbeitsbaum des Hosts.
- Dass die Kopie im Bare-Repo dem Skript im Repo entspricht — akzeptiertes Negativ (`SPEC-061` Punkt 1; Grund: die Kopie liegt außerhalb des Repos).
- Produkt-Code, Handbuch, `README.md` — Schicht-Abgrenzung: kein Verhalten des gebauten Binary.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Hook-Skript):** `tools/harness/gegenproben-nach-einlieferung.sh` erfüllt `SPEC-061` Punkte 1 bis 14 und 16: eine einzelne Datei, die als Kopie `post-receive` sein kann; Marke und Voraussetzungen in der Reihenfolge der Gründe; je eingeliefertem `refs/heads/*`-Zweig ein Lauf im Hintergrund auf genau dem Hash, nacheinander über die Sperre, im eigenen Arbeitsbaum, mit `make gegenproben` ohne `MAKEFLAGS` und `GIT_*`; Ergebnisdatei `gegenproben-<Name>-<Hash>` nur bei beendetem Lauf (Felder und Reihenfolge nach Punkt 9); Exit 0 in jedem Fall. Jede Randform aus §6 ist vor dem ersten Code-Commit in `SPEC-061` entschieden (Architect, `AGENTS.md` §3.12).
- [ ] **Liefer-Punkt 2 (Gegenprobe):** `tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh` mit dem Fragment `harness/mk/gegenproben-nach-einlieferung-gegenprobe.mk` (Eintrag in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS`, nicht in der Ausnahmeliste der schnellen Kette): je Zusage und je Element einer Aufzählung aus `SPEC-061` ein Fall und eine Mutation, die nur ihn (oder die genannten Fälle) rot färbt; die Gegenprobe verlangt genau diese Menge roter Fälle. Der Bericht nennt je Zusage: Zusage · Mutation · roter Test (`AGENTS.md` §3.10), vom Implementer selbst rot gesehen. Gegenprobe grün über `make gegenproben-nach-einlieferung-gegenprobe`; Docker braucht sie nicht.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `tools/` und `harness/mk/` ändern); die Gegenprobe des Skripts ist darin nicht enthalten (Namensregel), sie fährt einzeln.
- [ ] Doku-Update: `harness/README.md` §Sensors führt die Zeile der Gegenprobe (Vertrag `SPEC-061`, Ziel steht in `GEGENPROBEN_CONTAINER`, bash, `git`, `flock`, `setsid`, ohne Docker), §Werkzeuge die Zeile des Skripts (kein Gate, Vertrag `SPEC-061`). Aussagen der Zeilen nennen nur, was die Gegenprobe prüft (`AGENTS.md` §3.11).
- [ ] Beobachtungs-Register: Ausgänge *weiter offen* aus §6 geführt (`BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt` mit dem Beleg dieses Slice, Mitigation getragen: ja oder nein; `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, wenn die Größenmarke eintrat); kein Zähler wird gesetzt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review. Das Review fährt die Gegenprobe des Skripts gezielt selbst.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen; die Folge-Slice-Paarung prüft Existenz **und** Annahme: `grep` nach `slice-harness-gegenproben-aus-gates-fahrstelle` in `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` findet die Sendungen in §1 und §6 (`AGENTS.md` §3.13).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md` | neu (Architect, vor dem Code, geschrieben) | Vertrag des Hook-Skripts: Aufruf, Ergebnisdatei, Entscheidung je Randform |
| `tools/harness/gegenproben-nach-einlieferung.sh` | neu | Inhalt des `post-receive`-Hooks: liest die Zeilen von `git` auf stdin, startet je Zweig einen Lauf im Hintergrund, legt das Ergebnis ab |
| `tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh`, `harness/mk/gegenproben-nach-einlieferung-gegenprobe.mk` | neu | Gegenprobe in Temp-Repos mit Attrappe als `make` (kein Docker): je Zusage und je Element einer Aufzählung ein Fall und eine Mutation |
| `harness/README.md` | update | Zeile der Gegenprobe in §Sensors, Zeile des Skripts in §Werkzeuge |
| `docs/plan/planning/open/slice-harness-gegenproben-aus-gates-fahrstelle-betrieb.md` | neu (Planner, im Commit des Schnitts) | Nehmer der Sendungen aus §1; Zählung dort in §8 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. **Vorbedingung:** `slice-harness-gegenproben-aus-gates` und `slice-harness-gegenproben-aus-gates-sammelziele` liegen in `done/` (die Listen, die Ziele und ihre Namen stehen); der Vertrag des Hook-Skripts liegt in `harness/vertraege/` (`SPEC-061`, Commit des Architect `23db486c`), mit der Entscheidung je Randform aus Gruppe H (`AGENTS.md` §3.12); eine ADR entsteht nicht (Entscheidung des Architect). Reihenfolge: erste Stelle nach `slice-harness-gegenproben-aus-gates-sammelziele`; der Slice wartet nicht auf `welle-harness-regeln` (Drift-Log der Roadmap). `Verantwortlich:` trägt beim Start Person und Zweig; er läuft auf eigenem Zweig.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice>` über die Pfade am Ende dieses Absatzes, Tests eingeschlossen. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt. **Pfade:** `-- tools harness/mk harness/README.md`; `harness/vertraege` liegt nicht darin (vor der Basis geschrieben; eine Rückgabe an den Architect, die den Vertrag ändert, zählt nicht mit).

**Schätzung (Architect, 2026-10-18, nach `SPEC-061`, auf diesen Slice verteilt vom Planner; Metrik wie oben):**

| Teil | Zeilen |
|---|---|
| Hook-Skript (Marke, Werkzeuge, Zeilen, Kodierung, Sperre, Arbeitsbaum, Ergebnis, `SIGTERM`) | 140 bis 200 |
| Gegenprobe des Skripts: rund 75 Fälle und 100 Mutationen nach `SPEC-061` (bei `SPEC-060` stehen 53 Zeilen Produkt gegen 381 Zeilen Gegenprobe, Faktor 7; hier angesetzt 3 bis 4, weil Fälle ein Bare-Repo und eine Attrappe teilen) | 450 bis 800 |
| Fragment der Gegenprobe | 10 bis 15 |
| `harness/README.md` (zwei Zeilen) | 10 bis 20 |
| **Zusammen** | **610 bis 1035** |

Der Architect nannte für den Schnitt-Teil „Hook-Skript mit Gegenprobe“ 630 bis 1015; die Differenz von höchstens 20 Zeilen kommt aus der Zuordnung der Zeilen in `harness/README.md` (hier ganz diesem Slice). Meldeschwelle 700, **Marke 900**. Die Mitte (822) liegt bei 91 % der Marke, die Untergrenze darunter, die Obergrenze (1035) darüber. Der Schnitt, den der Nutzer entschieden hat, bringt die Mitte unter die Marke; die Obergrenze bleibt, weil Skript und Gegenprobe keine tragfähige weitere Naht haben (eine Gegenprobe, die nur einen Teil der Punkte von `SPEC-061` fängt, ließe Zusagen ohne Mutation, `AGENTS.md` §3.10). **Auflage zur Größe (Architect):** Die Fälle je Punkt-Gruppe zu einem Aufruf bündeln (ein Aufbau von Bare-Repo und Attrappe je Gruppe, nicht je Fall); der Bericht nennt die Messung in der Metrik. **Entscheidung nach Empfehlung (2026-10-18; Empfehlung des Planners, vom Nutzer nicht ausdrücklich bestätigt):** keine weitere Teilung. Marke 900, Obergrenze 1035; es gibt keine tragfähige Naht (Grund oben: Zusagen ohne Mutation), und der Vorbehalt der Rückführung bei einer Summe über 900 steht unten.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900. Der Planner entscheidet dann die Naht (nicht vorab geschnitten: Skript und Gegenprobe gehören zusammen); der Bericht nennt die Messung und die Punkte von `SPEC-061`, die den größten Teil tragen.
- `in-progress` → `open` (blockiert): Der Implementer findet eine Randform, die `SPEC-061` nicht entscheidet oder die dem Vertrag widerspricht (Rückgabe an den Architect, `.claude/commands/implement-slice.md`, Randform-Rückgabe); oder ein Werkzeug, das der Vertrag voraussetzt (`flock`, `setsid`, `pgrep`, `timeout`), fehlt im Image der Gegenprobe.

## 5. Closure-Trigger

DoD vollständig; `make gegenproben-nach-einlieferung-gegenprobe` grün, jede Mutation einzeln rot gesehen und im Bericht; `make gates` grün; Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag; die drei Paarungen getragen, darunter die Annahme der Sendungen aus §1 in `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (`grep` nach der Kennung dieses Slice dort findet §1 und §6).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Der Slice liefert einen neuen Vertrag im Sinne von `AGENTS.md` §3.10 (ein Skript mit Gegenprobe). Die Randformen sind in [`SPEC-061`](../../../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md) entschieden, vor dem ersten Code-Commit (`AGENTS.md` §3.12); die Tabelle *Randform → Punkt* im Kopf des Vertrags ist die Zuordnung. Entscheidungen, die kein Nutzer getroffen hat, stehen als solche.

**Gruppe H — Hook für die Einlieferung** (Randformen des neuen Vertrags; Entscheidungen des Architect in `SPEC-061`):

1. **Der Push wartet nicht.** `post-receive` läuft im Push; der Hook startet den Lauf im Hintergrund (`setsid`, Ausgabe in eine Logdatei) und endet sofort mit Exit 0; ein Fehler des Hooks darf den Push nicht blockieren. Entschieden: Architect, `SPEC-061` Punkt 5 (Hook endet mit Exit 0; der Lauf hat weder stdout noch stderr noch stdin des Hooks und eine eigene Sitzung).
2. **Welche Refs.** Nur `refs/heads/*`; Tags und andere Refs lösen keinen Lauf aus. Mehrere Zweige in einem Push: je Zweig ein Lauf, nacheinander in der Reihenfolge der Zeilen. Entschieden: Architect, `SPEC-061` Punkt 3; die Reihenfolge nacheinander ist eine **Empfehlung des Architect, nach Empfehlung entschieden** (nicht ausdrücklich vom Nutzer bestätigt).
3. **Zwei Einlieferungen kurz hintereinander.** Ein Lauf zur Zeit, der zweite wartet über `flock`; alle Läufe laufen nacheinander. Entschieden: Architect, `SPEC-061` Punkt 6 (die Reihenfolge verschiedener Einlieferungen ist nicht zugesagt, akzeptiertes Negativ).
4. **Zweig gelöscht** (neuer Hash aus Nullen) — kein Lauf, kein Fehler; vorhandene Ergebnisdateien bleiben. Entschieden: Architect, `SPEC-061` Punkt 4 (dazu Zeilenformat und Form des Hashes, 40 oder 64 Stellen).
5. **Derselbe Hash erneut:** neuer Lauf; ein vorhandenes Ergebnis desselben Namens wird beim Aufruf entfernt und nach dem Lauf neu geschrieben. Entschieden: Architect, `SPEC-061` Punkt 7.
6. **Name des Zweigs im Dateinamen:** `%XX` für jedes Byte außer `A-Za-z0-9._-`; Name über 120 Zeichen: Meldung, kein Lauf. Entschieden: Architect, `SPEC-061` Punkt 8.
7. **Ergebnis nur bei beendetem Lauf:** über Temp-Datei und `mv`; Marke `<N>.laeuft` beim Start, bei `SIGTERM`, `SIGKILL` oder Neustart kein Ergebnis. Entschieden: Architect, `SPEC-061` Punkt 9 (`SIGINT` und `SIGHUP` nicht zugesagt).
8. **Nur die Container-Liste.** Der Hook ruft `make gegenproben`, nie `make gegenproben-host`. Entschieden: Nutzer, 2026-10-18; Umsetzung `SPEC-061` Punkt 10.
9. **Der Hook gehört nur dem Gegenproben-Container.** Marke `<Bare-Repo>/gegenproben-ablage`; ohne sie Meldung und Exit 0, kein Lauf. Entschieden: Architect, `SPEC-061` Punkt 2.
10. **Wohin gehört der Hook-Code?** Das Skript liegt im Repo; der Hook im Bare-Repo ist eine **Kopie** (`install -m 755`), kein Symlink. Die **Ablage liegt außerhalb des Klons**, dort, wo die Marke sie nennt (**Empfehlung des Architect, nach Empfehlung entschieden**, nicht ausdrücklich vom Nutzer bestätigt; Abweichung vom Vorschlag des Planners „Pfad im Klon“). Entschieden: Architect, `SPEC-061` Punkt 1. Eine Abweichung der Kopie vom Stand des Repos erkennt keine Gegenprobe (akzeptiertes Negativ).
11. **Berechtigung und Zustand.** Reihenfolge der Gründe: Marke, Werkzeug, Ablage, Zeile; im Hintergrund Sperre, Arbeitsbaum, Lauf; Exit 0, nie ein Ergebnis `gruen` ohne beendetes `make`. Entschieden: Architect, `SPEC-061` Punkt 13.
12. **Welcher Hash zählt der Closure?** Regel an den Leser (`SPEC-061` Punkt 15) und deshalb Gegenstand von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`; dieser Slice liefert nur die Dateinamen, auf denen die Regel ruht (Bindung, §1). Das Lesen per `ssh` mit dem Hash als Muster ist eine **Empfehlung des Architect, nach Empfehlung entschieden** (nicht ausdrücklich vom Nutzer bestätigt).

**Randformen, die der Architect zu den zwölf ergänzt hat** (`SPEC-061`, Tabelle im Kopf):

13. **Umgebung des Hooks (`GIT_*`, `MAKEFLAGS`, `MFLAGS`, `MAKELEVEL`, `MAKEOVERRIDES`):** der Lauf ruft `make` ohne sie. `SPEC-061` Punkt 11.
14. **Arbeitsbaum je Lauf:** eigener, losgelöster Checkout des Hashes in der Ablage, nach dem Ende entfernt. `SPEC-061` Punkt 12.
15. **Zeilenformat und Form des Hashes:** `SPEC-061` Punkt 4.
16. **Länge des Namens:** `SPEC-061` Punkt 8.
17. **Kein Zeitlimit des Laufs:** ein hängender Lauf hält die Sperre bis zu seinem Ende oder `SIGTERM`; akzeptiertes Negativ. **Empfehlung des Architect, nach Empfehlung entschieden** (nicht ausdrücklich vom Nutzer bestätigt), bis die erste Fahrt (in `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`) die Dauer misst. `SPEC-061` Punkt 16 (1).

**Aufrufarten des Treibers (§3.12):** Der Hook ruft `make gegenproben` ohne Flags und ohne `MAKEFLAGS` aus der Umgebung des Pushs; ein geerbtes `-i` ließe das Sammelskript den Lauf ablehnen (`SPEC-060` Punkt 10 b). Es gibt kein Rezept, das `make` aufruft; `-n`, `-q`, `-t` entfallen für den Hook (`SPEC-061` Punkte 11 und 14). **Aufruf-Voraussetzung (§3.12):** Zusage nur für Läufe nacheinander über die Sperre (H3); gleichzeitige Aufrufe des Hooks mit gleichem Namen haben vor der Sperre keine Zusage (akzeptiertes Negativ, `SPEC-061` Punkt 14).

**Mit dem Betrieb verwandte Randformen** (Eingeliefert gegen Arbeitsbaum, zwei Container desselben Klons, Ressourcen, Wegfall der Dev-Container-Erkennung, Fahrt ohne Netz, wer die Pflicht auslöst) gehören zu `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (dessen §6).

Risiken:

- **Der Hook liegt als Kopie außerhalb des Repos und veraltet** (H10): Wer das Skript ändert, setzt die Kopie nicht nach, und der Container fährt den alten Stand. — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`; die Zeile zum Aufsetzen und der Hinweis „nach jeder Änderung wiederholen“ stehen im Dokument des Nehmers (§1 dort).
- **Das Skript ist nur gegen eine Attrappe von `make` geprüft**; ein echter `make gegenproben` im Container zeigt sich erst in der ersten Fahrt des Nehmers. — **Ausgang (vorgesehen):** entfallen, mit Begründung: Die Gegenprobe fährt einen echten `git push` in ein Bare-Repo mit dem Skript als Hook (`SPEC-061` Punkt 1, `p1-echter-push`) und prüft Argumente und Umgebung des Aufrufs; was das echte `make` tut, ist `SPEC-060`, und die erste Fahrt steht in der DoD des Nehmers.
- **Ein hängender Lauf hält die Sperre ohne Ende, die folgenden Einlieferungen stauen sich** (kein Zeitlimit, Punkt 17). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`; die Dauer misst die erste Fahrt des Nehmers.
- **Die Obergrenze der Schätzung liegt über der Marke** (Marke 900, Obergrenze 1035; keine weitere Teilung, Entscheidung nach Empfehlung vom 2026-10-18, §4). — **Ausgang (vorgesehen):** *eingetreten* → Rückführung nach `next` (§4) und Naht des Planners, falls die Messung über 900 liegt; sonst *entfallen* mit der gemessenen Zahl in der Closure-Notiz.
- **Der Nehmer nimmt die Sendungen nicht an** oder zieht sie nicht nach, wenn dieser Slice den Vertrag ändert. — **Ausgang (vorgesehen):** entfallen mit `grep` nach der Kennung in §1 und §6 des Nehmers (Closure-Trigger).

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Hintergrund), `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt` (Gegenmittel), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (je Zusage eine Mutation, `AGENTS.md` §3.10), `BEO-REPO/spec-randform-erst-im-review-entschieden` (die Randformen stehen in `SPEC-061`, `AGENTS.md` §3.12), `BEO-REPO/plan-folgt-korrektur-nicht` (die Sendungen stehen im Plan des Nehmers). Keiner erreicht mit diesem Slice 3× neu.

alle berührten Sub-Areas GF (`REPO`).

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-18, im Commit des Schnitts):** Liefer-Punkte 2 (Hook-Skript | Gegenprobe mit Fragment). Schichten 1: Harness-Werkzeug (Skript, Gegenprobe, Fragment, Gate-Index in `harness/README.md`); Planung, Vertrag und ADR sind keine Schicht. Schätzung 610 bis 1035 in der Metrik von §4, Marke 900, Meldeschwelle 700 (Obergrenze darüber, Auflage zur Bündelung der Fälle in §4). Die Sendungen und Bindungen des Ursprungs dieses Slice gehören dem Nehmer `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` und sind dort gezählt; der Plan trägt zwei Sendungen an ihn ein (§1), die dort keinen Liefer-Punkt und keine Schicht anlegen.

**Herkunft der Zählung:** Der ungeteilte Slice hatte 3 Liefer-Punkte (Rollen-Anweisungen · Hook-Skript mit Gegenprobe · Fahrstelle im Dokument) und 2 Schichten (Harness-Werkzeug und -Artefakte · Wartungs-Doku), Schätzung 740 bis 1345 (Architect nach `SPEC-061`). Der Schnitt teilt nach Lieferwert: das Werkzeug (dieser Slice) und sein Betrieb mit Beschreibung (Nehmer).
