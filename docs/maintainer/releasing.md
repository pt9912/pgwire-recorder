# Releasing: Release-Prozess für Maintainer

Version: 0.7
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
Binaries mit `make binaries` und legt den GitHub-Release mit den Binaries und `SHA256SUMS` an (§6);
`make release-gegenprobe` prüft sie mit Attrappen für `make` und `gh`, ein Lauf auf GitHub ist noch
nicht belegt. Das Image veröffentlicht sie noch nicht, und eine Formel im Tap gibt es nicht. Es ist noch kein Tag
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
- **nur bei einem stabilen Tag** eine Formel im Homebrew-Tap
  `pt9912/homebrew-pgwire-recorder` für macOS und Linux; eine Vorabversion ändert
  den Tap nicht (`LH-FA-19.a`, `SPEC-042`).

**Reproduzierbar ist nur das Binary:** Zwei Builds desselben Quellstands für
dieselbe Plattform liefern dieselbe Prüfsumme (`LH-QA-01`). Für das Image gilt das
nicht; es wird über seinen Digest identifiziert.

Die Pipeline, die beim Tag die Binaries baut und an den GitHub-Release hängt, liefert
`slice-erster-release-veroeffentlichung`, das Image in beiden Registries
`slice-erster-release-veroeffentlichung-image`, den Scan des Images vor dem Push und den ersten
echten Tag `slice-erster-release-freigabe`; das Verfahren der Formel
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

Zwei Begriffe dieses Dokuments: Ein **Probe-Tag** ist ein Vorabversions-Tag
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
   das ist eine Handlung des Betreibers. Wie der Mechanismus sie liest, liefert
   `slice-erster-release-veroeffentlichung`.

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
   Benutzerhandbuchs die Zeile der Version trägt (§9) und dass der Commit des Tags im
   Standard-Zweig liegt. Jede Abweichung endet mit Exit 2, bevor etwas gebaut wird.
2. `Binaries` ruft `make binaries` mit der Version des Tags und prüft `SHA256SUMS` und die
   Dateiliste: die sechs Binaries und `SHA256SUMS`, sonst nichts. Es gibt keinen zweiten Bau;
   die Gleichheit zweier Builds hält `make test-container`.
3. `GitHub-Release` legt `v<SemVer>` mit den sieben Dateien an, mit dem Text aus der Zeile der
   Änderungshistorie, als Vorabversion bei einem Tag mit `-…`.

Dieselbe Datei läuft als **Trockenlauf**, wenn man sie auf einem Zweig von Hand startet
(`workflow_dispatch`): `Eingaben` und `Binaries` mit der Version aus `VERSION` als Tag
`v<VERSION>`, ohne GitHub-Release. Ein Probe-Tag gibt es nicht. Der Trockenlauf belegt nicht,
dass das Token am Repository schreiben darf und dass der GitHub-Release angelegt wird; das zeigt
erst der Lauf für den ersten Tag. Das Image (Bau, Veröffentlichung, Vergleich der Digests,
Smoke, `latest`) liefert `slice-erster-release-veroeffentlichung-image`; den Scan des Images vor
dem Push `slice-erster-release-freigabe`.

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

Das veröffentlichte Image besteht den Smoke aus Abnahmeszenario 9 mit
`linux/amd64` und `linux/arm64` aus beiden Registries; das liefert
`slice-erster-release-veroeffentlichung`. Weitere Kontrollen nach dem Push liefert
derselbe Slice.

## 9. Release-Notes

Das Projekt führt kein `CHANGELOG.md`. Die Änderungshistorie des
Benutzerhandbuchs (`docs/user/benutzerhandbuch.md`, letztes Kapitel) ist die einzige
Quelle: eine Zeile je Version in Betreibersicht, ohne interne Kennungen. Vor dem Tag
trägt der Release dort seine Zeile ein; der Text des GitHub-Releases entsteht aus
dieser Zeile. Der Text des GitHub-Releases ist die dritte Zelle der Zeile, deren erste Zelle die
Version ohne `v` ist, in der Tabelle des Abschnitts, dessen Überschrift auf `Änderungshistorie`
endet (Spalten Version, Datum, Änderung). Fehlt die Zeile, endet `Eingaben` mit Exit 2.

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
Die Tabelle nennt die Schritte der Pipeline, soweit sie bestehen; die Zeilen für die Schritte des
Images liefert `slice-erster-release-veroeffentlichung-image`.

| Schritt | Fehler | hinterlässt | Wiederanlauf |
|---|---|---|---|
| `eingaben` | jede Eingabe | nichts | Einstellung beheben und wiederholen; liegt der Fehler im Tag-Commit (`VERSION`, Zeile der Änderungshistorie), neuer Tag |
| `binaries` | Bau, Summe, Dateiliste | nichts außerhalb des Runners | wiederholen |
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
