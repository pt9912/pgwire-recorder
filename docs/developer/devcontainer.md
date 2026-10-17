# Arbeiten im Dev-Container

Version: 0.2
Stand: 2026-10-17

## 1. Wofür der Dev-Container gut ist

Der Dev-Container ist ein privilegierter Docker-Container mit einem **eigenen Docker-Daemon**
(`dockerd`), `sshd`, Node.js und Claude Code. Das Repo liegt als Klon in einem Volume des
Containers, der aus einem Bare-Repo im selben Volume entsteht; der Host liefert mit `git push` ein und
holt mit `git fetch` zurück (§6, §9). Der Arbeitsbaum des Hosts ist nicht eingebunden. Die Images, Tags, Netze und Container, die `make` im Container anlegt, liegen im
Daemon des Containers, nicht im Daemon des Hosts: Derselbe Tag (`pgwire-recorder:dev`) hat dort eine
andere Image-ID als am Host, und der Container sieht weder die Container noch die Netze des Hosts.
Claude Code läuft darin als eigene Session, getrennt von der Session am Host.

Wofür er nicht gut ist:

- Er ist **keine Sandbox**. Der Container ist privilegiert, und `vscode` darf mit `sudo` ohne Passwort
  alles als `root` ausführen (§11); behandle ihn nicht als Schutz des Hosts vor fremdem Code.
- Er ersetzt weder den vollen Gate-Lauf am Host noch einen Release. Der volle `make gates` ist im
  Container **nicht gemessen** (§8).
- Er ist **nicht Teil des Repos**: Die Dateien, die ihn bauen, stehen in diesem Dokument und liegen
  in einem Verzeichnis außerhalb des Repos (§3). Wird dort etwas geändert, folgt dieses Dokument
  nicht von selbst.

Geprüft ist nur ein Linux-Host mit Docker 29.9.0, Compose 5.6.0 und `u-boot` 0.9.0 (Stand der
Proben). Windows- und macOS-Hosts sind nicht geprüft.

## 2. Voraussetzungen

