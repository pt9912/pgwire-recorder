# Review-Report: slice-extended-query-modell — 2026-10-04

**Review-Art:** Code, geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10 §Drei Review-Arten). Die DoD ist nicht Gegenstand dieses Reviews; sie prüft der Verifier.

**Gegenstand:** Commit `cd03e09` (`git show cd03e09`). Er enthält: Domain-Typen für Extended Query (`ClientMessage`, `Group`, Server-Nachrichten, `RequestExtended`), `Interaction.Validate`, YAML-Schreiber und YAML-Leser für `type: extended`, Unit-Tests, eine Abdeckungszeile, Abgleich von `LH-FA-18.a`, `SPEC-041` und `ARC-001`, den Slice-Plan sowie die Folge-Slices `slice-extended-query-record` und `slice-extended-query-replay`. Die reinen Moves davor (`10cbdef`, `bb7133b`) und die Verweis-Commits (`4390713`, `d2421e8`) sind nur auf Hard Rule 3.3 geprüft.

**Skill:** `.harness/skills/reviewer.md` @ `77a8c93`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-04

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-extended-query-modell.md` (§1 bis §8), `docs/plan/planning/open/slice-extended-query-record.md` und `slice-extended-query-replay.md` (§1 bis §4, §6)
- [ADR-0006](../plan/adr/0006-kanonisches-domain-model.md), [ADR-0007](../plan/adr/0007-strict-replay.md), [ADR-0012](../plan/adr/0012-extended-query-gruppen.md)
- `spec/lastenheft.md` [LH-FA-07](../../spec/lastenheft.md), [LH-FA-18](../../spec/lastenheft.md)
- `spec/spezifikation.md` LH-FA-18.a, `SPEC-001`, `SPEC-002`, `SPEC-041`, `SPEC-043` (Spalte `kind`); `spec/architecture.md` `ARC-001`, `ARC-008`
- `AGENTS.md` (Hard Rules 3.3 bis 3.9, §5), Beobachtungs-Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` und `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (je `Stand: offen`, zwei Belege)
- die Vorlauf-Reports am selben Modul: `2026-10-04-review-slice-walking-skeleton-replay.md` (F-273 bis F-282, insbesondere F-281) und `…-record-folge.md`

**Ausgeführte Läufe im Repo:** `make a-check` (0 Befunde; Hinweis auf `test/integration/*` wie F-248), `make docs-check` (0 Befunde), `make abdeckung-check` (grün). `make test` lief nicht im Arbeitsbaum, sondern als `docker build --target test` in einer Kopie (siehe unten). `make gates` lief nicht. Danach war der Arbeitsbaum unverändert (`git status` leer).

**Mutationen und Sonden:** Nur in Kopien (`git archive HEAD`) im Scratchpad, je Kopie `docker build --target test` mit eigenem Image-Tag. Ein Kontrolllauf ohne Mutation war grün. Alle Images und Kopien des Reviews sind gelöscht.

| # | Mutation | Unit |
|---|---|---|
| M1 | Replay: Prüfung `erwartet.Request.Type != model.RequestQuery` entfernt | **grün** |
| M2 | Leser: einfache Anfrage mit `groups` angenommen | rot |
| M3 | `Validate`: Zielart nur noch für `describe`, nicht für `close` geprüft | rot |
| M4 | Leser: Schlüssel-Whitelist je Client-Nachricht abgeschaltet | rot |
| M5 | `Validate`: Antworttyp einer einfachen Anfrage nicht geprüft | rot |
| M6 | `simpleResponses` ohne `parameter_status` | **grün** |
| M7 | `ready_for_query` am Ende jeder Gruppe statt nur der letzten zugelassen | rot |
| M8 | Leser: Extended-Interaktion mit `responses` angenommen | rot |
| M9 | `Validate`: Extended-Interaktion mit SQL angenommen | rot |
| M10 | Leser ruft `Validate` nicht | rot |
| M11 | Regel „sync genau am Ende der letzten Gruppe“ abgeschaltet | rot |
| M12 | Schreiber übernimmt `ParamTypes` einer Serverantwort nicht | rot |
| M13 | Gruppe ohne Client-Nachricht angenommen | rot |
| M14 | Leser: Prüfung „Client-Nachricht ist eine Abbildung“ entfernt | grün (äquivalent: ein Skalar führt zu „Client-Nachricht "" unbekannt“) |
| M15 | `simpleResponses` ohne `empty_query_response` | **grün** |
| M16 | `simpleResponses` ohne `notice_response` | rot |
| M17 | `simpleResponses` ohne `error_response` | **grün** |
| M18 | `simpleResponses` ohne `row_description` | rot |

*Sonden* (Testdatei in einer Kopie, `Unmarshal`/`Marshal` direkt):

- A: Eine Gruppe mit `- type: parse`, `- type: describe` mit `target: portal`, aber ohne `name`, `- type: execute` und `- type: sync`, also ohne die übrigen Felder ihres Typs: wird ohne Fehler gelesen.
- B: `param_types: [23]` an einem `command_complete` einer einfachen Anfrage: wird ohne Fehler gelesen. Vor `cd03e09` lehnte `KnownFields` den Schlüssel ab.
- C: Eine `parameter_description` ohne Parameter wird als `- type: parameter_description` geschrieben, ohne `param_types`.
- D: `param_types` an `parse_complete` in einer Extended-Gruppe: wird ohne Fehler gelesen.
- E: Eine Extended-Interaktion mit zusätzlichem `request: null`: wird ohne Fehler gelesen.
- F: `Marshal` einer Extended-Gruppe mit zwei `sync`: schreibt ohne Fehler.
- Rückwärtskompatibilität: `Marshal` einer Aufzeichnung mit einfachen Anfragen aller acht Antworttypen und `offset_ms` liefert vor (`cd03e09~1`) und nach `cd03e09` byteidentische Ausgabe (58 Zeilen). Diese Ausgabe liest der neue Leser.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-283 | MEDIUM | Plan §1 sagt zu: „eine Extended-Interaktion am Cursor ist für eine einfache Anfrage eine Abweichung (`PGR-E5001`)“. Seit `cd03e09` liest der Leser Extended-Interaktionen, damit ist dieser Pfad über eine Datei erreichbar. Kein Test hält die Prüfung des Anfrage-Typs im Matcher; sie kam mit der Antwort auf F-281. Ohne sie bleibt die Unit-Suite grün (M1). Erhält der Replay eine leere einfache Anfrage (`""`), während der Cursor auf einer Extended-Interaktion steht, passt dann der SQL-Text. Der Service liefert die leeren `Responses` ohne `ReadyForQuery` und rückt den Cursor vor. Das ist der kritische Pfad; heute hält er nur, weil die ungetestete Bedingung im Code steht. | [ADR-0007](../plan/adr/0007-strict-replay.md); [LH-FA-10](../../spec/lastenheft.md); Slice-Plan §1; Reviewer-Skill MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag“ | `internal/hexagon/services/replay.go:99`; `internal/hexagon/services/replay_test.go` (kein Fall mit `RequestExtended`) | ja — Mutation M1 | Protokollrand ohne Negativtest |
| F-284 | MEDIUM | Der Kommentar von `Validate` sagt zu: Eine einfache Anfrage trägt „nur Antworttypen des Simple Query Protocol“. Die Menge dafür, `simpleResponses`, ist mit dem Commit aus dem Leser ins Modell gewandert. Drei ihrer Mitglieder hält kein Test: Ohne `error_response`, `empty_query_response` oder `parameter_status` bleibt die Unit-Suite grün (M17, M15, M6). Dann wäre jede Aufzeichnung mit einer fehlgeschlagenen Anfrage, einer leeren Anfrage oder einem `SET` beim Laden beschädigt (`PGR-E3003`). `TestValidateGueltig` prüft für Extended jeden Typ, für die einfache Anfrage nur `command_complete` und `ready_for_query`. Die Lücke bestand im Leser schon vorher; mit dem Umzug ist sie jetzt Teil eines neuen Modell-Vertrags. | [LH-FA-07](../../spec/lastenheft.md) Happy Path; [LH-FA-18](../../spec/lastenheft.md) (LH-FA-18.a: `ParameterStatus`); Reviewer-Skill MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag“ | `internal/hexagon/model/response.go:34` (`simpleResponses`); `internal/hexagon/model/extended_test.go:55` (`TestValidateGueltig`) | ja — Mutationen M6, M15, M17 | Protokollrand ohne Negativtest |
| F-285 | MEDIUM | `SPEC-041` sagt: „Jede Client-Nachricht trägt neben `type` genau die Felder ihres Typs, auch wenn sie leer sind“. Als beschädigt nennt es aber nur „ein anderes Feld oder ein anderer Typ“. Der Leser liest eine Client-Nachricht, der Felder ihres Typs fehlen (Sonde A: `parse` ohne `statement`, `sql` und `param_types`, `execute` ohne `portal` und `max_rows`, `describe` ohne `name`). Die Werte werden dann als leer übernommen. Ob ein fehlendes Feld beschädigt ist oder leer bedeutet, entscheidet weder die Spezifikation noch der Plan. Der Feldvergleich in `slice-extended-query-replay` vergleicht dann gegen Werte, die nicht in der Datei stehen. | `SPEC-041`; [LH-FA-07](../../spec/lastenheft.md) Negative; Reviewer-Skill MEDIUM „unklare Fehlerbehandlung am Rand des Spec-Bereichs“ | `internal/adapters/driven/recording/yaml.go:240` (`clientMessageDTO.UnmarshalYAML`); `spec/spezifikation.md` §SPEC-041 | ja — Sonde A | Spec-Wortlaut deckt Leser-Entscheidung nicht |
| F-286 | MEDIUM | `SPEC-001` zählt `type: query` und `type: extended` zu Version 1. `SPEC-041` sagt für „jede Interaktion, einfach oder Extended“, `offset_ms` stehe „neben `sequence` und `type`“. Der Leser lehnt `type: query` auf Ebene der Interaktion mit „Interaktions-Typ "query" unbekannt“ ab (`PGR-E3003`); ein Test legt das fest. Das passt zum Beispiel in `SPEC-002`, wo `type: query` unter `request` steht. Wer aber `SPEC-001` und `SPEC-041` wörtlich liest, schreibt eine Datei, die abgelehnt wird, mit einer Meldung, die `SPEC-001` widerspricht. Der Commit hat `SPEC-041` geändert, diese Stelle nicht. | `SPEC-001`; `SPEC-041`; Reviewer-Skill MEDIUM „unklare Fehlerbehandlung am Rand des Spec-Bereichs“ | `internal/adapters/driven/recording/yaml.go:491`; `internal/adapters/driven/recording/yaml_test.go` (`TestUnmarshalExtendedFehler`, Fall „type: query auf Interaktionsebene“); `spec/spezifikation.md` §SPEC-001, §SPEC-041 | ja (Lesen) | Spec-Wortlaut deckt Leser-Entscheidung nicht |
| F-287 | LOW | `responseDTO` hat jetzt `param_types`, und der Leser nimmt den Schlüssel an jeder Serverantwort an: in einer einfachen Anfrage (Sonde B) und an jeder Extended-Server-Nachricht (Sonde D). Das Modell übernimmt den Wert, der Schreiber gibt ihn wieder aus. Eine einfache Anfrage mit `param_types`, die vor `cd03e09` beschädigt war, ist damit gültig. Laut `SPEC-041` trägt nur `parameter_description` dieses Feld. Andere Felder waren schon vorher nicht je Typ geprüft, etwa `tag` an `data_row`. | `SPEC-041`; [LH-FA-07](../../spec/lastenheft.md) Negative; Maintainability | `internal/adapters/driven/recording/yaml.go:281` (`responseDTO.ParamTypes`), `:536` (`responseFromDTO`) | ja — Sonden B, D | Feld ohne Typbindung im Leser |
| F-288 | LOW | Eine `parameter_description` ohne Parameter wird ohne den Schlüssel `param_types` geschrieben (Sonde C, `omitempty`). Für Client-Nachrichten schreibt derselbe Schreiber leere Felder aus (`param_types: []`). `SPEC-041` sagt, `parameter_description` „trägt die Typ-OIDs der Parameter als `param_types`“. Der Roundtrip bleibt gleich. Für den häufigen Fall `Describe` eines Statements ohne Parameter sieht die Datei aber anders aus, als die Spec-Formulierung erwarten lässt. | `SPEC-041`; Maintainability | `internal/adapters/driven/recording/yaml.go:281` | ja — Sonde C | leere Liste uneinheitlich serialisiert |
| F-289 | LOW | Der Kommentar von `interactionFromDTO` sagt: „jede Mischung der beiden Formen [ist] ein Fehler“. Geprüft wird `!= nil`, also werden `request: null` in einer Extended-Interaktion (Sonde E) sowie `responses: null` und `groups: null` angenommen; die letzten beiden folgen aus dem Lesen des Codes. Die Zusage ist weiter als die Prüfung. Laut Plan §8 nennen die Kommentare „nur, was geprüft wird“. Drittes Auftreten der Register-Klasse (Belege `slice-walking-skeleton-build-gates`, `slice-walking-skeleton-record`): Pflege-Schwelle des Skills. | `AGENTS.md` §3.7; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Maintainability | `internal/adapters/driven/recording/yaml.go:457` (Kommentar), `:465`, `:480` | ja — Sonde E | Zusage weiter als Code |
| F-290 | LOW | Die Commit-Message von `cd03e09` nennt die Struktur-IDs `SPEC-041`, `SPEC-001` und `ARC-001`. `AGENTS.md` §5 Regel 1 sagt, diese gehören nicht in die Commit-Message. Frühere Reviews haben das bei jedem Commit geprüft, bisher ohne Befund. | `AGENTS.md` §5 Regel 1 | Commit-Message `cd03e09` | ja (Lesen) | Struktur-ID in Commit-Message |
| F-291 | LOW | Plan §3 führt die Änderung an `slice-extended-query-record.md` als Zeile, die an `slice-extended-query-replay.md` nicht. Beide Folge-Slices sind im Commit geändert: Der Replay-Slice übernimmt in §6 das Risiko nil gegenüber leerer Liste. §1, §6 und der Kopf folgen dem Diff. Es fehlt nur die Tabellenzeile; das Muster von F-274 (§1 kleiner als der Diff) liegt nicht vor. | `AGENTS.md` §3.9; Maintainability | Slice-Plan §3 | ja (Lesen gegen Diff) | §3-Tabelle ohne Zeile für geänderte Planungsdatei |
| F-292 | INFO | Der Schreiber prüft nicht mit `Validate`. Sonde F: Eine Gruppe mit zwei `sync` wird ohne Fehler geschrieben, der Leser lehnt dieselbe Datei als `PGR-E3003` ab. Ein Fehler im künftigen Gruppieren fiele also erst beim Replay auf, nicht beim Record. Die DoD von `slice-extended-query-record` verlangt, dass jede aufgezeichnete Interaktion `Interaction.Validate` besteht. Hinweis an `slice-extended-query-record`. | [LH-FA-07](../../spec/lastenheft.md); [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) | `internal/adapters/driven/recording/yaml.go:119` (`Marshal`); `docs/plan/planning/open/slice-extended-query-record.md` §2 | ja — Sonde F | Schreiben ohne Formprüfung |
| F-293 | INFO | Plan §8 nennt einen „Bericht an den Reviewer“, in dem die rot färbende Mutation je Regel einmal gesehen ist. Bei diesem Review lag kein solches Artefakt vor. Eingang war nur der Commit. Die eigenen Mutationen M2 bis M13 bestätigen die Aussage für die Regeln von `Validate` und dem Leser. F-283 und F-284 liegen außerhalb davon. Ob die Aussage als Beleg für die DoD genügt, prüft der Verifier. | Modul 8 (kein Pfeil ohne benennbares Artefakt); Verweis Verifier | Slice-Plan §8 | ja (Lesen) | Übergabe ohne Artefakt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/model` — Core-Reinheit ([ADR-0006](../plan/adr/0006-kanonisches-domain-model.md)) | geprüft. `extended.go` importiert `errors` und `fmt`, `response.go` und `recording.go` nichts; keine YAML- oder Protokoll-Typen, keine Tags. `make a-check` 0 Befunde. |
| `internal/hexagon/model` — `Validate` gegen LH-FA-18.a/`SPEC-041` | geprüft. Mindestens eine Gruppe; jede Gruppe mit Client-Nachricht; `flush`/`sync` genau als letzte Nachricht; `sync` genau am Ende der letzten Gruppe; `ready_for_query` nur als letzte Server-Nachricht der letzten Gruppe und dort Pflicht; Flush-Gruppe ohne Server-Nachricht zulässig; Zielart von `describe`/`close` geprüft; kein SQL und keine `Responses` außerhalb der Gruppen; Mengen der Client- und Server-Nachrichten entsprechen LH-FA-18.a samt `ParameterStatus`. Jede Regel hat einen Negativfall mit Meldungstext (M3, M5, M7, M9, M11, M13 rot). Befund F-284. |
| `internal/adapters/driven/recording` — Schreiber | geprüft. Reihenfolge `type`, dann Felder aus `clientFelder`, auch leer (`TestExtendedLeereFelder`); deterministisch (zweites `Marshal` gleich); NULL- und Binärparameter nach `SPEC-003`; Ausgabe einfacher Anfragen byteidentisch zu `cd03e09~1`. Befund F-288, F-292. |
| `internal/adapters/driven/recording` — Leser | geprüft. `KnownFields` bleibt für alle DTOs außer `clientMessageDTO`; dort ersetzt die Whitelist je Typ die Prüfung (M4 rot). Doppelte Schlüssel lehnt `yaml.v3` ab (Test). Abgeschnittene Extended-Interaktionen fallen über `Validate` auf (M10 rot). Jeder Fehler wird `PGR-E3003` mit Session und Interaktion. Befund F-285, F-286, F-287, F-289. |
| `internal/hexagon/services` | geprüft. Nicht im Diff; die Typprüfung des Matchers greift (Lesen), ist aber ungetestet. Befund F-283. |
| Abdeckungs-Deklaration ([ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)) | geprüft, ohne Befund. Die neue Zeile LH-FA-07/Negative trifft das Kriterium „unvollständige oder beschädigte Datei“. LH-FA-18 ist nicht deklariert, so begründet Plan §1 das. `make abdeckung-check` grün. |
| Folge-Slices | geprüft. `slice-extended-query-record` nennt den gelieferten Format-Punkt in §1, die DoD verlangt `Validate` für jede aufgezeichnete Interaktion, §6 übernimmt die Feld-Grenze von `Validate`. `slice-extended-query-replay` übernimmt in §6 das Risiko nil gegenüber leerer Liste. Die Übergaben passen zu Plan §6. Befund F-291, F-292. |
| `spec/` — Strata und Hard Rule 3.4 | geprüft, ohne Befund. `spezifikation.md` nennt keine ADR, keinen Slice und keinen Commit; die Historienzeile nennt nur Spec-Kennungen. `ARC-001` nennt einen Code-Pfad, keine ADR, keine Sprach-Artefakte über den Pfad hinaus. Das Lastenheft ist unverändert. Abgleich der Inhalte: Befund F-285, F-286. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. `10cbdef` und `bb7133b` sind reine Moves, `4390713` und `d2421e8` ziehen Verweise nach, `cd03e09` verschiebt nichts. |
| Hard Rule 3.5, 3.8 — ADRs | geprüft, ohne Befund. Keine ADR im Diff; [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) bleibt unverändert, ihre Entscheidung (Gruppen je `Flush`/`Sync`) bildet das Modell ab. |
| Hard Rule 3.6 — Gates nicht lockern | geprüft, ohne Befund. Keine Gate-Konfiguration im Diff. |
| Hard Rule 3.7 — Kommentare | geprüft. Neue Kommentare sind Zusagen oder Abgrenzungen im Indikativ; keine verworfene Alternative, kein abwesender Text, kein abgebrochener Satz. Befund F-289. |
| Hard Rule 3.9 — Plan folgt Korrektur | geprüft. Kopf, §1, §2, §3, §6 und §8 folgen dem Umschnitt (Leser und Schreiber im Slice). Befund F-291. |
| Commit-Message `cd03e09` | geprüft. Trägt `slice-extended-query-modell`, [LH-FA-18](../../spec/lastenheft.md), [LH-FA-07](../../spec/lastenheft.md) und [ADR-0012](../plan/adr/0012-extended-query-gruppen.md). Befund F-290. |

## Summary

| Kategorie | Anzahl (neu) |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 5 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:**

- Protokollrand ohne Negativtest (F-283, F-284; Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, drittes Auftreten: Pflege-Schwelle)
- Spec-Wortlaut deckt Leser-Entscheidung nicht (F-285, F-286)
- Feld ohne Typbindung im Leser
- leere Liste uneinheitlich serialisiert
- Zusage weiter als Code (Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, drittes Auftreten: Pflege-Schwelle)
- Struktur-ID in Commit-Message
- §3-Tabelle ohne Zeile für geänderte Planungsdatei
- Schreiben ohne Formprüfung
- Übergabe ohne Artefakt

## Verdikt

**Merge-blockierend:** nein, es gibt kein HIGH. Modell und YAML-Format stimmen im Kern mit LH-FA-18.a und `SPEC-041` überein: Gruppenregeln, Sync/Flush, `ready_for_query`, Felder je Client-Nachricht beim Schreiben, Zielart. Die Regeln von `Validate` und dem Leser sind durch Mutationen belegt. Die Ausgabe einfacher Anfragen ist byteidentisch, die Core-Reinheit hält.

Vor der Closure verdienen die vier MEDIUM Aufmerksamkeit. F-283 und F-284 betreffen Zusagen am kritischen Pfad, die heute nur der Code hält und kein Test. F-285 und F-286 sind Leser-Entscheidungen, die der Wortlaut der Spezifikation nicht trägt.

**Übergabe:**

- F-283, F-284, F-287, F-288, F-289 und F-290 an den Implementer.
- F-285 und F-286 an den Implementer. Die Klärung liegt im Spec-Stratum 2 (`SPEC-001`, `SPEC-041`), dieser Slice berührt beide.
- F-291 an den Implementer (Plan §3).
- F-292 an `slice-extended-query-record` (Planner), F-293 an den Verifier.
- Steering-Loop: F-283/F-284 und F-289 erreichen je das dritte Auftreten ihrer Register-Klasse. Für beide geht an den Architect die Frage, ob ein Gate oder ein Eintrag in `AGENTS.md` sie verhindert hätte. Für Negativtests wäre das ein Mutationslauf, für Kommentar-Zusagen eine Regel zur Prüfung je Zusage.
- Widerspricht der Implementer einer MEDIUM-Einstufung, genügen Annahme oder Begründung. Der Konflikt-Pfad über den Architect gilt ab HIGH oder ab dem dritten gleichen Konflikttyp.
- Die Finding-Klassen gehen in die Slice-Closure §7. Die DoD-Konformität prüft der Verifier separat.
