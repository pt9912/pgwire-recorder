# Verifikation: slice-v1-abschluss-upstream-verbinden — 2026-10-09

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-09-review-slice-v1-abschluss-upstream-verbinden.md`, F-530 bis F-536). F-533, F-535 und F-536 hat er an mich übergeben.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-v1-abschluss-upstream-verbinden.md`, Kopf und §1 bis §8, gegen den Gesamt-Diff `9a04b12..b353056`. Architect: `9a04b12`. Implementer: `fe849ac`, `694eb99`, `e35fb56`, `8ba8caa`, `80cb896`, `74eb4b8`, `b353056`. Review: `a482303`.

**Abgrenzung:** Die Wirkung einer Verbindung bei `play` (Benutzer, Datenbank, Passwort, TLS nach `sslmode`, `PGR-E2005` in Benutzer, Passwort und Datenbank) liegt bei `slice-v1-abschluss-einspielen`. Dass sie fehlt, ist kein Befund.

**Eingang:**

- der Plan ganz, besonders §1, §2 (DoD), §3, §4, §6 (*Randformen*, A7 bis A12, F-521, F-528, U1 bis U8, *Risiken*) und §7 (*Belege des Implementers*, *Nacharbeit zum Review*)
- `spec/spezifikation.md` `LH-FA-17.a`: *Benannte Verbindungen* mit *Wirkung einer URL* und dem Absatz danach, *Geheimnisse*, *Fehler*, *Anzeige*
- am Stand `b353056`: `internal/adapters/driving/cli/upstream.go` ganz, der Diff von `cli.go`, `datei.go`, `verbindung.go`, `export_test.go`; `einsetzen_test.go` (`TestUpstreamVerbindung`, `TestUpstreamEinmalGelesen`, `TestEinsetzenAlleTeile`), `test/integration/verbindung_e2e_test.go` ganz; `tools/test/run-integration-tests.sh`
- `docs/user/benutzerhandbuch.md` §5 *Einstellungen* mit *Konfigurationsdatei* und §7 *Fehlercodes*, dazu der Diff des Slice daran
- der Review-Report ganz; den Bericht des Implementers habe ich nicht gelesen
- der Nehmer `slice-v1-abschluss-einspielen` §1 (*Übernommen aus* diesem Slice) und §3, dazu `AGENTS.md` §3.10 bis §3.13

Beim Start war der Arbeitsbaum sauber, HEAD `b353056`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-09

**So wurden die Proben gebaut.** Alles lief in einem eigenen Unterverzeichnis `ver-upv/` des Scratchpads, nie im Repo.

