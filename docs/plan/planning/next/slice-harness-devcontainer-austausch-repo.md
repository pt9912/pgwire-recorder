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
(„ja, Austausch-Repo-Umbau einplanen“); eingeordnet in `next/`, aber nicht beansprucht, solange
`slice-harness-integration-volumes` in `in-progress/` liegt (WIP-Limit 1, §4 *Start*).

**Bezug:** kein Anspruch des Produkts und kein Lastenheft-Bezug — der Slice ändert die Beschreibung der
Arbeitsumgebung der Entwickler. Vorgänger: `slice-harness-devcontainer-doku` (`done/`), dessen Kopf
ebenfalls keinen Lastenheft-Bezug führt; der Slice schreibt dessen Dokument um. Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

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

**Zielbild — zu prüfen, nicht zu behaupten** (Vorschlag des Nutzers, 2026-10-17: lokal einen Klon anlegen,
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

- [ ] `docs/developer/devcontainer.md` beschreibt den Austausch über das Bare-Repo nach §1; die Abschnitte
      §3.5 (Block `compose.yaml` ohne Mount und ohne `PGR_REPO`), §6, §7, §9, §10, §11, §12 und §13 sind
      nachgezogen; `/host-repo` und `PGR_REPO` kommen im Dokument nicht mehr vor (`grep`); jede Aussage ist
      durch einen Lauf am Container belegt (Tabelle *Aussage · Probe · Beobachtung* in §7) oder steht in §13.
- [ ] `make docs-check`, `make kopf-check` und `make gates-schnell` grün auf dem Endstand (nur
      `docs/developer/` und die Planung ändern; ein voller `make gates` ist nicht nötig, §5).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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
8. Am echten Container `/Development/pgr-dev` laufen die Schritte 1 bis 5 (Volumes unverändert, Datei `compose.yaml` des Messorts unverändert, Container nicht neu erzeugt);
   das Bare-Repo wird danach entfernt, wenn der Nutzer es nicht behalten will (Frage unten). Schritte 6 und 7 laufen im Scratch-Container.

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

**Start** (`next` → `in-progress`): Nach dem Schließen von `slice-harness-integration-volumes` (WIP-Limit 1 und Docker: die Proben
dieses Slice dürfen dessen Volume-Messungen nicht stören). Der Anspruch ist ein reiner `git mv` auf dem Hauptzweig vor der Arbeit.
Die Stelle in der Reihenfolge ist Empfehlung des Planners; der Nutzer kann sie ändern.

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
  `core.sshCommand` je Repo statt Umgebung (`GIT_SSH_COMMAND`), Wirkung auf andere Remotes (`origin`); Abbruch mitten im Push. — **Ausgang:** bei der Closure (Planner).
- **Das Zielbild trägt nicht ohne Eingriff in `Dockerfile` oder `start.sh`** (zum Beispiel `git` im Image fehlt, Rechte an `/workspaces`). — **Ausgang:** bei der Closure; Eintreten geht an den Nutzer
  (§1: keine Änderung dieser Dateien in diesem Slice), dann Rückführung nach §4.
- **Hostschlüssel-Frage wird berührt** (SSH bei jedem Push und Fetch; nach Neubau des Images andere Hostschlüssel, Warnung am Host). — **Ausgang:** bei der Closure; bleibt die Entscheidung über eine Einstellung
  (`StrictHostKeyChecking`, `UserKnownHostsFile`) offen, → Register `BEO-REPO/…` neu oder bestehender Eintrag, sonst Folge-Slice nach Entscheidung des Nutzers. Das Dokument erfindet keine Einstellung.
- **Docker-Last stört einen anderen Lauf** (Gate-Lauf oder Volume-Messung eines anderen Slice) oder umgekehrt. — **Ausgang:** bei der Closure; Gegenmittel in §3 (Vorbedingung vor jeder Probe).
- **Probe am echten Container ändert dessen Volume** (Bare-Repo, Arbeitsklon bleibt unberührt). — **Ausgang:** bei der Closure; der Zustand vor und nach (`ls /workspaces`) steht in §7, das Bare-Repo wird entfernt oder auf Wunsch des Nutzers behalten.
- **Das Dokument veraltet still gegen den Container** (kein Sensor; wie im Vorgänger) und **Aussagen über die Probe hinaus** (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; die Reviews des Vorgängers fanden sie zweimal).
  — **Ausgang:** bei der Closure; Gegenmittel: Quantoren-Durchsicht (`grep`) wie im Vorgänger, ein Gegenbeispiel je Verneinung.
- **Größe:** Der Umbau liegt über der Marke (Vorgänger 441 gegen Schätzung 150 bis 300). — **Ausgang:** bei der Closure; Gegenmittel §4 (Meldeschwelle 350, Marke 450, 560 Zeilen).

**Offene Fragen an den Nutzer** (keine Entscheidung des Planners):

1. Soll das Bare-Repo, das die Probe am echten Container `/Development/pgr-dev` anlegt, nach der Probe stehen bleiben (der Nutzer arbeitet dann damit weiter) oder entfernt werden?
2. Darf der Implementer den echten Container für Schritt 7 neu erzeugen (`u-boot up` nach Änderung von `compose.yaml` im Messort)? Der Plan sieht das nicht vor (Scratch-Container); mit dem echten Container wäre der Beleg
   *Überleben von `down`/`up`* am Bestand belegt.

## 7. Closure-Notiz

Noch nicht geschlossen. Der Implementer führt hier die Proben (Tabelle *Aussage · Probe · Beobachtung*), die Läufe, die Größe nach §4 und die Quantoren-Durchsicht;
der Planner schreibt bei der Closure Ausgang je Risiko, Lerneintrag und Paarungen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Zählung nach `AGENTS.md` §3.13:** 1 Liefer-Punkt (das umgebaute Dokument), 1 Schicht (Nutzer- und Wartungs-Doku; `docs/developer/` zählt wie
`docs/user/` und `docs/maintainer/` als Dokumentations-Schicht, die Planung zählt nicht). Unter den Grenzen von drei Punkten und zwei Schichten. Der Slice nimmt
keine Sendung eines anderen Slice auf und trägt keine in einen anderen ein; `grep` nach `slice-harness-devcontainer-austausch-repo` in `next/` und `open/` ist außerhalb dieser Datei leer.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist nur die Default-Sub-Area `*` (Kürzel `REPO`, GF) aus `harness/conventions.md`; sie erfüllt die Schwelle als Default für das gesamte Repo.

**Vorgelagert — offene Beobachtungen sichten:** Register unter `../observations/` für `REPO` durchgegangen. Treffer: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (Aussage über die Prüfung hinaus;
der Vorgänger trug ein Auftreten; Zähler bei der Closure aus den Dateien lesen), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (Zähler 11 laut Vorgänger; Stand *geplant*,
`slice-harness-schaetzung-planner-faktor`; dieser Slice bekommt keine Sendung, der Spielraum der Marke in §4 folgt dem Datenpunkt des Vorgängers), `BEO-REPO/stop-hook-positivliste-ohne-docs-developer`
(Pfad aufgenommen, kein Treffer für diesen Slice). Keine weiteren Treffer.

Alle berührten Sub-Areas GF.
