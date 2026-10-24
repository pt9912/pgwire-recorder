# Slice slice-test-unter-last-ursache: Die Ursache des Rots unter CPU-Last ist benannt (Frist, Reihenfolge oder Wettlauf im Code), bevor eine Regel oder Frist geändert wird, und bis dahin gilt eine Zwischenregel für den Implementer

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. **Empfehlung des Planners vom 2026-10-22** (Auftrag des Koordinators), **nicht Entscheidung des Nutzers**: Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist (`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht). Geprüft und nicht gewählt: [welle-gegenproben-und-mutanten](../welle-gegenproben-und-mutanten.md) trägt das Thema *Güte der Tests*, ihr Start-Trigger verlangt aber `welle-harness-regeln` in `done/`, und die Zwischenregel (Liefer-Punkt 3) und die Messung sollen nicht bis dahin warten; [welle-harness-regeln](../welle-harness-regeln.md) trägt Regeln, hier ist die Regel nur einer von drei Liefer-Punkten. Entschieden hat der Nutzer, dass es den Slice gibt: „deine empfehlungen“ vom 2026-10-22, so vom Koordinator übermittelt, auf die Vorschläge des Planners in `BEO-REPO/test-faellt-unter-cpu-last-aus`; über Welle, Stelle und Wortlaut der Zwischenregel ist nichts entschieden (`BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet`). Angelegt in `open/`, weil er nicht priorisiert ist; den Anspruch trägt der Nutzer.

**Bezug:** Kein Anspruch des Produkts — der Slice untersucht Tests und ändert eine Anweisung des Harness. Bindung an das Register: `BEO-REPO/test-faellt-unter-cpu-last-aus` (6×, Ausgang `geplant` mit diesem Slice seit 2026-10-22); abgegrenzt: `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` (dort fehlt die Frist; `geplant` mit `slice-harness-lint-warten-ohne-frist`), `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test`, `BEO-REPO/lastabhaengige-zaehlung` und `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (Gegenproben, die eine Größe der Umgebung trifft). Regelstelle: `.claude/commands/implement-slice.md`. Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912 (Zweig impl-last1)

**Autor:** Planner. **Datum:** 2026-10-22; nachgezogen 2026-10-24 (siebtes Auftreten, Marke nach `AGENTS.md` §3.14, Anspruch).

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Ursache des Rots unter CPU-Last ist benannt, bevor eine Regel oder eine Frist geändert wird: (1) Die zwei Go-Tests, die in den Belegen des Registers unter Last rot waren — `TestE2EPlayVergleichAbweichung` mit dem Unterfall `Serverfehler_ohne_Vorbild` (`test/integration/play_vergleich_e2e_test.go`, vier Belege) und `TestRunPlayZweitesSignal` (`internal/bootstrap/play_test.go`, ein Beleg) — und `TestE2EPlayVergleichAlle` mit dem Unterfall `Umgebungsvariable` (`test/integration/play_vergleich_e2e_test.go`, siebtes Auftreten 2026-10-24: Fehlertext `DROP DATABASE`, Frist 30 s, im ersten vollen `make gates` eines Doku-Slice, im zweiten Lauf am selben Stand grün) laufen unter einer künstlichen Last, die als Aufruf wörtlich in §7 steht, und die Ursache jedes Rots ist benannt als *Frist* (die 30 s sind für die Arbeit des Tests zu knapp), *Reihenfolge* (der Test wartet auf einen Zustand, den er selbst nicht abwartet) oder *Wettlauf im Code* (das Produkt hat einen Fehler, den Last sichtbar macht), mit der Messung, die sie trennt; (2) als zweite Probe der Fall `p17-budget-sigint` der Gegenprobe `container-einrichten-gegenprobe` (die Mutation `m-p17-signal-int`: ein Rot in 21 Läufen bei Last über 9, `slice-harness-devcontainer-einrichten`) unter derselben Last, mit Rotquote und Ursache oder dem Befund „nicht reproduziert“; (3) `.claude/commands/implement-slice.md` trägt bis dahin die Zwischenregel: Bei Rot der zwei Unterfälle `TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild` und `TestE2EPlayVergleichAlle/Umgebungsvariable` fährt der Implementer den Test allein und wiederholt `make gates` einmal, bevor er einen Befund schreibt; Rot von `TestRunPlayZweitesSignal` ist ein Befund (Wettlauf im Code, reproduzierbar), keine Wiederholung, bis ein Folge-Slice es behebt. Das Register spricht von „drei Tests“ (gemeint: drei Auftreten bis zur Schwelle); die Belege nennen drei Go-Tests (zwei davon in derselben Datei, `play_vergleich_e2e_test.go`), dazu den Fall der Gegenprobe. Der Plan zählt sie als *einen* Liefer-Punkt (Ursache der Go-Tests), nicht als drei.

**Beleg-Stand** (Stand 2026-10-24, gemergter Stand des Registers; die Belege stehen dort, nicht hier): Sieben Auftreten (das siebte: `TestE2EPlayVergleichAlle/Umgebungsvariable`, `slice-homebrew-nachweis-macos-arm64`; vom Planner als Sendung in diesen Slice eingetragen, §8); bis zum sechsten sechs Auftreten in `BEO-REPO/test-faellt-unter-cpu-last-aus/evidence/`; die Unterfälle liefen allein in 1,19 s und 19,85 s und rissen unter Last die Frist von 30 s (Last 6,0 bis 22 auf 20 Kernen). `TestRunPlayZweitesSignal` fiel mit 40 gleichzeitigen Rechenschleifen auf 20 Kernen in 3 von 15 und 2 von 12 Läufen aus; die Tests des Slice `slice-v1-abschluss-zeitangaben-uhr` liefen bei 60 Schleifen 75 Mal ohne Fehlschlag. Ein Gegenbeleg steht im Register (`slice-v1-abschluss-release-v0-1-1`: drei volle Läufe von `make gates` grün bei Last 13,7); er widerlegt die Beobachtung nicht, er sagt, dass Last allein sie nicht auslöst.

**Vermutungen des Planners, nicht gemessen, nicht Teil der Zusage** (sie sagen der Messung, wo sie suchen soll): Der Fehlertext der Belege von `TestE2EPlayVergleichAbweichung` und (siebtes Auftreten) `TestE2EPlayVergleichAlle/Umgebungsvariable` nennt `DROP DATABASE IF EXISTS vgl_abw` (`leereDatenbank`), nicht den Lauf von `play`; PostgreSQL wartet bei `DROP DATABASE`, solange eine Verbindung zur Datenbank besteht, und die Unterfälle des Tests nutzen eine Datenbank nacheinander. Beim Signal-Test wartet `typenBis` auf die Nachricht `X` nach dem zweiten Signal, und beide Signale sind Aufrufe im selben Test. Ob das die Ursache ist, entscheidet allein die Messung.

**Sendung von `slice-harness-randform-umgebung-der-gegenprobe`** (eingetragen vom Planner im Commit des Anlegens, `AGENTS.md` §3.13): Dessen Regel zur Fahrt unter Last braucht eine Last, die jemand wiederholbar beschrieben hat. Dieser Slice nennt in §7 den Aufruf der künstlichen Last wörtlich (Erzeuger, Anzahl, Dauer, Messung der Last vor und nach dem Lauf; DoD 1); `grep` nach der Kennung des Gebers findet die Sendung in §1 und DoD.

**Wortlaut und Entscheidung des Nutzers.** Die Zwischenregel steht in `.claude/commands/implement-slice.md` Schritt 16. Regeländerungen in `.claude/commands/` sind eine Entscheidung des Nutzers: Der Nutzer hat am 2026-10-22 („deine empfehlungen“, so vom Koordinator übermittelt) die Vorschläge des Planners im Register gewählt, den Slice zur Ursache und die Zwischenregel, und am 2026-10-24 den Wortlaut bestätigt („Vorschlag 2 (Recommended)“, Fundstelle und Wortlaut in §7): die Regel gilt nur für die zwei Unterfälle von A und B, Rot von `TestRunPlayZweitesSignal` ist ein Befund.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine **Änderung an einer Frist, einem Test oder dem Produktcode** — die Ursache wird benannt, die Behebung nicht: Eine Änderung vor der Ursache wäre Raten und ließe die Beobachtung ohne Beleg stehen. Proben ändern Fristen und Tests nur in einer Kopie außerhalb des Repos. Liegt die Ursache in Test oder Code, ist die Behebung ein eigener Slice; eine Adresse gibt es nicht, der Planner legt ihn nach der Ursache an (Ausgang der Risiken in §6).
- Das **Warten ohne Frist** (`slice-harness-lint-warten-ohne-frist`) — ein anderer Vorgang: Dort fehlt die Frist und der Lint soll es melden; hier ist sie da und reißt (`BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`).
- Ein **Werkzeug für Last** (ein `make`-Ziel, ein Skript im Repo) — anderer Vorgang (Werkzeug statt Untersuchung): Der Aufruf steht wörtlich in §7; braucht das Repo ihn wiederholt, ist das ein Slice mit Vertrag und Gegenprobe.
- Die **Regel zur Fahrt einer Gegenprobe unter Last** — `slice-harness-randform-umgebung-der-gegenprobe` schließt sie aus (dort §1) und nennt diesen Slice als Quelle der Last-Beschreibung (Sendung oben); die Regel selbst ist kein Gegenstand hier.
- Ein **Durchmustern aller Tests mit Frist** des Bestands — Bestand bleibt bewusst stehen: Untersucht wird, was das Register belegt; ein Durchmustern wäre ein Sensor (Lint) mit Vertrag und Gegenprobe.
- Eine **Änderung der Gegenprobe** `container-einrichten-gegenprobe` — sie gehört dem Slice, der sie lieferte (`slice-harness-devcontainer-einrichten`, `done/`); hier wird sie nur gefahren. Eine Behebung ihres Signalfalls ist ein Folge-Slice nach der Ursache.
- `AGENTS.md`, `.claude/agents/`, `.harness/skills/` und die übrigen Commands — Schicht-Abgrenzung und Größe: Die Zwischenregel steht in einer Datei; eine zweite Fundstelle (Verifier, Reviewer) brächte eine dritte Schicht.
- Produkt-Code und Benutzerhandbuch — Schicht-Abgrenzung: **kein Produktcode**, kein Verhalten des gebauten Binary ändert sich.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Liefer-Punkt 1 (Ursache der Go-Tests, drei Tests der Belege):** §7 trägt unter *Belege des Implementers* (a) den Aufruf der künstlichen Last wörtlich und wiederholbar (Erzeuger, Anzahl, Dauer, Messung der Last vor und nach dem Lauf; Sendung von `slice-harness-randform-umgebung-der-gegenprobe`), (b) je Test (`TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild`, `TestE2EPlayVergleichAlle/Umgebungsvariable`, `TestRunPlayZweitesSignal`) den Lauf allein (Dauer) und die Läufe unter Last (Anzahl, Zahl der roten Läufe), (c) die Proben in einer Kopie außerhalb des Repos, die Frist, Reihenfolge und Wettlauf trennen (der Implementer wählt sie und nennt je Probe, was sie trennt), mit ihrem Ergebnis, und (d) die Ursache je Test als *Frist*, *Reihenfolge* oder *Wettlauf im Code* oder, wenn das Rot nicht entsteht, den Befund „nicht reproduziert in N Läufen bei Last L“ (Randform 3). Die Tests laufen über das Docker-Gerüst des Repos (`AGENTS.md` §3.1); der Aufruf steht wörtlich in §7.
- [x] **Liefer-Punkt 2 (zweite Probe: Signalfall der Gegenprobe):** §7 trägt dasselbe für den Fall `p17-budget-sigint` der Gegenprobe `container-einrichten-gegenprobe`: der enge Lauf der Mutation `m-p17-signal-int` mit `GP_MUT` (nicht die volle Gegenprobe, [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)) unter derselben Last, mit der Zahl der Läufe, der Zahl der roten und der Ursache oder dem Befund „nicht reproduziert“; ein Rot durch das Zeitbudget der Gegenprobe selbst ist von einem Rot des Falls getrennt (Randform 6). Die Gegenprobe wird nicht geändert.
- [ ] **Liefer-Punkt 3 (Zwischenregel):** `.claude/commands/implement-slice.md` trägt in Schritt 16 die Zwischenregel: Bei Rot von `TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild` oder `TestE2EPlayVergleichAlle/Umgebungsvariable` fährt der Implementer den Test allein und wiederholt `make gates` einmal, bevor er einen Befund schreibt; §7 des Plans nennt beide Läufe und die Last; das erste Rot bleibt ein Auftreten im Register; Rot von `TestRunPlayZweitesSignal` ist ein Befund, keine Wiederholung. Der Wortlaut ist vom Nutzer am 2026-10-24 bestätigt, vor dem Commit; Fundstelle in §7.
- [x] Kein Code, keine Frist und kein Test im Repo ist geändert (Beleg: der Befehl aus §3.14 über `.claude internal test tools harness` liefert 1, die Regelzeile in `implement-slice.md`; unter `internal test tools harness` ist `numstat` leer; in §7 *Größe* genannt).
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `.claude/` ändert).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel: `implement-slice.md`, `liegt in`, soweit aufgenommen; sonst benannte Spec-Lücke: die Ursache ist nicht reproduzierbar, mit den Läufen).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: `BEO-REPO/test-faellt-unter-cpu-last-aus` von `geplant` auf `verkörpert` (Zwischenregel: Zielort und Anker) mit der benannten Ursache oder dem Befund „nicht reproduziert“ im `state.md`; liegt die Ursache in Test oder Code, steht der Folge-Slice mit Kennung als Datei in `open/` (Planner, Closure); kein Zähler wird gesetzt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/play_vergleich_e2e_test.go` (beide Tests), `internal/bootstrap/play_test.go` | gelesen und gefahren; Proben nur in einer Kopie außerhalb des Repos | Ursache des Rots unter Last (Frist, Reihenfolge, Wettlauf) |
| `tools/dev/harness/container-einrichten-gegenprobe.sh` | gefahren (`GP_MUT`), nicht geändert | zweite Probe: Signalfall |
| `.claude/commands/implement-slice.md` | update | Zwischenregel in Schritt 16 |
| `docs/plan/planning/observations/BEO-REPO/test-faellt-unter-cpu-last-aus/state.md` | update | Ausgang bei Closure |

Tests und Gates: keine neuen. Der Slice untersucht vorhandene Tests; die Zwischenregel ist Prosa (`AGENTS.md` §3.10 verlangt Mutationen für Verträge von Code und Gates).

**Wortlaut** (vom Nutzer am 2026-10-24 bestätigt, Vorschlag 2; Träger: `implement-slice.md` Schritt 16):

> *`implement-slice.md`, Schritt 16, Ergänzung (so bestätigt am 2026-10-24):* Endet `make gates` rot in `TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild` oder `TestE2EPlayVergleichAlle/Umgebungsvariable` (Frist von 30 s unter CPU-Last), fährst du den Test allein und wiederholst `make gates` einmal, bevor du einen Befund schreibst; §7 nennt beide Läufe und die Last des Hosts. Das erste Rot bleibt ein Auftreten (`BEO-REPO/test-faellt-unter-cpu-last-aus`), ein zweites Rot ist ein Befund. Rot von `TestRunPlayZweitesSignal` ist ein Befund (Wettlauf im Code, reproduzierbar), keine Wiederholung, bis ein Folge-Slice es behebt.

**Reihenfolge:** Empfehlung des Planners: Punkt 3 zuerst (der Nutzer hat den Wortlaut am 2026-10-24 bestätigt) (die Regel wirkt ab dem Commit auf jeden Lauf von `make gates`, auch während die Messung läuft), dann Punkt 1 und 2 auf einem freien Host; keiner der drei Punkte wartet auf einen anderen. Der Plan des Slice folgt jedem Commit (`AGENTS.md` §3.9).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice lag in `open/`, bis er am 2026-10-24 auf Auftrag des Koordinators priorisiert wurde (`open` → `next`, `Verantwortlich: pt9912`); der Anspruch ist der reine Move nach `in-progress/` auf `main`, `Verantwortlich:` trägt dann den Zweig (`pt9912 (Zweig impl-last1)`). Kein anderer Slice ist Start-Trigger. **Voraussetzung der Messung (Punkte 1 und 2):** Der Host ist frei — kein anderer Rollenlauf, keine Gegenprobe und kein `make gates` zur selben Zeit auf demselben Host und Docker-Daemon; die künstliche Last verfälscht fremde Läufe, und fremde Läufe verfälschen die Messung (`BEO-REPO/lastabhaengige-zaehlung`). Der Nutzer gibt den Host frei; §7 nennt, dass kein anderer Lauf lief. WIP-Limit 1. Der Implementer legt den Wortlaut der Zwischenregel dem Koordinator vor (dieser dem Nutzer) und committet `implement-slice.md` erst nach der Bestätigung; bis dahin ruht Punkt 3, nicht der Slice.

**Größe.** Metrik und Befehl stehen in `AGENTS.md` §3.14 (einmal, dort); gilt der Befehl von dort. **Pfade des Plans:** `.claude internal test tools harness` (mit den zwei Ausschlüssen des Befehls). **Marke 1500** (`AGENTS.md` §3.14, Entscheidung des Nutzers vom 2026-10-22); der Plan nennt keine Meldeschwelle und keine Schätzung (§3.14, Nachzug beim Anspruch 2026-10-24; die frühere Marke 100 mit Meldeschwelle 80 entfällt). **Basis:** der Commit, der diese Datei nach `in-progress/` bewegt; sein Hash steht hier: `ffdacf7d41a3d551a8d164470e070a8a0d2791e0` (Move nach `in-progress/`, 2026-10-24). Die Zeilenmarke misst diesen Slice schlecht — sein Umfang ist die Messung in §7 —, deshalb trägt die zweite Rückführung eine Zeitbedingung.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Messung der Punkte 1 und 2 passt nicht in einen Rollenlauf von einer Stunde (Entscheidung des Nutzers vom 2026-10-21 für jeden Rollenlauf): Der Implementer hält an, trägt den Stand in §7 ein, und der Planner schneidet nach (Punkt 2 als eigener Slice); oder die Summe nach der Metrik erreicht die Marke 1500.
- `in-progress` → `open` (blockiert): Der Host lässt sich für die Messung nicht freihalten, oder die Last lässt sich im Docker-Gerüst nicht erzeugen, ohne einen Teil des Repos zu ändern (Befund an den Planner, Randform-Rückgabe). Liegt die Ursache im Produktcode oder verlangt die Behebung eine Änderung an Test oder Code, hält der Implementer an und gibt es dem Planner zurück; behoben wird hier nichts (§1).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; in §7 stehen die Last wörtlich, die Läufe und die Ursache je Test (Frist, Reihenfolge oder Wettlauf im Code) oder der Befund „nicht reproduziert in N Läufen bei Last L“, dasselbe für den Signalfall der Gegenprobe, die Bestätigung des Wortlauts der Zwischenregel mit Datum; Review-Report liegt vor; `make gates` grün; Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12; der Slice liefert keinen Vertrag von Code, aber eine Messung und eine Regel, deren Ränder der Planner sonst im Einzelfall entscheidet). Entscheidungsort ist, wo es steht: der Plan für die Messung, der bestätigte Wortlaut für die Regel:

1. *Welche Tests:* die drei Go-Tests der Belege (`TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild`, `TestE2EPlayVergleichAlle/Umgebungsvariable`, `TestRunPlayZweitesSignal`) und ein Fall der Gegenprobe (§1); der Bestand mit Fristen ist ausgenommen.
2. *Die Last:* gleichzeitige Rechenschleifen im Docker-Gerüst, mindestens doppelt so viele wie Kerne (Beleg: 40 Schleifen auf 20 Kernen machten `TestRunPlayZweitesSignal` rot, 60 Schleifen ließen die Tests von `slice-v1-abschluss-zeitangaben-uhr` grün); der Aufruf steht wörtlich in §7, die Last ist vor und nach dem Lauf gemessen (Lastzahl und Kernzahl). Entschieden: der Plan; die Zahl der Schleifen wählt der Implementer und nennt sie.
3. *Zahl der Läufe und Aussagekraft:* je Test mindestens 15 Läufe unter Last (die Belege: 15 und 12); beim Signalfall der Gegenprobe sind es 8 (Laufzeit 27 s je Lauf, Abweichung genannt in §7). Ein Ergebnis ohne Rot nennt die Läufe und die Last; es widerlegt die Beobachtung nicht (der Gegenbeleg des Registers sagt dasselbe), und die Zwischenregel bleibt der Ausgang des Eintrags.
4. *Proben in einer Kopie:* Fristen, Tests und Reihenfolge werden nur in einer Kopie außerhalb des Repos geändert; nichts davon wird committet (§1).
5. *Ort des Wartepunkts:* Der Beleg nennt `DROP DATABASE` (`leereDatenbank`), nicht `play`; ob ein Unterfall der Datenbank `vgl_abw` den nächsten blockiert, ist eine Vermutung (§1), die eine Probe trennen muss (etwa ein Unterfall mit eigener Datenbank).
6. *Zeitbudget der Gegenprobe:* Die Gegenprobe `container-einrichten-gegenprobe` trägt ein Zeitbudget ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)); unter Last ist es ein zweiter Grund für Rot. Die Messung trennt ihn vom Rot des Falls `p17-budget-sigint` (Zeile mit Zeit und Last in der Ausgabe) und nennt beide.
7. *Zwischenregel, Geltung und Stelle:* Schritt 16 des Implementers; gilt für ihn, nicht für Reviewer und Verifier (§1). Der Vorschlag im Register sagt „bei Rot dieses Unterfalls“ (`Serverfehler_ohne_Vorbild`); der bestätigte Wortlaut nennt die zwei Unterfälle von A und B (`Serverfehler_ohne_Vorbild`, `Umgebungsvariable`, siebtes Auftreten), weil ihr Rot nicht erklärt ist; `TestRunPlayZweitesSignal` ist ausgenommen, weil die Messung dort einen reproduzierbaren Wettlauf im Code zeigt, den „einmal wiederholen“ verdeckte — dort ist Rot ein Befund (Entscheidung des Nutzers vom 2026-10-24).
8. *Wiederholung:* genau eine; ein zweites Rot ist ein Befund, das erste bleibt ein Auftreten im Register (jedes Rot ist eine Datei in `evidence/`, der Zähler folgt ihnen).
9. *Aufruf-Voraussetzung, Aufrufarten des Treibers, Plattform, Stelle einer Nachricht:* nicht berührt (`AGENTS.md` §3.12); die Messung ist eine Linux-Messung im Container und sagt nichts über macOS oder Windows (§7 nennt es nicht als Zusage).

**Risiken:**

- Das Rot lässt sich unter künstlicher Last nicht herstellen (der Gegenbeleg des Registers: drei volle Läufe von `make gates` grün bei Last 13,7) — **Ausgang:** weiter offen → `BEO-REPO/test-faellt-unter-cpu-last-aus` (Beleg: Läufe und Last in §7; die Zwischenregel bleibt der Ausgang des Eintrags).
- Die Ursache liegt im Produktcode (ein Wettlauf) — **Ausgang:** eingetreten → Folge-Slice mit Kennung, vom Planner angelegt; der Befund geht an den Planner (§4), behoben wird hier nichts.
- Die Ursache liegt in Frist oder Reihenfolge des Tests — **Ausgang:** eingetreten → Folge-Slice mit Kennung, vom Planner angelegt (Änderung am Test, mit der Marke aus der Messung).
- Die Zwischenregel verdeckt eine echte Wettlaufbedingung, weil „einmal wiederholen“ ein Rot grün färbt — **Ausgang:** weiter offen → `BEO-REPO/test-faellt-unter-cpu-last-aus` (jedes Rot bleibt ein Auftreten; ein zweites ist ein Befund, Randform 8; für `TestRunPlayZweitesSignal` gilt die Wiederholung nicht, schon das erste Rot ist ein Befund).
- Die Last der Messung stört einen anderen Lauf auf demselben Host — **Ausgang:** entfallen, wenn §7 nennt, dass kein anderer Lauf auf dem Host war (Voraussetzung in §4).
- Die Messung passt nicht in einen Rollenlauf von einer Stunde — **Ausgang:** eingetreten → Rückführung nach `next` (§4), Punkt 2 als eigener Slice (vom Planner).
- Die Zuordnung von `m-p17-signal-int` zu dieser Beobachtung ist ein Urteil des Planners (Beleg `evidence/slice-harness-devcontainer-einrichten.md`), kein Befund — **Ausgang:** offen bis Closure (entfallen, wenn die Ursache dieselbe Klasse trägt; sonst wandert die Beleg-Datei in einen eigenen Eintrag).
- Der Slice erreicht die Marke von §4 (1500, Metrik `AGENTS.md` §3.14) — **Ausgang:** offen bis Closure; Marke und Basis stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Hier stehen: *Belege des Implementers* (die Last wörtlich, je Test die Läufe allein und unter Last, die Proben und ihr Ergebnis, die Ursache oder der Befund „nicht reproduziert“; dasselbe für den Signalfall; dass kein anderer Lauf auf dem Host war), die Bestätigung des Wortlauts der Zwischenregel durch den Nutzer (Datum, Fundstelle), die Größe in der Metrik von §4 (Basis, letzter Commit, Summe; kein Faktor, keine Schätzung), und der Ausgang des Registers.

### Belege des Implementers (2026-10-24, Zweig impl-last1)

**Host und Voraussetzung.** Linux, 16 Kerne (`nproc`; der Plan nennt 20 aus den Belegen anderer Hosts), Docker; vor der Messung kein Container, Last 1,2 (`uptime` nach dem Start). Kein anderer Rollenlauf, keine Gegenprobe und kein `make gates` auf dem Host oder Daemon während der Messung (Voraussetzung §4). Zahl der Auftreten im Register zum Messzeitpunkt: 7×.

**(a) Die Last, wörtlich.** Je Lastlauf K Container, jeder mit einer Rechenschleife, im gepinnten Image des Repos (`postgres:17-alpine@sha256:b0f9560a…`, bereits vorhanden), ohne Netz, bis zum Ende der Serie:

```
for k in $(seq K); do docker run -d --name ml-<serie>-last-$k --network none \
  postgres:17-alpine@sha256:b0f9560a2de083e2cc7382e75f808c7381a32852a7ec49117deedb300e552b24 \
  sh -c 'while :; do :; done' >/dev/null; done; sleep 5
