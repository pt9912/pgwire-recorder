# Releasing: Release-Prozess für Maintainer

Version: 0.8
Stand: 2026-10-19

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
Binaries mit `make binaries`, baut das Image einmal für beide Plattformen, veröffentlicht es in
`ghcr.io` und `docker.io`, vergleicht die Digests, startet das Image aus beiden Registries,
setzt bei einem stabilen Tag `latest` und legt zuletzt den GitHub-Release mit den Binaries und
`SHA256SUMS` an (§6); `make release-gegenprobe` prüft sie mit Attrappen für `make`, `docker`,
`skopeo`, `gh` und den Smoke, ein Lauf auf GitHub für den Bau des Images und gegen die Registries
ist noch nicht belegt. Der Scan des Images fehlt, und eine Formel im Tap gibt es nicht. Es ist noch kein Tag
gesetzt worden. Was ein Slice erst liefert, nennt dieses Dokument mit seiner
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
nicht in dieser Pipeline (§7).

**Reproduzierbar ist nur das Binary:** Zwei Builds desselben Quellstands für
dieselbe Plattform liefern dieselbe Prüfsumme (`LH-QA-01`). Für das Image gilt das
nicht; es wird über seinen Digest identifiziert.

Die Pipeline, die beim Tag die Binaries baut und an den GitHub-Release hängt, liefert
`slice-erster-release-veroeffentlichung`, das Image in beiden Registries
`slice-erster-release-veroeffentlichung-image`; den Scan des Images vor dem Push und den ersten
echten Tag liefert `slice-erster-release-freigabe`, das Verfahren der Formel
`slice-v1-abschluss-homebrew`.

**Was lokal läuft:** `make binaries` baut je Zielplattform ein Binary in das Verzeichnis
`dist/` und schreibt `SHA256SUMS`; die Namen tragen `VERSION` (Vorgabe `dev`), unter Windows
mit `.exe`, und `dist/` muss leer sein oder fehlen. `make test-container` baut die Binaries
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
dieser Inhalt ist. `make binaries VERSION=<SemVer>` trägt die Version in die Namen der Binaries
und ins Binary; ohne Angabe bleibt es bei `dev`.

Zwei Begriffe, die nur das Verfahren der Formel (§7) benutzt; für die Veröffentlichung der
Binaries (§6) gibt es keinen Probe-Tag, dort ersetzt ihn der Trockenlauf: Ein **Probe-Tag** ist ein Vorabversions-Tag
(`v<SemVer>-probe.<N>`), mit dem das Release-Verfahren geprüft wird, ohne einen
Release zu veröffentlichen. Ein **Probe-Tap** ist ein Tap außerhalb von
`pt9912/homebrew-pgwire-recorder`, in den die Formel eines Probe-Tags geschrieben
wird. Beide benutzt `slice-v1-abschluss-homebrew`, um das Verfahren der Formel zu
prüfen.

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

Wer taggen darf, bestimmt die Konfiguration des Repositorys, nicht dieses
Dokument.

## 5. Vor dem Tag

1. `make gates` ist grün.
2. Die Änderungshistorie des Benutzerhandbuchs trägt die Zeile der neuen Version
   (§9).
3. Die Datei `VERSION` (§3) trägt die neue Version und ist committet.

Das Freshness-Audit der Baseline und von d-check vor dem Tag liefert
`slice-erster-release-freigabe`; die Prüfung des Tags gegen die Versionsdatei läuft in der
Pipeline (§6).

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
   einem Index und zwei Manifesten (Builder `docker-container`, Image des Builders per Digest
   gepinnt) und prüft, dass das Binary `/pgwire-recorder` jeder Plattform die Summe aus `SHA256SUMS` hat;
   der Digest `D` des Index steht danach fest. Das Image ist nicht reproduzierbar zugesagt.
4. `Veröffentlichen` kopiert genau dieses Layout ohne Neubau nach `ghcr.io`, dann nach `docker.io`, unter dem Tag
   `<SemVer>` ohne `v`. Ein vorhandener Tag mit demselben Digest wird übersprungen, einer mit einem anderen
   Digest ist rot und wird nie überschrieben.
5. `Vergleichen` liest den Index beider Registries ohne Anmeldung als rohe Bytes; `D`, der Digest aus
   `ghcr.io` und der aus `docker.io` müssen gleich sein, der Index nennt genau `linux/amd64` und
   `linux/arm64`. Ist das Paket nicht öffentlich, ist der Schritt rot.