- Docker Engine ≥ 24 und Docker Compose ≥ 2.20 (das fordert `u-boot doctor`, §12).
- `u-boot`, das Werkzeug, das Verzeichnis, Compose-Gerüst und Dev-Container-Gerüst erzeugt
  (<https://github.com/pt9912/u-boot>), in der Version 0.9.0, mit der die Proben liefen (`u-boot
  --version`). Die Optionen `--devcontainer` und `--no-git` aus §3.1 und die Aussage zu `u-boot add`
  gelten für diese Version; andere Versionen sind nicht geprüft. Wie es zu installieren ist, steht
  dort; dieses Dokument prüft es nicht.
- Git.
- Am Host ein SSH-Client mit `ssh-keygen` (§3.2, §5, §9).
- Schreibrecht in dem Verzeichnis, in dem `pgr-dev` angelegt wird.
- Netz für den Build des Images (§3.3, §4).
- Das Schlüsselpaar aus §3.2 liegt vor dem ersten `u-boot up` an dem Pfad, den `compose.yaml` nennt.
  Was Docker bei einer fehlenden Datei am Mount-Pfad tut, ist nicht geprüft.
- Freier Platz im Docker-Speicher des Hosts: Der Container baut ein Image und legt den Docker-Cache
  des Repos in einem Volume ab (§11).

## 3. Einrichten

Das Verzeichnis liegt **außerhalb** des Repos. Im Folgenden heißt es `pgr-dev`. Ein Pfad ist
Sache des Hosts und steht als Variable in den Befehlen:

| Variable | Bedeutung |
|---|---|
| `PGR_DEV` | absoluter Pfad des Verzeichnisses `pgr-dev` am Host (wird in §3.1 angelegt) |

```bash
export PGR_DEV=<Pfad, unter dem pgr-dev angelegt wird>/pgr-dev
```

`compose.yaml` (§3.5) enthält `PGR_DEV` nicht: Der Pfad des Schlüssels ist relativ zum
Verzeichnis der Datei. Das Host-Repo kommt in `compose.yaml` nicht vor (§6, §9).

### 3.1 Gerüst mit `u-boot`

```bash
mkdir "$PGR_DEV" && cd "$PGR_DEV"
u-boot init pgr-dev --devcontainer --no-git
```

`--devcontainer` erzeugt zusätzlich `.devcontainer/devcontainer.json` und
`.devcontainer/Dockerfile` und setzt `devcontainer.enabled` in `u-boot.yaml`; `--no-git` lässt das
Verzeichnis ohne Git-Repo.

Was `u-boot` erzeugt und was von Hand dazukommt:

| Datei | `u-boot init` | von Hand |
|---|---|---|
| `u-boot.yaml`, `.env.example`, `.gitignore`, `README.md`, `CHANGELOG.md` | erzeugt | — |
| `.devcontainer/devcontainer.json` | erzeugt | — |
| `.devcontainer/Dockerfile` | der Block `FROM … USER vscode` | alles danach (§3.3) |
| `.devcontainer/start.sh` | — | ganz (§3.4) |
| `compose.yaml` | der Kopf mit Projektname und Netz | der Service `dev` und das Volume-Verzeichnis (§3.5) |
| `ssh/` | — | Schlüsselpaar (§3.2) |
| `docker/`, `scripts/`, `docs/` | erzeugt, leer | — |

Die Blöcke zwischen `BEGIN U-BOOT MANAGED BLOCK` und `END U-BOOT MANAGED BLOCK` gehören
`u-boot`; die manuellen Teile stehen außerhalb davon. Die Kennungen, die in diesen Blöcken stehen
(sie beginnen mit `LH-`), sind Kennungen von `u-boot`, nicht dieses Repos; die Blöcke sind wörtlich
übernommen. `u-boot add` kennt nur das Add-on
`postgres`, und `u-boot generate devcontainer` hat keine Option für Docker-in-Docker; der Dev-Container
wird deshalb von Hand ergänzt.

### 3.2 Schlüsselpaar

```bash
mkdir -p ssh && chmod 700 ssh
ssh-keygen -t ed25519 -N '' -C pgr-dev -f ssh/id_ed25519
```

Der öffentliche Teil `ssh/id_ed25519.pub` wird als Mount `ro` in den Container eingebunden und ist
dort der einzige Schlüssel, mit dem man sich anmelden kann (§5). Der private Teil bleibt im
Verzeichnis `ssh/` auf dem Host. `-N ''` legt den Schlüssel **ohne Passphrase** an: Wer die Datei
`ssh/id_ed25519` lesen kann, kann sich am Container anmelden, und dort hat `vscode` `sudo` ohne
Passwort (§11). Schütze das Verzeichnis (Modus 700) oder setze eine Passphrase.

### 3.3 `.devcontainer/Dockerfile`

Der erste Block stammt von `u-boot`; alles ab `USER root` kommt von Hand dazu.

```dockerfile
# BEGIN U-BOOT MANAGED BLOCK: init
# Devcontainer image for pgr-dev — generated by `u-boot generate devcontainer`.
FROM mcr.microsoft.com/devcontainers/base:debian

# Non-root default user (LH-FA-DEV-004).
USER vscode
# END U-BOOT MANAGED BLOCK: init

# --- manuell ergänzt (außerhalb des u-boot-Blocks): Docker-in-Docker, sshd, Node, Claude Code ---
USER root
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates curl gnupg make git bash openssh-server iptables iproute2 procps \
 && install -m 0755 -d /etc/apt/keyrings \
 && curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc \
 && . /etc/os-release \
 && echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian ${VERSION_CODENAME} stable" > /etc/apt/sources.list.d/docker.list \
 && curl -fsSL https://deb.nodesource.com/setup_lts.x | bash - \
 && apt-get install -y --no-install-recommends docker-ce docker-ce-cli containerd.io docker-compose-plugin docker-buildx-plugin nodejs \
 && npm install -g @anthropic-ai/claude-code \
 && usermod -aG docker vscode \
 && mkdir -p /run/sshd /etc/pgr-dev \
 && printf '%s\n' 'PasswordAuthentication no' 'KbdInteractiveAuthentication no' 'PermitRootLogin no' 'AllowUsers vscode' 'PubkeyAuthentication yes' > /etc/ssh/sshd_config.d/pgr-dev.conf \
 && rm -rf /var/lib/apt/lists/*
COPY start.sh /usr/local/bin/start.sh
RUN chmod 755 /usr/local/bin/start.sh
EXPOSE 22
# Entrypoint läuft als root (dockerd, sshd); Anmeldung und Arbeit als vscode.
ENTRYPOINT ["/usr/local/bin/start.sh"]
```

Das Image installiert Docker Engine, Node.js (die LTS-Reihe von NodeSource) und Claude Code ohne feste
Version. Der Versuch und die Probe hatten Engine 29.9.1, Node.js v24.22.0 und Claude Code 2.1.302; ein
späterer Build kann andere Versionen holen. Der Build braucht Netz. Das Basis-Image
`mcr.microsoft.com/devcontainers/base:debian` ist nur mit einem Tag angegeben, nicht mit einem Digest,
und der Build führt das Skript `setup_lts.x` von NodeSource per `curl … | bash -` als `root` aus und
holt Docker-Schlüssel und Pakete aus dem Netz. Die Herkunft dieser Quellen ist nicht geprüft.

### 3.4 `.devcontainer/start.sh`

```bash
#!/bin/bash
set -euo pipefail
dockerd >/var/log/dockerd.log 2>&1 &
for i in $(seq 1 60); do docker info >/dev/null 2>&1 && break; sleep 1; done
docker info >/dev/null 2>&1 || { echo "dockerd nicht bereit"; tail -n 30 /var/log/dockerd.log; exit 1; }
install -d -m 700 -o vscode -g vscode /home/vscode/.ssh
install -m 600 -o vscode -g vscode /etc/pgr-dev/authorized_key.pub /home/vscode/.ssh/authorized_keys
chown vscode:vscode /workspaces
ssh-keygen -A
exec /usr/sbin/sshd -D -e
```

Das Skript startet `dockerd`, wartet bis zu 60 Sekunden auf ihn, trägt den Schlüssel aus §3.2 als
`authorized_keys` des Benutzers `vscode` ein, übergibt `/workspaces` an `vscode` und startet
`sshd` im Vordergrund. Es muss ausführbar im Image liegen; das Dockerfile setzt dafür
`chmod 755`.

### 3.5 `compose.yaml`

`u-boot init` schreibt den Block am Anfang und danach `services: {}` und `volumes: {}`. Beide
Zeilen werden durch den folgenden Teil ersetzt:

```yaml
# BEGIN U-BOOT MANAGED BLOCK: init
# Compose stack for pgr-dev. Service add-ons are appended by
# `u-boot add <service>` (LH-FA-ADD-*); manual edits outside the
# u-boot-managed blocks remain untouched (LH-SA-FILE-002).

name: pgr-dev

networks:
  default:
    name: pgr-dev-default
# END U-BOOT MANAGED BLOCK: init

services:
  dev:
    build:
      context: .devcontainer
      dockerfile: Dockerfile
    image: ${PGR_NAME:-pgr-dev}:latest
    container_name: ${PGR_NAME:-pgr-dev}
    hostname: ${PGR_NAME:-pgr-dev}
    privileged: true
    cpus: ${PGR_CPUS:-8}
    mem_limit: ${PGR_MEM:-12g}
    ports:
      - "127.0.0.1:${PGR_SSH_PORT:-2222}:22"
    volumes:
      - pgr-dev-docker:/var/lib/docker
      - pgr-dev-workspace:/workspaces
      - ./ssh/id_ed25519.pub:/etc/pgr-dev/authorized_key.pub:ro
    healthcheck:
      test: ["CMD-SHELL", "docker info >/dev/null 2>&1 && ss -ltn | grep -q ':22 '"]
      interval: 10s
      timeout: 10s
      retries: 12
      start_period: 20s

volumes:
  pgr-dev-docker:
  pgr-dev-workspace:
```

Der Block weicht an einer Stelle von der Datei ab, die `u-boot` mit `init` erzeugt: Der Pfad des
Schlüssels ist relativ (`./ssh/…`). Mit dem Block in dieser Form (`Dockerfile` und `start.sh` aus §3.3 und
§3.4, jeweils aus diesem Dokument übernommen) lief ein Scratch-Container mit eigenem Projektnamen, Namen
und Port aus einer `.env` (§10): `docker inspect` zählt drei Mounts, und `u-boot up` endet `healthy`.
Der Kopf oben trägt `name: pgr-dev` und das Netz `pgr-dev-default`; der Scratch-Container lief mit dem
Kopf, den `u-boot init pgr-aus` schrieb (`name: pgr-aus`, Netz `pgr-aus-default`). Dass `u-boot up` eine
Variable aus der Umgebung weiterreicht, ist nicht geprüft (§13).

Die Bestandteile:

- `privileged: true`: Der Container darf `dockerd` mit eigenem Speicher-Treiber und eigenen Netzen
  starten.
- `cpus` und `mem_limit`: Obergrenze des Containers (Vorgabe 8 CPUs, 12 GiB).
- `127.0.0.1:${PGR_SSH_PORT:-2222}:22` (Vorgabe `2222`): `sshd` ist nur auf der Loopback-Adresse des Hosts erreichbar.
- `pgr-dev-docker` (Compose setzt den Projektnamen davor) auf `/var/lib/docker`: Images und Cache des Daemons im Container.
- `pgr-dev-workspace` (ebenso) auf `/workspaces`: das Bare-Repo und der Klon (§6).
- `./ssh/id_ed25519.pub`: der Pfad ist relativ zum Verzeichnis von `compose.yaml`.
- Der Healthcheck ist grün, wenn der Daemon im Container antwortet und `sshd` auf Port 22 lauscht.

## 4. Start, Stopp, Löschen

Aus dem Verzeichnis `pgr-dev`:

```bash
u-boot up        # baut das Image beim ersten Mal und wartet, bis der Container healthy ist
docker compose ps
```

Die Zeile `pull access denied for pgr-dev…` beim ersten Start von `u-boot up` ist der Versuch, das
Image aus einer Registry zu holen; danach baut Compose es. Laut Hilfetext von `u-boot` ruft `u-boot up`
`docker compose up -d` auf; `docker compose up -d` allein lief nur für einen vorhandenen Container, den eine geänderte Variable aus der Umgebung (`PGR_SSH_PORT`) neu erzeugte (danach `healthy`), nicht als erster Start. Der Build lädt Pakete aus dem Netz. Die Probe baute mit warmem Layer-Cache
des Hosts, `u-boot up` endete nach 44 Sekunden mit `healthy`; die Zeit eines Builds ohne Cache ist nicht
gemessen. `u-boot up` wartet standardmäßig 60 Sekunden auf den Zustand `healthy` und endet bei
Überschreitung mit Exit 12 (`--timeout <Sekunden>` ändert das, gemessen mit `--timeout 5`); der
Container läuft danach weiter. `u-boot up` meldet zusätzlich `warn: Service "dev" port entry #0 (127.0.0.1:${PGR_SSH_PORT:-2222}:22)
is not TCP-probable; stabilization will rely on healthcheck or running-only`; der Container wurde danach `healthy`.

| Ziel | Befehl | Folge |
|---|---|---|
| Stoppen | `u-boot down` | Container und Netz entfallen; Volumes (Bare-Repo, Klon, Docker-Cache) und Image bleiben; nach `u-boot up` sind Bare-Repo, Branches und Klon wieder da, und der Host-Schlüssel von `sshd` ist derselbe Fingerabdruck wie vorher (gemessen, `git fetch` ohne Warnung). `/home/vscode` liegt im Dateisystem des Containers, in keinem der drei Mounts (§7): `~/.gitconfig` (§6) ist danach weg (gemessen); die Anmeldung von Claude Code (§7) entfällt damit ebenfalls (aus den Mounts gefolgert, nicht ausgeführt) |
| Alles löschen | `docker compose down -v` | zusätzlich die beiden Volumes (`Removed`): das Bare-Repo, der Klon und der Docker-Cache im Container sind weg; Stände, die nur dort lagen, auch |
| Image löschen | `docker rmi ${PGR_NAME:-pgr-dev}:latest` | danach baut der nächste Start neu |

## 5. Verbinden per SSH

```bash
ssh -i "$PGR_DEV/ssh/id_ed25519" -o IdentitiesOnly=yes -p "${PGR_SSH_PORT:-2222}" vscode@127.0.0.1
```

Benutzer ist `vscode` (Mitglied der Gruppe `docker`). Der Host-Port ist an `127.0.0.1` veröffentlicht
(`docker port` in der Probe mit Vorgabe: `127.0.0.1:2222`); von anderen Rechnern ist der Container darüber nicht erreichbar.
`sshd` selbst lauscht im Container auf allen Adressen (`ss -ltn`: `0.0.0.0:22` und `[::]:22`). `sshd` lässt nur Schlüssel zu
(`PasswordAuthentication no`, `KbdInteractiveAuthentication no`, `PermitRootLogin no`,
`AllowUsers vscode`); eine Anmeldung ohne Schlüssel endet mit `Permission denied (publickey)`.

Hostschlüssel: Ohne Eintrag in `~/.ssh/known_hosts` und ohne Terminal endet die Anmeldung mit
`Host key verification failed.` (Exit 255; `git fetch` mit Exit 128). Die Frage der
Erstverbindung am Terminal ist nicht geprüft. Der Fingerabdruck des Hostschlüssels nennt
`root@buildkitsandbox`, stammt also aus dem Build des Images: `u-boot down` und `u-boot up` ändern ihn
nicht, ein Neubau des Images ohne Cache (`docker compose build --no-cache`) ergab einen anderen Fingerabdruck. Dann meldet `ssh` bei
vorhandenem Eintrag für `[127.0.0.1]:<Port>` `WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!` und
`Host key verification failed.`, auch `git push` und `git fetch` enden ohne Verbindung. Das Dokument
empfiehlt dazu keine Einstellung (§13).

Im Container läuft ein eigener Daemon: `docker info` nennt den Namen `pgr-dev`, und `docker ps -a`
listet nicht die Container des Hosts, auch nicht `pgr-dev` selbst.

## 6. Bare-Repo und Klon im Container

Im Volume `pgr-dev-workspace` liegt ein Bare-Repo `/workspaces/austausch.git`; der Arbeitsklon
`/workspaces/pgwire-recorder` entsteht daraus. Die Reihenfolge ist: Bare-Repo anlegen, vom Host einliefern
(§9), im Container klonen.

Bare-Repo anlegen (im Container, nach §5 per `ssh` angemeldet, oder dem Befehl hinten angehängt):

```bash
git init --bare --initial-branch=main /workspaces/austausch.git
```

Nach dem Einliefern vom Host (`git push container main`, §9):

```bash
git clone /workspaces/austausch.git /workspaces/pgwire-recorder
cd /workspaces/pgwire-recorder
git config --global user.name "Dein Name"
git config --global user.email "du@example.org"
git switch -c thema
```

Gemessen am frischen Volume: `/workspaces` gehört `vscode` (UID 1000, `start.sh` übergibt es), das Bare-Repo
gehört `vscode` mit Modus 775, `git` und `git-receive-pack` liegen im Image (`git version 2.55.0`).
`git` läuft im Container als `vscode`, auch bei einem Push vom Host über `ssh`; der Host liest das
Bare-Repo nicht selbst, eine andere UID des Host-Benutzers (`dubious ownership`) betrifft diesen Weg
deshalb nicht (nicht mit einer anderen UID geprüft, §13).

Ohne `--initial-branch=main` zeigt `HEAD` des Bare-Repos auf `master`; nach einem Push von `main` meldet
`git clone` dann `warning: remote HEAD refers to nonexistent ref, unable to checkout` und checkt keine Dateien aus.
Ein Klon vor dem ersten Push endet mit `warning: You appear to have cloned an empty repository.`
(Exit 0, Branch `main`, keine Commits); ein zweiter `git init --bare` auf dem vorhandenen Verzeichnis
meldet `Reinitialized existing Git repository`, ein zweiter `git clone` in ein vorhandenes, nicht leeres Verzeichnis
bricht mit `fatal: destination path … already exists and is not an empty directory` ab (Exit 128), und ein
zweites `git remote add container …` am Host mit `externes Repository container existiert bereits` (Exit 3).

Name und E-Mail setzt du selbst; ohne sie weist `git commit` mit „Author identity unknown“ ab.
`--global` schreibt `~/.gitconfig` im Dateisystem des Containers, nicht in ein Volume: Nach `u-boot down`
und `u-boot up` ist die Datei weg (`git config --global -l`: `unable to read config file
'/home/vscode/.gitconfig'`), und du setzt beide Werte neu.

Im Container arbeitest du auf einem Branch, den du nicht `main` nennst, und pushst ihn in das Bare-Repo
(`origin` des Klons):

```bash
git push origin thema
```

Ein Stand des Hosts kommt mit `git switch main && git pull` im Klon.

## 7. Claude Code im Container

```bash
cd /workspaces/pgwire-recorder
claude
```

Die Anmeldung geschieht im Container selbst; du kopierst keine Zugangsdaten vom Host. Die Anmeldedaten
liegen danach im Container, nicht im Host-Repo und in keinem Volume (wo `claude` sie ablegt, ist nicht
geprüft). Jeder Prozess im Container kann sie mit `sudo` lesen, und nichts in
`compose.yaml` beschränkt den Netz-Ausgang (§11). Der Container
hat genau drei Mounts: den öffentlichen Schlüssel (Mount `ro`)
und die beiden Volumes aus §3.5 (`docker inspect` am Scratch-Container); Verzeichnisse mit Zugangsdaten
und der Arbeitsbaum des Hosts sind nicht darunter. `claude --version` läuft im Container. Der Start von `claude` und
die Anmeldung sind **nicht geprüft**, weil sie interaktiv sind.

## 8. Arbeiten im Container

Im Klon laufen die Targets wie am Host; sie benutzen den Docker-Daemon des Containers:

```bash
make docs-check
make gates-schnell
make test
```

Gemessene Zeiten im Container (Stand 2026-10-17, je Zeile die Läufe, die es gab). Die jeweils erste
Zahl stammt aus einem Skript des Nutzers am Versuchs-Container (`scripts/nachweis2.sh`,
`scripts/nachweis3.sh`); sie ist nicht wiederholt. Die zweite stammt aus der Probe in einem frischen
Container:

| Target | Zeit | Bedingung |
|---|---|---|
| `make docs-check` | 5 s; 1 s | erster Lauf im Versuchs-Container; Probe im Klon eines zweiten, frischen Containers. Der Zustand des Daemon-Caches (Image `d-check` geladen oder nicht) ist bei beiden nicht festgehalten; 1 s ist nur mit geladenem Image denkbar, ein kalter Daemon lädt es zuerst (nicht gemessen) |
| `make gates-schnell` | 201 s; 191 s | kalter Cache des Daemons: Versuchs-Container; frischer Container der Probe, Exit 0 |
| `make test` | 1 s; 1 s | warmer Cache: Versuchs-Container; derselbe frische Container nach `make gates-schnell` |

**Der volle `make gates` ist im Container nicht gemessen** und wird hier nicht zugesagt. Die Zahlen
hängen an Host, Cache und Last und sind keine Zusage.

Der Daemon im Container lädt die Images der Gates (zum Beispiel `d-check` und `a-check`) selbst und
legt sie im Volume `pgr-dev-docker` ab; die Images des Host-Daemons stehen ihm nicht zur Verfügung.

## 9. Austausch mit dem Host

Der Host schreibt in das Bare-Repo aus §6 und liest aus ihm; der Container sieht den Arbeitsbaum des
Hosts nicht. Die Befehle laufen im Host-Repo, `PGR_DEV` ist der Pfad aus §3:

```bash
git remote add container ssh://vscode@127.0.0.1:${PGR_SSH_PORT:-2222}/workspaces/austausch.git
git config core.sshCommand "ssh -i $PGR_DEV/ssh/id_ed25519 -o IdentitiesOnly=yes"
```

Stand einliefern (das Bare-Repo muss vorher bestehen, §6):

```bash
git push container main
```

Das liefert die **committeten** Stände von `main` ein, auch unveröffentlichte; uncommittete und
ignorierte Dateien kommen nicht mit. Im Container arbeitest du auf einem Branch mit eigenem Namen (§6) und
holst ihn am Host:

```bash
git fetch container
git branch -r
```

Der Branch `thema` des Containers liegt am Host als `container/thema`. Heißt er im Container schon
`container/thema`, liegt er am Host als `container/container/thema`. Er wird dort wie jeder andere
Branch geprüft und gemergt. Der Inhalt kommt aus einem Container, in dem `vscode` `root` werden kann:
Sieh die Commits durch, bevor du mergst.

Schreiben beide Seiten auf `main`, weist der zweite Push ab (`! [rejected] main -> main (fetch first)`).
Dann `git fetch container`, `container/main` in den eigenen `main` mergen und erneut pushen; im Container
holt `git pull` den Stand. Deshalb arbeitet der Container nicht auf `main`.

`core.sshCommand` gilt für das ganze Host-Repo, nicht für ein Remote: Ein `git ls-remote` auf eine
`ssh://`-URL, die kein Remote ist, rief denselben Befehl mit demselben Schlüssel auf (`origin` über `ssh` ist
nicht gelaufen). Für mehrere Container gilt §10.

