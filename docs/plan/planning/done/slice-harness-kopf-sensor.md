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
in [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) entschieden, §6 zeigt je Randform auf die Nummer. Sieben
Lesarten ihres Wortlauts hat der Architect bestätigt bzw. entschieden; sie stehen in §6
und im Kopfkommentar (ZUSAGE) von `tools/harness/kopf-check.sh`, nicht in der ADR, die
nach der Annahme gesperrt ist (§7, F-380). Der Code folgt ADR und Lesarten.
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

- [x] Prüfskript (bash, ohne Docker) mit einem Gate-Ziel an `GATE_CHECKS`, nach den
      Regeln von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md); der Bestand unter `open/`, `next/` und `in-progress/` ist
      grün, die Befunde aus §6 *Bestand* sind in ihren Plänen berichtigt — bestätigt an
      `7c29e5f` (Verifikation, Punkt 1); Reste V-40 und V-43 stehen in §6.
- [x] Gegenprobe als eigenes Gate-Ziel in Temp-Bäumen: je Kennungs-Klasse (`LH-XX-NN`,
      `LH-XX-NN.x`, `SPEC-NNN`, `ARC-NNN`, je in §1 und §2) ein Plan, dem der Kopf die
      Kennung nicht führt, wird abgelehnt; ein vollständiger Kopf wird angenommen;
      dazu je Nummer von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), die eine Mutation fangen kann, ein Fall. Je Zusage
      des Skripts ist die Mutation gesehen (`AGENTS.md` §3.10) — bestätigt an `7c29e5f`
      (Verifikation, Punkt 2 und Abschnitt 3: 70 Mutationen, 67 rot, 3 äquivalent).
- [x] `harness/README.md` §Sensors führt beide Ziele mit ihrem Vertrag;
      `AGENTS.md` §3.9 nennt den Sensor für die Kopf-Hälfte der Regel — bestätigt an
      `7c29e5f` (Verifikation, Punkt 3); den Anker in der Sensors-Zeile und im Fragment
      trägt die Closure nach (§7, Drei Paarungen).
- [x] `make gates` grün — an `7c29e5f` (Verifikation, Abschnitt 6) und an `1999c04`; die
      Closure ändert Planungsdokumente, Register und Anker-Zeilen.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8) — Review bis
      `3eb3555` (F-372 bis F-380); die Nacharbeit `6ad6319` bis `7c29e5f` hat die
      Verifikation geprüft, siehe §7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben; `state.md` von
      `BEO-REPO/plan-folgt-korrektur-nicht` nennt den Sensor als Zielort (Zeile
      `Sensor:` von `geplant` auf `verkörpert`, mit Anker `seit slice-harness-kopf-sensor`).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — bei Closure
      geprüft (§7); die nächste Welle-Closure prüft sie erneut.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/kopf-check.sh` | neu | liest Kopf, §1 und §2 jedes Slice-Plans und vergleicht die Kennungs-Mengen |
| `tools/harness/kopf-check-gegenprobe.sh` | neu | Fehlformen in Temp-Bäumen, wie `make abdeckung-gegenprobe`; Fälle je Nummer von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), dazu die akzeptierten Negative; Nr. 9 in einem Temp-Baum mit Makefile und Fragment: beide Ziele in `GATE_CHECKS`, beide enden über `make` bei rotem Skript mit Exit ≠ 0; Feldmarke nur am Absatzanfang, nicht lesbarer Plan als Befund |
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
  — zählt; keine Markierung. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 6. — **Ausgang:** entfallen: so umgesetzt, `nr6-abgrenzung` und
  `nr6-ausserhalb` halten es (Verifikation, Abschnitt 3).
- **Unterkennung gegen Hauptkennung** — exakte Gleichheit, keine deckt die andere.
  Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 4. — **Ausgang:** entfallen: so umgesetzt, `nr4-unter-deckt-haupt` und
  `nr4-haupt-deckt-unter` fangen beide Richtungen (Verifikation, Abschnitt 3).
- **Pläne in `done/`** — nicht geprüft, flach wie archiviert. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md)
  Nr. 1. — **Ausgang:** entfallen: so umgesetzt, `nr1-*` hält es; ein Plan in `done/`
  bleibt grün (Verifikation, Abschnitt 4).
- **Kennungen in Code-Spans und Links, Link-Ziele mit Anker** — Auszeichnung ohne
  Belang; Anker sind klein geschrieben und keine Nennung. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 2.
  — **Ausgang:** entfallen: so umgesetzt, `nr2-code-span`, `nr2-codeblock` und die
  Anker-Fälle halten es (Verifikation, Abschnitt 3).
- **Abschnitts-Erkennung und Regel-Absätze** — an der Nummer `## 1.`/`## 2.`; der
  Absatz `Regeln dieser Sektion` ist ausgenommen. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 6. —
  **Ausgang:** entfallen: so umgesetzt; die Lücken der ersten Gegenprobe (F-375) sind
  mit `nr6-ueberschrift-zaehlt` und `nr6-leerzeichen-nach-regel-absatz` geschlossen
  (Verifikation, Abschnitt 2).
