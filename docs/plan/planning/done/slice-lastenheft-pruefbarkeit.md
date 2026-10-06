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
kann (ein Link auf ihren Anker löst erst auf, wenn die Überschrift im Lastenheft steht,
sonst meldet das Doku-Gate `anchor-missing`; einen Link verlangt es nur für
ADR-Kennungen, nicht für Lastenheft-Kennungen), und
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
- Eine Stelle in Spezifikation oder Sicht für die Anforderung selbst — das Lastenheft
  nennt die Anforderung, ihre technische Schärfung tragen die ADRs der Gates (§6
  *Spezifikation*).
- Die Verträge der Gates und ihre Randformen in der Spezifikation (Abschnitt für
  Harness-Werkzeuge, Technik-Stratum) — ein Folge-Slice übernimmt sie:
  `slice-harness-vertraege-spezifikation` legt den Abschnitt an (nach der Antwort des
  Kurs-Repos auf den Change Request, 2026-10-06; §6 *Ort dieser Randformen*); die
  Folge-Slices der Reihe schreiben ihre Randformen dorthin.
- Produkt-Code und Harness-Werkzeuge — Schicht-Abgrenzung: Der Slice ändert das
  Lastenheft, die Roadmap, die Pläne von sechs Slices (dazu die Reihenfolge in §4 von
  `slice-harness-mutation`), legt `slice-harness-abdeckung-gate` an und schreibt
  [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) mit Index-Zeile. Aus
  Review F-406 ändert er dazu `AGENTS.md` §3.12 (die Spezifikation ist der Ort einer
  Randform auch für Harness-Werkzeuge), führt die befristete Adaption
  [MR-001](../../../../harness/conventions.md#mr-001) ein und löst sie nach der
  Kurs-Antwort wieder auf (Eintrag in `harness/conventions/done/`, Index-Zeile in
  `harness/conventions.md`).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Lastenheft: §4 führt die Anforderung im bestätigten Wortlaut aus §6 unter der
      nächsten freien Kennung der Qualitätsanforderungen, mit Priorität MUSS (§6
      *Priorität*); §7 führt das Abnahmeszenario 17 mit Bezug auf die Anforderung;
      Version und Historie nach §6 *Version* — bestätigt an `b098fb3` (Verifikation,
      Punkt 1: Wortlaut und Szenario per `diff` gleich, zwei Hunks, Proben PR1 bis PR3
      rot).
- [x] Roadmap: Der Trigger von M3 nennt das Abnahmeszenario 17 (vom Planner am
      2026-10-06 vorgezogen, mit Drift-Log-Zeile; bei Closure gegen §7 des
      Lastenhefts geprüft). Der Kopf (`Bezug`) von `slice-harness-lint`,
      `slice-harness-coverage` und der vier Umstellungs-Slices verlinkt die
      Anforderung, die Platzhalter-Sätze dort sind entfernt; `make kopf-check` und
      `make docs-check` sind grün — bestätigt an `b098fb3` (Verifikation, Punkt 2:
      Proben PK1 bis PK7, PD1, PD3 rot); dass die Köpfe *verlinken*, ist gelesen, nicht
      von einem Sensor gehalten (V-64, §1 (3) berichtigt).
- [x] Nachweis-Weg entschieden (Architect, vor dem Lastenheft-Commit, `AGENTS.md`
      §3.12): wie jede der Messmethoden 1 bis 3 in den Abdeckungstabellen belegt
      erscheint und wann die Anforderung als vollständig zählt (§6 *Nachweis*):
      [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) liegt vor und ergänzt
      [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md), mit dem
      Lastenheft-Commit `Accepted` und mit der Kennung im `Bezug`;
      `slice-harness-abdeckung-gate` liegt in `open/`; `slice-harness-lint`,
      `slice-harness-blackbox-einstieg` und `slice-harness-coverage` nennen in ihrer
      DoD, welchen Teil des Nachweises sie liefern — bestätigt an `b098fb3`
      (Verifikation, Punkt 3).
- [x] `make gates` grün — an `b098fb3` (Verifikation, Abschnitt 5); die Closure ändert
      Planungsdokumente, Register und eine Anker-Stelle in `AGENTS.md` §3.12
      (`make docs-check`, `make kopf-check` grün).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8) — Review bis
      `9d85fa4` (F-406 bis F-416); die Nacharbeit `b9eae96` bis `b098fb3` hat die
      Verifikation geprüft (V-62 bis V-66), siehe §7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/lastenheft.md` | update | neue Anforderung in §4 mit Priorität MUSS; Abnahmeszenario 17 in §7 |
| `docs/plan/planning/next/slice-harness-lint.md`, `slice-harness-coverage.md`, `slice-harness-blackbox-*.md` | update | Kopf `Bezug` mit Link auf die Anforderung; DoD-Zeile zum Nachweis nach dem entschiedenen Weg (lint, blackbox-einstieg, coverage); Reihenfolge in §4 *Start* mit `slice-harness-abdeckung-gate` |
| `docs/plan/planning/open/slice-harness-abdeckung-gate.md` | neu | aus der vendored Vorlage kopiert: Erweiterung von `make abdeckung` nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), erste Gate-Deklarationen, Randformen in §6 |
| `docs/plan/planning/open/slice-harness-mutation.md` | update | Reihenfolge in §4 *Start* mit `slice-harness-abdeckung-gate` |
| `docs/plan/planning/in-progress/roadmap.md` | geprüft / update | Trigger von M3 mit Abnahmeszenario 17 steht seit der Planung (2026-10-06); bei Closure gegen §7 abgeglichen. Drift-Log-Zeile 2026-10-06 zur Reihenfolge mit `slice-harness-abdeckung-gate` |
| `docs/plan/adr/0033-gate-nachweise-in-der-abdeckung.md`, `docs/plan/adr/README.md` | neu / update | [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) ergänzt [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) um Gate-Nachweise und geteilte Messung (Architect, Proposed); mit dem Lastenheft-Commit Kennung im `Bezug`, Status `Accepted`, Index nachgezogen |
| `AGENTS.md` §3.12 | update | aus Review F-406: Die Spezifikation (Technik-Stratum) ist der Ort einer Randform auch für Harness-Werkzeuge; Endfassung nach der Antwort des Kurs-Repos |
| `harness/conventions.md`, `harness/conventions/done/MR-001-spezifikations-ort-werkzeugvertraege.md` | neu / update | befristete Adaption [MR-001](../../../../harness/conventions.md#mr-001) bis zur Antwort auf den Change Request; aufgelöst am selben Tag, per `git mv` nach `done/`, Index-Zeile unter *Aufgelöste Adaptionen* |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-replay-semantik-fehlerreplay` liegt in
`done/` (der Wortlaut ist seit 2026-10-06 bestätigt). Erster Slice der Reihe
(Lastenheft, `slice-harness-vertraege-spezifikation`, Lint, die vier
Umstellungs-Slices, `slice-harness-abdeckung-gate`, Coverage, `slice-harness-mutation`). Vor dem Lastenheft-Commit
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
  - *Ort dieser Randformen* (F-406): Spezifikation, Abschnitt für Harness-Werkzeuge (Technik-Stratum), angelegt von `slice-harness-vertraege-spezifikation`, nicht
    [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), die Entscheidung und
    Gründe trägt. Bis zum Code stehen sie hier und in §6 von
    `slice-harness-abdeckung-gate`; jede zugesagte Randform braucht dort einen Fall in
    `make abdeckung-gegenprobe`.
