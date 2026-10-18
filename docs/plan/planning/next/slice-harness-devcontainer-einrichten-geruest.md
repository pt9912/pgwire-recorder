# Slice slice-harness-devcontainer-einrichten-geruest: Make-Ziel `container-start` legt das Gerüst des Dev-Containers an und startet ihn (Voraussetzungen, `u-boot`, Schlüsselpaar, `mk-override.sh`, `devcontainer up`, Stopp und Neustart)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein Abnahmeszenario und
kein Meilenstein hängt an ihm. Angelegt am 2026-10-18 durch den Schnitt von `slice-harness-devcontainer-einrichten` (dessen Schätzung 995 bis 1480 lag mit dem ganzen Spektrum über der Marke 900; Regel des Nutzers: schneiden, wenn die Schätzung die Marke selbst überschreitet).

**Bezug:** [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität) nur insoweit, als die Host-Voraussetzungen des Repos berührt sind (neu, noch nicht entschieden: `u-boot`, `devcontainer`-CLI, `node`, `ssh-keygen`, §1); kein Anspruch des Produkts.
Vorgänger: `slice-harness-devcontainer-austausch-repo` und `slice-harness-devcontainer-doku` (beide `done/`) und `slice-harness-devcontainer-u-boot-0-10` (`done/`; Geber des gemessenen Gerüsts auf `u-boot` v0.10.0 und der Messung von `make gates` im Container). Folge-Slices in dieser Reihenfolge: `slice-harness-devcontainer-einrichten` (Bare-Repo, Einlieferung, Klon, Zweig; setzt den Container dieses Slice voraus), dann gegebenenfalls `slice-harness-devcontainer-einrichten-zugang`. Keine aktive ADR; der Architect legt eine an (§6).

**Berührte Spec-Stellen:** neuer Vertrag eines Harness-Werkzeugs in `harness/vertraege/` (Nummer vergibt der Architect); Host-Voraussetzungen: [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität).

**Verantwortlich:** pt9912

**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Make-Ziel `container-start` (Name vorläufig, der Architect legt ihn fest) legt am Host in einem Aufruf das Gerüst des Dev-Containers nach `docs/developer/devcontainer.md` §3 an und startet den Container (§4): Voraussetzungen und Eingaben prüfen, `$PGR_DEV` mit `u-boot init … --devcontainer --sandbox --no-git`, die `u-boot config set`-Zeilen aus §3.1, `u-boot generate devcontainer`, Schlüsselpaar `ssh/id_ed25519` (§3.2), `mk-override.sh` (§3.5), `devcontainer up --workspace-folder … --override-config …`, warten bis `ssh` antwortet; dazu Stopp und Neustart (§4). Ergebnis ist ein Container, der per `ssh` erreichbar ist; was darin liegt, richtet `slice-harness-devcontainer-einrichten` ein.

**Arbeitsweise (Entscheidung des Nutzers vom 2026-10-18):** Die Hostsitzung treibt die Arbeit im Container per `ssh` und Befehl; Claude Code läuft nicht im Container (siehe Abgrenzung zum Zugang).

Ablauf des Ziels, in dieser Reihenfolge (die Schritte 1 bis 3 des Slice `slice-harness-devcontainer-einrichten` vor dessen Schnitt):

1. Eingaben und Voraussetzungen am Host prüfen, bevor etwas geändert wird: `PGR_DEV` (Randform A2), Werkzeuge `u-boot`, `devcontainer`, `node`, `ssh`, `ssh-keygen`, `docker` (Randform A1), Port `PGR_SSH_PORT` (belegt: Randform A4), Name `PGR_NAME`.
2. Gerüst anlegen, falls es fehlt (Randformen A3, A6, A7): Verzeichnis, `u-boot`-Aufrufe, Schlüsselpaar, `mk-override.sh`.
3. Container starten (Randformen A5, A8): `mk-override.sh`, `devcontainer up`, warten bis `ssh` antwortet (Zeitgrenze).

Stopp und Neustart (`docker stop <Name>`, dann Schritt 3) sind Dokument §4; ob dafür ein zweites Ziel entsteht oder der zweite Aufruf von `container-start` genügt (Randform A12), entscheidet der Architect.

