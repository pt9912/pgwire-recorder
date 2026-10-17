# Arbeiten im Dev-Container

Version: 0.3
Stand: 2026-10-17

## 1. Wofür der Dev-Container gut ist

Der Dev-Container ist ein privilegierter Docker-Container mit einem **eigenen Docker-Daemon**
(`dockerd`), `sshd`, Node.js und Claude Code. `u-boot` (Profil `sandbox`, `nestedRuntime: docker`,
Features `sshd`, `node` und `claude-code`) erzeugt seine Konfiguration, die `devcontainer`-CLI baut und
startet ihn. Das Repo liegt als Klon in einem Volume des
Containers; der Klon entsteht aus einem Bare-Repo im selben Volume, der Host liefert mit `git push` ein und
holt mit `git fetch` zurück (§6, §9). Der Arbeitsbaum des Hosts ist nicht eingebunden. Die Images, Tags, Netze und Container, die `make` im Container anlegt, liegen im
Daemon des Containers, nicht im Daemon des Hosts: Derselbe Tag (`pgwire-recorder:dev`) hat dort eine
andere Image-ID als am Host, und der Container sieht weder die Container noch die Netze des Hosts.
Claude Code läuft darin als eigene Session, getrennt von der Session am Host. Der Netz-Ausgang ist
auf eine Liste von Hosts eingeschränkt (§11).

Wofür er nicht gut ist:

- Er ist **keine Sandbox**. Der Container ist privilegiert, und `vscode` darf mit `sudo` ohne Passwort
  alles als `root` ausführen (§11); behandle ihn nicht als Schutz des Hosts vor fremdem Code.
- Er ersetzt keinen Release. Der volle `make gates` lief im Container grün (§8); die Zahlen sind
  Messwerte eines Laufs, keine Zusage.
- Er ist **nicht Teil des Repos**: Die Dateien, die ihn bauen, stehen in diesem Dokument und liegen
  in einem Verzeichnis außerhalb des Repos (§3). Wird dort etwas geändert, folgt dieses Dokument
  nicht von selbst.

Geprüft ist nur ein Linux-Host mit Docker 29.9.1 (Compose wird nicht gebraucht), `u-boot` 0.10.0 und
`devcontainer` 0.80.2 (Stand der Proben); im Container lief Docker 29.9.1, Node.js v24.22.0 und
Claude Code 2.1.302. `u-boot` führt nach seinem Changelog `nestedRuntime: docker` unter Docker auf Linux und auf
Colima zu; Docker Desktop sagt es nicht mehr zu (Zitat, nicht geprüft). Windows- und macOS-Hosts sind
nicht geprüft.

## 2. Voraussetzungen