- **Spezifikation** — **entschieden (Architect, 2026-10-06): keine
  Spezifikationsstelle für `LH-QA-07` selbst.** Die Anforderung wird nicht in der
  Spezifikation präzisiert; die ADRs der Gates nennen sie im `Bezug`. Die Verträge der
  Gates, die `LH-QA-07` nachweisen, liegen dagegen im Abschnitt für Harness-Werkzeuge
  der Spezifikation (Antwort des Kurs-Repos auf den Change Request, 2026-10-06), den
  `slice-harness-vertraege-spezifikation` anlegt; neue Gate-ADRs zeigen mit `Schärft:`
  dorthin. Wie das `Schärft: —`
  der angenommenen [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) und [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) dazu steht, entscheidet jener Slice. Die
  Rückführung aus §4 („zusätzlich eine Spezifikationsstelle“) ist nicht ausgelöst:
  Den Abschnitt schreibt nicht dieser Slice.

**Risiken:**

- **Vertrag bindet die Abnahme an Harness-Gates** — eine spätere Lockerung eines Gates
  berührt dann die Abnahme; das ist gewollt, macht aber jede Senkung zur
  Vertragsfrage. Da die Rot-Hälfte von Abnahmeszenario 17 auf den Gegenproben beruht,
  gilt das auch für das Entfernen eines Gegenfalls. — **Ausgang:** entfallen: Die
  Bindung ist die gewollte Folge der Team-Entscheidung vom 2026-10-06 und kein Ereignis,
  das noch eintreten kann. Eine Schwellen-Senkung braucht ohnehin eine ADR
  (`AGENTS.md` §3.6), und eine ADR darf das Lastenheft nicht ändern; wer Szenario 17
  damit unerfüllbar machte, braucht eine Lastenheft-Änderung. Dass ein einzelner
  Gegenfall still wegfällt, ist die Klasse von `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`
  (verkörpert in §3.10, Sensor geplant mit `slice-harness-mutation`), kein eigenes
  Risiko dieses Slice.
