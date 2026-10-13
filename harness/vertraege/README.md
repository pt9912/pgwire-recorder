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
