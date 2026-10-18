# Slice slice-harness-nachweis-uebernehmen: Ein Werkzeug übernimmt den Nachweis des vollen Laufs aus einer Quelle, wenn die Dateiliste des Baums, der hier liegt, mit der des Nachweises übereinstimmt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-gegenproben-und-mutanten, fünfte Stelle der Reihenfolge (vierte bis zum Schnitt von `slice-harness-gegenproben-aus-gates-fahrstelle` am 2026-10-18), hinter `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` und vor `slice-harness-gegenprobe-baumkopie` (**Empfehlung des Planners vom 2026-10-18, nicht entschieden**: ob der Slice aus der Reihenfolge der Wellen gelöst wird und ob er vor `-fahrstelle` rückt, sind Frage 1 in §6; bis zur Antwort gilt zusätzlich der Start-Trigger der Welle). Begründung der Empfehlung: gleiche Familie wie die drei Slices davor (Gate-Kette, Stop-Hook, Gegenproben; er setzt deren Listen und Namensregel voraus), und der Auftrag des Nutzers vom 2026-10-18 („Keine blockierenden Aktionen“ im Hauptbaum) macht ihn dringend. Angelegt am 2026-10-18 auf Auftrag des Nutzers („Nachweis aus dem Container oder Worktree übernehmen — das hört sich pragmatisch an“). Der Slice liegt in `open/`: Er ist nicht priorisiert, den Anspruch trägt der Nutzer.

**Bezug:** [LH-QA-07](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) nur insoweit, als das Werkzeug und seine Gegenprobe dessen Messmethode (4) tragen; kein Anspruch des Produkts, der Slice ändert die Gate-Kette des Harness. Aktive ADR: [ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md) (der Hook glaubt den Nachweis des vollen Laufs); eine neue ADR schreibt der Architect vor dem Code (Frage 5). Bindung an das Register: `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (Gegenmaßnahme: der Nachweis reist mit dem Stand), `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (2×, §8).

**Berührte Spec-Stellen:** neuer Vertrag `SPEC-062` in `harness/vertraege/` (Nummer: nächste freie, Architect; belegt sind `SPEC-058` und `SPEC-060`, `SPEC-061` ist an `slice-harness-gegenproben-aus-gates-fahrstelle` vergeben, `SPEC-059` für `slice-harness-gegenprobe-baumkopie` vorgemerkt); [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) Punkte 1, 11, 16, 19 und 21 (der Nachweis bekommt einen zweiten Schreiber; der Hook bleibt unverändert) und Punkt 12 (ein neues Skript unter `tools/` ist nicht inert); `SPEC-060` Teil 1 (die Gegenprobe des Werkzeugs steht nach der Namensregel in einer Liste), nur gelesen.

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

**Ziel:** `make nachweis-uebernehmen` (Name: Vorschlag, Architect) legt den Nachweis des vollen Laufs, den eine andere Umgebung geschrieben hat (Worktree, Klon im Dev-Container), im Arbeitsbaum ab, in dem es läuft, **genau dann**, wenn (a) die Dateiliste dieses Baums, frisch mit `dateiliste_erzeugen` berechnet, byte-gleich mit der Dateiliste der Quelle ist, (b) die Liste der Quelle unbeschädigt ist (`dateiliste_pruefen`) und (c) `gates-passed.diffsha` der Quelle der SHA-256 genau dieser Liste ist. Sonst lehnt es ab und ändert nichts. Der Hook bleibt unverändert: Er liest nur die Liste (`.claude/hooks/stop-require-gates.sh`, `cmp -s`, `SPEC-052` Punkt 16), und eine byte-gleiche Liste ist seine Freigabe-Bedingung; neu ist ein zweiter Schreiber des Standes, kein neuer Leser. Das Werkzeug ist kein Gate, hängt nicht an `GATE_CHECKS` und schreibt nichts außer den zwei Dateien des Nachweises.

**Befund am Stand `0a59a42f` (gelesen, nicht geändert):**

