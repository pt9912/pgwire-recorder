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

**Bezug:** [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)

**Berührte Spec-Stellen:** `LH-FA-01.a` · `LH-FA-03.b` · `LH-FA-13.b` · `LH-FA-17.a` · `LH-FA-10.a` · `SPEC-012` · `SPEC-018` · `SPEC-034` · `ARC-002` · `ARC-006`

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

Geliefert werden die Option, ihre Umgebungsvariable `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` mit den Werten aus [`LH-FA-17.a`](../../../../spec/spezifikation.md#lh-fa-17a--konfiguration) (nur `true`/`false`, leer heißt nicht gesetzt, CLI vor Umgebungsvariable), die Regeln aus [`LH-FA-03.b`](../../../../spec/spezifikation.md#lh-fa-03b--nicht-verbrauchte-interaktionen) für Verbrauch, Zeitpunkte, Meldung und Rangfolge, mit und ohne Option, und das Handbuch, das diese Regeln für den Anwender nennt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Meldungscodes und Fehlertext-Kopf — `slice-replay-semantik-meldungscodes`.
- Diagnose einer Abweichung — geliefert (siehe „Bereits geliefert“); dieser Slice ändert sie nicht.
- Warnung bei nicht verbrauchten Interaktionen — geliefert von `slice-walking-skeleton-replay` (`PGR-W2001`); dieser Slice ergänzt nur die Nummer der ersten nicht verbrauchten Interaktion in ihrem Text, weil Warnung und Fehler denselben Text tragen.
- Schlüssel `fail_on_unconsumed` der Konfigurationsdatei — `slice-v1-abschluss-betrieb`, der die Konfigurationsdatei für alle Optionen liefert; heute liest das Binary keine.
- Frist `--shutdown-timeout` und weiteres Signal — `slice-v1-abschluss-betrieb`; die Regel für eine dabei zwangsweise beendete Session steht in `LH-FA-03.b` und wird dort mit der Frist geprüft, weil es sie hier noch nicht gibt.
- Strenge Werte boolescher Optionen für `--force` — `slice-v1-abschluss-betrieb`, der `--force` hält; dieser Slice wendet die Regel nur auf die eigene Option an.

**Bereits geliefert** von `slice-walking-skeleton-replay`: Mismatch einfacher Anfragen mit Diagnose (Session, erwartete Nummer, erwartete und empfangene Anfrage), `ErrorResponse` mit `PGR-E5001` und Exit-Code 5 beim Herunterfahren. Von `slice-extended-query-replay`: Mismatch der Extended-Nachrichten mit Diagnose nach `LH-FA-10.a` (Session, Interaktion, Gruppe, Nachricht, erwarteter und empfangener Nachrichtentyp, abweichendes Feld, SQL der erwarteten und der empfangenen Anweisung, ohne Parameterwerte; `TestReplayExtendedDiagnoseAnweisung`), auch für eine einfache Anfrage, wo eine Extended-Nachricht erwartet ist, und umgekehrt (`TestReplayExtendedAbweichung`, `TestReplayExtendedFalscheArt`, `TestE2EReplayExtendedAbweichung`). Dieser Slice ergänzt `--fail-on-unconsumed` (`PGR-E5002`). Titel und Bezug folgen diesem Rest: Die Folgepflicht aus [ADR-0007](../../adr/0007-strict-replay.md), die Diagnose der Abweichung, ist für beide Protokollvarianten geliefert; `--fail-on-unconsumed` schärft [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), nicht das Matching.


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus): Mit `--fail-on-unconsumed` sind nicht verbrauchte Interaktionen und nie zugeordnete Sessions `PGR-E5002` mit Exit-Code 5 statt der Warnung `PGR-W2001` (Test).
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
| `internal/hexagon/model` | update | Meldungscode `PGR-E5002` |
| `internal/hexagon/ports/driving` (Replayer) | update | `CloseConnection` liefert Warnung oder Fehler; `Sent` meldet gesendete Antworten, weil eine Interaktion erst mit ihren gesendeten Antworten verbraucht ist (`LH-FA-03.b` §Verbraucht) |
| `internal/hexagon/services` (Replay-Service) | update | Option `FailOnUnconsumed`: nicht verbrauchte Interaktionen und Sessions als `PGR-E5002` bei gesetzter Option, sonst `PGR-W2001`, derselbe Text mit der Nummer der ersten nicht verbrauchten Interaktion (`LH-FA-03.b`); verbraucht erst nach `Sent` |
| `internal/adapters/driving/pgwire` | update | `Sent` nach jedem erfolgreichen Senden von Antworten; `PGR-E5002` am Ende jeder Verbindung als Verbindungsfehler merken, nach dem Fehler, der sie beendet, ohne `ErrorResponse` |
| `internal/bootstrap` | update | nie zugeordnete Sessions nach dem Ende aller Verbindungen als `PGR-E5002` merken, bevor der Exit-Code gebildet wird; nach einem Startfehler keine Prüfung |
| `internal/adapters/driving/cli` | update | Option `--fail-on-unconsumed` nur bei `replay` (bei `record` unbekannt), Umgebungsvariable, Werte `true`/`false`; eine gesetzte ungültige Umgebungsvariable ist auch neben der Option `PGR-E2001` (`LH-FA-17.a`) |
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