**Annahme (Entscheidung des Nutzers vom 2026-10-17, „ja“): ein Slice pro Container-Sitzung.** Ein Slice läuft in einer Sitzung, danach Neustart (`docker stop`, dann `devcontainer up …`; `u-boot up` startet das gemessene Gerüst nicht, `no service selected`). Die Volumes `<Projektname>-docker-<ID>` und `<Projektname>-workspace-<ID>` bleiben über den Neustart erhalten, der Neustart ist billig (1,7 bis 1,8 s bis `ssh` und `dockerd` antworten, gemessen von `slice-harness-devcontainer-u-boot-0-10`). Das Ziel ist deshalb wiederholbar (Randform A5) und baut ein vorhandenes Gerüst nicht neu; die Namen der Volumes tragen die `devcontainerId` (`<ID>`) und ändern sich, wenn das Verzeichnis des Gerüsts wechselt.

**Sendung des Gebers `slice-harness-devcontainer-u-boot-0-10` („Arbeit im Container“, eingetragen vom Planner am 2026-10-17, Entscheidung des Nutzers vom selben Tag)** teilt sich auf: Den Weg bis zum laufenden Container trägt dieser Slice, den Weg im Container (Bare-Repo, Klon, Zweig, Ablauf eines Implementer-Laufs im Dokument) `slice-harness-devcontainer-einrichten`. Der Geber liefert das Gerüst, die Messung und die Probe von `claude`; dieser Slice baut auf diesen Befunden auf und misst sie nicht neu.

**Bindung an `slice-harness-devcontainer-einrichten` (Sendung, eingetragen vom Planner am 2026-10-18, `AGENTS.md` §3.13):** Der Nehmer setzt voraus, was dieser Slice zusagt: `$PGR_DEV` mit `u-boot.yaml` (der Projektname bestimmt den Pfad `/workspaces/<Projektname>/` im Container), ein Container, der per `ssh` als `vscode` mit `$PGR_DEV/ssh/id_ed25519` auf `127.0.0.1:${PGR_SSH_PORT:-2222}` antwortet, und die Variablen `PGR_DEV`, `PGR_SSH_PORT`, `PGR_NAME` mit den Vorgaben aus Dokument §10. Ändert dieser Slice einen dieser Punkte, folgt der Nehmer in §1, §3 und §6.

**Bindung an `slice-harness-devcontainer-einrichten-zugang` (Sendung, eingetragen vom Planner am 2026-10-18):** `mk-override.sh` schreibt dieses Ziel, wie Dokument §3.5 es zeigt, auch die Zeile mit dem Verweis `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}` (nie ein Wert); ob sie bleibt, entscheidet der Nehmer. Dieses Ziel liest, löst und übergibt keinen Zugang.

**Entscheidung des Nutzers (2026-10-17): `ssh` ist Voraussetzung am Host** (für das Ziel des Nehmers `slice-harness-devcontainer-einrichten`; dieser Slice braucht `ssh` und `ssh-keygen` ebenfalls). **Die übrigen Voraussetzungen dieses Slice — `u-boot`, die `devcontainer`-CLI, `node` (braucht die CLI und `mk-override.sh`), `ssh-keygen` — hat der Nutzer nicht entschieden.** Der Architect legt sie vor dem Code vor (§6, E2), der Nutzer entscheidet; bis dahin ist es eine Empfehlung des Planners: aufnehmen, jede nur für dieses Ziel, kein Gate und kein anderes Ziel braucht sie. Die Voraussetzungen des Repos sind bisher `git`, `docker` und `make` (`AGENTS.md` §3.1).

