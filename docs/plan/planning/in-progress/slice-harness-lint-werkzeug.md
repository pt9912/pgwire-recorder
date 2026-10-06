# Slice slice-harness-lint-werkzeug: Werkzeug-Ziel `make lint` vor dem Gate

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Er trägt zum
Abnahmeszenario 17 (M3) bei, das mit dem letzten Slice der Reihe nachweisbar wird.
Eingesammelt wird er von der nächsten Welle-Closure. Angelegt nach Entscheidung des
Nutzers vom 2026-10-06 (Bereinigung vor dem Gate, ohne Stufen) als erster Slice der
Bereinigungs-Reihe; Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode 1; das Werkzeug misst, das Gate liefert `slice-harness-lint`). Bindung an Entscheidungen: [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Profil, dauerhafte Ausnahmen, Entscheidung 6: Werkzeug-Ziel vor dem Gate), [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Stufe des Multistage-`Dockerfile`, Image per Digest gepinnt, Build-Kontext als Allowlist).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--lint-profil-lint) (Punkte 1 bis 9 und die erste Hälfte von Punkt 10: Werkzeug ohne Gate)

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

**Ziel:** `make lint` gibt es als Werkzeug ohne Gate (`SPEC-049` Punkt 10, erste
Hälfte): Es prüft den Go-Code des Moduls mit golangci-lint `v2.14.0` als Stufe `lint`
des `Dockerfile` nach dem Profil `.golangci.yml` (`SPEC-049` Punkt 1 bis 5 und 8, mit
den dauerhaften Ausnahmen aus Entscheidung 2 in
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md), je mit `# Why:`) und
mit den drei eigenen Prüfungen (`//nolint`, Test in der Brücke, Regel ohne `Why:`,
Punkt 6 bis 8), meldet den ganzen Bestand ungekürzt mit Pfad und endet nach Punkt 9.
Es hängt nicht an `GATE_CHECKS`. Jeder folgende Slice der Reihe misst damit sein
Ergebnis unter seinen Pfaden, mit demselben Ziel, das `slice-harness-lint` später an
die Gate-Kette hängt (Entscheidung 6).

**Herkunft:** Entscheidung des Nutzers vom 2026-10-06 nach der Messung des Bestands:
Bereinigung vor dem Gate, ohne Stufen. Den Vorschlag des Architect, das Werkzeug mit
`slice-harness-blackbox-kern` zu liefern, hat der Planner geprüft und verworfen (§8,
*Schnitt*): Mit dem Werkzeug hätte jener Slice drei Liefer-Punkte in drei Schichten
(Harness und `Dockerfile`, Tests des Domain Model, Tests der Services) und zwei
Vorgänge, Arbeit am Werkzeug und Arbeit am Gegenstand, in einem Review.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Anschluss an `GATE_CHECKS`, die Gegenprobe `make lint-gegenprobe` und die Doku
  in `AGENTS.md` §3.2 und `harness/README.md` §Sensors — übernimmt `slice-harness-lint`,
  wenn der Bestand grün ist; vorher machte das Ziel `make gates` rot (103 Befunde).
- Die Bereinigung eines Befunds — Bestand bleibt bewusst stehen: Die Testdateien
  bereinigen die vier Umstellungs-Slices, den Produkt-Code `slice-lint-bestand-kern-driven`
  und `slice-lint-bestand-driving`. Dieser Slice ändert keine Datei `*.go`.
- Eine Ausnahme, die nur Bestand aussetzt, oder eine Stufe — gibt es nicht
  (Entscheidung 5); was der Lauf meldet, bleibt als Befund stehen.
