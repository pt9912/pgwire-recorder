# Verträge der Harness-Werkzeuge

Dieser Ordner legt die Werkzeuge der Prüfumgebung des Projekts fest, die
[`LH-QA-07`](../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) in Messmethode (4)
nennt: was ihre Prüfung zusagt, woran sie rot wird und was sie nicht prüft. Keine
Kennung dieses Ordners ist eine Zusage des Produkts, und keine erweitert eine
Anforderung des Lastenhefts. Die Schnittstelle eines Werkzeugs ist sein
`make`-Target. Die Punkte einer Kennung sind nummeriert; die Nummer ist ihre Adresse.

## Rang

Der Ordner teilt Rang 2 der Source Precedence mit `spec/spezifikation.md`
([`harness/README.md`](../README.md#source-precedence)): die Spezifikation für das Produkt,
dieser Ordner für die Werkzeuge. Das Lastenheft bleibt die Decke. Grundlage ist
[`MR-002`](../conventions.md#mr-002).

## Kennungen

Eine Datei je Kennung, `SPEC-<NNN>-<kurztitel>.md`; die Überschrift der Datei nennt die Kennung
und den Titel. Die Reihe der Kennungen zählt über beide Orte fortlaufend: Die nächste freie Nummer
ist die höchste vergebene in `spec/spezifikation.md` und in `harness/vertraege/`, plus eins.

## Verträge

| Kennung | Werkzeug | Datei |
|---|---|---|
| `SPEC-047` | `kopf-check` | [`SPEC-047-kopf-der-plaene.md`](SPEC-047-kopf-der-plaene.md) |
| `SPEC-048` | `abdeckung` | [`SPEC-048-abdeckung.md`](SPEC-048-abdeckung.md) |
| `SPEC-049` | `lint` | [`SPEC-049-lint-profil.md`](SPEC-049-lint-profil.md) |
| `SPEC-050` | `commit-msg` | [`SPEC-050-struktur-kennungen-im-commit.md`](SPEC-050-struktur-kennungen-im-commit.md) |
| `SPEC-052` | `gates-schnell`, `stop-hook-gegenprobe` | [`SPEC-052-gate-kette.md`](SPEC-052-gate-kette.md) |
| `SPEC-053` | `test-postgres-versionen` | [`SPEC-053-postgres-versionen.md`](SPEC-053-postgres-versionen.md) |
| `SPEC-054` | `meldungskatalog-check` | [`SPEC-054-meldungskatalog.md`](SPEC-054-meldungskatalog.md) |
| `SPEC-055` | `binaries`, `test-container` | [`SPEC-055-container-und-binaries.md`](SPEC-055-container-und-binaries.md) |
| `SPEC-056` | `test` (Schwelle der Anweisungsabdeckung), `coverage-gegenprobe` | [`SPEC-056-testabdeckung-schwelle.md`](SPEC-056-testabdeckung-schwelle.md) |
| `SPEC-057` | `test-integration` (Aufräumen), `integration-volumes-gegenprobe` | [`SPEC-057-integration-aufraeumen.md`](SPEC-057-integration-aufraeumen.md) |
| `SPEC-058` | `container-start`, `container-start-gegenprobe` | [`SPEC-058-container-start.md`](SPEC-058-container-start.md) |

## Historie

| Datum | Änderung |
|---|---|
| 2026-10-06 | Harness-Werkzeuge: Abschnitt angelegt; Prüfung des Kopfs lebender Pläne (`SPEC-047`) und Abdeckung je Anforderung und Pfad (`SPEC-048`) mit ihrem heutigen Vertrag übertragen (`LH-QA-07`, Messmethode 4) |
| 2026-10-06 | Harness-Werkzeuge: Lint-Profil mit Linter, Schwellen, Einstellungen, Verbot von `//nolint`, Export-Test-Brücke, dauerhafte Ausnahmen ohne Stufen, Ausgabe und Ausgang, Werkzeug vor dem Gate (`SPEC-049`) |
| 2026-10-06 | Lint-Profil: fehlendes Profil, Prüfung gegen das Schema, ungenutzte Regel als `lint:`-Zeile, jede Meldung ein Befund (`SPEC-049`) |
| 2026-10-06 | Lint-Profil: feste Form des Ausnahme-Abschnitts, zweites Kommentarzeichen vor `nolint`, Zeile der ungenutzten Regel, generierter Code nach Default (`SPEC-049`) |
| 2026-10-06 | Lint-Profil: Einzug der Regel-Einträge unter `rules`, Ablehnung beim Laden als Grenze (`SPEC-049`) |
| 2026-10-06 | Lint-Profil: Brücke und Test legen keinen Wert eines unexportierten Typs an; Kontext im Test bei einem Befund von `contextcheck` (`SPEC-049`) |
| 2026-10-07 | Lint-Profil: die Funktion, an die die Brücke zum Erzeugen weiterreicht, ruft auch der Produkt-Code; Eingaben exportierter Typen stellt der Test (`SPEC-049`) |
| 2026-10-07 | Lint-Profil: ungültiges YAML im Ausnahme-Abschnitt und Ergebnis aus dem Cache des Builds als Grenze (`SPEC-049`) |
| 2026-10-07 | Lint-Profil: Leerraum nach `exclusions:` und `rules:` ist Blockform, ein `-` allein ist ein Eintrag (`SPEC-049`) |
| 2026-10-08 | Lint-Profil: ob eine Regel zu den zulässigen gehört, als Grenze (`SPEC-049`) |
| 2026-10-08 | Harness-Werkzeuge: Commit-Träger lehnt Struktur-Kennungen in der Commit-Message ab; Lesebereich, Schreibweise und Wortgrenze, Merge und Revert, Vorrang vor der Annahme, Ausgabe und Ausgang (`SPEC-050`) |
| 2026-10-13 | Harness-Werkzeuge: Gate-Kette in voller und schneller Kette mit Ausnahmeliste, Stand der schnellen Kette, Wahl der Kette am geänderten Pfad (Positivliste, Dateiliste des vollen Laufs, Randformen), Hook-Teil nur zugesagt, soweit `stop-hook-gegenprobe` ihn prüft (`SPEC-052`) |
| 2026-10-13 | Harness-Werkzeuge: Schleifenschutz des Hooks als akzeptierte Grenze (Vorrang vor Konflikt, fehlendem Stand, `git`-Fehler), Umgebung des Hooks als leere Umgebung gefasst, Voraussetzung der Dateiliste (`SPEC-052`) |
| 2026-10-13 | Harness-Werkzeuge: Nachweis des vollen Laufs bei überlesenen Fehlern (`make -i gates`) und Make-Flags, Aufrufarten des Stop-Hooks, Beleg für `docs/user/` und `README.md` für jede Gegenprobe der Ausnahmeliste (`SPEC-052`) |
| 2026-10-14 | Harness-Werkzeuge: Beleg der Positivliste auch für `docs/plan/planning/` und `docs/reviews/` (Gegenprobe des Commit-Trägers ohne Bezug auf den Planungsbaum des Arbeitsbaums), `reason` bei fehlendem Stand der schnellen Kette, Pfad mit leerem Segment am Ende, Laufzeit des Hooks, Wortlaut zum älteren Nachweis (`SPEC-052`) |
| 2026-10-14 | Harness-Werkzeuge: Reihenfolge der Prüfungen des Hooks, Positivliste an einer Stelle und Selbständerung, Art `X`, Zulässigkeit der Pfade, Symlink-Einzelheiten (Spiegelregel mit Vorfahr, absolutes Ziel), Vergleich des Hashes der schnellen Kette, Freigabe in Unterverzeichnis und leerer Umgebung (`SPEC-052`) |
| 2026-10-14 | Harness-Werkzeuge: Versionsmatrix der Integrationstests mit Liste, Teilmenge, Eingabefehlern, Prüfung der Hauptversion des Servers, Ablauf, Aufräumen, Aufrufarten des Treibers und Grenzen (`SPEC-053`); Ausnahmeliste der schnellen Kette um `test-postgres-versionen` (`SPEC-052`) |
| 2026-10-17 | Harness-Werkzeuge: Abgleich von Code-Tabelle, Meldungskatalog der Spezifikation und Fehlercode-Tabellen des Handbuchs mit Quellen, Form, Mengenregeln (Rückfälle ohne Konstante zugelassen), Sperren, Ausgabe und Anschluss ohne Stufung (`SPEC-054`); Beleg des Handbuchs in der schnellen Kette (`SPEC-052`, Punkt 13) |
| 2026-10-17 | Harness-Werkzeuge: Binaries je Zielplattform mit `SHA256SUMS`, Gate über zwei Builds ohne Cache, Kopf der Binaries, Image je Plattform, Smoke im Container, Aufräumen, Aufrufarten des Treibers und Grenzen (`SPEC-055`) |
| 2026-10-17 | Harness-Werkzeuge: Aufräumen der Integrationstests mit Gegenstand (Mounts vom Typ `volume`), `rm -fv`, Lesen vor dem Entfernen, Rest und Vorrang des Codes, Ende des Laufs, Fälle und Mutationen der Gegenprobe, Docker-Version und Grenzen (`SPEC-057`); Beleg der Gegenprobe für `docs/user/` und `README.md` (`SPEC-052`, Punkt 13) |
| 2026-10-17 | Harness-Werkzeuge: Schwelle der Anweisungsabdeckung der Unit-Tests in der Stufe `test` mit Messung (nur Unit-Tests, ohne `-coverpkg`), Gesamtwert nach Anweisungen, Paketen ohne Test, Vergleich ohne Rundung, Eingaben des Prüfschritts, Ausgabe, Aufrufarten und Grenzen (`SPEC-056`); Beleg der Gegenprobe für `docs/user/` und `README.md` (`SPEC-052`, Punkt 13) |
| 2026-10-18 | Harness-Werkzeuge: Start des Dev-Containers mit Eingaben, Voraussetzungen am Host, Zustand von Verzeichnis, Schlüsselpaar und Container, `u-boot`-Aufrufen, `override.json`, Bereitschaft, Abbruch, Reihenfolge der Gründe und Exit-Codes, Aufrufarten des Treibers, Gegenprobe gegen Scratch-Container und Grenzen (`SPEC-058`) |
| 2026-10-18 | Start des Dev-Containers: `container-start-gegenprobe` ist ein Ziel für den Host, kein Teil von `GATE_CHECKS` und der Ausnahmeliste; Voraussetzung am Host ohne Überspringen (`SPEC-058`, Punkte 15 und 16) |