- **Binary:** eine frische Kopie per `git archive b353056 | tar -x`, ohne `.git` und ohne `cp -p`. Aus ihr wurde die Stufe `runtime` des `Dockerfile` mit eigenem Tag gebaut. Ein Skript legte ein internes Docker-Netz, ein Volume für die Aufzeichnungen und einen PostgreSQL-Container (gepinntes Image aus `harness/mk/integration.mk`) an. `record` lief als eigener Container in diesem Netz, der Client war `psql` aus dem PostgreSQL-Image. Jeder `docker`-Aufruf hatte eine Frist von 60 s, das ganze Skript 600 s. Ein `trap` entfernte am Ende Container, Netz und Volume; danach führte `docker ps -a`, `docker network ls` und `docker volume ls` mit dem Präfix nichts mehr. Für den IPv6-Fall lief `record` im Netz-Namensraum des PostgreSQL-Containers (`--network container:…`), damit `[::1]:5432` einen echten Server erreicht.
- **Unit-Mutanten:** je Mutant eine neue Kopie der Basis per `cp -r` ohne `-p`. Die Ersetzung machte ein Skript, das abbricht, wenn der Suchtext nicht genau einmal trifft. Danach lief `go test -count=1 ./internal/adapters/driving/cli/` im Image der Stufe `deps` mit der Kopie als Bind-Mount und ohne Netz; danach wurde nur das Verzeichnis des Mutanten gelöscht. Der Grundlauf ohne Mutation war grün.
- **Integrations-Mutant:** eine Kopie wie oben; in ihr `tools/test/run-integration-tests.sh` mit eigenem Image-Tag (statt `pgwire-recorder:integration`, damit das Image des Arbeitsbaums unberührt bleibt), Frist 900 s. Das Skript räumt Netz, Container und Volume selbst ab; danach waren keine mit dem Präfix `pgr-it` übrig. Kopie und Image habe ich danach gelöscht.
- **Danach:** Kopien und eigene Images (`runtime`, `deps`, `integration` mit eigenem Tag) sind entfernt. `make gates` lief im Arbeitsbaum; `git status` war danach sauber.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `--upstream`, `PGWIRE_RECORDER_UPSTREAM` und der Schlüssel `upstream` lösen den Namen auf, `record` verbindet zu Host und Port (`LH-FA-17`). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Name über die Option, `record` zeichnet auf | Binary: Datei mit `staging: "postgresql://GEHEIMUSER:${PGR_PW}@${PGR_H}:${PGR_P}/GEHEIMDB${PGR_DB}"`, `PGR_H`, `PGR_P` und `PGR_PW` gesetzt, `PGR_DB` leer, `--upstream staging`: `psql` erhält `42` auf `SELECT 41+1`, Startzeile `upstream=pgr-ver-upv-pg:5432`, Exit 0 nach `SIGTERM`, die Aufzeichnung enthält `sql: SELECT 41+1`. Die leere Variable in der Datenbank bleibt unbeachtet | bestätigt |
| Name über die Umgebungsvariable | Binary: `PGWIRE_RECORDER_UPSTREAM=staging`, `PGR_P=05432`: Abfrage aufgezeichnet, Startzeile `upstream=pgr-ver-upv-pg:05432` (U4, Port wie eingesetzt) | bestätigt |
| Name über den Schlüssel `upstream` | Binary: `record: upstream: staging` in der Datei, weder Option noch Umgebung: `SELECT 7` aufgezeichnet, Startzeile mit der Adresse | bestätigt |
| Name nur bei genauer Übereinstimmung, Auswahl nach dem genauen Namen (F-534) | Binary: Die Datei führt `Staging` (`falsch-host:1`) vor `staging`; `--upstream staging` verbindet zu `staging`. `--upstream STAGING` ist `PGR-E2001`. Mutanten **S**, **S2** (beide Namen klein), **S3** (`HasPrefix(name, v.name)`), **S4** (`HasPrefix(v.name, name)`): alle **rot**, `TestUpstreamVerbindung`, aus dem richtigen Grund (S: „`--upstream staging` zwischen ähnlichen Namen: `gross:5432`, erwartet `klein:5432`“, je für Option, Umgebung und Schlüssel) | bestätigt |
| Namen nur aus der gewählten Datei (U1) | Binary: `--upstream staging` ohne `--config` ist `PGR-E2001`. Mutant **U1** (`hatName` immer falsch): **rot**, neun Tests darunter `TestUpstreamNameNurAusDatei` | bestätigt |
| Jeder gesetzte Wert geprüft, auch der, der nicht gilt; Meldung ohne Wert (A7, A8, F-521) | Binary: `--upstream 'GEHEIM@h:5'` und `PGWIRE_RECORDER_UPSTREAM='GEHEIM h:5'` je `PGR-E2001` ohne `GEHEIM`; `PGWIRE_RECORDER_UPSTREAM=GEHEIM/x` neben `--upstream staging` ist `PGR-E2001` an der Umgebungsvariable. `:5432`, `h:0`, `h:65536` und `[fe80::1%a b]:5` (Leerraum in der Zone, V-122) `PGR-E2001`; `h:05432`, `[::1]:5`, `[fe80::1%GEHEIM]:5` gültig. Mutanten **O1** (Prüfung der Umgebung entfernt), **O2** (der Kommandozeile entfernt), **O3** (Umgebung nur ohne Kommandozeile): alle **rot**, `TestUpstreamUngueltig` und weitere | bestätigt |
| `record` verbindet zu Host und Port der Verbindung | Binary siehe oben. Mutant **A1** (`c.Record.Upstream = adresse` entfernt): **rot**, sechs Tests | bestätigt |
| IPv6-Host wieder in eckigen Klammern (U5, F-528, Integrationstest) | Binary: `sechs: "postgresql://GEHEIMUSER:${PGR_PW}@[::1]:5432/GEHEIMDB"` verbindet über `[::1]:5432` zum echten Server und zeichnet auf, Startzeile `upstream=[::1]:5432`; ebenso `${PGR_H6}` mit `::1`. Mutanten **K1** (nie geklammert) **rot** (`TestUpstreamEinsetzen`, `TestUpstreamVerbindung`), **K2** (immer geklammert) **rot** (sechs Tests). Integrations-Mutant **K1i**: **rot**, `TestE2ERecordVerbindungIPv6` mit „`Upstream ::1:1 nicht erreichbar: … too many colons in address`“ | bestätigt |
| Startzeile und Meldungen ohne Benutzer, Passwort, Datenbank (U6) | Binary: In keinem `record`-Lauf (vier Aufzeichnungen, neun Abbrüche beim Start) erscheint `GEHEIM` in `stdout`, `stderr` oder der Aufzeichnung; die einzigen Treffer im Protokoll sind die Namen meiner Fälle. Mutant **A2** (bei `record` auch Benutzer, Passwort, Datenbank eingesetzt): **rot**, `TestUpstreamVerbindung` | bestätigt |

