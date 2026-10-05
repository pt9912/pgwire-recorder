# Zweites Folge-Review: slice-extended-query-replay — 2026-10-05

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** Commit `b697462` (`git diff c54ed19..b697462`). Er soll die Findings F-330 bis F-338 aus dem [Folge-Review](2026-10-05-folge-review-slice-extended-query-replay.md) abarbeiten. Schwerpunkte laut Auftrag:

- Status je Finding.
- Die neue SQL-Herleitung der Diagnose bei bind, execute, describe und close, geprüft gegen LH-FA-10.a und LH-FA-18.a. Dazu gehören die Auflösung von Statement und Portal aus den Client-Nachrichten vor dem Cursor, unbenannte Objekte, Überschreiben, Close, „unbekannt“ und „keine“. Parameterwerte dürfen nie erscheinen.
- Der Test-Wrapper `spaeteFrist` und ob `TestReplayHerunterfahrenSpaeteFrist` deterministisch ist.
- Ob `TestE2EReplayExtendedSigtermMittenInFolge` den Vergleich „wie aufgezeichnet“ trägt.
- Der Zuschnitt der Folge-Slices und der Welle (`AGENTS.md` §3.9, Größenregel).
- Ob die Entscheidung trägt, den Folge-Slice nur im Plan und nicht im Code-Kommentar zu nennen.

**Ablage:** Eigene Datei, Nummerierung fortlaufend, wie beim ersten Folge-Review (Reviewer-Skill §Output-Schema: „Folgeläufe als neue Datei statt Überschreibung“).

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-extended-query-replay.md` (Kopf, §1 bis §3, §6), `docs/plan/planning/open/slice-v1-abschluss-betrieb.md` (Kopf bis §6), `docs/plan/planning/open/slice-replay-semantik-mismatch.md` (Kopf bis §4), `docs/plan/planning/welle-replay-semantik.md` (§1 bis §6)
- `spec/lastenheft.md` [LH-FA-10](../../spec/lastenheft.md), [LH-FA-13](../../spec/lastenheft.md), [LH-FA-18](../../spec/lastenheft.md)
- `spec/spezifikation.md`: LH-FA-10.a (Mismatch), LH-FA-18.a (Replay, Mismatch, `SPEC-033`), LH-FA-13.a (Frist, `PGR-E4006`), Meldungscode-Tabelle
- [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0007](../plan/adr/0007-strict-replay.md)
- `AGENTS.md` Hard Rules 3.3 bis 3.9 und §5; Register `BEO-REPO/plan-folgt-korrektur-nicht`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`
- Lebensdauer der Protokollobjekte nach der PostgreSQL-Protokolldokumentation (Extended Query). Das unbenannte Statement lebt bis zum nächsten Parse auf das unbenannte Statement oder bis zu einer einfachen Anfrage. Ein Portal lebt bis zum Ende der Transaktion, ein unbenanntes Portal zusätzlich nur bis zum nächsten Bind auf das unbenannte Portal. `Close` zerstört das Objekt.
- Vorherige Reports am Modul: [Review](2026-10-05-review-slice-extended-query-replay.md) (F-321 bis F-329), [Folge-Review](2026-10-05-folge-review-slice-extended-query-replay.md) (F-330 bis F-338)

**Ausgeführte Läufe im Repo:** `make docs-check` nach dem Anlegen dieser Datei. `make gates` lief nicht. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert (`git status --short` leer).

**Gegenproben und Mutationen:** Nur in Kopien (`git archive b697462`) im Scratchpad. Die Unit-Läufe liefen in einem Container aus der Stufe `deps` des `Dockerfile` (Tag `pgr-rev2-deps`, ohne Netz). `make test-integration` lief einmal in der unveränderten Kopie; es baut dabei das Tag `pgwire-recorder:integration` aus demselben Commit neu. Image `pgr-rev2-deps` und Kopien sind gelöscht.

