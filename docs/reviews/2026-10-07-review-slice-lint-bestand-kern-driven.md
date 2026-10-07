# Review-Report: slice-lint-bestand-kern-driven — 2026-10-07

**Review-Art:** Code-Review (Umbau ohne Verhaltensänderung, Charakterisierungstests, Kommentare). Geprüft wird gegen Plan §1, §3 und §6 (Stand Architect `7ac016b`), `SPEC-049`, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) und die Hard Rules (Modul 10). Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff 7ac016b..ac2e662`:

- `cedd891`: die Charakterisierungstests `TestValidateFehlerReihenfolge` und `TestUnmarshalFehlerReihenfolge`, dazu §1 und §3 des Plans.
- `d840888`: der Umbau von `Group.validate`, `(*cursor).objekte`, `toResponse` und `fromDTO`, die Doc-Kommentare, `ctx` → `_`, `QF1001`, der Kommentar V-81 und §3 des Plans.
- `ac2e662`: die Belege in §7.

**Skill:** `.harness/skills/reviewer.md` @ `de770f7` (unverändert bis `ac2e662`)
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-lint-bestand-kern-driven.md`: ganz gelesen am Stand `ac2e662`, dazu die Plan-Diffs von `cedd891`, `d840888` und `ac2e662`. §7 *Belege des Implementers* habe ich gelesen, den Bericht des Implementers habe ich nicht.
- `docs/plan/planning/open/slice-harness-coverage.md` §1 und §6, `docs/plan/planning/open/slice-harness-mutation.md` §1.
- `spec/spezifikation.md` `SPEC-049`, [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes).
- [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) (Accepted), [ADR-0001](../plan/adr/0001-hexagonale-architektur.md).
- `AGENTS.md` §3.2, §3.3, §3.6, §3.7, §3.9 bis §3.12.
- Produkt-Code an `7ac016b` und `ac2e662`: `internal/hexagon/model/extended.go`, `recording.go`, `internal/hexagon/services/replay.go` (auch `NewReplayService` und `zuordnen`), `internal/adapters/driven/postgres/upstream.go`, `internal/adapters/driven/recording/yaml.go`, `internal/hexagon/ports/driven/upstream.go`.
- Vorherige Reports: die Reviews zu `slice-harness-blackbox-driven`, `slice-harness-blackbox-pgwire` und `slice-harness-blackbox-einstieg`, dazu F-330 und F-431 als Präzedenz für Adressen, die nicht annehmen. Die höchste vergebene Nummer vor diesem Lauf war F-458.

**Ausgeführte Läufe:**

