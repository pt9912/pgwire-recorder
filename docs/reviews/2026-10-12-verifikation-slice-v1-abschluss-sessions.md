# Verifikation: slice-v1-abschluss-sessions — 2026-10-12

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-12-review-slice-v1-abschluss-sessions.md`, F-679 bis F-685); die Nacharbeit (`0d23ed8`, `60a5aaf`, `6df6976`) und die Entscheidungen des Architect (`8813bf5`, `937ad35`, `a5e47c5`) prüfe ich nach. Ich repariere nichts.

**Gegenstand:** Slice-Plan `slice-v1-abschluss-sessions` (Kopf, §1 bis §8) gegen den Gesamt-Diff `b12f4b4..6df6976` (22 Dateien, +2351 −87). Implementer: `6f428a1`, `ecacdd8`, `ed3a4c9`, `b19f33a`, `3817ed7`, `efc17ec`, `aae9097`, `2f9a696`, `0d23ed8`, `60a5aaf`, `6df6976`. Review: `d3d0b2a`. Architect: `8813bf5`, `937ad35`, `a5e47c5`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

**Eingang:** der Plan ganz am Stand `6df6976` (§7 *Belege des Implementers*; die Closure-Punkte sind offen, weil die Closure sie führt, das ist kein Befund); `LH-FA-12.a` (*Recording*, *Randformen des Recordings*, *Replay*, *Randformen der Zuordnung*), `LH-FA-13.a`, `LH-FA-13.b`, `LH-FA-09.a`, `LH-FA-03.b`, `LH-FA-17.a` in `spec/spezifikation.md`; der Review-Report; der Code `internal/hexagon/services/record.go` und `replay.go`, `internal/hexagon/ports/driving/replay.go`, `internal/adapters/driving/cli/cli.go`, `internal/bootstrap/bootstrap.go`, `test/integration/sessions_e2e_test.go`. Den Bericht des Implementers (außer §7 im Plan) habe ich nicht gelesen. Beim Start war der Arbeitsbaum sauber, HEAD `6df6976`.

**So wurden die Proben gebaut.** Alles Docker-only, Scratch im Scratchpad der Sitzung, nie im Repo. Binary: `make build` (Image `pgwire-recorder:dev`, Image-ID gleich dem vorhandenen, also HEAD). Ende-zu-Ende mit dem Binary: eigenes internes Docker-Netz, `postgres:17-alpine` mit dem gepinnten Digest aus `harness/mk/integration.mk`, Client `psql` und `nc` im PostgreSQL-Container. Mutanten: je Mutant eine frische Kopie des Baums (`git ls-files -co`, ohne `.git`, `tar -xm`), genau eine `perl`-Ersetzung (das Skript bricht ab, wenn sie nichts ändert), `gofmt -l` im selben Lauf (leer in allen Läufen), die mutierte Datei danach mit `touch`; Unit-Mutanten mit `go test -count=1` im Image `pgwire-recorder:test` ohne Netz, die Kopie als Bind-Mount nach `/src` (die Mutation kommt an: kein Build-Kontext); Mutanten am Binary mit `make test-integration` aus der Kopie. Die unveränderte Kopie (der Baum selbst, schreibgeschützt eingehängt) lief vorher grün für `services`, `cli` und `bootstrap`. Container und Netze trugen die Namen `vfy-*` und sind entfernt; `git status` am Ende leer.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: `LH-FA-12` (Sessions, `--record-empty-sessions`, `--session-assignment connection`). **Bestätigt, mit V-305 (LOW) und V-306 (LOW).**

Alle in §7 genannten Testnamen der Stichprobe existieren (die elf Funktionen von `test/integration/sessions_e2e_test.go` stehen im Lauf von `make test-integration`; die Unit-Tests tragen die Mutanten unten rot). Ohne den Fix aus dem richtigen Grund rot: 14 Mutanten selbst gefahren, unabhängig von der Tabelle neu gebaut, jeder bricht nur die Zusage im Wortlaut von `LH-FA-12.a`.

| Zusage (Wortlaut) | Mutant (nur sie) | rot, mit Test |
|---|---|---|
| Zwangsende im Aufbau: nur bei `EndForced` (mehrteilig, Bedingung 1) | `imAufbau := end == EndForced && …` ohne `end == EndForced &&` (M2) | `TestRecordLeereSessionJedesEnde`, `TestRecordLeereSessionNachSchreibfehler`, `TestRecordGleichzeitigesEnde` |
| … nur ohne abgeschlossene Interaktion (Bedingung 2) | `abgeschlossen == 0 &&` entfällt (M3) | `TestRecordZwangsende` |
| … nur ohne laufende Interaktion (Bedingung 3) | `&& !laeuft` entfällt (M4) | `TestRecordZwangsendeLeereSession` |
| … ohne Meldung `PGR-E4006` (verneinend) | die Nummer der verworfenen Interaktion wird bei jedem `EndForced` gesetzt (M5) | `TestRecordZwangsendeLeereSession`, `TestRecordZwangsende` |
| Ohne die Option wird eine Session ohne abgeschlossene Interaktion nicht geschrieben (verneinend) | `\|\| abgeschlossen == 0 && !leere` entfällt (R2) | u. a. `TestRecordSessionOhneAnfrage`, `TestRecordZwangsende` |
| Die `id` zählt ab 1, ohne Lücke | `len(…)+1` zu `len(…)` (R6i) | acht Tests, u. a. `TestRecordKennungNachReihenfolgeDerEnden`, `TestRecordSessionMitInteraktionen`, `TestRecordGleichzeitigeSessions` |
| `--session-assignment`: `first_request` ist kein Wert (Wertemenge, Alias, Schreibweise) | `first_request` zusätzlich gültig (C1b) | `TestSessionAssignmentWertemenge`, `TestDateiUngueltig`, `TestConfigShowFehler` |
| Nie zugeordnet ist nur eine Session mit `id` über der Zahl der Verbindungen | `>` zu `>=` in `Unassigned` (P3b) | `TestReplayConnectionZuordnungMitDerAnnahme` |
| Die n-te Verbindung erhält die Session mit der `id` n | `alle[next]` zu `alle[next+1]` (P1b) | sieben Tests, u. a. `TestReplayConnectionNteVerbindungNteSession`, `…SessionOhneInteraktionNichtZuMelden` |
| Handshake nach der zugeordneten Session | `vorlage = sess` entfällt (P2a) | `TestReplayConnectionNteVerbindungNteSession`, `…HandshakeUnabhaengigVonAnfragen` |
| Eine Verbindung ohne Session: `PGR-E5003` | `CodeReplayMismatch` statt `CodeReplaySession` (P5b) | `TestReplayConnectionMehrVerbindungenAlsSessions` |
| Der Bootstrap reicht `connection` durch, und nur `connection` | Bedingung immer wahr (B3) | `TestRunReplayAusDatei`, `TestRunLogLevel`, `TestRunReplayConnectionKennzeichnung` |
| `--record-empty-sessions` wirkt am Binary (yaml und sqlite, Kennzeichnung, Sessions ohne Interaktion) | `services.LeereSessions(false)` im Bootstrap (I1, `make test-integration`) | neun Tests rot: `TestE2ERecordEmptySessions`, `…ParalleleVerbindungen`, `…KeineVerbindung`, `…ZwangsendeImAufbau`, `…HalbeStartnachricht`, `TestE2EReplayConnection`, `…OhneKennzeichnung`, `TestE2EReplayZaehlungDerVerbindungen`, `TestE2ESessionsQuellen` (die Liste der Tabelle in §7 plus `TestE2ERecordHalbeStartnachricht`) |
| F-680: Startnachricht mit Kopf, Rest fehlt, ist `PGR-E6001`, Exit-Code 6 | `s.fail` durch `s.log.Debug` in `startup` (PGWire-Adapter, nur in der Kopie) | `TestE2ERecordHalbeStartnachricht` (`Exit-Code 0, erwartet 6`), sonst kein Test rot |

Sechs verlangte Arten sind abgedeckt: mehrteilige Bedingung (M2, M3, M4 je Bedingung), verneinend (M5, R2), Wertemenge von `--session-assignment` (C1b).

**Am gebauten Binary** (selbst gefahren, ohne Datenbank, `--network none`):

- `record --record-empty-sessions=maybe`, `replay --session-assignment Connection`, `""`, `first_request`: `PGR-E2001`, Exit 2; die Meldung nennt „erlaubt sind first-request und connection“. Kommandozeile beim falschen Kommando (`record --session-assignment`, `replay --record-empty-sessions`): `PGR-E2001`, `flag provided but not defined`, Exit 2.
- Umgebung: `PGWIRE_RECORDER_SESSION_ASSIGNMENT=connection` bei `record` und `PGWIRE_RECORDER_RECORD_EMPTY_SESSIONS=true` bei `replay`: `PGR-E2001` (`gehört zu einer Option, die record/replay nicht kennt`), Exit 2; ungültige Werte der Variablen: `PGR-E2001`, ohne den Wert.
- Datei: `replay.session_assignment: Connection` und `record.record_empty_sessions: maybe` sind `PGR-E2004` (Wert abgelehnt, ohne Wert in der Meldung); `replay.record_empty_sessions`, `record.session_assignment` und der Schlüssel oben sind `PGR-E2004 … unbekannter Schlüssel`. `config show` mit beiden gültigen Schlüsseln nennt sie.
- Hilfe: `record --help` führt `--record-empty-sessions[=true|false]` mit Variable, `replay --help` `--session-assignment first-request|connection` mit Variable; die allgemeine Hilfe nennt beide.
- Kennzeichnung fehlt bei `connection`: `PGR-E2001 … trägt keine Kennzeichnung empty_sessions`, Exit 2, mit `--network none` und ohne Datenbank (der Port wird nicht geöffnet); mit nicht vorhandener Eingabedatei kommt `PGR-E3001` (Exit 3) vor `PGR-E2001`.

**Am Binary gegen PostgreSQL 17:**

| Probe | Ergebnis |
|---|---|
| `record --record-empty-sessions`: `select 1`, `select 2` auf zwei `psql`-Verbindungen, dazu eine TCP-Probe ohne Bytes, eine mit 7 Bytes | Aufzeichnung mit `empty_sessions: true` und genau den Sessions `id 1`, `id 2` (Probe und `pg_isready` erzeugen keine) |
| volle Startnachricht (`user`, `database`), dann Ende | eine Session ohne Interaktion, `id 1` |
| halbe Startnachricht, 8 Bytes Kopf (Länge 100), dann Ende | `PGR-E6001 … Startnachricht nicht lesbar: unexpected EOF`, Exit 6, `sessions: []` mit `empty_sessions: true` |
| 7 Bytes, dann Ende | keine Meldung, Exit 0, `sessions: []` |
| Aufzeichnung dreier Verbindungen (`select 1`; Startnachricht ohne Anfrage; `select 3`), danach `replay --session-assignment connection --fail-on-unconsumed` mit denselben drei Verbindungen | Antworten `1` und `3`, die leere Session an der zweiten Verbindung, Exit 0, keine Warnung |
| derselbe Replay, die zweite Verbindung fragt `select 3` | `connection to server was lost`, die dritte Session bleibt nie zugeordnet: `PGR-W2001 … die erste mit Kennung 3`, Exit 5 |

Fehlende Beleg-Teile (§7 verlangt, DoD verlangt): keine.

### Punkt 2: `LH-FA-13` (Verbindungsfehler beendet nur die Verbindung, der Prozess merkt die Klasse). **Bestätigt.**

`TestE2EVerbindungsfehlerParallel` lief grün im vollen Lauf von `make test-integration` (Exit 0, 118 PASS, 0 FAIL). Der Mutant I2 („der letzte Fehler gilt“) und die Wartezeit (F-682) habe ich nicht neu gefahren, der Beleg ist die Tabelle in §7 und der Lauf; V-307.

### Punkt 3: `make gates` grün. Nicht von mir gefahren, siehe Verdikt.

---

## 2. Bindungen und Entscheidungen

- **Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando` erfüllt.** `nichtGelesen()` in `fremd_test.go` führt nur noch `tls-cert`, `tls-key`, `allow-plaintext`; die Leser stehen in `optionen("record")` und `optionen("replay")` in `cli.go`; `zuordnung()` nennt `record-empty-sessions` für `record` und `session-assignment` für `replay` (Spalte *Kommando*). Die Mutanten C2a, C4a, C4b sind im Plan; die Variablen beim anderen Kommando habe ich am Binary als `PGR-E2001` gesehen (oben). Kommentar und Zahl („drei“) in `fremd_test.go` sind nachgezogen.
- **F-679 bis F-685 sind behoben oder entschieden eingetragen.**
  - F-679: akzeptiertes Negativ, entschieden in §6 und in der Spezifikation (*Angenommene Verbindung*, Satz zur ruhenden Session), Kommentare von `wirdVerworfen` und des Tests nennen beide Fälle. Kein Code. Der Dienst unterscheidet die beiden Fälle nicht, das Verhalten ist zugesagt als „nicht geschrieben, keine Meldung“.
  - F-680: Spezifikation hält den Bestand (Tabelle zu `LH-FA-05.e` und `PGR-E6001` in der Meldungstabelle), Test `TestE2ERecordHalbeStartnachricht` beide Ausprägungen; am Binary und per Mutant bestätigt (oben).
  - F-681: Spezifikation (*Recording*, „Ein Sessionende ist der Zeitpunkt, zu dem der Recorder das Ende bemerkt“), kein Test, wie in §6 entschieden.
  - F-682: `warteAufSchlaf` ersetzt die 300 ms (`grep Sleep` nennt nur die 50 ms der Abfrageschleife in Zeile 54). Die Lastprobe in §7 sagt selbst, dass sie das Flackern der alten Fassung nicht zeigte; kein Befund, Aussage ist ehrlich.
  - F-683: §1 Nahtstelle 3, die Zeilen `internal/hexagon/services` und `test/integration` in §3 und §8 sind nachgezogen (Code gelesen: `wirdVerworfen`, `LeereSessions`, `Finish` stimmen mit §3). Offen bleibt die Eintragung bei den Nehmern, V-308.
  - F-684: Entscheidung des Architect in §8 (Verdrahtung). `bootstrap.go` ruft `NewRecordService`/`NewReplayService` mit einem Zweig auf einen Konstanten-Vergleich und liest keine Datei und stuft keinen Fehler ein: Die Lesart des Architect trägt; ein Mutant (B3) ist rot.
  - F-685: Kommentar an `frei` nennt jetzt, was die Tests der Zuordnung fangen (P1a bis P1c). Ein Lauf mit `-race` fehlt weiter, das sagt §7 selbst.
