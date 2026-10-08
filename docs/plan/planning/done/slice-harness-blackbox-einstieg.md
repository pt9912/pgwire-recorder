# Slice slice-harness-blackbox-einstieg: Black-Box-Tests von CLI, Bootstrap und Integrationstests

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

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode 3; Messmethode 1 für die Testdateien dieser Pakete). Bindung: [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Export-Test-Brücke, Entscheidung 4; Bereinigung vor dem Gate ohne Stufen, Entscheidung 5: die Befunde der Testdateien dieser Pakete behebt dieser Slice), [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (Abdeckungs-Deklarationen an den Tests bleiben unverändert), [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) (die Deklaration von Teil 3 an der Lint-Gegenprobe setzt `slice-harness-abdeckung-gate`).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--lint-profil-lint) (Punkt 7 und 8: Brücke und dauerhafte Ausnahmen, gegen die umgestellt wird; der Slice ändert die Stelle nicht)

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

**Ziel:** Die Unit-Tests unter `internal/adapters/driving/cli`, `internal/bootstrap` und `test/integration` laufen als Black-Box-Pakete
(`package <name>_test`) und prüfen über die exportierte Schnittstelle des Pakets;
keine Testdatei unter diesen Pfaden hat in `make lint` einen Befund. Dafür behebt der
Slice in den Testdateien, die er umschreibt, die 22 Befunde aus der Messung in
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Stand `79f40e1`; Aufteilung je Paketgruppe in der Fassung `92d1b86`
der ADR): `testpackage` 9, dazu `errcheck`, `gochecknoglobals`, `revive` und `staticcheck`, zusammen 13. Wo ein Test einen unexportierten Teil braucht, geht
er über die Export-Test-Brücke nach `SPEC-049` Punkt 7 (`export_test.go`) oder wird
gegen die exportierte Schnittstelle umgeschrieben. Die Fälle und ihre Prüfungen
bleiben erhalten; der Slice ist ein Umbau der Tests, keine neue Prüfung.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-05: Die Tests werden auf
Black-Box-Pakete umgestellt; am 2026-10-06, nach der Messung: Bereinigung vor dem
Gate, ohne Stufen, und die vier Umstellungs-Slices beheben dabei auch die übrigen
Befunde ihrer Testdateien
(`slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
`slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`). Geschnitten ist
nach Paketgruppe, je in einer Schicht des Hexagons, damit jeder Schnitt einzeln
lieferbar und in einer Review-Sitzung prüfbar bleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Tests der übrigen Pakete — sie übernehmen `slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`, `slice-harness-blackbox-pgwire`.
- Neue Exporte im Produkt-Code nur für Tests — sie verbreiterten die Schnittstelle des
  Pakets für einen Zweck außerhalb des Produkts; der Weg für unexportierte Teile ist
  die Brücke in einer `_test.go`-Datei (§6). Braucht ein Test mehr, ist das ein Befund
  für den Architect, keine Entscheidung im Umbau (§4).
- Neue Fälle oder geänderte Erwartungen — ein anderer Vorgang; der Umbau soll an
  derselben Testliste messbar sein (§2). Die Schwelle der Testabdeckung übernimmt
  `slice-harness-coverage`, gemessen erst nach allen vier Umstellungs-Slices.
- Die Abdeckungs-Deklaration von Teil 3 von LH-QA-07 — übernimmt
  `slice-harness-abdeckung-gate` nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md); das Abdeckungs-Skript kennt die
  Nachweisart Gate vor ihm nicht.
- Das Scharfschalten am Gate und der White-Box-Fall für diese Pfade in der
  Lint-Gegenprobe — übernimmt `slice-harness-lint`: Vor dem Gate gibt es weder eine
  Stufe noch eine Gegenprobe (Entscheidung 5); dieser Slice belegt sein Ergebnis mit
  dem Werkzeug `make lint`.
- Befunde im Produkt-Code dieser Pakete — übernimmt `slice-lint-bestand-driving` nach den
  Umstellungs-Slices.
- Produkt-Verhalten, Spezifikation, Lastenheft — Schicht-Abgrenzung: Der Slice ändert
  Testdateien; `.golangci.yml` und die Gegenprobe des Lint-Gates ändert er nicht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Die Testdateien unter `internal/adapters/driving/cli`, `internal/bootstrap` und `test/integration` gehören zu `_test`-Paketen; die Liste der Tests
      (`go test -list .` je Paket, im gepinnten Go-Image) ist vor und nach dem Umbau gleich, keine Prüfung ist
      entfallen, die Abdeckungs-Deklarationen nach [ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) stehen unverändert und
      `make abdeckung-check` ist grün. Unexportierte Teile erreicht ein Test nur über
      die Brücke nach `SPEC-049` Punkt 7. — Verifikation, Abschnitt 1 Punkt 1 und
      Abschnitt 2: `cli_test.go` ist `package cli_test`, `bootstrap_test.go`
      `package bootstrap_test`, die sieben Dateien unter `test/integration`
      `package integration_test`; im Paket des Codes liegt je nur `export_test.go`.
      `go test -list .` an `f28a2df` und `1e4381b`: `cli` 19, `bootstrap` 12,
      `integration` 39, `diff` je leer, dieselben 31 Zeilen `--- PASS` der Unit-Tests; nach
      Normierung unterscheiden sich die Unit-Testdateien nur in der Import-Zeile. 58
      Deklarationen Wort für Wort und in der Zuordnung gleich, `make abdeckung-check`
      Exit 0. Die Brücken halten Punkt 7 in der Fassung `f0ea7be`.
- [x] `make lint` meldet in den Testdateien unter `internal/adapters/driving/cli`, `internal/bootstrap` und `test/integration` keinen Befund: Die
      22 Befunde der Messung sind behoben, ohne `//nolint` und ohne Änderung an
      `.golangci.yml`; die dauerhaften Ausnahmen nach `SPEC-049` Punkt 8 blendet das
      Profil aus. Die Zeilen der Ausgabe unter diesen Pfaden vor und nach dem Umbau
      stehen im Bericht; Befunde im Produkt-Code bleiben für `slice-lint-bestand-driving`.
      Mit diesem Slice hat keine Testdatei des Moduls mehr einen Befund. Den Nachweis von
      Teil 3 (Lage der Unit-Tests) von LH-QA-07 tragen die Fälle der Lint-Gegenprobe aus
      `slice-harness-lint`; deklariert wird er nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) in
      `slice-harness-abdeckung-gate`. — Verifikation, Abschnitt 1 Punkt 2: die 22 Zeilen
      vorher in §7 *Belege des Implementers* sortiert byte-gleich nachgemessen
      (`testpackage` 9, dazu `errcheck` 9, `gochecknoglobals` 2, `revive` 1,
      `staticcheck` 1); 22 → 0 in `*_test.go`, modulweit 54 → 32 ohne neuen Befund
      (`comm -13` leer); die 2 Produkt-Zeilen unter diesen Pfaden an beiden Ständen gleich;
      kein `//nolint`, `.golangci.yml` unverändert, kein neues `_ =` (je 30
      Blank-Zuweisungen). In keiner Testdatei des Moduls steht noch ein Befund; alle 32
      liegen im Produkt-Code.
