# Verifikation: slice-harness-lint-werkzeug — 2026-10-06

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-lint-werkzeug.md` (Kopf, §1 bis §6, §8) gegen den Gesamt-Diff `2ee9527..8d9f0ee`. Darin: `4ab7e7b` ([`SPEC-049`](../../spec/spezifikation.md#spec-049--lint-profil-lint) Punkt 9), `2461551` (Umsetzung), `fa14758` ([Review](2026-10-06-review-slice-harness-lint-werkzeug.md), F-428 bis F-433), `f09e7c1` (Entscheidungen in `SPEC-049`) und `8d9f0ee` (Nacharbeit ohne eigenes Review).

**Eingang:**

- die DoD-Liefer-Punkte, §6 *Randformen* und *Risiken* des Plans und die Commit-Messages
- `spec/spezifikation.md` `SPEC-049` Punkt 1 bis 10 und *Grenze*; [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) (Accepted, im Diff unverändert) mit Messtabelle und Entscheidung 2, 5 und 6; [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- `.golangci.yml`, `tools/harness/lint.sh`, `Dockerfile`, `.dockerignore`, `harness/mk/lint.mk`, `Makefile`, `harness/README.md` am Stand `8d9f0ee`
- `docs/plan/planning/next/slice-harness-lint.md`, `docs/plan/planning/open/slice-lint-bestand-kern-driven.md`, `docs/plan/planning/open/slice-lint-bestand-driving.md`
- der Review-Report (F-428 bis F-433, Stand `2461551`)

Die Mutationstabelle des Implementers lag mir nicht vor (V-73). Alle Belege unten habe ich in diesem Kontext selbst erzeugt. Bei meinem Start war der Arbeitsbaum sauber, HEAD `8d9f0ee`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

**Wie die Proben gebaut wurden.** Jede Probe lief in einer frischen Kopie aus `git archive 8d9f0ee` unter `verify-lw/` im Scratchpad, ohne `cp -p`. Gebaut wurde mit `docker build --progress=plain --target lint`. Jede Mutation lief als Skript, das am Ende prüft, ob die Ersetzung gegriffen hat. „Aus dem richtigen Grund“ heißt: Die `lint:`-Zeile bzw. die Zählung, die die Zusage verlangt, steht in der Ausgabe. Bei einer Mutation des Werkzeugs heißt es: Dieselbe Verletzung bleibt ohne diese Zeile. Ein Image ist nicht entstanden, weil jeder Lauf rot war. Die Kopien sind weggeräumt.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — Profil und Stufe. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Linter nach `SPEC-049` Punkt 3 | Die Liste unter `linters.enable` ist mechanisch gleich der Liste in Punkt 3 (28 Linter), `default: none` | bestätigt |
| Schwellen nach Punkt 4 | `cyclop` 15, `funlen` 100/60, `gocognit` 20, `gocyclo` 15, `nestif` 5, `maintidx` 20, `dupl` 150, `interfacebloat` 10 | bestätigt |
| Einstellungen nach Punkt 5 | `errcheck` (fünf Funktionen), `ireturn` (Ports, `pgproto3`), `forbidigo` (drei Muster), `gomodguard_v2` (zwei Module), `revive` (25 Regeln, mechanisch gleich Punkt 5), `testpackage` `skip-regexp`, `gochecknoglobals` ohne Einstellung. Mutationen: Ohne `(net.Conn).Close` steigt `errcheck` von 14 auf 48 (`p5a`). `os.Stdout` unter `internal/` liefert `forbidigo` mit dem Text des Profils (`p5b`) | bestätigt |
| Punkt 1: `relative-path-mode: cfg`, Pfade mit `^`, `build-tags: integration` | im Profil; jeder `path` beginnt mit `^` | bestätigt (Wirkung von `cfg` nicht unterscheidbar, steht in §6 und in der DoD von `slice-harness-lint`) |
| Ausnahmen genau nach Entscheidung 2 bzw. Punkt 8, je mit `# Why:` | vier Testdatei-Regeln, `ST1005`, zehn `gochecknoglobals` je Datei und Name, `forbidigo` für `^(cmd\|test)/`. Unter `exclusions` stehen nur `warn-unused: true` und `rules`; `generated: lax` ist entfernt (F-433) | bestätigt |
| Stufe `lint` nach Punkt 2 | `Dockerfile:49`–`:54`: `--platform=$BUILDPLATFORM`, Pin `v2.14.0@sha256:ad862ba6…`, `COPY --from=deps /go/pkg/mod`, `RUN --network=none`, `GOTOOLCHAIN=local`, `GOFLAGS=-mod=readonly`. Den Digest führt außer dem `Dockerfile` nur der Plan | bestätigt |
| `.dockerignore` lässt Profil und Skript in den Kontext | `!.golangci.yml`, `!tools/harness/lint.sh` | bestätigt |
| Stufen `test`, `build`, `integration` unverändert | `git diff 2ee9527..8d9f0ee -- Dockerfile`: nur Kopfkommentar, Kommentar der Stufe `source` und die neue Stufe | bestätigt |