**Kein Gate**, wie `hooks-install`: Das Ziel steht im Abschnitt *Werkzeuge* von `harness/README.md`, nicht in `GATE_CHECKS`. Seine Gegenprobe ist ein Gate (§2).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Bare-Repo, keine Einlieferung, kein Klon, keine Identität, kein `make hooks-install`, kein `make traeger-fetch`, kein Zweig, kein Remote und kein `core.sshCommand` im Host-Repo: Das trägt `slice-harness-devcontainer-einrichten` (Klasse 1, Folge-Slice; dessen §1 nennt diese Schritte und die Bindung dieses Slice). Dieses Ziel berührt das Host-Repo nicht.
- Kein Zugang von Claude Code im Container (Quellenauflösung, `remoteEnv`-Weg über `devcontainer exec`, Variablen `ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`, Dokument §7): `slice-harness-devcontainer-einrichten-zugang` (Klasse 1; Entscheidung des Nutzers vom 2026-10-18, der Zugang wird nicht gebraucht).
- Kein Umbau des Gerüsts auf eine andere Version von `u-boot` und keine Messung von `make gates` im Container: Das liefert der Vorgänger `slice-harness-devcontainer-u-boot-0-10` (`done/`). Das Ziel automatisiert dessen gemessenen Ablauf (Dokument §3 bis §5) und setzt den Inhalt von `u-boot.yaml` und `mk-override.sh` aus dem Dokument um; trägt der Inhalt den Ablauf nicht, hält der Implementer an (§6).
- Keine Entscheidung zwischen `podman` und `docker` als `nestedRuntime` (`BEO-REPO/nested-runtime-podman-oder-docker-offen`): Das Ziel schreibt den Wert, den das Dokument nennt (`docker`).
- Keine Entscheidung zum Hostschlüssel (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`): Das Ziel nennt, was `ssh` ausgibt, und setzt keine Einstellung (`StrictHostKeyChecking`, `UserKnownHostsFile`) ohne Entscheidung des Nutzers (Randform A9).
- Kein Löschen von Container, Volumes oder Image (Dokument §4, *Container entfernen*, *Alles löschen*, *Image löschen*): Handlung des Nutzers; die Gegenprobe entfernt nur Objekte ihres Präfixes per Name.
- Kein Produktcode (Schicht-Abgrenzung): Es ändern sich Werkzeug, Gegenprobe und Dokumentation, nichts unter `cmd/` und `internal/`.
- Kein Eintrag in der Roadmap: Der Slice ist wellenlos. Die Release-Slices und die Closure von `welle-v1-abschluss` hängen nicht an ihm (Entscheidung des Nutzers vom 2026-10-17).

## 2. Definition of Done

Liefer-Punkte (drei):

- [ ] **Ziel `container-start`** (Fragment unter `harness/mk/` und Skript unter `tools/harness/`): erfüllt den Vertrag, den der Architect in `harness/vertraege/` festlegt (Ablauf aus §1, Randformen aus §6); kein Gate; in `harness/README.md` §Werkzeuge genannt.
- [ ] **Gegenprobe** `container-start-gegenprobe` (Gate, hängt an `GATE_CHECKS`, steht in der Ausnahmeliste der schnellen Kette): gegen einen Scratch-Container (Präfix, Entfernen nur per Name, `docker rm -fv`, nichts Fremdes, nie ein Container, in dem gearbeitet wird) je Randform aus §6 ein Fall und je Zusage und je Element einer Aufzählung eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die Gegenprobe steht im Gate-Index.
- [ ] **Dokument** `docs/developer/devcontainer.md` §3 bis §5 nennen `make container-start` als Weg und führen die Einzelbefehle nur, soweit das Dokument sie als Erklärung braucht; die Host-Voraussetzungen in §2 sind nachgeführt. Jede Aussage ist durch einen Lauf belegt oder steht in §13. §6 bis §9 gehören `slice-harness-devcontainer-einrichten`, §7 `slice-harness-devcontainer-einrichten-zugang`.

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
| `harness/vertraege/SPEC-<NNN>-container-start.md` | neu (Architect) | Vertrag des Werkzeugs: Ablauf, Randformen aus §6, Host-Voraussetzungen; Nummer vergibt der Architect |
| `docs/plan/adr/` (eine ADR) | neu (Architect) | Entscheidung: Ort des Ziels (Repo oder Host-Werkzeug), Host-Voraussetzungen (Entscheidung des Nutzers); `Schärft:` zeigt auf den Vertrag (`AGENTS.md` §3.8) |
| `harness/mk/container-start.mk` | neu | Ziel `container-start` (kein Gate); Gegenprobe hängt an `GATE_CHECKS` |
| `tools/harness/container-start.sh` | neu | Ablauf aus §1 (Prüfen, Gerüst, Start); Eingaben aus der Umgebung (`PGR_DEV`, `PGR_SSH_PORT`, `PGR_NAME`, `PGR_CPUS`, `PGR_MEM`); Fehlermeldungen mit Schritt; keine Zugangsvariable |
| `tools/harness/container-start-gegenprobe.sh`, `harness/mk/container-start-gegenprobe.mk` | neu | Gegenprobe gegen Scratch-Container; Fälle und Mutationen aus §6 |
| `harness/mk/gates-schnell.mk` | update | Gegenprobe in die Ausnahmeliste der schnellen Kette (Docker; Pflege nach `SPEC-052`, Architect prüft Punkt der Gate-Kette) |
| `harness/README.md` | update | Abschnitt *Werkzeuge* (`container-start`) und *Sensors* (Gegenprobe) |
| `docs/developer/devcontainer.md` | update | §2 und §3 bis §5 nennen das Ziel; §13 für Ungeprüftes |

- Das Skript trennt Prüfen und Handeln: erst alle Vorbedingungen am Host, dann Gerüst, dann Start. Die Reihenfolge der Gründe legt der Architect fest (Randform A11).
- **Docker-Last.** Die Gegenprobe startet Scratch-Container (Präfix `pgr-cs-` oder wie der Vertrag es festlegt; Entfernen nur per Name, `docker rm -fv` und Prüfung, dass keines der Volumes bleibt, wie `SPEC-057`). Sie läuft nie gleichzeitig mit `make gates`; die Vorbedingung `pgrep -fa "^make gates"` leer ist Teil der Gegenprobe (Abbruch oder Warten legt der Vertrag fest). Ein Erststart baut ein Image (Dokument §4: 65 s bei warmem Cache); die Laufzeit der Gegenprobe ist zu messen (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`).
- **Arbeit des Implementers im Dev-Container `pgr-dev2`** (Entscheidung des Nutzers vom 2026-10-18; Stand 2026-10-18: Container `pgr-dev2` läuft, Host-Verzeichnis `/Development/pgr-dev2`, Port 2223, Klon `/workspaces/pgr-dev2/pgwire-recorder`, Remote `container` im Host-Repo, Hostschlüssel-Datei `/Development/pgr-dev2/known_hosts`, kein Claude Code im Container). Die Hostsitzung treibt den Klon per `ssh`:
  1. Jeder Befehl: `ssh -i /Development/pgr-dev2/ssh/id_ed25519 -o IdentitiesOnly=yes -o UserKnownHostsFile=/Development/pgr-dev2/known_hosts -p 2223 vscode@127.0.0.1 'cd /workspaces/pgr-dev2/pgwire-recorder && <Befehl>'`. Die Datei hat der Nutzer bereitgestellt; sie ist eine Einstellung des Implementers für seine eigenen Befehle, das Ziel dieses Slice setzt sie nicht (Abgrenzung §1).
  2. Zweig `slice-harness-devcontainer-einrichten-geruest` im Klon, nicht `main`; Dateien als Patch über `stdin` (`git apply`) oder als Inhalt über `stdin`; Commits im Klon mit der Kennung im Text; `make hooks-install` und `make traeger-fetch` im Klon, falls noch nicht geschehen.
  3. `make gates-schnell` und `make test` im Klon; `make gates` (28 min 21 s laut Dokument §8, für diesen Stand ungeprüft) im Hintergrund (`setsid`, Log in `/tmp` im Container, Exit-Code in eine Datei), Abfrage per `ssh`; `pgrep -fa "^make gates"` im Container ist vor jedem Lauf der Gegenprobe leer.
  4. *Was nur der Host kann.* Die Gegenprobe legt Scratch-Gerüste an und braucht `u-boot`, `devcontainer`, `node`. Ob sie im Container von `pgr-dev2` vorliegen, ist ungemessen: Der Implementer prüft `command -v u-boot devcontainer node ssh ssh-keygen` im Container zuerst und nennt das Ergebnis im Bericht. Fehlen sie, läuft der Teil der Gegenprobe, der sie braucht, am Host mit Scratch-Verzeichnis, Container-Namen und Volumes im Präfix der Gegenprobe (nie `pgr-dev2` anfassen, nie stoppen, neu starten oder entfernen, auch nicht dessen Port 2223 belegen); der Bericht sagt, was im Container und was am Host lief.
  5. Ergebnis zum Host: im Klon `git push origin <Zweig>`; am Host `git fetch container`, der Zweig liegt als `container/<Zweig>` zur Prüfung (Dokument §9). Der Inhalt stammt aus einem Container, in dem `vscode` `root` werden kann: Review und Verifikation sehen die Commits durch, bevor sie nach `main` gehen. Der Übergang `next/` → `in-progress/` ist ein Commit des Nutzers auf `main` vor der Arbeit und wird vor dem ersten Klon-Stand eingeliefert (`git push container main`).
  6. Größenprobe: Die Metrik von §4 misst der Implementer im Klon mit `git diff --numstat <Basis>..<letzter Commit>`.
