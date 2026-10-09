# Verifikation: slice-v1-abschluss-konfigurationsdatei — 2026-10-09

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-09-review-slice-v1-abschluss-konfigurationsdatei.md`, F-506 bis F-518). F-516 und F-518 hat er an mich übergeben.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-konfigurationsdatei.md` im geschnittenen Zuschnitt (*Datei*), Kopf und §1 bis §8, gegen den Gesamt-Diff `f333934..c5b7d96`. Architect: `f19b440`, `096e51e`, `647f0ff`, `cfc4d6b`, `5cbfd1b`, `8d320cc` ([ADR-0037](../plan/adr/0037-yaml-bibliothek-fuer-die-anzeige-der-konfigurationsdatei.md) `Accepted`). Schnitt: `fdf1cd8`. Implementer: `7a80393`, `a063f2f`, `9332356`, `c5b7d96`. Review: `ef8bee5`.

**Abgrenzung:** Verbindungen, Grammatik der URL, Platzhalter und das Erzeugen von `PGR-E2005` und `PGR-E2006` liegen bei `slice-v1-abschluss-verbindungen-platzhalter`. Dass sie fehlen, ist kein Befund.

**Eingang:**

- der Plan ganz, besonders §1, §2 (DoD), §3, §4 (Schnitt), §6 (*Randformen*, *Rückgaben vom 2026-10-09*, *Befunde des Reviews vom 2026-10-09*, *Risiken*) und §7 (*Belege des Implementers*)
- `spec/spezifikation.md` `LH-FA-17.a` ganz, mit *Konfigurationsdatei*, *Fehler*, *Dauer*, *Anzeige* und der Optionstabelle
- [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md), [ADR-0036](../plan/adr/0036-yaml-bibliothek-fuer-die-konfigurationsdatei.md), [ADR-0037](../plan/adr/0037-yaml-bibliothek-fuer-die-anzeige-der-konfigurationsdatei.md), `.a-check.yml` (Abschnitt `tech`), der Diff von `spec/architecture.md` (`ARC-013`, §2)
- am Stand `c5b7d96`: `internal/adapters/driving/cli/datei.go` ganz, der Diff von `cli.go`, `export_test.go`, `internal/bootstrap/bootstrap.go` und `internal/hexagon/model/fehler.go`; von den Tests `datei_test.go` (Wahl, gültige und ungültige Dateien, Tag, Kodierung, Anzeige ohne BOM), `datei_unix_test.go`, `dreiQuellen` und `tabelle()` in `leser_test.go`, `TestE2EConfigShowStdin`
- `docs/user/benutzerhandbuch.md` §4 *Die gewählte Konfigurationsdatei anzeigen*, §5 *Konfigurationsdatei* und *Fehlercodes*
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen
- der Nehmer `slice-v1-abschluss-verbindungen-platzhalter` (§1, DoD, §6, nur die Stellen zu `config show`)

Beim Start war der Arbeitsbaum sauber, HEAD `c5b7d96`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-09

**So wurden die Proben gebaut.** Alles lief in einem eigenen Unterverzeichnis des Scratchpads, nie im Repo.

