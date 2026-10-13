# ADR-0040: Die Verträge der Harness-Werkzeuge stehen im Ordner harness, nicht in der Produktspezifikation

**Status:** Proposed

**Datum:** 2026-10-13

**Autor:** Architect; die Verlegung nach `harness/` ist die Entscheidung des Nutzers vom 2026-10-13, Ort im Einzelnen, Rang und Gate-Klasse sind Empfehlung des Architect bis zur Freigabe

**Bezug:** [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md), [ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0032](0032-kopf-sensor-fuer-slice-plaene.md), [ADR-0034](0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0035](0035-struktur-kennungen-im-commit-hook-abgelehnt.md), [ADR-0038](0038-gate-kette-voll-und-schnell.md); Adaption [`MR-002`](../../../harness/conventions.md#mr-002)

**Schärft:** — (Entscheidung über den Ort von Festlegungen; sie schärft keine einzelne Stelle. Die Verträge selbst bleiben von ihren ADRs geschärft, mit neuer Adresse nach dem Umzug.)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Verträge der Harness-Werkzeuge (`SPEC-047` Kopf der Pläne, `SPEC-048` Abdeckung, `SPEC-049` Lint-Profil, `SPEC-050` Struktur-Kennungen im Commit-Träger, `SPEC-052` Gate-Kette) stehen seit `MR-001` und seiner Auflösung in `spec/spezifikation.md` §11, weil die Baseline für die Festlegung eines Werkzeugs das Technik-Stratum nennt: „Der Ort ist deshalb die Spezifikation, kein viertes Stratum“ (`grundlagen-referenz-richtung.md` §Spec-Straten).

Das Ergebnis im Bestand: §11 hat rund 545 Zeilen, ein Fünftel davon `SPEC-052`, und wächst mit jedem Harness-Slice. Die Produktspezifikation ist mit der Historie über 3500 Zeilen lang; wer das Produkt liest, liest Werkzeugverträge mit, und wer einen Harness-Slice arbeitet, lädt die Produktspezifikation. Die Beobachtung `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum` zeigt außerdem, dass die Stelle für Werkzeuge ohne Vorbild leicht übersehen wird; ein eigener Ordner macht sie auffindbar.

Der Nutzer hat am 2026-10-13 entschieden, die Verträge nach `harness/` zu verlegen, mit ausdrücklicher Abweichung von der Baseline als Adaption `MR-<NNN>`.

## Entscheidung

Wir legen die Verträge der Harness-Werkzeuge **in den Ordner `harness/`** und halten die Abweichung von der Baseline als [`MR-002`](../../../harness/conventions.md#mr-002) fest. Die Einzelheiten (Ort, Kennung, Rang, Referenz-Richtung, Verweise) stehen in `MR-002`, nicht hier (`AGENTS.md` §3.8).

Was diese ADR festlegt:

1. **Der Ort der Verträge ist `harness/`, nicht `spec/`.** `spec/spezifikation.md` legt nur das Produkt fest.
2. **Es bleibt das Technik-Stratum.** Gemeint ist ein Ort für dieselbe Art Festlegung mit demselben Änderungs-Prozess (fortschreibbar ohne Change Request, ADR schärft, Lastenheft ist Decke), nicht ein viertes Stratum. Die Kennung `SPEC-<NNN>` und die Überschrift jedes Vertrags bleiben, damit die Reihe nicht bricht und Anker gleich bleiben.
3. **Die Referenz-Richtung der Spezifikation gilt auch dort** (Empfehlung, siehe unten): Ein Vertrag nennt keine ADR, keinen Slice, keine Welle und keine Adaption.
4. **Der Umzug geschieht nach dem zweiten der beiden Hook-Slices der Gate-Kette (Freigabe bei inerten Pfaden)**, weil beide `SPEC-052` ändern; bis dahin gilt `spec/spezifikation.md` §11 unverändert. Der Schnitt in Slices ist Sache des Planners.

**Was die Baseline-Regel „kein viertes Stratum“ für das Doku-Gate bedeutet.** Die Klasse `spec-straten` in `.d-check.yml` führt heute drei Pfade; `harness/` fällt in die Klasse `aussen`, und `spec-straten` darf nicht auf `aussen` verweisen. Für einen Vertrag in `harness/vertraege/` gibt es zwei Lesarten:

- **(a) Der Ordner gehört zur Klasse `spec-straten`** (Empfehlung des Architect). Ein Vertrag bleibt an die Decken-Regel gebunden; die Entscheidung des Nutzers verlegt die Datei, nicht den Änderungs-Prozess. Kosten: keine, denn die Verträge erfüllen die Regel heute schon (sie stehen in der Klasse); die Klasse bekommt einen Pfad dazu. Ob `order` zwei Pfade derselben Stufe ausdrückt, prüft der Umzug-Slice; kann es das nicht, steht der Ordner in `order` hinter `spec/spezifikation.md`, und die Spezifikation darf dann nicht auf ihn verweisen (sie tut es heute nicht).
- **(b) Der Ordner bleibt `aussen`.** Dann wäre ein Vertrag frei, auf ADRs und Slices zu zeigen, und die Zusage „die Grenze gilt dort so wie in der Spezifikation“ fiele weg. Das ist der Fall, den die Baseline mit „kein viertes Stratum“ ausschließt: eine Ablage, deren Festlegungen niemand an die Referenz-Richtung bindet, ist dieselbe stille Setzung, die ein nicht deklariertes Stratum wäre.

**Empfehlung: (a).** Gemessen am Aufwand ist (a) eine Zeile in `.d-check.yml` mit einem Mutanten in der Gegenprobe, (b) spart nichts und öffnet eine Grenze, die heute geschlossen ist. Die Entscheidung gehört dem Nutzer; bis er sie trifft, heißt sie hier Empfehlung.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — unverändert in `spec/spezifikation.md` §11 | Baseline-Regel unverändert, keine Adaption, kein Umzug, kein Nachzug der 23 Dateien mit Anker-Verweis | Das Produkt wächst um Werkzeugverträge; die Leser beider Gegenstände laden den jeweils anderen; Wartung der 3500 Zeilen Datei wird zum Engpass |
| B — zweite Technik-Datei `spec/harness.md` | bleibt unter `spec/`, ohne Abweichung von der Baseline (die Baseline nennt „eine ihr zugeordnete Datei“ als zulässig), Gate-Klasse und Rang ohne Zusatz | eine Datei von rund 550 Zeilen und wachsend; `spec/` führt Dateien, die das Produkt nicht betreffen; der Nutzer hat diesen Ort nicht gewählt |
| C — in die Sensor-Dateien `harness/sensors/<target>.md` | Vertrag und Lesart eines Gates in einer Datei | Die Baseline trennt Festlegung und Lesart ausdrücklich (`MR-001` dokumentiert die Antwort des Kurs-Repos); `SPEC-050` und `SPEC-052` haben keine Sensor-Datei, `SPEC-052` gehört zu mehreren Targets; die Sensor-Dateien müssten Festlegungen tragen, obwohl sie nach der Baseline keine tragen |
| **D — Ordner `harness/vertraege/`, eine Datei je Kennung (gewählt als Ort `harness/`, Einzelform Empfehlung)** | Entscheidung des Nutzers; Produktspezifikation bleibt Produkt; Vertrag, Lesart (`sensors/`), Index (`README.md`) und Fragmente (`mk/`) liegen unter einem Dach; eine Datei je Kennung hält den Ladeaufwand eines Harness-Slice klein; Überschriften und Anker bleiben | Abweichung von der Baseline als `MR-002`; Umzug mit rund 50 Dateien (Nachzug von Verweisen); Rang und Gate-Klasse brauchen Änderungen an `harness/README.md`, `AGENTS.md` und `.d-check.yml` |
| E — eine Datei `harness/vertraege.md` | ein Pfad statt Ordner | wächst wie §11 heute; ein Harness-Slice lädt alle Verträge |

## Konsequenzen

- Positiv: Die Produktspezifikation legt nur das Produkt fest; ein Harness-Slice lädt die Verträge, die er braucht; die Stelle für Werkzeugverträge ist ein Ordner mit eigenem Namen.
- Negativ: Eine zweite Abweichung von der Baseline (`MR-002`), die ein Baseline-Sprung neu bewerten muss; ein Umzug über rund 50 Dateien; die Kennung `SPEC-<NNN>` löst nicht mehr in einer Datei auf (die Reihe zählt über zwei Orte).
- **Anker und Verweise.** Die Überschriften bleiben wortgleich, die Anker darum gleich; es ändert sich der Pfad. Genannt sind die angenommenen ADRs mit einem Link der Form `spezifikation.md#spec-0…`: [ADR-0034](0034-lint-gate-mit-solid-nahem-profil.md) und [ADR-0035](0035-struktur-kennungen-im-commit-hook-abgelehnt.md); [ADR-0038](0038-gate-kette-voll-und-schnell.md) nennt `SPEC-049` und `SPEC-052` ohne Link auf die Datei, [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md), [ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0029](0029-benannte-welle-kennungen-im-commit-hook.md) und [ADR-0032](0032-kopf-sensor-fuer-slice-plaene.md) tragen `Schärft: —` und nennen keine Kennung (ihre fehlende Kante ist der Gegenstand der Beobachtung `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum`, nicht dieser ADR). Der Nachzug der beiden Links ist ein Referenz-Nachzug und keine Überschreibung (`modul-04-adrs.md` §Nachzug ist keine Überschreibung); die Entscheidung der ADRs bleibt unberührt, `AGENTS.md` §3.5 bleibt eingehalten. Ein Stub in `spec/spezifikation.md` ist **nicht** vorgesehen: Die Spezifikation darf nicht auf `harness/` verweisen (Klasse `spec-straten` gegen `aussen`), und eine Überschrift ohne Verweis ersetzt keinen Link. Ortsangaben im Fließtext angenommener ADRs („§11 der Spezifikation“, [ADR-0038](0038-gate-kette-voll-und-schnell.md)) bleiben stehen und lesen sich ab dieser ADR als der neue Ort; der Index (Beziehungen) trägt diesen Hinweis.
- Akzeptiertes Negativ: Reviews und Verifikationen unter `docs/reviews/` sowie archivierte Slice-Pläne, die den alten Pfad nennen, werden nur dort nachgezogen, wo ein Link sonst bricht (Gate); Nennungen im Fließtext ohne Link bleiben stehen. Der Grund: Sie sind Lauf-Belege und tragen den Stand ihrer Zeit; die Kennung, nicht die Adresse, ist ihr Verweis.
- Folgepflicht: Ein Folge-Slice (Kennung vergibt der Planner; Schnitt und Zeitpunkt im Bericht des Architect) liefert den Umzug mit `harness/vertraege/`, dem Nachzug der Verweise, der Gate-Klasse samt Gegenprobe, `harness/README.md` (Rang, Vertragszeilen), `AGENTS.md` §2 und §3.12 und den Command-Texten. Bis dahin ändert sich nichts an `spec/spezifikation.md`, `AGENTS.md`, `harness/README.md` und `.d-check.yml`.
- Folgepflicht für offene Pläne: Die Pläne in `open/` und `next/`, die §11 als Ort ihrer Randformen nennen, ziehen mit dem Umzug nach (`AGENTS.md` §3.9); bis dahin schreibt jeder Slice dort weiter und der Umzug trägt es mit.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| d-check | Klasse `spec-straten` nimmt `harness/vertraege/**` auf; ein Vertrag mit Verweis auf eine ADR, einen Slice oder ein Artefakt außerhalb der Straten ist rot (Empfehlung (a); die Gegenprobe liefert der Umzug-Slice) | `make docs-check` |

## Re-Evaluierungs-Trigger

Das Kurs-Repo ai-harness-course nennt einen Ort außerhalb der Produktspezifikation für Werkzeugverträge (dann löst sich `MR-002` auf), oder `harness/vertraege/` führt mehr als Verträge von Werkzeugen (dann ist der Ordner kein Ort für sie mehr).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-13 | Proposed (nach Entscheidung des Nutzers vom selben Tag, Verlegung nach `harness/`) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
