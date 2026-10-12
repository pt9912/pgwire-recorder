# Verifikation: slice-v1-abschluss-zeitangaben-aufzeichnen — 2026-10-12

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-12-review-slice-v1-abschluss-zeitangaben-aufzeichnen.md`, F-647 bis F-652); die Nacharbeit (`8ad1f3d`) und die Entscheidung des Architect (`fb88966`) prüfe ich nach. Ich repariere nichts.

**Gegenstand:** Slice-Plan `slice-v1-abschluss-zeitangaben-aufzeichnen` (Kopf, §1 bis §8) gegen den Gesamt-Diff `f4b7a97..8ad1f3d` (23 Dateien, 1074 Zeilen hinzu, 27 weg). Implementer: `9ed8d68`, `0bfbfe4`, `2e23f50`, `833960b`, `8ad1f3d`. Review: `8f040f3`. Architect: `fb88966`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

**Eingang:** der Plan ganz am Stand `8ad1f3d` (§7 *Belege des Implementers* nur als Beleg-Teile; §7 ist bis zur Closure unvollständig, das ist kein Befund); [`LH-FA-21`](../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen); [`LH-FA-21.a`](../../spec/spezifikation.md#lh-fa-21a--zeitangaben) *Aufzeichnen*, *Uhr*, *Messpunkt*; [`LH-FA-17.a`](../../spec/spezifikation.md#lh-fa-17a--konfiguration) (Optionstabelle); [ADR-0015](../plan/adr/0015-uhr-port.md); der Review-Report; der Code `internal/hexagon/services/record.go`, `internal/adapters/driving/cli/cli.go`, `internal/bootstrap`, `internal/hexagon/ports/driving/record.go`. Den Bericht des Implementers habe ich nicht gelesen. Beim Start war der Arbeitsbaum sauber, HEAD `8ad1f3d`.

**So wurden die Proben gebaut.** Alles Docker-only, Scratch im Scratchpad der Sitzung, nie im Repo. Mutanten: je Mutant eine frische Kopie des Baums ohne `.git` (`tar`, entpackt mit `-m`, keine Zeitstempel übernommen), genau eine `perl`-Ersetzung (das Skript bricht ab, wenn sie nichts ändert), `gofmt -l` im selben Lauf (Befunde stehen unten, wo sie auftraten), `go test -count=1 -run …` im Image `pgwire-recorder:test` mit `--network none`. Mutanten der Verdrahtung: `make test-integration` in der Kopie. Ende-zu-Ende: das mit `make build` gebaute Binary (`pgwire-recorder:dev`) in einem eigenen internen Docker-Netz, gegen `postgres:17-alpine` mit dem gepinnten Digest aus `harness/mk/integration.mk`, Client `psql` 17.11 mit `\! sleep` zwischen den Anfragen. Container und Netze trugen die Namen `vfy-*` und wurden namentlich entfernt.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `record --record-timing` trägt `offset_ms` ein. **Bestätigt.**

Alle Tests des Plans existieren (`grep "^func Test…("` je Name aus §7, genau eine Fundstelle; 28 Namen). Ohne den Fix aus dem richtigen Grund rot: Die Mutanten sind unabhängig von §7 neu gebaut, jeder bricht nur die Zusage im Wortlaut von `LH-FA-21.a`.

| Zusage (Wortlaut) | Mutant (nur sie) | rot, mit Meldung |
|---|---|---|
| Unterschied zweier Lesungen, Ankunft minus Beginn | `l.beginn - ankunft` | `TestRecordZeitEinfach` (`[-40 -45]`), `TestRecordZeitBeginnVorAufbau`, `TestRecordZeitAnkunftVorUpstream`, `TestRecordZeitAnkunftVorWarten` |
| Beginn ist der Aufruf von `OpenSession`, vor dem Aufbau der Upstream-Session | Lesung nach `upstream.Open` | `TestRecordZeitBeginnVorAufbau` (`[0]` statt `[30]`) |
| Jede Session misst von ihrem eigenen Beginn | Beginn konstant 0 | `TestRecordZeitEinfach` (`[140 145]`), `TestRecordZeitNullSteht`, `TestRecordZeitJedeSessionEigenerBeginn` (`Session b: [105]` statt `[5]`) |
| Es wird nicht gerundet | Differenz auf 10 ms aufgerundet | `TestRecordZeitEinfach` (`[40 50]` statt `[40 45]`), `TestRecordZeitLebendpruefung`, `TestRecordZeitJedeSessionEigenerBeginn` |
| Die erste Interaktion kann 0 tragen, das Feld steht dann da | `nil` bei Wert 0 | `TestRecordZeitNullSteht` (`[-1]` statt `[0]`), `TestRecordZeitAnkunftVorWarten` |
| Ankunft einer einfachen Anfrage ist der Aufruf von `Query` | Lesung beim Übernehmen der Antwort | `TestRecordZeitAnkunftVorUpstream` (`[510]` statt `[10]`), `TestRecordZeitAnkunftVorWarten` |
| Ankunft einer Extended-Interaktion ist die erste Client-Nachricht der ersten Gruppe, gleich welche Art | Lesung beim Abschluss der Interaktion | `TestRecordZeitAnkunftVorWarten`, `TestRecordZeitExtendedErsteNachricht` (Close, Flush, Sync, Parse, Bind, Describe, Execute: `[130]` statt `[20]`) |
| … der ersten Gruppe (nicht jeder) | Lesung bei der ersten Nachricht jeder Gruppe | `TestRecordZeitExtendedErsteNachricht` (`[25]` statt `[20]`, sechs Arten), `TestRecordZeitExtendedMehrereGruppen` |
| Die Werte einer Session steigen nicht ab | Wert der Extended-Interaktion um 100 ms je abgeschlossener Interaktion verschoben | `TestRecordZeitNichtAbsteigend` (`[10 -80 -170 240]` statt `[10 20 30 240]`) |
| Die Lebendprüfung trägt das Feld wie jede Anfrage | `nil` für Text, den `istLebendpruefung` annimmt | `TestRecordZeitLebendpruefung` (`[-1 -1]` statt `[3 7]`); die Fixture (`-- ping` und der leere Text) trifft den Pfad der Lebendprüfung |
| Eine unvollständig endende Interaktion wird nicht übernommen, hat keinen Wert | `CloseSession` hängt die offenen Interaktionen an | `TestRecordZeitUnvollstaendigOhneWert` (`[10 20]` statt `[10]`) |
| Ohne die Option fehlt das Feld | `offset` liefert immer einen Wert | `TestRecordOhneZeitangaben` (`[0 0]`, erwartet ohne Feld) |
| Ohne die Option liest der Service die Uhr nie | `lesung` liest immer | `TestRecordOhneZeitangaben` (`4 Lesungen ohne Zeitaufzeichnung`) |
| Dieser Slice ruft `WarteBis` nicht auf (§1) | `lesung` ruft `WarteBis` nach `Jetzt` | `TestRecordZeitWartetNie` (`WarteBis aufgerufen`) |
| Roundtrip sqlite: aufzeichnen, laden, Werte gleich | Schreiber schreibt Wert + 1 | `TestRecordZeitRoundtrip` (`sqlite, Interaktion 1: offset_ms 13, erwartet 12`) |
| Roundtrip yaml: aufzeichnen, laden, Werte gleich | Schreiber lässt `offset_ms` weg | `TestRecordZeitRoundtrip` (`yaml, Interaktion 1: ohne offset_ms, erwartet 12`) |
| Der Bootstrap verdrahtet die Systemuhr (ohne sie wäre der Unterschied 0) | `Jetzt()` der Systemuhr liefert 0 (`make test-integration`) | `TestE2ERecordTiming` (`rec.yaml: offset_ms [0 0 0]`, ebenso `rec.sqlite`) |
| Der Bootstrap reicht die Option durch | `Zeitangaben(NeueSystemuhr(), false)` | `TestE2ERecordTiming` (`[-1 -1 -1]`), `TestE2ERecordTimingQuellen` (Umgebung, Datei) |
| Ohne die Option ist die Ausgabe unabhängig vom Zeitverhalten | `Zeitangaben(NeueSystemuhr(), true)` | `TestE2ERecordOhneTiming` (`yaml: offset_ms [4 5 5], erwartet ohne Feld`, `yaml unterscheidet sich zwischen den Läufen`, sqlite `[5 5 6]`), `TestE2ERecordTimingQuellen` (`false` auf der Kommandozeile) |
| Kein Paket unter `./cmd/...` importiert die Testhilfe `offsetlese` (Kommentar von `doc.go`) | leerer Import der Testhilfe in `internal/bootstrap/bootstrap.go` | `TestBinaryOhneTesthilfen` (rot; der erste Versuch war nicht gofmt-sauber durch die Importreihenfolge, der Lauf mit richtiger Reihenfolge ist gofmt-sauber und rot) |

Die Zeile „Ohne die Option lädt die Aufzeichnung beider Formate ohne Feld“ in §7 (kein eigener Mutant) bestätigt sich: Der Zweig `mit=false` des Roundtrip-Tests ist durch `TestRecordOhneZeitangaben` und `TestE2ERecordOhneTiming` gebunden, beide rot in meinen Mutanten. Die grüne Zeile *Lebendprüfung* in §7 („der Code kennt keinen Sonderpfad“) stimmt: In `record.go` steht kein `istLebendpruefung`; mein Mutant legt einen an und wird rot.

Fehlende Beleg-Teile (§7 verlangt, DoD verlangt): keine. Die Mutationstabelle führt die Verbotsliste der Testhilfe jetzt (F-648 behoben, Zeile in §7).

### Punkt 2: Option `--record-timing`. **Bestätigt.**

Mutanten am CLI-Adapter (`internal/adapters/driving/cli/cli.go`), alle rot:

| Zusage | Mutant | rot |
|---|---|---|
| Default `false` | Default `true` | `TestParseRecord`, `TestLeserAlleOptionen` (`RecordTiming:true`), `TestRecordTimingQuellen` |
| Option setzt `RecordTiming` | `RecordTiming = false` | `TestLeserAlleOptionen`, `TestRecordTimingQuellen` |
| Reihenfolge nach `--format` | Option aus `optionen("record")` entfernt | `TestLeserOptionen` (`erwartet [… format record-timing shutdown-timeout …]`), `TestRecordTimingQuellen` |
| bei `replay` und `play` unbekannt | Option bei `replay` mit angemeldet | `TestDateiUngueltig` (`record_timing in replay`), `TestConfigShowFehler`, `TestLeserAlleOptionen`, `TestLeserOptionen` |
| Hilfe nennt die Umgebungsvariable | Hilfetext mit anderem Namen | `TestRecordTimingHilfe` (`Hilfe nennt "PGWIRE_RECORDER_RECORD_TIMING": false`) |

Am gebauten Binary (`pgwire-recorder:dev`, `record` gegen PostgreSQL 17):

- **Quellen.** Option `--record-timing` → `offset_ms` je Interaktion; Umgebungsvariable `PGWIRE_RECORDER_RECORD_TIMING=true` → ebenso; Konfigurationsdatei `record:\n  record_timing: true` (über `PGWIRE_RECORDER_CONFIG`) → ebenso. **Rangfolge:** Umgebung `true` mit `--record-timing=false` → Feld fehlt (Option vor Variable); Umgebung `false` mit Datei `true` → Feld fehlt (Variable vor Datei).
- **Ungültiger Wert `GEHEIMWERT`.** Option: `Konfiguration [PGR-E2001]: ungültige Verwendung von record: invalid boolean value "GEHEIMWERT" for -record-timing: erlaubt sind true und false`, Exit 2 (nennt den Wert, wie bei allen Wahrheitswert-Optionen; `LH-FA-17.a` erlaubt das für die Kommandozeile). Umgebung: `Umgebungsvariable PGWIRE_RECORDER_RECORD_TIMING: erlaubt sind true und false`, Exit 2, ohne Wert. Datei: `[PGR-E2004] Konfigurationsdatei: record.record_timing: erlaubt sind true und false`, Exit 2, ohne Wert.
- **`replay` und `play`.** `replay --record-timing` und `play --record-timing` enden mit `PGR-E2001` (`flag provided but not defined: -record-timing`), Exit 2; ihre `--help` nennen die Option nicht (0 Treffer). Der Schlüssel `record_timing` unter `replay:` ist `PGR-E2004` (`replay.record_timing: unbekannter Schlüssel`), Exit 2.
- **`--help` und `config show`.** `record --help` nennt `--record-timing[=true|false]`, den Text und die Variable; die allgemeine Hilfe nennt sie (2 Treffer). `config show --config <datei>` zeigt `record:\n  record_timing: true` und den Namen `PGWIRE_RECORDER_RECORD_TIMING` ohne Wert.

### Punkt 3: `make gates` grün. Siehe Verdikt (Lauf nach dem Commit dieses Berichts).

---

## 2. Ende-zu-Ende mit dem gebauten Binary

Szenario: ein `psql` pro Lauf in einer Sitzung, `select 1`, Pause 300 ms, `select 2`, Pause 200 ms, `select 3`, Pause 300 ms, `select $1::int \bind 5 \g` (Extended: Parse, Bind, Describe, Execute, Sync).

| Lauf | `offset_ms` je Interaktion | Beurteilung |
|---|---|---|
| `--record-timing`, yaml | `[5 307 509 812]` | Abstände 302, 202, 303 ≥ Pausen 300, 200, 300; nicht absteigend; erste Interaktion 5 (der Wert kann 0 sein, das belegt der Unit-Test) |
| `--record-timing --format sqlite` | `[3 305 507 809]` | geladen mit den geprüften Lesern (`offsetlese` in einer Kopie): Abstände 302, 202, 302 |
| Umgebungsvariable, yaml | `[4 306 508 809]` | Abstände 302, 202, 301 |
| Konfigurationsdatei, yaml | `[4 305 507 809]` | Abstände 301, 202, 302 |
| ohne Option, yaml (zwei Läufe mit verschiedenen Pausen und verschiedener Reihenfolge der Pausen: 0,5 s zwischen Anfrage 2 und 3, 0,3 s zwischen Anfrage 1 und 2) | `[-1 -1 -1]` (Feld fehlt, `grep offset_ms` leer) | `cmp ohne0.yaml ohne1.yaml`: **byte-gleich** |
| ohne Option, sqlite | `[-1 -1 -1 -1]` (Spalte NULL) | Feld fehlt |
| Option `false` vor Umgebung `true`; Umgebung `false` vor Datei `true` | `[-1 -1 -1 -1]` | Rangfolge bestätigt |

- **Extended-Interaktion:** Die vierte Interaktion (`type: extended`) trägt `offset_ms` auf der Ebene der Interaktion (nicht der Gruppe), gemessen vor der Pause bis zum `Parse`, nicht beim `Sync`. Mit den Nachrichtenarten Parse, Bind, Describe, Execute, Sync ist sie am Binary belegt; die Arten Close und Flush sowie mehrere Gruppen belegen nur die Unit-Tests (`TestRecordZeitExtendedErsteNachricht`, `TestRecordZeitExtendedMehrereGruppen`). Mein `psql` 17 kennt `\parse`, `\bind_named` und `\close_prepared` nicht (erst Version 18), so blieben diese Arten am Binary ungeprüft.
- **Abgebrochene Session und unvollständige Interaktion.** `record --record-timing --shutdown-timeout 300ms`, Client in `pg_sleep(5)`, `SIGTERM` an den Recorder: `PGR-E4006` (`Interaktion 2 nicht abgeschlossen und verworfen; die Session wird als Session 1 geschrieben`), die Datei enthält nur Interaktion 1 mit `offset_ms: 4`; die verworfene hat keinen Wert und keine Zeile. Bricht dagegen nur der Client ab und der Recorder läuft weiter, endet die Anfrage am Upstream, und die Interaktion wird vollständig mit Wert (`205`) aufgezeichnet; das ist das Verhalten von [`LH-FA-13.a`](../../spec/spezifikation.md#lh-fa-13a--signalbehandlung) und keine Abweichung.
- **Lebendprüfung am Binary:** nicht belegt. `psql` sendet keinen Text, den `istLebendpruefung` annimmt (nur Kommentar, nur Leerraum); Beleg ist der Unit-Test mit einer Fixture, die den Pfad trifft (V-191).
- **Round-Trip.** `play --upstream … --input opt.yaml` und `… opt.sqlite` (Dateien mit `offset_ms`, ohne `--keep-timing`): Exit 0, beide Dateien lassen sich laden, die Prüfung absteigender Werte besteht. `replay --input opt.yaml` und `opt.sqlite`: liefert die drei Antworten `1`, `2`, `3`; die vierte Interaktion blieb unverbraucht (`PGR-W2001`, weil der Client sie nicht schickte; erwartet).
- **Mutanten der Verdrahtung am Binary** (`make test-integration` in Kopien): `Jetzt()` liefert 0 → `TestE2ERecordTiming` rot (`[0 0 0]` in beiden Formaten); `Zeitangaben(NeueSystemuhr(), false)` → `TestE2ERecordTiming`, `TestE2ERecordTimingQuellen` rot; `Zeitangaben(NeueSystemuhr(), true)` → `TestE2ERecordOhneTiming`, `TestE2ERecordTimingQuellen` rot.

---

## 3. Spezifikation, Sicht, Plan gegen Code

- **`LH-FA-21.a`** *Aufzeichnen*, *Uhr*, *Messpunkt* stimmen mit dem Code überein: Der Use Case misst (`record.go`: `beginn := s.lesung()` als erste Anweisung von `OpenSession`, `ankunft := s.lesung()` als erste Anweisung von `Query` und `ClientMessage`); die Differenz ist nicht gerundet; ohne Option liest der Service die Uhr nie; ein Wert wird nicht für unvollständige Interaktionen vergeben. Die Bedingung des Architect (`fb88966`: Aufrufe der Client-Richtung kommen nacheinander) steht im Text von *Messpunkt* (mit der Grenze „ein Test für den Adapter gibt es nicht“), im Kommentar des Ports `Recorder`, im Kommentar von `Query` und `ClientMessage` (`8ad1f3d`) und in der Sicht (`spec/architecture.md`, Zeile *Record-Use-Case*). Die Formulierung sagt nicht mehr zu, als ein Test prüft (`AGENTS.md` §3.11): Die Eigenschaft ist eine Voraussetzung an den Aufrufer, der Use Case sagt für gleichzeitige Aufrufe nichts zu.
- **Sicht** (`spec/architecture.md`): meilensteinfrei, keine Kennung einer ADR; nennt die Voraussetzung als Zeile des Use Case. Der Code folgt.
- **Plan §1/§3/§6/§8 folgen dem Code.** Zwei Liefer-Punkte (Messen und Eintragen im Kern mit der Verdrahtung; Option), zwei Schichten (Kern, CLI-Adapter). `internal/bootstrap/bootstrap.go` ruft `services.Zeitangaben(NeueSystemuhr(), o.RecordTiming)` ohne Verzweigung: Verdrahtung, keine Schicht (`AGENTS.md` §3.13). Der PGWire-Adapter ist unberührt (`git diff --stat f4b7a97..HEAD` nennt keine Datei unter `internal/adapters/driving/pgwire`); der Port `Recorder` ändert nur seinen Kommentar. §3 nennt `services.Zeitangaben` und `offsetlese` (F-649 behoben, Kommentarzeilen zu F-647 stehen mit §3 im selben Commit `8ad1f3d`). Jeder Dateiname und jeder Name in §3, §7 und der DoD steht im Code (Namensabgleich per `grep`, 28 Testnamen je genau eine Fundstelle; `offsetlese`, `NeueSystemuhr`, `Zeitangaben`, `testuhr.Neu`, `Vorruecken` vorhanden).
- **`TestBinaryOhneTesthilfen`:** Die Testhilfe `internal/bootstrap/offsetlese` erreicht das Binary nicht; der Mutant (leerer Import im Bootstrap) macht den Test rot, ohne Mutant ist er grün.
- **Nehmer** nehmen an (`AGENTS.md` §3.13): `slice-v1-abschluss-zeitangaben` nennt in §1 die Abgabe des Aufzeichnens an diesen Slice mit der Kennung des Gebers und den Port mit `Jetzt` und `WarteBis` (V-186 aus der Verifikation von `slice-v1-abschluss-zeitangaben-uhr` ist damit für diesen Nehmer erledigt, er trägt die Namen in Zeile 39); `slice-v1-abschluss-zeitangaben-doku` nennt dieses Aufzeichnen in Ziel, DoD-Punkten und §3 mit der Kennung. Beide nehmen an, kein Ausschluss trifft die Sendungen; der Diff ändert sie nicht.
- **Handbuch und README unangetastet** (`git diff --stat f4b7a97..HEAD -- docs/user/benutzerhandbuch.md README.md` leer; `docs/user/abdeckung-*.md` sind erzeugte Tabellen). Keine `//nolint` im Diff (`grep` leer über die Code-Dateien; die Treffer des Diff-`grep` stehen im Review-Report, der den Begriff nennt). Keine Kennung eines Spezifikations- oder Sichtstratums in den Commit-Messages `f4b7a97..HEAD` (`git log --format=%B | grep "SPEC-\|ARC-"` leer); Kennung `slice-v1-abschluss-zeitangaben-aufzeichnen` und `LH-FA-21` in jeder ([ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md)).
- **Sensoren.** `make abdeckung-check`: Exit 0. `make kopf-check`: Exit 0. Die Abdeckungstabellen führen `LH-FA-21` als `vollständig`, obwohl `--keep-timing` aussteht; das ist die mechanische Folge der Deklarationen ([ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md), F-651) und kein Befund (V-190, vermerkt). `make gates`: Ergebnis im Verdikt.

