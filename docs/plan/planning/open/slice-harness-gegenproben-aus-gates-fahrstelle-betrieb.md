# Slice slice-harness-gegenproben-aus-gates-fahrstelle-betrieb: Pflicht und Betrieb der Gegenproben: Rollen-Anweisungen, Abschnitt über den Gegenproben-Container in `docs/developer/devcontainer.md` und die erste Fahrt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten, vierte Stelle der Reihenfolge, direkt hinter `slice-harness-gegenproben-aus-gates-fahrstelle` und vor `slice-harness-gegenprobe-baumkopie` (Entscheidung des Nutzers vom 2026-10-18 für die Stelle der Gegenproben-Slices vor der Baumkopie; Reihenfolge innerhalb: Schnitt vom Planner am 2026-10-18, Reihenfolge Hook, dann Betrieb). Wie die drei Slices davor aus der Reihenfolge der Wellen gelöst: Start-Trigger ist `slice-harness-gegenproben-aus-gates-fahrstelle` in `done/`, nicht `welle-harness-regeln`. **Entstanden am 2026-10-18 aus dem Schnitt von `slice-harness-gegenproben-aus-gates-fahrstelle`** (Entscheidung des Nutzers: Schnitt ja; Regel des Nutzers: schneiden, wenn die Schätzung die Marke selbst überschreitet): Der ungeteilte Slice war nach `SPEC-061` auf 740 bis 1345 Zeilen geschätzt (Marke 900). Dieser Slice trägt die Hälfte „Pflicht, Beschreibung, erste Fahrt“; der Hook mit seiner Gegenprobe bleibt in `slice-harness-gegenproben-aus-gates-fahrstelle`. Der Slice liegt in `open/`: Er ist nicht priorisiert, den Anspruch trägt der Nutzer.

