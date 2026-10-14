# Verifikation slice-v1-abschluss-postgres-versionen-matrix

**Rolle:** Verifier (Modul 11, frischer Kontext). **Datum:** 2026-10-14. **Spanne:** `978341e` (Übergang nach `in-progress/`) bis `c5c6ec0` (HEAD; Code-Stand `03d7c6a`, danach nur Plan und Review-Report). **Gegenstand:** Slice `slice-v1-abschluss-postgres-versionen-matrix` (Plan unter `docs/plan/planning/in-progress/`) gegen DoD (§2), Plan §1 bis §8, [`SPEC-053`](../../harness/vertraege/SPEC-053-postgres-versionen.md), [`SPEC-052`](../../harness/vertraege/SPEC-052-gate-kette.md), [ADR-0042](../plan/adr/0042-versionsmatrix-der-integrationstests.md) (Accepted), [LH-FA-05](../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) und `AGENTS.md`. Review: `docs/reviews/2026-10-14-review-slice-v1-abschluss-postgres-versionen-matrix.md` (1 HIGH, 3 MEDIUM, 3 LOW; Nacharbeit `03d7c6a`, `31cb509`).

Eine Verifier-Vorlage führt die vendored Baseline nicht; Aufbau wie ein Review-Report, Verdikt je DoD-Punkt. Alle Belege unten sind in diesem Lauf selbst erzeugt, nicht aus §7 übernommen. Mutationen liefen ausschließlich in Kopien von HEAD (`git archive`, Scratchpad); der Arbeitsbaum blieb unberührt (`git status` sauber, vom Report abgesehen). Läufe nacheinander, kein `pkill`; nach den Läufen 0 Container, 0 Netze, 0 Volumes des Laufs (`pgr-pv-*`; die zwei Volumes `pgr-anm-cache` und `pgr-gocache` sind Caches anderer Läufe).

