# Review-Report: slice-v1-abschluss-zeitangaben-uhr — 2026-10-11

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3, §6 und §8, gegen die Spezifikation `LH-FA-21.a` (*Uhr*) und die Sicht (Kern liest Zeit nur über den Port), gegen [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0015](../plan/adr/0015-uhr-port.md), [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), [ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md) und gegen `AGENTS.md` §3.

**Gegenstand:** `aa07a01..cf66288` (vier Commits: Roadmap `7312db2`, Port, Fake-Uhr, Systemuhr, Testhilfe und Migration `0ca3161`, a-check und Gegenprobe `685bf54`, Plan und Belege `cf66288`). Geändert: `internal/hexagon/ports/driven/uhr.go`, `internal/testuhr`, `internal/testzeit`, `internal/bootstrap/systemuhr.go` mit Tests und `export_test.go`, drei Testdateien in `internal/hexagon/services`, `.a-check.yml`, `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk`, `harness/README.md`, Roadmap, Plan.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-v1-abschluss-zeitangaben-uhr.md` ganz gelesen am Stand `cf66288`, mit §7 *Belege des Implementers* (Mutationstabelle, grüne Mutanten). Den Bericht des Implementers habe ich nicht gelesen; seine Einschätzungen zu den grünen Mutanten habe ich nachgefahren statt übernommen.
- `spec/spezifikation.md` `LH-FA-21.a` (*Uhr*), `spec/architecture.md` (Uhr-Port, Systemuhr), Lastenheft [`LH-FA-21`](../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen).
- Nehmer `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` (§1, DoD, §3, Start-Bedingung); der Diff ändert sie nicht.
- Vorherige Reports am selben Modul: `slice-v1-abschluss-zeitangaben-pruefen` (Nummern dieses Laufs ab F-640).
- `AGENTS.md` §3.3 bis §3.13.

**Ausgeführte Läufe** (Docker-only, Kopien unter dem Sitzungsverzeichnis, `cp -r`/`tar` ohne mtime, nichts im Repo geändert, Tests im Image `pgwire-recorder:test` ohne Netz, a-check im gepinnten Image):

- **Regel und Zulassungsliste** (vier Kopien von `.a-check.yml`, je Probe-Datei im Domain Model, im PGWire-Adapter, im Bootstrap und `runtime` im Domain Model). Eintrag entfernt: alle vier Proben grün (die Regel ist die Ursache der Ablehnung). `match: regex` entfernt: alle grün. `internal/adapters` aus der Liste: Proben rot, a-check nennt `einspielen.go` und `upstream.go` des Postgres-Adapters. `internal/bootstrap` aus der Liste: Probe im Bootstrap grün, die Liste ist dort wirkungslos. Liste ohne `cmd`, `test`, `internal/testuhr`, `internal/testzeit`, `internal/testpki`: **rot**, 22 `tech-leak` in `internal/testpki`, `internal/testuhr`, `internal/testzeit` und 17 Dateien unter `test/integration` (F-640).
- **Mutanten der Systemuhr und der Fake-Uhr** (eine Ersetzung je Kopie, `go test -run` des Pakets). Rot: Begrenzung der Stufe entfernt (`TestSystemuhrFernerZeitpunktDrehtNicht`, 350 ms CPU-Zeit); halbe Wartezeit ohne Nachlesen (`TestSystemuhrWartenUntergrenze`, Lesung 26 bei Zeitpunkt 51); Lesung gerundet (`TestSystemuhrLesungAbgeschnitten`); Vorab-Prüfung des Kontexts in der Fake-Uhr entfernt (`TestUhrBeendeterKontextVorDemWarten`). **Grün:** `zeitpunkt < jetzt` statt `<=` in `systemuhr.go` (F-642); Prüfung des Kontexts ab der zweiten Schleifenrunde entfernt (äquivalent bis auf das Rennen zwischen Zeitablauf und Abbruch); `return nil` nach dem Timer (äquivalent, wie §7 sagt: der Timer läuft nicht früher ab als seine Dauer, die Lesung ist danach abgeschnitten nicht kleiner als der Zeitpunkt).
- **Flackern.** `go test -c` der drei Pakete, 1 CPU (`--cpus=1`) mit sechs Schleifen, die die CPU belegen: `Systemuhr*` 15 Läufe, `testzeit` 30, `testuhr` 30, alle grün; das Paket `internal/bootstrap` ganz, fünf Läufe mit `--cpus=2`, grün. `TestSystemuhrFernerZeitpunktDrehtNicht`: der ehrliche Lauf verbraucht kaum CPU-Zeit, der Mutant 350 ms, die Grenze liegt bei 150 ms; das Paket enthält keine Tests mit `t.Parallel`, die die CPU-Zeit des Prozesses mitzählen.
- **Migration der drei Testdateien.** Jede Ersetzung gegen den Altwert gelesen: `2 * time.Second` und `2000`, `30 * time.Second` und `30000`, `100`, `50`, `1` Millisekunden; `time.Now().Add`/`After` und `NeueFrist`/`Abgelaufen`; die Zeit in Millisekunden ist dieselbe, kein Test ist sonst umgeschrieben.
- **Commit-Messages** `aa07a01..HEAD`: keine `SPEC-`/`ARC-`-Kennung. `//nolint`: keines im Diff. Handbuch und `README.md` unberührt. `make abdeckung-check` ohne Befund. Kein Move im Diff (§3.3).
- **Gates.** `make gates` auf dem sauberen Baum nach dem Commit dieses Reports: Ergebnis im Verdikt.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-640 | MEDIUM | Der Kommentar in `.a-check.yml` und die Zeile „grüne Mutanten“ in §7 sagen, Dateien in keiner Schicht (Testhilfen unter `internal/`, `test`, `cmd`) prüfe a-check nicht, die Einträge seien redundant. Für die Regel `tech` stimmt das nur für den Composition Root (`internal/bootstrap`, `cmd`); entfernt man `test`, `internal/testuhr`, `internal/testzeit` und `internal/testpki` aus der Liste, meldet a-check 22 `tech-leak`. §7 führt den Mutanten als grün und äquivalent. Ein Leser, der dem Kommentar folgt, hält tragende Einträge für Zierde. | `AGENTS.md` §3.11 (Prosa über ein Gate nennt keine Zusage ohne Gegenprobe) | `.a-check.yml` · `Dateien in keiner Schicht (Testhilfen unter internal/) prüft`; Plan · `a-check meldet sie als Dateien in keiner Schicht, ungeprüft` | ja — Kopie ohne die fünf Einträge, a-check rot | Aussage über Gate-Verhalten ohne Gegenprobe, falsch |
| F-641 | LOW | Fall `bootstrap-time` der Gegenprobe bleibt grün, auch wenn `internal/bootstrap` aus der Liste fehlt (Composition Root ist ausgenommen); die Beschreibung („im Composition Root zugelassen“ in `harness/README.md`, `harness/mk/arch-negativ.mk`, Kopf des Skripts) sagt damit etwas, das die Regel nicht trägt. Der Fall hält den Zustand fest, nicht den Eintrag; §7 sagt das, die drei Beschreibungen nicht. | `AGENTS.md` §3.11 | `harness/README.md` · `im PGWire-Adapter und im Composition Root zugelassen`; `tools/arch/a-check-negativ.sh` · `fall bootstrap-time` | ja — Kopie ohne `internal/bootstrap` in der Liste, Fall bleibt grün | Gegenprobe-Fall kann an der beschriebenen Regel nicht rot werden |
| F-642 | LOW | Der Mutant `zeitpunkt < jetzt` statt `<=` in `systemuhr.go` bleibt grün: ein erreichter Zeitpunkt wartet dann bis zu eine Millisekunde (Timer mit Dauer 0, nächste Lesung), die Tests unterscheiden das nicht (`u.Jetzt()-1` ist vergangen, nicht erreicht). §7 führt für die Fake-Uhr den Mutanten `<` und für die Systemuhr nur „Sofort-Rückkehr entfernt“, den grünen Mutanten der Systemuhr nennt es nicht. | `AGENTS.md` §3.10 | `internal/bootstrap/systemuhr.go` · `if zeitpunkt <= jetzt`; Plan · `Systemuhr: erreichter oder vergangener Zeitpunkt sofort` | ja — Kopie mit `<`, `go test ./internal/bootstrap/` grün | Grüner Mutant der Grenze nicht in §7 geführt |
| F-643 | LOW | §6 Risiko 3 nennt als Gegenmittel, dass `slice-v1-abschluss-zeitangaben-aufzeichnen` und `slice-v1-abschluss-zeitangaben` die Signatur in §1 nennen. Der Diff legt die Namen fest (`Jetzt() int64`, `WarteBis(ctx, zeitpunkt int64) error`, `testuhr.Neu`, `Vorruecken`, `NeueSystemuhr`), beide Nehmer sagen weiter nur „der Port“, „die Uhr des Ports“, ohne Namen; `grep` nach `WarteBis` findet sie im Nehmer nicht. Das Gegenmittel ist offen, die Namen liegen jetzt vor. Nehmer nehmen an (kein Ausschluss trifft), es ist keine Adresse, die nicht annimmt. | `AGENTS.md` §3.9 (betroffene Folge-Slices nachziehen), §3.13 | Plan · `nennen die Signatur in §1, nicht eine eigene`; `slice-v1-abschluss-zeitangaben-aufzeichnen` §1 · `Dieser Slice benutzt den Port, wie jener ihn liefert` | ja — `grep WarteBis` in den beiden Nehmern leer | Gegenmittel eines Risikos nicht an den Nehmern eingetragen |
| F-644 | INFO | `testuhr.Uhr.WarteBis`: sind der Wecker (`Vorruecken`) und das Ende des Kontexts gleichzeitig bereit, wählt `select` zufällig und kann `nil` liefern, obwohl der Kontext beendet ist; die Zusage „ein beendeter Kontext geht vor“ gilt dann nur vor und nach dem Warten, nicht im Gleichzeitigen. Die Systemuhr prüft den Kontext nach dem Timer zuerst. Nur Testhilfe, kein Test greift darauf. Zuständig bei Bedarf: Implementer. | `LH-FA-21.a` *Uhr* | `internal/testuhr/testuhr.go` · `case <-w.geweckt:` | nein | Fake-Uhr folgt der Rangfolge der Zusage im Gleichzeitigen nicht |
| F-645 | INFO | Zu den Abweichungen des Implementers, geprüft und nicht als Befund geführt: (1) `NeueSystemuhr` hat keinen Aufrufer außer den Tests; Plan §3 sagt „nur bereitgestellt“, `LH-FA-21.a` sagt „stellt der Composition Root bereit“, `make lint` und `make gates` laufen (Ergebnis im Verdikt). (2) Die Prüfung `var _ driven.Uhr = (*testuhr.Uhr)(nil)` steht in `internal/bootstrap/systemuhr_test.go`, nicht beim Fake; bricht der Fake den Port, fällt es in einem fremden Paket auf. (3) `Vorruecken` ohne Überlaufschutz: Testhilfe, die Zusage nennt nur die negative Zahl. (4) §3 des Plans nennt `internal/bootstrap/export_test.go` und die zwei `doc.go` nicht. (5) Der ferne Zeitpunkt (`math.MaxInt64`, `math.MinInt64`, Stufen der größten Dauer) ist im Code behandelt und in §6 nicht als Randform genannt; `LH-FA-21.a` bestimmt das Verhalten (warten, bis die Lesung den Zeitpunkt erreicht hat), der Code entscheidet nichts Neues. | Maintainability | `internal/bootstrap/systemuhr.go` · `func NeueSystemuhr`; `internal/bootstrap/systemuhr_test.go` · `var _ driven.Uhr = (*testuhr.Uhr)(nil)` | nein | Sammelhinweis ohne Aktion |
| F-646 | INFO | Zeitfenster der Tests: `TestFristLaeuftNachMillisekundenAb` (60 ms Frist, Prüfung unmittelbar nach der Erzeugung), `TestSystemuhrAbbruchWaehrendDesWartens` (Untergrenze 30 ms bei Abbruch nach 40 ms) und die 50-ms-Fenster in `bleibt` haben Zeitfenster; unter 1 CPU mit sechs Schleifen waren 30 Läufe grün. Die Fenster hängen an der Annahme, dass der Läufer zwischen zwei Anweisungen nicht länger als 10 bis 60 ms steht. `TestSystemuhrFernerZeitpunktDrehtNicht` liest die CPU-Zeit des ganzen Prozesses (`//go:build unix`); sein Abstand zum Mutanten ist groß (150 ms Grenze, 350 ms Mutant, nahe 0 ehrlich). | Maintainability | `internal/testzeit/testzeit_test.go` · `f := testzeit.NeueFrist(60)`; `internal/bootstrap/systemuhr_cpu_unix_test.go` · `cpu > 150*time.Millisecond` | nein | Zeitfenster in Tests, bei Last beobachtet ohne Flackern |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/ports/driven/uhr.go`: ganze Millisekunden als `int64`, kein Typ der Zeitbibliothek, Zusagen des Kommentars gegen `LH-FA-21.a` *Uhr* gelesen (beendeter Kontext vor und während, erreichter oder vergangener Zeitpunkt sofort `nil`) | geprüft, ohne Befund |
| `internal/bootstrap/systemuhr.go`: Ursprung bei Erzeugung, `time.Since` aus der monotonen Zeit, abgeschnitten, Überlauf (`rest` und Stufe) von Hand gerechnet: `rest` ist positiv und höchstens `MaxInt64/1e6`, `MinInt64` liefert nie eine Subtraktion; Abbruch gibt `ctx.Err()`; Stufen schlafen | geprüft, ohne Befund; Mutanten (Stufe, halbe Wartezeit, Rundung) rot, grüne siehe F-642 |
| `internal/testuhr/testuhr.go`: rückt nur auf Anweisung, negativ löst Panik aus (null erlaubt), weckt jeden erreichten Wartenden, Abbruch trägt den Wartenden aus, Sperre in `Vorruecken` hält Wecken und Löschen zusammen | geprüft, ohne Befund (Gleichzeitigkeit: F-644); Mutant ohne Vorab-Prüfung rot |
| `internal/testzeit/testzeit.go`: einzige Stelle der Tests des Kerns mit `time`, Einheit Millisekunden | geprüft, ohne Befund |
| `internal/hexagon/services` (drei Testdateien): nur die Aufrufe von `time` ersetzt, Werte in Millisekunden gleich dem Altwert, Importe und Reihenfolge `gofmt`-sauber; kein Verhalten geändert | geprüft, ohne Befund |
| `.a-check.yml`: Regel `^time$` mit `match: regex`, Eintrag trifft nicht `runtime`; die Zulassung ist **keine Lockerung** nach `AGENTS.md` §3.6: vor dem Diff war `time` überall erlaubt, die Regel ist neu und eine Verschärfung; die zwei wirkungslosen Einträge (`internal/bootstrap`, `cmd`) lockern nichts, weil der Composition Root ohnehin ausgenommen ist (Probe im Bootstrap grün ohne den Eintrag). Aussage im Kommentar: F-640 | geprüft; Befund F-640 |
| `tools/arch/a-check-negativ.sh`, `harness/mk/arch-negativ.mk`, `harness/README.md`: 23 Fälle (17 + 6) zählen, Kopfkommentar, Meldung am Ende und Beschreibung nennen dieselben Orte und Pakete; „Driven Ports“ in der Beschreibung ist das Geprüfte (Fall `ports-time` liegt in `ports/driven`); `ADR-0015` in Bindung und Kopf | geprüft, ohne Befund außer F-641; Mutanten der Regel (Eintrag, `match: regex`, `internal/adapters`) von der Gegenprobe gefangen |
| Core-Reinheit und Schichten: Kern bindet `time` weder in Code noch in Tests ein, `testuhr`/`testzeit` liegen außerhalb des Hexagons, keine Adapterkante | geprüft, ohne Befund |
| Plan §1, §3, §6, §7, §8 gegen den Code (`AGENTS.md` §3.9): drei Liefer-Punkte, zwei Schichten (Kern, Bootstrap mit Logik), Dateien und Namen in §3 stimmen bis auf F-645 (4); §6 hat sich im Code-Commit nicht geändert, keine neue Randform dort eingetragen (`AGENTS.md` §3.12); Größe 761 Zeilen, in einer Sitzung prüfbar | geprüft, ohne Befund außer F-643 |
| Hard Rules: kein `//nolint`, keine `SPEC-`/`ARC-`-Kennung in den vier Commit-Messages, Kommentare im Code sagen Zusage, Kopplung oder Grenze (§3.7), Handbuch und `README.md` unberührt, Roadmap-Zeile nennt Zustand ohne Chronik | geprüft, ohne Befund |
| Suppression und Gate-Lockerung: nichts in `.golangci.yml`, nichts in `.d-check.yml` geändert | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Aussage über Gate-Verhalten ohne Gegenprobe, falsch · Gegenprobe-Fall kann an der beschriebenen Regel nicht rot werden · Grüner Mutant der Grenze nicht in §7 geführt · Gegenmittel eines Risikos nicht an den Nehmern eingetragen · Fake-Uhr folgt der Rangfolge der Zusage im Gleichzeitigen nicht · Zeitfenster in Tests

## Verdikt

**Merge-blockierend:** nein (kein HIGH). Die Zusagen von `LH-FA-21.a` *Uhr* sind in Port, Fake-Uhr und Systemuhr gefangen; die nachgefahrenen Mutanten sind rot bis auf die in F-642 genannte Grenze. Die Zulassung in `.a-check.yml` lockert nichts.

**Übergabe:**

- **F-640** an den **Implementer**: Kommentar in `.a-check.yml` und die Zeile in §7 der Aussage nach dem Gate-Verhalten anpassen (`AGENTS.md` §3.11, §3.9); ein Eintrag der Liste, der etwas trägt, ist nicht redundant.
- **F-641** an den **Implementer**, gemeinsam mit F-640.
- **F-642** an den **Implementer**: Mutant in §7 als grün führen oder fangen.
- **F-643** an den **Planner**: Nehmer nennen die Namen; wer sie einträgt, zählt nach (`AGENTS.md` §3.13 *Nachzählen beim Eintragen*), beide Nehmer liegen bei zwei und drei Liefer-Punkten und zwei Schichten.
- **F-644** bis **F-646** ohne Blockade.
- Die **Finding-Klassen** gehen in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg; DoD- und Spec-Konformität prüft der Verifier.
