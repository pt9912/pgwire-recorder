# Slice slice-lint-bestand-driving: Lint-Bestand im Produkt-Code der treibenden Seite

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

**Bezug:** [`LH-QA-07`](../../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (Messmethode 1: Bestand ohne Befund vor dem Gate). Bindung an Entscheidungen: [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Messung des Bestands, Entscheidung 5: Bereinigung vor dem Gate, Driving; nach `SPEC-049` Punkt 9 16 Befunde, erwartete Bereinigung der Kontexte mit `context.WithoutCancel`), [ADR-0001](../../adr/0001-hexagonale-architektur.md) (Schichten des Hexagons), [ADR-0010](../../adr/0010-verwendung-von-pgproto3.md) (`pgproto3` im PGWire-Adapter).

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

**Ziel:** `make lint` meldet unter `internal/adapters/driving/pgwire`,
`internal/adapters/driving/cli` und `internal/bootstrap` keinen Befund; die dauerhaften
Ausnahmen nach `SPEC-049` Punkt 8 blendet das Profil aus. Damit hat der ganze Bestand
keinen Befund mehr, und `slice-harness-lint` kann das Gate anschließen. Dafür behebt der
Slice die 16 Befunde im Produkt-Code dieser Pakete, ohne das Verhalten zu ändern. Die
Messung in [ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) (Stand
`79f40e1`; Aufteilung je Paketgruppe in der Fassung `92d1b86` der ADR) zählt 12 mit
`uniq-by-line: true`; `make lint` zählt nach `SPEC-049` Punkt 9 jede Meldung, auch
mehrere auf derselben Zeile:

**Befundliste**, gemessen mit `make lint` am Stand `a453969` (Architect, 2026-10-07): 16
Befunde, dieselben wie in der Messung der ADR; 14 in `pgwire`, je einer in `cli` und
`internal/bootstrap`, `cmd/` ohne Befund. Pfade relativ zu `internal/`.

| # | Ort | Linter | Wert bzw. Gegenstand |
|---|---|---|---|
| 1 | `adapters/driving/pgwire/server.go:182` `(*Server).replaySitzung` | `gocognit` | 37, Schwelle 20 |
| 2 | ebenda | `gocyclo` | 20, Schwelle 15 |
| 3 | ebenda | `cyclop` | 21, Schwelle 15 |
| 4 | `adapters/driving/pgwire/server.go:435` `(*richtungen).clientRichtung` | `gocognit` | 22, Schwelle 20 |
| 5 | ebenda | `cyclop` | 17, Schwelle 15 |
| 6 | `adapters/driving/pgwire/server.go:569` `(*Server).startup` | `cyclop` | 17, Schwelle 15 |
| 7 | `adapters/driving/pgwire/server.go:739` `toMessage` | `gocyclo` | 19, Schwelle 15 |
| 8 | ebenda | `cyclop` | 20, Schwelle 15 |
| 9 | `adapters/driving/pgwire/server.go:302` | `containedctx` | Feld `ctx` im Struct `richtungen` |
| 10 | `adapters/driving/pgwire/server.go:440` | `contextcheck` | `r.s.recorder.Shutdown(r.ctx, …)` in `clientRichtung` |
| 11 | `adapters/driving/pgwire/server.go:462` | `contextcheck` | `r.s.recorder.Query(r.ctx, …)` in `clientRichtung` |
| 12 | `adapters/driving/pgwire/server.go:474` | `contextcheck` | `r.s.recorder.ClientMessage(r.ctx, …)` in `clientRichtung` |
| 13 | `bootstrap/bootstrap.go:71` | `contextcheck` | `service.Finish(context.Background())` in `record` |
| 14 | `adapters/driving/pgwire/server.go:62` | `noctx` | `net.Listen` in `pgwire.Listen` |
| 15 | `adapters/driving/pgwire/server.go:728` | `revive` | `unused-receiver`: `s` in `(*Server).send` |
| 16 | `adapters/driving/cli/cli.go:243` | `revive` | `unused-receiver`: `w` in `wahrheitswert.IsBoolFlag` |

Die drei `contextcheck`-Befunde in `pgwire` (10 bis 12) stehen an den Stellen, die den
Kontext aus dem Struct nehmen, während die Methode einen eigenen Parameter `ctx` hat; die
vier Stellen mit `context.WithoutCancel(ctx)` (`handle`, `recordSitzung`, `closeReplay`,
`closeRecord`) meldet der Lauf nicht. Maßgeblich bleibt der Lauf beim ersten
Code-Commit; weicht er von der Tabelle ab, folgt der Plan ihm (`AGENTS.md` §3.9).

**Herkunft:** Entscheidung des Nutzers vom 2026-10-06: Der Bestand wird vor dem Gate
bereinigt, ohne Stufen; die 32 Befunde im Produkt-Code (`SPEC-049` Punkt 9) tragen zwei Bereinigungs-Slices,
`slice-lint-bestand-kern-driven` für Kern und Driven, dieser für die treibende Seite.
Er berührt zwei Schichten, den Driving Adapter (`pgwire`, `cli`) und das Bootstrap, das
`pgwire.Listen` ruft.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Befunde unter `internal/hexagon/` und `internal/adapters/driven/` — übernimmt
  `slice-lint-bestand-kern-driven`, der vor diesem Slice liegt.
- Befunde in Testdateien — übernehmen `slice-harness-blackbox-pgwire` und
  `slice-harness-blackbox-einstieg` vor diesem Slice.
- Eine neue Ausnahme in `.golangci.yml`, eine Lockerung einer Schwelle oder ein
  `//nolint` — gibt es nicht (Entscheidung 5, `AGENTS.md` §3.6); lässt sich ein Befund
  nicht ohne sie beheben, geht er an den Architect (§4).
- Verhaltensänderung, insbesondere am Herunterfahren und an der Signalbehandlung, neue
  Fälle oder geänderte Erwartungen in Tests — ein anderer Vorgang; die Bereinigung ist
  ein Umbau, gemessen an denselben Tests (§6). Die Ausnahme für Charakterisierungstests
  aus `slice-lint-bestand-kern-driven` gilt hier **nicht** (Architect, 2026-10-07, vor
  dem ersten Code-Commit): Ihre Grenze ist die Reihenfolge zweier gleichzeitiger Fehler
  in einer umgebauten Funktion, die kein Test festhält. Hier hat nur `startup` eine
  solche Reihenfolge (Länge vor Startcode), und `TestFremdeErsteNachricht` hält sie
  fest: Eine HTTP-Zeile hat eine zu große Länge und einen unbekannten Startcode, der
  Test verlangt die Warnung ohne Antwort. `toMessage` liefert je Antworttyp höchstens
  einen Fehler; `replaySitzung` und `clientRichtung` ordnen Ereignisse der
  Nebenläufigkeit, keine zwei Fehler einer Eingabe. Findet der Implementer dennoch eine
  ungeprüfte Reihenfolge, legt er keinen Test an, sondern gibt sie dem Architect zurück.
- Ports und Services — Schicht-Abgrenzung: Eine Änderung an `internal/hexagon/` ist
  nicht Teil des Umbaus. Eine exportierte Signatur innerhalb dieser Pakete darf sich
  ändern, wenn alle Aufrufer in diesen Pfaden liegen (etwa `pgwire.Listen` mit
  Kontext, gerufen nur aus `internal/bootstrap`).
- Lastenheft, Spezifikation, Harness — kein Gate, kein Profil, keine Spec-Stelle wird
  geändert; den Anschluss an die Gate-Kette liefert `slice-harness-lint`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Komplexität: Die vier Funktionen aus §1 liegen unter den Schwellen aus
      `SPEC-049` Punkt 4, ohne Änderung des Verhaltens: Die Liste der Tests
      (`go test -list .` je Paket, im gepinnten Go-Image) ist vor und nach dem Umbau
      gleich, keine Erwartung ist geändert, `make test` und `make test-integration` sind
      grün. — Verifikation, Abschnitt 1 Punkt 1: die vier Funktionen und ihre neun
      Hilfsfunktionen höchstens bei 17 von 20 (`gocognit`, `replaySitzung`) bzw. 10 von 15
      (`gocyclo`) und 11 von 15 (`cyclop`), Werte gleich der Tabelle in §7; Testliste an
      `b7d4555` und `dd5d13d` byte-gleich (162 Unit-Tests, 39 Integrationstests); `git diff
      b7d4555 dd5d13d` ohne Testdatei, `docs/user/`, `.golangci.yml`, `Dockerfile`;
      Stichproben V7 bis V9 rot (Abschnitt 2).
- [x] Kontexte: Die sechs Befunde von `containedctx`, `contextcheck` und `noctx` sind
      behoben, ohne die Semantik des Herunterfahrens zu ändern (§6); die Tests und
      Integrationstests zum Abbruchsignal und zur Frist des Herunterfahrens laufen
      unverändert grün. — Verifikation, Abschnitt 1 Punkt 2: Form nach §6 *Kontexte*
      gelesen (F-470 bestätigt); V1 (`ctx.Err()` → `uc.Err()`) rot; die Äquivalenz von
      „`WithoutCancel` entfernt“ mit V5 (`Finish`, Unit und beide Integrationsphasen,
      darunter die Signaltests) und V6 bestätigt.
- [x] `make lint` meldet unter `internal/adapters/driving/pgwire`,
      `internal/adapters/driving/cli` und `internal/bootstrap` keinen Befund, ohne neue
      Ausnahme in `.golangci.yml` und ohne `//nolint`; über das ganze Modul meldet es
      keinen Befund mehr. Die Zeilen der Ausgabe unter diesen Pfaden vor und nach dem
      Slice stehen im Bericht. — Verifikation, Abschnitt 1 Punkt 3: an `b7d4555` 16
      Befunde, sortiert byte-gleich mit §7 *Vorher*; an `dd5d13d` Exit 0, `0 issues.`,
      keine `lint:`-Zeile; kein `//nolint`, `.golangci.yml`, `tools/harness/lint.sh` und
      `Dockerfile` unverändert. Die Zeilen stehen in §7 statt im Bericht.
- [x] `make gates` grün. — an `dd5d13d` (Verifikation, Abschnitt 5) und an `7cb298c`
      (Planner, vor dieser Closure; Go-Code gleich `893b54e`, das gegenüber `dd5d13d` nur
      den Doc-Kommentar von `replayWaechter` ändert). Diese Closure ändert nur Pläne,
      Roadmap und Register; `make docs-check` und `make kopf-check` grün am Stand dieser
      Closure.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      — `docs/reviews/2026-10-07-review-slice-lint-bestand-driving.md` (`934add1`; F-465
      MEDIUM, F-466 und F-467 LOW, F-468 bis F-470 INFO), Nacharbeit in `dd5d13d`;
      Verifikation `docs/reviews/2026-10-07-verifikation-slice-lint-bestand-driving.md`
      (`89fba63`; V-93 LOW, V-94 und V-95 Hinweise), V-93 umgesetzt in `893b54e`, V-94 und
      V-95 im Nehmer.
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
| `internal/adapters/driving/pgwire/server.go` | refactor | Befunde 1 bis 8: `replaySitzung`, `clientRichtung`, `startup` und `toMessage` unter die Schwellen, durch Herauslösen von Schritten in unexportierte Funktionen (§6 *Komplexitäts-Bereinigung*) |
| `internal/adapters/driving/pgwire/server.go` (`richtungen`) | refactor | Befunde 9 bis 12: Das Feld `ctx` entfällt; die Methoden, die den Kontext für den Use Case brauchen, bekommen ihn als Parameter, und `clientRichtung` reicht dem Use Case `context.WithoutCancel(ctx)` (§6 *Kontexte*); Befund 15: `send` mit unbenanntem Receiver `func (*Server) send` |
| `internal/adapters/driving/pgwire/server.go` (`Listen`) | update | Befund 14: `pgwire.Listen(ctx, address)` öffnet mit `(&net.ListenConfig{}).Listen(context.WithoutCancel(ctx), "tcp", address)`; der Meldungscode bleibt `PGR-E4001` (§6 *Listen*) |
| `internal/bootstrap/bootstrap.go` | update | `record` und `replay` reichen ihren `ctx` an `pgwire.Listen`; Befund 13: `service.Finish(context.WithoutCancel(ctx))` |
| `internal/adapters/driving/cli/cli.go` | update | Befund 16: `func (wahrheitswert) IsBoolFlag() bool` mit unbenanntem Receiver |
| Testdateien dieser Pakete und `test/integration` | unverändert | Beleg des unveränderten Verhaltens; ein Test, der wegen des Umbaus geändert werden müsste, ist ein Befund (§4). Keine Testdatei ruft `pgwire.Listen` (gemessen am Stand `a453969`). Ausgenommen ist `export_test.go`, falls eine dort weitergereichte unexportierte Funktion umbenannt wird; dann nur der Verweis |

**Größe** (Architect, 2026-10-07): ein Schnitt, kein zweiter Slice. Die 14 Befunde in
`server.go` sind vier Funktionen und ein Struct; die DoD hat drei Liefer-Punkte, der Slice
berührt den Driving Adapter und das Bootstrap. In einer Review-Sitzung prüfbar wird er
durch die Reihenfolge der Commits: zuerst die Kontexte und die beiden `revive`-Befunde
(9 bis 16), dann die Komplexität (1 bis 8), jeder mit `make test` grün. Die Rückführung
in §4 bleibt die Bedingung, falls der Umbau von `replaySitzung` allein den Rahmen sprengt.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-lint-bestand-kern-driven` liegt in `done/`,
und mit ihm die vier Umstellungs-Slices und `slice-harness-lint-werkzeug` (WIP-Limit 1).
Reihenfolge nach Entscheidung des Nutzers vom 2026-10-06:
`slice-harness-lint-werkzeug`, die vier Umstellungs-Slices,
`slice-lint-bestand-kern-driven`, dieser Slice, `slice-harness-lint`. Vor dem ersten
Code-Commit prüft der Architect §6, insbesondere ob `contextcheck` die Form
`context.WithoutCancel` annimmt (`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`).
Geprüft am Stand `a453969` (Architect, 2026-10-07); die Entscheidungen stehen in §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Umbau von
  `replaySitzung` allein füllt eine Review-Sitzung, oder Komplexität und Kontexte
  zusammen sind nicht in einer prüfbar; dann trägt ein eigener Slice
  `replaySitzung` (bzw. die Kontexte), dieser Slice den Rest.
- `in-progress` → `open` (blockiert — Carveout?): Ein Befund lässt sich nur mit einer
  Verhaltensänderung beheben, die Bereinigung eines Kontexts ändert die Semantik des
  Herunterfahrens, `contextcheck` nimmt `context.WithoutCancel` nicht an, oder ein
  Befund braucht eine neue Ausnahme; dann zuerst eine Entscheidung des Architect.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, `make lint` ohne Befund im ganzen Modul,
Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — der Slice liefert keinen neuen Vertrag; Profil,
Schwellen und Ausnahmen stehen in `SPEC-049`. Für den Umbau gilt, vor dem Code vom
Architect zu bestätigen:

- **Komplexitäts-Bereinigung ohne Verhaltensänderung** — `replaySitzung` (`gocognit`
  37) und die drei übrigen Funktionen werden zerlegt, nicht neu entworfen: gleiche
  Antworten an den Client in gleicher Reihenfolge, gleiche Meldungen und Codes, gleiche
  Behandlung von Verbindungsende, Abweichung und Abbruch. Die Nebenläufigkeit von
  Client- und Server-Richtung bleibt, wie sie ist. Ändert ein Umbau eines davon, geht
  er an den Architect, nicht in den Diff.
- **Kontexte** — **entschieden** (Architect, 2026-10-07, vor dem ersten Code-Commit):
  Die Bereinigung ist `context.WithoutCancel(ctx)` auf dem übergebenen Kontext
  ([ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 5); keine
  Folge-Entscheidung, keine Ausnahme. Je Befund:
  - *`richtungen` (9 bis 12).* Das Feld `ctx` entfällt. `recordSitzung` reicht
    `context.WithoutCancel(ctx)` als ersten Parameter an `serverRichtung`, und
    `beende`, `schreibe`, `fehler` und `wartenAufEnde` bekommen ihn ebenso.
    `clientRichtung(ctx)` behält den übergebenen Kontext allein für die Abfrage
    `ctx.Err()` und reicht an `Shutdown`, `Query`, `ClientMessage` und an die
    Methoden darunter `context.WithoutCancel(ctx)`. Das ist derselbe Kontext wie heute
    `r.ctx`: derselbe Elternkontext, dieselben Werte, kein Abbruch, keine Frist. Das
    Herunterfahren (`LH-FA-13.a`) hängt weiter allein an `ctx.Err()` in
    `clientRichtung` und am Wächter in `recordSitzung`; beide bleiben unverändert.
  - *`bootstrap.go:71` (13).* `service.Finish(context.WithoutCancel(ctx))`. An dieser
    Stelle ist `ctx` beendet; der gelöste Kontext hat wie `context.Background()` weder
    Abbruch noch Frist und trägt zusätzlich die Werte des Aufrufers, die kein Adapter
    liest.
  - *Gleich im Verhalten, belegt am Code:* Hinter den Ports liest auf diesem Weg kein
    Adapter den Kontext. `RecordService.Shutdown` und `Delivered` nehmen ihn als `_`,
    `Query`, `ClientMessage` und `AwaitServer` reichen ihn an `session.Query`, `Send`
    und `Receive` in `postgres/upstream.go`, die ihn nicht lesen, und `YAML.Write`
    nimmt ihn als `_`. Die Mutation „`WithoutCancel` entfernt“ an 10 bis 13 ist darum
    **äquivalent** und kein Befund (**akzeptiertes Negativ**): Kein Test kann sie
    fangen, solange kein Adapter den Kontext liest; liest ihn später einer, ist das
    eine Änderung dieses Adapters, die ihren eigenen Test bringt. An ihre Stelle tritt
    die Mutation, die das Herunterfahren trägt (Risiko *Herunterfahren ändert sich
    unbemerkt*).
- **`pgwire.Listen` mit Kontext (14)** — **entschieden** (Architect, 2026-10-07):
  `Listen(ctx, address)` öffnet mit
  `(&net.ListenConfig{}).Listen(context.WithoutCancel(ctx), "tcp", address)`. `net.Listen`
  tut heute dasselbe mit `context.Background()`; ein vor dem Öffnen beendeter `ctx`
  bricht das Öffnen also weiter nicht ab, und der Lauf endet wie bisher über `Serve`.
  Den Kontext ungelöst weiterzureichen wäre eine Verhaltensänderung: Bei einem Namen als
  Adresse (etwa `localhost:0`) bräche die Namensauflösung nach einem frühen Signal mit
  `PGR-E4001` ab statt mit dem regulären Ende. Der Doc-Kommentar von `Listen` sagt dazu
  nichts zu (`AGENTS.md` §3.11), denn kein Test prüft es. Mit einer IP-Adresse, wie sie
  alle Tests nutzen, liest `net` den Kontext nicht; die Mutation „`WithoutCancel`
  entfernt“ bleibt dann grün, ist aber über `bootstrap.Run` fangbar (Risiko *Verhalten
  ändert sich unbemerkt*).
- **Ungenutzte Receiver (15, 16)** — **entschieden** (Architect, 2026-10-07): unbenannter
  Receiver, `func (*Server) send` und `func (wahrheitswert) IsBoolFlag`. Die
  Methodenmengen bleiben gleich; `wahrheitswert` erfüllt weiter die Schnittstelle, an
  der `flag` eine Option ohne Wert erkennt. Keine Randform.
- **Rückgabewert von `weiterlesen`** (aus `slice-harness-blackbox-pgwire`, V-86) —
  **entschieden** (Architect, 2026-10-07): Der Umbau berührt ihn nicht. Der Lesefehler
  mit abgelaufener Frist ruft nach dem Umbau weiter `weiterlesen()`, verwirft den Wert
  und liest erneut; die Zerlegung von `clientRichtung` liest ihn nicht und entfernt ihn
  nicht. Der Testbedarf aus V-86 entsteht damit nicht, die Mutanten an ihm bleiben
  äquivalent (`slice-tests-ueberlebende-mutanten` §1 schließt sie aus, F-449). Liest ein
  Umbau den Wert doch, ist das eine Verhaltensänderung und geht an den Architect, nicht
  in den Diff.
- **Reihenfolge zweier gleichzeitiger Fehler** — keine offene Randform: §1 *Ausdrücklich
  NICHT*, Absatz Verhaltensänderung.

**Risiken:**

- **`contextcheck` nimmt `context.WithoutCancel` nicht an** — der Code nutzt die Form
  schon an vier Stellen (`handle`, `recordSitzung`, `closeReplay`, `closeRecord` in
  `server.go`), und die Messung meldet trotzdem drei Befunde in `pgwire`. Der Lauf am
  Stand `a453969` zeigt: Die drei stehen an den Zugriffen auf `r.ctx` (§1, 10 bis 12),
  keine an den vier Stellen mit `WithoutCancel`. Meldet der Lauf nach der Bereinigung
  dennoch einen, trägt die erwartete Bereinigung nicht, und der Weg ist eine
  Entscheidung des Architect (§4, `→ open`).
  — **Ausgang:** entfallen: Nach `78f28a1` meldet `make lint` keinen `contextcheck` mehr,
  an `dd5d13d` modulweit `0 issues.` (§7, Verifikation Abschnitt 1 Punkt 3);
  `contextcheck` nimmt die Form an.
- **Herunterfahren ändert sich unbemerkt** — ein Kontext, der bisher den Abbruch des
  Aufrufers überlebte, bricht nach der Bereinigung mit ihm ab (oder umgekehrt); das
  Recording wird dann beim Signal nicht mehr oder anders geschrieben. Die
  Integrationstests zum Abbruchsignal laufen vor und nach dem Umbau. Die Mutation
  „`WithoutCancel` entfernt“ ist an 10 bis 13 äquivalent (§6 *Kontexte*) und wird nicht
  gefahren. Gefahren wird nach dem Umbau die Mutation, die das Herunterfahren trägt:
  `clientRichtung` fragt `Err()` auf dem gelösten statt auf dem übergebenen Kontext; ein
  bestehender Test zum Herunterfahren einer Record-Session muss rot werden. Bleibt er
  grün, ist das ein Befund für den Architect. — **Ausgang:** entfallen: K1 ist vor und
  nach dem Umbau rot (§7), in der Verifikation als V1 nachgefahren; die Signaltests der
  Integration sind grün und unverändert, und „`WithoutCancel` entfernt“ ist an 10 bis 13
  äquivalent (V5 einschließlich beider Integrationsphasen, V6).
- **Verhalten von `replaySitzung` ändert sich unbemerkt** — ein Zweig ist von keinem
  Test erreicht. Je umgebaute Funktion eine Mutation in einem ihrer Zweige, die ein
  bestehender Test vor und nach dem Umbau fängt; fängt keiner, ist das kein neuer Test
  hier. Ändert der Mutant das Verhalten und ist er über die Schnittstelle fangbar,
  nennt §7 ihn mit Test-Idee und Grenze, und die Adresse ist
  `slice-tests-ueberlebende-mutanten` (Entscheidung des Nutzers vom 2026-10-07, dort §1
  *Sammelregel*). `slice-harness-coverage` nimmt ihn nicht an: Sein §1 schließt Tests
  über dem gemessenen Stand aus, und einen Mutanten auf einer abgedeckten Zeile zeigt
  seine Messung nicht (Review F-459 zu `slice-lint-bestand-kern-driven`). Ein Fund im
  PGWire-Adapter oder im Bootstrap ist für den Nehmer eine dritte Schicht, die sein §1
  ausschließt; nach seiner *Sammelregel* schneidet dann der Planner vor dessen Start
  einen zweiten Slice ab, dieser Slice nennt den Fund nur in §7. Erwartet ist ein Fund:
  die Mutation „`WithoutCancel` entfernt“ in `pgwire.Listen` (§6 *Listen*). Test-Idee:
  `bootstrap.Run` mit schon beendetem Kontext und `--listen localhost:0` endet mit
  Exit-Code 0 statt mit `PGR-E4001`. Grenze: nur mit einem Namen als Adresse sichtbar,
  und nur, weil die Namensauflösung von `net` den Kontext liest. Der Implementer fährt
  sie und trägt das Ergebnis in §7 ein. — **Ausgang:** eingetreten:
  `slice-tests-ueberlebende-mutanten-driving` — sieben grüne, verhaltensändernde Mutanten
  (G1 bis G7, §7), alle außer G1 schon vor dem Umbau grün; der Nehmer führt sie in §1
  *Gegenstand* und in drei DoD-Liefer-Punkten.
- **Aufzählung weicht vom Lauf ab** — die Liste in §1 stammt aus der Messung am Stand
  `79f40e1`. Maßgeblich ist der Lauf beim Start (§1); am Stand `a453969` ergab er
  dieselben 16 Befunde, jetzt mit Datei und Zeile in §1. — **Ausgang:** entfallen: Der
  Lauf an `b7d4555` ist sortiert byte-gleich mit §1 und §7 *Vorher* (Verifikation,
  Abschnitt 1 Punkt 3).

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

- **Was hat funktioniert:** Der Schnitt hielt: drei Liefer-Punkte, zwei Schichten
  (Driving-Adapter und Bootstrap), keine Rückführung aus §4, keine Ausnahme in
  `.golangci.yml`, kein `//nolint`, keine geänderte Testdatei. `make lint` sinkt modulweit
  von 16 auf **0** Befunde; damit ist die Vorbedingung von `slice-harness-lint` erfüllt.
  Die Randformen in §6 standen vor dem Code entschieden (`b7d4555`), der Code entschied
  keine weitere (Review, Negativbefund zu Hard Rule 3.12), und beide
  Äquivalenz-Behauptungen hielten der Probe stand (V5 auch in der Integration, V6). Die
  Reihenfolge der Commits aus §3 *Größe* (erst Kontexte und Receiver, dann Komplexität)
  machte den Umbau in einer Review-Sitzung prüfbar. Die Belege in §7 erreichten Review und
  Verifikation; jede nachgemessene Zahl stimmte.
- **Was ging anders als geplant:** Drei Punkte, keiner verlangt Nacharbeit am Verhalten.
  1. Die herausgelösten Funktionen bekamen Doc-Kommentare, die mehr zusagten, als ein
     Test hält: `replayLesefehler` und `replayZustellen` (F-465, MEDIUM), enger gefasst in
     `dd5d13d`; danach `replayWaechter` (V-93), gestrichen in `893b54e`. Dritte
     Wiederholung im Review nach F-401 und F-460.
  2. Die Liste der grünen Mutanten war unvollständig: X4 und Z3 fand das Review (F-466),
     V13 die Verifikation (V-93); nachgetragen als G5 bis G7. Wie in
     `slice-lint-bestand-kern-driven` (F-461) waren alle schon vor dem Umbau grün.
  3. Die erwartete Adresse der grünen Mutanten, `slice-tests-ueberlebende-mutanten`, nahm
     sie nicht an, weil ihr §1 den PGWire-Adapter ausschließt und G1 im Bootstrap liegt
     (F-469); §6 hatte das als Schnitt nach der *Sammelregel* vorgesehen. Der Planner hat
     `slice-tests-ueberlebende-mutanten-driving` in `7cb298c` angelegt. Dazu F-467 (Grenze
     von G2 weiter als der Code, berichtigt in `dd5d13d`) und F-468 (Mutant W nur über die
     Zeitüberschreitung des Pakets rot; an den Nehmer und an `slice-harness-mutation`).
  - **Summary-Zeilen:** Review
    `docs/reviews/2026-10-07-review-slice-lint-bestand-driving.md`: „0 HIGH · 1 MEDIUM · 2
    LOW · 3 INFO (F-465 neue Doc-Kommentare von `replayLesefehler` und `replayZustellen`
    sagen ungeprüft zu; F-466 grüner Mutant `sendFailed` in `replayZustellen` fehlt in der
    Liste; F-467 Grenze von G2/R4 weiter als der Code; F-468 Mutant W nur über
    Paket-Zeitüberschreitung rot; F-469 Übergabe der grünen Mutanten an den Planner nur in
    §7 getragen, nach §3.13 korrekt; F-470 Umbau und Kontexte gleichwertig).
    Wiederkehrende Klassen: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-465,
    drittes Mal nach F-401 und F-460); Liste grüner Mutanten unvollständig (F-466, wie
    F-461).“ Verifikation
    `docs/reviews/2026-10-07-verifikation-slice-lint-bestand-driving.md`, Urteil: Die drei
    Liefer-Punkte sind erfüllt und selbst gemessen; Lint modulweit 16 → 0, ohne
    `//nolint` und ohne neue Ausnahme, Vorbedingung für `slice-harness-lint` erfüllt;
    Komplexität unter den Schwellen, Werte gleich §7; Kontexte in der Form aus §6, das
    Herunterfahren hängt weiter an `ctx.Err()` (V1 rot), die äquivalenten Mutationen sind
    äquivalent; Testliste gleich (162 und 39), `make gates` grün an `dd5d13d`; F-465
    umgesetzt, Rest V-93; G1 bis G6 richtig eingeordnet, ein zweiter Sammel-Slice ist der
    richtige Weg; F-468 an *Laufzeit* von `slice-harness-mutation` und an den neuen
    Sammel-Slice. Kein Befund blockiert die DoD; V-93 LOW (umgesetzt in `893b54e`), V-94
    und V-95 Hinweise an den Planner (im Nehmer eingetragen).
