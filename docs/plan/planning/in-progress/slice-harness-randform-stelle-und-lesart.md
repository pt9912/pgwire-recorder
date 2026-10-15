# Slice slice-harness-randform-stelle-und-lesart: Die Stelle einer Nachricht in der Folge ist eine Randform, und wer „Lesart“ in §7 schreibt, hält an

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er von der nächsten Welle-Closure. Angelegt am 2026-10-15 nach Entscheidung des Nutzers vom selben Tag („ja deine Vorschläge“ auf die beiden Regeländerungen aus dem Register von `BEO-REPO/spec-randform-erst-im-review-entschieden` und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`); eingeordnet in `next/`, weil er mit einer Stelle in der Reihenfolge eingeplant ist (§4 *Start*, Drift-Log der Roadmap; die Stelle ist Empfehlung des Planners).

**Bezug:** kein Anspruch des Produkts — der Slice ändert Regeln des Harness. Bindung an das Register: `BEO-REPO/spec-randform-erst-im-review-entschieden` (Auslöser dieses Slice: `slice-v1-abschluss-verbindungsende-format`, F-1; `slice-v1-abschluss-protokollrand`, V-1) und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (18×); Regelstellen: `AGENTS.md` §3.12 (seit slice-harness-randformen-vor-code) und `.claude/commands/implement-slice.md` Randform-Rückgabe (seit slice-harness-blackbox-kern). Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** Planner. **Datum:** 2026-10-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Zwei Lücken der Randform-Regel schließen, die das Register mehrfach zählt: (1) `AGENTS.md` §3.12 und die Prüfung des Architect (`.claude/agents/architect.md`, Absatz *Randformen eines neuen Vertrags*) nennen die **Stelle einer Nachricht in der Folge** als Randform: Enthält ein Vertrag eine Nachricht, die eine Folge beendet oder ersetzt (Ende, `ready_for_query`, `sync`), nennt §6 je Nachrichtenart die zulässigen Stellen (Anfang, Mitte, direkt vor dem Ende, letzte Gruppe, nicht letzte Gruppe) und je Stelle eine Entscheidung, die Spezifikation trägt sie; dazu die **Reihenfolge mehrerer Ablehnungsgründe** in einer Folge im selben Absatz. (2) `.claude/commands/implement-slice.md` (Randform-Rückgabe) sagt dem Implementer: Wer beim Schreiben von §7 „Lesart“, „Auslegung“ oder „§6 nennt nicht“ formuliert, hat eine Randform gefunden und hält an; das Vermerken in §7 ersetzt die Rückgabe nicht.

**Entwurf des Wortlauts und Entscheidung des Nutzers.** Der Nutzer hat am 2026-10-15 die beiden Regeländerungen bestätigt („ja deine Vorschläge“). Der Wortlaut in §3 ist der Entwurf des Planners nach diesen Vorschlägen; der Implementer übernimmt ihn. Ändert er ihn inhaltlich (andere Stellen in der Aufzählung, anderer Satz zu §7), legt er die Fassung dem Nutzer vor und committet erst nach dessen Antwort; bis dahin ist die Abweichung eine Empfehlung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein Sensor, der „Lesart“ oder „Auslegung“ in §7 eines Plans findet, oder einer, der die Tabelle Nachricht × Stelle in §6 prüft — anderer Vorgang (Werkzeug statt Regel), und ein Grep auf Wörter ist unscharf (Treffer in Zitaten, Fehlstellen bei Umschreibung); der Nutzer hat die Regeländerung ohne Sensor bestätigt. Eine Adresse gibt es nicht; ein Sensor wird erst begründet, wenn das Register ein Auftreten trotz der Regel zählt.
- Die Schätzung des Planners in der Metrik von `AGENTS.md` §3.14 mit Faktor, Schnitt ab 70 % der Marke und Meldeschwelle — gehört `slice-harness-schaetzung-planner-faktor`; er ändert §3.14 und `plan-welle.md`, dieser Slice §3.12, `architect.md` und `implement-slice.md` (Reihenfolge und Konfliktfreiheit: §3).
- Änderung an `.claude/commands/plan-welle.md`, `.claude/agents/planner.md`, `.claude/agents/implementer.md`, `.harness/skills/` (auch `reviewer.md`) — Schicht-Abgrenzung und Größe: Der Planner prüft §6 über den bestehenden Satz zu Randformen in `plan-welle.md` Schritt 6, den die Aufzählung in §3.12 füllt; eine weitere Datei brächte einen weiteren Liefer-Punkt, und die Zeile `plan-welle.md` gehört dem anderen Slice.
- Rückwirkendes Nachziehen der Pläne in `next/` und `open/` um die Tabelle Nachricht × Stelle — Bestand bleibt bewusst stehen: Die Prüfung des Architect vor dem Code geht die Verträge mit Nachrichtenfolge ohnehin durch (`slice-v1-abschluss-verbindungsende-record`, `-replay`, `-sqlite`, `-upstream`, `-einspielen`, `-vergleich`), und die Regel gilt ab dem Commit.
- Änderung der Spezifikation oder der Meldungen zum Verbindungsende — anderer Vorgang: Die Entscheidung der beiden Lesarten aus F-1 liegt in `slice-v1-abschluss-verbindungsende-format-schreiber` bzw. beim Architect.
- Produkt-Code, Gates, Benutzerhandbuch — Schicht-Abgrenzung: Der Slice ändert drei Dokumente des Harness und kein Verhalten des gebauten Binary.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1:** `AGENTS.md` §3.12 trägt den Absatz *Die Stelle einer Nachricht gehört dazu* (Aufzählung der Stellen, eine Entscheidung je Stelle, die Spezifikation trägt sie; Reihenfolge mehrerer Ablehnungsgründe), mit einem Paar Falsch/Richtig (Anlass `slice-v1-abschluss-verbindungsende-format`, F-1), Begründung mit `BEO-REPO/spec-randform-erst-im-review-entschieden` und Retirement-Check mit `seit slice-harness-randform-stelle-und-lesart`; `.claude/agents/architect.md` nennt die Stelle in der Prüfung der Liste aus §6 mit Zeiger auf §3.12, ohne den Absatz zu wiederholen.
- [ ] **Liefer-Punkt 2:** `.claude/commands/implement-slice.md` (Randform-Rückgabe) trägt den Satz zu „Lesart“, „Auslegung“ und „§6 nennt nicht“ in §7, mit Zeiger auf `AGENTS.md` §3.12 und dem Herkunfts-Anker `seit slice-harness-randform-stelle-und-lesart`; kein Sensor.
- [ ] Kein Wortlaut lockert oder ersetzt eine Zusage von §3.12 oder der Randform-Rückgabe (Beleg: `git diff` der drei Dateien ohne Löschzeile in den bestehenden Regeln, in §7).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Herkunfts-Anker `seit slice-harness-randform-stelle-und-lesart` in `AGENTS.md` §3.12 und in `implement-slice.md` (`liegt in`).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: `BEO-REPO/spec-randform-erst-im-review-entschieden` und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` von `geplant` auf `verkörpert` (Zielort und Anker); der Zähler folgt den Dateien in `evidence/`, kein Auftreten durch die Lieferung.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3.12 | update (ergänzen) | neuer Absatz *Die Stelle einer Nachricht gehört dazu* nach dem Absatz *Aufrufarten des Treibers*, vor der Begründung; Paar Falsch/Richtig; Begründung und Retirement-Check um den Anlass und `seit slice-harness-randform-stelle-und-lesart` ergänzt |
| `.claude/agents/architect.md` | update | im Absatz *Randformen eines neuen Vertrags*: ein Satz, dass die Prüfung der Liste auch die Stelle jeder Nachricht, die eine Folge beendet oder ersetzt, und die Reihenfolge mehrerer Ablehnungsgründe umfasst (`AGENTS.md` §3.12) |
| `.claude/commands/implement-slice.md` | update | Randform-Rückgabe: ein Satz nach *Auch eine Lesart ist eine Entscheidung* |

**Entwurf des Wortlauts für `AGENTS.md` §3.12** (Planner, nach den bestätigten Vorschlägen). Neuer Absatz nach *Aufrufarten des Treibers*:

> **Die Stelle einer Nachricht gehört dazu** (seit slice-harness-randform-stelle-und-lesart, Entscheidung des Nutzers vom 2026-10-15): Enthält ein Vertrag eine Nachricht, die eine Folge beendet oder ersetzt (Ende, `ready_for_query`, `sync`), nennt §6 je Nachrichtenart die zulässigen Stellen in der Folge (Anfang, Mitte, direkt vor dem Ende, letzte Gruppe, nicht letzte Gruppe) und je Stelle eine Entscheidung; die Spezifikation trägt sie. Kann eine Folge aus mehreren Gründen abgelehnt werden, nennt §6 die Reihenfolge der Gründe und die Spezifikation bindet sie.
>
> **Falsch:** §6 nennt die Nachricht `ready_for_query` und `connection_closed`, aber nicht, ob sie vor dem Ende oder mitten in der Folge stehen dürfen; der Implementer entscheidet beide im Code, und das Review findet sie (F-1 in `slice-v1-abschluss-verbindungsende-format`).
> **Richtig:** §6 führt je Nachricht die Stellen auf und zeigt auf die Stelle der Spezifikation, die jede entscheidet; dann folgt der Code.

Begründung (Ergänzung der bestehenden): `BEO-REPO/spec-randform-erst-im-review-entschieden`; Retirement-Check um „für die Stelle einer Nachricht seit slice-harness-randform-stelle-und-lesart“ ergänzt.

**Entwurf für `.claude/commands/implement-slice.md`** (Randform-Rückgabe, nach *Auch eine Lesart ist eine Entscheidung*): „Schreibst du in §7 „Lesart“, „Auslegung“ oder „§6 nennt nicht“, hast du eine Randform gefunden und hältst an; ein Vermerk in §7 ersetzt die Rückgabe nicht (`AGENTS.md` §3.12, seit slice-harness-randform-stelle-und-lesart, `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`).“

**Entwurf für `.claude/agents/architect.md`:** an den Absatz *Randformen eines neuen Vertrags* anschließend: „Zur Liste gehört die Stelle jeder Nachricht, die eine Folge beendet oder ersetzt, und die Reihenfolge mehrerer Ablehnungsgründe (`AGENTS.md` §3.12, seit slice-harness-randform-stelle-und-lesart); fehlt sie, meldest du sie dem Planner.“

**Reihenfolge und Konfliktfreiheit zu `slice-harness-schaetzung-planner-faktor`:** Beide Slices ändern `AGENTS.md` und Dateien unter `.claude/`, aber verschiedene Stellen: dieser §3.12 (`architect.md`, `implement-slice.md`), jener §3.14 und `plan-welle.md`; keine Datei außer `AGENTS.md` ist gemeinsam, und dort liegen die Absätze getrennt (die Textmarken von §3.12 und §3.14 überlappen nicht). Dieser Slice steht in der Reihenfolge vor jenem (Schritt 35b liegt dazwischen, §4); WIP-Limit 1 verhindert parallele Zweige, und die Basis jeder Marke ist der Stand des Hauptzweigs vor dem eigenen Anspruch (hier `88f9c5e`), die Diffs zählen sich also nicht gegenseitig. Jener Slice zählt `architect.md` und `implement-slice.md` in seiner Abgrenzung als „nicht seines“; dieser Slice nimmt sie, ohne dass jener Gegenstand verliert (`grep` nach der Kennung dieses Slice in jenem bleibt leer, weil keine Sendung hin- oder hergeht).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Dieser Slice steht in der wellenlosen Reihe direkt **vor** `slice-v1-abschluss-verbindungsende-format-schreiber` (Schritt 35b) und damit vor `slice-harness-schaetzung-planner-faktor` (Empfehlung des Planners, nicht entschieden; Drift-Log der Roadmap vom 2026-10-15): Der Anlass F-1 betrifft genau die Verträge zur Nachrichtenfolge, die 35b und die sechs folgenden Slices zum Verbindungsende liefern, und der Architect prüft deren §6 mit der Regel. Der Nutzer kann ihn hinter welle-v1-abschluss legen. WIP-Limit 1. Kein weiterer Slice in `in-progress/`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat 88f9c5e..<letzter Commit des Slice> -- AGENTS.md .claude`, wobei die Basis `88f9c5e` der Stand des Hauptzweigs vor dem Anspruch ist; **über 100** geht der Slice zurück zur Zerlegung. **Meldeschwelle 80:** Erreicht die Summe 80, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet. Schätzung des Planners: 30 bis 50, Mitte 40, 40 % der Marke (anwendbar: Absatz mit Falsch/Richtig rund 15 bis 20 Zeilen, zwei Sätze in den Commands je 2 bis 4 Zeilen, Retirement-Check und Begründung 3 bis 6 Zeilen); Faktor Tests zu Produktcode entfällt (kein Code, kein Test).
- `in-progress` → `next`: Die Regel verlangt eine vierte Datei (`plan-welle.md`, `planner.md`, ein Skill) oder einen Sensor — zurück zur Zerlegung (§1).
- `in-progress` → `open`: Der Nutzer lehnt eine inhaltliche Abweichung des Implementers vom Wortlaut ab und hat keine Fassung, die er bestätigen würde — der Slice wartet auf die Entscheidung, ohne dass ein abweichender Wortlaut in `AGENTS.md` steht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben (geschärfte Regel: `AGENTS.md` §3.12 und `implement-slice.md`, `liegt in`).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen der Regel** (`AGENTS.md` §3.12; der Slice liefert keinen Vertrag von Code, aber eine Regel, deren Ränder sonst im Einzelfall entschieden werden). Entscheidungsort ist der Wortlaut in `AGENTS.md` §3.12:

- *Reichweite der Aufzählung:* „Anfang, Mitte, direkt vor dem Ende, letzte Gruppe, nicht letzte Gruppe“ ist die Liste der bestätigten Vorschläge; Stellen, die ein Vertrag nicht kennt (etwa „letzte Gruppe“ in einer Folge ohne Gruppen), entfallen mit dem Satz „nicht berührt“ in §6, nicht leer.
- *Reihenfolge der Ablehnungsgründe:* gilt, wo eine Folge aus mehr als einem Grund abgelehnt werden kann; bei genau einem Grund entfällt sie. Der Wortlaut nennt „wenn eine Folge aus mehreren Gründen abgelehnt werden kann“ (Anlass `slice-v1-abschluss-protokollrand`, V-1, `_pq_.` vor `PGR-E6003`).
- *„Lesart“ im Zitat:* Der Satz zu §7 gilt für den Text, den der Implementer als eigene Feststellung schreibt; ein Zitat der Spezifikation oder eines Reviews mit dem Wort ist keine Randform. Ein Sensor darauf ist ausgeschlossen (§1).
- *Aufruf-Voraussetzung, Aufrufarten des Treibers:* nicht berührt; der Slice liefert weder Werte über mehrere Aufrufe noch ein Ergebnis, das ein Hook oder Gate glaubt (`AGENTS.md` §3.12).

**Risiken:**

- Die Aufzählung der Stellen ist zu eng oder zu weit (ein Vertrag mit einer Stelle, die sie nicht nennt, bleibt ohne Eintrag in §6) — **Ausgang:** weiter offen: → `BEO-REPO/spec-randform-erst-im-review-entschieden` (ein Auftreten trotz der Regel zählt dort weiter).
- Der Satz in `implement-slice.md` wird gelesen, aber nicht befolgt, weil kein Sensor ihn prüft; das Review findet die Randform weiter erst nach dem Code — **Ausgang:** weiter offen: → `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (Auftreten trotz der Regel zählt weiter; ein Sensor wird erst dann ein Slice, §1).
- Ein Auftreten nach diesem Slice schwächt die Regel oder bestätigt sie — das ist Urteil des Planners beim Retirement-Check (`AGENTS.md` §3.12, §3.14-Muster), nicht des Implementers — **Ausgang:** entfallen, wenn bis zur Closure kein Auftreten im Register steht; sonst weiter offen: → die jeweilige Beobachtung.
- Der Slice wächst über Marke oder Meldeschwelle von §4 (drei Dateien, Entwurf mit Falsch/Richtig; Schätzung 30 bis 50 gegen Marke 100) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Der Beleg ohne Löschzeile (DoD), der Lerneintrag (geschärfte Regel: `AGENTS.md` §3.12 und `implement-slice.md`), das Register und die Risiko-Ausgänge stehen dann hier.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-15). Treffer in der Sub-Area `REPO`: `BEO-REPO/spec-randform-erst-im-review-entschieden` und `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` sind der Auslöser und durch Entscheidung des Nutzers vom 2026-10-15 diesem Slice als `geplant` zugewiesen (beide stehen über 3×; mit diesem Slice erreicht keine Beobachtung eine Schwelle neu). `BEO-REPO/plan-folgt-korrektur-nicht` und `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` werden nicht berührt: Der Slice ändert drei Regeldokumente und trägt die Änderungen in §1, §3 und §6 nach (`AGENTS.md` §3.9). Kein weiterer Treffer.

Nachgezählt (Planner, 2026-10-15, `AGENTS.md` §3.13): Liefer-Punkte: 2 (`AGENTS.md` §3.12 mit dem Satz in `architect.md`; `implement-slice.md`). Schichten: zwei, das Briefing (`AGENTS.md`) und die Rollen- und Command-Artefakte (`.claude/agents/`, `.claude/commands/`); Planung (Register, Plan, Roadmap) zählt keine Schicht; kein Produkt-Code, keine Nutzer-Doku. Der Slice nimmt die Zuweisung beider Beobachtungen aus dem Register an (Adresse nimmt an: §1 schließt sie nicht aus, die DoD führt sie) und trägt keine Sendung in einen anderen Nehmer ein: keine Nachzählung bei einem anderen Slice. Schätzung in der Metrik von §4: 30 bis 50 (Mitte 40, 40 % der Marke 100); Faktor Tests zu Produktcode entfällt (kein Test); der Plan liegt unter 70 % der Marke und bleibt ungeteilt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
