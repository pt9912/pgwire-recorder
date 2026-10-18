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
- Mehrere Container nebeneinander (Entscheidung des Nutzers, 2026-10-17: parametrisiert über `.env` mit `COMPOSE_PROJECT_NAME`,
  `PGR_NAME`, `PGR_SSH_PORT`, `PGR_CPUS`, `PGR_MEM`, Befehle für die zweite Instanz), Grenzen und Sicherheit (privilegiert, nur 127.0.0.1, nur
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

- [x] `docs/developer/devcontainer.md` liegt vor und deckt die Punkte aus §1 ab; jeder Schritt ist
      am laufenden Container nachgeführt oder als „nicht geprüft“ markiert, die Dateiinhalte
      stimmen mit denen des Versuchs überein.
- [x] `make kopf-check` und `make docs-check` grün; `make gates` grün (Stop-Hook, §6).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/developer/devcontainer.md` | neu | Liefer-Punkt: das Dokument nach §1 |

- Der Implementer führt die Schritte am laufenden Container unter `/Development/pgr-dev` nach
  (nur lesend und per SSH; den Container nicht verändern) und übernimmt die Dateiinhalte von dort.
- Zusätzlich baut der Implementer den Container in einem Scratch-Verzeichnis (anderer Name, Port 2223)
  aus den Codeblöcken des Dokuments neu auf; die Befehle und Beobachtungen stehen in §7. Das Dokument
  führt über die Liste aus §1 hinaus die Abschnitte Parallel arbeiten, Fehlersuche und Nicht geprüft.
- Nach dem Review (Nacharbeit): Sicherheitsaussagen sind auf das Belegte enger gefasst (Container
  privilegiert, `vscode` mit `sudo` ohne Passwort, `/host-repo` als Mount `ro` nur gegen Prozesse
  ohne `sudo`), Quellen und Lücken stehen im Dokument (§7, Abschnitt Nacharbeit).
- Pfade des Hosts stehen im Dokument nicht fest (Entscheidung des Nutzers, 2026-10-17): `PGR_REPO` und `PGR_DEV` in
  den Befehlen, im Block `compose.yaml` `${PGR_REPO:?…}` und ein relativer Pfad des Schlüssels. Der Block ist damit
  nicht mehr byte-gleich zur Datei im Messort `/Development/pgr-dev`; das Dokument nennt den Unterschied (§3.5, §13).
- Mehrere Container (Nutzer, 2026-10-17): `container_name`, `hostname`, `image`, Host-Port, `cpus`, `mem_limit` im Block
  sind `${…:-Vorgabe}`; Abschnitt 10 des Dokuments führt die zweite Instanz. Das Netz `pgr-dev-default` steht im Kopf von
  `u-boot` und bleibt (Kopf wörtlich, §3.1 des Dokuments).
- Das Dokument ist Ist-Zustand ohne Chronik (§3.11); es nennt keine Kennungen dieses Repos
  (Lastenheft, Spezifikation, ADRs, Slices, Wellen). Die Kennungen, die mit `LH-` beginnen und in den
  wörtlich übernommenen Blöcken von `u-boot` stehen (Codeblöcke von `Dockerfile` und `compose.yaml`), sind
  Kennungen von `u-boot`; ein Satz in §3.1 des Dokuments sagt es. Nennt es doch eine ADR-Kennung, ist sie ein Link
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
  des Nutzers). — **Ausgang (Planner, Closure):** weiter offen → Register
  `BEO-REPO/stop-hook-positivliste-ohne-docs-developer` (1×, Beleg
  `evidence/slice-harness-devcontainer-doku.md`). Die Frage an den Nutzer steht in §7.
- **Das Dokument beschreibt einen Versuch außerhalb des Repos.** Wird der Container umgebaut,
  veraltet das Dokument still; kein Sensor prüft es gegen den Container. — **Ausgang (Planner,
  Closure):** weiter offen → `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Beleg
  `evidence/slice-harness-devcontainer-doku.md`: zwei Reviews fanden Aussagen über die Probe hinaus;
  das stille Veralten selbst ist nicht eingetreten und bleibt ohne Sensor).
- **Nicht nachführbare Schritte** (etwa der volle `make gates` im Container, Zeit nicht
  gemessen). — **Ausgang:** entfallen: Das Dokument markiert sie als „nicht geprüft“ (§13 des
  Dokuments, von der Verifikation gegen den Plan abgeglichen: „Stimmt“).