### Punkt 2: Einsetzen, `PGR-E2005`, Port nach dem Einsetzen, `sslmode=require`, Reihenfolge am Ende (`LH-FA-17`). **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Nicht gesetzte Variable ist `PGR-E2005`, Stelle `connections.<Name>`, Name der Variable, kein Wert (A12, U3) | Binary: ohne `PGR_P`: „`Konfiguration [PGR-E2005]: Konfigurationsdatei: connections.staging: Umgebungsvariable PGR_P eines Platzhalters nicht gesetzt`“, Exit 2. Mutant **V1** (Code `PGR-E2004`): **rot**, drei Tests | bestätigt |
| Leere Variable gilt als nicht gesetzt | Binary: `PGR_H=` ergibt `PGR-E2005` mit `PGR_H`. Mutant **E1** (leerer Wert eingesetzt): **rot**, drei Tests | bestätigt |
| Wert unverändert, nicht gekürzt | Mutant **E2** (`strings.TrimSpace`, der zuerst grüne aus §7): **rot**, `TestUpstreamEinsetzen`, `TestUpstreamPortNachEinsetzen` | bestätigt |
| Port nach dem Einsetzen geprüft, `PGR-E2004` ohne Wert (Rückgabe 8, U3) | Binary: `PGR_P=GEHEIM` ergibt `PGR-E2004 … connections.staging: Port nach dem Einsetzen ist keine Zahl von 1 bis 65535` ohne `GEHEIM`. Mutanten **P1** (Prüfung entfernt), **P2** (Wert in die Meldung): beide **rot**, `TestUpstreamPortNachEinsetzen` | bestätigt |
| `sslmode=require` bei `record` `PGR-E2004` | Binary: `--upstream tls` ergibt `PGR-E2004 … connections.tls: sslmode=require ist bei record ungültig …`. Mutant **T1**: **rot**, `TestUpstreamReihenfolgeAmEnde` | bestätigt |
| Reihenfolge `--upstream`, `sslmode=require`, Variablen, Port | Binary: Bei `tls` fehlen keine Variablen in Host und Port, deshalb zeigt das Binary die Reihenfolge nur zwischen `--upstream` und dem Rest (`envUngueltigCliGut`: `PGR-E2001` vor jedem Verbindungsfehler). Mutant **R1** (`sslmode` nach dem Einsetzen): **rot**, `TestUpstreamReihenfolgeAmEnde`. Variablen vor dem Port: §7, mit `TestUpstreamVariableFehlt` (leere Variable im Port `PGR-E2005`) gelesen | bestätigt |
| `config show` meldet kein `PGR-E2005` | Binary: `config show` mit der Datei, deren Schlüssel `upstream` auf eine Verbindung mit nicht gesetzten Variablen zeigt: Exit 0, Platzhalter unaufgelöst | bestätigt |
| Einsetzen gilt für jeden Teil (U8) | `TestEinsetzenAlleTeile` über `SetzeEin` gelesen; E1 und V1 machen es rot | bestätigt |
| Variablen einmal gelesen (U2) | siehe Abschnitt 2, F-535 | bestätigt mit Grenze |
| Beleg in §7: je Zusage Zusage · Mutation · roter Test | Beide Tabellen sind vollständig gegen die Zusagen aus DoD-Punkt 1 und 2; meine 19 Unit-Mutanten und der Integrations-Mutant decken sich mit ihren Zeilen | bestätigt |