---

## 4. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-189 | LOW | Die Spezifikation nennt in *Aufzeichnen* als Beginn „die Annahme der Verbindung“, in *Messpunkt* den Aufruf von `OpenSession` nach TLS und Startup („was der Client davor tut, zählt nicht mit“). Das Binary misst nach *Messpunkt* (erster Wert 3 bis 5 ms bei sofortiger erster Anfrage). Ein Leser von *Aufzeichnen* allein erwartete einen früheren Beginn; derselbe Wortlaut steht im Ziel von §1 des Plans. *Aufzeichnen* bleibt die kürzere Aussage, *Messpunkt* gilt, aber der Satz widerspricht ihm im Wortlaut. Die Handbuch-Fassung (`slice-v1-abschluss-zeitangaben-doku`) sollte die Lesart von *Messpunkt* übernehmen. | `spec/spezifikation.md` · `LH-FA-21.a` *Aufzeichnen*, Satz „zwischen der Annahme der Verbindung“; Plan §1 *Ziel* | Planner (Wortlaut der Spezifikation und von §1 angleichen, ohne Blockade; ein Spezifikationsfall, nicht des Codes) |
| V-190 | INFO | Die Abdeckungstabellen führen `LH-FA-21` mit Stand `vollständig`, während `--keep-timing` (Einspielen, `slice-v1-abschluss-zeitangaben`) noch aussteht. Die Zeile folgt mechanisch aus den Deklarationen (Happy, Boundary, Negative belegt); das Werkzeug prüft nicht, ob die Anforderung als ganze belegt ist. Wie F-651; hier nur vermerkt. | `docs/user/abdeckung-gesamt.md`, `docs/user/abdeckung-vollstaendig.md` · `LH-FA-21` | ohne Aktion |
| V-191 | INFO | Die Lebendprüfung trägt `offset_ms` am Binary nicht belegt: `psql` sendet keinen Text, den `istLebendpruefung` annimmt, und `TestE2ERecordTiming` hat keinen Fall dafür. Der Beleg ist der Unit-Test `TestRecordZeitLebendpruefung` mit einer Fixture, die den Pfad trifft (`-- ping`, leerer Text; mein Mutant „`nil` für Lebendprüfung“ ist rot). Der Code kennt keinen Sonderpfad für sie, die Lücke ist klein. | `test/integration/zeitangaben_aufzeichnen_e2e_test.go` | ohne Blockade; Hinweis für den Planner, ob ein Binary-Fall verlangt wird |
| V-192 | INFO | Zwei Proben-Einschränkungen, kein Befund gegen den Slice: (a) Die Arten `Close` und `Flush` sowie mehrere Gruppen einer Extended-Interaktion habe ich am Binary nicht gefahren (`psql` 17 hat die Befehle nicht); Beleg Unit (`TestRecordZeitExtendedErsteNachricht`, `TestRecordZeitExtendedMehrereGruppen`, je Art ein Mutant rot). (b) Die ersten Fassungen der Mutanten „Beginn konstant 0“ (Ausrichtung des Literals) waren nicht gofmt-sauber und zählen nicht; Beleg ist die Fassung `beginn := int64(0)`, gofmt-sauber und rot (`[140 145]`). | — | ohne Aktion |

