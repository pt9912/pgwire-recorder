# Verifikation slice-v1-abschluss-verbindungsende-sqlite

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-15. **Spanne:** `fa8486b` (Basis) bis `654ad5a` (HEAD). **Gegenstand:** Slice `slice-v1-abschluss-verbindungsende-sqlite` (Plan unter `docs/plan/planning/in-progress/`, Schritt 36a, SQLite-Leser) gegen DoD (§2), Plan §1 (Abgrenzung), §3, §4, §6 (Randformen, Risiken), §7 (Belege des Implementers), §8, `SPEC-043` (*Ende des Servers*), `SPEC-001`, `SPEC-002`, `LH-FA-25` mit `LH-FA-25.a`, [`ADR-0039`](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md) und `AGENTS.md` §3.7 bis §3.14. Review: `docs/reviews/2026-10-15-review-slice-v1-abschluss-verbindungsende-sqlite.md` (F-1 LOW, F-2 und F-3 INFO, nachgearbeitet in `654ad5a`); gegen ihn ist nicht verifiziert, F-1 bis F-3 sind unabhängig nachgeprüft.

Eine Verifier-Vorlage führt die vendored Baseline nicht; Aufbau wie der Report von `slice-v1-abschluss-verbindungsende-format`. Belege sind in diesem Lauf selbst erzeugt: Quelle von `HEAD` (`git archive`: `internal`, `cmd`, `test`, `go.mod`, `go.sum`) in eine Kopie im Scratchpad, Läufe per Bind-Mount mit dem Image `pgwire-recorder:test` (`--network=none`, `go test` mit `-run 'TestLesenEnde|TestLesenMeta|TestLesenMeldung|TestLesenReihenfolge|TestLesenUserVersion|TestLesenGleiche'` im Paket `./internal/adapters/driven/recording/`), jeder Mutant einzeln aus einer unveränderten Kopie angewandt und danach verworfen. Ein Lauf zur Zeit, kein `pkill`, kein `make gates`. Der Arbeitsbaum des Repos ist unverändert (`git status` zeigt nur diesen Report).

**Verdikt: PASS mit Befunden.** Alle Liefer-Punkte der DoD sind bestätigt; jeder nachgefahrene Mutant färbt genau die in §7 genannten Tests, beide äquivalenten Mutanten sind tatsächlich äquivalent, die Abgrenzung ist eingehalten, der Plan folgt dem Diff, die Größe ist gemessen 387. 0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO. Offen sind allein die Closure-Pflichten (Notiz, Lerneintrag, Register, Risiko-Ausgang, Paarungen); das ist der erwartete Stand vor der Closure.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren oder gelesen) |
|---|---|---|
| 1: Leser nimmt `meta.version` 1 und 2 an (auch 2 ohne Ende); `connection_closed` in Version 1, falsche Stelle, nicht nach `FATAL`/`PANIC`, mit Feldern sind `PGR-E3003`; Ende vor der letzten Interaktion und zweites Ende `PGR-E3003` über `Session.ValidateEnd`; jede andere Version `PGR-E3002` (Test je Randform, Dateien von Hand) | bestätigt | Code gelesen: `sqlite_lesen.go` (`liesMeta`/`pruefeMeta` liefern den Text der Version, Vergleich mit `strconv.Itoa(model.Version)` und `strconv.Itoa(model.VersionMitEnde)`, `liesAufzeichnung` schreibt ihn in den Baum), `sqlite_lesen_baum.go` unverändert, `yaml.go` ruft `s.ValidateEnd()` in `sessionFromDTO`. Unveränderte Kopie: `./internal/hexagon/...` und das ganze Paket `recording` grün (52 s). Tests `TestLesenEndeAngenommen` (15 Fälle), `TestLesenEndeAbgelehntDarstellung` (11), `TestLesenEndeAbgelehntStelle` (17), `TestLesenEndeVersionen` (13), `TestLesenEndeJeSession` (3), `TestLesenEndeNummernOrdnen` (2), dazu `TestLesenMeta`; Dateien von Hand mit SQL gebaut (`endeSQL`), nicht vom Schreiber. Je Randform aus §6 ein Fall (Liste mit §6 Zeile für Zeile abgeglichen: Darstellung, `fields`-Schlüssel, `direction`, `kind`, Nummern statt Zeilenfolge, Version 1/2/ohne Ende, kanonische Zahlen, Schreibweisen, `user_version`, Rangfolge, 17 Stellen). |
| 2: je Zusage eine Mutation, vom Implementer rot gesehen; Tabelle in §7 | bestätigt, 24 Mutanten nachgefahren | Tabelle in §7 (zwei Tabellen: Umsetzung und Nacharbeit). Nachgefahren, alle wie behauptet: siehe nächster Abschnitt. Zu eng gefasst sind zwei Aussagen (V-1, V-2). |
| 3: `make gates` grün | durch den Hook getragen, nicht selbst gefahren | §7 *Läufe*: voller `make gates` grün (Exit 0, 12 min 56 s) nach der Nacharbeit. `echo '{}' \| bash .claude/hooks/stop-require-gates.sh` gibt `{"decision":"approve"}` (Exit 0); `.harness/state/gates-passed.dateiliste` 17:52:01, Commit `654ad5a` 17:52:16, Arbeitsbaum sauber. Ob der volle Lauf aus dem richtigen Grund grün war, belegt die Datei nicht; ich habe `make gates` nicht wiederholt (Auftrag). Eigene Läufe: `./internal/hexagon/...` und `./internal/adapters/driven/recording/` grün. |
| 4: Review durchgeführt, Report liegt vor | bestätigt | Report unter `docs/reviews/` (0 HIGH, 0 MEDIUM, 1 LOW, 2 INFO). F-1 bis F-3 in `654ad5a` nachgearbeitet, unten nachgeprüft. Kein Self-Review. |
| 5 bis 8: Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgang, Paarungen | offen (erwartet) | §7 trägt Belege des Implementers, noch keine Closure-Notiz; §2 ohne Häkchen; das Risiko in §6 trägt „offen bis Closure“. Handlung des Planners; Vorschläge unten. |

