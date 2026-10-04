# Mutationen: slice-extended-query-modell — 2026-10-04

**Zweck:** Beleg der Implementer-Rolle für den Verifier: welche Änderung am geprüften Code jede Zusage des Slice rot macht, und dass sie einmal rot gesehen wurde.

**Vorgehen:** Jede Mutation einzeln im Arbeitsbaum angewandt (`perl` auf genau eine Stelle), `make test` gelaufen (`go vet` und Unit-Tests im `Dockerfile`, Stufe `test`), Datei danach aus der Sicherung zurückgeschrieben. Ein Kontrolllauf ohne Mutation war jeweils grün. „Rot“ heißt: `make test` endet mit Exit-Code ungleich 0, und der genannte Test schlägt fehl, nicht ein Übersetzungsfehler.

**Stand:** M1 bis M11 gegen den Code von Commit `cd03e09`; N1 bis N10 gegen die Fassung nach den Review-Findings F-283 bis F-289 (Commit `6b4967d`); O1 bis O9 gegen die Fassung nach den Verifikations-Befunden V-16 bis V-18. Mit dieser Fassung heißt die Vorprüfung `formVorpruefung` (vorher `formSchluesselNull`), und die Ablehnung von `param_types` an anderen Antworten liegt dort statt in `responseFromDTO`; N7 und N10 treffen seither `formVorpruefung` (O1, O8).

## Erste Runde (Formregeln und Leser)

| # | Mutation | Ort | Rot in |
|---|---|---|---|
| M1 | `ready_for_query` überall in einer Gruppe zugelassen | `internal/hexagon/model/extended.go` (`Group.validate`) | `TestValidateFehler` (ready_for_query in der Flush-Gruppe, vor dem Ende der letzten Gruppe) |
| M2 | `parse_complete` in die Menge der einfachen Antworten aufgenommen | `internal/hexagon/model/response.go` (`simpleResponses`) | `TestValidateQueryFehler` (Extended-Antwort in einfacher Anfrage) |
| M3 | `parameter_status` aus der Extended-Menge entfernt | `internal/hexagon/model/response.go` (`extendedResponses`) | `TestValidateGueltig`, `TestExtendedRoundtrip` |
| M4 | Regel „sync genau am Ende der letzten Gruppe“ abgeschaltet | `internal/hexagon/model/extended.go` | `TestValidateFehler` (vordere Gruppe mit sync, letzte mit flush), `TestUnmarshalExtendedFehler` (abgeschnitten nach der Flush-Gruppe) |
| M5 | Schlüsselprüfung je Client-Nachricht im Leser abgeschaltet | `internal/adapters/driven/recording/yaml.go` (`clientMessageDTO.UnmarshalYAML`) | `TestUnmarshalExtendedFehler` (fremder Schlüssel an sync, Schlüssel von bind an parse) |
| M6 | Schreiber gibt alle Felder jeder Client-Nachricht aus | `internal/adapters/driven/recording/yaml.go` (`clientMessageDTO.MarshalYAML`) | `TestExtendedRoundtrip`, `TestExtendedLeereFelder` |
| M7 | Leser liest leere Listen nicht als nil | `internal/adapters/driven/recording/yaml.go` (`leerAlsNil`) | `TestExtendedLeereFelder` |
| M8 | Leser ruft `Validate` nicht | `internal/adapters/driven/recording/yaml.go` (`fromDTO`) | `TestUnmarshalFehler`, `TestUnmarshalExtendedFehler` (12 Fälle) |
| M9 | einfache Anfrage mit `groups` angenommen | `internal/adapters/driven/recording/yaml.go` (`interactionFromDTO`) | `TestUnmarshalExtendedFehler` (einfache Anfrage mit groups) |
| M10 | Zielart von describe/close nicht geprüft | `internal/hexagon/model/extended.go` | `TestValidateFehler`, `TestUnmarshalExtendedFehler` |
| M11 | `flush`/`sync` mitten in der Gruppe bzw. Gruppe ohne Abschluss zugelassen | `internal/hexagon/model/extended.go` | `TestValidateFehler`, `TestUnmarshalExtendedFehler` (abgeschnitten in den Client-Nachrichten) |

Zusätzlich: Ein Lauf, in dem der Schreiber nil-Listen nicht vorher durch leere ersetzte, blieb grün, weil yaml.v3 nil-Listen als `[]` schreibt. Der Kommentar, der das Gegenteil sagte, und der Code dazu sind entfernt.

## Zweite Runde (Review-Findings)

