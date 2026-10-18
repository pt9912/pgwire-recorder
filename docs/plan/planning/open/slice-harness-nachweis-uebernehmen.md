# Slice slice-harness-nachweis-uebernehmen: Ein Werkzeug übernimmt den Nachweis des vollen Laufs aus einer Quelle, wenn die Dateiliste des Baums, der hier liegt, mit der des Nachweises übereinstimmt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten, fünfte Stelle der Reihenfolge, hinter `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` und vor `slice-harness-gegenprobe-baumkopie`; **aus der Reihenfolge der Wellen gelöst, Start nach dem Kettenumbau** (Entscheidung des Nutzers vom 2026-10-18, so vom Koordinator übermittelt, der Wortlaut liegt dem Planner nicht vor): Start-Trigger ist `slice-harness-gegenproben-aus-gates` in `done/`, nicht `welle-harness-regeln`. Angelegt am 2026-10-18 auf Auftrag des Nutzers („Nachweis aus dem Container oder Worktree übernehmen — das hört sich pragmatisch an“). **Geschnitten am 2026-10-18** (Regel des Nutzers: schneiden, wenn die Schätzung die Marke selbst überschreitet; Schätzung des Architect für den ungeteilten Slice 830 bis 1180 Zeilen, Marke 900): Dieser Slice behält Teil 1 (Werkzeug, Quelle, Gleichheit, Ziel, `-i`), die Signale, die übrigen Aufrufarten und die Kopplung an den Hook sind in `slice-harness-nachweis-uebernehmen-aufrufarten` (neu, `open/`) aufgegangen, direkt dahinter. Der Slice liegt in `open/`: Er ist nicht priorisiert, den Anspruch trägt der Nutzer.

