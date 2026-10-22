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
ist die höchste vergebene in `spec/spezifikation.md` und in `harness/vertraege/`, plus eins. Ausgenommen ist
eine Nummer, die ein geplanter Vertrag freihält: `SPEC-059` (Kopie des Arbeitsbaums in den Gegenproben) ist noch nicht vergeben.
Frei sind außerdem die Lücken `SPEC-065` und `SPEC-066`; drei Pläne in `open/` nennen sie (CI-Nachtläufe, Freshness).

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
| `SPEC-060` | `gegenproben`, `gegenproben-host`, `gegenproben-gegenprobe`; Namensregel der Gegenproben | [`SPEC-060-gegenproben.md`](SPEC-060-gegenproben.md) |
| `SPEC-061` | `gegenproben-nach-einlieferung` (Hook im Bare-Repo des Gegenproben-Containers), `gegenproben-nach-einlieferung-gegenprobe` | [`SPEC-061-gegenproben-nach-einlieferung.md`](SPEC-061-gegenproben-nach-einlieferung.md) |
| `SPEC-062` | `nachweis-uebernehmen`, `nachweis-uebernehmen-gegenprobe` | [`SPEC-062-nachweis-uebernehmen.md`](SPEC-062-nachweis-uebernehmen.md) |
| `SPEC-063` | `workflows-check`, `workflows-check-gegenprobe` | [`SPEC-063-workflows-check.md`](SPEC-063-workflows-check.md) |
| `SPEC-064` | `ci-range`, `ci-gegenprobe`; Text von `ci.yml` | [`SPEC-064-ci.md`](SPEC-064-ci.md) |
| `SPEC-067` | `release.yml`, `release.sh`, `release-gegenprobe`; Datei `VERSION` | [`SPEC-067-release.md`](SPEC-067-release.md) |
| `SPEC-068` | `archiv-bezug`, `archiv-bezug-gegenprobe` | [`SPEC-068-archiv-bezug.md`](SPEC-068-archiv-bezug.md) |
| `SPEC-069` | `container-einrichten`, `container-einrichten-gegenprobe` | [`SPEC-069-container-einrichten.md`](SPEC-069-container-einrichten.md) |
| `SPEC-070` | `pretooluse-zeitgrenze` (Hook), `zeitgrenze-gegenprobe`; Zeitgrenze eines Rollenlaufs | [`SPEC-070-rollenlauf-zeitgrenze.md`](SPEC-070-rollenlauf-zeitgrenze.md) |

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
| 2026-10-18 | Gegenproben außerhalb der Gate-Kette: zwei Listen (Container, Host), Namens- und Eindeutigkeitsregel, zwei Sammelziele mit Reihenfolge, rotem Einzelziel, Ausgabe, Abbruch, leerer Liste, Aufrufarten von `make` und Reihenfolge der Gründe, Pflicht zum Fahren, Grenzen und Plattform (`SPEC-060`); Ausnahmeliste der schnellen Kette auf drei Gates, Nachweis ohne Gegenproben, Beleg der Positivliste als Regel an den Schreiber einer Gegenprobe (`SPEC-052`); Anschluss der Gegenproben an die Liste statt an `GATE_CHECKS` (`SPEC-053` Punkte 12 und 13, `SPEC-054` Punkt 11, `SPEC-055` Punkt 12, `SPEC-056` Punkt 11, `SPEC-057` Punkt 10, `SPEC-058` Punkt 16) |
| 2026-10-18 | Hook für die Einlieferung in den Gegenproben-Container: Marke, Refs, Löschung und Form der Zeilen, Hintergrundlauf, Sperre, derselbe Hash, Kodierung des Zweignamens, Ergebnis und Marke des laufenden Laufs, nur `make gegenproben`, Umgebung, Arbeitsbaum je Lauf, Reihenfolge der Gründe, Lesen des Ergebnisses und Grenzen (`SPEC-061`) |
| 2026-10-18 | Übernahme des Nachweises des vollen Laufs aus einer anderen Umgebung: Quelle als Verzeichnis, Bedingungen (unbeschädigte Liste, Hash neben der Liste, byte-gleiche Liste mit dem Baum des Ziels), Reihenfolge der Gründe, Kennwörter und Exit-Codes, Ersetzen und Ablehnen ohne Änderung, Signale, Make-Flags, Kopplung an den Hook, Grenzen und Plattform (`SPEC-062`); zweiter Schreiber des Nachweises (`SPEC-052` Punkte 1, 8 und 21) |
| 2026-10-18 | Gate der Workflows unter `.github/workflows/`: Prüfbereich, `actionlint` im Image per Digest, Einstufung, Lesegrammatik, Rechte, Zeitgrenzen, Bindung der Pin-Prüfung von d-check, Ausgabe und Reihenfolge der Gründe, Aufrufarten, Grenzen (`SPEC-063`; Empfehlung bis zur Annahme von [ADR-0050](../../docs/plan/adr/0050-ci-auf-github-actions.md)) |
| 2026-10-18 | Die CI auf GitHub: wörtlicher Text von `ci.yml`, Berechnung der Range aus dem Ereignis (`ci-range.sh`), Konfiguration der Range-Prüfungen von d-check (`commits`, `vcs`), Aufbau der Gegenprobe, Grenzen (`SPEC-064`; Empfehlung bis zur Annahme von [ADR-0050](../../docs/plan/adr/0050-ci-auf-github-actions.md)) |
| 2026-10-19 | Veröffentlichung bei einem Tag: wörtlicher Text von `release.yml`, Trockenlauf statt Probe-Tag, Eingaben (SemVer, `VERSION`, Notes, Geheimnisse), Binaries, Image als Manifestliste mit Prüfung der Summe je Plattform, Push-Reihenfolge ohne Überschreiben, Vergleich der Index-Digests aus den rohen Bytes, `latest`, GitHub-Release, Zugangsdaten, Fehlertabelle, Grenzen (`SPEC-067`; [ADR-0051](../../docs/plan/adr/0051-veroeffentlichung-auf-github-actions.md)); Smoke als Skript mit Image-Referenz (`SPEC-055` Punkt 7) |
| 2026-10-19 | Anforderungen im Stub archivierter Slices: Lauf-Menge, Eingaben und Reihenfolge der Gründe, Prüfung aller Stubs vor dem ersten Schreiben, Kopf des Volltexts und gelesene Kennungen, Zeile und Ordnung, vorhandene Zeile und Idempotenz, Stelle der Zeile, Zeilenende, Aufrufarten von `make`, kein Modus „prüfen“, Grenzen (`SPEC-068`) |
| 2026-10-19 | Veröffentlichung: Scan des Images vor dem Push mit Trivy (Pin, Aufruf, Auswertung, Fail-closed, Marke `scan.ok`, keine Ausnahmen, [ADR-0052](../../docs/plan/adr/0052-scan-des-images-vor-dem-push.md)), Freshness-Audit und Freigabe-Checkliste (Form des Belegs, Anti-Punkte, Incident-Klausel), der erste Tag (Belege, „rot“ nach dem Tag, Bezug im Handbuch) (`SPEC-067` Punkte 16 bis 18) |
| 2026-10-21 | Harness-Werkzeuge: Einrichtung im laufenden Dev-Container mit Eingaben, Werkzeugen und Projektname aus `u-boot.yaml`, einem Lesezug für den Zustand, Bare-Repo, Remote und Einlieferung, Klon, Identität des Klons nur `--local` über den Standardeingang, Commit-Träger, Zweig aus `origin/main`, Träger als letzter Schritt, Reihenfolge der Gründe und Exit-Codes, Abbruch, Aufrufarten des Treibers, Gegenprobe mit Attrappen ohne Container (Zeitbudget nach `SPEC-060` Punkt 16), Nachweis gegen einen laufenden Container als Handlung, Plattform und Grenzen (`SPEC-069`; [ADR-0057](../../docs/plan/adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md), Accepted 2026-10-21) |
| 2026-10-19 | Veröffentlichung: Das Freshness-Audit ist optional und kein Punkt der Freigabe-Checkliste (elf statt zwölf Punkte); die Tabelle als Form seines Belegs gilt nur, wenn es durchgeführt wird; ein Ergebnis sperrt den Tag nie (Entscheidung des Nutzers vom 2026-10-19; `SPEC-067` Punkte 17 und 18) |
| 2026-10-21 | Gegenproben: der Start mit ignoriertem `SIGINT` ist nicht zugesagt, Regel an den Aufrufer mit Rezept für den Hintergrundstart, Messung und Plattform (`SPEC-060` Punkt 17, Verweise in `SPEC-060` Punkt 8 und `SPEC-069` Punkte 17 und 19) |
| 2026-10-22 | Harness-Werkzeuge: Zeitgrenze eines Rollenlaufs mit Wächter im Werkzeug-Aufruf (Eintrag, Wer darunterfällt, Anker und Zeit, Reihenfolge der Gründe, Sperrmenge und Durchgelassenes, Ausgabe, Fail-open mit Selbstauskunft, Umgebung und Aufrufe, kein Weg über die Grenze, pflichtige Läufe, Nachweis im Agentenlauf, Grenzen, Gegenprobe und Aufnahme) (`SPEC-070`; Empfehlung bis zur Annahme von [ADR-0058](../../docs/plan/adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md)) |
| 2026-10-22 | Nachweise: Der Beleg liegt als eigene, ungekürzte Datei unter `docs/reviews/`, §7 des Plans nennt Pfad, Summe, Zeile mit dem Exit-Code und Kopf; Kopf in Prosa, Ausgabe im Codeblock (Entscheidung des Nutzers vom 2026-10-22; `SPEC-067` Punkt 20 (6) und (7), `SPEC-069` Punkt 18) |
| 2026-10-22 | Zeitgrenze eines Rollenlaufs: [ADR-0058](../../docs/plan/adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) angenommen (Wortlaut des Nutzers „adr accept“); Geltungsbereich bestätigt („2. ja“): Rollenläufe, Commands und Hauptsitzung nicht erzwungen; Zusagen über das Verhalten von Claude Code enger gefasst (`SPEC-070`) |
