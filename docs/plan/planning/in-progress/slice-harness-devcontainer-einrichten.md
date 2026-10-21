# Slice slice-harness-devcontainer-einrichten: Make-Ziel `container-einrichten` richtet im laufenden Dev-Container Bare-Repo, Einlieferung, Klon und Zweig ein (`docs/developer/devcontainer.md`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-dev-container (Entscheidung des Nutzers vom 2026-10-18; bis dahin ohne Welle). Angelegt am 2026-10-17 auf Wunsch des Nutzers vom selben Tag; am 2026-10-18 zweimal geschnitten: Der Zugang von Claude Code liegt in `slice-harness-devcontainer-einrichten-zugang` (Entscheidung des Nutzers vom 2026-10-18: nicht gebraucht, kein Login, die Hostsitzung treibt die Arbeit per `ssh` im Container), Gerüst und Start in `slice-harness-devcontainer-einrichten-geruest` (die Schätzung 995 bis 1480 lag mit dem ganzen Spektrum über der Marke 900); der Geber wurde am 2026-10-18 ein zweites Mal geschnitten, der Start (Container, Port, Bereitschaft, Abbruch) liegt in `slice-harness-devcontainer-einrichten-start`. Am 2026-10-21 ein drittes Mal geschnitten (Entscheidung des Nutzers vom selben Tag, dritter von drei Punkten): Der Dokument-Teil (`docs/developer/devcontainer.md` §6 und §9) liegt in `slice-harness-devcontainer-einrichten-doku`, direkt hinter diesem Slice in derselben Welle (`AGENTS.md` §3.11); das nimmt einen Liefer-Punkt und eine Schicht weg (§8).

**Bezug:** [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität) nur insoweit, als die Host-Voraussetzungen des Repos berührt sind (neu: `ssh`, Entscheidung des Nutzers, §1); kein Anspruch des Produkts.
Vorgänger: `slice-harness-devcontainer-austausch-repo` und `slice-harness-devcontainer-doku` (beide `done/`); der Slice automatisiert deren Ablauf. `slice-harness-devcontainer-u-boot-0-10` (`done/`; Gerüst auf `u-boot` v0.10.0, `make gates` im Container grün, 1701 s). Die unmittelbaren Vorgänger sind `slice-harness-devcontainer-einrichten-geruest` (legt das Gerüst an) und `slice-harness-devcontainer-einrichten-start` (startet den Container); dieser Slice setzt den laufenden Container voraus, §1 und §4. Folge-Slices: `slice-harness-devcontainer-einrichten-doku` (Dokument §6 und §9; der Doku-Slice direkt hinter diesem, `AGENTS.md` §3.11) und `slice-harness-devcontainer-einrichten-zugang` (Zugang von Claude Code; beide kein Teil dieses Slice, §1). Entscheidung: [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md) (Accepted 2026-10-21, Entscheidung des Nutzers), dazu [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Obergrenze 600 s je Gegenprobe, Accepted 2026-10-21). §4 rechnet in der berichtigten Größen-Metrik: [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) (Ausschluss von `docs/plan/planning` und `docs/reviews`, `--find-renames=10%`, Basis ist der Commit des Übergangs nach `in-progress/`), Status `Proposed`; `slice-harness-groessen-metrik` trägt sie nach der Annahme in `AGENTS.md` §3.14 ein, bis dahin gilt sie aus der ADR.

**Berührte Spec-Stellen:** [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) (neuer Vertrag des Werkzeugs; vom Architect am 2026-10-21 geschrieben, jede Randform aus §6 ist dort entschieden); gebunden an [`SPEC-058`](../../../../harness/vertraege/SPEC-058-container-start.md) (Punkte 2 und 11, Zusagen der Geber, §1); Host-Voraussetzungen: [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität); `SPEC-060` Punkte 2 und 3 (die Liste, an der die Gegenprobe hängt).

**Verantwortlich:** pt9912 (Zweig impl-einrichten)

**Autor:** Planner. **Datum:** 2026-10-17, geschnitten 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Make-Ziel `container-einrichten` am Host richtet für einen laufenden Dev-Container (`docs/developer/devcontainer.md` §6 und §9) in einem Aufruf ein: Bare-Repo `/workspaces/<Projektname aus u-boot.yaml>/austausch.git`, Einlieferung von `main`, Arbeitsklon `/workspaces/<Projektname>/pgwire-recorder`, Identität des Klons aus der Git-Konfiguration des Host-Repos, `make hooks-install` und `make traeger-fetch` im Klon, Zweig für den Slice. Arbeit im Dev-Container wird damit der Standardweg, und `main` am Host bleibt für den Stop-Hook frei.

**Arbeitsweise (Entscheidung des Nutzers vom 2026-10-18):** Die Hostsitzung treibt die Arbeit im Container per `ssh` und Befehl; Claude Code läuft nicht im Container, und das Ziel löst keinen Zugang auf (Abgrenzung unten).

**Bindung des Gebers `slice-harness-devcontainer-einrichten-geruest` (eingetragen vom Planner am 2026-10-18, `AGENTS.md` §3.13; am selben Tag geteilt):** Das Ziel setzt voraus, was die beiden Geber zusagen. `slice-harness-devcontainer-einrichten-geruest`: `$PGR_DEV` mit `u-boot.yaml` (der Projektname bestimmt den Pfad `/workspaces/<Projektname>/`) und die Variablen `PGR_DEV`, `PGR_SSH_PORT`, `PGR_NAME` mit den Vorgaben aus Dokument §10 (`SPEC-058` Punkt 2). `slice-harness-devcontainer-einrichten-start`: einen Container mit dem Namen `PGR_NAME`, der per `ssh` als `vscode` mit `$PGR_DEV/ssh/id_ed25519` auf `127.0.0.1:${PGR_SSH_PORT:-2222}` antwortet (`SPEC-058` Punkt 11); dass `dockerd` im Container bereit ist, sagt `slice-harness-devcontainer-einrichten-start` nicht zu (dessen Punkt 17 (4)), das Ziel prüft es selbst. Das Ziel legt kein Gerüst an, startet keinen Container und ändert nichts am Gerüst; ob `container-einrichten` das Ziel des Gebers aufruft oder nur auf einen laufenden Container prüft, entscheidet der Architect (E1). Ändert einer der Geber einen der Punkte, folgt dieser Slice in §1, §3 und §6.

**Sendung von `slice-harness-devcontainer-u-boot-0-10` („Arbeit im Container“, eingetragen vom Planner am 2026-10-17, Entscheidung des Nutzers vom selben Tag), Teil dieses Slice:** Lang laufende Slices laufen im Container, damit `main` am Host nicht über Stunden durch einen Implementer-Lauf gesperrt ist (`BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`). Dieser Slice trägt davon den Einrichtungsweg im Container (das Ziel). Die **Dokument-Hälfte** derselben Sendung — der Ablauf eines Implementer-Laufs im Container-Klon (Arbeit im Klon, `make gates` dort, Ergebnis per `git fetch container` am Host) — geht am 2026-10-21 weiter an `slice-harness-devcontainer-einrichten-doku` (Entscheidung des Nutzers vom 2026-10-21; die Adresse bleibt einen Hop länger, im Geber in `done/` wird kein Zeiger umgeschrieben). Den Weg bis zum laufenden Container trägt der Geber `slice-harness-devcontainer-einrichten-geruest`. Dieser Slice baut auf den Befunden von `slice-harness-devcontainer-u-boot-0-10` auf und misst sie nicht neu.

**Annahme (Entscheidung des Nutzers vom 2026-10-17, „ja“): ein Slice pro Container-Sitzung.** Ein Slice läuft in einer Sitzung, danach Neustart (Stopp und Neustart trägt der Geber). Je Slice ein Zweig aus `main`, das Ergebnis kommt per `git fetch container` an den Host. Das Ziel ist deshalb der **Einrichtungs-Ablauf einer Sitzung** und wiederholbar: es läuft bei jedem Neustart, nicht einmal je Container (Randformen 1 bis 3, 16). Ein Slice, der länger als rund 8 Stunden läuft, braucht einen Neustart mitten in der Arbeit; das ist kein Fehler des Werkzeugs (§6).

Ablauf des Ziels, in dieser Reihenfolge (die Schritte 4 bis 10 des Slice vor dem Schnitt):

1. Vorbedingungen prüfen, bevor etwas geändert wird, in der Reihenfolge von `SPEC-069` Punkten 2 bis 5 und 14: Eingaben (`PGR_DEV`, `PGR_SSH_PORT`, `PGR_REMOTE`, `PGR_ZWEIG`, `PGR_SSH_OPTS`; Exit 2), `ssh` und `git` am Host (Exit 3), dann der Zustand (Exit 4) als Gerüst (Projektname aus `u-boot.yaml`, Schlüssel), Host-Repo (`refs/heads/main`, Randform 12), Identität (Randform 5), Remote `PGR_REMOTE` (Randform 3), **ein** lesender `ssh`-Aufruf im Container (Randform 4), `git` und `make` im Container, Bare-Repo (Randform 1) und Klon (Randform 2).
2. Bare-Repo anlegen, falls es fehlt: `git init --bare --initial-branch=main /workspaces/<Projektname>/austausch.git` im Container, per `ssh` (Randform 1).
3. `main` einliefern: Remote `PGR_REMOTE` (`ssh://vscode@127.0.0.1:${PGR_SSH_PORT:-2222}/workspaces/<Projektname>/austausch.git`) anlegen, falls es fehlt, dann `git push <PGR_REMOTE> refs/heads/main:refs/heads/main` mit `GIT_SSH_COMMAND` in der Umgebung genau dieses Aufrufs, ohne `--force`; `core.sshCommand` des Host-Repos liest und schreibt das Ziel nie (`SPEC-069` Punkt 7, Entscheidung des Architect E1; Randform 13).
4. Klon `/workspaces/<Projektname>/pgwire-recorder` anlegen, falls er fehlt, sonst allein `git -C <Klon> fetch --prune origin` (Randform 2).
5. Name und E-Mail des Klons **lokal** aus der Git-Konfiguration des Host-Repos setzen: nie `--global`, nie ein im Dokument oder Skript fest stehender Name; fehlt am Host einer der beiden Werte, bricht das Ziel ab, bevor es etwas ändert (Schritt 1; Reihenfolge der Gründe: Randform 14); Werte mit Anführungszeichen oder Apostroph kommen unverändert an (Randformen 5 bis 7); der Transport läuft über den Standardeingang des `ssh`-Aufrufs, nie über die Kommandozeile und nie durch eine Shell, und Steuerzeichen sind abgelehnt (`SPEC-069` Punkt 9).
6. `make hooks-install` im Klon (`SPEC-069` Punkt 10).
7. Zweig für den Slice aus `origin/main` im Klon anlegen oder auf einen vorhandenen wechseln (`SPEC-069` Punkt 11, Randform 16).
8. `make traeger-fetch` im Klon als **letzter** Schritt (`SPEC-069` Punkt 12, Randform 10; Entscheidung des Architect vom 2026-10-21: Der Schritt braucht Netz; scheitert er, stehen Bare-Repo, Klon, Identität, Commit-Träger und Zweig schon, und die Wiederholung ist ein zweiter Aufruf).

**Entscheidungen des Architect vom 2026-10-21 ([`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md), [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md)):** (E1) Das Ziel ruft das Ziel des Gebers **nicht** auf; es prüft den laufenden Container mit **einem** lesenden `ssh`-Aufruf und verweist bei Misserfolg auf `make container-start` (`SPEC-069` Punkt 5). Der Projektname kommt aus `$PGR_DEV/u-boot.yaml`, nicht aus dem Wert der Variablen (Punkt 3). Eingaben sind `PGR_DEV`, `PGR_SSH_PORT`, `PGR_REMOTE` (Vorgabe `container`), `PGR_ZWEIG` (Pflicht, nie `main`) und `PGR_SSH_OPTS` (Punkt 2). Am Host schreibt das Ziel allein das Remote; `core.sshCommand` liest und schreibt es nie, der eigene Push läuft mit `GIT_SSH_COMMAND` (Punkt 7). Die Identität setzt es nur `--local` (Punkt 9). (E5) Reihenfolge der Gründe und Exit-Codes stehen in Punkt 14 (2 Eingaben, 3 Werkzeuge, 4 Zustand, 5 Fehler eines Schrittes, 130/143 Signal), die Identität vor dem Container. **Die Gegenprobe fährt ohne Container, ohne Docker und ohne Netz** (Attrappen von `ssh` und `make`, echtes `git` auf lokalen Pfaden; Punkt 17); der Lauf gegen einen echten Container ist der **Nachweis** aus Punkt 18 — kein Gate, keine Gegenprobe, kein eigener Treiber, Beleg in §7. Grund: die Obergrenze von 600 s je Gegenprobe ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)); die Rechnung steht in [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md) und in `SPEC-069` Punkt 17.

**Bindung an `slice-harness-devcontainer-einrichten-zugang` (Sendung, eingetragen vom Planner am 2026-10-18):** Das Skript liest, löst und übergibt keinen Zugang (weder `ANTHROPIC_API_KEY` noch `CLAUDE_CODE_OAUTH_TOKEN` noch eine Anmeldedatei) und prüft keinen; sein Ablauf lässt eine Stelle vor Schritt 1 frei, an der der Nehmer die Auflösung einfügt (Randform 25, Architect). Die Zeile mit dem Verweis `${localEnv:CLAUDE_CODE_OAUTH_TOKEN}` in `mk-override.sh` schreibt der Geber `slice-harness-devcontainer-einrichten-geruest`, nicht dieser Slice.

**Sendung „Quellenauflösung“ von `slice-harness-devcontainer-u-boot-0-10` (eingetragen vom Planner am 2026-10-17) geht am 2026-10-18 weiter an `slice-harness-devcontainer-einrichten-zugang`** (Entscheidung des Nutzers vom 2026-10-18). Der Geber nennt als Nehmer weiter diesen Slice; die Adresse bleibt einen Hop länger: Dieser Slice nennt den Nehmer, im Geber (`done/`) wird kein Zeiger umgeschrieben. Randform 15, Quellenauflösung, Übergabe über `devcontainer exec`, `remoteEnv` und beide Variablen stehen dort; dieser Slice führt sie nicht.

**Entscheidung des Nutzers (2026-10-17): `ssh` ist neue Voraussetzung am Host.** Auf die Frage, ob `ssh` als Voraussetzung für dieses Ziel in Ordnung ist, antwortete der Nutzer „ja“. Die Voraussetzungen des Repos sind bisher `git`, `docker` und `make` (`AGENTS.md` §3.1, [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität)). Der Architect trägt die Entscheidung in die Spezifikation des Vertrags und in die ADR ein (§3); `ssh` gilt nur für dieses Ziel, kein Gate und kein anderes Ziel braucht es. Die übrigen Voraussetzungen (`u-boot`, `devcontainer`-CLI, `node`, `ssh-keygen`) gehören dem Geber `slice-harness-devcontainer-einrichten-geruest`; der Nutzer hat sie mit [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md) am 2026-10-18 entschieden.

**Kein Gate**, wie `hooks-install`: Das Ziel steht im Abschnitt *Werkzeuge* von `harness/README.md`, nicht in `GATE_CHECKS`. Seine Gegenprobe ist ein Ziel einer der beiden Listen, kein Gate (§2).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Gerüst, kein Schlüsselpaar, kein `mk-override.sh`, kein Start, Stopp oder Neustart des Containers und keine Voraussetzungsprüfung für `u-boot`, `devcontainer`, `node`: Das tragen `slice-harness-devcontainer-einrichten-geruest` (Gerüst, Schlüsselpaar, `mk-override.sh`, Voraussetzungen) und `slice-harness-devcontainer-einrichten-start` (Start, Stopp, Neustart) (Klasse 1, Geber; deren §1 nennen diese Schritte und die Bindung an diesen Slice). Das Ziel setzt einen laufenden Container voraus (Dokument §3 bis §5) und fragt nur, ob er per `ssh` erreichbar ist.
- Kein Zugang von Claude Code im Container: keine Quellenauflösung, keine Randform 15, kein `remoteEnv`-Weg über `devcontainer exec`, keine Variable `ANTHROPIC_API_KEY` oder `CLAUDE_CODE_OAUTH_TOKEN` im Skript, kein Dokument §7 (Klasse 1: `slice-harness-devcontainer-einrichten-zugang`, dessen §1 diese Sendung annimmt; Entscheidung des Nutzers vom 2026-10-18). Das Ziel richtet den Container für Arbeit per `ssh` ein, nicht für `claude`.
- Keine Entscheidung zum Hostschlüssel (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`): Das Ziel nennt, was `ssh` ausgibt, und setzt keine Einstellung (`StrictHostKeyChecking`, `UserKnownHostsFile`) ohne Entscheidung des Nutzers (§6, Randform 9).
- Keine Messung von `make gates` im Container und keine Entscheidung zu `podman` oder `docker` als `nestedRuntime`: Das liefert `slice-harness-devcontainer-u-boot-0-10` (`done/`) bzw. bleibt offen (`BEO-REPO/nested-runtime-podman-oder-docker-offen`).
- Kein Aufruf von `docker` im Ziel **und in der Gegenprobe**: Das Ziel prüft den Container per `ssh` (`SPEC-069` Punkt 5), und die Gegenprobe fährt mit Attrappen ohne Container, ohne Docker und ohne Netz (`SPEC-069` Punkt 17, [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md); Entscheidung des Architect vom 2026-10-21 unter der Obergrenze von 600 s aus [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)). Der Lauf gegen einen echten Container ist der Nachweis aus `SPEC-069` Punkt 18: eine Handlung mit Beleg in §7, kein Gate und keine Gegenprobe.
- **Kein Dokument** `docs/developer/devcontainer.md`: §6 und §9 (das Ziel statt der Einzelbefehle, der Ablauf eines Implementer-Laufs im Container-Klon, §13 für Ungeprüftes) liegen in `slice-harness-devcontainer-einrichten-doku` (Klasse 1, `open/`; dessen §1 nennt diese Sendung mit der Kennung dieses Slice und nimmt sie an). Grund: Mit dem Dokument lag die Schätzung bei 815 bis 1210 und damit über der Marke 900 (§4); der Doku-Teil nimmt eine Schicht und 70 bis 120 Zeilen weg, und `AGENTS.md` §3.11 nennt für genau diesen Fall den Doku-Folge-Slice direkt dahinter in derselben Welle. Dieser Slice ändert am Dokument **nichts**; bis der Nehmer liefert, hinkt es dem Ziel hinterher (`AGENTS.md` §3.11, Entscheidung des Nutzers vom 2026-10-10 zu `slice-doku-ist-stand`).
- Kein Produktcode (Schicht-Abgrenzung): Es ändern sich Werkzeug und Gegenprobe, nichts unter `cmd/` und `internal/` und nichts unter `docs/user/` und `docs/developer/`.
- Kein Eintrag in der Roadmap: Der Slice ist wellenlos. Die Release-Slices und die Closure von `welle-v1-abschluss` hängen nicht an ihm (Entscheidung des Nutzers vom 2026-10-17).

## 2. Definition of Done

Liefer-Punkte (zwei; der Dokument-Punkt liegt seit dem 2026-10-21 in `slice-harness-devcontainer-einrichten-doku`, §1):

- [ ] **Ziel `container-einrichten`** (Fragment unter `harness/mk/` und Skript unter `tools/harness/`): erfüllt [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) (Ablauf aus §1, jede Randform aus §6 ist dort entschieden); kein Gate; in `harness/README.md` §Werkzeuge genannt. Dazu ist der **Nachweis** aus `SPEC-069` Punkt 18 gefahren (fünf Schritte gegen einen laufenden Container, Präfix `pgr-ce-`, Aufräumen mit Ausgabe), und sein Beleg steht in §7.
- [ ] **Gegenprobe** `container-einrichten-gegenprobe` nach `SPEC-069` Punkt 17: hängt an `GEGENPROBEN_CONTAINER` (sie ruft weder `u-boot` noch `devcontainer`, `SPEC-060` Punkte 2 und 3), nicht an `GATE_CHECKS`, und steht nicht in der Ausnahmeliste der schnellen Kette; **ohne Container, ohne Docker, ohne Netz und ohne echtes `ssh`** (Attrappen von `ssh` und `make`, echtes `git` auf lokalen Pfaden, Welten mit einem Commit in `mktemp -d`); je Randform aus §6 ein Fall und je Zusage und je Element einer Aufzählung eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); Aufräumen je Fall und je Mutation. Sie ist **ab ihrem Anlegen aufgenommen** im Sinn von `SPEC-060` Punkt 16: Sie trägt die Fälle von Punkt 16 (9), die Tabelle dort führt sie als `aufgenommen`, und der volle Lauf am Endstand liegt unter 600 s (Zeit im Bericht; Schätzung des Architect 200 bis 400 s bei `GP_PARALLEL` 4). Die Gegenprobe steht im Gate-Index.

