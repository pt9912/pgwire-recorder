# Releasing: Release-Prozess für Maintainer

Version: 0.12
Stand: 2026-10-20

## 1. Zweck und Zielgruppe

Dieses Dokument beschreibt den **Mechanismus**, mit dem ein Release von
`pgwire-recorder` entstehen soll: was ein Release liefert, wie er ausgelöst wird,
was vor und nach dem Tag zu prüfen ist und wie ein fehlerhafter Release behandelt
wird. Es richtet sich an Maintainer, die einen Release-Tag setzen; es liegt unter
`docs/maintainer/`, nicht unter `docs/user/`.

**Stand:** Der Quelltext von Record, Replay und Konfiguration existiert. Das
`Dockerfile` baut das Binary in mehreren Stufen, die Stufe `runtime` ist das
Produkt-Image; `make build` baut beide für die Plattform des Hosts
(`pgwire-recorder:dev`), `make binaries` baut die Binaries aller Zielplattformen mit
`SHA256SUMS`, und `make test-container` prüft sie (§2). Die Pipeline
`.github/workflows/release.yml` reagiert auf einen Tag `v*`: Sie prüft die Eingaben, baut die
Binaries mit `make binaries`, baut das Image einmal für beide Plattformen, scannt es vor dem Push,
veröffentlicht es in `ghcr.io` und `docker.io`, vergleicht die Digests, startet das Image aus beiden Registries,
setzt bei einem stabilen Tag `latest` und legt zuletzt den GitHub-Release mit den Binaries und
`SHA256SUMS` an (§6); `make release-gegenprobe` prüft sie mit Attrappen für `make`, `docker`,
`skopeo`, `gh`, den Scanner und den Smoke; der Scan, die Registries und der GitHub-Release hatten noch
keinen Lauf auf GitHub. Der Job `homebrew` legt bei einem stabilen Tag die Formel im Tap ab (§6, §7) und hatte ebenfalls noch keinen Lauf; was im echten Tap liegt, prüft kein Gate dieses Repos. Es ist noch kein Tag gesetzt worden. Was ein Slice erst liefert, nennt dieses Dokument mit seiner
Kennung („liefert `slice-…`“), nicht als Zusage. Der Beleg, dass die Pipeline mit einem echten
Tag durchgelaufen ist, steht hier erst, wenn es ihn gibt.

Es ersetzt keine Anwenderdokumentation: Betrieb, Optionen, Exit-Codes und
Meldungskatalog stehen in der Betriebsdokumentation für Anwender.

## 2. Was ein Release liefert

Ein Release liefert:

- das Binary `pgwire-recorder` für Linux, macOS und Windows, jeweils für `amd64`
  und `arm64` (`LH-QA-03`, `SPEC-035`), am GitHub-Release mit einer Datei
  `SHA256SUMS`, die die SHA-256-Summe jedes Binaries trägt; gegen sie prüft die
  Homebrew-Formel (`LH-FA-19.a`);
- ein Docker/OCI-Image für `linux/amd64` und `linux/arm64` in **einer**
  Manifestliste, das das Binary ohne weitere Laufzeitabhängigkeit enthält, in
  `ghcr.io/pt9912/pgwire-recorder` und `docker.io/pt9912/pgwire-recorder`
  (`LH-FA-16.a`, `SPEC-031`);

Eine Formel im Homebrew-Tap `pt9912/homebrew-pgwire-recorder` (`LH-FA-19.a`, `SPEC-042`) entsteht
bei einem stabilen Tag im eigenen Job `homebrew` derselben Pipeline, nach dem GitHub-Release (§6, §7).

**Reproduzierbar ist nur das Binary:** Zwei Builds desselben Quellstands für
dieselbe Plattform liefern dieselbe Prüfsumme (`LH-QA-01`). Für das Image gilt das
nicht; es wird über seinen Digest identifiziert.

Die Pipeline, die beim Tag die Binaries baut und an den GitHub-Release hängt, liefert
`slice-erster-release-veroeffentlichung`, das Image in beiden Registries
`slice-erster-release-veroeffentlichung-image`; den Scan des Images vor dem Push, die Freigabe-Checkliste (§10)
und den ersten echten Tag liefert `slice-erster-release-freigabe`, das Verfahren der Formel
`slice-v1-abschluss-homebrew`.

**Was lokal läuft:** `make binaries` baut je Zielplattform ein Binary in das Verzeichnis
`dist/` und schreibt `SHA256SUMS`; die Namen sind `pgwire-recorder-<os>-<arch>` ohne Version, unter Windows
mit `.exe`; `VERSION` (Vorgabe `dev`) steht im Binary (`version`), und `dist/` muss leer sein oder fehlen, auch
bei Dateien eines Laufs mit anderer `VERSION`. `make test-container` baut die Binaries
mit `VERSION=0.0.0-test` zweimal in Verzeichnisse außerhalb des Arbeitsbaums und vergleicht
die Summen, prüft den Dateikopf jedes Binaries, baut das Image für `linux/amd64` und
`linux/arm64` und prüft Betriebssystem, Architektur und Summe der Datei darin, und startet
das Image der Plattform des Docker-Servers: Aufzeichnen und Wiedergeben, mit entferntem
PostgreSQL beim Wiedergeben. Beide laufen nacheinander; gleichzeitige Läufe sind nicht
zugesagt. Eine Manifestliste entsteht dabei nicht.

