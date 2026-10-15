# Slice slice-harness-schaetzung-planner-faktor: Der Planner schätzt beim Anlegen in der Metrik von §3.14, mit benanntem Faktor, schneidet ab 70 % der Marke und nennt die Meldeschwelle

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er von der nächsten Welle-Closure. Angelegt am 2026-10-15 nach Entscheidung des Nutzers vom selben Tag (der Lücken-Slice zu `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` und `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` wird angelegt; Marken und Meldeschwellen 1000/800 und 700/550 der beiden Teile von `slice-v1-abschluss-verbindungsende-format` bleiben); eingeordnet in `next/`, weil er mit einer Stelle in der Reihenfolge eingeplant ist (§4 *Start*, Drift-Log der Roadmap; die Stelle ist Empfehlung des Planners).

**Bezug:** kein Anspruch des Produkts — der Slice ändert Regeln des Harness. Bindung an das Register: `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (4×), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (5×), einbezogen `BEO-REPO/meldeschwelle-ohne-artefakt` (1×); Regelstelle: `AGENTS.md` §3.14 (seit slice-harness-gate-kette-stop-hook-liste). Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Schätzung des Umfangs liegt beim Anlegen eines Plans nicht mehr unter der des Architect: `AGENTS.md` §3.14 und `.claude/commands/plan-welle.md` verlangen vom Planner, (1) beim Anlegen und beim Schnitt in der Metrik von §3.14 zu schätzen und den **Faktor Tests zu Produktcode** mit seinem Beleg zu nennen, (2) bei einer Tabelle von Randformen mit einer Mutation je Element (§3.10) mit mindestens 5 zu rechnen (`slice-v1-abschluss-protokollrand` maß rund 8), (3) einen Plan, dessen Schätzung über 70 % der Marke liegt, **vor dem Anspruch** zu schneiden statt der Prüfung des Architect danach, und (4) die **Meldeschwelle** neben der Marke in §4 zu nennen (`BEO-REPO/meldeschwelle-ohne-artefakt`, einbezogen, nicht abgeschlossen).

**Entwurf des Wortlauts und Entscheidung des Nutzers.** Der Wortlaut steht in §3 als Entwurf des Planners. Regeländerungen in `AGENTS.md` und in `.claude/commands/` sind eine Entscheidung des Nutzers: Der Implementer legt den Wortlaut vor und committet beide Dateien erst nach der Bestätigung; bis dahin ist er eine Empfehlung. Entschieden hat der Nutzer am 2026-10-15 nur, dass der Slice angelegt wird, den Gegenstand (Faktor, Mindestwert 5, Schnitt ab 70 %, Meldeschwelle in §4, die beiden Dateien) und die Marken der beiden Teile von `slice-v1-abschluss-verbindungsende-format`; den Wortlaut und die Stelle in der Reihenfolge hat er nicht entschieden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein Sensor, der Marke, Meldeschwelle oder Faktor in §4 eines Plans prüft (etwa eine Erweiterung von `kopf-check`) — anderer Vorgang (Werkzeug statt Regel), mit Vertrag in `harness/vertraege/` und Gegenprobe; er wäre eine Folge der Regel und wird erst begründet, wenn das Register ein Auftreten trotz der neuen Regel zählt (`BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` bleibt dafür der Zähler). Eine Adresse gibt es nicht; ein Slice wird dann vom Planner angelegt.
- Änderung an `.claude/agents/architect.md`, `.claude/agents/planner.md`, `.claude/commands/implement-slice.md` und `.harness/skills/` — Schicht-Abgrenzung und Größe: Die Schätzung des Architect nach §3.14 gilt unverändert, der Planner hat seinen Anweisungssatz in `plan-welle.md` (`planner.md` zeigt darauf und wiederholt ihn nicht); eine dritte Datei brächte einen dritten Liefer-Punkt.
- Rückwirkendes Nachschätzen bestehender Pläne in `next/` und `open/` — Bestand bleibt bewusst stehen: Deren Marken stehen mit Metrik in §4 (§3.14), und die Prüfung des Architect vor dem Code schätzt sie ohnehin neu; nachgeschätzt wird, wenn ein Plan angefasst oder geschnitten wird. Die beiden Teile von `slice-v1-abschluss-verbindungsende-format` behalten Marke und Meldeschwelle (Entscheidung des Nutzers vom 2026-10-15).
- Die Zahlen der Schätzung des Architect (Aufschlag „rund das Doppelte“ für Tabellen- und Gegenprobe-Code in §3.14) — bleiben: Der Entwurf ergänzt die Regel des Planners, er lockert die des Architect nicht (Entfernen oder Lockern setzt den Retirement-Check voraus, §3.14). Wie die beiden Zahlen zueinander stehen, ist ein offener Punkt für den Wortlaut (§6).
- Produkt-Code, Gates, Spezifikation, Benutzerhandbuch — Schicht-Abgrenzung: Der Slice ändert zwei Dokumente des Harness und kein Verhalten des gebauten Binary.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1:** `AGENTS.md` §3.14 trägt die Regel des Planners: Schätzung beim Anlegen und beim Schnitt in der Metrik der Marke, Faktor Tests zu Produktcode mit Beleg (mindestens 5 bei einer Tabelle von Randformen mit Mutation je Element, `slice-v1-abschluss-protokollrand` maß rund 8), Schnitt vor dem Anspruch bei einer Schätzung über 70 % der Marke, Meldeschwelle neben der Marke in §4; dazu ein Paar Falsch/Richtig und der Retirement-Check mit `seit slice-harness-schaetzung-planner-faktor`. Der Wortlaut ist vom Nutzer bestätigt, **bevor** der Commit entsteht; Datum und Fundstelle der Bestätigung stehen in §7.
- [ ] **Liefer-Punkt 2:** `.claude/commands/plan-welle.md` (Schritt 6, Slices bereitstellen) nennt dieselbe Regel für den Planner und zeigt auf `AGENTS.md` §3.14, ohne sie zu wiederholen; derselbe Vorbehalt der Bestätigung.
- [ ] Beleg der Regel (kein Liefer-Punkt), Ort §7: Die Regel, auf die Pläne vom 2026-10-15 angewandt, schneidet den ungeteilten Plan von `slice-v1-abschluss-verbindungsende-format` (Schätzung des Planners 650 bis 1000, Mitte 825, 82 % der Marke 1000; die Schätzung des Architect 850 bis 1250 liegt bei 105 %) und lässt die beiden Teile stehen (Mitte 650 bei Marke 1000 und Mitte 400 bei Marke 700, 65 % und 57 %). Der Implementer rechnet es mit dem Wortlaut nach, der bestätigt ist, und nennt abweichende Ergebnisse dem Nutzer.
- [ ] Kein Wortlaut lockert oder ersetzt eine Zusage von §3.14 (Metrik, Marke, Faktor „rund das Doppelte“ des Architect); der Diff von `AGENTS.md` ergänzt (Beleg: `git diff` ohne Löschzeile in der Regel selbst, in §7).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Herkunfts-Anker `seit slice-harness-schaetzung-planner-faktor` in `AGENTS.md` §3.14 (`liegt in`).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` und `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` von `geplant` auf `verkörpert` (Zielort und Anker); `BEO-REPO/meldeschwelle-ohne-artefakt` bekommt in `state.md` den Vermerk, dass die Regel mit diesem Slice ein Artefakt hat; der Zähler bleibt 1× (die Lieferung des Artefakts ist kein Auftreten, keine Datei in `evidence/`), der Eintrag bleibt `offen` (kein Ausgang unter 3×).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3.14 | update (ergänzen) | Regel des Planners nach dem Absatz zur Metrik; ein Paar Falsch/Richtig; Begründung und Retirement-Check um die Beobachtungen und `seit slice-harness-schaetzung-planner-faktor` ergänzt |
| `.claude/commands/plan-welle.md` | update | Schritt 6 (Slices bereitstellen): ein Satz zur Schätzung beim Anlegen, zu Faktor, Schnitt ab 70 % und Meldeschwelle, mit Zeiger auf §3.14 |
| `docs/plan/planning/observations/BEO-REPO/` (`state.md` dreier Einträge) | update | Ausgang `verkörpert` für zwei Einträge, Vermerk beim dritten (Closure) |