- [x] `make gates` grün. — an `0c1e9f7` (Go-Code gleich `1e4381b`; Verifikation,
      Abschnitt 1 Punkt 3 und Abschnitt 5); `3688bea` und diese Closure ändern nur
      Berichte, Pläne, Roadmap und das Register, `make docs-check` und `make kopf-check`
      grün am Stand dieser Closure.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      — `docs/reviews/2026-10-07-review-slice-harness-blackbox-einstieg.md` (bis
      `1e4381b`; F-454 LOW, F-455 bis F-458 INFO).
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
| `internal/adapters/driving/cli/cli_test.go` | refactor | `package cli_test` |
| `internal/bootstrap/bootstrap_test.go` | refactor | `package bootstrap_test` |
| `test/integration/*_test.go` | refactor | Paketname `integration` ohne `_test`-Suffix in einem Verzeichnis nur aus Testdateien (§6) |
| `export_test.go` je Paket, nur wo nötig | neu | Brücke zu unexportierten Teilen nach `SPEC-049` Punkt 7 |
| dieselben Testdateien | update | übrige Befunde aus `make lint` (dazu `errcheck`, `gochecknoglobals`, `revive` und `staticcheck`, zusammen 13) in den Dateien, die der Umbau ohnehin umschreibt |
| `Dockerfile` | unverändert | Stufe `integration` baut das Testbinary mit `package integration_test` unverändert (`go test -c -tags integration ./test/integration`, §7) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-blackbox-pgwire` liegt in `done/` (WIP-Limit 1,
Reihenfolge nach Entscheidung des Nutzers vom 2026-10-05 und 2026-10-06:
`slice-harness-lint-werkzeug`, die vier Umstellungs-Slices,
`slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`, `slice-harness-lint`,
`slice-harness-abdeckung-gate`, Coverage, Mutation). Gemessen wird mit dem Werkzeug
`make lint` aus `slice-harness-lint-werkzeug`. Vor dem ersten Code-Commit prüft der
Architect §6 gegen `SPEC-049` (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Umbau ist nicht in einer
  Review-Sitzung prüfbar, etwa weil mehr als eine Handvoll Tests umgeschrieben statt
  umgestellt werden muss; dann je Paket ein eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): Ein Test erreicht einen Teil, den er
  prüfen muss, weder über die exportierte Schnittstelle noch über die Brücke, die die
  ADR zulässt, oder ein Befund in einer Testdatei lässt sich nur mit einer neuen
  Ausnahme beheben; dann zuerst eine Entscheidung des Architect (Brücke erweitern,
  Schnittstelle des Pakets ändern, Ausnahme mit dauerhaftem Grund).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, `make lint` ohne Befund in den Testdateien dieser
Pfade,
Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — entschieden vom Architect am 2026-10-06 in
`SPEC-049` Punkt 7 und 8 und [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 und 5; was dort nicht
steht, gibt der Implementer an den Architect zurück.

- **Export-Test-Brücke** — entschieden (`SPEC-049` Punkt 7): einzige Datei
  `export_test.go` im Paket des Codes, nur Typ-Aliase, Konstanten und Funktionen oder
  Methoden, die an Unexportiertes weiterreichen; Zustand nur an einem übergebenen Wert,
  nie auf Paketebene (eine Frist für einen Test verkürzt sie nur an einem übergebenen
  Wert). Eine Variable oder eine Funktion `Test…` darin ist ein Befund. Einen Wert
  eines unexportierten Typs legt nur der Produkt-Code an, weder Brücke noch Test, auch
  nicht lokal oder über einen Alias; was nur an ihm zu sehen ist, prüft der Test über
  die exportierte Schnittstelle (Architect 2026-10-06, F-441/F-442 in
  `slice-harness-blackbox-kern`).
- **Übrige Befunde der Testdateien** — `contextcheck`: Kontext aus `t.Context()`, in einer Funktion für `t.Cleanup` aus `context.WithoutCancel(t.Context())` (`SPEC-049` Punkt 5); `errcheck` an `Close` einer Datei oder Datenbank: den Fehler prüfen (etwa in `t.Cleanup` mit `t.Error`); `gochecknoglobals`: Testdaten in Funktionen oder als Konstanten; `revive`, `staticcheck`: nach der Meldung. Kein `_ =` vor einem Fehler, den
  `errcheck` meldet, und keine Ausnahme, die nur Bestand aussetzt (Entscheidung 5).
- **White-Box-Zugriffe im Bestand** — per AST gemessen am Stand `f28a2df` (§7): in `cli`
  die Konstanten `envFailOnUnconsumed` und `envLogLevel`, in `bootstrap` die Funktionen
  `logger` und `fail`, in `test/integration` keine (es prüft das Binary als Prozess); kein
  Wert eines unexportierten Typs. Alle vier gehen über die Brücke (Konstanten und
  Funktionen, die weiterreichen); keiner ist Befund. Der Namensabgleich per Suche am Stand
  `ce50a10` nannte `replay` in `bootstrap`, ein Fehltreffer. `test/integration` hat keinen
  Produkt-Code; `testpackage` meldet dort je Testdatei (7 in der Messung), und
  `package integration_test` baut mit `go test -c -tags integration ./test/integration`
  (Dockerfile, Stufe `integration`) unverändert (§7).
- **Fakes der Ports** — die Ports sind exportiert (`internal/hexagon/ports/...`); ein
  Fake in einem `_test`-Paket implementiert sie unverändert. Ein Fake, der auf
  unexportierte Felder eines Produkt-Typs greift, fällt unter die Brücke.
- **Mutationstests auf unexportierte Teile** — eine Mutation, die bisher ein Test auf
  eine unexportierte Funktion fing, muss auch nach dem Umbau fangen; sonst fehlt eine
  Prüfung, obwohl die Testliste gleich ist.
- **`test/integration` ohne Produkt-Code** — `SPEC-049` Punkt 8 lässt für
  `testpackage` keine Ausnahme zu; der Weg ist die Umbenennung in
  `package integration_test`, gebaut weiter mit
  `go test -c -tags integration ./test/integration` (Stufe `integration`). Baut das
  nicht, geht der Befund an den Architect (§4).

**Risiken:**

- **Prüfung geht im Umbau verloren** — ein Test bleibt in der Liste, prüft aber
  weniger, weil ein unexportierter Vergleich wegfiel
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). Je umgeschriebener (nicht nur
  umgestellter) Test eine Mutation, die er weiter fängt. — **Ausgang:** entfallen: Kein
  Test ist umgeschrieben, alle 70 sind umgestellt. Nach Normierung (Paketname, Präfix,
  Brückennamen) unterscheiden sich die Unit-Testdateien nur in der Import-Zeile, in
  `test/integration` sind nur Paketname und die Befundstellen geändert, die Testliste ist
  gleich (Verifikation, Abschnitt 1 Punkt 1). Die neun nachgefahrenen Mutationen aus §7
  sind rot aus dem genannten Grund (Verifikation, Abschnitt 3), ebenso R2, R3 und R1i des
  Reviews. Die grünen Mutanten `envFailOnUnconsumed` (V5, R1) und PGR-E3003 an der
  Rückgabe nach `formVorpruefung` (P2) waren an `f28a2df` ebenso grün; beide fängt ein
  Test in einem anderen Paket (F-456, F-457).
- **Abdeckung sinkt** — Black-Box-Tests erreichen unexportierte Pfade seltener; die
  Zahl misst erst `slice-harness-coverage` nach allen vier Umstellungs-Slices.
  — **Ausgang:** entfallen: Die Testkörper sind unverändert und rufen über die Brücke
  dieselben Funktionen und Konstanten (`logger`, `fail`, `envFailOnUnconsumed`,
  `envLogLevel`); `test/integration` prüft das Binary als Prozess und erreichte nie
  Unexportiertes. Der erreichte Code kann damit nicht sinken (Verifikation, Abschnitt 6).
  Die Zahl für alle vier Umstellungs-Slices misst `slice-harness-coverage` nach seinem
  eigenen §1, nicht als Ausgang dieses Risikos.
- **Befund verdeckt statt behoben** — `_ =` vor einem ungeprüften Fehler, eine Globale
  als Funktion mit demselben geteilten Zustand: `make lint` ist grün, der Test nicht
  besser. Review prüft die Form, das Werkzeug nur die Zahl (Grenze von `SPEC-049`).
  — **Ausgang:** entfallen: Je 30 Blank-Zuweisungen an beiden Ständen, `diff` leer; kein
  `//nolint`, `.golangci.yml` unverändert (Verifikation, Abschnitt 1 Punkt 2). Die neun
  `Close`-Fehler gehen an `t.Error`, das `Close` steht nach LIFO weiter vor dem `cancel`
  (F-458). `lebendTexte()` und `vtTexte()` liefern je Aufruf ein neues Slice-Literal, kein
  geteilter Zustand (Review, Schwerpunkt 3).

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