**Abnahme-Umgebungen:** Gestartet und geprüft ist das Image für `linux/amd64` auf einem
Docker-Server mit `amd64`. Die Binaries für `darwin` und `windows` und das Image für
`linux/arm64` sind gebaut; belegt sind Bau, Dateikopf und Summe, ausgeführt sind sie nicht.

## 3. Versionierung, Probe-Tag, Probe-Tap

Die Version der Software folgt SemVer 2.0 (`MAJOR.MINOR.PATCH`). Der Tag hat die
Form `v<SemVer>`, zum Beispiel `v0.1.0` oder `v0.2.0-rc.1`. Ein **stabiler Tag**
hat keinen Vorabversions-Teil; jeder Tag mit `-…` ist eine **Vorabversion**.
Welche Versionsnummer der erste Release trägt, wird beim Release festgelegt; die
Produktstufe „v1“ des Lastenhefts ist ein Umfang und keine Versionsnummer.

Die Version des Lastenhefts ist davon unabhängig; ihre Zählregel steht in dessen
Kopf.

Die Version der Software steht in der Datei `VERSION` im Wurzelverzeichnis: genau eine Zeile
`<SemVer>` ohne `v`, Anfangswert `0.1.0-dev`. Die Pipeline lehnt einen Tag ab, der nicht `v` plus
dieser Inhalt ist. `make binaries VERSION=<SemVer>` trägt die Version ins Binary, nicht in die Namen der Binaries;
ohne Angabe bleibt es bei `dev`.

Es gibt keinen **Probe-Tag**: Für die Veröffentlichung (§6) ersetzt ihn der Trockenlauf, für die
Formel (§7) das Probe-Tap. Ein Probe-Tag wäre eine Vorabversion, und eine Vorabversion ändert den Tap nicht.
Ein **Probe-Tap** ist ein lokales Bare-Repository im Temp-Verzeichnis eines Falls von
`make release-gegenprobe`; die Gegenprobe legt die Formel dort ab und liest nach dem Lauf Zweig, Commits und
Dateien (§7). Sie greift weder auf `pt9912/homebrew-pgwire-recorder` noch auf das Netz zu und nutzt kein
echtes Geheimnis.

## 4. Vorbedingungen

Vor dem ersten Release gilt:

1. Alle Slices von `welle-v1-abschluss` liegen in `done/`.
2. `slice-harness-abdeckung-gate` und `slice-harness-coverage` liegen in `done/`
   (Start-Trigger von `welle-erster-release`; Abnahmeszenario 17 und damit M3 liegen
   vor dem ersten Release).
3. Der Release selbst ist Gegenstand von `welle-erster-release`.
4. Die Zugangsdaten für `ghcr.io` und `docker.io` sind im Repository hinterlegt;
   das ist eine Handlung des Betreibers: die Repository-Geheimnisse `DOCKERHUB_USERNAME` und
   `DOCKERHUB_TOKEN` (Zugriffstoken mit Lesen und Schreiben), das Repository auf Docker Hub und das
   Paket in `ghcr.io`, öffentlich und mit dem Repository verknüpft. Die Pipeline liest die Geheimnisse
   nur in den Schritten, die sie brauchen; fehlt oder ist leer eines der beiden, endet `Eingaben` mit Exit 2,
   bevor etwas veröffentlicht wird. Ob die Token schreiben dürfen und das Paket öffentlich ist, zeigt erst der
   Lauf für den ersten Tag.
5. Für die Formel (§7) ist das Repository-Geheimnis `HOMEBREW_TAP_GITHUB_TOKEN` hinterlegt, ebenfalls eine
   Handlung des Betreibers. Der Token braucht mindestens das Recht, in `pt9912/homebrew-pgwire-recorder`
   Inhalte zu schreiben und nichts darüber hinaus; bei einem Fine-grained Personal Access Token heißt das:
   nur dieses eine Repository, *Contents: Read and write*, mit Ablaufdatum. Der Job `homebrew` liest das Geheimnis
   als einziger, nur im Schritt `Homebrew`. Fehlt es oder ist es leer, endet der Schritt bei einem stabilen Tag mit
   Exit 2 und der Tap bleibt unverändert; ob der Token das Schreibrecht hat, zeigt erst der Lauf für einen stabilen
   Tag (Abnahmeszenario 11).

Wer taggen darf, bestimmt die Konfiguration des Repositorys, nicht dieses
Dokument.

## 5. Vor dem Tag

1. `make gates` ist grün; die Gegenproben laufen auf demselben Commit (§10, Punkte 2 und 3).
2. Die Änderungshistorie des Benutzerhandbuchs trägt die Zeile der neuen Version
   (§9).
