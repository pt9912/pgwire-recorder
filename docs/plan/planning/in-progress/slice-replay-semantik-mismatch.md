# Slice slice-replay-semantik-mismatch: Nicht verbrauchte Interaktionen als Fehler

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-replay-semantik.

**Bezug:** [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)

**Berührte Spec-Stellen:** `LH-FA-01.a` · `LH-FA-03.b` · `LH-FA-13.b` · `LH-FA-17.a` · `LH-FA-10.a` · `SPEC-012` · `SPEC-018` · `SPEC-034` · `ARC-002` · `ARC-003` · `ARC-005` · `ARC-006` · `ARC-009`

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

**Ziel:** Mit `--fail-on-unconsumed` sind nicht verbrauchte Interaktionen und nie zugeordnete Sessions im Replay ein Fehler (`PGR-E5002`, Exit-Code 5) statt der Warnung `PGR-W2001`. Erkennung und Diagnose der Abweichung sind für beide Protokollvarianten geliefert (siehe unten).

Geliefert werden die Option, ihre Umgebungsvariable `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` mit den Werten aus [`LH-FA-17.a`](../../../../spec/spezifikation.md#lh-fa-17a--konfiguration) (nur `true`/`false`, leer heißt nicht gesetzt, CLI vor Umgebungsvariable), die Regeln aus [`LH-FA-03.b`](../../../../spec/spezifikation.md#lh-fa-03b--nicht-verbrauchte-interaktionen) für Verbrauch, Zeitpunkte, Meldung und Rangfolge, mit und ohne Option, die Regeln aus [`LH-FA-01.a`](../../../../spec/spezifikation.md#lh-fa-01a--kommandos-und-hilfe) zu „Hilfe vor Prüfung“ und „Ende der Optionen“ für `record`, `replay` und den Aufruf ohne oder mit unbekanntem Kommando, und das Handbuch, das diese Regeln für den Anwender nennt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Meldungscodes und Fehlertext-Kopf — `slice-replay-semantik-meldungscodes`.
- Diagnose einer Abweichung — geliefert (siehe „Bereits geliefert“); dieser Slice ändert sie nicht.
- Warnung bei nicht verbrauchten Interaktionen — geliefert von `slice-walking-skeleton-replay` (`PGR-W2001`); dieser Slice ergänzt nur die Nummer der ersten nicht verbrauchten Interaktion in ihrem Text, weil Warnung und Fehler denselben Text tragen.
- Schlüssel `fail_on_unconsumed` der Konfigurationsdatei — `slice-v1-abschluss-betrieb`, der die Konfigurationsdatei für alle Optionen liefert; heute liest das Binary keine.
- Frist `--shutdown-timeout` und weiteres Signal — `slice-v1-abschluss-betrieb`; die Regel für eine dabei zwangsweise beendete Session steht in `LH-FA-03.b` und wird dort mit der Frist geprüft, weil es sie hier noch nicht gibt.
- Hilfe vor Prüfung für `config show` und `--config` — `slice-v1-abschluss-betrieb`, der beide liefert; heute gibt es sie nicht.
- Strenge Werte boolescher Optionen für `--force` — `slice-v1-abschluss-betrieb`, der `--force` hält; dieser Slice wendet die Regel nur auf die eigene Option an.
- `play` kennt `--fail-on-unconsumed` nicht und lässt ihre Umgebungsvariable unbeachtet (`LH-FA-03.b` §Andere Kommandos) — `slice-v1-abschluss-einspielen`, der das Kommando `play` liefert und den Test dafür übernimmt; heute gibt es `play` nicht, ein Aufruf ist ein unbekanntes Kommando.

**Bereits geliefert** von `slice-walking-skeleton-replay`: Mismatch einfacher Anfragen mit Diagnose (Session, erwartete Nummer, erwartete und empfangene Anfrage), `ErrorResponse` mit `PGR-E5001` und Exit-Code 5 beim Herunterfahren. Von `slice-extended-query-replay`: Mismatch der Extended-Nachrichten mit Diagnose nach `LH-FA-10.a` (Session, Interaktion, Gruppe, Nachricht, erwarteter und empfangener Nachrichtentyp, abweichendes Feld, SQL der erwarteten und der empfangenen Anweisung, ohne Parameterwerte; `TestReplayExtendedDiagnoseAnweisung`), auch für eine einfache Anfrage, wo eine Extended-Nachricht erwartet ist, und umgekehrt (`TestReplayExtendedAbweichung`, `TestReplayExtendedFalscheArt`, `TestE2EReplayExtendedAbweichung`). Dieser Slice ergänzt `--fail-on-unconsumed` (`PGR-E5002`). Titel und Bezug folgen diesem Rest: Die Folgepflicht aus [ADR-0007](../../adr/0007-strict-replay.md), die Diagnose der Abweichung, ist für beide Protokollvarianten geliefert; `--fail-on-unconsumed` schärft [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), nicht das Matching.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus): Mit `--fail-on-unconsumed` sind nicht verbrauchte Interaktionen und nie zugeordnete Sessions `PGR-E5002` mit Exit-Code 5 statt der Warnung `PGR-W2001` (Test).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/hexagon/model` | update | Meldungscode `PGR-E5002` |
| `internal/hexagon/ports/driving` (Replayer) | update | `CloseConnection` liefert Warnung oder Fehler; `Sent` meldet gesendete Antworten, weil eine Interaktion erst mit ihren gesendeten Antworten verbraucht ist (`LH-FA-03.b` §Verbraucht) |
| `internal/hexagon/services` (Replay-Service) | update | Option `FailOnUnconsumed`: nicht verbrauchte Interaktionen und Sessions als `PGR-E5002` bei gesetzter Option, sonst `PGR-W2001`, derselbe Text mit der Nummer der ersten nicht verbrauchten Interaktion (`LH-FA-03.b`); verbraucht erst nach `Sent` |
| `internal/adapters/driving/pgwire` | update | `Sent` nach jedem erfolgreichen Senden von Antworten; `PGR-E5002` am Ende jeder Verbindung als Verbindungsfehler merken, nach dem Fehler, der sie beendet, ohne `ErrorResponse` |
| `internal/bootstrap` | update | nie zugeordnete Sessions nach dem Ende aller Verbindungen als `PGR-E5002` merken, bevor der Exit-Code gebildet wird; nach einem Startfehler keine Prüfung |
| `internal/adapters/driving/cli` | update | Option `--fail-on-unconsumed` nur bei `replay` (bei `record` unbekannt), Umgebungsvariable, Werte `true`/`false`; eine gesetzte ungültige Umgebungsvariable ist auch neben der Option `PGR-E2001` (`LH-FA-17.a`); Hilfe-Angabe vor jeder Prüfung, Hilfe je Kommando oder global; Argumente am ersten `--` geteilt, bevor Hilfe-Suche und `fs.Parse` laufen; `--version` ist `PGR-E2001` (`LH-FA-01.a`) |
| `internal/bootstrap` (Test) | update | Hilfe endet mit Exit-Code 0, auch mit ungültiger Umgebungsvariable; `version` liest keine Umgebungsvariable; `--input --` liest keine Datei, `--input=--` scheitert erst beim Laden (`LH-FA-01.a`) |
| `test/integration` | update | Exit-Code 5 mit, Warnung und Exit-Code 0 ohne die Option, auch über die Umgebungsvariable; Herunterfahren mit offener Verbindung; Rangfolge nach `PGR-E6001` (Exit-Code 6) und nie zugeordnete Sessions zuletzt; keine Prüfung nach einem Startfehler. Die Rangfolge nach einem Mismatch (`PGR-E5001`) prüft der Adapter-Test, weil beide Exit-Code 5 tragen |
| `docs/user/benutzerhandbuch.md` | update | Exit-Code 5, `PGR-W2001` auch für nie geöffnete Sitzungen, Rangfolge, Werte der Option |
| `docs/user/abdeckung-*.md` | update | aus den Abdeckungs-Deklarationen der neuen Tests erzeugt (`make abdeckung`) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-extended-query` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: `--fail-on-unconsumed` verlangt eine Änderung am Recording-Format — zurück zur Zerlegung.
- `in-progress` → `open`: Der Replay-Pfad des Walking Skeleton ist unvollständig — Carveout.


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

