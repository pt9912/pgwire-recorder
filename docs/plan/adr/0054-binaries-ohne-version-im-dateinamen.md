# ADR-0054: Die Release-Binaries tragen keine Version im Dateinamen, das Image behält den Tag mit Version

**Status:** Proposed

**Datum:** 2026-10-20

**Autor:** Architect. Entschieden hat der Nutzer am 2026-10-20 (so vom Koordinator übermittelt): Die Release-Binaries tragen keine Version im Dateinamen, weil die Version im Namen dem Benutzer nur mehr Arbeit macht; das Image behält den Tag mit Version. Das Namensschema unten ist die Empfehlung des Architect innerhalb dieser Entscheidung; seine Annahme und damit die Annahme dieser ADR steht aus.

**Bezug:** [`LH-FA-01`](../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [ADR-0043](0043-binaries-und-container-gate.md), [ADR-0051](0051-veroeffentlichung-auf-github-actions.md), [ADR-0053](0053-homebrew-tap-in-der-release-pipeline.md)

**Schärft:** [`SPEC-055`](../../../harness/vertraege/SPEC-055-container-und-binaries.md), [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`make binaries` schreibt `pgwire-recorder_<VERSION>_<os>_<arch>` (unter Windows mit `.exe`); `v0.1.0` ist so veröffentlicht. Wer ein Binary lädt, muss den Namen mit der Version bilden, in Skripten und in der Homebrew-Formel die Version in den Namen einsetzen und sie nach dem Laden meist wieder entfernen, um das Werkzeug unter einem festen Namen zu starten. Die Version steht ohnehin im Release (Tag, Pfad der Adresse) und im Binary (`version`). Das Schwesterprojekt u-boot nennt seine Binaries `u-boot-<os>-<arch>` und ist das Vorbild des Nutzers.

Das Lastenheft legt kein Namensschema fest (`LH-FA-01`, `LH-QA-03` nennen Plattformen, keine Namen); das Schema steht allein im Vertrag der Werkzeuge. Die Formel für Homebrew ([ADR-0053](0053-homebrew-tap-in-der-release-pipeline.md)) ist noch nicht gebaut und hängt am Namen.

## Entscheidung

Wir wählen **Option A: `pgwire-recorder-<os>-<arch>`, unter Windows `pgwire-recorder-windows-<arch>.exe`, ohne Version; das Image behält den Tag `<Version>` (und `latest` bei einem stabilen Tag).**

1. **Schema.** Sechs Binaries (`darwin`, `linux`, `windows`, je `amd64` und `arm64`) plus `SHA256SUMS`, am Release also weiter genau sieben Anhänge ([`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md) Punkt 9 bleibt in der Zahl).
2. **Die Version steht im Binary und im Release, nicht im Namen.** `VERSION` geht weiter in `-X main.version` des Builds; `version` gibt sie aus. Der Name unterscheidet zwei Versionen nicht mehr; ein Verzeichnis `DIST` wird deshalb wie bisher nur leer oder fehlend angenommen.
3. **`v0.1.0` bleibt unverändert** mit den Namen `pgwire-recorder_0.1.0_<os>_<arch>`. Es bekommt keine Formel: Der Nachtrag von Hand aus [ADR-0053](0053-homebrew-tap-in-der-release-pipeline.md) Punkt 6 bleibt als Verfahren für einen Tag mit dem neuen Schema, dessen Lauf den Job nicht hatte oder dessen Formel die Vorlage nicht mehr trägt; das Beispiel `v0.1.0` in jenem Punkt entfällt. Die erste Formel entsteht mit dem nächsten stabilen Tag (neuer Tag nur mit neuem Ja des Nutzers).

Die Regeln im Einzelnen (Schema, Dateiliste, `SHA256SUMS`, Fälle der Gegenprobe) stehen in [`SPEC-055`](../../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 2 und 3 und [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md) Punkt 4 und 19.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — `pgwire-recorder-<os>-<arch>`** | feste Adresse je Plattform, ein Name in Skripten und Formel; wie u-boot | zwei Versionen im selben Verzeichnis überschreiben sich (Absicht); aus dem Namen ist die Version nicht lesbar |
| B — Schema bleibt, Version im Namen | nichts zu ändern; Name nennt die Version | die Arbeit des Benutzers, die der Nutzer vermeiden will |
| C — `pgwire-recorder_<os>_<arch>` (Unterstrich) | nah am bisherigen Schema | weicht vom Vorbild u-boot ohne Grund ab; der Name `pgwire-recorder` trägt selbst einen Bindestrich |
| D — Binaries in Archiven (`tar.gz`, `zip`) mit festem Namen darin | übliche Form, Version im Archivnamen | ein Archiv braucht Entpacken, eine zweite Prüfung der Reproduzierbarkeit (Archiv-Metadaten) und ändert Gate und Formel stärker; der Nutzer wollte weniger Arbeit |

## Konsequenzen

- Positiv: Die Adresse eines Binaries am Release hat die Form `…/download/v<Version>/pgwire-recorder-<os>-<arch>`; die Formel braucht die Version nur im Pfad und in der Zeile `version`.
- Negativ: `v0.1.0` und spätere Releases haben verschiedene Namen; das Handbuch nennt für `v0.1.0` die alten, bis der nächste Tag es ablöst. Pipeline und Gate binden die Version nicht mehr über den Namen an das Verzeichnis; sie bindet der Aufruf `make binaries VERSION=<Version>` und, für `linux/amd64`, die Kette Summe im Image (`SPEC-067` Punkt 5) und Smoke mit `version` (Punkt 7).
- Akzeptiertes Negativ: Das Werkzeug erkennt nicht, dass ein `DIST` mit Dateien einer anderen Version nur nach Name aussähe wie ein gültiger Lauf; es lehnt jedes nicht leere `DIST` ab. Grund: dieselbe Regel wie bisher, sie braucht keine Lesart des Namens.
- Akzeptiertes Negativ: Ein Name aus Kleinbuchstaben, Ziffern, `-` und `.` ist auf Linux, macOS und Windows gültig; erzeugt und geprüft ist er nur unter Linux (Docker). Grund: die Namen sind ASCII ohne Sonderzeichen, die Grenze ist die der Plattformen von `LH-QA-03` für alle Binaries (`SPEC-055` Punkt 13 (4)).
- Folgepflicht: [`SPEC-055`](../../../harness/vertraege/SPEC-055-container-und-binaries.md), [`SPEC-067`](../../../harness/vertraege/SPEC-067-release.md); `tools/harness/binaries.sh`, `run-container.sh`, `release.sh`, `container-gegenprobe.sh`, `release-gegenprobe.sh`; `harness/README.md` §Sensors (`make binaries`); Benutzerhandbuch (Abschnitt zu `make binaries`), `docs/maintainer/releasing.md`. Lastenheft unverändert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, make, Attrappe von `docker` | die Namen der Binaries und der Zeilen von `SHA256SUMS`, `DIST` nicht leer | `make test-container-gegenprobe` |
| bash, git, Attrappen | Dateiliste des Release, Namen in der Formel | `make release-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn ein Anwender Binaries mehrerer Versionen nebeneinander braucht und die Überschreibung zum Fehler führt; wenn ein Paketmanager (etwa Scoop) einen Namen mit Version verlangt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-20 | Proposed (Empfehlung des Architect innerhalb der Entscheidung des Nutzers) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
