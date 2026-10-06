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
deklariert wird an ihr erst in `slice-harness-abdeckung-gate`). Die ADR des neuen Gates schreibt der Architect im Slice, vor dem
Code; ihre Nummer vergibt der ADR-Index.

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** pt9912. **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Gate `make lint` meldet rot, wenn der Go-Code (`cmd/`, `internal/`,
`test/`) gegen das Lint-Profil in `.golangci.yml` verstößt oder eine
`//nolint`-Direktive trägt. Es läuft mit golangci-lint v2 im per Digest gepinnten
Image als Stufe `lint` des `Dockerfile` und hängt an `GATE_CHECKS`. Vorbild ist das
SOLID-nahe Profil des Schwester-Repos ai-harness-init (Default-Linter dazu
`cyclop`, `funlen`, `gocognit`, `gochecknoglobals`, `gomodguard_v2`, `ireturn`,
`interfacebloat`, `revive` u. a.), an dieses Repo angepasst: `gomodguard_v2` kennt
die Module aus `go.mod` (`github.com/jackc/pgx/v5`, darin `pgproto3`, und
`go.yaml.in/yaml/v3`), `forbidigo` und `exclusions` folgen der CLI und der
Testpraxis dieses Repos (§6). Die Schwellen sind die Werte des Vorbilds
(Entscheidung des Nutzers vom 2026-10-05). `testpackage` ist **gestuft**: Die ADR
des Gates führt es mit einem Hochschalt-Trigger je Paketgruppe auf die vier
Umstellungs-Slices; bis dahin ist es für deren Pfade ausgesetzt, als benannte Stufe,
nicht als stille Ausnahme. Ausnahmen stehen zentral in
`.golangci.yml`, je mit einem `Why:`-Kommentar; `AGENTS.md` §3.2 bekommt damit
seinen Träger. Die Gegenprobe `make lint-gegenprobe` zeigt, dass das Gate rot
werden kann (`AGENTS.md` §3.10).

**Herkunft:** Entscheidung des Nutzers vom 2026-10-05: ein SOLID-naher Lint vor dem
nächsten großen Slice. `harness/README.md` §Sensors führt Lint heute unter
„Nicht behauptet“, und `AGENTS.md` §3.2 ist ein Platzhalter mit einem Gate, das es
nicht gibt. Zum `testpackage` hat der Nutzer am selben Tag entschieden: Die Tests
werden auf Black-Box-Pakete umgestellt, in eigenen Slices.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Schwelle für Testabdeckung — ein anderer Gegenstand mit eigener ADR und
  eigener Messung; ihn übernimmt `slice-harness-coverage`.
- Die Umstellung der Tests auf Black-Box-Pakete — sie berührt alle Schichten und
  sprengte die Größenregel; sie übernehmen `slice-harness-blackbox-kern`,
  `slice-harness-blackbox-driven`, `slice-harness-blackbox-pgwire` und
  `slice-harness-blackbox-einstieg`, je mit dem Hochschalten von `testpackage` für
  ihre Pfade. Dieser Slice legt nur die Stufe und die Form der Export-Test-Brücke in
  der ADR fest (§6).
- LH-QA-07 im Lastenheft — liefert `slice-lastenheft-pruefbarkeit`, vor diesem Slice.
- Eine Abdeckungs-Deklaration an der Gegenprobe — übernimmt
  `slice-harness-abdeckung-gate` nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md): Das Abdeckungs-Skript kennt die
  Nachweisart Gate vor ihm nicht, und Teil 3 von LH-QA-07 gilt erst, wenn `testpackage`
  überall scharf ist.