- **Binary:** eine frische Kopie per `git archive c5b7d96 | tar -x`, ohne `.git` und ohne `cp -p`. Aus ihr wurde ein Image der Stufe `runtime` mit eigenem Tag gebaut und das Binary herauskopiert. Die Proben liefen im gepinnten Go-Image (es bringt `sh`, `mkfifo` und `timeout` mit), mit `--network none` und als Benutzer `65532`, das Binary schreibgeschützt eingebunden. Jeder Aufruf hatte eine Frist von 10 s.
- **Mutanten:** je Mutant eine neue Kopie per `git archive c5b7d96 | tar -x`. Die Ersetzung machte ein Skript, das abbricht, wenn der Suchtext nicht genau einmal trifft. Danach lief `go test -count=1` für `internal/adapters/driving/cli` und `internal/bootstrap` im Image der Stufe `deps`, mit der Kopie als Bind-Mount und ohne Netz. Danach wurde nur das Verzeichnis des Mutanten gelöscht. Der Grundlauf ohne Mutation war grün. Den E2E-Mutanten habe ich als Image der Stufe `integration` aus der mutierten Kopie gebaut und mit `-test.run '^TestE2EConfigShow|^TestE2EReplayKonfigurationsdatei'` laufen lassen. Die ungeänderte Kopie lief zum Vergleich genauso. `PGR_OHNE_POSTGRES=1` überspringt dabei nur das Warten auf PostgreSQL; die drei Tests verbinden sich nicht.
- **Danach:** Die eigenen Images und Kopien sind entfernt. `make gates` und `make a-check` liefen im Arbeitsbaum. `git status` war danach sauber.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: Hilfe vor jeder Prüfung, auch für `config show` und `--config`, samt `--`; `config show` in der Form aus `LH-FA-17.a` (`LH-FA-01`, `LH-FA-17`). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Die Hilfe geht jeder Prüfung vor, auch bei `config show` und `--config` | Binary: `config show --config gibtsnicht --help` gibt die Hilfe von `config show` aus, Exit 0. `PGWIRE_RECORDER_CONFIG=<FIFO> record --help` gibt die Hilfe aus, Exit 0, ohne zu warten. §7 DoD 1, Zeile 1 (`TestDateiHilfeVorPruefung`, `TestDateiNichtLesbar`, `TestDateiWahl`) | bestätigt |
| `config --help` ist die Hilfe von `config show` | Binary: Die erste Zeile lautet `Aufruf: pgwire-recorder config show [optionen]`. §7 DoD 1, Zeile 2 | bestätigt |
| `--` beendet die Optionen | Binary: `config show -- --help` ergibt `PGR-E2001: unerwartetes Argument "--help"`, Exit 2 | bestätigt |
| `config` ohne `show`, ein anderes Unterkommando und ein Argument nach `show` sind `PGR-E2001` | Binary: `config`, `config zeige` und `config show extra` enden je mit `PGR-E2001`, Exit 2. `config show --listen 1` ergibt `PGR-E2001` mit `flag provided but not defined`. `config show --config=` ergibt `PGR-E2001` mit `ein leerer Wert ist ungültig` | bestätigt |
| Erste Zeile: Pfad, wie gewählt; ohne Datei sagt sie das | Binary: `--config ./sub/../sub/hb.yaml` ergibt genau diesen Pfad, `/dev/stdin` bleibt `/dev/stdin`, `link.yaml` bleibt `link.yaml`. Ohne Datei lautet sie `keine Konfigurationsdatei gefunden`, Exit 0 | bestätigt |
| Inhalt als YAML, zwei Leerzeichen Einzug, Reihenfolge der Datei, ohne Kommentare, ohne BOM, mit `\n` | Binary: Eine Datei mit vier Leerzeichen Einzug unter `replay:`, Kopf- und Zeilenkommentaren kommt mit zwei Leerzeichen und ohne Kommentare heraus. Die Reihenfolge bleibt (`log_level`, `connections`, `replay`, `record`). `"null"` und `'true'` behalten ihre Anführungszeichen, `${STAGING_PASSWORD}` bleibt unaufgelöst. Dateien mit BOM und `\r` bzw. mit BOM und `\r\n`: `od -c` zeigt weder BOM noch `\r` | bestätigt |
| Aktive Variablen: nur Namen, sortiert, gesetzt und nicht leer, mit dem Präfix am Anfang, auch ohne Option | Binary: `PGWIRE_RECORDER_ZZZ`, `PGWIRE_RECORDER_AAA`, `PGWIRE_RECORDER_LEER=` und `XPGWIRE_RECORDER_MITTE` ergeben genau `PGWIRE_RECORDER_AAA` und `PGWIRE_RECORDER_ZZZ`, ohne Wert. Neben `--config` erscheint ein gesetztes `PGWIRE_RECORDER_CONFIG` mit Namen | bestätigt |
| Geprüft wird nur `PGWIRE_RECORDER_CONFIG` | Binary: `PGWIRE_RECORDER_LOG_LEVEL=GEHEIM` und `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT=GEHEIM` neben `config show` ergeben Exit 0, beide Namen stehen in der Liste | bestätigt |
| Ungültige Datei: Fehlercode des Ladens, keine Anzeige | Binary: Die Datei mit doppeltem `log_level` ergibt `PGR-E2004`, stdout ist leer, Exit 2 | bestätigt |

