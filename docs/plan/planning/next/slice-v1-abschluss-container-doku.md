# Slice slice-v1-abschluss-container-doku: Handbuch, README und Releasing-Dokument zu Container-Image und Binaries

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten)

**Berührte Spec-Stellen:** `LH-FA-16.a` · [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) (gelesen: das Gate des Gebers `slice-v1-abschluss-container` belegt, was die Doku zusagt) · [`SPEC-054`](../../../../harness/vertraege/SPEC-054-meldungskatalog.md) (gelesen: das Gate hält die zwei Handbuch-Tabellen)

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Benutzerhandbuch nennt die Abnahme-Umgebungen (§1 *Voraussetzungen*) und beschreibt in §2 *Installation* das Image mit seinen Plattformen, `make binaries` und `SHA256SUMS`; `README.md` und `docs/maintainer/releasing.md` sagen dasselbe, soweit sie es berühren. Jede Aussage beschreibt den Ist-Zustand des gebauten Binaries und des Images.

**Übernommen aus `slice-v1-abschluss-container`** (dort §1, *Abgegeben*, Schnitt vom 2026-10-17 beim Erreichen der Meldeschwelle): Handbuch §1 *Voraussetzungen* (Abnahme-Umgebungen, `LH-QA-03`), §2 *Installation* (*Container*, *Das Binary aus dem Image*, `make binaries`, `SHA256SUMS`), `README.md` und `docs/maintainer/releasing.md`. Aus `slice-doku-ist-stand` (dort §6, *Installation*, über den Geber) die Plattform-Zeile in §1 *Voraussetzungen* („Zum Bauen: … Zum Ausführen: Linux auf der Architektur des Rechners“), soweit das gebaute Image sie belegt.

**Voraussetzung `slice-v1-abschluss-container`** (Reihenfolge: Schritt 44a, dann dieser Slice als Schritt 44b): Der Geber liefert `make binaries`, die Stufe `runtime`, das Gate `make test-container` und den Smoke; die Beispiele dieses Slice laufen gegen sie.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Container-Bau, Gate, Smoke, Gegenprobe — `slice-v1-abschluss-container`: Das ist die Schicht Container-Bau/Harness; dieser Slice ändert weder `Dockerfile` noch `harness/` noch `tools/`.
- Registry, Tag, Manifestliste und Anhänge am GitHub-Release im Handbuch — `slice-erster-release-veroeffentlichung`: Das Handbuch nennt nur, was das Binary und das lokal gebaute Image heute zeigen (`AGENTS.md` §3.11); ein Pull aus einer Registry gibt es erst mit dem Release. `releasing.md` darf das Verfahren beschreiben, wo es dort bereits steht; neue Zusagen zur Pipeline schreibt jener Slice.
- Verweise des Handbuchs auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews — das Handbuch beschreibt den Ist-Zustand ohne Herkunft (`AGENTS.md` §3.11); `README.md` darf auf `spec/` zeigen, ohne Chronik.
- Beispiele für `darwin` und `windows` — kein Lauf in der Prüfumgebung (`SPEC-055` Punkt 13, Grenze 3); das Handbuch sagt zu ihnen nur, was das Gate belegt (Bau, Kopf, Summe).
- Änderungen an den Tabellen `### Fehlercodes` und `### Warnungen` — kein neuer Code in diesem Slice; berührt ein Satz sie dennoch, hält `make meldungskatalog-check` die Gleichheit mit dem Katalog.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität): `docs/user/benutzerhandbuch.md` nennt in §1 *Voraussetzungen* die Plattformen zum Bauen und Ausführen und in §2 *Installation* das Image (Plattformen `linux/amd64`, `linux/arm64`), `make binaries` mit `SHA256SUMS` und das Binary aus dem Image. Jedes Beispiel liegt als Datei vor und läuft gegen das gebaute Binary bzw. das gebaute Image (`AGENTS.md` §3.11); das Handbuch verweist nicht auf Spezifikation, ADRs, Slices oder Reviews und enthält keine Chronik; `make meldungskatalog-check` bleibt grün.
- [ ] [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten): `README.md` und `docs/maintainer/releasing.md` im Ist-Zustand (Image, `make binaries`, `make test-container`, Abnahme-Umgebungen); die Betriebsdokumentation benennt die Abnahme-Umgebungen; `README.md` ohne Chronik und ohne zu sagen, was noch kommt.
- [ ] `make gates` grün (nach `make gates-schnell`, wenn die Positivliste des Stop-Hooks genügt: nur `docs/user`, `docs/maintainer`, `README.md` und `docs/plan/planning/` ändern sich).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | §1 *Voraussetzungen* (Abnahme-Umgebungen), §2 *Installation* (*Container*, *Das Binary aus dem Image*, `make binaries`, `SHA256SUMS`) |
| `README.md` | update | Ist-Zustand von Image und Binaries |
| `docs/maintainer/releasing.md` | update | auf den Stand der umgesetzten Funktionen prüfen und anpassen; Abnahme-Umgebungen (`LH-QA-03`) |

