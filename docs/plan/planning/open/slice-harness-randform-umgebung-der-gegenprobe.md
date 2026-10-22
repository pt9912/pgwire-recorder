# Slice slice-harness-randform-umgebung-der-gegenprobe: Die Umgebung einer Gegenprobe ist eine Randform (`AGENTS.md` §3.12), und der Implementer fährt sie vor der Übergabe einmal unter fremder Locale, Zeitzone und fremdem Benutzer

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-harness-regeln, **hinter** `slice-harness-regeln-spezifikation-satzabgleich` (§4 und §5 der Welle). **Zuordnung und Stelle sind Empfehlung des Planners vom 2026-10-22** (Auftrag des Koordinators; Welle geprüft nach `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht: Die Welle trägt Regel-Slices mit Einträgen in `AGENTS.md` und `.claude/commands/`, dieser ist einer davon), **nicht Entscheidung des Nutzers**. Entschieden hat der Nutzer, dass es den Slice gibt: „deine empfehlungen“ vom 2026-10-22, so vom Koordinator übermittelt, auf die Vorschläge des Planners in `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`; über Welle, Stelle und Wortlaut ist nichts entschieden (`BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet`). Angelegt in `open/`, weil er nicht priorisiert ist; den Anspruch trägt der Nutzer.

**Bezug:** Kein Anspruch des Produkts — der Slice ändert Regeln des Harness. Bindung an das Register: `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (10×, Ausgang `geplant` mit diesem Slice seit 2026-10-22); verwandt, je mit einer Ursache der Klasse: `BEO-REPO/make-ausgabe-locale-abhaengig-gelesen`, `BEO-REPO/lastabhaengige-zaehlung`, `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig`. Regelstellen: `AGENTS.md` §3.12 (seit slice-harness-randformen-vor-code) und `.claude/commands/implement-slice.md`. Aktive ADR: [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (die Fahrt ist bei Gegenproben mit Zeitgrenze der engste Fall).

**Berührte Spec-Stellen:** der Slice ändert keinen Vertrag; gelesen sind [`SPEC-060`](../../../../harness/vertraege/SPEC-060-gegenproben.md) (Zeitbudget und Aufrufarten der Gegenproben) und [`SPEC-061`](../../../../harness/vertraege/SPEC-061-gegenproben-nach-einlieferung.md) (Beispiel einer Voraussetzung im Kopf: Benutzer root).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-22.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Umgebung, in der eine Gegenprobe läuft, ist eine benannte Randform, und der Implementer hat sie vor der Übergabe einmal gefahren: (1) `AGENTS.md` §3.12 trägt als weitere „gehört dazu“-Randform die *Umgebung der Gegenprobe*: Liefert ein Slice eine Gegenprobe oder ändert er eine, nennt §6 die Umgebung, in der sie läuft und in der sie gefahren wird (Locale, Zeitzone, Benutzer, Last, Werkzeuge des Hosts), und je Größe eine Entscheidung — im Aufruf der Gegenprobe festgelegt, oder als Voraussetzung im Kopf der Gegenprobe und im Vertrag genannt; (2) `.claude/commands/implement-slice.md` verlangt vom Implementer, die Gegenprobe vor der Übergabe **einmal** unter `LANG=de_DE.UTF-8`, `TZ=Europe/Berlin` und als ein Benutzer zu fahren, der nicht der des Schreibers ist, und die Läufe samt dem Beleg, dass die Umgebung wirkte, in §7 zu nennen; und vom Verifier, die betroffene Gegenprobe in seinem Lauf gezielt unter fremder Umgebung zu fahren. Das ist der Vorschlag des Planners aus dem Register, übernommen und klein geschnitten: zwei Liefer-Punkte, zwei Dateien.

**Beleg-Stand** (Stand 2026-10-22, gemergter Stand des Registers; die Belege stehen dort, nicht hier): Zehn Vorgänge in `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts/evidence/`, jeder war in der Umgebung des Schreibers grün und in einer anderen rot oder grün aus dem falschen Grund. Die Größen, die sie trafen: Sprache von `make`, Werkzeuge des Hosts, Last, Zeitzone, Benutzer root (fünf Vorgänge), Aufrufer des Hooks (Deskriptoren, Signale), Dateirechte und `umask`, freier Platz des Datenträgers und das Verhalten gerufener Werkzeuge. Locale und Zeitzone lassen sich mit den Variablen eines einzigen Aufrufs fremd stellen; der Benutzer verlangt einen Wechsel im Container. Gefunden hat die Fälle bisher das Review oder die Verifikation, nach der Übergabe.

**Entwurf des Wortlauts und Entscheidung des Nutzers.** Regeländerungen in `AGENTS.md` und in `.claude/commands/` sind eine Entscheidung des Nutzers. Der Nutzer hat am 2026-10-22 („deine empfehlungen“, so vom Koordinator übermittelt) die Vorschläge des Planners im Register gewählt, die Regel und die Fahrt; **den Wortlaut hat er nicht entschieden**. Der Implementer legt den Wortlaut aus §3 (oder seine Überarbeitung) dem Nutzer vor und committet `AGENTS.md` und `implement-slice.md` erst nach der Bestätigung; bis dahin ist er eine Empfehlung. Was der Nutzer nicht bestätigt, nimmt der Implementer nicht auf und nennt es in §7.

**Sendung an `slice-test-unter-last-ursache`** (eingetragen vom Planner im Commit des Anlegens, `AGENTS.md` §3.13): Die Fahrt unter Last ist unten ausgeschlossen, weil niemand eine Last wiederholbar beschrieben hat. Der Nehmer nennt in §7 den Aufruf der künstlichen Last wörtlich (Erzeuger, Stärke, Dauer; sein DoD 1), und dieser Slice nennt ihn als die Quelle, auf die eine Regel zur Fahrt unter Last sich stützen müsste. `grep` nach der Kennung dieses Slice findet die Sendung in §1 und DoD des Nehmers.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die **Fahrt unter Last**: Last ist als Größe in der Randform (1) genannt, mit einer Entscheidung des Plans je Gegenprobe; eine Fahrt der Gegenprobe unter Last gehört nicht in die Anweisung — ohne beschriebene Last ist sie nicht wiederholbar (die Beschreibung liefert `slice-test-unter-last-ursache`, Sendung oben), und sie verlängerte jede Übergabe um einen Lauf der Gegenprobe (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)). Eine Regel dazu entsteht, wenn der Nutzer sie nach dieser Beschreibung will; einen Slice dafür gibt es nicht, der Planner legt ihn auf Anweisung an.
- Ein **Sensor**, der prüft, ob eine Gegenprobe unter fremder Umgebung grün bleibt — anderer Vorgang (Werkzeug statt Regel): Eine solche Prüfung wäre ein weiterer Lauf der Gegenprobe selbst und verlängerte die Fahrt der Gegenproben; sie bekäme einen Vertrag in `harness/vertraege/` und eine Gegenprobe.
- Das **Nachfahren der Gegenproben des Bestands** unter fremder Umgebung — Bestand bleibt bewusst stehen: Die Regel gilt für die Gegenprobe, die ein Slice liefert oder ändert; nachgefahren wird, wenn ein Slice sie anfasst. Ein Rot des Bestands unter fremder Umgebung ist ein Auftreten im Register, kein Auftrag dieses Slice.
- Änderungen an `.claude/agents/` (Typen von Implementer und Verifier), `.harness/skills/reviewer.md`, `close-welle.md` und den Hooks — Schicht-Abgrenzung und Größe: Der Auftrag des Verifiers steht in `implement-slice.md` Schritt 22, die Typen zeigen darauf; eine dritte Datei brächte eine dritte Schicht (Agenten-Rollen) und einen dritten Liefer-Punkt.
- Die **Fahrt des Hooks im Gegenproben-Container** (Pflicht an Welle und Release, erste Fahrt) — `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`: ein anderer Vorgang. Dieser Slice betrifft die Fahrt des Implementers vor der Übergabe; die Fahrt im Container ist die Fahrstelle für den Stand der Welle.
- Produkt-Code, Gates, Gegenproben, Spezifikation, Benutzerhandbuch — Schicht-Abgrenzung: Der Slice ändert zwei Dokumente des Harness und kein Verhalten; **kein Produktcode**, keine Gegenprobe wird geändert.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Randform im Briefing):** `AGENTS.md` §3.12 trägt als weitere Randform die *Umgebung der Gegenprobe*: Wer eine Gegenprobe liefert oder ändert, nennt in §6 Locale, Zeitzone, Benutzer, Last und Werkzeuge des Hosts und je Größe eine Entscheidung (im Aufruf festgelegt, oder als Voraussetzung im Kopf der Gegenprobe und im Vertrag); dazu ein Paar Falsch/Richtig aus den Belegen, die Begründung mit `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` und der Retirement-Check `seit slice-harness-randform-umgebung-der-gegenprobe`. Der Wortlaut ist vom Nutzer bestätigt, **bevor** der Commit entsteht; Datum und Fundstelle stehen in §7.
- [ ] **Liefer-Punkt 2 (Fahrt in der Anweisung für Implementer und Verifier):** `.claude/commands/implement-slice.md` nennt in Schritt 15 oder 16 die Fahrt des Implementers (einmal unter `LANG=de_DE.UTF-8`, `TZ=Europe/Berlin` und als der andere Benutzer; bei den Gegenproben mit Zeitgrenze der engste Fall nach `AGENTS.md` §6 Schritt 6) und in Schritt 22 die gezielte Fahrt des Verifiers unter fremder Umgebung. Die Anweisung sagt die Regel an den Schreiber und zählt kein Verhalten der Gegenprobe auf (`AGENTS.md` §3.11). Belegt in §7 unter *Belege des Implementers*: je Fahrt Aufruf, Umgebung, der Beleg, dass sie wirkte (Ausgabe von `locale`, `date`, `id -un`), und das Ergebnis, gefahren an einer Gegenprobe des Bestands, die zur Umgebung einen Fall hat (Vorschlag: `make archiv-bezug-gegenprobe`, die einen Fall der Sprache trägt; der Implementer wählt und nennt sie). Derselbe Vorbehalt der Bestätigung des Wortlauts.
- [ ] Kein Wortlaut lockert oder ersetzt eine Zusage von §3.12 oder von `implement-slice.md`; der Diff ergänzt (Beleg: `git diff` ohne Löschzeile in den Regeln selbst, in §7).
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `.claude/` ändert).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Herkunfts-Anker `seit slice-harness-randform-umgebung-der-gegenprobe` in `AGENTS.md` §3.12 (`liegt in`), soweit aufgenommen.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` von `geplant` auf `verkörpert` (Zielort und Anker), soweit die Regel aufgenommen ist; ein nicht bestätigter Teil steht in §7 mit „nicht aufgenommen“ und der Eintrag bleibt `geplant` mit dem Rest; kein Zähler wird gesetzt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo mit Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit); die Folge-Slice-Paarung prüft Existenz **und** Annahme der Sendung an `slice-test-unter-last-ursache` (`grep` nach der Kennung dieses Slice in dessen §1 und DoD, `AGENTS.md` §3.13).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3.12 | update (ergänzen) | Randform *Umgebung der Gegenprobe* nach den bisherigen „gehört dazu“-Absätzen, mit Paar Falsch/Richtig, Begründung und Retirement-Check |
| `.claude/commands/implement-slice.md` | update | Schritt 15 oder 16 (Fahrt des Implementers vor der Übergabe) und Schritt 22 (Fahrt des Verifiers) |
| `docs/plan/planning/observations/BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts/state.md` | update | Ausgang bei Closure: `geplant` → `verkörpert` mit Zielort und Anker |

Tests und Gates: keine neuen. Der Slice liefert Prosa (`AGENTS.md` §3.10 verlangt Mutationen für Verträge von Code und Gates); die Wirkung zeigt das Register (Auftreten trotz der Regel).

**Entwurf des Wortlauts** (Empfehlung des Planners, keine Entscheidung des Nutzers; Nutzerentscheidung über den Wortlaut vor der Umsetzung):

> *`AGENTS.md` §3.12, nach den bisherigen Absätzen:* **Die Umgebung der Gegenprobe gehört dazu** (seit slice-harness-randform-umgebung-der-gegenprobe): Liefert ein Slice eine Gegenprobe oder ändert er eine, nennt §6 als Randform die Umgebung, in der sie läuft und in der sie gefahren wird (Locale, Zeitzone, Benutzer, Last, Werkzeuge des Hosts), und je Größe eine Entscheidung: im Aufruf der Gegenprobe festgelegt, oder als Voraussetzung im Kopf der Gegenprobe und im Vertrag genannt. Der Implementer fährt die Gegenprobe vor der Übergabe einmal unter `LANG=de_DE.UTF-8`, `TZ=Europe/Berlin` und als ein Benutzer, der nicht der des Schreibers ist, und nennt die Läufe mit dem Beleg, dass die Umgebung wirkte, in §7. Dazu ein Paar Falsch/Richtig und der Retirement-Check.
>
> *`implement-slice.md`, Schritt 15 oder 16, Ergänzung:* Eine Gegenprobe, die der Slice liefert oder ändert, fährst du vor der Übergabe einmal unter fremder Umgebung (`AGENTS.md` §3.12); bei den Gegenproben mit Zeitgrenze ist es der engste Fall (Schritt 15). §7 nennt Aufruf, Umgebung, den Beleg, dass sie wirkte (`locale`, `date`, `id -un`), und das Ergebnis.
>
> *Schritt 22, Ergänzung:* Der Verifier fährt die betroffene Gegenprobe gezielt unter fremder Locale, Zeitzone und fremdem Benutzer.

**Reihenfolge:** `AGENTS.md` zuerst (die Quelle), dann der Command, beide in einem Commit nach der Bestätigung; der Plan des Slice folgt in demselben Commit (`AGENTS.md` §3.9).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. Reihenfolge: **hinter** `slice-harness-regeln-spezifikation-satzabgleich` (**Empfehlung des Planners vom 2026-10-22**, nicht entschieden): Beide schreiben `AGENTS.md` §3.12, und der spätere rechnet den Wortlaut des ersten ein; die Stelle gegenüber den übrigen Slices der Reihe ist nicht festgelegt. WIP-Limit 1. Der Implementer legt den Wortlaut dem Nutzer vor und committet `AGENTS.md` und `implement-slice.md` erst nach der Bestätigung; bis dahin hält er an und gibt an den Aufrufer zurück.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice>` über `-- AGENTS.md .claude`; Basis ist der Elternteil des Commits, der diese Datei nach `in-progress/` bewegt. **Marke 100, Meldeschwelle 80:** Erreicht die Summe 80, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet; die Entscheidung steht als datierte Zeile in §7. Schätzung des Planners (2026-10-22, in dieser Metrik): `AGENTS.md` §3.12 mit Paar Falsch/Richtig, Begründung und Retirement-Check 15 bis 25, `implement-slice.md` 8 bis 15, zusammen **25 bis 40, das sind 25 bis 40 % der Marke**. Der Faktor Tests zu Produktcode entfällt: Es entsteht kein Code und keine Gegenprobe.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt über 100, oder die Regel verlangt eine dritte Datei (`.claude/agents/`, `.harness/skills/`) oder einen Sensor — zurück zur Zerlegung (§1).
- `in-progress` → `open` (blockiert): Der Nutzer bestätigt keinen Wortlaut und hat keine Fassung, die er bestätigen würde — der Slice wartet auf die Entscheidung, ohne dass ein Wortlaut in `AGENTS.md` steht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Bestätigung des Wortlauts durch den Nutzer in §7 mit Datum, die Fahrt aus DoD 2 in §7 belegt, Closure-Notiz mit Lerneintrag geschrieben (geschärfte Regel: `AGENTS.md` §3.12, `liegt in`, soweit aufgenommen).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen der Regel** (`AGENTS.md` §3.12; der Slice liefert keinen Vertrag von Code, aber eine Regel, deren Ränder der Planner sonst im Einzelfall entscheidet). Entscheidungsort ist der Wortlaut, vor dem Commit bestätigt:

1. *Welche Gegenprobe:* jede, die der Slice liefert oder ändert; der Bestand ist ausgenommen (§1).
2. *Welche Fahrt:* **ein** Lauf, der Locale, Zeitzone und Benutzer zugleich fremd stellt, nicht je Größe ein Lauf; bei den Gegenproben mit Zeitgrenze der engste Fall mit `GP_MUT` (`AGENTS.md` §6 Schritt 6, [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)), sonst die Gegenprobe ganz. Rot zeigt dann nur, dass eine der drei Größen versagt; der Wortlaut verlangt bei Rot, die Größen zu trennen.
3. *Benutzer:* der Vorschlag im Register sagt „als Benutzer ohne root“; die Belege zeigen root als die Größe, die ein Schreiber als `vscode` im Container nicht trägt (fünf der zehn Vorgänge). Der Entwurf nennt deshalb „ein Benutzer, der nicht der des Schreibers ist“; **Entscheidung des Nutzers beim Wortlaut.** Ist der andere Benutzer nicht herstellbar (kein Weg zu root im Container), nennt §7 es als nicht gefahren.
4. *Locale:* `de_DE.UTF-8` muss im Container vorhanden sein (`locale -a`); fehlt sie, fällt die Umgebung still auf die Standard-Locale zurück, und die Fahrt wäre keine fremde (`slice-harness-doc-trace-archiv`: ohne `en_US.UTF-8` entfielen Fälle). §7 trägt den Beleg (`locale`); fehlt die Locale, steht „nicht gefahren: Locale fehlt“.
5. *Zeitzone:* `TZ=Europe/Berlin` wirkt nur mit den Zeitzonendaten im Image; Beleg ist die Ausgabe von `date`.
6. *Last:* nicht Teil der Fahrt (§1); in der Randform (1) als Größe mit einer Entscheidung des Plans je Gegenprobe genannt.
7. *Werkzeuge des Hosts:* in (1) genannt; die Fahrt prüft sie nicht, die Gegenproben für den Host laufen mit `make gegenproben-host` am Host.
8. *Aufruf-Voraussetzung, Aufrufarten des Treibers, Plattform, Stelle einer Nachricht:* nicht berührt (`AGENTS.md` §3.12); die Regel liefert keinen Vertrag.