- **Zahl oder Werkzeug wandern ins Lastenheft** — eine Schwelle im Vertrag veraltet mit
  der ADR; der Wortlaut sagt darum „festgelegt“. — **Ausgang:** entfallen: Der
  gelieferte Wortlaut von `LH-QA-07` und Szenario 17 nennt keine Zahl, kein Werkzeug,
  keine ADR und keinen Slice (Verifikation, Punkt 1, gelesen); Links auf ADR, Slice
  oder Skript macht das Doku-Gate rot (Proben PR1 bis PR3). Ein Werkzeugname als
  blanker Text bliebe grün (PR4); eine solche Änderung wäre aber eine spätere
  Lastenheft-Änderung, ein anderer Vorgang mit eigener Prüfung, kein Rest dieses
  Slice.

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

- **Was hat funktioniert:** Jede Randform von Lastenheft und ADR war vor dem Lastenheft-Commit entschieden (`a238108` vor `daef83b`), Wortlaut und Szenariotext lagen dem Nutzer vor und stehen Byte für Byte im Lastenheft (Verifikation, Punkt 1). Die Referenz-Richtung hielt: Szenario 17 nennt weder Gate noch Werkzeug, die Proben PR1 bis PR3 werden rot. Der Weg aus F-406 trug in einem Tag: Die Regel des Repos (`AGENTS.md` §3.12) setzte einen Ort voraus, den die Baseline für Harness-Werkzeuge nicht ausdrücklich bot; statt ihn still zu setzen, ging ein Change Request an das Kurs-Repo, die Lücke überbrückte die befristete Adaption [MR-001](../../../../harness/conventions.md#mr-001) mit benanntem Auflösungs-Trigger (`c2b8d32`), und nach der Antwort vom 2026-10-06 wurde sie aufgelöst (`6ca20bc`, `566782c`, `3dfef35`) und der Folge-Slice `slice-harness-vertraege-spezifikation` angelegt (`8e98bc4`). Die Verifikation hat die drei Liefer-Punkte an `b098fb3` bestätigt; der Kopf-Sensor wurde in allen sieben Folge-Plänen rot, sobald `LH-QA-07` aus dem `Bezug` fiel (PK1 bis PK7).
- **Was ging anders als geplant:** Der Slice wuchs um eine Regel- und Adaptions-Arbeit, die §1 nicht nannte: `AGENTS.md` §3.12 und MR-001 (V-63), dazu ein neuer Slice in der Reihe, dem §1, §4 und §6 nicht folgten (V-62); beides mit dieser Closure nachgezogen. Die Begründung (3) in §1 berief sich auf eine Eigenschaft des Doku-Gates, die es nicht gibt (V-64; Probe PD2 grün); berichtigt. Das Review fand drei MEDIUM an den Folge-Plänen: den fehlenden Ort der Randformen (F-406), die zweite Bedingung von Messmethode 3 ohne Rot-Fall (F-407) und drei ungenannte Randformen von `slice-harness-abdeckung-gate` (F-408). Eine Prüfrunde mit Nacharbeit; die Nacharbeit `b9eae96` bis `b098fb3` hat kein eigenes Review, die Verifikation hat jedes Finding nachgeprüft (Abschnitt 2).
  - **V-65, Ausgänge:** (a) Der Satz von [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md), der Form der Deklaration, Pfad-Schreibweise und Fehlformen in den Kopf des Abdeckungs-Skripts legt, bleibt (Accepted, `AGENTS.md` §3.5); Adresse ist `slice-harness-vertraege-spezifikation` §6 *Lesart-Rangfolge bei Doppelung*, die ihn mit dieser Closure nennt: Der Abschnitt der Spezifikation geht vor. (b) F-413 („Messmethode“ in zwei Bedeutungen, Zuschreibung an [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md)) — entfallen: Der Text steht im Kontext und in der Begründung der Regel, nicht in ihrer Anwendung; die Teilung ist in §6 *Teilmessung* dieses Plans und in §6 von `slice-harness-abdeckung-gate` je Anforderung geschrieben (drei Teile für `LH-QA-07`), und ihre Randformen schreibt `slice-harness-abdeckung-gate` in den Abschnitt der Spezifikation, der im Rang über der ADR steht. Eine ersetzende ADR nur für den Wortlaut wäre der Fall, den §3.8 vermeiden soll.
  - **V-66, vermerkt:** Den Kopf dieses Plans und den von `slice-harness-vertraege-spezifikation` hält kein Sensor, weil ihr §1 und §2 `LH-QA-07` nicht nennen (Grenze von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md), Klasse von F-410). Für diesen Plan ohne Wirkung, er geht nach `done/`, wo der Sensor nicht liest; für den offenen Plan ist es dessen Sache bei seinem Start.
  - **F-415** (Messmethode 3 misst schwächer als die Anforderung) bleibt Frage an Validator und Nutzer; der Wortlaut ist bestätigt.