- **Der Nachweis** sind zwei Dateien in `.harness/state/` (ignoriert, also je Arbeitsbaum eigen, `bin/` darunter ebenso): `gates-passed.dateiliste` (Ausgabe von `dateiliste_erzeugen`: je Pfad Art, Wert, Pfad) und `gates-passed.diffsha` (SHA-256 der Liste, eine Zeile). `tools/harness/record-gates.sh` schreibt beide als letztes Prerequisite von `make gates`, also nur nach grünen Gates, und entfernt beide bei `-i` in `MAKEFLAGS`. **Die Dateien tragen weder Zeitpunkt noch Commit noch Herkunft**; ein „grüner Lauf“ ist an nichts weiter erkennbar als daran, dass beide da sind und zusammenpassen (§6, Gruppe S).
- **Der Hook liest nur die Liste.** `gates-passed.diffsha` liest kein Skript außer der Gegenprobe (`stop-hook-gegenprobe`, Fälle `p1-…`); die Zusammengehörigkeit prüft deshalb das neue Werkzeug, nicht der Hook.
- **Was die Liste enthält:** getrackte und untracked Pfade ohne ignorierte, je Eintrag Art (`F`, `L`, `G`, `X`) und SHA-256 bzw. Link-Ziel (`SPEC-052` Punkt 11). **Nicht** darin: `.harness/state/`, `dist/` und alles, was `.gitignore` ausblendet; Dateimodi (akzeptiertes Negativ, Punkt 11); Zeitstempel; `.git` (im Worktree eine Datei). Zwei Umgebungen desselben Commits unterscheiden sich in der Liste also nur durch **untracked, nicht ignorierte Dateien**, durch Inhalte, die `git` beim Auschecken umschreibt (`core.autocrlf`, `.gitattributes`), und durch lokale Ignore-Regeln (`.git/info/exclude`, `core.excludesFile`), die eine untracked Datei in der einen Umgebung ausblenden und in der anderen nicht.
- **Ein roter Lauf lässt einen älteren Nachweis stehen** (`SPEC-052` Punkt 21, akzeptiertes Negativ). Das schadet dem Werkzeug nicht: Es vergleicht die Liste der Quelle mit dem Baum, der hier liegt, nicht mit dem Stand der Quelle; ein älterer Nachweis ist für den Baum, den er beschreibt, weiter wahr, und für jeden anderen ist die Liste ungleich.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Änderung am Hook, an `dateiliste_erzeugen`, `dateiliste_pruefen`, der Positivliste oder an `record-gates.sh` (außer dem Kommentar, der den Schreiber nennt) — Bestand bleibt bewusst stehen: Dass der Hook bei byte-gleicher Liste freigibt, ist `SPEC-052` Punkt 16 mit Fall `p16-identisch`; der Slice stützt sich darauf, statt es zu ändern, und die Gegenprobe des Werkzeugs fährt den unveränderten Hook (§6, Gruppe H).
- Ein Transport im Werkzeug (`ssh`, `docker cp`, `git`-Abruf) — anderer Vorgang und eine Host-Abhängigkeit über `bash`, `git`, `docker` und `make` hinaus (`harness/conventions.md`); er brächte Verbindungsabbruch, Teilübertragung und Host-Schlüssel als Randformen mit. Die Quelle ist ein Verzeichnis mit den zwei Dateien; wie sie dorthin kommen, steht als Rezept in der Wartungs-Doku (Liefer-Punkt 3). Frage 2.
- Die Übernahme des Standes der schnellen Kette (`.harness/state/gates-schnell-passed.diffsha`) und eine Freigabe inerter Änderungen ohne diesen Stand — anderer Vorgang: Die erste Übernahme würde eine zweite Gleichheit prüfen (Hash des ganzen Baums), die zweite ändert, was der Hook glaubt ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), Lockerung nach `AGENTS.md` §3.6). Beides ist Entscheidung des Nutzers (Frage 3), kein Folge-Slice ist angelegt.
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