- Eine Änderung an `SPEC-049` oder an der ADR — Schicht-Abgrenzung: Der Vertrag ist
  vor dem Code geschrieben; weicht das Werkzeug davon ab, geht der Befund an den
  Architect (§4).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Profil und Stufe: `.golangci.yml` nach `SPEC-049` Punkt 1 bis 5 und 8, mit genau
      den dauerhaften Ausnahmen aus Entscheidung 2 der ADR, je mit einem Kommentarblock
      `# Why:` unmittelbar darüber; Stufe `lint` im `Dockerfile` nach Punkt 2
      (`golangci/golangci-lint:v2.14.0` per Digest, `$BUILDPLATFORM`, Module aus `deps`,
      `--network=none`, `GOTOOLCHAIN=local`, `GOFLAGS=-mod=readonly`); `.dockerignore`
      lässt `.golangci.yml` und, falls die eigenen Prüfungen als Skript liegen, dieses
      in den Build-Kontext. Die Stufen `test`, `build` und `integration` bleiben
      unverändert. — bestätigt an `8d9f0ee` und `17daf26` (Verifikation, Punkt 1 und
      Nachtrag): Linter und `revive`-Regeln mechanisch gleich `SPEC-049`, Schwellen und
      Ausnahmen nach Punkt 4, 5 und 8, `git diff` des `Dockerfile` berührt die übrigen
      Stufen nicht. Rot ohne die Einstellung: ohne `(net.Conn).Close` steigt `errcheck`
      von 14 auf 48 (`p5a`), `os.Stdout` unter `internal/` liefert `forbidigo` (`p5b`,
      Review `f1`), ohne das Ports-Muster drei `ireturn` mehr (Review `i1`), ein indirektes
      Modul liefert `gomodguard_v2` (Review `g1`), ohne `build-tags` 83 statt 103 (Review
      `t1`). Ohne Mutation, die der Lauf unterscheidet (`relative-path-mode: cfg`, Pfade
      mit `^`, `--network=none`, `GOFLAGS`, `GOTOOLCHAIN`, Pin und Plattform): übergeben an
      die Gegenprobe von `slice-harness-lint` (§6, §7).
- [x] Eigene Prüfungen, Ausgabe und Ausgang: die Prüfungen nach `SPEC-049` Punkt 6
      (`//nolint`), 7 (Funktion `Test…`, `Benchmark…`, `Example…`, `Fuzz…` in
      `export_test.go`) und 8 (Regel ohne `# Why:`) laufen in der Stufe, jede auch bei
      einem Befund einer anderen, mit Zeilen `lint: <pfad>:<zeile>: <befund>`; dazu nach
      Punkt 9 fehlendes Profil, Schema-Prüfung (`config verify`) und ungenutzte Regel
      (`warn-unused`) mit Zeilen `lint: .golangci.yml: <befund>` und der Ausgang nach
      Punkt 9. Je Zusage ist
      die Mutation einmal von Hand rot gesehen und im Bericht als Zusage · Mutation ·
      roter Lauf genannt (`AGENTS.md` §3.10); die bleibende Gegenprobe liefert
      `slice-harness-lint`. — bestätigt an `17daf26` (Verifikation, Punkt 2 und Nachtrag),
      je Zusage rot aus dem richtigen Grund: Punkt 6 `p6`, Werkzeug gebrochen `p6w` ohne
      Zeile; Punkt 7 `p7`; Punkt 8 `p8b`, `p8a`, `p8c`, `p8d`, `p8f`, `p8g`; Punkt 9 `p9a`,
      `p9b`, `p9c`, Werkzeug gebrochen `p9cw` ohne `Pfad außer`, `p9e`. Rot ohne den Fix
      von V-72: `p8d` am Stand `8d9f0ee` ohne `lint:`-Zeile, an `17daf26` mit
      `Form nicht erkannt` (Tabelle des Implementers in §7). Der Bericht des Implementers
      mit Zusage · Mutation · roter Lauf lag Review und Verifikation für `2461551` und
      `8d9f0ee` nicht vor (V-73); der Beleg je Zusage sind dort die Läufe der Verifikation.
- [x] Ziel und Messung: `harness/mk/lint.mk` führt `make lint` ohne Eintrag in
      `GATE_CHECKS`; `harness/README.md` nennt es in der Tabelle der Werkzeuge mit
      „kein Gate“ und Bindung an die ADR. Ein Lauf am Stand des Slice reproduziert die
      Messung der ADR: 103 Befunde nach den dauerhaften Ausnahmen (`SPEC-049` Punkt 9,
      `uniq-by-line: false`; mit `uniq-by-line: true` 96, je Linter und Paket wie die
      Messtabelle abzüglich der Ausnahmen), keine ungenutzte Regel; eine
      Abweichung ist im Bericht erklärt oder geht an den Architect (§4). — bestätigt an
      `8d9f0ee` und `17daf26` (Verifikation, Punkt 3 und Nachtrag): `GATE_CHECKS` aus
      `make -pn gates` ohne `lint`, Grundlauf 103 Befunde (32 Produkt-Code, 71
      Testdateien), keine `lint:`-Zeile; mit `uniq-by-line: true` 96, je Linter wie die
      Messtabelle abzüglich der Ausnahmen (`q1`).
