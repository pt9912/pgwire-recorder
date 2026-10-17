# Slice slice-harness-devcontainer-austausch-repo: Dev-Container tauscht Stände über ein Bare-Repo im Volume aus, ohne Mount des Host-Arbeitsbaums (`docs/developer/devcontainer.md`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein Abnahmeszenario und
kein Meilenstein hängt an ihm. Angelegt am 2026-10-17 nach Entscheidung des Nutzers vom selben Tag
(„ja, Austausch-Repo-Umbau einplanen“); eingeordnet in `next/`; beansprucht am 2026-10-17 nach dem Schließen von `slice-harness-integration-volumes`,
als paralleler Lauf neben `slice-harness-abdeckung-gate` (§4 *Start*, §6 *Paralleler Lauf*).

**Bezug:** kein Anspruch des Produkts und kein Lastenheft-Bezug — der Slice ändert die Beschreibung der
Arbeitsumgebung der Entwickler. Vorgänger: `slice-harness-devcontainer-doku` (`done/`), dessen Kopf
ebenfalls keinen Lastenheft-Bezug führt; der Slice schreibt dessen Dokument um. Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912 (Zweig slice/devcontainer-austausch-repo)

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `docs/developer/devcontainer.md` beschreibt den Austausch zwischen Host und Dev-Container
über ein Bare-Repo im Volume des Containers statt über einen Mount des Host-Arbeitsbaums, und
jede Aussage dazu ist am laufenden Container nachgeführt oder steht in §13 als nicht geprüft.

Heute (Vorgänger): §6 und §9 arbeiten mit dem Mount `/host-repo` (Mount `ro`) und `git clone /host-repo`.
Der Container sieht damit den ganzen Arbeitsbaum des Hosts (auch uncommittete und ignorierte Dateien wie
`.harness/state/` und `.claude/`), `PGR_REPO` ist Pflichtvariable, und der Rückweg läuft per `git fetch`
über SSH auf den ausgecheckten Klon oder per Bundle.

**Zielbild — von den Proben bestätigt, Befunde in §7** (Benennung am Host `container/container/<thema>` bei Präfix im Container; Hostschlüssel; Bare-Repo mit `--initial-branch=main`; Vorschlag des Nutzers, 2026-10-17: lokal einen Klon anlegen,
von dem aus gepullt und gepusht wird; gemountet wird nur das Volume):

- Im Volume `pgr-dev-workspace` liegt ein Bare-Repo `/workspaces/austausch.git`; daraus entsteht der
  Arbeitsklon `/workspaces/pgwire-recorder` (Clone aus dem Bare-Repo, Push und Pull lokal).
- Der Host ist Remote des Containers: `git remote add container ssh://vscode@127.0.0.1:${PGR_SSH_PORT:-2222}/workspaces/austausch.git`
  mit `core.sshCommand` (Schlüssel `$PGR_DEV/ssh/id_ed25519`, `IdentitiesOnly=yes`). `git push container main` liefert
  den Stand ein, auch unveröffentlichte Commits; `git fetch container` holt die Branches des Containers
  (`container/<thema>`, Benennung am Host: Randform in §6).
- Der Mount `/host-repo`, die Variable `PGR_REPO` im Block `compose.yaml` (und in den Befehlen, `.env`) und die Varianten A und B in §9 entfallen.
  §10 (mehrere Container: je ein Remote mit eigenem Port), §11 (der Container sieht den Host-Arbeitsbaum nicht mehr;
  Mounts nur noch drei), §13 (Nicht geprüft) und die Aussagen zu den „vier Mounts“ in §7 ziehen nach.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Umbau auf die neue Version von `u-boot`: Sie ist nicht erschienen; Dockerfile und Compose
  bleiben handgebaut, wie das Dokument sie heute nennt (anderer Vorgang).
- Keine Änderung der Dateiinhalte von `Dockerfile` und `start.sh` und keine Änderung am Messort außer dem
  additiven Anlegen des Bare-Repos (§3): Reicht der Inhalt für das Zielbild nicht (zum Beispiel fehlt `git`
  im Image), hält der Implementer an und gibt es zurück (§6), statt die Datei zu ändern; eine Änderung wäre
  eine Entscheidung des Nutzers.
- Kein Produktcode, keine Änderung an Hooks, Gates, Makefile, Positivliste oder Spezifikation
  (Schicht-Abgrenzung). `docs/developer/` steht seit `slice-harness-stop-hook-positivliste-developer` auf der
  Positivliste des Stop-Hooks; das genügt.
- Egress und Ablage der Anmeldedaten von `claude` (Betriebsfrage, offene Frage 2 des Vorgängers) und die
  Frage nach dem Hostschlüssel bei der Erstverbindung: nicht Gegenstand. Berührt der Zielablauf die
  Hostschlüssel-Frage unvermeidbar (jeder `git push container` ist eine SSH-Verbindung), nennt das Dokument,
  was gelaufen ist, und lässt den Rest in §13; die Entscheidung über eine Einstellung trifft der Nutzer (§6).
- Keine Aussagen zu Windows- oder macOS-Hosts und kein Schritt ohne Probe — wie im Vorgänger (§3.11).
- Kein Eintrag in der Roadmap: Der Slice ist wellenlos.

## 2. Definition of Done