- [ ] **Liefer-Punkt 1 (Werkzeug):** `make nachweis-uebernehmen NACHWEIS_QUELLE=<verzeichnis>` (Namen: Vorschlag des Planners, laut `SPEC-062`) übernimmt `gates-passed.dateiliste` und `gates-passed.diffsha` aus dem Verzeichnis in `.harness/state/` des Arbeitsbaums, in dem es läuft, genau unter den drei Bedingungen aus §1, und lehnt sonst ab, ohne eine Datei zu ändern oder eine Hilfsdatei zu hinterlassen; jede Entscheidung je Randform aus §6 (Gruppen Q, G, Z, A) steht im Vertrag `SPEC-062`, vor dem ersten Code-Commit; das Werkzeug benutzt `dateiliste_erzeugen` und `dateiliste_pruefen` aus `tools/harness/dateiliste.sh` und bildet die Liste nicht nach. Fragment `harness/mk/nachweis-uebernehmen.mk` und Skript `tools/harness/nachweis-uebernehmen.sh`; das Ziel steht nicht an `GATE_CHECKS`. Die Kommentare in `record-gates.sh` und `dateiliste.sh`, die den Schreiber des Nachweises nennen, nennen auch das Werkzeug.
- [ ] **Liefer-Punkt 2 (Gegenprobe):** `make nachweis-uebernehmen-gegenprobe` (Fragment unter `harness/mk/`, in der Liste, die die Namensregel von `SPEC-060` Teil 1 vorschreibt, und ausdrücklich nicht an `GATE_CHECKS`) trägt je Randform und je Element ihrer Aufzählungen aus §6 einen Fall und eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10), und je Fall der Gruppe H einen Lauf des unveränderten Hooks in einer Kopie; die Gegenprobe verlangt je Mutation genau die Menge roter Fälle, die der Vertrag nennt (Namen der Fälle tragen die Nummer des Punktes). Temp-Repos, kein Docker, kein Netz; der Arbeitsbaum bleibt unberührt.
- [ ] **Liefer-Punkt 3 (Gate-Index und Rezept):** `harness/README.md` führt die Zeile des Werkzeugs (Tabelle *Werkzeuge — kein Gate*) und die der Gegenprobe (§Sensors, Bindung `SPEC-062`); `docs/developer/devcontainer.md` trägt einen Abschnitt, der das Werkzeug mit dem Worktree und mit dem Klon im Dev-Container als Quelle zeigt (der Weg der zwei Dateien aus dem Container in ein Verzeichnis ist als Rezept gekennzeichnet; die Gegenprobe prüft eine lokale Kopie der Dateien, nicht `ssh`, `AGENTS.md` §3.11).
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `tools/` und `harness/` ändern); `make nachweis-uebernehmen-gegenprobe` einzeln grün, die Dauer im Bericht; `make gates-schnell-gegenprobe` grün, falls die Deckung der Ketten die neue Gegenprobe zählt (laut `SPEC-060` Teil 1).
- [ ] **Probe am Endstand (Verifier):** Quelle ist der Worktree des Slice mit grünem `make gates` auf dem Endstand, Ziel ein frischer Worktree desselben Commits ohne Nachweis; der unveränderte Hook antwortet im Ziel vorher `block` (Punkt 19), das Werkzeug endet mit Exit 0, der Hook antwortet nachher `approve`; mit einer untracked Datei im Ziel lehnt das Werkzeug ab und der Hook bleibt bei `block`. Der Hauptbaum wird dafür nicht gebraucht.
- [ ] Beobachtungs-Register: Vermerk (kein Zähler) in `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (Gegenmaßnahme geliefert); `BEO-REPO/nachweis-uebernommen-aus-fremder-umgebung` neu angelegt (Risiko 1 und 2, Beleg `evidence/slice-harness-nachweis-uebernehmen.md`, bei der Closure).
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
| `harness/vertraege/SPEC-062-nachweis-uebernehmen.md` (Nummer und Dateiname: Architect), `harness/vertraege/README.md`, `harness/vertraege/SPEC-052-gate-kette.md` Punkte 1 und 21 | Vertrag (Architect, liegt vor dem Code) | Entscheidung je Randform aus Gruppen Q, G, Z, A, H; der zweite Schreiber des Nachweises steht in `SPEC-052` (Punkt 21 sagt heute, wer ihn schreibt) |
| neue ADR (Architect, Annahme: Entscheidung des Nutzers, Frage 5) | Entscheidung | der Hook glaubt einen Nachweis, den auch eine andere Umgebung geschrieben haben kann; Alternativen und Gründe, `Schärft:` zeigt auf `SPEC-062` (`AGENTS.md` §3.8) |
| `tools/harness/nachweis-uebernehmen.sh`, `harness/mk/nachweis-uebernehmen.mk` | neu | liest die Quelle in eine Hilfsdatei neben dem Ziel, prüft Liste und Hash, vergleicht mit `dateiliste_erzeugen`, ersetzt beide Dateien (Reihenfolge wie `record-gates.sh`: Hash, dann Liste); lehnt `-i` in `MAKEFLAGS` ab |
| `tools/harness/record-gates.sh`, `tools/harness/dateiliste.sh` | update (nur Kommentare) | nennen den zweiten Schreiber (`AGENTS.md` §3.9) |
| `tools/harness/nachweis-uebernehmen-gegenprobe.sh`, `harness/mk/nachweis-uebernehmen-gegenprobe.mk` | neu | Temp-Repos mit Kopien von Skript, `dateiliste.sh` und Hook; je Randform ein Fall und eine Mutation; Attrappen für `git`, `cmp`, `mv` bei Abbruch und Fehler |
| `harness/README.md`, `docs/developer/devcontainer.md` | update | Gate-Index (zwei Zeilen) und Rezept (Worktree, Klon im Container) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `Verantwortlich:` setzt dieser Übergang; der Nutzer priorisiert und trägt den Anspruch. **Vorbedingungen:** (a) Vertrag `SPEC-062` und ADR liegen vor, mit der Entscheidung je Randform aus §6 (`AGENTS.md` §3.12); der Architect schreibt sie nach dem Start-Trigger (b), weil `SPEC-052` Punkt 9 dann den engeren Wortlaut trägt. (b) **Start-Trigger:** `slice-harness-gegenproben-aus-gates` liegt in `done/` (Listen und Namensregel für die Gegenprobe stehen; sonst hinge sie an `GATE_CHECKS` und verlängerte `make gates`). (c) **Reihenfolge:** Frage 1. (d) **WIP:** Der Slice läuft auf eigenem Zweig (Worktree); `Verantwortlich:` trägt beim Start Person und Zweig.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness docs/developer`. Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; Vertrag und ADR (Commits des Architect vor dem Start) zählen nicht.

