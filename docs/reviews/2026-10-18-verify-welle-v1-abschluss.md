# Verifikations-Bericht: welle-v1-abschluss — 2026-10-18

**Rolle:** Verifier (Modul 11), frischer Kontext. **Frage:** Bauen wir es richtig, gemessen an Plan, DoD und Entscheidungen?

**Gegenstand:** Closure-Schritt 1 der Welle `welle-v1-abschluss` (`modul-06-roadmap.md` §Wellen-Closure-Prozedur): Die Abnahmeszenarien 1 bis 10, 12 bis 16 und 18 des Lastenhefts (§7) sind je Zusage durch einen Test belegt, der rot wird, wenn der Code die Zusage bricht; dazu der Replay-Lauf.

**Stand:** `main` auf `d1691e7a` (Arbeitsbaum sauber). Mutationen liefen in einer Kopie (`git archive HEAD`) unter eigenem Temp-Pfad, nie im Hauptbaum.

> **Zitier-Form.** Dieser Bericht friert ein: Kennung, nicht Adresse.

## Ergebnis

**Konform, ohne Blocker.** Jedes der 15 Szenarien ist durch mindestens einen Test belegt, der unter einer Mutation der Zusage rot wird (29 Mutationen am Quellcode, davon 27 im E2E rot und 2 grün, dazu 2 am `Dockerfile`). Zwei Mutationen blieben im E2E-Test grün (Befund F-1, F-2); beide fängt ein Unit-Test, nachgewiesen. Zwei Szenarien (5, 9) sind über die Abdeckungstabellen nicht auffindbar (F-3). Der Replay-Lauf ist grün.

## 1. Replay-Lauf und Basislauf

| Lauf | Ergebnis |
|---|---|
| `make test-integration` (Hauptbaum, volle Suite, PostgreSQL 17) | grün, 2 m 20 s, 276 PASS-Zeilen, kein FAIL; Phase 2 (PostgreSQL gestoppt, `PGR_OHNE_POSTGRES=1`) führt `TestE2EOhnePostgresReplay`, `TestE2EOhnePostgresExtendedReplay`, `TestE2EOhnePostgresEnde` (je Format `rec.yaml`, `rec.sqlite`) gegen die Aufzeichnungen der Phase 1: alle PASS |
| `make test-container` (Hauptbaum) | grün, 1 m 04 s: Punkt 4 (zwei Builds, gleiche Summen), 5 (Kopf von 6 Binaries), 6 (Image je Plattform), 7 (Smoke auf `linux/amd64`) |
| `make test-container` (unveränderte Kopie) | grün |

Der Replay-Lauf ist damit das Replay-Szenario gegen das gebaute Binary, ohne PostgreSQL, Exit 0.

## 2. Je Szenario: Mutation und roter Test

Methode: Der Runner `tools/test/run-integration-tests.sh` der Kopie, mit `PGR_TESTS_1`/`PGR_TESTS_2` auf die Tests des Szenarios; das Testimage wird je Mutation neu gebaut. Ich habe geprüft, dass jeder Lauf an einer Assertion des Tests rot wurde und nicht am Bau (ein Mutant, der nicht baute, wurde ersetzt, siehe 4.). Die Zeile `Datei: Änderung` ist die ganze Mutation.

