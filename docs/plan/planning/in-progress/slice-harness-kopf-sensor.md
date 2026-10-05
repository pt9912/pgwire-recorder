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

**Bezug:** — (Harness-Arbeit; keine Produkt-Anforderung). Randformen des Vertrags
entschieden in [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md). Herkunft:
`BEO-REPO/plan-folgt-korrektur-nicht`, Lese-Schritt der Closure von
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

**Ziel:** Ein Gate meldet rot, wenn ein lebender Slice-Plan (`open/`, `next/`,
`in-progress/`) in §1 oder §2 eine Kennung aus Lastenheft (`LH-XX-NN`,
`LH-XX-NN.x`), Spezifikation (`SPEC-NNN`) oder Sicht (`ARC-NNN`) nennt, die sein
Kopf nicht führt — geprüft gegen die Vereinigung von `Bezug` und `Berührte
Spec-Stellen`, exakt und ohne Markierung für Nennungen ohne Anspruch. Die Regeln im
Einzelnen (Gegenstand, Kennungsform, Bereich, Abschnitte, Formfehler, Ausgabe) sind
in [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) entschieden; der Code folgt ihnen, §6 zeigt je Randform auf die Nummer.
Geliefert als `make kopf-check` (`tools/harness/kopf-check.sh`) mit der Gegenprobe
`make kopf-check-gegenprobe` (`tools/harness/kopf-check-gegenprobe.sh`), beide an
`GATE_CHECKS` über `harness/mk/kopf-check.mk`; dazu die Kopf-Berichtigung der fünf
Pläne aus §6 *Bestand*.

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
- Pläne in `done/`, flach wie archiviert (Stubs tragen weder §1 noch §2) — ein
  geschlossener Plan ist eingefroren und stand vor `done/` unter dem Sensor;
  [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 1.
- Ob eine Kennung existiert und ob sie im Feld steht, das ihre Klasse nahelegt —
  Urteil bzw. ein anderer Abgleich (gegen die Spec); akzeptiertes Negativ in
  [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) §Konsequenzen.
- Produkt-Code, Spezifikation, Lastenheft — Schicht-Abgrenzung: Der Slice ändert
  Harness-Werkzeuge und ihre Doku.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Prüfskript (bash, ohne Docker) mit einem Gate-Ziel an `GATE_CHECKS`, nach den
      Regeln von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md); der Bestand unter `open/`, `next/` und `in-progress/` ist
      grün, die Befunde aus §6 *Bestand* sind in ihren Plänen berichtigt.
- [ ] Gegenprobe als eigenes Gate-Ziel in Temp-Bäumen: je Kennungs-Klasse (`LH-XX-NN`,
      `LH-XX-NN.x`, `SPEC-NNN`, `ARC-NNN`, je in §1 und §2) ein Plan, dem der Kopf die
      Kennung nicht führt, wird abgelehnt; ein vollständiger Kopf wird angenommen;
      dazu je Nummer von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), die eine Mutation fangen kann, ein Fall. Je Zusage
      des Skripts ist die Mutation gesehen (`AGENTS.md` §3.10).
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
| `tools/harness/kopf-check.sh` | neu | liest Kopf, §1 und §2 jedes Slice-Plans und vergleicht die Kennungs-Mengen |
| `tools/harness/kopf-check-gegenprobe.sh` | neu | Fehlformen in Temp-Bäumen, wie `make abdeckung-gegenprobe`; Fälle je Nummer von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), dazu die akzeptierten Negative; Nr. 9 in einem Temp-Baum mit Makefile und Fragment: beide Ziele in `GATE_CHECKS`, beide enden über `make` bei rotem Skript mit Exit ≠ 0 |
| `harness/mk/kopf-check.mk` | neu | zwei Ziele, `kopf-check` und `kopf-check-gegenprobe`, beide an `GATE_CHECKS` |
| `harness/README.md` | update | zwei Zeilen in §Sensors, Bindung [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) |
| Kopf von fünf Plänen in `open/` (§6 *Bestand*) | update | Bestand gegen den Sensor berichtigt |
| `AGENTS.md` | update | §3.9 nennt den Sensor |
| `docs/plan/planning/observations/BEO-REPO/plan-folgt-korrektur-nicht/state.md` | update (Closure) | Sensor-Zeile verkörpert; geschrieben mit der Closure, nicht mit der Implementation |

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

