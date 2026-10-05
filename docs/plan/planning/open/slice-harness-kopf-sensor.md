# Slice slice-harness-kopf-sensor: Sensor für den Kopf des Slice-Plans

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
`BEO-REPO/plan-folgt-korrektur-nicht`, Lese-Schritt der Closure von
welle-extended-query.

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Gate meldet rot, wenn ein Slice-Plan in §1 oder §2 eine
Anforderungs-, Spezifikations- oder Architektur-Kennung nennt, die sein Kopf nicht
führt: eine Lastenheft-Kennung (Form `LH-XX-NN`) unter `Bezug`, eine Kennung der
Spezifikation (`SPEC-NNN`, `LH-XX-NN.x`) oder der Sicht (`ARC-NNN`) unter
`Berührte Spec-Stellen`.

**Herkunft:** `AGENTS.md` §3.9 verlangt seit welle-walking-skeleton, dass der Kopf
jeder Korrektur folgt. Die Beobachtung `BEO-REPO/plan-folgt-korrektur-nicht` trat
danach in jedem Slice von welle-extended-query wieder auf (6×), darunter der Kopf
ohne berührte Spec-Stelle in slice-extended-query-replay (F-325, F-334, F-343, V-28).
Nach dem vierten Auftreten über der Schwelle gilt die Prosa-Form als ausgeschöpft
(v6.13.0 · `regelwerk/modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3);
dieser Slice ist der mechanische Sensor für die Hälfte, die maschinell entscheidbar
ist.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- §3 und §6 gegen den Diff halten (welche geänderte Datei fehlt in §3, welches
  Risiko ist überholt) — das ist Urteil über Inhalt, kein Abgleich zweier Felder;
  es bleibt bei Review und Verifikation.
- Folge-Slices auf die Übernahme prüfen (ob der Nehmer nennt, was ihm zugewiesen
  ist) — ein anderer Vorgang mit eigener Randform (Zeiger in Prosa); wird er
  gebraucht, ist er ein eigener Slice aus dem Register.
- Archivierte Stubs unter `done/<welle-id>/` — sie tragen weder §1 noch §2; der
  Sensor hat dort keinen Gegenstand.
- Produkt-Code, Spezifikation, Lastenheft — Schicht-Abgrenzung: Der Slice ändert
  Harness-Werkzeuge und ihre Doku.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Prüfskript (bash, ohne Docker) mit einem Gate-Ziel an `GATE_CHECKS`; der
      Bestand unter `open/`, `next/`, `in-progress/` und flach in `done/` ist grün
      oder im Plan begründet berichtigt.
- [ ] Gegenprobe als eigenes Gate-Ziel in Temp-Bäumen: je Kennungs-Klasse (`LH-XX-NN`
      in §1, `LH-XX-NN.x`, `SPEC-NNN`, `ARC-NNN` in §2) ein Plan, dem der Kopf die
      Kennung nicht führt, wird abgelehnt; ein vollständiger Kopf wird angenommen. Je
      Zusage des Skripts ist die Mutation gesehen (`AGENTS.md` §3.10).
- [ ] `harness/README.md` §Sensors führt beide Ziele mit ihrem Vertrag;
      `AGENTS.md` §3.9 nennt den Sensor für die Kopf-Hälfte der Regel.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben; `state.md` von
      `BEO-REPO/plan-folgt-korrektur-nicht` nennt den Sensor als Zielort (Zeile
      `Sensor:` von `geplant` auf `verkörpert`, mit Anker `seit slice-harness-kopf-sensor`).
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
| Prüfskript unter `tools/harness/` (Name im ersten Lauf) | neu | liest Kopf, §1 und §2 jedes Slice-Plans und vergleicht die Kennungs-Mengen |
| Gegenprobe unter `tools/harness/` | neu | Fehlformen in Temp-Bäumen, wie `make abdeckung-gegenprobe` |
| Fragment unter `harness/mk/` | neu | zwei Ziele, beide an `GATE_CHECKS` |
| `harness/README.md` | update | zwei Zeilen in §Sensors |
| `AGENTS.md` | update | §3.9 nennt den Sensor |
| `docs/plan/planning/observations/BEO-REPO/plan-folgt-korrektur-nicht/state.md` | update | Sensor-Zeile verkörpert |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): jederzeit; keine Abhängigkeit. Je früher, desto
mehr Slices von welle-replay-semantik prüft er schon.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Bestand braucht
  Berichtigungen in mehr als einer Handvoll Pläne, oder die Kennungs-Regel braucht
  eine eigene Entscheidung über Kopf-Formen (dann Bestand und Sensor trennen).
- `in-progress` → `open` (blockiert — Carveout?): Die Kopf-Form ist so uneinheitlich,
  dass der Sensor sie nicht urteilsfrei lesen kann, ohne die Vorlage zu adaptieren
  (dann zuerst `MR-<NNN>` in `harness/conventions.md`).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün mit beiden neuen Zielen, Closure-Notiz mit
Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Randformen, vor dem Code zu entscheiden (Architect oder Nutzer):

- **Kennung im Fließtext ohne Anspruch** — §1 nennt eine Kennung als Abgrenzung
  („nicht die Fehlerantwort“) oder als Herkunft. Zählt sie? Vorschlag: §1 zählt ganz,
  wer abgrenzt, berührt; sonst braucht der Sensor eine Markierung. — **Ausgang:** offen bis Closure.
- **Unterkennung gegen Hauptkennung** — `LH-FA-18.a` in §2: Genügt `LH-FA-18` unter
  `Bezug`, oder muss `LH-FA-18.a` unter `Berührte Spec-Stellen` stehen? Der Bestand
  führt beides. — **Ausgang:** offen bis Closure.
- **Pläne in `done/`** — frieren ein; ein später geschärfter Sensor färbte sie rot.
  Prüfen oder nicht? — **Ausgang:** offen bis Closure.
- **Kennungen in Code-Spans und Links** — gelten gleich; ein Link-Ziel mit Anker
  (`#lh-fa-18--…`) ist keine zweite Nennung. — **Ausgang:** offen bis Closure.

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
welle-extended-query gesichtet. Treffer: `BEO-REPO/plan-folgt-korrektur-nicht` (6×,
verkörpert, Sensor geplant mit diesem Slice) ist der Gegenstand;
`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und
`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (je 6×, verkörpert in
`AGENTS.md` §3.10 und §3.11) treffen das neue Skript wie jeden neuen Vertrag;
`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (1×) trifft ein neues Gate.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
