# Verifikation slice-v1-abschluss-verbindungsende-replay-extended-flush

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-16. **Spanne:** `e982c7be4073c4e54be8c04a50125389af5ee660` (Basis der Marke) bis `498ace8` (HEAD): Umsetzung bis `245d245`, Belege `3b88f4a`, Review-Report `1adf9b0`, Nacharbeit `7dcbd12` (Architect) und `07faa02`, Gate-Beleg `498ace8`.

**Eigene Läufe statt Zusage.** Mutanten in Kopien (`go.mod`, `go.sum`, `internal/` im Scratchpad, Image `deps`, `--network=none`, `go test -count=1`), ein Lauf zur Zeit; ein E2E-Mutant im Arbeitsbaum mit `make test-integration`, danach `git checkout`; `git status` danach leer. Basis ohne Mutant: beide Pakete `ok`. Das Review wurde nicht neu gefahren.

**Verdikt: PASS mit Befunden.** Alle Liefer-Punkte der DoD sind bestätigt, die Nacharbeit ändert nur Kommentare, die Zusagen der `Flush`-Gruppe fängt je eine Mutation. Ein LOW (V-1: ein Kommentar im Port nennt noch „Sync“) und ein INFO (V-2). Keiner blockiert die Closure.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1: Replay-Service liefert die Gruppe mit Ende und `Flush`, Verbrauch, Abweichung `PGR-E5001`, `Terminate`/Abbruch unverbraucht, Randformen des Gebers, Fixture je Fall | bestätigt | sieben Tests in `replay_ende_extended_flush_test.go`; Produktcode unverändert (siehe Nacharbeit); Mutanten unten rot |
| 2: Binary (E2E, handgeschriebene Aufzeichnung) und Adapter-Test Lauf `flush` | bestätigt | `TestE2EReplayEndeExtendedFlush` mit eigener Fixture `testdata/ende-extended-flush.yaml` (eine Einbettung, nicht mit `…Sync` geteilt); Lauf `flush` als Unterfall von `TestReplayEndeGruppeDavorSchliesstNicht`; Zeilen je E2E-Fall in §7 (27 + 20 + 3 = 50) nachgezählt |
| 3: je Zusage eine Mutation, Tabelle in §7 | bestätigt | Tabelle vorhanden (rund 50 Zeilen); acht Gruppen davon selbst nachgebaut, alle rot wie behauptet (unten) |
| 4: `make gates` grün | bestätigt, durch den Nachweis getragen | `.harness/state/gates-passed.dateiliste` 10:37, 13 Minuten nach `07faa02` (10:24), passt zu „13 min 2 s“; `git diff --name-status 07faa02..498ace8` nennt nur die Plan-Datei (Positivliste). Vor meinem Lauf blockierte der Stop-Hook, weil der Stand der schnellen Kette (09:21) älter war; nach `make gates-schnell` (Exit 0) gibt er `approve` |
| 5: Review durchgeführt, Report liegt vor | bestätigt | `docs/reviews/2026-10-16-review-slice-v1-abschluss-verbindungsende-replay-extended-flush.md` (0 HIGH, 0 MEDIUM, 2 LOW, 2 INFO), F-1 und F-2 nachgearbeitet, kein Self-Review |
| 6 bis 9: Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen | offen (erwartet) | §2 ohne Häkchen, §7 trägt Belege und Nacharbeit, die drei Risiken in §6 tragen „offen bis Closure“; Closure ist Handlung des Planners |

## Entscheidungs-Konformität

| Entscheidung | Urteil | Beleg |
|---|---|---|
| `LH-FA-25.a` *Replay*, *Form*, *Stelle von `connection_closed`*, *Reihenfolge der Gründe* | konform | Tests prüfen: letzte Antwort `connection_closed`, kein `ready_for_query`, Lieferung erst beim abschließenden `Flush`, Abweichung davor `PGR-E5001` ohne Ende, Cursor bleibt stehen |
| `LH-FA-18.a` *Interaktion* (Abschluss = `Sync`, bei Ende und `Flush` dieses `Flush`; „vor ihrem Abschluss“), Commits `ef6ccc9`, `9ef7dfb`, `0f89aee` | konform | `mitten()` (`c.gruppe > 0 \|\| c.nachricht > 0`) entscheidet einfache Anfrage und Lebendprüfung mitten in und zwischen den Gruppen als `PGR-E5001`; Fälle und Mutanten M17e/M17f (hier A und B) rot |
| `LH-FA-13.a` „bis zu ihrem Abschluss“ (`7dcbd12`) | konform | `Shutdown` erlaubt nur außerhalb der Interaktion; `TestReplayEndeExtendedFlushHerunterfahren` Stände 0 bis 6 |
| `SPEC-001`/`validateAbschluss` | konform | Dargestellt sind nur Gruppen davor mit `Flush` und Ende-Gruppe mit `Flush`; die laut Format nicht darstellbaren Fälle (`PGR-E3003`) haben nach §6 keinen Fall |
| Randformen in §6 vor dem Code entschieden | konform | Commit `0f89aee` (Architect) liegt vor dem Test-Commit `3d02a29`; keine Randform im Code-Commit eingetragen |

