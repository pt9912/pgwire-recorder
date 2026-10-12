# Slice slice-v1-abschluss-antwortvergleich: Antwortvergleich beim Einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0020](../../adr/0020-antwortvergleich-beim-einspielen.md), [ADR-0023](../../adr/0023-antwortvergleich-entscheidung.md)

**Berührte Spec-Stellen:** `LH-FA-24.a` · `LH-FA-17.a` · `LH-FA-20.a` · `LH-FA-02.b` · `SPEC-018` · `SPEC-027` · `SPEC-034` · `SPEC-041` · `ARC-002`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `pgwire-recorder play --compare-responses` vergleicht die Struktur der Serverantworten und Fehler mit der Aufzeichnung und meldet Abweichungen als `PGR-E5004` mit Exit-Code 5.

**Gebunden durch `slice-v1-abschluss-umgebung-anderes-kommando`** (dort §7, Paarungen, offene Frage zu V-3; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): `--compare-responses` liest der Stand noch nicht; `nichtGelesen()` in `internal/adapters/driving/cli/fremd_test.go` führt sie. Dieser Slice liefert den Leser, kürzt `nichtGelesen()` um diese Option(en) und bindet ihn in `liesUmgebung` (`internal/adapters/driving/cli/cli.go`) an; die Zuordnung `zuordnung()` nennt sie für `play` (Spalte *Kommando* der Optionstabelle, `LH-FA-17.a`). Sonst wird `TestLeserUmgebungNichtGelesenerOptionen` rot, weil die Variable beim eigenen Kommando dann geprüft wird (`TestLeserZuordnungGleichTabelle` bleibt grün). Beim anderen Kommando (`record` und `replay`) ist die Variable schon jetzt `PGR-E2001`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-antwortvergleich-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- Vergleich von Zeilenwerten und Toleranzregeln für Felder — Out-of-Scope von LH-FA-24.
- Einspielen selbst und Fehlersemantik der Serverfehler — `slice-v1-abschluss-einspielen`; seit dessen zweitem Schnitt vom 2026-10-09 die drei Optionen der Laufsteuerung ohne Vergleich `slice-v1-abschluss-einspielen-laufsteuerung` und das Einspielen von Extended-Interaktionen `slice-v1-abschluss-einspielen-extended` (je §1, *Übernimmt*). Dieser Slice setzt alle drei voraus: DoD-Punkt 3 baut auf den Optionen auf, DoD-Punkt 1 und 2 vergleichen auch Extended-Interaktionen.
- Anmeldung und TLS beim Einspielen — `slice-v1-abschluss-einspielen-anmeldung` und `slice-v1-abschluss-einspielen-tls` (je §1, *Übernimmt*; seit dem Schnitt von `slice-v1-abschluss-einspielen` vom 2026-10-09); der Vergleich hängt an keinem der beiden.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen): Gegen eine Instanz mit gleicher Antwortstruktur endet das Einspielen mit Vergleich mit Erfolg, gegen eine abweichende mit `PGR-E5004` und Exit-Code 5; ein Serverfehler ohne Vorbild in der Aufzeichnung ist `PGR-E5004`, nicht `PGR-E4004`; eine Fehlerantwort `FATAL` oder `PANIC` ist `PGR-E5004`, auch bei gleichem SQLSTATE wie eine aufgezeichnete (die Aufzeichnung kennt kein Verbindungsende, `LH-FA-24.a` *Ende der Aufzeichnung*), ein Netzwerkabbruch ohne Fehlerantwort ist `PGR-E4003`; bei einer Fehlerantwort `FATAL`, `PANIC` oder einem Verbindungsverlust sendet `play` keine weiteren Nachrichten oder Gruppen der Interaktion, bei jeder anderen Abweichung wird die Interaktion bis zu ihrem `ReadyForQuery` verglichen und der Abbruch erfolgt vor der nächsten Interaktion (`LH-FA-20.a` Schritt 6); bei Extended wird nach einem Fehler bis zum `Sync` weitergesendet (Abnahmeszenario 16; Integrationstest).
- [ ] Jede Art der Abweichung (Nachrichtenart, Spaltenbeschreibung, Befehl ohne Zahlen, Fehler, Transaktionsstatus) wird erkannt und mit Session, Sequenznummer, Position und Art gemeldet (Position ab 1, bei fehlender oder zusätzlicher Nachricht die erste ohne Gegenstück, bei der Art Fehler der SQLSTATE beider Seiten und nie der Text, `LH-FA-24.a` *Randformen des Vergleichs*), Unterschiede nur in Zeilenwerten, Zeilenzahlen, Hinweisen und den abschließenden Zahlen des Befehls nicht (Test, einfach und Extended, auch `Flush`-Gruppen).
- [ ] `--continue-on-error` und `--allow-recorded-errors` wirken wie spezifiziert, `--continue-on-error` läuft nach einer Abweichung weiter, der Exit-Code am Ende ist 5, ein Verbindungsfehler beendet immer mit 4; ein Fehler mit gleichem SQLSTATE wie aufgezeichnet gilt als erwartet und bricht nicht ab, ein anderer SQLSTATE ist `PGR-E5004`; `--allow-recorded-errors` ist mit Vergleich ohne Wirkung; ein Abbruchsignal endet nach einer Abweichung mit 5, ein Verbindungsfehler nach einer Abweichung mit 4; bei jedem Ende des Laufs (auch Abbruch und Signal) trägt die Zeile zum Ende (`play beendet`, keine zweite Zeile) die Attribute eingespielt, verglichen, abweichend nach den Zählregeln von `LH-FA-24.a` und `LH-FA-20.a` Schritt 4; ohne `--compare-responses` findet kein Vergleich statt; der Leser der Option (Kommandozeile, Umgebungsvariable und Schlüssel `compare_responses` unter `play:`) ist angebunden und `nichtGelesen()` in `fremd_test.go` führt sie nicht mehr (aus `slice-v1-abschluss-umgebung-anderes-kommando`) (Test).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services` (Play-Service) | update | Normalisierung und Vergleich der Antworten, ohne PGWire-Typen |
| `internal/hexagon/model` | update | Katalogeintrag `PGR-E5004` (eine Konstante; zählt mit dem Play-Service als Kern, wie die Codes der früheren Play-Slices) |
| `internal/adapters/driving/cli` | update | Option `--compare-responses` (Kommandozeile, Umgebungsvariable, Schlüssel), Hilfetext, Abbildung auf Exit-Code 5 |
| `internal/bootstrap` | update | reicht die Zähler des Services an die Zeile zum Ende durch; reine Verdrahtung, keine Schicht (`AGENTS.md` §3.13), kein Lesen und kein Einstufen |
| `test/integration` | update | Happy/Boundary/Negative nach LH-FA-24 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen`, `slice-v1-abschluss-einspielen-laufsteuerung` und `slice-v1-abschluss-einspielen-extended` sind `done`, und [ADR-0020](../../adr/0020-antwortvergleich-beim-einspielen.md) ist `Accepted`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Normalisierung der Antworten erweist sich als größer als geplant — zurück zur Zerlegung. **Größenprobe** (`BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar`): Steht nach dem Commit, der DoD-Punkte 1 und 2 erfüllt, der Diff (ohne Plan, Spezifikation und erzeugte Tabellen) über 1400 Zeilen mit Tests, hält der Implementer an und meldet es dem Planner; Schnittlinie: DoD-Punkte 1 und 2 mit Option, Leser und Hilfetext | DoD-Punkt 3 (Laufsteuerung mit Vergleich, Zähler, Signal) als eigener Slice.
- `in-progress` → `open`: Die aufgezeichneten Antworten tragen zu wenig Felder für den Vergleich der Spaltenbeschreibung — Carveout.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Zusage, was nicht verglichen wird, ist schwer zu ändern; Treiber mit abweichenden Hinweisen oder Parameterstatus sind im Test zu belegen — **Ausgang:** offen bis Closure.
- Das Halten der Antworten bis zum Ende der Interaktion kostet Speicher bei großen Resultsets — **Ausgang:** offen bis Closure.
- Die Teile von `LH-FA-24.a` §Unvollständige Aufzeichnung, die aufgezeichnete Interaktionen ohne `ReadyForQuery` behandeln, sind seit der Entscheidung in `LH-FA-02.b` §Fehlerantwort vor dem Abbruch (`slice-replay-semantik-fehlerreplay`) nicht mehr erreichbar: Der Recorder schreibt keine solche Interaktion, und der Leser lehnt sie als beschädigt ab (`PGR-E3003`). Dieser Slice bereinigt sie in der Spezifikation, dazu in `LH-FA-20.a` den Teil der Tabellenzeile *Verbindung endet nach dem ersten `ReadyForQuery`* mit Vergleich, der ein aufgezeichnetes Ende behandelt, und den Satz „Wo die Aufzeichnung ein Verbindungsende zeigt, wartet der Recorder …“ (übergeben von `slice-v1-abschluss-einspielen`, das Schritt 3 bereinigt hat; die Klausel in `LH-FA-24` des Lastenhefts bleibt Sache des Nutzers); seine DoD sagt für sie nichts zu — **Ausgang:** entfallen: Der Architect hat die Spezifikation vor dem Code bereinigt (2026-10-12; `LH-FA-24.a` *Ende der Aufzeichnung*, Fehlerregeln und Katalogzeile in `LH-FA-20.a` und `SPEC-034`); der Implementer schreibt dazu nichts. Die Klauseln im Lastenheft bleiben Sache des Nutzers (Empfehlung des Architect: anpassen, siehe Bericht), sie sind ohne Wirkung auf diesen Slice, weil sie nur Fälle nennen, die keine Aufzeichnung enthält.