Randformen des Vertrags — vor dem Code entschieden in [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) (Architect,
2026-10-05; Status `Accepted`).
Je Randform die Nummer der Entscheidung; was dort nicht steht, entscheidet der
Implementer nicht, er gibt es zurück (`AGENTS.md` §3.12).

- **Kennung im Fließtext ohne Anspruch** (Abgrenzung, Herkunft, Bereits-Geliefertes)
  — zählt; keine Markierung. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 6. — **Ausgang:** offen bis Closure.
- **Unterkennung gegen Hauptkennung** — exakte Gleichheit, keine deckt die andere.
  Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 4. — **Ausgang:** offen bis Closure.
- **Pläne in `done/`** — nicht geprüft, flach wie archiviert. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md)
  Nr. 1. — **Ausgang:** offen bis Closure.
- **Kennungen in Code-Spans und Links, Link-Ziele mit Anker** — Auszeichnung ohne
  Belang; Anker sind klein geschrieben und keine Nennung. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 2.
  — **Ausgang:** offen bis Closure.
- **Abschnitts-Erkennung und Regel-Absätze** — an der Nummer `## 1.`/`## 2.`; der
  Absatz `Regeln dieser Sektion` ist ausgenommen. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 6. —
  **Ausgang:** offen bis Closure.
- **Kopf-Form: `—`, Mehrfachnennung, andere Schreibweise, Feldwahl** — leere Menge;
  ohne Belang; keine Kennung; Vereinigung beider Felder. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 2
  und 5. — **Ausgang:** offen bis Closure.
- **Bereich `SPEC-NNN bis SPEC-MMM`** — steht für jede Kennung dazwischen, im Kopf wie
  in §1/§2 (der Bestand führt diese Form an beiden Stellen). Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md)
  Nr. 3. — **Ausgang:** offen bis Closure.
- **Plan ohne Kopf-Feld oder ohne §1/§2; Ausgabe und Exit-Codes** — Befund, kein
  Überspringen; eine Zeile je Befund auf stderr, Exit 1/0/2. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md)
  Nr. 7 und 8. — **Ausgang:** offen bis Closure.
