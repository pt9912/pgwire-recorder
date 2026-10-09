# Review-Report: slice-v1-abschluss-schreiben — 2026-10-09

**Review-Art:** Code-Review. Geprüft wird gegen Plan §1, §3 und §6, gegen die Entscheidungen in [`LH-FA-07.a`](../../spec/spezifikation.md#lh-fa-07a--sicheres-schreiben-des-recordings) (*Zielpfad beim Start*, *Temporäre Datei*) und gegen die Hard Rules (Modul 10). Die Entscheidungen sind die Vorab-Prüfung des Architect `ad1b7ba` und seine Entscheidung der zwei Rückgaben `b74e16e`. Hinzu kommt [ADR-0005](../plan/adr/0005-recording-store-ist-driven-adapter.md). Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** `ad1b7ba..9bd65c5`. `ad1b7ba` selbst habe ich als Grundlage gelesen (`git show ad1b7ba`).

- `aad7085`: Zielpfad, der keine reguläre Datei ist, ist `PGR-E3001`. Die temporäre Datei wird nach einem Fehlschlag entfernt. Unit- und E2E-Tests, Abdeckung.
- `df4f591`: Handbuch §5 und Tabelle der Meldungscodes.
- `65968e6`: Probedatei über die Tabelle `dateiOps`, Tests zu Name und nicht beschreibbarem Verzeichnis.
- `0311f72`: §7 *Belege des Implementers*, mit zwei Rückgaben an den Architect.
- `b74e16e` (Architect): Setzen der Rechte und Schließen zählen zum *Fehlschlag*. Rechte bei einer ersetzten Verknüpfung.
- `24fc87c`: `Rechte` und `Schliessen` in `dateiOps` und `Eingriffe`, `TestWriteRechteVerknuepfung`.
- `9bd65c5`: §7 zur Nacharbeit.

Der Diff hat 970 Zeilen hinzu und 37 entfernt. Davon sind 78/25 Code in `yaml.go`, 507 Unit-Tests, 167 Integration und 161 Plan.

**Skill:** `.harness/skills/reviewer.md` @ `32afc69`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-09

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-v1-abschluss-schreiben.md`: ganz gelesen am Stand `9bd65c5`, mit §7 *Belege des Implementers*. Den Bericht des Implementers habe ich nicht gelesen.
- `spec/spezifikation.md`: [`LH-FA-07.a`](../../spec/spezifikation.md#lh-fa-07a--sicheres-schreiben-des-recordings) am Stand `9bd65c5`, dazu [`LH-FA-13.b`](../../spec/spezifikation.md#lh-fa-13b--fehlerebenen-und-prozessstatus) (Vorrang von Exit-Code 3), `SPEC-033`, `SPEC-034` und `SPEC-038` *Warten in Tests*.
- [`LH-FA-07`](../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-08`](../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-17`](../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration).
- [ADR-0005](../plan/adr/0005-recording-store-ist-driven-adapter.md).
- `AGENTS.md` §3.3, §3.7 und §3.9 bis §3.13.
- `slice-v1-abschluss-sqlite-format`: §1 und DoD, als Nehmer der Abgrenzung in §6 (§3.13).
- Vorherige Reports am selben Lauf der Welle: F-509, F-514 und F-526 als Vergleich für die Einstufung. Die Nummern dieses Laufs beginnen bei F-537.

**Ausgeführte Läufe:**

- **Arbeitsweise.** Je Lauf eine frische Kopie von `go.mod`, `go.sum`, `cmd`, `internal` und `test` per `cp -r` ohne `-p` unter dem Scratch-Pfad, als `rs-mut-<X>`, `rs-it-<X>` bzw. `rs-sonde-<X>`. Unit-Läufe: `go test -count=1 ./internal/adapters/driven/recording/` im eigenen Image `pgr-rev-schreiben:deps` (Stufe `deps`), `--network=none`, Bind-Mount der Kopie. Ein Skript ersetzt genau eine Stelle und bricht bei einer anderen Trefferzahl ab. Im Repo habe ich nichts geändert. Gelöscht habe ich nur die benannten Kopien und Images.
- **Unveränderte Kopie:** Unit grün. Die vier neuen E2E-Tests sind grün (`TestE2ERecordVorhandeneZieldatei` mit sieben Fällen, `…ZielIstVerzeichnis`, `…ZwangsendeSchreibtVollstaendig`, `…SchreibfehlerNachZwangsende`). Gelaufen sind sie in der Stufe `integration` mit eigenem Tag, eigenem internem Netz `pgr-rev-schreiben-<X>-net` und eigenem PostgreSQL-Container (Image aus `integration.mk`). Ein `trap` entfernt Container, Netz und Image.
- **Mutanten:**
  - A: Prüfung *keine reguläre Datei* nur mit `replace`. rot (`TestPrepareKeineRegulaereDatei`, drei Fälle `replace=false`)
  - D: Probedatei nur, wenn der Pfad fehlt (`case err == nil: return nil`). **grün** (F-537)
  - C: jeder Fehler beim Anlegen der temporären Datei führt zu einem neuen Zug, nicht nur `fs.ErrExist`. **grün** (F-544)
  - P: Fehler beim Entfernen der Probedatei verworfen. **grün** (F-543)
  - M2 (Integration, Gegenprobe zu §7): Entfernen im Zweig des Verschiebens weggelassen. rot (`TestE2ERecordSchreibfehlerNachZwangsende`, zwei `.tmp` im Verzeichnis), wie §7 sagt.
- **Sonde am Binary** (Stufe `integration`, `--user 65534`, `--network=none`): Eine vorhandene Datei liegt in einem Verzeichnis mit `0555`, gestartet wird `record … --force`. Das Original endet sofort mit `Recording [PGR-E3001]: Verzeichnis von /tmp/ro/rec.yaml nicht beschreibbar`, Exit 3. Mutant D startet, nimmt an, schreibt `Herunterfahren begonnen` und scheitert erst beim Schreiben am Ende (`temporäre Datei … nicht anzulegen`).
- `make a-check`: Exit 0, `gesamt: 0 Befund(e)`. `make kopf-check` und `make abdeckung-check`: Exit 0. `make docs-check` vor dem Commit.
- Nicht gefahren: `make gates`, `make lint`, die volle `make test-integration`. Von den Mutanten aus §7 habe ich nur M2 wiederholt; A fällt mit der Zeile *auch mit `--force`* zusammen.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-537 | MEDIUM | Die Zusage *„lässt sich darin keine Datei anlegen, ist das `PGR-E3001`“* beim Start hält eine Mutation nur für einen Pfad, der noch fehlt. Mutant D überspringt die Probedatei für eine vorhandene Datei, und alle Unit-Tests bleiben grün. `TestPrepareVerzeichnisNichtBeschreibbar` nutzt nur einen fehlenden Pfad mit `replace=false`. Kein E2E-Test fährt `--force` in einem Verzeichnis ohne Schreibrecht. *Failure-Szenario:* Die Sonde zeigt es am Binary. Mit `--force` über einer vorhandenen Datei in einem schreibgeschützten Verzeichnis startet der Mutant, nimmt Verbindungen an und scheitert erst beim Schreiben am Ende. Das widerspricht *„Alle diese Prüfungen laufen beim Start“*, und kein Gate wird rot. Die Klasse war in den vorigen Reports mehrfach LOW und zuletzt MEDIUM (F-534). | [`LH-FA-07.a`](../../spec/spezifikation.md#lh-fa-07a--sicheres-schreiben-des-recordings) *Zielpfad beim Start*, *Verzeichnis*; `AGENTS.md` §3.10 | `internal/adapters/driven/recording/yaml.go` · `case err == nil, errors.Is(err, fs.ErrNotExist):` | ja (Mutant D, Sonde, `go test`) | Zusage nur für einen Teil ihrer Fälle von einer Mutation gehalten |
| F-538 | MEDIUM | §6 *Abgrenzung zu `sqlite`* weist die Entscheidung, was von *Temporäre Datei* für `sqlite` gilt, `slice-v1-abschluss-sqlite-format` zu (*„entscheidet … vor seinem Code“*). Gesetzt hat das `ad1b7ba`. §1 und DoD des Nehmers führen die Sendung nicht. `grep slice-v1-abschluss-schreiben` in dessen Plan findet 0 Treffer, und sein §1 schließt nur Formatumwandlung und YAML aus. §7 behauptet trotzdem *„Kein anderer Slice ist als neue Adresse genannt (§3.13)“*. *Failure-Szenario:* `slice-v1-abschluss-sqlite-format` beginnt mit einem §6, das weder Zielpfad noch temporäre Datei nennt. Der Implementer entscheidet sie dann im Code, und das Review findet sie. Das ist genau das Muster aus `BEO-REPO/spec-randform-erst-im-review-entschieden`. | `AGENTS.md` §3.13, §3.12; `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` | Plan §6 · „`slice-v1-abschluss-sqlite-format` vor seinem Code (akzeptiertes Negativ hier“; Plan §7 · „Kein anderer Slice ist als neue Adresse genannt (§3.13).“ | ja (`grep` im Nehmer) | Adresse nimmt nicht an |
| F-539 | MEDIUM | Der Code-Commit `aad7085` entschied, dass auch nach einem Fehlschlag beim Setzen der Rechte (`merr`) und beim Schließen (`cerr`) die temporäre Datei entfernt wird. Neu ist das Entfernen im Zweig `errors.Join(werr, merr, serr, cerr)`. Damals nannte §6 nur *„Fehlschlag von Anlegen, Schreiben, Synchronisieren oder Verschieben“*. Zurückgegeben wurde die Frage erst danach, in `0311f72`, und der Architect bestätigte sie in `b74e16e`. Das Ergebnis ist jetzt entschieden und getestet. Der Weg ist aber der, den §3.12 ausschließt. *Failure-Szenario:* Hätte der Architect das Setzen der Rechte nicht zum Fehlschlag gezählt, wäre ein gelieferter und geprüfter Stand zurückzunehmen gewesen. Gleiche Klasse wie F-509. | `AGENTS.md` §3.12; `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` | `git show b74e16e` · „-- **Fehlschlag von Anlegen, Schreiben, Synchronisieren oder Verschieben**“; `yaml.go` · `errors.Join(err, ops.entfernen(tmp.Name())), "temporäre Datei für %s nicht zu schreiben"` | nein (Historie) | Randform im Code entschieden, danach zurückgegeben |
| F-540 | MEDIUM | §1, Schicht-Abgrenzung, sagt weiter *„Der Slice ändert den Recording-Adapter und den CLI-Adapter“*. §3 sagt seit `aad7085` für `internal/adapters/driving/cli` *„keine Änderung“*, und §7 nennt den CLI-Adapter unverändert. `aad7085` zog §3 nach, §1 nicht. *Failure-Szenario:* Die Schicht-Abgrenzung ist die Selbstbindung, die das Review prüft. Eine spätere Änderung am CLI-Adapter, etwa an der Anmeldung von `--force`, die §1 an `slice-v1-abschluss-konfiguration` abgibt, ginge gegen diese Zeile als im Umfang durch. Die Klasse ist zweimal LOW (F-514, F-526), deshalb MEDIUM. | `AGENTS.md` §3.9; `BEO-REPO/plan-folgt-korrektur-nicht` | Plan §1 · „den Recording-Adapter und den CLI-Adapter.“ (Z. 64) | nein (Lesen) | Plan-Zeile folgt dem Liefer-Commit nicht |
| F-541 | LOW | §7 sagt *„Kein Code-Commit ändert §6.“* Der Code-Commit `24fc87c` ändert aber eine Zeile in §6 (*„Der Code folgt schon (`os.Stat`); es fehlt der Test (§3).“* → *„… der Test ist `TestWriteRechteVerknuepfung`“*). Eine neue Randform kommt damit nicht hinzu, deshalb ist das kein Befund *Randform im Code-Commit*. Der Satz in §7 sagt aber mehr zu, als der Diff hält. | `AGENTS.md` §3.11 | Plan §7 · „**Randformen.** Kein Code-Commit ändert §6.“ | ja (`git show 24fc87c -- docs/plan`) | Zusage in der Plan-Zeile weiter als Prüfung |
| F-542 | LOW | Der neu geschriebene Kommentar an `Write` belegt die Rechte-Regel mit *„(SPEC-033)“*. `SPEC-033` ist *Sicherheit* und sagt nur *„Dateirechte … liegen in der Verantwortung des Anwenders“*. Die Regel *0666 nach der umask, ersetzt mit ihren Zugriffsrechten* steht in `LH-FA-07.a` *Temporäre Datei*, *Rechte*. Der Rang-Zeiger führt also an eine Stelle, die die Regel nicht trägt. Den gleichen Zeiger trägt `yaml_test.go` vor dem Slice; er liegt nicht im Diff. | `AGENTS.md` §3.7 (Rang-Zeiger) | `internal/adapters/driven/recording/yaml.go` · „Zugriffsrechte (SPEC-033).“ | nein (Lesen) | Rang-Zeiger auf eine Stelle, die die Regel nicht trägt |
| F-543 | LOW | Der neue Kommentar an `Prepare` sagt zu: *„legt … eine Probedatei an und entfernt sie sofort; gelingt das nicht, ist das PGR-E3001“*. Das schließt Schließen und Entfernen der Probedatei ein. Diese Hälfte nennen weder `LH-FA-07.a` (*„die sofort wieder entfernt wird“*) noch §6, und kein Test prüft sie: Mutant P verwirft den Fehler beim Entfernen und bleibt grün. Das Verhalten liegt vor dem Slice (`closeErr`, `removeErr`); neu ist die Zusage im Kommentar. | `AGENTS.md` §3.11, §3.12 | `internal/adapters/driven/recording/yaml.go` · „sofort; gelingt das nicht, ist das PGR-E3001.“ | ja (Mutant P) | Kommentar sagt zu, was kein Test prüft |
| F-544 | INFO | Mutant C bleibt grün: Jeder Fehler beim Anlegen der temporären Datei führt zu einem neuen Zug, nicht nur ein belegter Name. Der Code bleibt `PGR-E3001`. Unter dem Mutanten meldet ein zu langer Name aber *„kein freier Name … nach zehn Versuchen“* ohne die Ursache des Betriebssystems. `LH-FA-07.a` sagt nur *„ein vorhandener Name wird nie überschrieben, sondern ein neuer gezogen“*, und §7 führt keine Zusage, dass ein anderer Fehler sofort aufgibt. Eine Zusage ist deshalb nicht verletzt. Hinweis an den Architect, ob das Wort *vorhanden* diese Unterscheidung tragen soll. | `LH-FA-07.a` *Temporäre Datei*, *Name* | `internal/adapters/driven/recording/yaml.go` · `if errors.Is(err, fs.ErrExist) {` | ja (Mutant C) | — (Hinweis an den Architect) |
| F-545 | INFO | Die Zusammenfassung *Fehlermodi* unter `LH-FA-07.a` nennt für `PGR-E3001` *„Schreiben oder Verschieben gescheitert“*. Den ausführlichen Absatz *Fehlschlag* hat `b74e16e` um Setzen der Rechte und Schließen erweitert; sie nennen weder diesen noch das Anlegen und die zehn belegten Namen. Die Detailregel darüber ist vollständig, die Zusammenfassung ist enger. Hinweis an den Architect. | `LH-FA-07.a` *Fehlermodi* | `spec/spezifikation.md` · „Schreiben oder Verschieben gescheitert → Exit-Code `3`“ | nein (Lesen) | — (Hinweis an den Architect) |
| F-546 | INFO | Größe: Der Diff war in einer Sitzung prüfbar. Der Code umfasst 78 Zeilen in `yaml.go` mit einer kleinen Tabelle von Operationen. Die Tests sind lang, aber je Zusage einzeln lesbar, und §7 belegt sie Zeile für Zeile. Das Risiko aus §6 (*plattformabhängig atomar*) kann das Review nicht beurteilen. §7 *Grenzen* 2 nennt es offen. | Baseline `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice; Plan §6 *Risiken* | Plan §7 · „in eine Review-Sitzung passt, urteilt das Review.“ | nein | — (Hinweis an den Verifier) |

---

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driven/recording/yaml.go`, Prüfung beim Start | geprüft. Befunde F-537 und F-543. Die Reihenfolge stimmt mit *Zielpfad beim Start* überein: zuerst keine reguläre Datei (`PGR-E3001`, auch mit `replace`), dann vorhanden ohne `replace` (`PGR-E2002`), dann nicht prüfbar (`PGR-E3001`), zuletzt die Probedatei im Verzeichnis der Zieldatei, ohne `MkdirAll`. `os.Stat` folgt der Verknüpfung, eine Verknüpfung ins Leere fällt unter `fs.ErrNotExist`, und eine Schleife ist nicht prüfbar. Mutant A ist rot. |
| `yaml.go`, Schreiben | geprüft. Befunde F-539, F-542 und F-544. Der Name `.<Name>.<16 Hexziffern>.tmp` wird mit `O_EXCL` und `0666` angelegt, höchstens zehn Versuche. Rechte übernimmt der Code nur, wenn `os.Stat` gelingt (Grenze aus §6, akzeptiertes Negativ). Jeder Fehler nach dem Anlegen entfernt die temporäre Datei, und der Fehler beim Entfernen hängt per `errors.Join` an derselben Meldung `PGR-E3001` (`SPEC-034`). Das Verschieben mit `os.Rename` ersetzt die Verknüpfung. Ein Pfad, der nach dem Start entsteht, wird ohne Prüfung ersetzt. Eine Randform, die §6 heute nicht nennt, entscheidet der Code darüber hinaus nicht; die Probedatei ist F-543. |
| `export_test.go`, `Eingriffe` | geprüft, ohne Befund. `Eingriffe` liegt nur im Testbinary. Ein Feld mit `nil` ist die Operation des Betriebssystems, wie der Kommentar sagt. `PruefeMit` und `SchreibeMit` rufen dieselben `pruefe` und `schreibe` wie `Prepare` und `Write`. |
| [ADR-0005](../plan/adr/0005-recording-store-ist-driven-adapter.md), Hexagon | geprüft, ohne Befund. Geändert ist nur der Recording-Adapter. Kern, PGWire-Adapter, CLI-Adapter und Bootstrap liegen nicht im Diff. `make a-check`: Exit 0, `0 Befund(e)`. |
| Unit-Tests `schreiben_test.go` | geprüft. Befunde F-537 und F-543. Die Kommentare der Abdeckung sagen nicht mehr zu, als die Fälle prüfen. `TestWriteFehlschlag` legt die Zieldatei in jedem Fall an, also läuft auch der Fall `Rechte` (Kommentar im Fall). `TestWriteRechteVerknuepfung` unterscheidet mit umask `077` und Ziel `0640` beide Mutanten: `os.Lstat` ergäbe `0777`, ohne `Chmod` `0600`. `syscall.Umask` ist prozessweit; im Paket nutzt kein Test `t.Parallel`. |
| E2E `schreiben_e2e_test.go`, `startProzessIn`; `SPEC-038` *Warten in Tests* | geprüft, ohne Befund. Gewartet wird auf das Prozessende über `beendet`, je Prozess an einer Stelle (`warteEnde`) und mit Literalfristen (15 s, 30 s, Bereitschaft 10 s), höchstens 60 s. `starteBis` hat 15 s. `t.Setenv` im Fall `Umgebung` wirkt auf den gestarteten Prozess. Die Datei-Quelle liegt im Arbeitsverzeichnis `dir` des Prozesses. Exit 3 nach dem Zwangsende folgt dem Vorrang in [`LH-FA-13.b`](../../spec/spezifikation.md#lh-fa-13b--fehlerebenen-und-prozessstatus). |
| Handbuch §5 und Tabelle `PGR-E3000`/`PGR-E3001` | geprüft, ohne Befund. Jede neue Aussage hat einen Test: die Quellen von `--force`, keine reguläre Datei auch mit `--force`, kein angelegtes Verzeichnis, Verknüpfung nach ihrem Ziel und ersetzt, Name der temporären Datei, Entfernen nach Fehlschlag und Meldung des Fehlers beim Entfernen, keine Behandlung übrig gebliebener Dateien, Rechte und Rechte bei Verknüpfung. *SIGKILL* steht als Grenze. Beispiele sind nicht geändert. |
| Abdeckung `docs/user/abdeckung-*.md` | geprüft, ohne Befund. Die Zeilen sind gleich den Deklarationen der neuen Tests. `make abdeckung-check`: Exit 0. |
| `spec/spezifikation.md` (`ad1b7ba`, `b74e16e`) | geprüft. Befund F-545. Keine ADR, kein Slice, keine Welle und kein Commit wird im Spec-Stratum genannt. Die Änderungsliste trägt je eine Zeile. |
| Plan §1, §3, §6; `AGENTS.md` §3.9, §3.12, §3.13 | geprüft. Befunde F-538 bis F-541. Der Kopf führt die Kennungen aus §1 und §2 (`make kopf-check`: Exit 0). |
| Hard Rule 3.3, Commit-Messages | geprüft, ohne Befund. Im Diff gibt es keinen Move. Alle acht Messages nennen `slice-v1-abschluss-schreiben` und `LH-FA-07`, keine nennt eine `SPEC-` oder `ARC-`Kennung. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 3 |
| INFO | 3 |

**Summary:** 0 HIGH · 4 MEDIUM · 3 LOW · 3 INFO.

- F-537: Ein nicht beschreibbares Verzeichnis bei vorhandener Datei mit `--force` hält keine Mutation (Mutant D grün, Sonde am Binary).
- F-538: §6 weist die Entscheidung für `sqlite` einem Slice zu, der sie nicht führt; §7 sagt das Gegenteil.
- F-539: Entfernen nach Fehlschlag von Setzen der Rechte und Schließen ist im Code entschieden und erst danach zurückgegeben.
- F-540: §1 nennt den CLI-Adapter als geändert, §3 und der Diff nicht.
- F-541: §7 *„Kein Code-Commit ändert §6“* trifft für `24fc87c` nicht zu.
- F-542: Rang-Zeiger `SPEC-033` an der Rechte-Regel.
- F-543: Der Kommentar an `Prepare` sagt `PGR-E3001` für das Entfernen der Probedatei zu; das steht weder in der Spezifikation noch in einem Test.
- F-544: Nur ein belegter Name führt zu einem neuen Zug; das hält keine Mutation, und keine Zusage verlangt es.
- F-545: Die Zusammenfassung *Fehlermodi* ist enger als *Fehlschlag*.
- F-546: Der Diff war in einer Sitzung prüfbar.

Wiederkehrende Klassen: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (F-537), `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (F-538), `BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben` (F-539), `BEO-REPO/plan-folgt-korrektur-nicht` (F-540), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (F-541, F-543).

**Finding-Klassen dieses Laufs:** Zusage nur für einen Teil ihrer Fälle von einer Mutation gehalten · Adresse nimmt nicht an · Randform im Code entschieden, danach zurückgegeben · Plan-Zeile folgt dem Liefer-Commit nicht · Zusage in der Plan-Zeile weiter als Prüfung · Rang-Zeiger auf eine Stelle, die die Regel nicht trägt · Kommentar sagt zu, was kein Test prüft

## Verdikt

**Merge-blockierend:** ja, bis F-537, F-538 und F-540 beantwortet sind. Das gelieferte Verhalten hält am Code und an den Tests: Prüfung beim Start, Verknüpfungen, Name und Exklusivität der temporären Datei, Rechte, Entfernen nach jedem Fehlschlag mit dem Fehler beim Entfernen als Ursache, die drei Quellen von `--force` und das Zwangsende. Von vier eigenen Unit-Mutanten ist einer rot (A). Drei bleiben grün: D (F-537), C (F-544, ohne verletzte Zusage) und P (F-543). Die Gegenprobe M2 ist rot, wie §7 sagt. F-539 verlangt keine Änderung am Code; das Ergebnis ist entschieden und getestet. Es ist ein Beleg für die Closure.

**Übergabe:** F-537, F-540, F-541, F-542 und F-543 gehen an den Implementer. F-543 geht wegen der Randform (Entfernen der Probedatei scheitert) zuerst an den Architect, und der Implementer entscheidet sie nicht. F-538 geht an den Planner und Architect, der die Abgrenzung in `ad1b7ba` gesetzt hat (§3.13: Der Nehmer führt die Sendung, sonst ist er keine Adresse). Die Zeile in §7 geht an den Implementer. F-539 geht in die Closure §7 und ins Register. F-544 und F-545 gehen an den Architect, F-546 an den Verifier. Dieser Report ersetzt keine Verifikation.