Gelaufen ist das in einem Scratch-Klon des Host-Repos, nicht im echten Host-Repo.

## 10. Mehrere Container nebeneinander

Jede Instanz hat ein eigenes Verzeichnis (eigener Ordnername) mit eigenem `u-boot init`, eigenem
Schlüsselpaar und einer eigenen `.env`. Vier Werte unterscheiden die Instanzen; `compose.yaml` (§3.5) liest sie
aus der `.env` und nennt Vorgaben:

| Variable | Wirkung | Vorgabe |
|---|---|---|
| `COMPOSE_PROJECT_NAME` | Projektname; er präfixiert die Volumes (`<name>_pgr-dev-docker`, `<name>_pgr-dev-workspace`) | `pgr-dev` (`name:` im Kopf der Datei) |
| `PGR_NAME` | `container_name`, `hostname` und Tag des Images | `pgr-dev` |
| `PGR_SSH_PORT` | Host-Port von `sshd` | `2222` |
| `PGR_CPUS`, `PGR_MEM` | `cpus` und `mem_limit` | `8`, `12g` |

Zweite Instanz (Name `pgr-zwei`, Port 2223):

```bash
export PGR_DEV=<Pfad>/pgr-zwei
mkdir "$PGR_DEV" && cd "$PGR_DEV"
u-boot init pgr-zwei --devcontainer --no-git
# Schlüssel (§3.2), .devcontainer/ und compose.yaml (§3.3 bis §3.5) wie für die erste Instanz
cat > .env <<EOF
COMPOSE_PROJECT_NAME=pgr-zwei
PGR_NAME=pgr-zwei
PGR_SSH_PORT=2223
EOF
u-boot up
ssh -i "$PGR_DEV/ssh/id_ed25519" -o IdentitiesOnly=yes -p 2223 vscode@127.0.0.1
```