- Zählung nach `AGENTS.md` §3.13 siehe §8.

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-harness-devcontainer-u-boot-0-10` liegt in `done/` (erfüllt, Messung grün). Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang. Reihenfolge: dieser Slice, dann `slice-harness-devcontainer-einrichten`. Vorbedingung des Starts außerdem: Vertrag und ADR des Architect liegen vor (`Accepted` ist die Entscheidung des Nutzers), und der Nutzer hat die Host-Voraussetzungen entschieden (§1).

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; Vertrag und ADR (Commits des Architect vor dem Start) zählen nicht.

**Schätzung (Planner, 2026-10-18; Metrik wie oben, Gegenprobe mit dem Doppelten des Produktcodes, Faktor aus `AGENTS.md` §3.14):**

| Teil | Zeilen |
|---|---|
| Skript: Prüfen der Eingaben und Voraussetzungen | 25 bis 40 |
| Skript: Gerüst (Verzeichnis, `u-boot`-Aufrufe, Schlüsselpaar, `mk-override.sh`) | 60 bis 90 |
| Skript: Start (`devcontainer up`, Warten auf `ssh`) | 30 bis 50 |
| Skript: Stopp und Neustart | 15 bis 25 |
| Fragmente | 30 bis 40 |
| Gegenprobe (Doppeltes von Skript und Fragmenten, 160 bis 245) | 320 bis 490 |
| Dokument (§2 bis §5 gekürzt, Löschungen zählen, Index) | 80 bis 130 |
| **Zusammen** | **560 bis 865** |

Meldeschwelle 700, **Marke 900**. Die obere Hälfte liegt über der Meldeschwelle, die ganze Schätzung unter der Marke; ein Schnitt unterbleibt.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder der Vertrag verlangt mehr als drei Liefer-Punkte oder mehr als zwei Schichten (zum Beispiel eine Änderung an `.devcontainer/Dockerfile`), oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner (Vorschlag: Start mit Stopp/Neustart getrennt vom Gerüst).
- `in-progress` → `open` (blockiert): Die Host-Voraussetzungen `u-boot` und `devcontainer` bleiben unentschieden (Nutzer), oder das Gerüst aus dem Dokument trägt den Ablauf nicht ohne Eingriff am Host (Docker-Version, Netz; der Befund geht an den Planner).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün auf dem Endstand, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben. Das Mutationsprotokoll (Zusage · Mutation · roter Test, `AGENTS.md` §3.10) steht im Bericht des Implementers. Die Gegenprobe ist einmal gegen einen echten Scratch-Container grün gesehen, nicht nur gegen Attrappen.

## 6. Risiken und offene Punkte

Der Slice liefert einen **neuen Vertrag** (Werkzeug mit Zusagen); §3.12 gilt. Die Randformen stehen hier, entschieden werden sie **vor dem ersten Code-Commit** in `harness/vertraege/` (Architect); eine Form, die erst ein Commit mit Code oder Tests hier einträgt, ist nicht genannt. Der Implementer entscheidet keine still, er hält an und gibt sie dem Architect zurück.

**Entscheidungen des Architect vor dem Code:**

- **E1 Ort des Ziels.** Make-Ziel im Repo (`harness/mk/`, `tools/harness/`, wie §3 es annimmt) oder Host-Werkzeug außerhalb des Repos (wie `u-boot` selbst)? Das Ziel schreibt Verzeichnisse außerhalb des Repos (`$PGR_DEV`) und berührt das Host-Repo nicht. Die Antwort bestimmt Vertrag, Gate-Index, Gegenprobe und Größe und bindet `slice-harness-devcontainer-einrichten` (dessen E1).
- **E2 Host-Voraussetzungen.** `u-boot` (Version 0.10.0 oder höher?), `devcontainer`-CLI, `node` (ob `mk-override.sh` ohne `node` auskommt, ist eine Frage des Inhalts), `ssh-keygen`. Der Architect legt vor, der Nutzer entscheidet. **Empfehlung des Planners:** alle vier aufnehmen, jede nur für diese Ziele; streicht der Nutzer eine, entfällt der Automatismus für das Gerüst (Gerüst von Hand nach Dokument §3, das Ziel prüft nur und startet), und der Schnitt nach §4 ist neu zu schätzen.
- **E3 Randformen des Aufrufs** (A1 bis A10): fehlende Variable, vorhandenes Verzeichnis, belegter Port, zweiter Aufruf, Werkzeug fehlt, Schlüsselpaar.
- **E4 Aufrufarten von `make`** (Randform des Treibers, `AGENTS.md` §3.12), siehe unten.
- **E5 Reihenfolge der Gründe** (A11) und Exit-Codes.
- **E6 Stopp/Neustart** (A12): zweites Ziel oder zweiter Aufruf.
- **E7 Schätzung** (§4): Der Architect bestätigt sie in der Metrik vor dem Code.

**Randformen des Ablaufs** (Spalte *Wo entschieden*: Fundstelle im Vertrag, vom Architect einzutragen):

| # | Randform | Vorwissen aus `docs/developer/devcontainer.md` | Wo entschieden |
|---|---|---|---|
| A1 | Werkzeug am Host fehlt oder hat die falsche Version (`u-boot`, `devcontainer`, `node`, `ssh`, `ssh-keygen`, `docker`) | Dokument §2 nennt die Voraussetzungen; Frage: je Werkzeug ein eigener Fall und eine Meldung mit dem Namen; Version von `u-boot` (v0.10.0 gemessen); vor jeder Änderung | offen (E2, Entscheidung des Nutzers) |
| A2 | `PGR_DEV` fehlt, ist leer, relativ, enthält Leerzeichen, liegt im Host-Repo, ist kein Verzeichnis | Dokument §3: absoluter Pfad außerhalb des Repos; Frage: Pflicht oder Vorgabe, Prüfung der Form, Pfad im Arbeitsbaum des Repos abgelehnt | offen |
| A3 | Verzeichnis `$PGR_DEV` besteht schon (leer; mit Gerüst; mit fremdem Inhalt; mit Gerüst anderer Version von `u-boot`) | Dokument §3.1 legt mit `mkdir` an, bricht bei vorhandenem Verzeichnis ab; Frage: leer erlaubt, Gerüst erkennen und weiter, fremder Inhalt abgelehnt; Blöcke zwischen `BEGIN`/`END U-BOOT MANAGED BLOCK` nicht überschreiben | offen |
| A4 | Port `PGR_SSH_PORT` belegt (anderer Container, fremder Dienst), leer, nicht numerisch, außerhalb 1 bis 65535; Vorgabe 2222 | Dokument §10: ein belegter Port ist nicht gelaufen (§13); `PGR_SSH_PORT` liest nur `mk-override.sh` aus der Shell, keine `.env`; Frage: Prüfung vor dem Start oder Meldung von `devcontainer up`, Unterschied belegt durch denselben Container (zweiter Aufruf) und durch einen fremden | offen |
| A5 | Zweiter Aufruf: Container läuft schon; Container gestoppt; Container mit dem Namen besteht aus anderem Gerüst; `PGR_NAME` kollidiert | Dokument §4: Neustart 1,7 bis 1,8 s; `devcontainer up` auf laufendem Container ist ungemessen; Frage: erster und zweiter Aufruf unterscheiden, Neustart ohne Neubau, Kollision mit fremdem Namen abgelehnt | offen |
| A6 | Schlüsselpaar `ssh/id_ed25519` besteht schon (mit/ohne Passphrase, nur ein Teil, falsche Rechte) | Dokument §3.2: `ssh-keygen -N ''`; Überschreiben vernichtet den Zugang; Frage: nie überschreiben, Teil fehlt → Abbruch, Rechte 700/600 | offen |
| A7 | `u-boot` meldet Warnungen oder Fehler (Lockerungen des Profils, `doctor`-Warnung `no service selected`, Exit 12 von `u-boot up`) | Dokument §3.1, §12: Warnungen sind erwartet; Frage: welche Warnung und welcher Exit gelten als Erfolg, was bricht ab; die Ausgabe von `u-boot` erscheint unverändert | offen |
| A8 | Container startet nicht / `ssh` antwortet nicht: Zeitgrenze, `devcontainer up` meldet Erfolg, `sshd` antwortet nicht, `dockerd` nicht bereit | Dokument §4: Bereit ist der Container, wenn `ssh` antwortet; Erststart 65 s, Neustart 1,7 bis 1,8 s (Cache warm); ohne warmen Cache ungemessen; Frage: Zeitgrenze, Exit, ob `dockerd` Teil der Bereitschaft ist, was bleibt stehen (kein Rückbau) | offen |
| A9 | Hostschlüssel unbekannt oder geändert (Neubau) beim Warten auf `ssh` | Das Feature `sshd` erzeugt den Hostschlüssel im Build; Neustart und neuer Container aus demselben Image ändern ihn nicht, jeder Neubau (`--build-no-cache`, Änderung von `egress.allow`) ändert ihn (`slice-harness-devcontainer-u-boot-0-10`); keine Einstellung festgelegt (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`); Frage: Ziel setzt nichts und meldet, oder die Entscheidung des Nutzers liegt vor | offen (Entscheidung des Nutzers) |
| A10 | Abbruch mitten im Ablauf (`SIGINT`, `SIGTERM`, Verbindungsabbruch), auch mitten in `u-boot` oder `devcontainer up` | kein Rückbau; Frage: jeder Schritt wiederholbar (A3, A5, A6), Meldung nennt den letzten Schritt | offen |
| A11 | Reihenfolge der Gründe, wenn mehrere zutreffen (Werkzeug fehlt **und** `PGR_DEV` ungültig **und** Port belegt **und** Verzeichnis besteht) | Reihenfolge der Gründe bindet der Vertrag (`AGENTS.md` §3.12, Stelle einer Nachricht) | offen |
| A12 | Stopp und Neustart: Stopp eines nicht laufenden Containers; Neustart nach `docker rm` (Volumes bleiben, `~/.claude` und `~/.gitconfig` sind weg, Dokument §4); Neustart nach Änderung von `egress.allow` (Neubau, A9) | Dokument §4, gemessen: 1,7 bis 1,8 s Neustart, 6,4 und 7,9 s nach `docker rm`; Frage: zweites Ziel oder zweiter Aufruf (E6), Exit-Codes | offen |