- [x] `make gates` grün. — an `8d9f0ee` (Verifikation, Abschnitt 5) und an `17daf26`
      (Nachtrag, alle 12 Gates); die Closure ändert nur Planungsdokumente, Register und
      einen Satz in `.claude/commands/implement-slice.md` Schritt 19 (`make docs-check`,
      `make kopf-check` grün).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8). — Review bis
      `2461551` (F-428 bis F-433); die Nacharbeit `8d9f0ee` und `17daf26` hat kein eigenes
      Review, die Verifikation hat jedes Finding nachgeprüft (Abschnitt 2 und Nachtrag).
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
| `.golangci.yml` | neu | Profil nach `SPEC-049` Punkt 1 bis 5 und 8 (`relative-path-mode: cfg`, Pfade mit `^`, `build-tags: integration`, ungekürzte Ausgabe, `warn-unused: true`); dauerhafte Ausnahmen nach Entscheidung 2 der [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Testdatei-Ausnahmen des Vorbilds, `ST1005`, `gochecknoglobals` für die zehn Nachschlage-Tabellen und Sentinel-Werte je Datei und Name, `forbidigo` für `cmd/` und `test/`), je mit `# Why:`; Kopfkommentar nennt die Hard Rule aus `AGENTS.md` §3.2 nur so weit, wie die eigene Prüfung sie hält (`AGENTS.md` §3.11) |
| `Dockerfile` | update | Stufe `lint` aus `golangci/golangci-lint:v2.14.0@sha256:ad862ba6b3798cbe0fd9fd7408d498fd74fbd2623a92406b2fd3898faf0bf98f` auf `$BUILDPLATFORM`, Module aus `deps`, `RUN --network=none`, `GOTOOLCHAIN=local`, `GOFLAGS=-mod=readonly`; ruft `tools/harness/lint.sh`; Kommentar der Stufe `source` nennt die erweiterte Allowlist |
| `tools/harness/lint.sh` | neu | die Stufe als ein Skript: die drei eigenen Prüfungen (`SPEC-049` Punkt 6 bis 8, die Why-Prüfung mit der festen Form des Profils und `Form nicht erkannt`, auch für jede Zeile unter `rules` außer Eintrag `- ` mit 6 und Fortsetzung mit mindestens 8), dann nach Punkt 9 Profil vorhanden, `golangci-lint config verify`, `golangci-lint run -c .golangci.yml ./...` und die Warnungen von `warn-unused` als `lint:`-Zeilen; Ausgang nach Punkt 9; Kopf nennt, was es nicht prüft (*Grenze* von `SPEC-049`) |
| `.dockerignore` | update | `.golangci.yml` und `tools/harness/lint.sh` in die Allowlist ([ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)); ohne `.golangci.yml` im Kontext meldet die Stufe `lint: .golangci.yml: fehlt` und golangci-lint läuft nicht (`SPEC-049` Punkt 9), ohne das Skript läuft die Stufe nicht |
| `harness/mk/lint.mk` | neu | Ziel `lint` (`docker build --target lint`), Hilfe-Text „Werkzeug, kein Gate“; kein `GATE_CHECKS +=` |
| `harness/README.md` | update | Tabelle der Werkzeuge: Zeile `make lint`, kein Gate, Bindung an [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) |
| `slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`, `slice-harness-lint` | update | Zählungen nach `SPEC-049` Punkt 9 (16, 16, 32 im Produkt-Code, acht Funktionen) (`AGENTS.md` §3.9) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-lint` liegt nicht mehr in
`in-progress/` (zurück in `next/`, WIP-Limit 1), und
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) ist `Accepted`. Reihenfolge
nach Entscheidung des Nutzers vom 2026-10-06: dieser Slice, `slice-harness-upgrade-v6-16`
(eingeschoben nach Entscheidung des Nutzers vom 2026-10-06), `slice-harness-blackbox-kern`,
`slice-harness-blackbox-driven`, `slice-harness-blackbox-pgwire`,
`slice-harness-blackbox-einstieg`, `slice-lint-bestand-kern-driven`,
`slice-lint-bestand-driving`, `slice-harness-lint`, `slice-harness-abdeckung-gate`,
`slice-harness-coverage`, `slice-harness-mutation`. Vor dem ersten Code-Commit prüft
der Architect §6 gegen `SPEC-049`
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die eigenen Prüfungen sind
  mit Profil und Stufe nicht in einer Review-Sitzung prüfbar; dann trägt dieser Slice
  Profil, Stufe und Ziel, ein eigener Slice davor oder direkt danach die drei eigenen
  Prüfungen.
- `in-progress` → `open` (blockiert — Carveout?): Das gepinnte Image analysiert
  `go 1.27` bei `GOTOOLCHAIN=local` nicht mehr, ein Linter des Profils fehlt in
  `v2.14.0`, oder der Lauf weicht von der Messung der ADR ab, ohne dass die Abweichung
  sich am Profil erklären lässt; dann zuerst eine Entscheidung des Architect.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün ohne `lint` in der Gate-Kette, `make lint` meldet
die 103 Befunde der Messung (`SPEC-049` Punkt 9), Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — entschieden vom Architect am 2026-10-06 in
`SPEC-049` und [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md), gemessen am Stand `79f40e1`; die vier aus der Rückgabe des
Implementer (Profil fehlt, Schema, ungenutzte Regel, Zählregel) am 2026-10-06 in
`SPEC-049` Punkt 9. Dieser Slice entscheidet keine. Was dort nicht steht, gibt der Implementer an den Architect zurück.
Für das Werkzeug tragend:

- **Werkzeug, dann Gate** — `make lint` endet bei Befund mit Ausgang ungleich 0, auch
  als Werkzeug (Punkt 9 und 10); dass es `make gates` nicht rot macht, liegt allein
  daran, dass es nicht an `GATE_CHECKS` hängt.
- **Ort der eigenen Prüfungen** — in der Stufe `lint`, nicht auf dem Host (Punkt 2, 6
  bis 8); ob als `RUN`-Zeilen oder als Skript, ist eine Form, kein Vertrag.
- **Ungenutzte Regel** — `warn-unused` ist ein Befund (Punkt 8); jede dauerhafte
  Ausnahme muss im Bestand etwas ausblenden, sonst ist sie falsch geschrieben. Je
  Warnung eine Zeile `lint: .golangci.yml: Regel ohne Befund: …` ohne Zeilennummer
  (Punkt 9).
- **Profil fehlt** — Zeile `lint: .golangci.yml: fehlt`, Schema-Prüfung und
  golangci-lint laufen nicht, die eigenen Prüfungen nach Punkt 6 und 7 schon; das
  Profil nur über `-c` (Punkt 9).
- **Profil vom Schema abgelehnt** (etwa ein unbekannter Schlüssel, den `golangci-lint
  run` still übergeht) — `golangci-lint config verify` ohne Netz, Zeile `lint:
  .golangci.yml: von golangci-lint abgelehnt` mit der Meldung, golangci-lint läuft
  trotzdem (Punkt 9).
- **Form des Profils** — die Why-Prüfung liest `.golangci.yml` nur in fester Form
  (`exclusions:` mit 2, `rules:` mit 4, `- ` mit 6 Leerzeichen, Blockform); jede andere
  ist `Form nicht erkannt` mit ihrer Zeile. Unter `rules` gilt das für jede Zeile, die
  weder Kommentar, Leerzeile, Eintrag `- ` mit 6 noch dessen Fortsetzung mit mindestens
  8 ist; ein `- ` mit anderem Einzug ist es immer (Punkt 8, nach Review und
  Verifikation entschieden). Eine Regel, die `config verify` annimmt und golangci-lint
  erst beim Laden ablehnt, ist rot ohne `lint:`-Zeile (Grenze von `SPEC-049`; Fall der
  Gegenprobe in `slice-harness-lint`).
- **Zweites Kommentarzeichen** — maßgeblich ist jedes `//` und `/*` der Zeile:
  `// x //nolint` ist ein Befund, `// siehe nolint` nicht (Punkt 6).