- Exit-Code erst beim Prozessende (Verbindungsfehler beenden nur die Verbindung) — **Ausgang:** entfallen: so umgesetzt; der Exit-Code entsteht beim Prozessende aus dem ersten gemerkten Fehler, nach dem Verbindungsfehler laufen weitere Verbindungen (`TestE2EReplayNichtVerbrauchtRangfolge`, `TestErsterFehlerZaehlt`; VM04, VM05, VM10 rot; Verifikation, Punkt 1b).

**Randformen des Vertrags** (`AGENTS.md` §3.12) — je Randform die Stelle, die sie entscheidet; vom Architect vor dem ersten Code-Commit geprüft:

- Zeitpunkt der Prüfung: am Ende jeder Verbindung, gleich aus welchem Grund (Client schließt, Verbindungsfehler, Herunterfahren nach der laufenden Interaktion, Frist oder weiteres Signal), und für nie zugeordnete Sessions einmal beim Prozessende; nach einem Startfehler nicht — `LH-FA-03.b` §Zeitpunkte. — **Ausgang:** eingetreten → `slice-v1-abschluss-betrieb` für Frist und weiteres Signal: Dort nennen §1 („Übernommen aus `slice-replay-semantik-mismatch`“) und DoD-Punkt 1 die durch die Frist zwangsweise beendete Session mit `PGR-E4006` vor `PGR-E5002` und Exit-Code 4 (Test mit Signal, je Modus); das weitere Signal lässt dort die Frist sofort ablaufen (§1). Gliedert jener Slice die Frist nach seinem §4 aus, wandert der Punkt mit (`AGENTS.md` §3.9). Die übrigen Zeitpunkte halten `TestE2EReplayNichtVerbraucht`, `…Rangfolge`, `…Herunterfahren` und `…Startfehler` (VM22 rot; Verifikation, Punkt 1a).
- Verbindungsfehler beenden nur die Verbindung; `PGR-E5002` wird beim Entstehen gemerkt, der Exit-Code entsteht beim Prozessende — `LH-FA-03.b` §Fehlerebene, `LH-FA-13.b`. — **Ausgang:** entfallen: so umgesetzt; `closeReplay` merkt `PGR-E5002` über `note`, und in `TestE2EReplayNichtVerbrauchtRangfolge` laufen danach weitere Verbindungen (A3, VM05, VM10 rot; Verifikation, Punkt 1b).
- Rangfolge bei einem anderen Fehler (`PGR-E5001`, `PGR-E6001`, `PGR-E4006` durch die Frist): der Fehler, der die Verbindung beendet, wird zuerst gemerkt, beide Meldungen stehen im Log, der Exit-Code ist der des ersten; nie zugeordnete Sessions zuletzt — `LH-FA-03.b` §Fehlerebene, `LH-FA-13.b` §Prozessende. — **Ausgang:** eingetreten → `slice-v1-abschluss-betrieb` für `PGR-E4006` durch die Frist (DoD-Punkt 1, wie oben). `PGR-E5001` und `PGR-E6001` halten `TestReplayNichtVerbrauchtFehler` und `TestE2EReplayNichtVerbrauchtRangfolge`, nie zugeordnete Sessions zuletzt ebenso (A4, VM04, VM05, B1 bis B3 rot).
- Zustellung: `PGR-E5002` geht nicht als `ErrorResponse` an den Client — `LH-FA-03.b` §Fehlerebene, `LH-FA-13.b`. — **Ausgang:** entfallen: so umgesetzt; `TestReplayNichtVerbrauchtFehler` verlangt, dass der Client nach `Terminate` nichts empfängt (Verifikation, Punkt 1b).
- Verbraucht: einfache Interaktion mit gesendeten Antworten, Extended-Interaktion erst mit ihrer letzten Gruppe; eine vor dem letzten `Sync` abgebrochene zählt als nicht verbraucht — `LH-FA-03.b` §Verbraucht. — **Ausgang:** entfallen: so umgesetzt; `TestReplayVerbraucht`, `TestReplaySent` und `TestReplaySentJeVerbindung` halten es (VM06, VM07, VM08, VM21 rot; F-385 damit behoben).
- Aufgezeichnete Lebendprüfungen zählen nicht; Sessions ohne Interaktion oder nur aus Lebendprüfungen sind weder nicht verbraucht noch nie zugeordnet, auch bei `connection` — `LH-FA-03.b`, `LH-FA-09.a`, `LH-FA-12.a`. — **Ausgang:** eingetreten → `slice-v1-abschluss-sessions` für `connection`: Die Zuordnung `--session-assignment connection` gibt es hier nicht, ein Test kann den Teil deshalb nicht fassen; mit dieser Closure nennen §1 („Übernommen aus `slice-replay-semantik-mismatch`“) und DoD-Punkt 1 jenes Slice ihn, Kopf nachgezogen. Bei `first-request` halten es `TestReplayNichtsZuMelden` und `TestReplayLebendpruefungAufgezeichnet` („1 von 2“).
- Verbindung ohne zugeordnete Session (nur Handshake oder nur Lebendprüfungen) meldet nichts — `LH-FA-03.b` §Zeitpunkte. — **Ausgang:** entfallen: so umgesetzt; `TestReplayNichtsZuMelden` und `TestReplaySessionZuordnung` halten es (M7 im Review rot).
- Meldung: je Verbindung eine für ihre Session (Session, Zahl nicht verbraucht von allen, `sequence` der ersten nicht verbrauchten), eine Summe für nie zugeordnete Sessions (Zahl, `id` der ersten); derselbe Text mit und ohne Option, Stufe `warn` bzw. `error` — `LH-FA-03.b` §Meldung, `SPEC-034`. — **Ausgang:** entfallen: so umgesetzt; die Unit- und E2E-Tests vergleichen den Text exakt, mit und ohne Option (M4, VM01, VM02 rot; Verifikation, Punkt 1a).
- Andere Kommandos: bei `record` und `play` unbekannte Option (`PGR-E2001`), ihre Umgebungsvariable bleibt dort unbeachtet — `LH-FA-03.b` §Andere Kommandos, `LH-FA-17.a`. — **Ausgang:** eingetreten → `slice-v1-abschluss-einspielen` für `play`: Dort nennen §1 („Übernommen aus `slice-replay-semantik-mismatch`“), §3 und seit 29d02a1 DoD-Punkt 1 die unbekannte Option und die unbeachtete Variable, auch mit ungültigem Wert (Integrationstest; V-48 behoben). `record` hält `TestParseFailOnUnconsumedRecord` (C6 rot).
- Konfigurationsweg: Option und Umgebungsvariable, CLI vor Umgebungsvariable; Schlüssel der Konfigurationsdatei mit `slice-v1-abschluss-betrieb` — Optionstabelle in `LH-FA-17.a`, §1. — **Ausgang:** eingetreten → `slice-v1-abschluss-betrieb` für den Schlüssel `fail_on_unconsumed` im Abschnitt `replay:` (§1, DoD-Punkt 3, §3). Option und Umgebungsvariable, CLI vor Umgebungsvariable halten `TestParseFailOnUnconsumedUmgebung` und `TestE2EReplayNichtVerbraucht` (C3, VM11, VM20 rot).
- Werte: ohne Wert `true`; `=true`/`=false`; jeder andere Wert, auch leer und `1`, `PGR-E2001`; leere Umgebungsvariable gilt als nicht gesetzt; mehrfach auf der Kommandozeile gilt die letzte Angabe — `LH-FA-17.a`. — **Ausgang:** entfallen: so umgesetzt; `TestParseFailOnUnconsumed`, `…Werte` und `…Umgebung` halten es (C2, C4, C5 rot; Verifikation, Punkt 1c).
- Ungültige Umgebungsvariable neben gesetzter Option (etwa `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=1` mit `--fail-on-unconsumed=false`): `PGR-E2001`, Startfehler mit Exit-Code `2`, obwohl die Kommandozeile vorgeht; jeder gesetzte Wert wird geprüft — `LH-FA-17.a`. Vorgabe an den Implementer (Randform-Rückgabe, Stand 6e8d2dd): Verhalten bleibt, wie der Code es heute liest; dazu ein Test, der mit ungültiger Variable und gesetzter Option `PGR-E2001` erwartet und rot wird, wenn die Variable bei gesetzter Option übersprungen wird; daneben bleibt geprüft, dass eine leere Variable nicht gesetzt ist. — **Ausgang:** entfallen: entschieden in `LH-FA-17.a` (4e0aab0), `TestParseFailOnUnconsumedUmgebungNebenOption` hält es, die leere Variable `…Umgebung` (C1, C2, VM12 rot). Dass die Entscheidung erst nach dem Code fiel, ist F-381 (§7, Register).
- Hilfe neben ungültiger Konfiguration (F-382; etwa `replay --help` mit `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=1`, ebenso `replay --bogus --help` oder `replay --fail-on-unconsumed=1 -h`): Hilfe auf `stdout`, Exit-Code `0`, keine Prüfung von Optionen, Umgebungsvariablen oder Konfigurationsdatei — `LH-FA-01.a` §Hilfe vor Prüfung, `LH-FA-17.a`. Vorgabe an den Implementer: Der Parser erkennt die Hilfe-Angabe vor jeder Prüfung (vor dem Lesen der Variable und vor `fs.Parse`); Tests je Kommando `record` und `replay`: Hilfe mit ungültiger Variable, mit unbekannter Option vor der Hilfe-Angabe und mit ungültigem Optionswert erwartet Hilfe und Exit-Code `0`; rot, wenn die Variable vor der Hilfe geprüft wird (heutiger Stand) oder `fs.Parse` zuerst an der unbekannten Option scheitert; dazu `--input=--help` als Wert und `--` vor `--help` ohne Hilfe. — **Ausgang:** entfallen: entschieden in `LH-FA-01.a` §Hilfe vor Prüfung (235b9d1) vor dem Code (cf38617); `TestParseHilfeVorPruefung`, `TestRunHilfe` und `TestParseKeineHilfe` nehmen jeden Fall der Vorgabe (VM13, VM14, VM15, VM18 rot; Sonde S1 ergibt Hilfe mit Exit-Code `0`, Verifikation, Abschnitt 4). Eine Konfigurationsdatei liest das Binary hier nicht; die Hilfe vor ihrer Prüfung und für `config show` und `--config` führt `slice-v1-abschluss-betrieb` in §1 und DoD-Punkt 2.
- Formen der Hilfe-Angabe (Rückgabe des Implementers, Stand a205cbc), abgeschlossene Liste — `LH-FA-01.a` §Hilfe vor Prüfung. Vorgabe an den Implementer: Vor jeder anderen Verarbeitung durchsucht der Parser alle Argumente bis zum ersten `--` nach einer Hilfe-Angabe; je Punkt ein Testfall mit erwarteter Hilfe auf `stdout` und Exit-Code `0`, sofern nicht anders genannt:
  - `--h`: Hilfe-Angabe wie `-h`, `-help`, `--help` — Test `replay --h` mit ungültiger `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED`; Gegenfall `replay --Help` ist `PGR-E2001`.
  - Vor dem Kommando und nach unbekanntem Kommando: globale Hilfe — Tests `bogus --help` und `--listen x --help`.
  - `=`-Wert: jede der vier Formen mit beliebigem Wert, auch leer und `false`, ist Hilfe — Tests `replay --help=false`, `replay -h=x`, `replay --help=`.
  - Wertstelle: `--input -help` ist die Hilfe, nicht der Wert — Test `replay --listen x --input -help`; Gegenfall `replay --input=-help --listen x` ist keine Hilfe-Angabe.

  **Ausgang** (Formen der Hilfe-Angabe): entfallen: so umgesetzt; je Punkt nimmt `TestParseHilfeVorPruefung`, `TestParseHilfeGlobal` oder `TestParseKeineHilfe` den genannten Fall und Gegenfall (VM16, VM17, VM19, VM23 rot; Verifikation, Punkt 1d).

