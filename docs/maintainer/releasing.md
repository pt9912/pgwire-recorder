# Releasing: Release-Prozess für Maintainer

Version: 0.1
Stand: 2026-10-03

## 1. Zweck und Zielgruppe

Dieses Dokument beschreibt den **Mechanismus**, mit dem ein Release von
`pgwire-recorder` entstehen soll: wie ein Release ausgelöst wird, was dabei
entsteht und wie ein fehlerhafter Release behandelt wird. Es richtet sich an
Maintainer, die einen Release-Tag setzen; es liegt unter `docs/maintainer/`,
nicht unter `docs/user/`.

**Stand:** Das Projekt hat noch keinen Quelltext und keinen Release-Mechanismus.
Dieses Dokument legt fest, was ein Release liefert; die Automatisierung ist
noch nicht umgesetzt, und die Abschnitte nennen keinen Workflow und kein
`make`-Target, die es nicht gibt. Sobald ein Mechanismus existiert, wird er hier
beschrieben, einschließlich eines Belegs, dass er mit einem echten Tag
durchgelaufen ist.

Es ersetzt keine Anwenderdokumentation: Betrieb, Optionen, Exit-Codes und
Meldungskatalog stehen in der Betriebsdokumentation für Anwender.

## 2. Was ein Release liefert

Ein Release liefert:

- das Binary `pgwire-recorder` für Linux, macOS und Windows, jeweils für
  `amd64` und `arm64`,
- ein Docker/OCI-Image für `linux/amd64` und `linux/arm64` in einer
  Manifestliste, das das Binary ohne weitere Laufzeitabhängigkeit enthält und
  in `ghcr.io/pt9912/pgwire-recorder` sowie `docker.io/pt9912/pgwire-recorder`
  veröffentlicht wird,
- bei einem stabilen Release eine Formel im Homebrew-Tap
  `pt9912/homebrew-pgwire-recorder` für macOS und Linux; Vorabversionen ändern
  den Tap nicht.

Binary und Image müssen reproduzierbar gebaut werden können; das verlangt dieses
Release-Verfahren, das Lastenheft fordert nur die Bereitstellbarkeit
(`LH-QA-03`). Die zugesicherten Plattformen und das Image-Format legt die
Spezifikation fest (`SPEC-035`, `SPEC-031`).

## 3. Versionierung

Die Version der Software folgt SemVer 2.0 (`MAJOR.MINOR.PATCH`). Der Tag hat die
Form `v<SemVer>`, zum Beispiel `v0.1.0` oder `v0.2.0-rc.1`. Welche Versionsnummer
der erste Release trägt, wird beim Release festgelegt; die Produktstufe „v1“ des
Lastenhefts ist ein Umfang und keine Versionsnummer.

Die Version des Lastenhefts ist davon unabhängig; ihre Zählregel steht in dessen
Kopf.

Offen: Die Datei, die die Version der Software als Quelle der Wahrheit trägt,
wird festgelegt, wenn der Build existiert.

## 4. Einen Release auslösen

Festgelegt ist: Ein Release entsteht durch einen Git-Tag der Form `v<SemVer>`. Das
Verfahren ist noch nicht ausgeführt worden:

```bash
git tag v<SemVer>
git push origin v<SemVer>
```

Vor dem Tag gilt:

1. `make gates` ist grün.
2. Alle Slices der Welle, die den Release trägt, liegen in `done/`.
3. Die Versionsquelle (§3) trägt die neue Version und ist committet.

Wer taggen darf, bestimmt die Konfiguration des Repositorys, nicht dieses
Dokument.

## 5. Rollback

Es gibt keinen automatisierten Rollback. Ein fehlerhafter Release wird durch
einen neuen, höheren Tag mit einer korrigierten Version ersetzt; bereits
gesetzte Tags werden nicht verändert oder gelöscht.

### Änderungshistorie

| Version | Datum | Änderung |
|---|---|---|
| 0.1 | 2026-10-03 | Erste Fassung — beschreibt, was ein Release liefert; keine Automatisierung |
