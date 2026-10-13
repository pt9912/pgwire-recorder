# Verifikation: slice-v1-abschluss-sessions-doku — 2026-10-13

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-13-review-slice-v1-abschluss-sessions-doku.md`).

**Gegenstand:** Slice-Plan `slice-v1-abschluss-sessions-doku` (Kopf, §1 bis §8, mit *Entscheidung des Architect, 2026-10-13*) gegen `c7bce82..9b0a79a`. Commits: `1fe572a` (Handbuch 0.6, README), `a4f956f`/`1f34991` (Belege), `7cab1d4` (Review-Report), `597cf48` (Entscheidung des Architect, F-691/F-689), `de603fb` (Nacharbeit Handbuch), `9b0a79a` (Belege der Nacharbeit).

**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-13

**Eingang:** Plan ganz (§7 *Belege des Implementers*, *Nacharbeit*), Handbuch- und README-Diff, Spezifikation (`LH-FA-12.a` *Randformen des Recordings*, `LH-FA-13.a`, `LH-FA-17.a`), Review-Report (nur die Findings-Tabelle). Den Bericht des Implementers habe ich nicht gelesen.

## 1. Selbst gefahrene Sensoren

Alle Läufe gegen `pgwire-recorder:dev`, frisch mit `make build` gebaut (Image `sha256:deb48422487d…`, derselbe Stand wie im Plan; Produkt-Code seit `cd57943` und seit `c7bce82` unverändert: `git diff … -- . ':!docs' ':!README.md'` leer). Datenbank: `postgres:17-alpine@sha256:b0f9560a…` aus `harness/mk/integration.mk` in einem eigenen Docker-Netz, Anmeldung `trust`, Netzalias `postgres`; Container und Netz danach entfernt. Dateien im Scratch-Verzeichnis außerhalb des Repos.

| Nr. | Sensor / Probe | Ergebnis |
|---|---|---|
| 1 | `make docs-check` am Stand `9b0a79a` | 670 Dateien, 0 Befunde |
| 2 | `make abdeckung-check` | Exit 0 |
| 3 | Befehlsfolge aus `done/slice-doku-ist-stand` §7, Variablen-Zeile mit `grep -vx PGWIRE_RECORDER_PASSWORD` | Optionen: nur `--h`, `--help`, `--name`, `--rm`, `--version`; Variablen: leer; Codes: leer. Keine zweite Ausnahme. `--help` nennt beide Optionen bei `record` beziehungsweise `replay`, bei `play` und `config show` keine. |
| 4 | `grep -i` im Handbuch auf Spezifikation, Lastenheft, `ADR`, `SPEC-`, `LH-`, Slice, Welle, Review, `spec/` | kein Treffer (nur `Adresse` auf `ADR`) |
| 5 | `grep -i` auf `noch nicht`, `kommt`, `geplant`, `bisher`, `später`, `künftig`, `vorerst` | nur Bestand in anderem Sinn; neu: Zeilen 526 und 600 (`später begonnene Verbindung`, `später PGR-E4006`), beides Reihenfolge, keine Chronik. README-Diff: kein Treffer. |
| 6 | Beispiele des neuen Abschnitts als Datei: `config show` für beide YAML-Dateien (`record_empty_sessions`, `session_assignment`) | Exit 0, Inhalt wie geschrieben |
| 7 | Startfehler ohne Datenbank: `--session-assignment` `Connection`, `first_request`, leer; Variable `Connection`, `first_request`, `GEHEIM`; Variable leer (gilt als nicht gesetzt); Datei `Connection`, `""` | Option und Variable `PGR-E2001` (Variable nennt ihren Namen, nicht den Wert), Datei `PGR-E2004` (`replay.session_assignment`), Exit 2; leere Variable: kein Fehler |
| 8 | beide Optionen und Variablen beim anderen Kommando; Schlüssel im anderen Abschnitt; `--record-empty-sessions=1`, `=True`, `=`; Variable `vielleicht`, `1`, `True`; Schlüssel `yes`, `1` | Option `flag provided but not defined`, Variable `gehört zu einer Option, die <Kommando> nicht kennt` (auch `false`), alle `PGR-E2001`; Schlüssel im anderen Abschnitt `unbekannter Schlüssel` und falscher Wert `PGR-E2004`; ungültige Werte `PGR-E2001`; Exit 2 |
| 9 | `replay --session-assignment connection` mit einer ohne Option aufgezeichneten Datei, YAML und SQLite | `PGR-E2001` („trägt keine Kennzeichnung empty_sessions“), Exit 2; `first-request` startet mit derselben Datei |
| 10 | `record --record-empty-sessions`: Verbindungen ohne Anfrage, `select 'eins'`, ohne Anfrage, `select 'drei'`; YAML und SQLite | `empty_sessions: true`, vier Sessions (`id` 1 bis 4), leere mit `interactions: []`; SQLite `meta` mit `empty_sessions|true`, `session` 1 bis 4, `interaction` nur Sessions 2 und 4. Kopf der YAML-Datei gleich dem Block im Handbuch. |
| 11 | `replay` der Aufzeichnung aus 10 mit `first-request` und `--fail-on-unconsumed`, YAML und SQLite | `eins` und `drei` beantwortet, Exit 0, keine Meldung; dritte Anfrage `PGR-E5003`, Exit 5 |
| 12 | `replay --session-assignment connection --fail-on-unconsumed`, Verbindungen leer, `eins`, leer, `drei`, YAML und SQLite; Datei und Variable statt Option; Option vor Datei | Exit 0 ohne Meldung; zweite Verbindung `select 'zwei'` an Session 1 ohne Interaktionen: `PGR-E5001`, Exit 5; Variable und Datei wirken, die Option geht vor |
| 13 | `connection`: Lebendprüfung; Verbindung 5 ohne Session | `-- ping` ohne Fehler; Anfrage `PGR-E5003` („keine aufgezeichnete Session mit der Kennung 5“), Exit 5; leere Sessions nie in `PGR-W2001` genannt |
| 14 | Nur leere Sessions (zwei Verbindungen ohne Anfrage), und Aufzeichnung ohne Session mit Kennzeichnung (YAML, SQLite) | `replay` mit beiden Zuordnungen `PGR-E3004`, Exit 3 |
| 15 | Ohne Verbindung aufzeichnen, mit und ohne Option; Verbindung ohne Anfrage, YAML und SQLite, mit und ohne Option | mit Option `empty_sessions: true`, `sessions: []`; ohne: keine Kennzeichnung; Zieldatei (SQLite) nach einer leeren Verbindung nur mit Option da. `--record-empty-sessions=false` und Variable `true` verhalten sich wie beschrieben. |
| 16 | „Keine Session“ (mit Option): `SSLRequest`, `GSSENCRequest`, `CancelRequest`, 1, 4, 7, 8, 20 Bytes, TCP ohne Bytes (Roh-Bytes mit `printf \| nc`) | Antwort `N` auf beide Anfragen nach Verschlüsselung; `PGR-W3001` für den Abbruchwunsch; 1, 4, 7 Bytes und leer: Exit 0, keine Meldung; 8 und 20 Bytes: `PGR-E6001` („Startnachricht nicht lesbar“), Exit 6; stets `sessions: []` |
| 17 | Aufbau scheitert (`--upstream 127.0.0.1:1`); `COPY` auf Verbindung A, Verbindung B offen mit zwei Anfragen; Passwort (SCRAM) | `PGR-E4002`, Exit 4, `sessions: []`; B beantwortet und aufgezeichnet (`- id: 1`), A `PGR-E6001`, Exit 6; Passwort `PGR-E6001`, Exit 6, `sessions: []` mit Kennzeichnung |
| 18 | Zwangsende im Aufbau (stummer Upstream, `--shutdown-timeout 300ms`); `pg_sleep(20)` mit `--shutdown-timeout 500ms`; ruhende Verbindung bei `docker stop` vor der Frist; `COPY` dann `pg_sleep` bis zum Zwangsende | Aufbau: Exit 0, keine Meldung, keine Session (mit und ohne Option); `pg_sleep`: `PGR-E4006`, Exit 4, Session nur mit Option (Meldung nennt „die Session wird als Session 1 geschrieben“); ruhend vor der Frist: Exit 0, mit Option `- id: 1` mit `interactions: []`, ohne `sessions: []`; `PGR-E6001` dann `PGR-E4006`: Exit-Code 6 |
| 19 | `id`-Reihenfolge: A beginnt zuerst und endet 3 s nach B | `- id: 1` mit den Anfragen von B, `- id: 2` mit denen von A |
| 20 | `replay` mit zwei gleichzeitig offenen Verbindungen; B fragt falsch | B `PGR-E5001`, A beantwortet beide Anfragen, Exit 5 |
| 21 | `replay`: `SSLRequest`, `GSSENCRequest`, `CancelRequest`, leere Verbindung, `PGSSLMODE=prefer`, mit beiden Zuordnungen | zählen nicht als Verbindung (folgende Verbindungen erhalten die erwarteten Sessions), `prefer` zählt einmal, Exit 0 mit `--fail-on-unconsumed` |
| 22 | `play` einer Aufzeichnung mit zwei leeren und zwei Sessions mit Interaktionen; Zeilen `connection authorized` im Log der Datenbank | `eingespielt=2`, genau 2 Verbindungen je Lauf |
| 23 | Satz über die erste Anfrage nach COPY: `select 'vorher'`, dann `COPY` auf derselben Verbindung (ohne Option) | `PGR-E6001`, Exit 6, `sessions: []` (Grundlage für V-310) |

**Nicht selbst gefahren:** `make gates` (volle Kette, auf Anweisung; der Stand des Implementers und des Stop-Hooks für `9b0a79a` ist mir nur als Aussage des Aufrufers bekannt); `make lint`, `make test`, `make test-integration`, `a-check` (Produkt-Code unverändert); die Proben mit gleichzeitigem Ende und gleichzeitigem Fehler (die Negativsätze des Handbuchs haben bewusst keine Probe, Gleichzeitigkeit ist nicht deterministisch); das zweite Signal bei ruhender Session (akzeptiertes Negativ, das Handbuch sagt dazu nichts); der gleichzeitige Aufbau bei der Zuordnung; `PGR-E5002` mit `--fail-on-unconsumed` und nie zugeordneter Sitzung als Datei-Probe (nur `PGR-W2001` im Log gesehen, Zeilen 3 und 5 der Proben 12/13); die Tabelle §5 in den Zeilen außerhalb der beiden neuen Optionen; `GSSENCRequest` unter eigenem Namen (nicht im Handbuch); die Lücke des aufgezeichneten Serverendes (`LH-FA-25`): das Handbuch sagt dazu nichts (`grep` auf `FATAL`, `PANIC`, `Ende des Servers`: kein Treffer), laut Aufrufer akzeptiert, kein Befund.

## 2. DoD Punkt für Punkt

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| 1 Handbuch im Ist-Zustand (Abschnitte, Optionen gegen `--help`, Codes im Katalog und ausgelöst, Beispiele als Datei, keine Verweise, Grenz-Sätze ersetzt, Übernahme aus `…umgebung-anderes-kommando`) | erfüllt, mit V-310 (LOW) | Neuer Abschnitt *Mehrere Verbindungen aufzeichnen und wiedergeben*; geänderte Sätze in *Aufzeichnen*, *Wiedergeben*, *Einspielen*, SQLite, *Herunterfahren mit Frist*, §5 (Absatz, Tabelle), Exit-Code 5, Katalog, Glossar. Befehlsfolge grün (1.3), `grep` sauber (1.4, 1.5), Proben 6 bis 23 bestätigen jede Verhaltensaussage, die ich gefahren habe. Übernahme aus `…umgebung-anderes-kommando`: beide Variablen mit Kommandos in §5, Variable beim anderen Kommando `PGR-E2001` (Probe 8). Rest: V-310. |
| 2 README | erfüllt, mit V-310 (gleicher Satz) | Ein Absatz (Zeilen 47 bis 53), Ist-Zustand, kein Zielstand, kein Chronik-Wort; Verweise auf `spec/` und `docs/plan/` unverändert. Der Satz zu `record` trägt dieselbe „jede“-Aussage wie das Handbuch (V-310). Formatierung siehe V-312. |
| 3 `make gates` grün | erfüllt laut Aufrufer, Beleg in §7 nur für `a4f956f` | Mein Lauf: nur `docs-check` und `abdeckung-check` grün (1.1, 1.2); `make gates` habe ich nicht gefahren. §7 nennt Exit 0 für `a4f956f`; die Commits `de603fb`/`9b0a79a` ändern das Handbuch, und für sie nennt §7 nur `docs-check` (V-311). |
| 4 Review, Report unter `docs/reviews/` | erfüllt | Report liegt vor; 0 HIGH, 0 MEDIUM, 6 LOW, 2 INFO. |
| 5 Closure-Notiz mit Lerneintrag | offen (Closure) | §7 bis zur Closure unvollständig, kein Befund |
| 6 Beobachtungs-Register | offen (Closure) | Hinweis: die Klasse „aus einem Nachbarsatz übernommene Zusage ohne seine Voraussetzung/Ausnahme“ (F-686 bis F-689, hier V-310) gehört in den Lerneintrag |
| 7 Risiken mit Ausgang | offen (Closure) | zwei Risiken in §6, beide „offen bis Closure“ |
| 8 Paarungen | offen (Closure) | |

## 3. Entscheidungs-Konformität (Entscheidung des Architect, 2026-10-13)

| Punkt der Entscheidung | Urteil | Beleg |
|---|---|---|
| Handbuch nennt: Ende ohne Daten, Anfrage nach Verschlüsselung (unter diesem Wort), Abbruchwunsch, abgebrochene Startnachricht (ab 8 Bytes `PGR-E6001`, Exit 6; weniger still), gescheiterter Aufbau (`PGR-E4002`), `COPY`, Zwangsende im Aufbau | erfüllt | Handbuch Zeilen 540 bis 551 (*Hinweise zum Aufzeichnen*); je Fall eine Zeile der Probe-Tabelle (17, 19, 20) und meine Proben 16 bis 18 |
| `GSSENCRequest` nicht unter eigenem Namen | erfüllt | `grep GSSENC` im Handbuch: kein Treffer; Probe 16 zeigt die Gleichheit |
| Zwangsende bleibt eng, nur im Aufbau; „oder keine gestellt hatte“ entfällt | erfüllt | Zeile 532 bis 535: „… ihre erste Anfrage noch ausführte (`PGR-E4006`)“; `grep "keine gestellt"`: kein Treffer |
| Ruhende Verbindung vor der Frist als Sitzung nur mit Formulierung „beim Herunterfahren vor Ablauf der Frist“ | erfüllt | Zeile 534: „Endet eine Verbindung ohne Anfrage beim Herunterfahren vor Ablauf der Frist, ist sie mit der Option eine Sitzung.“ Probe 18 bestätigt (Exit 0, `- id: 1`). |
| Akzeptiertes Negativ: nichts zum Verlust beim zweiten Signal | erfüllt | kein Satz im Handbuch dazu |
| Fälle nur, soweit Zeilen der Probe-Tabelle sie decken | erfüllt | die Fälle der Proben 15 bis 20; „weniger als 8 Bytes“ an 1, 4, 7 Bytes (meine Probe 16 bestätigt auch leer) |

## 4. Plan gegen Diff (AGENTS.md §3.9) und `AGENTS.md` §3.11

- **Kopf:** `Bezug` (`LH-FA-12`, `-03`, `-13`, `-02`, `-09`) stimmt mit den Inhalten; `Berührte Spec-Stellen: —` stimmt (kein Eintrag in `spec/` im Diff); `Verantwortlich` gesetzt.
- **§1:** Abgrenzung eingehalten: Diff ändert nur `docs/user/benutzerhandbuch.md`, `README.md`, den Plan und den Review-Report; Produkt-Code, Hilfetexte, Tests, Gates, `abdeckung-*.md`, `benutzerhandbuch-standard.md` unverändert. Die Entscheidung des Architect steht in §1 (`597cf48`), §7 *Rückgaben* zeigt auf sie (F-691 behoben).
- **§3:** Handbuch und README genannt; stimmt.
- **§6:** Randformen und Risiken folgen dem Diff; die Randform „Verbindungen, die keine Session sind“ zeigt auf die Spezifikation, der Diff entscheidet keine weitere.
- **§7:** Probe-Tabelle (42 Zeilen) deckt jeden Verhaltenssatz, den ich geprüft habe; Zeile 41 trägt „gilt für Verbindungen, die … nacheinander aufbaut und nutzt“ (F-690), *Ungenannt geblieben* ist nachgezogen. Ein Widerspruch zwischen Plan und Diff blieb nicht stehen.
- **§8:** Zählung 2 Liefer-Punkte, 1 Schicht (Dokumentation) stimmt.
- **Commits:** Nachrichten nennen Slice-Kennung und `LH-FA-12`, keine `SPEC-`/`ARC-`-Kennung; kein `git mv` im Bereich.
- **§3.11:** keine Verweise nach oben im Handbuch (Probe 1.4); keine Chronik (1.5); Beispiele laufen als Datei (Proben 1.6, 10, 12); Negativsätze (gleichzeitiges Ende, gleichzeitige Fehler, gleichzeitige Nutzung) sind als „nicht zugesichert“ formuliert und im Plan als Negativsätze ohne Probe geführt.
- **Nachbarsätze (Schritt 17):** Merkmal „Allgemeinregel des Bestands, die der neue Abschnitt einschränkt“. Gefunden und trägt: Exit-Code 5, Absatz unter den Exit-Codes, `--fail-on-unconsumed`-Satz, `PGR-E3004`, `PGR-E5001`/`E5002`/`E5003`/`W2001`, Voraussetzung, Hinweis zum Beenden (Zeilen 224 bis 236: „Eine Verbindung ohne begonnene Anfrage … endet ohne Meldung“ bleibt wahr, Probe 18), Glossar, Tabelle §5. Nicht mitgelesen: V-310.

## 5. Review-Findings F-686 bis F-693

| Finding | Behoben? | Beleg |
|---|---|---|
| F-686 Folgefehler ohne Ausnahmen | ja | *Verbindungsfehler* nennt nur das Paar `PGR-E6001` dann `PGR-E4006` (Exit-Code 6, Probe 18) und verweist für Ausnahmen auf *Exit-Codes*; „ein später folgender Fehler einer anderen Klasse ändert ihn nicht“ entfällt (`grep`: kein Treffer) |
| F-687 „erster Fehler“ bei gleichzeitigen Fehlern | ja | Zeilen 597 bis 599 und 1392 bis 1394 tragen „gleichzeitig … nicht zugesichert, welcher der erste ist“; „gemerkt“ im neuen Abschnitt ersetzt (Bestand 234 unverändert) |
| F-688 `id`-Reihenfolge | ja | Zeile 525 bis 528: „Enden die Sitzungen nacheinander …; gleichzeitig: nicht zugesichert“; Probe 19 (3 s Abstand) bestätigt |
| F-689 Zwangsende, „oder keine gestellt hatte“ | ja | Zeile 532 bis 535, siehe Abschnitt 3 |
| F-690 Voraussetzung „nur zugesichert, wenn“ | ja | Zeile 452 bis 453: „gilt für Verbindungen, die Ihre Anwendung nacheinander aufbaut und nutzt“; Plan §7 *Ungenannt* nachgezogen; Probe-Zeile 41 vorhanden |
| F-691 Entscheidung des Architect ohne Artefakt | ja | §1 *Entscheidung des Architect, 2026-10-13*, §7 *Rückgaben* verweist auf `597cf48` |
| F-692 INFO | bleibt, Hinweis | keine Aktion |
| F-693 INFO `PGR-E4006`-Zeile | ja, wie §7 sagt (Satz trägt) | Katalogzeile „nennt die Sitzung und die Anfrage“; die Meldung nennt „Interaktion 1 … Session 1“ (Probe 18, Zeile 42 der Tabelle) |

## 6. Befunde

| ID | Schwere | Befund | Beleg |
|---|---|---|---|
| V-310 | LOW | Eine „jede“-Aussage ohne die Ausnahme des Nachbarsatzes, in derselben Klasse wie F-686 bis F-689: *Hinweise zum Aufzeichnen* (Zeile 215 bis 218) sagt „Jede Verbindung Ihrer Anwendung, die mindestens eine Anfrage abgeschlossen hat, wird als eigene Sitzung aufgezeichnet“; der Einleitungssatz des neuen Abschnitts (Zeile 446) sagt „zeichnet jede Verbindung als eigene Sitzung auf“, das README (Zeile 47) wiederholt den ersten Satz. Das Binary schreibt eine Verbindung, die erst eine Anfrage abschließt und dann `COPY` nutzt, nicht: `PGR-E6001`, Exit 6, `sessions: []`. Der neue Abschnitt nennt `COPY` als Verbindung, die keine Sitzung ist (Zeile 540), und der Passwort-Punkt direkt unter Zeile 215 nennt eine Ausnahme; der Satz selbst nicht. Die Ausnahme ist im Handbuch also vorhanden, nur nicht an der „jede“-Aussage; der ältere Satz (Bestand mit „mindestens einer Anfrage“) trug den Fehler schon, der Diff hat ihn umformuliert und ihn neben die neuen Sätze gestellt. Kein Zielstand-Problem, ein Verhaltenssatz, der eine Probe-Zeile braucht oder eng gefasst werden muss. | Probe 23: `printf "select 'vorher';\nCOPY (SELECT 1) TO STDOUT;\n" \| psql` gegen `record`, Aufzeichnung `sessions: []`; Handbuch Zeilen 215, 446, 540; README Zeile 47 |
| V-311 | LOW | DoD 3 (`make gates` grün) hat in §7 keinen Beleg für den Stand, den der Bericht nennt. Der Beleg *Stand der Gates* nennt `a4f956f`; danach ändern `de603fb` (Handbuch, Nacharbeit, +/− Sätze zu Exit-Code, `id`, Zwangsende, Voraussetzung) und `9b0a79a` (Plan) Dateien, die `docs-check` und `abdeckung-check` lesen, und §7 *Nacharbeit* nennt nur `make docs-check` mit 670 Dateien. Der Satz „`make gates` auf dem Stand nach dem letzten Commit ist die Messung, die der Bericht nennt“ verweist auf einen Lauf, der in §7 nicht steht. Das Urteil „erfüllt“ stützt sich auf die Aussage des Aufrufers (Implementer und Stop-Hook, Exit 0 für `9b0a79a`), nicht auf einen Beleg im Plan und nicht auf meinen Lauf. Dieselbe Klasse wie V-203 der letzten Doku-Verifikation. | `grep -n 'make gates'` im Plan: DoD, §5, *Stand der Gates* (`a4f956f`); §7 *Nacharbeit*: nur `docs-check` |
| V-312 | INFO | README Zeile 47: Der Absatz zu den Sitzungen beginnt mitten in einem Satz-Absatz zu den Zeitangaben und hat unregelmäßige Zeilenumbrüche (Zeile 47 trägt den Schluss des `PGR-E3003`-Satzes und den Anfang des neuen); `docs-check` ist grün, inhaltlich ist nichts falsch. Rein Lesbarkeit. | `git diff c7bce82..HEAD -- README.md` |
| V-313 | INFO | Der Bestandssatz *Hinweise zum Aufzeichnen* „Die abgeschlossenen Anfragen einer solchen Verbindung bleiben in der Aufzeichnung“ (Zeile 229 bis 231) steht neben der neuen Aussage, dass mit der Option auch eine Verbindung ohne abgeschlossene Anfrage nach `PGR-E4006` als Sitzung geschrieben wird; das ist kein Widerspruch (Probe 18 zeigt beide Fälle), nur zwei Sätze, die der Leser verbinden muss. Kein Handlungsbedarf. | Handbuch Zeilen 229 bis 236 und 530 bis 535; Probe 18 |

Offen, aber keine Befunde: Das aufgezeichnete Verbindungsende des Servers (`LH-FA-25`) ist vom Nutzer als Lücke akzeptiert (Entscheidung 2026-10-12); das Handbuch sagt dazu nichts.

## Verdikt

**Der Slice erfüllt die Liefer-Punkte der DoD (1, 2) mit einer Auflage (V-310).** Handbuch und README beschreiben das gebaute Binary im Ist-Zustand; jede Verhaltensaussage, die ich am Binary gefahren habe, stimmt (Proben 6 bis 22, einschließlich der Startfehler, beider Zuordnungen, der Fälle „keine Session“, des Verbindungsfehlers, `PGR-E5003` und `PGR-E6001` mit Exit 6); die Entscheidung des Architect ist umgesetzt, das Zwangsende bleibt eng; die Befehlsfolge ist grün; keine Verweise nach oben, keine Chronik; F-686 bis F-691 sind behoben. Keine HIGH- und keine MEDIUM-Befunde.

**Vor der Closure:** V-310 (Satz „jede Verbindung mit abgeschlossener Anfrage wird als Sitzung aufgezeichnet“ in Handbuch Zeile 215/446 und README eng fassen oder die Ausnahme `COPY` ergänzen, mit Probe-Zeile) und V-311 (Gate-Lauf auf dem Stand nach dem letzten Commit in §7 nennen). V-312 und V-313 sind Hinweise.

**Offen für die Closure (Implementer):** Notiz mit Lerneintrag (Finding-Klasse: „aus einem Nachbarsatz übernommene Zusage ohne seine Voraussetzung/Ausnahme“, F-686 bis F-689 und V-310; die Durchsicht nahm bei der Nacharbeit „erste Fehler“, „Reihenfolge“, „Zwangsende“, aber nicht „jede Verbindung“), Register, Ausgang der beiden Risiken aus §6, Paarungen, dann `git mv` nach `done/`. Die Closure habe ich nicht ausgeführt.

**Finding-Klasse dieses Laufs:** „Allgemeinaussage ‚jede …‘ im Bestand, die der neue Abschnitt einschränkt, ohne dass der Bestandssatz die Einschränkung trägt“ (V-310, wie F-686 bis F-689).

**Gate:** `make gates` habe ich nicht gefahren (auf Anweisung des Aufrufers); `make docs-check` (670 Dateien, 0 Befunde) und `make abdeckung-check` (Exit 0) liefen am Stand `9b0a79a` mit sauberem Arbeitsbaum. Ein Gate-Lauf belegt nur den Stand, auf dem er lief; für `make gates` gilt die Aussage des Implementers und des Stop-Hooks (Exit 0 für `9b0a79a`), nicht mein Lauf.
