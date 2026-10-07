# Slice slice-lint-bestand-kern-driven: Lint-Bestand im Produkt-Code von Kern und Driven

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
Nutzers vom 2026-10-06 (Bereinigung vor dem Gate, ohne Stufen); Reihenfolge in §4
*Start*.

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode 1: Bestand ohne Befund vor dem Gate). Bindung an Entscheidungen: [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Messung des Bestands, Entscheidung 5: Bereinigung vor dem Gate, Kern und Driven; nach `SPEC-049` Punkt 9 16 Befunde), [ADR-0001](../../adr/0001-hexagonale-architektur.md) (Schichten des Hexagons; die Bereinigung verschiebt nichts über eine Schichtgrenze).

**Berührte Spec-Stellen:** [`SPEC-049`](../../../../spec/spezifikation.md#spec-049--lint-profil-lint) (Profil und Schwellen, gegen die bereinigt wird; der Slice ändert die Stelle nicht)

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

**Ziel:** `make lint` meldet unter `internal/hexagon/model`, `internal/hexagon/services`,
`internal/adapters/driven/postgres` und `internal/adapters/driven/recording` keinen
Befund; die dauerhaften Ausnahmen nach `SPEC-049` Punkt 8 blendet das Profil aus. Dafür
behebt der Slice die 16 Befunde im Produkt-Code dieser Pakete, ohne das Verhalten zu
ändern. Die Messung in
[ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Stand `79f40e1`;
Aufteilung je Paketgruppe in der Fassung `92d1b86` der ADR) zählt 13 mit
`uniq-by-line: true`; `make lint` zählt nach `SPEC-049` Punkt 9 jede Meldung, auch
mehrere auf derselben Zeile:

- **Komplexität (7) in vier Funktionen** — `gocognit`, `gocyclo` bzw. `cyclop` über
  der Schwelle aus Punkt 4: `Group.validate` (`internal/hexagon/model/extended.go`,
  alle drei), `(*cursor).objekte`
  (`internal/hexagon/services/replay.go`), `toResponse`
  (`internal/adapters/driven/postgres/upstream.go`, `gocyclo` und `cyclop`), `fromDTO`
  (`internal/adapters/driven/recording/yaml.go`).
- **`revive` (8)** — sieben Doc-Kommentare in `internal/hexagon/model` ohne die Form
  von `exported` bzw. `package-comments`, ein ungenutzter Parameter in
  `internal/adapters/driven/postgres`.
- **`staticcheck` (1)** — `QF1001` an einer Stelle dieser Pakete.

Maßgeblich ist die Liste, die `make lint` beim Start unter diesen Pfaden ausgibt; weicht
sie von der Aufzählung ab, gilt der Lauf, und der Plan folgt ihm (`AGENTS.md` §3.9).

Dazu übernimmt der Slice aus `slice-harness-blackbox-kern` (Verifikation V-81) eine
Zusage ohne Prüfung im Produkt-Code dieser Pakete: Der Kommentar über
`(*cursor).letzteNummer` (`internal/hexagon/services/replay.go`) sagt „ohne erwartete
Interaktion 0“ zu. Kein Test prüft das, und über `ReplayService` ist der Zweig nicht
erreichbar, weil `NewReplayService` keine Session ohne erwartete Interaktion aufnimmt
(`AGENTS.md` §3.11). Der Slice streicht Zweig und Satz oder benennt den Fall als
Invariante; welches von beiden, steht vor dem Code in §6.

**Herkunft:** Entscheidung des Nutzers vom 2026-10-06: Der Bestand wird vor dem Gate
bereinigt, ohne Stufen; die 32 Befunde im Produkt-Code (`SPEC-049` Punkt 9) tragen zwei Bereinigungs-Slices,
dieser für Kern und Driven, `slice-lint-bestand-driving` für die treibende Seite.
Geschnitten ist nach Paketgruppe wie bei den Umstellungs-Slices: Domain Model und
Services gelten dort als eine Schicht, der Kern; dieser Slice berührt damit zwei, Kern
und Driven.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Befunde unter `internal/adapters/driving/`, `internal/bootstrap` und `cmd/` —
  übernimmt `slice-lint-bestand-driving`.
- Befunde in Testdateien — übernehmen die Umstellungs-Slices
  (`slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`), die vor diesem
  Slice liegen; beim Start sind die Testdateien dieser Pfade ohne Befund.
- Eine neue Ausnahme in `.golangci.yml`, eine Lockerung einer Schwelle oder ein
  `//nolint` — gibt es nicht (Entscheidung 5, `AGENTS.md` §3.6); lässt sich ein Befund
  nicht ohne sie beheben, geht er an den Architect (§4).
- Verhaltensänderung, neue Fälle, geänderte Erwartungen in Tests — ein anderer
  Vorgang; die Bereinigung ist ein Umbau, gemessen an denselben Tests.
- Schnittstellen der Ports und exportierte Signaturen — Schicht-Abgrenzung: Der Slice
  ändert unexportierte Funktionen und Kommentare; eine exportierte Signatur, die sich
  ändern müsste, ist ein Befund für den Architect.
- Lastenheft, Spezifikation, Harness — kein Gate, kein Profil, keine Spec-Stelle wird
  geändert.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Komplexität: Die vier Funktionen aus §1 liegen unter den Schwellen aus
      `SPEC-049` Punkt 4, ohne Änderung des Verhaltens: Die Liste der Tests
      (`go test -list .` je Paket, im gepinnten Go-Image) ist vor und nach dem Umbau
      gleich, keine Erwartung ist geändert, `make test` und `make test-integration` sind
      grün.
- [ ] `make lint` meldet unter `internal/hexagon/model`, `internal/hexagon/services`,
      `internal/adapters/driven/postgres` und `internal/adapters/driven/recording` keinen
      Befund, ohne neue Ausnahme in `.golangci.yml` und ohne `//nolint`; die Zeilen der
      Ausgabe unter diesen Pfaden vor und nach dem Slice stehen im Bericht.
- [ ] `(*cursor).letzteNummer` sagt nichts zu, was kein Test prüft (V-81 aus
      `slice-harness-blackbox-kern`): Zweig `len(c.session.Interactions) == 0` und der
      Satz „ohne erwartete Interaktion 0“ sind gestrichen, oder der Kommentar nennt den
      Fall als Invariante mit ihrem Grund (`NewReplayService` nimmt keine Session ohne
      erwartete Interaktion auf) statt als Zusage — nach der Entscheidung in §6; die
      Liste der Tests bleibt gleich.
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
| `internal/hexagon/model/extended.go` | refactor | `Group.validate` unter die Schwelle, etwa durch Herauslösen der Prüfung je Nachrichtenart in unexportierte Funktionen |
| `internal/hexagon/model/*.go` (Doc-Kommentare) | update | sieben Kommentare in der Form von `exported` bzw. `package-comments`; nur Kommentare |
| `internal/hexagon/services/replay.go` | refactor | `(*cursor).objekte` unter die Schwelle |
| `internal/hexagon/services/replay.go` | update | `(*cursor).letzteNummer`: Zweig und Satz „ohne erwartete Interaktion 0“ gestrichen oder als Invariante benannt (V-81, §6) |
| `internal/adapters/driven/postgres/upstream.go` | refactor | `toResponse` unter die Schwelle (`gocyclo`), etwa als Tabelle oder Aufteilung je Nachrichtenart; ungenutzter Parameter entfernt oder als `_` benannt, wo eine Signatur ihn verlangt |
| `internal/adapters/driven/recording/yaml.go` | refactor | `fromDTO` unter die Schwelle |
| Stelle des `QF1001` | update | Vereinfachung nach dem Vorschlag von `staticcheck`, ohne andere Wirkung |
| Testdateien dieser Pakete | unverändert | Beleg des unveränderten Verhaltens; ein Test, der wegen des Umbaus geändert werden müsste, ist ein Befund (§4). Ausgenommen ist `export_test.go`, falls eine dort weitergereichte unexportierte Funktion umbenannt wird; dann nur der Verweis |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-blackbox-einstieg` liegt in `done/`,
und mit ihm alle vier Umstellungs-Slices und `slice-harness-lint-werkzeug` (WIP-Limit
1). Grund: Komplexitäts-Umbauten an unexportierten Funktionen brechen erst nach der
Umstellung keine White-Box-Tests mehr ([ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 5), und gemessen wird mit
dem Werkzeug. Reihenfolge nach Entscheidung des Nutzers vom 2026-10-06:
`slice-harness-lint-werkzeug`, die vier Umstellungs-Slices, dieser Slice,
`slice-lint-bestand-driving`, `slice-harness-lint`. Vor dem ersten Code-Commit prüft
der Architect §6 (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die vier Komplexitäts-Umbauten
  sind zusammen nicht in einer Review-Sitzung prüfbar; dann trägt dieser Slice Kern
  (`model`, `services`), ein eigener Slice Driven (`postgres`, `recording`).
- `in-progress` → `open` (blockiert — Carveout?): Ein Befund lässt sich nur mit einer
  Verhaltensänderung, einer geänderten exportierten Signatur, einer geänderten
  Testerwartung oder einer neuen Ausnahme beheben; dann zuerst eine Entscheidung des
  Architect.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, `make lint` ohne Befund unter den vier Pfaden,
Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — der Slice liefert keinen neuen Vertrag; Profil,
Schwellen und Ausnahmen stehen in `SPEC-049`. Für den Umbau gilt:

- **Komplexitäts-Bereinigung ohne Verhaltensänderung** — eine Funktion wird zerlegt
  oder als Tabelle geschrieben, nicht neu entworfen: gleiche Ergebnisse, gleiche
  Fehler und Fehlertexte, gleiche Reihenfolge der Prüfungen, wo sie beobachtbar ist
  (etwa welcher von zwei Fehlern zuerst gemeldet wird). Ändert ein Umbau eines davon,
  geht er an den Architect, nicht in den Diff.
- **Zusage ohne Prüfung in `letzteNummer`** (V-81 aus `slice-harness-blackbox-kern`) —
  offen; der Architect entscheidet vor dem ersten Code-Commit (§4) zwischen zwei Wegen:
  Zweig und Satz streichen (ein Cursor ohne erwartete Interaktion entsteht im Produkt
  nicht, Verifikation von `slice-harness-blackbox-kern`, Abschnitt 3; der Zugriff auf
  das letzte Element setzt dann die Invariante voraus) oder den Fall im Kommentar als
  Invariante mit Grund benennen und den Zweig stehen lassen. Einen Test bekommt keiner
  der beiden Wege: Über die Schnittstelle ist der Fall nicht herzustellen, und eine
  Brücke oder ein Test legt keinen `cursor` an (`SPEC-049` Punkt 7). Kein
  beobachtbares Verhalten ändert sich; die Testliste bleibt gleich.
- **Kontexte** — betreffen diesen Slice nicht; die Befunde von `contextcheck`,
  `containedctx` und `noctx` im Produkt-Code liegen in `slice-lint-bestand-driving`.

**Risiken:**

- **Verhalten ändert sich unbemerkt** — ein Pfad der umgebauten Funktion ist von keinem
  Test erreicht, die Testliste bleibt gleich und grün. Je umgebaute Funktion eine
  Mutation in einem ihrer Zweige, die ein bestehender Test vor und nach dem Umbau fängt;
  fängt keiner, ist das ein Befund für `slice-harness-coverage`, kein neuer Test hier.
  — **Ausgang:** — (bei Closure)
- **Befund verschoben statt behoben** — eine Funktion fällt unter die Schwelle, weil
  ihre Logik in eine Hilfsfunktion wandert, die selbst knapp darunter liegt, oder ein
  Parameter heißt `_`, obwohl er überflüssig ist. Review prüft die Form, `make lint`
  nur die Zahl (Grenze von `SPEC-049`). — **Ausgang:** — (bei Closure)
- **Aufzählung weicht vom Lauf ab** — die Liste in §1 stammt aus der Messung am Stand
  `79f40e1`; Code seitdem kann Befunde verschoben haben. Maßgeblich ist der Lauf beim
  Start (§1). — **Ausgang:** — (bei Closure)

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
`docs/plan/planning/observations/BEO-REPO/` am Stand `0d464b4` gesichtet (Zähler =
Dateien unter `evidence/`). Treffer:

- `BEO-REPO/plan-folgt-korrektur-nicht` (12×, verkörpert in `AGENTS.md` §3.9) — die
  Aufzählung in §1 stammt aus einer älteren Messung; weicht der Lauf ab, folgt der Plan
  im selben Commit.
- `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (2×) — ein
  Komplexitäts-Umbau entscheidet leicht still, welcher Fehler zuerst kommt; §6 nennt
  das als Randform, die an den Architect geht. Ein drittes Auftreten in diesem Slice
  erreichte die Schwelle.
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — verwandt mit dem Risiko, dass
  ein Umbau eine Prüfung verliert (§6).
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — wellenlos; §4 nennt den
  Architect vor dem Code.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