- **Was hat funktioniert:** Der Schnitt hielt: drei Liefer-Punkte, neun Testdateien, zwei
  Brücken, kein Produkt-Code, kein umgeschriebener Test, keine Rückführung aus §4 und keine
  Entscheidung des Architect während der Arbeit. Die Brücken-Regel aus `SPEC-049` Punkt 7
  in der Fassung `f0ea7be` trug ohne Nachfrage. Alle vier Zugriffe sind Konstanten oder
  weiterreichende Funktionen mit Ein- und Ausgaben exportierter Typen. Die Belege in §7
  erreichten Review und Verifikation: AST-Messung, Testliste, Lint-Zeilen vorher und
  nachher, 17 Zeilen der Mutationstabelle und die grünen Mutanten unter *Ohne roten Test*.
  Die Verifikation fuhr neun Zeilen der Tabelle nach und stellte die zwei grünen Mutanten
  nach, beide sind richtig eingeordnet. Die Einordnung nach
  `.claude/commands/implement-slice.md` Schritt 19 trug für beide (F-456, F-457). Damit ist
  die Reihe der vier Umstellungs-Slices abgeschlossen. In keiner Testdatei des Moduls meldet
  `make lint` noch einen Befund, alle 32 liegen im Produkt-Code. Die Vorbedingung für
  `slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving` und `slice-harness-lint`
  ist erfüllt (Verifikation, Abschnitt 1).
