# Slice slice-v1-abschluss-einspielen-laufsteuerung: Laufsteuerung beim Einspielen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration), [`LH-FA-14`](../../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0014](../../adr/0014-konfigurationsdatei.md), [ADR-0017](../../adr/0017-einspielen-sequenziell-und-fehlersemantik.md)

**Berührte Spec-Stellen:** `LH-FA-20.a` · `LH-FA-17.a` · `LH-FA-14.a` · `SPEC-017` · `SPEC-034` · `ARC-002` · `ARC-003` · `ARC-005`

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `pgwire-recorder play` liest `--continue-on-error`, `--allow-recorded-errors` und `--finish-session-on-interrupt` aus Kommandozeile, Umgebung und dem Abschnitt `play:`, läuft damit nach einer Fehlerantwort des Servers weiter, lässt einen aufgezeichneten Fehler als erwartet gelten und beendet bei einem Abbruchsignal zuvor die laufende Session, mit den Exit-Codes nach `LH-FA-20.a`.

**Übernimmt:** `slice-v1-abschluss-einspielen` — dessen Teil *Laufsteuerung* (Marke [L] in
dessen §6) nach dem zweiten Schnitt vom 2026-10-09 (Option O3 des Architect, dort §6 Risiko
*Größe des Kerns*; Entscheidung des Nutzers; dort §1, *Abgegeben*). Im Einzelnen:

- die Optionen `--continue-on-error`, `--allow-recorded-errors` und
  `--finish-session-on-interrupt` am allgemeinen Leser, in der Reihenfolge der
  Optionstabelle, mit ihren Umgebungsvariablen und Schlüsseln im Abschnitt `play:`
  (`LH-FA-17.a`, Optionstabelle von `LH-FA-20.a`);
- die Fortsetzung nach einer Fehlerantwort mit `--continue-on-error`: die nächste Interaktion
  läuft, Exit-Code 4 am Ende, je `PGR-E4004` eine Log-Zeile `error` (`LH-FA-20.a` Schritt 6
  und *Meldungen*, `LH-FA-14.a`); ein Verbindungsfehler bricht auch mit der Option ab, und
  der Exit-Code beim Abbruch nach einem früheren `PGR-E4004` ist der des abbrechenden Fehlers
  (`LH-FA-20.a` *Exit-Code*);
- der erwartete Fehler mit `--allow-recorded-errors` (`LH-FA-20.a` *Interaktion*);
- `--finish-session-on-interrupt`: Das erste Signal, auch im Aufbau, beendet das Einspielen
  nach der laufenden Session; der Exit-Code beim Abbruchsignal nach einem früheren Fehler ist
  4 (`LH-FA-20.a` Schritt 7, *Abbruchsignal*);
- aus der DoD jenes Slice bis zum Schnitt: „`--continue-on-error` läuft weiter und endet mit
  Exit-Code 4; `--allow-recorded-errors` lässt aufgezeichnete Fehler zu“ (Punkt 2) und „mit
  `--finish-session-on-interrupt` nach der laufenden Session; der Exit-Code ist 0 ohne
  vorherigen Fehler, sonst 4“ (Punkt 3);
- die Randformen [L] aus §6 jenes Slice (§6 unten);
- aus §7 jenes Slice (*Beobachtungen für Review und Closure*, Befund V-129 aus dessen
  Verifikation) die Frage, wie der Exit-Code des abbrechenden Fehlers nach einem früheren
  `PGR-E4004` entsteht, wenn der Bootstrap den Exit-Code aus der ersten Meldung bildet (§4
  *Start*, §6 *Reihenfolge der Zeilen `error`* und Risiken unten).

Die Aussagen zu den drei Optionen und zu `play` in `README.md` und im Benutzerhandbuch §4
bringt der Slice auf den gelieferten Stand; sie sagen nur zu, was ein Test prüft
(`AGENTS.md` §3.11).

Dieser Slice ist auch die Adresse der Abgrenzung *Einspielen selbst und Fehlersemantik der
Serverfehler* von `slice-v1-abschluss-antwortvergleich` (dort §1), soweit sie die drei
Optionen ohne Vergleich betrifft.

**Aufsetzen.** Der Slice setzt auf dem Kern `slice-v1-abschluss-einspielen` auf: Dort sind
das Kommando `play`, die übrigen Optionen, das Einspielen einfacher Anfragen, der Abbruch bei
jeder Fehlerantwort (`PGR-E4004`, Exit-Code 4) und das Abbruchsignal nach der laufenden
Interaktion geliefert. Bis zu diesem Slice sind die drei Optionen bei `play` unbekannt
(`PGR-E2001`), ihre Schlüssel im Abschnitt `play:` unbekannt (`PGR-E2004`) und ihre
Umgebungsvariablen unbeachtet (§6 jenes Slice, *Optionen der Laufsteuerung*); Tests des
Kerns, die das prüfen, ändert dieser Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Fortsetzung und erwarteter Fehler innerhalb einer Extended-Interaktion (Marke [L·E]: der
  Server verwirft bis `Sync`, übrige Gruppen ohne Warten; eine `error_response` in jeder
  Gruppe zählt für `--allow-recorded-errors`) — `slice-v1-abschluss-einspielen-extended`
  (dort §1, *Übernimmt*); bis dahin spielt `play` keine Extended-Interaktion ein
  (Zwischenstand des Kerns), getestet wird hier mit einfachen Anfragen.
