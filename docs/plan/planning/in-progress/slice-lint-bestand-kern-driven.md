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
  der Schwelle aus Punkt 4: `Group.validate` (`internal/hexagon/model/extended.go:125`,
  alle drei: 22, 19, 19), `(*cursor).objekte`
  (`internal/hexagon/services/replay.go:320`, `gocognit` 28), `toResponse`
  (`internal/adapters/driven/postgres/upstream.go:226`, `gocyclo` 18 und `cyclop` 19),
  `fromDTO` (`internal/adapters/driven/recording/yaml.go:538`, `gocognit` 22).
- **`revive` (8)** — sieben Doc-Kommentare zu Konstanten in `internal/hexagon/model`,
  alle Regel `exported`: `TargetStatement` ohne Kommentar am Block
  (`extended.go:32`), sechs Kommentare zu `EndClosed` bis `EndFailed` nicht in der Form
  „`<Name> …`“ (`recording.go:80` bis `:93`); dazu `unused-parameter` für `ctx` in
  `(*session).Query` (`internal/adapters/driven/postgres/upstream.go:100`).
- **`staticcheck` (1)** — `QF1001` (De Morgan) in `Group.validate`
  (`internal/hexagon/model/extended.go:149`).

Gemessen mit `make lint` am Stand `c42afbe`, vor dem ersten Code-Commit: 32 Zeilen im
Produkt-Code, davon genau diese 16 unter den vier Pfaden, keine in Testdateien. Weicht
ein späterer Lauf ab, gilt der Lauf, und der Plan folgt ihm (`AGENTS.md` §3.9).