| # | Lauf oder Mutation | Ergebnis |
|---|---|---|
| K1 | Kopie ohne Änderung: `gofmt -l .`, `go vet -tags integration ./...`, Unit-Tests von `internal/hexagon/services` und `internal/adapters/driving/pgwire` | grün, `gofmt` leer |
| K2 | Kopie ohne Änderung: `make test-integration` | grün, darunter `TestE2EReplayExtendedSigtermMittenInFolge` und `TestE2EReplayExtendedPipeline` |
| K3 | `TestReplayHerunterfahrenSpaeteFrist` unverändert, `-count=300`, einmal mit Standard-`GOMAXPROCS`, einmal mit `GOMAXPROCS=1` | 600 von 600 grün |
| M1 | `<-geweckt` entfernt (wie im Folge-Review) | rot, 20 von 20 (`TestReplayHerunterfahrenSpaeteFrist`, `unexpected EOF`). F-337 behoben |
| T1 | Zeitprobe ohne Verhaltensänderung: `replaySitzung` schläft 20 ms, nachdem der Use Case eine Nachricht ohne Antwort angenommen hat | rot, 5 von 5 (`1 Nachrichten statt 2 beim Use Case`) (F-341) |
| P1 | Sonde: Close auf `s1` in eigener Interaktion, danach `bind s1` ohne neues Parse, wo Parse erwartet ist | Diagnose nennt `empfangen "SELECT 1"`, das SQL des geschlossenen Statements (F-339) |
| P2 | Sonde: unbenanntes Statement, dann eine einfache Anfrage, dann `bind ""` ohne Parse | Diagnose nennt das SQL des durch die einfache Anfrage zerstörten Statements (F-339) |
| P3 | Sonde: `execute` auf das unbenannte Portal in der nächsten Interaktion ohne neues Bind | Diagnose `Anweisung erwartet "SELECT 1", empfangen "SELECT 1"`, obwohl das Portal mit dem `Sync` endete (F-339) |
| P4 | Sonde: `close` auf ein unbekanntes Portal, wo `describe` des unbenannten Portals erwartet ist | `erwartet "SELECT 1", empfangen unbekannt`, wie beschrieben |
| D1 | `portalSQL` sucht das Statement zum Zeitpunkt des Fehlers statt zum Zeitpunkt des Bind | grün (F-340) |
| D2 | `vorher` nimmt die Nachricht am Cursor mit (`ni > c.nachricht` statt `>=`) | grün (F-340) |
| F1 | `fmt` `%q` auf `[][][]byte{{nil}}` und `[][][]byte{{[]byte{}}}` | beide `[[""]]` (F-342) |

---

## Status der Findings F-330 bis F-338