- Die Wirkung der Optionen mit `--compare-responses` (Exit-Code 5, `--allow-recorded-errors`
  ohne Wirkung, Fortsetzung nach einer Abweichung) — `slice-v1-abschluss-antwortvergleich`
  (dort DoD-Punkt 3); sie setzt den Vergleich voraus.
- Anmeldung mit Passwort und TLS zum Server — `slice-v1-abschluss-einspielen-anmeldung` und
  `slice-v1-abschluss-einspielen-tls`; getestet wird gegen einen Server ohne Passwort und
  ohne TLS.
- Code im Upstream-Adapter, im Bootstrap und im PGWire-Adapter — Schicht-Abgrenzung: Die
  Optionen liest der CLI-Adapter; Fortsetzung, erwarteten Fehler, Session-Ende und das
  Ergebnis des Laufs, aus dem der Bootstrap den Exit-Code bildet, entscheidet der
  Play-Service. Der Kern reicht die Optionen von `play` und die Signale an den Play-Service
  (dort §3). Braucht dieser Slice eine Änderung im Bootstrap oder im Upstream-Adapter, ist
  das eine dritte Schicht und die Rückführung in §4. Der Driving Port `Player` (Schicht
  *Ports*, `ARC-003`, nicht *Application*) ändert sich nur im Kommentar, der den Vertrag von
  `Play` nennt, den der Play-Service erfüllt (`AGENTS.md` §3.7); Signatur und Importe bleiben.
  Das zählt nicht als dritte Schicht: Die Größenregel misst, ob der Diff in einer
  Review-Sitzung prüfbar ist, und ein Kommentar, der dem Service folgt, ändert daran nichts;
  eine Änderung an Signatur oder Typen eines Ports wäre die Rückführung in §4.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung): `--continue-on-error`, `--allow-recorded-errors` und
      `--finish-session-on-interrupt` wirken bei `play` aus Kommandozeile, Umgebung und dem
      Abschnitt `play:` mit der Priorität und den booleschen Werten nach `LH-FA-17.a`; mit
      `--continue-on-error` läuft das Einspielen nach einer Fehlerantwort mit der nächsten
      Interaktion weiter und endet mit Exit-Code 4, je `PGR-E4004` eine Log-Zeile `error`;
      ein Verbindungsfehler bricht auch mit der Option ab, nach einem früheren `PGR-E4004`
      mit dem Exit-Code des abbrechenden Fehlers (4 oder 6) (Integrationstest).
- [ ] [`LH-QA-05`](../../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler): Mit `--allow-recorded-errors` gilt eine Fehlerantwort als erwartet, wenn
      die aufgezeichnete Interaktion eine `error_response` trägt: keine Meldung, kein
      Abbruch, kein Beitrag zum Exit-Code; trägt sie keine, ist es `PGR-E4004` wie ohne
      Option (Test).