**Bezug:** [LH-QA-07](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur insoweit, als das Werkzeug und seine Gegenprobe dessen Messmethode (4) tragen; kein Anspruch des Produkts, der Slice ändert die Gate-Kette des Harness. Aktive ADR: [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) (Accepted, Entscheidung des Nutzers vom 2026-10-18, Option B: der Hook glaubt auch einen Nachweis, der den Baum beschreibt und von einer anderen Umgebung stammt), [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md). Bindung an das Register: `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (Gegenmaßnahme: der Nachweis reist mit dem Stand), `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (2×; der Zähler gehört seit dem Schnitt zu slice-harness-nachweis-uebernehmen-aufrufarten, §8).

**Berührte Spec-Stellen:** [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Teil 1 und die Punkte 1 bis 14, 15 (die drei Fälle ohne Signal) und 16 (b) (der Vertrag liegt vor, `SPEC-062` ist vergeben; die Punkte 15 mit Signal, 16 (a), (c), (d), 17 und 18 liefert slice-harness-nachweis-uebernehmen-aufrufarten); [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) Punkte 1, 8 und 21 (der Nachweis bekommt einen zweiten Schreiber; der Hook bleibt unverändert; der Architect hat sie geändert, der Slice liest sie), dazu Punkt 12 (ein neues Skript unter `tools/` ist nicht inert); [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) Teil 1 (die Gegenprobe des Werkzeugs steht nach der Namensregel in einer Liste), nur gelesen.

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

**Ziel:** `make nachweis-uebernehmen` (Namen laut `SPEC-062`) legt den Nachweis des vollen Laufs, den eine andere Umgebung geschrieben hat (Worktree, Klon im Dev-Container), im Arbeitsbaum ab, in dem es läuft, **genau dann**, wenn (a) die Dateiliste dieses Baums, frisch mit `dateiliste_erzeugen` berechnet, byte-gleich mit der Dateiliste der Quelle ist, (b) die Liste der Quelle unbeschädigt ist (`dateiliste_pruefen`) und (c) `gates-passed.diffsha` der Quelle der SHA-256 genau dieser Liste ist. Sonst lehnt es ab und ändert nichts. Der Hook bleibt unverändert: Er liest nur die Liste (`.claude/hooks/stop-require-gates.sh`, `cmp -s`, `SPEC-052` Punkt 16), und eine byte-gleiche Liste ist seine Freigabe-Bedingung; neu ist ein zweiter Schreiber des Standes, kein neuer Leser. Das Werkzeug ist kein Gate, hängt nicht an `GATE_CHECKS` und schreibt nichts außer den zwei Dateien des Nachweises.

**Teil 1 und Teil 2.** Dieser Slice liefert Teil 1 des Vertrags: die Punkte 1 bis 14 (Ausgang, Reihenfolge der Gründe, Quelle, Gleichheit, Ziel), vom Punkt 15 die drei Fälle ohne Signal (Schreiben scheitert) und vom Punkt 16 den Buchstaben (b), `-i` in `MAKEFLAGS` (71 von 93 Fällen). `slice-harness-nachweis-uebernehmen-aufrufarten` liefert den Rest: `SIGINT` und `SIGTERM` (Punkt 15), `-n`, `-q`, `-t`, `-k`, `-j` und Langflags (Punkt 16 (a), (c), (d)), den Aufruf aus einem Unterverzeichnis (Punkt 17) und die Kopplung an den Hook (Punkt 18) (22 Fälle). `-i` gehört in Teil 1, weil das Skript es ablehnt (Punkt 2, Grund 1): Eine Ablehnung ohne Fälle für ihre Schreibweisen wäre eine Zusage ohne Mutation (`AGENTS.md` §3.10), und unter `make -i` sähe die Ablehnung wie Erfolg aus.

**Befund am Stand `0a59a42f` (gelesen, nicht geändert):**

- **Der Nachweis** sind zwei Dateien in `.harness/state/` (ignoriert, also je Arbeitsbaum eigen, `bin/` darunter ebenso): `gates-passed.dateiliste` (Ausgabe von `dateiliste_erzeugen`: je Pfad Art, Wert, Pfad) und `gates-passed.diffsha` (SHA-256 der Liste, eine Zeile). `tools/harness/record-gates.sh` schreibt beide als letztes Prerequisite von `make gates`, also nur nach grünen Gates, und entfernt beide bei `-i` in `MAKEFLAGS`. **Die Dateien tragen weder Zeitpunkt noch Commit noch Herkunft**; ein „grüner Lauf“ ist an nichts weiter erkennbar als daran, dass beide da sind und zusammenpassen (§6, Gruppe S).
- **Der Hook liest nur die Liste.** `gates-passed.diffsha` liest kein Skript außer der Gegenprobe (`stop-hook-gegenprobe`, Fälle `p1-…`); die Zusammengehörigkeit prüft deshalb das neue Werkzeug, nicht der Hook.
- **Was die Liste enthält:** getrackte und untracked Pfade ohne ignorierte, je Eintrag Art (`F`, `L`, `G`, `X`) und SHA-256 bzw. Link-Ziel (`SPEC-052` Punkt 11). **Nicht** darin: `.harness/state/`, `dist/` und alles, was `.gitignore` ausblendet; Dateimodi (akzeptiertes Negativ, Punkt 11); Zeitstempel; `.git` (im Worktree eine Datei). Zwei Umgebungen desselben Commits unterscheiden sich in der Liste also nur durch **untracked, nicht ignorierte Dateien**, durch Inhalte, die `git` beim Auschecken umschreibt (`core.autocrlf`, `.gitattributes`), und durch lokale Ignore-Regeln (`.git/info/exclude`, `core.excludesFile`), die eine untracked Datei in der einen Umgebung ausblenden und in der anderen nicht.
- **Ein roter Lauf lässt einen älteren Nachweis stehen** (`SPEC-052` Punkt 21, akzeptiertes Negativ). Das schadet dem Werkzeug nicht: Es vergleicht die Liste der Quelle mit dem Baum, der hier liegt, nicht mit dem Stand der Quelle; ein älterer Nachweis ist für den Baum, den er beschreibt, weiter wahr, und für jeden anderen ist die Liste ungleich.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `SIGINT` und `SIGTERM`, `-n`, `-q`, `-t`, `-k`, `-j`, Langflags, der Aufruf aus einem Unterverzeichnis und die Fälle der Kopplung an den Hook (Punkte 15 mit Signal, 16 (a), (c), (d), 17 und 18 von `SPEC-062`) — Klasse 1, `slice-harness-nachweis-uebernehmen-aufrufarten` übernimmt sie (dessen §1 und DoD nennen die Kennung dieses Slice; Start-Trigger dort: dieser Slice in `done/`). Dieser Slice sagt dazu nichts zu: Kein Kommentar, kein Text im Rezept und im Gate-Index nennt ein Verhalten bei Signal, bei diesen Flags, aus einem Unterverzeichnis oder bei einer Hook-Antwort als zugesagt (`AGENTS.md` §3.11); die Probe am Endstand (DoD) liest die beiden Antworten des Hooks von Hand und ist keine Gegenprobe.
- Eine Änderung am Hook, an `dateiliste_erzeugen`, `dateiliste_pruefen`, der Positivliste oder an `record-gates.sh` (außer dem Kommentar, der den Schreiber nennt) — Bestand bleibt bewusst stehen: Dass der Hook bei byte-gleicher Liste freigibt, ist `SPEC-052` Punkt 16 mit Fall `p16-identisch`; der Slice stützt sich darauf, statt es zu ändern, und die Gegenprobe des Werkzeugs fährt den unveränderten Hook (§6, Gruppe H; Fälle in `slice-harness-nachweis-uebernehmen-aufrufarten`).
- Ein Transport im Werkzeug (`ssh`, `docker cp`, `git`-Abruf) — anderer Vorgang und eine Host-Abhängigkeit über `bash`, `git`, `docker` und `make` hinaus (`harness/conventions.md`); er brächte Verbindungsabbruch, Teilübertragung und Host-Schlüssel als Randformen mit. Die Quelle ist ein Verzeichnis mit den zwei Dateien; wie sie dorthin kommen, steht als Rezept in der Wartungs-Doku (Liefer-Punkt 3). Frage 2 (Empfehlung des Planners).
- Die Übernahme des Standes der schnellen Kette (`.harness/state/gates-schnell-passed.diffsha`) und eine Freigabe inerter Änderungen ohne diesen Stand — anderer Vorgang: Die erste Übernahme würde eine zweite Gleichheit prüfen (Hash des ganzen Baums), die zweite ändert, was der Hook glaubt ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), Lockerung nach `AGENTS.md` §3.6). Der Nutzer hat am 2026-10-18 entschieden, die zweite Sperre nicht zu ändern (Frage 3, Option (a), Entscheidung des Nutzers, so vom Koordinator übermittelt; die Ablehnungen der Probe werden gemessen); kein Folge-Slice ist angelegt.
- Eine Übernahme bei **ungleicher** Liste, deren Unterschiede alle auf der Positivliste liegen — Schicht-Abgrenzung der Zuständigkeit: Die Positivliste steht an einer Stelle, im Hook (`SPEC-052` Punkt 12); ein Duplikat im Werkzeug wäre eine zweite Quelle, die driftet. Ein Hauptbaum mit unversionierter inerter Änderung lehnt das Werkzeug ab; der Hook entscheidet über inerte Unterschiede selbst.
- Der Beweis, dass die Quelle `make gates` gefahren hat, in derselben Umgebung und mit denselben Gates — nicht prüfbar: Der Nachweis trägt keine Herkunft, und jeder lokale Nachweis ist von Hand zu schreiben (`cp` einer frisch erzeugten Liste, `make record-gates` einzeln). Das Werkzeug ist der sanktionierte Weg und prüft Gleichheit; die Grenze bleibt offen und steht in §6 (Gruppe S).
- Rollen-Anweisungen (`implement-slice`, `close-welle`, Reviewer-, Verifier-, Implementer-Typ) und `AGENTS.md`: ob und wann die Rollen das Werkzeug aufrufen, ist eine Pflicht, die eine Entscheidung des Nutzers nach der Lieferung braucht (Frage 4) — Schicht-Abgrenzung (Rollen-Anweisungen sind eine weitere Schicht, `AGENTS.md` §3.13).
- Das Einrichten und Starten des Dev-Containers und der Klon darin — `slice-harness-devcontainer-einrichten` (Standardweg „Arbeit im Container“, `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`); dieser Slice liefert das Stück, das den Nachweis über die Umgebungsgrenze trägt, nicht den Standardweg.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Werkzeug):** `make nachweis-uebernehmen NACHWEIS_QUELLE=<verzeichnis>` (Namen laut `SPEC-062`) übernimmt `gates-passed.dateiliste` und `gates-passed.diffsha` aus dem Verzeichnis in `.harness/state/` des Arbeitsbaums, in dem es läuft, genau unter den drei Bedingungen aus §1, und lehnt sonst ab, ohne eine Datei zu ändern oder eine Hilfsdatei zu hinterlassen; jede Entscheidung je Randform aus §6 (Gruppen Q, G, Z und `-i` aus A) steht im Vertrag `SPEC-062`, vor dem ersten Code-Commit (liegt vor); das Werkzeug benutzt `dateiliste_erzeugen` und `dateiliste_pruefen` aus `tools/harness/dateiliste.sh` und bildet die Liste nicht nach. Fragment `harness/mk/nachweis-uebernehmen.mk` und Skript `tools/harness/nachweis-uebernehmen.sh`; das Ziel steht nicht an `GATE_CHECKS`. Die Kommentare in `record-gates.sh` und `dateiliste.sh`, die den Schreiber des Nachweises nennen, nennen auch das Werkzeug.
- [ ] **Liefer-Punkt 2 (Gegenprobe):** `make nachweis-uebernehmen-gegenprobe` (Fragment unter `harness/mk/`, in der Liste, die die Namensregel von `SPEC-060` Teil 1 vorschreibt, und ausdrücklich nicht an `GATE_CHECKS`) trägt die 71 Fälle des Teils 1 (Punkte 1 bis 14, `p15-state-ist-datei`, `p15-liste-ist-verzeichnis`, `p15-hash-ist-verzeichnis` und die fünf `-i`-Fälle `p16-i`, `p16-i-lang`, `p16-i-env`, `p16-i-k`, `p16-i-j`), je Randform und je Element ihrer Aufzählungen aus §6 einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die Gegenprobe verlangt je Mutation genau die Menge roter Fälle, die der Vertrag nennt (Namen der Fälle tragen die Nummer des Punktes; Mutationen, die Fälle von slice-harness-nachweis-uebernehmen-aufrufarten mitfärben, gelten dort). Temp-Repos, kein Docker, kein Netz; der Arbeitsbaum bleibt unberührt. Die Fälle der Signale, der übrigen Flags, des Unterverzeichnisses und der Kopplung an den Hook liefert `slice-harness-nachweis-uebernehmen-aufrufarten` (Sendung, §8).
- [ ] **Liefer-Punkt 3 (Gate-Index und Rezept):** `harness/README.md` führt die Zeile des Werkzeugs (Tabelle *Werkzeuge — kein Gate*) und die der Gegenprobe (§Sensors, Bindung `SPEC-062`); `docs/developer/devcontainer.md` trägt einen Abschnitt, der das Werkzeug mit dem Worktree und mit dem Klon im Dev-Container als Quelle zeigt (der Weg der zwei Dateien aus dem Container in ein Verzeichnis ist als Rezept gekennzeichnet; die Gegenprobe prüft eine lokale Kopie der Dateien, nicht `ssh`, `AGENTS.md` §3.11).
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `tools/` und `harness/` ändern); `make nachweis-uebernehmen-gegenprobe` einzeln grün, die Dauer im Bericht; die Messung der Größe nach der Metrik von §4 im Bericht; `make gates-schnell-gegenprobe` grün, falls die Deckung der Ketten die neue Gegenprobe zählt (laut `SPEC-060` Teil 1).
- [ ] **Probe am Endstand (Verifier):** Quelle ist der Worktree des Slice mit grünem `make gates` auf dem Endstand, Ziel ein frischer Worktree desselben Commits ohne Nachweis; der unveränderte Hook antwortet im Ziel vorher `block` (Punkt 19), das Werkzeug endet mit Exit 0, der Hook antwortet nachher `approve`; mit einer untracked Datei im Ziel lehnt das Werkzeug ab und der Hook bleibt bei `block`. Der Hauptbaum wird dafür nicht gebraucht.
- [ ] Beobachtungs-Register: Vermerk (kein Zähler) in `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (Gegenmaßnahme geliefert); `BEO-REPO/nachweis-uebernommen-aus-fremder-umgebung` neu angelegt (Risiko 1, Beleg `evidence/slice-harness-nachweis-uebernehmen.md`, bei der Closure).
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
| `harness/vertraege/SPEC-062-nachweis-uebernehmen.md`, `harness/vertraege/README.md`, `harness/vertraege/SPEC-052-gate-kette.md` Punkte 1, 8 und 21 | Vertrag (Architect, liegt vor dem Code, geschrieben) | Entscheidung je Randform aus Gruppen Q, G, Z, A, H; der zweite Schreiber des Nachweises steht in `SPEC-052` |
| `docs/plan/adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md` | ADR (Architect, `Accepted` am 2026-10-18, Entscheidung des Nutzers, geschrieben) | der Hook glaubt einen Nachweis, den auch eine andere Umgebung geschrieben haben kann; Alternativen und Gründe, `Schärft:` zeigt auf `SPEC-062` (`AGENTS.md` §3.8) |
| `tools/harness/nachweis-uebernehmen.sh`, `harness/mk/nachweis-uebernehmen.mk` | neu | liest die Quelle einmal in ein Hilfsverzeichnis unter `TMPDIR`, prüft Liste und Hash an den Kopien, vergleicht mit `dateiliste_erzeugen`, legt erst dann Hilfsdateien in `.harness/state/` an und ersetzt beide Dateien (Reihenfolge wie `record-gates.sh`: Hash, dann Liste, `mv -T`); lehnt `-i` in `MAKEFLAGS` ab (Exit 2); Kennwörter und Exit-Codes laut `SPEC-062`; keine Behandlung von Signalen im Vertrag dieses Slice (Teil 2) |
| `tools/harness/record-gates.sh`, `tools/harness/dateiliste.sh` | update (nur Kommentare) | nennen den zweiten Schreiber (`AGENTS.md` §3.9) |
| `tools/harness/nachweis-uebernehmen-gegenprobe.sh`, `harness/mk/nachweis-uebernehmen-gegenprobe.mk` | neu | Temp-Repos mit Kopien von Skript und `dateiliste.sh`; je Randform des Teils 1 ein Fall und eine Mutation; Attrappe für `git` (Punkt 11); der Hook liegt nicht in der Kopie, seine Fälle liefert Teil 2 |
| `harness/README.md`, `docs/developer/devcontainer.md` | update | Gate-Index (zwei Zeilen) und Rezept (Worktree, Klon im Container) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` setzt dieser Übergang; der Nutzer priorisiert und trägt den Anspruch. **Vorbedingungen:** (a) Vertrag `SPEC-062` und [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) (Accepted) liegen vor, mit der Entscheidung je Randform aus §6 (`AGENTS.md` §3.12). (b) **Start-Trigger:** `slice-harness-gegenproben-aus-gates` liegt in `done/` (Listen und Namensregel für die Gegenprobe stehen; sonst hinge sie an `GATE_CHECKS` und verlängerte `make gates`). (c) **Reihenfolge:** entschieden (Entscheidung des Nutzers vom 2026-10-18, Frage 1): aus der Reihenfolge der Wellen gelöst, Start nach dem Kettenumbau; die Stelle in der Liste der Welle (fünfte) ist Reihenfolge, kein Start-Trigger. Ob er vor `slice-harness-gegenproben-aus-gates-fahrstelle` läuft, hat der Nutzer nicht gesagt; es bindet ihn nichts daran. (d) **WIP:** Der Slice läuft auf eigenem Zweig (Worktree); `Verantwortlich:` trägt beim Start Person und Zweig. Er setzt `slice-harness-nachweis-uebernehmen-aufrufarten` nicht voraus; dieser setzt ihn voraus.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness docs/developer`. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; Vertrag und ADR (Commits des Architect vor dem Start) zählen nicht.

**Schätzung (Architect, 2026-10-18, für den ungeteilten Slice; Planner, nach dem Schnitt für Teil 1; Metrik wie oben; Vergleichswert `stop-hook-gegenprobe`: 1217 Zeilen für rund 150 Fälle, etwa 8 Zeilen je Fall mit Rahmen):** Der Vertrag zählt 93 Fälle in sechs Gruppen; der ungeteilte Slice lag bei 830 bis 1180 Zeilen, die Mitte (rund 1000) über der Marke 900. Der Architect empfahl den Schnitt, der Planner hat ihn nach der Regel des Nutzers (schneiden, wenn die Schätzung die Marke selbst überschreitet) gemacht.

| Teil 1 (dieser Slice, 71 Fälle) | Zeilen |
|---|---|
| `nachweis-uebernehmen.sh` (elf Gründe, Kennwörter, Hilfsverzeichnis, `-i`), Fragment, zwei Kommentare | 90 bis 140 |
| `nachweis-uebernehmen-gegenprobe` (71 Fälle mit Mutationen, Attrappe für `git`, Rahmen der Temp-Repos), Fragment | 535 bis 740 |
| `harness/README.md` (zwei Zeilen) und `docs/developer/devcontainer.md` (Abschnitt) | 40 bis 60 |
| **Zusammen** | **665 bis 940** |

Meldeschwelle 700, **Marke 900**. Die Mitte (rund 800) liegt bei 89 % der Marke, die Untergrenze darunter, die Obergrenze (940) darüber; die Architect-Schätzung für Teil 1 mit den 66 Fällen ohne `-i` (630 bis 890) liegt ganz darunter, die fünf `-i`-Fälle (35 bis 50 Zeilen) bringen die Obergrenze über die Marke. Entscheidung des Planners (Empfehlung, nicht vom Nutzer bestätigt): `-i` bleibt in Teil 1 (Grund in §1), eine weitere Teilung gibt es nicht — der Vorbehalt der Rückführung bei einer Summe über 900 steht unten. Teil 2 (`slice-harness-nachweis-uebernehmen-aufrufarten`): 22 Fälle, 165 bis 250 Zeilen. Teil 1 und Teil 2 zusammen 830 bis 1190 (die Rundung des Architect: 830 bis 1180).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen. Der Planner entscheidet dann die Naht (Vorschlag: die `-i`-Fälle zu Teil 2); der Bericht nennt die Messung.
- `in-progress` → `open` (blockiert): Der unveränderte Hook gibt bei byte-gleicher Liste nicht frei (Fall `p16-identisch` oder die Probe am Endstand zeigt es anders; dann ist die Annahme des Slice falsch, Befund an den Planner und den Architect), oder `dateiliste_erzeugen` liefert in der Quelle-Umgebung und hier für denselben Commit verschiedene Listen ohne untracked Datei (eine Quelle der Ungleichheit, die §1 nicht nennt).

## 5. Closure-Trigger

DoD vollständig; `make gates` grün auf dem Endstand; `make nachweis-uebernehmen-gegenprobe` grün mit den 71 Fällen des Teils 1; die Sendungen an `slice-harness-nachweis-uebernehmen-aufrufarten` stehen in dessen §1 und DoD (`grep` nach der Kennung dieses Slice dort); die Probe am Endstand (Verifier) mit den beiden Antworten des Hooks im Bericht; Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag geschrieben; die drei Paarungen getragen. Der Zähler von `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` ist Sache von `slice-harness-nachweis-uebernehmen-aufrufarten` (dort fährt die Gegenprobe den Hook); dieser Slice trägt keinen Beleg dazu.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des neuen Vertrags (`AGENTS.md` §3.12).** Je Randform steht die Empfehlung des Planners und, seit der Prüfung des Architect vom 2026-10-18, die Stelle der Entscheidung in [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) (und `SPEC-052` Punkte 1, 8 und 21); Namen, Exit-Codes und Kennwörter stehen dort. Die Entscheidungen sind entschieden: Der Nutzer hat [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) am 2026-10-18 angenommen (Accepted, Option B); `SPEC-062` ändert nur der Architect. Die Randformen sind mit Teil 1 und Teil 2 geschnitten (Nummern unverändert): dieser Slice trägt die Gruppen Q, G, Z, die Punkte 14 (Schreiben scheitert), 17 (nur `-i`), 19, 21, 22 und 23 bis 26; `slice-harness-nachweis-uebernehmen-aufrufarten` trägt 14 (Signale), 17 (übrige Flags), 18 und 20. Der Implementer entscheidet keine still (Randform-Rückgabe). Was aus dem Bestand entscheidbar war, steht als „Bestand“. Namen von Variable, Ziel und Dateien: Vorschlag des Planners.

**Gruppe Q — Quelle**

1. **Aufruf ohne Quelle.** `NACHWEIS_QUELLE` ungesetzt, leer: Exit 2 mit Meldung, nichts geändert. Das Verzeichnis existiert nicht oder ist eine Datei: Exit 2. Das Verzeichnis ist nicht lesbar: Exit 1. Ein Fall je Form. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 3 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
2. **Nachweis unvollständig.** `gates-passed.dateiliste` fehlt, `gates-passed.diffsha` fehlt, beide fehlen, eine von beiden ist leer: Exit 1, nichts geändert („ohne Nachweis“). Ein Fall je Form. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 5 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
3. **Liste beschädigt.** Je Schadensart von `SPEC-052` Punkt 20 (`dateiliste_pruefen`): leer, Feldzahl nicht durch 3 teilbar, unbekannte Art, `F` mit Wert, der nicht 64 Hexzeichen ist, `G` oder `X` mit Wert, doppelter Pfad: Exit 1. Ein Fall je Schadensart. Ein unzulässiger Pfad (Punkt 12) ist kein eigener Fall: Er ist nie gleich einer Liste, die `git` liefert, die Gleichheit lehnt ab. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 6 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
4. **Hash neben der Liste.** `gates-passed.diffsha` muss byte-gleich `<SHA-256 der Liste>\n` sein (die Form, die `record-gates.sh` schreibt: 64 Kleinbuchstaben-Hexzeichen und ein Zeilenumbruch). Abgelehnt, ein Fall je Form: anderer Hash; ohne Zeilenumbruch; Großbuchstaben; zusätzliche Zeile. Dies fängt auch eine halb geschriebene Quelle (`record-gates.sh` schreibt den Hash, dann ersetzt es die Liste; ein Lesen dazwischen sieht Hash neu, Liste alt). Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 7 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
5. **Quelle mit Leerzeichen oder führendem `-` im Pfad.** Wird angenommen wie jeder andere Pfad (Anführungszeichen im Skript, `--` vor dem Pfad); ein Fall je Form. Andere Sonderzeichen sind nicht zugesagt. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 8 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).

**Gruppe G — Gleichheit mit dem Baum, der hier liegt**

6. **Liste ungleich.** Abgelehnt, Exit 1, die Meldung nennt mindestens einen abweichenden Pfad (Zahl und Wortlaut darüber hinaus nicht zugesagt). Ein Fall je Form: anderer Inhaltshash; ein Pfad zu viel in der Quelle; ein Pfad zu wenig; Art `F` gegen `L`; anderes Link-Ziel; untracked Datei im Ziel; geänderte getrackte Datei im Ziel; im Ziel gelöschte getrackte Datei (Art `G`). Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 9 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
7. **Gleichheit trotz Sonderpfad.** Ein Pfad mit Zeilenumbruch und einer mit Leerzeichen in beiden Listen: übernommen (die Liste ist NUL-getrennt, `SPEC-052` Punkt 11). Zwei Fälle. Entschieden: Bestand für die Form, die Fälle sind neu ([`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 9).
8. **Baum ohne Liste.** Scheitert `git` im Ziel oder ist die frische Liste leer: Exit 1, nichts geändert; kein Git-Arbeitsbaum: Exit 2. Zwei Fälle. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 11 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
9. **Nicht aufgelöster Index-Eintrag im Ziel.** Das Werkzeug übernimmt bei gleicher Liste (die Liste trägt keine Index-Stufen), der Hook sperrt weiter (Punkt 17, vor der Identität). Ein Fall mit beiden Antworten; die Meldung des Werkzeugs sagt dazu nichts zu. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 12 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
10. **Nicht Teil der Liste, damit nicht geprüft** (Bestand, ohne Fall; `SPEC-052` Punkte 11 und 18): Dateimodi, ignorierte Dateien, Zeitstempel. Ein anderes Auschecken von Zeilenenden (`core.autocrlf`) gibt eine andere Liste und damit die Ablehnung, nie eine falsche Freigabe. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 20.

