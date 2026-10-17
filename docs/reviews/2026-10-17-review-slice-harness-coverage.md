# Review-Report: slice-harness-coverage — 2026-10-17

**Review-Art:** Code — geprüft gegen Plan (`slice-harness-coverage`), `ADR-0044` (Accepted), `SPEC-056`, `SPEC-052` Punkt 13 und `AGENTS.md` §3.9 bis §3.14.

**Gegenstand:** Diff `9ee6387d..201b9fe5` (Basis der Größenmetrik bis letzter Commit des Slice), nur die Teile von `slice-harness-coverage`.

**Skill:** `.harness/skills/reviewer.md` @ 201b9fe5 · 
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

> **Zitier-Form.** Kennung statt Adresse; Baseline-Stellen als Tag + Pfad (`v6.18.0` · `regelwerk/<datei>.md` §<Abschnitt>). Gelaufen sind nur Lesen, `grep`, `git diff --numstat`; kein `make gates`, kein Docker. Mutationen sind gedanklich gegen den Code gelesen, nicht gefahren.

**Eingangs-Kontext:**

- Slice-Plan `slice-harness-coverage` (§1 bis §8, §7 Belege des Implementers)
- `ADR-0044`, `ADR-0026`, `ADR-0033`; `SPEC-056`, `SPEC-052` Punkt 13
- `LH-QA-07` (Messmethode 2)
- `AGENTS.md` (Hard Rules, §3.9 bis §3.14)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Nachgemessen mit der Metrik aus §4 (`git diff --numstat 9ee6387d..201b9fe5 -- Dockerfile .dockerignore harness/mk tools/harness tools/test harness/vertraege harness/README.md docs/user docs/plan/adr`, Summe Einfügungen plus Löschungen): 903, nicht 873. Das liegt über der Marke 900, bei der §4 die Rückführung `in-progress` → `next` „ohne Ermessen" vorsieht; §7 nennt die Größe nicht. | `AGENTS.md` §3.14; Plan §4 Rückführungen | `docs/plan/planning/in-progress/slice-harness-coverage.md` · „**Marke 900:** darüber gilt die Rückführung unten, ohne Ermessen" | ja — `git diff --numstat` mit den Pfaden aus §4 | Größenmarke überschritten, Rückführung nicht gezogen |
| F-2 | MEDIUM | `SPEC-056` Punkt 1 sagt zu, die Zeile stehe „in der Quelle ohne Leerzeile davor". `build.mk` liest mit `sed -n 's/^\*\*Schwelle:\*\* ...'` jede Zeile am Zeilenanfang, gleich was davor steht; kein Fall der Gegenprobe (`p1-zeile-*`) setzt eine Leerzeile davor. Die Zusage fängt keine Mutation. | `AGENTS.md` §3.10 (Ortsangabe je Schritt), §3.11 | `harness/vertraege/SPEC-056-testabdeckung-schwelle.md` · „in der Quelle ohne Leerzeile davor" | ja — Fall in `coverage-gegenprobe` mit Leerzeile vor der Zeile | Zusage ohne Mutation (Ortsangabe in Punkt 1) |
| F-3 | LOW | Der Kopf der Gegenprobe führt unter OFFEN: „Die Reihenfolge, in der mehrere Gründe einer Eingabe gemeldet werden, ist nicht zugesagt; jeder Fall hat genau einen Grund". `SPEC-056` Punkt 6 sagt die Reihenfolge zu, und `p6-r-*` prüft sie mit zwei Gründen je Eingabe. Der Kommentar beschreibt den Stand vor dem Nachzug. | `AGENTS.md` §3.7, §3.11 | `tools/harness/coverage-gegenprobe.sh` · „ist nicht\n#     zugesagt; jeder Fall hat genau einen Grund" | nein — Lesen | Kommentar widerspricht Vertrag nach Nachzug |
| F-4 | LOW | Der Plan nennt ADR-0044 und die Schwelle weiter eine Empfehlung: Kopf „(Proposed, Empfehlung des Architect vom 2026-10-17, bis der Nutzer entschieden hat)", §1 „Empfehlung des Architect: 94 Prozent", §3 „Die Schwelle selbst ist eine Empfehlung, bis der Nutzer die ADR annimmt". Die ADR ist seit `6b4ecbd9` Accepted. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-harness-coverage.md` · „Proposed, Empfehlung des Architect" | nein — Lesen | Plan folgt Entscheidung nicht |
| F-5 | INFO | `ADR-0044` ist Accepted und trägt im Kopf weiter „Autor: Architect; Empfehlung, bis der Nutzer entschieden hat". Eine Accepted-ADR wird nicht überschrieben (§3.5); nur zur Kenntnis. Ebenso steht „94" in der Entscheidung der ADR und in Punkt 1 (Grund) von `SPEC-056`, während „die Schwelle steht an einer Stelle" gilt für Code und Gate (`Dockerfile`, `build.mk`, Prüfschritt führen keine Zahl). | `AGENTS.md` §3.5 | `docs/plan/adr/0044-schwelle-der-anweisungsabdeckung.md` · „Empfehlung, bis der Nutzer entschieden hat" | nein | Zustandsangabe in Accepted-ADR veraltet |
| F-6 | INFO | `p6-a-nicht-lesbar` und `p6-r-a-b` werden mit Hinweis auf stderr übersprungen, wenn der Aufrufer gesperrte Dateien liest (root); der Lauf bleibt grün. Dann läuft Punkt 6 (a) „nicht lesbar" ohne Fall. | `AGENTS.md` §3.10 | `tools/harness/coverage-gegenprobe.sh` · „p6-a-nicht-lesbar übersprungen" | ja — Lauf als root | Fall nach Umgebung übersprungen |
| F-7 | INFO | Der Bereich `9ee6387d..201b9fe5` enthält Commits fremder Slices (`slice-harness-devcontainer-*`: `docs/developer/devcontainer.md`, `open/slice-harness-devcontainer-einrichten.md`). Sie liegen außerhalb der Metrikpfade und wurden nicht geprüft; Zuständigkeit Orchestrator/andere Reviews. | Maintainability | `git log 9ee6387d..201b9fe5` · „slice-harness-devcontainer-einrichten" | ja — `git log` | Diff-Bereich mit fremden Commits |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `Dockerfile` Stufe `test` (gofmt, `go vet -tags integration`, Tests, Prüfschritt in eigener Schicht, `ARG` ohne Vorgabewert nach der Testschicht) | geprüft, ohne Befund; gofmt, vet und Tests bleiben, jeder Schritt bricht die Stufe; Schwelle nur als Build-Argument, kein `--target test` ohne Argument sonst im Repo (grep) |
| `.dockerignore`, `harness/vertraege/SPEC-052-gate-kette.md` Punkt 13, `harness/mk/gates-schnell.mk`, `tools/harness/gates-schnell-gegenprobe.sh` | geprüft, ohne Befund; Allowlist, Ausnahmeliste und `LANGSAM` stimmen mit `coverage.mk` überein |
| Schwelle nur in `SPEC-056` Punkt 1 | geprüft, ohne Befund; Wert kommt per `$(shell sed)` in `build.mk`; Rezept ohne `$(MAKE)` und `+`, `.PHONY: test` gesetzt |
| `tools/harness/coverage.sh` gegen Punkt 3 bis 7 | geprüft, ohne Befund; Reihenfolge (a), (b), (c), (g), (h), dann (d), (e), (f) stimmt mit dem Vertrag; `>=` in ganzen Zahlen, Anzeige abgerundet, Sortierung nach Importpfad |
| `tools/harness/coverage-gegenprobe.sh` gegen Punkt 10: Fälle je Element (4 a bis c, 6 a bis h, Paare der Reihenfolge, 9 je Flag und `MAKEFLAGS`) | geprüft, ohne Befund; einzelne Zusagen ohne Fall sind im Kopf unter OFFEN genannt (außer F-2, F-3) |
| `harness/README.md` §Sensors (`make test`, `make coverage-gegenprobe`, „Nicht behauptet") | geprüft, ohne Befund; sagt nur zu, was die Gegenprobe zeigt |
| `docs/user/abdeckung-gate.md`, `abdeckung-gesamt.md`, `abdeckung-vollstaendig.md`, `harness/vertraege/README.md` | geprüft, ohne Befund; Teil 2 von 3 mit Deklaration im Kopf der Gegenprobe |
| Core-Reinheit, Meldungscodes, Handbuch | nicht berührt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Größenmarke überschritten, Rückführung nicht gezogen · Zusage ohne Mutation (Ortsangabe in Punkt 1) · Kommentar widerspricht Vertrag nach Nachzug · Plan folgt Entscheidung nicht

## Verdikt

**Merge-blockierend:** ja — F-1 (Rückführung nach §4 ohne Ermessen; Entscheidung beim Planner) und F-2 (Zusage ohne Mutation). F-3 und F-4 sind LOW und blockieren nicht.

**Übergabe:** Findings gehen an den Implementer; F-1 zusätzlich an den Planner. F-2 verlangt entweder einen Fall mit Mutation oder eine engere Fassung der Zusage (Entscheidung des Architect). Die Finding-Klassen gehen in die Slice-Closure §7. DoD-Konformität prüft der Verifier.
