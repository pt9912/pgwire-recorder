# Verifikation: slice-v1-abschluss-zeitangaben-pruefen — 2026-10-11

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-11-review-slice-v1-abschluss-zeitangaben-pruefen.md`, F-637 bis F-639); die Übergaben prüfe ich nach. Ich repariere nichts.

**Gegenstand:** Slice-Plan `slice-v1-abschluss-zeitangaben-pruefen` (Kopf, §1 bis §8) gegen den Gesamt-Diff `8e430ce..e4bb247` (11 Dateien, 624 Zeilen hinzu, 26 weg). Implementer: `78a6b4c`, `c3ba98f`, `341b569`, `2a316b1`, `3cf1f57`, `e4bb247`. Architect: `ae1996d` (Rangfolge der Prüfung in [`SPEC-001`](../../spec/spezifikation.md#spec-001--recording-serialisierung-und-versionierung), `version`-Regeln, Fixture akzeptiert). Review: `b220377`.

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingang:** der Plan ganz am Stand `e4bb247` (§7 *Belege des Implementers* nur als Beleg-Teile; §7 ist bis zur Closure unvollständig, das ist kein Befund); `spec/spezifikation.md` `LH-FA-21.a` (*Prüfen*, *Einspielen*), `SPEC-001` (*Rangfolge der Prüfung*) und `SPEC-043`; der Review-Report ganz; der Diff an `yaml.go`; die Tests `zeitangaben_pruefen_test.go`, `zeitangaben_pruefen_e2e_test.go`; die Nehmer `slice-v1-abschluss-zeitangaben` (§1, DoD) und die Geschwister im Schnitt; `AGENTS.md` §3.9 bis §3.13. Den Bericht des Implementers habe ich nicht gelesen. Beim Start war der Arbeitsbaum sauber, HEAD `e4bb247`.

**So wurden die Proben gebaut.** Alles Docker-only, Scratch-Präfix `ver-zp-`, nie im Repo.

- **Binary:** `git archive HEAD` in eine frische Kopie, daraus die Stufe `runtime` des `Dockerfile`. Gegenstelle: ein echter PostgreSQL 17 (das gepinnte Image aus `harness/mk/integration.mk`, `trust`) in einem eigenen Netz. Zwei echte Aufzeichnungen (je zwei Sessions mit drei Anfragen, einmal `--format yaml`, einmal `--format sqlite`) nahm das Binary mit `psql` auf; daraus erzeugte ein Python-Skript (`python:3.12-slim`, rohes `sqlite3` und Textersetzung, unabhängig vom Leser des Adapters) 93 Varianten. Jede lief mit `replay`, ein Teil auch mit `play`.
- **Unit-Mutanten:** je Mutant eine frische Kopie (`cp -r` ohne `-p` von `go.mod`, `go.sum`, `cmd/`, `internal/`, `test/`), genau eine Ersetzung (das Skript bricht ab, wenn der Suchtext nicht genau einmal trifft), `gofmt -l` leer, `go test -count=1` im Image der Stufe `deps`, `--network none`. Der Grundlauf des Pakets `recording` war dreimal grün.
- **Mutanten am Binary:** Kopie des Baums, Binary gebaut, `TestE2ENichtAbsteigendeZeitangaben` mit `PGR_BINARY` im Netz des PostgreSQL.
- Danach alles entfernt (Ende des Berichts).

---

## 1. DoD-Liefer-Punkte

### Punkt 1: absteigende Werte, Boundary, Form in YAML. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Test existiert | `TestZeitangabenAbsteigendBeschaedigt`, `TestZeitangabenGueltig`, `TestZeitangabenFormYAML`, `TestZeitangabenFormSQLite`, `TestZeitangabenVersionVorForm`, `TestE2ENichtAbsteigendeZeitangaben`, `TestUnmarshalFehlerReihenfolge`, `TestFelderDerVersion1`, `TestUnmarshalFehler`: jeder `^func <Name>(` einmal unter `internal/` oder `test/` | ja |
| ohne den Fix rot, aus dem richtigen Grund | Mutant R0: `yaml.go` auf den Stand `8e430ce`, die neuen Tests unverändert: rot sind `TestZeitangabenAbsteigendBeschaedigt`, `TestZeitangabenFormYAML`, `TestZeitangabenVersionVorForm`, `TestUnmarshalFehlerReihenfolge`; am Binary R0 mit der Datei aus `sqlite`: `replay` endet nicht binnen 15 s | ja |
| Vergleich je Format | A1 (Vergleich entfernt): `TestZeitangabenAbsteigendBeschaedigt`; am Binary: `replay` auf der SQLite-Datei endet nicht binnen 15 s | ja |
| gleicher Wert gültig | A2 (`<=`): `TestZeitangabenAbsteigendBeschaedigt`, `TestZeitangabenGueltig`, `TestZeitangabenFormYAML` | ja |
| je Session, nicht über Sessions | A3 (Vorgänger über Sessions behalten): `TestZeitangabenGueltig`, `TestZeitangabenAbsteigendBeschaedigt`, dazu `TestUnmarshalFehlerReihenfolge`, `TestZeitangabenFormYAML` | ja |
| Mischform | A4 (fehlender Wert setzt zurück): `TestZeitangabenAbsteigendBeschaedigt/über_eine_Interaktion_ohne_Wert` | ja |
| Meldung ohne Wert (Vergleich) | A5: `TestZeitangabenAbsteigendBeschaedigt` | ja |
| Form: Anführungszeichen, Typ | B1 (`Style` ungeprüft): `TestZeitangabenFormYAML`, `TestZeitangabenVersionVorForm`; am Binary B1: der Fall „yaml Text“ meldet die Bibliothek mit Wert und Zeile, Test rot | ja |
| Form: führende Null, Bereich, Vorzeichen, `null` | B2, B3, B4, B5: je `TestZeitangabenFormYAML` (B4 auch `TestFelderDerVersion1`, `TestUnmarshalFehlerReihenfolge`) | ja |
| Form: jede Session, Ort, Wert, in beiden Formaten | B6 (nur erste Session), B7 (Wert in der Meldung), B8 (Session und Interaktion vertauscht): `TestZeitangabenFormYAML`, B7 und B8 auch `TestZeitangabenFormSQLite` | ja |

Alle 13 Mutanten A1 bis B8 rot, jeder mit dem Test, der die Zusage trägt; Teilergebnisse zu den Mutanten der Version in Punkt 4.

### Punkt 2: die Prüfung gilt beim Laden, `replay` und `play` am Binary. **Bestätigt.**

Das gebaute Binary, 93 Dateien, `replay` mit Exit-Code und Meldung (alle ungültigen enden mit Exit 3, ohne dass `replay` anläuft; `play` verbindet sich nicht):

| Fall | Ergebnis |
|---|---|
| absteigend in YAML und SQLite (Session 1, Session 2, Mischform mit Lücke) | Exit 3, `PGR-E3003`, „Session S, Interaktion I: offset_ms kleiner als der vorige der Session“; die Werte stehen nicht in der Meldung (SQLite nennt zusätzlich den Pfad der Datei) |
| gleiche Werte, 0, größter `int64`, kleinerer Wert am Anfang der nächsten Session, Interaktion ohne Wert dazwischen, keine Werte | läuft an (`replay` wartet; `play` Exit 0) |
| Form: `1.0`, `1.5`, `1e3`, `0x10`, `0o17`, `1_5`, `"15"`, `'15'`, `!!int 15`, `true`, `null`, `~`, leer, ` `, Liste, Abbildung, `015`, `00`, `+15`, `-1`, `-0` | Exit 3, `PGR-E3003`, „Session 1, Interaktion 2: offset_ms ist keine kanonische Dezimalzahl“, ohne Wert und ohne Zeile |
| `9223372036854775808` und 23 Ziffern | Exit 3, `PGR-E3003`, „offset_ms liegt außerhalb des Bereichs“ |
| `0` und `15` an Stelle 2 (kanonisch) | `0` ist hinter 10 absteigend, `15` läuft an |
| `-5` an Stelle 1 | Exit 3, Form |
| SQLite: Text `'15'` und `15.0` in die Spalte | laden, weil SQLite sie durch die Affinität der Spalte in `15` wandelt (kein Befund); `1.5`: Exit 3, „weder NULL noch eine ganze Zahl“; `-5`: Exit 3 |
| echte Aufzeichnungen beider Formate (Roundtrip) | `replay` beantwortet `select 1`, `2`, `3` in zwei Verbindungen; `play` Exit 0 |

Ältere Beispiele: Das Repo trägt keine weiteren Aufzeichnungen in `yaml` oder `sqlite` außer der neuen Fixture (`git ls-files`); Handbuch und `README.md` zeigen keine Datei mit `offset_ms` in YAML. Keine Regression erkennbar.

**Mutanten am Binary:** `make test-integration`-Test `TestE2ENichtAbsteigendeZeitangaben` unverändert grün; mit A1 und R0 rot (der Aufruf von `replay` auf der SQLite-Datei endet nicht binnen 15 s), mit B1 rot (die Meldung nennt Wert und Zeile der Bibliothek).

### Punkt 3: Gate, Review, Closure, Register, Risiken, Paarungen. **Nicht Gegenstand dieses Laufs** (Closure-Pflichten, §7 vor der Closure unvollständig). `make gates`: siehe Verdikt.

---

## 2. Rangfolge der Prüfung, `version`, Formatkennung (`SPEC-001`, Entscheidung `ae1996d`)

Am Binary (`replay`, mit und ohne Inhaltsfehler `offset_ms: 1.5` in derselben Datei):

| Datei | Erwartet nach `SPEC-001` | Ergebnis |
|---|---|---|
| `version: 2` (auch mit Inhaltsfehler, mit unbekanntem Schlüssel `bogus`, mit absteigenden Werten) | `PGR-E3002` | `PGR-E3002`, „Version 2 wird nicht unterstützt“ |
| `0`, `-1`, `9`, `99999999999999999999` | `PGR-E3002` | `PGR-E3002`, mit der gelesenen Version |
| `abc`, `01`, `1.0`, `"1"`, `null`, `~`, leer, fehlend, `0x2`, `+1`, `-0` | `PGR-E3003` ohne Wert | `PGR-E3003`, „Schlüssel version fehlt oder ist keine kanonische Dezimalzahl“, auch mit Inhaltsfehler |
| `version: 1` mit Inhaltsfehler | `PGR-E3003` am Inhalt | Session und Interaktion |
| falsche oder fehlende Formatkennung, auch mit `version: 2` | `PGR-E3003` | `PGR-E3003`, „Formatkennung stimmt nicht“ |
| leere Datei, nur Kommentar, Liste an der Wurzel | `PGR-E3003` | `PGR-E3003` |

Der Code folgt der Spezifikation; „Version 0“ aus nicht vorhandener Version tritt nirgends mehr auf (F-637 behoben). Mutanten: C1 (Kopfprüfung aus) rot (`TestUnmarshalFehler`, `TestZeitangabenVersionVorForm`); C2 (`abc` als Version) rot; C3 (`Style` der Version ungeprüft) rot (`TestZeitangabenVersionVorForm`); C4 (`-0`) rot; C5 (Formatkennung nicht verglichen) rot (`TestUnmarshalFehler`, `TestZeitangabenVersionVorForm`); C6 (`0` und negativ als `PGR-E3003`) rot; **C7 (der Wert der Datei in der Meldung des `PGR-E3003` zur Version) bleibt grün: V-182.**

---

## 3. Spezifikation, Plan, Fixture, Nehmer

- **`LH-FA-21.a` *Prüfen* gegen den Code:** Wortlaut und Verhalten stimmen (kanonische Dezimalzahl, `int64`, `null` ist nicht „fehlt“, je Session, gleich gültig, Mischform, Meldung mit Session und Interaktion ohne Wert und Zeile, `sqlite` über `Unmarshal`). Der Satz zu mehreren Fehlern („nicht zugesagt“) bleibt akzeptiertes Negativ und ist durch `SPEC-001` auf den Inhalt eingegrenzt.
- **Plan folgt dem Code (§3.9):** §3 nennt `Unmarshal`, `formVorpruefung`, `offsetForm`, `kanonischeDezimalzahl`, `sessionFromDTO`, `geprueftFromDTO` (alle im Diff); §1 sagt, was der Slice liefert und abgibt; §6 trägt die Rangfolge und die Fixture (F-637 und F-639 eingearbeitet, F-638 behoben); §8 nennt zwei Liefer-Punkte, eine Schicht (Diff: nur `internal/adapters/driven/recording`, Tests, Abdeckungstabellen, Fixture; Kern, PGWire-Adapter, CLI, Bootstrap unverändert). Nur ein Schönheitsrest: V-184.
- **Fixture** `test/integration/testdata/zeitangaben-absteigend.sqlite`: Herkunft im Kommentar von `absteigendSQLite`; Werte 4711 und 4709, sonst gültig, wird aus dem richtigen Grund abgelehnt (die Mutanten A1 und R0 lassen `replay` laufen, ein gültiger Rest ist daher die einzige Abweichung). Die behauptete Probe des Schema-Wechsels (Spalte umbenannt, Test rot) habe ich nicht nachgestellt; der Kommentar sagt nur zu, dass die Meldung ausbleibt, und der Test prüft „Session 1, Interaktion 2“.
- **`--keep-timing`:** Das gebaute Binary lehnt die Option ab (`PGR-E2001`, „flag provided but not defined“); Plan §7 sagt es; kein Test und keine Zusage nennt sie als geprüft (der Name in `internal/adapters/driving/cli/leser_test.go` steht in einer Namensliste der Optionstabelle, nicht als Verhalten).
- **Nehmer (§3.13):** `slice-v1-abschluss-zeitangaben` nennt in §1 den Halbsatz „ein ungültiger Wert ist `PGR-E3003`“ als abgegeben an diesen Slice und schließt die Prüfung aus; `slice-v1-abschluss-zeitangaben-aufzeichnen` nennt die Prüfung mit der Kennung dieses Slice; dieser Slice trägt die Sendung in §1 und DoD Punkt 1. Die Adresse nimmt an.

---

## 4. Namensabgleich, Hard Rules, Sensoren

- **Namen:** jeder `Test…`-Name im Plan hat genau eine `^func`-Zeile; die Teilnamen (`Version_2,_unbekannter_Schlüssel` und andere) stehen als Fälle in den Tabellen.
- **Handbuch und README:** `git diff 8e430ce..HEAD -- README.md docs/user/benutzerhandbuch.md` leer.
- **Commit-Messages** `8e430ce..HEAD`: keine `SPEC-`/`ARC-`-Kennung ([ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md)). Alle nennen `slice-v1-abschluss-zeitangaben-pruefen` und `LH-FA-21`.
- **`//nolint`:** keines unter `cmd/`, `internal/`, `test/`.
- **`make abdeckung-check`:** grün (Exit 0); `docs/user/abdeckung-*.md` im Diff.
- **Hexagon:** `yaml.go` nutzt `strconv`, `strings`, `errors`; `internal/hexagon/**` unberührt ([ADR-0001](../plan/adr/0001-hexagonale-architektur.md)).
- **Größe:** 90 Zeilen Code in `yaml.go` (mit Kommentaren), 359 Zeilen Tests; ein Review der Sitzung tragbar.

---

## 5. Befunde

| ID | Kategorie | Befund | Pfad | Übergabe |
|---|---|---|---|---|
| V-182 | MEDIUM | `SPEC-001` *Rangfolge der Prüfung* sagt zu: Ist `version` fehlend oder keine kanonische Dezimalzahl, ist die Datei `PGR-E3003` „und die Meldung nennt keinen Wert“. Der Test `TestZeitangabenVersionVorForm` prüft das nur gegen das Wort „Version“ (Spalte `nicht`); der Wert der Datei (`abc`, `01`, `0x2`, …) darf in der Meldung stehen. Mutant C7 (an die Meldung von `kopfVorpruefung` den Wert `v.Value` angehängt: „… keine kanonische Dezimalzahl %s“) bleibt **grün** (`TestZeitangabenVersionVorForm`, `TestUnmarshalFehler`, `TestUnmarshalFehlerReihenfolge`, `TestFelderDerVersion1` alle PASS). Plan §7 führt für die Zusage „Wert in der Meldung (V11)“ einen roten Mutanten; der Mutant, der diesen Wert trägt, ist nicht der in der Spezifikation zugesagte. `AGENTS.md` §3.10: je Zusage eine Mutation, die rot wird. Die Zusage ist am Binary wahr (die Meldung nennt heute keinen Wert), nur nicht gefangen. | `internal/adapters/driven/recording/zeitangaben_pruefen_test.go` · `func TestZeitangabenVersionVorForm` (Spalte `nicht`); `internal/adapters/driven/recording/yaml.go` · `func kopfVorpruefung` | Implementer: je Fall den Wert der Datei (`abc`, `01`, `1.0`, `1`, `null`, `0x2`, `-0`) als nicht enthalten prüfen; Mutant C7 rot sehen; Plan §7 in der Zeile zu V11 nachziehen |
| V-183 | LOW | Einmal rot gesehen, nicht reproduziert: Im ersten Lauf des Mutanten C3 (nur `Style` der Version ungeprüft) war neben `TestZeitangabenVersionVorForm` auch `TestZeitangabenAbsteigendBeschaedigt` rot; in drei weiteren Läufen derselben Kopie (gleicher Aufruf), zweimal 12 Wiederholungen von `TestZeitangabenAbsteigendBeschaedigt` an Kopie C3 und an der unveränderten Kopie sowie drei Grundläufen des Pakets nie wieder. Die Ausgabe des ersten Laufs ist nicht aufgehoben (auf die ersten zwölf Zeilen gekürzt), die Ursache ist unbekannt. Ich behaupte weder Flakiness noch einen Fehler; der Test durchläuft die Formate über eine Map (`ladeBeide`), deren Reihenfolge nicht festliegt. Eintrag für das Beobachtungs-Register, falls es noch einmal auftritt. | `internal/adapters/driven/recording/zeitangaben_pruefen_test.go` · `func TestZeitangabenAbsteigendBeschaedigt` | Planner (Beobachtung), kein Nachzug |
| V-184 | INFO | §1 enthält nach dem Nachzug von F-638 einen Rest: „… liefert dieser Slice ( das Verhalten mit Zeitangaben — hier das Absteigen — ist der andere Liefer-Punkt, den jener Slice an `slice-v1-abschluss-zeitangaben` abgab).“ (Klammer mit führendem Leerzeichen; der Halbsatz liest sich, als gebe der Nehmer das Absteigen ab, das dieser Slice liefert). Bei einer Korrektur des Plans nebenbei lesbar machen. | Plan · §1 *Übernommen aus* | Implementer, beim nächsten Plan-Commit |
| V-185 | INFO | Nicht zugesagt und nicht geprüft: `version: !!int 1` (ausdrücklicher Typ an der Version) ist nach dem Code `PGR-E3003` (der Stil ist nicht 0), die Spezifikation nennt die Form nicht (sie nennt `"1"`). Gleiche Linie wie `offset_ms: !!int 15`, den die Spezifikation ausdrücklich nennt. Kein Handlungsbedarf, nur damit die Randform nicht als geprüft gilt. | `SPEC-001` *Rangfolge der Prüfung* | — |

Falsifikationsversuche, die nichts fanden: Anker (`&a`) im Wert (abgelehnt als Anker mit Zeile, Bestand, `BEO-REPO/yaml-meldung-nennt-wert-und-zeile` bleibt offen); doppelter Schlüssel `offset_ms`; SQLite-Text und -Real mit Affinität; Wurzel als Liste; leere Datei; Version über `int64`; `play` auf ungültiger Datei verbindet sich nicht (Exit 3 vor jeder Verbindung).

---

## Verdikt

**Befunde:** 0 HIGH, 1 MEDIUM (V-182), 1 LOW (V-183), 2 INFO (V-184, V-185).

Die beiden Liefer-Punkte der DoD sind **bestätigt**: Die Tests existieren, werden ohne den Fix (Mutant R0) aus dem richtigen Grund rot, und das gebaute Binary verhält sich in beiden Formaten für `replay` und `play` wie `LH-FA-21.a` und `SPEC-001` zusagen, einschließlich der Rangfolge und der Fälle der Version; 13 eigene Mutanten des Vergleichs und der Form rot, sechs der Version rot, drei am Binary rot (R0, A1, B1). Eine Zusage der Spezifikation (Meldung zur Version ohne Wert) ist nicht von einem roten Test getragen (V-182).

**Closure-blockierend:** V-182 (Hard Rule `AGENTS.md` §3.10): vor der Closure fängt ein Test den Wert in der Meldung, oder der Architect fasst die Zusage enger. V-183 bis V-185 blockieren nicht.

**Finding-Klasse dieses Laufs:** Zusage der Spezifikation weiter als der Test (die Meldung „ohne Wert“ gegen das Wort geprüft, nicht gegen den Wert) — dieselbe Klasse wie V-175 im Lauf `slice-v1-abschluss-sqlite-format-lesen`.

**Übergabe an den Planner:** dieser Bericht; V-182 an den Implementer; V-183 als Beobachtung; für die Closure: die Ausgänge der zwei Risiken aus §6 (Zwei Leser, zwei Kopien: *entfallen* — eine Stelle, `Unmarshal`, trägt beide Formate, belegt durch A1 und R0 an beiden Formaten; die Randform im Code: *eingetreten* mit F-637, erledigt durch `ae1996d` und `3cf1f57`) und die Klassen aus Review und Verifikation.

**Aufgeräumt:** Container `ver-zp-pg` und die Läufe `ver-zp-run-*`, Netz `ver-zp-net`, Images `ver-zp-bin`, `ver-zp-deps`, die Verzeichnisse und Skripte mit dem Präfix `mut`, `src`, `p`, `rec` im Scratchpad. Der Arbeitsbaum des Repos war bis zu diesem Bericht unverändert.

**Gate:** `make gates` auf dem sauberen Baum nach dem Commit dieses Berichts: Ergebnis in der Übergabe-Nachricht (ein Gate-Lauf belegt nur den Stand, auf dem er lief).
