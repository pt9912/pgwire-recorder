# Verifikations-Bericht: slice-harness-coverage — 2026-10-17

**Rolle:** Verifier (Modul 11), frischer Kontext. **Frage:** Bauen wir es richtig, gemessen an Plan, DoD und Entscheidungen?

**Gegenstand:** `slice-harness-coverage`, Diff `9ee6387d..201b9fe5` (Code 7f1f6ead, 1655a6b8, 9e741444, 507a26e2, 201b9fe5; Architect 9f56eeec, 6b4ecbd9).

**Gegen:** Plan §1 bis §8; [`SPEC-056`](../../harness/vertraege/SPEC-056-testabdeckung-schwelle.md); [`[ADR-0044](../plan/adr/0044-schwelle-der-anweisungsabdeckung.md)`](../plan/adr/0044-schwelle-der-anweisungsabdeckung.md) (Accepted); `SPEC-052` Punkt 13; `LH-QA-07` Messmethode 2; Abnahmeszenario 17.

> **Zitier-Form.** Dieser Bericht friert ein: Kennung, nicht Adresse.

**Grenze dieser Verifikation.** Auf Wunsch des Auftraggebers lief kein Docker, kein `make gates`. Selbst gefahren: der Prüfschritt und der Teil der Gegenprobe ohne Docker (bash), `make abdeckung-check`, `make kopf-check`, `make -n test`, Mutationen auf Kopien unter dem Temp-Pfad. Nicht selbst gefahren: die Fälle mit Docker (`p1-schwelle-kommt-an`, `p1-zeile-fehlt-stufe`, `p2-paket-ohne-test`, `p8-*`), `make test`, `make gates`.

## Ergebnis

**Konform bis auf zwei Punkte, die offen bleiben: kein Blocker im Code.** Befunde: 1 MEDIUM (Nachweis `make gates` am Endstand nicht bestätigt), 1 MEDIUM (Größenmarke 900 überschritten, Wert fehlt in §7), 2 LOW.

## 1. DoD gegen den Diff

| DoD-Punkt | Befund |
|---|---|
| Gate mit Schwelle (Stufe `test`, `coverage.sh`, Schwelle aus SPEC-056 Punkt 1, Build-Argument ohne Vorgabe, `gofmt`/`vet`/Tests unverändert) | **bestätigt, Laufbeleg des Docker-Teils nicht selbst gesehen.** `Dockerfile`: `ARG COVERAGE_SCHWELLE` ohne Vorgabe, `--network=none`, gofmt, `go vet -tags integration`, `go test -covermode=set -coverprofile` ohne `-coverpkg`, danach `sh tools/harness/coverage.sh`. `make -n test` zeigt `--build-arg COVERAGE_SCHWELLE="94"`. `.dockerignore` hat `!tools/harness/coverage.sh`. Messung 4201/4448 steht in §7 und [ADR-0044](../plan/adr/0044-schwelle-der-anweisungsabdeckung.md). |
| Gegenprobe `coverage-gegenprobe` an `GATE_CHECKS` und `GATE_SCHNELL_AUSNAHME`, Fälle und Mutationen je Punkt 4, 6, 10, Abdeckungs-Deklaration | **bestätigt für die Fälle ohne Docker, für die mit Docker nur durch Lesen.** Ohne Docker: 110 Prüfungen wie erwartet (selbst gefahren). Anschluss: `coverage.mk` (`GATE_CHECKS +=`), `gates-schnell.mk`, `LANGSAM` in `gates-schnell-gegenprobe.sh`. Deklaration `LH-QA-07/Messung-2-von-3` im Kopf. `make abdeckung-check` Exit 0, `make kopf-check` Exit 0; `abdeckung-gesamt.md` führt `1/3 Gate; 2/3 Gate; 3/3 Gate`, `vollständig`; `abdeckung-vollstaendig.md` führt `LH-QA-07`. Abnahmeszenario 17 damit nachweisbar. |
| Doku `harness/README.md` | **bestätigt.** `make test` nennt die Schwelle mit `SPEC-056` und [ADR-0044](../plan/adr/0044-schwelle-der-anweisungsabdeckung.md); Zeile `make coverage-gegenprobe` mit Bindung; „Nicht behauptet“ grenzt gegen die Abdeckung je Anforderung ab. |
| `make gates` grün | **nicht bestätigt** (F-1). |
| Review, Report | nicht Gegenstand dieser Verifikation; ein Report zu diesem Slice lag in `docs/reviews/` nicht vor, als ich nachsah (Befund nur zur Kenntnis: Reihenfolge der Rolle liegt beim Auftraggeber). |
| Closure-Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen | erwartungsgemäß offen: §7 trägt noch Platzhalter `<…>`, §6 alle Risiken `Ausgang: —`. |

