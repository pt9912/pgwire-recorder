# Review-Report: slice-v1-abschluss-verbindungsende-replay — 2026-10-16

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-replay` (Schritt 37a), Diff `13c23a4..a65e8d7` (Basis `13c23a4559fb0df9e5718bf078c32053cdf10a5b`).

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

- Slice-Plan `slice-v1-abschluss-verbindungsende-replay` (§1, §3, §4, §6, §7 *Belege des Implementers*, §8)
- Spezifikation: `LH-FA-25.a` (*Meldung des Endes im Replay*, *Stelle von `connection_closed` im Replay*, *Reihenfolge der Gründe im Replay*), `LH-FA-09.a` (*Lebendprüfung mit Ende*), `LH-FA-12.a`, `LH-FA-03.b`, `SPEC-043`; Architect-Commit `94a3a24`
- `slice-v1-abschluss-verbindungsende-replay-extended` (§1, §6: Adresse für die Extended-Interaktion mit Ende)
- `AGENTS.md` §3.3, §3.7 bis §3.14 (Hard Rules), Hexagon-Richtung (`make a-check` nicht neu gefahren; der Diff berührt `internal/hexagon` nur ohne neue Imports)

**Eigene Prüfung statt Zusage.** Mutanten in Kopien von `go.mod`, `go.sum`, `internal/` außerhalb des Repos (Image `deps`, `--network=none`), Basis grün. Rot gesehen: `anfrageMitEnde` im Zweig `first-request` auf `false`, `frei[0]` durch die letzte freie Session ersetzt, beliebige statt erster Interaktion, Zweig `c.session` entfernt, `&& !endetEinfach(in)` und `&& !s.anfrageMitEnde(c, sql)` entfernt (Service); `Sent` nach dem Ende nicht gemeldet (Adapter); Schließen nach dem Ende entfernt (`make test-integration` rot im Gesamtlauf, der rote Einzeltest wurde nicht ausgelesen). Grün geblieben: siehe F-1, F-2, F-6. Arbeitsbaum nach dem Lauf mit `git checkout -- internal` zurückgesetzt.

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** (Single Source of Truth); Quelle der Spalten:
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Mutant `case !s.nachVerbindung && len(s.frei) > 0:` zu `case len(s.frei) > 0:` bleibt grün. Bei `--session-assignment connection` sendet eine Verbindung ohne Session n (Kennung über der Zahl der Sessions) eine Lebendprüfung, deren Text der ersten Interaktion einer Session mit Ende in `frei` gleicht; mit dem Mutanten antwortet das Replay `PGR-E5003` statt außer der Reihe. Die Randform (Verbindung ohne Session bei `connection` mit Lebendprüfung mit Ende) nennt §6 nicht; der Code entscheidet sie, `connection` kommt nur als zugeordnete Session vor, und die Spezifikation bindet nur „Bei `first-request`“. | `AGENTS.md` §3.10 (Ortsangabe zählt je Schritt), §3.12 (Randform nicht in §6) | `internal/hexagon/services/replay.go` · `case !s.nachVerbindung && len(s.frei) > 0:` | ja — Mutant in Kopie, `go test -run 'Ende\|Connection\|Lebend'` bleibt grün | Aufzählung/Ortsangabe: Mutation je Seite fehlt |
| F-2 | LOW | Der Mutant, der `erwartet[0].Request.Type == model.RequestQuery` aus `anfrageMitEnde` entfernt, bleibt grün. Zeigt der Cursor auf eine Extended-Interaktion mit Ende, wäre eine eingehende leere Lebendprüfung ohne die Prüfung „die Anfrage“; `-replay-extended` §6 verweist für Lebendprüfungen auf diesen Slice („hier wird nichts neu entschieden“), der Fall steht hier nicht. | `AGENTS.md` §3.10 | `internal/hexagon/services/replay.go` · `erwartet[0].Request.Type == model.RequestQuery &&` | ja — Mutant in Kopie | Zusage ohne Mutation (Randzweig) |
| F-3 | LOW | Plan §3 sagt „je Fall eine Datei“, die Einbettung im E2E-Test „jede gehört einem Fall“; `ende-fatal.yaml` trägt zwei Fälle (`TestE2EReplayEndeFatal` und `TestE2EReplayEndeTerminateDavor`, der „1 von 2, Nummer 2“ aus ihr liest). Eine Änderung für den einen Fall färbt den anderen. | `AGENTS.md` §3.9 (Plan folgt dem Diff), §3.11 | `test/integration/replay_ende_e2e_test.go` · `jede gehört einem Fall` | nein — kein Gate prüft den Plan gegen den Diff | Plan-Zeile weiter als der Diff |
| F-4 | LOW | Zwei Testkommentare sagen mehr zu, als der Test prüft: `TestReplayEndeNurAlsLetzteAntwort` nennt „kein Senden des Rests“, verwirft aber die gelesenen Nachrichten (`_ = bisEndeGelesen`); `TestReplayEndeLiestNichtsMehr` nennt „ein Terminate danach ändert nichts“, ohne dass der Test es unterscheidet (der Use Case wird so oder so einmal gefragt). | `AGENTS.md` §3.11 | `internal/adapters/driving/pgwire/server_ende_test.go` · `kein Senden des Rests, kein Sent` | ja — Mutant, der Nachrichten nach der Abbildung sendet, bleibt bei dieser Liste grün | Kommentar weiter als Prüfung |
| F-5 | INFO | Die Messung nach §3.14 liegt bei 897 (Summe aus Einfügungen und Löschungen, nachgerechnet gegen `13c23a4..a65e8d7` über `internal test docs/user`), über der Meldeschwelle 800 und der Schätzung 500 bis 850. Die Entscheidung „der Planner entschied: weiter“ steht nur als Satz in §7 *Belege des Implementers*, ohne eigenen Plan-Commit des Planners; nach der Closure ist es der dritte Datenpunkt zu `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`. Zuständig: Planner. | `AGENTS.md` §3.14 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-replay.md` · `Meldeschwelle 800 überschritten und gemeldet` | ja — `git diff --numstat` | Meldeschwelle ohne Artefakt |
| F-6 | INFO | Zwei Mutanten sind aus dem Diff heraus äquivalent oder unerreichbar und bleiben grün: `min(c.pos, len(...))` (der Cursor überschreitet die Länge nicht) und `len(responses) > 0`/`len(out) > 1` im Adapter (der Service liefert ein Ende nie ohne `error_response` davor, `PGR-E3003`). Kein Handlungsbedarf, wenn das Format die Form hält. | Maintainability | `internal/hexagon/services/replay.go` · `min(c.pos, len(c.session.Interactions))` | ja | äquivalente Mutanten |
| F-7 | INFO | `send` lässt `connection_closed` als letzte Antwort für jeden Aufrufer aus, auch im Record-Pfad (Beobachtung 2 des Implementers); dort liefert niemand sie vor `slice-v1-abschluss-verbindungsende-record`, ein Test des Record-Pfads dazu liegt nicht vor. Zuständig: der Slice `-record`, dessen §1 die Abbildung als gegeben nennt. | Maintainability | `internal/adapters/driving/pgwire/server.go` · `connection_closed als letzte` | nein | Gemeinsamer Code, Test im anderen Pfad folgt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services/replay.go` (Ende der einfachen Interaktion, `ohneLebendpruefungen`, `anfrageMitEnde`, Reihenfolge `PGR-E5003` vor `PGR-E5001` vor dem Ende, `first-request` nur erste Interaktion der ersten freien Session, `connection` zugeordnet) | geprüft, Mutanten rot bis auf F-1, F-2, F-6 |
| `internal/hexagon/ports/driving/replay.go` (Kommentar) | geprüft, ohne Befund; Zusage durch Adapter- und Service-Test getragen |
| `internal/adapters/driving/pgwire/server.go` (`send`, `replayZustellen`, `endetMitEnde`; Auslassung nur als letzte, `Sent` vor Schließen, Schließen für `Query` und `Sync`, kein weiteres Lesen) | geprüft, ohne Befund außer F-4 und F-7 |
| `internal/adapters/driving/pgwire/*_test.go`, `internal/hexagon/services/replay_ende_test.go` | geprüft; Fixtures je Test isoliert, Befund F-4 |
| `test/integration` (`replay_ende_e2e_test.go`, `testdata/ende-*.yaml`) | geprüft; Befund F-3 |
| `spec/spezifikation.md` (Architect-Commit `94a3a24`), Hexagon-Richtung (keine neuen Imports im Kern), Commit-Form (`git mv` rein in `ac2950c`) | geprüft, ohne Befund |
| `docs/user/abdeckung-*.md` | geprüft, ohne Befund; `make docs-check` grün |
| Randformen §3.12 (Stelle von `connection_closed`, Reihenfolge der Gründe, `PANIC`/`FATAL`, `Terminate` davor, Schreibfehler) | geprüft; bis auf F-1 in §6 genannt und in der Spezifikation entschieden |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Aufzählung/Ortsangabe: Mutation je Seite fehlt · Zusage ohne Mutation (Randzweig) · Plan-Zeile weiter als der Diff · Kommentar weiter als Prüfung · Meldeschwelle ohne Artefakt

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM) blockiert bis zum Fall und zur Entscheidung der Randform in §6 durch den Architect; F-2 bis F-4 sind LOW, F-5 bis F-7 INFO ohne Blockade.

**Übergabe:** Findings gehen an den Implementer (F-1 zuerst an den Architect wegen der Randform); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
