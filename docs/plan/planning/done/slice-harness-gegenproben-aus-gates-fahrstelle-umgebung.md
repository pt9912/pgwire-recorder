# Slice slice-harness-gegenproben-aus-gates-fahrstelle-umgebung: Der Lauf des Hooks erbt keine Deskriptoren über 2 und ignoriert SIGINT nicht; make und seine Nachkommen ignorieren weder SIGINT noch SIGQUIT

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten, vierte Stelle der Reihenfolge, direkt hinter `slice-harness-gegenproben-aus-gates-fahrstelle` und vor `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (Planner, 2026-10-18). Wie die vier Slices davor aus der Reihenfolge der Wellen gelöst: Start-Trigger ist `slice-harness-gegenproben-aus-gates-fahrstelle` in `done/` (erfüllt). **Entstanden am 2026-10-18 aus der ersten echten Fahrt des Hooks** im Gegenproben-Container (Arbeitsort `pgr-dev2`): zwei Fehler, die die Gegenprobe des Gebers nicht fand, weil sie mit Attrappen und einem lokalen Push fährt (siehe §1).

**Bezug:** [LH-QA-07](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur insoweit, als die Gegenproben dessen Messungen tragen; kein Anspruch des Produkts, der Slice ändert ein Harness-Werkzeug. Bindung an das Register: `BEO-REPO/hook-erbt-deskriptor-und-ignorierte-signale` (neu), `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (fünfter Beleg), `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt`. Keine aktive ADR berührt ([ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) trägt die Pflicht der Fahrt, der Slice ändert den Hook, nicht die Entscheidung).

**Berührte Spec-Stellen:** der Vertrag des Hooks für die Einlieferung [`SPEC-061`](../../../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md) Punkt 5 (erweitert vom Architect im Zweig `arch-fd`, nachgezogen nach dem Review in den Punkten 5, 11 und 13: keine geerbten Deskriptoren über 2 im Lauf, `SIGINT` im Lauf und `SIGINT` und `SIGQUIT` in `make` und seinen Nachkommen nicht ignoriert; `SIGQUIT` in der Shell des Laufs selbst ist ausgenommen, bash ignoriert es selbst; Plattform bash ≥ 5.1, Entscheidung des Nutzers vom 2026-10-18). Der Slice ändert den Vertrag nicht; der Architect trägt ihn vor dem ersten Code-Commit.

**Verantwortlich:** pt9912 (Zweig slice-harness-gegenproben-aus-gates-fahrstelle-umgebung)

**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Lauf, den der Hook `tools/harness/gegenproben-nach-einlieferung.sh` im Hintergrund startet, hat (1) keinen Deskriptor über 2, den er vom Hook oder von `git receive-pack` erbt, und (2) `SIGINT` im Lauf sowie `SIGINT` und `SIGQUIT` in `make` und seinen Nachkommen nicht auf „ignoriert“ (`SIGQUIT` in der Shell des Laufs selbst nimmt `SPEC-061` Punkt 11 aus: bash ignoriert es selbst); die Gegenprobe des Hooks zeigt beides mit je einem Fall und einer Mutation, die nur ihn bricht. Hook und Gegenprobe setzen bash ≥ 5.1 voraus (Entscheidung des Nutzers vom 2026-10-18); die Gegenprobe prüft die Version vorab. Zwei Liefer-Punkte: **(L1)** das Skript (`setsid -f` statt `&`, alle Deskriptoren über 2 im Kind schließen), **(L2)** die Gegenprobe (Fälle mit einem zusätzlichen geerbten Deskriptor an einer Pipe und mit der Prüfung von `SigIgn` in `/proc/<pid>/status` von Lauf und `make`; Mutationen).