`u-boot init pgr-zwei` schreibt `name: pgr-zwei` in den Kopf; `COMPOSE_PROJECT_NAME` in der `.env`
setzt denselben Namen. Gemessen (`docker compose config`): mit dieser `.env` ergeben sich
Projekt `pgr-zwei`, Container, Hostname und Image `pgr-zwei`, Volumes `pgr-zwei_pgr-dev-docker` und
`pgr-zwei_pgr-dev-workspace`, Port `2223`, 4 CPUs und 6 GiB bei `PGR_CPUS=4` und `PGR_MEM=6g`; ohne `.env`
die Vorgaben. `COMPOSE_PROJECT_NAME` in der Umgebung ersetzt `name:` im Kopf. Das Netz nennt der Kopf:
`u-boot init pgr-aus` schrieb `pgr-aus-default`, der Kopf in §3.5 trägt `pgr-dev-default`. Wer ihn unverändert in
eine zweite Instanz übernimmt, hängt beide am selben Netz; ob das Starten oder `u-boot down` der einen die
andere stört, ist nicht geprüft.

Gelaufen sind zwei Container gleichzeitig (Ports 2322 und 2323, Projekte `pgr-aus` und `pgr-aus2`), jeder in
einem eigenen Verzeichnis, beide `healthy`, die Werte über die `.env` wie oben (`COMPOSE_PROJECT_NAME`,
`PGR_NAME`, `PGR_SSH_PORT`); das zweite Verzeichnis war eine Kopie des ersten mit eigenem Schlüsselpaar, ohne
`u-boot init`, mit angepasstem Kopf (`name:`, Netz). Ein belegter
Port lässt den Start scheitern (§12). Jeder Container hat seinen eigenen Docker-Cache im eigenen Volume und
braucht den Platz dafür noch einmal (§11). `cpus` und `mem_limit` begrenzen jeden Container einzeln.