Dazu übernimmt der Slice aus `slice-harness-blackbox-kern` (Verifikation V-81) eine
Zusage ohne Prüfung im Produkt-Code dieser Pakete: Der Kommentar über
`(*cursor).letzteNummer` (`internal/hexagon/services/replay.go`) sagt „ohne erwartete
Interaktion 0“ zu. Kein Test prüft das, und über `ReplayService` ist der Zweig nicht
erreichbar, weil `NewReplayService` keine Session ohne erwartete Interaktion aufnimmt
(`AGENTS.md` §3.11). Der Slice benennt den Fall im Kommentar als Invariante mit ihrem
Grund und lässt den Zweig stehen; entschieden vor dem Code in §6.

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
  Vorgang; die Bereinigung ist ein Umbau, gemessen an denselben Tests. Ausgenommen
  sind Charakterisierungstests: Hält kein Test fest, welchen von zwei gleichzeitigen
  Fehlern eine umgebaute Funktion meldet, legt der Slice vor dem Umbau einen Test an,
  der den Bestand am alten Code festhält (gegen ihn grün); er sagt nichts Neues zu
  und steht in §7. **Bestätigt** (Architect, 2026-10-07, nach Review F-462; die
  Ausnahme schrieb der Implementer in `cedd891`): Ohne einen solchen Test wäre die
  Randform *gleiche Reihenfolge der Prüfungen* in §6 nicht prüfbar, und der Umbau
  ließe sich nur gegen Tests messen, die je Fall eine Regel verletzen. Grenze: nur
  die Reihenfolge zweier gleichzeitiger Fehler in einer umgebauten Funktion; im
  eigenen Commit vor dem Umbau, am alten Code grün; ohne Abdeckungs-Deklaration;
  der Kopfkommentar sagt, dass er den Bestand festhält und nichts zusagt (§6).
  Tests für überlebende Mutanten fallen nicht darunter, sie haben eine eigene Adresse
  (§6, Risiko *Verhalten ändert sich unbemerkt*).
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
      `slice-harness-blackbox-kern`): Der Satz „ohne erwartete Interaktion 0“ ist
      gestrichen; der Kommentar nennt stattdessen die Invariante mit ihrem Grund
      (`zuordnen` vergibt nur Sessions aus `frei`, `NewReplayService` nimmt dort nur
      Sessions mit erwarteter Interaktion auf) und sagt für den Zweig
      `len(c.session.Interactions) == 0` keinen Wert zu; der Zweig bleibt im Code
      (Entscheidung in §6), die Liste der Tests bleibt gleich.
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
| `internal/hexagon/model/extended.go` | refactor | `Group.validate` unter die Schwelle: die Prüfung der Client-Nachrichten in `validateClient`, die der Server-Nachrichten in `validateServer` herausgelöst; `validate` ruft sie in der bisherigen Reihenfolge der Prüfungen |
| `internal/hexagon/model/extended.go`, `recording.go` (Doc-Kommentare) | update | ein Kommentar am Block `TargetStatement`/`TargetPortal`, sechs Kommentare `EndClosed` bis `EndFailed` in der Form „`<Name> …`“ (`exported`); nur Kommentare, Inhalt gleich |
| `internal/hexagon/services/replay.go` | refactor | `(*cursor).objekte` unter die Schwelle: das Nachspielen einer Extended-Interaktion bis zum Cursor in `(objekte).extendedNachspielen` herausgelöst, die Haltebedingung am Cursor unverändert als Funktion übergeben |
| `internal/hexagon/services/replay.go` | update | `(*cursor).letzteNummer`: nur der Kommentar — Satz „ohne erwartete Interaktion 0“ ersetzt durch die Invariante mit Grund, Zweig unverändert (V-81, §6) |
| `internal/adapters/driven/postgres/upstream.go` | refactor | `toResponse` unter die Schwellen (`gocyclo`, `cyclop`): die Antworten ohne Felder in `ohneFelder`, die Spalten einer RowDescription in `spalten`, die Werte einer DataRow in `werte` herausgelöst; `ctx` in `(*session).Query` als `_` benannt — die Signatur verlangt `driven.UpstreamSession` (§6) |
| `internal/adapters/driven/recording/yaml.go` | refactor | `fromDTO` unter die Schwelle: die Prüfung je Session in `sessionFromDTO`, die je Interaktion in `geprueftFromDTO` herausgelöst, in der bisherigen Reihenfolge |
| `internal/hexagon/model/extended.go:149` (`QF1001`) | update | `!(letzte && si == len(g.Server)-1)` nach De Morgan zu `!letzte \|\| si != len(g.Server)-1` umgeschrieben, gleichwertig; steht nach dem Umbau in `validateServer` |
| Testdateien dieser Pakete | unverändert | Beleg des unveränderten Verhaltens; ein Test, der wegen des Umbaus geändert werden müsste, ist ein Befund (§4). Ausgenommen ist `export_test.go`, falls eine dort weitergereichte unexportierte Funktion umbenannt wird; dann nur der Verweis |
| `internal/hexagon/model/extended_test.go`, `internal/adapters/driven/recording/yaml_test.go` | add | Charakterisierungstests vor dem Umbau (§1): `TestValidateFehlerReihenfolge` hält fest, welchen von zwei Fehlern `Group.validate` meldet, `TestUnmarshalFehlerReihenfolge` dasselbe für `fromDTO`; im eigenen Commit vor dem Umbau, gegen den alten Code grün. Die Testliste „vor dem Umbau“ der DoD ist die mit ihnen; ohne Abdeckungs-Deklaration, die Abdeckungstabellen bleiben gleich. `toResponse` hat einen Fehlerpfad, `(*cursor).objekte` keinen; dort gibt es keine Reihenfolge zweier Fehler |

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
- **Reihenfolge zweier gleichzeitiger Fehler in `Group.validate` und `fromDTO`** —
  **entschieden** (Architect, 2026-10-07, nach Review F-463): kein Vertrag; die
  Charakterisierungstests halten den Bestand fest, sie sagen nichts zu. Nach außen
  geht nur der Text: Jeder Fall endet mit `PGR-E3003` und Exit-Code 3, der Text nennt
  die gemeldete Verletzung. Welche von zwei Verletzungen er nennt, entscheidet weder
  die Spezifikation noch eine ADR, und kein Nutzer braucht eine bestimmte; sie zum
  Vertrag zu machen, legte jede künftige Änderung der Prüfungen auf diese Reihenfolge
  fest, ohne Bedarf. Deshalb keine Stelle der Spezifikation. Die Kopfkommentare von
  `TestValidateFehlerReihenfolge` und `TestUnmarshalFehlerReihenfolge` sagen statt
  „und welche, steht fest“ ausdrücklich, dass sie den Bestand festhalten und keine
  Zusage sind; wer die Reihenfolge bewusst ändert, passt den Test an. Die
  Doc-Kommentare von `Group.validate`, `sessionFromDTO` und `geprueftFromDTO`
  beschreiben die Reihenfolge des Codes und bleiben (`AGENTS.md` §3.7); die
  Charakterisierungstests prüfen sie (§3.11). Nur Kommentare ändern sich, keine
  Erwartung.