- **Form der Ausgänge (Verifikation V-2):** erledigt; jedes der drei Risiken trägt oben genau einen der drei Ausgänge.
- **Randformen:** kein neuer Vertrag, keine Randformen im Sinn von `AGENTS.md` §3.12.

## 7. Closure-Notiz

**Closure (Planner, 2026-10-17).** Der Slice liefert `docs/developer/devcontainer.md` (ein Liefer-Punkt, eine
Schicht). Belege: Reviews `docs/reviews/2026-10-17-review-slice-harness-devcontainer-doku.md` (0 HIGH, 1 MEDIUM,
4 LOW, 2 INFO) und `…-review-2-…` (0 HIGH, 2 MEDIUM, 5 LOW, 2 INFO), Verifikation
`docs/reviews/2026-10-17-verify-slice-harness-devcontainer-doku.md` (bestanden mit Vorbehalt), Nacharbeit-Commits
`6ca83617` und `983cafdf`. Alle Findings der Reviews sind in der Nacharbeit behoben oder als akzeptierte Lücke
geführt (Nutzerlinie: LOW und INFO ohne Zusage ohne Test sind akzeptiert, Ausgang in dieser Notiz).

**Stand der Closure.** DoD-Punkt 2 (`make kopf-check`, `make docs-check` und `make gates` grün auf dem Endstand)
ist **nicht belegt**: Die Läufe stammen aus `531f357e`, vor der Nacharbeit; ein voller `make gates` lief während
der Closure (PID 2209712), deshalb fuhr der Planner weder `docs-check` noch `kopf-check`. Das Häkchen bleibt leer, der
Slice bleibt in `in-progress/`. Die Closure wartet auf den grünen Lauf auf dem Endstand; danach hakt der Planner
den Punkt ab, und der reine `git mv` nach `done/` folgt.

**Ausgang je Risiko (§6).** Stop-Hook und Positivliste: weiter offen → `BEO-REPO/stop-hook-positivliste-ohne-docs-developer`.
Veralten des Dokuments: weiter offen → `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`. Nicht nachführbare
Schritte: entfallen (Dokument §13 markiert sie, Verifikation bestätigt).

**Ausgang je Befund niedriger Schwere.** R2 F-3 bis F-7 (LOW) und F-8, F-9 (INFO), V-3 bis V-5: behoben im Dokument
oder Plan (Tabellen oben). Akzeptierte Lücken ohne Zusage im Dokument: Schreiben auf `/host-repo` mit `sudo`,
Hostschlüssel-Frage, Fetch ins echte Host-Repo, `docker compose up -d` allein, Zeiten des Versuchs-Containers, Zugang
eines zweiten Containers zu Port 22, Ablage der Anmeldedaten von `claude`; alle stehen in §13 des Dokuments.

**Meldeschwelle (V-3), Entscheidung des Planners, 2026-10-17: weiter, kein Schnitt.** Gemessen 441 (Metrik `AGENTS.md`
§3.14, Basis `bda6efef`), über der Meldeschwelle 400, unter der Marke 500. Gründe: ein Liefer-Punkt, eine Schicht, das
Mehr ist die Nacharbeit zweier Reviews im selben Dokument, Review und Verifikation liegen vor, ein Schnitt trennte nur
Abschnitte desselben Dokuments. Die Rückführung nach §4 trat nicht ein.

