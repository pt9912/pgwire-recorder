# Verifikationsbericht: slice-harness-devcontainer-austausch-repo — 2026-10-17

**Rolle:** Verifier (Modul 11) · **Frage:** Bauen wir es richtig, gegen Plan und DoD?
**Gegenstand:** Diff `38a9b20a..HEAD` (`9f8faef9` Dokument, `1e70abc3` nur Plan §7); Zweig `slice/devcontainer-austausch-repo`. Einen Review-Report gibt es zum Zeitpunkt dieses Berichts nicht (kein Lesen nach dem Urteil möglich).
**Maßstab:** DoD §2 des Plans, `AGENTS.md` §3.9, §3.11, §3.14.

**Verdikt: bestätigt.** Keine HIGH-, MEDIUM-Befunde; 2 LOW/INFO. Offen bleiben die Pflichten der Closure (Planner) und der Review-Report.

## Eigene Läufe

| Lauf | Ergebnis |
|---|---|
| `make docs-check` (Worktree) | 991 Dateien, 0 Befunde |
| `make kopf-check` | Exit 0 |
| Marke §3.14: `git diff --numstat 92cfaf23..HEAD -- docs/developer` | 151 + 100 = **251**, wie in §7; Meldeschwelle 350, Marke 450 nicht erreicht |
| Dokumentgrenze `wc -l docs/developer/devcontainer.md` | 549 (vorher 498), Grenze 560 nicht überschritten |
| `grep -n 'host-repo\|PGR_REPO'` im Dokument | leer |
| `make gates-schnell` | von mir nicht gefahren; Beleg in §7 (Exit 0 auf `9f8faef9`, 1 min 59,8 s); `1e70abc3` ändert nur den Plan (Positivliste). Beleg des Implementers |

## Kernprobe nach dem Dokument (eigener Container, nicht nach dem Bericht)

Verzeichnis im Scratchpad, Projekt `pgr-ver`, Port 2422, eigener Schlüssel. Dockerfile (Zeilen 110 bis 138), `start.sh` (151 bis 160) und der Compose-Block (174 bis 212) wörtlich aus dem Dokument gezogen; nur `pgr-dev` durch `pgr-ver` ersetzt, Kopf von Hand (kein `u-boot init`), `.env` mit Name und Port. Docker-Vorbedingung `pgrep -fa "^make gates"` leer.

| Zusage des Dokuments | Beobachtung |
|---|---|
| Build und Start, Container `healthy`; genau drei Mounts, Host-Arbeitsbaum nicht darunter (§3.5, §7, §11) | `healthy`, `docker inspect`: 3 (`volume:/workspaces`, `bind:/etc/pgr-dev/authorized_key.pub`, `volume:/var/lib/docker`) |
| `git init --bare --initial-branch=main /workspaces/austausch.git` als `vscode`; Rechte (§6) | angelegt, `vscode 1000 775`, git 2.55.0 |
| `remote add container …`, `core.sshCommand`, `git push container main` mit unveröffentlichtem Commit (§9) | `* [new branch] main -> main` (mit eingetragenem Hostschlüssel) |
| Clone im Container, `git config --global`, `git switch -c thema`, `git push origin thema` (§6) | alle Schritte ohne Fehler |
| `git fetch container` am Host (§9) | `thema` → `container/thema` |
| Präfix: im Container `container/thema2` → am Host `container/container/thema2` (§9) | bestätigt (`git branch -r`) |
| zweiter Push des Hosts, `git switch main && git pull` im Klon (§6) | Fast-Forward, Stand da |
| beide schreiben auf `main`: Host-Push abgewiesen, `fetch`, `merge container/main`, Push, `git pull` im Klon (§9, §12) | `rejected (fetch first)` (Hinweis `git pull` vor dem Push), nach Merge Push `d14e13f4..58b29bd3`, Pull im Klon holt `58b29bd3` |
| zweites `git remote add container …` (§6) | `externes Repository container existiert bereits`, Exit 3 |
| Bare-Repo und Klon überleben `down`/`up`; `~/.gitconfig` nicht; Hostschlüssel gleich; `fetch` ohne Warnung (§4, §5, §6) | `docker compose down` + `up -d`: `austausch.git` mit allen Branches, Klon `## main...origin/main`, `unable to read config file '/home/vscode/.gitconfig'`, gleicher Fingerabdruck, `git fetch` Exit 0 |

