# Verifikation slice-v1-abschluss-verbindungsende-format

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-15. **Spanne:** `6d0d697` (Basis) bis `46bd845` (HEAD). **Gegenstand:** Slice `slice-v1-abschluss-verbindungsende-format` (Plan unter `docs/plan/planning/in-progress/`, Schritt 35a) gegen DoD (§2), Plan §1 (Abgrenzung), §3, §4, §6 (Randformen, Risiken), §7 (Belege des Implementers), §8, `SPEC-001` und `SPEC-002` (Randformen aus `690e353` und `aed8a49`), `LH-FA-25` mit `LH-FA-25.a`, [`ADR-0039`](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md) und `AGENTS.md` §3.9 bis §3.14. Review: `docs/reviews/2026-10-15-review-slice-v1-abschluss-verbindungsende-format.md`; gegen ihn ist nicht verifiziert, F-1 und F-4 sind unabhängig nachgeprüft.

Eine Verifier-Vorlage führt die vendored Baseline nicht; Aufbau wie der Report von `slice-v1-abschluss-protokollrand-doku`. Belege sind in diesem Lauf selbst erzeugt: Quelle von `HEAD` (`git archive`: `internal`, `cmd`, `test`, `go.mod`, `go.sum`) in eine Kopie im Scratchpad, Läufe per Bind-Mount mit dem Image `pgwire-recorder:test` (`--network=none`, `go test` für `./internal/hexagon/...` und `./internal/adapters/driven/recording/`), Mutanten einzeln in der Kopie angewandt und zurückgesetzt. Ein Lauf zur Zeit, kein `pkill`, kein `make gates`. Der Arbeitsbaum des Repos ist unverändert (`git status` zeigt nur diesen Report).

