# Review-Report: slice-harness-integration-wait — 2026-10-08

**Review-Art:** Code-Review (Testgeschirr unter `test/integration`, Belege in §7). Geprüft wird gegen Plan §1, §3 und §6 (Randformen des Architect aus `2ebe058`), `SPEC-038` Absatz *Warten in Tests*, [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md) und die Hard Rules (Modul 10). Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff 2ebe058..6deea56`, ein Commit:

- `6deea56`: `test/integration/record_e2e_test.go` (Typ `recorder` mit `beendet` und `waitErr`, eine Goroutine mit `Wait` in `startProzess`, Cleanup mit `Kill` und 5 s Frist, neue Helfer `warteEnde` und `pruefeExit`, `stop` über beide), `test/integration/extended_e2e_test.go` und `test/integration/extended_replay_e2e_test.go` (je ein SIGTERM-Test über die Helfer), Plan §3 und §7 *Belege des Implementers*.

**Skill:** `.harness/skills/reviewer.md` am Stand `6deea56`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-08

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-harness-integration-wait.md`: ganz gelesen am Stand `6deea56`, dazu der Plan-Diff aus `6deea56`. §7 *Belege des Implementers* habe ich gelesen, den Bericht des Implementers nicht.
- `spec/spezifikation.md` `SPEC-038`, Absatz *Warten in Tests*.
- [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md); kein `LH-*` (der Plan nennt keinen Lastenheft-Bezug).
- `AGENTS.md` §3.3, §3.7, §3.9 bis §3.12.
- Bestand an `2ebe058`: die vier alten Wait-Stellen, alle Aufrufe von `stop`, alle Leser von `stderr.String()` und `ProcessState` unter `test/integration`; `tools/test/run-integration-tests.sh`, Stufe `integration` im `Dockerfile`.
- Register: `BEO-REPO/liste-gruener-mutanten-unvollstaendig`, `BEO-REPO/spec-randform-erst-im-review-entschieden` (Name aus `AGENTS.md` §3.12).
- Vorheriger Report am Modul: `slice-harness-lint` (F-471 bis F-479). Die höchste vergebene Nummer vor diesem Lauf war F-479.

**Ausgeführte Läufe:**

- **I2-Messung nachgefahren**, je Stand eine frische Kopie per `git archive <stand> | tar -x -m` unter dem Scratch-Pfad der Sitzung, Mutation I2 per Skript mit geprüfter Trefferzahl 1 (Block `if l.herunterfahren { … return model.ErrShutdown }` im Zweig *neue Interaktion* von `ClientMessage`), danach `touch`. Je Stand ein eigenes Image-Tag der Stufe `integration`; je Lauf ein eigenes `--internal`-Netz, ein eigener PostgreSQL-Container (Image und Digest aus `harness/mk/integration.mk`) und ein eigenes Volume, nach jedem Lauf entfernt, nach dem Stand Image und Kopie entfernt. Gestartet mit `-test.run '^TestE2ERecordExtendedSigtermBeimPipelining$' -test.timeout 60s`. Nichts lief im Repo; nach dem Lauf sind keine Container, Netze, Volumes oder Images mit den Lauf-Namen übrig.
- **Probe zu F-480** im gepinnten Go-Image (`golang:1.27-alpine`, Digest aus dem `Dockerfile`, `--network none`): `sh -c 'exit 4'` starten, 300 ms warten, `Process.Signal(SIGTERM)`; einmal ohne frühes `Wait`, einmal mit `Wait` in einer Goroutine direkt nach `Start`.
- `make docs-check` vor dem Commit.
- Nicht gefahren: `make gates`, `make lint`, die Mutationen aus §7 (deren Roh-Logs habe ich für die I2-Läufe stichprobenweise gelesen: Hänger *vorher* Lauf 8 steht in `awaitGoroutines` unter `record_e2e_test.go:373`).

**Messung I2** (eigene Läufe):

| Stand | Läufe | rot | hängend | Testdauer je roter Lauf | Lauf gesamt | Zeitlimit |
|---|---|---|---|---|---|---|
| vorher `2ebe058` | 20 | 18 | 2 (Läufe 5 und 12) | 6,11–6,12 s | 7,6–8,8 s; hängend 61,6–61,7 s | 60 s |
| nachher `6deea56` | 30 | 30 | 0 | 6,11 s | 7,6–8,9 s | 60 s |