### Punkt 2: Wahl genau einer Datei, Priorität über drei Quellen, `PGR-E2004`, Konstanten `PGR-E2005` und `PGR-E2006` (`LH-FA-17`). **Bestätigt.** Zur Zeilennummer in der Meldung siehe V-118 und V-119.

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `--config` vor `PGWIRE_RECORDER_CONFIG` vor Standarddatei, genau eine Datei | Binary: Mit `PGWIRE_RECORDER_CONFIG=gibtsnicht.yaml` und `--config regulaer.yaml` wird `regulaer.yaml` gezeigt, Exit 0; die Datei aus der Umgebung wird nicht geprüft. Allein ergibt `PGWIRE_RECORDER_CONFIG=gibtsnicht.yaml` `PGR-E2004 … (PGWIRE_RECORDER_CONFIG) nicht vorhanden`. Eine leere `PGWIRE_RECORDER_CONFIG` gilt als nicht gesetzt | bestätigt |
| Die Standarddatei existiert, wenn ihr Pfad Links gefolgt auflöst (R2) | Binary: Die Standarddatei als Link ins Leere ergibt bei `config show` die Zeile `keine Konfigurationsdatei gefunden`, Exit 0. Bei `replay` folgt `Pflichtoption --input fehlt`; es wurde keine Datei gelesen. Eine Schleife von Links ergibt `PGR-E2004 … (.pgwire-recorder.yaml) nicht lesbar: too many levels of symbolic links`. `--config` auf einen Link ins Leere ergibt `PGR-E2004 … nicht vorhanden` | bestätigt |
| Priorität Kommandozeile vor Umgebung vor Datei vor Default für jede angemeldete Option, auch `fail_on_unconsumed`, Frist und `log_level` | `dreiQuellen` in `leser_test.go` gelesen: Datei allein, dann Umgebung neben Datei, dann Kommandozeile neben Datei, je für jede Option aus `cli.Optionen`. Die Defaults kommen aus `tabelle()`, der Ort des Schlüssels aus `ortDerDatei`, unabhängig vom Feld `oben`. **V-OBEN nachgefahren** (`log_level` im Abschnitt statt oben): **rot**, `TestLeserAlleOptionen`, `TestDateiGueltig`, `TestDateiWahl` und drei weitere | bestätigt |
| `log_level` nur oben, Wertemenge und Strenge von `--log-level` | Binary: `replay:` mit `log_level` ergibt `PGR-E2004 … replay.log_level: unbekannter Schlüssel`. `log_level: GEHEIM` ergibt `PGR-E2004 … log_level: erlaubt sind error, warn, info und debug` | bestätigt |
| Eine ungültige Datei ist `PGR-E2004` (Schlüssel, Abschnitt, Wertform, YAML) | Binary, je Exit 2 mit `PGR-E2004`: `Log_Level`, `replay.output`, `replay.config`, `play: {}`, `record:` ohne Inhalt, `---` allein, zwei Dokumente, `Null`, `NULL`, leerer Wert, Anker, Alias, Merge-Schlüssel `<<`, Liste als Wert, ungültige Dauer und ungültiger Wahrheitswert. Gültig sind die leere Datei, eine nur mit Kommentar und `output: "~"` | bestätigt |
| Die Code-Tabelle führt `PGR-E2004` bis `PGR-E2006`, ohne Erzeuger für die beiden letzten | `internal/hexagon/model/fehler.go`: drei Konstanten ohne Logik. `grep` nach `CodeConfigVariable` und `CodeConfigPassword` außerhalb der Tests findet nur die Konstanten | bestätigt |

### Punkt 3: Vorgaben des Architect vom 2026-10-09, je Test mit roter Mutation (`LH-FA-17`). **Bestätigt.**

| Vorgabe | Mutation, nachgefahren | Ergebnis |
|---|---|---|
| R1: `output: "null"` und `output: '~'` sind Text | V-R1a: Text `null` oder `~` auch in Anführungszeichen als `null` (der früher grüne Mutant) | **rot**, `TestDateiGueltig` |
| R1: `Null` und `NULL` ohne Anführungszeichen sind `PGR-E2004` | V-R1b: Prüfung auf `!!null` durch Textvergleich mit leer, `null` und `~` ersetzt | **rot**, `TestDateiUngueltig`, `TestConfigShowFehler` |
| R2: Link ins Leere wählt keine Datei, die erste Zeile sagt das | V-R2: `os.Lstat` statt `os.Stat` in `waehleDatei` | **rot**, `TestConfigShowOhneDatei` (`Standarddatei als Link auf ein fehlendes Ziel: "", … nicht vorhanden`) |
| Lesart 3: Die erste Zeile ist der Pfad aus `--config`, auch relativ | V-L3: `filepath.Abs` auf den Pfad | **rot**, `TestConfigShowOhneDatei` (`--config mit relativem Pfad: "/tmp/…/relativ.yaml…"`) |
| Lesart 2: Ein unbekanntes Unterkommando ist `PGR-E2001` | V-L2: jedes Unterkommando als `show` | **rot**, `TestConfigShowFehler` |

Die Belege in §7 stehen für alle drei Punkte in der Form *Zusage · Mutation · roter Test* (`AGENTS.md` §3.10). Meine Läufe decken sich mit den Zeilen dort.

### Punkt 4: `make gates` grün. **Bestätigt durch eigenen Lauf, Beleg in §7 vorhanden.**

Eigener Lauf im Arbeitsbaum am Stand `c5b7d96`: Exit 0. Darin `d-check: 382 Datei(en) geprüft, 0 Befund(e)`, `gesamt: 0 Befund(e)` (a-check), `a-check-negativ: gruen`, `run-integration-tests: gruen` mit `TestE2EConfigShow`, `TestE2EReplayKonfigurationsdateiUngueltig` und `TestE2EConfigShowStdin`, `abdeckung-gegenprobe: gruen`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `commit-msg-gegenprobe: gruen`, `baseline-verify: v6.16.0 OK`. Das deckt sich mit der Zeile *Läufe zur Nacharbeit* in §7 (382 Dateien).