- **Was ging anders als geplant:** Drei Punkte, keiner verlangt Nacharbeit am Code.
  1. Vor dem Code stand die White-Box-Liste in §6 nur als Namenssuche am Stand `ce50a10`
     (`envFailOnUnconsumed`, dazu `replay` als Fehltreffer). Gemessen per AST wurde erst im
     Code-Commit `1e4381b`, mit drei weiteren Zugriffen (`envLogLevel`, `logger`, `fail`).
     Alle vier fallen in vorab entschiedene Klassen, deshalb ist keine Randform im Code
     entschieden (F-455). Dasselbe Muster stand schon in `slice-harness-blackbox-kern` und
     `slice-harness-blackbox-pgwire` (Register unten).
  2. Die Commit-Message von `1e4381b` nennt eine Struktur-Kennung, gegen `AGENTS.md` §5
     Regel 1 (F-454). Der Commit ist gepusht und bleibt so, die Historie wird nicht
     umgeschrieben.
  3. In §7 *Ohne roten Test* war die Einordnung des grünen PGR-E3003-Mutanten ungenau
     formuliert. Die Fakten stimmten. Die Closure hat den Satz nach dem Vorschlag der
     Verifikation präzisiert (F-457, Verifikation Abschnitt 4).

  Aus dem Bestand kam V-88 dazu. Unter einem roten Lauf kann
  `TestE2ERecordExtendedSigtermBeimPipelining` bis zum Zeitlimit hängen, weil zweimal
  `Wait` auf demselben `exec.Cmd` läuft (Folge-Slices).
  - **Summary-Zeilen:** Review
    `docs/reviews/2026-10-07-review-slice-harness-blackbox-einstieg.md`: „0 HIGH · 0 MEDIUM
    · 1 LOW (F-454: Commit-Message `1e4381b` nennt die Struktur-ID `SPEC-049`) · 4 INFO
    (F-455 White-Box-Liste erst im Code-Commit gemessen, alle Zugriffe in vorab
    entschiedenen Klassen; F-456 grüner Mutant `envFailOnUnconsumed` richtig eingeordnet,
    Namens-Asymmetrie aus dem Bestand; F-457 PGR-E3003-Einordnung stimmt, Formulierung
    ungenau; F-458 `Close`-Prüfungen korrekt, Reihenfolge vor `cancel` unverändert).
    Wiederkehrende Klasse: „Struktur-ID in der Commit-Message“ (zweites Auftreten nach
    F-336).“ Verifikation
    `docs/reviews/2026-10-07-verifikation-slice-harness-blackbox-einstieg.md`, Urteil: Die
    drei Liefer-Punkte sind erfüllt und selbst belegt. Die Brücken halten Punkt 7, kein
    Produkt-Code ist geändert. Die neun nachgefahrenen Mutationen sind rot aus dem genannten
    Grund, die grünen Mutanten sind tragfähig eingeordnet. F-456 bis F-458 sind bestätigt,
    bei F-457 stimmen die Fakten und die Formulierung war ungenau. Kein Befund blockiert die
    Closure. V-88 ist ein Hinweis aus dem Bestand. `make gates` grün an `0c1e9f7`.