Tests und Gates: keine neuen. Der Slice liefert Prosa, und §3.10 verlangt Mutationen für Verträge von Code und Gates, nicht für eine Regel des Planners; die Wirkung der Regel zeigt das Register (Auftreten trotz der Regel). Ein Sensor ist abgegrenzt (§1).

**Entwurf des Wortlauts für `AGENTS.md` §3.14** (Empfehlung des Planners, keine Entscheidung des Nutzers; der Nutzer ändert oder bestätigt ihn vor dem Commit). Neuer Absatz nach dem Absatz zur Metrik, vor dem Paar Falsch/Richtig:

> Der Planner schätzt beim Anlegen eines Plans und bei jedem Schnitt in derselben Metrik und nennt in §4 den Faktor Tests zu Produktcode mit seinem Beleg (der gemessene Wert eines Vorgängers mit seiner Kennung, oder „keine Messung“). Für einen Plan mit einer Tabelle von Randformen und einer Mutation je Element (§3.10) ist der Faktor mindestens 5 (`slice-v1-abschluss-protokollrand` maß rund 8). Liegt die Mitte der Schätzung über 70 % der Marke, schneidet der Planner vor dem Anspruch; die Prüfung des Architect ist dafür nicht der erste Ort. Neben der Marke nennt §4 beim Anlegen die Meldeschwelle, mit derselben Metrik; ein Plan ohne sie ist nicht angelegt.

