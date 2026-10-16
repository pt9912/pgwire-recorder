# Review-Report: slice-v1-abschluss-verbindungsende-record — 2026-10-16

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-record` (Schritt 40), Diff `9c19eba6..ae01b2a0`; Maßstab daneben die Architect-Commits `23374c20` und `d5c42925` (Spezifikation `LH-FA-25.a`, *Stelle von `connection_closed` beim Aufzeichnen*).

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec08 · <!-- d-check:ignore (Skill-Pfad des Repos, Zitat des Laufs) -->
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

- Slice-Plan `slice-v1-abschluss-verbindungsende-record` (§1, §3, §4, §6, §7 *Belege des Implementers*, §8)
- Spezifikation: `LH-FA-25.a` (*Aufzeichnen*, *Stelle*, *Meldung des Endes beim Aufzeichnen*, *Stelle von `connection_closed` beim Aufzeichnen*, *Reihenfolge der Gründe beim Aufzeichnen*, *Herunterfahren beim Aufzeichnen*), `LH-FA-18.a` *Abbruch*, `LH-FA-13.a`, `LH-FA-12.a`, `SPEC-002`, `SPEC-043`
- `AGENTS.md` §3.7 bis §3.14 (Hard Rules); Hexagon-Richtung (`.a-check.yml`: Imports im Diff unverändert, kein neuer Import in Produktcode)
- Eigene Läufe: Mutanten in Kopien des Baums (Image `deps`, `--network=none`, ein Lauf zur Zeit, Arbeitsbaum unberührt) — Service 12 (n1, n1b bis n1d, n3a, n3b, n4, n2, M9b, M17 bis M19, M21, M26), Adapter 1 (A7); Stabilität `go test -count=60 -run TestRecordEnde` über `services` und `pgwire` grün; Messung des Plans nachgerechnet: `git diff --numstat 9c19eba6..ae01b2a0 -- internal test docs/user` = 825 (800 Einfügungen, 25 Löschungen), Produktcode 130 (davon 51 hinzugefügte und 15 entfernte Kommentarzeilen = 66), Tests 674, Tabelle 21; stimmt mit §7. Port und Format unverändert außer Kommentaren in `ports/driving/record.go` und dem Ereignis `EndServer` in `model/recording.go`

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | `LH-FA-25.a` *Herunterfahren beim Aufzeichnen* sagt zu: Hat der Recorder das Ende gelesen, läuft keine Interaktion mehr, es entsteht keine Meldung, und die Session wird geschrieben, „auch wenn die Zustellung noch aussteht“. Jeder Test ruft vor `CloseSession(EndForced)` schon `Delivered`; der Zustand zwischen Lesen und Zustellen ist ungeprüft. Selbst gefahren grün: `CloseSession` mit `end == model.EndForced && (l.laeuft() \|\| l.serverEnde && l.unzugestellt > 0)` (Zwangsende meldete dann `PGR-E4006` für ein gelesenes Ende) und `laeuft()` mit `\|\| l.serverEnde && l.unzugestellt > 0`; die ganze Suite des Pakets `services` bleibt grün. Die Tabelle in §7 führt den Satz als Zusage („`EndForced` nach der Interaktion mit Ende: geschrieben, kein `PGR-E4006`“), die Mutation hat nur den Zustand nach `Delivered`. | `AGENTS.md` §3.10 (Ortsangabe/Aufzählung zählt je Element), `LH-FA-25.a` | `internal/hexagon/services/record_ende_test.go` · "if !l.s.Delivered(context.Background(), l.id) {" (`TestRecordEndeEnden`, `Delivered` vor `l.schliessen(end)`); `internal/hexagon/services/record.go` · "if end == model.EndForced && l.laeuft() {" | ja — `go test ./internal/hexagon/services` mit der genannten Mutation bleibt grün | Zusage der Spezifikation ohne Mutation: Zustand zwischen Lesen und Zustellen |
| F-2 | MEDIUM | Die Grenze (1) in §7 nennt als Adresse für den Test des Schreibfehlers in `serverRichtung` (Extended) `slice-v1-abschluss-verbindungsende-upstream` („Test-Idee … Adresse“). Dessen §1 und DoD führen die Sendung nicht (kein Treffer für `Schreibfehler` oder `serverRichtung`; die Übernahme aus diesem Slice nennt nur Fakes und Abnahmeszenario 18), kein Eintrag mit der Kennung des Gebers, keine Nachzählung in §8 des Nehmers; der E2E des Nehmers läuft zudem gegen das Binary, nicht gegen einen Schreibfehler am Client. Der Pfad bleibt ohne Test und ohne annehmende Adresse. | `AGENTS.md` §3.13 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-record.md` · "Adresse: `slice-v1-abschluss-verbindungsende-upstream` (E2E)"; `docs/plan/planning/next/slice-v1-abschluss-verbindungsende-upstream.md` · "Übernommen aus `slice-v1-abschluss-verbindungsende-record`" | ja — `grep -n "Schreibfehler" docs/plan/planning/next/slice-v1-abschluss-verbindungsende-upstream.md` ohne Treffer | Adresse nimmt nicht an |
| F-3 | LOW | §7 ordnet M17 (Prüfung `serverEnde` am Anfang von `Query` entfernt) ohne `Shutdown` als „äquivalent, doppelte Absicherung“ ein. Das gilt nur nach `Delivered`: Ist das Ende gelesen und `Delivered` steht aus (`empfaengt` wahr), liefert `Query` mit der Prüfung sofort `ErrSessionEnded`, ohne sie blockiert es bis `Delivered`. Selbst gefahren: ein Test „Anfrage nach dem Lesen des Endes, vor `Delivered`, kehrt binnen 300 ms zurück“ ist im Original grün und mit dem Mutanten rot; `TestRecordEndeDanachNichts` hat keinen solchen Fall. Ein Ausfall ist nicht zu erzählen (der Adapter ruft `Delivered` immer), daher LOW. | `AGENTS.md` §3.10, Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-record.md` · "ohne Shutdown trägt die Prüfung nach dem Warten, äquivalent, doppelte Absicherung" | ja — Mutant (Prüfung entfernt) plus Aufruf von `Query` zwischen `AwaitServer` und `Delivered` | Grüner Mutant als äquivalent eingeordnet, Bedingung ungenau |
| F-4 | INFO | Nach dem Ende liefern `Query` und `ClientMessage` `ErrSessionEnded` auch bei gemeldetem `Shutdown` (die Prüfung steht vor der auf `ErrShutdown`); `TestRecordEndeDanachNichts` („Anfrage, herunterfahren true“) legt das fest. Weder §6 noch `LH-FA-25.a` nennen den Vorrang der beiden Fehlerwerte, die Spezifikation sagt „nimmt keine Anfrage … mehr an“; der Port-Kommentar nennt `ErrSessionEnded`. Zuständig: Architect, ohne erwartete Aktion, wenn der Port-Kommentar als Stelle der Entscheidung genügt. | `AGENTS.md` §3.12 | `internal/hexagon/services/record.go` · "if l.serverEnde {" (Anfang von `Query`) | nein — Lesart | Randform im Code entschieden, in §6 nicht genannt |
| F-5 | INFO | Der Schreibfehler in der Extended-Richtung des Adapters (`serverRichtung`) und `EndTerminate` nach dem Ende haben keinen Test; die Kommentare an `schreibe` und `CloseSession` sagen dazu nichts zu (der Satz in `CloseSession` nennt `EndClosed` und `EndForced`, nicht `EndTerminate`) — kein Verstoß gegen `AGENTS.md` §3.11. `schreibe` ist dieselbe Funktion wie im einfachen Pfad; ein `Terminate` zwischen Lesen des Endes und `beende(EndServer)` führt in `CloseSession` mit `laeuft() == false` zum selben Ergebnis. Zur Adresse siehe F-2. Zuständig: Implementer, ohne erwartete Aktion. | `AGENTS.md` §3.11 | `internal/adapters/driving/pgwire/server.go` · "func (r *richtungen) schreibe(ctx context.Context, rs []model.Response) bool {" | nein — kein Gate-Lauf hat den Pfad | Pfad ohne eigenen Test, Kommentar zusagt nichts |
| F-6 | INFO | Äquivalent bestätigt: A7 (`schreibe` liefert nach `beende(EndServer)` wahr, Wache in `clientRichtung` bleibt) selbst gefahren, `go test` des Pakets `pgwire` grün; die Invariante trägt (`beende` setzt `r.beendet` vor `CloseSession`, `AwaitServer` liefert danach `ErrSessionEnded`). Die Doppelabsicherung von M9 (`pruefeEnde` und `ohneReadyForQuery`) ist ebenso äquivalent, weil beide denselben Fall `[Fehlerantwort, ready_for_query, connection_closed]` mit `PGR-E6001` ablehnen; `ohneReadyForQuery` allein trägt „`ready_for_query` davor“ (M8, rot). Kein Befund gegen den Diff. | `AGENTS.md` §3.10 | `internal/adapters/driving/pgwire/server.go` · "r.beende(ctx, model.EndServer, nil)"; `internal/hexagon/model/response.go` · "func ohneReadyForQuery" | ja — Mutant A7, Lauf von `go test` über `pgwire` | — |
| F-7 | INFO | §8 zählt `docs/user/abdeckung-unit.md` weiter „konservativ“ als dritte Schicht, die Entscheidung des Nutzers steht aus; mit der Messung 825 liegt der Slice unter der Marke 1400 und der Meldeschwelle, die Schichtzählung bleibt eine offene Entscheidung für die Closure. Zuständig: Planner. | Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-record.md` · "konservativ als dritte; die Entscheidung des Nutzers über die Abdeckungstabelle steht aus" | nein — Entscheidung | Schichtzählung der erzeugten Tabelle offen |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services/record.go` — Mutationstabelle §7 (Stichprobe aus jeder Gruppe: `bisEnde` in `Query` und in `AwaitServer` einzeln, `zuordnen` ohne `unsupported` beim Ende ohne laufende Interaktion, M9b, M17 bis M19, M21, M26; Schweregrad-, Stellen-, Gruppen-, Schreibfehler-Zeilen gelesen gegen die Tests) | geprüft; außer F-1 und F-3 jeder Mutant rot, die Tests der Zusage „danach nichts mehr“ trennen die Wächter (Session bleibt offen, mit und ohne `Shutdown`) |
| `internal/hexagon/services/record_ende_test.go` — Stelle a bis e (Anfang, Mitte, `ERROR` davor, Nachricht dazwischen, `ready_for_query` davor/dazwischen im Extended-Pfad = `PGR-E1000`-Grenze, Sync-/Flush-Gruppe, nicht abgeschlossene Gruppe, `ERROR` vor `FATAL`, Fehler des Upstream-Ports vor der Form) | geprüft, ohne Befund über F-1 hinaus |
| `internal/adapters/driving/pgwire/server.go` und `server_record_ende_test.go` — `schreibe`, `beende(EndServer)`, `send` lässt `connection_closed` aus, `Delivered` vor dem Schließen, kein Eintrag im Log, Verhalten bei `Shutdown`; Tests 60-mal stabil | geprüft, ohne Befund über F-5 hinaus |
| `internal/hexagon/model/recording.go`, `ports/driving/record.go` (nur `EndServer` ohne Logik und Kommentare; Port und Format sonst unverändert) | geprüft, ohne Befund |
| Kommentare im Diff (`AGENTS.md` §3.7, §3.11): Indikativ, keine verworfene Alternative, jede Zusage hat einen Test (`serverEnde`, `serverBeendet`, `EndServer`, `schreibe`, Port); `EndTerminate` nicht zugesagt | geprüft, ohne Befund |
| Wort „Sync“ → „Abschluss“ über `internal/` einschließlich `ports/` (`grep -rn "Sync" internal --include=*.go` ohne Tests): Treffer in `record.go` (Zeilen 73, 81, 207, 279, 282), `ports/driving/record.go:49`, `ports/driving/replay.go:32`, `model/extended.go:62`, `replay.go:384`, `play.go:312/315`, `yaml.go:685`, `server.go:334` meinen das `Sync` selbst oder nennen den Abschluss mit Ende im Satz | geprüft, ohne Befund |
| Hexagon-Richtung, Core-Reinheit (`internal/hexagon/**`): keine neuen Imports im Produktcode des Diffs, `pgproto3` nur im PGWire-Adapter | geprüft, ohne Befund |
| Suppression (`AGENTS.md` §3.2): kein `//nolint` im Diff | geprüft, ohne Befund |
| Plan folgt Diff (`AGENTS.md` §3.9): §1, §3, §6, §7, Kopf; Beobachtung (3) in §7 trägt die Entscheidung `d5c42925` nicht nach — Aufgabe des Planners in der Closure, kein Befund | geprüft, ohne Befund |
| Messung §7 (825, 130/66, 674, 21; Marke 1400, Meldeschwelle 1000, Basis `9c19eba6`) | nachgerechnet, stimmt |
| Randformen (`AGENTS.md` §3.12), soweit §6 sie nennt | geprüft; nicht genannte Entscheidung siehe F-4 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Zusage der Spezifikation ohne Mutation: Zustand zwischen Lesen und Zustellen · Adresse nimmt nicht an · Grüner Mutant als äquivalent eingeordnet, Bedingung ungenau · Randform im Code entschieden, in §6 nicht genannt · Pfad ohne eigenen Test, Kommentar zusagt nichts · Schichtzählung der erzeugten Tabelle offen

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 (MEDIUM) blockieren: F-1 verlangt vom Implementer einen Fall mit Mutation für den Zustand zwischen Lesen des Endes und `Delivered` (Zwangsende und `Shutdown`), F-2 verlangt die Sendung samt Nachzählung beim Nehmer (Planner) oder einen Test in diesem Slice. F-3 bis F-7 blockieren nicht.

**Übergabe:** Findings gehen an den Implementer (F-1, F-3, F-5), an den Planner (F-2, F-7) und an den Architect (F-4); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