**Risiken:**

- Der Nutzer ändert den Wortlaut, streicht die Fahrt oder bestätigt keinen — **Ausgang:** offen bis Closure (erwartet: entfallen, der bestätigte Wortlaut steht in `AGENTS.md`; bei Ablehnung geht der Slice zurück nach `open/`, §4).
- Die Regel ist Prosa ohne Sensor: Eine Gegenprobe, die in fremder Umgebung rot ist, bleibt bis zum Review unbemerkt, wenn der Implementer die Fahrt nicht fährt — **Ausgang:** weiter offen → `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (ein Auftreten trotz der Regel zählt weiter; ein Sensor ist ausgeschlossen, §1).
- Die Fahrt verlängert die Übergabe um einen Lauf der Gegenprobe — **Ausgang:** weiter offen → `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, falls die Läufe in §7 Zeit zeigen, die der Nutzer als Zumutung liest; sonst entfallen mit der Zahl.
- Ein Lauf mit allen drei Größen zugleich färbt rot, ohne die Größe zu nennen (Randform 2) — **Ausgang:** entfallen, wenn der Wortlaut bei Rot das Trennen der Größen verlangt.
- Locale oder Zeitzone wirken im Container nicht, und die Fahrt ist still keine fremde (Randformen 4 und 5) — **Ausgang:** entfallen, wenn der Wortlaut den Beleg (`locale`, `date`) verlangt und §7 ihn trägt.
- Die Änderungen treffen auf `slice-harness-regeln-spezifikation-satzabgleich` (`AGENTS.md` §3.12, anderer Absatz), auf `slice-harness-regeln-nach-v1-abschluss` (`implement-slice.md` Schritte 17 und 19, hier 15 oder 16 und 22) und auf `slice-test-unter-last-ursache` (`implement-slice.md`, Zwischenregel) — **Ausgang:** offen bis Closure (erwartet: entfallen, die Absätze und Schritte überlappen nicht; der Slice, der später schließt, rechnet den Wortlaut des ersten ein).
- Der Slice wächst über Marke oder Meldeschwelle von §4 (Schätzung 25 bis 40 gegen Marke 100) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Hier stehen: die Bestätigung des Wortlauts durch den Nutzer (Datum, Fundstelle, je Teil), *Belege des Implementers* (die Fahrt aus DoD 2 mit Aufruf, Umgebung, Beleg der Wirkung und Ergebnis), die Größe in der Metrik von §4 (Basis, letzter Commit, Summe) mit dem Ergebnis der Schätzung gegen die Messung, und der Ausgang des Registers.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-22). Treffer in der Sub-Area `REPO`: `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (10×, Schwelle erreicht; Ausgang `geplant` mit diesem Slice) — der Auslöser; verwandt, je eine Ursache der Klasse: `BEO-REPO/make-ausgabe-locale-abhaengig-gelesen` (1×), `BEO-REPO/lastabhaengige-zaehlung` (2×) und `BEO-REPO/gegenprobe-mutation-faerbt-zusaetzlichen-fall-zeitabhaengig` (5×, `geplant` mit `slice-harness-gegenprobe-matrix-p7-flanke`); `BEO-REPO/spec-randform-erst-im-review-entschieden` (49×, `verkörpert` in `AGENTS.md` §3.12) — die Regel ergänzt dessen Zielort um eine Randform; `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (14×, `verkörpert`) — die Kostenseite der Fahrt. Eine Schwelle wird mit diesem Slice nicht neu erreicht (keine Evidence-Datei entsteht).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-22, im Commit, der diesen Plan anlegt und die Sendung an `slice-test-unter-last-ursache` einträgt):** Liefer-Punkte 2 (Randform im Briefing | Fahrt in der Anweisung für Implementer und Verifier). Schichten 2: Briefing (`AGENTS.md`) und Commands (`.claude/commands/`); Planung (Register, Plan) zählt keine Schicht; kein Produkt-Code, keine Nutzer-Doku. Schätzung 25 bis 40 in der Metrik von §4 (Marke 100, Meldeschwelle 80); der Faktor Tests zu Produktcode entfällt (kein Code, keine Gegenprobe). Der Slice trägt eine Sendung in `slice-test-unter-last-ursache` ein (Nachzählung dort in §8) und keine in einen anderen Nehmer.