## Mutationen selbst nachgefahren

„nur“ heißt: keine anderen Tests rot. Quelle der Mutanten: die Zeile in §7 (Umsetzung) oder *Mutationen der Nacharbeit*.

| Mutant | Rot | Urteil |
|---|---|---|
| Baum trägt fest `"1"` | `TestLesenEndeAngenommen` (alle Fälle mit Ende `PGR-E3003`; „Version 2 ohne Ende“ liest Version 1), `…JeSession`, `…NummernOrdnen`; dazu Fall `user_version 1, meta.version 2` | wie §7 |
| Baum trägt fest `"2"` | `TestLesenEndeAngenommen` („Version 1 ohne Ende“, „user_version 2, meta.version 1 ohne Ende“ lesen Version 2), `TestLesenEndeAbgelehntDarstellung` („Version 1, einfach/extended“, „user_version 2, meta.version 1, Ende“ angenommen), `TestLesenMeta` („gelockerte Datei“), `TestLesenUserVersionIrrelevant`, `TestLesenGleicheAblehnungBeiderFormate`, `TestLesenReihenfolgeDerZeilen` | wie §7 (mehr als genannt, §7 nennt „weitere Fälle“ nicht; kein Widerspruch) |
| nur `model.Version` angenommen | `TestLesenEndeAngenommen` (alle Fälle der Version 2), `TestLesenEndeAbgelehntDarstellung`/`…Stelle` (Code `PGR-E3002` statt `PGR-E3003`) | wie §7 |
| Menge in `pruefeMeta` entfällt (`if version == strconv.Itoa(-99)`) | `TestLesenEndeVersionen` (Meldung nennt die Zahl, 0/3/-1/10/20, beide „Ende vor der Fehlerantwort“-Fälle; Zahl über 64 Bit: **Code** `PGR-E3003` statt `PGR-E3002`); `TestLesenMeta`: „Version vor Inhalt“, „Version vor fehlender Tabelle“, „Version vor unbekanntem Schlüssel“ (**Code** `PGR-E3003` statt `PGR-E3002`), beide Zahlen über 64 Bit (Code); `TestLesenMeldungOhneWertUndZeile` (Meldung nennt `4711`) | rot, aber **mehr am Code als §7 sagt** (V-1) |
| `pruefeMeta` meldet `Fassung <Zahl>` statt „Version der Datei …“ | `TestLesenEndeVersionen` (alle acht E3002-Fälle), `TestLesenMeldungOhneWertUndZeile` | wie §7 (F-1, Ausprägung 3); nur Meldungstext, Code bleibt `PGR-E3002` |
| Menge entfällt und `kopfVorpruefung` meldet `Fassung <Zahl>` | `TestLesenEndeVersionen`, `TestLesenMeta`, `TestLesenMeldungOhneWertUndZeile` | wie §7 |
| nur `kopfVorpruefung` meldet `Fassung <Zahl>` (`yaml.go`) | keiner, `ok` | **grün, äquivalent, bestätigt:** `pruefeMeta` lässt nur `"1"` und `"2"` an `Unmarshal`, `kopfVorpruefung` lehnt diese nie ab; die Meldung ist für den SQLite-Leser nicht erreichbar |
| Vergleich als Präfix (`HasPrefix` `"1"`/`"2"`) | `TestLesenEndeVersionen` „Version 10“, „Version 20“ (zwei Zeilen; die Meldung nennt die Zahl, Code bleibt `PGR-E3002`) | wie §7 |
| Version erst nach dem Inhalt geprüft (fremde Version trägt `2` im Baum, `PGR-E3002` folgt nach `Unmarshal`) | `TestLesenEndeVersionen` „Version 3, Ende vor der Fehlerantwort“ (und „Zeilen umgekehrt“), `TestLesenMeta` „Version vor Inhalt/fehlender Tabelle/unbekanntem Schlüssel“, `TestLesenMeldungOhneWertUndZeile` | wie §7 |
| Leser liest `PRAGMA user_version` | `TestLesenEndeAngenommen` („user_version 1, meta.version 2, Ende“, „user_version 2, meta.version 1 ohne Ende“), `TestLesenEndeAbgelehntDarstellung` („user_version 2, meta.version 1, Ende“), `TestLesenUserVersionIrrelevant` („user_version 2“) | wie §7 |
| Aufruf von `Session.ValidateEnd` in `sessionFromDTO` entfällt | nur `TestLesenEndeAbgelehntStelle` („Ende in der Interaktion 1 von 2“, „Ende in zwei Interaktionen“), `TestLesenEndeJeSession` (dieselben zwei Fälle) | wie §7 |
| Sortierung der Interaktionen entfällt | `TestLesenEndeNummernOrdnen` „Zeilen von interaction umgekehrt“ (`PGR-E3003` über die Prüfung der Nummern), `TestLesenReihenfolgeDerZeilen/interaction` | wie §7, mit der Grenze aus §7 (V-4) |
| Sortierung der Sessions entfällt | `TestLesenEndeNummernOrdnen` „Zeilen von session umgekehrt“, `TestLesenReihenfolgeDerZeilen/session` | wie §7 |
| Sortierung der Nachrichten entfällt | `TestLesenEndeAngenommen` (beide „Zeile des Endes vor der Fehlerantwort“), `TestLesenReihenfolgeDerZeilen/message` | wie §7 |
| Zeile `connection_closed` ausgelassen, einfach (`knoten`) | `TestLesenEndeAngenommen`, `…AbgelehntStelle`, `…JeSession`, `…NummernOrdnen` | wie §7 |
| Zeile ausgelassen, Extended (`gruppen`) | `TestLesenEndeAngenommen`, `…AbgelehntStelle` (sechs Fälle) | wie §7 |
| `fields` des Endes verworfen | nur `TestLesenEndeAbgelehntDarstellung`: „fields mit tag“, „mit Wert null“, „mit leerem Namen“, „mit type“ | wie §7 (die Einzelmutanten je Schlüssel habe ich nicht nachgefahren) |
| `kind` in Kleinbuchstaben gewandelt | nur „kind in Großbuchstaben“ | wie §7 |
| `kind` von Leerraum befreit | nur „kind mit Leerzeichen“ | wie §7 |
| `direction client` bei `connection_closed` in Extended routet in die Server-Liste (`gruppen`) | keiner, `ok` | **grün, äquivalent, bestätigt** (siehe Bewertung 2) |