- **Kopf-Form: `—`, Mehrfachnennung, andere Schreibweise, Feldwahl** — leere Menge;
  ohne Belang; keine Kennung; Vereinigung beider Felder. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 2
  und 5. — **Ausgang:** entfallen: so umgesetzt, `nr5-*` und `nr7-leer` halten es; die
  Position der Feldmarke (F-372) ist Lesart (a) unten.
- **Bereich `SPEC-NNN bis SPEC-MMM`** — steht für jede Kennung dazwischen, im Kopf wie
  in §1/§2 (der Bestand führt diese Form an beiden Stellen). Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md)
  Nr. 3. — **Ausgang:** entfallen: so umgesetzt, `nr3-*` hält es; am echten Bestand
  wird ein verkürzter Bereich rot (Verifikation, Abschnitt 4).
- **Plan ohne Kopf-Feld oder ohne §1/§2; Ausgabe und Exit-Codes** — Befund, kein
  Überspringen; eine Zeile je Befund auf stderr, Exit 1/0/2. Entschieden: [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md)
  Nr. 7 und 8. — **Ausgang:** entfallen: so umgesetzt, `nr7-*` und `nr8-*` halten es; den
  nicht lesbaren Plan regelt Lesart (b) unten, das nicht lesbare Verzeichnis nicht (eigener
  Eintrag unten).
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
  eine Handvoll Pläne) ist mit fünf Kopfzeilen nicht erreicht. — **Ausgang:** entfallen:
  berichtigt, der Bestand ist grün und jeder der fünf Pläne wird rot, sobald eine Kennung
  aus dem Kopf fällt (Verifikation, Punkt 1 und Abschnitt 4); F-376 ist nachgezogen.
