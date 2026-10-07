# Review-Report: slice-lint-bestand-driving — 2026-10-07

**Review-Art:** Code-Review (Umbau ohne Verhaltensänderung, Kontexte, Kommentare, Belege in §7). Geprüft wird gegen Plan §1, §3 und §6 (Stand Architect `b7d4555`), `SPEC-049`, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 5 und die Hard Rules (Modul 10). Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** `git diff b7d4555..76ce087`:

- `78f28a1`: Kontexte und Receiver (Befunde 9 bis 16), `internal/adapters/driving/pgwire/server.go`, `internal/bootstrap/bootstrap.go`, `internal/adapters/driving/cli/cli.go`.
- `e76b25f`: Komplexität (Befunde 1 bis 8), nur `server.go`.
- `76ce087`: die Belege des Implementers in §7.

**Skill:** `.harness/skills/reviewer.md` am Stand `76ce087` (mit MEDIUM *Randform im Code-Commit* und *Adresse nimmt nicht an*)
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-lint-bestand-driving.md`: ganz gelesen am Stand `76ce087`. Der Plan-Diff im Gegenstand ist ein einziger Hunk, angehängt in §7. §7 *Belege des Implementers* habe ich gelesen, den Bericht des Implementers nicht.
- `docs/plan/planning/open/slice-tests-ueberlebende-mutanten.md` §1 (Gegenstand, *Sammelregel*, *Ausdrücklich NICHT*), §2 und §4.
- `spec/spezifikation.md` `SPEC-049`, [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes).
- [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) (Entscheidung 5), [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md).
- `AGENTS.md` §3.3, §3.6, §3.7, §3.9 bis §3.13.
- `server.go` an `b7d4555` und `e76b25f` gegeneinander gelesen, dazu `TestReplayHerunterfahrenSpaeteFrist` und die Quelle von `net.Listen` und `(*net.ListenConfig).Listen` im gepinnten Go-Image (`/usr/local/go/src/net/dial.go`).
- Vorherige Reports: das Review zu `slice-lint-bestand-kern-driven` (F-459 bis F-464), dazu F-401 als Präzedenz für eine Kommentar-Zusage ohne Test. Die höchste vergebene Nummer vor diesem Lauf war F-464.

**Ausgeführte Läufe:**

- Je Mutant eine frische Kopie aus `git archive` (`76ce087` bzw. `b7d4555`) unter dem Scratch-Pfad, kein `cp -p`. Die Ersetzung lief per Skript mit Prüfung der Trefferzahl, danach `touch` auf die Datei. Unit-Tests liefen ohne Netz per Bind-Mount (nur lesend) in einem eigenen Image der Stufe `deps` mit eigenem Tag, `go vet` und `go test -count=1`.
- Integrationsläufe wie `make test-integration` mit beiden Phasen, aber mit eigenem Image-Tag, eigenem internen Netz, eigenem PostgreSQL-Container (gepinntes Image aus `harness/mk/integration.mk`) und eigenem Volume. Danach wurden Container, Netz, Volume, Image und Kopie entfernt. Ein hostweites Docker-Aufräumen gab es nicht.
- Stufe `lint` aus einer Kopie von `76ce087` mit eigenem Tag: `0 issues.`
- Ohne Mutation sind alle Unit-Tests der Kopie grün (Lauf X0).
- `make docs-check` vor dem Commit (siehe Ende).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-465 | MEDIUM | Zwei neue Doc-Kommentare sagen ein Verhalten zu, das kein Test hält. `replayLesefehler`: „Jeder andere Lesefehler geht als PGR-E6001 an den Client.“ Der Satz ist neu im Diff, und R3 (`s.fail` entfernt) ist laut §7 grün; ich habe das in Unit und Integration nachgefahren, beide grün. Im selben Kommentar steht der verschobene Satz „Ein Ende der Client-Verbindung ist im Replay regulär“; R2 dazu ist laut §7 ebenfalls grün. `replayZustellen`: „scheitert das Schreiben, merkt es den Fehler“. Mein Mutant X4 (ohne `s.sendFailed(err)`) ist in Unit und Integration grün (F-466). §7 ordnet G2 und G3 als Lücken des Bestands ein, nimmt aber keinen der beiden Wege aus §3.11 für die Kommentare. *Failure-Szenario:* Eine spätere Änderung lässt die Fehlerantwort oder das Merken weg, `make gates` bleibt grün, und der Kommentar sagt weiter zu, was niemand prüft. Wer die Stelle liest, hält sie für geprüft. Dasselbe Muster war schon zweimal LOW (F-401, F-460), deshalb MEDIUM nach Reviewer-Skill. | `AGENTS.md` §3.11; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Reviewer-Skill MEDIUM „Wiederholung eines Musters, das schon zweimal LOW war“ | `internal/adapters/driving/pgwire/server.go` · „Jeder andere Lesefehler geht als PGR-E6001 an den Client.“ (Zeile 234, Lesehilfe); · „scheitert das Schreiben, merkt es den Fehler“ (Zeile 272) | ja (Mutanten R3, X4: Unit und Integration grün) | Kommentar-Zusage ohne Test, verschoben statt gefasst |
| F-466 | LOW | Die Liste grüner Mutanten in §7 ist unvollständig. Ohne `s.sendFailed(err)` in `replayZustellen` (X4) bleiben Unit-Tests und beide Integrationsphasen an `76ce087` grün. Am Stand vor dem Umbau (`b7d4555`) ist der gleiche Mutant an beiden Stellen in `replaySitzung` ebenfalls grün (X4vor, Unit). Er ändert das Verhalten: Ein Schreibfehler an den Client im Replay wird nicht als PGR-E4003 gemerkt, und der Lauf endet mit Exit-Code 0. Fangbar ist er über die Schnittstelle, mit `pgwire.Handle` über eine Verbindung, deren Schreiben scheitert, und `FirstErrorCode()`. Der Plan verlangt je umgebauter Funktion eine Mutation, nicht jeden Zweig. Gefunden hat ihn der Implementer deshalb nicht pflichtwidrig, aber die Sammlung für den Planner fehlt ihm. | Plan §6 Risiko *Verhalten von `replaySitzung` ändert sich unbemerkt*; `slice-tests-ueberlebende-mutanten` §1 *Sammelregel* | `internal/adapters/driving/pgwire/server.go` · `func (s *Server) replayZustellen(`; `docs/plan/planning/in-progress/slice-lint-bestand-driving.md` · „*Grüne Mutanten, eingeordnet*“ | ja (X4, X4vor) | Liste grüner Mutanten unvollständig |
| F-467 | LOW | Die Einordnung von G2 sagt zu R4: „die Sitzung endet nicht, `closeReplay` läuft nicht“. Nach dem Code stimmt das nur, solange `ctx` läuft. `closeReplay` ist per `defer` registriert. Unter R4 bleibt der Ausgang über den Shutdown-Zweig erhalten: Endet `ctx` und gibt `Shutdown` das Ende frei, kehrt `replaySitzung` zurück und `closeReplay` läuft. Bis dahin liest die Schleife nach dem Lesefehler erneut. Die Test-Idee („`Handle` kehrt binnen einer Frist zurück“, ohne Abbruch) bleibt richtig. Die Grenze, die der Planner in den Nehmer überträgt, ist aber weiter gefasst, als der Code sie trägt. | `AGENTS.md` §3.11 (Plan-Zeile) | `docs/plan/planning/in-progress/slice-lint-bestand-driving.md` · „nicht, `closeReplay` läuft nicht“ | ja (Lesen; Mutant R4 mit Abbruch) | Zusage im Plan weiter als der Code |
| F-468 | INFO | Mutant W ist rot, aber nur über die Zeitüberschreitung des Pakets. Nachgefahren mit `-timeout 40s`: `panic: test timed out`, „running tests: TestReplayHerunterfahrenSpaeteFrist“. Der Test blockiert in `server_extended_test.go:846` bei `<-conn.gesetzt`, ohne eigene Frist. In `make test` (ohne `-timeout`) sind das 10 Minuten, und statt einer Zusicherung kommt ein Goroutinen-Dump. Faktisch grün ist er nicht, denn das Rot ist deterministisch und nennt den richtigen Test. Teuer ist er aber, und wie ein Mutations-Gate mit einer Frist je Mutant ihn zählt, ist offen. Der Test liegt außerhalb dieses Slice (§1 schließt geänderte Tests aus). `replaySitzung` hält mit R1 unabhängig davon eine rote Mutation. Für den Architect bzw. die Randformen von `slice-harness-mutation`. | Plan §1 *Ausdrücklich NICHT* (Tests); Maintainability | `internal/adapters/driving/pgwire/server_extended_test.go` · `<-conn.gesetzt` | ja (W mit `-timeout 40s`) | Mutant nur über Paket-Zeitüberschreitung rot |
| F-469 | INFO | Zur Adresse der grünen Mutanten G1 bis G4 (und X4 aus F-466). Dass §7 keine Adresse nennt, folgt `AGENTS.md` §3.13. Der Nehmer `slice-tests-ueberlebende-mutanten` schließt in §1 „Der PGWire-Adapter (`internal/adapters/driving/pgwire`) und `test/integration`“ aus. Damit trifft ein Ausschluss die Sendung, er ist keine Adresse, und der Punkt geht an den Planner. Der Diff nennt keinen Slice neu als Adresse. Die Nennung in §6 stammt aus `b7d4555` und sieht den Schnitt selbst vor. Kein MEDIUM *Adresse nimmt nicht an*. Für den Planner, drei Punkte. (a) Getragen wird die Übergabe nur von der Zeile „**Nehmer: offen**“ in §7. Bindend wird sie bei der Closure dieses Slice über die Paarung *Folge-Slice*. (b) Der Satz im Nehmer „seine Funde nimmt dieser Slice nach dieser Regel an“ steht neben dem Ausschluss. Er gilt nur über den Schnitt der *Sammelregel*. (c) G4 beantwortet einen offenen Punkt des Nehmers: In der Grenze von P8/P2 nennt er den Rückweg im PGWire-Adapter „nicht gemessen“. Jetzt ist er gemessen, T3 und T5 überleben dort. | `AGENTS.md` §3.13; `slice-tests-ueberlebende-mutanten` §1 | `docs/plan/planning/in-progress/slice-lint-bestand-driving.md` · „**Nehmer: offen**“; `docs/plan/planning/open/slice-tests-ueberlebende-mutanten.md` · „seine Funde nimmt dieser Slice nach dieser Regel an“ | nein (Lesen) | Übergabe an Planner nur in §7 getragen |
| F-470 | INFO | Zur Gleichwertigkeit, für den Verifier. **Kontexte:** Wo vorher `r.ctx` stand, steht jetzt überall ein `context.WithoutCancel(ctx)` über demselben `ctx`. Es gibt zwei Instanzen (`serverRichtung` und `uc` in `clientRichtung`) statt einer, aber kein Code vergleicht Kontexte. Werte, Abbruch und Frist sind gleich. `clientRichtung` fragt `Err()` allein auf dem ungelösten `ctx` und reicht `uc` an `Shutdown`, an `beende` in jedem Zweig und an `nachricht`. Der Wächter in `recordSitzung` und der in `replayWaechter` sehen denselben `ctx` wie vorher. Die Goroutine von `serverRichtung` bildet den gelösten Kontext selbst, statt ihn aus dem Struct zu lesen; der Unterschied ist nicht beobachtbar. K1 und mein X7 (`EndShutdown` → `EndClosed`) sind rot. In `nachricht`, `fehler`, `wartenAufEnde`, `schreibe`, `beende` und `serverRichtung` heißt der Parameter `ctx`, ist aber immer gelöst; heute fragt dort niemand `Err()`. **`Listen`:** `net.Listen` ist im gepinnten Go-Image wörtlich `var lc ListenConfig; return lc.Listen(context.Background(), network, address)`. Der Diff nutzt dasselbe Nullwert-`ListenConfig` (kein `KeepAlive`, `Control` oder MPTCP-Unterschied) und Netzwerk `"tcp"`, also Dual-Stack wie vorher. Der Unterschied liegt nur in den Werten des Kontexts, und die liest `net` nur über eigene interne Schlüssel. **Herausgelöste Funktionen:** Reihenfolge der Fälle, Zweige, Meldungscodes und Fehlertexte sind gleich. Das gilt für `replayAntwort`/`nachricht` (Fall-Reihenfolge wie vorher, `continue` → `true`), `startkopf` (Länge vor Startcode, S1 rot) und `toMessage`/`toExtendedMessage`: `ReadyForQuery` wandert, aber der Switch geht über disjunkte Konstanten, und der Text „Antworttyp %q ohne Abbildung“ bleibt (`TestToMessageErfindetNichts`). `weiterlesen()` wird weiter aufgerufen, sein Wert verworfen. **`replaySitzung` mit `gocognit` 17** ist eine echte Zerlegung: Wächter, Klassifikation der Lesefehler, Antwort auf eine Nachricht und Zustellung sind eigene Funktionen. In der Schleife bleibt nur das Herunterfahren mit seiner Reihenfolge (Shutdown fragen, auf `geweckt` warten, Frist zurücksetzen), und die gehört zusammen. Eigene rote Mutanten: X5, X6, X7, X8 (Tabelle unten). | Plan §6 *Kontexte*, *`pgwire.Listen` mit Kontext*, *Komplexitäts-Bereinigung*; [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 5 | `internal/adapters/driving/pgwire/server.go` · `uc := context.WithoutCancel(ctx)`; · `(&net.ListenConfig{}).Listen(context.WithoutCancel(ctx), "tcp", address)` | ja (Mutanten, Lesen) | — (Negativbefund) |

**Eigene Mutationen** (jede in einer frischen Kopie; ohne Mutation sind alle genannten Tests grün):

| ID | Stand | Mutation | Lauf | Ergebnis |
|---|---|---|---|---|
| X4 | `76ce087` | `replayZustellen` ohne `s.sendFailed(err)` | Unit, beide Integrationsphasen | grün (F-465, F-466) |
| X4vor | `b7d4555` | `replaySitzung` ohne `s.sendFailed(err)` an beiden Stellen (2 Treffer) | Unit `pgwire` | grün (F-466) |
| X5 | `76ce087` | `nachricht`: Extended-Fall `return true` → `return false` | Unit | rot, u. a. `TestExtendedSyncGruppe`, `TestExtendedHerunterfahren`, `TestWeiterlesenNachWecken` |
| X6 | `76ce087` | `startkopf`: bekannter Startcode `return true` → `return false` | Unit | rot, u. a. `TestReplayExtended`, `TestExtendedHerunterfahren` |
| X7 | `76ce087` | `clientRichtung`: `beende(uc, model.EndShutdown, nil)` → `EndClosed` | Unit | rot, `TestExtendedHerunterfahren`, `TestHerunterfahrenWeckenNichtVerloren` |
| X8 | `76ce087` | `replayZustellen` ohne `s.replayer.Sent(ctx, id)` | Unit | rot, `TestReplaySent` |
| R3 | `76ce087` | `replayLesefehler` ohne `s.fail(…)` (nachgefahren aus §7) | Unit, beide Integrationsphasen | grün (bestätigt §7, F-465) |
| W | `76ce087` | `replayWaechter` ohne `SetReadDeadline(time.Now())` (nachgefahren aus §7) | Unit `pgwire`, `-timeout 40s` | rot nur über Zeitüberschreitung (F-468) |

Nicht nachgefahren habe ich G1 (temporärer Test des Implementers), R2, R4 und die T-Mutanten. Ihre Einordnung habe ich am Code gelesen. Bei R4 weicht sie ab (F-467).

---

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driving/pgwire/` | geprüft. Befunde F-465 bis F-468 und F-470. Keine neue Ausnahme, kein `//nolint`, kein neuer Import. |
| `internal/bootstrap/` | geprüft, ohne Befund. `Listen(ctx, …)` in `record` und `replay`; `Finish(context.WithoutCancel(ctx))` hat wie `context.Background()` weder Abbruch noch Frist (K2 rot). |
| `internal/adapters/driving/cli/` | geprüft, ohne Befund. Unbenannter Receiver, die Methodenmenge von `wahrheitswert` ist gleich (K5 rot). |
| `.golangci.yml`, Testdateien, `docs/user/` | geprüft, ohne Befund. Im Diff unverändert. |
| `docs/plan/planning/in-progress/slice-lint-bestand-driving.md` | geprüft. Nur §7 ist angehängt (ein Hunk). Befunde F-466, F-467, F-469. |
| Hard Rule 3.3 | geprüft, ohne Befund. Im Diff gibt es keinen Move. |
| Hard Rule 3.6 | geprüft, ohne Befund. Keine Schwelle gesenkt, das Profil ist unverändert. Die Stufe `lint` meldet `0 issues.` |
| Hard Rule 3.7 (Kommentar-Klassen) | geprüft, ohne Befund. Die neuen Doc-Kommentare beschreiben den Ist-Zustand, keiner nennt eine verworfene Alternative oder einen abwesenden Text. Zu §3.11 siehe F-465. |
| Hard Rule 3.9 | geprüft, ohne Befund. Die Code-Commits brauchten keine Planänderung: Die Befunde sind dieselben 16 wie in §1, und die Kontexte sind so umgesetzt, wie §6 sie entscheidet. |
| Hard Rule 3.12 / MEDIUM *Randform im Code-Commit* | geprüft, ohne Befund. §6 ist in keinem Commit des Diffs geändert. Die Randformen hat der Architect in `b7d4555` vor dem Code entschieden, und der Diff entscheidet keine weitere. |
| Hard Rule 3.13 / MEDIUM *Adresse nimmt nicht an* | geprüft, ohne Befund. Der Diff nennt keinen Slice neu als Adresse (F-469). |
| Core-Reinheit, [ADR-0001](../plan/adr/0001-hexagonale-architektur.md), [ADR-0010](../plan/adr/0010-verwendung-von-pgproto3.md) | geprüft, ohne Befund. `internal/hexagon/` ist unberührt, `pgproto3` bleibt im PGWire-Adapter. Die geänderte exportierte Signatur `pgwire.Listen` ruft nur `internal/bootstrap`. |
| Commit-Messages `78f28a1`, `e76b25f`, `76ce087` | geprüft, ohne Befund. Jede nennt `slice-lint-bestand-driving`, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) und `LH-QA-07`, keine nennt eine Struktur-ID. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

