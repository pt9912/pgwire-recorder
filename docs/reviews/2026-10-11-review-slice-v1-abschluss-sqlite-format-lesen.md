# Review-Report: slice-v1-abschluss-sqlite-format-lesen — 2026-10-11

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3, §6, §7 und §8, gegen die Spezifikation `LH-FA-22.a` (*Wahl und Erkennung*, *Lesen von `sqlite`*, *Grenzen*), `SPEC-043`, gegen [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0013](../plan/adr/0013-sqlite-recording-backend.md), [ADR-0024](../plan/adr/0024-sqlite-bibliothek.md), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md) und gegen `AGENTS.md` §3.

**Gegenstand:** `97f3ca6..a6bcc45` (sechs Commits: Lesen `45af81d`, Belege `60a448f`, Randformen `ea1ef40`, Surrogat/Meldungen/`ctx` `18695d4`, Belege `aee68b3`, Speicher-Entscheidung `a6bcc45`; +2455 −46 über 20 Dateien). Neu: `internal/adapters/driven/recording/` (`leser.go`, `sqlite_lesen.go`, `sqlite_lesen_baum.go`, `sqlite_json_lesen.go`, drei Testdateien) und `test/integration/sqlite_lesen_e2e_test.go`. Geändert: `sqlite.go` (`uriVon`, Platzhalter `Load` entfällt), `export_test.go`, `internal/bootstrap/bootstrap.go`, zwei Integrationstests, `spec/spezifikation.md`, der Slice-Plan, `docs/user/abdeckung-*.md`.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `slice-v1-abschluss-sqlite-format-lesen` ganz gelesen am Stand `a6bcc45`, mit §7 *Belege des Implementers* und *Nacharbeit*. Den Bericht des Implementers habe ich nicht gelesen.
- `spec/spezifikation.md` `LH-FA-22.a` (*Wahl und Erkennung*, *Lesen von `sqlite`*, *Grenzen*), `SPEC-043` und die Änderungen des Diffs daran.
- Nehmer `slice-v1-abschluss-sqlite-format-doku` (liegt in `next/`; nennt `slice-v1-abschluss-sqlite-format-lesen` in Kopf, §1 und DoD).
- `AGENTS.md` §3.3 bis §3.13; vorheriger Report am selben Modul: `slice-v1-abschluss-sqlite-format` (F-614 bis F-622), die Nummern dieses Laufs beginnen bei F-623.

**Ausgeführte Läufe** (Docker-only, Scratch unter dem Sitzungsverzeichnis, nichts im Repo geändert):

- **Arbeitsweise der Mutanten.** Frische Kopie (`cp -r` ohne `-p`: `go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`), genau eine Ersetzung, `gofmt -l` leer, `go test -count=1` des Pakets `recording` im Image `pgwire-recorder:test` mit Bind-Mount der Kopie (`--network=none`).
  - **Rot:** SG1 (Prüfung des einzelnen Surrogats entfällt: `TestLesenSurrogate`); ME1 (Ursache der YAML-Bibliothek durchgereicht: `TestLesenMeldungOhneWertUndZeile`); CX1 (Abbruch von `ctx` als `PGR-E3003`: `TestLesenNichtLesbarUndKopfOhneInhalt`); `SQLITE_PROTOCOL` aus der Stufung `PGR-E3001` entfernt (`TestLesenStufeNachErgebniscode`).
  - **Grün, wie in §7 als äquivalent geführt, bestätigt:** I9 (Art von `group_no` ungeprüft: Wert 0 lehnen die Regeln der Gruppen an jeder Stelle ab); E11 (jeder Fehler beim Lesen der 16 Bytes gilt als kurze Datei: `YAML.Load` meldet danach denselben `PGR-E3001`); J13 (Strings ohne Anführungszeichen-Stil: eine Aufzeichnung mit 44 Startup-Schlüsseln und -Werten wie `~`, `null`, `yes`, `0x10`, `1e3`, `.inf`, `- a`, `a: b`, `# c`, `&a`, `*a`, `!t`, Zeilenumbruch, `U+0085`, `U+2028`, NUL, `base64:` läuft mit dem Mutanten unverändert hin und zurück). E12 (Schließen der gelesenen Datei) habe ich nicht gefahren; ein Aufbau, in dem es scheitert, sehe ich nicht.
  - **Grün, nicht äquivalent:** O5 (`filepath.Abs` entfällt) — siehe F-624. U7 (`rows.Err()` in `zeilen` entfällt) — siehe F-625.
