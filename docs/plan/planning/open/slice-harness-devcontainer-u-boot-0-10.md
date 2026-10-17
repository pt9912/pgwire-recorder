# Slice slice-harness-devcontainer-u-boot-0-10: Dev-Container auf `u-boot` v0.10.0 umstellen (Features `sshd` und `claude-code`, Schlüssel als Umgebungsvariable, `make gates` im Container gemessen; `docs/developer/devcontainer.md`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein Abnahmeszenario und
kein Meilenstein hängt an ihm. Angelegt am 2026-10-17 auf Entscheidung des Nutzers vom selben Tag
(Reihenfolge: nach `slice-harness-coverage`, vor der Closure von `welle-v1-abschluss` und den Release-Slices).

**Bezug:** kein Anspruch des Produkts und keine Anforderung des Lastenhefts; die Host-Voraussetzungen des Repos (`git`, `docker`, `make`) bleiben unverändert (`ssh` und `u-boot` gelten nur für den Dev-Container). Vorgänger: `slice-harness-devcontainer-austausch-repo` und `slice-harness-devcontainer-doku` (beide `done/`). Keine aktive ADR; eine ADR entsteht nur, wenn der Nutzer die Laufzeit des Containers entscheidet (§1, §6).

**Berührte Spec-Stellen:** —

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der Dev-Container (`docs/developer/devcontainer.md`, Gerüst `pgr-dev` außerhalb des Repos) läuft auf `u-boot` v0.10.0 mit den Features `sshd`, `node` und `claude-code` statt handgebautem sshd, Node und Claude Code, der Zugangsschlüssel von Claude Code kommt als Umgebungsvariable, und es ist gemessen, ob und wie lange der volle `make gates` im Container läuft und ob `claude` dort startet und sich mit der Variablen anmeldet. Damit ist belegt, ob lang laufende Slices im Container laufen können und der Stop-Hook am Hauptbaum frei bleibt (`BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`).

Anlass und Stand von `u-boot` (Entscheidung des Nutzers vom 2026-10-17: „wir sollten zuerst dev-container angehen, dann können wir die lang laufenden slices darin laufen lassen und der main wird nicht mehr blockiert“): v0.10.0 liefert nach dem Auftrag die Features `sshd` (Port 2222, Konten gesperrt bis Schlüssel oder Passwort, Root-Anmeldung erlaubt) und `claude-code` (braucht das Feature `node`; im Sandbox-Profil `remoteEnv` mit `${localEnv:ANTHROPIC_API_KEY}`, Egress nur `registry.npmjs.org` und `api.anthropic.com`, Anmelde-Hosts trägt man selbst ein) sowie `nestedRuntime: podman` mit `podman build` (braucht `/dev/net/tun` am Host, gemessen nur mit Colima auf arm64; Docker Desktop ist nicht mehr zugesagt). Der Implementer liest den Wortlaut im Changelog und in der Dokumentation von `u-boot` v0.10.0; diese Angaben sind der Stand des Auftrags, nicht geprüft.

Der Schlüssel kommt als Umgebungsvariable (Entscheidung des Nutzers vom 2026-10-17, Randform 15 in `slice-harness-devcontainer-einrichten`, „da reicht vielleicht die Env-Var“): kein Mount der Datei. Der Container braucht den Schlüssel nicht besser zu schützen als der Host; das Token erneuert der Host, danach startet der Container neu (Entscheidung des Nutzers vom 2026-10-17). Dieser Slice klärt die offenen Punkte aus jener Randform: woher die Variable kommt, was `claude` bei fehlendem oder abgelaufenem Token tut.

