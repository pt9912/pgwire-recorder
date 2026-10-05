# Folge-Review: slice-extended-query-replay — 2026-10-05

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** Commit `c8ebf08` (`git diff 8b80a17..c8ebf08`), der die Findings F-321 bis F-329 aus dem [Review zu `b895565`/`09405f0`](2026-10-05-review-slice-extended-query-replay.md) abarbeiten soll. Schwerpunkte laut Auftrag: Status je Finding; Nebenläufigkeit beim Herunterfahren (Wächter, Lesefrist, Kanal `geweckt`); die neue Port-Methode `Shutdown` gegen `ARC-003` und [ADR-0001](../plan/adr/0001-hexagonale-architektur.md); die Startprüfung mit `Validate`/`PGR-E3003` gegen LH-FA-13.b; Kommentare gegen Prüfungen; Zuschnitt von `slice-replay-semantik-mismatch` nach dem Fix (`AGENTS.md` §3.9) und die Adresse der Replay-Frist in `slice-v1-abschluss-betrieb`.

**Entscheidung zu F-321 (Eingang, nicht Gegenstand dieses Reviews):** Der Nutzer hat entschieden, dass der Code der Spezifikation folgt. Beim Herunterfahren wird eine begonnene Extended-Interaktion bis zu ihrem `Sync` beantwortet. Spezifikation und Handbuch bleiben unverändert. Die Frist `--shutdown-timeout` wird in `slice-v1-abschluss-betrieb` verdrahtet.

**Ablage:** Eigene Datei statt Anhang an den Vorlauf-Report. Der Reviewer-Skill (§Output-Schema) verlangt „ein Report pro Lauf, Folgeläufe als neue Datei statt Überschreibung“, und auch beim Vorgänger-Slice steht das Folge-Review in einer eigenen Datei ([Folge-Review `slice-extended-query-record`](2026-10-05-folge-review-slice-extended-query-record.md)). Die Nummerierung läuft weiter.

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-extended-query-replay.md` (Kopf, §1 bis §3, §6), `docs/plan/planning/open/slice-replay-semantik-mismatch.md` (ganz), `docs/plan/planning/open/slice-v1-abschluss-betrieb.md` (ganz), `docs/plan/planning/welle-replay-semantik.md` (Slice-Tabelle)
- `spec/lastenheft.md` [LH-FA-03](../../spec/lastenheft.md), [LH-FA-10](../../spec/lastenheft.md), [LH-FA-13](../../spec/lastenheft.md), [LH-FA-18](../../spec/lastenheft.md)
- `spec/spezifikation.md`: LH-FA-13.a (Herunterfahren, Frist), LH-FA-10.a, LH-FA-18.a „Mismatch“, Optionstabelle (`--shutdown-timeout` für `record` und `replay`), `SPEC-025`, `SPEC-046`, `PGR-E3003`, `PGR-E4006`
- `spec/architecture.md` §1, §2 (`ARC-002`, `ARC-003`, `ARC-006`, Schichtregeln)
- [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0007](../plan/adr/0007-strict-replay.md)
- `AGENTS.md` Hard Rules 3.3 bis 3.9 und §5 (Dokumentations-Regeln); Register `BEO-REPO/plan-folgt-korrektur-nicht`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`
- Quelltext von `pgproto3` v5.11.0 (`backend.go`, `chunkreader.go`) im Modul-Cache der Stufe `deps`: Ein `Receive`, das an einer Lesefrist mitten in einer Nachricht abbricht, setzt beim nächsten Aufruf fort (`partialMsg`).
- Vorherige Reports am Modul: [Review](2026-10-05-review-slice-extended-query-replay.md) (F-321 bis F-329), [Review](2026-10-05-review-slice-extended-query-record.md) und [Folge-Review](2026-10-05-folge-review-slice-extended-query-record.md) zu `slice-extended-query-record`

**Ausgeführte Läufe im Repo:** `make docs-check` nach dem Anlegen dieser Datei. `make gates` lief nicht. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert (`git status --short` leer).

**Gegenproben und Mutationen:** Nur in Kopien (`git archive c8ebf08`) im Scratchpad, in einem Container aus der Stufe `deps` des `Dockerfile` (Tag `pgr-rev-fr-deps`, ohne Netz). Unit-Tests der Pakete `internal/hexagon/services` und `internal/adapters/driving/pgwire`; Integrationstests liefen nicht. Image und Cache-Volume sind gelöscht.