- **Kommentar über die Bestätigungen in `extendedNachspielen`** — **entschieden**
  (Architect, 2026-10-07, nach Review F-460): enger fassen, kein Test hier. Der Satz
  „Die ersten n Nachrichten einer Art mit n Bestätigungen in der Interaktion gelten
  als angenommen“ entfällt, denn ohne `bestaetigt[m.Type]--` bleibt jeder Test grün.
  An seiner Stelle steht nur, was ein Test fängt: dass die Bestätigungen einer Art in
  der ganzen Interaktion gezählt werden, mit dem Grund (späte Bestätigung in der
  folgenden Gruppe, `LH-FA-18.a`); gefangen von der Mutation *nur die erste Gruppe
  gezählt* (§7). Ein weiterer Satz, etwa dass eine Art ohne Bestätigung nicht
  nachgespielt wird, bleibt nur, wenn eine Mutation an der Bedingung
  `bestaetigt[m.Type] > 0` einen Test rot macht. Die Zusage über die ersten n kommt
  mit dem Test aus `slice-tests-ueberlebende-mutanten` zurück.
- **Zusage ohne Prüfung in `letzteNummer`** (V-81 aus `slice-harness-blackbox-kern`) —
  **entschieden** (Architect, 2026-10-07, vor dem ersten Code-Commit): Der Kommentar
  benennt den Fall als Invariante mit Grund, der Zweig bleibt stehen. Der Satz „ohne
  erwartete Interaktion 0“ entfällt; an seiner Stelle steht sinngemäß: *Die Session
  eines Cursors hat mindestens eine erwartete Interaktion — `zuordnen` vergibt nur
  Sessions aus `frei`, und `NewReplayService` nimmt dort nur solche auf; die Prüfung auf
  die leere Liste schützt allein den Index.* Kein Wert für den leeren Fall wird
  zugesagt (`AGENTS.md` §3.11). Grund der Wahl: Streichen (a) machte aus einer
  verletzten Invariante einen Indexfehler mitten in der Diagnose einer Abweichung;
  ein Panic kommt im Produkt-Code nirgends vor, und ein interner Fehler
  (`CodeInternal`, wie in `zuordnen`) verlangte eine geänderte Signatur von
  `letzteNummer` an zwei Aufrufstellen samt einem weiteren unerreichbaren Zweig — mehr
  Umbau für denselben Schutz. Ort der Entscheidung ist dieser Abschnitt, keine Stelle
  der Spezifikation: Der Fall ist an keiner Schnittstelle beobachtbar, und die
  beobachtbare Seite — Sessions ohne Interaktion werden übersprungen, eine Aufzeichnung
  ohne Session mit Interaktion ist `PGR-E3004` — steht schon in der Spezifikation und
  ist getestet (`replay_test.go`, `replay_unverbraucht_test.go`,
  `replay_lebendpruefung_test.go`). **Akzeptiertes Negativ:** Den Zweig hält kein Test,
  und eine Mutation darin bleibt grün; über die Schnittstelle ist der Fall nicht
  herzustellen, und eine Brücke oder ein Test legt keinen `cursor` an (`SPEC-049`
  Punkt 7). Kein Befund für `slice-harness-coverage`. Kein beobachtbares Verhalten
  ändert sich; die Testliste bleibt gleich.
- **Ungenutzter Parameter `ctx` in `(*session).Query`** — **entschieden** (Architect,
  2026-10-07): Er heißt `_`. Entfernen ginge nicht ohne geänderte Signatur, die
  `driven.UpstreamSession` vorgibt (§1, Schicht-Abgrenzung); der Port sagt für `Query`
  keine Beachtung des Kontexts zu, `_` beschreibt also den Bestand. Den Kontext zu
  beachten wäre eine Verhaltensänderung und gehört nicht in diesen Slice. Das ist kein
  Fall von *Befund verschoben* (Risiko unten): Der Parameter ist nicht überflüssig,
  die Schnittstelle verlangt ihn.
- **`QF1001` und Doc-Kommentare** — keine Randform: De Morgan ist gleichwertig, die
  Kommentare ändern nur ihre Form; eine Ausnahme in `.golangci.yml` braucht keiner der
  16 Befunde.
- **Kontexte** — betreffen diesen Slice nicht; die Befunde von `contextcheck`,
  `containedctx` und `noctx` im Produkt-Code liegen in `slice-lint-bestand-driving`.

**Risiken:**

