# Review-Report: slice-v1-abschluss-verbindungsende-replay-extended — 2026-10-16

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-replay-extended` (Schritt 37b), Diff `1896db3ea5103e462a7aec895ea877f43b7b6338..6f8a6ee`.

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec0 · <!-- d-check:ignore (Skill-Pfad des Repos, Zitat des Laufs) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-v1-abschluss-verbindungsende-replay-extended` (§1, §3, §4, §6, §7 *Belege des Implementers*, §8)
- Spezifikation: `LH-FA-25.a` (*Meldung des Endes im Replay*, *Stelle von `connection_closed` im Replay*, *Reihenfolge der Gründe im Replay*; Architect-Commit `ef6ccc9`), `LH-FA-09.a`, `LH-FA-12.a`, `LH-FA-18.a`, `LH-FA-03.b`
- `slice-v1-abschluss-verbindungsende-replay-extended-flush` (Nehmer der Sendung des Architect)
- `AGENTS.md` §3.7 bis §3.14 (Hard Rules)
- Eigene Läufe: Mutanten in Kopien des Baums (Image `deps`, `--network=none`, kein Eingriff in den Arbeitsbaum) — Service 19, Adapter 4, Binary 5 (Binary ohne PostgreSQL, `PGR_OHNE_POSTGRES=1`); Produktcode `internal/hexagon/services/replay.go` und `internal/adapters/driving/pgwire/server.go` im Diff: 0 Zeilen (bestätigt mit `git diff --stat`); Messung des Plans nachgerechnet: `git diff --numstat` über `internal test docs/user` = 502 (Einfügungen 502, Löschungen 0), stimmt mit §7

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Spezifikation (`LH-FA-25.a`, Zuständigkeit Architect) und §6 des Plans (Zuständigkeit Planner) sagen zu, dass Gruppen vor der Gruppe mit Ende „mit `Sync` wie mit `Flush`“ laufen und „beide Formen Fälle sind“. `Interaction.Validate` (`validateAbschluss`) lässt `Sync` nur am Ende der letzten Gruppe zu, eine Gruppe mit `Sync` vor der Gruppe mit Ende kann das Format nicht darstellen (`PGR-E3003` beim Laden); getestet ist nur die frühere Interaktion mit `Sync` (`TestReplayEndeExtendedSyncDavor`) und die Gruppe mit `Flush` davor. Der Wortlaut sagt mehr zu, als es einen Fall geben kann; §6 und die Spezifikation stehen nach dem Lauf des Implementers unverändert. | `AGENTS.md` §3.11, §3.9, §3.12 | `spec/spezifikation.md` · "`Sync` wie mit `Flush`, laufen wie ohne Ende"; `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-replay-extended.md` · "nach Gruppen davor, mit `Sync` und mit `Flush`: beide Formen sind Fälle" | ja — `go test` des Modells (`Interaction.Validate`) zeigt die Ablehnung einer Gruppe mit `Sync` vor der letzten; ein Fall mit der Zusage lässt sich nicht schreiben | Zusage im Wortlaut weiter als Format und Test |
| F-2 | LOW | Das Gegenbeispiel zur Zusage aus F-1 steht als „Beobachtung (1)“ in §7 mit „zur Kenntnis“ für den Architect; eine Rückgabe an den Architect (Artefakt, das §6 oder die Spezifikation ändert) liegt im Diff nicht vor. `AGENTS.md` §3.12 verlangt für ein bemerktes Gegenbeispiel Anhalten und Rückgabe, ein Satz im Bericht gilt dort nicht als Rückgabe (Wortlaut steht im Absatz zur Aufruf-Voraussetzung). | `AGENTS.md` §3.12 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-replay-extended.md` · "bleibt für den Architect zur Kenntnis (keine Änderung am Code, keine Randform entschieden)" | nein — Prozessbefund, kein Gate | Gegenbeispiel nur als Satz im Bericht |
| F-3 | LOW | Die Abdeckungs-Deklaration des E2E-Falls sagt, die Interaktion mit `Sync` und die Gruppe mit `Flush` davor „laufen wie sonst“. Der Test liest sie nur bis zum Typ `Z` bzw. `1` (`bisTyp` überspringt alles davor): Ein Mutant im Adapter, der bei `ParseComplete` + `ReadyForQuery` das `ParseComplete` auslässt, bleibt im Binary-Test grün (selbst gefahren: `TestE2EReplayEndeExtendedSync` ok) und färbt auch im Paket des Adapters keinen Test. Der Inhalt der Gruppen davor ist durch die Unit-Tests des Service getragen, nicht durch die Aussage dieses Falls. | `AGENTS.md` §3.11 | `test/integration/replay_ende_extended_e2e_test.go` · "Gruppe mit Flush davor laufen wie sonst"; `docs/user/abdeckung-e2e.md` · `TestE2EReplayEndeExtendedSync` | ja — Mutant in `send` (Auslassen der ersten Antwort bei zwei Antworten `parse_complete`, `ready_for_query`) + `make test-integration` bleibt grün | Abdeckungs-Deklaration weiter als Prüfung |
| F-4 | INFO | Der Mutant „Cursor rückt bei Abweichung vor“ (M23) färbt `TestReplayEndeExtendedAbweichung` nur über eine Panik (`index out of range [2] with length 2`), die den Lauf des Pakets abbricht; eine Prüfung (`Fatalf`) färbt ihn nicht. §7 nennt das (Laufzeitfehler, kein Fehlschlag der Prüfung). Rot ist er, die übrigen Tests des Pakets melden bei diesem Mutanten nichts. Zuständig: Implementer, ohne erwartete Aktion. | Maintainability | `internal/hexagon/services/replay_ende_extended_test.go` · "das Ende kommt nach der Abweichung noch" | ja — Mutant `c.nachricht++` im Zweig `default` der Abweichung | Mutant nur über Panik rot |
| F-5 | INFO | Der als äquivalent bezeichnete Mutant M17d (`!s.anfrageMitEnde(c, sql)` aus der Bedingung in `Query`) ist für eine Extended-Interaktion äquivalent (`endetEinfach` liest `Responses`, die `Interaction.Validate` einer Extended-Interaktion verbietet). Die Behauptung, `slice-v1-abschluss-verbindungsende-replay` decke ihn ab, trägt: selbst gefahren, rot sind `TestReplayEndeLebendpruefungMitEndeBleibtInteraktion`, `TestReplayEndeLebendpruefungAndererTextAusserDerReihe`, `TestReplayEndeFirstRequestLebendpruefungVorZuordnung`, `TestReplayConnectionEndeBleibtZugeordnet`. Kein Finding gegen den Diff. | `AGENTS.md` §3.10 | `internal/hexagon/services/replay.go` · "!s.anfrageMitEnde(c, sql)" | ja — Mutant, Lauf von `go test` über das Paket | — |
| F-6 | INFO | `LH-FA-25.a` nennt als Abweichung „eine Lebendprüfung mitten in der Gruppe“; Code und Test (Zeile „Lebendprüfung nach der Gruppe davor“) entscheiden auch die Lebendprüfung zwischen zwei Gruppen als `PGR-E5001`. Die Entscheidung steht in `LH-FA-18.a` („vor deren `Sync`, auch eine Lebendprüfung“), nicht in `LH-FA-25.a`; §6 nennt „zwischen den Gruppen“ nicht. Zuständig: Architect, ohne erwartete Aktion, solange `LH-FA-18.a` die Stelle trägt. | `AGENTS.md` §3.12 | `spec/spezifikation.md` · "Anfrage oder eine Lebendprüfung mitten in der Gruppe" | nein — Lesart | Stelle einer Nachricht nur über anderen Absatz entschieden |
| F-7 | INFO | Das akzeptierte Negativ „`Terminate` oder Abbruch mitten in der Gruppe: durch vorhandene Tests des Adapters getragen“ benennt keinen Test. Der Service-Teil ist belegt (`TestReplayEndeExtendedTerminateUnverbraucht`, `CloseConnection` für beides); ein Replay-Test des Adapters für `Terminate` mitten in einer Extended-Gruppe ist im Diff nicht benannt und von mir nicht gefunden (die `Terminate`-Fälle in `server_extended_test.go` gehören zum Record-Pfad). Zuständig: Planner/Implementer, ohne erwartete Aktion. | `AGENTS.md` §3.11 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-replay-extended.md` · "durch dessen vorhandene Tests getragen" | nein — Suche, kein Gate | Akzeptiertes Negativ ohne benannten Test |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services/replay.go`, `internal/adapters/driving/pgwire/server.go` (Produktcode, 0 Zeilen laut Plan) | geprüft: im Diff unverändert; Mutanten dort gefahren (siehe unten), ohne grünen Mutanten |
| `internal/hexagon/services/replay_ende_extended_test.go` (Service-Tests) | geprüft; 19 Mutanten (M5a, M5b, M8, M9, M11, M12, M15, M17a, M17d, M17e, M17f, M18a, M18b, M20b, M23 sowie C1 bis C4: `ungesendet` in `CloseConnection`, Nummer der ersten nicht verbrauchten, Index in `Forced`, Grenze `verbraucht >= alle`) rot; Konjunktions-Glieder von `mitten` und `Shutdown`, Reihenfolge der Gründe (`PGR-E5003` vor Abweichung), Abweichung je Stelle (erste, mittlere, abschließende; `Sync` statt `Flush`), Fixture je Fall: ohne Befund |
| `internal/adapters/driving/pgwire/server_ende_test.go` (Adapter-Tests) | geprüft; MAx (Ende im Pfad `ClientMessage` ignoriert), MAlen0 (leere Antwort nicht übersprungen) rot; Ablauf `gruppe sent gruppe sent close`: ohne Befund |
| `test/integration/replay_ende_extended_e2e_test.go`, `testdata/ende-extended-sync.yaml` (Binary) | geprüft; M13, M5a, M5b, E2 (Adapter schließt nicht, Zeitüberschreitung nach 10 s) rot; eine Aufzeichnung je Fall; Befund nur F-3 |
| `spec/spezifikation.md` (`LH-FA-25.a`, Commit `ef6ccc9`) | geprüft; Befund F-1 (Wortlaut), F-6 (Hinweis) |
| `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-replay-extended.md` (§1, §3, §4, §6, §7, §8) | geprüft; Plan folgt dem Diff (Produktcode 0, ein E2E-Fall, Messung 502 nachgerechnet, Meldeschwelle 550 und Marke 700 nicht erreicht); Befund F-1 (§6), F-2, F-7 |
| `docs/plan/planning/next/slice-v1-abschluss-verbindungsende-replay-extended-flush.md` (Nehmer) | geprüft; Sendung mit Kennung des Gebers in §1/DoD, §8 nachgezählt (Schichten drei konservativ, Nutzerentscheidung steht aus, unverändert), ohne neuen Befund |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | geprüft; `make abdeckung-check` und `make kopf-check` ohne Ausgabe und Exit 0; Befund F-3 zur Aussage der E2E-Zeile |
| `docs/plan/planning/in-progress/roadmap.md` | geprüft, ohne Befund |

Der Arbeitsbaum war nach den Läufen sauber (alle Mutanten liefen in Kopien unter dem Scratchpad); `make docs-check` und die Gates `make gates` fuhr dieser Lauf nicht (Verifier-Aufgabe, `make docs-check` siehe Commit).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Zusage im Wortlaut weiter als Format und Test · Gegenbeispiel nur als Satz im Bericht · Abdeckungs-Deklaration weiter als Prüfung · Mutant nur über Panik rot · Akzeptiertes Negativ ohne benannten Test

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM): Der Wortlaut in der Spezifikation und in §6 sagt einen Fall zu, den das Format nicht darstellt; Spezifikation (Architect) und §6 (Planner) folgen dem Befund des Implementers nicht. Kein Produktcode betroffen; die Tests des Diffs tragen die darstellbaren Fälle (Gruppe mit `Flush` davor, frühere Interaktion mit `Sync`). F-2 und F-3 sind LOW und blockieren nicht.

**Übergabe:** F-1 an den Architect (Spezifikation) und den Planner (§6, Randform-Zeile „Ende-Gruppe als einzige Gruppe …“), F-2, F-3, F-4 und F-7 an den Implementer, F-5 und F-6 zur Kenntnis; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
