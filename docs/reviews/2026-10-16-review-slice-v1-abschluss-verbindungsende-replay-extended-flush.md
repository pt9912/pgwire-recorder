# Review-Report: slice-v1-abschluss-verbindungsende-replay-extended-flush — 2026-10-16

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-replay-extended-flush` (Schritt 37c), Diff `e982c7be4073c4e54be8c04a50125389af5ee660..3b88f4a`.

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

- Slice-Plan `slice-v1-abschluss-verbindungsende-replay-extended-flush` (§1, §3, §4, §6, §7 *Belege des Implementers*, §8)
- Geber `slice-v1-abschluss-verbindungsende-replay-extended` (in `done/`, Tests und Report des Gebers)
- Spezifikation: `LH-FA-25.a` (*Replay*, *Form*, *Stelle von `connection_closed` im Replay*, *Reihenfolge der Gründe im Replay*, *Herunterfahren*), `LH-FA-18.a` (*Interaktion*, *Nicht unterstützt*, *Replay*; Architect-Commits `ef6ccc9`, `9ef7dfb`, `0f89aee`), `LH-FA-03.b`, `LH-FA-12.a`, `LH-FA-13.a`; `SPEC-001` mit `validateAbschluss` (`internal/hexagon/model/extended.go`)
- `AGENTS.md` §3.7 bis §3.14 (Hard Rules)
- Eigene Läufe: Mutanten in Kopien des Baums (Image `deps`, `--network=none`, Arbeitsbaum unberührt) — Service 19, Adapter 3 (ein vierter hängt bis zum Zeitlimit des Pakets: `Terminate` ignoriert, Geber-Pfad), Binary ohne PostgreSQL 2; Produktcode im Diff: `git diff --stat e982c7be4073c4e54be8c04a50125389af5ee660..3b88f4a -- internal ':!*_test.go'` ist leer (bestätigt, 0 Zeilen); Messung des Plans nachgerechnet: `git diff --numstat` über `internal test docs/user` = 342 (Einfügungen 306, Löschungen 36) an `245d245` wie an `3b88f4a`, Höchststand 377 an `eaec1f3`; `make docs-check` 0 Befunde vor dem Report

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Der Kommentar von `mitten()` nennt die Grenze der Extended-Interaktion „vor ihrem Sync (`LH-FA-18.a`)“. `LH-FA-18.a` sagt seit `0f89aee` „vor ihrem Abschluss“, der bei einer Interaktion mit Ende und `Flush` dieses `Flush` ist; die Funktion entscheidet das für die `Flush`-Gruppe, und die Tests dieses Diffs (Abweichung, Herunterfahren) belegen es. Der Kommentar verweist auf die Stelle, deren Wortlaut er nicht mehr trägt. Zuständig: Architect oder Planner entscheiden das Wort (der Implementer nennt es in §7 *Beobachtungen* (3) als offen), der Implementer ändert den Kommentar; keine Logik, Plan §3 zieht nach (`AGENTS.md` §3.9). | `AGENTS.md` §3.7, §3.11 | `internal/hexagon/services/replay.go` · "nach ihrer ersten Nachricht und vor ihrem Sync (LH-FA-18.a)" | nein — kein Gate prüft Kommentarprosa | Kommentar nennt Sync, Abschluss kann Flush sein |
| F-2 | LOW | Der Kopf von `replaySitzung` sagt, bei laufender Extended-Interaktion beantworte die Sitzung ihre Nachrichten „bis zu ihrem Sync (`LH-FA-13.a`)“. Bei einer Interaktion, die mit `Flush` und Ende schließt, endet die Sitzung mit dem Senden des Endes, nicht an einem `Sync`; getestet ist das Verhalten durch `TestReplayEndeGruppeDavorSchliesstNicht` (Lauf `flush`). Die Leseregel in der Spezifikation („Wo diese Spezifikation „vor deren `Sync`“ oder „vor dem `Sync` der Interaktion“ sagt …“) nennt „bis zu ihrem `Sync`“ nicht; der Wortlaut der Spezifikation und der Kommentar stehen gleich, weder Stelle noch Kommentar folgen der Neufassung des Abschlusses. Zuständig: Architect (Lesart in der Spezifikation), danach Implementer (Kommentar). | `AGENTS.md` §3.7, §3.11 | `internal/adapters/driving/pgwire/server.go` · "beantwortet ihre Nachrichten bis zu ihrem Sync" | nein — kein Gate prüft Kommentarprosa | Kommentar nennt Sync, Abschluss kann Flush sein |
| F-3 | INFO | Der Mutant M17c (`endetEinfach` aus `anfrageMitEnde` entfernt) färbt nur Tests des Gebers (selbst gefahren: `TestReplayEndeExtendedLebendpruefungVorDerGruppe`, `TestReplayEndeExtendedIstKeineEinfacheAnfrage`, `TestReplayExtendedAmCursor`), keinen Test dieses Slice. Die Begründung des Implementers trägt: `endetEinfach` liest `Responses`, die eine Extended-Interaktion nicht trägt, die Art des Abschlusses wird dort nie erreicht; ein Mutant, der nur die `Flush`-Gruppe bricht, ist an dieser Stelle nicht möglich. Zuständig: Implementer, ohne erwartete Aktion. | `AGENTS.md` §3.10 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-replay-extended-flush.md` · "Nicht gespiegelt, über die Schnittstelle fangbar: M17c" | ja — Mutant `&& endetEinfach(erwartet[0])` entfernt, `go test` über das Paket | Nicht gespiegelter Fall mit tragender Begründung |
| F-4 | INFO | Das akzeptierte Negativ „`Terminate` und Abbruch mitten in einer `Flush`-Gruppe am Adapter ohne Fall“ trägt: Es benennt einen Test (`TestReplayEndeTerminateMittenInDerGruppe`, Geber), `replayAntwort` behandelt `Terminate` und den Lesefehler vor der abschließenden Nachricht ohne Blick auf deren Art, und die einzige Stelle, an der die Art in den Adapter eingeht (Abbildung `pgproto3.Flush` auf `ClientFlush`), färbt der Lauf `flush` (selbst gefahren: Mutant `Flush` als `Sync` abgebildet rot in `/flush`, `/sync`, `TestExtendedFlush`, `TestToClientMessage`). Dasselbe gilt für die Randform „`Sync` nach dem `Flush` mit Ende, der Adapter liest nichts mehr“ aus §6: Der Pfad ist derselbe wie beim Geber (`replayZustellen` → `folgeEnde`), `TestReplayEndeLiestNichtsMehr` fährt ihn nur mit einfachen Anfragen, ein Fall mit einer Extended-Nachricht nach dem Ende ist nicht geschrieben. Zuständig: Implementer/Planner, ohne erwartete Aktion. | `AGENTS.md` §3.10, §3.11 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-replay-extended-flush.md` · "kein Mutant bricht nur die `Flush`-Gruppe" | nein — Aussage über die Menge möglicher Mutanten, kein Gate | Akzeptiertes Negativ mit Begründung statt Fall |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services/replay.go`, `internal/adapters/driving/pgwire/server.go` (Produktcode, 0 Zeilen laut Plan) | geprüft: im Diff unverändert (`git diff --stat` leer); Mutanten dort gefahren (siehe unten), kein grüner Mutant; Befund nur zu zwei Kommentaren (F-1, F-2) |
| `internal/hexagon/services/replay_ende_extended_flush_test.go` (neu, sieben Tests) | geprüft; selbst gefahren, jeweils rot in den `Flush`-Tests: `mitten` auf `c.gruppe > 0` und auf `c.nachricht > 0` verkürzt (M17f, M17e), Lieferung bei der ersten Nachricht der Gruppe (M3), `Sync` statt `Flush` akzeptiert (nur `…FlushAbweichung` und Geber-Abweichung), `Flush` statt `Sync` akzeptiert (nur Geber-Abweichung, erwartet), `Shutdown` auf `c.gruppe == 0` verkürzt (M18b), `Forced` ohne Zweig `mitten` und ohne Zweig `ungesendet` (M20), Cursor nach der Gruppe davor und vor dem abschließenden `Flush` als verbraucht gezählt (T2, T4), `connection_closed` an die Gruppe davor gehängt (M5b); Konjunktions-Glieder von `mitten` und `Shutdown`, Aufzählungen je Element (Verbrauch je Seite, Lebendprüfung mitten in und zwischen den Gruppen, `first-request`, `connection`), Abweichung je Stelle, Herunterfahren-Stände: ohne Befund |
| `internal/hexagon/services/replay_ende_extended_test.go` (Parameter `abschluss` im Geber) | geprüft; die Fälle des Gebers sind mit `mitSync()` unverändert scharf: `sendeBis`, `endeNachrichten()[2:]` und `syncNachricht()` sind dieselben Werte wie zuvor (`ab.nachrichten()[4]` ist `syncNachricht()`), Aussagen und Fehlertexte unverändert, die Mutanten M17e/M17f/M18b/M20/T2/T4/M5b färben die `Sync`-Läufe des Gebers weiterhin; nichts verloren, ohne Befund |
| `internal/adapters/driving/pgwire/server_ende_test.go` (Lauf `flush`, Fake `gruppenReplayer`) | geprüft; Mutanten „jede gelieferte Folge schließt“ (MA1), „`Sent` fehlt“ (MA3), „`Flush` als `Sync` abgebildet“ (MA6) rot in `…/flush`; der Fake liefert im Lauf `sync` wie zuvor (Zähler `flushs` ändert nur den Lauf `flush`); Befund nur F-4 (INFO) |
| `test/integration/replay_ende_extended_e2e_test.go`, `testdata/ende-extended-flush.yaml` (Binary) | geprüft; eine Fixture und eine Einbettung je Fall, nicht mit `…Sync` geteilt; Mutanten „Adapter schließt nicht“ (E2a, Zeitüberschreitung nach 10 s) und „`Flush` als `Sync` abgebildet“ (E2c) rot in `TestE2EReplayEndeExtendedFlush` und `…Sync`, unmutiert beide grün (ohne PostgreSQL gefahren); ohne Befund |
| `spec/spezifikation.md` (`LH-FA-18.a`, `LH-FA-25.a`; Commits `ef6ccc9`, `9ef7dfb`, `0f89aee`) | geprüft gegen `validateAbschluss` (die letzte Gruppe schließt mit `Sync`, bei Ende auch mit `Flush`, jede andere mit `Flush`: genau die Fälle des Diffs); Leseregel „vor ihrem Abschluss“ trägt für „vor deren `Sync`“ und „vor dem `Sync` der Interaktion“, nicht für „bis zu ihrem `Sync`“ (siehe F-2) |
| `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-replay-extended-flush.md` (§1, §3, §4, §6, §7, §8) | geprüft; Plan folgt dem Diff (Produktcode 0, `replay.go` und `server.go` „unverändert“, zwei Liefer-Punkte, Messung 342 nachgerechnet, Marke 500 und Meldeschwelle 400 nicht erreicht, Höchststand 377 nachgerechnet, Zeilen je E2E-Fall 27 + 20 + 3 = 50); die Randformen aus §6 sind vor dem Code entschieden und nicht erst im Code-Commit eingetragen; Abdeckung, Adresse und Nachzählung: kein neuer Nehmer im Diff; ohne Befund |
| `docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md` | geprüft, Aussagen nur so weit wie der Test: „an jeder Stelle“ der Abweichung ist durch fünf Zeilen für die fünf Nachrichtenstellen gedeckt, `Terminate` „vor, in und nach der Gruppe davor und mitten in der Gruppe mit Ende“ durch die Stände 0 bis 4, „erst mit dem abschließenden `Flush`“ durch die Prüfung auf keine Antwort bei den Nachrichten davor; ohne Befund |
| `docs/plan/planning/in-progress/roadmap.md` | geprüft: Zustandsfeld „In Arbeit“ nennt Zustand und Verantwortlichen ohne Chronik; ohne Befund |
| Commit-Folge `5a8b97a`, `2db9d56` | geprüft: Der Move `2db9d56` ist ein reiner Rename ohne Inhaltsänderung (`AGENTS.md` §3.3), die Basis der Marke steht vor dem Anspruch; ohne Befund |

Der Arbeitsbaum war nach den Läufen sauber (alle Mutanten liefen in Kopien unter dem Scratchpad); `make gates` fuhr dieser Lauf nicht (Verifier-Aufgabe), `make docs-check` vor dem Report: 0 Befunde.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Kommentar nennt Sync, Abschluss kann Flush sein · Nicht gespiegelter Fall mit tragender Begründung · Akzeptiertes Negativ mit Begründung statt Fall

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. F-1 und F-2 sind Kommentare an unveränderter Logik (Produktcode bleibt 0 Zeilen Logik); ein grüner Mutant wurde in 24 eigenen Läufen nicht gefunden. Der Wortlaut „Abschluss“ statt „Sync“ liegt beim Architect, die Kommentarkorrektur ist minimal und zieht Plan §3 nach (`AGENTS.md` §3.9).

**Übergabe:** F-1 und F-2 an den Architect (Wortlaut, Lesart für „bis zu ihrem `Sync`“) und danach an den Implementer (Kommentare in `replay.go` und `server.go`, Plan §3), F-3 und F-4 zur Kenntnis; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