- **Steering-Loop-Eintrag:** Geschärfte Regel: Der Ort einer Randform ist die Spezifikation (Technik-Stratum) auch für Harness-Werkzeuge; eine angenommene Gate-ADR verweist darauf, Randformen nach `Accepted` werden dort fortgeschrieben — liegt in `AGENTS.md §3.12`.
  Auslöser: `BEO-REPO/harness-lesart-ohne-entscheidungsort` (slice-harness-kopf-sensor; mit diesem Slice 2×), beantwortet vom Kurs-Repo auf den Change Request (2026-10-06). Herkunfts-Anker `seit slice-lastenheft-pruefbarkeit` im Satz von §3.12, der die Harness-Werkzeuge nennt (mit der Closure nachgetragen; die Überschrift trägt weiter `seit slice-harness-randformen-vor-code`). **Warum dieser Eintrag:** Er ist der einzige der drei Kandidaten, der eine Regel dieses Repos ändert, und er beseitigt die Ursache einer Beobachtung, die zweimal Nacharbeit gekostet hat (F-380/V-42, F-406). Die Ablehnung eines neuen Register-Ausgangs durch das Kurs-Repo bestätigt nur die bestehende Form (drei Ausgänge, Prosa-Begründung, wo kein Sensor möglich ist); sie ändert hier nichts. V-66 ist eine bekannte, in [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md) entschiedene Grenze ohne Auftreten eines Fehlers; ein Sensor dafür bräuchte eine ersetzende ADR. Retirement-Check von §3.12: wieder aufgetreten (F-408, unten), die Regel bleibt.