**Gruppe Z — Ziel**

11. **Erfolg ersetzt.** Nach Erfolg sind beide Dateien byte-gleich denen der Quelle; ein älterer Nachweis (gültig, beschädigt, fehlend) und ein fehlendes `.harness/state/` ändern das Ergebnis nicht (das Verzeichnis wird angelegt). Vier Fälle. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 13 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
12. **Ablehnung ändert nichts.** Alte Dateien bleiben byte-gleich, `.harness/state/` trägt keine Hilfsdatei (auch nicht nach Ablehnung in Q, G und Z); zwei Fälle (Dateien, Hilfsdateien), je Gruppe geprüft. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 14 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
13. **Wiederholung und Alter.** Ein zweites Übernehmen derselben Quelle endet mit Exit 0 und gleichen Bytes. Ein Nachweis kennt keinen Zeitpunkt; ein „älterer“ Nachweis desselben Baums ist so gültig wie ein neuerer, ein Nachweis eines anderen Baums ist ungleich. Ein Fall. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 13 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
14. **Schreiben scheitert.** `.harness/state` ist eine Datei, oder eine der beiden Zieldateien ist ein Verzeichnis: Exit 1, alter Nachweis unverändert, keine Hilfsdatei. Drei Fälle (`p15-state-ist-datei`, `p15-liste-ist-verzeichnis`, `p15-hash-ist-verzeichnis`). Ein Verzeichnis ohne Schreibrecht ist unter `root` nicht herstellbar (akzeptiertes Negativ, `SPEC-062` Punkt 15). Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 15 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted). **`SIGINT` und `SIGTERM` zwischen Prüfung und Ersetzen (vier Fälle) liefert `slice-harness-nachweis-uebernehmen-aufrufarten`** (Sendung, §8).
15. **Quelle ist das Ziel.** Gleicher Pfad, mit `./` und mit Schrägstrich am Ende, und ein Symlink auf das Verzeichnis des Ziels: Exit 2, nichts geändert (Vergleich der Dateien mit `-ef`). Drei Fälle. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 4 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).
16. **Quelle im Arbeitsbaum des Ziels.** Liegt die Quelle als untracked Verzeichnis im Baum, sind ihre Dateien Teil der Liste und der Vergleich schlägt fehl (Fall aus 6, eigene Zeile in der Meldung nicht zugesagt); liegt sie unter `.harness/state/` (ignoriert), wird sie angenommen. Zwei Fälle. Das Rezept (Liefer-Punkt 3) nennt den zweiten Ort. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 10 (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted).