```

K = 32 (zweimal die Kerne) und K = 48 (dreimal). Last gemessen mit `cut -d' ' -f1-3 /proc/loadavg` vor, nach und nach jedem Lauf. Wichtig: **Eine Schleife je Container ist nötig.** Ein erster Versuch mit 32 und 48 Schleifen in *einem* Container (`sh -c 'for i in $(seq K); do (while :; do :; done) & done; wait'`) belastete die Tests nicht (Läufe 0,9 bis 1,3 s wie allein): Der Scheduler teilt die Kerne zwischen Container-cgroups, nicht zwischen Threads, und ein Container mit 48 Schleifen bekommt nicht mehr als ein Test-Container. Die Läufe dieses Versuchs zählen nicht (eine Rotung darin, Test C, Last 35 bis 42, fiel mit derselben Meldung wie unten).

**Aufruf der Tests** (Image `docker build --target integration -t pgwire-recorder:integration .`, PostgreSQL wie `make test-integration`: `postgres:17-alpine@sha256:b0f9560a…` mit `POSTGRES_HOST_AUTH_METHOD=trust` in einem `--internal`-Netz; je Lauf ein frischer Container des Testimages):

```
docker run --rm --network <netz> -e PGR_UPSTREAM=<pg>:5432 pgwire-recorder:integration \
  -test.run '^TestE2EPlayVergleichAbweichung$/^Serverfehler_ohne_Vorbild$'   # Test A
  -test.run '^TestE2EPlayVergleichAlle$/^Umgebungsvariable$'                  # Test B