| ID | Stand | Beleg |
|---|---|---|
| F-330 | behoben | `slice-v1-abschluss-betrieb` übernimmt die Frist für `replay` in §1 (`:34`, `:37`), in DoD-Punkt 1 (`:57`, „je Modus“), in §3 (`:78`, `replaySitzung`, `PGR-E4006`) und als Risiko in §6 (`:114`). Plan §1 (`:42`) und §6 (`:117`) dieses Slice zeigen dorthin, und das Risiko trägt den Ausgang „eingetreten“ mit Kennung. Offen sind zwei Reste: Dem Kopf fehlt die Spec-Stelle (F-343), und ein Spiegelstrich ist nicht nachgezogen (F-344). |
| F-331 | behoben, mit Rest | Die Diagnose nennt jetzt das SQL der Anweisungen (LH-FA-10.a). Eine Adresse für eine offene Prüfung braucht es deshalb nicht mehr. Ausschluss (`:37`), Rückführung (`:82`) und „Bereits geliefert“ (`:40`) von `slice-replay-semantik-mismatch` sind nachgezogen, die Welle-Zeile ebenfalls. Der Titel in Zeile 1 und der Bezug passen nicht zur Welle-Zeile (F-344). Bei der neuen Herleitung fehlt die Lebensdauer der Protokollobjekte (F-339). |
| F-332 | behoben | Der Funktionskommentar (`server.go:176`–`:179`) sagt jetzt drei Dinge: Das Warten ist nicht begrenzt, die Frist ist nicht verdrahtet, und wer `conn` schließt, merkt `PGR-E4006` selbst. Das trifft zu: `verbindungsende` erkennt `net.ErrClosed`, die Session endet also regulär. |
| F-333 | behoben | Laut Port-Kommentar ändert `Shutdown` keinen Zustand, darf mehrfach gerufen werden, und das Verbot einer neuen Interaktion hält der Aufrufer ein. Das entspricht `ReplayService.Shutdown`. |
| F-334 | behoben | Der Kopf nennt `LH-FA-13.b` und `SPEC-025` (`slice-extended-query-replay.md:16`). |
| F-335 | behoben, mit Rest | Der Test vergleicht jetzt die ganze Sicht des Clients mit der beim Aufzeichnen (`:259`), K2 grün. NULL und leerer Text sind in dieser Sicht nicht unterscheidbar (F-342). |
| F-336 | keine Nacharbeit erwartet | Die Message von `b697462` nennt keine `SPEC-*`- oder `ARC-*`-Kennung. Das Muster ist nicht wiedergekehrt. |
| F-337 | behoben | `TestReplayHerunterfahrenSpaeteFrist` erreicht die Lage, in der die Ordnung trägt; M1 ist rot (20 von 20). Plan §6 (`:119`) nennt den Beleg. Der Test hängt an einer Zeitannahme (F-341). |
| F-338 | aufgenommen | DoD-Punkt 1 nennt Herunterfahren und `PGR-E3003`, DoD-Punkt 2 nennt die SQL-Texte und das Verbindungsende. Ob das trägt, entscheidet der Verifier (F-346). |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-339 | MEDIUM | Die Herleitung sucht das letzte Parse beziehungsweise Bind gleichen Namens in allen aufgezeichneten Client-Nachrichten der Session vor dem Cursor, ohne Rücksicht auf die Lebensdauer des Objekts. Sie übersieht damit drei Fälle: `Close` (P1), die einfache Anfrage, die das unbenannte Statement zerstört (P2), und das Ende der Transaktion, mit dem das Portal endet (P3). Die Diagnose nennt dann als „empfangen“ das SQL eines Objekts, das es auf dem Draht nicht mehr gibt. Das wäre nach dem eigenen Schema „unbekannt“. In P3 lautet sie `Anweisung erwartet "SELECT 1", empfangen "SELECT 1"` und legt damit nahe, der Client meine dieselbe Anweisung. LH-FA-10.a verlangt die „tatsächlich empfangene Query“. Eine Abweichung wird weiter richtig erkannt (Klasse 5); falsch ist nur der Inhalt der Diagnose. Parameterwerte erscheinen in keinem Fall: `anweisung` liest nur `SQL`, `Statement`, `Portal`, `Name` und `Target`. | LH-FA-10.a; LH-FA-18.a (Mismatch); Reviewer-Skill MEDIUM „unklare Fehlerbehandlung am Rand des Spec-Bereichs“ | `internal/hexagon/services/replay.go:181`–`:188`, `:207`–`:220`, `:241`–`:258` | ja (Sonden P1 bis P3) | Diagnose leitet den Bezug ohne Lebensdauer der Protokollobjekte her |
| F-340 | LOW | Laut Abdeckungszeile nennt die Diagnose das SQL „in bind, execute, describe oder close“. `TestReplayExtendedDiagnoseAnweisung` hat keinen Fall mit `close`, keinen mit dem Ziel Portal bei describe oder close, keinen mit einem nach dem Bind überschriebenen unbenannten Statement und keinen an der Cursorgrenze. D1 (Statement zum falschen Zeitpunkt) und D2 (Nachricht am Cursor zählt mit) bleiben grün. Der Code ist in beiden Fällen richtig, aber ungeschützt. | `AGENTS.md` §3.7; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `docs/user/abdeckung-unit.md:34`; `internal/hexagon/services/replay_extended_test.go:328`–`:353` | ja (D1, D2) | Zusage in der Abdeckungszeile weiter als die Prüfung |
| F-341 | LOW | `TestReplayHerunterfahrenSpaeteFrist` ruft `cancel()`, sobald der Use Case das Parse gesehen hat (`:746`–`:747`). Der Test setzt dabei voraus, dass die Sitzung schon wieder in `Receive` steht. Prüft sie `ctx.Err()` erst nach `cancel()`, wartet sie auf `geweckt`, und damit auf die Freigabe. Die gibt der Test aber erst frei, wenn die zweite Nachricht angekommen ist. Der Test bricht dann nach 2 s ab, obwohl der Code richtig ist. T1 belegt das, eine Verzögerung ohne Verhaltensänderung: rot, 5 von 5. Unverändert war er in 600 Läufen grün; das Fenster ist klein, aber nicht geschlossen. Dazu kommt: Beim Schreiben über `sende` hängt der Test nicht, weil `Flush` in einer Goroutine läuft. | Maintainability; LH-FA-13.a | `internal/adapters/driving/pgwire/server_extended_test.go:723`–`:760`, besonders `:746`–`:749` | ja (T1) | Test hängt an einer Zeitannahme, die er nicht herstellt |
| F-342 | LOW | Laut Kommentar und Abdeckungszeile gleicht die Sicht des Clients nach SIGTERM der beim Aufzeichnen. Die Sicht schreibt die Zeilen mit `%q`, und darin sind NULL und leerer Text gleich (F1: beide `[[""]]`). Die zweite Ausführung liefert NULL. Ein Replay, das stattdessen `''` liefert, bliebe in diesem Test und in `TestE2EReplayExtendedPipeline` (`pipelineAblauf`, `:161`) grün. | `AGENTS.md` §3.7; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `test/integration/extended_replay_e2e_test.go:207`–`:210`, `:254`, `:161`; `docs/user/abdeckung-e2e.md:31` | ja (F1) | Zusage im Test-Kommentar weiter als die Prüfung |
| F-343 | MEDIUM | `slice-v1-abschluss-betrieb` bindet seit dem Fix die Frist `--shutdown-timeout` für beide Modi in §1 und DoD-Punkt 1. Sein Kopf nennt `SPEC-046` nicht, obwohl §1 (`:42`) die Stelle selbst als spezifiziert anführt. `AGENTS.md` §3.9 verlangt, dass der Kopf betroffener Folge-Slices im selben Commit nachgezogen wird. Die Klasse tritt zum dritten Mal auf, nach F-325 und F-334, die beide LOW waren. Daher MEDIUM nach der Skill-Regel. | Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice („Der Kopf nennt die berührten Spec-Stellen“); `AGENTS.md` §3.9; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“ | `docs/plan/planning/open/slice-v1-abschluss-betrieb.md:16`, `:42`, `:57` | ja (Lesen) | Kopf nennt berührte Spec-Stelle nicht |
| F-344 | LOW | Die Folge-Slices sind an drei Stellen nicht nachgezogen. (1) `slice-replay-semantik-mismatch` heißt in Zeile 1 weiter „Mismatch-Diagnose und Exit-Codes“ und nennt im Bezug [LH-FA-10](../../spec/lastenheft.md). Die Welle-Zeile führt einen anderen Titel mit Klammerzusatz und einen Bezug ohne LH-FA-10. Beide Fassungen bezeichnen denselben Slice verschieden. (2) In `slice-v1-abschluss-betrieb` sagt §1 `:38` weiter „Eigener Meldungscode der Klasse 4 (Vergabe nach `SPEC-034` …)“. Der Spiegelstrich darüber (`:37`), Zeile `:42` und DoD-Punkt 1 nennen dagegen den schon vergebenen `PGR-E4006`. Dritter Lauf in Folge mit dieser Klasse (F-323, F-331). | `AGENTS.md` §3.9; `BEO-REPO/plan-folgt-korrektur-nicht` | `docs/plan/planning/open/slice-replay-semantik-mismatch.md:1`, `:14`; `docs/plan/planning/welle-replay-semantik.md:52`; `docs/plan/planning/open/slice-v1-abschluss-betrieb.md:37`, `:38` | ja (Lesen) | Folge-Slice nach Teillieferung nur teilweise nachgezogen |
| F-345 | INFO | Für den Planner, zur Größe von `slice-v1-abschluss-betrieb`. Formal bleibt es bei drei Liefer-Punkten. DoD-Punkt 1 bündelt aber atomares Schreiben, den Exit-Code aller Klassen und jetzt die Frist in zwei Modi samt `PGR-E4006`. §3 führt vier Komponenten: `cli` und `pgwire` (treibend), `recording` (getrieben) und `bootstrap`. §1 `:42` verschiebt die Schnittfrage auf den Start des Slice. Das Wachstum ist damit benannt (Modul 5: „macht Wachstum benennbar, nicht unmöglich“). Ob der Slice noch in eine Review-Sitzung passt, ist eine Entscheidung des Planners. | Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice | `docs/plan/planning/open/slice-v1-abschluss-betrieb.md:42`, `:57`, `:75`–`:80` | — | — (Verweis an Planner) |
| F-346 | INFO | Für den Verifier: F-338 ist dadurch gelöst, dass die neuen Zusagen in die bestehenden DoD-Punkte eingefügt wurden. DoD-Punkt 1 trägt jetzt drei Aussagen: Abnahmeszenario 7, Herunterfahren nach LH-FA-13 und `PGR-E3003`. Ob jede davon belegt ist, prüft der Verifier (Modul 11). | Reviewer-Skill §Was dieser Skill NICHT macht | `docs/plan/planning/in-progress/slice-extended-query-replay.md:52`–`:53` | — | — (Verweis an Verifier) |

