# Slice slice-harness-lint: Lint-Gate mit SOLID-nahem Profil

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Er trägt zum
Abnahmeszenario 17 (M3) bei, das `slice-lastenheft-pruefbarkeit` anlegt und das mit
dem letzten Slice der Reihe nachweisbar wird. Eingesammelt wird er von der
nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-05: nach `slice-replay-semantik-fehlerreplay` und vor dem nächsten großen
Slice (WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode 1; Messmethode 3 mit den vier Umstellungs-Slices).
Bindung an Entscheidungen: [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)
(Gates als Stufen des Multistage-`Dockerfile`, Images per Digest gepinnt),
[ADR-0001](../../adr/0001-hexagonale-architektur.md) (Abgrenzung zum
Architektur-Gate), [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (die Gegenprobe ist Nachweis der Nachweisart Gate,
deklariert wird an ihr erst in `slice-harness-abdeckung-gate`),
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (das Gate, seine dauerhaften
Ausnahmen, die Messung des Bestands, Einführung nach Bereinigung ohne Stufen und das
Werkzeug-Ziel davor, Entscheidung 5 und 6).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--lint-profil-lint) (`spezifikation.md` §11: Vertrag des Gates, vom Architect vor dem Code geschrieben; dieser Slice liefert die zweite Hälfte von Punkt 10, den Anschluss an die Gate-Kette) · `spezifikation.md` §12 (*Historie*)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make lint` wird Gate: Es hängt an `GATE_CHECKS` und macht `make gates`
rot, wenn der Go-Code (`cmd/`, `internal/`, `test/`) gegen das Lint-Profil in
`.golangci.yml` verstößt oder eine `//nolint`-Direktive trägt (`SPEC-049` Punkt 10,
zweite Hälfte). Das Ziel selbst — Profil, Stufe `lint` des `Dockerfile` mit
golangci-lint v2 im per Digest gepinnten Image, die drei eigenen Prüfungen — liefert
`slice-harness-lint-werkzeug` als Werkzeug ohne Gate (Entscheidung 6 in
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md)). Den Vertrag hält
`SPEC-049`; Entscheidung, Gründe und die Messung des Bestands hält
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md). Der Bestand ist beim
Start dieses Slice **bereinigt, ohne Stufen** (Entscheidung des Nutzers vom
2026-10-06): `.golangci.yml` trägt nur die dauerhaften Ausnahmen aus Entscheidung 2,
je mit `Why:`, und `testpackage` ist vom ersten Gate-Lauf an für jeden Pfad scharf.
Die Gegenprobe `make lint-gegenprobe` zeigt, dass das Gate rot werden kann
(`AGENTS.md` §3.10), einschließlich der White-Box-Fälle je Paketgruppe, die bisher die
vier Umstellungs-Slices tragen sollten. `AGENTS.md` §3.2 bekommt seinen Träger.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-05: ein SOLID-naher Lint vor dem
nächsten großen Slice. `harness/README.md` §Sensors führt Lint heute unter
„Nicht behauptet“, und `AGENTS.md` §3.2 ist ein Platzhalter mit einem Gate, das es
nicht gibt. Zum `testpackage` hat der Nutzer am selben Tag entschieden: Die Tests
werden auf Black-Box-Pakete umgestellt, in eigenen Slices. Nach der Messung hat der
Nutzer am 2026-10-06 entschieden: Der Bestand wird vor dem Gate bereinigt, ohne
Stufen; dieser Slice ging dafür nach `next/` zurück (§4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Schwelle für Testabdeckung — ein anderer Gegenstand mit eigener ADR und
  eigener Messung; ihn übernimmt `slice-harness-coverage`.
- Das Werkzeug-Ziel `make lint` (Profil `.golangci.yml`, Stufe `lint` im
  `Dockerfile`, die drei eigenen Prüfungen, `.dockerignore`) — übernimmt
  `slice-harness-lint-werkzeug`, als erster der Reihe, damit jeder Bereinigungs-Slice
  sein Ergebnis mit demselben Ziel misst, das hier Gate wird (Entscheidung 6 in
  [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md)).
- Die Umstellung der Tests auf Black-Box-Pakete samt den übrigen 47 Befunden in
  Testdateien — sie berührt alle Schichten und sprengte die Größenregel; sie
  übernehmen `slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
  `slice-harness-blackbox-pgwire` und `slice-harness-blackbox-einstieg`, je in den
  Testdateien, die sie umschreiben. Die Form der Export-Test-Brücke steht in
  `SPEC-049` Punkt 7.
- LH-QA-07 im Lastenheft — liefert `slice-lastenheft-pruefbarkeit`, vor diesem Slice.
- Eine Abdeckungs-Deklaration an der Gegenprobe — übernimmt
  `slice-harness-abdeckung-gate` nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md): Das Abdeckungs-Skript kennt die
  Nachweisart Gate vor ihm nicht. Teil 3 von LH-QA-07 gilt ab diesem Slice, weil
  `testpackage` erst mit dem Gate für alle Pfade scharf ist.
- Bereinigung des Produkt-Codes — gemessen (§6 *Bestand*,
  [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md)): 25 Befunde nach den
  dauerhaften Ausnahmen, in Kern, Driven und Driving. Sie übernehmen
  `slice-lint-bestand-kern-driven` (13) und `slice-lint-bestand-driving` (12), nach den
  Umstellungs-Slices. Dieser Slice ändert keinen Produkt-Code außer dem Kommentar in
  `internal/hexagon/model/fehler.go` (§3).
- Importrichtungen im Hexagon und Bibliotheken je Adapter — das hält
  `make a-check` ([ADR-0001](../../adr/0001-hexagonale-architektur.md)).
  `gomodguard_v2` prüft Module gegen eine Liste, keine Richtung; eine zweite Quelle
  für dieselbe Regel driftete gegen `.a-check.yml`.
- Standardbibliotheks-Importe im Domain Model (`depguard` o. Ä.) — anderer Vorgang:
  Die Lücke ist als `BEO-REPO/kern-fremdimporte-nur-ueber-tech-liste` (1×) im
  Register und betrifft das Architektur-Gate; wird sie Muster, ist sie ein eigener
  Slice aus dem Register.
- Shell-Lint für `tools/` (shellcheck) — anderer Gegenstand (andere Sprache, anderes
  Image); bleibt bewusst außen vor, bis ein Slice ihn nennt.
- Ein Freshness-Sensor für die golangci-lint-Version — er bräuchte Netz über
  `bash`, `git` und `docker` hinaus (`harness/conventions.md` §Baseline, dieselbe
  Begründung); die Anhebung bleibt ein bewusster Commit wie bei jedem gepinnten Image.
- Lastenheft und Spezifikation — Schicht-Abgrenzung: `SPEC-049` und die
  Historien-Zeile hat der Architect vor dem Code geschrieben; dieser Slice ändert das
  Gate-Fragment, die Gegenprobe und die Doku, Produkt-Code nur den Kommentar aus §3.
  Braucht das Gate eine Änderung am Profil oder an der Stufe, geht sie an den
  Architect zurück, weil `slice-harness-lint-werkzeug` sie geliefert hat.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Gate `make lint`: Das Fragment `harness/mk/lint.mk` hängt `lint` an
      `GATE_CHECKS` (`SPEC-049` Punkt 10, zweite Hälfte), und `make gates` ist mit ihm
      grün. Der Bestand ist dafür grün, ohne `//nolint`, und `.golangci.yml` trägt nur
      die dauerhaften Ausnahmen nach `SPEC-049` Punkt 8, je mit `Why:`, keine Regel, die
      nur Bestand aussetzt, und keine ungenutzte Regel; `testpackage` ist für jeden Pfad
      scharf.
- [ ] Gegenprobe `make lint-gegenprobe` an `GATE_CHECKS`: je Zusage aus `SPEC-049`, die
      eine Mutation fangen kann, ein Fall, der rot wird (mindestens ein Verstoß je
      aktivierter Linter-Gruppe, ein `//nolint`, ein Modul außerhalb der Liste von
      `gomodguard_v2`, ein Integrationstest hinter dem Build-Tag, eine ungenutzte Regel,
      eine Regel ohne `Why:`; für die erste Bedingung von Messmethode 3 je Paketgruppe
      der Umstellungs-Slices — Kern, Driven, PGWire, Einstieg — ein White-Box-Test, der
      rot wird; für die zweite Bedingung ein Test in `internal_test.go` im Paket des
      Codes, eine Funktion `Test…` und eine Variable in `export_test.go`), dazu ein
      grüner Fall je dauerhafter Ausnahme; je Zusage ist die Mutation gesehen
      (`AGENTS.md` §3.10). Die Gegenprobe liegt unter `tools/harness/lint-gegenprobe.sh`
      und ist der Nachweis von Teil 1 (statische Analyse) und Teil 3 (Lage der
      Unit-Tests) von LH-QA-07; nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) trägt sie in diesem Slice keine
      Abdeckungs-Deklaration, LH-QA-07 steht bis zu `slice-harness-abdeckung-gate` in
      keiner Abdeckungstabelle (`make abdeckung-check` grün).