**Anlass (erste Fahrt, 2026-10-18):** (1) Der `ssh`-Push hing rund 70 s: Der Hintergrundlauf (`setsid … &`) von `git receive-pack` erbt einen weiteren Deskriptor (`/proc/<pid>/fd`: `5 -> pipe`), den nur der Lauf, `make` und `gegenproben.sh` tragen; `git` wartet auf das Ende der Leitung. Die Gegenprobe des Gebers (`SPEC-061` Punkt 5, Fälle `p5-fd-0` bis `p5-fd-2`, `p5-leitung-*`) kennt nur die Deskriptoren 0 bis 2. (2) Ein Lauf mit `&` in einer Shell ohne Jobsteuerung bekommt `SIGINT` auf „ignoriert“; `container-gegenprobe` wird rot (`p9-binINT-signal …`), weil ein ignoriertes Signal in der Kindkette bleibt. Beide sind Auftreten der Klasse `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (Ort der Ausführung, hier der Aufrufer `git receive-pack` statt des Test-Treibers) und von `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt` (Regression nach `done/`, bemerkt erst durch die Fahrt).

**Sendung des Gebers `slice-harness-gegenproben-aus-gates-fahrstelle` (Hook, Skript und Gegenprobe; `AGENTS.md` §3.13):** Dieser Slice setzt dessen `done/` voraus und ändert dessen Skript und Gegenprobe. Der Geber nimmt nichts mehr an (er liegt in `done/`); die Adresse ist dieser Slice, den der Planner dafür anlegt.

**Übernimmt:** nichts.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Änderung am Vertrag `SPEC-061` — Architect, vor dem Code (§3.12); der Slice liefert, was Punkt 5 danach zusagt, und ändert ihn nicht.
- Pflicht, Fahrstelle im Dokument und erste Fahrt — `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`; er setzt diesen Slice in `done/` voraus (Start-Trigger, im Commit des Anlegens dort eingetragen), weil die erste Fahrt sonst am selben Fehler hinge. Beleg, dass der `ssh`-Push nicht mehr hängt, ist dort die Fahrt im frischen Container (DoD 3 dort).
- Andere Randformen des Hooks (Sperre, Ablage, Ergebnis, Aufrufarten von `make`) — Bestand des Gebers, geprüft und in `done/`; die Fahrt hat sie nicht gebrochen.
- Ein Sensor, der die Umgebung des Aufrufers (Deskriptoren, Signale) für jede Gegenprobe prüft — wäre der erste Zuschnitt eines neuen Gates (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`); der Slice trägt nur die Fälle des Hooks. Ob die Klasse in `AGENTS.md` §3.12 eine Regel bekommt, entscheidet der Nutzer beim Lese-Schritt der Welle-Closure (Register: `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`, Ausgang offen).
- Produkt-Code, Handbuch, `README.md` — Schicht-Abgrenzung: kein Verhalten des gebauten Binary.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Liefer-Punkt 1 (Skript):** `tools/harness/gegenproben-nach-einlieferung.sh` startet den Lauf mit `setsid -f` statt `setsid … &`; der Lauf hat keinen Deskriptor über 2 (alle im Kind geschlossen, der Weg ist der des Architect in `SPEC-061` Punkt 5), und `SIGINT` ist im Lauf, `SIGINT` und `SIGQUIT` sind in `make` und seinen Nachkommen nicht ignoriert (`SIGQUIT` in der Shell des Laufs selbst ausgenommen, `SPEC-061` Punkt 11). Verhalten nach `SPEC-061` Punkt 5 in der Fassung des Architect.
- [x] **Liefer-Punkt 2 (Gegenprobe):** `tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh` führt Fälle für (a) einen zusätzlich geerbten Deskriptor an einer Pipe, deren Leser bis zum Ende der Leitung wartet (der Leser endet vor dem Ende des Laufs; Deskriptor 3 und ein höherer je ein Fall), (b) `SigIgn` in `/proc/<pid>/status` des Laufs (`SIGINT`) und von `make`, Kind und Enkel (`SIGINT` und `SIGQUIT` je ein Fall; der Hook wird mit beiden auf „ignoriert“ aufgerufen) und je eine Mutation, die nur ihren Fall färbt (`&` statt `setsid -f`, Deskriptoren bleiben offen, Signale bleiben ignoriert; je Signal und je Deskriptor eine, `AGENTS.md` §3.10); die unveränderte Kopie bleibt grün. Der Bericht nennt je Zusage: Zusage · Mutation · roter Test.
- [x] `make gates` grün (Stop-Hook; der Slice ändert `tools/`), belegt in §7 *Läufe* am Stand `07e45e1e` (Exit 0); danach ändert `3d8ba0a0` unter `tools/` nur einen Kommentar; der Lauf am Integrationsstand folgt auf dem Host. Die Gegenprobe des Hooks ist im Gegenproben-Container gefahren oder im Bericht als nicht gefahren benannt.
- [x] Doku-Update: `harness/README.md` §Sensors (Zeile der Gegenprobe) nennt die neuen Fälle nur, soweit die Zeile Aufzählungen trägt; sonst unverändert.
- [x] Beobachtungs-Register: `BEO-REPO/hook-erbt-deskriptor-und-ignorierte-signale` und `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` mit dem Beleg dieses Slice geführt; kein Zähler wird gesetzt.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen; die Folge-Slice-Paarung prüft Existenz **und** Annahme der Sendung an `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`: `grep` nach der Kennung dieses Slice findet sie in dessen §1, §4 und §6 (`AGENTS.md` §3.13).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/gegenproben-nach-einlieferung.sh` | update | Start des Laufs mit `setsid -f`; Deskriptoren über 2 im Kind schließen; Signale nicht ignoriert (L1) |
| `tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh` | update | Fälle und Mutationen nach `SPEC-061` Punkt 5 (L2) |
| `harness/README.md` | update, nur falls die Zeile eine Aufzählung trägt | Gate-Index |

Plattform: Hook und Gegenprobe setzen bash ≥ 5.1 voraus (Entscheidung des Nutzers vom 2026-10-18; die Gegenprobe prüft die Version vorab, Exit 2), der Hook `env` mit `--default-signal` (coreutils ≥ 8.31, `SPEC-061` Punkt 13) und `/proc` (nur Linux, Kopf des Vertrags).

Ansatz: Der Architect legt in `SPEC-061` Punkt 5 den Weg zum Schließen der Deskriptoren fest (etwa eine Schleife über `/proc/self/fd` statt über Zahlen; Randform in §6, „Wo entschieden“ trägt der Architect nach). Die Gegenprobe ruft den Hook mit einem zusätzlich geöffneten Deskriptor an einer Pipe und mit `trap '' INT QUIT` in der aufrufenden Shell, liest `/proc/<pid>/fd` und `SigIgn` des Laufs und der Attrappe von `make` (die Attrappen des Gebers schreiben, was sie sehen, in die Ablage). Der Implementer weist die Fälle am Stand vor der Änderung als rot nach (Bewusstes Brechen, Modul 11).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-gegenproben-aus-gates-fahrstelle` liegt in `done/` (erfüllt). Der Architect hat `SPEC-061` Punkt 5 (keine geerbten Deskriptoren über 2, `SIGINT` und `SIGQUIT` nicht ignoriert, Fälle) im Vertrag festgehalten und die Randformen aus §6 mit „Wo entschieden“ versehen (Zweig `arch-fd`). Den Anspruch trägt der Nutzer; `Verantwortlich:` setzt `open` → `next`.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice>` über die Pfade am Ende dieses Absatzes, Tests eingeschlossen. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt: `4219bab6`. **Pfade:** `-- tools harness/mk harness/README.md`; `harness/vertraege` liegt nicht darin (der Architect schreibt dort).

**Schätzung (Planner, Metrik wie oben; Gegenprobe mit dem Doppelten des Produktcodes, `AGENTS.md` §3.14):**

| Teil | Zeilen |
|---|---|
| Skript | 10 bis 25 |
| Gegenprobe (Fälle und Mutationen) | 80 bis 150 |
| `harness/README.md` | 0 bis 5 |
| **Zusammen** | **90 bis 180** |

Meldeschwelle 700, **Marke 900**; die Schätzung liegt weit darunter.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder die Gegenprobe braucht mehr als die Deskriptor- und Signal-Fälle mit ihren Mutationen, weil der Architect weitere Randformen des Aufrufers festlegt (dann Schnitt durch den Planner).
- `in-progress` → `open` (blockiert): Der Weg zum Schließen der Deskriptoren trägt im Container nicht (etwa kein `/proc`), oder `setsid -f` ändert eine Zusage von `SPEC-061` Punkt 5 oder 9, die der Architect nicht vorsah (Befund an den Architect, Rückgabe an den Planner).

## 5. Closure-Trigger

DoD vollständig; `make gates` grün; Review- und Verifikations-Report liegen vor; die Gegenprobe des Hooks ist an diesem Stand grün und die Mutationen aus DoD 2 sind vom Implementer rot gesehen; Closure-Notiz mit Lerneintrag; die drei Paarungen getragen, darunter die Annahme der Sendung durch `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen (§3.12), jede mit der Stelle, an der sie entschieden ist** (die Spalte „Wo entschieden“ trägt der Architect nach, vor dem ersten Code-Commit; was dort offen bleibt, entscheidet der Implementer nicht, er gibt es zurück):