| Sz. | Zusage | Mutation | roter Test (Grund) |
|---|---|---|---|
| 1 | Aufzeichnung entsteht | `services/record.go`: `Responses: responses` → `nil` | `TestE2ERecordSelect1`, `TestE2ERecordMehrereInteraktionen` rot |
| 2 | Replay ohne PostgreSQL liefert das Aufgezeichnete | `services/replay.go` `Query`: letzte Antwort der Interaktion weggelassen | `TestE2EReplaySelect1` rot (`context deadline exceeded`) |
| 3 | Client ohne Programmänderung | `pgwire/server.go`: Wert jedes `ParameterStatus` → `"x"` | `TestE2ERecordExtendedPgx` (`standard_conforming_strings=on` fehlt), `TestE2EReplayExtendedPgx` rot |
| 4 | Abweichung erkannt | `replay.go`: SQL-Vergleich `!= sql` → `!= sql && false` | `TestE2EReplayAbweichung` rot |
| 5 | keine interaktive Eingabe | `cli/datei.go`: Prüfung „reguläre Datei“ aus | `TestE2EConfigShowStdin` rot (`endet nicht binnen 15 s`) |
| 6 | Fehlerreplay | `pgwire/errorfields.go`: `Code` (SQLSTATE) → `""` | `TestE2EFehlerreplayEinfach` rot (Sicht beim Aufzeichnen) |
| 7 | Extended Query | `replay.go` `ClientMessage`: Server-Gruppe auf die halbe Länge gekürzt | `TestE2EReplayExtendedPgx` rot (`conn closed`) |
| 8 | drei Störungen unterscheidbar | `model/fehler.go`: (a) `CodeReplayMismatch` → Code von `CodeUpstream`; (b) `CodeRecordingBroken` → Code von `CodeReplayMismatch`; (c) `CodeUpstream` → Code von `CodeRecordingBroken` | (a) `TestE2EReplayAbweichung`; (b) `TestE2EReplayBeschaedigt`; (c) `TestE2ERecordUpstreamNichtErreichbar` — je eine Kollision färbt genau den Test ihres Falls |
| 10 | Argument, Umgebung, Datei, Priorität | `cli/cli.go` `zusammenfuehren`: (a) Umgebung vor Argument; (b) Umgebung ignoriert; (c) Datei ignoriert | (a) `TestE2ESessionsQuellen`, `TestE2ERecordTimingQuellen`, `TestE2EPlayKeepTimingQuellen/Option_false_vor_Umgebung`; (b) `TestE2EPlayVorrang`, die drei `…Quellen`; (c) `TestE2EPlayKonfiguration`, die drei `…Quellen` |
| 12 | Einspielen: Wirkung in der Datenbank | `postgres/einspielen.go`: gesendetes SQL → `SELECT 1` | `TestE2EPlayDDLDML` rot (`relation "log" does not exist`) |
| 13 | zeitgetreu / ohne Wunsch ohne Wartezeit | `play_zeit.go`: (a) Abstand halbiert; (b) `!KeepTiming`-Zweig aus | (a) `TestE2EPlayKeepTiming`, `TestE2EPlayKeepTimingQuellen/Option`; (b) `TestE2EPlayDDLDML` rot (ohne Wunsch, ohne Offsets) |
| 14 | SQLite-Datei | `recording/sqlite_lesen.go`: `FROM [message]` + `LIMIT 3` (Lesen); `sqlite_zeilen.go`: `fields` → `"{}"` (Schreiben) | Lesen: `TestE2EReplaySQLite`, `TestE2EReplayExtendedSQLite`, `TestE2EPlaySQLite` (`PGR-E3003`); Schreiben: `TestE2EReplaySQLite` |
| 15 | verschlüsselt / Klartext abgelehnt | `pgwire/server.go`: Klartextprüfung aus; `pgwire/tls.go`: `s.tls == nil` → `true` (immer „N“) | `TestE2EKlartextAbgewiesen`; `TestE2ERecordTLS`, `TestE2EReplayTLS` |
| 16 | Vergleich bei Struktur und Fehler; erwartetes Ende | `play_vergleich.go`: (a) SQLSTATE-Vergleich aus; (b) Spaltenvergleich aus; `play.go` `fatal`: (d) erwartetes Ende als Fehler, (e) Ende nie abweichend | (a) `TestE2EPlayVergleichAbweichung/anderer_SQLSTATE`; (b) `…/Spaltenbeschreibung`; (d) `TestE2EPlayEndeVergleichGleich`; (e) `TestE2EPlayEndeVergleichAnders` (beide Läufe) |
| 18 | Ende durch den Server; Netzwerkabbruch kein Ende | `postgres/upstream.go`: (d) `connection_closed` nach FATAL nicht angefügt; (c) Abbruch ohne Fehlerantwort als `connection_closed`; `pgwire/server.go`: (b) Verbindung nach `connection_closed` nicht beendet | (d) `TestE2ERecordEndeTerminate` (yaml und sqlite); (c) `TestE2ERecordNetzwerkabbruch` (Exit 6 statt 4); (b) `TestE2EOhnePostgresEnde` (yaml und sqlite) |
| 9 | Record und Replay im Image, Replay ohne PostgreSQL | `Dockerfile` Stufe `runtime`: Binary nach `/pgwire-recorder-x`; Stufe `build`: Zeit in `-X main.version` | `make test-container` rot: Image (`Could not find the file /pgwire-recorder`, Exit 2); Summen (`die Summen zweier Builds weichen ab`, alle sechs Binaries) |

Wiederholbarkeit: Skripte und Protokolle liegen im Scratchpad der Sitzung (`…/scratchpad/v18/`: `mut.sh`, `alle.sh`, `zwei.sh`, `log/*.log`, `log/*.diff`).

## 3. Befunde

