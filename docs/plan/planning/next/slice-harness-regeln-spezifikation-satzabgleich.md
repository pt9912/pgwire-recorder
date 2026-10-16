# Slice slice-harness-regeln-spezifikation-satzabgleich: Der Architect gleicht vor dem Code die Sätze einer Spezifikationsstelle ab (Einleitung, Aufzählung, Tabelle)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er von der nächsten Welle-Closure. Angelegt am 2026-10-16 als Schnitt von `slice-harness-regeln-nach-v1-abschluss` (Planner, beim Eintragen der Sendung, `AGENTS.md` §3.13) nach dem Beschluss des Nutzers vom selben Tag („ja deine Empfehlungen“ zu den Fragen der Closure von `slice-v1-abschluss-verbindungsende-einspielen-extended`: der Architect nimmt den Abgleich der Sätze einer Spezifikationsstelle in seine Liste vor dem Code auf); eingeordnet in `next/` hinter `slice-harness-regeln-nach-v1-abschluss` und vor `slice-harness-schaetzung-planner-faktor`; seit dem 2026-10-16 hinter dem Release (§4 *Start*). Den **Wortlaut** hat der Nutzer nicht entschieden (§1).

**Bezug:** kein Anspruch des Produkts — der Slice ändert Regeln des Harness. Bindung an das Register: `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` (3×); Regelstellen: `AGENTS.md` §3.12 (seit slice-harness-randformen-vor-code) und `.claude/agents/architect.md`. Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-16.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Architect gleicht vor dem Code die Sätze jeder Spezifikationsstelle, die der Slice berührt oder neu schreibt (Einleitungssatz, Aufzählung, Tabelle), gegeneinander ab, und das steht in seiner Liste vor dem Code: in `.claude/agents/architect.md` und in `AGENTS.md` §3.12 (Randformen vor dem Code; „Der Architect prüft die Liste vor dem Code“).

**Entwurf des Wortlauts und Entscheidung des Nutzers.** Regeländerungen in `AGENTS.md` und in `.claude/agents/` sind eine Entscheidung des Nutzers. Der Nutzer hat am 2026-10-16 („ja deine Empfehlungen“) bestätigt, dass der Abgleich in die Liste des Architect aufgenommen wird; **den Wortlaut hat er nicht entschieden**. Der Implementer legt den Wortlaut aus §3 (oder seine Überarbeitung) dem Nutzer vor und committet beide Dateien erst nach der Bestätigung, **Nutzerentscheidung über den Wortlaut vor der Umsetzung**; bis dahin ist er eine Empfehlung. Was der Nutzer nicht bestätigt, nimmt der Implementer nicht auf und nennt es in §7.

**Beleg-Stand** (Stand 2026-10-16, gemergter Stand des Registers; die Belege stehen dort und in den Closure-Notizen der Slices, nicht hier): `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` steht bei 3×: `slice-v1-abschluss-protokollrand` (F-2: die Zeile zu abgelehnten Verbindungen gegen die *Zählung der Verbindungen* der Spezifikation), `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-verbindungsende-einspielen-extended` (F-9 und V-3: Einleitungssatz gegen Aufzählungspunkt in den Stellen *Gruppen* und *Ende der Aufzeichnung* der Spezifikation). Der Fund kam jedes Mal erst vom Review oder von der Verifikation. Offen und nicht mitgezählt: `slice-v1-abschluss-verbindungsende-replay-extended` (F-1, als Randform geführt; die Frage, ob er mitzählt, steht beim Nutzer).

