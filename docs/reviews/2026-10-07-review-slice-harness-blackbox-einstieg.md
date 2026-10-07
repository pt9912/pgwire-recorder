# Review-Report: slice-harness-blackbox-einstieg — 2026-10-07

**Review-Art:** Code-Review (Testumbau, kein Produkt-Code) gegen Plan, Entscheidungen (`SPEC-049` Punkt 5, 7 in der Fassung `f0ea7be` und 8, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 4 und 5) und Hard Rules (Modul 10). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git show 1e4381b`. Die Testdateien unter `internal/adapters/driving/cli`, `internal/bootstrap` und `test/integration` gehören jetzt zu `cli_test`, `bootstrap_test` und `integration_test`. Neu sind die Brücken `internal/adapters/driving/cli/export_test.go` und `internal/bootstrap/export_test.go`. Produkt-Code ist nicht geändert. Schwerpunkte laut Auftrag: Brücken nach Punkt 7, `errcheck` an neun `Close`-Stellen, `lebendTexte`/`vtTexte`, die geteilte SA4000-Bedingung, Prüfung im Umbau verloren, die grünen Mutanten nach Schritt 19 und eigene Mutationen.

**Skill:** `.harness/skills/reviewer.md` (Stand `de770f7`, unverändert bis `1e4381b`)
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-07

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-harness-blackbox-einstieg.md` (ganz, Stand `1e4381b`) und sein Diff in `1e4381b`. §7 *Belege des Implementers* habe ich gelesen, den Bericht des Implementers nicht.
- `spec/spezifikation.md` `SPEC-049` Punkt 1 bis 9 (seit `f0ea7be` unverändert, `git diff f0ea7be 1e4381b -- spec/` leer); [`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes)
- [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) Entscheidung 2, 4 und 5; [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md)
- `AGENTS.md` §3.2, §3.3, §3.6, §3.7, §3.9 bis §3.12, §5 Regel 1; `.claude/commands/implement-slice.md` Schritt 19
- Produkt-Code `internal/adapters/driving/cli/cli.go` (Konstanten, Hilfetext, Lesen der Umgebung), `internal/bootstrap/bootstrap.go` (`Run`, `logger`, `fail`), `internal/adapters/driven/recording/yaml.go` (`Unmarshal`-Pfad, `formVorpruefung`, `fromDTO`); `.golangci.yml` (`errcheck.exclude-functions`)
- Testdateien an `1e4381b^` und `1e4381b`, je Datei der ganze Diff
- Vorherige Reports: Reviews zu `slice-harness-blackbox-kern`, `slice-harness-blackbox-driven`, `slice-harness-blackbox-pgwire`; für die Klasse „Struktur-ID in der Commit-Message“ F-336. Die höchste vergebene Nummer vor diesem Lauf war F-453.

**Ausgeführte Läufe:**

- `make lint` im Repo: Exit 2, `32 issues:`. Unter den Pfaden dieses Slice stehen genau zwei Befunde, `internal/bootstrap/bootstrap.go:71:26` (`contextcheck`) und `internal/adapters/driving/cli/cli.go:243:7` (`revive` `unused-receiver`), beide Produkt-Code; keine `*_test.go` und keine `lint:`-Zeile in der Ausgabe.
- `make abdeckung-check` und `make kopf-check`: beide Exit 0.
- `1e4381b^` und `1e4381b` per `git archive` in je eine frische Kopie unter dem Scratch-Pfad; Läufe in einem eigenen Image der Stufe `deps` (eigener Tag), ohne Netz, per Bind-Mount. `go test -list .` für `cli` und `bootstrap`, für `./test/integration` mit `-tags integration` und `PGR_OHNE_POSTGRES=1`: vorher und nachher 70 Tests, die Listen gleich.
- In der Kopie von `1e4381b`: `gofmt -l` leer, `go vet -tags integration` der drei Pakete grün, `go test` von `cli` und `bootstrap` grün.
- Integrationsläufe in einem eigenen internen Docker-Netz mit eigener PostgreSQL-Instanz (das gepinnte Image aus `harness/mk/integration.mk`) und eigenem Volume, Testimage der Stufe `integration` aus der jeweiligen Kopie unter eigenem Tag. Ohne Mutation (E0): 37 bestanden, 2 übersprungen (die Tests der zweiten Phase), Exit 0. Die sechs Tests mit geprüftem `Close`, deren Gegenseite die Verbindung selbst beenden kann (`BeendenMitOffenerVerbindung`, `NichtVerbrauchtHerunterfahren`, `SigtermMittenInFolge`, `SigtermBeimPipelining`, `LebendpruefungPgxpool`, `LebendpruefungDatabaseSQL`), mit `-test.count=5`: 30 von 30 bestanden. Danach Container, Netz, Volume und die eigenen Images entfernt; nichts hostweit aufgeräumt.
- Vier eigene Mutationen, je in einer frischen Kopie (`cp -r` ohne `-p`, Ankunft per Prüfsumme kontrolliert), Tabelle unten. Das Repo blieb bis auf diese Datei unverändert.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-454 | LOW | Die Commit-Message von `1e4381b` nennt eine Struktur-ID: „Unexportiertes erreichen die Unit-Tests nur über export_test.go nach SPEC-049 Punkt 7“. Struktur-IDs adressieren innerhalb der Spec und gehören nicht in die Commit-Message. Zweites Auftreten der Klasse nach F-336. | `AGENTS.md` §5 Regel 1 | Commit `1e4381b` (Message, Rumpf Zeile 4) · „nach SPEC-049“ | ja (`git show -s 1e4381b`) | Struktur-ID in der Commit-Message |
| F-455 | INFO | Die White-Box-Liste in §6 stand vor dem Code nur als Namenssuche (`envFailOnUnconsumed`, dazu `replay` als Fehltreffer) und ist erst im Code-Commit per AST gemessen: vier Zugriffe, darunter drei neue (`envLogLevel`, `logger`, `fail`). Eine Randform ist damit nicht im Code entschieden: Alle vier fallen in die Klassen, die der Punkt *Export-Test-Brücke* vor dem Code entschieden hat (Konstanten und weiterreichende Funktionen mit Ein- und Ausgaben exportierter Typen, kein Wert eines unexportierten Typs), deshalb kein MEDIUM nach `AGENTS.md` §3.12. Dasselbe Muster, eine im Code-Commit nachgezogene Zugriffsliste, steht auch in Reviews früherer Umstellungs-Slices (`slice-harness-blackbox-kern`: „Plan nachgezogen“ an *White-Box-Zugriffe*; `slice-harness-blackbox-pgwire`: Messung erst am Code-Stand). Für die Closure (Steering Loop). | `AGENTS.md` §3.12; Plan §4 *Start* (Architect prüft §6 vor dem Code) | `docs/plan/planning/in-progress/slice-harness-blackbox-einstieg.md` · „per AST gemessen am Stand `f28a2df` (§7)“ | nein (Lesen) | White-Box-Liste vor dem Code nur per Namenssuche |
| F-456 | INFO | Der grüne Mutant `envFailOnUnconsumed` ist richtig eingeordnet: In `cli` und `bootstrap` grün (R1), in `TestE2EReplayNichtVerbraucht` rot (R1i, „Exit-Code 0, erwartet 5“); vorher nannten die Tests in `cli` dieselbe Konstante direkt, ein Verlust durch den Umbau ist es nicht. Asymmetrie aus dem Bestand: `TestParseLogLevelHilfe` hält den Namen von `--log-level` über den Hilfetext fest (R2 rot), für `--fail-on-unconsumed` gibt es in `cli` keinen solchen Vergleich, obwohl der Hilfetext den Namen als Literal führt. Ein solcher Test wäre ein neuer Fall, den §1 ausschließt. | Plan §6 *Mutationstests auf unexportierte Teile*; `.claude/commands/implement-slice.md` Schritt 19; Rolle Verifier | `internal/adapters/driving/cli/cli.go` · „const envFailOnUnconsumed = \"PGWIRE_RECORDER_FAIL_ON_UNCONSUMED\"“ | ja (R1, R1i, R2) | — (Hinweis) |
| F-457 | INFO | Die Einordnung des grünen PGR-E3003-Platzes stimmt (R3): Die Eingabe von `TestE2EReplayBeschaedigt` (`format`, `version`, keine `sessions`) läuft durch `formVorpruefung` ohne Fehler, weil `wert` für fehlende `sessions` einen leeren Knoten liefert, und scheitert erst in `fromDTO`; die Rückgabe nach `formVorpruefung` fangen `TestUnmarshalExtendedFehler` und `TestVorpruefungNenntOrt`. Die Formulierung „Der erste Platz derselben Mutation“ ist ungenau: Gemutet ist dort die Rückgabe „Aufzeichnung beschädigt“ nach `formVorpruefung`, nicht die Meldung „Liste der Sessions fehlt“ aus der Tabellenzeile. | `.claude/commands/implement-slice.md` Schritt 19 | `docs/plan/planning/in-progress/slice-harness-blackbox-einstieg.md` · „Der erste Platz derselben Mutation von PGR-E3003 nach PGR-E3001 (`formVorpruefung`)“ | ja (R3) | — (Hinweis) |
| F-458 | INFO | Die neun `Close`-Prüfungen verdecken nichts: Der Fehler wurde vorher verworfen und wird jetzt mit `t.Error` gemeldet. An allen neun Stellen steht `defer cancel()` vor dem `defer func()` mit `Close`, nach LIFO schließt die Verbindung also wie vorher vor dem `cancel` und am Ende des Helfers; `t.Cleanup` hätte das Schließen hinter `rec.stop`/`rep.stop` des Aufrufers verschoben. 30 Läufe der Tests, in denen die Gegenseite die Verbindung beenden kann, lieferten keinen `Close`-Fehler. Aus dem Bestand: In `TestE2ERecordExtendedSigtermBeimPipelining` ist die Goroutine, die `pc` benutzt, beim `Close` nicht beendet; ein Fehler aus dieser Überschneidung würde jetzt sichtbar, beobachtet habe ich keinen, und ohne `-race` (Alpine, `CGO_ENABLED=0`) bleibt das Lesen. Dass die Meldung bei einem Fehler wirkt, prüft kein Test; §7 nennt das als Grenze. Für den Verifier. | Plan §6 *Übrige Befunde der Testdateien* und Risiko *Befund verdeckt statt behoben*; `SPEC-049` Punkt 5 (`errcheck`); Rolle Verifier | `test/integration/extended_e2e_test.go` · „if err := pc.Close(context.Background()); err != nil {“ | ja (E0, 30 Läufe) | — (Hinweis) |

**Eigene Mutationen** (frische Kopie je Mutation; ohne Mutation grün):

| ID | Mutation | Weg im Test | Ergebnis |
|---|---|---|---|
| R1 | `envFailOnUnconsumed = "PGWIRE_RECORDER_FAILONUNCONSUMED"` | Konstante `EnvFailOnUnconsumed` | grün in allen Tests von `cli` und `bootstrap`, wie §7 angibt (F-456) |
| R1i | wie R1, Integrationstest | Binary, Literal `PGWIRE_RECORDER_FAIL_ON_UNCONSUMED` | rot, `TestE2EReplayNichtVerbraucht` („Exit-Code 0, erwartet 5“) |
| R2 | `envLogLevel = "PGWIRE_RECORDER_LOGLEVEL"` | Konstante `EnvLogLevel`, Hilfetext | rot, `TestParseLogLevelHilfe` (`cli`), dazu `TestRunLogLevel` und `TestRunStartfehlerJeStufe` (`bootstrap`) |
| R3 | Rückgabe nach `formVorpruefung` als PGR-E3001 statt PGR-E3003 | Paket `recording` | rot, `TestUnmarshalExtendedFehler`, `TestVorpruefungNenntOrt`, wie §7 angibt (F-457) |

## Antwort auf die Schwerpunkte

1. **Brücken nach Punkt 7.**
   - Je Paket ist `export_test.go` die einzige Testdatei im Paket des Codes (`cli`: `cli.go`, `doc.go`, `cli_test.go` in `cli_test`; `bootstrap`: `bootstrap.go`, `doc.go`, `bootstrap_test.go` in `bootstrap_test`).
   - `cli/export_test.go` enthält zwei Konstanten, `EnvFailOnUnconsumed = envFailOnUnconsumed` und `EnvLogLevel = envLogLevel`. `bootstrap/export_test.go` enthält zwei Funktionen, `Logger` und `Fail`, die genau einmal an `logger` bzw. `fail` weiterreichen und deren Ergebnis unverändert liefern.
   - Keine Variable, kein Typ-Alias, keine Funktion `Test…`/`Benchmark…`/`Example…`/`Fuzz…`, kein Zustand auf Paketebene. Signaturen nur aus exportierten Typen (`io.Writer`, `string`, `error`, `*slog.Logger`, `int`); einen Wert eines unexportierten Typs legt weder Brücke noch Test an.
   - `logger` und `fail` ruft auch der Produkt-Code (`bootstrap.go` in `Run`, `record`, `replay`); die Bedingung aus Punkt 7 für erzeugende Funktionen ist damit ohnehin erfüllt.
   - Die Doc-Kommentare beschreiben, was da ist (`AGENTS.md` §3.7).
2. **`errcheck` an neun `Close`-Stellen.** Korrekt und kein Verdecken; die Reihenfolge vor dem `cancel` ist dieselbe wie vorher und richtig begründet, siehe F-458. `(net.Listener).Close` an `belegt` bleibt ungeprüft, das lässt `SPEC-049` Punkt 5 zu (Entscheidung 2). Ein neues `_ =` gibt es im Diff nicht.
3. **`lebendTexte` und `vtTexte`.** Beide sind Funktionen, die je Aufruf ein neues Slice-Literal liefern; geteilten Zustand gibt es nicht. `append(lebendTexte(), vtTexte()...)` schreibt in ein frisches Slice (Länge gleich Kapazität, also neue Allokation), das frühere `append(append([]string{}, …), …)` war nur wegen der Globalen nötig. Inhalt und Reihenfolge der Texte sind unverändert.
4. **SA4000.** Dieselbe Prüfung, nicht mehr und nicht weniger. `senden() != nil || senden() != nil` rief `senden` zuerst einmal und nur ohne Fehler ein zweites Mal; die zwei Anweisungen tun genau das. SA4000 war ein Fehlalarm (die Ausdrücke sind gleich, haben aber Wirkung), kein Fehler im Test; vorher ging keine Prüfung verloren.
5. **Prüfung im Umbau verloren.** Keine gefunden. Stichproben gegen `1e4381b^`: alle Hunks von `test/integration` (sieben Dateien) einzeln, dazu `TestParseFailOnUnconsumedUmgebung`, `TestParseHilfeVorPruefung`, `TestParseLogLevelHilfe`, `TestLoggerSchwelle`, `TestLoggerOrtszeit`, `TestFailJeKlasse`, `TestRunStartfehlerJeStufe`. Geändert sind nur Paketname, Qualifizierung `cli.`/`bootstrap.`, Brücke und die Befundstellen; Erwartungen, Fristen und Abläufe sind gleich. In den geänderten Zeilen steht keine Deklaration `Abdeckung:`. Testliste gleich (70).
6. **Grüne Mutanten nach Schritt 19.** Beide richtig eingeordnet und nachgestellt: `envFailOnUnconsumed` (F-456), PGR-E3003 (F-457). Weitere grüne Mutanten habe ich nicht gefunden.
7. **Eigene Mutationen.** R1 bis R3 und R1i, Tabelle oben. Die Tabelle in §7 deckt jede Brücke und jede Befundstelle ab; R2 zeigt zusätzlich, dass `EnvLogLevel` auch in `bootstrap` gefangen ist. Ob die Abdeckung die DoD trägt, prüft der Verifier.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driving/cli/export_test.go`, `internal/bootstrap/export_test.go` | geprüft, ohne Befund nach `SPEC-049` Punkt 7. |
| `internal/adapters/driving/cli/cli_test.go` | geprüft, ohne Befund. Nur Paketname, Qualifizierung und Brücke; kein Befund von `make lint`. |
| `internal/bootstrap/bootstrap_test.go` | geprüft, ohne Befund. Nur Paketname, Qualifizierung und Brücke; kein Befund von `make lint`. |
| `test/integration/*_test.go` | geprüft, ohne Befund. Hinweis F-458. Paket `integration_test`, baut mit der unveränderten Stufe `integration`; `fuehreAus` ist eine reine Umbenennung von `exec_`. |
| Produkt-Code unter den drei Pfaden, `Dockerfile`, `.golangci.yml` | geprüft, ohne Befund. Unverändert in `1e4381b`. |
| `docs/plan/planning/in-progress/slice-harness-blackbox-einstieg.md` | geprüft. §1, §3 und §6 folgen dem Diff (`AGENTS.md` §3.9; Zahlen 9 und 13 stimmen mit §7); Hinweise F-455 (§6) und F-457 (§7). |
| Hard Rule 3.2, 3.6 (Suppression, Gate-Lockerung) | geprüft, ohne Befund. Kein `//nolint`, `.golangci.yml` unverändert. |
| Hard Rule 3.3 | geprüft, ohne Befund. `5c85717` und `ca46bf2` sind reine Renames (100 %); `1e4381b` enthält keinen Move. |
| Hard Rule 3.7, 3.11 (Kommentare) | geprüft, ohne Befund. Die geänderten Doc-Kommentare (`lebendTexte`, `vtTexte`, Paket-Kommentar) und die Kommentare der Brücken beschreiben, was da ist, und sagen nichts zu, was kein Test prüft. |
| Hard Rule 3.12 | geprüft, ohne Befund. Der Diff entscheidet keine Randform, die §6 nicht vor dem Code nannte; zur nachgezogenen Zugriffsliste F-455. |
| Commit-Message `1e4381b` | geprüft. Sie nennt `slice-harness-blackbox-einstieg`, [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) und `LH-QA-07`; Struktur-ID siehe F-454. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 4 |

**Summary:** 0 HIGH · 0 MEDIUM · 1 LOW (F-454: Commit-Message `1e4381b` nennt die Struktur-ID `SPEC-049`) · 4 INFO (F-455 White-Box-Liste erst im Code-Commit gemessen, alle Zugriffe in vorab entschiedenen Klassen; F-456 grüner Mutant `envFailOnUnconsumed` richtig eingeordnet, Namens-Asymmetrie aus dem Bestand; F-457 PGR-E3003-Einordnung stimmt, Formulierung ungenau; F-458 `Close`-Prüfungen korrekt, Reihenfolge vor `cancel` unverändert). Wiederkehrende Klasse: „Struktur-ID in der Commit-Message“ (zweites Auftreten nach F-336).

**Finding-Klassen dieses Laufs:** Struktur-ID in der Commit-Message · White-Box-Liste vor dem Code nur per Namenssuche

## Verdikt

**Merge-blockierend:** nein. Die Brücken halten Punkt 7 ein, Produkt-Code ist unverändert, und die Tests prüfen dasselbe wie vorher; die vier Befundklassen sind behoben, ohne etwas zu verdecken.

**Übergabe:** F-454 geht an den Implementer (keine Korrektur am Commit möglich, ohne Historie umzuschreiben; Annahme oder Begründung genügt). F-455 geht in die Closure (Steering Loop). F-456 bis F-458 gehen an den Verifier. Die Finding-Klassen gehen in §7 des Slice und von dort in den Zähler. Dieser Report ersetzt keine Verifikation.