Die Beispiele des Handbuchs liegen als Dateien vor und laufen gegen das gebaute Binary (`make binaries`) bzw. das Image der Stufe `runtime`; der Weg dafür (vorhandene Beispiel-Dateien des Handbuchs) wird vor dem Schreiben gelesen, kein neuer Test-Mechanismus.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-container` liegt in `done/` (der Geber liefert, was die Beispiele ruft). Die Basis der Metrik ist der HEAD vor dem Anspruch und wird beim Anspruch in diesem Absatz eingetragen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- docs/user docs/maintainer README.md` liegt über 500 (Marke); Basis: beim Anspruch. Schätzung des Planners (in dieser Metrik, `AGENTS.md` §3.14): 180 bis 280, Meldeschwelle 400. Wer 400 überschreitet, meldet es vor dem nächsten Abschnitt. Der natürliche Schnitt liegt zwischen Handbuch und `README.md` mit `releasing.md`.
- `in-progress` → `open`: Ein Beispiel läuft gegen das gebaute Binary oder Image nicht, weil der Geber etwas anderes liefert, als sein Gate belegt — zurück an den Geber (Review-Befund dort), kein Carveout im Handbuch.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen:** Der Slice liefert keinen neuen Vertrag (§3.12), sondern Prosa über einen; die Randformen stehen im Vertrag des Gebers ([`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md)). Für die Sätze des Handbuchs gilt (§3.11):

- Plattformen: Das Image gibt es für `linux/amd64` und `linux/arm64`; belegt ist der Smoke nur auf der Plattform des Docker-Servers, die andere ist gebaut und geprüft (Kopf, Summe), nicht ausgeführt. Das Handbuch sagt es so (Grenze 3 und 4 in `SPEC-055` Punkt 13; Bindung an den Geber, §1).
- `darwin` und `windows`: gebaut, mit `SHA256SUMS`, nicht ausgeführt.
- Registry, Tag, Manifestliste: kommen mit `slice-erster-release-veroeffentlichung`; das Handbuch sagt dazu nichts (§1).
- Aufruf-Voraussetzung: `make binaries` und `make test-container` laufen nacheinander, nicht mit `make -j gates` (Grenze 5, `SPEC-055` Punkt 13); steht im Handbuch nur, wenn es den Aufruf zeigt.

**Risiken:**

- Handbuch und Releasing-Dokument beschreiben Verhalten vor der Software; Abweichungen zeigen sich erst im Bau (aus `slice-v1-abschluss-container` §6, dort *entfallen* mit Verweis hierher) — **Ausgang:** offen bis Closure.
- Ein Handbuch-Beispiel gilt gelesen gegen den Code als richtig und wird vom Binary oder Image abgelehnt (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) — **Ausgang:** offen bis Closure; Gegenmittel: jedes Beispiel als Datei gegen das gebaute Binary bzw. Image.
- Die Voraussetzung `slice-v1-abschluss-container` liegt noch in `in-progress/`; ändert das Review dort einen Namen oder Aufruf, zieht dieser Slice nach — **Ausgang:** offen bis Closure.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (gemergter Stand, 2026-10-17): Kein Eintrag nennt Container, Image oder reproduzierbaren Build; einschlägig sind `BEO-REPO/verhalten-nur-unter-linux-geprueft` und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` für Aussagen des Handbuchs über Umgebungen sowie `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` für die Beispiele; keiner erreicht mit diesem Slice die dritte Berührung als Lücke ohne Folge-Slice, der Stand wird bei der Closure aus den Dateien des Registers gelesen.

**Nachgezählt beim Eintragen der Sendung** (2026-10-17, Planner, `AGENTS.md` §3.13, Geber `slice-v1-abschluss-container`): zwei Liefer-Punkte (§2: Handbuch | `README.md` mit `releasing.md`), kein dritter; Schichten: eine (Nutzer- und Wartungs-Doku: `docs/user/`, `docs/maintainer/`, `README.md`; die Planung zählt nicht). Der Geber zählt im selben Commit nach: ein Liefer-Punkt, eine Schicht.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (die Sub-Area `REPO` ist GF).