- **Verhalten ändert sich unbemerkt** — ein Pfad der umgebauten Funktion ist von keinem
  Test erreicht, die Testliste bleibt gleich und grün. Je umgebaute Funktion eine
  Mutation in einem ihrer Zweige, die ein bestehender Test vor und nach dem Umbau fängt;
  fängt keiner, ist das kein neuer Test hier: Adresse `slice-tests-ueberlebende-mutanten`
  (Entscheidung des Nutzers vom 2026-10-07; Architect nach Review F-459 und F-461).
  `slice-harness-coverage` nimmt nicht an, sein §1 schließt Tests über dem gemessenen
  Stand aus, und seine Messung zeigt einen Mutanten auf einer abgedeckten Zeile nicht;
  `slice-harness-mutation` schließt Tests für heute überlebende Mutanten aus. Die
  Adresse trägt vier verhaltensändernde grüne Mutanten, alle schon am Code vor dem
  Umbau grün: `extendedNachspielen` ohne `bestaetigt[m.Type]--`, `spalten` ohne
  `TableOID`, `spalten` ohne `ColumnNumber`, die Session-Nummer in den Fehlern einer
  Interaktion ab Session 2 (`geprueftFromDTO`). Sie trägt, wenn der Slice bei der
  Closure dieses Slice angelegt und nicht geschlossen ist und die vier nicht
  ausschließt.
  — **Ausgang:** — (bei Closure)
- **Befund verschoben statt behoben** — eine Funktion fällt unter die Schwelle, weil
  ihre Logik in eine Hilfsfunktion wandert, die selbst knapp darunter liegt, oder ein
  Parameter heißt `_`, obwohl er überflüssig ist. Review prüft die Form, `make lint`
  nur die Zahl (Grenze von `SPEC-049`). — **Ausgang:** — (bei Closure)
- **Aufzählung weicht vom Lauf ab** — die Liste in §1 stammt aus der Messung am Stand
  `79f40e1`; Code seitdem kann Befunde verschoben haben. Maßgeblich ist der Lauf beim
  Start (§1). — **Ausgang:** entfallen — der Lauf am Start (`c42afbe`) deckt sich mit
  der Aufzählung: 16 Befunde unter den vier Pfaden, keiner in Testdateien (§1).

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

**Belege des Implementers** (Stand `d840888`; Charakterisierungstests in `cedd891`, Umbau in
`d840888`):

*`make lint`, Zeilen unter den vier Pfaden.* Vorher (`7ac016b`, Exit 2, `32 issues:`
modulweit):

```text
internal/adapters/driven/postgres/upstream.go:226:1: calculated cyclomatic complexity for function toResponse is 19, max is 15 (cyclop)
internal/hexagon/model/extended.go:125:1: calculated cyclomatic complexity for function validate is 19, max is 15 (cyclop)
internal/adapters/driven/recording/yaml.go:538:1: cognitive complexity 22 of func `fromDTO` is high (> 20) (gocognit)
internal/hexagon/model/extended.go:125:1: cognitive complexity 22 of func `(Group).validate` is high (> 20) (gocognit)
internal/hexagon/services/replay.go:320:1: cognitive complexity 28 of func `(*cursor).objekte` is high (> 20) (gocognit)
internal/adapters/driven/postgres/upstream.go:226:1: cyclomatic complexity 18 of func `toResponse` is high (> 15) (gocyclo)
internal/hexagon/model/extended.go:125:1: cyclomatic complexity 19 of func `(Group).validate` is high (> 15) (gocyclo)
internal/adapters/driven/postgres/upstream.go:100:25: unused-parameter: parameter 'ctx' seems to be unused, consider removing or renaming it as _ (revive)
internal/hexagon/model/extended.go:32:2: exported: exported const TargetStatement should have comment (or a comment on this block) or be unexported (revive)
internal/hexagon/model/recording.go:80:2: exported: comment on exported const EndClosed should be of the form "EndClosed ..." (revive)
internal/hexagon/model/recording.go:82:2: exported: comment on exported const EndTerminate should be of the form "EndTerminate ..." (revive)
internal/hexagon/model/recording.go:84:2: exported: comment on exported const EndWriteFailed should be of the form "EndWriteFailed ..." (revive)
internal/hexagon/model/recording.go:86:2: exported: comment on exported const EndShutdown should be of the form "EndShutdown ..." (revive)
internal/hexagon/model/recording.go:89:2: exported: comment on exported const EndUnsupported should be of the form "EndUnsupported ..." (revive)
internal/hexagon/model/recording.go:93:2: exported: comment on exported const EndFailed should be of the form "EndFailed ..." (revive)
internal/hexagon/model/extended.go:149:41: QF1001: could apply De Morgan's law (staticcheck)
```

