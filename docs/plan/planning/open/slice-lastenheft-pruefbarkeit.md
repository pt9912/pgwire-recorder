# Slice slice-lastenheft-pruefbarkeit: Qualitätsanforderung an die Prüfbarkeit des Quellcodes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein
Abnahmeszenario und kein Meilenstein hängt an ihm. Eingesammelt wird er von der
nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-05: nach `slice-replay-semantik-fehlerreplay` und vor dem nächsten großen
Slice (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** — (dieser Slice schreibt die Anforderung; ihre Kennung entsteht mit ihm). Entscheidung des Nutzers vom 2026-10-05.

**Berührte Spec-Stellen:** `lastenheft.md §4`

**Verantwortlich:** —

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
Prüfbarkeit des Quellcodes, im Wortlaut, den der Nutzer bestätigt hat (Vorschlag in
§6); die sechs Slices, die sie umsetzen, führen sie danach im Kopf unter `Bezug`.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-05: Das Lastenheft bekommt eine
Anforderung zur Wartbarkeit bzw. Prüfbarkeit des Quellcodes, damit Lint, Black-Box-Tests
und Testabdeckung einen Vertragsbezug haben. Der Planner legt einen Wortlaut vor; er
schreibt ihn nicht ins Lastenheft.

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
- Die Gates selbst — übernehmen `slice-harness-lint`, `slice-harness-coverage` und die
  vier Umstellungs-Slices.
- Eine Stelle in Spezifikation oder Sicht — das Lastenheft nennt die Anforderung,
  ihre technische Schärfung tragen die ADRs; eine Spezifikationsstelle entsteht nur,
  wenn §6 *Spezifikation* das anders entscheidet.
- Produkt-Code und Harness — Schicht-Abgrenzung: Der Slice ändert das Lastenheft und
  die Köpfe von sechs Slice-Plänen.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] §4 des Lastenhefts führt die Anforderung im vom Nutzer bestätigten Wortlaut,
      unter der nächsten freien Kennung der Qualitätsanforderungen, mit Anforderung und
      Messmethode in der Form der übrigen; Version und Historie des Lastenhefts nach §6
      *Version*.
- [ ] Der Kopf (`Bezug`) von `slice-harness-lint`, `slice-harness-coverage` und der vier
      Umstellungs-Slices verlinkt die Anforderung; die Platzhalter-Sätze dort sind
      entfernt; `make kopf-check` und `make docs-check` sind grün.
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
| `spec/lastenheft.md` | update | neue Anforderung in §4; gegebenenfalls §7 (§6 *Abnahmeszenario*) |
| `docs/plan/planning/open/slice-harness-lint.md`, `slice-harness-coverage.md`, `slice-harness-blackbox-*.md` | update | Kopf `Bezug` mit Link auf die Anforderung |
| `docs/plan/planning/in-progress/roadmap.md` | update, nur falls §6 *Abnahmeszenario* eines anlegt | Trigger von M3 nennt die Abnahmeszenarien einzeln |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-replay-semantik-fehlerreplay` liegt in
`done/`, und der Nutzer hat den Wortlaut aus §6 bestätigt oder geändert. Erster Slice
der Reihe (Lastenheft, Lint, die vier Umstellungs-Slices, Coverage).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Nutzer entscheidet
  zusätzlich ein Abnahmeszenario und eine Spezifikationsstelle; dann trägt ein eigener
  Slice die Abnahme-Seite.
- `in-progress` → `open` (blockiert — Carveout?): Der Wortlaut ist nicht bestätigt.
  Verwirft der Nutzer die Vertragsbindung ganz, geht der Slice aus `open/` nach `done/`
  mit `Gegenstand: entfallen` und Grund; die sechs Köpfe behalten `—`.

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
> - **Anforderung:** Der Quellcode soll so gehalten werden, dass jede Änderung
>   automatisch auf Verstöße gegen festgelegte Struktur- und Komplexitätsregeln und auf
>   nicht durch Tests ausgeführten Code geprüft werden kann. Tests sollen das Verhalten
>   einer Einheit über deren öffentliche Schnittstelle prüfen.
> - **Messmethode:** (1) Eine statische Analyse nach einem festgelegten Regelprofil
>   meldet für den Quellstand keinen Befund; Ausnahmen stehen zentral und begründet,
>   nicht im Code. (2) Die Anweisungsabdeckung der Unit-Tests liegt nicht unter einer
>   festgelegten Schwelle. (3) Jeder Unit-Test liegt außerhalb der Einheit, die er
>   prüft, und greift nur über eine ausdrücklich für Tests vorgesehene Stelle auf
>   Internes zu. (4) Alle drei Prüfungen laufen ohne Eingaben in der Prüfumgebung des
>   Projekts, und jede Prüfung wird nachweislich rot, wenn ihre Bedingung verletzt ist.

**Randformen** (vor dem Lastenheft-Commit vom Nutzer zu entscheiden):

- **Wortlaut** — der Vorschlag oben oder eine Fassung des Nutzers; der Planner schreibt
  ihn nicht ins Lastenheft.
- **Verbindlichkeit** — „soll“ wie die übrigen Qualitätsanforderungen oder „muss“. Das
  Produkt ist laut Lastenheft fertig, wenn alle MUSS- und SOLL-Anforderungen umgesetzt
  sind; die Anforderung gehört damit in jedem Fall zum Umfang von M3.
- **Version** — geprüft: Das Lastenheft steht auf `Draft` (Version 0.1.0); vor
  `Accepted` ist es frei änderbar, ohne Change Request und ohne Historie-Zeile. Ein
  Versionssprung ist darum nicht nötig. Nach `Accepted` wäre eine neue Anforderung ein
  Minor-Sprung (`harness/conventions.md` §Versionierung des Lastenhefts).
- **Abnahmeszenario** — ein eigenes Szenario in §7 (dann ändert sich der Trigger von M3
  in der Roadmap, der die Szenarien einzeln nennt) oder keins, weil die Messmethode
  über Gates prüfbar ist.
- **Nachweis in der Abdeckungs-Tabelle** — `make abdeckung` kennt nur
  Test-Deklarationen (`<Kennung>/Messung`, [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md)); ein Gate trägt keine. Ohne
  Deklaration zeigt die Gesamtsicht die Anforderung nicht als belegt. Offen: eine
  Deklaration an einem Test der Gegenproben (die Bash-Skripte deklarieren heute nicht),
  eine Erweiterung des Abdeckungs-Werkzeugs (dann eine ADR) oder bewusst „belegt durch
  Gate“ ohne Tabellenzeile.
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