**Gruppe A — Aufrufarten (`AGENTS.md` §3.12, Aufrufarten des Treibers)**

17. **Make-Flags, hier `-i`.** `-i`: Die Einzelschritte erben `-i`, ein abgelehntes Übernehmen sähe aus wie Erfolg, deshalb lehnt das Skript `-i` in `MAKEFLAGS` ab (Exit 2, `make-i`, nichts geändert; ein vorhandener Nachweis entfällt nicht); maßgeblich ist `MAKEFLAGS` auch aus der Umgebung, in jeder Schreibweise und zusammen mit `-k` oder `-j`. Fünf Fälle: `p16-i`, `p16-i-lang` (`--ignore-errors`), `p16-i-env`, `p16-i-k`, `p16-i-j`. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 16 (b) (Architect, 2026-10-18; entschieden, [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) Accepted). **`-n`, `-q`, `-t` (auch aus der Umgebung), die Kombinationen mit `-i`, `-k`, `-j` und das Langflag liefert `slice-harness-nachweis-uebernehmen-aufrufarten`** (Punkt 16 (a), (c), (d); 12 Fälle); das Rezept ohne `$(MAKE)` und ohne `+` steht im Fragment dieses Slice, geprüft wird es dort (Sendung, §8).
18. **Aufruf aus einem Unterverzeichnis** (zwei Fälle, `SPEC-062` Punkt 17): liefert `slice-harness-nachweis-uebernehmen-aufrufarten` (Sendung, §8). Das Skript dieses Slice arbeitet im Normalfall aus der Wurzel; was es sonst tut, sagt dieser Slice nicht zu (§1).
19. **Gleichzeitige Läufe** (zwei Übernahmen, eine Übernahme zusammen mit `make gates`, mit einer Änderung des Baums oder mit dem Hook): nicht zugesagt (`SPEC-052` Punkt 8, Aufruf-Voraussetzung nacheinander), ohne Fall. Ändert sich der Baum nach der Übernahme, vergleicht der Hook neu und sperrt oder gibt nach seinen eigenen Punkten frei; das Werkzeug kann nie öfter freigeben als der Hook, weil dieser den Baum liest, nicht das Werkzeug. **Aufruf-Voraussetzung:** Das Werkzeug sagt keine Ordnung über mehrere Aufrufe zu. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 19 (ohne Fall).