- **Steering-Loop-Eintrag:** Benannte Spec-Lücke: Ob der Commit-Träger
  `.githooks/commit-msg` eine Struktur-Kennung (`SPEC-*`, `ARC-*`) in der Message ablehnt,
  hat keine Entscheidungsstelle. Die Spezifikation führt für den Träger keine Festlegung,
  [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md) und
  [ADR-0029](../../adr/0029-benannte-welle-kennungen-im-commit-hook.md) tragen
  `Schärft: —` (`BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum`, 1×). Der Träger
  prüft nur, ob eine zugelassene Kennung dasteht, und für `AGENTS.md` §5 Regel 1 gibt es
  keinen Sensor. Die Regel hält deshalb nur das Lesen von Review und Verifikation, und beide
  fanden sie erst am fertigen Commit (F-454 hier, V-61 in
  `slice-replay-semantik-meldungscodes`).
  Nach der Entscheidung des Nutzers vom 2026-10-06 kommt die Festlegung eines älteren
  Werkzeugs in die Spezifikation, sobald für es eine Randform zu entscheiden ist. Die Frage
  „Struktur-Kennung ablehnen oder annehmen“ ist eine solche Randform. Ein Sensor im Träger
  setzt zuerst diese Stelle voraus. Adresse: `BEO-REPO/commit-nennt-struktur-kennung`
  (2×, offen, unter der Schwelle). Retirement-Checks: `AGENTS.md` §3.9 (seit
  welle-walking-skeleton) ist nicht wieder aufgetreten (Review, Negativbefund zum Plan),
  die Regel bleibt. §3.10 entfällt, kein neuer Vertrag, die Mutationen in §7 sind trotzdem
  je Zusage gefahren. §3.11 ohne Befund. §3.12 und die Randform-Rückgabe sind ohne Befund
  im Code-Commit (F-455 INFO, alle Zugriffe in vorab entschiedenen Klassen) und bleiben.
  `.claude/commands/implement-slice.md` Schritt 19 (seit slice-harness-blackbox-driven)
  trug für beide grünen Mutanten und bleibt.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `3688bea` (Zähler =
  Dateien unter `evidence/`).
  - `BEO-REPO/commit-nennt-struktur-kennung`: **Beleg**, 1× → 2× (F-454), offen. F-336
    (`c8ebf08` in `slice-extended-query-replay`) ist ein früheres Auftreten derselben
    Klasse. Bei dessen Closure gab es den Eintrag noch nicht, und nachgetragen ist der
    Beleg nicht. Zählte er, stünde der Eintrag bei 3×.
  - `BEO-REPO/white-box-liste-vor-code-nur-namenssuche`: **neu, mit drei Belegen**,
    `slice-harness-blackbox-kern`, `slice-harness-blackbox-pgwire` und dieser Slice
    (F-455; in kern und pgwire aus deren Review und Verifikation belegt). **Der Eintrag
    erreicht die Schwelle 3×.** Ausgang **gestrichen**, die Begründung steht in seinem
    `state.md`: Die Ursache ist weggefallen. Mit diesem Slice gehört jede Testdatei des
    Moduls zu einem `_test`-Paket. Was ein Test an Unexportiertem erreicht, steht je Paket
    in genau einer `export_test.go`, und der Compiler erzwingt das. `slice-harness-lint`
    schaltet `testpackage` am Gate scharf, und ein weiterer Umstellungs-Slice ist nicht
    geplant. `slice-harness-blackbox-driven` ist das Gegenbeispiel der Reihe, dort wurde
    vor dem Code gemessen.
  - Ohne Beleg: `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (bleibt 1×, keine Prüfung
    verloren, keine Gate-Regel geändert), `BEO-REPO/randform-wellenlos-ohne-architect-vor-code`
    (bleibt 1×, die Randformen standen vor dem Code entschieden in §6),
    `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (bleibt 3×, verkörpert; keine
    Randform im Code entschieden), `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht`
    (bleibt 3×, verkörpert; die Belege in §7 erreichten beide Prüfer),
    `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (bleibt 1×, jede Mutation lief in einem
    frischen Pfad), `BEO-REPO/plan-folgt-korrektur-nicht` (bleibt 15×; §1, §3 und §6 folgten
    dem Diff, F-457 ist eine ungenaue Formulierung bei richtigen Fakten),
    `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (bleibt 16×, Review Hard Rule 3.11
    ohne Befund), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (bleibt 13×, kein neuer
    Vertrag), `BEO-REPO/werkzeug-festlegung-ausserhalb-technik-stratum` (bleibt 1×; der
    Slice ändert kein Werkzeug, der Lerneintrag nennt ihn nur).

  Einmalig und nicht eingetragen: F-456 (Namens-Asymmetrie aus dem Bestand, der
  Integrationstest hält den Namen als Literal), F-457 (Formulierung, präzisiert), F-458
  (keine Aktion) und V-88 (Fehler im Testgeschirr aus dem Bestand; Adresse unter
  *Folge-Slices*). Mit diesem Slice erreicht **ein** Eintrag die Schwelle 3× neu,
  `white-box-liste-vor-code-nur-namenssuche`, mit Ausgang gestrichen. Über der Schwelle
  stehen sonst nur verkörperte Einträge (`implementer-bericht-erreicht-pruefer-nicht` 3×,
  `randform-im-code-entschieden-dann-zurueckgegeben` 3×,
  `spec-randform-erst-im-review-entschieden` 11×,
  `negativtests-fehlen-bei-neuem-vertrag` 13×, `plan-folgt-korrektur-nicht` 15×,
  `zusage-im-kommentar-weiter-als-pruefung` 16×).
- **Folge-Slices:** `slice-v1-abschluss-betrieb` (V-88 als Risiko in §6 mit
  Herkunfts-Anker). Er nimmt an: Er liegt in `open/`, und DoD-Punkt 1 verlangt Tests mit
  Signal je Modus über dieselben Helfer `startProzess` und `stop`. Ihr Fehlerfall, ein
  Prozess, der nach dem Signal nicht endet, ist genau der, unter dem der Test hängt. §1
  schließt nur Container-Image und Meldungscodes aus. Nicht `slice-harness-mutation`, sein
  §1 schließt Mutationen in den Integrationstests aus. Weiter:
  `slice-lint-bestand-driving` (die zwei Befunde im Produkt-Code unter diesen Pfaden,
  `bootstrap.go:71:26` und `cli.go:243:7`, Zeilen in §7 *Nachher*),
  `slice-harness-abdeckung-gate` (Deklaration von Teil 3 von
  [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)),
  `slice-harness-lint` (Scharfschalten am Gate, White-Box-Fall dieser Pfade in der
  Gegenprobe) und `slice-harness-coverage` (Zahl der Abdeckung nach allen vier
  Umstellungs-Slices). Als nächster in der Reihe folgt `slice-lint-bestand-kern-driven`
  (§4 *Start*: dieser Slice in `done/`; die Reihenfolge in §4 von
  `slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving` und `slice-harness-lint`
  stimmt überein).
- **Risiken aus §6:** drei, alle **entfallen**, mit Begründung: *Prüfung geht im Umbau
  verloren* (kein Test umgeschrieben, Unit-Testdateien nach Normierung bis auf die
  Import-Zeile gleich, Mutationen rot, die grünen Mutanten waren vorher ebenso grün),
  *Abdeckung sinkt* (dieselben Funktionen über die Brücke, `test/integration` ohne
  Produkt-Code), *Befund verdeckt statt behoben* (Blank-Zuweisungen gleich, kein
  `//nolint`, `Close`-Fehler an `t.Error`, Testdaten je Aufruf neu). Die Randformen in §6
  sind Entscheidungen, keine Risiken.