**Verdikt: PASS mit Befunden.** Kein DoD-Punkt der Lieferung scheitert. 0 HIGH, 1 MEDIUM, 2 LOW, 1 INFO. Die Closure-Pflichten (Notiz, Register, Risiko-Ausgänge, Paarungen, Nahtstelle im Nehmer) sind offen, wie vor der Closure erwartet.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren) |
|---|---|---|
| 1: Bereich 14 bis 18 bestätigt, Liste im Fragment, `make test-postgres-versionen` an `make gates`, Szenarien 1, 2, 7 und Fall SCRAM-SHA-256 je Version außer 17 mit Digest, README, Ausnahmeliste, Gegenprobe an `GATE_CHECKS`; bewusst gebrochen: Image ohne Digest, fehlende Version, Auswahl ohne Treffer | bestätigt, Vorbehalt V-1 | `harness/mk/postgres-versionen.mk`: Referenz 17, fünf Einträge mit 64 Hexziffern, `GATE_CHECKS += test-postgres-versionen test-postgres-versionen-gegenprobe`, Rezept eine Zeile ohne `$(MAKE)` und ohne `+`. `make test-postgres-versionen` Exit 0 in 22 s: 14, 15, 16, 18 grün, danach nichts übrig. `make test-postgres-versionen-gegenprobe` Exit 0 (141 s; 88 `mut`-Zeilen plus 8 Fälle der Schleife `p3-element-*` = 96, wie §7 sagt). Gebrochen und rot gesehen: Mutanten `p2g` (ohne Digest), `p2b` (Referenz fehlt), `p3-treffer` (Auswahl ohne Treffer) färben genau ihre Fälle (`p2g-ohne-digest`, `p2b-referenz-fehlt`, `p3-ohne-treffer`). `harness/README.md` §Sensors: zwei neue Zeilen und die geänderte Zeile von `make gates-schnell` stimmen mit Fragment, Vertrag und Gegenprobe (bash, make, setsid: `setsid` steht in `signalfall`); `GATE_SCHNELL_AUSNAHME` und `LANGSAM`, Namensmuster und zwei Mutanten der Gegenprobe der schnellen Kette nennen beide Ziele. Offen: V-1 (Streichen eines Tests aus der Teilmenge bleibt grün). |
| 2: `TestE2EReplayLebendpruefungWiePostgres` gegen jede Version der Liste grün, `\v` diesseits und jenseits der Grenze; bei Abweichung Rückführung | bestätigt | Lauf oben: 14, 15, 16 (unter der Grenze) und 18 grün, 17 über `make test-integration` im Nachweis des vollen Laufs. Nachfahrbarer Aufruf aus §7 (Nahtstelle 4, Version 16 und 17) selbst gefahren: Exit 0, `Version 16 grün`. Bewusst gebrochen (Kopie): Grenze `vtAbVersion = 17` auf `16`: 14 und 15 grün, `16 (Phase 1)` rot aus dem richtigen Grund („PostgreSQL 16.15 beantwortet `\v` mit fehler 42601 …; `\v` Leerraum laut LH-FA-09.a: true“), 18 grün. Keine Abweichung gefunden, keine Rückführung. |
| 3: `make gates` grün | bestätigt (Nachweis, nicht neu gefahren) | `.harness/state/gates-passed.diffsha` (19:37) ist gleich `tools/harness/working-tree-hash.sh` des heutigen Arbeitsbaums (`40c11632…c6c775`); der Nachweis gilt also für diesen Stand. Selbst gefahren davon: Gegenprobe der Matrix, `make test-postgres-versionen`. §7 trägt für den Stand nach der Nacharbeit nur den roten Lauf bei `docs-check` und den Lauf mit `-k`, keine Zeile des grünen vollen Laufs (V-3). |
| 4: Review durchgeführt, Report liegt vor | bestätigt | Report liegt vor, Ausgang jedes Findings unten geprüft. |
| 5 bis 8: Closure-Notiz mit Lerneintrag, Register, Risiko-Ausgänge, drei Paarungen | offen (erwartet) | §7 trägt bisher „Wird bei Closure gefüllt“; kein Häkchen in §2 gesetzt. Handlung des Planners. |
| 9: Nahtstelle in `slice-v1-abschluss-postgres-versionen-zusage` | offen (erwartet), Namen genannt | Der Nehmer führt die Nahtstelle als Adresse ohne Namen (`grep -n` nach der Kennung findet nur den Platzhalter). §7 nennt die Namen: Datei `harness/mk/postgres-versionen.mk`, Variablen `POSTGRES_REFERENZ`, `POSTGRES_VERSIONEN`, `POSTGRES_TESTS_PHASE1`, `POSTGRES_TESTS_PHASE2`, Teilmenge im Wortlaut (gleich dem Fragment, geprüft), Aufruf `make test-postgres-versionen POSTGRES_VERSIONEN='16=… 17=…'` (selbst nachgefahren, grün), Berichtigung der Grenze: keine. Die Namen decken sich mit dem Code. |

## Ausgang jedes Review-Findings (selbst gefahren)

Mutationen in einer Kopie von HEAD gegen die Fälle der betroffenen Gruppe der Gegenprobe (die Prüfung der unveränderten Kopie, `nur` gesetzt).