**Gruppe H — Kopplung an den Hook** (zugeordnet: `slice-harness-nachweis-uebernehmen-aufrufarten`)

20. **Der unveränderte Hook liest das Ergebnis** (vier Fälle `p18-…`, `SPEC-062` Punkt 18, Mutation: das Werkzeug bildet die Liste selbst, färbt `p9-untracked` des Teils 1 und die Fälle der Kopplung): liefert `slice-harness-nachweis-uebernehmen-aufrufarten` (Sendung, §8). Dieser Slice liest die beiden Antworten des Hooks in der Probe am Endstand von Hand (DoD) und sagt sie nicht zu.

**Gruppe S — Grenze, offen genannt (ohne Fall)**

21. **Jeder lokale Nachweis ist von Hand zu schreiben** (`dateiliste_erzeugen > .harness/state/gates-passed.dateiliste`, `make record-gates` einzeln, `enforce.mk`); das gilt heute schon. Das Werkzeug ist der einzige sanktionierte Weg und prüft Gleichheit, Unversehrtheit und Zusammengehörigkeit; es kann nicht prüfen, dass die Quelle `make gates` fuhr, mit denselben Gates, im selben Image-Stand und mit allen Gegenproben, die dann gelten. Akzeptiertes Negativ; Grund: Der Nachweis trägt keine Herkunft, und eine Herkunft ohne Signatur wäre von Hand ebenso zu schreiben. Die Grenze steht im Vertrag ([`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 21) und im Rezept (§3.11).
22. **Der Nachweis gilt einem Baum, nicht einem Zweig.** Ein Merge, der den Baum des Hauptzweigs ändert (Merge-Commit mit anderem Inhalt), macht die Listen ungleich; die Regel an den Aufrufer lautet, den Zweig vor dem Merge auf den Stand des Hauptzweigs zu bringen und dort `make gates` zu fahren, oder `make gates` im Hauptbaum. Das Werkzeug lehnt ab, es rät nicht. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 22 (ohne Fall).

**Vom Architect ergänzt (Prüfung vom 2026-10-18; §6 des Plans nannte sie nicht, `AGENTS.md` §3.12):**

23. **Reihenfolge der Gründe.** Elf Gründe in fester Reihenfolge, je Nachbarschaft ein Fall. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 2.
24. **Plattform** (`LH-QA-03`). Linux, bash ≥ 4, GNU-Werkzeuge; macOS und Windows nicht zugesagt. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Kopf.
25. **Hilfsverzeichnis und `TMPDIR`.** Die Quelle wird in ein Hilfsverzeichnis unter `${TMPDIR:-/tmp}` gelesen, nicht neben das Ziel; Ablehnung legt nichts im Arbeitsbaum an; ein `TMPDIR` im Baum endet in `ungleich`, nie in einer Übernahme. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkte 14 und 15.
26. **Quelle und Stand der schnellen Kette bleiben unverändert**, im Erfolg wie bei Ablehnung. Entschieden: [`SPEC-062`](../../../../harness/vertraege/SPEC-062-nachweis-uebernehmen.md) Punkt 14.

### Risiken

- **Der Nachweis gilt Umgebungen, die niemand vergleicht** (Randform 21; ein Gate, das nur in einer Umgebung rot wird, bleibt in der Quelle unbemerkt). — **Ausgang (vorgesehen):** weiter offen → neues `BEO-REPO/nachweis-uebernommen-aus-fremder-umgebung` (Beleg bei der Closure).
- **Der Hauptbaum ist selten gleich der Quelle** (unversionierte inerte Änderung, untracked Datei; dann lehnt das Werkzeug ab, und die Sperre bleibt). Nach einer Übernahme verlangt die nächste inerte Änderung im Hauptbaum wieder den Stand der schnellen Kette, also Docker; der Nutzer hat entschieden, das nicht zu ändern (Frage 3, Option (a)). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (Vermerk mit der Zahl der Ablehnungen und Gründe aus der Probe).
- **Die Gegenprobe hängt an `GATE_CHECKS`, weil die Listen von `slice-harness-gegenproben-aus-gates` fehlen.** — **Ausgang (vorgesehen):** entfallen, wenn beim Start `slice-harness-gegenproben-aus-gates` in `done/` liegt und `grep -l nachweis-uebernehmen-gegenprobe harness/mk/*.mk` das Fragment in der Liste zeigt (nicht an `GATE_CHECKS`); sonst eingetreten → Rückführung nach `open` (§4).
- **Das Rezept für den Weg der Dateien aus dem Container (`ssh`, `tar`) ist keine geprüfte Zusage** (`AGENTS.md` §3.11). — **Ausgang (vorgesehen):** entfallen, wenn der Abschnitt das Rezept als nicht geprüft kennzeichnet und die Gegenprobe den Fall mit einer lokal durch eine `tar`-Pipe angelegten Kopie führt (Review prüft den Wortlaut).
- **Die Obergrenze der Schätzung liegt über der Marke** (665 bis 940, Marke 900; die Mitte rund 800); sie hängt an der Zahl der Fälle der Gegenprobe. — **Ausgang (vorgesehen):** entfallen, wenn die Messung nach der Metrik unter 900 bleibt, mit der Zahl in der Closure-Notiz; darüber *eingetreten* → Rückführung *zu groß* (§4) und Naht des Planners.
- **Teil 1 liefert ein Skript, dessen Verhalten bei Signal, bei `-n`, `-q`, `-t`, aus einem Unterverzeichnis und gegenüber dem Hook ohne Fall und ohne Mutation bleibt, bis `slice-harness-nachweis-uebernehmen-aufrufarten` schließt** (`AGENTS.md` §3.10, §3.11). — **Ausgang (vorgesehen):** entfallen, mit Begründung: §1 verbietet jede Zusage dazu in Kommentar, Rezept und Gate-Index, das Review prüft den Wortlaut, und `slice-harness-nachweis-uebernehmen-aufrufarten` schließt direkt dahinter (Start-Trigger: dieser Slice in `done/`).
- **Der Nehmer `slice-harness-nachweis-uebernehmen-aufrufarten` nimmt die Sendungen nicht an** oder zieht sie nicht nach, wenn dieser Slice den Vertrag oder die Namen ändert (`AGENTS.md` §3.13). — **Ausgang (vorgesehen):** entfallen mit `grep` nach der Kennung dieses Slice in §1 und §6 des Nehmers (Closure-Trigger).

**Fragen an den Nutzer** (jede mit Empfehlung des Planners; eine Empfehlung ist keine Entscheidung):

1. **Stelle in der Welle — entschieden** (Entscheidung des Nutzers vom 2026-10-18, so vom Koordinator übermittelt, der Wortlaut liegt dem Planner nicht vor): aus der Reihenfolge der Wellen lösen, Start nach dem Kettenumbau (`slice-harness-gegenproben-aus-gates` in `done/`). Offen bleibt, ob er vor `slice-harness-gegenproben-aus-gates-fahrstelle` rückt (Empfehlung des Planners: nein, kein Grund aus dem Auftrag; ändert sich nichts, solange es ihn nicht bindet).
2. **Quelle nur als Verzeichnis, Transport als Rezept** (`ssh` mit `tar` oder `docker cp`, in der Wartungs-Doku) statt eines Transports im Werkzeug. Empfehlung: ja (§1, zweiter Ausschluss); steht so in `SPEC-062` (Schnittstelle).
3. **Zweite Sperre des Hauptbaums — entschieden** (Entscheidung des Nutzers vom 2026-10-18, so vom Koordinator übermittelt): Option (a), nichts ändern; die Ablehnungen der Probe werden gemessen. Option (b) (Übernahme des Standes der schnellen Kette) und (c) (Hook gibt inerte Änderungen frei, ADR nach `AGENTS.md` §3.6) sind nicht gewählt; kein Folge-Slice angelegt.
4. **Pflicht in den Rollen-Anweisungen** (wann das Werkzeug zu rufen ist, etwa in `implement-slice` nach dem Merge). Empfehlung: erst nach der Lieferung und einer Probe, dann als Satz in der Fahrstelle oder einem Folge-Slice.
5. **ADR — entschieden:** [ADR-0049](../../adr/0049-nachweis-aus-anderer-umgebung-uebernehmen.md) ist `Accepted` (Entscheidung des Nutzers vom 2026-10-18, Option B).
6. **Folge-Slice für die Umgebung des Hooks** (`BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung`, dritter Beleg): gehört seit dem Schnitt zu `slice-harness-nachweis-uebernehmen-aufrufarten` (Gruppe H). Empfehlung: Der Planner entscheidet bei der Closure dort.

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
- **Beobachtungs-Register (`../observations/`):** Vermerk in `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`; `BEO-REPO/nachweis-uebernommen-aus-fremder-umgebung/` neu, Beleg `evidence/slice-harness-nachweis-uebernehmen.md` (bei der Closure)
- **Folge-Slices:** `slice-harness-nachweis-uebernehmen-aufrufarten` (Teil 2, angelegt beim Schnitt)
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (1×; der Slice ist die Gegenmaßnahme für den Nachweis, kein Auftreten), `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (**2×**; die Gruppe H, die den Hook in einer Kopie fährt und seine Antwort liest, liefert seit dem Schnitt `slice-harness-nachweis-uebernehmen-aufrufarten`; dieser Slice fährt den Hook nicht in der Gegenprobe und trägt keinen Beleg; den Zähler prüft der Planner dort bei der Closure), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (DoD 2 trägt je Element einen Fall und eine Mutation), `BEO-REPO/spec-randform-erst-im-review-entschieden` (§6 nennt die Randformen vor dem Code), `BEO-REPO/slice-waechst-durch-uebernahmen` (der Slice übernimmt nichts und trägt keine Sendung ein), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groessenprobe-ohne-messbaren-stand` (§4 rechnet in der Metrik mit dem Vergleichswert der Gegenprobe des Hooks), `BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet` (jede Wahl der Stelle, des Namens und der Option steht als Empfehlung des Planners; entschieden ist nur der Auftrag des Nutzers). Keiner der übrigen Treffer erreicht mit diesem Slice 3× neu.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-18, im Commit des Schnitts, mit den Sendungen an `slice-harness-nachweis-uebernehmen-aufrufarten`):** Liefer-Punkte 3 (Werkzeug | Gegenprobe | Gate-Index und Rezept); der Slice liegt an der Grenze, ein vierter Punkt hat keinen Platz. Schichten 2: Harness-Werkzeug mit Gegenprobe und Gate-Index (eine Schicht) und Wartungs-Doku `docs/developer/` (eine Schicht); `AGENTS.md` bleibt unberührt, Planung, Vertrag und ADR sind keine Schicht. Schätzung in der Metrik von §4: 665 bis 940 (Marke 900, Meldeschwelle 700, die Obergrenze über der Marke, die Mitte rund 800; Risiko in §6).
**Sendungen an `slice-harness-nachweis-uebernehmen-aufrufarten`** (im selben Commit in dessen §1, DoD, §6 und §8 eingetragen, mit der Kennung dieses Slice): die Fälle und Mutationen der Punkte 15 (vier Fälle mit Signal), 16 (a), (c), (d) (12 Fälle), 17 (zwei Fälle) und 18 (vier Fälle) von `SPEC-062`; dazu die **Bindungen**: das Fragment dieses Slice trägt im Rezept weder `$(MAKE)` noch `+` (Punkt 16 (a) prüft es dort), das Skript dieses Slice arbeitet aus der Wurzel des Repos (Punkt 17), und die Kennwörter und Exit-Codes sind die von Punkt 1 (Fälle in Teil 2 lesen sie). Nachgezählt im Nehmer: Liefer-Punkte 3, Schicht 1, Schätzung 165 bis 250 (§8 dort). Voraussetzung, keine Sendung, an die dieser Slice sich bindet: `slice-harness-gegenproben-aus-gates` (Listen und Namensregel für die Gegenprobe, Wortlaut von `SPEC-052` Punkt 9), der Geber liegt in `done/`, wenn dieser Slice startet (Start-Trigger in §4).
