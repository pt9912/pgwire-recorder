# Review-Report: slice-v1-abschluss-verbindungsende-vergleich — 2026-10-16

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-vergleich` (Schritt 39a, Ende im Vergleich, Service), Diff `48325eae..18b828ab` (Basis `48325eae` laut Plan §4; Pfade `internal`, `test`, `docs/user` laut Plan §4, dazu Plan und Roadmap).

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec08 · <!-- d-check:ignore (Skill-Pfad des Repos, Zitat des Laufs) -->
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

- Slice-Plan `slice-v1-abschluss-verbindungsende-vergleich` (§1, §3, §4, §6, §7 Belege des Implementers, §8); Geber `slice-v1-abschluss-verbindungsende-einspielen`, `slice-v1-abschluss-verbindungsende-einspielen-extended`, `slice-v1-abschluss-verbindungsende-einspielen-grenzen`
- [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md), [ADR-0023](../plan/adr/0023-antwortvergleich-entscheidung.md), `SPEC-034`
- `LH-FA-24` (`LH-FA-24.a` *Ende der Aufzeichnung* mit dem vierten Punkt, *Meldung*, *Randformen des Vergleichs*; Spezifikationsstände `51145c0`, `cc563f8d`), `LH-FA-20` (`LH-FA-20.a` Tabelle der Fehlerregeln, *Gruppen*), `LH-FA-25` (`LH-FA-25.a`)
- `AGENTS.md` (Hard Rules, §3.7 bis §3.14), Hexagon-Richtung (`make a-check`, vom Implementer gefahren, nicht neu gefahren)

**Eigene Läufe des Reviewers** (Unit-Mutanten in Kopien von `go.mod`, `go.sum`, `internal`, `cmd` im Image der Stufe `deps`, `--network=none`, je Mutant genau eine Ersetzung mit erzwungenem Treffer, ein Lauf zur Zeit, Arbeitsbaum unberührt; unveränderte Kopie grün):

- Aus der Tabelle in §7, rot wie behauptet: M01 (72 Unterfälle `Erwartet`, nur diese), M02, M03, M06, M07, M08, M09a, M09b, M10, M12, M13 (zusätzlich `TestVergleichStruktur`), M14, M15, M18, M19 (nur Extended Flush, drei Fälle), M20, M21, M24, M32. M22 und M23 in eigener Form (Fehler beim Lesen beendet die Session ohne Fehler beziehungsweise zählt als verglichen): rot.
- Eigene Mutanten, rot: M01b (`connection_closed` doppelt angefügt: alle 72 Fälle `Erwartet`), `u.position <= ende` (nur `Server früher` und `ohne Fehlerantwort der Aufzeichnung`), `ende := len(l.v.ist) + 1`, `Abweichend` nur ohne `--continue-on-error`, erwartetes Ende nur mit `--continue-on-error`, `SchliesseOhneTerminate` auch bei nicht-fatalen Fehlerantworten mit Vergleich, `ende == 0` entfällt in `erste` (Panik beim ersten Vergleich ohne Ende, ein roter Test, der Rest des Pakets läuft nicht).
- **Grün:** (a) die beiden Mutanten E1 und E2 aus §7 wie dort; (b) `session` setzt nach `errSessionEnde` mit `continue` statt `return nil` fort (die Session endet nicht, die nächste Interaktion derselben Session wird gesendet): **kein Test rot**; in `48325eae` färbte ihn allein `TestVergleichFatalFortsetzung`, den dieser Diff löscht (F-1).
- E1 außerhalb der Schnittstelle der Tabelle: ein Wegwerf-Test mit einer Aufzeichnung, deren `command_complete` ein Feld `C` trägt, zeigt mit entferntem Typ-Check „SQLSTATE aufgezeichnet „99999““ statt „“ (F-3).
- `TestE2EPlayVergleichFatal` am Binary gegen PostgreSQL 17 (`tools/test/run-integration-tests.sh` mit `PGR_TESTS_1`, `PGR_ROLLEN=scram`): `PASS` in Phase 1 (die rote Phase 2 ist die Folge der Auswahl eines Tests, nicht des Diffs); Datei unberührt (`git diff --stat` über `test/` leer).
- Messung `497` nachgerechnet: `git diff --numstat 48325eae..85434400 -- internal test docs/user` gibt Einfügungen 355 und Löschungen 142, bis `18b828ab` dieselbe Summe (die Commits danach ändern nur Plan); Produktcode `play.go` 39 und `play_vergleich.go` 50 gleich 89, davon 48 Kommentarzeilen und 40 Zeilen Code; Tests 396 (`play_vergleich_ende_test.go` 288, `play_vergleich_test.go` 108); Faktor 396 / 89 = 4,4 und 396 / 40 = 9,9; 62 % der Marke 800, Meldeschwelle 640 nicht erreicht. Stimmt.
- 208 Unterfälle nachgezählt: 72 (`Erwartet`, 6 Fälle mal 3 Formen mal 4) plus 80 (`Abweichung`, 20 Fall-Formen mal 4) plus 44 (`Gewoehnlich`, 11 mal 4) plus 12 (`Verbindungsverlust`).
- `make abdeckung-check` und `make kopf-check`: Exit 0.

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Zusage „die Session endet, die nächste läuft“ (DoD Punkt 1, §6, Kommentare von `fatal`, `interaktion` und der Tests) hat keine Mutation und keinen roten Test: Setzt `session` nach `errSessionEnde` mit `continue` fort, bleibt jeder Test grün, weil B in jedem Fall der Tabelle die letzte Interaktion der Session ist. Der einzige Test, der das fing (`TestVergleichFatalFortsetzung`, FATAL in A, B wird nicht gesendet), ist gelöscht, obwohl §3 sagt, er gehe in die Tabellen ein; Failure-Szenario: nach dem Ende sendet `play` die nächste Anfrage derselben Session in die geschlossene Verbindung (`PGR-E4003`, Exit-Code 4 statt erwartetem Ende mit Exit-Code 0), und kein Gate wird rot. | `AGENTS.md` §3.10, §3.11 | `internal/hexagon/services/play_vergleich_ende_test.go` · `{ID: 1, Interactions: []model.Interaction{aufgezeichnet(1, "A", ok...), f.interaktion(c.soll)}},`; Mutationsstelle `internal/hexagon/services/play.go` · `if errors.Is(err, errSessionEnde) {` | ja — der Mutant `continue` statt `return nil` an dieser Stelle lässt `go test ./internal/hexagon/services/` grün; auf `48325eae` rot durch `TestVergleichFatalFortsetzung` | Zusage ohne Mutation (Session endet nach Ende) |
| F-2 | LOW | §3 des Plans sagt „je Fall in den Formen einfach, Extended Sync und Extended Flush“, der Diff fährt vier der neun Fälle von `TestVergleichEndeAbweichung` nur in zwei Formen (`ohneFlush()`); der Grund (Flush gibt es nur bei aufgezeichnetem Ende, und das Lesen endet dort vor der Fehlerantwort) steht weder im Plan noch im Kommentar von `ohneFlush`. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-vergleich.md` · `je Fall in den Formen einfach, Extended Sync und Extended Flush`; `internal/hexagon/services/play_vergleich_ende_test.go` · `func ohneFlush() []vform` | nein — Plan-Text gegen Diff, kein Gate | Plan nennt mehr Formen als der Diff trägt |
| F-3 | LOW | Mutant E1 (Typ-Check in `erste` entfällt) ist nicht äquivalent: `Validate` und der YAML-Leser lassen ein Feld `C` an einer `command_complete` zu, und die Meldung nennt dann diesen Wert statt „leer“; §7 begründet die Äquivalenz mit „Nur eine `error_response` trägt das Feld `C`“, was das Modell dokumentiert, aber nicht erzwingt. Der Kommentar von `erste` („der aufgezeichnete leer, wenn dort keine Fehlerantwort steht“) sagt mehr zu, als ein Test prüft. | `AGENTS.md` §3.10, §3.11 | `internal/hexagon/services/play_vergleich.go` · `ende <= len(v.soll) && v.soll[ende-1].Type == model.ResponseErrorResponse` | ja — Test mit einer Aufzeichnung, deren `command_complete` `notice: {C: …}` trägt, färbt E1 | Grüner Mutant als äquivalent eingeordnet, Voraussetzung nicht erzwungen |
| F-4 | INFO | Befund M01 stimmt: Das angefügte `connection_closed` ist über die Schnittstelle nur im erwarteten Ende beobachtbar, weil die Umwandlung zur Art Fehler jedes andere Fehlen deckt; die zwei Zeugen des Architect (Ende an früherer oder späterer Position, Server-FATAL ohne Ende) bleiben auch ohne Anfügen rot und können es nicht ändern. Die Zusage trägt: M01 färbt alle 72 Fälle `Erwartet`, das doppelte Anfügen (M01b) ebenso; ein Fall fehlt nicht. | Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-vergleich.md` · `M01 färbt nur`; `internal/hexagon/services/play.go` · `l.v.nimm(model.Response{Type: model.ResponseConnectionClosed})` | ja — M01 und M01b | — |
| F-5 | INFO | Mutant E2 (Schutz `ende <= len(v.soll)`) ist mit `Validate` äquivalent: Eine gültige Aufzeichnung endet auf `ready_for_query` oder `connection_closed`, die Folge des Servers trägt beides nicht vor der Fehlerantwort (das Lesen endet am `ready_for_query`, das Ende hängt `play` selbst an), also reicht `soll` immer bis zur Position der Fehlerantwort; der Schutz steht ohne Kommentar, der das sagt. | Maintainability | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-vergleich.md` · `E2: Schutz` | nein — kein Test kann den Zweig ohne Umgehung von `Validate` erreichen | — |
| F-6 | INFO | Nicht Teil dieses Diffs, aber Grundlage der neuen Entscheidung „Ende“: `fehlerantwort` fällt bei leerem Feld `V` auf `S` zurück, `beendetVerbindung` im Modell lässt ein vorhandenes, auch leeres `V` allein entscheiden (SPEC-002, „genau“); die Spezifikation sagt „Feld `V`, ohne es `S`“ und nennt das vorhandene leere `V` nicht. Der Server sendet `V` in der Praxis nie leer; für den Architect. | Maintainability | `internal/hexagon/services/play.go` · `schwere := r.Fields["V"]`; `internal/hexagon/model/response.go` · `schweregrad, ok := r.Fields["V"]` | nein — ein Test mit vorhandenem leerem `V` würde die Lesarten trennen, keine Zusage verlangt ihn | Schweregrad-Lesart in Service und Modell verschieden |

**Zuständigkeit.** F-1: Implementer (Fall mit einer Interaktion hinter B in Session 1 im Läufer oder ein eigener Test, dazu die Zeile der Zusage in §7). F-2: Implementer (Plan §3 oder Test nachziehen, `AGENTS.md` §3.9). F-3: Implementer (Test oder Einordnung in §7 berichtigen). F-4 und F-5: keine Aktion. F-6: Architect.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/hexagon/services` Produktcode (`play.go`, `play_vergleich.go`): `fatal`, `erste`, `gewoehnlich`, `abschluss`, Zähler | geprüft, ohne Befund außer F-1 und F-3; 32 Mutanten gefahren, alle bis auf E1, E2 und `continue` statt `return nil` rot; Verhalten entspricht `LH-FA-24.a` *Ende der Aufzeichnung* (vier Punkte), *Meldung*, *Zähler* und der Tabelle der Fehlerregeln |
| `internal/hexagon/services` Tests (`play_vergleich_ende_test.go`, `play_vergleich_test.go`) | geprüft, ohne Befund außer F-1, F-2; 208 Unterfälle sind trennscharf: Die drei Formen fahren verschiedene Pfade (`bisBereit`, `extended` mit Sync, Lesen bis zur Antwortzahl), die Schweregrade ergeben je Element einen eigenen roten Mutanten (M14 bis M16), die Optionen `continue` und `allow` je einen (M09a, M09b, M12) |
| Gelöschte Tests `TestVergleichFatal`, `TestVergleichFatalRand` | geprüft, ohne Befund; ihre Aussagen (Ablauf, Zähler, Meldung, V ERROR bei S FATAL, FATAL an der Stelle von `command_complete`, frühere Abweichung) liegen in `Erwartet`, `Abweichung` und `Gewoehnlich` und färben die Mutanten; `TestVergleichFatalFortsetzung` siehe F-1 (Reihenfolge der Meldungen fängt `TestVergleichFortsetzung`, Mutant mit vorangestellter Meldung rot) |
| `TestVergleichFatalExtendedGruppen` (nicht in die Tabelle überführt) | geprüft, ohne Befund; die Begründung trägt: Der Fall braucht eine frühere Gruppe mit `Flush` und eine weitere Gruppe dahinter, die der Läufer nicht baut; M03 färbt ihn; der Kommentar nennt den Fall eng („Gruppe mit Flush, die in der Aufzeichnung eine Fehlerantwort ERROR mit demselben SQLSTATE ohne Ende zeigt“) |
| Kommentare in `play.go`, `play_vergleich.go` (`AGENTS.md` §3.7, §3.11) | geprüft, ohne Befund außer dem Teil „Session endet“ (F-1) und dem Satz zu `erste` (F-3); die drei Kommentare der Sendung V-1 (`interaktion`, `abschluss`, `Zaehler`) tragen M19 bis M21; keine verworfene Alternative, keine Chronik |
| Altes Wort „bis zum ReadyForQuery“ über `internal/` einschließlich `ports/` | geprüft, ohne Befund; verbliebene Stellen (`ports/driving/record.go`, `ports/driven/upstream.go`, `adapters/driven/postgres/upstream.go`, `play.go` bei `anfrage`, `extended`, `bisBereit`) betreffen Aufnahme oder die Pfade, die tatsächlich bis zum `ReadyForQuery` lesen; der Kommentar von `Play` nennt die letzte Gruppe mit `Flush` |
| Hexagon-Richtung, Port `Einspielziel`, Einspiel-Adapter, CLI, Hilfetext, Bootstrap | geprüft, ohne Befund; `git diff --stat 48325eae..18b828ab` nennt nur `internal/hexagon/services`, `docs/user/abdeckung-unit.md`, Plan und Roadmap; Imports unverändert |
| `test/integration`, `TestE2EPlayVergleichFatal` | geprüft, ohne Befund; unberührt und grün (siehe Läufe) |
| Plan §1, §3, §4, §6, §8 gegen den Diff (`AGENTS.md` §3.9), Randformen (`AGENTS.md` §3.12) | geprüft, ohne Befund außer F-2; der Zähler-Satz in §6 ist die Angleichung an die Randform, die der Extended-Punkt von §6 und `LH-FA-24.a` *Zähler* schon nannten, keine neue Randform im Code-Commit; der Code entscheidet keine Randform still (Stelle von `connection_closed`, Reihenfolge der Gründe, Schweregrad und leerer SQLSTATE stehen in §6 und der Spezifikation) |
| Sendungen der Geber (drei Kommentare V-1, zwei Fälle der letzten `Flush`-Gruppe, Einleitungssatz V-3) und Nachzählung in §8 | geprüft, ohne Befund; Kommentare geändert und belegt, Fälle vorhanden (M19 bis M21), Einleitungssatz steht in `spec/spezifikation.md`; Liefer-Punkte 1, Schichten eine (`docs/user` konservativ als zweite), Marke 800 nicht erreicht |
| `docs/user/abdeckung-unit.md` | geprüft, ohne Befund; `make abdeckung-check` Exit 0 |
| Roadmap (Zustandsfeld „In Arbeit“) | geprüft, ohne Befund; nennt Zustand und Schritt, keine Chronik |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zusage ohne Mutation (Session endet nach Ende) · Plan nennt mehr Formen als der Diff trägt · Grüner Mutant als äquivalent eingeordnet, Voraussetzung nicht erzwungen · Schweregrad-Lesart in Service und Modell verschieden

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM): Der Diff löscht den einzigen Test, der „die Session endet“ nach einem Ende mit Vergleich prüfte, und die Tabelle trägt die Zusage nicht (Mutant grün, weil B stets die letzte Interaktion ist). F-2 und F-3 sind LOW, F-4 bis F-6 INFO. Produktcode, die Mutationstabelle (bis auf die genannten Stellen), die Messung 497 und die Unberührtheit von Port, Adapter, CLI und `TestE2EPlayVergleichFatal` tragen.

**Übergabe:** Findings gehen an den Implementer (F-6 an den Architect); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
