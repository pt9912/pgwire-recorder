# MR-001 — Spezifikations-Ort für Verträge der Harness-Werkzeuge

Regeln dieser Datei: Pflichtfelder sind Datum, Geltungsbereich,
**Ersetzt-Baseline-Regel**, Adaption, Begründung und Auflösungs-Trigger;
`Löst auf` und `Ausgelöst durch Baseline-Stand` nur, wenn dieser Eintrag einen
früheren ablöst. `Ersetzt-Baseline-Regel` nennt **genau eine** Regel der
Baseline, an deren Stelle dieser Eintrag tritt — als Link mit
Abschnitts-Anker in die vendored Fassung; ein Datei-Link benennt keine Regel.
**Dieser Link trägt zwei Dinge, die sich bewegen, und für beide gibt es einen
Wächter statt einer Formregel:** die *Tiefe* (nach dem `git mv` nach `done/`
zeigt der relative Pfad eine Ebene zu hoch — der Umzug zieht die
Pfad-Berichtigung nach sich, als eigener Commit; die Existenzprüfung des Links
meldet sie, wenn sie ausbleibt) und die *Version* (jeder Baseline-Bump
entwertet `<tag>` — der adoptierte Stand steht einmal im Adaptions-Block, ein
Versions-Sensor prüft jeden Pin dagegen; Muster in `.d-check.yml`).
Wer mehrere Regeln ersetzen will, schreibt mehrere Einträge. Ein Eintrag, der
keine benannte Regel ersetzt, ist ein **Fork**, keine Adaption.

- **Datum:** 2026-10-06
- **Geltungsbereich:** Harness-Werkzeuge und Gates unter `tools/` und `harness/mk/` samt ihren Gegenproben; nicht das Produkt.
- **Ersetzt-Baseline-Regel:** [`grundlagen-referenz-richtung.md` §Spec-Straten](../../../.harness/baseline/v6.13.0/regelwerk/grundlagen-referenz-richtung.md#spec-straten-mehr-als-ein-spec-dokument) — das Technik-Stratum als einziger zulässiger Ort einer selbst gesetzten technischen Festlegung. Die Spec-Straten der Baseline (dort und in `modul-03-spec.md` §Spec-Stratifizierung) kennen nur das Produkt. `AGENTS.md` §3.12 verlangt die Entscheidung einer Randform „in der Spezifikation oder einer ADR“, und eine ADR ist nach `Accepted` gesperrt (`modul-04-adrs.md` §Hard Rule für Accepted-ADRs); für einen Werkzeugvertrag bleibt damit kein Ort.
- **Adaption:** Für Harness-Werkzeuge gilt als Spezifikations-Ort im Sinn von `AGENTS.md` §3.12 der ZUSAGE-Kopf des Skripts zusammen mit der Vertragszeile in `harness/README.md` §Sensors. Die ADR des Gates trägt Entscheidung und Gründe (`AGENTS.md` §3.8). Eine Randform, die nach Annahme der ADR entschieden wird, steht im ZUSAGE-Kopf. §6 des Slice-Plans zeigt auf diesen Ort. Jede dort zugesagte Randform braucht einen Fall in der Gegenprobe (`AGENTS.md` §3.10, §3.11).
- **Begründung:** `BEO-REPO/harness-lesart-ohne-entscheidungsort`, zwei Auftreten: `slice-harness-kopf-sensor` (Review F-380, Verifikation V-42; im Register belegt) und `slice-lastenheft-pruefbarkeit` (Review F-406; Beleg folgt mit dessen Closure). Dazu ein Change Request an das Kurs-Repo ai-harness-course, den der Nutzer stellt (Entscheidung des Nutzers vom 2026-10-06).
- **Auflösungs-Trigger:** Das Kurs-Repo entscheidet den Change Request: Eine neue Baseline mit einem Ort für Werkzeugverträge wird adoptiert, oder der Change Request wird abgelehnt und das Repo entscheidet neu.
- **Aufgelöst:** 2026-10-06. Das Kurs-Repo ai-harness-course hat den Change Request beantwortet: Der Ort eines Werkzeugvertrags ist das Technik-Stratum, ein Abschnitt der Spezifikation oder eine ihr zugeordnete Datei; dort werden Randformen auch nach Annahme der ADR fortgeschrieben, ohne Folge-ADR. Die Lesart eines Laufs steht in einer Sensor-Datei `harness/sensors/<target>.md` nach der Baseline-Vorlage, die den Vertrag verlinkt; womit das Werkzeug selbst geprüft ist, bleibt in Skriptkopf und Tests. Damit ersetzt dieser Eintrag keine Baseline-Regel mehr. Den Abschnitt legt `slice-harness-vertraege-spezifikation` an.
