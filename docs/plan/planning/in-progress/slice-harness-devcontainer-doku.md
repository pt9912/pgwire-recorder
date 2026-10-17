# Slice slice-harness-devcontainer-doku: Dokument zum Arbeiten im Dev-Container (`docs/developer/devcontainer.md`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein Abnahmeszenario und
kein Meilenstein hängt an ihm. Angelegt am 2026-10-17 auf Wunsch des Nutzers, ohne Voraussetzung
eines Schritts von welle-v1-abschluss oder welle-erster-release; er steht nach
`slice-v1-abschluss-container-doku` und vor den wellenlosen Slices vor dem Release oder
parallel dazu.

**Bezug:** kein Anspruch des Produkts — der Slice dokumentiert die Arbeitsumgebung der
Entwickler. Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Dokument `docs/developer/devcontainer.md` beschreibt, wie man mit einem
Dev-Container arbeitet, den das Werkzeug `u-boot` (`pt9912/u-boot`) erzeugt, vom Einrichten bis
zum Aufräumen, mit den Dateiinhalten, die außerhalb des Repos liegen.

Inhalt des Dokuments, je Punkt nur mit Probe am laufenden Container (§3.11):

- Voraussetzungen: Docker ≥ 24, Compose ≥ 2.20, `u-boot`.
- Einrichten: `u-boot init <name> --devcontainer --no-git`.
- Manuelle Ergänzungen **mit den Dateiinhalten**: `Dockerfile` (Docker Engine, sshd, Node, Claude
  Code), `start.sh` (`dockerd`, `authorized_keys`, `chown /workspaces`, `sshd -D`), Service `dev`
  in `compose.yaml` (`privileged`, `cpus`, `mem_limit`, Port `127.0.0.1:2222:22`, Volumes,
  Healthcheck).
- Schlüsselpaar erzeugen; Start, Stopp und Löschen (`u-boot up/down`, `docker compose`);
  Verbinden per SSH; Klon im Container.
- Claude Code starten und selbst anmelden (keine Zugangsdaten kopieren).
- Arbeiten mit `make` und den Gates im Container (eigener Docker-Daemon); gemessene Zeiten:
  `docs-check` 5 s, `gates-schnell` 201 s kalt; der volle `make gates` ist im Container nicht
  gemessen und steht als solcher da.
- Ergebnisse zurück ins Host-Repo: Branch `container/<thema>`, `git fetch` über SSH oder
  `git bundle`.
- Parallel arbeiten (mehrere Container), Grenzen und Sicherheit (privilegiert, nur 127.0.0.1, nur
  Schlüssel, Ressourcen), Fehlersuche und Aufräumen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Umbau des Repos und kein Eingriff in Makefile, Gates oder Dockerfile des Produkts — der
  Slice beschreibt eine Umgebung um das Repo, nicht das Repo (Schicht-Abgrenzung).
- Kein `.devcontainer/` im Repo: Die Dateien stehen im Dokument, weil sie zum Versuch außerhalb
  des Repos gehören und ein Verzeichnis im Repo ein Gate und eine Pflege verlangte, die niemand
  zugesagt hat (anderer Vorgang).
- Keine Aussagen zu Windows- oder macOS-Hosts — nicht geprüft, und §3.11 erlaubt nur Belegtes.
- Kein Schritt, den der Implementer am laufenden Container nicht nachgeführt hat: Er ist als
  „nicht geprüft“ markiert oder fehlt (§3.11).
- Keine Änderung der Positivliste des Stop-Hooks um `docs/developer/` — das ist eine
  Regeländerung und Entscheidung des Nutzers (§6).

## 2. Definition of Done

- [ ] `docs/developer/devcontainer.md` liegt vor und deckt die Punkte aus §1 ab; jeder Schritt ist
      am laufenden Container nachgeführt oder als „nicht geprüft“ markiert, die Dateiinhalte
      stimmen mit denen des Versuchs überein.