### Übrige DoD-Punkte (konstant je Slice)

- **Review:** Der Report liegt vor (`ef8bee5`), aus einem anderen Lauf. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen:** noch offen. Sie gehören in die Closure nach dieser Verifikation. In §7 stehen sie als `<…>`; alle drei Risiken in §6 tragen „offen bis Closure“. Das ist kein Befund, sondern die Reihenfolge. In die Closure gehören die beiden Beobachtungen, die der Plan dort schon nennt: F-509 als Beleg zu `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (Review) und `BEO-REPO/slice-waechst-durch-uebernahmen` mit dann drei Belegen (§8).

---

## 2. Review-Befunde F-506 bis F-518 und Vorgaben 1 bis 8: nachgefahren

| Befund / Vorgabe | Prüfung | Urteil |
|---|---|---|
| F-506 Tag (Vorgabe 1) | Binary, je `PGR-E2004 … Tag ist ungültig`: `! info`, `!!str info`, `!<tag:yaml.org,2002:str> info`, `!` mit führendem BOM, `!` nach `\r`, `!!str` nach `\r\n`, `!x` am Abschnitt, `!!map` am Dokument, `!` am Schlüssel oben und am Namen einer Verbindung. Die gültige Datei aus §6 (`!` im Kommentar in der Spalte nach `\r`) läuft mit Exit 0. **V-TS** (Prüfung über `TaggedStyle` entfernt): **rot**, `TestDateiTag` (`TaggedStyle ohne Text: <nil>`). **V-BOMZ** (Zeilen aus dem Text mit BOM gezählt): **rot**, `TestDateiTag`. **V-CR** (nur nach `\n` gezählt): **rot**, `TestDateiTag` | erledigt |
| F-507 (a) Kodierung (Vorgabe 2) | Binary: UTF-16LE mit BOM, UTF-16BE mit BOM und UTF-32BE mit BOM ergeben je `PGR-E2004 … ungültiges YAML in Zeile 1`. UTF-16LE ohne BOM ergibt `PGR-E2004 … ungültiges YAML`, **ohne Zeile** (V-119). Ein BOM in einem Wert ergibt `… in Zeile 3`, wo es steht; `0xFF` in einem Kommentar `… in Zeile 2` ohne Byte der Datei. **V-UTF8** (Prüfung auf UTF-8 entfernt), **V-BOM2** (BOM nach dem ersten Zeichen zugelassen): beide **rot**, `TestDateiKodierung` | erledigt bis auf V-119 |
| F-507 (b) Zeilenende (Vorgabe 3) | Binary: `U+2028` in einem Wert in Anführungszeichen in Zeile 1 ergibt `… in Zeile 1`. **V-2028** (Prüfung der drei Zeichen entfernt): **rot**, `TestDateiKodierung` | erledigt |
| F-507 (c) keine reguläre Datei (Vorgabe 4) | Binary, je `PGR-E2004 … nicht lesbar: keine reguläre Datei` mit der Quelle und ohne Pfad: `--config` auf ein Verzeichnis, `/dev/null` und eine FIFO; `PGWIRE_RECORDER_CONFIG` auf ein Verzeichnis; die Standarddatei als Verzeichnis und als FIFO; `replay --config <FIFO>`. **Pipe auf `/dev/stdin`:** `sleep 8 \| config show --config /dev/stdin` endet nach 0 s mit Exit 2, ebenso `replay` und `PGWIRE_RECORDER_CONFIG=/dev/stdin`. Mit `< regulaer.yaml` zeigt `config show` den Inhalt, Exit 0. Ohne Leserecht: `… nicht lesbar: permission denied`, ohne Pfad. **V-REG** (Prüfung nur auf Verzeichnis): **rot**, `TestDateiKeineRegulaere`. Als E2E-Mutant **rot**, `TestE2EConfigShowStdin` (`endet nicht binnen 15 s`); ungeändert grün. **V-LSTAT** (`os.Lstat` in `ladeDatei`): **rot**, `TestDateiKeineRegulaere` | erledigt |
| F-507 (d), F-511 leerer und nicht skalarer Schlüssel (Vorgabe 5) | Binary: `"": 1` ergibt `oberste Ebene: unbekannter Schlüssel`; `? [GEHEIM]` in `record:` ergibt `record: unbekannter Schlüssel`. **V-F511** (Stelle wieder aus dem Text): **rot**, `TestDateiUngueltig` | erledigt |
| F-508, F-510 doppelte Schlüssel (Vorgabe 6) | Binary, je nur `ungültiges YAML in Zeile N, Schlüssel doppelt`, ohne `GEHEIM`: unter `connections.a` (Zeile 4), als Name einer Verbindung (Zeile 3), in einer Liste an einer Wertstelle (Zeile 3), in vierter Ebene (Zeile 5) und oben (Zeile 2). **V-F508** (Text des Schlüssels in der Meldung): **rot**, `TestDateiUngueltig`. **M1 als V-M1** (kein Abstieg in Listen): **rot**, `TestDateiUngueltig` (`doppelt in einer Liste … erhalten … erwartet ein Skalar`). **M2 als V-M2** (Doppelte nur bis zur zweiten Ebene): **rot**, `TestDateiUngueltig` (`doppelt an Wertstelle mit Wert`, `doppelt in dritter Ebene`) | erledigt |
| Ausgabe ohne BOM und mit `\n` (Vorgabe 7) | Binary mit `od -c`, siehe Punkt 1. `TestConfigShowOhneBOM` vergleicht die ganze Ausgabe. Ohne Mutation, wie §6 sagt | erledigt |
| Kommentare an `form`, `zeilen`, `kodierung`, `doppelte` (Vorgabe 8) | gelesen: Jeder sagt zu, was die Tests oben halten. Zum Kopfkommentar von `ladeDatei` („nennt die Stelle“) siehe V-119 | erledigt |
| F-509 | Beleg für die Closure (Register), siehe Abschnitt 1 | offen bis Closure, kein Befund hier |
| F-512, F-513 | **V-M3** (Steuerzeichen nur unter `0x20`): **rot**, `TestDateiUngueltig`, `TestConfigShowFehler` (DEL, C1). **V-M4** (`Contains` statt `HasPrefix`): **rot**, `TestConfigShow`. **V-ENVLEER** (leere Variablen mitgezählt): **rot**, `TestConfigShow` und drei weitere | erledigt |
| F-514, F-515 | §3 nennt die Tests als „geliefert in `9332356`“. §1 führt TLS über `slice-v1-abschluss-verbindungen-platzhalter`, der die Wirkung bei `play` weitergibt | erledigt |
| F-516 Zwischenrisiko Klartext-Passwort | Das Risiko trägt: kein Tag im Repo, kein Release dazwischen; der Nehmer liegt in `next/` und führt in seiner DoD `PGR-E2006` auch bei `config show`. Dass das Handbuch §5 `PGR-E2006` schon zusagt, deckt §3 als Zielstand. Zu §4 des Handbuchs siehe V-120, zu einer Lücke im Nehmer V-121 | trägt |
| F-517 Kodierer in `config show` | [ADR-0037](../plan/adr/0037-yaml-bibliothek-fuer-die-anzeige-der-konfigurationsdatei.md) ist `Accepted`. Eingabe des Kodierers ist nur `d.inhalt`, der Baum aus dem Lesen derselben Datei, verändert nur durch `ohneKommentare`. `ARC-013` und §2 der Sicht nennen „Lesen und Anzeigen“ | erledigt |
| F-518 Negativbefund | bestätigt, siehe Abschnitt 3 | bestätigt |

---

## 3. Keine Meldung zu Datei oder Umgebung nennt einen Wert (`647f0ff`)

Jede Meldung in `datei.go` habe ich gelesen und, wo möglich, mit `GEHEIM` an der Stelle eines Werts am Binary ausgelöst.

| Meldung | Inhalt | Probe |
|---|---|---|
| `nichtLesbar` | Quelle und Grund des Betriebssystems ohne Pfad (`pe.Err`) | `permission denied`, `too many levels of symbolic links`, `nicht vorhanden`, je ohne Pfad |
| `keine reguläre Datei` | Quelle | Verzeichnis, FIFO, `/dev/null`, Pipe |
| `fehlerZeile`, `ungueltigesYAML` | nur die Zeile; vom Text der Bibliothek nur die Zahl | unterminiertes `"GEHEIM`, Tab vor `GEHEIM`, `U+2028` neben `GEHEIM`, `0xFF` hinter `GEHEIM`: kein `GEHEIM` |
| `mehr als ein Dokument` | fester Text | zweites Dokument mit `log_level: GEHEIM`: kein `GEHEIM` |
| `doppelte` | nur die Zeile | siehe F-508 oben |
| `fehlerDatei` | Stelle (Schlüssel bzw. `connections` und Name) und ein fester Grund | `connections.a: [GEHEIM]`, Anker `&GEHEIM`, Alias `*GEHEIM`, `? [GEHEIM]`: kein `GEHEIM` |
| `setze` | Stelle und die Ursache aus `art.pruefe`; deren vier Texte in `cli.go` enthalten keinen Wert | `log_level`, `record.shutdown_timeout`, `record.force` mit `GEHEIM`: kein `GEHEIM` |
| `anzeige` (`PGR-E1000`) | hängt den Fehler des Kodierers an | mit einem Baum aus dem Lesen nicht ausgelöst; kein Befund |

