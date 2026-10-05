# Viertes Folge-Review: slice-extended-query-replay — 2026-10-05

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** Die Commits `f528c90` und `1ab9db2` (`git diff 7ad8f5b..1ab9db2` ohne die Report-Dateien unter `docs/reviews/`). Sie sollen F-347 und F-348 aus dem [dritten Folge-Review](2026-10-05-folge-review-3-slice-extended-query-replay.md) und V-25 bis V-28 aus der [Verifikation](2026-10-05-verifikation-slice-extended-query-replay.md) abarbeiten. Frühere Teile des Slice prüft dieses Review nicht erneut. Schwerpunkte laut Auftrag:

- Status je Befund.
- Der Bestätigungsfilter in `objekte` (`ParseComplete`, `BindComplete`, `CloseComplete`; Zählung je Interaktion; „die ersten n Nachrichten einer Art“), geprüft gegen das Protokoll: Fehler mitten in der Gruppe, Close auf ein nicht vorhandenes Objekt, späte Bestätigungen.
- `nachrichtText` und das Abbrechen von `objekte` an einer einfachen Interaktion am Cursor.
- Parameterwerte erscheinen nie.

Maßstab: nur echte Fehler oder nicht eingelöste Zusagen. Gegenstand ist der Text der Diagnose; die Erkennung der Abweichung ist bestätigt.