## 2. SPEC-056, zwölf Punkte gegen `coverage.sh` und die Dateien

| Punkt | Befund |
|---|---|
| 1 Schwelle, Zeile am Zeilenanfang, `make test` liest, kein Wert im `Dockerfile` | erfüllt (`sed` mit `^\*\*Schwelle:\*\*`; Fälle mitten im Text, eingerückt, doppelt, fehlend laufen grün) |
| 2 Messung `set`, ohne `-coverpkg`, nur Unit-Tests | erfüllt im `Dockerfile`; Wirkung nur im Docker-Fall `p2-paket-ohne-test` |
| 3 Gesamtwert nach Anweisungen | erfüllt (awk summiert `$2`) |
| 4 (a) bis (d) | erfüllt; (d) hat laut Kopf bewusst keinen Fall |
| 5 `>=` ganzzahlig, Anzeige abgerundet | erfüllt (`ausgefuehrt*100 >= schwelle*gesamt`) |
| 6 (a) bis (h) und Reihenfolge | erfüllt: Skript prüft (a), (b), (c), (g), (h), danach Zeilen (d)/(e) nach Zeilenfolge, (f) zuletzt |
| 7 Ausgabe | erfüllt (Liste sortiert, Gesamtzeile zuletzt, Differenz nur bei rot, Meldungen nur auf stderr) |
| 8 Reihenfolge der Stufe, Allowlist | erfüllt im `Dockerfile` und `.dockerignore` |
| 9 Rezept ohne `$(MAKE)` | erfüllt (Einzeiler); `-n`, `-q`, `-t`, `MAKEFLAGS` laufen ohne Docker grün |
| 10 Gegenprobe | erfüllt, siehe Abschnitt 3 |
| 11 Ausnahmeliste, Beleg SPEC-052 Punkt 13 | erfüllt (`p11-*` grün; Gegenprobe baut aus Kopie ohne `docs/user/` und `README.md`) |
| 12 Grenzen | Prosa, kein Code |

## 3. Bewusstes Brechen (Stichprobe, selbst gefahren)

Treiber: Kopie des Baums unter dem Temp-Pfad, Gegenprobe ohne Docker darauf gerichtet.

| Zusage | Mutation | Ergebnis |
|---|---|---|
| Vergleich genau an der Schwelle (Punkt 5) | `>` statt `>=` | rot: `p5-genau`, `p6-h-grenze-100`, `p4-anzahl-null`, `p3-zaehler-ist-nicht-anzahl` |
| Paket ohne Test zählt (Punkt 4a) | Blöcke unter `/cmd/` übersprungen | rot: `p4-main-ohne-test` (nur dieser) |
| Schwelle kommt als Build-Argument an (Punkt 1, Attrappen-Teil) | `--build-arg` aus `build.mk` entfernt | rot: `p1-wert-94`, `p1-wert-aus-datei`, `p1-zeile-*`, `p9-ohne-flag` |
| Ausnahmeliste (Punkt 11) | `coverage-gegenprobe` aus `GATE_SCHNELL_AUSNAHME` | rot: `p11-ausnahmeliste` (nur dieser) |
| Reihenfolge (b) vor (c) | Prüfung (c) vor (b) | rot: `p6-r-b-c`, dazu `p6-b-leer` (zweiter Fall folgt aus der Mutation, kein Fremdfall) |
| Eingabeform Schwelle (g)/(h) | Prüfung (g) entfernt | rot: `p6-g-fehlt`, `p6-g-leer`, `p6-r-g-h`, `p6-r-h-d` und `p6-h-*` (Mutation war breiter als eine Zusage; Richtung stimmt) |