**Verdikt: PASS mit Befunden.** Alle Liefer-Punkte der DoD sind bestätigt, jeder nachgefahrene Mutant färbt genau die in §7 genannten Tests, die Abgrenzung ist eingehalten, der Plan folgt dem Diff. 0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO. Offen sind allein die Closure-Pflichten (Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen); das ist der erwartete Stand vor der Closure.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren oder gelesen) |
|---|---|---|
| 1: `connection_closed` ist Antworttyp; `Validate` nimmt ihn nur als letzte Server-Nachricht nach `error_response` `FATAL`/`PANIC` (V, ohne V S), ohne Felder; einfach statt `ready_for_query`, Extended in der letzten Gruppe mit `Sync`/`Flush` ohne `ready_for_query`; `Session.ValidateEnd` (Test) | bestätigt | Code gelesen (`response.go`: `ResponseConnectionClosed`, `pruefeEnde`, `ohneFelder`, `beendetVerbindung`, `ohneReadyForQuery`; `extended.go`: `validateQuery`, `Group.validate`, `validateAbschluss`; `recording.go`: `Session.ValidateEnd`). Unveränderte Kopie: `./internal/hexagon/...` grün (`model`, `services`). Tests `TestValidateEndeGueltig`, `TestValidateEndeFehler`, `TestSessionValidateEnd` rot unter den Mutanten unten. |
| 2: Leser nimmt `version: 1` und `2` an (auch 2 ohne Ende); `connection_closed` in Version 1, falsche Stelle, nicht nach FATAL/PANIC, mit Feldern ist `PGR-E3003`; `version: 3` und andere ganze Zahlen `PGR-E3002`; Leser ruft `Session.ValidateEnd` je Session; Tests mit Text, vorhandene Tests grün | bestätigt | `yaml.go` gelesen (`kopfVorpruefung`, `formVorpruefung`, `antwortVorpruefung`, `sessionFromDTO`). `./internal/adapters/driven/recording/` grün in der unveränderten Kopie (37 s). Mutanten Leser: siehe unten. Die drei vorhandenen Tests (`yaml_test.go`, `leser_test.go`, `zeitangaben_pruefen_test.go`) nehmen die 3 statt der 2; `TestZeitangabenVersionVorForm` führt den Fall „Version 2“ mit `PGR-E3003` aus dem Inhalt (offset_ms), der Leser nimmt Version 2 also an. |
| 3: je Zusage eine Mutation, vom Implementer rot gesehen; Tabelle in §7 | bestätigt, 18 Mutanten nachgefahren | Tabelle liegt in §7 (Spalten Zusage · Mutation · roter Test). Nachgefahren, alle wie behauptet: M1 bis M3, M4a, M8 (Stellen von `ohneReadyForQuery`), Stellenprüfung in `pruefeEnde`, `Interaction.endet`, `Group.endet`, `!letzte && g.endet()`, `pruefeEnde` je Art entfernt, `Session.ValidateEnd`, Aufruf im Leser, Version-1-Prüfung, Schlüssel neben `type`, Version 2 nicht angenommen; Ergebnisse im nächsten Abschnitt. Eine Zeile ist zu weit gefasst (V-1). |
| 4: `make gates` grün | durch den Hook getragen, nicht selbst gefahren | §7 *Letzter Lauf*: Exit 0 auf dem Arbeitsbaum der Nacharbeit; danach ändern nur `46bd845` (nur Plan) und der Commit der Nacharbeit mit demselben Stand. `bash .claude/hooks/stop-require-gates.sh` mit `{}` gibt `{"decision":"approve"}` (Exit 0): Der Hook hält den Stand für getragen (`.harness/state/gates-passed.dateiliste` 12:30, schnelle Kette 12:31). Ob der volle Lauf *aus dem richtigen Grund* grün war, belegt die Datei nicht; ich habe `make gates` nicht wiederholt. Eigene Läufe: die beiden Unit-Pakete grün, `go test` (Bind-Mount) wie in §7 beschrieben. |
| 5: Review durchgeführt, Report liegt vor | bestätigt | Report unter `docs/reviews/` (0 HIGH, 1 MEDIUM, 1 LOW, 3 INFO). F-1 (MEDIUM) durch Architect-Entscheidung `aed8a49` und Nacharbeit `91a5228` behoben (unten nachgeprüft), F-2 und F-4 behoben, F-3 beim Planner, F-5 keine Aktion. Kein Self-Review: Der Report stammt aus einem eigenen Lauf. |
| 6 bis 9: Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen | offen (erwartet) | §7 trägt Belege des Implementers, noch keine Closure-Notiz; §2 ohne Häkchen; die Risiken in §6 tragen „offen bis Closure“. Handlung des Planners; Vorschläge unten. |

## Mutationen selbst nachgefahren

Rote Tests aus dem Lauf `go test ./internal/hexagon/model/ ./internal/adapters/driven/recording/ -v` (Namen gekürzt). „nur“ heißt: keine anderen Tests rot.