- **Lesarten des ADR-Wortlauts, die der Code trägt** (Stand von 1dd6ad6) — zur Bestätigung beim Architect
  zurückgegeben (`AGENTS.md` §3.12), nicht hier entschieden: „ganzes Wort“ mit den
  Wortzeichen Buchstabe, Ziffer, Unterstrich (ein Punkt endet das Wort,
  `LH-FA-03.ab` zählt als `LH-FA-03`); nur Dateien flach in den drei Verzeichnissen;
  eine Zeile nur aus Leerzeichen ist eine Leerzeile; ein Bereich darf innerhalb eines
  Absatzes über einen Zeilenumbruch reichen; eine Ablage ohne Lifecycle-Verzeichnisse
  ist grün, Exit 2 nur ohne `docs/plan/planning/`. — **Stand:** überholt; der Architect
  hat alle fünf in 3eb3555 bestätigt, siehe den folgenden Eintrag. Dass sie vor der
  Rückgabe schon im Code standen, ist F-373 (§7). — **Ausgang:** entfallen: bestätigt,
  nichts mehr offen.
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
- **Absatzanfang im Kopf und nicht lesbarer Plan** (Randform-Rückgabe nach Review
  F-372 und F-377, `AGENTS.md` §3.12; Architect 2026-10-05, Stand 6ad6319) — beide
  entschieden als Lesart des Wortlauts von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), keine ergänzende ADR:
  (a) **Ein Absatz beginnt nur in der ersten Zeile der Datei oder nach einer
  Leerzeile** (Leerzeile wie Lesart 3). Eine Zeile direkt nach `# …`, nach `---` oder
  nach einer anderen Kopfzeile ist kein Absatzanfang; eine Feldmarke dort ist kein
  Feld, und fehlt das Feld sonst, ist das der Befund aus Nr. 7. Grund: Die ADR nennt
  die Leerzeile als einzige Absatzgrenze (Nr. 5 „bis zur nächsten Leerzeile“); beide
  Enden eines Absatzes an derselben Grenze zu messen ist die Lesart, die keine zweite
  Grenze einführt. Eine Abweichung wird laut (Feld fehlt) und mit einer Leerzeile
  behoben; Bestand und Vorlage stehen nach einer Leerzeile.
  (b) **Ein nicht lesbarer Plan ist ein Befund (Exit 1), gleich an welcher Position**;
  die übrigen Pläne werden weiter gelesen. Befund-Zeile
  `kopf-check: <pfad>: Datei: nicht lesbar`, einzige Zeile dieses Plans, sortiert wie
  jede andere (unter `LC_ALL=C` steht `Datei` vor `Kopf` und `§1`). Grund: Nr. 8
  behält Exit 2 der fehlenden Ablage vor, und ein Plan, den der Sensor nicht
  bestätigen kann, ist nicht bestätigt (Nr. 7: Befund statt Überspringen). Exit 2
  dafür ginge über die ADR hinaus.
  **Vorgabe an den Implementer:** Feldmarke nur, wenn sie am Zeilenanfang steht und
  die Zeile die erste der Datei ist oder auf eine Leerzeile folgt; sonst gehört die
  Zeile zum laufenden Absatz (zum Feld, wenn er eines ist). Vor dem Lesen `-r`
  prüfen, und ein Fehlschlag von `awk` an einem Plan ist derselbe Befund; Exit 2 nur
  aus Nr. 8. ZUSAGE-Kopf um (a) und (b) ergänzen, die zwei Sätze dazu aus GRENZE
  streichen — im selben Commit wie der Code (§3.11). Gegenprobe, je rot mit Befund:
  `**Welle:** x` und in der Folgezeile `**Bezug:** …` → `Kopf: Feld Bezug fehlt`;
  `# Titel` bzw. `---` direkt vor `**Bezug:**` → dasselbe; `**Bezug:** …` direkt
  vor `**Berührte Spec-Stellen:** …` → `Kopf: Feld Berührte Spec-Stellen fehlt`;
  Feldmarke mitten in einer Zeile → Feld fehlt (fängt `index(…) == 1` → `> 0`);
  grün: Feld in der ersten Zeile der Datei. Nicht lesbarer Plan (`chmod 000`) vor
  einem lesbaren und als letzter: je Exit 1 mit der Zeile oben, der lesbare Plan
  meldet seine Befunde trotzdem. Läuft die Gegenprobe als root (Datei bleibt
  lesbar), ist der Fall nicht herstellbar: Er meldet das auf stderr und wird
  übersprungen, ohne die Gegenprobe rot zu machen — laut, nicht still. — **Stand:** umgesetzt;
  Fälle `nr5-welle-vor-bezug`, `nr5-titel-vor-bezug`, `nr5-linie-vor-bezug`,
  `nr5-bezug-vor-stellen`, `nr5-mitten-in-zeile`, `nr5-mitten-in-zeile-stellen`,
  `nr5-erste-zeile`, `nr8-unlesbar-davor`, `nr8-unlesbar-zuletzt`, `nr8-awk-scheitert`
  (`awk`-Ersatz im `PATH`). — **Ausgang:** entfällt mit Closure (umgesetzt, im Skript verkörpert).