- **Drei Paarungen:** Anker: Der Lerneintrag ist eine benannte Spec-Lücke, keine
  verkörperte Regel, deshalb steht kein `liegt in`. Der Herkunfts-Anker
  `seit slice-harness-blackbox-einstieg` steht an der Adresse von V-88,
  `grep -n "seit slice-harness-blackbox-einstieg" docs/plan/planning/done/slice-v1-abschluss-betrieb.md`
  findet ihn in §6. Folge-Slice: `slice-v1-abschluss-betrieb`,
  `slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`,
  `slice-harness-abdeckung-gate` und `slice-harness-coverage` liegen in `open/`,
  `slice-harness-lint` in `next/`. Register: `BEO-REPO/commit-nennt-struktur-kennung` trägt
  `evidence/slice-harness-blackbox-einstieg.md`, `BEO-REPO/white-box-liste-vor-code-nur-namenssuche`
  trägt drei Dateien und `state.md` mit Begründung. Die übrigen genannten Einträge bestehen
  als Verzeichnis, jedes mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft
  erneut.
- **Belege:** Review `docs/reviews/2026-10-07-review-slice-harness-blackbox-einstieg.md`
  (bis `1e4381b`; F-454 bis F-458), Verifikation
  `docs/reviews/2026-10-07-verifikation-slice-harness-blackbox-einstieg.md` (bis `0c1e9f7`,
  `make gates` an `0c1e9f7`; V-88), Entscheidung
  [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) unverändert.
  Validierung: n/a, der Slice ändert Tests, kein End-Nutzer-Verhalten.

**Belege des Implementers** (Arbeitsbaum auf `f28a2df` mit dem Diff des Commits, der
diesen Abschnitt anlegt):

*White-Box-Zugriffe, per AST gemessen* (go/types im Image der Stufe `deps`, ohne Netz;
gezählt sind Bezeichner in Testdateien der Pakete `cli`, `bootstrap` und `integration`,
die auf ein unexportiertes Objekt einer Produkt-Datei desselben Pakets zeigen, dazu
Literale, `var`, `new` und `:=` unexportierter Typen in Testdateien). Vorher
(`f28a2df`), je Zugriff mit Zuordnung:

| Paket | Zugriff | Stellen | Zuordnung |
|---|---|---|---|
| `cli` | Konstante `envFailOnUnconsumed` | 15 | Brücke, Konstante `EnvFailOnUnconsumed` |
| `cli` | Konstante `envLogLevel` | 6 | Brücke, Konstante `EnvLogLevel` |
| `bootstrap` | Funktion `logger` | 2 | Brücke `Logger` (reicht an `logger` weiter; Writer und Stufe stellt der Test, Ergebnis ist `*slog.Logger`) |
| `bootstrap` | Funktion `fail` | 1 | Brücke `Fail` (reicht an `fail` weiter; Writer und Fehler stellt der Test, Ergebnis ist `int`) |
| `integration` | — | 0 | kein Produkt-Code im Verzeichnis |

Werte unexportierter Typen in Testdateien vorher und nachher 0. Nachher stehen die vier
Objekte je einmal in `export_test.go` ihres Pakets und sonst in keiner Testdatei; kein
Zugriff blieb Befund, keiner wurde gegen die Schnittstelle umgeschrieben. Fehltreffer der
Suche in §6: `replay` in `bootstrap` (kein Test greift darauf zu).

Die Brücken: `internal/adapters/driving/cli/export_test.go` mit zwei Konstanten,
`internal/bootstrap/export_test.go` mit zwei Funktionen, die nur weiterreichen; keine
Variable, kein Zustand, kein Typ-Alias. Kein Produkt-Code geändert; `Dockerfile`
unverändert, die Stufe `integration` baut `package integration_test` mit
`go test -c -tags integration` (Lauf unten).

*Testliste* (`go test -list .` je Paket im Image der Stufe `deps`, für
`./test/integration` mit `-tags integration` und `PGR_OHNE_POSTGRES=1`, damit `TestMain`
ohne Upstream listet): vorher und nachher `cli` 19, `bootstrap` 12, `integration` 39
Tests, `diff` leer. Die Abdeckungs-Deklarationen stehen unverändert: Der Diff ändert
genau drei Kommentare, die Doc-Kommentare von `lebendTexte` und `vtTexte` und den
Paket-Kommentar (`Package integration_test`), keiner davon eine Deklaration
`Abdeckung:` oder ihre Fortsetzung; `make abdeckung-check` grün.

*`make lint`.* Vorher (`f28a2df`, Exit 2): golangci-lint `54 issues:` im Modul, davon 22
in Testdateien, alle unter den Pfaden dieses Slice; je Paket `cli` 2 (1 in Testdateien),
`bootstrap` 2 (1), `test/integration` 20 (20):