## Antwort auf die Prüffragen

1. **F-330 bis F-338.** Siehe Statustabelle. F-330, F-332, F-333, F-334 und F-337 sind behoben. F-331 und F-335 sind behoben, mit Rest (F-344, F-339 und F-342). F-336 braucht keine Nacharbeit, und F-338 liegt beim Verifier (F-346).
2. **SQL-Herleitung gegen LH-FA-10.a und LH-FA-18.a.**
   - **Was trägt:** Die Auflösung über das Portal sucht das Statement zum Zeitpunkt des Bind (`statementSQL(vorher, i, …)`). Ein später überschriebenes unbenanntes Statement wird deshalb richtig zugeordnet, aber nicht geprüft (D1, F-340). „keine“ steht für Flush und Sync, „unbekannt“ für einen Namen ohne Parse beziehungsweise Bind. Aufgerufen wird die Herleitung nur im Zweig `default`, wenn der Cursor schon geprüft innerhalb einer Extended-Interaktion steht; `Interactions[:c.pos+1]` kann daher nicht überlaufen. Für die empfangene Nachricht ist die Vorgeschichte der Aufzeichnung auch die des Clients, weil bis zum Cursor jede Nachricht strikt gleich war.
   - **Was nicht trägt:** Die Lebensdauer der Objekte wird nicht beachtet, also weder `Close` noch das Zerstören durch eine einfache Anfrage noch das Ende der Transaktion (F-339).
   - **Parameterwerte:** erscheinen in keinem Fall. Der bestehende Test mit `geheim-wert` prüft die ganze Meldung wörtlich.
