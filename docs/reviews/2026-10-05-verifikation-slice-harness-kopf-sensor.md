# Verifikation: slice-harness-kopf-sensor — 2026-10-05

**Rolle:** Verifier (Modul 11). Ich prüfe, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer bis `3eb3555` geprüft; die Nacharbeit danach hatte kein eigenes Review, ihre Wirkung auf die Findings prüft dieser Bericht mit.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-kopf-sensor.md` (Kopf, §1 bis §4, §6) gegen den Gesamt-Diff `af67ba2..7c29e5f` bei HEAD `7c29e5f`. Commits: `2d45101` und `5d0b598` ([ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) Proposed, Accepted), `1dd6ad6` (Skript, Gegenprobe, Fragment, Bestand), `3eb3555` (Architect: fünf Lesarten), `6fae5e3` (Review-Report F-372 bis F-380), `6ad6319` (Nacharbeit F-374 bis F-378, Rückgabe F-372/F-377), `436a385` (Architect: Lesarten (a) und (b)), `7c29e5f` (Umsetzung (a) und (b)).

**Eingang:**

- die DoD-Liefer-Punkte 1 bis 3, §1, §3, §6 und die Commit-Messages
- [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 1 bis 9 und §Konsequenzen; `AGENTS.md` §3.6, §3.9 bis §3.12
- der [Review-Report](2026-10-05-review-slice-harness-kopf-sensor.md) (F-372 bis F-380) mit seinen damals grün gebliebenen Mutationen

Der Arbeitsbaum war bei meinem Start sauber, HEAD `7c29e5f`. Als Sensor-Beleg des Implementers lagen nur die Commit-Messages vor; eine Liste *Zusage · Mutation · roter Test* (§3.10) habe ich dort nicht gesehen. Die Behauptungen zählen darum nicht als Beleg: Gate, Gegenprobe, Mutationen und `make gates` habe ich selbst gefahren. Alle Mutationen und Proben liefen in Kopien aus `git archive 7c29e5f` im Scratchpad, der Arbeitsbaum blieb unberührt (`git status --porcelain` leer davor und danach).

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — Prüfskript am Gate, Bestand grün und berichtigt. **Bestätigt.**

- `tools/harness/kopf-check.sh` (bash, ohne Docker) hängt über `harness/mk/kopf-check.mk:16` an `GATE_CHECKS`; `make -pq help` zeigt `kopf-check` und `kopf-check-gegenprobe` in `GATE_CHECKS`.
- Bestand: `make kopf-check` im Arbeitsbaum Exit 0, ebenso in der Kopie mit `gawk`, `mawk` und `busybox awk` als `awk`. Laufzeit 0,08 s.
- Die fünf Pläne aus §6 *Bestand* tragen die 13 Kennungen (acht Nennungen, ein Bereich) in `Berührte Spec-Stellen`; nur die Kopfzeile ist geändert (`git diff af67ba2..7c29e5f -- docs/plan/planning/open/`).
- Gegen die Regeln von [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md): Nr. 1 bis 9 sind im Skript umgesetzt und halten jeder Mutation in Abschnitt 3 stand. Reste außerhalb dessen, was ADR und §6 entscheiden: V-40 (nicht lesbares Lifecycle-Verzeichnis), V-43 (Symlink).

### Punkt 2 — Gegenprobe als eigenes Gate-Ziel. **Bestätigt.**

- Je Kennungs-Klasse in §1 und §2 ein abgelehnter Plan: Schleife `nr2-<Kennung>-§1/§2` über `LH-FA-01`, `LH-QA-05.b`, `SPEC-001`, `ARC-002` (`kopf-check-gegenprobe.sh:104`–`:107`).
- Vollständiger Kopf angenommen: `vollstaendig` (`:98`).
- Je Nummer von [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) Fälle, Nr. 9 in einem Temp-Baum mit Makefile und Fragment (`:286`–`:307`); die akzeptierten Negative als grüne Fälle (`:309`–`:317`).
- „Je Zusage des Skripts ist die Mutation gesehen“: von mir gesehen, 70 Mutationen, 67 rot, 3 äquivalent (Abschnitt 3). Jede Zeile des ZUSAGE-Kopfs (`kopf-check.sh:5`–`:36`) hat mindestens eine rote Mutation.

### Punkt 3 — `harness/README.md` §Sensors und `AGENTS.md` §3.9. **Bestätigt.**

- `harness/README.md:73`–`:74`: beide Ziele mit Vertrag und Bindung [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md). Der Vertragstext folgt dem Stand `7c29e5f` (Feldmarke am Absatzanfang, nicht lesbarer Plan, Nr. 9 im Temp-Baum mit Makefile). Er sagt „jeder Slice-Plan in `open/` …“; bei nicht lesbarem Verzeichnis gilt das nicht (V-40).
- `AGENTS.md:184`–`:188`: §3.9 nennt den Sensor für die Kopf-Hälfte mit Herkunfts-Anker und grenzt §1, §3 und §6 als Urteil ab.

### Punkt 4 — `make gates` grün. **Bestätigt** (eigener Lauf, Abschnitt 6).

### Punkt 5 — Review-Report liegt vor. **Bestätigt.** [Report](2026-10-05-review-slice-harness-kopf-sensor.md) von `6fae5e3`, in anderem Kontext erstellt. Die Nacharbeit `6ad6319`..`7c29e5f` hat kein eigenes Review; Abschnitt 2 prüft, was sie an den Findings ändert.

Die Closure-Punkte (Notiz, Register, Risiko-Ausgänge, Paarungen) sind offen. Das ist vor der Closure richtig; Häkchen habe ich keine gesetzt.

## 2. Review-Findings F-372 bis F-380

| Finding | Ergebnis | Beleg |
|---|---|---|
| F-372 (HIGH) | **behoben.** Vorher an den Architect zurückgegeben (`6ad6319` ändert am Skript nur den Kommentar), entschieden als Lesart (a) in §6 (`436a385`), dann umgesetzt (`7c29e5f`). Feldmarke zählt nur in der ersten Zeile der Datei oder nach einer Leerzeile. Sieben Fälle; die Review-Mutation `== 1` → `> 0` ist jetzt rot (M5a), ebenso das Streichen der Absatz-Bedingung (M5b, M5c, M5d, M5j). | `kopf-check.sh:104`–`:108`; `kopf-check-gegenprobe.sh:155`–`:177` |
| F-373 (HIGH) | **nicht durch Code behebbar, für die Nacharbeit eingehalten.** Der Befund betrifft die Reihenfolge in `1dd6ad6`. In der Nacharbeit lief §3.12 richtig: `6ad6319` gibt F-372 und F-377 zurück, ohne sie im Code zu entscheiden; `436a385` (Architect) entscheidet; `7c29e5f` setzt um. Die Register-Zuordnung (eigener Slug oder Evidence) bleibt bei der Closure. | `git show 6ad6319 -- tools/harness/kopf-check.sh` (nur Kommentar); `git log af67ba2..7c29e5f` |
| F-374 (MEDIUM) | **behoben.** Alle vier damals grünen Fragment-Mutationen sind rot: `\|\| true` (M9a), `-@bash` (M9b), Rezept `@true` (M9c), `kopf-check-gegenprobe` aus `GATE_CHECKS` (M9d). Ebenso dieselben an der Gegenprobe und eine `ifdef`-Stufung (M9e bis M9h). | `kopf-check-gegenprobe.sh:286`–`:307` |
| F-375 (MEDIUM) | **behoben.** (a) Leerzeile aus Leerzeichen nach dem Regel-Absatz: Mutation `leer` → `$0 == ""` rot im Abschnitt (M6d, `nr6-leerzeichen-nach-regel-absatz`) und im Kopf (M5f). (b) Überschriftzeile zählt: Streichen der Zeile rot (M6a, `nr6-ueberschrift-zaehlt`). | `kopf-check-gegenprobe.sh:184`, `:192`–`:195`; `kopf-check.sh:100`, `:113` |
| F-376 (MEDIUM) | **behoben.** `LH-FA-01` (homebrew) und `LH-FA-20` (zeitangaben) stehen in `Berührte Spec-Stellen`, `Bezug` führt nur noch den Scope (`LH-FA-19` bzw. `LH-FA-21`). §6 *Bestand* nennt den Grund. | `open/slice-v1-abschluss-homebrew.md:16`; `open/slice-v1-abschluss-zeitangaben.md:16`; Plan §6 *Bestand* |
| F-377 (LOW) | **behoben für den Plan.** Lesart (b): ein nicht lesbarer Plan oder ein Fehlschlag von `awk` ist genau der Befund `Datei: nicht lesbar`, Exit 1, gleich an welcher Position. Fälle davor, zuletzt und mit `awk`-Ersatz; Mutationen M8g, M8h, M8i, M8m rot. Rest eine Ebene höher: V-40. | `kopf-check.sh:53`–`:54`, `:130`–`:134`; `kopf-check-gegenprobe.sh:240`–`:272` |
| F-378 (LOW) | **behoben.** Der Fragment-Kommentar nennt `make`, `mktemp`, `grep`, GNU-`sed` und die coreutils der Gegenprobe; `sed -z` aus `7c29e5f` ist von „GNU sed“ gedeckt. | `harness/mk/kopf-check.mk:5`–`:7` |
| F-379 (INFO) | **unverändert, kein Handlungsbedarf.** `a < b` steht noch; die Mutation ist äquivalent (M3f grün), wie das Review sagte. | `kopf-check.sh:73` |
| F-380 (INFO) | **offen beim Architect und im Steering Loop.** Die Nacharbeit wendet dieselbe Wahl zweimal mehr an (V-42). | Plan §6 |

## 3. Bewusst brechen (Modul 11)

Treiber: je Mutation eine frische Kopie, literale Ersetzung, Abbruch, wenn die Stelle nicht genau einmal vorkommt, dann `bash tools/harness/kopf-check-gegenprobe.sh` in der Kopie. Gezeigt sind die Fälle, an denen die Gegenprobe rot wurde. Eine Mutation zählt als „aus dem richtigen Grund rot“, wenn der rote Fall die gebrochene Zusage benennt.

| Nr. | Mutation (Auswahl) | Ergebnis, roter Fall |
|---|---|---|
| 1 | `-maxdepth 1` gestrichen (M1a); `done` in die Liste (M1b); `-name '*.md'` (M1c); `in-progress` gestrichen (M1d); `slice*.md` ohne Bindestrich (M1e) | alle rot: `nr1-nicht-geprueft` bzw. `nr1-in-progress` |
| 2 | `_` kein Wortzeichen (M2a); Unterkennung `(\.[a-z])?` gestrichen (M2b); Vorzeichen-Prüfung gestrichen (M2c); Lesart (1) `.ab`-Behandlung gestrichen (M2d); `SPEC-NN` zugelassen (M2e); Ziffer kein Wortzeichen (M2f); Nachzeichen-Prüfung gestrichen (M2g); Codeblock übersprungen (M2h); Code-Spans entfernt (M2i) | alle rot: `nr2-schreibweise-*`, `nr2-LH-QA-05.b-§1/§2`, `nr2-punkt-endet-wort`, `nr2-codeblock`, `nr2-code-span` |
| 3 | Bereichsauflösung aus (M3a); Klassengleichheit gestrichen (M3b); Wortgrenze vorn/hinten gestrichen (M3c, M3d); Kopfzeilen ohne Leerraum verbunden, Lesart (4) (M3e); Backticks nicht entfernt (M3h) | alle rot: `nr3-bereich-*`, `nr3-gemischt`, `nr3-kein-ganzes-wort`, `nr3-bereich-zeilenumbruch` |
| 3 | `a < b` gestrichen (M3f); `i <= b` → `i < b` (M3g) | **grün, äquivalent**: F-379; der obere Endpunkt ist stets auch Einzel-Treffer (gleiche Wortgrenzen) |
| 4 | Unterkennung trägt Hauptkennung ein (M4a); Hauptkennung im Kopf deckt Unterkennung (M4b) | rot: `nr4-unter-deckt-haupt`, `nr4-haupt-deckt-unter` |
| 5 | `== 1` → `> 0` (M5a, Review); Absatz-Bedingung an `Bezug` bzw. `Stellen` gestrichen (M5b, M5c); erste Zeile nicht als Absatzanfang (M5d); Feld-Reset an Leerzeile gestrichen (M5e); Leerzeile nur `""` (M5f, Review); `imkopf` bleibt nach `## ` (M5g); nur `Bezug` bzw. nur `Stellen` zählt (M5h, M5i); `nachleer` nicht zurückgesetzt (M5j); Folgezeilen des Felds verworfen (M5k) | alle rot: `nr5-mitten-in-zeile`, `nr5-welle/titel/linie-vor-bezug`, `nr5-bezug-vor-stellen`, `nr5-erste-zeile`, `nr5-nach-leerzeile`, `nr5-leerzeile-mit-leerzeichen`, `nr5-nur-stellen`, `nr5-nur-bezug`, `nr5-mehrzeilig` |
| 6 | Überschriftzeile nicht gezählt (M6a, Review); Regel-Absatz auch mitten im Absatz (M6b); Regel-Reset gestrichen (M6c); Leerzeile nur `""` im Abschnitt (M6d, Review); Erkennung am Titel (M6e); `### ` beendet Abschnitt (M6f); Regel-Absatz zählt (M6g); `absatzanfang` nie 0 (M6h) | alle rot: `nr6-ueberschrift-zaehlt`, `nr6-regel-mitten-im-absatz`, `nr6-nach-regel-absatz`, `nr6-leerzeichen-nach-regel-absatz`, `nr6-titel`, `nr6-unterueberschrift` |
| 7 | je Formfehler-Meldung gestrichen (M7a, M7b, M7c, M7e); `## 20.` als §2 (M7d) | alle rot: `nr7-ohne-bezug/-stellen/-1/-2`, `nr7-leer` |
| 8 | Exit 2 → 0 (M8a); `LC_ALL=C` gestrichen bei `LANG=de_DE.UTF-8` (M8b); Schluss-`sort` gestrichen (M8c); Menge je Abschnitt nicht geleert (M8d); stdout statt stderr (M8e); Datei-Befund geschluckt bzw. Exit 2 (M8g, M8h); `awk`-stderr nicht unterdrückt (M8i); fehlendes Lifecycle-Verzeichnis Exit 2, Lesart (5) (M8j); Exit-Codes vertauscht (M8k, M8l); `awk`-Fehlschlag mit `\|\| true` (M8m); Default-Wurzel falsch (M8n) | alle rot: `nr8-ohne-ablage`, `nr8-sortiert`, `nr8-unlesbar-davor/-zuletzt`, `nr8-awk-scheitert`, `nr8-leere-ablage`, `nr8-wurzel-aktuell` |
| 8 | `[ -r "$plan" ]` gestrichen (M8f) | **grün, äquivalent**: `gawk` und `mawk` enden an einer nicht lesbaren Datei mit 2, `busybox awk` mit 1; der `else`-Zweig liefert denselben Befund |
| 9 | Fragment: `\|\| true`, `-@bash`, `@true` (M9a bis M9c, Review); `kopf-check-gegenprobe` bzw. `kopf-check` aus `GATE_CHECKS` (M9d, M9e; M9d Review); Gegenprobe-Rezept mit `-@` bzw. `\|\| true` (M9f, M9g); `GATE_CHECKS` hinter `ifdef` (M9h) | alle rot: `nr9-scharf` bzw. `nr9-gate-checks` |