Beide Hänger *vorher* melden zuerst die Meldung des Tests und hängen danach in `os/exec.(*Cmd).awaitGoroutines` unter `record_e2e_test.go:373` (Cleanup), wie V-88. *Nachher* meldet jeder Lauf die Meldung des Tests mit `stderr` darunter.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-480 | MEDIUM | Mit dem `Wait` direkt nach `Start` erntet die Goroutine einen Prozess, der von selbst endet, sofort; `Process.Signal` liefert danach `os: process already finished`. `stop` und die beiden SIGTERM-Tests enden dann mit `t.Fatalf("SIGTERM: …")`, ohne Exit-Code und ohne `stderr`. An `2ebe058` blieb der Prozess bis zum `Wait` in `stop` ein Zombie, das Signal ging durch, und der Exit-Code wurde verglichen (Probe: alt `Signal: <nil>`, `Exit: 4`; neu `Signal: os: process already finished`, `Exit: 4`). Die Randform *Prozess ist beim Signal schon beendet* nennt §6 nicht; *Test scheitert vorher* deckt nur den Cleanup. *Failure-Szenario:* Eine Regression lässt den Recorder vor `stop` abstürzen (Panic, Exit 2). Der Test meldet nur `SIGTERM: os: process already finished`, der Stack in `stderr` fehlt in der Meldung. Endet der Prozess dagegen von selbst mit dem erwarteten Code, hängt das Ergebnis vom Zeitpunkt ab: rot, wenn die Goroutine vor dem Signal geerntet hat, sonst grün. `slice-v1-abschluss-herunterfahren` schreibt seine Tests mit Signal und Frist über diese Helfer. | `AGENTS.md` §3.12 („Ein Diff, der eine Randform entscheidet, die §6 nicht nennt“); `BEO-REPO/spec-randform-erst-im-review-entschieden` | `test/integration/record_e2e_test.go` · `if err := r.cmd.Process.Signal(syscall.SIGTERM); err != nil {`; · `r.waitErr = cmd.Wait()` | ja (Probe oben; im Geschirr mit einem Prozess, der vor `stop` endet) | Randform vom Code entschieden, die §6 nicht nennt |
| F-481 | LOW | Zwei Punkte zur Einordnung der Mutanten in §7. (a) „Wartet der Cleanup ohne Frist“ steht unter *äquivalent*. Der Mutant ist nicht äquivalent: Endet der Prozess nach `Kill` nicht, hängt er bis zur Zeitgrenze von `go test`. Rot wird er nur zusammen mit der Mutation *`Kill` im Cleanup entfernt*, und die steht in §7 als Prüfung der Frist, nicht als Prüfung dieses Mutanten. (b) Die Cleanup-Mutation fügt ein `t.Fatal` ein; der Test ist damit ohnehin rot. Ob der Cleanup den Test über `t.Errorf` rot macht, unterscheidet der Lauf nicht, und ein Mutant `t.Errorf` → `t.Logf` erzeugt dieselbe Ausgabe. Der Mutant fehlt in der Liste. Der Kommentar an `startProzess` sagt „danach t.Errorf“ zu. | `AGENTS.md` §3.10, §3.11; `BEO-REPO/liste-gruener-mutanten-unvollstaendig` | `docs/plan/planning/in-progress/slice-harness-integration-wait.md` · „(nicht gefahren, äquivalent nach §6)“; `test/integration/record_e2e_test.go` · „wartet höchstens 5 s auf das Ende; danach t.Errorf.“ | nein (Lesen; die Ausgabe von (b) ist mit und ohne Mutant gleich) | Liste grüner Mutanten unvollständig |
| F-482 | INFO | Leser von `stderr` vor dem Ende von `Wait`, außerhalb des Diffs. `startProzess` liest bei Ablauf der 10-s-Frist zum Lauschen `stderr.String()`, während der Prozess läuft und `os/exec` in denselben `strings.Builder` kopiert; diese Zeile ist unverändert. Der Diff selbst beseitigt zwei solche Leser: das alte `stop` las `stderr` direkt nach `Kill`, der alte Pipelining-Test nach Ablauf der Frist ohne `Kill`. Die übrigen Leser in den Tests lesen nach `stop`. Die Stufe `integration` baut ohne `-race`, der Wettlauf fällt nicht auf. Für Architect bzw. Planner. | Maintainability | `test/integration/record_e2e_test.go` · `t.Fatalf("Recorder lauscht nicht auf %s: %v\n%s", listen, err, stderr.String())` | nein (kein Lauf mit `-race`) | Ausgabe eines laufenden Prozesses vor `Wait` gelesen |
| F-483 | INFO | Für den Verifier, zu den Schwerpunkten Nebenläufigkeit und gleiche Prüfungen. **happens-before:** `waitErr` und `cmd.ProcessState` schreibt die Goroutine vor `close(r.beendet)`. Nach dem Speichermodell von Go ist das Schließen vor jedem Empfang synchronisiert, der wegen des Schließens zurückkehrt. `pruefeExit` läuft an allen drei Aufrufstellen erst nach `warteEnde`, also nach dem Empfang. `warteEnde` liest `stderr` nur im Zweig nach `<-r.beendet`. Der Cleanup liest nur den Kanal. Einen Leser von `waitErr` oder `ProcessState` vor dem Schließen gibt es unter `test/integration` nicht. **Ein Wait:** `grep -rn "\.Wait()" test/integration` findet nur `record_e2e_test.go:379`. **Gleiche Prüfungen:** Fristen 15 s, 5 s und 10 s und die Wartemeldungen sind wörtlich gleich, die Exit-Erwartungen (`wantExit`, 0, 0) ebenso. Anders ist die Meldung bei falschem Exit-Code (`Exit-Code %d, erwartet %d (%s, Wait: %v)`) und das `Kill` im Pipelining-Test nach Ablauf der Frist. Beides hat §6 entschieden. Abweichung nur bei der Randform aus F-480. **`SPEC-038`:** Die Fristen stehen als Literal am Aufruf (`15*time.Second`, `5*time.Second`, `10*time.Second`), die Nachfrist und die Frist im Cleanup als Literal 5 s. Zusammen höchstens 20 s. Jede Meldung nennt das ausgebliebene Ende des Prozesses. | `SPEC-038` *Warten in Tests*; Plan §6 *Ergebnis nach dem Ende*, *Frist und Ablauf*, *Mehrere Leser* | `test/integration/record_e2e_test.go` · `close(r.beendet)`; · `func (r *recorder) pruefeExit(` | ja (Lesen, `grep`) | — (Negativbefund) |
| F-484 | INFO | Zur I2-Messung. Zehn Läufe je Stand belegen den Fix allein nur schwach. Mit allen bekannten Läufen *vorher* (an `1e4381b` 2 von 5, an `f28a2df` 1 von 5, `2ebe058` 1 von 10 laut §7, 2 von 20 eigene) liegt die Hängerate bei etwa 10–20 %. Auch ohne Wirkung des Umbaus bleiben zehn Läufe dann mit 11–35 % Wahrscheinlichkeit ohne Hänger. Meine Läufe (Tabelle oben: vorher 2 von 20 hängend, nachher 0 von 30) ergeben mit §7 *nachher* 0 von 40; bei 10 % wäre das eine Wahrscheinlichkeit von etwa 1,5 %. Tragender ist der strukturelle Beleg: Der Hänger braucht zwei `Wait` auf demselben `exec.Cmd`, und davon gibt es nur noch eines (F-483). Die Bedingung der Rückführung in §4 ist nicht eingetreten. | Plan §6 Risiko *Messung unterscheidet nicht*; DoD-Punkt 2 | `docs/plan/planning/in-progress/slice-harness-integration-wait.md` · „Mindestens zehn Läufe je Stand“ | ja (Messung oben) | — (Negativbefund) |

