# Slice slice-harness-devcontainer-einrichten: Make-Ziel `container-einrichten` richtet im laufenden Dev-Container Bare-Repo, Einlieferung, Klon und Zweig ein (`docs/developer/devcontainer.md`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-dev-container (Entscheidung des Nutzers vom 2026-10-18; bis dahin ohne Welle). Angelegt am 2026-10-17 auf Wunsch des Nutzers vom selben Tag; am 2026-10-18 zweimal geschnitten: Der Zugang von Claude Code liegt in `slice-harness-devcontainer-einrichten-zugang` (Entscheidung des Nutzers vom 2026-10-18: nicht gebraucht, kein Login, die Hostsitzung treibt die Arbeit per `ssh` im Container), Gerüst und Start in `slice-harness-devcontainer-einrichten-geruest` (die Schätzung 995 bis 1480 lag mit dem ganzen Spektrum über der Marke 900); der Geber wurde am 2026-10-18 ein zweites Mal geschnitten, der Start (Container, Port, Bereitschaft, Abbruch) liegt in `slice-harness-devcontainer-einrichten-start`.

**Bezug:** [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität) nur insoweit, als die Host-Voraussetzungen des Repos berührt sind (neu: `ssh`, Entscheidung des Nutzers, §1); kein Anspruch des Produkts.
Vorgänger: `slice-harness-devcontainer-austausch-repo` und `slice-harness-devcontainer-doku` (beide `done/`); der Slice automatisiert deren Ablauf. `slice-harness-devcontainer-u-boot-0-10` (`done/`; Gerüst auf `u-boot` v0.10.0, `make gates` im Container grün, 1701 s). Die unmittelbaren Vorgänger sind `slice-harness-devcontainer-einrichten-geruest` (legt das Gerüst an) und `slice-harness-devcontainer-einrichten-start` (startet den Container); dieser Slice setzt den laufenden Container voraus, §1 und §4. Folge-Slice: `slice-harness-devcontainer-einrichten-zugang` (Zugang von Claude Code; kein Teil dieses Slice, §1). Keine aktive ADR; der Architect legt eine an (§6).

