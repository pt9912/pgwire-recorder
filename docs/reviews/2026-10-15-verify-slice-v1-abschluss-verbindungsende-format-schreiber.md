# Verifikation slice-v1-abschluss-verbindungsende-format-schreiber

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-15. **Spanne:** `6c6dcd9` (Basis) bis `827763c` (HEAD). **Gegenstand:** Slice `slice-v1-abschluss-verbindungsende-format-schreiber` (Plan unter `docs/plan/planning/in-progress/`, Schritt 35b) gegen DoD (§2), Plan §1 (Abgrenzung), §3, §4, §6 (Randformen, Mutationstabelle, Risiken), §7 (Belege des Implementers), §8, `SPEC-001` (*Schreiber in `yaml`*) und `SPEC-002`, `LH-FA-25` (.a), [`ADR-0039`](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md) und `AGENTS.md` §3.7 bis §3.14. Review: `docs/reviews/2026-10-15-review-slice-v1-abschluss-verbindungsende-format-schreiber.md` (0 HIGH, 0 MEDIUM, 3 LOW, 4 INFO); gegen ihn ist nicht verifiziert, die Nacharbeit (`827763c`) ist unabhängig nachgeprüft.

Eine Verifier-Vorlage führt die vendored Baseline nicht; Aufbau wie der Report von `slice-v1-abschluss-verbindungsende-format`. Belege sind in diesem Lauf selbst erzeugt: Quelle von `HEAD` (`git archive`: `internal`, `cmd`, `test`, `go.mod`, `go.sum`) in eine Kopie im Scratchpad, Läufe per Bind-Mount mit dem Image `pgwire-recorder:test` (`--network=none`, `go test` für `./internal/hexagon/model/` und `./internal/adapters/driven/recording/`), Mutanten einzeln in der Kopie angewandt und zurückgesetzt (am Ende `diff` gegen die Originale: gleich). Ein Lauf zur Zeit, kein `pkill`, kein `make gates`. Der Arbeitsbaum des Repos ist unverändert; `git status` zeigt nur diesen Report.

**Verdikt: PASS mit Befunden.** Alle Liefer-Punkte der DoD sind bestätigt, jeder nachgefahrene Mutant färbt genau die in §7 genannten Fälle (Abweichungen unten, alle INFO), die Abgrenzung ist eingehalten, der Plan folgt dem Diff, die Größe liegt weit unter Meldeschwelle und Marke. 0 HIGH, 0 MEDIUM, 0 LOW, 5 INFO. Offen sind allein die Closure-Pflichten (Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen), der erwartete Stand vor der Closure.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren oder gelesen) |
|---|---|---|
| 1: `Recording.HasEnd()` wahr bei einem Ende irgendwo; Schreiber schreibt 2 mit Ende, sonst 1, nach dem Inhalt (Feld 1 mit Ende schreibt 2, Feld 2 ohne Ende schreibt 1); ohne Sessions / nur `empty_sessions` 1; Schreiber prüft die Form nicht (Test je Stelle); vorhandene Tests grün | bestätigt | Code gelesen: `recording.go` (`HasEnd` über `Interaction.endet`, Responses und Gruppen), `yaml.go` (`formatVersion`: `model.VersionMitEnde` / `model.Version`, kein Literal, kein Lesen von `rec.Version`; `toDTO` ruft sie). Unveränderte Kopie: `go test -run 'Schreiben\|Roundtrip\|HasEnd\|VersionMitEnde'` grün; die volle Paketfolge der Mutanten-Läufe ohne Mutation (37 s) grün, also auch alle vorhandenen Tests. Mutanten: siehe unten. |
| 2: Roundtrip (Laden, Schreiben, Laden) für Ende einfach, `Sync`-Gruppe, `Flush`-Gruppe, `offset_ms`, `empty_sessions`; ohne Ende unverändert `version: 1` (Modell im Test, Text für Darstellung und den Fall ohne Ende) | bestätigt | `TestRoundtripMitEnde` (sechs Formen, je zwei Durchgänge), `TestSchreibenEndeDarstellung/ohne_Ende,_Bytes_bleiben` (Text, Byte-gleich). Mutanten je Form rot (unten). Die Fixtures sind Modelle; §1 und DoD nennen das jetzt (Review F-2 nachgezogen). |
| 3: je Zusage eine Mutation, vom Implementer rot gesehen; Tabelle in §7 | bestätigt, 34 Mutanten nachgefahren | Tabelle liegt in §7 (Zusage · Mutation · roter Test). Alle Zeilen, die ich gefahren habe, färben, was §7 sagt; Abweichungen sind V-1 und V-3. |
| 4: `make gates` grün | durch den Hook getragen, nicht selbst gefahren | §7: Exit 0 (ohne Flanke), davor `make lint` einzeln (erster Lauf rot, behoben). `bash .claude/hooks/stop-require-gates.sh` mit `{}` gibt `{"decision":"approve"}` (Exit 0). Der volle Nachweis (`gates-passed.dateiliste`, 15:46) liegt vor dem Commit der Nacharbeit (15:54); der Hook gibt frei, weil die schnelle Kette (`gates-schnell-passed.diffsha`, 15:54) für den Arbeitsbaum grün war. Ob der volle Lauf aus dem richtigen Grund grün war, belegt die Datei nicht; `make abdeckung-check` und `make kopf-check` selbst gefahren: beide Exit 0. |
| 5: Review durchgeführt, Report liegt vor | bestätigt | Report unter `docs/reviews/` (s. o.); F-1 bis F-7 in `827763c` bearbeitet (unten). Kein Self-Review: Der Report stammt aus einem eigenen Lauf. |
| 6 bis 9: Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen | offen (erwartet) | §7 trägt die Belege des Implementers, noch keine Closure-Notiz; §2 ohne Häkchen; die Risiken in §6 tragen „offen bis Closure“. Handlung des Planners; Vorschläge unten. |