3. Die Datei `VERSION` (§3) trägt die neue Version und ist committet; dieser Commit ist der, den der Tag trägt.
4. Der Trockenlauf auf diesem Stand ist grün (§6), ohne Hinweis zur Zeile der Änderungshistorie.

Die Prüfung des Tags gegen die Versionsdatei läuft in der Pipeline (§6).

**Freshness-Audit** (optional, empfohlen; kein Punkt der Freigabe-Checkliste und keine Vorbedingung des Tags): Der Maintainer
darf von Hand am Netz die Release-Listen der beiden Quellen abfragen, deren Stand das Repository pinnt, die der Baseline
(`pt9912/ai-harness-course`, gepinnter Tag in `harness/conventions.md` §Baseline) und die von d-check (`pt9912/d-check`,
gepinnter Tag in `d-check.mk`, `DCHECK_IMAGE`). Es gibt dafür kein Gate und keinen Fall der Gegenprobe. Der Tag hängt weder an
der Durchführung noch an einem Zeitpunkt noch an einem Ergebnis: **ein nicht durchgeführtes Audit sperrt den Tag nicht, und ein
Ergebnis sperrt ihn nie.** **Frisch** heißt: Der gepinnte Stand ist der neueste stabile Release der Liste.

Nur wenn das Audit durchgeführt wird, steht es als eigener Absatz ohne Punktnummer im Eintrag der Checkliste (§10), mit einer
Tabelle, je Quelle eine Zeile:

| Quelle | gepinnter Stand | neuester stabiler Stand | Datum der Abfrage | Ergebnis |
|---|---|---|---|---|
| Baseline | Tag aus `harness/conventions.md` | Tag der Liste | Datum | `aktuell` oder `neuerer Stand <Tag>` mit einem Satz zur Auswirkung |
| d-check | Tag aus `d-check.mk` | Tag der Liste | Datum | wie oben |

Eine Frist für das Datum der Abfrage gibt es nicht. Ein neuerer Stand wird in der Tabelle dokumentiert, und sein Nachziehen ist
ein eigener Slice. Nennen die Release-Notizen des neueren Standes eine Sicherheitskorrektur, legt der Maintainer die Frage dem
Nutzer vor; der Nutzer entscheidet, ob vor oder nach dem Tag nachgezogen wird, und ohne seine Antwort bleibt der Tag, wie er ist.

## 6. Was der Tag auslöst

Ein Release entsteht durch einen Git-Tag der Form `v<SemVer>`:

```bash
git tag v<SemVer>
git push origin v<SemVer>
```

Der Push löst `.github/workflows/release.yml` aus; ein Tag läuft nur dort (`ci.yml` lässt Tags
aus). Die Schritte, die Logik steht in `tools/harness/release.sh`:

1. `Eingaben` prüft vor jedem Schritt mit Außenwirkung, in dieser Reihenfolge: Ereignis und Ref,
   die Form des Tags `v<SemVer>` (Build-Metadaten sind nicht zugelassen), die Form der Datei
   `VERSION`, dass der Tag `v` plus ihr Inhalt ist, dass die Änderungshistorie des
   Benutzerhandbuchs die Zeile der Version mit nicht leerer dritter Zelle trägt (§9) und dass
   der Commit des Tags im Standard-Zweig liegt. Jede Abweichung endet mit Exit 2, bevor etwas
   gebaut wird; jede Meldung beginnt mit `release: ` und einem Schlüsselwort, das den Grund nennt.
2. `Binaries` ruft `make binaries` mit der Version des Tags und prüft `SHA256SUMS` und die
   Dateiliste: die sechs Binaries und `SHA256SUMS`, sonst nichts. Es gibt keinen zweiten Bau;
   die Gleichheit zweier Builds hält `make test-container`.
3. `Image` baut das Image einmal für `linux/amd64` und `linux/arm64` in ein OCI-Layout mit genau
   einem Index und zwei Manifesten (Treiber `docker-container`, das Image des Builders mit Digest gepinnt) und prüft in einem zweiten Bau je Plattform (`type=local`, derselbe Builder), dass das Binary `/pgwire-recorder` die Summe aus `SHA256SUMS` hat;
   das kopierte Layout liest diese Prüfung nicht. Der Digest `D` des Index steht danach fest. Das Image ist nicht reproduzierbar zugesagt.
4. `Scan` scannt dieses Layout vor dem Push mit Trivy (Image per Digest in `release.sh` gepinnt, `docker run` ohne
   Ausnahmedatei) je Plattform, `linux/amd64` und `linux/arm64`, und lässt `D` unverändert. Trivy liest ein Layout nur als
   Verzeichnis und bei zwei Manifesten das erste; der Schritt entpackt das Layout deshalb in ein Verzeichnis je Plattform, dessen
   `index.json` nur das Manifest dieser Plattform nennt (fehlt es oder steht es mehrfach im Index: `plattform`, vor dem Aufruf von
   `docker`; lässt sich das Layout nicht entpacken oder lesen: `scan-gescheitert`). Ein Eintrag mit dem
   Schweregrad `CRITICAL` oder `HIGH` und einer behobenen Version beendet den Lauf mit Exit 1 und nennt alle Befunde
   (`release: scan: schwachstelle: …`); Einträge ohne behobene Version stehen nur in der Summenzeile
   (`<m> ohne Behebung`) und sperren nicht. Ein Scan, der nicht lief, endet mit Exit 2 und einem Stichwort: Trivy scheitert
   (`scan-gescheitert`), kein JSON oder kein Ergebnis (`pruefmenge`), die Datenbank ist älter als 72 Stunden oder ihr
   Alter ist nicht lesbar (`datenbank`), die gescannte Architektur weicht ab (`plattform`), das Layout fehlt
   (`vorbedingung`). Ausnahmen gibt es nicht. `Veröffentlichen` verlangt den grünen `Scan`.