Nachher (`d840888`, Exit 2, `16 issues:` modulweit): unter den vier Pfaden keine Zeile.
Die 16 übrigen Zeilen sind dieselben wie vorher, Zeile für Zeile gleich (Vergleich der
sortierten Befundzeilen: keine neue), alle unter `internal/adapters/driving/`,
`internal/bootstrap` und `cmd/` (`slice-lint-bestand-driving`). Kein `//nolint`, keine
Änderung an `.golangci.yml` (der Diff berührt die Datei nicht).

*Komplexität nachher* (golangci-lint v2.14.0 im gepinnten Image, in einer Kopie außerhalb
des Repos mit Schwelle 1 je Linter, damit jede Funktion ihren Wert nennt; Schwellen nach
`SPEC-049` Punkt 4: `gocognit` > 20, `gocyclo` > 15, `cyclop` > 15):

| Funktion | `gocognit` | `gocyclo` | `cyclop` |
|---|---|---|---|
| `(Group).validate` | 6 | 7 | 7 |
| `(Group).validateClient` | 11 | 10 | 10 |
| `(Group).validateServer` | 7 | 6 | 6 |
| `(*cursor).objekte` | 12 | 10 | 10 |
| `(objekte).extendedNachspielen` | 12 | 7 | 7 |
| `toResponse` | 2 | 10 | 10 |
| `ohneFelder` | — | 7 | 7 |
| `spalten` | — | 2 | 2 |
| `werte` | 4 | 3 | 3 |
| `fromDTO` | 4 | 4 | 4 |
| `sessionFromDTO` | 6 | 6 | 6 |
| `geprueftFromDTO` | 5 | 6 | 6 |

Keine der herausgelösten Funktionen liegt knapp unter einer Schwelle (Risiko *Befund
verschoben*); ein Strich heißt, der Linter meldet für die Funktion auch bei Schwelle 1
nichts.

*Reihenfolge der Fehler vor dem Umbau geprüft.* `Group.validate` und `fromDTO` melden bei
zwei gleichzeitigen Verletzungen genau einen Fehler; kein Test hielt fest, welchen
(`TestValidateFehler` und `TestUnmarshalFehler` verletzen je Fall eine Regel). Dafür vor
dem Umbau, in `cedd891`, zwei **Charakterisierungstests** am alten Code, gegen ihn grün:

- `TestValidateFehlerReihenfolge` (`internal/hexagon/model/extended_test.go`), 13 Fälle,
  vergleicht den ganzen Fehlertext: Typ vor Stellung vor Zielart in derselben
  Client-Nachricht, eine frühere Client-Nachricht vor einer späteren, Client-Nachrichten
  vor sync/flush am Ende, sync/flush am Ende vor den Server-Nachrichten, Typ vor Stellung
  von ready_for_query über die Server-Nachrichten hinweg, Server-Nachrichten vor
  ready_for_query am Ende der letzten Gruppe.
- `TestUnmarshalFehlerReihenfolge` (`internal/adapters/driven/recording/yaml_test.go`),
  5 Fälle: Kennung vor Session ohne Interaktion, Nummer vor `offset_ms`, `offset_ms` vor
  der Form, eine frühere Interaktion vor einer späteren Nummer, eine frühere Session vor
  einer späteren Kennung.

`toResponse` hat einen Fehlerpfad (die nicht unterstützte Antwort), `(*cursor).objekte`
keinen; dort gibt es keine Reihenfolge zweier Fehler und keinen Charakterisierungstest.
Beide Tests tragen keine Abdeckungs-Deklaration; die Abdeckungstabellen bleiben gleich
(`make abdeckung-check` im Gate-Lauf grün).

*Testliste* (`go test -list .` je Paket im Image der Stufe `deps`, ohne Netz):

| Paket | `7ac016b` | `cedd891` (vor dem Umbau) | `d840888` (nach dem Umbau) |
|---|---|---|---|
| `internal/hexagon/model` | 9 | 10 (+ `TestValidateFehlerReihenfolge`) | 10 |
| `internal/hexagon/services` | 53 | 53 | 53 |
| `internal/adapters/driven/postgres` | 14 | 14 | 14 |
| `internal/adapters/driven/recording` | 13 | 14 (+ `TestUnmarshalFehlerReihenfolge`) | 14 |