---

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `test/integration/record_e2e_test.go` | geprüft. Befunde F-480, F-481 (b), F-482. Ein `Wait`, Synchronisation über den Kanal korrekt (F-483). |
| `test/integration/extended_e2e_test.go` | geprüft. Frist, Meldung und Exit-Erwartung gleich; das `Kill` bei Ablauf hat §6 entschieden. Signal vor dem Warten wie `stop`, siehe F-480. Der `defer pc.Close` hängt in keinem meiner 50 Läufe. |
| `test/integration/extended_replay_e2e_test.go` | geprüft. Frist, Meldung und Exit-Erwartung gleich. Signal siehe F-480. |
| übrige Dateien unter `test/integration` | geprüft, ohne Befund aus dem Diff. Keine weitere Stelle mit `Wait`, `ProcessState` oder eigenem Warten auf den Prozess. `replay_e2e_test.go` nutzt `CombinedOutput` mit eigenem Kontext, wie §1 sagt. Zu Lesern von `stderr` siehe F-482. |
| Abdeckungs-Deklarationen, Testliste | geprüft, ohne Befund. Im Diff steht keine Zeile `Abdeckung:` und keine neue Funktion `Test…`. |
| `docs/plan/planning/in-progress/slice-harness-integration-wait.md` | geprüft. Befund F-481 (a). §3 folgt dem Code (§3.9). §6 ist in `6deea56` unverändert. |
| Hard Rule 3.3 | geprüft, ohne Befund. Im Diff gibt es keinen Move. |
| Hard Rule 3.7 (Kommentar-Klassen) | geprüft, ohne Befund. Die neuen Kommentare an `recorder`, `startProzess`, `warteEnde` und `pruefeExit` beschreiben den Ist-Zustand. Zur Zusage „danach t.Errorf“ siehe F-481. |
| Hard Rule 3.12 / MEDIUM *Randform im Code-Commit* | geprüft, ohne Befund für das MEDIUM. `6deea56` ändert §6 nicht. Zur nicht genannten Randform siehe F-480. |
| Hard Rule 3.13 / MEDIUM *Adresse nimmt nicht an* | geprüft, ohne Befund. Der Diff nennt keinen Slice neu als Adresse; die Annahmen in §1 stammen aus `2ebe058`. |
| Produkt-Code, Core-Reinheit ([ADR-0001](../plan/adr/0001-hexagonale-architektur.md)) | geprüft, ohne Befund. Der Diff berührt nichts außerhalb von `test/integration` und dem Plan; I2 lief nur in Kopien. |
| Commit-Message `6deea56` | geprüft, ohne Befund. Sie nennt `slice-harness-integration-wait` und [ADR-0026](../plan/adr/0026-build-und-test-im-multistage-dockerfile.md), keine Struktur-ID. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 3 |