5. `Veröffentlichen` kopiert genau dieses Layout ohne Neubau nach `ghcr.io`, dann nach `docker.io`, unter dem Tag
   `<SemVer>` ohne `v`. Ein vorhandener Tag mit demselben Digest wird übersprungen, einer mit einem anderen
   Digest ist rot und wird nie überschrieben.
6. `Vergleichen` liest den Index beider Registries ohne Anmeldung als rohe Bytes; `D`, der Digest aus
   `ghcr.io` und der aus `docker.io` müssen gleich sein, der Index nennt genau `linux/amd64` und
   `linux/arm64`. Ist das Paket nicht öffentlich, ist der Schritt rot.
7. `Smoke` startet das Image `<Registry>/pt9912/pgwire-recorder@<D>` aus jeder der beiden Registries
   (Aufzeichnen und Wiedergeben, `tools/harness/container-smoke.sh`).
8. `Latest` setzt `latest` auf `D`, in `ghcr.io` zuerst, nur nach grünem `Smoke` und nur bei einer stabilen Version, die höher ist als jede
   andere stabile Version unter den Tags `v*`. Eine Vorabversion und der Patch eines älteren Zweigs lassen
   `latest` unverändert.
9. `GitHub-Release` legt `v<SemVer>` mit den sieben Dateien an, mit dem Text aus der Zeile der
   Änderungshistorie, als Vorabversion bei einem Tag mit `-…`; der Schritt läuft nur, wenn `Smoke` grün war.
10. Der Job `homebrew` (eigener Job mit `needs: release`, nur bei einem Push, `contents: read`) erzeugt die Formel aus
    der Vorlage `tools/harness/homebrew-formula.rb.tmpl` und der Datei `SHA256SUMS` des GitHub-Release und legt sie im
    Tap ab (§7). Ein roter Job `homebrew` lässt den Release und die Registries unberührt; den Wiederanlauf
    beschreibt §11.

Dieselbe Datei läuft als **Trockenlauf**, wenn man sie auf einem Zweig von Hand startet
(`workflow_dispatch`): `Eingaben`, `Binaries`, `Image` und `Scan` mit der Version aus `VERSION` als Tag
`v<VERSION>`, ohne Anmeldung, ohne Veröffentlichung und ohne GitHub-Release; fehlt dort die Zeile der Änderungshistorie oder liegt der Commit
nicht im Standard-Zweig, steht das als Hinweis `release: hinweis: …`, kein Fehler. Ein Probe-Tag gibt es nicht. Der Trockenlauf belegt nicht,
dass die Token schreiben dürfen, dass das Paket in `ghcr.io` öffentlich ist und das Repository auf Docker Hub
besteht, dass der Smoke aus den Registries gelingt und dass der GitHub-Release angelegt wird; das zeigt
erst der Lauf für den Tag. Der Scan braucht Netz für das Image des Scanners und die Schwachstellen-Datenbank und läuft deshalb
auch dort; der Job `homebrew` läuft dort nicht; ein grüner Scan sagt „am Tag des Laufs ist nichts Behebbares mit `HIGH` oder `CRITICAL` bekannt“, nicht „sicher“, und dieselbe
Quelle kann an zwei Tagen verschieden bewertet werden. Die Binaries am GitHub-Release werden nicht gescannt.

## 7. Homebrew

Die Formel installiert das veröffentlichte Binary des Release-Tags für macOS und
Linux (`amd64`, `arm64`) und prüft es gegen die SHA-256-Summe desselben Releases;
sie baut nicht aus dem Quelltext (`LH-FA-19.a`). Sie entsteht nur bei einem stabilen
Tag; Windows bedient Homebrew nicht.

Der Schritt `homebrew` (`tools/harness/release.sh homebrew`, Job `homebrew`) prüft in dieser Reihenfolge und endet
beim ersten Treffer: Ereignis und Ref (Exit 2, `ereignis`), die Form des Tags (Exit 2, `tag-form`), dann eine Vorabversion
(Exit 0, die Zeile `release: homebrew: <Tag> ist eine Vorabversion, der Tap bleibt unverändert`, ohne Aufruf von `gh` und
`git` und ohne Prüfung der Token), `GH_TOKEN` (Exit 2, `token`) und `HOMEBREW_TAP_GITHUB_TOKEN` (Exit 2, `geheimnis`; die
Meldung nennt den Namen, nie einen Wert). Danach:

1. Die Datei `SHA256SUMS` kommt aus dem GitHub-Release (`gh release download`). Sie muss für `darwin` und `linux`, je `amd64`
   und `arm64`, genau eine Zeile `<Summe>  pgwire-recorder-<os>-<arch>` mit einer SHA-256-Summe aus 64 Zeichen `0-9a-f`
   tragen, sonst Exit 1 (`summen`, `summen-laden`), bevor der Tap angefasst wird. Das alte Schema mit der Version im Namen
   (so trägt es `v0.1.0`) wird nicht angenommen: Für `v0.1.0` gibt es keine Formel, der Schritt endet dort mit Exit 1 `summen`.
2. Die Formel ist die Vorlage mit der Version ohne `v` und den vier Summen.
3. Der Tap wird in ein neues Temp-Verzeichnis geklont (Exit 1 `tap-klon`). Ein leerer Tap bekommt den Zweig `main`, ein Tap
   mit Commits muss auf `main` ausgecheckt sein (Exit 1 `tap-zweig`). Der Schritt schreibt nur
   `Formula/pgwire-recorder.rb` (Exit 1 `tap-pfad`, wenn `Formula` oder die Formel keine Datei bzw. kein Verzeichnis ist; Exit 1
   `tap-formel`, wenn die Formel im Tap nicht genau eine Zeile `  version "<X.Y.Z>"` mit einer stabilen Version hat). Trägt der
   Tap eine höhere Version, bleibt er unverändert (Exit 0); trägt er dieselbe Version mit demselben Text, ebenfalls (Exit 0,
   „schon aktuell“); sonst, also bei derselben Version mit anderem Text, bei einer niedrigeren Version und bei fehlender
   Formel, entsteht ein Commit `pgwire-recorder <Version>` mit Autor und Committer `github-actions[bot]` und ein
   `git push origin main` ohne `--force`. Weist der Tap den Push ab, endet der Schritt mit Exit 1 `tap-push` und der Tap
   bleibt unverändert; der Schritt wiederholt den Push nicht.
4. Der Token von `HOMEBREW_TAP_GITHUB_TOKEN` gelangt nur über `GIT_ASKPASS` zu `git` (`GIT_TERMINAL_PROMPT=0`,
   `-c credential.helper=` leer); er steht nie in einem Argument von `git` oder `gh`, nie in `.git/config` des Klons, in
   keiner Datei unter dem Temp-Verzeichnis und nie in der Ausgabe.

**Nachtrag von Hand** (Handlung des Betreibers, mit seinem eigenen Token): Wer einen Tag nachziehen muss, dessen Lauf den Job
nicht hatte oder dessen Formel die Vorlage nicht mehr trägt, ruft in einem Klon des Hauptzweigs
`bash tools/harness/release.sh homebrew` mit `RELEASE_EVENT=push`, `RELEASE_REF=refs/tags/v<Version>`,
`RELEASE_REPOSITORY=pt9912/pgwire-recorder`, `GH_TOKEN` und `RELEASE_TAP_TOKEN` in der Umgebung, nie als Argument. Der
Schritt unterscheidet das nicht vom Lauf des Jobs. Das Skript ist unter Linux zugesagt (bash ab 5.1, `git`, `gh`, GNU `awk`
und `sort`); macOS und Windows als Ausführungsort sind es nicht.

**Was `make release-gegenprobe` belegt:** den Text der Formel und der Vorlage, die Ablage und die Randformen des Taps gegen
ein lokales Probe-Tap (§3), die Reihenfolge der Gründe und den Umgang mit dem Token. **Was sie nicht belegt:** dass
Homebrew die Formel annimmt, dass die URLs am GitHub-Release erreichbar sind, dass das Binary auf macOS und Linux läuft, dass
das Tap-Token das Schreibrecht hat, dass `https://github.com/pt9912/homebrew-pgwire-recorder.git` die Form des Klons
trägt und den Wortlaut von `gh release download` auf dem Runner. Das zeigt erst `slice-erster-release-homebrew-nachweis` mit
dem echten Tap (Abnahmeszenario 11), die Installation aus dem Tap und der Abschnitt im Benutzerhandbuch gehören dorthin.

## 8. Nach dem Push kontrollieren

Der Smoke aus Abnahmeszenario 9 läuft in der Pipeline aus beiden Registries auf dem Docker-Server des
Runners (`ubuntu-24.04`, `amd64`; §6, Schritt `Smoke`); für `linux/arm64` prüft die Pipeline die Plattformen des
Index und die Summe des Binaries, sie startet das Image nicht. Nach dem Push ist zu kontrollieren, dass der Lauf grün endete und `latest` bei einem stabilen Tag
auf den Digest des Laufs zeigt; die Belege des Laufs trägt die Checkliste (§10).

## 9. Release-Notes