- [ ] `SIGINT` und `SIGTERM` beenden das Einspielen mit `--finish-session-on-interrupt` nach
      der laufenden Session, ein erstes Signal im Aufbau nach der ganzen Session; der
      Exit-Code beim Abbruchsignal ist 4 nach einem früheren Fehler, sonst 0 (Test). Beleg in
      §7 für Punkt 1 bis 3: je Zusage Zusage · Mutation · roter Test (`AGENTS.md` §3.10).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/cli` | update | Optionen `--continue-on-error`, `--allow-recorded-errors` und `--finish-session-on-interrupt` am allgemeinen Leser mit Umgebung und Schlüsseln im Abschnitt `play:`, je ein Feld in `Einspielvorgaben`, Hilfetext von `play` (sagt nur zu, was ein Test prüft: kein Exit-Code bei `--continue-on-error`, erwartet nur eine Fehlerantwort mit dem Schweregrad `ERROR`; Review F-565); ersetzt den Stand *unbekannt* des Kerns |
| `internal/adapters/driving/cli` (Tests) | update | `TestParsePlayLaufsteuerung` neu; Optionstabelle, Reihenfolge und Hilfe in den Tests des Lesers nachgezogen; die Fälle des Stands *unbekannt* (Umgebungsvariablen unbeachtet, Schlüssel `continue_on_error` ungültig) entfernt |
| `internal/hexagon/services` (Play-Service) | update | dieselben Felder in derselben Reihenfolge in `PlayOptions` (der Bootstrap des Kerns konvertiert den Wert, ohne ein Feld zu nennen); Fortsetzung nach `PGR-E4004`, erwarteter Fehler, Session-Ende nach dem Signal; `Play` liefert die Fehler des Laufs per `errors.Join`, den abbrechenden zuerst, danach die früheren in der Reihenfolge ihres Auftretens (§6, *Reihenfolge der Zeilen `error`*) |
| `internal/hexagon/ports/driving` | update | nur der Kommentar des Ports `Player`: Ende mit `--finish-session-on-interrupt` und die Fehler des Laufs als Zusammenfassung |
| `internal/hexagon/services` (Tests) | update | `play_laufsteuerung_test.go`: Fortsetzung, Abbruch nach früheren Fehlern, erwarteter Fehler, Signal im Aufbau, in und nach der letzten Interaktion mit `--finish-session-on-interrupt`, Signal nach einer Fehlerantwort, zweites Signal nach einer Fehlerantwort, zweites Signal mit `--finish-session-on-interrupt` in der Interaktion, nach früherem Fehler und im Aufbau (`TestPlayZweitesSignalFinishSession`, Review F-564); je Zusage eine Mutation |
| `test/integration` | update | `play_laufsteuerung_e2e_test.go`: Happy/Boundary/Negative nach LH-FA-20 für die drei Optionen aus Kommandozeile, Umgebung und `play:`, Abbruch mit `PGR-E6001` und `PGR-E4003` nach früherem `PGR-E4004`, `SIGINT`/`SIGTERM` am Binary, zwei Signale am Binary ohne und mit `--finish-session-on-interrupt` (`TestE2EPlayZweitesSignal`, Review F-564), gegen einen Server ohne Passwort und ohne TLS |
| `docs/user/abdeckung-*.md` | update | erzeugt mit `make abdeckung` |
| `docs/user/benutzerhandbuch.md` | update | §4 *Eine Aufzeichnung in eine Datenbank einspielen*, *Hinweise*: Exit-Code nach einem abbrechenden Fehler bei `--continue-on-error`, `FATAL` bei `--allow-recorded-errors`, zweites Signal; nur Geprüftes (Review F-566) |
| `README.md` | update | *Was kann ich heute tun?*: `play` mit dem gelieferten Stand statt „fehlt noch“, ohne Zusage für vorbereitete Anweisungen, Passwort-Anmeldung und TLS |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen` liegt in `done/`
(Kommando `play`, Einspielen einfacher Anfragen, Abbruch bei einer Fehlerantwort,
Abbruchsignal). Schritt 9 der Reihenfolge in §5 von
[welle-v1-abschluss](../welle-v1-abschluss.md) (zweiter Schnitt vom 2026-10-09). Die
Randformen aus §6 entschied der Architect am 2026-10-09 vor dem Code in `LH-FA-20.a`; vor dem
ersten Code-Commit prüft er die Liste gegen den gelieferten Kern, besonders gegen die Form,
in der der Kern Optionen und Signale an den Play-Service reicht (`AGENTS.md` §3.12). Dazu prüft er
vor dem Code die Randform *Reihenfolge der Zeilen `error`* (§6), übergeben von
`slice-v1-abschluss-einspielen` (dort §7, *Beobachtungen für Review und Closure*; V-129 aus
dessen Verifikation): Der gelieferte Bootstrap bildet den Exit-Code aus der **ersten** Meldung
des Fehlers, den der Play-Service liefert, und schreibt die Zeilen `error` erst nach dem
Einspielen, in dieser Reihenfolge. `LH-FA-20.a` *Exit-Code* verlangt nach einem früheren
`PGR-E4004` den Code des **abbrechenden** Fehlers (4 oder 6); Reihenfolge und Zeitpunkt der
Zeilen `error` legen weder `LH-FA-20.a` *Meldungen* noch `LH-FA-14.a` fest. Geprüft und
entschieden am 2026-10-10 (§6, *Reihenfolge der Zeilen `error`*): ohne Code im Bootstrap
umsetzbar, die Rückführung *dritte Schicht* tritt nicht ein.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer
  Review-Sitzung prüfbar, oder der Slice braucht eine dritte Schicht (§1). Schnitt dann:
  `--finish-session-on-interrupt` als eigener Slice; Fortsetzung und erwarteter Fehler sind
  ohne ihn lieferbar.
- `in-progress` → `open` (blockiert — Carveout?): Der gelieferte Kern entscheidet das Ende
  nach einem Signal außerhalb des Play-Service, sodass das Session-Ende dort nicht
  entschieden werden kann; dann zuerst die Entscheidung des Architect, wo es liegt.
- `in-progress` → `next` (dritte Schicht, Bedingung oben): Die Entscheidung zur Randform
  *Reihenfolge der Zeilen `error`* (§6) lässt sich nur mit Code im Bootstrap umsetzen, etwa
  weil die Zeilen in zeitlicher Reihenfolge stehen müssen und der Exit-Code dann nicht mehr
  aus der ersten Meldung folgt. Schnitt dann: `--continue-on-error` mit dem Exit-Code nach
  einem früheren `PGR-E4004` als eigener Slice mit Bootstrap und Play-Service;
  `--allow-recorded-errors` und `--finish-session-on-interrupt` sind ohne ihn lieferbar.

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

