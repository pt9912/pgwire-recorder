# Verifikation: slice-v1-abschluss-schreiben — 2026-10-09

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-09-review-slice-v1-abschluss-schreiben.md`, `5bff9bd`, F-537 bis F-546). F-546 hat er an mich übergeben.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-schreiben.md`, Kopf und §1 bis §8, gegen den Gesamt-Diff `8e24973..3a7018b`. Architect: `ad1b7ba`, `b74e16e`, `278d929`, `9996a8c`. Implementer: `aad7085`, `df4f591`, `65968e6`, `0311f72`, `24fc87c`, `9bd65c5`, `ca5124f`, `9086fd7`, `3a7018b`. Review: `5bff9bd`. Die Nacharbeit nach dem Review (`ca5124f`, `3a7018b`) hat kein eigenes Review; ihre Zusagen habe ich unten mit Mutanten nachgefahren.

**Abgrenzung:** Die Anmeldung von `--output` und `--force` am allgemeinen Leser und ihre Werte liegen bei `slice-v1-abschluss-konfiguration`, die Schlüssel der Datei bei `slice-v1-abschluss-konfigurationsdatei`, das Schreiben von `sqlite` bei `slice-v1-abschluss-sqlite-format` (§1). Dass sie hier fehlen, ist kein Befund.

**Eingang:**

