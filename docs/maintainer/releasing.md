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

Binary und Image müssen reproduzierbar gebaut werden können. Die zugesicherten
Plattformen und die Anforderungen an Container und Reproduzierbarkeit stehen im
Lastenheft und in der Spezifikation (`LH-FA-16`, `LH-QA-03`).

## 3. Versionierung

Die Version der Software folgt SemVer 2.0 (`MAJOR.MINOR.PATCH`). Der Tag hat die
Form `v<SemVer>`, zum Beispiel `v1.0.0` oder `v1.1.0-rc.1`.

Die Version des Lastenhefts ist davon unabhängig; ihre Zählregel steht in
`harness/conventions.md`.

Offen: Die Datei, die die Version der Software als Quelle der Wahrheit trägt,
wird festgelegt, wenn der Build existiert.

## 4. Einen Release auslösen

Festgelegt ist: Ein Release entsteht durch einen Git-Tag der Form `v<SemVer>`. Das
Verfahren ist noch nicht ausgeführt worden. Die erste Version `v1.0.0` wird
getaggt, wenn der Meilenstein M3 der Roadmap erreicht ist (alle Anforderungen
umgesetzt):

```bash
git tag v1.0.0
git push origin v1.0.0
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