### Punkt 2 — Eigene Prüfungen, Ausgabe und Ausgang. **Bestätigt mit einer Lücke** (V-72); der Beleg „im Bericht genannt“ lag nicht vor (V-73).

| Zusage | Mutation | Lauf | Urteil |
|---|---|---|---|
| Punkt 6: jedes `//` der Zeile zählt | `// x //nolint:errcheck` und `// siehe nolint` in `internal/hexagon/model/doc.go:4`–`:5` (`p6`) | `lint: internal/hexagon/model/doc.go:4: Direktive nolint`, keine Zeile für `:5`, Exit 1 | bestätigt |
| ebenso, Werkzeug gebrochen | dieselbe Verletzung, Muster in `lint.sh` auf das erste Kommentarzeichen verankert (`p6w`) | keine `lint:`-Zeile mehr | Prüfung trägt |
| Punkt 7: Testfunktion in der Brücke | `export_test.go` mit `func TestX` und einem Alias (`p7`) | `lint: internal/hexagon/model/export_test.go:5: Testfunktion in der Brücke export_test.go`, der Alias ohne Zeile | bestätigt |
| Punkt 8: Regel ohne `# Why:` | `# Why:` der zweiten Regel durch `# Grund:` ersetzt (`p8b`) | `lint: .golangci.yml:168: Regel ohne Kommentarblock "# Why:" unmittelbar darüber`, weiter 103 Befunde | bestätigt |
| Punkt 8: feste Form, Einzug (F-428) | jede eingerückte Zeile +2, `# Why:` der zweiten Regel ersetzt; `config verify` nimmt das Profil an (`p8a`) | `lint: .golangci.yml:155: Form nicht erkannt`, `:157` ebenso, 103 Befunde | bestätigt |
| Punkt 8: feste Form, Flussform | `rules: [{linters: [staticcheck], text: "^ST1005"}]` (`p8c`) | `lint: .golangci.yml:157: Form nicht erkannt` | bestätigt |
| Punkt 8: Eintrag unter `rules` mit anderem Einzug | nur die Regeln um zwei Leerzeichen tiefer (`- ` auf 8), `# Why:` der zweiten Regel ersetzt (`p8d`) | **keine `lint:`-Zeile**, 103 Befunde: Die Ausnahmen wirken, die Why-Prüfung sieht nichts | **Lücke V-72** |
| Punkt 9: Profil fehlt | `.golangci.yml` gelöscht, ein `//nolint` ergänzt (`p9a`) | `lint: …doc.go:4: Direktive nolint`, `lint: .golangci.yml: fehlt`, golangci-lint läuft nicht, Exit 1 allein durch `lint:`-Zeilen | bestätigt |
| Punkt 9: Schema (`config verify`) | Tippfehler `max-same-isues` (`p9b`) | `lint: .golangci.yml: von golangci-lint abgelehnt`, dazu die `jsonschema`-Meldung, danach läuft golangci-lint (88 Befunde, weil der Tippfehler die Kürzung wieder einschaltet) | bestätigt |
| Punkt 9: Feldfolge der ungenutzten Regel | Regel mit `source`, `text`, `path-except`, `linters` in umgekehrter Folge, ohne Treffer (`p9c`). golangci-lint meldet `Text, Source, Path Except, Linters` | `lint: .golangci.yml: Regel ohne Befund: Linter: errcheck, Pfad außer: ^(cmd\|internal\|test)/, Text: ^gibtesnicht, Quelle: ^gibtesnicht`; die Warnung steht daneben | bestätigt |
| ebenso, Werkzeug gebrochen | Schlüssel `"Path Except"` zurück auf `"PathExcept"` (Stand vor F-429b, `p9cw`) | `Pfad außer` fehlt in der Zeile | Prüfung trägt |
| Punkt 9: Felder nicht erkannt | dieselbe Regel, die Feldnamen im awk verfremdet; das steht für einen geänderten Logtext (`p9e`) | `lint: .golangci.yml: Regel ohne Befund: Felder nicht erkannt` | bestätigt |
| Punkt 9: Zählregel | `uniq-by-line: true` (`q1`) | 96 statt 103 | bestätigt, Punkt 3 |

