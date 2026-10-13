# ADR-0040: Die Verträge der Harness-Werkzeuge stehen im Ordner harness, nicht in der Produktspezifikation

**Status:** Accepted

**Datum:** 2026-10-13

**Autor:** Architect; Entscheidung des Nutzers vom 2026-10-13 (Verlegung nach `harness/`) und vom selben Tag (Ablage, Rang, Gate-Klasse, Zeitpunkt des Umzugs)

**Bezug:** [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md), [ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0032](0032-kopf-sensor-fuer-slice-plaene.md), [ADR-0034](0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0035](0035-struktur-kennungen-im-commit-hook-abgelehnt.md), [ADR-0038](0038-gate-kette-voll-und-schnell.md); Adaption [`MR-002`](../../../harness/conventions.md#mr-002)

**Schärft:** — (Entscheidung über den Ort von Festlegungen; sie schärft keine einzelne Stelle. Der neue Ort `harness/vertraege/` entsteht erst mit dem Umzug-Slice und ist darum noch kein Ziel eines `Schärft:`; die Verträge selbst bleiben von ihren ADRs geschärft, mit neuer Adresse nach dem Umzug.)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Die Verträge der Harness-Werkzeuge (`SPEC-047` Kopf der Pläne, `SPEC-048` Abdeckung, `SPEC-049` Lint-Profil, `SPEC-050` Struktur-Kennungen im Commit-Träger, `SPEC-052` Gate-Kette) stehen seit `MR-001` und seiner Auflösung in `spec/spezifikation.md` §11, weil die Baseline für die Festlegung eines Werkzeugs das Technik-Stratum nennt: „Der Ort ist deshalb die Spezifikation, kein viertes Stratum“ (`grundlagen-referenz-richtung.md` §Spec-Straten).

Das Ergebnis im Bestand: §11 hat rund 545 Zeilen, ein Fünftel davon `SPEC-052`, und wächst mit jedem Harness-Slice. Die Produktspezifikation ist mit der Historie über 3500 Zeilen lang; wer das Produkt liest, liest Werkzeugverträge mit, und wer einen Harness-Slice arbeitet, lädt die Produktspezifikation. Die Beobachtung `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum` zeigt außerdem, dass die Stelle für Werkzeuge ohne Vorbild leicht übersehen wird; ein eigener Ordner macht sie auffindbar.

Der Nutzer hat am 2026-10-13 entschieden, die Verträge nach `harness/` zu verlegen, mit ausdrücklicher Abweichung von der Baseline als Adaption `MR-<NNN>`. Zu den offenen Einzelfragen hat er am selben Tag entschieden: (1) Ablage eine Datei je Kennung unter `harness/vertraege/SPEC-<NNN>-<kurztitel>.md`, (2) Rang 2 der Source Precedence, geteilt mit `spec/spezifikation.md`, (3) der Ordner bleibt in der Doku-Gate-Klasse `aussen`, (4) ADR und Adaption werden angenommen und der Umzug direkt hinter `slice-harness-gate-kette-stop-hook-freigabe` eingeplant. <!-- d-check:status-provenance -->

## Entscheidung

Wir legen die Verträge der Harness-Werkzeuge **in den Ordner `harness/`** und halten die Abweichung von der Baseline als [`MR-002`](../../../harness/conventions.md#mr-002) fest. Die Einzelheiten (Ort, Kennung, Rang, Referenz-Richtung, Verweise) stehen in `MR-002`, nicht hier (`AGENTS.md` §3.8).

Was diese ADR festlegt:

1. **Der Ort der Verträge ist `harness/`, nicht `spec/`.** `spec/spezifikation.md` legt nur das Produkt fest.
2. **Es bleibt das Technik-Stratum.** Gemeint ist ein Ort für dieselbe Art Festlegung mit demselben Änderungs-Prozess (fortschreibbar ohne Change Request, ADR schärft, Lastenheft ist Decke), nicht ein viertes Stratum. Die Kennung `SPEC-<NNN>` und die Überschrift jedes Vertrags bleiben, damit die Reihe nicht bricht und Anker gleich bleiben.
3. **Die Referenz-Richtung der Spezifikation gilt dort nicht** (Entscheidung des Nutzers, abweichend von der Empfehlung des Architects, siehe unten): Der Ordner ist keine Klasse der Spec-Straten; die Verträge dürfen auf ADRs und Slice-Pläne zeigen.
4. **Der Umzug geschieht direkt hinter `slice-harness-gate-kette-stop-hook-freigabe`**, vor den übrigen Harness-Slices der Welle und vor `slice-v1-abschluss-tls-client` (Entscheidung des Nutzers); bis dahin gilt `spec/spezifikation.md` §11 unverändert. Der Schnitt in Slices ist Sache des Planners. <!-- d-check:status-provenance -->
5. **Ablage und Rang** (Entscheidung des Nutzers): eine Datei je Kennung unter `harness/vertraege/`; Rang 2 der Source Precedence trägt beide Orte, die Gegenstände sind disjunkt. Einzelheiten in `MR-002`.

**Was die Baseline-Regel „kein viertes Stratum“ für das Doku-Gate bedeutet.** Die Klasse `spec-straten` in `.d-check.yml` führt drei Pfade; `harness/` fällt in die Klasse `aussen`, und `spec-straten` darf nicht auf `aussen` verweisen. Für einen Vertrag in `harness/vertraege/` gab es zwei Lesarten:

- **(a) Der Ordner gehört zur Klasse `spec-straten`** (Empfehlung des Architects, nicht gewählt). Ein Vertrag bliebe an die Decken-Regel gebunden und nennte keine ADR, keinen Slice, keine Welle, keine Adaption; kein Gate-Befund wäre nötig, weil das Gate sie fängt. Kosten: eine Zeile in `.d-check.yml` mit einem Mutanten in der Gegenprobe; ob `order` zwei Pfade derselben Stufe ausdrückt, wäre zu prüfen gewesen.
- **(b) Der Ordner bleibt `aussen`** (**Entscheidung des Nutzers, abweichend von der Empfehlung des Architects**). Der Nutzer hat sich am 2026-10-13 für (b) entschieden; die Gründe, die er nannte, sind in dieser ADR nicht festgehalten, der Architect legt keine unter.

**Folge von (b).** Der Ordner ist keine Klasse der Spec-Straten; die Richtung `no-downward` gilt für ihn nicht. Die Verträge **dürfen** damit auf ADRs und Slice-Pläne zeigen, was einem Technik-Vertrag in `spec/spezifikation.md` heute verboten ist. Die Gegenrichtung bleibt geschlossen: Lastenheft, Spezifikation und Sicht dürfen nicht auf `harness/vertraege/` verweisen (`spec-straten` gegen `aussen`); Querverweise der Produktspezifikation auf `SPEC-047` bis `SPEC-052` löst der Umzug-Slice auf.

**Die Grenze, die sich öffnet, und was dagegen steht.** Geöffnet ist die Referenz-Richtung für Festlegungen: Ein Vertrag, der Rang 2 trägt, kann sich auf eine ADR (Rang 4) oder einen Slice (Zeitschicht) stützen, und seine Geltung hinge dann an einem Text, der unter ihm steht. Dagegen stehen die Baseline-Regel „kein viertes Stratum“ (`MR-002` weist die Abweichung aus) und der Sinn der Decken-Regel: eine Ablage, deren Festlegungen niemand an die Referenz-Richtung bindet, ist die stille Setzung, die ein nicht deklariertes Stratum wäre. Die Abweichung ist damit größer als die Verlegung allein: sie betrifft Ort und Grenze.

**Setzung des Architects zum Verhältnis zu `spec-straten`: die Verträge dürfen zeigen; der Ordner ist ein Ort außerhalb der Straten, und die Aussage „Rang 2“ gilt ohne Gate.** Eine Selbstbeschränkung („die Verträge zeigen dennoch nicht auf ADRs und Slices“) wäre nur mit einer eigenen Gate-Klasse für `harness/vertraege/` prüfbar; das ist eine vierte Klasse neben dem Beschluss, den Ordner in `aussen` zu belassen, und eine Regel ohne Prüfung läse jeder Lauf als geprüft (`AGENTS.md` §3.11). Die Option mit dem geringsten Aufwand, die real trägt, ist daher die offene: keine Änderung an `.d-check.yml`, keine Regel, die nur der Wortlaut hält. Was bleibt, ist Urteil von Review und Verifikation: das Lastenheft bleibt Decke (ein Vertrag präzisiert und erweitert es nie), und ein Vertrag nennt eine ADR oder einen Slice nur als Herkunft oder Begründung, nie als Bedingung seiner Geltung. Das ist eine Konvention, kein Gate, und sie heißt so.

**Rückholbar.** Wird die Öffnung zum Problem (ein Vertrag, dessen Geltung an einer ADR hängt; Verweise auf archivierte Slices, die brechen), ist (a) der Weg zurück: der Ordner in `spec-straten`, eine Gegenprobe, und die vorhandenen Verweise werden herausgenommen. Das ist eine Folge-ADR mit `Supersedes: ADR-0040` (oder `MR-002` wird um die Klasse ergänzt, wenn nur die Gate-Klasse wechselt); der Auflösungs-Trigger in `MR-002` nennt den Anlass.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — unverändert in `spec/spezifikation.md` §11 | Baseline-Regel unverändert, keine Adaption, kein Umzug, kein Nachzug der 23 Dateien mit Anker-Verweis | Das Produkt wächst um Werkzeugverträge; die Leser beider Gegenstände laden den jeweils anderen; Wartung der 3500 Zeilen Datei wird zum Engpass |
| B — zweite Technik-Datei `spec/harness.md` | bleibt unter `spec/`, ohne Abweichung von der Baseline (die Baseline nennt „eine ihr zugeordnete Datei“ als zulässig), Gate-Klasse und Rang ohne Zusatz | eine Datei von rund 550 Zeilen und wachsend; `spec/` führt Dateien, die das Produkt nicht betreffen; der Nutzer hat diesen Ort nicht gewählt |
| C — in die Sensor-Dateien `harness/sensors/<target>.md` | Vertrag und Lesart eines Gates in einer Datei | Die Baseline trennt Festlegung und Lesart ausdrücklich (`MR-001` dokumentiert die Antwort des Kurs-Repos); `SPEC-050` und `SPEC-052` haben keine Sensor-Datei, `SPEC-052` gehört zu mehreren Targets; die Sensor-Dateien müssten Festlegungen tragen, obwohl sie nach der Baseline keine tragen |
| **D — Ordner `harness/vertraege/`, eine Datei je Kennung, Gate-Klasse `aussen` (gewählt, Entscheidung des Nutzers)** | Entscheidung des Nutzers; Produktspezifikation bleibt Produkt; Vertrag, Lesart (`sensors/`), Index (`README.md`) und Fragmente (`mk/`) liegen unter einem Dach; eine Datei je Kennung hält den Ladeaufwand eines Harness-Slice klein; Überschriften und Anker bleiben | Abweichung von der Baseline als `MR-002`, über den Ort hinaus auch in der Referenz-Richtung (die Verträge dürfen auf ADRs und Slice-Pläne zeigen, kein Gate bindet sie); Umzug mit rund 50 Dateien (Nachzug von Verweisen); der Rang braucht Änderungen an `harness/README.md` und `AGENTS.md` |
| D′ — wie D, Ordner in der Klasse `spec-straten` (Empfehlung des Architects, nicht gewählt) | Referenz-Richtung gilt für die Verträge wie bisher, ein Gate fängt Verweise auf ADR und Slice; die Abweichung bleibt der Ort | eine Zeile in `.d-check.yml` samt Gegenprobe; `order` für zwei Pfade derselben Stufe ungeklärt; der Nutzer hat sie nicht gewählt |
| E — eine Datei `harness/vertraege.md` | ein Pfad statt Ordner | wächst wie §11 heute; ein Harness-Slice lädt alle Verträge |

## Konsequenzen

- Positiv: Die Produktspezifikation legt nur das Produkt fest; ein Harness-Slice lädt die Verträge, die er braucht; die Stelle für Werkzeugverträge ist ein Ordner mit eigenem Namen.
- Negativ: Eine zweite Abweichung von der Baseline (`MR-002`), die ein Baseline-Sprung neu bewerten muss; die Verträge sind nicht mehr durch die Referenz-Richtung gebunden, und die Aussage „Rang 2“ trägt kein Gate; ein Umzug über rund 50 Dateien; die Kennung `SPEC-<NNN>` löst nicht mehr in einer Datei auf (die Reihe zählt über zwei Orte).
- **Anker und Verweise.** Die Überschriften bleiben wortgleich, die Anker darum gleich; es ändert sich der Pfad. Genannt sind die angenommenen ADRs mit einem Link der Form `spezifikation.md#spec-0…`: [ADR-0034](0034-lint-gate-mit-solid-nahem-profil.md) und [ADR-0035](0035-struktur-kennungen-im-commit-hook-abgelehnt.md); [ADR-0038](0038-gate-kette-voll-und-schnell.md) nennt `SPEC-049` und `SPEC-052` ohne Link auf die Datei, [ADR-0025](0025-benannte-slice-kennungen-im-commit-hook.md), [ADR-0028](0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0029](0029-benannte-welle-kennungen-im-commit-hook.md) und [ADR-0032](0032-kopf-sensor-fuer-slice-plaene.md) tragen `Schärft: —` und nennen keine Kennung (ihre fehlende Kante ist der Gegenstand der Beobachtung `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum`, nicht dieser ADR). Der Nachzug der beiden Links ist ein Referenz-Nachzug und keine Überschreibung (`modul-04-adrs.md` §Nachzug ist keine Überschreibung); die Entscheidung der ADRs bleibt unberührt, `AGENTS.md` §3.5 bleibt eingehalten. Ein Stub in `spec/spezifikation.md` ist **nicht** vorgesehen: Die Spezifikation darf nicht auf `harness/` verweisen (Klasse `spec-straten` gegen `aussen`; diese Grenze bleibt), und eine Überschrift ohne Verweis ersetzt keinen Link. Ortsangaben im Fließtext angenommener ADRs („§11 der Spezifikation“, [ADR-0038](0038-gate-kette-voll-und-schnell.md)) bleiben stehen und lesen sich ab dieser ADR als der neue Ort; der Index (Beziehungen) trägt diesen Hinweis.
- Akzeptiertes Negativ: Reviews und Verifikationen unter `docs/reviews/` sowie archivierte Slice-Pläne, die den alten Pfad nennen, werden nur dort nachgezogen, wo ein Link sonst bricht (Gate); Nennungen im Fließtext ohne Link bleiben stehen. Der Grund: Sie sind Lauf-Belege und tragen den Stand ihrer Zeit; die Kennung, nicht die Adresse, ist ihr Verweis.
- Folgepflicht: Ein Folge-Slice (Kennung vergibt der Planner; Zeitpunkt: direkt hinter `slice-harness-gate-kette-stop-hook-freigabe`) liefert den Umzug mit `harness/vertraege/`, dem Nachzug der Verweise, `harness/README.md` (Rang, Vertragszeilen), `AGENTS.md` §2 und §3.12 und den Command-Texten. `.d-check.yml` ändert sich nicht: Der Ordner fällt ohne Eintrag in `aussen`; der Slice belegt nur, dass das Doku-Gate bei den umgezogenen Dateien grün bleibt und die Querverweise der Produktspezifikation auf die Verträge aufgelöst sind. Bis dahin ändert sich nichts an `spec/spezifikation.md`, `AGENTS.md` und `harness/README.md`. <!-- d-check:status-provenance -->
- Folgepflicht für offene Pläne: Die Pläne in `open/` und `next/`, die §11 als Ort ihrer Randformen nennen, ziehen mit dem Umzug nach (`AGENTS.md` §3.9); bis dahin schreibt jeder Slice dort weiter und der Umzug trägt es mit.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| d-check | Keine Regel für `harness/vertraege/`: der Ordner steht in `aussen`, die Klasse `spec-straten` darf nicht auf ihn verweisen (die Regel gilt heute; ein Verweis der Spezifikation auf den Ordner ist rot). Dass die Verträge selbst nicht auf ADR und Slice zeigen, prüft kein Gate | `make docs-check` |

## Re-Evaluierungs-Trigger

Das Kurs-Repo ai-harness-course nennt einen Ort außerhalb der Produktspezifikation für Werkzeugverträge (dann löst sich `MR-002` auf), oder `harness/vertraege/` führt mehr als Verträge von Werkzeugen (dann ist der Ordner kein Ort für sie mehr), oder ein Vertrag stützt seine Geltung auf eine ADR oder einen Slice-Plan, oder ein Verweis aus einem Vertrag auf einen archivierten Slice bricht (dann wird Lesart (a) geprüft, siehe oben).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-13 | Proposed (nach Entscheidung des Nutzers vom selben Tag, Verlegung nach `harness/`) | — |
| 2026-10-13 | Accepted nach Entscheidung des Nutzers vom selben Tag: Ablage eine Datei je Kennung, Rang 2 geteilt, Gate-Klasse `aussen` (abweichend von der Empfehlung des Architects), Umzug direkt hinter `slice-harness-gate-kette-stop-hook-freigabe` | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