**Neue Regeln:** Absatzanfang (M5a bis M5d, M5j), nicht lesbarer Plan (M8g bis M8i, M8m), Nr. 9 scharf (M9a bis M9c, M9f, M9g) — alle rot. **Review-Mutationen, die damals grün blieben:** alle acht gefahren; sieben sind jetzt rot, `a < b` (M3f) bleibt äquivalent grün, wie F-379 sagt.

Grenze: Läuft die Gegenprobe als root, überspringt sie `nr8-unlesbar-davor` und `nr8-unlesbar-zuletzt` laut (§6 sagt das zu). M8g und M8h fängt dann nur noch `nr8-awk-scheitert`. Mein Lauf war als uid 1000; beide Fälle liefen.

## 4. Gate auf dem Bestand und am gemeinten Fall

Bestand: Exit 0 (Abschnitt 1). Am gemeinten Fall, in Kopien echter Pläne, nur die Kopfzeile geändert:

| Probe | Ergebnis |
|---|---|
| `slice-v1-abschluss-betrieb`: `SPEC-034` aus dem Kopf | Exit 1, `§1: SPEC-034 fehlt im Kopf` |
| `slice-v1-abschluss-betrieb`: Kopf-Bereich `SPEC-013` bis `SPEC-018` statt `…019` | Exit 1, `§2: SPEC-019 fehlt im Kopf` (der Bereich in §2 wird aufgelöst) |
| `slice-v1-abschluss-homebrew`: `LH-FA-01` (Abgrenzungs-Nennung) aus dem Kopf | Exit 1, `§1: LH-FA-01 fehlt im Kopf` |
| `slice-replay-semantik-mismatch`: Kopf `LH-FA-10.a` → `LH-FA-10` | Exit 1, `§1: LH-FA-10.a fehlt im Kopf` (Nr. 4) |
| `slice-replay-semantik-fehlerreplay` und `slice-v1-abschluss-zeitangaben` zugleich gebrochen | Exit 1, zwei Zeilen, nach Pfad sortiert |
| `slice-v1-abschluss-homebrew` nach `next/`, `LH-FA-19` aus beiden Kopf-Feldern | Exit 1, `§1` und `§2: LH-FA-19 fehlt im Kopf` |
| derselbe Plan in `done/` | Exit 0 (Nr. 1) |
| dieser Plan, Leerzeile vor `**Bezug:**` entfernt | Exit 1, `Kopf: Feld Bezug fehlt` (Lesart (a)) |
| Plan aus der vendored Vorlage `slice.template.md` in `open/` | Exit 0; Platzhalter `<SPEC-NNN>` sind keine Kennung, die Felder stehen nach Leerzeilen |

