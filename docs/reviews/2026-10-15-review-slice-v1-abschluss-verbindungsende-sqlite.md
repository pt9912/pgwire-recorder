# Review-Report: slice-v1-abschluss-verbindungsende-sqlite — 2026-10-15

**Review-Art:** Code — Code-Review gegen Plan (§1, §3, §4, §6, §7, §8) und Konventionen
(`AGENTS.md` §3.7 bis §3.14).

**Gegenstand:** `slice-v1-abschluss-verbindungsende-sqlite` (Schritt 36a, SQLite-Leser),
Diff `fa8486b..dbd2c3a` (Architect `f449ea1`, Umsetzung `dbd2c3a`).

**Skill:** `.harness/skills/reviewer.md` @ `dbd2c3a` · <!-- d-check:ignore (Skill-Pfad des Laufs) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-15

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

- `slice-v1-abschluss-verbindungsende-sqlite` (§1, §3, §4, §6 mit Fällen und Mutationen, §7 Belege des Implementers, §8)
- [`ADR-0039`](../plan/adr/0039-aufzeichnung-kennt-das-verbindungsende.md) (Aufzeichnung kennt das Verbindungsende), [`ADR-0013`](../plan/adr/0013-sqlite-recording-backend.md)
- `LH-FA-25` (.a), `LH-FA-22`, `LH-QA-06`; `SPEC-001`, `SPEC-002`, `SPEC-041`, `SPEC-043` (neuer Punkt *Ende des Servers*)
- `AGENTS.md` (Hard Rules, insbesondere §3.7 bis §3.14)

**Nachgefahren (Kopie des Baums im Scratchpad, Bind-Mount, `--network=none`, Image `pgwire-recorder:test`, kein `make gates`):**
Tests `TestLesenEnde*`, `TestLesenMeta`, `TestLesenMeldung*`, `TestLesenReihenfolgeDerZeilen`. Jede Mutation einzeln, danach zurückgesetzt:

| Mutation | Ergebnis |
|---|---|
| Zahlenmenge in `pruefeMeta` entfällt | rot: `TestLesenEndeVersionen` (0, 3, -1, 10, 20 nur über die Meldung `Version <Zahl>`; die Zahl über 64 Bit mit `PGR-E3003` statt `PGR-E3002`), `TestLesenMeta` (Version vor Inhalt, Überlauf), `TestLesenMeldungOhneWertUndZeile` |
| Meldung in `pruefeMeta` trägt `Version %s` | rot: `TestLesenEndeVersionen` |
| nur `model.Version` angenommen | rot: `TestLesenEndeAngenommen` |
| Baum trägt fest `"1"` bzw. fest `"2"` | rot: `TestLesenEndeAngenommen`, `TestLesenEndeAbgelehntDarstellung` |
| `+`, Leerraum, führende Nullen vor der Prüfung entfernt | je rot: `TestLesenEndeVersionen`, `TestLesenMeta` |
| Sortierung der Nachrichten entfällt | rot: `TestLesenEndeAngenommen` (beide Fälle „Zeile des Endes vor der Fehlerantwort“), `TestLesenReihenfolgeDerZeilen` |
| Sortierung der Interaktionen bzw. Sessions entfällt | rot nur `TestLesenReihenfolgeDerZeilen` (siehe F-3) |
| Aufruf von `Session.ValidateEnd` in `sessionFromDTO` entfällt | rot: `TestLesenEndeAbgelehntStelle` und `TestLesenEndeJeSession` (die zwei Fälle) |
| `beendetVerbindung` vergleicht `FATAL` ohne Beachtung der Schreibweise | rot: nur „Fatal in Kleinschreibung“ |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Zusage „die Meldung zu `PGR-E3002` nennt keine Version“ und, über sie, die Zahlenmenge in `pruefeMeta` (0, 3, -1, 10, 20) hängen allein an der Prüfung auf das Wortlaut-Muster `"Version "+<Zahl>`: `Unmarshal` lehnt dieselben Zahlen mit demselben Code ab, nur seine Meldung trägt die Zahl. Eine Mutation, die die Zahl in anderer Form nennt (`version=3`, `Fassung 3`), oder eine geänderte Meldung in `kopfVorpruefung` lässt die Mutationen an der Menge grün. | `AGENTS.md` §3.10 (je Zusage eine Mutation, die nur sie bricht) und §3.11 | `internal/adapters/driven/recording/sqlite_lesen_ende_test.go` · `strings.Contains(err.Error(), "Version "+c.version)` | ja — `go test -run TestLesenEndeVersionen` mit der Mutation „`Version %s`“ in anderer Schreibweise | Zusage hängt an einem Wortlaut der Meldung |
| F-2 | INFO | Zur Frage (1): Die Zusage „Version 0/3/-1/10/20/Überlauf = `PGR-E3002`“ ist geprüft, der Überlauf-Fall trägt keine falsche Zusage. Nachgefahren: Fällt die Mengenprüfung in `pruefeMeta` weg, bleibt der Code für 0, 3, -1, 10, 20 `PGR-E3002` (so die Anmerkung in §7), und die Tests werden über die Meldung rot. Die Zahl über 64 Bit wird dagegen `PGR-E3003`, denn der Baum, den der Leser an `Unmarshal` gibt, trägt sie so, dass `kopfVorpruefung` sie als nicht kanonisch ablehnt; hier ist `pruefeMeta` die einzige Stelle für `PGR-E3002`, und der Test ist am Code-Ergebnis rot (SPEC-konform: `SPEC-001` nennt die zu große Zahl `PGR-E3002`). Im Plan §7 nennt die Zeile „Überlauf der Umwandlung als angenommen“ eine Umwandlung, die der Code nicht hat (er vergleicht Texte), und die Anmerkung „die Prüfung in `Unmarshal` hält die Menge zusätzlich“ gilt nicht für die Zahl über 64 Bit. | Maintainability (Genauigkeit der Belege in §7) | `docs/plan/planning/in-progress/slice-v1-abschluss-verbindungsende-sqlite.md` · `Überlauf der Umwandlung als angenommen` | ja — Mutation „Mengenprüfung entfällt“ (siehe Tabelle oben) | Beleg beschreibt Mutation anders als der Code sie erlaubt |
| F-3 | INFO | Zur Frage (2): Der grüne Mutant „Richtung `client` in Extended“ ist äquivalent. Es gibt im Diff keinen Zweig, der die Richtung einer Extended-Zeile `connection_closed` prüft; die Zusage trägt die vorhandene Prüfung der Gruppen (Client-Nachrichten vor Server-Nachrichten, `connection_closed` nur als letzte Server-Nachricht). Der Fall „direction client, extended“ belegt den Code `PGR-E3003`, nicht eine eigene Richtungsprüfung; die Begründung in §7 trägt. Die Ordnung nach den Nummern (`sequence` der Interaktionen, `id` der Sessions) hat keinen Fall mit Ende in `sqlite_lesen_ende_test.go` (`rueck` kehrt nur Zeilen von `message` um); sie ist durch `TestLesenReihenfolgeDerZeilen` aus einem Vorgänger belegt (Mutation nachgefahren, rot). | Maintainability | `internal/adapters/driven/recording/sqlite_lesen_ende_test.go` · `direction client, extended` | ja — Mutation „Sortierung der Interaktionen entfällt“ rot nur in `TestLesenReihenfolgeDerZeilen` | äquivalenter Mutant belegt; Sortierung ohne Ende-Fall |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/adapters/driven/recording/sqlite_lesen.go` (Version aus `liesMeta`/`pruefeMeta` mit `strconv.Itoa(model.Version)` und `strconv.Itoa(model.VersionMitEnde)`, kein Literal; Version geht dem Inhalt vor; Meldung ohne Version; kein eigener Aufruf von `Session.ValidateEnd`) | geprüft, ohne Befund (`grep` nach `Value: "1"` und `version != "1"`: kein Treffer) |
| `internal/adapters/driven/recording/sqlite_lesen_ende_test.go` (Fälle isoliert, Ende als `message`-Zeile mit `kind` genau, `direction` `server`, `fields` `{}`; Stelle nach `position`/`group_no`/`sequence`; je Stelle der Liste in §6 ein Fall; Doc-Kommentare der vier Tests gegen die Fälle gelesen: jede genannte Form hat einen Fall) | geprüft, ohne Befund (Fixtures „Fatal in Kleinschreibung“ und „vorletzte Gruppe mit Sync“ färben sich nur durch ihre Regel) |
| `internal/adapters/driven/recording/sqlite_lesen_inhalt_test.go`, `leser_test.go`, `test/integration/sqlite_lesen_e2e_test.go` (Änderung 2 auf 3) | geprüft, ohne Befund (Fälle prüfen weiter „unbekannte Version“, jetzt mit 3; Version 2 steht als angenommen in `TestLesenMeta`) |
| `spec/spezifikation.md` (`SPEC-043` Punkt *Ende des Servers*, `SPEC-001`) | geprüft, ohne Befund (keine Nennung von ADR, Slice oder Welle; Aussagen decken sich mit den Tests; Drift-Log-Zeile nennt Umplanung) |
| `docs/user/abdeckung-unit.md` | geprüft, ohne Befund (Zeilen entsprechen den Deklarationen der Tests; nur Aussagen, die ein Test prüft) |
| Planungsdokumente (Plan §1 bis §8, Roadmap, `welle-v1-abschluss`, Nehmer `-doku`, `-record`, `-replay`, `-upstream`, `-sqlite-schreiber`) | geprüft, ohne Befund (Sendungen mit Kennung des Gebers, Nachzählung in §8 der Nehmer im selben Commit, kein Move mit Inhaltsänderung) |
| Suppression und Hilfsdateien | geprüft, ohne Befund (kein `//nolint` im Diff, keine Hilfsdateien) |
| Größe (Metrik §4: `git diff --numstat fa8486b..HEAD -- internal docs/user`) | geprüft, ohne Befund: 302 Einfügungen plus 28 Löschungen = 330, wie gemeldet; Meldeschwelle 550 und Marke 700 nicht erreicht |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusage hängt an einem Wortlaut der Meldung · Beleg beschreibt Mutation anders als der Code sie erlaubt · äquivalenter Mutant belegt; Sortierung ohne Ende-Fall

## Verdikt

**Merge-blockierend:** nein — kein HIGH und kein MEDIUM; das LOW (F-1) ist ein Hinweis zur Schärfe eines Tests, die Zusagen selbst sind belegt (alle nachgefahrenen Mutationen rot, der Mutant „Richtung `client` in Extended“ nachweislich äquivalent). `make gates` wurde in diesem Lauf nicht gefahren.

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
