# Slice slice-harness-lh-links-bestand: Lastenheft- und MR-Kennungen in ADRs und Wartungs-Doku als Links

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er
von der nächsten Welle-Closure. Geschnitten aus `slice-harness-lh-links-pflicht` am
2026-10-09 (Review F-556 zu `slice-harness-d-check-v0-85`); er steht in der Harness-Reihe
unmittelbar vor jenem (§4 *Start*).

**Bezug:** [MR-000](../../../../harness/conventions.md#mr-000--baseline-aussage) (Baseline-Aussage mit dem ID-Schema, dessen Lastenheft- und Adaptions-Kennungen dieser Slice als Links schreibt). Keine Anforderung des Lastenhefts im Scope: Der Slice ändert Verweise in der Dokumentation, nicht das Produkt.

**Berührte Spec-Stellen:** — (die Links zeigen auf das Lastenheft und den Adaptions-Block; keine Spec-Stelle ändert sich)

**Verantwortlich:** —

**Autor:** pt9912 (Planner). **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Jede Lastenheft-Kennung (Präfix `LH-` mit Klasse und Nummer) und jede Kennung des
Adaptions-Blocks (Präfix `MR-` mit drei Ziffern) in den ADRs unter `docs/plan/adr/` und in
der Nutzer- und Wartungs-Doku (`docs/user/` außer den Abdeckungstabellen,
`docs/maintainer/`) ist ein Link auf ihre Definition, auch in Inline-Code; das bestehende
Doku-Gate prüft jedes Linkziel und jeden Anker.

**Sendung aus `slice-harness-lh-links-pflicht`** (`AGENTS.md` §3.13; Schnitt des Planners vom
2026-10-09 nach Review F-556 zu `slice-harness-d-check-v0-85`): Dessen DoD-Punkt *Bestand
verlinkt* nannte alle gescannten, nicht ausgenommenen Dateien. Mit der Schichtteilung in §8
berührten ADRs und Wartungs-Doku dort eine dritte und vierte Schicht neben Spezifikation und
Harness. Dieser Slice verlinkt sie vorab, dazu die Nutzer-Doku außer den
Abdeckungstabellen (heute ohne Treffer); jener behält Muster, Spezifikation, Commands und
den Bestand in Spezifikation, Harness und Planung, und sein Start verlangt diesen Slice in
`done/`. Die Teil-Zusage „bares `LH-`-Token ist rot“ aus `slice-harness-d-check-v0-85` bleibt
dort.

**Messung des Planners** (2026-10-09, Stand `13b17dc`, Suche nach Kennungen ohne
vorangehende Link-Klammer und ohne Anker-Präfix, also auch in Inline-Code): elf
Lastenheft-Kennungen in vier ADRs (0026: 2, 0031: 5, 0034: 1, 0037: 3), sechs in
`docs/maintainer/releasing.md`, keine in `docs/user/` außer den Abdeckungstabellen, keine
MR-Kennung. Gemessen wird beim Start neu (§4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das `ids`-Muster in `.d-check.yml`, die neue Stelle der Spezifikation und der Block
  *Strenges Doc-Gate* der Commands — Schicht-Abgrenzung: Das ist Spezifikation und Harness;
  beides liefert `slice-harness-lh-links-pflicht` (DoD-Punkte *Muster aktiv* und *Aussage
  über das Gate*, startet nach diesem). Ohne Muster ist dieser Slice trotzdem lieferbar:
  Jeder Link löst auf, und das Gate prüft heute Linkziel und Anker.
- Der Bestand in Spezifikation, `harness/`, lebenden Plänen und Roadmap — Schicht-Abgrenzung:
  Er liegt in den Schichten von `slice-harness-lh-links-pflicht` (DoD-Punkt *Bestand
  verlinkt*), der dort mit dem Muster zugleich nachzieht, was bis dahin neu entsteht.
- Die Abdeckungstabellen `docs/user/abdeckung-*.md` — ein anderer Vorgang: `make abdeckung`
  schreibt sie aus den Abdeckungs-Deklarationen der Tests; Links dort sind Arbeit am
  Generator (`slice-harness-lh-links-pflicht` §1).
- Zeitdokumente (`docs/reviews/**`, `docs/plan/planning/done/**`,
  `docs/plan/planning/observations/**`) — Bestand bleibt bewusst stehen: Sie frieren mit
  ihrem Vorgang ein (`slice-harness-lh-links-pflicht` §1).
- Den Wortlaut einer ADR ändern — Schicht-Abgrenzung: Der Slice setzt einen Link um eine
  schon genannte Kennung und ändert keine Aussage; ob das bei einer ADR mit Status
  `Accepted` trägt (`AGENTS.md` §3.5), entscheidet der Nutzer vor dem Code (§6).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **ADRs verlinkt** (Sendung aus `slice-harness-lh-links-pflicht`): Jede Lastenheft- und
      MR-Kennung in `docs/plan/adr/` ist ein Link auf ihren Anker in `spec/lastenheft.md`
      bzw. `harness/conventions.md`, auch in Inline-Code und Tabellen; außer dem Link ändert
      sich in keiner ADR ein Zeichen (Beleg in §7: Messung vorher und nachher, Wort-Diff
      ohne andere Änderung). Entfällt, wenn der Nutzer vor dem Code entscheidet,
      dass angenommene ADRs keine Links bekommen (§6).
- [ ] **Nutzer- und Wartungs-Doku verlinkt** (Sendung aus `slice-harness-lh-links-pflicht`):
      dasselbe für `docs/user/` außer den Abdeckungstabellen und für `docs/maintainer/`
      (Beleg in §7: Messung vorher und nachher).
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
| `docs/plan/adr/*.md` mit Treffer der Messung (heute 0026, 0031, 0034, 0037) | update | Kennung als Link auf ihren Anker; in Inline-Code wird der Link um den Code gelegt, wie im ADR-Index |
| `docs/maintainer/releasing.md`, Dateien unter `docs/user/` mit Treffer beim Start (heute keine) | update | dasselbe |
| dieser Plan, §7 | update | Messung beim Start und nach der Änderung, Entscheidung des Nutzers zu §3.5 |

- Der Anker einer Kennung mit Unterpunkt (Form `.a`) ist der Anker ihrer Anforderung im
  Lastenheft; der Linktext bleibt die genannte Kennung mit Unterpunkt. Gibt es den Anker
  nicht, meldet das Gate `anchor-missing`, und die Stelle geht an den Planner.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Der Nutzer hat beim Übergang `open` → `next` entschieden,
ob angenommene ADRs Links bekommen (§6); WIP-Limit 1. Platz in der Harness-Reihe
unmittelbar vor `slice-harness-lh-links-pflicht`, also nach welle-v1-abschluss und nach
`slice-harness-mutation`. Beim Start misst der Implementer die Treffer aus §1 neu und trägt
sie in §7 ein.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Messung beim Start findet
  Treffer in einer Datei, die weder in den Schichten dieses Slice noch in denen von
  `slice-harness-lh-links-pflicht` liegt (§8, Schichtteilung); dann zuerst die Zuordnung
  durch den Planner.
- `in-progress` → `open` (blockiert — Carveout?): Eine Kennung hat keinen Anker im
  Lastenheft oder im Adaptions-Block; dann zuerst die Klärung, ob die Kennung stimmt —
  eine falsche Kennung in einer angenommenen ADR berichtigt nur eine neue ADR (`AGENTS.md`
  §3.5).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (Messung in §7 nachher ohne Treffer in `docs/plan/adr/`, `docs/user/` außer
den Abdeckungstabellen und `docs/maintainer/`, `make docs-check` 0 Befunde), Review-Report
liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12): keine. Der Slice liefert keinen neuen Vertrag; die
Zusage des Gates über Linkziel und Anker besteht. Die Linkpflicht und ihre Randformen
entscheidet `slice-harness-lh-links-pflicht`.

**Offener Punkt vor dem Code:** Ob ein Link um eine schon genannte Kennung in einer ADR mit
Status `Accepted` eine inhaltliche Änderung ist (`AGENTS.md` §3.5). Entscheidet der Nutzer
beim Übergang `open` → `next`; bei *ja* entfällt DoD-Punkt *ADRs verlinkt*, und
`slice-harness-lh-links-pflicht` nimmt angenommene ADRs aus (dort §6, Randform
*Ausnahmen*).

**Risiken:**

- Ein Link verändert eine Aussage, etwa weil eine Kennung in einer Tabelle als Datum steht —
  **Ausgang:** offen bis Closure; der Beleg in §7 zeigt per `--word-diff`, dass nur Links
  hinzukamen.
- Bis zum Start von `slice-harness-lh-links-pflicht` entstehen neue blanke Kennungen in
  ADRs oder Nutzer- und Wartungs-Doku — **Ausgang:** offen bis Closure; der Start jenes
  Slice misst neu, und seine Rückführung nennt diesen Fall.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel
`REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `13b17dc` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/slice-waechst-durch-uebernahmen` (3×, verkörpert) und
  `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (5×, verkörpert) — dieser Slice ist Nehmer
  einer Sendung aus `slice-harness-lh-links-pflicht` (§1). **Nachgezählt mit der Sendung**
  (`AGENTS.md` §3.13, im Commit des Schnitts): zwei Liefer-Punkte (ADRs verlinkt, Nutzer-
  und Wartungs-Doku verlinkt), zwei Schichten nach der Schichtteilung unten (Entscheidungen;
  Nutzer- und Wartungs-Doku); der eigene Plan zählt nicht.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — nicht berührt: Der
  Slice ändert keine Gate-Konfiguration; kein Beleg erwartet.

Keiner der Einträge erreicht mit diesem Plan neu die Schwelle 3×.

**Schichtteilung** (gleich in `slice-harness-d-check-v0-85`, `slice-harness-lh-links-pflicht`
und diesem Slice; festgelegt am 2026-10-09 nach Review F-556): Eine Schicht im Sinn von
Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice ist ein Bereich mit
eigenem Maßstab im Review — die Spezifikation (`spec/`) · die Entscheidungen
(`docs/plan/adr/`) · der Harness (Gate-Konfiguration und Gate-Werkzeuge, `Dockerfile`,
Agenten-Anweisungen: `AGENTS.md`, Skills, Commands, `harness/`) · die Nutzer- und
Wartungs-Doku (`docs/user/`, `docs/maintainer/`) · im Produkt-Code je Schicht des Hexagons,
ebenso deren Tests. So zählten `slice-lastenheft-pruefbarkeit` (Harness und
Lastenheft mit Spezifikation als verschiedene Schichten) und `slice-harness-lint-werkzeug` (Harness mit
`Dockerfile` als eine Schicht, Tests je Schicht des Hexagons). Keine Schicht ist die
Planung (Pläne, Roadmap, Register): In ihr schreibt jeder Slice (eigener Plan, Sendungen,
Drift-Log); was sie an Umfang trägt, misst das dritte Kriterium, eine Review-Sitzung.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