Austausch mit zwei Containern: Jeder Container hat sein eigenes Bare-Repo (§6) und am Host ein Remote
mit eigenem Namen und eigenem Port. `core.sshCommand` (§9) kennt einen Befehl je Host-Repo; mit einem Schlüssel
je Instanz nennt er beide:

```bash
git remote add container2 ssh://vscode@127.0.0.1:2223/workspaces/austausch.git
git config core.sshCommand "ssh -i <PGR_DEV der ersten>/ssh/id_ed25519 -i $PGR_DEV/ssh/id_ed25519 -o IdentitiesOnly=yes"
```

Mit dem Schlüssel der ersten Instanz allein wies der zweite Container ab (`git push container2 main`
endete mit „Bitte stellen Sie sicher, dass die korrekten Zugriffsberechtigungen bestehen“); mit beiden Schlüsseln
liefen `git push container2 main`, `git fetch container` und `git fetch container2`, und `git branch -r` führte
`container/…` und `container2/…` nebeneinander. Die Befehle der übrigen Abschnitte (`ssh`, `git push`, `git fetch`)
nehmen `-p "$PGR_SSH_PORT"` und `PGR_DEV` der Instanz; `docker compose` und `u-boot` laufen im Verzeichnis der Instanz.

## 11. Grenzen und Sicherheit

