# Slice slice-v1-abschluss-umgebung-anderes-kommando: Die Umgebungsvariable der Option eines anderen Kommandos ist `PGR-E2001`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein Abnahmeszenario
hängt an ihm. Eingesammelt wird er von der nächsten Welle-Closure. Eingeordnet am 2026-10-12 nach
Entscheidung des Nutzers vom 2026-10-12 (Ziel: ein Release): als Nächstes, direkt hinter den Slices von
welle-v1-abschluss bis `slice-v1-abschluss-zeitangaben-doku` und vor `slice-v1-abschluss-antwortvergleich`
(WIP-Limit 1); die Gate-Kette (`slice-harness-gate-kette-schnell`,
`slice-harness-gate-kette-stop-hook-liste`, `slice-harness-gate-kette-stop-hook-freigabe`) steht seit
derselben Entscheidung nach dem Release, nicht davor. Reihenfolge in §4 *Start* und im Drift-Log der Roadmap. **Wellenlos** ist Entscheidung des Nutzers vom 2026-10-12
(„Deine Empfehlungen“; die Empfehlung des Planners, ihn nicht in `welle-v1-abschluss` zu führen, weil
seine Closure-Bedingung keine andere als die DoD ist).

**Bezug:** [`LH-FA-17`](../../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)
(Scope), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus) (die Option
`--fail-on-unconsumed` gibt es nur bei `replay`), [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung)
(die Hilfe geht jeder Prüfung vor).

**Berührte Spec-Stellen:** `LH-FA-17.a` (*Einstellungen*: Absatz zur Umgebungsvariable der Option
eines anderen Kommandos, Reihenfolge der Fehler, Optionstabelle; Entscheidung des Nutzers vom
2026-10-12, Spezifikation seit Commit `9bf8f06`) · `LH-FA-03.b` (*Andere Kommandos*) · `LH-FA-01.a`
(Hilfe und `version` lesen nichts) · `SPEC-034` (Zeile `PGR-E2001` des Katalogs nennt die neue Ursache seit der Prüfung des
Architect vom 2026-10-12; die Änderung ist nicht Teil dieses Slice)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der CLI-Adapter lehnt beim Start jede gesetzte, nicht leere Umgebungsvariable einer Option
mit `PGR-E2001` ab, die das Kommando laut Spalte *Kommando* der Optionstabelle (`LH-FA-17.a`) nicht
kennt — gleich welchen Wert sie hat —, und die Meldung nennt den Namen der Variable, nie ihren Wert;
das Benutzerhandbuch beschreibt die Regel im Ist-Zustand. Es gilt:

- leer gilt als nicht gesetzt;
- `config show` prüft von den Variablen nur `PGWIRE_RECORDER_CONFIG`, die Hilfe und `version` lesen
  nichts;
- Ausnahmen sind `PGWIRE_RECORDER_PASSWORD` (hat keine Option und bleibt bei `record` und `replay`
  unbeachtet; bei `play` gilt sie als Passwort) und jede Variable mit dem Präfix, deren Name keiner
  Option der Tabelle entspricht (ein Tippfehler ist von der Variable eines anderen Werkzeugs nicht
  zu unterscheiden).

**Anlass:** Entscheidung des Nutzers vom 2026-10-12 (Fehler statt Warnung: der Katalog führt keinen
Meldungscode für Warnungen des Starts, und eine Variable, die der Lauf nicht beachtet, soll nicht
unbemerkt bleiben); die Spezifikation trägt sie seit Commit `9bf8f06`. Bis dahin bleiben diese
Variablen unbeachtet, vier Tests erklären das (§3).

**Ergebnis der Lektüre der Tests** (Planner-Auftrag, 2026-10-12; der Architect hatte sie nicht
gelesen): Vier Tests im Paket `internal/adapters/driving/cli` sagen „unbeachtet“ und folgen dem
neuen Vertrag — `TestParseFailOnUnconsumedRecord` (`cli_test.go`), `TestLeserFremdeUmgebung`
(`leser_test.go`), `TestParsePlayFailOnUnconsumed` und `TestParsePlayFremdeUmgebung`
(`play_test.go`); ihre Zeilen in `docs/user/abdeckung-unit.md` sind 22 und 224 (der zweite Test, mit
`LH-FA-03` und `LH-FA-17`), 172, 176 und 219. Von dem, was sie setzen, bleibt manches unbeachtet und
manches wird `PGR-E2001`: In `TestLeserFremdeUmgebung` bleiben `PGWIRE_RECORDER_BOGUS` und
`PGWIRE_RECORDER_` unbeachtet, `PGWIRE_RECORDER_FORMAT`, `_RECORD_TIMING` (bei `replay` und `play`)
und `_KEEP_TIMING`, `_TIMING_MODE`, `_TIMING_REFERENCE` (bei `record` und `replay`) werden Fehler,
und `PGWIRE_RECORDER_SESSION_ASSIGNMENT` ist bei `record` und `play` ein Fehler, bei `replay` nicht
(die Tabelle nennt es dort; der Stand liest die Option noch nicht, §6). In `TestParsePlayFremdeUmgebung`
bleibt `PGWIRE_RECORDER_COMPARE_RESPONSES` bei `play` unbeachtet (die Tabelle nennt `play`; der Stand
liest sie noch nicht), `_LISTEN` und `_SHUTDOWN_TIMEOUT` werden Fehler. **Nicht betroffen**, gelesen:
`TestPlayPasswort` (Zeilen 180, 299 und 330, die Ausnahme `PGWIRE_RECORDER_PASSWORD` bei `play`);
`TestRunHilfe` und `TestRunVersion` in `internal/bootstrap/bootstrap_test.go` (Hilfe und `version`
lesen nichts, die Variable gehört dort zum eigenen Kommando); die Integrationstests setzen
Variablen nur für das Kommando, das sie liest (`FORMAT`, `RECORD_TIMING`, `FORCE` bei `record`;
`FAIL_ON_UNCONSUMED` bei `replay`; die Variablen der Laufsteuerung, der Zeitangaben und des
Passworts bei `play`), `config show` mit `RECORD_TIMING` oder `TIMING_MODE` bleibt grün. Der
Implementer wiederholt die Suche (`grep -rn "PGWIRE_RECORDER_" --include=*_test.go`) und nennt das
Ergebnis in §7.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Der Leser der Optionen, die der Stand noch nicht liest (`--tls-cert`, `--tls-key`,
  `--allow-plaintext`, `--session-assignment`, `--record-empty-sessions`, `--compare-responses`) —
  die Slices, die die Option liefern (`slice-v1-abschluss-tls-client`, `slice-v1-abschluss-sessions`,
  `slice-v1-abschluss-antwortvergleich`; ihr §1 nennt die Option je), geben ihr Leser und Schlüssel.
  Dieser Slice führt sie in der Zuordnung *Option → Kommandos* mit (die Spalte *Kommando* der Tabelle
  gilt „für jede Option der Tabelle“), liest sie aber nicht: bei den eigenen Kommandos bleibt die
  Variable unbeachtet, bis der Slice den Leser liefert. Dafür braucht keiner der drei Slices einen
  Eintrag; die Prüfung der Zuordnung gegen den Leser steht in §6 *Mutationen*.