```text
test/integration/extended_e2e_test.go:367:16: Error return value of `pc.Close` is not checked (errcheck)
test/integration/extended_replay_e2e_test.go:35:18: Error return value of `conn.Close` is not checked (errcheck)
test/integration/extended_replay_e2e_test.go:135:18: Error return value of `conn.Close` is not checked (errcheck)
test/integration/extended_replay_e2e_test.go:245:18: Error return value of `conn.Close` is not checked (errcheck)
test/integration/lebendpruefung_e2e_test.go:65:16: Error return value of `db.Close` is not checked (errcheck)
test/integration/lebendpruefung_e2e_test.go:208:18: Error return value of `conn.Close` is not checked (errcheck)
test/integration/record_e2e_test.go:233:18: Error return value of `conn.Close` is not checked (errcheck)
test/integration/replay_e2e_test.go:41:18: Error return value of `conn.Close` is not checked (errcheck)
test/integration/unverbraucht_e2e_test.go:219:18: Error return value of `conn.Close` is not checked (errcheck)
test/integration/lebendpruefung_e2e_test.go:148:5: lebendTexte is a global variable (gochecknoglobals)
test/integration/lebendpruefung_e2e_test.go:155:5: vtTexte is a global variable (gochecknoglobals)
test/integration/replay_e2e_test.go:134:6: var-naming: don't use underscores in Go names; func exec_ should be exec (revive)
test/integration/extended_e2e_test.go:376:6: SA4000: identical expressions on the left and right side of the '||' operator (staticcheck)
internal/adapters/driving/cli/cli_test.go:1:9: package should be `cli_test` instead of `cli` (testpackage)
internal/bootstrap/bootstrap_test.go:1:9: package should be `bootstrap_test` instead of `bootstrap` (testpackage)
test/integration/einfach_e2e_test.go:3:9: package should be `integration_test` instead of `integration` (testpackage)
test/integration/extended_e2e_test.go:3:9: package should be `integration_test` instead of `integration` (testpackage)
test/integration/extended_replay_e2e_test.go:3:9: package should be `integration_test` instead of `integration` (testpackage)
test/integration/lebendpruefung_e2e_test.go:3:9: package should be `integration_test` instead of `integration` (testpackage)
test/integration/record_e2e_test.go:10:9: package should be `integration_test` instead of `integration` (testpackage)
test/integration/replay_e2e_test.go:3:9: package should be `integration_test` instead of `integration` (testpackage)
test/integration/unverbraucht_e2e_test.go:3:9: package should be `integration_test` instead of `integration` (testpackage)
```

Nachher (Arbeitsbaum, Exit 2): `32 issues:` im Modul, 22 weniger, keiner neu (`diff` der
sortierten Befundzeilen: 22 entfallen, 0 hinzu); **in keiner Testdatei des Moduls ein
Befund** und keine `lint:`-Zeile der eigenen Prüfungen (auch nicht an den beiden
`export_test.go`). Unter den Pfaden dieses Slice bleiben 2, beide im Produkt-Code und
dieselben Zeilen wie vorher, für `slice-lint-bestand-driving`:

```text
internal/bootstrap/bootstrap.go:71:26: Non-inherited new context, use function like `context.WithXXX` instead (contextcheck)
internal/adapters/driving/cli/cli.go:243:7: unused-receiver: method receiver 'w' is not referenced in method's body, consider removing or renaming it as _ (revive)
```

Behoben ohne `//nolint`, ohne neues `_ =` und ohne Änderung an `.golangci.yml`:
`testpackage` durch `package cli_test`, `package bootstrap_test` und
`package integration_test`; `errcheck` an `Close` (neun Stellen, `pgconn`, `pgx` und
`database/sql`) durch ein `defer func()`, das den Fehler mit `t.Error` meldet. Das `defer`
bleibt, statt nach `t.Cleanup` zu wandern: Die Verbindung schließt weiter am Ende des
Helfers, vor dem `cancel` des Kontexts, und die Tests, die danach auf das Ende der Session
warten, sehen denselben Ablauf. `gochecknoglobals` durch die Funktionen `lebendTexte()`
und `vtTexte()`, die je Aufruf eine neue Liste liefern; `revive` `var-naming` durch den
Namen `fuehreAus` statt `exec_`; `staticcheck` SA4000 durch zwei Anweisungen
`if senden() != nil { return }`, dieselben zwei Aufrufe in derselben Reihenfolge.

*Umgeschrieben:* keiner der 70 Tests gegen die Schnittstelle. Umgestellt sind alle
(Paketname, Brücke, Präfix `cli.` und `bootstrap.`); geändert über die Umstellung hinaus
sind nur die Stellen der Befunde oben, in den Tests der Tabelle unten.

*Mutationen* (`AGENTS.md` §3.10, Risiko *Prüfung geht im Umbau verloren*): je Mutant ein
frischer Pfad außerhalb des Repos (`cp -r` ohne `-p` von `go.mod`, `go.sum`, `cmd/`,
`internal/`, `test/` aus dem Arbeitsbaum), Tests per Bind-Mount im Image der Stufe `deps`;
Unit-Tests mit `go test -count=1 -run`, Integrationstests mit `-tags integration` gegen
eine eigene PostgreSQL-17-Instanz (dasselbe gepinnte Image wie `make test-integration`)
in einem eigenen internen Docker-Netz, das danach entfernt wurde. Ohne Mutation ist jeder
genannte Test grün (`U0`, `E0`), mit ihr rot.

