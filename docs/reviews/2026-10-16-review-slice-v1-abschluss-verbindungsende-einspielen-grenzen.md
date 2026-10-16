# Review-Report: slice-v1-abschluss-verbindungsende-einspielen-grenzen — 2026-10-16

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-einspielen-grenzen` (Schritt 38c), Diff `fa149c68..63dbbdb3`; Lieferung `670ec850`, danach Belege und Roadmap.

**Skill:** `.harness/skills/reviewer.md` @ `c5c6ec08`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-16

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-v1-abschluss-verbindungsende-einspielen-grenzen` (§1, §3, §4, §6, §7 *Belege des Implementers*, §8)
- Geber `slice-v1-abschluss-verbindungsende-einspielen` (38a) und `slice-v1-abschluss-verbindungsende-einspielen-extended` (38b)
- Spezifikation `LH-FA-20.a` (*Interaktion*, Einstufung; Architect-Commits `526a50f`, `5c8c592`, `e92825a`, `51145c0`), `SPEC-034` (`PGR-E4003`)
- `ADR-0039` und `ADR-0017` als Inline-Code gelesen: [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md)
- berührte `LH-*`-IDs: `LH-FA-20`, `LH-FA-25`
- `AGENTS.md` (Hard Rules, §3.7 bis §3.14)

**Eigene Läufe.** Basis bestätigt: `git diff --stat fa149c68..63dbbdb3 -- internal cmd` ist leer (Produktcode 0 Zeilen). Fünf Mutanten der Tabelle in §7 selbst nachgefahren, je ein voller `make test-integration`, ein Lauf zur Zeit, Mutant im Arbeitsbaum **nicht** committed (das Image baut aus dem Arbeitsbaum; der Nachweis der Gates ist hier nicht Gegenstand), danach `git checkout internal`:

| Mutant | Stelle | Ergebnis (selbst gesehen) |
|---|---|---|
| Ende immer erwartet | `services/play.go` Z. 261 `l.ende = traegtEnde(in)` | nur `TestE2EPlayEndeGrenzeOhneOption` rot (Exit-Code 0, erwartet 4); `…OhneEnde` grün |
| Option allein genügt | Z. 261 `l.ende = s.optionen.AllowRecordedErrors` | nur `TestE2EPlayEndeGrenzeOhneEnde` rot (Exit-Code 0, `eingespielt=2`, also am Grund, nicht nur am Zähler); `…OhneOption` grün |
| Stelle der Interaktion fehlt | Z. 203 ohne `Interaktion %d` | beide Grenzfälle rot (Exit 4, aber ohne die Stelle), dazu sechs Bestandstests |
| Code `PGR-E4003` | Z. 417 `CodeServerError` | beide Grenzfälle rot, dazu `TestE2EPlayFortsetzungAbbruch` |
| Exit-Code 4 | `model/fehler.go` `exitCode` liefert für `PGR-E4003` 5 | beide Grenzfälle rot (Exit-Code 5, erwartet 4), dazu drei Bestandstests |