- [ ] `make kopf-check` und `make docs-check` grün; `make gates` grün (Stop-Hook, §6).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/developer/devcontainer.md` | neu | Liefer-Punkt: das Dokument nach §1 |

- Der Implementer führt die Schritte am laufenden Container unter `/Development/pgr-dev` nach
  (nur lesend und per SSH; den Container nicht verändern) und übernimmt die Dateiinhalte von dort.
- Zusätzlich baut der Implementer den Container in einem Scratch-Verzeichnis (anderer Name, Port 2223)
  aus den Codeblöcken des Dokuments neu auf; die Befehle und Beobachtungen stehen in §7. Das Dokument
  führt über die Liste aus §1 hinaus die Abschnitte Parallel arbeiten, Fehlersuche und Nicht geprüft.
- Das Dokument ist Ist-Zustand ohne Chronik (§3.11); es nennt keine Kennungen von Lastenheft,
  Spezifikation, ADRs, Slices oder Wellen. Nennt es doch eine ADR-Kennung, ist sie ein Link
  (`.d-check.yml`: `ids` mit `link-policy: always`).
- `.d-check.yml` schließt `docs/developer/` nicht aus (`scan.roots: ["."]`, `ignore` ohne den
  Pfad): Links und Anker im Dokument werden geprüft. `harness/README.md` braucht keinen Eintrag
  (kein neues Target, keine Source-Precedence-Zeile).

## 4. Trigger

**Start** (`next` → `in-progress`): ohne Voraussetzung. Der Anspruch ist ein reiner `git mv` auf
dem Hauptzweig vor der Arbeit.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt. Schätzung 150 bis 300,
Meldeschwelle 400, Marke 500.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 500, oder die Schritte
  verlangen mehr als dieses eine Dokument.
- `in-progress` → `open` (blockiert): Der Container unter `/Development/pgr-dev` läuft nicht
  mehr und die Schritte lassen sich nicht nachführen, ohne ihn neu zu bauen.

## 5. Closure-Trigger

DoD vollständig, `make kopf-check`, `make docs-check` und `make gates` grün, Review- und
Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Stop-Hook und Positivliste.** `docs/developer/` steht nicht auf der Positivliste des
  Stop-Hooks (`docs/plan/planning/`, `docs/reviews/`, `docs/user/`, Wurzel-`README.md`); der Hook
  verlangt für diesen Slice einen vollen `make gates` vor dem Push. Dieser läuft ohnehin mit dem
  nächsten Code-Slice; die Positivliste ändert dieser Slice nicht (Regeländerung, Entscheidung
  des Nutzers). — **Ausgang:** weiter offen: → Entscheidung des Nutzers über `docs/developer/`
  auf der Positivliste; bis dahin trägt der volle Lauf.
- **Das Dokument beschreibt einen Versuch außerhalb des Repos.** Wird der Container umgebaut,
  veraltet das Dokument still; kein Sensor prüft es gegen den Container. — **Ausgang:** weiter
  offen: → BEO-REPO/zusage-im-kommentar-weiter-als-pruefung im Register, falls der Fall bei der
  Closure auftritt; sonst entfallen mit Begründung.
- **Nicht nachführbare Schritte** (etwa der volle `make gates` im Container, Zeit nicht
  gemessen). — **Ausgang:** entfallen: Das Dokument markiert sie als „nicht geprüft“ (§1).
- **Randformen:** kein neuer Vertrag, keine Randformen im Sinn von `AGENTS.md` §3.12.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

### Belege des Implementers

**Probe-Aufbau (2026-10-17).** Host: Docker 29.9.0, Compose 5.6.0, `u-boot` 0.9.0, Git 2.43.0; im
Container Engine 29.9.1, Node.js v24.22.0, Claude Code 2.1.302, Debian 13. Der Versuchs-Container
`pgr-dev` (Port 2222) blieb unverändert; im Versuchs-Container legte ich zeitweise
`/workspaces/probe-clone`, `/tmp/c`, `/tmp/x.bundle` und `/tmp/dc.log` an und entfernte sie wieder
(`ls /workspaces` danach: `pgwire-recorder`). Neuaufbau im Scratch-Verzeichnis `…/scratchpad/pgr-dev-probe`
(Name `pgr-dev-probe`, Port 2223, `container_name`, `hostname`, `image` mit `-probe`, Pfade für Schlüssel
angepasst, sonst die Dateien aus den Codeblöcken des Dokuments; `diff` Dokument gegen Versuch:
`Dockerfile`, `start.sh`, `compose.yaml` ohne Unterschied). Vorbedingungen vor dem Aufbau: `df -h /docker`
159 GB frei, `pgrep -af "make gates"` ohne Treffer (ein früherer Treffer, `make gates-schnell` eines anderen
Laufs, war beendet). Zählung des Docker-Hosts vor und nach dem Aufräumen (`docker ps -aq`, `docker images -q`,
`docker volume ls -q`, `docker network ls -q`): vorher 1 / 316 / 50 / 4, nachher 1 / 316 / 50 / 4; der Build-Cache
des Hosts wuchs von 3,528 GB auf 6,517 GB (nicht aufgeräumt: er gehört dem gemeinsamen Host-Daemon, ein
`docker builder prune` träfe fremde Läufe; `docker system df`). `df -h /docker`: 289 GB belegt vorher, 292 GB nachher.
Aufgeräumt mit `docker compose down -v` (Volumes und Netz entfernt) und `docker rmi pgr-dev-probe:latest`;
Verzeichnis `pgr-dev-probe` gelöscht.

**Aussage · Probe · Beobachtung** (`V` = Versuchs-Container, `P` = Probe-Container; `SSH-V`/`SSH-P` =
`ssh -i … -o IdentitiesOnly=yes -p 2222|2223 vscode@127.0.0.1`).

| Aussage im Dokument | Probe | Beobachtung |
|---|---|---|
| §1 eigener Daemon, Name `pgr-dev` | `SSH-V 'docker info'` | `Name: pgr-dev`, `Server Version: 29.9.1`, `Storage Driver: overlayfs` |
| §1/§5 `docker ps -a` listet keine Host-Container, auch nicht `pgr-dev` | `SSH-V 'docker ps -a'`; Host: `docker ps -a` | Container: leere Liste; Host: `pgr-dev` |
| §1 derselbe Tag, andere Image-ID | `docker images` Host und `SSH-V` | `pgwire-recorder:dev` Host `2f13004ff698`, V `0c0b64b3febf` |
| §1 Container sieht die Netze des Hosts nicht | `docker network ls` Host und V | Host: `bridge host none pgr-dev-default`; V: `bridge host none` |
| §1/§11 privilegiert; §12 ohne `privileged` scheitert `dockerd` | `docker inspect pgr-dev` (`priv=true`); `docker run --rm pgr-dev-probe:latest` ohne `--privileged` | nach 64 s „dockerd nicht bereit“, letzte Zeile `failed to create NAT chain DOCKER: iptables failed` |
| §2 Voraussetzungen | `u-boot doctor` | 16 ok, 0 warn, 0 error; „docker 29.9.0 available (≥ 24.0)“, „docker compose 5.6.0 available (≥ 2.20)“ |
| §3.1 erzeugte Dateien, Block des `Dockerfile`, `services: {}`/`volumes: {}` | `u-boot init pgr-dev-probe --devcontainer --no-git` in leerem Verzeichnis | Ausgabe nennt `docker/ scripts/ docs/ .devcontainer/ README.md CHANGELOG.md compose.yaml .env.example .gitignore devcontainer.json Dockerfile u-boot.yaml`; kein `.git`; `compose.yaml` endet mit `services: {}` und `volumes: {}`; `u-boot.yaml`: `devcontainer.enabled: true` |
| §3.1 `u-boot add` nur `postgres`; `generate devcontainer` ohne Docker-in-Docker-Option | `u-boot add --help`, `u-boot generate devcontainer --help` | „today MVP supports only postgres“; Flags nur `--allow-external-feature-sources`, `--diff`, `--dry-run` |
| §3.2 Schlüsselpaar, Verzeichnis 700 | die Befehle des Dokuments in `P` | Exit 0, `stat -c %a ssh` = 700, `id_ed25519` 600 |
| §3.3–3.5 Dateiinhalte | `diff` der Codeblöcke gegen `/Development/pgr-dev/…` | kein Unterschied (`D-gleich`, `S-gleich`, `C-gleich`) |
| §3.4 trägt den Schlüssel als `authorized_keys` ein, 1 Schlüssel | `SSH-P 'wc -l < ~/.ssh/authorized_keys; ls -ld ~/.ssh'` | 1; `drwx------ vscode` |
| §3.4 wartet bis 60 s auf `dockerd` | `docker run` ohne `--privileged` | Meldung nach 64 s (60 s Schleife + Start) |
| §3.4 übergibt `/workspaces` an `vscode`; §12 frisches Volume gehört `root` | `ls -ld /workspaces` in P; `docker run --rm --entrypoint bash -v neu:/workspaces … 'su vscode -c "mkdir /workspaces/x"'` | P: `vscode vscode`; frisches Volume `root root`, `mkdir: Permission denied` |
| §3.4 `sshd` im Vordergrund | `ps` in V | PID 1 `sshd: /usr/sbin/sshd -D -e` |
| §3.5 `cpus`, `mem_limit` | `cat /sys/fs/cgroup/cpu.max memory.max` in V und P | `800000 100000`, `12884901888` (beide) |
| §3.5 `127.0.0.1:2222:22`; §11 nur Loopback | `docker port pgr-dev`; `ss -ltn` am Host | `22/tcp -> 127.0.0.1:2222`; `LISTEN 127.0.0.1:2222` |
| §3.5/§7/§11 vier Mounts, Host-Repo ro, Schlüssel ro, kein `docker.sock` | `docker inspect pgr-dev` (Mounts); `touch /host-repo/x` in V | genau: `/host-repo` rw=false, `/workspaces`, `authorized_key.pub` rw=false, `/var/lib/docker`; `touch`: „Read-only file system“; `mount | grep docker.sock` leer |
| §3.5 Healthcheck-Befehl | Text der `compose.yaml` (Diff); `docker compose ps` | gleich; `Up … (healthy)` in V und P |
| §4 `u-boot up` baut, endet `healthy`; Zeile `pull access denied` | `u-boot up` in P | Exit 0 nach 44 s; Ausgabe `pull access denied for pgr-dev-probe`, danach `Building`/`Built`; `dev running 2223:22 healthy` |
| §4 Exit 12 bei Zeitüberschreitung, Container läuft weiter | `u-boot down`, dann `u-boot up --timeout 5` | Exit 12, „stabilization timeout after 5s“; `docker compose ps` danach `Up … (healthy)` |
| §4 `u-boot down`: Container und Netz weg, Volumes und Image bleiben, Klon wieder da | `u-boot down`, `docker ps -a`, `docker network ls`, `docker volume ls`, `docker images`; nach `up` `ls /workspaces` | Container `pgr-dev-probe` und Netz fehlen; beide Volumes und Image da; `pgwire-recorder` da |
| §4 `docker compose down -v`; `docker rmi` | beide Befehle in P | drei Volume/Netz-Zeilen „Removed“, Exit 0; „Untagged/Deleted“; Zählung vor/nachher gleich |
| §5 SSH-Befehl, Gruppe `docker`, Benutzer `vscode` | `SSH-V 'id'` | `uid=1000(vscode) … groups=…,997(docker)` |
| §5 nur Schlüssel; anderer Schlüssel, `root` | Anmeldung mit Passwort ohne Schlüssel (V); Schlüssel von V an P; `root` mit Schlüssel an P | je `Permission denied (publickey).` |
| §6 Klon, `git config --global` Form | `git clone /host-repo …` in P; `git config --global` mit `HOME=/tmp/h`, danach `--list` | Klon HEAD `241423d4`; Liste `safe.directory`, `user.name`, `user.email` |
| §6 `safe.directory` nicht nötig | `HOME=/tmp/h git clone /host-repo /tmp/h/c` ohne Eintrag (V) | Exit 0, keine Meldung; `stat -c %U /host-repo` = `vscode` (UID 1000) |
| §6 „Author identity unknown“ | `git commit --allow-empty` ohne Identität (V-Klon) | „Author identity unknown“ |
| §6 Klon enthält committeten Stand, nicht Uncommittetes | `git status` Host: `?? docs/developer/`; frischer Klon in V | `ls /tmp/c/docs/developer`: „No such file“ |
| §6 `git pull` holt späteren Stand | Klon in V, `git reset --hard HEAD~3`, `git pull` | Fast-Forward auf `7bd8dd82` (`origin` = `/host-repo`) |
| §7 `claude --version` im Container | `SSH-V`, `SSH-P` | `2.1.302 (Claude Code)` |
| §8 Zeiten | `make docs-check`, `make gates-schnell`, `make test` im Klon von P (frischer Daemon) | 1 s Exit 0; 191 s Exit 0; 1 s Exit 0. Werte aus V (5 s, 201 s, 1 s) sind Angaben des Nutzers aus `scripts/nachweis2.sh`/`nachweis3.sh`, von mir nicht wiederholt |
| §8 Images der Gates im Daemon des Containers | `docker images` in P nach dem Lauf | `ghcr.io/pt9912/d-check`, `a-check`, `pgwire-recorder:dev/test` |
| §9 Remote nur `origin`; Variante A | `git remote -v` am Host; Zweig `container/probe` im Klon von V, `git fetch ssh://… container/probe:container/probe` in Scratch-Klon des Host-Repos mit `GIT_SSH_COMMAND` | nur `origin` GitHub; neuer Branch, `4fc4483a probe` (ich fetchte nicht ins Host-Repo, um dessen Refs nicht anzufassen) |
| §9 Variante B, fehlender Stand | `git bundle create /tmp/x.bundle main..container/probe`; `scp -P 2222 …`; `git bundle verify`; `git fetch ./x.bundle container/probe:container/probe`; Fetch in leeres Repo | verify ok; Branch `4fc4483a probe`; leeres Repo: „Dem Repository fehlen folgende vorausgesetzte Commits“ |
| §10 zwei Container gleichzeitig; Volumes mit Projektnamen | `docker compose ps` beide; `docker volume ls` | `pgr-dev` und `pgr-dev-probe` beide healthy; `pgr-dev-probe_pgr-dev-docker`, `…-workspace` |
| §10/§12 Port belegt | Port in P auf 2222, `docker compose up -d` | „Bind for 127.0.0.1:2222 failed: port is already allocated“, Exit 1; Port zurück, V unberührt |
| §11 Volumes überleben `down`, `down -v` löscht | siehe `u-boot down`, `docker compose down -v` | siehe oben |
| §11 Platz | `docker system df -v` | Image 2,05 GB; `pgr-dev_pgr-dev-docker` 5,87 GB, `…-workspace` 37 MB (V); P nach `make gates-schnell`: 5,87 GB, 27,5 MB |
| §11 `docker info` zeigt CPUs/Speicher des Hosts | `docker info` in V; `nproc`, `free -g` am Host | `CPUs: 20`, `Total Memory: 31.07GiB`; Host 20, 31 GiB |
| §12 `docker logs`, `/var/log/dockerd.log` | `docker logs pgr-dev-probe`; `docker exec … tail /var/log/dockerd.log` | `sshd`-Zeilen („Accepted publickey“); `dockerd`-Zeilen |
| §12 `No space left on device` | nicht ausgelöst | siehe „nicht geprüft“ |

