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
`slice-harness-coverage`, dem letzten Slice der Reihe vor M3; die Reihe endet mit
`slice-harness-mutation`. Eingesammelt wird er von der
nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-05: nach `slice-replay-semantik-fehlerreplay` und vor dem nächsten großen
Slice (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (von diesem Slice geschrieben). Entscheidungen des Nutzers vom 2026-10-05 und 2026-10-06; Abnahmeszenario vom Team entschieden (Koordinator, 2026-10-06), sein Text vom Nutzer bestätigt (2026-10-06). Nachweis-Weg: [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Architect, 2026-10-06; vom Nutzer am 2026-10-06 angenommen, Accepted mit dem Lastenheft-Commit), ergänzt [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md).

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
Anforderung in den Abdeckungstabellen nicht als unbelegt erscheinen darf. Ebenfalls am
2026-10-06 hat der Nutzer den Text des Szenarios bestätigt, [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) angenommen
und `slice-harness-abdeckung-gate` mit seinem Platz in der Reihe entschieden.

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
- Die Erweiterung von `make abdeckung` nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) und die ersten
  Gate-Deklarationen — ein anderer Vorgang am Werkzeug; sie übernimmt
  `slice-harness-abdeckung-gate` (von diesem Slice in `open/` angelegt, zwischen
  `slice-harness-blackbox-einstieg` und `slice-harness-coverage`, §6 *Nachweis*).
  Dieser Slice legt den Weg und die ADR fest.
- Eine Stelle in Spezifikation oder Sicht — das Lastenheft nennt die Anforderung,
  ihre technische Schärfung tragen die ADRs; eine Spezifikationsstelle entsteht nur,
  wenn §6 *Spezifikation* das anders entscheidet.
- Produkt-Code und Harness-Werkzeuge — Schicht-Abgrenzung: Der Slice ändert das
  Lastenheft, die Roadmap, die Pläne von sechs Slices (dazu die Reihenfolge in §4 von
  `slice-harness-mutation`), legt `slice-harness-abdeckung-gate` an und schreibt
  [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) mit Index-Zeile.

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
      erscheint und wann die Anforderung als vollständig zählt (§6 *Nachweis*):
      [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) liegt vor und ergänzt
      [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md), mit dem
      Lastenheft-Commit `Accepted` und mit der Kennung im `Bezug`;
      `slice-harness-abdeckung-gate` liegt in `open/`; `slice-harness-lint`,
      `slice-harness-blackbox-einstieg` und `slice-harness-coverage` nennen in ihrer
      DoD, welchen Teil des Nachweises sie liefern.
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
| `docs/plan/planning/open/slice-harness-lint.md`, `slice-harness-coverage.md`, `slice-harness-blackbox-*.md` | update | Kopf `Bezug` mit Link auf die Anforderung; DoD-Zeile zum Nachweis nach dem entschiedenen Weg (lint, blackbox-einstieg, coverage); Reihenfolge in §4 *Start* mit `slice-harness-abdeckung-gate` |
| `docs/plan/planning/open/slice-harness-abdeckung-gate.md` | neu | aus der vendored Vorlage kopiert: Erweiterung von `make abdeckung` nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), erste Gate-Deklarationen, Randformen in §6 |
| `docs/plan/planning/open/slice-harness-mutation.md` | update | Reihenfolge in §4 *Start* mit `slice-harness-abdeckung-gate` |
| `docs/plan/planning/in-progress/roadmap.md` | geprüft / update | Trigger von M3 mit Abnahmeszenario 17 steht seit der Planung (2026-10-06); bei Closure gegen §7 abgeglichen. Drift-Log-Zeile 2026-10-06 zur Reihenfolge mit `slice-harness-abdeckung-gate` |
| `docs/plan/adr/0033-gate-nachweise-in-der-abdeckung.md`, `docs/plan/adr/README.md` | neu / update | [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) ergänzt [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) um Gate-Nachweise und geteilte Messung (Architect, Proposed); mit dem Lastenheft-Commit Kennung im `Bezug`, Status `Accepted`, Index nachgezogen |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-replay-semantik-fehlerreplay` liegt in
`done/` (der Wortlaut ist seit 2026-10-06 bestätigt). Erster Slice der Reihe
(Lastenheft, Lint, die vier Umstellungs-Slices, `slice-harness-abdeckung-gate`,
Coverage, `slice-harness-mutation`). Vor dem Lastenheft-Commit
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

**Wortlaut** — Vorschlag des Planners, vom Nutzer am 2026-10-06 bestätigt und so in
§4 des Lastenhefts geschrieben (Überschrift `LH-QA-07`, siehe *Nummer und Ort*):

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
- **Nummer und Ort** — **entschieden (Architect, 2026-10-06):** `LH-QA-07`, die nächste
  freie Kennung der Qualitätsanforderungen; Überschrift `### LH-QA-07 — Prüfbarkeit des
  Quellcodes` nach `LH-QA-06` und vor `LH-RB-01`, Anker
  `#lh-qa-07--prüfbarkeit-des-quellcodes` (Umlaut wie bei `LH-QA-03`).