- Das Passwort (`PGWIRE_RECORDER_PASSWORD`) bei `record` und `replay` — Entscheidung des Nutzers vom
  2026-10-12, bestätigt: bleibt unbeachtet, weil es keine Option hat.
- Eine Warnung statt des Fehlers — Entscheidung des Nutzers vom 2026-10-12 (Fehler); der Katalog hat
  keinen Meldungscode dafür, und einen neuen Code anzulegen wäre ein Eingriff in `SPEC-034`, den dieser
  Slice nicht vornimmt.
- Änderung der Hilfe (`internal/adapters/driving/cli/cli.go`, Texte der Kommandos) — Bestand bleibt
  bewusst stehen: `LH-FA-01.a` verlangt keinen Satz dazu, und `TestLeserHilfe` hält fest, dass die Hilfe
  keine Ausnahme nennt; vom Architect am 2026-10-12 bestätigt (§6).
- Änderung der Spezifikation, einschließlich der Zeile `PGR-E2001` in `SPEC-034` — Schicht-Abgrenzung:
  Die Spezifikation führt, der Slice folgt; die Katalogzeile nennt die neue Ursache seit der Prüfung
  des Architect vom 2026-10-12, der Slice ändert sie nicht (§6).
- `README.md` — sagt zu Umgebungsvariablen nichts und bleibt unverändert (Suche am 2026-10-12: keine
  Stelle, die „unbeachtet“ für Umgebungsvariablen anderer Kommandos behauptet).
- Eine Prüfung der Variablen mit Kleinschreibung, mit anderem Präfix oder auf anderen Plattformen —
  Der Name einer Option der Tabelle ist die einzige Form; was ihr nicht entspricht, bleibt unbeachtet
  (§6; `LH-FA-17.a` seit dem 2026-10-12: der Name gilt genau).
## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1:** Der CLI-Adapter prüft die Umgebungsvariablen der Optionen, die das
      Kommando nicht kennt, nach den Zusagen aus §1 und der Entscheidung der Randformen in
      `LH-FA-17.a` (§6): gesetzt, nicht leer → `PGR-E2001` mit dem Namen der Variable und ohne Wert;
      leer gilt als nicht gesetzt; `config show`, Hilfe und `version` nach §1; Ausnahmen nach §1. Die
      vier Tests aus §1 folgen dem Vertrag, jede Zusage hat einen Test, der unter ihrer Mutation rot
      wird (§6 *Mutationen*); die Abdeckungstabellen (`make abdeckung`) folgen den Deklarationen
      der Tests, Zeilen 22, 172, 176, 219 und 224 von `docs/user/abdeckung-unit.md` sagen nicht mehr
      „unbeachtet“ für Variablen, die jetzt Fehler sind.
- [ ] **Liefer-Punkt 2:** `docs/user/benutzerhandbuch.md` beschreibt die Regel im Ist-Zustand
      (`AGENTS.md` §3.11): im Abschnitt *Einstellungen* (gesetzte, nicht leere Umgebungsvariable der
      Option eines anderen Kommandos, Name statt Wert in der Meldung, die zwei Ausnahmen, `config show`)
      und in der Zeile `PGR-E2001` der Meldungstabelle; ohne Chronik, ohne Verweis auf Spezifikation,
      Lastenheft, ADRs, Slices, Wellen oder Reviews. Jede Aussage dort hat einen Test aus Punkt 1;
      ein Beispiel als Datei gegen das gebaute Binary ist nicht nötig, solange das Handbuch keines
      nennt.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Die Mutationen und ihre roten Tests stehen als Tabelle in §7 (Zusage · Mutation · roter Test),
      vom Implementer selbst gefahren (`AGENTS.md` §3.10), dazu das Ergebnis der Suche nach Tests,
      die Variablen setzen (§1).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben, falls eine Beobachtung anfiel; kein
      Zähler wird gesetzt. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7
      notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).
## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/cli/cli.go` | update | Zuordnung *Option → Kommandos* für alle Optionen der Optionstabelle, auch die der Stand noch nicht liest; die Prüfung der Variablen anderer Kommandos in `lies` an der Stelle, die `LH-FA-17.a` (§6) festlegt; Namen und Zuschnitt wählt der Implementer |
| `internal/adapters/driving/cli/cli_test.go`, `leser_test.go`, `play_test.go` | update | `TestParseFailOnUnconsumedRecord`, `TestLeserFremdeUmgebung`, `TestParsePlayFailOnUnconsumed`, `TestParsePlayFremdeUmgebung` auf den neuen Vertrag; neue Tests je Zusage aus §6 (Happy/Boundary/Negative nach `LH-FA-17.a`) |
| `docs/user/abdeckung-unit.md`, `abdeckung-gesamt.md`, `abdeckung-vollstaendig.md` | erzeugt (`make abdeckung`) | folgen den Deklarationen; in `-unit.md` ändern sich die Zeilen 22, 172, 176, 219, 224, die Zeilen von `TestPlayPasswort` (180, 299, 330) bleiben; `TestLeserUmgebungConfigShowUndHilfe` deklariert zusätzlich `LH-FA-01/Boundary` (die Hilfe liest nichts, `LH-FA-01.a`), womit `LH-FA-01` in `-gesamt.md` von „teilweise“ auf „vollständig“ geht und `-vollstaendig.md` die Zeile `LH-FA-01` bekommt |
| `docs/user/benutzerhandbuch.md` | update | Abschnitt *Einstellungen* und Zeile `PGR-E2001` der Meldungstabelle, im Ist-Zustand |

- Kein Bootstrap: `internal/bootstrap` reicht den Befehl durch und bekommt keine Logik (`AGENTS.md`
  §3.13); kein Hexagon-Kern, kein Eintrag im Katalog der Meldungscodes, keine neue Option.
- Der Implementer wiederholt vor dem Commit die Suche aus §1 und zieht jeden Test nach, der eine
  Variable eines anderen Kommandos setzt und im selben Test das nächste Kommando liest (Variablen,
  die ein `t.Setenv` in einer Schleife über `leserKommandos()` stehen lässt).
- Der Implementer prüft, ob eine weitere lebende Stelle „unbeachtet“ für Umgebungsvariablen anderer
  Kommandos behauptet (`grep -rn "unbeacht" spec/ docs/user/ README.md internal/ harness/`) und zieht
  sie im selben Commit nach, soweit sie nicht die Ausnahmen aus §1 meint.
## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` ist leer (WIP-Limit 1; `slice-harness-lint-gegenprobe-last`
und `slice-v1-abschluss-zeitangaben-doku` liegen in `done/`); dieser Slice steht direkt hinter den Slices von
welle-v1-abschluss bis `slice-v1-abschluss-zeitangaben-doku` und vor `slice-v1-abschluss-antwortvergleich`
(Drift-Log der Roadmap vom 2026-10-12). Technisch hängt er an nichts, auch nicht an der Gate-Kette: seine
Dateien (§3) und Zusagen (§6) berühren weder `SPEC-052` noch `tools/harness`, `harness/mk` oder die Hooks,
und seine DoD verlangt `make gates`, die volle Kette, die heute schon läuft. Die Reihenfolge — dieser Slice
zuerst, die Gate-Kette (`slice-harness-gate-kette-schnell`, `slice-harness-gate-kette-stop-hook-liste`,
`slice-harness-gate-kette-stop-hook-freigabe`) nach dem Release — ist Entscheidung des Nutzers vom
2026-10-12 (Ziel: ein Release; vorher stand die Gate-Kette davor, ebenfalls Entscheidung des Nutzers vom
2026-10-12). Die Randformen aus §6 sind vom Architect am 2026-10-12 vor dem Code in `LH-FA-17.a` entschieden
(`BEO-REPO/randform-wellenlos-ohne-architect-vor-code`, `AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Prüfung des Architect schätzt den Diff
  über eine Review-Sitzung (Schätzung des Planners: 200 bis 400 Zeilen mit Tests und Handbuch, nicht
  gemessen). Schnitt dann nach Lieferwert: der Handbuch-Teil in einen Doku-Folge-Slice direkt dahinter
  (`AGENTS.md` §3.11; Kennung vergibt der Planner), die Prüfung bleibt hier.
- `in-progress` → `open` (blockiert — Carveout?): Der Nutzer nimmt die Entscheidung vom 2026-10-12
  zurück (Warnung statt Fehler, oder `PGWIRE_RECORDER_PASSWORD` auch bei `record` und `replay` ein
  Fehler); dann ändert der Architect `LH-FA-17.a`, und der Plan folgt (`AGENTS.md` §3.9).
## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün (vom Auftraggeber gefahren), die Tabelle der Mutationen und die
Suche in §7, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12): Der Slice liefert einen neuen Vertrag. Ort aller Zeilen ist
`LH-FA-17.a` (*Einstellungen*, Spezifikation `spezifikation.md` §1). Eine nicht genannte oder offene
Randform entscheidet der Implementer nicht, er gibt sie zurück. Stand heute:

| Randform | Stand | Ort / Frage |
|---|---|---|
| gesetzt und nicht leer, Option eines anderen Kommandos → `PGR-E2001`; gleich welcher Wert | entschieden | `LH-FA-17.a` (Entscheidung des Nutzers 2026-10-12) |
| gültiger Wert der fremden Variable (`true`, `5s`) und ungültiger (`GEHEIM`) | entschieden: beides Fehler | `LH-FA-17.a` („gleich welchen Wert sie hat“) |
| leere Variable | entschieden: gilt als nicht gesetzt | `LH-FA-17.a` |
| Meldung: nennt den Namen der Variable, nie ihren Wert | entschieden; der Wortlaut ist frei, solange er den Namen nennt und den Wert nicht (auch nicht Teile davon) | `LH-FA-17.a` |
| `config show` | entschieden: prüft von den Variablen nur `PGWIRE_RECORDER_CONFIG`; die Liste der aktiven Variablen am Ende bleibt | `LH-FA-17.a` *Anzeige* |
| Hilfe; `version` | entschieden: lesen und prüfen nichts | `LH-FA-17.a`, `LH-FA-01.a` |
| `PGWIRE_RECORDER_PASSWORD` | entschieden: bei `record` und `replay` unbeachtet (Ausnahme, vom Nutzer bestätigt), bei `play` das Passwort | `LH-FA-17.a` |
| Variable mit Präfix, deren Name keiner Option der Tabelle entspricht (`PGWIRE_RECORDER_BOGUS`, `PGWIRE_RECORDER_`) | entschieden: unbeachtet. andere Schreibweisen (`pgwire_recorder_listen`, `PGWIRE_RECORDER_Listen`) ebenso: der Name gilt genau, Groß- und Kleinschreibung unterscheidet (Architect 2026-10-12) | `LH-FA-17.a` |
| Optionen der Tabelle, die der Stand noch nicht liest (sechs, §1), beim eigenen Kommando | entschieden (Architect 2026-10-12): „kennt“ heißt: steht in der Spalte *Kommando*; beim eigenen Kommando bleibt die Variable unbeachtet, bis ihr Slice den Leser liefert (Stand des Codes, die Spezifikation sagt die Prüfung zu), beim anderen Kommando ist sie `PGR-E2001`. Keine Sendung: die drei Slices liefern ihren Leser nach ihrem eigenen §1; der Test der Gleichheit (Mutationen) wird rot, wenn ein Leser eine Option ohne Zuordnung anmeldet | `LH-FA-17.a` |
| Variable gesetzt, und die Kommandozeile setzt dieselbe Option des eigenen Kommandos (eigene Option, ungültige Variable) | entschieden: die Variable wird geprüft, unabhängig von der Priorität. Die Option eines anderen Kommandos auf der Kommandozeile ist schon eine unbekannte Option und geht der Umgebung vor, weil die Kommandozeile zuerst geprüft wird | `LH-FA-17.a` |
| **Rangfolge** der Prüfung gegenüber den anderen `PGR-E2001`-Prüfungen | entschieden (Architect 2026-10-12): Kommandozeile; dann ein Durchlauf über die Variablen aller Optionen der Tabelle in ihrer Reihenfolge, bei einer Option des Kommandos der Wert, bei einer fremden die Variable selbst; dann die Datei; die Reihenfolge der eigenen Variablen bleibt | `LH-FA-17.a` *Fehler* |
| mehrere fremde Variablen zugleich | entschieden: nur die erste in der Reihenfolge der Tabelle, eine Meldung (der Start endet beim ersten Fehler) | `LH-FA-17.a` *Fehler* |
| fremde Variable und eine ungültige eigene Variable zugleich | entschieden: die in der Reihenfolge der Tabelle erste wird gemeldet (z. B. `play`: fremde `PGWIRE_RECORDER_LISTEN` steht vor `_INPUT`, `_SHUTDOWN_TIMEOUT` nach `_INPUT`) | `LH-FA-17.a` *Fehler* |
| fremde Variable und ungültige Datei zugleich | entschieden: die fremde Variable (`PGR-E2001`), die Datei wird nicht gelesen | `LH-FA-17.a` *Fehler* |
| `PGWIRE_RECORDER_CONFIG` und `PGWIRE_RECORDER_LOG_LEVEL` | entschieden: `CONFIG` kennen alle Kommandos (nie fremd; ihre Prüfung ist die der Datei), `LOG_LEVEL` kennen `record`, `replay`, `play`; bei `config show` nur `CONFIG` | `LH-FA-17.a` Tabelle, *Anzeige* |
| Quelle der Zuordnung *Option → Kommandos* | entschieden (Architect 2026-10-12): Weg a. Der Test hält die Erwartung hart kodiert aus der Spalte *Kommando*; ein Test, der `spec/` liest, ist nicht möglich, weil der Build-Kontext der Stufe `test` `spec/` ausschließt (`.dockerignore`, [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)). Driftet die Spalte von der Zuordnung, fängt das das Review gegen die Spezifikation; **akzeptiertes Negativ**, weil nur ein Slice, der eine Option hinzufügt oder einem Kommando zuordnet, die Spalte ändert, und der setzt Zuordnung und Test dann mit (kein Beobachtungs-Eintrag) | `LH-FA-17.a` |
| Hilfetext unverändert | entschieden: unverändert (Architect 2026-10-12); `LH-FA-01.a` verlangt keinen Satz, der Hilfetext sagt „jede Option ist auch über ihre Umgebungsvariable setzbar“ und bleibt wahr | `LH-FA-01.a` |
| Zeile `PGR-E2001` in `SPEC-034` | entschieden: angepasst vom Architect am 2026-10-12 („gesetzte Umgebungsvariable einer Option, die das Kommando nicht kennt“); nicht Teil dieses Slice | `SPEC-034` |

**Mutationen** (`AGENTS.md` §3.10): Der Implementer fährt je Zusage die Mutation selbst und trägt die
Tabelle in §7 ein. Vorgegeben:

| Zusage | Mutation | roter Test |
|---|---|---|
| Gesetzte, nicht leere Variable eines anderen Kommandos ist `PGR-E2001` | Prüfung entfernen | Test der Zusage |
| Jeder Wert, gültig oder nicht, ist ein Fehler | nur ungültige Werte der Wertemenge der Option ablehnen | Fall mit `true`/`5s` |
| Leer gilt als nicht gesetzt | „gesetzt“ mit `os.LookupEnv` statt „nicht leer“ | Fall mit leerer Variable |
| Die Meldung nennt den Namen, nie den Wert | Wert in die Meldung schreiben; Name weglassen | Fall mit `GEHEIM` |
| Jede Option der Tabelle ist erfasst, mit ihren Kommandos | eine Option aus der Zuordnung nehmen; eine Option dem falschen Kommando zuordnen | Test über alle Optionen der Tabelle |
| Jede Option, die ein Leser liest, steht in der Zuordnung mit diesem Kommando | `optionen("play")` um eine Option erweitern, die die Zuordnung nicht kennt | Test der Gleichheit |
| `PGWIRE_RECORDER_PASSWORD` bleibt bei `record` und `replay` unbeachtet | das Passwort wie eine Option behandeln | Fall `record`/`replay` mit Passwort |
| Eine Variable mit Präfix ohne passende Option bleibt unbeachtet | jede Variable mit Präfix als fremd behandeln | Fall `PGWIRE_RECORDER_BOGUS` |
| `config show` prüft nur `PGWIRE_RECORDER_CONFIG` | die Prüfung auch in `config show` | Fall `config show` mit fremder Variable |
| Die Hilfe liest nichts | die Prüfung vor die Hilfe | Fall `record --help` mit fremder Variable |
| Die Rangfolge nach `LH-FA-17.a`: fremde Variable vor der Datei | Prüfung der fremden Variablen nach der Datei | fremde Variable und ungültige Datei zugleich |
| Ein Durchlauf in der Reihenfolge der Tabelle, nur der erste Fehler | fremde Variablen in einem eigenen Durchlauf vor oder nach den eigenen; alle fremden melden | zwei fremde Variablen; fremde nach einer ungültigen eigenen, die in der Tabelle vor ihr steht, und davor |
| Der Name gilt genau | Name ohne Beachtung der Schreibweise vergleichen | Fall `pgwire_recorder_listen` bei `play` unbeachtet |

**Risiken:**

- Tests, die Variablen eines Kommandos setzen und im selben Test ein anderes lesen, werden rot, ohne
  dass der Vertrag verletzt ist — **Ausgang:** entfallen: Der Implementer sucht sie (§1, §3) und
  räumt sie im Test-Diff; `make gates` läuft in der Stufe `test` des `Dockerfile` ohne Variablen des
  Aufrufers (`AGENTS.md` §3.1).
- Die Zuordnung *Option → Kommandos* im Code driftet von der Optionstabelle der Spezifikation (ein
  späterer Slice liefert eine Option und ändert die Spalte, oder der Test hält eine Erwartung, die die
  Spezifikation nicht mehr trägt) — **Ausgang:** entfallen als akzeptiertes Negativ (Architect
  2026-10-12, §6 *Quelle der Zuordnung*): Ein Test kann `spec/` nicht lesen (Build-Kontext), die Spalte
  ändert nur ein Slice, der die Option liefert und Zuordnung und Test dann mit nachzieht, und das Review
  liest gegen die Spezifikation; kein Beobachtungs-Eintrag.
- Der Slice ändert Verhalten, das ein Nutzer mit exportierten Variablen anderer Kommandos bisher
  unbemerkt hatte (CI-Umgebung mit allen `PGWIRE_RECORDER_*` für `record` und `play` in einer
  Shell) — **Ausgang:** entfallen: Entscheidung des Nutzers vom 2026-10-12 (Fehler), die
  Spezifikation sagt „wer `record` und `play` in derselben Umgebung fährt, setzt die Variablen der
  Optionen des jeweils anderen Kommandos dort nicht“; das Handbuch nennt es (Liefer-Punkt 2).
- Die Option kommt in einem späteren Slice zu einem Kommando, das die Zuordnung dort noch nicht führt
  (`slice-v1-abschluss-tls-client`, `slice-v1-abschluss-sessions`, `slice-v1-abschluss-antwortvergleich`)
  — **Ausgang:** entfallen: Die Zuordnung führt jede Option der Tabelle schon, und der Test der
  Gleichheit (Mutationen oben) wird rot, wenn ein Leser eine Option ohne Zuordnung anmeldet; die drei
  Slices brauchen keinen Eintrag.
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

### Belege des Implementers

**Läufe** (2026-10-12): `make test` und `make lint` grün; `make gates` grün am Stand `f1203f8` (Code, Tests,
Abdeckungstabellen, Handbuch; 5 min 41 s), erneut nach diesem Eintrag. Ein erster Lauf am Stand `b91131c` war rot
in `make lint` (`gochecknoglobals` für `tabelleSpec` und `nichtGelesen` in `fremd_test.go`); beide sind Funktionen.
`make abdeckung` hat `docs/user/abdeckung-unit.md`, `-gesamt.md` und `-vollstaendig.md` geschrieben: Die vier
Deklarationen sagen nicht mehr „unbeachtet“ für Variablen, die Fehler sind; die Zeilen von `TestPlayPasswort`
(`PGWIRE_RECORDER_PASSWORD` bei `play`) sind unverändert.

**Suche nach Tests, die Variablen setzen** (`grep -rln "PGWIRE_RECORDER_" --include=*_test.go .`): 29 Dateien in
`internal/adapters/driving/cli`, `internal/bootstrap` und `test/integration`. Umgebaut sind die vier Tests aus §1;
kein weiterer Test wurde bei `make test` oder `make test-integration` rot, denn alle anderen setzen Variablen nur für das Kommando, das sie liest, oder leer. Die Tests, die
`leserKommandos()` schleifen, räumen mit `leere` auf. `grep -rn "unbeacht" spec/ docs/user/ README.md internal/ harness/`
fand keine weitere lebende Stelle, die „unbeachtet“ für die Variablen anderer Kommandos behauptet; die übrigen
Treffer meinen `PGWIRE_RECORDER_PASSWORD`, Platzhalter nicht benutzter Verbindungen, PEM-Text und die Ausnahmen aus §1.

**Mutationen** (`AGENTS.md` §3.10). Weg: Bind-Mount des Arbeitsbaums in das Image `pgwire-recorder:test`
(`go test ./internal/adapters/driving/cli/`), der Mutant wird je Lauf in `cli.go` geschrieben und danach
zurückgesetzt (`diff` gegen die Vorlage leer); der Mutant kommt so im Lauf an, kein Build-Kontext. Der Mutant
baute in jedem Fall; „rot“ nennt die Tests der Zusage.

| Zusage | Mutation | roter Test |
|---|---|---|
| gesetzte, nicht leere Variable eines anderen Kommandos ist `PGR-E2001` | Prüfung entfernt (`if false`) | `TestLeserUmgebungAndererKommandos`, `TestParseFailOnUnconsumedRecord`, `TestParsePlayFailOnUnconsumed`, `TestParsePlayFremdeUmgebung` |
| jeder Wert, gültig oder nicht | nur Werte ablehnen, die die Wertemenge der Option ablehnt | `TestLeserUmgebungAndererKommandos` (Fälle `true`, `5s`), `TestParsePlayFremdeUmgebung` |
| leer gilt als nicht gesetzt (fremde Variable) | `LookupEnv` nur im fremden Zweig (`gesetzt` statt nicht leer) | `TestLeserUmgebungAndererKommandos` (Fall mit leerer Variable) |
| leer gilt als nicht gesetzt (allgemein) | `LookupEnv` für alle Variablen | 68 Tests, darunter `TestLeserUmgebungAndererKommandos` (Aufbau bricht; belegt die Zusage nur zusammen mit der Zeile darüber) |
| Meldung nennt nie den Wert | Wert in die Meldung (`%s=%s`) | `TestLeserUmgebungAndererKommandos` (Fall `GEHEIM`), `TestParseFailOnUnconsumedRecord`, `TestParsePlayFailOnUnconsumed`, `TestParsePlayFremdeUmgebung` |
| Meldung nennt den Namen | Name weggelassen | dieselben Tests |
| jede Option der Tabelle ist erfasst | `user` aus der Zuordnung genommen | `TestLeserUmgebungAndererKommandos` (`record`, `replay`) |
| jede Option mit ihren Kommandos | `output` auch `replay` zugeordnet | `TestLeserUmgebungAndererKommandos` |
| jede Option, die ein Leser liest, steht in der Zuordnung | `optionen("play")` um `extra-option` erweitert | `TestLeserZuordnungGleichTabelle` |
| sechs nicht gelesene Optionen: beim eigenen Kommando unbeachtet | nicht gelesene Option beim eigenen Kommando als fremd | `TestLeserUmgebungNichtGelesenerOptionen`, `TestParsePlayFremdeUmgebung` |
| `PGWIRE_RECORDER_PASSWORD` bleibt bei `record`/`replay` unbeachtet | Passwort wie Option (`password` für `play` in der Zuordnung) | `TestLeserFremdeUmgebung` |
| Variable mit Präfix ohne Option bleibt unbeachtet | jede Variable mit Präfix (außer Passwort, Config) fremd | `TestLeserFremdeUmgebung` |
| der Name gilt genau | Schreibweise ignoriert (`EqualFold` über `os.Environ`) | `TestLeserFremdeUmgebung` (`pgwire_recorder_listen`) |
| `config show` prüft nur `PGWIRE_RECORDER_CONFIG` | Prüfung auch in `parseConfig` | `TestLeserUmgebungConfigShowUndHilfe` |
| die Hilfe liest nichts | Prüfung vor der Hilfe in `Parse` | `TestLeserUmgebungConfigShowUndHilfe`, `TestParseHilfeVorPruefung` |
| Rangfolge: Variable vor Datei | Prüfung nach `ladeGewaehlte` | `TestLeserUmgebungReihenfolge` (Fall fremde Variable und ungültige Datei) |
| Durchlauf in Tabellenreihenfolge: fremde vor eigenen | zwei Durchläufe, fremde zuerst | `TestLeserUmgebungReihenfolge` |
| Durchlauf in Tabellenreihenfolge: eigene vor fremden | zwei Durchläufe, eigene zuerst | `TestLeserUmgebungReihenfolge` |
| nur der erste Fehler (eine Meldung) | alle fremden Variablen in einer Meldung | `TestLeserUmgebungReihenfolge` |

Keine grüne Mutation. Nicht als Mutation gefahren: die Zusage „Kommandozeile geht der Umgebung vor“ und die
Priorität der Option des eigenen Kommandos (`TestLeserUmgebungNebenKommandozeile`); sie ist Bestand
und hier nur durch einen Fall ergänzt.

**Probe-Tabelle des Handbuchs** (gebautes Binary aus `make build`, Skript `probe.sh` im Scratchpad; Aufrufe mit
`record --listen=127.0.0.1:0 --upstream=127.0.0.1:1 --output=x.yaml`, `play --upstream=127.0.0.1:1 --input=x.yaml`,
`replay --listen=127.0.0.1:0 --input=x.yaml`, Umgebung nur mit den genannten Variablen):

| Satz im Handbuch | Aufruf | Ergebnis |
|---|---|---|
| Variable einer unbekannten Option ist `PGR-E2001`, gleich welcher Wert (`true`) | `record` mit `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=true` | Exit 2, `Konfiguration [PGR-E2001]: Umgebungsvariable PGWIRE_RECORDER_FAIL_ON_UNCONSUMED gehört zu einer Option, die record nicht kennt` |
| dasselbe bei `play` (`LISTEN`) | `play` mit `PGWIRE_RECORDER_LISTEN=127.0.0.1:1` | Exit 2, `PGR-E2001`, nennt `PGWIRE_RECORDER_LISTEN` |
| Wert `5s` | `play` mit `PGWIRE_RECORDER_SHUTDOWN_TIMEOUT=5s` | Exit 2, `PGR-E2001`, nennt den Namen |
| Wert, den die Option nicht annimmt, und Name statt Wert | `play` mit `PGWIRE_RECORDER_FORCE=GEHEIM` | Exit 2, `PGR-E2001`, Meldung ohne `GEHEIM` |
| `false` bei `play` | `play` mit `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=false` | Exit 2, `PGR-E2001` |
| leere Variable gilt als nicht gesetzt | `record` ohne Optionen, `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED=` | Exit 2, `Pflichtoption --listen fehlt` |
| `PGWIRE_RECORDER_PASSWORD` bei `record` und `replay` unbeachtet | `record` / `replay` ohne Optionen mit der Variable | Exit 2, `Pflichtoption --listen fehlt` |
| `PGWIRE_RECORDER_BOGUS`, die zu keiner Option gehört, bleibt bei `record`, `replay` und `play` unbeachtet | `record`, `replay`, `play` je ohne Optionen mit der Variable | Exit 2, `Pflichtoption --listen fehlt` (`record`, `replay`), `Pflichtoption --upstream fehlt` (`play`) |
| `pgwire_recorder_listen`, andere Schreibweise des Namens einer Option, bleibt unbeachtet | `play` ohne Optionen, `pgwire_recorder_listen=x` | Exit 2, `Pflichtoption --upstream fehlt` |
| `config show` prüft nur `PGWIRE_RECORDER_CONFIG` | `config show` mit `PGWIRE_RECORDER_LISTEN=x`, `PGWIRE_RECORDER_FORCE=GEHEIM`; dann `PGWIRE_RECORDER_CONFIG=/gibt/es/nicht.yaml` | Exit 0, listet beide Namen; Exit 2, `PGR-E2004` |
| (Satz der Hilfe, unverändert) Hilfe liest nichts | `play --help` mit `PGWIRE_RECORDER_LISTEN=x` | Exit 0, Hilfetext |
| Zeile `PGR-E2001` der Meldungstabelle | wie die ersten Zeilen | Meldung nennt die Variable, nicht den Wert |

Der Satz „Wer `record` und `play` in derselben Umgebung ausführt, setzt dort die Variablen der Optionen der jeweils
anderen Betriebsart nicht“ folgt aus den ersten beiden Zeilen. Das Handbuch nennt weder Spezifikation noch Slice;
die Tabelle der Optionen (ohne die sechs noch nicht gelesenen) bleibt unverändert, `README.md` auch.

**Nacharbeit nach dem Review** (2026-10-12, Stand `e1a33d1`; Befunde F-1 und F-2 des Reports):

- F-1 (Handbuch, enge Fassung): Der Satz „eine Variable mit dem Präfix, die zu keiner Option gehört, bleibt bei jeder
  Betriebsart unbeachtet“ las sich für die sechs noch nicht gelesenen Optionen mit. Vorher, gebautes Binary: `play` mit
  `PGWIRE_RECORDER_TLS_CERT=x` endet mit Exit 2 und `Konfiguration [PGR-E2001]: Umgebungsvariable PGWIRE_RECORDER_TLS_CERT
  gehört zu einer Option, die play nicht kennt`, der Satz sagte dort Falsches. Gewählt ist die enge Fassung, weil die Probe das
  Handbuch um Optionen erweitern müsste, die das Binary noch nicht liest. Der Satz nennt jetzt nur die geprüften Fälle:
  `PGWIRE_RECORDER_BOGUS` (gehört zu keiner Option) bei `record`, `replay` und `play`, `pgwire_recorder_listen` (andere
  Schreibweise) und `PGWIRE_RECORDER_PASSWORD` bei `record` und `replay`. Keine allgemeine Aussage über Variablen mit dem
  Präfix mehr; die Probe-Zeilen oben tragen die drei Namen (BOGUS jetzt bei allen drei Betriebsarten, aus `probe.sh`).
  Der Test dazu: `TestLeserFremdeUmgebung` (BOGUS und Schreibweise, je Kommando); Mutation „jede Variable mit Präfix fremd“
  in der Tabelle oben, rot dort.
- F-2 (Plan §3): führt jetzt `abdeckung-gesamt.md` und `abdeckung-vollstaendig.md` und die Deklaration `LH-FA-01/Boundary`.
  `LH-FA-01/Boundary` trifft zu: `TestLeserUmgebungConfigShowUndHilfe` ruft `record --help`, `replay -h`, `play --help` und
  `config show --help` mit ungültigen Variablen jeder Option und erwartet die Hilfe ohne Fehler (`LH-FA-01.a`: die Hilfe geht
  jeder Prüfung vor); die Mutation „Prüfung vor der Hilfe“ der Tabelle oben ist dort rot. Kopf, §1 und §6 nennen `LH-FA-01`
  und die Hilfe schon; nichts weiter nachzuziehen.
- Läufe: `make docs-check` (633 Dateien, 0 Befunde), `make abdeckung-check` grün, `make build` (Image unverändert zwischengespeichert,
  Code seit `f1203f8` unberührt), `make gates` grün am Stand `e1a33d1` (4 min 45 s), erneut nach diesem Eintrag.

**Randformen:** keine nicht genannte oder offene Randform gefunden, keine Rückgabe an den Architect. Gelesen wurde
„kennt“ für die sechs noch nicht gelesenen Optionen nach der Zeile in §6 (Architect 2026-10-12). Die Meldung lautet
`Umgebungsvariable <NAME> gehört zu einer Option, die <Kommando> nicht kennt`; der Wortlaut ist nach §6 frei.

**Verbleibende Risiken:** das Handbuch nennt für den Satz zu Variablen ohne Option nur Beispiele
(`PGWIRE_RECORDER_BOGUS`, `pgwire_recorder_listen`), weil die Handbuch-Tabelle die sechs noch nicht gelesenen
Optionen nicht führt und deren Variablen bei `play`/`record` je nach Kommando Fehler sind; die Slices, die diese
Optionen liefern, schreiben ihren Teil. Der Slice schätzte den Diff auf 200 bis 400 Zeilen; gemessen (`git diff --shortstat`)
sind es 372 hinzugefügte und 43 entfernte Zeilen in `internal/` und im Handbuch, ohne die erzeugten Tabellen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/`
gesichtet (Zähler = Dateien unter `evidence/`, gelesen am Stand `f1f4415`, 2026-10-12). Treffer:

- `BEO-REPO/randform-wellenlos-ohne-architect-vor-code` (1×, offen) — der Slice ist wellenlos; die
  Randformen, die §6 als offen führt, entscheidet der Architect vor dem Code (§4 *Start*).
- `BEO-REPO/begriff-der-meldung-weicht-vom-handbuch-ab` (1×, offen) — berührt: Das Handbuch nennt die
  Meldung; Liefer-Punkt 2 hält Handbuch und Wortlaut der Meldung gleich (nur Name, kein Wert).
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — berührt, aber nicht ausgelöst: Der
  Handbuch-Teil liegt in diesem Slice, kein Doku-Folge-Slice; sein Zählweg steht in der Zählung unten.
- `BEO-REPO/plan-folgt-korrektur-nicht` (verkörpert, `AGENTS.md` §3.9),
  `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (verkörpert, `AGENTS.md` §3.10),
  `BEO-REPO/spec-randform-erst-im-review-entschieden` (verkörpert, `AGENTS.md` §3.12) und
  `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) — neuer Vertrag mit
  Mutationen und Randformen in §6; das Handbuch sagt nur zu, was ein Test prüft.