- **Steering-Loop-Eintrag:** Benannte Spec-Lücke, drei Stellen, alle mit Adresse in §6
  von `slice-tests-ueberlebende-mutanten-driving`, der Architect entscheidet sie vor dessen
  Code:
  1. *Lese- und Versandfehler an der Client-Verbindung im Replay.* `LH-FA-02.b` regelt
     Verbindungsende, nicht lesbare Nachricht und Sendefehler nur für `record`. Für
     `replay` sagt keine Stelle, dass ein Verbindungsende des Clients regulär ist, dass
     eine unlesbare Nachricht `PGR-E6001` an den Client ist und dass ein gescheiterter
     Versand `PGR-E4003` ist und die Sitzung beendet; das tun nur Code und Kommentar, und
     die Kommentare mussten darum enger werden (F-465). Sichtbar an R2, R3, R4, X4 und Z3.
     Adresse: Randformen *Verbindungsende im Replay*, *Unlesbare Nachricht im Replay*,
     *Versand im Replay*.
  2. *Signal vor dem Öffnen des Listeners.* Keine Stelle sagt, ob ein Lauf, dessen
     Kontext vor dem Öffnen schon beendet ist, regulär endet oder als Startfehler; sichtbar
     an G1. Adresse: Randform *Signal vor dem Öffnen*.
  3. *Frist eines Tests, der auf einen Kanal wartet.* Keine Stelle in Spezifikation oder
     ADR sagt, ob ein Test eine eigene Frist trägt; ohne sie wird ein roter Lauf erst an
     der Zeitgrenze des Pakets rot (F-468). Adresse: Randform *Warten auf einen Kanal*;
     die Zählung durch ein Gate ist die Randform *Laufzeit* von `slice-harness-mutation`.

  Retirement-Checks: `AGENTS.md` §3.9 (seit welle-walking-skeleton) ist für §1, §3 und §6
  nicht wieder aufgetreten (Review, Negativbefund); die Regel bleibt. §3.10 entfällt, kein
  neuer Vertrag; die Mutationen sind trotzdem je Umbau gefahren. §3.11 (seit
  welle-extended-query) ist wieder aufgetreten (F-465, V-93, F-467); die Regel bleibt, ihr
  geplanter Sensor `slice-harness-mutation` trägt weiter. §3.12 und die
  Randform-Rückgabe ohne Befund im Code-Commit, sie bleiben. §3.13 (seit
  slice-lint-bestand-kern-driven) hat getragen: Der Slice nannte keine Adresse, die nicht
  annimmt, und gab den Punkt an den Planner (F-469, kein MEDIUM); die Regel bleibt.