Nicht nachgefahren: die Einzelmutanten je Schlüssel von `fields` (`tag`, `null`, leerer Name, `type`), `direction client` einfach, die Schreibweisen `02`, `+2`, ` 2`, `2 `, der Mutant `0`/`3`/`-1` angenommen, `beendetVerbindung` ohne Beachtung der Schreibweise; sie liegen im Bestand oder im Gemeinsamen Code und sind aus den Tabellen der Fälle abgeleitet, nicht gefahren.

## Bewertung der vier Fragen

**(1) Genügt die Grenze „Menge in `pruefeMeta` für kleine Zahlen nur über den Meldungstext rot“ (§3.10, §3.11), oder ist die Menge redundanter Code?** Die Menge ist **nicht redundant**, und die Grenze ist enger zu fassen als geschrieben. Gefahren (Mutant „Menge entfällt“): Ohne sie ändert sich der Code für 0, 3, -1, 10, 20 mit sauberem Inhalt tatsächlich nicht (`Unmarshal` lehnt sie mit `PGR-E3002` ab) und nur der Text wird rot; die Zusage dazu („Meldung nennt keine Version“, `SPEC-043`) ist durch `TestLesenEndeVersionen` und `TestLesenMeldungOhneWertUndZeile` geprüft, nicht an einem Wortlaut (der Test prüft die Ziffernfolge nach Abzug von Pfad und Code, und drei Ausprägungen sowie die Variante in `kopfVorpruefung` sind rot). Aber die Menge trägt zwei weitere Zusagen, die **am Code** rot werden: die Rangfolge *Version vor Inhalt* (ohne Menge laufen `liesSessions` und die Tabellen vor `kopfVorpruefung`; `TestLesenMeta` „Version vor Inhalt“, „…fehlender Tabelle“, „…unbekanntem Schlüssel“ liefern `PGR-E3003` statt `PGR-E3002`) und die Zahl über 64 Bit (`PGR-E3003` statt `PGR-E3002`). Damit erfüllt die Menge §3.10 (jede Zusage hat eine Mutation, die nur sie bricht) und §3.11 (die Prosa sagt zu, was ein Test prüft). Was nicht stimmt, ist der Satz der Anmerkung in §7 „für diese Zahlen ist die Menge … am Code nicht zu isolieren … sie wird nur über den Text der Meldung rot“: Für `3` mit unsauberem Inhalt ist sie am Code isolierbar (V-1). Für `10`, `20`, `0`, `-1` gibt es keinen Fall *Version vor Inhalt*; eine Mutation, die nur diese Zahlen an `Unmarshal` durchlässt, bleibt für die Rangfolge unentdeckt (aus dem Code abgeleitet, nicht gefahren; `SPEC-043` und §6 nennen die Rangfolge allgemein, den Fall nur mit `3`).