| Finding | Ausgang laut Plan §7 | Ergebnis | Urteil |
|---|---|---|---|
| F-1 HIGH (Konjunktiv im Kommentar `vordergrund`) | Kommentar indikativ gefasst | gelesen: „Lange Aufrufe laufen im Hintergrund, damit ein Signal das wait beendet …; … endet dieser Runner auf SIGTERM, das jenes Skript weiterreicht“; keine verworfene Alternative mehr | trägt |
| F-2 MEDIUM (Weiterreichen des Signals) | Fälle `p7-weiter-*`, Mutation `p7-weiter` | Mutation M1 (`kill -TERM "$laeuft"` → `:`): rot sind `p7-weiter-INT-gruppe`, `p7-weiter-INT-skript`, `p7-weiter-TERM-skript` | trägt |
| F-3 MEDIUM (Server vor Phase 2 stoppen) | Fall `p3-phase2-ohne-postgres`, zwei Mutationen | M2 (`stop "$pg"` → `:`): rot ist `p3-phase2-ohne-postgres` | trägt |
| F-4 MEDIUM (Einrichtung der Anmeldung) | Fälle `p5-sql-scheitert`, `p5-hba-datei-scheitert`, Mutationen `p5-on-error-stop`, `p5-set-e`, echte Server | M3 (ohne `ON_ERROR_STOP=1`): rot nur `p5-sql-scheitert`; M4 (ohne `set -e`): rot nur `p5-hba-datei-scheitert` | trägt |
| F-5 LOW (Mutationen je Element) | `p4-abweicht`, `p7-namen-lauf`, `p7-namen-haupt` | M5 (`PGR_LAUF="$lauf"`): rot `p7-namen`, endet in unter einer Sekunde statt zu hängen; `p4-abweicht` und `p7-namen-lauf` laufen in der unveränderten Gegenprobe (grün, 96 Mutationen) | trägt |
| F-6 LOW (Signalgeber ohne Zeitgrenze) | Wächter 20 s, Mutation `p7-sig-ohne-testlauf` | M5 endet ohne Hängen (siehe oben); die Mutation `p7-sig-ohne-testlauf` ist Teil der grünen Gegenprobe | trägt |
| F-7 LOW (`Fatalf` ohne Mutationsbeleg) | `TestHauptversionAus` | Mutant `Fatalf` → `_ = err` in einer Kopie, `make test-postgres-versionen` mit `TestHauptversionAus` in Phase 1 (Liste 16, 17): `TestHauptversionAus` rot, vier Meldungen „endet ohne Fatalf“, `ROT — 16 (Phase 1)` | trägt |
| F-8 INFO („wie zuvor“) | Kommentar berichtigt | `grep` nach „wie zuvor“, „liefe“, „wäre“, „würde“ in Fragment, Treiber, Runner und Gegenprobe: kein Treffer | trägt |
| F-9, F-10 INFO | keine Aktion | Grenze der Attrappe steht im Kopf der Gegenprobe; Sendung an den Nehmer vorgemerkt | — |

## Punkte von SPEC-053 (Fall und Mutation)

Maßstab: jede Zusage mit Fall und einer Mutation, die nur ihre Gruppe rot färbt (`AGENTS.md` §3.10, Aufzählung je Element).