**Lerneintrag.** Form: benannte Lücke, im Register geführt. (a) Die Schätzung 150 bis 300 lag unter der Messung
(407 nach dem ersten Review, 441 am Endstand; Obergrenze um 141 überschritten, 47 %). Ursache: Sie rechnete mit einem Dokument,
das die Schritte beschreibt; gemessen wurde ein Dokument, das drei Dateiinhalte wörtlich trägt und zu jeder Aussage
die Grenze der Probe, und die Nacharbeit zweier Reviews fügte Zeilen hinzu (385 → 407 → 441, nie Löschungen). Nach der
Zählregel vom 2026-10-16 ist das ein Auftreten von `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Zähler 11);
`slice-harness-schaetzung-planner-faktor` trägt den Faktor nur für Tests zu Produktcode, nicht für Prosa-Dokumente mit
wörtlichen Dateiinhalten; ob er dafür einen Ansatz nennt, entscheidet der Nutzer (Frage unten, keine Zuweisung).
(b) Der zweite Review fand zwei MEDIUM (Bindung von `sshd`, Anmeldedaten) nach dem ersten, der ein MEDIUM zur
Sicherheit gefunden und behoben hatte: Die Nacharbeit des ersten Reviews engte die Aussage *Mount `ro`* ein, prüfte
aber die übrigen Sicherheitsaussagen nicht gegen die Bindung im Container und gegen die Lage der Zugangsdaten. Das ist
die Klasse *Aussage über die Prüfung hinaus* (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, Auftreten Nr. 72 in
den Dateien): Quantoren („nur Loopback“) standen für den veröffentlichten Port und wurden als Aussage über den
Container gelesen. Die Regel `AGENTS.md` §3.11 gilt unverändert; kein neuer Satz. (c) Ein Dokument unter
`docs/developer/` hat keinen Sensor gegen den Container und keine Positivliste im Stop-Hook; beides ist als
Beobachtung geführt, nicht als Regel.

**Paarungen.** Anker: keine neue verkörperte Regel, deshalb kein neuer `seit slice-harness-devcontainer-doku`;
`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` behält den Stand `verkörpert` mit seinen Ankern. Folge-Slice: keiner
neu; `slice-harness-schaetzung-planner-faktor` (`next/`) bekommt keine Sendung, der Datenpunkt steht im Register.
Register: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Evidence, Auftreten), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`
(Evidence, Zähler 11), `BEO-REPO/stop-hook-positivliste-ohne-docs-developer` (neu, 1×), `BEO-REPO/meldeschwelle-ohne-artefakt`
(Gegenbeleg: Schwelle in §4, Entscheidung als datierte Zeile hier; Zähler bleibt 2). Wellenlos: kein Eintrag in der Roadmap.

**Fragen an den Nutzer (offen, nichts geändert).**

1. Soll `docs/developer/` auf die Positivliste des Stop-Hooks (Regeländerung in `tools/harness/dateiliste.sh` und
   `SPEC-052` samt Gegenprobe)? Bis dahin verlangt jeder Commit unter `docs/developer/` einen vollen `make gates`.
2. Braucht der Betrieb Schritte für die Zugangsdaten von `claude` im Container (Abmelden vor `u-boot down`, ein
   Netz-Ausgang nur zu bestimmten Zielen, Verzicht auf die Anmeldung im Container)? Das Dokument nennt nur die Lage.
3. Soll `slice-harness-schaetzung-planner-faktor` einen Ansatz für Prosa-Dokumente mit wörtlichen Dateiinhalten nennen
   (hier: 441 gemessen gegen 150 bis 300 geschätzt)?


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
| §3.3–3.5 Dateiinhalte | `diff` der Codeblöcke gegen `/Development/pgr-dev/…` | kein Unterschied (`D-gleich`, `S-gleich`, `C-gleich`); gilt für `compose.yaml` nur bis zur Portabilität unten |
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
voller `make gates` läuft nicht (Auftrag; §6 Risiko Stop-Hook). Größe vor der Nacharbeit (Stand `531f357e`, `git diff --numstat bda6efef..531f357e -- docs/developer`): 385 Einfügungen, 0 Löschungen, Summe 385; spätere Stände siehe Nacharbeit unten. `make docs-check` und `make kopf-check` stammen aus diesem Stand.

**Nicht geprüft** (im Dokument §13 benannt): Start von `claude` und Anmeldung; voller `make gates` im
Container; Windows/macOS; Installation von `u-boot`; `u-boot down --volumes`; Fehler bei vollem
Docker-Speicher im Container; andere UID als 1000 am Host (`safe.directory`); neue SSH-Hostschlüssel nach
Neubau des Images (nach `u-boot down`/`up` meldete `ssh` mit gespeichertem Schlüssel keine Warnung); seit der Nacharbeit
zusätzlich Hostschlüssel-Frage der Erstverbindung, Fetch ins echte Host-Repo, Schreiben auf `/host-repo` mit `sudo`,
`docker compose up -d` allein.

### Nacharbeit nach dem Review (2026-10-17)

Review `docs/reviews/2026-10-17-review-slice-harness-devcontainer-doku.md` (0 HIGH, 1 MEDIUM, 4 LOW, 2 INFO).
Nur Text in `docs/developer/devcontainer.md` und diesem Plan.

