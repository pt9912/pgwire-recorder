# Review slice-v1-abschluss-postgres-versionen-matrix

**Rolle:** Reviewer. **Datum:** 2026-10-14.
**Gegenstand:** Diff `978341e..c5ecc5f` (Basis = Übergang nach `in-progress/`; Code-Commits `fc3e987`, `8f97a05`; Plan-Commits `4c4cfcc`, `c5ecc5f`; Architect-Commit `f0a2280` mit den Randform-Entscheidungen) gegen den Plan `docs/plan/planning/in-progress/slice-v1-abschluss-postgres-versionen-matrix.md`, [`SPEC-053`](../../harness/vertraege/SPEC-053-postgres-versionen.md), [`SPEC-052`](../../harness/vertraege/SPEC-052-gate-kette.md), [ADR-0042](../plan/adr/0042-versionsmatrix-der-integrationstests.md), [LH-FA-05](../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [LH-FA-09](../../spec/lastenheft.md#lh-fa-09--reproduzierbares-replay) und `AGENTS.md` §3.2, §3.7, §3.9 bis §3.14.
**Summary:** 1 HIGH · 3 MEDIUM · 3 LOW · 3 INFO.

Vorgehen: Diff, Vertrag, Plan §1 bis §8 und §7 (Belege des Implementers) gelesen. `make test-postgres-versionen-gegenprobe` im Arbeitsbaum gefahren (grün, 86 s). Mutationen eigener Wahl (Kennung M, nicht die des Implementers) in einer Kopie des HEAD-Stands (`git archive`, Temp-Verzeichnis) gegen die ganze Gegenprobe gefahren: 20 Mutanten, davon 4 grün und 1 hängend (F-2 bis F-5); die Kopie danach verworfen, der Arbeitsbaum blieb unberührt. Ein echter Lauf `make test-postgres-versionen` (einmal, nacheinander): Exit 0, 22 s, Versionen 14, 15, 16, 18 grün, danach keine Container, Netze und Volumes des Laufs (`docker ps -a`, `network ls`, `volume ls`); `git status` sauber.

## Findings

### F-1 HIGH

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist), Reviewer-Skill (HIGH: Kommentar beschreibt die verworfene Alternative)
- `pfad`: `tools/test/run-integration-tests.sh` · "im Vordergrund liefe sie erst nach dem Aufruf"
- `befund`: Der Kommentar vor `vordergrund()` begründet den Hintergrundlauf mit der verworfenen Alternative im Konjunktiv ("im Vordergrund liefe sie erst nach dem Aufruf"). Der erste Satz ("damit ein Signal das wait beendet und die Falle EXIT läuft") trägt die Kopplung allein; der zweite beschreibt den Zustand, der nicht gewählt wurde. Dasselbe Muster ist in `AGENTS.md` §3.7 als "Falsch" ausgeführt.
- `verifizierbar`: nein (kein Gate prüft Kommentare)
- `klasse`: Konjunktiv über die verworfene Alternative im Kommentar

### F-2 MEDIUM

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.10 und §3.11 (Zusage braucht Test), `SPEC-053` Punkt 7 ("Es räumt auf, bevor es endet, ob das Signal an das Skript oder an die Prozessgruppe geht")
- `pfad`: `tools/test/run-postgres-versionen.sh` · `kill -TERM "$laeuft" 2> /dev/null || true` (in `beenden`)
- `befund`: Der Kopf des Skripts sagt zu, es gebe `SIGINT` und `SIGTERM` als `SIGTERM` an den laufenden Runner weiter. Ersetzt man das `kill` durch `:`, bleibt die ganze Gegenprobe grün (M1, nachgefahren: die unveränderte Kopie besteht, kein Fall rot): Die Attrappe lässt den Testlauf höchstens 10 s warten und endet dann, und das `wait` in `beenden` fängt das auf. Szenario: `Strg-C` oder `SIGTERM` während Phase 1 einer echten Version; der Runner ist ein Hintergrundjob und ignoriert `SIGINT`, ohne das Weiterreichen läuft er seine Tests (Minuten) zu Ende, und das Skript endet erst danach; die Zusage "Aufräumen" gilt, die Abbruchzusage bricht still. Die Mutationen `p7-int-treiber`, `p7-term-treiber` und `p7-term-runner` ersetzen die Fallen ganz, nicht das Weiterreichen.
- `verifizierbar`: ja (Mutation, `make test-postgres-versionen-gegenprobe`)
- `klasse`: Hälfte einer Zusage ohne Fall und Mutation

### F-3 MEDIUM

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.10, `SPEC-053` Punkt 3 ("Phase 2 (ohne PostgreSQL, gegen die Aufzeichnungen der Phase 1)") und Punkt 8
- `pfad`: `tools/test/run-integration-tests.sh` · `"$docker" stop "$pg" >/dev/null`
- `befund`: Ersetzt man die Zeile durch `:`, bleibt die Gegenprobe vollständig grün (M2, nachgefahren). Die Attrappe behandelt `stop` als Nullaktion und kein Fall prüft, dass der Server vor Phase 2 gestoppt ist. Szenario: Läuft PostgreSQL in Phase 2 weiter, bestehen die Tests "ohne PostgreSQL" (Szenarien 2 und 7) auch dann, wenn sie den Server brauchen; das Gate meldet grün für eine Aussage, die es nicht belegt. Die Zusage ist im Vertrag neu für die Matrix gefasst, die Zeile stammt aus dem Runner.
- `verifizierbar`: ja (Mutation, `make test-postgres-versionen-gegenprobe`)
- `klasse`: Zusage des neuen Vertrags ohne Fall und Mutation

### F-4 MEDIUM

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.10, `SPEC-053` Punkt 5 ("Scheitert die Einrichtung, ist die Version rot") mit Punkt 4 und 6 (Fehler der Einrichtung steht als `(Aufbau)`)
- `pfad`: `tools/test/run-integration-tests.sh` · "psql -U postgres -v ON_ERROR_STOP=1 -q"
- `befund`: Der Exit-Code der Einrichtung hängt an `set -e` im Skript im Container und an `ON_ERROR_STOP=1`. Entfernt man `-v ON_ERROR_STOP=1` (M3) oder die Zeile `set -e` (M4), bleibt die Gegenprobe grün; `p5-einrichtung` und `p5-aufbau-einrichtung` mutieren nur den Exit-Code des `docker exec` (`|| true`), und die Attrappe erzeugt den Fehler über `AT_HBA_ROT` selbst. Szenario: `CREATE ROLE` oder das Umschreiben von `pg_hba.conf` scheitert auf einer Version; `psql` endet mit 0, die Einrichtung gilt als gelungen, und die Version wird erst in Phase 1 rot (Anmeldung scheitert), mit der Stufe "Phase 1" statt "Aufbau". Gegen einen echten Server ist der Pfad in §7 nicht gesehen (die Tabelle "Gegen echte Server gesehen" hat keine Zeile zu Punkt 5).
- `verifizierbar`: ja (Mutation, `make test-postgres-versionen-gegenprobe`)
- `klasse`: Hälfte einer Zusage ohne Fall und Mutation

### F-5 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.10 (Aufzählung zählt je Element: ein Fall und eine Mutation, die nur ihn bricht)
- `pfad`: `tools/test/postgres-versionen-gegenprobe.sh` · "mut p4-version" und Plan §7 · "Die Namen tragen die Kennung des Laufs und die Hauptversion"
- `befund`: Zwei Elemente haben einen Fall, aber keine Mutation, die nur sie bricht. (a) `SPEC-053` Punkt 4 nennt "von der des Eintrags abweicht": `p4-version-abweichend` hat nur `p4-version`, und die färbt alle drei Fälle der Hauptversion. (b) Punkt 7 letzter Satz: `p7-namen` hat keine Mutation in der Gegenprobe oder in der Tabelle von §7. Eigene Mutanten: Kennung des Laufs aus dem Namen (`lauf="pgr-pv-x"`) färbt `p7-namen` und `p7-vor-naechster` (der Fall fängt es, die Zeile fehlt in der Tabelle); die Hauptversion aus dem Namen (`PGR_LAUF="$lauf"`) lässt die Gegenprobe hängen (F-6) statt rot werden.
- `verifizierbar`: ja
- `klasse`: Element einer Aufzählung ohne eigene Mutation

### F-6 LOW

- `kategorie`: LOW
- `quelle`: Maintainability, `AGENTS.md` §3.10 (Mutation muss rot werden, nicht hängen)
- `pfad`: `tools/test/postgres-versionen-gegenprobe.sh` · "until [ -e "$marke" ] && [ -s "$pidf" ]; do /bin/sleep 0.05; done"
- `befund`: Der Signalgeber in `signalfall` wartet ohne Zeitlimit auf die Marke des Testlaufs, und `wait "$wp"` wartet auf ihn. Erreicht der Lauf den Testlauf nicht (Mutant `PGR_LAUF="$lauf-$haupt"` → `PGR_LAUF="$lauf"` im Treiber, M5), hängt die Gegenprobe; nach 27 Minuten von mir abgebrochen (`ps`: nur der Signalgeber lief noch, ohne Kindprozess außer `sleep`; die Ursache ist hier nicht weiter isoliert). Das Gate wird dann nicht rot, sondern endlos. Der Vertrag sagt zu einem Zeitlimit des ganzen Laufs nichts zu (Punkt 14 (2)); es geht um die Gegenprobe selbst.
- `verifizierbar`: ja (Mutation M5)
- `klasse`: Wartebedingung ohne Zeitgrenze im Prüfwerkzeug

### F-7 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.10, `SPEC-053` Punkt 11 ("Ein Wert ohne führende Ziffer macht den Test rot")
- `pfad`: `test/integration/lebendpruefung_e2e_test.go` · "t.Fatalf("server_version %q: %v", v, err)" (in `hauptversion`) und Plan §7 · "ein Mutant, der den `Fatalf` entfernt, ist nicht gefahren"
- `befund`: Die Zusage von Punkt 11 ist neu im Vertrag; die Gegenprobe führt keinen Fall (Kopf: "ohne Fall hier"), und der echte Lauf in §7 mutiert nur die Eingabe (`"x" + server_version`). Der Implementer räumt den fehlenden Mutanten selbst ein. Entfällt der `Fatalf`, liefert `Atoi("")` den Wert 0 und `0 >= 17` ist falsch, also die Erwartung für 14 bis 16; dort bliebe der Test grün. Szenario nicht gefahren (Einzeiler, ein Lauf `make test-postgres-versionen` mit Mutation).
- `verifizierbar`: ja (echter Lauf mit Mutation)
- `klasse`: Zusage des neuen Vertrags ohne Mutationsbeleg

### F-8 INFO

- `quelle`: `AGENTS.md` §3.7 · `pfad`: `tools/test/postgres-versionen-gegenprobe.sh` · "# Punkt 8: der Runner ohne Auswahl wie zuvor."
- `befund`: "wie zuvor" verweist auf abwesenden Text; `SPEC-053` Punkt 8 sagt es als Zusage ("wie vor diesem Vertrag" mit den benannten Fällen). Der Satz trägt keine Information, die der Fall `p8-*` nicht trägt. Keine Aktion nötig.
- `klasse`: Verweis auf abwesenden Text

### F-9 INFO

- `quelle`: Maintainability · `pfad`: `tools/test/postgres-versionen-gegenprobe.sh` · "GRENZE der Gegenprobe: Mit der Attrappe läuft kein PostgreSQL und kein Test"
- `befund`: Die Attrappe antwortet auf jedes `exec` mit einer Version. Mutanten am Text der Abfrage (`SHOW server_version` → `SELECT 1`, M6) und am Start des Servers (`-e POSTGRES_HOST_AUTH_METHOD=trust` entfernt, M7) bleiben grün; das benennt der Kopf als Grenze, und der echte Lauf deckt beide (Version wird gelesen, Server startet). Der echte Lauf dieses Reviews: 14, 15, 16, 18 grün in 22 s, Aufräumen sauber. Keine Aktion.
- `klasse`: Grenze der Attrappe

### F-10 INFO

- `quelle`: `AGENTS.md` §3.13 (Nehmer nachzählen), Plan §1 *Nahtstelle* · `pfad`: Plan §1 · "Nahtstelle für `slice-v1-abschluss-postgres-versionen-zusage`"
- `befund`: Der Diff trägt keine Sendung in `slice-v1-abschluss-postgres-versionen-zusage` ein; die Nahtstelle ist als DoD-Punkt vor `done/` vorgemerkt, die Namen stehen in §7. Keine Adresse neu im Diff, keine Nachzählung fällig; fällig wird sie im Commit, der die Namen einträgt (DoD-Punkt 9, Zählung in §8 des Nehmers).
- `klasse`: Sendung vorgemerkt, nicht eingetragen

## Eigene Prüfung

Mutanten dieses Laufs (Kennung M), in einer Kopie gegen die ganze Gegenprobe:

| M | Mutation | Ergebnis |
|---|---|---|
| M1 | Treiber: `kill -TERM "$laeuft"` in `beenden` → `:` | grün (F-2) |
| M2 | Runner: `docker stop "$pg"` → `:` | grün (F-3) |
| M3 | Runner: `-v ON_ERROR_STOP=1` entfernt | grün (F-4) |
| M4 | Runner: `set -e` im Skript des Containers entfernt | grün (F-4) |
| M5 | Treiber: `PGR_LAUF="$lauf-$haupt"` → `"$lauf"` | Gegenprobe hängt (F-6) |
| M6 | Runner: `SHOW server_version` → `SELECT 1` | grün (F-9, Grenze der Attrappe) |
| M7 | Runner: `POSTGRES_HOST_AUTH_METHOD=trust` entfernt | grün (F-9) |
| — | Treiber: `wait "$laeuft"` in `beenden` → `:` | rot: `p7-sig*` (4 Fälle) |
| — | Treiber: `lauf="pgr-pv-$$"` → `"pgr-pv-x"` | rot: `p7-namen`, `p7-vor-naechster` (F-5) |
| — | Runner: `set -euo pipefail` → `set -eu` | rot: `p6-exit`, `p6-keine-phase2`, `p6-meldung-*`, `p6-zusammenfassung`, `p9-i` |
| — | Runner: `--internal` am Netz entfernt | rot: `p7-namen`, `p7-vor-naechster` |
| — | Runner: Marke `phase 1` entfernt | rot: alle `p3-element-*`, `p3-skip`, `p3-ohne-treffer`, `p3-fehlt-phase2` |
| — | Runner: Marke `phase aufbau` entfernt | grün, äquivalent (die Vorgabe der Stufe im Treiber ist `Aufbau`) |
| — | Treiber: Image der Version → `POSTGRES_IMAGE` der Referenz | rot: `p4-start`, `p6-referenz-nicht`, `p6-reihenfolge` |
| — | Runner: Passwort des Benutzers verändert | rot: `p5-benutzer` |

Lesart: Wird die unveränderte Kopie (mit meiner Mutation in der Wurzel der Kopie) rot, fängt ein Fall die Mutation. Zeilen "Mutation '…': erwartet rot" daneben sind Folgen meiner Mutation in den Kopien der Gegenprobe und kein Befund.

Stichproben aus der Tabelle in §7: `make test-postgres-versionen-gegenprobe` im Arbeitsbaum grün (86 s; 79 `mut`-Zeilen plus 8 Fälle der Schleife `p3-element-*` = 87, wie §7 sagt). Größe in der Metrik von §4: `git diff --numstat 978341e..c5ecc5f -- harness tools test internal spec docs/plan/adr docs/user` ergibt 786 Einfügungen und 37 Löschungen, Summe 823, wie §7 sagt; unter Meldeschwelle 1200 und Marke 1500; `c5ecc5f` ändert nur den Plan.

## Geprüft ohne Befund

- `harness/mk/postgres-versionen.mk`: Liste (14 bis 18, Referenz 17, je Digest 64 Hexziffern), Teilmenge je Phase, Ziel in einer Zeile ohne `$(MAKE)` und ohne `+`, `GATE_CHECKS` für beide Ziele; Referenz-Image gleich `POSTGRES_IMAGE` (Fall `p1-liste-gueltig`, Mutanten `p9-make`, `p9-plus`, `p12-*` rot).
- `tools/test/run-postgres-versionen.sh`: neun Formen der Eingabe vor dem ersten Aufruf von `docker` (Fälle `p2*`, Mutation je Form, plus `p2-vor-docker`), Reihenfolge, Zusammenfassung mit Stufe, Exit-Code, roter Bau ohne weiteren Aufruf, ein Bau je Lauf; Randformen `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS=n` und `-n -i` in `p9-*` (kein Docker, kein Stand, `-i` belegt nichts).
- `tools/test/run-integration-tests.sh` ohne Variablen: volle Suite in Phase 1, `^TestE2EOhnePostgres` in Phase 2, vier Benutzer, keine Prüfung der PASS-Zeilen, Bau einmal (Fälle `p8-*`); Auswahl mit Unterfall (`TestE2ERecordAnmeldung/scram-sha-256`) trifft im echten Lauf; Test mit Unterfällen unter den übrigen Namen der Teilmenge gibt es nicht (`grep t.Run` in `test/integration`: nur `TestE2ERecordAnmeldung` ist gewählt). Hauptversion am laufenden Server gelesen und verglichen (Fälle `p4-version-*`, im echten Lauf gesehen).
- Aufräumen: Fälle für grünes Ende, rote Phase, Aufbau (nicht bereit, nicht startbar) und je `SIGINT`/`SIGTERM` an Skript und Prozessgruppe; Mutanten `p7-container`, `p7-netz`, `p7-volume` rot; nach dem echten Lauf nichts übrig. Verhalten von `make` bei Signalen als akzeptiertes Negativ in `SPEC-053` Punkt 14 (7) benannt.
- Ausnahmeliste der schnellen Kette: `harness/mk/gates-schnell.mk` und `tools/harness/gates-schnell-gegenprobe.sh` führen `test-postgres-versionen` und `test-postgres-versionen-gegenprobe`; die zwei neuen Mutanten (`m-test-postgres-versionen*`) laut §7 gefahren; Namensmuster in `deckung()` ergänzt.
- `harness/README.md` §Sensors: die zwei neuen Zeilen und die geänderte Zeile von `make gates-schnell` stimmen mit Fragment, Vertrag und Gegenprobe überein; Bindungen `ADR-0042` vorhanden; kein Ziel genannt, das das Makefile nicht kennt (`make test-postgres-versionen`, `-gegenprobe`).
- Plan: §1 (Ziel, Abgrenzung mit Begründung je Punkt), §3 (Dateien des Diffs; `test/integration` und `docs/user/abdeckung-*.md` als bedingt, im Diff unberührt), §4 (Marke mit Metrik, 823 nachgemessen), §6 (Randformen vor dem Code, `f0a2280` vor `8f97a05`; keine Randform im Code-Commit neu in §6), §7 (Belege, Läufe, Namen der Nahtstelle), §8 (Zählung: zwei Liefer-Punkte, zwei Schichten) folgen dem Diff. Der Lauf `make gates` in §7 (553 s) steht auf `8f97a05`, dem letzten Code-Commit.
- `harness/vertraege/SPEC-053-postgres-versionen.md`: Punkte 1 bis 14 gegen Fragment, Treiber und Runner gelesen; Entscheidungen des Architect (Stufe "Aufbau", roter Bau, Phase 2, Kanäle, Signal über `make`, "starten") im Code umgesetzt; der Vertragstext nennt weder ADR noch Slice noch Welle.
- Suppression-Verbot (§3.2) und Gate-Lockerung (§3.6): keine `nolint`-Direktive im Diff, `.golangci.yml`, `.a-check.yml` und `.d-check.yml` unverändert; die Referenz fährt weiter `make test-integration` mit der vollen Suite.
- Commit-Regeln (§3.3, §3.5): kein Move im Diff, keine `Accepted`-ADR inhaltlich überschrieben (`ADR-0042` liegt vor der Basis).
- Kein Produkt-Code unter `internal/` und `cmd/` im Diff (Schicht-Abgrenzung aus §1); `make a-check` nicht gefahren, weil nichts unter `internal/` geändert ist.
