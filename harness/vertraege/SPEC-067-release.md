# SPEC-067 — Veröffentlichung: Binaries und Image bei einem Tag (`release.yml`, `release.sh`, `release-gegenprobe`)

Ein Tag `v<SemVer>` löst auf GitHub einen Lauf aus, der Binaries und Image prüft, baut und veröffentlicht: die sechs Binaries mit
`SHA256SUMS` aus `make binaries` am GitHub-Release, das Image als **eine** Manifestliste für `linux/amd64` und `linux/arm64` mit demselben
Index-Digest in `ghcr.io/pt9912/pgwire-recorder` und `docker.io/pt9912/pgwire-recorder`. Der Lauf ist **fail-closed**: jede Eingabe, die nicht stimmt,
endet vor dem ersten Schritt mit Außenwirkung; ein Fehler danach hinterlässt einen benannten Zustand und lässt sich ohne Veränderung eines Tags
wieder anstoßen (Punkt 11). Dieselbe Datei läuft auch als **Trockenlauf** auf einem Zweig und tut dann alles außer dem, was nach außen wirkt (Punkt 2); damit ist
der erste Tag der erste Lauf mit Außenwirkung, ohne einen Probe-Tag in den Registries zu hinterlassen. Warum die Veröffentlichung auf GitHub Actions läuft, wie
die Registries und Tags heißen und warum kein Probe-Tag nötig ist, steht in [ADR-0051](../../docs/plan/adr/0051-veroeffentlichung-auf-github-actions.md);
diese Kennung wiederholt es nicht. Das Gate der Datei ist [`SPEC-063`](SPEC-063-workflows-check.md), die Konventionen der CI sind [`SPEC-064`](SPEC-064-ci.md).

**Stand der Entscheidung.** Der Nutzer hat [ADR-0051](../../docs/plan/adr/0051-veroeffentlichung-auf-github-actions.md) am 2026-10-19 angenommen und dabei entschieden: **kein Probe-Tag** (Trockenlauf statt dessen), Repository-Geheimnisse statt eines Environments, Tag-Satz nur `<Version>` und `latest`. Rahmen des Nutzers vom 2026-10-18: der erste Tag `v0.1.0` und das Veröffentlichen sind erlaubt; weitere Tags nur mit neuem Ja; Zugangsdaten kommen nur als Geheimnisse von GitHub und nur vom Nutzer. Die übrigen Festlegungen dieser Kennung sind Entscheidung des Architect, vom Nutzer mit der ADR bestätigt.

**Zwei Slices** (Schnitt vom 2026-10-19, Entscheidung des Architect nach der Größe). **Slice A** `slice-erster-release-veroeffentlichung` (Tag und Binaries) liefert die Eingaben ohne Geheimnisse, die Binaries, den GitHub-Release, die Datei `VERSION` und den Trockenlauf; seine Fassung von `release.yml` ist Punkt 1a und setzt nichts von B voraus. **Slice B** `slice-erster-release-veroeffentlichung-image` (Image) liefert Bau, Veröffentlichung, Vergleich, Smoke und `latest`; seine Fassung ist die Endfassung Punkt 1b, die Fassung A um Zeilen ergänzt, ohne eine Zeile von A zu ändern oder zu streichen, und hängt die Zeilen an der in 1b gezeigten Stelle ein. Welcher Punkt zu welchem Slice gehört, zeigt die Tabelle nach der Randform-Liste; in jedem Punkt steht die Marke **(A)**, **(B)** oder **(A, ergänzt von B)**.

**Schnittstelle.** `.github/workflows/release.yml` hat den Text von Punkt 1. Die Logik steht in `tools/harness/release.sh <Schritt>` (Schritte:
`eingaben`, `binaries`, `image`, `veroeffentlichen`, `vergleichen`, `smoke`, `latest`, `github-release`), nicht im YAML, wo sie ungeprüft bliebe. Die Version steht in
der Datei `VERSION` im Wurzelverzeichnis: genau eine Zeile, `<SemVer>` ohne `v`, LF, keine Leerzeichen. Der Smoke ist das Skript `tools/harness/container-smoke.sh <Image-Referenz>`, das
der Smoke von [`SPEC-055`](SPEC-055-container-und-binaries.md) Punkt 7 und dieser Lauf gemeinsam rufen (SPEC-055 behält seinen Text; die Quelle der Image-Referenz ist ein Parameter). Das Fragment
`harness/mk/release.mk` definiert `release-gegenprobe` (`.PHONY`), das `tools/harness/release-gegenprobe.sh` ruft und in `GEGENPROBEN_CONTAINER` steht
([`SPEC-060`](SPEC-060-gegenproben.md) Punkte 1 bis 3). Die Pipeline läuft nie unter `make`.

**Plattform** (`AGENTS.md` §3.12). Der Lauf nutzt `ubuntu-24.04`, `linux/amd64`, mit den vorinstallierten Werkzeugen `docker` (mit `buildx`), `skopeo`, `gh`, `git`, `bash`
und GNU-Werkzeugen. Gestartet und geprüft wird das Image für `linux/amd64`; das für `linux/arm64` ist nach Struktur, Architektur und Summe des Binaries belegt, nicht ausgeführt
(kein Emulator vorausgesetzt; wie [`SPEC-055`](SPEC-055-container-und-binaries.md) Punkt 13 (3)). Binaries für `darwin` und `windows` sind gebaut und gesummt, nicht ausgeführt. Andere Runner sind nicht
zugesagt. Die Gegenprobe läuft auf Linux mit bash ≥ 5.1, `git`, `make` und Attrappen, ohne Docker und ohne Netz.

Der Name jedes Falls der Gegenprobe trägt die Nummer des Punktes. Je Zusage und je Element einer Aufzählung gibt es einen Fall und eine Mutation, die nur ihn färbt
(`AGENTS.md` §3.10); die unveränderte Kopie ist grün. Wo ein Fall mit „Attrappe“ steht, zeichnet die Attrappe von `docker`, `skopeo`, `gh` ihre Aufrufe auf und liefert Ausgabe und
Exit-Code des Falls.

**Randform → Punkt** (§6 des Plans):