Das Projekt führt kein `CHANGELOG.md`. Die Änderungshistorie des
Benutzerhandbuchs (`docs/user/benutzerhandbuch.md`, letztes Kapitel) ist die einzige
Quelle: eine Zeile je Version in Betreibersicht, ohne interne Kennungen. Vor dem Tag
trägt der Release dort seine Zeile ein; der Text des GitHub-Releases entsteht aus
dieser Zeile. Der Text des GitHub-Releases ist die dritte Zelle der Zeile, deren erste Zelle die
Version ohne `v` ist, in der Tabelle des Abschnitts, dessen Überschrift auf `Änderungshistorie`
endet (Spalten Version, Datum, Änderung). Die Zeile hat genau drei Zellen; ein weiteres `|` im Text, auch als `\|`, ist ein
Fehler. Fehlt die Zeile oder ist die Zelle leer, endet `Eingaben`
mit Exit 2; im Trockenlauf ist eine fehlende Zeile ein Hinweis, eine leere Zelle ein Fehler.

## 10. Freigabe-Checkliste

Ein Release ist freigegeben, wenn jeder Punkt einen Beleg trägt. Der **Eintrag** eines Releases ist eine Kopie der Tabelle mit
ausgefüllter Spalte *Beleg*: für den ersten Release steht er in der Closure von `welle-erster-release`
(`welle-erster-release-results.md`), für spätere Releases ohne Welle in `docs/maintainer/freigaben/<Version>.md`. Ein Punkt ohne
Beleg ist nicht erfüllt; was ein Lauf nicht zeigt, steht als Grenze im Eintrag, nie als Behauptung. Das Freshness-Audit (§5) ist
kein Punkt der Tabelle; ein durchgeführtes steht als eigener Absatz im Eintrag, ohne Punktnummer.

Der **getaggte Commit** ist der Commit, der `VERSION` und die Zeile der Änderungshistorie trägt; die Punkte 1 bis 3 gelten für genau
diesen Hash.

| Punkt | Form des Belegs | Beleg |
|---|---|---|
| 1. `make gates` grün auf dem getaggten Commit | Hash, Datum | |
| 2. Gegenproben-Container grün für diesen Hash | die Dateien `gegenproben-*-<Hash>` der Ablage des Gegenproben-Containers per `ssh`, jede mit `ergebnis=gruen`, `beginn` und `ende`, keine Marke `….laeuft`; fehlt das Ergebnis oder ist es rot, fährt die Freigabe `make gegenproben` selbst | |
| 3. `make gegenproben-host` am Host grün auf demselben Hash | Datum, Hash, die Zusammenfassung der Ausgabe | |
| 4. `VERSION` und die Zeile der Änderungshistorie im getaggten Commit | Ausgabe von `grep` | |
| 5. Trockenlauf auf dem Stand | Lauf-Kennung, Dauer, ohne Hinweis zur Zeile der Änderungshistorie | |
| 6. Scan grün | die Zeile `release: scan: 0 Befund(e) über der Schwelle, <m> ohne Behebung` aus dem Log des Trockenlaufs und des Laufs des Tags | |
| 7. Ausdrückliches Ja des Nutzers zum Tag | Wortlaut, Datum | |
| 8. Lauf des Tags | Lauf-Kennung, Index-Digest je Registry (gleich), Smoke aus beiden Registries, `latest` | |
| 9. Anhänge am GitHub-Release | die sieben Dateien (sechs Binaries und `SHA256SUMS`) | |
| 10. Beobachtungen des ersten Laufs (unten) | je Beobachtung Lauf-Kennung und Zeile im Log, oder die Grenze | |
| 11. Bezug nach dem Tag (unten) | Befehl, Lauf und Ausgabe je Befehl im Handbuch | |

**Anti-Punkte** (was ein Release nicht tut): Ein Tag wird nicht verändert oder gelöscht. Kein Image geht ohne grünen Scan in eine
Registry. Eine Vorabversion setzt kein `latest`. Es gibt keinen Tag ohne ausdrückliches Ja des Nutzers (Punkt 7). Kein Tag hängt an einem Freshness-Audit, weder an seiner Durchführung noch an seinem Ergebnis. Keine Datei kommt von Hand
an den GitHub-Release. Kein Wert eines Geheimnisses steht in einem Log, einem Eintrag oder einer Nachricht.

**Incident-Klausel:** Ein veröffentlichter Release mit einem Fehler (abweichender Digest, bekannte behebbare Schwachstelle, falsches
Binary) wird durch einen neuen, höheren Tag ersetzt. Auslöser ist, wer den Fehler bemerkt; entschieden wird mit dem Nutzer. Festgehalten
wird er im Eintrag des fehlerhaften Release (Datum, Befund, Ersatz-Tag) und als Zeile in der Änderungshistorie des Handbuchs; der
fehlerhafte Tag bleibt stehen.

**Den Tag setzt, wer taggen darf, und nur nach dem Ja des Nutzers (Punkt 7).** Ein Implementer startet ihn nicht. Das Ja zum ersten
Tag schließt ein, dass `v0.1.1` ihm folgen kann, wenn der Lauf einen Fehler im Skript oder im Vertrag zeigt (unten); jeder weitere Tag braucht ein
neues Ja.

