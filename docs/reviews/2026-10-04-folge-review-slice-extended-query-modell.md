# Review-Report: slice-extended-query-modell, Folge-Review — 2026-10-04

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Das ist ein Folge-Lauf zu [`2026-10-04-review-slice-extended-query-modell.md`](2026-10-04-review-slice-extended-query-modell.md) (F-283 bis F-293). Er schließt auch die Befunde aus [`2026-10-04-verifikation-slice-extended-query-modell.md`](2026-10-04-verifikation-slice-extended-query-modell.md) ein (V-16 bis V-21). Die DoD ist nicht Gegenstand dieses Reviews; sie prüft der Verifier.

**Gegenstand:** die Korrektur-Commits `6b4967d` (Review-Findings) und `d2efe37` (Verifikations-Befunde), `git show` je Commit, HEAD `d2efe37`. Die mitcommitteten Reports unter `docs/reviews/` sind nur als Eingang gelesen. Inhaltlicher Schwerpunkt ist die neue Vorprüfung am YAML-Baum `formVorpruefung`.

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-04

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-extended-query-modell.md` (Kopf, §1 bis §3, §6, §8), `docs/plan/planning/open/slice-extended-query-record.md` §6
- [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md), [ADR-0007](../plan/adr/0007-strict-replay.md), [ADR-0012](../plan/adr/0012-extended-query-gruppen.md)
- `spec/lastenheft.md` [LH-FA-07](../../spec/lastenheft.md), [LH-FA-18](../../spec/lastenheft.md)
- `spec/spezifikation.md` LH-FA-18.a, `SPEC-001`, `SPEC-002`, `SPEC-041`, Historienzeilen
- `AGENTS.md` (Hard Rules 3.3 bis 3.9, §5), Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (drei Belege, einer davon aus diesem Slice)
- Review F-283 bis F-293, Verifikation V-16 bis V-21, Mutationstabelle [`2026-10-04-mutationen-slice-extended-query-modell.md`](2026-10-04-mutationen-slice-extended-query-modell.md) (M1 bis M11, N1 bis N10, O1 bis O9, „Nicht abgedeckt“)

**Ausgeführte Läufe im Repo:** `make docs-check` und `make abdeckung-check` nach dem Anlegen dieser Datei (siehe unten). `make gates` lief nicht. Der Arbeitsbaum war vor dem Anlegen dieser Datei unverändert (`git status --short` leer).

**Mutationen und Sonden:** Nur in Kopien (`git archive`) im Scratchpad, je Kopie `docker build --target test` mit dem Tag `pgr-review-fr:*`. Sonden liefen als zusätzliche Testdatei mit `go test -v` in einem netzlosen Container der Stufe `source`. Ein Kontrolllauf ohne Mutation war grün. Images und Kopien sind gelöscht.

Mutationen an HEAD, einzeln:

| # | Mutation | Unit |
|---|---|---|
| Q1 | Gruppe: `server: []` wie fehlend abgelehnt | rot (`TestExtendedRoundtrip`, „abgeschnitten nach der Flush-Gruppe“) |
| Q2 | `type`: nur leerer Wert abgelehnt, `null` nicht | rot („type null an Extended“, „… an einfacher Anfrage“) |
| Q3 | `param_types` nur in Gruppen geprüft, nicht in `responses` | rot („param_types in einfacher Anfrage“, „… null …“) |
| Q4 | `param_types` an jeder Antwort abgelehnt, auch an `parameter_description` | rot (`TestExtendedRoundtrip`, `TestExtendedLeereFelder`) |
| Q5 | `client`/`server`: nur fehlend abgelehnt, `null` nicht | rot („Gruppe mit client: null“, „Flush-Gruppe mit server: null“) |
| Q6 | Vorprüfung auf `null` ohne `responses` | rot („extended mit responses: null“) |
| Q7 | Replay: Typprüfung des Matchers entfernt (F-283) | rot (`TestReplayExtendedAmCursor`) |
| Q8 | Client-Nachricht: Feld mit `null` als belegt gezählt (F-285) | rot („parse mit sql: null“) |

*Sonden* (`Unmarshal` direkt, an HEAD und an `6b4967d`):

| # | Eingabe | HEAD `d2efe37` | `6b4967d` |
|---|---|---|---|
| A1 | Flush-Gruppe mit Anker `&fg`, zweite Gruppe `- *fg`, dann Sync-Gruppe | **`PGR-E3003`: Gruppe 2 ohne client** | angenommen |
| A2 | wie A1, zweite Gruppe `- <<: *fg` (Merge-Key) | **`PGR-E3003`: Gruppe 2 ohne client** | angenommen |
| A2b | dieselben zwei Flush-Gruppen ausgeschrieben (Gegenstück) | angenommen | angenommen |
| A3 | `{<<: *pd, type: close_complete}`, `&pd` an einer gültigen `parameter_description` | **angenommen**, Modell: `close_complete` mit `ParamTypes [23]` | `PGR-E3003`: param_types gehört nicht zu close_complete |
| A10 | wie A3, in einer einfachen Anfrage an `command_complete` | **angenommen**, Modell: `command_complete` mit `ParamTypes [23]` | `PGR-E3003` |
| A4 | `request: *n`, `&n` an `offset_ms: null` einer anderen Interaktion, an einer Extended-Interaktion | angenommen | angenommen |
| A5 | `type: *e`, `&e` an einem leeren `sql` | angenommen, als einfache Anfrage | angenommen |
| A6 | `server: *n` (null) an einer Flush-Gruppe | angenommen | angenommen |
| A7 | `{type: *pdt, param_types: [23]}`, `&pdt` am Wert `parameter_description` | **`PGR-E3003`: param_types gehört nicht zu pdt** | angenommen |
| A8 | Interaktion per Merge-Key aus einer einfachen Anfrage | angenommen | angenommen |
| A9 | Session 2 mit Sync-Gruppe ohne `server` | `PGR-E3003`: **Gruppe 1 ohne server** (ohne Ort) | `… Session 2, Interaktion 1: …` |
| A11 | `Marshal` einer einfachen Anfrage mit allen acht Antworttypen, NULL- und Binärwert, `offset_ms`, dann `Unmarshal` | angenommen | — |
| — | `TestUnmarshalSpec041` (Beispiel aus `SPEC-041` wörtlich) | PASS | — |

Ein Alias-Knoten trägt in `yaml.v3` `Kind=AliasNode`, `Tag=""` und als `Value` den Namen des Ankers, keinen Inhalt; `wert` liefert für ihn `ohneWert`.

---

## Stand der Findings aus dem Vorlauf

| ID | Vorlauf | Stand | Beleg |
|---|---|---|---|
| F-283 | MEDIUM | **behoben** | `TestReplayExtendedAmCursor` mit leerer und nicht leerer Anfrage, Cursor bleibt stehen; Q7 rot. |
| F-284 | MEDIUM | **behoben** | `TestValidateGueltig`, Fall „query mit jedem Antworttyp“, enthält alle acht Typen von `simpleResponses` (Lesen). |
| F-285 | MEDIUM | **behoben** | `SPEC-041` nennt fehlendes Feld und `null` beschädigt; `clientMessageDTO.UnmarshalYAML` verlangt jeden Schlüssel mit Wert; Q8 rot. |
| F-286 | MEDIUM | **behoben** | `SPEC-001` nennt die Stelle von `type` je Art, `SPEC-041` stellt `offset_ms` nur bei Extended neben `type`. |
| F-287 | LOW | **behoben für ausgeschriebene Werte, mit `d2efe37` wieder offen über Merge-Keys** | Q3, Q4 rot; Sonden A3, A10, siehe F-295. |
| F-288 | LOW | **behoben** | `responseToDTO` schreibt an jeder `parameter_description` `param_types`, auch `[]`; `TestExtendedLeereFelder`. |
| F-289 | LOW | **behoben für ausgeschriebenes `null`** | Q6 rot. Über Aliase siehe F-295, Kommentare F-296. |
| F-290 | LOW | **nicht wiederholt** | Die Messages von `6b4967d` und `d2efe37` nennen keine `SPEC-*`- oder `ARC-*`-Kennung. |
| F-291 | LOW | **behoben** | §3 führt beide Folge-Slices, die Reports und die Register-Belege. |
| F-292 | INFO | **übergeben** | `slice-extended-query-record` §6. |
| F-293 | INFO | **behoben** | Die Mutationstabelle liegt vor; Q1 bis Q8 bestätigen ihre Aussagen für O1 bis O9 und N1, N6. |
| V-16 | MEDIUM | **behoben für ausgeschriebene Werte** | Q3, Q4 rot. Über Merge-Keys siehe F-295. |
| V-17 | LOW | **behoben** | `SPEC-041` macht `client`/`server` zur Pflicht, `server: []` zulässig; Q1, Q5 rot. Falsch-Ablehnung über Aliase siehe F-294. |
| V-18 | LOW | **behoben für ausgeschriebene Werte** | Q2 rot. Über Aliase siehe F-295 (A5). |
| V-19 | LOW | **teilweise** | Kommentare und Abdeckungszeile nennen die Formen; die Grenze der Vorprüfung nennen sie nicht, siehe F-296. |
| V-20 | LOW | **behoben** | Plan §6 trägt das Risiko einmal. |
| V-21 | INFO | **behoben** | Plan §6 Risiko 4 und `slice-extended-query-record` §6 nennen `Response.ParamTypes`. |

## Neue Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-294 | MEDIUM | `formVorpruefung` liest den YAML-Baum, ohne Aliase und Merge-Keys aufzulösen, und lehnt damit gültige Aufzeichnungen ab, die der Decoder liest und `Validate` besteht. Betroffen sind eine Gruppe als Alias oder per Merge-Key (A1, A2: „Gruppe 2 ohne client“) und ein Antworttyp als Alias (A7: „gehört nicht zu pdt“). Bei `6b4967d` wurden alle drei gelesen. `SPEC-041` und `SPEC-001` regeln Aliase nicht; eine Aufzeichnung wird laut [LH-FA-07](../../spec/lastenheft.md) mit Testcode versioniert und von Hand bearbeitet. Die Recorder-Ausgabe ist nicht betroffen (A11, `TestUnmarshalSpec041`). | [LH-FA-07](../../spec/lastenheft.md) Happy Path; `SPEC-041`; Reviewer-Skill MEDIUM „unklare Fehlerbehandlung am Rand des Spec-Bereichs“ | `internal/adapters/driven/recording/yaml.go:164` (`formVorpruefung`), `:182`, `:189`, `:202` (`wert`) | ja — Sonden A1, A2, A7 | Prüfung am YAML-Baum sieht den aufgelösten Wert nicht |
| F-295 | MEDIUM | Umgekehrt nimmt der Leser über Aliase und Merge-Keys an, was `SPEC-041` und `SPEC-001` beschädigt nennen. `d2efe37` hat die Ablehnung von `param_types` an anderen Antworten aus `responseFromDTO` (dekodierter Wert) in `formVorpruefung` (Baum) verlegt. Seitdem trägt das Modell `ParamTypes` an `close_complete` und `command_complete`, wenn der Schlüssel per Merge-Key kommt (A3, A10); bei `6b4967d` war das `PGR-E3003`. Dazu kommen `request: null` (A4) und `server: null` (A6) über einen Alias sowie ein leerer `type` über einen Alias (A5). Die Tabelle „Nicht abgedeckt“ nennt nur `null` über einen Alias, nicht Merge-Keys und nicht den leeren `type`. | `SPEC-041`; `SPEC-001`; [LH-FA-07](../../spec/lastenheft.md) Negative; Reviewer-Skill MEDIUM „unklare Fehlerbehandlung am Rand des Spec-Bereichs“ | `internal/adapters/driven/recording/yaml.go:176`, `:188` (`formVorpruefung`), `:618` (`responseFromDTO`) | ja — Sonden A3, A4, A5, A6, A10 | Prüfung am YAML-Baum sieht den aufgelösten Wert nicht |
| F-296 | LOW | Die Zusagen gehen weiter als die Prüfung. Der Kommentar von `formVorpruefung` sagt ohne Einschränkung, `request`/`responses`/`groups` stünden nicht mit `null` und `param_types` an keiner anderen Serverantwort. Der von `responseFromDTO` sagt, „dass es an keiner anderen Antwort steht, prüft formVorpruefung“. Der von `responseDTO.ParamTypes` sagt „an jeder anderen Antwort fehlt der Schlüssel“. Die Abdeckungszeile LH-FA-07/Negative sagt dasselbe. Alle vier gelten nicht über Aliase und Merge-Keys (F-295). Die Begründung im Kopf von `formVorpruefung` („was der Decoder nicht unterscheidet, weil er null wie einen fehlenden Schlüssel liest“) trifft auf `param_types` mit Wert nicht zu, den unterschied der Decoder. Die Message von `d2efe37` sagt zu: „Kommentare und Abdeckungs-Deklaration sagen genau das Geprüfte zu“. Die Klasse tritt in diesem Slice zum dritten Mal auf (F-289, V-19). | `AGENTS.md` §3.7; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Maintainability | `internal/adapters/driven/recording/yaml.go:153`, `:353`, `:614`; `docs/user/abdeckung-unit.md:25`; `internal/adapters/driven/recording/yaml_test.go` (Kommentar von `TestUnmarshalExtendedFehler`) | ja — Sonden A3, A4, A10 | Zusage weiter als Code |
| F-297 | LOW | Die bekannte Grenze der Vorprüfung (Aliase) steht nur in der Mutationstabelle unter `docs/reviews/`. Plan §6 Risiko 1 sagt zu den Randformen V-16 bis V-18: „Alle in der Spezifikation geschlossen.“ Kein Risiko und kein Folge-Slice nennt, dass der Leser Aliase nicht prüft. Die Mutationstabelle ist ein Beleg der Implementer-Rolle, kein Planungsdokument. | `AGENTS.md` §3.9; Maintainability | Slice-Plan §6 (Risiko 1); `docs/reviews/2026-10-04-mutationen-slice-extended-query-modell.md` §Nicht abgedeckt | ja (Lesen) | Bekannte Grenze nur im Beleg, nicht im Plan |
| F-298 | LOW | Meldungen aus `formVorpruefung` nennen Session und Interaktion nicht (A9: „Gruppe 1 ohne server“). Jeder Fehler aus `fromDTO` trägt „Session %d, Interaktion %d“. Bei `6b4967d` trug derselbe Fall (eine am Ende abgeschnittene Sync-Gruppe) den Ort. In einer Aufzeichnung mit vielen Sessions ist die Stelle nicht mehr zu finden. Das Review zu `cd03e09` hatte unter den Negativbefunden „Jeder Fehler wird `PGR-E3003` mit Session und Interaktion“ festgehalten. | [LH-FA-07](../../spec/lastenheft.md) Negative; Maintainability | `internal/adapters/driven/recording/yaml.go:173`, `:177`, `:183`, `:190` | ja — Sonde A9 | Fehlermeldung ohne Ort |
| F-299 | LOW | Der Umbruch in der Spezifikation ist an zwei geänderten Stellen zerfallen: `SPEC-001` hat eine überlange Zeile, gefolgt von einer halben (`spezifikation.md:1031`–`1032`); `SPEC-041` hat eine überlange Zeile, danach steht „Simple-Interaktionen“ allein in einer Zeile (`:1143`–`1145`). Inhaltlich ohne Folge. | Maintainability | `spec/spezifikation.md:1031`, `:1143` | ja (Lesen) | Zeilenumbruch nach Teilersetzung zerfallen |
| F-300 | INFO | `SPEC-001` sagt, abwärtskompatible Ergänzungen (neue optionale Felder) ändern `version` nicht. `SPEC-041` nennt „ein anderes Feld“ einer Client-Nachricht beschädigt, und der Leser lehnt jeden unbekannten Schlüssel ab (`KnownFields`, Schlüssel-Liste je Typ). Ein Leser der Version 1 meldet eine Datei mit einem künftigen optionalen Feld also als beschädigt, obwohl die Version gleich bleibt. Die Spannung liegt nicht im Lastenheft und bestand vor diesen Commits; mit den Präzisierungen ist sie deutlicher. Hinweis an den Planner. | `SPEC-001`; `SPEC-041` | `spec/spezifikation.md:1026`, `:1130` | ja (Lesen) | Spec-Sätze ohne gemeinsame Lesart |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driven/recording` — `formVorpruefung` gegen `SPEC-001`/`SPEC-041` | geprüft. Die vier Regelgruppen des Kommentars entsprechen den Sätzen der Spezifikation; für ausgeschriebene Werte hat jede Regel einen Negativfall mit Meldung und ein gültiges Gegenstück (`server: []`, `parameter_description` mit und ohne Parameter), Q1 bis Q6 rot. Recorder-Ausgabe einfacher Anfragen und Beispiel aus `SPEC-041` werden gelesen (A11, `TestUnmarshalSpec041`). Ein `type: query` auf Interaktionsebene besteht die Vorprüfung und wird danach in `interactionFromDTO` abgelehnt; doppelte Schlüssel lehnt der Decoder vorher ab. Befund F-294, F-295, F-296, F-298. |
| `internal/adapters/driven/recording` — Schreiber | geprüft, ohne Befund. `responseToDTO` schreibt `param_types` nur und stets an `parameter_description`; `groupDTO.Server` ohne `omitempty`, eine Flush-Gruppe ohne Server-Nachricht wird als `server: []` geschrieben und gelesen (`TestExtendedRoundtrip`). |
| `internal/adapters/driven/recording` — Tests | geprüft. Neue Fälle tragen die Regel im Namen und prüfen den Meldungstext; die Umbenennung „abgeschnitten in den Client-Nachrichten“ → „Gruppe ohne flush oder sync“ beschreibt jetzt, was der Fall prüft; abgeschnittene Gruppen fallen über die Vorprüfung auf. Kein Fall mit Alias oder Merge-Key (F-294, F-295). |
| `internal/hexagon/model` | geprüft, ohne Befund. Nur Test und Testkommentar geändert; Core-Reinheit unverändert ([ADR-0006](../plan/adr/0006-kanonisches-domain-model.md)). |
| `internal/hexagon/services` | geprüft, ohne Befund. Nur `replay_test.go`; der Test hält die Abweichung nach [ADR-0007](../plan/adr/0007-strict-replay.md) und das Stehenbleiben des Cursors (Q7 rot). |
| `spec/` — Strata und Hard Rule 3.4 | geprüft. Historienzeilen nennen nur Spec-Kennungen, keine ADR, keinen Slice, keinen Commit. Die Präzisierungen machen nur Formen beschädigt, die keine Aufzeichnung des Recorders trägt; kein Satz von [LH-FA-07](../../spec/lastenheft.md) oder [LH-FA-18](../../spec/lastenheft.md) wird eingeschränkt. Das Lastenheft ist unverändert. Befund F-299, F-300. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. Kein Move in `6b4967d` und `d2efe37`. |
| Hard Rule 3.5, 3.8 — ADRs | geprüft, ohne Befund. Keine ADR im Diff. |
| Hard Rule 3.6 — Gates nicht lockern | geprüft, ohne Befund. Keine Gate-Konfiguration im Diff; die Abdeckungstabelle ist erzeugt. |
| Hard Rule 3.7 — Kommentare | geprüft. Indikativ, Zusage oder Kopplung; keine verworfene Alternative, kein abwesender Text, kein abgebrochener Satz. Befund F-296. |
| Hard Rule 3.9 — Plan folgt Korrektur | geprüft. §3, §6 und §8 folgen beiden Commits; §1 und der Kopf bleiben gültig; `slice-extended-query-record` §6 übernimmt F-292 und V-21. Befund F-297. |
| Abdeckungs-Deklaration ([ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)) | geprüft. Die Zeile LH-FA-07/Negative trifft das Kriterium „beschädigte Datei“; Tabelle und Testkommentar sind gleich. Befund F-296. |
| Commit-Messages `6b4967d`, `d2efe37` | geprüft, ohne Befund. Beide tragen `slice-extended-query-modell`, [LH-FA-18](../../spec/lastenheft.md) und [LH-FA-07](../../spec/lastenheft.md), keine Struktur-ID. |
| Register `BEO-REPO/*` | geprüft. Die Belege zu F-283/F-284 und F-289/V-19 liegen. Befund F-296 ist ein weiterer Fall derselben Klasse. |