- **Probe `//`-Pfad** (Testprogramm in der Kopie): `Leser.Load("/" + <absoluter Pfad>)` und Pfade mit `./` und `../` lesen mit dem Code des Diffs; mit dem Mutanten O5 scheitert der Pfad `//tmp/…` mit `PGR-E3003 … invalid uri authority: tmp`.
- **Probe Meldungen** (Testprogramm in der Kopie, fünf Verstöße): Formatkennung, Richtung und Zahl `1.5` liefern eine Meldung ohne Pfad der Datei; ein unbekannter Schlüssel in `fields` liefert Pfad und „Inhalt der SQLite-Datei beschädigt“ ohne Art; die unbekannte Version nennt den Pfad nicht (F-623).
- **Probe Seitenbruch:** je Seite 0xFF über eine Datei mit 60 Sessions, Original und Mutant U7: identische Ergebnisse (die Beschädigung fällt vor oder in `Scan`, nicht in `rows.Err()`); den Zweig habe ich nicht erreicht.
- **Zählung.** `git diff --numstat 97f3ca6..HEAD -- '*.go'`: 863 Zeilen Code, 1375 Zeilen Tests (stimmt mit §7 *Nacharbeit*), 2455 Zeilen hinzu insgesamt; Commit-Messages: keine `SPEC-`/`ARC-`-Kennung.
- **Gates.** `make gates` auf dem sauberen Baum nach dem Commit dieses Reports: Ergebnis im Bericht der Übergabe.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-623 | MEDIUM | `LH-FA-22.a` *Lesen von `sqlite`* sagt zu: Die Meldung zu einem Fehler des Inhalts „nennt die Datei und die Art des Fehlers“. Die eigenen Prüfungen des Lesers (Formatkennung, `meta`, Art, Richtung, `fields`, JSON, Zahl, unbekannte Version `PGR-E3002`) nennen die Art, aber nicht den Pfad der Datei; der Weg über `Unmarshal` nennt den Pfad, aber nur „Inhalt der SQLite-Datei beschädigt“, nicht die Art (unbekannter Schlüssel, Zählung, Gruppe). Szenario: `replay --input` mit mehreren Aufzeichnungen im Skript meldet `PGR-E3003: meta: Formatkennung ist nicht pgwire-recorder` ohne zu sagen, welche Datei; die Dienste ergänzen den Pfad nicht. `TestLesenMeldungOhneWertUndZeile` prüft die Abwesenheit von Wert und Zeile, nicht die Anwesenheit von Datei und Art. | `LH-FA-22.a` *Lesen von `sqlite`*; Plan §3 (Auftrag „der Leser nennt Datei und Art“) | `internal/adapters/driven/recording/sqlite_lesen.go` · `return false, beschaedigt("meta ohne %s", metaFormat)` (Z. 216); `internal/adapters/driven/recording/sqlite_lesen.go` · `Inhalt der SQLite-Datei beschädigt` (Z. 45) | ja (Test auf Pfad in jeder Meldung des Inhalts; Mutation: Pfad entfällt) | Zusage der Meldung weiter als Code und Test |
| F-624 | MEDIUM | §7 führt O5 (`filepath.Abs` entfällt) als äquivalent; das trifft nicht zu: Ein Eingabepfad mit führendem `//` (unter Linux derselbe Pfad wie mit einem `/`) wird ohne `Abs` als `file://tmp/…` mit Authority gelesen und mit `PGR-E3003` abgelehnt, mit `Abs` gelesen. `LH-FA-22.a` sagt zu, der Pfad gelte wörtlich; `TestLesenSonderzeichenImPfad` fährt `?`, `#`, `%`, Leerzeichen und einen relativen Pfad, nicht diese Form. Szenario: `Abs` wird bei einer Änderung von `uriVon` entfernt, kein Test wird rot, und `replay --input //srv/rec.sqlite` bricht. | `LH-FA-22.a` *Wahl und Erkennung*; `AGENTS.md` §3.10 | `internal/adapters/driven/recording/sqlite.go` · `abs, err := filepath.Abs(path)` (Z. 260); Plan §7 · `(O5) \`filepath.Abs\` entfällt` | ja (Fall `//`-Pfad in `TestLesenSonderzeichenImPfad`; Mutant O5 wird rot) | Mutation grün, als äquivalent geführt, nicht äquivalent |
| F-625 | LOW | Die Spezifikation sagt zu, der Abbruch des Laufs „während des Lesens“ sei `PGR-E3001`; der Test bricht `ctx` vor dem Aufruf ab, `BeginTx` scheitert. Der Mutant, der `rows.Err()` in `zeilen` entfernt, bleibt grün; ein Abbruch oder Fehler mitten in einer Tabelle würde dann als vollständig gelesene Tabelle durchgehen. §7 führt den Mutanten nicht. Die Wirkung ist klein (der Lauf endet mit dem Abbruch), der Zweig ist ungeprüft. | `AGENTS.md` §3.11; `LH-FA-22.a` *Lesen von `sqlite`* | `internal/adapters/driven/recording/sqlite_lesen.go` · `if err := errors.Join(rows.Err(), rows.Close()); err != nil {` (Z. 176) | ja (Mutant U7) | Zusage im Vertrag weiter als der Test |
| F-626 | LOW | Zwei Kommentare beschreiben nicht, was da ist. `SQLite.Load`: „Schlüssel, Formen und Regeln der Aufzeichnung prüft nur dieser eine Weg“ — der Leser prüft Art, Richtung, `group_no`, Reihenfolge der Gruppen, Positionen und verwaiste Zeilen selbst (`interaktionVon`, `nachrichtVon`, `gruppen`, `liesNachrichten`). `export_test.go`, `JSONKnoten`: liefert „`""` und `""` für einen Knoten, der kein Skalar ist“ — der Code liefert `n.Tag` und `n.Value`, also `!!map` und `""` für ein Objekt. | `AGENTS.md` §3.7, §3.11 | `internal/adapters/driven/recording/sqlite_lesen.go` · `prüft nur dieser eine Weg` (Z. 21); `internal/adapters/driven/recording/export_test.go` · `oder "" und "" für einen Knoten, der kein Skalar ist` (Z. 141) | nein | Kommentar sagt mehr zu als der Code |
| F-627 | LOW | Der Plan folgt dem Code an einer Stelle nicht: §8 nennt als „Gemessen“ 2137 Zeilen hinzu, davon 2090 Go (799/1291), und „Liefer-Punkte und Schichten wie gezählt“; nach `18695d4` sind es 863 Code und 1375 Tests (2238 Go), wie §7 *Nacharbeit* selbst führt. | `AGENTS.md` §3.9; `BEO-REPO/plan-folgt-korrektur-nicht` | Plan §8 · `Gemessen: 2137 Zeilen hinzu, davon 2090 Go (799 Code, 1291 Tests)` | nein | Plan folgt der Korrektur nicht |
| F-628 | INFO | Größe. Der Diff ist in einer Sitzung prüfbar gewesen: 863 Zeilen Code in einer Schicht (Recording-Adapter, der Bootstrap tauscht einen Wert), davon `sqlite_lesen.go` 493 und `sqlite_json_lesen.go` 185 Zeilen, die ich ganz gelesen habe; die 1375 Zeilen Tests sind überwiegend Tabellen mit rohem SQL und kurzen Erwartungen, von mir stichprobenartig Zeile für Zeile gelesen (Meldungs-, Surrogat-, Stufungs-, Sperrfrist- und Binary-Tests), der Rest über die Deklarationen und die nachgefahrenen Mutanten geprüft. Die Schwelle aus Plan §4 („rund 2000“, gemessen 2238 Go) ist formal überschritten; die Rückführung ist eine Entscheidung des Planner, in §7 steht nur die Empfehlung des Architect (hinnehmen), kein Eintrag der Entscheidung. Der in §4 genannte Schnitt (DoD-Punkt 3) trennt `stufeCode` und seine Tests vom Leser, an dem sie in `bibliotheksfehler` hängen. Ein Fehler, den ein kleinerer Schnitt gefangen hätte, ist in diesem Lauf nicht aufgefallen; F-624 und F-625 sind Lücken der Mutationsliste, nicht der Größe. | Modul 5 (Slice in einer Review-Sitzung prüfbar); Plan §4 | Plan §7 · `**Größe:** \`git diff --numstat 97f3ca6..45af81d\` zählt 2137 Zeilen hinzu` | nein | Größe über der Schätzung |
| F-629 | INFO | `uriVon` trägt weiterhin den Zweig `runtime.GOOS == "windows"` (`"/" + abs`); er gilt jetzt für Lesen und Schreiben, hat keinen Test und keine Zusage. Plan §6 führt das Risiko (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, offen bis Closure). Zuständig: Planner bei der Closure. | Maintainability | `internal/adapters/driven/recording/sqlite.go` · `abs = "/" + abs` | nein | Plattformzweig ohne Prüfung |
| F-630 | INFO | Außerhalb des Diffs: `SQLite.Write` lehnt einen Schlüssel `<<` im Startup mit `PGR-E1000` („Tag !!merge ist kein String“) ab, während der Leser denselben Schlüssel aus einer fremden Datei als String liest. Das Verhalten des Schreibens liegt in `slice-v1-abschluss-sqlite-format` (`done`); ob ein solcher Parametername vorkommt, ist offen. Zuständig: Architect. | Maintainability | `internal/adapters/driven/recording/sqlite_json.go` · `ist kein String` | ja (Aufruf in der Kopie) | Randform des Schreibens ohne Entscheidung |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driven/recording/leser.go`: Erkennung an den ersten 16 Bytes, kürzere und leere Datei als `yaml`, Fehler beim Öffnen oder Lesen `PGR-E3001` mit Ursache, Pfad wörtlich, `Prepare`/`Write` `PGR-E1000` ohne Pfad; E1 bis E10 und K1/K2 laut §7, E11 bestätigt äquivalent | geprüft, ohne Befund außer F-624 (Pfad) |
| `sqlite_lesen.go`: `mode=ro`, Sperrfrist fünf Sekunden, eine Lesetransaktion (`SetMaxOpenConns(1)`, Rollback, Schließen nach dem Lesen ohne Meldung), Reihenfolge Formatkennung → Version → Inhalt, `meta.version` kanonisch (`-0`, führende Null, `+1` `PGR-E3003`, jede andere kanonische Zahl `PGR-E3002`), `empty_sessions` nur `true`, doppelter Schlüssel in `meta`, Typ-Prüfung je Spalte, verwaiste Zeilen, Sortierung nach Nummern, `{}` als fehlendes Feld, Fremdschlüssel nicht erzwungen, Zusätze ignoriert | geprüft, ohne Befund außer F-623, F-625, F-626 |
| `stufeCode`, `bibliotheksfehler`: zehn primäre Codes `PGR-E3001` über `& 0xff`, `IOERR_SHORT_READ` `PGR-E3003`, Fehler ohne Typ der Bibliothek (Abbruch von `ctx`) `PGR-E3001`; Mutanten PROTOCOL und CX1 rot | geprüft, ohne Befund |
| `sqlite_json_lesen.go`: JSON-Zahlen nicht über `float64` (`UseNumber`, kanonisch, 64 Bit), Umkehr von `\` und `base64:` (Standard-Alphabet mit Auffüllung), Bytes ohne UTF-8 als `!!binary`, einzelnes Surrogat (`ohneEinzelnesSurrogat`: Paar, `�`, wörtliches U+FFFD gültig, `\\ud800` Text; SG1 rot), Index-Arithmetik bei Escapes am Textende nicht erreichbar (jeder Text endet vor dem Ende mit `"`) | geprüft, ohne Befund |
| `sqlite_lesen_baum.go`: `group_no` ab 1 lückenlos, Client vor Server je Gruppe, `query` nur Server-Nachrichten ohne `group_no`, Nullzeiger bei Gruppe 0 durch `aktuell != 0` verhindert | geprüft, ohne Befund |
| Meldungen der Dateien `sqlite_lesen*.go`, `sqlite_json_lesen.go`: kein `%q`/`%v`/`%s` eines Werts aus der Datei, keine Zeilennummer (Base64 nennt nur einen Byte-Versatz im Wert); kein Wert aus Konfiguration oder Umgebung (der Diff berührt keine); ME1 rot | geprüft, ohne Befund außer F-623 |
| `internal/bootstrap/bootstrap.go`: zwei Stellen tauschen `recording.YAML{}` gegen `recording.Leser{}`, keine Logik; K1/K2 laut §7 am Binary rot | geprüft, ohne Befund |
| `internal/hexagon/**`: nicht berührt; `.a-check.yml`: `modernc.org/sqlite` und `modernc.org/sqlite/lib` nur im Recording-Adapter (Muster ohne Anker deckt den Unterpfad); `.golangci.yml`: unberührt, kein `//nolint` im Diff | geprüft, ohne Befund |
| `spec/spezifikation.md` (`LH-FA-22.a`, `SPEC-043`, Drift-Log): keine ADR-, Slice- oder Wellen-Kennung, Drift-Log nennt nur die Umplanung; Entscheidungen (Surrogat, Meldungen, Speicher, `ctx`, Gruppen, `meta`) decken sich mit dem Code | geprüft, ohne Befund |
| `test/integration/sqlite_lesen_e2e_test.go` und die zwei geänderten Aufzeichner: Roundtrip beider Formate über das Binary für `replay` und `play`, Datei unverändert (Bytes, Zeit, nichts daneben), Ablehnung in beiden Formaten mit Exit-Code 3 | geprüft, ohne Befund |
| `docs/user/abdeckung-*.md` (erzeugt); Benutzerhandbuch und `README.md` unangetastet (Doku-Slice) | geprüft, ohne Befund |
| §3.13: Nehmer `slice-v1-abschluss-sqlite-format-doku` führt `slice-v1-abschluss-sqlite-format-lesen` in Kopf, Ziel, DoD und §3; der Diff trägt keine neue Sendung ein; Nachzählung in §8 (drei Liefer-Punkte, eine Schicht) | geprüft, ohne Befund |
| §3.12: jede Randform des Lesens hat eine Zeile in §6 und eine Entscheidung vor `45af81d` bzw. in `ea1ef40`/`a6bcc45`; die Rückgaben des Implementers kamen vor der Nacharbeit | geprüft, ohne Befund |
| §3.3: kein Move im Diff; Commit-Messages ohne `SPEC-`/`ARC-`-Kennung ([ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md)) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zusage der Meldung weiter als Code und Test · Mutation grün, als äquivalent geführt, nicht äquivalent · Zusage im Vertrag weiter als der Test · Kommentar sagt mehr zu als der Code · Plan folgt der Korrektur nicht · Größe über der Schätzung · Plattformzweig ohne Prüfung · Randform des Schreibens ohne Entscheidung

## Verdikt

**Merge-blockierend:** nein für HIGH (keines); die zwei MEDIUM (F-623, F-624) gehen vor der Closure an den Implementer. Der Lesepfad (Erkennung, `mode=ro`, Lesetransaktion, Stufung, Version, `base64:`, Zahlen, Surrogat, Gruppen) ist durch die nachgefahrenen Mutanten gefangen, bis auf die in F-624 und F-625 benannten.

**Übergabe:**

- **F-623** an den **Implementer** (Pfad der Datei in jeder Meldung des Inhalts und die Art auf dem Weg über `Unmarshal`, oder die Zusage in der Spezifikation enger fassen — Entscheidung beim Architect, `AGENTS.md` §3.11); **F-624** an den **Implementer** (Test oder Mutation für `//`-Pfad, §7 berichtigen); **F-625, F-626** an den **Implementer**; **F-627** an den **Implementer** (`AGENTS.md` §3.9).
- **F-628** an den **Planner** (Entscheidung zu §4 festhalten); **F-629** an den **Planner** (Register bei der Closure); **F-630** an den **Architect**.
- Die **Finding-Klassen** gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg; DoD- und Spec-Konformität prüft der Verifier.
