# Slice slice-harness-devcontainer-einrichten-geruest: Make-Ziel `container-start` prüft Eingaben und Voraussetzungen und legt das Gerüst des Dev-Containers an (Teil I von `SPEC-058`: `u-boot`, Schlüsselpaar, `mk-override.sh`; Docker ändert nichts)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-dev-container (Entscheidung des Nutzers vom 2026-10-18; bis dahin ohne Welle). Angelegt am 2026-10-18 durch den Schnitt von `slice-harness-devcontainer-einrichten` (dessen Schätzung 995 bis 1480 lag mit dem ganzen Spektrum über der Marke 900; Regel des Nutzers: schneiden, wenn die Schätzung die Marke selbst überschreitet). **Am 2026-10-18 ein zweites Mal geschnitten:** Nach dem Vertrag `SPEC-058` schätzte der Architect das ganze Ziel auf rund 1000 bis 1450 Zeilen (Metrik `AGENTS.md` §3.14), über der Marke 900. Der Schnitt liegt an der Naht des Vertrags: Teil I (Punkte 2 bis 8) bleibt in diesem Slice, Teil II (Punkte 9 bis 12: Container, Port, Bereitschaft, Abbruch) liegt in `slice-harness-devcontainer-einrichten-start`. Der Anspruch dieses Slice bleibt (`in-progress/`, kein Move).

**Bezug:** [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität) nur insoweit, als die Host-Voraussetzungen des Repos berührt sind (`u-boot`, `devcontainer`-CLI, `node`, `ssh`, `ssh-keygen`, nur für dieses Ziel; entschieden, §1); kein Anspruch des Produkts. [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md) (Accepted 2026-10-18).
Vorgänger: `slice-harness-devcontainer-austausch-repo` und `slice-harness-devcontainer-doku` (beide `done/`) und `slice-harness-devcontainer-u-boot-0-10` (`done/`; Geber des gemessenen Gerüsts auf `u-boot` v0.10.0 und der Messung von `make gates` im Container). Folge-Slices in dieser Reihenfolge: `slice-harness-devcontainer-einrichten-start` (Teil II desselben Ziels), `slice-harness-devcontainer-einrichten` (Bare-Repo, Einlieferung, Klon, Zweig; setzt den Container voraus), dann gegebenenfalls `slice-harness-devcontainer-einrichten-zugang`.

**Berührte Spec-Stellen:** [`SPEC-058`](../../../../harness/vertraege/SPEC-058-container-start.md) Teil I (Punkte 2 bis 8) und, soweit sie Teil I betreffen, die Punkte 1, 4, 13, 14, 15, 16 und 17; [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) (Ausnahmeliste der schnellen Kette, Punkt 8 für die Regel an den Aufrufer); [`SPEC-057`](../../../../harness/vertraege/SPEC-057-integration-aufraeumen.md) (Vorbild des Aufräumens); [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität).

**Verantwortlich:** pt9912 (Zweig slice-harness-devcontainer-einrichten-geruest)

**Autor:** Planner. **Datum:** 2026-10-18, zweiter Schnitt 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Make-Ziel `container-start` ([ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md), Ort im Repo; Vertrag `SPEC-058`) prüft am Host in einem Aufruf Eingaben und Voraussetzungen und legt das Gerüst des Dev-Containers nach `docs/developer/devcontainer.md` §3 an, ohne dass Docker etwas ändert: die Schritte `eingaben`, `werkzeuge`, `zustand` (Verzeichnis `$PGR_DEV` und Schlüsselpaar nur lesen), `geruest` (`u-boot init`, `config set`, `generate devcontainer`), `schluessel` (`ssh/id_ed25519`) und `override` (`mk-override.sh`, `override.json`), mit den Exit-Codes 2 bis 5 und der Reihenfolge der Gründe aus `SPEC-058` Punkt 13 (die Nachbarpaare ohne Container). Dazu der Rahmen der Gegenprobe `container-start-gegenprobe` (Gruppe Attrappen, Ersetzen genau einer Stelle für Mutanten, Aufräumen der Temp-Verzeichnisse, Gate-Kette) und der Fall gegen das echte `u-boot` aus Punkt 7. Der Container selbst (Zustand, `devcontainer up`, Port, Bereitschaft, Abbruch) ist Teil II.

**Arbeitsweise (Entscheidung des Nutzers vom 2026-10-18):** Die Hostsitzung treibt die Arbeit im Container per `ssh` und Befehl; Claude Code läuft nicht im Container (siehe Abgrenzung zum Zugang).