- **Nicht lesbares Lifecycle-Verzeichnis** (Verifikation V-40, mit der Closure
  eingetragen) — weder [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) noch §6 nennen es; der Code entscheidet es still und je
  nach Position verschieden: `open/` mit Rechten 000 ergibt Exit 0, ein roter Plan darin
  bleibt ungesehen; `in-progress/` ergibt Exit 1 ohne Befund-Zeile. Die Sensors-Zeile in
  `harness/README.md` („jeder Slice-Plan in `open/` …“) sagt dort mehr zu, als geprüft
  wird. — **Ausgang:** weiter offen → Register
  `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (1×). *Entfallen* wäre falsch, das
  Risiko ist reproduziert. *Eingetreten* ist es nicht: Verzeichnisrechte verfolgt git
  nicht, im Bestand und in jedem Klon sind die Verzeichnisse lesbar. Ein Folge-Slice
  bräuchte zuerst eine Entscheidung des Architect, und die ADR ist angenommen (eine
  ergänzende oder ersetzende ADR); das ist für eine Randform ohne Auftreten über
  git-Inhalt unverhältnismäßig. Der Zähler entscheidet, ob es Muster ist.
- **Symlink als Plan** (Verifikation V-43, mit der Closure eingetragen) — ein
  Symlink `slice-*.md` in `open/` wird nicht geprüft (`find -type f`), Exit 0 auch bei
  rotem Ziel; ZUSAGE und Nr. 1 sagen dazu nichts Eindeutiges. — **Ausgang:** weiter
  offen → Register `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (1×, derselbe
  Vorgang). Nicht als akzeptiertes Negativ, obwohl git Symlinks verfolgt und der Bestand
  keinen führt: Es so zu setzen hieße, eine Randform des Vertrags bei der Closure im
  Plan zu entscheiden, außerhalb von Spezifikation und ADR, also genau die Reihenfolge,
  die `AGENTS.md` §3.12 verbietet und F-373 gerade zeigte. Das entscheidet der Architect,
  wenn der Eintrag aufgenommen wird.
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

- **Was hat funktioniert:** Der Vertrag war vor dem ersten Code-Commit in einer ADR entschieden (2d45101, 5d0b598 vor 1dd6ad6), und §6 zeigte je Randform auf die Nummer. Am Bestand fand der Sensor die 13 Kennungen in fünf Plänen, die kein Kopf führte; berichtigt wurde nur die Kopfzeile. Die Verifikation hat die drei Liefer-Punkte an `7c29e5f` bestätigt: 70 Mutationen, 67 rot aus dem richtigen Grund, 3 äquivalent (V-44), jede Zeile des ZUSAGE-Kopfs mit mindestens einer roten Mutation; an echten Plänen wird das Gate rot, sobald eine Kennung aus dem Kopf fällt (Verifikation, Abschnitt 4). Die Nacharbeit hielt die Reihenfolge von `AGENTS.md` §3.12 ein: Rückgabe von F-372 und F-377 in 6ad6319 ohne Code-Entscheidung, Entscheidung des Architect in 436a385, Code in 7c29e5f.
- **Was ging anders als geplant:** Die erste Runde hat fünf Lesarten des ADR-Wortlauts erst im Code umgesetzt und im selben Commit an den Architect zurückgegeben (1dd6ad6, bestätigt in 3eb3555; Review F-373, HIGH). Inhaltlich folgte daraus nichts, aber §3.12 baut auf die umgekehrte Reihenfolge. Das ist das erste Auftreten unter §3.12, gleich im ersten Slice nach der Regel. Das Review fand zwei Randformen, die vorher niemand genannt hatte: die Position der Feldmarke im Kopf (F-372, HIGH) und den nicht lesbaren Plan (F-377); dazu Lücken der Gegenprobe an Nr. 9 und Nr. 6 (F-374, F-375) und zwei Abgrenzungs-Kennungen im Scope-Feld `Bezug` (F-376). Die Verifikation fand eine dritte Randform, das nicht lesbare Lifecycle-Verzeichnis (V-40), dazu den Symlink als Plan (V-43); beide stehen in §6, weiter offen. V-41 (§1 und ein überholter §6-Eintrag folgten dem Stand nicht) hat diese Closure nachgezogen; der neue Kopf-Sensor war auf diesem Plan grün, denn der Fund liegt in der Hälfte von §3.9, die Urteil bleibt. Eine Prüfrunde mit Nacharbeit; die Nacharbeit 6ad6319 bis 7c29e5f hat kein eigenes Review, die Verifikation hat ihre Wirkung auf jedes Finding geprüft (Abschnitt 2). Zwei Läufe brachen durch API-Fehler 529 ab; ohne Folgen für den Inhalt.
  - **F-380 und V-42, Kandidat für den Steering Loop:** Für Harness-Gates gibt es keine Spezifikation, und die ADR ist nach der Annahme gesperrt (`AGENTS.md` §3.5). §3.12 verlangt die Entscheidung einer Randform in Spezifikation oder ADR, §3.8 hält Regeldetails aus der ADR heraus. Lesarten landen deshalb in §6 und im Skript-Kopf; dauerhaft ist davon nur der Skript-Kopf, §6 wandert ins Archiv. Für Präzisierungen des Wortlauts trägt der Ort, für Lesarten, die eher eine Wahl sind (Lesart (1), F-380; Lesart (b), V-42), weicht er von §3.12 ab, gedeckt durch die Einstufung des Architect. Eingetragen als `BEO-REPO/harness-lesart-ohne-entscheidungsort` (1×); ob die Antwort ein Entscheidungsort für Harness-Verträge ist (etwa ein Abschnitt in `harness/README.md` oder eine ergänzende ADR je Lesart), entscheidet der Zähler, nicht dieser Slice.