Dazu ein Paar: *Falsch:* §4 nennt Marke 1000 und die Schätzung 650 bis 1000 „mit Verhältnis Test zu Code etwa 3 zu 1“, die Mitte liegt bei 82 % der Marke, der Architect schätzt 850 bis 1250, und der Slice geht vor dem ersten Code-Commit zurück. *Richtig:* §4 nennt den Faktor 5 mit Beleg, die Mitte der Schätzung über 70 % führt zum Schnitt beim Anlegen, jede Hälfte hat Schätzung, Marke und Meldeschwelle. Begründung: `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` und `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Retirement-Check: Ist eine der beiden Beobachtungen seit slice-harness-schaetzung-planner-faktor wieder aufgetreten?

**Entwurf für `.claude/commands/plan-welle.md`** (Schritt 6, nach dem Satz zu Randformen): „Schätze den Umfang beim Anlegen und beim Schnitt in der Metrik von `AGENTS.md` §3.14, nenne den Faktor Tests zu Produktcode mit Beleg und trage Marke und Meldeschwelle in §4 ein; liegt die Mitte der Schätzung über 70 % der Marke, schneide vor dem Anspruch (`AGENTS.md` §3.14, seit slice-harness-schaetzung-planner-faktor).“

**Reihenfolge:** `AGENTS.md` zuerst (die Quelle), dann der Command (zeigt auf sie), beide in einem Commit nach der Bestätigung; der Plan des Slice folgt in demselben Commit (`AGENTS.md` §3.9).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-format-schreiber` liegt in `done/` (Schritt 35b, Voraussetzung); dieser Slice steht in der wellenlosen Reihe hinter Schritt 42 von welle-v1-abschluss (`slice-v1-abschluss-verbindungsende-doku`) und vor Schritt 43 (umgestellt am 2026-10-15, weil er auf die Bestätigung des Wortlauts durch den Nutzer wartet; Empfehlung des Planners, nicht entschieden; Drift-Log der Roadmap vom 2026-10-15). WIP-Limit 1. Der Implementer legt den Wortlaut aus §3 (oder seine Überarbeitung) dem Nutzer vor und committet `AGENTS.md` und `plan-welle.md` erst nach der Bestätigung; bis dahin hält er an und gibt an den Aufrufer zurück.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- AGENTS.md .claude`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 120** geht der Slice zurück zur Zerlegung. **Meldeschwelle 100:** Erreicht die Summe 100, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`). Schätzung des Planners: 40 bis 80, Mitte 60, 50 % der Marke; Faktor Tests zu Produktcode entfällt (kein Code, kein Test, §3); die Schätzung ist die erste Anwendung der Regel, die der Slice liefert.
- `in-progress` → `next`: Die Regel verlangt eine Änderung an einer dritten Datei (`architect.md`, `planner.md`, `implement-slice.md`, ein Skill) oder einen Sensor — zurück zur Zerlegung (§1).
- `in-progress` → `open`: Der Nutzer bestätigt den Wortlaut nicht und hat keine Fassung, die er bestätigen würde (etwa weil die Schwelle von 70 % oder der Mindestwert 5 entfällt) — der Slice wartet auf die Entscheidung, ohne dass ein Wortlaut in `AGENTS.md` steht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Bestätigung des Wortlauts durch den Nutzer in §7 mit Datum, Closure-Notiz mit Lerneintrag geschrieben (geschärfte Regel: `AGENTS.md` §3.14, `liegt in`).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen der Regel** (`AGENTS.md` §3.12; der Slice liefert keinen Vertrag von Code, aber eine Regel, deren Ränder der Planner sonst im Einzelfall entscheidet). Entscheidungsort ist der Wortlaut in `AGENTS.md` §3.14, vor dem Commit bestätigt:

- *Schätzung als Spanne:* maßgeblich ist die Mitte der Spanne; mit der Obergrenze lägen die beiden Teile von `slice-v1-abschluss-verbindungsende-format` (750 bei Marke 1000, 500 bei Marke 700: 75 % und 71 %) selbst über 70 % und müssten weiter geschnitten werden, was die Entscheidung des Nutzers vom 2026-10-15 zu ihren Marken nicht meint. Die Wahl der Mitte ist Entwurf des Planners; der Nutzer bestätigt sie oder nennt die Obergrenze.
- *Kein gemessener Vorgänger für den Faktor:* der Plan nennt „keine Messung“ und rechnet mit dem Mindestwert; ohne Tabelle von Randformen bleibt der Faktor eine Angabe des Planners mit Beleg, kein festgelegter Wert.
- *Plan ohne Tests (Prosa, Doku):* der Faktor entfällt und steht so im Plan, nicht leer.
- *Geteilter Slice:* jede Hälfte hat eigene Schätzung, Marke und Meldeschwelle (wie bei `slice-v1-abschluss-verbindungsende-format` und `-format-schreiber`); die Schätzung der Hälften liegt zusammen nicht unter der des ungeteilten Plans.
- *Bestehende Pläne:* nicht rückwirkend (§1); gilt für neu angelegte und neu geschnittene Pläne ab der Bestätigung.
- *Schätzung des Architect höher als die des Planners:* sie gilt für den Anspruch (unverändert nach §3.14); liegt ihre Mitte über 70 % der Marke, schneidet der Planner, bevor der Implementer beginnt.
- *Wert der Meldeschwelle:* die Regel fordert, dass sie in §4 steht, mit Metrik; ihr Abstand zur Marke legt sie nicht fest (bisher 80 % bis 90 %, hier 100 von 120).
- *Aufruf-Voraussetzung, Aufrufarten des Treibers:* nicht berührt; der Slice liefert weder Werte über mehrere Aufrufe noch ein Ergebnis, das ein Hook oder Gate glaubt (`AGENTS.md` §3.12).

**Risiken:**

- Der Nutzer ändert den Wortlaut oder die Zahlen (70 %, Mindestwert 5, Mitte statt Obergrenze), oder bestätigt ihn nicht — **Ausgang:** offen bis Closure (erwartet: entfallen, der bestätigte Wortlaut steht in `AGENTS.md`; bei Ablehnung geht der Slice zurück nach `open/`, §4).
- Das Verhältnis des Faktors zum Satz „rund das Doppelte des Produktcodes“ (Architect, §3.14) bleibt im Wortlaut unklar und der Planner und der Architect rechnen verschieden — **Ausgang:** offen bis Closure (erwartet: entfallen, der bestätigte Wortlaut nennt beide Sätze und ihren Zusammenhang oder streicht keinen).
- Die Regel ist Prosa ohne Sensor: ein Plan ohne Faktor oder Meldeschwelle bleibt unbemerkt, bis das Review ihn findet (`BEO-REPO/meldeschwelle-ohne-artefakt` als Fundstelle der fehlenden Prüfbarkeit) — **Ausgang:** weiter offen: → `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (Auftreten trotz der Regel zählt weiter; ein Sensor wird erst dann ein Slice, §1).
- Die Meldeschwelle bleibt trotz Artefakt ohne Ausgang im Register, weil die Zählung unter 3× liegt — **Ausgang:** weiter offen: → `BEO-REPO/meldeschwelle-ohne-artefakt`.
- Der Slice wächst über Marke oder Meldeschwelle von §4 (zwei Dateien, Entwurf mit Zahlen, Falsch/Richtig; Schätzung 40 bis 80 gegen Marke 120) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Der Beleg der Regel (DoD, Rechnung an den Plänen vom 2026-10-15), die Bestätigung des Wortlauts durch den Nutzer (Datum, Fundstelle), die Größe in der Metrik von §4 (Basis, letzter Commit, Summe) und das Ergebnis der Schätzung dieses Plans gegen die Messung (erste Anwendung der Regel) stehen hier.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-15). Treffer in der Sub-Area `REPO`: `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (4×) und `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (5×) — beide sind der Auslöser und durch Entscheidung des Nutzers vom 2026-10-15 diesem Slice als `geplant` zugewiesen; `BEO-REPO/meldeschwelle-ohne-artefakt` (1×) ist einbezogen und bleibt offen; `BEO-REPO/groessenprobe-ohne-messbaren-stand` ist verkörpert (`AGENTS.md` §3.14, Metrik der Marke) und wird nicht berührt. Mit diesem Slice erreicht keine Beobachtung eine Schwelle neu; kein weiterer Treffer.

Nachgezählt (Planner, 2026-10-15, `AGENTS.md` §3.13): Liefer-Punkte: 2 (`AGENTS.md` §3.14, `plan-welle.md`). Schichten: zwei, das Briefing (`AGENTS.md`) und die Commands (`.claude/commands/`); Planung (Register, Plan) zählt keine Schicht; kein Produkt-Code, keine Nutzer-Doku. Der Slice trägt keine Sendung in einen Nehmer ein und nimmt keine an: keine Nachzählung bei einem anderen Slice. Schätzung in der Metrik von §4: 40 bis 80 (Mitte 60, 50 % der Marke 120); Faktor Tests zu Produktcode entfällt (kein Test); der Plan liegt unter 70 % der Marke und bleibt ungeteilt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
