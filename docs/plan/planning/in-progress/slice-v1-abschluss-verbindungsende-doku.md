# Slice slice-v1-abschluss-verbindungsende-doku: Handbuch und README zu: Verbindungsende durch den Server

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-25`](../../../../spec/lastenheft.md#lh-fa-25--verbindungsende-durch-den-server), [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-FA-20`](../../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-FA-03`](../../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-02`](../../../../spec/lastenheft.md#lh-fa-02--record-modus), [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-verbindungsende-format` und `slice-v1-abschluss-verbindungsende-format-schreiber` bis `slice-v1-abschluss-verbindungsende-upstream`)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das die Slices `slice-v1-abschluss-verbindungsende-format`, `-format-schreiber`, `-sqlite`, `-sqlite-schreiber`, `-replay`, `-replay-extended`, `-replay-extended-flush`, `-einspielen`, `-einspielen-extended`, `-einspielen-grenzen`, `-vergleich`, `-hilfetext-binary`, `-record` und `-upstream` liefern, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice der sieben Code-Slices des Verbindungsendes, gebündelt hinter dem letzten, nach Entscheidung des Nutzers vom 2026-10-12 (Vorschlag des Planners; Begründung: Ein Anwender kann ein Ende erst mit `record` erzeugen; bis dahin kennt kein Weg des Binaries eine Aufzeichnung mit Ende außer handgeschriebenen Dateien, die das Handbuch nicht beschreibt. Der Doku-Teil ist ungeteilt, wie bei `slice-v1-abschluss-zeitangaben-doku` und `slice-v1-abschluss-sqlite-format-doku`.) Der Teil, den die Code-Slices in DoD und §3 trügen, liegt hier, weil die Dokumentation eine eigene Schicht ist (`AGENTS.md` §3.13). Zwischen den Code-Slices und diesem Slice hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen. Die Welle schließt erst mit diesem Slice (§3.11).

**Bindung aus `slice-v1-abschluss-verbindungsende-upstream-extended`** (Verschiebung hinter das Release, Entscheidung des Nutzers vom 2026-10-16; Sendung des Planners, `AGENTS.md` §3.13): Der Slice liegt hinter dem Release, nach diesem Slice, und ist keine Voraussetzung (Schritt 42 setzt 35a bis 41 voraus). Das Handbuch und das README sagen zum Ende einer Extended-Interaktion beim Aufzeichnen gegen PostgreSQL nichts und lassen es weg, auch wenn ein Test mit Fake-Server (`-upstream` DoD-Punkt 1) es zeigt; sie nennen das Ende einer Extended-Interaktion nur für das, was die Probe am Binary zeigt (Wiedergeben und Einspielen handgeschriebener Aufzeichnungen, Aufzeichnen gegen PostgreSQL nur für die einfache Anfrage). Die Probe sucht den Wortlaut im Handbuch und im README nach; ein Satz, der das Aufzeichnen einer Extended-Interaktion mit Ende gegen PostgreSQL zusagt, ist ein Befund.

**Handbuch-Lücke akzeptiert** (Entscheidung des Nutzers vom 2026-10-12, Vorschlag des Planners): Zwischen `slice-v1-abschluss-verbindungsende-vergleich` und diesem Slice stimmt der Satz des Handbuchs zu `FATAL`/`PANIC` (aus `slice-v1-abschluss-antwortvergleich-doku`) für eine handgeschriebene Aufzeichnung mit `connection_closed` nicht mehr; das Handbuch hinkt dem Binary bis zu diesem Slice hinterher, ein Release liegt nicht dazwischen. Die Lücke steht in `slice-v1-abschluss-verbindungsende-einspielen`, `-einspielen-extended` und `-vergleich` §1 als derselbe Vermerk.

**Übernommen aus `slice-v1-abschluss-antwortvergleich-doku`** (dort §1, Abgrenzung; Sendung des Planners, 2026-10-12, `AGENTS.md` §3.13): die Sätze zu `FATAL` und `PANIC` bei `--compare-responses` (dort „sind `PGR-E5004`“). Sie beschreiben bis `slice-v1-abschluss-verbindungsende-vergleich` den Ist-Zustand; ab dann gilt die Ausnahme eines aufgezeichneten Endes an derselben Stelle mit demselben SQLSTATE, und dieser Slice ersetzt den Satz.

**Übernommen aus `slice-v1-abschluss-sessions-doku`** (dort §1/§7, *Ungenannt geblieben*; Sendung des Planners bei dessen Closure, 2026-10-13, `AGENTS.md` §3.13): das aufgezeichnete Ende des Servers (`FATAL`, `PANIC`) im Zusammenhang mehrerer Verbindungen und Sitzungen. Das Handbuch nennt im Abschnitt *Mehrere Verbindungen aufzeichnen und wiedergeben* dazu keinen Satz (Lücke vom Nutzer akzeptiert, Entscheidung vom 2026-10-12, wie oben); dieser Slice prüft die Sätze jenes Abschnitts und der Sitzungs-Hinweise bei *Aufzeichnen* und *Wiedergeben* gegen das Verbindungsende (Sitzung mit Ende des Servers, Zuordnung der Sitzungen, Exit-Code, Verbindungsfehler) und sagt dazu nur, was die Probe am Binary zeigt. Sie gehört zu DoD-Punkt 1 (Handbuch) und 2 (README, soweit es Sitzungen nennt).

**Übernommen aus den Slices des Verbindungsendes** (Sendung des Planners, 2026-10-12): die Sätze des Handbuchs, die ein Ende des Servers berühren — zu `--allow-recorded-errors` (der Satz, ein Fehler mit dem Schweregrad `FATAL` breche das Einspielen auch dann ab), zur Version der Aufzeichnung in YAML und in SQLite (`version`, `meta`), zu `record` (Verbindung des Servers endet nach einer Fehlerantwort), zu `replay` und zu den Meldungscodes `PGR-E3002`, `PGR-E3003`, `PGR-E4003`, `PGR-E5004` und `PGR-E6001`. Welche Sätze ins Handbuch kommen, bestimmt die Probe am gebauten Binary; was sie nicht zeigt, bleibt draußen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur das Ende des Servers.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` in den Code-Slices.
- Der Satz „ältere Programmversionen lehnen Version 2 ab“ — nicht belegbar ohne ältere Programmversion (sie zu bauen, ohne das Image `pgwire-recorder:dev` zu überschreiben, wäre ein eigener Aufwand, `AGENTS.md` §3.11); die Pflicht aus dem Produktvertrag trägt `PGR-E3002` (jede Version außer 1 und 2), belegt. Gestrichen aus DoD-Punkt 1 (Review F-1, Entscheidung des Orchestrators 2026-10-17), kein Folge-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Das Benutzerhandbuch beschreibt, was die Slices des Verbindungsendes liefern, im Ist-Zustand des gebauten Binaries (Probe am Binary, gegen PostgreSQL und gegen handgeschriebene Aufzeichnungen): wann `record` ein Ende aufzeichnet und wann nicht (Fehlerantwort `FATAL`/`PANIC` in laufender Anfrage; Ende im Leerlauf nur mit nachfolgender Anfrage; Netzwerkabbruch, Ende im Aufbau und `ERROR` mit Abbruch nicht), was der Client sieht, das Format (`version` 1 und 2 in YAML und SQLite, ein Ende macht die Version 2; jede andere Version als 1 und 2 lehnt das Programm mit `PGR-E3002` ab), was `replay` wiedergibt, wie `play` ein aufgezeichnetes Ende behandelt (`--allow-recorded-errors` mit und ohne `--compare-responses`, `PGR-E4003`, `PGR-E5004`, Zähler) und die Zeilen der Meldungscodes, dazu das Ende des Servers im Zusammenhang mehrerer Verbindungen und Sitzungen (aus `slice-v1-abschluss-sessions-doku`, §1); die Sätze, die dieser Stand überholt, sind ersetzt. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11).
- [ ] `README.md` nennt, was die Slices des Verbindungsendes liefern, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Stand überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu Aufzeichnen, Wiedergeben, Einspielen (Hinweise zu `--allow-recorded-errors`), Format von YAML und SQLite und Meldungscodes |
| `README.md` | update | Funktionsumfang im Ist-Zustand |
| Beispiel-Dateien für die Probe (`test/integration/testdata` oder Probe-Verzeichnis) | entfällt | Die Probe nutzt die Aufzeichnungen, die `record` im Probe-Verzeichnis außerhalb des Repos erzeugt, und die vorhandenen Fixtures unter `test/integration/testdata` (`ende-*.yaml`, `play-ende-*.yaml`); jedes Beispiel des Handbuchs läuft als Datei gegen das Binary (§7). Unter `test/integration` entsteht nichts. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-verbindungsende-upstream` (Schritt 41) liegt in `done/` (Voraussetzung: das Binary zeigt das ganze Verhalten, das diese Welle liefert; das Handbuch nennt das Ende einer Extended-Interaktion gegen PostgreSQL erst nach dem Beleg in `slice-v1-abschluss-verbindungsende-upstream-extended`, der seit dem 2026-10-16 hinter dem Release liegt und keine Voraussetzung ist, Bindung aus dessen §1); `slice-v1-abschluss-antwortvergleich-doku` liegt davor. Dieser Slice ist Schritt 42 der Reihenfolge in §5 von welle-v1-abschluss, hinter Schritt 41. WIP-Limit 1.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: Die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- docs/user README.md test/integration` liegt über 600 (`AGENTS.md` §3.14; **Basis:** `e9f61167`, der Commit vor dem Anspruch; die Plan- und der Move-Commit des Anspruchs berühren keinen der Pfade, die Basis ändert sich danach nicht; die Beispiel-Dateien der Probe liegen unter `test/integration` als `testdata/*.yaml` und Tests, die sie gegen das gebaute Binary fahren, und zählen mit). **Meldeschwelle:** über 480 meldet der Implementer dem Architect. Schätzung des Planners in dieser Metrik: 150 bis 300 Zeilen Text, mit den Beispiel-Dateien der Probe 200 bis 450, Mitte 300 (50 % der Marke); Schnitt vor dem Anspruch: nein (die Obergrenze 450 liegt unter der Marke). Oder: Die Probe am Binary zeigt Verhalten, das die Spezifikation anders regelt, oder eines der Code-Slices liefert etwas, das §1 nicht trägt — Befund an den Planner, zurück zur Zerlegung.
- `in-progress` → `open`: Ein Code-Slice des Verbindungsendes wird zurückgeführt und liefert sein Verhalten nicht — Carveout.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen:** keine; der Slice liefert keinen neuen Vertrag (`AGENTS.md` §3.12), er beschreibt das Verhalten der Code-Slices.

**Risiken:**

- Das Handbuch sagt das Ende einer Extended-Interaktion beim Aufzeichnen gegen PostgreSQL zu, obwohl `slice-v1-abschluss-verbindungsende-upstream-extended` hinter dem Release liegt und es nicht belegt hat — **Ausgang:** offen bis Closure (erwartet: entfallen, die Probe sucht den Wortlaut nach und der Befund bleibt aus; Gegenmittel ist die Bindung in §1).
- Das Handbuch sagt zu, was nur gegen PostgreSQL 17 und handgeschriebene Aufzeichnungen geprüft ist (`BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft`, 1×); es nennt keine Version, die die Probe nicht gefahren hat — **Ausgang:** offen bis Closure.
- Bis `slice-v1-abschluss-verbindungsende-vergleich` (Schritt 39a) beschreiben die Sätze von `slice-v1-abschluss-antwortvergleich-doku` zu `FATAL`/`PANIC` den Ist-Zustand; zwischen den Code-Slices und diesem Slice ist der Satz für eine handgeschriebene Aufzeichnung mit Ende nicht mehr wahr; die Lücke ist akzeptiert (Entscheidung des Nutzers vom 2026-10-12, §1; `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** offen bis Closure (erwartet: eingetreten und mit diesem Slice ersetzt).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

### Belege des Implementers

**Probe-Aufbau.** `make build` am Stand `d2b3059c` (Image `pgwire-recorder:dev`, `sha256:facf1c729079…`, Produktcode unverändert; geändert sind nur `docs/user/benutzerhandbuch.md`, `README.md` und dieser Plan). Proben in einem Scratch-Verzeichnis außerhalb des Repos; ein Docker-Netz `pgr-doku`, PostgreSQL 17.11 aus dem gepinnten Image von `harness/mk/integration.mk` (`postgres:17-alpine@sha256:b0f9560a…`, Servermodus `trust`, Rolle und Datenbank `app`), `record`, `replay` und `play` als `docker run … pgwire-recorder:dev` im selben Netz, Client `psql` 17.11 aus demselben PostgreSQL-Image, `sqlite3` 3.45.1 auf dem Host. Ein Lauf zur Zeit, kein `pkill`. Die Beispiele des neuen Abschnitts liefen **wörtlich**: Der Wortlaut jeder Kommandozeile des Abschnitts wurde mit `grep -F` im Handbuch bestätigt (alle `ok`); `psql` lief im Netzraum des Recorders (`--network container:rec`), damit `host=localhost port=15432` unverändert gilt, der Upstream hieß `postgres`, die Dateien lagen unter `./recordings/` im Arbeitsverzeichnis des Containers. Fixtures: die vorhandenen unter `test/integration/testdata` (`ende-panic.yaml`, `play-ende-vergleich-gleich.yaml`, `-vergleich-anders.yaml`, `-grenze-ohne-ende.yaml`, `-grenze-ohne-option.yaml`), schreibgeschützt eingehängt. Zusätzlich im Scratch (nicht im Repo): `error-statt-fatal.yaml` (Sitzung 1 von `ende.yaml` mit Schweregrad `ERROR` und `ready_for_query` statt `connection_closed`, `version: 1`), Kopien mit geänderter Version (`sed` bzw. `update meta`), ein Ersatz-Server (`fake.py`, Python 3.12 im Image `python:3.12-slim`: Startup, `AuthenticationOk`, `ReadyForQuery`, auf die erste Anfrage `PANIC` bzw. `ERROR`, Verbindung schließen) und ein Client (`raw.py`), der die Nachrichtentypen bis zum Ende der Verbindung zählt. Container und Netz sind danach entfernt (`docker ps -a`, `docker network ls` ohne Rest). Nur Linux im Container (`BEO-REPO/verhalten-nur-unter-linux-geprueft`); PostgreSQL nur 17.11, `PANIC` und `ERROR` mit Abbruch nur mit dem Ersatz-Server (das Handbuch sagt es an den Stellen).

**Beispiele des Handbuchs: Datei · Aufruf · beobachtete Ausgabe.**

| Beispiel | Datei · Aufruf | Beobachtet |
|---|---|---|
| Aufzeichnen mit Ende | `./recordings/ende.yaml`; `record --listen 0.0.0.0:15432 --upstream postgres:5432 --output ./recordings/ende.yaml`, `psql … -c "select pg_terminate_backend(pg_backend_pid())"`, `psql … -c "select 2"`, `docker stop` | erster `psql`: die fünf Zeilen des Handbuchs, Exit 2 (mit Tabulator; `cat -A`); zweiter: Zeile `2`, Exit 0; `record` Exit 0, Log ohne `PGR-`; Datei `version: 2`, zwei Sitzungen, ein `connection_closed`; Auszug Zeilen 41 bis 51 gleich dem Block im Handbuch |
| Wiedergeben | `replay --listen 0.0.0.0:15432 --input ./recordings/ende.yaml`, dieselben zwei `psql`-Aufrufe | erster: gleiche Meldung, Exit 2; zweiter: Zeile `2`; `replay` Exit 0, Log: gestartet, Herunterfahren, beendet |
| Einspielen mit Option | `play --upstream postgres:5432 --input ./recordings/ende.yaml --allow-recorded-errors` | Exit 0, `eingespielt=2 verglichen=0 abweichend=0` (Zeile im Handbuch wörtlich, Zeitstempel des Laufs) |
| Einspielen ohne Option (Fehlerbild `PGR-E4003`) | wie oben ohne `--allow-recorded-errors` | Exit 4, beide Zeilen wörtlich im Handbuch (Zeitstempel des Laufs), `eingespielt=1` |
| SQLite-Abfragen | `./recordings/ende.sqlite` (`record … --format sqlite`, dieselben zwei `psql`-Aufrufe); `sqlite3 ./recordings/ende.sqlite "select key, value from meta"` und `"select session_id, sequence, kind from message where kind = 'connection_closed'"` | `format|pgwire-recorder`, `version|2`; `1|1|connection_closed` |
| Fehlerbild `PGR-E5004` bei Ende | `./recordings/leer.yaml` (Aufnahme von `set idle_session_timeout = '1s'`, `\! sleep 3`, `select 1;` mit `psql -f`; `play --upstream postgres:5432 --input ./recordings/leer.yaml --compare-responses` | Exit 5, beide Zeilen wörtlich im Handbuch (Zeitstempel des Laufs): `Position 1, Art Nachrichtenart`, `eingespielt=2 verglichen=2 abweichend=1` |

**Zusagen: Satz · Probe · beobachtet · Gegenprobe (welcher abweichende Wortlaut die Probe widerlegte).** `V` ist der Lauf mit den Dateien aus dem Scratch (`hb1.yaml` ist `ende.yaml` des Handbuchs), `F` die Fixtures.

| # | Satz im Handbuch/README | Probe | Beobachtet | Gegenprobe |
|---|---|---|---|---|
| 1 | Administrator-Ende (`57P01`) und Leerlauf (`57P05`) senden `FATAL` und schließen die Verbindung (PostgreSQL 17.11) | `pg_terminate_backend`; `idle_session_timeout`; `raw.py` direkt am Server | `C: 57P01`/`57P05`, `S/V: FATAL`; Typen `TDE`, dann Ende der Verbindung | wäre „`ERROR`“ oder „`ReadyForQuery` folgt“ gesagt, widerlegte es `TDE` ohne `Z` |
| 2 | Das Werkzeug behandelt `PANIC` ebenso | `fake.py panic`: `record`; `ende-panic.yaml`: `replay`; `fake.py panic`: `play` ohne Option, mit `--allow-recorded-errors`, mit `--compare-responses` | `record`: `version: 2`, `error_response`, `connection_closed`; `replay`: `PANIC:  server panicked`; `play`: `PGR-E4003` Exit 4 / Exit 0 / Exit 0 `verglichen=1` | „nur `FATAL`“ widerlegt der aufgezeichnete `PANIC` |
| 3 | Geprüft mit PostgreSQL 17.11 und `psql` 17.11; `PANIC` nur mit Ersatz-Server | `select version()`, `psql --version` | `17.11` beide | Satz „alle Versionen“ hätte keinen Lauf |
| 4 | Aufzeichnen: Antworten der Anfrage bis einschließlich Fehlerantwort, dann `connection_closed` statt `ready_for_query`, letzte Anfrage der Sitzung; die Verbindung, deren erste Anfrage so endet, ist eine Sitzung | `ende.yaml`: `grep` auf `sequence`, `connection_closed`; Sitzung 1 hat `sequence: 1` mit Ende | Sitzung 1 mit genau dieser Anfrage; `connection_closed` einmal | „zeichnet die Verbindung nicht auf“ widerlegt Sitzung 1 |
| 5 | Anwendung erhält Antworten bis Fehlerantwort, kein `ReadyForQuery`, danach endet die Verbindung | `raw.py` an `record`, an `replay`, direkt an den Server; normale Anfrage zum Vergleich | `TDE` + `EOFError` jeweils; normale Anfrage `TDCZ` (Verbindung bleibt offen) | „kein `ReadyForQuery`“ wäre falsch, wenn `Z` käme: `TDCZ` zeigt, dass `raw.py` es sähe |
| 6 | Keine Meldung, Exit-Code unverändert (0), andere Verbindungen laufen weiter und werden aufgezeichnet, auch eine offene | `record`-Log; `record` mit `psql -f` (`pg_sleep(3)`; `select 5`) in A und Ende in B | Log ohne `PGR-`, Exit 0; Sitzung 2 (A) mit `select 5` aufgezeichnet | „stoppt das Werkzeug“ widerlegt A |
| 7 | Die Aufzeichnung kann mehrere Sitzungen mit Ende enthalten | drei Verbindungen (Ende, `select 2`, Ende), `--format sqlite` | `connection_closed` in Sitzung 1 und 3 | „höchstens eine“ widerlegt es |
| 8 | Leerlauf: Ende nur mit nachfolgender Anfrage; sonst kein Ende, `version: 1` | `psql -f` mit `set idle_session_timeout`, `\! sleep 3`, mit und ohne `select 1;` | mit: `select 1;` trägt `C: 57P05` und `connection_closed`, `version: 2`; ohne: kein `error_response`, kein `connection_closed`, `version: 1` | „immer aufgezeichnet“ widerlegt der Fall ohne Anfrage |
| 9 | Abbruch ohne Fehlerantwort: `PGR-E4003`, Exit 4 | `docker kill` des PostgreSQL-Containers während `pg_sleep(30)` | Log `PGR-E4003`, `record` Exit 4, `psql`: `FATAL:  Netzwerk [PGR-E4003]` | „Ende“ widerlegt der Exit-Code |
| 10 | `ERROR` mit Abbruch: `PGR-E6001`, Exit 6, keine Sitzung (nur Ersatz-Server) | `fake.py error` | Log `PGR-E6001`, Exit 6, Datei `version: 1` ohne Sitzung | „als Ende aufgezeichnet“ widerlegt die Datei |
| 11 | `FATAL` in der Anmeldung ist kein Ende: Anwendung erhält sie, keine Meldung, keine Sitzung | `psql … dbname=nirgends` über `record` | `FATAL:  database "nirgends" does not exist`, Exit 0, Log ohne `PGR-`, `sessions: []`, `version: 1` | |
| 12 | `version: 2` mit Ende, `version: 1` ohne; SQLite `meta` 2, sobald eine Sitzung mit Ende in der Datei steht | `ende.yaml`, `leer2.yaml`, `ohneende.sqlite`, `spaet.sqlite` (erst `select 2`, dann Ende, `meta` nach jedem Schritt) | `version: 2`/`1`; `spaet.sqlite`: `version|1`, nach der zweiten Verbindung `version|2` | „2 von Anfang an“ widerlegt `spaet.sqlite` |
| 13 | SQLite: letzte Nachricht `kind = 'connection_closed'`, `direction = 'server'`, `fields = '{}'` | `select session_id, sequence, position, direction, kind, fields from message where kind='connection_closed'` | `1|1|4|server|connection_closed|{}` | |
| 14 | `replay`: Antworten bis Fehlerantwort, dann Verbindung zu, ohne `ReadyForQuery`, Meldung, Fehler (Exit 0); `FATAL` wie `PANIC`; YAML wie SQLite | `replay` mit `ende.yaml`, `hb2.sqlite`, `ende-panic.yaml` (`psql -c "SELECT 1;"`), `raw.py` | siehe 5; `PANIC:  server panicked`; SQLite: `FATAL:  terminating connection due to administrator command`, `2`, `FATAL…`; Exit 0, Log ohne `PGR-` | |
| 15 | Sitzung gilt als abgeschlossen (kein `PGR-W2001`, auch nicht mit `--fail-on-unconsumed`); nächste Verbindung erhält nächste Sitzung | `replay --fail-on-unconsumed` mit `hb1.yaml`, zwei `psql`-Aufrufe | Exit 0, Log ohne Warnung; zweiter Aufruf liefert `2` | „bricht ab“ widerlegt Exit 0 |
| 16 | `play` ohne Option, ohne Vergleich: `PGR-E4003`, Exit 4, bricht ab; auch mit `--continue-on-error` | `ende.yaml`, `F` `play-ende-grenze-ohne-option.yaml` mit `--continue-on-error` | Exit 4, `eingespielt=1` (Sitzung 2 nicht gespielt); mit `--continue-on-error` Exit 4 | „läuft mit `--continue-on-error` weiter“ widerlegt es |
| 17 | `--allow-recorded-errors`: erwartet, wenn die aufgezeichnete Anfrage ein Ende trägt, gleich welcher SQLSTATE; nächste Sitzung läuft ohne `--continue-on-error`; ohne Ende in der Aufzeichnung `PGR-E4003` | `F` `play-ende-vergleich-anders.yaml` (aufgezeichnet `57P02`, Server `57P01`), `play-ende-grenze-ohne-ende.yaml` (aufgezeichnet `ERROR` mit `ready_for_query`), `error-statt-fatal.yaml`; Tabellen vorher gelöscht | `anders`: Exit 0, `eingespielt=3`, beide Tabellen da; `ohne-ende` und `error-statt-fatal`: `PGR-E4003`, Exit 4 | „nur bei gleichem SQLSTATE“ widerlegt `anders`; „bei jedem Fehler in der Aufzeichnung“ widerlegt `error-statt-fatal` |
| 18 | `--compare-responses`: Ende erwartet bei gleichem SQLSTATE (verglichen gezählt); anderer SQLSTATE `PGR-E5004` Art `Fehler` mit beiden SQLSTATE; Vorbild ein Fehler ohne Ende: `PGR-E5004`; mit `--continue-on-error` läuft die nächste Sitzung; `--allow-recorded-errors` ohne Wirkung | `ende.yaml`, `play-ende-vergleich-gleich.yaml`, `-anders.yaml` (auch mit `--continue-on-error`), `error-statt-fatal.yaml`, `play-ende-grenze-ohne-ende.yaml`, je mit und ohne `--allow-recorded-errors` | `gleich`: Exit 0, `verglichen=3`; `anders`: `Position 1, Art Fehler (SQLSTATE aufgezeichnet „57P02“, Server „57P01“)`, Exit 5; mit `--continue-on-error` `eingespielt=3`, beide Tabellen; `error-statt-fatal`: `Position 2, Art Fehler (… „57P01“, … „57P01“)`; `ohne-ende`: `Position 1, Art Nachrichtenart` mit und ohne die Option gleich | „gleicher SQLSTATE genügt“ widerlegt `error-statt-fatal` |
| 19 | Datenbank ohne Ende, Aufzeichnung mit Ende: ohne Vergleich nichts, mit Vergleich `PGR-E5004` | `leer.yaml` ohne und mit `--compare-responses` | Exit 0 / Exit 5 | |
| 20 | `PGR-E3003` für Version 1 mit `connection_closed`, `PGR-E3002` für andere Version als 1 und 2 (beide Formate) | Kopien von `ende.yaml`/`ende.sqlite` mit Version 1 bzw. 3 | `Recording [PGR-E3003]: … connection_closed in Version 1` (YAML), `PGR-E3003` (SQLite); `PGR-E3002` ×2; Exit 3 | |
| 21 | Tabellenzeilen `PGR-E3002`/`E3003`/`E4003`/`E5004`/`E6001`, Glossar, Zählerzeile | wie 9, 10, 16 bis 20 | siehe dort | |
| 22 | README: ein Absatz zum Ende; „Fehlerantwort der Datenbank, die die Verbindung nicht beendet“ | wie 16 bis 18 | siehe dort | |

**Durchsicht des Bestands (Quantoren über den Gegenstand, je Treffer die Entscheidung).** Gesucht je Absatz mit „jede/alle/immer/stets/nie/nur/kein/jeweils“ zusammen mit Fehlerantwort, beenden, Ende, FATAL, PANIC, Version (Handbuch und README, ohne den neuen Abschnitt). Treffer mit Bezug zum Ende:

| Stelle | Entscheidung |
|---|---|
| Record-Hinweis „Eine Verbindung … mindestens eine Anfrage abgeschlossen hat“ (Handbuch) und derselbe Satz im Abschnitt zu mehreren Verbindungen und im README | verbunden: „(auch eine, die mit dem Ende der Verbindung durch die Datenbank endet)“ an allen drei Stellen gleich; Zeile 4 |
| „Ohne `--record-empty-sessions` ist eine Verbindung ohne abgeschlossene Anfrage keine Sitzung“ und die Liste der Verbindungen, die auch mit der Option keine Sitzung sind | unverändert: Die Anfrage mit Ende zählt als abgeschlossen (Zeile 4); ein `FATAL` in der Anmeldung steht dort schon (Zeile 11) |
| „wartet, bis jede Verbindung ihre laufende Anfrage … abgeschlossen hat“ (Herunterfahren, `record`) und der Abschnitt *Herunterfahren mit Frist* | unverändert: Sie sagen nichts zum Ende; das Handbuch sagt zum Herunterfahren mit Ende nichts zu, keine Probe |
| „Endet eine Verbindung, bevor alle Anfragen ihrer Sitzung … beantwortet sind, warnt das Werkzeug (`PGR-W2001`)“ | unverändert: Ein Ende ist die letzte Interaktion, `replay` warnt nicht (Zeile 15) |
| „Ihre Anwendung erhält die aufgezeichneten Antworten … auch Fehlerantworten“ (*Wiedergeben*, Ergebnis) | verbunden: Satz zum Ende der Verbindung mit Verweis (Zeile 14) |
| Einspielen: „Antwortet die Datenbank auf eine Anfrage mit einem Fehler, bricht das Einspielen ab (`PGR-E4004`)“, derselbe Satz beim erweiterten Protokoll, „Ein Fehler mit dem Schweregrad `FATAL` … bricht … ab (`PGR-E4003`)“ | eng gefasst („der die Verbindung nicht beendet“) bzw. ersetzt (`FATAL` oder `PANIC`, erwartet nur mit Ende in der Aufzeichnung); Zeilen 16, 17 |
| Vergleichen: „Trägt die Aufzeichnung … eine Fehlerantwort mit demselben SQLSTATE, ist das keine Abweichung“, „Mit Vergleich hat `--allow-recorded-errors` keine Wirkung … einer mit demselben SQLSTATE ist auch ohne die Option erwartet“, der `FATAL`-Satz | eng gefasst („die kein Ende ist“, „Schweregrad `ERROR`“) bzw. ersetzt; Zeilen 17, 18; der `FATAL`-Satz nennt jetzt erwartet/`PGR-E5004` wie die Zeilen belegen |
| Vergleichen: „Ohne `--continue-on-error` endet `play` nach der Interaktion mit der Abweichung … Mit … macht es mit der nächsten Interaktion weiter“ | unverändert: Gilt auch für ein abweichendes Ende (Zeile 18, `eingespielt=3`) |
| „Bricht die Verbindung ohne Fehlerantwort der Datenbank ab, ist das `PGR-E4003` … auch nach einer Abweichung“ | unverändert: früher belegt (`slice-v1-abschluss-antwortvergleich-doku`); der neue Abschnitt nennt den Fall für `record` (Zeile 9) |
| Log-Tabelle „`error`: Fehler, die eine Verbindung beenden“, Exit-Code 0 „wenn zuvor kein Fehler auftrat“ | unverändert: Ein Ende ist kein Fehler, keine Log-Zeile, Exit 0 (Zeile 6) |
| SQLite: „Jede weitere beendete Sitzung ergänzt das Werkzeug danach in der Datei“ | unverändert: gilt auch für die Sitzung mit Ende (`spaet.sqlite`, Zeile 12) |
| `PGR-E3002`/`E3003`/`E4003`/`E5004`/`E6001` in der Fehlercode-Tabelle | ergänzt, Zeilen 9, 10, 16, 18, 20 |

**Nicht aufgenommen, weil nicht belegbar (Abweichung von DoD-Punkt 1, an Review und Planner).** (a) „Ältere Programmversionen lehnen Version 2 ab“: Es gibt keine ältere Programmversion zum Lauf, ohne das Image `pgwire-recorder:dev` zu überschreiben; das Handbuch nennt nur `PGR-E3002` für jede andere Version als 1 und 2 (Zeile 20). (b) „Ende im Aufbau“ als eigener Satz: Der Bestand zur Anmeldung trägt ihn (Zeile 11), der neue Abschnitt verweist darauf. (c) Das Ende einer Extended-Interaktion beim Aufzeichnen gegen PostgreSQL und „nur noch Schließen“ nach dem Ende: nicht genannt (Bindung aus `-upstream-extended`); die Probe suchte den Wortlaut im Diff: `grep -i 'erweitert|Extended|Sync|Flush'` im neuen Abschnitt ohne Treffer, im Diff keine Treffer für Chronik, Spezifikation, ADR, Slice, Welle, Kennungen. (d) Herunterfahren mit Ende (Spezifikation: „Herunterfahren beim Aufzeichnen“): kein Satz, keine Probe.

**Läufe.** `make docs-check`: 932 Dateien, 0 Befunde (nach den Änderungen). `make kopf-check`: Exit 0. `make abdeckung-check`: Exit 0 (keine Tests berührt). `make gates-schnell`: Exit 0 am Arbeitsbaum mit Commit `1fbd217b` (geändert sind nur `docs/user/`, `README.md` und der Plan, alle auf der Positivliste des Stop-Hooks; ein voller `make gates` ist dafür nicht nötig und lief nicht, `test/integration` ist unberührt). Messung der Marke: `git diff --numstat e9f61167 HEAD -- docs/user README.md test/integration` ergibt 287 (9+4 und 244+30; Einfügungen plus Löschungen) bei Marke 600, Meldeschwelle 480 nicht erreicht.

**Restrisiken.** Die Ersatz-Server-Proben (`PANIC`, `ERROR` mit Abbruch) belegen das Verhalten des Binaries, nicht das von PostgreSQL; das Handbuch sagt es. `psql` ist die einzige geprüfte Anwendung (die Ausgabe der Meldungen ist die von libpq 17.11). Die Zeitstempel in den Log-Blöcken sind die der Läufe.

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen. Treffer in der Sub-Area `REPO`: `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×; geliefert bei den bisherigen Doku-Folge-Slices, ohne Auftreten) und `BEO-REPO/serververhalten-nur-gegen-eine-version-geprueft` (1×; §6) — die Probe prüft die ersetzten Sätze mit der Suche nach den alten Wortlauten. Keine weiteren Treffer.

Nachgezählt beim Anlegen (Planner, 2026-10-12, `AGENTS.md` §3.13): Liefer-Punkte: 2 (Handbuch, README). Schichten: eine (Nutzer- und Wartungs-Doku). Schätzung des Planners: 150 bis 300 Zeilen Text.

Nachgezählt beim Eintragen der Übernahme aus `slice-v1-abschluss-sessions-doku` (Planner, 2026-10-13, `AGENTS.md` §3.13): Sie gehört zu DoD-Punkt 1 (Handbuch) und 2 (README), kein neuer Liefer-Punkt: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku).

Nachgezählt beim Eintragen der Kennung `slice-v1-abschluss-verbindungsende-format-schreiber` (Ziel und Berührte Spec-Stellen; Planner, 2026-10-15, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-format`): Der Teil beschreibt, was das Binary zeigt, und gehört zu DoD-Punkt 1 (Handbuch) und 2 (README); kein neuer Liefer-Punkt: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Spezifikation und [ADR-0039](../../adr/0039-aufzeichnung-kennt-das-verbindungsende.md) gehen dem Code voraus).

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-sqlite-schreiber` (Ziel; Planner, 2026-10-15, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-sqlite`): Der Teil beschreibt, was das Binary zeigt, und gehört zu DoD-Punkt 1 (Handbuch) und 2 (README); kein neuer Liefer-Punkt: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku).

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-replay-extended` (Ziel; Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-replay`): Der Teil beschreibt, was das Binary zeigt (Ende einer Extended-Interaktion im Replay), und gehört zu DoD-Punkt 1 (Handbuch) und 2 (README); kein neuer Liefer-Punkt: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Schritte 35a bis 41 der Voraussetzung (Schritt 42) schließen 37b ein.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-upstream-extended` (Ziel, Start und Bindung; Planner, 2026-10-16, `AGENTS.md` §3.13, Schritt 41b hinter `-upstream`): Der Teil beschreibt, was das Binary zeigt (Ende einer Extended-Interaktion gegen PostgreSQL 17, Probe am Binary), gehört zum selben Satz des Handbuchs wie 37b und 37c und zu DoD-Punkt 1 und 2; kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Voraussetzung (Schritt 42) schloss 41b ein.

Nachgezählt beim Herausnehmen der Kennung `slice-v1-abschluss-verbindungsende-upstream-extended` (Ziel, Start und Bindung; Planner, 2026-10-16, `AGENTS.md` §3.13, Verschiebung hinter das Release, Entscheidung des Nutzers): Der Teil entfällt aus Ziel, Start und Voraussetzung (Schritt 42 setzt 35a bis 41 voraus); die Bindung (Handbuch und README sagen zum Aufzeichnen einer Extended-Interaktion mit Ende gegen PostgreSQL nichts) schränkt den Text ein und liefert nichts. Liefer-Punkte weiter 2 (Handbuch, README), Schichten weiter eine (Nutzer- und Wartungs-Doku); die Schätzung 150 bis 300 Zeilen Text sinkt nicht und steigt nicht.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-replay-extended-flush` (Ziel; Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-replay-extended`): Der Teil beschreibt, was das Binary zeigt (Ende einer Extended-Interaktion mit `Flush` im Replay), gehört zum selben Satz des Handbuchs wie 37b und zu DoD-Punkt 1 und 2; kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Voraussetzung (Schritt 42) schließt 37c ein.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-einspielen-extended` (Ziel; Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt von `slice-v1-abschluss-verbindungsende-einspielen`): Der Teil beschreibt, was das Binary zeigt (Ende einer Extended-Interaktion beim Einspielen ohne Vergleich), gehört zum selben Satz des Handbuchs wie 38a und zu DoD-Punkt 1 und 2; kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Voraussetzung (Schritt 42) schließt 38b ein.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-einspielen-grenzen` (Ziel; Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt beim Anspruch aus `-einspielen-extended`): Der Teil beschreibt, was das Binary zeigt (ohne `--allow-recorded-errors` und bei aufgezeichnet ohne Ende an der Stelle ist die `FATAL` `PGR-E4003`), gehört zum selben Satz des Handbuchs wie 38a und zu DoD-Punkt 1 und 2; kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Voraussetzung (Schritt 42) schließt 38c ein.

Nachgezählt beim Nachziehen der Kennung `slice-v1-abschluss-verbindungsende-hilfetext-binary` (Ziel; Planner, 2026-10-16, `AGENTS.md` §3.13, Schnitt beim Anspruch aus `-vergleich`, Schritt 39b): Der Teil beschreibt, was das Binary zeigt (Ende im Vergleich am Binary, Hilfetext von `--allow-recorded-errors`), gehört zum selben Satz des Handbuchs wie 39a (`FATAL`/`PANIC` bei `--compare-responses`) und zu DoD-Punkt 1 und 2; kein neuer Liefer-Punkt, keine neue Schicht: Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku). Die Voraussetzung (Schritt 42) schließt 39b ein. Die Handbuch-Lücke (§1, §6) beginnt mit 39a, dem Schritt, der das Verhalten ändert; 39b ändert am Verhalten nichts.

Nachgezählt beim Anspruch (Planner, 2026-10-16, `AGENTS.md` §3.13, §3.14): Liefer-Punkte 2 (Handbuch, README), Schichten eine (Nutzer- und Wartungs-Doku; `test/integration` trägt nur die Beispiel-Dateien der Probe, keine Logik des Produkts). §1 (Ausschlüsse, Bindung aus `-upstream-extended`) und §6 tragen unverändert; die Metrik und die Marke stehen in §4.

Nachgezählt beim Streichen des Satzes „ältere Programmversionen lehnen Version 2 ab“ (Planner, 2026-10-17, `AGENTS.md` §3.13, Review F-1): DoD-Punkt 1 verlangt ihn nicht mehr, §1 nennt den Ausschluss mit Grund, §7 trägt die Abweichung (a) schon. Liefer-Punkte weiter 2, Schichten weiter eine (Nutzer- und Wartungs-Doku); der Text wird enger, die Schätzung bleibt.