## Mutationen selbst nachgefahren

Rot sind die genannten Fälle; „nur“ heißt: keine anderen Tests rot. Kürzel: SV = `TestSchreibenVersionNachInhalt`, NF = `TestSchreibenNimmtFormverletzungAn`, SD = `TestSchreibenEndeDarstellung`, RT = `TestRoundtripMitEnde`.

| Mutant | Rot | Urteil |
|---|---|---|
| Feld 1 mit Ende schreibt 1 (`HasEnd() && Version != 1`) | nur `SV/Feld_1_mit_Ende` | F-1 behoben, wie §7 |
| Feld 2 ohne Ende schreibt 2 | `SV/Feld_2_ohne_Ende`, `RT/Version_2_ohne_Ende` | wie §7 (dasselbe Element) |
| Feld 0 mit Ende schreibt 1 / Feld 0 ohne Ende schreibt 2 | je nur `SV/Feld_0_mit_Ende` / `SV/Feld_0_ohne_Ende` | wie §7 |
| Feld 3 mit Ende schreibt 1 / Feld 3 ohne Ende schreibt 2 | je nur `SV/Feld_3_mit_Ende` / `SV/Feld_3_ohne_Ende` | wie §7 |
| negatives Feld mit Ende schreibt 1 (`Version >= 0`) / ohne Ende schreibt 2 (`Version < 0`) | je nur `SV/negatives_Feld_mit_Ende` / `…_ohne_Ende` | F-5 behoben, wie §7 |
| Feld 2 mit Ende schreibt 1 | fünf Fälle `SV` (mit Ende, 2 Sessions, Gruppe, einfach) und alle NF-, SD-, RT-Fälle mit Ende | nicht isoliert, weil Feld 2 der Träger aller übrigen Fixtures ist (V-3) |
| immer `model.Version` | alle Fälle mit Ende in SV, NF, SD, RT | wie §7 |
| immer `model.VersionMitEnde` | alle Fälle ohne Ende (SV, SD „Bytes bleiben“, RT „Version 2 ohne Ende“, `TestLesenRoundtrip`, `TestRoundtrip`, `TestExtendedRoundtrip`) | wie §7 |
| ohne Ende: Feld übernehmen | `SV/Feld_2_/Feld_0_/Feld_3_/negatives_Feld_ohne_Ende`, `RT/Version_2_ohne_Ende` | wie §7 |
| `VersionMitEnde = 3` | `TestVersionMitEnde`, alle Fälle mit Ende in SV und NF | wie §7 |
| `return 2` statt `model.VersionMitEnde` | keiner, alle grün | äquivalent, wie §7 (Schreibweise ist nicht beobachtbar) |
| `HasEnd` nur erste / nur letzte Session | `TestHasEnd` (zwei bzw. zwei Fälle), `SV` je ein Fall; „nur letzte“ zusätzlich `NF/richtige_Stelle,_nicht_letzte_Session`, `RT/empty_sessions` | wie §7 |
| `HasEnd` nur erste / nur letzte Interaktion | `TestHasEnd/Ende_in_der_letzten_Interaktion` bzw. `…_in_einer_nicht_letzten_Interaktion` und die dazu passenden SV-, NF-, RT-Fälle | wie §7 |
| `HasEnd` wahr ohne Sessions | `TestHasEnd/ohne_Sessions`, `SV/ohne_Sessions`, `TestLesenRoundtrip/ohne_Sessions` | wie §7 |
| `Interaction.endet` nur `Groups` / nur `Responses` | `TestHasEnd` die Fälle des jeweils anderen Ortes, dazu SV, NF, SD, RT der Fälle des Ortes | wie §7 |
| Schreiber lehnt ab: Ende nicht letzte Nachricht | nur `NF/Anfang_der_Folge`, `NF/Mitte_der_Folge` | wie §7 |
| Schreiber lehnt ab: Ende an Index 0 | `NF/Anfang_der_Folge` und `SD/mit_Feld_neben_type` (das Ende ist dort die einzige Antwort) | unschädlich (V-4) |
| Schreiber lehnt ab: `ready_for_query` neben dem Ende | nur `NF/ready_for_query_neben_dem_Ende` | wie §7 |
| Schreiber lehnt ab: zwei Enden in einer Session | nur `NF/zwei_Enden_in_einer_Session` | wie §7 |
| Schreiber lehnt ab: Ende in nicht letzter Session | `SV/Ende_nur_in_der_ersten_von_zwei_Sessions`, `NF/richtige_Stelle,_nicht_letzte_Session`, `RT/empty_sessions` | wie §7 (die richtige Stelle lädt) |
| Roundtrip: `empty_sessions` mit Ende verloren | nur `RT/empty_sessions` | wie §7 |
| Roundtrip: `offset_ms` der Interaktion mit Ende verloren | nur `RT/offset_ms` | wie §7 |
| Roundtrip: Ende in `Sync`-Gruppe verloren | `RT/Sync-Gruppe`, `NF/richtige_Stelle,_letzte_Gruppe`, `SD/in_server_einer_Gruppe` | wie §7 (je Form der eigene RT-Fall) |
| Roundtrip: Ende in `Flush`-Gruppe verloren | `RT/Flush-Gruppe`, `NF/nicht_letzte_Gruppe` | wie §7 |
| Roundtrip: Ende in einfacher Interaktion verloren | `RT/einfach`, `RT/offset_ms`, `RT/empty_sessions` (gleiche Bauform), NF- und SD-Fälle | wie §7 |
| Darstellung: Feld neben `type` verworfen | nur `NF/Feld_neben_type`, `SD/mit_Feld_neben_type` | wie §7 |
| Ohne Ende: ein Zeilenumbruch mehr | nur `SD/ohne_Ende,_Bytes_bleiben` | wie §7 |
| **Zweiter Durchgang:** Feld 1 ohne Ende schreibt 2 (`HasEnd() \|\| Version == 1`) | `RT/Version_2_ohne_Ende` mit der Meldung „Laden, Schreiben, Laden weicht ab“ (erster Durchgang grün, weil das gebaute Feld 2 trägt); dazu `SV/Feld_1_ohne_Ende`, `SV/ohne_Sessions`, `SV/nur_empty_sessions`, `SD/ohne_Ende,_Bytes_bleiben` und die Lesertests | siehe V-1 |