- der Plan ganz, besonders §1, §2 (DoD), §3, §6 (*Randformen*, *Risiken*) und §7 (*Belege des Implementers*)
- `spec/spezifikation.md` [`LH-FA-07.a`](../../spec/spezifikation.md#lh-fa-07a--sicheres-schreiben-des-recordings) ganz (*Zielpfad beim Start*, *Temporäre Datei*, *Fehlermodi*), [`LH-FA-08.a`](../../spec/spezifikation.md#lh-fa-08a--auswahl-des-recordings), [`LH-FA-13.a`](../../spec/spezifikation.md#lh-fa-13a--signalbehandlung) (Frist, Zwangsende, Schreiben danach) und die Änderungsliste
- [`LH-FA-07`](../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-08`](../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [ADR-0005](../plan/adr/0005-recording-store-ist-driven-adapter.md)
- am Stand `3a7018b`: der Diff von `internal/adapters/driven/recording/yaml.go` ganz, `export_test.go`, `schreiben_test.go` (Testnamen und die Fehlerzeilen der Mutanten), `test/integration/schreiben_e2e_test.go` ganz, `startProzessIn` in `record_e2e_test.go`, die Zeile `force` in `internal/adapters/driving/cli/cli.go`
- `docs/user/benutzerhandbuch.md`, Diff des Slice
- der Nehmer `slice-v1-abschluss-sqlite-format` (§1, DoD), als Adresse zu F-538
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen

Beim Start war der Arbeitsbaum sauber, HEAD `3a7018b`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-09

**So wurden die Proben gebaut.** Alles lief unter dem Scratchpad der Sitzung, nie im Repo.

- **Binary:** eine frische Kopie per `git archive 3a7018b | tar -x`, ohne `.git`. Daraus ein Image der Stufe `runtime` mit eigenem Tag. Ein eigenes internes Netz mit dem gepinnten PostgreSQL-Image aus `integration.mk` als Upstream; der Recorder lief im Image `runtime` (Benutzer `nonroot`, für die Fälle ohne Schreibrecht `--user 65534`), das Arbeitsverzeichnis als Bind-Mount. Jeder Start lief mit `timeout 30`, jedes `docker wait` mit `timeout 30`. Für die umask-Fälle habe ich dasselbe Binary aus dem Image `runtime` kopiert und in der Shell des PostgreSQL-Images gestartet, weil `runtime` keine Shell hat. Als Client für das Zwangsende diente `psql` aus demselben Image.
- **Unit-Mutanten:** je Mutant eine frische Kopie von `go.mod`, `go.sum`, `cmd`, `internal` und `test` per `cp -r` ohne `-p`. Ein Skript ersetzt genau eine Stelle und bricht bei einer anderen Trefferzahl ab. Dann im Image der Stufe `deps` (eigener Tag), Bind-Mount, `--network=none`: `gofmt -l ./internal` (leer, sonst Abbruch), `go vet` und `go test -count=1` des Pakets `internal/adapters/driven/recording`. Die Fehlerzeilen jedes roten Tests habe ich gelesen.
- **E2E-Mutanten:** je Mutant eine frische `cp -r`-Kopie der Archiv-Kopie, `gofmt -l` leer, Image der Stufe `integration` mit eigenem Tag, eigenes internes Netz und eigener PostgreSQL-Container, gefahren nur der genannte Test (`-test.run`), Frist 300 s.
- **Danach:** eigene Container, Netze und Images (`pgr-ver-schreiben:*`) entfernt, die Kopien und Proben-Verzeichnisse gelöscht. `make gates` lief im Arbeitsbaum; `git status` war danach sauber.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: atomares Schreiben über eine temporäre Datei, keine teilweise Datei nach einem Abbruch, auch nach dem Zwangsende ([`LH-FA-07`](../../spec/lastenheft.md#lh-fa-07--persistente-recordings)). **Bestätigt.** Siehe V-127 (Plan-Text) und V-128 (Grenze 5).

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Schreiben in eine temporäre Datei, danach Verschieben; Zieldatei vollständig oder unverändert | Code: `schreibe` legt `.<Name>.<16 Hex>.tmp` mit `O_EXCL` an, schreibt, synchronisiert, schließt und verschiebt mit `os.Rename`. Mutant *Verschieben durch Entfernen ersetzt* als E2E: **rot**, `TestE2ERecordZwangsendeSchreibtVollstaendig` (Zieldatei bleibt `alt`) | bestätigt |
| Nach dem Zwangsende vollständig ([`LH-FA-13.a`](../../spec/spezifikation.md#lh-fa-13a--signalbehandlung)) | Binary: `--force` über einer vorhandenen Datei, `--shutdown-timeout 1s`, `psql` mit `select 42` und offenem `pg_sleep(60)`, dann `SIGTERM`: `PGR-E4006`, Exit 4, `record beendet`. Unter `--output` liegt eine vollständige Aufzeichnung (45 Zeilen, `select 42 as antwort`), `replay --input` startet sie ohne Fehler. Im Verzeichnis liegt nur `rec.yaml` | bestätigt |
| Ein Abbruch während des Schreibens lässt keine teilweise Datei | `TestWriteFehlschlag` (Schreiben nach einem Teil der Daten) hält es über `Eingriffe`; ein harter Abbruch (`SIGKILL`) ist §7 *Grenzen* 3, die Spezifikation erlaubt dann eine liegen gebliebene temporäre Datei. Das Verschieben in einem Schritt trägt die Zusage; über einem Dateisystem unter Linux geprüft, das Risiko in §6 bleibt | bestätigt, mit Grenze 3 |
| Zielpfad beim Start: keine reguläre Datei `PGR-E3001`, auch mit `--force` | Binary: Verzeichnis an `--output`, ohne und mit `--force`: `Recording [PGR-E3001]: /w/rec.yaml ist keine reguläre Datei`, Exit 3, Verzeichnis leer. Mutant *`IsRegular` entfernt*: unit **rot** `TestPrepareKeineRegulaereDatei`; als E2E **rot** `TestE2ERecordZielIstVerzeichnis` (Exit 2, `PGR-E2002`) | bestätigt |
| Fehlendes Verzeichnis `PGR-E3001`, kein Verzeichnis angelegt | Binary: `--output /w/fehlt/rec.yaml`, ohne und mit `--force`: `PGR-E3001 … nicht beschreibbar: … no such file or directory`, Exit 3; `fehlt` wurde nicht angelegt | bestätigt |
| Nicht beschreibbares Verzeichnis `PGR-E3001` beim Start, auch mit vorhandener Datei und `--force` (F-537) | Binary mit `--user 65534`, Verzeichnis `0555`: Datei fehlt; vorhandene Datei mit `--force`; Verknüpfung auf sie mit `--force` — je sofort `PGR-E3001 … nicht beschreibbar: … permission denied`, Exit 3, Datei bleibt `alt`. Ohne `--force` über der vorhandenen Datei geht `PGR-E2002` vor (Exit 2), wie die Reihenfolge in *Zielpfad beim Start*. Mutant D (Probedatei nur für fehlenden Pfad): **rot**, `TestPrepareVerzeichnisNichtBeschreibbar` (`vorhanden mit --force`, `Verknüpfung mit --force`) und `TestPrepareProbedateiScheitert` (`vorhanden=true`) | bestätigt |
| Verknüpfung: vorhanden nach ihrem Ziel, ins Leere nicht vorhanden, ersetzt wird die Verknüpfung | Binary: Verknüpfung auf eine Datei ohne `--force`: `PGR-E2002`, Exit 2. Verknüpfung ins Leere ohne `--force`: Start, `SIGTERM`, Exit 0, an ihrer Stelle eine reguläre Datei `sessions: []`. Verknüpfung auf `ziel.yaml` (`0640`) mit `--force`: `rec.yaml` ist danach eine reguläre Datei mit `0640`, `ziel.yaml` bleibt `alt`. Mutant `os.Lstat` in `pruefe`: **rot**, `TestPrepareVerknuepfung`, `TestPrepareNichtPruefbar`, dazu `TestPrepareVerzeichnisNichtBeschreibbar` (die Verknüpfung gilt als keine reguläre Datei) | bestätigt |
| Rechte: neu `0666` nach der umask, ersetzt mit ihren Zugriffsrechten, bei Verknüpfung die des Ziels | Binary: neue Datei unter umask `022` `0644`, unter `077` `0600`, unter `027` `0640`; ersetzte Datei mit `0600` behält `0600`; Verknüpfung siehe oben. Mutant `0o600` statt `0o666`: **rot**, `TestWriteRechteUnterUmask022`. Mutant `os.Lstat` in `schreibe`: **rot**, `TestWriteRechteVerknuepfung` | bestätigt |
| Temporäre Datei: zehn Versuche, nur ein belegter Name führt zu einem neuen (F-544) | Mutant neun statt zehn Versuche: **rot**, `TestWriteZehnBelegteNamen`. Mutant C (`err != nil` statt `fs.ErrExist`): **rot**, `TestWriteAndererFehlerBeimAnlegen` (beide Fälle: zehn Züge, Ursache des Betriebssystems fehlt) | bestätigt |
| Fehlschlag: temporäre Datei entfernt, Fehler des Entfernens als Ursache derselben Meldung | Mutant *Entfernen im Zweig des Verschiebens weggelassen*: **rot**, `TestWriteFehlschlag` (`.tmp` liegt), `TestWriteEntfernenScheitert` (Ursache fehlt) | bestätigt |
| Probedatei: Schließen oder Entfernen scheitert `PGR-E3001` mit Ursache; nach gescheitertem Schließen dennoch entfernt, eine Meldung (F-543, Randform-Rückgabe) | Mutant *nach gescheitertem Schließen nicht entfernt*: **rot**, `TestPrepareProbedateiScheitert` (`.probe` liegt; Ursache des Entfernens fehlt). Mutant *Ursache des Entfernens `nil`*: **rot**, dort `Entfernen`. Mutant *Ursache des Schließens verworfen*: **rot**, dort `Schliessen`, `Schliessen und Entfernen` | bestätigt |
| Keine `.tmp`- oder `.probe`-Reste | Binary: `find` nach `*.tmp` und `*.probe` über alle Proben-Verzeichnisse: keine | bestätigt |

### Punkt 2: vorhandenes `--output` ohne `--force` abgelehnt, Exit 2, gleich aus welcher Quelle `--force` kommt ([`LH-FA-08`](../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Ohne `--force` und mit `--force=false` `PGR-E2002`, Exit 2, Datei unverändert | Binary: beide Fälle `Konfiguration [PGR-E2002]: /w/rec.yaml existiert bereits; --force ersetzt die Datei`, Exit 2, Inhalt `alt` | bestätigt |
| Mit `--force` ersetzt | Binary: `--force` über `alt.yaml`, `SIGTERM`, Exit 0, Datei ersetzt | bestätigt |
| `true` und `false` aus Kommandozeile, Umgebung und Datei | `TestE2ERecordVorhandeneZieldatei` im eigenen `make gates`, sieben Fälle PASS. E2E-Mutant *`replace` übergangen* (`case err == nil:`): **rot** genau in `Kommandozeile/true`, `Umgebung/true`, `Datei/true` (Start scheitert, der Recorder lauscht nicht), die vier `false`-Fälle grün. E2E-Mutant *Standardwert von `force` `true`*: **rot** genau in `keine/false` (endet nicht binnen 15 s), die übrigen sechs grün | bestätigt |

Die Belege in §7 stehen für beide Punkte in der Form *Zusage · Mutation · roter Test* (`AGENTS.md` §3.10). Meine 11 Unit- und 4 E2E-Mutanten decken sich mit den Zeilen dort; jeder ist aus dem genannten Grund rot, keiner aus `gofmt` oder `vet`. Abweichung nur nach oben: `os.Lstat` in `pruefe` ist zusätzlich in `TestPrepareVerzeichnisNichtBeschreibbar` rot, was §7 nicht nennt; das ist kein Befund.

### Punkt 3: `make gates` grün. **Bestätigt durch eigenen Lauf, Beleg in §7 vorhanden.**

Eigener Lauf im Arbeitsbaum am Stand `3a7018b`: Exit 0. Darin `d-check: 408 Datei(en) geprüft, 0 Befund(e)`, `gesamt: 0 Befund(e)` (a-check), `a-check-negativ: gruen`, `baseline-verify: v6.16.0 OK`, `run-integration-tests: gruen` mit `TestE2ERecordVorhandeneZieldatei` (sieben Fälle), `…ZielIstVerzeichnis`, `…ZwangsendeSchreibtVollstaendig` und `…SchreibfehlerNachZwangsende` PASS, `abdeckung-gegenprobe`, `kopf-check-gegenprobe`, `commit-msg-gegenprobe` und `lint-gegenprobe` je `gruen`. Das deckt sich mit der letzten Zeile *Läufe* in §7 (`make gates`, Commit zur Randform-Rückgabe, Exit 0).

### Übrige DoD-Punkte (konstant je Slice)

- **Review:** Der Report liegt vor (`5bff9bd`), aus einem anderen Lauf. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen:** noch offen; in §7 stehen sie als `<…>`, das Risiko in §6 trägt „offen bis Closure“. Das ist die Reihenfolge, kein Befund. In die Closure gehören F-539 als Beleg zu `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (Übergabe des Reviews) und die wiederkehrenden Klassen der Summary des Reviews. Das Risiko *Atomarität plattformabhängig* ist nicht eingetreten und nicht entfallen: Geprüft ist nur Linux mit einem Dateisystem (§7 *Grenzen* 2, meine Proben ebenso).

---

## 2. Review-Befunde F-537 bis F-546: nachgefahren

| Befund | Prüfung | Urteil |
|---|---|---|
| F-537 | Mutant D am Stand `3a7018b` **rot** (siehe Punkt 1); am Binary mit `--user 65534` für vorhandene Datei und Verknüpfung mit `--force` sofort `PGR-E3001`, Exit 3, ohne `Herunterfahren begonnen` | erledigt |
| F-538 | `slice-v1-abschluss-sqlite-format` §1 führt *Übernommen von `slice-v1-abschluss-schreiben`* (`278d929`), Kopf mit `LH-FA-07.a`; §1 *Ausdrücklich NICHT* schließt die Sendung nicht aus; der Nehmer liegt in `next/` (`AGENTS.md` §3.13). §7 nennt die Adresse richtig | erledigt |
| F-539 | Historie, keine Änderung am Code verlangt | offen bis Closure (Register), kein Befund hier |
| F-540 | §1 *Schicht-Abgrenzung* nennt jetzt nur den Recording-Adapter; `git diff --stat 8e24973 3a7018b` ändert an Code nur `yaml.go` (dazu Tests und `export_test.go`) | erledigt |
| F-541 | §7 *Randformen*: Die Code-Commits `24fc87c`, `ca5124f` und `3a7018b` ändern in §6 nur Testverweise bzw. den Stand schon entschiedener Randformen; `aad7085` und `65968e6` ändern nur §3. Der Satz stimmt mit dem Diff | erledigt |
| F-542 | Kommentar an `Write` zeigt auf `LH-FA-07.a` *Rechte* | erledigt |
| F-543 | `LH-FA-07.a` *Verzeichnis* entscheidet Schließen und Entfernen (`278d929`, `9996a8c`); Mutant P und die Mutanten an `pruefe` **rot** (Punkt 1) | erledigt, siehe V-127 zum Plan-Text |
| F-544 | `LH-FA-07.a` *Name* entscheidet (`278d929`); Mutant C **rot** | erledigt |
| F-545 | *Fehlermodi* nennen jeden Fehlschlag aus *Zielpfad beim Start* und *Temporäre Datei* | erledigt |
| F-546 (an mich) | Zwei Liefer-Punkte; Code nur im Recording-Adapter, eine Schicht (dazu Tests); der Diff ist in einer Sitzung zu prüfen | bestätigt |

---

## 3. Hexagon und Architektur-Gate ([ADR-0005](../plan/adr/0005-recording-store-ist-driven-adapter.md))

- Die Dateioperationen leben im Driven-Adapter `internal/adapters/driven/recording`; Kern, PGWire-Adapter, CLI-Adapter, Bootstrap und `cmd/` liegen nicht im Diff. `make a-check` im eigenen `make gates`: `gesamt: 0 Befund(e)`.
- `Eingriffe`, `PruefeMit` und `SchreibeMit` stehen nur in `export_test.go` und rufen dieselben `pruefe` und `schreibe` wie `Prepare` und `Write`.

---

## 4. Plan gegen Code

- **§1:** folgt dem Diff (F-540 erledigt). *Abgegeben* und *Ausdrücklich NICHT* halten: kein Code am Leser, an der Datei oder an `sqlite`.
- **§3:** nennt jede geänderte Datei; die Spezifikation ändern nur Commits des Architect, wie die Zeile sagt.
- **§6:** Jede Randform, die ich am Binary gefahren habe, folgt ihrem Ort in `LH-FA-07.a`. Eine Zeile beschreibt den Code vor `3a7018b` (V-127).
- **§7:** Die Belege decken beide Liefer-Punkte und den Gate-Lauf; es fehlt kein Beleg, den die DoD verlangt. Zwei Zeilen der Tabelle zu Punkt 1 nennen Mutationen an Funktionen, die `3a7018b` entfernt hat (V-127); die Zusagen dahinter halten am heutigen Code.
- **Handbuch:** Jede neue Aussage stimmt mit dem Binary überein (Quellen von `--force`, Verzeichnis als Ziel, fehlendes und nicht beschreibbares Verzeichnis auch mit `--force`, Verknüpfung, Name der temporären Datei, Rechte unter umask `022`).

---

## 5. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-127 | LOW | `3a7018b` hat `closeErr` und `removeErr` aus `yaml.go` entfernt (§7 sagt das selbst: „`closeErr` und `removeErr` entfallen“). §6 sagt zur Randform *Schließen oder Entfernen der Probedatei scheitert* aber weiter „Der Code folgt (`schliessen` aus der Tabelle, `closeErr`, `removeErr`)“. Zwei Zeilen der Tabelle zu DoD-Punkt 1 beschreiben Mutationen an diesen Funktionen („Ursache `nil` in `closeErr`“, „Ursache `nil` in `removeErr`“), die am heutigen Code nicht mehr gefahren werden können. Die Zusagen halten trotzdem: Ursache des Entfernens `nil` und Ursache des Schließens verworfen sind am Stand `3a7018b` **rot** in `TestPrepareProbedateiScheitert`. *Failure-Szenario:* Wer eine Mutation aus §7 wiederholt, findet die Stelle nicht und kann den Beleg nicht nachfahren. `AGENTS.md` §3.9, `BEO-REPO/plan-folgt-korrektur-nicht`. | Plan §6 · „Der Code folgt (`schliessen` aus der Tabelle, `closeErr`, `removeErr`)“; Plan §7 · „Ursache `nil` in `closeErr`“, „Ursache `nil` in `removeErr`“ | Implementer, nur Plan, vor der Closure |
| V-128 | INFO | §7 *Grenzen* 5: Ein nicht beschreibbares Verzeichnis ist nur über `Eingriffe` getestet, nicht als E2E. Am Binary habe ich es mit `--user 65534` für drei Formen gezeigt (Datei fehlt, vorhandene Datei mit `--force`, Verknüpfung mit `--force`), je sofort `PGR-E3001`, Exit 3. Die Grenze steht offen im Plan, und der Unit-Test hält die Mutation; das verletzt die DoD nicht. | Plan §7 *Grenzen* 5 | kein Auftrag; Hinweis für die Closure |

---

## Verdikt

**DoD-Liefer-Punkte 1 und 2 und `make gates`: bestätigt.** Die Belege in §7 sind vollständig, und meine Läufe decken sich mit ihnen. Das Binary aus `git archive 3a7018b` (Stufe `runtime`) hält jede Stichprobe aus `LH-FA-07.a`: vorhandene Datei ohne und mit `--force`, Verzeichnis als Ziel auch mit `--force`, fehlendes und nicht beschreibbares Verzeichnis auch mit vorhandener Datei und `--force`, Verknüpfung auf eine Datei und ins Leere, Rechte unter drei umasks, eine vollständige Aufzeichnung nach dem Zwangsende, keine Reste. 15 Mutanten sind aus dem richtigen Grund rot, keiner blieb grün. Das Hexagon bleibt rein. Die Review-Befunde F-537, F-538 und F-540 bis F-545 sind erledigt, F-539 gehört in die Closure.

**Vor der Closure:** V-127 an den Implementer (§6 und zwei Zeilen in §7 auf den Code nach `3a7018b`, ohne Code). Danach stehen die Closure-Pflichten in §7 aus, mit F-539 im Register und einem Ausgang für das Risiko der Atomarität.

**Summary:** 0 HIGH · 0 MEDIUM · 1 LOW · 1 INFO (V-127: §6 und §7 nennen `closeErr` und `removeErr`, die `3a7018b` entfernt hat; V-128: das nicht beschreibbare Verzeichnis ist am Binary bestätigt, ein E2E-Test fehlt weiter als benannte Grenze).