- **Version** — geprüft: Das Lastenheft steht auf `Draft` (Version 0.1.0); vor
  `Accepted` ist es frei änderbar, ohne Change Request und ohne Historie-Zeile. Ein
  Versionssprung ist darum nicht nötig, Kopf (Version, Datum) und §8 bleiben
  unverändert. Nach `Accepted` wäre eine neue Anforderung ein Minor-Sprung
  (`harness/conventions.md` §Versionierung des Lastenhefts).
- **Glossar** — **entschieden (Architect, 2026-10-06): kein Eintrag.** „Statische
  Analyse“, „Anweisungsabdeckung“, „Unit-Test“, „Einheit“ und „öffentliche
  Schnittstelle“ sind Fachbegriffe ohne produkteigene Bedeutung; was sie in diesem Repo
  technisch heißen (Paket als Einheit, Prüfumgebung als Gate-Kette im Container),
  legen die ADRs der Gates fest, nicht der Vertrag (Lastenheft §1). Das Glossar führt
  Begriffe der Produkt-Domäne. Akzeptiertes Negativ: Wer „Einheit“ anders liest, findet
  die Lesart erst in der ADR von `slice-harness-lint`.
- **Verweis aus der Messmethode auf das Szenario** — **entschieden (Architect,
  2026-10-06): keiner.** `LH-QA-02` und `LH-QA-04` nennen ihr Szenario in der
  Messmethode, `LH-QA-01`, `-05`, `-06` nicht; der bestätigte Wortlaut bleibt
  unverändert, die Verbindung trägt der `Bezug` des Szenarios.
- **Abnahmeszenario** — **entschieden (Team, 2026-10-06):** ein eigenes, Nummer 17 in
  §7 nach Abnahmeszenario 16, Teil von M3. **Rot-Hälfte — entschieden (Architect,
  2026-10-06):** Sie beruht auf den Gegenproben der Gates. Eine eigens vorgeführte
  Verletzung wiederholte, was die Gegenprobe in jedem Gate-Lauf zeigt, und wäre ein
  Handschritt, wo Messmethode (4) „ohne Eingaben“ verlangt. Das Szenario nennt darum
  keine eigenen Quellstände, sondern einen automatisierten Gegenfall je Prüfung, der
  mit den Prüfungen läuft; Gate und Gegenprobe nennt es nicht (Lastenheft §1,
  Referenz-Richtung). **Text — entschieden (Nutzer, 2026-10-06):** der Vorschlag des
  Architects, unverändert übernommen:

  > **Abnahmeszenario 17 — Prüfbarkeit des Quellcodes**
  >
  > Für den abzunehmenden Quellstand laufen die drei Prüfungen aus LH-QA-07 — die
  > statische Analyse, die Prüfung der Anweisungsabdeckung der Unit-Tests und die
  > Prüfung der Lage der Unit-Tests — in der Prüfumgebung des Projekts ohne Eingaben und
  > ohne Befund. Zu jeder der drei Prüfungen gehört ein automatisierter Gegenfall, der
  > mit ihnen läuft: Er legt der Prüfung einen Quellstand vor, der genau ihre Bedingung
  > verletzt, und die Prüfung meldet den Verstoß und endet mit einem Fehlerstatus.
  > Bezug: LH-QA-07.

  Im Lastenheft wird `LH-QA-07` im `Bezug` blank geschrieben wie in den Szenarien 1
  bis 16.