**Schätzung (Planner, vorläufig; der Architect rechnet bei der Prüfung neu; Metrik wie oben; Gegenprobe-Code gerechnet mit dem Vergleichswert `stop-hook-gegenprobe`: 1217 Zeilen für rund 150 Fälle, etwa 8 Zeilen je Fall mit Rahmen):**

| Teil | Zeilen |
|---|---|
| `nachweis-uebernehmen.sh`, Fragment, zwei Kommentare | 90 bis 140 |
| `nachweis-uebernehmen-gegenprobe` (rund 65 Fälle in den Gruppen Q 20, G 13, Z 15, A 14, H 3, mit Mutationen, Attrappen und Rahmen) | 380 bis 600 |
| `harness/README.md` (zwei Zeilen) und `docs/developer/devcontainer.md` (Abschnitt) | 25 bis 45 |
| **Zusammen** | **495 bis 785** |

Meldeschwelle 700, **Marke 900**; die Obergrenze liegt über der Meldeschwelle und unter der Marke, deshalb ist der Slice nicht geschnitten. Zählt der Architect mehr als 65 Fälle oder mehr als vier Gruppen, rechnet er neu; liegt die neue Schätzung über 900, schneidet der Planner (Vorschlag: Werkzeug mit den Gruppen Q, G, Z getrennt von den Aufrufarten A und der Kopplung H).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen.
- `in-progress` → `open` (blockiert): Der unveränderte Hook gibt bei byte-gleicher Liste nicht frei (Fall `p16-identisch` oder die Probe am Endstand zeigt es anders; dann ist die Annahme des Slice falsch, Befund an den Planner und den Architect), oder `dateiliste_erzeugen` liefert in der Quelle-Umgebung und hier für denselben Commit verschiedene Listen ohne untracked Datei (eine Quelle der Ungleichheit, die §1 nicht nennt).

## 5. Closure-Trigger

