# Slice slice-harness-devcontainer-einrichten: Make-Ziel `container-einrichten` richtet den Dev-Container ein und startet ihn (Gerüst, Schlüsselpaar, Start, Bare-Repo, Einlieferung, Klon; `docs/developer/devcontainer.md`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein Abnahmeszenario und
kein Meilenstein hängt an ihm. Angelegt am 2026-10-17 auf Wunsch des Nutzers vom selben Tag; am 2026-10-18 geschnitten (Entscheidung des Nutzers vom 2026-10-18: Der Zugang von Claude Code im Container wird nicht gebraucht, kein Login, die Hostsitzung treibt die Arbeit per `ssh` und Befehl im Container). Der Zugang liegt seitdem in `slice-harness-devcontainer-einrichten-zugang`; Gerüst und Start gehören zu diesem Slice.

**Bezug:** [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität) nur insoweit, als die Host-Voraussetzungen des Repos berührt sind (neu: `ssh`, Entscheidung des Nutzers, §1; dazu `u-boot`, `devcontainer`-CLI, `node`, `ssh-keygen`, noch nicht entschieden, §1); kein Anspruch des Produkts.
Vorgänger: `slice-harness-devcontainer-austausch-repo` und `slice-harness-devcontainer-doku` (beide `done/`); der Slice automatisiert deren Ablauf. Weiterer Vorgänger seit dem 2026-10-17: `slice-harness-devcontainer-u-boot-0-10` (`done/`; stellt den Container auf `u-boot` v0.10.0 um, gebaut und gestartet mit der `devcontainer`-CLI statt mit `u-boot up`, und misst `make gates` im Container: grün, 1701 s); dieser Slice baut auf dessen Ergebnis auf (§4). Folge-Slice: `slice-harness-devcontainer-einrichten-zugang` (Zugang von Claude Code; kein Teil dieses Slice, §1). Keine aktive ADR; der Architect legt eine an (§6).

