# Slice slice-v1-abschluss-sessions: Mehrere Sessions und Verbindungsfehler

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-09`](../../../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay), [ADR-0004](../../adr/0004-postgresql-upstream-ist-driven-adapter.md), [ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md)

**Berührte Spec-Stellen:** `LH-FA-12.a` · `LH-FA-17.a` · `LH-FA-09.a` · `LH-FA-03.b` · `LH-FA-13.b` · `LH-FA-13.a` · `LH-FA-10.a` · `LH-FA-25.a` · `LH-FA-02.b` · `LH-FA-21.a` · `SPEC-017` · `SPEC-028` · `SPEC-034` · `SPEC-043` · `ARC-002`

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

**Bereits geliefert** vom Walking Skeleton: gleichzeitige Verbindungen im Record, Übernahme in der Reihenfolge der Sessionenden, Zuordnung `first-request` im Replay. Dieser Slice ergänzt `--record-empty-sessions`, `--session-assignment connection` und belegt LH-FA-12 vollständig. Von `slice-extended-query-lebendpruefung` ([ADR-0031](../../adr/0031-lebendpruefungen-im-replay.md), `LH-FA-09.a`): Lebendprüfungen (Anfragen nur aus Leerraum und Kommentaren) beantwortet das Replay außerhalb der Reihe; sie lösen keine Zuordnung aus, auch nicht, wenn keine Session mehr frei ist, und eine Verbindung, die nur Lebendprüfungen stellt, ist eine Verbindung ohne Anfrage. Aufgezeichnete Lebendprüfungen überspringt das Replay; eine Session nur aus ihnen (ohne Ende des Servers, `LH-FA-09.a`) ist für das Replay eine Session ohne Interaktion und bei `first-request` übersprungen.

**Kopplung an `--session-assignment connection`:** `NewReplayService` nimmt Sessions ohne Interaktion, auch solche nur aus Lebendprüfungen, schon beim Laden aus der Liste der freien Sessions (`frei`). `LH-FA-12.a` verlangt für `connection` dagegen, dass die n-te Verbindung die Session mit der `id` n erhält, auch eine Session ohne Interaktion und eine nur aus Lebendprüfungen. Die Zuordnung nach `id` muss diese Sessions deshalb weiter führen; die Liste `frei` reicht dafür nicht. Für eingehende Lebendprüfungen gilt bei `connection` vom Verbindungsaufbau an die Serverversion der zugeordneten Session (`server_version`, `LH-FA-09.a`).

**Übernommen aus `slice-replay-semantik-mismatch`:** Eine Session ohne Interaktion, auch eine nur aus Lebendprüfungen, ist mit und ohne `--fail-on-unconsumed` weder nicht verbraucht noch nie zugeordnet, auch bei `connection` (`LH-FA-03.b` §Verbraucht). Für `first-request` prüft das `TestReplayNichtsZuMelden`; für `connection` gehört der Test hierher, weil die Zuordnung nach `id` diese Sessions weiter führt und `Unassigned` sie dann nicht als nie zugeordnet zählen darf.

**Übernommen von `slice-v1-abschluss-sqlite-format`** (dort §1, *Abgegeben*, und DoD; Verifikation V-174): Die Wirkung von `--record-empty-sessions` mit `--format sqlite` am Binary. Die Option gibt es erst mit diesem Slice, und `meta.empty_sessions` des Schreibens von `sqlite` (`SPEC-043`) ist bis dahin nur im Adapter-Test erreichbar. Der Test steht in DoD-Punkt 1; geschrieben wird mit `--format sqlite`, gelesen über `replay` (`slice-v1-abschluss-sqlite-format-lesen` liegt in der Reihenfolge davor).

**Gebunden durch `slice-v1-abschluss-zeitangaben-aufzeichnen`** (dort §6 *Nicht absteigend* und §7; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13):

1. Die Option `--record-empty-sessions` und ihr Schlüssel `record_empty_sessions` entstehen mit diesem Slice. Der Fall „Option, die der Stand nicht kennt“ in `TestDateiUngueltig` (`internal/adapters/driving/cli/datei_test.go`) benutzt heute den Schlüssel `record_empty_sessions` als Platzhalter einer unbekannten Option und wird mit der Lieferung laut rot. Dieser Slice stellt den Fall auf einen Schlüssel um, den auch spätere Slices nicht liefern, und führt `record_empty_sessions` mit eigenen Fällen (gültiger Wert, Ebene `replay:`/`play:`/oben, ungültiger Wert nach der Optionstabelle), wie `slice-v1-abschluss-zeitangaben-aufzeichnen` es für `record_timing` tat. Zum Test in DoD-Punkt 1.
2. Die Zusage „Die Werte einer Session steigen nicht ab“ (`LH-FA-21.a` *Messpunkt*, `--record-timing`) gilt nur, wenn die Aufrufe der Client-Richtung einer Session (`Query`, `ClientMessage`) nacheinander kommen. Nebenläufig werden die Verbindungen (je eine Session mit eigener Schleife), nicht die Aufrufe einer Session. Dieser Slice hält das ein; kein Test prüft es (akzeptiertes Negativ im Geber, ein Adapter-Test wäre dort eine dritte Schicht). Zeigt sich im Code, dass die Nebenläufigkeit Aufrufe einer Session überlappen lässt, ist das ein Plan-Wechsel (§4).

**Gebunden durch `slice-v1-abschluss-umgebung-anderes-kommando`** (dort §7, Paarungen, offene Frage zu V-3; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): `--session-assignment` (bei `replay`) und `--record-empty-sessions` (bei `record`) liest der Stand noch nicht; `nichtGelesen()` in `internal/adapters/driving/cli/fremd_test.go` führt sie. Dieser Slice liefert den Leser, kürzt `nichtGelesen()` um diese Option(en) und bindet ihn in `liesUmgebung` (`internal/adapters/driving/cli/cli.go`) an; die Zuordnung `zuordnung()` nennt sie für `replay` (`--session-assignment`) und `record` (`--record-empty-sessions`) (Spalte *Kommando* der Optionstabelle, `LH-FA-17.a`). Sonst wird `TestLeserUmgebungNichtGelesenerOptionen` rot, weil die Variable beim eigenen Kommando dann geprüft wird (`TestLeserZuordnungGleichTabelle` bleibt grün). Beim anderen Kommando (bei `--session-assignment` `record` und `play`, bei `--record-empty-sessions` `replay` und `play`) ist die Variable schon jetzt `PGR-E2001`.

**Ziel:** Mehrere Client-Verbindungen werden im Record parallel als eigene Sessions aufgezeichnet (Verbindungen ohne Anfrage nicht); im Replay erhält die n-te Verbindung mit einer Anfrage die n-te Session; Verbindungsende und Verbindungsfehler verhalten sich wie spezifiziert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Benutzerhandbuch und `README.md` — Doku-Folge-Slice `slice-v1-abschluss-sessions-doku` (dort §1 und DoD), direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11, §3.13; Entscheidung des Nutzers vom 2026-10-10): Mit der Dokumentation als eigener Schicht läge dieser Plan über zwei Schichten. Zwischen beiden Slices hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen.
- Eine deterministische Zuordnung bei gleichzeitigem Verbindungsaufbau — Out-of-Scope von LH-FA-12.
- Eine Änderung am PGWire-Adapter (`internal/adapters/driving/pgwire`) — entschieden vom Architect (2026-10-12): Er nimmt Verbindungen schon nebenläufig an, ruft `OpenSession` beziehungsweise `OpenConnection` nach der Startnachricht und lässt `SSLRequest`, `GSSENCRequest` und `CancelRequest` vor dem Aufruf enden; beide Optionen ändern nur Use Cases, Optionen und Verdrahtung. Seine Tests bleiben unberührt, die Tests dieses Slice liegen in den Services, im CLI-Adapter, im Bootstrap und in `test/integration`. Zeigt sich im Code, dass er sich ändern muss, ist das die Rückführung in §4.
- Signalbehandlung — `slice-v1-abschluss-herunterfahren` (aus `slice-v1-abschluss-betrieb` hervorgegangen).
- Das Ende des Servers (Fehlerantwort `FATAL` oder `PANIC`, danach das Verbindungsende) als Ereignis einer Session — `slice-v1-abschluss-verbindungsende-record` (dort §1, *Übernommen aus `slice-v1-abschluss-sessions`*) für das Aufzeichnen und `slice-v1-abschluss-verbindungsende-replay` für die Wiedergabe und die Zuordnung (Entscheidung des Nutzers vom 2026-10-12, Option B, [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md); Sendung des Planners, `AGENTS.md` §3.13): „Verbindungsende und Verbindungsfehler“ im Ziel dieses Slice meint das Ende der Client-Verbindung und die Verbindungsfehler nach `LH-FA-13.b`, nicht ein Ende, das der Server nach einer Fehlerantwort herbeiführt. Dieser Slice liefert die Verwaltung der Sessions (Kennung nach Sessionende, `--record-empty-sessions`, `connection`), die jene Slices unverändert benutzen; sie stehen in der Reihenfolge dahinter.

**Nahtstelle zu den Verbindungsende-Slices** (Architect, 2026-10-12; `slice-v1-abschluss-verbindungsende-replay` und `-record` stehen in der Reihenfolge hinter diesem Slice, Schritt 37 und 40 gegen 26; Sendung nach `AGENTS.md` §3.13, Eintrag in deren §1 jeweils durch den Planner): Diese Slice legt fest, was jene nicht brechen dürfen:

1. *Eine Stelle entscheidet, was für das Replay eine Interaktion ist.* `NewReplayService` nimmt Interaktionen, die Lebendprüfungen sind, in einem Schritt aus den Sessions; ob eine Session Interaktionen hat, folgt daraus, und aus nichts anderem. `-replay` ändert diesen einen Schritt (eine Lebendprüfung mit Ende bleibt, `LH-FA-09.a`); `first-request`, `connection`, `Unassigned` und der Handshake lesen das Ergebnis.
2. *Bei `connection` ist die Session mit der Annahme zugeordnet* (`LH-FA-12.a`): `OpenConnection` setzt sie, nicht die erste Anfrage. Das Schließen der Verbindung nach einer Interaktion mit Ende (`LH-FA-25.a`) ändert die Zuordnung nicht.
3. *Der Record-Service entscheidet beim Ende einer Session, ob und unter welcher `id` sie geschrieben wird* (`CloseSession`, `id` = Zahl der geschriebenen Sessions + 1, unter dem Sperrzustand des Services): `-record` fügt die Art des Endes (Ende des Servers) hinzu und ändert nicht, dass eine Session mit nicht unterstützter Interaktion und das Zwangsende im Aufbau (`EndForced` ohne laufende und ohne abgeschlossene Interaktion, `wirdVerworfen`) verworfen und jede andere mit `--record-empty-sessions` geschrieben wird. Eine Session mit Ende trägt mindestens die Interaktion mit dem Ende und ist deshalb nie eine leere Session.
4. *Eine Verbindung wird erst mit eröffneter Session eine Session* (Randform in §6): Ein Ende des Servers vor dem ersten `ReadyForQuery` (`LH-FA-25.a` *Kein Ende*) bleibt ein Aufbaufehler und erzeugt keine Session, auch mit `--record-empty-sessions`.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-12`](../../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen): Zwei parallele Verbindungen erzeugen zwei Sessions mit je geordneten Interaktionen, eine Verbindung ohne Anfrage keine; im Replay erhält die n-te Verbindung mit einer Anfrage die n-te Session, eine Anfrage darüber hinaus ist ein Mismatch, eine Lebendprüfung nicht (`LH-FA-09.a`); bei `--session-assignment connection` erhält die n-te Verbindung die Session mit der `id` n auch dann, wenn diese nur Lebendprüfungen enthält; `--record-empty-sessions` und `--session-assignment connection` verhalten sich wie spezifiziert; mit `--fail-on-unconsumed` ist eine solche Session bei `connection` weder nicht verbraucht noch nie zugeordnet (`LH-FA-03.b`); `--record-empty-sessions` mit `--format sqlite` schreibt `meta.empty_sessions` mit dem Wert `true` und die leere Session, und `replay` liest die Datei wie dieselbe Aufzeichnung als YAML (am Binary; aus `slice-v1-abschluss-sqlite-format`, V-174); die Leser beider Optionen (auch ihrer Umgebungsvariablen) sind angebunden und `nichtGelesen()` in `fremd_test.go` führt sie nicht mehr (aus `slice-v1-abschluss-umgebung-anderes-kommando`) (Test).
- [ ] [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus): Ein Verbindungsfehler beendet nur die Verbindung, der Prozess merkt sich die Klasse (Test).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/services` | update | geliefert: `RecordService` mit der Option `LeereSessions` (`CloseSession` entscheidet allein über das Schreiben, die Bedingung des Verwerfens steht in `wirdVerworfen`: eine Session ohne abgeschlossene Interaktion wird mit der Option geschrieben, eine nicht unterstützte nie, ebenso nie das Zwangsende im Aufbau (`EndForced` ohne laufende und ohne abgeschlossene Interaktion); `Finish` schreibt die Kennzeichnung `empty_sessions` auch ohne Session; `nichtGeschrieben` kennt den Grund mit der Option) und `ReplayService` mit der Option `ZuordnungNachVerbindung` (`OpenConnection` setzt die Session mit der Kennung n und liefert deren Serverparameter und Version, `alle` führt auch Sessions ohne Interaktion und nur aus Lebendprüfungen, `Unassigned` zählt nach der Zahl der Verbindungen, `NewReplayService` prüft die Kennzeichnung nach den Ladefehlern: `PGR-E2001`, `zuordnen` liefert bei einer Verbindung ohne Session `PGR-E5003`, `keineWeitere` ersetzt `letzteNummer` und greift bei null Interaktionen nicht außerhalb der Liste); der Schritt `ohneLebendpruefungen` in `NewReplayService` bleibt die eine Stelle, die entscheidet, was eine Interaktion ist; die Kommentare an `ReplayService`, `NewReplayService`, `OpenConnection`, `zuordnen`, `cursor` und der Port `driving.Replayer` sind nachgezogen (`AGENTS.md` §3.7) |
| `internal/adapters/driving/cli` | update | geliefert: `--record-empty-sessions` (Wahrheitswert, nur `record`) und `--session-assignment first-request\|connection` (nur `replay`) mit Umgebungsvariable, Schlüssel, `config show`, Hilfe; `ReplayOptions.SessionAssignment` und die Konstanten `ZuordnungErsteAnfrage`, `ZuordnungVerbindung`; `nichtGelesen()` in `fremd_test.go` führt beide Optionen nicht mehr; der Platzhalter in `TestDateiUngueltig` ist `unbekannte_option`, `record_empty_sessions` und `session_assignment` haben eigene Fälle |
| `internal/adapters/driving/pgwire` | unverändert | Verbindungen laufen dort schon nebenläufig, die Client-Richtung je Session bleibt eine Schleife (Bindung aus `slice-v1-abschluss-zeitangaben-aufzeichnen`); Entscheidung des Architect in §1; `git diff --stat` des Slice nennt keine Datei dort |
| `internal/bootstrap` | update | reine Verdrahtung: `services.LeereSessions(o.RecordEmptySessions)` an den `RecordService`, `services.ZuordnungNachVerbindung` an den `ReplayService`, wenn `o.SessionAssignment` `connection` ist; keine Logik, zählt nicht als Schicht (`AGENTS.md` §3.13); die Prüfung der Kennzeichnung liegt im Service, der Fehler läuft über `fail` |
| `test/integration` | update | `sessions_e2e_test.go` am Binary gegen PostgreSQL: Aufzeichnen mit der Option in `yaml` und `sqlite`, parallele Verbindungen und Kennungen nach der Reihenfolge der Enden, `connection` mit und ohne Kennzeichnung, Verbindungen, die nicht zählen, Zwangsende im Aufbau, halbe Startnachricht (`PGR-E6001`), Verbindungsfehler im Replay und im Record, Quellen und ungültige Werte |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-replay-semantik` ist `done`, ebenso die Schritte 1 bis 25 der Reihenfolge in §5 von welle-v1-abschluss (geprüft vom Architect am 2026-10-12: Stand `main` 061365a, alle in `done/`, darunter `slice-v1-abschluss-sqlite-format-lesen`, `-umgebung-anderes-kommando`, `-zeitangaben-aufzeichnen`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Parallelität verlangt eine Änderung des Recording-Formats — zurück zur Zerlegung.
- `in-progress` → `next`: Der PGWire-Adapter ändert sich (dritte Schicht neben Kern und CLI-Adapter, siehe §8) oder die Nebenläufigkeit lässt Aufrufe einer Session überlappen (Bindung aus `slice-v1-abschluss-zeitangaben-aufzeichnen`) — zurück zur Zerlegung, Schnitt vor dem Anspruch.
- `in-progress` → `next`: Der Diff liegt über 2200 Zeilen mit Tests (Schätzung des Architect in §8: 1400 bis 1900) — zurück zur Zerlegung; Schnittlinie: `--session-assignment` mit `connection` (Kern und CLI-Adapter, Randformen der Zuordnung in §6) als eigener Slice hinter dem Teil `--record-empty-sessions` und Parallelität.
- `in-progress` → `open`: Die Zuordnungsregel (n-te Verbindung mit Anfrage, n-te Session) ändert sich durch eine Entscheidung des Auftraggebers — Carveout.

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

**Randformen** (`AGENTS.md` §3.12), vom Architect geprüft und entschieden am 2026-10-12, vor dem Code. Ort der Entscheidung ist die Spezifikation (`LH-FA-12.a` *Randformen des Recordings* und *Randformen der Zuordnung*, `LH-FA-13.b`, `LH-FA-17.a`, `LH-FA-03.b`, `LH-FA-09.a`); eine ADR ist nicht nötig (keine neue Entscheidung über Wirkung oder Alternativen, nur Randformen eines Vertrags aus `LH-FA-12`). Zu jeder Zeile gehört im Bericht des Implementers nach `AGENTS.md` §3.10 *Zusage · Mutation · roter Test*. Findet der Implementer eine Randform, die hier nicht steht, hält er an und gibt sie zurück.

| Randform | Entscheidung | Ort |
|---|---|---|
| *Aufrufe einer Session* (Bindung aus `slice-v1-abschluss-zeitangaben-aufzeichnen`) | nacheinander, nie gleichzeitig; kein Test (akzeptiertes Negativ im Geber) | `LH-FA-21.a` *Messpunkt*, Sicht Zeile *Record-Use-Case* |
| *Variable beim eigenen und beim anderen Kommando* (Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando`) | `--session-assignment` nur `replay`, `--record-empty-sessions` nur `record`; beim anderen Kommando `PGR-E2001`; Test: `TestLeserUmgebungNichtGelesenerOptionen` mit gekürztem `nichtGelesen()`, `TestLeserUmgebungAndererKommandos` | `LH-FA-17.a` |
| *Wertemenge, Alias, leerer und unbekannter Wert* von `--session-assignment` | genau `first-request` und `connection`, klein geschrieben, kein Alias; jeder andere Wert, auch `Connection`, `first_request`, der leere, ist `PGR-E2001` (Datei `PGR-E2004`, Schlüssel `session_assignment` im Abschnitt `replay:`, `record_empty_sessions` im Abschnitt `record:`); die booleschen Formen der Option folgen `LH-FA-17.a` | `LH-FA-12.a`, `LH-FA-17.a` |
| *Angenommene Verbindung*, Record und Replay | zählt ab gelesener Startnachricht und (Record) eröffneter Upstream-Session; `SSLRequest`, `GSSENCRequest`, `CancelRequest`, Ende vor der Startnachricht, gescheiterter Aufbau zum Upstream und Zwangsende im Aufbau erzeugen keine Session, auch mit `--record-empty-sessions`, und zählen im Replay nicht; Test in `test/integration` am Binary (nicht im PGWire-Adapter) | `LH-FA-12.a` |
| *Verbindung ohne Anfrage mit `--record-empty-sessions`* und die anderen Enden einer Session ohne abgeschlossene Interaktion | jede eröffnete Session wird geschrieben (Ende durch den Client, Verbindungsfehler, Zwangsende; `PGR-E4006` nennt die `id`); eine Session mit nicht unterstützter Interaktion wird verworfen, mit und ohne Option; ohne Option bleibt es bei mindestens einer abgeschlossenen Interaktion | `LH-FA-12.a`, `LH-FA-13.a` |
| *Kennzeichnung ohne Session* | `empty_sessions` steht im Recording, sobald der Lauf die Option hat, auch bei null Sessions, in `yaml` und `sqlite` | `LH-FA-12.a`, `SPEC-043` |
| *Zwangsende im Aufbau, wenn der Aufbau danach noch gelingt* (Rückgabe des Implementers, Option (a)) | `CloseSession` mit `EndForced` ohne begonnene und ohne abgeschlossene Interaktion ist das Zwangsende im Aufbau: nie geschrieben, auch mit der Option, keine Meldung `PGR-E4006`; mit `EndForced` und abgeschlossenen oder laufender Interaktion unverändert. Der PGWire-Adapter bleibt unverändert (Option (b) wäre die dritte Schicht). Ruhende Session ohne Interaktion beim Herunterfahren: `EndShutdown` vor der Frist, mit der Option geschrieben. Code: Bedingung `verwerfen` in `CloseSession` (`internal/hexagon/services/record.go`) um `end == model.EndForced && len(interactions) == 0 && !l.laeuft()` ergänzen (vor dem Schreiben der laufenden Interaktion auswerten); Tests: `TestRecordZwangsendeLeereSession` anpassen (Fall „Kennung genannt“ gilt nur mit laufender Interaktion), neuer Fall `EndForced` ohne Interaktion mit Option: nicht geschrieben, `rec.Sessions` leer, keine Meldung; Mutante: Bedingung entfällt, roter Test = dieser Fall. Schichten und Liefer-Punkte unverändert (Kern, CLI; 2 Punkte) | `LH-FA-12.a` *Angenommene Verbindung*, `LH-FA-13.a` |
| *Zwangsende trifft eine ruhende Session ohne Interaktion* (F-679, Architect 2026-10-12) | Auf dem Pfad gibt es sie: `richtungen.zwangsende` (`server.go`) ruft `CloseSession` mit `EndForced` für jede Record-Session, deren Ende nicht schon bemerkt ist; ein zweites Signal lässt die Frist sofort ablaufen und kann dem Ende durch das Herunterfahren (`Shutdown`, `EndShutdown`) zuvorkommen. Der Dienst unterscheidet sie nicht von der im Aufbau; sie wird nicht geschrieben, auch mit der Option, ohne Meldung. Akzeptiertes Negativ: keine Interaktion geht verloren, nur eine leere Session, und nur in einem Lauf, den der Betreiber erzwingt; ein eigenes `SessionEnd` oder ein Aufruf von `Shutdown` im Adapter wäre die dritte Schicht (§4). Kein neuer Code; der vorhandene Fall `EndForced` ohne Interaktion mit Option (nicht geschrieben, keine Meldung) deckt die Zusage, und sein Kommentar nennt beide Fälle | `LH-FA-12.a` *Angenommene Verbindung*, `LH-FA-13.a` |
| *Startnachricht mit vollständigem Kopf, Rest fehlt* (F-680, Architect 2026-10-12) | Bestand: `PGR-E6001`, Exit-Code `6`, keine Session; weniger als 8 Bytes endet still (`PGR-E6001` ist für einen Abbruch mitten in der Nachricht eine unscharfe Klasse, `PGR-E4003` wäre es nicht: die Verbindung ist noch nicht angenommen). Die Spezifikation hält den Bestand fest; ihn auf „still“ zu ändern wäre eine Änderung des PGWire-Adapters (`startup`, `verbindungsende(err)`), also die dritte Schicht, und gehört in einen eigenen Slice, falls der Nutzer es will. Test (Implementer): Integrationstest am Binary in `test/integration`: 8 Bytes Kopf mit gültigem Code und Länge über 8, dann Ende; erwartet Exit-Code `6`, `PGR-E6001`, keine Session auch mit `--record-empty-sessions` | `LH-FA-12.a`, `LH-FA-05.e` |
| *Reihenfolge der Sessionenden bei Abbruch mitten in der Anfrage* (F-681, Architect 2026-10-12) | gemeint ist der Zeitpunkt, zu dem der Recorder das Ende bemerkt (Rückkehr der Anfrage an den Upstream), nicht der Abbruch des Clients; die `id` folgt ihm. Kein Test: der Zeitpunkt der Rückkehr hängt am Upstream, die Aussage ist die der Spezifikation | `LH-FA-12.a` |
| *Gleichzeitiges Ende mehrerer Sessions* (Aufruf-Voraussetzung `CloseSession`) | jedes Ende vergibt genau eine `id`, eindeutig und ohne Lücke; welche Session welche erhält, ist bei gleichzeitigem Ende nicht zugesichert. Der Test der `id`-Reihenfolge erzwingt die Reihenfolge der Enden (die zweite Verbindung endet erst, wenn das Ende der ersten beobachtet ist); ein getrennter Test mit gleichzeitigen Enden prüft nur Eindeutigkeit und Lückenlosigkeit | `LH-FA-12.a` |
| *Aufruf-Voraussetzung der Zuordnung im Replay* | die n-te Verbindung ist bei nacheinander aufgebauten und genutzten Verbindungen deterministisch, bei gleichzeitigem Aufbau nicht zugesichert (Out-of-Scope von `LH-FA-12`); kein Test auf die Reihenfolge bei gleichzeitigem Aufbau | `LH-FA-12.a` |
| *Kennzeichnung fehlt bei `connection`* | Startfehler `PGR-E2001` nach den Fehlern des Ladens (auch `PGR-E3004`), vor dem Öffnen des Ports; die Kennzeichnung genügt, eine leere Session wird nicht verlangt; Test: Aufzeichnung ohne und mit Kennzeichnung, Exit-Code `2` ohne geöffneten Port | `LH-FA-12.a` |
| *`first-request` mit gekennzeichnetem Recording* | Sessions ohne Interaktion werden übersprungen, wie ohne Kennzeichnung | `LH-FA-12.a` |
| *Zuordnung bei `connection` mit der Annahme* | `OpenConnection` setzt die Session; Handshake, Version für Lebendprüfungen und die Meldung am Ende der Verbindung (`LH-FA-03.b`) folgen ihr; eine Session mit Interaktionen, deren Verbindung nichts anfragt, ist zugeordnet und nicht verbraucht, nie zugeordnet nur eine mit `id` über der Zahl der Verbindungen; Test mit `--fail-on-unconsumed` und ohne | `LH-FA-12.a`, `LH-FA-03.b` |
| *Verhalten je Serverversion* (`\v` als Leerraum einer Lebendprüfung ab PostgreSQL 17) | bei `connection` nach `server_version` der zugeordneten Session vom Verbindungsaufbau an; fehlt sie oder ist sie nicht lesbar, ist `\v` kein Leerraum; Test mit einer Session der Version 16 und einer der Version 17 an einer mittleren `id` | `LH-FA-09.a`, `LH-FA-12.a` |
| *Mehr Verbindungen als Sessions* | Handshake der letzten Session; Anfrage, die keine Lebendprüfung ist: `PGR-E5003`; Lebendprüfung wird beantwortet (Version nicht gesetzt) | `LH-FA-12.a` |
| *Anfrage an eine Session ohne Interaktion* (auch nur Lebendprüfungen), je Nachrichtenart | `PGR-E5001` für `Query` und für Extended-Nachrichten (`ClientMessage`); die Diagnose nennt die Session und dass die Aufzeichnung keine (weitere) Interaktion erwartet; kein Index außerhalb der Liste bei null Interaktionen; Test für beide Wege | `LH-FA-12.a`, `LH-FA-10.a` |
| *Herunterfahren und Prozessende bei `connection`* | Prüfung beim Ende der Verbindung und beim Prozessende wie bei `first-request`; Sessions ohne Interaktion melden nichts | `LH-FA-03.b`, `LH-FA-12.a` |
| *Verbindungsfehler paralleler Verbindungen* (`LH-FA-13`) | der Fehler einer Verbindung beendet nur sie, die anderen laufen weiter, der Prozess merkt die Klasse des ersten Fehlers; entstehen Fehler verschiedener Verbindungen gleichzeitig, ist nicht zugesichert, welcher der erste ist (der Test erzwingt die Reihenfolge); Tests: Replay mit einer Abweichung (`PGR-E5001`) an der ersten Verbindung und einer vollständigen zweiten, Record mit einer nicht unterstützten Interaktion (`PGR-E6001`) an einer Verbindung und einer aufgezeichneten zweiten; Exit-Code der Klasse am Ende | `LH-FA-13.b` |

Von der Lieferung nicht berührt (Architect): das Format (`empty_sessions` schreiben und lesen kann `yaml` wie `sqlite` seit `slice-v1-abschluss-sqlite-format` und `-lesen`; kein Schreiben oder Lesen ändert sich), der Recording-Adapter und die Meldungskataloge (`PGR-E2001`, `PGR-E5001`, `PGR-E5003` und `PGR-E4006` gibt es).

**Risiken:**

- Ein Connection-Pool baut Verbindungen gleichzeitig auf; die Zuordnung ist dann nicht zugesichert — **Ausgang: entfallen** (Architect, 2026-10-12): Out-of-Scope von `LH-FA-12`, die Spezifikation bindet die Zusage an nacheinander aufgebaute Verbindungen (`LH-FA-12.a`), kein Test; akzeptiertes Negativ.
- Die Zuordnung `connection` übersieht Sessions nur aus Lebendprüfungen, weil das Replay sie beim Laden aus `frei` nimmt; die n-te Verbindung erhielte dann eine falsche Session. Ein Test mit einer solchen Session an einer mittleren `id` fängt das — **Ausgang:** offen bis Closure.
- Die Ablösung von `frei` als einziger Quelle (Zuordnung nach `id` neben Zuordnung nach Reihenfolge) lässt `Unassigned` eine bei `connection` zugeordnete Session als frei zählen oder umgekehrt; der Test mit `--fail-on-unconsumed` und `connection` fängt das (Randform *Zuordnung bei `connection` mit der Annahme*) — **Ausgang:** offen bis Closure.

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

**Läufe** (Stand `b19f33a` für den Code, dazu die Commits dieses Abschnitts; `make gates` auf sauberem Baum nennt der Bericht): `make test` grün (gofmt, `go vet -tags integration`, alle Unit-Tests, 62 s); `make lint` `0 issues.`; `make a-check` 0 Befunde; `make docs-check` 657 Dateien, 0 Befunde; `make kopf-check` und `make abdeckung-check` Exit 0; `make abdeckung` hat die vier Tabellen neu geschrieben (`b19f33a`); `make test-integration` grün (alle Tests des Repos, 2 min 24 s). Diff `b12f4b4..b19f33a`: 18 Dateien, +2027 −73, ohne die vier Abdeckungstabellen +1958 −68 (Produktcode +170 −56, Tests +1788 −12); unter der Schwelle von 2200 in §4, die Rückführung trat nicht ein. Der PGWire-Adapter ist unberührt (`git diff --stat` nennt keine Datei unter `internal/adapters/driving/pgwire`), das Format und der Recording-Adapter ebenso.

**Mutanten:** Je Mutant eine frische Kopie des Baums (`git ls-files -co`, ohne `.git`, im Scratchpad, nicht im Repo) unter einem neuen Pfad, die mutierte Datei danach mit `touch`; die Unit-Mutanten mit `go test -count=1` im Image `pgwire-recorder:test` ohne Netz, die Kopie als Bind-Mount nach `/src` (kein Build-Kontext, die Mutation kommt an: ein Mutant mit Übersetzungsfehler zeigt es, er wäre sonst grün); die Mutanten am Binary bauen das Image der Stufe `integration` aus der Kopie und führen die Tests gegen eine frische PostgreSQL-Instanz im eigenen Netz. Vor jedem Satz Mutanten lief die unveränderte Kopie grün (`basis-*`).

**Zusage · Mutation (nur sie) · roter Test** (Wortlaut `LH-FA-12.a`, `LH-FA-13.b`, `LH-FA-17.a`; Kern, Record):

| Zusage | Mutation | roter Test |
|---|---|---|
| Mit der Option wird jede eröffnete Session ohne abgeschlossene Interaktion geschrieben | `len(interactions) == 0 && !s.leere` zu `len(interactions) == 0` (R1) | `TestRecordLeereSessionJedesEnde` (fünf Fälle „mit Option“), `TestRecordLeereSessionNachSchreibfehler`, `TestRecordGleichzeitigesEnde`, `TestRecordZwangsendeLeereSession/Kennung_genannt` |
| Ohne die Option wird sie nicht geschrieben (verneinend) | die Bedingung entfällt, jede Session ohne Interaktion wird geschrieben (R2) | `TestRecordLeereSessionJedesEnde` (fünf Fälle „ohne Option“), `TestRecordSessionOhneAnfrage`, `TestRecordParalleleSessions` (dritte Verbindung ohne Anfrage), `TestRecordZwangsende/ohne_abgeschlossene_Interaktion` |
| … nach dem Ende durch den Client ohne Terminate | zusätzlich verwerfen bei `EndClosed` (R3) | `TestRecordLeereSessionJedesEnde/Client_ohne_Terminate_mit_Option` |
| … nach dem Ende durch den Client mit Terminate | zusätzlich verwerfen bei `EndTerminate` | `…/Client_mit_Terminate_mit_Option` |
| … nach einem Verbindungsfehler (Schreibfehler zum Client) | zusätzlich verwerfen bei `EndWriteFailed` | `…/Schreibfehler_zum_Client_mit_Option`, `TestRecordLeereSessionNachSchreibfehler` |
| … nach einem Verbindungsfehler (Fehler des Use Case) | zusätzlich verwerfen bei `EndFailed` | `…/Fehler_des_Use_Case_mit_Option` |
| … beim Herunterfahren | zusätzlich verwerfen bei `EndShutdown` | `…/Herunterfahren_mit_Option` |
| … nach dem Zwangsende einer laufenden Interaktion | zusätzlich verwerfen bei `EndForced` | `TestRecordZwangsendeLeereSession/Kennung_genannt` |
| Eine Session mit nicht unterstützter Interaktion wird verworfen, mit und ohne Option: (a) meldet der Adapter `EndUnsupported` | `end == model.EndUnsupported` entfällt (R4a) | `TestRecordLeereSessionNichtUnterstuetztVerworfen/EndUnsupported,_Option_true`, `TestRecordSessionNichtUnterstuetzt/Client-Nachricht` |
| … (b) hat der Use Case sie selbst als nicht unterstützt markiert | `l.unsupported` entfällt (R4b) | `TestRecordLeereSessionNichtUnterstuetztVerworfen/Serverantwort_nicht_unterstützt…`, `TestRecordExtendedNichtDarstellbar` |
| Die Kennzeichnung `empty_sessions` steht mit der Option | die Option setzt `rec.EmptySessions` nicht (R5a) | `TestRecordKennzeichnung`, `TestRecordLeereSessionJedesEnde` |
| … steht ohne die Option nicht (verneinend) | `rec.EmptySessions` immer wahr (R5b) | `TestRecordKennzeichnung` (zweiter Teil) |
| … steht auch ohne Session | `Finish` setzt sie ohne Session zurück (R5c) | `TestRecordKennzeichnung` (erster Teil) |
| Die Kennung zählt ab 1 | `len(…)+1` zu `len(…)` (R6i) | `TestRecordKennungNachReihenfolgeDerEnden`, `TestRecordSessionMitInteraktionen` |
| … ohne Lücken: eine verworfene Session erhält keine | verworfene Sessions zählen mit (R6ii) | `TestRecordKennungNachReihenfolgeDerEnden` |
| … in der Reihenfolge der Enden (die Enden kommen im Test nacheinander) | die Session wird vorn eingefügt, die Kennung bleibt `len+1` (R6iii) | `TestRecordKennungNachReihenfolgeDerEnden` (Reihenfolge der Kennungen umgekehrt), `TestRecordGleichzeitigeSessions` |
| Jedes Ende vergibt genau eine Kennung, eindeutig und ohne Lücken (gleichzeitiges Ende) | Kennung berechnet, Sperre freigegeben, 5 ms gewartet, Sperre genommen (R7) | `TestRecordGleichzeitigesEnde`, `TestRecordGleichzeitigeSessions` (doppelte Kennungen) |
| Zwei parallele Verbindungen sind zwei Sessions mit je geordneten Interaktionen | `laufende` liefert immer die zuletzt geöffnete Session (R8) | `TestRecordParalleleSessions`, `TestRecordKennungNachReihenfolgeDerEnden` |
| `PGR-E4006` nennt die Kennung, unter der die Session geschrieben wird | die Meldung nennt `0` (R9) | `TestRecordZwangsendeLeereSession/Kennung_genannt` (`Session 2`), `TestRecordZwangsende` |
| … ist die laufende Interaktion nicht unterstützt, nennt sie diesen Grund | `nichtGeschrieben` ignoriert die Option (R10) | `TestRecordZwangsendeLeereSession/nicht_unterstützt` |
| Scheitert der Aufbau zum Upstream (Fehler oder Fehlerantwort), ist die Kennung 0 und keine Session, auch mit der Option | `OpenSession` führt die Session trotz Fehler (A1) | `TestRecordAufbauFehlerKeineSession` (beide Fälle); am Binary `TestE2ERecordKeineVerbindung` |

**Kern, Replay** (`ZuordnungNachVerbindung`):

| Zusage | Mutation | roter Test |
|---|---|---|
| Die n-te Verbindung erhält die Session mit der Kennung n, auch ohne Interaktion und nur aus Lebendprüfungen | Session ohne Interaktion wird nicht zugeordnet (P1a) | `TestReplayConnectionNteVerbindungNteSession`, `TestReplayConnectionAnfrageAnSessionOhneInteraktion` |
| … nach der Kennung, nicht nach der Stelle | Kennung um eins verschoben (P1b) / die n-te Session mit Interaktion (P1c) | `TestReplayConnectionNteVerbindungNteSession`, `…HandshakeUnabhaengigVonAnfragen`, `…Herunterfahren` |
| Die Zuordnung geschieht mit der Annahme: Handshake nach der Session | der Handshake bleibt der der letzten Session (P2a) | `TestReplayConnectionNteVerbindungNteSession` (`Handshake von "s5"`), `…HandshakeUnabhaengigVonAnfragen` |
| … Version für Lebendprüfungen vom Aufbau an | `vt` erst bei der ersten Anfrage (P2b) | `TestReplayConnectionServerversionJeSession` |
| … Meldung am Ende der Verbindung nach der Session | eine Verbindung ohne Anfrage meldet nichts (P2c) | `TestReplayConnectionZuordnungMitDerAnnahme`, `…Herunterfahren` |
| Nie zugeordnet ist nur eine Session mit Interaktionen über der Zahl der Verbindungen: (a) nicht jede freie | Filter entfällt (P3a) | `TestReplayConnectionZuordnungMitDerAnnahme`, `…SessionOhneInteraktionNichtZuMelden` |
| … (b) nicht die der letzten Verbindung | `>` zu `>=` (P3b) | `TestReplayConnectionZuordnungMitDerAnnahme` |
| … (c) nicht eine Session ohne Interaktion (verneinend; auch mit `--fail-on-unconsumed`) | `Unassigned` zählt alle Sessions (P3c) | `TestReplayConnectionSessionOhneInteraktionNichtZuMelden` (beide Läufe) |
| Eine Session ohne Interaktion ist nicht „nicht verbraucht“ (verneinend) | `CloseConnection` meldet auch bei null Interaktionen (P9) | `TestReplayConnectionSessionOhneInteraktionNichtZuMelden` (Panik beim Index) |
| Eine Anfrage an eine Session ohne Interaktion ist `PGR-E5001`, ohne Zugriff außerhalb der Liste | die Prüfung auf die leere Liste entfällt (P4a) | `TestReplayConnectionAnfrageAnSessionOhneInteraktion/nur_Lebendprüfungen` (Panik) |
| … die Diagnose nennt die Session (Query) | `0` statt der Kennung (P4b) | `…AnfrageAnSessionOhneInteraktion` (beide Fälle) |
| … die Diagnose nennt die Session (Extended-Nachricht) | `0` statt der Kennung (P4c) | `…AnfrageAnSessionOhneInteraktion` (beide Fälle) |
| … die Diagnose sagt, dass die Aufzeichnung keine Interaktion erwartet | anderer Text (P4d) | `…AnfrageAnSessionOhneInteraktion` |
| Mehr Verbindungen als Sessions: Handshake der letzten Session der Aufzeichnung | `letzte` bleibt die letzte mit Interaktion (P5a) | `TestReplayConnectionMehrVerbindungenAlsSessions` (Handshake einer anderen Session) |
| … eine Anfrage, die keine Lebendprüfung ist, ist `PGR-E5003` | `PGR-E5001` statt `PGR-E5003` (P5b) | `TestReplayConnectionMehrVerbindungenAlsSessions` |
| … eine Lebendprüfung wird beantwortet, `\v` nicht (Version nicht gesetzt); bei Version 16 ist `\v` keine Lebendprüfung, bei 17 eine | `istLebendpruefung(sql, true)` (P5c) | `TestReplayConnectionMehrVerbindungenAlsSessions`, `TestReplayConnectionServerversionJeSession` |
| Die Kennzeichnung wird bei `connection` verlangt | die Prüfung entfällt (P6a) | `TestReplayConnectionKennzeichnung` |
| … nach den Fehlern des Ladens (`PGR-E3004` geht vor) | Prüfung vor dem Fehler `PGR-E3004` (P6b) | `TestReplayConnectionKennzeichnung` (zwei Fälle) |
| … sie genügt, eine leere Session wird nicht verlangt | die Prüfung verlangt eine Session ohne Interaktion (P6c) | `TestReplayConnectionKennzeichnung`, `…ServerversionJeSession` |
| … der Fehler ist `PGR-E2001` mit Exit-Code 2 | andere Klasse (P6d) | `TestReplayConnectionKennzeichnung` |
| … mit `first-request` ist die Aufzeichnung ohne Kennzeichnung gültig (verneinend) | die Prüfung gilt für beide Zuordnungen (P6e) | `TestReplayConnectionKennzeichnung` (`first-request ohne Kennzeichnung`) |
| Mit `first-request` werden Sessions ohne Interaktion übersprungen, auch in einer gekennzeichneten Aufzeichnung | leere Sessions kommen in die Liste der freien (P7) | `TestReplayFirstRequestGekennzeichneteAufzeichnung`, `TestReplayConnectionKennzeichnung` |
| Zwangsende bei `connection` wie bei `first-request` (`PGR-E4006` mit Kennung, nichts für eine Session ohne begonnene Interaktion) | `Forced` meldet bei `connection` nichts (P8) | `TestReplayConnectionHerunterfahren` |

**CLI-Adapter und Bootstrap:**

| Zusage | Mutation | roter Test |
|---|---|---|
| `--session-assignment` nimmt nur `first-request` und `connection`: (a) kein Alias | `conn` zusätzlich gültig (C1d) | `TestSessionAssignmentWertemenge` |
| … (b) klein geschrieben | `Connection` zusätzlich gültig (C1a) | `TestSessionAssignmentWertemenge`, `TestLeserAlleOptionen`, `TestDateiUngueltig` |
| … (c) `first_request` ist ungültig | `first_request` zusätzlich gültig (C1b) | `TestSessionAssignmentWertemenge`, `TestDateiUngueltig` |
| … (d) der leere Wert ist ungültig | `""` zusätzlich gültig (C1c) | `TestSessionAssignmentWertemenge`, `TestLeserAlleOptionen`, `TestDateiUngueltig` |
| … (e) Leerraum um den Wert ist ungültig | der Wert wird beschnitten (C1e) | `TestSessionAssignmentWertemenge` |
| `--session-assignment` gibt es nur bei `replay` | Option auch bei `record` (C2a) | `TestSessionsOptionenNurBeiEigenemKommando`, `TestLeserZuordnungGleichTabelle` |
| `--record-empty-sessions` gibt es nur bei `record` | Option auch bei `replay` (C2b) | `TestSessionsOptionenNurBeiEigenemKommando`, `TestLeserZuordnungGleichTabelle` |
| Der Standard ist `first-request` / `false` | `connection` (C3a) / `true` (C3b) | `TestSessionAssignmentQuellen`, `TestParseReplay` / `TestRecordEmptySessionsQuellen`, `TestParseRecord` |
| Die Leser sind angebunden: die Option steht in `optionen` | `record-empty-sessions` (C4a) / `session-assignment` (C4b) entfernt | `TestLeserOptionen`, `TestLeserZuordnungGleichTabelle`, `TestRecordEmptySessionsQuellen` / `TestSessionAssignmentQuellen` |
| … ihre Umgebungsvariable wird beachtet | Variable bleibt unbeachtet (C6a, C6b) | `TestLeserAlleOptionen`, `TestRecordEmptySessionsQuellen` / `TestSessionAssignmentQuellen` |
| … der Wert gelangt ins `Command` | fester Wert (C5a, C5b) | `TestLeserAlleOptionen`, `TestRecordEmptySessionsQuellen` / `TestSessionAssignmentQuellen` |
| `nichtGelesen()` führt beide nicht mehr (Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando`) | `nichtGelesen()` führt beide wieder (C7, Mutant im Test) | `TestLeserUmgebungNichtGelesenerOptionen` (die Variable ist beim eigenen Kommando nicht mehr unbeachtet, sondern `PGR-E2001`) |
| Der Bootstrap reicht `--record-empty-sessions` durch | `services.LeereSessions(false)` (B1) | `TestRunRecordEmptySessionsKennzeichnung` (beide Formate); am Binary `TestE2ERecordEmptySessions` u. a. (I1) |
| Der Bootstrap reicht `connection` durch, und nur `connection` | Option entfällt (B2) / immer gesetzt (B3) | `TestRunReplayConnectionKennzeichnung` / `TestRunReplayAusDatei`, `TestRunLogLevel` |
| Die Prüfung der Kennzeichnung liegt vor dem Öffnen des Ports | Port zuerst geöffnet (B4) | `TestRunReplayConnectionKennzeichnung` (Port belegt: `PGR-E4001` statt `PGR-E2001`) |

**Am Binary** (`test/integration`, Mutanten gegen eine frische PostgreSQL-Instanz):

| Zusage | Mutation | roter Test |
|---|---|---|
| `--record-empty-sessions` wirkt am Binary (yaml und sqlite), Kennzeichnung und Session ohne Interaktion | `LeereSessions(false)` (I1) | `TestE2ERecordEmptySessions`, `TestE2ERecordKeineVerbindung`, `TestE2ERecordParalleleVerbindungen`, `TestE2ERecordZwangsendeImAufbau`, `TestE2EReplayConnection`, `TestE2EReplayConnectionOhneKennzeichnung`, `TestE2EReplayZaehlungDerVerbindungen`, `TestE2ESessionsQuellen` |
| Mit `sqlite` steht `meta.empty_sessions` mit dem Wert `true` | der Schreiber schreibt `false` (I5) | `TestE2ERecordEmptySessions` (der Leser lehnt den Wert ab, `sqlite_lesen.go`) |
| Der Prozess merkt die Klasse des ersten Verbindungsfehlers (`PGR-E6001` vor `PGR-E4006`, Exit-Code 6) | der letzte Fehler gilt (I2) | `TestE2EVerbindungsfehlerParallel` |
| Eine Verbindung zählt im Replay erst mit der gelesenen Startnachricht (`SSLRequest`, `GSSENCRequest`, `CancelRequest`, Ende ohne Bytes zählen nicht) | der Server zählt jede angenommene Verbindung (I3) | `TestE2EReplayZaehlungDerVerbindungen` |
| Im Record sind sie keine Session, auch mit der Option | für jede Nichtverbindung eine eröffnete und beendete Session (I4) | `TestE2ERecordKeineVerbindung` |

**Grüne Mutanten:** keiner der Tabelle. Zwei Entwürfe waren unwirksam und sind ersetzt, nicht als grüne Mutanten gezählt: Der erste Mutant zur Reihenfolge der Enden sortierte nach der Öffnung, nachdem er die früheren Kennungen schon neu vergeben hatte, und änderte das Verhalten nicht (grün, äquivalent, ersetzt durch R6iii); der erste Mutant zu I4 öffnete die Session ohne `user` und scheiterte am Upstream (grün, unwirksam, mit `user` und `database` rot).

**Zusagen ohne eigenen Mutanten, mit Grund:**

1. *Aufrufe einer Session nacheinander* und *Aufruf-Voraussetzung der Zuordnung im Replay*: kein Test, entschieden als akzeptiertes Negativ in §6.
2. *Zwangsende im Aufbau ist keine Session* (`LH-FA-12.a`): `TestE2ERecordZwangsendeImAufbau` belegt den Fall, in dem der Aufbau zum Upstream beim Zwangsende noch läuft (der Upstream des Tests antwortet nie): Der Lauf endet nach der Frist ohne `PGR-E4006` mit Exit-Code 0, die Aufzeichnung hat keine Session und trägt die Kennzeichnung. Der Fall, dass der Aufbau nach dem Zwangsende noch gelingt, ist entschieden (Architect, 2026-10-12, Option (a), Randform-Tabelle in §6); sein Test ist ein Unit-Test des Dienstes.
3. *Aufbau scheitert, kein Session*: der Mutant A1 fängt den Dienst; im Adapter liefert der Fehler die Kennung 0, und `CloseSession` mit 0 ist wirkungslos (`TestRecordAufbauFehlerKeineSession`), ein Mutant im Adapter, der eine Kennung ≠ 0 erfindet, wäre derselbe wie A1.

**Nacharbeit nach der Entscheidung des Architects** (Stand `aae9097`; Option (a), Randform-Tabelle in §6): Code `efc17ec` (Hilfsfunktion `wirdVerworfen` in `internal/hexagon/services/record.go`, die Bedingung `EndForced` und keine abgeschlossene und keine laufende Interaktion verwirft die Session, auch mit der Option; sie ist aus `CloseSession` herausgezogen, weil die Zyklomatik sonst 17 statt höchstens 15 ist), Test `TestRecordZwangsendeLeereSession/im_Aufbau` (keine Meldung, keine Schreibung, Aufzeichnung ohne Session mit Kennzeichnung); der Fall „Kennung genannt“ hat eine laufende Interaktion und blieb unverändert. Der PGWire-Adapter ist unberührt. `make test`, `make lint` und `make gates` (Exit 0, auf sauberem Baum `aae9097`) grün. Mutanten wie oben (frische Kopie, `touch`, Bind-Mount in `pgwire-recorder:test`, ohne Netz); die unveränderte Kopie lief vorher grün.

| Zusage | Mutation (nur sie) | roter Test |
|---|---|---|
| `EndForced` ohne laufende und ohne abgeschlossene Interaktion schreibt die Session nie, auch mit der Option | die ganze Bedingung entfällt (`imAufbau := false`, M1) | `TestRecordZwangsendeLeereSession/im_Aufbau` |
| … nur bei `EndForced` | `EndForced` entfällt aus der Bedingung (M2) | `TestRecordLeereSessionJedesEnde` (fünf Fälle „mit Option“), `TestRecordLeereSessionNachSchreibfehler`, `TestRecordGleichzeitigesEnde` |
| … nur ohne abgeschlossene Interaktion | `len(interactions) == 0` entfällt (M3) | `TestRecordZwangsende/ohne_laufende_Interaktion`, `…/abgeschlossen,_nicht_zugestellt` |
| … nur ohne laufende Interaktion | `!l.laeuft()` entfällt (M4) | `TestRecordZwangsendeLeereSession/Kennung_genannt` |
| … ohne Meldung `PGR-E4006` (verneinend) | die Nummer der verworfenen Interaktion wird bei jedem `EndForced` gesetzt (M5) | `TestRecordZwangsendeLeereSession/im_Aufbau`, `TestRecordZwangsende/ohne_laufende_Interaktion`, `…/abgeschlossen,_nicht_zugestellt` |

Grüne Mutanten: keiner. Der Fall, dass der Aufbau nach dem Zwangsende noch gelingt, ist am Dienst belegt (`EndForced` auf einer Session ohne Interaktion); am Binary belegt `TestE2ERecordZwangsendeImAufbau` den Fall, in dem der Aufbau nie endet.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Der Handbuch- und README-Teil liegt im Doku-Folge-Slice `slice-v1-abschluss-sessions-doku` direkt hinter diesem Plan, in derselben Welle; die Dokumentation zählt hier nicht mehr mit. Liefer-Punkte: 2. Schichten: zwei Schichten (Kern, CLI-Adapter).

Nachgezählt beim Eintragen der Sendung aus `slice-v1-abschluss-sqlite-format` (2026-10-11, `AGENTS.md` §3.13): Der Test mit `--format sqlite` liegt in DoD-Punkt 1 und in `test/integration`, ohne Code im Recording-Adapter; Liefer-Punkte weiter 2, Schichten weiter zwei (Kern, CLI-Adapter).

Nachgezählt beim Eintragen der Sendungen aus `slice-v1-abschluss-zeitangaben-aufzeichnen` (Planner, 2026-10-12, `AGENTS.md` §3.13): Beide Sendungen gehören zu DoD-Punkt 1 (Test mit `--record-empty-sessions`, Fälle der Option) und fügen keinen Liefer-Punkt hinzu: Liefer-Punkte 2. Schichten: Die Optionen `--record-empty-sessions` und `--session-assignment` gibt es im CLI-Adapter erst mit diesem Slice (nur die Tests der Optionstabelle nennen sie); der CLI-Adapter ist damit eine Schicht, die §3 bisher nicht führte. Kern und CLI-Adapter sind zwei; der PGWire-Adapter ist die dritte, wenn er sich ändert. Die gleichzeitigen Verbindungen lieferte der Walking Skeleton (§1, *Bereits geliefert*), die Zeile des PGWire-Adapters in §3 trägt damit höchstens Tests; der Architect prüft das vor dem Code. Ändert er den Adapter, liegt der Plan über der Grenze, und der Slice wird vor dem Anspruch geschnitten (§4). Eine Sendung schiebt ihn nicht darüber: Der Fall in `TestDateiUngueltig` liegt im CLI-Adapter, den die Option ohnehin berührt.

Nachgezählt beim Eintragen der Bindung aus `slice-v1-abschluss-umgebung-anderes-kommando` (Planner, 2026-10-12, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 1 (Leser der Option und Kürzung von `nichtGelesen()` im Test des CLI-Adapters) und fügt keinen Liefer-Punkt hinzu: Liefer-Punkte 2. Schichten: unverändert (Kern und CLI-Adapter; der PGWire-Adapter als dritte nur, wenn er sich ändert, §4); der CLI-Adapter war schon gezählt, `fremd_test.go` liegt in ihm.

Nachgezählt beim Eintragen der Abgrenzung zum Ende des Servers (Planner, 2026-10-12, `AGENTS.md` §3.13): Die Abgrenzung nimmt dem Slice nichts und fügt nichts hinzu; Liefer-Punkte weiter 2, Schichten weiter zwei (Kern und CLI-Adapter; der PGWire-Adapter als dritte nur, wenn er sich ändert, §4).

Nachgezählt vom Architect (2026-10-12, Rückgaben F-679 bis F-684 des Reviews): Der Zweig `if o.SessionAssignment == cli.ZuordnungVerbindung` im Bootstrap ist Verdrahtung und zählt nicht als Schicht. Gründe: Er bildet den Wert einer Option eins zu eins auf einen Konstruktor-Parameter des Dienstes ab, wie `if o.FailOnUnconsumed` daneben; er liest keine Datei, stuft keinen Fehler ein und entscheidet nichts nach Daten (die Fälle, die §3.13 als Logik nennt). Kommt ein Zweig hinzu, der etwas anderes tut, als eine Option durchzureichen, zählt der Bootstrap. F-679 und F-680 ändern den PGWire-Adapter nicht (akzeptierte Negative, §6); Schichten bleiben zwei (Kern, CLI-Adapter), Liefer-Punkte 2.

Nachgezählt und geschätzt vom Architect (2026-10-12, `AGENTS.md` §3.13 *Nachzählen*, Größenprüfung vor dem Code): DoD-Punkte 2, die Arbeit darin sind drei Teile: (a) `--record-empty-sessions` und Parallelität im Record (Kern, CLI), (b) `--session-assignment` mit `connection` im Replay (Kern, CLI), (c) die Prüfung der Verbindungsfehler (nur Tests). Liefer-Punkte höchstens 3. Schichten: zwei (Kern, CLI-Adapter); der PGWire-Adapter wird nicht berührt (Entscheidung in §1, Zeile in §3), der Bootstrap ist reine Verdrahtung und zählt nicht, Recording-Adapter und Format sind unverändert. Schätzung des Diffs in Zeilen mit Tests: Produktcode etwa 200 (Kern 120, CLI 60, Bootstrap 20), Tests etwa 1200 bis 1700 (Kern 450, CLI 250 mit `TestDateiUngueltig` und Lesern, `test/integration` 450 bis 700, Bootstrap 50), dazu die erzeugten Abdeckungstabellen; zusammen 1400 bis 1900, vergleichbar `slice-v1-abschluss-zeitangaben-aufzeichnen` (875 Zeilen im Hauptcommit, etwa das 1,5- bis 2-fache hier), unter der Schwelle von 2200 in §4. Das Register führt `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Tests und Nacharbeit liegen gemessen über der Schätzung): Der Spielraum zur Schwelle ist dafür gemeint, die Schnittlinie steht in §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