Gates und Closure:

- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `tools/`, `harness/` und das Makefile-Fragment ändern); die Regel, nicht gleichzeitig mit dem Lauf eines anderen Slice auf demselben Daemon zu laufen, richtet sich an den Aufrufer (§3).
- [ ] `make container-einrichten-gegenprobe` am Endstand **in einem Zug** grün, mit ihrer Zeit im Bericht (`SPEC-060` Punkt 16); während der Arbeit fährt der Implementer nur den engsten Fall mit `GP_MUT`, den vollen Lauf fährt der Verifier einmal am Endstand ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md), *Übergang*). Der Lauf findet in einem frischen Gegenproben-Container **ohne vorherigen `make gates`** statt und belegt damit die Zusage von `SPEC-060` Punkt 3 für diese Gegenprobe (`SPEC-069` Punkt 17); `make gegenproben` und `make gegenproben-host` gehören der Welle und dem Release (`AGENTS.md` §6), nicht diesem Slice.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md) | **liegt vor** (Architect, 2026-10-21) | Vertrag des Werkzeugs: Ablauf, jede Randform aus §6, Host-Voraussetzungen `ssh` und `git`, Gegenprobe ohne Container, Nachweis, Plattform, Grenzen |
| [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md) | **liegt vor** (Architect, 2026-10-21; Accepted 2026-10-21) | Entscheidung und Gründe: Ort des Ziels, Host-Voraussetzungen `ssh` und `git`, Gegenprobe mit Attrappen ohne Container, Nachweis als Handlung; `Schärft:` zeigt auf `SPEC-069` (`AGENTS.md` §3.8) |
| `harness/mk/container-einrichten.mk` | neu | Ziel `container-einrichten` (kein Gate), Gegenprobe hängt an einer der beiden Listen |
| `tools/harness/container-einrichten.sh` | neu | Ablauf aus §1 (Prüfen, Einrichtung, Zweig); Eingaben aus der Umgebung (`PGR_DEV`, `PGR_SSH_PORT`, Remote-Name), Fehlermeldungen mit Schritt; keine Zugangsvariable (Bindung, §1) |
| `tools/harness/container-einrichten-gegenprobe.sh`, `harness/mk/container-einrichten-gegenprobe.mk` | neu | Gegenprobe mit Attrappen von `ssh` und `make` und echtem `git`, ohne Container, ohne Docker, ohne Netz; Fälle und Mutationen nach `SPEC-069` Punkt 17, dazu die Fälle von `SPEC-060` Punkt 16 (9); `GEGENPROBEN_CONTAINER` |
| `harness/vertraege/SPEC-060-gegenproben.md` | update | Die Tabelle in Punkt 16 nimmt die Zeile der neuen Gegenprobe als `aufgenommen` auf (sie trägt die Fälle von Punkt 16 (9)); kein anderer Punkt ändert sich |
| `harness/mk/gates-schnell.mk` | entfällt | Die Ausnahmeliste der schnellen Kette nennt keine Gegenprobe mehr (`SPEC-060` Punkt 2) |
| `harness/README.md` | update | Abschnitt *Werkzeuge* von §Sensors: beide Zeilen (Ziel und Gegenprobe, je „kein Gate"), hinter `container-start` und `container-start-gegenprobe` — der Ort, an dem dieses Repo die Gegenprobe eines Werkzeugs führt |

- Das Skript trennt Prüfen und Handeln: erst alle Vorbedingungen am Host (Identität, `ssh` vorhanden, Schlüssel lesbar, Container per `ssh` erreichbar), dann die Schritte im Container. Die Reihenfolge der Gründe legt der Architect fest (Randform 14).
- Die Form des Transports der Identität ist entschieden: Standardeingang des `ssh`-Aufrufs, je ein Wert in einer Zeile, Lesen mit `IFS= read -r`, Steuerzeichen und Nicht-UTF-8 abgelehnt (`SPEC-069` Punkte 5 (3) und 9). Der Implementer setzt sie um und entscheidet keine Quoting-Form selbst; `printf %q` aus dem Dokument §6 ist damit abgelöst.
- **Keine Docker-Last.** Die Gegenprobe braucht keinen Container und keinen Docker-Daemon (`SPEC-069` Punkt 17): Sie läuft in Temp-Verzeichnissen mit Attrappen von `ssh` und `make` und echtem `git`. Damit entfällt die Frage nach gleichzeitigen Läufen auf einem Daemon für sie; sie bleibt für den **Nachweis** (`SPEC-069` Punkt 18), der einen echten Container startet: Er läuft am Host, nie gegen `pgr-dev2`, nie gleichzeitig mit dem Lauf eines anderen Slice auf demselben Daemon (Regel an den Aufrufer, `SPEC-052` Punkt 8, `SPEC-058` Punkt 15), und räumt mit `docker rm -fv` nach Namen des Präfixes `pgr-ce-` auf, mit der Ausgabe im Beleg. Die Laufzeit der Gegenprobe ist zu messen und liegt unter 600 s (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `SPEC-060` Punkt 16).
- **Arbeit des Implementers im Dev-Container `pgr-dev2`** (Entscheidung des Nutzers vom 2026-10-18). Stand 2026-10-18: Der Container `pgr-dev2` läuft (Host-Verzeichnis `/Development/pgr-dev2`, Port `127.0.0.1:2223`), der Klon `/workspaces/pgr-dev2/pgwire-recorder` besteht, das Host-Repo hat das Remote `container` (`ssh://vscode@127.0.0.1:2223/workspaces/pgr-dev2/austausch.git`; Remote heißt `container`, nicht `container2`), und die Hostschlüssel-Datei `/Development/pgr-dev2/known_hosts` liegt vor; kein Claude Code im Container. Die Hostsitzung treibt den Klon per `ssh`:
  1. *Befehl.* `ssh -i /Development/pgr-dev2/ssh/id_ed25519 -o IdentitiesOnly=yes -o UserKnownHostsFile=/Development/pgr-dev2/known_hosts -p 2223 vscode@127.0.0.1 'cd /workspaces/pgr-dev2/pgwire-recorder && <Befehl>'`. Die Datei hat der Nutzer bereitgestellt; die Einstellung gilt für die eigenen Befehle des Implementers, das Ziel dieses Slice setzt sie nicht (Abgrenzung §1, Randform 9).
  2. *Vorbefund.* Bare-Repo, Klon und Remote bestehen für `pgr-dev2` schon (Handarbeit nach Dokument §6 und §9). Der Implementer prüft vor dem Start, dass Bare-Repo und Klon antworten (`git -C … status`, `git ls-remote container`), Identität des Klons lokal gesetzt, `core.hooksPath` gesetzt, Träger `ai-harness-init` vorhanden; was fehlt, trägt er nach dem Dokument nach und nennt es im Bericht. Nach der Lieferung des Ziels ruft er es für diese Instanz auf (Probe am eigenen Arbeitsplatz: zweiter Aufruf auf bestehendem Bare-Repo und Klon, Randformen 1 bis 3).
  3. *Arbeiten.* Zweig `impl-einrichten` im Klon, nicht `main` (derselbe Name, den `Verantwortlich:` im Kopf nennt; die Zweignamen dieses Repos tragen das Rollen-Praefix, nicht die Slice-Kennung). Dateien kommen als Patch über `stdin` (`git apply`) oder als Inhalt über `stdin` (`cat > Datei`), nicht über eine eingehängte Arbeitskopie (der Container sieht den Arbeitsbaum des Hosts nicht). Commits im Klon, jeder mit der Kennung `slice-harness-devcontainer-einrichten` im Text (Hook).
  4. *Gates.* `make gates-schnell` und `make test` im Klon; der volle Lauf `make gates` dauert 28 min 21 s (Dokument §8, für diesen Stand ungeprüft) und läuft im Hintergrund (`setsid`, Log in `/tmp` im Container, Exit-Code in eine Datei), Abfrage per `ssh`; vor einem eigenen Lauf der Gegenprobe außerhalb von `make gates` prüft der Aufrufer selbst, dass auf demselben Daemon kein anderer Lauf läuft (Regel an den Aufrufer, Docker-Last).
  5. *Was nur der Host kann.* Die Gegenprobe läuft vollständig im Dev-Container (`git`, `bash`, `make`, keine Attrappe braucht Docker oder Netz). Allein der **Nachweis** (`SPEC-069` Punkt 18) braucht `u-boot`, die `devcontainer`-CLI und Docker und läuft deshalb am Host, mit Präfix `pgr-ce-`, eigenem Port (nie 2223), eigenem Temp-Verzeichnis und `docker rm -fv` am Ende; `pgr-dev2` wird nie angefasst, gestoppt, neu gestartet oder entfernt. Der Bericht sagt, was im Container und was am Host lief.
  6. *Ergebnis zum Host.* Im Klon `git push origin impl-einrichten`; am Host `git fetch container`, der Zweig liegt als `container/impl-einrichten` zur Prüfung vor (Dokument §9). Der Inhalt stammt aus einem Container, in dem `vscode` `root` werden kann: Review und Verifikation sehen die Commits durch, bevor sie nach `main` gehen. Der Übergang `next/` → `in-progress/` ist ein Commit des Nutzers auf `main` vor der Arbeit und wird vor dem ersten Stand im Klon eingeliefert (`git push container main`).
  7. *Größenprobe.* Die Metrik von §4 misst der Implementer im Klon mit `git diff --numstat <Basis>..<letzter Commit>`.
- Zählung nach `AGENTS.md` §3.13 siehe §8.

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten-start` liegen in `done/` (Reihenfolge: beide Geber, dann dieser Slice, dann `slice-harness-devcontainer-einrichten-doku` direkt dahinter, dann gegebenenfalls `slice-harness-devcontainer-einrichten-zugang`); `slice-harness-devcontainer-u-boot-0-10` liegt in `done/` (erfüllt, Messung grün: zwei Läufe zu je 1701 s). Den Anspruch trägt der Nutzer; `Verantwortlich:` setzt der Übergang `open` → `next`, den Zweig der Übergang nach `in-progress` (Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine). Die Release-Slices und die Closure von `welle-v1-abschluss` verlangen nur `slice-harness-devcontainer-u-boot-0-10` in `done/`, diesen Slice nicht (Entscheidung des Nutzers vom 2026-10-17). Vorbedingung des Starts außerdem: Vertrag und ADR des Architect liegen vor und sind entschieden — **erfüllt** (`SPEC-069` und [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md) im Status `Accepted`, Entscheidung des Nutzers vom 2026-10-21).

**Größe.** Metrik in der berichtigten Fassung ([ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md), Status
`Proposed`; `slice-harness-groessen-metrik` trägt sie in `AGENTS.md` §3.14 ein, bis dahin gilt sie aus
der ADR). Die Pfade sind seit dem Schnitt vom 2026-10-21 `tools harness` — `docs/developer` liegt beim
Nehmer `slice-harness-devcontainer-einrichten-doku`:

```
git diff --numstat --find-renames=10% <Basis>..HEAD -- tools harness \
  ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' \
  | awk '{a+=$1; b+=$2} END {print a+b+0}'
```

Basis ist der Commit, der diese Datei nach `in-progress/` bewegt: **`28bb1f01`** (Planner, 2026-10-21; der Elternteil `96be282a` ergibt dieselbe Zahl, weil beide nur ausgeschlossene Pfade berühren);
Vertrag und ADR (Commits des Architect vor dem Start) liegen vor der Basis und zählen nicht. Für diesen
Plan sind alte und neue Fassung **gleichwertig**: Seine Pfade (`tools harness`) reichen
nicht in den Planungsbaum und nicht in `docs/reviews`, und `--find-renames=10%` wirkt nur auf bewegte
Dateien, die der Slice nicht hat. Messpunkt ist der committete Stand zu jeder Zeit der Arbeit, nicht
erst der letzte Commit.

**Schätzung nach dem Schnitt vom 2026-10-18 (Planner; Metrik wie oben, Gegenprobe mit dem Doppelten des Produktcodes, Faktor aus `AGENTS.md` §3.14):**

| Teil | Zeilen |
|---|---|
| Skript: Prüfen (Identität, `ssh`, `git`, `main`, Erreichbarkeit, Remote) | 20 bis 30 |
| Skript: Einrichtung (Bare-Repo, Einlieferung, Klon, Identität, Hook, Träger, Zweig) | 115 bis 170 |
| Fragmente | 30 bis 40 |
| Gegenprobe (Doppeltes von Skript und Fragmenten, 165 bis 240) | 330 bis 480 |
| Dokument (§6, §9 gekürzt, Löschungen zählen, dazu Ablauf des Implementer-Laufs, Index) | 70 bis 110 |
| **Zusammen** | **565 bis 830** |

Meldeschwelle 700, **Marke 900**. Lesart des Planners vom 2026-10-18 zu der Tabelle darüber: Die obere Hälfte liegt über der Meldeschwelle, die ganze Schätzung unter der Marke; ein Schnitt unterbleibt. **Diese Lesart trägt seit dem 2026-10-21 nicht mehr** — der Absatz darunter rechnet neu.

**Bestätigung und Berichtigung der Schätzung (Architect, 2026-10-21, E7; dieselbe Metrik, neue Fassung):** Die Schätzung des Planners trägt nicht mehr, weil der Vertrag seit dem 2026-10-18 zwei Teile mehr verlangt und der Faktor für Gegenprobe-Code im Bestand dieses Repos über 2 liegt. Gemessen am nächsten Verwandten: `tools/harness/container-start.sh` hat 278 Zeilen, `tools/harness/container-start-gegenprobe.sh` 886 — **Faktor 3,2**; `archiv-bezug.sh` 104 zu 477 — Faktor 4,6 (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`).