**Randformen** (`AGENTS.md` §3.12) — die mit Marke [L] aus §6 von
`slice-v1-abschluss-einspielen`, geprüft vom Architect am 2026-10-09 vor dem ersten
Code-Commit jenes Slice und für den zweiten Schnitt markiert. Neu entschieden heißt: im Commit
jener Prüfung in `LH-FA-20.a`, sonst an der genannten Stelle. Die Teile mit [L·E] liefert
`slice-v1-abschluss-einspielen-extended` (§1). Vom Architect am 2026-10-10 vor dem ersten
Code-Commit dieses Slice gegen den gelieferten Kern geprüft (Bootstrap konvertiert
`cli.Einspielvorgaben` per Typumwandlung in `services.PlayOptions`, reicht erstes und zweites
Signal als `ctx` und `ablauf` an den Play-Service, schreibt je Meldung des gelieferten Fehlers
nach dem Einspielen eine Zeile `error` und bildet den Exit-Code aus der ersten); „neu
entschieden am 2026-10-10“ heißt: im Commit dieser Prüfung in `LH-FA-20.a`. Offen ist keine.

- **Optionen der Laufsteuerung** [L] — am allgemeinen Leser, Abschnitt `play:`; bestätigt,
  `LH-FA-17.a` und `LH-FA-20.a` (Optionstabelle).
- **Erwarteter Fehler bei `--allow-recorded-errors`** [L] — entscheidet allein, ob die
  aufgezeichnete Interaktion eine `error_response` trägt; keine Meldung; neu entschieden in
  `LH-FA-20.a` *Interaktion*. „In jeder Gruppe“ ist [L·E].
- **Fortsetzung** (`--continue-on-error`, erwarteter Fehler) [L] — Rest der Interaktion wie
  aufgezeichnet, Exit-Code 4 am Ende nach einem `PGR-E4004`; bestätigt, Schritt 6 und
  Tabelle. Bei einer Extended-Interaktion ist sie [L·E].
- **Je `PGR-E4004` bei `--continue-on-error` eine Log-Zeile** [L] — `error` mit `code` und
  `error`, nennt `id`, `sequence`, SQLSTATE und `M` der Fehlerantwort, keine weiteren Felder,
  nie Passwort oder Parameterwerte; neu entschieden in `LH-FA-20.a` *Meldungen* und
  `LH-FA-14.a`.
- **Erstes Signal im Aufbau mit `--finish-session-on-interrupt`** [L] — der Aufbau läuft zu
  Ende, ein Fehler darin zählt, danach die ganze Session; neu entschieden in `LH-FA-20.a`
  *Abbruchsignal*.
- **Exit-Code beim Abbruchsignal nach einem früheren Fehler** [L] — 4; bestätigt,
  Tabellenzeile *Abbruchsignal*.
- **Exit-Code beim Abbruch nach einem früheren `PGR-E4004`** [L] — der des abbrechenden
  Fehlers (4 oder 6); neu entschieden in `LH-FA-20.a` *Exit-Code*.
- **Reihenfolge der Zeilen `error`** [L] — übergeben von `slice-v1-abschluss-einspielen`
  (dort §7, V-129 aus dessen Verifikation); neu entschieden am 2026-10-10 in `LH-FA-20.a`
  *Meldungen* und *Exit-Code*: Die Zeilen stehen nach dem Ende des Einspielens (`time` ist der
  Zeitpunkt des Schreibens), der abbrechende Fehler zuerst, danach die früheren, nach denen
  das Einspielen weiterlief, in der Reihenfolge ihres Auftretens; ohne Abbruch alle in dieser
  Reihenfolge. Der Exit-Code ist der der Klasse der ersten Zeile, ohne Zeile `0`. Die Fehler
  eines Laufs sind keine gleichrangigen Fehler (`SPEC-034`). Umsetzbar ohne Bootstrap: Der
  Play-Service liefert die Fehler als reine Zusammenfassung in dieser Reihenfolge, der
  gelieferte Bootstrap schreibt sie und nimmt den Code der ersten.
- **Fehlerantwort vor dem zweiten Signal** [L] — eine Fehlerantwort, die mit
  `--continue-on-error` eintraf, bevor das zweite Signal ihre Interaktion unterbricht, ist eine
  Zeile und zählt; die Unterbrechung selbst nicht; neu entschieden am 2026-10-10 in
  `LH-FA-20.a` *Meldungen*.
- **Abbrechender Fehler nach dem ersten Signal** [L] — im laufenden Aufbau, in der laufenden
  Interaktion oder mit `--finish-session-on-interrupt` in der laufenden Session: Exit-Code
  seiner Klasse, nicht die Zeile *Abbruchsignal*; diese zählt Fehler bis zum Ende des
  Einspielens. Neu entschieden am 2026-10-10 in `LH-FA-20.a` *Exit-Code* und Tabellenzeile
  *Abbruchsignal* (der Kern verhält sich schon so, ohne dass eine Regel es sagte).
- **Fehlerregeln in der zu Ende laufenden Session; Signal zwischen Sessions** [L] — mit
  `--finish-session-on-interrupt` gelten für die übrigen Interaktionen die Fehlerregeln wie
  ohne Signal; trifft das Signal zwischen zwei Sessions ein, beginnt keine Session mehr, auch
  mit der Option; neu entschieden am 2026-10-10 in `LH-FA-20.a` *Abbruchsignal*. Ein Signal im
  Start: keine Session, bestätigt, *Start*.
