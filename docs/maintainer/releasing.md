# Releasing: Release-Prozess für Maintainer

Version: 0.4
Stand: 2026-10-09

## 1. Zweck und Zielgruppe

Dieses Dokument beschreibt den **Mechanismus**, mit dem ein Release von
`pgwire-recorder` entstehen soll: was ein Release liefert, wie er ausgelöst wird,
was vor und nach dem Tag zu prüfen ist und wie ein fehlerhafter Release behandelt
wird. Es richtet sich an Maintainer, die einen Release-Tag setzen; es liegt unter
`docs/maintainer/`, nicht unter `docs/user/`.

**Stand:** Der Quelltext von Record, Replay und Konfiguration existiert. Das
`Dockerfile` baut das Binary in mehreren Stufen, die Stufe `runtime` ist das
Produkt-Image; `make build` baut beide für die Plattform des Hosts
(`pgwire-recorder:dev`). Einen Release-Mechanismus gibt es nicht: keinen Workflow,
der auf einen Tag reagiert, keinen Build für alle Zielplattformen, keine
Veröffentlichung in einer Registry und keine Formel im Tap. Es ist noch kein Tag
gesetzt worden. Was ein Slice erst liefert, nennt dieses Dokument mit seiner
Kennung („liefert `slice-…`“), nicht als Zusage. Sobald ein Mechanismus existiert,
wird er hier beschrieben, einschließlich eines Belegs, dass er mit einem echten Tag
durchgelaufen ist.

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

Den Build für alle Zielplattformen, die Prüfsummen und die Veröffentlichung in beiden
Registries liefert `slice-erster-release-veroeffentlichung`; das Image selbst
`slice-v1-abschluss-container`; das Verfahren der Formel `slice-v1-abschluss-homebrew`.

## 3. Versionierung, Probe-Tag, Probe-Tap

Die Version der Software folgt SemVer 2.0 (`MAJOR.MINOR.PATCH`). Der Tag hat die
Form `v<SemVer>`, zum Beispiel `v0.1.0` oder `v0.2.0-rc.1`. Ein **stabiler Tag**
hat keinen Vorabversions-Teil; jeder Tag mit `-…` ist eine **Vorabversion**.
Welche Versionsnummer der erste Release trägt, wird beim Release festgelegt; die
Produktstufe „v1“ des Lastenhefts ist ein Umfang und keine Versionsnummer.

Die Version des Lastenhefts ist davon unabhängig; ihre Zählregel steht in dessen
Kopf.

Welche Datei die Version der Software als Quelle der Wahrheit trägt, liefert
`slice-erster-release-veroeffentlichung`.

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
3. Die Versionsquelle (§3) trägt die neue Version und ist committet.

Weitere Prüfungen vor dem Tag liefert `slice-erster-release-veroeffentlichung`.

## 6. Was der Tag auslöst

Ein Release entsteht durch einen Git-Tag der Form `v<SemVer>`:

```bash
git tag v<SemVer>
git push origin v<SemVer>
```

Was der Push eines Tags auslöst — Build der Binaries für alle Zielplattformen mit
`SHA256SUMS`, Build und Veröffentlichung des Images als Manifestliste in beiden
Registries, GitHub-Release —, liefert `slice-erster-release-veroeffentlichung`.
Heute löst ein Tag nichts aus.

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
dieser Zeile. Das Übernehmen in den GitHub-Release liefert
`slice-erster-release-veroeffentlichung`.

## 10. Freigabe-Checkliste

Die Checkliste, mit der ein Release freigegeben wird, liefert
`slice-erster-release-veroeffentlichung`.

## 11. Fehler, Wiederanlauf, Rollback

Es gibt keinen automatisierten Rollback. Ein fehlerhafter Release wird durch
einen neuen, höheren Tag mit einer korrigierten Version ersetzt; bereits
gesetzte Tags werden nicht verändert oder gelöscht.

Wie ein abgebrochener Lauf wieder angestoßen wird, beschreibt dieses Dokument, sobald
der Mechanismus existiert (§1); heute gibt es keinen Lauf.

### Änderungshistorie

| Version | Datum | Änderung |
|---|---|---|
| 0.1 | 2026-10-03 | Erste Fassung — beschreibt, was ein Release liefert; keine Automatisierung |
| 0.2 | 2026-10-03 | Abschnitt Release-Notes: Quelle ist die Änderungshistorie des Handbuchs, kein `CHANGELOG.md` |
| 0.3 | 2026-10-04 | Reproduzierbarkeit nur für das Binary; Plattformen und Registries aus dem Lastenheft |
| 0.4 | 2026-10-09 | Stand auf Code, `Dockerfile` und `make build` gezogen; neue Gliederung; `SHA256SUMS` und Formel nur bei stabilen Tags; Vorbedingungen um die Start-Bedingung von `welle-erster-release` ergänzt; was erst ein Slice liefert, mit Kennung |