| Teil | Zeilen (Architect) |
|---|---|
| Skript: Eingaben, Werkzeuge, Projektname (`SPEC-069` Punkte 2, 3) | 45 bis 65 |
| Skript: Zustand mit einem Lesezug (Punkt 5) | 55 bis 80 |
| Skript: die sieben handelnden Schritte, Meldungen, Signale (Punkte 6 bis 16) | 75 bis 110 |
| Fragmente (Ziel, Gegenprobe, Listen-Anschluss, `help`) | 25 bis 35 |
| Gegenprobe: Rahmen, Attrappen, Welten, Mutationsmechanik | 150 bis 200 |
| Gegenprobe: rund 65 Fälle und rund 120 Mutationen | 300 bis 450 |
| Gegenprobe: Fälle des Zeitbudgets (`SPEC-060` Punkt 16 (9)) | 80 bis 120 |
| `harness/README.md` (zwei Zeilen) und `SPEC-060` Punkt 16 (eine Zeile) | 15 bis 30 |
| Dokument (§6, §9 gekürzt, Löschungen zählen, Ablauf des Implementer-Laufs, Nachweis-Hinweis) | 70 bis 120 — **seit dem 2026-10-21 beim Nehmer `slice-harness-devcontainer-einrichten-doku`** |
| **Zusammen, mit Dokument** | **815 bis 1210** |
| **Zusammen, nach dem Schnitt vom 2026-10-21 (gilt)** | **745 bis 1090** |

Mit dem Dokument lag die **ganze** Schätzung über der Meldeschwelle 700 und die obere Hälfte über der **Marke 900**. Das war eine Meldung an den Planner, keine Rückführung: Gemessen wird am committeten Stand. **Empfehlung des Architect, in dieser Reihenfolge:** (1) Den Dokument-Punkt **vorab** in einen eigenen Slice `slice-harness-devcontainer-einrichten-doku` direkt dahinter schneiden (`AGENTS.md` §3.11 und §3.13 nennen diese Form für genau diesen Fall); das nimmt 70 bis 120 Zeilen und eine Schicht weg und ergibt 745 bis 1090. (2) Reißt die Marke während der Arbeit, liegt die Naht im Vertrag: **Teil 1** (`SPEC-069` Punkte 2 bis 5: Eingaben, Werkzeuge, Zustand — nichts geändert) und **Teil 2** (Punkte 6 bis 13: die handelnden Schritte), je mit den eigenen Fällen und einem Zwischenstand nach dem Muster von `SPEC-058` Punkt 13; derselbe Schnitt trug schon die beiden Geber. (3) Die Gegenprobe allein in einen Folge-Slice zu schneiden ist **keine** Option: Ein neuer Vertrag ohne die Mutation je Zusage verstieße gegen `AGENTS.md` §3.10. Die Summe von Geber und diesem Slice (1125 bis 1695) liegt über der des ungeteilten Slice (995 bis 1480), weil Fragmente, Gegenprobe-Rahmen und Index doppelt anfallen.

**Entscheidung des Planners vom 2026-10-21 (Entscheidung des Nutzers vom selben Tag, dritter von drei Punkten):** Stufe (1) ist ausgeführt — der Dokument-Punkt liegt in `slice-harness-devcontainer-einrichten-doku`. Es gilt **745 bis 1090** bei Meldeschwelle 700 und **Marke 900**. Die **obere Hälfte liegt weiter über der Marke**, und das steht hier als Zahl statt als Beruhigung: Erreicht die Summe 700, hält der Implementer an und meldet dem Planner, bevor er weiterarbeitet; reißt sie 900, greift die Rückführung unten mit der **vorab benannten Naht** aus Stufe (2) — Teil 1 (`SPEC-069` Punkte 2 bis 5) und Teil 2 (Punkte 6 bis 13). Ein zweiter Schnitt **vor** der Arbeit unterbleibt, weil die Naht erst trägt, wenn der Zwischenstand von Teil 1 beschrieben ist (Muster `SPEC-058` Punkt 13), und die untere Hälfte der Schätzung unter der Marke liegt; ob stattdessen jetzt geschnitten wird, ist eine Entscheidung des Nutzers, die der Planner ihm am 2026-10-21 vorgelegt hat.