Zwei Läufe meiner ersten Probe blieben grün, weil mein `sed` auch die Nennung in §1/§2 änderte; mit Ersetzung nur in der Kopfzeile wurden sie rot (oben). Ein Kopf-Streichen von `SPEC-042` in homebrew bleibt zu Recht grün: §1/§2 nennen `SPEC-042` nicht (Nr. 5, Kopf-Kennung ohne Nennung).

Robustheit: Gegenprobe und Bestand grün mit `mawk` und `busybox awk`. Ein nicht lesbares Lifecycle-Verzeichnis ergibt V-40, ein Symlink als Plan V-43.

## 5. Plan gegen Code (§3.9)

- **Kopf:** `Bezug: —` mit Zeiger auf [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md), `Berührte Spec-Stellen: —`. Stimmt: kein Spec-Stratum im Diff, §1/§2 nennen keine Kennung (das Gate bestätigt es).
- **§1:** Ziel, Gate-Namen, Skripte, Fragment und Bestand stimmen; die Abgrenzungen treffen zu (kein Produkt-Code, keine Spezifikation, `done/` ungeprüft). Der Satz „Die Regeln im Einzelnen … sind in [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) entschieden; … §6 zeigt je Randform auf die Nummer“ ist enger als der Stand: Sieben Lesarten stehen in §6 und im Skript-Kopf, nicht in der ADR (V-41).
- **§3:** Jede Zeile entspricht dem Diff; die Gegenprobe-Zeile nennt Nr. 9 im Temp-Baum, Absatzanfang und nicht lesbaren Plan. `state.md` ist als Closure-Änderung geführt und im Diff unberührt. Das `Makefile` ist unverändert und in §3 nicht geführt, richtig so.
- **§4:** Keine Rückführungs-Bedingung eingetreten (fünf Kopfzeilen, Kopf-Form urteilsfrei lesbar).
- **§6:** Die Einträge zu Nr. 1 bis 9, *Bestand* und Lesarten (a)/(b) folgen dem Code; die genannten Fallnamen existieren alle in der Gegenprobe. Der Eintrag „Lesarten des ADR-Wortlauts, die der Code trägt“ ist überholt und widerspricht dem folgenden (V-41).