- `make lint` im Repo: Exit 2, `16 issues:`. Unter `internal/hexagon/model`, `internal/hexagon/services`, `internal/adapters/driven/postgres` und `internal/adapters/driven/recording` steht keine Zeile.
- `git archive` von `ac2e662` und von `cedd891` in je eine frische Kopie unter dem Scratch-Pfad (`tar -x --touch`, kein `cp -p`). Daraus ein eigenes Image der Stufe `deps` mit eigenem Tag. Die vier Pakete laufen dort mit `go test -count=1 .` ohne Netz grün.
- Die Charakterisierungstests am alten Code (`cedd891`) mit `-run Reihenfolge -v`: 13 von 13 Fällen von `TestValidateFehlerReihenfolge` und 5 von 5 Fällen von `TestUnmarshalFehlerReihenfolge` sind grün.
- 31 eigene Mutationen. Jede lief in einer frischen Kopie aus `git archive` und wurde per `perl` angewendet. Ob die Ersetzung gegriffen hat, habe ich über eine Prüfsumme kontrolliert. Unit-Tests liefen per Bind-Mount im eigenen `deps`-Image. Die Integrationstests liefen per `make test-integration` in der Kopie. Im Runner der Kopie stand dafür ein eigener Image-Tag, den ich danach entfernt habe. Danach habe ich auch die Kopien gelöscht. Im Repo hat sich bis auf diese Datei nichts geändert. Ein hostweites Docker-Aufräumen gab es nicht.
- `make docs-check` vor dem Commit (siehe Ende).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-459 | MEDIUM | Die Adresse `slice-harness-coverage` nimmt die Sendung nicht an. §7 weist ihr zwei grüne Mutanten mit Test-Idee zu (`bestaetigt[m.Type]--` und `TableOID`). Plan §6 weist ihr jeden ungefangenen Pfad zu: „ist das ein Befund für `slice-harness-coverage`, kein neuer Test hier“. Ihr §1 schließt aber „Neue Tests, um eine Zielzahl über dem gemessenen Stand zu erreichen“ aus. Ihr Gate misst die Anweisungs-Abdeckung, und beide Zeilen laufen unter Tests. Der Mutant überlebt also auf einer abgedeckten Zeile, und die Messung zeigt ihn nie. `slice-harness-mutation` schließt „Neue Tests für Mutanten, die heute überleben“ ebenfalls aus. Folge: Mit der Closure liegen die Test-Ideen in `done/`, kein offener Slice führt sie, und die Zusage über `extendedNachspielen` (F-460) bleibt dauerhaft ungeprüft. Die Route steht seit `175fd2d` in §6 und betrifft damit den Architect. | Baseline-Regelwerk `v6.13.0` `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice, Klasse 1 („Die Adresse muss die Sendung annehmen“); Plan §6 Risiko *Verhalten ändert sich unbemerkt*; Präzedenz F-330 | `docs/plan/planning/in-progress/slice-lint-bestand-kern-driven.md` · „fängt keiner, ist das ein Befund für `slice-harness-coverage`, kein neuer Test hier“; `docs/plan/planning/open/slice-harness-coverage.md` · „Neue Tests, um eine Zielzahl über dem gemessenen Stand zu erreichen“ | ja (Lesen; M3/M3i und P8 grün) | Folge-Slice-Adresse nimmt die Sendung nicht an |
| F-460 | LOW | Der Satz „Die ersten n Nachrichten einer Art mit n Bestätigungen in der Interaktion gelten als angenommen“ steht jetzt im Doc-Kommentar der neuen Funktion `extendedNachspielen`, also in den hinzugefügten Zeilen des Diffs. Ein Test hält das „die ersten n“ nicht. Mein Mutant M3 ohne `bestaetigt[m.Type]--` ist in den Unit-Tests von `services` grün und auch in `make test-integration` (M3i). §3.11 nennt zwei Wege: enger fassen oder ein Test. §7 nimmt keinen davon. Die Begründung „weil er den Bestand beschreibt“ ist kein Kriterium von §3.11: Der Satz beschreibt den Code richtig, ein Test prüft ihn trotzdem nicht. Anders als V-81, das der Architect in §6 entschieden hat, steht diese Entscheidung nur in §7 beim Implementer. §1 führt sie auch nicht als bewusst stehenden Bestand (Klasse 2). Die Einstufung ist LOW wie bei F-401: Der Text ist alt, der Code unverändert, und die Wirkung liegt im Diagnosetext einer Abweichung, nicht in der Antwort des Replays. | `AGENTS.md` §3.11; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Plan §1 (Klasse 2), §6 | `internal/hexagon/services/replay.go` · „Die ersten n Nachrichten einer Art mit n Bestätigungen in der Interaktion gelten als angenommen“; `docs/plan/planning/in-progress/slice-lint-bestand-kern-driven.md` · „der Slice fasst ihn nicht enger, weil er den Bestand beschreibt“ | ja (M3, M3i) | Kommentar-Zusage ohne Test, verschoben statt gefasst |
| F-461 | LOW | Die Liste der grünen Mutanten in §7 ist unvollständig. Sie nennt nur `bestaetigt[m.Type]--` und `TableOID`. Zwei weitere Mutanten in umgebauten Funktionen bleiben vor und nach dem Umbau grün. Erstens `spalten` ohne `ColumnNumber`: grün in den Unit-Tests von `postgres` und in `make test-integration`, sowohl an `ac2e662` (P2) als auch an `7ac016b` (P2alt). Zweitens die Session-Nummer in den Fehlern einer Interaktion: Mit `geprueftFromDTO(1, ii, id)` bleibt `recording` grün (Y2). Am alten Code bleibt die gleichwertige Konstante in beiden Meldungen ebenfalls grün (Y2alt, Y2alt2). Den neuen Parameter `sessionID` hält damit kein Test. Kein Test erzeugt einen Interaktionsfehler in einer Session ab 2. Das Risiko *Verhalten ändert sich unbemerkt* verlangt, solche Fälle als Befund zu benennen. Die Adresse selbst ist in F-459 behandelt. | Plan §6 Risiko *Verhalten ändert sich unbemerkt*; `AGENTS.md` §3.10 | `docs/plan/planning/in-progress/slice-lint-bestand-kern-driven.md` · „*Grüne Mutanten, eingeordnet*“; `internal/adapters/driven/postgres/upstream.go` · „ColumnNumber: f.TableAttributeNumber,“; `internal/adapters/driven/recording/yaml.go` · „i, err := geprueftFromDTO(sd.ID, ii, id)“ | ja (P2, P2alt, Y2, Y2alt) | Liste grüner Mutanten unvollständig |
| F-462 | LOW | Die Ausnahme für die Charakterisierungstests in §1 hat der Implementer im Test-Commit `cedd891` eingeführt. In der Fassung des Architect (`7ac016b`) schloss §1 „neue Fälle“ aus, und §6 sagte „kein neuer Test hier“. Ebenfalls in `cedd891` legt §3 die Testliste „vor dem Umbau“ der DoD auf den Stand mit den neuen Tests fest. Ein Vermerk des Architect fehlt. Inhaltlich halten die Tests den Bestand fest: Sie sind an `cedd891` grün, und bei Vertauschungen werden sie rot (S1 bis S3 und die Tabelle in §7). Nach Modul 5 hat aber den Plan geändert, wer mitnimmt, was ausgeschlossen war. Im Plan steht jetzt nicht mehr getrennt, was der Architect entschieden hat und was der Implementer nachgetragen hat. Messen Review und Verifikation „gegen §1“, messen sie gegen eine Fassung, die der Implementer selbst erweitert hat. | Baseline-Regelwerk `v6.13.0` `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice („Wer später mitnimmt, was hier ausgeschlossen war, hat den Plan geändert“); Plan §4 Rückführung `in-progress → open`; `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-lint-bestand-kern-driven.md` · „Ausgenommen sind Charakterisierungstests“ und „Die Testliste „vor dem Umbau“ der DoD ist die mit ihnen“ | ja (`git show cedd891 -- docs/plan`) | Abgrenzung in §1 vom Implementer im Code-Commit erweitert |
| F-463 | INFO | Die Kopfkommentare beider Charakterisierungstests schreiben die Reihenfolge als Festlegung: „meldet … genau eine, und welche, steht fest“. Dazu nennt der neue Doc-Kommentar von `Group.validate` die Reihenfolge („die erste verletzte meldet ihren Fehler“). Weder die Spezifikation noch eine ADR entscheidet, welcher von zwei Fehlern zuerst kommt. §1 sagt, die Tests sagten „nichts Neues zu“. §3.11 ist eingehalten, denn der Test prüft genau das, was Kommentar und Test zusagen. Wer die Reihenfolge später ändert, liest den Test aber als Vertrag und nicht als Bestandsaufnahme. Hinweis an den Architect, der entscheidet, ob die Reihenfolge Vertrag ist. | Plan §1, §6 Randform *Komplexitäts-Bereinigung ohne Verhaltensänderung*; `AGENTS.md` §3.12; Rolle Architect | `internal/hexagon/model/extended_test.go` · „meldet Validate genau eine, und welche, steht fest“; `internal/adapters/driven/recording/yaml_test.go` · „meldet Unmarshal genau eine, und welche, steht fest“ | nein | — (Hinweis) |
| F-464 | INFO | Zur Gleichwertigkeit des Umbaus. Ich habe jede Funktion gegen `7ac016b` gelesen, Ergebnisse und Fehlertexte sind gleich. Die Prüfungen laufen in derselben Reihenfolge. Das gilt in `Group.validate` (Client, sync/flush am Ende, Server, ready_for_query am Ende) und in `fromDTO` (Kennung, leere Session, je Interaktion Nummer, `offset_ms`, Form). `toResponse` ordnet die Fälle nur neu, aber ein Type-Switch auf disjunkte Zeigertypen hängt nicht von der Reihenfolge ab. `%T` an `msg` statt an `m` zeigt im `default`-Zweig denselben dynamischen Typ. `objekte` hat einen Wert-Empfänger, das ist richtig, denn es mutiert nur Maps. Den Halt ohne `gi > c.gruppe` hält §7 für gleichwertig, und das bestätige ich (M5). Gleichwertig sind außerdem `break` statt `return o` (M6), `si+1` statt `sd.ID` (Y1), `ii+1` statt `id.Sequence` in der Formmeldung (Y5) und `!= ClientFlush` statt `== ClientSync` (V1). Jede dieser Varianten hängt an einer Prüfung, die vorher läuft. Das ist der Bestand, gegen den der Verifier die DoD prüft. | Plan §6 Randform *Komplexitäts-Bereinigung ohne Verhaltensänderung*; Rolle Verifier | `internal/hexagon/services/replay.go` · „return pi == c.pos && (gi > c.gruppe \|\| gi == c.gruppe && ni >= c.nachricht)“ | ja (M5, M6, Y1, Y5, V1) | — (Hinweis) |