Falsifikationsversuche, die nichts fanden: Basis je Session (zwei Sessions gleichzeitig, `TestRecordZeitJedeSessionEigenerBeginn`); `offset_ms: 0` bleibt stehen (nicht `nil`); Wert mit Option `false` auf allen drei Quellen; Wert bei `--format sqlite` mit NULL ohne Option; Extended-Interaktion nach einfacher Anfrage (Wert nicht kleiner); `replay`/`play` kennen die Option nicht; Meldung zu Umgebung und Datei enthält den Wert nie.

---

## Verdikt

**Befunde:** 0 HIGH, 0 MEDIUM, 1 LOW (V-189), 3 INFO (V-190, V-191, V-192).

Die zwei Liefer-Punkte der DoD sind **bestätigt**: Die Tests existieren, werden ohne den Fix aus dem richtigen Grund rot (25 Mutanten: 14 im Service, 5 im CLI-Adapter, 2 am Schreiber für den Roundtrip, 3 an der Verdrahtung mit `make test-integration`, 1 an der Verbotsliste, alle neu gebaut im Wortlaut von `LH-FA-21.a` und der Optionstabelle, keiner aus falschem Grund rot), das Binary zeichnet mit `--record-timing` in `yaml` und `sqlite` Werte, deren Abstände die Clientpausen einhalten (302, 202, 303), ohne die Option fehlt das Feld und zwei Läufe mit verschiedenem Zeitverhalten sind byte-gleich, die Quellen und die Rangfolge der Option stimmen, die Aufzeichnung lädt in `replay` und `play`. Der Plan folgt dem Code (zwei Liefer-Punkte, zwei Schichten, der Bootstrap ist reine Verdrahtung); die Entscheidung des Architect zu F-647 steht in Spezifikation, Sicht, Port und Plan.