**Gegen die Hard Rules:**

- **§3.6:** Kein Gate gelockert; ein neues Gate kommt voll scharf hinzu, `Makefile` und andere Fragmente unberührt.
- **§3.9:** Code- und Plan-Änderungen der Nacharbeit stehen je im selben Commit (`6ad6319`, `7c29e5f`).
- **§3.10:** Je Zusage eine rote Mutation (Abschnitt 3). Die Liste des Implementers lag mir nicht vor; der Nachweis stammt aus meinem Lauf.
- **§3.11:** ZUSAGE-Kopf, README-Zeilen und Fragment-Kommentar sagen nichts zu, was die Gegenprobe nicht hält. Ausnahme in der README-Zeile: „jeder Slice-Plan“ bei nicht lesbarem Verzeichnis (V-40).
- **§3.12:** Für F-372/F-377 eingehalten (Abschnitt 2, F-373). V-40 ist eine weitere Randform, die weder ADR noch §6 nennt; der Code entscheidet sie heute still.

## 6. Gate-Lauf

`make gates` auf `7c29e5f`, sauberer Arbeitsbaum, eigener Lauf: **Exit 0**, 56 s.

- `abdeckung-gegenprobe`, `a-check-negativ`, `commit-msg-gegenprobe`, `run-integration-tests`, `kopf-check-gegenprobe`: gruen
- `a-check`: 0 Befunde; `baseline-verify`: v6.13.0 OK, 54 Dateien; `d-check`: 197 Dateien, 0 Befunde
- `kopf-check` gibt bei Erfolg nichts aus; `make -pq help` zeigt es in `GATE_CHECKS`, `make -s kopf-check` endet mit 0
- `git status --porcelain` danach leer