**Eigene Mutationen** (jede in einer frischen Kopie; ohne Mutation sind alle genannten Tests grün):

| ID | Mutation | Lauf | Ergebnis |
|---|---|---|---|
| M1 | Halt ohne `pi == c.pos &&` | Unit `services` | rot, `TestReplayExtendedDiagnoseLebensdauer`, `TestReplayExtendedDiagnoseArt` |
| M2 | `extendedNachspielen` liefert am Ende `false` | Unit `services` | rot, `TestReplayExtendedDiagnoseLebensdauer` |
| M3 | ohne `bestaetigt[m.Type]--` | Unit `services` | grün (F-460) |
| M3i | wie M3 | `make test-integration` | grün (F-460) |
| M4 | `bestaetigt[m.Type] >= 0` | Unit `services` | rot, `TestReplayExtendedDiagnoseLebensdauer` |
| M5 | Halt ohne `gi > c.gruppe` | Unit `services` | grün, gleichwertig (F-464) |
| M6 | `break` statt `return o` nach dem Halt | Unit `services` | grün, gleichwertig |
| L1 | `NewReplayService` nimmt auch leere Sessions in `frei` auf | Unit `services` | rot, `TestReplayLebendpruefungAufgezeichnet`, `TestReplayNichtsZuMelden` (Grund der Invariante V-81 ist getestet) |
| P0 | `Columns: nil` | `make test-integration` | rot, u. a. `TestE2EErgebnisartenEinfach`, `TestE2EReplayExtendedPgx` |
| P1 | ohne `TypeModifier` | Unit `postgres` grün; `make test-integration` | rot, `TestE2EFehlerreplayEinfach`, `TestE2EErgebnisartenEinfach` |
| P2 | ohne `ColumnNumber` | Unit `postgres`; `make test-integration` | grün (F-461) |
| P2alt | wie P2 an `7ac016b` | `make test-integration` | grün (F-461) |
| P3 | ohne `TypeSize` | Unit grün; `make test-integration` | rot, `TestE2EFehlerreplayEinfach`, `TestE2EErgebnisartenEinfach` |
| P4 | ohne `TypeOID` | Unit grün; `make test-integration` | rot, u. a. `TestE2EReplaySelect1` |
| P5 | `Name: ""` | Unit grün; `make test-integration` | rot, u. a. `TestE2ERecordMehrereInteraktionen` |
| P6 | Fall `ParameterDescription` entfernt | Unit `postgres` | rot, `TestSendUndReceive` |
| P7 | `ParamTypes: nil` | Unit `postgres` | rot, `TestSendUndReceive` |
| P8 | ohne `TableOID` | Unit `postgres` | grün, wie §7 angibt |
| Y1 | `geprueftFromDTO(si+1, …)` | Unit `recording` | grün, gleichwertig |
| Y2 | `geprueftFromDTO(1, …)` | Unit `recording` | grün (F-461) |
| Y2alt, Y2alt2 | an `cedd891` `1` statt `sd.ID` in der Form- bzw. `offset_ms`-Meldung | Unit `recording` | grün (F-461) |
| Y3 | Session ohne `Startup` | Unit `recording` | rot, `TestRoundtrip`, `TestExtendedRoundtrip` |
| Y4 | Session ohne `ServerParameters` | Unit `recording` | rot, `TestRoundtrip`, `TestExtendedRoundtrip` |
| Y5 | `ii+1` statt `id.Sequence` in der Formmeldung | Unit `recording` | grün, gleichwertig |
| Y6 | `offset_ms < -1` | Unit `recording` | rot, `TestFelderDerVersion1`, `TestUnmarshalFehlerReihenfolge/offset_ms_vor_der_Form` |
| V1 | `!= ClientFlush` statt `== ClientSync` | Unit `model` | grün, gleichwertig |
| V2 | `len(g.Server) == 1` statt `== 0` | Unit `model` | rot, `TestValidateGueltig`, `TestValidateFehler/letzte_Gruppe_ohne_Server-Nachricht` |
| S1 | ready_for_query am Ende vor `validateClient` geprüft | Unit `model` | rot, drei Fälle von `TestValidateFehlerReihenfolge` |
| S2 | Zielart vor Stellung in `validateClient` | Unit `model` | rot, `TestValidateFehlerReihenfolge/Stellung_vor_Zielart_in_derselben_Nachricht` |
| S3 | leere Session vor Kennung in `sessionFromDTO` | Unit `recording` | rot, `TestUnmarshalFehlerReihenfolge/Kennung_vor_Session_ohne_Interaktion` |