- **Beobachtungs-Register (`../observations/`):** gesichtet am Stand `7cb298c` (Zähler =
  Dateien unter `evidence/`).
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: **Beleg**, 17× → 18× (F-465,
    V-93, F-467), Stand verkörpert, bleibt; Sensor geplant mit `slice-harness-mutation`.
  - `BEO-REPO/liste-gruener-mutanten-unvollstaendig`: **neu, mit zwei Belegen**,
    `slice-lint-bestand-kern-driven` (F-461, bei dessen Closure als einmalig nicht
    eingetragen) und dieser Slice (F-466, V-93). 2×, Stand offen. Eigener Eintrag statt
    Beleg unter `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: Beide Slices liefern
    keinen Vertrag, die Lücke liegt im Bestand, und was fehlt, ist ihre vollständige
    Übergabe an den Planner.
  - `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit`: **neu, mit zwei Belegen**,
    `slice-harness-blackbox-einstieg` (V-88, bei dessen Closure nicht eingetragen, Adresse
    `slice-harness-integration-wait`) und dieser Slice (F-468). 2×, Stand offen.
  - Ohne Beleg: `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (bleibt 3×, verkörpert;
    F-469 ist nach §3.13 korrekt, keine Adresse ohne Annahme),
    `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (bleibt 3×, verkörpert;
    keine Randform im Code entschieden), `BEO-REPO/randform-wellenlos-ohne-architect-vor-code`
    (bleibt 1×, der Architect entschied §6 vor dem Code in `b7d4555`),
    `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (bleibt 1×, der Umbau verschob
    Lesen und Antworten nicht), `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (bleibt 1×,
    keine Gate-Regel geändert), `BEO-REPO/plan-folgt-korrektur-nicht` (bleibt 16×,
    verkörpert; kein Befund zu §1, §3 oder §6), `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`
    (bleibt 1×, jede Mutation lief in einer frischen Kopie),
    `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (bleibt 13×, kein neuer Vertrag).

  Einmalig und nicht eingetragen: F-469 (Übergabe nur in §7 getragen, nach §3.13
  korrekt), F-470 (keine Aktion), V-94 (Test-Idee unnötig schwer, im Nehmer), V-95
  (Kopplung an `slice-v1-abschluss-betrieb`, im Nehmer). **Kein Eintrag erreicht mit
  diesem Slice neu 3×.** Über der Schwelle stehen nur Einträge mit Ausgang
  (`commit-nennt-struktur-kennung` 3× geplant; `folge-slice-adresse-nimmt-nicht-an` 3×,
  `implementer-bericht-erreicht-pruefer-nicht` 3×,
  `randform-im-code-entschieden-dann-zurueckgegeben` 3×,
  `spec-randform-erst-im-review-entschieden` 11×, `negativtests-fehlen-bei-neuem-vertrag`
  13×, `plan-folgt-korrektur-nicht` 16×, `zusage-im-kommentar-weiter-als-pruefung` 18×,
  verkörpert; `white-box-liste-vor-code-nur-namenssuche` 3× gestrichen).