| Randform | Wo entschieden |
|---|---|
| Deskriptor 3, ein höherer Deskriptor (etwa 5 wie in der Fahrt) und mehrere zugleich; Deskriptor an einer Pipe, einer Datei, einem Socket | Architect, `SPEC-061` Punkt 5 (alle Deskriptoren über 2 vor dem Lauf geschlossen, ohne Nummernliste; Mutation „schließt erst nach dem Start von `make`“ entfällt mit Grund) |
| Deskriptor, den der Hook selbst für die Sperre öffnet (`flock`): im Lauf nötig oder geschlossen | Architect, `SPEC-061` Punkt 5 (Sperre nur im Lauf, nicht in `make` und Nachkommen) |
| Obergrenze der Deskriptoren (`ulimit -n` sehr hoch): Schleife über `/proc/self/fd` statt über Zahlen | Architect, `SPEC-061` Punkt 5 (Schleife über `/proc/self/fd`) |
| `SIGINT` und `SIGQUIT` ignoriert beim Aufruf (`trap '' INT QUIT`, Aufruf mit `&` in einer Shell ohne Jobsteuerung); `SIGHUP`, `SIGPIPE`, `SIGTERM` | Architect, `SPEC-061` Punkt 11: der Lauf `SIGINT` nicht ignoriert, `make` und Nachkommen `SIGINT` und `SIGQUIT`; `SIGQUIT` im Lauf selbst Grenze (bash ignoriert es selbst); der Hook setzt beide für den Lauf zurück, wenn der Aufrufer sie ignoriert (`env --default-signal=INT,QUIT`, Fälle `p11-erbt-*`); `SIGHUP`, `SIGPIPE` nicht zugesagt, `SIGTERM` Punkt 9 |
| Wo `SigIgn` gelesen wird: Lauf, `make`, Kinder von `make` | Architect, `SPEC-061` Punkt 11 (Lauf, `make`, Kind und Enkel; Fälle `p11-sigint-enkel`, `p11-sigquit-enkel`) |
| Mehrere Zweige in einem Push: ein Lauf je Zweig, jeder mit denselben Zusagen | Bestand, `SPEC-061` Punkt 16 (4) |
| Kein `/proc` (macOS, minimale Umgebung) | Architect, `SPEC-061` Kopf, Plattform (nur Linux zugesagt; der Gegenproben-Container ist Linux, der Hook nur dort aufgesetzt) |
| Aufruf von Hand oder aus `make` statt durch `git receive-pack` | Bestand, `SPEC-061` Punkt 5 |
| Mindestversion von bash für den Hook und die Gegenprobe; Äquivalenz des Mutanten ohne `exec` | Architect, `SPEC-061` Punkt 5 (Hook und Gegenprobe bash ≥ 5.1, Entscheidung des Nutzers vom 2026-10-18: festgelegt; Messung der Äquivalenz mit `bash:4.4` bis `bash:5.2`: das Auslassen des `fork` beginnt mit 5.1; die Gegenprobe prüft die Version vorab, Exit 2) |
| `env --default-signal` (coreutils ab 8.31) und `/proc` fehlen | Architect, `SPEC-061` Punkt 13 (`env` wird geprüft, Fall und Mutation) und Kopf, Plattform (`/proc` nicht geprüft, akzeptiertes Negativ) |
| Echter `git push` über `ssh` | Architect, `SPEC-061` Punkt 16 (11) |

**Aufrufarten des Treibers (§3.12):** nicht Gegenstand dieses Slice (`SPEC-061` Punkte 10, 11 und 14 sind Bestand des Gebers). **Aufruf-Voraussetzung:** Die Zusagen des Punktes 5 gelten für jeden Aufrufer, der den Hook startet, gleich welche Deskriptoren und Signalzustände er vererbt; gleichzeitig laufende Hooks sind Punkt 6.

Risiken:

- **Weitere vererbte Zustände, die die Gegenprobe nicht kennt** (Umgebungsvariablen von `git`, `umask`, Arbeitsverzeichnis, Ressourcenlimits, Signalmaske `SigBlk`): Die zwei Funde der ersten Fahrt sind zwei von mehreren Größen des Aufrufers. — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/hook-erbt-deskriptor-und-ignorierte-signale`; die erste Fahrt von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` ist der Beleg, ob weitere eintreten.
- **Die Fälle ersetzen die Fahrt nicht:** Die Gegenprobe läuft mit Attrappen und lokalem Push; ein Rot der Fahrt, das hier nicht vorkommt, bleibt möglich. — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`; die Fahrt in `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (DoD 3 dort) zeigt, dass der `ssh`-Push nicht mehr hängt.
- **`setsid -f` ändert eine Zusage von Punkt 5 oder 9** (Sitzung, Signal an die Prozessgruppe, Log, Zeile `gestartet`). — **Ausgang (vorgesehen):** entfallen, wenn die Gegenprobe des Gebers am Endstand ohne geänderten Fall der Punkte 5 und 9 grün ist (Bericht); sonst eingetreten → Rückführung nach `open` (§4).
- **Der Slice wächst durch Randformen des Aufrufers**, die der Architect nachträgt (§1 schließt weitere aus). — **Ausgang (vorgesehen):** eingetreten → Rückführung nach `next` (§4) und Schnitt durch den Planner; sonst entfallen mit der Messung.

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

