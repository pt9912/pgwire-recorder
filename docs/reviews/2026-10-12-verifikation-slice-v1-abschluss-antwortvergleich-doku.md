# Verifikation: slice-v1-abschluss-antwortvergleich-doku — 2026-10-12

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-12-review-slice-v1-abschluss-antwortvergleich-doku.md`, F-673 bis F-678). Ich repariere nichts.

**Gegenstand:** Slice-Plan `slice-v1-abschluss-antwortvergleich-doku` (Kopf, DoD, §1, §3, §6, §7, §8) gegen `6c4d53c..e0acab9` für `docs/user`, `README.md` und den Plan. Commits: `ef0c823`, `6f10a41`, `485aa94`, `029dad5`, `934c48a`, `d88b5e2`, `0e70ecb`, `564294d`, `59eab7c` (Review-Report), `f106e12` (Nacharbeit Handbuch und README), `baafed2`, `e0acab9` (Belege). Produkt-Code, Tests, Gates und Hilfetexte sind im Bereich unverändert (`git diff --stat 6c4d53c..HEAD`: README, Handbuch, Plan, Review-Report).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-12

**Eingang:** Plan ganz (§7 *Belege des Implementers*, *Nacharbeit nach dem Review*), Diff von Handbuch und README, `LH-FA-24.a` mit *Randformen des Vergleichs*, Review-Report. Den Bericht des Implementers habe ich nicht gelesen.

## 1. Selbst gefahrene Sensoren

Binary: `make build` am Stand `e0acab9` (Image `sha256:63da44cc86ce…`; die Kennung weicht von `e6c5ef3178…` in §7 ab, der Quellstand des Produkts ist derselbe, `6c4d53c`). Proben im Scratch-Verzeichnis außerhalb des Repos; PostgreSQL `postgres:17-alpine@sha256:b0f9560a…` (Pin aus `harness/mk/integration.mk`, Anmeldung `trust`, Netzalias `postgres`) in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt (`docker ps -a --filter name=vfy` leer).

| Sensor / Probe | Ergebnis |
|---|---|
| `make docs-check` | 655 Dateien, 0 Befunde |
| `make abdeckung-check` | Exit 0 |
| Beispiel *Vorgehen*, Option, wörtlich (`--upstream postgres:5432 --input ./recordings/users.yaml --compare-responses`, aufgezeichnet mit `record` und `psql`) | Exit 0, `play beendet eingespielt=3 verglichen=3 abweichend=0`, Form wie im Handbuch |
| Variable `PGWIRE_RECORDER_COMPARE_RESPONSES=true` ohne Option; Schlüssel `compare_responses: true` unter `play:` | beide Exit 0, `verglichen=3`; `config show` gibt die Datei wie im Handbuch aus |
| Abweichung: Aufzeichnung gegen Datenbank mit anderer Spaltenzahl (`d_cnt`) | Exit 5; beide Zeilen wie im Handbuch (`code=PGR-E5004 … Session 1, Interaktion 1: … Position 1, Art Spaltenbeschreibung`; `play beendet eingespielt=1 verglichen=1 abweichend=1`) |
| Ohne `--continue-on-error` endet `play` nach der Abweichung; mit der Option beide abweichenden Interaktionen gemeldet | `eingespielt=1` bzw. `eingespielt=4 verglichen=4 abweichend=2`, Exit 5 |
| Gleicher SQLSTATE (`1/0` aufgezeichnet und beim Server), ohne `--continue-on-error`, ohne `--allow-recorded-errors` | Exit 0, `eingespielt=4 verglichen=4 abweichend=0`; ohne `--compare-responses` derselbe Lauf `PGR-E4004`, Exit 4 |
| Fehler ohne Vorbild (aufgezeichnetes `select 2` zu `select 1/0` geändert) | mit Vergleich `PGR-E5004`, Position 1, Art `Nachrichtenart`, Exit 5, auch mit `--allow-recorded-errors`; ohne Vergleich `PGR-E4004`, Exit 4 |
| Werte und Zeilenzahl nicht verglichen (drei Zeilen aufgezeichnet gegen zwei andere Zeilen) | Exit 0, `verglichen=4` |
| `--compare-responses=false`; `--log-level warn` mit und ohne Abweichung | `verglichen=0` (Exit 4 durch `PGR-E4004`, derselbe Lauf wie ohne Option); mit Abweichung nur die Zeile `PGR-E5004`, ohne Abweichung keine Ausgabe |
| `FATAL` (handgeschriebene Aufzeichnung `select pg_sleep(8)` ohne Fehler; `pg_terminate_backend` aus zweiter Sitzung) | mit Vergleich Exit 5, `Position 2, Art Fehler (SQLSTATE aufgezeichnet „“, Server „57P01“)`, auch mit `--allow-recorded-errors --continue-on-error`; ohne Vergleich `PGR-E4003`, Exit 4 |
| Verbindungsverlust ohne Fehlerantwort nach einer Abweichung (`--continue-on-error`, `docker kill` des Servers in `pg_sleep(8)`) | Exit 4; Zeilen `PGR-E4003` und `PGR-E5004`; `eingespielt=2 verglichen=1 abweichend=1` |
| `SIGTERM` in `pg_sleep(4)`: ohne Abweichung, mit `--finish-session-on-interrupt`, mit Abweichung in der laufenden Interaktion | Exit 0 mit `eingespielt=1 verglichen=1`; Exit 0 mit `eingespielt=2 verglichen=2`; Exit 5, `Position 2, Art Befehl` |
| `record`, `replay`: Option, Variable, Schlüssel | Option `PGR-E2001` (`flag provided but not defined`), Variable `PGR-E2001` (`gehört zu einer Option, die … nicht kennt`), Schlüssel `record.compare_responses`/`replay.compare_responses` `PGR-E2004` „unbekannter Schlüssel“, alle Exit 2; `compare_responses: yes` bei `play`: `PGR-E2004`, `erlaubt sind true und false`; Variable `vielleicht`: `PGR-E2001`, ebenso |
| `--help`: `grep -c compare` bei `record`, `replay`, `play`, `config show` | 0, 0, 1, 0 |
| Befehlsfolge aus `done/slice-doku-ist-stand` §7, Variablen-Zeile mit `grep -vx PGWIRE_RECORDER_PASSWORD` | Optionen: `--h`, `--help`, `--name`, `--rm`, `--version` (Nicht-Optionen und `docker run`); Variablen: leer; Codes: leer. Ohne die benannte Ausnahme nennt die Variablen-Zeile genau `PGWIRE_RECORDER_PASSWORD` |
| `grep -n -i` im Handbuch auf `spec/`, Spezifikation, Lastenheft, ADR, SPEC-, LH-, Slice, Welle, Review | nur Wortteile: „Adresse“ (Treffer auf `adr`, kein Verweis); kein Verweis |
| `grep` im Handbuch auf „noch nicht“, „kommt“, „geplant“, „bisher“, „später“, „künftig“ | nur Bestand: Zeile 14 „später“ (Einleitung), 164 „existiert noch nicht“ (Zieldatei), 515 und 1035 „kommt“ (TLS, Passwort), 908 „Ob sie kommt“ (Ruhefrist); keine Zeile im Diff |
| `grep` im README auf Chronik- und Zielstand-Wörter | kein Treffer im Diff |
| `PGWIRE_RECORDER_BOGUS` und `PGWIRE_RECORDER_PASSWORD` je bei `record`, `replay`, `play` ohne Optionen | alle sechs Exit 2, nur `Pflichtoption … fehlt`, keine Meldung zur Variable |
| COPY (handgeschriebene Aufzeichnung `copy t to stdout`) und abgelehnte Anmeldung (`--user niemand`), je mit und ohne `--compare-responses` | identisch: COPY `PGR-E6001`, Exit 6; Anmeldung `PGR-E4005`, Exit 4 |

Nicht selbst gefahren: `make gates` (Auftrag; der Auftraggeber nennt es grün an `e0acab9`, Exit 0, 19:31:40–19:39:23; die Belege in §7 nennen `d88b5e2` und `baafed2`); `PANIC` (nur als Unit-Test `TestVergleichFatal`, `TestPlayFehlerantwort`, nicht gelesen oder gefahren); `FATAL` mit aufgezeichnetem Fehler an Position 2 oder 1; Befehl-Tags mit anderen Zahlen und `DELETE 1` gegen `UPDATE 1`; Transaktionsstatus; Parameterbeschreibung im erweiterten Protokoll; Art und Reihenfolge im erweiterten Protokoll (`\bind`); Meldungen über Einstellungen und `NOTICE`; Signal in der Wartezeit von `--keep-timing`; zweites Signal; alle Zeilen der Probe-Tabelle in §7, die editierte Aufzeichnungen nutzen (6 bis 10, 14, 16, 21, 23 zum Teil); TLS, Passwort- und SCRAM-Anmeldung; Docker-Netz mit anderer Plattform als Linux. Auch die Sätze, die ich nicht nachgefahren habe, haben je eine Zeile in der Tabelle von §7 und widersprechen keiner meiner Proben.

## 2. DoD Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1 Handbuch im Ist-Zustand: Optionen gegen `--help`, Codes im Katalog und ausgelöst, Beispiele als Datei, keine Verweise, keine Chronik; `--compare-responses` mit Variable und Verhalten beim anderen Kommando | erfüllt | Option, Variable und Schlüssel am Binary bestätigt; `PGR-E5004` im Katalog und ausgelöst (Exit 5); beide Beispielblöcke (Option, Konfigurationsdatei) und beide Ergebnis-Zeilen laufen wie gedruckt; `PGR-E2001` bei `record`/`replay` für Option und Variable, `PGR-E2004` für den Schlüssel; Befehlsfolge grün; `grep` ohne Verweis und ohne neue Chronik. Der Satz zu Variablen mit Präfix nennt `BOGUS` nicht mehr (Befund (a) behoben, keine zweite Ausnahme). |
| 2 `README.md` im Ist-Zustand, ohne Chronik; überholte Sätze ersetzt; Verweise nach oben ohne Stand | erfüllt | Sätze zu `play` (Vergleich, `PGR-E5004`, Exit 5, COPY mit „Ohne `--compare-responses`“) und *Kerngedanke* ergänzt; kein Chronik-Wort im Diff; die Aussage „Exit-Code 5 mit `PGR-E5004`“ ist durch meine Probe belegt. Siehe V-302 zur Lesart des COPY-Satzes. |
| 3 `make gates` grün | erfüllt laut Auftraggeber, selbst nicht gefahren | §7: `d88b5e2` Exit 0 (5 min 0 s) und `baafed2` Exit 0; `e0acab9` ändert gegenüber `baafed2` nur die Zeile zum Stand der Gates (`git diff baafed2 e0acab9`: 2 Zeilen). §7 nennt den Lauf am letzten Commit `e0acab9` nicht; er steht im Auftrag. |
| 4 Review, Report unter `docs/reviews/` | erfüllt | `2026-10-12-review-slice-v1-abschluss-antwortvergleich-doku.md`, 0 HIGH, 0 MEDIUM, 2 LOW, 4 INFO |
| 5 Closure-Notiz mit Lerneintrag | offen (führt die Closure) | §7 nennt keine Notiz |
| 6 Beobachtungs-Register fortgeschrieben | offen (Closure) | |
| 7 Jedes Risiko aus §6 mit Ausgang | offen (Closure) | beide Risiken „offen bis Closure“ |
| 8 Drei Paarungen | offen (Closure) | |

## 3. Probe-Tabelle in §7 gegen die Verhaltenssätze

Ich bin die Sätze des neuen Abschnitts *Antworten beim Einspielen vergleichen*, des Vorspanns der *Hinweise*, der beiden Exit-Code-Sätze, der Zeile 5 der Exit-Code-Tabelle, der Zeilen `PGR-E2004`, `PGR-E4004`, `PGR-E5004`, des Satzes zu Variablen und der README-Sätze durchgegangen. Jeder Satz mit Verhaltensaussage hat eine Zeile in der Tabelle (1 bis 39) oder, wo keine Probe stand, die ausdrückliche Zeile 38/39 mit „keine Probe, enger gefasst“. Kein Satz ohne Zeile gefunden. Die Probe der Zeile 36/37 (Variablen) habe ich wiederholt, ebenso die Zeilen 1 bis 4, 11, 17 bis 20, 22, 23, 25 bis 31, 34.

## 4. AGENTS.md §3.11

| Forderung | Urteil |
|---|---|
| nur Ist-Zustand des gebauten Binaries, keine Chronik, kein Zielstand | erfüllt (`grep`, Proben) |
| keine Verweise auf Spezifikation, Lastenheft, ADRs, Slices, Wellen, Reviews im Handbuch | erfüllt |
| README: Verweise nach oben ohne Stand | erfüllt |
| Beispiele als Datei gegen das Binary | beide Blöcke und beide Ergebnis-Zeilen selbst gefahren |
| zu einem aufgezeichneten Verbindungsende sagt das Handbuch nichts (Entscheidung des Nutzers vom 2026-10-12) | erfüllt: `FATAL`/`PANIC` mit Vergleich ist durchgehend `PGR-E5004`, so wie das Binary es zeigt; kein Befund |
| ein unveränderter Satz neben einem neuen liest sich mit | zwei Lesarten offen, V-301 und V-302 |

## 5. Review-Findings F-673 bis F-678

| Finding | Behoben? | Beleg |
|---|---|---|
| F-673 (LOW) Sätze zu COPY, Anmeldung, TLS, `--finish-session-on-interrupt` | ja, enger gefasst | Der Vorspann der *Hinweise* beschreibt die vier Gruppen ausdrücklich als Einspielen ohne `--compare-responses`; README: „Ohne `--compare-responses` endet `play`, wenn die Datenbank mit einem COPY-Datenstrom antwortet“. Die Zusage ist damit nicht weiter als die Probe. Zusätzlich habe ich COPY, Anmeldung und `--finish-session-on-interrupt` mit Option gefahren: Das Binary verhält sich mit der Option wie ohne (siehe V-302). TLS nicht gefahren. |
| F-674 (LOW) Satz zu Variablen mit Präfix, `PGWIRE_RECORDER_PASSWORD` bei `play` | ja | Satz lautet jetzt „… deren Name zu keiner Option gehört und die nicht `PGWIRE_RECORDER_PASSWORD` heißt“; die Befehlsfolge bleibt grün; Probe: `BOGUS` und `PASSWORD` je Kommando ohne Optionen Exit 2 mit `Pflichtoption … fehlt` (die Nutzung von `PASSWORD` bei `play` regelt der Bestand danach, nicht neu gelesen) |
| F-675 (INFO) Klammer zu `Nachrichtenart`, leere SQLSTATE der Aufzeichnung | nicht nachgezogen, ungenannt (Entscheidung des Implementers, begründet in §7) | Die leere SQLSTATE habe ich in der Probe gesehen (`aufgezeichnet „“`); das Handbuch sagt „nennt die SQLSTATE der Aufzeichnung und der Datenbank“, was wahr bleibt. Kein Befund |
| F-676 (INFO) Beleg `PANIC` | ja | Zeile 24 nennt `TestVergleichFatal` und `TestPlayFehlerantwort`; Tests nicht von mir gefahren |
| F-677 (INFO) Probe 6 | kein Handlungsbedarf; ungenannt | |
| F-678 (a) Exit-Code-Sätze | kein Handlungsbedarf laut Implementer | trägt für den Vergleich nicht mehr, siehe V-301 |
| F-678 (b) README *Kerngedanke*, Zählung der Commits | ja | README-Satz „mit `--compare-responses` vergleicht er die Antworten mit der Aufzeichnung“; Zeile zum Stand der Gates nennt `0e70ecb` und `564294d` |

## 6. Plan gegen Diff (AGENTS.md §3.9)

- **Kopf:** `Bezug` (`LH-FA-24`, `LH-FA-20`, `LH-QA-05`) stimmt; `Berührte Spec-Stellen: —` stimmt (kein Eintrag in `spec/` im Diff); `Verantwortlich` gesetzt.
- **§1 gegen Diff:** Abgrenzung eingehalten. Geändert sind Handbuch, README und der Plan; Produkt-Code, Hilfetexte, Tests, Gates, `docs/user/abdeckung-*.md` und `benutzerhandbuch-standard.md` unverändert. Die Sendung „Gebunden durch“ ist Satz für Satz im Handbuch getragen; zum Verbindungsende schweigt es.
- **§3:** Handbuch und README genannt, stimmt. **§6:** Randformen und Risiken folgen dem Diff; keine Randform im Diff entschieden, die §6 nicht nennt.
- **§7:** *Stand* nennt Version 0.5, die geänderten Abschnitte und die Commits `ef0c823`, `6f10a41`; die Nacharbeit nennt `f106e12`. Die Befehlsfolge in §7 stimmt mit meinem Lauf überein. Der Plan nennt Gate-Läufe an `d88b5e2` und `baafed2`, nicht am letzten Commit; der Auftrag nennt `e0acab9` (V-303).
- **§8:** Zählung 2 Liefer-Punkte, 1 Schicht (Nutzerdokumentation) stimmt mit dem Diff.
- **Commits:** keine Moves im Bereich (§3.3); Messages tragen `slice-v1-abschluss-antwortvergleich-doku` und `LH-FA-24`.
- **Roadmap:** Zeile 23 führt „In Arbeit: nichts. Als Nächstes: `slice-v1-abschluss-antwortvergleich-doku`“, obwohl der Slice in `in-progress/` liegt; der Implementer hat sie nach Befund (b) nicht angefasst (Planung). Hinweis an den Planner, kein Befund gegen den Slice.

## 7. Befunde

| ID | Schwere | Befund | Beleg |
|---|---|---|---|
| V-301 | LOW | Zwei unverändert gebliebene Sätze lesen sich mit dem neuen Abschnitt falsch: *Exit-Codes* (Handbuch Zeile 1205): „liefert den Exit-Code des ersten aufgetretenen Fehlers“, und *Fehlerbehebung* (Zeile 1238): „der erste zählt für den Exit-Code“. Der neue Abschnitt sagt: Verbindungsverlust ohne Fehlerantwort ist `PGR-E4003` mit Exit-Code 4, „auch nach einer Abweichung“. Meine Probe: Interaktion 1 weicht ab (`PGR-E5004`), der Server stirbt in Interaktion 2, Exit 4. Wer die Sätze in *Exit-Codes* wörtlich liest, erwartet 5 (die Abweichung trat zuerst auf). Der Reviewer hat es als INFO (F-678 a) geführt, weil es ohne Vergleich schon galt; mit dem Vergleich widerspricht das Handbuch sich jetzt in der Sache, nicht nur in der Lesart. Nicht blockierend: Die Zeile zum Vergleich nennt den Exit-Code richtig. | `docs/user/benutzerhandbuch.md` Zeilen 641–642, 1205, 1238; Probe: `play --compare-responses --continue-on-error`, `docker kill` des Servers: Exit 4, Zeilen `PGR-E4003` dann `PGR-E5004`, `eingespielt=2 verglichen=1 abweichend=1` |
| V-302 | INFO | Das Enger-Fassen (F-673) ist wahr, aber der README-Satz „Ohne `--compare-responses` endet `play`, wenn die Datenbank mit einem COPY-Datenstrom antwortet, mit `PGR-E6001`“ legt einen Gegensatz nahe: mit der Option sei es anders. Das Binary zeigt keinen: COPY mit Option `PGR-E6001`, Exit 6, wie ohne; abgelehnte Anmeldung `PGR-E4005`, Exit 4, wie ohne; `--finish-session-on-interrupt` mit Option setzt die laufende Sitzung fort (`eingespielt=2 verglichen=2`). Die Probe war billig (handgeschriebene Aufzeichnung mit `copy t to stdout`); die Begründung in §7 („Aufwand steht in keinem Verhältnis“) trägt für COPY und Anmeldung nicht. Kein Widerspruch zum Binary, keine Aktion erwartet; ein Satz „gilt auch mit `--compare-responses`“ wäre mit den Proben tragbar. | `README.md` Zeile 32–33; Handbuch Zeilen 485–489; Proben oben |
| V-303 | INFO | §7 nennt für DoD-Punkt 3 den Lauf an `baafed2`, nicht den letzten Commit `e0acab9`; die Auskunft „grün an `e0acab9`“ steht nur im Auftrag. Inhaltlich gleich (`e0acab9` ändert nur die Zeile), aber ein Beleg des Implementers fehlt für den letzten Stand. | Plan §7 *Stand der Gates nach der Nacharbeit*; `git diff baafed2 e0acab9` |
| V-304 | INFO | Hinweis an den Planner: Roadmap-Zeile 23 „In Arbeit: nichts“ passt nicht zum Zustand `in-progress/`. | `docs/plan/planning/in-progress/roadmap.md` Zeile 23 |

## Verdikt

**Der Slice erfüllt die Liefer-Punkte der DoD (1, 2); Punkt 3 (`make gates`) ist nach Auskunft des Auftraggebers erfüllt, von mir nicht gefahren.** Handbuch und README beschreiben das Binary im Ist-Zustand; alle Beispiele und Ergebnis-Zeilen des neuen Abschnitts laufen so, wie sie gedruckt sind; Befehlsfolge grün; `make docs-check` 0 Befunde, `make abdeckung-check` Exit 0; keine Verweise auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews; F-673 und F-674 sind behoben (F-673 enger gefasst). Keine HIGH- und keine MEDIUM-Befunde. Das Schweigen zum aufgezeichneten Verbindungsende ist die vom Nutzer akzeptierte Lücke und kein Befund.

**Vor der Closure erwägen:** V-301 (zwei Sätze zum Exit-Code „erster Fehler“ neben der Aussage „Exit 4 auch nach einer Abweichung“; enger fassen oder die Ausnahme nennen). V-302 bis V-304 sind Hinweise.

**Offen für die Closure (Implementer):** Notiz mit Lerneintrag (Finding-Klassen aus dem Review und V-301), Register (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`: hier nicht eingetreten; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`: Nachbarsatz neben neuem Satz), Ausgang der beiden Risiken aus §6, Paarungen, dann `git mv` nach `done/`. Die Closure habe ich nicht ausgeführt.