- [ ] Doku: `AGENTS.md` §3.2 nennt den Träger (kein `//nolint`, Ausnahmen zentral in
      `.golangci.yml` mit `Why:`) statt der Platzhalter; `harness/README.md` führt
      `lint` nicht mehr unter den Werkzeugen, sondern unter §Sensors, dazu
      `lint-gegenprobe`, je mit Vertrag und Bindung an die ADR; die Zeile
      „Nicht behauptet“ nennt Lint nicht mehr.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0034-lint-gate-mit-solid-nahem-profil.md`, `docs/plan/adr/README.md` | neu / update | geschrieben vom Architect vor dem Code: Entscheidung, Gründe, Messung des Bestands, dauerhafte Ausnahmen, Einführung nach Bereinigung, Werkzeug-Ziel davor; Index-Zeile. Den Status `Accepted` nach der Wahl des Nutzers vom 2026-10-06 setzt der Architect vor dem Start von `slice-harness-lint-werkzeug` |
| `spec/spezifikation.md` §11, §12 | update | geschrieben vom Architect vor dem Code: `SPEC-049` mit den Randformen aus §6; Historien-Zeile |
| `harness/mk/lint.mk` | update | aus `slice-harness-lint-werkzeug`: `lint` an `GATE_CHECKS`, dazu das Ziel `lint-gegenprobe`, ebenfalls an `GATE_CHECKS` |
| `tools/harness/lint-gegenprobe.sh` | neu | Mutanten in einer Kopie des Arbeitsbaums unter einem Temp-Pfad, je Fall `make lint` bzw. der Docker-Build dort mit erwartetem Exit; Vorbild `make kopf-check-gegenprobe` und `make a-check-negativ`. Je Zusage aus `SPEC-049` ein Mutant, der rot wird; Rot-Fälle für die erste Bedingung von Messmethode 3 je Paketgruppe (ein White-Box-Test unter `internal/hexagon/`, `internal/adapters/driven/`, `internal/adapters/driving/pgwire` und unter `internal/adapters/driving/cli`, `internal/bootstrap` oder `test/integration`), für die zweite Bedingung Zugriff an der Brücke vorbei (`internal_test.go` im Paket des Codes), Test in der Brückendatei, Variable in der Brückendatei; dazu eine ungenutzte Regel und eine Regel ohne `Why:`. Den grünen Fall je Ausnahme trägt der Bestands-Lauf mit `warn-unused` (`SPEC-049` Punkt 8). Ohne Abdeckungs-Deklaration (§1) |
| `internal/hexagon/model/fehler.go` | update | Kommentar an `Meldungen` (Zeilen 171 bis 173) nennt die Tiefensuche der Spezifikation statt „außen nach innen und in der Reihenfolge seiner Ursachen“ (`AGENTS.md` §3.11); übernommen aus der Closure von welle-replay-semantik, Nebenbefund 2. Nur der Kommentar, kein Verhalten |
| `AGENTS.md` | update | §3.2 mit echtem Träger, Falsch/Richtig mit diesem Repo |
| `harness/README.md` | update | Zeile `make lint` aus den Werkzeugen nach §Sensors, dazu `make lint-gegenprobe`; „Nicht behauptet“ ohne Lint |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-lint-werkzeug`, die vier
Umstellungs-Slices (`slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
`slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`) und beide
Bereinigungs-Slices (`slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`)
liegen in `done/`, und `make lint` meldet am Stand des Starts keinen Befund (WIP-Limit
1). Reihenfolge nach Entscheidung des Nutzers vom 2026-10-05 und 2026-10-06:
`slice-harness-lint-werkzeug`, die vier Umstellungs-Slices, die beiden
Bereinigungs-Slices, dieser Slice, `slice-harness-abdeckung-gate`,
`slice-harness-coverage`, `slice-harness-mutation`. Vor dem ersten Code-Commit prüft
der Architect, ob die Randformen aus §6 noch dem Werkzeug entsprechen
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführung `in-progress` → `next` am 2026-10-06 (eingetreten).** Vorab benannt war:
Die Messung des Bestands ergibt Befunde, deren Bereinigung Produkt-Code in mehr als
zwei Schichten berührt oder nicht in einer Review-Sitzung prüfbar ist. **Grund:** Die
Bedingung trat ein. Die Messung des Architect (Stand `79f40e1`, Messtabelle in
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md)) ergab 166 Befunde, nach
den dauerhaften Ausnahmen 96: 25 im Produkt-Code in allen Schichten, darunter acht
Komplexitätsbefunde, die schwerste `(*Server).replaySitzung` mit `gocognit` 37, und 71 in
Testdateien. Der Bestand ist zu groß für diesen Slice. Der Nutzer hat am 2026-10-06
zwischen gestufter Einführung und Bereinigung vor dem Gate gewählt: Bereinigung, ohne
Stufen (Entscheidung 5 der ADR). Daraus folgen das Werkzeug-Ziel vorab
(`slice-harness-lint-werkzeug`), der Mehrumfang der vier Umstellungs-Slices (47 Befunde
außer `testpackage`) und zwei neue Bereinigungs-Slices; dieser Slice behält Gate,
Gegenprobe und Doku.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Gegenprobe ist mit den
  Fällen je Zusage, je Paketgruppe und je Brücken-Fall nicht in einer Review-Sitzung
  prüfbar; dann trägt dieser Slice Gate, Doku und die Fälle zu Messmethode 1, und ein
  eigener Slice vor `slice-harness-abdeckung-gate` die Fälle zu Messmethode 3.
- `in-progress` → `open` (blockiert — Carveout?): `make lint` meldet beim Start oder
  im Lauf einen Befund außerhalb der dauerhaften Ausnahmen, etwa aus Code, der nach
  der Bereinigung kam. Eine Stufe dafür gibt es nicht (Entscheidung 5); der Befund geht
  an den Slice, der ihn einbrachte, oder an einen neuen Bereinigungs-Slice.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün mit `lint` und `lint-gegenprobe`, Closure-Notiz
mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen des Vertrags** (`AGENTS.md` §3.12) — offen, soweit nicht als
Entscheidung des Nutzers markiert. Entschieden werden sie vor dem ersten Code-Commit vom
Architect und festgehalten im Abschnitt für Harness-Werkzeuge der Spezifikation (Technik-Stratum; angelegt von `slice-harness-vertraege-spezifikation`, ein bestehendes Werkzeug schreibt seine Kennung fort, ein neues bekommt die nächste freie Kennung, vergeben vom Slice, der es liefert); die ADR des Gates trägt Entscheidung und Gründe
und verweist mit `Schärft:` auf die Stelle. Die Entscheidungen des Nutzers vom
2026-10-05 werden dort festgehalten. Was dort nicht steht, entscheidet der Implementer
nicht, er gibt es zurück.

Entschieden am 2026-10-06 vom Architect vor dem ersten Code-Commit, festgehalten in
`SPEC-049` (Punktnummern unten) und [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md); gemessen gegen den
Bestand am Stand `79f40e1`.

- **`testpackage`** — ohne Stufe: Die vier Umstellungs-Slices beheben seine Befunde
  vor dem Gate, und ab dem Gate ist es für jeden Pfad scharf (Punkt 8 und 10). Nicht
  ganz aus. `skip-regexp` nur `export_test.go`: Der Default ließ auch
  `internal_test.go` White-Box durch (gemessen, Punkt 5). `cmd/` ist nicht
  ausgenommen; `testpackage` lässt Tests im Paket `main` selbst zu, `cmd/` hat keine
  (Grenze). Entschieden (Nutzer): Umstellung in den vier Slices; die White-Box-Fälle
  je Paketgruppe trägt die Gegenprobe dieses Slice.
- **Export-Test-Brücke** — zulässig, nur als `export_test.go` im Paket des Codes, mit
  Aliasen, Konstanten und weiterreichenden Funktionen; Zustand nur an übergebenen
  Werten (Punkt 7). Keine Ausnahme für `gochecknoglobals` oder `revive` in der Brücke:
  eine Variable dort ist ein Befund. Rot-Fälle nach F-407: `internal_test.go` im
  Paket des Codes (an der Brücke vorbei, `testpackage`) und `func Test…` in
  `export_test.go` (eigene Prüfung), dazu eine Variable in der Brücke.
- **`forbidigo`** — `fmt.Print…`, `print`, `println`, dazu `os.Stdout` und `os.Stderr`
  außerhalb von `cmd/` und `test/` (Punkt 5, dauerhafte Regel nach Punkt 8); Bestand
  ohne Befund.
- **`gochecknoglobals`** — `main.version` braucht keine Ausnahme (der Linter lässt
  `version` zu, gemessen); `Err…` und `var _ Port = …` ebenso. Die zehn
  Nachschlage-Tabellen und Sentinel-Werte im Produkt-Code sind dauerhafte Ausnahmen je
  Datei und Name; die 23 Globalen in Testdateien beheben die Umstellungs-Slices (Punkt 8,
  [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 2 und 5).
- **`gomodguard_v2`** — Erlaubnisliste `github.com/jackc/pgx/v5` und
  `go.yaml.in/yaml/v3`; ein indirektes Modul aus `go.mod` ist beim direkten Import ein
  Befund (Mutant mit `golang.org/x/sync` rot gesehen). Ein neues Modul ändert die Liste
  im selben Commit wie `go.mod` (Punkt 5).
- **Durchsetzung des `//nolint`-Verbots** — eigene Prüfung in der Stufe `lint`: `//`
  oder `/*`, dann Leerraum oder `/`, dann `nolint` in beliebiger Schreibweise, überall
  in der Zeile, auch in Strings; Prosa mit einem Wort davor zählt nicht (Punkt 6).
- **Ausnahmen mit `Why:`** — Kommentarblock unmittelbar über jeder Regel, erste Zeile
  `# Why:`; eine eigene Prüfung meldet das Fehlen; `warn-unused` macht eine Regel ohne
  Wirkung zum Befund (Punkt 8). Die `_test.go`-Ausnahmen des Vorbilds sind übernommen
  (sie blenden 39 Befunde aus, Grund gilt für jeden Test).
- **Integrationstests hinter dem Build-Tag** — mitgelintet (`run.build-tags:
  integration`, Punkt 1); in der Messung 28 Befunde unter `test/integration`.
- **Schwellen des Profils** — entschieden (Nutzer): die Werte des Vorbilds (Punkt 4);
  die `_test.go`-Ausnahmen für Komplexität und `funlen` übernommen.
- **`revive`-Regeln `exported` und `package-comments`** — wie im Vorbild, auch für
  `internal/` (Punkt 5); Bestand: sieben Kommentare in `model` ohne die Form, behoben
  von `slice-lint-bestand-kern-driven`.
- **Image und Version** — `golangci/golangci-lint:v2.14.0@sha256:ad862ba6…` (gebaut mit
  go1.27.0, Image go1.27.1, analysiert `go 1.27` bei `GOTOOLCHAIN=local`, gemessen); Pin
  nur im `Dockerfile`, `$BUILDPLATFORM`, Module aus `deps`, `--network=none` (Punkt 2).
- **Ausgabe und Exit** — ungekürzt (`max-issues-per-linter: 0`, `max-same-issues: 0`,
  `uniq-by-line: false`); eigene Prüfungen `lint: <pfad>:<zeile>: <befund>`; alle
  Prüfungen laufen, Exit ≠ 0 bei Befund, ungenutzter Regel oder abgelehnter
  Konfiguration (Punkt 9).
- **Pfade in `.golangci.yml`** (neu, aus der Messung) — golangci-lint liest Pfade
  relativ zur Konfiguration; eine verankerte Regel griff in der Messung nicht, bis die
  Konfiguration an der Modulwurzel lag. Entschieden: `relative-path-mode: cfg`, Pfade
  mit `^` (Punkt 1).
- **Einführung des Bestands** — entschieden (Nutzer, 2026-10-06): Bereinigung vor dem
  Gate, ohne Stufen ([ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 5, `SPEC-049` Punkt 10); die
  gestufte Einführung, die der Architect empfahl, ist als Alternative G verworfen. Bis
  zum Gate ist `make lint` ein Werkzeug ohne Gate (Entscheidung 6), geliefert von
  `slice-harness-lint-werkzeug`; dieser Slice ging nach `next/` zurück (§4).

**Risiken:**

- **Bestand** (Hauptrisiko) — gemessen ([ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md)): 142 Befunde ohne `testpackage`,
  166 mit; nach den dauerhaften Ausnahmen 96, davon 25 im Produkt-Code in allen
  Schichten (sieben Funktionen über den Komplexitäts-Schwellen) und 71 in Testdateien.
  Eingetreten: Rückführung nach `next/` (§4); den Bestand bereinigen die vier
  Umstellungs-Slices, `slice-lint-bestand-kern-driven` und `slice-lint-bestand-driving`
  vor dem Start dieses Slice. — **Ausgang:** — (bei Closure)
- **Bestand wächst bis zum Gate nach** — zwischen der Bereinigung und dem Gate prüft
  niemand neuen Code automatisch ([ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md), Konsequenzen, negativ akzeptiert); ein
  Befund aus dieser Zeit stünde beim Start ohne Stufe da. Die Start-Bedingung in §4
  verlangt `make lint` ohne Befund, die Rückführung `→ open` gibt einen Befund an den
  Slice, der ihn einbrachte. — **Ausgang:** — (bei Closure)
- **Gegenprobe sieht den Mutanten nicht** — BuildKit überträgt eine Datei gleicher
  Größe und mtime nicht neu (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`, 1×;
  `.claude/commands/implement-slice.md` Schritt 19). Die Gegenprobe baut aus einer
  Kopie unter eigenem Temp-Pfad, nicht aus dem Arbeitsbaum. — **Ausgang:** — (bei
  Closure)
- **Konfiguration wirkt anders, als sie gelesen wird** — ein Pfad-Regex in
  `exclusions`, ein Modulname in `gomodguard_v2` oder ein fehlendes `build-tags`
  liest sich als Zusage und greift nicht
  (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, 1×). Je Ausnahme und je
  Liste ein Fall der Gegenprobe. — **Ausgang:** — (bei Closure)
- **Bestehende Prüfung fällt weg** — `gofmt` und `go vet` in der Stufe `test` bleiben;
  die Stufe `lint` tritt daneben, nicht an ihre Stelle
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). — **Ausgang:** — (bei Closure)
- **Laufzeit von `make gates`** — eine weitere Docker-Stufe mit eigener Analyse;
  `make -j` fährt sie parallel. Wird sie zum Engpass, ist das eine Beobachtung, kein
  Grund für eine Lockerung. — **Ausgang:** — (bei Closure)

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `8b399ba` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (8×, verkörpert in `AGENTS.md`
  §3.10) — trifft das Gate unmittelbar: Ein Lint-Gate, das nur grün gesehen wurde,
  belegt keine Zusage; daher der zweite Liefer-Punkt, je Zusage eine Mutation.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (9×, verkörpert in §3.11) — der
  Kopfkommentar von `.golangci.yml`, die Sensors-Zeile und §3.2 sagen nur zu, was die
  Gegenprobe prüft.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (6×, verkörpert in §3.12) und
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (2×) — darum stehen die
  Randformen in §6 offen und werden vor dem Code entschieden. Mit einem dritten
  Auftreten in diesem Slice erreichte der zweite Eintrag die Schwelle.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — dieser Slice ist
  wellenlos; §4 *Start* nennt den Architect-Schritt vor dem Code ausdrücklich.
- `BEO-REPO/harness-lesart-ohne-entscheidungsort` (1×) — Entscheidungsort der
  Randformen ist der Abschnitt für Harness-Werkzeuge der Spezifikation, angelegt von
  `slice-harness-vertraege-spezifikation`; die ADR trägt Entscheidung und Gründe.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (1×),
  `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) und
  `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — je ein Risiko in §6.
- `BEO-REPO/kern-fremdimporte-nur-ueber-tech-liste` (1×) — berührt, nicht
  aufgenommen (§1).

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