- **Folge-Slices:** `slice-tests-ueberlebende-mutanten-driving` (G1 bis G7, Risiko
  *Verhalten von `replaySitzung` ändert sich unbemerkt*, dazu die Frist aus F-468 als
  Randform). Er nimmt an: Er liegt in `open/`, §1 *Herkunft* nennt diesen Slice als Geber,
  §1 *Gegenstand* führt alle sieben mit Test-Idee, Mutation und Grenze (V-94 in G5, V-95
  in G2), drei DoD-Liefer-Punkte verlangen die roten Mutationen, und kein Punkt unter
  *Ausdrücklich NICHT* trifft einen der sieben. `slice-harness-mutation` (F-468: Zählung
  eines Mutanten, der nur über die Zeitüberschreitung rot wird); er nimmt an: §1 nennt die
  Sendung mit dieser Kennung, §6 *Laufzeit* führt den Fall, und sein Ausschluss neuer
  Tests trifft die Zählung nicht. `slice-harness-lint` (Anschluss an die Gate-Kette, §1
  hier *Ausdrücklich NICHT*); er nimmt an: Sein §1 nennt diesen Slice, sein §4 *Start*
  verlangt ihn in `done/` und `make lint` ohne Befund, beides erfüllt. Als nächster in der
  Reihe folgt `slice-harness-lint`.