Liefer-Punkte (Einzelheiten in §2): das Gerüst auf v0.10.0, die Messung und die Probe im Container, das Dokument.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Bare-Repo, keine Einlieferung, kein Klon, kein Make-Ziel `container-einrichten` und nicht „Arbeit im Container als Standardweg“: Das übernimmt `slice-harness-devcontainer-einrichten` (Klasse 1; der Nehmer nennt die Sendung dieses Slice in seinem §1 und zählt nach, §8 dort). Dieser Slice lässt den Ablauf aus §6 und §9 des Dokuments, wie er ist, und belegt nur, dass der Container ihn trägt (Punkt 2 der DoD).
- Keine Entscheidung zwischen `podman` und `docker` als `nestedRuntime` und keine ADR dazu: Die Messung gilt nur an einem Linux-Host, `u-boot` hat `podman` nur mit Colima auf arm64 gemessen, und die Wahl braucht den Nutzer. Der Slice misst mit der Laufzeit, die das Gerüst nach `u-boot` v0.10.0 für Docker-in-Docker braucht, und führt die Frage mit den Befunden der Probe als offene Entscheidung in §13 des Dokuments. Wählt der Nutzer `podman`, legt der Planner einen Folge-Slice an (§6).
- Keine Entscheidung zum Hostschlüssel (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`): Der Slice notiert nur, wann das Feature `sshd` die Hostschlüssel erzeugt (im Build oder beim Start), und setzt keine Einstellung (`StrictHostKeyChecking`, `UserKnownHostsFile`) ohne Entscheidung des Nutzers.
- Kein Umbau des laufenden Containers `pgr-dev` des Nutzers: Die Probe läuft in einer Scratch-Instanz (eigener Name, eigener Port, Entfernen nur per Name mit `-v`); der Umbau des echten Containers samt Volumes ist Handarbeit des Nutzers und steht im Dokument als Ablauf, nicht als gelaufen (§13).
- Kein Produkt-Code und kein Harness-Werkzeug (Schicht-Abgrenzung): Es ändern sich das Gerüst außerhalb des Repos und `docs/developer/`; nichts unter `cmd/`, `internal/`, `tools/`, `harness/`, kein Makefile-Fragment. Verlangt die Messung eine Änderung dort, ist das ein anderer Slice (Rückführung in §4).
- Keine Windows- und macOS-Hosts: Geprüft wird ein Linux-Host; die Aussage zu Docker Desktop und Colima übernimmt das Dokument aus dem Changelog von `u-boot` als Zitat, nicht als Probe.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Liefer-Punkte (drei):

- [ ] **Gerüst auf `u-boot` v0.10.0** (`$PGR_DEV`, außerhalb des Repos): `u-boot.yaml` mit den Features `sshd`, `node` und `claude-code`; `.devcontainer/Dockerfile` und `start.sh` verlieren, was die Features leisten (sshd, Node, Claude Code), und behalten, was sie nicht leisten (Docker-Daemon im Container und was dazugehört, nach Probe); der Schlüssel von Claude Code kommt als Umgebungsvariable des Hosts, ohne Mount (Randform 15). Beleg: Build und Start einer Scratch-Instanz, `ssh` mit dem Schlüssel aus `ssh/`, `claude --version`, die Variable im Container sichtbar (nur Name und Länge im Bericht, nie der Wert); die Datei-Inhalte des Gerüsts stehen im Dokument (Punkt 3) und in §7.
- [ ] **Messung und Probe im Container** (Scratch-Instanz, Klon wie in §6 des Dokuments von Hand angelegt): (a) volle `make gates` im Klon, mit Dauer, Ergebnis, Speicher und Docker-Speicher, bei vorher leerem `pgrep -fa "^make gates"` am Host; (b) `claude` startet im Container und meldet sich mit der Variablen an, einschließlich des Verhaltens bei fehlender und bei ungültiger Variable (Randformen 1 und 4 in §6); (c) der Wechsel der Laufzeit (`docker` oder `podman`) ist als Frage mit den Befunden der Probe formuliert, nicht entschieden. Beleg in §7 als Tabelle (Probe · Befehl · Ergebnis), die Ergebnisse stehen im Dokument.
- [ ] **Dokument** `docs/developer/devcontainer.md` auf `u-boot` v0.10.0 gebracht: Version und Voraussetzungen (§1, §2), Gerüst und Starten (§3 bis §5, Port und Konten des Features `sshd`), Claude Code und Schlüssel (§7), Arbeiten und Messung (§8), Grenzen (§11, die Vorbehalte der Vorfassung bleiben, geändert nur durch Befund), Nicht geprüft (§13: Entscheidung `podman` oder `docker`, Hostschlüssel, Umbau des laufenden Containers). Jede Aussage ist durch einen Lauf aus Punkt 1 und 2 belegt oder steht in §13 (`AGENTS.md` §3.11).

Gates und Closure:

- [ ] `make docs-check` und `make kopf-check` grün; `make gates` grün oder vom Stop-Hook freigegeben, weil nur Pfade der Positivliste ändern (`docs/developer/` steht auf der Positivliste des Hooks).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `$PGR_DEV/u-boot.yaml`, `$PGR_DEV/.devcontainer/Dockerfile`, `$PGR_DEV/.devcontainer/start.sh`, `$PGR_DEV/compose.yaml` (außerhalb des Repos) | update | Features statt handgebautem sshd, Node und Claude Code; Schlüssel als Umgebungsvariable (DoD Punkt 1); Dateien liegen nicht im Repo, ihr Stand steht im Dokument und in §7 |
| `docs/developer/devcontainer.md` | update | DoD Punkt 3; jede Aussage durch Lauf aus Punkt 1 oder 2 belegt oder in §13 |

- Probe zuerst am Scratch: Der Implementer legt das Gerüst für eine zweite Instanz an (Namen und Port nach §10 des Dokuments), baut und startet sie, misst darin und entfernt sie nur per Name mit Volumes (`docker rm -fv`, `docker volume rm` der benannten Volumes). Das laufende `pgr-dev` des Nutzers wird nicht angefasst (§1).
- **Docker-Last.** Die Messung von `make gates` im Container und ein `make gates` am Host laufen nicht gleichzeitig (Laufzeit und Speicher verfälscht, `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert`); `pgrep -fa "^make gates"` am Host ist vor der Messung leer, und das steht im Bericht. Der Klon stammt aus `main` nach dem Schließen von `slice-harness-coverage` (Start-Trigger, §4); die Messung gilt für diesen Stand, den der Bericht mit dem Commit nennt.
- Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-17): **3 Liefer-Punkte** (Gerüst, Messung und Probe, Dokument). **Schichten: 2**, nämlich das Dev-Container-Gerüst außerhalb des Repos (Betrieb des Werkzeugs, kein Produkt-Code und kein Harness-Werkzeug des Repos) und Nutzer- und Wartungs-Doku (`docs/developer/`); Planung und Register sind keine Schicht. Sendungen an andere Slices: eine, an `slice-harness-devcontainer-einrichten` (§1, erster Ausschluss); der Nehmer trägt sie im selben Commit in §1 und §6 ein und zählt in seinem §8 nach. Aus dem Nehmer kommt eine Sendung zurück: Der Risiko-Ausgang „`make gates` im Container ungemessen“ von `slice-harness-devcontainer-einrichten` wird durch DoD Punkt 2 aufgefangen (dessen §6, Eintrag mit der Kennung dieses Slice).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-coverage` liegt in `done/`, und `u-boot --version` am Host nennt `0.10.0` (Entscheidung des Nutzers vom 2026-10-17: erst Coverage, dann der Dev-Container, dann die Closure von `welle-v1-abschluss` und die Release-Slices). WIP-Limit 1 verbietet den Start neben `slice-harness-coverage`. Den Start der Release-Slices und die Closure von `welle-v1-abschluss` hält dieser Slice zusammen mit `slice-harness-devcontainer-einrichten` (Reihenfolge in der Roadmap, Drift-Log vom 2026-10-17).

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt. Das Gerüst liegt außerhalb des Repos und fehlt in dieser Zahl; sein Umfang steht im Bericht als Summe aus `git diff --no-index --numstat` zwischen altem und neuem Stand der vier Dateien. Schätzung (Planner, 2026-10-17): Dokument 200 bis 350 (drei Abschnitte mit Dockerfile, Skript und Compose-Block werden ersetzt, mehrere Abschnitte nachgezogen, eine Löschung zählt wie eine Einfügung), Gerüst 40 bis 100; Tests und Gegenprobe entfallen, der Faktor aus §3.14 greift nicht. Meldeschwelle 350, **Marke 450** (nur Dokument).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Summe nach der Metrik liegt über 450, oder das Gerüst verlangt mehr als die Features leisten (zum Beispiel eine eigene Laufzeit für Docker-in-Docker, die das Dokument über zwei Abschnitte neu schreibt), oder die Probe verlangt eine Änderung an `Makefile`, `tools/` oder `harness/` (anderer Vorgang, Folge-Slice mit Kennung). Schnitt dann: *Gerüst und Dokument* (Punkte 1 und 3) zuerst, die Messung (Punkt 2) als Folge-Slice `slice-harness-devcontainer-u-boot-0-10-messung` (Name vorläufig, vom Planner bei der Rückführung vergeben).
- `in-progress` → `open` (blockiert): Die Features von `u-boot` v0.10.0 laufen an diesem Host nicht ohne Eingriff am Host (Docker-Version, Gerätedateien, Netz), oder der volle Lauf von `make gates` im Container ist unlösbar (Speicher, Docker-in-Docker, Netz-Ausgang, Zeit); dann entscheidet der Nutzer, ob der Weg über den Container aufgegeben wird (der Nehmer `slice-harness-devcontainer-einrichten` geht dann mit Grund nach `done/`, und die Start-Trigger der Release-Slices verlieren beide Slices).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make docs-check` und `make kopf-check` grün auf dem Endstand (oder `make gates` grün), Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben. Die Tabelle der Proben in §7 trägt je Probe einen Lauf, der Scratch-Container und seine Volumes sind per Name entfernt und das ist im Bericht belegt. Das Ergebnis der Messung (Dauer, grün oder rot, Grenze) steht im Dokument und in §7, nicht nur im Bericht.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Kein neuer Vertrag eines Werkzeugs.** Der Slice liefert kein Make-Ziel, kein Gate und keine Zusage in `harness/vertraege/`; `AGENTS.md` §3.12 (Randformen vor dem Code in der Spezifikation) greift nicht. Die folgenden Randformen sind **Fragen der Probe**: Der Implementer beantwortet sie mit einem Lauf und schreibt die Antwort ins Dokument (oder nach §13); eine Antwort, die eine Entscheidung des Nutzers verlangt, gibt er dem Planner zurück und entscheidet sie nicht.

| # | Randform | Vorwissen (Auftrag, nicht geprüft) | Wo beantwortet |
|---|---|---|---|
| 1 | Variable `ANTHROPIC_API_KEY` am Host fehlt, ist leer oder nur Leerzeichen, wenn der Container startet (`${localEnv:…}`) | Das Profil reicht `${localEnv:ANTHROPIC_API_KEY}` durch; was `u-boot up` und `claude` bei fehlender Variable tun, ist offen | Dokument §7, Probe (b) |
| 2 | Welche Variable: `ANTHROPIC_API_KEY` (Feature) oder `CLAUDE_CODE_OAUTH_TOKEN` (Randform 15 in `slice-harness-devcontainer-einrichten` nennt beide) | Das Feature nennt `ANTHROPIC_API_KEY`; welches Konto der Nutzer hat, weiß nur er | Entscheidung des Nutzers (Rückgabe), Dokument §7 |
| 3 | Wo der Wert sichtbar ist (`docker inspect`, `/proc`, Kindprozesse, Logs, Shell-History) | Entscheidung des Nutzers vom 2026-10-17: kein Sicherheitsargument gegen das Durchreichen; das Dokument nennt die Orte, schützt sie nicht | Dokument §7 und §11, Probe (a) |
| 4 | Das Token ist abgelaufen oder ungültig | Beobachtung des Nutzers: die Anmeldemeldung kommt beim Start von `claude`, mitten in einer Sitzung ist sie ihm nicht aufgefallen; im Container nicht geprüft; das Token erneuert der Host, danach startet der Container neu | Dokument §7 und §13, Probe (b) |
| 5 | Egress: das Sandbox-Profil erlaubt nur `registry.npmjs.org` und `api.anthropic.com`; `make gates` braucht Hosts des Docker-Registry, des Go-Proxys und des Doku-Gate-Images | Ein Profil mit eingeschränktem Egress macht den Lauf unlösbar, ein Profil ohne ihn ändert die Aussage zur Sicherheit in §11; welches Profil, ist Entscheidung des Nutzers | Dokument §3 und §11, Probe (a) (Rückgabe bei Zielkonflikt) |
| 6 | Port: `sshd` des Features lauscht auf 2222; bisher bildet `compose.yaml` `127.0.0.1:2222:22` ab, und `PGR_SSH_PORT` aus dem Nehmer-Plan heißt der Port am Host | Das Mapping und die URL der Remote in §5, §9, §10 des Dokuments hängen daran; der Nehmer übernimmt den Port aus dem Ergebnis | Dokument §3 bis §5, §9, §10 |
| 7 | Konten gesperrt bis Schlüssel oder Passwort; Root-Anmeldung erlaubt | Bisher: `PermitRootLogin no`, einziger Schlüssel für `vscode` (§5, §11 des Dokuments); die Aussage ändert sich mit dem Feature | Dokument §5 und §11, Probe (a) |
| 8 | Hostschlüssel: erzeugt das Feature sie im Build oder beim Start; ändert ein Neubau den Fingerabdruck | `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`; keine Einstellung ohne Entscheidung des Nutzers | Dokument §5 und §13 (nur Befund) |
| 9 | Docker-Daemon im Container: leistet das Gerüst nach v0.10.0 Docker-in-Docker (`privileged`, `dockerd`), oder nur `nestedRuntime` | Bisher: eigener `dockerd` in `start.sh`; `nestedRuntime: podman` mit `podman build`, nicht mit `docker`; ob `make` (`docker build`, `docker run` im `Makefile`) mit `podman` läuft, ist offen | Dokument §3, §8 und §13 (Frage, nicht Entscheidung) |
| 10 | Host ohne `/dev/net/tun` oder anderer Host als Colima auf arm64 | Gemessen hat `u-boot` nur Colima auf arm64; dieser Host ist Linux; Docker Desktop ist nicht mehr zugesagt | Dokument §1 und §13 |
| 11 | Wechsel eines laufenden Containers (Volumes `pgr-dev-docker` und `pgr-dev-workspace` mit Klon und Daten) auf das neue Gerüst | Der Umbau des echten Containers ist Handarbeit des Nutzers; die Probe nutzt eine Scratch-Instanz | Dokument §4 und §13 (nicht gelaufen) |

**Weitere Risiken:**

- **Der volle Lauf von `make gates` im Container ist ungemessen** (`docs/developer/devcontainer.md` §13; Risiko aus `slice-harness-devcontainer-einrichten`, das hier aufgefangen wird). Läuft er nicht oder dauert er zu lang, bleibt der Container ohne Gate-Nachweis, und der Stop-Hook am Hauptbaum bleibt unter Last des Implementers gesperrt. — **Ausgang (Closure):** entfallen mit der gemessenen Dauer und grün, oder eingetreten → Rückführung `in-progress` → `open` (§4) mit Entscheidung des Nutzers und neuem Slice mit Kennung.
- **Die Laufzeit-Frage (`podman` oder `docker`) bleibt offen.** — **Ausgang (Closure):** weiter offen → `BEO-REPO/nested-runtime-podman-oder-docker-offen` (neu, 1×, Evidence dieser Slice), bis der Nutzer entscheidet; entscheidet er für `podman`, legt der Planner einen Folge-Slice an.
- **Hostschlüssel nach Neubau** (Randform 8). — **Ausgang (Closure):** weiter offen → `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (Evidence dieser Slice, wenn der Befund das Wiederauftreten zeigt).
- **Die Messung läuft neben einem anderen Lauf am Host** (Docker-Daemon und CPU geteilt; `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert`). — **Ausgang (Closure):** entfallen, wenn `pgrep -fa "^make gates"` vor der Messung leer war und das im Bericht steht; sonst eingetreten → Messung wiederholt.
- **Das Dokument veraltet still gegen das Gerüst** (liegt außerhalb des Repos, kein Sensor) und Aussagen über die Probe hinaus (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). — **Ausgang (Closure):** weiter offen im selben Eintrag, oder eingetreten (Evidence dieses Slice).
- **Umbau des Dokuments verliert Vorbehalte der Vorfassung** (§11, „Behandle den Container nicht als Sandbox“, Aussagen zur Anmeldung und zum Netz). — **Ausgang (Closure):** entfallen mit Vergleich der Vorbehalte der Vorfassung mit dem Diff im Review; ein Befund ist *eingetreten* und geht in `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.
- **Größe** über Schätzung, weil das Dokument drei Blöcke ersetzt (§4). — **Ausgang (Closure):** entfallen mit gemessener Summe, oder eingetreten → Folge-Slice `slice-harness-devcontainer-u-boot-0-10-messung`.
- **Der Start der Release-Slices hängt an diesem Slice und an `slice-harness-devcontainer-einrichten`.** Scheitert die Messung, blockiert der Weg zum ersten Release ohne Produktgrund. — **Ausgang (Closure):** entfallen, sobald beide Slices in `done/` liegen; sonst entscheidet der Nutzer über die Start-Trigger (Rückführung `in-progress` → `open`, §4).
- **Implementer-Läufe im Hauptbaum sperren den Stop-Hook** (`BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`, Anlass dieses Slice). — **Ausgang (Closure):** weiter offen im selben Eintrag; die Beobachtung wird erst mit dem Standardweg im Container (`slice-harness-devcontainer-einrichten`) beantwortet, nicht mit diesem Slice.

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; das Gerüst außerhalb des Repos hat keine eigene Deklaration und gehört zu `REPO` als Betrieb des Werkzeugs, keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-17). Treffer für Container, Stop-Hook und Doku: `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (1×, offen; Randform 8), `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook` (0×, neu angelegt am 2026-10-17, benannt, nicht gezählt), `BEO-REPO/oom-kill-beim-gate-lauf-ursache-ungeklaert` (1×, offen; Risiko der Messung), `BEO-REPO/stop-hook-positivliste-ohne-docs-developer` (2×, offen; `docs/developer/` steht auf der Positivliste, das Risiko eines Gates, das den Pfad liest, gilt unverändert), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (75×, verkörpert in `AGENTS.md` §3.11), `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (7×, verkörpert; hier gemessen, nicht verlängert). Keine offene Klasse erreicht mit diesem Plan 3× neu ohne Regel.

alle berührten Sub-Areas GF (`REPO`: das Dokument beschreibt, was die Probe zeigt; die Proben gehen dem Text voraus).

**Nachzählung der Sendung an `slice-harness-devcontainer-einrichten` (Planner, 2026-10-17):** Der Geber trägt die Sendung nicht ein; sie steht im Nehmer (§1, §6, §8 dort). Dieser Slice zählt 3 Liefer-Punkte und 2 Schichten (§3) und nimmt aus dem Nehmer den Risiko-Ausgang „`make gates` im Container ungemessen“ auf; er hat damit keine Sendung zu tragen, die die Zählung hebt.