**Berührte Spec-Stellen:** neuer Vertrag eines Harness-Werkzeugs in `harness/vertraege/` (Nummer vergibt der Architect); gebunden an [`SPEC-058`](../../../../harness/vertraege/SPEC-058-container-start.md) (Punkte 2 und 11, Zusagen der Geber, §1); Host-Voraussetzungen: [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-17, geschnitten 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Make-Ziel `container-einrichten` am Host richtet für einen laufenden Dev-Container (`docs/developer/devcontainer.md` §6 und §9) in einem Aufruf ein: Bare-Repo `/workspaces/<Projektname aus u-boot.yaml>/austausch.git`, Einlieferung von `main`, Arbeitsklon `/workspaces/<Projektname>/pgwire-recorder`, Identität des Klons aus der Git-Konfiguration des Host-Repos, `make hooks-install` und `make traeger-fetch` im Klon, Zweig für den Slice. Arbeit im Dev-Container wird damit der Standardweg, und `main` am Host bleibt für den Stop-Hook frei.

**Arbeitsweise (Entscheidung des Nutzers vom 2026-10-18):** Die Hostsitzung treibt die Arbeit im Container per `ssh` und Befehl; Claude Code läuft nicht im Container, und das Ziel löst keinen Zugang auf (Abgrenzung unten).

**Bindung des Gebers `slice-harness-devcontainer-einrichten-geruest` (eingetragen vom Planner am 2026-10-18, `AGENTS.md` §3.13; am selben Tag geteilt):** Das Ziel setzt voraus, was die beiden Geber zusagen. `slice-harness-devcontainer-einrichten-geruest`: `$PGR_DEV` mit `u-boot.yaml` (der Projektname bestimmt den Pfad `/workspaces/<Projektname>/`) und die Variablen `PGR_DEV`, `PGR_SSH_PORT`, `PGR_NAME` mit den Vorgaben aus Dokument §10 (`SPEC-058` Punkt 2). `slice-harness-devcontainer-einrichten-start`: einen Container mit dem Namen `PGR_NAME`, der per `ssh` als `vscode` mit `$PGR_DEV/ssh/id_ed25519` auf `127.0.0.1:${PGR_SSH_PORT:-2222}` antwortet (`SPEC-058` Punkt 11); dass `dockerd` im Container bereit ist, sagt `slice-harness-devcontainer-einrichten-start` nicht zu (dessen Punkt 17 (4)), das Ziel prüft es selbst. Das Ziel legt kein Gerüst an, startet keinen Container und ändert nichts am Gerüst; ob `container-einrichten` das Ziel des Gebers aufruft oder nur auf einen laufenden Container prüft, entscheidet der Architect (E1). Ändert einer der Geber einen der Punkte, folgt dieser Slice in §1, §3 und §6.

**Sendung von `slice-harness-devcontainer-u-boot-0-10` („Arbeit im Container“, eingetragen vom Planner am 2026-10-17, Entscheidung des Nutzers vom selben Tag), Teil dieses Slice:** Lang laufende Slices laufen im Container, damit `main` am Host nicht über Stunden durch einen Implementer-Lauf gesperrt ist (`BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`). Dieser Slice trägt davon den Einrichtungsweg im Container (das Ziel) und im Dokument den Ablauf eines Implementer-Laufs im Container-Klon: Arbeit im Klon, `make gates` dort, Ergebnis per `git fetch container` am Host. Den Weg bis zum laufenden Container trägt der Geber `slice-harness-devcontainer-einrichten-geruest`. Dieser Slice baut auf den Befunden von `slice-harness-devcontainer-u-boot-0-10` auf und misst sie nicht neu.

**Annahme (Entscheidung des Nutzers vom 2026-10-17, „ja“): ein Slice pro Container-Sitzung.** Ein Slice läuft in einer Sitzung, danach Neustart (Stopp und Neustart trägt der Geber). Je Slice ein Zweig aus `main`, das Ergebnis kommt per `git fetch container` an den Host. Das Ziel ist deshalb der **Einrichtungs-Ablauf einer Sitzung** und wiederholbar: es läuft bei jedem Neustart, nicht einmal je Container (Randformen 1 bis 3, 16). Ein Slice, der länger als rund 8 Stunden läuft, braucht einen Neustart mitten in der Arbeit; das ist kein Fehler des Werkzeugs (§6).

Ablauf des Ziels, in dieser Reihenfolge (die Schritte 4 bis 10 des Slice vor dem Schnitt):

1. Vorbedingungen am Host prüfen, bevor etwas geändert wird: `ssh` und `git` vorhanden, `main` im Host-Repo (Randform 12), Identität des Host-Repos (Randform 5), Container per `ssh` erreichbar (Randform 4), Remote `container` zeigt auf die erwartete URL oder fehlt (Randform 3).
2. Bare-Repo anlegen, falls es fehlt: `git init --bare --initial-branch=main /workspaces/<Projektname>/austausch.git` im Container, per `ssh` (Randform 1).
3. `main` einliefern: Remote `container` (`ssh://vscode@127.0.0.1:${PGR_SSH_PORT:-2222}/workspaces/<Projektname>/austausch.git`) und `core.sshCommand` im Host-Repo, dann `git push container main` (Randform 13).
4. Klon `/workspaces/<Projektname>/pgwire-recorder` anlegen, falls er fehlt, sonst aktualisieren (Randform 2).
5. Name und E-Mail des Klons **lokal** aus der Git-Konfiguration des Host-Repos setzen: nie `--global`, nie ein im Dokument oder Skript fest stehender Name; fehlt am Host einer der beiden Werte, bricht das Ziel ab, bevor es etwas ändert (Schritt 1; Reihenfolge der Gründe: Randform 14); Werte mit Anführungszeichen oder Apostroph kommen unverändert an (Randformen 5 bis 7).
6. `make hooks-install` und `make traeger-fetch` im Klon (Randform 10).
7. Zweig für den Slice aus `main` im Klon anlegen (Randform 16).

**Bindung an `slice-harness-devcontainer-einrichten-zugang` (Sendung, eingetragen vom Planner am 2026-10-18):** Das Skript liest, löst und übergibt keinen Zugang (weder `ANTHROPIC_API_KEY` noch `CLAUDE_CODE_OAUTH_TOKEN` noch eine Anmeldedatei) und prüft keinen; sein Ablauf lässt eine Stelle vor Schritt 1 frei, an der der Nehmer die Auflösung einfügt (Randform 25, Architect). Die Zeile mit dem Verweis `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}` in `mk-override.sh` schreibt der Geber `slice-harness-devcontainer-einrichten-geruest`, nicht dieser Slice.

**Sendung „Quellenauflösung“ von `slice-harness-devcontainer-u-boot-0-10` (eingetragen vom Planner am 2026-10-17) geht am 2026-10-18 weiter an `slice-harness-devcontainer-einrichten-zugang`** (Entscheidung des Nutzers vom 2026-10-18). Der Geber nennt als Nehmer weiter diesen Slice; die Adresse bleibt einen Hop länger: Dieser Slice nennt den Nehmer, im Geber (`done/`) wird kein Zeiger umgeschrieben. Randform 15, Quellenauflösung, Übergabe über `devcontainer exec`, `remoteEnv` und beide Variablen stehen dort; dieser Slice führt sie nicht.

**Entscheidung des Nutzers (2026-10-17): `ssh` ist neue Voraussetzung am Host.** Auf die Frage, ob `ssh` als Voraussetzung für dieses Ziel in Ordnung ist, antwortete der Nutzer „ja“. Die Voraussetzungen des Repos sind bisher `git`, `docker` und `make` (`AGENTS.md` §3.1, [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität)). Der Architect trägt die Entscheidung in die Spezifikation des Vertrags und in die ADR ein (§3); `ssh` gilt nur für dieses Ziel, kein Gate und kein anderes Ziel braucht es. Die übrigen Voraussetzungen (`u-boot`, `devcontainer`-CLI, `node`, `ssh-keygen`) gehören dem Geber `slice-harness-devcontainer-einrichten-geruest`; der Nutzer hat sie mit [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md) am 2026-10-18 entschieden.

**Kein Gate**, wie `hooks-install`: Das Ziel steht im Abschnitt *Werkzeuge* von `harness/README.md`, nicht in `GATE_CHECKS`. Seine Gegenprobe ist ein Gate (§2).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Gerüst, kein Schlüsselpaar, kein `mk-override.sh`, kein Start, Stopp oder Neustart des Containers und keine Voraussetzungsprüfung für `u-boot`, `devcontainer`, `node`: Das tragen `slice-harness-devcontainer-einrichten-geruest` (Gerüst, Schlüsselpaar, `mk-override.sh`, Voraussetzungen) und `slice-harness-devcontainer-einrichten-start` (Start, Stopp, Neustart) (Klasse 1, Geber; deren §1 nennen diese Schritte und die Bindung an diesen Slice). Das Ziel setzt einen laufenden Container voraus (Dokument §3 bis §5) und fragt nur, ob er per `ssh` erreichbar ist.
- Kein Zugang von Claude Code im Container: keine Quellenauflösung, keine Randform 15, kein `remoteEnv`-Weg über `devcontainer exec`, keine Variable `ANTHROPIC_API_KEY` oder `CLAUDE_CODE_OAUTH_TOKEN` im Skript, kein Dokument §7 (Klasse 1: `slice-harness-devcontainer-einrichten-zugang`, dessen §1 diese Sendung annimmt; Entscheidung des Nutzers vom 2026-10-18). Das Ziel richtet den Container für Arbeit per `ssh` ein, nicht für `claude`.
- Keine Entscheidung zum Hostschlüssel (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`): Das Ziel nennt, was `ssh` ausgibt, und setzt keine Einstellung (`StrictHostKeyChecking`, `UserKnownHostsFile`) ohne Entscheidung des Nutzers (§6, Randform 9).
- Keine Messung von `make gates` im Container und keine Entscheidung zu `podman` oder `docker` als `nestedRuntime`: Das liefert `slice-harness-devcontainer-u-boot-0-10` (`done/`) bzw. bleibt offen (`BEO-REPO/nested-runtime-podman-oder-docker-offen`).
- Keine Docker-Daemon-Probe im Ziel selbst: Es prüft den Container per `ssh`, nicht per `docker`. Die Gegenprobe darf den Daemon nur nutzen, wenn kein `make gates` läuft (§3, *Docker-Last*).
- Kein Produktcode (Schicht-Abgrenzung): Es ändern sich Werkzeug, Gegenprobe und Dokumentation, nichts unter `cmd/` und `internal/`.
- Kein Eintrag in der Roadmap: Der Slice ist wellenlos. Die Release-Slices und die Closure von `welle-v1-abschluss` hängen nicht an ihm (Entscheidung des Nutzers vom 2026-10-17).

## 2. Definition of Done

Liefer-Punkte (drei):

- [ ] **Ziel `container-einrichten`** (Fragment unter `harness/mk/` und Skript unter `tools/harness/`): erfüllt den Vertrag, den der Architect in `harness/vertraege/` festlegt (Ablauf aus §1, Randformen aus §6); kein Gate; in `harness/README.md` §Werkzeuge genannt.
- [ ] **Gegenprobe** `container-einrichten-gegenprobe` (Gate, hängt an `GATE_CHECKS`, steht in der Ausnahmeliste der schnellen Kette): gegen einen Scratch-Container (Präfix, Entfernen nur per Name, `docker rm -fv`, nichts Fremdes, nie der Container, in dem der Implementer arbeitet) je Randform aus §6 ein Fall und je Zusage und je Element einer Aufzählung eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die Gegenprobe steht im Gate-Index.
- [ ] **Dokument** `docs/developer/devcontainer.md` §6 und §9 nennen `make container-einrichten` als Weg und führen die Einzelbefehle nur, soweit das Dokument sie als Erklärung braucht; dazu der Ablauf eines Implementer-Laufs im Container-Klon (Sendung von `slice-harness-devcontainer-u-boot-0-10`, §1): Arbeit im Klon per `ssh`, `make gates` dort mit der Dauer, die jener Slice misst, Ergebnis per `git fetch container`. §3 bis §5 gehören `slice-harness-devcontainer-einrichten-geruest`, §7 `slice-harness-devcontainer-einrichten-zugang`. Jede Aussage ist durch einen Lauf belegt oder steht in §13.

Gates und Closure:

- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `tools/`, `harness/` und das Makefile-Fragment ändern); die Regel, nicht gleichzeitig mit dem Lauf eines anderen Slice auf demselben Daemon zu laufen, richtet sich an den Aufrufer (§3).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-<NNN>-container-einrichten.md` | neu (Architect) | Vertrag des Werkzeugs: Ablauf, Randformen aus §6, Host-Voraussetzung `ssh`; Nummer vergibt der Architect |
| `docs/plan/adr/` (eine ADR) | neu (Architect) | Entscheidung: Ort des Ziels (folgt dem Geber), Host-Voraussetzung `ssh` (Entscheidung des Nutzers vom 2026-10-17); `Schärft:` zeigt auf den Vertrag (`AGENTS.md` §3.8) |
| `harness/mk/container-einrichten.mk` | neu | Ziel `container-einrichten` (kein Gate), Gegenprobe hängt an `GATE_CHECKS` |
| `tools/harness/container-einrichten.sh` | neu | Ablauf aus §1 (Prüfen, Einrichtung, Zweig); Eingaben aus der Umgebung (`PGR_DEV`, `PGR_SSH_PORT`, Remote-Name), Fehlermeldungen mit Schritt; keine Zugangsvariable (Bindung, §1) |
| `tools/harness/container-einrichten-gegenprobe.sh`, `harness/mk/container-einrichten-gegenprobe.mk` | neu | Gegenprobe gegen Scratch-Container; Fälle und Mutationen aus §6 |
| `harness/mk/gates-schnell.mk` | update | Gegenprobe in die Ausnahmeliste der schnellen Kette (Docker; Pflege nach `SPEC-052`, Architect prüft Punkt der Gate-Kette) |
| `harness/README.md` | update | Abschnitt *Werkzeuge* (`container-einrichten`) und *Sensors* (Gegenprobe) |
| `docs/developer/devcontainer.md` | update | §6 und §9 nennen das Ziel; Ablauf eines Implementer-Laufs; §13 für Ungeprüftes |

- Das Skript trennt Prüfen und Handeln: erst alle Vorbedingungen am Host (Identität, `ssh` vorhanden, Schlüssel lesbar, Container per `ssh` erreichbar), dann die Schritte im Container. Die Reihenfolge der Gründe legt der Architect fest (Randform 14).
- Die Identität reist ohne Auswertung durch eine zweite Shell, soweit der Architect keine andere Form festlegt: Der Implementer entscheidet die Quoting-Form nicht, sondern setzt die des Vertrags um (§6, Sonderzeichen).
- **Docker-Last.** Die Gegenprobe braucht einen laufenden Container und startet dafür einen Scratch-Container (Präfix, Entfernen nur per Name, `docker rm -fv` und Prüfung, dass keines der Volumes bleibt, wie `SPEC-057`). Eine Vorbedingung „kein `make gates` läuft“ gibt es nicht: Die Gegenprobe läuft selbst in `make gates`, und gleichzeitige Läufe auf einem Daemon sind nicht zugesagt (`SPEC-052` Punkt 8, `SPEC-058` Punkt 15); die Regel „nie gleichzeitig mit dem Lauf eines anderen Slice“ richtet sich an den Aufrufer. Ihre Laufzeit ist zu messen (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`). Wie die Gegenprobe den Scratch-Container erzeugt (über das Ziel des Gebers oder einen Rahmen), legt der Architect fest (E1).
- **Arbeit des Implementers im Dev-Container `pgr-dev2`** (Entscheidung des Nutzers vom 2026-10-18). Stand 2026-10-18: Der Container `pgr-dev2` läuft (Host-Verzeichnis `/Development/pgr-dev2`, Port `127.0.0.1:2223`), der Klon `/workspaces/pgr-dev2/pgwire-recorder` besteht, das Host-Repo hat das Remote `container` (`ssh://vscode@127.0.0.1:2223/workspaces/pgr-dev2/austausch.git`; Remote heißt `container`, nicht `container2`), und die Hostschlüssel-Datei `/Development/pgr-dev2/known_hosts` liegt vor; kein Claude Code im Container. Die Hostsitzung treibt den Klon per `ssh`:
  1. *Befehl.* `ssh -i /Development/pgr-dev2/ssh/id_ed25519 -o IdentitiesOnly=yes -o UserKnownHostsFile=/Development/pgr-dev2/known_hosts -p 2223 vscode@127.0.0.1 'cd /workspaces/pgr-dev2/pgwire-recorder && <Befehl>'`. Die Datei hat der Nutzer bereitgestellt; die Einstellung gilt für die eigenen Befehle des Implementers, das Ziel dieses Slice setzt sie nicht (Abgrenzung §1, Randform 9).
  2. *Vorbefund.* Bare-Repo, Klon und Remote bestehen für `pgr-dev2` schon (Handarbeit nach Dokument §6 und §9). Der Implementer prüft vor dem Start, dass Bare-Repo und Klon antworten (`git -C … status`, `git ls-remote container`), Identität des Klons lokal gesetzt, `core.hooksPath` gesetzt, Träger `ai-harness-init` vorhanden; was fehlt, trägt er nach dem Dokument nach und nennt es im Bericht. Nach der Lieferung des Ziels ruft er es für diese Instanz auf (Probe am eigenen Arbeitsplatz: zweiter Aufruf auf bestehendem Bare-Repo und Klon, Randformen 1 bis 3).
  3. *Arbeiten.* Zweig `slice-harness-devcontainer-einrichten` im Klon, nicht `main`. Dateien kommen als Patch über `stdin` (`git apply`) oder als Inhalt über `stdin` (`cat > Datei`), nicht über eine eingehängte Arbeitskopie (der Container sieht den Arbeitsbaum des Hosts nicht). Commits im Klon, jeder mit der Kennung `slice-harness-devcontainer-einrichten` im Text (Hook).
  4. *Gates.* `make gates-schnell` und `make test` im Klon; der volle Lauf `make gates` dauert 28 min 21 s (Dokument §8, für diesen Stand ungeprüft) und läuft im Hintergrund (`setsid`, Log in `/tmp` im Container, Exit-Code in eine Datei), Abfrage per `ssh`; vor einem eigenen Lauf der Gegenprobe außerhalb von `make gates` prüft der Aufrufer selbst, dass auf demselben Daemon kein anderer Lauf läuft (Regel an den Aufrufer, Docker-Last).
  5. *Was nur der Host kann.* Die Gegenprobe braucht `ssh` und `git` und einen Scratch-Container; ob der Daemon im Container von `pgr-dev2` Container mit `sshd` starten kann (Image, Egress), ist ungemessen. Der Implementer prüft das zuerst im Klon und nennt das Ergebnis im Bericht. Trägt der Daemon im Container es nicht, läuft der Scratch-Teil am Host mit Namen, Port und Volumes im Präfix der Gegenprobe (nie `pgr-dev2` anfassen, nie stoppen, neu starten, entfernen, nie Port 2223 belegen); der Bericht sagt, was im Container und was am Host lief.
  6. *Ergebnis zum Host.* Im Klon `git push origin slice-harness-devcontainer-einrichten`; am Host `git fetch container`, der Zweig liegt als `container/slice-harness-devcontainer-einrichten` zur Prüfung vor (Dokument §9). Der Inhalt stammt aus einem Container, in dem `vscode` `root` werden kann: Review und Verifikation sehen die Commits durch, bevor sie nach `main` gehen. Der Übergang `next/` → `in-progress/` ist ein Commit des Nutzers auf `main` vor der Arbeit und wird vor dem ersten Stand im Klon eingeliefert (`git push container main`).
  7. *Größenprobe.* Die Metrik von §4 misst der Implementer im Klon mit `git diff --numstat <Basis>..<letzter Commit>`.
- Zählung nach `AGENTS.md` §3.13 siehe §8.

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten-start` liegen in `done/` (Reihenfolge: beide Geber, dann dieser Slice, dann gegebenenfalls `slice-harness-devcontainer-einrichten-zugang`); `slice-harness-devcontainer-u-boot-0-10` liegt in `done/` (erfüllt, Messung grün: zwei Läufe zu je 1701 s). Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. Die Release-Slices und die Closure von `welle-v1-abschluss` verlangen nur `slice-harness-devcontainer-u-boot-0-10` in `done/`, diesen Slice nicht (Entscheidung des Nutzers vom 2026-10-17). Vorbedingung des Starts außerdem: Vertrag und ADR des Architect liegen vor (`Accepted` ist die Entscheidung des Nutzers, keine Annahme des Planners).

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; Vertrag und ADR (Commits des Architect vor dem Start) liegen vor der Basis und zählen nicht.

**Schätzung nach dem Schnitt vom 2026-10-18 (Planner; Metrik wie oben, Gegenprobe mit dem Doppelten des Produktcodes, Faktor aus `AGENTS.md` §3.14):**

| Teil | Zeilen |
|---|---|
| Skript: Prüfen (Identität, `ssh`, `git`, `main`, Erreichbarkeit, Remote) | 20 bis 30 |
| Skript: Einrichtung (Bare-Repo, Einlieferung, Klon, Identität, Hook, Träger, Zweig) | 115 bis 170 |
| Fragmente | 30 bis 40 |
| Gegenprobe (Doppeltes von Skript und Fragmenten, 165 bis 240) | 330 bis 480 |
| Dokument (§6, §9 gekürzt, Löschungen zählen, dazu Ablauf des Implementer-Laufs, Index) | 70 bis 110 |
| **Zusammen** | **565 bis 830** |

Meldeschwelle 700, **Marke 900**. Die obere Hälfte liegt über der Meldeschwelle, die ganze Schätzung unter der Marke; ein Schnitt unterbleibt. Die Summe von Geber und diesem Slice (1125 bis 1695) liegt über der des ungeteilten Slice (995 bis 1480), weil Fragmente, Gegenprobe-Rahmen und Index doppelt anfallen.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder der Vertrag verlangt mehr als drei Liefer-Punkte oder mehr als zwei Schichten (zum Beispiel eine Änderung am Gerüst des Gebers), oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner (Vorschlag: Ziel mit Dokument zuerst, die Gegenprobe der Randformen als Folge-Slice mit eigener Kennung).
- `in-progress` → `open` (blockiert): Der Ablauf trägt nicht ohne Eingriff in das Image (zum Beispiel fehlt `git` im Container, Rechte an `/workspaces`), oder `ssh` am Host ist nicht zugesagt (Entscheidung des Nutzers wird zurückgenommen), oder das Gerüst, das der Geber hinterlässt, trägt den Ablauf nicht (der Befund geht an den Planner).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün auf dem Endstand, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben. Das Mutationsprotokoll (Zusage · Mutation · roter Test, `AGENTS.md` §3.10) steht im Bericht des Implementers. Die Gegenprobe ist einmal gegen einen echten Scratch-Container grün gesehen, nicht nur gegen Attrappen.

## 6. Risiken und offene Punkte

Der Slice liefert einen **neuen Vertrag** (Werkzeug mit Zusagen); §3.12 gilt. Die Randformen stehen hier, entschieden werden sie **vor dem ersten Code-Commit** in `harness/vertraege/` (Architect); eine Form, die erst ein Commit mit Code oder Tests hier einträgt, ist nicht genannt. Der Implementer entscheidet keine still, er hält an und gibt sie dem Architect zurück. Vorläufiges Wissen aus dem Dokument steht dabei, nicht als Entscheidung.

**Entscheidungen des Architect vor dem Code** (die Liste E1 bis E7 des ungeteilten Slice ist auf diesen Slice und den Geber verteilt):

- **E1 Ort des Ziels.** Folgt der Entscheidung des Gebers (`slice-harness-devcontainer-einrichten-geruest`, dessen E1). Dieses Ziel ändert das Host-Repo (`remote`, `core.sshCommand`) und gehört deshalb ins Repo; ob es das Ziel des Gebers aufruft oder nur einen laufenden Container prüft, und wie die Gegenprobe ihren Scratch-Container erzeugt, legt der Architect fest.
- **E2 Host-Voraussetzungen.** Dieser Slice braucht `ssh` (entschieden, §1) und `git`; `u-boot`, `devcontainer`, `node`, `ssh-keygen` gehören dem Geber (dessen E2). Der Nutzer entscheidet dort; **Empfehlung des Planners:** alle aufnehmen.
- **E3 Randformen des Aufrufs** (Zeilen unten).
- **E4 Aufrufarten von `make`** (unten).
- **E5 Reihenfolge der Gründe** (Randform 14) und Exit-Codes.
- **E7 Schätzung** (§4): Der Architect bestätigt sie in der Metrik vor dem Code.
- E6 (Stopp/Neustart) liegt beim Geber.

**Randformen des Ablaufs** (Spalte *Wo entschieden*: Fundstelle im Vertrag, vom Architect einzutragen). Die Zeilen 17 bis 24 des ungeteilten Slice (Werkzeug, `PGR_DEV`, Verzeichnis, Port, zweiter Aufruf, Schlüsselpaar, `u-boot`, Warten auf `ssh`) liegen als A1 bis A12 beim Geber:

| # | Randform | Vorwissen aus `docs/developer/devcontainer.md` (§6, §9, §10) | Wo entschieden |
|---|---|---|---|
| 1 | Bare-Repo besteht schon | `git init --bare` auf vorhandenem Verzeichnis meldet `Reinitialized`, Exit 0; Frage: Wiederholung ist erlaubt oder bricht ab? Besteht es mit anderem `HEAD` (nicht `main`)? | offen |
| 2 | Klon besteht schon (leer oder nicht leer, anderer Remote) | `git clone` in nicht leeres Verzeichnis: Exit 128. Frage: überspringen, aktualisieren oder abbrechen; was geschieht mit uncommittetem Stand im Klon? (`pgr-dev2` hat einen Klon, Stand 2026-10-18) | offen |
| 3 | Remote `container` besteht schon (gleiche URL, andere URL, anderer Port) | zweites `git remote add` bricht mit Exit 3 ab; das Host-Repo hat `container` auf Port 2223 (`pgr-dev2`). Frage: bei gleicher URL weiter, bei anderer abbrechen oder setzen | offen |
| 4 | Container nicht erreichbar beim Aufruf (gestoppt, Port, Schlüssel, Zeitgrenze) | Der Start gehört dem Geber (A8). Frage: Meldung, Exit-Code, Zeitgrenze für die Verbindung (`ConnectTimeout`), vor jeder Änderung, Verweis auf das Ziel des Gebers | offen |
| 5 | Identität am Host leer (nur Name, nur E-Mail, beide, nur Leerzeichen) | der Aufruf in §6 bricht ab, wenn eine fehlt. Frage: Reihenfolge gegen Randform 4 (welcher Grund zuerst), Meldung nennt den fehlenden Wert | offen |
| 6 | Identität mit Sonderzeichen: Anführungszeichen `"`, Apostroph `'`, Leerzeichen, `$`, Backtick, Backslash, Zeilenumbruch, Nicht-UTF-8 | `printf %q` am Host; Gegenprobe am Wert nötig. Frage: erlaubte Menge, Quoting-Form (die Shell im Container ist die Anmelde-Shell von `vscode`), ob Zeilenumbruch abgelehnt wird | offen |
| 7 | Identität nur lokal, nie global | Zusage: `--global` und `--system` werden nie gesetzt; Frage: gilt auch, wenn das Host-Repo die Werte nur aus `--global` bezieht (`git config user.name` liest alle Ebenen) | offen |
| 8 | Anderer Port, mehrere Instanzen (zweiter Remote-Name) | Dokument §10: je Instanz eigener Port und Remote (`container2` im Dokument; für `pgr-dev2` heißt das Remote `container`); der Pfad im Container trägt den Projektnamen (`/workspaces/<Projektname>/`). Frage: wie wählt das Ziel Remote-Name und Bare-Repo-Pfad je Instanz (Parameter, Ableitung aus `u-boot.yaml`); Port aus der Umgebung des Gebers (A4); Instanz im Pfad | offen |
| 9 | Hostschlüssel unbekannt, geändert (Neubau) bei den `ssh`-Aufrufen des Ziels | Warnung am Host bekannt; keine Einstellung festgelegt (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`); Neubau ändert den Fingerabdruck (A9 beim Geber). Frage: Ziel setzt nichts und bricht ab, oder die Entscheidung des Nutzers liegt vor | offen (Entscheidung des Nutzers) |
| 10 | Netz fehlt bei `make traeger-fetch` im Klon | braucht Netz im Container; das Egress bleibt eingeschränkt (Entscheidung des Nutzers vom 2026-10-17), der Host des Klons muss in der `egress.allow`-Liste stehen; `slice-harness-devcontainer-u-boot-0-10` hat für `make gates` zusätzlich `gcr.io` und `pkg-containers.githubusercontent.com` gemessen, `make traeger-fetch` im Klon ist dort nicht gelaufen (Dokument §13); Frage: Exit-Code des Ziels, ob die vorigen Schritte stehen bleiben (ja, kein Rückbau) und wie die Wiederholung aussieht | offen |
| 11 | Abbruch mitten im Ablauf (`SIGINT`, `SIGTERM`, Verbindungsabbruch) und Wiederholung | kein Rückbau; Frage: jeder Schritt wiederholbar (Randformen 1 bis 3), Meldung nennt den letzten Schritt | offen |
| 12 | `main` am Host fehlt, ist leer, steht nicht im Host-Repo oder nur auf anderem Zweig; Host-Repo im Zustand `detached HEAD` | Frage: einzuliefernder Zweig ist `main`, das Ziel hängt nicht am ausgecheckten Zweig; Meldung bei fehlendem `main` | offen |
| 13 | Push abgewiesen (Bare-Repo hat schon einen `main` mit anderer Historie) | `! [rejected]` bei nicht-fast-forward. Frage: kein `--force`, Meldung und Exit | offen |
| 14 | Ausgang mit Exit-Codes: Reihenfolge der Gründe, wenn mehrere zutreffen (Identität fehlt **und** `main` fehlt **und** Container nicht erreichbar) | Reihenfolge der Gründe bindet der Vertrag (`AGENTS.md` §3.12, Stelle einer Nachricht) | offen |
| 15 | *(verschoben)* Quelle des Zugangs von Claude Code | steht in `slice-harness-devcontainer-einrichten-zugang`, Randform 15 dort; die Nummer bleibt vergeben | — |
| 16 | Zweig je Slice im Klon (Schritt 7) | Entscheidung des Nutzers vom 2026-10-17: je Slice ein Zweig aus `main`. Fragen: woher der Name (Parameter, Pflicht?), Zweig besteht schon (Wiederholung des Ziels nach Neustart: weiter auf dem Zweig oder Abbruch?), Klon hat uncommittete Änderungen, `main` im Klon steht hinter dem Bare-Repo | offen |
| 25 | Stelle des Zugangs im Ablauf (Bindung an `slice-harness-devcontainer-einrichten-zugang`) | Der Nehmer fügt die Auflösung vor Schritt 1 ein; Frage: wie das Skript die Stelle frei hält, ohne dass dieser Slice Zugangs-Code enthält (Funktion, Reihenfolge der Prüfungen) | offen |

**Aufrufarten von `make`** (Randform des Treibers, `AGENTS.md` §3.12 *Aufrufarten des Treibers*): Das Ziel ändert Zustand außerhalb des Arbeitsbaums (Container, Remote, Konfiguration des Host-Repos). Je Flag eine Entscheidung im Vertrag und je Flag ein Fall und eine Mutation in der Gegenprobe: `-n` (Rezept mit `$(MAKE)` darf nichts ausführen: kein `ssh`, kein Push, keine Konfiguration), `-q`, `-t` (Ziel ist `.PHONY`, trägt kein Zeitstempel-Ziel), `-i` (Fehler dürfen den Ablauf nicht überlesen: kein „Einrichtung abgeschlossen“ nach einem Fehler), `MAKEFLAGS` aus der Umgebung, auch in Kombination und mit Vorrang. Dazu `-k` und `-j` (parallel laufende Schritte dürfen die Reihenfolge nicht brechen).

**Weitere Risiken:**

- **Neue Host-Abhängigkeit `ssh`.** Entscheidung des Nutzers vom 2026-10-17 („ja“); Eintrag in Vertrag und ADR durch den Architect, kein Plan-Eintrag in `AGENTS.md` §3.1 und nicht in Lastenheft oder Gate-Index, solange der Architect nichts anderes festlegt. Fehlt `ssh` am Host, meldet das Ziel es vor jeder Änderung (eigener Fall). — **Ausgang (Closure):** offen bis zur Closure; erwartet *entfallen* mit Verweis auf Vertrag und ADR.
- **Die Geber `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten-start` liefern nicht, was §1 voraussetzt** (Container, Variablen, Pfad). — **Ausgang (Closure):** entfallen, weil der Start beide Geber in `done/` verlangt (§4); ein Befund geht an den Planner.
- **Der volle Lauf von `make gates` im Container ist von `slice-harness-devcontainer-u-boot-0-10` gemessen** (grün, zwei Läufe zu je 1701 s, höchstens 5,42 GiB; Dokument §8); dieser Slice misst ihn nicht neu. Trägt der Lauf im Klon von `pgr-dev2` nicht (anderes Gerüst, andere Last), ist das ein Befund an den Planner. — **Ausgang (Closure):** entfallen, weil `slice-harness-devcontainer-u-boot-0-10` den Lauf misst und der Start seine grüne Messung verlangt (§4).
- **Gegenprobe stört einen anderen Lauf** (Docker-Daemon geteilt, `make gates` eines anderen Slice) oder den Container, in dem der Implementer arbeitet (`pgr-dev2`), oder hinterlässt Volumes (`BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`). — **Ausgang (Closure):** entfallen mit Beleg, dass die Gegenprobe nur Objekte ihres Präfixes entfernt und nach Ende keines bleibt, oder eingetreten → Folge-Slice mit Kennung.
- **Gegenprobe wächst** (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`): zwei Zusagen je Aufzählungselement. — **Ausgang (Closure):** entfallen mit gemessener Laufzeit und Mutationsprotokoll.
- **Größe über Schätzung wegen der Gegenprobe** (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Faktor 2 eingerechnet, Rückführung §4). — **Ausgang (Closure):** entfallen mit gemessener Summe unter 900, oder eingetreten → `slice-harness-devcontainer-einrichten-gegenprobe` (Name vorläufig, vom Planner bei Rückführung vergeben).
- **Konflikt an `harness/mk/gates-schnell.mk` und `harness/README.md`** mit dem Geber (beide ändern sie). — **Ausgang (Closure):** entfallen, weil der Start den Geber in `done/` verlangt (§4).
- **Die Sendung „Arbeit im Container“ wächst den Dokument-Punkt** (`BEO-REPO/slice-waechst-durch-uebernahmen`). — **Ausgang (Closure):** entfallen mit gemessener Summe nach der Metrik aus §4 (unter der Marke 900), oder eingetreten → `slice-harness-devcontainer-einrichten-gegenprobe`.
- **Das Dokument veraltet still gegen das Ziel** (kein Sensor) und Aussagen über die Probe hinaus (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). — **Ausgang (Closure):** weiter offen im selben Eintrag, oder eingetreten (Evidence dieses Slice).
- **Umbau des Dokuments verliert Vorbehalte der Vorfassung** (Sicherheits- und Isolationsaussagen aus §11, „Behandle den Container nicht als Sandbox“). — **Ausgang (Closure):** entfallen mit Vergleich der Vorbehalte der Vorfassung mit dem Diff im Review; ein Befund ist *eingetreten* und geht in `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.
- **Gegenprobe gleichzeitig mit dem Lauf eines anderen Slice auf demselben Daemon** (Abgrenzung §1; die Regel richtet sich an den Aufrufer, `SPEC-052` Punkt 8, `SPEC-058` Punkt 15; die Gegenprobe prüft sie nicht). — **Ausgang (Closure):** entfallen, wenn der Bericht nennt, dass kein anderer Lauf auf dem Daemon lief; sonst eingetreten (Befund an den Planner).
- **Arbeit des Implementers im Klon von `pgr-dev2` bricht am Egress** (Randform 10; `make traeger-fetch` im Klon ist ungemessen). — **Ausgang (Closure):** entfallen mit Beleg im Bericht, dass die Arbeit im Klon lief, oder eingetreten (Befund an den Planner, Eintrag im Register).

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie ist so grob wie das Repo, trägt aber eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer für Werkzeug und Container: `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (Randform 9), `BEO-REPO/nested-runtime-podman-oder-docker-offen`, `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Zählerstände liest der Implementer aus den Dateien in `evidence/` vor dem Start (kein Zähler wird hier gesetzt). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel; die verkörperten (§3.10 bis §3.14) gelten.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Schnitt einträgt):** Liefer-Punkte 3 (Ziel mit Ablauf, Gegenprobe, Dokument; die Sendung „Arbeit im Container“ liegt im Dokument-Punkt, die Sendung „Quellenauflösung“ liegt beim Zugang, Gerüst und Start beim Geber). Schichten 2 (Harness-Werkzeug mit Gegenprobe und Gate-Index; Nutzer- und Wartungs-Doku `docs/developer/`; `harness/README.md` ist Index desselben Werkzeugs, Planung und Vertrag sind keine Schicht). Schätzung in der Metrik von §4: 565 bis 830, **unter der Marke 900** (Meldeschwelle 700 in der oberen Hälfte überschritten). Vor dem Schnitt 995 bis 1480 (ungeteilt, mit Gerüst und Start). **Nachgezählt am 2026-10-18 im Commit, der die Aufteilung des Gebers in `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten-start` einträgt:** Die Bindung an den Container kommt jetzt von zwei Gebern; sie ist keine neue Lieferung dieses Slice. Liefer-Punkte 3 und Schichten 2 unverändert, Schätzung 565 bis 830 unverändert, unter der Marke 900. `grep` nach der Kennung dieses Slice findet ihn in den Gebern `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten-start` (§1, nehmen die Bindung auf), im Folge-Slice `slice-harness-devcontainer-einrichten-zugang`, in `slice-harness-devcontainer-u-boot-0-10` (Geber, `done/`), in den Plänen der Release-Welle (keine Bindung) und in den Beobachtungen.