| # | Lauf oder Mutation | Ergebnis |
|---|---|---|
| K1 | Kopie ohne Änderung: `gofmt -l .`, `go vet -tags integration ./...`, Unit-Tests beider Pakete | grün, `gofmt` leer |
| A7 | wie im Vorlauf: `replaySitzung` liest nach einem Fehler von `ClientMessage` weiter (`continue` statt `return`) | rot (`TestReplayExtendedAbweichungImAdapter`, Frist des Clients statt Schließen) — F-322 behoben |
| M1 | `<-geweckt` entfernt (Zurücksetzen der Lesefrist nicht mehr nach dem Setzen durch den Wächter geordnet) | **grün** (F-337) |
| M4 | Lesefrist durch den Wächter beendet die Sitzung (`return` statt `continue`) | rot (`TestReplayHerunterfahren/laufende_Interaktion`) |
| M5 | Startprüfung mit `Validate` abgeschaltet | rot (`TestReplayStartformExtended`, beide Fälle) |
| M6 | `ReplayService.Shutdown` prüft nur `c.nachricht == 0`, nicht die Gruppe | rot (`TestReplayExtendedHerunterfahren`, nach dem `Flush`) |
| M7 | `replaySitzung` fragt `Shutdown` nur einmal, nicht vor jedem weiteren Lesen | rot (`TestReplayHerunterfahren/laufende_Interaktion`, Verbindung nach dem `Sync` offen) |

---

## Status der Findings F-321 bis F-329

