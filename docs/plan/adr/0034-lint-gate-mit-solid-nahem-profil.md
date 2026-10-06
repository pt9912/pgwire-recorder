# ADR-0034: Lint-Gate mit SOLID-nahem Profil, eingeführt nach Bereinigung

**Status:** Accepted

**Datum:** 2026-10-06

**Autor:** pt9912

**Bezug:** [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md), [ADR-0001](0001-hexagonale-architektur.md)

**Schärft:** [`SPEC-049`](../../../spec/spezifikation.md#spec-049--lint-profil-lint)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`LH-QA-07` verlangt in Messmethode (1) eine statische Analyse nach festgelegtem Profil ohne Befund, mit Ausnahmen zentral und begründet, und in Messmethode (3) Unit-Tests außerhalb der geprüften Einheit, die nur über eine für Tests vorgesehene Stelle auf Internes zugreifen. Ein Lint-Gate gibt es nicht; `AGENTS.md` §3.2 ist ein Platzhalter. Vorbild ist das Profil des Schwester-Repos ai-harness-init mit golangci-lint v2. Der Nutzer hat am 2026-10-05 entschieden: die Schwellen des Vorbilds gelten unverändert, die Tests werden in vier Umstellungs-Plänen je Paketgruppe auf Black-Box-Pakete umgestellt, der Vertrag des Gates steht in der Spezifikation; am 2026-10-06, nach der Messung: Bereinigung vor dem Gate, ohne Stufen.

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
5. **Einführung nach Bereinigung, ohne Stufen** (Entscheidung des Nutzers vom 2026-10-06). `make lint` hängt erst an der Gate-Kette, wenn der Bestand bis auf die dauerhaften Ausnahmen aus Entscheidung 2 keinen Befund mehr hat; `.golangci.yml` trägt nie eine Ausnahme, die nur Bestand aussetzt. Die 71 Befunde in Testdateien beheben die vier Umstellungs-Pläne in den Dateien, die sie ohnehin umschreiben, die 25 im Produkt-Code zwei Bereinigungs-Pläne (Kern und Driven mit 13, Driving mit 12), nach den Umstellungen, weil Komplexitäts-Umbauten an unexportierten Funktionen erst dann keine White-Box-Tests mehr brechen. Die erwartete Bereinigung der neuen Kontexte ist `context.WithoutCancel`; ändert sie die Semantik des Herunterfahrens, geht der Befund an den Architect zurück.
6. **Werkzeug-Ziel vor dem Gate.** Damit jeder Bereinigungs-Plan sein Ergebnis messen kann, gibt es `make lint` mit Profil, Stufe `lint` und den eigenen Prüfungen schon vorher, als Werkzeug ohne Gate; es meldet den ganzen Bestand, und ein Plan liest die Befunde unter seinen Pfaden. Geliefert wird es vom ersten Plan der Bereinigungs-Reihe; das Gate (Gate-Kette, Gegenprobe, Doku) bleibt Gegenstand des Plans, der das Lint-Gate liefert. Ein direkter Aufruf des gepinnten Images in jedem Plan verworfen: Er dupliziert Pin, Profil und eigene Prüfungen in jedem Plan, und was gemessen wird, wiche vom späteren Gate ab.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — Bestand vor dem Gate bereinigen, Werkzeug-Ziel ohne Gate bis dahin** | Gate startet ohne Stufe; `.golangci.yml` trägt nur begründete dauerhafte Ausnahmen; die Testbefunde behebt der Umbau, der dieselben Dateien umschreibt | neuer Code bleibt bis zum Gate ungeprüft (akzeptiert) |
| B — alles in diesem Plan bereinigen | eine Lieferung | mehr als zwei Schichten, sieben Komplexitäts-Umbauten, nicht in einer Review-Sitzung prüfbar |
| C — Schwellen oder Linter an den Bestand anpassen | sofort grün | Lockerung nach `AGENTS.md` §3.6; gegen die Entscheidung des Nutzers |
| D — `testpackage` ganz aus bis zum letzten Umstellungs-Plan | eine Regel weniger | eine Ausnahme, die nur Bestand aussetzt |
| E — `errcheck`-Preset `std-error-handling` | weniger Konfiguration | blendet auch `Close` von Dateien aus, deren Fehler Daten kostet |
| F — `//nolint` mit `nolintlint` begrenzen | kein eigener Schritt | `nolintlint` verlangt Form, verbietet nicht |
| G — Gate jetzt, Rest als enge Stufen mit Hochschalt-Trigger (Empfehlung des Architect) | neuer Code ab sofort geprüft; jede Stufe hat Adresse | verworfen nach Entscheidung des Nutzers vom 2026-10-06: rund zwanzig Ausnahmen, die nur Bestand tragen, stünden im Profil, und Messmethode (1) gälte bis zu ihrer Aufhebung nur mit ihnen |

## Konsequenzen

- Positiv: Ab dem Gate macht jeder Befund und jedes `//nolint` den Lauf rot; jede Ausnahme im Profil hat einen Grund, der auch für neuen Code gilt.
- Negativ (akzeptiert): Bis zum Gate prüft niemand neuen Code automatisch; das Werkzeug-Ziel misst nur, wer es aufruft. Ob ein `Why:` zutrifft, prüft Review. Testdateien im Paket `main` lässt `testpackage` zu; `cmd/` hat keine. Ein `//nolint` in einem String-Literal ist ein Befund.
- Folgepflicht: Eine Gegenprobe hält jede Zusage aus [`SPEC-049`](../../../spec/spezifikation.md#spec-049--lint-profil-lint), die eine Mutation fangen kann, darunter die drei Rot-Fälle der Brücke und eine ungenutzte Regel.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| golangci-lint und eigene Prüfungen | Profil nach [`SPEC-049`](../../../spec/spezifikation.md#spec-049--lint-profil-lint) ohne Befund | `make lint`, bis zur Bereinigung Werkzeug ohne Gate, danach Gate-Ziel; genannt in `harness/README.md` |
| Gegenprobe | je Zusage ein Mutant, der rot wird | Gegenproben-Ziel, genannt in `harness/README.md` §Sensors |

## Re-Evaluierungs-Trigger

Wenn golangci-lint angehoben wird und ein Linter seine Meldung, seine Einstellungen oder seine eingebauten Ausnahmen ändert; wenn ein neues Modul in `go.mod` kommt; wenn die Bereinigung den Bestand nicht ohne neue Ausnahme grün bekommt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-06 | Proposed | — |
| 2026-10-06 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