- **Risiken aus §6:** vier, je mit Ausgang: *`contextcheck` nimmt `context.WithoutCancel`
  nicht an* **entfallen**, *Herunterfahren ändert sich unbemerkt* **entfallen**,
  *Verhalten von `replaySitzung` ändert sich unbemerkt* **eingetreten** →
  `slice-tests-ueberlebende-mutanten-driving`, *Aufzählung weicht vom Lauf ab*
  **entfallen**; die Gründe stehen in §6. Die Randformen in §6 sind Entscheidungen, keine
  Risiken.
- **Drei Paarungen:** Anker: Kein Eintrag trägt ein `liegt in`; die Spec-Lücke ist
  benannt, nicht verkörpert. Folge-Slice: `slice-tests-ueberlebende-mutanten-driving` und
  `slice-harness-mutation` liegen in `open/`, `slice-harness-lint` in `next/`; `grep -n
  "slice-lint-bestand-driving"` findet die Kennung in §1 und §2 (DoD) des ersten, in §1
  (*Übernommen aus*) und §6 des zweiten und in §1 und §4 des dritten; ob sie die Sendung
  inhaltlich führen, ist oben beurteilt. Register:
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` trägt
  `evidence/slice-lint-bestand-driving.md`; `BEO-REPO/liste-gruener-mutanten-unvollstaendig`
  und `BEO-REPO/roter-lauf-haengt-bis-zum-zeitlimit` tragen `observation.md`, `state.md`
  und je zwei Dateien in `evidence/`. Die übrigen genannten Einträge bestehen als
  Verzeichnis, jedes mit nicht leerem `evidence/`. Die nächste Welle-Closure prüft erneut.

**Belege des Implementers** (Stand: Code in `78f28a1` (Kontexte und Receiver, Befunde 9
bis 16) und `e76b25f` (Komplexität, Befunde 1 bis 8), beide auf `b7d4555`; Kommentare
nach F-465 und die Einordnung nach F-466 und F-467 in `dd5d13d`, Kommentar von
`replayWaechter` und G7 nach V-93 in dem Commit, der diesen Kopf schreibt; jede Zeile nennt den Stand, an dem sie gemessen ist):

*`make lint`, Zeilen unter `internal/adapters/driving/pgwire`,
`internal/adapters/driving/cli` und `internal/bootstrap`.* Vorher (`b7d4555`, Exit 2,
`16 issues:` modulweit, alle unter diesen Pfaden, dieselben 16 wie in §1):

```text
internal/adapters/driving/pgwire/server.go:302:2: found a struct that contains a context.Context field (containedctx)
internal/adapters/driving/pgwire/server.go:440:47: Non-inherited new context, use function like `context.WithXXX` instead (contextcheck)
internal/adapters/driving/pgwire/server.go:462:34: Non-inherited new context, use function like `context.WithXXX` instead (contextcheck)
internal/adapters/driving/pgwire/server.go:474:40: Non-inherited new context, use function like `context.WithXXX` instead (contextcheck)
internal/bootstrap/bootstrap.go:71:26: Non-inherited new context, use function like `context.WithXXX` instead (contextcheck)
internal/adapters/driving/pgwire/server.go:182:1: calculated cyclomatic complexity for function replaySitzung is 21, max is 15 (cyclop)
internal/adapters/driving/pgwire/server.go:435:1: calculated cyclomatic complexity for function clientRichtung is 17, max is 15 (cyclop)
internal/adapters/driving/pgwire/server.go:569:1: calculated cyclomatic complexity for function startup is 17, max is 15 (cyclop)
internal/adapters/driving/pgwire/server.go:739:1: calculated cyclomatic complexity for function toMessage is 20, max is 15 (cyclop)
internal/adapters/driving/pgwire/server.go:182:1: cognitive complexity 37 of func `(*Server).replaySitzung` is high (> 20) (gocognit)
internal/adapters/driving/pgwire/server.go:435:1: cognitive complexity 22 of func `(*richtungen).clientRichtung` is high (> 20) (gocognit)
internal/adapters/driving/pgwire/server.go:182:1: cyclomatic complexity 20 of func `(*Server).replaySitzung` is high (> 15) (gocyclo)
internal/adapters/driving/pgwire/server.go:739:1: cyclomatic complexity 19 of func `toMessage` is high (> 15) (gocyclo)
internal/adapters/driving/pgwire/server.go:62:22: net.Listen must not be called. use (*net.ListenConfig).Listen (noctx)
internal/adapters/driving/cli/cli.go:243:7: unused-receiver: method receiver 'w' is not referenced in method's body, consider removing or renaming it as _ (revive)
internal/adapters/driving/pgwire/server.go:728:7: unused-receiver: method receiver 's' is not referenced in method's body, consider removing or renaming it as _ (revive)
```

Nach `78f28a1` (Exit 2, `8 issues:`): die acht Komplexitäts-Befunde 1 bis 8, keine
Zeile von `containedctx`, `contextcheck`, `noctx` oder `revive`; `contextcheck` nimmt
`context.WithoutCancel(ctx)` an (Risiko *`contextcheck` nimmt `context.WithoutCancel`
nicht an*: nicht eingetreten). Nachher (`e76b25f`, Exit 0, `0 issues.`, keine
`lint:`-Zeile): **modulweit kein Befund**. Kein `//nolint` im Bestand (`grep` unter
`cmd`, `internal`, `test`: 0 Treffer), `.golangci.yml` unverändert (`git diff b7d4555
e76b25f -- .golangci.yml` leer).

