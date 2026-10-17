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

- [x] [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität): `docs/user/benutzerhandbuch.md` nennt in §1 *Voraussetzungen* die Plattformen zum Bauen und Ausführen und in §2 *Installation* das Image (Plattformen `linux/amd64`, `linux/arm64`), `make binaries` mit `SHA256SUMS` und das Binary aus dem Image. Jedes Beispiel läuft im Handlauf gegen das gebaute Binary bzw. das gebaute Image (`AGENTS.md` §3.11; als Datei nicht, §7); das Handbuch verweist nicht auf Spezifikation, ADRs, Slices oder Reviews und enthält keine Chronik; `make meldungskatalog-check` bleibt grün.
- [x] [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit), [`LH-RB-01`](../../../../spec/lastenheft.md#lh-rb-01--umgang-mit-sensiblen-daten): `README.md` und `docs/maintainer/releasing.md` im Ist-Zustand (Image, `make binaries`, `make test-container`, Abnahme-Umgebungen); die Betriebsdokumentation benennt die Abnahme-Umgebungen; `README.md` ohne Chronik und ohne zu sagen, was noch kommt.
- [x] `make gates` grün (nach `make gates-schnell`, wenn die Positivliste des Stop-Hooks genügt: nur `docs/user`, `docs/maintainer`, `README.md` und `docs/plan/planning/` ändern sich).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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

Umgesetzt (Stand des Implementers): Handbuch §1 *Voraussetzungen* ein Punkt zu den Plattformen; §2 *Installation* ein neuer Abschnitt *Binaries für alle Plattformen bauen* (`make binaries`, Namen, `SHA256SUMS`, `VERSION`) und im Abschnitt *Container* die Plattformen des Images mit dem Bau für `linux/arm64`; `README.md` ein Punkt der Liste; `docs/maintainer/releasing.md` Stand in §1, ein Absatz *Was lokal läuft* und *Abnahme-Umgebungen* in §2, Änderungshistorie 0.6. Registry, Tag, Manifestliste und Anhänge des GitHub-Release stehen in keinem der drei Texte als vorhanden. Nach dem Review (F-1 bis F-3, F-5): Bezugsmenge „fünf Binaries außer dem für Linux auf `amd64`“, Summenprüfung in Unter-Shell, die Plattformsätze in §1 und *Das Binary aus dem Image* auf den Ist-Zustand gezogen, „Docker 29.9.0“ gestrichen.

Die Beispiele des Handbuchs laufen im Handlauf gegen das gebaute Binary (`make binaries`) bzw. das Image der Stufe `runtime`; als Datei liegen sie nicht vor, kein neuer Test-Mechanismus (akzeptierte Lücke, §7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-container` liegt in `done/` (der Geber liefert, was die Beispiele ruft). **Basis:** `2e89f594` (HEAD vor dem Anspruch; die Commits des Anspruchs ändern keinen der gemessenen Pfade).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- docs/user docs/maintainer README.md` liegt über 500 (Marke); Basis: `2e89f594`. Schätzung des Planners (in dieser Metrik, `AGENTS.md` §3.14): 180 bis 280, Meldeschwelle 400. Wer 400 überschreitet, meldet es vor dem nächsten Abschnitt. Der natürliche Schnitt liegt zwischen Handbuch und `README.md` mit `releasing.md`.
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

- Handbuch und Releasing-Dokument beschreiben Verhalten vor der Software; Abweichungen zeigen sich erst im Bau (aus `slice-v1-abschluss-container` §6, dort *entfallen* mit Verweis hierher) — **Ausgang: *entfallen*** — jedes Beispiel lief im Handlauf gegen das gebaute Binary und das Image (Implementer §7, Review, Verifikation V-1 bis V-3), keine Abweichung zum Bau; die Doku folgte dem Geber in `done/`.
- Ein Handbuch-Beispiel gilt gelesen gegen den Code als richtig und wird vom Binary oder Image abgelehnt (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) — **Ausgang:** offen bis Closure; Gegenmittel: jedes Beispiel als Datei gegen das gebaute Binary bzw. Image. **Ausgang: *weiter offen*** — Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Beleg `evidence/slice-v1-abschluss-container-doku.md`, Zähler verkörpert/über 3×, Regel bleibt). Im Handlauf lief kein Beispiel an Binary oder Image vorbei (Implementer, Review, Verifikation), die Gegenmaßnahme „als Datei“ fehlt als akzeptierte Lücke (Nutzerlinie).
- Die Voraussetzung `slice-v1-abschluss-container` lag in `in-progress/`; ändert das Review dort einen Namen oder Aufruf, zieht dieser Slice nach — Voraussetzung erfüllt (der Geber liegt in `done/`), **Ausgang: *entfallen*** — der Geber liegt in `done/`, Namen mit `VERSION` übernommen (Verifikation V-2). Aus der Closure-Notiz des Gebers (§7): Die Dateinamen von `make binaries` tragen `VERSION` (Vorgabe `dev`; das Gate nutzt `0.0.0-test`); die Beispiele zeigen die Namen so, wie `make binaries` sie erzeugt.

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

Closure am 2026-10-17 (Planner), Belege: Review `docs/reviews/2026-10-17-review-slice-v1-abschluss-container-doku.md` (0 HIGH, 0 MEDIUM, 4 LOW, 2 INFO, behandelt), Verifikation `docs/reviews/2026-10-17-verify-slice-v1-abschluss-container-doku.md` (bestanden).

### Closure-Notiz

**Ergebnis:** Handbuch §1 und §2, `README.md` und `docs/maintainer/releasing.md` beschreiben Image (`linux/amd64`, `linux/arm64`), `make binaries` mit `SHA256SUMS` und die Abnahme-Umgebungen im Ist-Zustand; Messung der Marke 111 von 500 (§4). Kriterien: Verifikation bestanden (DoD 1 und 2 gegen Diff und eigene Läufe), `make docs-check`, `make kopf-check`, `make meldungskatalog-check` grün.

**Risiko-Ausgänge** (§6): Doku vor der Software *entfallen* (Handlauf ohne Abweichung); Beispiel vom Binary abgelehnt *weiter offen* (Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`); Voraussetzung Geber *entfallen* (Geber in `done/`).

**Ausgänge der Findings** (Nutzerlinie: LOW/INFO akzeptiert mit Ausgang): F-1 bis F-3, F-5, F-6 behoben (Nacharbeit `535e0dfe`, `c2505b44`), F-4 behoben (`c2505b44`); V-1 bis V-3 (INFO) ohne Handlungsbedarf (voller `make gates` vom Hook gefahren, Schätzung 180 bis 280 lag über der Messung 111).

**Lerneintrag — neuer Sensor (geplant), plus Register:** Ein Doku-Slice, der das Gate seines Gebers zum ersten Mal gegen ein volles Docker laufen lässt, findet den Fehler des Gates: `run-container.sh` hinterließ ein Volume je Lauf (47 → 48; 669 Volumes, 169,6 GB, `test-container` rot an `initdb`). Kein Test des Gebers sah es, weil `SPEC-055` Punkt 9 Volumes nicht nannte. Der Sensor ist der Fall `p9-*-volume` mit der Mutation `rm -f` in `test-container-gegenprobe` — Kennung `slice-harness-integration-volumes` (nimmt an: §1 Ziel, DoD und §8 dort nennen `test-container` und diesen Slice; Architect-Empfehlung, Nutzer nicht entschieden). Dazu die Finding-Klassen des Reviews (Quantor „übrigen“ ohne Bezugsmenge · Beispielfolge nicht als Folge belegt · Bestandssatz widerspricht neuer Zusage · Beobachtung als Zusage lesbar · Beleg-Tabelle unvollständig): alle in `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Regel `AGENTS.md` §3.11 gilt, Zähler weit über 3×, Regel bleibt); „Plan folgt der Korrektur nicht“ in `BEO-REPO/plan-folgt-korrektur-nicht` (Regel `AGENTS.md` §3.9 bleibt). `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes` trägt die Evidenz dieses Slice schon (Zähler 2×, offen; Zuweisung ab 3×, bis dahin Verweis auf `slice-harness-integration-volumes`).

**Offen für den Nutzer:** (1) Ob das Merkmal „übrige/andere/restliche“ ohne benannte Bezugsmenge eine Zeile in Schritt 17 von `implement-slice` bekommt (Vorschlag des Planners, nicht verkörpert). (2) Ob die Lücke „Handbuch-Beispiele nur per Handlauf, nicht als Datei“ (`AGENTS.md` §3.11 verlangt die Datei) einen Sensor bekommt; heute kein Nehmer.

**Paarungen:** Anker — nichts verkörpert, kein `liegt in`; Folge-Slice — `slice-harness-integration-volumes` (Volume-Leck von `test-container`, `grep` nach `slice-v1-abschluss-container-doku` findet die Sendung dort in §1 und §8); Register — fortgeschrieben (zwei Evidence-Dateien dieses Slice, eine bereits vorhanden).

### Belege des Implementers

Stand: Basis `2e89f594`, Docker Server 29.9.0 `linux/amd64`, Probe-Läufe am 2026-10-17 nach dem Ende des vollen `make gates` (Nachweis 06:28); Messung der Marke am Ende (§4).

**Hinweis an Architect und Planner (Befund gegen den Geber `slice-v1-abschluss-container`, Ausgang dort):** Zwei Läufe von `make test-container` waren rot (`PostgreSQL nicht bereit`, je 2 m 01 s), weil der Docker-Speicher voll war (`initdb`: `No space left on device`; 668 Volumes ohne Container, 169,6 GB). Ein Container aus dem PostgreSQL-Image hinterlässt beim `docker rm -f` ein Volume (`docker create` + `docker rm -f`: 668 → 669), und `tools/harness/run-container.sh` Zeile 37 räumt ohne `-v`. Gemessen am Lauf nach dem Aufräumen: Volumes ohne Container vor `make test-container` 47, danach 48, also ein Volume je Lauf; `SPEC-055` Punkt 9 nennt Volumes seit `240562bd` (`rm -fv`), und `slice-harness-integration-volumes` übernimmt sie von `test-container` (Nehmer); das Leck bleibt, bis jener Slice liefert. Dieser Slice ändert nichts daran (Gate-Code, §1).

**Läufe (nacheinander)**

| Lauf | Ergebnis |
|---|---|
| `make binaries` (kein `dist/`) | sechs Binaries und `SHA256SUMS`, 3 s (Cache) |
| `sha256sum -c SHA256SUMS` in `dist/` | sechs Zeilen, jede `OK` |
| `./dist/pgwire-recorder_dev_linux_amd64 version` | `pgwire-recorder dev`, Exit 0 |
| `make binaries` bei vorhandenem `dist/` | `binaries: Eingabe: DIST 'dist' ist nicht leer`, `make` Exit 2, `dist/` weiter sieben Dateien |
| `make binaries VERSION=1.2.3` (`dist/` entfernt) | Namen `pgwire-recorder_1.2.3_…`, unter Windows `.exe`; `version` gibt `pgwire-recorder 1.2.3`; `sha256sum -c` `OK` |
| `docker build --target runtime --platform linux/arm64 -t pgwire-recorder:dev-arm64 .` und `docker image inspect --format '{{.Os}}/{{.Architecture}}'` | `linux/arm64`; Tag danach entfernt |
| `make build`, `docker image inspect …`, `docker run --rm pgwire-recorder:dev version` | `linux/amd64`, `pgwire-recorder dev` |
| `make test-container` (nach dem Aufräumen des Speichers, kalter Cache) | grün, 2 m 09 s: Punkte 4 bis 7 grün; danach keine Container außer `pgr-dev`, keine Netze/Tags `pgr-ct-*`, keine `pgr-container.*`; ein Volume ohne Container mehr (47 → 48). Zuvor zwei rote Läufe wegen vollen Speichers (Hinweis oben); grün auch im vollen `make gates` |
| `make docs-check` | 959 Dateien, 0 Befunde |
| `make kopf-check`, `make meldungskatalog-check`, `make abdeckung-check` | je Exit 0 |
| `make gates-schnell` | grün, 10 s |
| `dist/` entfernt | ja; Container `pgr-ct-*`, Netze und Temp-Verzeichnisse `pgr-container.*` nach `test-container` nicht vorhanden; zurück bleibt ein Volume je Lauf (siehe Hinweis) |

**Zusagen** (Satz · Probe · Beobachtung)

| Satz | Probe | Beobachtung |
|---|---|---|
| Sechs Plattformen lassen sich bauen (Linux, macOS, Windows, je `amd64`, `arm64`) | `make binaries`, `ls dist` | sechs Namen `darwin`/`linux`/`windows` × `amd64`/`arm64` |
| Name `pgwire-recorder_<VERSION>_<os>_<arch>`, unter Windows `.exe`, Vorgabe `dev` | `ls dist` mit und ohne `VERSION=1.2.3` | wie im Handbuch gezeigt |
| `version` gibt `pgwire-recorder 1.2.3` aus | `./dist/…_linux_amd64 version` nach `VERSION=1.2.3` | `pgwire-recorder 1.2.3` |
| `SHA256SUMS` eine Zeile je Binary im Format von `sha256sum`; `sha256sum -c` in `dist/`; jede Zeile `OK` | `sha256sum -c SHA256SUMS` | sechs `OK` |
| `dist/` muss leer sein oder fehlen, sonst bricht `make binaries` mit einer Meldung ab | zweiter Aufruf | Meldung, Exit 2 |
| Auf Linux mit `amd64` startet `./dist/…_linux_amd64 version` | Aufruf | `pgwire-recorder dev` |
| `make build` baut das Image für die Plattform des Docker-Servers | `docker image inspect` | `linux/amd64` |
| `docker build --target runtime --platform linux/arm64 …` baut das Image für `linux/arm64`; die Zeile `inspect` gibt `linux/arm64` aus | Aufruf | `linux/arm64` |
| Gestartet und mit Aufzeichnen und Wiedergeben geprüft ist nur das Image für `linux/amd64`, auf einem Docker-Server mit `amd64` | `make test-container`, Punkt 7 | grün, Smoke auf `linux/amd64`; „Docker 29.9.0“ steht nicht mehr im Handbuch (Beobachtung des Laufs, vom Gate nicht gepinnt) |
| Die fünf Binaries außer dem für Linux auf `amd64` und das Image für `linux/arm64` sind gebaut, Kopf und Metadaten gezeigt, nicht ausgeführt | `make test-container` Punkte 5 und 6 (grün); Handlauf `inspect` | wie behauptet; „nicht ausgeführt“ belegt kein Lauf, es ist die Grenze des Gates (`SPEC-055` Punkt 13) |
| `releasing.md`: `make test-container` baut zweimal, vergleicht Summen, prüft Kopf und Image, startet das Image der Plattform des Servers, PostgreSQL beim Wiedergeben entfernt | `make test-container` grün (Punkte 4 bis 7) | Prosa folgt dem Vertrag und dem Lauf |

**Quantoren-Durchsicht** (neue und angrenzende Sätze): „jede Zeile der Ausgabe endet auf `OK`“ — sechs Zeilen beobachtet, eng gefasst auf `sha256sum -c`; „sechs Plattformen“, „alle sechs“ — „alle sechs Binaries“ ist nicht Text geblieben (ersetzt durch „die fünf Binaries außer dem für Linux auf `amd64`“, Review F-1); „muss leer sein oder fehlen“ — beide Fälle: fehlend (erster Lauf), nicht leer (zweiter Lauf), leeres vorhandenes `dist/` (`mkdir dist`, `make binaries`): sechs Binaries und `SHA256SUMS`, sieben Dateien; „nur das Image für `linux/amd64`“ — Grenze des Gates; `releasing.md` „Beide laufen nacheinander; gleichzeitige Läufe sind nicht zugesagt“ — unverändert, deckt `SPEC-055` Punkt 6 und Grenze 5, kein Lauf des Gates; „Eine Manifestliste entsteht dabei nicht“ — unverändert, deckt `SPEC-055` Punkt 6; die zwei Bestandssätze neben den neuen Zusagen (Review F-3) sind eng gefasst: „Zum Ausführen: das Binary Ihrer Plattform aus `make binaries` oder ein Container-Laufzeitsystem“ und „Das Binary im Image ist … für die Plattform des Images, das Sie gebaut haben (bei `make build` die des Docker-Servers)“ (Probe: `make build` gibt `linux/amd64`, `--platform linux/arm64` gibt `linux/arm64`).

**Nacharbeit nach dem Review** (`docs/reviews/2026-10-17-review-slice-v1-abschluss-container-doku.md`): F-1 Bezugsmenge „fünf Binaries außer dem für Linux auf `amd64`“ in Handbuch §1 und README — Ausgang: erledigt. F-2 Summenprüfung als `( cd dist && sha256sum -c SHA256SUMS )`, danach die Befehle in Folge in einer Shell vom Wurzelverzeichnis gefahren (`make binaries`, Summenprüfung sechs `OK`, `./dist/…_linux_amd64 version` gibt `pgwire-recorder dev`; `dist/` danach entfernt) — Ausgang: erledigt. F-3 siehe oben — Ausgang: erledigt. F-4 Hinweis oben berichtigt — Ausgang: erledigt. F-5 „Docker 29.9.0“ aus dem Handbuch gestrichen — Ausgang: erledigt. F-6 Messung der Marke: `git diff --numstat 2e89f594..535e0dfe -- docs/user docs/maintainer README.md` = 5+2, 25+9, 65+5 = 111 (Marke 500, Meldeschwelle 400 nicht erreicht); Quantoren-Durchsicht ergänzt. **Akzeptierte Lücke:** die Shell-Beispiele des Handbuchs sind per Handlauf belegt, nicht als Datei, die gegen das Binary läuft — Ausgang: weiter offen, kein neuer Test in diesem Slice (Nutzerlinie: LOW ohne „Zusage ohne Test“ nicht mit Code); die Closure trägt sie ins Register (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`).

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
