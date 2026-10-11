# Verifikation: slice-v1-abschluss-sqlite-format-doku — 2026-10-11

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob er richtig gebaut ist. Ob das Richtige gebaut wurde, prüft der Validator; den Diff gegen Plan, Entscheidungen und Hard Rules hat der Reviewer geprüft (Report `docs/reviews/2026-10-11-review-slice-v1-abschluss-sqlite-format-doku.md`, F-631 bis F-636); die Übergaben prüfe ich nach. Ich repariere nichts.

**Review-Art:** Verifikation (DoD- und Entscheidungs-Konformität, Plan gegen Handbuch und README als Diff) in frischem Kontext.

**Gegenstand:** Slice-Plan `slice-v1-abschluss-sqlite-format-doku` (Kopf und §1 bis §8) gegen den Gesamt-Diff `221ee96..ab22809` (5 Dateien: `README.md`, `docs/user/benutzerhandbuch.md`, `docs/plan/planning/in-progress/roadmap.md`, der Plan, der Review-Report; +297 −26). Implementer `f2e4894`, `3234506`; Review `24edefb`; Nacharbeit `ab22809`. Bezug: [LH-FA-22](../../spec/lastenheft.md#lh-fa-22--wählbares-aufzeichnungsformat), [ADR-0013](../plan/adr/0013-sqlite-recording-backend.md).

**Skill:** `.harness/skills/reviewer.md` (nicht angewandt; Rolle ist Verifier)
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-11

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link. Ein
> `pfad`-Feld auf den **geprüften Gegenstand** zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- Plan ganz am Stand `ab22809`, besonders DoD (§2), §3, §6 und §7 (*Belege des Implementers*: Probe-Matrix, Befehlsfolge, `grep`, Namensabgleich, Hinweise).
- `AGENTS.md` §3.9, §3.11 (*Handbuch und README beschreiben den Ist-Zustand*), §3.13; `slice-doku-ist-stand` §6 und §7 (Befehlsfolge).
- Review-Report F-631 bis F-636.

## 1. DoD-Liefer-Punkte

**Methode.** Das Binary stammt aus einer frischen Kopie (`git archive HEAD`, `ab22809`, Stufe `runtime`, Image `ver-sqldoku-bin:dev`, Image-ID gleich der des Implementers, `d221c756b469…`). Docker-only, Scratch-Präfix `ver-sqldoku-`: eigenes Netz, ein `postgres:17-alpine` mit dem Digest aus `harness/mk/integration.mk` (`trust`, per Netz-Alias `postgres` für die Handbuchbeispiele), Client `psql` aus demselben Image (einfache Anfragen, `\bind`, Fehlerantwort, Nicht-ASCII), `sqlite3` 3.45 und `python3` auf dem Host für Abfragen, Sperre und das abgebrochene Journal. Die Beispiele sind **wörtlich per `awk` aus den Blöcken des Handbuchs geschnitten**, nicht abgetippt. Container, Netz und Image sind danach per Name entfernt (Zähler 0/0/0). Dateien im Scratch außerhalb des Repos.

### Punkt 1: Handbuch beschreibt, was `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` liefern — **bestätigt**

Jede neue oder geänderte Verhaltensaussage einzeln, in beide Richtungen (Kettenprobe):

| Aussage | Probe (selbst) | Ergebnis |
|---|---|---|
| `--format sqlite` schreibt SQLite, ohne `--format` YAML | `record` mit und ohne Option, Ablauf (einfach, `\bind`, Fehler, `ä€`) | Datei beginnt mit `SQLite format 3` und NUL; die andere mit `format: pgwire-recorder`; Ausgabe von `psql` durch `record` byte-gleich der direkten (`cmp`) |
| `replay` und `play` verhalten sich mit sqlite wie mit yaml | `replay` aus beiden Dateien, `cmp` gegen die direkte Ausgabe; Gegenprobe `select 99;` in neuer Verbindung; `play` aus der sqlite-Datei | Ausgaben gleich; Gegenprobe `PGR-E5003` in beiden Formaten (kein falsches Grün); `play` Exit 0 |
| Wert `yaml`/`sqlite`, sonst `PGR-E2001` (Option, Variable); Datei `PGR-E2004`; Standard `yaml` | `SQLite`, `sqlite3`, leer als Option; Variable `sqlite3`; Variable `SQLite` mit gültiger Option; `format: SQLite` in der Datei | je Exit 2: Option und Variable `PGR-E2001`, auch die Variable mit gültiger Option; Datei `PGR-E2004` |
| leere Variable gilt als nicht gesetzt | `PGWIRE_RECORDER_FORMAT=` mit `config show` | Exit 0 ohne Befund |
| Priorität Option > Variable > Datei | Datei `format: sqlite` gegen Variable und Option, vier Kombinationen, Format der entstehenden Datei gelesen | nichts/nichts: SQLite; Variable `yaml`: YAML (Variable geht der Datei vor); Variable `sqlite` + Option `yaml`: YAML; Variable `yaml` + Option `sqlite`: SQLite (Option geht vor) |
| `replay` und `play` kennen `--format` nicht | beide mit `--format sqlite` | je `PGR-E2001`, „flag provided but not defined“, Exit 2 |
| vorhandene Zieldatei: `PGR-E2002`, gleich welches Format | yaml→sqlite, sqlite→yaml, sqlite→sqlite, yaml→yaml, je ohne `--force` | je `PGR-E2002`, Exit 2 |
| `--force` ersetzt, ergänzt nie; bis zur ersten Aktualisierung unverändert; Format wechselt | sqlite→sqlite, yaml→sqlite, sqlite→yaml; `sha256sum` vor Verbindung und nach der ersten | vor der Verbindung byte-gleich; danach genau die neue Sitzung (1, nicht 1+n); Format wie gewählt |
| Datei entsteht bei sqlite mit dem ersten Schreibvorgang | `ls` vor/nach der ersten Verbindung | vorher nicht vorhanden |
| erster Schreibvorgang über temporäre Datei, danach in der Zieldatei ergänzt (sqlite); yaml jedes Mal ersetzt | Inode der Zieldatei nach 1. und 2. beendeter Verbindung | sqlite: Inode gleich, Sitzungen 1 → 2; yaml: Inode wechselt |
| Aufzeichnung ohne Verbindung: gültige Datei ohne Sitzung, `replay` `PGR-E3004` | `record` ohne Verbindung, `sqlite3`, `replay` | `meta` mit `format`, `version`; 0 Sitzungen; `PGR-E3004` |
| Datei lesbar während des Aufzeichnens | `sqlite3` und `replay` nach der ersten, `sqlite3` nach der zweiten Verbindung | 1 Sitzung, `replay` antwortet; später 2 |
| Format nach Inhalt, nicht nach Name | sqlite-Inhalt als `falsch.yaml`, yaml-Inhalt als `falsch.sqlite` | beide gelesen, Antworten gleich der direkten |
| Datei ohne SQLite-Kennung ist YAML; als YAML unlesbar `PGR-E3003`, auch leer | leere Datei; 15 Byte `SQLite format 3`; gültige YAML-Datei ohne Kennung | `PGR-E3003` (`EOF`; `yaml: unmarshal errors`); die YAML-Datei wird gelesen |
| nur lesend, auch auf `:ro`; `play` ändert nichts | `replay` und `play` mit Verzeichnis `:ro` und Datei `444`, Prüfsumme vorher/nachher | gleich; kein `-journal` entstanden |
| Sperre höchstens fünf Sekunden, dann `PGR-E3001` | `begin exclusive` 12 s: `replay`; 3 s: `replay` | 12 s: `PGR-E3001` nach 5,4 s (inklusive Containerstart); 3 s: „replay gestartet“ nach 2,3 s |
| `PGR-E3001` bei `<Zieldatei>-journal`, auch mit `--force`, auch ohne Zieldatei; mit Zieldatei ohne `--force` zuerst `PGR-E2002`; bei yaml kein Befund | leeres Journal: ohne Zieldatei, mit `--force`; Zieldatei ohne `--force`; `j3.yaml-journal` bei yaml | `PGR-E3001` („Rollback-Journal eines früheren Laufs“), Exit 3; `PGR-E2002`; yaml startet |
| F-631 Hinweis, Richtung 1: Journal mit Inhalt macht `sqlite3` rückgängig und entfernt es, danach startet `record` | echtes Journal (Python, `os._exit` mitten in der Transaktion, `cache_size=1`) neben `d.sqlite`: `record --force`, `sqlite3 "select …"`, `ls`, `record --force` | `PGR-E3001`; Abfrage liefert den Stand vor der Transaktion (300 Zeilen `x…`), Journal weg; `record` startet |
| F-631 Hinweis, Richtung 2: leeres Journal entfernt `sqlite3` nicht; selbst löschen, dann startet `record` | 0-Byte-Journal neben gültiger Datei: `sqlite3`, `ls`, `record --force`; `rm`, `record --force` | Abfrage ohne Fehler, Journal bleibt, `PGR-E3001`; nach `rm` startet `record` |
| `PGR-E3001` weiter: Datei fehlt, keine Rechte, Verzeichnis; bei `record` Verzeichnis als Ziel, Verzeichnis fehlt, nicht beschreibbar (temporäre Datei) | Aufrufe wie genannt | je `PGR-E3001`, Exit 3 (beim Schreiben ohne Rest im Verzeichnis) |
| `PGR-E3002` für beide Formate | `meta.version = 2`; `version: 2` | je `PGR-E3002` |
| `PGR-E3003` bei sqlite: keine/fremde Datenbank, Tabelle oder Spalte fehlt, abgeschnitten, Inhalt ungültig (Formatkennung, Nummerierung, Eintrag ohne Sitzung) | Kopf 16/40/300 Byte; fremde DB; `drop table message`; `drop column sql`; Formatkennung; Lücke in `sequence`; `session.id` verschoben | je `PGR-E3003`; die Lücke meldet „Interaktion 2 trägt die Nummer 3“, der verwaiste Eintrag „ohne Session“ |
| Aufbau und Spalten der vier Tabellen | `.schema`; `select` auf Gruppen von `direction`, `kind`, `type`, `meta.key` | stimmt mit der Tabelle überein; bei einfacher Anfrage in `message` nur `server` (10 Zeilen, keine `client`), bei `extended` beide |
| Beispiele des Abschnitts | Blöcke 1 bis 4 wörtlich: `record …` (Block 1, 5 Zeilen), `record` mit `.pgwire-recorder.yaml` aus Block 2, `replay` und `play` aus Block 3, beide `sqlite3`-Aufrufe aus Block 4 | `record` beide Exit 0 mit Datei `users.sqlite`; `replay` beantwortet beide Anfragen; `play` Exit 0; Abfrage 1 `1|1|select 1 as a, 'x' as b;`, Abfrage 2 `1|2|select $1::int + 1 as n, $2::text as t ` |
| FAQ: Datei ohne rechnerspezifische Angaben, beide Formate | `grep -a` nach Pfad, Netzname, Containername, Dateiname | 0 Treffer in beiden |
| Verbindung ohne Anfrage wird nicht aufgezeichnet | `psql` ohne Eingabe, dann `select 1` | 1 Sitzung |
| README: `record` mit `--format sqlite`, `replay` und `play` erkennen das Format selbst | wie oben | wahr |

**Urteil Punkt 1: bestätigt.** Keine Aussage musste entfallen; zwei Formulierungs-Befunde unten (V-179, V-180).

### Punkt 2: `README.md` nennt den Stand im Ist-Zustand — **bestätigt**

Der Satz in *Was kann ich heute tun?* nennt Format `sqlite` und Erkennung durch `replay` und `play`; jede Teilaussage ist durch die Zeilen oben belegt. Keine Chronik („noch nicht“, „kommt“, „geplant“, „bisher“: kein Treffer in den neuen Zeilen), die Verweise auf `spec/` und `docs/plan/` sind unverändert und ohne Stand-Aussage.

### Punkt 3: Aktive Suche nach Chronik, Zielstand, Verweisen, Zusagen

- `grep -n -E 'ADR-|SPEC-|LH-|ARC-|[Ss]lice|Welle|Spezifikation|Lastenheft|[Rr]eview'` im ganzen Handbuch: kein Treffer.
- Neue Zeilen des Diffs gegen `noch nicht|kommt|geplant|bisher|später|künftig|zukünftig|demnächst|Windows|macOS|groß|SQLITE_BUSY`: kein Treffer (ein Muster-Treffer auf „rechnerspezifischen“, falsch positiv). `SQLITE_BUSY` im Handbuch: 0 Treffer (F-632).
- Keine Zusage zu Windows, macOS oder großen Dateien; die Grenze nennt „unter Linux im Container auf einem lokalen Dateisystem geprüft, für andere Systeme und Netzwerk-Dateisysteme nicht geprüft“, in Präsens und mit Probe belegt (alle Proben dieser Verifikation: Linux, Container, lokales Dateisystem).
- Veraltete Grenz-Sätze: `grep -i 'yaml-aufzeichnung\|nur yaml\|yaml-datei'` über Handbuch und README zeigt nur die neuen Sätze; kein Satz sagt noch, die Aufzeichnung sei immer YAML.
- Meldungen nennen nie einen Wert (Spezifikation: für Datei und Umgebungsvariable; eine Meldung zur Kommandozeile darf ihn nennen): Variable `geheimwert` und `format: geheimwert` in der Datei lösen `PGR-E2001` bzw. `PGR-E2004` ohne den Wert aus; die Option nennt ihn (`invalid value "SQLite"`), was die Spezifikation erlaubt. Das Handbuch macht dazu keine Zusage. Restfall siehe V-181.

### Punkt 4: Mutationen am Dokument (nur in Kopien)

Die Befehlsfolge aus `slice-doku-ist-stand` §7 mit der benannten Ausnahme `grep -vx PGWIRE_RECORDER_PASSWORD` am Stand `ab22809`: Optionen außerhalb der Hilfe `--h --help --name --rm --version` (Nicht-Optionen bzw. `docker run`); Variablen leer; Codes außerhalb des Katalogs leer. Mutanten in Kopien des Handbuchs:

| Mutant | Ergebnis |
|---|---|
| Option `--format-bogus` angehängt (falsche Option) | rot: `--format-bogus` in der ersten Ausgabe |
| `--format sqlite` zu `--formt sqlite` (Tippfehler im Text) | rot: `--formt` in der ersten Ausgabe |
| `PGR-E9999` und `PGR-E3005` angehängt (falscher Code) | rot: beide in der dritten Ausgabe |
| `PGR-E3003` zu `PGR-E3033` | rot: `PGR-E3033` in der dritten Ausgabe |
| `PGWIRE_RECORDER_FORMAT` zu `…FORMT` | rot: in der zweiten Ausgabe |
| `--format` aus der Hilfe entfernt | rot: `--format` und `PGWIRE_RECORDER_FORMAT` in Ausgabe 1 und 2 |
| Code-Zuordnung im Text: `ist er PGR-E2004` → `PGR-E2001`; `kennen --format nicht (PGR-E2001)` → `PGR-E2004`; `bricht ab (PGR-E2002)` → `PGR-E2003` | je rot gegen die Codes, die das Binary in den Proben nannte (Original grün) |
| Beispiel Block 1 mit `--formt`, `--outpt`, `--format sqlit` | je Exit 2 `PGR-E2001` (das Beispiel läuft nicht mehr) |

Die Befehlsfolge prüft Namen, nicht Zuordnung; die Zuordnungs-Mutanten (vorletzte Zeile der Tabelle) laufen gegen Probe-Codes und sind keine Wiederholung der Befehlsfolge.

### Übrige DoD-Punkte

| Punkt | Urteil | Beleg |
|---|---|---|
| `make gates` grün | **bestätigt** | eigener Lauf am sauberen Baum nach dem Commit dieses Reports; das Ergebnis steht im Bericht an den Planner (der Report selbst kann es nicht tragen, es folgt kein Commit) |
| Review durchgeführt, Report unter `docs/reviews/` | **bestätigt** | `docs/reviews/2026-10-11-review-slice-v1-abschluss-sqlite-format-doku.md`, anderer Lauf (`24edefb`), kein Self-Review |
| Closure-Notiz mit Lerneintrag | **offen, Pflicht des Planners** | §7 trägt *Belege des Implementers*, keine Closure-Notiz |
| Beobachtungs-Register fortgeschrieben | **offen, bei Closure** | keine Datei unter `evidence/` für diesen Slice |
| Jedes Risiko aus §6 trägt einen Ausgang | **offen, bei Closure** | beide Risiken stehen mit „offen bis Closure“ |
| Drei Paarungen | **offen, bei Closure** | nicht Teil der Liefer-Punkte |

Die DoD-Häkchen im Plan sind noch leer; das gehört zur Closure und ist keine Abweichung (Eingang 5 des Auftrags).

## 2. Übergaben des Reviews (F-631 bis F-636)

| ID | Übergabe | Stand am HEAD | Urteil |
|---|---|---|---|
| F-631 (MEDIUM) | Journalhinweis nur in einer Richtung geprobt | Satz nennt Journal mit Inhalt und leeres getrennt; beide Richtungen selbst gefahren (echtes Journal, 0 Byte) | behoben |
| F-632 (LOW) | `SQLITE_BUSY` im Handbuch weiter als der Test | Zeile `PGR-E3001` nennt Klasse und Code; `grep` 0 Treffer | behoben |
| F-633 (LOW) | „sie ist dann beschädigt“ mehrdeutig | Satz: „meldet `PGR-E3003`, wenn sie als YAML nicht lesbar ist“; gültige YAML-Datei ohne Kennung selbst gelesen | behoben |
| F-634 (INFO) | `play` mit Aufzeichnung ohne Sitzung: Exit 0 ohne Meldung | Handbuch sagt dazu nichts; an den Planner | offen beim Planner |
| F-635 (INFO) | Meldung nennt Dateiinhalt | Handbuch sagt dazu nichts; siehe V-181 | offen beim Planner |
| F-636 (INFO) | Satz „während des Ergänzens eine Datei `-journal`“ ohne Probezeile | Satz gestrichen (`grep` 0 Treffer) | behoben |

## 3. Hexagon und Architektur-Gate

Nicht berührt: Der Diff enthält keinen Produkt-Code, keinen Test, kein Gate, keine Spec-Stelle und keine `docs/user/abdeckung-*.md` (`git diff --name-only 221ee96..HEAD`). `make abdeckung-check` Exit 0 am unveränderten Baum; `make a-check` und `make a-check-negativ` sind Teil des Gate-Laufs unten.

## 4. Plan gegen Handbuch und Code

| Prüfung | Ergebnis |
|---|---|
| Kopf: `Berührte Spec-Stellen: —` | stimmt; kein Spec-Diff; `make kopf-check` im Gate-Lauf |
| §1: nur Handbuch, README, Plan (und Roadmap-Zeile) | stimmt; `benutzerhandbuch-standard.md` und Abdeckungstabellen unverändert |
| §3: Handbuch-Abschnitte (§1, §3, §4 samt zwei neuer Abschnitte, Tabelle §5, Zeilen §7, FAQ, Glossar) und README | stimmen mit dem Diff überein; die Zeilen `PGR-E2001`, `PGR-E2004`, `PGR-E3001` bis `E3003` sind geändert, `E3004` unverändert |
| §3 Ansatz: Befehlsfolge mit Ausnahme | gefahren, siehe Punkt 4 |
| §6 Randformen (*Teilweise geliefert*, *Verhalten neben der Spezifikation*, *Verweise im README*) | die Grenze steht im Präsens und ist beobachtbar; zwei Funde an den Planner stehen in §7 (*Hinweise*), keine im Handbuch |
| §8: 2 Liefer-Punkte, 1 Schicht (Dokumentation) | stimmt: Handbuch, README; der Diff berührt nur Doku und Planung |
| Roadmap-Zeile *In Arbeit* nennt den Slice | stimmt (`roadmap.md`, Schritt 18) |
| Namensabgleich §7 | `docs/user/benutzerhandbuch.md`, `README.md`, `internal/hexagon/model/fehler.go`, `harness/mk/integration.mk`, die Pläne der drei Nachbar-Slices und die Anker `#das-format-der-aufzeichnung-festlegen`, `#eine-sqlite-aufzeichnung-abfragen` vorhanden (`make docs-check` im Gate-Lauf) |
| Commit-Messages `221ee96..HEAD` | alle vier nennen `slice-v1-abschluss-sqlite-format-doku` und `LH-FA-22`; keine `SPEC-`/`ARC-`-Kennung ([ADR-0035](../plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md)) |
| `AGENTS.md` §3.3 | kein Move im Diff |
| Beleg „Probe je Aussage“ im Plan §7 | vorhanden und von mir nachgefahren (Abschnitt 1); die Zahlen stimmen (12 s Sperre: 5,5 s beim Implementer, 5,4 s bei mir) |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-179 | LOW | Die Tabelle nennt für `interaction.sql` bei `extended` „leer“. Die Datei trägt dort `NULL` (`quote(sql)` liefert `NULL`, nicht `''`); wer `where sql = ''` abfragt, findet nichts. Die Aussage ist als Umgangssprache wahr, als SQL-Aussage zu weit. Der Plan-Beleg (§7) prüft `sql` bei `extended` nicht. | `AGENTS.md` §3.11 (Zusage weiter als die Probe) | `docs/user/benutzerhandbuch.md` · „(der Text einer einfachen Anfrage, bei `extended` leer)“ | ja (`sqlite3 … "select quote(sql) from interaction where type = 'extended'"`) | Zusage im Handbuch enger zu fassen |
| V-180 | LOW | Der Satz, bei `sqlite` lege der erste Schreibvorgang die Datei über die temporäre Datei an und jede weitere beendete Verbindung werde „in der Zieldatei selbst ergänzt“, hat in der Probe-Matrix (§7) keine eigene Zeile; die Matrix zählt nur Sitzungen nach 1 und 2 Verbindungen. Der Satz ist wahr (Inode der Zieldatei bleibt, bei `yaml` wechselt er); der Beleg liegt hier, nicht im Plan. Dieselbe Klasse wie F-636. | `AGENTS.md` §3.11 · Review F-636 | `docs/user/benutzerhandbuch.md` · „Das gilt für `yaml` bei jeder Aktualisierung; bei `sqlite` legt der erste Schreibvorgang die Datei auf diesem Weg an“ | ja (`stat -c %i` nach der 1. und 2. Verbindung) | Satz ohne eigene Probezeile |
| V-181 | INFO | Die Meldung `PGR-E3003` zu einer Datei aus den ersten 15 Bytes einer SQLite-Kennung nennt ein Stück des Dateiinhalts („cannot unmarshal !!str `SQLite ...`“). Das ist die bekannte Lücke `BEO-REPO/yaml-meldung-nennt-wert-und-zeile` (Review F-635); das Handbuch sagt nichts über den Inhalt der Meldungen und macht keine Zusage, die dadurch falsch würde. Zuständig: Planner (Register). | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · Zeile `PGR-E3003` (keine Aussage) | ja | Meldung nennt Dateiinhalt, Handbuch schweigt |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| DoD-Punkte 1 und 2 (Liefer-Punkte), jede Verhaltensaussage in beide Richtungen | geprüft, ohne Befund außer V-179 und V-180 |
| Priorität Option > Variable > Datei, leere Variable, ungültige Variable mit gültiger Option | geprüft, vier Kombinationen am Format der entstehenden Datei |
| `--force` (drei Richtungen), Datei bis zur ersten Aktualisierung unverändert | geprüft, ohne Befund |
| `PGR-E2002` vor `PGR-E3001`; Journal mit Inhalt und ohne (echtes Journal) | geprüft, ohne Befund |
| Sperre fünf Sekunden (12 s und 3 s), lesen auf `:ro`, Prüfsumme gleich, kein Journal | geprüft, ohne Befund |
| Beispiele (Blöcke 1 bis 4 wörtlich), beide `sqlite3`-Aufrufe | geprüft, laufen |
| Tabellen- und Spaltenbeschreibung | geprüft, ohne Befund außer V-179 |
| Befehlsfolge und Mutanten (falscher Code, falsche Option, Tippfehler, Variable, Hilfe ohne `--format`) | geprüft, je Mutant rot |
| Chronik, Zielstand, Verweise auf Spezifikation, Lastenheft, ADRs, Slices, Wellen, Reviews (`grep`) | geprüft, ohne Befund |
| Keine Zusage zu Windows, macOS, großen Dateien | geprüft, ohne Befund |
| Meldungen ohne Wert bei Datei und Umgebungsvariable | geprüft, ohne Befund (Restfall V-181) |
| Review-Übergaben F-631 bis F-636 | geprüft; F-634 und F-635 offen beim Planner |
| Commit-Messages ohne `SPEC-`/`ARC-` | geprüft, ohne Befund |
| `make abdeckung-check` am unveränderten Baum | Exit 0; keine Abdeckungsdatei im Diff |
| Nicht geprüft | Windows, macOS, Netzwerk-Dateisysteme, WAL-Modus, `sqlite3` ohne JSON-Funktionen (der Handbuchsatz dazu ist nicht probiert), Zielpfade mit `?`, `#`, `%` (nicht im Handbuch); alle Proben unter Linux im Container (`BEO-REPO/verhalten-nur-unter-linux-geprueft`) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Zusage im Handbuch weiter als die Probe (`NULL` als „leer“) · Satz ohne eigene Probezeile · Meldung nennt Dateiinhalt, Handbuch schweigt

## Verdikt

**Merge-blockierend:** nein. Die zwei Liefer-Punkte der DoD sind bestätigt, jede Verhaltensaussage trägt eine Probe in beide Richtungen, die Mutanten machen die Proben rot, `make gates` läuft nach dem Commit dieses Reports am sauberen Baum (Ergebnis im Bericht an den Planner). Offen sind nur die Closure-Pflichten, die der Plan auf die Closure schiebt.

**Übergabe:**

- **V-179** und **V-180** (LOW) an den Implementer: „leer“ durch „ohne Wert (`NULL`)“ ersetzen oder enger fassen; für V-180 eine Probezeile (Inode) in §7 aufnehmen. Keine Blockade der Closure.
- **V-181** (INFO) und **F-634**, **F-635** an den Planner (Register: `BEO-REPO/yaml-meldung-nennt-wert-und-zeile`; Verhalten von `play` bei einer Aufzeichnung ohne Sitzung).
- Für die Closure: Closure-Notiz mit Lerneintrag; Register fortschreiben (Klassen der Review und dieses Laufs); Ausgänge der zwei Risiken aus §6 (`BEO-REPO/verhalten-nur-unter-linux-geprueft`: alle Proben unter Linux im Container, Handbuch nennt die Grenze; `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`: der Teil von `slice-v1-abschluss-sqlite-format` und `slice-v1-abschluss-sqlite-format-lesen` ist geliefert); DoD-Häkchen.
- Dieser Report ist ein Lauf-Beleg; Aufbau (Container, Netz, Image mit Präfix `ver-sqldoku-`) ist entfernt.