**(2) Die äquivalenten Mutanten.** Beide sind äquivalent, beide sind gefahren und grün.

- *„Richtung `client` in Extended“:* In einer Extended-Gruppe stehen die Client-Nachrichten vor den Server-Nachrichten (`gruppen`; eine Client-Nachricht nach einer Server-Nachricht ist `PGR-E3003`). Routet man `connection_closed` mit `direction client` in die Server-Liste, landet sie dort als **erste** Server-Nachricht und damit vor der `error_response`; `pruefeEnde` lehnt ein Ende ohne unmittelbar vorangehendes `FATAL`/`PANIC` mit demselben Code ab. Die Aussage in §7 („der Fall belegt den Code, nicht die Prüfung der Richtung“) trägt, die Begründung ist die Ordnung der Gruppe, nicht eine Eingabe, an die jemand gedacht hat (Gegenprobe zu `BEO-REPO/gruener-mutant-faelschlich-aequivalent`: kein Fall gefunden, der Original und Mutant unterscheidet).
- *„Nur `kopfVorpruefung` meldet Fassung“:* Die Meldung von `kopfVorpruefung` ist für den SQLite-Leser nicht erreichbar, weil `pruefeMeta` jede andere Zahl als `1` und `2` vorher ablehnt. Das ist wahr und bestätigt (kein Fall nötig). Es ist zugleich der Grund, warum die YAML-Aussage („die Meldung nennt die Version“) und die SQLite-Aussage („nennt keine“) nebeneinander bestehen (`SPEC-001`/`SPEC-043`, Drift-Log-Zeile).

**(3) „Richtung client in Extended“ steht in §7 statt in §6: Befund nach `AGENTS.md` §3.12 oder nicht?** **Kein Befund.** §3.12 verlangt, dass die *Randform* in §6 genannt und vor dem Code in der Spezifikation entschieden ist. Das ist sie: §6 nennt „`direction` `client` mit `kind` `connection_closed` ist `PGR-E3003`“ (Abschnitt *Darstellung des Endes*), `SPEC-043` (*Ende des Servers*) entscheidet sie, beides im Architect-Commit `f449ea1` vor `dbd2c3a`. Was in §7 steht, ist keine Randform, sondern eine **Eigenschaft des Belegs** (die Mutation an der Richtung in Extended ist äquivalent); sie fiel erst beim Fahren der Mutation an und gehört nach `AGENTS.md` §3.10 in die Mutationstabelle. Ein Eintrag der Äquivalenz schon in §6 wäre die Vorbild-Form „akzeptiertes Negativ“ (§6 führt sie an mehreren Stellen), war aber vor dem Code nicht zu wissen und keine Randform, die der Code still entscheidet.

