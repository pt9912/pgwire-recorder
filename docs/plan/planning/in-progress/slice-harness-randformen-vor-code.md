# Slice slice-harness-randformen-vor-code: Randformen vor dem Code entscheiden

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
nächsten Welle-Closure.

**Bezug:** — (Harness-Arbeit; keine Produkt-Anforderung). Herkunft:
`BEO-REPO/spec-randform-erst-im-review-entschieden`, Lese-Schritt der Closure von
welle-extended-query.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Ablauf von Planner, Architect und Implementer hat einen Schritt
„Randformen vor dem Code entscheiden“: Ein Slice, der einen neuen Vertrag liefert
(Format, Leser, Protokollrand, Diagnose, Option), nennt in §6 dessen Randformen;
jede ist entschieden — in der Spezifikation, einer ADR oder als Entscheidung des
Nutzers —, bevor der erste Code-Commit entsteht. Trifft der Implementer auf eine
nicht genannte, entscheidet er sie nicht still, sondern gibt sie zurück.

**Herkunft:** In allen vier Slices von welle-extended-query entschied der Code eine
Randform still, und erst Review oder Verifikation fanden sie, Runde um Runde eine
weitere (`BEO-REPO/spec-randform-erst-im-review-entschieden`, 4×: fehlende und
`null`-Felder, Aliase und Merge-Keys im Format; Herunterfahren; „erwartete Query“
in der Extended-Diagnose; `\v` je Serverversion). Der Gegenbeleg steht in
slice-extended-query-lebendpruefung: Die Randformen, die §6 vor dem Code nannte und
eine ADR vor der ersten Code-Zeile entschied, öffnete kein Review wieder; die nicht
genannten fielen erst dort auf.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die vendored Slice-Vorlage unter `.harness/baseline/` — Bestand bleibt: Der Baum
  ist Kurs-Inhalt, byte-genau gehalten (`make baseline-verify`); der Schritt lebt in
  den repo-eigenen Rollen-Typen und Commands. Braucht er doch eine Vorlagen-Abweichung,
  ist das eine Adaption `MR-<NNN>` in `harness/conventions.md`, kein Edit am Baum.
- Ein Sensor, der die Randformen-Liste prüft — ob eine Randform fehlt, ist Urteil;
  maschinell prüfbar wäre nur, dass §6 eine Liste trägt, und das erzeugte
  Pflichterfüllung (v6.13.0 · `regelwerk/modul-05-planning-harness.md`
  §Ziel-Form: Slice, „Keine Mindestzahl, und kein Sensor darauf“).
- Randformen bestehender Verträge nachträglich sammeln — ein anderer Vorgang; die
  offenen gehören in die Slices, die die Verträge berühren.
- Produkt-Code, Spezifikation, Lastenheft — Schicht-Abgrenzung: Der Slice ändert
  Rollen-Typen, Commands und `AGENTS.md`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Planner und Architect: `.claude/commands/plan-welle.md` (Slice bereitstellen)
      verlangt die Randformen-Liste in §6 für jeden neuen Vertrag;
      `.claude/agents/architect.md` prüft sie vor dem Code und führt jede zu einer
      Entscheidung (Spezifikation, ADR oder Frage an den Nutzer).
- [ ] Implementer: `.claude/commands/implement-slice.md` Schritt 12/13 — eine nicht
      genannte Randform wird zurückgegeben (Plan-Defekt-Rücksprung), nicht still
      entschieden.
- [ ] Regel mit Herkunfts-Anker `seit slice-harness-randformen-vor-code` (oder
      `seit welle-<Kennung>`, falls eine Welle ihn einsammelt) in `AGENTS.md` §3;
      `state.md` von `BEO-REPO/spec-randform-erst-im-review-entschieden` von
      `geplant` auf `verkörpert` mit Zielort.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — geprüft von
      der nächsten Welle-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.claude/commands/plan-welle.md` | update | Randformen-Liste in §6 beim Bereitstellen eines Slice |
| `.claude/agents/architect.md` | update | Prüfung und Entscheidung der Liste vor dem Code |
| `.claude/commands/implement-slice.md` | update | Rückgabe einer nicht genannten Randform |
| `AGENTS.md` | update | Hard Rule mit Herkunfts-Anker |
| `docs/plan/planning/observations/BEO-REPO/spec-randform-erst-im-review-entschieden/state.md` | update | Ausgang verkörpert |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): jederzeit; keine Abhängigkeit. Er wirkt erst auf
Slices, die nach ihm geplant werden; vor dem ersten Slice von welle-replay-semantik
geliefert, trifft er deren Pläne.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Schritt verlangt eine
  Abweichung von der vendored Vorlage (dann zuerst eine Adaption, getrennt).
- `in-progress` → `open` (blockiert — Carveout?): keine Bedingung absehbar; der
  Slice ändert nur repo-eigene Prozess-Texte.

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

- Der Schritt bläht Pläne auf: Eine Liste möglicher Randformen ohne Entscheidung ist
  Text, keine Grenze. Gegenmittel ist die Pflicht, jede genannte zu entscheiden, nicht
  die Liste lang zu machen; der Architect hält sie nach seiner Regel „Lösungen, nicht
  neue Hürden“ klein — **Ausgang:** offen bis Closure.
- Was „neuer Vertrag“ heißt, ist Urteil; zu eng gefasst, greift der Schritt bei einer
  Diagnose-Änderung nicht (die Klasse von slice-extended-query-replay) — **Ausgang:**
  offen bis Closure.

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

**Vorgelagert — offene Beobachtungen sichten:** Register am Stand der Closure von
welle-extended-query gesichtet. Treffer: `BEO-REPO/spec-randform-erst-im-review-entschieden`
(4×, geplant mit diesem Slice) ist der Gegenstand; `BEO-REPO/plan-folgt-korrektur-nicht`
(6×, verkörpert, Sensor geplant in slice-harness-kopf-sensor) trifft jede Plan-Änderung
dieses Slice.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
