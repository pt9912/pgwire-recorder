# Review-Report: slice-harness-devcontainer-doku (zweiter, unabhängiger Lauf) — 2026-10-17

**Review-Art:** Code — das Dokument `docs/developer/devcontainer.md` gegen Plan und Konventionen
(Doku-Review; Maßstab `AGENTS.md` §3.7, §3.9, §3.11), mit lesenden Proben am Versuchs-Container.

**Gegenstand:** `docs/developer/devcontainer.md` im Stand von Commit `d725a22c` (seither unverändert;
HEAD `affbd5a2`); Plan `slice-harness-devcontainer-doku`.

**Skill:** `.harness/skills/reviewer.md` @ `affbd5a2`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

> **Zitier-Form.** Kennung statt Adresse; Baseline-Stellen als Tag + Pfad in Inline-Code.

**Eingangs-Kontext:**

- Slice-Plan `slice-harness-devcontainer-doku`
- `/Development/pgr-dev` (`compose.yaml`, `.devcontainer/Dockerfile`, `.devcontainer/start.sh`, `u-boot.yaml`, `scripts/nachweis*.sh`) und der laufende Container `pgr-dev` (nur lesende Proben per SSH und `docker inspect`)
- `u-boot` 0.9.0 (`init`, `add`, `generate`, `up`, `down`, `doctor` mit `--help`; `u-boot init` in einem Scratch-Verzeichnis)
- `AGENTS.md` (Hard Rules, §3.7, §3.9, §3.11)

**Unabhängigkeit.** Die Befunde wurden vor dem Abgleich am Dokument und am Container gebildet. Ein
`grep` nach Pfadverweisen zeigte beim Lauf Zeilen des ersten Reviews (F-1 bis F-5); ich habe es danach
nicht gelesen. Der Abgleich unten stützt sich nur auf diese Zeilen.

