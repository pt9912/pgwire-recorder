# Verifikation slice-v1-abschluss-verbindungsende-sqlite-schreiber

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-15. **Spanne:** `b5933e8` (Basis der Marke) bis `e21eea8` (HEAD): Architect `1ce3d3a`, Umsetzung `e63b184`, Review-Report `27b69aa`, Nacharbeit `e21eea8`. **Gegenstand:** Slice `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Plan unter `docs/plan/planning/in-progress/`, Schritt 36b, SQLite-Schreiber und Roundtrip) gegen DoD (§2), Plan §1 (Abgrenzung), §3, §4, §6 (Randformen, Zählung der Mutationstabelle, Risiken), §7 (Belege des Implementers), §8, `LH-FA-22.a` (*Ende und Version beim Schreiben von `sqlite`*), `SPEC-001` (*Schreiber in `sqlite`*), `SPEC-043`, `LH-FA-25` mit `LH-FA-25.a`, [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md) Punkt 4, [ADR-0013](../plan/adr/0013-sqlite-recording-backend.md) und `AGENTS.md` §3.7 bis §3.14. Review: `docs/reviews/2026-10-15-review-slice-v1-abschluss-verbindungsende-sqlite-schreiber.md` (F-1 und F-2 LOW, F-3 und F-4 INFO, nachgearbeitet in `e21eea8`); gegen ihn ist nicht verifiziert, F-1 bis F-4 sind unabhängig nachgeprüft.

Eine Verifier-Vorlage führt die vendored Baseline nicht; Aufbau wie der Report von `slice-v1-abschluss-verbindungsende-sqlite`. Belege sind in diesem Lauf selbst erzeugt: Quelle von `HEAD` (`git archive`: `go.mod`, `go.sum`, `internal`) in eine Kopie im Scratchpad, Läufe per Bind-Mount mit dem Image `pgr-36b:deps` (`--network=none`, `go test -c` einmal je Mutant im Paket `./internal/adapters/driven/recording/`, Lauf von `TestSQLite(Version|Ende|SchreibenNimmt|Roundtrip)`), jeder Mutant einzeln aus einer unveränderten Kopie angewandt, danach verworfen. Die Mutanten sind **selbst gebaut** nach den Zeilen der Tabelle in §7, nicht die des Implementers. Ein Lauf zur Zeit, kein `pkill`, kein `make gates`. Der Arbeitsbaum des Repos ist unverändert (`git status` zeigt nur diesen Report).

**Verdikt: PASS mit Befunden.** Alle Liefer-Punkte der DoD sind bestätigt. 42 selbst nachgebaute Einzelmutanten und 36 Läufe der nach `user` schlüsselnden Mutanten (28 der Stellen, acht des Roundtrips) färben genau die in §7 genannten Tests; die Zeile `davor` ist nach F-3 richtig gefasst. Die Abgrenzung ist eingehalten, der Plan folgt dem Diff bis auf Kleinigkeiten, die Größe ist gemessen 595. **0 HIGH, 0 MEDIUM, 1 LOW, 5 INFO.** Der LOW (V-1): Die Zusage „genau eine Zeile“ ist nur für die Seite „null Zeilen“ geprüft; der Mutant `n != 1` zu `n == 0` bleibt grün, die Seite „mehr als eine Zeile“ ist erreichbar (Sonde unten). Offen sind allein die Closure-Pflichten; das ist der erwartete Stand vor der Closure.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren oder gelesen) |
|---|---|---|
| 1: Schreiber legt `meta.version` = 1 an, setzt sie in der Transaktion der ersten Session mit Ende auf 2; spätere Sessions ändern sie nicht, eine Datei ohne Ende bleibt bei 1; das Ende steht als `connection_closed` mit `fields` `{}` in der Reihenfolge der Interaktion; scheitert die Transaktion, bleibt Version 1 und nur die zuvor ergänzten Sessions, der nächste Schreibvorgang ergänzt (Test mit eingespeistem Fehler) | bestätigt, mit V-1 | Code gelesen: `legeSchemaAn` schreibt `model.Version`, `zeilenDerSession` setzt `mitEnde` über `Recording{Sessions: {s}}.HasEnd()` (Modell unberührt), `ergaenze` bildet `setzeVersion := z.mitEnde && !s.endeGesetzt`, `ergaenzeSession` setzt die Version vor `tx.Commit()` in derselben Transaktion, `setzeVersionMitEnde` prüft `RowsAffected() == 1` (Fehler wird `PGR-E3001` über `schreibe`), das Flag `endeGesetzt` steht nach `ergaenzeSession` und vor dem Schließen der Datenbank. Tests: `TestSQLiteVersionNachInhalt` (17 Fälle), `…WerSetzt` (6), `…ZeileUnpassend` (3), `…FehlerInDerTransaktion` (2: Auslöser in der Datei über `rohSchreibe`, doppelte `sequence`), `…SchliessenScheitert`, `…AbbruchUnterbrichtNicht`, `TestSQLiteEndeMitFalscherNummer`, `TestSQLiteEndeDarstellung`. Unveränderte Kopie: PASS. Die Darstellung des Endes ist keine neue Abbildung (`serverZeile` unverändert), der Test belegt Stelle und Form (`typeof`, `length`). |
| 2: Roundtrip in beiden Formaten (einfach, `Sync`, `Flush`, `offset_ms`); ein späterer Schreibvorgang an Version 2 prüft `TestSQLiteVersionWerSetzt`; verletzte Form nimmt der Schreiber an, der Leser lehnt mit `PGR-E3003` ab | bestätigt, mit V-2 | `TestSQLiteRoundtripMitEnde` (fünf Formen + „Version 2 ohne Ende“): YAML und SQLite lesen gleich (`reflect.DeepEqual` auf der geladenen Aufzeichnung, auch `Version`), gleich der geschriebenen mit Version 2, Laden–Schreiben–Laden gleich. `TestSQLiteSchreibenNimmtFormverletzungAn`: 14 Stellen, je `err == nil`, `meta.version` 2, Laden nach Muster. Roundtrip-Mutanten je Form gefahren (unten). |
| 3: je Zusage eine Mutation, vom Implementer rot gesehen; Tabelle in §7 | bestätigt, 42 Einzelmutanten und 36 schlüsselnde Läufe nachgefahren | siehe nächster Abschnitt; V-1 (nicht geprüfte Hälfte einer Zusage), V-6 (Fall „Zeile fehlt“ ohne eigene Mutation, benannt). |
| 4: `make gates` grün | durch den Hook getragen, nicht selbst gefahren | §7 *Lauf nach der Nacharbeit*: `make gates` vollständig Exit 0. `echo '{}' \| bash .claude/hooks/stop-require-gates.sh` gibt `{"decision":"approve"}` (Exit 0); `.harness/state/gates-passed.dateiliste` 20:09:01, Commit `e21eea8` 20:09:24, Arbeitsbaum sauber; `make abdeckung-check` und `make kopf-check` selbst: Exit 0. Ob der volle Lauf aus dem richtigen Grund grün war, belegt die Datei nicht; ich habe `make gates` nicht wiederholt (Auftrag). |
| 5: Review durchgeführt, Report liegt vor | bestätigt | Report unter `docs/reviews/` (0 HIGH, 0 MEDIUM, 2 LOW, 2 INFO); F-1 bis F-4 in `e21eea8` nachgearbeitet, unten nachgeprüft. Kein Self-Review. |
| 6 bis 9: Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen | offen (erwartet) | §7 trägt Belege des Implementers, noch keine Closure-Notiz; §2 ohne Häkchen; beide Risiken in §6 tragen „offen bis Closure“. Handlung des Planners; Vorschläge unten. |

## Mutationen selbst nachgefahren

„nur“ heißt: keine anderen Tests rot. Quelle der Mutanten: die Zeile in §7. Namen wie in §7, Nachbau von mir. Kurzform der Tests: `NachInhalt` = `TestSQLiteVersionNachInhalt`, `WerSetzt`, `ZeileUnpassend`, `Fehler` = `…FehlerInDerTransaktion`, `Schliessen` = `…SchliessenScheitert`, `Darstellung` = `TestSQLiteEndeDarstellung`, `Stellen` = `TestSQLiteSchreibenNimmtFormverletzungAn`, `Roundtrip` = `TestSQLiteRoundtripMitEnde`.

| Mutant | Rot | Urteil |
|---|---|---|
| `keine` (Version nie gesetzt) | `NachInhalt` alle Fälle „mit Ende“, `WerSetzt`, `ZeileUnpassend`, `Fehler`, `Schliessen`, `Abbruch`, `Stellen` alle 14 | wie §7 |
| `immer` (erste Transaktion setzt auch ohne Ende) | `NachInhalt` „Feld N ohne Ende“ (fünf) und „nur empty_sessions“; `WerSetzt`/„spätere Session setzt, frühere nicht“ und „Session ohne Ende ändert meta nicht“; `Darstellung`/„ohne Ende“; `Roundtrip`/„Version 2 ohne Ende“; dazu `ZeileUnpassend`, `Fehler`/Auslöser | wie §7 (mehr rot als genannt) |
| `feld1` (Ende wirkt nicht bei Modellfeld 1) | nur `NachInhalt`/„Feld 1 mit Ende“ | wie §7 |
| `feld1ohne` (Anlegen 2 bei Feld 1 ohne Ende) | nur `NachInhalt`/„Feld 1 ohne Ende“ | wie §7 (F-2 getrennt) |
| `feldMit` K = 2, 0, 3, -1 (Anlegen schreibt `rec.Version`, wenn K und Ende) | je K nur `NachInhalt`/„Feld K mit Ende“; bei K = 2 zusätzlich alle fünf Formen von `Roundtrip` | wie §7 |
| `feldOhne` K = 2, 0, 3, -1 | je K nur `NachInhalt`/„Feld K ohne Ende“; bei K = 2 zusätzlich `Roundtrip`/„Version 2 ohne Ende“ | wie §7 |
| `leer2` · `empty2` | nur `NachInhalt`/„ohne Sessions“ · `NachInhalt`/„nur empty_sessions“ und `Roundtrip`/„empty_sessions“ | wie §7 |
| `nurErste` | `NachInhalt`/„Ende nur in einer späteren Session“, `Abbruch` | wie §7 |
| `nurLetzte` | `NachInhalt` (zwei Fälle), `WerSetzt`/„zweite Session mit Ende im selben Schreibvorgang“, `Stellen`/„nicht letzte Session“, `Roundtrip`/„empty_sessions“ | wie §7 |
| `hasEndeResponses` · `hasEndeGroups` | Gruppen-Fälle in `NachInhalt`, `Stellen`, `Roundtrip` · Fälle mit einfachem Ende überall | wie §7 |
| `anlegen2` | `Fehler`/„doppelte sequence im ersten Schreibvorgang“ und alle Fälle mit Ende | wie §7 |
| `jede` · `nachV2` | `WerSetzt` (drei bzw. vier Fälle), `Schliessen` (beide); `nachV2` zusätzlich `NachInhalt`, `Stellen`/„nicht letzte Session“, `Roundtrip`/„empty_sessions“ | wie §7 |
| `zeileFehltOk` (null Zeilen kein Fehler, wenn keine Zeile `version`) | nur `ZeileUnpassend`/„Zeile fehlt“ und „Zeile fehlt, andere Zeile mit Wert 1“ | wie §7 |
| `wert2Ok` (null Zeilen kein Fehler, wenn eine Zeile `version` steht) | nur `ZeileUnpassend`/„Wert 2“ | wie §7 |
| `zeilenzahl` (`n < 0`) | `ZeileUnpassend` alle drei | wie §7 |
| `ohneWert` · `ohneKey` | nur `ZeileUnpassend`/„Wert 2“ · nur `ZeileUnpassend`/„Zeile fehlt, andere Zeile mit Wert 1“ | wie §7 |
| **`nNull`** (`n != 1` zu `n == 0`: „mehr als eine Zeile“ kein Fehler) | **keiner, PASS** | **grün, nicht in §7; V-1** |
| `userVersion` | `NachInhalt` alle zehn Fälle „mit Ende“ | wie §7 |
| `commitTrotz` · `ausserhalb` | `ZeileUnpassend` (drei), `Fehler`/„Auslöser beim Setzen“ | wie §7 |
| `davor` (eigene Transaktion vor dem Einfügen, Zeilenzahl geprüft) | nur `Fehler`/„doppelte sequence im ersten Schreibvorgang“; `ZeileUnpassend` und der Fall mit Auslöser bleiben grün | wie §7 nach F-3 (die Zeile nennt nur diesen Fall) |
| `flagVorCommit` · `nachSchliessen` | beide Fälle von `Fehler` · `Schliessen` und `WerSetzt`/„zweite Session mit Ende im selben Schreibvorgang“ | wie §7 |
| `abbruch` · `nummerNach` | nur `Abbruch` · nur `EndeMitFalscherNummer` | wie §7 |
| `endePos` · `endeKind` · `endeFelder` (`{ }`) · `endeVerwerfen` · `endeGruppeNull` · `endeGruppeEins` · `endeEinfachGruppe` · `endeSeq` · `endeRichtung` | `Darstellung` (jeweils die Fälle der Zeile in §7), `Roundtrip`, `Stellen`; `endeVerwerfen` nur `Darstellung`/„Feld neben type steht in fields“ und `Stellen`/„Feld neben type“; `endeGruppeEins` nur `Darstellung`/„Flush-Gruppe“ und `Roundtrip`/„Flush-Gruppe“ | wie §7 |
| Stellen: `ablehnen` und `ohneVersion`, je 14 Läufe (Schlüssel: Startup-Parameter `user`) | je Lauf **genau** der Fall dieser Stelle in `Stellen`, kein anderer (28 von 28) | wie §7 |
| Roundtrip: `ohneEndeZeile` je Form (fünf Läufe) · `ohneOffset` · `ohneEmptyMeta` | je Lauf genau `Roundtrip`/<Form> · nur „offset_ms“ (für „empty_sessions“ grün, dort kein Offset) · nur „empty_sessions“ | wie §7 |

Nicht nachgefahren: `metaExtra`, `offsetEnde`, `endeEinfachGruppe` in der Variante an der Gruppe einer Extended-Interaktion (liegen im Bestand der Darstellung und sind aus den Tabellen der Fälle abgeleitet); die Mutanten der Aufruf-Voraussetzung (keine, §6). Ohne Mutation läuft die Kopie grün (PASS).

## Bewertung der Fragen

**(1) Die Zusage „genau eine Zeile, sonst `PGR-E3001` mit Rollback“ und die nicht trennbaren Fälle.**

- *Null Zeilen, beide Ursachen:* geprüft. Fehlt die Zeile oder steht sie auf 2, ist das `PGR-E3001`, keine Zeile der Session, `meta` unverändert (drei Fälle, Mutanten `zeilenzahl`, `zeileFehltOk`, `wert2Ok`, `ohneWert`, `ohneKey` je rot wie in §7). Der Rollback ist geprüft (`commitTrotz` rot in allen drei Fällen und im Fall mit Auslöser). Die Zusage hält auf dieser Seite.
- *Mehr als eine Zeile:* **nicht geprüft** (V-1). `if n != 1` ist mit `if n == 0` austauschbar, ohne dass ein Test rot wird. Im Schema des Schreibers ist `key` Primärschlüssel (`tools/schema/schema.yaml`), also ist die Seite dort unerreichbar; in einer fremden Datei ohne Primärschlüssel ist sie erreichbar. Sonde (nur in der Kopie, kein Repo-Eintrag): `meta` neu ohne Primärschlüssel mit zwei Zeilen `version`/`1`; das Original lehnt mit `PGR-E3001` ab („2 Zeilen mit dem Wert 1“), der Mutant nimmt die Session an und lässt zwei Zeilen auf 2 setzen. Der Kommentar von `setzeVersionMitEnde` („jede andere Zahl geänderter Zeilen ist ein Fehler“) und die Fehlermeldung nennen die Seite, `SPEC-001` sagt „genau die eine Zeile“. §6 und §7 nennen sie nicht (dort „Zeile fehlt, Wert ist 2“); die Zusage im Code ist weiter als der Test.
- *„Zeile fehlt“ und „Zeile fehlt, andere Zeile mit Wert 1“:* **nicht getrennt belegbar, und die Tabelle sagt es richtig** (V-6). Der Code unterscheidet beide nicht (kein Zweig auf die Existenz der Zeile). `zeileFehltOk` färbt beide, `ohneKey` nur den zweiten; ein Mutant, der nur „Zeile fehlt“ färbt, ist mit den natürlichen Eingriffen nicht zu bauen (`meta` des ersten Falls ist eine Teilmenge des zweiten). Der erste Fall ist durch den zweiten dominiert; er trägt als Eingang („Zeile fehlt“ ohne Nebenbedingung) die Lesart der Spezifikation und stört nicht. Nach `AGENTS.md` §3.10 („Fall und Mutation, die nur ihn bricht“) ist das eine benannte Grenze, kein Mangel.

**(2) Enge der Zusage im Roundtrip-Doc nach F-1.** F-1 ist behoben: Weder das Doc von `TestSQLiteRoundtripMitEnde` noch die drei Zeilen in `abdeckung-unit.md` (`LH-FA-22`, `LH-FA-25`, `LH-QA-06`) sagen noch zu, dass ein späterer Schreibvorgang Sessions ergänzt; `grep` auf „späterer Schreibvorgang“ im Roundtrip-Text ohne Treffer, DoD und §1 verweisen auf `TestSQLiteVersionWerSetzt` (dort Mutanten `jede`, `nachV2` rot, oben gefahren). Eine Restweite bleibt (V-2): „liest gleich (auch in `version`) und gleich der geschriebenen“. Die geschriebene Aufzeichnung trägt `version` 5 (`feldNeutral`), die gelesene 2; der Test vergleicht mit `want.Version = model.VersionMitEnde`. „Gleich der geschriebenen“ gilt ohne diese Ausnahme nicht, `SPEC-001` nennt sie („ausgenommen das Feld `version`“). Kleiner Befund, kein Mutant bleibt grün.

**(3) Abgrenzung.** Eingehalten, siehe Tabelle unten.

**(6) MUT-Technik des Implementers (Umgebungsvariable in einer Kopie).** Keine Spur im Repo: `grep -rn 'Getenv' internal/adapters/driven/recording` ohne Treffer, kein `MUT` im Code, in `Makefile`, `harness/`, `tools/`, `test/`; Treffer für `MUT` nur im Prosa-Text von §7 dieses Plans. Der Diff enthält keine Datei außer den genannten. Die Tabelle ist aus den Zeilen nachbaubar: alle Nachbauten ergaben dasselbe Bild wie die Zeile (`davor` in der nach F-3 berichtigten Fassung; bei `immer` und `keine` färbt der Mutant mehr Tests, als die Zeile nennt, ohne Widerspruch). Die Mutanten selbst liegen nicht im Repo; wer sie reproduzieren will, baut sie aus der Beschreibung (siehe Lerneintrag-Kandidat 3).

## Nacharbeit (F-1 bis F-4, Stichprobe)

| Prüfung | Ergebnis |
|---|---|
| F-1: Roundtrip-Doc und drei Abdeckungszeilen enger | bestätigt, Rest V-2 |
| F-2: Paare getrennt (`feldMit`, `feldOhne`, `feld1ohne`; `zeileFehltOk`, `wert2Ok`) | bestätigt: jede Zeile färbt nur ihre Fälle (14 Läufe nachgefahren, Tabelle oben); „Zeile fehlt“ ohne eigene Mutation ist ehrlich benannt (V-6) |
| F-3: `davor` färbt nur „doppelte sequence“ | bestätigt (gefahren), `ZeileUnpassend` und Auslöser-Fall grün |
| F-4: Stellen-Tabelle nicht geteilt, Begründung in §7 | Grenze, keine Adresse; die Begründung (Prüfung der Version gegen Kopf-Text im YAML, `meta.version` im SQLite, Schlüssel `user`) ist in den beiden Tests nachzulesen. Siehe Lerneintrag-Kandidat 4. |

## Abgrenzung, Plan und Größe

| Prüfung | Ergebnis |
|---|---|
| §1 Abgrenzung | `git diff --name-status b5933e8..HEAD`: Code nur `sqlite.go`, `sqlite_zeilen.go`, `sqlite_schreiben_ende_test.go` (neu), `abdeckung-unit.md`; dazu `spec/spezifikation.md` (Architect), Plan, Roadmap, Review-Report, `-record`. **Leser (36a) unverändert:** keine `sqlite_lesen*.go`, `yaml.go` oder `sqlite_json*.go` im Diff. **Kein Handbuch, kein README:** `docs/user/benutzerhandbuch.md` und `README.md` nicht im Diff (`abdeckung-unit.md` ist Werkzeug-Ausgabe, §3 nennt sie). **Kein Modell-Eingriff:** nichts unter `internal/hexagon/` im Diff (`HasEnd` auf einer Aufzeichnung mit einer Session, keine neue Methode); Port unberührt. **Schema:** `tools/schema/schema.yaml` und `export_test.go` unberührt, `make schema-validate` nicht betroffen. Eingehalten. |
| Bindung an `-record` (`AGENTS.md` §3.13) | Commit `1ce3d3a` trägt die Bindung im Nehmer: §1 („für `sqlite` bestätigt … `SPEC-001` *Schreiber in `sqlite`*“), §6 *Form beim Schreiben* mit beiden Kennungen, §8 mit Zählung („Liefer-Punkte weiter 3, Schichten weiter zwei“, Kennung des Gebers). Die Adresse nimmt an: `-record` liegt in `next/`, ihr Use Case ruft `Session.ValidateEnd` einmal für beide Formate. |
| Bindung an `-upstream` | Die Übernahme des E2E-Belegs („`record --format sqlite` zeichnet ein echtes Ende auf, die Datei trägt `meta.version` 2“) steht in §1 des Nehmers (Absatz „Übernommen aus `slice-v1-abschluss-verbindungsende-sqlite-schreiber`“), der Start nennt den Geber, §8 trägt die Nachzählung des Planners; der Nehmer liegt in `next/`. In diesem Diff keine Änderung nötig (Eintrag vom Planner aus dem Schnitt). |
| §3 (Plan) gegen `git diff --stat` | Jede Code-Datei des Diffs steht in §3; `sqlite_json.go` und `export_test.go` stehen als „unberührt“ und sind es. `spec/spezifikation.md` und `-record` stehen in keiner Zeile von §3, sind aber Gegenstand des Architect-Commits und in §6/§8 benannt; kein Mismatch. |
| §6 Randformen | Jede vor dem Code entschieden (`1ce3d3a` vor `e63b184`) in `LH-FA-22.a` (*Ende und Version beim Schreiben von `sqlite`*), `SPEC-001` (*Schreiber in `sqlite`*), `SPEC-043` und im Drift-Log (eine Zeile 2026-10-15). Der Diff entscheidet keine Randform, die §6 nicht nennt (Stichprobe: Modellfeld 0/1/2/3/-1, Reihenfolge `PGR-E1000` vor `PGR-E3001`, Flag vor dem Schließen, `user_version`, Abbruch, Felder neben `type`, 14 Stellen). Die Seite „mehr als eine Zeile“ ist in §6 nicht genannt und im Code entschieden (V-1). |
| §7 gegen Diff | Tabellen, Läufe, Größe, Nacharbeit folgen dem Diff. Abweichungen: V-1 (Seite der Zeilenzahl), V-5 (596 in §6). |
| Größe (Metrik §4: `git diff --numstat b5933e8..HEAD -- internal docs/user`) | Selbst addiert: `sqlite.go` 14/6, `sqlite_zeilen.go` 32/3, `sqlite_schreiben_ende_test.go` 519/0, `abdeckung-unit.md` 21/0: Einfügungen 586 plus Löschungen 9 = **595**, gleich dem Wert nach der Nacharbeit in §7 und §8. Produktcode 55, Tests 519, Tabelle 21. **60 % der Marke 1000**, unter der Meldeschwelle 800; in der Spanne des Architect (480 bis 760, Mitte 620; die Messung liegt 4 % unter der Mitte). Faktor Tests zu Produktcode 519 / 55 ≈ 9,4, in der Spanne der Datenpunkte 8 und 11. Liefer-Punkte 2 (Schreiber mit Version in der Transaktion; Roundtrip), Schichten 2 (Recording-Adapter, `docs/user`); der Schnitt „Version in der Transaktion \| Roundtrip“ war nicht nötig. |
| §8 | Zählung (Planner, Architect, Implementer, Nacharbeit) steht und stimmt mit der Messung (596 beim Umsetzen, 595 nach der Nacharbeit). |
| `AGENTS.md` §3.9 | Plan folgt dem Diff in §1, §3, §6, §7. Kopf: `LH-QA-06` fehlt (V-4). |
| `AGENTS.md` §3.11 | Docs der Tests und Abdeckungszeilen entsprechen den Fall-Listen (gelesen); Ausnahmen V-2 und V-3. |
| `AGENTS.md` §3.7 | Kommentare im Produktcode: Zusage und Kopplung im Indikativ, keine Chronik, kein Konjunktiv; Ausnahme V-1 (Zusage weiter als der Test). |
| Nicht im Diff | `.harness/`, `harness/`, `Makefile`, `.golangci.yml`: unberührt; kein `//nolint`. |