- **Felder der ungenutzten Regel** — Reihenfolge Linter, Pfad, Pfad außer, Text,
  Quelle, nur die vorhandenen; ohne eines davon `Felder nicht erkannt` (Punkt 9).
- **Zählregel** — jede Meldung ist ein Befund, auch mehrere auf derselben Zeile
  (`uniq-by-line: false`, Punkt 9). [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) Kontext: 96 gemessen mit `uniq-by-line`
  true; nach `SPEC-049` Punkt 9 sind es 103 bei acht Komplexitäts-Funktionen;
  Entscheidung unberührt, keine Folge-ADR.

**Risiken:**

- **Konfiguration wirkt anders, als sie gelesen wird** — ein Pfad-Regex in
  `exclusions`, ein Modulname in `gomodguard_v2` oder ein fehlendes `build-tags`
  liest sich als Zusage und greift nicht
  (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, 1×). Der Abgleich mit der
  Messtabelle der ADR (DoD, Punkt 3) und `warn-unused` fangen es für diesen Stand.
  — **Ausgang:** eingetreten, an den Folge-Slice `slice-harness-lint`. Viermal las sich
  die Konfiguration anders, als sie wirkt: `exclusions.generated: lax` ohne Wirkung
  (F-433, entfernt in `8d9f0ee`); ein Eintrag für das Modul selbst in `gomodguard_v2`
  wäre ohne Wirkung, weil der Linter Importe aus dem eigenen Modul nicht meldet
  (Implementer, Kommentar im Profil; nicht aufgenommen); Regeln unter `rules` mit
  anderem Einzug wandte golangci-lint an, die Why-Prüfung las sie nicht (V-72, behoben
  in `17daf26`); `config verify` nimmt Regeln an, die golangci-lint erst beim Laden
  ablehnt (V-74, Grenze von `SPEC-049`). Offen bleiben V-74, die Einzüge 8 und 2 und
  V-75; sie stehen als Fälle in der DoD der Gegenprobe von `slice-harness-lint`.
  Register: Beleg in `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (§7).
- **Eigene Prüfungen bis zum Gate ohne Gegenprobe** — die Mutationen sieht der
  Implementer einmal von Hand; die bleibende Gegenprobe kommt erst mit
  `slice-harness-lint`. Bricht eine Prüfung dazwischen still, misst jeder
  Bereinigungs-Slice mit einem stumpfen Werkzeug. Die Gegenprobe von
  `slice-harness-lint` fährt jede Zusage, bevor das Gate scharf wird. Dort ausdrücklich
  zu übernehmen, weil sie hier per Mutation nicht rot wurden oder erst nach dem ersten
  Lauf entschieden sind:
  - ohne Mutation, die der Lauf unterscheidet: `relative-path-mode: cfg` und Pfade mit
    `^` (Profil, Modul und Arbeitsverzeichnis fallen auf `/src` zusammen),
    `--network=none`, `GOFLAGS=-mod=readonly`, `GOTOOLCHAIN=local` (ohne es scheitert der
    Download am fehlenden Netz), `-c` statt Default-Suche, Pin und Plattform der Stufe,
    die untere Grenze von `dupl` (Probe nur auf eine Anweisung genau), die Erkennung der
    ungenutzten Regel am Logtext von golangci-lint;
  - neu nach dem ersten Lauf entschieden (`SPEC-049` Punkt 6, 8 und 9): Profil fehlt,
    Profil vom Schema abgelehnt, Zeile der ungenutzten Regel mit Feldfolge und ohne
    erkannte Felder, feste Form des Profils (Einzug, Flussform, Einträge unter `rules` mit
    Einzug 8 und 2), Ablehnung erst beim Laden (Grenze), `// x //nolint` als
    Befund und `// siehe nolint` als keiner, Ausgang ungleich 0 bei einer `lint:`-Zeile
    allein.
  — **Ausgang:** eingetreten, an den Folge-Slice `slice-harness-lint`. Die
  Why-Prüfung war zweimal stumpf, ohne dass ein Lauf des Implementers es zeigte: bei
  Einzug +2 (F-428, Review) und bei tieferen Einträgen unter `rules` (V-72,
  Verifikation); beide behoben (`8d9f0ee`, `17daf26`). Bis zum Gate misst jeder
  Bereinigungs-Slice mit dem Stand `17daf26`; die Gegenprobe von `slice-harness-lint`
  fährt jede Zusage der beiden Listen oben, dazu V-74 und V-75 (DoD dort).