- [x] `docs/developer/devcontainer.md` beschreibt den Austausch über das Bare-Repo nach §1; die Abschnitte
      §3.5 (Block `compose.yaml` ohne Mount und ohne `PGR_REPO`), §6, §7, §9, §10, §11, §12 und §13 sind
      nachgezogen; `/host-repo` und `PGR_REPO` kommen im Dokument nicht mehr vor (`grep`); jede Aussage ist
      durch einen Lauf am Container belegt (Tabelle *Aussage · Probe · Beobachtung* in §7) oder steht in §13.
- [x] `make docs-check`, `make kopf-check` und `make gates-schnell` grün auf dem Endstand (nur
      `docs/developer/` und die Planung ändern; ein voller `make gates` ist nicht nötig, §5).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/developer/devcontainer.md` | update | Liefer-Punkt: das Dokument nach §1 (Umbau der Abschnitte aus §2) |

**Probe (Implementer, vor dem Schreiben, am Container; Ergebnis je Zeile in §7):**

1. Bare-Repo im Container anlegen (`git init --bare /workspaces/austausch.git` als `vscode`, per SSH); Besitzer, Rechte
   und UID prüfen (`stat`; UID 1000 am Host und im Container); `git` und `git-receive-pack` im Image vorhanden?
2. Push vom Host über SSH in das leere Bare-Repo (`git remote add container …`, `core.sshCommand`, `git push container main`),
   auch mit einem Commit, der nicht veröffentlicht ist (Scratch-Klon des Host-Repos, das echte Host-Repo bleibt unberührt, wie im Vorgänger).
3. Clone im Container aus dem Bare-Repo nach `/workspaces/pgwire-recorder`; `git config --global` für Name und E-Mail; Commit auf
   `container/<thema>`, Push in das Bare-Repo (lokal).
4. `git fetch container` am Host (Scratch-Klon): wie heißen die Branches am Host (`container/<thema>` oder `container/container/<thema>`)?
   Zweiter Push des Hosts nach dem Push des Containers; Pull im Arbeitsklon.
5. Überleben von `u-boot down` und `u-boot up`: Bare-Repo und Arbeitsklon im Volume, Remote am Host, SSH ohne Warnung (Hostschlüssel nach Neubau).
6. Zweite Instanz mit anderem Port und eigenem Remote-Namen am selben Host-Repo; zwei Remotes nebeneinander.
7. Der Block `compose.yaml` ohne Mount `/host-repo` und ohne `PGR_REPO`: Start eines Scratch-Containers (anderer Name, Port, Volumes, wie im Vorgänger) aus den
   Codeblöcken des Dokuments; `docker inspect` zählt die Mounts.
8. Am echten Container `/Development/pgr-dev` waren die Schritte 1 bis 5 vorgesehen. Beim Start des Laufs (2026-10-17) lief dort kein Container
   (`docker ps -a` leer; Image `pgr-dev:latest` und die Volumes `pgr-dev_pgr-dev-docker`, `pgr-dev_pgr-dev-workspace` vorhanden). Ihn zu starten hieße, ihn neu zu erzeugen (Festlegung 2);
   deshalb liefen alle Schritte 1 bis 7 im Scratch-Container (Projekt `pgr-aus`, Port 2322, Dockerfile und `start.sh` aus den Codeblöcken des Dokuments, Compose-Block ohne Mount),
   Schritt 6 mit einem zweiten (Projekt `pgr-aus2`, Port 2323). Das echte Verzeichnis, seine Volumes und das Image blieben unberührt (§7).

**Docker-Last.** Die Proben starten Container und Daemon-Läufe (Scratch-Container, Volumes, Port). Sie laufen nie
gleichzeitig mit einem Gate-Lauf oder einer Volume-Messung eines anderen Slice (`pgrep -af "make gates"` ohne Treffer, `docker volume ls -q` vorher und
nachher gleich; `docker ps -a` zeigt nur `pgr-dev` und die Scratch-Container), Aufräumen wie im Vorgänger (`docker compose down -v`, `docker rmi`).

- Das Dokument bleibt Ist-Zustand ohne Chronik (§3.11) und nennt keine Kennungen dieses Repos; Pfade des Hosts stehen nicht fest (`PGR_DEV` bleibt, `PGR_REPO` entfällt).
  Codeblöcke, die der Probe nicht entsprechen, stehen nicht im Dokument (Unterschied zum Messort nennt das Dokument, wie §3.5 es heute tut).
- Das Dokument macht aus dem Bestandscontainer keine Anleitung zum Umbau: Wer schon einen Container mit `/host-repo` hat, findet in §4 den Weg über `docker compose`
  (nur beschrieben, was gelaufen ist, sonst §13).
- Der Satz „Behandle den Container nicht als Sandbox“ aus §11 bleibt wahr und bleibt: Das Bare-Repo nimmt Pushes von `vscode` an; der Host liest Branches
  aus dem Container und prüft sie vor dem Merge. Was entfällt, ist nur die Sicht des Containers auf den Arbeitsbaum (lesbare ignorierte Dateien).

## 4. Trigger

**Start** (`next` → `in-progress`): erfolgt nach dem Schließen von `slice-harness-integration-volumes` (`done/`). Der Anspruch ist ein reiner `git mv` auf dem Hauptzweig vor der Arbeit.
Der Lauf geht parallel zu `slice-harness-abdeckung-gate` (§6 *Paralleler Lauf*); die Proben beachten den geteilten Docker-Daemon (§6, Festlegung 4).

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt. Schätzung 220 bis 340 (Planner, 2026-10-17; ein Umbau schreibt
Abschnitte um und löscht dabei: grob 110 Löschungen und 150 Einfügungen, dazu der Aufschlag der Nacharbeit; der Vorgänger lag bei 441 gegen 150 bis 300
geschätzt, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, deshalb der Spielraum nach oben). Meldeschwelle 350, **Marke 450**.
Zusätzlich gilt die Länge des Dokuments: Es hat vor dem Umbau 498 Zeilen (`wc -l`); über **560** Zeilen am Endstand ist es zu lang für eine Review-Sitzung
und geht zurück zur Zerlegung. Die Marke des Vorgängers (500) galt dem Diff; hier ist sie niedriger, weil der Umbau löscht und Löschungen mitzählen.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 450, oder das Dokument hat am Endstand mehr als 560 Zeilen, oder die Schritte
  verlangen mehr als dieses eine Dokument (zum Beispiel eine Änderung an `start.sh` oder `Dockerfile`).
- `in-progress` → `open` (blockiert): Das Zielbild trägt nicht (zum Beispiel Push in das Bare-Repo scheitert an Rechten, die nur eine Änderung des Containers behebt) und die Entscheidung
  darüber liegt beim Nutzer; oder `/Development/pgr-dev` läuft nicht mehr und lässt sich ohne Neubau nicht nachführen.

## 5. Closure-Trigger

DoD vollständig, `make docs-check`, `make kopf-check` und `make gates-schnell` grün auf dem Endstand, Review- und Verifikations-Report liegen
vor, Closure-Notiz mit Lerneintrag geschrieben. Ein voller `make gates` ist nicht verlangt: Es ändern sich nur `docs/developer/` und die Planung,
beides auf der Positivliste des Stop-Hooks; ändert der Slice dennoch einen Pfad außerhalb, verlangt der Hook den vollen Lauf.

## 6. Risiken und offene Punkte

Dieser Slice liefert eine Beschreibung, keinen Vertrag (§3.12 gilt nicht im engen Sinn). Die Randformen des Ablaufs sind dennoch vor dem Schreiben
zu entscheiden; entschieden heißt hier: durch die Probe aus §3 belegt, sonst in §13 des Dokuments. Eine, die der Implementer erst beim Schreiben bemerkt
und nicht probt, hält er an und gibt sie an den Architect zurück (`.claude/commands/implement-slice.md`, Randform-Rückgabe).

- **Randformen des Austauschs** (Liste, je Form eine Probe oder ein Eintrag in §13): leeres Bare-Repo (Clone vor dem ersten Push: „leeres Repository“, Reihenfolge: erst
  Push, dann Clone); Standardbranch des Bare-Repos und `main`; Wiederholung der Befehle (`git init --bare` auf vorhandenem Verzeichnis, `git clone` in vorhandenes Verzeichnis,
  `git remote add` bei vorhandenem Remote); Push eines anderen Branches als `main`, Push nicht fast-forward (Host und Container schreiben beide auf `main`: Entscheidung
  gesucht, etwa Container arbeitet nur auf `container/<thema>`); `git fetch container` ohne Branch des Containers; Benennung am Host (`refs/remotes/container/…` und der Branch
  `container/<thema>` ergeben möglicherweise `container/container/<thema>`); Besitzer und Rechte im Volume (frisches Volume gehört `root`; Bare-Repo als `vscode` angelegt, UID 1000;
  andere UID am Host, `dubious ownership`); Überleben von `u-boot down` und `up`; zweiter Container mit gleichem Remote-Namen oder gleichem Port; Hostschlüssel nach Neubau des Images;
  `core.sshCommand` je Repo statt Umgebung (`GIT_SSH_COMMAND`), Wirkung auf andere Remotes (`origin`); Abbruch mitten im Push. — **Stand des Laufs:** jede Form hat eine Zeile in §7 oder steht in §13 des
  Dokuments (nicht gelaufen: andere UID am Host, `origin` über `ssh`, Fetch aus leerem Bare-Repo über `ssh`). — **Ausgang (Closure):** entfallen — jede Form hat eine Zeile in §7 oder steht in §13 des Dokuments; die nicht gelaufenen (andere UID am Host, `origin` über `ssh`, Fetch aus leerem Bare-Repo über `ssh`) sind dort als nicht geprüft geführt. Eine Randform, die das Review nachträglich fand, ist nicht darunter (F-3, F-4 betrafen Befehle, keine Randform des Ablaufs).
- **Das Zielbild trägt nicht ohne Eingriff in `Dockerfile` oder `start.sh`** (zum Beispiel `git` im Image fehlt, Rechte an `/workspaces`). — **Ausgang (Closure):** entfallen — `git`, `git-receive-pack` und die Rechte an `/workspaces` reichen (Probe 1); Dockerfile und `start.sh` blieben unverändert (`diff` leer), die Rückführung nach §4 trat nicht ein.
- **Hostschlüssel-Frage wird berührt** (SSH bei jedem Push und Fetch; nach Neubau des Images andere Hostschlüssel, Warnung am Host). — **Ausgang (Closure):** eingetreten, Entscheidung offen → weiter offen im Register `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (neu, 1×). Die Warnung trat in der Probe auf (§7); das Dokument nennt sie in §5 und §12 und erfindet keine Einstellung; die Entscheidung trifft der Nutzer (Frage in §7).
- **Docker-Last stört einen anderen Lauf** (Gate-Lauf oder Volume-Messung eines anderen Slice) oder umgekehrt. — **Ausgang (Closure):** entfallen — die Vorbedingung `pgrep -fa "^make gates"` war vor jeder Probe leer (sechsmal), der Zustand von Volumes, Containern und Netzen war nach dem Aufräumen gleich dem davor (§7).
- **Probe am echten Container ändert dessen Volume** (Bare-Repo, Arbeitsklon bleibt unberührt). — **Ausgang (Closure):** entfallen — es lief kein echter Container, die Proben liefen im Scratch-Container und `/Development/pgr-dev` blieb unberührt (§7); der Zustand vor und nach (`ls /workspaces`) steht in §7, das Bare-Repo wird entfernt oder auf Wunsch des Nutzers behalten. Im Lauf lief kein echter Container (§3 Schritt 8); die Proben liefen im Scratch-Container, der Zustand vor und nach steht in §7.
- **Das Dokument veraltet still gegen den Container** (kein Sensor; wie im Vorgänger) und **Aussagen über die Probe hinaus** (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; die Reviews des Vorgängers fanden sie zweimal).
  — **Ausgang (Closure):** eingetreten → Auftreten im Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Evidence dieses Slice); das Review fand drei MEDIUM und fünf LOW dieser Klasse, alle im Dokument eingegrenzt oder nach §13 (Nacharbeit `00959d52`). Das Veralten des Dokuments gegen den Container selbst: weiter offen, im selben Eintrag geführt (kein Sensor).