- Docker Engine ≥ 24 (das fordert `u-boot doctor`, §12).
- `u-boot`, das Werkzeug, das Konfiguration und Dev-Container-Gerüst erzeugt
  (<https://github.com/pt9912/u-boot>), in der Version 0.10.0, mit der die Proben liefen (`u-boot
  --version`). Die Features `sshd` und `claude-code` und das Profil `sandbox` mit `nestedRuntime: docker`
  gelten für diese Version; andere Versionen sind nicht geprüft. Wie es zu installieren ist, steht
  dort; dieses Dokument prüft es nicht.
- Die `devcontainer`-CLI (npm-Paket `@devcontainers/cli`, Version 0.80.2 in den Proben, mit
  Node.js 22.18.0 am Host). Sie baut die Features ein, die `docker build` allein nicht kennt.
- Git.
- Am Host ein SSH-Client mit `ssh-keygen` (§3.2, §5, §9).
- Schreibrecht in dem Verzeichnis, in dem `pgr-dev` angelegt wird.
- Netz für den Build des Images (§3.4, §4).
- Das Schlüsselpaar aus §3.2 liegt vor dem ersten `devcontainer up` an dem Pfad, den `mk-override.sh`
  nennt ; fehlt die Datei, endet `devcontainer up` mit `invalid mount config … bind source path does not exist`, ohne dass Container oder Volumes entstehen (gemessen).
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

Das Host-Repo kommt in der Konfiguration nicht vor (§6, §9).

### 3.1 Gerüst mit `u-boot`

```bash
mkdir "$PGR_DEV" && cd "$PGR_DEV"
u-boot init pgr-dev --devcontainer --sandbox --no-git
u-boot config set devcontainer.features.node.enabled true
u-boot config set devcontainer.features.sshd.enabled true
u-boot config set devcontainer.features.claude-code.enabled true
u-boot config set devcontainer.sandbox.nestedRuntime docker
u-boot config set devcontainer.sandbox.egress.enabled true
u-boot config set devcontainer.sandbox.egress.allow gcr.io,pkg-containers.githubusercontent.com
u-boot generate devcontainer
```

`--sandbox` setzt das Profil, `--no-git` lässt das Verzeichnis ohne Git-Repo (ohne Remote `origin` erzeugt
`u-boot` keinen Klon-Schritt und warnt; den Klon legt §6 von Hand an). Die Hosts in `egress.allow`
sind die Messung aus §8 und §11. `u-boot generate devcontainer` meldet die Lockerungen des Profils
(`--privileged`, `--cap-add=NET_ADMIN`), den aktiven Login-Dienst `sshd` und das Drittanbieter-Feature `claude-code`.

Was `u-boot` erzeugt und was von Hand dazukommt:

| Datei | von `u-boot` | von Hand |
|---|---|---|
| `u-boot.yaml` | erzeugt und durch die `config set`-Befehle ergänzt (§3.3) | — |
| `.devcontainer/devcontainer.json`, `.devcontainer/Dockerfile`, `.devcontainer/docker-init.sh`, `.devcontainer/egress-init.sh` | `u-boot generate devcontainer` (§3.4) | — |
| `compose.yaml`, `.env.example`, `.gitignore`, `README.md`, `CHANGELOG.md`, `docker/`, `scripts/`, `docs/` | erzeugt | nicht gebraucht: `compose.yaml` hat keinen Dienst (`u-boot doctor` warnt, §12), die übrigen Dateien nutzt der Dev-Container nicht |
| `mk-override.sh` | — | ganz (§3.5) |
| `ssh/` | — | Schlüsselpaar (§3.2) |

Die Blöcke zwischen `BEGIN U-BOOT MANAGED BLOCK` und `END U-BOOT MANAGED BLOCK` gehören `u-boot`;
`u-boot generate devcontainer` überschreibt sie. Die Kennungen, die dort stehen (sie beginnen mit
`LH-`), sind Kennungen von `u-boot`, nicht dieses Repos.

### 3.2 Schlüsselpaar

```bash
mkdir -p ssh && chmod 700 ssh
ssh-keygen -t ed25519 -N '' -C pgr-dev -f ssh/id_ed25519
```

Der öffentliche Teil `ssh/id_ed25519.pub` wird als Mount `ro` auf `~/.ssh/authorized_keys` von `vscode`
eingebunden (§3.5) und ist dort der einzige Schlüssel, mit dem man sich anmelden kann (§5). Der private Teil bleibt im
Verzeichnis `ssh/` auf dem Host. `-N ''` legt den Schlüssel **ohne Passphrase** an: Wer die Datei
`ssh/id_ed25519` lesen kann, kann sich am Container anmelden, und dort hat `vscode` `sudo` ohne
Passwort (§11). Schütze das Verzeichnis (Modus 700) oder setze eine Passphrase.

### 3.3 `u-boot.yaml`

```yaml
schemaVersion: 1
project:
    name: pgr-dev
devcontainer:
    enabled: true
    features:
        claude-code:
            enabled: true
        node:
            enabled: true
        sshd:
            enabled: true
    profile: sandbox
    sandbox:
        nestedRuntime: docker
        egress:
            enabled: true
            allow:
                - gcr.io
                - pkg-containers.githubusercontent.com
```

`project.name` bestimmt den Namen der Volumes (`pgr-dev-workspace-<ID>`, `pgr-dev-docker-<ID>`) und den
Pfad `/workspaces/pgr-dev` im Container. `<ID>` ist die `devcontainerId` der `devcontainer`-CLI; sie
hängt am Ordner (§10).

### 3.4 Was `u-boot generate devcontainer` erzeugt

Gemessen mit dem `u-boot.yaml` aus §3.3 (die Datei trägt dort den Namen des Verzeichnisses):

- `.devcontainer/devcontainer.json`: `runArgs` `--privileged` und `--cap-add=NET_ADMIN`, zwei benannte
  Volumes (`/workspaces/pgr-dev` und `/var/lib/docker`), `postStartCommand` (startet die Beschränkung des Netzes und
  den Daemon), `remoteEnv` mit `GIT_TOKEN` und `ANTHROPIC_API_KEY` als Verweis `${localEnv:…}` (nie ein Wert), die drei Features
  (`ghcr.io/anthropics/devcontainer-features/claude-code:1`, `ghcr.io/devcontainers/features/node:1`,
  `ghcr.io/devcontainers/features/sshd:1`) und `remoteUser` `vscode`. Es gibt weder `forwardPorts` noch `appPort`.
- `.devcontainer/Dockerfile`: Basis `mcr.microsoft.com/devcontainers/base:debian` (Tag, kein Digest), `docker-ce` aus dem
  Docker-APT-Repository (`u-boot` prüft den Fingerabdruck des Schlüssels), `dnsmasq-base` und `nftables`.
- `.devcontainer/docker-init.sh` startet `dockerd` und wartet 60 Sekunden; `.devcontainer/egress-init.sh` setzt
  die Beschränkung des Netzes (§11).

Die Features laufen im Build: Claude Code 2.1.302 und Node.js v24.22.0 (die Features pinnen auf den Tag `1`; ein späterer
Build kann andere Versionen holen), `sshd` mit Port 2222. Der Build braucht Netz; die Herkunft der Quellen (Docker-APT-Repository,
Features, `mcr.microsoft.com`) ist nicht geprüft.

### 3.5 `mk-override.sh`

`devcontainer up` kennt keine Option für einen veröffentlichten Port, eine Obergrenze oder einen Namen, und `u-boot`
erzeugt dafür nichts (Port 2222 des Features wird nicht veröffentlicht). `--override-config` **ersetzt** die
`devcontainer.json`; das Skript schreibt deshalb eine Kopie mit Zusätzen:

```bash
#!/bin/bash
# erzeugt .devcontainer/override.json aus der von u-boot erzeugten devcontainer.json
set -euo pipefail
cd "$(dirname "$0")"
node -e '
const fs = require("fs");
const j = JSON.parse(fs.readFileSync(".devcontainer/devcontainer.json", "utf8").replace(/^\s*\/\/.*$/mg, ""));
const port = process.env.PGR_SSH_PORT || "2222";
j.runArgs = [...j.runArgs, "--name=" + (process.env.PGR_NAME || "pgr-dev"),
  "--cpus=" + (process.env.PGR_CPUS || "8"), "--memory=" + (process.env.PGR_MEM || "12g")];
j.appPort = ["127.0.0.1:" + port + ":2222"];
j.mounts = [...j.mounts, "type=bind,source=${localWorkspaceFolder}/ssh/id_ed25519.pub,target=/home/vscode/.ssh/authorized_keys,readonly"];
j.remoteEnv.CLAUDE_CODE_OAUTH_TOKEN = "${localEnv:CLAUDE_CODE_OAUTH_TOKEN}";
fs.writeFileSync(".devcontainer/override.json", JSON.stringify(j, null, 2) + "\n");
'
```

Die Zusätze: Name, `cpus` und `mem_limit` (Vorgabe 8 CPUs, 12 GiB; `docker inspect` zeigt `8000000000` und `12884901888`), den
Port `127.0.0.1:<PGR_SSH_PORT>:2222` (Vorgabe `2222` am Host), den Schlüssel als Mount `ro` und die zweite Variable
`CLAUDE_CODE_OAUTH_TOKEN` (§7). Das Skript läuft nach jedem `u-boot generate devcontainer` neu; `.devcontainer/override.json` ist abgeleitet
und gehört nicht in die Pflege. Das Skript braucht `node` am Host (die CLI braucht es ohnehin).

## 4. Start, Stopp, Löschen

Aus dem Verzeichnis `pgr-dev` (die Variablen für Claude Code: §7):

```bash
./mk-override.sh
devcontainer up --workspace-folder . --override-config .devcontainer/override.json
```

`devcontainer up` baut das Image beim ersten Mal, startet den Container und führt `postStartCommand` aus; es endet mit
`{"outcome":"success",…}`. Es gibt kein `healthy` wie bei Compose: Bereit ist der Container, wenn `ssh` (§5) antwortet. Der Build
lädt Pakete aus dem Netz. Gemessen (Docker-Cache des Hosts warm, Basis-Image und Features lokal vorhanden): 65 Sekunden (`devcontainer up`, erster Start);
mit `--build-no-cache` 70 Sekunden. Ein Build ohne warmen Cache des Hosts (Basis-Image nicht lokal) ist nicht gemessen.

| Ziel | Befehl | Folge |
|---|---|---|
| Stoppen und neu starten | `docker stop pgr-dev`, dann `devcontainer up …` | derselbe Container; 1,7 bis 1,8 Sekunden bis `ssh` und `dockerd` antworten; Bare-Repo, Klon (auch eine nicht versionierte Datei), die Images im Daemon des Containers und der Fingerabdruck des Hostschlüssels waren danach unverändert (gemessen) |
| Container entfernen | `docker rm -f pgr-dev`, dann `devcontainer up …` | neuer Container aus dem vorhandenen Image; die Volumes blieben und mit ihnen Bare-Repo, Klon und Images (gemessen, 6,4 und 7,9 Sekunden bis `ssh` und `dockerd` antworten); `/home/vscode` liegt im Dateisystem des Containers: `~/.gitconfig` und die Ablage von `claude` sind danach weg (`~/.claude` fehlte, gemessen) |
| Alles löschen | `docker rm -fv pgr-dev`, dann `docker volume rm <Namen der beiden Volumes>` (`docker volume ls`) | Bare-Repo, Klon und Docker-Cache im Container sind weg |
| Image löschen | `docker rmi <Image>` (`docker images`, Name beginnt mit `vsc-`) | danach baut der nächste Start neu |

`u-boot up` startet diesen Container nicht: Es steuert den Compose-Stack, und `compose.yaml` hat keinen Dienst
(`no service selected`, Exit 12, gemessen).

## 5. Verbinden per SSH

```bash
ssh -i "$PGR_DEV/ssh/id_ed25519" -o IdentitiesOnly=yes -p "${PGR_SSH_PORT:-2222}" vscode@127.0.0.1
```

Benutzer ist `vscode` (Mitglied der Gruppe `docker`). Der Host-Port ist an `127.0.0.1` veröffentlicht
(`docker ps` in der Probe: `127.0.0.1:2423->2222/tcp` bei `PGR_SSH_PORT=2423`); von anderen Rechnern ist der Container darüber nicht erreichbar.
`sshd` des Features lauscht im Container auf Port 2222 auf allen Adressen (`ss -ltn`: `0.0.0.0:2222` und `[::]:2222`).
Die Konfiguration des Features erlaubt Schlüssel und Passwort und `PermitRootLogin yes`; die Konten `root` und `vscode` haben kein Passwort
(`/etc/shadow`: `*` und `!`), und nur `vscode` hat einen Schlüssel. Gemessen: `ssh root@…` endet mit `Permission denied (publickey,password)`.
Mit dem Schlüssel einer anderen Instanz endet die Anmeldung als `vscode` ebenso.

Hostschlüssel: Das Feature erzeugt ihn beim Build des Images (Fingerabdruck nennt `root@buildkitsandbox`, Zeitstempel der Datei
im Build). Neustart und neuer Container aus demselben Image ändern ihn nicht; ein Neubau ohne Cache
(`devcontainer up … --build-no-cache`) und der Neubau nach einer Änderung von `egress.allow` ergaben einen anderen Fingerabdruck. Dann meldet `ssh` bei vorhandenem Eintrag für
`[127.0.0.1]:<Port>` `WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!`. Die Probe lief mit `-o UserKnownHostsFile=<Datei im Scratch>` und
`StrictHostKeyChecking=accept-new`; das Dokument empfiehlt dazu keine Einstellung (§13). Ohne Eintrag und ohne Terminal endet die Anmeldung mit
`Host key verification failed.`

Im Container läuft ein eigener Daemon: `docker info` nennt den Namen der Container-ID, und `docker ps -a` listet nicht die Container des Hosts,
auch nicht `pgr-dev` selbst.

## 6. Bare-Repo und Klon im Container

Im Volume `pgr-dev-workspace-<ID>`, eingehängt als `/workspaces/pgr-dev` (Name des Projekts in `u-boot.yaml`), liegt ein Bare-Repo `/workspaces/pgr-dev/austausch.git`; der Arbeitsklon
`/workspaces/pgr-dev/pgwire-recorder` entsteht daraus. Die Reihenfolge ist: Bare-Repo anlegen, vom Host einliefern
(§9), im Container klonen.

Bare-Repo anlegen (im Container, nach §5 per `ssh` angemeldet, oder dem Befehl hinten angehängt):

```bash
git init --bare --initial-branch=main /workspaces/pgr-dev/austausch.git
```

Nach dem Einliefern vom Host (`git push container main`, §9):

```bash
git clone /workspaces/pgr-dev/austausch.git /workspaces/pgr-dev/pgwire-recorder
cd /workspaces/pgr-dev/pgwire-recorder
git switch -c thema
make hooks-install
```

Gemessen am frischen Volume: `/workspaces` gehört `root`, das Verzeichnis des Volumes darin (`/workspaces/pgr-dev`) `vscode`
(UID 1000); das Bare-Repo darin gehört `vscode`. Wer außerhalb des Volumes anlegt, scheitert mit `fatal: cannot mkdir … Permission denied`. `git` und `git-receive-pack` liegen im Image (`git version 2.55.0`).
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

Meldungen von `git` und `ssh` stehen im Dokument, wie sie der Host (deutsche Locale) oder der Container
(englische Locale) ausgab; in einer anderen Locale lauten sie anders.

Name und E-Mail des Klons kommen aus der Git-Konfiguration deines Host-Repos, nicht aus dem Dokument.
Dieser Aufruf am Host bricht ab, wenn dort eine der beiden fehlt, und setzt sie im Klon lokal; `--global`
(`~/.gitconfig` im Container-Dateisystem) ist nach dem Entfernen des Containers (§4) weg, der Klon liegt im Volume.
Ohne Identität weist `git commit` mit „Author identity unknown“ ab:

```bash
n=$(git config user.name) e=$(git config user.email)
test -n "$n" -a -n "$e" && ssh -p ${PGR_SSH_PORT:-2222} -i $PGR_DEV/ssh/id_ed25519 vscode@127.0.0.1 \
  "cd /workspaces/pgr-dev/pgwire-recorder && git config user.name $(printf %q "$n") && git config user.email $(printf %q "$e")"
```

`make hooks-install` aktiviert den Commit-Wächter `.githooks/commit-msg` (`core.hooksPath` ist pro Klon); der
Träger `ai-harness-init` fehlt im frischen Klon (gitignoriert) und kommt mit `make traeger-fetch` (Netz nötig).

Im Container arbeitest du auf einem Branch, den du nicht `main` nennst, und pushst ihn in das Bare-Repo
(`origin` des Klons):

```bash
git push origin thema
```

Ein Stand des Hosts kommt mit `git switch main && git pull` im Klon.

## 7. Claude Code im Container

Claude Code 2.1.302 ist im Image (`claude --version` läuft im Container). Der Zugang kommt als Umgebungsvariable des Hosts,
ohne Mount und ohne Datei: `ANTHROPIC_API_KEY` (API-Schlüssel) oder `CLAUDE_CODE_OAUTH_TOKEN` (Abo-Token). Beide stehen als
Verweis `${localEnv:…}` in der `remoteEnv` von `.devcontainer/override.json` (§3.5), nie als Wert. `u-boot` trägt nur
`ANTHROPIC_API_KEY` ein; die zweite Variable ergänzt `mk-override.sh`.

**`remoteEnv` gilt nur für Prozesse, die die `devcontainer`-CLI startet.** Gemessen (Länge der Variablen, Platzhalterwerte):

| Weg | `ANTHROPIC_API_KEY` | `CLAUDE_CODE_OAUTH_TOKEN` |
|---|---|---|
| `devcontainer exec …` mit beiden Variablen in der Shell des Hosts | gesetzt (23) | gesetzt (28) |
| `devcontainer exec …` ohne die Variablen in der Shell | leer | leer |
| `ssh … vscode@…`, auch mit `-o SendEnv=…` | nicht gesetzt | nicht gesetzt |
| `docker exec -u vscode …` | nicht gesetzt | nicht gesetzt |

Der Wert wird bei jedem `devcontainer exec` neu aus der Shell des Hosts gelesen und steht nicht in `docker inspect` (`Config.Env`) und
nicht im Log des Containers. Das `sshd` des Features nimmt nur `LANG`, `LC_*`, `COLORTERM` und `NO_COLOR` an (`AcceptEnv`). Claude Code
startest du deshalb nicht in einer `ssh`-Sitzung, sondern mit `devcontainer exec`:

```bash
cd "$PGR_DEV"
devcontainer exec --workspace-folder . --override-config .devcontainer/override.json \
  bash -c 'cd /workspaces/pgr-dev/pgwire-recorder && claude'
```

`devcontainer exec` belegt eine Pseudo-Terminal-Sitzung, wenn der Host ein Terminal hat (`tty` meldete `/dev/pts/0`, gemessen unter `script`);
die interaktive Oberfläche von `claude` selbst lief nicht (§13). Für einen Aufruf ohne Terminal über `ssh` liest die Variable von `stdin`, damit der Wert
in keiner Kommandozeile steht: `printf %s "$ANTHROPIC_API_KEY" | ssh … vscode@127.0.0.1 'ANTHROPIC_API_KEY=$(cat) claude -p ok'`.

**Wo der Wert sichtbar ist** (gemessen mit einem Platzhalter): in der Umgebung des Prozesses im Container (`/proc/<PID>/environ`), und in der
Prozessliste des Hosts: `devcontainer exec` ruft `docker exec -e ANTHROPIC_API_KEY=<Wert> …` auf, `ps -eo args` am Host zeigt den Wert, solange der
Aufruf läuft. Nicht gefunden wurde er in `docker inspect`, im Log des Containers und in Dateien unter `/home/vscode`. Behandle beide Variablen wie das Konto, aus dem sie stammen.

**Verhalten von `claude -p ok` im Container** (Version 2.1.302, Egress wie in §11, jeweils Exit 1):

| Zustand der Variablen | Ausgabe |
|---|---|
| `ANTHROPIC_API_KEY` mit absichtlich falschem Wert (23 Zeichen) | `Invalid API key · Fix external API key` |
| beide fehlen | `Not logged in · Please run /login` |
| beide leer (`ANTHROPIC_API_KEY=`) | `Not logged in · Please run /login` |
| `ANTHROPIC_API_KEY` nur aus Leerzeichen | `Not logged in · Please run /login` |
| `CLAUDE_CODE_OAUTH_TOKEN` mit Platzhalter (28 Zeichen) | `Failed to authenticate. API Error: 401 Invalid bearer token` |
| `ssh`-Sitzung ohne Variable | `Not logged in · Please run /login` |

Die Meldung kommt beim Start des Aufrufs. Eine erneuerte Variable gilt ab dem nächsten `devcontainer exec` aus der Shell des Hosts; ein Neustart
des Containers ist dafür nicht nötig. Dass ein abgelaufenes Token mitten in einer Sitzung eine Meldung zeigt, ist nicht geprüft (§13).

Der echte Login mit einem gültigen Schlüssel oder Token **lief nicht**: Die Proben hatten keinen. Der Nutzer führt je Weg einen Befehl aus (Platzhalter ersetzen,
Wert nicht in die Shell-History legen):

```bash
cd "$PGR_DEV"
ANTHROPIC_API_KEY=<Schlüssel> devcontainer exec --workspace-folder . --override-config .devcontainer/override.json claude -p ok
CLAUDE_CODE_OAUTH_TOKEN=<Token> devcontainer exec --workspace-folder . --override-config .devcontainer/override.json claude -p ok
```

Die Abo-Anmeldung per `/login` im Container und der Weg über ein Abo-Token sind bei `u-boot` nicht zugesagt. Die Hosts der Anmeldung (`claude.ai`,
`platform.claude.com`, `claude.com`) stehen nicht in `egress.allow`; das Token wurde von `api.anthropic.com` geprüft (`401 Invalid bearer token`). Meldet ein
gültiges Token einen Netzfehler, trage den Host nach §11 ein.

Der Container hat genau drei Mounts: die beiden Volumes aus §3.4 und den öffentlichen Schlüssel (Mount `ro`; `docker inspect` am Scratch-Container);
Verzeichnisse mit Zugangsdaten und der Arbeitsbaum des Hosts sind nicht darunter.

## 8. Arbeiten im Container

**Annahme: pro Sitzung ein Slice** (Entscheidung des Nutzers vom 2026-10-17). Eine Sitzung ist die Zeit zwischen zwei Starts des
Containers; nach dem Slice stoppst du den Container und startest ihn neu (§4, rund 2 Sekunden, Klon und Images bleiben), das Zugangstoken
ist beim Start frisch. Dauert ein Slice länger als rund 8 Stunden, startest du ebenfalls neu. Die 8 Stunden sind die Grenze des Nutzers, nicht gemessen.

Im Klon laufen die Targets wie am Host; sie benutzen den Docker-Daemon des Containers:

```bash
make docs-check
make gates-schnell
make test
make gates
```

Gemessen im Klon (Stand 2026-10-17, `main` bei `9007c044`; Container mit 8 CPUs und 12 GiB, Host mit 31 GiB; `pgrep -fa "^make gates"` am Host war vor
dem Lauf leer und nach ihm leer):

| Target | Zeit | Bedingung |
|---|---|---|
| `make gates` | 1701 s (28 min 21 s), Exit 0 | frisches Volume, Daemon im Container ohne Images; höchstens 4,2 GiB Speicher des Containers (`docker stats`, alle 10 s) |
| `make gates` | 1701 s, Exit 0 | wiederholt in einem Container, dessen Daemon die Images aus zwei abgebrochenen Läufen hatte (beide Läufe sind unabhängig, die Gleichheit der Sekunden ist Zufall) |
| `make docs-check` | 2 s | warmer Daemon nach `make gates` |
| `make test` | 2 s | warmer Daemon nach `make gates` |
| `make gates-schnell` | 12 s | warmer Daemon nach `make gates` |

Ein Lauf mit leerem Cache des Hosts, ein zweiter Lauf von `make gates` mit warmem Daemon und ein Lauf während einer Last am Host sind nicht gemessen.
Die Zahlen hängen an Host, Cache und Last und sind keine Zusage. Der Speicher nach dem Lauf steht in §11.

Der Daemon im Container lädt die Images der Gates (zum Beispiel `d-check` und `a-check` von `ghcr.io`, `golang` von Docker Hub, `distroless` von `gcr.io`) selbst und
legt sie im Volume `pgr-dev-docker-<ID>` ab; die Images des Host-Daemons stehen ihm nicht zur Verfügung.

## 9. Austausch mit dem Host

Der Host schreibt in das Bare-Repo aus §6 und liest aus ihm; der Container sieht den Arbeitsbaum des
Hosts nicht. Die Befehle laufen im Host-Repo, `PGR_DEV` ist der Pfad aus §3:

```bash
git remote add container ssh://vscode@127.0.0.1:${PGR_SSH_PORT:-2222}/workspaces/pgr-dev/austausch.git
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
Dann `git fetch container`, `container/main` in den eigenen `main` mergen und erneut pushen; im Klon holt
`git switch main && git pull` den Stand (gelaufen auf `main`; auf `thema` aus §6, ohne Upstream, endet `git pull` mit Exit 1 und `There is no tracking information for the current branch.`). Deshalb arbeitet der Container nicht auf `main`.

`core.sshCommand` gilt für das ganze Host-Repo, nicht für ein Remote: Ein `git ls-remote` auf eine
`ssh://`-URL, die kein Remote ist, rief denselben Befehl mit demselben Schlüssel auf (`origin` über `ssh` ist
nicht gelaufen). Für mehrere Container gilt §10.

Gelaufen ist das in einem Scratch-Klon des Host-Repos, nicht im echten Host-Repo.

## 10. Mehrere Container nebeneinander

Jede Instanz hat ein eigenes Verzeichnis (eigener Ordnername) mit eigenem `u-boot init` (der Projektname bestimmt den Pfad im Container und den Namen
der Volumes), eigenem Schlüsselpaar und eigenen Werten für vier Variablen, die `mk-override.sh` (§3.5) aus der Shell liest und mit Vorgaben nennt:

| Variable | Wirkung | Vorgabe |
|---|---|---|
| `PGR_NAME` | Name des Containers (`--name`) | `pgr-dev` |
| `PGR_SSH_PORT` | Host-Port von `sshd` | `2222` |
| `PGR_CPUS`, `PGR_MEM` | `--cpus` und `--memory` | `8`, `12g` |

Zweite Instanz (Name `pgr-zwei`, Port 2423):

```bash
export PGR_DEV=<Pfad>/pgr-zwei
mkdir "$PGR_DEV" && cd "$PGR_DEV"
# u-boot init pgr-zwei --devcontainer --sandbox --no-git, Konfiguration (§3.1) und Schlüssel (§3.2) wie für die erste Instanz
cp <Pfad der ersten>/mk-override.sh .
PGR_NAME=pgr-zwei PGR_SSH_PORT=2423 ./mk-override.sh
devcontainer up --workspace-folder . --override-config .devcontainer/override.json
ssh -i "$PGR_DEV/ssh/id_ed25519" -o IdentitiesOnly=yes -p 2423 vscode@127.0.0.1
```

Gemessen mit den Namen `pgr-probe-dev` und `pgr-probe-zwei` (Ports 2423 und 2424): beide Container liefen gleichzeitig, jeder mit eigenen Volumes (`pgr-probe-zwei-docker-<ID>`,
`pgr-probe-zwei-workspace-<ID>`, die `<ID>` der beiden Ordner waren verschieden) und eigenem Hostschlüssel; `docker inspect` zeigt für `pgr-probe-zwei` ohne `PGR_CPUS` und `PGR_MEM` `8000000000` und `12884901888`.
Der Schlüssel der einen Instanz meldete an der anderen `Permission denied (publickey,password)` beziehungsweise `Host key verification failed.`
(der Eintrag in `known_hosts` fehlte). `PGR_CPUS` und `PGR_MEM` mit anderen Werten sind nicht gelaufen. Ein belegter Port ist nicht gelaufen (§13). Jeder Container hat seinen eigenen
Docker-Cache im eigenen Volume und braucht den Platz dafür noch einmal (§11). `--cpus` und `--memory` begrenzen jeden Container einzeln.

Austausch mit zwei Containern: Jeder Container hat sein eigenes Bare-Repo (§6) und am Host ein Remote
mit eigenem Namen und eigenem Port. `core.sshCommand` (§9) kennt einen Befehl je Host-Repo; mit einem Schlüssel
je Instanz nennt er beide. Die Befehle laufen im Host-Repo (nicht im Verzeichnis der Instanz):

```bash
git remote add container2 ssh://vscode@127.0.0.1:2423/workspaces/pgr-zwei/austausch.git
git config core.sshCommand "ssh -i <PGR_DEV der ersten>/ssh/id_ed25519 -i <PGR_DEV der zweiten>/ssh/id_ed25519 -o IdentitiesOnly=yes"
```

Mit dem Schlüssel der zweiten Instanz allein wies der erste Container ab (`git fetch container` endete mit „Bitte stellen Sie sicher, dass die korrekten
Zugriffsberechtigungen bestehen“); mit beiden Schlüsseln liefen `git fetch container` und `git fetch container2`, und `git branch -r` führte `container/main` und
`container2/main` nebeneinander. Die Befehle der übrigen Abschnitte (`ssh`, `git push`, `git fetch`)
nehmen Port und Schlüssel der Instanz. `PGR_SSH_PORT` steht nur in der Shell, die `mk-override.sh` ruft: Die Befehle
in §5 und §9 lösen `${PGR_SSH_PORT:-2222}` ohne gesetzte Variable zu `2222` auf. Für die zweite Instanz setzt du
`export PGR_SSH_PORT=2423` in der Shell oder schreibst `-p 2423` bzw. `:2423` in die URL.

## 11. Grenzen und Sicherheit

- **Privilegiert, `sudo` ohne Passwort:** `--privileged` (`docker inspect`: `Privileged=true`, `CapAdd=[CAP_NET_ADMIN]`) gibt dem Container weitgehenden Zugriff
  auf den Kernel des Hosts, und `sudo -n -l` zeigt für `vscode` `(root) NOPASSWD: ALL`; dazu ist
  `vscode` in der Gruppe `docker`. Was ein Prozess darin am Host ändern könnte, ist nicht geprüft.
  Behandle den Container deshalb nicht als Sandbox: Lass darin nur Arbeit laufen, der du traust, und
  prüfe Ergebnisse, die du zurückholst (§9), als Branch, bevor du sie mergst. `vscode` gehört das
  Bare-Repo und darf dort jeden Stand ablegen (gelaufen: Push von `main` und von Branchen); was du am Host
  übernimmst, entscheidest du beim Merge.
- **Login-Dienst `sshd`:** Das Feature lässt Schlüssel, Passwort und `PermitRootLogin yes` zu und lauscht im Container auf allen Adressen; die Konten haben
  kein Passwort (§5). Veröffentlicht ist nur der Host-Port, und der ist an `127.0.0.1` gebunden. Ein anderer Container auf dem Host erreichte Port 2222
  über die Adresse des Containers im Bridge-Netz (gemessen mit einem `alpine`-Container, `nc -z`); geschützt ist er dort nur durch den einen öffentlichen Schlüssel aus §3.2.
- **Netz-Ausgang, eine Leitplanke:** `devcontainer.sandbox.egress` startet beim Containerstart einen lokalen DNS-Server (`dnsmasq`), der nur die Namen der Liste
  auflöst, und eine Firewall (`nftables`), die nur deren Adressen durchlässt. `example.com` löste im Container nicht auf. Die Liste ist die Standardliste von `u-boot`
  (GitHub, Debian, mit `nestedRuntime: docker` Docker Hub und `ghcr.io`, `registry.npmjs.org`, `api.anthropic.com`; `egress-init.sh` führt sie) und `egress.allow` aus §3.3.
  **Gemessen, was `make gates` braucht:** die Auflösungen im Container beim Lauf (`registry-1.docker.io`, `auth.docker.io`, `production.cloudfront.docker.com`,
  `ghcr.io`) sind in der Standardliste; zusätzlich fehlten `pkg-containers.githubusercontent.com` (Blobs von `ghcr.io`, ohne es `docker run … a-check` mit
  `failed to copy: httpReadSeeker … pkg-containers.githubusercontent.com`) und `gcr.io` (das Runtime-Image `distroless`, ohne es `make build` mit
  `failed to resolve source metadata for gcr.io/distroless/…`); mit beiden lief `make gates` grün. Die Anmeldung von Claude Code erreichte `api.anthropic.com`; `http-intake.logs.us5.datadoghq.com`
  (Telemetrie) wurde angefragt und blieb gesperrt, `claude -p` lief trotzdem. Die Hosts der Anmeldung im Browser sind nicht gemessen (§7). **Die Beschränkung gilt nicht für verschachtelte Container:** ein
  `alpine`-Container im Daemon des Containers erreichte `example.com` (gemessen), und mit `--privileged` kann jeder Prozess die Regeln lösen. `u-boot doctor` warnt
  bei `nestedRuntime: docker` mit Beschränkung genau davor; ein Build, der im Container einen Befehl mit `RUN` (Netz) ausführt, ist also nicht auf die Liste beschränkt.
- **Zugangsdaten:** Name und E-Mail von Git (§6) liegen im Dateisystem des Containers; `~/.gitconfig` ist nach dem Entfernen des Containers weg (gemessen). Die Variablen aus §7
  stehen in der Umgebung der Prozesse, die `devcontainer exec` startet, und in der Prozessliste des Hosts. Wohin `claude` nach einem Login im Container schreibt,
  ist nicht geprüft; jeder Prozess darin, ein Build, ein Test oder ein Paket-Skript, kann Dateien des Containers mit `sudo` lesen.
- **Docker-Socket des Hosts:** nicht gemountet (die drei Mounts stehen in §7). Der Socket im Container gehört dem Daemon im
  Container. Bei einem privilegierten Container ist das keine Grenze zum Host (siehe oben).
- **Host-Repo:** kein Mount; der Arbeitsbaum des Hosts ist in keinem der drei Mounts (`docker inspect`), ein
  Prozess ohne `sudo` sieht ihn also nicht. Ob ein Prozess mit `root` ihn bei `--privileged` über
  Gerätedateien lesen kann, ist nicht geprüft (§13). Das Bare-Repo enthält, was du mit
  `git push` einlieferst (§9).
- **Volumes:** `pgr-dev-workspace-<ID>` (Bare-Repo und Klon) und `pgr-dev-docker-<ID>` (Daemon im Container) überleben
  Stoppen, Neustart und das Entfernen des Containers (§4); `docker volume rm` löscht sie.
- **Platz:** Das Image belegt 2,15 GB, das Volume `pgr-dev-docker-<ID>` 12,59 GB nach einem Lauf von `make gates` (`docker system df -v` am Host;
  im Container `du -sh /var/lib/docker`: 13 GB, `docker system df`: 15 Images mit 6,7 GB und 216 Einträge des Build-Caches mit 7,6 GB), das Volume
  `pgr-dev-workspace-<ID>` 27,76 MB; der Klon darin 32 MB, das Bare-Repo 124 KB (`du -sh` im Container). Die Größe wächst mit dem Host-Repo und mit dem Docker-Cache
  des Repos; beobachte den Speicher des Hosts (`docker system df`, `df -h` auf dem Docker-Verzeichnis). Dass ein voller Speicher Builds und Tests scheitern lässt,
  ist im Container nicht ausgelöst (§12).
- **Ressourcen:** `--cpus=8` und `--memory=12g` sind Obergrenzen des Containers (`docker inspect`: `8000000000`, `12884901888`; `cgroup memory.max` im Container: 12884901888);
  im Container zeigt `docker info` trotzdem die CPUs (20) und den Speicher des Hosts. `make gates` brauchte höchstens 4,2 GiB (§8); der Kernel hat den Container dabei nicht beendet (`OOMKilled=false`).

## 12. Fehlersuche

| Beobachtung | Ursache und Schritt |
|---|---|
| Voraussetzungen unklar | `u-boot doctor` im Verzeichnis `pgr-dev`: meldet `compose.yaml has no services` und, bei `nestedRuntime: docker` mit Beschränkung, die Warnung zum Netz-Ausgang (§11); beide sind erwartet (0 Fehler, 2 Warnungen, gemessen) |
| `devcontainer up` endet mit `invalid mount config … bind source path does not exist` | das Schlüsselpaar aus §3.2 fehlt am Pfad `ssh/id_ed25519.pub` (§2) |
| Start nicht bereit | `docker logs pgr-dev`; das Log von `dockerd` liegt im Container unter `/var/log/u-boot-dockerd.log` (die Datei gibt es; ein Start ohne Daemon ist nicht ausgelöst) |
| `ssh` meldet `Permission denied (publickey,password)` | anderer Schlüssel als `ssh/id_ed25519` oder Benutzer `root` (beides gemessen); Benutzer ist `vscode` |
| `claude` meldet `Not logged in` in einer `ssh`-Sitzung | die Variablen aus §7 kommen nur über `devcontainer exec` in den Prozess |
| `git init` meldet `fatal: cannot mkdir … Permission denied` | der Pfad liegt außerhalb des Volumes; `/workspaces` gehört `root`, `/workspaces/pgr-dev` `vscode` (§6) |
| `git clone` aus dem Bare-Repo meldet `warning: remote HEAD refers to nonexistent ref, unable to checkout` | `HEAD` des Bare-Repos zeigt auf `master` (angelegt ohne `--initial-branch=main`); mit `--initial-branch=main` anlegen (§6) |
| `git push container main` meldet `! [rejected] main -> main (fetch first)` | der Container hat `main` weitergeschrieben; holen, mergen, erneut pushen (§9) |
| `ssh` oder `git` meldet `REMOTE HOST IDENTIFICATION HAS CHANGED!` | der Container läuft mit einem anderen Hostschlüssel als der Eintrag in `~/.ssh/known_hosts` (Neubau des Images, §5); die Entscheidung darüber liegt bei dir (§13) |
| Docker meldet `driver failed programming external connectivity … ` beim Start | der Host-Port ist belegt (gemessen mit `docker run -p` auf einem belegten Port); anderen Port wählen (§10) |
| `docker pull` im Container scheitert mit einem Fehler zu einem Host, den die Liste nicht nennt | der Host fehlt in `egress.allow` (§11): eintragen, `u-boot generate devcontainer`, `./mk-override.sh`, Container entfernen (`docker rm -f`) und `devcontainer up`: Der Build läuft neu (45 s gemessen) und der Host löst auf; ein `devcontainer up` bei laufendem Container ändert nichts (gemessen mit `example.org`). Der Hostschlüssel ändert sich dabei (§5) |
| `No space left on device` beim Bau oder in einem Test, der PostgreSQL startet | der Docker-Speicher ist voll: `docker system df`, ungenutzte Volumes und Images räumen. Im Container **nicht ausgelöst**; am Host wurde der Fehler bei vollem Docker-Speicher beobachtet |

## 13. Nicht geprüft

- **Der echte Login von Claude Code mit gültigem Zugang**: API-Schlüssel und Abo-Token lagen den Proben nicht vor. Der Nutzer führt die zwei Befehle aus §7 aus. Auch der
  Start der interaktiven Oberfläche von `claude` (nur `claude -p` lief) und das Verhalten eines Tokens, das mitten in einer Sitzung abläuft.
- **Entscheidung `podman` oder `docker` als `nestedRuntime` (offen, Nutzer):** gemessen ist `docker` (Daemon im Container, `make gates` grün). `u-boot` 0.10.0 nennt für `podman`
  nach seinem Changelog `podman build` auf Colima (arm64), mit `/dev/net/tun` und ohne `buildx`, `docker compose` und `# syntax=`; hier nicht gelaufen. Das `Dockerfile` des Repos nutzt
  `RUN --network=none`, `--platform=$BUILDPLATFORM` und mehrstufige Builds (gelesen, nicht gelaufen); ob `make` mit `podman` läuft, ist nicht geprüft. Der Host hier ist Linux/amd64, nicht Colima.
- Der Hostschlüssel: Das Feature erzeugt ihn im Build (§5); ob und welche Einstellung (`StrictHostKeyChecking`, `UserKnownHostsFile`) sinnvoll ist, entscheidet der Nutzer. Die Frage der Erstverbindung am Terminal ist nicht geprüft.
- Der Umbau des laufenden Containers `pgr-dev` des Nutzers auf dieses Gerüst (Volumes `pgr-dev-docker` und `pgr-dev-workspace` des alten Gerüsts, Handarbeit des Nutzers): nicht gelaufen. Die neuen Volumes tragen `<ID>` im Namen und sind andere als die des alten Gerüsts.
- Die Aufrufe aus §6 für `make traeger-fetch` im Container-Klon (Identität und `make hooks-install` liefen).
- Die Meldungen von `git` in §6 (außer den dort als gemessen genannten) und alle Meldungen in §9 stammen aus der Probe des Vorgängers; im neuen Gerüst liefen Push von `main`, `git fetch` aus zwei Containern, die Bare-Repo-Meldungen in §6 und der Klon ohne Fehler.
- Windows- und macOS-Hosts, Docker Desktop, andere Docker-Versionen als in §1 genannt, andere Versionen von `u-boot` und `devcontainer`.
- Wie `u-boot` und die `devcontainer`-CLI zu installieren sind.
- Der Fehler bei vollem Docker-Speicher im Container (§12).
- Andere Benutzerkennungen als 1000 am Host (§6).
- Der Austausch im echten Host-Repo (§9); geprüft ist ein Scratch-Klon.
- Ob ein Prozess mit `root` im Container bei `--privileged` über Gerätedateien den Arbeitsbaum des Hosts lesen
  kann (§11).
- `git fetch` vom Host aus einem leeren Bare-Repo im Container, ein `git push`, der mitten in der Übertragung
  abbricht (gelaufen ist in der Probe des Vorgängers ein `SIGKILL` an `git` nach 0,4 s: `main` im Bare-Repo blieb unverändert, die
  Wiederholung lieferte ein; ob der Abbruch mitten in der Übertragung lag, ist nicht festgestellt), und `origin`
  des Host-Repos über `ssh` bei gesetztem `core.sshCommand` (§9).
- Ein Lauf von `make gates` mit leerem Cache des Hosts, ein zweiter Lauf mit warmem Daemon und ein Lauf unter Last am Host (§8); die 8-Stunden-Grenze (§8).
- `PGR_CPUS` und `PGR_MEM` mit anderen Werten, ein belegter Port bei `devcontainer up` (gemessen nur mit `docker run`), ein Start ohne Daemon (`nestedRuntime` wirkungslos, §12).
- Ein Build ohne warmen Cache des Hosts, bei dem das Basis-Image nicht lokal vorhanden ist (§4), und die Herkunft des Basis-Images, der Features und des Docker-APT-Repositorys (§3.4).
- Die Hosts, die Claude Code für die Anmeldung im Browser braucht (`claude.ai`, `platform.claude.com`, `claude.com`), und die Updates, Plugins und MCP-Connectors (§7, §11).
- Die Meldung des Hostschlüssels bei der Erstverbindung ohne `UserKnownHostsFile`-Option (die Proben liefen mit einer Datei im Scratch).
