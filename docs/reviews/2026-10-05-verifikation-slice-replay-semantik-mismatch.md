# Verifikation: slice-replay-semantik-mismatch — 2026-10-05

**Rolle:** Verifier (Modul 11). Ich prüfe, ob die Umsetzung Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules prüft der Reviewer.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-replay-semantik-mismatch.md` (Kopf, §1 bis §6, §8) gegen den Gesamt-Diff `8538d1a..cf38617` bei HEAD `cf38617`. Spezifikation: `001403f`, `4e0aab0`, `235b9d1`, `56732d6`. Code: `6e8d2dd`, `19c9951`, `a205cbc`, `cf38617`. Review-Report: `ea75e75`.

**Eingang:**

- der DoD-Liefer-Punkt und die Commit-Messages
- der [Review-Report](2026-10-05-review-slice-replay-semantik-mismatch.md) (F-381 bis F-389, Stand `19c9951`). Die Nacharbeit `a205cbc` und `cf38617` hat kein eigenes Review gesehen. Ich habe sie mit eigenen Mutationen geprüft (Abschnitt 2).
- [`LH-FA-01`](../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-03`](../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration); in der Spezifikation [`LH-FA-01.a`](../../spec/spezifikation.md#lh-fa-01a--kommandos-und-hilfe), [`LH-FA-03.b`](../../spec/spezifikation.md#lh-fa-03b--nicht-verbrauchte-interaktionen), [`LH-FA-13.b`](../../spec/spezifikation.md#lh-fa-13b--fehlerebenen-und-prozessstatus), [`LH-FA-17.a`](../../spec/spezifikation.md#lh-fa-17a--konfiguration)
- Folge-Slices `docs/plan/planning/open/slice-v1-abschluss-betrieb.md` und `docs/plan/planning/open/slice-v1-abschluss-einspielen.md`; [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)

Der Arbeitsbaum war bei meinem Start sauber.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

Die Sensoren unten habe ich in diesem Kontext selbst gefahren. Die Behauptungen des Implementers zählen nicht als Beleg. Mutationen liefen nur in Kopien von `git archive cf38617` im Scratchpad, je Mutation eine frische Kopie mit eigenem Pfad.

**Wie die Mutanten gebaut wurden, und warum so.** Ich habe den Docker-Build-Kontext für die Mutationen umgangen (Abschnitt 3). Ein Image der Stufe `deps` des `Dockerfile` (`verify-mm-deps:cf38617`) trägt nur die Module. Es hat keinen Quellstand.

- **Unit-Tests:** Die Kopie ist schreibgeschützt in einen Container dieses Images eingebunden. Darin laufen netzlos `go test -count=1 ./internal/...`.
- **E2E:** Binary und Testbinary entstehen ebenso aus der eingebundenen Kopie. Das Testbinary läuft gegen `postgres:17-alpine` mit dem Digest aus `harness/mk/integration.mk`, in einem eigenen internen Docker-Netz.

Was der Compiler sieht, ist damit genau der Inhalt der Kopie. Der Go-Build-Cache ist inhaltsadressiert und kann keine alte Fassung liefern.

Die unveränderte Kopie ist auf beiden Wegen grün: alle Unit-Pakete und die vier `TestE2EReplayNichtVerbraucht*`.

Eigene Container, das Netz, das Image, die Cache-Volumes und die Kopien habe ich danach gelöscht. Andere Images des Hosts habe ich nicht angefasst. `git status --short` war vor dem Anlegen dieser Datei leer.

---

## 1. DoD-Liefer-Punkte

### Punkt 1 — [LH-FA-03](../../spec/lastenheft.md#lh-fa-03--replay-modus): Mit `--fail-on-unconsumed` sind nicht verbrauchte Interaktionen und nie zugeordnete Sessions `PGR-E5002` mit Exit-Code 5 statt der Warnung `PGR-W2001` (Test). **Bestätigt.**

Die Teilbehauptungen habe ich einzeln gegen die vier Spezifikationsstellen geprüft. Die Kennungen VM.. verweisen auf Abschnitt 2.

**1a — [LH-FA-03.b](../../spec/spezifikation.md#lh-fa-03b--nicht-verbrauchte-interaktionen)**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Mit Option: nicht verbrauchte Interaktionen sind `PGR-E5002`, Exit-Code 5 | `TestReplayNichtVerbrauchtMeldung`, `TestE2EReplayNichtVerbraucht`, `…Herunterfahren`; VM01, VM02, VM03 rot | bestätigt |
| Mit Option: nie zugeordnete Sessions sind `PGR-E5002`, Exit-Code 5 | `TestReplayNieZugeordnetMeldung`, `TestE2EReplayNichtVerbrauchtRangfolge` (erster Teil); VM09, VM10 rot | bestätigt |
| „statt“: nie Warnung und Fehler zugleich | `zerlege` bricht bei beiden ab; der E2E-Test verlangt, dass der jeweils andere Code fehlt | bestätigt |
| Ohne Option: `PGR-W2001` mit demselben Text, Exit-Code 0 | Tabellenfälle beider Tests; VM01 zeigt, dass der Text in beiden Fällen gleich geprüft wird | bestätigt |
| *Verbraucht* erst mit gesendeten Antworten, Extended erst mit der letzten Gruppe, vor dem letzten `Sync` nicht verbraucht | `TestReplayVerbraucht`, `TestReplaySent`, `TestReplaySentJeVerbindung`; VM06, VM07, VM08, VM21 rot | bestätigt |
| *Zeitpunkte:* Ende jeder Verbindung (Client schließt, Fehlerende, Herunterfahren); nie zugeordnete beim Prozessende; nach Startfehler keine Prüfung | E2E: `…NichtVerbraucht`, `…Rangfolge`, `…Herunterfahren`, `…Startfehler`; VM22 rot. Frist und weiteres Signal sind abgegrenzt (Plan §1), sie gehen an `slice-v1-abschluss-betrieb`. | bestätigt |
| *Meldung:* Session, Zahl der nicht verbrauchten und aller Interaktionen ohne Lebendprüfungen, `sequence` der ersten; Summe der nie zugeordneten mit `id` der ersten | exakte Textvergleiche in den Unit- und E2E-Tests; `TestReplayNichtsZuMelden` | bestätigt |
| Lebendprüfungen zählen nicht; eine Verbindung ohne Session meldet nichts | `TestReplayNichtsZuMelden`, `TestReplayLebendpruefungAufgezeichnet` („1 von 2“) | bestätigt |
| *Fehlerebene:* nicht zugestellt, der erste gemerkte Fehler bestimmt den Exit-Code, nie zugeordnete zuletzt | `TestReplayNichtVerbrauchtFehler`, `TestE2EReplayNichtVerbrauchtRangfolge`; VM04, VM05 rot | bestätigt |
| *Andere Kommandos:* `record` | `TestParseFailOnUnconsumedRecord` | bestätigt |
| *Andere Kommandos:* `play` | abgegrenzt nach `slice-v1-abschluss-einspielen` (Plan §1). Dort steht es in §1 und §3, nicht in der DoD (V-48). | nicht Gegenstand dieses Slice |

**1b — [LH-FA-13.b](../../spec/spezifikation.md#lh-fa-13b--fehlerebenen-und-prozessstatus)**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `PGR-E5002` ist ein Verbindungsfehler und beendet nur die Verbindung, der Prozess läuft weiter | `closeReplay` → `note`; in `…Rangfolge` laufen nach dem gemerkten Fehler weitere Verbindungen | bestätigt |
| nicht zugestellt | `TestReplayNichtVerbrauchtFehler` (der Client empfängt nach `Terminate` nichts) | bestätigt |
| *Prozessende:* nie zugeordnete werden nach allen Verbindungsfehlern gemerkt; Exit-Code der gemerkten Klasse | VM05 (Exit-Code 5 statt 6), VM10 (Exit-Code 0 statt 5) | bestätigt |

**1c — [LH-FA-17.a](../../spec/spezifikation.md#lh-fa-17a--konfiguration)**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| ohne Wert `true`; `=true`/`=false`; jeder andere Wert, auch leer und `1`, ist `PGR-E2001` | `TestParseFailOnUnconsumed`, `…Werte` | bestätigt (C2, C4 und C5 aus dem Review sind rot, und die Tests sind seitdem unverändert) |
| Umgebungsvariable mit derselben Wertemenge, leer heißt nicht gesetzt, CLI geht vor | `…Umgebung`; VM11, VM20 rot, VM20 auch im E2E (Exit-Code 0 statt 5) | bestätigt |
| ungültige Variable ist auch neben gesetzter Option `PGR-E2001` | `…UmgebungNebenOption`; VM12 rot | bestätigt |
| mehrfach angegeben, gilt die letzte Angabe | Tabellenfälle in `TestParseFailOnUnconsumed` | bestätigt |
| die Variable einer Option, die das Kommando nicht kennt, bleibt unbeachtet | `TestParseFailOnUnconsumedRecord` (`record` mit Variable `1` läuft) | bestätigt |
| bei Hilfe wird keine Quelle gelesen oder geprüft | `TestParseHilfeVorPruefung`, `TestRunHilfe`; VM13 rot | bestätigt |

**1d — [LH-FA-01.a](../../spec/spezifikation.md#lh-fa-01a--kommandos-und-hilfe) §Hilfe vor Prüfung**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| vier Formen `-h`, `--h`, `-help`, `--help`; andere Schreibweisen (`---help`, `-H`, `--Help`) sind keine | `TestParseHilfeVorPruefung`, `TestParseKeineHilfe`; VM17, VM18 rot | bestätigt |
| mit `=` und beliebigem Wert, auch leer und `false` | VM16 rot | bestätigt |
| an jeder Stelle vor dem ersten `--`, auch vor dem Kommando und nach einem unbekannten Kommando | `TestParseHilfeGlobal`; VM14 rot | bestätigt |
| Hilfe des Kommandos, wenn das erste Argument ein bekanntes Kommando ist, sonst die globale | VM19 rot | bestätigt |
| es wird nichts weiter geprüft (unbekannte Option, ungültiger Wert, fehlende Pflichtoption, Umgebungsvariable) | VM13, VM14 rot; Sonde am Binary: `replay --help`, `replay --bogus --help`, `replay --fail-on-unconsumed=1 -h`, `record --force=yes -h` mit `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=1` geben je Hilfe auf `stdout` aus, Exit-Code 0, `stderr` leer | bestätigt |
| an der Stelle eines Optionswerts ist sie die Hilfe; als Wert geht sie nur mit `=` | `--input -help` (VM23 rot), `--input=--help` als Wert (VM18 rot) | bestätigt |
| nach `--` ist sie ein gewöhnliches Argument | `TestParseKeineHilfe`; VM15 rot | bestätigt, bis auf die Randform `--` an der Wertstelle (V-45) |
| `version` liest keine Konfiguration; eine Option `--version` gibt es nicht (`PGR-E2001`) | kein Test. Die Sonde am Binary zeigt das zugesagte Verhalten (`--version` und `version --version` ergeben Exit-Code 2). | Zusage ohne Test (V-49) |

### Punkt 2 — `make gates` grün. **Bestätigt.**

`make gates` habe ich selbst im Repo bei `cf38617` gefahren, **Exit-Code 0**, Laufzeit 57 s.

- `a-check`: 0 Befunde
- `a-check-negativ`, `abdeckung-gegenprobe`, `commit-msg-gegenprobe`, `kopf-check-gegenprobe`: grün
- `baseline-verify`: v6.13.0 OK
- d-check: 212 Dateien, 0 Befunde
- `run-integration-tests`: grün, beide Phasen, darunter alle vier `TestE2EReplayNichtVerbraucht*`
- `make kopf-check` einzeln: Exit-Code 0

Weil die Gate-Kette über den Docker-Build-Kontext baut (Abschnitt 3), habe ich Unit- und E2E-Tests zusätzlich über den Bind-Mount-Weg gefahren. Beide sind grün.

### Prozess-Punkte der DoD (nur vermerkt, kein Häkchen gesetzt)

Der Review-Report liegt vor. Closure-Notiz, Register, Risiko-Ausgänge und Paarungen stehen aus. §7 ist leer, und das ist vor der Closure richtig.

---

## 2. Bewusstes Brechen

Die Rot-Spalte gibt die Meldung des Tests an, damit der Grund prüfbar ist. Alle 23 Mutationen sind rot.

| # | Mutation | Rot in | Grund (Testmeldung) |
|---|---|---|---|
| VM01 | **Option im Service wirkungslos** (`meldung` liefert immer die Warnung) | Unit: `TestReplayNichtVerbrauchtMeldung`, `…Verbraucht`, `…NieZugeordnetMeldung`, `…SentJeVerbindung`. E2E: `…NichtVerbraucht`, `…Rangfolge`, `…Herunterfahren` | `mit Option: PGR-W2001 … Exit 0, erwartet PGR-E5002 … Exit 5`; E2E `Exit-Code 0, erwartet 5` |
| VM02 | `PGR-E5002` als `PGR-E2002` (falsche Klasse) | dieselben vier Unit-Tests; E2E wie VM01 | `Fehlertext "Konfiguration [PGR-E2002]: …", erwartet Kopf "Replay [ … ]"`; E2E `Exit-Code 2, erwartet 5` |
| VM03 | Composition Root reicht die Option nicht an den Service | Unit grün (erwartet: kein Unit-Test der Verdrahtung). E2E: `…NichtVerbraucht`, `…Rangfolge`, `…Herunterfahren` | `Exit-Code 0, erwartet 5` |
| VM04 | **Erster Fehler gewinnt nicht:** `note` merkt den letzten | `TestErsterFehlerZaehlt`, `TestReplayNichtVerbrauchtFehler`. E2E: `…Rangfolge` | `erster Fehler "PGR-E5002", erwartet PGR-E5001`; E2E `Exit-Code 5, erwartet 6` |
| VM05 | nie zugeordnete Sessions überschreiben den gemerkten Code | E2E: `…Rangfolge` | `Exit-Code 5, erwartet 6` |
| VM06 | **Verbraucht schon bei Lieferung** (einfache Anfrage setzt `ungesendet` nicht) | `TestReplayVerbraucht`, `TestReplaySentJeVerbindung` | `Antworten nicht gesendet:  ""` (keine Meldung) |
| VM07 | wie VM06, Extended, letzte Gruppe | `TestReplayVerbraucht` | `Antworten der letzten Gruppe nicht gesendet:  ""` |
| VM08 | Adapter meldet `Sent` auch, wenn das Senden scheitert | `TestReplaySent` | `Sent nach gescheitertem Senden gemeldet (1)` |
| VM09 | **Nie zugeordnete Sessions nie gemeldet** (`Unassigned` liefert nichts) | `TestReplaySessionZuordnung`, `TestReplayNieZugeordnetMeldung`. E2E: `…Rangfolge` | `nie zugeordnete Session ohne Warnung: <nil>`; E2E `Exit-Code 0, erwartet 5` |
| VM10 | Composition Root wertet den Fehler aus `Unassigned` nicht aus | E2E: `…Rangfolge` | `Exit-Code 0, erwartet 5` |
| VM11 | **Ungültige Variable ignoriert** | `…Umgebung`, `…UmgebungNebenOption` | `Umgebung "1": erwartet PGR-E2001, erhalten <nil>` |
| VM12 | Variable übersprungen, wenn die Kommandozeile die Option nennt | `…UmgebungNebenOption` | `Umgebung "1", [--fail-on-unconsumed=false]: erwartet PGR-E2001, erhalten <nil>` |
| VM13 | **Hilfe nach der Prüfung:** Variable von `replay` vor der Hilfe-Suche geprüft (Stand vor `cf38617`) | `TestParseHilfeVorPruefung`, `TestRunHilfe` | `["replay" "--help"]: erwartet Hilfe, erhalten … PGR-E2001`; `Exit-Code 2, stderr "Konfiguration [PGR-E2001]: …"` |
| VM14 | Hilfe-Suche nur im ersten Argument nach dem Kommando, also scheitert `fs.Parse` zuerst an der unbekannten Option | `TestParseHilfeVorPruefung`, `TestParseHilfeGlobal` | `["replay" "--bogus" "--help"]: erwartet Hilfe …`; `["--listen" "x" "--help"]: …` |
| VM15 | `--` beendet die Hilfe-Suche nicht | `TestParseKeineHilfe` | `[… "--" "--help"]: Hilfe statt Prüfung` |
| VM16 | `=`-Wert ist keine Hilfe-Angabe | `TestParseHilfeVorPruefung` | `["replay" "--help=false"]: erwartet Hilfe` |
| VM17 | `--h` keine Hilfe-Angabe | `TestParseHilfeVorPruefung`, `TestParseHilfeGlobal` | `["replay" "--h"]: erwartet Hilfe …`; `["--h" "replay"]: …` |
| VM18 | Hilfe auch als Teil eines Arguments (Endung `help`) | `TestParseKeineHilfe` | `["replay" "---help"]: Hilfe statt Prüfung` |
| VM19 | immer die globale Hilfe | `TestParseHilfeVorPruefung` | `replay ["--help"]: nicht die Hilfe von replay` |
| VM20 | Umgebungsvariable gar nicht gelesen | `…Umgebung`, `…UmgebungNebenOption`. E2E: `…NichtVerbraucht` | E2E `Exit-Code 0, erwartet 5` (Fall `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=true`) |
| VM21 | `Sent` setzt den Zustand aller Verbindungen zurück (Mutation M6 des Reviews) | `TestReplaySentJeVerbindung` | `Sent der zweiten Verbindung hat die erste verbraucht` |
| VM22 | nach dem Startfehler (Port belegt) werden nie zugeordnete Sessions doch geprüft | E2E: `…Startfehler` | `nach dem Startfehler geprüft` |
| VM23 | Wert nach `--input`/`--listen` wird bei der Hilfe-Suche übersprungen | `TestParseHilfeVorPruefung` | `["replay" "--listen" "x" "--input" "-help"]: erwartet Hilfe` |

Die sechs verlangten Fälle:

- **Option macht `PGR-E5002` mit Exit-Code 5:** VM01, VM02, VM03
- **erster Fehler gewinnt:** VM04, VM05
- **verbraucht erst nach Senden:** VM06, VM07, VM08, VM21
- **nie zugeordnete Sessions:** VM09, VM10, VM22
- **ungültige Variable:** VM11, VM12, VM20
- **Hilfe vor Prüfung:** VM13 bis VM19, VM23

Jeder Fall wird aus dem Grund rot, den die DoD behauptet.

---

## 3. Docker-Build-Kontext nach dem Zurücksetzen einer Mutation

**Reproduzierbar, aber unter einer engeren Bedingung als berichtet.** Die Sonde lief in einem eigenen Verzeichnis mit `FROM scratch`, `COPY f /f` und `docker build -o type=local`. Ich habe verglichen, was im Kontext ankommt, mit dem, was auf der Platte liegt. Umgebung: Docker 29.8.2, buildx v0.37.1, ext4.

| Fall | Datei vorher → nachher | Im Kontext |
|---|---|---|
| gleiche Größe, gleiche Sekunde, andere Nanosekunden | `AAAA` → `BBBB` | `BBBB` (richtig) |
| **gleiche Größe, mtime auf die Nanosekunde gleich** | `BBBB` → `CCCC` | **`BBBB` (veraltet)** |
| gleiche Größe, andere Sekunde | → `DDDD` | richtig |
| andere Größe, gleiche mtime | → `EEEEE` | richtig |
| Original aus `tar` (mtime ohne Sekundenbruchteil), Mutation mit `sed -i`, Zurücksetzen aus `tar` | `ORIG` → `MUTA` → `ORIG` | jedes Mal richtig |
| Mutation und Zurücksetzen mit `sed -i` in derselben Sekunde | `MUTB` → `ORIG` | richtig |
| **Mutation mit erhaltener mtime** (`sed -i` und `touch -r`) | `ORIG` → `MUTC` | **`ORIG`: Die Mutation kommt nicht an** |
| wie davor, mit `--no-cache` | → `MUTE` | **`ORIG`**: `--no-cache` hilft nicht |
| wie davor, mit `touch` nach dem Zurücksetzen | → `ORIG` | richtig |
| dieselbe Datei mit erhaltener mtime in einem **frischen Pfad** | `MUTE` | richtig |

**Was daraus folgt:**

- BuildKit überträgt eine Datei des Kontexts nicht neu, wenn Größe und mtime dem zuletzt übertragenen Stand **desselben Kontext-Pfads** gleichen.
- Auf ext4 ist „gleiche Sekunde“ nicht genug, die mtime muss auf die Nanosekunde gleich sein. Gleiche Sekunde reicht nur, wenn beide mtimes keinen Sekundenbruchteil tragen. Das ist der Fall nach `tar -x`, `git archive | tar -x` und `touch -d` ohne Bruchteil, und bei `cp -p`, `rsync -a` und `touch -r` von solchen Dateien.
- Die Gefahr geht in **beide Richtungen.** Eine Mutation, die die mtime erhält, kommt nicht im Image an. Der Mutant „überlebt“ dann zu Unrecht, und der Lauf meldet eine Testlücke, die es nicht gibt. Ebenso kommt ein Zurücksetzen, das die mtime erhält, nicht an.
- `--no-cache` hilft nicht. Was hilft: ein frischer Pfad je Lauf, `touch` nach jeder Änderung, oder Bind-Mount statt Build-Kontext (so lief diese Verifikation).
- Die Gate-Läufe im Repo sind nur betroffen, wenn ein Werkzeug dort die mtime erhält. `git checkout` und Editoren setzen die aktuelle Zeit.

Befund V-50.

---

## 4. Review-Findings F-381 bis F-389

| Finding | Status | Beleg |
|---|---|---|
| F-381 (HIGH) | **ohne Code-Aktion, offen für die Closure.** Die Reihenfolge lässt sich nicht nachträglich heilen. Das Ergebnis ist entschieden (`4e0aab0`) und geprüft (VM12). Die Register-Zuordnung (zweiter Fall `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`) gehört in §7. | `spec/spezifikation.md` LH-FA-17.a |
| F-382 (HIGH) | **behoben.** Entschieden in der Spezifikation vor dem Code: `235b9d1` um 16:55, `56732d6` um 17:02, Code `cf38617` um 17:08. `a205cbc` dazwischen ändert `cli.go` nicht. Die Hilfe-Suche läuft vor jeder Prüfung (`cli.go:85`–`:96`, `:210`). Alle Vorgaben aus §6 haben einen Testfall, VM13 bis VM19 und VM23 sind rot, und die Sonde S1 des Reviews ergibt jetzt Hilfe mit Exit-Code 0. Offen bleibt eine neue Randform (V-45). | `internal/adapters/driving/cli/cli.go`, `cli_test.go`, `internal/bootstrap/bootstrap_test.go` |
| F-383 (MEDIUM) | **behoben.** [LH-FA-03](../../spec/lastenheft.md#lh-fa-03--replay-modus) und [LH-FA-17](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) stehen nicht mehr in `abdeckung-vollstaendig.md`. Die Deklarationen auf die Boundary von LH-FA-03 sind entfernt, ebenso die auf die Happy von LH-FA-03 und LH-FA-17. Ein Rest bleibt, mit [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) vereinbar: V-51. | `docs/user/abdeckung-gesamt.md`, `abdeckung-vollstaendig.md` |
| F-384 (MEDIUM) | **behoben für diesen Slice.** Plan §1 grenzt `play` an `slice-v1-abschluss-einspielen` ab. Dort nennen §1 („Übernommen aus …“), §3, Kopf und `Bezug` den Punkt. Die DoD jenes Slice nennt ihn nicht (V-48). | Plan §1; `slice-v1-abschluss-einspielen.md` |
| F-385 (LOW) | **behoben.** `TestReplaySentJeVerbindung`; die frühere Überlebende (M6) ist jetzt rot (VM21). | `replay_unverbraucht_test.go` |
| F-386 (LOW) | **behoben.** Der Kopf nennt `ARC-003`, `ARC-005`, `ARC-009` und `LH-FA-01.a`. Zu `Bezug` siehe V-46. | Plan Kopf |
| F-387 (LOW) | **behoben.** Das Handbuch sagt, dass eine ungültige Variable auch neben der Option `PGR-E2001` ist. | `docs/user/benutzerhandbuch.md` §4 |
| F-388 (LOW) | **behoben.** In `slice-v1-abschluss-betrieb` nennen DoD-Punkte 1 bis 3 und §3 die übernommenen Punkte: `PGR-E4006` vor `PGR-E5002` mit Exit-Code 4, strenge Werte für `--force`, Prüfung jeder gesetzten Variable, Schlüssel `fail_on_unconsumed`, Hilfe vor Prüfung. `Bezug` führt LH-FA-03. Rest V-47. | `slice-v1-abschluss-betrieb.md` |
| F-389 (INFO) | keine Aktion erwartet. `make a-check` meldet 0 Befunde. | — |

---

## 5. Offene Lesart: `--` an der Stelle eines Optionswerts

Gefragt war, ob ein `--` an der Wertstelle (`replay --input -- --help`) die Hilfe-Suche beendet und ob LH-FA-01.a das trägt.

**Für die Hilfe-Suche trägt es der Wortlaut.** LH-FA-01.a sagt „Steht unter den Argumenten vor einem `--` eine Hilfe-Angabe“. Plan §6 sagt „bis zum ersten `--`“. Beide meinen ein Argument, nicht eine Position. Der Code folgt dem: `hilfeVerlangt` bricht bei jedem Argument `--` ab.

**Nicht entschieden ist, was das `--` an der Wertstelle sonst ist.** LH-FA-01.a entscheidet die symmetrische Frage für die Hilfe-Angabe positionsbezogen: „Eine Hilfe-Angabe an der Stelle eines Optionswerts ist die Hilfe, nicht der Wert“. Für `--` an derselben Stelle sagt er nichts. Der Code gibt dem `--` zwei Bedeutungen. Sonde am Binary bei `cf38617`, ohne Umgebungsvariable:

| Aufruf | Ergebnis |
|---|---|
| `replay --listen x --input --` | `Recording [PGR-E3001]: -- nicht lesbar`, Exit-Code 3. `fs.Parse` nimmt `--` als Wert von `--input`. |
| `replay --input -- --help` | `Konfiguration [PGR-E2001]: ungültige Verwendung von replay: flag: help requested`, Exit-Code 2 |
| `replay --listen x --input -- -h` | ebenso, Exit-Code 2 |
| `replay --listen x --input --help` | Hilfe, Exit-Code 0 |

Für die Hilfe-Suche ist das `--` also Ende der Optionen. Für den Optionsparser ist es der Wert von `--input`, und `--help` danach ist für ihn eine Option. Das ist kein „gewöhnliches Argument“, wie LH-FA-01.a es für die Stelle nach `--` zusagt. Die Meldung „help requested“ nennt eine Hilfe, die nicht kommt.

Unter keiner Lesart passt das zusammen:

- **`--` beendet die Optionen** (wörtlich): Dann fehlt `--input` der Wert, und `--help` wäre ein unerwartetes Argument. Der Exit-Code 2 passt dazu, Meldung und Grund nicht.
- **`--` ist an der Wertstelle der Wert** (analog zur Hilfe-Angabe): Dann steht `--help` vor keinem `--`, und es wäre die Hilfe mit Exit-Code 0.

Kein Test fährt `--` an der Wertstelle. Die Randform ist damit im Code entschieden und nicht in LH-FA-01.a. Befund V-45.

---

## 6. Plan gegen Code

### Kopf

- `Berührte Spec-Stellen` deckt den Diff: `LH-FA-01.a`, `LH-FA-03.b`, `LH-FA-13.b`, `LH-FA-17.a`, `SPEC-034`; Port (`ARC-003`), CLI (`ARC-005`) und Composition Root (`ARC-009`) sind nachgetragen.
- `Bezug` nennt LH-FA-03, LH-FA-13 und LH-FA-17, nicht aber LH-FA-01. Die Commits `235b9d1`, `56732d6` und `cf38617` nennen LH-FA-01 (V-46).
- `make kopf-check`: Exit-Code 0. Der Sensor prüft nur Kennungen aus §1 und §2 gegen den Kopf und sieht diese Lücke deshalb nicht.

### §1 Ziel und Abgrenzung

- Ziel und „Geliefert werden“ decken die Option, die Umgebungsvariable, LH-FA-03.b und das Handbuch.
- Die Abgrenzungen sind nachgezogen: `play` an `slice-v1-abschluss-einspielen`; Konfigurationsdatei, Frist und `--force` an `slice-v1-abschluss-betrieb`.
- **Nicht nachgezogen** ist die Lieferung aus `235b9d1`, `56732d6` und `cf38617`: die Hilfe vor jeder Prüfung als abgeschlossene Liste nach LH-FA-01.a, für `record`, `replay`, `version` und global. Ebenso fehlt die Abgrenzung, dass die Hilfe vor der Prüfung von `config show`, `--config` und der Konfigurationsdatei an `slice-v1-abschluss-betrieb` geht. Diese Übernahme steht nur im Folge-Slice. Hard Rule §3.9 verlangt beides in §1 (V-46).

### §3 Plan

Folgt dem Code. Die Zeilen für CLI (Hilfe-Angabe) und Composition Root (Test) sind nachgetragen. Keine Datei im Diff fehlt in §3.

### §6 Risiken

- Das Risiko „Exit-Code erst beim Prozessende“ hat den Ausgang „offen bis Closure“. Das ist richtig vor der Closure.
- Die Randformen „Hilfe neben ungültiger Konfiguration“ und „Formen der Hilfe-Angabe“ stehen mit Vorgaben und Testfällen da. Jeder genannte Testfall existiert.
- Nicht in §6 steht `--` an der Wertstelle (V-45). Ebenso fehlt die Zusage zu `version` und `--version` aus `235b9d1` (V-49).

### Folge-Slices

- **`slice-v1-abschluss-betrieb`:** stimmig nachgezogen. Kopf, `Bezug` (LH-FA-03), §1, DoD-Punkte 1 bis 3 und §3 nennen alle Übernahmen. Ausnahme: Der Kopf nennt LH-FA-01.a, `Bezug` aber nicht LH-FA-01 (V-47).
- **`slice-v1-abschluss-einspielen`:** Kopf, `Bezug`, §1 und §3 nennen die `play`-Regel. Die DoD nennt sie nicht, und der Verifier jenes Slice misst an der DoD (V-48).

---

## 7. Befunde

| ID | Kategorie | Befund | An |
|---|---|---|---|
| V-45 | MEDIUM | `--` an der Stelle eines Optionswerts hat im Code zwei Bedeutungen: Ende der Hilfe-Suche (`cli.go:210`) und Wert der Option für `fs.Parse` (`cli.go:155`). `replay --input -- --help` endet mit `PGR-E2001 … flag: help requested`, Exit-Code 2. `replay --listen x --input --` liest eine Datei `--`. LH-FA-01.a entscheidet die Wertstelle nur für die Hilfe-Angabe, nicht für `--`. Keine der beiden Lesarten ergibt das heutige Verhalten vollständig (Abschnitt 5). Kein Test und kein Punkt in §6. Die Hilfe-Suche selbst trägt der Wortlaut. Muster: Randform im Code entschieden (`BEO-REPO/spec-randform-erst-im-review-entschieden`, Hard Rule §3.12). Die DoD-Behauptung ist nicht verletzt. | Architect (Entscheidung in LH-FA-01.a), danach Implementer (Test) |
| V-46 | LOW | Plan §1 nennt die Lieferung „Hilfe vor Prüfung, abgeschlossene Liste“ (LH-FA-01.a) weder im Ziel noch unter „Geliefert werden“. Er grenzt die Hilfe für `config show`, `--config` und die Konfigurationsdatei nicht an `slice-v1-abschluss-betrieb` ab. `Bezug` führt LH-FA-01 nicht, obwohl der Kopf `LH-FA-01.a` nennt und drei Commits LH-FA-01 nennen. Ob die Lieferung einen eigenen DoD-Punkt braucht, entscheidet der Planner. Die DoD hat einen Liefer-Punkt und damit Raum. Muster: `BEO-REPO/plan-folgt-korrektur-nicht` (Hard Rule §3.9). Das Register zählt diesen Fall in seinem Retirement-Check. | Planner |
| V-47 | LOW | `slice-v1-abschluss-betrieb`: `Berührte Spec-Stellen` nennt seit `a205cbc` `LH-FA-01.a`, `Bezug` aber nicht LH-FA-01. Das ist dasselbe Muster wie der `Bezug`-Teil von F-388. | Planner |
| V-48 | LOW | `slice-v1-abschluss-einspielen` führt den übernommenen Test „`play` kennt `--fail-on-unconsumed` nicht, Variable unbeachtet“ nur in §1 und §3, nicht in einem DoD-Punkt. Das ist die Klasse von F-388 an einem zweiten Folge-Slice. | Planner |
| V-49 | LOW | `235b9d1` fügt in LH-FA-01.a die Zusage ein: „`version` liest weder Umgebungsvariablen noch Konfigurationsdatei; eine Option `--version` gibt es nicht (`PGR-E2001`)“. Kein Test prüft sie, und §6 nennt sie nicht. Das Verhalten stimmt laut Sonde: `--version` und `version --version` ergeben Exit-Code 2, und `version` liest keine Variable. Hard Rule §3.10 verlangt je Zusage einen Test, der rot werden kann. | Implementer |
| V-50 | MEDIUM | **Harness, nicht dieser Slice.** BuildKit überträgt eine Datei des Build-Kontexts nicht neu, wenn Größe und mtime auf die Nanosekunde gleich dem zuletzt übertragenen Stand desselben Pfads sind (Abschnitt 3). `--no-cache` hilft nicht. Gleiche Sekunde genügt, wenn beide mtimes keinen Sekundenbruchteil tragen: nach `tar -x` oder `git archive`, und bei `cp -p`, `rsync -a` oder `touch -r` von solchen Dateien. Ein Mutationslauf über `docker build` in einem wiederverwendeten Pfad kann damit einen Mutanten zu Unrecht überleben lassen oder ein Zurücksetzen verfehlen. Das betrifft die Glaubwürdigkeit jedes „grün gebliebenen“ Mutanten in Reports, die so gearbeitet haben. Abhilfe: ein frischer Pfad je Mutant, `touch` nach jeder Änderung oder Bind-Mount statt Build-Kontext. Die Regel gehört als Lerneintrag in die Closure, eine geschärfte Regel für Mutationsläufe. | Closure-Notiz, Steering Loop |
| V-51 | INFO | Rest von F-383: `LH-FA-17/Boundary` ist durch Tests der Priorität deklariert (`TestParseFailOnUnconsumedUmgebung`, `TestE2EReplayNichtVerbraucht`). Geheimnisse und Anzeige der Konfiguration deckt keiner. [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) lässt Teilfälle zu, und heute steht LH-FA-17 auf „teilweise“, weil die Happy fehlt. Kommt eine Happy-Deklaration hinzu, wird LH-FA-17 „vollständig“, ohne dass Geheimnisse und Anzeige geprüft sind. Das ist der Ort für `slice-v1-abschluss-betrieb`. | Planner von `slice-v1-abschluss-betrieb` |

---

## 8. Urteil

- **DoD-Liefer-Punkt 1** ([LH-FA-03](../../spec/lastenheft.md#lh-fa-03--replay-modus)): **bestätigt**, jede Teilbehauptung gegen LH-FA-03.b, LH-FA-13.b, LH-FA-17.a und LH-FA-01.a einzeln. Alle 23 Mutationen werden aus dem behaupteten Grund rot.
- **`make gates`:** **bestätigt**, selbst gefahren bei `cf38617`, Exit-Code 0. `make kopf-check` ist grün.
- **Review-Findings:** F-382 bis F-388 sind behoben. F-381 wartet auf die Closure. F-389 erwartet keine Aktion.
- **Offen vor der Closure:** V-45 ist eine Randform, die der Architect entscheiden muss. Sie verletzt keine DoD-Behauptung, ist aber nach Hard Rule §3.12 nicht im Code zu belassen, ohne dass LH-FA-01.a sie nennt. V-46 bis V-49 gehen an Planner und Implementer. V-50 gehört in den Lerneintrag.

Häkchen habe ich keine gesetzt.