- Bereinigung des Bestands über das hinaus, was in eine Review-Sitzung passt — ob
  der heutige Code das Profil hält, misst der Architect vor dem Code (§4, §6
  *Bestand*). Reicht die Bereinigung weiter, geht der Slice zurück zur Zerlegung
  (§4), statt still zu wachsen; die Folge-Slices entstehen dann aus der Messung,
  nicht vorab.
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
- Spezifikation und Lastenheft — Schicht-Abgrenzung: Der Slice ändert Harness,
  `Dockerfile`, Lint-Konfiguration und ihre Doku; Produkt-Code nur, soweit die
  Bestands-Bereinigung aus §6 im Umfang bleibt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Gate `make lint`: Stufe `lint` im `Dockerfile` (golangci-lint v2, Image per
      Digest gepinnt, netzlos wie die übrigen Stufen nach `deps`), `.golangci.yml`
      nach den Randformen aus §6 und der ADR des Gates, Fragment unter `harness/mk/`
      an `GATE_CHECKS`; der Bestand ist grün, ohne `//nolint` und ohne Ausnahme
      ohne `Why:`. Die Schwellen sind die des Vorbilds; was der Bestand davon und von
      den übrigen Linter verletzt, ist bereinigt oder steht als Stufe mit
      Hochschalt-Trigger in der ADR. `testpackage` ist gestuft: ausgesetzt je Pfad
      einer Umstellungs-Paketgruppe, jede Stufe nennt in ihrem `Why:` den Slice, der
      sie aufhebt.
- [ ] Gegenprobe `make lint-gegenprobe` an `GATE_CHECKS`: je Zusage des Profils, die
      eine Mutation fangen kann, ein Fall, der rot wird (mindestens ein Verstoß je
      aktivierter Linter-Gruppe, ein `//nolint`, ein Modul außerhalb der Liste von
      `gomodguard_v2`, ein Integrationstest hinter dem Build-Tag, ein White-Box-Test in
      einem Pfad außerhalb der `testpackage`-Stufe, und für die zweite Bedingung von
      Messmethode 3 je ein Fall, in dem ein Test an der Brücke nach der ADR vorbei auf
      Internes zugreift und in dem ein Test in der Brückendatei selbst steht, beide in
      der Form, die die ADR für die Brücke festlegt), dazu ein grüner Fall je zentraler
      Ausnahme und je Stufe; je Zusage ist die Mutation gesehen (`AGENTS.md` §3.10).
      Die Gegenprobe liegt unter `tools/harness/lint-gegenprobe.sh` und ist der Nachweis
      von Teil 1 (statische Analyse) von LH-QA-07, mit den Fällen der Umstellungs-Slices
      auch von Teil 3 (Lage der Unit-Tests); nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) trägt sie in diesem Slice
      keine Abdeckungs-Deklaration, LH-QA-07 steht bis zu `slice-harness-abdeckung-gate`
      in keiner Abdeckungstabelle (`make abdeckung-check` grün).