- **Privilegiert, `sudo` ohne Passwort:** `privileged: true` gibt dem Container weitgehenden Zugriff
  auf den Kernel des Hosts, und `sudo -n -l` zeigt für `vscode` `(root) NOPASSWD: ALL`; dazu ist
  `vscode` in der Gruppe `docker`. Was ein Prozess darin am Host ändern könnte, ist nicht geprüft.
  Behandle den Container deshalb nicht als Sandbox: Lass darin nur Arbeit laufen, der du traust, und
  prüfe Ergebnisse, die du zurückholst (§9), als Branch, bevor du sie mergst. `vscode` gehört das
  Bare-Repo und darf dort jeden Stand ablegen (gelaufen: Push von `main` und von Branchen); was du am Host
  übernimmst, entscheidest du beim Merge.
- **Schlüssel, Loopback am Host-Port:** Veröffentlicht ist nur der Host-Port, und der ist an
  `127.0.0.1` gebunden. Im Container lauscht `sshd` auf `0.0.0.0:22` und `[::]:22`; andere Container auf
  dem Host (das Netz `pgr-dev-default` ist ein Bridge-Netz) können Port 22 des Containers über dessen
  Adresse erreichen (aus der Bindung gefolgert, mit einem zweiten Container nicht ausgeführt). Geschützt
  ist er dort nur durch den einen öffentlichen Schlüssel aus §3.2.
