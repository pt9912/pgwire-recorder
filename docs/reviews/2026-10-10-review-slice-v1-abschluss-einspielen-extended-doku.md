# Review-Report: slice-v1-abschluss-einspielen-extended-doku — 2026-10-10

**Review-Art:** Code-Review (Dokumentation). Geprüft wird gegen Plan §1, §3, §6 und §7, gegen `AGENTS.md` §3.9, §3.11 (inkl. *Handbuch und README beschreiben den Ist-Zustand*), §3.12, §3.13, gegen [ADR-0012](../plan/adr/0012-extended-query-gruppen.md) und [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md) und gegen `LH-FA-20`. Gegen die DoD prüft dieses Review nicht, das ist Aufgabe des Verifiers.

**Gegenstand:** Diff `93382d1..8e72827`: Randformen und Probe-Matrix `2e010b5` (Architect), Handbuch, README und Belege in Plan §7 `880522e`, Gate-Beleg `8e72827` (Implementer). Geändert sind `docs/user/benutzerhandbuch.md`, `README.md` und der Slice-Plan.

**Skill:** `.harness/skills/reviewer.md` @ `0e61bed`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-10

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-v1-abschluss-einspielen-extended-doku` (§1, §3, §6 mit den Randformen des Architect, §7 *Belege des Implementers*)
- `slice-doku-ist-stand` §6 und §7 (Randformen und Befehlsfolge für Doku-Slices)
- [ADR-0012](../plan/adr/0012-extended-query-gruppen.md), [ADR-0017](../plan/adr/0017-einspielen-sequenziell-und-fehlersemantik.md)
- `LH-FA-20`, `LH-FA-18`
- `AGENTS.md` (Hard Rules)

**Selbst gefahren** (gebautes Binary `sha256:b8abba88772c`, `make build`; gepinntes PostgreSQL 17 aus `harness/mk/integration.mk`, Anmeldung `trust`; `psql` `\bind` aus demselben Image; Scratch-Präfix `rev-extdoku-`, Container, Netz und Dateien danach entfernt): Probe 1 (Extended aufzeichnen, einspielen: Exit 0, Zeilen vorhanden; ohne Tabelle `PGR-E4004` Exit 4); Probe 2 (Fehler in Gruppe 1, spätere Gruppe der Interaktion nicht ausgeführt, ohne Option Exit 4 und nächste Interaktion nicht gelaufen, mit `--continue-on-error` nächste Interaktion gelaufen, Exit 4); Probe 3 (`--allow-recorded-errors`: Fehlerantwort nur in der Sync-Gruppe aufgezeichnet, Fehler in der Flush-Gruppe: Exit 0; Fehlerantwort nur in der Flush-Gruppe aufgezeichnet: Exit 0; ohne Option Exit 4; ohne aufgezeichnete Fehlerantwort Exit 4; anderer SQLSTATE `XX000` Exit 0); Probe 4 (`SIGTERM` in `pg_sleep` der ersten Interaktion: Exit 0, Zeile der Interaktion vorhanden, nächste Interaktion und zweite Sitzung fehlen; mit `--finish-session-on-interrupt` zweite Interaktion der Sitzung vorhanden, zweite Sitzung fehlt; Sleep in der Flush-Gruppe: Sync-Gruppe der Interaktion läuft noch; zweites Signal nach 0,5 s: Ende 0,73 s nach dem ersten, Exit 0, kein `level=ERROR`, mit und ohne Option); Probe 5 (`COPY … FROM STDIN` und `COPY … TO STDOUT`, mit und ohne `--continue-on-error`: `PGR-E6001`, Exit 6, Meldung nennt Sitzung und Interaktion, Zeilen der Folge vor dem Copy vorhanden, die danach nicht); zusätzlich `FATAL` (57P01) in einer Extended-Interaktion: `PGR-E4003`, Exit 4, auch mit `--allow-recorded-errors` und `--continue-on-error`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-587 | LOW | Der neue Satz zu `--allow-recorded-errors` bei Extended steht unmittelbar vor dem unveränderten `FATAL`-Satz; dieser liest sich damit auch für Extended. Plan §6 und §7 nennen für `FATAL` bei Extended keine Probe und keinen Test (§7: „eine Extended-Probe ist nicht gefahren“). Die Aussage ist wahr (Probe dieses Reviews: `PGR-E4003`, Exit 4), die Belege des Implementers tragen sie nicht. | `AGENTS.md` §3.11 (Handbuch sagt nur zu, was ein Test oder eine Probe prüft) | docs/user/benutzerhandbuch.md · „Fehlerantwort an irgendeiner Stelle der aufgezeichneten Folge. Ein Fehler mit dem" | ja — Probe oder E2E-Test mit `FATAL` in einer Extended-Interaktion | Handbuch-Satz liest sich weiter als sein Beleg |
| F-588 | INFO | Einordnung des Fundes „Interaktion“ gegen „Nummer der Anfrage“ bestätigt: `PGR-E4004` und `PGR-E6001` nennen bei Extended „Session 1, Interaktion 2“, das Handbuch (Hinweis zu `PGR-E6001`, Katalogzeile `PGR-E4004`) sagt „Nummer der Anfrage“. Das Handbuch beschreibt das Binary richtig genug (die Nummer ist die der Aufzeichnung); ein Begriff „Interaktion“ steht im Handbuch nicht. Entscheidung gehört dem Planner. | Maintainability | docs/user/benutzerhandbuch.md · „die Meldung nennt Sitzung und Nummer der Anfrage" | nein | Begriff der Meldung weicht vom Begriff des Handbuchs ab |
| F-589 | INFO | Einordnung des Fundes README „Slices“ bestätigt als kein Verstoß: `AGENTS.md` §3.11 verbietet den Verweis auf Slices im Handbuch (dort keiner, `grep` ohne Treffer), dem README verbietet sie Chronik und Zielstand. Die Zeile erklärt die Ausgabe von `make doc-trace` und trägt weder Stand noch Chronik; sie steht unverändert im Bestand. Ob das README solche Wörter führen soll, ist Planner-Sache. | `AGENTS.md` §3.11 | README.md · „(Anforderung, Entscheidungen, Slices)" | nein | — |
| F-590 | INFO | Fund „spielt nichts ein“ war falsch: richtig eingeordnet und im Diff behoben. Probe dieses Reviews: bei Copy laufen die Folgen vor dem Copy, danach nichts (Exit 6, auch mit `--continue-on-error`); der ersetzte Satz und die neue Fassung („kann `play` ihn nicht verarbeiten und endet mit `PGR-E6001`“) stimmen mit dem Binary überein. | Maintainability | docs/user/benutzerhandbuch.md · „kann `play` ihn nicht verarbeiten und endet mit `PGR-E6001`" | ja — Probe 5 | Grenz-Satz stand weiter als das Binary (behoben) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Regel Handbuch (kein Verweis auf Spezifikation, Lastenheft, ADR, Slice, Welle, Review; keine Chronik, kein Zielstand): `grep` über das Handbuch | geprüft, ohne Befund (Treffer nur Bestand ohne Bezug: „existiert noch nicht“, „kommt“ in der Ruhefrist, „bis zu deren `Sync`“ als Ablaufangabe) |
| Regel README (Ist-Zustand, keine Stand-Aussage, keine Chronik) | geprüft, ohne Befund |
| Nicht belegbare Aussagen (Warten je Antwort, Gegendruck, Serverversion, „übrige Gruppen ohne Warten“, `CopyBothResponse`, Aufzeichnung endet mitten in der Interaktion): Suche in Handbuch, README und Diff | geprüft, ohne Befund (keine dieser Aussagen steht im Text) |
| Probe-Matrix aus Plan §6 (Zeilen 1 bis 5) gegen Binary und PostgreSQL 17 | geprüft, ohne Befund; jede neue Aussage in §4 und §7 des Handbuchs und im README trägt eine bestandene Probe, bis auf F-587 |
| Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen, Codes) | geprüft, ohne Befund (Implementer-Beleg in Plan §7; `play --help` unverändert zum Diff, Optionen im Handbuch gegen Hilfetext gehalten) |
| Bestandsstellen zu `play` im Handbuch (§1 Betriebsarten, Voraussetzungen, §4 Ergebnis und Hinweise, §5 Optionen, §7 Katalog) und im README | geprüft, ohne Befund (kein „einfache Anfragen“, „nichts eingespielt“ oder Ablehnungsgrund „Folge des erweiterten Protokolls“ mehr; Reihenfolge der Meldungen unverändert) |
| Plan folgt dem Diff (`AGENTS.md` §3.9): §1, §3, §6, §7 | geprüft, ohne Befund; DoD 1 und 2 vom Architect vor dem Code nachgezogen, §7 vom Implementer, Randformen vor dem ersten Doku-Commit (`2e010b5`) entschieden (§3.12) |
| Adressen (`AGENTS.md` §3.13): `slice-v1-abschluss-einspielen-anmeldung-doku`, `slice-v1-abschluss-einspielen-tls-doku` | geprüft, ohne Befund (Diff nennt keine Adresse neu oder geändert; Passwort und `sslmode` bleiben unangefasst) |
| Commit-Reihenfolge und Moves (`AGENTS.md` §3.3) | geprüft, ohne Befund (kein Move im Diff) |
| Mutation (`AGENTS.md` §3.10) | entfällt wie in Plan §6 entschieden (kein neuer Vertrag, nur Beschreibung) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Handbuch-Satz liest sich weiter als sein Beleg · Begriff der Meldung weicht vom Begriff des Handbuchs ab

## Verdikt

**Merge-blockierend:** nein — kein HIGH oder MEDIUM. F-587 ist LOW: die Aussage ist wahr, es fehlt der Beleg im Plan.

**Übergabe:** F-587 geht an den Implementer (Probe nachfahren und in §7 belegen oder den Satz auf den belegten Umfang fassen; Entscheidung beim Implementer). F-588 und F-589 gehen an den Planner zur Kenntnis, ohne Auftrag an diesen Slice. F-590 ist eine Bestätigung und geht an den Verifier zur Kenntnis. Die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist ein Lauf-Beleg und ersetzt keine Verifikation.
