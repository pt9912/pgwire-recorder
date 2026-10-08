# Slice slice-harness-integration-wait: Ein `Wait` je Prozess im Integrations-Testgeschirr

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt
wird er von der nächsten Welle-Closure. Eingeschoben nach Entscheidung des Nutzers vom
2026-10-07: nach `slice-harness-lint` und vor `slice-harness-abdeckung-gate`
(WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** kein Lastenheft-Bezug — der Slice ändert nur das Testgeschirr, keine Zusage
des Produkts und keinen Nachweis; die Abdeckungs-Deklarationen der Tests bleiben
unverändert. Bindung an Entscheidungen:
[ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md)
(Integrationstests in der Stufe `integration` gegen ein gepinntes PostgreSQL-Image),
[ADR-0028](../../adr/0028-abdeckung-je-anforderung-und-pfad.md) (Abdeckungs-Deklaration
je Test, hier nur unverändert gehalten).

**Berührte Spec-Stellen:** `SPEC-038` (Absatz *Warten in Tests*, angewandt; vom Architect vor
dem Code eingetragen)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Im Integrations-Testgeschirr unter `test/integration` ruft für jeden
gestarteten Prozess genau eine Stelle `Wait` auf dem `exec.Cmd` auf; `stop`, die
Goroutinen der Tests und der `t.Cleanup` aus `startProzess` lesen ihr Ergebnis, statt
selbst zu warten. Ein Test, dessen Prozess nach dem Signal nicht endet, wird damit
nach seiner eigenen Frist rot, statt bis zum Zeitlimit von `go test` zu hängen.

Die Regel, nach der das Testgeschirr wartet, steht in `SPEC-038` *Warten in Tests*:
eigene Frist als Literal, höchstens 60 s, rot mit Klartext statt Hängen, ein `Wait` je
Prozess. Der Slice wendet sie in `test/integration` an; die Einzelheiten der Helfer
stehen in §6 *Randformen*.

**Übernimmt:** aus `slice-v1-abschluss-betrieb` das Risiko *Testgeschirr wartet zweimal
auf denselben Prozess* (Verifikation V-88 zu `slice-harness-blackbox-einstieg`), nicht
dessen Gegenstand; dort steht der Ausgang *übernommen von*
`slice-harness-integration-wait`.

**Aus `slice-tests-ueberlebende-mutanten-driving`** (dort §1, Abgrenzung; Review F-468):
die Fristen in `test/integration`, soweit sie das Warten im Integrations-Testgeschirr
betreffen. Nach `SPEC-038` *Warten in Tests* bekommt dieser Slice dafür die Frist im
`t.Cleanup` und die Nachfrist nach `Kill` (§6); die Fristen der drei Tests bleiben.

**Bestand** (gemessen vom Architect am Stand `86db551`, `grep -n "Wait()"` und
`grep -n "startProzess("` unter `test/integration`): `Wait` steht an vier Stellen —
`record_e2e_test.go:373` (`t.Cleanup` in `startProzess`, nach `Kill`, ohne Frist),
`record_e2e_test.go:396` (`stop`, Goroutine, Frist 15 s, danach `Kill` und `Fatalf`),
`extended_e2e_test.go:410` (`TestE2ERecordExtendedSigtermBeimPipelining`, Goroutine,
Frist 5 s, ohne `Kill`) und `extended_replay_e2e_test.go:297`
(`TestE2EReplayExtendedSigtermMittenInFolge`, Goroutine, Frist 10 s, danach `Kill`).
Jeder Prozess entsteht über `startProzess` (`record_e2e_test.go:360`), direkt oder über
`startRecorder` (`:353`); `stop` (`:390`) ist der einzige weitere Helfer, der wartet.
`exec.Command` steht sonst nur in `replay_e2e_test.go:142` (`CombinedOutput` mit
Kontext-Frist 15 s, ein eigenes `Wait` ohne zweiten Leser) und bleibt. Weitere Stellen
mit dem Muster gibt es nicht; die Rückführung *zu groß* aus §4 tritt nicht ein.

**Befund** (`docs/reviews/2026-10-07-verifikation-slice-harness-blackbox-einstieg.md`,
V-88; Review-Report `docs/reviews/2026-10-07-review-slice-harness-blackbox-einstieg.md`,
F-458): Die Goroutine `beendet <- rec.cmd.Wait()` in
`TestE2ERecordExtendedSigtermBeimPipelining` (`extended_e2e_test.go`) und der
`t.Cleanup` aus `startProzess` (`Kill`, dann `cmd.Wait()`, `record_e2e_test.go`) rufen
`Wait` auf demselben `exec.Cmd` nebeneinander. Den Kanal der Kopier-Goroutinen liest
nur einer der beiden, der andere wartet in `awaitGoroutines` für immer. Unter der
Mutation I2 der Verifikation (`ClientMessage` ohne die Sperre `herunterfahren`) hing der
Test an `1e4381b` in 2 von 5 Läufen bis zum Zeitlimit, an `f28a2df` in 1 von 5. Dasselbe
Muster steht in `stop` (`record_e2e_test.go`, `done <- r.cmd.Wait()`) und in
`TestE2EReplayExtendedSigtermMittenInFolge` (`extended_replay_e2e_test.go`,
`done <- rep.cmd.Wait()`). Im Grünfall wirkt es nicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code — Schicht-Abgrenzung: Der Slice ändert nur Testcode unter
  `test/integration`. Die Mutation I2 wird nur in einer Kopie des Arbeitsbaums gefahren
  und nie committet.
- Neue Tests, geänderte Erwartungen, die Fristen der drei Tests (15 s, 5 s, 10 s) oder
  Abdeckungs-Deklarationen — Bestand bleibt bewusst stehen: Testliste und Deklarationen
  sind die Messlatte, an der sich zeigt, dass nur das Warten umgebaut ist. Neu sind nur
  die Frist im `t.Cleanup` und die Nachfrist nach `Kill` (§6), die `SPEC-038` verlangt.
- `defer pc.Close(context.Background())` in `TestE2ERecordExtendedSigtermBeimPipelining`
  (F-458) — Bestand bleibt bewusst stehen: `Close` wartet nicht auf ein Ereignis des
  Prüflings und fällt nicht unter `SPEC-038`; ob es unter I2 hängt, zeigt die Messung
  (Risiko *Neuer Fund* in §6).
- Warten ohne Frist in Tests anderer Schichten — anderer Vorgang:
  `slice-tests-ueberlebende-mutanten-driving` setzt `SPEC-038` im PGWire-Adapter und im
  Bootstrap um; für den Kern siehe dort §6 *Warten auf einen Kanal*.
- Die Tests mit Signal und Frist je Modus (`--shutdown-timeout`) — übernimmt
  `slice-v1-abschluss-herunterfahren` (aus `slice-v1-abschluss-betrieb` hervorgegangen); er
  schreibt sie über die Helfer, die dieser Slice
  umbaut.
- Mutationen in den Integrationstests als Gate — anderer Vorgang; `slice-harness-mutation`
  schließt sie in seinem §1 aus.
- Ein Zeitlimit für `go test` in `make test-integration` — anderer Vorgang am Gate:
  Ein kürzeres Limit verkürzte das Hängen, beseitigte es nicht.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Ein `Wait` je Prozess: Unter `test/integration` ruft für jeden gestarteten
      Prozess genau eine Stelle `Wait` auf; `stop`, die Goroutinen in
      `TestE2ERecordExtendedSigtermBeimPipelining` und
      `TestE2EReplayExtendedSigtermMittenInFolge` und der `t.Cleanup` aus
      `startProzess` lesen ihr Ergebnis. Beleg in §7: die Fundstellen von `.Wait()`
      unter `test/integration` vor und nach dem Umbau (Suche mit Pfad und Zeile).
- [ ] Rot statt hängend: Unter der Mutation I2 der Verifikation (`ClientMessage` ohne
      die Sperre `herunterfahren`, nur in einer Kopie des Arbeitsbaums) wird
      `TestE2ERecordExtendedSigtermBeimPipelining` in jedem von mindestens zehn Läufen
      rot mit der Meldung des Tests („Recorder endet nicht binnen 5 s nach SIGTERM“) und
      endet weit vor einem Zeitlimit von `go test`, das ein Hängen sichtbar macht; am
      Stand vor dem Umbau zeigt dieselbe Messung mindestens einen hängenden Lauf. Beleg
      in §7: je Stand Quellstand, Zahl der Läufe, rote und hängende Läufe, Dauer je Lauf,
      gesetztes Zeitlimit. Dazu die Mutation des Cleanups aus §6 *Mutationen*: ein Lauf,
      rot mit der Meldung des Cleanups nach etwa 5 s; Beleg in §7 mit Dauer und Meldung.
- [ ] Testliste und Deklarationen unverändert: Die Liste der Tests unter
      `test/integration` (Namen, Zahl) ist vor und nach dem Umbau gleich, keine Zeile
      `Abdeckung:` ist geändert, und `make abdeckung-check` ist grün ohne neu
      geschriebene Tabellen. Beleg in §7: beide Listen oder ihr Vergleich und der Diff
      ohne Deklarationszeile.
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
| `test/integration/record_e2e_test.go` | refactor | `startProzess` startet nach `Start` genau eine Goroutine mit `Wait`, die das Ergebnis ablegt und danach einen Kanal schließt; `recorder` trägt Kanal und Ergebnis. Der Helfer `warteEnde` wartet mit Frist auf das Ende (§6 *Frist und Ablauf*), `pruefeExit` vergleicht danach den Exit-Code und nennt bei Abweichung `ProcessState.String()` und den Fehler von `Wait` (§6 *Ergebnis nach dem Ende*). `stop` geht über beide; der `t.Cleanup` ruft bei offenem Kanal `Kill` und wartet mit Frist auf den Kanal, statt `Wait` aufzurufen |
| `test/integration/extended_e2e_test.go` | refactor | `TestE2ERecordExtendedSigtermBeimPipelining` wartet über `warteEnde` und `pruefeExit` statt `beendet <- rec.cmd.Wait()`; Erwartungen und Fristen gleich |
| `test/integration/extended_replay_e2e_test.go` | refactor | `TestE2EReplayExtendedSigtermMittenInFolge` ebenso statt `done <- rep.cmd.Wait()` |

- Wer den Exit-Code liest (`stop`, die beiden Tests), liest `ProcessState` erst nach dem
  Ende, das der Kanal meldet; danach schreibt niemand mehr daran.
- Die Suche nach `Wait()` unter `test/integration` hat der Architect vor dem Code
  gefahren; die vier Fundstellen stehen in §1 *Bestand*, weitere gibt es nicht. Der
  Implementer wiederholt sie für den Beleg in §7.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-lint` liegt in `done/` (WIP-Limit 1);
dieser Slice ist Schritt 1 der Reihenfolge. Reihenfolge nach Entscheidung des Nutzers vom 2026-10-08 (Wellen vor Harness, M3 vor M4):
die Slices von welle-v1-abschluss, darin `slice-harness-integration-wait` direkt vor
`slice-v1-abschluss-herunterfahren` und `slice-harness-meldungskatalog-gate` direkt vor
`slice-v1-abschluss-container`, dann `slice-harness-abdeckung-gate` und
`slice-harness-coverage`, dann die Slices von welle-erster-release, danach
`slice-harness-commit-struktur-id`, `slice-tests-ueberlebende-mutanten`,
`slice-tests-ueberlebende-mutanten-driving`, `slice-harness-mutation`; die Reihenfolge
der Wellen-Slices steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md). Technisch hängt dieser Slice nur an
`slice-harness-lint`: Ab ihm ist `make lint` Gate auch für `test/integration`. Er steht
direkt vor `slice-v1-abschluss-herunterfahren`, der die Tests mit Signal und Frist über dieselben
Helfer schreibt; deren Mutationen (Prozess endet nach dem Signal nicht) sind genau der
Fall, der heute bis zum Zeitlimit hängt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Suche aus §3 findet das
  Muster in so vielen weiteren Tests, dass der Umbau nicht in einer Review-Sitzung
  prüfbar ist; dann trägt dieser Slice den Helfer und die drei genannten Stellen, die
  übrigen ein eigener Slice vor `slice-v1-abschluss-herunterfahren`.
- `in-progress` → `open` (blockiert — Carveout?): Am Stand vor dem Umbau zeigt die
  Messung aus DoD-Punkt 2 in keinem Lauf ein Hängen, und auch mehr Läufe oder eine
  andere Mutation, unter der der Prozess nach dem Signal nicht endet, führen nicht dazu;
  dann belegt die Messung den Fix nicht, und der Architect entscheidet einen anderen
  Nachweis.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün mit `test-integration` und `abdeckung-check`, die
Messung aus DoD-Punkt 2 in §7 mit null hängenden Läufen nach dem Umbau, Closure-Notiz
mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** — vom Architect vor dem Code entschieden. Der Slice liefert keinen
Vertrag nach `AGENTS.md` §3.10 (keine Zusage des Produkts oder eines Gates); die Helfer
sind Testgeschirr, das `slice-v1-abschluss-herunterfahren` weiter nutzt. Die Regel, die
für alle Tests gilt, steht darum im Technik-Stratum (`SPEC-038` *Warten in Tests*); die
Einzelheiten der Helfer sind Plan-Entscheidungen hier, weil sie nur dieses Geschirr
betreffen. Entscheidet der Umbau etwas, das hier nicht steht, gibt der Implementer es dem
Architect zurück.

- **Ergebnis nach dem Ende** — **entschieden:** Die eine Goroutine aus `startProzess`
  legt den Fehler von `Wait` im `recorder` ab und schließt danach den Kanal; der Zustand
  ist `cmd.ProcessState`. Leser lesen beides erst nach dem Schließen. Der Exit-Code ist
  `ProcessState.ExitCode()`; ein Ende durch Signal ergibt dort `-1`, das Signal steht in
  `ProcessState.String()` (etwa `signal: killed`). Der Helfer deutet nichts: Wer einen
  Exit-Code erwartet, vergleicht ihn und nennt bei Abweichung `ProcessState.String()`.
  Der Fehler von `Wait` ist bei einem Exit-Code ungleich 0 ein `*exec.ExitError` und
  keine eigene Zusicherung (Erwartungen unverändert, §1); er steht nur in der Meldung.
- **Frist und Ablauf** — **entschieden:** Der Helfer zum Warten nimmt die Frist und den
  Text der Meldung vom Aufrufer; die drei Tests und `stop` behalten ihre Fristen und
  Meldungen (15 s, 5 s, 10 s). Läuft die Frist ab, ruft der Helfer `Kill`, wartet eine
  Nachfrist von 5 s auf das Schließen des Kanals und endet mit `t.Fatalf`: der Text des
  Aufrufers, dazu `stderr`, das erst nach dem Ende gelesen wird. Schließt der Kanal auch in
  der Nachfrist nicht, nennt die Meldung das statt `stderr` („endet auch nach Kill nicht
  binnen 5 s“). Rot also spätestens nach Frist plus 5 s, nie an der Zeitgrenze von
  `go test`.
- **Mehrere Leser** — **entschieden:** Der geschlossene Kanal ist das Signal für jeden
  Leser, beliebig oft; Ergebnis und `ProcessState` sind vor dem Schließen geschrieben und
  danach unverändert. Kein `sync.Once`, kein zweites `Wait`, kein Wert, den nur ein Leser
  bekommt. Ein Test, der prüfen will, dass der Prozess **noch nicht** endet (etwa
  `slice-v1-abschluss-herunterfahren` mit `--shutdown-timeout 0`), liest denselben Kanal
  in einem `select` mit eigener Frist.
- **Test scheitert vorher** — **entschieden:** Der `t.Cleanup` wird nach `Start` und nach
  dem Start der Goroutine registriert; scheitert `Start`, gibt es weder Prozess noch
  Goroutine noch Cleanup. Ist der Kanal beim Cleanup offen (Test endete vor `stop`,
  `Fatalf` beim Lauschen, Ablauf im Test), ruft er `Kill` und wartet mit 5 s Frist auf
  das Schließen; bei Ablauf `t.Errorf` mit Klartext, kein Hängen. Ist er geschlossen, tut
  der Cleanup nichts. `Kill` auf einen schon beendeten Prozess ist harmlos, sein Fehler
  wird verworfen. Der Cleanup läuft nach den `defer` des Tests und vor dem Entfernen von
  `t.TempDir` (in `startProzess` später registriert).
- **Mutationen** (`AGENTS.md` §3.10, sinngemäß für das Geschirr) — **entschieden:**
  Zusage *rot statt hängend* · Mutation I2 der Verifikation · rot wird
  `TestE2ERecordExtendedSigtermBeimPipelining` nach höchstens 5 s plus Nachfrist (DoD
  Punkt 2). Zusage *Cleanup beendet und wartet mit Frist* · in einer Kopie des
  Arbeitsbaums `Kill` im Cleanup entfernt und in einem Test direkt nach `startRecorder`
  ein `t.Fatal` eingefügt · derselbe Test endet rot mit der Meldung des Cleanups nach
  etwa 5 s, nicht an der Zeitgrenze (DoD Punkt 2, ein Lauf genügt, die Mutation ist
  deterministisch). Im grünen Lauf endet der Prozess fast immer über `stop` (48 Aufrufstellen
  bei 50 Startstellen am Stand `86db551`), den Kill-Pfad des Cleanups trägt also erst diese
  Mutation.

**Risiken:**

- **Messung unterscheidet nicht** — das Hängen trat an `1e4381b` in 2 von 5 Läufen auf;
  bei wenigen Läufen kann der Stand vor dem Umbau zufällig nicht hängen, und dann
  belegt der Stand nach dem Umbau nichts. Mindestens zehn Läufe je Stand, das
  Zeitlimit so gesetzt, dass ein Hängen als Abbruch erscheint (Rückführung in §4). —
  **Ausgang:** — (bei Closure)
- **Mutant kommt im Build-Kontext nicht an** — die Mutation I2 liegt im Produkt-Code
  und läuft über die Stufe `integration` des `Dockerfile`; BuildKit überträgt eine Datei
  gleicher Größe und mtime nicht neu (`BEO-REPO/mutant-kommt-im-build-kontext-nicht-an`,
  1×). Je Lauf eine frische Kopie unter eigenem Pfad, wie in der Verifikation. —
  **Ausgang:** — (bei Closure)
- **Erste Phase verdeckt die zweite** — der Integrations-Runner endet nach einer roten
  ersten Phase (`BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar`, 1×); die
  Messung läuft darum gezielt für den einen Test (`-run`), nicht über
  `make test-integration`. — **Ausgang:** — (bei Closure)
- **`t.Cleanup` nach fehlgeschlagenem Start oder nach `Fatalf`** — der Cleanup läuft
  auch, wenn `stop` den Prozess schon getötet hat oder der Test vor dem Lauschen
  abbricht; er darf dann weder ein zweites Mal warten noch auf einen Kanal warten, der
  nie geschlossen wird. — **Ausgang:** — (bei Closure)
- **Neuer Fund am Testgeschirr** — F-458 nennt aus dem Bestand die Goroutine in
  `TestE2ERecordExtendedSigtermBeimPipelining`, die `pc` beim `Close` noch benutzt; der
  Umbau ändert daran nichts und darf es nicht verdecken. — **Ausgang:** — (bei Closure)

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

**Belege des Implementers** (Stand vor dem Umbau `2ebe058`; Umbau im Commit dieses Abschnitts):

- **Weg der Läufe:** je Lauf eine frische Kopie unter eigenem Pfad außerhalb des Repos
  (`git archive 2ebe058 | tar -x -m` bzw. die Dateien des Arbeitsbaums über `git ls-files`
  und `tar -x -m`, also mtime beim Entpacken neu, kein `cp -p`), die Mutation per Skript
  in der Kopie, danach `touch` auf die geänderte Datei; je Lauf ein eigenes Image-Tag der
  Stufe `integration`, ein eigenes `--internal`-Netz, ein eigener PostgreSQL-Container
  (Image wie `make test-integration`) und ein eigenes Volume, nach dem Lauf alle vier und
  die Kopie entfernt. Gestartet wird gezielt
  `-test.run '^TestE2ERecordExtendedSigtermBeimPipelining$' -test.timeout 60s`
  (Risiko *Erste Phase verdeckt die zweite*); ein Hängen erscheint als
  `panic: test timed out after 1m0s`. Die Läufe nach dem Umbau liefen am Arbeitsbaum vor
  dem letzten Kommentar-Schliff am Typ `recorder` (nur Kommentarzeilen, Verhalten gleich).
- **DoD 1, Fundstellen von `.Wait()`** (`grep -rn "\.Wait()" test/integration`):
  vorher vier — `record_e2e_test.go:373` (`_ = cmd.Wait()` im `t.Cleanup`),
  `record_e2e_test.go:396` (`stop`), `extended_e2e_test.go:410`,
  `extended_replay_e2e_test.go:297`; nachher eine — `record_e2e_test.go:379`
  (`r.waitErr = cmd.Wait()` in der Goroutine aus `startProzess`). `stop`, die beiden Tests
  und der `t.Cleanup` lesen den Kanal `beendet` (`warteEnde`, `pruefeExit`, Cleanup).
- **DoD 2, Messung unter I2** (Mutation: in `ClientMessage`,
  `internal/hexagon/services/record.go`, der Block `if l.herunterfahren { … return
  model.ErrShutdown }` entfernt; nur in den Kopien):

  | Stand | Läufe | rot | hängend | Testdauer je roter Lauf | Lauf gesamt (mit Containerstart) | Zeitlimit |
  |---|---|---|---|---|---|---|
  | vorher `2ebe058` | 10 | 9 | 1 (Lauf 8) | 6,11–6,12 s | 7–8 s; Lauf 8: 61 s | 60 s |
  | nachher (Umbau) | 10 | 10 | 0 | 6,11–6,15 s | 7,7–9,3 s | 60 s |

  Vorher meldet jeder rote Lauf „Recorder endet nicht binnen 5 s nach SIGTERM, obwohl der
  Client weiter pipelinet“; Lauf 8 meldet dasselbe und hängt danach bis zum Zeitlimit, der
  Stack steht in `os/exec.(*Cmd).awaitGoroutines` unter `Wait` aus `record_e2e_test.go:373`
  (Cleanup) — der Befund V-88. Nachher meldet jeder Lauf dieselbe Meldung mit `stderr`
  des Recorders darunter und endet nach 6,1 s (1 s Pipelinen vor dem Signal, 5 s Frist;
  `Kill` beendet den Prozess sofort, die Nachfrist läuft nicht aus). Der `defer pc.Close`
  hängt in keinem der zehn Läufe (Risiko *Neuer Fund*).
- **Mutationen am Geschirr** (§6 *Mutationen*, je ein Lauf, Weg wie oben):

  | Zusage | Mutation | Ergebnis |
  |---|---|---|
  | rot statt hängend | I2 (oben) | rot, 10 von 10, 6,11–6,15 s, Meldung des Tests |
  | Cleanup beendet und wartet mit Frist | `_ = cmd.Process.Kill()` im Cleanup entfernt, `t.Fatal` direkt nach `startRecorder` in `TestE2ERecordExtendedSigtermBeimPipelining` | rot nach 5,11 s: „record_e2e_test.go:391: Prozess record endet im Cleanup auch nach Kill nicht binnen 5 s“, kein Hängen |
  | Gegenlauf zur Zeile davor: Cleanup mit `Kill` | nur das `t.Fatal` nach `startRecorder` | rot nach 0,10 s mit der Meldung des `t.Fatal`, keine Meldung des Cleanups — der Kill-Pfad beendet den Prozess |
  | Nachfrist nach `Kill` (§6 *Frist und Ablauf*) | I2 und `_ = r.cmd.Process.Kill()` in `warteEnde` entfernt | rot nach 11,11 s: „Recorder endet nicht binnen 5 s nach SIGTERM, obwohl der Client weiter pipelinet; endet auch nach Kill nicht binnen 5 s“; danach beendet der Cleanup den Prozess, kein Hängen |
  | Abweichender Exit-Code nennt `ProcessState.String()` und den Fehler von `Wait` (§6 *Ergebnis nach dem Ende*) | `rec.pruefeExit(t, 0)` → `rec.pruefeExit(t, 1)` | rot nach 1,20 s: „Exit-Code 0, erwartet 1 (exit status 0, Wait: <nil>)“ |
  | Grundlauf | keine | grün, 1,19 s |

  **Grüne Mutanten, eingeordnet** (nicht gefahren, äquivalent nach §6): Entfällt im
  Cleanup die Abfrage des geschlossenen Kanals, ruft er `Kill` auf einen beendeten
  Prozess, dessen Fehler verworfen wird (§6 *Test scheitert vorher*), und der Kanal ist
  geschlossen — gleiches Verhalten, kein Test fängt ihn. Wartet der Cleanup ohne Frist,
  ändert das nur den Fall, dass der Prozess nach `Kill` nicht endet; diesen Fall erzeugt
  nur die Mutation *Cleanup* oben, die die Frist rot zeigt. Weitere grüne Mutanten fielen
  nicht an.
- **DoD 3, Testliste und Deklarationen:** `grep -hn "^func Test" test/integration/*.go`
  ergibt vorher und nachher dieselben 40 Namen (`diff` leer); `git diff -U0
  test/integration | grep -c Abdeckung` ergibt 0, keine Deklarationszeile geändert.
  `make abdeckung-check` läuft in `make gates` (unten), ohne neu geschriebene Tabellen.
- **Gates:** `make test` und `make lint` (0 issues) am Umbau grün; `make gates` grün am
  Commit dieses Abschnitts (Lauf vor der Übergabe).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `cceb17f` gesichtet, dazu der
Nachtrag von F-336 im selben Commit wie dieser Plan (Zähler = Dateien unter
`evidence/`). V-88 selbst steht nicht im Register; seine Adresse ist dieser Slice.
Treffer:

- `BEO-REPO/mutant-kommt-im-build-kontext-nicht-an` (1×) — die Messung unter I2 baut
  einen Mutanten über den Build-Kontext; ein Risiko in §6.
- `BEO-REPO/abnahme-ohne-postgres-nicht-einzeln-brechbar` (1×) — der Runner endet nach
  der ersten roten Phase; ein Risiko in §6.
- `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht` (3×, verkörpert) — darum stehen
  die Belege der DoD in §7, nicht im Bericht.
- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (16×, verkörpert in `AGENTS.md`
  §3.11) — ein neuer Kommentar am Helfer sagt nur zu, was die Messung zeigt.
- `BEO-REPO/plan-folgt-korrektur-nicht` (15×, verkörpert in §3.9) — findet die Suche
  aus §3 weitere Stellen, folgen §1 und §3 im selben Commit.

Keiner der Einträge erreicht mit diesem Slice die Schwelle 3×; keine neue Lücke vor
dem Code.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