| Punkt | Urteil | Beleg |
|---|---|---|
| 1 Liste | bestätigt | Fall `p1-liste-gueltig` (vier Versionen neben der Referenz), Mutation `p1-liste` (Digest mit Großbuchstaben); Referenz-Image gleich `POSTGRES_IMAGE` |
| 2 Eingabe, neun Formen | bestätigt | 17 Fälle `p2*` (a bis i, Digest in fünf Formen), je eine Zeilen-Mutation, dazu `p2-vor-docker` (alle Fälle) |
| 3 Teilmenge | bestätigt, Vorbehalt V-1 | `p3-skip`, `p3-ohne-treffer`, `p3-element-*` (8, umbenannt), `p3-fehlt-phase2` mit Phase-Mutation, `p3-auswahl`; V5 (Phase 2 ohne Prüfung der erwarteten Tests) färbt `p3-fehlt-phase2`; echter Lauf mit `-test.run` trifft alle acht Tests (alle `--- PASS`) |
| 4 Server | bestätigt | Fälle `p4-nicht-bereit`, `p4-60-sekunden`, `p4-start`, `p4-version-*` (abweichend, unlesbar, ohne Ziffer), je eine Mutation und eine Mutation „Stufe Aufbau fällt auf Phase 1“; eigene Mutation V4 (Prüfung der Hauptversion entfällt) färbt alle vier Fälle der Gruppe |
| 5 Anmeldung | bestätigt | Fälle `p5-benutzer`, `p5-methode`, `p5-hba-vor-allgemeiner`, `p5-kein-md5-klartext`, `p5-*-scheitert` (3); eigene Mutation V3 (`PGR_ROLLEN=scram` entfällt) färbt `p5-kein-md5-klartext`; echter Lauf: SCRAM-Fall `TestE2ERecordAnmeldung/scram-sha-256` grün auf 14, 15, 16, 18 |
| 6 Ablauf und Meldung | bestätigt | Fälle `p6-*` (Reihenfolge, Referenz nicht, nie parallel, Phase 1 rot lässt Phase 2 entfallen, alle Versionen, drei Stufen, Zusammenfassung, Exit, ein Bau, roter Bau mit drei Mutationen); Mutation je Zusage in der Gegenprobe |
| 7 Aufräumen | bestätigt, Vorbehalt V-2 | Container, Netz, Volume je eigene Mutation; Fälle Ende grün, rot, Aufbau (nicht bereit, nicht ziehbar), `SIGINT` und `SIGTERM` je an Skript und Gruppe, Weiterreichen (M1), Namen (Kennung, Hauptversion). Nach den echten Läufen nichts übrig. |
| 8 Lauf der Referenz | bestätigt | Fälle `p8-phase1`, `p8-phase2`, `p8-rollen`, `p8-ohne-pruefung` mit vier Mutationen; `make test-integration` im Nachweis des vollen Laufs |
| 9 Aufrufarten | bestätigt | Fälle für `-n`, `-q`, `-t`, `MAKEFLAGS=n`, `-n -i`, `-i`; Mutationen `p9-make`, `p9-plus`, `p9-stand` (siehe V-4) |
| 10, 14 Voraussetzung, Grenzen | bestätigt | akzeptierte Negative ohne Fall, im Vertrag mit Grund; keine Zusage darüber hinaus im README oder Plan |
| 11 Versionsabhängige Aussage | bestätigt | echte Server: Grenze 16 rot (siehe DoD 2); `TestHauptversionAus` mit `Fatalf`-Mutant rot (F-7); Vertrag und README sagen es zu, ohne Fall in der Gegenprobe, wie dort als Grenze benannt |
| 12 Gegenprobe | bestätigt | `p12-gate-checks-ziel`, `p12-gate-checks-gegenprobe`, je Mutation |
| 13 Ausnahmeliste | bestätigt | `GATE_SCHNELL_AUSNAHME` und `SPEC-052` Punkt 2, 4, 13 nennen `test-postgres-versionen`; die Gegenprobe der schnellen Kette führt `m-test-postgres-versionen*`; Fixture der Matrix enthält keine Dokumente unter `docs/` |

## Findings

### V-1 MEDIUM

- `quelle`: DoD 1 (die Matrix fährt Szenarien 1, 2, 7, Fall SCRAM-SHA-256, Lebendprüfungs-Abgleich), [`SPEC-053`](../../harness/vertraege/SPEC-053-postgres-versionen.md) Punkt 3, `AGENTS.md` §3.10
- `pfad`: `tools/test/postgres-versionen-gegenprobe.sh` · `namen1="$(make_variable "$wurzel" POSTGRES_TESTS_PHASE1 …)"` und `tests=($namen1 $namen2)`
- `befund`: Die Gegenprobe leitet die erwarteten Tests aus dem Fragment der Kopie ab; ein Test, der aus `POSTGRES_TESTS_PHASE1` gestrichen wird (Mutation V1: ` TestE2ERecordExtendedPgx` entfällt), lässt die Gegenprobe grün (nichts rot). Das Gate bleibt grün und fährt Szenario 7 nicht mehr; nur das Umbenennen im Go-Code fällt (PASS-Zeile). Der Vertrag nennt die Tests nicht namentlich, nur die Szenarien.
- `verifizierbar`: ja (Mutation V1, `make test-postgres-versionen-gegenprobe` mit Teilmenge in der Kopie)
- `klasse`: Zusage über den Inhalt einer Liste ohne Fall, der Streichen fängt