**Quantoren-Durchsicht** (`grep -nE "jede|alle|immer|stets|nie|nichts|kein|jeweils|nur|einzige|weder|ohne"` über
das Dokument). Entscheidungen: „sieht weder die Container noch die Netze“ (§1) verbunden mit den Zeilen
oben; „keine Sandbox“ (§1), „ersetzt weder“ (§1) und „keine Zusage“ (§8) sind Regeln des Schreibers ohne
Verhaltensaussage, unverändert; „nur ein Linux-Host“ (§1) benennt die Probe; „kennt nur das Add-on postgres“
(§3.1) verbunden mit `u-boot add --help` (Wortlaut „today MVP“, deshalb im Dokument ohne „künftig“);
„der einzige Schlüssel“ (§3.2) verbunden mit `wc -l` = 1; „nur auf der Loopback-Adresse“ (§3.5, §11)
verbunden mit `docker port`, `ss`; „lässt nur Schlüssel zu“ (§5) verbunden mit den drei abgewiesenen
Anmeldungen (Passwort, falscher Schlüssel, `root`); „genau vier Mounts“ (§7) verbunden mit `docker inspect`;
„nur `origin`“ (§9) verbunden mit `git remote -v`; „jeden Container einzeln“ (§10) verbunden mit `cgroup`
in V und P; „kein Remote zeigt auf ihn“ (§9) eng gefasst auf das Host-Repo am Prüftag; „Docker-Socket nicht
gemountet“ (§11) verbunden mit der Mount-Liste. Ausnahmen, die die Seite für „nur Schlüssel“ schon aufzählt
(Benutzer `root`, anderer Benutzer): `AllowUsers vscode` in der Konfiguration, `root` probiert, anderer
Benutzer nicht; das Dokument sagt „Benutzer `vscode`“ ohne Zusage über weitere Benutzer.