Nicht nachgefahren: Flow-Form und Einrückung des Endes (Review und Implementer haben sie gefahren; die Fälle in SD sind Byte-genaue Texte).

## Randformen und Entscheidungen (§3.12)

| Prüfung | Ergebnis |
|---|---|
| Entscheidung vor dem Code | `4940cce` (Architect: `SPEC-001` *Schreiber in `yaml`*, `SPEC-002`-Blockform, Plan §6, Bindung an `-record`) liegt vor `48d50a3` (Code und Tests). Reihenfolge eingehalten. |
| Stellen je Fall (§6: Anfang, Mitte, richtig, letzte Gruppe, nicht letzte Gruppe, nicht letzte Interaktion, nicht letzte Session) | ein Fall in NF je Stelle, jeder erwartet Version 2 ohne Fehler beim Schreiben, `PGR-E3003` beim Laden (richtige Stellen laden); Mutanten oben zeigen den isolierten Fall je Stelle (einzige Überschneidungen: „nicht letzte Interaktion“ gegen „zwei Enden“, prinzipbedingt, wie im Review). |
| Randformen des Diffs gegen §6 | Der Diff entscheidet keine Randform, die §6 nicht nennt (Stichprobe: Feld 0/3/negativ, `empty_sessions`, `connection_closed` mit Feld, Version 2 ohne Ende, Byte-Gleichheit ohne Ende). Keine Rückgabe an den Architect nötig. |
| Frage: ändert die Straffung von `SPEC-001` in `827763c` („liest die Aufzeichnung einmal je Schreiben“ gestrichen) eine Entscheidung des Architect? | **Nein.** `4940cce` führt den Halbsatz als Beschreibung der Umsetzung (`toDTO` ruft `HasEnd` einmal) in der Aussage, die Zusagen seien Zustandsprüfungen; entschieden sind darin nur: die Version folgt dem Inhalt, die Stelle zählt nicht, die Form prüft der Aufrufer, Darstellung, Byte-Gleichheit, Roundtrip. Der Rest des Satzes bleibt („sagt nichts über die Reihenfolge oder Gleichzeitigkeit mehrerer Aufrufe zu“), und er trägt die Randform „Aufruf-Voraussetzung: keine“ aus §6 weiter. Eine Zusage ohne Test enger zu fassen ist der Weg, den `AGENTS.md` §3.11 nennt. Formal bleibt es eine Änderung des Technik-Stratums durch den Implementer (V-2). |
| Aufrufarten des Treibers | keine (kein Hook oder Gate liest ein Ergebnis dieses Schreibers); §6 sagt es, der Diff berührt keine. |