**Beobachtungen des ersten Laufs (Punkt 10).** Was der Trockenlauf nicht belegt (§6), zeigt der Lauf des Tags im Log; der Eintrag nennt
je Beobachtung die Lauf-Kennung und die Zeile:

- die Wortlaute von `skopeo inspect` für einen nicht vorhandenen Tag in `ghcr.io` und in `docker.io` (Zeilen `release: skopeo: `
  in den Schritten `Veröffentlichen` und `Vergleichen`) und die Verfügbarkeit von `skopeo` und `gh` auf dem Runner;
- der Wortlaut von `gh release view` für einen nicht vorhandenen Release (Schritt `GitHub-Release`) und das Anlegen des Release;
- das Schreibrecht der Token in beiden Registries, die Sichtbarkeit des Pakets in `ghcr.io` (anonymer Zugriff im Schritt `Vergleichen`)
  und der Smoke aus beiden Registries mit der Referenz `<Registry>/pt9912/pgwire-recorder@<D>` (Digest je Registry);
- die Zeitstempel des veröffentlichten Images als Beobachtung zur Wirkung von `SOURCE_DATE_EPOCH` (keine Zusage);
- der Wiederanstoß eines roten Laufs nach §11, nur wenn der Lauf rot endete; endet er grün, nennt der Eintrag den Wiederanstoß als
  nicht belegt.

**Rot nach dem ersten Tag.** Ein Lauf führt `release.yml` und `release.sh` des Commits aus, den der Tag trägt; ein Wiederanstoß desselben
Tags führt dieselbe Fassung aus. Liegt der Fehler in der **Umgebung** (Token, Schreibrecht, Sichtbarkeit des Pakets, Repository auf Docker
Hub, Limit anonymer Pulls, Netz), behebt der Betreiber ihn und stößt denselben Tag wieder an (§11). Liegt er im **Skript, im Text der Datei
oder im Vertrag** (ein Wortlaut weicht ab, eine Zusage war falsch), wird die Korrektur auf dem Hauptzweig gemacht und ein **neuer, höherer
Tag** gesetzt; der Tag des fehlerhaften Laufs bleibt, und jede Registry, die ihn trägt, behält ihn. Ein Digest, der abweicht, und ein
abweichender Anhang verlangen immer einen neuen Tag.

**Bezug nach dem Tag (Punkt 11).** Die Zeile der Version in der Änderungshistorie steht vor dem Tag (§9) und bleibt. Nach dem Lauf,
in einem eigenen Commit, der nur `docs/user/`, `README.md` und Planung ändert, nennt der Kopf des Benutzerhandbuchs als
*Software-Version* die Ausgabe von `pgwire-recorder version` des **veröffentlichten** Binaries (nicht von `make build`), und §2
*Installation* bekommt die Abschnitte zum Bezug der Binaries (Plattformen, `SHA256SUMS`, Prüfung mit `sha256sum -c`) und der Images
(`<Registry>/pt9912/pgwire-recorder:<Version>` aus beiden Registries) mit nur den Befehlen, die nach dem Tag einmal gelaufen sind. Ein Befehl
ohne diesen Lauf steht nicht im Handbuch. `latest` wird genannt, weil der Lauf es gesetzt und der Smoke es belegt hat. `darwin` und
`windows` sind gebaut, nicht ausgeführt (§2); das Handbuch sagt nur das.

## 11. Fehler, Wiederanlauf, Rollback

Es gibt keinen automatisierten Rollback. Ein fehlerhafter Release wird durch
einen neuen, höheren Tag mit einer korrigierten Version ersetzt; bereits
gesetzte Tags werden nicht verändert oder gelöscht.

Wiederanlauf eines roten Laufs ist „Re-run failed jobs“ derselben Ausführung oder ein erneuter
Lauf desselben Tags; die Schritte sind so geschrieben, dass der zweite Lauf nichts verändert,
was der erste hinterließ. Ein Fehler, den nur ein höherer Tag behebt, bekommt einen neuen Tag.
Die Tabelle nennt die Schritte der Pipeline.