- **Bestand** — Prototyp-Abgleich nach [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) am Stand af67ba2: fünf Pläne in
  `open/`, acht Kennungen (13 mit aufgelöstem Bereich); dieser Plan grün:
  `slice-replay-semantik-fehlerreplay` (`SPEC-003`), `slice-replay-semantik-mismatch`
  (`LH-FA-10.a`), `slice-v1-abschluss-betrieb` (`LH-FA-02.b`, `SPEC-013` bis
  `SPEC-019`, `SPEC-034`), `slice-v1-abschluss-homebrew` (`LH-FA-01`),
  `slice-v1-abschluss-zeitangaben` (`LH-FA-20`). Kein gestufter Start
  ([ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 9): Der Slice ergänzt je Plan `Berührte Spec-Stellen` als Code-Span
  (ein Bereich bleibt ein Bereich). `LH-FA-01` (homebrew) und `LH-FA-20` (zeitangaben)
  nennt §1 nur als Abgrenzung; sie stehen darum nicht im Scope-Feld `Bezug`
  ([ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) §Verglichene Alternativen, C; Review F-376); umformuliert würde
  nur, wo die Nennung in §1/§2 selbst falsch ist — geprüft, keine der acht ist es. Das
  bleibt im ersten Liefer-Punkt; die Rückführungs-Bedingung aus §4 (mehr als
  eine Handvoll Pläne) ist mit fünf Kopfzeilen nicht erreicht. — **Ausgang:** offen bis Closure.
- **Lesarten des ADR-Wortlauts, die der Code trägt** — zur Bestätigung beim Architect
  zurückgegeben (`AGENTS.md` §3.12), nicht hier entschieden: „ganzes Wort“ mit den
  Wortzeichen Buchstabe, Ziffer, Unterstrich (ein Punkt endet das Wort,
  `LH-FA-03.ab` zählt als `LH-FA-03`); nur Dateien flach in den drei Verzeichnissen;
  eine Zeile nur aus Leerzeichen ist eine Leerzeile; ein Bereich darf innerhalb eines
  Absatzes über einen Zeilenumbruch reichen; eine Ablage ohne Lifecycle-Verzeichnisse
  ist grün, Exit 2 nur ohne `docs/plan/planning/`. — **Ausgang:** offen bis Closure.
- **Lesarten des Wortlauts** (Randform-Rückgabe des Implementers nach 1dd6ad6,
  Architect 2026-10-05) — alle fünf bestätigt; sie folgen aus dem Wortlaut von
  [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) und sind Regeldetails im Sinn von `AGENTS.md` §3.8, keine neue
  Entscheidung. Eine ergänzende ADR wiederholte nur Regeltext. Festgehalten sind sie
  im Kopfkommentar von `tools/harness/kopf-check.sh` (ZUSAGE), je mit einem Fall der
  Gegenprobe:
  (1) Wortzeichen sind Buchstabe, Ziffer und Unterstrich; `.ab` hinter einer
  Hauptkennung ist keine Unterkennung `.x`, die Hauptkennung zählt (Nr. 2; laut
  statt still) — `nr2-punkt-endet-wort`;
  (2) „in `open/` …“ heißt flach; die Lifecycle-Position ist das Verzeichnis selbst,
  Unterordner gibt es nur in `done/` (Nr. 1) — `nr1-nicht-geprueft`;
  (3) eine Zeile nur aus Leerzeichen und Tabs ist eine Leerzeile, wie in Markdown
  (Nr. 5, 6) — `nr5-leerzeile-mit-leerzeichen`, `nr6-leerzeichen-nach-regel-absatz`;
  (4) ein Bereich reicht über einen Zeilenumbruch im Absatz; ein Umbruch ist dort
  Leerraum (Nr. 3) — `nr3-bereich-zeilenumbruch`;
  (5) Exit 2 nur ohne `docs/plan/planning/`; ein fehlendes Lifecycle-Verzeichnis
  enthält keinen Plan, Exit 0 ist wahr, und ein Exit 2 machte das Gate rot, sobald ein
  leeres `next/` seine `.gitkeep` verliert (Nr. 8) — `nr8-leere-ablage`.
  Keine Code-Änderung. — **Ausgang:** entfällt mit Closure (bestätigt, im Skript
  verkörpert).
- **Offen beim Architect** (Randform-Rückgabe nach Review F-372 und F-377, `AGENTS.md`
  §3.12; kein Code dazu): (a) Feldmarke nur am Absatzanfang
  ([ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 5/7) — ist eine Zeile direkt nach einer Überschriftzeile (`# …`) oder
  einer Linie `---` ohne Leerzeile ein Absatzanfang? (b) Nicht lesbarer Plan — Befund
  mit Exit 1 steht in keiner Nummer der ADR; welche Zeile meldet ihn (Abschnitt, Text)?
  Bis zur Entscheidung erkennt das Skript eine Feldmarke an jedem Zeilenanfang im Kopf,
  und ein Lesefehler hängt im Ausgang von der Position ab. — **Ausgang:** offen bis Closure.
- **Bereich über eine Absatzgrenze in §1/§2** — das Skript fügt die Absätze eines
  Abschnitts zusammen; `SPEC-013 bis` am Absatzende und `SPEC-015` am Anfang des
  nächsten gelten als Bereich (am 2026-10-05 gegen 1dd6ad6 probiert). Akzeptiertes
  Negativ: Das macht §1/§2 nur strenger, nie still grün, und der Kopf kann es nicht,
  weil die Feldmarke zwischen seinen Absätzen steht. Kein Risiko mit eigenem Ausgang.
- **Akzeptierte Negative** — Existenz der Kennung, Feldzuordnung, Kennungen außerhalb
  von §1/§2, Zeilen `## ` in Codeblöcken: nicht geprüft, Grund in [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md)
  §Konsequenzen. Kein Risiko mit eigenem Ausgang.

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