### Punkt 3: Handbuch §5 *Konfigurationsdatei* und §7 *Fehlercodes* wie geliefert, Abdeckung nachgezogen. **Bestätigt mit Einschränkung (V-125).**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Die drei Absätze zu Verbindungen, `record` mit `sslmode` und Platzhaltern | gegen das Binary gelesen: Form der URL, Name nur in genauer Schreibweise und nur aus der gewählten Datei, jeder Wert von Option und Umgebung geprüft, Schlüssel `upstream` `PGR-E2004`, bei `record` nur Host und Port, Klammern bei `:`, Startzeile mit der Adresse, `sslmode=require` bei `record` `PGR-E2004`, Einsetzen einmal und unverändert, leere Variable, `PGR-E2005`, Port nach dem Einsetzen. Jede Aussage stimmt mit einer Probe in Punkt 1 oder 2 überein. Keiner dieser Absätze sagt etwas über `play` | bestätigt |
| Zeilen `PGR-E2005` und `PGR-E2006` in §7 | `PGR-E2005`: Meldung nennt Verbindung und erste fehlende Variable, bei `record` nur Host und Port, wie am Binary gesehen. `PGR-E2006` ist Stand von `slice-v1-abschluss-verbindungen-platzhalter`, dort verifiziert | bestätigt |
| Beispiel und Abschnittsliste in §5 *Konfigurationsdatei* | Das Beispiel im selben Abschnitt wird vom gelieferten Binary abgelehnt, siehe V-125 | Einschränkung |
| Abdeckung über `make abdeckung` | `make abdeckung-check` im eigenen `make gates` grün | bestätigt |

### Punkt 4: `make gates` grün. **Bestätigt durch eigenen Lauf.**

Eigener Lauf im Arbeitsbaum am Stand `b353056`: Exit 0. Darin `baseline-verify: v6.16.0 OK`, `d-check: 404 Datei(en) geprüft, 0 Befund(e)`, `gesamt: 0 Befund(e)` (a-check), `a-check-negativ: gruen`, `--- PASS: TestE2ERecordVerbindungAusDatei`, `--- PASS: TestE2ERecordVerbindungIPv6`, `run-integration-tests: gruen`, `abdeckung-gegenprobe: gruen`, `kopf-check-gegenprobe: gruen`, `lint-gegenprobe: gruen`, `commit-msg-gegenprobe: gruen`. Der Beleg in §7 (*Nacharbeit*, „`make gates` · Commit der Nacharbeit · Exit 0“) deckt sich damit.

### Übrige DoD-Punkte (konstant je Slice)

- **Review:** Der Report liegt vor (`a482303`), aus einem anderen Lauf. Bestätigt.
- **Closure-Notiz, Register, Risiko-Ausgänge, Paarungen:** noch offen; §7 trägt `<…>`. Das ist die Reihenfolge, kein Befund. Für die Closure aus meiner Sicht: Das erste Risiko (Größe) trat nicht ein (F-536, 712 Zeilen in der Schätzung); das dritte (Handbuch im Zielstand) trat für das Beispiel in §5 ein (V-125); das vierte (Name ohne Auflösen) schließt DoD-Punkt 1, am Binary bestätigt. In die Register-Fortschreibung gehören die Klassen aus der Summary-Zeile des Reviews (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` für F-534, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` für F-530).

---

## 2. Übergaben des Reviews