- **Größe:** Der Umbau liegt über der Marke (Vorgänger 441 gegen Schätzung 150 bis 300). — **Ausgang (Closure):** entfallen — gemessen 279 (Marke 450, Meldeschwelle 350), Dokument 559 Zeilen (Grenze 560); die Schätzung 220 bis 340 hielt (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, kein Auftreten).

**Festlegungen zu den offenen Punkten** (2026-10-17, beim Anspruch; Empfehlung, die der Nutzer nicht als Entscheidung bestätigt hat, er kann sie ändern;
einzige Ausnahme: der Hostschlüssel, der offen bleibt):

1. Das Bare-Repo, das die Probe am echten Container `/Development/pgr-dev` anlegt, wird nach der Probe entfernt (im Lauf gegenstandslos: Es gab keinen Container, §3 Schritt 8).
2. Der echte Container wird nicht neu erzeugt und nicht mit `down`/`up` angefasst. Das Überleben von `down`/`up` (Probe 5) läuft im Scratch-Container; der Bestand bleibt unangetastet.
3. Hostschlüssel: Das Dokument nennt nur, was gelaufen ist, und empfiehlt keine Einstellung (`StrictHostKeyChecking`, `UserKnownHostsFile`). Die Entscheidung darüber liegt beim Nutzer und bleibt offen; sie steht in §13 des Dokuments.
4. Vorbedingung jeder Container-Probe: `pgrep -fa "^make gates"` ist leer, weil der Docker-Daemon mit dem Lauf von `slice-harness-abdeckung-gate` geteilt ist. Das Dokument selbst darf jederzeit geschrieben werden.