6. `Smoke` startet das Image `<Registry>/pt9912/pgwire-recorder@<D>` aus jeder der beiden Registries
   (Aufzeichnen und Wiedergeben, `tools/harness/container-smoke.sh`).
7. `Latest` setzt `latest` auf `D`, in `ghcr.io` zuerst, nur bei einer stabilen Version, die höher ist als jede
   andere stabile Version unter den Tags `v*`. Eine Vorabversion und der Patch eines älteren Zweigs lassen
   `latest` unverändert.
8. `GitHub-Release` legt `v<SemVer>` mit den sieben Dateien an, mit dem Text aus der Zeile der
   Änderungshistorie, als Vorabversion bei einem Tag mit `-…`; der Schritt läuft nur, wenn `Vergleichen` grün war.

Dieselbe Datei läuft als **Trockenlauf**, wenn man sie auf einem Zweig von Hand startet
(`workflow_dispatch`): `Eingaben`, `Binaries` und `Image` mit der Version aus `VERSION` als Tag
`v<VERSION>`, ohne Anmeldung, ohne Veröffentlichung und ohne GitHub-Release; fehlt dort die Zeile der Änderungshistorie oder liegt der Commit
nicht im Standard-Zweig, steht das als Hinweis `release: hinweis: …`, kein Fehler. Ein Probe-Tag gibt es nicht. Der Trockenlauf belegt nicht,
dass die Token schreiben dürfen, dass das Paket in `ghcr.io` öffentlich ist und das Repository auf Docker Hub
besteht, dass der Smoke aus den Registries gelingt und dass der GitHub-Release angelegt wird; das zeigt
erst der Lauf für den ersten Tag. Den Scan des Images vor dem Push liefert `slice-erster-release-freigabe`.

## 7. Homebrew

Die Formel installiert das veröffentlichte Binary des Release-Tags für macOS und
Linux (`amd64`, `arm64`) und prüft es gegen die SHA-256-Summe desselben Releases;
sie baut nicht aus dem Quelltext (`LH-FA-19.a`). Sie entsteht nur bei einem stabilen
Tag; Windows bedient Homebrew nicht.

- Das Verfahren, die Formel aus den Release-Binaries zu erzeugen und im Tap
  abzulegen, liefert `slice-v1-abschluss-homebrew`, geprüft mit einem Probe-Tag und
  einem Probe-Tap (§3).
- Den Nachweis mit dem echten Tap (Abnahmeszenario 11) liefert
  `slice-erster-release-homebrew-nachweis`.

## 8. Nach dem Push kontrollieren

Der Smoke aus Abnahmeszenario 9 läuft in der Pipeline aus beiden Registries auf dem Docker-Server des
Runners (`ubuntu-24.04`, `amd64`; §6, Schritt `Smoke`); für `linux/arm64` prüft die Pipeline die Plattformen des
Index und die Summe des Binaries, sie startet das Image nicht. Nach dem Push ist zu kontrollieren, dass der Lauf grün endete und `latest` bei einem stabilen Tag
auf den Digest des Laufs zeigt; weitere Kontrollen liefert `slice-erster-release-freigabe`.

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

Die Checkliste, mit der ein Release freigegeben wird, liefert
`slice-erster-release-freigabe`.

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
| `veröffentlichen` | `ghcr.io` ok, `docker.io` rot | Tag `<SemVer>` in `ghcr.io` | wiederholen: `ghcr.io` wird übersprungen (gleicher Digest), `docker.io` folgt; weicht der Digest ab (neuer Bau), neuer Tag |
| `veröffentlichen` | Tag vorhanden mit anderem Digest | die vorhandene Registrierung | neuer, höherer Tag |
| `vergleichen`, `smoke` | Digest, Zugriff, Plattformen, Smoke | beide Registries beschrieben, kein `latest`, kein GitHub-Release | Ursache beheben (zum Beispiel das Paket öffentlich stellen), wiederholen; bei Abweichung der Digests neuer Tag |
| `latest` | zweite Registry rot | `latest` in einer Registry neu | wiederholen: setzt beide auf `D` |
| `github-release` | Anlegen oder Anhang | der Release fehlt oder trägt einen Teil der Anhänge | wiederholen: ergänzt fehlende Anhänge; weicht `SHA256SUMS` am Release von der lokalen Datei ab, neuer Tag |

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