| Mutant | Rot | Urteil |
|---|---|---|
| M1 `ohneReadyForQuery` überspringt Index 0 | nur `TestValidateEndeFehler/query_ready_for_query_am_Anfang,…`, `TestUnmarshalEndeAbgelehnt/ready_for_query_am_Anfang,…` | wie §7 |
| M2 überspringt die Mitte (Index 1 bis drei vor dem Ende) | nur `…ready_for_query_in_der_Mitte,…` (Validate, Leser) | wie §7 |
| M3 überspringt Index drei vor dem Ende (direkt vor der Fehlerantwort) | nur `…ready_for_query_direkt_vor_der_Fehlerantwort,…` (Validate, Leser) | wie §7 |
| M4b überspringt Index zwei vor dem Ende allein (zwischen Fehlerantwort und Ende) | keiner, grün | wie §7: äquivalent, weil `pruefeEnde` („unmittelbar“) vorher anschlägt |
| M4a `pruefeEnde` lässt `ready_for_query` als Vorgänger durch | nur `…ready_for_query_zwischen_Fehlerantwort_und_Ende` (Validate, Leser) | rot **nur über den Text**: `ohneReadyForQuery` meldet weiter einen Fehler, der Text ist ein anderer (siehe V-1) |
| M4 beides zugleich (M4a und M4b) | nur `…ready_for_query_zwischen_Fehlerantwort_und_Ende` (Validate, Leser) | **rot durch die Zusage**: Ohne Fehler lädt die Folge; dies ist die Mutation, die nur diese Stelle bricht (V-1) |
| M8 `ohneReadyForQuery` aus `validateQuery` entfernt | die Fälle Anfang, Mitte, direkt davor (Validate, Leser); „zwischen“ bleibt grün (doppelt gedeckt durch `pruefeEnde`) | wie §7 |
| Stellenprüfung in `pruefeEnde` entfernt | sechs Fälle `TestValidateEndeFehler` (nicht letzte Nachricht, einfach/Extended, Ende mitten, Nachricht danach) und vier `TestUnmarshalEndeAbgelehnt` | wie §7 (Mutation 3) |
| `pruefeEnde` in `validateQuery` entfernt | die einfachen Fälle (Stelle, Felder, Schweregrad, unmittelbar) | wie §7 |
| `pruefeEnde(g.Server)` je Gruppe entfernt | die Extended-Fälle (Stelle, Feld, ERROR, ohne Fehlerantwort, Hinweis, `ready_for_query` nach dem Ende) | wie §7 |
| `!letzte && g.endet()` entfernt | nur die zwei Fälle „Ende in einer vorderen Gruppe“ (Validate), `TestUnmarshalEndeAbgelehnt/Ende_in_einer_vorderen_Gruppe` | wie §7 |
| `Interaction.endet` erkennt in `Responses` nur die letzte Nachricht (Mutation 1) | nur `TestSessionValidateEnd/falsch_platziertes_Ende_in_der_ersten_von_zwei,_einfach` und `…am_Anfang_der_ersten_von_zwei,_einfach` | wie §7; die Leser-Tests bleiben grün, weil `Interaction.Validate` zuerst anschlägt |
| `Group.endet` erkennt nur die letzte Server-Nachricht (Mutation 2) | nur `TestSessionValidateEnd/…letzten_Gruppe…` und `…vorderen_Gruppe…` (Extended) | wie §7 |
| `Session.ValidateEnd`: Prüfung `ii != letzte && i.endet()` abgeschaltet | zehn Fälle `TestSessionValidateEnd` und `TestUnmarshalEndeAbgelehnt/nicht_letzte_Interaktion,_einfach`, `…_extended`, `zwei_Enden_in_einer_Session` | wie §7 |
| Leser ruft `Session.ValidateEnd` nicht | nur `TestUnmarshalEndeAbgelehnt/nicht_letzte_Interaktion,_einfach`, `…_extended`, `zwei_Enden_in_einer_Session` | wie §7 |
| Leser: Version-1-Prüfung für `connection_closed` entfernt | nur `TestUnmarshalEndeAbgelehnt/Version_1,_einfach`, `…_extended` | wie §7 |
| Leser: Schlüssel neben `type` nicht geprüft | die fünf Fälle `Schlüssel_…` | wie §7 |
| Leser: Version 2 nicht angenommen | alle Version-2-Fälle von `TestUnmarshalEnde` und `TestUnmarshalEndeAbgelehnt` | wie §7 |

Der erste Versuch, `Session.ValidateEnd` mit `false` zu mutieren, brach den Bau (unbenutzte Variable); der Mutant `… && false` ist der in der Tabelle. Kein Mutant blieb grün, bei dem §7 „rot“ sagt.

## Nacharbeit (F-1, Stichprobe)