| Schritt | Fehler | hinterlässt | Wiederanlauf |
|---|---|---|---|
| `eingaben` | jede Eingabe | nichts | Geheimnis oder Einstellung beheben und wiederholen; liegt der Fehler im Tag-Commit (`VERSION`, Zeile der Änderungshistorie), neuer Tag |
| `binaries`, `image` | Bau, Summe, Dateiliste | nichts außerhalb des Runners | wiederholen |
| `scan` | Befund über der Schwelle (Exit 1), Scan nicht gelaufen (Exit 2) | nichts außerhalb des Runners | Exit 2 durch Netz, Zwischenspeicher oder Datenbank (`scan-gescheitert` von Trivy, `datenbank`): wiederholen; `scan-gescheitert` für das Layout (nicht entpackbar, Index nicht lesbar) ist ein Fehler des Schrittes `image`: den Lauf ab `image` wiederholen, der Schritt `scan` allein liest dasselbe Layout; bleibt dieselbe Meldung, oder meldet er `plattform` oder `pruefmenge`, war eine Annahme über Trivy oder das Skript falsch: Korrektur auf dem Hauptzweig und ein neuer, höherer Tag (§10, „Rot nach dem ersten Tag“); Befund: beheben (Basis-Image oder Abhängigkeit heben) und einen neuen Tag setzen, es gibt keine Ausnahme |
| `veröffentlichen` | `ghcr.io` ok, `docker.io` rot | Tag `<SemVer>` in `ghcr.io` | wiederholen: `ghcr.io` wird übersprungen (gleicher Digest), `docker.io` folgt; weicht der Digest ab (neuer Bau), neuer Tag |
| `veröffentlichen` | Tag vorhanden mit anderem Digest | die vorhandene Registrierung | neuer, höherer Tag |
| `vergleichen`, `smoke` | Digest, Zugriff, Plattformen, Smoke | beide Registries beschrieben, kein `latest`, kein GitHub-Release | Ursache beheben (zum Beispiel das Paket öffentlich stellen), wiederholen; bei Abweichung der Digests neuer Tag |
| `latest` | zweite Registry rot | `latest` in einer Registry neu | wiederholen: setzt beide auf `D` |
| `github-release` | Anlegen oder Anhang | der Release fehlt oder trägt einen Teil der Anhänge | wiederholen: ergänzt fehlende Anhänge; weicht `SHA256SUMS` am Release von der lokalen Datei ab, neuer Tag |
| `homebrew` | Summen nicht ladbar oder unbrauchbar, Klon, Zweig, Pfad, Formel im Tap unlesbar, Push abgewiesen | Release und Registries vollständig, der Tap unverändert | Ursache beheben (zum Beispiel das Geheimnis, sein Recht oder die Formel im Tap), den Job wiederholen („Re-run failed jobs“ derselben Ausführung; ob GitHub dabei nur diesen Job wiederholt, zeigt kein Test dieses Repos) oder den Nachtrag von Hand fahren (§7); kein neuer Tag, auch nicht bei einem Fehler im Skript oder in der Vorlage: die Formel ist kein veröffentlichtes Artefakt, die Korrektur geschieht auf dem Hauptzweig, danach der Nachtrag |

### Änderungshistorie

| Version | Datum | Änderung |
|---|---|---|
| 0.1 | 2026-10-03 | Erste Fassung — beschreibt, was ein Release liefert; keine Automatisierung |
| 0.2 | 2026-10-03 | Abschnitt Release-Notes: Quelle ist die Änderungshistorie des Handbuchs, kein `CHANGELOG.md` |
| 0.3 | 2026-10-04 | Reproduzierbarkeit nur für das Binary; Plattformen und Registries aus dem Lastenheft |
| 0.4 | 2026-10-09 | Stand auf Code, `Dockerfile` und `make build` gezogen; neue Gliederung; `SHA256SUMS` und Formel nur bei stabilen Tags; Vorbedingungen um die Start-Bedingung von `welle-erster-release` ergänzt; was erst ein Slice liefert, mit Kennung |
| 0.5 | 2026-10-09 | Nach dem Schnitt des Release-Slice: Scan, Freshness-Audit, Checkliste und erster echter Tag bei `slice-erster-release-freigabe`; Fehlertabelle und Wiederanlauf bei `slice-erster-release-veroeffentlichung` |
| 0.6 | 2026-10-17 | `make binaries` und `make test-container` im Ist-Zustand; Abnahme-Umgebungen |
| 0.7 | 2026-10-19 | Pipeline für Tag und Binaries (`release.yml`, `release.sh`), Datei `VERSION`, Trockenlauf, Fehlertabelle und Wiederanlauf für `eingaben`, `binaries` und `github-release` |
| 0.8 | 2026-10-19 | Image: Bau, Veröffentlichung in beiden Registries, Vergleich der Digests, Smoke, `latest`; Geheimnisse der Registries; Fehlertabelle um die Schritte des Images ergänzt; die Formel gehört nicht zur Pipeline |
| 0.9 | 2026-10-19 | Scan des Images vor dem Push (Schritt `Scan`, Zeile der Fehlertabelle); Freshness-Audit mit der Form seines Belegs (§5); Freigabe-Checkliste mit Anti-Punkten und Incident-Klausel (§10) |
| 0.10 | 2026-10-19 | Freshness-Audit optional und kein Punkt der Freigabe-Checkliste mehr (§5); Checkliste mit elf Punkten neu nummeriert, Anti-Punkt zum Audit (§10) |
| 0.11 | 2026-10-20 | Namen der Binaries ohne Version: `pgwire-recorder-<os>-<arch>` (§2, §3) |
| 0.12 | 2026-10-20 | Formel im Tap: Job `homebrew` nach dem GitHub-Release, Geheimnis `HOMEBREW_TAP_GITHUB_TOKEN` als Vorbedingung (§4), Probe-Tap statt Probe-Tag (§3), Verfahren und Nachtrag von Hand (§7), Zeile der Fehlertabelle (§11) |
