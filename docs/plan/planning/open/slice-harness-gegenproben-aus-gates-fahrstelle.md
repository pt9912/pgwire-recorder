# Slice slice-harness-gegenproben-aus-gates-fahrstelle: Pflicht und Fahrstelle der Gegenproben: Rollen-Anweisungen und der Gegenproben-Container in `docs/developer/devcontainer.md`

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

**Berührte Spec-Stellen:** der Vertrag der Sammelziele aus `slice-harness-gegenproben-aus-gates` (`SPEC-060`, Nummer vergibt der Architect dort; dieser Slice ändert ihn nicht und liest die Namen der Ziele daraus); sonst —.

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

**Ziel:** Die Pflicht, die Gegenproben zu fahren, steht dort, wo die Rollen sie lesen, und die Fahrstelle ist beschrieben und einmal belegt: (1) Die Rollen-Anweisungen (`implement-slice`, `close-welle`, Reviewer- und Verifier-Typ, Implementer-Typ, `AGENTS.md` §6) nennen, wann welche Gegenprobe zu fahren ist: bei einer Änderung am geprüften Werkzeug oder Gate, an der Gegenprobe oder an deren Vertrag; am Ende jeder Welle; vor dem Release; Review und Verifikation fahren die betroffene Gegenprobe gezielt (Regel des engsten Sensors). (2) `docs/developer/devcontainer.md` beschreibt in einem Abschnitt den Gegenproben-Container (Name vorläufig `pgr-gp`; eigener Port, eigenes Volume, eigener Klon über sein Bare-Repo; Schritte aus §3 bis §6 und §9 des Dokuments): `make gegenproben` (Name laut Vertrag in `slice-harness-gegenproben-aus-gates`) läuft dort im Hintergrund auf dem eingelieferten Stand eines Zweigs, während Implementer, Planner und Hook am Host und im Arbeitscontainer ungestört weiterarbeiten; `make gegenproben-host` läuft am Host.

**Zahlen, die der Abschnitt trägt, sind gemessen:** eine Fahrt von `make gegenproben` im Gegenproben-Container auf einem benannten Commit (Dauer, höchster Speicher des Containers, Last am Host, soweit der Implementer sie liest); was nicht gelaufen ist, steht als nicht gemessen (§13 des Dokuments).

**Sendung an `slice-erster-release-freigabe`** (eingetragen vom Planner im Commit, der diesen Plan anlegt, `AGENTS.md` §3.13): Die Freigabe-Checkliste des ersten Releases (`docs/maintainer/releasing.md` §10) verlangt als Punkt mit Beleg: `make gegenproben` im Gegenproben-Container und `make gegenproben-host` am Host grün auf dem Commit, der getaggt wird. Der Punkt gehört zu seinem dritten Liefer-Punkt (Checkliste); ein neuer Liefer-Punkt und eine neue Schicht entstehen nicht (Nachzählung in dessen §8).