| Randform | Punkt |
|---|---|
| Träger der Automatisierung, Pins | 1, [ADR-0051](../../docs/plan/adr/0051-veroeffentlichung-auf-github-actions.md) |
| Reihenfolge der Schritte, erste Außenwirkung | 2, 6 |
| Manifestliste und Image-Bau, Binary im Image | 5 |
| Registry-Namen, Tag-Satz, Sichtbarkeit | 6, 7 |
| Zugangsdaten | 10 |
| Anhänge am GitHub-Release, Notes | 9 |
| `:latest` | 8 |
| Fail-closed-Digest-Vergleich | 7 |
| Versionsdatei | Schnittstelle, 3 |
| SemVer-Prüfung | 3 |
| Wiederanlauf und Fehlertabelle | 11 |
| Probe-Tag, Trockenlauf | 2, 12 |
| Beispiel und `:latest` in `LH-FA-16.a` | 8 (Zusage), Spezifikation `LH-FA-16.a` |
| Plattform | Kopf dieser Kennung, 15 |
| Aufrufarten von `make` | 13 |
| Abbruch und Herunterfahren | 13 |
| Reihenfolge mehrerer Gründe | 3, 11 |
| Stelle des Scans vor dem Push | 5 |

**Zuordnung zu den Slices.**

| Punkt | Slice | Anteil |
|---|---|---|
| 1a | A | Text von `release.yml` ohne Image-Schritte |
| 1b | B | Endfassung (Text 1a plus die Zeilen von B); die Fälle `p1-*` halten die Zeilenliste, B hängt seine Zeilen an |
| 2 | A, ergänzt von B | Trockenlauf, Auslöser; A: Außenwirkung nur `GitHub-Release`; B: die weiteren vier Schritte (Fälle je Schritt) |
| 3 (1) bis (6) | A | Eingaben ohne Geheimnisse |
| 3 (7) | B | Geheimnisse; Fälle `p3-secret-*` |
| 4 | A | Binaries |
| 5 | B | Image und Stelle des Scans |
| 6, 7, 8 | B | Veröffentlichen, Vergleichen und Smoke, `latest` |
| 9 | A, ergänzt von B | GitHub-Release; B ergänzt die Vorbedingung `smoke` (Marke `smoke.ok`, Fall `p9-nicht-vor-registries`) |
| 10 | A, ergänzt von B | A: `GITHUB_TOKEN` und das Verbot von Ausgabe und `set -x`; B: Geheimnisse der Registries und Anmeldung |
| 11 | A, ergänzt von B | Fehlertabelle: A die Zeilen `eingaben`, `binaries`, `github-release`; B die übrigen |
| 12 | A, ergänzt von B | Trockenlauf: A belegt Eingaben, Binaries, Dauer; B ergänzt Bau und Prüfung des Images |
| 13, 14, 15 | A, ergänzt von B | Aufrufarten, Abbruch, Aufbau der Gegenprobe, Grenzen |

1. **Der Text (A: 1a; B: 1b).** `.github/workflows/release.yml` besteht, außer leeren Zeilen und Kommentarzeilen (erstes Zeichen außer Leerzeichen `#`), aus genau den Zeilen der Fassung seines Slice in dieser Reihenfolge
   (Einzug und Anführungszeichen sind Teil des Textes; für `<SHA>` und `<Tag>` gilt dasselbe wie in [`SPEC-064`](SPEC-064-ci.md) Punkt 1).

   **1a: Fassung von Slice A** (Tag und Binaries; kein Geheimnis außer `GITHUB_TOKEN`, kein `packages: write`):

   ```
   name: release
   on:
     push:
       tags:
         - 'v*'
     workflow_dispatch:
   permissions: {}
   concurrency:
     group: release-${{ github.event_name }}
     cancel-in-progress: false
   jobs:
     release:
       name: release
       runs-on: ubuntu-24.04
       timeout-minutes: 60
       permissions:
         contents: write
       env:
         RELEASE_EVENT: ${{ github.event_name }}
         RELEASE_REF: ${{ github.ref }}
         RELEASE_SHA: ${{ github.sha }}
         RELEASE_DEFAULT_BRANCH: ${{ github.event.repository.default_branch }}
         RELEASE_REPOSITORY: ${{ github.repository }}
       steps:
         - name: Checkout
           uses: actions/checkout@<SHA> # <Tag>
           with:
             fetch-depth: 0
             persist-credentials: false
         - name: Eingaben
           run: bash tools/harness/release.sh eingaben
         - name: Binaries
           run: bash tools/harness/release.sh binaries
         - name: GitHub-Release
           if: ${{ github.event_name == 'push' }}
           env:
             GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
           run: bash tools/harness/release.sh github-release
   ```

   **1b: Endfassung von Slice B** (Fassung 1a plus die Zeilen von B; in 1a stehende Zeilen bleiben unverändert und in ihrer Reihenfolge, B fügt hinzu: die Zeile `packages: write` hinter `contents: write`, den Block `env:` mit den zwei Geheimnissen im Schritt `Eingaben`, und die fünf Schritte `Image`, `Veröffentlichen`, `Vergleichen`, `Smoke` und `Latest` zwischen `Binaries` und `GitHub-Release`):

   ```
   name: release
   on:
     push:
       tags:
         - 'v*'
     workflow_dispatch:
   permissions: {}
   concurrency:
     group: release-${{ github.event_name }}
     cancel-in-progress: false
   jobs:
     release:
       name: release
       runs-on: ubuntu-24.04
       timeout-minutes: 60
       permissions:
         contents: write
         packages: write
       env:
         RELEASE_EVENT: ${{ github.event_name }}
         RELEASE_REF: ${{ github.ref }}
         RELEASE_SHA: ${{ github.sha }}
         RELEASE_DEFAULT_BRANCH: ${{ github.event.repository.default_branch }}
         RELEASE_REPOSITORY: ${{ github.repository }}
       steps:
         - name: Checkout
           uses: actions/checkout@<SHA> # <Tag>
           with:
             fetch-depth: 0
             persist-credentials: false
         - name: Eingaben
           env:
             RELEASE_DOCKERHUB_USER: ${{ secrets.DOCKERHUB_USERNAME }}
             RELEASE_DOCKERHUB_TOKEN: ${{ secrets.DOCKERHUB_TOKEN }}
           run: bash tools/harness/release.sh eingaben
         - name: Binaries
           run: bash tools/harness/release.sh binaries
         - name: Image
           run: bash tools/harness/release.sh image
         - name: Veröffentlichen
           if: ${{ github.event_name == 'push' }}
           env:
             RELEASE_GHCR_USER: ${{ github.actor }}
             RELEASE_GHCR_TOKEN: ${{ secrets.GITHUB_TOKEN }}
             RELEASE_DOCKERHUB_USER: ${{ secrets.DOCKERHUB_USERNAME }}
             RELEASE_DOCKERHUB_TOKEN: ${{ secrets.DOCKERHUB_TOKEN }}
           run: bash tools/harness/release.sh veroeffentlichen
         - name: Vergleichen
           if: ${{ github.event_name == 'push' }}
           run: bash tools/harness/release.sh vergleichen
         - name: Smoke
           if: ${{ github.event_name == 'push' }}
           run: bash tools/harness/release.sh smoke
         - name: Latest
           if: ${{ github.event_name == 'push' }}
           env:
             RELEASE_GHCR_USER: ${{ github.actor }}
             RELEASE_GHCR_TOKEN: ${{ secrets.GITHUB_TOKEN }}
             RELEASE_DOCKERHUB_USER: ${{ secrets.DOCKERHUB_USERNAME }}
             RELEASE_DOCKERHUB_TOKEN: ${{ secrets.DOCKERHUB_TOKEN }}
           run: bash tools/harness/release.sh latest
         - name: GitHub-Release
           if: ${{ github.event_name == 'push' }}
           env:
             GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
           run: bash tools/harness/release.sh github-release
   ```

   Die Gegenprobe hält die Zeilenliste mit der Herkunft A oder B je Zeile; B hängt seine Zeilen an die Liste an und lässt die Fälle von A unverändert (die Fälle `p1-fehlt-<n>` und `p1-geaendert-<n>` entstehen in einer Schleife über die Liste). Jede Änderung der Datei ist eine Änderung dieser Kennung; der Scan-Schritt von `slice-erster-release-freigabe` kommt zwischen `Image` und `Veröffentlichen` (Punkt 5) und ändert diesen Text.
   Ein Schritt mit Außenwirkung (`Veröffentlichen`, `Vergleichen`, `Smoke`, `Latest`, `GitHub-Release`) trägt die Bedingung auf das Ereignis `push`; die Geheimnisse stehen nur in den
   Schritten, die sie lesen. Fälle und Mutationen wie in [`SPEC-064`](SPEC-064-ci.md) Punkt 1: `p1-gleich`, je Zeile `p1-fehlt-<n>` und `p1-geaendert-<n>`, Zusatzzeilen
   (`pull_request_target`, `secrets.ACTIONS_RUNNER_DEBUG`, ein weiterer Schritt mit `uses:` außer dem Checkout, `environment:`, `continue-on-error: true`), die Pin-Fälle `p1-pin-*`, die Reihenfolge der Schritte, und je Schritt mit
   Außenwirkung ein Fall `p1-bedingung-<Schritt>` (Bedingung entfernt: rot) mit Mutation.