| Prüfung | Ergebnis |
|---|---|
| Entscheidung steht vor dem Code | Commit `aed8a49` (Architect: `SPEC-002` und Plan §6) liegt vor `91a5228` (Code und Tests). Reihenfolge nach `AGENTS.md` §3.12 eingehalten. |
| `ohneReadyForQuery` je Stelle | Anfang, Mitte, direkt vor der Fehlerantwort je ein Fall (Validate und Leser) mit eigener Mutation, siehe M1 bis M3; „zwischen“ ist über `pruefeEnde` gedeckt (M4). |
| Gegenfall ohne Ende | `TestValidateEndeGueltig` („query ohne Ende, ready_for_query mitten“) und `TestUnmarshalEnde/Version_2_ohne_Ende,_ready_for_query_mitten…` halten das akzeptierte Negativ der Spezifikation fest. |
| `endet` und `ValidateEnd` unverändert | Diff `91a5228` an `response.go`: 19 Zeilen, nur `ohneReadyForQuery` und Aufruf in `validateQuery`; `endet` und `ValidateEnd` unberührt, wie §7 sagt. |
| `pruefeEnde`-Stelle | siehe Tabelle oben (Stellenprüfung, `!letzte && g.endet()`, Mutation 1 und 2 färben nur die Fälle ihrer Gruppe). |
| F-4 (Test mit Zufallspfad) | Unabhängig: Der Test der Basis (`6d0d697`) mit `TMPDIR=/tmp/t4709x` ist rot (`Meldung nennt den Wert 4709: … /tmp/t4709x/…/rec.sqlite…`), der geänderte Test (`ohnePfad`) mit demselben `TMPDIR` grün. Die Ursache ist der Pfad, wie §7 sagt. Dass der Test den Wert weiter fängt, zeigen laut §7 die Mutanten (nicht nachgefahren). |

## Die Einschränkung zur Äquivalenz („zwischen Fehlerantwort und Ende“, doppelt gedeckt)

§7 sagt: Die Überspringung der Stelle in `ohneReadyForQuery` allein (M4b) bleibt grün und ist äquivalent; eine Mutation, die nur diese Stelle der neuen Prüfung bricht und den Fall rot färbt, gebe es nicht.

- **M4b ist äquivalent:** bestätigt. `pruefeEnde` läuft in `validateQuery` vor `ohneReadyForQuery` und meldet jedes Ende, dessen Vorgänger keine `error_response` FATAL/PANIC ist; ein `ready_for_query` an Index zwei vor dem Ende erreicht `ohneReadyForQuery` nie. Der Mutant ist grün und lässt kein beobachtbares Verhalten ändern.
- **Die Zusage ist nicht verletzt:** Die Zusage aus `SPEC-002` lautet „`ready_for_query` zwischen Fehlerantwort und Ende ist `PGR-E3003`“. Der Fall prüft sie, und es gibt eine Mutation, die sie bricht und nur den Fall rot färbt: **M4** (beide Verteidigungen entfernt, `pruefeEnde` lässt `ready_for_query` als Vorgänger zu und `ohneReadyForQuery` überspringt die Stelle). Gefahren: rot sind nur `TestValidateEndeFehler/query_ready_for_query_zwischen_Fehlerantwort_und_Ende` und `TestUnmarshalEndeAbgelehnt/ready_for_query_zwischen_Fehlerantwort_und_Ende`. Nach `AGENTS.md` §3.10 (Aufzählung und Ortsangabe zählen je Element: ein Fall und eine Mutation, die nur ihn bricht) ist das erfüllt; dass die Mutation zwei Stellen im Code anfassen muss, folgt aus der Doppelung, nicht aus einer Lücke im Test.
- **Was §7 zu eng sagt:** „Eine Mutation, die nur diese Stelle … bricht … gibt es nicht“ gilt für eine Mutation an *einer* Stelle des Codes. M4 (zwei Stellen) gibt es, M4a allein färbt den Fall rot nur über den Text der Meldung (V-1). Der Satz in §7 ist vor der Closure zu berichtigen oder zu streichen; die Aussage über die Zusage ändert er nicht.