**Closure-blockierend:** nichts im Code. Vor der Closure: V-189 (Planner), die Closure-Pflichten aus §2 des Plans (Review-Häkchen, Notiz mit Lerneintrag, Register, Ausgänge der zwei Risiken: Risiko 1 erwartet *entfallen* (Untergrenzen statt Dauer, kein Flackern nach oben prüfbar, Grund steht in §6 und in der Grenze (a) von §7), Risiko 2 *entfallen* (Messpunkt im Use Case, zwei Schichten bestätigt)), die Finding-Klassen des Reviews (F-647 bis F-652) in der Notiz.

**Finding-Klasse dieses Laufs:** keine neue; „Zusage weiter als der Test“ (V-175, V-182) trat nicht auf, die Fixtures treffen den Pfad.

**Aufgeräumt:** nur benannte Verzeichnisse unter `m-*` und `b-*` im Scratchpad der Sitzung; Container `vfy-*`, `pgr-it-*` und die Netze `vfy-zeit-net` und die des Runners sind entfernt. Der Arbeitsbaum des Repos war bis zu diesem Bericht unverändert.

**Gate:** `make gates` auf dem sauberen Baum nach dem Commit dieses Berichts: Ergebnis in der Übergabe-Nachricht (ein Gate-Lauf belegt nur den Stand, auf dem er lief).