2. **Auslöser, Trockenlauf, Außenwirkung (A, ergänzt von B).** Ein Push eines Tags `v*` löst den **Lauf** aus; `workflow_dispatch` (nur im Standard-Zweig oder in einem Zweig, den der Aufrufer wählt) löst den **Trockenlauf** aus.
   `ci.yml` hat `tags-ignore`, `release.yml` nimmt nur `v*`: ein Tag läuft nur hier. Der Trockenlauf fährt die Schritte `Eingaben`, `Binaries` und `Image` (und später `Scan`) mit der Version aus `VERSION`
   als Tag `v<VERSION>`, schreibt nichts in eine Registry, legt keinen GitHub-Release an, setzt kein `latest` und nimmt keine Anmeldung vor. **Außenwirkung** haben genau die Schritte `Veröffentlichen`, `Vergleichen`, `Smoke`, `Latest` und `GitHub-Release`; in Fassung 1a hat nur `GitHub-Release` Außenwirkung, die vier anderen kommen mit Slice B;
   `release.sh` lehnt sie zusätzlich selbst ab, wenn `RELEASE_EVENT` nicht `push` ist oder `RELEASE_REF` nicht `refs/tags/v…` (Exit 2, vor jedem Aufruf von `docker`, `skopeo` und `gh`). Die erste Außenwirkung ist der Push nach `ghcr.io`. Fälle:
   `p2-trockenlauf-ohne-aussenwirkung` (Ereignis `workflow_dispatch`: alle Schritte bis `image` laufen, die Attrappen von `skopeo` und `gh` werden nie gerufen, kein `login`), `p2-aussenwirkung-ohne-push` (je Schritt mit Außenwirkung bei `workflow_dispatch`: Exit 2), `p2-ref-kein-tag` (Ereignis `push`, Ref
   `refs/heads/main`: Exit 2). Mutationen: der Trockenlauf ruft `skopeo`; ein Schritt mit Außenwirkung prüft das Ereignis nicht (je Schritt).