- Exit-Code erst beim Prozessende (Verbindungsfehler beenden nur die Verbindung) — **Ausgang:** offen bis Closure.

**Randformen des Vertrags** (`AGENTS.md` §3.12) — je Randform die Stelle, die sie entscheidet; vom Architect vor dem ersten Code-Commit geprüft:

- Zeitpunkt der Prüfung: am Ende jeder Verbindung, gleich aus welchem Grund (Client schließt, Verbindungsfehler, Herunterfahren nach der laufenden Interaktion, Frist oder weiteres Signal), und für nie zugeordnete Sessions einmal beim Prozessende; nach einem Startfehler nicht — `LH-FA-03.b` §Zeitpunkte.
- Verbindungsfehler beenden nur die Verbindung; `PGR-E5002` wird beim Entstehen gemerkt, der Exit-Code entsteht beim Prozessende — `LH-FA-03.b` §Fehlerebene, `LH-FA-13.b`.
- Rangfolge bei einem anderen Fehler (`PGR-E5001`, `PGR-E6001`, `PGR-E4006` durch die Frist): der Fehler, der die Verbindung beendet, wird zuerst gemerkt, beide Meldungen stehen im Log, der Exit-Code ist der des ersten; nie zugeordnete Sessions zuletzt — `LH-FA-03.b` §Fehlerebene, `LH-FA-13.b` §Prozessende.
- Zustellung: `PGR-E5002` geht nicht als `ErrorResponse` an den Client — `LH-FA-03.b` §Fehlerebene, `LH-FA-13.b`.
- Verbraucht: einfache Interaktion mit gesendeten Antworten, Extended-Interaktion erst mit ihrer letzten Gruppe; eine vor dem letzten `Sync` abgebrochene zählt als nicht verbraucht — `LH-FA-03.b` §Verbraucht.
- Aufgezeichnete Lebendprüfungen zählen nicht; Sessions ohne Interaktion oder nur aus Lebendprüfungen sind weder nicht verbraucht noch nie zugeordnet, auch bei `connection` — `LH-FA-03.b`, `LH-FA-09.a`, `LH-FA-12.a`.
- Verbindung ohne zugeordnete Session (nur Handshake oder nur Lebendprüfungen) meldet nichts — `LH-FA-03.b` §Zeitpunkte.
- Meldung: je Verbindung eine für ihre Session (Session, Zahl nicht verbraucht von allen, `sequence` der ersten nicht verbrauchten), eine Summe für nie zugeordnete Sessions (Zahl, `id` der ersten); derselbe Text mit und ohne Option, Stufe `warn` bzw. `error` — `LH-FA-03.b` §Meldung, `SPEC-034`.
- Andere Kommandos: bei `record` und `play` unbekannte Option (`PGR-E2001`), ihre Umgebungsvariable bleibt dort unbeachtet — `LH-FA-03.b` §Andere Kommandos, `LH-FA-17.a`.
- Konfigurationsweg: Option und Umgebungsvariable, CLI vor Umgebungsvariable; Schlüssel der Konfigurationsdatei mit `slice-v1-abschluss-betrieb` — Optionstabelle in `LH-FA-17.a`, §1.
- Werte: ohne Wert `true`; `=true`/`=false`; jeder andere Wert, auch leer und `1`, `PGR-E2001`; leere Umgebungsvariable gilt als nicht gesetzt; mehrfach auf der Kommandozeile gilt die letzte Angabe — `LH-FA-17.a`.
- Ungültige Umgebungsvariable neben gesetzter Option (etwa `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=1` mit `--fail-on-unconsumed=false`): `PGR-E2001`, Startfehler mit Exit-Code `2`, obwohl die Kommandozeile vorgeht; jeder gesetzte Wert wird geprüft — `LH-FA-17.a`. Vorgabe an den Implementer (Randform-Rückgabe, Stand 6e8d2dd): Verhalten bleibt, wie der Code es heute liest; dazu ein Test, der mit ungültiger Variable und gesetzter Option `PGR-E2001` erwartet und rot wird, wenn die Variable bei gesetzter Option übersprungen wird; daneben bleibt geprüft, dass eine leere Variable nicht gesetzt ist.
- Hilfe neben ungültiger Konfiguration (F-382; etwa `replay --help` mit `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=1`, ebenso `replay --bogus --help` oder `replay --fail-on-unconsumed=1 -h`): Hilfe auf `stdout`, Exit-Code `0`, keine Prüfung von Optionen, Umgebungsvariablen oder Konfigurationsdatei — `LH-FA-01.a` §Hilfe vor Prüfung, `LH-FA-17.a`. Vorgabe an den Implementer: Der Parser erkennt die Hilfe-Angabe vor jeder Prüfung (vor dem Lesen der Variable und vor `fs.Parse`); Tests je Kommando `record` und `replay`: Hilfe mit ungültiger Variable, mit unbekannter Option vor der Hilfe-Angabe und mit ungültigem Optionswert erwartet Hilfe und Exit-Code `0`; rot, wenn die Variable vor der Hilfe geprüft wird (heutiger Stand) oder `fs.Parse` zuerst an der unbekannten Option scheitert; dazu `--input=--help` als Wert und `--` vor `--help` ohne Hilfe.

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