- **Steering-Loop-Eintrag:** Neuer Sensor: `make kopf-check` meldet rot, wenn §1 oder §2 eines Slice-Plans in `open/`, `next/` oder `in-progress/` eine Lastenheft-, Spezifikations- oder Sicht-Kennung nennt, die der Kopf nicht führt, und hält damit die Kopf-Hälfte von `AGENTS.md` §3.9 mechanisch
  — liegt in `make kopf-check` (`harness/mk/kopf-check.mk`, `tools/harness/kopf-check.sh`).
  Auslöser: `BEO-REPO/plan-folgt-korrektur-nicht` (slice-walking-skeleton-build-gates, slice-walking-skeleton-record, slice-walking-skeleton-replay, slice-extended-query-record, slice-extended-query-replay, slice-extended-query-lebendpruefung — 6× bei Planung; mit diesem Slice 7×). Entschieden in [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md); Herkunfts-Anker `seit slice-harness-kopf-sensor` im Sensor-Absatz von `AGENTS.md` §3.9, in der Sensors-Zeile von `harness/README.md` und im Kopf von `harness/mk/kopf-check.mk` (die letzten zwei mit der Closure nachgetragen). Retirement-Check von §3.9: wieder aufgetreten (V-41), und zwar in der Urteils-Hälfte, die der Sensor ausdrücklich nicht prüft; die Regel bleibt, der Sensor ersetzt sie nicht.