| Befund | Ausgang |
|---|---|
| F-1 MEDIUM Sicherheitsaussagen | behoben: `sudo -n -l` am Versuchs-Container (nur lesend) zeigt `(root) NOPASSWD: ALL`; §1, §3.2, §7, §9, §11 sagen „Mount `ro`“ statt „schreibgeschützt“ und nennen, dass ein Prozess mit `sudo` den Mount ändern könnte (nicht ausgeführt, nur gefolgert), dazu den Umgang (nur vertrauenswürdige Arbeit, Rückholen als Branch/Bundle prüfen). Durchsicht der übrigen Quantoren: „nur Loopback“ (belegt durch `docker port`, `ss`) und „nur Schlüssel“ (belegt durch drei Abweisungen) unverändert |
| F-2 Quelle der Zeiten | behoben: §8 nennt `scripts/nachweis2.sh`/`nachweis3.sh` des Nutzers als Quelle der jeweils ersten Zahl, nicht wiederholt |
| F-3 Lücken | behoben: §13 um Hostschlüssel-Frage, Fetch ins echte Host-Repo, Schreiben mit `sudo` und `docker compose up -d` ergänzt; §5 und §9 erwähnen die ersten beiden |
| F-4 `compose up -d` | behoben: auf „laut Hilfetext von `u-boot`“ eingegrenzt, `docker compose up -d` allein in §13 |
| F-5 Tabelle §3.1 | behoben: Zeile `docker/`, `scripts/`, `docs/` (leer erzeugt; `ls` am Versuchs-Container gelesen) |
| F-6 `safe.directory` | behoben: aus dem Codeblock von §6 genommen, im Absatz als Eintrag bei fremder UID beschrieben |
| F-7 Passphrase / Kennung | behoben für das Dokument: §3.2 nennt Schlüssel ohne Passphrase samt Folge. Beobachtung: der Plan-Kopf trägt die Commit-Kennung-Frage; kein History-Eingriff |

**Größe nach der Nacharbeit** (§4, `git diff --numstat bda6efef..d725a22c -- docs/developer`): 407 Einfügungen,
0 Löschungen, Summe 407; über der Meldeschwelle 400, unter der Marke 500.

**Akzeptierte Lücken:** Schreiben auf `/host-repo` mit `sudo` nicht ausgeführt; Hostschlüssel-Frage, Fetch ins echte
Host-Repo, `docker compose up -d` allein, Zeiten des Versuchs-Containers nicht geprüft bzw. nicht wiederholt.

### Nacharbeit nach Review 2 und Verifikation (2026-10-17)

Reports `docs/reviews/2026-10-17-review-2-slice-harness-devcontainer-doku.md` (0 HIGH, 2 MEDIUM, 5 LOW, 2 INFO)
und `docs/reviews/2026-10-17-verify-slice-harness-devcontainer-doku.md` (V-1 bis V-5). Nur Text in
`docs/developer/devcontainer.md` und diesem Plan; Dokument-Commit `6ca83617`. Neue lesende Proben am
Versuchs-Container (SSH, nichts verändert): `ss -ltn` = `0.0.0.0:22`, `[::]:22` (und `127.0.0.11` für den
Docker-DNS); `findmnt`: die vier Mounts, `/home/vscode` liegt im Overlay des Containers, mit `.gitconfig`;
`ls -a /host-repo` zeigt `.git`, `.claude`, `.harness`, `recordings`, `.harness/state/` (ignoriert) ist
lesbar; `du -sh` Klon 47 MB, ohne `.git` 18 MB; `docker system df -v` Volume `pgr-dev_pgr-dev-workspace`
37,14 MB; `u-boot --version` 0.9.0; `ssh-keygen`, `scp` am Host vorhanden.

