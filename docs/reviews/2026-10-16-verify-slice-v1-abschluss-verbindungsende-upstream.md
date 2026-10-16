# Verifikation `slice-v1-abschluss-verbindungsende-upstream` (Schritt 41)

Rolle: Verifier (Modul 11). Diff `c5089444..1ed4c0cc`, Plan `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-upstream.md`, Review `docs/reviews/2026-10-16-review-slice-v1-abschluss-verbindungsende-upstream.md`. Kein Produktcode geändert.

**Ergebnis: bestanden, 0 blockierende Befunde.** Ein LOW (Plan folgt dem Port-Kommentar nicht), drei INFO.

## 1. DoD gegen den Diff

| DoD-Punkt | Befund | Beleg |
|---|---|---|
| 1 Adapter `Query`/`Send`/`Receive` | erfüllt | `upstream.go`: `lies` liefert `ende`, `Query` und `Receive` kehren mit Fehlerantwort und `connection_closed` zurück, vor `ReadBufferLen`; `Send` liefert bei `endeGelesen` nil; `beendetVerbindung` liest `V`, leer dann `S`, genau `FATAL`/`PANIC`; Merker wird bei jeder `ErrorResponse` neu gesetzt, `ReadyForQuery` leert ihn; `ERROR` plus Abbruch bleibt `PGR-E6001`, ohne Fehlerantwort `PGR-E4003` (`TestEnde*`, `TestSendNachDemEnde`, `TestFehlerantwortVorDemAbbruch`) |
| 2 Leerlauf | erfüllt | `TestEndeImLeerlauf`, `TestEndeSendefehlerImLeerlauf` (Fake), `TestE2ERecordEndeLeerlauf` mit `idle_session_timeout` (mit und ohne Anfrage danach) |
| 3 Szenario 18, Formvergleich, `SPEC-029` | erfüllt | `TestE2ERecordEndeTerminate` (YAML und SQLite, `version: 2` / `meta.version` 2, Exit 0, keine Meldung), `TestE2EOhnePostgresEnde` (Replay, Phase 2), `TestE2ERecordNetzwerkabbruch` (`PGR-E4003`, `version` 1), `gleicheForm` gegen `ende-fatal`, `ende-terminate`, `ende-lebend`; der Satz in `SPEC-029` Punkt 2 steht in `51193b48`, dem Commit des Tests (nicht davor); Numerus-Satz steht seit `179ab424`, nicht erneut angefasst |
| Mutationstabelle in §7 | erfüllt | Tabellen Unit (S1 bis S7, E1, E2, V1, V2, M1, Q1 bis Q5, R1 bis R5, N1, N2, X1, X2, D1 bis D7) und E2E (EM1 bis EM10, RM1, RM2) je Zusage mit rotem Test |
| `make gates` grün | durch Beleg erfüllt | §7: voller Lauf an `bf6487cb`, Exit 0; `git diff --name-only bf6487cb..1ed4c0cc` nennt nur den Slice-Plan (Positivliste), der Stop-Hook antwortet `approve`. Nicht wiederholt, der Code-Stand entspricht dem Nachweis |
| Review, Closure, Register, Risiko-Ausgänge, Paarungen | offen (Closure des Planners) | Review liegt vor (1 MEDIUM, 3 LOW, 3 INFO, F-1 behoben); die übrigen Punkte sind Sache der Closure |

Handbuch und `README.md` sind nicht berührt (`git diff --name-only`), wie §1 und `AGENTS.md` §3.11 verlangen.

## 2. Bewusstes Brechen (Stichprobe, selbst gefahren)

Unit in einer Kopie, Image `pgr-verify40:deps`, `--network=none`, ein Textersatz je Mutant, Original danach zurück (Kopie und Original byte-gleich); Ausgangslage `TestEnde|TestSendNach|TestSendLiest|TestFehler` grün.

| Zusage | Mutation | Ergebnis |
|---|---|---|
| `FATAL` ist ein Ende (S1) | `FATAL` gestrichen | rot (`TestSendNachDemEnde`, `TestEnde*`) |
| genau `FATAL`/`PANIC` (S5-Art) | `fatal` zusätzlich akzeptiert | rot (`TestEndeSchweregrad`, `TestSendNachDemEnde`) |
| leerer Schweregrad ist kein Ende (S7) | `""` zählt als Ende | rot (`TestEndeSchweregrad`, `TestSendNachDemEnde`) |
| letzte Fehlerantwort entscheidet, `ERROR` dann `FATAL` (E1) | Merker behält die erste | rot (`TestSendNachDemEnde`, `TestFehlerantwortVorDemAbbruch`) |
| Merker-Reset bei `ReadyForQuery` (M1) | Reset entfernt | rot (`TestFehlerantwortVorDemAbbruch`, `TestSendNachFehlerantwort`) |
| `Send` liest nach dem Ende nichts (X1: Lesen mit 30 ms Frist; X2: ohne Frist) | Lesen nach dem Sendefehler | beide rot, nur `TestSendLiestNichtNachDemEnde` |
| Netzwerkabbruch ist `PGR-E4003` (EM6, E2E, frischer Baum aus `git archive HEAD`, PostgreSQL 17) | Verbindungsende ohne Fehlerantwort als Ende gewertet | rot nur `TestE2ERecordNetzwerkabbruch` (`PGR-E6001`, "connection_closed nur unmittelbar nach einer error_response"), `TestE2ERecordEndeTerminate` bleibt grün |