- **Mehrere Fehlerantworten in einer Interaktion** [L] — beim Weiterlesen bis zum
  `ReadyForQuery` ist jede eine Zeile (*Meldungen*, „je Fehlerantwort eine“), erwartet ist
  jede, wenn die Aufzeichnung der Interaktion eine `error_response` trägt (*Interaktion*);
  `FATAL` oder `PANIC` ist `PGR-E4003`, auch bei erwartetem Fehler (Tabelle, „immer“);
  bestätigt.
- **Verbindungsende beim Weiterlesen nach einer Fehlerantwort** [L] — `PGR-E4003` bricht ab,
  Exit-Code 4, die Fehlerantwort davor ist eine frühere Zeile; bestätigt, Tabelle und
  *Meldungen*.
- **Kombinationen der Optionen** [L] — `--continue-on-error` mit `--allow-recorded-errors`:
  ein erwarteter Fehler zählt nicht, jeder andere läuft weiter und zählt; ein erwarteter
  Fehler ohne `--continue-on-error` läuft ebenso weiter (Schritt 6);
  `--finish-session-on-interrupt` ohne Signal ohne Wirkung; bestätigt, Tabelle und Schritt 6.
  Mit `--compare-responses`: `slice-v1-abschluss-antwortvergleich` (§1 oben).
- **Quellen und ungültige Werte** [L] — Kommandozeile, Umgebung, `play:`; boolesch nach
  `LH-FA-17.a` (ohne Wert `true`, `=true` oder `=false`, jeder andere Wert `PGR-E2001`, in der
  Datei `PGR-E2004`; leere Umgebungsvariable nicht gesetzt; eine ungültige Umgebungsvariable
  auch bei gesetzter Option `PGR-E2001`); bestätigt, `LH-FA-17.a` und Optionstabelle.
- **Akzeptierte Negative** — (a) Die Zeilen `error` erscheinen erst am Ende des Laufs, nicht
  beim Fehler: Sofortiges Schreiben verlangte einen Port für Meldungen und Code im Bootstrap
  (dritte Schicht, §4) für eine Ausgabe, die das Lastenheft nicht verlangt (`LH-FA-14`,
  `LH-FA-20`); `time` ist kein Vertragsschlüssel (`LH-FA-14.a` *Zeilenform*). (b) Der
  Play-Service hält die Fehler bis zum Ende; ihr Speicher wächst mit ihrer Zahl, höchstens
  einer je Fehlerantwort einer Aufzeichnung, die ohnehin ganz im Speicher liegt. (c) Ein
  Fehler, der mit dem zweiten Signal zusammenfällt, kann als Unterbrechung gelten (Bestand
  des Kerns, *Abbruchsignal*). (d) Der Hilfetext von `play` nennt die drei Optionen; sein
  Wortlaut ist kein Vertrag (`LH-FA-01.a`) und sagt nur zu, was ein Test prüft (`AGENTS.md`
  §3.11). (e) Ein abbrechender Fehler ohne Meldungscode nach einem früheren `PGR-E4004` hinge
  nach `SPEC-034` *Ausgabe* als Ursache an die erste Meldung, Exit-Code 4 statt 1; entschieden
  vom Architect am 2026-10-10 auf Rückgabe des Implementers (§7): keine Regel in `LH-FA-20.a`
  oder `SPEC-034`, kein Code, kein Test. Grund: Am Code `de00222` entsteht er nicht. Laden und
  `Validate` laufen im Start vor jeder Interaktion (Fehler beim Laden ohne Code ist dort der
  Fall `PGR-E1000` von `SPEC-034`, ohne früheren Fehler); die Fehler des Play-Service selbst
  tragen einen Code, `eingeordnet` behält ihn; der Upstream-Adapter liefert nur klassifizierte
  Fehler, auch beim abgebrochenen Aufbau; `ctx` (erstes Signal) erreicht den Adapter nicht, und
  Fehler nach dem zweiten Signal verwirft `Play`. Ausgeschlossen ist der Fall durch den
  Vertrag des Ports `Einspielziel` und `EinspielSession` (je Fehler sein Code), an den auch
  `slice-v1-abschluss-einspielen-anmeldung` und `slice-v1-abschluss-einspielen-tls` mit
  neuen Fehlern im Aufbau gebunden sind; bräche ein Adapter ihn, bliebe der Text als Ursache
  sichtbar und der Exit-Code ungleich 0. Eine eigene Regel für `play` teilte die allgemeine
  Regel von `SPEC-034`.
- **Ablösung des Stands *unbekannt*** — bis zu diesem Slice sind die drei Optionen bei `play`
  unbekannt, ohne eigene Regel (§6 von `slice-v1-abschluss-einspielen`, *Optionen der
  Laufsteuerung*); dieser Slice ändert die Tests des Kerns dazu, die Spezifikation nicht.

**Risiken:**

- Der gelieferte Kern reicht Optionen oder Signale nicht so an den Play-Service, dass dieser
  Slice ohne Änderung im Bootstrap auskommt; dann berührt er drei Schichten (§1, §4). Vorab
  geprüft am 2026-10-10: Optionen per Typumwandlung, Signale als `ctx` und `ablauf` (§6) —
  **Ausgang:** offen bis Closure.