## Antwort auf die Schwerpunkte

1. **Ohne Verhaltensänderung.** Gleichwertig, siehe F-464. Die Reihenfolge bei mehreren gleichzeitigen Fehlern ist in `Group.validate` und `fromDTO` unverändert. Die Charakterisierungstests sind am alten Code grün, und bei Vertauschungen werden sie rot (S1 bis S3 und die Mutationstabelle in §7). Nicht beobachtbar sind zwei Vertauschungen, und deshalb deckt sie kein Test: Typ und Zielart in derselben Client-Nachricht (eine unbekannte Art ist nie describe/close) und Typ und Stellung in derselben Server-Nachricht (ready_for_query ist bekannt).
2. **Charakterisierungstests.** Sie halten den Bestand fest, mehr nicht. Erlaubt sind sie aber erst durch die Ausnahme, die der Implementer in `cedd891` in §1 geschrieben hat (F-462). In der Fassung des Architect waren sie ausgeschlossen. Zum Wortlaut ihrer Kopfkommentare siehe F-463.
3. **Kommentar V-81.** Er entspricht §6 Wort für Wort und sagt für den leeren Zweig keinen Wert zu. Sein Grund ist getestet: L1 macht drei Tests rot. „schützt allein den Index“ beschreibt den Zweig und sagt nichts zu. Kein Befund.
4. **Verschobener Kommentar über die Bestätigungen.** Die Grenze in §7 reicht nach §3.11 nicht (F-460). Den Slice-Plan unterstützt dort keine Entscheidung des Architect, und die Adresse nimmt nicht an (F-459).
5. **Grüne Mutanten.** `slice-harness-coverage` nimmt nicht an (F-459). Der Halt in `objekte` ist gleichwertig (M5, F-464). Zwei weitere grüne Mutanten fehlen in §7 (F-461).
6. **Doc-Kommentare.** `EndClosed` bis `EndFailed` tragen jetzt die Form „`<Name> heißt: …`“, der Inhalt ist gleich. „Die Zielarten von describe und close.“ trifft auf `TargetStatement` und `TargetPortal` zu. Die neuen Kommentare an `toResponse`, `ohneFelder`, `werte`, `sessionFromDTO`, `geprueftFromDTO`, `validateClient` und `validateServer` beschreiben den Code, und je eine Mutation in §7 oder hier fängt sie. Kein Befund außer F-460 und F-463.
7. **Eigene Mutationen.** Siehe die Tabelle oben.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/model/` (`extended.go`, `recording.go`, `extended_test.go`) | geprüft. Gleichwertig (F-464), De Morgan richtig umgeformt; Hinweis F-463. |
| `internal/hexagon/services/replay.go` | geprüft. Gleichwertig, V-81 nach §6; Befund F-460. |
| `internal/adapters/driven/postgres/upstream.go` | geprüft. Gleichwertig. `_` für `ctx` passt zum Port, denn `driven.UpstreamSession` sagt für `Query` nichts über den Kontext zu. Befund F-461 (`ColumnNumber`). |
| `internal/adapters/driven/recording/` (`yaml.go`, `yaml_test.go`) | geprüft. Gleichwertig; Befund F-461 (Session-Nummer). |
| Plan `docs/plan/planning/in-progress/slice-lint-bestand-kern-driven.md` | geprüft. Befunde F-459, F-461, F-462. §3 folgt dem Diff. |
| Hard Rule 3.2, 3.6 (Suppression, Gate-Lockerung) | geprüft, ohne Befund. Kein `//nolint` im Diff, `.golangci.yml` ist unverändert. |
| Hard Rule 3.3 | geprüft, ohne Befund. Im Diff gibt es keinen Move. |
| Hard Rule 3.7 (Kommentar-Klassen) | geprüft, ohne Befund. Die neuen Kommentare beschreiben den Ist-Zustand, keiner nennt eine verworfene Alternative oder einen abwesenden Text. |
| Hard Rule 3.9 | geprüft. §3 ist in `d840888` im selben Commit nachgezogen. Zu §1 siehe F-462. |
| Hard Rule 3.12 | geprüft, ohne Befund. Die Randformen in §6 hat der Architect vor dem Code entschieden (`7ac016b`). Die Code-Commits bringen in §6 keine neue Randform. |
| Core-Reinheit, [ADR-0001](../plan/adr/0001-hexagonale-architektur.md) | geprüft, ohne Befund. Kein neuer Import, über keine Schichtgrenze wird etwas verschoben, keine exportierte Signatur ist geändert. |
| Commit-Messages `cedd891`, `d840888`, `ac2e662` | geprüft, ohne Befund. Jede nennt `slice-lint-bestand-kern-driven`, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) und `LH-QA-07`, keine nennt eine Struktur-ID. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 2 |

