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

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--lint-profil-lint) (`spezifikation.md` §11: Vertrag des Gates, vom Architect vor dem Code geschrieben, die Grenze am 2026-10-07 vor dem Code um ungültiges YAML und den Cache des Builds fortgeschrieben, Punkt 8 am 2026-10-07 nach dem ersten Code-Commit um Leerraum nach dem Doppelpunkt und den Eintrag `-` allein (Rückgaben aus `221845d`, §6), die Grenze am 2026-10-08 um die zulässigen Regeln (V-100, §6); dieser Slice liefert die zweite Hälfte von Punkt 10, den Anschluss an die Gate-Kette) · `spezifikation.md` §12 (*Historie*)

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
vier Umstellungs-Slices tragen sollten. `AGENTS.md` §3.2 bekommt seinen Träger, und die
Sensor-Datei `harness/sensors/lint.md` sagt, wie ein Lauf zu lesen ist, nach dem Muster
von `kopf-check`.

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
  `Dockerfile`, die drei eigenen Prüfungen, `.dockerignore`) — Bestand bleibt bewusst
  stehen: geliefert von `slice-harness-lint-werkzeug` (in `done/`), damit jeder
  Bereinigungs-Slice sein Ergebnis mit demselben Ziel maß, das hier Gate wird
  (Entscheidung 6 in [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md)).
  Dieser Slice ändert daran nur Kommentare, die das Werkzeug als Nicht-Gate beschreiben
  (§3), kein Verhalten.
- Die Umstellung der Tests auf Black-Box-Pakete samt den übrigen 47 Befunden in
  Testdateien — Bestand: geliefert von `slice-harness-blackbox-kern`,
  `slice-harness-blackbox-driven`, `slice-harness-blackbox-pgwire` und
  `slice-harness-blackbox-einstieg` (alle in `done/`). Die Form der Export-Test-Brücke
  steht in `SPEC-049` Punkt 7.
- LH-QA-07 im Lastenheft — Bestand: geliefert von `slice-lastenheft-pruefbarkeit` (in
  `done/`).
- Eine Abdeckungs-Deklaration an der Gegenprobe — übernimmt
  `slice-harness-abdeckung-gate` nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md): Das Abdeckungs-Skript kennt die
  Nachweisart Gate vor ihm nicht. Teil 3 von LH-QA-07 gilt ab diesem Slice, weil
  `testpackage` erst mit dem Gate für alle Pfade scharf ist.
- Bereinigung des Produkt-Codes — Bestand: Die 32 Befunde nach `SPEC-049` Punkt 9
  (§6 *Bestand*) haben `slice-lint-bestand-kern-driven` (16) und
  `slice-lint-bestand-driving` (16) behoben (beide in `done/`); `make lint` meldet am
  Start keinen Befund (§4). Dieser Slice ändert keinen Produkt-Code außer dem Kommentar
  in `internal/hexagon/model/fehler.go` (§3).
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
  Historien-Zeilen schreibt nur der Architect: vor dem Code den Vertrag und die Grenze,
  nach dem ersten Code-Commit Punkt 8 auf die Rückgaben aus `221845d` und die Grenze auf
  V-100 der Verifikation (§6); dieser Slice ändert das
  Gate-Fragment, die Gegenprobe, die Doku und die Kommentare aus §3, Produkt-Code nur
  den Kommentar an `Meldungen`.
  Braucht das Gate eine Änderung am Profil oder an der Stufe, geht sie an den
  Architect zurück, weil `slice-harness-lint-werkzeug` sie geliefert hat.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Gate `make lint`: Das Fragment `harness/mk/lint.mk` hängt `lint` an
      `GATE_CHECKS` (`SPEC-049` Punkt 10, zweite Hälfte), und `make gates` ist mit ihm
      grün. Der Bestand ist dafür grün, ohne `//nolint`, und `.golangci.yml` trägt nur
      die dauerhaften Ausnahmen nach `SPEC-049` Punkt 8, je mit `Why:`, keine Regel, die
      nur Bestand aussetzt, und keine ungenutzte Regel; `testpackage` ist für jeden Pfad
      scharf. — Verifikation, Abschnitt 1 Punkt 1: `GATE_CHECKS += lint lint-gegenprobe`;
      in einer Kopie mit `.git` macht ein `//nolint` `make gates` rot (Exit 2, rot nur
      `lint` und die Gegenprobe als Folgerot, `record-gates` lief nicht); `make lint` in
      einer frischen Kopie `0 issues.` ohne `lint:`-Zeile; 16 Regeln, je mit `# Why:`,
      keine ungenutzt (`p0-grundlauf`), `.golangci.yml` im Slice unverändert;
      `skip-regexp` nur `export_test.go` (V10 und M2 rot).
- [x] Gegenprobe `make lint-gegenprobe` an `GATE_CHECKS`: je Zusage aus `SPEC-049`, die
      eine Mutation fangen kann, ein Fall, der rot wird, und zwar aus seinem Grund: Jeder
      rote Fall prüft neben dem Ausgang die Zeile, die er auslösen soll (Pfad und Linter
      bzw. Text der `lint:`-Zeile); vorweg ein Grundlauf der unveränderten Kopie mit
      Ausgang 0 und ohne `lint:`-Zeile; dazu ein Fall, dass `lint` und `lint-gegenprobe`
      an `GATE_CHECKS` hängen und `make gates` bei rotem `make lint` mit Fehlerstatus endet
      (`SPEC-049` Punkt 10). Die Fälle sind (mindestens ein Verstoß je
      aktivierter Linter-Gruppe, ein `//nolint`, ein Modul außerhalb der Liste von
      `gomodguard_v2`, ein Integrationstest hinter dem Build-Tag, eine ungenutzte Regel,
      eine Regel ohne `Why:`; nach `SPEC-049` Punkt 6, 8 und 9 zudem ein fehlendes
      Profil, ein vom Schema abgelehntes Profil, ein Profil außerhalb der festen Form
      (Einzug, Flussform, Einträge unter `rules` mit Einzug 8 und 2; bei Einzug 2 prüft
      der Fall nur `Form nicht erkannt` an der ersten Eintragszeile und Ausgang 1, nicht
      die vollständige Liste der Zeilen, weil die Grenze von `SPEC-049` für ungültiges
      YAML nur diese Zeile zusagt), eine Regel, die
      `config verify` annimmt und golangci-lint erst beim Laden ablehnt (rot ohne
      `lint:`-Zeile), die Zeile der ungenutzten Regel mit Feldfolge (ohne erkannte
      Felder offen im Kopf der Gegenprobe, §6 *Randform-Rückgaben aus `221845d`*), `// x //nolint` als Befund und `// siehe nolint` als keiner, eine
      `lint:`-Zeile allein mit Ausgang ungleich 0; für die erste Bedingung von Messmethode 3 je Paketgruppe
      der Umstellungs-Slices — Kern, Driven, PGWire, Einstieg — ein White-Box-Test, der
      rot wird; für die zweite Bedingung ein Test in `internal_test.go` im Paket des
      Codes, eine Funktion `Test…` und eine Variable in `export_test.go`), dazu ein
      grüner Fall je dauerhafter Ausnahme; je Zusage ist die Mutation gesehen
      (`AGENTS.md` §3.10). Ausdrücklich behandelt sie die Zusagen, die der Lauf von
      `slice-harness-lint-werkzeug` per Mutation nicht unterschied: `relative-path-mode:
      cfg` und Pfade mit `^`, `--network=none`, `GOFLAGS=-mod=readonly`,
      `GOTOOLCHAIN=local`, `-c` statt Default-Suche, Pin und Plattform, die untere
      Grenze von `dupl` und die Erkennung der ungenutzten Regel am Logtext — je mit
      einem Fall, der sie unterscheidet, oder benannt als offen im Kopf der Gegenprobe. Die Gegenprobe liegt unter `tools/harness/lint-gegenprobe.sh`
      und ist der Nachweis von Teil 1 (statische Analyse) und Teil 3 (Lage der
      Unit-Tests) von LH-QA-07; nach [ADR-0033](../../adr/0033-gate-nachweise-in-der-abdeckung.md) trägt sie in diesem Slice keine
      Abdeckungs-Deklaration, LH-QA-07 steht bis zu `slice-harness-abdeckung-gate` in
      keiner Abdeckungstabelle (`make abdeckung-check` grün). — Verifikation, Abschnitt 1
      Punkt 2 und Abschnitt 3: Grundlauf, eigene Zeile je rotem Fall, `p10-gates-rot` (V1
      rot), die Fälle der Aufzählung, die Zusagen aus `slice-harness-lint-werkzeug` als
      Fall oder OFFEN; F-471 bis F-473 und F-475 nachgefahren (M1, M2, M3, V13 rot); keine
      Abdeckungs-Deklaration. Die beiden Lücken aus V-96 schließt `5f9fef5`:
      `p6-leerzeichen-dann-strich` (`// /nolint`) und `p8-why-ohne-doppelpunkt` (`# Why
      kein Doppelpunkt.`), je in `sammel` mit eigener Zeile; die Mutationen V7 und V11,
      vorher grün, sind in §7 *Nachtrag zur Verifikation* je in einer frischen Kopie rot
      gesehen. Den Satz zur Zeile ohne erkannte Felder zog `d3473b9` auf §6 nach (V-97).