**Gelaufen:** Codeblöcke gegen `/Development/pgr-dev` per `diff` (Dockerfile Zeilen 90 bis 118,
`start.sh`, `compose.yaml`: identisch); `u-boot init pgr-dev --devcontainer --no-git` im Scratch-Verzeichnis
(Tabelle §3.1 stimmt, `services: {}` und `volumes: {}` stimmen); `docker inspect pgr-dev` (vier Mounts,
`/host-repo` und Schlüssel `ro`, privilegiert, 8 CPUs, 12 GiB); im Container `ss -ltn`, `id`,
`sudo -n -l`, Versionen (Engine 29.9.1, Node v24.22.0, Claude Code 2.1.302), Host-Docker 29.9.0 und
Compose 5.6.0, `docker system df -v`. Keine Make-Gates, keine Builds.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Das Dokument sagt „`sshd` ist an `127.0.0.1` gebunden“. Im Container lauscht `sshd` auf `0.0.0.0:22` und `[::]:22` (`ss -ltn`); an Loopback gebunden ist nur der veröffentlichte Host-Port. Andere Container auf dem Host (Netz `pgr-dev-default`, Standard-Bridge) erreichen Port 22 des Containers über dessen Adresse; geschützt sind sie nur durch den Schlüssel. §5 und §3.5 formulieren dagegen richtig („der Port ist an `127.0.0.1` gebunden“, „erreichbar auf der Loopback-Adresse des Hosts“). | `AGENTS.md` §3.11 (Zusage nur, was die Probe belegt) | `docs/developer/devcontainer.md` · „`sshd` ist an `127.0.0.1` gebunden und lässt nur den einen“ (§11) | ja — `ss -ltn` im Container; Anmeldeversuch von einem zweiten Container aus | Sicherheitsaussage weiter als die Messung (Bindung des Prozesses statt des Ports) |
| F-2 | MEDIUM | §7 sagt „du kopierst keine Zugangsdaten vom Host“ und „Verzeichnisse mit Zugangsdaten sind nicht darunter“, §11 führt aber die Claude-Anmeldedaten nicht als Gefahr auf. Die Anmeldung im Container legt Zugangsdaten unter `/home/vscode` ab (kein Volume, `~/.claude` fehlt noch im Versuchs-Container). Jeder Prozess in der Session (Build- oder Testlauf, Paket-Skript) kann sie mit `sudo` lesen; ein Netz-Ausgang ist nicht eingeschränkt. Dasselbe gilt für Git-Name und -E-Mail aus §6. Dass `u-boot down` sie mit dem Container entfernt, steht nirgends. | `AGENTS.md` §3.11; Maintainability (Gefahr fehlt) | `docs/developer/devcontainer.md` · „Die Anmeldung geschieht im Container selbst; du kopierst keine Zugangsdaten vom Host.“ (§7); §11 „Grenzen und Sicherheit“ | nein — Anmeldung ist interaktiv, im Dokument selbst als ungeprüft geführt | Sicherheitsabschnitt nennt eine Gefahr nicht (Zugangsdaten im Container, Ablage und Lebensdauer) |
| F-3 | LOW | §4 sagt zu `u-boot down`: „Container und Netz entfallen; Volumes … bleiben, nach `u-boot up` ist der Klon wieder da“. Das Home-Verzeichnis liegt nicht im Volume: `~/.gitconfig` (§6, `git config --global`) und eine Claude-Anmeldung sind nach `down` und `up` weg, `git commit` meldet dann wieder „Author identity unknown“. Der Satz ist für den Klon wahr, schweigt zum Rest, und §6 schreibt `--global` ohne Hinweis. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „nach `u-boot up` ist der Klon wieder da“ (§4) | ja — `u-boot down`, `u-boot up`, `git config --global -l` im Container | Lebensdauer des Zustands außerhalb der Volumes nicht genannt |
| F-4 | LOW | „Klon 28 bis 37 MB“ (§11) trägt keine Quelle für die 28 und mischt Messarten: `docker system df -v` zeigt 37,14 MB für das Volume `pgr-dev_pgr-dev-workspace`, `du -sh` auf dem Klon im Versuchs-Container 47 MB (mit `.git`), 18 MB ohne. Die Spanne wächst mit dem Host-Repo. Die übrigen Größen (Image 2,05 GB, Volume 5,87 GB) stimmen mit `docker system df -v` überein. | `AGENTS.md` §3.11 (Messwerte mit Quelle) | `docs/developer/devcontainer.md` · „der Klon 28 bis 37 MB“ (§11) | ja — `docker system df -v`, `du -sh` im Container | Messwert ohne Quelle und mit gemischter Messart |
| F-5 | LOW | Die Tabelle in §8 nennt für `make docs-check` im frischen Container „1 s“ und lässt (anders als die Zeilen darunter) den Zustand des Daemon-Caches offen. Ein kalter Daemon müsste zuerst das Image `d-check` laden; 1 s ist nur mit warmem Cache denkbar. Wer die Zeit nachmisst, liest sie als kalt. | `AGENTS.md` §3.11 (Bedingung des Messwerts) | `docs/developer/devcontainer.md` · „erster Lauf im Versuchs-Container; Probe im Klon eines zweiten, frischen Containers“ (§8) | nein — Läufe nicht wiederholt (§13 deckt den Ablauf nicht) | Messbedingung unvollständig |
| F-6 | LOW | „Docker-Socket des Hosts: nicht gemountet“ steht als Schutzaussage neben „Host-Repo: `ro`“. `/host-repo` zeigt aber den ganzen Arbeitsbaum des Host-Repos lesbar (`.git`, `.claude`, `.harness`, uncommittete und ignorierte Dateien), nicht nur die committeten Stände, die §6 nennt. Bei privilegiertem Container ist „kein Socket“ keine Grenze zum Host-Daemon oder Host-Dateisystem. §11 sagt Letzteres im ersten Punkt nur als „nicht geprüft“; das Lesen des Arbeitsbaums nennt es nicht. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „**Docker-Socket des Hosts:** nicht gemountet“ (§11) | ja — `ls -a /host-repo` im Container | Schutzaussage ohne Hinweis auf Lesezugriff auf den Arbeitsbaum |
| F-7 | LOW | §2 nennt keine `u-boot`-Version; „0.9.0“ steht nur in §1 („Stand der Proben“). §3.1 setzt `--devcontainer` und `--no-git` voraus, §3.1 zu `u-boot add` („kennt nur das Add-on `postgres`“) hängt am Stand 0.9.0 (Hilfetext: „today MVP“). Weitere Voraussetzungen fehlen in §2: ein SSH-Client mit `ssh-keygen` und `scp` am Host, Schreibrecht in `/Development`, Netz für den Build, ein vorhandenes Schlüsselpaar vor dem ersten `u-boot up` (fehlt die Datei, legt Docker an ihrer Stelle ein Verzeichnis an; nicht geprüft). | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „`u-boot`, das Werkzeug, das Verzeichnis, Compose-Gerüst und Dev-Container-Gerüst erzeugt“ (§2) | nein | Voraussetzungen unvollständig (Version, Werkzeuge) |
| F-8 | INFO | Die eingebetteten Blöcke von `u-boot` tragen Kennungen wie `LH-FA-DEV-004`, `LH-FA-ADD-*`, `LH-SA-FILE-002`. Sie sehen aus wie Kennungen dieses Repos, gehören aber zu `u-boot`; der Text sagt es an der Stelle nicht. Die Blöcke sind byte-gleich zur Ausgabe von `u-boot init` und deshalb zu belassen. | Maintainability | `docs/developer/devcontainer.md` · „Non-root default user (LH-FA-DEV-004).“ (§3.3) | nein | Fremde Kennung ohne Zuordnung |
| F-9 | INFO | Das Image holt `mcr.microsoft.com/devcontainers/base:debian` ohne Digest und führt `curl … setup_lts.x \| bash -` als root aus. Das Dokument nennt die fehlende Versionsfestlegung (§3.3), nicht die Herkunft der Skripte im Build eines privilegierten Containers. Keine erwartete Aktion. | Maintainability | `docs/developer/devcontainer.md` · „ohne feste Version“ (§3.3) | nein | Lieferkette des Builds nicht benannt |