Die Tabelle in §7 stimmt mit den Läufen überein, einschließlich der Isolation der beiden ersten Mutanten.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | INFO | Der Helfer prüft `code=PGR-E4003` und `Session 1, Interaktion 2` als zwei Teilstrings des ganzen stderr, nicht in derselben Zeile, und weder SQLSTATE `57P01` noch Schweregrad `FATAL`. Ein Mutant, der SQLSTATE und Meldung aus der Fehlerzeile streicht, bliebe grün; die Zusage dieses Slice ist `PGR-E4003`, Exit 4 und Stelle (§6 *Zähler und Nachricht*), SQLSTATE und Meldung nennt `LH-FA-20.a` *Meldungen* für Fehlerantworten allgemein. Zufallspfad-Risiko besteht nicht: der Pfad der Aufzeichnung steht nur in der Zeile „play gestartet“ und enthält keinen der beiden Teilstrings. | `AGENTS.md` §3.10 (Maintainability) | `test/integration/play_ende_grenzen_e2e_test.go` · `!strings.Contains(stderr, "Session 1, Interaktion 2")` | ja — Mutant ohne SQLSTATE in `play.go` „Fehlerantwort des Servers %s „%s“, die Verbindung endet“ | Teilstring ohne Bindung an dieselbe Zeile |
| F-2 | INFO | Der Wortlaut „Session 1, Interaktion 2“ stammt aus dem Code (`eingeordnet(err, "Session %d, Interaktion %d", …)`); die Spezifikation (*Meldungen*) sagt nur, dass die Zeile die `id` der Session und die `sequence` nennt. Der Test bindet den Wortlaut und färbt bei einer Umformulierung rot, ohne dass eine Zusage bricht. | `AGENTS.md` §3.11 (Maintainability) | `test/integration/play_ende_grenzen_e2e_test.go` · `"Session 1, Interaktion 2"` | nein | Test bindet Wortlaut, den die Spezifikation nicht fixiert |
| F-3 | INFO | Die Abdeckungs-Deklaration von `TestE2EPlayEndeGrenzeOhneEnde` beginnt mit „wie TestE2EPlayEndeGrenzeOhneOption, mit …“ und steht damit in `docs/user/abdeckung-e2e.md` als Zeile, die ohne die andere nicht lesbar ist; alphabetisch steht sie vor der Zeile, auf die sie verweist. Inhaltlich trägt sie nur, was der Test prüft. | Maintainability | `docs/user/abdeckung-e2e.md` · „wie TestE2EPlayEndeGrenzeOhneOption, mit --allow-recorded-errors“ | nein | Abdeckungs-Zeile verweist auf Nachbarzeile |
| F-4 | INFO | `LH-FA-25.a` *Negative* nennt für ein Ende, das die Aufzeichnung „mit einem anderen Fehlercode“ zeigt, einen Verbindungsfehler; `spec/spezifikation.md` (Einspielen, Schritt 6) lässt das Ende mit `--allow-recorded-errors` „gleich welcher SQLSTATE“ erwartet sein. Beide Fixtures tragen `57P01`; der Slice entscheidet die Frage nicht und beansprucht sie nicht. Zuständig: Architect (Spezifikation) bei Gelegenheit, nicht dieser Diff. | `LH-FA-25` / `LH-FA-20.a` | `spec/lastenheft.md` · „oder mit einem anderen Fehlercode zeigt“ | nein | Lastenheft und Spezifikation lesbar verschieden bei SQLSTATE |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/`, `cmd/` (Produktcode, `fa149c68..63dbbdb3`) | geprüft, ohne Befund (Diff leer, wie in §1 zugesagt; keine Rückführung nötig) |
| `test/integration/play_ende_grenzen_e2e_test.go` | geprüft, ohne Befund außer F-1, F-2: jeder Fall prüft Exit 4, Code und Stelle; beide Mutanten aus §6 färben nur den Fall ihrer Zusage; `-OhneEnde` rot am richtigen Grund (Exit 0), nicht nur am Zähler |
| `test/integration/testdata/play-ende-grenze-ohne-option.yaml` | geprüft, ohne Befund: gültig (Format 2 lädt; beide Bestandsläufe mit Mutanten laden sie), zweite Interaktion trägt `error_response` `FATAL` `57P01` und als letzte Antwort `connection_closed`: eine Aufzeichnung mit Ende an der Stelle |
| `test/integration/testdata/play-ende-grenze-ohne-ende.yaml` | geprüft, ohne Befund: gültig, zweite Interaktion trägt `ERROR` `57P01` und `ready_for_query`, kein Ende; der Server sendet `FATAL`: Aufzeichnung ohne Ende an der Stelle, mit Option erwartet nur der Mutant „Option allein genügt“ das Ende |
| Isolation der Fälle | geprüft, ohne Befund: eine Datei je Fall, keine geteilte Fixture mit `play-ende-erwartet.yaml`, getrennte Datenbanken; kein grüner Mutant der fünf Zusagen gefunden |
| `docs/user/abdeckung-e2e.md` | geprüft, ohne Befund außer F-3: vier Zeilen (je Test unter `LH-FA-20` und `LH-FA-25`), jede sagt zu, was der Test prüft; `Server 17` deckt `make test-integration` (Referenz) |
| Slice-Plan (§3.9) | geprüft, ohne Befund: §1, §3 (Dateinamen), §4 (Basis `fa149c68`), §6 und §7 folgen dem Diff; Zeilenangaben der Mutanten (203, 261, 417) stimmen mit `play.go`; Zählung der Schichten (eine, konservativ zwei) steht in §8 |
| Messung (§4, §7) | geprüft, ohne Befund: `git diff --numstat fa149c68..63dbbdb3 -- internal test docs/user` ergibt 52 + 18 + 18 + 4 = 92 (Marke 400, Meldeschwelle 300 nicht erreicht); je Fall 46 ist 92 durch zwei, die Hilfe `spieleEndeNichtErwartet` und der Kopf (rund 30 Zeilen) sind geteilt |
| Randformen (§3.12) | geprüft, ohne Befund: Der Code und die Tests entscheiden keine Randform still; Stelle des `connection_closed` (a), ein Grund (d) und die Fälle (b), (c) stehen in §6 und in der Spezifikation |
| Adressen (§3.13) | geprüft, ohne Befund: Der Diff weist keinem Nehmer etwas Neues zu (`-vergleich`, `-doku`, `-einspielen-extended` unverändert in diesem Diff); `Übernommen aus …-extended` ist die Annahme dieses Slice |
| Suppression (§3.2), Core-Reinheit | geprüft, ohne Befund (kein `nolint`, kein Import im Hexagon geändert) |
| Roadmap | geprüft, ohne Befund: Zeile „In Arbeit“ nennt Zustand und Verantwortlichen, keine Chronik |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Teilstring ohne Bindung an dieselbe Zeile · Test bindet Wortlaut, den die Spezifikation nicht fixiert · Abdeckungs-Zeile verweist auf Nachbarzeile · Lastenheft und Spezifikation lesbar verschieden bei SQLSTATE

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM; die vier INFO sind ohne erwartete Aktion des Implementers (F-4 an den Architect).

**Übergabe:** Findings gehen an den Implementer (F-1 bis F-3 zur Kenntnis) und F-4 an den Architect; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** — er wird über Läufe hinweg nicht wieder gelesen. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