## Abgrenzung, Plan und Größe

| Prüfung | Ergebnis |
|---|---|
| §1 Abgrenzung | `git diff --name-only 6d0d697..HEAD` außerhalb von `docs/plan/` und `docs/reviews/`: `internal/…` (Modell, Leser, Tests), `docs/user/abdeckung-*.md` (3), `spec/spezifikation.md` (Architect, `aed8a49`), `test/integration/sqlite_lesen_e2e_test.go` (2 Zeilen). Kein Schreiber: `Marshal`/`toDTO` unberührt, `grep -rn HasEnd internal` leer. Kein SQLite-Code (`sqlite*.go` nicht im Diff, `model.Version` bleibt 1, der SQLite-Leser lehnt Version 2 weiter ab: `TestE2ELesenAbgelehnt`-Fall „sqlite, Version 2“ steht). Kein Handbuch, kein README (`docs/user/benutzerhandbuch.md` und `README.md` nicht im Diff; die drei Dateien `abdeckung-*.md` sind Werkzeug-Ausgabe, §3 nennt sie). Kein `schema.yaml`. Eingehalten. |
| §3 (Plan) gegen `git diff --stat` | Jede Datei des Diffs steht in §3 (`response.go`, `recording.go`, `extended.go`, beide Testpakete, `yaml.go`, drei vorhandene Tests, `sqlite_lesen_e2e_test.go`, drei `abdeckung-*.md`); F-2 des Reviews (`abdeckung-e2e.md` statt `abdeckung-gesamt.md`/`-vollstaendig.md`) ist nachgezogen. Stimmt. |
| §6 Randformen | Jede entschieden vor dem Code (`690e353`, `aed8a49`); die zwei Lesarten aus F-1 stehen in §6 und `SPEC-002`. Der Diff entscheidet keine Randform, die §6 nicht nennt (Stichprobe: Version-Reihenfolge, Schlüssel neben `type`, Schweregrad `V`/`S`, Stelle des Endes, Version 2 ohne Ende, `offset_ms`/`empty_sessions`). |
| §7 gegen Diff | Tabelle, Läufe, Größe und *Grenzen* folgen dem Diff. Abweichung nur V-1 (Satz zu eng). |
| Größe (Metrik §4) | `git diff --numstat 6d0d697..HEAD -- internal docs/user`: 627 Einfügungen und 28 Löschungen, Summe **655** (selbst addiert aus den 13 Zeilen), gleich dem Wert des Implementers. Unter Meldeschwelle 800 und Marke 1000; in der Schätzung 550 bis 750. Zwei Liefer-Punkte, zwei Schichten (Kern, Recording-Adapter). Außerhalb der Metrik 2 Zeilen in `test/integration/…`. |
| §8 | Zählung und Messung stehen (Planner 2026-10-15; Implementer nach der Nacharbeit); keine Sendung an einen Nehmer im Diff. |
| `AGENTS.md` §3.9 | Der Plan folgt dem Diff in §1, §3, §6; Kopf (`Bezug`, `Berührte Spec-Stellen`) trägt `LH-FA-25`, `LH-FA-22.a`, `SPEC-001`, `SPEC-002`, `SPEC-041`, `ARC-001`, `ARC-008`. |
| `AGENTS.md` §3.11 | Abdeckungszeilen sagen nur zu, was Tests prüfen: die vier Zeilen zu `LH-FA-25` in `abdeckung-unit.md` entsprechen `TestUnmarshalEnde*`/`TestUnmarshalVersionen` (Fälle gelesen). Die Zeile der Fehlerfälle zählt „(am Anfang, in der Mitte, direkt vor und nach der Fehlerantwort)“ – alle vier Fälle gibt es. |
| `AGENTS.md` §3.13 (Adressen) | Die Risiko-Adressen nehmen an: `slice-v1-abschluss-verbindungsende-replay` (`next/`) nennt in §1 „Übernommen aus `slice-v1-abschluss-verbindungsende-format`“ (die Wiedergabe einer Aufzeichnung mit Ende) und in §4 den Start nach diesem Slice; `slice-v1-abschluss-verbindungsende-format-schreiber` (`next/`) nennt in §1 und DoD den Schreiber mit `HasEnd`, Roundtrip und die Voraussetzung dieses Slice. Beide liegen nicht in `done/`. |