## Abgleich mit dem ersten Review (nur Klassen aus den Zeilen, die der `grep` zeigte)

| Klasse im ersten Lauf | Stand im Dokument |
|---|---|
| Schutzzusage trotz privilegiertem Container mit `sudo` (F-1) | behoben: „nicht ausgeführt, nur gefolgert“, „keine Zusage der Isolation“ (§9, §11); Rest siehe F-6 |
| Messwert ohne Quelle im Dokument (F-2) | behoben: Skripte benannt; Rest siehe F-4, F-5 |
| „Nicht geprüft“ unvollständig (F-3) | behoben: Fetch in das echte Host-Repo und Hostschlüssel-Frage stehen in §13 |
| Zusage aus Hilfetext (F-4) | behoben: §4 nennt `docker compose up -d` als ungeprüft |
| Liste liest sich vollständig (F-5) | behoben: Tabelle §3.1 trägt `docker/`, `scripts/`, `docs/` |

Neue Fehler der Nacharbeit: F-1 (Wortlaut in §11), F-4 (neue Größenspanne), F-5 (Tabellenzeile); F-2, F-3, F-6, F-7 sind Lücken, die das Dokument von Anfang an hatte.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Codeblöcke Dockerfile, `start.sh`, `compose.yaml` gegen `/Development/pgr-dev` (`diff`) | geprüft, ohne Befund (identisch) |
| §3.1 Tabelle „Was `u-boot` erzeugt“ gegen `u-boot init` im Scratch-Verzeichnis | geprüft, ohne Befund |
| Reihenfolge §3.2 (Schlüssel) vor §3.5 (Mount) vor §4 (Start); Pfade `/etc/pgr-dev/authorized_key.pub`, `/workspaces`, Modi 700/600 | geprüft, ohne Befund (kein Widerspruch zwischen §3.3, §3.4, §3.5) |
| Versionen (Docker 29.9.0, Compose 5.6.0, `u-boot` 0.9.0, Engine 29.9.1, Node v24.22.0, Claude Code 2.1.302) | geprüft, ohne Befund |
| Mounts („genau vier“), `privileged`, `sudo -n -l`, Gruppe `docker`, `cpus`, `mem_limit` | geprüft, ohne Befund |
| Host-Schlüssel des Containers (liegen im Image, überleben `down`/`up`; §13 „nach einem Neubau des Images“) | geprüft, ohne Befund |
| Image-IDs `pgwire-recorder:dev` verschieden (Host `2f13004ff698`, Container `0c0b64b3febf`) | geprüft, ohne Befund |
| Größen Image 2,05 GB und Volume 5,87 GB | geprüft, ohne Befund |
| `u-boot up`/`down`-Semantik (Timeout 60 s, Exit 12, `--volumes`) gegen Hilfetexte | geprüft, ohne Befund (die gemessenen 44 s sind aus den Skripten nicht nachvollziehbar, im Dokument als Probe geführt) |
| §13 „Nicht geprüft“: Vollständigkeit | geprüft; Ergänzungen siehe F-2, F-3, F-5 |
| Anker, Tabellen, interne Querverweise (§-Verweise) | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 5 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Sicherheitsaussage weiter als die Messung · Sicherheitsabschnitt nennt eine Gefahr nicht · Lebensdauer des Zustands außerhalb der Volumes nicht genannt · Messwert ohne Quelle oder mit unvollständiger Bedingung · Voraussetzungen unvollständig · Fremde Kennung ohne Zuordnung

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 (MEDIUM) sind Sicherheitsaussagen und blockieren. Die LOW- und INFO-Punkte blockieren nicht.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg. DoD- und Spec-Konformität prüft der Verifier separat.

**Zuständigkeit:** Dokumenttext und Slice-Plan §7: Implementer. Ob F-2 zusätzlich Schritte für den Betrieb verlangt (Zugangsdaten, Netz-Ausgang), ist eine Entscheidung des Architect bzw. des Nutzers; F-6 und F-9 sind Hinweise ohne Aktion an den Reviewer-Skill.