Alle Mutanten wie in §7 beschrieben rot; die Tabelle ist in der Stichprobe nicht widerlegt. Nicht wiederholt: EM1 (SQLite-Version), Q/R/N/D-Reihen, RM1/RM2 (Phase 2).

## 3. Äquivalenz-Aussagen in §7

`grep` nach "äquivalent" im Plan: kein Treffer außer der Überschrift "Äquivalente". Die drei Aussagen der Rubrik tragen Invariante und Ort: EM12 (Proxy schließt unmittelbar nach dem Schreiben, `kappProxy` im E2E-Test, Kommentar nennt es als Grenze), "Leerlauf ohne Anfrage" (`lies` nur aus `Query`/`Receive`, Ort `services/record.go`), `PANIC` (nicht herstellbar, akzeptiertes Negativ aus §6). Die frühere Aussage zu `Send` ist gestrichen und durch X1/X2 ersetzt, selbst rot gesehen. Kein Befund. Der Zähler `BEO-REPO/gruener-mutant-faelschlich-aequivalent` (F-1) ist Sache der Closure.

## 4. Kommentare (`AGENTS.md` §3.7, §3.11)

Der Port-Kommentar sagt zu, was Tests prüfen (`TestEndeSchweregrad`, `TestEndeLiestNichtWeiter`, `TestSendNachDemEnde`, `TestSendLiestNichtNachDemEnde`) und nennt "nicht zugesagt" für `Query`/`Receive` nach `connection_closed`; "nur noch `Close`" steht dort nicht mehr. Der Satz steht noch in der Spezifikation (V-2).

## 5. Messung

`git diff --numstat c5089444..1ed4c0cc -- internal/adapters/driven/postgres internal/hexagon/ports/driven test/integration docs/user`: 884 Einfügungen, 57 Löschungen, **941**, nachgerechnet (Aufschlüsselung 65 + 14 + 469 + 370 + 23 = 941). Über der Meldeschwelle 900 (gemeldet, Entscheidung des Orchestrators: kein Schnitt), unter der Marke 1200, in der Spanne des Architect (780 bis 1050), über der Mitte 910. Die Rückführung `in-progress` nach `next` ist nicht eingetreten.

## 6. Sensoren

`make abdeckung-check` Exit 0; `make gates-schnell` Exit 0 (schreibt den Stand); Stop-Hook `approve`; `make docs-check` vor dem Commit (siehe Commit).

## Findings

| # | Schwere | Befund | Zuständigkeit |
|---|---|---|---|
| V-1 | LOW | Plan §6, Zeile "Aufruf von `Query`, `Send` oder `Receive` nach `connection_closed`" sagt "der Use Case ruft nur `Close` (der Port-Kommentar sagt es)". Der Port-Kommentar sagt es seit `cf06cfdd` nicht mehr (Review F-2: `ClientMessage` ruft `Send` neben dem lesenden `Receive`); der Plan folgt der Korrektur nicht (`AGENTS.md` §3.9, `BEO-REPO/plan-folgt-korrektur-nicht`). Der Vermerk "akzeptiertes Negativ: kein Aufrufer" stimmt für `Send` nicht. | Implementer oder Planner, vor der Closure: Zeile auf den Ist-Stand ziehen ("nicht zugesagt für `Query` und `Receive`; `Send` zugesagt") |
| V-2 | INFO | `spec/spezifikation.md`, `LH-FA-25.a` *Senden nach dem Ende*: "Nach dem Ende ruft der Use Case den Upstream nur noch zum Schließen auf" hat keinen Test und widerspricht F-2 (`Send` kann nach dem Ende laufen). Gleicher Satz, den der Implementer aus dem Port-Kommentar gestrichen hat. Kein Code. | Architect (Spezifikation enger fassen oder den Test liefern, `AGENTS.md` §3.11) |
| V-3 | INFO | Review F-7: Der Formvergleich gegen die Fixtures läuft direkt nur für YAML; die SQLite-Fehlerantwort prüft `V`, `S`, `C`, `M` erst im Replay der Phase 2 (`endeSicht`). DoD 3 ist damit gedeckt (Wortlaut: "aufgezeichnete Interaktion", Replay liefert dieselben Antworten); ich habe Phase 2 nicht selbst gefahren, sie stützt sich auf den vollen Lauf in §7. | keine Aktion; Planner nimmt es zur Kenntnis |
| V-4 | INFO | Review F-5/F-6: Schweregrad-Regel an drei Stellen, SQLite-Versionsprobe liest Rohbytes. Mit Grenze in §7 benannt. | Planner (Beobachtung), keine Aktion im Slice |

## Vorschlag Risiko-Ausgänge (nur vorbereitet, Closure schreibt der Planner)

- Risiko 1 (Verhalten nur gegen 17): *weiter offen* nach `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (bereits so in §6).
- Risiko 2 (Weg für den Netzwerkabbruch): *entfallen*, Mittel entschieden und geliefert (Proxy im Testprozess), `TestE2ERecordNetzwerkabbruch` mit EM6 rot gesehen, hier erneut.
- Risiko 3 (Diff an der Grenze): *entfallen*, 941 unter der Marke 1200, eine Review-Sitzung hat ihn getragen; für `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` ist 941 gegen die Spanne bis 1050 keine Unterschätzung (kein Zählerzuwachs aus diesem Slice).
- Register: F-1 erhöht `BEO-REPO/gruener-mutant-faelschlich-aequivalent` auf 5×; F-4 des Reviews `BEO-REPO/plan-folgt-korrektur-nicht`; V-1 ist ein weiterer Fall dieser Klasse, vom Planner zu zählen.