Für die Umgebung prüft `config show` nur `PGWIRE_RECORDER_CONFIG`, und dessen Meldungen nennen die Quelle, nicht den Pfad. Die Meldungen des Lesers zu den übrigen Variablen stammen aus `slice-v1-abschluss-konfiguration` und nennen keinen Wert (dort verifiziert).

Eine Stelle nennt Text, den die Spezifikation als Stelle erlaubt, der aber ein Geheimnis tragen kann: der Name einer Verbindung (V-121).

---

## 4. Hexagon und Architektur-Gate

- `make a-check`: Exit 0, `gesamt: 0 Befund(e)`. Der Hinweis auf neun Dateien unter `test/integration/` ohne Schicht ist Bestand.
- Die YAML-Bibliothek importieren im CLI-Adapter nur `datei.go` und `export_test.go` (Test). `datei` und `inhalt` sind nicht exportiert; `Command` trägt nur `Anzeige string`. Kein Typ der Bibliothek verlässt den Adapter.
- Schicht-Abgrenzung: `git diff --stat f333934 c5b7d96` ändert im Kern nur `internal/hexagon/model/fehler.go`, drei Zeilen Konstanten. Driven-Adapter, PGWire-Adapter, `cmd/`, `go.mod` und `go.sum` sind unverändert; der Bootstrap bekommt drei Zeilen für `config show`. Das ist genau die Ausnahme aus §1.