- [x] Doku: `AGENTS.md` §3.2 nennt den Träger (kein `//nolint`, Ausnahmen zentral in
      `.golangci.yml` mit `Why:`) statt der Platzhalter; `harness/README.md` führt
      `lint` nicht mehr unter den Werkzeugen, sondern unter §Sensors, dazu
      `lint-gegenprobe`, je mit Vertrag und Bindung an die ADR, die Zelle von `make lint`
      verlinkt die Sensor-Datei `harness/sensors/lint.md` (Vertrag `SPEC-049`, Grenze,
      Ausgabe und Ausgänge; per `cp` aus der Vorlage
      `.harness/baseline/v6.16.0/templates/harness/sensors/gate.template.md`); die Zeile
      „Nicht behauptet“ nennt Lint nicht mehr. Kein Kommentar nennt `make lint` noch ein
      Werkzeug ohne Gate (`harness/mk/lint.mk` mit Hilfetext, Kommentar der Stufe `lint` im
      `Dockerfile`), der Kopf von `tools/harness/lint.sh` zeigt auf die Sensor-Datei und
      ordnet unter `GEPRÜFT DURCH tools/harness/lint-gegenprobe.sh:` jedem Punkt von
      `SPEC-049` seine Fälle zu, und der Kommentar an `Meldungen` in
      `internal/hexagon/model/fehler.go` nennt die Tiefensuche (§3). — Verifikation,
      Abschnitt 1 Punkt 3; die Zuordnung der Fälle aus F-471 und F-472 in den Köpfen
      (V-98) und der Hilfetext ohne „netzlos“ (V-99) in `5f9fef5`, die Grenze der
      zulässigen Regeln in `harness/sensors/lint.md` (V-100) ebenda.
- [x] `make gates` grün. — an `29bcb16` (Verifikation, Abschnitt 4: Exit 0, 2 min 27 s)
      und an `5f9fef5` (Planner, vor dieser Closure: Exit 0, 2 min 29 s, `lint-gegenprobe: gruen` mit den beiden Fällen aus V-96; Arbeitsbaum mit Code und Skripten gleich `5f9fef5`). Diese Closure ändert nur
      Plan, Roadmap und Register; `make docs-check` und `make kopf-check` grün am Stand
      dieser Closure.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      — `docs/reviews/2026-10-08-review-slice-harness-lint.md` (`c3c8ead`; F-471 bis
      F-474 MEDIUM, F-475 und F-476 LOW, F-477 bis F-479 INFO), Nacharbeit in `53c70fc`
      (Architect, F-475) und `29bcb16`; Verifikation
      `docs/reviews/2026-10-08-verifikation-slice-harness-lint.md` (`6fd2306`; V-96
      MEDIUM, V-97 bis V-100 LOW, V-101 Hinweis), V-97 und V-100 entschieden in
      `d3473b9`, V-96, V-98, V-99 und V-100 umgesetzt in `5f9fef5`, V-101 im Risiko
      *Laufzeit*.
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
| `docs/plan/adr/0034-lint-gate-mit-solid-nahem-profil.md`, `docs/plan/adr/README.md` | neu / update | geschrieben vom Architect vor dem Code: Entscheidung, Gründe, Messung des Bestands, dauerhafte Ausnahmen, Einführung nach Bereinigung, Werkzeug-Ziel davor; Index-Zeile. Den Status `Accepted` nach der Wahl des Nutzers vom 2026-10-06 setzt der Architect vor dem Start von `slice-harness-lint-werkzeug` |
| `spec/spezifikation.md` §11, §12 | update | geschrieben vom Architect vor dem Code: `SPEC-049` mit den Randformen aus §6; Historien-Zeile. Am 2026-10-07 vor dem Code die Grenze um ungültiges YAML und den Cache des Builds fortgeschrieben (§6, *Prüfung vor dem Code*); am 2026-10-07 nach dem ersten Code-Commit Punkt 8 um Leerraum nach `exclusions:` und `rules:` und den Eintrag `-` allein (`7a32d57`, §6, *Randform-Rückgaben aus `221845d`*); am 2026-10-08 die Grenze um die zulässigen Regeln (§6, *V-100*), je mit Historien-Zeile |
| `harness/mk/lint.mk` | update | aus `slice-harness-lint-werkzeug`: `lint` an `GATE_CHECKS`, dazu das Ziel `lint-gegenprobe`, ebenfalls an `GATE_CHECKS`; Kopfkommentar und Hilfetext nennen das Gate statt des Werkzeugs |
| `Dockerfile`, `tools/harness/lint.sh` | update | nur Kommentare: Die Stufe `lint` ist Teil der Gate-Kette; der Kopf von `lint.sh` zeigt auf die Sensor-Datei und führt `GEPRÜFT DURCH tools/harness/lint-gegenprobe.sh:` mit einer Zeile je Punkt von `SPEC-049`, Muster `tools/harness/kopf-check.sh`. Kein Verhalten, kein Profil |
| `harness/sensors/lint.md` | neu | per `cp` aus `.harness/baseline/v6.16.0/templates/harness/sensors/gate.template.md`: Vertrag als Link auf `SPEC-049`, Grenze aus ihr, Ausgabe und Ausgänge, keine Sperre (die Stufe hat keine benannte Abbruch-Meldung), Bindung [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md); entscheidet nichts neu |
| `tools/harness/lint-gegenprobe.sh` | neu | Mutanten in Kopien des Arbeitsbaums unter `mktemp -d` (`cp -R` ohne `-p`, `touch` je mutierter Datei), je Kopie der Docker-Build der Stufe `lint` mit erwartetem Exit und erwarteter Zeile, bis zu sechs gleichzeitig; Grundlauf zuerst; Vorbild `make kopf-check-gegenprobe` und `make a-check-negativ`. Die Fälle mit eigener Zeile teilen sich eine Kopie, allein läuft ein Fall, dessen Zusage der Ausgang ist oder der andere ausschließt (§6 *Ort in der Gate-Kette*); den Fall nach Punkt 10 fährt `make gates` in einer Kopie mit nur dem Fragment `lint.mk`. Je Zusage aus `SPEC-049` ein Mutant, der rot wird: je aktivem Linter einer (statt je Gruppe), je Schwelle ein grüner und ein roter an der Grenze, je Einstellung und je Regel von `revive`, je Schreibweise von `nolint`, je Form des Profils, je Regel einer benannten Globale, je Art einer `lint:`-Zeile allein mit Ausgang 1; Rot-Fälle für die erste Bedingung von Messmethode 3 je Paketgruppe (ein White-Box-Test unter `internal/hexagon/model`, `internal/adapters/driven/recording`, `internal/adapters/driving/pgwire` und `internal/adapters/driving/cli`), für die zweite Bedingung Zugriff an der Brücke vorbei (`internal_test.go` im Paket des Codes), Test in der Brückendatei, Variable in der Brückendatei; dazu ungenutzte Regeln und Regeln ohne `Why:`. Den grünen Fall je Ausnahme trägt der Bestands-Lauf mit `warn-unused` (`SPEC-049` Punkt 8). Was sich nicht unterscheiden lässt, steht als offen im Kopf. Ohne Abdeckungs-Deklaration (§1) |
| `internal/hexagon/model/fehler.go` | update | Kommentar an `Meldungen` (Zeilen 171 bis 173) nennt die Tiefensuche der Spezifikation statt „außen nach innen und in der Reihenfolge seiner Ursachen“ (`AGENTS.md` §3.11); übernommen aus der Closure von welle-replay-semantik, Nebenbefund 2. Nur der Kommentar, kein Verhalten |
| `AGENTS.md` | update | §3.2 mit echtem Träger, Falsch/Richtig mit diesem Repo |
| `harness/README.md` | update | Zeile `make lint` aus den Werkzeugen nach §Sensors, Target-Zelle verlinkt `sensors/lint.md`, Vertrag ein Satz mit `SPEC-049`; dazu `make lint-gegenprobe` mit seinem Vertrag in einem Satz (keine Sensor-Datei, wie bei den übrigen Gegenproben); „Nicht behauptet“ ohne Lint |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-lint-werkzeug`, die vier
Umstellungs-Slices (`slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
`slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`) und beide
Bereinigungs-Slices (`slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`)
liegen in `done/`, und `make lint` meldet am Stand des Starts keinen Befund (WIP-Limit
1). Reihenfolge nach Entscheidung des Nutzers vom 2026-10-05, 2026-10-06 und
2026-10-07: `slice-harness-lint-werkzeug`, die vier Umstellungs-Slices, die beiden
Bereinigungs-Slices, dieser Slice, `slice-harness-commit-struktur-id`,
`slice-harness-integration-wait`, `slice-harness-meldungskatalog-gate`, `slice-harness-abdeckung-gate`,
`slice-harness-coverage`, `slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`, `slice-harness-mutation`. Vor dem ersten Code-Commit prüft
der Architect, ob die Randformen aus §6 noch dem Werkzeug entsprechen
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Start-Bedingungen erfüllt** (Architect, 2026-10-07, Stand `9fe02ca`): Die sieben
Slices liegen in `done/`. `make lint` endet mit Ausgang 0, die Stufe meldet `0 issues.`
und keine `lint:`-Zeile; `.golangci.yml` ist seit `8d9f0ee` unverändert, die Bereinigung
hat also keine Ausnahme ergänzt. Die Randformen sind gegen das Werkzeug geprüft (§6,
*Prüfung vor dem Code*).

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
  `export_test.go` (eigene Prüfung), dazu eine Variable in der Brücke. Dass Brücke und
  Test keinen Wert eines unexportierten Typs anlegen (Punkt 7), prüft das Werkzeug
  nicht (Grenze); dafür gibt es keinen Rot-Fall, es bleibt Urteil des Review.
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