## Beurteilung der genannten Grenze: zweiter Roundtrip-Durchgang

Frage: Der zweite Durchgang (Laden, Schreiben, Laden) hat nach §7 keine eigene Mutation. Genügt das `AGENTS.md` §3.10 bzw. §3.11?

- **§3.11 (Zusage nur, was ein Test prüft): erfüllt.** Die Deklaration sagt „Schreiben und Laden einer Aufzeichnung und, von der geladenen Aufzeichnung aus, Schreiben und erneutes Laden“; der Test tut genau das (`schreibeUndLade(t, got)`, Vergleich `again` mit `got`). Die frühere Fassung („Laden, Schreiben …“ vom gebauten Modell aus) ist als F-3 behoben. Der Satz in `SPEC-001` (Laden, Schreiben und erneutes Laden) ist durch denselben Test gedeckt, auch für die Datei der Version 2 ohne Ende (zweiter Durchgang liest Version 1).
- **§3.10 (je Zusage eine Mutation): erfüllt, und die Aussage in §7 ist zu streng.** §7 sagt, es gebe keinen Mutanten, den nur der zweite Durchgang färbt, „jeder Mutant färbt schon im ersten Durchgang“. Das stimmt nicht für jeden: „Feld 1 ohne Ende schreibt 2“ lässt den ersten Durchgang von `RT/Version_2_ohne_Ende` grün (das gebaute Feld ist 2, der Schreiber schreibt 1) und wird im zweiten Durchgang rot (die geladene Aufzeichnung trägt Feld 1). Selbst gefahren (Meldung „Laden, Schreiben, Laden weicht ab“). Der zweite Durchgang hat also eine Mutation, die in diesem Test nur er färbt. Dass derselbe Mutant im Gesamtlauf zusätzlich von `SV/Feld_1_ohne_Ende` und anderen gefangen wird, ändert daran nichts; nach §3.10 genügt eine rote Prüfung je Zusage.
- **Wo der zweite Durchgang tatsächlich ohne Mutation bleibt:** für Aufzeichnungen mit Ende. Die geladene Form unterscheidet sich von der gebauten dort in nichts außer dem Feld (beide 2); ein Schreiber, der geladene und gebaute Aufzeichnungen verschieden behandelte, ist nicht konstruierbar, ohne ein beobachtbares Merkmal zu lesen (Startup und Sessions sind in beiden Formen gefüllt). Das ist die Grenze, die §7 meint, und sie ist ehrlich benannt; sie bricht keine Zusage.
- **Folge:** keine Handlung am Code. Der Satz in §7 vor der Closure berichtigen (V-1).