- **PGWire-Adapter unberührt.** `git diff --name-only b12f4b4..HEAD | grep -c pgwire/` ist 0; ebenso `internal/adapters/driven`, `spec/architecture.md`, Format und Recording-Adapter. Mein F-680-Mutant liegt nur in einer Kopie.
- **Handbuch und README unberührt** (`git diff --stat b12f4b4..HEAD -- docs/user/benutzerhandbuch.md README.md` leer; `docs/user/abdeckung-*.md` sind erzeugte Tabellen). Der Doku-Folge-Slice `slice-v1-abschluss-sessions-doku` trägt sie (§1).
- **Aufruf-Voraussetzung der Zuordnung in der Spezifikation gebunden.** `LH-FA-12.a` *Replay*: „Die Zuordnung ist deterministisch, wenn die Clients ihre Verbindungen nacheinander aufbauen und nutzen; bei gleichzeitiger Nutzung … nicht zugesichert“; *Gleichzeitiges Ende* im Record nennt dieselbe Voraussetzung; die Aufruf-Voraussetzung der Client-Richtung (nacheinander) steht in `LH-FA-21.a` *Messpunkt* aus dem Geber. Kein Test, wie in §6 entschieden.
- **Der Plan folgt dem Diff (§3.9).** §1, §3 (sechs Zeilen), §6 (Randformen) und §8 stimmen mit dem Code überein: Namen `LeereSessions`, `ZuordnungNachVerbindung`, `wirdVerworfen`, `alle`, `keineWeitere`, `ReplayOptions.SessionAssignment`, `ZuordnungErsteAnfrage`/`ZuordnungVerbindung` stehen im Code; `keineWeitere` hat `letzteNummer` ersetzt. `make kopf-check` Exit 0, `make abdeckung-check` Exit 0, `make docs-check` 658 Dateien, 0 Befunde (selbst gefahren).
- **Größenprobe.** Der Bericht in §7 nennt für `b12f4b4..b19f33a` „18 Dateien, +2027 −73, ohne die vier Tabellen +1958 −68 (Produktcode +170 −56, Tests +1788 −12)“: stimmt (`git diff --shortstat`, `--numstat`). Der Stand nach der Nacharbeit ist nicht berichtet, siehe V-306.