| Befund | Prüfung | Urteil |
|---|---|---|
| F-533 Handbuch und `play` | In den Absätzen aus DoD-Punkt 3 steht nach `b353056` nichts mehr über `play`. Die drei verbliebenen älteren Stellen: (1) die Einleitung von §5 *Einstellungen* („Ausnahmen sind das Passwort …“) steht vor *Konfigurationsdatei*, also außerhalb von DoD-Punkt 3; sie beschreibt eine Einstellung von `play` wie die ganze Optionstabelle und das Kapitel *Eine Aufzeichnung in eine Datenbank einspielen*, die den Zielstand tragen — kein Befund gegen diesen Slice. (2) Das Beispiel mit `play:` und die Liste „(`record:`, `replay:`, `play:`)“ stehen **in** §5 *Konfigurationsdatei*, und das Binary lehnt das Beispiel ab: V-125. (3) Die Zeile `PGR-E4005` in §7 nennt `PGWIRE_RECORDER_PASSWORD` und `--upstream-tls`; DoD-Punkt 3 nennt in §7 nur `PGR-E2005` und `PGR-E2006` — kein Befund gegen diesen Slice. Übernahme nach `AGENTS.md` §3.13: `slice-v1-abschluss-einspielen` §1 nennt unter *Übernommen aus* `slice-v1-abschluss-upstream-verbinden` die Wirkung einer Verbindung bei `play` im Handbuch §5 *Konfigurationsdatei* mit der Kennung dieses Slice und F-533, §3 hat die Zeile `docs/user/benutzerhandbuch.md`; sein §1 *Ausdrücklich NICHT* schließt kein Handbuch aus, er liegt in `next/`. Die Adresse nimmt an — für die drei herausgenommenen Sätze. Das Beispiel nennt sie nicht (V-125) | Sendung angenommen; Beispiel V-125 |
| F-535 U2, `TestUpstreamEinmalGelesen` | Die Grenze ist ehrlich benannt: §6 U2 und §7 *Grenzen* sagen, dass die Mutation *Einsetzen beim Verbinden* im Zuschnitt nicht baubar ist, und dass der Test nur die Form der Übergabe sichert. Die erste Hälfte des Tests trägt eine Mutation: Mein Mutant **A1** (Adresse nicht in `c.Record.Upstream`, das Einsetzen damit nicht in `Parse`) macht `TestUpstreamEinmalGelesen` rot. Die zweite Hälfte (`Setenv` danach) kann nicht rot werden, weil `RecordOptions.Upstream` ein Text ist; das sagt §7 auch. Die Tabelle in §7 führt bei U2 als Mutation nur „Grenze, unten“, obwohl A1 die erste Hälfte hält (V-126) | ehrlich; V-126 |
| F-536 Größe | 712 hinzugefügte Zeilen im Diff `9a04b12..80cb896` bei geschätzten 660 bis 870, das Review lief in einer Sitzung; mit der Nacharbeit `b353056` 1033 Zeilen über alle 16 Dateien einschließlich Plan. Für die Closure: Das erste Risiko trat nicht ein | Hinweis für die Closure |

---

## 3. Hexagon und Architektur-Gate

- `make a-check` im eigenen `make gates`: `gesamt: 0 Befund(e)`, `a-check-negativ: gruen`.
- `git diff --name-only 9a04b12..b353056` ändert außerhalb von `docs/` nur `internal/adapters/driving/cli/` und `test/integration/verbindung_e2e_test.go`. Kern, Code-Tabelle im Model, PGWire-Adapter, Driven-Adapter, `internal/bootstrap`, `cmd/`, `go.mod`, `go.sum` und `spec/` sind unberührt.
- `upstream.go` importiert nur `errors`, `os` und `internal/hexagon/model`; die neuen Funktionen in `verbindung.go` nutzen `strings` und `model`. Kein Typ der YAML-Bibliothek ist beteiligt.

---

## 4. Plan gegen Code