**Summary:** 0 HIGH · 1 MEDIUM · 3 LOW · 2 INFO (F-459 Folge-Slice-Adresse `slice-harness-coverage` nimmt die grünen Mutanten nicht an; F-460 verschobener Kommentar über die Bestätigungen sagt ungeprüft zu; F-461 Liste grüner Mutanten unvollständig, `ColumnNumber` und Session-Nummer; F-462 Ausnahme in §1 vom Implementer im Test-Commit eingeführt; F-463 Charakterisierungstests als Festlegung formuliert; F-464 Umbau gleichwertig). Wiederkehrende Klassen: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-460); Folge-Slice-Adresse nimmt nicht an (F-459, wie F-330 und F-431).

**Finding-Klassen dieses Laufs:** Folge-Slice-Adresse nimmt die Sendung nicht an · Kommentar-Zusage ohne Test, verschoben statt gefasst · Liste grüner Mutanten unvollständig · Abgrenzung in §1 vom Implementer im Code-Commit erweitert

## Verdikt

**Merge-blockierend:** nein. Der Umbau verhält sich wie `7ac016b`, und `make lint` meldet unter den vier Pfaden nichts. Vor der Closure braucht F-459 einen Ausgang, denn ohne ihn hat das Risiko *Verhalten ändert sich unbemerkt* keine tragende Adresse.

**Übergabe:** F-459, F-462 und F-463 gehen an den Architect, der §1 und §6 führt. F-460 und F-461 gehen an den Implementer. F-464 geht an den Verifier. Dieser Report ersetzt keine Verifikation.