- [ ] Doku: `AGENTS.md` §3.2 nennt den Träger (kein `//nolint`, Ausnahmen zentral in
      `.golangci.yml` mit `Why:`) statt der Platzhalter; `harness/README.md` §Sensors
      führt `lint` und `lint-gegenprobe` mit Vertrag und Bindung an die ADR, die Zeile
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
| `docs/plan/adr/<NNNN>-lint-gate.md`, `docs/plan/adr/README.md` | neu / update | ADR des neuen Gates (Architect, vor dem Code): Profil mit den Schwellen des Vorbilds, Durchsetzung des `//nolint`-Verbots, Stufen mit Hochschalt-Trigger (`testpackage` je Paketgruppe auf die vier Umstellungs-Slices, dazu jede Stufe aus der Bestands-Messung), Form der Export-Test-Brücke, Randformen aus §6; Index-Zeile |
| `.golangci.yml` | neu | Profil nach Vorbild ai-harness-init, angepasst nach §6; Stufen als benannte Ausnahmen mit `Why:` und Kennung des aufhebenden Slice; Kopfkommentar nennt die Hard Rule aus `AGENTS.md` §3.2 |
| `Dockerfile` | update | Stufe `lint` aus dem per Digest gepinnten golangci-lint-Image, Modul-Cache aus `deps`, `RUN --network=none`; Build-Tag `integration` für `test/integration` |
| `.dockerignore` | update | `.golangci.yml` in die Allowlist des Build-Kontexts ([ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)); ohne sie läuft das Image mit Default-Profil |
| `harness/mk/lint.mk` | neu | `lint` und `lint-gegenprobe`, beide an `GATE_CHECKS` |
| `tools/harness/lint-gegenprobe.sh` | neu | Mutanten in einer Kopie des Arbeitsbaums unter einem Temp-Pfad, je Fall `make lint` bzw. der Docker-Build dort mit erwartetem Exit; Vorbild `make kopf-check-gegenprobe` und `make a-check-negativ`; Rot-Fälle für beide Bedingungen von Messmethode 3 (White-Box-Test, Zugriff an der Brücke vorbei, Test in der Brückendatei); ohne Abdeckungs-Deklaration (§1) |
| Produkt-Code unter `cmd/`, `internal/`, `test/` | refactor | nur die Bestands-Befunde, die nach der Messung (§6 *Bestand*) im Umfang dieses Slice bleiben |
| `internal/hexagon/model/fehler.go` | update | Kommentar an `Meldungen` (Zeilen 171 bis 173) nennt die Tiefensuche der Spezifikation statt „außen nach innen und in der Reihenfolge seiner Ursachen“ (`AGENTS.md` §3.11); übernommen aus der Closure von welle-replay-semantik, Nebenbefund 2. Nur der Kommentar, kein Verhalten |
| `AGENTS.md` | update | §3.2 mit echtem Träger, Falsch/Richtig mit diesem Repo |
| `harness/README.md` | update | §Sensors: zwei Zeilen; „Nicht behauptet“ ohne Lint |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-replay-semantik-fehlerreplay` und
`slice-lastenheft-pruefbarkeit` liegen in `done/` (WIP-Limit 1). Reihenfolge nach
Entscheidung des Nutzers vom 2026-10-05 und 2026-10-06: `slice-lastenheft-pruefbarkeit`,
dieser Slice, `slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
`slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`,
`slice-harness-abdeckung-gate`, `slice-harness-coverage`. Die Umstellungs-Slices folgen dem Lint-Gate, damit jede
Umstellung ihr Hochschalten sofort am Gate belegt. Erster Schritt nach
dem Start, vor jedem Code-Commit: Der Architect misst den Bestand gegen das
Kandidaten-Profil (Befunde je Linter und je Paket), entscheidet die Randformen aus
§6 und schreibt die ADR. Wellenlos heißt hier nicht ohne Architect
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Messung des Bestands
  ergibt Befunde, deren Bereinigung Produkt-Code in mehr als zwei Schichten berührt
  (Domain Model, Services, Adapter, Bootstrap/`cmd`) oder nicht in einer
  Review-Sitzung prüfbar ist. `testpackage` zählt dafür nicht mit, es ist schon
  gestuft. Nach Entscheidung des Nutzers gibt es zwei Wege, die der Architect je
  Linter wählt: Bestand in eigenen Slices bereinigen (je einzeln lieferbar, dieser
  Slice danach), oder gestufte Einführung — ein Teil des Profils sofort, der Rest mit
  Hochschalt-Trigger auf eine Slice-Kennung, beides in der ADR (Baseline-Regelwerk
  `modul-13-quality-gates.md`, bootstrap-aware Gate). Eine Ausnahme in
  `exclusions`, die nur Bestand verschluckt, ist kein dritter Weg, sondern eine
  stille Lockerung (`AGENTS.md` §3.6).
- `in-progress` → `open` (blockiert — Carveout?): Für Go 1.27 gibt es kein
  golangci-lint-Release mit Image, das den Code dieses Repos analysieren kann, oder
  ein Linter des Profils ist in v2 nicht verfügbar und ein Ersatz verlangt eine
  eigene Entscheidung.

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
Entscheidung des Nutzers markiert. Entschieden werden sie vor dem ersten Code-Commit in
der ADR des Gates (Architect; die Entscheidungen des Nutzers vom 2026-10-05 werden
dort festgehalten). Was dort nicht steht, entscheidet der Implementer nicht, er gibt
es zurück.