**Summary:** 0 HIGH · 1 MEDIUM · 2 LOW · 3 INFO (F-465 neue Doc-Kommentare von `replayLesefehler` und `replayZustellen` sagen ungeprüft zu; F-466 grüner Mutant `sendFailed` in `replayZustellen` fehlt in der Liste; F-467 Grenze von G2/R4 weiter als der Code; F-468 Mutant W nur über Paket-Zeitüberschreitung rot; F-469 Übergabe der grünen Mutanten an den Planner nur in §7 getragen, nach §3.13 korrekt; F-470 Umbau und Kontexte gleichwertig). Wiederkehrende Klassen: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-465, drittes Mal nach F-401 und F-460); Liste grüner Mutanten unvollständig (F-466, wie F-461).

**Finding-Klassen dieses Laufs:** Kommentar-Zusage ohne Test, verschoben statt gefasst · Liste grüner Mutanten unvollständig · Zusage im Plan weiter als der Code · Mutant nur über Paket-Zeitüberschreitung rot · Übergabe an Planner nur in §7 getragen

## Verdikt

**Merge-blockierend:** nein. Der Umbau verhält sich wie `b7d4555`, und das Herunterfahren hängt weiter allein an `ctx.Err()` in `clientRichtung` und an den Wächtern. Die Stufe `lint` meldet modulweit nichts. F-465 braucht vor der Closure einen Ausgang, weil der Kommentar sonst eine ungeprüfte Zusage in den Bestand einführt, den `slice-harness-lint` dann als bereinigt anschließt.

**Übergabe:** F-465 und F-467 gehen an den Implementer. F-466 und F-469 gehen an den Planner, der den Nehmer der grünen Mutanten nachträgt. F-468 geht an den Architect. F-470 geht an den Verifier. Dieser Report ersetzt keine Verifikation.