**Berührte Spec-Stellen:** neuer Vertrag eines Harness-Werkzeugs in `harness/vertraege/` (Nummer vergibt der Architect); Host-Voraussetzungen: [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-17, geschnitten 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Make-Ziel `container-einrichten` am Host erzeugt in einem Aufruf den Dev-Container nach `docs/developer/devcontainer.md` §3 bis §6 und §9 und richtet ihn ein: Gerüst mit `u-boot` (`$PGR_DEV`, `u-boot.yaml`, die erzeugten Dateien unter `.devcontainer/`), Schlüsselpaar, `mk-override.sh`, Start mit `devcontainer up`, Bare-Repo `/workspaces/<Projektname aus u-boot.yaml>/austausch.git`, Einlieferung von `main`, Arbeitsklon `/workspaces/<Projektname>/pgwire-recorder`, Identität des Klons aus der Git-Konfiguration des Host-Repos, `make hooks-install` und `make traeger-fetch` im Klon, Zweig für den Slice. Dazu Stopp und Neustart (Dokument §4). Arbeit im Dev-Container wird damit der Standardweg, und `main` am Host bleibt für den Stop-Hook frei.

**Arbeitsweise (Entscheidung des Nutzers vom 2026-10-18):** Die Hostsitzung treibt die Arbeit im Container per `ssh` und Befehl; Claude Code läuft nicht im Container, und das Ziel löst keinen Zugang auf (Abgrenzung unten).

**Sendung von `slice-harness-devcontainer-u-boot-0-10` („Arbeit im Container“, eingetragen vom Planner am 2026-10-17, Entscheidung des Nutzers vom selben Tag):** Lang laufende Slices laufen im Container, damit `main` am Host nicht über Stunden durch einen Implementer-Lauf gesperrt ist (`BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`). Dieser Slice trägt davon den Einrichtungsweg (das Ziel) und im Dokument den Ablauf eines Implementer-Laufs im Container-Klon: Arbeit im Klon, `make gates` dort, Ergebnis per `git fetch container` am Host. Der Geber liefert das Gerüst auf `u-boot` v0.10.0, die gemessene Dauer von `make gates` im Container und die Probe von `claude`; dieser Slice baut auf diesen Befunden auf und misst sie nicht neu. Hat der Geber die Messung nicht grün belegt, startet dieser Slice nicht (§4).

**Annahme (Entscheidung des Nutzers vom 2026-10-17, „ja“): ein Slice pro Container-Sitzung.** Statt einer Welle pro Sitzung gilt: Ein Slice läuft in einer Sitzung, danach Neustart (`docker stop <Name>`, dann `devcontainer up …`; `u-boot up` startet das gemessene Gerüst nicht, `no service selected`). Gründe: Die Volumes `<Projektname>-docker-<ID>` und `<Projektname>-workspace-<ID>` bleiben über den Neustart erhalten (Cache, Bare-Repo, Klon), der Neustart ist billig (1,7 bis 1,8 s bis `ssh` und `dockerd` antworten, gemessen von `slice-harness-devcontainer-u-boot-0-10`); Prozesse und Zustand des Containers beginnen neu. Je Slice ein Zweig aus `main`, das Ergebnis kommt per `git fetch container` an den Host. Das Ziel ist deshalb der **Start-Ablauf einer Sitzung** und wiederholbar (Randform 21): es läuft bei jedem Neustart, nicht einmal je Container, und baut ein vorhandenes Gerüst nicht neu; die Namen der Volumes tragen die `devcontainerId` (`<ID>`) und ändern sich, wenn das Verzeichnis des Gerüsts wechselt. Ein Slice, der länger als rund 8 Stunden läuft, braucht einen Neustart mitten in der Arbeit; das ist kein Fehler des Werkzeugs (§6).

Start-Ablauf des Ziels, in dieser Reihenfolge (Nutzerwunsch 2026-10-17, am 2026-10-18 um Gerüst und Start erweitert, ohne Zugang):

1. Eingaben und Voraussetzungen am Host prüfen, bevor etwas geändert wird: `PGR_DEV` (Randform 18), Werkzeuge `u-boot`, `devcontainer`, `node`, `ssh`, `ssh-keygen`, `docker`, `git` (Randform 17), Port `PGR_SSH_PORT` (belegt: Randform 20), Identität des Host-Repos (Randform 5).
2. Gerüst anlegen, falls es fehlt (Randformen 19, 22, 23): `$PGR_DEV`, `u-boot init … --devcontainer --sandbox --no-git`, die `u-boot config set`-Zeilen aus Dokument §3.1, `u-boot generate devcontainer`, Schlüsselpaar `ssh/id_ed25519` (Dokument §3.2), `mk-override.sh` (Dokument §3.5).
3. Container starten (Randformen 21, 24): `mk-override.sh`, `devcontainer up --workspace-folder … --override-config …`, warten bis `ssh` antwortet (Zeitgrenze).
4. Remote und Erreichbarkeit prüfen: Container per `ssh` erreichbar, Remote `container` zeigt auf die erwartete URL.
5. Bare-Repo anlegen, falls es fehlt: `git init --bare --initial-branch=main /workspaces/<Projektname>/austausch.git` im Container, per `ssh`.
6. `main` einliefern: Remote `container` (`ssh://vscode@127.0.0.1:${PGR_SSH_PORT:-2222}/workspaces/<Projektname>/austausch.git`) und `core.sshCommand` im Host-Repo, dann `git push container main`.
7. Klon `/workspaces/<Projektname>/pgwire-recorder` anlegen, falls er fehlt, sonst aktualisieren (Randform 2).
8. Name und E-Mail des Klons **lokal** aus der Git-Konfiguration des Host-Repos setzen: nie `--global`, nie ein im Dokument oder Skript fest stehender Name; fehlt am Host einer der beiden Werte, bricht das Ziel ab, bevor es etwas ändert (Schritt 1; Reihenfolge der Gründe: Randform 14); Werte mit Anführungszeichen oder Apostroph kommen unverändert an.
9. `make hooks-install` und `make traeger-fetch` im Klon.
10. Zweig für den Slice aus `main` im Klon anlegen (Randform 16).

Stopp und Neustart (`docker stop`, dann Schritt 3 und folgende) sind Dokument §4; ob dafür ein zweites Ziel entsteht oder der zweite Aufruf von `container-einrichten` genügt (Randform 21), entscheidet der Architect.

**Bindung an `slice-harness-devcontainer-einrichten-zugang` (Sendung, eingetragen vom Planner am 2026-10-18, `AGENTS.md` §3.13):** Das Skript liest, löst und übergibt keinen Zugang (weder `ANTHROPIC_API_KEY` noch `CLAUDE_CODE_OAUTH_TOKEN` noch eine Anmeldedatei) und prüft keinen; sein Ablauf lässt eine Stelle vor Schritt 1 frei, an der der Nehmer die Auflösung einfügt (Randform 25, Architect). `mk-override.sh` schreibt das Ziel, wie Dokument §3.5 es zeigt, auch die Zeile mit dem Verweis `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}`; was aus dem Verweis wird, entscheidet der Nehmer.

**Sendung „Quellenauflösung“ von `slice-harness-devcontainer-u-boot-0-10` (eingetragen vom Planner am 2026-10-17) geht am 2026-10-18 weiter an `slice-harness-devcontainer-einrichten-zugang`** (Entscheidung des Nutzers vom 2026-10-18). Der Geber nennt als Nehmer weiter diesen Slice; die Adresse bleibt einen Hop länger: Dieser Slice nennt den Nehmer, im Geber (`done/`) wird kein Zeiger umgeschrieben. Randform 15, Quellenauflösung, Übergabe über `devcontainer exec`, `remoteEnv` und beide Variablen stehen dort; dieser Slice führt sie nicht.

**Entscheidung des Nutzers (2026-10-17): `ssh` ist neue Voraussetzung am Host.** Auf die Frage, ob `ssh` als Voraussetzung für dieses Ziel in Ordnung ist, antwortete der Nutzer „ja“. Die Voraussetzungen des Repos sind bisher `git`, `docker` und `make` (`AGENTS.md` §3.1, [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität)). Der Architect trägt die Entscheidung in die Spezifikation des Vertrags und in die ADR ein (§3); `ssh` gilt nur für dieses Ziel, kein Gate und kein anderes Ziel braucht es. **Mit dem Start im Ziel kommen `u-boot`, die `devcontainer`-CLI, `node` (braucht die CLI und `mk-override.sh`) und `ssh-keygen` dazu; der Nutzer hat diese Voraussetzungen nicht entschieden.** Der Architect legt sie vor dem Code vor (§6, Randform 17), der Nutzer entscheidet; bis dahin ist es eine Empfehlung des Planners, kein Plan-Inhalt.

**Kein Gate**, wie `hooks-install`: Das Ziel steht im Abschnitt *Werkzeuge* von `harness/README.md`, nicht in `GATE_CHECKS`. Seine Gegenprobe ist ein Gate (§2).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Zugang von Claude Code im Container: keine Quellenauflösung, keine Randform 15, kein `remoteEnv`-Weg über `devcontainer exec`, keine Variable `ANTHROPIC_API_KEY` oder `CLAUDE_CODE_OAUTH_TOKEN` im Skript, kein Dokument §7 (Klasse 1: `slice-harness-devcontainer-einrichten-zugang`, dessen §1 diese Sendung annimmt; Entscheidung des Nutzers vom 2026-10-18). Das Ziel richtet den Container für Arbeit per `ssh` ein, nicht für `claude`.
- Kein Umbau des Gerüsts auf eine andere Version von `u-boot` und keine Messung von `make gates` im Container: Das liefert der Vorgänger `slice-harness-devcontainer-u-boot-0-10` (`done/`; Gerüst auf `u-boot` v0.10.0, Messung grün). Das Ziel automatisiert dessen gemessenen Ablauf (Dokument §3 bis §5) und setzt den Inhalt von `u-boot.yaml` und `mk-override.sh` aus dem Dokument um; trägt der Inhalt den Ablauf nicht, hält der Implementer an (§6).
- Keine Entscheidung zwischen `podman` und `docker` als `nestedRuntime` (`BEO-REPO/nested-runtime-podman-oder-docker-offen`): Das Ziel schreibt den Wert, den das Dokument nennt (`docker`), und entscheidet nichts darüber.
- Keine Entscheidung zum Hostschlüssel (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`): Das Ziel nennt, was `ssh` ausgibt, und setzt keine Einstellung (`StrictHostKeyChecking`, `UserKnownHostsFile`) ohne Entscheidung des Nutzers (§6).
- Kein Löschen von Container, Volumes oder Image (Dokument §4, *Container entfernen*, *Alles löschen*, *Image löschen*): Das ist Handlung des Nutzers; das Ziel entfernt nichts, die Gegenprobe nur Objekte ihres Präfixes per Name.
- Keine Docker-Daemon-Probe außer dem Start: Das Ziel fährt `devcontainer up` (das ruft `docker`) und prüft sonst den Container per `ssh`, nicht per `docker`. Die Gegenprobe darf den Daemon nur nutzen, wenn kein `make gates` läuft (§3, *Docker-Last*).
- Kein Produktcode (Schicht-Abgrenzung): Es ändern sich Werkzeug, Gegenprobe und Dokumentation, nichts unter `cmd/` und `internal/`.
- Kein Eintrag in der Roadmap: Der Slice ist wellenlos. Die Release-Slices und die Closure von `welle-v1-abschluss` hängen nicht an ihm (Entscheidung des Nutzers vom 2026-10-17).

## 2. Definition of Done

Liefer-Punkte (drei):

- [ ] **Ziel `container-einrichten`** (Fragment unter `harness/mk/` und Skript unter `tools/harness/`): erfüllt den Vertrag, den der Architect in `harness/vertraege/` festlegt (Ablauf aus §1, Randformen aus §6); kein Gate; in `harness/README.md` §Werkzeuge genannt.
- [ ] **Gegenprobe** `container-einrichten-gegenprobe` (Gate, hängt an `GATE_CHECKS`, steht in der Ausnahmeliste der schnellen Kette): gegen einen Scratch-Container (Präfix, Entfernen nur per Name, `docker rm -fv`, nichts Fremdes, nie der Container, in dem der Implementer arbeitet) je Randform aus §6 ein Fall und je Zusage und je Element einer Aufzählung eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die Gegenprobe steht im Gate-Index.
- [ ] **Dokument** `docs/developer/devcontainer.md` §3 bis §6 und §9 nennen `make container-einrichten` als Weg und führen die Einzelbefehle nur, soweit das Dokument sie als Erklärung braucht; dazu der Ablauf eines Implementer-Laufs im Container-Klon (Sendung von `slice-harness-devcontainer-u-boot-0-10`, §1): Arbeit im Klon per `ssh`, `make gates` dort mit der Dauer, die jener Slice misst, Ergebnis per `git fetch container`. §7 bleibt unberührt (gehört `slice-harness-devcontainer-einrichten-zugang`). Jede Aussage ist durch einen Lauf belegt oder steht in §13.

Gates und Closure:

- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `tools/`, `harness/` und das Makefile-Fragment ändern); Ausführung nicht gleichzeitig mit den Läufen der Scratch-Container (§3).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-<NNN>-container-einrichten.md` | neu (Architect) | Vertrag des Werkzeugs: Ablauf, Randformen aus §6, Host-Voraussetzungen (`ssh`, und je nach Entscheidung `u-boot`, `devcontainer`, `node`, `ssh-keygen`); Nummer vergibt der Architect |
| `docs/plan/adr/` (eine ADR) | neu (Architect) | Entscheidung: Ort des Ziels (Repo oder Host-Werkzeug), Host-Voraussetzungen (Entscheidung des Nutzers für `ssh` vom 2026-10-17, übrige offen); `Schärft:` zeigt auf den Vertrag (`AGENTS.md` §3.8) |
| `harness/mk/container-einrichten.mk` | neu | Ziel `container-einrichten` (kein Gate), Gegenprobe hängt an `GATE_CHECKS` |
| `tools/harness/container-einrichten.sh` | neu | Start-Ablauf aus §1 (Prüfen, Gerüst, Start, Einrichtung, Zweig); Eingaben aus der Umgebung (`PGR_DEV`, `PGR_SSH_PORT`, `PGR_NAME`, Remote-Name), Fehlermeldungen mit Schritt; keine Zugangsvariable (Bindung, §1) |
| `tools/harness/container-einrichten-gegenprobe.sh`, `harness/mk/container-einrichten-gegenprobe.mk` | neu | Gegenprobe gegen Scratch-Container; Fälle und Mutationen aus §6 |
| `harness/mk/gates-schnell.mk` | update | Gegenprobe in die Ausnahmeliste der schnellen Kette (Docker; Pflege nach `SPEC-052`, Architect prüft Punkt der Gate-Kette) |
| `harness/README.md` | update | Abschnitt *Werkzeuge* (`container-einrichten`) und *Sensors* (Gegenprobe) |
| `docs/developer/devcontainer.md` | update | §3 bis §6 und §9 nennen das Ziel; Ablauf eines Implementer-Laufs; §13 für Ungeprüftes |

- Das Skript trennt Prüfen und Handeln: erst alle Vorbedingungen am Host (Werkzeuge, Eingaben, Port, Identität, Schlüssel lesbar), dann Gerüst und Start, dann die Schritte im Container. Die Reihenfolge der Gründe, aus denen es abbricht, legt der Architect fest (Randform 14).
- Die Identität reist ohne Auswertung durch eine zweite Shell, soweit der Architect keine andere Form festlegt: Der Implementer entscheidet die Quoting-Form nicht, sondern setzt die des Vertrags um (§6, Sonderzeichen).
- **Docker-Last.** Die Gegenprobe startet Scratch-Container (Präfix, Entfernen nur per Name, `docker rm -fv` und Prüfung, dass keines der Volumes bleibt, wie `SPEC-057`). Sie läuft nie gleichzeitig mit `make gates`; die Vorbedingung `pgrep -fa "^make gates"` leer ist Teil der Gegenprobe (ob als Abbruch oder als Warten, legt der Vertrag fest). Ihre Laufzeit ist zu messen (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`).
- **Arbeit des Implementers im Dev-Container `pgr-dev2`** (Entscheidung des Nutzers vom 2026-10-18; Host-Verzeichnis `/Development/pgr-dev2`, Port 2223, kein Claude Code im Container). Der Implementer läuft in der Hostsitzung und treibt den Klon im Container per `ssh`:
  1. *Instanz.* `PGR_DEV=/Development/pgr-dev2`, `PGR_SSH_PORT=2223`, `PGR_NAME=pgr-dev2`; der Pfad im Container ist `/workspaces/<Projektname aus /Development/pgr-dev2/u-boot.yaml>/` (Annahme `pgr-dev2`, der Implementer prüft den Wert mit `grep` in der Datei und `ls /workspaces` im Container). Ob der Container läuft, ist ungemessen: `ssh -i /Development/pgr-dev2/ssh/id_ed25519 -o IdentitiesOnly=yes -o ConnectTimeout=10 -p 2223 vscode@127.0.0.1 true`; die Antwort entscheidet, ob er startet (Dokument §4: `./mk-override.sh` mit den Variablen, dann `devcontainer up …`) oder weitermacht. Der erste Kontakt zeigt einen unbekannten Hostschlüssel: Der Implementer setzt keine Einstellung (`StrictHostKeyChecking`, `UserKnownHostsFile`) und hält an, bis der Nutzer den Eintrag gesetzt oder entschieden hat (Abgrenzung §1).
  2. *Bare-Repo, Klon, Remote.* Das Ziel gibt es vor diesem Slice nicht; für `pgr-dev2` richtet der Implementer Bare-Repo, Einlieferung und Klon nach Dokument §6 und §9 mit den Werten dieser Instanz von Hand ein (Remote am Host `container2`, Dokument §10; Identität lokal aus der Git-Konfiguration des Host-Repos, nie `--global`, nie fest im Befehl). Nach der Lieferung des Ziels ruft er es für diese Instanz auf und hält so die Probe des Ziels am eigenen Arbeitsplatz fest.
  3. *Arbeiten.* Zweig `slice-harness-devcontainer-einrichten` im Klon `/workspaces/<Projektname>/pgwire-recorder`, nicht `main`. Jeder Befehl: `ssh -i /Development/pgr-dev2/ssh/id_ed25519 -o IdentitiesOnly=yes -p 2223 vscode@127.0.0.1 'cd /workspaces/<Projektname>/pgwire-recorder && <Befehl>'`. Dateien kommen als Patch über `stdin` (`git apply` im Klon) oder als Inhalt über `stdin` (`cat > Datei`), nicht über eine eingehängte Arbeitskopie (der Container sieht den Arbeitsbaum des Hosts nicht, Dokument §7). Commits im Klon, jeder mit der Kennung `slice-harness-devcontainer-einrichten` im Text (Hook; `make hooks-install` und `make traeger-fetch` vorher im Klon).
  4. *Gates.* `make gates-schnell` und `make test` im Klon; der volle Lauf `make gates` dauert 28 min 21 s (Dokument §8, ungeprüft für diesen Stand) und läuft im Hintergrund (`setsid`, Log in `/tmp` im Container, Exit-Code in eine Datei), die Abfrage per `ssh`; `pgrep -fa "^make gates"` im Container ist vor jedem Lauf der Gegenprobe leer (Docker-Last).
  5. *Was nur der Host kann.* Die Gegenprobe braucht zum Erzeugen eines Scratch-Gerüsts `u-boot`, `devcontainer` und `node`. Ob sie im Container von `pgr-dev2` vorliegen, ist ungemessen: Der Implementer prüft `command -v u-boot devcontainer node ssh ssh-keygen` im Container zuerst und nennt das Ergebnis im Bericht. Fehlen sie, läuft nur der Teil der Gegenprobe, der sie braucht, am Host mit Scratch-Verzeichnis, Container-Namen und Volumes im Präfix `pgr-ceg-` (nie `pgr-dev2` anfassen, nie stoppen oder entfernen); der Bericht sagt, was im Container und was am Host lief.
  6. *Ergebnis zum Host.* Im Klon `git push origin slice-harness-devcontainer-einrichten`; am Host `git fetch container2` und der Zweig `container2/slice-harness-devcontainer-einrichten` liegt zur Prüfung vor (Dokument §9). Der Inhalt stammt aus einem Container, in dem `vscode` `root` werden kann: Review und Verifikation sehen die Commits durch, bevor sie nach `main` gehen (Dokument §9). Der Übergang `next/` nach `in-progress/` ist ein Commit des Nutzers auf `main` vor der Arbeit und wird vor dem ersten Klon eingeliefert (`git push container2 main`).
  7. *Größenprobe.* Die Metrik von §4 misst der Implementer im Klon mit `git diff --numstat <Basis>..<letzter Commit>`.
- Zählung nach `AGENTS.md` §3.13 siehe §8.

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-harness-devcontainer-u-boot-0-10` liegt in `done/` (erfüllt), und seine Messung von `make gates` im Container ist grün (erfüllt: zwei Läufe zu je 1701 s). Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. Reihenfolge laut Entscheidung des Nutzers vom 2026-10-17: Die Release-Slices und die Closure von `welle-v1-abschluss` verlangen nur `slice-harness-devcontainer-u-boot-0-10` in `done/`, diesen Slice nicht. Vorbedingung des Starts außerdem: Vertrag und ADR des Architect liegen vor (`Accepted` ist die Entscheidung des Nutzers, keine Annahme des Planners), und der Nutzer hat die Host-Voraussetzungen entschieden (§1).

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; der Vertrag und die ADR (Commits des Architect vor dem Start) liegen vor der Basis und zählen nicht.

**Schätzung nach dem Schnitt vom 2026-10-18 (Planner; Metrik wie oben, Gegenprobe mit dem Doppelten des Produktcodes, Faktor aus `AGENTS.md` §3.14):**

| Teil | Zeilen |
|---|---|
| Skript: Prüfen der Eingaben und Voraussetzungen | 30 bis 50 |
| Skript: Gerüst (Verzeichnis, `u-boot`-Aufrufe, Schlüsselpaar, `mk-override.sh`) | 60 bis 90 |
| Skript: Start (`devcontainer up`, Warten auf `ssh`) | 30 bis 50 |
| Skript: Einrichtung (Remote, Bare-Repo, Einlieferung, Klon, Identität, Hook, Träger, Zweig) | 115 bis 170 |
| Fragmente | 50 bis 60 |
| Gegenprobe (Doppeltes von Skript und Fragmenten, 285 bis 420) | 570 bis 840 |
| Dokument (§3 bis §6, §9 gekürzt, Löschungen zählen, dazu Ablauf des Implementer-Laufs, Index) | 140 bis 220 |
| **Zusammen** | **995 bis 1480** |

Meldeschwelle 700, **Marke 900**. Die Schätzung liegt **mit dem ganzen Spektrum über der Marke**; Vergleich: der Schnitt ohne Zugang, aber ohne Gerüst und Start, lag bei 550 bis 820. Praxis vom 2026-10-16: geschnitten wird, wenn die Schätzung die Marke überschreitet. **Vorschlag des Planners (Empfehlung, Entscheidung des Nutzers offen, solange der Slice in `open/` liegt):** Gerüst und Start (Schritte 1 bis 3 und Stopp/Neustart, Dokument §3 bis §5) in einen Slice `slice-harness-devcontainer-einrichten-geruest` (Name vorläufig, **nicht angelegt**, daher noch keine Adresse im Sinn von `AGENTS.md` §3.13), der vor diesem Slice schließt und ein laufendes Gerüst liefert; dieser Slice behält die Einrichtung im laufenden Container (Schritte 4 bis 10, Dokument §6 und §9, mit dem Ablauf des Implementer-Laufs). Der Schnitt nach Lieferwert hält: Der Gerüst-Slice liefert einen Container, der per `ssh` erreichbar ist; dieser liefert Bare-Repo, Klon und Zweig darin. Bis zur Entscheidung gilt der Umfang aus §1, und die Rückführung unten bleibt.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder der Vertrag verlangt mehr als drei Liefer-Punkte oder mehr als zwei Schichten (zum Beispiel eine Änderung an `.devcontainer/Dockerfile` oder an dem Inhalt von `mk-override.sh` über das Dokument hinaus), oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner (Vorschlag: der Schnitt oben).
- `in-progress` → `open` (blockiert): Der Ablauf trägt nicht ohne Eingriff in das Image (zum Beispiel fehlt `git` im Container, Rechte an `/workspaces`), oder `ssh` am Host ist nicht zugesagt (Entscheidung des Nutzers wird zurückgenommen), oder die Host-Voraussetzungen `u-boot` und `devcontainer` bleiben unentschieden (Nutzer), oder das Gerüst, das `slice-harness-devcontainer-u-boot-0-10` hinterlässt, trägt den Ablauf nicht (der Befund geht an den Planner).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün auf dem Endstand, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben. Das Mutationsprotokoll (Zusage · Mutation · roter Test, `AGENTS.md` §3.10) steht im Bericht des Implementers. Die Gegenprobe ist einmal gegen einen echten Scratch-Container grün gesehen, nicht nur gegen Attrappen.

## 6. Risiken und offene Punkte

Der Slice liefert einen **neuen Vertrag** (Werkzeug mit Zusagen); §3.12 gilt. Die Randformen stehen hier, entschieden werden sie **vor dem ersten Code-Commit** in `harness/vertraege/` (Architect); eine Form, die erst ein Commit mit Code oder Tests hier einträgt, ist nicht genannt. Der Implementer entscheidet keine still, er hält an und gibt sie dem Architect zurück. Vorläufiges Wissen aus dem Dokument steht dabei, nicht als Entscheidung.

**Entscheidungen des Architect vor dem Code** (Liste für die Prüfung vor dem ersten Code-Commit; die Randformen darunter sind die Einzelfälle):

- **E1 Ort des Ziels.** Ein Make-Ziel im Repo (`harness/mk/`, `tools/harness/`, wie §3 es annimmt) oder ein Host-Werkzeug außerhalb des Repos (wie `u-boot` selbst)? Das Ziel schreibt Verzeichnisse außerhalb des Repos (`$PGR_DEV`) und ändert Konfiguration des Host-Repos (`remote`, `core.sshCommand`); er entscheidet, ob das in ein Repo-Makefile gehört. Die Antwort bestimmt Vertrag (`harness/vertraege/` oder nicht), Gate-Index, Gegenprobe und Größe.
- **E2 Host-Voraussetzungen.** `u-boot` (Version 0.10.0 oder höher?), `devcontainer`-CLI, `node` (nur durch CLI und `mk-override.sh`; ob `mk-override.sh` ohne `node` auskommt, ist eine Frage des Inhalts), `ssh-keygen`. `ssh` ist entschieden (§1), die übrigen nicht. Der Architect legt vor, der Nutzer entscheidet; eine Voraussetzung, die der Nutzer streicht, verlangt einen anderen Weg (Gerüst von Hand, Ziel nur Einrichtung) und damit den Schnitt aus §4.
- **E3 Randformen des Aufrufs** (Zeilen 17 bis 25): fehlende Variable, vorhandenes Verzeichnis, belegter Port, zweiter Aufruf, Werkzeug fehlt.
- **E4 Aufrufarten von `make`** (Randform des Treibers, `AGENTS.md` §3.12): `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` aus der Umgebung, auch in Kombination, dazu `-k` und `-j`.
- **E5 Reihenfolge der Gründe** (Randform 14, ergänzt um Voraussetzungen, Eingaben, Port, Gerüst, Start) und Exit-Codes.
- **E6 Stopp/Neustart** (§1): zweites Ziel oder zweiter Aufruf.
- **E7 Schnitt** (§4): ob geschnitten wird, bevor der Architect schätzt; der Architect bestätigt die Schätzung in der Metrik von §4 vor dem Code.

**Randformen des Ablaufs** (Spalte *Wo entschieden*: Fundstelle im Vertrag, vom Architect einzutragen):

| # | Randform | Vorwissen aus `docs/developer/devcontainer.md` (§3 bis §6, §9, §10) | Wo entschieden |
|---|---|---|---|
| 1 | Bare-Repo besteht schon | `git init --bare` auf vorhandenem Verzeichnis meldet `Reinitialized`, Exit 0; Frage: Wiederholung ist erlaubt oder bricht ab? Besteht es mit anderem `HEAD` (nicht `main`)? | offen |
| 2 | Klon besteht schon (leer oder nicht leer, anderer Remote) | `git clone` in nicht leeres Verzeichnis: Exit 128. Frage: überspringen, aktualisieren oder abbrechen; was geschieht mit uncommittetem Stand im Klon? | offen |
| 3 | Remote `container` besteht schon (gleiche URL, andere URL, anderer Port) | zweites `git remote add` bricht mit Exit 3 ab. Frage: bei gleicher URL weiter, bei anderer abbrechen oder setzen | offen |
| 4 | Container startet nicht / `ssh` nicht erreichbar (Port, Schlüssel, Zeitgrenze) | Der Start läuft im Ziel (Schritt 3). Frage: Meldung, Exit-Code, Zeitgrenze für das Warten und für die Verbindung (`ConnectTimeout`), was bleibt stehen (kein Rückbau) | offen |
| 5 | Identität am Host leer (nur Name, nur E-Mail, beide, nur Leerzeichen) | der Aufruf in §6 bricht ab, wenn eine fehlt. Frage: Reihenfolge gegen Randform 4 (welcher Grund zuerst; jetzt vor dem Gerüst, Schritt 1), Meldung nennt den fehlenden Wert | offen |
| 6 | Identität mit Sonderzeichen: Anführungszeichen `"`, Apostroph `'`, Leerzeichen, `$`, Backtick, Backslash, Zeilenumbruch, Nicht-UTF-8 | `printf %q` am Host; Gegenprobe am Wert nötig. Frage: erlaubte Menge, Quoting-Form (die Shell im Container ist die Anmelde-Shell von `vscode`), ob Zeilenumbruch abgelehnt wird | offen |
| 7 | Identität nur lokal, nie global | Zusage: `--global` und `--system` werden nie gesetzt; Frage: gilt auch, wenn das Host-Repo die Werte nur aus `--global` bezieht (`git config user.name` liest alle Ebenen) | offen |
| 8 | Anderer Port, mehrere Instanzen (`PGR_SSH_PORT`, zweiter Remote-Name) | Dokument §10: je Instanz eigener Port und Remote (`container2`); `PGR_SSH_PORT` liest nur `mk-override.sh` aus der Shell (gemessen), nicht eine `.env`, und der Pfad im Container trägt den Projektnamen (`/workspaces/<Projektname>/`). Frage: wie wählt das Ziel Remote-Name und Bare-Repo-Pfad je Instanz; `PGR_SSH_PORT` leer, nicht numerisch, außerhalb 1 bis 65535; Instanz im Pfad. Der Fall der Arbeit des Implementers: `pgr-dev2`, Port 2223 | offen |
| 9 | Hostschlüssel unbekannt, geändert (Neubau) | Warnung am Host bekannt; keine Einstellung festgelegt (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`). Befund von `slice-harness-devcontainer-u-boot-0-10`: Das Feature `sshd` erzeugt den Hostschlüssel im Build des Images; Neustart und neuer Container aus demselben Image ändern ihn nicht, jeder Neubau (`--build-no-cache`, Änderung von `egress.allow`) ändert den Fingerabdruck. Frage: Ziel setzt nichts und bricht ab, oder die Entscheidung des Nutzers liegt vor | offen (Entscheidung des Nutzers) |
| 10 | Netz fehlt bei `make traeger-fetch` im Klon | braucht Netz im Container; das Egress des Containers bleibt eingeschränkt (Entscheidung des Nutzers vom 2026-10-17), der Host des Klons muss in der `egress.allow`-Liste stehen; `slice-harness-devcontainer-u-boot-0-10` hat für `make gates` zusätzlich `gcr.io` und `pkg-containers.githubusercontent.com` gemessen und eingetragen, `make traeger-fetch` im Klon ist dort nicht gelaufen (Dokument §13), sein Host ist ungemessen; Frage: Exit-Code des Ziels, ob die vorigen Schritte stehen bleiben (ja, kein Rückbau) und wie die Wiederholung aussieht | offen |
| 11 | Abbruch mitten im Ablauf (`SIGINT`, `SIGTERM`, Verbindungsabbruch) und Wiederholung | kein Rückbau; Frage: jeder Schritt wiederholbar (Randformen 1 bis 3, 19, 21, 22), Meldung nennt den letzten Schritt; Abbruch mitten in `u-boot` oder `devcontainer up` | offen |
| 12 | `main` am Host fehlt, ist leer, steht nicht im Host-Repo oder nur auf anderem Zweig; Host-Repo im Zustand `detached HEAD` | Frage: einzuliefernder Zweig ist `main`, das Ziel hängt nicht am ausgecheckten Zweig; Meldung bei fehlendem `main` | offen |
| 13 | Push abgewiesen (Bare-Repo hat schon einen `main` mit anderer Historie) | `! [rejected]` bei nicht-fast-forward. Frage: kein `--force`, Meldung und Exit | offen |
| 14 | Ausgang mit Exit-Codes: Reihenfolge der Gründe, wenn mehrere zutreffen (Identität fehlt **und** Werkzeug fehlt **und** Port belegt **und** Container aus) | Reihenfolge der Gründe bindet der Vertrag (`AGENTS.md` §3.12, Stelle einer Nachricht) | offen |
| 15 | *(verschoben)* Quelle des Zugangs von Claude Code | steht in `slice-harness-devcontainer-einrichten-zugang`, Randform 15 dort; die Nummer bleibt vergeben | — |
| 16 | Zweig je Slice im Klon (Schritt 10) | Entscheidung des Nutzers vom 2026-10-17: je Slice ein Zweig aus `main`. Fragen: woher der Name (Parameter, Pflicht?), Zweig besteht schon (Wiederholung des Ziels nach Neustart: weiter auf dem Zweig oder Abbruch?), Klon hat uncommittete Änderungen, `main` im Klon steht hinter dem Bare-Repo | offen |
| 17 | Werkzeug am Host fehlt oder hat die falsche Version (`u-boot`, `devcontainer`, `node`, `ssh`, `ssh-keygen`, `docker`, `git`) | Dokument §2 nennt die Voraussetzungen des Dokuments; Frage: je Werkzeug ein eigener Fall und eine Meldung, die den Namen nennt; Version von `u-boot` (v0.10.0 im Dokument gemessen); vor jeder Änderung (Schritt 1) | offen (E2, Entscheidung des Nutzers) |
| 18 | `PGR_DEV` fehlt, ist leer, relativ, enthält Leerzeichen, liegt im Host-Repo, ist kein Verzeichnis | Dokument §3: absoluter Pfad außerhalb des Repos; Frage: Pflicht oder Vorgabe, Prüfung der Form, Pfad im Arbeitsbaum des Repos abgelehnt | offen |
| 19 | Verzeichnis `$PGR_DEV` besteht schon (leer; mit Gerüst; mit fremdem Inhalt; mit Gerüst anderer Version von `u-boot`) | Dokument §3.1 legt mit `mkdir` an, bricht bei vorhandenem Verzeichnis ab; Frage: leer erlaubt, Gerüst erkennen und weiter, fremder Inhalt abgelehnt; Dateien zwischen `BEGIN`/`END U-BOOT MANAGED BLOCK` nicht überschreiben | offen |
| 20 | Port `PGR_SSH_PORT` belegt (anderer Container, fremder Dienst), Port 2222 gilt als Vorgabe | Dokument §10: ein belegter Port ist nicht gelaufen (§13); Frage: Prüfung vor dem Start oder Meldung von `devcontainer up`, Exit, Unterschied belegt durch denselben Container (zweiter Aufruf) und durch einen fremden | offen |
| 21 | Zweiter Aufruf: Container läuft schon; Container gestoppt; Container mit dem Namen besteht aus anderem Gerüst; `PGR_NAME` kollidiert | Dokument §4: Stopp und Neustart messen 1,7 bis 1,8 s; `devcontainer up` auf laufendem Container: ungemessen; Frage: erster und zweiter Aufruf unterscheiden, Neustart ohne Neubau, Kollision mit fremdem Namen abgelehnt | offen |
| 22 | Schlüsselpaar `ssh/id_ed25519` besteht schon (mit/ohne Passphrase, nur ein Teil, falsche Rechte) | Dokument §3.2: `ssh-keygen -N ''`; überschreiben vernichtet den Zugang; Frage: nie überschreiben, Teil fehlt → Abbruch, Rechte 700/600 | offen |
| 23 | `u-boot` meldet Warnungen oder Fehler (Lockerungen des Profils, `doctor`-Warnung `no service selected`, Exit 12 von `u-boot up`) | Dokument §3.1, §12: Warnungen sind erwartet; Frage: welche Warnung und welcher Exit gelten als Erfolg, was bricht ab; die Ausgabe von `u-boot` erscheint unverändert | offen |
| 24 | Warten auf `ssh` nach `devcontainer up`: Zeitgrenze, `devcontainer up` meldet Erfolg, `sshd` antwortet nicht, `dockerd` nicht bereit | Dokument §4: Bereit ist der Container, wenn `ssh` antwortet; Erststart 65 s, Neustart 1,7 bis 1,8 s (Cache warm); ohne warmen Cache ungemessen; Frage: Zeitgrenze, Exit, ob `dockerd` Teil der Bereitschaft ist (`make traeger-fetch` braucht ihn nicht, `make gates` im Klon schon) | offen |
| 25 | Stelle des Zugangs im Ablauf (Bindung an `slice-harness-devcontainer-einrichten-zugang`) | Der Nehmer fügt die Auflösung vor Schritt 1 ein; Frage: wie das Skript die Stelle frei hält, ohne dass dieser Slice Zugangs-Code enthält (Funktion, Reihenfolge der Prüfungen); der Nehmer ändert dann dieselbe Datei | offen |

**Aufrufarten von `make`** (Randform des Treibers, `AGENTS.md` §3.12 *Aufrufarten des Treibers*): Das Ziel ändert Zustand außerhalb des Arbeitsbaums (Verzeichnisse am Host, Container, Remote, Konfiguration des Host-Repos). Je Flag eine Entscheidung im Vertrag und je Flag ein Fall und eine Mutation in der Gegenprobe: `-n` (Rezept mit `$(MAKE)` darf nichts ausführen: kein `u-boot`, kein `ssh`, kein `devcontainer up`, kein Push, keine Konfiguration), `-q`, `-t` (Ziel ist `.PHONY`, trägt kein Zeitstempel-Ziel), `-i` (Fehler dürfen den Ablauf nicht überlesen: kein „Einrichtung abgeschlossen“ nach einem Fehler), `MAKEFLAGS` aus der Umgebung, auch in Kombination und mit Vorrang. Dazu `-k` und `-j` (parallel laufende Schritte dürfen die Reihenfolge nicht brechen).

**Weitere Risiken:**

- **Neue Host-Abhängigkeiten.** `ssh` ist entschieden (Nutzer, 2026-10-17); `u-boot`, `devcontainer`-CLI, `node`, `ssh-keygen` nicht (§1, Randform 17); Eintrag in Vertrag und ADR durch den Architect, kein Plan-Eintrag in `AGENTS.md` §3.1 und nicht in Lastenheft oder Gate-Index, solange der Architect nichts anderes festlegt. Fehlt ein Werkzeug am Host, meldet das Ziel es vor jeder Änderung (eigener Fall je Werkzeug). — **Ausgang (Closure):** offen bis zur Closure; erwartet *entfallen* mit Verweis auf Vertrag und ADR.
- **Der Slice wächst über die Marke durch Gerüst und Start** (§4: 995 bis 1480 gegen Marke 900). — **Ausgang (Closure):** eingetreten → `slice-harness-devcontainer-einrichten-geruest` (Name vorläufig, Schnitt wie in §4, vom Planner vergeben), oder entfallen mit gemessener Summe unter 900.
- **Der volle Lauf von `make gates` im Container ist von `slice-harness-devcontainer-u-boot-0-10` gemessen** (grün, zwei Läufe zu je 1701 s, höchstens 5,42 GiB; Dokument §8); dieser Slice misst ihn nicht neu. Trägt der Lauf im Klon von `pgr-dev2` nicht (anderes Gerüst, andere Last), ist das ein Befund an den Planner. — **Ausgang (Closure):** entfallen, weil `slice-harness-devcontainer-u-boot-0-10` den Lauf misst und der Start dieses Slice seine grüne Messung verlangt (§4).
- **Gegenprobe stört einen anderen Lauf** (Docker-Daemon geteilt, `make gates` des Implementers eines anderen Slice) oder den Container, in dem der Implementer arbeitet (`pgr-dev2`), oder hinterlässt Volumes (`BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`). — **Ausgang (Closure):** entfallen mit Beleg, dass die Gegenprobe nur Objekte ihres Präfixes entfernt und nach Ende keines bleibt, oder eingetreten → Folge-Slice mit Kennung.
- **Gegenprobe wächst** (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`): zwei Zusagen je Aufzählungselement, dazu Gerüst und Start je Randform 17 bis 24. — **Ausgang (Closure):** entfallen mit gemessener Laufzeit und Mutationsprotokoll.
- **Größe über Schätzung wegen der Gegenprobe** (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Faktor 2 eingerechnet, Rückführung §4). — **Ausgang (Closure):** entfallen mit gemessener Summe, oder eingetreten → `slice-harness-devcontainer-einrichten-geruest` (Name vorläufig, Schnitt wie im Risiko darüber).
- **Konflikt mit dem Lauf an `harness/mk/gates-schnell.mk`** (Stand 2026-10-17 änderte ein Implementer die Datei für `slice-harness-coverage`, jetzt `done/`). — **Ausgang (Closure):** entfallen, wenn die Datei zum Start keinen offenen Lauf hat; sonst Abstimmung der Reihenfolge durch den Planner.
- **Die Sendung „Arbeit im Container“ wächst den Dokument-Punkt** (`BEO-REPO/slice-waechst-durch-uebernahmen`). — **Ausgang (Closure):** entfallen mit gemessener Summe nach der Metrik aus §4 (unter der Marke 900), oder eingetreten → `slice-harness-devcontainer-einrichten-geruest`.
- **Das Dokument veraltet still gegen das Ziel** (kein Sensor) und Aussagen über die Probe hinaus (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). — **Ausgang (Closure):** weiter offen im selben Eintrag, oder eingetreten (Evidence dieses Slice).
- **Umbau des Dokuments verliert Vorbehalte der Vorfassung** (Sicherheits- und Isolationsaussagen aus §11, „Behandle den Container nicht als Sandbox“). — **Ausgang (Closure):** entfallen mit Vergleich der Vorbehalte der Vorfassung mit dem Diff im Review; ein Befund ist *eingetreten* und geht in `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.
- **Docker-Daemon-Probe nur ohne laufendes `make gates`** (Abgrenzung §1). — **Ausgang (Closure):** entfallen, wenn die Vorbedingung `pgrep -fa "^make gates"` vor jeder Probe leer war und das im Bericht steht.
- **Arbeit des Implementers im Container-Klon bricht am Hostschlüssel oder am Egress** (§3, Randformen 9 und 10; erster Kontakt zu `pgr-dev2` ist ungemessen). — **Ausgang (Closure):** entfallen mit Beleg im Bericht, dass die Arbeit im Klon lief, oder eingetreten (Befund an den Planner, Eintrag in `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` oder `BEO-REPO/nested-runtime-podman-oder-docker-offen`).

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie ist so grob wie das Repo, trägt aber eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer für Werkzeug und Container: `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (Randform 9), `BEO-REPO/nested-runtime-podman-oder-docker-offen`, `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Zählerstände liest der Implementer aus den Dateien in `evidence/` vor dem Start (kein Zähler wird hier gesetzt). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel; die verkörperten (§3.10 bis §3.14) gelten.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Schnitt einträgt):** Liefer-Punkte 3 (Ziel mit Start-Ablauf, Gegenprobe, Dokument; die Sendung „Arbeit im Container“ liegt im Dokument-Punkt, die Sendung „Quellenauflösung“ liegt nicht mehr hier). Schichten 2 (Harness-Werkzeug mit Gegenprobe und Gate-Index; Nutzer- und Wartungs-Doku `docs/developer/`; das Gerüst unter `$PGR_DEV` liegt außerhalb des Repos und ist keine Schicht dieses Slice, `harness/README.md` ist Index desselben Werkzeugs, Planung und Vertrag sind keine Schicht). Schätzung in der Metrik von §4: 995 bis 1480 (vor dem Schnitt, mit Zugang und ohne Gerüst: 830 bis 1235; die Zahlen sind nicht vergleichbar, weil das Gerüst neu hineinkommt und der Zugang hinausgeht). Der Nehmer liegt nach Punkten und Schichten innerhalb der Größenregel und **mit dem ganzen Spektrum über der Marke 900**; vorgemerkter Schnitt: Gerüst und Start in `slice-harness-devcontainer-einrichten-geruest` (Name vorläufig, nicht angelegt, §4), vom Nutzer zu entscheiden, bevor der Slice nach `next/` geht. `grep` nach der Kennung dieses Slice findet ihn in `slice-harness-devcontainer-einrichten-zugang` (Folge-Slice, nimmt die Sendung an), in `slice-harness-devcontainer-u-boot-0-10` (Geber, `done/`), in den Plänen und der Roadmap der Release-Welle (keine Bindung) und in den Beobachtungen.