- **Nachweis in den Abdeckungstabellen** — **entschieden (Team, 2026-10-06):** Die
  Anforderung wird belegt und erscheint dort nicht als unbelegt. **Weg — entschieden
  (Architect, 2026-10-06): Kandidat (b), [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), ergänzt
  [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md), keine Ablösung.** (a) trägt nicht: Kein Go-Test prüft eine
  Lint-Regel oder eine Abdeckungsschwelle, eine Deklaration daran sagte zu, was ihr Test
  nicht prüft (`AGENTS.md` §3.11). Die Entscheidung:
  - *Nachweisart Gate:* Eine Gegenprobe an `GATE_CHECKS` trägt die Deklaration im Kopf,
    in derselben Form wie ein Test; nur `LH-QA`/`LH-RB`, kein Produktverhalten.
  - *Teilmessung:* Die Messung ist geteilt; jeder Nachweis deklariert seinen Teil mit
    Nummer und Gesamtzahl, vollständig ist die Anforderung mit allen Teilen. Für
    `LH-QA-07` drei Teile: (1) statische Analyse und (3) Lage der Unit-Tests aus der
    Gegenprobe des Lint-Gates, (2) Anweisungsabdeckung aus der Gegenprobe des
    Coverage-Gates. Messmethode (4) ist kein eigener Teil: Die Nachweisart Gate erfüllt
    sie, weil eine Gegenprobe ohne Eingaben in `make gates` läuft und den Verstoß rot
    zeigt. Der ungeteilte Pfad `Messung` (`LH-QA-02`, `LH-QA-06`) bleibt gültig.
  - *Wer deklariert wann:* `slice-harness-lint` und die vier Umstellungs-Slices
    liefern die Gegenprobe, deklarieren aber nicht — das Werkzeug kennt die Form noch
    nicht, und Teil (3) gilt erst, wenn `testpackage` überall scharf ist.
    `slice-harness-abdeckung-gate` (nach `slice-harness-blackbox-einstieg`; Platz vom
    Architect vorgeschlagen, vom Nutzer am 2026-10-06 entschieden) erweitert
    `tools/test/abdeckung.sh` und `make abdeckung-gegenprobe` und deklariert Teil 1
    und 3 an der Gegenprobe des Lint-Gates → Stand *teilweise*.
    `slice-harness-coverage` deklariert Teil 2 an seiner Gegenprobe → *vollständig*,
    die RTM zählt die Anforderung, Abnahmeszenario 17 ist nachweisbar. Bis zu
    `slice-harness-abdeckung-gate` steht die Anforderung in keiner Tabelle und in der
    RTM als unbelegt — Zwischenstand der Reihe, kein Endzustand.
  - *Rückführung aus §4 nicht ausgelöst:* Die Erweiterung muss nicht vor
    `slice-harness-lint` stehen, weil Lint nicht deklariert; dieser Slice liefert
    weiter Lastenheft, Roadmap, Köpfe und ADR, die Erweiterung bekommt ihren eigenen
    Harness-Slice ohnehin. Lint und Coverage haben je drei Liefer-Punkte, die
    Erweiterung passt in keinen von beiden.
  - *Vorgaben für §6 von `slice-harness-abdeckung-gate`* (Randformen, dort vor dem
    Code zu bestätigen): Pfad-Schreibweise `Messung-<i>-von-<n>` mit 2 ≤ n und
    1 ≤ i ≤ n; abgelehnt werden i > n, n < 2, verschiedene n für dieselbe Anforderung
    und `Messung` neben einer geteilten Messung derselben Anforderung; gelesen werden
    Dateien `tools/**/*-gegenprobe.sh`, die Deklaration als `#`-Zeilen im Kopf;
    Nachweis-Spalte ist der Skriptpfad; eigene Tabelle der Nachweisart Gate neben E2E
    und Unit; ob die Gegenprobe an `GATE_CHECKS` hängt, prüft das Skript nicht
    (Grenze im Skriptkopf, Review).
  - *Ort dieser Randformen* (F-406): Für Harness-Werkzeuge ist der Spezifikations-Ort
    im Sinn von §3.12 nach [MR-001](../../../../harness/conventions.md#mr-001) der ZUSAGE-Kopf von `tools/test/abdeckung.sh` mit der
    Vertragszeile von `make abdeckung-check` in `harness/README.md` §Sensors, nicht
    [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), die Entscheidung und
    Gründe trägt. Bis zum Code stehen sie hier und in §6 von
    `slice-harness-abdeckung-gate`; jede zugesagte Randform braucht dort einen Fall in
    `make abdeckung-gegenprobe`.
- **Spezifikation** — **entschieden (Architect, 2026-10-06): keine
  Spezifikationsstelle.** Die Spezifikation beschreibt Produktverhalten und externe
  Verträge; Regelprofil, Schwelle und Testanordnung betreffen den Quellstand und die
  Gate-Kette. Die ADRs der Gates nennen die Anforderung im `Bezug` und tragen
  `Schärft: —` wie [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) und [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md). Die Rückführung aus §4
  („zusätzlich eine Spezifikationsstelle“) ist damit nicht ausgelöst.

**Risiken:**

- **Vertrag bindet die Abnahme an Harness-Gates** — eine spätere Lockerung eines Gates
  berührt dann die Abnahme; das ist gewollt, macht aber jede Senkung zur
  Vertragsfrage. Da die Rot-Hälfte von Abnahmeszenario 17 auf den Gegenproben beruht,
  gilt das auch für das Entfernen eines Gegenfalls. — **Ausgang:** — (bei Closure)
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
  gibt den Harness-Gates einen Bezug, aber keinen Entscheidungsort für Lesarten. Der
  Ort einer Randform eines Werkzeugvertrags ist nach [MR-001](../../../../harness/conventions.md#mr-001) der ZUSAGE-Kopf des
  Skripts mit der Vertragszeile in `harness/README.md` §Sensors; die ADR trägt
  Entscheidung und Gründe. Mit diesem Slice ist es das zweite Auftreten (Review F-406),
  der Beleg folgt mit der Closure.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