Nicht blockierend für den Wortlaut der DoD (die dort verlangten Brüche sind rot), aber der Satz „fährt die Abnahmeszenarien 1, 2 und 7“ ist nur gegen Umbenennen belegt. Der Planner entscheidet: ein Fall, der die Namen im Vertrag oder in der Gegenprobe festschreibt, oder ein akzeptiertes Negativ in Punkt 14.

### V-2 LOW

- `quelle`: [`SPEC-053`](../../harness/vertraege/SPEC-053-postgres-versionen.md) Punkt 7 („auch nach … einem Fehler des Aufbaus … je Fall des Endes ein Fall“), Punkt 4 und 5 (fünf Formen des Aufbaus), `AGENTS.md` §3.10
- `pfad`: `tools/test/run-integration-tests.sh` · `exit 1` nach „der Server meldet server_version“ und `SELECT pg_reload_conf();"`
- `befund`: Aufräumen nach einem Fehler des Aufbaus hat Fälle für „nicht bereit“ (`p7-ende-bereit`) und „nicht ziehbar“ (`p7-ende-pull`); für „Hauptversion weicht ab“ (Mutation W1: `exit 1` → `kill -KILL $$`) und „Einrichtung scheitert“ (W2: `|| kill -KILL $$`) bleibt die Gegenprobe grün. Alle Pfade laufen durch dieselbe Falle `EXIT`, die Lücke ist ein Fall, kein Defekt.
- `verifizierbar`: ja (Mutationen W1, W2)
- `klasse`: Element einer Aufzählung ohne eigene Mutation

### V-3 LOW

- `quelle`: DoD 3 und §6 *Laufzeit des Gates* (Closure-Notiz nennt beide Zahlen), Plan §7 *Läufe*
- `pfad`: Plan §7 · „zusammen rund 105 s“ und die Zeile `make gates` nach dem Slice
- `befund`: Der Zuwachs „rund 105 s“ rechnet mit 80 s für die Gegenprobe; sie läuft nach der Nacharbeit 141 s (selbst gemessen; §7 selbst: 151 s) mit 96 statt 87 Mutationen, die Gegenprobe der schnellen Kette 24 s statt 3 s, zusammen rund 175 s. Zudem fehlt in §7 die Zeile des grünen vollen Laufs für den Stand nach der Nacharbeit (der Nachweis liegt vor und passt zum Arbeitsbaum, ist aber nicht genannt).
- `verifizierbar`: ja (`make test-postgres-versionen-gegenprobe`, `make gates`)
- `klasse`: Plan folgt der Nacharbeit nicht

### V-4 INFO

- `quelle`: [`SPEC-053`](../../harness/vertraege/SPEC-053-postgres-versionen.md) Punkt 9
- `befund`: Die Mutation `p9-make` färbt `p9-n`, `p9-q`, `p9-t`, `p9-makeflags-n` und `p9-n-und-i` zusammen; der Fall `p9-i` hat keine Mutation, die ihn allein bricht (er prüft Verhalten von `make`, nicht des Rezepts; mit `set -o pipefail` entfernt wird er rot). Der Vertrag legt „die Mutation ist ein Rezept mit `$(MAKE)`“ selbst fest; keine Aktion.
- `klasse`: Mutation bricht mehrere Elemente gemeinsam, vom Vertrag entschieden

## Plan-vs-Code