3. **Eingaben (Schritt `eingaben`; (1) bis (6) A, (7) B).** Vor jedem Schritt mit Außenwirkung; Exit 2 und eine Meldung, die den Grund nennt, beim ersten Treffer in dieser Reihenfolge: (1) Ereignis und Ref (Punkt 2); (2) die **Form** des Tags
   `v<SemVer>` nach SemVer 2.0: `v`, dann `MAJOR.MINOR.PATCH` aus Zahlen ohne führende Null (außer `0`), danach optional `-` und Vorabversions-Bezeichner nach SemVer 2.0 (durch `.` getrennt, je Bezeichner nicht leer: entweder rein numerisch ohne führende Null (außer `0`) oder aus `0-9A-Za-z-` mit mindestens einem Zeichen, das keine Ziffer ist; `-01` ist abgelehnt, `-0a` zugelassen); Build-Metadaten (`+…`)
   sind nicht zugelassen (ein Docker-Tag trägt kein `+`); Länge bis 100 Zeichen; (3) die Datei `VERSION` hat die Form (eine Zeile, LF, `<SemVer>` ohne `v` nach derselben Regel); (4) der Tag ist `v` plus der Inhalt von `VERSION`, Zeichen für Zeichen; (5) die Änderungshistorie des
   Benutzerhandbuchs trägt die Zeile der Version (Punkt 9) mit nicht leerer dritter Zelle (eine leere Zelle gilt wie eine fehlende Zeile; die Tabellenzeile hat genau drei Zellen: ein weiteres `|` in der Zelle, auch maskiert als `\|`, macht die Zeile für die Prüfung ungültig und ist ein Eingabefehler `notes`, damit kein abgeschnittener Text in den Release geht), sonst Eingabefehler; **im Trockenlauf** ist eine fehlende Zeile ein Hinweis (die Zeile schreibt der Slice der Freigabe mit dem Tag, die Handbuchfassung ist bis dahin Ist-Zustand und trägt keine Zeile für eine Version, die es nicht gibt), eine vorhandene Zeile mit leerer dritter Zelle bleibt ein Fehler, und eine vorhandene Zeile wird ausgewertet und ihr Text als Hinweis ausgegeben; die Freigabe fährt den Trockenlauf mit Zeile und `VERSION` des Tags einmal vor dem Tag, ohne Hinweis zur Zeile; (6) nur beim Lauf: der Commit des Tags liegt im Standard-Zweig (`git merge-base --is-ancestor $RELEASE_SHA origin/<Standard-Zweig>`); im Trockenlauf steht die Antwort als Hinweis und ist kein Fehler; (7) **(B)** die Geheimnisse `DOCKERHUB_USERNAME`
   und `DOCKERHUB_TOKEN` sind gesetzt und nicht leer (Trockenlauf und Lauf; die Meldung nennt den Namen, nie den Wert); fehlt eines, ist es ein Eingabefehler **vor dem ersten Push**. Der Trockenlauf ersetzt (1) durch den Trockenlauf, leitet den Tag als `v` plus die erste Zeile von `VERSION` ab und lässt (4), (6) als Hinweis, (5) wie beschrieben und (7) wirksam. **Meldungen:** Jede Meldung eines Eingabefehlers beginnt mit `release: ` und einem Schlüsselwort, das den Grund nennt: `ereignis`, `tag-form`, `version-datei`, `tag-ungleich`, `notes`, `standard-zweig` und mit Slice B `geheimnis` (mit dem Namen); ein Hinweis beginnt mit `release: hinweis: `. Der übrige Wortlaut ist nicht zugesagt. Ein unbekannter Schritt und ein fehlendes `RUNNER_TEMP` enden mit Exit 2 (nicht als Fall zugesagt; akzeptiertes Negativ, Grund: Fehlbedienung des Aufrufers, den die Pipeline nie so ruft); ein fehlendes `GH_TOKEN` im Schritt `github-release` endet mit Exit 2 vor dem ersten Aufruf von `gh` (Punkt 9). Fälle: `p3-tag-gueltig` (`v1.2.3`, `v0.1.0-rc.1`, `v1.0.0-x.7.z.92`),
   `p3-ohne-v`, `p3-zwei-teile` (`v1.2`), `p3-fuehrende-null` (`v01.2.3`), `p3-leere-vorabversion` (`v1.2.3-`), `p3-build-metadaten` (`v1.2.3+x`), `p3-zu-lang`, `p3-version-datei-form` (zwei Zeilen; `v` in der Datei; Leerzeichen; CRLF), `p3-tag-ungleich-datei`, `p3-notes-fehlen`, `p3-notes-zelle-leer`, `p3-notes-pipe` (Zelle mit `|` und mit `\|`: Exit 2 `notes`, im Lauf und im Trockenlauf), `p3-vorabversion-fuehrende-null` (`v1.2.3-01` abgelehnt; `v1.2.3-0a` zugelassen), `p3-trockenlauf-notes-fehlen` (grün mit Hinweis), `p3-trockenlauf-notes-zelle-leer` (rot), `p3-trockenlauf-notes-vorhanden` (grün, Hinweis mit dem Text), `p3-schluesselwort-<Grund>` (je Schlüsselwort ein Fall), `p3-nicht-im-standard-zweig`, `p3-trockenlauf-nicht-im-standard-zweig` (grün mit Hinweis), `p3-secret-fehlt-user`,
   `p3-secret-fehlt-token`, `p3-secret-fehlen-beide` (eine Meldung nennt beide Namen; Mutation: nur der erste), `p3-reihenfolge-<Paar>` (je Paar benachbarter Gründe ein Fall mit beiden Verstößen: die Meldung nennt den früheren). Mutationen: je Form-Element, je Prüfung (2) bis (7) entfernt, je Paar vertauscht; Vergleich ohne `v`; Build-Metadaten zugelassen; leere Vorabversion zugelassen; Vorabversion mit führender Null zugelassen; leere dritte Zelle angenommen (Lauf und Trockenlauf); eine Zelle mit `|` wird am ersten `|` abgeschnitten angenommen; der Trockenlauf wertet die Zeile nicht aus oder lässt eine leere Zelle durch; ein Schlüsselwort weicht ab.
4. **Binaries (Schritt `binaries`; A).** `make binaries VERSION=<Version> DIST=<neues Verzeichnis außerhalb des Arbeitsbaums>`; danach `sha256sum -c SHA256SUMS` in diesem Verzeichnis grün, genau die sechs Dateien [`SPEC-055`](SPEC-055-container-und-binaries.md) Punkt 3 plus `SHA256SUMS`, sonst Exit 1. Es gibt keinen zweiten Bau und keinen Vergleich zweier Builds (das hält `make test-container`).
   Fälle: `p4-aufruf` (der Aufruf von `make` trägt `VERSION=<Version>` und ein `DIST` außerhalb des Arbeitsbaums), `p4-summe-falsch` (Attrappe schreibt eine falsche Summe: rot), `p4-datei-fehlt`, `p4-datei-zu-viel`. Mutationen: `VERSION` fehlt im Aufruf; `sha256sum -c` entfällt; die Dateiliste wird nicht gezählt.