- `--` an der Stelle eines Optionswerts (V-45): `--` beendet die Optionen überall, für Hilfe-Suche und Parser gleich; der Wert `--` nur mit `=`; die Meldung des Parsers zur Hilfe-Anforderung entsteht nach außen nie — `LH-FA-01.a` §Ende der Optionen. Vorgabe an den Implementer: Die Argumente werden am ersten `--` geteilt, bevor Hilfe-Suche und `fs.Parse` laufen; `fs.Parse` erhält nur den Teil davor, der Rest sind Argumente. Testfälle: `replay --input -- --help` ist `PGR-E2001` (Option ohne Wert), Exit `2`, und der Fehlertext enthält nicht `help requested`; `replay --listen x --input --` ist `PGR-E2001` (Option ohne Wert), ohne eine Datei zu lesen; `replay --listen x --input=--` nimmt `--` als Wert (Fehler erst beim Laden, `PGR-E3001`); `replay --listen x --input y -- z` ist `PGR-E2001` (unerwartetes Argument). — **Ausgang:** entfallen: entschieden in `LH-FA-01.a` §Ende der Optionen (bd2c0a9) vor dem Code (29d02a1); `TestParseEndeDerOptionen` und `TestRunEndeDerOptionen` nehmen jeden Testfall der Vorgabe. Die roten Mutationen dazu hat nur der Implementer gefahren; die Nacharbeit bd2c0a9 und 29d02a1 hat weder Review noch Verifikation gesehen (§7).

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