**Bezug:** [LH-QA-07](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur insoweit, als die Gegenproben dessen Messungen tragen; kein Anspruch des Produkts, der Slice ändert Rollen-Anweisungen und Wartungs-Doku des Harness. Bindung an das Register: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` und das mit `slice-harness-gegenproben-aus-gates` entstehende `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`. Keine aktive ADR berührt ([ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) trägt die Pflicht, der Slice setzt sie in den Anweisungen um).

**Berührte Spec-Stellen:** der Vertrag des Hooks für die Einlieferung [`SPEC-061`](../../../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md) Punkt 15 (Lesen des Ergebnisses; welcher Hash zählt) und Punkt 16 (1), (2), (3), (6) und (7) (Grenzen, die das Dokument nennt); der Vertrag der Sammelziele [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Punkte 3, 14 und 15 (7), nur gelesen; [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) Punkt 12 (Positivliste der inerten Pfade, nur gelesen). Der Slice ändert keinen Vertrag.

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

**Ziel:** Die Pflicht, die Gegenproben zu fahren, steht dort, wo die Rollen sie lesen, und die Fahrstelle ist beschrieben und einmal belegt: (1) Die Rollen-Anweisungen (`implement-slice`, `close-welle`, Reviewer- und Verifier-Typ, Implementer-Typ, Reviewer-Skill, `AGENTS.md` §6) nennen, wann welche Gegenprobe zu fahren ist: bei einer Änderung am geprüften Werkzeug oder Gate, an der Gegenprobe oder an deren Vertrag; am Ende jeder Welle; vor dem Release; Review und Verifikation fahren die betroffene Gegenprobe gezielt (Regel des engsten Sensors). Die Closure von Welle und Release prüft, dass für den Stand ein grünes Ergebnis des Gegenproben-Containers vorliegt (Hash des Commits, Regel nach `SPEC-061` Punkt 15); fehlt es, fährt die Closure `make gegenproben` selbst. (2) `docs/developer/devcontainer.md` beschreibt in einem Abschnitt den Gegenproben-Container (Name vorläufig `pgr-gp`; eigener Port, eigenes Volume, eigener Klon über sein Bare-Repo; Schritte aus §3 bis §6 und §9 des Dokuments) mit der Zeile zum Aufsetzen des Hooks aus `slice-harness-gegenproben-aus-gates-fahrstelle` (Kopie mit `install -m 755`, Marke), dem Einliefern eines Zweigs, dem Lesen des Ergebnisses und des Logs per `ssh`, dem Aufräumen von Hand und den Grenzen aus `SPEC-061` Punkt 16; `make gegenproben-host` läuft am Host. (3) Eine Einlieferung in `pgr-gp` hat `make gegenproben` ausgelöst und ein Ergebnis für den Hash abgelegt; sie ist die erste Fahrt im frischen Container und belegt oder widerlegt das Risiko 1 von `slice-harness-gegenproben-aus-gates-sammelziele`.

**`pgr-gp` ist ein spezieller Container, der nur die Gegenproben fährt** (Ergänzung des Nutzers vom 2026-10-18); in ihm wird kein Slice abgearbeitet. Implementer, Planner und Hook arbeiten am Host und in den Arbeits-Containern (`pgr-dev`, `pgr-dev2`, …) ungestört weiter, während `pgr-gp` im Hintergrund läuft. Die Arbeits-Container bekommen den Hook nicht (unten, *Ausdrücklich NICHT*).

**Zahlen, die der Abschnitt trägt, sind gemessen:** eine Fahrt von `make gegenproben` im Gegenproben-Container auf einem benannten Commit (Dauer, höchster Speicher des Containers, Last am Host, soweit der Implementer sie liest); was nicht gelaufen ist, steht als nicht gemessen (§13 des Dokuments).

**Sendung von `slice-harness-gegenproben-aus-gates-fahrstelle` (Geber des Hooks; eingetragen vom Planner im Commit des Schnitts, `AGENTS.md` §3.13):** (a) *Bindung:* Der Abschnitt, die Zeile zum Aufsetzen und die Regel an den Leser ruhen auf den Zusagen des Gebers: Name und Ort des Skripts `tools/harness/gegenproben-nach-einlieferung.sh`, Kopie mit `install -m 755`, Marke `<Bare-Repo>/gegenproben-ablage`, Pfade `gegenproben-<Name>-<Hash>` und Felder des Ergebnisses (`SPEC-061` Punkte 1, 2, 8, 9 und 15). Wer sie ändert, zieht den Abschnitt nach oder meldet es dem Planner. (b) *Die erste Fahrt:* Die Einlieferung, die den Hook an echtem `make gegenproben` zeigt, und die Messung der Dauer (`SPEC-061` Punkt 16 (1), Zeitlimit-Entscheidung als Empfehlung des Architect, nach Empfehlung entschieden) gehören diesem Slice (DoD 3). `grep` nach der Kennung des Gebers findet beide in §1, DoD und §6.

**Nachtrag zur Sendung des Gebers `slice-harness-gegenproben-aus-gates-fahrstelle`** (eingetragen vom Planner bei dessen Closure, 2026-10-18, `AGENTS.md` §3.13; der Geber liegt danach in `done/`): (c) *Bindung: Benutzer ohne root.* Die Gegenprobe des Hooks (`gegenproben-nach-einlieferung-gegenprobe`, steht in `GEGENPROBEN_CONTAINER`) endet als root mit Exit 2 und einer Meldung, ohne Überspringen (`SPEC-061` Kopf, Entscheidung des Architect; Grund: der Fall `p13-ablage-weg` macht die Ablage mit `chmod` unschreibbar, und root schließt `chmod` nicht aus). Der Gegenproben-Container läuft als `vscode`; der Abschnitt nennt `ssh` als `vscode` als Weg und sagt, dass `docker exec` ohne `-u vscode` root ergibt und `make gegenproben` dann rot endet. Der Hook selbst läuft als der Benutzer, der pusht (`SPEC-061` Punkt 5 und 13). (d) *Dauer der Fahrt.* `make gegenproben` im Container fährt seit dem Geber zusätzlich diese Gegenprobe (94 Fälle, 131 Mutationen; 5 min 30 s im Container am Stand `7da1b158` des geschlossenen Gebers, ohne Last und vor der Umstellung auf Fristen; mit künstlicher Last 9 min 08 s, am Host 6 min 10 s, vom Implementer des Gebers gemessen); die erste Fahrt misst die Dauer des ganzen Laufs (DoD 3), die Zahl des Gebers ist kein Beleg dafür. (e) *Die Gegenprobe des Hooks hat unter Last einmal ein Rot gezeigt, das der Geber nicht erklärt* (`p13-werkzeug-make`, einmalig in vier Lastläufen; Register: `BEO-REPO/lastabhaengige-zaehlung`). Färbt die erste Fahrt sie rot, nennt der Bericht den Fall und die Last des Hosts; das ist ein Auftreten im Register und kein Grund, den Hook zu ändern. Die drei Punkte stehen in diesem Slice in DoD 2 und 3 und in §6 (Randformen 10 und 11); sie legen keinen Liefer-Punkt und keine Schicht an (Nachzählung in §8).

**Sendung von `slice-harness-gegenproben-aus-gates-sammelziele` (eingetragen vom Planner bei dessen Closure, `AGENTS.md` §3.13; der Geber liegt in `done/`, die Adresse ist dieser Slice, früher geführt unter der Kennung `slice-harness-gegenproben-aus-gates-fahrstelle`):** Der Beleg für sein Risiko 1 (*Eine Gegenprobe ist ohne Kette nur grün, wenn sie von selbst hält*) fehlt: `make gegenproben` lief dort im Container `pgr-dev2` nach `make gates` im selben Container, nicht im frischen. Die erste Fahrt von `make gegenproben` im Gegenproben-Container dieses Slice ist dieser Beleg, sofern in dem Klon und dem Docker-Daemon des Containers vorher kein `make gates` lief und kein Image des Repos gebaut wurde; der Bericht nennt, was im Container vor der Fahrt lief (Ausgang des Risikos *Frischer Container* in §6, Nachzählung in §8). Wird sie rot, weil eine Gegenprobe von einem Zustand lebt, den ein anderes Gate herstellt, geht der Slice nach `open` (§4, Rückführung) und der Befund an den Planner.

**Sendung an `slice-erster-release-freigabe`** (eingetragen vom Planner im Commit des Schnitts, früher unter der Kennung des ungeteilten Slice, `AGENTS.md` §3.13): Die Freigabe-Checkliste des ersten Releases (`docs/maintainer/releasing.md` §10) verlangt als Punkt mit Beleg: ein grünes Ergebnis des Gegenproben-Containers für den Commit, der getaggt wird (fehlt es, `make gegenproben` im Gegenproben-Container), und `make gegenproben-host` am Host grün auf demselben Commit. Der Punkt gehört zu seinem dritten Liefer-Punkt (Checkliste); ein neuer Liefer-Punkt und eine neue Schicht entstehen nicht (Nachzählung in dessen §8; der Nehmer nennt die Kennung dieses Slice).

**Sendung an die Welle-Closure** (`close-welle`, Liefer-Punkt 1 dieses Slice): Sie prüft das Ergebnis des Gegenproben-Containers für den Stand und fährt sonst `make gegenproben`; die Sendung an `slice-erster-release-freigabe` trägt dasselbe für den Release (oben).

**Bindung an `slice-harness-devcontainer-einrichten` und `slice-harness-devcontainer-einrichten-start`** (eingetragen vom Planner im Commit des Schnitts, Entscheidung des Nutzers vom 2026-10-18, `AGENTS.md` §3.13): Der Abschnitt über den Gegenproben-Container verweist auf §3 bis §6 und §9 des Dokuments und auf die Limits (`PGR_CPUS`, `PGR_MEM`) und nennt deren Einzelbefehle. Wer §3 bis §5 (`slice-harness-devcontainer-einrichten-start`) oder §6 und §9 (`slice-harness-devcontainer-einrichten`) mit einem Ziel umschreibt, erhält den Abschnitt gültig oder meldet es dem Planner. Beide Nehmer nennen die Kennung dieses Slice in §1 und §6 (Nachzählung dort in §8).

**Übernimmt:** nichts.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Hook-Skript, seine Gegenprobe, Fragment und Gate-Index — `slice-harness-gegenproben-aus-gates-fahrstelle`; dieser Slice setzt dessen `done/` voraus (die Anweisung und der Abschnitt nennen ein Skript, das dann existiert). Ändert die erste Fahrt das Skript, geht der Slice nach `open` und der Befund an den Planner (§4).
- Der Kettenumbau, die Listen und die Namensregel, die erste Hälfte des Satzes in `AGENTS.md` §6 Schritt 6 — `slice-harness-gegenproben-aus-gates`; die Sammelziele, ihre Namen und Gegenprobe und die zweite Hälfte des Satzes — `slice-harness-gegenproben-aus-gates-sammelziele`.
- Ein Scheduler, eine CI-Pipeline, ein Warteschlangen-Dienst oder eine Benachrichtigung — Entscheidung des Nutzers: Der Container wird von Hand oder von einem Agenten gestartet; der Nachtlauf in der CI ist `slice-harness-ci-gegenproben`, er ergänzt die Pflicht, er ersetzt sie nicht.
- Der Hook in den Bare-Repos der Arbeits-Container (`pgr-dev`, `pgr-dev2`, …) — Schicht-Abgrenzung und Last: Der Arbeits-Container soll nicht ausgebremst werden; eine Fahrt von `make gegenproben` belegt Docker-Daemon und Kerne (§6, Ressourcen). Das Skript lehnt den Lauf ab, wo die Marke fehlt (`SPEC-061` Punkt 2); das Dokument sagt, dass die Zeile zum Aufsetzen nur im Gegenproben-Container steht.
- Ein Ziel, das den Gegenproben-Container anlegt, einliefert und startet (`make container-einrichten` und ein Ziel für die Fahrt) — `slice-harness-devcontainer-einrichten` automatisiert Bare-Repo, Einlieferung, Klon und Zweig; dieser Slice beschreibt das Aufsetzen mit den Einzelbefehlen und ruft kein Ziel, das es noch nicht gibt. Ein Ziel für das Aufsetzen selbst ist ein anderer Vorgang (Werkzeug statt Beschreibung); eine Adresse gibt es nicht, ein Slice legt der Planner an, wenn der Nutzer es entscheidet.
- `u-boot` und `devcontainer` im Gegenproben-Container — nicht gebraucht: Die Gegenproben, die sie brauchen, laufen am Host (`make gegenproben-host`, Entscheidung des Nutzers vom 2026-10-18); sie in das Image zu bringen verlangte einen Umbau des Gerüsts (die ersetzte ADR zum Überspringen im Dev-Container, Option C dort, nicht gewählt).
- Ein Stand oder Nachweis der Gegenproben, den der Stop-Hook am Host liest — `slice-harness-gegenproben-aus-gates` §1 und §6 (Frage 4 dort, Empfehlung des Planners: nein); `SPEC-060` Punkt 14. Das Ergebnis liegt im Container, nicht im Arbeitsbaum des Hosts; nur die Closure liest es, per `ssh`.
- Eine Messung der Last auf dem Host bei zwei Containern jenseits der einen Fahrt — Bestand bleibt ungemessen; der Abschnitt sagt es.
- Produkt-Code, Handbuch, `README.md` — Schicht-Abgrenzung: kein Verhalten des gebauten Binary.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Rollen-Anweisungen):** `.claude/commands/implement-slice.md`, `.claude/commands/close-welle.md`, `.claude/agents/implementer.md`, `.claude/agents/reviewer.md`, `.claude/agents/verifier.md`, `.harness/skills/reviewer.md` und `AGENTS.md` §6 nennen die Pflicht in einer Fassung: die Auslöser (Änderung am geprüften Werkzeug oder Gate, an der Gegenprobe, am Vertrag; Ende der Welle; Release), die Zuordnung Gegenprobe zu geprüftem Werkzeug über `harness/README.md` §Sensors, die gezielte Fahrt durch Review und Verifikation, und dass die Fahrt der eingelieferte, committete Stand ist; die Closure von Welle und Release prüft, dass für den Stand ein grünes Ergebnis des Gegenproben-Containers vorliegt, fährt sonst `make gegenproben` selbst und trägt Ergebnis und Hash (beide Hashes, wenn nach dem Lauf nur inerte Pfade änderten, `SPEC-052` Punkt 12) in die Closure-Notiz, dazu die Fahrt von `make gegenproben-host`. Die Anweisungen sagen die Regel an den Schreiber und zeigen auf `harness/README.md` §Sensors als Quelle, sie zählen kein Verhalten der Gates auf (`AGENTS.md` §3.11).
- [ ] **Liefer-Punkt 2 (Fahrstelle im Dokument):** `docs/developer/devcontainer.md` führt einen Abschnitt (Ort §10 *Mehrere Container nebeneinander* oder ein neuer Abschnitt dahinter, dazu eine Zeile in §8; der Implementer wählt) über den Gegenproben-Container mit dem Ablauf: Anlegen wie §3 bis §6 mit eigenen Werten (`PGR_NAME`, `PGR_SSH_PORT`, `PGR_CPUS`, `PGR_MEM`), die Zeile zum Aufsetzen des Hooks (Marke und Kopie des Skripts im Bare-Repo, mit dem Hinweis, sie nach jeder Änderung des Skripts zu wiederholen), Einliefern eines Zweigs (§9), Lesen des Ergebnisses für einen Hash und des Logs per `ssh` (Muster `gegenproben-*-<Hash>`, Regel nach `SPEC-061` Punkt 15), die Probe, ob der Hook läuft (Testeinlieferung und Ergebnisdatei), Aufräumen von Hand (verwaiste Marke, Log, Rest des Arbeitsbaums, Container und Netze der Einzelziele nach `SPEC-061` Punkt 16 (2) und (3)); dass die Arbeits-Container den Hook nicht bekommen; dass die Fahrt als `vscode` läuft (Sendung (c) des Gebers `slice-harness-gegenproben-aus-gates-fahrstelle`: `ssh` als `vscode`, `docker exec` nur mit `-u vscode`). Jede Aussage über Dauer, Speicher und Last steht mit „gemessen“ und Datum oder mit „nicht gemessen“; jeder Befehl des Abschnitts ist gelaufen (`AGENTS.md` §3.11, Beispiele gelten erst als belegt, wenn sie gelaufen sind).
- [ ] **Liefer-Punkt 3 (erste Fahrt als Beleg):** Eine Einlieferung eines Zweigs in den Gegenproben-Container hat `make gegenproben` ausgelöst und ein Ergebnis für den Hash abgelegt (Ergebnis, Dauer und höchster Speicher stehen im Bericht und im Abschnitt); `make gegenproben-host` am Host auf demselben Commit ebenso. **Die erste Fahrt ist die im frischen Container** (Sendung von `slice-harness-gegenproben-aus-gates-sammelziele`): Vor ihr lief im Klon und im Docker-Daemon des Containers kein `make gates`; der Bericht nennt, was vorher lief, und das Ergebnis der Fahrt belegt oder widerlegt das Risiko 1 jenes Slice. Der Zweig der Fahrt trägt den Stand, an dem der Hook (Kopie) und der Abschnitt geprüft sind. Der Bericht nennt den Benutzer der Fahrt, die Dauer des ganzen Laufs mit der Gegenprobe des Hooks darin (Sendung (d) des Gebers `slice-harness-gegenproben-aus-gates-fahrstelle`) und, falls sie rot wird, den Fall und die Last (Sendung (e)).
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `.claude/` ändert).
- [ ] Doku-Update: `harness/README.md` §Sensors nennt in der Zeile der Sammelziele die Fahrstelle (Verweis auf den Abschnitt im Dokument); `docs/developer/devcontainer.md` §8 trägt die Zeile des neuen Abschnitts in der Tabelle der Messungen.
- [ ] Beobachtungs-Register: `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt` mit dem Beleg dieses Slice geführt (Mitigation getragen: ja oder nein); kein Zähler wird gesetzt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review. Das Review fährt die Stichprobe nach der neuen Regel selbst.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen; die Folge-Slice-Paarung prüft Existenz **und** Annahme der Sendungen an `slice-erster-release-freigabe`, `slice-harness-devcontainer-einrichten` und `slice-harness-devcontainer-einrichten-start`: `grep` nach der Kennung dieses Slice findet sie in deren §1 und §6 (`AGENTS.md` §3.13).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/commands/implement-slice.md`, `.claude/commands/close-welle.md` | update | Pflicht im Workflow (Schritt „Gates“ und Sensor-Belege) und in der Welle-Closure (Fahrten beider Sammelziele als Beleg, Ergebnis des Containers nach `SPEC-061` Punkt 15) |
| `.claude/agents/implementer.md`, `reviewer.md`, `verifier.md`, `.harness/skills/reviewer.md` | update | Stichprobe nach der Regel des engsten Sensors; Reviewer und Verifier fahren die betroffene Gegenprobe gezielt |
| `AGENTS.md` | update | §6 Schritt 6 (der Satz aus `slice-harness-gegenproben-aus-gates` und `slice-harness-gegenproben-aus-gates-sammelziele` wird um die Auslöser ergänzt, falls dort kürzer) |
| `docs/developer/devcontainer.md` | update | Abschnitt über den Gegenproben-Container, Zeile in §8, §13 für Ungeprüftes |
| `harness/README.md` | update | Fahrstelle in der Zeile der Sammelziele |
| `docs/plan/planning/next/slice-erster-release-freigabe.md`, `docs/plan/planning/open/slice-harness-devcontainer-einrichten.md`, `docs/plan/planning/done/slice-harness-devcontainer-einrichten-start.md` | update (Planner, im Commit des Schnitts) | Sendung der Checkliste und Bindung an den Abschnitt (§1, §6, §8 der Nehmer) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. **Vorbedingung:** `slice-harness-gegenproben-aus-gates-fahrstelle`, `slice-harness-gegenproben-aus-gates` und `slice-harness-gegenproben-aus-gates-sammelziele` liegen in `done/` (der Hook, die Listen, die Ziele und ihre Namen stehen); der Dev-Container-Ablauf von `slice-harness-devcontainer-einrichten-geruest` liegt in `done/` (der Abschnitt baut auf `make container-start`, soweit es dann steht; sonst auf §3 bis §6 des Dokuments). Reihenfolge: nach `slice-harness-gegenproben-aus-gates-fahrstelle`; der Slice wartet nicht auf `welle-harness-regeln` (Drift-Log der Roadmap). `Verantwortlich:` trägt beim Start Person und Zweig; er läuft auf eigenem Zweig. Die Fahrt braucht einen laufenden Gegenproben-Container; ihn anzulegen ist Teil dieses Slice (Einzelbefehle des Dokuments).

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice>` über die Pfade am Ende dieses Absatzes. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt. **Pfade:** `-- AGENTS.md .claude .harness/skills docs/developer harness/README.md`.

**Schätzung (Architect, 2026-10-18, nach `SPEC-061`, auf diesen Slice verteilt vom Planner; Metrik wie oben):**

| Teil | Zeilen |
|---|---|
| Rollen-Anweisungen, sieben Dateien, je 6 bis 20 Zeilen | 50 bis 130 |
| Abschnitt im Dokument mit Befehlen, Zeile zum Aufsetzen und Messung | 80 bis 160 |
| `harness/README.md` (Fahrstelle in der Zeile der Sammelziele), `AGENTS.md` | 10 bis 40 |
| **Zusammen** | **140 bis 330** |

Meldeschwelle 700, **Marke 900**; die Schätzung liegt weit darunter. Die Spanne hängt an der Länge des Abschnitts und daran, wie viele Zeilen der Anweisungen die Auslöser wiederholen statt auf `harness/README.md` §Sensors zu zeigen.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder die Fahrt verlangt mehr als eine Einlieferung zur Klärung (Dauer, Speicher, Last) und der Abschnitt wächst über 160 Zeilen.
- `in-progress` → `open` (blockiert): Die Fahrt von `make gegenproben` ist im Gegenproben-Container nicht grün (Befund an den Planner, Rückgabe an `slice-harness-gegenproben-aus-gates-sammelziele`, Risiko *Zustand eines anderen Gates* dort), oder der Hook aus `slice-harness-gegenproben-aus-gates-fahrstelle` legt in der echten Einlieferung kein Ergebnis ab (Befund an den Planner, Rückgabe an den Architect).

## 5. Closure-Trigger

DoD vollständig; eine Einlieferung in den Gegenproben-Container hat ein Ergebnis für den Hash abgelegt, `make gegenproben-host` auf demselben Commit am Host belegt; `make gates` grün; Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag; die drei Paarungen getragen, darunter die Annahme der Sendungen durch die drei Nehmer.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Der Slice liefert keinen neuen Vertrag im Sinne von `AGENTS.md` §3.10 (kein Code); die Randformen der Fahrstelle sind Aussagen des Dokuments und werden dort belegt oder als nicht gemessen geführt (§3.11). Sie stehen hier als Randformen (§3.12), jede mit Entscheidung:

1. **Eingelieferter Stand gegen Arbeitsbaum.** Die Fahrt prüft den committeten Stand, den der Container im Bare-Repo hat, nicht den Arbeitsbaum des Hosts (`devcontainer.md` §9: uncommittete und ignorierte Dateien kommen nicht mit). Entscheidung (Empfehlung): Der Bericht nennt den Commit; eine Fahrt auf einem anderen Commit als dem getaggten oder zu schließenden gilt nicht. Entschieden: offen (Architect/Implementer im Dokument).
2. **Zwei Container desselben Klons** (Arbeitscontainer und Gegenproben-Container teilen ihren Klon nicht; zwei Fahrten im selben Klon gleichzeitig sind nicht zugesagt, `SPEC-052` Punkt 8). Entscheidung: eine Fahrt je Container; der Abschnitt sagt es.
3. **Ressourcen.** Je Container 8 CPUs und die Vorgabe `PGR_MEM` (`12g` im Stand des Dokuments, `8g` nach der Sendung `arch-limits` in `slice-harness-devcontainer-einrichten-start`), Host 20 CPUs und 31 GiB. Gemessen ist (Dokument §8) `make gates` mit 1701 s und höchstens 5,42 GiB; für die Gegenproben ist nichts gemessen. Empfehlung: Standardwerte, `PGR_CPUS` und `PGR_MEM` des Gegenproben-Containers senken, wenn die Messung des Implementers Last am Host zeigt; der Abschnitt nennt die Messung oder „nicht gemessen“.
4. **Dev-Container-Erkennung nach der ersetzten ADR zum Überspringen entfällt** ([ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md), Accepted) (Entscheidung des Nutzers vom 2026-10-18): Alle Gegenproben in der Container-Liste laufen im Gegenproben-Container ohne Überspringen; was `u-boot` und `devcontainer` braucht, läuft nur am Host (`make gegenproben-host`). Vorschlag zur Frage „u-boot im Gegenproben-Container?“: nein, weil die Aufteilung den Bedarf aufhebt und das Image sonst vom Gerüst von `u-boot` abwiche. Entschieden: Nutzer (Aufteilung); der Vorschlag zum Image offen.
5. **Hintergrundlauf.** Der Hook startet die Fahrt mit `setsid` und Log in der Ablage, sie überlebt das Ende der `ssh`-Sitzung; ein Neustart des Containers beendet sie (§8 des Dokuments: Prozesse beginnen neu). Das Verhalten ist Zusage des Hooks (`SPEC-061` Punkte 5 und 9, Geber `slice-harness-gegenproben-aus-gates-fahrstelle`); das Dokument sagt, was nach einem Neustart zu tun ist (verwaiste Marke, Fahrt wiederholen). Entschieden: Architect im Vertrag; Aussage im Dokument offen.
6. **Fahrt ohne Netz.** Der Daemon im Container lädt die Images selbst (Egress eingeschränkt, `egress.allow`); fehlt ein Host der Liste, endet die Fahrt rot mit der Meldung der Gegenprobe, nicht still. Entschieden: Bestand (§11 des Dokuments).
7. **Wer die Pflicht auslöst.** Das Werkzeug zur Zuordnung „geändertes Werkzeug → Gegenprobe“ gibt es nicht (eine Tabelle wäre die Eingabeliste, die der Nutzer verworfen hat); der Schreiber der Anweisung nennt die Quelle, `harness/README.md` §Sensors. Akzeptiertes Negativ; Grund: ein Sensor dafür wäre der erste Zuschnitt.
8. **Welcher Hash zählt der Closure** (H12 des ungeteilten Slice). Der Commit, der den letzten nicht inerten Stand trägt (Positivliste von `SPEC-052` Punkt 12); ein Commit, der nur inerte Pfade ändert, braucht keinen neuen Lauf, die Closure-Notiz nennt beide Hashes. Fehlt ein Ergebnis oder ist es rot, fährt die Closure `make gegenproben` selbst. Kein Sensor dafür (Urteil von Closure und Review). Entschieden: Architect, `SPEC-061` Punkt 15 (Regel an den Leser; Lesen per `ssh` mit dem Muster `gegenproben-*-<Hash>`: **Empfehlung des Architect, nach Empfehlung entschieden**, nicht ausdrücklich vom Nutzer bestätigt).
9. **Mehrere Zweige in einem Push** fahren nacheinander, die Dauer ist die Summe (`SPEC-061` Punkt 16 (4); **Empfehlung des Architect, nach Empfehlung entschieden**, nicht ausdrücklich vom Nutzer bestätigt). Das Dokument nennt es; die Fahrt belegt einen Zweig.
10. **Benutzer der Fahrt** (Sendung (c) des Gebers `slice-harness-gegenproben-aus-gates-fahrstelle`). Die Gegenprobe des Hooks verlangt einen Benutzer ohne root, der Container läuft als `vscode`, ein `docker exec` ohne `-u` ergibt root. Entschieden: Architect, `SPEC-061` Kopf (als root Exit 2, ohne Überspringen); die Aussage im Dokument ist offen (Implementer): `ssh` als `vscode`, `docker exec -u vscode`.
11. **Dauer und Last der Gegenprobe des Hooks** (Sendungen (d) und (e)). Gemessen hat sie der Geber nur für sich (siehe §1); die Dauer des ganzen Laufs und ein Rot unter Last sind Gegenstand der ersten Fahrt. Entschieden: kein Zeitlimit (`SPEC-061` Punkt 16 (1)); der Bericht nennt die Zahl.

**Aufrufarten des Treibers (§3.12):** nicht Gegenstand dieses Slice (der Hook ist `slice-harness-gegenproben-aus-gates-fahrstelle`, `SPEC-061` Punkte 10, 11 und 14). **Aufruf-Voraussetzung:** Zusagen des Abschnitts gelten für Fahrten nacheinander über die Sperre des Hooks; zwei Fahrten im selben Klon gleichzeitig sind nicht zugesagt (Randform 2).

Risiken:

- **Die Pflicht wird nicht gelesen oder nicht befolgt, und die Regression bleibt unbemerkt** (Risiko 1 von `slice-harness-gegenproben-aus-gates`). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`.
- **Die Beschreibung der Fahrstelle veraltet**, sobald `slice-harness-devcontainer-einrichten-start` §3 bis §5 und `slice-harness-devcontainer-einrichten` §6 und §9 des Dokuments mit Zielen umschreibt, auf die dieser Abschnitt verweist. — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`; die Bindung steht als Sendung in §1 und §6 beider Nehmer (eingetragen vom Planner im Commit des Schnitts, Nachzählung dort in §8): Wer die Abschnitte umschreibt, zieht den Abschnitt nach oder meldet es dem Planner.
- **Der Hook läuft bei jeder Einlieferung, auch bei unfertigen Zweigen, und belegt den Docker-Daemon des Gegenproben-Containers lange** (28 min 21 s für `make gates` im Container sind gemessen, die Gegenproben nicht). Die Läufe stauen sich (Sperre des Hooks, `SPEC-061` Punkt 6). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, falls die Messung des Implementers Staus zeigt; sonst entfallen mit der Zahl.
- **Das Ergebnis liegt im Container und der Container ist weg** (Volume gelöscht, Neustart, `devcontainer.md` §4: ein Neustart behält das Volume). — **Ausgang (vorgesehen):** entfallen, weil die Closure bei fehlendem Ergebnis selbst fährt (Randform 8); das Ergebnis ist ein Nachweis zum Nachlesen, kein Stand.
- **Eine Fahrt kostet Zeit und Platz im Container** (28 min 21 s für `make gates` im Container sind gemessen, die Gegenproben nicht; ein zweites Volume mit eigenem Docker-Cache, `devcontainer.md` §10 und §11). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, falls die Messung des Implementers die Dauer als Zumutung zeigt; sonst entfallen mit der Zahl. Ein Zeitlimit des Laufs gibt es nicht (`SPEC-061` Punkt 16 (1)); die gemessene Dauer ist die Grundlage, falls später eines entschieden wird.
- **Frischer Container** (Sendung von `slice-harness-gegenproben-aus-gates-sammelziele`): Eine Gegenprobe der Container-Liste ist nur grün, wenn vorher `make gates` oder ein anderer Lauf im selben Container einen Zustand herstellte (Image, Cache, Datei). — **Ausgang (vorgesehen):** entfallen, wenn die erste Fahrt im Gegenproben-Container ohne vorheriges `make gates` grün endet (Beleg im Bericht); sonst eingetreten → Rückführung nach `open` (§4).
- **Die Sendung an `slice-erster-release-freigabe` nimmt der Nehmer nicht an**, weil er vor diesem Slice schließt. — **Ausgang (vorgesehen):** eingetreten → die Checkliste trägt den Punkt nachträglich, Folge-Slice mit Kennung vom Planner; die Reihenfolge der Wellen spricht dagegen (welle-erster-release liegt vor den Harness-Wellen) — der Planner prüft sie beim Start dieses Slice und trägt das Ergebnis in §7 ein.
- **Der Hook (Kopie im Bare-Repo) steht beim Einliefern nicht auf dem Stand des Skripts im Repo** (`SPEC-061` Punkt 1: kein Sensor). Die Fahrt prüft dann ein altes Skript. — **Ausgang (vorgesehen):** entfallen, wenn der Bericht den Hash der Kopie (`sha256sum`) neben dem des Skripts im Repo nennt und beide gleich sind; sonst eingetreten → die Zeile zum Aufsetzen wird wiederholt.

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

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Hintergrund), `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt` (Gegenmittel; zweites Auftreten bei `slice-harness-gegenproben-aus-gates-sammelziele`, mit diesem Slice kein weiteres), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Dokument und Anweisungen sagen nur, was gelaufen ist; `AGENTS.md` §3.11), `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (die Fahrt im eigenen Container ist die Antwort darauf für die Gegenproben, nicht Gegenstand), `BEO-REPO/plan-folgt-korrektur-nicht` (die Sendungen stehen in den Plänen der Nehmer). Keiner erreicht mit diesem Slice 3× neu.

alle berührten Sub-Areas GF (`REPO`).

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-18, im Commit des Schnitts):** Liefer-Punkte 3 (Rollen-Anweisungen mit Prüfung des Ergebnisses in der Closure | Fahrstelle im Dokument | erste Fahrt als Beleg und Messung). Schichten 2: Harness-Artefakte (Rollen-Anweisungen, `AGENTS.md`, `harness/README.md` als Gate-Index) und Wartungs-Doku (`docs/developer/`, nach der Teilung von `slice-harness-d-check-v0-85` §8 und `AGENTS.md` §3.13 eine Schicht); Planung und Vertrag sind keine. Schätzung 140 bis 330 in der Metrik von §4, Marke 900, Meldeschwelle 700 (weit darunter). Der Plan trägt im Commit des Schnitts ein: die Bindung und die Sendung des Gebers `slice-harness-gegenproben-aus-gates-fahrstelle` (Abschnitt und erste Fahrt: DoD 2 und 3; keine neue Zahl), die Sendung von `slice-harness-gegenproben-aus-gates-sammelziele` (Bedingung an die Fahrt in DoD 3 und Risiko in §6, kein Liefer-Punkt, keine Schicht), die Sendung an `slice-erster-release-freigabe` (Nachzählung dort in §8), die Sendung an die Welle-Closure (Liefer-Punkt 1) und die Bindung an `slice-harness-devcontainer-einrichten` und `slice-harness-devcontainer-einrichten-start` (Nachzählung dort in §8).

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit der Closure von `slice-harness-gegenproben-aus-gates-fahrstelle`): Der Nachtrag (c) bis (e) des Gebers `slice-harness-gegenproben-aus-gates-fahrstelle` steht in DoD 2 und 3 und in §6 (Randformen 10 und 11) als je ein Satz; Liefer-Punkte bleiben 3 (Rollen-Anweisungen · Fahrstelle im Dokument · erste Fahrt), Schichten bleiben 2 (Harness-Artefakte · Wartungs-Doku). Die Schätzung 140 bis 330 in der Metrik von §4 wächst um höchstens 5 Zeilen (zwei Sätze im Abschnitt des Dokuments), Marke 900 und Meldeschwelle 700 bleiben weit entfernt.

**Herkunft der Zählung:** Der ungeteilte Slice `slice-harness-gegenproben-aus-gates-fahrstelle` hatte 3 Liefer-Punkte (Rollen-Anweisungen · Hook-Skript mit Gegenprobe · Fahrstelle im Dokument mit einer Einlieferung als Beleg) und 2 Schichten, Schätzung 740 bis 1345. Der Schnitt gibt das Hook-Skript an den Geber ab; dieser Slice behält Rollen-Anweisungen, Dokument und Fahrt.
