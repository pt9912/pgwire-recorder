# ADR-0034: Lint-Gate mit SOLID-nahem Profil, gestuft eingeführt

**Status:** Proposed

**Datum:** 2026-10-06

**Autor:** pt9912

**Bezug:** [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md), [ADR-0001](0001-hexagonale-architektur.md)

**Schärft:** [`SPEC-049`](../../../spec/spezifikation.md#spec-049--lint-profil-lint)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`LH-QA-07` verlangt in Messmethode (1) eine statische Analyse nach festgelegtem Profil ohne Befund, mit Ausnahmen zentral und begründet, und in Messmethode (3) Unit-Tests außerhalb der geprüften Einheit, die nur über eine für Tests vorgesehene Stelle auf Internes zugreifen. Ein Lint-Gate gibt es nicht; `AGENTS.md` §3.2 ist ein Platzhalter. Vorbild ist das Profil des Schwester-Repos ai-harness-init mit golangci-lint v2. Der Nutzer hat am 2026-10-05 entschieden: die Schwellen des Vorbilds gelten unverändert, `testpackage` wird gestuft eingeführt und von vier Umstellungs-Plänen je Paketgruppe scharf geschaltet, der Vertrag des Gates steht in der Spezifikation.

**Messung des Bestands** (Stand `79f40e1`, Profil des Vorbilds unverändert mit Build-Tag `integration` und ungekürzter Ausgabe, golangci-lint `v2.14.0` im Image des Vorbilds, ohne Netz). Befunde je Linter und Paket; `cmd/pgwire-recorder` hat keinen.

| Linter | bootstrap | postgres | recording | cli | pgwire | model | services | integration | Σ |
|---|---|---|---|---|---|---|---|---|---|
| `errcheck` | 4 | 18 | 0 | 0 | 18 | 0 | 0 | 17 | 57 |
| `gochecknoglobals` | 1 | 4 | 2 | 1 | 1 | 4 | 18 | 2 | 33 |
| `staticcheck` | 0 | 0 | 9 | 1 | 0 | 3 | 0 | 1 | 14 |
| `revive` | 0 | 1 | 0 | 1 | 3 | 7 | 0 | 1 | 13 |
| `contextcheck` | 1 | 0 | 0 | 0 | 3 | 0 | 5 | 0 | 9 |
| `gocognit` | 0 | 0 | 1 | 0 | 2 | 1 | 1 | 0 | 5 |
| `ireturn` | 0 | 3 | 0 | 0 | 0 | 0 | 2 | 0 | 5 |
| `gocyclo` | 0 | 1 | 0 | 0 | 1 | 0 | 0 | 0 | 2 |
| `containedctx`, `cyclop`, `noctx`, `unused` | 0 | 0 | 1 | 0 | 3 | 0 | 0 | 0 | 4 |
| **Σ ohne `testpackage`** | 6 | 27 | 13 | 3 | 31 | 15 | 26 | 21 | **142** |
| `testpackage` (je Testdatei) | 1 | 2 | 1 | 1 | 3 | 2 | 7 | 7 | 24 |
| **Σ mit `testpackage`** | 7 | 29 | 14 | 4 | 34 | 17 | 33 | 28 | **166** |

Ohne `testpackage` stehen 57 Befunde im Produkt-Code und 85 in Testdateien. Die übrigen 15 Linter des Profils, darunter `funlen`, `nestif`, `maintidx`, `dupl`, `interfacebloat` und `gomodguard_v2`, melden nichts. Die Testdatei-Ausnahmen des Vorbilds blenden weitere 39 Befunde aus (`noctx` 15, `gocognit` 7, `revive` 12, `gocyclo` 4, `unparam` 1).

Nach den dauerhaften Entscheidungen unten (`ST1005` 12, `ireturn` 5, `errcheck` auf Netz-`Close` 43, Nachschlage-Tabellen 10) bleiben 96: 25 im Produkt-Code, 71 in Testdateien, davon 24 `testpackage`. Die 25 im Produkt-Code verteilen sich auf alle vier Schichten (Domain Model, Services, Driven und Driving Adapter, dazu Bootstrap); acht davon sind Komplexitätsbefunde in sieben Funktionen, die schwerste `(*Server).replaySitzung` mit `gocognit` 37.

## Entscheidung

Wir führen `make lint` als Stufe `lint` des `Dockerfile` mit dem Profil des Vorbilds und seinen Schwellen ein, an dieses Repo angepasst; die Einzelregeln stehen in [`SPEC-049`](../../../spec/spezifikation.md#spec-049--lint-profil-lint). Im Einzelnen entschieden ist:

1. **Durchsetzung statt Vertrauen.** golangci-lint wertet `//nolint` aus, statt es zu verbieten; eine eigene Prüfung in derselben Stufe meldet jede Direktive. Zwei weitere eigene Prüfungen halten die Form der Ausnahmen (`# Why:` über jeder Regel) und der Brücke (kein Test in `export_test.go`). `exclusions.warn-unused` wird zum Befund: Eine Regel, die nichts ausblendet, ist falsch gelesen oder überholt. In der Messung blendete eine verankerte Pfad-Regel nichts aus, weil golangci-lint Pfade relativ zur Konfiguration liest; die Warnung zeigte es, darum gilt `relative-path-mode: cfg` ausdrücklich.
2. **Dauerhafte Ausnahmen mit Grund, der auch für neuen Code gilt.** Die Testdatei-Ausnahmen des Vorbilds; `ST1005`, weil deutsche Fehlertexte mit einem Substantiv beginnen; `ireturn` für Ports und `pgproto3`-Nachrichten, weil ein Adapter im Hexagon seinen Port liefert ([ADR-0001](0001-hexagonale-architektur.md)) und `pgproto3` Nachrichten als Interface führt ([ADR-0010](0010-verwendung-von-pgproto3.md)); `errcheck` für `Close` an Netz-Verbindung und Listener, weil ihr Fehler nach dem Ende nichts mehr ändert; `gochecknoglobals` je Datei und Name für zehn Nachschlage-Tabellen und Sentinel-Werte, die nur gelesen werden — Go kennt keine konstante Map. `main.version` braucht keine Ausnahme; der Linter lässt den Namen `version` selbst zu.
3. **`gomodguard_v2` als Erlaubnisliste** der direkten Module aus `go.mod` statt der Sperrliste des Vorbilds. Ein neues Modul ändert die Liste im selben Commit wie `go.mod` und wird damit sichtbar; welches Paket welche Bibliothek nutzt, hält weiter `make a-check`.
4. **Export-Test-Brücke zulässig**, als einzige Datei `export_test.go` je Paket, nur mit Aliasen, Konstanten und weiterreichenden Funktionen, Zustand nur an übergebenen Werten. `testpackage` überspringt nur diesen Dateinamen; sein Default ließe auch `internal_test.go` als White-Box-Test durch. Rot-Fälle für die zweite Bedingung von Messmethode (3): ein Test in `internal_test.go` im Paket des Codes (an der Brücke vorbei, `testpackage`), eine Funktion `Test…` in `export_test.go` (eigene Prüfung), eine Variable in `export_test.go` (`gochecknoglobals`).
5. **Gestufte Einführung, kein Wert an den Bestand angepasst** (Empfehlung; die Wahl trifft der Nutzer, darum `Proposed`). Das Gate ist ab dem ersten Lauf für jeden neuen Befund scharf. Die 96 verbleibenden Befunde stehen als Stufen mit Hochschalt-Trigger in `.golangci.yml`, je so eng wie der Befund ([`SPEC-049`](../../../spec/spezifikation.md#spec-049--lint-profil-lint) Punkt 8):

   | Stufe | Pfade | Linter | Befunde | aufgehoben von |
   |---|---|---|---|---|
   | Tests Kern | Testdateien unter `internal/hexagon/model`, `internal/hexagon/services` | `testpackage`, `contextcheck`, `gochecknoglobals` | 31 | Umstellungs-Plan Kern |
   | Tests Driven | Testdateien unter `internal/adapters/driven/postgres`, `internal/adapters/driven/recording` | `testpackage`, `errcheck`, `gochecknoglobals`, `unused` | 12 | Umstellungs-Plan Driven |
   | Tests PGWire | Testdateien unter `internal/adapters/driving/pgwire` | `testpackage`, `gochecknoglobals`, `revive` | 6 | Umstellungs-Plan PGWire |
   | Tests Einstieg | Testdateien unter `internal/adapters/driving/cli`, `internal/bootstrap`, `test/integration` | `testpackage`, `errcheck`, `gochecknoglobals`, `revive`, `staticcheck` | 22 | Umstellungs-Plan Einstieg |
   | Code Kern und Driven | `Group.validate`, `cursor.objekte`, `toResponse`, `fromDTO` (Komplexität); Doc-Kommentare in `model`; ungenutzter Parameter in `postgres`; `QF1001` | `gocognit`, `gocyclo`, `revive`, `staticcheck` | 13 | Bereinigungs-Plan Kern und Driven |
   | Code Driving | `pgwire/server.go` (`replaySitzung`, `clientRichtung`, `startup`, `toMessage`, Kontext im Struct, drei neue Kontexte, `net.Listen`, ungenutzter Receiver); `cli` ungenutzter Receiver; `bootstrap` neuer Kontext | `gocognit`, `cyclop`, `gocyclo`, `containedctx`, `contextcheck`, `noctx`, `revive` | 12 | Bereinigungs-Plan Driving |

   Die Test-Stufen heben die vier Umstellungs-Pläne auf, die jede dieser Testdateien ohnehin umschreiben; ihr Umfang wächst um die Befunde außer `testpackage`. Die Code-Stufen folgen den Umstellungen, weil Komplexitäts-Umbauten an unexportierten Funktionen erst dann keine White-Box-Tests mehr brechen. Die erwartete Bereinigung der neuen Kontexte ist `context.WithoutCancel`; ändert sie die Semantik des Herunterfahrens, geht der Befund an den Architect zurück.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Bestand vor dem Gate bereinigen (Plan zurück nach `next`) | Gate startet ohne Stufe | 96 Befunde in allen Schichten; die Testbefunde kollidieren mit den Umstellungs-Plänen, die dieselben Dateien danach umschreiben; neuer Code bleibt bis dahin ungeprüft |
| B — alles in diesem Plan bereinigen | eine Lieferung | mehr als zwei Schichten, sieben Komplexitäts-Umbauten, nicht in einer Review-Sitzung prüfbar |
| C — Schwellen oder Linter an den Bestand anpassen | sofort grün | Lockerung nach `AGENTS.md` §3.6; gegen die Entscheidung des Nutzers |
| D — `testpackage` ganz aus bis zum letzten Umstellungs-Plan | eine Regel weniger | keine Zwischenstufe belegt ihr Hochschalten am Gate |
| E — `errcheck`-Preset `std-error-handling` | weniger Konfiguration | blendet auch `Close` von Dateien aus, deren Fehler Daten kostet |
| F — `//nolint` mit `nolintlint` begrenzen | kein eigener Schritt | `nolintlint` verlangt Form, verbietet nicht |
| **G — Gate jetzt, dauerhafte Ausnahmen mit Grund, Rest als enge Stufen mit Hochschalt-Trigger** | neuer Code ab sofort geprüft; jede Stufe hat Adresse; ungenutzte Regeln fallen auf | `.golangci.yml` trägt bis zur Bereinigung rund zwanzig Stufen-Regeln |

## Konsequenzen

- Positiv: Jeder neue Befund und jedes `//nolint` macht das Gate rot; eine Stufe, deren Befund bereinigt ist, muss gelöscht werden, sonst ist das Gate rot.
- Negativ (akzeptiert): Bis die Stufen aufgehoben sind, gilt Messmethode (1) mit Ausnahmen, die der Bestand trägt; ob ein `Why:` zutrifft und ob der Plan einer Stufe existiert, prüft Review. Testdateien im Paket `main` lässt `testpackage` zu; `cmd/` hat keine. Ein `//nolint` in einem String-Literal ist ein Befund.
- Folgepflicht: Eine Gegenprobe hält jede Zusage aus [`SPEC-049`](../../../spec/spezifikation.md#spec-049--lint-profil-lint), die eine Mutation fangen kann, darunter die drei Rot-Fälle der Brücke und eine ungenutzte Regel. Mit Option G entstehen zwei Bereinigungs-Pläne in `open/`, bevor eine Code-Stufe sie im `Why:` nennt.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| golangci-lint und eigene Prüfungen | Profil nach [`SPEC-049`](../../../spec/spezifikation.md#spec-049--lint-profil-lint) ohne Befund | Gate-Ziel, genannt in `harness/README.md` §Sensors |
| Gegenprobe | je Zusage ein Mutant, der rot wird | Gegenproben-Ziel, genannt in `harness/README.md` §Sensors |

## Re-Evaluierungs-Trigger

Wenn golangci-lint angehoben wird und ein Linter seine Meldung, seine Einstellungen oder seine eingebauten Ausnahmen ändert; wenn ein neues Modul in `go.mod` kommt; wenn eine Stufe drei Pläne lang nicht aufgehoben wird.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-06 | Proposed | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