| ID | Stand | Beleg |
|---|---|---|
| F-321 | behoben, wie entschieden | `replaySitzung` fragt nach dem Ende von ctx vor jedem weiteren Lesen `Replayer.Shutdown`; läuft eine Extended-Interaktion, liest sie ohne Frist weiter bis zum `Sync`, danach endet sie ohne weiteres Lesen. `TestReplayHerunterfahren`, `TestReplayExtendedHerunterfahren`, `TestE2EReplayExtendedSigtermMittenInFolge` (E2E nicht gelaufen); M4, M6, M7 rot. Spezifikation und Handbuch unverändert, wie entschieden. Die Adresse der Frist in `slice-v1-abschluss-betrieb` nimmt die Sendung nicht an: F-330; die Kopplung zur Frist im Kommentar: F-332. |
| F-322 | behoben | `TestReplayExtendedAbweichungImAdapter` verlangt EOF statt irgendeines Fehlers (`geschlossen`) und prüft, dass nur eine Nachricht den Use Case erreichte; A7 jetzt rot. Abdeckungszeile angepasst. |
| F-323 | teilweise behoben | DoD-Punkt 2 des Plans nennt die Diagnose als hier geliefert; `slice-replay-semantik-mismatch` hat Ziel, DoD und §3 auf `--fail-on-unconsumed` geschnitten. Mit dem Schnitt verlor die offene Prüfung gegen LH-FA-10.a ihre Adresse, und Titel, Ausschluss, §4 und Welle-Tabelle tragen den alten Gegenstand weiter: F-331. |
| F-324 | behoben | Funktionskommentar beschreibt Wächter, `Shutdown` und Weiterlesen bis zum `Sync`; der Kommentar im Lesezweig sagt ausdrücklich, dass ein Verbindungsende auch mitten in einer Extended-Interaktion im Replay regulär ist. Zum letzten Satz des Funktionskommentars F-332. |
| F-325 | behoben | Kopf nennt `LH-FA-13.a` und `SPEC-033`, Bezug nennt [LH-FA-13](../../spec/lastenheft.md). Für die neue Startprüfung fehlt die Stelle wieder: F-334. |
| F-326 | behoben | `NewReplayService` prüft jede Interaktion mit `Validate` und meldet `PGR-E3003` (Klasse 3, Startfehler nach LH-FA-13.b, `SPEC-025`); der Kommentar nennt die Kopplung an `Query` und `ClientMessage`. `TestReplayStartformExtended`, M5 rot. Die Prüfung steht damit im Core und hängt nicht mehr an einem bestimmten Leser (`ARC-008`). |
| F-327 | beim Planner, keine Aktion erwartet | INFO; im Fix nicht berührt, wie vorgesehen. |
| F-328 | beim Planner, keine Aktion erwartet | INFO; im Fix nicht berührt. |
| F-329 | als Risiko aufgenommen | Plan §6 nennt den ungeprüften Fall (große sofortige Ausgabe einer Flush-Gruppe bei gleichzeitig großer nächster Gruppe) mit Ausgang „offen bis Closure“. |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-330 | MEDIUM | Plan §1 weist „Frist `--shutdown-timeout` und `PGR-E4006` im Replay“ an `slice-v1-abschluss-betrieb`. Dieser Slice nimmt die Sendung nicht an, und `c8ebf08` hat ihn nicht berührt. Sein §1 übernimmt die Obergrenze nur „aus `slice-extended-query-record`“ und begründet sie mit „Heute wartet `record` ohne Frist“; sein Risiko in §6 spricht vom Verlust der Aufzeichnung; sein §3 nennt `internal/adapters/driving/pgwire` nicht, wo `replaySitzung` liegt; DoD-Punkt 1 verlangt nur, dass „Herunterfahren das Recording atomar schreibt“. Die Spezifikation verlangt die Frist für `record` und `replay` (Optionstabelle) samt `PGR-E4006`, wenn eine Interaktion unvollständig endet. Der Fix macht das dringlicher: Vorher endete `replay` auf `SIGTERM` sofort, jetzt hält ein Client, der mitten in einer Extended-Interaktion pausiert, Sitzung und Prozess ohne Grenze (`Serve` wartet auf alle Verbindungen). | `AGENTS.md` §3.9 („betroffene Folge-Slices im selben Commit“); Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice, Klasse 1 („Die Adresse muss die Sendung annehmen“); LH-FA-13.a, `SPEC-046`; `BEO-REPO/plan-folgt-korrektur-nicht`; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“ | `docs/plan/planning/in-progress/slice-extended-query-replay.md:42`; `docs/plan/planning/open/slice-v1-abschluss-betrieb.md:34`, `:77`, `:111`; `spec/spezifikation.md:459`–`:466`, `:647`; `internal/adapters/driving/pgwire/server.go:196`–`:205` | ja (Lesen) | Folge-Slice nimmt die zugewiesene Sendung nicht an |
| F-331 | MEDIUM | Beim Zuschnitt von `slice-replay-semantik-mismatch` auf `--fail-on-unconsumed` fiel der Satz „prüft die Diagnose beider Protokollvarianten gegen `LH-FA-10.a`“ weg, ohne neue Adresse. Das Ziel sagt jetzt „Erkennung und Diagnose der Abweichung sind für beide Protokollvarianten geliefert“. Die offene Frage aus dem Vorlauf (Antwort 7) bleibt aber stehen: LH-FA-10.a verlangt „erwartete Query“ und „tatsächlich empfangene Query“; die Diagnose einer Extended-Abweichung nennt SQL-Texte nur, wenn das Feld `sql` abweicht, bei einer Abweichung in `Bind` oder `Execute` keine. Ob das genügt, entscheidet niemand mehr. Zudem trägt der Slice den alten Gegenstand an vier Stellen weiter: Titel „Mismatch-Diagnose und Exit-Codes“, Ausschluss „hier genügt der Klartext der Diagnose“, Rückführung „die Diagnose verlangt eine Änderung am Recording-Format“ und dieselbe Zeile in der Slice-Tabelle der Welle mit [LH-FA-10](../../spec/lastenheft.md). Weiterer Fall der Klasse aus F-323, im selben Nachziehen. | LH-FA-10.a; `AGENTS.md` §3.9; Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice; `BEO-REPO/plan-folgt-korrektur-nicht`; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“ | `docs/plan/planning/open/slice-replay-semantik-mismatch.md:1`, `:32`, `:36`, `:81`; `docs/plan/planning/welle-replay-semantik.md:52`; `internal/hexagon/services/replay.go:141`–`:146`; `spec/spezifikation.md:376`–`:381` | ja (Lesen) | Folge-Slice nach Teillieferung nur teilweise nachgezogen |
| F-332 | LOW | Der Funktionskommentar an `replaySitzung` und Plan §1 legen fest, wie die Frist anschließt: „Wer das Warten begrenzt, schließt conn; das beendet das Lesen wie jedes Verbindungsende.“ Ein Verbindungsende behandelt derselbe Zweig als regulär, und `closeReplay` meldet nur `PGR-W2001`. Für das Ende durch die Frist mitten in einer Interaktion verlangt die Spezifikation `PGR-E4006` mit Klasse 4. Die Kopplung sagt dem, der die Frist verdrahtet, nicht, dass er den Fehler selbst merken muss. So wie beschrieben, endete der Lauf mit Exit-Code 0. | LH-FA-13.a (Frist, `PGR-E4006`), LH-FA-13.b; `AGENTS.md` §3.7 (Kopplung); `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `internal/adapters/driving/pgwire/server.go:176`–`:177`, `:212`–`:217`; `docs/plan/planning/in-progress/slice-extended-query-replay.md:34`, `:42` | ja (Lesen) | Kopplung nennt den Anschluss, nicht dessen Einstufung |
| F-333 | LOW | Der Port-Kommentar sagt: „Shutdown meldet den Beginn des Herunterfahrens“. `ReplayService.Shutdown` ändert keinen Zustand und ist eine reine Abfrage. Der Adapter ruft sie nach dem Ende von ctx vor jedem Lesen, also mehrfach. Am `Recorder` heißt dasselbe Wort „meldet“, dass der Service das Herunterfahren vormerkt und danach `ErrShutdown` liefert. Ein zweiter Adapter, der den Kommentar liest, ruft `Shutdown` womöglich nur einmal oder verlässt sich darauf, dass der Use Case danach keine neue Interaktion annimmt. Das tut am Replay nur der Adapter. | `ARC-003`; `AGENTS.md` §3.7; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `internal/hexagon/ports/driving/replay.go:28`–`:31`; `internal/hexagon/services/replay.go:193`–`:200`; `internal/hexagon/ports/driving/record.go:52`–`:57` | ja (Lesen) | Port-Kommentar sagt Zustandswechsel zu, den die Implementierung nicht hat |
| F-334 | LOW | Plan §1 liefert jetzt eine Startprüfung mit `PGR-E3003`. Der Kopf nennt unter „Berührte Spec-Stellen“ weder `LH-FA-13.b` noch `SPEC-025`, die diese Einstufung tragen. Dieselbe Klasse wie F-325 im selben Slice. | Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice („Der Kopf nennt die berührten Spec-Stellen“); `AGENTS.md` §3.9 (Kopf) | `docs/plan/planning/in-progress/slice-extended-query-replay.md:16`, `:34` | ja (Lesen) | Kopf nennt berührte Spec-Stelle nicht |
| F-335 | LOW | `TestE2EReplayExtendedSigtermMittenInFolge` sagt in Kommentar und Abdeckungstabelle zu: „beantwortet replay die Sync-Gruppe noch wie aufgezeichnet“. Geprüft wird je Ausführung nur, dass kein Fehler kommt und genau eine Zeile; die Werte (`a`, NULL) vergleicht der Test nicht. | `AGENTS.md` §3.7; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` | `test/integration/extended_replay_e2e_test.go:207`–`:209`, `:245`–`:250`; `docs/user/abdeckung-e2e.md:31` | ja (Lesen) | Zusage im Test-Kommentar weiter als die Prüfung |
| F-336 | LOW | Die Commit-Message nennt `SPEC-033` („F-325: Kopf nennt SPEC-033 und LH-FA-13.a“). `AGENTS.md` §5 Regel 1: Struktur-IDs (`SPEC-<NNN>`, `ARC-<NNN>`) adressieren innerhalb der Spec und gehören nicht in die Commit-Message. | `AGENTS.md` §5 Regel 1 | Commit `c8ebf08` (Message, Zeile 12 des Rumpfs) | ja (`git show -s c8ebf08`) | Struktur-ID in der Commit-Message |
| F-337 | INFO | Nebenläufigkeit beim Herunterfahren gelesen: Der Wächter setzt die Lesefrist und schließt danach `geweckt`; die Sitzung setzt die Frist erst nach `<-geweckt` zurück. Das Zurücksetzen ist damit gegen das Setzen geordnet, anders als in F-312 beim Record. `<-geweckt` kann nicht hängen, weil `fertig` erst beim Verlassen der Funktion schließt. Eine Lesefrist mitten in einer Nachricht verliert keine Bytes (`partialMsg`). Kein Test erreicht aber die Lage, in der die Ordnung trägt: ctx endet, während die Sitzung eine Nachricht verarbeitet, und `Shutdown` sagt nein. In allen drei Tests steht die Sitzung beim Signal in `Receive`. M1 (ohne `<-geweckt`) bleibt grün. | LH-FA-13.a; Maintainability | `internal/adapters/driving/pgwire/server.go:180`–`:205` | nein (zeitabhängig; M1 grün) | Ordnung einer Nebenläufigkeit ohne fangenden Test |
| F-338 | INFO | Für den Verifier: Der Fix liefert zwei neue Zusagen (Herunterfahren mitten in einer Extended-Interaktion, Startprüfung `PGR-E3003`), die in Plan §1 stehen, aber in keinem DoD-Punkt. Ob die DoD sie tragen muss, ist Verifier-Frage (Modul 11), nicht dieses Reviews. | Reviewer-Skill §Was dieser Skill NICHT macht | `docs/plan/planning/in-progress/slice-extended-query-replay.md:34`, `:52`–`:53` | — | — (Verweis an Verifier) |