**Prüfung vor dem Code** (Architect, 2026-10-07, Stand `9fe02ca`): Die Randformen oben
entsprechen dem Werkzeug. Für das Gate und seine Gegenprobe gilt:

- **Gate bei rotem Lint** — entschieden in `SPEC-049` Punkt 9 und 10: Die Stufe endet
  mit Ausgang ungleich 0, `docker build` und damit `make lint` ebenso, `make gates` endet
  mit Fehlerstatus, und `record-gates` läuft nicht. Die Meldung ist die Ausgabe der Stufe
  im Build-Log (`--progress=plain`): die `lint:`-Zeilen und die Befunde von
  golangci-lint. Dass `make -j` ohne `-k` nach dem ersten roten Ziel keine weiteren
  startet, ist Verhalten von `make` und gilt für jedes Gate (Grenze).
- **Ort in der Gate-Kette und Laufzeit** — `GATE_CHECKS += lint lint-gegenprobe` in
  `harness/mk/lint.mk`, ohne Ordnungskante zu anderen Zielen; der Nachweis läuft wie
  immer zuletzt. Gemessen an `9fe02ca`: der Schritt der Stufe 8,4 s bei warmem Cache der
  Stufe `deps`. Die Gegenprobe darf Fälle in einer Kopie zusammenfassen, wenn jeder Fall
  an seiner eigenen Zeile erkannt wird; allein läuft ein Fall, dessen Zusage der Ausgang
  selbst ist (eine `lint:`-Zeile allein, Ablehnung erst beim Laden) oder der andere
  ausschließt (Profil fehlt, Form des Profils). Laufzeit bleibt Beobachtung, keine
  Schwelle (Risiko unten).
