# Review-Report: slice-extended-query-replay — 2026-10-05

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff d0e5d3b..09405f0`, also die Commits `b895565` (Extended Query im Replay) und `09405f0` (Diagnose der Parameter-Abweichung exakt geprüft) für `slice-extended-query-replay`. Schwerpunkte laut Auftrag: drei vom Implementer gemeldete Randformen, Mutations-Abdeckung der Zusagen aus LH-FA-18.a und [LH-FA-10](../../spec/lastenheft.md), Kommentare gegen Prüfungen, Plan gegen Code (`AGENTS.md` §3.9), Überschneidung mit `slice-replay-semantik-mismatch`.

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-extended-query-replay.md` (Kopf, §1 bis §6, §8), `docs/plan/planning/open/slice-replay-semantik-mismatch.md` (Kopf, §1, §2, §4)
- `spec/lastenheft.md` [LH-FA-09](../../spec/lastenheft.md), [LH-FA-10](../../spec/lastenheft.md), [LH-FA-13](../../spec/lastenheft.md), [LH-FA-18](../../spec/lastenheft.md)
- `spec/spezifikation.md`: LH-FA-02.b, LH-FA-09.a, LH-FA-10.a, LH-FA-13.a, LH-FA-13.b, LH-FA-18.a (Interaktion, Abbruch, Record, Gruppen, Replay, Fehler, Mismatch), `SPEC-018`, `SPEC-033`
- `spec/architecture.md` `ARC-002` bis `ARC-008` (Komponenten und Schichtregeln)
- [ADR-0007](../plan/adr/0007-strict-replay.md), [ADR-0012](../plan/adr/0012-extended-query-gruppen.md), [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md)
- `AGENTS.md` Hard Rules 3.3 bis 3.9; Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/spec-randform-erst-im-review-entschieden`, `BEO-REPO/plan-folgt-korrektur-nicht`
- `docs/user/benutzerhandbuch.md` §Mit einem Datenbanktreiber arbeiten
- Vorherige Reports am Modul: [Review](2026-10-05-review-slice-extended-query-record.md) und [Folge-Review](2026-10-05-folge-review-slice-extended-query-record.md) zu `slice-extended-query-record` (bis F-320)

**Ausgeführte Läufe im Repo:** `make docs-check` nach dem Anlegen dieser Datei. `make gates` lief nicht. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert (`git status --short` leer).

**Mutationen und Sonden:** Nur in Kopien (`git archive 09405f0`) im Scratchpad, in einem Container aus der Stufe `deps` des `Dockerfile` (Tag `pgr-rev-eqr-deps`, ohne Netz). Unit-Tests der Pakete `internal/hexagon/services` und `internal/adapters/driving/pgwire`; Integrationstests liefen nicht. Image, Cache-Volume und Kopien sind gelöscht.

| # | Lauf oder Mutation | Ergebnis |
|---|---|---|
| K1 | Kopie ohne Änderung: `gofmt -l .`, `go vet -tags integration ./...`, Unit-Tests beider Pakete | grün, `gofmt` leer |
| A1 | `ClientMessage` setzt `c.gruppe` nach der letzten Gruppe nicht auf 0 zurück | rot (`TestReplayExtendedWiederholung`, Index außerhalb) |
| A2 | `ClientMessage` gibt die Server-Nachrichten der Gruppe schon vor deren letzter Client-Nachricht frei | rot (fünf Tests über `sendeAlle`) |
| A3 | Matcher vergleicht `portal` nicht | rot (`TestReplayExtendedFelder`) |
| A4 | `gleicheWerte` ignoriert `Null` | rot (`TestReplayExtendedFelder`) |
| A5 | Matcher vergleicht `max_rows` nicht | rot (`TestReplayExtendedFelder`) |
| A6 | Diagnose der Abweichung hängt `m.Params` mit `%v` an | rot (`TestReplayExtendedAbweichung`, exakter Vergleich aus `09405f0`) |
| A7 | `replaySitzung` liest nach einem Fehler von `ClientMessage` weiter (`continue` statt `return` nach `s.fail`) | **grün**; `TestReplayExtendedAbweichungImAdapter` läuft 5 s statt unter 2 s (F-322) |
| A8 | `ClientMessage` prüft nicht, ob am Cursor eine Extended-Interaktion steht | rot (`TestReplayExtendedFalscheArt`, Index außerhalb) |
| S1 | Sonde im Paket `pgwire` mit `fakeReplayer`: `Parse` senden, dann `ctx` beenden, dann `Bind`, `Execute`, `Sync` senden | Client liest `unexpected EOF` ohne `ErrorResponse`; der Use Case sah eine Nachricht; `CloseConnection` lief (F-321) |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-321 | MEDIUM | Beim Herunterfahren mitten in einer Extended-Interaktion endet die Replay-Session vor dem `Sync`. Der Client erhält keine Fehlerantwort, die Verbindung wird geschlossen, und die Interaktion bleibt unverbraucht (`PGR-W2001`, Sonde S1). LH-FA-13.a sagt ohne Unterscheidung nach Modus, laufende Sessions endeten „nach Abschluss ihrer laufenden Interaktion“; diese nehme „bis zu ihrem `Sync` noch ihre Client-Nachrichten an“. LH-FA-18.a „Abbruch“ sagt dasselbe für „eine Session“. Das Benutzerhandbuch sagt es in einem Abschnitt, der `record` und `replay` gleichermaßen behandelt. Die Lesart in Plan §6 („ohne den Replay-Modus für Extended zu nennen“) trägt für den Satz aus LH-FA-13.a nicht, weil dieser keinen Modus ausnimmt. Für LH-FA-18.a „Abbruch“ ist sie vertretbar, weil der Absatz sonst mit Record-Begriffen spricht. §1 schließt die Randform als „anderer Vorgang“ aus, ohne Folge-Slice. Sie entsteht aber erst mit diesem Slice, weil `replaySitzung` Extended-Nachrichten vorher mit `PGR-E6001` ablehnte. „Nicht im Code entschieden“ (§6) stimmt nicht: Der Code hat ein bestimmtes Verhalten, und das widerspricht dem Handbuch. | LH-FA-13.a; LH-FA-18.a „Abbruch“; [LH-FA-13](../../spec/lastenheft.md); `BEO-REPO/spec-randform-erst-im-review-entschieden`; Reviewer-Skill MEDIUM „unklare Fehlerbehandlung am Rand des Spec-Bereichs“ | `internal/adapters/driving/pgwire/server.go:172`–`:191`; `spec/spezifikation.md:446`–`:454`, `:688`–`:690`; `docs/user/benutzerhandbuch.md:446`–`:449`; `docs/plan/planning/in-progress/slice-extended-query-replay.md:42`, `:113` | ja — Sonde S1 | Randform still im Code entschieden |
| F-322 | MEDIUM | `TestReplayExtendedAbweichungImAdapter` sagt in Kommentar und Abdeckungstabelle zu: „und die Verbindung endet“. Geprüft wird das mit `fe.Receive()` und `err == nil`. Läuft der Adapter nach der Abweichung weiter (A7), endet `Receive` an der 5-s-Frist des Clients mit einem Fehler, und der Test ist grün. Auch `TestE2EReplayExtendedAbweichung` fängt A7 nicht, weil pgx nach der `FATAL`-Antwort selbst schließt. Dass die Verbindung nach `PGR-E5001` endet (`SPEC-018`), und damit, dass LH-FA-10.a „springt nicht zur nächsten Interaktion“ nicht von der Disziplin des Clients abhängt, prüft für Extended kein Test. | `SPEC-018` („Verbindung endet“); LH-FA-10.a; `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Reviewer-Skill MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag“ | `internal/adapters/driving/pgwire/server_extended_test.go:356`–`:369`, `:318`; `internal/adapters/driving/pgwire/server.go:212`–`:216`; `docs/user/abdeckung-unit.md:35` | ja — A7 | Prüfung besteht durch Zeitüberschreitung statt durch das zugesagte Verhalten |
| F-323 | MEDIUM | Der Plan und der überschneidende Folge-Slice sagen Verschiedenes darüber, wer die Diagnose der Extended-Abweichung liefert. Im Plan dieses Slice sagen §1 („ihre Diagnose nennt Session, Interaktion, Gruppe …“) und der Ausschluss („dort als geliefert vermerkt“), sie sei hier geliefert. DoD-Punkt 2 sagt dagegen: „Diagnose und Klasse folgen `slice-replay-semantik-mismatch`“. In `slice-replay-semantik-mismatch` steht unter „Bereits geliefert“ dasselbe wie hier in §1. Sein Ziel sagt aber weiter: „die Diagnose und Klasse für beide Protokollvarianten liefert dieser Slice“. Sein erster DoD-Liefer-Punkt verlangt eine abweichende Extended-Nachricht „mit Gruppe und Nachricht in der Diagnose“, also genau das, was dieser Diff liefert. Hard Rule 3.9 verlangt, dass betroffene Folge-Slices im selben Commit nachgezogen werden. Nachgezogen ist nur ein Absatz. Das Muster stand schon viermal als LOW in Reports (F-291, F-297, F-307, F-315); nach Skill ist es deshalb MEDIUM. Weil `BEO-REPO/plan-folgt-korrektur-nicht` als `AGENTS.md` §3.9 verkörpert ist, betrifft das Wiederauftreten auch dessen Retirement-Check. | `AGENTS.md` §3.9; Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice, Klasse 1 („Die Adresse muss die Sendung annehmen“); `BEO-REPO/plan-folgt-korrektur-nicht`; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“ | `docs/plan/planning/in-progress/slice-extended-query-replay.md:34`, `:41`, `:53`; `docs/plan/planning/open/slice-replay-semantik-mismatch.md:32`, `:39`, `:49` | ja (Lesen) | Folge-Slice nach Teillieferung nur teilweise nachgezogen |
| F-324 | LOW | Die Kommentare an `replaySitzung` beschreiben nur den einfachen Fall. Seit Extended-Nachrichten den Use Case erreichen, deckt derselbe Zweig mehr ab. Der Funktionskommentar sagt „eine laufende Anfrage läuft zu Ende“. Eine laufende Extended-Interaktion läuft nicht zu Ende (F-321). Der Kommentar im Lesezweig sagt „Ende nach einem ReadyForQuery, mit oder ohne Terminate, ist regulär (LH-FA-02.b)“. Der Zweig behandelt aber auch ein Verbindungsende vor dem `Sync` einer laufenden Extended-Interaktion still als regulär. Das kann für das Replay so gewollt sein, denn LH-FA-02.b gehört zum Record; der Kommentar sagt es aber nicht. | `AGENTS.md` §3.7; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Maintainability | `internal/adapters/driving/pgwire/server.go:167`–`:170`, `:187`–`:188` | ja (Lesen) | Kommentar deckt den erweiterten Zweig nicht |
| F-325 | LOW | Der Kopf nennt unter „Berührte Spec-Stellen“ `SPEC-033` nicht. §1 („nie Parameterwerte“), der Kommentar an `ClientMessage` („Parameterwerte nennt sie nicht (SPEC-033)“) und der Test aus `09405f0` setzen genau diese Stelle um. LH-FA-13.a fehlt ebenfalls, obwohl §1 und §6 sie behandeln (F-321). | Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice („Der Kopf nennt die berührten Spec-Stellen“); `AGENTS.md` §3.9 (Kopf) | `docs/plan/planning/in-progress/slice-extended-query-replay.md:16`; `internal/hexagon/services/replay.go:113` | ja (Lesen) | Kopf nennt berührte Spec-Stelle nicht |
| F-326 | LOW | `ClientMessage` und `erwarteteNachricht` greifen ohne Längenprüfung auf `Groups[c.gruppe].Client[c.nachricht]` zu. Sie verlassen sich darauf, dass jede Extended-Interaktion mindestens eine Gruppe mit mindestens einer Client-Nachricht hat. Das hält heute `Interaction.Validate`, aufgerufen im YAML-Leser. Der Port sagt dazu nur „Load liest eine Aufzeichnung“, und `NewReplayService` prüft die Form nicht. Ein zweiter Leser ohne `Validate` (`ARC-008` nennt auch SQLite) ließe eine beschädigte Aufzeichnung im Replay-Prozess als Panic enden statt als Startfehler der Klasse 3. Die Kopplung steht an keiner der beiden Stellen. | `ARC-004`, `ARC-008`; LH-FA-13.b (Startfehler); Maintainability | `internal/hexagon/services/replay.go:141`, `:173`–`:177`; `internal/hexagon/ports/driven/recording.go:17`; `internal/adapters/driven/recording/yaml.go:560` | ja (Lesen; A1/A8 zeigen die Panic-Form) | Invariante des Lesers am Port ungenannt |
| F-327 | INFO | Randform 2 (falsche Protokollart): Im Replay ist sie `PGR-E5001`, im Record `PGR-E6001`. Das ist in sich stimmig und vom Lastenheft gedeckt. [LH-FA-18](../../spec/lastenheft.md) Negative verlangt für eine nicht passende Extended-Interaktion eine Abweichung „wie in LH-FA-10“. LH-FA-18.a „Replay“ verlangt für jede eingehende Client-Nachricht Gleichheit im Nachrichtentyp. Die `PGR-E6001`-Liste in LH-FA-18.a gilt ausdrücklich „im Record“. Eine `Query` mitten in einer Extended-Interaktion kann nie in einer Aufzeichnung stehen, weil der Record sie ablehnt. Im Replay weicht sie deshalb notwendig von der Aufzeichnung ab. Ausdrücklich nennt die Spezifikation den Replay-Fall nicht: Schritt 1 in LH-FA-09.a vergleicht nur SQL-Texte. Plan §1 leitet die Regel offen ab, statt sie still zu setzen; ein Fall für das Register ist sie deshalb nicht. Hinweis an den Planner, falls ein Satz in LH-FA-18.a gewünscht ist. | LH-FA-18.a „Replay“, „Record“; LH-FA-09.a; [LH-FA-18](../../spec/lastenheft.md) Negative | `internal/hexagon/services/replay.go:97`–`:100`, `:127`–`:130`; `docs/plan/planning/in-progress/slice-extended-query-replay.md:34` | ja (Lesen; A8, `TestReplayExtendedFalscheArt`) | — (kein Fehlermuster) |
| F-328 | INFO | Randform 3 (Parameterwerte auch bei `debug` nie in der Diagnose): Der Code erfüllt beide Stellen der Spezifikation. LH-FA-18.a „Mismatch“ verbietet Klartext-Parameterwerte, „wenn der Log-Level nicht `debug` ist“; das ist ein bedingtes Verbot, kein Gebot für `debug`. `SPEC-033` sagt ohne Bedingung „verborgen werden nur die Parameterwerte“. Die beiden Sätze sind aber verschieden gefasst. Ein späterer Slice, der LH-FA-18.a als Erlaubnis für `debug` liest, verletzte `SPEC-033`. Zudem geht die Diagnose als `ErrorResponse` an den Client, nicht nur ins Log; ein Log-Level als Schalter wirkte also auch auf den Client. Hinweis an den Spec-Verantwortlichen. | LH-FA-18.a „Mismatch“; `SPEC-033` | `spec/spezifikation.md:759`–`:762`, `:1413`–`:1414`; `internal/hexagon/services/replay.go:132`–`:139` | ja (Lesen; A6) | Zwei Spec-Sätze zu derselben Regel verschieden gefasst |
| F-329 | INFO | [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) nennt als Re-Evaluierungs-Trigger, dass der Replay-Pfad denselben Zwei-Richtungs-Ablauf braucht. Plan §1 begründet, warum er hier nicht ausgelöst ist: `replaySitzung` antwortet abwechselnd wie ein PostgreSQL-Server. Die Begründung trägt. Das Replay gibt eine Gruppe erst nach `Flush` oder `Sync` frei und verlangt vom Client damit nicht mehr als PostgreSQL. `TestE2EReplayExtendedGegendruck` prüft aber nur die Aufzeichnung, in der die verzögerte große Ausgabe in der Sync-Gruppe landet. Den Fall „große sofortige Ausgabe einer Flush-Gruppe, während der Client eine große nächste Gruppe schreibt“ prüft er nicht. | [ADR-0030](../plan/adr/0030-full-duplex-im-record-pfad.md) Re-Evaluierungs-Trigger; [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) | `docs/plan/planning/in-progress/slice-extended-query-replay.md:40`; `test/integration/extended_replay_e2e_test.go` (`TestE2EReplayExtendedGegendruck`) | nein (Integrationstests liefen nicht) | — (kein Fehlermuster) |