**(4) Abgrenzung.** Eingehalten (Tabelle unten).

## Nacharbeit (F-1 bis F-3, Stichprobe)

| Prüfung | Ergebnis |
|---|---|
| F-1: neue Prüfung „Meldung nennt die Zahl in keiner Form“ | `pruefeEndeFaelle` prüft bei `PGR-E3002` die Meldung nach Abzug von `ohnePfad` und dem Code auf die Ziffernfolge der Version. Gefahren: `Fassung <Zahl>` in `pruefeMeta` rot (alle acht E3002-Fälle); Menge entfällt (Meldung aus `Unmarshal`) rot; Menge entfällt und `kopfVorpruefung` meldet `Fassung <Zahl>` rot; nur `kopfVorpruefung` meldet `Fassung <Zahl>` grün und äquivalent. Wie §7. „in keiner Form“ ist weiter als die Prüfung (V-2). |
| F-2: Anmerkung zur Zahl über 64 Bit an den Code angeglichen | Bestätigt: `pruefeMeta` vergleicht Texte (`strconv.Itoa` der Konstanten, keine Umwandlung der Datei); fällt die Menge weg, ist der Code für die Zahl über 64 Bit `PGR-E3003` (`kopfVorpruefung` lehnt den als `!!int` gesetzten Text ab), der Test ist am Code rot. |
| F-3: Fall mit Ende für die Ordnung von Interaktionen und Sessions | `TestLesenEndeNummernOrdnen` (zwei Fälle) rot unter beiden Sortier-Mutanten, im Gutfall grün. Die Session-Zeilen stehen über `CREATE TABLE … AS SELECT … ORDER BY id DESC` umgekehrt, wirkt (Mutant rot). Grenze wie in §7: Der Test wird über die Prüfung der Nummern rot, bevor `ValidateEnd` läuft (V-4). |
| `Befund vorher und nachher` in §7 | `grep -n 'Value: "1"\|version != "1"' sqlite_lesen.go`: kein Treffer; `grep -n 'strconv.Itoa(model.Version'`: Zeile 256 (beide Konstanten). Wie behauptet. |

## Abgrenzung, Plan und Größe