---

## 5. Plan gegen Code

- **§3:** Die Tabelle nennt jede geänderte Code-Datei: `internal/adapters/driving/cli` mit Tests, `internal/bootstrap`, das Model, `test/integration`, das Handbuch und die Spezifikation. Nicht in §3 stehen die Abdeckungstabellen `docs/user/abdeckung-*.md`. Sie sind vom Werkzeug geschrieben und gehören zu jedem Slice mit Tests; kein Befund.
- **§1, Abgrenzung:** Kein Code für URL, `sslmode`, Platzhalter, `PGR-E2005` oder `PGR-E2006` (siehe Punkt 2). `output` und `force` sind Schlüssel von `record:` (Rückgabe 12). Ein Schlüssel einer Option, die der Stand nicht kennt, ist unbekannt (`record.format`, `play`).
- **§6, Randformen:** Jede Randform, die ich am Binary gefahren habe, folgt ihrem Ort in `LH-FA-17.a`. Abweichungen: die Zeile in der Meldung (V-118, V-119).
- **§7:** Die Belege decken alle drei Liefer-Punkte und den Gate-Lauf. Es fehlt kein Beleg, den die DoD verlangt.
- **Handbuch:** §5 *Konfigurationsdatei* stimmt mit dem Binary überein, soweit dieser Slice liefert. Der Rest ist der Zielstand, den §3 nennt. §4 weicht von der Spezifikation ab (V-120).

---