Basis ohne Mutation: 110 Prüfungen grün. Die Tests wären ohne den Fix aus dem richtigen Grund rot. Ungefahren: dass die Schwelle in der Stufe wirklich ankommt (Docker-Fall `p1-schwelle-kommt-an`) und dass ein Mutant im `Dockerfile` (`ARG … =1`) ihn rot färbt; das belegt nur die Tabelle des Implementers in §7, nicht ich.

## 4. Befunde

- **F-1 (MEDIUM) `make gates` am Endstand nicht bestätigt.** `.harness/state/gates-passed.diffsha` stammt von 19:51; die Commits 507a26e2 (19:52) und 201b9fe5 (19:54) ändern `harness/mk/build.mk`, `tools/harness/coverage-gegenprobe.sh` und `SPEC-056`. `working-tree-hash.sh` liefert heute `48b5674b…`, der Stand lautet `d18f7665…`. Die Abweichung kann auch von den nicht verfolgten Dateien unter `docs/plan/planning/` kommen, aber gesichert ist der grüne Lauf am Endstand nicht. §7 sagt „`make gates` einmal am Endstand (siehe Bericht)“, nennt aber keinen Beleg. Verlangt: ein `make gates` am Endstand.
- **F-2 (MEDIUM) Größenmarke 900.** Die Metrik aus §4 (`git diff --numstat 9ee6387d..201b9fe5`, genannte Pfade) ergibt **903**; 603 davon die Gegenprobe. Die Rückführung bei „über 900“ gilt „ohne Ermessen“. Der Wert fehlt in §7 (dort steht nur der Platzhalter „wird mit der Übergabe genannt“). Die Summe enthält Architect-Commits (SPEC-056, [ADR-0044](../plan/adr/0044-schwelle-der-anweisungsabdeckung.md), SPEC-052, Indizes); Entscheidung zur Zählung liegt beim Planner. Kein Code-Mangel.
- **F-3 (LOW) Kopf der Gegenprobe widerspricht dem Vertrag.** Der Abschnitt OFFEN sagt: „Die Reihenfolge, in der mehrere Gründe einer Eingabe gemeldet werden, ist nicht zugesagt; jeder Fall hat genau einen Grund.“ SPEC-056 Punkt 6 und die Fälle `p6-r-*` sagen das Gegenteil (`AGENTS.md` §3.11, Kommentar sagt Falsches zu).
- **F-4 (LOW) Stale Text.** [ADR-0044](../plan/adr/0044-schwelle-der-anweisungsabdeckung.md) trägt `Status: Accepted`, aber `Autor: Architect; Empfehlung, bis der Nutzer entschieden hat`; der Plan Kopf nennt die ADR noch „(Proposed, …)“ und §7 trägt einen alten Block „Rückgaben an den Architect“, den §7 selbst als beantwortet kennzeichnet. Die ADR ist nach `Accepted` unveränderlich (`AGENTS.md` §3.5); Korrektur im Plan.

Außerhalb des Gegenstands: Der Diff `9ee6387d..201b9fe5` enthält auch Commits eines anderen Slice (`slice-harness-devcontainer-*`, `docs/developer/devcontainer.md`); sie zählen nicht für diese Verifikation, verfälschen aber die Metrik nicht (andere Pfade, bis auf `docs/developer`, das nicht in der Metrik liegt).

## 5. Offen für die Closure

F-1 und F-2 klären, F-3 und F-4 nachziehen, danach Closure-Notiz, Risiko-Ausgänge, Register, Lerneintrag.