**Entscheidungen vor dem Code (Architect, Nutzer):** [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md) ist `Accepted` (Entscheidung des Nutzers vom 2026-10-18: Ort im Repo; Host-Voraussetzungen `u-boot` ab 0.10.0, `devcontainer`-CLI, `node`, `ssh-keygen`; `ssh` hatte er am 2026-10-17 entschieden). `SPEC-058` legt Ablauf, Randformen, Gegenprobe und Grenzen fest. Die Voraussetzungen gelten nur für `container-start` und seine Gegenprobe; kein Gate und kein anderes Ziel braucht sie.

**Bindung an `slice-harness-devcontainer-einrichten-start` (Sendung, eingetragen vom Planner am 2026-10-18, `AGENTS.md` §3.13):** Der Nehmer liefert Teil II und erweitert Skript, Fragment und Gegenprobe dieses Slice. Er setzt voraus, was dieser Slice zusagt: die Schritte bis `override` in der Reihenfolge von Punkt 4, die Exit-Codes 2 bis 5 (Punkt 13), `$PGR_DEV` mit Gerüst, Schlüsselpaar und `mk-override.sh` (Punkte 5 bis 8), die Eingaben mit ihren Vorgaben und Formen einschließlich `PGR_WARTEN` (Punkt 2), und den Rahmen der Gegenprobe (Attrappen, Ersetzen genau einer Stelle, Aufräumen der Temp-Verzeichnisse bei Ende, `SIGINT`, `SIGTERM`). Ändert dieser Slice einen dieser Punkte, folgt der Nehmer in §1, §3 und §6.

**Zwischenstand bis `slice-harness-devcontainer-einrichten-start`:** `SPEC-058` sagt nicht, was das Ziel nach `override` tut, solange Teil II fehlt (Exit-Code 0 heißt „bereit“, Punkt 13). Entschieden (Z1, `SPEC-058` Punkt 13): Der Lauf endet nach `override` mit Exit-Code 6 und der Zeile `container-start: override: Teil II fehlt, Gerüst liegt, kein Container gestartet`, nie mit 0 und nie mit `bereit`. Der Zwischenstand ist nicht zugesagt: kein Abbruchfall (Punkt 12, Z3), kein Container, keine Zeile `bereit`.

**Bindung an `slice-harness-devcontainer-einrichten` (Sendung, eingetragen vom Planner am 2026-10-18, `AGENTS.md` §3.13):** Der Nehmer setzt voraus, was dieser Slice zusagt: `$PGR_DEV` mit `u-boot.yaml` (der Projektname, der Basisname von `$PGR_DEV`, bestimmt den Pfad `/workspaces/<Projektname>/` im Container) und die Variablen `PGR_DEV`, `PGR_SSH_PORT`, `PGR_NAME` mit den Vorgaben aus Dokument §10 (`SPEC-058` Punkt 2). Den Container, der per `ssh` als `vscode` mit `$PGR_DEV/ssh/id_ed25519` auf `127.0.0.1:${PGR_SSH_PORT:-2222}` antwortet, sagt `slice-harness-devcontainer-einrichten-start` zu. Ändert dieser Slice einen dieser Punkte, folgt der Nehmer in §1, §3 und §6.

**Bindung an `slice-harness-devcontainer-einrichten-zugang` (Sendung, eingetragen vom Planner am 2026-10-18):** `mk-override.sh` schreibt dieses Ziel, wie Dokument §3.5 und `SPEC-058` Punkt 8 (f) es zeigen, auch die Zeile mit dem Verweis `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}` (nie ein Wert); ob sie bleibt, entscheidet der Nehmer. Dieses Ziel liest, löst und übergibt keinen Zugang.

**Sendung des Gebers `slice-harness-devcontainer-u-boot-0-10` („Arbeit im Container“, eingetragen vom Planner am 2026-10-17, Entscheidung des Nutzers vom selben Tag)** teilt sich auf: Den Weg bis zum laufenden Container tragen dieser Slice (Gerüst) und `slice-harness-devcontainer-einrichten-start` (Start), den Weg im Container (Bare-Repo, Klon, Zweig, Ablauf eines Implementer-Laufs im Dokument) `slice-harness-devcontainer-einrichten`. Der Geber liefert das Gerüst, die Messung und die Probe von `claude`; dieser Slice baut auf diesen Befunden auf und misst sie nicht neu.