- **Was hat funktioniert:** Die Randformen von `--fail-on-unconsumed` standen vor dem ersten Code-Commit in `LH-FA-03.b`, `LH-FA-13.b` und `LH-FA-17.a` (001403f vor 6e8d2dd), und der Kern hielt in beiden Prüfrunden: Das Review sah 23 von 24 Mutationen rot, die Verifikation alle 23 ihrer eigenen, aus dem behaupteten Grund (VM01 bis VM23, Bind-Mount, je Mutant eine frische Kopie). DoD-Liefer-Punkt 1 ist an `cf38617` bestätigt, jede Teilbehauptung einzeln gegen die vier Spezifikationsstellen (Verifikation, Abschnitt 1). Die Nacharbeit danach ändert den Parser (29d02a1), fügt aber nur Tests hinzu; die Tests, die VM11 bis VM20 fingen, sind unverändert, und `make gates` ist an `29d02a1` grün. `AGENTS.md` §3.12 hat jede neue Randform vor ihren Code gebracht: F-382 (235b9d1, cf38617), die vier Formen der Hilfe-Angabe (Rückgabe a205cbc, 56732d6, cf38617) und V-45 (bd2c0a9, 29d02a1). Beide Folge-Slices führen ihre Übernahmen jetzt in §1, DoD und Kopf (F-384, F-388, V-47, V-48).
- **Was ging anders als geplant:** Der Slice ist über seine Option hinaus gewachsen. §1 nannte anfangs nur `--fail-on-unconsumed`; geliefert hat er dazu die Regeln „Hilfe vor Prüfung“ und „Ende der Optionen“ aus `LH-FA-01.a` für `record`, `replay` und den globalen Aufruf. Benannt wurde das Wachstum erst durch V-46 (bd2c0a9). Es kam aus einer Randform, die §6 nicht führte: Eine ungültige Umgebungsvariable der neuen Option stieß an die Hilfe (F-382). Daraus wurden drei Runden an derselben Regel: F-382 im Review, vier Formen der Hilfe-Angabe als Rückgabe des Implementers, `--` an der Wertstelle in der Verifikation (V-45). Jede Runde lief in der Reihenfolge von §3.12, aber die erste Liste in §6 kannte die Nachbarregel nicht, an die die neue Option rührt. Dazu der zweite Fall „im Code entschieden, dann zurückgegeben“ (F-381, 6e8d2dd vor 4e0aab0). Die Abdeckungs-Deklarationen sagten mehr zu, als die Tests prüfen (F-383), und der Plan folgte dreimal nicht (F-384, F-386, F-388; V-46 bis V-48). Zwei Prüfrunden; die Nacharbeit bd2c0a9 (Architect) und 29d02a1 (Implementer) hat weder Review noch Verifikation gesehen, ihre roten Mutationen hat nur der Implementer gefahren.
- **Steering-Loop-Eintrag:** Geschärfte Regel: Ein Mutationsnachweis zählt nur, wenn der Mutant im Build angekommen ist; Mutation und Zurücksetzen laufen über den Build-Kontext nie mit einem Werkzeug, das die mtime erhält, sondern mit `touch`, frischem Pfad oder Bind-Mount, und der Bericht nennt den Weg — liegt in `.claude/commands/implement-slice.md Schritt 19`.
  Auslöser: V-50 (MEDIUM), neu eingetragen als `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×). Herkunfts-Anker `seit slice-replay-semantik-mismatch` am Absatz „Der Mutant muss im Build ankommen“.
  **Warum dieser Eintrag und nicht die beiden anderen Kandidaten:** V-50 trifft das Beweismittel, auf dem `AGENTS.md` §3.10 und Modul 11 (bewusstes Brechen) ruhen. Ein Mutant, der zu Unrecht grün bleibt, meldet eine Testlücke, die es nicht gibt. Ein Zurücksetzen, das nicht ankommt, lässt einen Gate-Lauf den Mutanten prüfen. Beides sieht man dem Bericht nicht an. Die Ursache ist mechanisch und reproduziert (Verifikation, Abschnitt 3), die Abhilfe bekannt und klein. Auf den Zähler zu warten hieße, drei unbemerkt falsche Nachweise abzuwarten. Ob sie unbemerkt blieben, kann der Zähler zudem nicht sehen. Kandidat (b), die unvollständige §6-Liste, ist `BEO-REPO/spec-randform-erst-im-review-entschieden`. Er ist in §3.12 verkörpert, und §3.12 hat hier jedes Mal gegriffen. Eine Schärfung der Art „§6 nennt auch die Nachbarregeln desselben Eingangs“ bliebe Urteil ohne Prüfpunkt; sie gehört in den Retirement-Check der Welle-Closure, die den Eintrag mit 6× liest. Kandidat (c), F-381, steht bei 2×; den Ausgang entscheidet der Zähler, nicht dieser Slice.
- **Beobachtungs-Register (`../observations/`):** je Eintrag `evidence/slice-replay-semantik-mismatch.md`; Zuordnung gegen die Summary-Zeilen des Reviews und die Befund-Tabelle der Verifikation geprüft.
  - `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben/` — F-381 (Summary-Zeile nennt den Eintrag). **2×**, `offen`. Retirement-Check von §3.12: wieder aufgetreten.
  - `BEO-REPO/spec-randform-erst-im-review-entschieden/` (verkörpert in `AGENTS.md` §3.12 seit slice-harness-randformen-vor-code) — F-382 (Summary-Zeile), V-45 (Verifikation nennt den Eintrag), dazu die vier zurückgegebenen Formen der Hilfe-Angabe. Es zählt einmal. **6×**. Retirement-Check von §3.12: wieder aufgetreten, die Regel bleibt.
  - `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung/` (verkörpert in `AGENTS.md` §3.11 seit welle-extended-query) — F-383 (Summary-Zeile); F-387 nennt nur die Zeile *Register* der Negativbefunde, die Summary-Zeile keinen Eintrag. **9×**. Retirement-Check von §3.11: wieder aufgetreten.
  - `BEO-REPO/plan-folgt-korrektur-nicht/` (verkörpert in `AGENTS.md` §3.9 seit welle-walking-skeleton, Sensor `make kopf-check` seit slice-harness-kopf-sensor) — F-384 und F-388 (Summary-Zeilen), F-386 (nur Zeile *Register* der Negativbefunde; die Summary-Zeile nennt keinen Eintrag, der Fund ist §3.9-Kopf), V-46 (nennt den Eintrag), V-47 und V-48 (Klasse von F-388). **8×**. Retirement-Check von §3.9: wieder aufgetreten; `make kopf-check` war jedes Mal grün, die Funde lagen in `Bezug`, an Sicht-Stellen und in Folge-Slices, die er nicht prüft.
  - `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag/` (verkörpert in `AGENTS.md` §3.10 seit welle-extended-query) — F-384 und F-385 (Summary-Zeilen), V-49 (Zusage ohne Test, §3.10). **8×**. Retirement-Check von §3.10: wieder aufgetreten.
  - `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an/` — neu, aus V-50, **1×**, `offen`; die Regel steht schon (Steering-Loop-Eintrag oben), ein weiterer Beleg ist ihr Retirement-Check.

  Einmalig und nicht eingetragen: F-389 (INFO, kein Fehlermuster), V-51 (INFO, Rest von F-383; `slice-v1-abschluss-betrieb` prüft Geheimnisse und Anzeige in DoD-Punkt 2 und 3). Über der Schwelle stehen nur verkörperte Einträge: `spec-randform-erst-im-review-entschieden` 6×, `zusage-im-kommentar-weiter-als-pruefung` 9×, `plan-folgt-korrektur-nicht` 8×, `negativtests-fehlen-bei-neuem-vertrag` 8×. Diese Closure gibt ihnen keinen Ausgang, den Lese-Schritt führt die Closure von welle-replay-semantik.
- **Folge-Slices:** keine neuen. Übernahmen mit Kennung: `slice-v1-abschluss-betrieb` (Frist und weiteres Signal mit `PGR-E4006` vor `PGR-E5002`; Schlüssel `fail_on_unconsumed`; Hilfe vor Prüfung für `config show` und `--config`; strenge Werte für `--force`), `slice-v1-abschluss-einspielen` (`play`), `slice-v1-abschluss-sessions` (`connection` bei Sessions ohne Interaktion, mit dieser Closure in §1, DoD-Punkt 1 und Kopf aufgenommen).
- **Risiken aus §6:** sechzehn, jedes mit genau einem Ausgang. Elf sind entfallen, entschieden in der Spezifikation und von einem Test gehalten: Exit-Code beim Prozessende, Fehlerebene, Zustellung, Verbraucht, Verbindung ohne Session, Meldung, Werte, ungültige Variable neben Option, Hilfe neben ungültiger Konfiguration, Formen der Hilfe-Angabe, `--` an der Wertstelle. Fünf sind eingetreten: Zeitpunkt (Frist), Rangfolge (`PGR-E4006`) und Konfigurationsweg (Schlüssel) an `slice-v1-abschluss-betrieb`; Andere Kommandos (`play`) an `slice-v1-abschluss-einspielen`; Lebendprüfungen bei `connection` an `slice-v1-abschluss-sessions`. Weiter offen: keines. Geprüft, dass jeder Nehmer den Punkt führt: `slice-v1-abschluss-betrieb` in §1 und DoD-Punkt 1 bis 3, `slice-v1-abschluss-einspielen` in §1, §3 und DoD-Punkt 1, `slice-v1-abschluss-sessions` in §1 und DoD-Punkt 1. Alle drei liegen in `open/`.
- **Drei Paarungen:** Anker — `liegt in` nennt `.claude/commands/implement-slice.md` Schritt 19; `grep -rn "seit slice-replay-semantik-mismatch" .claude/ harness/ AGENTS.md docs/plan/planning/observations/` findet ihn dort und in `state.md` von `mutant-kommt-im-build-kontext-nicht-an`. Folge-Slice — keiner neu genannt; die drei Nehmer liegen als Datei in `open/`. Register — die sechs genannten Kennungen bestehen als Verzeichnis mit nicht leerem `evidence/`. Die Closure von welle-replay-semantik prüft erneut.
- **Belege:** Review `docs/reviews/2026-10-05-review-slice-replay-semantik-mismatch.md` (bis `19c9951`; F-381 bis F-389), Verifikation `docs/reviews/2026-10-05-verifikation-slice-replay-semantik-mismatch.md` (bis `cf38617`; V-45 bis V-51; DoD-Liefer-Punkt 1 bestätigt, 23 Mutationen rot, `make gates` grün an `cf38617`); Nacharbeit bd2c0a9 (V-45, V-46, V-47) und 29d02a1 (V-45, V-48, V-49), ohne eigenes Review; `make gates` grün an `29d02a1`. Validierung: n/a in diesem Slice; die Rollen-Sequenz sieht sie vor größeren Wellen vor, nicht je Slice.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; es trägt nur seine `README.md` — keine Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält noch keinen Produktionscode).