5. **Image (Schritt `image`; B).** Ein Bau für beide Plattformen in **ein** OCI-Layout (`buildx build` mit Treiber `docker-container`, dessen Image `moby/buildkit` in `release.sh` per Digest gepinnt steht, in der Form `moby/buildkit@sha256:<64 Hexziffern>` mit dem Tag als Kommentar hinter der Zeile, `# v0.33.1`; Version und Digest sind eine Messung vom 2026-10-19 (Index-Digest von `v0.33.1`, älter als zwei Wochen) und altern: der Pin ist eine Handlung bei der Pflege, kein Wert dieser Kennung; ein Platzhalter aus Nullen ist ungültig; `--platform linux/amd64,linux/arm64`, `--target runtime`, `--build-arg VERSION=<Version>`,
   `--build-arg SOURCE_DATE_EPOCH=<Committer-Zeit des Commits `RELEASE_SHA` in Sekunden (`git log -1 --format=%ct`)>`, `--provenance=false --sbom=false`, `--output type=oci,dest=…,rewrite-timestamp=true`). Das Layout enthält genau einen Index mit genau zwei Manifesten (`linux/amd64`, `linux/arm64`); `release.sh` berechnet seinen
   Index-Digest `D` (SHA-256 des Index-Blobs). Das Binary `/pgwire-recorder` jedes Manifests hat die Summe des Eintrags für `linux/<Arch>` in `SHA256SUMS` aus Punkt 4 (Export des Dateisystems je Plattform aus demselben Builder), sonst Exit 1 mit Plattform. Das Image trägt kein Binary aus `dist/`, sondern wird aus demselben Quellstand
   neu gebaut; die Gleichheit der Summen ist die Prüfung (das Binary ist reproduzierbar, `LH-QA-01`). Das Image wird nicht als reproduzierbar zugesagt; das Setzen von `SOURCE_DATE_EPOCH` ist ein Versuch, der den Wiederanlauf (Punkt 11) erleichtert; ob BuildKit den Wert mit `rewrite-timestamp=true` ohne `ARG` im `Dockerfile` anwendet, ist nicht gemessen und nicht zugesagt (Grenze). **Grenze zur Form des Layouts:** Angenommen ist, dass `index.json` des Layouts genau einen Index-Blob nennt und dieser zwei Manifeste trägt; gemessen ist nur der Bau für eine Plattform im Treiber `docker`, der Bau für zwei Plattformen im Treiber `docker-container` zeigt erst der Trockenlauf auf dem Runner; jede andere Form ist rot (Exit 1), nie raten. **Stelle des Scans:** zwischen Bau und Push; der Scan liest das Layout, und der Push in
   Punkt 6 veröffentlicht genau dieses Layout (`D` ist unverändert); `slice-erster-release-freigabe` trägt den Schritt ein. Fälle: `p5-aufruf` (alle genannten Argumente, je Element: Plattformen, Ziel, `VERSION`, `SOURCE_DATE_EPOCH`, kein Provenance), `p5-zwei-manifeste` (Attrappe liefert drei Manifeste: rot), `p5-plattform-fehlt`,
   `p5-summe-amd64`, `p5-summe-arm64` (Binary im Image weicht ab: rot mit Plattform), `p5-digest` (`D` wird aus dem Layout berechnet und an die folgenden Schritte weitergereicht), `p5-builder-pin-form` (die Zeile trägt `@sha256:` und 64 Hexziffern und einen Tag-Kommentar), `p5-builder-platzhalter` (alle Ziffern Null: rot), `p5-epoch-committer` (Commit mit verschiedener Autor- und Committer-Zeit: der Wert ist die Committer-Zeit), `p5-layout-form` (`index.json` mit zwei Indexen oder ohne Index: Exit 1), `p5-ohne-binaries` (Aufruf ohne `binaries.ok`: Exit 2). Mutationen: je Argument entfällt; die Prüfung der Manifeste entfällt; die Prüfung der Summe je Plattform entfällt (je Plattform); `D` wird aus einem anderen Ort gelesen; Platzhalter oder fehlender Tag-Kommentar wird angenommen; die Autor-Zeit statt der Committer-Zeit; eine abweichende Layout-Form wird angenommen; die Marke `binaries.ok` wird nicht verlangt.
6. **Veröffentlichen (Schritt `veröffentlichen`; B).** Reihenfolge: Anmeldung bei `ghcr.io` (Benutzer `github.actor`, Token `GITHUB_TOKEN`, Rechte `packages: write`), Push nach `ghcr.io`, Anmeldung bei `docker.io` (`DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`), Push nach `docker.io`. Der Push kopiert das Layout aus Punkt 5 ohne Neubau
   (`skopeo copy --all --preserve-digests oci-archive:… docker://<Registry>/pt9912/pgwire-recorder:<Version>`), das heißt, beide Registries bekommen dieselben Bytes und dieselbe Manifestliste. **Tag-Satz je Version:** `<Version>` (SemVer ohne `v`), bei einem stabilen Tag zusätzlich `latest` (Punkt 8); keine Tags `<MAJOR>` oder `<MAJOR>.<MINOR>`.
   **Ein vorhandener Tag wird nie überschrieben:** liegt `<Version>` in einer Registry schon mit dem Index-Digest `D`, wird sie übersprungen (Wiederanlauf); liegt er mit einem anderen Digest da, ist der Lauf rot (Exit 1, Registry und beide Digests in der Meldung), und der Weg ist ein höherer Tag (Punkt 11). Das Anlegen des Repositorys auf Docker Hub und die Sichtbarkeit des Pakets
   in `ghcr.io` (öffentlich, mit dem Repository verknüpft) sind Handlungen des Betreibers; Punkt 7 prüft die Folge. **Ob ein Tag fehlt**, entscheidet die Meldung von `skopeo inspect`: `manifest unknown`, `name unknown` oder `not found` (ohne Beachtung der Schreibweise) heißt „fehlt“, jede andere Meldung ist rot (Exit 1) vor dem Push dorthin. **Ausgabe:** Bei jedem Fehler von `skopeo inspect` schreibt das Skript dessen Fehlerausgabe mit dem Präfix `release: skopeo: ` in das Log des Schritts (ein Fehler enthält keine Zugangsdaten; das Skript gibt keinen Wert eines Geheimnisses aus, Punkt 10); Fall `p6-inspect-meldung-im-log` (die Attrappe meldet einen Text, er steht im Log, auch im Fall „Tag fehlt“), Mutation: die Meldung wird verworfen. **Grenze:** Die Wortlaute sind nicht gegen `ghcr.io` und `docker.io` gemessen (`skopeo` ist im Dev-Container und am Host nicht installiert); der Trockenlauf ruft `skopeo` nicht, der erste Lauf mit einem Tag zeigt sie; weicht einer ab, endet der Lauf vor dem Push rot (fail-closed), und der Vertrag ist nachzuziehen. Fälle: `p6-tag-fehlt-manifest-unknown`, `p6-tag-fehlt-name-unknown`, `p6-tag-fehlt-not-found`, `p6-inspect-anderer-fehler` (Exit 1, kein Copy), `p6-reihenfolge` (Aufrufe der Attrappe: login ghcr, copy ghcr, login docker, copy docker), `p6-all-und-digests` (`--all --preserve-digests`), `p6-tag-vorhanden-gleich` (übersprungen, Exit 0), `p6-tag-vorhanden-anders` (Exit 1, kein Copy dorthin),
   `p6-zweite-registry-faellt` (Copy nach `docker.io` scheitert: Exit 1, `ghcr.io` bleibt beschrieben, Meldung nennt den Zustand), `p6-tag-satz` (nur `<Version>`, keine weiteren Tags). Mutationen: Reihenfolge vertauscht; `--all` entfällt; `--preserve-digests` entfällt; ein vorhandener Tag wird überschrieben; ein vorhandener gleicher Tag ist rot; zusätzlicher Tag `<MAJOR>.<MINOR>`; jeder Fehler von `skopeo inspect` gilt als „fehlt“; ein Wortlaut der drei wird nicht erkannt.