- **Zugangsdaten und Netz:** Die Anmeldung von Claude Code (§7) und Name und E-Mail von Git (§6) liegen
  im Dateisystem des Containers (kein Volume). Jeder Prozess darin, ein Build, ein Test oder ein
  Paket-Skript, kann sie mit `sudo` lesen, und der Netz-Ausgang ist nicht eingeschränkt. Sie entfallen mit dem
  Container (§4).
- **Docker-Socket des Hosts:** nicht gemountet (die drei Mounts stehen in §7). Der Socket im Container gehört dem Daemon im
  Container. Bei einem privilegierten Container ist das keine Grenze zum Host (siehe oben).
- **Host-Repo:** kein Mount; der Arbeitsbaum des Hosts ist in keinem der drei Mounts (`docker inspect`), also
  liegen darin keine uncommitteten oder ignorierten Dateien des Hosts. Das Bare-Repo enthält, was du mit
  `git push` einlieferst (§9).
- **Volumes:** `pgr-dev-workspace` (Bare-Repo und Klon) und `pgr-dev-docker` (Daemon im Container) überleben
  `u-boot down` und Neustart; `docker compose down -v` löscht sie.
- **Platz:** Das Image `pgr-dev:latest` belegt 2,05 GB, das Volume `pgr-dev-docker`
  5,87 GB nach einem Lauf von `make gates-schnell` (`docker system df -v` am Host; im Versuchs-Container und
  im frischen Container der Probe gleich), das Volume `pgr-dev-workspace` 37,14 MB (`docker system df -v` am Host, Versuchs-Container);
der Klon darin 47 MB mit `.git` und 18 MB ohne (`du -sh` im Versuchs-Container; Messung vor dem Bare-Repo, dessen Platz nicht gemessen ist). Die Größe wächst mit dem Host-Repo. Der Docker-Cache des Repos
  wächst im Volume; beobachte den Speicher des Hosts (`docker system df`, `df -h` auf dem
  Docker-Verzeichnis). Dass ein voller Speicher Builds und Tests scheitern lässt, ist im Container nicht
  ausgelöst (§12).