- Die Entscheidung zur Randform *Reihenfolge der Zeilen `error`* verlangt Code im Bootstrap
  (übergeben von `slice-v1-abschluss-einspielen`, V-129): Ohne ihn endet `--continue-on-error`
  nach `PGR-E4004` und einem abbrechenden `PGR-E6001` mit Exit-Code 4 statt 6, oder der
  Slice trifft seinen Ausschluss *Code im Bootstrap* und die Rückführung in §4. Entschieden am
  2026-10-10 so, dass der gelieferte Bootstrap genügt (§6) —
  **Ausgang:** offen bis Closure.
- Ein Signal zwischen einer Fehlerantwort und der Fortsetzung: Der Exit-Code 4 folgt aus der
  Tabellenzeile *Abbruchsignal*, belegt ist er erst mit einem Test, der das Signal genau dort
  setzt — **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
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

**Belege des Implementers** (Code-Commit `1b63ef8` auf `85caec4`):

*Größe und Schichten.* `git diff --shortstat 85caec4 1b63ef8`: 13 Dateien, +881 −54. Produkt-Code
+100 −31 (`services/play.go` +72 −26, CLI `cli.go` +23 −3, Kommentar des Ports `Player` +5 −2),
Tests +753 −17, erzeugt (`make abdeckung`) +21 −2, Plan §3 +7 −4. Zwei Schichten (CLI-Adapter,
Play-Service mit dem Kommentar seines Ports); `git diff --stat 85caec4 1b63ef8 --
internal/bootstrap internal/adapters/driven internal/adapters/driving/pgwire
internal/hexagon/model cmd` ist leer. Die Reihenfolge der Zeilen `error` trägt `errors.Join` im
Play-Service; der gelieferte Bootstrap schreibt je Meldung eine Zeile und nimmt den Code der
ersten. Keine Rückführung aus §4 ist eingetreten.

*Randformen.* Keine Operation des Diffs entscheidet eine Randform außerhalb von §6; jede
Zusage unten zeigt auf ihre Zeile in §6. Eine Frage an den Architect, die §6 nicht nennt, steht
am Ende dieses Abschnitts.

*Weg der Mutanten.* Je Mutant eine frische Kopie aller Dateien von `git ls-files -co
--exclude-standard` (Python `shutil.copy`, ohne Übernahme der mtime, gleichwertig `cp -r` ohne
`-p`) unter dem Scratch-Verzeichnis, genau eine Ersetzung im Wortlaut (gofmt-sauber, die Stufe
`test` meldet keine nicht formatierte Datei), dann `make test` oder `make test-integration` in
der Kopie, die Kopie danach gelöscht; die Integrationsläufe räumen Container, Netz und Volume
selbst ab (`docker ps -a` und `docker network ls` mit `pgr-it` danach leer). Rot heißt: der
genannte Test schlug mit der genannten Meldung fehl, nicht Build, Vet oder gofmt, und nicht
über eine Gesamtfrist. Alle 21 Unit- und 6 E2E-Mutanten wurden rot; kein grüner Mutant.