- **Beobachtungs-Register (`../observations/`):**
  - `BEO-REPO/harness-lesart-ohne-entscheidungsort/` — Beleg `evidence/slice-lastenheft-pruefbarkeit.md` (F-406), **2×**. Ausgang **gestrichen** in `state.md`, mit Begründung: Die Ursache ist weggefallen — der Ort ist entschieden (Kurs-Antwort, `AGENTS.md` §3.12), die Spezifikation ist fortschreibbar und braucht nach `Accepted` keine Folge-ADR. `gestrichen` ist an die Schwelle nicht gebunden (Register-README §Die drei Ausgänge). Landet eine Randform trotzdem im Skriptkopf oder im Plan, ist das ein Verstoß gegen §3.12 und zählt bei `spec-randform-erst-im-review-entschieden` oder `randform-im-code-entschieden-dann-zurueckgegeben`. Damit ist der Zusatz in §5 von `slice-harness-vertraege-spezifikation` („dazu den Ausgang … im Register“) vorweg erfüllt; der Plan zieht das bei seinem Start nach (§3.9).
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` (verkörpert in `AGENTS.md` §3.10 seit welle-extended-query) — Beleg `evidence/slice-lastenheft-pruefbarkeit.md` (F-407), **11×**. Retirement-Check von §3.10: wieder aufgetreten.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` (verkörpert in `AGENTS.md` §3.11 seit welle-extended-query) — Beleg `evidence/slice-lastenheft-pruefbarkeit.md` (F-407 laut Review-Summary, V-64), **12×**. Retirement-Check von §3.11: wieder aufgetreten.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden/` (verkörpert in `AGENTS.md` §3.12 seit slice-harness-randformen-vor-code) — Beleg `evidence/slice-lastenheft-pruefbarkeit.md` (F-408), **9×**.
  - `BEO-REPO/plan-folgt-korrektur-nicht/` (verkörpert in `AGENTS.md` §3.9 seit welle-walking-skeleton, Sensor `make kopf-check` seit slice-harness-kopf-sensor) — Beleg `evidence/slice-lastenheft-pruefbarkeit.md` (F-410, F-411, F-412, F-414, V-62, V-63), **11×**. Retirement-Check von §3.9: wieder aufgetreten; `make kopf-check` war jedes Mal grün, die Funde liegen in der Urteils-Hälfte.

  Einmalig und nicht eingetragen: F-409 (Zählweise der Liefer-Punkte, behoben), F-413 (Wortlaut einer angenommenen ADR, Ausgang oben), V-65 (Ausgang oben), V-66 (Grenze von [ADR-0032](../../adr/0032-kopf-sensor-fuer-slice-plaene.md)), F-415 und F-416 (INFO). Über der Schwelle stehen nur verkörperte Einträge (`plan-folgt-korrektur-nicht` 11×, `negativtests-fehlen-bei-neuem-vertrag` 11×, `zusage-im-kommentar-weiter-als-pruefung` 12×, `spec-randform-erst-im-review-entschieden` 9×); ihnen gibt diese Closure keinen Ausgang, den Lese-Schritt führt die nächste Welle-Closure.
- **Folge-Slices:** `slice-harness-vertraege-spezifikation` (Abschnitt der Spezifikation für Harness-Werkzeuge; dazu der Satz von [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) aus V-65) und `slice-harness-abdeckung-gate` (Nachweisart Gate, Deklaration von Teil 1 und 3); die übrigen der Reihe (`slice-harness-lint`, die vier Umstellungs-Slices, `slice-harness-coverage`, `slice-harness-mutation`) führen `LH-QA-07` im `Bezug` bzw. ihren Platz in der Reihe, Lint, Einstieg und Coverage dazu ihren Nachweis-Teil in der DoD.
- **Risiken aus §6:** zwei, beide **entfallen** mit Begründung (Bindung der Abnahme an die Gates als gewollte Folge; Wortlaut ohne Zahl und Werkzeug, Proben PR1 bis PR3). Die Randformen in §6 sind Entscheidungen, keine Risiken; die Rückführungen aus §4 sind begründet nicht ausgelöst (§6 *Rückführung*, *Spezifikation*).
- **Drei Paarungen:** Anker — `liegt in` nennt `AGENTS.md §3.12`; `grep -n "seit slice-lastenheft-pruefbarkeit" AGENTS.md` findet ihn in §3.12. Folge-Slice — `slice-harness-vertraege-spezifikation`, `slice-harness-abdeckung-gate` und die übrigen genannten liegen als Datei in `open/`. Register — die fünf genannten Kennungen bestehen als Verzeichnis, jedes mit nicht leerem `evidence/` und einer Datei `slice-lastenheft-pruefbarkeit.md`. Die nächste Welle-Closure prüft erneut.
- **Belege:** Review `docs/reviews/2026-10-06-review-slice-lastenheft-pruefbarkeit.md` (bis `9d85fa4`; F-406 bis F-416), Verifikation `docs/reviews/2026-10-06-verifikation-slice-lastenheft-pruefbarkeit.md` (bis `b098fb3`; V-62 bis V-66; DoD 1 bis 3 bestätigt, `make gates` grün an `b098fb3`), Entscheidung [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (Accepted), Adaption [MR-001](../../../../harness/conventions.md#mr-001) (aufgelöst). Validierung: n/a, der Slice ändert Vertrag und Planung, kein End-Nutzer-Verhalten; F-415 liegt beim Validator.

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
  Ort einer Randform eines Werkzeugvertrags ist die Spezifikation, Abschnitt für Harness-Werkzeuge (Technik-Stratum), angelegt von `slice-harness-vertraege-spezifikation`; die ADR trägt
  Entscheidung und Gründe. Mit diesem Slice ist es das zweite Auftreten (Review F-406),
  der Beleg folgt mit der Closure.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