## Abgrenzung, Plan und Größe

| Prüfung | Ergebnis |
|---|---|
| §1 Abgrenzung | `git diff --name-only 6c6dcd9..HEAD` außerhalb von `docs/plan/` und `docs/reviews/`: `internal/hexagon/model/{recording.go,hasend_test.go}`, `internal/adapters/driven/recording/{yaml.go,yaml_schreiben_ende_test.go}`, `docs/user/abdeckung-unit.md` (Werkzeug-Ausgabe), `spec/spezifikation.md` (Architect, `4940cce`; Straffung `827763c`). Kein Handbuch (`docs/user/benutzerhandbuch.md`), kein `README.md`, kein SQLite-Code (`sqlite*.go` nicht im Diff, `model.Version` unverändert), kein `-record`-/`-upstream`-Code, kein `tools/schema/schema.yaml`, kein Modell der Formregeln und kein Leser (`yaml.go` ändert nur `formatVersion` und `toDTO`). Eingehalten. |
| §3 (Plan) gegen `git diff --stat` | jede Datei des Diffs steht in §3 (fünf Zeilen); `abdeckung-e2e.md` bleibt, es ist kein E2E-Test dabei. Stimmt. |
| §6, §7 gegen Diff | Tabelle, Läufe und Grenzen folgen dem Diff. Abweichung nur V-1 (Aussage zum zweiten Durchgang) und V-3 (Feld 2 mit Ende). |
| Größe (Metrik §4: `git diff --numstat 6c6dcd9..HEAD -- internal docs/user`) | Einfügungen 319, Löschungen 1, Summe **320** (7 + 12 + 236 + 51 + 14 selbst addiert aus den fünf Zeilen), gleich dem Wert des Implementers. Meldeschwelle 550 und Marke 700 nicht erreicht. Schätzung des Architect 340 bis 500: leicht unterschritten (V-5). Liefer-Punkte 2, Schichten zwei (Kern, Recording-Adapter). |
| §8 | Zählung (Planner) und Messung (Implementer) stehen; die Bindung an `-record` trug der Architect in dessen §1, §6, §8 ein (`4940cce`), die Messzeile nennt es (F-6 nachgezogen). |
| `AGENTS.md` §3.9 | Plan folgt dem Diff in §1 (Modell im Test), §3, §6, §7, §8; Kopf (`Bezug`, `Berührte Spec-Stellen`) trägt `LH-FA-25`, `LH-FA-22`, `LH-QA-06`, `SPEC-001`, `SPEC-002`, `ARC-001`, `ARC-008`; `make kopf-check` Exit 0. |
| `AGENTS.md` §3.7, §3.11 | Kommentare von `HasEnd`, `formatVersion` und den Testkopfzeilen sagen nur, was Tests prüfen (Indikativ, keine Chronik, keine Konjunktive); die Abdeckungs-Zeilen in `abdeckung-unit.md` entsprechen den Deklarationen (6 Zeilen, `make abdeckung-check` Exit 0). |
| `AGENTS.md` §3.13 (Bindungen) | `-record`: Bindung in Kopf, §1 (Zeile 39), §6 (Zeile 135) und §8 (Zeile 178), mit Kennung des Gebers. `-sqlite` und `-upstream` nennen den Schreiber als Voraussetzung („Übernommen aus … format-schreiber“, Start, §8 nachgezählt); beide liegen in `next/`, nehmen an. Die Bindung „Schreiber prüft die Form nicht“ liegt dort über `-record` (der Use Case sichert die Form), wie §8 sagt; das trägt, solange `-record` vor beiden schließt (Schritt-Reihenfolge der Welle). |
| Lifecycle | Der Slice liegt in `in-progress/`; `git mv` (`0cdcf84`) und Verweise (`f559dfb`) sind eigene Commits; Roadmap „In Arbeit“ stimmt. |

