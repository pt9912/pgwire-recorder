# Review-Report: slice-harness-lint-werkzeug — 2026-10-06

**Review-Art:** Harness-Werkzeug, Profil und Planung, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff 2ee9527..2461551`, also `4ab7e7b` (Entscheidung des Architects in `SPEC-049` Punkt 9: fehlendes Profil, Schema, ungenutzte Regel, Zählregel) und `2461551` (`.golangci.yml`, Stufe `lint` im `Dockerfile`, `.dockerignore`, `tools/harness/lint.sh`, `harness/mk/lint.mk`, Werkzeug-Zeile in `harness/README.md`, Plan und drei Folge-Slices). Schwerpunkte laut Auftrag: Profil gegen `SPEC-049` und das Vorbild, eigene Prüfungen in `lint.sh`, die Stufe und die Allowlist, kein Anschluss an `GATE_CHECKS`, Stichprobe der Mutationen, Plan und Folge-Slices. Den Wechsel von 96 auf 103 Befunde in DoD und §5 hat der Koordinator als Folge der Zählregel `uniq-by-line: false` bestätigt; er wird hier nur nachgemessen, nicht gewertet.

**Skill:** `.harness/skills/reviewer.md` @ `2461551`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-harness-lint-werkzeug.md` (ganz, Stand `2461551`), die Diffs von `docs/plan/planning/next/slice-harness-lint.md`, `docs/plan/planning/open/slice-lint-bestand-kern-driven.md` und `docs/plan/planning/open/slice-lint-bestand-driving.md`, dazu die DoD von `slice-harness-lint`
- `spec/spezifikation.md` [`SPEC-049`](../../spec/spezifikation.md#spec-049--lint-profil-lint) Punkt 1 bis 10 und Grenze; [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) (Accepted), [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md)
- Vorbild `/Development/KI/ai-harness-init/.golangci.yml`
- `Dockerfile`, `.dockerignore`, `Makefile`, `harness/mk/build.mk`, `harness/mk/lint.mk`, `harness/README.md` §Sensors
- `AGENTS.md` §3.2, §3.3, §3.5, §3.6, §3.8 bis §3.12
- Vorheriger Report: [Review zu `slice-harness-vertraege-spezifikation`](2026-10-06-review-slice-harness-vertraege-spezifikation.md); höchste vergebene Nummer vor diesem Lauf F-427

**Ausgeführte Läufe:** Alle Läufe fanden in frischen Kopien statt, je aus `git archive 2461551` entpackt (neue mtime, kein `cp -p`), unter dem Scratch-Pfad des Auftrags. Gebaut wurde mit `docker build --progress=plain --target lint`, ohne Tag. Der Grundlauf ergab Exit 1 mit 103 Befunden: 32 im Produkt-Code, 71 in Testdateien, keine `lint:`-Zeile. Je Paket und Linter stimmt das mit den Folge-Slices überein (Punkt 6 unten). Danach liefen 23 Mutationen der Stufe `lint`. Zwei davon (Leerzeile *vor* dem Why-Block; Einzug verdoppelt, daran scheiterte das YAML) waren falsch angesetzt und wurden durch `w2b` und `w4b` ersetzt. Dazu kam ein Lauf der Stufe `test` in einer Kopie mit geändertem Profil; er war grün. Im Image: bash 5.2.37, go1.27.1. `make gates` lief nicht. Die Kopien sind weggeräumt. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-428 | MEDIUM | Die Prüfung nach `SPEC-049` Punkt 8 (Regel ohne `# Why:`) erkennt den Abschnitt `rules` nur, wenn `exclusions` mit genau zwei Leerzeichen unter `linters` steht (`pfad[2] == "exclusions"`). In der Kopie bekam jede eingerückte Zeile des Profils zwei Leerzeichen mehr (`exclusions` auf 4), und über der zweiten Regel stand `# Grund:` statt `# Why:`. `config verify` nahm das Profil an, golangci-lint lief mit denselben 103 Befunden, und eine `lint:`-Zeile gab es nicht. Die Prüfung blieb also falsch grün (`w4b`). Mit Einzug 2 ist dieselbe Mutation rot (`w1`). Dasselbe gilt für eine Regel-Liste in Flow-Form. Der Kopf von `.golangci.yml` und der Kommentar zu (8) in `lint.sh` sagen die Meldung ohne Bedingung zu. Die Einzug-Annahme steht nur als Kommentar im awk-Block; Grenze von `SPEC-049` und Kopf von `lint.sh` nennen sie nicht. Nach dem Anschluss an das Gate macht ein Umformatieren des Profils die Prüfung blind, ohne dass ein Lauf es zeigt. | `SPEC-049` Punkt 8; `AGENTS.md` §3.10, §3.11; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Reviewer-Skill MEDIUM „unklare Fehlerbehandlung am Rand des Spec-Bereichs“ | `tools/harness/lint.sh:47`–`:48`, `:68`–`:75`; `.golangci.yml:4`–`:8`; `spec/spezifikation.md:2038`–`:2040` | ja (Mutation `w4b`: Einzug +2, `# Why:` entfernt, Stufe ohne `lint:`-Zeile) | Prüfung hängt an der Formatierung statt an der Struktur |
| F-429 | LOW | Zwei Randformen entscheidet der Code, ohne dass `SPEC-049` oder Plan §6 sie nennen. (a) **Zweites Kommentarzeichen im Kommentar:** `// x //nolint:errcheck` wird gemeldet (`n1`, Zeile 237), denn das Muster greift am zweiten `//`. Punkt 6 sagt dagegen: „Steht vor `nolint` im Kommentar ein anderes Wort, ist das kein Befund.“ Beide Lesarten passen auf den Wortlaut. golangci-lint selbst wertet die Direktive dort nicht aus. (b) **Felder der ungenutzten Regel:** Punkt 9 verlangt „Linter, Pfad und Text der Regel“. Das Skript schreibt zusätzlich `Pfad außer:` und `Quelle:`. Führt eine Regel keines dieser fünf Felder, bleibt die Zeile nach `Regel ohne Befund: ` leer. Beides ist streng statt locker und macht nichts falsch grün. | `AGENTS.md` §3.12; `SPEC-049` Punkt 6 und 9 | `tools/harness/lint.sh:30`, `:122`–`:127`; `spec/spezifikation.md:2024`–`:2027`, `:2064`–`:2066` | ja (Mutation `n1`) | Randform im Code entschieden |
| F-430 | LOW | Plan §3 begründet die Zeile `.dockerignore` mit „ohne sie läuft das Image mit Default-Profil“. Nach `SPEC-049` Punkt 9 und `lint.sh` stimmt das nicht: Ohne `.golangci.yml` im Kontext schreibt die Stufe `lint: .golangci.yml: fehlt`, und golangci-lint läuft nicht (`p1`). Der Plan gibt den Stand vor der Entscheidung `4ab7e7b` wieder. | `AGENTS.md` §3.9; `BEO-REPO/plan-folgt-korrektur-nicht` | `docs/plan/planning/in-progress/slice-harness-lint-werkzeug.md:122`; `tools/harness/lint.sh:86`–`:89` | ja (Mutation `p1`) | Plan folgt Korrektur nicht |
| F-431 | LOW | Laut Auftrag sind die Zusagen, die sich nicht mutieren lassen, „als offen für `slice-harness-lint` benannt“. Im Repo finde ich das nicht. Plan §6 (*Eigene Prüfungen bis zum Gate ohne Gegenprobe*) und §7 nennen keine einzelne Zusage. Die DoD von `slice-harness-lint` verlangt allgemein „je Zusage aus `SPEC-049`, die eine Mutation fangen kann“. Ihre Aufzählung enthält keinen der Punkte aus `4ab7e7b` (Profil fehlt, Schema, Zeilenform der ungenutzten Regel, Ausgang bei einer `lint:`-Zeile allein), auch keinen, der nicht mutierbar wäre. Die Liste steht damit nur im Bericht des Implementers. Der ist Lauf-Beleg, und spätere Läufe lesen ihn nicht. Der Folge-Slice hat also keine benannte Sendung. | `AGENTS.md` §3.10, §3.11; Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice, Klasse 1 | `docs/plan/planning/in-progress/slice-harness-lint-werkzeug.md:205`–`:210`, `:233`–`:238`; `docs/plan/planning/next/slice-harness-lint.md:128`–`:137` | ja (Lesen) | Übergabe an Folge-Slice ohne Artefakt |
| F-432 | INFO | Ob eine Regel ungenutzt ist, erkennt `lint.sh` nur am Logtext `[runner/exclusion_rules] Skipped 0 issues by rules: [`. golangci-lint ändert seinen Exit-Status bei dieser Warnung nicht. Ändert eine neue Version den Text, steht die Warnung zwar noch auf stderr, aber es gibt keine `lint:`-Zeile und keinen Exit 1. Der Re-Evaluierungs-Trigger der ADR (Anhebung von golangci-lint) deckt das ab, und `slice-harness-lint` führt „eine ungenutzte Regel“ als Fall der Gegenprobe. Hinweis für jenen Slice: Dieser Fall ist der einzige Halt der Kopplung. | [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) §Re-Evaluierungs-Trigger | `tools/harness/lint.sh:104`–`:110`, `:129`–`:133` | ja (Mutation `u1` heute rot) | — (Hinweis) |
| F-433 | INFO | `exclusions.generated: lax` ist eine Ausnahme, die weder `SPEC-049` Punkt 5 noch Punkt 8 nennt. Sie stammt aus dem Vorbild und ist auch der Default von golangci-lint v2. Im Repo trägt keine Datei `*.go` die Markierung `Code generated`. Die Einstellung ändert also nichts. Ob Punkt 8 („nur als Einstellung nach Punkt 5 oder als Regel“) sie abdecken soll, entscheidet der Architect. | `SPEC-049` Punkt 8; Rolle Architect | `.golangci.yml:153`; `spec/spezifikation.md:2037`–`:2038` | nein | — (Hinweis) |

## Antwort auf die Schwerpunkte

1. **`.golangci.yml` gegen `SPEC-049` und das Vorbild.**
   - Die Linter-Menge entspricht genau Punkt 3 (`default: none`, 28 Einträge), die Schwellen genau Punkt 4.
   - Die Einstellungen entsprechen Punkt 5: `errcheck` mit `fmt.Fprint*` und `Close` an `net.Conn` und `net.Listener`, `ireturn` mit Ports und `pgproto3`, `forbidigo`, Erlaubnisliste in `gomodguard_v2`, die 25 Regeln von `revive`, `skip-regexp` von `testpackage`. Dazu kommen `relative-path-mode: cfg`, `build-tags: integration` und die ungekürzte Ausgabe nach Punkt 1 und 9.
   - Die Regeln unter `exclusions` sind genau die aus Punkt 8 bzw. Entscheidung 2 der ADR: vier für Testdateien, `ST1005`, zehn `gochecknoglobals` je Datei und Name (die ADR zählt zehn), `forbidigo` für `cmd/` und `test/`. Jede hat `# Why:`, und keine blendet nur Bestand aus. Dass die zehn Werte nach der Initialisierung nur gelesen werden, zeigt ein `grep` nach Zuweisungen: Außer der Deklaration gibt es keine.
   - Gegenüber dem Vorbild fehlen `testpackage` für `cmd/` und `gochecknoglobals` für die Fassung. Die ADR begründet beides: `cmd/` hat keine Tests, `version` lässt der Linter zu. Die Sperrliste ist zur Erlaubnisliste geworden (Entscheidung 3). Die erweiterten `forbidigo`-Muster stehen in Punkt 5. Die Testregel `unused-parameter` hat jetzt ein `Why:`, das ihr im Vorbild fehlte.
   - Wirkung in Mutationen geprüft: Ohne das Ports-Muster gibt es drei zusätzliche `ireturn` (`i1`). Ein eigenes Interface im Modell liefert `ireturn` (`i2`). Ein Import von `github.com/jackc/pgpassfile` (indirekt in `go.mod`) wird von `gomodguard_v2` gemeldet (`g1`). `fmt.Println` und `os.Stdout` unter `internal/` liefern `forbidigo` (`f1`). `Close` an einer Datei liefert `errcheck` (`c1`). `uniq-by-line: true` ergibt 96 (`q1`). Ohne `build-tags` sind es 83 (`t1`).
   - Ohne Befund bis auf F-433.
2. **`lint.sh`.**
   - `//nolint`: `/* NoLint */`, `/// NOLINT begruendung`, `//<Tab>nolint`, eine Zeile mit CR LF und ein String-Literal werden gemeldet (`n1`, `n3`). Nicht gemeldet werden `// see nolint`, `//nolintx`, `// nolint_foo` und `/** nolint */` (`n2`); das entspricht Punkt 6. Eine Datei unter einem Pfad mit Leerzeichen wird mit ihrem Pfad gemeldet (`s1`). Zum zweiten `//` siehe F-429 (a).
   - Brücke: `func TestX` und `func Benchmarkfoo` in `export_test.go` werden gemeldet, der Alias nicht (`e1`).
   - Why-Block: Gemeldet werden ein fehlendes `Why:` (`w1`), eine Leerzeile zwischen Block und Regel (`w2b`), eine andere erste Zeile im Block (`w3`) und ein Kommentar hinter `exclusions:` bei fehlendem `Why:` (`w5`). Zum Einzug siehe F-428.
   - `config verify`: Bei einem unbekannten Schlüssel kommen die `lint:`-Zeile, die Meldung des Schemas und danach trotzdem der Lauf von golangci-lint (`v1`), ohne Netz. Ein leeres Profil wird abgelehnt (`l1`).
   - Fehlendes Profil: Es kommt `fehlt`, golangci-lint läuft nicht, `//nolint` wird weiter gemeldet (`p1`).
   - Ungenutzte Regel: genau eine Zeile im Format von Punkt 9, Exit 1 (`u1`); zum Format siehe F-429 (b).
   - Ausgang: in jedem Lauf 1, im Grundlauf durch golangci-lint.
   - Robustheit: `set -uo pipefail` ohne `-e` ist Absicht (jede Prüfung läuft). Die Pfade sind gequotet und werden zeilenweise gelesen. `LC_ALL=C` steht da. `wait` wartet in bash 5.2 auf die Prozess-Substitution von `tee`. Dass ein Exit ungleich 0 ohne Befund als rot zählt, passt zu Punkt 2 („scheitert die Stufe“).
   - Ein Fall, der falsch rot wäre, ist nicht aufgetreten. Falsch grün war nur F-428.
3. **Dockerfile-Stufe und Allowlist.**
   - Die Stufe hat Digest-Pin, `$BUILDPLATFORM`, `RUN --network=none`, `CGO_ENABLED=0`, `GOTOOLCHAIN=local` und `GOFLAGS=-mod=readonly`; die Module kommen per `COPY --from=deps`. Der Pin steht im Code nur im `Dockerfile`.
   - Zur Allowlist: Die Stufen `test`, `build` und `integration` lesen nur `*.go` bzw. die Pakete (`gofmt -l`, `go vet ./...`, `go test ./...`, `go build ./cmd/...`), und `tools/harness/` enthält keine Datei `*.go`. In der Kopie mit geändertem Profil lief die Stufe `test` grün durch.
   - Die Aussage des Implementers trifft zu: Es ändert sich nur der Cache, weil eine Änderung an Profil oder Skript die Stufe `source` neu baut. Der Kommentar in `.dockerignore` verspricht nur für Doku, Spec und den übrigen Harness, dass der Cache hält. Das stimmt.
4. **`lint.mk` und Werkzeug-Zeile.** Kein `GATE_CHECKS +=`. `$(DOCKER_BUILD)` aus `build.mk` trägt `--progress=plain`, also kommt die Ausgabe ungekürzt; damit trägt „meldet alle Befunde mit Pfad“. Die README-Zeile sagt nichts zu, was der Lauf nicht zeigt: Die Stufe schreibt nichts in den Arbeitsbaum, und es gibt kein Tag und keinen Output.
5. **Mutationen.** Die Tabelle des Implementers lag mir nicht vor, nur der Hinweis auf sie. Die 23 eigenen Mutationen (Kürzel oben) wurden alle rot, wo eine Zusage es verlangt, mit einer Ausnahme: F-428. Dass die Zusagen ohne Mutation als offen benannt seien, ist im Repo nicht belegt (F-431). Ob die Tabelle die DoD trägt, prüft der Verifier.
6. **Plan und Folge-Slices.**
   - Nachgemessen am Stand `2461551`: `slice-lint-bestand-kern-driven` hat 16 (Komplexität 7 in vier Funktionen, `revive` 8, `staticcheck` 1), `slice-lint-bestand-driving` 16 (Komplexität 8 in vier Funktionen, Kontexte 6, `revive` 2). `slice-harness-lint` hat 32 im Produkt-Code mit acht Funktionen und 71 in Testdateien. Die Zuordnung der Linter zu den Funktionen in beiden Bereinigungs-Slices stimmt mit der Ausgabe überein.
   - Die Zahlen 96 und 25 stehen nur noch dort, wo sie die gemessene ADR-Lage oder den Grund einer schon eingetretenen Rückführung wiedergeben (`slice-harness-lint.md:195`, Drift-Log der Roadmap). Das ist richtig so.
   - Kopf, §1, §5 und §6 des eigenen Plans folgen dem Code. §3 tut es an einer Stelle nicht (F-430).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `spec/spezifikation.md` `SPEC-049` Punkt 9, Historie (`4ab7e7b`) | geprüft, ohne Befund. Keine ADR, kein Slice, keine Welle, kein Commit genannt. Die vier Randformen der Rückgabe sind entschieden, bevor der Code-Commit kam. |
| `docs/plan/adr/` — Hard Rule 3.5, 3.8 | geprüft, ohne Befund. Keine ADR im Diff. Die Zählregel steht in der Spezifikation, nicht in der ADR. |
| Hard Rule 3.3 | geprüft, ohne Befund. Kein Move im Diff; der Move `94c25c6` liegt davor und ist rein. |
| Hard Rule 3.2, 3.6 | geprüft, ohne Befund. Keine Ausnahme ohne `Why:`, keine Schwelle gesenkt, kein `//nolint` im Bestand. |
| Hard Rule 3.12 | geprüft. Befund F-429. |
| `.golangci.yml` | geprüft. Befund F-433; Kopfkommentar siehe F-428. |
| `tools/harness/lint.sh` | geprüft. Befund F-428, F-429, F-432. |
| `Dockerfile`, `.dockerignore` | geprüft, ohne Befund. |
| `harness/mk/lint.mk`, `harness/README.md` | geprüft, ohne Befund. |
| `docs/plan/planning/in-progress/` | geprüft. Befund F-430, F-431. |
| `docs/plan/planning/next/`, `docs/plan/planning/open/` | geprüft. Zählungen stimmen; Befund F-431 (`slice-harness-lint`). |

**Summary:** 1 MEDIUM (F-428: Why-Prüfung falsch grün bei anderem Einzug), 3 LOW (F-429 Randform im Code entschieden, F-430 Plan folgt Korrektur nicht, F-431 Übergabe an Folge-Slice ohne Artefakt), 2 INFO; keine HIGH. Wiederkehrende Klassen: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-428), `BEO-REPO/plan-folgt-korrektur-nicht` (F-430).
