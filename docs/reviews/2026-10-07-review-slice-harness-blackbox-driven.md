# Review-Report: slice-harness-blackbox-driven — 2026-10-07

**Review-Art:** Code-Review (Testumbau und eine Stelle Produkt-Code) gegen Plan, Entscheidungen (`SPEC-049` Punkt 7 und 8, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 und 5) und Hard Rules (Modul 10). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git show 1c30e43`. Die Testdateien unter `internal/adapters/driven/postgres` (2) und `internal/adapters/driven/recording` (1) sind jetzt `_test`-Pakete. Neu sind `internal/adapters/driven/postgres/export_test.go` und der unexportierte Konstruktor `newSession` in `upstream.go`. Die Entscheidungen dazu stehen in `58e863c` (`SPEC-049` Punkt 7, Plan §1, §3 und §6). Schwerpunkte laut Auftrag: `newSession` und das Verhalten von `Open`, die Brücken nach Punkt 7, verdeckt statt behoben, Prüfung im Umbau verloren, Mutationen.

**Skill:** `.harness/skills/reviewer.md` @ `de770f7` (unverändert bis `1c30e43`)
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-harness-blackbox-driven.md` (ganz, Stand `1c30e43`), dazu die Plan-Diffs von `58e863c` und `1c30e43`. §7 *Belege des Implementers* ist gelesen; den Bericht des Implementers habe ich nicht.
- `spec/spezifikation.md` `SPEC-049` Punkt 1 bis 8, dazu der Diff von `58e863c`; [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) (Accepted), [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)
- `AGENTS.md` §3.2, §3.3, §3.6, §3.7, §3.9 bis §3.12
- Produkt-Code: `internal/adapters/driven/postgres/upstream.go` an `1c30e43^` und `1c30e43`, `internal/adapters/driven/recording/yaml.go`
- Vorherige Reports: das Review zu `slice-harness-blackbox-kern` (F-441 bis F-443). Die höchste vergebene Nummer vor diesem Lauf war F-443.

**Ausgeführte Läufe:**