## Findings

### V-1 INFO

- `quelle`: `AGENTS.md` §3.9, §3.11
- `pfad`: `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-format-schreiber.md` · §7, Tabellenzeile „Laden, Schreiben, erneutes Laden (Nacharbeit `F-3`)“
- `befund`: Die Zeile sagt, der zweite Durchgang habe keinen Mutanten, „den nur dieser zweite Durchgang färbt“, weil „jeder Mutant schon im ersten Durchgang“ rot werde. Gefahren: „Feld 1 ohne Ende schreibt 2“ lässt den ersten Durchgang von `RT/Version_2_ohne_Ende` grün und färbt dort allein den zweiten (Meldung „Laden, Schreiben, Laden weicht ab“). Die Grenze gilt für Aufzeichnungen mit Ende, nicht für die Datei der Version 2 ohne Ende.
- `verifizierbar`: ja (Mutant in der Kopie)
- `klasse`: Zusage in §7 weiter als die Prüfung (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, kleine Form)
- Vorschlag: Die Zeile vor der Closure ersetzen durch „Mutation: Feld 1 ohne Ende schreibt 2 · rot: `RT/Version_2_ohne_Ende`, nur im zweiten Durchgang; für Aufzeichnungen mit Ende kein Mutant, der nur der zweite Durchgang färbt (geladene und gebaute Form gleich bis aufs Feld, beide 2)“.

### V-2 INFO

- `quelle`: `AGENTS.md` §3.12 (Spezifikation im Technik-Stratum), Modul 8 (Rollen)
- `pfad`: `spec/spezifikation.md` · `SPEC-001` *Schreiber in `yaml`*, Satz „Die Zusagen sind Zustandsprüfungen: …“
- `befund`: Die Straffung (Halbsatz gestrichen) stammt vom Implementer, die Fassung vom Architect. Sie ändert keine Entscheidung (Tabelle oben), ist aber eine Änderung am Technik-Stratum außerhalb von dessen Rolle. Der Satz steht jetzt als eine lange Zeile mit anschließendem Umbruch („… mehrerer\nAufrufe zu. Ein Leser“); Markdown rendert es gleich, es liest sich uneben.
- `verifizierbar`: ja (`git show 827763c -- spec/spezifikation.md`)
- `klasse`: Rollengrenze, formal
- Vorschlag: Der Architect bestätigt die Straffung bei der Closure (ein Satz in §7 genügt); den Zeilenumbruch beim nächsten Berühren der Stelle glätten.

### V-3 INFO