- **§1:** Geliefert ist das *Benutzen*: Prüfung von Option und Umgebung, Auflösen, `sslmode=require`, Einsetzen, `PGR-E2005`, Port, Zusammensetzen, Handbuch. Die Wirkung bei `play` fehlt, wie §1 sagt. `PGR-E2005` hat jetzt einen Erzeuger (`einsetzen`). Die Schicht-Abgrenzung hält (Abschnitt 3).
- **§3:** Jede geänderte Datei steht dort. Der Bootstrap ist „keine Änderung erwartet“ und unverändert; `RecordOptions.Upstream` trägt die zusammengesetzte Adresse, wie die Startzeile am Binary zeigt. `spec/spezifikation.md` ändert nur der Architect-Commit `9a04b12`.
- **§4:** Keine Rückführung trat ein; der Zerleger ist unverändert.
- **§6:** Jede Randform, die ich am Binary gefahren habe, folgt ihrem Ort in `LH-FA-17.a`: U1, U3 bis U6, A7 bis A9, A12, F-521 mit der Zone aus V-122, F-528, Rückgabe 8. Keine Randform außerhalb von §6 ist im Code entschieden.
- **§7:** Die Tabellen decken DoD-Punkt 1 und 2 in der Form *Zusage · Mutation · roter Test*; jede Zeile, die ich nachgefahren habe, deckt sich: 19 eigene Unit-Mutanten, alle rot, dazu ein Integrations-Mutant, rot. Die Läufe sind mit Stand belegt, auch `make gates` am Stand der Nacharbeit. Die Einordnung der verbliebenen `play`-Stellen als „außerhalb der Absätze aus DoD-Punkt 3“ trägt für die Einleitung und `PGR-E4005`, für das Beispiel nur halb (V-125).

---

## 5. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-125 | LOW | Das Beispiel in §5 *Konfigurationsdatei* — das einzige, das Verbindungen, einen Platzhalter und `sslmode` zeigt — wird vom gelieferten Binary abgelehnt. Probe: das Beispiel wörtlich als Datei, `config show --config` und `record --config … --upstream lokal` enden beide mit „`Konfiguration [PGR-E2004]: Konfigurationsdatei: play: unbekannter Schlüssel`“, Exit 2. Ebenso nennt der Satz davor einen Abschnitt `play:`, den das Laden nicht kennt. DoD-Punkt 3 verlangt, dass §5 *Konfigurationsdatei* Verbindungen, `sslmode` und Platzhalter *wie geliefert* beschreibt; wer das Beispiel kopiert, startet nicht. Die Stelle stammt nicht aus diesem Slice, und mit `slice-v1-abschluss-einspielen` (dort §1, *Übernommen aus* `slice-v1-abschluss-konfigurationsdatei`: der Abschnitt `play:`) wird sie wahr; dessen Sendung aus F-533 nennt das Beispiel aber nicht, und schneidet jemand den Nehmer, bleibt es falsch. Das ist der eingetretene Teil des dritten Risikos in §6. | `docs/user/benutzerhandbuch.md` §5 *Konfigurationsdatei* · „übrigen in einem Abschnitt je Betriebsart (`record:`, `replay:`, `play:`)“ und der YAML-Block mit `play:` | Implementer: das Beispiel auf einen Abschnitt umstellen, den das Binary heute liest (etwa `record:` mit `upstream: lokal`), und `play:` aus der Liste nehmen, oder Planner: die Sendung in `slice-v1-abschluss-einspielen` §1 um Beispiel und Liste erweitern (§3.13) und den Ausgang des dritten Risikos so eintragen |
| V-126 | INFO | Die Zeile U2 in §7 *DoD-Punkt 2* führt als Mutation nur „Grenze, unten“. Die erste Hälfte von `TestUpstreamEinmalGelesen` hält aber eine Mutation: Ohne die Übernahme der Adresse in `Parse` (mein Mutant A1) wird der Test rot. Der Beleg ist damit schwächer notiert, als er ist; die Grenze selbst ist richtig benannt (F-535). | `docs/plan/planning/in-progress/slice-v1-abschluss-upstream-verbinden.md` §7 · „Variablen einmal gelesen (U2) \| Grenze, unten“ | Implementer: bei der Closure A1 als Mutation der ersten Hälfte eintragen, die Grenze für die zweite stehen lassen |

Kein HIGH, kein MEDIUM. DoD-Liefer-Punkte 1, 2 und 4 sind bestätigt, Punkt 3 mit der Einschränkung V-125.