| Befund | Ausgang |
|---|---|
| R2 F-1 MEDIUM `sshd` an Loopback | behoben: §5, §11 trennen veröffentlichten Host-Port (Loopback) und Bindung im Container (alle Adressen); Zugang anderer Container auf Port 22 aus der Bindung gefolgert, nicht ausgeführt (§13) |
| R2 F-2 MEDIUM Zugangsdaten | behoben im Dokument: §7 und §11 (neuer Punkt) nennen Ablage im Container außerhalb der Volumes, Lesbarkeit mit `sudo`, nicht eingeschränkten Netz-Ausgang, Entfallen mit dem Container; wo `claude` ablegt, ist nicht geprüft (§13). `u-boot down` entfernt Container und Netz (gemessen); Home und Anmeldung: gefolgert. Frage an den Nutzer unten |
| R2 F-3 LOW Lebensdauer | behoben: §4 (Zeile Stoppen), §6 (`--global`), §13 |
| R2 F-4 LOW Klon-Größe | behoben: §11 nennt 37,14 MB (`docker system df -v`), 47 MB mit `.git` und 18 MB ohne (`du -sh`); „28“ gestrichen |
| R2 F-5 LOW Cache | behoben: §8 nennt die offene Bedingung, §13 |
| R2 F-6 LOW Socket, Arbeitsbaum | behoben: §11 (Socket keine Grenze; Mount zeigt Arbeitsbaum), §6 |
| R2 F-7 LOW Voraussetzungen | behoben: §2 (Version 0.9.0, `ssh-keygen`/`scp`, Schreibrecht, Netz, Schlüsselpaar vorab; fehlendes Paar nicht geprüft) |
| R2 F-8 INFO fremde Kennungen | behoben: ein Satz in §3.1 des Dokuments, Blöcke unverändert; Plan §3 nachgezogen (V-4) |
| R2 F-9 INFO Lieferkette | behoben: §3.3 nennt Tag ohne Digest, `curl … \| bash -` als `root`, Herkunft nicht geprüft |
| V-1 MEDIUM Beleg der Gates | behoben: Beleg im Absatz „Läufe“ unten (Endstand `22f34d82`) |
| V-2 MEDIUM Form der Ausgänge | an den Planner: Hinweis in §6; Ausgang wird dort nicht umgeschrieben |
| V-3 LOW Meldeschwelle | **Meldung:** Größe 441 (Einfügungen 441, Löschungen 0; `git diff --numstat bda6efef..6ca83617 -- docs/developer`), über der Meldeschwelle 400, unter der Marke 500; vor der Nacharbeit 407 (`d725a22c`), 385 (`531f357e`). Folge benennt der Planner bei der Closure |
| V-4 INFO Kennungen | behoben: Plan §3 |
| V-5 INFO Zeiten | behoben: §13 des Dokuments nennt „Zeiten des Versuchs-Containers nicht wiederholt“ |

**Mutation je Aussage der Nacharbeit** (Dokument, kein Code): die neuen Sätze sind Aussagen über den
Versuchs-Container; ein Gegenbeispiel ist die Probe, die sie widerlegte: `sshd` lauscht an Loopback ↔ `ss -ltn`
zeigt `0.0.0.0:22`; Klon 28 MB ↔ `du` 47 MB/18 MB; Home im Volume ↔ `findmnt` (vier Mounts, `/home/vscode` im Overlay).
Nicht gesehen: das Verschwinden von `~/.gitconfig` nach `down`/`up` und der Zugang eines zweiten Containers; beide
tragen im Dokument „gefolgert“ und stehen in §13.

**Frage an den Nutzer (aus R2 F-2, keine Entscheidung des Implementers):** Braucht der Betrieb Schritte für die
Zugangsdaten im Container (zum Beispiel Abmelden vor `u-boot down`, ein Netz-Ausgang, der nur bestimmte Ziele
erlaubt, Verzicht auf die Claude-Anmeldung im Container)? Das Dokument erfindet keine; es nennt nur die Lage.

**Läufe (Beleg V-1, 2026-10-17).** Endstand `22f34d82`, Arbeitsbaum sauber: `make docs-check` (974 Dateien, 0 Befunde) und
`make kopf-check` grün; voller `make gates` grün, Ende `EXIT 0` (Log `gates4.log`, 08:47). Zwei volle Läufe davor endeten
rot an Gegenproben, die keine Änderung dieses Slice berühren (nur `docs/developer/` und der Plan): `lint-gegenprobe`
(07:49, mehrere Fälle der Punkte 8 und 9 rot, Ausgang 2) und `test-postgres-versionen-gegenprobe` (08:26, Mutation
`p7-int-treiber`: zusätzlich zu den erwarteten Fällen war `p7-sigTERM-gruppe` rot, Ausgang 2). Beide Gegenproben liefen im
folgenden Lauf ohne Änderung des Codes grün; als Flake gewertet, Ursache nicht ermittelt (zeitabhängige Signal- und
Last-Fälle, gefolgert, nicht gemessen).

### Portabilität der Pfade (2026-10-17, Wunsch des Nutzers)

Das Dokument trug `/Development/…` in neun Zeilen (38, 48, 54, 191, 192, 213, 242, 341, 351). Jetzt: `PGR_REPO`
(Host-Repo) und `PGR_DEV` (Verzeichnis `pgr-dev`) als `export` in §3; `mkdir "$PGR_DEV"`, `ssh`, `GIT_SSH_COMMAND`, `scp` mit
`$PGR_DEV`; `compose.yaml` mit `${PGR_REPO:?PGR_REPO setzen}:/host-repo:ro` und `./ssh/id_ed25519.pub`. `grep /Development`
im Dokument: keine Treffer. Der Messort `/Development/pgr-dev` steht nur hier.