- **Mutant kommt im Build-Kontext nicht an** — BuildKit überträgt eine Datei gleicher
  Größe und mtime nicht neu (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`, 1×);
  die Mutationen von Hand laufen in einer Kopie unter eigenem Temp-Pfad.
  — **Ausgang:** entfallen: Jede Mutation von Review und Verifikation lief in einer
  frischen Kopie aus `git archive` ohne `cp -p` (Review, *Ausgeführte Läufe*;
  Verifikation, *Wie die Proben gebaut wurden* und Nachtrag); kein Lauf dieses Slice
  baut mehr. Die Gegenprobe von `slice-harness-lint` führt dasselbe Risiko in ihrem §6
  selbst (*Gegenprobe sieht den Mutanten nicht*).
- **Bestehende Prüfung fällt weg** — `gofmt` und `go vet` in der Stufe `test` bleiben;
  die Stufe `lint` tritt daneben, nicht an ihre Stelle
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). — **Ausgang:** entfallen: `git diff 2ee9527..8d9f0ee -- Dockerfile` berührt nur
  Kopfkommentar, Kommentar der Stufe `source` und die neue Stufe (Verifikation, Punkt
  1); `make gates` mit `gofmt`, `go vet` und Unit-Tests der Stufe `test` grün an
  `17daf26`.

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

- **Was hat funktioniert:** Der Schnitt aus §8 hielt: drei Liefer-Punkte in einer Schicht, keine Datei `*.go` im Diff, keine Rückführung aus §4. Die vier Randformen aus der Rückgabe des Implementers (Profil fehlt, Schema, ungenutzte Regel, Zählregel) entschied der Architect vor dem ersten Code-Commit (`4ab7e7b` vor `2461551`); die Abweichung 96 gegen 103 löste damit keine Rückführung `in-progress` → `open` aus, sie ist Folge von `uniq-by-line: false` und in `SPEC-049` Punkt 9 festgelegt. Die Messung trug: 103 Befunde, je Paket und Linter gleich den Zählungen der Folge-Slices (Review, Schwerpunkt 6; Verifikation, Abschnitt 3). `make gates` blieb grün, `make lint` rot, wie Punkt 10 es für das Werkzeug vorsieht.
- **Was ging anders als geplant:** Zwei Prüfrunden. Das Review fand ein MEDIUM, die Why-Prüfung war falsch grün bei anderem Einzug (F-428), dazu zwei Randformen im Code entschieden (F-429), einen Plan-Satz auf altem Stand (F-430) und eine Übergabe an `slice-harness-lint` ohne Artefakt (F-431), F-432 und F-433 als Hinweis. Die Nacharbeit `8d9f0ee` behob F-428 nur für die gemeldete Ausprägung; die Verifikation fand dieselbe Klasse an den Einträgen unter `rules` (V-72). Der Architect entschied sie in `SPEC-049` Punkt 8 (`42d90e8`), `17daf26` schloss sie; die Nacharbeit hat kein eigenes Review, der Nachtrag der Verifikation hat sie geprüft. V-74 steht als Grenze in `SPEC-049` und als Fall in `slice-harness-lint`. V-75 (Hinweis) ist mit dieser Closure in den Gegenproben-Punkt von `slice-harness-lint` eingetragen.
  - **V-73, Tabelle des Implementers zu V-72** (am Stand `17daf26` gegen `8d9f0ee`; Zusage: `SPEC-049` Punkt 8, jeder Eintrag unter `rules` mit anderem Einzug ist `Form nicht erkannt`):

    | Mutation | mit Fix (`17daf26`) | ohne Fix (`8d9f0ee`) |
    |---|---|---|
    | `p8d`, Einträge unter `rules` tiefer, am Prüfstand | Exit 1, `lint: .golangci.yml:159`, `:160`, `:161` `Form nicht erkannt` | Exit 0 |
    | `p8d` am echten Baum | 66 Zeilen `lint:` | keine |
    | Einzug 2 (YAML ungültig) | eine Zeile `Form nicht erkannt`, dazu `von golangci-lint abgelehnt` | nur `von golangci-lint abgelehnt`; der Lauf endet ohnehin rot, weil das YAML ungültig ist |

    Die Verifikation hat dieselben Fälle unabhängig gefahren (Nachtrag, `p8d`, `p8f`, `p8g`). Für `2461551` und `8d9f0ee` lag die Tabelle Review und Verifikation nicht vor; dort sind ihre Läufe der Beleg (DoD, Punkt 2).
  - **Summary-Zeilen der Review-Reports:** Review `2026-10-06-review-slice-harness-lint-werkzeug.md`: „1 MEDIUM (F-428 …), 3 LOW (F-429, F-430, F-431), 2 INFO; keine HIGH. Wiederkehrende Klassen: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-428), `BEO-REPO/plan-folgt-korrektur-nicht` (F-430).“ Verifikation `2026-10-06-verifikation-slice-harness-lint-werkzeug.md`, Urteil des Nachtrags: drei Liefer-Punkte erfüllt, V-72 behoben, V-74 getragen, V-73 und V-75 Hinweise.
- **Steering-Loop-Eintrag:** Geschärfte Regel: Ein Befund aus Review oder Verifikation wird als Klasse behoben; die Nacharbeit nennt das Merkmal, an dem der Wächter scheiterte, und fährt je weitere Ausprägung eine Mutation — liegt in `.claude/commands/implement-slice.md Schritt 19`.
  Auslöser: F-428 und V-72, dieselbe Klasse über zwei Prüfrunden (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, mit diesem Slice 13×). Herkunfts-Anker `seit slice-harness-lint-werkzeug` im Satz nach dem zur mehrteiligen Zusage. **Warum dieser Eintrag:** §3.10 und der Satz zur mehrteiligen Zusage verlangen je Zusage und je Bedingung eine Mutation; die Nacharbeit hat das getan und die gemeldete Mutation (`w4b`, alles +2) rot gesehen. Durch ging nur die Nachbar-Ausprägung desselben Merkmals (`p8d`, nur `rules` tiefer). Der neue Satz setzt an der Nacharbeit an, nicht an der ersten Lieferung. Retirement-Check von §3.9, §3.10, §3.11 und §3.12: alle wieder aufgetreten, alle bleiben.
- **Beobachtungs-Register (`../observations/`):**
  - `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen/` — Beleg `evidence/slice-harness-lint-werkzeug.md` (F-433, das Modul selbst in `gomodguard_v2`, V-72, V-74), **2×**. Erreicht die Schwelle nicht; der nächste Beleg macht ihn zur Lücke, die einen Folge-Slice braucht.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` (verkörpert in `AGENTS.md` §3.11 seit welle-extended-query) — Beleg (F-428, V-72), **14×**.
  - `BEO-REPO/plan-folgt-korrektur-nicht/` (verkörpert in `AGENTS.md` §3.9 seit welle-walking-skeleton, Sensor `make kopf-check` seit slice-harness-kopf-sensor) — Beleg (F-430, F-431), **13×**; `make kopf-check` war grün, §3 und §6 liest er nicht.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` (verkörpert in `AGENTS.md` §3.10 seit welle-extended-query) — Beleg (F-428, V-72), **13×**.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden/` (verkörpert in `AGENTS.md` §3.12 seit slice-harness-randformen-vor-code) — Beleg (F-429, V-72), **11×**. Nicht `randform-im-code-entschieden-dann-zurueckgegeben` (bleibt 2×): Zurückgegeben hat der Implementer die vier Randformen vor dem Code; F-429 kam nicht als Rückgabe, sondern im Review.
  - `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht/` — **neu**, Beleg (V-73), **1×**.
  - Ohne Beleg: `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (bleibt 1×; §4 nannte den Architect vor dem Code, und die Rückgabe erreichte ihn vor `2461551`), `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` und `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (bleiben 1×; Ausgänge in §6).

  Einmalig und nicht eingetragen: F-432 (Hinweis, übergeben an `slice-harness-lint`), V-75 (Hinweis, eingetragen dort). Kein Eintrag erreicht mit diesem Slice die Schwelle 3× neu; über ihr stehen nur verkörperte Einträge (`zusage-im-kommentar-weiter-als-pruefung` 14×, `plan-folgt-korrektur-nicht` 13×, `negativtests-fehlen-bei-neuem-vertrag` 13×, `spec-randform-erst-im-review-entschieden` 11×); ihnen gibt diese Closure keinen Ausgang, den Lese-Schritt führt die nächste Welle-Closure.
- **Folge-Slices:** `slice-harness-lint` (Gate, Gegenprobe mit den Fällen aus §6, V-74 und V-75), `slice-harness-upgrade-v6-16` (als nächster in der Reihe, §4), dann `slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`, `slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`, `slice-lint-bestand-kern-driven` und `slice-lint-bestand-driving`, die je ihr Ergebnis mit `make lint` messen.
- **Risiken aus §6:** vier. *Konfiguration wirkt anders, als sie gelesen wird* und *Eigene Prüfungen bis zum Gate ohne Gegenprobe*: **eingetreten**, an `slice-harness-lint`. *Mutant kommt im Build-Kontext nicht an* und *Bestehende Prüfung fällt weg*: **entfallen**, mit Begründung. Die Randformen in §6 sind Entscheidungen, keine Risiken.
- **Drei Paarungen:** Anker — `liegt in` nennt `.claude/commands/implement-slice.md Schritt 19`; `grep -n "seit slice-harness-lint-werkzeug" .claude/commands/implement-slice.md` findet ihn in Schritt 19. Folge-Slice — `slice-harness-lint` liegt in `next/` und führt V-74, die Einzüge 8 und 2 und V-75 im Gegenproben-Punkt seiner DoD; `slice-harness-upgrade-v6-16` und die übrigen genannten liegen in `open/`. Register — die sechs Kennungen mit Beleg und die drei ohne bestehen als Verzeichnis, jedes mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft erneut.
- **Belege:** Review `docs/reviews/2026-10-06-review-slice-harness-lint-werkzeug.md` (bis `2461551`; F-428 bis F-433), Verifikation `docs/reviews/2026-10-06-verifikation-slice-harness-lint-werkzeug.md` (bis `8d9f0ee`, Nachtrag bis `17daf26`; V-72 bis V-75; `make gates` grün an beiden Ständen), Entscheidung [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) unverändert. Validierung: n/a, der Slice ändert die Prüfumgebung, kein End-Nutzer-Verhalten.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Schnitt** (Größenregel, geprüft vom Planner am 2026-10-06): Der Architect schlug vor,
das Werkzeug mit `slice-harness-blackbox-kern` zu liefern. Jener Slice hätte damit drei
Liefer-Punkte (Umbau der Tests, keine Befunde unter seinen Pfaden, Werkzeug) — an der
Grenze —, aber drei Schichten (Harness mit `Dockerfile`, Tests des Domain Model, Tests
der Services) und ein Review, das Profil, Stufe, drei eigene Prüfungen und acht
umgeschriebene Testdateien mit 31 Befunden zugleich prüft; dazu mischte er Arbeit am
Werkzeug mit Arbeit am Gegenstand. Eigener Slice: drei Liefer-Punkte in einer Schicht
(Harness), einzeln lieferbar, ohne Datei `*.go`.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `7ea7a30` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (12×, verkörpert in `AGENTS.md`
  §3.10) — je Zusage der eigenen Prüfungen eine Mutation, von Hand (DoD, Punkt 2).
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (13×, verkörpert in §3.11) — der
  Kopfkommentar von `.golangci.yml` und die Werkzeug-Zeile sagen nur zu, was der Lauf
  zeigt (§3).
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (1×),
  `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) und
  `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — je ein Risiko in §6.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — wellenlos; §4 nennt den
  Architect vor dem Code.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