3. **`spaeteFrist` und Determinismus.** Der Wrapper hält nur das Setzen einer Frist an, die nicht null ist. Das tut in `replaySitzung` allein der Wächter; das Zurücksetzen und `SetDeadline` laufen durch. M1 ist zuverlässig rot. Deterministisch ist der Test nicht: Er setzt voraus, dass die Sitzung vor `cancel()` wieder in `Receive` steht (F-341, T1).
4. **`TestE2EReplayExtendedSigtermMittenInFolge`.** Der Test vergleicht die volle Sicht mit derselben Formatierung wie `pipelineAblauf`, also Prepare, Zeilen, Befehls-Tag und Fehler je Ausführung (K2 grün). Für „wie aufgezeichnet“ trägt das bis auf NULL gegen leeren Text (F-342). Weicht die Formatierung der beiden Stellen später auseinander, wird der Test rot statt still grün.
5. **Zuschnitt der Folge-Slices und der Welle.** `slice-replay-semantik-mismatch` hat einen Liefer-Punkt, den Gegenstand `--fail-on-unconsumed`, Ausschluss und Rückführung passen dazu. Titel und Bezug weichen von der Welle-Zeile ab (F-344). `slice-v1-abschluss-betrieb` nimmt die Sendung jetzt an (F-330 behoben). Es bleiben formal drei Liefer-Punkte; das Wachstum steckt in DoD-Punkt 1 (F-345). Der Kopf nennt `SPEC-046` nicht (F-343), und ein Spiegelstrich in §1 ist veraltet (F-344).
6. **Bewertung, ohne zu entscheiden: Folge-Slice nur im Plan, nicht im Code-Kommentar.** Die Begründung trägt. Der Kommentar an `replaySitzung` trägt zwei der Klassen aus `AGENTS.md` §3.7:
   - **Grenze:** Das Warten ist nicht begrenzt, die Frist ist hier nicht verdrahtet.
   - **Kopplung:** Wer `conn` schließt, merkt `PGR-E4006` selbst.

   Beides steht im Indikativ und beschreibt den Ist-Zustand. Ein Slice-Name im Kommentar wäre ein Verweis auf künftige Arbeit, also zeitliche Schicht. Die gehört in `docs/plan/planning/`, und `AGENTS.md` §3.7 lässt als Herkunftsfeld nur die Rückwärts-Form `· seit …` zu. Der Name würde zudem veralten, sobald der Slice nach `done/` und ins Archiv wandert. Die Adresse ist trotzdem auffindbar, in Gegenrichtung: §3 von `slice-v1-abschluss-betrieb` nennt `replaySitzung` und die Pflicht, `PGR-E4006` zu merken. Wer die Frist verdrahtet, ändert also genau diese Funktion und ihren Kommentar. Ein Rest bleibt: Wer nur den Code liest, erfährt nicht, wo die Frist geplant ist. Er erfährt aber, dass sie fehlt und was ihr Anschluss leisten muss, und das ist die Zusage an den, der die Stelle ändert.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services` | geprüft. `anweisungen` läuft unter `mu`, nur im Mismatch-Zweig, ohne Zustandsänderung; Parameterwerte werden nicht gelesen. Befund F-339, F-340. |
| `internal/hexagon/ports/driving` | geprüft, ohne Befund. Der Port-Kommentar entspricht der Implementierung (F-333 behoben). |
| Core-Reinheit | geprüft, ohne Befund. `replay.go` importiert unverändert `context`, `fmt`, `sort`, `sync`, Domain und Ports. `make a-check` lief nicht. |
| `internal/adapters/driving/pgwire` | geprüft. Am Produktionscode ändert sich nur der Kommentar, und er beschreibt den Code richtig. Befund F-341 (Test). |
| `test/integration` | geprüft, K2 grün. Befund F-342. |
| `docs/user/` | geprüft. Die Abdeckungszeilen entsprechen den Deklarationen der Tests. Befund F-340, F-342. |
| `docs/plan/planning/` — Hard Rule 3.9, Größenregel | geprüft. Plan §1, §3 und §6 dieses Slice folgen dem Code. Befund F-343, F-344, F-345. |
| `spec/` — Strata, Hard Rule 3.4 | geprüft, ohne Befund. Keine Spec-Datei im Diff. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. `b697462` enthält keinen Move. |
| Hard Rules 3.5, 3.8 — ADRs | geprüft, ohne Befund. Keine ADR im Diff. |
| Hard Rule 3.6 — Gates nicht lockern | geprüft, ohne Befund. Keine Gate-Konfiguration im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft. Die neuen Kommentare stehen im Indikativ, ohne verworfene Alternative und ohne Satzrest. Befund F-340, F-342 (Zusagen in Test-Kommentaren und Abdeckungszeilen). |
| Commit-Message `b697462` | geprüft, ohne Befund. Sie nennt `slice-extended-query-replay`, [LH-FA-10](../../spec/lastenheft.md), [LH-FA-13](../../spec/lastenheft.md) und [LH-FA-18](../../spec/lastenheft.md), dazu keine Struktur-ID. |
| Register `BEO-REPO/*` | geprüft. F-343 und F-344 betreffen `BEO-REPO/plan-folgt-korrektur-nicht`, F-340 und F-342 `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:**

- Diagnose leitet den Bezug ohne Lebensdauer der Protokollobjekte her (F-339)
- Zusage in der Abdeckungszeile weiter als die Prüfung (F-340; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Test hängt an einer Zeitannahme, die er nicht herstellt (F-341)
- Zusage im Test-Kommentar weiter als die Prüfung (F-342; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, wieder nach F-335)
- Kopf nennt berührte Spec-Stelle nicht (F-343; drittes Mal nach F-325 und F-334)
- Folge-Slice nach Teillieferung nur teilweise nachgezogen (F-344; Register `BEO-REPO/plan-folgt-korrektur-nicht`, dritter Lauf in Folge nach F-323 und F-331)

## Verdikt

**Merge-blockierend nach dem Skill:** nein; kein HIGH. Die Findings des Folge-Reviews sind abgearbeitet, Reste stehen oben. Neu ist der Inhalt der Diagnose an der Lebensdauer der Protokollobjekte (F-339). Die Erkennung der Abweichung und ihre Klasse sind davon nicht berührt.

**Übergabe:**

- F-339 bis F-342 an den Implementer.
- F-343 und F-344 an Planner und Implementer (Kopf und Text der Folge-Slices).
- F-345 an den Planner, F-346 an den Verifier.
- **Pflege-Schritt des Skills fällig:** Zwei Klassen treten zum dritten Mal auf, „Kopf nennt berührte Spec-Stelle nicht“ (F-325, F-334, F-343) und „Folge-Slice nach Teillieferung nur teilweise nachgezogen“ (F-323, F-331, F-344). Nach Reviewer-Skill §Pflege ist zu prüfen, ob die Kategorie stimmt, ob eine Regel in `AGENTS.md` oder eine ADR das verhindert hätte und ob ein Sensor es prüfen kann. `AGENTS.md` §3.9 nennt Kopf und Folge-Slices bereits; die Wiederholung zeigt, dass die Regel allein nicht greift. Den Konflikt-Pfad über den Architect braucht es nicht, weil der Implementer keiner Einstufung widersprochen hat.
- Widerspricht der Implementer einer MEDIUM-Einstufung, genügt eine Begründung im Plan oder in der Closure-Notiz.
- Die Finding-Klassen gehen in die Closure-Notiz §7. Die DoD-Konformität prüft der Verifier separat.