**Paralleler Lauf.** Dieser Slice läuft neben `slice-harness-abdeckung-gate` (`in-progress/`, anderer Implementer-Lauf), nach WIP-Limit pro Lauf
(Baseline-Regelwerk `modul-05-planning-harness.md`; `Verantwortlich:` unterscheidet die Läufe per Zweig). Eigener Worktree (`../pgwire-recorder-austausch`) und Zweig
`slice/devcontainer-austausch-repo`; Fast-Forward-Push auf `main` erst nach Rebase. Konflikte unter `docs/plan/planning/` sind nicht zu erwarten (verschiedene Dateien).

## 7. Closure-Notiz

**Closure (Planner, 2026-10-17).** Der Slice liefert `docs/developer/devcontainer.md` auf Bare-Repo im Volume (ein Liefer-Punkt, eine Schicht). Belege:
Umbau `9f8faef9`, Gate-Läufe und Größe `1e70abc3`, Verifikation `97aea4f9` (bestätigt, 2 Findings), Review `ed7cd83e` (0 HIGH, 3 MEDIUM, 6 LOW, 3 INFO, merge-blockierend),
Nacharbeit `00959d52`. `make docs-check`, `make kopf-check` und `make gates-schnell` waren auf dem Endstand im Worktree grün (Angabe des Aufrufers, vom Planner nicht wiederholt, dazu `make docs-check` und `make kopf-check` vor den Closure-Commits).

