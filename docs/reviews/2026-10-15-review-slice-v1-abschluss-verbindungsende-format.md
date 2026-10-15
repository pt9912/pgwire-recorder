# Review-Report: slice-v1-abschluss-verbindungsende-format — 2026-10-15

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-format`, Diff `6d0d697..e4ae4d4` (Schritt 35a: Modell, `Session.ValidateEnd`, YAML-Leser).

**Skill:** `.harness/skills/reviewer.md` @ e4ae4d4 ·
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

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

- Slice-Plan `slice-v1-abschluss-verbindungsende-format` (§1, §3, §6, §7 Belege des Implementers)
- `ADR-0039` (Aufzeichnung kennt das Verbindungsende); Randformen des Architect in `SPEC-001`, `SPEC-002` (Commit `690e353`)
- `LH-FA-25` (mit `LH-FA-25.a`), `LH-QA-06`
- `AGENTS.md` (Hard Rules) §3.7 bis §3.14

Eigene Läufe dieses Reviews (Bind-Mount in einer Kopie des Arbeitsbaums, `--network=none`, `go test ./internal/hexagon/... ./internal/adapters/driven/recording/...`): unveränderte Kopie grün; Mutanten einzeln, alle rot wie in §7 behauptet: `Session.ValidateEnd` ohne Prüfung „nicht letzte Interaktion“ (`TestSessionValidateEnd/…`, `TestUnmarshalEndeAbgelehnt/nicht_letzte_Interaktion,_einfach`) · `idx == 0` entfernt (`TestValidateEndeFehler/query_nur_connection_closed`, `TestUnmarshalEndeAbgelehnt/nur_Ende`; der Mutant endet in einer Panik, nicht in einer Fehlermeldung) · `flush` mit Ende abgelehnt (`TestValidateEndeGueltig`, `TestUnmarshalEnde/extended,_flush`) · Version-1-Prüfung des Lesers entfernt (`TestUnmarshalEndeAbgelehnt/Version_1,_einfach` und `…_extended`) · „unmittelbar“ zu „irgendein FATAL davor“ (`…Hinweis_zwischen_Fehler_und_Ende`, drei Fälle) · `Columns` aus `ohneFelder` (`TestValidateEndeFehler/query_mit_Spalten`) · Groups-Zweig von `Interaction.endet` (`TestSessionValidateEnd/Extended-…`, `TestUnmarshalEndeAbgelehnt/nicht_letzte_Interaktion,_extended`). Metrik von §4 nachgemessen: `git diff --numstat 6d0d697..HEAD -- internal docs/user` ergibt 590.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Zwei Randformen entscheidet der Code, ohne dass §6 oder die Spezifikation sie nennt; §7 vermerkt sie nur als „Lesarten“: (a) eine einfache Anfrage `[command_complete, ready_for_query, error_response FATAL, connection_closed]` lädt als gültig (`pruefeEnde` prüft kein `ready_for_query` vor dem Ende, obwohl das Ende laut DoD „statt `ready_for_query`“ steht); (b) `Interaction.endet` zählt `connection_closed` an jeder Stelle als Ende. Der Implementer hielt nicht an und gab nichts an den Architect zurück. | `AGENTS.md` §3.12 (Randformen vor dem Code; Rückgabe an den Architect) | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-format.md` · „Zwei Lesarten, die §6 nicht eigens nennt“ | nein — kein Gate sieht eine ungenannte Randform; ein Test mit der Antwortfolge (a) bestätigt, dass sie lädt | Randform im Code entschieden, nicht in §6 genannt |
| F-2 | LOW | §3 des Plans nennt `docs/user/abdeckung-e2e.md` als geändert, der Diff ändert sie nicht; geändert sind `abdeckung-gesamt.md` und `abdeckung-vollstaendig.md`, die §3 nicht nennt. Der Plan folgt dem Diff an dieser Stelle nicht. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-format.md` · „`docs/user/abdeckung-unit.md`, `docs/user/abdeckung-e2e.md`“ | nein — Abgleich von §3 mit `git diff --stat` | Plan-Tabelle §3 weicht vom Diff ab |
| F-3 | INFO | `LH-FA-25` steht in `abdeckung-gesamt.md` und `abdeckung-vollstaendig.md` als „vollständig“ (nur Unit), obwohl Aufzeichnen, Wiedergabe, Einspielen und SQLite in sieben weiteren Slices folgen; `make doc-trace` liest die Datei. Die Tabelle ist Werkzeug-Ausgabe über die Deklarationen, die Deklarationen der Tests selbst sagen nur zu, was sie prüfen. Der Implementer hat den Punkt für das Review gemeldet; Zuständigkeit liegt beim Planner (Bedeutung von „vollständig“ bei teilweise gelieferter Anforderung). | `AGENTS.md` §3.11 | `docs/user/abdeckung-vollstaendig.md` · „`LH-FA-25`“ | ja — `make abdeckung-check` bestätigt die Tabelle, nicht ihre Lesart | Abdeckungs-Stand „vollständig“ vor Lieferung der Anforderung |
| F-4 | INFO | Zum Flake `TestZeitangabenAbsteigendBeschaedigt/…sqlite`: Die Ursache ist nach Code-Lesung nicht Last, sondern Zufall im Pfad. Die Meldung des SQLite-Lesers enthält den Dateipfad (`"%s: Inhalt der SQLite-Datei beschädigt"`), der Test lehnt jede Meldung ab, die einen der Werte `4709`, `4710`, `4300`, `4900`, `4800` als Teilstring enthält, und `t.TempDir()` trägt eine Zufallszahl im Pfad. Nicht reproduziert; vorbestehend, die Datei ändert der Diff nur in einem Fall (`Version 3`). Der Mutant „S gilt, wenn V fehlt“ ist dem Rot nicht zuzurechnen (Einschätzung des Implementers insoweit bestätigt, die Begründung „Last“ nicht). | Maintainability | `internal/adapters/driven/recording/zeitangaben_pruefen_test.go` · „Meldung nennt den Wert %s“ | nein — Wiederholung des Tests mit Pfaden, die die Ziffernfolge enthalten | Test prüft Teilstring gegen Meldung mit Zufallspfad |
| F-5 | INFO | Grüne Mutanten in §7: (a) „zweites Ende“ ohne eigene Mutation und (b) `ready_for_query` in der Gruppe mit Ende ohne eigenen Code. Beide Begründungen tragen: (a) jede Session mit zwei Enden hat ein Ende außerhalb der letzten Interaktion, ein Mutant, der nur das zweite bricht, wäre äquivalent, und die Fälle `zwei_Enden_…` stehen im Test und werden von `ii < letzte-1` gefärbt; zwei Enden innerhalb einer Interaktion fängt `pruefeEnde` („nur als letzte Nachricht“, Fall vorhanden); (b) die vorhandene Klausel in `validateServer` trägt die Regel, die neuen Fälle färben ihre Mutation. Kein Befund, nur Beurteilung der Implementer-Aussage. | `AGENTS.md` §3.10 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-format.md` · „Zweites Ende in einer Session“ | ja — Mutanten nachgefahren (siehe Eingangs-Kontext) | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/model` (Produktcode: `response.go`, `recording.go`, `extended.go`) | geprüft, ohne weiteren Befund; Core-Reinheit gewahrt (Importe nur `errors`, `fmt`), Kommentare tragen Zusage/Kopplung und haben je einen Test, Reihenfolge in `Group.validate` entspricht dem Kommentar und §6 |
| `internal/hexagon/model` (Tests: `ende_test.go`, `extended_test.go`) | geprüft, ohne Befund; je Element der Aufzählungen (Schweregrade, Felder, Stellen, Ortsschritte der Extended-Form) ein Fall, sieben Mutanten selbst rot gesehen |
| `internal/adapters/driven/recording` (`yaml.go`, `yaml_ende_test.go`, drei geänderte Tests) | geprüft, ohne Befund; Version geht dem Inhalt voran, `PGR-E3002` für 3, 0, negativ, über 64 Bit, Version 1 und 2 angenommen, Schlüssel neben `type` (Wert, leer, `null`) abgelehnt, Abdeckungs-Deklarationen entsprechen den Tests |
| `test/integration/sqlite_lesen_e2e_test.go` | geprüft, ohne Befund; Fall „yaml, Version 3“, Fall „sqlite, Version 2“ bleibt |
| `docs/user/abdeckung-*.md` | geprüft, Befunde F-2 und F-3 |
| Größe und Schichten (`AGENTS.md` §3.13, §3.14) | geprüft, ohne Befund; Summe 590 (Einfügungen plus Löschungen), unter Meldeschwelle 800 und Marke 1000, zwei Liefer-Punkte, zwei Schichten (Kern, Recording-Adapter); keine Sendung an einen Nehmer im Diff |
| Suppression (`//nolint`) und Gate-Lockerung | geprüft, ohne Befund |
| Plan §1, §6, §7 gegen Diff | geprüft; Abgrenzung (kein Schreiber, kein SQLite) eingehalten, Befund F-1 und F-2 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Randform im Code entschieden, nicht in §6 genannt · Plan-Tabelle §3 weicht vom Diff ab · Abdeckungs-Stand „vollständig“ vor Lieferung der Anforderung · Test prüft Teilstring gegen Meldung mit Zufallspfad

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM): Der Architect entscheidet die beiden Randformen (a) und (b) in `SPEC-002` und §6, danach folgt Code und Test; F-2 bis F-5 blockieren nicht.

**Übergabe:** Findings gehen an den Implementer (F-1 über den Architect, F-3 an den Planner); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