| Prüfung | Ergebnis |
|---|---|
| §1 Abgrenzung | `git diff --name-only fa8486b..HEAD` außerhalb `docs/plan/` und `docs/reviews/`: `internal/adapters/driven/recording/{leser_test,sqlite_lesen,sqlite_lesen_ende_test,sqlite_lesen_inhalt_test}.go`, `docs/user/abdeckung-unit.md`, `spec/spezifikation.md` (Architect, `f449ea1`), `test/integration/sqlite_lesen_e2e_test.go` (eine Zeile, `2` wird `3`). **Schreiber (36b) unverändert:** keine `sqlite_schreiben*.go` im Diff. **Kein Handbuch, kein README:** `docs/user/benutzerhandbuch.md` und `README.md` nicht im Diff (`abdeckung-unit.md` ist Werkzeug-Ausgabe, §3 nennt sie). **Kein Aufruf von `ValidateEnd` im SQLite-Leser:** `grep -n ValidateEnd internal/adapters/driven/recording/sqlite_lesen*.go` ohne Treffer außerhalb der Tests; der Aufruf liegt in `sessionFromDTO` (`yaml.go`, unverändert). Kein Eingriff in `tools/schema`, `sqlite_lesen_baum.go`, `sqlite_json_lesen.go`, `model`. Eingehalten. |
| §3 (Plan) gegen `git diff --stat` | Jede Code-Datei des Diffs steht in §3 (`sqlite_lesen.go`, `sqlite_lesen_ende_test.go`, `sqlite_lesen_inhalt_test.go`, `leser_test.go`, `sqlite_lesen_e2e_test.go`, `abdeckung-unit.md`); `sqlite_lesen_baum.go` und `sqlite_json_lesen.go` stehen als „bleiben unverändert“ und sind es. `spec/spezifikation.md` steht in keiner Zeile von §3, ist aber Gegenstand des Architect-Commits und in §6 benannt; kein Mismatch (Review: Negativbefund). |
| §6 Randformen | Jede vor dem Code entschieden (`f449ea1` vor `dbd2c3a`), in `SPEC-043` (*Ende des Servers*), `SPEC-001` (Satz zur Meldung bei `sqlite`), Drift-Log-Zeile. Der Diff entscheidet keine Randform, die §6 nicht nennt (Stichprobe: `fields` mit `null`/leerem Namen/`type`, `user_version`, Zahlen `10`/`20`/über 64 Bit, Rangfolge, 17 Stellen); §7 *Randformen* nennt keine neue. |
| §7 gegen Diff | Tabellen, Läufe, Größe, *Befund vorher und nachher* folgen dem Diff. Abweichungen: V-1 (Anmerkung zu eng), V-2 („in keiner Form“ weiter als die Prüfung). |
| Größe (Metrik §4: `git diff --numstat fa8486b..HEAD -- internal docs/user`) | Selbst addiert aus den fünf Zeilen: 359 Einfügungen plus 28 Löschungen = **387**, gleich dem Wert des Implementers. Produktcode 39 (21 + 18), Tests 338 (322 in der neuen Datei, 14 in `sqlite_lesen_inhalt_test.go`, 2 in `leser_test.go`), Tabelle 10. **55 % der Marke 700**, unter der Meldeschwelle 550, in der Schätzung des Architect (220 bis 420, Mitte 320; die Messung liegt 21 % über der Mitte, nahe am oberen Ende). Faktor Tests zu Produktcode 338 / 39 ≈ 8,7, in der Spanne der Datenpunkte 8 und 11 und über dem Ansatz 5. Ein Liefer-Punkt, zwei Schichten (Recording-Adapter, Nutzer-Doku-Tabelle). Außerhalb der Metrik: 1 Zeile in `test/integration/`. |
| §8 | Zählung (Planner, Architect, Implementer) steht; Messung 330 beim Umsetzen, 387 nach der Nacharbeit in §7; konsistent. |
| `AGENTS.md` §3.9 | Der Plan folgt dem Diff in §1, §3, §6, §7; Kopf trägt `LH-FA-25`, `LH-FA-22`, `LH-QA-06`, [`ADR-0039`](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [`ADR-0013`](../plan/adr/0013-sqlite-recording-backend.md), `LH-FA-25.a`, `LH-FA-22.a`, `SPEC-001`, `SPEC-002`, `SPEC-043`, `ARC-008`. `SPEC-041` fehlt im Kopf (V-3). |
| `AGENTS.md` §3.11 | Die sechs Abdeckungs-Deklarationen und -Zeilen in `abdeckung-unit.md` entsprechen den Fall-Listen der Tests (Fälle gelesen); Ausnahme „in keiner Form“ (V-2). |
| `AGENTS.md` §3.13 / Bindungen der Nehmer | Alle fünf Nehmer tragen die Sendung im selben Commit (`ead8132`) und haben ihre Zählung in §8 mit der Kennung des Gebers: `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Schreiber und Roundtrip; Bindungen aus `-format` und `-format-schreiber`; Liefer-Punkte 2, Schichten 1; „Verantwortlich: —“ in `next/`), `-record` (Kennung in der Abgrenzung, Liefer-Punkte 3, Schichten 2), `-upstream` (Übernahme aus `-sqlite-schreiber`, Start nennt beide Hälften, Liefer-Punkte 3, Schichten 2), `-replay` (Leser als Vorgänger 36a, Start nach 36a und 36b, Liefer-Punkte 3, Schichten 2), `-doku` (Ziel nennt `-sqlite-schreiber`; Liefer-Punkte 2, Schichten 1). Die Adressen nehmen an: keiner liegt in `done/`. `welle-v1-abschluss` führt beide Zeilen und die Schritte 36a/36b; die Roadmap („In Arbeit: Schritt 36a, Als Nächstes Schritt 36b“) und das Drift-Log tragen den Schnitt. Der Geber `-format-schreiber` (in `done/`) nennt weiter die Kennung `-sqlite` für den Roundtrip: ein Hop (die Datei des Slice nennt in §1 den Nehmer `-sqlite-schreiber`), kein Befund. |

## Risiko aus §6 und Vorschlag für den Ausgang

| Risiko | Vorschlag | Begründung |
|---|---|---|
| Der Diff liegt über einer Review-Sitzung (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`) | **entfallen** | Gemessen 387 in der Metrik von §4 (55 % der Marke 700, unter der Meldeschwelle 550, innerhalb der Schätzung 220 bis 420 des Architect); das Review lief in einer Sitzung ohne Größen-Finding (Summary 0/0/1/2). Der Planner hat vor dem Anspruch geschnitten, der Schreiber liegt in `-sqlite-schreiber`. Der Faktor Tests zu Produktcode (8,7) und die Lage nahe dem oberen Ende der Spanne gehören als Datenpunkt ins Register zu `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (kein Zählerzuwachs: die Schätzung hat getragen). |

Kein weiteres Risiko steht in §6; ein neues entsteht nicht (V-1 bis V-4 sind Berichtigungen, keine Risiken).

## Lerneintrag-Kandidaten (Entscheidung beim Planner)

1. **Zusage hängt am Wortlaut der Meldung** (Review F-1; hier die Klasse „Test enger als die Zusage“): Die Zusage „die Meldung nennt keine Version“ wurde zuerst gegen den Wortlaut `Version <Zahl>` geprüft; eine andere Form (`Fassung`, `version=`) ließ die Mutanten grün. Behoben mit der Abwesenheitsprüfung der Ziffernfolge. Form: **geschärfte Regel**, Ort `.claude/commands/implement-slice.md` Schritt 19 (Beleg einer Zusage) mit dem Satz „eine Abwesenheitszusage über eine Meldung prüft die Form des Werts, nicht den Wortlaut der Meldung, und der Beleg fährt mindestens zwei andere Wortlaute“. Register: neuer Eintrag oder Zähler an `BEO-REPO/test-prueft-teilstring-gegen-meldung-mit-zufallspfad` (Teilstring gegen Meldung); das Planner-Urteil, ob sie dieselbe Beobachtung sind (hier ist der Test zu eng, dort fehlerempfindlich), steht aus.
2. **Beleg beschreibt Mutation anders als der Code** (Review F-2; in diesem Lauf eine zweite Instanz, V-1): Die Zeile „Überlauf der Umwandlung“ nannte eine Umwandlung, die es nicht gibt; die Anmerkung hält die Menge nach der Nacharbeit weiter für „nur über die Meldung“ rot. Die Nähe zu `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (§7-Form) und `BEO-REPO/beleg-traegt-zusage-aus-anderem-grund` ist belegt; Zähler dort, mit diesem Slice als Beleg (Form: Beleg, der eine Mutation beschreibt, die der Verfasser nicht gefahren hat, sondern aus dem Code abgeleitet; Sensor wäre: der Beleg nennt je Zeile, welcher Test mit welchem Grund rot wurde, nicht „über die Meldung“ allgemein).
3. **Äquivalenter Mutant belegt** (Review F-3): Beide Äquivalenzen halten (gefahren). Das ist ein Gegenbeleg zu `BEO-REPO/gruener-mutant-faelschlich-aequivalent` und im Register nur als Hinweis „in diesem Slice nicht aufgetreten“ zu führen; die Begründung über eine Invariante der Struktur (Client vor Server in der Gruppe) statt über erdachte Eingaben ist die tragfähige Form. Keine neue Regel.