7. **Vergleichen und Smoke (Schritte `vergleichen`, `smoke`; B), fail-closed.** Nach dem Push lesen ohne Anmeldung (`--no-creds`) den Index beider Registries: Der Digest ist die SHA-256 der **rohen** Bytes des Index (`skopeo inspect --raw … | sha256sum`), nicht eine Anzeige eines Werkzeugs. Rot (Exit 1), wenn einer der drei Werte
   `D`, Digest `ghcr.io`, Digest `docker.io` abweicht, wenn der anonyme Zugriff scheitert (Paket nicht öffentlich: die Meldung nennt die Registry und die Handlung des Betreibers), wenn der Index nicht genau die zwei Plattformen enthält oder wenn der Tag `<Version>` fehlt. Kein Löschen, kein Zurücksetzen (Tags werden nicht verändert). Dann der Smoke (`container-smoke.sh`) mit der Referenz `<Registry>/pt9912/pgwire-recorder@<D>` aus **jeder** Registry für die Plattform
   des Docker-Servers (`amd64`): `version` gibt `pgwire-recorder <Version>`, Record und Replay nach [`SPEC-055`](SPEC-055-container-und-binaries.md) Punkt 7. Rot: weder `latest` noch GitHub-Release. **Zustand:** `vergleichen` legt bei Grün `vergleichen.ok` an, `smoke` bei Grün `smoke.ok`; `smoke` verlangt `vergleichen.ok`, `latest` und `github-release` verlangen `smoke.ok` (sonst Exit 2), unabhängig vom Workflow-Text. Das PostgreSQL-Image des Smokes steht in `harness/mk/integration.mk` (`POSTGRES_IMAGE`) und wird von dort gelesen; fehlt es, Exit 2 mit dem Schlüsselwort `vorbedingung`. Es wird anonym vom Docker Hub gezogen (Limit anonymer Pulls); ein Fehler dort ist rot (Exit 1) und bleibt es (kein Anmelden, keine Geheimnisse), die Registries sind dann beschrieben, `latest` und Release stehen aus, ein Wiederanlauf wiederholt Vergleich und Smoke (Punkt 11). Fälle: `p7-gleich` (grün), `p7-digest-ghcr`, `p7-digest-docker`, `p7-digest-lokal` (je einer weicht ab: rot mit Namen), `p7-anonym-ghcr`, `p7-anonym-docker` (Zugriff ohne Anmeldung scheitert), `p7-eine-plattform` (Index mit einem Manifest),
   `p7-tag-fehlt`, `p7-rohbytes` (Attrappe liefert zwei Anzeigen desselben Index in verschiedener Formatierung: gleich, weil die rohen Bytes gleich sind), `p7-smoke-<Registry>` (der Smoke scheitert für eine Registry: rot, kein `latest`), `p7-postgres-fehlt` (Quelle fehlt: Exit 2 `vorbedingung`), `p7-smoke-ohne-vergleichen` (Exit 2), `p7-smoke-ok-marke` (Marke nur bei Grün). Mutationen: je Vergleich entfällt (drei); der anonyme Zugriff nutzt die Anmeldung; der Digest kommt aus einer formatierten Anzeige; die Plattformzahl wird nicht geprüft;
   der Smoke läuft nur für eine Registry; die Marke `smoke.ok` entfällt oder wird bei Rot angelegt; `smoke` verlangt `vergleichen.ok` nicht.
8. **`latest` (Schritt `latest`; B).** Nur wenn (a) die Version stabil ist (kein `-`), (b) Punkt 7 für beide Registries grün war und (c) die Version höher ist als jede andere stabile Version unter den Tags `v*` des Repositorys (`git tag -l 'v*'` mit allen Tags, Vergleich nach SemVer, Vorabversionen zählen nicht). Dann setzt der Schritt `latest` auf `D` in beiden Registries
   (`skopeo copy --all --preserve-digests docker://…@<D> docker://…:latest`), `ghcr.io` zuerst. Eine Vorabversion und der Patch eines älteren Zweigs lassen `latest` unverändert (Exit 0 mit Meldung); scheitert das Setzen in der zweiten Registry, ist der Lauf rot und die Meldung nennt, welche Registry auf welchen Digest zeigt (Punkt 11). Das Beispiel in `LH-FA-16.a` nennt eine Referenz mit Registry und Version
   und sagt nur, dass `latest` auf den neuesten stabilen Release zeigt. Es zählen nur Tags der Form `v` plus stabiles SemVer (Punkt 3 (2) ohne Vorabversion); alle anderen Tags `v*` (etwa `vfoo`, `v1.2`) werden ignoriert. Der Schritt kopiert `latest` immer, auch wenn die Registry schon auf `D` zeigt (der Wiederanlauf ändert dann nichts); zugesagt ist, dass danach beide Registries `latest` auf `D` führen. Fälle: `p8-fremder-tag` (`vfoo`, `v1.2` ignoriert), `p8-ohne-smoke` (Exit 2), `p8-stabil-hoechste` (gesetzt, beide Registries), `p8-vorabversion` (nicht gesetzt), `p8-aelterer-zweig` (`v1.0.1` bei vorhandenem `v1.1.0`: nicht gesetzt), `p8-erster-stabiler` (kein anderer stabiler Tag: gesetzt), `p8-vorabversion-hoeher-zaehlt-nicht` (`v1.2.0-rc.1` vorhanden, `v1.1.1`
   stabil und höchste stabile: gesetzt), `p8-nach-rotem-vergleich` (Aufruf nach rotem Punkt 7: Exit 2), `p8-zweite-registry-faellt`. Mutationen: ein fremder Tag zählt; `latest` ohne `smoke.ok`; Vorabversion setzt `latest`; Vorabversionen zählen als höher; die Prüfung auf die höchste Version entfällt; `latest` vor dem Vergleich; Reihenfolge der Registries vertauscht.