**Aufrufarten von `make`** (Randform des Treibers, `AGENTS.md` §3.12 *Aufrufarten des Treibers*): Das Ziel ändert Zustand außerhalb des Arbeitsbaums (Verzeichnisse am Host, Container). Je Flag eine Entscheidung im Vertrag und je Flag ein Fall und eine Mutation in der Gegenprobe: `-n` (Rezept mit `$(MAKE)` darf nichts ausführen: kein `u-boot`, kein `ssh-keygen`, kein `devcontainer up`), `-q`, `-t` (Ziel ist `.PHONY`, trägt kein Zeitstempel-Ziel), `-i` (Fehler dürfen den Ablauf nicht überlesen: kein „Container läuft“ nach einem Fehler), `MAKEFLAGS` aus der Umgebung, auch in Kombination und mit Vorrang. Dazu `-k` und `-j`.

**Weitere Risiken:**

- **Neue Host-Abhängigkeiten** `u-boot`, `devcontainer`-CLI, `node`, `ssh-keygen` (§1, A1). — **Ausgang (Closure):** entfallen mit Verweis auf Vertrag und ADR, in denen der Nutzer sie entschieden hat; streicht er eine, eingetreten → Rückführung nach `open/` (§4).
- **Gegenprobe stört einen anderen Lauf** (Docker-Daemon geteilt, `make gates` eines anderen Slice) oder den Container, in dem der Implementer arbeitet (`pgr-dev2`), oder hinterlässt Volumes (`BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`). — **Ausgang (Closure):** entfallen mit Beleg, dass die Gegenprobe nur Objekte ihres Präfixes entfernt und nach Ende keines bleibt, oder eingetreten → Folge-Slice mit Kennung.
- **Gegenprobe wächst und dauert** (Image-Bau je Fall; `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`). — **Ausgang (Closure):** entfallen mit gemessener Laufzeit und Mutationsprotokoll.
- **Größe über Schätzung wegen der Gegenprobe** (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Faktor 2 eingerechnet, Rückführung §4). — **Ausgang (Closure):** entfallen mit gemessener Summe unter 900, oder eingetreten → Folge-Slice mit Kennung (vom Planner bei Rückführung vergeben).
- **Konflikt an `harness/mk/gates-schnell.mk` und `harness/README.md`** mit `slice-harness-devcontainer-einrichten` (beide ändern sie). — **Ausgang (Closure):** entfallen, weil die Reihenfolge den Nehmer hinter diesen Slice setzt (Start-Trigger dort).
- **Das Dokument veraltet still gegen das Ziel** (kein Sensor; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). — **Ausgang (Closure):** weiter offen im selben Eintrag, oder eingetreten (Evidence dieses Slice).
- **Umbau des Dokuments verliert Vorbehalte der Vorfassung** (Sicherheitsaussagen aus §11, „Behandle den Container nicht als Sandbox“). — **Ausgang (Closure):** entfallen mit Vergleich der Vorbehalte der Vorfassung mit dem Diff im Review; ein Befund ist *eingetreten* und geht in `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.
- **Docker-Daemon-Probe nur ohne laufendes `make gates`.** — **Ausgang (Closure):** entfallen, wenn `pgrep -fa "^make gates"` vor jeder Probe leer war und das im Bericht steht.
- **Arbeit des Implementers im Klon von `pgr-dev2` bricht am Egress oder am Hostschlüssel** (A9; Hosts für `make traeger-fetch` und `make gates` sind in `egress.allow` gemessen, Dokument §8). — **Ausgang (Closure):** entfallen mit Beleg im Bericht, dass die Arbeit im Klon lief, oder eingetreten (Befund an den Planner, Eintrag im Register).

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (A9), `BEO-REPO/nested-runtime-podman-oder-docker-offen`, `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Zählerstände liest der Implementer aus den Dateien in `evidence/` vor dem Start (kein Zähler wird hier gesetzt). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel; die verkörperten (§3.10 bis §3.14) gelten.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Schnitt einträgt):** Liefer-Punkte 3 (Ziel mit Ablauf, Gegenprobe, Dokument). Schichten 2 (Harness-Werkzeug mit Gegenprobe und Gate-Index; Nutzer- und Wartungs-Doku `docs/developer/`; das Gerüst unter `$PGR_DEV` liegt außerhalb des Repos und ist keine Schicht, Planung und Vertrag sind keine Schicht). Schätzung in der Metrik von §4: 560 bis 865, **unter der Marke 900** (Meldeschwelle 700 in der oberen Hälfte überschritten). Dieser Slice ist Geber der Bindungen an `slice-harness-devcontainer-einrichten` und `slice-harness-devcontainer-einrichten-zugang`; beide Nehmer nennen diese Kennung in §1 und zählen in ihrem §8 nach.