Nicht wiederholt (dem Dokument gemäß ohnehin Beleg des Implementers): Neubau ohne Cache, zweite Instanz, Abbruch im Push, `u-boot up`-Warnung. `u-boot down`/`up` nicht gefahren, weil mein Verzeichnis kein `u-boot.yaml` hat; `docker compose down`/`up` ersetzt sie, die Aussage des Dokuments (Container und Netz entfallen, Volumes bleiben) ist damit gedeckt.

Aufräumen: `docker compose down -v`, `docker rmi pgr-ver:latest`; `docker volume ls -q` vorher/nachher gleich, `docker ps -a` leer, kein Netz `pgr-ver*`, `pgr-dev:latest` unverändert (`d7d697d39d6d`); `~/.ssh/known_hosts` byte-gleich der Sicherung (`cmp`); `/Development/pgr-dev` und der Hauptbaum nicht angefasst.

## DoD-Konformität

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| Dokument beschreibt den Austausch nach §1; §3.5, §6, §7, §9, §10, §11, §12, §13 nachgezogen; `/host-repo`, `PGR_REPO` weg; jede Aussage belegt (Tabelle in Plan-§7) oder in §13 | erfüllt | Diff gelesen; `grep` leer; Kernprobe oben bestätigt die Kette; Compose-Block ohne Mount, §10 mit Remote je Container, §11 mit drei Mounts, §12 mit zwei neuen Zeilen, §13 mit nicht geprüften Randformen |
| `make docs-check`, `kopf-check`, `gates-schnell` grün | erfüllt, `gates-schnell` als Beleg | siehe oben |
| Review durchgeführt, Report liegt vor | offen | Datei zum Zeitpunkt nicht vorhanden |
| Closure-Notiz mit Lerneintrag | offen für den Planner | Plan-§7: „Noch nicht geschlossen“ |
| Register fortgeschrieben oder „keine Beobachtung“ | offen für den Planner | |
| Risiko-Ausgänge §6 | offen für den Planner | Jedes Risiko trägt „Ausgang: bei der Closure“; der Stand des Laufs steht dabei (Randformen belegt, Rest in Dokument-§13) |
| Drei Paarungen | offen für den Planner | |

## Plan-vs-Code (§3.9, §3.11, §3.14)

- Der Diff berührt nur das Dokument (Version 0.1 auf 0.2, Stand gesetzt) und den Plan; Dockerfile und `start.sh` im Dokument unverändert (Abgrenzung §1 gehalten); kein Produktcode, keine Hooks.
- §1, §3 und §6 des Plans folgen dem Diff (Festlegungen, Stand des Laufs, Probe 8: Scratch statt echter Container). Zählung §8 steht: 1 Liefer-Punkt, 1 Schicht.
- §3.11: Dokument ohne Kennungen dieses Repos (`grep` nach `slice-`, `ADR-`, Welle: nur Treffer in `u-boot`-Blöcken `LH-…`, im Dokument als solche genannt) und ohne Chronik; Aussagen reichen nicht weiter als die Probe, wo ich sie wiederholt habe. Quantoren-Durchsicht in Plan-§7 vorhanden.
- §3.14: Metrik in §4 genannt, Messung stimmt (251 gegen Schätzung 220 bis 340).

## Findings

| # | Schwere | Befund |
|---|---|---|
| F-1 | INFO | Das Dokument gibt als gelaufen an: Zeit „44 Sekunden mit `healthy`“ über `u-boot up`; ich fuhr `docker compose up -d` (Build 37 s, `healthy` im ersten Abfragetakt). Nicht widersprüchlich, aber `u-boot up` selbst ist von mir nicht nachgefahren. |
| F-2 | LOW | §1, Satz zum Repo: „…als Klon in einem Volume des Containers, der aus einem Bare-Repo im selben Volume entsteht“ — Bezug von „der“ (Container oder Klon) mehrdeutig; gemeint ist der Klon. Redaktionell. |

Keine Streudateien im Repo-Wurzelverzeichnis; nichts gepusht oder gemergt.