**Übernimmt:** nichts.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Kettenumbau, die Sammelziele, ihre Namen, Listen und Gegenprobe, `AGENTS.md` §6 Schritt 6 als Satz — `slice-harness-gegenproben-aus-gates`; dieser Slice setzt dessen `done/` voraus (die Ziele müssen existieren, bevor eine Anweisung sie nennt).
- Ein Scheduler, eine CI-Pipeline oder ein Hook, der die Gegenproben anstößt — Schicht-Abgrenzung und Entscheidung des Nutzers: Der Container wird von Hand oder von einem Agenten gestartet; was ihn anstößt, ist Pflicht der Rollen, kein Werkzeug.
- Ein Ziel, das den Gegenproben-Container anlegt, einliefert und startet (`make container-einrichten` und ein Ziel für die Fahrt) — `slice-harness-devcontainer-einrichten` automatisiert Bare-Repo, Einlieferung, Klon und Zweig; dieser Slice beschreibt die Fahrt mit den Einzelbefehlen und ruft kein Ziel, das es noch nicht gibt. Ein Ziel für die Fahrt selbst ist ein anderer Vorgang (Werkzeug statt Beschreibung); eine Adresse gibt es nicht, ein Slice legt der Planner an, wenn der Nutzer es entscheidet.
- `u-boot` und `devcontainer` im Gegenproben-Container — nicht gebraucht: Die Gegenproben, die sie brauchen, laufen am Host (`make gegenproben-host`, Entscheidung des Nutzers vom 2026-10-18); sie in das Image zu bringen verlangte einen Umbau des Gerüsts ([ADR-0046](../../adr/0046-gegenprobe-ueberspringt-im-dev-container.md), Option C dort, nicht gewählt).
- Ein Stand oder Nachweis der Gegenproben und ein Hook, der ihn liest — `slice-harness-gegenproben-aus-gates` §1 und §6 (Frage 4 dort).
- Eine Messung der Last auf dem Host bei zwei Containern jenseits der einen Fahrt — Bestand bleibt ungemessen; der Abschnitt sagt es.
- Produkt-Code, Handbuch, `README.md` — Schicht-Abgrenzung: kein Verhalten des gebauten Binary.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Rollen-Anweisungen):** `.claude/commands/implement-slice.md`, `.claude/commands/close-welle.md`, `.claude/agents/implementer.md`, `.claude/agents/reviewer.md`, `.claude/agents/verifier.md`, `.harness/skills/reviewer.md` und `AGENTS.md` §6 nennen die Pflicht in einer Fassung: die Auslöser (Änderung am geprüften Werkzeug oder Gate, an der Gegenprobe, am Vertrag; Ende der Welle; Release), die Zuordnung Gegenprobe zu geprüftem Werkzeug über `harness/README.md` §Sensors, die gezielte Fahrt durch Review und Verifikation, und dass die Fahrt der eingelieferte, committete Stand ist; die Welle-Closure trägt die Fahrten beider Sammelziele als Beleg in die Closure-Notiz. Die Anweisungen sagen die Regel an den Schreiber und zeigen auf `harness/README.md` §Sensors als Quelle, sie zählen kein Verhalten der Gates auf (`AGENTS.md` §3.11).
- [ ] **Liefer-Punkt 2 (Fahrstelle):** `docs/developer/devcontainer.md` führt einen Abschnitt (Ort §10 *Mehrere Container nebeneinander* oder ein neuer Abschnitt dahinter, dazu eine Zeile in §8; der Implementer wählt) über den Gegenproben-Container mit dem Ablauf: Anlegen wie §3 bis §6 mit eigenen Werten (`PGR_NAME`, `PGR_SSH_PORT`, `PGR_CPUS`, `PGR_MEM`), Einliefern eines Zweigs (§9), Fahrt im Hintergrund auf dem benannten Commit mit Log und Exit-Datei im Volume, Lesen des Ergebnisses, Aufräumen. Jede Aussage über Dauer, Speicher und Last steht mit „gemessen“ und Datum oder mit „nicht gemessen“; jeder Befehl des Abschnitts ist gelaufen (`AGENTS.md` §3.11, Beispiele gelten erst als belegt, wenn sie gelaufen sind).
- [ ] Eine Fahrt von `make gegenproben` im Gegenproben-Container auf einem benannten Commit ist gelaufen; Ergebnis, Dauer und höchster Speicher stehen im Bericht und im Abschnitt; `make gegenproben-host` am Host auf demselben Commit ebenso.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `.claude/` ändert).
- [ ] Doku-Update: `harness/README.md` §Sensors nennt in der Zeile der Sammelziele die Fahrstelle; `docs/developer/devcontainer.md` §8 trägt die Zeile des neuen Ziels in der Tabelle der Messungen.
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
| `.claude/agents/implementer.md`, `reviewer.md`, `verifier.md`, `.harness/skills/reviewer.md` | update | Stichprobe nach der Regel des engsten Sensors; Reviewer und Verifier fahren die betroffene Gegenprobe gezielt |
| `AGENTS.md` | update | §6 Schritt 6 (der Satz aus `slice-harness-gegenproben-aus-gates` wird um die Auslöser ergänzt, falls dort kürzer) |
| `docs/developer/devcontainer.md` | update | Abschnitt über den Gegenproben-Container, Zeile in §8 |
| `harness/README.md` | update | Fahrstelle in der Zeile der Sammelziele |
| `docs/plan/planning/next/slice-erster-release-freigabe.md` | update (Planner, im Commit des Plans) | Sendung der Checkliste (§1, §6, §8) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. **Vorbedingung:** `slice-harness-gegenproben-aus-gates` liegt in `done/` (die Ziele und ihre Namen stehen); der Dev-Container-Ablauf von `slice-harness-devcontainer-einrichten-geruest` liegt in `done/` (der Abschnitt baut auf `make container-start`, soweit es dann steht; sonst auf §3 bis §6 des Dokuments). Reihenfolge: erste Stelle nach `slice-harness-gegenproben-aus-gates`; der Slice wartet nicht auf `welle-harness-regeln` (Drift-Log der Roadmap). `Verantwortlich:` trägt beim Start Person und Zweig; er läuft auf eigenem Zweig.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- AGENTS.md .claude .harness/skills docs/developer harness/README.md`. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt.

**Schätzung (Planner, 2026-10-18):** Rollen-Anweisungen sieben Dateien, je 6 bis 20 Zeilen: 50 bis 130; Abschnitt im Dokument mit Befehlen und Messung: 80 bis 160; `harness/README.md`, `AGENTS.md`: 10 bis 30. **Zusammen 140 bis 320**; Meldeschwelle 700, Marke 900, weit darunter. Es ist kein Code im Slice, also kein Faktor für Tests.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder die Fahrt braucht ein Werkzeug, das nicht existiert (dann Schnitt: Beschreibung, Ziel in eigenem Slice).
- `in-progress` → `open` (blockiert): Die Fahrt von `make gegenproben` ist im Gegenproben-Container nicht grün (Befund an den Planner, Rückgabe an `slice-harness-gegenproben-aus-gates`, Risiko *Zustand eines anderen Gates* dort).

## 5. Closure-Trigger

DoD vollständig; eine Fahrt je Sammelziel auf einem benannten Commit belegt; `make gates` grün; Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag; die drei Paarungen getragen.

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

Risiken:

- **Die Pflicht wird nicht gelesen oder nicht befolgt, und die Regression bleibt unbemerkt** (Risiko 1 von `slice-harness-gegenproben-aus-gates`). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`.
- **Die Beschreibung der Fahrstelle veraltet**, sobald `slice-harness-devcontainer-einrichten` die Schritte in ein Ziel legt (`make container-einrichten`) und §6 und §9 des Dokuments umschreibt, auf die dieser Abschnitt verweist. — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`; die Bindung steht als Sendung in `slice-harness-devcontainer-einrichten` §1 und §6 (eingetragen vom Planner im Commit, der diesen Plan anlegt, Nachzählung dort in §8): Wer §6 und §9 umschreibt, zieht den Abschnitt nach oder meldet es dem Planner.
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

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Liefer-Punkte 2 (Rollen-Anweisungen | Fahrstelle im Dokument mit einer Fahrt als Beleg). Schichten 2: Harness-Artefakte (Rollen-Anweisungen, `AGENTS.md`, `harness/README.md` als Gate-Index) und Wartungs-Doku (`docs/developer/`, nach der Teilung von `slice-harness-d-check-v0-85` §8 und `AGENTS.md` §3.13 eine Schicht); Planung und ADR sind keine. Schätzung 140 bis 320, weit unter der Marke 900. Der Plan trägt eine Sendung in einen Nehmer ein (`slice-erster-release-freigabe`, Nachzählung dort in §8) und eine Bindung (`slice-harness-devcontainer-einrichten`, §1, §6 und §8 dort).