DoD vollständig; `make gates` grün auf dem Endstand; `make nachweis-uebernehmen-gegenprobe` grün; die Probe am Endstand (Verifier) mit den beiden Antworten des Hooks im Bericht; Review- und Verifikations-Report liegen vor; Closure-Notiz mit Lerneintrag geschrieben; die drei Paarungen getragen. Bei der Closure prüft der Planner den Zähler von `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (§8).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des neuen Vertrags (`AGENTS.md` §3.12).** Je Randform steht die Empfehlung des Planners; entschieden ist keine, der Architect legt sie vor dem Code in `SPEC-062` (und `SPEC-052` Punkt 21), und der Implementer entscheidet keine still (Randform-Rückgabe). Was aus dem Bestand entscheidbar war, steht als „Bestand“. Namen von Variable, Ziel und Dateien: Vorschlag des Planners.

**Gruppe Q — Quelle**

1. **Aufruf ohne Quelle.** `NACHWEIS_QUELLE` ungesetzt, leer: Exit 2 mit Meldung, nichts geändert. Das Verzeichnis existiert nicht oder ist eine Datei: Exit 2. Das Verzeichnis ist nicht lesbar: Exit 1. Ein Fall je Form. Entschieden: offen.
2. **Nachweis unvollständig.** `gates-passed.dateiliste` fehlt, `gates-passed.diffsha` fehlt, beide fehlen, eine von beiden ist leer: Exit 1, nichts geändert („ohne Nachweis“). Ein Fall je Form. Entschieden: offen.
3. **Liste beschädigt.** Je Schadensart von `SPEC-052` Punkt 20 (`dateiliste_pruefen`): leer, Feldzahl nicht durch 3 teilbar, unbekannte Art, `F` mit Wert, der nicht 64 Hexzeichen ist, `G` oder `X` mit Wert, doppelter Pfad: Exit 1. Ein Fall je Schadensart. Ein unzulässiger Pfad (Punkt 12) ist kein eigener Fall: Er ist nie gleich einer Liste, die `git` liefert, die Gleichheit lehnt ab. Entschieden: offen.
4. **Hash neben der Liste.** `gates-passed.diffsha` muss byte-gleich `<SHA-256 der Liste>\n` sein (die Form, die `record-gates.sh` schreibt: 64 Kleinbuchstaben-Hexzeichen und ein Zeilenumbruch). Abgelehnt, ein Fall je Form: anderer Hash; ohne Zeilenumbruch; Großbuchstaben; zusätzliche Zeile. Dies fängt auch eine halb geschriebene Quelle (`record-gates.sh` schreibt den Hash, dann ersetzt es die Liste; ein Lesen dazwischen sieht Hash neu, Liste alt). Entschieden: offen.
5. **Quelle mit Leerzeichen oder führendem `-` im Pfad.** Wird angenommen wie jeder andere Pfad (Anführungszeichen im Skript, `--` vor dem Pfad); ein Fall je Form. Andere Sonderzeichen sind nicht zugesagt. Entschieden: offen.

**Gruppe G — Gleichheit mit dem Baum, der hier liegt**

6. **Liste ungleich.** Abgelehnt, Exit 1, die Meldung nennt mindestens einen abweichenden Pfad (Zahl und Wortlaut darüber hinaus nicht zugesagt). Ein Fall je Form: anderer Inhaltshash; ein Pfad zu viel in der Quelle; ein Pfad zu wenig; Art `F` gegen `L`; anderes Link-Ziel; untracked Datei im Ziel; geänderte getrackte Datei im Ziel; im Ziel gelöschte getrackte Datei (Art `G`). Entschieden: offen.
7. **Gleichheit trotz Sonderpfad.** Ein Pfad mit Zeilenumbruch und einer mit Leerzeichen in beiden Listen: übernommen (die Liste ist NUL-getrennt, `SPEC-052` Punkt 11). Zwei Fälle. Entschieden: Bestand für die Form, die Fälle sind neu.
8. **Baum ohne Liste.** Scheitert `git` im Ziel oder ist die frische Liste leer: Exit 1, nichts geändert; kein Git-Arbeitsbaum: Exit 2. Zwei Fälle. Entschieden: offen.
9. **Nicht aufgelöster Index-Eintrag im Ziel.** Das Werkzeug übernimmt bei gleicher Liste (die Liste trägt keine Index-Stufen), der Hook sperrt weiter (Punkt 17, vor der Identität). Ein Fall mit beiden Antworten; die Meldung des Werkzeugs sagt dazu nichts zu. Entschieden: offen.
10. **Nicht Teil der Liste, damit nicht geprüft** (Bestand, ohne Fall; `SPEC-052` Punkte 11 und 18): Dateimodi, ignorierte Dateien, Zeitstempel. Ein anderes Auschecken von Zeilenenden (`core.autocrlf`) gibt eine andere Liste und damit die Ablehnung, nie eine falsche Freigabe.

**Gruppe Z — Ziel**

11. **Erfolg ersetzt.** Nach Erfolg sind beide Dateien byte-gleich denen der Quelle; ein älterer Nachweis (gültig, beschädigt, fehlend) und ein fehlendes `.harness/state/` ändern das Ergebnis nicht (das Verzeichnis wird angelegt). Vier Fälle. Entschieden: offen.
12. **Ablehnung ändert nichts.** Alte Dateien bleiben byte-gleich, `.harness/state/` trägt keine Hilfsdatei (auch nicht nach Ablehnung in Q, G und Z); zwei Fälle (Dateien, Hilfsdateien), je Gruppe geprüft. Entschieden: offen.
13. **Wiederholung und Alter.** Ein zweites Übernehmen derselben Quelle endet mit Exit 0 und gleichen Bytes. Ein Nachweis kennt keinen Zeitpunkt; ein „älterer“ Nachweis desselben Baums ist so gültig wie ein neuerer, ein Nachweis eines anderen Baums ist ungleich. Ein Fall. Entschieden: offen.
14. **Schreiben scheitert oder wird abgebrochen.** `.harness/state/` nicht schreibbar (unter root nicht herstellbar, übersprungen mit Meldung wie in `stop-hook-gegenprobe`): Exit 1, alter Nachweis unverändert. `SIGINT` und `SIGTERM` zwischen Prüfung und Ersetzen: keine Hilfsdatei, die Liste ist alt oder neu, nie halb (Ersetzen mit `mv`). Drei Fälle. Entschieden: offen.
15. **Quelle ist das Ziel.** Gleicher Pfad, mit `./` und mit Schrägstrich am Ende, und ein Symlink auf das Verzeichnis des Ziels: Exit 2, nichts geändert (Vergleich der Dateien mit `-ef`). Drei Fälle. Entschieden: offen.
16. **Quelle im Arbeitsbaum des Ziels.** Liegt die Quelle als untracked Verzeichnis im Baum, sind ihre Dateien Teil der Liste und der Vergleich schlägt fehl (Fall aus 6, eigene Zeile in der Meldung nicht zugesagt); liegt sie unter `.harness/state/` (ignoriert), wird sie angenommen. Zwei Fälle. Das Rezept (Liefer-Punkt 3) nennt den zweiten Ort. Entschieden: offen.

**Gruppe A — Aufrufarten (`AGENTS.md` §3.12, Aufrufarten des Treibers)**

17. **Make-Flags.** `-n`, `-q`, `-t`: Das Rezept trägt weder `$(MAKE)` noch `+`, `make` führt es nicht aus, nichts ändert sich; `-i`: Die Einzelschritte erben `-i`, ein abgelehntes Übernehmen sähe aus wie Erfolg, deshalb lehnt das Skript `-i` in `MAKEFLAGS` ab (Exit 2, nichts geändert, wie `record-gates.sh` den Nachweis bei `-i` nicht schreibt); `-n`, `-q` oder `-t` zusammen mit `-i`: `-n`, `-q`, `-t` gelten; `-k`, `-j`: ohne Wirkung; maßgeblich ist `MAKEFLAGS` auch aus der Umgebung. Ein Fall je Flag und je Kombination (`-n`, `-q`, `-t`, `-i`, `-ni`, `-qi`, `-ti`, `-i` aus der Umgebung, `-k`, `-j`). Entschieden: offen.
18. **Aufruf aus einem Unterverzeichnis.** Das Werkzeug urteilt über den Baum der Wurzel (`git rev-parse --show-toplevel`); eine relative Quelle gilt gegen das Verzeichnis des Aufrufs. Zwei Fälle. Entschieden: offen.
19. **Gleichzeitige Läufe** (zwei Übernahmen, eine Übernahme zusammen mit `make gates`, mit einer Änderung des Baums oder mit dem Hook): nicht zugesagt (`SPEC-052` Punkt 8, Aufruf-Voraussetzung nacheinander), ohne Fall. Ändert sich der Baum nach der Übernahme, vergleicht der Hook neu und sperrt oder gibt nach seinen eigenen Punkten frei; das Werkzeug kann nie öfter freigeben als der Hook, weil dieser den Baum liest, nicht das Werkzeug. **Aufruf-Voraussetzung:** Das Werkzeug sagt keine Ordnung über mehrere Aufrufe zu.

**Gruppe H — Kopplung an den Hook**

20. **Der unveränderte Hook liest das Ergebnis.** In einer Kopie des Temp-Repos mit dem Hook aus dem Arbeitsbaum antwortet er nach Erfolg bei diesem Baum `approve` ohne Stand der schnellen Kette (Punkt 16); nach einer Ablehnung bei fehlendem Nachweis `block` (Punkt 19), bei vorhandenem Nachweis nach dessen eigenem Befund. Drei Fälle. Die Mutation, das Werkzeug bilde die Liste selbst (`git ls-files` ohne `--others`), färbt den Fall mit der untracked Datei (Randform 6). Entschieden: offen.

**Gruppe S — Grenze, offen genannt (ohne Fall)**

21. **Jeder lokale Nachweis ist von Hand zu schreiben** (`dateiliste_erzeugen > .harness/state/gates-passed.dateiliste`, `make record-gates` einzeln, `enforce.mk`); das gilt heute schon. Das Werkzeug ist der einzige sanktionierte Weg und prüft Gleichheit, Unversehrtheit und Zusammengehörigkeit; es kann nicht prüfen, dass die Quelle `make gates` fuhr, mit denselben Gates, im selben Image-Stand und mit allen Gegenproben, die dann gelten. Akzeptiertes Negativ; Grund: Der Nachweis trägt keine Herkunft, und eine Herkunft ohne Signatur wäre von Hand ebenso zu schreiben. Die Grenze steht im Vertrag und im Rezept (§3.11).
22. **Der Nachweis gilt einem Baum, nicht einem Zweig.** Ein Merge, der den Baum des Hauptzweigs ändert (Merge-Commit mit anderem Inhalt), macht die Listen ungleich; die Regel an den Aufrufer lautet, den Zweig vor dem Merge auf den Stand des Hauptzweigs zu bringen und dort `make gates` zu fahren, oder `make gates` im Hauptbaum. Das Werkzeug lehnt ab, es rät nicht.

### Risiken

- **Der Nachweis gilt Umgebungen, die niemand vergleicht** (Randform 21; ein Gate, das nur in einer Umgebung rot wird, bleibt in der Quelle unbemerkt). — **Ausgang (vorgesehen):** weiter offen → neues `BEO-REPO/nachweis-uebernommen-aus-fremder-umgebung` (Beleg bei der Closure).
- **Der Hook wird zum dritten Mal berührt, ohne dass die Gegenprobe seine Umgebung prüft** (`BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung`, 2×; Gruppe H fährt den Hook und liest seine Antwort, ob die Umgebung von Claude Code sie so liest, prüft sie nicht; der Kopf der Gegenprobe nennt es). — **Ausgang (vorgesehen):** weiter offen → dieselbe Beobachtung; legt der Planner bei der Closure einen dritten Beleg an, braucht sie nach §8 einen eigenen Folge-Slice (Frage 6).
- **Der Hauptbaum ist selten gleich der Quelle** (unversionierte inerte Änderung, untracked Datei; dann lehnt das Werkzeug ab, und die Sperre bleibt). Nach einer Übernahme verlangt die nächste inerte Änderung im Hauptbaum wieder den Stand der schnellen Kette, also Docker (Frage 3). — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (Vermerk mit der Zahl der Ablehnungen und Gründe aus der Probe), bei Entscheidung des Nutzers zu Frage 3 eingetreten → ein Folge-Slice mit Kennung.
- **Die Gegenprobe hängt an `GATE_CHECKS`, weil die Listen von `slice-harness-gegenproben-aus-gates` fehlen.** — **Ausgang (vorgesehen):** entfallen, wenn beim Start `slice-harness-gegenproben-aus-gates` in `done/` liegt und `grep -l nachweis-uebernehmen-gegenprobe harness/mk/*.mk` das Fragment in der Liste zeigt (nicht an `GATE_CHECKS`); sonst eingetreten → Rückführung nach `open` (§4).
- **Das Rezept für den Weg der Dateien aus dem Container (`ssh`, `tar`) ist keine geprüfte Zusage** (`AGENTS.md` §3.11). — **Ausgang (vorgesehen):** entfallen, wenn der Abschnitt das Rezept als nicht geprüft kennzeichnet und die Gegenprobe den Fall mit einer lokal durch eine `tar`-Pipe angelegten Kopie führt (Review prüft den Wortlaut).
- **Die Größe liegt unter der Marke, die Spanne ist breit** (495 bis 785, Marke 900); sie hängt an der Zahl der Fälle der Gegenprobe. — **Ausgang (vorgesehen):** entfallen, wenn die Messung nach der Metrik unter 900 bleibt; darüber meldet der Implementer (§4 Rückführung *zu groß*).

**Fragen an den Nutzer** (jede mit Empfehlung des Planners; eine Empfehlung ist keine Entscheidung):

1. **Stelle in der Welle.** Eingetragen ist die fünfte Stelle (seit dem Schnitt von `slice-harness-gegenproben-aus-gates-fahrstelle` am 2026-10-18; vorher die vierte), hinter `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`, ohne Lösung aus der Reihenfolge der Wellen. Empfehlung: Lösung aus der Reihenfolge wie bei den drei Slices davor (Start-Trigger `slice-harness-gegenproben-aus-gates` in `done/`), und die dritte Stelle, vor `-fahrstelle`, wenn die Sperre des Hauptbaums dringender ist als die Pflicht der Fahrstelle; das Vorziehen ändert die Sätze „direkt hinter“ in `-sammelziele`, `-fahrstelle`, `-fahrstelle-betrieb`, in `welle-gegenproben-und-mutanten` §5 und im Drift-Log, deshalb hat der Planner es nicht ohne Antwort getan.
2. **Quelle nur als Verzeichnis, Transport als Rezept** (`ssh` mit `tar` oder `docker cp`, in der Wartungs-Doku) statt eines Transports im Werkzeug. Empfehlung: ja (§1, zweiter Ausschluss).
3. **Zweite Sperre des Hauptbaums.** Das Werkzeug löst die Sperre nach einem Pull mit Code, wenn der Hauptbaum danach unverändert ist. Die nächste inerte Änderung dort (zum Beispiel eine Plan-Datei) verlangt den Stand der schnellen Kette, also `make gates-schnell` mit Docker. Optionen: (a) so lassen; (b) das Werkzeug übernimmt bei Gleichheit auch `gates-schnell-passed.diffsha` (rund 8 Fälle und 80 bis 120 Zeilen mehr, der Hash des ganzen Baums muss dann gleich sein, was bei einer inerten Änderung im Hauptbaum nie der Fall ist); (c) der Hook gibt inerte Änderungen ohne den Stand der schnellen Kette frei (ADR, Lockerung nach `AGENTS.md` §3.6, eigener Slice). Empfehlung: (a) jetzt, die Ablehnungen der Probe messen; (c) ist die Antwort auf „Keine blockierenden Aktionen“ für Planung im Hauptbaum und braucht eine eigene Entscheidung.
4. **Pflicht in den Rollen-Anweisungen** (wann das Werkzeug zu rufen ist, etwa in `implement-slice` nach dem Merge). Empfehlung: erst nach der Lieferung und einer Probe, dann als Satz in der Fahrstelle oder einem Folge-Slice.
5. **ADR.** Empfehlung: ja, vor dem Code. Ein zweiter Schreiber ändert die Lesart „der Hook glaubt den Nachweis, den `make gates` in dieser Umgebung schrieb“ ([ADR-0038](../../adr/0038-gate-kette-voll-und-schnell.md), `SPEC-052` Punkt 21); die Annahme ist eine Entscheidung des Nutzers.
6. **Folge-Slice für die Umgebung des Hooks** (`BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung`, dritter Beleg). Empfehlung: Der Planner entscheidet bei der Closure nach dem Befund des Reviews und legt ihn an, wenn die Grenze erneut als offen stehen bleibt.

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
- **Folge-Slices:** — (Frage 3 und Frage 6 können je einen anlegen)
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (1×; der Slice ist die Gegenmaßnahme für den Nachweis, kein Auftreten), `BEO-REPO/stop-hook-gegenprobe-prueft-nicht-die-umgebung` (**2×**; Gruppe H fährt den Hook in einer Kopie und liest seine Antwort, die Umgebung von Claude Code prüft sie nicht; erhält die Beobachtung bei der Closure einen dritten Beleg, ist sie nach Modul 5 keine Notiz mehr, sondern eine Lücke und braucht einen eigenen Folge-Slice, Frage 6), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (DoD 2 trägt je Element einen Fall und eine Mutation), `BEO-REPO/spec-randform-erst-im-review-entschieden` (§6 nennt die Randformen vor dem Code), `BEO-REPO/slice-waechst-durch-uebernahmen` (der Slice übernimmt nichts und trägt keine Sendung ein), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` und `BEO-REPO/groessenprobe-ohne-messbaren-stand` (§4 rechnet in der Metrik mit dem Vergleichswert der Gegenprobe des Hooks), `BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet` (jede Wahl der Stelle, des Namens und der Option steht als Empfehlung des Planners; entschieden ist nur der Auftrag des Nutzers). Keiner der übrigen Treffer erreicht mit diesem Slice 3× neu.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Liefer-Punkte 3 (Werkzeug | Gegenprobe | Gate-Index und Rezept); der Slice liegt an der Grenze, ein vierter Punkt hat keinen Platz (Frage 3, Option (b), wäre einer und gehört in einen Folge-Slice). Schichten 2: Harness-Werkzeug mit Gegenprobe und Gate-Index (eine Schicht) und Wartungs-Doku `docs/developer/` (eine Schicht); `AGENTS.md` bleibt unberührt, Planung, Vertrag und ADR sind keine Schicht. Schätzung in der Metrik von §4: 495 bis 785 (Marke 900, Meldeschwelle 700). **Sendungen:** Der Plan trägt keine Sendung in einen anderen Slice ein (`grep` nach `slice-harness-nachweis-uebernehmen` unter `docs/plan/planning/` findet nur diese Datei, `welle-gegenproben-und-mutanten` und das Drift-Log der Roadmap, keinen Slice); die Ausschlüsse in §1 sind Klasse 2, 3 und 4 (Bestand, anderer Vorgang, Schicht), keiner nennt einen Nehmer. Voraussetzungen, keine Sendungen, an die der Slice sich bindet: `slice-harness-gegenproben-aus-gates` (Listen und Namensregel für die Gegenprobe, Wortlaut von `SPEC-052` Punkt 9), der Geber liegt in `in-progress/` und der Start-Trigger in §4 verlangt `done/`.