**Zu jeder Zusage das gesehene Gegenbeispiel** (Dokument statt Test: ein Gegenfall je Aussage mit
Verneinung): `privileged` ↔ Lauf ohne `--privileged` rot; Schlüssel ↔ drei abgewiesene Anmeldungen;
Loopback ↔ `ss` zeigt nur `127.0.0.1`; `/host-repo` schreibgeschützt ↔ `touch` scheitert; Port ↔
Kollision mit 2222; `/workspaces` ↔ frisches Volume `root`, `mkdir` scheitert; Bundle-Voraussetzung ↔ Fetch in
leeres Repo scheitert; Timeout ↔ `--timeout 5` Exit 12.

**Läufe.** `make docs-check` 0 (966 Dateien), `make kopf-check` 0, `make meldungskatalog-check` 0; ein
voller `make gates` läuft nicht (Auftrag; §6 Risiko Stop-Hook). Größe nach §4: siehe unten.

**Nicht geprüft** (im Dokument §13 benannt): Start von `claude` und Anmeldung; voller `make gates` im
Container; Windows/macOS; Installation von `u-boot`; `u-boot down --volumes`; Fehler bei vollem
Docker-Speicher im Container; andere UID als 1000 am Host (`safe.directory`); neue SSH-Hostschlüssel nach
Neubau des Images (nach `u-boot down`/`up` meldete `ssh` mit gespeichertem Schlüssel keine Warnung; ein Satz dazu steht
nicht im Dokument).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Zählung nach `AGENTS.md` §3.13:** 1 Liefer-Punkt (das Dokument), 1 Schicht (Nutzer- und
Wartungs-Doku; `docs/developer/` zählt wie `docs/user/` und `docs/maintainer/` als
Dokumentations-Schicht, die Planung zählt nicht). Unter den Grenzen von drei Punkten und zwei
Schichten.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist nur die Default-Sub-Area `*` (Kürzel
`REPO`, GF) aus `harness/conventions.md`; sie erfüllt die Schwelle als Default für das gesamte Repo.

**Vorgelagert — offene Beobachtungen sichten:** Register unter `../observations/` durchgegangen
für `REPO`: der einschlägige Eintrag ist `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`
(Zusage über die Prüfung hinaus); das Dokument sagt nur Nachgeführtes zu (§1), der Zähler-Stand
ist bei der Closure aus den Dateien zu lesen.

Alle berührten Sub-Areas GF.