`diff` der Listen `cedd891` gegen `d840888` leer; gegen `7ac016b` nur die beiden
Charakterisierungstests. Keine bestehende Erwartung ist geändert: Der Umbau-Commit
`d840888` berührt keine Testdatei.

*Mutationen* (`AGENTS.md` §3.10, Risiko *Verhalten ändert sich unbemerkt*): je Mutant ein
frischer Pfad außerhalb des Repos (Arbeitsbaum bzw. `git archive` per `tar -x --touch`,
kein `cp -p`), Ersetzung per Text, mtime der mutierten Datei neu gesetzt; Unit-Tests per
Bind-Mount im Image der Stufe `deps` (`go test -count=1 -v .` im Paket), Integrationstests
per `make test-integration` in der frischen Kopie (eigener Build-Kontext). Ohne Mutation
ist jeder genannte Test grün, mit ihr rot.

| Umbau | Zusage | Mutation | roter Test |
|---|---|---|---|
| `Group.validate` | sync/flush am Ende vor den Server-Nachrichten | Prüfung sync/flush nach `validateServer` | `TestValidateFehlerReihenfolge/sync_am_Ende_einer_vorderen_Gruppe_vor_einer_unbekannten_Server-Nachricht` |
| `Group.validate` | Server-Nachrichten vor ready_for_query am Ende | Prüfung am Ende vor `validateServer` | `TestValidateFehlerReihenfolge/unbekannte_Server-Nachricht_vor_fehlendem_ready_for_query`, `…/Stellung_von_ready_for_query_vor_fehlendem_ready_for_query_am_Ende` |
| `Group.validate` | Client- vor Server-Nachrichten | `validateServer` vor `validateClient` | `TestValidateFehlerReihenfolge/Client-Nachricht_vor_Server-Nachricht`, `…/ohne_Client-Nachricht_vor_einer_unbekannten_Server-Nachricht` |
| `validateClient` | Stellung vor Zielart in derselben Nachricht | Zielart-Prüfung vor Stellungs-Prüfung | `TestValidateFehlerReihenfolge/Stellung_vor_Zielart_in_derselben_Nachricht` |
| `validateClient` | Nachricht für Nachricht, nicht Regel für Regel | Typen aller Client-Nachrichten in einer eigenen Schleife zuerst | `TestValidateFehlerReihenfolge/Stellung_einer_früheren_vor_Typ_einer_späteren_Nachricht`, `…/Zielart_einer_früheren_vor_Typ_einer_späteren_Nachricht` |
| `validateServer` | Nachricht für Nachricht, nicht Regel für Regel | Typen aller Server-Nachrichten in einer eigenen Schleife zuerst | `TestValidateFehlerReihenfolge/Stellung_einer_früheren_vor_Typ_einer_späteren_Server-Nachricht` |
| `QF1001`, Bedingung *letzte Gruppe* | ready_for_query nur in der letzten Gruppe | `!letzte \|\|` gestrichen | `TestValidateFehler/ready_for_query_in_der_Flush-Gruppe` |
| `QF1001`, Bedingung *letzte Stelle* | ready_for_query nur an letzter Stelle | `\|\| si != len(g.Server)-1` gestrichen | `TestValidateFehler/ready_for_query_vor_dem_Ende_der_letzten_Gruppe` |
| `Group.validate` | sync genau am Ende der letzten Gruppe | `!=` zu `==` | `TestValidateFehler/letzte_Gruppe_endet_mit_flush_(abgeschnitten)` u. a. |
| `(*cursor).objekte` | die Nachricht am Cursor wird nicht nachgespielt | Halt erst bei `ni > c.nachricht` | `TestReplayExtendedDiagnoseLebensdauer` |
| `(*cursor).objekte` | Halt vor dem Nachspielen | Nachspielen vor der Halt-Prüfung | `TestReplayExtendedDiagnoseLebensdauer` |
| `(*cursor).objekte` | nach dem Halt kein Transaktionsende der laufenden Interaktion | Ergebnis von `extendedNachspielen` ignoriert | `TestReplayExtendedAbweichung`, `TestReplayExtendedDiagnoseAnweisung`, `TestReplayExtendedDiagnoseLebensdauer` |
| `(*cursor).objekte` | Bestätigungen je Interaktion gezählt | nur die erste Gruppe gezählt | `TestReplayExtendedAbweichung`, `TestReplayExtendedDiagnoseSpaeteBestaetigung` |
| `(*cursor).objekte` | eine Art ohne Bestätigung wird nicht nachgespielt (Doc-Kommentar `extendedNachspielen`) | `bestaetigt[m.Type] > 0` zu `>= 0` (frische Kopie aus `git archive 8d765ba`, Code gleich `d840888`) | `TestReplayExtendedDiagnoseLebensdauer` |
| `toResponse` | NULL bleibt NULL (`werte`) | NULL als leere Bytes | `TestOpenUndQuery` |
| `toResponse` | Werte sind Kopien (`werte`) | ohne Kopie | Integration: `TestE2EErgebnisartenEinfach`, `TestE2ERecordExtendedPgx` u. a. |
| `toResponse` | Format der Spalte (`spalten`) | `Format` nicht übernommen | Integration: `TestE2EErgebnisartenEinfach` |
| `toResponse` | Spalten der RowDescription | `Columns: nil` | Integration: `TestE2EReplaySelect1` u. a. |
| `ohneFelder` | empty_query_response | als no_data | Integration: `TestE2EReplayLebendpruefungDatabaseSQL`, `…Pgxpool`, `…WiePostgres` |
| `ohneFelder` | parse_complete | als bind_complete | `TestSendUndReceive`, `TestNachrichtenZwischenInteraktionen/Extended`, `TestSendNachFehlerantwort/ohne_ErrorResponse` |
| `ohneFelder` | bind_complete | als parse_complete | `TestSendUndReceive` |
| `ohneFelder` | close_complete | als parse_complete | `TestSendUndReceive` |
| `ohneFelder` | no_data | als empty_query_response | `TestSendUndReceive` |
| `ohneFelder` | portal_suspended | als command_complete | `TestSendUndReceive` |
| `toResponse` | eine andere Antwort ist nicht unterstützt | `ohneFelder` nimmt jede an | `TestReceiveFehler/COPY` |
| `toResponse` | Code der nicht unterstützten Antwort | `CodeInternal` statt `CodeUnsupported` | `TestNichtVermittelbar/COPY`, `TestReceiveFehler/COPY` |
| `fromDTO` | Kennung vor Session ohne Interaktion | Prüfungen getauscht | `TestUnmarshalFehlerReihenfolge/Kennung_vor_Session_ohne_Interaktion` |
| `fromDTO` | Nummer vor `offset_ms` | Prüfungen getauscht | `TestUnmarshalFehlerReihenfolge/Nummer_vor_offset_ms` |
| `fromDTO` | `offset_ms` vor der Form | Form zuerst | `TestUnmarshalFehlerReihenfolge/offset_ms_vor_der_Form` |
| `fromDTO` | Session für Session | alle Kennungen in einer eigenen Schleife zuerst | `TestUnmarshalFehlerReihenfolge/frühere_Session_vor_späterer_Kennung` |
| `fromDTO` | Interaktion für Interaktion | alle Nummern in einer eigenen Schleife zuerst | `TestUnmarshalFehlerReihenfolge/frühere_Interaktion_vor_späterer_Nummer` |
| `fromDTO` | Session ohne Interaktion nur mit `empty_sessions` | `emptySessions` immer falsch | `TestFelderDerVersion1` |
| `fromDTO` | die Form prüft `Validate` | `Validate` nicht gerufen | `TestUnmarshalExtendedFehler` (sechs Fälle, etwa `…/Zielart_leer`), `TestUnmarshalFehler` (drei Fälle, etwa `…/abgeschnitten_nach_request`) |
| `fromDTO` | die Interaktionen gehen in die Session | Interaktion nicht angehängt | `TestRoundtrip`, `TestUnmarshalNullUngequotet` |