- **Cache der Docker-Stufe** — entschieden in der Grenze von `SPEC-049` (2026-10-07): Mit
  gleichem Image, gleichen Modulen und gleichem Inhalt des Build-Kontexts nimmt der Build
  das Ergebnis aus dem Cache (an `9fe02ca` gesehen: zweiter Lauf `CACHED`). Das ist das
  Ergebnis derselben Eingaben, denn ein roter Lauf liegt nie im Cache, und Image, Netz
  und Toolchain sind fest. Kein `--no-cache`: Es baute auch `deps` neu und bräuchte Netz.
  Akzeptiertes Negativ: Eine Datei mit gleicher Größe und Änderungszeit überträgt der
  Build nicht neu (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`); im Arbeitsbaum
  ändert ein Editor die Änderungszeit, nur `cp -p` oder `touch -r` halten sie.
- **Gegenprobe als eigenes Gate-Ziel** — Muster `kopf-check-gegenprobe` und
  `a-check-negativ`: `lint-gegenprobe` an `GATE_CHECKS`, Skript
  `tools/harness/lint-gegenprobe.sh`, auf dem Host mit `bash` und `docker`. Sie kopiert
  den Arbeitsbaum, nicht `HEAD`, denn `make gates` steht für den Arbeitsbaum. Die Kopie
  liegt unter `mktemp -d`, ohne `cp -p`, mit eigenem Pfad je Lauf; die Gegenprobe
  schreibt nichts in den Arbeitsbaum und räumt ihre Temp-Pfade ab. Netz braucht wie bei
  `make lint` nur `deps` bei leerem Cache. Keine Sensor-Datei, und keine Kennung in der
  Spezifikation: Die Gegenprobe prüft, ob das Werkzeug richtig ist, und ihr Vertrag ist
  ein Satz in ihrer Zeile in `harness/README.md` (Entscheidung von
  `slice-harness-vertraege-spezifikation`, §6 *Welche Targets eine Sensor-Datei
  bekommen*). Akzeptiertes Negativ: Jeder Lauf, auch jeder Lauf von `make lint`, lässt
  ungetaggte Images im Docker-Cache liegen; Aufräumen ist Sache des Hosts.
- **Übernommene Fälle aus `slice-harness-lint-werkzeug`** — alle entschieden: Einzug 8
  ist ein Eintrag mit anderem Einzug (Punkt 8). Bei Einzug 2 und ungültigem YAML sagt
  die Grenze von `SPEC-049` seit 2026-10-07 nur die erste Eintragszeile zu (vorher stand
  das nur in V-75 der Verifikation). V-74 ist die Ablehnung erst beim Laden (Grenze). Die
  Punkte 6, 8 und 9 haben ihre Fälle in der DoD. Für die Zusagen ohne unterscheidende
  Mutation gilt die DoD: Fall oder benannt als offen im Kopf der Gegenprobe.
- **`contextcheck` in Testdateien** (Hinweis aus `slice-harness-blackbox-kern`) —
  entschieden in Punkt 5: ohne Ausnahme, Kontext aus `t.Context()` bzw.
  `context.WithoutCancel(t.Context())`; der Fall der Linter-Gruppe deckt ihn.
- **Wert eines unexportierten Typs** (Punkt 7) — Grenze von `SPEC-049`, wie oben unter
  *Export-Test-Brücke*: kein Rot-Fall, Urteil des Review.
- **Sensor-Datei und `harness/README.md`** — `harness/sensors/lint.md` nach der Vorlage,
  wie `kopf-check.md`; die Zeile `make lint` wandert aus den Werkzeugen nach §Sensors,
  `make lint-gegenprobe` kommt dazu, „Nicht behauptet“ nennt nur noch Testabdeckung
  (DoD, §3).
- **Kommentare des Werkzeugs** — Die Kommentare in `harness/mk/lint.mk`, im `Dockerfile`
  (Stufe `lint`: „Kein Teil der Gate-Kette“) und im Kopf von `tools/harness/lint.sh`
  ändert dieser Slice, weil sie mit dem Gate falsch würden (`AGENTS.md` §3.11). Das ist
  keine Änderung an Profil oder Stufe im Sinn von §1.
- **[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md): keine Folge-ADR** —
  Gate nach Bereinigung (Entscheidung 5), Gegenprobe (Folgepflicht) und beide Ziele in
  `harness/README.md` (Fitness Function) stehen dort. Die Bereinigung kam ohne neue
  Ausnahme aus, also ist kein Re-Evaluierungs-Trigger eingetreten. Die beiden neuen
  Sätze der Grenze sind Randformen in der Spezifikation und ändern die Entscheidung nicht.

**Randform-Rückgaben aus `221845d`** (Architect, 2026-10-07; Belege in §7, *Grüne
Mutanten*). Vorgabe an den Implementer je Rückgabe:

- **Leerraum nach `exclusions:` bzw. `rules:` ohne Kommentar** — ist Blockform;
  entschieden in `SPEC-049` Punkt 8 („nach dem Doppelpunkt höchstens Leerraum und ein
  Kommentar“). YAML liest beide Formen gleich; ein Befund dort wäre eine Form, die
  golangci-lint annimmt und die Prüfung ablehnt. **Fall: ja**, grün: `exclusions:` und
  `rules:` je mit Leerzeichen und Tab am Zeilenende, sonst Bestand; Erwartung Ausgang 0,
  keine `lint:`-Zeile. Der Mutant ohne Abschneiden des Leerraums wird damit rot.
- **Ein `-` allein auf der Zeile unter `rules`** — ist ein Eintrag; entschieden in
  `SPEC-049` Punkt 8 (mit sechs Leerzeichen, Inhalt auf den Fortsetzungszeilen mit
  mindestens acht, Kommentarblock darüber wie bei jedem Eintrag). Das ist gültiges YAML,
  und golangci-lint wendet die Regel an; ein `Form nicht erkannt` meldete eine Regel, die
  wirkt. **Fälle: ja**, zwei: (a) grün: eine Regel des Bestands als `-` allein, darüber
  ihr `# Why:`-Block, Inhalt mit acht Leerzeichen darunter; Erwartung Ausgang 0, keine
  `lint:`-Zeile. (b) rot: dieselbe Regel ohne Kommentarblock; Erwartung Ausgang 1 und die
  Zeile `lint: .golangci.yml:<zeile des ->: Regel ohne Kommentarblock "# Why:"
  unmittelbar darüber`, keine Zeile `Form nicht erkannt`. Der Mutant `/^- /` wird an
  beiden rot. Ein `-` allein mit anderem Einzug als sechs bleibt `Form nicht erkannt`.
- **Pin und Plattform des Images** — **kein Fall, die Stufe bleibt unverändert**; offen
  im Kopf der Gegenprobe, mit Grund. Eine Ausgabe der Version in der Stufe wäre Verhalten,
  das nur der Gegenprobe dient, und ein Textvergleich mit der `FROM`-Zeile wäre eine
  zweite Quelle für den Pin, den `SPEC-049` Punkt 2 nur im `Dockerfile` führt. Der Pin
  ändert sich nur durch einen Commit an dieser Zeile, und den fängt der
  Re-Evaluierungs-Trigger von [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md)
  (Anhebung von golangci-lint) und das Review. Akzeptiertes Negativ, keine Folgepflicht.
  Dasselbe gilt für die Zeile `Felder nicht erkannt`, die die gepinnte Version nicht
  auslöst (§7): offen im Kopf, kein Fall.

**Randform aus dem Review, F-475** (Architect, 2026-10-08): `LINT_GEGENPROBE_PARALLEL`.
Entschieden hier und nicht in `SPEC-049`: Die Einstellung gehört zur Gegenprobe, und die
Gegenprobe hat keinen Vertrag in der Spezifikation (§6 oben, *Gegenprobe als eigenes
Gate-Ziel*). Vorgabe an den Implementer:

- **Werte.** Nicht gesetzt: 6. Gesetzt: eine positive ganze Zahl in Dezimalform ohne
  Vorzeichen, Leerraum und führende Null (`^[1-9][0-9]*$`); sie ist die Höchstzahl
  gleichzeitiger Läufe der Stufe. Keine Obergrenze.
- **Jeder andere Wert**, auch der leere, `0`, eine negative Zahl, ein Wert mit Leerraum
  oder einer, der keine Zahl ist: Abbruch vor der ersten Kopie und dem ersten Lauf, mit
  der Zeile `lint-gegenprobe: LINT_GEGENPROBE_PARALLEL ist keine positive ganze Zahl:
  '<wert>'` auf stderr und Ausgang 2. Leer zählt nicht als nicht gesetzt: Wer die
  Variable setzt, meint einen Wert, und `make … LINT_GEGENPROBE_PARALLEL=` ist ein
  Schreibfehler, keine Bitte um den Default.
- **Fälle: ja**, je ein Aufruf des Skripts mit `0`, leer, `-1`, `abc` und ` 3`, jeder
  unter `timeout` (wenige Sekunden). Erwartung je Fall: Ausgang 2 (nicht 124), die Zeile
  oben mit dem Wert, und kein Lauf der Stufe hat begonnen (keine Zeile eines Falls und kein
  `docker build` in der Ausgabe). Ohne die Prüfung hängt `0` bis zum `timeout` und wird
  rot; dasselbe gilt für die übrigen Werte.
- **Kein Fall** für einen gültigen Wert außer dem Default: Er führe die ganze Gegenprobe
  ein zweites Mal. Den Default fährt jeder Lauf von `make gates`. Akzeptiertes Negativ:
  Ein Muster, das einzelne gültige Zahlen ablehnt, fände erst ein Aufruf mit diesem Wert.
- **Ort der Zusage.** Kopf von `tools/harness/lint-gegenprobe.sh`. `harness/README.md`
  und der Hilfetext nennen die Einstellung nicht, sie ist ein Hilfsmittel des Laufs und
  kein Teil des Gates.

**Randform aus der Verifikation, V-100** (Architect, 2026-10-08): Ob eine Regel unter
`rules` zu den zulässigen nach `SPEC-049` Punkt 8 gehört, prüft das Werkzeug nicht.
Entschieden als **Grenze**, festgehalten in der Grenze von `SPEC-049`; **keine
Werkzeugänderung, kein Fall.** Eine Prüfung bräuchte eine zweite Liste der zulässigen
Regeln neben Punkt 8 und urteilte doch nicht über den Grund. Das Profil enthält heute
nur zulässige Regeln (Verifikation, 16 Regeln gelesen). Eine neue Ausnahme ist eine
Lockerung und braucht nach `AGENTS.md` §3.6 eine ADR; das Review findet sie am Diff von
`.golangci.yml`. `warn-unused` und die Why-Prüfung bleiben, was das Werkzeug dazu
beiträgt. Vorgabe an den Implementer: (a) im Kopf von `tools/harness/lint-gegenprobe.sh`
unter OFFEN die Zeile „Punkt 8, zulässige Regeln: ob eine Regel zu den zulässigen gehört
und mehr als Bestand ausblendet, prüft das Werkzeug nicht (Grenze von `SPEC-049`); kein
Fall“; (b) unter *Grenze* in `harness/sensors/lint.md` derselbe Satz wie in der
Spezifikation, weil die Sensor-Datei deren Grenze wiedergibt. Keine Änderung an
`lint.sh`, Profil oder Stufe.

**Risiken:**

- **Bestand** (Hauptrisiko) — gemessen ([ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md)): 142 Befunde ohne `testpackage`,
  166 mit; nach den dauerhaften Ausnahmen 96 mit `uniq-by-line: true`, nach `SPEC-049`
  Punkt 9 103, davon 32 im Produkt-Code in allen Schichten (acht Funktionen über den
  Komplexitäts-Schwellen) und 71 in Testdateien (`make lint`, Stand des
  Werkzeug-Slice).
  Eingetreten: Rückführung nach `next/` (§4); den Bestand bereinigen die vier
  Umstellungs-Slices, `slice-lint-bestand-kern-driven` und `slice-lint-bestand-driving`
  vor dem Start dieses Slice. — **Ausgang:** **eingetreten** → `slice-harness-lint-werkzeug`,
  `slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`,
  `slice-harness-blackbox-pgwire`, `slice-harness-blackbox-einstieg`,
  `slice-lint-bestand-kern-driven`, `slice-lint-bestand-driving`. Zugewiesen bei der
  Rückführung am 2026-10-06, als alle in `open/` lagen; alle liegen in `done/` und haben
  geliefert: `make lint` modulweit `0 issues.` (Architect an `9fe02ca`, Verifikation
  Abschnitt 1 Punkt 1 an `29bcb16`).
- **Bestand wächst bis zum Gate nach** — zwischen der Bereinigung und dem Gate prüft
  niemand neuen Code automatisch ([ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md), Konsequenzen, negativ akzeptiert); ein
  Befund aus dieser Zeit stünde beim Start ohne Stufe da. Die Start-Bedingung in §4
  verlangt `make lint` ohne Befund, die Rückführung `→ open` gibt einen Befund an den
  Slice, der ihn einbrachte. — **Ausgang:** **entfallen.** Begründung: `make lint` war
  beim Start ohne Befund (§4, `9fe02ca`) und an `29bcb16` in einer frischen Kopie ebenso
  (Verifikation); die Zeit ohne Prüfung ist vorbei, denn ab diesem Slice hängt `lint` an
  `GATE_CHECKS`, und neuer Code mit Befund macht `make gates` rot.
- **Gegenprobe sieht den Mutanten nicht** — BuildKit überträgt eine Datei gleicher
  Größe und mtime nicht neu (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`, 1×;
  `.claude/commands/implement-slice.md` Schritt 19). Die Gegenprobe baut aus einer
  Kopie unter eigenem Temp-Pfad, nicht aus dem Arbeitsbaum. — **Ausgang:** **entfallen.**
  Begründung: Die Gegenprobe kopiert mit `cp -R` ohne `-p` und ruft `touch` auf jede
  mutierte Datei (Review F-479, Verifikation Abschnitt 5 zu §3); von 16 Mutanten der Verifikation
  wurden 14 rot, die zwei grünen (V7, V11) waren Lücken der Fälle, kein Cache-Effekt,
  und sind seit `5f9fef5` rot. Kein Beleg für
  `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`.
- **Konfiguration wirkt anders, als sie gelesen wird** — ein Pfad-Regex in
  `exclusions`, ein Modulname in `gomodguard_v2` oder ein fehlendes `build-tags`
  liest sich als Zusage und greift nicht
  (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, 2×; ein weiterer Beleg aus
  diesem Slice erreicht die Schwelle 3× und braucht einen Folge-Slice). Je Ausnahme und
  je Liste ein Fall der Gegenprobe. — **Ausgang:** **entfallen.** Begründung: Je Liste
  und je Ausnahme hat die Gegenprobe einen Fall, und die Mutanten daran wurden rot (V5
  `errcheck`, V6 `gomodguard_v2`, V9 `forbidigo`, V10 `testpackage`, V12 `ireturn`;
  `p1-integration-build-tag` und `p1-pfad-anker-*` aus §7). Was Review und Verifikation
  fanden, waren Prüf-Lücken bei Konfiguration, die so wirkt, wie sie gelesen wird: F-472
  (der Anker von `skip-regexp` wirkt, nur fehlte der Fall) und V-100 (die Zulässigkeit
  einer Regel prüft das Werkzeug nicht, jetzt Grenze von `SPEC-049`). Beide zählen unter
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` bzw. sind entschieden; kein dritter
  Beleg für `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, der Eintrag bleibt
  2×.
- **Bestehende Prüfung fällt weg** — `gofmt` und `go vet` in der Stufe `test` bleiben;
  die Stufe `lint` tritt daneben, nicht an ihre Stelle
  (`BEO-REPO/gate-regel-ersetzt-statt-ergaenzt`, 1×). — **Ausgang:** **entfallen.**
  Begründung: Die Stufe `test` ist unverändert, der Diff am `Dockerfile` ändert nur den
  Kommentar der Stufe `lint` (Verifikation, Abschnitt 1 Punkt 3); `gofmt -l` und
  `go vet` laufen in `make test` wie zuvor.
- **Laufzeit von `make gates`** — eine weitere Docker-Stufe mit eigener Analyse (8,4 s
  an `9fe02ca`) und die Gegenprobe mit einem Lauf der Stufe je Kopie; `make -j` fährt sie
  parallel. Wird sie zum Engpass, ist das eine Beobachtung, kein
  Grund für eine Lockerung. — **Ausgang:** **weiter offen** →
  `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (neu, 1×). Gemessen (V-101, je Lauf
  allein, warmer Cache, 20 Kerne): `make lint-gegenprobe` 1 min 24 s, `make gates`
  2 min 27 s; die Gegenprobe ist mit rund 60 % der größte Posten. Tragbar heute, aber
  `slice-harness-coverage`, `slice-harness-abdeckung-gate` und `slice-harness-mutation`
  hängen weitere Läufe an dieselbe Kette.

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

- **Was hat funktioniert:** Das Gate hängt an der Kette und hält: `lint` und
  `lint-gegenprobe` an `GATE_CHECKS`, ein `//nolint` macht `make gates` rot, `record-gates`
  läuft dann nicht (Verifikation, Punkt 1). Die Reihenfolge aus §4 trug: Werkzeug ohne
  Gate, vier Umstellungs-Slices, zwei Bereinigungs-Slices, dann dieser; beim Start war
  `make lint` modulweit ohne Befund, und `.golangci.yml` blieb im ganzen Slice
  unverändert. Die Prüfung vor dem Code (§6, Architect an `9fe02ca`) entschied Gate-Ort,
  Cache und Form der Gegenprobe vorab; die beiden Randformen, die der erste Code-Commit
  fand (Leerraum nach dem Doppelpunkt, `-` allein), gingen ohne neue Entscheidung im Code an
  den Architect zurück (`221845d` → `7a32d57`), wie `AGENTS.md` §3.12 es verlangt. Der
  Schnitt hielt: drei Liefer-Punkte, Harness und Doku, ein Kommentar im Produkt-Code; die
  Gegenprobe ist lang (1177 Zeilen), aber als Tabelle aus Fall und Erwartung in einer
  Sitzung prüfbar (Review F-477, Verifikation §5 zu §4). Die Belege in §7 erreichten
  Review und Verifikation; die nachgemessenen Zahlen stimmten.
- **Was ging anders als geplant:**
  1. Die Gegenprobe deckte nicht jede Zusage von `SPEC-049`, obwohl die DoD es verlangte
     und §7 33 Mutationen rot zeigte. Das Review fand drei grüne Mutanten (F-471 `nolint`
     in Testdateien, F-472 Anker von `skip-regexp`, F-473 untere Grenze von `dupl`), die
     Verifikation nach deren Nacharbeit zwei weitere (V-96: `/` vor `nolint`, Doppelpunkt
     in `# Why:`). Zwei Runden für dieselbe Klasse, wie schon in
     `slice-harness-lint-werkzeug`. Gesucht hat der Implementer je Zusage eine rote
     Mutation, nicht je Zeichen eines Musters die Mutation, die es entfernt.
  2. Kommentare und Hilfetext sagten „ohne Netz“ und „nichts in den Arbeitsbaum“ als
     geprüft zu, obwohl beides OFFEN stand (F-474, MEDIUM als dritte Wiederholung; Rest im
     Hilfetext V-99).
  3. Die Einstellung `LINT_GEGENPROBE_PARALLEL` kam mit `221845d`, ohne dass §6 sie
     nannte; ihre Randformen (`0` hängt, ein Wort hebt die Grenze auf) entschied der Code
     still, erst das Review fand sie (F-475). Der Architect entschied sie in `53c70fc` in
     §6, nicht in `SPEC-049`, weil die Gegenprobe keinen Vertrag in der Spezifikation hat.
  4. Der Plan folgte zwei Korrekturen nicht sofort: Kopf, §1 und §3 nannten die Änderung
     an Punkt 8 aus `7a32d57` nicht (F-476), und die DoD verlangte weiter einen Fall für
     `Felder nicht erkannt`, den §6 als offen entschieden hatte (V-97).
  5. `SPEC-049` Punkt 8 sagte zu, dass nur die genannten Regeln zulässig sind; das
     Werkzeug prüft das nicht. Der Architect machte es zur Grenze (`d3473b9`, V-100).
  - **Summary-Zeilen:** Review `docs/reviews/2026-10-08-review-slice-harness-lint.md`:
    „0 HIGH · 4 MEDIUM · 2 LOW · 3 INFO (F-471 `nolint` in Testdateien ohne Fall, M1
    grün; F-472 Anker von `skip-regexp` ohne Fall, M2 grün; F-473 untere Grenze von
    `dupl` weder Fall noch offen, M3 grün; F-474 Kopf von `lint.sh` und `lint.mk` sagen
    „ohne Netz“ bzw. „nichts in den Arbeitsbaum“ als geprüft zu; F-475 Randformen von
    `LINT_GEGENPROBE_PARALLEL` vom Code entschieden; F-476 §1, §3 und Kopf folgen der
    Änderung an Punkt 8 nicht; F-477 Größe trägt, keine Rückführung nötig; F-478
    Folgerot der Gegenprobe ohne Lesehinweis; F-479 Verdrahtung und Hygiene der
    Gegenprobe bestätigt). Wiederkehrende Klassen:
    `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (F-471 bis F-473),
    `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-474),
    `BEO-REPO/plan-folgt-korrektur-nicht` (F-476).“ Verifikation
    `docs/reviews/2026-10-08-verifikation-slice-harness-lint.md`, Urteil: „Fast fertig,
    ein Liefer-Punkt ist noch nicht ganz bestätigt.“ Gate (Punkt 1) bestätigt und selbst
    rot gezeigt, `make lint` modulweit 0; Gegenprobe (Punkt 2) mit Grundlauf und
    Verdrahtung, F-471 bis F-473 und F-475 behoben, von 13 eigenen Mutationen 11 rot,
    zwei grüne (V-96, MEDIUM); Doku (Punkt 3) bestätigt mit V-98 und V-99; `make gates`
    grün. V-97 und V-100 an den Architect, V-101 Laufzeit an den Planner. Umgesetzt:
    V-97 und V-100 in `d3473b9`, V-96, V-98, V-99 und V-100 in `5f9fef5`; V-96 ist
    vor dem Haken an DoD-Punkt 2 nachgeprüft (§2).
- **Steering-Loop-Eintrag:** Neuer Sensor: das Lint-Gate `make lint` mit dem Profil
  `.golangci.yml` und den eigenen Prüfungen nach `SPEC-049`, dazu seine Gegenprobe
  `make lint-gegenprobe`, beide an `GATE_CHECKS`. Damit hat `AGENTS.md` §3.2
  (Suppression-Verbot) einen Träger statt eines Platzhalters, und Teil 1 und 3 von
  `LH-QA-07` sind durch ein Gate belegt (die Abdeckungs-Deklaration folgt mit
  `slice-harness-abdeckung-gate`). liegt in `harness/README.md §Sensors` — Herkunfts-Anker
  `seit slice-harness-lint`.

  Retirement-Checks: `AGENTS.md` §3.9 (seit welle-walking-skeleton) ist wieder
  aufgetreten (F-476, V-97); die Regel bleibt, `make kopf-check` fängt die Kopf-Hälfte,
  §1, §3 und DoD bleiben Urteil. §3.10 (seit welle-extended-query) ist wieder aufgetreten,
  in zwei Runden (F-471 bis F-473, V-96); die Regel bleibt, ihr geplanter Sensor
  `slice-harness-mutation` trägt weiter. §3.11 (seit welle-extended-query) ist wieder
  aufgetreten (F-474, V-99); die Regel bleibt. §3.12 (seit slice-harness-randformen-vor-code)
  ist wieder aufgetreten (F-475: eine Einstellung der Gegenprobe ohne Randformen in §6);
  die Regel bleibt, die beiden Rückgaben aus `221845d` zeigen, dass die Randform-Rückgabe
  trägt, wo der Implementer die Randform sieht. §3.13 (seit slice-lint-bestand-kern-driven):
  Dieser Slice wies keine neue Adresse zu; nicht berührt. `implement-slice` Schritt 19
  (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`): nicht wieder aufgetreten, jede
  Mutation lief in einer frischen Kopie mit `touch`.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `5f9fef5` (Zähler =
  Dateien unter `evidence/`).
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: **Beleg**, 13× → 14× (F-471 bis
    F-473, V-96), Stand verkörpert, bleibt; Sensor geplant mit `slice-harness-mutation`.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: **Beleg**, 18× → 19× (F-474,
    V-99), Stand verkörpert, bleibt.
  - `BEO-REPO/plan-folgt-korrektur-nicht`: **Beleg**, 16× → 17× (F-476, V-97), Stand
    verkörpert, bleibt.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden`: **Beleg**, 11× → 12× (F-475),
    Stand verkörpert, bleibt. F-475 war eine Randform, die der Code entschied (`0` hängt,
    ein Wort hebt die Grenze auf), aber der Implementer gab sie nicht zurück; das Review
    fand sie. Deshalb kein Beleg für
    `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`: Dessen Merkmal ist die
    Rückgabe nach der Entscheidung im Code, und eine Rückgabe gab es nicht.
  - `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`: **neu**, 1× (Risiko *Laufzeit*, V-101),
    Stand offen.
  - Ohne Beleg: `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (bleibt 2×; F-472
    und V-100 sind Prüf-Lücken, keine Konfiguration, die anders wirkt, siehe Risiko in
    §6), `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (bleibt 3×,
    verkörpert; oben), `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (bleibt 2×; V-96
    betrifft einen neuen Vertrag und zählt nach der Abgrenzung des Eintrags unter
    `negativtests-fehlen-bei-neuem-vertrag`), `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`
    (bleibt 1×), `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (bleibt 1×, Stufe `test`
    unverändert), `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (bleibt 1×, der
    Architect prüfte §6 vor dem Code an `9fe02ca`),
    `BEO-REPO/kern-fremdimporte-nur-ueber-tech-liste` (bleibt 1×, §1),
    `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (bleibt 3×, keine neue Adresse).

  Einmalig und nicht eingetragen: F-477 und F-479 (Negativbefunde), F-478 (Lesehinweis
  zum Folgerot, umgesetzt in `29bcb16`), V-98 (Zuordnung der Fälle in den Köpfen
  unvollständig, umgesetzt in `5f9fef5`). **Kein Eintrag erreicht mit diesem Slice neu
  3×.** Über der Schwelle stehen nur Einträge mit Ausgang (`commit-nennt-struktur-kennung`
  3× geplant; `folge-slice-adresse-nimmt-nicht-an` 3×,
  `implementer-bericht-erreicht-pruefer-nicht` 3×,
  `randform-im-code-entschieden-dann-zurueckgegeben` 3×,
  `spec-randform-erst-im-review-entschieden` 12×, `negativtests-fehlen-bei-neuem-vertrag`
  14×, `plan-folgt-korrektur-nicht` 17×, `zusage-im-kommentar-weiter-als-pruefung` 19×,
  verkörpert; `white-box-liste-vor-code-nur-namenssuche` 3× gestrichen).
- **Folge-Slices:** `slice-harness-abdeckung-gate` (die Abdeckungs-Deklaration von
  `LH-QA-07` Teil 1 und 3 an der Gegenprobe, §1 hier *Ausdrücklich NICHT*); er nimmt an:
  Er liegt in `open/`, sein §1 *Ziel* nennt die Deklaration an der Gegenprobe des
  Lint-Gates, seine DoD verlangt sie im Kopf der Gegenprobe, und sein §1 *Ausdrücklich
  NICHT* nennt diesen Slice nur als Lieferanten der Gegenprobe und ihrer Fälle.
  `slice-harness-coverage` (Schwelle für Testabdeckung, §1 hier *Ausdrücklich NICHT*); er
  nimmt an: Er liegt in `open/`, die Schwelle ist sein Gegenstand, und sein §1 schließt
  nur das Lint-Gate aus, mit dieser Kennung. Die sieben Nehmer des Risikos *Bestand*
  liegen in `done/` und haben geliefert (§6). Als nächster in der Reihe aus §4 folgt
  `slice-harness-commit-struktur-id`.
- **Risiken aus §6:** sechs, je mit Ausgang: *Bestand* **eingetreten** → die sieben
  Slices aus §4 *Start* (alle in `done/`), *Bestand wächst bis zum Gate nach*
  **entfallen**, *Gegenprobe sieht den Mutanten nicht* **entfallen**, *Konfiguration
  wirkt anders, als sie gelesen wird* **entfallen**, *Bestehende Prüfung fällt weg*
  **entfallen**, *Laufzeit von `make gates`* **weiter offen** →
  `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`; die Gründe stehen in §6. Die
  Randformen in §6 sind Entscheidungen, keine Risiken.
- **Drei Paarungen:** Anker: Der Steering-Loop-Eintrag trägt `liegt in
  harness/README.md §Sensors`; die Zeile `make lint` dort trägt in der Bindung
  `· seit slice-harness-lint`. Folge-Slice: `slice-harness-abdeckung-gate` und
  `slice-harness-coverage` liegen in `open/`; `grep -n "slice-harness-lint"` findet die
  Kennung in §1 beider (*Ausdrücklich NICHT*); ob sie die Sendung inhaltlich führen, ist
  oben beurteilt. Für die sieben Nehmer des Risikos *Bestand* (alle in `done/`) findet
  der `grep` die Kennung in §1 von `slice-harness-lint-werkzeug`, der vier
  Umstellungs-Slices und von `slice-lint-bestand-driving`; in
  `slice-lint-bestand-kern-driven` steht sie nur in §4 *Start* (Reihenfolge), sein §1
  *Ziel* führt die Bereinigung von Kern und Driven inhaltlich. Er entstand vor
  `AGENTS.md` §3.13, und geliefert hat er. Register:
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`,
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/plan-folgt-korrektur-nicht`
  und `BEO-REPO/spec-randform-erst-im-review-entschieden` tragen
  `evidence/slice-harness-lint.md`; `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` trägt
  `observation.md`, `state.md` und `evidence/slice-harness-lint.md`. Die übrigen genannten
  Einträge bestehen als Verzeichnis, jedes mit nicht leerem `evidence/`. Die nächste
  Welle-Closure prüft erneut.

**Belege des Implementers** (Stand: Arbeitsbaum auf `fdbcfaf` mit den Änderungen des
Commits, der diesen Abschnitt schreibt; jede Zeile nennt, woran sie gemessen ist).

*Läufe.* `make lint` an diesem Stand in einer frischen Kopie (`cp -R`, dazu eine Datei
`internal/frisch.txt`, damit der Build nicht aus dem Cache nimmt): Ausgang 0, keine
`lint:`-Zeile, `0 issues.`, Schritt der Stufe 8,6 s. `make lint-gegenprobe` allein:
grün, 1 min 37 s (22 Läufe der Stufe, sechs gleichzeitig, auf 20 Kernen). `make gates`
(ohne `-j`) mit `lint` und `lint-gegenprobe`: grün, 2 min 49 s; `lint` kam dort aus dem
Cache des Grundlaufs der Gegenprobe (Grenze von `SPEC-049`: ein grüner Lauf derselben
Eingaben). `git status --porcelain` vor und nach `make lint-gegenprobe` gleich: Die
Gegenprobe schreibt nichts in den Arbeitsbaum.

*Fälle je Punkt von `SPEC-049`* (Namen im Kopf von `tools/harness/lint-gegenprobe.sh`;
ein Fall, der eine Schleife fährt, zählt je Durchlauf): Grundlauf 2 · Punkt 1: 5 ·
Punkt 3: 19 (17 Linter ohne Schwelle, die Liste, `default: none`; die übrigen elf über
4, 5 und 7) · Punkt 4: 18 (neun Schwellen, je grün an und rot über der Grenze) · Punkt
5: 55 · Punkt 6: 12 · Punkt 7: 13 · Punkt 8: 43 · Punkt 9: 72 (davon 51 für
`max-issues-per-linter`) · Punkt 10: 3 · Grenze: 2. Die Fälle aus der DoD stehen darin:
Einzug 8 als vollständige Liste (`p8-form-eintrag-8`), Einzug 2 nur mit der ersten
Eintragszeile und Ausgang 1 (`p8-form-eintrag-2`, V-75), V-74 als `p9-laden`, Punkt 10
als `p10-*`, die White-Box-Fälle der vier Umstellungs-Slices als `p7-whitebox-kern`,
`-driven`, `-pgwire`, `-einstieg`. Je Linter statt je Gruppe, weil die DoD keine
Gruppen nennt und ein Fall je Linter jede Gruppe deckt.

*Umfang und Rückführung aus §4.* Das Skript hat 1074 Zeilen, davon rund 150 Zeilen
Hilfsfunktionen; der Rest sind Go- und YAML-Schnipsel der Fälle und je Fall eine Zeile
Erwartung. Die Fälle zu Messmethode 3 (`p7-*`) sind 13 Erwartungen in rund 60 Zeilen.
Die Rückführung `in-progress → next` ist nicht gezogen: Sie nähme den kleinsten Teil
heraus, und der Rest bleibt eine Tabelle aus Fall und Erwartung. Ob das in eine
Review-Sitzung passt, entscheidet das Review; trägt es nicht, ist der Schnitt nach
Punkten von `SPEC-049` naheliegender als nach Messmethode.

*Weg der Mutanten in den Build* (`implement-slice` Schritt 19): Die Gegenprobe kopiert
den Arbeitsbaum je Fall mit `cp -R` ohne `-p` in ein eigenes Verzeichnis unter
`mktemp -d` und ruft nach jeder Änderung `touch` auf die Datei. Die Mutationen am
Werkzeug unten liefen je in einer eigenen Kopie des Repos (`cp -R`, Ersetzung mit
genau einem Treffer, dann `touch`), in der `make lint-gegenprobe` erneut kopiert.

*Mutationen am Werkzeug: ohne die Prüfung wird der Fall rot* (§3.10; 33 Mutationen,
jede rot gesehen; Spalte 3 nennt die roten Fälle des Laufs):

| Zusage | Mutation | roter Fall |
|---|---|---|
| Punkt 7: keine Testfunktion in der Brücke | `lint.sh`: `(Test\|Benchmark\|Example\|Fuzz)` → `(KeinTest)` | `p7-test-in-bruecke-*` (4), `p9-profil-fehlt-bruecke` |
| Punkt 7: auch mit Tab nach `func` | `lint.sh`: `^func[ ${tab}]+` → `^func[ ]+` | `p7-test-in-bruecke-tab` |
| Punkt 1: Pfade mit `^` | `.golangci.yml`: `path: ^(cmd\|test)/` → `(cmd\|test)/` | `p1-pfad-anker-forbidigo` |
| Punkt 1: Pfade mit `^` | `.golangci.yml`: `text: ^zielart is …` → `text: zielart is …` | `p8-ausnahme-name-zielart` |
| Punkt 1: Build-Tag `integration` | `.golangci.yml`: `build-tags` entfernt | `p1-integration-build-tag` |
| Punkt 3: genau diese Linter | `.golangci.yml`: `- godot` ergänzt | `p3-genau-diese` |
| Punkt 3: `gochecknoinits` aktiv | `.golangci.yml`: `- gochecknoinits` entfernt | `p3-gochecknoinits`, `p9-golangci-allein` |
| Punkt 4: `dupl` 150 | `threshold: 155` | `p4-dupl-rot` |
| Punkt 4: `cyclop` 15 | `max-complexity: 16` | `p4-cyclop-rot` |
| Punkt 5: Regeln von `revive` | `- name: time-naming` entfernt | `p5-revive-genau-diese`, `p5-revive-time-naming` |
| Punkt 5: `testpackage` überspringt nur `export_test.go` | `skip-regexp: (^\|/)(export\|internal)_test\.go$` | `p7-internal-test` |
| Punkt 6: `/*` vor `nolint` | `nolint_muster="(//)…"` | `p6-block` |
| Punkt 8: feste Form, Einzug 6 | `else if (!strich \|\| e != 6)` → `else if (!strich)` | `p8-form-eintrag-8`, `p8-form-eintrag-2` |
| Punkt 8: Fortsetzung mit mindestens 8 | `e >= 8` → `e >= 7` | `p8-form-fortsetzung-7` |
| Punkt 8: jede andere Zeile unter `rules` | `(!strich \|\| e != 6)` → `(e != 6)` | `p8-form-ohne-strich` |
| Punkt 8: höchstens ein Kommentar nach dem Doppelpunkt | `sub(/[ \t]+#.*$/, "", t)` → `t = t` | `p8-kommentar-nach-exclusions`, `-rules`, `p8-ohne-why`, `p8-why-*` |
| Punkt 8: `rules` endet bei Einzug höchstens 4 | `if (!strich && e <= 4)` → `if (0)` | `p8-schluessel-nach-rules` |
| Punkt 8: `exclusions` endet bei Einzug höchstens 2 | `if (in_ex && e <= ex_e …` → `if (0 && …` | `p8-reihenfolge` |
| Punkt 8: Schlüssel `exclusions` in Anführungszeichen | `["\047]?exclusions["\047]?` → `exclusions` | `p8-form-quote-exclusions` |
| Punkt 8: Schlüssel `rules` in Anführungszeichen | `["\047]?rules["\047]?` → `rules` | `p8-form-quote-rules` |
| Punkt 8: `exclusions` in Flussform an anderer Stelle | `(^\|[{,]\|- )` → `(^\|- )` | `p8-form-fluss-anderswo` |
| Punkt 8: `exclusions` als Listenpunkt | `(^\|[{,]\|- )` → `(^\|[{,])` | `p8-form-listenpunkt` |
| Punkt 8: Ausnahmen der Testdateien für jede Testdatei | `path: ^.+_test\.go$` → `^internal/.+_test\.go$` (Komplexität) | `p0-grundlauf`, `p8-ausnahme-test-*` (6) und sechs Fälle, deren Kopie die Testdateien des Bestands trägt |
| Punkt 9: ungekürzt je Text | `max-same-issues: 3` | `p9-max-same-issues-4` |
| Punkt 9: ungekürzt je Linter | `max-issues-per-linter: 50` | `p9-max-issues-per-linter-48` bis `-51` und 15 weitere |
| Punkt 9: Ausgang 1 bei `lint: … fehlt` allein | `befund=1` nach `fehlt` entfernt | `p9-fehlt-allein` |
| Punkt 9: Ausgang 1 bei `von golangci-lint abgelehnt` allein | `befund=1` dort entfernt | `p9-schema-allein` |
| Punkt 9: Ausgang 1 bei `Regel ohne Befund` allein | `befund=1` dort entfernt | `p9-ungenutzt-allein` |
| Punkt 9: Ausgang 1 bei einer Testfunktion in der Brücke allein | `befund=1` dort entfernt | `p9-bruecke-allein` |
| Punkt 9: Ausgang 1 bei `Form nicht erkannt` oder fehlendem `Why:` allein | `befund=1` dort entfernt | `p9-form-allein` und fünf `p8-form-*` |
| Punkt 9: Ausgang 1 bei einem Befund von golangci-lint allein | `if [ "$status" -ne 0 ]` → `if false` | `p9-golangci-allein`, `p9-laden` |
| Punkt 9: Werte der ungenutzten Regel ohne Quoting | `gsub(/\\\\/, "\001", t);` entfernt | `p8-ungenutzt-feldfolge` |
| Punkt 10: `lint` an `GATE_CHECKS` | `lint.mk`: `GATE_CHECKS += lint-gegenprobe` | `p10-gate-checks`, `p10-gates-rot`, `p9-lint-zeile-allein` |

*Grüne Mutanten, eingeordnet* (je ein Lauf von `make lint-gegenprobe` in einer Kopie
mit der Mutation, jeder grün):

| Mutation | Einordnung |
|---|---|
| `Dockerfile`, Stufe `lint`: ohne `GOFLAGS=-mod=readonly` | äquivalent: Ohne `vendor/` (der Build-Kontext lässt es nicht zu) wählt go denselben Modus. Offen im Kopf der Gegenprobe. |
| `Dockerfile`, Stufe `lint`: ohne `GOTOOLCHAIN=local` | äquivalent für den Ausgang: Eine `go`-Zeile über der Version des Images lässt schon `deps` scheitern; offen im Kopf. |
| `Dockerfile`: `RUN` ohne `--network=none` | äquivalent: Keine Prüfung lädt etwas, die Module kommen aus `deps`; offen im Kopf. |
| `.golangci.yml`: `relative-path-mode: wd` | äquivalent: Profil, Modulwurzel und Arbeitsverzeichnis sind in der Stufe `/src`; offen im Kopf. |
| `lint.sh`: `golangci-lint run` ohne `-c` | äquivalent: Der Build-Kontext führt nur `.golangci.yml`, die Default-Suche findet dieselbe Datei; offen im Kopf. |
| `lint.sh`: Leerraum am Ende von `exclusions:`/`rules:` nicht abgeschnitten | an `221845d` grün und zurückgegeben; seit `7a32d57` entschieden (Blockform) und mit `p8-leerraum` gefangen, siehe *Nachtrag zu den Rückgaben*. |
| `lint.sh`: ein `-` allein ist kein Eintrag (`/^- /` statt `/^-( \|$)/`) | an `221845d` grün und zurückgegeben; seit `7a32d57` entschieden (ein Eintrag) und mit `p8-strich-allein` und `p8-strich-allein-ohne-why` gefangen, siehe *Nachtrag zu den Rückgaben*. |

Nicht als Mutation gefahren und offen im Kopf, nach der Entscheidung in §6
(*Randform-Rückgaben aus `221845d`*) mit Grund und ohne Fall: Pin und Plattform des
Images und die Zeile `Regel ohne Befund: Felder nicht erkannt`. Der Kommentar an `Meldungen` in `internal/hexagon/model/fehler.go` sagt die
Tiefensuche zu, die der Fall `Tiefensuche` in `fehler_test.go` prüft (§3.11); er ändert
kein Verhalten.

*Nachtrag zu den Rückgaben* (Stand: `7a32d57` mit den Änderungen des Commits, der
diesen Absatz schreibt). Das Werkzeug entsprach beiden Entscheidungen schon:
`tools/harness/lint.sh` ist unverändert. Neue Fälle, je in einer eigenen Kopie:
`p8-leerraum` (Leerzeichen und Tab nach `exclusions:`, Tab und Leerzeichen nach
`rules:`; Ausgang 0, keine `lint:`-Zeile), `p8-strich-allein` (die erste Regel als `-`
allein, `# Why:`-Block darüber, Inhalt mit acht Leerzeichen; Ausgang 0, keine
`lint:`-Zeile), `p8-strich-allein-ohne-why` (dieselbe Regel ohne Kommentarblock;
Ausgang 1, genau die Zeile `lint: .golangci.yml:158: Regel ohne Kommentarblock "# Why:"
unmittelbar darüber`, kein `Form nicht erkannt`), `p8-form-strich-allein-8` (`-`
allein mit Einzug 8 direkt unter `rules:`; Ausgang 1, `Form nicht erkannt` an seiner
Zeile). Punkt 8 hat damit 47 Fälle, zusammen 248; die Gegenprobe fährt 26 Läufe der
Stufe. Kopf der Gegenprobe: Pin und Plattform sowie `Felder nicht erkannt` offen mit
dem Grund aus §6.

| Zusage | Mutation (je in einer frischen Kopie des Repos) | roter Fall |
|---|---|---|
| Punkt 8: Leerraum nach dem Doppelpunkt ist Blockform | `lint.sh`: `sub(/[ \t]+$/, "", rest)` entfernt | `p8-leerraum` |
| Punkt 8: `-` allein ist ein Eintrag | `lint.sh`: `/^-( \|$)/` → `/^- /` | `p8-strich-allein`, `p8-strich-allein-ohne-why` |

*Nachtrag zum Review* (`docs/reviews/2026-10-08-review-slice-harness-lint.md`, F-471
bis F-475 und F-478; F-475 nach §6 *Randform aus dem Review*; Stand: `53c70fc` mit den
Änderungen des Commits, der diesen Absatz schreibt). Neue Fälle:

- F-471: `p6-testdatei`, ein `//nolint` in `internal/gegenprobe/nolint/fall_test.go`.
- F-472: `p7-endet-auf-export-test`, eine Datei `fooexport_test.go` im Paket des Codes
  ist ein Befund von `testpackage`.
- F-473: `dupl` an der Grenze. Gemessen mit vier Schwellen (148, 149, 150, 151) an
  Paaren aus 45 Anweisungen und einem Schluss: mit `{}` ist ein Paar 149 groß, mit
  `a++` 150, ohne Schluss 148. `p4-dupl-gruen` ist jetzt das Paar mit 149 (kein Befund
  bei 150), `p4-dupl-rot` das Paar mit 150. Die alten Paare (45 und 46 Anweisungen)
  lagen bei 148 und darüber, daher blieb M3 des Review grün.
- F-475: `einstellung-0`, `einstellung-` (leer), `einstellung--1`, `einstellung-abc`,
  `einstellung- 3`. Je ein Aufruf einer Kopie des Skripts in einem Baum ohne Profil,
  unter `timeout 20`; erwartet Ausgang 2 und als einzige Ausgabe die Zeile
  `lint-gegenprobe: LINT_GEGENPROBE_PARALLEL ist keine positive ganze Zahl: '<wert>'`.
  Der Baum ohne Profil hält einen Lauf ohne die Prüfung davon ab, etwas zu bauen: Er
  endet dann an der Profilprüfung mit Ausgang 1 statt zu hängen. Damit ist die Zeile
  ohne Lauf der Stufe gezeigt, `124` kommt in keinem der Mutanten vor. Die Zusage steht
  nur im Kopf des Skripts.

F-474: Der Kopf von `tools/harness/lint.sh` ordnet Punkt 2 keinen Fall mehr zu und sagt
nicht mehr „schreibt nichts in den Arbeitsbaum“. Der Kopf der Gegenprobe sagt weder
„ohne Netz“ noch „Der Arbeitsbaum bleibt unberührt“. Unter OFFEN stehen jetzt
`--network=none` mit den Modulen aus `deps` sowie „nichts in den Arbeitsbaum“ für
`make lint` und für die Gegenprobe; dazu, dass sie ihre Kopien löscht und bei einem
Abbruch Builds beendet. `harness/mk/lint.mk` verweist für das Ungeprüfte auf OFFEN.

F-478: In `harness/sensors/lint.md` steht, dass ein roter Bestand auch die Gegenprobe
rot macht und `p0-grundlauf` die Ursache nennt. Die Gegenprobe beendet bei einem
Abbruch über `set -e` ihre Hintergrund-Läufe: Jeder Lauf wartet auf seinen Build und
beendet ihn bei `TERM`. Sie schreibt dann `lint-gegenprobe: ROT — Abbruch mit Ausgang
<n>`. Das ist einmal von Hand gesehen: eine Kopie des Skripts mit `false` nach `starte
sammel` endete mit dieser Zeile und Ausgang 1, danach lief kein `docker build` mehr. Es
gibt dafür keinen Fall, OFFEN im Kopf.

Zählung: Punkt 6 hat 13 Fälle, Punkt 7 hat 14, dazu kommen 5 Fälle zu
`LINT_GEGENPROBE_PARALLEL`; zusammen 255. `make lint-gegenprobe` allein braucht 1 min
44 s.

| Zusage | Mutation (je in einer frischen Kopie des Repos) | roter Fall |
|---|---|---|
| Punkt 6: auch Testdateien | `lint.sh`: `go_dateien '*.go' \| grep -v '_test\.go$'` | `p6-testdatei` |
| Punkt 5: nur eine Datei namens `export_test.go` | `.golangci.yml`: `skip-regexp: export_test\.go$` | `p7-endet-auf-export-test` |
| Punkt 4: `dupl` nicht unter 150 | `threshold: 149` | `p4-dupl-gruen` |
| Punkt 4: `dupl` nicht über 150 | `threshold: 151` | `p4-dupl-rot` |
| Kopf: ungültiger Wert bricht ab | Prüfung `^[1-9][0-9]*$` durch `true` ersetzt | `einstellung-*` (alle fünf) |
| Kopf: leer ist nicht ungesetzt | `${…+gesetzt}` → `${…:+gesetzt}` | `einstellung-` |

*Nachtrag zur Verifikation* (`docs/reviews/2026-10-08-verifikation-slice-harness-lint.md`,
V-96, V-98, V-99, V-100 nach §6; Stand: `d3473b9` mit den Änderungen des Commits, der
diesen Absatz schreibt):

- V-96: zwei neue rote Fälle in `sammel`. `p6-leerzeichen-dann-strich` legt `// /nolint`
  in `internal/gegenprobe/nolint/fall.go` (Punkt 6: nach `//` beliebig viele
  Leerzeichen, Tabs oder `/`). `p8-why-ohne-doppelpunkt` ist eine Regel, deren
  Kommentarblock mit `# Why kein Doppelpunkt.` beginnt; erwartet ist die Zeile `Regel
  ohne Kommentarblock "# Why:" unmittelbar darüber` an ihrem Eintrag.
- V-98: Der Kopf der Gegenprobe führt `p6-testdatei` unter (6),
  `p7-endet-auf-export-test` unter (5) und (7) und die neuen Fälle unter (6) und (8).
  `lint.sh` (5) nennt die Fälle zu `testpackage`.
- V-99: Der Hilfetext von `make lint` sagt „Stufe lint des Dockerfile“ statt
  „netzlos ausser deps“. Die Zelle in `harness/README.md` sagt, dass die Gegenprobe
  den Lauf ohne Netz nicht prüft.
- V-100: Im Kopf der Gegenprobe steht unter OFFEN die Zeile zu den zulässigen Regeln
  nach §6. Unter *Grenze* in `harness/sensors/lint.md`, Punkt 1, steht der Satz aus der
  Grenze von `SPEC-049`. `lint.sh`, Profil und Stufe sind unverändert.

Punkt 6 hat damit 14 Fälle, Punkt 8 hat 48; zusammen 257.

| Zusage | Mutation (je in einer frischen Kopie des Repos) | roter Fall |
|---|---|---|
| Punkt 6: `/` zwischen Kommentarzeichen und `nolint` | `lint.sh`: `[ ${tab}/]*` → `[ ${tab}]*` (V7) | `p6-leerzeichen-dann-strich` |
| Punkt 8: erste Zeile beginnt mit `# Why:` | `lint.sh`: `/^[ \t]*# Why:/` → `/^[ \t]*# Why/` (V11) | `p8-why-ohne-doppelpunkt` |

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `9fe02ca` gesichtet (Architect,
2026-10-07; Zähler = Dateien unter `evidence/`). Treffer:

- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (13×, verkörpert in `AGENTS.md`
  §3.10) — trifft das Gate unmittelbar: Ein Lint-Gate, das nur grün gesehen wurde,
  belegt keine Zusage; daher der zweite Liefer-Punkt, je Zusage eine Mutation.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (18×, verkörpert in §3.11) — der
  Kopf von `lint.sh`, die Sensor-Datei, die Sensors-Zeilen und §3.2 sagen nur zu, was
  die Gegenprobe prüft.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (11×, verkörpert in §3.12) und
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (3×, verkörpert) — darum
  die Prüfung vor dem Code in §6.
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — dieser Slice ist
  wellenlos; §4 *Start* nennt den Architect-Schritt vor dem Code, er ist erfolgt.
- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×) — Risiko in §6; ein
  dritter Beleg aus diesem Slice ist eine Lücke.
- `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (2×) — verwandt: Die Liste der
  Zusagen ohne unterscheidende Mutation im Kopf der Gegenprobe muss vollständig sein;
  findet das Review eine fehlende, ist es dieselbe Klasse.
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) und
  `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — je ein Risiko in §6.
- `BEO-REPO/kern-fremdimporte-nur-ueber-tech-liste` (1×) — berührt, nicht
  aufgenommen (§1).
- `BEO-REPO/harness-lesart-ohne-entscheidungsort` — gestrichen; der Ort ist die
  Spezifikation, die Lesart die Sensor-Datei, wie hier geplant.

Keiner der Einträge erreicht vor dem Code die Schwelle 3× neu; keine neue Lücke vor dem
Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