## Nacharbeit (Review F-1, F-2; `07faa02`, `498ace8`, `7dcbd12`)

- `git diff e982c7b..498ace8 -- internal ':!*_test.go'` zeigt genau zwei Hunks, beide nur Kommentarzeilen: `replay.go` (`mitten()`: „vor ihrem Abschluss (LH-FA-18.a)“) und `server.go` (Kopf von `replaySitzung`: „bis zu ihrem Abschluss (Sync, bei einem Ende des Servers das abschließende Flush; LH-FA-13.a, LH-FA-18.a)“). Keine Logikzeile. Die Kommentare stimmen mit `LH-FA-18.a` *Interaktion* und `LH-FA-13.a` (`spec/spezifikation.md`, `7dcbd12`) überein.
- Grep nach „Sync“ in Kommentaren des Produktcodes und der Tests: Treffer, bei denen `Flush` gemeint sein kann, nur in V-1 und V-2. Kein Fund: `replay.go:384` („bis zum Sync verworfenes parse“, Serververhalten bei Fehlern, Sync ist dort richtig), Tests des Gebers und `server_extended_test.go:724/805` (Interaktionen mit `Sync`, ohne Ende; dort richtig, Review nahm sie bewusst aus), `model/extended.go`, `ende_test.go`, Recording-Adapter (nennen `sync oder flush`).
- Plan §3 und §7 sind nachgezogen (`AGENTS.md` §3.9).
- Messung nachgerechnet: `git diff --numstat e982c7b..498ace8 -- internal test docs/user` = Einfügungen 310, Löschungen 39, **Summe 349**, wie behauptet (342 plus 7); Marke 500 (70 %) und Meldeschwelle 400 nicht erreicht; Schätzung 160 bis 330 um 19 überschritten. Aufschlüsselung deckt sich mit §7 (Service-Tests 162 + 41/19, Adapter 36/14, E2E 27 + 20, Tabellen 20/3, Kommentare 4/3 im Produktcode).

## Bewusstes Brechen (selbst gefahren)

| Zusage | Mutant (Kopie) | Rot |
|---|---|---|
| Abweichung/Lebendprüfung mitten in der Interaktion (Glied `c.nachricht > 0`) | A: `mitten` auf `c.gruppe > 0` | `…FlushAbweichung`, `…FlushHerunterfahren`; Geber: `…Abweichung`, `…Herunterfahren`, `TestReplayLebendpruefungExtended`, `TestReplayZwangsende`, `…ZwangsendeNichtVerbraucht` |
| … (Glied `c.gruppe > 0`, zwischen zwei Gruppen) | B: `mitten` auf `c.nachricht > 0` | `…FlushAbweichung`, `…FlushHerunterfahren`; Geber wie oben |
| Verbrauch erst mit `Sent` | M8: `c.ungesendet = true` entfernt | `…FlushVerbrauch`, `…FlushHerunterfahren`; Geber `…Verbrauch`, `…Herunterfahren` |
| `Shutdown` erlaubt nur außerhalb der Interaktion (Glied `c.nachricht == 0`) | S: auf `c.gruppe == 0` verkürzt | `…FlushHerunterfahren`; Geber `…Herunterfahren`, `TestReplayExtendedHerunterfahren` |
| `Sync` statt abschließendem `Flush` ist Abweichung | T: `Sync` für erwartetes `Flush` akzeptiert | nur `…FlushAbweichung` und `…Abweichung` (Zeile Sync statt Flush) |
| `Flush` statt `Sync` ist Abweichung (Sync-Fälle nicht abgeschwächt) | T2: `Flush` für erwartetes `Sync` akzeptiert | `TestReplayEndeExtendedAbweichung` (Geber, `Sync`-Lauf): rot, also unverändert scharf |
| Adapter bildet `Flush` als `Flush` ab | MA6: `pgproto3.Flush` als `Sync` | `TestReplayEndeGruppeDavorSchliesstNicht` rot (Zeitüberschreitung 5 s) |
| Binary: `Flush` des Clients wird als `Flush` verglichen | E2: Abbildung `Flush` auf `Sync` im Arbeitsbaum, `make test-integration` | 17 Tests rot, darunter `TestE2EReplayEndeExtendedFlush` und `TestE2EReplayEndeExtendedSync` |

Parametrisierung des Gebers: `mitSync()` setzt `extendedEnde`, `endeNachrichten`, `sendeBis` ein, also dieselben Werte wie vor dem Diff; `ab.nachrichten()[4]` ist `syncNachricht()`. Die Sync-Fälle sind nicht abgeschwächt (T2 und A bis S färben sie weiter).