## 6. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-118 | MEDIUM | Die Zeile, die eine Meldung zu ungültigem YAML nennt, ist bei Fehlern des Parsers falsch. `LH-FA-17.a` *Fehler* sagt: „Zu ungültigem YAML nennt die Meldung die Zeile, gezählt ab 1“. `ungueltigesYAML` übernimmt die Zahl aus dem Text der Bibliothek. Bei Fehlern des Parsers ist das die Zeile, in der das umgebende Konstrukt beginnt, gezählt ab 0. Das zeigt die Rohausgabe der Bibliothek in einer Sonde: `d: [x` in Zeile 4 ergibt `yaml: line 3: did not find expected ',' or ']'`. Binary: offene Folge in Zeile 4 ergibt `Zeile 3`; falsche Einrückung in Zeile 5 ergibt `Zeile 4`; `  - y` in Zeile 5 unter `record:` ergibt `Zeile 2`; `]` zu viel in Zeile 3 ergibt `Zeile 2`; offene Abbildung in Zeile 5 ergibt `Zeile 4`. Unter keiner Lesart von „die Zeile“ stimmt die Zeile vor dem Konstrukt. Fehler des Scanners stimmen (Tab und Einrückung in Zeile 3: `Zeile 3`; `mapping values` in Zeile 2: `Zeile 2`). Kein Test prüft die Zahl: Der Fall *Syntax* in `ungueltigeDateien` erwartet nur „ungültiges YAML in Zeile“. *Failure-Szenario:* Wer eine fehlende `]` sucht, wird an die Zeile davor geschickt, und dort steht nichts Falsches. | `internal/adapters/driving/cli/datei.go` · `yamlZeile = regexp.MustCompile(` und `ungueltigesYAML`; `datei_test.go` · `{"Syntax", "replay:\n  listen: \"GEHEIM\n", "ungültiges YAML in Zeile"}` | Architect: welche Zeile bei einem Konstrukt über mehrere Zeilen gilt, Anfang oder Stelle der Erkennung (`AGENTS.md` §3.12). Danach Implementer: Test mit der Zahl und Mutation (§3.10) |
| V-119 | LOW | Drei Klassen ungültigen YAMLs ergeben eine Meldung **ohne** Zeile, weil der Text der Bibliothek keine trägt: (1) ein Steuerzeichen aus C0, auch NUL und in Anführungszeichen (`yaml: control characters are not allowed`); (2) daraus folgend UTF-16LE ohne BOM, das in ASCII-Text aus NUL-Bytes besteht; (3) ein Alias ohne Anker (`yaml: unknown anchor 'x' referenced`). `LH-FA-17.a` *Fehler* verlangt die Zeile ohne Ausnahme. §6 sagt für UTF-16 ausdrücklich „`PGR-E2004` mit der Zeile“, und der Kopfkommentar von `ladeDatei` sagt „nennt die Stelle“. Der Test zu UTF-16 prüft nur „ungültiges YAML“. Den Fall *Alias ohne Anker* hat der Implementer bewusst so getestet (§7, „ohne Text der Bibliothek“); dass dabei die Zeile fehlt, nennt weder §6 noch die Spezifikation. | `internal/adapters/driving/cli/datei.go` · `return model.Errorf(model.CodeConfigFile, nil, "Konfigurationsdatei: ungültiges YAML")`; `datei_test.go` · `"UTF-16LE ohne BOM": utf16(text, true, false),` | Architect: Randform „die Bibliothek nennt keine Zeile“ entscheiden, etwa C0 in `kodierung` mit Zeile ablehnen (§3.12). Danach Implementer: Kommentar und Test nachziehen (§3.11) |
| V-120 | LOW | Handbuch §4 sagt zu `config show`: „Ist die Datei ungültig, zeigt das Werkzeug nichts und meldet `PGR-E2004` bis `PGR-E2006`.“ `LH-FA-17.a` *Anzeige* sagt: „mit dem Fehlercode des Ladens (`PGR-E2004` oder `PGR-E2006`) …; `PGR-E2005` kommt nicht vor“. Das Handbuch nennt einen Code, den der Befehl auch im Zielstand nie meldet. Nach §3 gehört §4 des Handbuchs zu diesem Slice, nicht zum Nehmer. | `docs/user/benutzerhandbuch.md` §4 *Die gewählte Konfigurationsdatei anzeigen* · „meldet `PGR-E2004` bis `PGR-E2006`“ | Implementer |
| V-121 | INFO | Der Name einer Verbindung erscheint in der Meldung, und er kann ein Passwort tragen. Binary: `connections:` mit `postgresql://u:GEHEIM@h/db:` ohne Wert ergibt `PGR-E2004 … connections.postgresql://u:GEHEIM@h/db: ein leerer Wert und null sind ungültig`. Mit Wert `x` ist die Datei gültig, und `config show` zeigt den Namen. Beides ist spezifikationskonform: *Fehler* erlaubt den Verbindungsnamen als Stelle, und ein Name ist jeder Text ohne Steuerzeichen. Der Grund derselben Regel („die Geheimnisse tragen können“) trifft hier aber zu. Das Klartext-Passwort des Nehmers (`PGR-E2006`) ist als Passwortteil einer URL im *Wert* definiert; eine URL im Namen fängt es nicht. | `internal/adapters/driving/cli/datei.go` · `skalar(n.Content[i+1], unter("connections", name), z)` | Architect: Randform für §6 von `slice-v1-abschluss-verbindungen-platzhalter` (Name in der Form einer URL mit Passwort). Kein Befund gegen diesen Slice |

---

## 7. Gefahrene Mutanten

Alle am Stand `c5b7d96`, je in einer frischen Kopie, `go test -count=1` für `internal/adapters/driving/cli` und `internal/bootstrap`. Ausnahme ist V-REG-E2E.