## 7. Befunde

| ID | Kategorie | Befund | Pfad |
|---|---|---|---|
| V-40 | LOW | Ein nicht lesbares Lifecycle-Verzeichnis entscheidet das Skript still und je nach Position verschieden. `open/` mit Rechten 000: `find` meldet auf stderr, Exit **0**, ein roter Plan darin bleibt ungesehen (falsch grün, auch über `make kopf-check`). `next/`: Exit 1 mit den übrigen Befunden. `in-progress/`: Exit 1 **ohne** Befund-Zeilen, weil `set -e` die Ersetzung abbricht. Grund: Der Status der `for`-Schleife ist der des letzten `find`, `pipefail` sieht nur ihn. Das ist F-377 eine Ebene höher; weder [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) noch §6 nennen die Randform (§3.12). Über git-verfolgten Inhalt nicht herstellbar, darum LOW wie F-377. Für den Architect, nicht für einen Fix ohne Entscheidung. | `tools/harness/kopf-check.sh:48`–`:52`; `harness/README.md:73` |
| V-41 | LOW | Plan folgt dem Stand nicht ganz (§3.9). (a) §6 führt zwei Einträge zu den fünf Lesarten: der erste sagt „zur Bestätigung beim Architect zurückgegeben …, nicht hier entschieden“ mit Ausgang *offen bis Closure*, der zweite „alle fünf bestätigt“ mit Ausgang *entfällt*. Der erste ist überholt und bekäme bei der Closure einen Ausgang für etwas, das schon entschieden ist. (b) §1 sagt, die Regeln im Einzelnen seien in der ADR entschieden und §6 zeige auf die Nummer; sieben Lesarten stehen aber in §6 und im Skript-Kopf. | `docs/plan/planning/in-progress/slice-harness-kopf-sensor.md:194`–`:200`, `:201`–`:220`, `:42`–`:44` |
| V-42 | INFO | Lesart (b) führt eine Befund-Form ein, die Nr. 8 nicht nennt: Abschnitt `Datei`, neben „Abschnitt und Kennung bzw. fehlendes Feld oder Abschnitt“. Das ist eher eine Wahl als eine Präzisierung des Wortlauts, wie Lesart (1) in F-380. Der Architect hat sie als Lesart eingestuft, die Reihenfolge nach §3.12 ist diesmal eingehalten. Wie F-380 gehört die Ortsfrage (§3.8 gegen §3.12 bei Harness-Verträgen) in den Steering Loop, nicht in diesen Slice. | Plan §6 Lesart (b); `kopf-check.sh:24`–`:26`; [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) Nr. 8 |
| V-43 | INFO | Ein Symlink `slice-*.md` in `open/` wird nicht geprüft (`find -type f`), Exit 0, obwohl er auf einen roten Plan zeigt. git verfolgt Symlinks; im Bestand gibt es keinen. ZUSAGE („Dateien slice-*.md“) und Nr. 1 („Dateiname“) sagen dazu nichts Eindeutiges. Bei Bedarf eine Zeile in GRENZE oder eine Rückgabe; kein Handlungsdruck. | `kopf-check.sh:51` |
| V-44 | INFO | Drei äquivalente Mutationen, keine Lücke der Gegenprobe: `a < b` (F-379), `i <= b` → `i < b` (Endpunkt ist stets Einzel-Treffer) und `[ -r "$plan" ]` (der `awk`-Fehlschlag liefert denselben Befund mit `gawk`, `mawk` und `busybox awk`). Die `-r`-Prüfung ist damit Absicherung, keine eigene Zusage. | `kopf-check.sh:54`, `:73` |

**Summary:** 0 HIGH · 0 MEDIUM · 2 LOW · 3 INFO. DoD-Liefer-Punkte 1 bis 3 bestätigt, `make gates` grün auf `7c29e5f`. F-372, F-374 bis F-378 sind behoben, F-373 ist für die Nacharbeit eingehalten und bleibt für die Closure-Zuordnung offen, F-379 und F-380 bleiben INFO. Die Gegenprobe fängt 67 von 70 Mutationen über alle neun Nummern, die drei übrigen sind äquivalent. Am echten Bestand wird das Gate rot, sobald eine Kennung aus dem Kopf fällt. Merge-blockierend ist keiner der Befunde; V-40 ist eine neue Randform für den Architect.