**Ausgang je Befund des Reviews.** F-1 (Vorbehalt zu `privileged`/`sudo`), F-2 (Ablage der Anmeldedaten), F-3 (`$PGR_SSH_PORT`), F-4 (Arbeitsverzeichnis des Blocks), F-5 (Fetch gegen Push), F-6 (Wortlaut zu `u-boot init`), F-9 (`git pull` auf `main` gefasst): **behoben** bzw. **eng gefasst** im Dokument.
F-7 (Aussagen ohne Zeile in der Tabelle): **eng gefasst**, die Sätze aus der Probe des Vorgängers tragen die Kennzeichnung (siehe *Aussagen aus der Probe des Vorgängers* unten). F-8 (Plan zitiert Sätze, die nicht stehen): **behoben** im Plan.
F-10 (gefolgert): **gekennzeichnet** im Dokument. F-11 (Locale der Meldungen): **gekennzeichnet** in §6. F-12 (Relativsatz) und V-F2: **behoben**. V-F1 (`u-boot up` selbst nicht nachgefahren): **eng gefasst** („Messung von `u-boot`“).
Was das Dokument nicht belegt (Root-Zugriff über Gerätedateien des Hosts, Ablage der Anmeldung von `claude`, Bestandscontainer), steht in §13 des Dokuments.

**Die Nacharbeit durchlief keinen erneuten Review und keine erneute Verifikation.** Sie enthält engere Fassungen, Kennzeichnungen und Einträge nach §13, aber auch einzelne neue ausführbare Zeilen, die kein Lauf belegt hat:
§10 `export PGR_SSH_PORT=2223` bzw. `-p 2223` und der Remote `container2` im Host-Repo, §9 `git switch main && git pull`. Dazu keine neue Container-Probe nach der Nacharbeit.
Das ist ein **bewusstes Restrisiko**; Ausgang *weiter offen* → `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` (neu, 1×). Die Aussage dazu ist weder durch Review noch durch Verifikation gedeckt; sie gilt als gelesen, nicht gelaufen.

**Ausgang je Risiko (§6).** Jedes Risiko trägt seinen Ausgang an Ort und Stelle: Randformen, Eingriff in `Dockerfile`/`start.sh`, Docker-Last, Probe am echten Container, Größe: entfallen mit Begründung;
Hostschlüssel: weiter offen → `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`; Veralten und Aussagen über die Probe hinaus: eingetreten → `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.

**Marke §4.** `git diff --numstat 92cfaf23..00959d52 -- docs/developer` = 170 + 109 = **279** (Meldeschwelle 350, Marke 450 nicht erreicht; Schätzung 220 bis 340 hielt). Dokument **559** Zeilen (Grenze 560, ein Zeilenpuffer). Die Rückführung nach §4 trat nicht ein.

**Lerneintrag.** Form: benannte Lücke, im Register geführt (keine Zeile in `AGENTS.md` oder `.claude/`). Beobachtung: Ein Dokument über Sicherheit und Isolation verlor beim Umbau einen Vorbehalt der Vorfassung
(„gegen einen Prozess mit `sudo` ist es keine Zusage“); der Satz „also liegen darin keine Dateien des Hosts“ stand danach ohne ihn. Das Review fand es (F-1, MEDIUM), die Verifikation nicht: Sie fuhr die Kernprobe
(drei Mounts, Austausch, `down`/`up`) und bestätigte, was das Dokument zusagt, aber nicht, was beim Umbau aus einer Sicherheitsaussage gestrichen wurde. Probe gegen Zusage fängt Gestrichenes nicht; dafür braucht es den Vergleich der Aussagen
der Vorfassung mit der neuen. Zweitens bleibt die Nacharbeit nach dem Review unbelegt, wo sie Neues schreibt (Restrisiko oben). Vorschlag des Planners, nicht eingetragen: Beim Umbau eines Sicherheits- oder Isolationsdokuments nennt der Plan die
Vorbehalte der Vorfassung und der Review-Auftrag vergleicht sie mit dem Diff. Ob das eine Regel wird, entscheidet der Nutzer (Frage unten). Register: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Evidence, Auftreten Nr. 74 in den Dateien; Stand `verkörpert`, Regel bleibt),
`BEO-REPO/plan-folgt-korrektur-nicht` (Evidence, F-8), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (kein Auftreten, kein Evidence), neu `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` und `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (je 1×).
Keine Klasse erreicht mit diesem Slice 3× neu: `zusage-im-kommentar-weiter-als-pruefung` und `plan-folgt-korrektur-nicht` sind längst verkörpert (`AGENTS.md` §3.11 bzw. §3.9) und bleiben es; es wurde keine Regel geändert.

**Paarungen.** Anker: keine neue verkörperte Regel, deshalb kein neuer `seit slice-harness-devcontainer-austausch-repo`. Folge-Slice: keiner genannt; keine Sendung an einen anderen Slice (`grep` nach der Kennung in `next/` und `open/` leer).
Register: alle genannten Kennungen existieren als Verzeichnis und führen ein nicht leeres `evidence/`. Wellenlos: kein Eintrag in der Roadmap.

**Fragen an den Nutzer (offen, nichts geändert).**

1. Hostschlüssel: Soll das Dokument eine Einstellung nennen (`StrictHostKeyChecking`, `UserKnownHostsFile`) oder bleibt es bei der Beschreibung der Warnung? Die Probe am echten `/Development/pgr-dev` entfiel (Container lief nicht, Neuerzeugung ausgeschlossen); das Verhalten dort ist nicht gemessen.
2. Soll beim Umbau eines Sicherheits- oder Isolationsdokuments der Vergleich der Vorbehalte mit der Vorfassung eine Regel werden (Vorschlag oben)?
3. Soll die Nacharbeit mit den neuen Zeilen (§9, §10) vor dem Merge einen Lauf oder einen zweiten Review bekommen?

### Belege des Implementers

