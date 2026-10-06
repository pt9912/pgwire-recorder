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

**Verantwortlich:** —

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

- **Komplexität (8) in vier Funktionen** in `internal/adapters/driving/pgwire/server.go`
  — `gocognit`, `gocyclo` bzw. `cyclop` über der Schwelle aus Punkt 4:
  `(*Server).replaySitzung` (alle drei, `gocognit` 37, Schwelle 20),
  `(*richtungen).clientRichtung` (`gocognit`, `cyclop`), `(*Server).startup`
  (`cyclop`), `toMessage` (`gocyclo`, `cyclop`).
- **Kontexte (6)** — `containedctx` für den Kontext im Struct `richtungen`;
  `contextcheck` an drei Stellen in `pgwire` und an einer in `internal/bootstrap`
  (neuer Kontext statt des übergebenen, etwa `service.Finish(context.Background())`);
  `noctx` für `net.Listen` in `pgwire.Listen`.
- **`revive` (2)** — `unused-receiver` je einmal in `pgwire` und in `cli`.

Maßgeblich ist die Liste, die `make lint` beim Start unter diesen Pfaden ausgibt; weicht
sie von der Aufzählung ab, gilt der Lauf, und der Plan folgt ihm (`AGENTS.md` §3.9).
`cmd/` hat in der Messung keinen Befund.

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
  ein Umbau, gemessen an denselben Tests (§6).
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

- [ ] Komplexität: Die vier Funktionen aus §1 liegen unter den Schwellen aus
      `SPEC-049` Punkt 4, ohne Änderung des Verhaltens: Die Liste der Tests
      (`go test -list .` je Paket, im gepinnten Go-Image) ist vor und nach dem Umbau
      gleich, keine Erwartung ist geändert, `make test` und `make test-integration` sind
      grün.
- [ ] Kontexte: Die sechs Befunde von `containedctx`, `contextcheck` und `noctx` sind
      behoben, ohne die Semantik des Herunterfahrens zu ändern (§6); die Tests und
      Integrationstests zum Abbruchsignal und zur Frist des Herunterfahrens laufen
      unverändert grün.
- [ ] `make lint` meldet unter `internal/adapters/driving/pgwire`,
      `internal/adapters/driving/cli` und `internal/bootstrap` keinen Befund, ohne neue
      Ausnahme in `.golangci.yml` und ohne `//nolint`; über das ganze Modul meldet es
      keinen Befund mehr. Die Zeilen der Ausgabe unter diesen Pfaden vor und nach dem
      Slice stehen im Bericht.
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
| `internal/adapters/driving/pgwire/server.go` | refactor | `replaySitzung`, `clientRichtung`, `startup` und `toMessage` unter die Schwellen, durch Herauslösen von Schritten in unexportierte Funktionen; Kontext nicht mehr im Struct `richtungen`, sondern als Parameter; die drei neuen Kontexte nach §6; ungenutzter Receiver entfernt |
| `internal/adapters/driving/pgwire/server.go` (`Listen`) | update | `net.ListenConfig.Listen` mit Kontext statt `net.Listen`; die Signatur von `pgwire.Listen` bekommt den Kontext, der Meldungscode bleibt `PGR-E4001` |
| `internal/bootstrap/bootstrap.go` | update | Aufrufer von `pgwire.Listen` reichen ihren Kontext; der neue Kontext an `service.Finish` nach §6 |
| `internal/adapters/driving/cli/*.go` | update | ungenutzter Receiver |
| Testdateien dieser Pakete und `test/integration` | unverändert | Beleg des unveränderten Verhaltens; ein Test, der wegen des Umbaus geändert werden müsste, ist ein Befund (§4). Ausgenommen sind Aufrufe von `pgwire.Listen` (neuer Parameter) und `export_test.go`, falls eine dort weitergereichte unexportierte Funktion umbenannt wird; dann nur der Aufruf bzw. Verweis |

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
- **Kontexte** — die erwartete Bereinigung eines neuen Kontexts ist
  `context.WithoutCancel(ctx)` auf dem übergebenen Kontext
  ([ADR-0034](../../adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 5): Er
  behält die Werte des Aufrufers und überlebt dessen Abbruch, wie ein neuer
  Hintergrund-Kontext heute. Ändert die Bereinigung die Semantik des Herunterfahrens
  (Signalbehandlung nach [`LH-FA-13.a`](../../../../spec/spezifikation.md#lh-fa-13a--signalbehandlung), Frist `--shutdown-timeout`), geht der Befund an den
  Architect zurück. Der Kontext im Struct `richtungen` wird Parameter der Methoden, die
  ihn brauchen; `pgwire.Listen` bekommt einen Kontext für `net.ListenConfig`.

**Risiken:**

- **`contextcheck` nimmt `context.WithoutCancel` nicht an** — der Code nutzt die Form
  schon an vier Stellen (`handle`, `recordSitzung`, `closeReplay`, `closeRecord` in
  `server.go`), und die Messung meldet trotzdem drei Befunde in `pgwire`; ob die
  Befunde an genau diesen Stellen stehen, zeigt erst der Lauf. Dann trägt die erwartete
  Bereinigung nicht, und der Weg ist eine Entscheidung des Architect (§4, `→ open`).
  — **Ausgang:** — (bei Closure)
- **Herunterfahren ändert sich unbemerkt** — ein Kontext, der bisher den Abbruch des
  Aufrufers überlebte, bricht nach der Bereinigung mit ihm ab (oder umgekehrt); das
  Recording wird dann beim Signal nicht mehr oder anders geschrieben. Die
  Integrationstests zum Abbruchsignal laufen vor und nach dem Umbau; dazu je
  bereinigter Stelle eine Mutation (`WithoutCancel` entfernt), die ein bestehender Test
  fängt. Fängt keiner, ist das ein Befund für den Architect. — **Ausgang:** — (bei
  Closure)
- **Verhalten von `replaySitzung` ändert sich unbemerkt** — ein Zweig ist von keinem
  Test erreicht. Je umgebaute Funktion eine Mutation in einem ihrer Zweige, die ein
  bestehender Test vor und nach dem Umbau fängt; fängt keiner, ist das ein Befund für
  `slice-harness-coverage`, kein neuer Test hier. — **Ausgang:** — (bei Closure)
- **Aufzählung weicht vom Lauf ab** — die Liste in §1 stammt aus der Messung am Stand
  `79f40e1`. Maßgeblich ist der Lauf beim Start (§1). — **Ausgang:** — (bei Closure)

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