## Antwort auf die Prüffragen

1. **F-321 bis F-329.** Siehe Statustabelle. F-321, F-322, F-324, F-325 und F-326 sind behoben, F-329 ist als Risiko aufgenommen, F-327 und F-328 liegen ohne erwartete Aktion beim Planner. F-323 ist nur teilweise behoben (F-331). F-321 folgt der Entscheidung des Nutzers; Spezifikation und Handbuch sind nicht im Diff.
2. **Nebenläufigkeit beim Herunterfahren.** Kein Wettlauf gefunden (F-337). Wächter, Lesefrist und `geweckt` sind richtig geordnet. `Shutdown` hält `ReplayService.mu` nur kurz und ruft nichts Blockierendes. Nach dem `Sync` gibt `Shutdown` das Ende frei, und die Sitzung liest nicht weiter, auch wenn der Client schon die nächste Interaktion gesendet hat (M7 rot). Eine einfache Anfrage gibt immer frei, weil der Cursor bei ihr auf Gruppe 0, Nachricht 0 steht. Offen bleibt die Obergrenze: Ein Client, der mitten in einer Extended-Interaktion pausiert, hält den Prozess ohne Frist, und die Adresse dafür trägt nicht (F-330, F-332).
3. **Port `Shutdown` gegen `ARC-003` und [ADR-0001](../plan/adr/0001-hexagonale-architektur.md).** Konform: Die Signatur trägt nur `context.Context`, `model.SessionID` und `bool`. Die Entscheidung, ob eine Interaktion läuft, liegt im Use Case, der Adapter übersetzt nur (`ARC-006`). Der Kommentar verspricht mehr, als die Implementierung tut (F-333).
4. **Startprüfung gegen LH-FA-13.b.** Konform: `PGR-E3003` ist Klasse 3 (`SPEC-025`) und entsteht vor dem Annehmen von Verbindungen. Die Prüfung läuft vor der Prüfung auf eine leere Aufzeichnung (`PGR-E3004`), und eine YAML-Aufzeichnung, die der Leser schon mit `Validate` angenommen hat, wird nicht neu abgelehnt. Der Kopf nennt die Stelle nicht (F-334).
5. **Kommentare gegen Prüfungen.** Der Kommentar an `TestReplayExtendedAbweichungImAdapter` und die neuen Unit-Kommentare sagen nur zu, was ihr Test prüft (A7, M4 bis M7 rot). Weiter als ihre Prüfung gehen der Port-Kommentar (F-333), die Kopplung zur Frist (F-332) und der E2E-Kommentar (F-335).
6. **Zuschnitt von `slice-replay-semantik-mismatch`.** Ziel, DoD und §3 sind auf `--fail-on-unconsumed` geschnitten, der Kopf ist angepasst. Die Prüfung gegen LH-FA-10.a hat keine Adresse mehr, und an vier Stellen steht der alte Gegenstand (F-331).
7. **Frist im Replay in `slice-v1-abschluss-betrieb`.** Die Adresse fehlt dort. Die Spezifikation verlangt die Frist für beide Modi, aber §1, §3, §6 und die DoD jenes Slice sprechen nur vom Record (F-330). Übernehmen müsste er für das Replay: die Frist ab dem ersten Signal, das Schließen der Client-Verbindung, `PGR-E4006` mit Klasse 4 statt eines regulären Endes (F-332), das zweite Signal und die Info-Zeile.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services` | geprüft. `Shutdown` unter `mu`, ohne blockierenden Aufruf; Startprüfung vor dem Befüllen von `frei`. Befund F-333 (Implementierung gegen Port-Kommentar). |
| `internal/hexagon/ports/driving` | geprüft. Nur Domain-Typen ([ADR-0001](../plan/adr/0001-hexagonale-architektur.md), `ARC-003`). Befund F-333. |
| Core-Reinheit | geprüft, ohne Befund. `replay.go` importiert unverändert `bytes`, `fmt`, Domain und Ports; kein `pgproto3`, kein Adapter, kein Dateisystem. `make a-check` lief nicht. |
| `internal/adapters/driving/pgwire` | geprüft. Der Wächter ist mit `fertig` und `geweckt` eingesammelt; `fakeReplayer.Shutdown` bildet die Regel des Service nach. Befund F-332, F-337. |
| `test/integration` | geprüft. Neuer E2E-Test mit Abdeckungs-Deklaration (nicht gelaufen). Befund F-335. |
| `docs/user/` | geprüft. Abdeckungstabellen entsprechen den neuen Deklarationen; das Handbuch ist nicht im Diff, wie entschieden. Befund F-335. |
| `docs/plan/planning/` — Hard Rule 3.9 | geprüft. Plan §1, §2, §3, §6 folgen dem Code. Befund F-330, F-331, F-334. |
| `spec/` — Strata, Hard Rule 3.4 | geprüft, ohne Befund. Keine Spec-Datei im Diff. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. `c8ebf08` enthält keinen Move. |
| Hard Rules 3.5, 3.8 — ADRs | geprüft, ohne Befund. Keine ADR im Diff. |
| Hard Rule 3.6 — Gates nicht lockern | geprüft, ohne Befund. Keine Gate-Konfiguration, kein `.a-check.yml`, kein `harness/mk/*` im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft. Neue Kommentare im Indikativ, ohne verworfene Alternative und ohne Satzrest. Befund F-332, F-333, F-335. |
| Commit-Message `c8ebf08` | geprüft. Nennt `slice-extended-query-replay`, [LH-FA-13](../../spec/lastenheft.md), [LH-FA-18](../../spec/lastenheft.md), [LH-FA-10](../../spec/lastenheft.md). Befund F-336. |
| Register `BEO-REPO/*` | geprüft. F-330 und F-331 betreffen `BEO-REPO/plan-folgt-korrektur-nicht` (verkörpert als `AGENTS.md` §3.9), F-332, F-333 und F-335 `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 5 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:**

- Folge-Slice nimmt die zugewiesene Sendung nicht an (F-330; Register `BEO-REPO/plan-folgt-korrektur-nicht`)
- Folge-Slice nach Teillieferung nur teilweise nachgezogen (F-331; Register `BEO-REPO/plan-folgt-korrektur-nicht`, zweiter Lauf in Folge nach F-323)
- Kopplung nennt den Anschluss, nicht dessen Einstufung (F-332; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Port-Kommentar sagt Zustandswechsel zu, den die Implementierung nicht hat (F-333; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Kopf nennt berührte Spec-Stelle nicht (F-334; zweites Mal im Slice nach F-325)
- Zusage im Test-Kommentar weiter als die Prüfung (F-335; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Struktur-ID in der Commit-Message (F-336)
- Ordnung einer Nebenläufigkeit ohne fangenden Test (F-337)

## Verdikt

**Merge-blockierend nach dem Skill:** nein; kein HIGH. Das Herunterfahren im Replay folgt jetzt LH-FA-13.a, wie der Nutzer es entschieden hat. Abweichung im Adapter und Startprüfung tragen Tests, die die Mutationen A7, M4 bis M7 rot färben.

**Übergabe:**

- F-330 an Planner und Implementer: `slice-v1-abschluss-betrieb` muss die Frist im Replay annehmen (§1, §3, §6, DoD) oder Plan §1 dieses Slice eine andere Adresse nennen. Bis dahin hängt `replay` auf `SIGTERM` ohne Grenze, sobald ein Client mitten in einer Extended-Interaktion pausiert.
- F-331 an Planner und Implementer (Zuschnitt und Adresse der Prüfung gegen LH-FA-10.a).
- F-332 bis F-335 an den Implementer; F-336 ohne Nacharbeit am Commit, als Hinweis für die nächsten Messages.
- F-337 ohne erwartete Aktion; F-338 an den Verifier.
- Widerspricht der Implementer einer MEDIUM-Einstufung, genügt eine Begründung im Plan oder in der Closure-Notiz. Den Konflikt-Pfad über den Architect braucht es erst ab HIGH oder beim dritten gleichen Konflikttyp. F-331 ist der zweite Lauf in Folge mit derselben Klasse; ein dritter wäre Anlass für den Pflege-Schritt des Skills.
- Die Finding-Klassen gehen in die Closure-Notiz §7. Die DoD-Konformität prüft der Verifier separat.