Jede Prüfung läuft auch bei einem Befund einer anderen: `p9a` zeigt die Zeile nach Punkt 6 und die nach Punkt 9 zusammen, `p8a` zeigt Form und 103 Befunde zusammen.

### Punkt 3 — Ziel und Messung. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `harness/mk/lint.mk` ohne Eintrag in `GATE_CHECKS` | `GATE_CHECKS` aus `make -pn gates`: `abdeckung-check abdeckung-gegenprobe a-check a-check-negativ baseline-verify build test docs-check hook-gegenprobe test-integration kopf-check kopf-check-gegenprobe`, ohne `lint`; `make -n gates` baut die Stufe `lint` nicht | bestätigt |
| `make gates` grün, obwohl `make lint` rot ist | Am Stand `8d9f0ee` im Arbeitsbaum: `make gates` Exit 0, kein Aufruf von golangci-lint im Log. `make lint` Exit 2 (golangci-lint 103 Befunde). Danach war der Arbeitsbaum unverändert | bestätigt |
| 103 Befunde nach der Zählregel, keine ungenutzte Regel | Grundlauf: 103 Zeilen `<pfad>:<zeile>:<spalte>: …` gleich der Summe von golangci-lint, 32 im Produkt-Code, 71 in Testdateien, keine `lint:`-Zeile | bestätigt |
| mit `uniq-by-line: true` 96, je Linter wie die Messtabelle abzüglich der Ausnahmen | `q1`: 96, davon 25 im Produkt-Code und 71 in Testdateien. Je Linter: `errcheck` 57−43=14, `gochecknoglobals` 33−10=23, `staticcheck` 14−12=2, `ireturn` 5−5=0, `revive` 13, `contextcheck` 9, `gocognit` 5, `gocyclo` 2, `containedctx`/`cyclop`/`noctx`/`unused` je 1, `testpackage` 24. Die sieben Befunde mehr bei `false` sind `cyclop` und `gocyclo` auf derselben Zeile wie ein anderer Komplexitätsbefund | bestätigt |
| Werkzeug-Zeile in `harness/README.md` sagt nur zu, was gilt | „prüft … nach `.golangci.yml` und mit eigenen Prüfungen“ (Grundlauf), „netzlos außer der Download-Stufe“ (`RUN --network=none`), „meldet alle Befunde mit Pfad“ (`--progress=plain`, ungekürzt), „endet bei einem Befund mit Fehlerstatus“ (Exit 2), „hängt nicht an `make gates`“ (oben); Bindung „kein Gate“ und [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) | bestätigt |

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| `make gates` grün | bestätigt, Abschnitt 5 |
| Review mit Report | liegt vor (F-428 bis F-433). Die Nacharbeit hat kein eigenes Review; ihre Wirkung prüft Abschnitt 2 |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen. §7 trägt Platzhalter, §6 *Risiken* vier Ausgänge „— (bei Closure)“. Das ist Sache der Closure |

