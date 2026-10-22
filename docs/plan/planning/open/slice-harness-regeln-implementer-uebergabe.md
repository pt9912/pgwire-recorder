# Slice slice-harness-regeln-implementer-uebergabe: Regeln der Übergabe des Implementers: Ein Folgelauf nach Review oder Verifikation, wenn die Nacharbeit neue ausführbare Zeilen trägt, und der Ort des Belegs eines Nachweises in den Anweisungen; dazu, welche Abschnitte des Plans der Implementer ändert

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-harness-regeln, **hinter** `slice-harness-randform-umgebung-der-gegenprobe` (§4 und §5 der Welle; beide ändern `.claude/commands/implement-slice.md`). **Zuordnung und Stelle sind Empfehlung des Planners vom 2026-10-22**, nicht Entscheidung des Nutzers. Entschieden hat der Nutzer die **Sache** in drei Teilen: am 2026-10-22 den Ort des Belegs eines Nachweises (auslagern statt kürzen; der Wortlaut steht in [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 (6) und [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 18), und am 2026-10-22 („deine empfehlungen“, so vom Koordinator übermittelt), dass `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` einen neuen Slice bekommt, weil kein vorhandener ohne Verstoß gegen den Schnitt annimmt; am 2026-10-22 („ja“ auf die Vorlage mit den Empfehlungen des Koordinators, so übermittelt), dass `BEO-REPO/implementer-ueberschreibt-plan-abschnitte` (3×) denselben Slice als Adresse bekommt, mit der Regel zu den Abschnitten des Plans. Den **Wortlaut** der Regeln hat der Nutzer nicht entschieden. Angelegt am 2026-10-22 in `open/` auf Anweisung der Koordination; nicht priorisiert.

**Bezug:** Kein Anspruch des Produkts (`LH-*`) — der Slice schreibt Regeln für die Übergabe zwischen Rollen. Bindung an das Register: `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` (3×, Ausgang `geplant` mit diesem Slice seit 2026-10-22), `BEO-REPO/ungekuerzter-beleg-laesst-den-plan-wachsen` (1×, Beleg-Ort); `BEO-REPO/implementer-ueberschreibt-plan-abschnitte` (3×, Ausgang `geplant` mit diesem Slice seit 2026-10-22). Regelstellen: `.claude/commands/implement-slice.md` (Schritt 18, Übergaben 21 und 22), `.claude/commands/plan-welle.md` (Schritt 6), `.claude/agents/implementer.md`, `.claude/agents/verifier.md`. Keine aktive ADR.

**Berührte Spec-Stellen:** [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 (6) und [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 18 (gelesen: der Ort des Belegs eines Nachweises; die Anweisungen zeigen dorthin und wiederholen den Vertrag nicht). Der Slice ändert keinen Vertrag.

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

**Ziel:** Die Übergabe des Implementers an Review und Verifikation trägt zwei Regeln in den Anweisungen: (1) **Ein Folgelauf nach der Nacharbeit** — Trägt die Nacharbeit nach einem Review oder einer Verifikation **neue ausführbare Zeilen** (Befehle, Variablen, Remote-Namen, Einträge in einen Abschnitt, den kein Lauf gelesen hat), bricht die Verifikation jede Behebung von Hand, und auf die Nacharbeit folgt ein **zweiter Review-Durchgang (`-r2`) oder eine Nachverifikation**, nicht der Schluss (`BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf`: drei Belege, in denen die Nacharbeit neue Zeilen schrieb, die weder ein Reviewer noch ein Verifier sah; sieben Gegenbelege in Folge zeigen, dass das Verhalten gilt, nicht, dass es gelten muss). (2) **Der Ort des Belegs eines Nachweises** — Die Sätze „Die Belege stehen in §7“ in `implement-slice.md` (Schritt 18), `implementer.md`, `verifier.md` und `plan-welle.md` (Schritt 6) bekommen die Ausnahme: Der Beleg eines Nachweises liegt als eigene, ungekürzte Datei unter `docs/reviews/`, und §7 nennt Pfad, Summe nach `sha256sum`, die Zeile mit dem Exit-Code und den Kopf (Entscheidung des Nutzers vom 2026-10-22; der Wortlaut des Vertrags steht in [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 (6) und [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 18). (3) **Welche Abschnitte des Plans der Implementer ändert** — Er ändert im Plan §7 und was `AGENTS.md` §3.9 nachzieht (§1, §3, §6, den Kopf, betroffene Folge-Slices), nicht §2, §4, §5, §8 und keinen Satz, der eine Entscheidung des Nutzers trägt; Bedarf geht im Bericht an den Planner (`BEO-REPO/implementer-ueberschreibt-plan-abschnitte`: drei Belege, in denen ein Commit des Implementers Abschnitte des Plans ersetzte oder einen Satz mit einer Entscheidung des Nutzers in §4 zufügte).

**Sendung von `slice-harness-devcontainer-einrichten`** (Geber des Beleg-Orts; eingetragen vom Planner im Commit des Anlegens, 2026-10-22, `AGENTS.md` §3.13): Die Ausnahme in den vier Dateien oben ist dieser Slice (Liefer-Punkt 2). Der Geber liegt in `done/`; die Adresse ist dieser Slice, `grep` nach der Kennung des Gebers findet die Sendung in §1 und DoD.

**Mitgenommen: die Regel zu den Abschnitten des Plans** (Ausgang `geplant` von `BEO-REPO/implementer-ueberschreibt-plan-abschnitte`, eingetragen vom Planner im Commit der Zuweisung, 2026-10-22, `AGENTS.md` §3.13; Entscheidung des Nutzers vom 2026-10-22, „ja“ auf die Vorlage mit den Empfehlungen des Koordinators, so übermittelt). Die Adresse ist dieser Slice: Die Regel steht in derselben Datei, `implement-slice.md`, und legt keine dritte Schicht an. Sein §1 schließt sie nicht aus (die Ausschlüsse treffen `AGENTS.md`, den Reviewer-Skill, Hooks und einen Sensor, nicht den Satz in `implement-slice.md`). Der Wortlaut ist Entwurf des Planners (§3); die Zählung steht in §8.

**Entwurf des Wortlauts und Entscheidung des Nutzers.** Regeländerungen in `.claude/commands/` und `.claude/agents/` sind eine Entscheidung des Nutzers; den Wortlaut hat er nicht entschieden. Der Implementer legt ihn aus §3 (oder seine Überarbeitung) dem Nutzer vor und committet erst nach der Bestätigung; was der Nutzer nicht bestätigt, nimmt er nicht auf und nennt es in §7.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die **Regel zur Stunde** (`AGENTS.md`, sechs Rollen-Typen) — Klasse 1: `slice-harness-rollenlauf-zeitgrenze-regel`. Grund: dieselben Dateien `implementer.md` und `verifier.md`, aber ein anderer Satz mit anderer Herkunft; der später schließende rechnet den Wortlaut des ersten ein.
- Die **Plan-Vorlage** (`slice.template.md`) und ihr Ort des Belegs — Bestand bleibt bewusst stehen (Klasse 2): Der Baum unter `.harness/baseline/` ist Kurs-Inhalt, byte-genau und in diesem Repo nicht autoritativ (`harness/conventions.md`); eine Änderung gehört dem Kurs-Repo.
- Die **Zeile `make homebrew-nachweis` in `harness/README.md`** und die **Kommentare** in `harness/mk/homebrew-nachweis.mk`, `tools/harness/homebrew-nachweis.sh` und `tools/harness/container-einrichten-gegenprobe.sh`, die „§7“ als Ort des Belegs nennen — Klasse 4 (Schicht-Abgrenzung): Gate-Index und Werkzeug-Kommentare sind keine Anweisung der Rollen. Die README-Zeile nimmt `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` an (dort §1, Sendung von `slice-harness-devcontainer-einrichten`); für die Kommentare gibt es keine Adresse.
- Ein **Sensor** für Folgeläufe oder für Abschnitte eines Plans — anderer Vorgang (Werkzeug statt Regel), mit Vertrag unter `harness/vertraege/` und Gegenprobe.
- Änderungen an `AGENTS.md`, `.harness/skills/reviewer.md` (dessen Namensregel für `-r2` bleibt), Hooks und Produkt-Code — Schicht-Abgrenzung und Größe: Beide Regeln stehen in `implement-slice.md` und in den Rollen-Dateien; eine weitere Schicht brächte einen vierten Liefer-Punkt.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Folgelauf nach der Nacharbeit):** `.claude/commands/implement-slice.md` nennt in den Übergaben (Schritte 21 und 22) die Regel: Trägt die Nacharbeit nach Review oder Verifikation neue ausführbare Zeilen, bricht die Verifikation jede Behebung von Hand, und ein zweiter Review-Durchgang oder eine Nachverifikation folgt; der Bericht der Nacharbeit nennt je Behebung den Lauf, der sie sah. Der Wortlaut ist vom Nutzer bestätigt, **bevor** der Commit entsteht; Datum und Fundstelle stehen in §7.
- [ ] **Liefer-Punkt 2 (Ort des Belegs eines Nachweises):** `.claude/commands/implement-slice.md` (Schritt 18), `.claude/agents/implementer.md`, `.claude/agents/verifier.md` (Eingang) und `.claude/commands/plan-welle.md` (Schritt 6) nennen neben „Die Belege stehen in §7“ die Ausnahme für den Beleg eines Nachweises (eigene Datei unter `docs/reviews/`; §7 nennt Pfad, Summe, Exit-Zeile und Kopf), (Sendung von `slice-harness-devcontainer-einrichten`, §1) mit dem Zeiger auf [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 (6) und [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 18; kein Satz sagt mehr zu als der Vertrag (`AGENTS.md` §3.11), namentlich nicht, dass `make archive-welle` die Datei archiviert (nicht zugesagt). Derselbe Vorbehalt der Bestätigung.
- [ ] **Liefer-Punkt 3 (Abschnitte des Plans):** `.claude/commands/implement-slice.md` nennt in Schritt 18 die Regel: Im Plan ändert der Implementer §7 und was `AGENTS.md` §3.9 nachzieht (§1, §3, §6, den Kopf, betroffene Folge-Slices); §2, §4, §5 und §8 und jeder Satz, der eine Entscheidung des Nutzers trägt, bleiben beim Planner, und Bedarf steht im Bericht (Sendung von `BEO-REPO/implementer-ueberschreibt-plan-abschnitte`, §1). Der Wortlaut ist vom Nutzer bestätigt, **bevor** der Commit entsteht; Datum und Fundstelle stehen in §7.
- [ ] Kein Wortlaut lockert oder ersetzt eine Zusage der vier Dateien; der Diff ergänzt (Beleg: `git diff` ohne Löschzeile in den Regeln selbst, in §7).
- [ ] `make docs-check` und `make kopf-check` grün.
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `.claude/` ändert).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Herkunfts-Anker `seit slice-harness-regeln-implementer-uebergabe` in `implement-slice.md` (`liegt in`), soweit aufgenommen.
- [ ] Beobachtungs-Register (`../observations/`): `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` und `BEO-REPO/implementer-ueberschreibt-plan-abschnitte` je von `geplant` auf `verkörpert` (Zielort und Anker), soweit die Regel aufgenommen ist; bei `BEO-REPO/ungekuerzter-beleg-laesst-den-plan-wachsen` ein Vermerk, dass die Anweisungen den Ort tragen; kein Zähler wird gesetzt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der nächsten Welle-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/commands/implement-slice.md` | update (ergänzen) | Übergaben (21, 22): Regel zum Folgelauf; Schritt 18: Ausnahme für den Beleg eines Nachweises und Regel zu den Abschnitten des Plans |
| `.claude/agents/implementer.md`, `.claude/agents/verifier.md` | update (ergänzen) | der Satz zum Ort der Belege bekommt die Ausnahme |
| `.claude/commands/plan-welle.md` | update (ergänzen) | Schritt 6: der Ort des Belegs, den ein DoD-Punkt verlangt |
| `docs/plan/planning/observations/BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf/` und `…/implementer-ueberschreibt-plan-abschnitte/` | update bei Closure | Ausgang `geplant` → `verkörpert` (kein Liefer-Punkt) |

Tests und Gates: keine neuen. Der Slice liefert Prosa (`AGENTS.md` §3.10 verlangt Mutationen für Verträge von Code und Gates); die Wirkung zeigt das Register (eine Nacharbeit mit neuen Zeilen ohne Folgelauf).

**Entwurf des Wortlauts** (Empfehlung des Planners, keine Entscheidung des Nutzers; Nutzerentscheidung über den Wortlaut vor der Umsetzung):

> *`implement-slice.md`, Übergaben (21, 22), Ergänzung:* **Eine Nacharbeit mit neuen ausführbaren Zeilen bekommt einen Folgelauf.** Trägt die Nacharbeit nach einem Review oder einer Verifikation neue ausführbare Zeilen (Befehle, Variablen, Remote-Namen, Einträge in einen Abschnitt, den kein Lauf gelesen hat), bricht die Verifikation jede Behebung von Hand, und auf die Nacharbeit folgt ein zweiter Review-Durchgang (`-r2`) oder eine Nachverifikation, nicht der Schluss. Der Bericht der Nacharbeit nennt je Behebung den Lauf, der sie sah. Reine Verengung von Prosa braucht keinen Folgelauf.
>
> *Schritt 18, `implementer.md`, `verifier.md`, `plan-welle.md`, Ergänzung zu „Die Belege stehen in §7“:* Der Beleg eines Nachweises (`SPEC-067` Punkt 20, `SPEC-069` Punkt 18) liegt als eigene Datei unter `docs/reviews/`; §7 nennt ihren Pfad, die Summe nach `sha256sum`, die Zeile mit dem Exit-Code und den Kopf, nicht die Ausgabe.
>
> *Schritt 18, nach dem Satz zu den Belegen:* **Im Plan ändert der Implementer §7 und was `AGENTS.md` §3.9 nachzieht** (§1, §3, §6, den Kopf, betroffene Folge-Slices). §2, §4, §5 und §8 gehören dem Planner, ebenso jeder Satz, der eine Entscheidung des Nutzers trägt; der Implementer ersetzt und ergänzt keinen von ihnen. Braucht eine dieser Stellen eine Änderung, steht sie im Bericht an den Planner. Ein Commit, der eine Plandatei ersetzt, behält jeden Abschnitt, den er nicht ändern darf.

**Reihenfolge:** `implement-slice.md` zuerst, dann die drei übrigen Dateien, alles in einem Commit nach der Bestätigung; der Plan des Slice folgt in demselben Commit (`AGENTS.md` §3.9).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. Reihenfolge: **hinter** `slice-harness-randform-umgebung-der-gegenprobe` (**Empfehlung des Planners**): beide ändern `implement-slice.md`, der spätere rechnet den Wortlaut des ersten ein. WIP-Limit 1. Der Implementer legt den Wortlaut dem Nutzer vor und committet erst nach der Bestätigung; bis dahin hält er an und gibt an den Aufrufer zurück.

**Größe.** Metrik nach [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (der Befehl steht in `AGENTS.md` §3.14 nach `slice-harness-groessen-metrik`): Summe aus Einfügungen und Löschungen über den Pfad `.claude`, mit den Ausschlüssen des Befehls; Basis ist der Commit, der diese Datei nach `in-progress/` bewegt. **Marke 100, Meldeschwelle 80.** Schätzung des Planners (2026-10-22, in dieser Metrik): Regel zum Folgelauf 8 bis 15, Ort des Belegs in vier Dateien je 4 bis 8, Regel zu den Abschnitten des Plans 10 bis 20, zusammen **34 bis 67, das sind 34 bis 67 % der Marke**. Der Faktor Tests zu Produktcode entfällt: kein Code, keine Gegenprobe.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt über 100, oder eine Regel verlangt eine dritte Schicht (`AGENTS.md`, einen Hook oder einen Sensor).
- `in-progress` → `open` (blockiert): Der Nutzer bestätigt keinen Wortlaut und hat keine Fassung, die er bestätigen würde — der Slice wartet auf die Entscheidung, ohne dass ein Wortlaut in den Dateien steht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, `make gates` grün, Bestätigung des Wortlauts durch den Nutzer in §7 mit Datum, `grep` nach der Kennung des Gebers der Sendung findet die Ausnahme in den vier Dateien, Closure-Notiz mit Lerneintrag geschrieben (geschärfte Regel: `implement-slice.md`, `liegt in`, soweit aufgenommen).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen der Regeln** (`AGENTS.md` §3.12; der Slice liefert keinen Vertrag von Code, aber Regeln, deren Ränder der Planner sonst im Einzelfall entscheidet). Entscheidungsort ist der Wortlaut, vor dem Commit bestätigt:

1. *„Neue ausführbare Zeile“:* ein Befehl, eine Variable, ein Remote-Name, ein Pfad in einem Beispiel; die Verengung einer Aussage (der Befund war „zu weit“, die Behebung fasst es enger) ist keine. Der Wortlaut nennt die Beispiele; ob eine Behebung zählt, urteilt der Verifier.
2. *Welcher Folgelauf:* zweiter Review-Durchgang **oder** Nachverifikation; wer wählt (der Koordinator, der Planner), ist **Entscheidung des Nutzers beim Wortlaut**.
3. *Name des Berichts:* `-r2` nach dem Reviewer-Skill; für die Nachverifikation gilt die Namensform der Berichte unter `docs/reviews/` (bisherige Berichte tragen `nachverifikation` im Namen).
4. *Welcher Beleg die Ausnahme trägt:* der Beleg eines Nachweises nach [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 und [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) Punkt 18; Mutationstabellen, Befundzeilen und Läufe bleiben in §7.
5. *Abschnitte des Plans (Liefer-Punkt 3):* (a) „was `AGENTS.md` §3.9 nachzieht“ ist der Nachzug von §1, §3, §6 und Kopf, den die Änderung am Code, an Gates, Tests, Spezifikation oder Harness-Artefakten verlangt, kein neuer Inhalt; (b) die Häkchen in §2 setzt der Planner bei der Closure (`AGENTS.md` §3.3), der Implementer behauptet die DoD im Bericht (Schritt 18) — **Entscheidung des Nutzers beim Wortlaut**; (c) „ein Satz, der eine Entscheidung des Nutzers trägt“ ist jeder Satz mit „Entscheidung des Nutzers“ und Datum im Plan; (d) Bedarf an einer Änderung eines geschützten Abschnitts geht im Bericht an den Planner, ein Satz im Plan ersetzt die Rückgabe nicht (wie bei der Randform-Rückgabe).
6. *Aufruf-Voraussetzung, Aufrufarten des Treibers, Plattform, Stelle einer Nachricht:* nicht berührt (`AGENTS.md` §3.12); keine der beiden Regeln liefert einen Vertrag.

**Risiken:**

- Der Nutzer ändert den Wortlaut, streicht einen Teil oder bestätigt keinen — **Ausgang:** offen bis Closure (erwartet: entfallen, der bestätigte Wortlaut steht in den Dateien; bei Ablehnung geht der Slice zurück nach `open/`, §4).
- Die Regel zum Folgelauf ist Prosa ohne Sensor: Eine Nacharbeit mit neuen Zeilen bleibt ohne Folgelauf, wenn der Lauf sie nicht liest — **Ausgang:** weiter offen → `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` (ein Auftreten trotz der Regel zählt weiter; ein Sensor ist ausgeschlossen, §1).
- Ein Folgelauf nach jeder Nacharbeit verlängert die Closure um einen Lauf — **Ausgang:** entfallen, wenn der bestätigte Wortlaut den Folgelauf auf Nacharbeit mit neuen ausführbaren Zeilen begrenzt (Randform 1).
- Der Satz zum Ort des Belegs widerspricht dem Satz „Die Belege stehen in §7“ (seit slice-harness-blackbox-kern) — **Ausgang:** entfallen, wenn der Wortlaut die Ausnahme als Ergänzung zum Satz fasst und `git diff` keine Löschzeile in der Regel zeigt.
- Die Änderungen an `implement-slice.md` treffen auf `slice-harness-randform-umgebung-der-gegenprobe`, `slice-test-unter-last-ursache` und `slice-harness-regeln-nach-v1-abschluss` (andere Schritte), an `implementer.md` und `verifier.md` auf `slice-harness-rollenlauf-zeitgrenze-regel` — **Ausgang:** offen bis Closure (erwartet: entfallen, andere Schritte und Sätze; der später schließende rechnet den Wortlaut des ersten ein).
- Die Kommentare in `harness/mk/homebrew-nachweis.mk` und zwei Skripten nennen „§7“ weiter und haben keine Adresse — **Ausgang:** weiter offen → `BEO-REPO/ungekuerzter-beleg-laesst-den-plan-wachsen`.
- Die Regel zu den Abschnitten ist Prosa ohne Sensor: Ein Commit des Implementers ersetzt weiter Abschnitte, wenn der Lauf die Regel nicht liest — **Ausgang:** weiter offen → `BEO-REPO/implementer-ueberschreibt-plan-abschnitte` (ein Auftreten trotz der Regel zählt weiter; ein Sensor ist ausgeschlossen, §1: der Diff gegen `main` müsste die Menge aus `AGENTS.md` §3.9 zulassen, nicht „nur §7“).
- Der Slice wächst über Marke oder Meldeschwelle von §4 (Schätzung 34 bis 67 gegen Marke 100) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Hier stehen: die Bestätigung des Wortlauts durch den Nutzer (Datum, Fundstelle, je Teil), die Größe in der Metrik von §4 (Basis, letzter Commit, Summe) mit dem Ergebnis der Schätzung gegen die Messung, der Ausgang des Registers und die Ausgänge der Risiken.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-22). Treffer in der Sub-Area `REPO`: `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` (3×, Schwelle erreicht; Ausgang `geplant` mit diesem Slice) — der Auslöser; `BEO-REPO/ungekuerzter-beleg-laesst-den-plan-wachsen` (1×) — der Beleg-Ort; `BEO-REPO/implementer-ueberschreibt-plan-abschnitte` (3×, Schwelle erreicht; Ausgang `geplant` mit diesem Slice) — der Auslöser von Liefer-Punkt 3; `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht` (`verkörpert`, die Regel „Die Belege stehen in §7“ ist ihr Zielort) — die Ausnahme ergänzt sie; `BEO-REPO/slice-waechst-durch-uebernahmen` (`verkörpert`) — die Sendung und die Regel zu den Abschnitten sind Übernahmen, deshalb die Zählung unten. Eine Schwelle wird mit diesem Slice nicht neu erreicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 und §3.14 (Planner, 2026-10-22, im Commit, der diesen Plan anlegt und die Sendung von `slice-harness-devcontainer-einrichten` einträgt; nachgezählt am 2026-10-22 mit der Regel zu den Abschnitten des Plans):** **Liefer-Punkte 3** (Folgelauf nach der Nacharbeit | Ort des Belegs in den Anweisungen | Abschnitte des Plans); die Sendung liegt in Liefer-Punkt 2 und legt keinen an, die Regel zu den Abschnitten ist der dritte. Schichten 2, konservativ gezählt: Commands (`.claude/commands/`) und Agenten-Rollen (`.claude/agents/`); die Regel zu den Abschnitten steht in `implement-slice.md` und legt keine dritte an. Schätzung **34 bis 67** in der Metrik von §4 (Marke 100, Meldeschwelle 80), das sind 34 bis 67 % der Marke und unter 70 %; der Faktor Tests zu Produktcode entfällt (kein Code, keine Gegenprobe). Der Slice ist mit drei Liefer-Punkten **voll**: Jede weitere Übernahme geht an den Planner (`AGENTS.md` §3.13).