| Zusage (§6) | Mutation | roter Test |
|---|---|---|
| `--continue-on-error` setzt sein Feld (*Optionen der Laufsteuerung*) | setzt `AllowRecordedErrors` | `TestParsePlayLaufsteuerung/continue-on-error` (Feld `false` bei `--continue-on-error`, `=true`, Datei, Umgebung), `TestLeserAlleOptionen` |
| `--allow-recorded-errors` ist boolesch, ohne Wert `true` (*Quellen und ungültige Werte*) | Art Text statt Wahrheitswert | `TestParsePlayLaufsteuerung/allow-recorded-errors` (ohne Wert `PGR-E2001`; `=1` angenommen), `TestLeserAlleOptionen` |
| Standard `false` (Optionstabelle) | Standard von `--finish-session-on-interrupt` `true` | `TestLeserAlleOptionen` („ohne … erwartet Default“), `TestParsePlay`, `TestParsePlayLaufsteuerung` (ohne Quelle) |
| `=false` ist `false` (*Quellen und ungültige Werte*) | `v != ""` statt `v == "true"` | `TestLeserAlleOptionen` („true und false ergeben dasselbe“), `TestParsePlayLaufsteuerung` |
| Hilfe nennt die Option (akzeptiertes Negativ (d)) | Zeile `--finish-session-on-interrupt` umbenannt | `TestParsePlayHilfe` |
| Fortsetzung nach `PGR-E4004` (*Fortsetzung*) | `if !ContinueOnError` → `if true` | `TestPlayFortsetzung` (nur die erste Meldung), `TestPlayFortsetzungAbbruch`, `TestPlaySignalNachFehlerantwort`; E2E `TestE2EPlayFortsetzung` (alle drei Quellen: eine statt zwei Zeilen `PGR-E4004`), `TestE2EPlayFortsetzungAbbruch`, `TestE2EPlayFinishSession/Umgebung_nach_früherem_Fehler` |
| Rest der Interaktion bis `ReadyForQuery` gelesen (*Fortsetzung*; *Mehrere Fehlerantworten*) | nach `weiter` `return nil` | `TestPlayFortsetzung` (Ablauf: `anfrage B` vor dem `ReadyForQuery` von A; die zweite Fehlerantwort von A fehlt), `TestPlayFortsetzungAbbruch` |
| Abbrechender Fehler zuerst, Exit-Code seiner Klasse (*Reihenfolge der Zeilen `error`*, *Exit-Code beim Abbruch nach einem früheren `PGR-E4004`*) | `errors.Join(append(frueher, err)...)` | `TestPlayFortsetzungAbbruch` (alle vier Fälle; `PGR-E6001` mit Exit-Code 4 statt 6); E2E `TestE2EPlayFortsetzungAbbruch/PGR-E6001` (Exit-Code 4 statt 6, Zeilen `PGR-E4004` vor `PGR-E6001`), `/PGR-E4003` |
| Frühere `PGR-E4004` nach einem Abbruch gemeldet (*Meldungen*) | `errors.Join(err)` | `TestPlayFortsetzungAbbruch` (alle vier Fälle, nur die abbrechende Meldung) |
| Frühere `PGR-E4004` ohne Abbruch gemeldet, Exit-Code 4 (*Fortsetzung*) | `return nil` statt `errors.Join(frueher...)` | `TestPlayFortsetzung`, `TestPlayErwarteterFehler/mit_--continue-on-error`, `TestPlayFinishSessionFehler/mit_--continue-on-error`, `TestPlaySignalNachFehlerantwort`, `TestPlayZweitesSignalNachFehlerantwort` |
| Frühere in der Reihenfolge ihres Auftretens (*Reihenfolge der Zeilen `error`*) | neue vor die alten gestellt | `TestPlayFortsetzung` (d, c, a2, a1), `TestPlayFortsetzungAbbruch/PGR-E6001_in_Session_2` (b vor a) |
| Fehlerantwort vor dem zweiten Signal ist eine Zeile, die Unterbrechung keine (*Fehlerantwort vor dem zweiten Signal*) | nach dem zweiten Signal `return nil` | `TestPlayZweitesSignalNachFehlerantwort` (keine Meldung, Exit-Code 0) |
| erwartet nur mit `--allow-recorded-errors` (*Erwarteter Fehler*) | `erwartet := mitFehlerantwort(in)` | `TestPlayErwarteterFehler/ohne_Option` |
| erwartet nur mit aufgezeichneter `error_response`; sonst `PGR-E4004` wie ohne Option (*Erwarteter Fehler*) | `erwartet := AllowRecordedErrors` | `TestPlayErwarteterFehler/ohne_aufgezeichnete_error_response`, `/mit_--continue-on-error`; E2E `TestE2EPlayErwarteterFehler` (Exit-Code 0 ohne Zeile statt 4) |
| an irgendeiner Stelle der Aufzeichnung, nicht SQLSTATE (*Erwarteter Fehler*) | nur die erste aufgezeichnete Antwort geprüft | `TestPlayErwarteterFehler/erwartet`, `/mit_--continue-on-error`; E2E `TestE2EPlayErwarteterFehler` (alle drei Quellen: Exit-Code 4) |
| `FATAL` bleibt `PGR-E4003`, auch erwartet (*Mehrere Fehlerantworten*) | Prüfung auf erwartet vor `FATAL` | `TestPlayErwarteterFehler/FATAL` (keine Meldung) |
| erwarteter Fehler ohne Meldung (*Erwarteter Fehler*) | erwarteter Fehler an `weiter` gereicht | `TestPlayErwarteterFehler/erwartet`, `/mit_--continue-on-error` |
| mit `--finish-session-on-interrupt` die übrigen Interaktionen der Session (*Erstes Signal im Aufbau*; Schritt 7) | Option in `session` nicht beachtet | `TestPlayFinishSession` (Signal in der Interaktion), `TestPlayFinishSessionFehler` (beide); E2E `TestE2EPlayFinishSession/Kommandozeile`, `/Umgebung_nach_früherem_Fehler` (Wirkung `1` statt `1,2`) |
| nach dem Signal keine weitere Session, auch mit der Option (*Signal zwischen Sessions*) | Option auch in der Schleife der Sessions | `TestPlayFinishSession` (Session 2 läuft), `TestPlayFinishSessionFehler/mit_--continue-on-error` |
| Fehlerregeln in der zu Ende laufenden Session (*Fehlerregeln in der zu Ende laufenden Session*) | Fehler nach dem ersten Signal verworfen | `TestPlayFinishSessionFehler/ohne_--continue-on-error` (Exit-Code 0 statt 4), `TestPlayErstesSignal` |
| Signal nach einer Fehlerantwort: Exit-Code 4 (*Exit-Code beim Abbruchsignal nach einem früheren Fehler*; Risiko 3) | nach dem ersten Signal frühere Fehler verworfen | `TestPlaySignalNachFehlerantwort`, `TestPlayFinishSessionFehler/mit_--continue-on-error`; E2E `TestE2EPlayFinishSession/Umgebung_nach_früherem_Fehler` (Exit-Code 0 statt 4) |