---

## 2. Review-Findings

| Finding | Status | Beleg |
|---|---|---|
| F-428 (Why-Prüfung blind bei anderem Einzug) | **teilweise behoben** | Der gemeldete Fall (alles +2) und die Flussform sind jetzt `Form nicht erkannt` (`p8a`, `p8c`). Rücken nur die Einträge unter `rules` tiefer ein, bleibt die Prüfung blind (`p8d`, V-72). Dieselbe Klasse bleibt also offen |
| F-429 (a) zweites Kommentarzeichen | **behoben** | `SPEC-049` Punkt 6 entscheidet es (`f09e7c1`); `p6` und `p6w` |
| F-429 (b) Felder der ungenutzten Regel | **behoben** | `SPEC-049` Punkt 9 legt Folge und `Felder nicht erkannt` fest; `p9c`, `p9cw`, `p9e`. Der alte Schlüssel `PathExcept` traf die Warnung von golangci-lint nie, die schreibt `Path Except` |
| F-430 (Plan §3 zu `.dockerignore`) | **behoben** | §3 nennt `lint: .golangci.yml: fehlt` und dass golangci-lint nicht läuft; `p9a` zeigt genau das |
| F-431 (nicht mutierbare Zusagen ohne Artefakt) | **behoben** | Plan §6 *Eigene Prüfungen bis zum Gate ohne Gegenprobe* nennt beide Listen. Die DoD von `docs/plan/planning/next/slice-harness-lint.md` übernimmt die neuen Zusagen aus Punkt 6, 8 und 9 als Fälle und die nicht unterscheidbaren als „je mit einem Fall … oder benannt als offen im Kopf der Gegenprobe“ |
| F-432 (Logtext als einziger Halt) | **übergeben** | als nicht mutierbare Zusage in Plan §6 und in der DoD von `slice-harness-lint` |
| F-433 (`generated: lax`) | **behoben** | Zeile entfernt; *Grenze* von `SPEC-049` nennt den Default; weiter 103 Befunde |

---

## 3. Plan gegen Code

- **Kopf:** `Berührte Spec-Stellen` (Punkt 1 bis 9 und erste Hälfte von Punkt 10) und `Bezug` stimmen mit dem Diff überein. `make kopf-check` ist grün (Exit 0, auch im Gate-Lauf).
- **§1:** Ziel und Abgrenzung stimmen. Keine Datei `*.go` im Diff, kein `GATE_CHECKS +=`, ADR unverändert.
- **§3:** Jede Zeile hat ihre Datei im Diff. Die Zeile `.dockerignore` ist nachgezogen (F-430).
- **§5:** Die beiden beobachtbaren Kriterien „`make gates` grün ohne `lint`“ und „103 Befunde“ sind gemessen.
- **§6 *Randformen*:** Die Einträge zu Form, zweitem Kommentarzeichen und Feldern stehen dort. Der Satz „jede andere ist `Form nicht erkannt`“ geht über den Code hinaus (V-72).
- **§6 *Risiken*:** Die Liste der nicht mutierbaren und der neu entschiedenen Zusagen steht im Plan. Sie entspricht dem Zuwachs der DoD in `slice-harness-lint` (F-431).
- **Folge-Slices:** `slice-lint-bestand-kern-driven` mit 16 gemessen: Komplexität 7 (`toResponse`, `fromDTO`, `(Group).validate`, `(*cursor).objekte`), `revive` 8, `staticcheck` 1. `slice-lint-bestand-driving` mit 16 gemessen: Komplexität 8 (`replaySitzung`, `clientRichtung`, `startup`, `toMessage`), Kontexte 6, `revive` 2. `slice-harness-lint` nennt 32 im Produkt-Code mit acht Funktionen und 103 gesamt; die acht Funktionen sind genau die oben genannten. Die Zahlen 96 und 25 stehen dort nur als gemessene Lage der ADR.