*Grüne Mutanten, eingeordnet* (Schritt 19):

- **äquivalent** — Halt in `objekte` ohne `gi > c.gruppe` (nur `gi == c.gruppe && ni >=
  c.nachricht`): Am Cursor gilt `c.nachricht < len(Groups[c.gruppe].Client)`, weil
  `ClientMessage` die Nachricht nach dem letzten Client-Eintrag einer Gruppe auf 0 setzt
  und `Validate` jeder Gruppe eine Client-Nachricht verlangt; die Bedingung `gi ==
  c.gruppe && ni == c.nachricht` greift daher vor jeder späteren Gruppe. Kein Test kann
  ihn fangen. Die Grenze trägt die Randform *Komplexitäts-Bereinigung ohne
  Verhaltensänderung* in §6: Die Bedingung steht unverändert im Umbau. Am Code vor dem
  Umbau (`cedd891`) ebenso grün.
- **verhaltensändernd, über die Schnittstelle fangbar, ohne Test** — alle vier am Code
  vor dem Umbau (`7ac016b` bzw. `cedd891`) ebenso grün, also keine Lücke des Umbaus,
  sondern des Bestands; §1 schließt neue Fälle aus, darum je eine Test-Idee mit Grenze,
  Adresse `slice-tests-ueberlebende-mutanten` (Risiko *Verhalten ändert sich
  unbemerkt*, §6):
  - `extendedNachspielen` ohne `bestaetigt[m.Type]--`: Hat eine Interaktion ein
    bestätigtes und danach ein abgelehntes parse, gilt auch das abgelehnte als angenommen.
    Test-Idee: Replay einer Aufzeichnung mit zwei parse desselben Statement-Namens, das
    zweite mit ErrorResponse statt parse_complete, danach eine Abweichung an einem bind
    auf dieses Statement; die Diagnose nennt das SQL des ersten parse, nicht das des
    zweiten. Grenze: Kein Test erzeugt ein abgelehntes parse, bind oder close nach einem
    angenommenen derselben Art in einer Interaktion. Der Doc-Kommentar von
    `extendedNachspielen` ist dafür enger gefasst (§6, nach Review F-460): Der Satz über
    die ersten n Nachrichten entfällt; es bleiben die Zählung über die ganze Interaktion
    mit Grund (gefangen von *nur die erste Gruppe gezählt*) und „eine Art ohne
    Bestätigung in der Interaktion wird nicht nachgespielt“ (gefangen von
    `bestaetigt[m.Type] >= 0`, Mutationstabelle). Die Zusage über die ersten n kommt mit
    dem Test aus `slice-tests-ueberlebende-mutanten` zurück.
  - `spalten` ohne `TableOID`: Die Tabellen-OID einer Spalte ginge verloren. Test-Idee:
    Unit-Test in `postgres_test` mit einer RowDescription mit `TableOID` ungleich 0 über
    `Query`, Vergleich der Spalte; Grenze: Die Integrationstests vergleichen die
    Tabellen-OID nicht.
  - `spalten` ohne `ColumnNumber` (Review-Belege P2 an `ac2e662`, P2alt an `7ac016b`,
    je Unit-Tests von `postgres` und `make test-integration` grün): Die Spaltennummer
    in der Tabelle ginge verloren. Test-Idee: derselbe Unit-Test wie bei `TableOID`, mit
    `TableAttributeNumber` ungleich 0; Grenze: Kein Test vergleicht die Spaltennummer,
    die Integrationstests nicht und die Unit-Tests nicht.
  - Die Session-Nummer in den Fehlern einer Interaktion (Review-Belege Y2:
    `geprueftFromDTO(1, ii, id)` an `ac2e662`; Y2alt: `1` statt `sd.ID` in der
    Formmeldung an `cedd891`; je Unit-Tests von `recording` grün): Eine Meldung zu einer
    Interaktion in Session 2 oder später nennte Session 1. Test-Idee: `Unmarshal` einer
    Aufzeichnung mit zwei Sessions, deren zweite eine abgeschnittene Interaktion trägt,
    erwartet `PGR-E3003` mit „Session 2, Interaktion 1“; ebenso mit negativem `offset_ms`
    und falscher Nummer. Grenze: Kein Test erzeugt einen Interaktionsfehler aus
    `fromDTO` in einer Session ab 2; `TestVorpruefungNenntOrt` nennt Session 2, aber die
    Meldung stammt aus der Vorprüfung, nicht aus `fromDTO`.
- `ctx` heißt `_`, die Doc-Kommentare und der Kommentar über `letzteNummer`: keine
  Verhaltensänderung, keine Mutation; für den Zweig in `letzteNummer` gilt das
  akzeptierte Negativ aus §6.

*Läufe:* `gofmt -l internal` (leer), `go vet ./internal/...` und `go test -count=1
./internal/...` im Image der Stufe `deps` nach dem Umbau grün; `make lint` an `7ac016b`
(32), `cedd891` (32, Zeilen gleich) und `d840888` (16); `make gates` grün an `d840888`, an
`ac2e662` und am Stand des Commits, der die Kommentare nach den Entscheidungen zu F-460
und F-463 fasst (nur Kommentare, keine Erwartung und kein Testfall geändert).

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