Nicht eigens mutiert, weil der allgemeine Leser es für jede Option trägt und `TestLeserAlleOptionen`
die drei neuen Optionen mit ihrem Eintrag in der Optionstabelle des Tests einschließt: Vorrang
Kommandozeile vor Umgebung vor Datei, `PGR-E2001` für eine ungültige Umgebungsvariable auch neben
der Option, `PGR-E2004` mit dem Schlüssel und ohne den Wert, Schlüssel auf der falschen Ebene.
Am Binary belegt `TestE2EPlayFortsetzung` die Umgebungsvariable mit `1` (`PGR-E2001`) und den
Schlüssel mit `1` (`PGR-E2004`).

*Läufe.* Am Arbeitsbaum vor dem Code-Commit: `make test` grün, `make test-integration` grün
(darunter `TestE2EPlayFortsetzung`, `TestE2EPlayFortsetzungAbbruch`, `TestE2EPlayErwarteterFehler`,
`TestE2EPlayFinishSession` mit je allen Unterfällen), `make lint` grün, `make abdeckung`
geschrieben. `make gates` am sauberen Stand `1b63ef8`: grün (Exit-Code 0, Nachweis `record-gates` gestempelt; darunter `make test-integration` „run-integration-tests: gruen“, `make lint` und `make lint-gegenprobe` grün, d-check 444 Dateien ohne Befund, `make kopf-check`, `make abdeckung-check`, `make a-check` und die Gegenproben grün).

*Handbuch und Hilfe.* Das Benutzerhandbuch beschreibt die drei Optionen schon im Zielstand
(§4 *Hinweise*, Optionstabelle in §5); der Diff ändert es nicht, und jede seiner Aussagen zu den
drei Optionen prüft ein E2E-Test oben (Fortsetzung mit Exit-Code 4, erwarteter Fehler,
Session-Ende nach `SIGINT`/`SIGTERM`). Der Hilfetext von `play` sagt nur zu, was die Tests oben
prüfen: Fortsetzung mit der nächsten Anfrage und Exit-Code 4, erwarteter Fehler bei
aufgezeichneter Fehlerantwort, Ende nach der laufenden Session, die Umgebungsvariablen
(`TestParsePlayHilfe`).

*Frage an den Architect* (Randform · Frage, nicht entschieden, kein Code dazu): Ein abbrechender
Fehler ohne Meldungscode nach einem früheren `PGR-E4004` · Der Port `Einspielziel` sagt für
jeden Fehler einen Code zu, und der Upstream-Adapter liefert nur klassifizierte Fehler; käme doch
einer ohne Code, hinge `model.Meldungen` ihn als Ursache an die erste klassifizierte Meldung
(`SPEC-034` *Ausgabe*), und der Exit-Code wäre 4 statt 1. Genügt der Port-Vertrag als
Ausschluss, oder braucht `LH-FA-20.a` *Exit-Code* eine Zeile dazu? — Entschieden in §6,
akzeptiertes Negativ (e).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `19f5512` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/slice-waechst-durch-uebernahmen`, `BEO-REPO/rueckfuehrung-ohne-verzeichniswechsel`
  und `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` — betreffen den Schnitt des Gebers und
  werden dort gezählt (§8 von `slice-v1-abschluss-einspielen`). Nachgezählt beim Eintragen
  der Sendung: drei Liefer-Punkte, zwei Schichten (CLI-Adapter, Play-Service).
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (`AGENTS.md` §3.12),
  `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`,
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (§3.10),
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (§3.11),
  `BEO-REPO/plan-folgt-korrektur-nicht` (§3.9) und
  `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (§3.13) — verkörpert; die Randformen in §6
  sind vor dem Code entschieden, je Zusage eine Mutation mit Beleg in §7, der Geber
  `slice-v1-abschluss-einspielen` zeigt im selben Commit hierher, und
  `slice-v1-abschluss-einspielen-extended` nennt die Teile [L·E] mit der Kennung dieses
  Slice.

Nachgezählt beim Eintragen der Frage aus V-129 (Closure von `slice-v1-abschluss-einspielen`,
`AGENTS.md` §3.13): Sie ist eine Randform und ein Risiko, kein Liefer-Punkt; DoD-Punkt 1 sagt
den Exit-Code des abbrechenden Fehlers schon zu. Drei Liefer-Punkte, zwei Schichten
(CLI-Adapter, Play-Service), solange die Entscheidung ohne Code im Bootstrap auskommt; sonst
die Rückführung *dritte Schicht* in §4.

Nachgezählt bei der Vorab-Prüfung des Architect am 2026-10-10: Die neu entschiedenen
Randformen in §6 sind Randformen, keine Liefer-Punkte; keine neue Zuweisung an einen anderen
Slice. Drei Liefer-Punkte, zwei Schichten (CLI-Adapter, Play-Service), kein Code im Bootstrap.

Nachgezählt bei der Rückgabe des Implementers am 2026-10-10: Der Kommentar des Driving Ports
`Player` liegt in der Schicht *Ports*, nicht im Play-Service; §1 nennt ihn jetzt und warum er
keine dritte Schicht ist. Drei Liefer-Punkte, zwei Schichten mit Code (CLI-Adapter,
Play-Service), kein Code im Bootstrap.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