---

## 4. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-72 | Lücke (Zusage weiter als Prüfung, F-428 nur teilweise behoben) | Die Why-Prüfung erkennt eine Regel nur an `- ` mit genau sechs Leerzeichen. Stehen `exclusions:` (2) und `rules:` (4) richtig, die Einträge darunter aber auf 7 oder mehr, greift kein Zweig des awk-Blocks. Es gibt weder `Form nicht erkannt` noch `Regel ohne Kommentarblock`. golangci-lint und `config verify` nehmen das Profil an, und die Ausnahmen wirken weiter. Eine Regel ohne `# Why:` bleibt damit still grün. Zugesagt ist das Gegenteil an drei Stellen: im Kommentar von `lint.sh` („ein Eintrag unter `rules` mit anderem Einzug ergeben `Form nicht erkannt`“), im Kopf von `.golangci.yml` („eine andere Form meldet `make lint` als `Form nicht erkannt`“) und in Plan §6 *Form des Profils* („jede andere ist `Form nicht erkannt`“). `SPEC-049` Punkt 8 nennt die sechs Leerzeichen als feste Form, zählt als Befund aber nur `exclusions` und `rules` auf. Ob ein tieferer Eintrag Befund ist, entscheidet der Architect; danach ziehen Code oder die drei Kommentare nach (`AGENTS.md` §3.11, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) | `tools/harness/lint.sh:54`–`:56`, `:74`–`:83`; `.golangci.yml:8`–`:11`; Plan §6 *Form des Profils* | `p8d`: 103 Befunde, keine `lint:`-Zeile |
| V-73 | Beleg fehlt im Eingang | Liefer-Punkt 2 verlangt, dass jede Mutation „im Bericht als Zusage · Mutation · roter Lauf genannt“ ist. Dieser Bericht lag weder dem Review noch mir vor, auch nicht für die Nacharbeit `8d9f0ee`. Abschnitt 1 ersetzt ihn durch eigene Läufe. Ob der Implementer V-72 gesehen hat, ist damit nicht belegt | DoD Punkt 2 | — |
| V-74 | Hinweis an `slice-harness-lint` | Eine Regel, die golangci-lint erst beim Laden ablehnt, nimmt `config verify` an: `path` zusammen mit `path-except` oder nur ein Feld (`at least 2 of (text, source, path[-except], linters)`). Der Lauf endet dann rot mit der Meldung von golangci-lint, aber ohne `lint:`-Zeile. Das ist nicht falsch grün und nach Punkt 9 zulässig; die Gegenprobe sollte den Fall aber kennen, damit „Profil vom Schema abgelehnt“ nicht als vollständige Profilprüfung gelesen wird | `tools/harness/lint.sh:106`–`:118` | Probe mit `path` und `path-except` sowie Probe nur mit `path`: je Exit 1, `can't load config` |

---

## 5. Gate-Ergebnis

`make gates` am Stand `8d9f0ee` im Arbeitsbaum: **grün** (Exit 0). Gelaufen sind: `abdeckung-check`, `abdeckung-gegenprobe`, `a-check` (0 Befunde), `a-check-negativ`, `baseline-verify` (v6.13.0, 54 Dateien), `build`, `test` (gofmt, vet, Unit-Tests), `docs-check` (269 Dateien, 0 Befunde), `hook-gegenprobe`, `test-integration`, `kopf-check`, `kopf-check-gegenprobe`. `make lint` am selben Stand: rot, 103 Befunde, wie es Punkt 10 für das Werkzeug vorsieht.

**Urteil:** Die drei Liefer-Punkte sind erfüllt. Ausnahme ist die Lücke V-72 in Punkt 2: F-428 ist nur für den gemeldeten Fall behoben, nicht für die Klasse. Vor der Closure braucht V-72 eine Entscheidung des Architects und danach Code oder Kommentar, die ihr folgen.