## Risiken aus §6 und Vorschlag für den Ausgang

| Risiko | Vorschlag | Begründung |
|---|---|---|
| Ein Fehler in der Transaktion ist ohne Einspeisung nicht erreichbar | **entfallen** | Belegt: beide Mittel tragen, ohne Änderung am Produktcode oder an `export_test.go`. Gefahren: `ausserhalb` (Version nach dem Commit in einer eigenen Transaktion) färbt den Fall mit dem Auslöser, `davor` (vor dem Einfügen in einer eigenen Transaktion) färbt allein den Fall mit der doppelten `sequence`; `flagVorCommit` färbt beide. Der Nehmer ist nicht berührt. Die Seite „mehr als eine Zeile“ (V-1) ist ein anderer Gegenstand, kein Rest dieses Risikos. |
| Der Diff liegt über der Marke (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`) | **entfallen** | Gemessen 595 in der Metrik von §4 (60 % der Marke 1000, unter der Meldeschwelle 800, in der Spanne 480 bis 760 des Architect); das Review lief in einer Sitzung ohne Größen-Finding. Der Faktor Tests zu Produktcode (9,4) gehört als Datenpunkt ins Register zu `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (kein Zählerzuwachs, die Schätzung hat getragen). |

Bei V-1 zu entscheiden (Planner): Entweder ein Fall „zwei Zeilen `version` mit Wert 1“ (Datei ohne Primärschlüssel, `rohSchreibe`) samt Mutation `n == 0` im Slice, oder der Kommentar und die Spezifikation schreiben „keine Zeile“ statt „genau eine“ / „jede andere Zahl“; falls die Seite offen bleibt, ist es ein **weiter offenes** Risiko ins Register (Datei ohne Primärschlüssel ist Fremddatei, kein Weg des Schreibers). Welche Variante gilt, steht nicht in dieser Rolle.

## Lerneintrag-Kandidaten (Entscheidung beim Planner)

1. **Zusage im Doc oder Kommentar weiter als die Prüfung** (Review F-1; hier erneut V-1, V-2, V-3): vier Instanzen in einem Slice, in jeder die Zusage in Prosa weiter als der Test. Form: **Zähler** an `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Beleg: dieser Slice, drei Orte: Roundtrip-Doc, Kommentar „jede andere Zahl“, „dieselben Zeilen wie bisher“). Wenn der Planner eine geschärfte Regel will: Zusagen mit Zahlwort („genau eine“, „jede andere Zahl“, „alle Zeilen“) bekommen im Beleg je Seite einen Fall, in `AGENTS.md` §3.10 neben *Aufzählung und Ortsangabe zählen je Element*.
2. **Mutationstabelle zählt Paare statt Elemente** (Review F-2): Die Zeile fasste „mit Ende“ und „ohne Ende“ zusammen; erst das Review trennte sie. Wiederkehr von `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` in der Form „Aufzählung zählt als eins“ (`AGENTS.md` §3.10 *Aufzählung* gilt seit `slice-harness-gate-kette-stop-hook-liste`): Zähler, mit diesem Slice als Beleg, und die Rückfrage, ob §6 *Zählung der Mutationstabelle* als Anker des Planners (steht hier bereits, vom Architect) die Klasse nicht verhindert hat, weil die Tabelle in §7 der Zählung nicht gegen gelesen wurde.
3. **Mutantenzeile nicht reproduzierbar** (Review F-3, hier gegengeprüft): Die Zeile `davor` nannte rote Tests, die der Nachbau nicht färbt. Die übrigen Nachbauten stimmen. Form: **neuer Sensor oder Beleg-Form**: Die Tabelle nennt je Zeile, welcher Test aus welchem Grund rot wird, und das Skript der Mutanten liegt (mit Absicht ohne Spur im Repo) nur als Beschreibung vor; ein Sensor ist nicht billig. Register-Eintrag `weiter offen` oder Zähler an `BEO-REPO/beleg-traegt-zusage-aus-anderem-grund`; Urteil des Planners, ob es dieselbe Beobachtung ist.
4. **Doppelte Fall-Tabelle je Format** (Review F-4): Grenze, kein Mangel am Vertrag. Beobachtung für das Register mit Adresse, wenn ein Folge-Slice (`slice-v1-abschluss-verbindungsende-replay`, `-einspielen`) wieder eine dritte Tabelle braucht; bis dahin „weiter offen“ als Register-Zeile ohne Zähler-Zuwachs. Nach `AGENTS.md` §3.13 braucht eine Adresse den Eintrag im Nehmer; diese Eintragung ist nicht erfolgt (der Plan sagt es), also kein Ausgang „geplant“.

## Findings

### V-1 LOW

- `quelle`: `AGENTS.md` §3.10, §3.11
- `pfad`: `internal/adapters/driven/recording/sqlite_zeilen.go` · Kommentar und Bedingung von `setzeVersionMitEnde` („jede andere Zahl geänderter Zeilen ist ein Fehler“, `if n != 1`); `spec/spezifikation.md` · `LH-FA-22.a` *Ende und Version beim Schreiben von `sqlite`* („ändert genau die eine Zeile `version` mit dem Wert 1“)
- `befund`: Geprüft ist „null Zeilen“ (drei Fälle), nicht „mehr als eine“. Mutant `n != 1` zu `n == 0`: alle Tests `TestSQLite(Version|Ende|SchreibenNimmt|Roundtrip)` grün (PASS). Die Seite ist erreichbar: In einer Datei, deren `meta` ohne Primärschlüssel mit zwei Zeilen `version`/`1` steht, lehnt das Original mit `PGR-E3001` ab, der Mutant ergänzt die Session und setzt beide Zeilen. Im Schema des Schreibers ist `key` Primärschlüssel; ein Fall der Seite ist nur eine fremde Datei. Der Code und die Spezifikation sagen die Seite zu, §6 und die Tabelle in §7 nennen sie nicht.
- `verifizierbar`: ja (Mutant und Sonde, Ausgabe oben)
- `klasse`: Zusage im Kommentar weiter als die Prüfung (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Vorschlag: Fall „zwei Zeilen `version` mit Wert 1“ in `TestSQLiteVersionZeileUnpassend` (Datei per `rohSchreibe` ohne Primärschlüssel neu angelegt) samt Mutation `n == 0` in §7, oder Kommentar und Spezifikation auf „keine Zeile“ enger fassen; Entscheidung beim Planner/Architect, nicht beim Verifier.

### V-2 INFO

- `quelle`: `AGENTS.md` §3.11
- `pfad`: `internal/adapters/driven/recording/sqlite_schreiben_ende_test.go` · Doc von `TestSQLiteRoundtripMitEnde`; `docs/user/abdeckung-unit.md` · drei Zeilen zum Roundtrip (`LH-FA-22`, `LH-FA-25`, `LH-QA-06`)
- `befund`: „liest gleich (auch in `version`) und gleich der geschriebenen“. Der Test vergleicht SQLite und YAML ohne Ausnahme, die geschriebene Aufzeichnung aber mit `Version = model.VersionMitEnde` (geschrieben 5, gelesen 2). „Gleich der geschriebenen“ gilt ohne diese Ausnahme nicht.
- `verifizierbar`: ja (Test gelesen, Zeile `want.Version = model.VersionMitEnde`)
- `klasse`: Zusage im Doc weiter als die Prüfung
- Vorschlag: „… gleich der geschriebenen, in `version` mit dem Wert 2“ in Doc und Abdeckungszeilen (`make abdeckung`).

### V-3 INFO

- `quelle`: `AGENTS.md` §3.11
- `pfad`: Doc von `TestSQLiteEndeDarstellung` („eine Aufzeichnung ohne Ende schreibt dieselben Zeilen wie bisher“); Plan §6 *Zählung der Mutationstabelle* („alle Zeilen der vier Tabellen wie bisher“)
- `befund`: Der Fall „ohne Ende“ prüft `meta` (format, version), von `interaction` die Spalten `sequence` und `type`, von `message` `sequence`, `position`, `direction`, `group_no`, `kind`. Nicht geprüft sind die Tabelle `session` und die Spalten `sql`, `offset_ms`, `fields`. Der Bestand (`sqlite_schreiben_test.go`) deckt sie, der neue Test sagt sie aber zu.
- `verifizierbar`: ja (Test gelesen)
- `klasse`: Zusage im Doc weiter als die Prüfung
- Vorschlag: „dieselben Schlüsselspalten“ statt „dieselben Zeilen“, oder die Spalten in den Fall aufnehmen.

### V-4 INFO

- `quelle`: `AGENTS.md` §3.9
- `pfad`: `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-sqlite-schreiber.md` · Kopf, `Bezug`
- `befund`: Test und Abdeckungstabelle deklarieren `LH-QA-06` (Messung), §7 nennt die Zeile; der Kopf führt `LH-QA-06` nicht (`Bezug`: `LH-FA-25`, `LH-FA-22`, [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0013](../plan/adr/0013-sqlite-recording-backend.md)). `make kopf-check` prüft nur §1 und §2 und ist nicht rot.
- `verifizierbar`: ja (`grep -n 'LH-QA-06'` im Plan)
- `klasse`: Kopf folgt dem Diff nicht (`BEO-REPO/plan-folgt-korrektur-nicht`, kleine Form)
- Vorschlag: `LH-QA-06` in `Bezug` aufnehmen oder die Deklaration im Test streichen, wenn der Slice die Messung nicht zusagt.

### V-5 INFO

- `quelle`: `AGENTS.md` §3.9
- `pfad`: Plan §6, zweites Risiko („gemessen: 596 in der Metrik von §4 (§7), 60 % der Marke“)
- `befund`: Nach der Nacharbeit sind es 595 (§7 erster Absatz und §8); §6 trägt den Stand vor der Nacharbeit. Differenz eine Zeile im Doc des Roundtrips.
- `verifizierbar`: ja (`git diff --numstat b5933e8..HEAD -- internal docs/user`)
- `klasse`: Plan folgt der Nacharbeit nicht ganz (kleine Form)
- Vorschlag: 596 zu 595 in §6 bei der Closure.

### V-6 INFO

- `quelle`: `AGENTS.md` §3.10 (Fall und Mutation, die nur ihn bricht)
- `pfad`: Plan §7, Zeile `zeileFehltOk`; Test `TestSQLiteVersionZeileUnpassend`/„Zeile fehlt“
- `befund`: Für „Zeile fehlt“ allein gibt es keine Mutation, die nur ihn färbt (Code unterscheidet ihn nicht von „Zeile fehlt, andere Zeile mit Wert 1“). §7 benennt das. Kein Mangel, eine benannte Grenze; der Fall bleibt als Eingang der Lesart stehen.
- `verifizierbar`: ja (`zeileFehltOk` und `ohneKey` gefahren)
- `klasse`: akzeptiertes Negativ, benannt
- Vorschlag: keine Aktion.

## Zusammenfassung

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Zusage im Kommentar oder Doc weiter als die Prüfung (V-1 bis V-3) · Kopf und Plan folgen der Nacharbeit nicht ganz (V-4, V-5) · benannte Grenze (V-6).