**Ablage:** Eigene Datei, Nummerierung fortlaufend ab F-349 (Reviewer-Skill §Output-Schema: „Folgeläufe als neue Datei statt Überschreibung“).

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-extended-query-replay.md` (Kopf, §1, §3, §6)
- `spec/lastenheft.md` [LH-FA-10](../../spec/lastenheft.md), [LH-FA-18](../../spec/lastenheft.md); `spec/spezifikation.md` LH-FA-10.a, LH-FA-18.a (§Gruppen, §Fehler, §Mismatch)
- `internal/hexagon/model` (`ClientMessageType` und `ResponseType` als Zeichenketten, `validateExtended`, `Group.validate`)
- `AGENTS.md` Hard Rules 3.3 bis 3.9; Register `BEO-REPO/plan-folgt-korrektur-nicht`
- PostgreSQL-Protokolldokumentation (Extended Query): Parse, Bind und Close werden bei Erfolg mit `ParseComplete`, `BindComplete` und `CloseComplete` bestätigt; Close auf ein nicht vorhandenes Objekt ist kein Fehler und wird ebenfalls bestätigt; nach einem Fehler verwirft der Server alle Nachrichten bis zum `Sync`.
- Vorherige Reports am Modul: [Review](2026-10-05-review-slice-extended-query-replay.md) (F-321 bis F-329), [Folge-Review](2026-10-05-folge-review-slice-extended-query-replay.md) (F-330 bis F-338), [zweites Folge-Review](2026-10-05-folge-review-2-slice-extended-query-replay.md) (F-339 bis F-346), [drittes Folge-Review](2026-10-05-folge-review-3-slice-extended-query-replay.md) (F-347, F-348), [Verifikation](2026-10-05-verifikation-slice-extended-query-replay.md) (V-25 bis V-30)

**Ausgeführte Läufe im Repo:** `make abdeckung-check` (Exit 0) und `make docs-check` nach dem Anlegen dieser Datei. `make gates` lief nicht. Vor dem Anlegen dieser Datei war der Arbeitsbaum unverändert.

**Gegenproben und Mutationen:** Nur in einer Kopie (`git archive 1ab9db2`) im Scratchpad. Die Unit-Läufe liefen in einem Container aus der Stufe `deps` des `Dockerfile` (Hilfs-Tag ohne Netz). Kopie und Hilfs-Image sind gelöscht.

| # | Lauf oder Mutation | Ergebnis |
|---|---|---|
| K1 | Kopie ohne Änderung: `gofmt -l internal/`, `go vet ./internal/hexagon/...`, Unit-Tests von `internal/hexagon/services` | grün, `gofmt` leer |
| M1 | Abbruch an der einfachen Interaktion am Cursor entfernt (`replay.go:235` nur `pi > c.pos`) | **grün** (F-349) |
| M2 | Bestätigungen je Gruppe statt je Interaktion (Stand `f528c90`) | rot in `TestReplayExtendedDiagnoseSpaeteBestaetigung` |
| M3 | `Query` ohne `nachrichtText` für die erwartete Nachricht | rot in `TestReplayExtendedDiagnoseArt` (P4a) |
| M4 | Ende der Aufzeichnung ohne `nachrichtText` | rot in `TestReplayExtendedDiagnoseArt` (P4c, P4c bind) |
| M5 | Art-Abweichung in `ClientMessage` ohne `nachrichtText` | rot in `TestReplayExtendedDiagnoseArt` (P4b, P4b bind) |
| M6 | Bestätigungsfilter entfernt (jede Nachricht gilt als angenommen) | rot in `TestReplayExtendedDiagnoseLebensdauer` (S1, S2, verworfenes bind) |

Eigene Sonden in der Kopie (zusätzliche Testdatei, nicht im Repo):

| # | Lage | Diagnose |
|---|---|---|
| Q1 | Interaktion 1: `parse ""`, `bind ""`, `Sync` (`T`); am Cursor die einfache Anfrage `SELECT 7`, empfangen `bind ""` bzw. `execute` | `empfangen Client-Nachricht bind (Anweisung "SELECT 1")`, ebenso `execute`; richtig. Unter M1 `Anweisung unbekannt` |
| Q2 | Fehler mitten in der Gruppe: `parse s1` bestätigt, `parse s2` abgelehnt, `bind s1` verworfen; danach Abweichung `bind s2` statt `bind s1` | `Anweisung erwartet "A", empfangen unbekannt`; richtig |
| Q3 | Fehler in einer Flush-Gruppe auf `parse s1`; die Sync-Gruppe danach (`parse s2`, `bind s2`) verworfen; danach `bind s2` | `empfangen unbekannt`; richtig |
| Q4 | `close` auf ein nicht vorhandenes Statement, danach `close s1`, beide mit `CloseComplete`; danach `bind s1` | `empfangen unbekannt`; richtig, die Zählung verschiebt sich nicht |
| Q5 | späte Bestätigung von `parse s1` in der folgenden Flush-Gruppe, dort dann Fehler auf `parse s2`; am Cursor `bind s2` statt `bind s1` | `Anweisung erwartet "A", empfangen unbekannt`; richtig |
| Q6 | wie Q1, empfangen `bind ""` mit dem Parameterwert `GEHEIM` | Wert erscheint nicht in der Meldung |

---

## Status der Befunde

| ID | Stand | Beleg |
|---|---|---|
| F-347 | behoben | `objekte` wendet ein `parse`, `bind` oder `close` nur an, wenn eine Bestätigung des Servers in der Interaktion dafür steht (`replay.go:249`–`:263`). Neue Fälle S1, S2 und „verworfenes bind“ in `TestReplayExtendedDiagnoseLebensdauer`, unter M6 rot. Plan §1 (`:34`) nennt die Regel. |
| F-348 | behoben (Hinweis aufgenommen) | Der Kommentar an `objekte` (`replay.go:223`–`:231`) und Plan §6 (`:120`) nennen die Grenze; §1 sagt, dass SQL-Befehle, die Objekte beenden, nicht nachgebildet werden. |
| V-25 | behoben, mit Rest bei der Prüfung | `nachrichtText` (`replay.go:185`) an allen drei Stellen (`:106`, `:132`, `:137`). M3 bis M5 sind rot. Rest: Der Abbruch an der einfachen Interaktion am Cursor, auf den die Herleitung für „erwartet Anfrage“ angewiesen ist, hält kein Test (F-349). |
| V-26 | behoben | Gezählt wird je Interaktion (`replay.go:249`–`:254`); `TestReplayExtendedDiagnoseSpaeteBestaetigung` ist unter M2 rot. |
| V-27 | behoben | Plan §3 (`:71`) nennt `replay_extended_test.go`, eine neue Zeile nennt die Berichte unter `docs/reviews/`. |
| V-28 | behoben | Der Kopf nennt [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) im `Bezug` (`:14`) und `LH-FA-11.a` in den berührten Spec-Stellen (`:16`). |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-349 | LOW | Die neue Bedingung `pi == c.pos && in.Request.Type == model.RequestQuery` beendet die Nachbildung vor einer einfachen Anfrage am Cursor, die noch nicht gelaufen ist. Der Code ist richtig (Q1), aber kein Test bricht, wenn die Bedingung fehlt (M1 grün). `TestReplayExtendedDiagnoseArt` prüft „erwartet Anfrage“ mit `bind` nur ohne vorangehende Interaktion („unbekannt“, P4b bind). Fehlt die Bedingung, löscht die Nachbildung das unbenannte Statement und Portal und wertet das ReadyForQuery der nicht gelaufenen Anfrage als Transaktionsende; die Diagnose nennt dann „unbekannt“ für ein bestehendes Objekt. Die Abdeckungszeile (`abdeckung-unit.md:35`) sagt für diesen Test „bei bind das hergeleitete oder ‚unbekannt‘“; der Fall „erwartet Anfrage“ ist nur in der zweiten Hälfte belegt. | LH-FA-10.a; Plan §1 („soweit es nach dem Verlauf der Aufzeichnung vor dem Cursor besteht“) | `internal/hexagon/services/replay.go:235`; `internal/hexagon/services/replay_extended_test.go:481`–`:497`; `docs/user/abdeckung-unit.md:35` | ja (M1, Sonde Q1) | Bedingung der Diagnose-Herleitung ohne Test, der sie bricht |

## Antwort auf die Prüffragen

1. **Befunde.** F-347, F-348, V-26, V-27 und V-28 sind behoben. V-25 ist im Code behoben, mit einem Rest in der Prüfung (F-349).
2. **Bestätigungsfilter gegen das Protokoll.**
   - **Präfix-Eigenschaft:** Nach einem Fehler verwirft der Server alles bis zum `Sync`, und eine Extended-Interaktion endet genau mit ihrem `Sync` (`Group.validate`). Die angenommenen Nachrichten einer Interaktion sind also ein Präfix ihrer Client-Nachrichten, und damit auch je Art. „Die ersten n einer Art mit n Bestätigungen“ trifft genau diese. Das gilt auch für die Interaktion am Cursor: Die Bestätigungen späterer Gruppen der Aufzeichnung fallen zuerst auf die Nachrichten vor dem Cursor, und eine dort abgelehnte Nachricht ließe keine spätere Bestätigung zu.
   - **Fehler mitten in der Gruppe:** richtig (Q2, Q3, Q5; S1, S2 im Repo).
   - **Close auf ein nicht vorhandenes Objekt:** Der Server bestätigt es; die Zählung bleibt richtig, und das Entfernen eines nicht vorhandenen Eintrags ändert nichts (Q4).
   - **Späte Bestätigungen:** Bei einer Sync-Gruppe gehören alle Server-Nachrichten bis zum `ReadyForQuery` zu ihr (LH-FA-18.a §Gruppen); über die Grenze einer Interaktion wandert keine Bestätigung. Innerhalb einer Interaktion zählt die Zählung je Interaktion sie mit (Q5, `TestReplayExtendedDiagnoseSpaeteBestaetigung`).
   - **Nicht zugeordnete Antworttypen:** `bestaetigung[r.Type]` liefert für jede andere Antwort `""`; `ClientMessageType` ist eine Zeichenkette ohne leeren Wert, ein Zusammenfallen mit einer Client-Art ist ausgeschlossen.
3. **`nachrichtText` und der Abbruch an einfachen Interaktionen.** `nachrichtText` nennt Typ und, wo es eine gibt, die Anweisung über `anweisung`, also dieselbe Herleitung wie `anweisungen`. Am Ende der Aufzeichnung läuft die Schleife über alle Interaktionen; der frühere Ausdruck `Interactions[:c.pos+1]` hätte dort den Bereich überschritten, der neue kann es nicht. Eine einfache Anfrage am Cursor ist noch nicht gelaufen und wird richtig nicht nachgespielt (Q1); ungeprüft ist das in den Tests (F-349).
4. **Parameterwerte.** `nachrichtText` und `anweisung` lesen nur Typ, Namen und SQL, nie `Params` (Q6). Bei `parse` erscheint das empfangene SQL wie auf dem einfachen Pfad.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services` | geprüft. Befund F-349. `objekte` läuft weiter nur im Mismatch-Zweig unter `mu` und ändert den Cursor nicht. |
| Core-Reinheit | geprüft, ohne Befund. Keine neuen Imports; `make a-check` lief nicht. |
| `docs/user/` | geprüft. Die zwei neuen Abdeckungszeilen entsprechen den Deklarationen (`make abdeckung-check` Exit 0); Reichweite der Zeile 35 siehe F-349. |
| `docs/plan/planning/` — Hard Rule 3.9 | geprüft, ohne Befund. §1 nennt Bestätigungsfilter, Zählung je Interaktion, Diagnose der Art-Abweichung und die Grenze bei SQL-Befehlen; §3 und Kopf sind nachgezogen; §6 führt die Grenze aus F-348 mit Ausgang. |
| `spec/` — Strata, Hard Rule 3.4 | geprüft, ohne Befund. Keine Spec-Datei im Diff. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. Kein Move in beiden Commits. |
| Hard Rules 3.5, 3.6, 3.8 | geprüft, ohne Befund. Keine ADR und keine Gate-Konfiguration im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft, ohne Befund. Die neuen Kommentare an `nachrichtText`, `objekte`, `bestaetigung`, `nachspielen` und den Test-Hilfen stehen im Indikativ; der Satz zur Grenze bei SQL-Befehlen ist eine Grenze. |
| Commit-Messages `f528c90`, `1ab9db2` | geprüft, ohne Befund. Beide nennen `slice-extended-query-replay` und [LH-FA-10](../../spec/lastenheft.md), [LH-FA-18](../../spec/lastenheft.md), keine Struktur-ID. |

## Summary

1 Finding (0 HIGH, 0 MEDIUM, 1 LOW, 0 INFO); F-347, F-348, V-25 bis V-28 behoben, V-25 mit Prüf-Rest; wiederkehrende Klasse: Diagnose-Herleitung ohne Test, der ihre Bedingung bricht.