## Summary

| Kategorie | Anzahl (neu) |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 1 |

Vorlauf: F-283 bis F-286, F-288, F-291, F-293, V-17, V-20, V-21 behoben; F-287, F-289, V-16, V-18 nur für ausgeschriebene Werte behoben (Rest F-295); V-19 teilweise; F-290 nicht wiederholt; F-292 übergeben.

**Finding-Klassen dieses Laufs:**

- Prüfung am YAML-Baum sieht den aufgelösten Wert nicht (F-294, F-295)
- Zusage weiter als Code (F-296; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, im Slice zum dritten Mal nach F-289 und V-19)
- Bekannte Grenze nur im Beleg, nicht im Plan
- Fehlermeldung ohne Ort
- Zeilenumbruch nach Teilersetzung zerfallen
- Spec-Sätze ohne gemeinsame Lesart

## Verdikt

**Merge-blockierend:** nein, es gibt kein HIGH. Die Korrekturen tragen die Findings des Reviews und die Befunde der Verifikation für jede ausgeschriebene Form; jede neue Zusage hat einen Test, den eine Mutation an ihr rot färbt (Q1 bis Q8). Die Spezifikations-Präzisierungen widersprechen keinem Satz des Lastenhefts.

Vor der Closure verdienen die zwei MEDIUM Aufmerksamkeit. Sie haben eine gemeinsame Ursache: `d2efe37` prüft am unaufgelösten YAML-Baum und hat eine Prüfung vom dekodierten Wert dorthin verlegt. Deshalb lehnt der Leser seitdem gültige Dateien mit Aliasen ab (F-294) und nimmt per Merge-Key eine Form an, die bei `6b4967d` noch abgelehnt wurde (F-295). Ob Aliase und Merge-Keys im Format zulässig sind, regelt die Spezifikation nicht.

**Übergabe:**

- F-294, F-295, F-296, F-298 an den Implementer. Die Frage, ob `SPEC-001`/`SPEC-041` Aliase und Merge-Keys zulassen, liegt im Spec-Stratum 2, das dieser Slice berührt; ohne diese Entscheidung hat der Leser keinen Maßstab.
- F-297, F-299 an den Implementer (Plan §6, Spezifikation).
- F-300 an den Planner.
- Steering-Loop: F-296 ist innerhalb eines Slice das dritte Auftreten von `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, jeweils im Korrektur-Commit nach dem vorigen Fund. Die Frage aus dem Review an den Architect (Regel „Prüfung je Zusage“) ist weiter offen.
- Widerspricht der Implementer einer MEDIUM-Einstufung, genügen Annahme oder Begründung. Der Konflikt-Pfad über den Architect gilt ab HIGH oder ab dem dritten gleichen Konflikttyp.
- Die Finding-Klassen gehen in die Slice-Closure §7. Die DoD-Konformität prüft der Verifier separat.