---

## 6. Gefahrene Mutanten

Alle am Stand `b353056`, je in einer frischen Kopie, `go test -count=1 ./internal/adapters/driving/cli/`. Ausnahme ist K1i.

| Mutant | Änderung | Ergebnis |
|---|---|---|
| basis | keine | grün |
| S | `datei.verbindung`: `strings.EqualFold(v.name, name)` | rot: `TestUpstreamVerbindung` |
| S2 | beide Namen kleingeschrieben verglichen | rot: `TestUpstreamVerbindung` |
| S3 | `strings.HasPrefix(name, v.name)` | rot: `TestUpstreamVerbindung` |
| S4 | `strings.HasPrefix(v.name, name)` | rot: `TestUpstreamVerbindung` |
| U1 | `hatName` immer falsch | rot: neun Tests, darunter `TestUpstreamNameNurAusDatei`, `TestDateiUpstreamUngueltig` |
| O1 | Prüfung der Umgebungsvariable entfernt | rot: `TestUpstreamNameNurAusDatei`, `TestUpstreamReihenfolgeAmEnde`, `TestUpstreamUngueltig` |
| O2 | Prüfung der Kommandozeile entfernt | rot: `TestUpstreamNameNurAusDatei`, `TestUpstreamUngueltig` |
| O3 | Umgebung nur ohne Kommandozeile geprüft | rot: `TestUpstreamReihenfolgeAmEnde`, `TestUpstreamUngueltig` |
| A1 | Adresse nicht in `c.Record.Upstream` übernommen | rot: sechs Tests, darunter `TestUpstreamEinmalGelesen` |
| A2 | bei `record` auch Benutzer, Passwort, Datenbank eingesetzt | rot: `TestUpstreamVerbindung` |
| K1 | nie geklammert | rot: `TestUpstreamEinsetzen`, `TestUpstreamVerbindung` |
| K2 | immer geklammert | rot: sechs Tests |
| E1 | leerer Wert eingesetzt | rot: `TestEinsetzenAlleTeile`, `TestUpstreamReihenfolgeAmEnde`, `TestUpstreamVariableFehlt` |
| E2 | eingesetzter Wert mit `strings.TrimSpace` | rot: `TestUpstreamEinsetzen`, `TestUpstreamPortNachEinsetzen` |
| V1 | `PGR-E2005` durch `PGR-E2004` ersetzt | rot: drei Tests |
| T1 | `sslmode=require` nicht geprüft | rot: `TestUpstreamReihenfolgeAmEnde` |
| R1 | `sslmode=require` nach dem Einsetzen geprüft | rot: `TestUpstreamReihenfolgeAmEnde` |
| P1 | Port nach dem Einsetzen nicht geprüft | rot: `TestUpstreamPortNachEinsetzen` |
| P2 | Wert des Ports in die Meldung | rot: `TestUpstreamPortNachEinsetzen` |
| K1i | wie K1, `tools/test/run-integration-tests.sh` in der Kopie | rot: `TestE2ERecordVerbindungIPv6` |

---

## Summary

**Summary:** 0 HIGH · 0 MEDIUM · 1 LOW · 1 INFO.

- V-125: Das Beispiel in §5 *Konfigurationsdatei* wird vom gelieferten Binary mit `PGR-E2004` (`play: unbekannter Schlüssel`) abgelehnt; DoD-Punkt 3 hält nur mit dieser Einschränkung.
- V-126: §7 notiert für U2 keine Mutation, obwohl A1 die erste Hälfte des Tests hält.

**Urteil:** Der Slice ist richtig gebaut. DoD-Liefer-Punkte 1, 2 und 4 sind am echten Binary gegen PostgreSQL, an 19 Unit-Mutanten und einem Integrations-Mutanten bestätigt, alle Belege aus §7, die ich nachgefahren habe, decken sich. Punkt 3 ist bestätigt bis auf V-125, das vor der Closure behoben oder als Sendung im Nehmer eingetragen werden sollte. F-534 ist erledigt (S bis S4 rot), F-533 eingeordnet, F-535 ehrlich benannt.