**Randformen** (`AGENTS.md` §3.12) — geprüft und entschieden vom Architect am 2026-10-12, vor dem Code, in der Spezifikation; der Implementer entscheidet keine:

| Randform | entschieden in |
|---|---|
| Variable beim eigenen und beim anderen Kommando (Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando`, §1): `--compare-responses` nur `play`; beim eigenen Kommando liest der Leser die Variable, beim anderen (`record`, `replay`) bleibt sie `PGR-E2001`; Test `TestLeserUmgebungNichtGelesenerOptionen` mit gekürztem `nichtGelesen()`, `TestLeserUmgebungAndererKommandos`. Wert der Option, Variable und Schlüssel: leer, `1`, `True` wie bei jeder booleschen Option | `LH-FA-17.a` (Tabelle, Absatz zu Umgebungsvariablen und Wertemenge, Konfigurationsdatei) |
| Normalisierung: `data_row`, `notice_response`, `parameter_status` entfallen; was je Nachricht verglichen wird | `LH-FA-24.a` (Tabelle) |
| Befehl mit und ohne Zahlen, Tag aus Ziffern allein | `LH-FA-24.a` *Randformen des Vergleichs*, *Befehl* |
| Fehlerantwort ohne Feld `C`, Schweregrad, `FATAL`/`PANIC` | `LH-FA-24.a` *Randformen des Vergleichs*, *Fehler*, und *Ende der Aufzeichnung* |
| Fehlende oder zusätzliche Nachricht, Position (ab 1), Art | `LH-FA-24.a` *Randformen des Vergleichs*, *Position*, *Art* |
| Rückgabe des Implementers (Architect, 2026-10-12): `FATAL`/`PANIC` ohne `ErrorResponse` an der Stelle der Aufzeichnung: Art Fehler, aufgezeichneter SQLSTATE leer, eine frühere Abweichung geht vor (Ausnahme von *Art*). Code `erste()` und `TestVergleichFatalRand` bleiben, keine Nacharbeit | `LH-FA-24.a` *Ende der Aufzeichnung*, *Art* |
| Rückgabe des Implementers (Architect, 2026-10-12): Position ist die erste Position mit verschiedener Nachricht; „Länge der kürzeren Folge plus 1“ nur bei gleichem Anfang. Code `erste()` und `TestVergleichStruktur` bleiben, keine Nacharbeit; ein Testfall mit fehlender Nachricht in der Mitte (Position 1) wäre eine Ergänzung nach Ermessen des Implementers, kein Muss. `PGR-E5004` in `SPEC-034` und `LH-FA-24` nennen die Position nicht anders | `LH-FA-24.a` *Randformen des Vergleichs*, *Position* |
| Inhalt der Meldung: SQLSTATE beider Seiten ja, Text und Felder nie, Werte nie | `LH-FA-24.a` *Meldung* |
| Extended: Folge aller Server-Nachrichten ohne Gruppengrenzen; Fehler mitten in der Interaktion; Abbruch erst am Ende der Interaktion, außer `FATAL`/`PANIC` | `LH-FA-24.a` (Einleitung), `LH-FA-20.a` Schritt 6 und *Gruppen* |
| Aufgezeichnete Interaktion ohne `ReadyForQuery`, aufgezeichnetes Verbindungsende: gibt es nicht (Leser lehnt ab) | `LH-FA-24.a` *Ende der Aufzeichnung*, `LH-FA-02.b` *Replay* |
| Verbindungsverlust ohne Fehlerantwort: `PGR-E4003`, nach einer Abweichung ebenfalls 4 | `LH-FA-20.a` (Tabelle, Absatz darunter) |
| Abbruch und Herunterfahren: erstes Signal (Interaktion wird verglichen), zweites Signal (ohne Vergleich), `--finish-session-on-interrupt` | `LH-FA-20.a` Schritt 7 und *Abbruchsignal*, `LH-FA-24.a` *Randformen des Vergleichs*, *Signal* |
| Zähler: wann eine Interaktion `eingespielt`, `verglichen`, `abweichend` zählt; die Zeile zum Ende trägt sie, keine zweite Zeile; Stufe `info`; Lauf ohne Interaktion; ohne Option alle `verglichen`/`abweichend` 0 | `LH-FA-24.a` *Meldung*, `LH-FA-20.a` Schritt 4 und *Meldungen* |
| Reihenfolge der Zeilen `error` bei Abweichungen: je Interaktion eine Zeile, in der Reihenfolge des Auftretens, abbrechende zuerst | `LH-FA-20.a` *Meldungen*, *Exit-Code* |
| `--allow-recorded-errors` mit Vergleich ohne Wirkung, keine Fehlermeldung für die Kombination | `LH-FA-20.a` (Optionstabelle, Fehlertabelle) |
| Verhalten je Serverversion: kein Ausgleich | `LH-FA-24.a` *Randformen des Vergleichs*, *Serverversion* |
| Aufruf-Voraussetzung: keine Ordnungszusage über mehrere Aufrufe; Zähler gelten für einen Lauf mit nacheinander eingespielten Interaktionen | `LH-FA-24.a` *Randformen des Vergleichs*, *Aufruf-Voraussetzung* |

**Hilfetext** (`AGENTS.md` §3.11): Der Hilfetext von `--compare-responses` sagt nur zu, was ein Test dieses Slice prüft (Vergleich der Antwortstruktur, Exit-Code 5); die Wirkung auf `--continue-on-error` und `--allow-recorded-errors` steht im Hilfetext nur, wenn ein Test sie prüft.

**Akzeptierte Negative** (Grund im Satz, `AGENTS.md` §3.12): (a) Speicher bei großen Resultsets (Risiko oben): `data_row` entfällt bei der Normalisierung, der Service muss Zeilen beim Lesen nicht halten; ob er es tut, prüft kein Test, der Plan sagt es nicht zu. (b) Der Vergleich über Serverversionen hinweg (Tag, OID) ist eine Abweichung, kein Fehler des Vergleichs; kein Test mit zwei Versionen in diesem Slice — `slice-v1-abschluss-postgres-versionen` fährt mehrere Versionen, ist aber keine Adresse dafür (kein Eintrag dort).

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

### Belege des Implementers

**Stand.** Commits `6ef3cf6` (Kern, Leser, Zähler, Tests, Abdeckungstabellen), `79a380c` (Folgen verschiedener Länge mit gleichem Anfang), `b0122fd` (Kommentar von `Zaehler()`), `25ccf77` und `8e8b60d` (Integrationstests Extended-Fehler und Hinweise, Abdeckungstabelle). Geändert: `internal/hexagon/services/play.go`, neu `play_vergleich.go`, `internal/hexagon/model/fehler.go` (`CodeResponseMismatch`), `internal/adapters/driving/cli/cli.go` (Option, Hilfetext), `internal/bootstrap/bootstrap.go` (Zähler in der Zeile `play beendet`), Tests in `services`, `cli`, `bootstrap` und `test/integration`. Handbuch und `README.md` unverändert (`grep -i compare` in beiden ohne Treffer). `liesUmgebung` und `zuordnung()` sind unverändert: Die Option steht in `optionen("play")` an der Stelle der Tabelle (hinter `upstream-ca`), `zuordnung()` nannte sie schon.

**Größe (§4).** Gegenüber `08787ed`, ohne `docs/` und `spec/`: 1441 Zeilen hinzugefügt, 54 entfernt. Einen Commit nur für DoD-Punkte 1 und 2 gab es nicht (Kern, Zähler und Option sind in `6ef3cf6` zusammen). Der Anteil von DoD-Punkt 3 (Fortsetzung, Zähler, Signal, Option, Leser, Zeile `play beendet` und ihre Tests) ist etwa 500 Zeilen: 274 Zeilen sind die Funktionen seiner Tests in `play_vergleich_test.go` und `play_vergleich_e2e_test.go` (ohne ihre Kommentare, dazu etwa 80), dazu Option, Leser, Hilfetext, die Tests in `cli` und `bootstrap` und die Zähler im Service. Für DoD-Punkte 1 und 2 bleiben etwa 940. Die Probe von 1400 Zeilen ist damit nicht eingetreten; kein Halt und kein Hinweis an den Planner aus diesem Grund (der Gesamtstand liegt unter der Schätzung des Architect in §8).

**Läufe** (Stand `8e8b60d`, Arbeitsbaum sauber).

| Sensor | Ergebnis |
|---|---|
| `make test` (gofmt, vet, Unit-Tests) | grün; alle Pakete `ok` |
| `make lint` | `0 issues.` |
| `make a-check` | `gesamt: 0 Befund(e)` |
| `make gates` (07:08 min, Lauf nach `8e8b60d`) | Exit 0; darin `abdeckung-check`, `kopf-check`, `d-check` (636 Dateien, 0 Befunde), `baseline-verify`, `a-check-negativ`, `lint-gegenprobe`, `run-integration-tests: gruen` (beide Phasen) |
| Integrationstests von `play --compare-responses` allein (`TestE2EPlayVergleich*`, zehn Tests) | grün gegen `postgres:17-alpine` |

**Mutationen (Weg).** Jeder Mutant ändert eine Stelle in `play.go`, `play_vergleich.go`, `cli.go`, `bootstrap.go`, `fehler.go` oder `fremd_test.go`, das Testimage `pgwire-recorder:test` läuft mit dem Arbeitsbaum als Bind-Mount (`-v "$PWD":/src`, kein Build-Kontext, kein mtime-Problem) und `go test` über die Pakete des Mutanten; danach `git checkout` der Datei. Alle 70 Mutanten (Werkzeug im Scratch-Verzeichnis) sind **rot**, keiner grün; Stand `b0122fd` (Quellcode wie `8e8b60d`). Die Integrationstests liefen für `C2`, `M30`, `M33a`, `M1`, `M31`, `M3`, `M4`, `M5` nach `touch` der Datei im Image-Build (jeder rot, gelistet unten). Die Spalte „roter Test“ nennt aus der Ausgabe des Laufs die Tests (bei Tabellenfällen den Fall), an denen die Zusage hängt; rot sind bei vielen Mutanten weitere Tests.

| # | Zusage | Mutation | roter Test |
|---|---|---|---|
| M1 | ohne `--compare-responses` findet kein Vergleich statt | Vergleich immer an (`if true`) | `TestVergleichOhneOption`; E2E `TestE2EPlayVergleichQuellen` |
| M55 | ohne Option `verglichen` und `abweichend` 0 | `verglichen` zählt auch ohne Vergleich | `TestVergleichOhneOption` |
| M2 | Position zählt ab 1 | Position ab 0 | `TestVergleichStruktur` (jeder Fall mit Abweichung) |
| M3 | `data_row` entfällt vor dem Vergleich | `data_row` bleibt in der Folge | `TestVergleichStruktur/Zeilenwerte_und_Zeilenzahl`, `…/keine_Zeile`; E2E `TestE2EPlayVergleichGleich` |
| M4 | `notice_response` entfällt | bleibt in der Folge | `TestVergleichStruktur/Hinweise_und_Parameterstatus_*`; E2E `TestE2EPlayVergleichHinweise` |
| M5 | `parameter_status` entfällt | bleibt in der Folge | `TestVergleichStruktur/Hinweise_und_Parameterstatus_*`; E2E `TestE2EPlayVergleichHinweise` |
| M50 | bei Extended zählt die Folge aller Gruppen ohne Grenzen | nur die Gruppe am Ende | `TestVergleichExtendedGruppen` |
| M22 | der Nachrichtentyp wird verglichen (jede Nachricht) | Typvergleich aus | `TestVergleichStruktur/Nachrichtenart`, `…/fehlende_Nachricht` |
| M56 | verschieden lange Folgen mit gleichem Anfang: Position hinter der kürzeren | Zweig entfernt | `TestVergleichStruktur/Aufzeichnung_länger,_Anfang_gleich` |
| M12 | `row_description`: Anzahl der Spalten | nur die gemeinsamen Spalten | `TestVergleichStruktur/Spaltenzahl` |
| M10 | … je Spalte der Name | Name nicht verglichen | `TestVergleichStruktur/Spaltenname` |
| M11 | … je Spalte die Typ-OID | OID nicht verglichen | `TestVergleichStruktur/Typ-OID_der_Spalte`, `…/zweite_Spalte` |
| M13 | … jede Spalte, nicht nur die erste | nur die erste Spalte | `TestVergleichStruktur/zweite_Spalte` |
| M14 | … übrige Felder der Spalte nicht (Tabelle, Format, Größe) | ganze Spalte verglichen | `TestVergleichStruktur/übrige_Felder_der_Spalte` |
| M15 | `parameter_description`: die Liste der Typ-OIDs | nur ihre Länge | `TestVergleichStruktur/Parameterbeschreibung,_OID` (Extended), `TestVergleichExtendedGruppen` |
| M6 | `command_complete`: Tag ohne die abschließenden Zahlen | Tag ganz verglichen | `TestVergleichStruktur/Zahlen_am_Ende_des_Befehls`, `…/Befehl_mit_und_ohne_Zahl` |
| M7 | … nur Wörter aus Ziffern sind Zahlen | `TrimRight` über Ziffern und Leerzeichen | `TestVergleichStruktur/Wort_mit_Ziffer_ist_keine_Zahl` |
| M8 | … ein Tag aus Ziffern allein gilt als leer | mindestens ein Wort bleibt | `TestVergleichStruktur/Tag_aus_Ziffern_allein` |
| M9 | … der ganze Befehl, nicht nur das erste Wort | nur das erste Wort | `TestVergleichStruktur/Wort_des_Befehls`, `…/Befehl_kürzer_um_ein_Wort` |
| M16 | `error_response`: der SQLSTATE, nicht der Text | `M` statt `C` verglichen | `TestVergleichStruktur/gleicher_SQLSTATE,_anderer_Text_und_Felder`, `TestVergleichErwarteterFehler` |
| M17 | … weitere Felder nicht | alle Felder verglichen | `TestVergleichStruktur/gleicher_SQLSTATE,_anderer_Text_und_Felder` |
| M18 | … der Schweregrad nicht | `V` mit verglichen | `TestVergleichStruktur/anderer_Schweregrad,_gleicher_SQLSTATE` |
| M19 | … der Text nicht | `M` mit verglichen | `TestVergleichStruktur/gleicher_SQLSTATE,_anderer_Text_und_Felder` |
| M20 | … zwei Fehlerantworten ohne `C` sind gleich | ohne `C` immer verschieden | `TestVergleichStruktur/Fehler_ohne_Feld_C_auf_beiden_Seiten` |
| M21 | `ready_for_query`: der Transaktionsstatus | nicht verglichen | `TestVergleichStruktur/Transaktionsstatus`, `TestVergleichFortsetzung` |
| M23a | Meldung bei der Art Fehler: SQLSTATE beider Seiten, richtig zugeordnet | beide vertauscht | `TestVergleichStruktur/anderer_SQLSTATE`, `TestVergleichErwarteterFehler` |
| M23b | … der SQLSTATE beider Seiten | nur der des Servers | `TestVergleichStruktur/anderer_SQLSTATE`, `TestVergleichFatal` |
| M24 | Text und weitere Felder erscheinen nie in der Meldung | `M` statt `C` in der Meldung | `TestVergleichStruktur/anderer_SQLSTATE`, `TestVergleichErwarteterFehler` |
| M25 | je Interaktion die erste Abweichung | die letzte Position gemeldet | `TestVergleichStruktur/erste_von_mehreren_Abweichungen`, `TestVergleichFatalRand` |
| M26 | bei jeder Abweichung außer `FATAL`/`PANIC` wird bis zum `ReadyForQuery` gelesen | Fehler bei der ersten abweichenden Nachricht | `TestVergleichAbbruchNachInteraktion` (zählt die `naechste` im Ablauf), `TestVergleichExtendedFehlerMitten` |
| M27 | der Abbruch erfolgt vor der nächsten Interaktion | Lauf geht ohne `--continue-on-error` weiter | `TestVergleichAbbruchNachInteraktion` |
| M28 | `--continue-on-error` läuft nach einer Abweichung weiter | bricht immer ab | `TestVergleichFortsetzung`, `TestVergleichZweitesSignal` |
| M29 | die Abweichung nach `--continue-on-error` ist eine Meldung, Exit-Code 5 | Abweichung verschluckt | `TestVergleichFortsetzung` (drei Meldungen in der Reihenfolge), `TestVergleichErstesSignal` |
| M30 | ein Serverfehler ohne Vorbild ist `PGR-E5004`, nicht `PGR-E4004` | Vergleich-Zweig der Fehlerantwort aus | `TestVergleichFehlerOhneVorbild`; E2E `TestE2EPlayVergleichAbweichung/Serverfehler_ohne_Vorbild` |
| M31 | ein Fehler mit gleichem SQLSTATE wie aufgezeichnet gilt als erwartet und bricht nicht ab | jeder Fehler bricht beim Lesen ab | `TestVergleichErwarteterFehler`, `TestVergleichExtendedFehlerMitten`; E2E `TestE2EPlayVergleichGleich`, `TestE2EPlayVergleichExtendedFehler` |
| M32 | `--allow-recorded-errors` ist mit Vergleich ohne Wirkung | erwarteter Fehler gilt auch mit Vergleich | `TestVergleichErwarteterFehler` (anderer SQLSTATE bei aufgezeichnetem Fehler) |
| M33a | `FATAL`/`PANIC` mit Vergleich ist `PGR-E5004`, nicht `PGR-E4003` | `PGR-E4003` | `TestVergleichFatal`, `TestVergleichFatalExtendedGruppen`; E2E `TestE2EPlayVergleichFatal` |
| M34 | … auch bei gleichem SQLSTATE wie die aufgezeichnete | gleicher SQLSTATE ist keine Abweichung | `TestVergleichFatal` |
| M33b | … sofort | Fehlerantwort geht in den Vergleich am Ende ein | `TestVergleichFatal` |
| M35 | … play liest nach ihr nichts mehr (verneinend) | liest bis zum `ReadyForQuery` weiter | `TestVergleichFatal` (der Server liefert danach zwei Nachrichten; der Ablauf zählt vier `naechste`), `TestVergleichFatalExtendedGruppen` |
| M52 | … in der Antwort auf eine Gruppe mit `Flush` sendet play keine weitere Gruppe (verneinend) | Fehlerantwort im Warten verschluckt | `TestVergleichFatalExtendedGruppen` |
| M36a | die Interaktion zählt als verglichen | zählt nicht | `TestVergleichFatal`, `TestVergleichFatalFortsetzung` |
| M36b | … und als abweichend | zählt nicht | `TestVergleichFatal`, `TestVergleichFatalFortsetzung` |
| M37 | mit `--continue-on-error` endet die Session, die nächste läuft | die Session sendet die nächste Interaktion | `TestVergleichFatalFortsetzung` |
| M38 | … die Session endet ohne Fehler | `errSessionEnde` wird als Fehler gemeldet | `TestVergleichFatalFortsetzung` |
| M39 | … der Lauf bricht nicht ab | `PGR-E5004` bricht ab | `TestVergleichFatalFortsetzung` |
| M40 | eine frühere Abweichung der Interaktion geht der Fehlerantwort vor | Fehlerantwort immer gemeldet | `TestVergleichFatalRand` |
| M41 | die Position der Fehlerantwort in der normalisierten Folge | Position minus eins | `TestVergleichFatal`, `TestVergleichFatalFortsetzung` |
| M42 | der Schweregrad kommt aus `V`, ohne `V` aus `S` | nur `S` | `TestVergleichFatal` (V FATAL bei S SCHWER), `TestVergleichFatalRand` (S FATAL, V ERROR), `TestPlayFehlerantwort` |
| M43 | `PANIC` wie `FATAL` | `PANIC` entfernt | `TestVergleichFatal` (V PANIC, S PANIC) |
| M44 | ein Verbindungsverlust ohne Fehlerantwort ist `PGR-E4003`, keine Abweichung | `PGR-E5004` | `TestVergleichVerbindungsverlust` |
| M46 | … auch mit `--continue-on-error` sofort, Exit-Code 4 nach einer früheren Abweichung | Verlust wird mit der Option verschluckt | `TestVergleichVerbindungsverlust`, `TestVergleichFatalFortsetzung` |
| M45 | eine Interaktion, die ein Verbindungsverlust oder das zweite Signal unterbricht, zählt nicht als verglichen (verneinend) | `verglichen` zählt am Beginn | `TestVergleichVerbindungsverlust`, `TestVergleichZweitesSignal` |
| M54 | das zweite Signal unterbricht ohne Vergleich (verneinend) | die gelesenen Antworten werden verglichen und gezählt | `TestVergleichZweitesSignal`, `TestVergleichVerbindungsverlust` |
| M53 | das erste Signal beendet die Interaktion mit Vergleich, eine Abweichung darin wird festgestellt | Fehler nach dem ersten Signal verschluckt | `TestVergleichErstesSignal` |
| M47 | `eingespielt` zählt ab dem Beginn des Sendens, auch wenn es scheitert | zählt nur, wenn die Interaktion ohne Fehler endet | `TestVergleichZaehlerEingespielt` (einfach und Extended, ohne und mit Vergleich) |
| M48 | `abweichend` zählt nur Interaktionen mit Abweichung, höchstens eine je Interaktion | zählt jede verglichene | `TestVergleichStruktur` (Zähler jedes Falls), `TestVergleichFortsetzung` |
| M49 | `verglichen` zählt am Ende der Interaktion | zählt nie | `TestVergleichStruktur`, `TestVergleichAbbruchNachInteraktion` |
| M51 | nach einer Fehlerantwort mitten in einer Extended-Interaktion werden die übrigen Gruppen ohne Warten gesendet | `fehlerGesehen` entfernt | `TestVergleichExtendedFehlerOhneWarten`, `TestPlayExtendedFortsetzung` |
| C1 | Leser: die Option steht in `optionen("play")` | Eintrag entfernt | `TestLeserOptionen`, `TestLeserZuordnungGleichTabelle`, `TestParsePlayLaufsteuerung` |
| C2 | Leser: Kommandozeile, Umgebungsvariable und Schlüssel `compare_responses` setzen das Feld | Feld bleibt `false` | `TestParsePlayLaufsteuerung/compare-responses`, `TestLeserAlleOptionen`; E2E `TestE2EPlayVergleichQuellen`, `…Gleich` |
| C3 | Leser: Standard `false` | `true` | `TestLeserAlleOptionen`, `TestParsePlay` |
| C6 | Option gehört zu `play` (Spalte *Kommando*): bei `record` und `replay` ist die Variable `PGR-E2001` | Zuordnung auf `record` und `replay` | `TestLeserAlleOptionen`, `TestLeserUmgebungAndererKommandos` |
| C7 | Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando`: `nichtGelesen()` führt die Option nicht, die Variable wird geprüft | `compare-responses` wieder in `nichtGelesen()` | `TestLeserUmgebungNichtGelesenerOptionen` (der Test verlangt die Variable einer Option in `nichtGelesen()` als unbeachtet, auch mit ungültigem Wert; der Leser prüft sie) |
| C4 | Hilfetext nennt `--compare-responses` und die Variable | Option im Text umbenannt | `TestParsePlayHilfe` |
| C5 | Hilfetext: Abweichung ist `PGR-E5004`, Exit-Code 5 | Exit-Code 4 im Text | `TestParsePlayHilfe` |
| C8 | die Zeile `play beendet` trägt `verglichen` aus dem Zähler | `eingespielt` an dieser Stelle | `TestRunPlayVergleich` |
| C10 | … `abweichend` aus dem Zähler | konstant 0 | `TestRunPlayVergleich` (Fall „abweichende Struktur“) |
| C9 | … es entsteht keine zweite Zeile | zweite Zeile `play beendet` | `TestRunPlayVergleich`, `TestRunPlay` |
| C11 | der Katalog nennt `PGR-E5004` | anderer Code | `TestRunPlayVergleich`, `TestVergleichStruktur` |

**Nicht erreichbar.** Eine aufgezeichnete Folge, die vor der Position der Fehlerantwort `FATAL` endet, gibt es nicht: Sie endet in `ready_for_query`, das der Server vorher nicht schickte, die Folgen weichen dort schon ab; der Code hat den Zweig nicht. Der Zweig *Folgen verschiedener Länge mit gleichem Anfang* ist nur über eine einfache Aufzeichnung mit einem `ready_for_query` vor dem Ende erreichbar (`Validate` verbietet das nur bei Extended); `TestVergleichStruktur/Aufzeichnung_länger,_Anfang_gleich` fährt ihn (M56).

**Geänderte bestehende Tests.** `fremd_test.go`: `nichtGelesen()` ohne `compare-responses` (Befund vorher: die Variable blieb beim eigenen Kommando unbeachtet; nachher: `TestLeserUmgebungNichtGelesenerOptionen` grün mit gekürzter Liste, rot mit der Option darin, C7). `play_test.go`: `TestParsePlayFremdeUmgebung` verliert den Teil zu einer Option, die `play` kennt und der Stand nicht liest (es gibt keine mehr); `TestParsePlayLaufsteuerung` führt die vierte Option, `TestParsePlayHilfe` ihren Hilfetext. `leser_test.go`: Standard und Optionsliste von `play`. `bootstrap/play_test.go`: `leerePlay` setzt die neue Variable leer.

**Hilfetext (`AGENTS.md` §3.11).** Jeder Satz hat einen Test: *vergleicht die Struktur der Serverantworten jeder Interaktion mit der Aufzeichnung* (`TestVergleichStruktur`, `TestE2EPlayVergleichGleich`), *eine Abweichung ist PGR-E5004 (Exit-Code 5)* (`TestParsePlayHilfe` für den Text; `TestVergleichStruktur`, `TestRunPlayVergleich`, `TestE2EPlayVergleichAbweichung` für Code und Exit-Code), *Standard false* (`TestLeserAlleOptionen`), *Umgebungsvariable* (`TestParsePlayLaufsteuerung`, `TestE2EPlayVergleichQuellen`). Die Wirkung auf `--continue-on-error` und `--allow-recorded-errors` steht nicht im Hilfetext.

**Akzeptierte Negative (§6).** (a) Der Service hält von der Folge des Servers nur die normalisierten Nachrichten (`data_row` kommt gar nicht in die Folge, `nimmNormalisiert`); kein Test, keine Zusage. (b) Kein Test mit zwei Serverversionen, wie dort. Beide unverändert.

**Risiken (§6) für die Closure.** *Treiber mit abweichenden Hinweisen oder Parameterstatus*: gegen den realen Server belegt für beide Richtungen und beide Nachrichten (`TestE2EPlayVergleichHinweise`, Aufzeichnung mit `SET application_name` und `RAISE NOTICE`; dazu `pgx` im Extended-Modus in `TestE2EPlayVergleichExtended`); ob das für andere Treiber gilt, belegt der Test nicht. *Speicher bei großen Resultsets*: siehe (a), kein Test. *Aufgezeichnete Teile ohne `ReadyForQuery`*: entfallen (Architect).

**Rückgaben und Hinweise.**

1. **An den Architect, zur Bestätigung (Lesart in den Commits).** `LH-FA-24.a` *Ende der Aufzeichnung* sagt: `FATAL`/`PANIC` ist „stets eine Abweichung der Art Fehler“; *Art* sagt: bei verschiedenem Typ an einer Position ist die Art Nachrichtenart. Der Code liest den Satz zu `FATAL`/`PANIC` als den besonderen: Hat die Aufzeichnung an der Position keine `ErrorResponse` (etwa `command_complete`), ist die Art trotzdem Fehler, der SQLSTATE der Aufzeichnung ist leer („aufgezeichnet „“, Server „57P01““); eine frühere Abweichung der Interaktion geht vor. Festgehalten in `TestVergleichFatalRand` (Fall „FATAL an der Stelle von command_complete“). Wählt die Spezifikation anders, ändert sich `erste()` in `play_vergleich.go` um zwei Zeilen und dieser Fall.
2. **An den Architect, zur Bestätigung.** `LH-FA-24.a` *Position*: „die Länge der kürzeren Folge plus 1“. Der Code meldet die Position der ersten Stelle, an der die Nachrichten verschieden sind (Stück für Stück, wie die Einleitung sagt); die Formel stimmt damit nur bei gleichem Anfang (`TestVergleichStruktur/Aufzeichnung_länger,_Anfang_gleich`, Position 3). Bei fehlender Nachricht in der Mitte (`TestVergleichStruktur/fehlende_Nachricht`: Aufzeichnung `row_description`, `command_complete`, `ready_for_query`, Server ohne `row_description`) meldet er Position 1, nicht Länge plus 1.
3. **An den Nutzer (Lastenheft, kein Befund gegen den Slice).** Abnahmeszenario 16, dritter Satz („Endet die Aufzeichnung mit einem Verbindungsende, das der Server genauso verursacht, gilt es als erwartet“), ist nach der Bereinigung der Spezifikation unerreichbar: Die Aufzeichnung kennt kein Verbindungsende, ein Verbindungsverlust ohne Fehlerantwort ist `PGR-E4003`, `FATAL`/`PANIC` stets `PGR-E5004`. Er ist durch keinen Test belegt und nicht belegbar; das Lastenheft ist nicht geändert.
4. **Für Review und Verifikation.** `errSessionEnde` ist ein Wert im Service, den `session()` in `nil` verwandelt; er verlässt den Service nicht. Der Service zählt `eingespielt` auch ohne Vergleich (DoD: Zeile `play beendet` trägt die Zähler immer); ohne Vergleich sind `verglichen` und `abweichend` 0.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-antwortvergleich-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 3. Schichten: zwei Schichten (Kern mit Play-Service, CLI-Adapter).

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando` (Planner, 2026-10-12, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 3 (Leser der Option und Kürzung von `nichtGelesen()` im Test des CLI-Adapters) und fügt keinen Liefer-Punkt hinzu: Liefer-Punkte 3. Schichten: zwei (Kern mit Play-Service, CLI-Adapter); der CLI-Adapter war schon gezählt, `fremd_test.go` liegt in ihm.

Nachgezählt bei der Prüfung des Architect (2026-10-12): Liefer-Punkte 3 (DoD 1 bis 3). Schichten zwei: Kern (Play-Service mit Katalog) und CLI-Adapter; der Bootstrap reicht nur die Zähler durch (Verdrahtung, `AGENTS.md` §3.13, keine Schicht; liest er oder stuft er ein, zählt er, und der Plan wird neu geprüft). Die Sendung aus `slice-v1-abschluss-umgebung-anderes-kommando` liegt in DoD-Punkt 3 und ist nachgezählt. **Größe:** Schätzung des Architect 1800 bis 2300 Zeilen mit Tests und Abdeckungstabelle (Produkt-Code etwa 450 bis 550: Normalisierung und Vergleich 200, Einbindung in `interaktion`/`extended` 120, Zähler und Verdrahtung 80, Option, Hilfetext und Katalog 60; Tests der Zusagen je Mutation, Unit und Integration, 1300 bis 1700); der Slice liegt an der Obergrenze einer Sitzung, nicht darüber: Er liegt unter den 2200 bis 3000 Zeilen der beiden Belege von `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` und über `slice-v1-abschluss-einspielen-laufsteuerung` (etwa 900). Kein dritter Beleg; die Größenprobe in §4 fängt eine Fehlschätzung vor dem Rest der Arbeit. Tritt sie ein, ist das der dritte Beleg, den die Closure einträgt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