**Nicht gespiegelte Fälle und akzeptierte Negative.** M17c färbt nur Tests des Gebers; die Begründung trägt (`anfrageMitEnde` liest `Responses`, die eine Extended-Interaktion nicht trägt), Test des Gebers `…IstKeineEinfacheAnfrage` benannt. Das akzeptierte Negativ „`Terminate` mitten in der `Flush`-Gruppe am Adapter“ benennt `TestReplayEndeTerminateMittenInDerGruppe` (existiert, `server_ende_test.go:286`); die einzige `Flush`-abhängige Stelle im Adapter (Abbildung) färbt MA6. Das zweite akzeptierte Negativ (teilweises Senden) verweist auf `TestReplayEndeSchreibfehlerBleibtUnverbraucht` des Gebers.

**Abdeckung.** `make abdeckung-check` grün (Exit 0). Aussagen der Zeilen reichen nicht weiter als die Tests: „an jeder Stelle“ ist durch fünf Zeilen der Abweichungstabelle, `Terminate` „vor, in und nach der Gruppe davor und mitten in der Gruppe mit Ende“ durch die Stände 0 bis 4 gedeckt (§3.11).

## Läufe

- `make gates-schnell`: Exit 0 (an `498ace8`); der Stop-Hook gibt danach `approve`.
- `make abdeckung-check`: Exit 0.
- `make docs-check`: siehe Commit (0 Befunde vor dem Commit).
- Voller `make gates`: nicht neu gefahren (Auftrag); Nachweis `.harness/state/gates-passed.dateiliste` an `07faa02`, Code seitdem unverändert.

## Findings

| ID | Schwere | Befund | Zuständig |
|---|---|---|---|
| V-1 | LOW | `internal/hexagon/ports/driving/replay.go` Zeile 37: Der Kommentar von `Shutdown` sagt „wenn keine Extended-Interaktion läuft, deren Sync noch aussteht (LH-FA-13.a)“. `LH-FA-13.a` sagt seit `7dcbd12` „bis zu ihrem Abschluss“; bei einer Interaktion mit Ende und `Flush` ist das ihr `Flush`, und der Service (`Shutdown`, `mitten`-Logik) entscheidet das so. Dieselbe Klasse wie F-1/F-2 (Kommentar nennt Sync, Abschluss kann Flush sein), an einer Stelle, die Review und die Grep-Suche der Nacharbeit nicht erfassten (Port statt Service/Adapter). Keine Logik. | Implementer (Kommentar auf „deren Abschluss“, Plan §3 nachziehen); Planner entscheidet, ob in diesem Slice oder bei der Closure; nicht blockierend |
| V-2 | INFO | `internal/hexagon/services/record.go` (Zeilen 48, 198, 262) und `ports/driving/record.go` Zeile 40 sagen „vor deren Sync“ bzw. „deren Sync schon gesendet ist“. Es ist der Record-Pfad, den dieser Slice laut §1 ausschließt (`slice-v1-abschluss-verbindungsende-record`). Ob dort eine Interaktion mit Ende und `Flush` entsteht, entscheidet jener Slice; die Lesart „Abschluss“ aus `LH-FA-18.a` gilt für `Query` „vor deren Abschluss“ (Spezifikation 1556). Eine Sendung an jenen Slice (`AGENTS.md` §3.13: Eintrag in §1 oder DoD mit Kennung des Gebers, Nachzählen in §8) ist nicht eingetragen. | Planner (Sendung prüfen, ob der Nehmer annimmt) |

Kein Finding zu: Plan-vs-Code (Plan folgt dem Diff: Produktcode 0 Logik, zwei Kommentare; zwei Liefer-Punkte; Rückführungen §4 nicht eingetreten), Messung, Fixture je Fall, Abdeckungstabellen, Parametrisierung des Gebers.

## Vorschlag für die Risiko-Ausgänge zur Closure (nicht eingetragen)

1. *Diff über einer Review-Sitzung*: **entfallen** mit Begründung: Messung 349, 70 % der Marke 500, Meldeschwelle 400 nicht erreicht, Produktcode 0 Logik. Die Klasse `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` trat dennoch auf (19 über dem oberen Ende 330); das ist eine Beobachtung für das Register, kein offenes Risiko.
2. *Pfad des Gebers trägt die `Flush`-Gruppe nicht*: **entfallen** mit Begründung: Der Service lieferte sie ohne Codeänderung, 55 Mutanten des Implementers und die hier gefahrenen sind rot, Port und Abbildung unverändert.
3. *Handgeschriebene Aufzeichnungen bilden nicht ab, was PostgreSQL liefert*: **weiter offen**, ins Beobachtungs-Register, bis der E2E-Test von `slice-v1-abschluss-verbindungsende-upstream` die Entsprechung belegt.

Dazu für den Lerneintrag: Finding-Klasse „Kommentar nennt Sync, Abschluss kann Flush sein“ (F-1, F-2, V-1) gehört in den Zähler; V-1 zeigt, dass die Grep-Suche der Nacharbeit die Ports übersah.