| Zusage | Mutation | roter Test |
|---|---|---|
| eine Stufe zeigt ihre und die strengeren Zeilen (über `Logger`) | `Level: slog.LevelInfo` statt `stufen[stufe]` | `TestLoggerSchwelle` |
| `time` ist Ortszeit (über `Logger`) | `ReplaceAttr` setzt die Zeit auf UTC | `TestLoggerOrtszeit` |
| je Meldung eine Zeile beim Prozessende (über `Fail`) | `fail` schreibt nur die erste Meldung | `TestFailJeKlasse` |
| Exit-Code der ersten Meldung (über `Fail`) | `fail` liefert den Exit-Code der letzten | `TestFailJeKlasse` |
| der Name `PGWIRE_RECORDER_LOG_LEVEL` (Konstante über die Brücke) | `envLogLevel = "PGWIRE_RECORDER_LOGLEVEL"` | `TestRunStartfehlerJeStufe` (`bootstrap`, Literal), `TestParseLogLevelHilfe` (`cli`) |
| der Name `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` (Konstante über die Brücke) | `envFailOnUnconsumed = "PGWIRE_RECORDER_FAILONUNCONSUMED"` | `TestE2EReplayNichtVerbraucht` (Literal); in `cli` grün, unten eingeordnet |
| `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` wird gelesen | `Parse` liest die Variable nicht | `TestParseFailOnUnconsumedUmgebung`, `TestParseFailOnUnconsumedUmgebungNebenOption` |
| `PGWIRE_RECORDER_LOG_LEVEL` wird gelesen | `Parse` liest die Variable nicht | `TestParseLogLevelUmgebung` |
| nach SIGTERM beginnt keine neue Interaktion, auch beim Pipelinen (SA4000 umgeschrieben, `Close` geprüft) | `ClientMessage` ohne die Sperre `herunterfahren` | `TestE2ERecordExtendedSigtermBeimPipelining` (endet nicht binnen 5 s) |
| verschachtelte Blockkommentare sind Lebendprüfungen (`lebendTexte()`) | `/*` öffnet nur auf Tiefe 0 eine Ebene | `TestE2EReplayLebendpruefungWiePostgres` (PGR-E5001 an `"/* a /* b */ c */"`) |
| ein Zeilenkommentar reicht bis zum Textende (`lebendTexte()`, `sqlAblauf` und Pool mit geprüftem `Close`) | `--` ohne Zeilenende ist keine Lebendprüfung | `TestE2EReplayLebendpruefungPgxpool`, `…DatabaseSQL`, `…WiePostgres` |
| `\v` ist ab Version 17 Leerraum (`vtTexte()`) | `vtAbVersion = 18` | `TestE2EReplayLebendpruefungWiePostgres` (PGR-E5001 an `"\v"`) |
| eine beschädigte Aufzeichnung ist PGR-E3003 (über `fuehreAus`) | „Liste der Sessions fehlt“ als PGR-E3001 | `TestE2EReplayBeschaedigt` |
| ein Startfehler ist PGR-E4001 (über `fuehreAus`) | `Listen` meldet PGR-E4000 | `TestE2EReplayNichtVerbrauchtStartfehler` |
| die Interaktion einer offenen Verbindung bleibt beim Beenden (`Close` geprüft) | `CloseSession` verwirft Sessions mit `EndShutdown` | `TestE2ERecordBeendenMitOffenerVerbindung` |
| eine beim Herunterfahren nicht verbrauchte Session ist PGR-E5002 (`Close` geprüft) | `CloseConnection` meldet nur ohne verbrauchte Interaktion | `TestE2EReplayNichtVerbrauchtHerunterfahren` |
| Replay liefert die Zeilen einer Extended-Interaktion (`extendedAblauf`, `pipelineAblauf`, Sigterm-Test, je `Close` geprüft) | `ClientMessage` lässt `DataRow` weg | `TestE2EReplayExtendedPgx`, `TestE2EReplayExtendedPipeline`, `TestE2EReplayExtendedSigtermMittenInFolge` |
| Replay liefert die Zeilen einer einfachen Anfrage (`beobachten`, `Close` geprüft) | `Query` lässt `DataRow` weg | `TestE2EReplaySelect1` |

*Ohne roten Test*, eingeordnet nach `.claude/commands/implement-slice.md` Schritt 19:

- `envFailOnUnconsumed = "PGWIRE_RECORDER_FAILONUNCONSUMED"` bleibt in allen Tests von
  `cli` grün, vorher (`git archive` von `f28a2df` in einem frischen Pfad) wie nachher:
  Die Tests nennen dieselbe Konstante, vorher direkt, nachher über die Brücke. Der Mutant
  ändert das Verhalten und ist über die Schnittstelle gefangen, in einem anderen Paket:
  `TestE2EReplayNichtVerbraucht` setzt die Variable als Literal auf `true` und wird rot.
  Kein Verlust durch den Umbau.
- Dieselbe Ersetzung PGR-E3003 → PGR-E3001 an der Rückgabe „Aufzeichnung beschädigt“
  nach `formVorpruefung` blieb in `TestE2EReplayBeschaedigt` grün. Die Eingabe besteht
  `formVorpruefung` und scheitert erst in `fromDTO` (Zeile oben, rot). Die Rückgabe selbst
  fangen `TestUnmarshalExtendedFehler` und `TestVorpruefungNenntOrt` im Paket
  `recording`, mit dieser Ersetzung rot (Review R3, Verifikation P2; Fassung nach F-457).
- Für die Prüfung des `Close`-Fehlers selbst ist kein Mutant gefahren: `Close` schließt
  die Verbindung des Tests, nicht eine des Produkts, und lieferte in allen Läufen nil,
  auch in den Tests, in denen das Produkt die Verbindung vorher beendet. Einen Mutanten
  im Produkt, der hier einen Fehler erzeugt, kenne ich nicht; die Prüfung ist ein Zusatz
  am Testgeschirr, keine Zusage über das Produkt. Grenze: Ob sie bei einem Fehler meldet,
  prüft kein Test.

*Läufe:* `gofmt -l` (leer), `go vet -tags integration` der drei Pakete, `go test` von
`cli` und `bootstrap` und `go test -c -tags integration ./test/integration` im Image der
Stufe `deps`; `make test-integration` grün (37 Tests bestanden und 2 übersprungen in der ersten Phase, dazu 2 bestanden in der zweiten Phase ohne PostgreSQL); `make lint` vorher
und nachher wie oben; `make abdeckung-check` grün; `make gates` grün am Stand dieses
Commits.

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `175fd2d` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (12×, verkörpert in `AGENTS.md`
  §3.10) — der rote Fall für diese Pfade gehört zum Gate und liegt bei
  `slice-harness-lint`; dieser Slice liefert keinen neuen Vertrag.
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — ein Risiko in §6.
  `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) betrifft diesen Slice nicht
  mehr, er führt keine Gegenprobe.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — wellenlos; §4 nennt den
  Architect vor dem Code.
- `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (2×) — die Brücke ist
  genau eine Randform, die ein Umbau still entscheiden würde; sie steht in §6.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
