# Review-Report: slice-v1-abschluss-verbindungsende-sqlite-schreiber — 2026-10-15

**Review-Art:** Code — gegen Plan (§1, §3, §4, §6 mit *Zählung der Mutationstabelle*, §7, §8), die Entscheidungen der Spezifikation und `AGENTS.md` §3.7 bis §3.14.

**Gegenstand:** `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Schritt 36b), Diff `b5933e8..e63b184` (Architect `1ce3d3a`, Umsetzung `e63b184`).

**Skill:** `.harness/skills/reviewer.md` @ `e63b184` · <!-- d-check:ignore (Skill-Pfad als Inline-Code) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

> **Zitier-Form.** Kennung, nicht Adresse: `slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill). Das `pfad`-Feld zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- Slice-Plan `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (in `in-progress/`), §7 *Belege des Implementers*; Nehmer-Plan `slice-v1-abschluss-verbindungsende-record` (Bindung)
- [ADR-0039](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md) (Punkt 4), [ADR-0013](../plan/adr/0013-sqlite-recording-backend.md)
- `LH-FA-25` (`.a`), `LH-FA-22.a` (*Ende und Version beim Schreiben von sqlite*), `SPEC-001` (*Schreiber in sqlite*), `SPEC-043`
- `AGENTS.md` (Hard Rules), insbesondere §3.7 bis §3.14

**Eigene Läufe (Belege des Reviews).** Kopie von `go.mod`, `go.sum`, `internal/` im Scratchpad, Testbinary einmal gebaut (Image `pgr-36b:deps`, `--network=none`, Bind-Mount), je Mutant ein Lauf gegen `TestSQLite(Version|Ende|SchreibenNimmt|Roundtrip)`, Kopie danach zurückgesetzt. Unverändert: PASS. Kein `make gates`, kein Eingriff im Repo außer diesem Report.

| Mutant (eigener Nachbau) | Ergebnis |
|---|---|
| `ausserhalb` (Version in eigener Transaktion nach dem Commit, Zeilenzahl geprüft) | rot: `FehlerInDerTransaktion`/Auslöser beim Setzen, `ZeileUnpassend` (drei); wie §7 |
| `davor` (Version in eigener Transaktion vor dem Einfügen, Zeilenzahl geprüft) | rot: nur `FehlerInDerTransaktion`/doppelte sequence; `ZeileUnpassend` bleibt grün (§7 nennt sie, siehe F-3) |
| `nachSchliessen` (Flag nur bei erfolgreichem Schließen) | rot: `SchliessenScheitert`, `WerSetzt`/zweite Session mit Ende im selben Schreibvorgang; wie §7 |
| `flagVorCommit` | rot: beide Fälle von `FehlerInDerTransaktion`; wie §7 |
| `commitTrotz` | rot: Auslöser beim Setzen, `ZeileUnpassend` (drei); wie §7 |
| `feld1` | rot: nur `NachInhalt`/Feld 1 mit Ende; wie §7 |
| `nurErste` | rot: u. a. `NachInhalt`/spätere Session, `WerSetzt`/spätere Session, `AbbruchUnterbrichtNicht`; wie §7 |
| `ohneWert` / `ohneKey` | rot: nur `ZeileUnpassend`/Wert 2 bzw. nur /andere Zeile mit Wert 1; wie §7 |
| `zeilenzahl` (`n < 0`) | rot: `ZeileUnpassend` alle drei; wie §7 |
| `userVersion` (`PRAGMA user_version = 2`) | rot: `NachInhalt` alle Fälle mit Ende |
| `feld` K=2 (Anlegen schreibt `rec.Version`, wenn 2) | rot: `NachInhalt`/Feld 2 mit und ohne Ende, alle sechs `Roundtrip`-Fälle; wie §7 |
| Stellen, `ohneVersion` und `ablehnen`, je 14 Läufe (Schlüssel: Startup-Parameter `user`) | je Lauf genau der Fall der Stelle rot, kein anderer; wie §7 |
| eigen: Anlegen schreibt `rec.Version`, wenn die Aufzeichnung kein Ende trägt | rot: `Feld -1/0/2/3 ohne Ende`, `ohne Sessions`, `nur empty_sessions` u. a.; `Feld K mit Ende` bleibt grün (F-2) |

Die Prüfung `feldNeutral` (Version 5) wirkt als Isolierung: Mutanten, die an einem festen Wert K des Modellfelds hängen, färben nur die Fälle „Feld K“ (und, bei K = 2, den Rückweg des `Roundtrip`).

---

## Findings

Spalten gespiegelt aus dem §Output-Schema des Reviewer-Skills; bei Abweichung gilt der Skill (`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill).

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Das Doc des Tests und seine Zeilen in `abdeckung-unit.md` (unter `LH-FA-22`, `LH-FA-25`, `LH-QA-06`) sagen zu, „ein späterer Schreibvorgang an eine Datei der Version 2 ergänzt Sessions ohne Änderung der Version“. `TestSQLiteRoundtripMitEnde` schreibt nicht ein zweites Mal; diesen Fall prüft `TestSQLiteVersionWerSetzt`. | `AGENTS.md` §3.11 | `internal/adapters/driven/recording/sqlite_schreiben_ende_test.go` · „ein späterer Schreibvorgang an eine Datei der Version 2 ergänzt Sessions ohne Änderung der Version.“ (Doc von `TestSQLiteRoundtripMitEnde`) | ja — Mutant `jede` (jede Session mit Ende setzt) lässt `Roundtrip` grün, `WerSetzt` wird rot | Zusage im Doc weiter als die Prüfung des Tests |
| F-2 | LOW | Die Tabelle in §7 führt für „Feld K bestimmt die Version nicht“ einen Mutanten (`feld`), der „Feld K mit Ende“ und „Feld K ohne Ende“ zusammen färbt (§7 nennt es „beide“); §6 *Zählung der Mutationstabelle* verlangt je Element eine Mutation, die nur es bricht. Die Tests trennen die Fälle (eigener Mutant: Anlegen schreibt `rec.Version` nur ohne Ende, färbt die Fälle „ohne Ende“, „mit Ende“ bleibt grün); die Tabelle weist es nicht aus. Ebenso trennt `zeilenzahl` die Fälle „Zeile fehlt“ und „Wert 2“ nicht. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-sqlite-schreiber.md` · „je K `VersionNachInhalt`/„Feld K mit Ende“ und „Feld K ohne Ende“ (beide, weil beide dasselbe Anlegen treffen)“ | ja — eigener Mutant (siehe Tabelle oben) | Mutationstabelle zählt Paare statt Elemente |
| F-3 | INFO | Die Zeile „Version in derselben Transaktion (doppelte `sequence`)“ nennt als roten Test auch `VersionZeileUnpassend` (drei). Der Nachbau von `davor` mit geprüfter Zeilenzahl färbt nur den Fall der doppelten `sequence`; die drei Fälle von `ZeileUnpassend` bleiben grün (es wird nichts eingefügt, `meta` bleibt unverändert). Die Aussage trägt nur für einen Mutanten ohne Zeilenprüfung; die Zusage des Risikos in §6 (beide Wege rot) bleibt belegt. | Maintainability (Beleg in §7) | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-sqlite-schreiber.md` · „`UPDATE` in einer eigenen Transaktion vor dem Einfügen (`davor`)“ | ja — Mutant `davor` mit Zeilenprüfung | Mutantenzeile der Tabelle nicht reproduzierbar |
| F-4 | INFO | Die 14 Stellen-Fälle stehen als zweite Tabelle neben `TestSchreibenNimmtFormverletzungAn` (YAML); §6 nennt eine geteilte Tabelle als Möglichkeit, die Zeilen spart. Kein Verstoß, Rückbezug für den Format-Folge-Slice (`slice-v1-abschluss-verbindungsende-replay`, `-einspielen`). Zuständig: Implementer bei Bedarf. | Maintainability | `internal/adapters/driven/recording/sqlite_schreiben_ende_test.go` · `TestSQLiteSchreibenNimmtFormverletzungAn` | nein | Doppelte Fall-Tabelle je Format |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Anlegen schreibt `meta.version` immer mit `model.Version` (nicht `rec.Version`); Kommentar an `legeSchemaAn` benennt es | geprüft, ohne Befund |
| „Session mit Ende“ über `Recording{Sessions: {s}}.HasEnd()` im Adapter; `internal/hexagon/model` im Diff unberührt, keine neue Methode, Port unberührt | geprüft, ohne Befund |
| `UPDATE … WHERE key='version' AND value=?` genau eine Zeile, sonst Fehler (`PGR-E3001` über `schreibe`) mit Rollback; Mutanten `ohneWert`, `ohneKey`, `zeilenzahl` je rot | geprüft, ohne Befund |
| Flag `endeGesetzt` nach dem Commit der Session, vor dem Schließen der Datenbank; `flagVorCommit` und `nachSchliessen` rot | geprüft, ohne Befund |
| Kein `PRAGMA user_version` im Schreiber; `userVersion`-Mutant rot | geprüft, ohne Befund |
| Ende-Zeile (kind, direction, sequence, position, group_no `NULL` bei einfacher Interaktion, `fields` `{}`): Darstellung unverändert aus der bestehenden Abbildung, `EndeDarstellung` prüft Stelle und Form (`typeof`, `length`); Schreiber prüft die Form des Endes nicht (14 Fälle, je Fall isoliert) | geprüft, ohne Befund |
| Reihenfolge E1000 vor E3001 (`EndeMitFalscherNummer`: keine Datei); Aufrufarten und Abbruch (`WithoutCancel`, `AbbruchUnterbrichtNicht`) | geprüft, ohne Befund |
| `contextcheck`-Behebung: `context.Background()` am Aufruf in Tests, wie in den übrigen Test-Dateien des Pakets; kein `//nolint`, `.golangci.yml` im Diff unberührt, keine Suppression | geprüft, ohne Befund |
| Größe (§3.14): Summe der Metrik von §4 nachgemessen mit `git diff --numstat b5933e8..HEAD -- internal docs/user` = 21+14+6+520+32+3 = 596 (Marke 1000, Meldeschwelle 800); Schichten zwei (Recording-Adapter, `docs/user`), Liefer-Punkte zwei; Plan §8 nennt die Zählung | geprüft, ohne Befund |
| Bindung in `slice-v1-abschluss-verbindungsende-record` (§1, §6, Zählung in §8, Liefer-Punkte 3, Schichten 2) im selben Commit `1ce3d3a`; Adresse nimmt an (liegt in `next/`, schließt die Sendung nicht aus) | geprüft, ohne Befund |
| `spec/spezifikation.md` (`LH-FA-22.a` (*Ende und Version beim Schreiben von sqlite*), `SPEC-001`, `SPEC-043`, Änderungstabelle): keine ADR-, Slice- oder Wellen-Kennung im Stratum (§3.4); Aussagen decken sich mit Code und Tests | geprüft, ohne Befund |
| Kommentare im Produktcode (§3.7): Zusage/Kopplung, kein Konjunktiv, keine Chronik | geprüft, ohne Befund |
| `internal/adapters/driven/recording/` Core-Reinheit: Imports nur `database/sql`, `fmt`, `strconv`, `model` | geprüft, ohne Befund |
| `docs/reviews/`, `.harness/`, `harness/`, `Makefile`: im Diff unberührt | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusage im Doc weiter als die Prüfung des Tests · Mutationstabelle zählt Paare statt Elemente · Mutantenzeile der Tabelle nicht reproduzierbar · Doppelte Fall-Tabelle je Format

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM; F-1 und F-2 sind LOW mit Anker, F-3 und F-4 INFO.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