| # | Stufe | Befund |
|---|---|---|
| F-1 | INFO | Mutation `play_vergleich.go` `erste`: `if u == nil \|\| ende == 0 \|\| u.position < ende` → `if true` (Szenario 16, dritter Satz) bleibt im E2E grün (`TestE2EPlayEndeVergleichGleich`, `…Anders`). Es ändert nur den Text der Abweichung (Art `Nachrichtenart` statt `Fehler`), nicht Zähler oder Exit-Code. Die Unit-Tests fangen es: `make test` in der Kopie rot (`TestVergleichFatalExtendedGruppen`, `play_vergleich_ende_test.go:175`). Kein Mangel der Zusage; der E2E-Test deckt die Zusage „erwartet / abweichend“, nicht den Text der Art. |
| F-2 | INFO | Mutation `record.go` `endetMitEnde` → `return false` (Szenario 18, Aufzeichnen: nach dem Ende nimmt die Session nichts mehr an) bleibt im E2E grün (`TestE2ERecordEndeTerminate`, `TestE2EOhnePostgresEnde`). Der Unit-Test `TestRecordEndeDanachNichts` fängt es (`record_ende_test.go:361`). Das Szenario selbst (Antworten, Aufzeichnung enthält das Ende) ist durch (d) und (b) in 2. belegt. |
| F-3 | LOW | Die Szenarien 5 und 9 sind über die Abdeckungstabellen `docs/user/abdeckung-*.md` nicht auffindbar: `LH-FA-15`, `LH-QA-04`, `LH-FA-16`, `LH-QA-03` haben keine Zeile (weder E2E noch Gate). Szenario 9 belegt `make test-container` (Vertrag `SPEC-055`), Szenario 5 `TestE2EConfigShowStdin` (Deklaration auf `LH-FA-17`). Folge: `make abdeckung-check` kann hier keine Lücke melden. Kein Verstoß gegen die Welle-Trigger (sie verlangen Tests mit roter Mutation, keine Tabellenzeile). Vorschlag an den Planner: Register-Beobachtung oder Folge-Slice, falls die Tabelle auch diese Anforderungen führen soll. |
| F-4 | INFO | Szenario 8 hat keinen Test, der die drei Störungen nebeneinander zeigt; die Unterscheidbarkeit ist über je einen Test mit festem Code belegt (Mutationen 8a bis 8c färben je genau einen). `abdeckung-gesamt.md` führt `LH-FA-14` mit Happy `—` („teilweise“). |
| F-5 | INFO | Szenario 12 spielt in `TestE2EPlayDDLDML` eine handgebaute Aufzeichnung ein (`einspielAufzeichnung`), nicht die aus Szenario 1 aufgezeichnete. Die Wirkung in der Datenbank ist belegt, die Kette Aufzeichnen → Einspielen im selben Test nicht. |

## 4. Verworfene Läufe (ehrlich benannt)

- Mutation 16(e) baute zuerst nicht (`declared and not used: ende`); ersetzt durch eine bauende Fassung (`u = nil` hinter dem Aufruf), die rot wurde. Zählt nur die bauende.
- `Dockerfile`-Mutation „Binary umbenannt“ traf beim ersten Versuch die Stufe `dist` (erste Fundstelle) statt `runtime`; der Lauf war rot am Bau der Binaries, nicht am Image. Wiederholt mit Anker `AS runtime`.
- Mutation 16(c), erster Wurf (`if true` in `erste`): Phase 2 lief rot, weil die gewählte Auswahl der Phase 1 die Aufzeichnungen der Phase 2 nicht erzeugte; gewertet ist nur Phase 1 (grün, siehe F-1).

## 5. Was nicht geprüft werden konnte

- **Szenario 11** (Homebrew): keine Bedingung dieser Welle (`welle-v1-nachtrag`).
- **Szenario 17**: laut Auftrag nicht Teil dieser Prüfung (Gates `make lint`, `make test`, Gegenproben liefen im `make gates` des Koordinators).
- **Plattformen:** nur Linux (Docker-Server `linux/amd64`); das Image `linux/arm64`, macOS und Windows prüft kein Lauf außer Kopf und Summe der Binaries (`make test-container` Punkt 5).
- **Szenario 9 Eingabefreiheit des Replay im Container:** belegt nur der Smoke (Replay ohne PostgreSQL, Exit 0 nach `SIGTERM`); eine Mutation, die den Replay im Container etwas verlangen ließe, habe ich nicht gefahren (zwei von mehreren möglichen `Dockerfile`-Mutationen).
- **Mutationen je Aufzählungselement** der Zusagen (etwa jede Pfadform, jedes Umgebungsvariablen-Paar): Stichprobe je Gruppe, nicht lückenlos. Szenario 10 deckt Argument, Umgebung, Datei und Vorrang am Kommando `play`, `record`, `replay` über die vorhandenen Tests, nicht jeden Schlüssel einzeln.
- **Laufzeitwirkung:** der volle `make gates` lief nicht erneut (Koordinator: Exit 0 um 04:07).

## 6. Aufräumen

Alle Container und Netze der Läufe (`pgr-mut-*`, `pgr-ct-*`) sind entfernt; die Runner prüfen das selbst (`SPEC-057`). Fremde Volumes (`pgr-anm-cache`, `pgr-dev_*`) habe ich nicht angefasst. Das Testimage `pgwire-recorder:integration` trug zuletzt eine Mutation; ich habe es aus dem Hauptbaum neu gebaut (Tag überschrieben). Der Hauptbaum ist unverändert (`git status` leer, bis auf diese Datei).