## Findings

### V-1 INFO

- `quelle`: `AGENTS.md` §3.9, §3.10
- `pfad`: `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-sqlite.md` · §7, *Anmerkung zu den Zahlen* und *Nacharbeit zum Review*, Satz „die Menge in `pruefeMeta` für die kleinen Zahlen bleibt am Code nicht isolierbar“
- `befund`: Der Satz sagt, für 0, 3, -1, 10, 20 werde die Menge nur über den Text der Meldung rot. Gefahren: Fällt die Menge weg, werden `TestLesenMeta` „Version vor Inhalt“, „Version vor fehlender Tabelle“ und „Version vor unbekanntem Schlüssel“ am **Code** rot (`PGR-E3003` statt `PGR-E3002`); die Tabelle in §7 führt diese Fälle bei der Zeile `3` selbst, Anmerkung und Nacharbeit widersprechen ihr. Zutreffend ist: Mit sauberem Inhalt ist die Menge für diese Zahlen im Code nicht isolierbar; als Ganzes ist sie es (Rangfolge Version vor Inhalt). Die Menge ist nicht redundant. Außerdem trägt der Fall *Version vor Inhalt* nur die Zahl `3`.
- `verifizierbar`: ja (Mutant „Menge entfällt“, Ausgabe oben)
- `klasse`: Beleg beschreibt Mutation anders als der Code sie erlaubt (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, §7-Form)
- Vorschlag: Den Satz vor der Closure ersetzen durch „mit sauberem Inhalt wird die Menge für diese Zahlen nur über die Meldung rot; als Ganzes ist sie am Code rot (`TestLesenMeta` *Version vor Inhalt*, Zahl über 64 Bit)“; ob ein Fall *Version vor Inhalt* mit `10` oder `20` fehlt, entscheidet der Planner (Rangfolge-Zusage in `SPEC-001` ist allgemein).