**Kein Gate**, wie `hooks-install`: Das Ziel steht im Abschnitt *Werkzeuge* von `harness/README.md`, nicht in `GATE_CHECKS`. Seine Gegenprobe ist ein Gate (§2).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Container, kein `devcontainer up`, kein Port, keine Probe der Bereitschaft, kein Abbruchfall, kein Stopp und Neustart, keine Fälle gegen reales Docker mit Container (`SPEC-058` Punkte 9 bis 12 und 15 (b) bis (d)): `slice-harness-devcontainer-einrichten-start` (Klasse 1, Nehmer; dessen §1 nimmt die Punkte an). Docker ruft dieser Slice nur mit `info` und `version` (Punkt 4).
- Kein Dokument: `docs/developer/devcontainer.md` §2 bis §5 nennt das Ziel erst, wenn es startet (`slice-harness-devcontainer-einrichten-start`, Klasse 1; Doku-Folge-Slice nach `AGENTS.md` §3.11 und §3.13). Bis dahin beschreibt das Dokument den Weg von Hand, wie er gemessen ist.
- Kein Bare-Repo, keine Einlieferung, kein Klon, keine Identität, kein `make hooks-install`, kein `make traeger-fetch`, kein Zweig, kein Remote und kein `core.sshCommand` im Host-Repo: `slice-harness-devcontainer-einrichten` (Klasse 1; dessen §1 nennt diese Schritte und die Bindung dieses Slice). Dieses Ziel berührt das Host-Repo nicht (`SPEC-058` Punkt 1).
- Kein Zugang von Claude Code im Container (Quellenauflösung, `remoteEnv`-Weg über `devcontainer exec`, Variablen `ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`, Dokument §7): `slice-harness-devcontainer-einrichten-zugang` (Klasse 1; Entscheidung des Nutzers vom 2026-10-18, der Zugang wird nicht gebraucht).
- Kein Umbau des Gerüsts auf eine andere Version von `u-boot` und keine Messung von `make gates` im Container: Das liefert der Vorgänger `slice-harness-devcontainer-u-boot-0-10` (`done/`). Das Ziel automatisiert dessen gemessenen Ablauf (Dokument §3) und setzt den Inhalt von `u-boot.yaml` und `mk-override.sh` aus dem Dokument um; trägt der Inhalt den Ablauf nicht, hält der Implementer an (§6).
- Keine Entscheidung zwischen `podman` und `docker` als `nestedRuntime` (`BEO-REPO/nested-runtime-podman-oder-docker-offen`): Das Ziel schreibt den Wert, den der Vertrag nennt (`docker`).
- Keine Entscheidung zum Hostschlüssel (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`): Sache der Probe in Teil II (`slice-harness-devcontainer-einrichten-start`).
- Kein Löschen von Container, Volumes oder Image und kein Rückbau eines Gerüsts (`SPEC-058` Punkte 1 und 12): Handlung des Nutzers; die Gegenprobe entfernt nur Objekte ihres Präfixes per Name.
- Kein Produktcode (Schicht-Abgrenzung): Es ändern sich Werkzeug, Gegenprobe und Index, nichts unter `cmd/` und `internal/`.
- Kein Eintrag in der Roadmap: Der Slice ist wellenlos. Die Release-Slices und die Closure von `welle-v1-abschluss` hängen nicht an ihm (Entscheidung des Nutzers vom 2026-10-17).

## 2. Definition of Done

Liefer-Punkte (zwei):

- [ ] **Ziel `container-start`, Teil I** (Fragment `harness/mk/container-start.mk` und Skript `tools/harness/container-start.sh`): erfüllt `SPEC-058` Punkte 1 bis 8 mit den Teilen der Punkte 13 und 14, die ohne Container gelten (Schritte `eingaben` bis `override`, Exit-Codes 2 bis 5, Rezept ohne `$(MAKE)`); kein Gate; in `harness/README.md` §Werkzeuge genannt (Zustand: legt das Gerüst an).
- [ ] **Gegenprobe** `container-start-gegenprobe`, Gruppe Attrappen und der Fall gegen das echte `u-boot` (Gate, hängt an `GATE_CHECKS`, steht in der Ausnahmeliste der schnellen Kette, Fall `p16-gate-kette`; im Gate-Index): je Randform aus §6 ein Fall und je Zusage und je Element einer Aufzählung eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); Rahmen nach `SPEC-058` Punkt 15 (Attrappen mit Protokoll, Ersetzen genau einer Stelle, Aufräumen der Temp-Verzeichnisse bei Ende, `SIGINT`, `SIGTERM`, kein Name ohne Präfix `pgr-cs-`).

Gates und Closure:

- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `tools/`, `harness/` und das Makefile-Fragment ändern). Die Regel, nicht gleichzeitig mit dem Lauf eines anderen Slice auf demselben Daemon zu laufen, richtet sich an den Aufrufer (`SPEC-058` Punkt 15; §3).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-058-container-start.md`, `docs/plan/adr/0045-container-start-als-werkzeug-im-repo.md` | liegen vor (Architect) | Vertrag und Entscheidung; zählen nicht zur Größe |
| `harness/mk/container-start.mk` | neu | Ziel `container-start` (`.PHONY`, kein Gate, Rezept ohne `$(MAKE)` und `+`, `SPEC-058` Punkt 14) |
| `tools/harness/container-start.sh` | neu | Schritte `eingaben` bis `override` (Punkte 2 bis 8), Fehlermeldung `container-start: <Schritt>: <Grund>`, Exit-Codes 2 bis 5; schreibt `mk-override.sh` nach Punkt 8; keine Zugangsvariable |
| `tools/harness/container-start-gegenprobe.sh`, `harness/mk/container-start-gegenprobe.mk` | neu | Gegenprobe: Rahmen, Attrappen (`docker` für `info` und `version`, `u-boot`), Fälle und Mutationen nach §6, Fall gegen echtes `u-boot`, `p16-gate-kette` |
| `harness/mk/gates-schnell.mk` | update | Gegenprobe in die Ausnahmeliste der schnellen Kette (`SPEC-052`, Architect prüft den Punkt der Gate-Kette) |
| `tools/harness/gates-schnell-gegenprobe.sh` | update (Zeile 32) | Die Liste `LANGSAM` ihrer Fixture spiegelt die Ausnahmeliste; ohne den Eintrag `container-start-gegenprobe` ist `gates-schnell-gegenprobe` an der unveränderten Kopie rot (`p4-eintrag`); eine Zeile, kein weiterer Liefer-Punkt |
| `harness/README.md` | update | Abschnitt *Werkzeuge* (`container-start`) und *Sensors* (Gegenprobe, Bindung an [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md)) |

- Das Skript trennt Prüfen und Handeln: erst alle Vorbedingungen am Host (bis einschließlich `zustand` ändert es nichts), dann Gerüst. Die Reihenfolge der Gründe bindet `SPEC-058` Punkt 13.
- **Docker-Last.** Die Gegenprobe dieses Slice fährt in Temp-Verzeichnissen mit Attrappen und dem echten `u-boot` ohne Container (`SPEC-058` Punkt 15, Gruppe Attrappen; Gruppe Docker mit Container liefert `slice-harness-devcontainer-einrichten-start`). Eine Vorbedingung „kein `make gates` läuft“ gibt es nicht: Die Gegenprobe läuft selbst in `make gates`, und gleichzeitige Läufe auf einem Daemon sind nicht zugesagt (`SPEC-058` Punkt 15, `SPEC-052` Punkt 8); die Regel „nie gleichzeitig mit dem Lauf eines anderen Slice“ richtet sich an den Aufrufer. Die Laufzeit der Gegenprobe ist zu messen (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`).
- **Arbeit des Implementers im Dev-Container `pgr-dev2`** (Entscheidung des Nutzers vom 2026-10-18; Stand 2026-10-18: Container `pgr-dev2` läuft, Host-Verzeichnis `/Development/pgr-dev2`, Port 2223, Klon `/workspaces/pgr-dev2/pgwire-recorder`, Remote `container` im Host-Repo, Hostschlüssel-Datei `/Development/pgr-dev2/known_hosts`, kein Claude Code im Container). Die Hostsitzung treibt den Klon per `ssh`:
  1. Jeder Befehl: `ssh -i /Development/pgr-dev2/ssh/id_ed25519 -o IdentitiesOnly=yes -o UserKnownHostsFile=/Development/pgr-dev2/known_hosts -p 2223 vscode@127.0.0.1 'cd /workspaces/pgr-dev2/pgwire-recorder && <Befehl>'`. Die Datei hat der Nutzer bereitgestellt; sie ist eine Einstellung des Implementers für seine eigenen Befehle, das Ziel dieses Slice setzt sie nicht.
  2. Zweig `slice-harness-devcontainer-einrichten-geruest` im Klon, nicht `main`; Dateien als Patch über `stdin` (`git apply`) oder als Inhalt über `stdin`; Commits im Klon mit der Kennung im Text; `make hooks-install` und `make traeger-fetch` im Klon, falls noch nicht geschehen.
  3. `make gates-schnell` und `make test` im Klon; `make gates` (28 min 21 s laut Dokument §8, für diesen Stand ungeprüft) im Hintergrund (`setsid`, Log in `/tmp` im Container, Exit-Code in eine Datei), Abfrage per `ssh`. Vor einem eigenen Lauf der Gegenprobe außerhalb von `make gates` prüft der Aufrufer selbst, dass auf demselben Daemon kein anderer Lauf läuft (Regel an den Aufrufer; die Gegenprobe prüft es nicht).
  4. *Was nur der Host kann.* Die Gegenprobe braucht `u-boot`, `node`, `ssh-keygen` (und für Teil II `devcontainer`, `docker`). Ob sie im Container von `pgr-dev2` vorliegen, ist ungemessen: Der Implementer prüft `command -v u-boot devcontainer node ssh ssh-keygen` im Container zuerst und nennt das Ergebnis im Bericht. Fehlen sie, läuft der Teil der Gegenprobe, der sie braucht, am Host mit Scratch-Verzeichnis im Präfix der Gegenprobe (nie `pgr-dev2` anfassen, nie stoppen, neu starten oder entfernen, auch nicht dessen Port 2223 belegen); der Bericht sagt, was im Container und was am Host lief.
  5. Ergebnis zum Host: im Klon `git push origin <Zweig>`; am Host `git fetch container`, der Zweig liegt als `container/<Zweig>` zur Prüfung (Dokument §9). Der Inhalt stammt aus einem Container, in dem `vscode` `root` werden kann: Review und Verifikation sehen die Commits durch, bevor sie nach `main` gehen.
  6. Größenprobe: Die Metrik von §4 misst der Implementer im Klon mit `git diff --numstat <Basis>..<letzter Commit>`.
- Zählung nach `AGENTS.md` §3.13 siehe §8.

## 4. Trigger

**Start** (`next` → `in-progress`): erfolgt (`slice-harness-devcontainer-u-boot-0-10` in `done/`; Vertrag und ADR liegen vor und sind `Accepted`; Host-Voraussetzungen entschieden, §1). Reihenfolge: dieser Slice, `slice-harness-devcontainer-einrichten-start`, `slice-harness-devcontainer-einrichten`, dann gegebenenfalls `slice-harness-devcontainer-einrichten-zugang`.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt, hier `8b19dc66` (der Commit vor dem Move); Vertrag und ADR (Commits des Architect vor dem Start) zählen nicht.

**Schätzung (Planner, 2026-10-18, nach dem zweiten Schnitt; Metrik wie oben, Gegenprobe mit dem Doppelten des Produktcodes, Faktor aus `AGENTS.md` §3.14; Grundlage: die Schätzung des Architect für das ganze Ziel von 1000 bis 1450, aufgeteilt an der Naht von `SPEC-058`; der Rahmen der Gegenprobe liegt ganz in diesem Slice):**

| Teil | Zeilen |
|---|---|
| Skript: `eingaben` und `werkzeuge` (Formen, Version, Reihenfolge, alle Verstöße) | 55 bis 80 |
| Skript: `zustand` (Verzeichnis, Schlüsselpaar), `geruest`, `schluessel`, `override` mit `mk-override.sh`, Rahmen (Schritt, Meldung, Exit) | 115 bis 170 |
| Fragmente (Ziel, Gegenprobe, Ausnahmeliste) | 30 bis 40 |
| Gegenprobe (Doppeltes von Skript und Fragmenten, 200 bis 290): Rahmen, Attrappen, rund 90 Fälle der Punkte 2 bis 8, Mutanten | 400 bis 540 |
| Index in `harness/README.md` | 10 bis 15 |
| **Zusammen** | **610 bis 845** |

Meldeschwelle 700, **Marke 900**. Die obere Hälfte liegt über der Meldeschwelle, die ganze Schätzung unter der Marke; ein weiterer Schnitt unterbleibt. Gegenüber der Schätzung des Architect für das ganze Ziel (1000 bis 1450) liegen die Teile zusammen (dieser Slice 610 bis 845 und `slice-harness-devcontainer-einrichten-start` 400 bis 620, also 1010 bis 1465) im selben Spektrum.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder der Vertrag verlangt mehr als zwei Liefer-Punkte oder mehr als zwei Schichten, oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner (Vorschlag: `eingaben` und `werkzeuge` mit Gegenprobe getrennt von `zustand`, `geruest`, `schluessel` und `override`).
- `in-progress` → `open` (blockiert): Das Gerüst aus dem Dokument trägt den Ablauf nicht ohne Eingriff am Host (Docker-Version, Netz; der Befund geht an den Planner), oder `SPEC-058` ist an der Naht nicht schneidbar (Z1 ist in `SPEC-058` entschieden; trägt der Zwischenstand nicht, geht der Befund an den Planner).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün auf dem Endstand, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben. Das Mutationsprotokoll (Zusage · Mutation · roter Test, `AGENTS.md` §3.10) steht im Bericht des Implementers. Der Fall gegen das echte `u-boot` ist einmal grün gesehen, nicht nur gegen Attrappen; die Laufzeit der Gegenprobe ist gemessen und im Bericht genannt.

## 6. Risiken und offene Punkte

Der Slice liefert einen **neuen Vertrag** (Werkzeug mit Zusagen); §3.12 gilt. Die Randformen stehen hier mit ihrer Fundstelle in `SPEC-058`; sie sind entschieden, bevor der erste Code-Commit steht (Architect). Eine Form, die erst ein Commit mit Code oder Tests hier einträgt, ist nicht genannt. Der Implementer entscheidet keine still, er hält an und gibt sie dem Architect zurück.

**Entschieden vor dem Code (Architect, Nutzer):** E1 Ort im Repo ([ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md)), E2 Host-Voraussetzungen (`SPEC-058` Punkt 3; Entscheidung des Nutzers vom 2026-10-18), E3 Randformen des Aufrufs (Punkte 2 bis 8), E4 Aufrufarten von `make` (Punkt 14), E5 Reihenfolge der Gründe und Exit-Codes (Punkt 13), E6 Stopp und Neustart als zweiter Aufruf (Punkt 9; Teil II), E7 Schätzung (§4, aufgeteilt).

**Entschieden nach dem Schnitt (Architect, 2026-10-18):** **Z1** Zwischenstand: `SPEC-058` Punkt 13 (Exit-Code 6 mit der Zeile `container-start: override: Teil II fehlt, Gerüst liegt, kein Container gestartet`, nie Exit 0, nie `bereit`; nicht zugesagt, entfällt mit Teil II) und Punkt 14 (Erfolgsfall in Teil I ist dieser Zwischenstand; Fall `p13-zwischenstand`). **Z3** Abbruch: `SPEC-058` Punkt 12 gilt erst mit Teil II; dieser Slice sagt weder Falle noch Abbruchfall zu und hat keinen (akzeptiertes Negativ mit Grund im Vertrag). **Z4** Gruppe Docker: Regeln (a) und (e), (f) für Verzeichnisse und Temp-Pfade gelten in diesem Slice, (b) bis (d) und die Regeln für Container, Image und Volumes in `slice-harness-devcontainer-einrichten-start` (`SPEC-058` Punkt 15).

**Randformen von Teil I** (Spalte *Wo entschieden*: Fundstelle in `SPEC-058`; die Randformen von Teil II stehen in `slice-harness-devcontainer-einrichten-start` §6):

| # | Randform | Wo entschieden |
|---|---|---|
| A1 | Werkzeug fehlt (`docker`, `u-boot`, `devcontainer`, `node`, `ssh`, `ssh-keygen`), je Werkzeug ein Fall; `u-boot` unter 0.10.0; alle fehlen (Reihenfolge der Meldungen) | Punkt 3 |
| A2 | `PGR_DEV` fehlt, ist leer, relativ, mit Leerzeichen, Tabulator, Zeilenumbruch, Nicht-ASCII, mit Segment `..`, `.`, leerem Segment, `/` am Ende, Basisname mit Großbuchstaben oder Unterstrich; Pfad im Repo (auch über Symlink), Elternpfad fehlt, Pfad ist kein Verzeichnis | Punkt 2 (Form), Punkt 5 (Lage im Repo, Elternpfad) |
| A3 | Verzeichnis `$PGR_DEV` besteht schon: leer, mit Gerüst, mit `u-boot.yaml` aber ohne `devcontainer.json`, fremder Inhalt, anderer Name in `u-boot.yaml` | Punkt 5 |
| A4 | Form von `PGR_SSH_PORT` (leer, nicht numerisch, führende Null, Vorzeichen, 0, 65536; Grenzen 1 und 65535); belegt prüft Teil II | Punkt 2 (Form) |
| A6 | Schlüsselpaar: fehlt, gültig, nur ein Teil, Rechte von `ssh/` und privater Datei, Passphrase, `.pub` eines anderen Schlüssels, leere Datei; nie überschreiben | Punkt 6 |
| A7 | `u-boot` meldet Warnung bei Exit 0 (Erfolg), scheitert (Exit 5, kein Folgeaufruf), Ausgabe unverändert; Fortsetzen ohne `init`; kein `up`, kein `doctor` | Punkt 7 |
| A11 | Reihenfolge der Gründe: Nachbarpaare ohne Container (Eingabe vor Werkzeug; Werkzeug vor Pfad im Repo; Pfad im Repo vor fremdem Verzeichnis; fremdes Verzeichnis vor Schlüsselmangel) und Exit-Code jedes Grundes | Punkt 13 |
| A13 | Form von `PGR_NAME` (leer, Großbuchstaben) und `PGR_WARTEN` (leer, `abc`, 0, 601; Grenzen 1 und 600; Verwendung in Teil II) | Punkt 2 |
| A14 | `PGR_CPUS`, `PGR_MEM` nicht geprüft, unverändert an `mk-override.sh` gereicht; ein ungültiger Wert zeigt sich erst bei Docker (Teil II) | Punkt 2, Punkt 8, Punkt 17 (3) |
| A15 | Projektname aus dem Basisname von `PGR_DEV` (`a-z0-9-`, erstes Zeichen `a-z0-9`); Vorgaben `pgr-dev`, `2222`, `60`, `8`, `12g` sichtbar in `override.json` | Punkt 2 (Fall `p2-vorgaben`), Punkt 8 |
| A16 | Vorhandenes `mk-override.sh` wird nicht überschrieben, sein Inhalt nicht geprüft | Punkt 8, Punkt 17 (2) |
| A17 | Kein Wert eines Zugangs (`CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY`, `GIT_TOKEN`) in `override.json`, in einer Datei unter `$PGR_DEV` und in der Ausgabe (Wert-Sentinel in `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` und `GIT_TOKEN`; das Skript entfernt die drei aus der Umgebung aller Kinder, Fall `p8-umgebung`); die Zeile mit dem Verweis `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}` bleibt Verweis | Punkt 8 (Z2, Umgebung der Kinder) |
| A18 | Unlesbare Version von `u-boot` (keine Zeichenfolge `<x>.<y>.<z>`); Vergleich Teil für Teil als Zahl (`0.9.0` unter, `0.100.0` nicht unter 0.10.0) | Punkt 3 (a) |
| A19 | Docker nicht erreichbar (`docker info` Exit ungleich 0) oder Server unter Hauptversion 24 | Punkt 3 (b) |
| A20 | Bis einschließlich `zustand` ändert das Ziel nichts: weder Verzeichnis noch Datei noch Container; von `u-boot` nur `--version`, von `docker` nur `info`, `version`, `ps`, `inspect`, `port` | Punkt 4 (Fall `p4-nichts-geaendert` je Verstoß der Schritte 2, 3, 5, 6) |
| A21 | Der Elternpfad des Verzeichnisses wird nie angelegt (kein `mkdir -p`); das Ziel schreibt nur unter `$PGR_DEV`, berührt das Host-Repo, `~/.ssh` und `known_hosts` nicht | Punkt 1, Punkt 5 |
| A22 | Plattform: Linux-Host, bash ≥ 4, GNU-Werkzeuge zugesagt; andere Hosts nicht zugesagt, nicht abgelehnt | Kopf von `SPEC-058`, Punkt 17 |

**Aufrufarten von `make`** (Randform des Treibers, `AGENTS.md` §3.12 *Aufrufarten des Treibers*; Fundstelle `SPEC-058` Punkt 14): Das Ziel ändert Zustand außerhalb des Arbeitsbaums (Verzeichnisse am Host). Je Flag eine Entscheidung im Vertrag und je Flag ein Fall `p14-*` und eine Mutation in der Gegenprobe: `-n` (nichts läuft: kein Verzeichnis, kein `u-boot`, kein `ssh-keygen`; Rezept ohne `$(MAKE)` und `+`), `-q`, `-t` (Ziel ist `.PHONY`, kein Zeitstempel), `-i` (Fehler dürfen den Ablauf nicht überlesen: kein `bereit` nach einem Fehler, die Zeile des Schrittes auf stderr), `MAKEFLAGS` aus der Umgebung mit demselben Flag, `-n` zusammen mit `-i` (gilt `-n`), `-k` und `-j` (ändern nichts). Der Erfolgsfall von `-k` und `-j` hängt an Z1.

**Weitere Risiken:**

- **Neue Host-Abhängigkeiten** `u-boot`, `devcontainer`-CLI, `node`, `ssh-keygen`, `ssh` (§1, A1). — **Ausgang (Closure):** entfallen mit Verweis auf [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md) und `SPEC-058` Punkt 3, in denen der Nutzer sie entschieden hat (2026-10-17 und 2026-10-18); streicht er eine, eingetreten → Rückführung nach `open/` (§4).
- **Gegenprobe stört einen anderen Lauf** (Docker-Daemon geteilt, `make gates` eines anderen Slice) oder den Container, in dem der Implementer arbeitet (`pgr-dev2`), oder hinterlässt Verzeichnisse oder Volumes (`BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`). Dieser Slice startet keinen Container; die Regel „nie gleichzeitig mit dem Lauf eines anderen Slice“ richtet sich an den Aufrufer (`SPEC-058` Punkt 15, `SPEC-052` Punkt 8), die Gegenprobe prüft sie nicht. — **Ausgang (Closure):** entfallen mit Beleg, dass die Gegenprobe nur Objekte ihres Präfixes `pgr-cs-` anlegt und entfernt, kein Name ohne Präfix angefasst wird und nach Ende, `SIGINT` und `SIGTERM` kein Temp-Verzeichnis bleibt, oder eingetreten → Folge-Slice mit Kennung.
- **Gegenprobe wächst und dauert** (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`). — **Ausgang (Closure):** entfallen mit gemessener Laufzeit und Mutationsprotokoll.
- **Größe über Schätzung wegen der Gegenprobe** (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Faktor 2 eingerechnet, Rückführung §4). — **Ausgang (Closure):** entfallen mit gemessener Summe unter 900, oder eingetreten → Folge-Slice mit Kennung (vom Planner bei Rückführung vergeben).
- **Der Zwischenstand ohne Teil II wird für das fertige Ziel gehalten** (Z1; das Ziel endet nach `override`). — **Ausgang (Closure):** entfallen, weil `slice-harness-devcontainer-einrichten-start` Teil II liefert und seinen Start in `done/` dieses Slice verlangt, oder eingetreten → Befund an den Planner.
- **Konflikt an `harness/mk/gates-schnell.mk` und `harness/README.md`** mit `slice-harness-devcontainer-einrichten-start` und `slice-harness-devcontainer-einrichten` (alle ändern sie). — **Ausgang (Closure):** entfallen, weil die Reihenfolge die Nehmer hinter diesen Slice setzt (Start-Trigger dort).
- **Das Dokument veraltet still gegen das Ziel** (kein Sensor; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`); bis `slice-harness-devcontainer-einrichten-start` beschreibt es den Weg von Hand. — **Ausgang (Closure):** weiter offen im selben Eintrag, oder eingetreten (Evidence dieses Slice).
- **Arbeit des Implementers im Klon von `pgr-dev2` bricht am Egress oder am Hostschlüssel** (Hosts für `make traeger-fetch` und `make gates` sind in `egress.allow` gemessen, Dokument §8; `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`). — **Ausgang (Closure):** entfallen mit Beleg im Bericht, dass die Arbeit im Klon lief, oder eingetreten (Befund an den Planner, Eintrag im Register).

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (Teil II), `BEO-REPO/nested-runtime-podman-oder-docker-offen`, `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Zählerstände liest der Implementer aus den Dateien in `evidence/` vor dem Start (kein Zähler wird hier gesetzt). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel; die verkörperten (§3.10 bis §3.14) gelten.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den zweiten Schnitt einträgt):** Liefer-Punkte 2 (Ziel Teil I, Gegenprobe Teil I). Schichten 1 (Harness-Werkzeug mit Gegenprobe und Gate-Index; `docs/developer/` liegt in `slice-harness-devcontainer-einrichten-start`; das Gerüst unter `$PGR_DEV` liegt außerhalb des Repos und ist keine Schicht, Planung und Vertrag sind keine Schicht). Schätzung in der Metrik von §4: 610 bis 845, **unter der Marke 900** (Meldeschwelle 700 in der oberen Hälfte überschritten). Dieser Slice ist Geber der Bindungen an `slice-harness-devcontainer-einrichten-start`, `slice-harness-devcontainer-einrichten` und `slice-harness-devcontainer-einrichten-zugang`; alle drei Nehmer nennen diese Kennung in §1 und zählen in ihrem §8 nach (`grep` nach der Kennung dieses Slice findet jeden von ihnen).