9. **GitHub-Release (Schritt `github-release`; A, Vorbedingung `smoke` von B).** Zuletzt, nach den Registries: `gh release create v<Version> --verify-tag --title v<Version> --notes-file <Datei>` mit `--prerelease` genau bei einer Vorabversion, nicht als Entwurf, mit den Anhängen **genau**: die sechs Binaries und `SHA256SUMS` aus Punkt 4, sonst nichts. Der Text ist die dritte Zelle der Zeile (genau drei Zellen, kein `|` im Text, Punkt 3 (5)) der Änderungshistorie des Benutzerhandbuchs
   (die Tabelle im Abschnitt, dessen Überschrift auf `Änderungshistorie` endet, Zeile mit erster Zelle gleich der Version ohne `v`); ohne Zeile Eingabefehler (Punkt 3 (5)). Jeder Aufruf von `gh` trägt `--repo $RELEASE_REPOSITORY`. Ein fehlendes `GH_TOKEN` endet mit Exit 2 vor dem ersten Aufruf. Ob der Release schon besteht, entscheidet allein die Meldung `release not found` von `gh release view`; **jeder andere Fehler** (Netz, Rechte, unbekannte Meldung) ist Exit 1 und legt nichts an. Bei mehreren Zeilen derselben Version gilt die erste (nicht zugesagt, die Änderungshistorie trägt eine Zeile je Version). Existiert der Release schon (Wiederanlauf), werden fehlende Anhänge hochgeladen und `SHA256SUMS` am Release mit der lokalen Datei verglichen; eine abweichende Summe ist rot (Exit 1), nie überschrieben. Fälle: `p9-anhaenge` (genau sieben Dateien), `p9-stabil-ohne-prerelease`,
   `p9-vorabversion-prerelease`, `p9-notes-aus-zelle`, `p9-repo-argument` (jeder Aufruf mit `--repo`), `p9-gh-token-fehlt` (Exit 2, `gh` nie gerufen), `p9-view-anderer-fehler` (Exit 1, kein `create`), `p9-release-existiert-gleich` (fehlende Anhänge werden ergänzt), `p9-release-existiert-abweichend` (rot), `p9-nicht-vor-registries` **(B)** (Aufruf, bevor `smoke` grün war: Exit 2; die Zustandsdateien legen `release.sh vergleichen` und `release.sh smoke` an, B ergänzt die Fixture der Fälle von A um `smoke.ok` und ändert keinen Fall von A). Mutationen: ein achter Anhang; `--prerelease` immer oder nie; Notes aus der falschen Zelle; überschreiben statt vergleichen; Release vor den Registries; `--repo` fehlt; ohne `GH_TOKEN` wird `gh` gerufen; jeder Fehler von `gh release view` gilt als „nicht vorhanden“. **Grenze:** Der Wortlaut `release not found` ist der von `gh` auf dem Runner-Image; er ist im Dev-Container nicht gemessen (kein `gh`); der Trockenlauf ruft `gh` nicht, der erste Lauf mit einem Tag zeigt ihn, und weicht er ab, endet der Lauf rot (fail-closed) und der Wiederanlauf hilft nicht; dann ist der Vertrag nachzuziehen.
10. **Zugangsdaten (A: `GITHUB_TOKEN`, Ausgabe; B: Registries).** Die Namen sind `DOCKERHUB_USERNAME` und `DOCKERHUB_TOKEN` (Zugriffstoken mit Lesen und Schreiben) als **Repository-Geheimnisse**; `ghcr.io` nutzt `GITHUB_TOKEN` mit `packages: write`; der GitHub-Release nutzt `GITHUB_TOKEN` mit `contents: write`. Kein anderes Geheimnis, kein Eintrag im Repository, kein `environment:`
   (Punkt 1; Entscheidung des Nutzers vom 2026-10-19: Repository-Geheimnisse statt eines Environments, [ADR-0051](../../docs/plan/adr/0051-veroeffentlichung-auf-github-actions.md)). Anmeldung nur per `--password-stdin`; das Skript gibt keinen Wert eines Geheimnisses aus (kein `set -x`, kein `env`, keine Echo-Zeile über Zugangsdaten), auch nicht bei einem Fehler. Der Token wird nur den Schritten übergeben, die sich anmelden (Punkt 1). Fälle: `p10-nicht-im-log` (die Attrappe merkt sich die Token-Zeichenketten; weder stdout noch stderr von `release.sh` tragen sie, in keinem
   Schritt, auch nicht im Fehlerfall `p6-tag-vorhanden-anders`), `p10-stdin` (Anmeldung der Attrappe mit `--password-stdin`, nie als Argument), `p10-fehlendes-geheimnis` (Punkt 3 (7)). Mutationen: das Skript gibt den Token aus (bei Fehler); die Anmeldung übergibt den Token als Argument; `set -x` im Skript. Ob die Rechte der Token ausreichen (Schreibrecht auf Docker Hub, Verknüpfung des Pakets mit dem
   Repository), belegt nur der erste Lauf (Punkt 12).
11. **Wiederanlauf und Fehlertabelle (A: Zeilen `eingaben`, `binaries`, `github-release`; B: die übrigen).** Wiederanlauf eines roten Laufs ist „Re-run failed jobs“ derselben Ausführung oder ein erneuter Lauf desselben Tags; die Schritte sind so geschrieben, dass der zweite Lauf nichts verändert, was der erste hinterließ (Punkte 6, 8, 9). Ein bereits gesetzter Git-Tag wird nie verändert oder gelöscht ([`docs/maintainer/releasing.md`](../../docs/maintainer/releasing.md) §11); ein Fehler, der nur ein höherer Tag behebt (ein vorhandener Tag mit anderem Digest, ein abweichender Anhang), bekommt einen neuen Tag mit höherer Version.
   Die Tabelle steht im Dokument `releasing.md` §11 und deckt sich mit dieser:

   | Schritt | Fehler | hinterlässt | Wiederanlauf |
   |---|---|---|---|
   | `eingaben` | jede Eingabe | nichts | Geheimnis oder Einstellung beheben und wiederholen; liegt der Fehler im Tag-Commit (`VERSION`, Zeile der Änderungshistorie), neuer Tag |
   | `binaries`, `image` | Bau, Summe | nichts außerhalb des Runners | wiederholen |
   | `veröffentlichen` | `ghcr.io` ok, `docker.io` rot | Tag `<Version>` in `ghcr.io` | wiederholen: `ghcr.io` wird übersprungen (gleicher Digest), `docker.io` folgt; weicht der Digest ab (neuer Bau), neuer Tag |
   | `veröffentlichen` | Tag vorhanden mit anderem Digest | die vorhandene Registrierung | neuer, höherer Tag |
   | `vergleichen`, `smoke` | Digest, Zugriff, Plattformen, Smoke (auch Limit anonymer Pulls von PostgreSQL) | beide Registries beschrieben, kein `latest`, kein Release | Ursache beheben (z. B. Paket öffentlich stellen), wiederholen; bei Abweichung der Digests neuer Tag |
   | `latest` | zweite Registry rot | `latest` in einer Registry neu | wiederholen: setzt beide auf `D` |
   | `github-release` | Anlegen oder Anhang | Registries vollständig, Release fehlt oder teilweise | wiederholen: ergänzt fehlende Anhänge; abweichende Summe: neuer Tag |

   Fälle: `p11-wiederanlauf-ghcr-vorhanden` (zweiter Lauf nach `p6-zweite-registry-faellt`: `ghcr.io` übersprungen, `docker.io` kopiert, Exit 0), `p11-wiederanlauf-latest`, `p11-wiederanlauf-release-teilweise`, `p11-reihenfolge-gruende` (Fehler in zwei Schritten: der frühere Schritt endet den Lauf; ein späterer Schritt läuft nie nach einem roten). Mutationen: ein zweiter Lauf kopiert erneut über den vorhandenen Tag; er endet rot bei gleichem Digest; ein Schritt läuft nach rotem Vorgänger.