## Antwort auf die Schwerpunkte

1. **Randform 1, Herunterfahren mitten in einer Extended-Interaktion.** Die Randform ist echt, und der Code entscheidet sie bereits: Die Session endet vor dem `Sync`, ohne Fehlerantwort (S1). Gemessen am Wortlaut ist das die schwächere Lesart. LH-FA-13.a unterscheidet nicht nach Modus. Das Handbuch sagt für beide Modi zu, dass das Werkzeug auf das `Sync` wartet. Für den Ausschluss als „anderer Vorgang“ fehlt eine Adresse. Wer entscheidet: der Planner beziehungsweise der Spec-Verantwortliche, entweder als Satz in LH-FA-13.a/18.a, der das Replay ausnimmt, samt Handbuch, oder als Verhalten in `replaySitzung`. Das Review entscheidet das nicht (F-321).
2. **Randform 2, falsche Protokollart.** Tragfähig und vom Lastenheft gedeckt; die Asymmetrie zum Record ist gewollt, weil `PGR-E6001` dort eine Form ist, die keine Aufzeichnung trägt (F-327).
3. **Randform 3, Parameterwerte auch bei `debug` nicht.** Konform mit beiden Spec-Stellen. Offen ist nur, dass die Spezifikation dieselbe Regel an zwei Stellen verschieden fasst (F-328).
4. **Mutations-Abdeckung.** Die Zusagen aus LH-FA-18.a „Replay“ fangen die Unit-Tests: Feldvergleich je Feld, keine Freigabe vor dem Gruppenende, Cursor bleibt bei Abweichung stehen, Gruppen- und Interaktionswechsel, falsche Art, keine Parameterwerte in der Diagnose (A1 bis A6, A8). Ohne fangenden Test bleibt die Zusage „die Verbindung endet“ nach einer Abweichung (A7, F-322).
5. **Kommentare gegen Prüfungen.** Matcher- und Port-Kommentare sagen nur zu, was geprüft ist. Abweichend: der Test-Kommentar aus F-322 und die Kommentare an `replaySitzung` (F-324).
6. **Plan folgt Code (§3.9).** §1, §3 und §6 folgen dem Code; ausgenommen sind der Kopf (F-325), DoD-Punkt 2 und der Folge-Slice (F-323).
7. **Überschneidung mit `slice-replay-semantik-mismatch`.** Was dort noch bleibt: `--fail-on-unconsumed` und der Abgleich der Diagnose beider Varianten mit LH-FA-10.a. Ziel und erster DoD-Punkt sind noch nicht auf diesen Rest geschnitten (F-323). Für den Abgleich mit LH-FA-10.a ist offen, ob „erwartete Query“ bei einer Bind-Abweichung ohne `Parse` in derselben Gruppe erfüllt ist. Diese Prüfung liegt bei jenem Slice und beim Verifier, nicht hier.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services` | geprüft. Strict Matcher nach LH-FA-18.a und [ADR-0007](../plan/adr/0007-strict-replay.md) ohne Normalisierung; Freigabe der Server-Nachrichten nach dem Gruppenende nach [ADR-0012](../plan/adr/0012-extended-query-gruppen.md); Sperre über den ganzen Aufruf, ohne blockierenden Aufruf darunter. Befund F-326. |
| `internal/hexagon/ports/driving` | geprüft, ohne Befund. `ClientMessage` trägt nur Domain-Typen ([ADR-0001](../plan/adr/0001-hexagonale-architektur.md), `ARC-003`); der Kommentar sagt nur zu, was `TestReplayExtendedGruppen` und `TestReplayExtendedFalscheArt` prüfen. |
| Core-Reinheit | geprüft, ohne Befund. `matcher.go` und `replay.go` importieren `bytes`, `fmt`, Domain und Ports; kein `pgproto3`, kein Adapter. |
| `internal/adapters/driving/pgwire` | geprüft. `replaySitzung` übersetzt nur und trifft keine Replay-Entscheidung (`ARC-006`); Zielart außer `S`/`P` erreicht den Use Case nicht. Befund F-321, F-322, F-324. |
| `test/integration` | geprüft. Abdeckungs-Deklarationen sagen zu, was die Tests prüfen; Abnahmeszenario 7 in der Phase ohne PostgreSQL vergleicht die volle Sicht. Befund F-329. |
| `tools/test/run-integration-tests.sh` | geprüft, ohne Befund. Nur der Kopfkommentar, im Indikativ. |
| `docs/user/` | geprüft. Abdeckungstabellen entsprechen den Deklarationen. Das Handbuch ist nicht im Diff; seine Zusage zum Herunterfahren gilt jetzt auch für das Replay (F-321). |
| `README.md` | geprüft, ohne Befund. Der Satz zum Stand nennt die Wiedergabe beider Protokolle. |
| `docs/plan/planning/` — Hard Rule 3.9 | geprüft. Befund F-323, F-325. |
| `spec/` — Strata, Hard Rule 3.4 | geprüft, ohne Befund. Keine Spec-Datei im Diff. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. Beide Commits ohne Move; der Move `986035c` und die Verweise `73e0ccf` liegen davor in eigenen Commits. |
| Hard Rule 3.5, 3.8 — ADRs | geprüft, ohne Befund. Keine ADR im Diff. |
| Hard Rule 3.6 — Gates nicht lockern | geprüft, ohne Befund. Keine Gate-Konfiguration, kein `.a-check.yml`, kein `harness/mk/*` im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft. Neue Kommentare sind Zusagen oder Kopplungen im Indikativ; keine verworfene Alternative, kein Satzrest. Befund F-324. |
| Commit-Messages `b895565`, `09405f0` | geprüft, ohne Befund. Beide nennen `slice-extended-query-replay` und `LH-*`-Kennungen, `b895565` dazu [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) und [ADR-0007](../plan/adr/0007-strict-replay.md); keine `SPEC-*`- oder `ARC-*`-Kennung. |
| Register `BEO-REPO/*` | geprüft. F-321 ist ein weiterer Fall von `BEO-REPO/spec-randform-erst-im-review-entschieden` (Plan §8: bisher 2×). F-322 betrifft `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, F-323 `BEO-REPO/plan-folgt-korrektur-nicht` (verkörpert). |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:**

- Randform still im Code entschieden (F-321; Register `BEO-REPO/spec-randform-erst-im-review-entschieden`)
- Prüfung besteht durch Zeitüberschreitung statt durch das zugesagte Verhalten (F-322; Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Folge-Slice nach Teillieferung nur teilweise nachgezogen (F-323; Register `BEO-REPO/plan-folgt-korrektur-nicht`)
- Kommentar deckt den erweiterten Zweig nicht (F-324)
- Kopf nennt berührte Spec-Stelle nicht (F-325)
- Invariante des Lesers am Port ungenannt (F-326)
- Zwei Spec-Sätze zu derselben Regel verschieden gefasst (F-328)

## Verdikt

**Merge-blockierend:** nein. Kein HIGH. Der Replay-Pfad liefert für Extended Query keine geratene Antwort, hält die Reihenfolge und gibt Gruppen erst nach ihrem Ende frei. Diese Zusagen tragen Tests, die die Mutationen A1 bis A6 und A8 rot färben.

**Übergabe:**

- F-321 an den Planner beziehungsweise Spec-Verantwortlichen mit diesem Report als Artefakt: Der Gegenstand ist eine Spec-Entscheidung, keine Implementierungsfrage. Bis dahin gehört ihm in Plan §1 eine Adresse oder eine Begründung, die trägt.
- F-322, F-324, F-325, F-326 an den Implementer; F-323 an Implementer und Planner, weil `slice-replay-semantik-mismatch` mitbetroffen ist.
- F-327 und F-328 an den Planner, F-329 ohne erwartete Aktion.
- Widerspricht der Implementer einer MEDIUM-Einstufung, genügt eine Begründung im Plan oder in der Closure-Notiz. Den Konflikt-Pfad über den Architect braucht es erst ab HIGH oder beim dritten gleichen Konflikttyp.
- Die Finding-Klassen gehen in die Closure-Notiz §7. Die DoD-Konformität prüft der Verifier separat.
