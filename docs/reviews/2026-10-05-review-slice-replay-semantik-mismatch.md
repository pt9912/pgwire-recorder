# Review-Report: slice-replay-semantik-mismatch — 2026-10-05

**Review-Art:** Code und Spezifikation, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff 8538d1a..19c9951`, also die Commits `001403f` (Randformen in der Spezifikation, Plan §1/§3/§6, Folge-Slice), `6e8d2dd` (Code, Tests, Handbuch, Abdeckung, Plan §3), `4e0aab0` (Entscheidung zur ungültigen Umgebungsvariable neben gesetzter Option, Spezifikation, Plan §6, Folge-Slice) und `19c9951` (Test dazu, Abdeckung, Plan §3) für `slice-replay-semantik-mismatch`. Schwerpunkte laut Auftrag: Code gegen die Randformen, Nebenläufigkeit, Option und Umgebungsvariable, §3.10 bis §3.12, Folge-Slice `slice-v1-abschluss-betrieb`.

**Skill:** `.harness/skills/reviewer.md` @ `19c9951`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-replay-semantik-mismatch.md` (Kopf, §1 bis §6, §8), `docs/plan/planning/open/slice-v1-abschluss-betrieb.md` (Kopf, §1, §2, §3), `docs/plan/planning/open/slice-v1-abschluss-einspielen.md` (Kopf, §2, §3)
- `spec/lastenheft.md` [LH-FA-03](../../spec/lastenheft.md#lh-fa-03--replay-modus), [LH-FA-13](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [LH-FA-17](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)
- `spec/spezifikation.md`: LH-FA-01.a, LH-FA-03.b, LH-FA-13.b, LH-FA-17.a, `SPEC-034` (Meldungscodes)
- `spec/architecture.md` `ARC-002`, `ARC-003`, `ARC-005`, `ARC-006`, `ARC-009` und §2 (Schichten)
- [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0007](../plan/adr/0007-strict-replay.md), [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md)
- `AGENTS.md` §3.3 bis §3.12; `.claude/commands/implement-slice.md` *Randform-Rückgabe*; Register `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/plan-folgt-korrektur-nicht`, `BEO-REPO/spec-randform-erst-im-review-entschieden`
- `docs/user/benutzerhandbuch.md` (Wiedergabe, Testautomatisierung, Exit-Codes, Meldungscodes), `docs/user/abdeckung-*.md`
- Vorherige Reports am Modul: [Review zu `slice-extended-query-lebendpruefung`](2026-10-05-review-slice-extended-query-lebendpruefung.md), [Review zu `slice-harness-kopf-sensor`](2026-10-05-review-slice-harness-kopf-sensor.md) (F-371, F-373); höchste vergebene Nummer vor diesem Lauf F-380

**Ausgeführte Läufe im Repo:** `make a-check` vor dem Anlegen dieser Datei (0 Befunde), `make docs-check` und `make abdeckung-check` danach. `make gates` lief nicht. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert (`git status --short` leer).

**Mutationen und Sonden:** Nur in Kopien (`git archive 19c9951`) im Scratchpad. Unit-Tests in einem Container aus der Stufe `deps` des `Dockerfile` ohne Netz. Integrationstests `-test.run NichtVerbraucht` aus einem Image der Stufe `integration` gegen `postgres:17-alpine` (Digest aus `harness/mk/integration.mk`) in einem internen Docker-Netz. Images, Container, Netz, Cache-Volume und Kopien sind gelöscht.

| # | Lauf oder Mutation | Ergebnis |
|---|---|---|
| K0 | Kopie ohne Änderung: Unit-Tests `./internal/...` | grün |
| I0 | Kopie ohne Änderung: die vier Integrationstests `TestE2EReplayNichtVerbraucht*` | grün |
| M1 | `Sent` setzt `ungesendet` nicht zurück | rot (acht Tests, darunter `TestReplayNichtVerbrauchtMeldung`, `TestReplayVerbraucht`) |
| M2 | `Query` setzt `ungesendet` nicht | rot (`TestReplayVerbraucht`) |
| M3 | `ClientMessage` setzt `ungesendet` nach der letzten Gruppe nicht | rot (`TestReplayVerbraucht`) |
| M4 | Meldung nennt den Index statt der aufgezeichneten `sequence` | rot (`TestReplayNichtVerbrauchtMeldung`) |
| M5 | `Unassigned` liefert auch mit Option die Warnung | rot (`TestReplayNieZugeordnetMeldung`) |
| M6 | `Sent` setzt `ungesendet` aller Verbindungen zurück | **grün** (F-385) |
| M7 | `CloseConnection` meldet auch ohne zugeordnete Session | rot (`TestReplaySessionZuordnung`, Panic) |
| M8 | `CloseConnection` liefert auch mit Option die Warnung | rot (zwei Tests) |
| M9 | `ungesendet` nach jeder Extended-Nachricht gelöscht | rot (`TestReplayVerbraucht`) |
| C1 | Umgebungsvariable übersprungen, wenn die Kommandozeile die Option nennt | rot (`TestParseFailOnUnconsumedUmgebungNebenOption`) |
| C2 | leere Umgebungsvariable gilt als gesetzt | rot (vier Tests) |
| C3 | Umgebungsvariable nach der Kommandozeile ausgewertet | rot (`TestParseFailOnUnconsumedUmgebung`) |
| C4 | `1` als wahr angenommen | rot (drei Tests) |
| C5 | Option nicht als boolesche Option (`IsBoolFlag` falsch) | rot (drei Tests) |
| C6 | `record` kennt `--fail-on-unconsumed` | rot (`TestParseFailOnUnconsumedRecord`) |
| A1 | kein `Sent` nach einer Anfrage | rot (`TestReplaySent`) |
| A2 | `Sent` auch nach einer Extended-Nachricht ohne Antwort | rot (`TestReplaySent`) |
| A3 | `PGR-E5002` nur protokolliert, nicht gemerkt | rot (`TestReplayNichtVerbrauchtFehler`) |
| A4 | der letzte statt des ersten Fehlers wird gemerkt | rot (`TestReplayNichtVerbrauchtFehler`, Rangfolge `PGR-E5001` vor `PGR-E5002`) |
| A6 | `Sent` vor statt nach dem Senden | rot (`TestReplaySent`, Fall mit gescheitertem Senden) |
| A7 | kein `Sent` nach einer abschließenden Extended-Nachricht | rot (`TestReplaySent`) |
| B1 | nie zugeordnete Sessions überschreiben den Exit-Code eines Verbindungsfehlers | rot (`TestE2EReplayNichtVerbrauchtRangfolge`, Exit-Code 5 statt 6) |
| B2 | nie zugeordnete Sessions nur geprüft, wenn kein Verbindungsfehler auftrat | rot (`TestE2EReplayNichtVerbrauchtRangfolge`) |
| B3 | nie zugeordnete Sessions bestimmen nie den Exit-Code | rot (`TestE2EReplayNichtVerbrauchtRangfolge`, Exit-Code 0 statt 5) |
| S1 | `replay --help` mit `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=1` | **`PGR-E2001`**, keine Hilfe; `--help` ohne Kommando gibt die Hilfe aus (F-382) |
| S2 | `--fail-on-unconsumed false` (Wert mit Leerzeichen) | `PGR-E2001`, „unerwartetes Argument "false"“ |
| S3 | `play --fail-on-unconsumed` | `PGR-E2001`, „unbekanntes Kommando "play"“ (F-384) |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-381 | HIGH | Die Randform „ungültige Umgebungsvariable neben gesetzter Option“ ist in `6e8d2dd` im Code entschieden. `parseReplay` prüft die Variable vor der Kommandozeile und bricht bei einem ungültigen Wert ab, gleich ob die Option gesetzt ist. Die Rückgabe an den Architect kam danach („Randform-Rückgabe, Stand 6e8d2dd“), die Entscheidung in `4e0aab0`. Sie lautet „Verhalten bleibt, wie der Code es heute liest“. §6 nannte die Randform vor `6e8d2dd` nicht. Der Satz zu Werten in LH-FA-17.a trug sie nicht eindeutig (Bewertung unter Schwerpunkt 6). Dass der Architect einen eigenen Satz in LH-FA-17.a einfügen musste, bestätigt, dass sie offen war. Die Randform-Rückgabe verlangt, anzuhalten, statt im Code zu entscheiden. Inhaltlich folgt daraus nichts: Die Entscheidung passt zur Negative von [LH-FA-17](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), und `19c9951` prüft sie (C1 rot). Der Befund betrifft die Reihenfolge, wie F-373. Für das Register ist das der zweite Fall nach `slice-harness-kopf-sensor`. | `AGENTS.md` §3.12; `.claude/commands/implement-slice.md` *Randform-Rückgabe*; LH-FA-17.a; Reviewer-Skill HIGH „ADR-Verstoß (… Hard Rule)“; `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` | `internal/adapters/driving/cli/cli.go:126`–`:131` (seit `6e8d2dd`); `spec/spezifikation.md:673`–`:676` (seit `4e0aab0`); `docs/plan/planning/in-progress/slice-replay-semantik-mismatch.md:126` | ja (`git show 6e8d2dd -- internal/adapters/driving/cli/cli.go`, `git show 4e0aab0`) | Randform im Code entschieden, danach zurückgegeben |
| F-382 | HIGH | `replay --help` mit gesetzter ungültiger `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` endet mit `PGR-E2001`, ohne Hilfe auszugeben (Sonde S1). Die Variable wird vor `fs.Parse` geprüft, `--help` erst darin. Ohne Kommando gibt `--help` die Hilfe aus, und bei fehlender Pflichtoption gewinnt `--help` ebenfalls. LH-FA-01.a nennt als Ausgabe „Hilfetext bei `--help`“. Ob eine ungültige Konfiguration der Hilfe vorgeht, nennen weder §6 noch LH-FA-17.a. Der Code entscheidet es still und ohne Rückgabe. Die Wirkung ist gering: Wer eine ungültige Variable gesetzt hat, sieht den Fehler statt der Hilfe. Nach Skill-Wortlaut ist es dennoch HIGH, siehe F-371 zur Proportionalität. | `AGENTS.md` §3.12; LH-FA-01.a; LH-FA-17.a; Reviewer-Skill HIGH „ADR-Verstoß (… Hard Rule)“; `BEO-REPO/spec-randform-erst-im-review-entschieden` | `internal/adapters/driving/cli/cli.go:126`–`:136`; `spec/spezifikation.md:32`–`:33`; `docs/plan/planning/in-progress/slice-replay-semantik-mismatch.md:112`–`:126` | ja (Sonde S1) | Randform im Code entschieden, ohne in §6 zu stehen |
| F-383 | MEDIUM | Durch die neuen Abdeckungs-Deklarationen werden [LH-FA-03](../../spec/lastenheft.md#lh-fa-03--replay-modus) und [LH-FA-17](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) „vollständig“ und stehen in `abdeckung-vollstaendig.md`, die `make doc-trace` liest. Die Deklarationen verfehlen dabei die Kriterien des Lastenhefts. Die Boundary von LH-FA-03 lautet „versucht keine Verbindung zu einem Server“. Sie war vorher „—“ und wird jetzt von `TestReplayVerbraucht`, `TestReplayNichtsZuMelden`, `TestReplaySent` und `TestE2EReplayNichtVerbrauchtStartfehler` belegt. Keiner davon prüft das Kriterium. Die Happy von LH-FA-03 belegt `TestParseFailOnUnconsumed`, ein Parser-Test. Die Happy von LH-FA-17 („alle für Record, Replay und Einspielen nötigen Einstellungen … setzbar“) belegen ein Parser-Test und ein E2E-Test über eine einzige Option. Konfigurationsdatei, benannte Verbindungen, `play` und `config show` fehlen. Für die Boundary von LH-FA-17 gilt dasselbe: Geheimnisse und Anzeige fehlen. [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) nennt genau diesen Fall: „Deklarationen können ein Kriterium verfehlen; das fängt nur das Review.“ §3.11 zählt die Abdeckungs-Deklaration ausdrücklich zu den Texten, die nur zusagen dürfen, was geprüft ist. | `AGENTS.md` §3.11; [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) §Konsequenzen; LH-FA-03, LH-FA-17 Akzeptanzkriterien; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Reviewer-Skill MEDIUM „Wiederholung eines Musters“ | `internal/adapters/driving/cli/cli_test.go:72`, `:107`; `internal/hexagon/services/replay_unverbraucht_test.go:86`, `:171`; `internal/adapters/driving/pgwire/server_test.go:701`; `test/integration/unverbraucht_e2e_test.go:46`, `:183`; `docs/user/abdeckung-gesamt.md:13`, `:22`; `docs/user/abdeckung-vollstaendig.md:13`, `:17` | ja (Lesen; `make doc-trace`) | Abdeckungs-Deklaration verfehlt das Kriterium des Lastenhefts |
| F-384 | MEDIUM | Die Randform „bei `play` unbekannte Option, ihre Umgebungsvariable unbeachtet“ steht in §6 und in LH-FA-03.b §Andere Kommandos. Kein Test prüft sie, und kein Slice übernimmt sie. Heute gibt es `play` nicht. `play --fail-on-unconsumed` ist `PGR-E2001`, weil das Kommando unbekannt ist, nicht die Option (Sonde S3). Plan §1 grenzt `play` nicht ab, §3 nennt nur `record`. `slice-v1-abschluss-einspielen` liefert `play` mit „Kommando `play`, Optionen“, nennt aber weder die Option noch die Variable. §3.10 verlangt je Zusage eine Mutation. §3.9 verlangt, den Folge-Slice nachzuziehen, und die Abgrenzung in §1 wäre dessen Adresse. Dasselbe Muster war zuletzt MEDIUM (F-351). Das Handbuch sagt „Die Option gibt es nur bei `replay`“. Das stimmt, sagt für `play` aber nichts Geprüftes zu. | `AGENTS.md` §3.9, §3.10; LH-FA-03.b §Andere Kommandos; LH-FA-17.a; `BEO-REPO/plan-folgt-korrektur-nicht`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` | `docs/plan/planning/in-progress/slice-replay-semantik-mismatch.md:36`–`:43`, `:78`, `:123`; `spec/spezifikation.md:210`–`:212`; `docs/plan/planning/open/slice-v1-abschluss-einspielen.md:16`, `:69`; `internal/adapters/driving/cli/cli_test.go:153`–`:165` | ja (Sonde S3; Lesen) | Zusage ohne Test und ohne übernehmenden Slice |
| F-385 | LOW | Der Port sagt zu, dass `Sent` die Antworten meldet, die `Query` oder `ClientMessage` zuletzt für diese Verbindung geliefert hat. Die Mutation, nach der `Sent` den Zustand aller Verbindungen zurücksetzt, bleibt grün (M6). Jeder Test von `Sent` nutzt eine Verbindung oder ruft `Sent` direkt nach der eigenen Anfrage. Mit der Mutation zählte eine Interaktion als verbraucht, deren Senden scheiterte, wenn zwischen Lieferung und Fehler eine andere Verbindung `Sent` meldet. | `AGENTS.md` §3.10; LH-FA-03.b §Verbraucht; `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` | `internal/hexagon/ports/driving/replay.go:35`–`:38`; `internal/hexagon/services/replay.go:438`–`:447`; `internal/hexagon/services/replay_unverbraucht_test.go` | ja (M6) | Zusage je Verbindung nur mit einer Verbindung geprüft |
| F-386 | LOW | `Berührte Spec-Stellen` nennt `ARC-002` und `ARC-006`. Der Diff erweitert aber auch den Driving Port um `Sent` und ändert die Signatur von `CloseConnection` (`ARC-003`). Er ändert außerdem den CLI-Adapter (`ARC-005`) und die Composition Root (`ARC-009`). Die Zeile für `internal/hexagon/ports/driving` kam erst mit dem Code (`6e8d2dd`) in §3. Die Prüfung vor dem Code (`001403f`) hat die Port-Erweiterung deshalb nicht gesehen. §3.9 erlaubt das Nachziehen im selben Commit, verlangt dann aber auch den Kopf. `make kopf-check` prüft nur Kennungen aus §1 und §2, deshalb meldet es nichts. | `AGENTS.md` §3.9 (Kopf); Maintainability | `docs/plan/planning/in-progress/slice-replay-semantik-mismatch.md:16`, `:73`–`:74` | ja (Lesen) | Kopf nennt die berührten Sicht-Stellen nicht vollständig |
| F-387 | LOW | Das Handbuch sagt zur Umgebungsvariable „die Option geht ihr vor“. Seit `4e0aab0` ist eine ungültige Variable auch neben gesetzter Option `PGR-E2001`. Wer `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=1` gesetzt hat und `--fail-on-unconsumed=false` übergibt, erwartet nach dem Text einen Lauf ohne die Option und erhält einen Startfehler. `19c9951` zieht Plan und Abdeckung nach, das Handbuch nicht. | LH-FA-17.a; `AGENTS.md` §3.9 (sinngemäß, öffentlicher Vertrag); Maintainability | `docs/user/benutzerhandbuch.md:503`–`:505` | ja (Lesen; `TestParseFailOnUnconsumedUmgebungNebenOption`) | Handbuch folgt der späteren Entscheidung nicht |
| F-388 | LOW | `slice-v1-abschluss-betrieb` führt die übernommenen Punkte nur im Absatz „Übernommen aus `slice-replay-semantik-mismatch`“ in §1. Das sind der Schlüssel `fail_on_unconsumed`, strenge Werte für `--force`, die Prüfung jeder gesetzten Variable im allgemeinen Leser und der Test „`PGR-E4006` vor `PGR-E5002`, Exit-Code 4“. DoD und §3 nennen keinen davon ausdrücklich. DoD-Punkt 2 nennt `--force` nur als Ablehnung ohne `--force`. DoD-Punkt 1 nennt `PGR-E4006`, aber nicht die Rangfolge zu `PGR-E5002`. `Bezug` führt LH-FA-03 nicht, obwohl `Berührte Spec-Stellen` jetzt LH-FA-03.b nennt. Der Verifier jenes Slice misst gegen die DoD. Was nur in §1 steht, kann dort ungeprüft durchgehen. | `AGENTS.md` §3.9 (Folge-Slice); `BEO-REPO/plan-folgt-korrektur-nicht` | `docs/plan/planning/open/slice-v1-abschluss-betrieb.md:14`–`:16`, `:44`, DoD-Punkte 1 und 2 | ja (Lesen) | Übernahme nur in §1, nicht in DoD und Kopf |
| F-389 | INFO | `Replayer.Sent` trägt und verletzt weder `ARC-003` noch [ADR-0001](../plan/adr/0001-hexagonale-architektur.md). Die Methode ist fachlich formuliert („Antworten gesendet“, der Begriff aus LH-FA-03.b §Verbraucht) und bringt keinen Import in den Kern. `make a-check` meldet 0 Befunde. Sie setzt §6 „Verbraucht: … mit gesendeten Antworten“ um und entscheidet keine eigene Randform. Ein Senden, das scheitert, verbraucht nicht, und genau das folgt aus dem Spec-Satz. Seit `6e8d2dd` steht die Methode in Plan §3 (zum Kopf siehe F-386). Sie koppelt Lieferung und Meldung zeitlich: `Query`/`ClientMessage`, dann Senden, dann `Sent`. Diese Reihenfolge hält der Adapter, der Port-Kommentar nennt sie, und A1, A2, A6 und A7 prüfen sie. `make a-check` nennt `test/integration/unverbraucht_e2e_test.go` wie die fünf älteren E2E-Dateien als „in keiner Schicht“. Das ist bestehendes Muster, kein Befund dieses Slice. | [ADR-0001](../plan/adr/0001-hexagonale-architektur.md); `ARC-003`; LH-FA-03.b | `internal/hexagon/ports/driving/replay.go:35`–`:38`; `internal/adapters/driving/pgwire/server.go:235`, `:254` | ja (`make a-check`, A1 bis A7) | — (kein Fehlermuster) |

## Antwort auf die Schwerpunkte

1. **Code gegen die Randformen.** Der Code folgt §6 bzw. LH-FA-03.b in allen genannten Punkten.
   - *Zeitpunkt:* `closeReplay` läuft als `defer` am Ende jeder Replay-Verbindung, gleich aus welchem Grund. Ausgenommen ist die Frist, die es noch nicht gibt (Plan §1). Nie zugeordnete Sessions prüft die Composition Root nach `<-done`, also nach dem Ende aller Verbindungen. Ein Startfehler kehrt vorher zurück (B1 bis B3 rot, Startfehler-E2E grün).
   - *Rangfolge:* Der Fehler, der die Verbindung beendet, wird in `fail`/`sendFailed` gemerkt, bevor das `defer` `PGR-E5002` merkt. Der erste gemerkte Fehler bestimmt den Code (A4 rot). Nie zugeordnete Sessions bestimmen den Code nur ohne vorherigen Verbindungsfehler (B1 rot).
   - *Zustellung:* `closeReplay` hat kein Backend und sendet nichts. `TestReplayNichtVerbrauchtFehler` prüft das.
   - *Verbraucht:* Einfache Interaktion nach `Sent`, Extended nach dem `Sent` der letzten Gruppe. Abbruch vor dem letzten `Sync` ist nicht verbraucht (M1 bis M3, M9 rot). Lebendprüfungen fallen beim Laden heraus und zählen nicht mit (M4, `TestReplayNichtsZuMelden`).
   - *`Sent`:* Steht seit `6e8d2dd` in Plan §3. Bewertung in F-389, Lücke F-385, Kopf F-386.
   - *Meldungstext:* Text, Zahl und `sequence` wie in LH-FA-03.b, gleicher Text mit und ohne Option. `zerlege` prüft den Kopf aus `SPEC-034`.
   - *`record`/`play`:* `record` ist geprüft (C6 rot). Für `play` siehe F-384.
2. **Nebenläufigkeit.** `note` merkt den ersten Code unter der Sperre des Servers, verbindungsübergreifend richtig. Innerhalb einer Verbindung ist die Reihenfolge im Log fest, weil die Sitzung sequenziell läuft und `closeReplay` als `defer` zuletzt kommt. Über Verbindungen hinweg protokolliert `note` erst nach dem Lösen der Sperre. Zwei gleichzeitige Fehler können deshalb in anderer Reihenfolge im Log stehen, als sie gemerkt wurden. LH-FA-03.b verlangt eine Reihenfolge nur je Verbindung, und das Verhalten ist älter als dieser Slice. `Unassigned` läuft nach `<-done`, also nach `wg.Wait` aller Verbindungen. Seine Meldung steht deshalb immer zuletzt im Log (B2 rot). `Sent` und `CloseConnection` arbeiten unter der Sperre des Service. Die Zuordnung je Verbindung prüft kein Test (F-385).
3. **Option und Umgebungsvariable.** Strenge Werte für Option und Variable (C2, C4, C5 rot), die letzte Angabe der Kommandozeile gewinnt (Tabellenfälle), CLI vor Variable (C3 rot), leere Variable nicht gesetzt (C2 rot), ungültige Variable neben gesetzter Option `PGR-E2001` (C1 rot). Ein Wert mit Leerzeichen (`--fail-on-unconsumed false`) ist ein unerwartetes Argument (S2). Das folgt aus „ohne Wert `true`“ und ist kein neuer Fall. Offen ist die Hilfe bei ungültiger Variable (F-382), das Handbuch folgt `4e0aab0` nicht (F-387).
4. **§3.10, Stichprobe.** 24 eigene Mutationen, davon 23 rot. Grün blieb M6 (F-385). Die beiden genannten Kandidaten:
   - *„`play` kennt die Option nicht“* hat keinen Test und kann heute keinen haben (F-384).
   - *Rangfolge `PGR-E5001` vor `PGR-E5002`* prüft nur der Adapter-Test. Das trägt: Die Rangfolge entsteht im Adapter (`fail` vor dem `defer`), nicht im Use Case. Der Test fährt den echten Adapter mit einem Fake, der beide Meldungen liefert, und A4 macht ihn rot. Plan §3 begründet die Wahl zutreffend: Beide Codes tragen Exit-Code 5.
5. **§3.11.** Hilfetext, Port- und Service-Kommentare sagen zu, was geprüft ist. Das gilt auch für „nie beide“ (`zerlege` bricht bei beiden ab) und für „Eine Verbindung ohne Session meldet nichts“ (M7 rot). Zu weit reichen die Abdeckungs-Deklarationen (F-383). Im Handbuch sagt „die Option geht ihr vor“ weniger, als seit `4e0aab0` gilt (F-387). Die übrigen Handbuch-Sätze decken Tests: Herunterfahren, Fehlerende, Startfehler, Rangfolge, Werte.
6. **§3.12.**
   - *`Sent`* ist keine im Code entschiedene Randform, sondern die Umsetzung von §6 „Verbraucht“ (F-389).
   - *Die Lesart zur ungültigen Variable* ist eine: Sie stand vor dem Code weder in §6 noch eindeutig in LH-FA-17.a. Der Implementer stützte sich auf „Für die Umgebungsvariable gilt dieselbe Wertemenge … jeder andere Wert … ist `PGR-E2001`“. Dieser Satz legt fest, welche Werte gültig sind, nicht, ob ein durch die Priorität überdeckter Wert geprüft wird. Mit „CLI-Argument > Umgebungsvariable“ daneben lässt er beide Lesarten zu: jeden gesetzten Wert prüfen oder nur den, der wirkt. Die Stütze trägt deshalb nicht. Richtig war, die Randform zurückzugeben. Falsch war, dass `6e8d2dd` sie vorher im Code festlegte (F-381). Die Entscheidung des Architects selbst ist schlüssig: Eine ungültige Konfiguration ist nach der Negative von LH-FA-17 ein Fehler, gleich ob ein anderer Wert vorgeht.
   - *Eine weitere, nicht genannte Randform* entscheidet der Code ohne Rückgabe: `--help` bei ungültiger Variable (F-382).
7. **Folge-Slice `slice-v1-abschluss-betrieb`.** Kopf (`LH-FA-03.b`) und §1 sind nachgezogen, auch mit der Regel aus `4e0aab0`. DoD, §3 und `Bezug` sind es nicht (F-388). Nicht nachgezogen ist zudem `slice-v1-abschluss-einspielen` für `play` (F-384).

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/model` | geprüft, ohne Befund. Nur die Konstante `PGR-E5002`. |
| `internal/hexagon/ports/driving` | geprüft. Fachlich formuliert, keine neuen Importe. Befund F-385 (Zusage je Verbindung), F-389 (INFO). |
| `internal/hexagon/services` | geprüft. Verbrauch, Meldung, Option und Zählung ohne Lebendprüfungen folgen LH-FA-03.b. [ADR-0031](../plan/adr/0031-lebendpruefungen-im-replay.md) ist eingehalten, aufgezeichnete Lebendprüfungen zählen nicht. Befund F-385. |
| `internal/adapters/driving/pgwire` | geprüft. `Sent` nur nach erfolgreichem Senden mit Antwort, `PGR-E5002` gemerkt und nicht zugestellt, Rangfolge über `defer`. Ohne Befund außer F-389 (INFO). |
| `internal/adapters/driving/cli` | geprüft. Befund F-381, F-382. |
| `internal/bootstrap` | geprüft, ohne Befund. Nie zugeordnete Sessions zuletzt, kein Prüfen nach Startfehler (B1 bis B3 rot). |
| Core-Reinheit, [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) | geprüft, ohne Befund. `make a-check` 0 Befunde; `internal/hexagon/**` importiert weder Adapter noch `pgproto3`. |
| [ADR-0007](../plan/adr/0007-strict-replay.md) | geprüft, ohne Befund. Das Matching ist unverändert, die Option ändert nur Code und Stufe der Meldung über nicht Verbrauchtes. |
| `test/integration` | geprüft. Stabil (I0), Rangfolge, Herunterfahren und Startfehler belegt. Befund F-383 (Deklarationen). |
| `docs/user/` | geprüft. Abdeckungstabellen entsprechen den Deklarationen (`make abdeckung-check`). Befund F-383, F-384, F-387. |
| `spec/spezifikation.md` — Strata | geprüft, ohne Befund. Keine ADR-, Slice- oder Commit-Nennung im Diff, Historie nachgetragen. |
| `spec/lastenheft.md`, `spec/architecture.md` | geprüft, ohne Befund. Nicht im Diff. |
| `docs/plan/planning/` — Hard Rule 3.9 | geprüft. §1, §3 und §6 folgen dem Code. Befund F-384, F-386, F-388. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. Kein Move im Diff. |
| Hard Rule 3.5, 3.6 | geprüft, ohne Befund. Keine ADR und keine Gate-Konfiguration im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft, ohne Befund. Neue Kommentare sind Zusagen oder Kopplungen im Indikativ. Der geänderte Kommentar in `replaySitzung` („meldet closeReplay“) beschreibt den Zustand. |
| Hard Rule 3.10 | geprüft. 23 von 24 Mutationen rot. Befund F-384, F-385. |
| Hard Rule 3.11 | geprüft. Befund F-383, F-387. |
| Hard Rule 3.12 | geprüft. Befund F-381, F-382. |
| Commit-Messages `001403f` bis `19c9951` | geprüft, ohne Befund. Jede nennt `slice-replay-semantik-mismatch` und `LH-*`-Kennungen, keine nennt `SPEC-*` oder `ARC-*`. |
| Register `BEO-REPO/*` | geprüft. F-381 ist der zweite Fall von `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`. F-382 betrifft `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, §3.12). F-383 und F-387 betreffen `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. F-384 und F-385 betreffen `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`. F-384, F-386 und F-388 betreffen `BEO-REPO/plan-folgt-korrektur-nicht` (verkörpert, §3.9; Retirement-Check betroffen). |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:**

- Randform im Code entschieden, danach zurückgegeben (F-381; Register `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`)
- Randform im Code entschieden, ohne in §6 zu stehen (F-382; Register `BEO-REPO/spec-randform-erst-im-review-entschieden`)
- Abdeckungs-Deklaration verfehlt das Kriterium des Lastenhefts (F-383; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Zusage ohne Test und ohne übernehmenden Slice (F-384; Register `BEO-REPO/plan-folgt-korrektur-nicht`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`)
- Zusage je Verbindung nur mit einer Verbindung geprüft (F-385; Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`)
- Kopf nennt die berührten Sicht-Stellen nicht vollständig (F-386)
- Handbuch folgt der späteren Entscheidung nicht (F-387)
- Übernahme nur in §1, nicht in DoD und Kopf (F-388; Register `BEO-REPO/plan-folgt-korrektur-nicht`)

## Verdikt

**Merge-blockierend:** ja, wegen F-382. Der Kern trägt: Verbrauch, Rangfolge, Zustellung, Meldung und Werte fangen Mutationen auf Unit- und Integrationsebene (23 von 24 rot). F-381 lässt sich nicht durch Code beheben. Der Befund betrifft die Reihenfolge, das Ergebnis ist entschieden und geprüft. F-382 ist eine offene Randform, die ein Architect entscheiden muss, bevor der Code bleiben kann.

**Übergabe:**

- F-382 an den Architect, mit diesem Report und Sonde S1 als Artefakt. Er entscheidet, ob `--help` einer ungültigen Konfiguration vorgeht, und hält das in LH-FA-01.a oder LH-FA-17.a fest. Danach folgt Code oder Test. Widerspricht der Implementer der Einstufung HIGH, läuft der Konflikt-Pfad über den Architect als Sequenz mit Übergabe-Artefakten (Modul 8).
- F-381 an Architect und Closure: keine Code-Aktion. Die Register-Zuordnung (zweiter Fall) entscheidet die Closure.
- F-383 an den Implementer: Deklarationen der neuen Tests gegen die Kriterien des Lastenhefts prüfen.
- F-384 an Implementer und Planner, weil `slice-v1-abschluss-einspielen` mitbetroffen ist; F-388 an den Planner (`slice-v1-abschluss-betrieb`).
- F-385, F-386 und F-387 an den Implementer. F-389 ohne erwartete Aktion.
- Bei MEDIUM und darunter genügt eine Begründung im Plan oder in der Closure-Notiz.
- Die Finding-Klassen gehen in die Closure-Notiz §7. Die DoD-Konformität prüft der Verifier separat.