**Summary:** 0 HIGH · 1 MEDIUM · 1 LOW · 3 INFO (F-480 frühes `Wait` macht das Signal an einen schon beendeten Prozess zum Fehler, Randform nicht in §6; F-481 Cleanup ohne Frist als äquivalent eingeordnet, `t.Errorf` → `t.Logf` nicht unterschieden; F-482 `stderr` beim Lausch-Timeout vor `Wait` gelesen, außerhalb des Diffs; F-483 happens-before, ein `Wait`, gleiche Fristen und Meldungen, `SPEC-038` bestätigt; F-484 I2 nachgemessen, vorher 2 von 20 hängend, nachher 0 von 30). Wiederkehrende Klassen: `BEO-REPO/spec-randform-erst-im-review-entschieden` (F-480), `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (F-481).

**Finding-Klassen dieses Laufs:** Randform vom Code entschieden, die §6 nicht nennt · Liste grüner Mutanten unvollständig · Ausgabe eines laufenden Prozesses vor `Wait` gelesen

## Verdikt

**Merge-blockierend:** nein, aber F-480 braucht vor der Closure einen Ausgang. Der Umbau hält, was der Slice verspricht: ein `Wait` je Prozess, korrekt synchronisierte Leser, gleiche Fristen und Meldungen, und unter I2 kein Hänger mehr (eigene Messung 0 von 30 gegen 2 von 20 vorher). Offen ist eine Randform, die der Code still entschieden hat. Sie trifft gerade den Folge-Slice, der die Helfer weiter nutzt.

**Übergabe:** F-480 geht nach `AGENTS.md` §3.12 an den Architect: die Randform in §6 entscheiden, danach der Implementer. F-481 geht an den Implementer (Einordnung in §7, Liste der grünen Mutanten für den Planner). F-482 geht an Architect bzw. Planner, F-483 und F-484 an den Verifier. Dieser Report ersetzt keine Verifikation.