| Aussage | Probe (Kopie des Blocks im Scratchpad, Datei im Versuchs-Container unverändert) | Beobachtung |
|---|---|---|
| Variable aus der Umgebung ersetzt | `PGR_REPO=/x/y docker compose config` | Quelle `/x/y`; Schlüssel aufgelöst relativ zum Verzeichnis der Datei |
| fehlende Variable bricht ab | `docker compose config` ohne Variable | `required variable PGR_REPO is missing a value: PGR_REPO setzen` |
| `.env` im Verzeichnis ersetzt | `.env` mit `PGR_REPO=/e/f` | Quelle `/e/f` |
| `.env` und relativer Pfad mit `-f` aus anderem Verzeichnis | `docker compose -f <Kopie>/compose.yaml config` aus `/` | Quelle `/e/f`, Schlüssel unter dem Verzeichnis der Datei |

**Unterschied zur getesteten Datei:** `compose.yaml` im Dokument ist an den zwei Zeilen der Mounts nicht mehr gleich der
Datei, mit der die Proben liefen; das Dokument sagt es in §3.5. **Nicht geprüft** (§13 des Dokuments): Start des Containers
mit dem neuen Block; `u-boot up` mit `PGR_REPO` nur aus der Umgebung (laut Hilfetext ruft es `docker compose up -d`, die
Weitergabe der Umgebung ist nicht gemessen). Kein Container gestartet. **Größe** nach dem Commit `bf07eea6`
(`git diff --numstat bda6efef..bf07eea6 -- docs/developer`): 466 Einfügungen, 0 Löschungen, Summe 466; über der Meldeschwelle
400, unter der Marke 500. `make docs-check` und `make kopf-check` siehe Bericht.

### Mehrere Container (2026-10-17, Zusatzwunsch des Nutzers)

Probe wie oben, Kopie im Scratchpad, kein Container, kein Docker-Rest (`docker ps -a` danach nur `pgr-dev` des Versuchs).

| Aussage | Probe | Beobachtung |
|---|---|---|
| Volumes folgen dem Projektnamen | `docker compose config` ohne `.env` / mit `COMPOSE_PROJECT_NAME=zwei` / `-p drei` | `pgr-dev_pgr-dev-docker` / `zwei_pgr-dev-docker` / Projekt `drei` |
| `container_name`, `hostname`, `image`, Port, `cpus`, `mem_limit` waren fest | vor der Änderung | `pgr-dev`, `pgr-dev`, `pgr-dev:latest`, `2222`, 8, 12 GiB auch bei `COMPOSE_PROJECT_NAME=zwei` |
| Variablen aus `.env` wirken | `.env` mit `COMPOSE_PROJECT_NAME=pgr-zwei`, `PGR_NAME=pgr-zwei`, `PGR_SSH_PORT=2223`, `PGR_CPUS=4`, `PGR_MEM=6g` | Projekt, Container, Hostname, Image `pgr-zwei`; Volumes `pgr-zwei_…`; Port `2223`; 4; 6442450944 |
| Vorgaben ohne `.env` | nur `PGR_REPO` gesetzt | `pgr-dev`, `2222`, 8, 12 GiB |
| Netz | beide Läufe | `pgr-dev-default` fest (Kopf von `u-boot`, nicht geändert) |

**Nicht geprüft** (im Dokument §13): Start einer zweiten Instanz und des geänderten Blocks; Wirkung des gemeinsamen Netzes
bei `up`/`down` der Instanzen; `u-boot up` mit Variablen nur aus der Umgebung. Das Dokument nennt diese Lücken in §10 und §13.
Block `compose.yaml` im Dokument weicht von der Datei im Messort nun in neun Zeilen ab (Mounts, Image, Name, Host-Name, Port, Grenzen).
**Größe** (`git diff --numstat bda6efef..e397805b -- docs/developer`): 498 Einfügungen, 0 Löschungen, Summe 498; über der
Meldeschwelle 400, unter der Marke 500 um 2; jede weitere Zeile im Dokument überschreitet die Marke und führt nach §4 zurück zu
`next` (Meldung an den Planner). `make docs-check` (974 Dateien, 0 Befunde) und `make kopf-check` (0) liefen nach `e397805b`
(kein `make gates` aktiv).

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