**Umgebung der Proben** (2026-10-17). Docker-Vorbedingung vor jeder Container-Probe: `pgrep -fa "^make gates"` ohne Treffer (sechsmal geprüft, nie ein Lauf).
Objekte nur mit Präfix `pgr-aus`: Projekte `pgr-aus` (Container `pgr-aus-dev`, Port 2322) und `pgr-aus2` (Container `pgr-aus2-dev`, Port 2323), Verzeichnisse im Scratchpad,
die aus dem Dokument gelesen wurden: `Dockerfile`-Ergänzung (Zeilen 120 bis 140 der Fassung vor dem Umbau, ohne Änderung; `diff` gegen `/Development/pgr-dev/.devcontainer/Dockerfile` leer bis auf den Kommentar mit dem Projektnamen),
`start.sh` (`diff` leer) und der Compose-Block ohne Mount. Kopf der Compose-Datei von `u-boot init pgr-aus` (0.9.0). Aufräumen: `docker compose down -v` je Projekt, `docker rmi pgr-aus-dev:latest pgr-aus2-dev:latest`;
`docker volume ls -q` vorher und nachher gleich (`diff` leer), `docker ps -a` leer, keine Netze `pgr*`, Image `pgr-dev:latest` unverändert (`d7d697d39d6d`).
Der Host-Eintrag für `known_hosts` entstand mit `ssh-keyscan` (Marke am Zeilenende) und wurde danach entfernt; `~/.ssh/known_hosts` ist wieder byte-gleich der Kopie vor dem Lauf (`cmp`).
Der Host-Klon der Proben liegt im Scratchpad (`host`, Zweig `main` des Worktree-Repos plus ein unveröffentlichter Commit); das echte Host-Repo blieb unberührt.

**Echter Container (§3 Schritt 8).** Beim Start des Laufs lief kein Container (`docker ps -a` leer); Image und Volumes `pgr-dev_*` lagen vor. Ihn zu starten hätte ihn neu erzeugt, was die Festlegung 2 ausschließt.
Keine Probe am echten Container; er, `compose.yaml` des Messorts und seine Volumes sind nicht angefasst (`ls -la /Development/pgr-dev` ohne Änderung). Das Dokument sagt deshalb nichts über den Bestandscontainer (§13).

**Befund an den Randformen aus §6** (Probe · Beobachtung; „Gegenfall“ ist die Form, die das Gegenteil des Satzes im Dokument zeigt):

