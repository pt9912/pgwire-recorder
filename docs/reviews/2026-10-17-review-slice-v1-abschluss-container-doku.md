# Review-Report: slice-v1-abschluss-container-doku — 2026-10-17

**Review-Art:** Code (Doku) — gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** `slice-v1-abschluss-container-doku`, Diff `2e89f594..b303a342` über `docs/user docs/maintainer README.md` (Code-Commit `04edd0a8`, Plan `b303a342`).

**Skill:** `.harness/skills/reviewer.md` @ `b303a342`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

> **Zitier-Form.** Kennung, nicht Adresse: `slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code.

**Eingangs-Kontext:**

- Slice-Plan `slice-v1-abschluss-container-doku` (§1, §2, §3, §4, §6, §7, §8)
- Vertrag `SPEC-055` (`harness/vertraege/SPEC-055-container-und-binaries.md`), Gate `make binaries` und `make test-container`
- `LH-FA-16`, `LH-QA-03`
- `AGENTS.md` (§3.7, §3.9, §3.11, §3.13, §3.14)

**Eigene Läufe (nacheinander):** `make binaries` (leeres `dist/`): sechs Binaries und `SHA256SUMS`; `sha256sum -c` sechs Zeilen `OK`; `./dist/pgwire-recorder_dev_linux_amd64 version` gibt `pgwire-recorder dev`; zweiter `make binaries` bei nicht leerem `dist/`: `binaries: Eingabe: DIST 'dist' ist nicht leer`, Exit 2; `VERSION=1.2.3 make binaries`: Namen `pgwire-recorder_1.2.3_…`, `.exe` unter Windows, `version` gibt `pgwire-recorder 1.2.3`, `sha256sum -c` ohne Abweichung; `docker build --target runtime --platform linux/arm64` und `inspect` gibt `linux/arm64`; `make build` baut ohne `--platform`. `make meldungskatalog-check`, `make docs-check` (962 Dateien, 0 Befunde), `make kopf-check`, `make abdeckung-check`, `make gates-schnell`: Exit 0. Messung `git diff --numstat 2e89f594..b303a342 -- docs/user docs/maintainer README.md`: 7 + 34 + 63 = 104 (Marke 500, Meldeschwelle 400: nicht erreicht). Danach `dist/` und das Tag `dev-arm64` entfernt; keine Container, Netze, Tags `pgr-ct-*` oder `pgr-container.*` zurückgeblieben (`pgr-dev` und `pgr-dev-default` gehören dem anderen Lauf). `make test-container` selbst nicht gefahren; die Belege dazu liegen in §7 des Plans.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Das Handbuch sagt in §1, von „den übrigen Binaries“ sei geprüft, dass sie gebaut sind, „Ausgeführt sind sie nicht“; §2 startet das Binary `linux_amd64` aus `dist/` selbst. Die Bezugsmenge von „übrigen“ steht nirgends, die README-Zeile „die übrigen Binaries … gebaut, nicht ausgeführt“ hat ebenso keinen Bezug. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „Von den übrigen“; `README.md` · „die übrigen Binaries und das Image für `linux/arm64` sind gebaut, nicht ausgeführt“ | nein — kein Gate liest Prosa | Quantor „übrigen“ ohne benannte Bezugsmenge |
| F-2 | LOW | Der Block zur Prüfung der Summen endet mit `cd dist`; der folgende Block ruft `./dist/pgwire-recorder_dev_linux_amd64 version` vom Wurzelverzeichnis aus. Wer die Blöcke nacheinander in derselben Shell tippt, steht in `dist/` und findet den Pfad nicht. Der Lauf des Implementers startete jeden Block vom Wurzelverzeichnis. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „./dist/pgwire-recorder_dev_linux_amd64 version“ | ja — Handlauf der Blöcke in Folge | Beispielfolge nicht als Folge belegt |
| F-3 | LOW | Unverändert stehen zwei Sätze neben den neuen Zusagen: „Zum Ausführen: Linux auf der Architektur des Rechners, auf dem Sie gebaut haben, oder ein Container-Laufzeitsystem“ (das Handbuch bietet jetzt Binaries für macOS und Windows an) und „Das Binary im Image ist … für die Architektur des Rechners, auf dem Sie gebaut haben“ (der Abschnitt darunter baut ein Image für `linux/arm64` auf beliebigem Server). Der Plan nahm die Plattform-Zeile in §1 ausdrücklich auf. | `AGENTS.md` §3.9, §3.11 | `docs/user/benutzerhandbuch.md` · „Zum Ausführen: Linux auf der Architektur des“ und „Das Binary im Image ist ein statisch gebundenes Linux-Programm“ | nein | Bestandssatz widerspricht neuer Zusage |
| F-4 | LOW | §7 des Plans („Hinweis an Architect und Planner“) sagt „`SPEC-055` Punkt 9 nennt Volumes nicht“. Seit `240562bd` (vor `b303a342`) nennt Punkt 9 Volumes mit `rm -fv`, und `slice-harness-integration-volumes` übernimmt sie von `test-container`; der Hinweis nennt diesen Nehmer nicht. Behauptet wird das Leck nicht als gelöst (`tools/harness/run-container.sh` räumt weiter mit `rm -f`; offen bis jener Slice liefert). | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-v1-abschluss-container-doku.md` · „Punkt 9 nennt Volumes nicht“ | ja — `grep Volumes` in `SPEC-055` | Plan folgt der Korrektur nicht |
| F-5 | INFO | Das Handbuch nennt „Docker 29.9.0“ als Server des Laufs. Das Gate pinnt die Docker-Version nicht, der Satz ist eine Beobachtung des Implementers (§7, `docker version`) und veraltet mit dem Docker des Rechners. Er ist belegt, sagt aber nichts zu, was `make test-container` prüft. | `AGENTS.md` §3.11 | `docs/user/benutzerhandbuch.md` · „(Docker 29.9.0)“ | nein | Beobachtung als Zusage lesbar |
| F-6 | INFO | §7 des Plans trägt die Messung der Marke (104, Basis `2e89f594`) noch nicht, die Quantoren-Durchsicht nennt die `releasing.md`-Sätze „Beide laufen nacheinander; gleichzeitige Läufe sind nicht zugesagt“ und „Eine Manifestliste entsteht dabei nicht“ nicht; beide decken `SPEC-055` Punkt 6 und Grenze 5, kein Lauf des Gates. Die DoD verlangt jedes Beispiel „als Datei“; §7 belegt die neuen Shell-Beispiele durch Handläufe. Für die Closure zu füllen. | `AGENTS.md` §3.14, §3.11 | `docs/plan/planning/in-progress/slice-v1-abschluss-container-doku.md` · „Quantoren-Durchsicht“ | nein | Beleg-Tabelle unvollständig gegenüber den Sätzen |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/maintainer/releasing.md` (Version 0.6, Stand, Änderungshistorie, „keinen Build für alle Zielplattformen“ entfernt, Absätze *Was lokal läuft* und *Abnahme-Umgebungen* gegen `SPEC-055` Punkte 3 bis 7, Grenzen 3 bis 5) | geprüft, ohne Befund |
| Chronik und Verweise im Diff von Handbuch, README und `releasing.md` (grep: Slice, Welle, `ADR-`, `SPEC-`, `LH-`, Spezifikation, Lastenheft, Review, „noch nicht“, geplant, künftig, bisher) | geprüft, ohne Befund; Treffer nur die schon vorhandenen Kennungen `slice-erster-release-…` und `slice-v1-abschluss-homebrew` in `releasing.md` (dort erlaubt) |
| Handbuch: keine Aussage zu Registry, Tag `latest`, Manifestliste, Release-Anhängen | geprüft, ohne Befund (§1 Ausschluss eingehalten) |
| Handbuch-Beispiele zu `make binaries`, Namen, `VERSION`, `SHA256SUMS`, Fehlerfall, `docker build --platform linux/arm64` und `inspect` | geprüft, ohne Befund (eigener Lauf, Ausgaben wie im Text) |
| Quantoren („jede Zeile … `OK`“, „sechs“, „muss leer sein oder fehlen“, „nur das Image für `linux/amd64`“, „jedes Binaries“ in `releasing.md`) | geprüft, ohne Befund außer F-1 |
| Tabellen *Fehlercodes* und *Warnungen* | `make meldungskatalog-check` grün, ohne Befund |
| Plan §1, §3, §6, §8 gegen den Diff; Randformen aus §6 (Plattformen, `darwin`/`windows`, Registry, Aufruf-Voraussetzung) | geprüft, ohne Befund außer F-4 |
| `README.md` (kein Zielstand, keine Chronik) | geprüft, ohne Befund außer F-1 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Quantor „übrigen“ ohne benannte Bezugsmenge · Beispielfolge nicht als Folge belegt · Bestandssatz widerspricht neuer Zusage · Plan folgt der Korrektur nicht · Beobachtung als Zusage lesbar · Beleg-Tabelle unvollständig gegenüber den Sätzen

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Das Volume-Leck des Gates ist kein Befund dieses Diffs; es liegt bei `slice-harness-integration-volumes` und wird im Text nicht als gelöst behauptet.

**Übergabe:** F-1 bis F-3 an den Implementer dieses Slice, F-4 und F-6 an den Implementer (Plan, §7), F-5 zur Kenntnis; die Finding-Klassen gehen in die Slice-Closure §7. Die Konformität mit der DoD prüft der Verifier.