- **`testpackage`** — der Bestand testet ausnahmslos White-Box: alle Unit-Tests unter
  `internal/` liegen im Paket des Codes (`package services`, `package pgwire`, …),
  `test/integration` ist ein reines Testpaket (`package integration`). **Entschieden
  (Nutzer):** Umstellung auf `_test`-Pakete in vier Slices
  (`slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
  `slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`); bis dahin
  gestuft, mit Hochschalt-Trigger auf diese Slices, in der ADR. Offen: Form der Stufe
  (Ausnahme je Pfad in `exclusions` mit `Why:` und Slice-Kennung, damit jeder
  Umstellungs-Slice genau seine Zeile löscht, oder Linter ganz aus bis zum letzten
  Slice — dann belegt keine Zwischenstufe ihr Hochschalten), und ob `cmd/` wie im
  Vorbild dauerhaft ausgenommen ist (`main` ohne Test heute).
- **Export-Test-Brücke** (`export_test.go`) — der Bestand greift in Tests auf
  unexportierte Teile zu (Namensabgleich per Suche am Stand `ce50a10`, ungemessen):
  `pgwire` etwa zehn Namen (u. a. `meldeFrist`, `wecke`, `toClientMessage`),
  `services` etwa sieben (u. a. `cursor`, `vtLeerraum`, `istLebendpruefung`),
  `postgres` drei, `cli` und `bootstrap` je einen, `model` und `recording` keinen.
  Darunter Mutationstests auf unexportierte Teile und Zustandseingriffe (Fristen).
  Offen: Brücke zulässig ja/nein; Form (Datei `export_test.go` im Paket, exportierte
  Aliase, Namensregel); ob sie nur lesen oder auch Zustand setzen darf; wie `revive`
  und `gochecknoglobals` sie behandeln (eine Brücke aus `var X = x` ist eine globale
  Variable in einer `_test.go`-Datei). Die Umstellungs-Slices folgen dieser
  Entscheidung, sie treffen sie nicht.
- **`forbidigo`** — die CLI schreibt über injizierte Writer (`bootstrap.Run` erhält
  `os.Stdout` und `os.Stderr` aus `cmd/pgwire-recorder/main.go`); `fmt.Print*` kommt
  im Bestand nicht vor. Offen: Muster wie im Vorbild, und ob zusätzlich `os.Stdout`
  und `os.Stderr` außerhalb von `cmd/` verboten werden.
- **`gochecknoglobals`** — Bestand: `main.version` (gesetzt per `-ldflags -X`,
  Dockerfile-Stufe `build`, wie die Ausnahme im Vorbild), Nachschlage-Tabellen auf
  Paketebene (u. a. in `model`, `services`, `cli`, `recording`, `postgres`), ein
  `var`-Block in `bootstrap`, Sentinel-Fehler `Err…` und Zusicherungen
  `var _ Port = …`. Offen: welche davon Ausnahme mit `Why:` sind, welche Bereinigung.
- **`gomodguard_v2`** — Liste erlaubter Module (dann braucht jedes neue Modul eine
  Konfig-Änderung im selben Commit wie `go.mod`) oder Liste verbotener wie im Vorbild.
  Erlaubt nach `go.mod`: `github.com/jackc/pgx/v5` (darin `pgproto3`) und
  `go.yaml.in/yaml/v3`; offen ist, wie indirekte Module zählen. d-migrate ist kein
  Go-Modul, sondern das Werkzeug-Image von `make schema-validate`, und gehört nicht in
  die Liste. Abgrenzung zur Bibliotheks-Regel je Adapter in `.a-check.yml` (§1).
- **Durchsetzung des `//nolint`-Verbots** — golangci-lint wertet `//nolint` aus, statt
  es zu verbieten. Offen: eigener Prüfschritt in der Stufe `lint` oder ein Linter;
  was als Direktive zählt (Zeilenanfang, nach Code, Groß-/Kleinschreibung,
  `//nolint:x` mit Begründung) und was nicht (das Wort in Prosa, in einem String,
  in einer Testvorlage).
- **Ausnahmen mit `Why:`** — Form des Kommentars (Zeile über der Regel, Pflicht je
  Regel) und ob ein Sensor sie prüft oder Review. Übernahme der `_test.go`-Ausnahmen
  des Vorbilds (Komplexität, `funlen`, `noctx`, `unparam`, ungenutzte Parameter und
  Receiver) gegen die Testdateien dieses Repos prüfen.
- **Integrationstests hinter dem Build-Tag** — `test/integration` übersetzt nur mit
  `-tags integration`; ohne `run.build-tags` lintet golangci-lint diese Dateien nicht.
  Offen: mitlinten (wie `go vet -tags integration` in der Stufe `test`) oder
  ausdrücklich nicht.
- **Schwellen des Profils** (`cyclop`, `funlen`, `gocognit`, `gocyclo`, `nestif`,
  `maintidx`, `dupl`, `interfacebloat`) — **entschieden (Nutzer):** die Werte des
  Vorbilds (`cyclop` 15, `funlen` 100 Zeilen / 60 Anweisungen, `gocognit` 20,
  `gocyclo` 15, `nestif` 5, `maintidx` 20, `dupl` 150, `interfacebloat` 10). Was der
  Bestand verletzt, wird bereinigt oder mit ADR und Hochschalt-Trigger gestuft; kein
  Wert wird an den Bestand angepasst. Eine spätere Senkung braucht eine neue ADR
  (`AGENTS.md` §3.6). Offen: ob die `_test.go`-Ausnahmen des Vorbilds für Komplexität
  und `funlen` übernommen werden (Ausnahme, keine Schwelle).
- **`revive`-Regeln `exported` und `package-comments`** — der Bestand kommentiert
  deutsch; die Regel prüft die Form (Kommentar beginnt mit dem Namen), nicht die
  Sprache. Offen: Regel wie im Vorbild, und ob die Doc-Comment-Pflicht auch für
  `internal/` gilt.
- **Image und Version** — golangci-lint v2 im Image per Digest gepinnt; die
  Go-Version des Images muss `go 1.27` aus `go.mod` bei `GOTOOLCHAIN=local` tragen
  (Vorbild: v2.14.0). Offen: Pin-Ort (Makefile-Variable wie im Vorbild oder nur im
  `Dockerfile`), Plattform (`$BUILDPLATFORM` wie `deps`), Modul-Cache aus `deps` und
  Lauf `--network=none`.
- **Ausgabe und Exit** — Befunde im Format von golangci-lint auf der Build-Ausgabe
  (`--progress=plain`), Exit ≠ 0 bei mindestens einem Befund, der Befund des
  `//nolint`-Schritts nennt Datei und Zeile. Offen: ob das Gate Befunde begrenzt
  (`max-issues-per-linter`, `max-same-issues`; Default kürzt die Ausgabe) — ein
  gekürzter Bericht sagt nicht, wie viele Befunde es gibt.

**Risiken:**

- **Bestand** (Hauptrisiko) — unbekannt, wie viele Befunde der heutige Code gegen das
  Profil liefert. Sichtbar ohne Messung: `testpackage` träfe jede Testdatei
  (gestuft, siehe oben), `gochecknoglobals` mehrere Pakete, und die größten Dateien
  (`internal/adapters/driving/pgwire/server.go`,
  `internal/adapters/driven/recording/yaml.go`, `internal/hexagon/services/replay.go`)
  sind Kandidaten für `funlen` und `cyclop`. Der Architect misst vor dem Code (§4);
  zu viel für diesen Slice heißt Rückführung `in-progress` → `next`, nicht
  Ausnahme. — **Ausgang:** — (bei Closure)
- **Stufe wird zum Dauerzustand** — eine `testpackage`-Stufe, deren Umstellungs-Slice
  nie startet, ist eine stille Ausnahme mit Aufschrift. Jede Stufe nennt in ihrem
  `Why:` die Kennung des Slice, der sie aufhebt; die Kennung liegt als Datei in
  `open/`. — **Ausgang:** — (bei Closure)
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
- `BEO-REPO/harness-lesart-ohne-entscheidungsort` (1×) — für Harness-Gates gibt es
  keine Spezifikation; Entscheidungsort der Randformen ist die ADR des Gates.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (1×),
  `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) und
  `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — je ein Risiko in §6.
- `BEO-REPO/kern-fremdimporte-nur-ueber-tech-liste` (1×) — berührt, nicht
  aufgenommen (§1).

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