12. **Trockenlauf statt Probe-Tag; was er nicht belegt (A, ergänzt von B).** Der Trockenlauf (Punkt 2) belegt auf GitHub: Checkout, Werkzeuge des Runners (`docker buildx`, `skopeo`, `gh`), Eingaben, mit Slice B dass die Geheimnisse gesetzt sind, `make binaries`, mit Slice B den Bau und die Prüfung des Images, die Dauer. Er belegt **nicht**: das Schreibrecht der Token, die Verknüpfung und die Sichtbarkeit des Pakets in `ghcr.io`, die Existenz des Repositorys auf Docker Hub, das Anlegen des
    GitHub-Release, den Smoke aus den Registries. Diese Stellen zeigt erst der erste Lauf mit einem Tag; ein Fehler dort ist nach Punkt 11 behebbar, ohne einen Tag zu verändern. Der Trockenlauf wird einmal auf dem Stand der Lieferung jedes Slice gefahren und in §7 belegt (Lauf-Kennung, Dauer). **Ein Probe-Tag ist nicht vorgesehen** (Entscheidung des Nutzers vom 2026-10-19); die erste Außenwirkung in Registries und am GitHub-Release ist der Lauf für `v0.1.0`. Slice A belegt Eingaben, Binaries und Dauer; Slice B ergänzt den Bau und die Prüfung des Images (Punkt 5).
13. **Aufrufarten, Abbruch, Gleichzeitigkeit, Zustand zwischen den Schritten (A, ergänzt von B).** Der Zustand zwischen den Schritten liegt in `$RUNNER_TEMP/release`: das Verzeichnis `dist/` (Punkt 4) und die Marken `eingaben.ok`, `binaries.ok`, mit Slice B `vergleichen.ok`, `smoke.ok` und die Dateien des Images. `binaries` verlangt `eingaben.ok`, `github-release` verlangt `binaries.ok` und mit Slice B zusätzlich `smoke.ok` (sonst Exit 2); fehlt `RUNNER_TEMP`, Exit 2. Fälle `p13-binaries-ohne-eingaben`, `p13-release-ohne-binaries`; Mutation je Schritt: die Prüfung der Marke entfällt. `release.sh` läuft nie unter `make`; `-n`, `-q`, `-t`, `-i` und `MAKEFLAGS` haben keinen Fall (Grund: kein Hook und kein Stand liest das Ergebnis, die Pipeline ist der Aufrufer). Zwei Tags gleichzeitig laufen nacheinander: Die Gruppe `release-${{ github.event_name }}` trennt den Lauf (`push`) vom Trockenlauf (`workflow_dispatch`), sodass ein Trockenlauf einen Lauf nie verdrängt, und `cancel-in-progress: false` bricht einen **laufenden** Lauf nie ab. **Grenze (Verhalten von GitHub):** Gruppen halten höchstens einen laufenden und einen wartenden Lauf; kommt ein dritter hinzu, wird der wartende ersetzt (er endet als abgebrochen) — schiebt jemand drei Tags kurz hintereinander, kann der mittlere ohne Lauf bleiben. Das ist ein akzeptiertes Negativ (Grund: der Nutzer setzt Tags einzeln und nach dem Lauf des vorigen); der Weg ist „Re-run“ des abgebrochenen Laufs in der Oberfläche von GitHub (der Lauf ist dann derselbe Tag, Punkt 11). Die Zusage der Datei ist der Text von Punkt 1 (Fälle `p1-*`); ein Verhalten von GitHub prüft die Gegenprobe nicht. Bricht der Runner den Lauf ab (`SIGTERM`, Zeitgrenze 60 Minuten),
    gilt die Fehlertabelle von Punkt 11 für den Schritt, in dem er stand; das Skript räumt sein Temp-Verzeichnis bei Ende, `SIGINT` und `SIGTERM` und beendet seine Unterprozesse zuerst (wie [`SPEC-064`](SPEC-064-ci.md) Punkt 18; Fälle `p13-sigint`, `p13-sigterm` mit Mutationen). Die Gegenprobe selbst räumt ebenso.
14. **Aufbau der Gegenprobe (A, ergänzt von B).** `make release-gegenprobe` legt je Fall einen Temp-Baum mit einem Git-Repository (Commits, Tags, `origin/main`), Kopien von `release.sh`, `VERSION`, `docs/user/benutzerhandbuch.md`-Auszug und einem `PATH` an, in dem `docker`, `skopeo`, `gh`, `make`, `sha256sum` Attrappen sind, die ihre Aufrufe aufzeichnen; ohne Netz, ohne Docker. Der Arbeitsbaum bleibt unberührt. Der Text von `release.yml` wird von
    einer Funktion gehalten wie in [`SPEC-064`](SPEC-064-ci.md) Punkt 14. Keine Abdeckungs-Deklaration.
15. **Grenzen (A, ergänzt von B)** (akzeptierte Negative, jeweils mit Grund). Ein grüner Lauf auf GitHub ist ein Beleg, kein Gate. Die Geheimnisse, das Paket auf `ghcr.io`, das Repository auf Docker Hub und der Branch-Schutz sind Einstellungen des Betreibers, aus dem Klon nicht auditierbar. Das Image ist nicht reproduzierbar zugesagt (Punkt 5); der Wiederanlauf eines teilweise veröffentlichten Tags kann daran scheitern und endet dann rot mit dem Weg „neuer Tag“.
    Die Pipeline führt Werkzeuge des Runners (`skopeo`, `gh`, `buildx`) in der Version des Runner-Images aus; eine Änderung des Runner-Images kann das Verhalten ändern (Punkt 1: `ubuntu-24.04` fest, die Werkzeuge nicht). Ein Tag auf einem Commit, dessen `release.yml` verändert ist, führt die veränderte Datei aus (Verhalten von GitHub für Tag-Ereignisse): Wer taggen darf, bestimmt die Konfiguration des Repositorys. Das Homebrew-Tap, die Formel und der Scan des Images sind nicht Teil dieser Kennung.