- `quelle`: `AGENTS.md` §3.10 (Aufzählung zählt je Element)
- `pfad`: `internal/adapters/driven/recording/yaml_schreiben_ende_test.go` · `TestSchreibenVersionNachInhalt` · Fälle „Feld 2 mit Ende“ und die Fixtures mit Ende
- `befund`: §7 sagt, die Mutanten der Aufzählung färbten „je genau ein Fall“. Das gilt für die Elemente von §6 (1 mit Ende, 2 ohne Ende, 0 und 3 je mit und ohne Ende, negativ je mit und ohne Ende), alle selbst gefahren. „Feld 2 mit Ende schreibt 1“ färbt dagegen alle Fälle mit Ende (Feld 2 ist der Träger der übrigen Fixtures). Das Element steht nicht in §6, die Zusage „Version folgt dem Inhalt, nicht dem Feld“ ist durch die anderen Felder isoliert; es ist kein Mangel, nur die Grenze der Isolierung.
- `verifizierbar`: ja
- `klasse`: Element einer Aufzählung ohne isolierenden Fall (Grenzfall, kein Befund)
- Vorschlag: keine Handlung; wer die Zeile „je genau ein Fall“ in §7 liest, sieht „Ausnahme `Feld_2_ohne_Ende`“ schon; „Feld 2 mit Ende“ ergänzen oder streichen.

### V-4 INFO

- `quelle`: `AGENTS.md` §3.10
- `pfad`: `yaml_schreiben_ende_test.go` · `NF/Anfang_der_Folge`, `SD/mit_Feld_neben_type`
- `befund`: Ein Schreiber-Mutant, der nur „Ende an Index 0“ ablehnt, färbt neben dem Fall „Anfang der Folge“ auch `SD/mit_Feld_neben_type`, weil dort das Ende die einzige Antwort ist. Die Stelle „Anfang“ ist darum nicht streng isoliert. Für die Zusage unerheblich (der Mutant lehnt einen Fall ab, den der Darstellungstest ebenfalls annimmt), und der Mutant ist nicht der, den §7 für „Anfang“ nennt (Ende nicht letzte Nachricht färbt nur Anfang und Mitte).
- `verifizierbar`: ja
- `klasse`: Überschneidung zweier Fälle, akzeptiert
- Vorschlag: keine Handlung.

### V-5 INFO

- `quelle`: `AGENTS.md` §3.14
- `pfad`: Plan §4, §8 · Schätzung des Architect 340 bis 500 (zuvor 300 bis 500)
- `befund`: Gemessen 320 (311 vor der Nacharbeit): unter der unteren Grenze der Schätzung. Die Schätzung mit dem Faktor 1,3 bis 1,5 für Tabellen-Tests lag diesmal darüber, die Messung nicht darunter wie in `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`. Ein Datenpunkt für `slice-harness-schaetzung-planner-faktor` (Schritt danach in der Roadmap): Produktcode 26, Tests 287, Tabelle 7, Verhältnis Tests zu Produktcode rund 11 bei einem Produktcode, den der Architect mit rund 40 ansetzte.
- `verifizierbar`: ja (`git diff --numstat`)
- `klasse`: Schätzung gegen Messung
- Vorschlag: als Beobachtung ins Register, wenn der Planner den Faktor-Slice schneidet; kein Rücksprung.

## Risiko-Ausgänge (Vorschläge für die Closure)

| Risiko (§6) | Vorschlag | Begründung |
|---|---|---|
| Zwischenstand: Der Schreiber kann `version: 2` schreiben, aber kein Produktpfad erzeugt ein Ende | *eingetreten* → Folge-Slices `slice-v1-abschluss-verbindungsende-record` und `slice-v1-abschluss-verbindungsende-upstream` | Der Zwischenstand besteht (nur der Schreiber ist geliefert). Beide Adressen nehmen an: `-record` führt das Aufzeichnen des Endes im Use Case und die Bindung des Schreibers (§1, §6), `-upstream` nennt den Schreiber als Voraussetzung und zeichnet in YAML auf. Beide liegen in `next/`. |
| Die Wahl der Version ist Voraussetzung des Roundtrips beider Formate in `-sqlite` und des Abnahmeszenarios 18 in YAML in `-upstream` | *entfallen*, mit Begründung (alternativ *eingetreten* → `-sqlite`, `-upstream`, wie §6 erwartet) | Die Voraussetzung ist mit dem Schreiber geliefert und belegt (`TestSchreibenVersionNachInhalt`, `TestRoundtripMitEnde`); beide Nehmer tragen den Slice in §1 und im Start (`done/` als Bedingung). Als Risiko bleibt nichts offen; wer die Lesart „eingetreten“ der Erwartung in §6 folgt, nennt dieselben zwei Kennungen. |
| Der Diff liegt über der Marke (Schätzung 300 bis 500) | *entfallen*, mit Begründung | Gemessen 320 in der Metrik von §4, unter Meldeschwelle 550 und Marke 700; Liefer-Punkte 2, Schichten zwei. |