- **Was hat funktioniert:** Der Vertrag trug: Der Architect entschied die Randformen von `SPEC-061` Punkte 5, 11 und 13 vor dem Code und nach dem Review ein zweites Mal; die Gegenprobe färbt je Mutation genau die benannte Menge (113 Fälle, 156 Mutationen, im Gegenproben-Container grün). Das geforderte Vorher/Nachher ist belegt: Der Hook aus `02a606eb` ist mit der neuen Gegenprobe rot an genau den Fällen der zwei Funde der ersten Fahrt (Deskriptoren, `SigIgn`), der Hook am Kopf grün; die Verifikation fuhr 51 Mutationen einzeln und zwei eigene Mutanten des Hooks mit derselben Menge. Review (0 HIGH, 1 MEDIUM, 3 LOW, 4 INFO) und Verifikation (bestanden; 0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO) fanden die Lücken vor der Closure. Größe 224 nach der Metrik von §4 (Summe aus Einfügungen und Löschungen, `git diff --numstat 4219bab6..HEAD -- tools harness/mk harness/README.md`), Marke 900 und Meldeschwelle 700 nicht berührt.
- **Was ging anders als geplant:** (1) Die Schätzung 90 bis 180 lag unter der Messung 224 (Faktor knapp 1,3 über der Obergrenze): Das Review verlangte eine Vorab-Prüfung der Version von bash mit Fall und sieben Mutationen und eine Probe für `env --default-signal`, die Nutzer-Entscheidung „bash ≥ 5.1“ kam dazu. (2) Fünf Rückgaben der ersten Runde entschied der Architect nach dem Code (Punkte 5, 9, 11); eine davon änderte den Plan: `SIGQUIT` in der Shell des Laufs gilt nicht (bash ignoriert es selbst, Messung M1), §1, DoD und Titel sagten es anders (F-3, bei dieser Closure nachgezogen). (3) Die Gegenprobe ist unter hoher Last (30 bis 46, parallel `make gates`) für zwei Fälle des Bestands einmal zusätzlich rot (`p2-kein-verzeichnis`, `p13-werkzeug-git`), Ursache ungeklärt; das Rot von `m-p13-ohne-git` im Stand `a918ef2e` war dagegen ein Fehler des PATH der Werkzeugfälle, behoben in `07e45e1e`. (4) Die Äquivalenz von `m-p5-lauf-ohne-exec` hängt an der Version von bash; ihr Ort ist die Versionsprüfung der Gegenprobe, nicht der Hook (V-3).
- **Steering-Loop-Eintrag:** *neuer Sensor* und *benannte Lücke*. Sensor: Fälle und Mutationen der Gegenprobe `make gegenproben-nach-einlieferung-gegenprobe` für den Zustand, den der Aufrufer des Hooks vererbt: `p5-eof-fd3`, `p5-eof-fd5`, `p5-eof-fd7`, `p5-eof-fd12` (Deskriptor an einer Pipe, Leser wartet bis zum Ende der Leitung; Sonde, dass der Deskriptor eine offene Pipe ist), `p5-fd-lauf|make|kind`, `p11-sigint-*`, `p11-sigquit-*`, `p11-erbt-int|quit` (Hook aufgerufen mit ignoriertem Signal), `p5-bash-version` und `p13-werkzeug-env` mit ihren Mutationen; liegt in `tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh` und `tools/harness/gegenproben-nach-einlieferung.sh`, getragen von `SPEC-061` Punkte 5, 11, 13, Zeile in `harness/README.md` §Sensors; Herkunfts-Anker `seit slice-harness-gegenproben-aus-gates-fahrstelle-umgebung`. Lücken (benannt, nicht verkörpert; Vorschläge des Planners, nicht entschieden): (a) `AGENTS.md` §3.12 nennt die Umgebung, in der eine Gegenprobe läuft, nicht als Randform (`BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`, 6×, Ausgang beim Nutzer; hier kam der Aufrufer, die Last und die Version von bash als Größen hinzu); (b) die Ort-Klausel in `AGENTS.md` §3.10 für einen grünen Mutanten nennt eine Eigenschaft des Interpreters nicht: Gilt die Äquivalenz nur ab einer Version, nennt §7 die Version als Voraussetzung und als Ort die Prüfung dieser Version im Repo (`BEO-REPO/gruener-mutant-faelschlich-aequivalent`, Zähler 8).
- **Beobachtungs-Register (`../observations/`):** neu keine Beobachtung; `BEO-REPO/hook-erbt-deskriptor-und-ignorierte-signale/` (1×, Stand aktualisiert: Gegenmittel getragen, Wirkung steht aus), `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts/` (Beleg des Vorgangs um Last und bash-Version ergänzt, Zähler bleibt 6), `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt/` (viertes Auftreten, neue Datei unter `evidence/`), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` (F-1; Zähler 56), `BEO-REPO/gruener-mutant-faelschlich-aequivalent/` (F-2; Zähler 8), `BEO-REPO/plan-folgt-korrektur-nicht/` (F-3; Zähler 68), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` (V-1, F-4; Zähler 83), jeweils mit `evidence/slice-harness-gegenproben-aus-gates-fahrstelle-umgebung.md`.
- **Folge-Slices:** `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (eine Datei in `open/`; die Sendung (Beleg der ersten Fahrt, Dauer des Pushs, `container-gegenprobe` nicht rot an `SIGINT`) steht in dessen §1, §4 und §6; der Planner trägt bei dieser Closure den Nachtrag zu den Versionen von bash und coreutils in dessen §1 und §8 ein, Liefer-Punkte bleiben 3, Schichten 2).
- **Risiken aus §6:** (1) weitere vererbte Zustände (Umgebung von `git`, `umask`, Arbeitsverzeichnis, Limits, `SigBlk`): *weiter offen* → `BEO-REPO/hook-erbt-deskriptor-und-ignorierte-signale`. (2) die Fälle ersetzen die Fahrt nicht: *weiter offen* → `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`. (3) `setsid -f` ändert eine Zusage von Punkt 5 oder 9: *entfallen*, weil von den Fällen des Bestands keiner entfernt oder geändert ist (Verifikation, Plan-vs-Code-Diff) und die Gegenprobe am Kopf grün ist. (4) der Slice wächst durch Randformen: *entfallen*, gemessen 224 gegen Marke 900; die Randformen standen vor dem Code in §6 (die Erweiterungen nach dem Review, Versionsprüfung und `env`-Probe, sind in der Messung).
- **Findings je Ausgang (Review F-1 bis F-8, Verifikation V-1 bis V-4):** *entfallen* (behoben im Slice): F-1 (Mutationen je Element, `9230dbab`), F-2 (bash ≥ 5.1, `9230dbab`, `58f047ba`), F-3 (Plan nachgezogen, Closure), F-4 (Verweis „Punkt 16 (11)“, `9230dbab`), F-5 (Kommentar im Hook trägt eine Kopplung, Test `m-p11-hook-mit-und`; Urteil: Abgrenzung, kein Befund), F-6 (Probe für `env`, `9230dbab`), F-7 (Listen im Vertrag an die Messung angeglichen, `b393f1cc`), V-1 (Kommentar `g11s`, `3d8ba0a0`), V-2 (Closure-Pflichten, mit dieser Closure erledigt), V-3 (kein Handlungsbedarf: Ort der Äquivalenz ist die Versionsprüfung der Gegenprobe; ob die Prüfung in den Hook gehört, entscheidet der Nutzer, Vertrag unverändert). *Weiter offen* → `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`: F-8 und V-4 (Fälle des Bestands unter Last einmal zusätzlich rot, Ursache ungeklärt).
- **Drei Paarungen:** Anker: der Sensor liegt in `tools/harness/gegenproben-nach-einlieferung-gegenprobe.sh` und `tools/harness/gegenproben-nach-einlieferung.sh`, getragen von `SPEC-061` Punkte 5, 11, 13, Zeile in `harness/README.md` §Sensors; eine Regel wurde nicht verkörpert (die Vorschläge oben sind keine). Folge-Slice: `grep` nach `fahrstelle-umgebung` im Nehmer `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` findet die Kennung in §1, §4 und §6; der Nehmer liegt in `open/`, ist keine geschlossene Adresse und führt die Sendung in §1. Register: jede Beobachtung des Lerneintrags ist ein Verzeichnis unter `../observations/BEO-REPO/` mit `evidence/slice-harness-gegenproben-aus-gates-fahrstelle-umgebung.md` (fünf neue Belege, einer ergänzt, einer im Stand geführt).
- **Offene Punkte für den Nutzer:** (1) `make gates` lief am Stand `07e45e1e` grün (Exit 0); seitdem ändert sich unter `tools/` nur ein Kommentar (`3d8ba0a0`, Zeile `g11s`; `make abdeckung-check` liest nur den Kopf der Gegenproben, SPEC-048 Punkt 7), das Urteil des Planners ist, dass der Lauf den Stand des Codes belegt; der volle `make gates` am Integrationsstand (nach dem Zusammenführen mit `origin/main`) fährt der Host. (2) Regel „Umgebung einer Gegenprobe“ (`AGENTS.md` §3.12): `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` steht bei 6× ohne Ausgang; Last und Version von bash sind weitere Größen. (3) Ort-Klausel von `AGENTS.md` §3.10 für Äquivalenz an einer Werkzeugversion (Vorschlag oben). (4) V-3: Soll der Hook seine bash-Version selbst prüfen? (5) Die Wirkung (der `ssh`-Push kehrt zurück) zeigt erst die erste Fahrt in `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`.

### Belege des Implementers

Stand: Zweig `slice-harness-gegenproben-aus-gates-fahrstelle-umgebung`, Basis `02a606eb`. Code-Commits `7b4eabcc`, `51b5ba8d`, `c26c9092`, danach nach der Entscheidung des Architect (`SPEC-061` Punkte 5, 9, 11, Commit `597f2ed0`) `6c740c67` und `3d7b32e8`. Dieser Abschnitt ist der einzige Teil des Plans, den der Implementer geändert hat (`git diff 02a606eb..HEAD -- docs/plan` zeigt neben den §6-Zeilen des Architect nur §7). Nach dem Review (`0761194d`, Findings F-1 bis F-8) und der Entscheidung des Architect zu ihm (`b393f1cc`) sowie der Festlegung des Nutzers „bash ≥ 5.1“ (`58f047ba`): Code `9230dbab`.

**Rückgaben:** Die fünf Rückgaben der ersten Runde sind vom Architect entschieden und umgesetzt. Offen ist eine Anmerkung zu den Listen im Vertrag (kein Randfall): `SPEC-061` Punkt 11 nennt für die Mutation „der Hook startet den Lauf mit `cmd &`“ die sieben Fälle ohne `erbt` und für `trap '' INT`, `trap '' QUIT` und „`make` mit `cmd &`“ Mengen ohne `p11-erbt-*`. Gemessen färbt jede dieser Mutationen auch den passenden Fall `p11-erbt-*` (er fährt `make` ebenfalls; der Mutant „Hook mit `cmd &`“ enthält zudem kein `env --default-signal`, das die ignorierten Signale zurücksetzte, ein `cmd &` nach dem `env` gibt es nicht). Die Mengen unten sind die gemessenen; der Architect möge die Listen im Vertrag angleichen. Zum Review (0 HIGH, 1 MEDIUM, 3 LOW, 4 INFO) steht jeder Befund mit Ausgang unter *Findings des Reviews*.

**Messungen (Container `pgr-dev2`, bash 5.2.37, util-linux 2.41.5, coreutils 9.7):**

- M1 — bash ignoriert `SIGQUIT` selbst: `env --default-signal=INT,QUIT bash z.sh`, `z.sh` liest `/proc/$$/status` mit dem Builtin `read` (kein Kindprozess): `SigIgn: 0000000000000004`; ein Handler (`trap 'echo q' QUIT`) setzt es auf `0`; `env --default-signal=QUIT sh -c '…'` (dash): `0`.
- M2 — `timeout` setzt `SIGINT` zurück: `setsid -w sh -c 'trap "" INT; exec env -i sh -c "grep SigIgn /proc/\$\$/status"'` zeigt `0000000000000002`; mit `timeout 30` vor dem inneren `sh`: `0000000000000000`. Die Fälle `p11-erbt-*` rufen den Hook deshalb ohne `timeout`.
- M3 — Ein Kind der `ssh`-Sitzung erbt kein ignoriertes `SIGQUIT`: `ssh … 'env -i sh -c "grep SigIgn /proc/\$\$/status"'` zeigt `0`; die `bash -c`-Shell der Sitzung selbst zeigt Bit 3 (M1).
- M4 — `bash -c '( exec bash -c "grep SigIgn …" ) & wait'` aus einer Shell ohne Jobsteuerung: Kind `0`; `sleep 5 & grep SigIgn /proc/$!/status`: `0000000000000006`.

**Befunde vorher und nachher:** (1) Hook aus `02a606eb` mit der Gegenprobe des Standes `51b5ba8d`: die unveränderte Kopie ist rot an genau `p11-sigint-kind p11-sigint-lauf p11-sigint-make p11-sigquit-kind p11-sigquit-make p5-eof-fd12 p5-eof-fd3 p5-eof-fd5 p5-eof-fd7 p5-fd-kind p5-fd-lauf p5-fd-make`, sonst nichts. (2) Hook aus `405999db` (ohne `env --default-signal`) mit der Gegenprobe dieses Standes: rot an genau `p11-erbt-int p11-erbt-quit`. Nach der Änderung: unveränderte Kopie grün, 113 Fälle, 156 Mutationen.

**Mutationstabelle** (rote Menge gemessen mit `GP_ENTDECKE=1`; der volle Lauf verlangt genau diese Mengen):

| Zusage (SPEC-061) | Mutation | rote Fälle |
|---|---|---|
| 5: kein Deskriptor über 2 im Lauf, ohne Liste fester Nummern | `eval "exec $n>&-"` entfällt | `p5-eof-fd3 p5-eof-fd5 p5-eof-fd7 p5-eof-fd12 p5-fd-lauf p5-fd-make p5-fd-kind` |
| 5: … Deskriptor 3 (untere Grenze) | geschlossen wird nur über 3 | `p5-eof-fd3` |
| 5: … nicht nur 5 | nur Deskriptor 5 | `p5-eof-fd3 p5-eof-fd7 p5-eof-fd12` |
| 5: … nicht nur 3 bis 9 | nur 3 bis 9 | `p5-eof-fd12` |
| 5: … 0 bis 2 bleiben offen | geschlossen wird ab 2 | `p5-fd-2` |
| 5: `make` und Nachkommen tragen die Sperre nicht | `9>&-` entfällt | `p5-fd-make p5-fd-kind` |
| 5: der Hook kehrt zurück, bevor `make` endet; `setsid -f` | `setsid -w` | `p1-echter-push p5-ein-prozess p5-eof-fd3 p5-eof-fd5 p5-eof-fd7 p5-eof-fd12 p5-fd-0 p5-fd-1 p5-fd-2 p5-leitung-stderr p5-leitung-stdout p5-vor-ende p7-altes-weg p9-marke p9-nicht-vorher` |
| 5, 16 (10): ein Hintergrundprozess je Push | die `bash -c`-Stufe bleibt als Elternprozess (`"$BASH" "$@"; :`) | `p5-ein-prozess` |
| 5: eigene Sitzung | ohne `setsid`, mit `&`, Signale per `env --default-signal` zurückgesetzt | `p5-gruppe` |
| 11: `SIGINT`/`SIGQUIT` in `make`, Kind, Enkel; `SIGINT` im Lauf (Start durch den Hook) | `setsid` ohne `-f`, Start mit `cmd &`, ohne `env --default-signal` | `p11-sigint-lauf p11-sigint-make p11-sigint-kind p11-sigint-enkel p11-sigquit-make p11-sigquit-kind p11-sigquit-enkel p11-erbt-int p11-erbt-quit` |
| 11: … Start von `make` durch den Lauf | `env -C "$wt" make … &` statt `( … exec make …) &` | `p11-sigint-make p11-sigint-kind p11-sigint-enkel p11-sigquit-make p11-sigquit-kind p11-sigquit-enkel p11-erbt-int p11-erbt-quit` |
| 11: `SIGINT` | `trap '' INT` im Lauf | `p11-sigint-lauf p11-sigint-make p11-sigint-kind p11-sigint-enkel p11-erbt-int` |
| 11: `SIGQUIT` | `trap '' QUIT` im Lauf | `p11-sigquit-make p11-sigquit-kind p11-sigquit-enkel p11-erbt-quit` |
| 11: der Aufrufer ändert nichts daran | `env --default-signal=INT,QUIT` entfällt | `p11-erbt-int p11-erbt-quit` |
| 11: … nur `SIGINT` / nur `SIGQUIT` wird zurückgesetzt (Element der Aufzählung) | `--default-signal=INT` allein / `=QUIT` allein | `p11-erbt-quit` / `p11-erbt-int` |
| 5: … jeder Deskriptor über 2, nicht „alle außer 5“ / „außer 7“ | geschlossen wird jeder außer 5 / außer 7 | `p5-eof-fd5 p5-fd-lauf p5-fd-make p5-fd-kind` / `p5-eof-fd7` |
| 13: `env` mit `--default-signal` wird geprüft (Meldung `werkzeug`, Name `env`) | die Prüfung entfällt | `p13-werkzeug-env` |
| Kopf, 5: die Gegenprobe prüft bash ≥ 5.1 vorab (Exit 2, `bash >= 5.1 nötig`) | `m-bash-ab-5-0` (`-ge 0`), `m-bash-ab-5-2` (`-ge 2`), `m-bash-major-ab-4` (Hauptversion über 3), `m-bash-major-nicht-6` (über 6), `m-bash-nicht-5` (`-eq 4`), `m-bash-ohne-meldung`, `m-bash-exit-1` an einer Kopie der Gegenprobe; der Fall fährt sie mit vorgetäuschter Version (`GP_BASH_VERSION` 4.4.0, 5.0.19, 5.1.0, 5.2.37, 6.0.0) | je Mutation genau `p5-bash-version` |

Die Mengen der übrigen Mutationen des Bestands sind unverändert bis auf vier Zeilen mit zusätzlichen roten Fällen, je erklärt: `m-p1-liest-stdin-nicht` (der Hook liest keine Zeile, kein Lauf entsteht: alle Fälle, die einen Lauf brauchen, sind rot, auch die neuen), `m-p5-wartet` (Zeile oben: ein wartender Hook hält die Pipe, bis der Lauf endet, und `setsid -w` ist ein zweiter Prozess mit `--lauf` in der Befehlszeile), `m-p3-gleichzeitig` (+ `p5-ein-prozess`: der Lauf startet seine Läufe als Hintergrund-Subshell, ein zweiter Prozess mit `--lauf`) und `m-p12-nicht-losgeloest` (+ `p11-erbt-quit`: der Fall fährt den Hook zweimal im selben Bare-Repo mit demselben Hash, der Zweig `gp-<Name>` besteht beim zweiten Mal schon).

**Grüne Mutanten:** `exec "$BASH" "$@"` → `"$BASH" "$@"` (`m-p5-lauf-ohne-exec`, erwartet: nichts rot): äquivalent **ab bash 5.1**, der vom Nutzer festgelegten Mindestversion für Hook und Gegenprobe. (a) Invariante: `bash -c` führt den letzten einfachen Befehl eines `-c`-Strings ab 5.1 im selben Prozess aus, ohne `fork`. (b) Ort: das Verhalten von bash; im Repo erzwingt es die Versionsprüfung der Gegenprobe (`bash >= 5.1 nötig`, Fall `p5-bash-version`, sieben Mutationen `m-bash-*`), die Äquivalenz ist an die Version gebunden. (c) Sonde: Messung des Architect am 2026-10-18 mit den Images `bash:4.4`, `bash:5.0`, `bash:5.1`, `bash:5.2` (ab 5.1 gleicher Prozess, bis 5.0 mit `fork`); eigener Lauf im Container (bash 5.2.37): der Mutant bleibt grün, `"$BASH" "$@"; :` (`m-p5-lauf-als-kind`) ist rot an `p5-ein-prozess`. Weitere grüne Mutanten gibt es nicht; `SigIgn` eines Enkels ist gedeckt (`p11-sigint-enkel`, `p11-sigquit-enkel`). **Keine Mutation, die nur Kind oder nur Enkel bricht:** Ein ignoriertes Signal erbt nach unten; wer es für das Kind setzt, setzt es für den Enkel mit, und jeder Weg, `make` zu brechen, trifft das Kind mit. Die Fälle `-kind` und `-enkel` färben deshalb immer gemeinsam (F-8 des Reviews zur Kenntnis).

**Sonde:** Der Fall `p5-eof-*` prüft vorab, dass der Deskriptor eine offene Pipe ist. Eine Kopie, in der `eof_gruppe` den Deskriptor an `/dev/null` legt, ist rot mit `p5-eof-fdn: Deskriptor n ist im Fall keine offene Pipe` (gefahren für 5, 7 und 12 vor dem Fall `-fd3`, derselbe Code).

**Akzeptierte Negative:** Dass ein echter `git push` über `ssh` zurückkehrt, prüft die Gegenprobe nicht (`SPEC-061` Punkt 16 (11)); den Beleg trägt die erste Fahrt in `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`. Für `SIGQUIT` in der Shell des Laufs gibt es keine Zusage (Punkt 11, M1). Fehlendes `/proc` prüft der Hook nicht (Plattform, Kopf des Vertrags).

**Größe:** Summe aus Einfügungen und Löschungen von `git diff --numstat 4219bab6..HEAD -- tools harness/mk harness/README.md`: 224 (199 Einfügungen, 25 Löschungen; Schätzung des Plans 90 bis 180) (Meldeschwelle 700, Marke 900).

**Findings des Reviews und Ausgang:** F-1 (LOW) behoben: Mutationen `--default-signal=INT` allein, `=QUIT` allein, „außer 5“, „außer 7“ (Tabelle). F-2 (MEDIUM) behoben: Mindestversion bash ≥ 5.1 (Nutzer), Vorab-Prüfung der Gegenprobe mit Fall und sieben Mutationen, Äquivalenz mit Version, Messung und Ort (Grüne Mutanten). F-3 (LOW) beim Planner: §1 und DoD sagen noch „weder `SIGINT` noch `SIGQUIT`“ für den Lauf, der Implementer ändert dort nichts (Closure). F-4 (LOW) behoben: der Verweis lautet „Punkt 16 (11)“ in der README-Zeile und hier. F-5 (INFO) behoben: Kommentar im Hook auf den Wortlaut des Architect. F-6 (INFO) behoben: Probe `env --default-signal=INT,QUIT true` im Werkzeug-Schritt, Fall `p13-werkzeug-env`, Mutation `m-p13-ohne-env`; kein `/proc` ist akzeptiertes Negativ im Vertrag. F-7 (INFO) behoben: die erwarteten Mengen entsprechen den Listen im Vertrag (der volle Lauf verlangt sie). F-8 (INFO) zur Kenntnis: Ursache von `p2-kein-verzeichnis` unter Last ungeklärt, siehe Läufe.

**Doku:** `harness/README.md` §Sensors, Zeile `make gegenproben-nach-einlieferung-gegenprobe`: Zahlen, die neuen Fälle, die Voraussetzungen (`sh`, `env` mit `--default-signal`) und das Nicht-Geprüfte (`ssh`-Push, `SIGQUIT` in der Shell des Laufs). Quantoren der neuen Sätze: „ohne geerbten Deskriptor über 2“ — Ausnahmen derselben Zeile: die Sperre (9) und der Skript-Deskriptor (255) im Lauf, beide im Muster von `p5-fd-lauf`; „`SIGINT` und `SIGQUIT` in `make` und seinen Nachkommen nicht ignoriert“ — Ausnahme: die Shell des Laufs ignoriert `SIGQUIT` selbst, in der Zeile unter „nicht geprüft“ genannt.

**Läufe:** Stand `07e45e1e` (die Dateien unter `tools/` und `harness/README.md` sind seitdem unverändert). Gegenprobe im Container, UTC: 18:15:18 bis 18:22:37 UTC, Exit 0, `gruen — Punkte 1 bis 13 von SPEC-061, gefahren: 113 Fälle und 156 Mutationen, je Mutation genau die benannten Fälle rot, unveränderte Kopie grün` (bash 5.2.37). `make gates`: 18:08:42 bis 18:13:33 UTC, Exit 0. `make gates-schnell-gegenprobe`: Exit 0. `make docs-check`: Exit 0, 1156 Dateien, 0 Befunde. Zeitgrenzen `timeout 1500`, `3000`, `900`, `300`, keine ausgelöst; Prozesse im Container 94 bis 276. Auf einen Lauf mit Rot kamen zwei grüne Wiederholungen: Ein voller Lauf des Standes `a918ef2e` war rot an `m-p13-ohne-git` (der PATH der Werkzeugfälle trug `env` und `true` noch nicht; behoben in `07e45e1e`), ein voller Lauf des Standes `07e45e1e` (18:01:21 bis 18:08:42, parallel zu `make gates`) war rot, weil `m-p5-nur-ueber-3` zusätzlich `p13-werkzeug-git` rot hatte; die Wiederholung allein ist grün (Zeile oben). Mit `p2-kein-verzeichnis` (zwei frühere Läufe) sind das drei Fälle des Bestands, die bei hoher Last einmal zusätzlich rot waren (Last 30 bis 46 und parallel laufendes `make gates`); die Ursache ist ungeklärt (F-8 des Reviews). Ein Lauf mit `TZ=Europe/Berlin LANG=de_DE.UTF-8` am Stand `3d7b32e8` war grün (17:12:06 bis 17:18:38 UTC), nach der letzten Änderung nicht wiederholt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien**, vier und nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (4× zuvor, mit diesem Vorgang 5×; Ausgang offen, Vorschlag des Planners in `state.md`), `BEO-REPO/gegenprobe-nicht-in-kette-regression-unbemerkt` (Regression des Gebers erst durch die Fahrt bemerkt; der Beleg dieses Slice wird bei der Closure geführt), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (Gegenmittel: Mutation je Zusage, DoD 2), `BEO-REPO/hook-erbt-deskriptor-und-ignorierte-signale` (neu, 1×). Keiner erreicht mit diesem Slice neu 3×.

alle berührten Sub-Areas GF (`REPO`).

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-18, im Commit des Anlegens):** Liefer-Punkte 2 (Skript · Gegenprobe). Schichten 1: Harness-Werkzeuge (`tools/harness`); Vertrag und Planung sind keine Schicht. Schätzung 90 bis 180 in der Metrik von §4, Marke 900, Meldeschwelle 700 (weit darunter). In diesen Slice ist keine Sendung eingetragen außer der des Gebers (§1; kein Liefer-Punkt, keine Schicht).