- `make lint`: Exit 2. Unter beiden Pfaden stehen nur die vier Befunde im Produkt-Code, die §7 nennt (`upstream.go:100`, zweimal `upstream.go:226`, `yaml.go:538`). In einer `*_test.go` gibt es keinen.
- `make kopf-check` und `make abdeckung-check`: beide Exit 0.
- Für die Testliste habe ich `1c30e43^` und `1c30e43` per `git archive` in je eine frische Kopie unter dem Scratch-Pfad entpackt. Dort liefen die folgenden Läufe in einem eigenen Image der Stufe `deps` (eigener Tag), ohne Netz. `go test -list .` ist je Paket gleich: `postgres` 14 Tests, `recording` 13. Ein `go test -v` lieferte vorher und nachher dieselben 108 Zeilen `--- PASS` für Tests und Subtests und kein `FAIL`.
- In der Kopie von `1c30e43` ist `gofmt -l` leer und `go vet` grün.
- In den geänderten Zeilen des Diffs steht keine `Abdeckung:`-Zeile.
- `go test -race` läuft im gepinnten Alpine-Image nicht, es fehlt gcc. Die Synchronisation habe ich deshalb nur lesend beurteilt (Schwerpunkt 3).
- Sieben eigene Mutationen, je in einer frischen Kopie (`cp -r` ohne `-p`) der Kopie von `1c30e43`, mit `go test -count=1 -run` (Tabelle unten). Danach habe ich die Kopien gelöscht. Das Repo blieb bis auf diese Datei unverändert.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-444 | INFO | Urteil zu `newSession`, wie der Plan es dem Review zuweist. `Open` legt die Session vor dem Startup mit `newSession(conn)` an. Den Aufbau liest es über `s.fe` und gibt an ReadyForQuery genau dieses `s` zurück. Frontend, Reihenfolge und Fehlerpfade sind dieselben wie an `1c30e43^`: Jeder Fehlerpfad ruft weiter nur `conn.Close()`, verwirft `s` und schickt kein Terminate. Neu ist nur, dass das Struct früher angelegt wird. Ein Verhalten ändert sich nicht. Meine Mutation R1 (`return newSession(conn), …`) bleibt in allen Tests von `postgres` grün, wie §7 es angibt. Dass die zurückgegebene Session das Frontend des Aufbaus trägt, prüfte auch vorher kein Test, denn das Literal `&session{conn: conn, fe: fe}` war genauso ungeprüft. Die Lücke liegt im kritischen Pfad, weil ein neues Frontend den Lesepuffer des Aufbaus verwirft. Sie ist aber nicht neu und gehört nicht zu diesem Umbau. | Plan §6 *Session über `net.Pipe`* („das ist Urteil des Review“); `SPEC-049` Punkt 7; Rolle Verifier | `internal/adapters/driven/postgres/upstream.go` · „s := newSession(conn)“ und „return s, responses, nil“ | ja (R1 grün, R6 rot) | — (Hinweis) |
| F-445 | INFO | Durch den Einschub in `SPEC-049` Punkt 7 (`58e863c`) ist ein Bezug unklar geworden. Der Satz „Was nur an einem solchen Wert zu sehen ist, prüft der Test über die exportierte Schnittstelle“ folgte bisher direkt auf den Wert eines unexportierten Typs. Jetzt steht unmittelbar davor „Werte exportierter Typen, die sie annimmt, etwa eine Verbindung, stellt der Test“. „Ein solcher Wert“ lässt sich damit auch als die Verbindung lesen. Bei dieser Lesart stünde `Verbindung`, das die `conn` einer Session herausgibt, gegen den Wortlaut, obwohl der Architect es am 2026-10-07 bestätigt hat (Plan §6). Die Zeile ist zudem nicht umbrochen. Für den Architect, der Punkt 7 führt. Gemeldet als Hinweis, weil die Bestätigung im Plan die Auslegung für diesen Slice trägt. | `SPEC-049` Punkt 7; Rolle Architect | `spec/spezifikation.md` · „Werte exportierter Typen, die sie annimmt, etwa eine Verbindung, stellt der Test. Was nur an einem solchen Wert zu sehen ist“ | nein | — (Hinweis) |
| F-446 | INFO | Nicht im Diff: An drei Stellen in `upstream_test.go` steht weiter `_ = s.Close()` an einer `driven.UpstreamSession`, darunter die Cleanup-Funktion von `oeffne`. Das ist derselbe Typ, dessen Fehler `schliesse` in `upstream_extended_test.go` jetzt prüft. `errcheck` meldet diese Stellen nicht, weil der Fehler ausdrücklich verworfen wird. Sie standen schon an `1c30e43^`, und die Regel des Plans („kein `_ =` vor einem Fehler, den `errcheck` meldet“) ist eingehalten. Folge: Meine Mutation R3 (`Close` liefert immer einen Fehler) machen nur die Tests rot, die `schliesse` nutzen. Die Tests über `oeffne` bleiben grün. Eine Aktion erwarte ich nicht. | Plan §6 Risiko *Befund verdeckt statt behoben*; `SPEC-049` Punkt 5 (`errcheck`) | `internal/adapters/driven/postgres/upstream_test.go` · „t.Cleanup(func() { _ = s.Close() })“ | ja (R3) | — (Hinweis) |
| F-447 | INFO | Zur Frage, ob die Mutationen für §3.10 reichen. Die Tabelle in §7 deckt jede der drei Brücken-Funktionen, den Helfer `schliesse` und die drei Testdaten-Funktionen ab. Für `recording` nennt sie eine Zeile. Meine Mutationen R1 bis R7 decken die Lücken ab: die Rückgabe von `Open`, den Fehlerpfad beim Aufbau, den Meldungstext nach `beendet()` und die Testfälle mit lokalem `syncGruppe`. Bis auf die zwei erwarteten Fälle waren sie rot. Erwartet grün war R1 (F-444). R7 entfernt die Lesefrist in `oeffne`; sie verhindert nur ein Hängen, und §7 sagt das auch. Nach meinem Urteil reicht die Abdeckung. Ob sie die DoD trägt, prüft der Verifier. | `AGENTS.md` §3.10; Plan §6 Risiko *Prüfung geht im Umbau verloren*; Rolle Verifier | `docs/plan/planning/in-progress/slice-harness-blackbox-driven.md` · „*Mutationen* (`AGENTS.md` §3.10“ | ja (R1 bis R7) | — (Hinweis) |

**Eigene Mutationen** (frische Kopie je Mutation, ohne Mutation grün):

| ID | Mutation | Weg im Test | Ergebnis |
|---|---|---|---|
| R1 | `Open` gibt an ReadyForQuery `newSession(conn)` statt `s` zurück | öffentliches `Open` | grün (alle Tests von `postgres`), siehe F-444 |
| R2 | `SetWriteDeadline` in `Close` entfernt | Brücke `NeueSession` auf `net.Pipe` | rot, `TestCloseMitSchreibfrist` |
| R3 | `Close` liefert nach `conn.Close()` immer `net.ErrClosed` | Helfer `schliesse` | rot, `TestSendUndReceive`, `TestReceiveEndetMitReadyForQuery`, `TestReceiveFehler/COPY`, `…/Verbindungsende` (ausgewertet sind nur die ersten Fehlzeilen); über `oeffne` grün, siehe F-446 |
| R4 | `verbindungsende` nennt die Meldung `M` nicht mehr | Testdaten `beendet()`, `pruefeAbbruch` | rot, u. a. `TestFehlerantwortVorDemAbbruch/nach_Ergebnissen,_FATAL`, `…/Schweregrad_ERROR`, `…/Schweregrad_PANIC`, `…/zwei_Fehlerantworten,_die_letzte_zählt` (ausgewertet sind nur die ersten Fehlzeilen) |
| R5 | Extended-Interaktion: `\|\|` zu `&&` bei request und responses | `recording_test`, Testfälle mit lokalem `syncGruppe` | rot, `TestUnmarshalExtendedFehler/extended_mit_request`, `…/extended_mit_responses` |
| R6 | `Open` gibt bei ErrorResponse im Aufbau `s` statt `nil` zurück | öffentliches `Open` | rot, `TestNichtVermittelbar/Fehlerantwort_im_Aufbau` |
| R7 | in `oeffne` die Lesefrist über `Verbindung` entfernt (Testcode) | — | grün, wie §7 angibt |

## Antwort auf die Schwerpunkte

1. **`newSession`.**
   - `Open` ruft `newSession(conn)` für genau die Session, die es zurückgibt, und liest den Aufbau über deren Frontend.
   - Startup-Reihenfolge, Fehlercodes und Fehlerpfade sind unverändert. Geprüft mit `git show 1c30e43^:internal/adapters/driven/postgres/upstream.go`: Der Unterschied liegt nur in den drei Zeilen `s := newSession(conn)`, `fe := s.fe` und `return s, …`.
   - `newSession` ruft sonst nur die Brücke. Es ist kein Export.
   - Urteil und Grenze stehen in F-444.
2. **Brücken nach Punkt 7.**
   - `export_test.go` gehört zum Paket `postgres` und ist dort die einzige Testdatei. Sie enthält nur drei Funktionen. Es gibt keine Variable, keinen Typ-Alias, keine Funktion `Test…`/`Benchmark…`/`Example…`/`Fuzz…` und keinen Zustand auf Paketebene.
   - `ToFrontendMessage` reicht nur weiter.
   - `NeueSession` reicht an `newSession` weiter, das auch `Open` ruft. Den `net.Conn` stellt der Test. Die Brücke legt keinen `session`-Wert selbst an.
   - `Verbindung` liest per Typ-Assertion das Feld `conn` einer übergebenen Session und setzt nichts. Fristen und `CloseWrite` setzt der Test an der Verbindung. Diese Zuordnung hat der Architect am 2026-10-07 in §6 bestätigt (zum Wortlaut siehe F-445).
   - Im Test gibt es kein `session`-Literal, kein `var` und kein `new` eines unexportierten Typs mehr. `TestCloseMitSchreibfrist` geht über `NeueSession`.
   - `recording` hat keine Brücke und braucht keine.
3. **Verdeckt statt behoben.**
   - `schliesse` prüft den Fehler von `Close` mit `t.Error`. R3 bestätigt das.
   - `bereit()`, `beendet()` und `ergebnis()` liefern bei jedem Aufruf neue Literale. Es gibt keine Closure und keinen Zustand auf Paketebene.
   - Die Server-Goroutinen bekommen ihre Werte entweder beim Start (`fakeServer(t, bereit(), …)`, das ist eine happens-before-Kante) oder erzeugen sie selbst (`range bereit()` in der Goroutine). Zwei Goroutinen teilen sich keinen Wert mehr. Vorher teilten sich Tests und Server dieselben Zeiger nur lesend.
   - Wo derselbe Fehler gesendet und erwartet wird, kommt jetzt je ein frischer Wert aus `beendet()`. `pruefeAbbruch` vergleicht die Felder `Code` und `Message`, nicht die Zeiger. Die Prüfung bleibt also gleich stark (R4).
   - Mit `syncGruppe` ist keine Prüfung weggefallen. Jede Verwendung an `1c30e43^` war eine lokale Variable gleichen Namens (Zeilen 477 und 546). Der Compiler und `unused` belegen, dass die Konstante unbenutzt war. R5 zeigt, dass die Fälle mit dem lokalen Namen weiter fangen.
   - Im Diff kommt kein `_ =` hinzu, es gibt kein `//nolint` und `.golangci.yml` ist unverändert. Zum Bestand siehe F-446.
4. **Prüfung im Umbau verloren.**
   - Stichproben gegen `1c30e43^`: `oeffne`, `schreibhaelfteZu`, `TestCloseMitSchreibfrist`, `TestFehlerantwortVorDemAbbruch`, `TestSendNachFehlerantwort`, `TestReceiveFehler/Zielart`, `TestUnmarshalExtendedFehler` und `TestUnmarshalVerweise`.
   - `oeffne` prüft den Fehler von `SetReadDeadline` jetzt mit `t.Fatal`, vorher wurde er mit `_ =` verworfen. Das ist strenger, nicht schwächer.
   - Die Session in `TestCloseMitSchreibfrist` ist feldgleich mit dem früheren Literal (`conn`, `fe` aus `NewFrontend(a, a)`).
   - Testliste und Subtests sind gleich (*Ausgeführte Läufe*). Ich habe keinen Test gefunden, der weniger prüft als vorher.
5. **Mutationen.** Siehe F-447 und die Tabelle oben.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driven/postgres/upstream_test.go`, `upstream_extended_test.go` | geprüft, ohne Befund. Paket `postgres_test`, Helfer mit `driven.UpstreamSession`, keine Globale, kein Befund von `make lint`; Hinweis F-446. |
| `internal/adapters/driven/postgres/export_test.go` | geprüft, ohne Befund nach Punkt 7; Hinweis zum Wortlaut von Punkt 7 F-445. |
| `internal/adapters/driven/postgres/upstream.go` | geprüft, ohne Befund. Eine Stelle wie in Plan §1, kein neuer Export und keine Verhaltensänderung; Urteil F-444. |
| `internal/adapters/driven/recording/yaml_test.go` | geprüft, ohne Befund. Paket `recording_test`, nur Qualifizierung und Entfernen von `syncGruppe`. |
| `spec/spezifikation.md` `SPEC-049` Punkt 7 (Diff `58e863c`) | geprüft. Hinweis F-445. |
| Hard Rule 3.2, 3.6 (Suppression, Gate-Lockerung) | geprüft, ohne Befund. Kein `//nolint`, `.golangci.yml` ist unverändert. |
| Hard Rule 3.3 | geprüft, ohne Befund. Die Commits `8ebb210` und `45dc13b` sind reine Renames (100 %). `1c30e43` enthält keinen Move. |
| Hard Rule 3.7, 3.11 (Kommentare) | geprüft, ohne Befund. Die neuen Kommentare an `newSession`, `schliesse`, `bereit`, `ergebnis` und an der Brücke beschreiben, was da ist, und sagen nichts zu, was kein Test prüft. |
| Hard Rule 3.9 (Plan folgt Korrektur) | geprüft, ohne Befund. §1, §3 und §6 stimmen mit dem Diff überein (Stelle in `upstream.go`, drei Brücken-Funktionen). `make kopf-check` ist grün. |
| Hard Rule 3.12 | geprüft, ohne Befund. `1c30e43` ändert im Plan nur §7. Die Randformen (Brücke, `newSession`, `Verbindung`) stehen seit `58e863c` in §6 und sind dort vor dem Code-Commit vom Architect entschieden. Keine Randform ist dem Review statt dem Architect zugewiesen. Dem Review zugewiesen ist das Urteil, ob `Open` `newSession` ruft, und das ist keine Randform. |
| Commit-Message `1c30e43` | geprüft, ohne Befund. Sie nennt `slice-harness-blackbox-driven`, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) und `LH-QA-07`, aber keine Struktur-ID. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 4 |

**Summary:** 0 HIGH · 0 MEDIUM · 0 LOW · 4 INFO (F-444 `Open` ruft `newSession` für die zurückgegebene Session, Verhalten unverändert, Rückgabe weiter ungeprüft wie vorher; F-445 Bezug in `SPEC-049` Punkt 7 nach dem Einschub mehrdeutig; F-446 `_ = s.Close()` im Bestand neben `schliesse`; F-447 Mutationen reichen). Wiederkehrende Klasse: keine.

**Finding-Klassen dieses Laufs:** keine

## Verdikt

**Merge-blockierend:** nein. Die Tests prüfen nicht weniger als vorher, die Brücke hält Punkt 7 ein, und `Open` verhält sich wie an `1c30e43^`.

**Übergabe:** F-445 geht an den Architect, der den Wortlaut von `SPEC-049` Punkt 7 führt. F-444 und F-447 gehen an den Verifier. F-446 erwartet keine Aktion. Dieser Report ersetzt keine Verifikation.