**Geschnitten aus `slice-harness-regeln-nach-v1-abschluss`** (Planner, 2026-10-16, `AGENTS.md` §3.13): Der Regel-Slice zählt drei Liefer-Punkte in zwei Schichten (Briefing, Commands); mit dem Abgleich wären es vier Liefer-Punkte in drei Schichten (Briefing, Commands, Agenten-Rollen), über beiden Grenzen. Einzeln lieferbar: Dieser Slice setzt `slice-harness-regeln-nach-v1-abschluss` nicht voraus und er nicht ihn; beide ändern `AGENTS.md` in verschiedenen Abschnitten (§6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Regeln aus dem Verbindungsende-Block (erzeugte Abdeckungstabellen als Schicht, Wortlaut-Suche, Arbeitsbaum vor E2E-Mutanten, vier offene Fragen zu `AGENTS.md` §3.10) — `slice-harness-regeln-nach-v1-abschluss` (Klasse 1): andere Abschnitte, anderer Nutzer-Entscheid, andere Rolle.
- Der Abgleich zwischen Lastenheft und Spezifikation (Stratum gegen Stratum) — `BEO-REPO/spezifikation-weicht-vom-lastenheft-ab` (1×, Frage an den Nutzer offen): Die Angleichung nach unten ist keine Entscheidung des Architect, und die Beobachtung hat noch keinen Ausgang (Klasse 3, anderer Vorgang als der Abgleich innerhalb einer Stelle).
- Ein Sensor, der Widersprüche in der Spezifikation findet — anderer Vorgang (Werkzeug statt Regel): ein Sensor bekäme einen Vertrag in `harness/vertraege/` und eine Gegenprobe; ob der Wortlaut ein Widerspruch ist, ist Urteil.
- Rückwirkendes Bereinigen von Stellen der Spezifikation — Bestand bleibt bewusst stehen: Der Abgleich gilt für Stellen, die ein Slice berührt oder neu schreibt; niemand hat den Rest entschieden.
- Änderung an `.claude/agents/planner.md`, `reviewer.md`, `.harness/skills/`, `plan-welle.md`, `close-welle.md` und an den Hooks — Schicht-Abgrenzung und Größe: Die Regel gehört in die Liste des Architect und in das Briefing; eine dritte Datei brächte einen dritten Liefer-Punkt.
- Produkt-Code, Gates, Spezifikation, Benutzerhandbuch — Schicht-Abgrenzung: Der Slice ändert zwei Dokumente des Harness und kein Verhalten des gebauten Binary; **kein Produktcode**, die Spezifikation ändert der Architect im jeweiligen Slice, nicht hier.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt: Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan **geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1:** `AGENTS.md` §3.12 trägt den Satz: Der Architect gleicht vor dem Code die Sätze einer Stelle der Spezifikation, die der Slice berührt oder neu schreibt (Einleitung, Aufzählung, Tabelle), gegeneinander ab; ein Einleitungssatz, den ein Aufzählungspunkt oder eine Tabellenzeile einschränkt, nennt die Ausnahme oder sagt nicht allgemein zu. Dazu ein Paar Falsch/Richtig (Beleg: `slice-v1-abschluss-verbindungsende-einspielen-extended`, F-9, V-3), die Begründung mit `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` und der Retirement-Check `seit slice-harness-regeln-spezifikation-satzabgleich`. Der Wortlaut ist vom Nutzer bestätigt, **bevor** der Commit entsteht; Datum und Fundstelle stehen in §7.
- [ ] **Liefer-Punkt 2:** `.claude/agents/architect.md` nennt den Abgleich in dem Absatz zur Liste vor dem Code (neben der Stelle einer Nachricht und der Reihenfolge der Ablehnungsgründe), mit Zeiger auf `AGENTS.md` §3.12 und Herkunfts-Anker. Derselbe Vorbehalt der Bestätigung.
- [ ] Kein Wortlaut lockert oder ersetzt eine Zusage von §3.12; der Diff von `AGENTS.md` ergänzt (Beleg: `git diff` ohne Löschzeile in den Regeln selbst, in §7).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; Herkunfts-Anker `seit slice-harness-regeln-spezifikation-satzabgleich` in `AGENTS.md` §3.12 (`liegt in`), soweit aufgenommen.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` von `geplant` auf `verkörpert` (Zielort und Anker); bei Ablehnung des Wortlauts der Zähler unverändert und der Beleg-Stand genannt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` §3.12 | update (ergänzen) | Satz zum Abgleich nach dem Absatz zur Stelle einer Nachricht; Paar Falsch/Richtig; Begründung und Retirement-Check |
| `.claude/agents/architect.md` | update | Der Abgleich in dem Absatz zur Liste vor dem Code |
| `docs/plan/planning/observations/BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut/` (`state.md`) | update | Ausgang bei Closure |

Tests und Gates: keine neuen. Der Slice liefert Prosa (§3.10 verlangt Mutationen für Verträge von Code und Gates); die Wirkung zeigt das Register (Auftreten trotz der Regel).

**Entwurf des Wortlauts** (Empfehlung des Planners, keine Entscheidung des Nutzers; Nutzerentscheidung über den Wortlaut vor der Umsetzung):

> *`AGENTS.md` §3.12, nach der Stelle einer Nachricht:* **Die Sätze einer Stelle gehören dazu**: Schreibt oder berührt ein Slice eine Stelle der Spezifikation, gleicht der Architect vor dem Code ihre Sätze gegeneinander ab (Einleitungssatz, Aufzählung, Tabelle). Ein Einleitungssatz, den ein Aufzählungspunkt oder eine Tabellenzeile einschränkt, nennt die Ausnahme oder sagt nicht allgemein zu.
>
> *Falsch:* Der Einleitungssatz von `LH-FA-24.a` verspricht den Vergleich „bis zum `ReadyForQuery`“, der Aufzählungspunkt nimmt die letzte `Flush`-Gruppe aus; ein Leser, der nur den ersten liest, erwartet anderes, und die Verifikation findet es.
> *Richtig:* Der Einleitungssatz nennt die Ausnahme oder zeigt auf den Punkt, der sie nennt; der Architect hat beide Sätze vor dem Code nebeneinander gelesen.
>
> *`architect.md`, Ergänzung im Absatz zur Liste vor dem Code:* Zur Liste gehört der Abgleich der Sätze jeder Stelle der Spezifikation, die der Slice berührt oder neu schreibt (`AGENTS.md` §3.12, seit slice-harness-regeln-spezifikation-satzabgleich).

**Reihenfolge:** `AGENTS.md` zuerst (die Quelle), dann die Rolle, beide in einem Commit nach der Bestätigung; der Plan des Slice folgt in demselben Commit (`AGENTS.md` §3.9).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [welle-v1-nachtrag](../welle-v1-nachtrag.md) liegt in `done/` (Voraussetzung der Stelle, nicht der Technik); der Slice steht in der wellenlosen Reihe hinter `slice-harness-regeln-nach-v1-abschluss` und vor `slice-harness-schaetzung-planner-faktor` (Empfehlung des Planners) und ist seit dem 2026-10-16 hinter das Release verschoben (Entscheidung des Nutzers vom 2026-10-16, Drift-Log der Roadmap): wellenlos hinter [welle-v1-nachtrag](../welle-v1-nachtrag.md) und vor der übrigen Harness-Reihe (`slice-harness-commit-struktur-id` zuerst). Bis dahin prüft der Architect den Abgleich der Sätze einer Stelle der Spezifikation nicht aus dieser Regel (`BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut`, 3×, Ausgang `geplant` bleibt). WIP-Limit 1. Der Implementer legt den Wortlaut dem Nutzer vor und committet `AGENTS.md` und `architect.md` erst nach der Bestätigung; bis dahin hält er an und gibt an den Aufrufer zurück.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): **Marke mit Metrik** (`AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- AGENTS.md .claude`, wobei die Basis der Commit des Übergangs nach `in-progress/` ist; **über 100** geht der Slice zurück zur Zerlegung. **Meldeschwelle 80:** Erreicht die Summe 80, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet (`BEO-REPO/meldeschwelle-ohne-artefakt`); die Entscheidung des Planners steht als datierte Zeile in §7. Schätzung des Planners vor dem Anspruch: **20 bis 45, Mitte 32, 32 % der Marke**; Faktor Tests zu Produktcode entfällt (kein Code, kein Test, §3); Basis der Schätzung: ein Satz mit einem Paar Falsch/Richtig, Begründung und Retirement-Check in `AGENTS.md` (rund 15 bis 25 Zeilen) und zwei Zeilen in der Rolle.
- `in-progress` → `next`: Die Regel verlangt eine Änderung an einer dritten Datei (`planner.md`, `reviewer.md`, ein Skill, ein Command) oder einen Sensor — zurück zur Zerlegung (§1).
- `in-progress` → `open`: Der Nutzer bestätigt keinen Wortlaut und hat keine Fassung, die er bestätigen würde — der Slice wartet auf die Entscheidung, ohne dass ein Wortlaut in `AGENTS.md` steht.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Bestätigung des Wortlauts durch den Nutzer in §7 mit Datum, Closure-Notiz mit Lerneintrag geschrieben (geschärfte Regel: `AGENTS.md` §3.12, `liegt in`, soweit aufgenommen).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen der Regel** (`AGENTS.md` §3.12; der Slice liefert keinen Vertrag von Code, aber eine Regel, deren Ränder der Planner sonst im Einzelfall entscheidet). Entscheidungsort ist der Wortlaut, vor dem Commit bestätigt:

- *Reichweite:* Stellen, die der Slice berührt oder neu schreibt; nicht der Bestand der Spezifikation (§1).
- *Gegenstand des Abgleichs:* Sätze derselben Stelle (Einleitungssatz, Aufzählungspunkte, Tabellenzeilen, Randformen-Absätze). Lastenheft gegen Spezifikation ist ein anderer Vorgang (§1, `BEO-REPO/spezifikation-weicht-vom-lastenheft-ab`).
- *Wer entscheidet den Widerspruch:* Der Architect gleicht an, in der Spezifikation, an dem Ort, den §3.12 nennt; was das Lastenheft betrifft, legt er dem Nutzer vor.
- *Abgrenzung zur Randform:* Ist eine Randform unentschieden, bleibt es die Liste von §3.12 (`BEO-REPO/spec-randform-erst-im-review-entschieden`); der Abgleich gilt entschiedenen Sätzen, die sich widersprechen.
- *Aufruf-Voraussetzung, Aufrufarten des Treibers, Stelle einer Nachricht:* nicht berührt.

**Risiken:**

- Der Nutzer ändert den Wortlaut oder bestätigt keinen — **Ausgang:** offen bis Closure (erwartet: entfallen, der bestätigte Wortlaut steht in `AGENTS.md`; bei Ablehnung geht der Slice zurück nach `open/`, §4).
- Die Regel ist Prosa ohne Sensor: ein Widerspruch zwischen Sätzen einer Stelle bleibt unbemerkt, bis das Review ihn findet — **Ausgang:** weiter offen → `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` (Auftreten trotz der Regel zählt weiter; ein Sensor ist abgegrenzt, §1).
- Die Änderungen in `AGENTS.md` treffen auf `slice-harness-regeln-nach-v1-abschluss` (§3.13, §3.10, derselbe Rang) und `slice-harness-schaetzung-planner-faktor` (§3.14) — **Ausgang:** offen bis Closure (erwartet: entfallen, die Abschnitte überlappen nicht; der Slice, der später schließt, rechnet den Wortlaut des ersten ein).
- Der Slice wächst über Marke oder Meldeschwelle von §4 (Schätzung 20 bis 45 gegen Marke 100) — **Ausgang:** offen bis Closure; Marke und Meldeschwelle stehen in §4.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`). Die Bestätigung des Wortlauts durch den Nutzer (Datum, Fundstelle), die Größe in der Metrik von §4 (Basis, letzter Commit, Summe) und das Ergebnis der Schätzung gegen die Messung stehen hier.

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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-16). Treffer in der Sub-Area `REPO`: `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` (3×, Stand `geplant` mit diesem Slice) — sie ist der Auslöser (§1, Beleg-Stand); `BEO-REPO/spec-randform-erst-im-review-entschieden` und `BEO-REPO/spezifikation-weicht-vom-lastenheft-ab` (1×) sind abgegrenzt (§1, §6). Eine Schwelle wird mit diesem Slice nicht neu erreicht (keine Evidence-Datei entsteht).

Nachgezählt (Planner, 2026-10-16, `AGENTS.md` §3.13): Liefer-Punkte: 2 (§3.12-Satz, Rolle `architect.md`). Schichten: zwei, das Briefing (`AGENTS.md`) und die Agenten-Rollen (`.claude/agents/`); Planung (Register, Plan) zählt keine Schicht; kein Produkt-Code, keine Nutzer-Doku. Der Slice ist die Adresse der Sendung aus `slice-harness-regeln-nach-v1-abschluss` (Schnitt, siehe dort §8); er nimmt keinem anderen Nehmer etwas ab. Schätzung in der Metrik von §4: 20 bis 45 (Mitte 32, 32 % der Marke 100); der Plan liegt unter 70 % der Marke und bleibt ungeteilt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

Nachgezählt beim Verschieben hinter das Release (Planner, 2026-10-16, `AGENTS.md` §3.13): Start und Stelle geändert, Gegenstand und DoD unverändert; der Slice nimmt die Sendung aus `slice-harness-regeln-nach-v1-abschluss` weiter an (nicht in `done/`, §1, DoD und §3 nennen sie): Liefer-Punkte weiter 2, Schichten weiter zwei. Kein Slice vor dem Release nennt ihn als Voraussetzung.