## Die Zeile im Review-Report (Auftrag)

`git diff 3d14ef7 HEAD -- docs/reviews/2026-10-15-review-slice-v1-abschluss-verbindungsende-format.md`: genau eine Zeile, in der Eingangs-Kontext-Liste. Aus „[`ADR-0039`](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md) (Aufzeichnung kennt das Verbindungsende)“ wird „[`ADR-0039`](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md) (Aufzeichnung kennt das Verbindungsende)“. Bewertung:

- **Inhalt unverändert:** Befunde, Verdikt, Summary und Negativbefunde sind byte-gleich; es ist kein Finding und keine Zahl berührt.
- **Grund trägt:** Das Doku-Gate verlangt für jede ADR-Kennung einen auflösenden Link (`harness/conventions.md`, Zeile `modul-04-adrs.md`); §7 *Letzter Lauf* nennt den Befund (`id-unlinked`). Der Pfad löst auf (`docs/plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md` liegt vor, relativ von `docs/reviews/`). Andere Reports führen ADR-Kennungen ebenso als Link (`docs/reviews/2026-10-03-review-extended-query-planung.md`).
- **Zitier-Form des Reports:** Sie schreibt „Kennung, nicht Adresse“ für Slice, Sensor und Baseline vor, nicht für ADRs; die Ausnahme der ADR-Kennung durch das Gate ist keine Verletzung.
- **Vorbehalt (INFO, V-4):** Der Implementer hat ein Artefakt der Review-Rolle verändert. Das war hier rein formal und ist in §7 benannt, und ich halte es für zulässig; die sauberere Form wäre gewesen, die Änderung im Commit-Text der Nacharbeit (`91a5228`) zu nennen („Link [`ADR-0039`](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md) im Report“) – der Text nennt sie nicht, nur §7 des Plans. Das ist kein Grund für einen Rücksprung.

## Findings

### V-1 INFO

