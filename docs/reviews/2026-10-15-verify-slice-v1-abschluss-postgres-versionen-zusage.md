# Verifikation slice-v1-abschluss-postgres-versionen-zusage

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-15. **Spanne:** `510e7ab` (Übergang nach `in-progress/`) bis `d6dd49b` (HEAD). **Gegenstand:** Slice `slice-v1-abschluss-postgres-versionen-zusage` (Plan unter `docs/plan/planning/in-progress/`) gegen DoD (§2), Plan §1 bis §8 (§6: 13 Entscheidungen des Architect; Entscheidung des Nutzers vom 2026-10-15, Frage B), `SPEC-029`, `SPEC-038`, `LH-FA-05.e`, [LH-FA-05](../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [LH-RB-02](../../spec/lastenheft.md#lh-rb-02--einsatzumgebung), [ADR-0042](../plan/adr/0042-versionsmatrix-der-integrationstests.md) und `AGENTS.md` §3.11. Review: `docs/reviews/2026-10-15-review-slice-v1-abschluss-postgres-versionen-zusage.md` (F-1 bis F-7; Nacharbeit `1bc7dc7`, `d6dd49b`).

Eine Verifier-Vorlage führt die vendored Baseline nicht; Aufbau wie der Report der Matrix. Belege sind in diesem Lauf selbst erzeugt. Mutationen liefen in einer Kopie von HEAD (`git archive`, Scratchpad, danach gelöscht); der Arbeitsbaum blieb unberührt. Läufe nacheinander, kein `pkill`; danach 0 Container und 0 Netze des Laufs (`pgr-*`, `impl-zus-*`).

**Verdikt: PASS mit Befunden.** Kein Liefer-Punkt der DoD scheitert an einem Lauf. 0 HIGH, 1 MEDIUM, 1 LOW, 2 INFO. Der MEDIUM-Befund (V-1) verlangt eine Rückgabe an den Architect, bevor die Closure den Satz zu Abnahmeszenario 18 stehen lässt.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren) |
|---|---|---|
| 1: `SPEC-029`, `LH-FA-05.e`, `SPEC-038`, `LH-RB-02`, Historie; jede Aussage durch einen Lauf der Matrix belegt | bestätigt für Liste, Szenarien 1, 2, 7, SCRAM, Lebendprüfung, Klartext, MD5; **Vorbehalt V-1** für Abnahmeszenario 18 und Szenario 16 Satz 3 | Diff gelesen (`git diff 510e7ab..HEAD -- spec`): `SPEC-029` Tabellenzeile nennt 14, 15, 16, 17, 18 und Referenz 17 (nicht „fünf jüngste“), Absatz mit Punkt 1 und 2, `LH-FA-05.e` behält 3.0 und `PGR-E6002`, `SPEC-038` ein Satz ohne Gate-Namen, Historie-Zeile 2026-10-15, Lastenheft Draft 0.1.0 ohne Bump. `make test-postgres-versionen`: Exit 0, 22 s, 14, 15, 16, 18 grün, Phase 1 sechs Tests, Phase 2 zwei. `make test-integration`: Exit 0, 129 s, 127 `--- PASS` auf oberster Ebene, 0 FAIL, 2 SKIP (die beiden `OhnePostgres*`-Tests in Phase 1, laufen in Phase 2), `TestE2ERecordAnmeldung/{scram-sha-256,md5,password}` PASS. Bewusst gebrochen (Kopie): Eintrag `14=` mit dem Image von 16 gibt Exit 2, `server_version '16.15'`, `ROT — 14 (Aufbau)`; Eintrag `14=` mit erfundenem Digest gibt Exit 2, `manifest unknown`, `ROT — 14 (Aufbau)`. Die Zusage „Hauptversion N“ greift also. |
| 2: Handbuch nennt Bereich; vier Stellen zur Anmeldung; je Version und Verfahren ein nachfahrbarer Lauf in §7; Ist-Zustand ohne Verweise | bestätigt | *Voraussetzungen* (Punkt „Geprüft ist das Werkzeug mit PostgreSQL 14 bis 18“), §4 Hinweis, Abschnitt *Voraussetzung* der Passwort-Anmeldung, §6 *Rollen und Rechte*: je SCRAM-SHA-256 14 bis 18, Klartext und MD5 17. `grep` über den Diff von Handbuch und README nach `SPEC-`, `ADR-`, `slice-`, `LH-`, `Review`: 0 Treffer. §7 führt je Version (14, 15, 16, 18 Matrix, 17 `make test-integration`) und je Verfahren (SCRAM in allen; md5 und password nur 17) einen Aufruf. Der Matrix-Aufruf und `make test-integration` sind oben nachgefahren. |
| 3: README nennt je Verfahren dieselben Versionen | bestätigt | README-Zeile 42 bis 43: SCRAM 14 bis 18, Klartext und MD5 17; deckt sich mit den vier Stellen. |
| 4: `make gates` grün | bestätigt (Nachweis), `make docs-check` und `make kopf-check` selbst gefahren | `make docs-check` 770 Dateien, 0 Befunde; `make kopf-check` Exit 0. `make gates` (10 min 25 s) nicht erneut gefahren; der Stand `1bc7dc7` ist belegt, `d6dd49b` ändert nur den Plan (Positivliste). |
| 5: Review durchgeführt, Report liegt vor | bestätigt | Report liegt vor; Ausgang jedes Findings unten geprüft. |
| 6 bis 9: Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, Paarungen | offen (erwartet) | §7 trägt „Wird bei Closure gefüllt“; kein Häkchen in §2. Handlung des Planners. Register-Ablage `BEO-REPO/postgres-versionsliste-altert` liegt vor. |

## Ausgang jedes Review-Findings

| Finding | Ausgang laut Plan §7 | Ergebnis | Urteil |
|---|---|---|---|
| F-1 MEDIUM (ADR-Kennung ungelinkt, `docs-check` rot) | Link gesetzt | `make docs-check` 0 Befunde am HEAD | trägt |
| F-2 LOW (Trefferstand des `grep`) | Prüfung genau gefasst (Diff-Zählung 0, Bestand benannt) | `git diff 510e7ab -- spec \| grep '^+' \| grep -c "harness/vertraege\|make \|test-postgres\|test-integration"` gibt 0; Bestand `make a-check` nicht im Diff | trägt |
| F-3 LOW (§6 „noch nicht geschrieben“) | §6 nachgezogen | §6 sagt „sind geschrieben (Commit `9458483`, `1bc7dc7`)“ | trägt |
| F-4 LOW (Nachweis-Zeile `LH-RB-02`) | enger gefasst auf die Teilmenge | siehe Beurteilung unten; die Fassung ist selbst fehlerhaft für 17 (V-2) | trägt nur teilweise |
| F-5 INFO (Beleg im Scratchpad) | Aufruf im Wortlaut in §7 | Skript aus §7 in den Scratchpad kopiert und gefahren: `psql (PostgreSQL) 17.11`, drei Anmeldungen plus `trust` liefern ihre Zeile, `Treffer Geheimnis: 0`, `4` `sql:`-Einträge; Container und Netz danach entfernt | trägt |
| F-6, F-7 INFO | keine Aktion | `spec/` führt keinen Link auf `harness/vertraege/` (`grep vertraege spec/*.md` leer) | — |

### Beurteilung F-4 (Einengung gegenüber dem Architect-Wortlaut)

Der Architect schrieb (`7d4534f`): „Die Zeile *Nachweis* nennt zusätzlich die Kompatibilitätstests je unterstützter Hauptversion (Text).“ Die erste Fassung des Lastenhefts („die Kompatibilitätstests laufen je unterstützter Hauptversion“) las sich als alle Tests. Die Einengung auf die Teilmenge deckt die Entscheidung in der Sache: Der Architect selbst fasst „unterstützt heißt genau“ als Teilmenge je Hauptversion und volle Suite nur gegen 17 (§6, Punkt 3), und der Plan nennt den Satz „alle Abnahmeszenarien laufen gegen jede“ als weiterreichend als das Gate (§1). Sie verfälscht die Entscheidung des Architect und des Nutzers also nicht; sie ist die Richtung, die `AGENTS.md` §3.11 verlangt. Zwei Mängel bleiben: (a) die enge Fassung steht erst nachträglich in §6 (`d6dd49b`), ohne dass ein Architect-Vermerk sie bestätigt; (b) das Wort „nur“ ist für die Referenz falsch (V-2). Eine Rückgabe ist wegen V-1 ohnehin nötig; beide Punkte gehen in dieselbe.

## Findings

### V-1 MEDIUM

- `quelle`: `AGENTS.md` §3.11 (Zusage nur, was ein Test prüft), Plan §6 *Zusage reicht weiter als das Gate*, DoD 1 („Jede Aussage der Zusage ist durch einen Lauf … belegt“)
- `pfad`: `spec/spezifikation.md` · `SPEC-029`, Punkt 2: „das Abnahmeszenario 18 und das erwartete Verbindungsende beim Einspielen (Abnahmeszenario 16, dritter Satz) sind nur dort geprüft“; Plan §6 *„Unterstützt“ heißt genau* (2) und §7 Zeile „Referenz 17: volle Suite; … Szenario 18, Szenario 16 dritter Satz nur dort“
- `befund`: Die Zeile sagt, beide Szenarien seien gegen 17 geprüft. Es gibt dafür keinen Lauf: Die Slices `slice-v1-abschluss-verbindungsende-upstream`, `-record`, `-replay`, `-vergleich`, `-einspielen`, `-sqlite`, `-format`, `-doku` liegen in `next/`; keiner liegt in `done/`. `grep -rn "Abnahmeszenario 18" test internal docs/user` hat keinen Treffer; der Test `TestE2EPlayVergleichFatal` (Szenario 16, `FATAL` als Abweichung `PGR-E5004`) belegt nicht „erwartetes Verbindungsende“ (Satz 3), und eine Aufzeichnung mit echtem Ende des Servers im Record-Modus (Szenario 18) hat keinen Integrationstest. §7 nennt als Beleg für diese Zeile nur `md5` und `password`; für 18 und 16 Satz 3 steht dort kein Test. Die Ausnahme ist als Entscheidung des Nutzers vom 2026-10-12 und des Architect eingetragen, beschreibt aber einen Zielstand (die Szenarien laufen erst, wenn die Verbindungsende-Slices schließen). Das Handbuch ist nicht betroffen (nennt die Szenarien nicht).
- `verifizierbar`: ja (`ls docs/plan/planning/done | grep verbindungsende` leer; `grep` oben; `make test-integration` zeigt keinen Test dazu)
- `klasse`: Zusage in der Spezifikation nennt Prüfung, die noch nicht existiert
- Rückgabe an den Architect (der Implementer hielt die Vorgabe ein, die Entscheidung liegt dort): die zwei Elemente aus Punkt 2 streichen, solange die Verbindungsende-Slices offen sind, und die Sendung „Abnahmeszenario 18 und Szenario 16 Satz 3 laufen nur gegen die Referenz 17“ an `slice-v1-abschluss-verbindungsende-upstream` und `-vergleich` (jeweils §1 mit Kennung dieses Gebers, §3.13) legen, die sie liefern; oder der Nutzer bestätigt den Zielstand ausdrücklich als Ausnahme von §3.11. Nicht blockierend für Klartext, MD5 und die Matrix.

### V-2 LOW

- `quelle`: `AGENTS.md` §3.11, Quantor „jede … nur“
- `pfad`: `spec/lastenheft.md` · `LH-RB-02` *Nachweis*: „gegen jede unterstützte Hauptversion läuft nur die Teilmenge der Kompatibilitätstests“
- `befund`: Für die Referenz 17 läuft die volle Suite (`SPEC-029` Punkt 2, selbst gefahren: 127 Tests). „Jede … nur die Teilmenge“ ist für 17 falsch, in die untertreibende Richtung. Zugleich steht „Teilmenge“ nicht als Wort in der Spezifikation (`SPEC-029` Punkt 1 nennt sie als Aufzählung).
- `verifizierbar`: nein (Lesart)
- `klasse`: Quantor „jede“ mit „nur“
- Vorschlag an den Architect: „gegen jede unterstützte Hauptversion läuft mindestens die Teilmenge der Kompatibilitätstests, die die Spezifikation nennt“. Draft, kein Versionssprung.

### V-3 INFO

- `quelle`: Quantoren in den neuen Sätzen
- `befund`: Geprüft, ohne Befund außer V-1 und V-2. „Gegen jede Hauptversion der Liste sind geprüft: …“ (`SPEC-029` Punkt 1) deckt sich mit `POSTGRES_TESTS_PHASE1` und `PHASE2` (Szenarien 1, 2, 7, `scram-sha-256`, Lebendprüfung) und für 17 mit der vollen Suite. „Mit PostgreSQL 17 laufen zusätzlich alle weiteren Integrationstests“ (Handbuch) stimmt für `make test-integration` (beide Phasen, die zwei SKIP sind Phasenwechsel). „Andere Hauptversionen und Vorabversionen sind nicht geprüft“: verneinend, wahr. „Klartext … nur dort“ (Matrix fährt nur `scram-sha-256`): wahr.
- `klasse`: —

### V-4 INFO

- `quelle`: Plan §8 / §4
- `befund`: Größe nach Metrik §4, selbst gemessen: `git diff --numstat 510e7ab..HEAD -- spec docs/user README.md` = 3 + 24 + 7 + 23 = 57; Marke 400 und Meldeschwelle 300 nicht berührt; §7 nennt 57 und das stimmt. Die Schicht-Zählung (zwei) und die Liefer-Punkte (drei) stimmen mit dem Diff (kein Produkt-Code, kein Gate, keine ADR). Der Plan `slice-v1-abschluss-protokollrand` (next/) führt `LH-FA-05.e` und `SPEC-029` in der Kopfzeile; er nennt den entfallenen Satz nicht (`grep` leer), keine Anpassung nötig.
- `klasse`: —

## Entscheidungen des Architect (§6, 13 Punkte) gegen den Diff

| Entscheidung | Urteil |
|---|---|
| Wortlaut `SPEC-029` (Liste, nicht „fünf jüngste“, Vertrag-Datei `—`) | eingehalten |
| Matrix als Text, kein Gate, kein `make`, kein Dateiname, kein Link nach `aussen`; Teilmenge in Worten | eingehalten (`grep` leer; Szenarionummern stehen, Testnamen nicht) |
| „Unterstützt“ heißt genau (1) bis (3) | eingehalten bis auf V-1 (Szenario 18, 16 Satz 3); „alle übrigen Szenarien“ durch „volle Integrationssuite“ ersetzt, dokumentiert in §7, enger als vorgegeben |
| `LH-FA-05.e` | eingehalten |
| `SPEC-038` | eingehalten (ein Satz ohne Gate-Namen) |
| `LH-RB-02` Satz und *Nachweis* | eingehalten; *Nachweis* enger gefasst, siehe F-4 und V-2 |
| Kein Bump des Lastenhefts, Historie nur in der Spezifikation | eingehalten |
| Neue Hauptversion und Ende von 14 ohne Datum oder „bis“ | eingehalten (Handbuch und Spezifikation ohne Chronik) |
| Handbuch und README im Ist-Zustand, je Verfahren | eingehalten |
| Acht prüfbare Aussagen und Träger | (a) bis (e), (h) belegt durch Läufe; (g) für Klartext, MD5 belegt, für Szenario 18 und 16 Satz 3 nicht (V-1) |
| Aufruf-Voraussetzung, Aufrufarten | nicht einschlägig; stimmt |
| Vorgabe der Größe | eingehalten |
| Sätze je Verfahren (Nutzer, Frage B) | eingehalten: SCRAM 14 bis 18, Klartext und MD5 nur 17, kein Lauf von Hand |

## Plan-vs-Code(Text)

| Plan | Text | Urteil |
|---|---|---|
| §1/§3: Spezifikation, Lastenheft, Handbuch, README, kein Produkt-Code | `git diff --stat 510e7ab..HEAD`: nur `spec/`, Handbuch, README, Plan, Review-Report | deckungsgleich |
| §1 Sendung: vier Stellen im Handbuch plus README-Satz, je Verfahren | alle fünf Stellen geändert; „geprüft mit PostgreSQL 17“ ohne Verfahren steht nirgends mehr | deckungsgleich |
| §1 Nahtstelle (1) bis (4) | Fragment `harness/mk/postgres-versionen.mk`: `POSTGRES_REFERENZ`, `POSTGRES_VERSIONEN`, `POSTGRES_TESTS_PHASE1/2` im Wortlaut gleich Nahtstelle (2); Aufruf (4) nachgefahren | deckungsgleich |
| §3 Abdeckungs-Tabellen bedingt | `docs/user/abdeckung-*.md` unberührt; kein Test geschrieben | deckungsgleich |
| §6 Randform „kein Wort zu Pooler und anderen Servern“ | `SPEC-029` nennt sie als „nicht zugesagt“; keine Aussage zu Ablehnung oder Annahme | eingehalten |
| Handbuch-Satz „Als Anwendung dienten ein Go-Programm mit pgx 5.11 und … `psql`“ | `go.mod` pgx 5.11.0, `psql` 17.11 nachgefahren; Satz in der Vergangenheit, aber Aussage über den Lauf, keine Chronik der Entwicklung | eingehalten |

## Ausgänge der Risiken aus §6 (Kandidaten, Urteil beim Planner)

| Risiko | Kandidat | Beleg |
|---|---|---|
| Feste Liste altert | *weiter offen* → `BEO-REPO/postgres-versionsliste-altert` (Ablage liegt; Closure hängt Datei an) | Register |
| Satz im Lastenheft | *entfallen* (Satz geschrieben, kein Versionssprung); die Restfrage aus V-2 geht an den Architect | Diff |
| Zusage reicht weiter als das Gate | **nicht entfallen**: V-1 und V-2 sind die Realisierung; *eingetreten* mit Rückgabe an den Architect, Folge-Adresse `slice-v1-abschluss-verbindungsende-upstream` und `-vergleich` (müssen die Sendung in §1 annehmen, §3.13) | V-1, V-2 |
| Handbuch-Satz ohne nachfahrbaren Lauf | *entfallen*: jede genannte Version und jedes Verfahren hat einen Lauf (selbst nachgefahren); kein Satz nennt mehr | DoD 2 |
| MD5 gegen die jüngsten Versionen | *entfallen*: Handbuch und README nennen MD5 nur für 17 | Diff |
| Namen der Nahtstelle fehlen | *entfallen*: Namen stimmen mit Fragment und Aufruf | Plan-vs-Code |

## Gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make test-postgres-versionen` | Exit 0, 22 s, 14, 15, 16, 18 grün |
| `make test-integration` | Exit 0, 129 s, 127 PASS, 0 FAIL, 2 SKIP (Phasenwechsel), `TestE2ERecordAnmeldung/{scram-sha-256,md5,password}` PASS |
| `psql`-Probe aus §7 (`make build`, Image der Referenz, `psqlprobe.sh` aus dem Plan) | `psql` 17.11, vier Anmeldungen, `Treffer Geheimnis: 0`, 4 `sql:`-Einträge |
| Matrix mit `14=` auf dem Image von 16 (Kopie) | Exit 2, `server_version '16.15'`, `ROT — 14 (Aufbau)` |
| Matrix mit `14=` auf erfundenem Digest (Kopie) | Exit 2, `manifest unknown`, `ROT — 14 (Aufbau)` |
| `make docs-check` / `make kopf-check` | 770 Dateien, 0 Befunde / Exit 0 |
| `make gates` | nicht neu gefahren; Nachweis des Implementers für `1bc7dc7` |
| Aufräumen | 0 Container, 0 Netze (`pgr-*`, `impl-zus-*`); Volumes `pgr-anm-cache` und `pgr-gocache` sind Caches anderer Läufe |