*Komplexität nachher* (golangci-lint v2.14.0 im gepinnten Image, in einer Kopie von
`e76b25f` außerhalb des Repos mit Schwelle 1 je Linter; Schwellen nach `SPEC-049`
Punkt 4: `gocognit` > 20, `gocyclo` > 15, `cyclop` > 15; ein Strich heißt, der Linter
meldet die Funktion auch bei Schwelle 1 nicht):

| Funktion | `gocognit` | `gocyclo` | `cyclop` |
|---|---|---|---|
| `(*Server).replaySitzung` | 17 | 10 | 10 |
| `replayWaechter` | 2 | 3 | 3 |
| `(*Server).replayLesefehler` | — | 2 | 2 |
| `(*Server).replayAntwort` | 7 | 7 | 8 |
| `(*Server).replayZustellen` | — | 2 | 2 |
| `(*richtungen).clientRichtung` | 13 | 10 | 11 |
| `(*richtungen).nachricht` | 5 | 6 | 7 |
| `(*Server).startup` | 12 | 9 | 10 |
| `(*Server).startkopf` | 3 | 8 | 9 |
| `toMessage` | 3 | 10 | 11 |
| `toExtendedMessage` | — | 7 | 8 |
| `rowDescription` | — | 2 | 2 |
| `dataRow` | 3 | 3 | 3 |

Am nächsten an einer Schwelle liegt `replaySitzung` mit `gocognit` 17 (Schwelle 20).

*Testliste* (`go test -list .` je Paket, für `test/integration` mit Build-Tag
`integration`, im Image der Stufe `source` ohne Netz): an `b7d4555`, `78f28a1` und
`e76b25f` je Paket gleich, `postgres` 14, `recording` 14, `cli` 19, `pgwire` 40,
`bootstrap` 12, `model` 10, `services` 53, `test/integration` 39; die sortierten
Namenslisten sind an allen drei Ständen byte-gleich (SHA-256 `40f57928…` über die
Unit-Tests, `6b829c2d…` über die Integrationstests). Keine Testdatei und keine
Abdeckungs-Deklaration geändert (`git diff b7d4555 e76b25f --stat -- '*_test.go' test/
docs/user` leer; `export_test.go` unberührt).

*Läufe:* `make test` grün an `78f28a1` und an `e76b25f` (je vor dem Commit). `make gates`
grün an `e76b25f` (Exit 0, darin `make test`, `make test-integration` mit
`run-integration-tests: gruen`, `make abdeckung-check`, `make kopf-check`, `make a-check`
und `make docs-check`); ein weiterer Lauf an dem Commit, der diese Belege schreibt.

*Mutationstabelle* (§3.10 von `AGENTS.md`, Schritt 19). Weg: je Mutant ein frischer Pfad
außerhalb des Repos (Kopie per `git archive` des genannten Stands, Ersetzung genau einer
Textstelle mit Treffer-Prüfung, danach `touch` auf die Datei), Lauf über
`docker build --target test` mit eigenem Tag; Integrationsläufe mit eigenem Image, eigenem
internen Netz, eigenem PostgreSQL-Container und eigenem Volume, beide Phasen wie
`make test-integration`. Image, Netz, Container, Volume und Kopie danach entfernt. „vor“
ist für die Kontexte `b7d4555`, für die Komplexität `78f28a1`; „nach“ ist `e76b25f`.