### V-2 INFO

- `quelle`: `AGENTS.md` §3.11
- `pfad`: `internal/adapters/driven/recording/sqlite_lesen_ende_test.go` · Kommentar von `TestLesenEndeVersionen`; `docs/user/abdeckung-unit.md` · Zeile zu `TestLesenEndeVersionen`; Plan §7, Zeile „Meldung ohne die Zahl“
- `befund`: Die Zusage lautet „enthält die Zahl der Datei in keiner Form“. Geprüft wird, dass die Ziffernfolge der Version (nach Abzug von Pfad und Code) nicht als Teilstring in der Meldung steht. Eine Form ohne diese Ziffernfolge (ausgeschrieben, umgerechnet, `Fassung <Zahl+1>`) bliebe grün; für `-1`, `10`, `20` auch die gedrehte Ziffernfolge. Aus dem Code abgeleitet, nicht mit einer Sonde gefahren. Die drei gefahrenen Ausprägungen und die Variante in `kopfVorpruefung` sind rot, der Befund ist klein.
- `verifizierbar`: ja (Code des Tests gelesen)
- `klasse`: Zusage in Kommentar und Tabelle weiter als die Prüfung (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Vorschlag: „… die Ziffernfolge der Zahl der Datei steht nirgends in der Meldung“ in Kommentar, Abdeckungszeile (`make abdeckung`) und §7 statt „in keiner Form“.

### V-3 INFO

- `quelle`: `AGENTS.md` §3.9; `kopf-check`
- `pfad`: `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-sqlite.md` · Kopf, `Berührte Spec-Stellen`
- `befund`: §6 benennt `SPEC-041` als Ort der Stellen des Endes (Absatz *Stellen des Endes*, `SPEC-002`, `SPEC-041`); der Kopf führt `SPEC-041` nicht. `make kopf-check` prüft nur §1 und §2 und ist deshalb nicht rot. Der Vorgänger `-format` führte `SPEC-041`.
- `verifizierbar`: ja (`grep -n SPEC-041` im Plan)
- `klasse`: Kopf folgt §6 nicht (kleine Form von `BEO-REPO/plan-folgt-korrektur-nicht`)
- Vorschlag: `SPEC-041` in `Berührte Spec-Stellen` aufnehmen, oder in §6 streichen, wenn die Stellen nur in `SPEC-002` liegen.

### V-4 INFO

- `quelle`: `AGENTS.md` §3.10, §3.11
- `pfad`: `internal/adapters/driven/recording/sqlite_lesen_ende_test.go` · `TestLesenEndeNummernOrdnen`; `docs/user/abdeckung-unit.md` · Zeile Happy zur Ordnung
- `befund`: Der Test ist, wie §7 es sagt, kein Beleg einer eigenen Zusage des Endes: Beide Sortier-Mutanten werden `PGR-E3003` über die Prüfung der Nummern („trägt die Nummer 2“, „trägt die Kennung 2“), bevor `ValidateEnd` läuft, und `TestLesenReihenfolgeDerZeilen` färbt dieselben Mutanten. Es gibt keinen Mutanten, der „die letzte Interaktion ist die mit der höchsten `sequence`“ bricht und die Nummern-Prüfung überlebt. Die Abdeckungszeile (Happy, `LH-FA-25`) liest sich wie eine Zusage des Endes; die Grenze steht nur in §7. Kein falscher Satz, aber ein Test, der zwei Dinge belegt, die der Bestand schon belegt.
- `verifizierbar`: ja (beide Mutanten gefahren, Ausgabe oben)
- `klasse`: Beleg trägt die Zusage aus einem anderen Grund (`BEO-REPO/beleg-traegt-zusage-aus-anderem-grund`, kleine Form, im Plan benannt)
- Vorschlag: keine Aktion; bei der Closure die Grenze in einem Satz neben der Abdeckungszeile lassen oder den Test als Gegenfall zur Ordnung führen, nicht als Zusage des Endes.

## Zusammenfassung

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Beleg beschreibt Mutation anders als der Code · Zusage in Kommentar weiter als die Prüfung · Kopf folgt §6 nicht · Test belegt Zusage aus anderem Grund (Grenze benannt).
