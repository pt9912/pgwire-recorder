# Slice slice-lastenheft-pruefbarkeit: Qualitätsanforderung an die Prüfbarkeit des Quellcodes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Der Slice legt
das neue Abnahmeszenario 17 an, das zu M3 gehört; nachweisbar wird es erst mit
`slice-harness-coverage`, dem letzten Slice der Reihe. Eingesammelt wird er von der
nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-05: nach `slice-replay-semantik-fehlerreplay` und vor dem nächsten großen
Slice (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** — (dieser Slice schreibt die Anforderung; ihre Kennung entsteht mit ihm). Entscheidungen des Nutzers vom 2026-10-05 und 2026-10-06; Abnahmeszenario vom Team entschieden (Koordinator, 2026-10-06).

**Berührte Spec-Stellen:** `lastenheft.md §4` · `lastenheft.md §7`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** §4 des Lastenhefts führt eine neue Qualitätsanforderung an die
Prüfbarkeit des Quellcodes mit Priorität MUSS, im Wortlaut aus §6, den der Nutzer am
2026-10-06 unverändert bestätigt hat; §7 führt dazu das Abnahmeszenario 17 (alle drei
Prüfungen grün, jede nachweislich rot bei verletzter Bedingung — Messmethode 4), und
M3 umfasst es. Der Weg, auf dem die Anforderung in den Abdeckungstabellen belegt
erscheint, ist vor dem Code entschieden (§6 *Nachweis*); die sechs Slices, die sie
umsetzen, führen sie im Kopf unter `Bezug`.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-05: Das Lastenheft bekommt eine
Anforderung zur Wartbarkeit bzw. Prüfbarkeit des Quellcodes, damit Lint, Black-Box-Tests
und Testabdeckung einen Vertragsbezug haben. Der Planner legt einen Wortlaut vor; er
schreibt ihn nicht ins Lastenheft. Am 2026-10-06 hat der Nutzer den Wortlaut bestätigt
und MUSS gewählt; das Abnahmeszenario hat das Team entschieden, ebenso, dass die
Anforderung in den Abdeckungstabellen nicht als unbelegt erscheinen darf.

**Warum ein eigener Slice und nicht der erste Schritt von `slice-harness-lint`:**
(1) andere Schicht — der Lint-Slice ändert Harness und `Dockerfile` und grenzt
Lastenheft und Spezifikation ausdrücklich aus; eine Vertragsänderung darin hebt diese
Selbstbindung auf; (2) Größe — der Lint-Slice hat schon drei Liefer-Punkte; (3)
Reihenfolge — die Kennung muss im Lastenheft stehen, bevor ein Kopf sie verlinken
kann (das Doku-Gate verlangt für jede Lastenheft-Kennung einen auflösenden Link), und
der Wortlaut braucht die Bestätigung des Nutzers vor dem Start des Lint-Slice; (4) ein
Wortlaut, den der Nutzer verwirft, lässt diesen Slice entfallen, ohne den Lint-Slice
zu berühren.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Werkzeuge, Regelprofil, Schwellen und Zahlen — das Lastenheft legt keine
  technischen Entscheidungen fest (Lastenheft §1); sie stehen in den ADRs von
  `slice-harness-lint` und `slice-harness-coverage`. Die Anforderung sagt „festgelegt“,
  nicht wo.
- Die Gates selbst und ihre Nachweise — übernehmen `slice-harness-lint` (Messmethode
  1), `slice-harness-coverage` (Messmethode 2) und die vier Umstellungs-Slices
  (Messmethode 3, vollständig mit `slice-harness-blackbox-einstieg`); jeder liefert den
  Nachweis seiner Messmethode auf dem Weg, den dieser Slice festlegt.
- Die Erweiterung von `make abdeckung` selbst, falls der Architect diesen Weg wählt —
  ein anderer Vorgang am Werkzeug; dieser Slice legt den Weg und, bei einer
  Erweiterung, die ADR fest (§4 Rückführung, §6 *Nachweis*).
- Eine Stelle in Spezifikation oder Sicht — das Lastenheft nennt die Anforderung,
  ihre technische Schärfung tragen die ADRs; eine Spezifikationsstelle entsteht nur,
  wenn §6 *Spezifikation* das anders entscheidet.
- Produkt-Code und Harness-Werkzeuge — Schicht-Abgrenzung: Der Slice ändert das
  Lastenheft, die Roadmap und die Pläne von sechs Slices; eine ADR nur bei §6
  *Nachweis*.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Lastenheft: §4 führt die Anforderung im bestätigten Wortlaut aus §6 unter der
      nächsten freien Kennung der Qualitätsanforderungen, mit Priorität MUSS (§6
      *Priorität*); §7 führt das Abnahmeszenario 17 mit Bezug auf die Anforderung;
      Version und Historie nach §6 *Version*.
- [ ] Roadmap: Der Trigger von M3 nennt das Abnahmeszenario 17 (vom Planner am
      2026-10-06 vorgezogen, mit Drift-Log-Zeile; bei Closure gegen §7 des
      Lastenhefts geprüft). Der Kopf (`Bezug`) von `slice-harness-lint`,
      `slice-harness-coverage` und der vier Umstellungs-Slices verlinkt die
      Anforderung, die Platzhalter-Sätze dort sind entfernt; `make kopf-check` und
      `make docs-check` sind grün.
- [ ] Nachweis-Weg entschieden (Architect, vor dem Lastenheft-Commit, `AGENTS.md`
      §3.12): wie jede der Messmethoden 1 bis 3 in den Abdeckungstabellen belegt
      erscheint und wann die Anforderung als vollständig zählt (§6 *Nachweis*). Wählt
      er eine Erweiterung von `make abdeckung`, liegt deren ADR vor, die
      [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) ergänzt; die
      Folge-Slices nennen in ihrer DoD, welchen Nachweis sie liefern.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | neue Anforderung in §4 mit Priorität MUSS; Abnahmeszenario 17 in §7 |
| `docs/plan/planning/open/slice-harness-lint.md`, `slice-harness-coverage.md`, `slice-harness-blackbox-*.md` | update | Kopf `Bezug` mit Link auf die Anforderung; DoD-Zeile zum Nachweis nach dem entschiedenen Weg |
| `docs/plan/planning/in-progress/roadmap.md` | geprüft | Trigger von M3 mit Abnahmeszenario 17 steht seit der Planung (2026-10-06); bei Closure gegen §7 abgeglichen |
| `docs/plan/adr/<NNNN>-…md`, `docs/plan/adr/README.md` | neu / update, nur bei Werkzeug-Erweiterung | ADR, die [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) um Gate-Nachweise ergänzt (Architect) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-replay-semantik-fehlerreplay` liegt in
`done/` (der Wortlaut ist seit 2026-10-06 bestätigt). Erster Slice der Reihe
(Lastenheft, Lint, die vier Umstellungs-Slices, Coverage). Vor dem Lastenheft-Commit
entscheidet der Architect §6 *Nachweis*.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Architect wählt für den
  Nachweis eine Erweiterung von `make abdeckung`, und sie soll vor `slice-harness-lint`
  gebaut sein; dann trägt ein eigener Harness-Slice Erweiterung und Gegenprobe, und
  dieser Slice liefert nur Lastenheft, Roadmap und ADR. Ebenso, wenn zusätzlich eine
  Spezifikationsstelle entschieden wird.
- `in-progress` → `open` (blockiert — Carveout?): Der Architect findet keinen
  Nachweis-Weg, der ohne Änderung der Lesart von
  [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (eine
  Qualitätsanforderung hat genau einen Pfad *Messung*) eine Teilmessung von der
  vollständigen unterscheidet, und eine ersetzende ADR ist nötig.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Vorschlag des Planners für den Wortlaut** — zur Entscheidung durch den Nutzer,
nicht geschrieben (Überschrift unter der nächsten freien Nummer in §4, heute die
siebte Qualitätsanforderung):

> **— Prüfbarkeit des Quellcodes**
>
> - **Anforderung:** Der Quellcode muss so gehalten werden, dass jede Änderung
>   automatisch auf Verstöße gegen festgelegte Struktur- und Komplexitätsregeln und auf
>   nicht durch Tests ausgeführten Code geprüft werden kann. Tests müssen das Verhalten
>   einer Einheit über deren öffentliche Schnittstelle prüfen.
> - **Messmethode:** (1) Eine statische Analyse nach einem festgelegten Regelprofil
>   meldet für den Quellstand keinen Befund; Ausnahmen stehen zentral und begründet,
>   nicht im Code. (2) Die Anweisungsabdeckung der Unit-Tests liegt nicht unter einer
>   festgelegten Schwelle. (3) Jeder Unit-Test liegt außerhalb der Einheit, die er
>   prüft, und greift nur über eine ausdrücklich für Tests vorgesehene Stelle auf
>   Internes zu. (4) Alle drei Prüfungen laufen ohne Eingaben in der Prüfumgebung des
>   Projekts, und jede Prüfung wird nachweislich rot, wenn ihre Bedingung verletzt ist.

**Randformen** (vor dem Lastenheft-Commit entschieden; `AGENTS.md` §3.12):

- **Wortlaut** — **entschieden (Nutzer, 2026-10-06):** der Vorschlag oben; „soll“ ist
  dort zu „muss“ und „sollen“ zu „müssen“ geworden, sonst unverändert.
- **Priorität** — **entschieden (Nutzer, 2026-10-06):** MUSS, getragen allein vom
  Wortlaut („muss“, „müssen“); die Anforderung bekommt kein eigenes Feld `Priorität`,
  wie die übrigen Qualitätsanforderungen in §4. Das Produkt ist laut Lastenheft erst fertig, wenn alle MUSS- und
  SOLL-Anforderungen umgesetzt sind; die Anforderung gehört damit zu M3.
- **Version** — geprüft: Das Lastenheft steht auf `Draft` (Version 0.1.0); vor
  `Accepted` ist es frei änderbar, ohne Change Request und ohne Historie-Zeile. Ein
  Versionssprung ist darum nicht nötig. Nach `Accepted` wäre eine neue Anforderung ein
  Minor-Sprung (`harness/conventions.md` §Versionierung des Lastenhefts).
- **Abnahmeszenario** — **entschieden (Team, 2026-10-06):** ein eigenes, Nummer 17 in
  §7, Teil von M3. Vorschlag des Planners für den Text — er bleibt Vorschlag bis zum
  Start dieses Slice; Wortlaut beim Lastenheft-Commit
  dem Nutzer vorzulegen: *„Abnahmeszenario 17 — Prüfbarkeit des Quellcodes. Für den
  abzunehmenden Quellstand laufen die statische Analyse, die Prüfung der
  Testabdeckung und die Prüfung der Testanordnung ohne Eingaben und ohne Befund. Für
  jede der drei Prüfungen wird ein Quellstand vorgelegt, der genau ihre Bedingung
  verletzt; die Prüfung meldet den Verstoß und endet mit einem Fehlerstatus. Bezug:
  die neue Anforderung.“* Offen: ob die Rot-Hälfte auf den Gegenproben der Gates
  beruhen darf (dann nennt das Szenario keine eigenen Quellstände) oder bei der
  Abnahme eigens vorgeführt wird.
- **Nachweis in den Abdeckungstabellen** — **entschieden (Team, 2026-10-06):** Die
  Anforderung wird belegt und erscheint dort nicht als unbelegt. **Offen, entscheidet
  der Architect vor dem Code:** der Weg. `make abdeckung` kennt heute nur
  Deklarationen über `func Test…` in Go-Tests
  (`<Kennung>/Messung`, [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md));
  die Gegenproben sind Bash-Skripte und deklarieren nicht. Kandidaten: (a) eine
  Deklaration an einem Go-Test (welcher Test belegt eine Lint-Regel?); (b) eine
  Erweiterung von `make abdeckung` um Gate-Nachweise, etwa eine Deklaration im Kopf
  einer Gegenprobe mit eigener Nachweisart „Gate“ — braucht eine ADR, die
  [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) ergänzt, und eine
  Erweiterung von `make abdeckung-gegenprobe`; (c) ein anderer Weg. Dazu offen:
  Teilmessung — eine Qualitätsanforderung hat nach [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) genau einen Pfad
  *Messung*, also zählte schon der Lint-Nachweis allein als vollständig, obwohl
  Messmethode 2 und 3 fehlen; entweder je Messmethode ein Pfad, oder die Deklaration
  kommt erst mit dem letzten Slice. Bis zum ersten Nachweis steht die Anforderung in
  keiner Tabelle (die Gesamtsicht führt nur deklarierte Anforderungen) und in der RTM
  als unbelegt — ein Zwischenstand der Reihe, kein Endzustand.
- **Spezifikation** — ob eine Spezifikationsstelle die Anforderung schärft oder die
  ADRs der Gates direkt (wie [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) für Build und Test). Referenz-Richtung: Das
  Lastenheft nennt weder Gate noch ADR.

**Risiken:**

- **Vertrag bindet die Abnahme an Harness-Gates** — eine spätere Lockerung eines Gates
  berührt dann die Abnahme; das ist gewollt, macht aber jede Senkung zur
  Vertragsfrage. — **Ausgang:** — (bei Closure)
- **Zahl oder Werkzeug wandern ins Lastenheft** — eine Schwelle im Vertrag veraltet mit
  der ADR; der Wortlaut sagt darum „festgelegt“. — **Ausgang:** — (bei Closure)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `ce50a10` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/plan-folgt-korrektur-nicht` (8×, verkörpert in `AGENTS.md` §3.9, Sensor
  `make kopf-check`) — der zweite Liefer-Punkt zieht sechs Köpfe nach; der Sensor
  prüft, dass §1 und §2 keine Kennung nennen, die der Kopf nicht führt.
- `BEO-REPO/harness-lesart-ohne-entscheidungsort` (1×) — eine Vertragsanforderung
  gibt den Harness-Gates einen Bezug, aber keinen Entscheidungsort für Lesarten; der
  bleibt die ADR.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