| # | Finding | Mutation | Ort | Rot in |
|---|---|---|---|---|
| N1 | F-283 | Typprüfung des Matchers entfernt (`erwartet.Request.Type != model.RequestQuery`) | `internal/hexagon/services/replay.go` | `TestReplayExtendedAmCursor` |
| N2 | F-284 | `error_response` aus `simpleResponses` entfernt | `internal/hexagon/model/response.go` | `TestValidateGueltig` (query mit jedem Antworttyp) |
| N3 | F-284 | `empty_query_response` aus `simpleResponses` entfernt | `internal/hexagon/model/response.go` | `TestValidateGueltig` |
| N4 | F-284 | `parameter_status` aus `simpleResponses` entfernt | `internal/hexagon/model/response.go` | `TestValidateGueltig` |
| N5 | F-285 | fehlendes Feld des Typs angenommen | `internal/adapters/driven/recording/yaml.go` (`clientMessageDTO.UnmarshalYAML`) | `TestUnmarshalExtendedFehler` (parse ohne sql, parse mit sql: null, bind ohne params, execute ohne max_rows, describe ohne name, close ohne target) |
| N6 | F-285 | Feld mit Wert `null` als belegt gezählt | ebenda | `TestUnmarshalExtendedFehler` (parse mit sql: null) |
| N7 | F-287 | `param_types` an jeder Antwort angenommen | `internal/adapters/driven/recording/yaml.go` (`responseFromDTO`) | `TestUnmarshalExtendedFehler` (param_types an parse_complete, in einfacher Anfrage) |
| N8 | F-288 | `param_types` an `parameter_description` nicht verlangt | ebenda | `TestUnmarshalExtendedFehler` (parameter_description ohne param_types) |
| N9 | F-288 | leere `param_types` einer `parameter_description` nicht geschrieben | `internal/adapters/driven/recording/yaml.go` (`responseToDTO`) | `TestExtendedLeereFelder` |
| N10 | F-289 | Vorprüfung auf `null` an `request`/`responses`/`groups` abgeschaltet | `internal/adapters/driven/recording/yaml.go` (Vorprüfung, heute `formVorpruefung`) | `TestUnmarshalExtendedFehler` (extended mit request: null, mit responses: null, einfache Anfrage mit groups: null) |

## Dritte Runde (Verifikations-Befunde)

| # | Befund | Mutation | Ort | Rot in |
|---|---|---|---|---|
| O1 | V-16 | Prüfung von `param_types` an anderen Antworten entfernt | `internal/adapters/driven/recording/yaml.go` (`formVorpruefung`) | `TestUnmarshalExtendedFehler` (param_types an parse_complete, in einfacher Anfrage, je auch mit null) |
| O2 | V-16 | `param_types: null` an anderen Antworten zugelassen | ebenda | `TestUnmarshalExtendedFehler` (param_types: null an parse_complete, in einfacher Anfrage) |
| O3 | V-17 | fehlendes `client`/`server` einer Gruppe zugelassen | ebenda | `TestUnmarshalExtendedFehler` (Gruppe ohne client, Flush-Gruppe ohne server, abgeschnitten nach den Client-Nachrichten) |
| O4 | V-17 | `client: null`/`server: null` zugelassen | ebenda | `TestUnmarshalExtendedFehler` (Gruppe mit client: null, Flush-Gruppe mit server: null) |
| O5 | V-18 | Prüfung von `type` auf null oder leer entfernt | ebenda | `TestUnmarshalExtendedFehler` (type null und leer, je an Extended und einfacher Anfrage) |
| O6 | V-18 | nur `type: null` abgelehnt, leerer `type` zugelassen | ebenda | `TestUnmarshalExtendedFehler` (type leer an Extended, an einfacher Anfrage) |
| O7 | V-17 | nur `server` geprüft, `client` nicht | ebenda | `TestUnmarshalExtendedFehler` (Gruppe ohne client, mit client: null) |
| O8 | V-16 | `param_types` nur in `responses`, nicht in den Gruppen geprüft | ebenda | `TestUnmarshalExtendedFehler` (param_types an parse_complete, auch mit null) |
| O9 | V-18 | nur leerer `type` abgelehnt, `type: null` zugelassen | ebenda | `TestUnmarshalExtendedFehler` (type null an Extended, an einfacher Anfrage) |

Gültige Gegenstücke im selben Test: Sync-Gruppe, Flush-Gruppe mit `server: []`, `parameter_description` mit und ohne Parameter.

## Nicht abgedeckt

- Felder einer Serverantwort außer `param_types` sind nicht je Typ gebunden (etwa `tag` an `data_row`); das galt schon vor diesem Slice.
- Anker und Aliase im YAML-Baum sieht `formVorpruefung` nicht aufgelöst; eine Aufzeichnung, die `null` über einen Alias setzt, prüft sie nicht.