- **Ressourcen:** `cpus: 8` und `mem_limit: 12g` sind Obergrenzen des Containers; im Container zeigt
  `docker info` trotzdem die CPUs und den Speicher des Hosts.

## 12. Fehlersuche

| Beobachtung | Ursache und Schritt |
|---|---|
| Voraussetzungen unklar | `u-boot doctor` im Verzeichnis `pgr-dev`: prüft Docker, Compose, Git, `compose.yaml`, `devcontainer.json` und `Dockerfile` |
| Container nicht `healthy` | `docker compose ps`, dann `docker logs pgr-dev`; das Log von `dockerd` liegt auch im Container unter `/var/log/dockerd.log` |
| `docker logs` endet mit „dockerd nicht bereit“ und `failed to create NAT chain DOCKER: iptables failed` | der Container läuft ohne `privileged: true` (so gestartet in der Probe); `compose.yaml` prüfen. `start.sh` wartet 60 Sekunden, bevor es das meldet |
| `ssh` meldet `Permission denied (publickey)` | anderer Schlüssel als `ssh/id_ed25519` oder Benutzer `root` (beides in der Probe abgewiesen); Benutzer ist `vscode` |
| `git clone` nach `/workspaces` meldet `Permission denied` | ein frisches Volume gehört `root`; `start.sh` übergibt `/workspaces` bei jedem Start an `vscode`. In der Probe trat es auf, als der Entrypoint übergangen wurde |
| `git clone` aus dem Bare-Repo meldet `warning: remote HEAD refers to nonexistent ref, unable to checkout` | `HEAD` des Bare-Repos zeigt auf `master` (angelegt ohne `--initial-branch=main`); mit `--initial-branch=main` anlegen (§6) |
| `git push container main` meldet `! [rejected] main -> main (fetch first)` | der Container hat `main` weitergeschrieben; holen, mergen, erneut pushen (§9) |
| `ssh` oder `git` meldet `REMOTE HOST IDENTIFICATION HAS CHANGED!` | der Container läuft mit einem anderen Hostschlüssel als der Eintrag in `~/.ssh/known_hosts` (Neubau des Images, §5); die Entscheidung darüber liegt bei dir (§13) |
| Start meldet `Bind for 127.0.0.1:2222 failed: port is already allocated` | der Host-Port ist belegt; anderen Port wählen (§10) |
| `No space left on device` beim Bau oder in einem Test, der PostgreSQL startet | der Docker-Speicher ist voll: `docker system df`, ungenutzte Volumes und Images räumen. Im Container **nicht ausgelöst**; am Host wurde der Fehler bei vollem Docker-Speicher beobachtet |

## 13. Nicht geprüft

- Der Start von `claude` und die Anmeldung (interaktiv).
- Der volle `make gates` im Container.
- Windows- und macOS-Hosts, andere Docker-Versionen als in §1 genannt.
- Wie `u-boot` zu installieren ist.
- Die Wirkung von `u-boot down --volumes`.
- Der Fehler bei vollem Docker-Speicher im Container (§12).
- Andere Benutzerkennungen als 1000 am Host (§6).
- Die Frage der Erstverbindung am Terminal und jede Einstellung dazu (`StrictHostKeyChecking`,
  `UserKnownHostsFile`): Die Entscheidung darüber liegt beim Nutzer und ist offen (§5).
- Der Austausch im echten Host-Repo (§9); geprüft ist ein Scratch-Klon.
- `git fetch` vom Host aus einem leeren Bare-Repo im Container, ein `git push`, der mitten in der Übertragung
  abbricht (gelaufen ist ein `SIGKILL` an `git` nach 0,4 s: `main` im Bare-Repo blieb unverändert, die
  Wiederholung lieferte ein; ob der Abbruch mitten in der Übertragung lag, ist nicht festgestellt), und `origin`
  des Host-Repos über `ssh` bei gesetztem `core.sshCommand` (§9).
- Der Umbau eines schon laufenden Containers, der mit einem anderen Compose-Block erzeugt wurde (§4).
- `u-boot up` mit einer Variable nur aus der Umgebung (§3.5).
- Der Start einer zweiten Instanz mit `u-boot init pgr-zwei` und dem Kopf aus §3.5 (§10).
- Die Zeiten des Versuchs-Containers in §8 (5 s, 201 s, 1 s): nicht wiederholt.
- Der Zustand des Daemon-Caches bei den Zeiten von `make docs-check` (§8).
- Was nach `u-boot down` und `u-boot up` mit der Anmeldung von Claude Code geschieht (§4); wo `claude`
  seine Anmeldedaten ablegt (§7).
- Ein Zugriff auf Port 22 des Containers von einem anderen Container aus (§11).
- Das Verhalten bei fehlendem Schlüsselpaar vor dem ersten Start (§2) und andere `u-boot`-Versionen als 0.9.0.