| Mutant | Änderung | Ergebnis |
|---|---|---|
| V-GRUND | keine | grün |
| V-TS | `n.Style&yaml.TaggedStyle != 0 \|\|` aus `form` entfernt | rot: `TestDateiTag` |
| V-BOMZ | `zeilen(roh)` statt `zeilen(data)` (BOM vor dem Zählen nicht entfernt) | rot: `TestDateiTag` |
| V-CR | `zeilen` trennt nur nach `\r\n` und `\n` | rot: `TestDateiTag` |
| V-M1 | `doppelte` kehrt bei einer Liste sofort zurück | rot: `TestDateiUngueltig` |
| V-M2 | Doppelte nur bis zur zweiten Ebene, nur über Werte von Abbildungen | rot: `TestDateiUngueltig` |
| V-F508 | Text des Schlüssels in der Meldung zu Doppelten | rot: `TestDateiUngueltig` |
| V-F511 | `stelleVon` nennt immer den Text des Schlüssels | rot: `TestDateiUngueltig` |
| V-R1a | `null` und `~` auch in Anführungszeichen als `null` | rot: `TestDateiGueltig` |
| V-R1b | `!!null` durch Textvergleich mit leer, `null`, `~` ersetzt | rot: `TestDateiUngueltig`, `TestConfigShowFehler` |
| V-R2 | `os.Lstat` statt `os.Stat` in `waehleDatei` | rot: `TestConfigShowOhneDatei` |
| V-L3 | `filepath.Abs` auf den Pfad aus `--config` | rot: `TestConfigShowOhneDatei` |
| V-L2 | `parseConfig` nimmt jedes Unterkommando als `show` | rot: `TestConfigShowFehler` |
| V-REG | `info.IsDir()` statt `!info.Mode().IsRegular()` | rot: `TestDateiKeineRegulaere` |
| V-REG-E2E | dieselbe Mutation, Image der Stufe `integration` | rot: `TestE2EConfigShowStdin` (`endet nicht binnen 15 s`); ungeändert grün |
| V-LSTAT | `os.Lstat` statt `os.Stat` in `ladeDatei` | rot: `TestDateiKeineRegulaere` |
| V-UTF8 | Prüfung auf ungültiges UTF-8 entfernt | rot: `TestDateiKodierung` |
| V-BOM2 | BOM nach dem ersten Zeichen zugelassen | rot: `TestDateiKodierung` |
| V-2028 | Prüfung von `U+0085`, `U+2028`, `U+2029` entfernt | rot: `TestDateiKodierung` |
| V-M3 | Steuerzeichen im Namen nur unter `0x20` | rot: `TestDateiUngueltig`, `TestConfigShowFehler` |
| V-M4 | `strings.Contains` statt `strings.HasPrefix` | rot: `TestConfigShow` |
| V-ENVLEER | leere Variablen mitgezählt | rot: `TestConfigShow`, `TestConfigShowOhneBOM`, `TestConfigShowOhneDatei`, `TestRunConfigShow` |
| V-OBEN | `log_level` im Abschnitt statt oben | rot: `TestLeserAlleOptionen`, `TestDateiGueltig`, `TestDateiWahl`, `TestDateiTag`, `TestDateiKeineRegulaere`, `TestRunReplayAusDatei` |

Jeder rote Mutant scheitert aus dem richtigen Grund; die gelesenen Fehlerzeilen nennen die mutierte Zusage. Kein Mutant blieb grün. Für die Zahl in der Zeile (V-118) gibt es keinen Test, der rot werden könnte.

---

## Verdikt

**DoD-Liefer-Punkte 1 bis 3 und `make gates`: bestätigt.** Die Belege in §7 sind vollständig, und meine Läufe decken sich mit ihnen. Das Binary hält jede Stichprobe aus `LH-FA-17.a`: BOM und `\r`, UTF-16, `!` und `!!str`, Doppelte ohne Wert, `--config` auf ein Verzeichnis, FIFO und Pipe auf `/dev/stdin` ohne Warten, die Standarddatei als Link ins Leere und `config show` gegen das Handbuch. Keine Meldung zu Datei oder Umgebung nennt einen Wert. Das Hexagon bleibt rein, `make a-check` ist grün. Die Review-Befunde F-506 bis F-515 und F-517 sind erledigt, F-516 trägt.

**Vor der Closure:** V-118 und V-119 an den Architect, weil die Zeile bei mehrzeiligen Konstrukten und das Fehlen einer Zeile Randformen sind (`AGENTS.md` §3.12); danach an den Implementer. V-120 an den Implementer (Handbuch, ohne Code). V-121 an den Architect für §6 des Nehmers. Danach stehen die Closure-Pflichten in §7 aus.

**Summary:** 0 HIGH · 1 MEDIUM · 2 LOW · 1 INFO (V-118 die Zeile zu ungültigem YAML ist bei Parser-Fehlern die Zeile vor dem Konstrukt; V-119 C0-Steuerzeichen, UTF-16LE ohne BOM und ein Alias ohne Anker ergeben eine Meldung ohne Zeile; V-120 Handbuch §4 nennt `PGR-E2005` für `config show`; V-121 ein Verbindungsname in der Form einer URL trägt ein Passwort in Meldung und Anzeige, Randform für den Nehmer).