| Plan | Code | Urteil |
|---|---|---|
| §1/§3: Fragment, Treiber, Runner mit Auswahl je Phase, Gegenprobe, Ausnahmeliste, README | alle genannten Dateien im Diff; `docs/user/abdeckung-*.md` unberührt (bedingt, `make abdeckung-check` Teil des Nachweises) | deckungsgleich |
| §3: „Prüfung der Eingabe (neun Formen) vor dem ersten Aufruf von docker“, „ein Bau“, „Signale weitergeleitet“ | Treiber: `P2a` bis `P2i` vor `docker build`, ein `build`, `kill -TERM` in `beenden` | deckungsgleich |
| §1: Gegenstand `TestE2EReplayLebendpruefungWiePostgres` je Version; Berichtigung nur nach Rückführung | Testdatei: `hauptversionAus`, `TestHauptversionAus`; kein Produkt-Code (`git diff` unter `internal/` und `cmd/` leer) | deckungsgleich |
| §4 Marke: Summe 1500, Meldeschwelle 1200 | selbst gemessen `git diff --numstat 978341e..HEAD -- harness tools test internal spec docs/plan/adr docs/user`: 899 + 38 = 937 | unter Meldeschwelle |
| §6 Randformen vor dem Code | SPEC-053 vor `fc3e987` im Basisstand, Rückgaben in `f0a2280` vor `8f97a05`; keine Randform im Code-Commit neu | eingehalten |
| §7 Nahtstelle (Namen) | Fragment, Treiber, Aufruf stimmen | deckungsgleich |
| §8 Zählung: zwei Liefer-Punkte, zwei Schichten | Harness (mit `tools/`, `test/integration`) und Entscheidungen ([ADR-0042](../plan/adr/0042-versionsmatrix-der-integrationstests.md)); eine Textdatei (README) zählt zur Schicht Harness | eingehalten |

[ADR-0042](../plan/adr/0042-versionsmatrix-der-integrationstests.md) (Accepted, unverändert nach `Accepted`): Option C umgesetzt (volle Suite gegen 17 unverändert, Teilmenge gegen die übrigen, Liste an einer Stelle, Gate liest die Release-Liste nicht); Index und `Schärft:` stimmen. `harness/README.md` ist wahr (`AGENTS.md` §3.11): kein genanntes Ziel fehlt im Makefile, keine Zusage über Gegenprobe und Gate hinaus.

## Ausgänge der Risiken aus §6 (Kandidaten, Urteil beim Planner)

| Risiko | Kandidat | Beleg |
|---|---|---|
| Versionsbereich unbestätigt | entfallen (bestätigt 2026-10-14), Pflege von Hand (19 erscheint, 14 endet am 2026-11-12) als *weiter offen* im Register | §7 Bereich bestätigt; Grenze 14 des Vertrags |
| Versionsgrenze für `\v` | entfallen für 14 bis 18 (selbst gefahren, Grenze 17 stimmt); 19 *weiter offen* | DoD 2 |
| Teilmenge fällt still leer | entfallen für die Form „Auswahl ohne Treffer, SKIP, umbenannt“; Rest V-1 *weiter offen* oder Fall | Findings V-1 |
| Laufzeit des Gates | Zahlen in der Closure-Notiz (V-3), kein Carveout nötig | V-3 |
| Anmeldung je Image | entfallen: gleiche Einrichtung wirkt auf 14, 15, 16, 18 | echter Lauf |
| Zweite Phase verdeckt | *weiter offen* (geerbt, Meldung nennt die Stufe) | `p6-meldung-phase1` |
| Namen der Nahtstelle | *eingetreten* mit der Kennung des Nehmers, sobald der Planner sie dort einträgt (DoD 9) | §7 nennt die Namen |
| Abweichung hebt die Schichten | entfallen: keine Abweichung gefunden | DoD 2 |

## Gefahrene Sensoren

| Sensor | Ergebnis |
|---|---|
| `make test-postgres-versionen` | Exit 0, 22 s; Version 14, 15, 16, 18 grün; kein Rest an Containern, Netzen, Volumes |
| dasselbe mit `POSTGRES_VERSIONEN` nur 16 und 17 (Aufruf aus §7) | Exit 0 |
| `make test-postgres-versionen-gegenprobe` | Exit 0, 141 s |
| Mutationen in Kopien (Gegenprobe, `nur`) | V1 bis V5, M1 bis M5, W1, W2, Aufbau-Stufe: Ergebnisse oben |
| Mutationen gegen echte Server in Kopien | `vtAbVersion = 16`: `ROT — 16 (Phase 1)`; `Fatalf` entfernt: `TestHauptversionAus` rot |
| Nachweis des vollen Laufs | `gates-passed.diffsha` gleich Hash des Arbeitsbaums; `make gates` nicht erneut gefahren |