```

Test C (`TestRunPlayZweitesSignal`) ist ein Unit-Test und braucht kein PostgreSQL: das Testbinary einmal gebaut (`go test -c -trimpath -buildvcs=false -o bootstrap.test ./internal/bootstrap` im Image der Stufe `source`), je Lauf `docker run --rm --network none -v <dir>:/t:ro -w /t <source-image> /t/bootstrap.test -test.count=1 -test.v -test.run '^TestRunPlayZweitesSignal$'`. Die Skripte liegen außerhalb des Repos (Scratchpad der Sitzung) und sind nicht committet (§1: kein Werkzeug für Last). **Grenze der Wiederholbarkeit (Review F-5):** Der Aufruf steht mit Platzhaltern (`<serie>`, `K`, `<netz>`, `<pg>`) und Bereichen der Last (16 bis 32, 18 bis 48; kein Wert vor und nach je Serie), der Start von PostgreSQL nur als „wie `make test-integration`"; der Messhost ist abgeschaltet, die Messung ist nicht wörtlich wiederholbar. DoD 1 (a) ist damit nur als Beschreibung erfüllt, nicht als Skript; das Risiko, dass das Rot unter dieser Last nicht herzustellen ist, bleibt *weiter offen* → `BEO-REPO/test-faellt-unter-cpu-last-aus`.

**(b) Läufe allein und unter Last** (Dauer je Lauf; rot = Exit ungleich 0):

| Test | allein | K=32 | K=48 |
|---|---|---|---|
| A `TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild` | 5 Läufe, 0,82 bis 0,98 s, 0 rot | 15 Läufe, 1,2 bis 1,6 s, 0 rot | 15 Läufe, 0 rot |
| B `TestE2EPlayVergleichAlle/Umgebungsvariable` | 5 Läufe, 0,91 bis 1,16 s, 0 rot | 15 Läufe, 1,3 bis 2,1 s, 0 rot | 15 Läufe, 0 rot |
| C `TestRunPlayZweitesSignal` | 5 Läufe, 0,29 bis 0,35 s, 0 rot | 20 Läufe, **8 rot** | 20 Läufe, **9 rot** |

Last in den Serien K=32/48: 16 bis 32 bzw. 18 bis 48 (Lastzahl 1 min, 16 Kerne). Dazu die Läufe aus dem ersten Versuch (ein Container, 32 und 48 Schleifen: je 15 Läufe, A und B 0 rot, C 0 rot bei 32 und 1 rot bei 48).

**Rotes Bild von C**, in allen 18 roten Läufen gleich: `play_test.go:341: binnen 30 s keine Client-Nachricht 'X', bisher ""`, `--- FAIL: TestRunPlayZweitesSignal (30.04s)`. Die erste Wartefrist (`play` endet binnen 30 s nach dem zweiten Signal) wird *nicht* gerissen — `Run` kehrt mit Exit-Code 0 zurück, der Lauf dauert 0,5 s, danach wartet `typenBis` 30 s auf `X`, und es kommt keine Nachricht: Der Server sieht nach `Q` kein `Terminate`.

**(c) Proben in einer Kopie außerhalb des Repos** (Kopie von `go.mod`, `go.sum`, `cmd`, `internal`, `test` im Scratchpad; Testbinary dort gebaut; der Arbeitsbaum des Repos blieb unberührt):

| Probe | trennt | Ergebnis |
|---|---|---|
| P1: in `einspielSession.Anfrage` (`internal/adapters/driven/postgres/einspielen.go`) 300 ms `time.Sleep` zwischen `Flush` und dem Ende der Funktion, also die Sperre `schreiben` nach dem Senden gehalten; Test C allein, ohne Last | Wettlauf im Code (Sperre noch gehalten, wenn das Signal kommt) von Frist | 5 Läufe, **5 rot**, dieselbe Meldung wie unter Last |
| P2: wie P1, zusätzlich `s.schreiben.TryLock()` in `schliesse` durch `Lock()` ersetzt | dass `TryLock` der Mechanismus ist | 5 Läufe, 0 rot |
| P3: `DROP DATABASE` auf eine Datenbank, zu der eine Sitzung noch verbunden ist (`psql` mit `pg_sleep(60)`), im Image des Repos | Vermutung des Planners zu A und B: wartet `DROP DATABASE` bis zur Frist von 30 s? | Fehler `database "vgl_abw" is being accessed by other users` nach 5,16 s; kein 30-s-Warten |
| P4: `DROP DATABASE IF EXISTS` und `CREATE DATABASE` je 8-mal, allein und neben 8 Containern mit `dd … oflag=dsync` | Frist durch E/A-Last | 0,21 s allein, 0,12 s mit E/A-Last (dominiert vom `docker exec`): kein Unterschied |

**(d) Ursache.**

- **C `TestRunPlayZweitesSignal`: Wettlauf im Code.** `schliesse` sendet `Terminate` nur, wenn `s.schreiben.TryLock()` gelingt (`einspielen.go`, „sendet Terminate nur, wenn gerade nichts sendet“); `Anfrage` hält `schreiben` bis nach `Flush` (`defer Unlock`). Der Test wartet, bis der Server `Q` gelesen hat, und löst dann beide Signale aus; zwischen dem Schreiben auf den Socket und dem `Unlock` des Sendenden kann der Scheduler ihn lange genug aufhalten, dass `Schliesse` die Sperre gehalten findet und die Verbindung ohne `Terminate` schließt. Das Rot ist reproduzierbar ohne Last (P1), die Last erhöht nur die Wahrscheinlichkeit des Fensters (8 und 9 von 20). Es ist keine Frist: 30 s sind für die Arbeit des Tests ausreichend (0,3 s). Die Lesart der Zusage ist entschieden: Commit 36b6387c schärft `LH-FA-20.a` (nach vollständigem Schreiben wird `Terminate` gesendet); der Code bricht sie, das Rot ist ein Befund. Behoben wird hier nichts (§1), das ist ein Folge-Slice.
- **A `…/Serverfehler_ohne_Vorbild` und B `…/Umgebungsvariable`: nicht reproduziert** in je 5 Läufen allein und je 30 Läufen unter Last (je 15 bei K=32 und K=48 mit einer Schleife je Container, Tabelle (b); die 30 Läufe des ersten Versuchs mit einem Container zählen nicht, siehe (a)), Last bis 48 auf 16 Kernen. Die Läufe unter Last dauern höchstens 2,1 s, das 30-s-Gebot liegt um den Faktor 15 darüber; Last allein auf der CPU erzeugt das Rot nicht. Die Vermutung des Planners (`DROP DATABASE` wartet auf die Verbindung der Vorgänger-Unterfälle) ist für die beobachtete Meldung nicht tragend: P3 zeigt 5 s und eine Fehlermeldung statt 30 s und „context deadline exceeded“. Ursache bleibt offen. **Reichweite der Zwischenregel (Hinweis des Reviews F-3):** Der Klammersatz „(Frist von 30 s unter CPU-Last)“ im Wortlaut der Regel beschreibt das beobachtete Rot der Belege; dass Last allein dieses Rot für A und B erzeugt, belegt §7 nicht (0 von 30 unter Last, Ursache offen). Der Wortlaut ist nutzerbestätigt und bleibt; die Reichweite der Aussage für A und B ist damit nicht belegt.

**Weiter offen für A und B:** das Rot der Belege (Last auf einem Host, der zugleich `docker build` und ein volles `make gates` fährt) ist mit reiner CPU-Last nicht herzustellen; eine Messung neben einem laufenden `make gates` wäre die nächste Probe und hätte den Host belegt (Voraussetzung §4 verbietet es hier).

**Signalfall der Gegenprobe (Punkt 2).** `GP_MUT=m-p17-signal-int make container-einrichten-gegenprobe`: allein 1 Lauf, 30 s, grün. Unter K=32 Containern mit je einer Schleife: 8 Läufe, 26 bis 27 s je Lauf, Zeile „Zeit 27 s, Last 16,6 bis 32,0“, 0 rot. **Nicht reproduziert in 8 Läufen bei Last 16 bis 32 auf 16 Kernen.** Abweichung von Randform 3 (mindestens 15 Läufe, Review F-4): 8 statt 15, weil ein Lauf 27 s dauert; der Befund trägt entsprechend weniger. Einschränkung: Die Gegenprobe läuft als Prozesse des Hosts, nicht in einer Container-cgroup; die Schleifen in Containern nehmen ihnen nicht denselben Anteil, die Läufe wurden nicht langsamer (27 s statt 30 s). Ein Rot durch das Zeitbudget ist nicht aufgetreten (mit `GP_MUT` wird es nicht geprüft, SPEC-060 Punkt 16 (4)); die 21 Läufe aus `slice-harness-devcontainer-einrichten` bleiben der einzige Beleg. Die Gegenprobe wurde nicht geändert.

**Zwischenregel (Punkt 3).** `implement-slice.md` Schritt 16 trägt die Regel in dem Wortlaut, den der Nutzer am 2026-10-24 bestätigt hat. **Bestätigung des Nutzers:** 2026-10-24, Wortlaut des Nutzers „Vorschlag 2 (Recommended)“, so vom Koordinator im Auftrag an den Implementer übermittelt; Vorschlag 2 aus dem Bericht der Messsession: (1) die Regel nur für `TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild` und `TestE2EPlayVergleichAlle/Umgebungsvariable`, (2) „Rot von `TestRunPlayZweitesSignal` ist ein Befund (Wettlauf im Code, reproduzierbar), keine Wiederholung, bis ein Folge-Slice es behebt.“ Die Regel ist Prosa ohne Test (§3; `AGENTS.md` §3.10 gilt für Code und Gates) und gilt für den Implementer, nicht für Reviewer und Verifier. Risiko 4 (Regel verdeckt Wettlauf): für `TestRunPlayZweitesSignal` durch den Wortlaut aufgefangen; den Ausgang setzt die Closure.

**Größe** (`AGENTS.md` §3.14, Basis `ffdacf7d41a3d551a8d164470e070a8a0d2791e0`, Pfade `.claude internal test tools harness`): der Befehl aus §3.14 liefert am Kopf des Zweigs (Stand nach dem Review-Report, e1382c88) **1**: die eine Regelzeile in `.claude/commands/implement-slice.md` (Commit 00bf41a6); `.claude` zählt. Unter `internal test tools harness` ist nichts geändert. Der Wert 0 des Eintrags davor galt vor diesem Commit. Marke 1500 nicht erreicht.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-22). Treffer in der Sub-Area `REPO`: `BEO-REPO/test-faellt-unter-cpu-last-aus` (6×, Schwelle erreicht; Ausgang `geplant` mit diesem Slice) — der Auslöser; abgegrenzt `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` (3×, `geplant` mit `slice-harness-lint-warten-ohne-frist`: dort fehlt die Frist) und `BEO-REPO/feste-wartezeit-statt-beobachtung-im-test` (2×, verwandt: eine feste Wartezeit im Test); `BEO-REPO/lastabhaengige-zaehlung` (2×) und `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (10×, `geplant` mit `slice-harness-randform-umgebung-der-gegenprobe`) — Last als Größe der Gegenprobe; `BEO-REPO/lauf-ignoriert-sigint-nach-nohup` (2×) — der Fall `p17-budget-sigint` ist ein Signalfall, die Zuordnung zu dieser Beobachtung ist Urteil des Planners (§6). Eine Schwelle wird mit diesem Slice nicht neu erreicht (keine Evidence-Datei entsteht); die Läufe der Messung sind Beleg des Slice, kein Auftreten im Register, solange kein Rot ohne Last entsteht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-22, im Commit, der diesen Plan anlegt und die Sendung von `slice-harness-randform-umgebung-der-gegenprobe` einträgt):** Liefer-Punkte 3 (Ursache der zwei Go-Tests | zweite Probe: Signalfall der Gegenprobe | Zwischenregel). Schichten 2: Tests des Produkts (Bootstrap und Integrationstests, gelesen und gefahren, Proben in einer Kopie, kein Diff; konservativ als eine Schicht gezählt) und Harness-Artefakte (Gegenprobe gefahren, `implement-slice.md` geändert); Planung zählt keine Schicht, kein Produkt-Code im Diff. Größe: Marke 1500 nach §4 (die Schätzung und die Meldeschwelle der Anlage vom 2026-10-22 entfallen nach `AGENTS.md` §3.14). Die Sendung von `slice-harness-randform-umgebung-der-gegenprobe` (Aufruf der Last wörtlich in §7) steht in DoD 1 und legt keinen Liefer-Punkt und keine Schicht an. Läge der Slice mit der Zeit der Messung über einem Rollenlauf, ist der Schnitt vorgemerkt: Punkt 2 als eigener Slice (§4).

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der die Ortsangabe aus `slice-harness-gegenproben-umzug-dev` einträgt): Die Gegenproben liegen seit jenem Slice unter `tools/dev/<bereich>/`, die Pfade der Gegenproben in diesem Plan sind nachgezogen (`tools/dev/` ist wie `tools/harness/` Harness-Werkzeug, keine neue Schicht). Ortsangabe der Gegenprobe geändert, Liefer-Punkte (3) und Schichten (2) unverändert.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-24, im Commit, der den Plan zum Anspruch nachzieht; Sendung des siebten Auftretens `BEO-REPO/test-faellt-unter-cpu-last-aus`: `TestE2EPlayVergleichAlle/Umgebungsvariable`): Die Sendung erweitert den Gegenstand von Liefer-Punkt 1 um einen dritten Test derselben Datei und derselben Ursache-Klasse (Frist 30 s, `DROP DATABASE`); ein neuer Liefer-Punkt entsteht nicht. Liefer-Punkte weiter 3 (Ursache der Go-Tests | zweite Probe: Signalfall der Gegenprobe | Zwischenregel); Schichten weiter 2 (Tests des Produkts, gelesen und gefahren, kein Diff; Harness-Artefakte). Der Nehmer liegt nicht darüber. Die Größenmarke ist auf 1500 nach §3.14 nachgezogen (§4). Bestand des Registers beim Start neu zu messen: Der Implementer nennt in §7 die Zahl der Auftreten zum Messzeitpunkt (Stand dieses Nachzugs: 7×).