| # | Zusage (bewahrt) | Mutation | vor | nach | roter Test |
|---|---|---|---|---|---|
| K1 | Das Herunterfahren einer Record-Session hängt an `ctx.Err()` in `clientRichtung` (Pflichtmutation, §6 Risiko *Herunterfahren ändert sich unbemerkt*) | `Err()` auf dem gelösten Kontext (`uc.Err()`, vor: `r.ctx.Err()`) | rot | rot | `TestExtendedHerunterfahren` (4 Untertests), `TestHerunterfahrenWeckenNichtVerloren` |
| K2 | `record` schreibt die Aufzeichnung am Ende über `Finish` | Aufruf `service.Finish(…)` durch `error(nil)` ersetzt | rot | rot | `TestRunRecordSchreibfehlerJeStufe` |
| K3 | Eine nicht nutzbare Adresse ist `PGR-E4001` (`Listen`) | Code `CodeListen` → `CodeNetwork` | Unit grün, Integration rot | Unit grün, Integration rot | `TestE2EReplayNichtVerbrauchtStartfehler` |
| K4 | Ein beendeter `ctx` bricht das Öffnen in `Listen` nicht ab | `context.WithoutCancel` entfernt | — | **grün** (Unit und Integration) | keiner, eingeordnet unten (G1) |
| K5 | `wahrheitswert` bleibt eine Option ohne Wert (`IsBoolFlag`) | `return false` | rot | rot | `TestParseFailOnUnconsumed`, `TestParseFailOnUnconsumedUmgebung`, `TestParseFailOnUnconsumedUmgebungNebenOption`, `TestRunLogLevel` |
| K6 | `send` bleibt Methode von `*Server` | Methode umbenannt | — | rot (Build) | `go vet`: `s.send undefined` |
| R1 | Eine Extended-Nachricht ohne Antwort beendet die Replay-Sitzung nicht | `continue` (nach: `return true`) → Ende | rot | rot | `TestReplayExtended`, `TestReplaySent`, `TestReplayHerunterfahren` (2), `TestReplayHerunterfahrenSpaeteFrist`, `TestInfoOhneZeileJeVerbindung` |
| W | Der Wächter der Replay-Sitzung weckt das Lesen bei Ende von `ctx` | `SetReadDeadline(time.Now())` entfernt | rot | rot | `TestReplayHerunterfahrenSpaeteFrist` (Zeitüberschreitung des Pakets nach 10 min, der Lauf nennt diesen Test) |
| R2 | Ein Verbindungsende im Replay ist regulär, ohne Fehlerantwort | Bedingung `verbindungsende(err) && false` | **grün** | **grün** (Unit und Integration) | keiner, eingeordnet unten (G2) |
| R3 | Eine unlesbare Nachricht im Replay ist `PGR-E6001` an den Client | `s.fail(…)` entfernt | **grün** | **grün** (Unit und Integration) | keiner, eingeordnet unten (G3) |
| R4 | Nach einem Lesefehler endet die Replay-Sitzung | `return` → `continue` | **grün** | **grün** (Unit und Integration) | keiner, eingeordnet unten (G2) |
| C1 | Terminate beendet die Record-Session mit `EndTerminate` | `EndTerminate` → `EndClosed` | rot | rot | `TestQueryUndTerminate`, `TestEndeSchliesstVerbindung`, `TestExtendedEreignisse/Terminate`, `TestExtendedSyncGruppe`, `TestInfoOhneZeileJeVerbindung` |
| C2 | Nach einer zugestellten Antwort liest die Client-Richtung weiter | Ende nach `schreibe` | rot | rot | `TestQueryUndTerminate`, `TestExtendedHerunterfahren/gelesene_Anfrage`, `TestInfoOhneZeileJeVerbindung` |
| S1 | `startup` prüft die Länge vor dem Startcode | Länge nur bei bekanntem Startcode geprüft | rot | rot | `TestFremdeErsteNachricht` |
| T1 | `NoData` wird `NoData` | → `PortalSuspended` | rot | rot | `TestToMessageExtended`, `TestExtendedFlush` |
| T2 | Ein Wert NULL bleibt in der DataRow `nil` | Bedingung `true` | rot | rot | `TestQueryUndTerminate` |
| T4 | Die Spaltenbeschreibung trägt den Namen | `Name` → `nil` | rot | rot | `TestQueryUndTerminate` |
| T6 | … den Typ-OID | `DataTypeOID` → 0 | Unit grün, Integration rot | Unit grün, Integration rot | `TestE2EReplaySelect1`, `TestE2EErgebnisartenEinfach`, `TestE2EFehlerreplayEinfach` |
| T7 | … die Typgröße | `DataTypeSize` → 0 | Unit grün, Integration rot | Unit grün, Integration rot | `TestE2EErgebnisartenEinfach`, `TestE2EFehlerreplayEinfach` |
| T8 | … den Typmodifikator | `TypeModifier` → 0 | Unit grün, Integration rot | Unit grün, Integration rot | `TestE2EErgebnisartenEinfach`, `TestE2EFehlerreplayEinfach` |
| T9 | … das Format | `Format` → 0 | Unit grün, Integration rot | Unit grün, Integration rot | `TestE2EErgebnisartenEinfach` |
| T3 | … die Tabellen-OID | `TableOID` → 0 | **grün** | **grün** (Unit und Integration) | keiner, eingeordnet unten (G4) |
| T5 | … die Spaltennummer | `TableAttributeNumber` → 0 | **grün** | **grün** (Unit und Integration) | keiner, eingeordnet unten (G4) |
| X4 | Ein gescheiterter Versand im Replay wird gemerkt (Review F-466) | `s.sendFailed(err)` entfernt (vor: an beiden Stellen) | **grün** (Unit) | **grün** (Unit und Integration) | keiner, eingeordnet unten (G5) |
| Z3 | Nach einem gescheiterten Versand endet die Replay-Sitzung | `return false` → `return true` (vor: `return` → `continue`, beide Stellen) | **grün** (Unit) | **grün** (Unit und Integration) | keiner, eingeordnet unten (G6) |
| V13 | Ein Schließen von `fertig` beendet den Wächter der Replay-Sitzung (Verifikation V-93) | `case <-fertig:` entfernt | **grün** (Unit) | **grün** (Unit) | keiner, eingeordnet unten (G7) |

*Kommentare nach F-465* (Satz für Satz von einer roten Mutation gedeckt; gefahren an
der Fassung dieses Commits, als `git stash create` vor dem Commit, je Mutant ein
frischer Pfad wie oben; ohne Mutant `make`-Stufe `test` grün):

| Kommentar | Satz | Mutation | roter Test |
|---|---|---|---|
| `replayZustellen` | „schreibt Antworten an den Client“ | `s.send(…)` durch `error(nil)` ersetzt (Z1) | `TestReplaySent`, `TestReplayExtended`, `TestReplayHerunterfahren` (2), `TestReplayHerunterfahrenSpaeteFrist`, `TestInfoOhneZeileJeVerbindung` |
| `replayZustellen` | „meldet sie dem Use Case … als gesendet“ | `s.replayer.Sent(…)` entfernt (X8) | `TestReplaySent` |
| `replayZustellen` | „danach …, nicht, wenn das Schreiben scheitert“ | `Sent` vor `send` (Z2) | `TestReplaySent` |
| `replayLesefehler` | „behandelt einen Lesefehler der Replay-Sitzung“ | — | nennt nur den Gegenstand, sagt kein Verhalten zu |

Entfallen sind „Ein Ende der Client-Verbindung ist im Replay regulär …; was unverbraucht
bleibt, meldet closeReplay“ und „Jeder andere Lesefehler geht als PGR-E6001 an den
Client“ (`replayLesefehler`, R2, R3, R4 grün) sowie „scheitert das Schreiben, merkt es den
Fehler und meldet, dass die Sitzung endet“ (`replayZustellen`, X4, Z3 grün). Sie stehen
als Test-Ideen in G2, G3, G5 und G6.

Je Umbau ist mindestens eine Mutation rot, vor und nach dem Umbau aus demselben Test:
`replaySitzung` (R1, W), `clientRichtung` (C1, C2, K1), `startup` (S1), `toMessage`
(T1, T2, T4, T6 bis T9), die Kontexte (K1, K2, K3), die Receiver (K5, K6). Eine
Mutation an `bootstrap.go:71` und an `contextcheck` 10 bis 12 („`WithoutCancel`
entfernt“) ist nicht gefahren: äquivalent nach §6 *Kontexte*. Der Rückgabewert von
`weiterlesen` bleibt verworfen; an ihm ist nichts gefahren (§6, entschieden).

