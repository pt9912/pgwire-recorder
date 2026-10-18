# Slice slice-harness-gegenproben-aus-gates-fahrstelle: Pflicht und Fahrstelle der Gegenproben: Rollen-Anweisungen, Hook für die Einlieferung und der Gegenproben-Container in `docs/developer/devcontainer.md`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten, zweite Stelle der Reihenfolge, direkt hinter `slice-harness-gegenproben-aus-gates` und vor `slice-harness-gegenprobe-baumkopie` (Entscheidung des Nutzers vom 2026-10-18; wie jener Slice aus der Reihenfolge der Wellen gelöst: Start-Trigger ist `slice-harness-gegenproben-aus-gates` in `done/`, nicht `welle-harness-regeln`; der Schnitt vom Planner am selben Tag, damit der Kettenumbau unter der Marke bleibt). Angelegt am 2026-10-18 als Folge-Slice von `slice-harness-gegenproben-aus-gates`. Der Slice liegt in `open/`: Er ist nicht priorisiert, den Anspruch trägt der Nutzer.

**Bezug:** [LH-QA-07](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur insoweit, als die Gegenproben dessen Messungen tragen; kein Anspruch des Produkts, der Slice ändert Rollen-Anweisungen und Wartungs-Doku des Harness. Bindung an das Register: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` und das mit `slice-harness-gegenproben-aus-gates` entstehende `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`. Keine aktive ADR berührt.

**Berührte Spec-Stellen:** der Vertrag der Sammelziele aus `slice-harness-gegenproben-aus-gates` (`SPEC-060`, Nummer vergibt der Architect dort; dieser Slice ändert ihn nicht und liest die Namen der Ziele daraus); neuer Vertrag des Hooks für die Einlieferung in `harness/vertraege/` (nächste freie Nummer nach `SPEC-060`, also `SPEC-061`; der Architect vergibt sie vor dem Code).

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

**Ziel:** Die Pflicht, die Gegenproben zu fahren, steht dort, wo die Rollen sie lesen, und die Fahrstelle läuft ohne Zutun und ist beschrieben und einmal belegt: (1) Die Rollen-Anweisungen (`implement-slice`, `close-welle`, Reviewer- und Verifier-Typ, Implementer-Typ, `AGENTS.md` §6) nennen, wann welche Gegenprobe zu fahren ist: bei einer Änderung am geprüften Werkzeug oder Gate, an der Gegenprobe oder an deren Vertrag; am Ende jeder Welle; vor dem Release; Review und Verifikation fahren die betroffene Gegenprobe gezielt (Regel des engsten Sensors). Die Closure von Welle und Release prüft, dass für den Stand ein grünes Ergebnis des Gegenproben-Containers vorliegt (Hash des Commits); fehlt es, fährt die Closure `make gegenproben` selbst. (2) Ein Skript im Repo (Vorschlag: `tools/harness/gegenproben-nach-einlieferung.sh`, Name laut Vertrag) ist der Inhalt des `post-receive`-Hooks im Bare-Repo des Gegenproben-Containers: Jede Einlieferung eines Zweigs startet dort `make gegenproben` (nur die Container-Liste) im Hintergrund auf genau diesem Stand und legt das Ergebnis ab. (3) `docs/developer/devcontainer.md` beschreibt in einem Abschnitt den Gegenproben-Container (Name vorläufig `pgr-gp`; eigener Port, eigenes Volume, eigener Klon über sein Bare-Repo; Schritte aus §3 bis §6 und §9 des Dokuments) mit der Zeile zum Aufsetzen des Hooks; `make gegenproben-host` läuft am Host.

**`pgr-gp` ist ein spezieller Container, der nur die Gegenproben fährt** (Ergänzung des Nutzers vom 2026-10-18); in ihm wird kein Slice abgearbeitet. Implementer, Planner und Hook arbeiten am Host und in den Arbeits-Containern (`pgr-dev`, `pgr-dev2`, …) ungestört weiter, während `pgr-gp` im Hintergrund läuft. Die Arbeits-Container bekommen den Hook nicht (unten, *Ausdrücklich NICHT*).

**Zahlen, die der Abschnitt trägt, sind gemessen:** eine Fahrt von `make gegenproben` im Gegenproben-Container auf einem benannten Commit (Dauer, höchster Speicher des Containers, Last am Host, soweit der Implementer sie liest); was nicht gelaufen ist, steht als nicht gemessen (§13 des Dokuments).

**Sendung an `slice-erster-release-freigabe`** (eingetragen vom Planner im Commit, der diesen Plan anlegt, `AGENTS.md` §3.13): Die Freigabe-Checkliste des ersten Releases (`docs/maintainer/releasing.md` §10) verlangt als Punkt mit Beleg: ein grünes Ergebnis des Gegenproben-Containers für den Commit, der getaggt wird (fehlt es, `make gegenproben` im Gegenproben-Container), und `make gegenproben-host` am Host grün auf demselben Commit. Der Punkt gehört zu seinem dritten Liefer-Punkt (Checkliste); ein neuer Liefer-Punkt und eine neue Schicht entstehen nicht (Nachzählung in dessen §8).

**Sendung an die Welle-Closure** (`close-welle`, Liefer-Punkt 1 dieses Slice): Sie prüft das Ergebnis des Gegenproben-Containers für den Stand und fährt sonst `make gegenproben`; die Sendung an `slice-erster-release-freigabe` trägt dasselbe für den Release (oben).

**Übernimmt:** nichts.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Kettenumbau, die Sammelziele, ihre Namen, Listen und Gegenprobe, `AGENTS.md` §6 Schritt 6 als Satz — `slice-harness-gegenproben-aus-gates`; dieser Slice setzt dessen `done/` voraus (die Ziele müssen existieren, bevor eine Anweisung sie nennt).
- Ein Scheduler, eine CI-Pipeline, ein Warteschlangen-Dienst oder eine Benachrichtigung — Entscheidung des Nutzers: Der Container wird von Hand oder von einem Agenten gestartet; der `post-receive`-Hook im Bare-Repo von `pgr-gp` ist die kleinste Form, die Einlieferung und Fahrt verbindet, und nicht mehr.
- Der Hook in den Bare-Repos der Arbeits-Container (`pgr-dev`, `pgr-dev2`, …) — Schicht-Abgrenzung und Last: Der Arbeits-Container soll nicht ausgebremst werden; eine Fahrt von `make gegenproben` belegt Docker-Daemon und Kerne (§6, Ressourcen), und der Arbeits-Container ist der, in dem der Implementer arbeitet. Das Skript lehnt den Lauf ab, wo die Marke des Gegenproben-Containers fehlt (§6, H9).
- Ein Ziel, das den Gegenproben-Container anlegt, einliefert und startet (`make container-einrichten` und ein Ziel für die Fahrt) — `slice-harness-devcontainer-einrichten` automatisiert Bare-Repo, Einlieferung, Klon und Zweig; dieser Slice beschreibt das Aufsetzen mit den Einzelbefehlen und einer Zeile für den Hook und ruft kein Ziel, das es noch nicht gibt. Ein Ziel für das Aufsetzen selbst ist ein anderer Vorgang (Werkzeug statt Beschreibung); eine Adresse gibt es nicht, ein Slice legt der Planner an, wenn der Nutzer es entscheidet.
- `u-boot` und `devcontainer` im Gegenproben-Container — nicht gebraucht: Die Gegenproben, die sie brauchen, laufen am Host (`make gegenproben-host`, Entscheidung des Nutzers vom 2026-10-18); sie in das Image zu bringen verlangte einen Umbau des Gerüsts ([ADR-0046](../../adr/0046-gegenprobe-ueberspringt-im-dev-container.md), Option C dort, nicht gewählt).
- Ein Stand oder Nachweis der Gegenproben, den der Stop-Hook am Host liest — `slice-harness-gegenproben-aus-gates` §1 und §6 (Frage 4 dort, Empfehlung des Planners: nein). Das Ergebnis, das der `post-receive`-Hook im Gegenproben-Container ablegt, liegt im Container (ignorierter Pfad im Klon), nicht im Arbeitsbaum des Hosts; der Hauptbaum und der Stop-Hook bleiben unberührt, nur die Closure liest das Ergebnis, per `ssh`.
- Eine Messung der Last auf dem Host bei zwei Containern jenseits der einen Fahrt — Bestand bleibt ungemessen; der Abschnitt sagt es.
- Produkt-Code, Handbuch, `README.md` — Schicht-Abgrenzung: kein Verhalten des gebauten Binary.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Rollen-Anweisungen):** `.claude/commands/implement-slice.md`, `.claude/commands/close-welle.md`, `.claude/agents/implementer.md`, `.claude/agents/reviewer.md`, `.claude/agents/verifier.md`, `.harness/skills/reviewer.md` und `AGENTS.md` §6 nennen die Pflicht in einer Fassung: die Auslöser (Änderung am geprüften Werkzeug oder Gate, an der Gegenprobe, am Vertrag; Ende der Welle; Release), die Zuordnung Gegenprobe zu geprüftem Werkzeug über `harness/README.md` §Sensors, die gezielte Fahrt durch Review und Verifikation, und dass die Fahrt der eingelieferte, committete Stand ist; die Closure von Welle und Release prüft, dass für den Stand ein grünes Ergebnis des Gegenproben-Containers vorliegt, fährt sonst `make gegenproben` selbst und trägt Ergebnis und Hash (beide Hashes, wenn nach dem Lauf nur inerte Pfade änderten) in die Closure-Notiz, dazu die Fahrt von `make gegenproben-host`. Die Anweisungen sagen die Regel an den Schreiber und zeigen auf `harness/README.md` §Sensors als Quelle, sie zählen kein Verhalten der Gates auf (`AGENTS.md` §3.11).
- [ ] **Liefer-Punkt 2 (Hook für die Einlieferung):** Das Skript (Name laut Vertrag) ist der Inhalt des `post-receive`-Hooks im Bare-Repo von `pgr-gp`: Es startet je eingeliefertem Zweig `make gegenproben` im Hintergrund auf genau dem eingelieferten Hash (eigener Arbeitsbaum je Lauf, der Klon bleibt unberührt), legt bei beendetem Lauf eine Ergebnisdatei ab (Zweig, Hash, Exit, Beginn, Ende, rote Ziele; Pfad nach dem Muster `.harness/state/gegenproben-<Zweig>-<Hash>` im Klon), nie bei abgebrochenem Lauf, und endet selbst sofort mit Exit 0. Entscheidung je Randform aus §6 (Gruppe H) steht im Vertrag vor dem ersten Code-Commit; eine Gegenprobe des Skripts (Fragment unter `harness/mk/`, an der Container-Liste aus `slice-harness-gegenproben-aus-gates`, Zeile im Gate-Index) arbeitet in Temp-Repos mit einer Attrappe als `make` und trägt je Randform und je Element ihrer Aufzählungen einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10).
- [ ] **Liefer-Punkt 3 (Fahrstelle im Dokument):** `docs/developer/devcontainer.md` führt einen Abschnitt (Ort §10 *Mehrere Container nebeneinander* oder ein neuer Abschnitt dahinter, dazu eine Zeile in §8; der Implementer wählt) über den Gegenproben-Container mit dem Ablauf: Anlegen wie §3 bis §6 mit eigenen Werten (`PGR_NAME`, `PGR_SSH_PORT`, `PGR_CPUS`, `PGR_MEM`), die Zeile zum Aufsetzen des Hooks (Marke und Kopie des Skripts im Bare-Repo), Einliefern eines Zweigs (§9), Lesen des Ergebnisses für einen Hash und des Logs per `ssh`, Aufräumen; dass die Arbeits-Container den Hook nicht bekommen. Jede Aussage über Dauer, Speicher und Last steht mit „gemessen“ und Datum oder mit „nicht gemessen“; jeder Befehl des Abschnitts ist gelaufen (`AGENTS.md` §3.11, Beispiele gelten erst als belegt, wenn sie gelaufen sind).
- [ ] Eine Einlieferung eines Zweigs in den Gegenproben-Container hat `make gegenproben` ausgelöst und ein Ergebnis für den Hash abgelegt (Ergebnis, Dauer und höchster Speicher stehen im Bericht und im Abschnitt); `make gegenproben-host` am Host auf demselben Commit ebenso.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `.claude/` ändert).
- [ ] Doku-Update: `harness/README.md` §Sensors nennt in der Zeile der Sammelziele die Fahrstelle und führt die Zeile der Gegenprobe des Hook-Skripts; `docs/developer/devcontainer.md` §8 trägt die Zeile des neuen Ziels in der Tabelle der Messungen.
- [ ] Beobachtungs-Register: `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt` mit dem Beleg dieses Slice geführt (Mitigation getragen: ja oder nein); kein Zähler wird gesetzt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review. Das Review fährt die Stichprobe nach der neuen Regel selbst.
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
| `.claude/commands/implement-slice.md`, `.claude/commands/close-welle.md` | update | Pflicht im Workflow (Schritt „Gates“ und Sensor-Belege) und in der Welle-Closure (Fahrten beider Sammelziele als Beleg) |
| `harness/vertraege/SPEC-<NNN>-gegenproben-nach-einlieferung.md` | neu (Architect, vor dem Code) | Vertrag des Hook-Skripts: Aufruf, Ergebnisdatei, Entscheidung je Randform aus Gruppe H |
| `tools/harness/gegenproben-nach-einlieferung.sh` | neu | Inhalt des `post-receive`-Hooks: liest die Zeilen von `git` auf stdin, startet je Zweig einen Lauf im Hintergrund, legt das Ergebnis ab |
| `tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh`, `harness/mk/gegenproben-nach-einlieferung-gegenprobe.mk` | neu | Gegenprobe in Temp-Repos mit Attrappe als `make` (kein Docker): je Randform ein Fall und eine Mutation |
| `.claude/agents/implementer.md`, `reviewer.md`, `verifier.md`, `.harness/skills/reviewer.md` | update | Stichprobe nach der Regel des engsten Sensors; Reviewer und Verifier fahren die betroffene Gegenprobe gezielt |
| `AGENTS.md` | update | §6 Schritt 6 (der Satz aus `slice-harness-gegenproben-aus-gates` wird um die Auslöser ergänzt, falls dort kürzer) |
| `docs/developer/devcontainer.md` | update | Abschnitt über den Gegenproben-Container, Zeile in §8 |
| `harness/README.md` | update | Fahrstelle in der Zeile der Sammelziele |
| `docs/plan/planning/next/slice-erster-release-freigabe.md` | update (Planner, im Commit des Plans) | Sendung der Checkliste (§1, §6, §8) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. **Vorbedingung:** `slice-harness-gegenproben-aus-gates` liegt in `done/` (die Ziele und ihre Namen stehen); der Architect hat den Vertrag des Hook-Skripts in `harness/vertraege/` geschrieben, mit der Entscheidung je Randform aus Gruppe H (`AGENTS.md` §3.12), und entscheidet, ob eine ADR entsteht (Commits des Architect vor dem Move nach `in-progress/`); der Dev-Container-Ablauf von `slice-harness-devcontainer-einrichten-geruest` liegt in `done/` (der Abschnitt baut auf `make container-start`, soweit es dann steht; sonst auf §3 bis §6 des Dokuments). Reihenfolge: erste Stelle nach `slice-harness-gegenproben-aus-gates`; der Slice wartet nicht auf `welle-harness-regeln` (Drift-Log der Roadmap). `Verantwortlich:` trägt beim Start Person und Zweig; er läuft auf eigenem Zweig.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- AGENTS.md .claude .harness/skills docs/developer harness/README.md`. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt.

**Schätzung (Planner, 2026-10-18; Metrik wie oben, Metrik-Pfade um `tools` erweitert; Gegenprobe-Code mit dem Doppelten des Produktcodes):**

| Teil | Zeilen |
|---|---|
| Hook-Skript (Zeilen von `git` lesen, Zweig und Hash prüfen, Marke, Hintergrund, Ergebnis) | 70 bis 110 |
| Gegenprobe des Hook-Skripts (11 Randformen mit Mutation, Attrappen für `make` und `git`-Repos) | 200 bis 330 |
| Rollen-Anweisungen, sieben Dateien, je 6 bis 20 Zeilen | 50 bis 130 |
| Abschnitt im Dokument mit Befehlen, Zeile zum Aufsetzen und Messung | 80 bis 160 |
| `harness/README.md` (Zeile der Gegenprobe, Fahrstelle), `AGENTS.md` | 10 bis 40 |
| **Zusammen** | **410 bis 770** |

Meldeschwelle 700, **Marke 900**. Die Mitte (590) liegt bei 66 % der Marke, die Obergrenze (770) bei 86 %, über der Meldeschwelle; ein weiterer Schnitt unterbleibt. Die Spanne hängt an der Zahl der Randformen der Gegenprobe; zählt der Architect mehr als 11, rechnet er neu.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder der Vertrag des Hooks verlangt mehr als 11 Randformen (dann Schnitt: Hook-Skript mit Gegenprobe getrennt von Rollen-Anweisungen und Dokument).
- `in-progress` → `open` (blockiert): Die Fahrt von `make gegenproben` ist im Gegenproben-Container nicht grün (Befund an den Planner, Rückgabe an `slice-harness-gegenproben-aus-gates`, Risiko *Zustand eines anderen Gates* dort).

## 5. Closure-Trigger

DoD vollständig; eine Einlieferung in den Gegenproben-Container hat ein Ergebnis für den Hash abgelegt, `make gegenproben-host` auf demselben Commit am Host belegt; `make gates` grün; Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag; die drei Paarungen getragen.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Der Slice liefert keinen neuen Vertrag im Sinne von `AGENTS.md` §3.10 (kein Code); die Randformen der Fahrstelle sind Aussagen des Dokuments und werden dort belegt oder als nicht gemessen geführt (§3.11). Sie stehen hier als Randformen (§3.12), jede mit Entscheidung:

1. **Eingelieferter Stand gegen Arbeitsbaum.** Die Fahrt prüft den committeten Stand, den der Container im Bare-Repo hat, nicht den Arbeitsbaum des Hosts (`devcontainer.md` §9: uncommittete und ignorierte Dateien kommen nicht mit). Entscheidung (Empfehlung): Der Bericht nennt den Commit; eine Fahrt auf einem anderen Commit als dem getaggten oder zu schließenden gilt nicht. Entschieden: offen (Architect/Implementer im Dokument).
2. **Zwei Container desselben Klons** (Arbeitscontainer und Gegenproben-Container teilen ihren Klon nicht; zwei Fahrten im selben Klon gleichzeitig sind nicht zugesagt, `SPEC-052` Punkt 8). Entscheidung: eine Fahrt je Container; der Abschnitt sagt es.
3. **Ressourcen.** Je Container 8 CPUs und 12 GiB (Vorgabe `PGR_CPUS`, `PGR_MEM`), Host 20 CPUs und 31 GiB: Arbeits- und Gegenproben-Container zusammen 16 CPUs und 24 GiB, am Host bleiben 4 CPUs und 7 GiB. Gemessen ist (Dokument §8) `make gates` mit 1701 s und höchstens 5,42 GiB; für die Gegenproben ist nichts gemessen. Empfehlung: Standardwerte, `PGR_CPUS` und `PGR_MEM` des Gegenproben-Containers senken, wenn die Messung des Implementers Last am Host zeigt; der Abschnitt nennt die Messung oder „nicht gemessen“.
4. **Dev-Container-Erkennung nach [ADR-0046](../../adr/0046-gegenprobe-ueberspringt-im-dev-container.md) entfällt** ([ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md), Accepted) (Entscheidung des Nutzers vom 2026-10-18): Alle Gegenproben in der Container-Liste laufen im Gegenproben-Container ohne Überspringen; was `u-boot` und `devcontainer` braucht, läuft nur am Host (`make gegenproben-host`). Vorschlag zur Frage „u-boot im Gegenproben-Container?“: nein, weil die Aufteilung den Bedarf aufhebt und das Image sonst vom Gerüst von `u-boot` abwiche. Entschieden: Nutzer (Aufteilung); der Vorschlag zum Image offen.
5. **Hintergrundlauf.** Die Fahrt läuft mit `setsid` und Log im Volume, überlebt das Ende der `ssh`-Sitzung; ein Neustart des Containers beendet sie (§8 des Dokuments: Prozesse beginnen neu). Entschieden: offen (Dokument).
6. **Fahrt ohne Netz.** Der Daemon im Container lädt die Images selbst (Egress eingeschränkt, `egress.allow`); fehlt ein Host der Liste, endet die Fahrt rot mit der Meldung der Gegenprobe, nicht still. Entschieden: Bestand (§11 des Dokuments).
7. **Wer die Pflicht auslöst.** Das Werkzeug zur Zuordnung „geändertes Werkzeug → Gegenprobe“ gibt es nicht (eine Tabelle wäre die Eingabeliste, die der Nutzer verworfen hat); der Schreiber der Anweisung nennt die Quelle, `harness/README.md` §Sensors. Akzeptiertes Negativ; Grund: ein Sensor dafür wäre der erste Zuschnitt.

**Gruppe H — Hook für die Einlieferung** (Randformen des neuen Vertrags, `AGENTS.md` §3.12; Empfehlungen des Planners, entschieden ist keine, der Architect legt sie vor dem Code in den Vertrag):

1. **Der Push wartet nicht.** `post-receive` läuft im Push; der Hook startet den Lauf im Hintergrund (`setsid`, Ausgabe in eine Logdatei) und endet sofort mit Exit 0. `post-receive` kann den Push nicht ablehnen; ein Fehler des Hooks darf ihn nicht blockieren. Fall: Hook endet, bevor der Lauf endet (Attrappe mit `sleep`). Entschieden: offen.
2. **Welche Refs.** Nur `refs/heads/*`; Tags und andere Refs lösen keinen Lauf aus. Mehrere Zweige in einem Push: je Zweig ein Lauf, nacheinander in der Reihenfolge der Zeilen. Entschieden: offen.
3. **Zwei Einlieferungen kurz hintereinander.** Der laufende Lauf wird nicht abgebrochen; der zweite wartet (eine Sperre mit `flock` im Volume, ein Lauf zur Zeit), und alle Läufe laufen nacheinander. Empfehlung: nacheinander statt Abbruch, weil das Ergebnis je Hash gebraucht wird (Closure) und Docker-Daemon und Kerne nicht doppelt belastet werden sollen; ein Abbruch des älteren Laufs ließe dessen Hash ohne Ergebnis. Alternative: ein wartender Lauf desselben Zweigs entfällt, wenn ein jüngerer Hash nachkommt; nicht empfohlen (Ergebnis für den älteren Hash fehlt). Fall: zwei Einlieferungen, die Läufe überlappen nicht. Entschieden: offen.
4. **Zweig gelöscht** (`git push --delete`, neuer Hash aus Nullen) — kein Lauf, kein Fehler; vorhandene Ergebnisdateien bleiben. Entschieden: offen.
5. **Derselbe Hash erneut** (erneuter Push, Force-Push auf einen früheren Stand, anderer Zweigname): ein neuer Lauf; ein vorhandenes Ergebnis wird nicht als Stand wiederverwendet (der Stand-für-den-Hook bleibt „nein“) und nach dem Lauf überschrieben. Entschieden: offen.
6. **Name des Zweigs im Dateinamen.** Ein Zweig mit Schrägstrich (`thema/teil`) und andere Sonderzeichen werden im Dateinamen kodiert (Vorschlag: `/` wird `%2F`); jede Form, die `git check-ref-format` zulässt, landet in genau einer Datei. Je Form ein Fall (Schrägstrich, `%` im Namen, Punkt am Anfang). Entschieden: offen.
7. **Ergebnis nur bei beendetem Lauf.** Datei nach Ende, über eine Temp-Datei und `mv`, mit Exit des Laufs; bei `SIGTERM`, `SIGKILL` (kein Fall mit `SIGKILL`, kein Handler) oder Neustart des Containers existiert keine Ergebnisdatei. Eine Marke „läuft“ entsteht beim Start und entfällt am Ende; eine verwaiste Marke bedeutet „abgebrochen“ und ist kein Ergebnis. Fälle: Exit 0, Exit 1, `SIGTERM` mitten im Lauf. Entschieden: offen.
8. **Nur die Container-Liste.** Der Hook ruft `make gegenproben`, nie `make gegenproben-host`; die Host-Liste läuft am Host. Entschieden: Nutzer, 2026-10-18.
9. **Der Hook gehört nur dem Gegenproben-Container.** Das Skript läuft nur, wenn die Marke des Containers vorliegt (Datei, die die Zeile zum Aufsetzen anlegt); sonst Meldung auf stderr (`remote:` beim Pusher) und Exit 0, kein Lauf. Das fängt den Hook im Bare-Repo eines Arbeits-Containers, falls er versehentlich dort aufgesetzt wird (der Arbeits-Container würde ausgebremst). Fall: Marke fehlt, es startet nichts. Entschieden: offen.
10. **Wohin gehört der Hook-Code?** Das Skript liegt im Repo (`tools/harness/`), versioniert und prüfbar; der Hook im Bare-Repo (`hooks/post-receive`, eine Datei außerhalb des Repos) ist eine **Kopie** davon, kein Symlink: Der Lauf checkt den eingelieferten Hash in einen eigenen Arbeitsbaum aus, ein Symlink auf den Klon zeigte auf eine Datei, die ein Wechsel des Klons ersetzt, und `bash` liest ein laufendes Skript fortlaufend. Das Dokument trägt die Zeile zum Aufsetzen (`install -m 755`, dazu die Marke); eine Abweichung der Kopie vom Stand des Repos (nach einer Änderung am Skript) erkennt keine Gegenprobe (akzeptiertes Negativ; Grund: die Kopie liegt außerhalb des Repos, die Zeile zum Aufsetzen wiederholt man nach jeder Änderung des Skripts; das Dokument sagt es). Entschieden: offen.
11. **Berechtigung und Zustand.** Der Hook läuft als `vscode` (der Benutzer von `git-receive-pack`, `devcontainer.md` §6): Ist die Datei nicht ausführbar, ignoriert `git` den Hook mit einer Warnung, der Push gelingt und es läuft nichts (das Dokument nennt die Probe: Testeinlieferung und Ergebnisdatei). Fehlen Schreibrecht für die Ablage, der Klon, der Docker-Daemon oder das Netz, meldet das Skript es auf stderr und im Log und endet mit Exit 0; ein Lauf, der deshalb nicht startet, hat kein Ergebnis (nicht „grün“). Reihenfolge der Gründe: Marke, Ref, Sperre, Arbeitsbaum, Lauf. Entschieden: offen.
12. **Welcher Hash zählt der Closure?** Die Closure sucht das Ergebnis für den Commit, der den letzten nicht inerten Stand trägt (inert nach `SPEC-052` Punkt 12: Planung, Reviews, Dokumentation); ein Commit, der nur inerte Pfade ändert, braucht keinen neuen Lauf, die Closure-Notiz nennt beide Hashes. Fehlt ein Ergebnis oder ist es rot, fährt die Closure `make gegenproben` selbst. Der Planner kennt keinen Sensor dafür (Urteil von Closure und Review). Entschieden: offen.

**Aufrufarten des Treibers (§3.12), Gruppe H:** Der Hook ruft `make gegenproben` ohne Flags und ohne `MAKEFLAGS` aus der Umgebung des Pushs (die Umgebung eines `git`-Hooks ist leer von `MAKEFLAGS`; das Skript setzt `MAKEFLAGS=` ausdrücklich zurück, sonst erbte ein Lauf `-i`, und das Sammelskript lehnte ihn ab, G4). **Aufruf-Voraussetzung (§3.12):** Zusage nur für Läufe nacheinander (H3, Sperre); die Zusagen von `git` über gleichzeitige Pushes in dasselbe Bare-Repo gelten nicht für die Ergebnisdateien.

Risiken:

- **Die Pflicht wird nicht gelesen oder nicht befolgt, und die Regression bleibt unbemerkt** (Risiko 1 von `slice-harness-gegenproben-aus-gates`). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`.
- **Die Beschreibung der Fahrstelle veraltet**, sobald `slice-harness-devcontainer-einrichten` die Schritte in ein Ziel legt (`make container-einrichten`) und §6 und §9 des Dokuments umschreibt, auf die dieser Abschnitt verweist. — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`; die Bindung steht als Sendung in `slice-harness-devcontainer-einrichten` §1 und §6 (eingetragen vom Planner im Commit, der diesen Plan anlegt, Nachzählung dort in §8): Wer §6 und §9 umschreibt, zieht den Abschnitt nach oder meldet es dem Planner.
- **Der Hook läuft bei jeder Einlieferung, auch bei unfertigen Zweigen, und belegt den Docker-Daemon des Gegenproben-Containers lange** (28 min 21 s für `make gates` im Container sind gemessen, die Gegenproben nicht). Die Läufe stauen sich (H3). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, falls die Messung des Implementers Staus zeigt; sonst entfallen mit der Zahl.
- **Der Hook liegt als Kopie außerhalb des Repos und veraltet** (H10). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`.
- **Das Ergebnis liegt im Container und der Container ist weg** (Volume gelöscht, Neustart, `devcontainer.md` §4: ein Neustart behält das Volume). — **Ausgang (vorgesehen):** entfallen, weil die Closure bei fehlendem Ergebnis selbst fährt (H12); das Ergebnis ist ein Nachweis zum Nachlesen, kein Stand.
- **Eine Fahrt kostet Zeit und Platz im Container** (28 min 21 s für `make gates` im Container sind gemessen, die Gegenproben nicht; ein zweites Volume mit eigenem Docker-Cache, `devcontainer.md` §10 und §11). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, falls die Messung des Implementers die Dauer als Zumutung zeigt; sonst entfallen mit der Zahl.
- **Die Sendung an `slice-erster-release-freigabe` nimmt der Nehmer nicht an**, weil er vor diesem Slice schließt. — **Ausgang (vorgesehen):** eingetreten → die Checkliste trägt den Punkt nachträglich, Folge-Slice mit Kennung vom Planner; die Reihenfolge der Wellen spricht dagegen (welle-erster-release liegt vor den Harness-Wellen) — der Planner prüft sie beim Start dieses Slice und trägt das Ergebnis in §7 ein.

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

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Hintergrund), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Dokument und Anweisungen sagen nur, was gelaufen ist; `AGENTS.md` §3.11), `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (die Fahrt im eigenen Container ist die Antwort darauf für die Gegenproben, nicht Gegenstand), `BEO-REPO/plan-folgt-korrektur-nicht` (die Sendungen stehen in den Plänen der Nehmer). Keiner erreicht mit diesem Slice 3× neu.

alle berührten Sub-Areas GF (`REPO`).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der die Ergänzung „nur Gegenproben, Hook bei Einlieferung“ einträgt):** Liefer-Punkte 3 (Rollen-Anweisungen mit Prüfung des Ergebnisses in der Closure | Hook-Skript mit Gegenprobe | Fahrstelle im Dokument mit einer Einlieferung als Beleg). Schichten 2: Harness-Werkzeug und -Artefakte (Hook-Skript mit Gegenprobe, Rollen-Anweisungen, `AGENTS.md`, `harness/README.md` als Gate-Index) und Wartungs-Doku (`docs/developer/`, nach der Teilung von `slice-harness-d-check-v0-85` §8 und `AGENTS.md` §3.13 eine Schicht); Planung, Vertrag und ADR sind keine. Schätzung 410 bis 770 in der Metrik von §4, Marke 900, Meldeschwelle 700 (Obergrenze darüber). Der Plan trägt zwei Sendungen und eine Bindung ein: `slice-erster-release-freigabe` (Ergebnis des Gegenproben-Containers oder Fahrt in der Checkliste; Nachzählung dort in §8, im selben Commit), die Welle-Closure (Liefer-Punkt 1 dieses Slice; `welle-gegenproben-und-mutanten` §3 nennt sie) und `slice-harness-devcontainer-einrichten` (Bindung an §6 und §9 des Dokuments, unverändert).