**Kein weiterer Schnitt** (Entscheidung des Nutzers vom 2026-10-21, wörtlich vom Koordinator übermittelt: „ok deine empfehlugnen" auf eine Vorlage mit sechs Punkten, deren dritter dieser war): Nach dem Abspalten des Dokument-Teils in `slice-harness-devcontainer-einrichten-doku` wird dieser Slice nicht erneut geteilt. Marke und Meldeschwelle dieses §4 bleiben unverändert: Erreicht die Messung die Meldeschwelle, hält der Implementer an und meldet dem Planner; reißt sie die Marke, ist das ein Befund an den Planner und **keine** Rückführung nach `next` — über einen Schnitt hat der Nutzer entschieden, und die Entscheidung bleibt bei ihm.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik oben liegt **über 900** (Meldeschwelle 700, dann anhalten und melden), oder der Vertrag verlangt mehr als zwei Liefer-Punkte oder mehr als eine Schicht (zum Beispiel eine Änderung am Gerüst des Gebers oder am Dokument), oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner an der **vorab benannten Naht**: Teil 1 (`SPEC-069` Punkte 2 bis 5, Eingaben, Werkzeuge, Zustand, mit Zwischenstand nach dem Muster von `SPEC-058` Punkt 13) und Teil 2 (Punkte 6 bis 13, die handelnden Schritte), je mit den eigenen Fällen. Die Gegenprobe allein abzutrennen ist keine Option (`AGENTS.md` §3.10).
- `in-progress` → `open` (blockiert): Der Ablauf trägt nicht ohne Eingriff in das Image (zum Beispiel fehlt `git` im Container, Rechte an `/workspaces`), oder `ssh` am Host ist nicht zugesagt (Entscheidung des Nutzers wird zurückgenommen), oder das Gerüst, das der Geber hinterlässt, trägt den Ablauf nicht (der Befund geht an den Planner).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün auf dem Endstand, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben. Das Mutationsprotokoll (Zusage · Mutation · roter Test, `AGENTS.md` §3.10) steht im Bericht des Implementers. Der **Nachweis** aus `SPEC-069` Punkt 18 ist gegen einen laufenden Container gefahren, und sein Beleg (Schritte 1 bis 5 mit ihren Ausgaben, Versionen, Exit-Codes) steht in §7; die **Gegenprobe** fährt ohne Container (`SPEC-069` Punkt 17) und ist am Endstand in einem Zug grün, ihre Zeit liegt unter 600 s und steht im Bericht (`SPEC-060` Punkt 16).

## 6. Risiken und offene Punkte

Der Slice liefert einen **neuen Vertrag** (Werkzeug mit Zusagen); §3.12 gilt. Die Randformen stehen hier, entschieden werden sie **vor dem ersten Code-Commit** in `harness/vertraege/` (Architect); eine Form, die erst ein Commit mit Code oder Tests hier einträgt, ist nicht genannt. Der Implementer entscheidet keine still, er hält an und gibt sie dem Architect zurück. Vorläufiges Wissen aus dem Dokument steht dabei, nicht als Entscheidung.

**Entscheidungen des Architect vor dem Code — am 2026-10-21 getroffen** (Vertrag [`SPEC-069`](../../../../harness/vertraege/SPEC-069-container-einrichten.md), Entscheidung und Gründe [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md), Accepted 2026-10-21; die Liste E1 bis E7 des ungeteilten Slice ist auf diesen Slice und den Geber verteilt). Vor dem Code war keine Randform offen; **sieben kamen in der Umsetzung und im Review hinzu und sind am 2026-10-21 vom Architect entschieden** (Zeilen 26 bis 32 der Tabelle unten; die Vertragsstellen stehen dort, die Belege in §7). Eine Form, die der Vertrag nicht nennt, entscheidet der Implementer nicht, er hält an und gibt sie zurück (`AGENTS.md` §3.12):

- **E1 Ort des Ziels — entschieden.** Ziel im Repo (Fragment unter `harness/mk/`, Skript unter `tools/harness/`), kein Gate, Abschnitt *Werkzeuge* (`SPEC-069` Punkt 1). Es ruft das Ziel des Gebers **nicht** auf, sondern prüft den laufenden Container mit einem lesenden `ssh`-Aufruf (Punkt 5) und verweist bei Misserfolg auf `make container-start`. Es ändert am Host **allein** das Remote; `core.sshCommand` liest und schreibt es nie (Punkt 7). Die Gegenprobe erzeugt **keinen** Scratch-Container (Punkt 17); der Container-Lauf ist der Nachweis aus Punkt 18.
- **E2 Host-Voraussetzungen — entschieden.** `ssh` (Entscheidung des Nutzers vom 2026-10-17) und `git`, je ohne Versionsprüfung (`SPEC-069` Punkte 3 und 20 (1)); im Container `git` und `make`, Vorhandensein in `zustand` geprüft (Punkt 5 (6)). `u-boot`, `devcontainer`, `node`, `ssh-keygen` braucht dieses Ziel **nicht**; sie gehören dem Geber und dem Nachweis (Punkt 18).
- **E3 Randformen des Aufrufs — entschieden** (Spalte *Wo entschieden* unten, je Zeile ein Punkt des Vertrags).
- **E4 Aufrufarten von `make` — entschieden** (`SPEC-069` Punkt 16, je Flag ein Fall und eine Mutation).
- **E5 Reihenfolge der Gründe und Exit-Codes — entschieden** (`SPEC-069` Punkt 14: 2 Eingaben, 3 Werkzeuge am Host, 4 Zustand in der Reihenfolge Gerüst, Host-Repo, Identität, Remote, Lesezug, Werkzeuge im Container, Bare-Repo, Klon; 5 Fehler eines handelnden Schrittes; 130 und 143 Signal; 0 eingerichtet, 1 nicht belegt). Die Identität steht **vor** dem Container (Randform 5 gegen 4), weil sie lokal und ohne Container behebbar ist.
- **E7 Schätzung — berichtigt** (§4, *Bestätigung und Berichtigung*): 815 bis 1210 in der neuen Metrik; die obere Hälfte liegt über der Marke 900, mit Schnitt-Empfehlung an den Planner.
- E6 (Stopp/Neustart) liegt beim Geber.

**Randformen des Ablaufs** (Spalte *Wo entschieden*: Fundstelle im Vertrag, vom Architect einzutragen). Die Zeilen 17 bis 24 des ungeteilten Slice (Werkzeug, `PGR_DEV`, Verzeichnis, Port, zweiter Aufruf, Schlüsselpaar, `u-boot`, Warten auf `ssh`) liegen als A1 bis A12 beim Geber:

| # | Randform | Vorwissen aus `docs/developer/devcontainer.md` (§6, §9, §10) | Wo entschieden |
|---|---|---|---|
| 1 | Bare-Repo besteht schon | `git init --bare` auf vorhandenem Verzeichnis meldet `Reinitialized`, Exit 0; Frage: Wiederholung ist erlaubt oder bricht ab? Besteht es mit anderem `HEAD` (nicht `main`)? | `SPEC-069` Punkte 5 (7), 6 |
| 2 | Klon besteht schon (leer oder nicht leer, anderer Remote) | `git clone` in nicht leeres Verzeichnis: Exit 128. Frage: überspringen, aktualisieren oder abbrechen; was geschieht mit uncommittetem Stand im Klon? (`pgr-dev2` hat einen Klon, Stand 2026-10-18) | `SPEC-069` Punkte 5 (8), 8 |
| 3 | Remote `container` besteht schon (gleiche URL, andere URL, anderer Port) | zweites `git remote add` bricht mit Exit 3 ab; das Host-Repo hat `container` auf Port 2223 (`pgr-dev2`). Frage: bei gleicher URL weiter, bei anderer abbrechen oder setzen | `SPEC-069` Punkte 5 (4), 7 |
| 4 | Container nicht erreichbar beim Aufruf (gestoppt, Port, Schlüssel, Zeitgrenze) | Der Start gehört dem Geber (A8). Frage: Meldung, Exit-Code, Zeitgrenze für die Verbindung (`ConnectTimeout`), vor jeder Änderung, Verweis auf das Ziel des Gebers | `SPEC-069` Punkte 5 (5), 14 |
| 5 | Identität am Host leer (nur Name, nur E-Mail, beide, nur Leerzeichen) | der Aufruf in §6 bricht ab, wenn eine fehlt. Frage: Reihenfolge gegen Randform 4 (welcher Grund zuerst), Meldung nennt den fehlenden Wert | `SPEC-069` Punkte 5 (3), 14 |
| 6 | Identität mit Sonderzeichen: Anführungszeichen `"`, Apostroph `'`, Leerzeichen, `$`, Backtick, Backslash, Zeilenumbruch, Nicht-UTF-8 | `printf %q` am Host; Gegenprobe am Wert nötig. Frage: erlaubte Menge, Quoting-Form (die Shell im Container ist die Anmelde-Shell von `vscode`), ob Zeilenumbruch abgelehnt wird | `SPEC-069` Punkte 5 (3), 9 |
| 7 | Identität nur lokal, nie global | Zusage: `--global` und `--system` werden nie gesetzt; Frage: gilt auch, wenn das Host-Repo die Werte nur aus `--global` bezieht (`git config user.name` liest alle Ebenen) | `SPEC-069` Punkt 9 |
| 8 | Anderer Port, mehrere Instanzen (zweiter Remote-Name) | Dokument §10: je Instanz eigener Port und Remote (`container2` im Dokument; für `pgr-dev2` heißt das Remote `container`); der Pfad im Container trägt den Projektnamen (`/workspaces/<Projektname>/`). Frage: wie wählt das Ziel Remote-Name und Bare-Repo-Pfad je Instanz (Parameter, Ableitung aus `u-boot.yaml`); Port aus der Umgebung des Gebers (A4); Instanz im Pfad | `SPEC-069` Punkte 2, 3, 7 |
| 9 | Hostschlüssel unbekannt, geändert (Neubau) bei den `ssh`-Aufrufen des Ziels | Warnung am Host bekannt; keine Einstellung festgelegt (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`); Neubau ändert den Fingerabdruck (A9 beim Geber). Frage: Ziel setzt nichts und bricht ab, oder die Entscheidung des Nutzers liegt vor | `SPEC-069` Punkte 5 (5), 20 (5) |
| 10 | Netz fehlt bei `make traeger-fetch` im Klon | braucht Netz im Container; das Egress bleibt eingeschränkt (Entscheidung des Nutzers vom 2026-10-17), der Host des Klons muss in der `egress.allow`-Liste stehen; `slice-harness-devcontainer-u-boot-0-10` hat für `make gates` zusätzlich `gcr.io` und `pkg-containers.githubusercontent.com` gemessen, `make traeger-fetch` im Klon ist dort nicht gelaufen (Dokument §13); Frage: Exit-Code des Ziels, ob die vorigen Schritte stehen bleiben (ja, kein Rückbau) und wie die Wiederholung aussieht | `SPEC-069` Punkt 12 |
| 11 | Abbruch mitten im Ablauf (`SIGINT`, `SIGTERM`, Verbindungsabbruch) und Wiederholung | kein Rückbau; Frage: jeder Schritt wiederholbar (Randformen 1 bis 3), Meldung nennt den letzten Schritt | `SPEC-069` Punkte 13, 15 |
| 12 | `main` am Host fehlt, ist leer, steht nicht im Host-Repo oder nur auf anderem Zweig; Host-Repo im Zustand `detached HEAD` | Frage: einzuliefernder Zweig ist `main`, das Ziel hängt nicht am ausgecheckten Zweig; Meldung bei fehlendem `main` | `SPEC-069` Punkt 5 (2) |
| 13 | Push abgewiesen (Bare-Repo hat schon einen `main` mit anderer Historie) | `! [rejected]` bei nicht-fast-forward. Frage: kein `--force`, Meldung und Exit | `SPEC-069` Punkt 7 |
| 14 | Ausgang mit Exit-Codes: Reihenfolge der Gründe, wenn mehrere zutreffen (Identität fehlt **und** `main` fehlt **und** Container nicht erreichbar) | Reihenfolge der Gründe bindet der Vertrag (`AGENTS.md` §3.12, Stelle einer Nachricht) | `SPEC-069` Punkt 14 |
| 15 | *(verschoben)* Quelle des Zugangs von Claude Code | steht in `slice-harness-devcontainer-einrichten-zugang`, Randform 15 dort; die Nummer bleibt vergeben | — |
| 16 | Zweig je Slice im Klon (Schritt 7) | Entscheidung des Nutzers vom 2026-10-17: je Slice ein Zweig aus `main`. Fragen: woher der Name (Parameter, Pflicht?), Zweig besteht schon (Wiederholung des Ziels nach Neustart: weiter auf dem Zweig oder Abbruch?), Klon hat uncommittete Änderungen, `main` im Klon steht hinter dem Bare-Repo | `SPEC-069` Punkte 2, 11 |
| 25 | Stelle des Zugangs im Ablauf (Bindung an `slice-harness-devcontainer-einrichten-zugang`) | Der Nehmer fügt die Auflösung vor Schritt 1 ein; Frage: wie das Skript die Stelle frei hält, ohne dass dieser Slice Zugangs-Code enthält (Funktion, Reihenfolge der Prüfungen) | `SPEC-069` Punkt 13 |
| 26 | `PGR_ZWEIG` ist allein nach `check-ref-format` ungültig **und** `git` fehlt am Host | In der Umsetzung aufgefallen, in §7 als Rückgabe gestellt (Review F-7): Punkt 14 ordnet „Eingabe vor Werkzeug“ (Exit 2), Punkt 2 bindet die Prüfung an `git` | `SPEC-069` Punkt 14 (**Ausnahme**: diese eine Prüfung bleibt aus, Exit 3), Punkt 2 (Tabelle), Punkt 3 (Fall `p3-fehlt-git-zweig-ungueltig`) |
| 27 | Aufräumen des Nachweises: die benannten Volumes und die Images des Gerüsts überleben `docker rm -fv` | Im Nachweis gemessen (§7): drei Kommandos des Punktes genügten nicht, die Zusage „kein Objekt des Präfixes bleibt“ war nicht erfüllt | `SPEC-069` Punkt 18 Schritt 5 (Entfernen per Name, Beleg zeigt beide Zustände) und *Akzeptiertes Negativ* |
| 28 | Eine Option in `PGR_SSH_OPTS` mit einem Shell-Metazeichen im Push | Für die eigenen `ssh`-Aufrufe sind die Optionen einzelne Argumente (Fall `p2-opts-keine-shell`); `GIT_SSH_COMMAND` ist nach Punkt 7 eine Zeichenkette, die `git` an eine Shell gibt | `SPEC-069` Punkt 20 (7) (**Grenze**, nicht zugesagt und nicht gemessen) und Punkt 2 (Einschub) |
| 29 | Der Exit-Code, mit dem `git push` eine gescheiterte Verbindung meldet | Gemessen 128, nicht 1; die Gegenprobe prüft die Form der Zeile | `SPEC-069` Punkt 20 (11) (**akzeptiertes Negativ**, Zugesagt ist Exit 5 des Ziels und die Form der Meldung) |
| 30 | Zeichenvorrat von `PGR_DEV` (Metazeichen, Leerzeichen) | Review F-3: der Pfad reist in `GIT_SSH_COMMAND` und wird dort von einer Shell ausgewertet, im Schritt, der schreibt; Punkt 2 lagerte die Form an ein anderes Ziel aus | `SPEC-069` Punkt 2 (Tabelle, Absatz *Warum `PGR_DEV` hier erneut geprüft wird*, Fälle `p2-dev-metazeichen`, `p2-dev-leerzeichen`) und Punkt 20 (7) |
| 31 | bash-Fassung am Host (4.0 bis 4.3 gegen die Zusage „≥ 4“) | Review F-4: `kind()` wartet über `$!` auf eine Prozess-Substitution — das führt bash ab 4.4 | `SPEC-069` Kopf *Plattform* (Zusage **≥ 4.4**, akzeptiertes Negativ mit Grund, kein Selbsttest) und Punkt 19 |
| 32 | `ssh` und `git` am Host gegen `AGENTS.md` §3.1 („Host braucht nur Docker und GNU `make`“) | Review F-17: die Lage besteht seit `make container-start`; der Grund stand an keiner auffindbaren Stelle | `SPEC-069` Punkt 3 (Absatz *Warum zwei Werkzeuge am Host*); `AGENTS.md` bleibt unverändert (Entscheidung des Architect) |

**Aufrufarten von `make`** (Randform des Treibers, `AGENTS.md` §3.12 *Aufrufarten des Treibers*; **entschieden in `SPEC-069` Punkt 16**): Das Ziel ändert Zustand außerhalb des Arbeitsbaums (Container, Remote, Konfiguration des Host-Repos). Je Flag eine Entscheidung im Vertrag und je Flag ein Fall und eine Mutation in der Gegenprobe: `-n` (Rezept mit `$(MAKE)` darf nichts ausführen: kein `ssh`, kein Push, keine Konfiguration), `-q`, `-t` (Ziel ist `.PHONY`, trägt kein Zeitstempel-Ziel), `-i` (Fehler dürfen den Ablauf nicht überlesen: kein „Einrichtung abgeschlossen“ nach einem Fehler), `MAKEFLAGS` aus der Umgebung, auch in Kombination und mit Vorrang. Dazu `-k` und `-j` (parallel laufende Schritte dürfen die Reihenfolge nicht brechen).

**Bindung von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`** (eingetragen vom Planner am 2026-10-18, Entscheidung des Nutzers vom selben Tag, `AGENTS.md` §3.13): Der Abschnitt über den Gegenproben-Container in `docs/developer/devcontainer.md` verweist auf §3 bis §6 und §9 dieses Dokuments und nennt deren Einzelbefehle. Diese Bindung **geht am 2026-10-21 weiter an `slice-harness-devcontainer-einrichten-doku`** (Entscheidung des Nutzers vom selben Tag), weil dieser Slice das Dokument nicht mehr anfasst (§1); dessen §1 und §6 nennen sie mit der Kennung des Gebers `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`. Im Geber wird kein Zeiger umgeschrieben — die Adresse bleibt einen Hop länger. Für diesen Slice folgt daraus nichts: kein Liefer-Punkt, keine Schicht.

**Plattform** (`AGENTS.md` §3.12 *Die Plattform gehört dazu*, **entschieden**): Geprüft wird **Linux** am Host gegen einen Linux-Container; für macOS- und Windows-Hosts ist keine Zusage geprüft, das Ziel lehnt sie nicht ab, und die Grenze steht an der Zusage im Kopf von `SPEC-069` und in dessen Punkt 19. Plattformabhängig sind die Optionsformen von `ssh`, der Transport der Identität über den Standardeingang (ein Zeilenende `CRLF` liefert einen Wert mit `CR`; deshalb sind Steuerzeichen abgelehnt) und die GNU-Optionen der gerufenen Werkzeuge. Ein Satz im Dokument über eine dieser Zusagen nennt allein Linux (`AGENTS.md` §3.11). **Nachgezogen am 2026-10-21** (Entscheidung des Architect, Review F-4): Die Zusage lautet bash **ab 4.4** — `kind()` wartet über `$!` auf eine Prozess-Substitution, und das führt bash ab 4.4; unter 4.0 bis 4.3 endet der Lauf im Lesezug des Schrittes `zustand`, also vor der ersten Änderung an Host und Container. Das Ziel prüft seine bash-Fassung nicht selbst (akzeptiertes Negativ mit Grund im Kopf von `SPEC-069`); der Kopf des Skripts nennt die Fassung.

**Weitere Risiken:**

- **`u-boot` 0.13.0 ist seit dem 2026-10-21 am Host installiert** (Angabe des Nutzers), `SPEC-058` Punkt 3 verlangt ab 0.10.0 und nennt 0.10.0 als gemessene Umgebung. Für dieses Ziel folgt daraus nichts: Es ruft `u-boot` nie auf, und die Version ist für seine Zusagen ohne Belang (`SPEC-069` Punkt 20 (3)). Zwei Berührungen bleiben: (a) Der **Nachweis** (Punkt 18) startet über `make container-start` ein Gerüst mit der Version, die am Host liegt, und nennt sie im Beleg; (b) das Egress des Gerüsts entsteht unter 0.13.0 in einer `nft -f`-Transaktion, und DNS gilt nur dem lokalen Resolver-Nutzer — das kann den Schritt `traeger` treffen (Randform 10), der ohnehin als letzter läuft und mit Exit 5 endet. `SPEC-058` ändert dieser Slice **nicht** (fremder Vertrag); die gemessene Umgebung dort nachzuziehen ist ein eigener Vorgang und liegt seit dem 2026-10-21 bei `slice-harness-container-start-u-boot-0-13` (`open/`; Angabe des Koordinators vom 2026-10-21, keine Sendung dieses Slice: der Nehmer nennt die Herkunft selbst). — **Ausgang (Closure):** entfallen mit dem Beleg des Nachweises, der die Versionen nennt, oder eingetreten → Befund an den Planner.
- **Die Gegenprobe reißt die Obergrenze von 600 s** (`SPEC-060` Punkt 16). Geschätzt 200 bis 400 s bei `GP_PARALLEL` 4; die Schätzung ruht auf gemessenen 25 bis 60 ms je Fall mit Attrappen und auf rund 65 Fällen und 120 Mutationen. Wächst eine der beiden Zahlen stark, reißt sie. — **Ausgang (Closure):** entfallen mit der gemessenen Zeit eines vollen Laufs unter 600 s, oder eingetreten → Befund an den Planner (eine Lockerung der Obergrenze ist Entscheidung des Nutzers, [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) Punkt 5).
- **Der Nachweis bleibt ungefahren oder hinterlässt einen Container** (`SPEC-069` Punkt 18, akzeptiertes Negativ; kein Gate erzwingt ihn, wie [ADR-0047](../../adr/0047-gegenprobe-des-container-starts-ist-ein-host-ziel.md) es schon für die Gegenprobe des Gebers hinnimmt). — **Ausgang (Closure):** entfallen mit dem Beleg in §7, dessen Schritt 5 zeigt, dass kein Objekt des Präfixes `pgr-ce-` bleibt; sonst eingetreten (Befund an den Planner).
- **Neue Host-Abhängigkeit `ssh`.** Entscheidung des Nutzers vom 2026-10-17 („ja“); Eintrag in Vertrag und ADR durch den Architect, kein Plan-Eintrag in `AGENTS.md` §3.1 und nicht in Lastenheft oder Gate-Index, solange der Architect nichts anderes festlegt. Fehlt `ssh` am Host, meldet das Ziel es vor jeder Änderung (eigener Fall). — **Ausgang (Closure):** offen bis zur Closure; erwartet *entfallen* mit Verweis auf Vertrag und ADR.
- **Die Geber `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten-start` liefern nicht, was §1 voraussetzt** (Container, Variablen, Pfad). — **Ausgang (Closure):** entfallen, weil der Start beide Geber in `done/` verlangt (§4); ein Befund geht an den Planner.
- **Der volle Lauf von `make gates` im Container ist von `slice-harness-devcontainer-u-boot-0-10` gemessen** (grün, zwei Läufe zu je 1701 s, höchstens 5,42 GiB; Dokument §8); dieser Slice misst ihn nicht neu. Trägt der Lauf im Klon von `pgr-dev2` nicht (anderes Gerüst, andere Last), ist das ein Befund an den Planner. — **Ausgang (Closure):** entfallen, weil `slice-harness-devcontainer-u-boot-0-10` den Lauf misst und der Start seine grüne Messung verlangt (§4).
- **Gegenprobe stört einen anderen Lauf** (Docker-Daemon geteilt, `make gates` eines anderen Slice) oder den Container, in dem der Implementer arbeitet (`pgr-dev2`), oder hinterlässt Volumes (`BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`). — **Ausgang (Closure):** entfallen mit Beleg, dass die Gegenprobe nur Objekte ihres Präfixes entfernt und nach Ende keines bleibt, oder eingetreten → Folge-Slice mit Kennung.
- **Gegenprobe wächst** (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`): zwei Zusagen je Aufzählungselement. — **Ausgang (Closure):** entfallen mit gemessener Laufzeit und Mutationsprotokoll.
- **Größe über Schätzung wegen der Gegenprobe** (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; der Architect hat am 2026-10-21 mit Faktor 3,2 nachgerechnet, §4, und die obere Hälfte liegt über der Marke 900; Rückführung und Schnitt-Naht siehe §4). — **Ausgang (Closure):** entfallen mit gemessener Summe unter 900, oder eingetreten → `slice-harness-devcontainer-einrichten-gegenprobe` (Name vorläufig, vom Planner bei Rückführung vergeben).
- **Konflikt an `harness/mk/gates-schnell.mk` und `harness/README.md`** mit dem Geber (beide ändern sie). — **Ausgang (Closure):** entfallen, weil der Start den Geber in `done/` verlangt (§4).
- **Das Dokument hinkt dem gelieferten Ziel hinterher**, bis `slice-harness-devcontainer-einrichten-doku` liefert: `docs/developer/devcontainer.md` §6 und §9 nennen nach diesem Slice weiter die Einzelbefehle, obwohl `make container-einrichten` sie ersetzt. Das ist die Form, die `AGENTS.md` §3.11 für den Doku-Folge-Slice ausdrücklich hinnimmt (Entscheidung des Nutzers vom 2026-10-10, `slice-doku-ist-stand`): Die Welle schließt erst mit beiden, und ein Release entsteht erst danach. — **Ausgang (Closure):** *entfallen*, weil `slice-harness-devcontainer-einrichten-doku` in `open/` liegt, die Sendung in §1 annimmt und in derselben Welle direkt dahinter steht (§4); liegt er bei der Closure nicht vor, ist es *eingetreten* → Befund an den Planner.
- **Die drei Dokument-Risiken liegen beim Nehmer** (das Dokument veraltet still gegen das Ziel, der Umbau verliert die Vorbehalte der Vorfassung aus §11, die Sendung „Arbeit im Container“ wächst den Dokument-Punkt): Sie sind am 2026-10-21 nach `slice-harness-devcontainer-einrichten-doku` §6 gewandert, weil dieser Slice das Dokument nicht mehr ändert. — **Ausgang (Closure):** *entfallen* mit dem Schnitt vom 2026-10-21 (die Risiken stehen im Plan des Nehmers, `grep` nach dieser Kennung findet sie dort), ohne eigenen Ausgang hier.
- **Gegenprobe gleichzeitig mit dem Lauf eines anderen Slice auf demselben Daemon** (Abgrenzung §1; die Regel richtet sich an den Aufrufer, `SPEC-052` Punkt 8, `SPEC-058` Punkt 15; die Gegenprobe prüft sie nicht). — **Ausgang (Closure):** entfallen, wenn der Bericht nennt, dass kein anderer Lauf auf dem Daemon lief; sonst eingetreten (Befund an den Planner).
- **Arbeit des Implementers im Klon von `pgr-dev2` bricht am Egress** (Randform 10; `make traeger-fetch` im Klon ist ungemessen). — **Ausgang (Closure):** entfallen mit Beleg im Bericht, dass die Arbeit im Klon lief, oder eingetreten (Befund an den Planner, Eintrag im Register).

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

### Belege des Implementers

**Arbeitsort.** Klon im Dev-Container `pgr-dev2` (ssh-Port 2223), eigener Worktree
`/workspaces/pgr-dev2/wt-einr`, Zweig `impl-einrichten` aus `origin/main` (`f1a1993b`); Dateien über den
Standardeingang des `ssh`-Befehls, Commits im Klon, `git push origin impl-einrichten` in das Bare-Repo, der
Host holt mit `git fetch container`. Allein der **Nachweis** (`SPEC-069` Punkt 18) lief am Host (Docker,
`u-boot`, `devcontainer`); `pgr-dev2` wurde dabei nie angefasst. Der Plan ist auf dem Stand von `main`;
dieser Lauf ändert §1, §3, §4 und §6 (die Korrekturen nach `AGENTS.md` §3.9, je im Commit der Lieferung) und
diesen Abschnitt.

**Läufe.**

| Lauf | Ort und Stand | Ergebnis |
|---|---|---|
| `make container-einrichten-gegenprobe`, **voller Lauf in einem Zug** | **Gegenproben-Container `pgr-gp`** (Port 2225, frischer Klon aus einem `git bundle`, **ohne** `.harness/state`, also ohne vorherigen `make gates`), Stand `4d145d66` | **grün**, 225 Fälle, **Zeit 176 s von höchstens 600 s, Last 4.22** (Last des Hosts zum Start 4.42); der Lauf am Endstand steht im Bericht des Implementers (er ist der letzte Lauf dieses Slice) |
| derselbe Lauf, `SPEC-060` Punkt 16 (8): **fremder Benutzer, fremde Locale, unter Last** | `pgr-gp`, **Endstand `a322b844`**, als `root`, `LC_ALL=en_US.utf8`, sechs Rechenlasten daneben, `GP_BUDGET=0` | **grün**, 225 Fälle, **193 s** Wanduhr (Zeile: `Zeit 193 s`, `Last 10.65`; Last zum Start 10.72); die Schlusszeile nennt den als root entfallenen Fall `p5-schluessel-unlesbar` und die mit ihm entfallene Mutation `m-p5-schluessel-lesbar` |
| voller Lauf, Stände davor | `pgr-gp` `407dae88` (176 s, Last 10.60) · `pgr-dev2` `34c1f1a0` (174 s, Last 4.12) und zweimal mit Diagnose (grün) · `pgr-dev2` `5c8a17d2` / `6c7f635d` (159 s bei Last 4.75 / 157 s bei Last 5.58, damals 216 Fälle) |
| `make gates` im Klon | `pgr-dev2`, Stand `5c8a17d2` | **grün, Exit 0, 4 min 20 s** (Start 14:07 UTC; a-check, baseline-verify, build, test, lint, docs-check, kopf-check, meldungskatalog-check, workflows-check, abdeckung-check, test-integration, test-container, test-postgres-versionen 15/16/18); der Lauf am Endstand steht im Bericht des Implementers |
| `make docs-check`, `make kopf-check` | `pgr-dev2`, Endstände der Nacharbeit | grün, 1340 Dateien, 0 Befunde (2 s); `kopf-check` Exit 0 |
| Nachweis gegen einen echten Container (Punkt 18) | Host, Stand `6c7f635d` | **durchgelaufen**, Beleg unten; zwei Läufe davor brachen vor Schritt 2 ab (unten, *Zwei Läufe, die nicht durchliefen*) |

**Zur Zusage von `SPEC-060` Punkt 3 und Punkt 16 (1)** (grün ohne vorherigen Lauf von `make gates`, Messung
im Gegenproben-Container): Nachgeholt nach Review F-9. Der Gegenproben-Container `pgr-gp` lief nicht; er ist
mit dem Werkzeug des Repos gestartet (`PGR_DEV=/Development/pgr-gp PGR_NAME=pgr-gp PGR_SSH_PORT=2225 make
container-start`, 2 min 2 s, Zeile `container-start: bereit pgr-gp 127.0.0.1:2225`, Limits des Gerüsts
unberührt: `--cpus=8 --memory=8g`). Der Zweig kam **nicht** per `git push` in sein Bare-Repo — dort hängt der
`post-receive`-Hook, der `make gegenproben` über die ganze Liste fährt —, sondern als `git bundle` über
`scp` in einen frischen Klon `/workspaces/pgr-gp/ce-pruefung`, in dem `.harness/state` fehlt: in diesem Baum
ist `make gates` nie gelaufen. Dort lief **nur** `make container-einrichten-gegenprobe`. Der eigene
Hostschlüssel-Eintrag des Containers änderte sich mit dem Neubau; die eigenen `ssh`-Aufrufe des Implementers
nutzten eine eigene `known_hosts`-Datei unter `/tmp`, `/Development/pgr-gp/known_hosts` des Nutzers blieb
unberührt (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`).

**Größe** (Metrik §4, [ADR-0056](../../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md), Basis `28bb1f01`):

```
git diff --numstat --find-renames=10% 28bb1f01..HEAD -- tools harness \
  ':(exclude)docs/plan/planning' ':(exclude)docs/reviews' | awk '{a+=$1; b+=$2} END {print a+b+0}'
```

| Datei | Zeilen (Endstand `a322b844`) |
|---|---|
| `tools/harness/container-einrichten-gegenprobe.sh` | 898 |
| `tools/harness/container-einrichten.sh` | 318 |
| `harness/vertraege/SPEC-069-container-einrichten.md` | 85 + 14 (die sieben Entscheidungen des Architect vom 2026-10-21) |
| `harness/mk/container-einrichten-gegenprobe.mk` | 12 |
| `harness/mk/container-einrichten.mk` | 11 |
| `harness/README.md` | 2 |
| `harness/vertraege/SPEC-060-gegenproben.md` | 1 |
| **Summe** | **1341** |

Gemessen am Endstand `a322b844`: **1341**; nach der ersten Runde der Nacharbeit 1326 (`48ada9b2`), vor der
Nacharbeit 1164 (`19567b8c`), davor 1162 (`6c7f635d`). Die Nacharbeit trägt 177 Zeilen: 99 Vertragstext der
sieben Entscheidungen des Architect, 63 Gegenprobe (neun Fälle, sechzehn Mutationen, Härtung der
Momentaufnahme, Diagnose), 15 Skript (Zeichenvorrat, `set -f`, Weitergabe auf stderr, Kopf). **Die Summe liegt über der Marke 900** (Meldeschwelle 700). Nach der Entscheidung des Nutzers vom
2026-10-21 (§4, *Kein weiterer Schnitt*) ist das ein **Befund an den Planner** und keine Rückführung: Der
Implementer meldet die Zahl und arbeitet weiter. Die Schätzung des Architect lautete 745 bis 1090 (§4); über
ihrer oberen Hälfte liegen 236 Zeilen (vor der Nacharbeit 74). Woran es liegt, in der Reihenfolge des Gewichts: (1) Die Gegenprobe
trägt 898 statt der geschätzten 530 bis 770 Zeilen — 225 Fälle und 194 Mutationen statt „rund 65 Fälle und
rund 120 Mutationen"; der Faktor gegen das Skript ist **2,8** (898 zu 318), also unter dem Faktor 3,2 des
nächsten Verwandten (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`). (2) Die Zahl der Fälle wächst mit
den Aufzählungen des Vertrags, nicht mit dem Umfang des Ziels: allein Punkt 1 verlangt 24 Fälle (zwei
Aufzählungen), Punkt 9 verlangt 26 Zeichen-Fälle, Punkt 2 verlangt 36 Eingabe-Fälle, Punkt 14 zehn
Nachbarpaare. (3) Das Skript liegt mit 303 Zeilen über der Schätzung von 175 bis 255; 65 davon sind
Kommentar (Punkt-Verweise und Gründe).

**Mutationsprotokoll (Zusage · Mutation · roter Fall).** Eine Zeile je Zusage oder Element einer Aufzählung;
`GP_ENTDECKE=1` nennt je Mutation die roten Fälle, und `mut` verlangt genau diese Menge (mehr und weniger
machen den Lauf rot). Die Tabelle der Mutationen in `tools/harness/container-einrichten-gegenprobe.sh`
(Abschnitt *Tabelle der Mutationen*) führt sie in 174 Zeilen (zwei davon Schleifen über 17 und 3 Elemente, zusammen **194** Mutationen) mit Text und roter Menge; hier stehen sie
nach Zusage gebündelt.

| Zusage (`SPEC-069`) | Mutation | roter Fall |
|---|---|---|
| Punkt 1: ändert nie `core.sshCommand` (schreibt und liest) | Schritt setzt sie; Schritt liest sie | `p1-sshcommand`, `p1-andere-einstellung`; `p1-sshcommand` |
| Punkt 1: ändert nie eine andere Einstellung des Host-Repos außer dem Remote | Schritt setzt `gc.auto` | `p1-andere-einstellung` |
| Punkt 1: ändert nie `~/.ssh`, nie eine `known_hosts`-Datei | Schritt legt `~/.ssh/config` an; Schritt schreibt in `known_hosts` | `p1-ssh-verzeichnis`; `p1-known-hosts`, `p1-ssh-verzeichnis` |
| Punkt 1: ändert nie eine Datei im Arbeitsbaum des Host-Repos | Schritt legt `datei.txt` an | `p1-arbeitsbaum` |
| Punkt 1: setzt nie `--global` oder `--system` am Host, nie im Container | Schritt setzt `--global` am Host; Schritt setzt `--global` im Container | `p1-global-host`, `p1-global-aufruf`; `p1-global-aufruf` |
| Punkt 1: ändert nie `$PGR_DEV` samt Gerüst und Schlüsselpaar | Schritt legt `$PGR_DEV/spur` an | `p1-pgr-dev` |
| Punkt 1: ruft nie auf — `docker`, `u-boot`, `devcontainer`, `git push --force`, `--mirror`, `--delete`, `git reset`, `git checkout -f`, `git clean`, `git merge`, `git pull`, `git rebase`, `git remote remove`, `git remote set-url`, `git branch -D`, `rm -rf` (16 Elemente) | je Element ein Schritt, der das Kommando ruft (`m-p1-nie-*`) | je Element `p1-nie-<Element>`; bei den drei Push-Formen dazu `p4-reihenfolge-der-schritte` |
| Punkt 2: `PGR_DEV` Pflicht, nicht leer, absolut, vorhanden, Verzeichnis (5 Formen) | je Form die Prüfung entfällt | `p2-dev-nicht-gesetzt`, `-leer`, `-relativ`, `-fehlt`, `-datei` |
| Punkt 2: `PGR_SSH_PORT` ganze Zahl 1 bis 65535, keine führende Null, kein Vorzeichen; Grenzen angenommen | je Prüfung entfällt; Grenze verschoben | `p2-port-leer`, `-abc`, `-plus22`, `-02222`, `-null`, `-ueber`, `p2-port-grenze-1`, `-grenze-65535` |
| Punkt 2: `PGR_REMOTE` Zeichenmenge, erstes Zeichen, Länge 64, gesetzt und leer ungültig | je Form die Prüfung geweitet (`/`, Leerzeichen, führendes `-`, 65) bzw. `:-` statt `-` | `p2-remote-slash`, `-leerzeichen`, `-minus`, `-65`, `-leer` |
| Punkt 2: `PGR_ZWEIG` Pflicht, nicht leer, nicht `main`, nicht `-…`, nicht `refs/…`, von `check-ref-format` angenommen | je Prüfung entfällt | `p2-zweig-nicht-gesetzt`, `-leer`, `-main`, `-minus`, `-refs`, `-punkte`+`-leerzeichen` |
| Punkt 2: `PGR_SSH_OPTS` an Leerraum geteilt, einzelne Argumente, Tabulator zulässig, Zeilenumbruch ungültig, von keiner Shell ausgewertet | Prüfung der Steuerzeichen entfällt; Tabulator abgelehnt; ein Argument statt mehrerer; `eval` auf den Wert | `p2-opts-zeilenumbruch`+`p2-alle-gemeldet`+`p2-reihenfolge`; `p2-opts-tabulator`; `p2-opts-leer`+`-nicht-gesetzt`+`-tabulator`; `p2-opts-keine-shell`+`p4-nichts-geaendert` |
| Punkt 2: Vorgaben `2222` und `container` | je Vorgabe eine andere | `p2-vorgaben` (+ die `p2-opts-*`-Fälle beim Port) |
| Punkt 2: alle Verstöße, je einer in einer Zeile, in der Reihenfolge der Tabelle | nach dem ersten Verstoß abbrechen; Reihenfolge der Prüfungen vertauscht | `p2-alle-gemeldet`+`p2-reihenfolge`; `p2-reihenfolge` |
| Punkt 3: `ssh` und `git` am Host, Reihenfolge der Meldungen | je Werkzeug die Prüfung entfällt; Reihenfolge vertauscht | `p3-fehlt-ssh`+`p3-fehlt-beide`+`p14-werkzeug-vor-geruest`; `p3-fehlt-git`+`p3-fehlt-beide`+`p3-fehlt-git-zweig-ungueltig`; `p3-fehlt-beide` |
| Punkt 3: Projektname aus `u-boot.yaml` (erste Zeile, vier Leerzeichen Einzug), Zeichenmenge, gleich dem Basisnamen; Datei oder Zeile fehlt heißt „kein Gerüst" | Name aus dem Basisnamen gebildet; Einzug nicht geprüft; Zeichenmenge geweitet; Vergleich entfällt; zweite Zeile gewinnt; Prüfung der Datei entfällt | `p3-name-zeile-fehlt`+`-gross`+`-unterstrich`+`-ungleich-basis`+`-anderer-einzug`; `p3-name-anderer-einzug`; `p3-name-gross`+`-unterstrich`; `p3-name-ungleich-basis`; `p3-name-zwei-zeilen`; `p3-name-datei-fehlt`+`p4-meldung-form` |
| Punkt 4: stdout trägt vom Ziel nur die Schlussmeldung | ein `echo` in einem Schritt | `p4-stdout-nur-erfolg` |
| Punkt 4: die Ausgabe eines gerufenen Werkzeugs erscheint unverändert auf ihrem Strom, je Strom | stdout umgeleitet; stderr umgeleitet | `p4-werkzeug-stdout`; `p4-werkzeug-stderr` |
| Punkt 4: jede Fehlermeldung nennt den Schritt | Meldung ohne Schritt | `p4-meldung-form`+`p3-fehlt-ssh`+`p3-fehlt-git`+`p3-fehlt-beide` |
| Punkt 4: Reihenfolge der zehn Schritte | `traeger` vor `zweig` | `p4-reihenfolge-der-schritte`+`p12-stelle`+`p12-fehler` |
| Punkt 5 (1): Gerüst — Schlüssel vorhanden, Schlüssel lesbar (zwei Elemente) | je Prüfung entfällt | `p5-schluessel-fehlt`; `p5-schluessel-unlesbar` (entfällt als root) |
| Punkt 5 (2): Host-Repo, `refs/heads/main`; detached HEAD und anderer Zweig sind kein Grund | Prüfung des Repos entfällt; Prüfung von `main` entfällt; `HEAD` statt `refs/heads/main`; Prüfung auf den ausgecheckten Zweig ergänzt | `p5-kein-repo`; `p5-kein-main`+`p14-repo-vor-identitaet`; `p5-kein-main`+`p14-repo-vor-identitaet`; `p5-detached`+`p5-anderer-zweig` |
| Punkt 5 (3): Identität — Name fehlt, E-Mail fehlt, beide in einer Zeile, leer, nur Leerraum, Steuerzeichen (Tabulator, Zeilenumbruch, CR), kein UTF-8, Meldung ohne den Wert | je Prüfung entfällt bzw. wird geweitet; nur der erste Mangel gemeldet; Leerraum nicht getrimmt; Wert in die Meldung gehängt | `p5-ident-name-fehlt`…`-leerraum`; `p5-ident-mail-fehlt`; `p5-ident-beide-fehlen`; `p5-ident-tabulator`+`-zeilenumbruch`+`-cr`; `p5-ident-nicht-utf8`+`-mail-nicht-utf8`; die vier Steuerzeichen-Fälle + `p4-meldung-form` |
| Punkt 5 (4): Remote fehlt oder hat genau die URL; nichts wird gesetzt | Vergleich der URL entfällt; `set-url` statt `get-url` | `p5-remote-andere-url`+`p5-remote-nichts-gesetzt`+`p14-remote-vor-lesezug`; dazu `p5-remote-gleich`+`p1-nie-remote-set-url`+`p4-nichts-geaendert` |
| Punkt 5 (5): **ein** `ssh`-Aufruf, Form des Aufrufs, `BatchMode`, `ConnectTimeout`, keine Hostschlüssel-Option, Lesebefehl schreibt nichts, Meldung nennt Adresse, Port und die letzte Ausgabe | zweiter Lesezug; `BatchMode` entfernt; `ConnectTimeout` entfernt; `StrictHostKeyChecking` ergänzt; Lesebefehl schreibt eine Datei; Prüfung des Exit-Codes entfällt; Meldung ohne die Ausgabe von `ssh` | `p5-lesezug-ein-aufruf`+`p4-nichts-geaendert`+`p4-reihenfolge-der-schritte`; `p5-lesezug-aufruf`; `p5-lesezug-aufruf`; `p5-lesezug-kein-hostschluessel`+`p5-lesezug-aufruf`; `p5-lesezug-schreibt-nichts`+`p4-nichts-geaendert`; `p5-nicht-erreichbar`+`p14-lesezug-vor-werkzeugen`; `p5-nicht-erreichbar` |
| Punkt 5 (6): `git` und `make` im Container, je Werkzeug | je Werkzeug die Prüfung entfällt | `p5-fern-git-fehlt`; `p5-fern-make-fehlt`+`p14-werkzeuge-vor-bare` |
| Punkt 5 (7) und (8): Bare-Repo mit anderem `HEAD`, kein Bare-Repo; Klon schmutzig, anderer `origin`, kein Git-Repo | je Prüfung entfällt | `p5-bare-master`+`p5-lesezug-schreibt-nichts`+`p14-bare-vor-klon`; `p5-bare-fremd`; `p5-klon-schmutzig`+`-unversioniert`+`p14-klon-vor-schritt`; `p5-klon-anderer-origin`; `p5-klon-fremd` |
| Punkt 6: `git init --bare --initial-branch=main` nur bei fehlendem Bare-Repo, Fehler endet mit 5 | `--initial-branch` entfällt; `init` läuft immer; Fehler überlesen | `p6-fehlt`; `p6-vorhanden`+`p13-zweiter-aufruf`; `p6-fehler` |
| Punkt 7: Remote nur anlegen, wenn es fehlt; Refspec wörtlich; kein `--force`; `GIT_SSH_COMMAND` in der Umgebung und nicht in der Konfiguration; Abweisung erkannt; kein zweiter Versuch; Fehler endet mit 5 | Bedingung invertiert; `--force` ergänzt; Refspec `main`; `core.sshCommand` gesetzt; Abweisung nie erkannt; Abweisung immer; zweiter Versuch; Fehler überlesen | `p7-remote-fehlt` und die übrigen `p7-*`+`p13-*`; `p7-push-aufruf`+`p7-abgewiesen`+`p1-nie-push-force`; `p7-push-aufruf`+`p7-fehler`; `p7-git-ssh-command`+`-opts`+`p1-sshcommand`+`p1-andere-einstellung`; `p7-abgewiesen`; `p7-fehler`; `p7-abgewiesen`; `p7-fehler` |
| Punkt 8: `clone` nur bei fehlendem Klon, sonst allein `fetch --prune origin`, Fehler endet mit 5 | `clone` läuft immer; `--prune` entfällt; `pull` statt `fetch`; Fehler überlesen | `p8-besteht`+`p8-fehler`+`p13-zweiter-aufruf`; `p8-besteht`+`p8-fehler`; `p8-besteht`+`p8-fehler`; `p8-fehler` |
| Punkt 9: Identität nur `--local`; Wert nie als Argument, nie in einer Shell-Zeichenkette; `IFS= read -r`; Reihenfolge der zwei Zeilen; 26 Zeichen kommen byteweise gleich an | `--global`; `--system`; Wert aus dem Skript; Wert als Argument; `eval` auf den Wert; `read` ohne `-r`; `read` ohne `IFS=`; Werte vertauscht | `p9-lokal`+`p4-reihenfolge-der-schritte`+alle 26 Zeichen-Fälle (je Mutation); `p9-kein-argument`+`p9-lokal`+13 Zeichen-Fälle; 13 Zeichen-Fälle; `p9-zeichen-backslash`; `p9-zeichen-leerraum-am-rand`; `p9-lokal`+`p9-reihenfolge`+26 Zeichen-Fälle |
| Punkt 10: `make hooks-install` im Klon, bei jedem Aufruf; Fehler endet mit 5 | anderes Ziel; anderes Verzeichnis; Fehler überlesen | `p10-aufruf`+`p10-zweiter-aufruf`+`p13-zweiter-aufruf`+`p10-fehler`; `p10-aufruf`+`p10-zweiter-aufruf`; `p10-fehler` |
| Punkt 11: `switch -c <Zweig> origin/main`, sonst `switch <Zweig>`; kein `--force`, kein `checkout -f`; Quelle immer `origin/main`; der Name reist über den Standardeingang | Quelle `main`; `-c` auch beim vorhandenen Zweig; `--force` ergänzt; `checkout -f`; Name als Argument; Fehler überlesen | `p11-aus-origin-main`; `p11-besteht`; `p11-besteht`; `p11-besteht`; `p11-sonderzeichen`; `p11-fehler` |
| Punkt 12: `make traeger-fetch` als letzter Schritt, Fehler endet mit 5 und ohne Schlussmeldung | anderes Ziel; Fehler überlesen; Schritt vor `zweig` | `p12-aufruf`+`p12-stelle`+`p12-fehler`; `p12-fehler`; `p4-reihenfolge-der-schritte`+`p12-stelle`+`p12-fehler` |
| Punkt 13: Schlussmeldung mit Klon-Pfad und Zweig, Exit 0 beim ersten wie beim zweiten Aufruf | Meldung ohne Pfad und Zweig | `p13-erfolg`+`p13-zweiter-aufruf`+`p13-umgebung` |
| Punkt 13: die drei Zugangs-Variablen verlassen die eigene Umgebung, je Variable | je Variable das `unset` entfällt | `p13-umgebung` |
| Punkt 14: Reihenfolge der Gründe, je Nachbarpaar | je Paar die beiden Prüfungen vertauscht bzw. eine entfernt | `p14-eingabe-vor-werkzeug` … `p14-klon-vor-schritt` (zehn Fälle, je Paar genannt) |
| Punkt 15: `SIGINT` und `SIGTERM` beenden Kind und Lauf, Meldung mit Schritt, 130 und 143, kein Rückbau | Falle entfällt (je Signal); Exit 1 statt 130/143; Kind wird nicht beendet; ohne Meldung; Rückbau ergänzt | `p15-sigint`+`-sigint-zweig`+`-kein-rueckbau`+`-wiederholung`; `p15-sigterm`+`-sigterm-traeger`; dieselben je Signal; alle sechs `p15-*`; alle sechs; `p15-kein-rueckbau` |
| Punkt 16: Rezept ohne `$(MAKE)`, Ziel `.PHONY` ohne Voraussetzungen | Rezept mit `$(MAKE)`; `.PHONY` umbenannt; Voraussetzung, die fehlt | `p16-n`+`-q`+`-t`+`-makeflags-*`+`-n-und-i`; `p16-phony`+`p17-phony`; `p16-phony`+`p16-i`+`p16-k`+`p16-j` |
| Punkt 17: Gegenprobe in `GEGENPROBEN_CONTAINER`, nicht in `GEGENPROBEN_HOST`, nicht in `GATE_CHECKS`, nicht in der Ausnahmeliste, `.PHONY` | Liste entfernt; in die Host-Liste; an `GATE_CHECKS` (Ziel und Gegenprobe); in die Ausnahmeliste (in `gates-schnell.mk`, weil ein `+=` eines früher eingebundenen Fragments vom `:=` dort verworfen wird) | `p17-liste-container`; `p17-liste-container`+`p17-nicht-host`; `p17-kein-gate`; `p17-keine-ausnahme` |
| `SPEC-060` Punkt 16 (9): Obergrenze, Zeile mit Zeit und Last, `GP_BUDGET`, Reihenfolge der Gründe, Aufräumen je Fall und je Mutation | Prüfung der Obergrenze entfällt; `0` prüft doch; Zeile entfällt; Vorgabe auch mit `GP_MUT`; Last als `0`; Lesen der Last invertiert; `GP_BUDGET` lax geprüft; Prüfung erst nach den Fällen; Falle für `INT` entfällt; EXIT-Falle räumt nicht; `rm` je Mutation entfällt; `rm` je Fall entfällt | `p17-budget-ueber`+`-gruende`; `p17-budget-null`; `p17-budget-unter`+`-null`+`-last-fehlt`+`-gruende`; `p17-budget-vorgabe-nur-ohne-mut`; `p17-budget-unter`+`-null`+`-vorgabe-nur-ohne-mut`; dieselben drei; `p17-budget-ungueltig-leer`+`-minus`+`-bruch`; die vier `-ungueltig-*`; `p17-budget-sigint`; `p17-budget-sigint`+`-sigterm`; `p17-ordner-mutation`; `p17-ordner-fall` |

**Drei Mutanten ohne roten Fall — äquivalent, je mit Invariante, Ort und Sonde** (`AGENTS.md` §3.10):

1. `m-p2-erster-verstoss` — `(( bad == 0 )) || exit 2` wird zu `(( bad == 0 )) || { (( bad > 1 )) && exit 2; exit 2; }`.
   *Invariante:* Beide Zweige enden mit `exit 2`, und nichts dazwischen wirkt nach außen.
   *Ort:* der Text des Mutanten selbst (beide Pfade enden in `exit 2`; `(( bad > 1 ))` hat keine Wirkung außer
   dem Rückgabewert). *Sonde:* Der volle Lauf fährt alle 36 Eingabe-Fälle; jede Abweichung im Exit-Code oder
   in der Zahl der Zeilen wäre an `p2-alle-gemeldet`, `p2-reihenfolge` und jedem `p2-*`-Fall sichtbar — alle
   bleiben grün, auch der Fall mit **fünf** Verstößen.
2. `m-p16-voraussetzung` — `container-einrichten: ##` wird zu `container-einrichten: Makefile ##`.
   *Invariante:* Eine Voraussetzung, die eine vorhandene Datei ohne eigene Regel ist, ist immer aktuell und
   ändert weder Rezept noch Reihenfolge noch Exit-Code des Ziels. *Ort:* `make` selbst (ein vorhandenes Ziel
   ohne Regel gilt als aktuell). *Sonde:* dieselbe Stelle mit einer Voraussetzung, die **nicht** existiert
   (`m-p16-voraussetzung-fehlt`, `container-einrichten: gibtsnicht ##`): gemessen rot in `p16-phony`,
   `p16-i`, `p16-k` und `p16-j` — die Fälle trennen Original und Mutant also, sobald die Invariante
   („vorhanden, ohne Regel") nicht gilt.
3. `m-p17-signal-term` — `trap 'exit 143' TERM` der Gegenprobe wird zu `:`.
   *Invariante:* Ein nicht gefangenes `SIGTERM` beendet die `bash` dieses Skripts mit Exit 143, führt dabei
   die EXIT-Falle aus (der Arbeitsordner entfällt) und erreicht `abschluss` nicht (keine Zeile `Zeit …`).
   *Ort:* `bash` selbst (Vorgabe-Behandlung von `SIGTERM`), nicht das Skript — für `SIGINT` gilt das
   **nicht**: die Vorgabe beendet dort nur das wartende Kind, das Skript läuft weiter.
   *Sonde* (Lauf im Dev-Container, je eine Kopie der Gegenprobe mit einer Warteschleife, Signal an die
   Kopie):

   | Kopie | Signal | Exit | Zeile `Zeit …` | Arbeitsordner |
   |---|---|---|---|---|
   | unverändert | `SIGINT` | 130 | nein | weg |
   | unverändert | `SIGTERM` | 143 | nein | weg |
   | Falle für `INT` entfällt | `SIGINT` | **0** | **ja** | weg |
   | Falle für `TERM` entfällt | `SIGTERM` | 143 | nein | weg |
   | Falle für `INT` entfällt | `SIGTERM` | 143 | nein | weg |
   | Falle für `TERM` entfällt | `SIGINT` | 130 | nein | weg |

   Die dritte Zeile ist der rote Fall von `m-p17-signal-int` (`p17-budget-sigint`); die vierte zeigt, dass
   die Falle für `TERM` am Ergebnis nichts ändert. Entfällt die EXIT-Falle, bleibt der Ordner stehen — das
   fängt `m-p17-ordner-ende` (`p17-budget-sigint`, `p17-budget-sigterm`).

**Nachweis gegen einen laufenden Container** (`SPEC-069` Punkt 18; Form wie
[`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 (6); am Host, Präfix `pgr-ce-`,
Port 2253, nie gegen `pgr-dev2`, kein anderer Lauf auf demselben Daemon):

```text
### Umgebung
Datum (UTC): 2026-10-21 13:57
uname -s: Linux
uname -m: x86_64
docker: 29.9.2
u-boot: u-boot version 0.13.0
devcontainer: 0.80.2
git: git version 2.43.0
Port: 2253, PGR_DEV: /tmp/pgr-ce-nachweis.XIg9Fw/pgr-ce-nachweis

### Vorbereitung (Klon des Host-Repos, Zweig des Slice)
git clone: Exit 0
git fetch: Exit 0
git checkout: Exit 0
Identitaet des Klons: pt9912 / 66826472+pt9912@users.noreply.github.com
HEAD: 6c7f635d slice-harness-devcontainer-einrichten: Gegenprobe container-einrichten-gegenprobe ohne Container
refs/heads/main: 606a60f7

### Schritt 1: make container-start
Initialized u-boot project "pgr-ce-nachweis".

Created:
[gekürzt: 152 Zeilen Ausgabe von `u-boot generate devcontainer` und `devcontainer up` (Bau des Images) — Ausgabe des Gebers `SPEC-058`, keine Aussage dieses Ziels; kein Geheimnis darin]
[2026-10-21T13:57:54.585Z] Start: Run: docker run --sig-proxy=false -a STDOUT -a STDERR -p 127.0.0.1:2253:2222 --mount source=pgr-ce-nachweis-workspace-0vaqmrlv9mu7d9tju7037lv8pkmodgioac84vr9vicc55af1b4te,target=/workspaces/pgr-ce-nachweis,type=volume --mount source=pgr-ce-nachweis-docker-0vaqmrlv9mu7d9tju7037lv8pkmodgioac84vr9vicc55af1b4te,target=/var/lib/docker,type=volume --mount type=bind,source=/tmp/pgr-ce-nachweis.XIg9Fw/pgr-ce-nachweis/ssh/id_ed25519.pub,target=/home/vscode/.ssh/authorized_keys,readonly -l devcontainer.local_folder=/tmp/pgr-ce-nachweis.XIg9Fw/pgr-ce-nachweis -l devcontainer.config_file=/tmp/pgr-ce-nachweis.XIg9Fw/pgr-ce-nachweis/.devcontainer/devcontainer.json --privileged --cap-add=NET_ADMIN --name=pgr-ce-nachweis --cpus=2 --memory=2g --entrypoint /bin/sh vsc-pgr-ce-nachweis-68e8a303c582202f497fa7524d9d16c2b55e23c52cfb277bf9cf47c07a51052b-uid -c echo Container started
[2026-10-21T13:57:54.978Z] Container started
[2026-10-21T13:57:55.164Z] Running the postStartCommand from devcontainer.json...

{"outcome":"success","containerId":"7f94f3bacb864ce4046c0f1226b9a74e726907718f9da7b437f4d5b415410d96","remoteUser":"vscode","remoteWorkspaceFolder":"/workspaces/pgr-ce-nachweis"}
container-start: bereit pgr-ce-nachweis 127.0.0.1:2253
Exit: 0


### Schritt 2: erster Aufruf von make container-einrichten
PGR_SSH_OPTS: -o StrictHostKeyChecking=accept-new -o UserKnownHostsFile=/tmp/pgr-ce-nachweis.XIg9Fw/known_hosts   (Punkt 20 (5): das Ziel setzt keine Einstellung fuer den Hostschluessel)
Initialized empty Git repository in /workspaces/pgr-ce-nachweis/austausch.git/
To ssh://127.0.0.1:2253/workspaces/pgr-ce-nachweis/austausch.git
 * [new branch]        main -> main
Cloning into '/workspaces/pgr-ce-nachweis/pgwire-recorder'...
done.
hooks-install: core.hooksPath=.githooks — .githooks/commit-msg laeuft ab dem naechsten Commit (Umgehung: git commit --no-verify).
branch 'impl-einrichten' set up to track 'origin/main'.
Switched to a new branch 'impl-einrichten'
Unable to find image 'curlimages/curl@sha256:463eaf6072688fe96ac64fa623fe73e1dbe25d8ad6c34404a669ad3ce1f104b6' locally
docker.io/curlimages/curl@sha256:463eaf6072688fe96ac64fa623fe73e1dbe25d8ad6c34404a669ad3ce1f104b6: Pulling from curlimages/curl
bd9ddc54bea9: Pulling fs layer
9824c27679d3: Pulling fs layer
e4200f789582: Pulling fs layer
bd9ddc54bea9: Download complete
9824c27679d3: Download complete
9824c27679d3: Pull complete
e4200f789582: Download complete
bd9ddc54bea9: Pull complete
e4200f789582: Pull complete
Digest: sha256:463eaf6072688fe96ac64fa623fe73e1dbe25d8ad6c34404a669ad3ce1f104b6
Status: Downloaded newer image for curlimages/curl@sha256:463eaf6072688fe96ac64fa623fe73e1dbe25d8ad6c34404a669ad3ce1f104b6
traeger-fetch: Traeger abgelegt (/workspaces/pgr-ce-nachweis/pgwire-recorder/.harness/state/bin/ai-harness-init) — ai-harness-init-linux-amd64 aus Release v0.2.6, Digest verifiziert.
container-einrichten: eingerichtet /workspaces/pgr-ce-nachweis/pgwire-recorder impl-einrichten
Exit: 0

### Schritt 3: zweiter Aufruf, derselbe Aufruf
Everything up-to-date
hooks-install: core.hooksPath=.githooks — .githooks/commit-msg laeuft ab dem naechsten Commit (Umgehung: git commit --no-verify).
Your branch is up to date with 'origin/main'.
Already on 'impl-einrichten'
traeger-fetch: Traeger abgelegt (/workspaces/pgr-ce-nachweis/pgwire-recorder/.harness/state/bin/ai-harness-init) — ai-harness-init-linux-amd64 aus Release v0.2.6, Digest verifiziert.
container-einrichten: eingerichtet /workspaces/pgr-ce-nachweis/pgwire-recorder impl-einrichten
Exit: 0

### Schritt 4: Ablesen im Container
$ git -C /workspaces/pgr-ce-nachweis/austausch.git symbolic-ref HEAD
refs/heads/main
  Exit: 0
$ git -C /workspaces/pgr-ce-nachweis/pgwire-recorder remote get-url origin
/workspaces/pgr-ce-nachweis/austausch.git
  Exit: 0
$ git -C /workspaces/pgr-ce-nachweis/pgwire-recorder rev-parse --abbrev-ref HEAD
impl-einrichten
  Exit: 0
$ git -C /workspaces/pgr-ce-nachweis/pgwire-recorder config --local --get user.name
pt9912
  Exit: 0
$ git -C /workspaces/pgr-ce-nachweis/pgwire-recorder config --local --get user.email
66826472+pt9912@users.noreply.github.com
  Exit: 0
$ git -C /workspaces/pgr-ce-nachweis/pgwire-recorder config --get core.hooksPath
.githooks
  Exit: 0
$ ls /workspaces/pgr-ce-nachweis/pgwire-recorder/.harness/ 2> /dev/null || echo '(kein .harness)'
baseline
skills
state
  Exit: 0
$ git -C /workspaces/pgr-ce-nachweis/pgwire-recorder status --porcelain | head -3
  Exit: 0
$ git -C /workspaces/pgr-ce-nachweis/pgwire-recorder log --oneline -1
606a60f7 slice-harness-container-start-u-boot-0-13: Slice in open/ angelegt (gemessene Umgebung von SPEC-058 auf u-boot 0.13.0, Lauf der Gruppe Docker am Host; Konflikt mit ADR-0055 in 1 und 6); Sendung in slice-harness-gegenproben-grenze-rest 1 und 6 mit Nachzaehlung in dessen 8; Drift-Log der Roadmap
  Exit: 0
Identitaet des Host-Repos (Klon): pt9912 / 66826472+pt9912@users.noreply.github.com
Remote im Klon des Hosts: ssh://vscode@127.0.0.1:2253/workspaces/pgr-ce-nachweis/austausch.git

### Schritt 5: Aufraeumen
$ docker rm -fv pgr-ce-nachweis
pgr-ce-nachweis
  Exit: 0
$ docker ps -a --filter name=pgr-ce-
CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
$ docker volume ls | grep pgr-ce-  (benannte Volumes ueberleben docker rm -fv)
local     pgr-ce-nachweis-docker-0vaqmrlv9mu7d9tju7037lv8pkmodgioac84vr9vicc55af1b4te
local     pgr-ce-nachweis-workspace-0vaqmrlv9mu7d9tju7037lv8pkmodgioac84vr9vicc55af1b4te
$ docker volume ls -q | grep pgr-ce- | xargs -r docker volume rm
pgr-ce-nachweis-docker-0vaqmrlv9mu7d9tju7037lv8pkmodgioac84vr9vicc55af1b4te
pgr-ce-nachweis-workspace-0vaqmrlv9mu7d9tju7037lv8pkmodgioac84vr9vicc55af1b4te
$ docker images | grep vsc-pgr-ce-
vsc-pgr-ce-nachweis-68e8a303c582202f497fa7524d9d16c2b55e23c52cfb277bf9cf47c07a51052b-uid:latest
vsc-pgr-ce-nachweis-68e8a303c582202f497fa7524d9d16c2b55e23c52cfb277bf9cf47c07a51052b:latest
$ docker images ... | xargs -r docker rmi
  Exit: 0
$ docker ps -a --filter name=pgr-ce- ; docker volume ls -q | grep pgr-ce- ; docker images | grep pgr-ce-   (nach dem Entfernen)
(kein Volume des Praefixes)
(kein Image des Praefixes)
git remote remove: Exit 0
rm -rf /tmp/pgr-ce-nachweis.XIg9Fw: Exit 0
### Ende
```

**Zwei Läufe, die nicht durchliefen — und was sie belegen** (dieselben fünf Schritte, Ports 2251 und 2252):

1. Der Klon des Host-Repos trug keine Identität (`git clone` kopiert die lokale Konfiguration nicht), und der
   erste Aufruf endete mit der Meldung `container-einrichten: zustand: die Identität des Host-Repos fehlt:
   user.name user.email` (über `make` Exit 2, des Skripts 4) — Randform 5 am echten Container.
2. Mit Identität, aber ohne eigene `known_hosts`-Datei endete der Aufruf mit
   `container-einrichten: zustand: der Container antwortet nicht auf vscode@127.0.0.1, Port 2252 (ssh: Host
   key verification failed.) — make container-start startet ihn` — Randform 9 und Punkt 20 (5) am echten
   Container: Das Ziel setzt keine Einstellung für den Hostschlüssel. Der dritte Lauf gibt sie in
   `PGR_SSH_OPTS` (`-o StrictHostKeyChecking=accept-new -o UserKnownHostsFile=…`), wie Punkt 20 (5) es
   vorsieht.

**Vier Beobachtungen zum Vertrag — Hinweise an den Architect** (alle vier sind am 2026-10-21 entschieden
und im Vertrag eingearbeitet, siehe *Nacharbeit nach dem Review* unten; der Text darunter ist der Stand der
Lieferung, der zur Entscheidung führte), **keine Randform-Rückgabe** (jede Zeile ist eine
Lücke im *Text*, nicht im entschiedenen Verhalten; der Implementer hat nichts davon selbst entschieden):

1. **Punkt 18, Schritt 5 räumt nicht vollständig.** `docker rm -fv` entfernt die **anonymen** Volumes; das
   Gerüst aus `SPEC-058` legt zwei **benannte** an (`pgr-ce-nachweis-workspace-…`,
   `pgr-ce-nachweis-docker-…`) und ein Image-Paar (`vsc-pgr-ce-nachweis-…`, `…-uid`). Nach den drei
   Kommandos des Punktes bleiben sie stehen, und die Zusage „kein Objekt des Präfixes bleibt" ist nicht
   erfüllt. Der Beleg oben entfernt sie **per Name** (`docker volume rm`, `docker rmi`) und zeigt beide
   Zustände; `container-start-gegenprobe` hält es mit ihrem Präfix `pgr-cs-` ebenso (`SPEC-058` Punkt 15).
   Vorschlag: Punkt 18 Schritt 5 nennt die beiden Kommandos.
2. **Punkt 2 gegen Punkt 4 beim Wort „Werkzeug".** Punkt 2 sagt, `eingaben` „ruft kein Werkzeug auf";
   dieselbe Tabelle verlangt für `PGR_ZWEIG` `git check-ref-format --branch`. Punkt 4 löst es auf (bis
   einschließlich `zustand` sind die lesenden `git`-Kommandos einschließlich `check-ref-format` erlaubt). Das
   Skript ruft in `eingaben` nur dieses eine Kommando und kein `ssh` und kein `make`; die Gegenprobe prüft
   genau das (`p4-nichts-geaendert`: höchstens ein `ssh`-Aufruf, von `git` nur die lesenden Kommandos).
   **Daran hängt eine Randform, die der Vertrag nicht entscheidet, und sie geht als Rückgabe an den
   Architect** (`AGENTS.md` §3.12; die Liste *Randform · Frage* steht im Bericht des Implementers): Was gilt,
   wenn `PGR_ZWEIG` ungültig ist **und** `git` am Host fehlt? Punkt 14 ordnet „Eingabe vor Werkzeug" (Exit 2),
   Punkt 2 bindet die Prüfung des Zweignamens aber an `git check-ref-format`, und Punkt 3 verlangt einen Fall
   `je Werkzeug fehlt` mit Exit 3. Beides zugleich ist nicht erfüllbar: Entweder die Prüfung bleibt ohne `git`
   aus (dann endet die Kombination mit Exit 3), oder `eingaben` lehnt ohne `git` jeden Zweignamen ab (dann ist
   `p3-fehlt-git` mit gültigem Zweig nicht erreichbar). Die Lieferung trägt die **erste** Antwort, weil sie
   beide Fälle des Punktes 3 erreichbar hält; der Fall `p3-fehlt-git-zweig-ungueltig` hält sie fest, und ein
   Satz im Vertrag kann sie bestätigen oder umdrehen (dann ändert sich eine Zeile im Skript und ein Fall).
3. **`PGR_SSH_OPTS` im Push.** Für die eigenen `ssh`-Aufrufe reisen die Optionen als einzelne Argumente,
   die keine Shell sieht (Punkt 2, Fall `p2-opts-keine-shell`). `GIT_SSH_COMMAND` ist nach Punkt 7 eine
   **Zeichenkette**, und `git` gibt sie an `sh -c`: eine Option mit Shell-Metazeichen wird dort ausgewertet.
   Das folgt aus den beiden wörtlich festgelegten Formen; Punkt 20 nennt als Grenze nur eine Option mit
   Leerraum. Vorschlag: Punkt 20 (7) nennt auch Metazeichen.
4. **Der Exit-Code von `git push` bei einer gescheiterten Verbindung** ist nicht Gegenstand des Vertrags
   (gemessen 128, nicht 1); die Gegenprobe prüft deshalb die **Form** der Zeile
   (`einliefern: git push <Remote> endete mit Exit <Zahl>`), nicht die Zahl.

**Zwei Zusagen, die der Lauf nicht belegt** (`SPEC-069` Punkt 17, akzeptiertes Negativ): dass ein echtes
`ssh` die Argumente so annimmt, dass der Befehl im Container läuft, und dass `make hooks-install` und
`make traeger-fetch` im Klon gelingen. Beides belegt der Nachweis oben (Schritt 2: `hooks-install:
core.hooksPath=.githooks`, `traeger-fetch: Traeger abgelegt … Digest verifiziert`, Schritt 4:
`ls …/.harness/` nennt `baseline skills state`).

### Nacharbeit nach dem Review vom 2026-10-21

**Sieben Entscheidungen des Architect vom 2026-10-21** (alle in diesem Slice, keine ADR berührt; Randformen
26 bis 32 in §6):

| Entscheidung | Vertragsstelle |
|---|---|
| Ausnahme von „Eingabe vor Werkzeug“: fehlt `git` am Host, bleibt allein die Prüfung von `check-ref-format` aus, der Lauf endet mit **Exit 3** (Review F-7, die Rückgabe des Implementers) | `SPEC-069` Punkt 14 (neuer Absatz), Punkt 2 (Tabelle `PGR_ZWEIG`), Punkt 3 (Fall und Mutation) |
| Aufräumen des Nachweises: die benannten Volumes und das Image-Paar werden **per Name** entfernt, der Beleg zeigt beide Zustände (Hinweis 1 der Lieferung) | `SPEC-069` Punkt 18 Schritt 5 und *Akzeptiertes Negativ* |
| Eine Option in `PGR_SSH_OPTS` mit Shell-Metazeichen: **Grenze**, nicht zugesagt und nicht gemessen (Hinweis 3) | `SPEC-069` Punkt 20 (7), Punkt 2 (Einschub) |
| Der Exit-Code von `git push` bei gescheiterter Verbindung bleibt **ungebunden**, akzeptiertes Negativ (Hinweis 4) | `SPEC-069` Punkt 20 (11) |
| `PGR_DEV` trägt denselben Zeichenvorrat wie `SPEC-058` Punkt 2 und wird **hier erneut geprüft** (Review F-3) | `SPEC-069` Punkt 2 (Tabelle, Absatz *Warum `PGR_DEV` hier erneut geprüft wird*, Fälle), Punkt 20 (7) |
| Die Plattform-Zusage lautet bash **≥ 4.4**, mit akzeptiertem Negativ und **ohne** Selbsttest (Review F-4) | `SPEC-069` Kopf *Plattform*, Punkt 19; Kopf des Skripts; §6 (Randform *Plattform*) |
| `ssh` und `git` am Host: die Lage bleibt, der Grund steht im Vertrag; `AGENTS.md` bleibt unverändert (Review F-17) | `SPEC-069` Punkt 3 (Absatz *Warum zwei Werkzeuge am Host*) |

**Was die Befunde geändert haben:**

- **F-1** (Ausgabe von `ssh` verschluckt): `zust_lesezug` ordnet die Zeilen des Lesebefehls den Schlüsseln
  zu und schreibt **jede andere Zeile auf stderr** (Punkt 4); der Fehlerpfad nennt sie weiter in der Meldung
  (Punkt 5 (5)). Die Attrappe von `ssh` schreibt auf Wunsch (`AT_SSH_LAUT`) im **Erfolgsfall** auf stderr.
- **F-2** (Pfadnamen-Expansion): `set -f` um die Teilung von `PGR_SSH_OPTS`.
- **F-3** (`PGR_DEV` in `GIT_SSH_COMMAND`): Zeichenvorrat geprüft, Kopf des Skripts berichtigt.
- **F-5** (`MAKEFLAGS` nur für `n`, `q`, `t`): drei Fälle ergänzt.
- **F-6** (Fälle ohne Mutation): neun Fälle haben jetzt je eine Mutation, dazu ein Fall, der liest, **wo der
  lokale `main` des Klons steht** (`p11-main-unbewegt`).
- **F-8** (§6 trug den Stand vor der Lieferung): §6 nennt die sieben Randformen und ihre Entscheidung.
- **F-10** (Zahlen): 225 Fälle, 194 Mutationen, Metrik 1326 — alle Zahlen am Stand `48ada9b2`.
- **F-11** (`--force-with-lease`, URL ohne Port): je ein eigener Fall bzw. eine eigene Mutation; das Muster
  von `p1-nie-push-force` endet jetzt auf `( |$)` und trennt beide Formen.
- **F-12** (`pgrep`): aus Voraussetzung, Kopf, Fragment-Kommentar und Pfad der Welten entfernt.
- **F-13** (Fahrt unter fremdem Benutzer, fremder Locale, unter Last): gefahren, siehe Läufe-Tabelle.
- **F-9** (Messung im Gegenproben-Container): gefahren, siehe Läufe-Tabelle und der Absatz zu `SPEC-060`.
- **F-14** (der Absatz *Kein weiterer Schnitt* in §4 kam mit einem Code-Commit): §4 bleibt inhaltlich
  unverändert — die Entscheidung ist die des Nutzers vom 2026-10-21, vom Koordinator übermittelt; dass ein
  Commit des Implementers sie in §4 setzte, ist unten als Beobachtung für den Planner vermerkt.
- **F-15, F-16** (INFO, Zuständigkeit Verifier): Das Ergebnis von `make gates` steht jetzt in der
  Läufe-Tabelle; `make gates-schnell-gegenprobe` (Namensregel von `SPEC-060` im Arbeitsbaum) ist **nicht**
  gefahren — der Auftrag dieses Laufs nennt nur die eigene Gegenprobe ([ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) *Übergang*), und der
  Reviewer hat `namensregel .` ohne Verletzung gefahren. Die vier Listen des Arbeitsbaums sind geprüft:
  `container-einrichten-gegenprobe` steht genau einmal in `GEGENPROBEN_CONTAINER`, in keiner Hälfte von
  `GATE_CHECKS`, nicht in `GATE_SCHNELL_AUSNAHME` und nicht in `GEGENPROBEN_HOST`.

**Mutationsprotokoll der Nacharbeit (Zusage · Mutation · roter Fall).** Gemessen mit `GP_ENTDECKE=1`, je
Mutation verlangt `mut` genau diese Menge.

| Zusage | Mutation | roter Fall |
|---|---|---|
| Punkt 4: die Ausgabe eines gerufenen Werkzeugs erscheint auf ihrem Strom — auch die von `ssh` im Lesezug | `m-p4-lesezug-stderr`: die Weitergabe der fremden Zeilen entfällt | `p4-werkzeug-stderr-ssh` |
| Punkt 2: `PGR_SSH_OPTS` wird an Leerraum geteilt, als einzelne Argumente übergeben, **ohne Pfadnamen-Expansion** | `m-p2-opts-glob`: `set -f` entfällt | `p2-opts-muster` |
| Punkt 2: `PGR_DEV` trägt nur `A-Za-z0-9._+-` und `/` | `m-p2-dev-zeichen`: die Prüfung entfällt | `p2-dev-metazeichen`, `p2-dev-leerzeichen` |
| Punkt 14, Ausnahme: ohne `git` bleibt allein die Prüfung von `check-ref-format` aus | `m-p2-zweig-ohne-git`: der Vorbehalt `command -v git` entfällt | `p3-fehlt-git`, `p3-fehlt-beide`, `p3-fehlt-git-zweig-ungueltig` |
| Punkt 2: Grenze `a` angenommen | `m-p2-zweig-kurz`: `a` abgelehnt | `p2-zweig-a` |
| Punkt 2: Grenze `slice-harness-devcontainer-einrichten` angenommen | `m-p2-zweig-lang`: Länge über 20 abgelehnt | `p2-zweig-kennung` |
| Punkt 2: `PGR_REMOTE` mit 64 Zeichen angenommen | `m-p2-remote-laenge-63`: Grenze 63 | `p2-remote-64` |
| Punkt 5 (3): die Identität kommt aus **allen Ebenen**, auch `--global` | `m-p5-ident-nur-lokal`: `config --local --get` | `p5-ident-global`, `p4-nichts-geaendert` |
| Punkt 5 (3): gültiges UTF-8 wird angenommen | `m-p5-ident-nur-ascii`: die Folge aus zwei Bytes entfällt | `p5-ident-utf8-gut`, `p9-zeichen-umlaut` |
| Punkt 7: das Remote heißt `PGR_REMOTE`, nicht `container` | `m-p5-remote-name-fest`: `remote add container` | `p5-remote-anderer-name` |
| Punkt 7: die URL trägt den Port | `m-p7-url-ohne-port`: Port aus der URL | `p2-vorgaben`, `p5-remote-andere-url`, `p5-remote-anderer-name`, `p5-remote-gleich`, `p7-remote-fehlt` |
| Punkt 8: der Klon entsteht am Pfad des Vertrags | `m-p8-clone-ziel`: anderes Ziel | `p8-fehlt`, `p13-erfolg`, `p13-umgebung`, `p13-zweiter-aufruf` |
| Punkt 11: ein fehlender Zweig entsteht mit `-c` | `m-p11-nie-c`: `switch` ohne `-c` | `p11-fehlt`, `p11-aus-origin-main`, `p11-main-unbewegt`, `p11-sonderzeichen` |
| Punkt 11: der lokale `main` des Klons wird **nie bewegt** | `m-p11-main-bewegt`: `branch -f main origin/main` nach dem Wechsel | `p11-main-unbewegt` |
| Punkt 14: je Nachbarpaar vertauscht — erstes Paar `eingaben`/`werkzeuge` | `m-p14-eingaben-nach-werkzeuge`: die zwei Schritte vertauscht | `p14-eingabe-vor-werkzeug` |
| Punkt 1: ruft nie `git push --force-with-lease` auf | `m-p1-nie-push-force-with-lease`: ein Schritt ruft es | `p1-nie-push-force-with-lease`, `p4-reihenfolge-der-schritte` |
| Punkt 16: je `MAKEFLAGS` aus der Umgebung mit demselben Flag | `m-p16-rezept-make` (Rezept mit `$(MAKE)`) und `m-p16-voraussetzung-fehlt` | `p16-makeflags-n/-q/-t` bzw. `p16-makeflags-i/-k/-j` |
| `SPEC-060` Punkt 16 (8): als root entfällt der Fall mit Dateirechten und **mit ihm seine Mutation** | eine Kopie ohne das Überspringen in `mut` meldet `rot sind [nichts]` (gemessen am 2026-10-21 im Lauf als root, vor dem Fix) | `m-p5-schluessel-lesbar` (die Mutation selbst, nicht ein Fall) |

**Die Läufe unter Last haben drei Fehler der Gegenprobe gefunden — alle drei behoben.** Die Fahrt nach
`SPEC-060` Punkt 16 (8) endete zweimal rot mit **je einem zusätzlichen** `p4-nichts-geaendert` (beim ersten
Lauf unter `m-p2-alle-gemeldet`, beim zweiten unter `m-p5-ident-steuerzeichen` und `m-p5-ident-utf8`), und
ein voller Lauf im Gegenproben-Container bei Last 4.86 ebenso; dieselben Mutationen einzeln waren grün.
Gesucht und gefunden mit der Diagnose, die derselbe Lauf eingeführt hat (`GP_ENTDECKE=1` nennt seither den
Fall, in dem ein roter Name gesetzt wurde) und einem Vergleich der zwei Momentaufnahmen:

1. **Die Momentaufnahme der Welt verlor unter Last Zeilen.** `snap` lief als
   `find … -printf … -exec cksum {} +` mit unterdrückter Fehlerausgabe; schlug unter Last ein `fork` für
   `cksum` fehl, fehlte die Zeile, und der Vergleich meldete eine Änderung, die es nicht gab. `snap` und
   `snapw` vergleichen jetzt Pfad, Rechte, Art, Link-Ziel, Größe und Änderungszeit (Nanosekunden) aus
   **einem** `find`-Aufruf ohne weiteren Prozess je Datei — robuster und schärfer als vorher, weil nun jeder
   Schreibzugriff auffällt, auch einer mit gleichem Inhalt (`BEO-REPO/test-faellt-unter-cpu-last-aus`).
2. **Die Attrappe von `ssh` schrieb ohne Not.** Sie setzte den Pfad in der Konfiguration des Klons nach
   jedem Befehl mit `sed -i` zurück, auch wenn er schon die Form des Containers trug — gleicher Inhalt, neue
   Änderungszeit. Sie schreibt jetzt nur, wenn der Pfad die andere Form trägt.
3. **`git status --porcelain` im Lesebefehl schreibt `.git/index` neu**, wenn die Stat-Daten im selben
   Sekundenfenster liegen (*racily clean*). Das ist eine Eigenschaft des gerufenen `git`, kein Zustand, den
   Punkt 4 meint; `.git/index` bleibt in der Momentaufnahme außen vor, mit Grund im Kopf der Funktion.
4. **`git` startet nach einem Commit seine Wartung im Hintergrund** und legt `.git/objects/maintenance.lock`
   an, die es wieder entfernt. Unter Last liegt die Sperre in der **ersten** Momentaufnahme (gleich nach dem
   Aufbau der Welt) und in der zweiten nicht mehr — der Vergleich meldete eine Änderung, die das Ziel nicht
   gemacht hat. Gefunden mit der Diagnose aus 1. bis 3. (`GP_ENTDECKE=1` nennt seit diesem Lauf den Fall und
   die ersten Zeilen des Unterschieds); jede Datei `*.lock` bleibt jetzt außen vor, mit Grund im Kopf der
   Funktion.

**Was die Zusage „bis einschließlich `zustand` ändert das Ziel nichts" nach der Härtung trägt:** Die
Momentaufnahme vergleicht Pfad, Rechte, Art, Link-Ziel, Größe und Änderungszeit **ohne** `.git/index` und
ohne `*.lock`; was in diesen beiden Ausnahmen geschieht, fängt sie nicht. Getragen wird die Zusage dort von
der **Positivliste der lesenden Kommandos** in derselben Prüfung (`^ssh ` höchstens einmal und von `git` nur
`rev-parse`, `config --get`, `remote get-url`, `check-ref-format` und die lesenden `-C`-Formen) — ein
schreibendes `git`-Kommando fällt an der Liste auf, auch wenn es nur den Index oder eine Sperre berührt
(gemessen: `m-p5-remote-gesetzt` mit `remote set-url` und `m-p5-ident-nur-lokal` mit `config --local --get`
färben `p4-nichts-geaendert` rot).

Nach den vier Korrekturen sind grün: der volle Lauf im Gegenproben-Container (176 s), die Fahrt nach
Punkt 16 (8) am Endstand unter Last als root (193 s, Last 10.72 zum Start — dieselbe Fahrt war vor der
vierten Korrektur rot) und zwei volle Läufe im Dev-Container mit Diagnose. **Zahl der Läufe:** Dieser Slice
hat die Gegenprobe oft gefahren (46 volle Sätze zu je rund 175 s); auf Anweisung des Nutzers vom 2026-10-21
wird jede Wiederholung seither eng ausgewählt (`GP_MUT`, `GP_NUR`), und der volle Satz läuft einmal am
Endstand.

**Für die Closure vorgemerkt (Beobachtungen dieses Laufs, keine Register-Einträge — die schreibt der
Planner): (1) Die Größe liegt mit 1341 über der Marke 900, der Faktor Gegenprobe zu Skript bei 2,8
(`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`). (2) Ein Kind-Lauf der Gegenprobe mit 42 Fällen fiel
unter Last (Last 10) in `p17-ordner-fall` und `p17-budget-*` aus; seit `GP_NUR` je Mutation fährt die Kopie
einen Fall (`BEO-REPO/test-faellt-unter-cpu-last-aus`, `BEO-REPO/lastabhaengige-zaehlung`). (3) Eine
Attrappe von `ssh` ohne Rechner und Befehl hing am Standardeingang, den `git` offen hält — seit dem Lauf
endet sie wie `ssh` mit 255; ein Fall endet zusätzlich nach zwei Minuten (`timeout`). (4) Der Pfad der Welt
trug den Namen der Mutation, und eine Prüfung fand ihren Text im Pfad statt in der Meldung (falsches Grün);
die Welten heißen jetzt neutral. (5) `u-boot` am Host steht auf 0.13.0, `SPEC-058` nennt 0.10.0 als
gemessene Umgebung; für dieses Ziel ohne Belang (Punkt 20 (3)), der Nehmer ist
`slice-harness-container-start-u-boot-0-13`. (6) Sieben Randformen kamen erst in der Umsetzung und im
Review auf und wurden am 2026-10-21 entschieden; eine davon (F-7) hatte der Implementer im Code
beantwortet und zurückgegeben (`BEO-REPO/randform-im-code-entschieden-dann-zurueckgegeben`), drei kamen
aus Messungen des Reviewers (`BEO-REPO/spec-randform-erst-im-review-entschieden`). (7) Der Absatz *Kein
weiterer Schnitt* in §4 kam mit dem ersten Code-Commit in den Plan, obwohl `AGENTS.md` §3.9 dem Implementer
§1, §3, §6 und den Kopf nennt (Review F-14); die Entscheidung selbst ist die des Nutzers vom 2026-10-21.
(8) Neun von 216 Fällen hatten keine Mutation, drei vom Vertrag verlangte Mutationen fehlten, und zwei
Aufzählungselemente waren nur mittelbar geprüft (Review F-5, F-6, F-11,
`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`). (9) Die Messung lief zuerst im Dev-Container statt im
Gegenproben-Container, und die Tabelle in `SPEC-060` stand schon auf `aufgenommen` (Review F-9).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie ist so grob wie das Repo, trägt aber eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer für Werkzeug und Container: `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (Randform 9), `BEO-REPO/nested-runtime-podman-oder-docker-offen`, `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Zählerstände liest der Implementer aus den Dateien in `evidence/` vor dem Start (kein Zähler wird hier gesetzt). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel; die verkörperten (§3.10 bis §3.14) gelten.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Schnitt einträgt):** Liefer-Punkte 3 (Ziel mit Ablauf, Gegenprobe, Dokument; die Sendung „Arbeit im Container“ liegt im Dokument-Punkt, die Sendung „Quellenauflösung“ liegt beim Zugang, Gerüst und Start beim Geber). Schichten 2 (Harness-Werkzeug mit Gegenprobe und Gate-Index; Nutzer- und Wartungs-Doku `docs/developer/`; `harness/README.md` ist Index desselben Werkzeugs, Planung und Vertrag sind keine Schicht). Schätzung in der Metrik von §4: 565 bis 830, **unter der Marke 900** (Meldeschwelle 700 in der oberen Hälfte überschritten). Vor dem Schnitt 995 bis 1480 (ungeteilt, mit Gerüst und Start). **Nachgezählt am 2026-10-18 im Commit, der die Aufteilung des Gebers in `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten-start` einträgt:** Die Bindung an den Container kommt jetzt von zwei Gebern; sie ist keine neue Lieferung dieses Slice. Liefer-Punkte 3 und Schichten 2 unverändert, Schätzung 565 bis 830 unverändert, unter der Marke 900. `grep` nach der Kennung dieses Slice findet ihn in den Gebern `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten-start` (§1, nehmen die Bindung auf), im Folge-Slice `slice-harness-devcontainer-einrichten-zugang`, in `slice-harness-devcontainer-u-boot-0-10` (Geber, `done/`), in den Plänen der Release-Welle (keine Bindung) und in den Beobachtungen.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit, der die Bindung aus `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` einträgt): Die Bindung ändert weder Liefer-Punkte noch Schichten dieses Slice; sie verlangt im Dokument dieselben Abschnitte, die DoD „Dokument“ schon umschreibt. Zählung bleibt wie im Schnitt vom 2026-10-18 (§8 oben).

**Nachgezählt** (`AGENTS.md` §3.9, §3.13, Planner, 2026-10-18, im Commit, der den Kettenumbau `slice-harness-gegenproben-aus-gates` mit seinem Vertrag (`SPEC-060`) nachzieht): Die Gegenprobe hängt an einer der beiden Listen statt an `GATE_CHECKS` und steht nicht in der Ausnahmeliste; die Zeile zu `gates-schnell.mk` entfällt. Liefer-Punkte und Schichten ändern sich nicht (die Zählung in §8 oben bleibt).

**Nachgezählt** (`AGENTS.md` §3.13, Architect, 2026-10-21, im Commit, der `SPEC-069` und [ADR-0057](../../adr/0057-einrichtung-im-dev-container-gegenprobe-ohne-container.md) einträgt): Liefer-Punkte **3** unverändert (Ziel mit Ablauf — der Nachweis aus `SPEC-069` Punkt 18 ist kein eigener Liefer-Punkt: er hat kein Artefakt, keinen Treiber und keine Zusage, sondern fünf Schritte aus vorhandenen Zielen und einen Beleg in §7 —, Gegenprobe, Dokument). Schichten **2** unverändert (Harness-Werkzeug mit Gegenprobe und Gate-Index; Wartungs-Doku `docs/developer/`; `SPEC-060` Punkt 16 und `SPEC-069` sind derselbe Ort wie das Werkzeug, Planung und Vertrag sind keine Schicht, `AGENTS.md` §3.13). **Schätzung berichtigt:** 815 bis 1210 statt 565 bis 830 (§4); die obere Hälfte liegt über der Marke 900. Der Architect schneidet nicht selbst: Das ist die Entscheidung des Planners (§4, Empfehlung in drei Stufen). `grep` nach `SPEC-069` findet den Vertrag im Kopf, in §1, §2, §3, §4, §5 und §6 dieses Plans und in `harness/vertraege/README.md`.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit, der `slice-harness-gegenproben-aus-gates-fahrstelle` schneidet): Die Bindung steht jetzt unter der Kennung des Gebers `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (Abschnitt über den Gegenproben-Container im Dokument und erste Fahrt); der Hook mit seiner Gegenprobe (`slice-harness-gegenproben-aus-gates-fahrstelle`) berührt das Dokument nicht. Die Bindung ändert weder Liefer-Punkte noch Schichten dieses Slice; Zählung bleibt wie im Schnitt vom 2026-10-18 (§8 oben). `slice-harness-devcontainer-einrichten-start` trägt dieselbe Bindung für §3 bis §5.

**Nachgezählt** (`AGENTS.md` §3.9, §3.13, Planner, 2026-10-21, im Commit, der den Dokument-Teil nach `slice-harness-devcontainer-einrichten-doku` schneidet; Entscheidung des Nutzers vom selben Tag): **Liefer-Punkte 2** statt 3 — Ziel mit Ablauf und Gegenprobe; der Dokument-Punkt liegt beim Nehmer, der ihn in §1 und DoD mit der Kennung dieses Slice annimmt. Mit ihm gingen die Dokument-Hälfte der Sendung „Arbeit im Container“ (von `slice-harness-devcontainer-u-boot-0-10`), die Bindung von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` und drei Dokument-Risiken aus §6. **Schichten 1** statt 2 — Harness-Werkzeug mit Gegenprobe und Gate-Index (`tools/`, `harness/mk/`, `harness/README.md`, `SPEC-060` Punkt 16); die Wartungs-Doku `docs/developer/` ist die Schicht, die wegfällt (Teilung nach `slice-harness-d-check-v0-85` §8, `AGENTS.md` §3.13). **Schätzung 745 bis 1090** statt 815 bis 1210 (§4), Pfade der Marke `tools harness` statt `tools harness docs/developer`; Meldeschwelle 700 und **Marke 900** unverändert, die obere Hälfte liegt weiter über der Marke, mit der Naht im Vertrag vorab benannt (§4). `grep` nach der Kennung dieses Slice findet ihn im Nehmer `slice-harness-devcontainer-einrichten-doku` (Kopf, §1, §2, §4, §6, §8), in den Gebern `slice-harness-devcontainer-einrichten-geruest` und `slice-harness-devcontainer-einrichten-start` (§1), im Folge-Slice `slice-harness-devcontainer-einrichten-zugang`, in `welle-dev-container` (§4 und §5) und in den Beobachtungen.