*Grüne Mutanten, eingeordnet* (Schritt 19). Keiner ist äquivalent; jeder ändert das
Verhalten und ist über die Schnittstelle fangbar. §1 schließt neue Tests aus, darum je
Mutant eine Test-Idee mit Grenze. Alle außer G1 sind schon am Stand vor dem Umbau grün
(`78f28a1`; X4 und Z3 dort an beiden Stellen in `replaySitzung`),
sind also Lücken des Bestands, nicht des Umbaus. **Nehmer:**
`slice-tests-ueberlebende-mutanten-driving` (angelegt vom Planner in `7cb298c` nach der
*Sammelregel* von `slice-tests-ueberlebende-mutanten`, Entscheidung des Nutzers vom
2026-10-07); er führt G1 bis G7 in §1 *Gegenstand* mit dieser Kennung als Geber, G5 mit
der Test-Idee aus V-94, G2 mit der Grenze aus V-95.

- **G1 — `Listen` ohne `WithoutCancel` (K4)**, der in §6 erwartete Fund. Test-Idee nach
  §6: `bootstrap.Run` mit schon beendetem Kontext und `record --listen localhost:0` endet
  mit Exit-Code 0 statt mit `PGR-E4001`. Gefahren in einer Kopie von `e76b25f` als
  temporärer Test (nicht im Repo), 50 Läufe je Stand: ohne Mutant 50 grün, mit Mutant 50
  rot mit Exit-Code 4 und `Netzwerk [PGR-E4001]: Adresse localhost:0 nicht nutzbar:
  listen tcp: lookup localhost: operation was canceled`. Grenze: nur mit einem Namen als
  Adresse sichtbar, weil allein die Namensauflösung von `net` den Kontext liest; mit
  einer IP-Adresse, wie sie alle bestehenden Tests nutzen, bleibt der Mutant grün.
- **G2 — Ende der Replay-Sitzung nach Verbindungsende (R2, R4).** Nach dem Code schreibt R2 beim
  Verbindungsende eine Fehlerantwort und merkt `PGR-E6001` als ersten Verbindungsfehler
  (Exit-Code ungleich 0); R4 liest nach dem Verbindungsende erneut, und solange `ctx`
  läuft, endet die Sitzung nicht und `closeReplay` läuft nicht. Endet `ctx` und gibt
  `Shutdown` das Ende frei, kehrt `replaySitzung` auch unter R4 zurück, und das per
  `defer` registrierte `closeReplay` läuft (Review F-467). Test-Idee: Replay-Verbindung über
  `pgwire.Handle` mit Replay-Fake, nach dem Startup schließt der Client ohne Terminate;
  `Handle` kehrt binnen einer Frist zurück (fängt R4), `CloseConnection` ist einmal
  gerufen und `FirstErrorCode()` bleibt leer (fängt R2). Grenze: Die Fehlerantwort an den
  geschlossenen Client ist nicht lesbar, beobachtbar ist nur der gemerkte Code; das Ende
  der Sitzung nur über eine Frist, und nur, solange `ctx` nicht endet.
- **G3 — unlesbare Nachricht im Replay (R3).** Der Mutant beendet die Sitzung ohne
  Fehlerantwort und ohne gemerkten Code. Test-Idee: Replay-Verbindung, nach dem Startup
  ein Nachrichtenkopf mit einem Typ, den `pgproto3` nicht kennt; der Client erhält eine
  ErrorResponse mit SQLSTATE `0A000` und `PGR-E6001`, `FirstErrorCode()` ist `PGR-E6001`.
  Grenze: Welcher Lesefehler eine unlesbare Nachricht ist, legt `pgproto3` fest; der Test
  hält nur die Art fest, die er sendet.
- **G4 — Tabellen-OID und Spaltennummer der Spaltenbeschreibung (T3, T5).** Kein
  Unit- und kein Integrationstest fängt die beiden Felder. Test-Idee: `pgwire.ToMessage` mit einer RowDescription, deren
  Spalte in jedem Feld einen Wert ungleich 0 trägt, vergleicht jedes Feld der
  `FieldDescription`. Grenze: Die Felder sind ein Durchreichen ohne Logik; der Test
  fängt das Vertauschen oder Weglassen eines Feldes, nicht dessen Bedeutung.
- **G5 — Merken eines gescheiterten Versands im Replay (X4, Review F-466).** Ohne
  `sendFailed` wird ein Schreibfehler an den Client nicht als `PGR-E4003` gemerkt, und
  der Lauf endet mit Exit-Code 0. Test-Idee: Replay über `pgwire.Handle` mit einer
  Verbindung, deren Schreiben nach dem Startup scheitert (Wrapper um `net.Conn`), dann
  eine Anfrage; `FirstErrorCode()` ist `PGR-E4003`. Für eine nicht abbildbare Antwort
  des Replay-Fakes ist `FirstErrorCode()` der Code des internen Fehlers. Grenze: Den
  Schreibfehler muss der Test erzwingen; ein geschlossener TCP-Peer lässt den ersten
  Schreibvorgang oft gelingen.
- **G6 — Ende der Replay-Sitzung nach gescheitertem Versand (Z3).** Nach einem
  Schreibfehler scheitert in der Regel auch das nächste Lesen an derselben Verbindung,
  und die Sitzung endet über den Lesefehler; dort ist Z3 nicht beobachtbar. Sichtbar ist
  er bei einer nicht abbildbaren Antwort: `send` scheitert, bevor es schreibt, und unter
  Z3 liest die Sitzung weiter und beantwortet die nächste Anfrage. Test-Idee: Replay-Fake
  liefert auf die erste Anfrage einen Antworttyp ohne Abbildung; `Handle` kehrt zurück,
  ohne eine zweite Anfrage zu lesen, und `Query` ist einmal gerufen. Grenze: nur über
  den Pfad der nicht abbildbaren Antwort fangbar.
- **G7 — Ende des Wächters der Replay-Sitzung (V13, Verifikation V-93).** Ohne
  `case <-fertig:` läuft der Wächter nach dem Ende der Sitzung weiter, bis `ctx` endet,
  je Replay-Verbindung eine Goroutine; danach setzt er die Lesefrist einer geschlossenen
  Verbindung. Gefahren an der Fassung dieses Commits (`git stash create`) und an
  `78f28a1`, je ein frischer Pfad: Stufe `test` grün. Der Satz „Ein Schließen von fertig
  beendet ihn.“ ist darum aus dem Kommentar von `replayWaechter` gestrichen; die übrigen
  Sätze deckt W. Test-Idee: Replay-Verbindung über `pgwire.Handle`, Client endet mit
  Terminate, `ctx` läuft weiter; nach der Rückkehr von `Handle` kehrt die Zahl der
  Goroutinen (`runtime.NumGoroutine`) binnen einer Frist auf den Stand vor der Verbindung
  zurück. Grenze: über das Protokoll nicht beobachtbar, nur über die Zahl der Goroutinen
  des Prozesses; der Test darf nicht parallel zu anderen laufen und braucht eine Frist.

*Randformen:* Keine neue gefunden. Keine ungeprüfte Reihenfolge zweier Fehler: Die einzige
Reihenfolge, Länge vor Startcode in `startup`, hält `TestFremdeErsteNachricht` (S1 rot).
Kein Umbau ändert ein Verhalten; §6 ist in keinem Code-Commit geändert.

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
- `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (2×) — die Semantik des
  Herunterfahrens ist genau eine Randform, die eine Kontext-Bereinigung still
  entscheiden würde; §6 nennt sie, der Architect bestätigt vor dem Code. Ein drittes
  Auftreten in diesem Slice erreichte die Schwelle.
- `BEO-REPO/replay-haengt-am-zeitverhalten-des-clients` (1×) — `replaySitzung` und
  `clientRichtung` tragen das Zeitverhalten des Replays; ein Umbau darf die Reihenfolge
  von Lesen und Antworten nicht verschieben (§6).
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×) — verwandt mit dem Risiko, dass
  ein Umbau eine Prüfung verliert (§6).
- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×) — wellenlos; §4 nennt den
  Architect vor dem Code.

Keiner der Einträge erreicht mit diesem Slice allein die Schwelle 3×; keine neue
Lücke vor dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