---

## 3. Kommentare nach `AGENTS.md` §3.7 und §3.11

Geprüft gegen die Tests: `ReplayService`, `NewReplayService`, `OpenConnection`, `zuordnen`, `cursor`, `frei`, der Port `Replayer`, `CloseSession`, `wirdVerworfen`, `nichtGeschrieben`, `EndForced`, der Hilfetext. Alle im Indikativ, ohne verworfene Alternative und ohne abwesenden Text; jede Zusage hat einen Test, den ich als roten Mutanten gesehen habe (`OpenConnection`: P1b, P2a; `zuordnen`: P5b; `CloseSession`: R6i, M2 bis M5; `wirdVerworfen`: M2 bis M4, R2; `frei`: P1a bis P1c laut Tabelle). Ausnahme: `LeereSessions`, V-305. Die Aussage in `wirdVerworfen`, die zweite Bedingung treffe auch eine ruhende Session, ist die dokumentierte Folge von F-679 (kein Zusatz-Versprechen: Test und Spezifikation nennen das Verhalten).

---

## 4. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-305 | LOW | Der Kommentar von `LeereSessions` sagt: „Mit an wird jede eröffnete Session geschrieben, auch eine ohne abgeschlossene Interaktion“. Das ist weiter als der Code und als die Tests: Eine Session mit nicht unterstützter Interaktion (`TestRecordLeereSessionNichtUnterstuetztVerworfen`) und das Zwangsende im Aufbau (`TestRecordZwangsendeLeereSession/im_Aufbau`) werden mit der Option nicht geschrieben. Der Satz in der Spezifikation („Jede eröffnete Session wird am Ende geschrieben …“) trägt die Ausnahmen im selben Absatz; der Kommentar nicht. Die Hilfe von `record` („zeichnet auch eine Verbindung ohne Anfrage … auf“) ist enger und stimmt. | `internal/hexagon/services/record.go` · `// LeereSessions stellt --record-empty-sessions ein: Mit an wird jede eröffnete Session geschrieben` | Implementer (Kommentar auf „auch eine ohne abgeschlossene Interaktion, außer …“ verweisen oder enger fassen, `AGENTS.md` §3.11); ohne Blockade |
| V-306 | LOW | Die Größenprobe in §7 gilt für `b19f33a` (+2027 −73, 2100 mit Löschungen). Nach der Nacharbeit liegt der Diff ohne Plan und Review-Report bei +2127 −79 (2206 mit Löschungen, mit den Tabellen; +2056 −74 ohne Tabellen), mit Plan und Review +2351 −87. §4 nennt „über 2200 Zeilen mit Tests“ ohne Metrik; nach Einfügungen bleibt der Slice unter 2200, nach Einfügungen plus Löschungen liegt er knapp darüber. Die Rückführung trat nicht ein, die Probe nennt aber den Stand nicht, an dem die Schwelle gilt, und nicht die Metrik. | `docs/plan/planning/in-progress/slice-v1-abschluss-sessions.md` §4 (Zeile *Der Diff liegt über 2200 Zeilen*), §7 *Läufe* | Planner (Metrik der Schwelle festlegen: Einfügungen oder Summe, mit oder ohne erzeugte Tabellen); Implementer: Stand nach der Nacharbeit in §7 nachtragen; ohne Blockade |
| V-307 | INFO | Nicht selbst gefahren: `make gates`, `make test` (volle Kette mit `gofmt` und `go vet`; gefahren habe ich die drei Pakete `services`, `cli`, `bootstrap` und die Mutanten), `make lint`, `make a-check`, die Lastprobe zu F-682, der Mutant I2 („der letzte Fehler gilt“), I3 bis I5, die restlichen rund 55 Mutanten der Tabelle in §7 (gefahren: 14, jeder mit einer Ersetzung in einer frischen Kopie, Auswahl nach dem Wortlaut der Spezifikation), ein Lauf mit `-race`. Gefahren: `make build`, `make docs-check`, `make kopf-check`, `make abdeckung-check`, volle `make test-integration` (grün, 118 PASS, 0 FAIL). | — | ohne Aktion |
| V-308 | LOW | Die Nahtstelle aus §1 (Punkte 1 bis 4, „Sendung nach `AGENTS.md` §3.13, Eintrag in deren §1 jeweils durch den Planner“) steht in den Nehmern nur teilweise. `slice-v1-abschluss-verbindungsende-replay` nennt `slice-v1-abschluss-sessions` nicht (`grep -c` = 0), obwohl Punkt 1 (`ohneLebendpruefungen` ist die eine Stelle) und Punkt 2 (`OpenConnection` setzt die Zuordnung) es binden. `slice-v1-abschluss-verbindungsende-record` nennt in §1 „`CloseSession`, `--record-empty-sessions` und die Kennung … unverändert zu benutzen“, nicht aber, dass das Zwangsende im Aufbau nie geschrieben wird und ein Ende des Servers vor dem ersten `ReadyForQuery` keine Session erzeugt (Punkte 3 und 4). Der Implementer verweist auf den Planner (§7, F-683). Die Sendung steht damit im Geber, nicht im Nehmer (`AGENTS.md` §3.13: im selben Commit); der Nehmer ist ein Slice in `next/` oder `open/` und kann annehmen. | `docs/plan/planning/*/slice-v1-abschluss-verbindungsende-replay.md` §1, `…-record.md` §1 | Planner (Sendung in beide §1 eintragen, Liefer-Punkte und Schichten nachzählen); vor der Closure |
| V-309 | INFO | Die dritte Zeile des Belegs zu F-680 („keine Session, auch mit der Option“) hat keinen eigenen Mutanten, mit Grund (Mutant läge im PGWire-Adapter, außerhalb des Slice). Das Binary zeigt es: Bei 8 Bytes Kopf und bei 7 Bytes steht `sessions: []` mit `empty_sessions: true` in der Datei. Der Grund trägt. | `docs/plan/planning/in-progress/slice-v1-abschluss-sessions.md` §7 *Nacharbeit nach dem Review* | ohne Aktion |