Weitere Punkte für die Closure:

- Review F-1 bis F-7 sind in `827763c` bearbeitet; die nachgefahrenen Stichproben (F-1, F-3, F-5, F-7) stimmen mit §7 überein. F-4 bleibt als äquivalenter Mutant `return 2` dokumentiert (selbst gefahren: grün).
- Lerneintrag, Kandidaten (Urteil des Planners):
  1. **Element einer Aufzählung ohne isolierenden Fall** (Review F-1, F-5): trat trotz der Erweiterung von §3.10 („Aufzählung und Ortsangabe zählen je Element“, seit `slice-harness-gate-kette-stop-hook-liste`) in der ersten Fassung wieder auf, weil Fixtures das Feld der Aufzählung als Träger nutzten. Zähler `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`; Form: geschärfte Regel („das Feld einer Aufzählung trägt keine Fixture, die ein anderes Element prüft“) oder Beobachtung unter dem Zähler; der Retirement-Check von §3.10 ist damit ausgelöst (Beobachtung ist wieder aufgetreten).
  2. **Zusage weiter als der Test** (Review F-3, hier V-1): Deklaration „Laden, Schreiben …“ ohne Laden; dann die Aussage in §7 zum zweiten Durchgang. `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, zwei Auftreten in einem Slice.
  3. **Plan folgt dem Diff** (Review F-2, F-6): §1/DoD „Text“ gegen §3 „Modell“, §8 „kein Nehmer“; `BEO-REPO/plan-folgt-korrektur-nicht`, ein Auftreten.
  4. **Nicht aufgetreten:** `BEO-REPO/spec-randform-erst-im-review-entschieden` (kein Befund zu einer Randform im Review, Randformen standen vor dem Code); `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` und `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Metrik genannt, Messung unter der Schätzung). Nützlich für die Retirement-Checks von §3.12 und §3.14 (ein Slice ohne Auftreten).
  5. **Spec-Lücke:** „Zusage in der Spezifikation ohne Test“ (Review F-7): eine Beschreibung der Umsetzung („liest einmal je Schreiben“) im Technik-Stratum; vorgeschlagene Form: benannte Lücke oder Beobachtung, kein neuer Sensor.
- Hinweis zur Kenntnis: Der Implementer meldete einen versehentlichen `pkill` in einem Hilfsbefehl, der nur die eigene Shell traf. Ohne Wirkung auf Arbeitsbaum, Gates oder Prozesse anderer Läufe nachweisbar; `AGENTS.md`/Projektregel „kein `pkill`“ bleibt unberührt, eine Handlung ist nicht nötig.

## Gates und Zustand

- Eigene Läufe: unveränderte Kopie von `HEAD`, `go test` der beiden Pakete (gefiltert und, in den Mutanten-Läufen ohne Mutation, ungefiltert) grün; 34 Mutanten einzeln, alle wie in §7 behauptet rot, `return 2` grün (äquivalent); dazu `make abdeckung-check` und `make kopf-check` (Exit 0) und der Stop-Hook (`approve`).
- Voller `make gates` nicht gefahren; `make lint` nicht gefahren. §7 trägt den Lauf; der Hook gibt auf Grund der schnellen Kette frei (siehe DoD-Punkt 4).
- Aufräumen: Kopie und Skripte im Scratchpad, keine Container und kein Netz angelegt (`docker run --rm`), `git status` zeigt nur diesen Report (nicht committiert).