- `BEO-REPO/yaml-meldung-nennt-wert-und-zeile` (2×, offen) — verwandt über „Meldung nennt den Wert nicht“,
  aber nicht der Gegenstand (dort die Konfigurationsdatei); der Fall mit `GEHEIM` in §6 trägt die Zusage.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Zählung nach `AGENTS.md` §3.13:** Liefer-Punkte: 2 (die Prüfung im CLI-Adapter mit Tests und
Abdeckungs-Deklarationen; der Handbuch-Teil). Gate-Läufe, Review, Closure-Notiz und Register zählen
nicht. Schichten: 2 — CLI-Adapter (mit seinen Tests) und Nutzer-Doku (`docs/user/benutzerhandbuch.md`;
die Abdeckungstabellen unter `docs/user/` sind erzeugte Folge der Test-Deklarationen und zählen auch
dann nicht dazu, wenn man sie als Doku-Schicht rechnete, denn es bleibt dieselbe Schicht);
`README.md` bleibt unverändert; Bootstrap ohne Logik (reine Durchreichung), kein Kern, keine andere
Schicht des Hexagons. Der Handbuch-Teil bleibt damit im Slice (`AGENTS.md` §3.11, §3.13: höchstens
zwei Schichten); ein Doku-Folge-Slice ist nicht nötig. Der Slice trägt keine Sendung in einen anderen
Slice ein (die drei Slices mit den noch nicht gelesenen Optionen brauchen keinen Eintrag, §1 und §6),
also gibt es keinen Nehmer nachzuzählen.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