| Aussage im Dokument (Abschnitt) | Probe | Beobachtung |
|---|---|---|
| `git` und `git-receive-pack` im Image; `/workspaces` gehört `vscode`; Bare-Repo `vscode` 1000, 775 (§6) | `ssh … 'command -v git git-receive-pack; git --version; ls -ld /workspaces'`, `stat -c "%U %u %a"` | `/usr/local/bin/git`, `git-receive-pack`, `git version 2.55.0`, `drwxr-xr-x vscode vscode /workspaces`, `vscode 1000 775 /workspaces/austausch.git` |
| ohne `--initial-branch=main` zeigt `HEAD` auf `master`; Clone danach: Warnung, keine Dateien (§6, §12) | `git init --bare` ohne Option, Push von `main`, Clone | `HEAD` `refs/heads/master`; `warning: remote HEAD refers to nonexistent ref, unable to checkout`, `ls` leer. Gegenfall: `--initial-branch=main` → `ref: refs/heads/main`, Clone checkt aus |
| Clone vor erstem Push: leeres Repo, Exit 0 (§6) | Clone des leeren Bare-Repos | `warning: You appear to have cloned an empty repository.`, Exit 0, Branch `main`, `No commits yet` |
| Wiederholung: `init --bare`, `clone`, `remote add` (§6) | zweiter Aufruf je Befehl | `Reinitialized existing Git repository` (Exit 0); `fatal: destination path … already exists and is not an empty directory` (Exit 128); `externes Repository container existiert bereits` (Exit 3) |
| `git push container main` liefert ein, auch unveröffentlichte Commits (§9) | Push des Scratch-Klons mit einem lokalen Commit (`dc5081bf`) | `* [new branch] main -> main`; `git log` im Bare-Repo zeigt `dc5081bf` an der Spitze |
| `core.sshCommand` mit Schlüssel und `IdentitiesOnly=yes` genügt (§9) | `git config core.sshCommand "ssh -i … -o IdentitiesOnly=yes"`, Push, Fetch | Push und Fetch ohne weitere Optionen erfolgreich (bei vorhandenem Hostschlüssel-Eintrag) |
| `core.sshCommand` gilt für das ganze Repo, nicht für ein Remote (§9) | `GIT_TRACE=1 git ls-remote ssh://git@127.0.0.1:1/x` | `run_command: … 'ssh -i …/id_ed25519 -o IdentitiesOnly=yes' … -p 1 git@127.0.0.1`. `origin` über `ssh`: nicht gelaufen (§13) |
| Clone im Container, `git config --global`, Commit ohne Identität weist ab (§6) | Clone nach dem Push; `git commit --allow-empty` vor der Konfiguration | Clone Exit 0, Branch `main`; `Author identity unknown`. Mit Konfiguration: Commit und Push von `container/thema` in das Bare-Repo |
| Benennung am Host: `container/thema` im Container → `container/container/thema` (§9) | `git fetch container` nach Push von `container/thema` (Gegenfall: Push von `thema2`) | `* [neuer Branch] container/thema -> container/container/thema`; Gegenfall: `thema2 -> container/thema2` |
| zweiter Push des Hosts, Pull im Klon (§6, §9) | Host-Commit `115edfc3`, `git push container main`, im Klon `git switch main; git pull` | `dc5081bf..115edfc3 main -> main`; `Fast-forward`, `h2.txt` im Klon |
| beide schreiben auf `main`: Host-Push weist ab, Erholung (§9, §12) | Container pusht `6d772aad` auf `main`, Host pusht `107d5b11` | `! [rejected] main -> main (fetch first)`, Exit 1; danach `git fetch container`, `git merge container/main`, Push `6d772aad..440c7bd3` erfolgreich, `git pull` im Klon holt `440c7bd3` |
| Push ohne vorhandene Quelle (nicht dokumentiert, Randform) | `git push container thema-gibts-nicht` | `error: Src-Refspec … entspricht keiner Referenz`, Exit 1 |
| Abbruch im Push (§13) | `timeout -s KILL 0.4 git push container main` mit 30-MB-Commit | Spitze von `main` im Bare-Repo unverändert (`440c7bd3`), `git fsck` ohne Ausgabe; Wiederholung `440c7bd3..1b55652d` erfolgreich. Ob der Abbruch mitten in der Übertragung lag: nicht festgestellt (§13) |
| Bare-Repo, Branches, Klon überleben `u-boot down`/`up`; `~/.gitconfig` nicht (§4, §6) | `u-boot down`, `u-boot up`, dann `ls /workspaces`, `git branch`, `status`, `git config --global -l` | Volumes `pgr-aus_*` bleiben; `austausch.git` mit `container/thema`, `main`, `thema2`; `## main...origin/main`; `unable to read config file '/home/vscode/.gitconfig'`. Nach `up` Push des Hosts (`1b55652d..1a8f3a73`) ohne Warnung |
| Hostschlüssel gleich nach `down`/`up` (§4, §5) | `ssh-keygen -lf /etc/ssh/ssh_host_ed25519_key.pub` vor und nach | beide Male `SHA256:6XSFi4LleM3Ae3hAbZHwuftWqvtnG7ueUzBKSK4aJ/U root@buildkitsandbox` |
| Hostschlüssel nach Neubau ohne Cache: andere Schlüssel, Warnung, Abbruch (§5, §12) | `docker compose build --no-cache` (35 s), zweiter Container auf Port 2322 nach `down` des ersten, `git fetch container`, `git push container main` | Fingerabdruck `SHA256:IiuSe6Bc…`; `WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!`, `Host key verification failed.`, Fetch Exit 128, `ssh` Exit 255 |
| ohne Eintrag und ohne Terminal: Abbruch (§5) | `ssh -i … -p 2322 vscode@127.0.0.1 true < /dev/null` vor dem Eintrag | `Host key verification failed.`, Exit 255; mit Terminal (Frage): nicht geprüft |
| zweiter Container, eigener Port, eigenes Remote (§10) | Projekt `pgr-aus2`, Remote `container2`, Bare-Repo im zweiten Container | beide `healthy` gleichzeitig; Push nach `container2` mit dem Schlüssel der ersten Instanz allein: Zugriff abgewiesen; mit beiden `-i`: Push, Fetch beider Remotes, `git branch -r` zeigt `container/…` und `container2/…` |
| Compose-Block ohne Mount: drei Mounts (§3.5, §7, §11) | `docker inspect pgr-aus-dev --format '{{len .Mounts}} …'`; `docker compose config` | `3 volume:/workspaces:true bind:/etc/pgr-dev/authorized_key.pub:false volume:/var/lib/docker:true`; Gegenfall: `docker compose config` mit zusätzlicher Mount-Zeile: 4 |
| `u-boot up` Warnung zum Port (§4) | `u-boot up` im Scratch | `warn: Service "dev" port entry #0 (127.0.0.1:${PGR_SSH_PORT:-2222}:22) is not TCP-probable …`; danach `healthy` nach 43,7 s (Image im Cache) |
| `docker compose up -d` mit `PGR_SSH_PORT` aus der Umgebung (§4) | `PGR_SSH_PORT=2322 docker compose up -d` im Verzeichnis von `pgr-aus2` | Container neu erzeugt, `healthy`, `127.0.0.1:2322->22/tcp` |

**Mutation (Gegenbeispiel) für die Sätze, die eine Verneinung tragen.** Zusage „`/host-repo` und `PGR_REPO` kommen im Dokument nicht mehr vor“: `grep -n 'host-repo\|PGR_REPO' docs/developer/devcontainer.md` vor dem Umbau
21 Zeilen (`git grep -c` auf `92cfaf23`), nach dem Umbau leer (Exit 1). Zusage „der Arbeitsbaum des Hosts ist in keinem der drei Mounts“: `docker compose config` auf dem Block mit einer zusätzlich eingefügten Mount-Zeile (nur Konfiguration, kein Container) zählt vier `type:`-Einträge, der Block des Dokuments drei; `docker inspect` des laufenden Scratch-Containers drei.