- **Beobachtungs-Register (`../observations/`):**
  - `BEO-REPO/plan-folgt-korrektur-nicht/` — Zeile `Sensor:` von `geplant` auf `verkörpert`, Zielort und Anker in `state.md`; Beleg `evidence/slice-harness-kopf-sensor.md` (V-41), **7×**.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden/` (verkörpert in `AGENTS.md` §3.12 seit slice-harness-randformen-vor-code) — Beleg `evidence/slice-harness-kopf-sensor.md`: F-372 und F-377 hat erst das Review gefunden, V-40 erst die Verifikation; vorher hatte sie niemand genannt, weder ADR noch §6. **5×**. Retirement-Check von §3.12: wieder aufgetreten, die Regel bleibt.
  - `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben/` — neu, eigener Slug nach Empfehlung des Reviews (Schwerpunkt 5): Der Wortlaut des Eintrags darüber trifft F-373 nur halb, denn der Implementer hat die Lesarten gemeldet, nicht das Review gefunden. Beleg `evidence/slice-harness-kopf-sensor.md` (F-373), **1×**, `offen`. Retirement-Check von §3.12 auch hier: wieder aufgetreten.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` (verkörpert in `AGENTS.md` §3.10 seit welle-extended-query) — Beleg `evidence/slice-harness-kopf-sensor.md` (F-374, F-375; Zuordnung laut Review-Summary). **7×**. Retirement-Check von §3.10: wieder aufgetreten.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` (verkörpert in `AGENTS.md` §3.11 seit welle-extended-query) — Beleg `evidence/slice-harness-kopf-sensor.md` (F-375 laut Review-Summary, dazu F-372 und F-374 an ZUSAGE-Kopf und Fragment-Kommentar; Rest V-40 an der Sensors-Zeile). **8×**. Retirement-Check von §3.11: wieder aufgetreten.
  - `BEO-REPO/gate-uebergeht-ablage-eintrag-still/` — neu, aus den Risiken V-40 und V-43 (weiter offen), Beleg `evidence/slice-harness-kopf-sensor.md` (F-377 behoben, V-40, V-43), **1×**, `offen`.
  - `BEO-REPO/harness-lesart-ohne-entscheidungsort/` — neu, aus F-380 und V-42 (oben), **1×**, `offen`.
  - `BEO-REPO/randform-wellenlos-ohne-architect-vor-code/` — kein Auftreten: Der Slice ist wellenlos, erreichte den Architect aber vor dem Code ([ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) vor 1dd6ad6); die Randformen, die das Review fand, hatte auch der Architect nicht genannt. Das zählt bei `spec-randform-erst-im-review-entschieden`, nicht hier, auch wenn §6 von slice-harness-randformen-vor-code beide Zählungen in Aussicht stellte: Die Identität dieses Eintrags ist der fehlende Weg zum Architect, und der war da.

  Einmalig und nicht eingetragen: F-376 (Abgrenzungs-Kennung im Scope-Feld, behoben; die Feldwahl prüft der Sensor nach [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) bewusst nicht), F-378 (Werkzeugliste im Kommentar, behoben), F-379 und V-44 (äquivalente Mutationen, kein Fehlermuster), V-41 (gezählt bei `plan-folgt-korrektur-nicht`). Über der Schwelle stehen außer dem Sensor dieses Slice nur verkörperte Einträge (`plan-folgt-korrektur-nicht` 7×, `negativtests-fehlen-bei-neuem-vertrag` 7×, `zusage-im-kommentar-weiter-als-pruefung` 8×, `spec-randform-erst-im-review-entschieden` 5×); ihnen gibt diese Closure keinen Ausgang, den Lese-Schritt führt die nächste Welle-Closure.
- **Folge-Slices:** keine. V-40 und V-43 gehen ins Register, nicht in einen Slice (Begründung in §6).
- **Risiken aus §6:** vierzehn, jedes mit genau einem Ausgang — zwölf entfallen (die acht Randformen-Einträge zu Nr. 1 bis 9, *Bestand*, die überholten Lesarten aus 1dd6ad6, die bestätigten Lesarten (1) bis (5) und die Lesarten (a)/(b); die letzten zwei trugen ihren Ausgang schon), zwei weiter offen → `BEO-REPO/gate-uebergeht-ablage-eintrag-still` (V-40, V-43). Die zwei Einträge *Bereich über eine Absatzgrenze* und *Akzeptierte Negative* sind entschiedene Negative ohne eigenen Ausgang; siehe §6.
- **Drei Paarungen:** Anker — `liegt in` nennt `make kopf-check` mit Fragment und Skript; `grep -rn "seit slice-harness-kopf-sensor" AGENTS.md harness/` findet ihn in `AGENTS.md` §3.9 (Sensor-Absatz), in der Zeile `make kopf-check` von `harness/README.md` §Sensors und im Kopf von `harness/mk/kopf-check.mk`; `state.md` von `plan-folgt-korrektur-nicht` nennt ihn. Folge-Slice — keiner genannt. Register — die sieben genannten Kennungen bestehen als Verzeichnis mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft erneut.
- **Belege:** Review `docs/reviews/2026-10-05-review-slice-harness-kopf-sensor.md` (bis `3eb3555`; F-372 bis F-380), Verifikation `docs/reviews/2026-10-05-verifikation-slice-harness-kopf-sensor.md` (bis `7c29e5f`; V-40 bis V-44; DoD 1 bis 3 bestätigt, `make gates` grün an `7c29e5f`), Entscheidung [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) (Accepted) mit den Lesarten in §6; `make gates` grün an `1999c04`. Validierung: n/a, der Slice liefert ein Harness-Gate und keinen End-Nutzer-Wert.

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