Falsifikationsversuche, die nichts fanden: Verbindung ohne Anfrage und `connection` mit `--fail-on-unconsumed` (nichts gemeldet); `pg_isready` und TCP-Probe erzeugen im Record keine Session; mehr Verbindungen als Sessions; Wertemenge mit Leerer, Großschreibung, Unterstrich; Schlüssel in falschem Abschnitt und oben; Variable beim falschen Kommando; fehlende Kennzeichnung geht vor dem Öffnen des Ports (`--network none`, kein Listen-Fehler); Datei fehlt geht vor Kennzeichnung (Exit 3 vor 2).

---

## Verdikt

**Befunde:** 0 HIGH, 0 MEDIUM, 3 LOW (V-305, V-306, V-308), 2 INFO (V-307, V-309).

Die zwei Liefer-Punkte der DoD sind **bestätigt**: Die Tests existieren und laufen (`make test-integration` grün); 14 selbst gebaute Mutanten (drei Bedingungen des Zwangsendes im Aufbau je einzeln, verneinend „ohne Option nicht“ und „keine Meldung“, Wertemenge von `--session-assignment`, `id`-Zählung, `Unassigned`, Zuordnung nach `id`, Handshake, `PGR-E5003`, Bootstrap-Verdrahtung, `--record-empty-sessions` am Binary, halbe Startnachricht) sind aus dem richtigen Grund rot, keiner grün oder aus falschem Grund. Das Binary verhält sich wie spezifiziert: Startfehler ohne Datenbank (`PGR-E2001`, `PGR-E2004`), Kennzeichnung vor dem Öffnen des Ports, Nichtverbindungen sind keine Session, halbe Startnachricht Exit 6 und 7 Bytes still, `connection` mit leerer Session in der Mitte. Die Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando` ist erfüllt, der PGWire-Adapter ist unberührt, Handbuch und README sind unberührt, die Aufruf-Voraussetzung der Zuordnung steht in der Spezifikation, F-679 bis F-685 sind behoben oder entschieden eingetragen, der Plan folgt dem Diff.

**Closure-blockierend:** nichts im Code. Vor der Closure: V-308 (Planner trägt die Nahtstelle in die Nehmer ein); V-305 (Kommentar) und V-306 (Metrik, Stand der Größenprobe) ohne Blockade. Offen, weil die Closure sie führt: DoD-Häkchen, Review-Häkchen, Closure-Notiz mit Lerneintrag, Beobachtungs-Register (die Finding-Klassen des Reviews: Randform einer Erkennung per Zustand nicht genannt, beobachtete Randform nicht zurückgegeben, Zeitpunkt des Bemerkens, feste Wartezeit, Plan folgt der Korrektur nicht an allen Stellen, Verdrahtung/Logik im Bootstrap, Zusage im Kommentar ohne Test), die Ausgänge der zwei offenen Risiken aus §6 (beide erwartet *entfallen*: die Tests `TestReplayConnectionNteVerbindungNteSession` mit der Session nur aus Lebendprüfungen an mittlerer `id` und `TestReplayConnectionZuordnungMitDerAnnahme` mit `--fail-on-unconsumed` sind rot unter P1a bis P1c, P3a bis P3c), die drei Paarungen.

**Finding-Klasse dieses Laufs:** „Zusage im Kommentar weiter als der Test“ (V-305, enger Fall: eine Ausnahme der Zusage steht im Kommentar nicht).

**Gate:** `make gates` auf dem Baum nach dem Commit dieses Berichts: Ergebnis in der Übergabe-Nachricht (ein Gate-Lauf belegt nur den Stand, auf dem er lief); der Implementer und der Stop-Hook melden Exit 0 für `6df6976`, ich habe die volle Kette nicht gefahren.
