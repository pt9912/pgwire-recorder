# Slice slice-v1-abschluss-sqlite-format-doku: Handbuch und README zu: SQLite als Aufzeichnungsformat

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-abschluss.

**Bezug:** [`LH-FA-22`](../../../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [`LH-QA-06`](../../../../spec/lastenheft.md#lh-qa-06--wartbarkeit-des-recording-formats), [`LH-FA-07`](../../../../spec/lastenheft.md#lh-fa-07--persistente-recordings)

**Berührte Spec-Stellen:** — (der Slice ändert keine Spec-Stelle; er beschreibt das gelieferte Verhalten von `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen`)

**Verantwortlich:** pt9912

**Autor:** pt9912. **Datum:** 2026-10-10.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/user/benutzerhandbuch.md` und `README.md` beschreiben das Verhalten, das `slice-v1-abschluss-sqlite-format` (Schreiben) und `slice-v1-abschluss-sqlite-format-lesen` (Lesen) liefern, im Ist-Zustand des gebauten Binaries nach beiden Slices: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews (`AGENTS.md` §3.11).

**Herkunft:** Doku-Folge-Slice von `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` (jener Slice wurde am 2026-10-11 nach Entscheidung des Nutzers in *Schreiben* und *Lesen* geschnitten; der Doku-Teil blieb ungeteilt hier, weil das Handbuch das Format als Ganzes beschreibt), abgeschnitten nach Entscheidung des Nutzers vom 2026-10-10 („Dafür haben wir die Welle – man kann dafür einen weiteren Slice anlegen“) bei der Closure von `slice-doku-ist-stand` (Verifikation V-151): Mit dem Handbuch- und README-Teil als Dokumentation, einer eigenen Schicht (`AGENTS.md` §3.13), läge `slice-v1-abschluss-sqlite-format` (damals noch ungeteilt) über zwei Schichten. Der Teil, den jener Slice in DoD und §3 trug, liegt hier; er setzt beide Slices voraus (§4) und steht in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md) direkt hinter `slice-v1-abschluss-sqlite-format-lesen`. Zwischen den Code-Slices und diesem hinkt das Handbuch dem Binary hinterher; ein Release liegt nicht dazwischen, und die Welle schließt erst mit allen dreien. Das Handbuch beschreibt das Format nach `slice-v1-abschluss-sqlite-format-lesen`: Es sagt nichts über ein Format, das sich schreiben, aber nicht wiedergeben lässt (Entscheidung des Nutzers vom 2026-10-11, Zwischenstand der Code-Slices).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produkt-Code, Hilfetexte im Binary, Tests des Produkts und Gates — Schicht-Abgrenzung: Der Slice ändert nur `docs/user/benutzerhandbuch.md` und `README.md` und diesen Plan. Zeigt das Binary ein Verhalten, das die Spezifikation anders regelt, beschreibt das Handbuch das Binary (`slice-doku-ist-stand` §6, *Verhalten neben der Spezifikation*); der Fund geht als Befund an den Planner.
- Die Abschnitte des Handbuchs, die andere Slices liefern — jeder Slice führt seinen Teil in seinem Plan oder in seinem Doku-Folge-Slice; dieser Slice beschreibt nur, was `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` liefern.
- `docs/user/benutzerhandbuch-standard.md` und die Abdeckungstabellen `docs/user/abdeckung-*.md` — Bestand bleibt bewusst stehen: Der Standard ist die Schreibanleitung, kein Ist-Zustand des Binaries; die Tabellen schreibt `make abdeckung` im Code-Slice.
- Ein Sensor, der Handbuch oder README gegen das Binary prüft — ein anderer Vorgang: Er wäre ein neues Gate mit eigenem Vertrag in der Spezifikation; dieser Slice prüft von Hand gegen das gebaute Binary.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.
- [x] Das Benutzerhandbuch beschreibt, was `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` liefern, im Ist-Zustand des gebauten Binaries: Die Abschnitte, die Kommandos, Optionen, Meldungscodes, Beispiele und Grenzen dieses Slice betreffen, stehen im Handbuch, und die Grenz-Sätze, die er überholt, sind ersetzt. Jede Option steht gegen `--help` des Binaries, jeder Meldungscode im Katalog (`internal/hexagon/model/fehler.go`) und ausgelöst durch einen Test oder eine Probe, jedes Beispiel läuft als Datei gegen das Binary (`AGENTS.md` §3.11); das Handbuch enthält keinen Verweis auf Spezifikation, Lastenheft, ADRs, Slices, Wellen oder Reviews und kein „noch nicht“, „kommt“, „geplant“ über das Produkt (`grep`, Beleg in §7).
- [x] `README.md` nennt, was `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` liefern, im Ist-Zustand des gebauten Binaries, ohne Chronik und ohne Zielstand; die Sätze, die dieser Slice überholt, sind ersetzt, Verweise auf `spec/` und `docs/plan/` bleiben ohne Aussage über einen Stand (`slice-doku-ist-stand` §6, *Verweise im README*).
- [x] `make gates` grün — Beleg in §7, *Closure*, Stand der Gates.
- [x] Review durchgeführt, Report `docs/reviews/2026-10-11-review-slice-v1-abschluss-sqlite-format-doku.md` und Verifikation `docs/reviews/2026-10-11-verifikation-slice-v1-abschluss-sqlite-format-doku.md` liegen vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag — §7, *Closure*.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | Abschnitte zu dem, was `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` liefern: Optionen gegen `--help`, Beispiele als Datei, Codes, Grenz-Sätze ersetzt. Betroffen sind §1 (Zweck), §3 (Hinweis auf das Format), §4 (Aufzeichnen, Wiedergeben, Einspielen und der neue Abschnitt *Das Format der Aufzeichnung festlegen* samt *Eine SQLite-Aufzeichnung abfragen*), die Tabelle in §5, die Zeilen `PGR-E2001`, `PGR-E2004`, `PGR-E3001` bis `PGR-E3003` in §7, die FAQ und das Glossar |
| `README.md` | update | Der Satz in *Was kann ich heute tun?*, der `record` und `replay` beschreibt: Format `sqlite` und Erkennung durch `replay` und `play` |

**Ansatz:**

- Binary mit `make build` bauen, Proben netzlos in einem Scratch-Verzeichnis außerhalb des Repos (Beispieldateien, Aufzeichnungen), für `record` und `play` gegen das gepinnte PostgreSQL-Image aus `harness/mk/integration.mk` in einem eigenen Docker-Netz; Container, Netz und Dateien danach entfernt. Die Proben stehen als Liste in §7 (Aufruf, Ergebnis), nicht im Repo.
- Die Befehlsfolge aus `slice-doku-ist-stand` §7 (Optionen, Variablen und Codes des Handbuchs gegen `--help` und Katalog) läuft vor der Übergabe mit **einer benannten Ausnahme**: `PGWIRE_RECORDER_PASSWORD` ist keine Option und steht nicht in der Tabelle von §5; die Variablen-Zeile lautet `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`, und mit ihr ist die Ausgabe leer (Beleg und Mutanten: `slice-v1-abschluss-einspielen-anmeldung-doku` §7).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` liegen in `done/` (Voraussetzung; Schritt 18 der Reihenfolge in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md), direkt hinter dem Lesen; WIP-Limit 1). Vor dem ersten Commit am Handbuch entscheidet der Architect die Randformen aus §6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar; Schnitt dann: Handbuch und README getrennt.
- `in-progress` → `open` (blockiert — Carveout?): Eine Aussage lässt sich nur mit einer Änderung am Binary wahr machen, und der Nutzer will sie nicht streichen; dann zuerst die Entscheidung des Nutzers, welcher Slice das Verhalten liefert.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` grün, die Proben und `grep`-Läufe in §7, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12) — entschieden für diesen Typ Slice in `slice-doku-ist-stand` §6 (Ort ist dieser Plan, nicht die Spezifikation: Der Slice legt keinen Vertrag an): *Software-Version* und *Installation* (unberührt), *Teilweise geliefert* (eine Grenze steht im Präsens, wenn das Binary an ihr beobachtbar reagiert, kein Satz sagt, dass etwas kommt), *Verhalten neben der Spezifikation* (eine Zeile in §7 an den Planner), *Verweise im README*, *Wie geprüft wird*. Der Architect prüft vor dem ersten Commit, ob `slice-v1-abschluss-sqlite-format` oder `slice-v1-abschluss-sqlite-format-lesen` eine weitere Randform liefert, die dort nicht steht.

**Risiken:**

- Ein Beispiel läuft nur mit einer Umgebung, die die Probe nicht nachbaut (`BEO-REPO/verhalten-nur-unter-linux-geprueft`, 1×) — **Ausgang:** *weiter offen* — alle Proben liefen unter Linux im Container; das Register führt den Eintrag (3× mit `slice-v1-abschluss-sqlite-format-lesen`, nicht zugewiesen) und vermerkt diesen Slice ohne neues Auftreten.
- Der Slice liefert seinen Teil von Handbuch und README nicht vollständig, weil nur die DoD ihn trägt und das Review ihn übersieht (`BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`, 1×) — **Ausgang:** *entfallen* — Die Verifikation bestätigt Handbuch- und README-Teil beider Code-Slices (Punkt 1 und 2 der DoD); die Suche nach den ersetzten Grenz-Sätzen findet nichts Verbliebenes (`grep` nach „nur `yaml`“, „unbekannt“ zu `--format`: nur die Zeilen `PGR-E2001` mit dem Wertebereich und `PGR-E3002`). Das Register vermerkt „Teil geliefert“.

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

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

### Belege des Implementers

**Stand.** Geändert sind `docs/user/benutzerhandbuch.md`, `README.md` und dieser Plan (§3, §7, §8, DoD); die Roadmap-Zeile *In Arbeit* nennt den Slice. Produkt-Code, Tests, Gates und `docs/user/abdeckung-*.md` sind unverändert. Das Handbuch trägt Version 0.3.

**Probe-Aufbau.** `make build` am Stand `221ee96` (Image `pgwire-recorder:dev`, `sha256:d221c756b469…`, Produktcode unverändert). Proben in einem Scratch-Verzeichnis außerhalb des Repos; Netz `impl-sqldoku-net`, ein PostgreSQL 17.11 aus dem gepinnten Image von `harness/mk/integration.mk` (`impl-sqldoku-pg`, Servermodus `trust`), `record`, `replay` und `play` als `docker run … pgwire-recorder:dev` im selben Netz, Client `psql` aus dem PostgreSQL-Image (einfache und erweiterte Anfragen mit `\bind`, ein Fehler der Datenbank, Nicht-ASCII, `PGCLIENTENCODING=LATIN1`). Der Ablauf der Hauptproben: `select 1 as a, 'x' as b;`, `select $1::int + 1 as n, $2::text as t` mit `\bind 41 hallo`, `select * from nirgends;` (Fehlerantwort), `select 'ä€' as u;`. Die Beispiele des neuen Abschnitts liefen **wörtlich** (die Blöcke zogen ein Skript aus dem Handbuch, `record` mit `--format sqlite`, `record --config` mit dem YAML-Block, `replay`, `play`, beide `sqlite3`-Aufrufe) mit dem Binary aus dem Image in einem Container desselben Netzes (Host `postgres`, Port 15432 über `localhost`). `sqlite3` 3.45.1 auf dem Host für die Abfragen. Container, Netz und Scratch-Dateien sind danach entfernt. Nur Linux im Container, lokales Dateisystem (`BEO-REPO/verhalten-nur-unter-linux-geprueft`).

**Proben je Aussage** (`Z` = `record --listen 0.0.0.0:15432 --upstream impl-sqldoku-pg:5432`; „Ablauf“ = die vier Anweisungen oben):

| Aussage im Handbuch | Probe | Ergebnis |
|---|---|---|
| `--format sqlite` schreibt eine SQLite-Datei; ohne `--format` YAML | `Z --output users.sqlite --format sqlite`, `Z --output users.yaml`, Ablauf | `users.sqlite` beginnt mit `SQLite format 3` und NUL; `users.yaml` mit `format: pgwire-recorder` |
| `replay` und `play` verhalten sich mit sqlite wie mit yaml (Gegenprobe: eine abweichende Anfrage) | Ablauf gegen `replay` aus beiden Dateien, Ausgabe von `psql` byte-gleich mit der direkten; `select 99;` gegen beide; `play` beider Dateien nach `drop table` | Ausgaben gleich (`cmp`); `select 99;` in beiden `PGR-E5001`, Exit 5; die Tabelle `pt` hat nach `play` aus beiden Formaten dieselben Zeilen (`1:ä€`, `2:zwei`) |
| Bytes ohne UTF-8 kommen zurück | `PGCLIENTENCODING=LATIN1`, `select 'café'` mit dem Byte `0xE9`; `record` in beide Formate, `replay` | `interaction.sql` trägt `…63 61 66 E9…`; die Wiedergabe liefert `caf\351` in beiden Formaten |
| Datei wird nach der ersten beendeten Verbindung aktualisiert | `Z --format sqlite`, `ls` vor der ersten Verbindung, nach der ersten, nach der zweiten | vorher nicht vorhanden; danach 1, dann 2 Zeilen in `session` |
| `--force` ersetzt, ergänzt nie; vorher unverändert; Format wechselt | vorhandene Datei mit `--force` (sqlite → sqlite, yaml → sqlite, sqlite → yaml), `sha256sum` vor und nach der ersten Verbindung | eine Sitzung (die neue); Datei vor der ersten Verbindung unverändert, danach ersetzt; Format wie gewählt |
| vorhandene Zieldatei ist `PGR-E2002`, gleich welches Format | sqlite auf yaml, yaml auf sqlite, sqlite auf sqlite, je ohne `--force`; Verknüpfung auf sqlite | je `PGR-E2002`, Exit 2 |
| Rechte einer ersetzten Datei bleiben; neue Datei `0666` nach umask | vorhandene Datei `640` und `666`, Verknüpfung auf `600`, je `--force`; neue Datei unter umask `022` | `640`, `666`, `600` bleiben; neu `0644`; die Verknüpfung ist danach eine reguläre Datei, ihr Ziel unverändert |
| `--format`: nur `yaml`, `sqlite`; Priorität Option > Variable > Datei | `SQLite`, `sqlite3`, leer als Option; Variable `sqlite`/`SQLite`/leer; Datei `format: sqlite` mit Variable `yaml`; Variable `sqlite` mit Option `yaml`; Variable `SQLite` mit gültiger Option; `format: SQLite` in der Datei; Option leer mit Variable `sqlite` | ungültige Werte `PGR-E2001` (Option, Variable), `PGR-E2004` (Datei), Exit 2; leere Variable gilt als nicht gesetzt (YAML); Option geht vor; Variable `SQLite` ist auch mit gültiger Option `PGR-E2001`; leere Option `PGR-E2001` |
| `replay` und `play` kennen `--format` nicht | `replay … --format sqlite`, `play … --format sqlite` | je `PGR-E2001` (`flag provided but not defined`) |
| `config show` nennt `PGWIRE_RECORDER_FORMAT` und `format:` | `config show --config` mit `format: sqlite` und gesetzter Variable | Datei, `record:` mit `format: sqlite`, Name der Variable ohne Wert |
| Erkennung nach Inhalt, nicht nach Endung | sqlite-Inhalt als `falsch.yaml`, yaml-Inhalt als `falsch.sqlite`, `replay` | beide gelesen, Ausgabe gleich der direkten |
| Datei ohne SQLite-Kennung ist YAML, auch eine leere, dann `PGR-E3003` | leere Datei, 15 Bytes `SQLite format 3` ohne NUL | `PGR-E3003` (`EOF`; `yaml: unmarshal errors …`) |
| Zielpfad mit `?`, `#`, `%`, Leerzeichen | `--output "/rec/a b?#%x.sqlite"`, dann `replay` | Datei unter genau diesem Namen; Wiedergabe gleich der direkten |
| Datei nur lesend, auch auf schreibgeschütztem Verzeichnis; `play` ändert sie nicht | `replay` mit `-v …:ro`, `sha256sum` vor und nach; `play` mit `sha256sum` | Wiedergabe gleich; Prüfsumme gleich; kein `-journal` entstanden |
| Lesen wartet höchstens fünf Sekunden auf die Sperre, dann `PGR-E3001` | Sperre mit `begin exclusive` (Python, `sqlite3`) für 12 s: `replay`; für 3 s: `replay` | 12 s: nach 5,5 s `PGR-E3001` `database is locked (5) (SQLITE_BUSY)`, Exit 3; 3 s: „replay gestartet“ 30 ms nach dem Ende der Sperre (Zeitstempel beider Seiten) |
| Datei während des Aufzeichnens lesbar | `Z --format sqlite`, eine Verbindung, dann `replay` auf die Datei und `sqlite3` | `replay` beantwortet die Anfrage; `sqlite3` zählt 1, nach der zweiten Verbindung 2 |
| `PGR-E3001`: Journal neben dem Zielpfad, auch mit `--force`, auch ohne Zieldatei; Zieldatei ohne `--force` zuerst `PGR-E2002`; kein Journal-Befund bei yaml | leere Datei `j.sqlite-journal` ohne Zieldatei, mit `--force`; Zieldatei und Journal ohne `--force`; `j3.yaml-journal` bei yaml | `PGR-E3001` („Rollback-Journal eines früheren Laufs“), Exit 3; `PGR-E2002`; yaml startet |
| `sqlite3` macht den Abbruch rückgängig und entfernt das Journal, wenn es Inhalt hat (Nacharbeit F-631) | ein echtes Journal (Python, `os._exit` mitten in einer Transaktion bei `cache_size=1`, 9728 Bytes) neben `d.sqlite`: `record --force` (`PGR-E3001`), dann `sqlite3 "select count(*) from t"`, `ls`, dann `record --force` | `PGR-E3001`, Exit 3; Abfrage liefert `0` (Stand vor der Transaktion), kein Journal mehr; `record --force` startet (Exit 0) |
| eine leere Journal-Datei entfernt `sqlite3` nicht; selbst löschen, dann startet `record` (Nacharbeit F-631, Gegenrichtung) | leere `e.sqlite-journal` neben gültiger `e.sqlite`: `sqlite3`-Abfrage, `ls`, `record --force`; leere Journal-Datei neben anderer Zieldatei und neben fehlender Zieldatei: `record --force` / `record`, dann `rm` des Journals, erneut | Abfrage ohne Fehler, Journal bleibt (`ls` zählt 1), `record --force` weiter `PGR-E3001`; nach dem Löschen startet `record` (Exit 0) in beiden Fällen |
| `interaction.sql` ist bei `extended` `NULL` (Nacharbeit V-179) | `select $1::int` mit `\bind 5 \g` und zwei einfache Anfragen in `sqlite`, dann `select type, sql is null, count(*) … group by 1,2`, `where type='extended' and sql=''`, `typeof(sql)` | `extended|1|1`, `query|0|2`; kein Treffer für `sql=''`; `typeof` ist `null` |
| `sqlite`: erster Schreibvorgang legt über die temporäre Datei an, weitere Verbindungen ergänzen in der Zieldatei selbst (Nacharbeit V-180); bei `yaml` ersetzt jede Aktualisierung die Datei | `stat -c %i` der Zieldatei nach jeder von drei Verbindungen und nach dem Beenden, beide Formate; das Anlegen über die temporäre Datei zeigt die Fehlermeldung `temporäre Datei … nicht anzulegen` (Zeile oben) | `sqlite`: Inode nach allen drei Verbindungen und nach dem Beenden gleich (`15519485`); `yaml`: Inode wechselt (`…486`, `…487`, `…486`, `…487`) |
| Zeile `PGR-E3001` nennt Klasse und Code, nicht den Wortlaut der Sperre (Nacharbeit F-632) | `grep SQLITE_BUSY` im Handbuch | kein Treffer |
| Satz zur Erkennung: nur als YAML Unlesbares ist `PGR-E3003` (Nacharbeit F-633) | leere Datei und 15-Byte-Datei (`PGR-E3003`, Zeilen oben), gültige YAML-Datei ohne SQLite-Kennung (`users.yaml`, `replay` gelesen) | `PGR-E3003` bzw. gelesen |
| „Neben der Zieldatei liegt während des Ergänzens eine Datei `…-journal`“ (F-636) | nicht beobachtbar mit einer Probe von außen | Satz gestrichen |
| leere Aufzeichnung: gültige Datei ohne Sitzung, `replay` `PGR-E3004` | `Z --format sqlite` ohne Verbindung, `sqlite3`, `replay` | `meta` mit `format`, `version`, 0 Sitzungen; `PGR-E3004`, Exit 3 (wie bei yaml) |
| Verbindung ohne Anfrage wird nicht aufgezeichnet | `psql -f /dev/null`, dann `select 1`; beide Formate | je eine Sitzung |
| Beenden mit Frist: abgeschlossene Anfragen bleiben, `PGR-E4006`, Exit 4 | `--shutdown-timeout 1s`, Verbindung mit `select 1;` und `select pg_sleep(30);`, `docker stop` | in `interaction` nur `select 1 as erste;`; `PGR-E4006`, Exit 4 |
| temporäre Datei: bei Fehler `PGR-E3001` ohne Rest; ein alter Rest bleibt unbemerkt | Verzeichnis nach dem Start auf `555`, eine Verbindung, beide Formate; ein Rest `.stale.sqlite.0123456789abcdef.tmp` neben dem Ziel | `PGR-E3001` („temporäre Datei … nicht anzulegen“), Verzeichnis leer, Exit 3; der Rest bleibt, das Log nennt ihn nicht |
| Ziel ist Verzeichnis, Verzeichnis fehlt: `PGR-E3001` | `--output` auf ein Verzeichnis mit `--force`; `--output /rec/fehlt/x.sqlite` | je `PGR-E3001`, Exit 3 |
| `PGR-E3001` beim Lesen: Datei fehlt, keine Rechte, Verzeichnis | `replay` auf `fehlt.db`, Datei `000`, Verzeichnis | je `PGR-E3001` |
| `PGR-E3002` bei unbekannter Version, beide Formate | `meta.version = 2`; `version: 2` in yaml | je `PGR-E3002`, Exit 3 |
| `PGR-E3003` bei sqlite: keine Datenbank, fremde Datei, fehlende Tabelle oder Spalte, abgeschnitten, ungültiger Inhalt | Kopf mit 16 und 40 Bytes, fremde Datenbank, `drop table message`, `drop column sql`, halbe und 300-Byte-Datei, Formatkennung, `version = 'abc'`, Lücke in `sequence`, verwaiste Zeile, JSON `{kaputt`, Schlüssel in `meta` | je `PGR-E3003`, Exit 3 |
| Aufbau der Tabellen und Beispielabfragen | `.schema`, `select` auf alle vier Tabellen; die zwei `sqlite3`-Aufrufe aus dem Handbuch, wörtlich | `meta`, `session`, `interaction`, `message` mit den genannten Spalten; Abfrage 1 liefert `1|1|select 1 as a, 'x' as b;`, Abfrage 2 `1|2|select $1::int + 1 as n, $2::text as t ` |
| Beispiele des neuen Abschnitts | wörtlich ausgeführt (siehe Aufbau) | `record` Exit 0 mit 1 Sitzung, `record --config` Exit 0 mit 1 Sitzung, `replay` beantwortet beide Anfragen, `play` Exit 0 |
| keine rechnerspezifischen Angaben | `grep -a` nach Hostname des Containers, `/rec`, Netzname, Dateiname in `.sqlite` und `.yaml` | 0 Treffer in beiden |

**Nicht in das Handbuch übernommen** (Probe zeigt es nicht oder die Grenze steht nur in der Spezifikation): dass zwei Läufe byte-gleiche SQLite-Dateien liefern oder nicht (in der Probe waren zwei Dateien gleich; zugesagt ist es nicht), die Obergrenze des Speichers beim Lesen, das Lesen einer Datei im WAL-Modus, ein Journal neben der Datei beim Lesen, Windows und macOS, Netzwerk-Dateisysteme, Verhalten von `play` mit einer Aufzeichnung ohne Sitzung.

**Befehlsfolge aus `slice-doku-ist-stand` §7** (mit der benannten Ausnahme `PGWIRE_RECORDER_PASSWORD`, Variablen-Zeile `comm -23 $W/v-hb $W/v-help | grep -vx PGWIRE_RECORDER_PASSWORD`), am Stand dieses Commits:

| Prüfung | Ergebnis |
|---|---|
| Optionen des Handbuchs gegen `--help` | nur `--h`, `--help`, `--version` (Nicht-Optionen im Handbuch), `--rm`, `--name` (`docker run`) |
| Variablen ohne Option | leer |
| Codes außerhalb des Katalogs | leer |
| Mutant: `--format-bogus` angehängt | `--format-bogus` in der ersten Ausgabe |
| Mutant: `PGWIRE_RECORDER_BOGUS` angehängt | in der zweiten Ausgabe |
| Mutant: `PGR-E9999` angehängt | in der dritten Ausgabe |
| Mutant: `--format` aus der Hilfe entfernt | `--format` in der ersten Ausgabe (die Option steht im Handbuch) |

**`grep` über Handbuch und README.** Im Handbuch: kein Treffer für `ADR-`, `SPEC-`, `LH-`, `[Ss]lice`, `Welle`, `Spezifikation`, `Lastenheft`, `[Rr]eview`; kein neuer Satz mit „noch nicht“, „kommt“, „geplant“, „bisher“ (die vorhandene „später“ in der Einleitung gehört zum Ablauf der Wiedergabe, nicht zum Produktstand). Im README bleiben die Verweise auf `spec/` und `docs/plan/` ohne Aussage über einen Stand.

**Namensabgleich.** `docs/user/benutzerhandbuch.md`, `README.md`, `internal/hexagon/model/fehler.go`, `harness/mk/integration.mk`, `done/slice-doku-ist-stand.md` und die zwei Slices `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` in `done/`: per `ls` und `grep` gefunden; die Anker `#das-format-der-aufzeichnung-festlegen` und `#eine-sqlite-aufzeichnung-abfragen` löst `make docs-check` auf (`d-check: 567 Datei(en) geprüft, 0 Befund(e)`).

**Hinweise an Architect und Planner** (keine Randform dieses Slice, kein Befund gegen die Spezifikation):
1. *Meldungen von `yaml` bei einer SQLite-Datei ohne gültige Kopfzeile:* Eine 15-Byte-Datei `SQLite format 3` meldet `PGR-E3003` mit einem Stück ihres Inhalts (`cannot unmarshal !!str …`); das ist die bekannte Lücke `BEO-REPO/yaml-meldung-nennt-wert-und-zeile`. Das Handbuch sagt über den Inhalt der Meldungen nichts.
2. *`play` mit einer Aufzeichnung ohne Sitzung:* endet mit Exit 0 und den Zeilen „play gestartet“, „play beendet“, ohne Warnung. Die Spezifikation sagt, `play` führe Sessions ohne Interaktion nicht aus; ob eine Aufzeichnung ganz ohne Sitzung eine Meldung verdient, ist nicht entschieden. Das Handbuch nennt das Verhalten nicht.

### Closure

**Review und Verifikation.** Review (`24edefb`): F-631 (MEDIUM), F-632, F-633 (LOW), F-634 bis F-636 (INFO); F-631 bis F-633 und F-636 erledigt in `ab22809`. Verifikation (`f9e4417`): DoD-Punkte 1 und 2 bestätigt; V-179 und V-180 (LOW) erledigt in `4932680` mit den Proben in der Tabelle oben; V-181 (INFO) siehe unten.

**Entscheidungen zu den Übergaben.**

- **F-634** (`play` mit einer Aufzeichnung ohne Sitzung: Exit 0 ohne Meldung, `replay` meldet `PGR-E3004`): Das Handbuch sagt nichts Falsches. Die Spezifikation regelt es nicht: `PGR-E3004` ist für `replay` definiert (`LH-FA-03.a`), für `play` steht nur „führt Sessions ohne Interaktion nicht aus“. Eine **benannte Spezifikationslücke**, im Register als `BEO-REPO/play-ohne-sitzung-ohne-meldung` (neu, 1×, offen). Keine Zuweisung an einen Slice: Kein Slice, der `play` berührt, führt den Gegenstand (`slice-v1-abschluss-antwortvergleich` ändert den Ablauf von `play`, schließt aber die Zusage über Sitzungen nicht an), und die Entscheidung (ablehnen, warnen, schweigen) gehört zuerst in die Spezifikation, bei Lastenheft-Bezug dem Nutzer. Eine Zuweisung ohne annehmenden Nehmer wäre eine Adresse, die ins Leere zeigt (`AGENTS.md` §3.13); sie folgt bei 3× durch den Lese-Schritt.
- **F-635, V-181** (Meldung zu einer 15-Byte-Datei nennt ein Stück des Dateiinhalts): `BEO-REPO/yaml-meldung-nennt-wert-und-zeile`, zweites Auftreten (1× → 2×); die Entscheidung des Nutzers vom 2026-10-11 (Regel vorerst nur für `sqlite`) gilt, das Handbuch macht dazu keine Zusage. Keine Handlung an diesem Slice.
- **F-633** (Satzbezug „sie ist dann beschädigt“): Einzelfall der Verständlichkeit, behoben; kein Register-Eintrag.

**Risiken aus §6 — Ausgang je Risiko** (geschlossene Menge): siehe §6; Beispiel läuft nur unter Linux: *weiter offen*; Teil von Handbuch und README nicht geliefert: *entfallen*.

**Register** (Zähler = Dateien unter `evidence/`, gelesen am Stand dieser Closure):

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` — neue Datei `evidence/slice-v1-abschluss-sqlite-format-doku.md` (F-631, F-632, F-636, V-179, V-180): 38 → 39. Verkörpert (`AGENTS.md` §3.11); `state.md` führt zusätzlich die zwei Absätze unten. Die Klasse „Satz ohne eigene Probezeile“ (F-636, V-180) und „Zusage weiter als die Probe“ (F-631, V-179) sind hier dieselbe Beobachtung (ein Satz sagt zu, was keine Probe zeigt); ein eigener Eintrag entsteht nicht. `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` meint etwas anderes (der Code-Slice liefert den Handbuch-Teil nicht).
- `BEO-REPO/yaml-meldung-nennt-wert-und-zeile` — neue Datei `evidence/slice-v1-abschluss-sqlite-format-doku.md` (F-635, V-181): 1 → 2. Offen, nicht zugewiesen.
- `BEO-REPO/play-ohne-sitzung-ohne-meldung` — **neu**, `evidence/slice-v1-abschluss-sqlite-format-doku.md` (F-634): 1×, offen.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` — kein neues Auftreten; vermerkt in `state.md`, bleibt 3×, nicht zugewiesen (der Lese-Schritt der Closure von welle-v1-abschluss weist zu).
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` — kein Auftreten; vermerkt „Teil geliefert“, bleibt 1×, offen.
- Kein nicht verkörperter Eintrag erreicht mit diesem Slice 3× neu. Retirement-Check: `AGENTS.md` §3.11 (Handbuch gilt nur mit Probe) ist wieder aufgetreten (F-631, F-632, V-179, V-180); die Regel bleibt, ihr Ergebnis steht im Lerneintrag.

**Lerneintrag — geschärfte Regel (zwei Absätze) und benannte Spezifikationslücke.**

(a) Der Hinweis zum Journal sagte zu: „`sqlite3` macht den Abbruch rückgängig und entfernt die Datei, danach startet `record`“. Geprobt war ein Journal mit Inhalt; ein Journal der Länge 0 bleibt liegen, und der Text führte in eine Sackgasse (F-631). Die Kettenregel in Schritt 19 nannte nur Aussteller, Herkunft und Zuständigkeit; eine Anleitung mit Folgeschritten fiel nicht darunter. Sie nennt jetzt auch sie: Ausgangszustand mit Inhalt und ohne sind zwei Fälle, jeder Schritt der Kette hat seine Probe. `liegt in` `.claude/commands/implement-slice.md` Schritt 19 — Herkunfts-Anker `seit slice-v1-abschluss-sqlite-format-doku`.

(b) Zwei Sätze zur Art des Schreibens hatten keine Zeile in der Probe-Tabelle (F-636, V-180), ein Wort („leer“) nannte einen Zustand anders, als die Probe ihn zeigt (V-179), und ein Wortlaut der Bibliothek stand im Handbuch (F-632, Wiederholung von F-615). Die Proben des Implementers waren je Aussage geführt, aber nicht je Satz; die Prüfung am Satz fand erst Review und Verifikation. Schritt 17 sagt jetzt: Jeder Satz mit einer Verhaltensaussage hat seine eigene Zeile, auch ein Nebensatz; das Zustandswort folgt der Probe. `liegt in` `.claude/commands/implement-slice.md` Schritt 17 — Herkunfts-Anker `seit slice-v1-abschluss-sqlite-format-doku`.

(c) Benannte Spezifikationslücke: Die Spezifikation regelt `play` mit einer Aufzeichnung ohne Sitzung nicht (`BEO-REPO/play-ohne-sitzung-ohne-meldung`); `PGR-E3004` ist nur für `replay` definiert. Zuständig: Architect mit dem Nutzer, bevor ein Slice, der `play` berührt, das Verhalten festlegt.

**Paarungen.** (a) Anker: Beide Zielorte (Schritt 17 und Schritt 19 in `.claude/commands/implement-slice.md`) tragen `seit slice-v1-abschluss-sqlite-format-doku` (`grep`), `state.md` von `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` nennt sie. (b) Folge-Slice: Diese Closure weist keinem Slice etwas zu; `slice-v1-abschluss-zeitangaben-doku` und die übrigen Doku-Folge-Slices bleiben, wie sie sind. (c) Register: die fünf genannten Kennungen existieren als Verzeichnis mit nicht leerem `evidence/`.

**Nachzählen.** Keine Sendung in einen Nehmer, daher kein Nachzählen von Liefer-Punkten und Schichten. Dieser Slice hat 2 Liefer-Punkte und eine Schicht (Dokumentation), wie in §8.

**Stand der Gates.** `make gates` auf dem sauberen Baum `4932680` (Stand nach der Nacharbeit zu V-179 und V-180, vor diesem Eintrag): Exit 0. Die Commits dieser Closure ändern nur Plan, Register, `implement-slice.md` und Roadmap; `make gates` auf dem Stand nach dem letzten Commit ist die Messung, die der Bericht des Planners nennt.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` gesichtet (Zähler = Dateien unter `evidence/`). Treffer:

- `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (verkörpert, `AGENTS.md` §3.11) — genau der Gegenstand dieses Slice: Das Handbuch sagt nur zu, was das Binary zeigt.
- `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht` (1×, offen) — dieser Slice ist die Adresse des Teils von `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen`; seine Beobachtung ist das Risiko in §6.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (1×, offen) — Proben laufen nur unter Linux im Container (§6).
- `BEO-REPO/schichtteilung-je-plan-verschieden` (verkörpert, `AGENTS.md` §3.13) — dieser Slice ist die Folge der Zählung.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Liefer-Punkte: 2. Schichten nach der Zählung in `AGENTS.md` §3.13: eine (Dokumentation), kein Code. Nachgezählt bei der Umsetzung: Handbuch (ein Liefer-Punkt, mit dem neuen Abschnitt) und README (ein weiterer); der Diff berührt nur `docs/user/benutzerhandbuch.md`, `README.md`, diesen Plan und die Roadmap-Zeile.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