- `quelle`: `AGENTS.md` §3.9, §3.10
- `pfad`: `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-format.md` · §7, *Grüne Mutanten und Grenzen*, Eintrag „`ready_for_query` zwischen Fehlerantwort und Ende“
- `befund`: Der Eintrag sagt, eine Mutation, die nur diese Stelle der neuen Prüfung bricht und den Fall rot färbt, gebe es nicht. Es gibt sie: M4 (beide Verteidigungen entfernt) färbt nur den Fall rot, durch die Zusage und nicht durch den Text. Zudem ist M4a allein (nur `pruefeEnde` gelockert) rot, ohne dass die Zusage bricht (der Fehler kommt dann von `ohneReadyForQuery`), also nur über den Meldungstext des Falles; die Tabelle in §7 (Zeile „einfach mit Ende: kein `ready_for_query` zwischen Fehlerantwort und Ende“) nennt beide, der Text der *Grenzen* widerspricht ihr.
- `verifizierbar`: ja (Mutanten oben gefahren)
- `klasse`: Zusage in §7 weiter als die Prüfung (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`), kleine Form
- Vorschlag: Den Satz vor der Closure ersetzen durch „M4b allein ist äquivalent; die Zusage bricht erst M4 (beide Stellen), die nur den Fall zur Stelle rot färbt“.

### V-2 INFO

- `quelle`: `AGENTS.md` §3.12 (Meldung nicht zugesagt), §3.10
- `pfad`: `internal/hexagon/model/ende_test.go` · `TestValidateEndeFehler` und `yaml_ende_test.go` · `TestUnmarshalEndeAbgelehnt`
- `befund`: Beide Tests prüfen je Fall den Teilstring der Meldung („der erwartete Text belegt, dass die gemeinte Regel greift“), obwohl §6 und `SPEC-002` zusagen, dass nur der Fehler zugesagt ist und nicht, welche Prüfung meldet. Folge: Mutationen, die nur die *Reihenfolge* der Prüfungen ändern (M4a), färben Fälle rot, ohne die Zusage zu brechen; die Tests binden die Reihenfolge von `pruefeEnde`, `ohneReadyForQuery` und `validateServer` mehr, als der Vertrag es tut. Das ist eine bewusste Stärke (jeder Fall belegt seine Regel) und kein Fehler, aber eine Kopplung, die ein Refactor der Reihenfolge bezahlt. Der Kommentar am Test nennt sie.
- `verifizierbar`: ja
- `klasse`: Test enger als die Zusage (kein Befund; Grenze zur Kenntnis)
- Vorschlag: keine Aktion; bei einem Refactor der Prüfreihenfolge die Texte nachziehen.

### V-3 INFO

- `quelle`: `AGENTS.md` §3.11
- `pfad`: `docs/user/abdeckung-gesamt.md` und `abdeckung-vollstaendig.md` · Zeile `LH-FA-25`; `internal/hexagon/model/ende_test.go`
- `befund`: Zweierlei zur Abdeckung. (a) `LH-FA-25` steht als „vollständig“ (Unit), obwohl Aufzeichnen, Wiedergabe, Einspielen und SQLite in sieben weiteren Slices folgen (Review F-3, beim Planner; unverändert). (b) Die Tests des Modells (`TestValidateEndeGueltig`, `TestValidateEndeFehler`, `TestSessionValidateEnd`, `TestVersionMitEnde`) tragen keine Abdeckungs-Deklaration; die vier Zeilen zu `LH-FA-25` stammen allein aus den Lesertests. `make abdeckung-check` verlangt für Unit-Tests keine, die Tabelle ist Ausgabe der Deklarationen. Kein falscher Satz, aber „Unit“ belegt `LH-FA-25` heute nur über den Leser.
- `verifizierbar`: ja (`grep -n "Abdeckung:"` in `ende_test.go` leer)
- `klasse`: Abdeckungs-Stand „vollständig“ vor Lieferung der Anforderung
- Vorschlag: (a) Planner entscheidet die Lesart; (b) keine Aktion, oder eine Deklaration an `TestValidateEndeFehler` (Spec-Stelle `SPEC-002`).

### V-4 INFO

- `quelle`: Modul 8 (Rollen), `AGENTS.md` §6
- `pfad`: `docs/reviews/2026-10-15-review-slice-v1-abschluss-verbindungsende-format.md` · Eingangs-Kontext
- `befund`: Der Implementer hat eine Zeile im Report der Review-Rolle geändert (Link auf [`ADR-0039`](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md)). Rein formal, siehe oben; der Commit-Text der Nacharbeit nennt es nicht, §7 des Plans schon. Zulässig, kein Rücksprung.
- `verifizierbar`: ja (`git diff 3d14ef7 HEAD -- <Report>`)
- `klasse`: Rollengrenze, formal
- Vorschlag: künftig die Änderung eines Review-Reports im Commit-Text nennen.

## Risiko-Ausgänge (Vorschläge für die Closure)

| Risiko (§6) | Vorschlag | Begründung |
|---|---|---|
| Zwischenstand: `replay` lädt eine Aufzeichnung mit Ende, der PGWire-Adapter kennt `connection_closed` noch nicht | *eingetreten* → Folge-Slice `slice-v1-abschluss-verbindungsende-replay` | Der Zwischenstand besteht (nur das Laden ist geliefert, `Interaction.Validate` nimmt das Ende an jeder Interaktion an, kein Adapter gibt es wieder). Die Adresse nimmt an: §1 „Übernommen aus … format“ (Wiedergabe) und §4 „liegt in `done/`“ als Start. |
| Der ungeteilte Plan lag über einer Review-Sitzung (850 bis 1250) | *eingetreten* → Folge-Slice `slice-v1-abschluss-verbindungsende-format-schreiber` | Steht schon so in §6; die Adresse nimmt den Schreiber mit Roundtrip in §1 und DoD an. |
| Dieser Teil überschreitet Meldeschwelle oder Marke (Schätzung 550 bis 750) | *entfallen*, mit Begründung | Gemessen 655 in der Metrik von §4: unter Meldeschwelle 800 und Marke 1000, in der Schätzung. Das ist die Metrik; die Beobachtung `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` trat hier nicht auf (kein Auftreten; die Schätzung lag in der Spanne). Hinweis für die Closure: die erste Fassung maß 590, die Nacharbeit aus dem Review 65 Zeilen mehr; die Randform-Rückgabe an den Architect hat den Zuwachs getragen. |
| Ort der Formregel einer Session | *entfallen*, mit Begründung | Entschieden vor dem Code (`Session.ValidateEnd`), geliefert und nachgewiesen (`TestSessionValidateEnd`, Leser ruft sie). Das steht in §6 schon so. |

Weitere Punkte für die Closure:

- Review F-3 (Abdeckung „vollständig“) → Planner: Ausgang *weiter offen* (Register) oder Lesart festlegen, bevor die Anforderung geliefert ist.
- Lerneintrag, Kandidaten (Urteil des Planners):
  1. `BEO-REPO/spec-randform-erst-im-review-entschieden` ist in diesem Slice trotz `AGENTS.md` §3.12 wieder aufgetreten (Review F-1: zwei Lesarten im Code entschieden, §6 nannte sie nicht): löst den Retirement-Check der Regel aus (kein Lockern). Klasse: Stelle einer Nachricht in der Nachrichtenfolge (`ready_for_query` vor dem Ende, `connection_closed` mitten in der Folge); die Nacharbeit zählte die Ausprägungen der Klasse und deckt vier plus sechs davon.
  2. Neue Beobachtung: ein Test, der eine Meldung mit Teilstring gegen einen Wert prüft, den ein Zufallspfad enthalten kann (`t.TempDir()`); die Wirkung ist ein seltener Fehlalarm im Gate (Review F-4, von mir mit `TMPDIR` reproduziert). Vorschlag: Pfade aus Meldungen streichen, bevor gegen Werte geprüft wird – Form des Lerneintrags: geschärfte Regel oder Test-Hilfe (`ohnePfad`).
  3. §3 des Plans weicht vom Diff ab (Review F-2): die bekannte Klasse `BEO-REPO/plan-folgt-korrektur-nicht`; ein Auftreten, kein neuer Eintrag.
- Die Metrik der Sendung nach dem Schnitt (§3.13): keine Sendung ist eingegangen; die Zählung in §8 hat der Planner vor der Arbeit gemacht.

## Gates und Zustand

- Eigene Läufe: unveränderte Kopie von `HEAD`, `go test ./internal/hexagon/... ./internal/adapters/driven/recording/` grün; 18 Mutanten einzeln, alle wie in §7 behauptet rot (M4b grün, äquivalent); der F-4-Test mit und ohne Pfadfilter unter `TMPDIR=/tmp/t4709x`.
- Voller `make gates` nicht gefahren; der Stop-Hook gab `approve` (siehe DoD-Punkt 4). `make lint` nicht gefahren; §7 trägt die Behebung der Lint-Befunde (`b2b26fb`), der letzte volle Lauf war grün.
- Aufräumen: Kopie und Skripte im Scratchpad, keine Container und kein Netz angelegt, eigene Wartekommandos per PID beendet, `git status` zeigt nur diesen Report (nicht committiert).