**Aussagen aus der Probe des Vorgängers.** Diese Sätze des Dokuments stammen aus `slice-harness-devcontainer-doku` und haben hier keine Zeile in der Tabelle; das Dokument nennt sie „Probe des Vorgängers“ oder „Versuchs-Container“:
Zeiten in §8 (5 s, 201 s, 1 s, 191 s), Platzangaben in §11 (2,05 GB, 5,87 GB, 37,14 MB, 47 MB), `ss -ltn` und `docker port` in §5, die Fehlerzeilen in §12 (`iptables`, `Permission denied (publickey)`, `Permission denied` beim Klonen), `u-boot up` nach 44 s. Dieser Umbau hat sie nicht wiederholt.
Aus diesem Umbau gemessen sind `u-boot up` 43,7 s, `docker compose build --no-cache` 35 s (nicht im Dokument), drei Mounts und alles in der Tabelle.

**Quantoren-Durchsicht** (`grep -nEi 'jede|alle|immer|stets|nie|nichts|kein|jeweils|genau|nur'` über das Dokument; Entscheidung je neuem Treffer): „genau drei Mounts“ (§7): gemessen, unverändert; „Jeder Container hat sein eigenes Bare-Repo“ (§10): im zweiten
Container eigenes Bare-Repo angelegt, im ersten nicht sichtbar, unverändert; §9 „`core.sshCommand` gilt für das ganze Host-Repo, nicht für ein Remote“: eng gefasst auf einen `ls-remote` einer `ssh://`-URL, die kein Remote ist;
§11 Host-Repo: „ein Prozess ohne `sudo` sieht ihn also nicht“, Root über Gerätedateien nicht geprüft (§13 des Dokuments); §4 „Stände, die nur dort lagen, auch“: als gefolgert gekennzeichnet; §6 „checkt keine Dateien aus“: eng gefasst;
§7/§11 Ablageort der Anmeldung von Claude Code: keine Verneinung mehr, nicht geprüft (§13); §9 „im Klon holt `git switch main && git pull` den Stand“: auf den gelaufenen Fall `main` gefasst. Übrige Treffer sind Bestand ohne neuen Satz zum Gegenstand.

**Nacharbeit nach Review und Verifikation** (Commit nach `1e70abc3`): zitierte Meldungen als Locale des Probe-Hosts gekennzeichnet (§6), `$PGR_SSH_PORT`-Sätze in §10 auf den Fall der zweiten Instanz gefasst, Ort der `git remote add container2`-Befehle genannt,
Wortlaut zu `u-boot init` und Messort in §3.5, §4 Fetch → Push, Relativsatz in §1, `u-boot up` 44 s als Messung von `u-boot` benannt. Keine neue Container-Probe; jede Änderung ist eine engere Fassung oder ein Eintrag in §13.

**Gates.** `make docs-check`: `d-check: 991 Datei(en) geprüft, 0 Befund(e)`. `make kopf-check`: Exit 0. `make gates-schnell` auf `9f8faef9`: Exit 0, 1 min 59,8 s (Stufe `lint`: `0 issues.`). Ein voller `make gates` ist nicht gelaufen (§5: nur `docs/developer/` und Planung ändern).

**Größe nach §4.** `git diff --numstat 92cfaf23..9f8faef9 -- docs/developer`: 151 Einfügungen, 100 Löschungen, Summe 251 (Meldeschwelle 350, Marke 450; Schätzung 220 bis 340). Länge des Dokuments: 498 Zeilen vorher, 549 nachher (Grenze 560). Basis `92cfaf23` ist der Elternteil des Anspruchs-Commits `38a9b20a`.

**Hostschlüssel.** Das Dokument nennt, was lief, und empfiehlt keine Einstellung; die Entscheidung bleibt beim Nutzer (§6 Festlegung 3, Dokument §13).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Zählung nach `AGENTS.md` §3.13:** 1 Liefer-Punkt (das umgebaute Dokument), 1 Schicht (Nutzer- und Wartungs-Doku; `docs/developer/` zählt wie
`docs/user/` und `docs/maintainer/` als Dokumentations-Schicht, die Planung zählt nicht). Unter den Grenzen von drei Punkten und zwei Schichten. Der Slice nimmt
keine Sendung eines anderen Slice auf und trägt keine in einen anderen ein; `grep` nach `slice-harness-devcontainer-austausch-repo` in `next/` und `open/` ist außerhalb dieser Datei leer.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist nur die Default-Sub-Area `*` (Kürzel `REPO`, GF) aus `harness/conventions.md`; sie erfüllt die Schwelle als Default für das gesamte Repo.

**Vorgelagert — offene Beobachtungen sichten:** Register unter `../observations/` für `REPO` durchgegangen. Treffer: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Aussage über die Prüfung hinaus;
der Vorgänger trug ein Auftreten; Zähler bei der Closure aus den Dateien lesen), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Zähler 11 laut Vorgänger; Stand *geplant*,
`slice-harness-schaetzung-planner-faktor`; dieser Slice bekommt keine Sendung, der Spielraum der Marke in §4 folgt dem Datenpunkt des Vorgängers), `BEO-REPO/stop-hook-positivliste-ohne-docs-developer`
(Pfad aufgenommen, kein Treffer für diesen Slice). Keine weiteren Treffer.

**Paralleler Lauf** neben `slice-harness-abdeckung-gate`: eigener Worktree und Zweig, Rebase vor dem Fast-Forward-Push, geteilter Docker-Daemon (Vorbedingung §6).

Alle berührten Sub-Areas GF.
