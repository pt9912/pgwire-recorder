# Slice slice-harness-devcontainer-einrichten: Make-Ziel `container-einrichten` richtet Bare-Repo, Einlieferung und Klon im laufenden Dev-Container ein (`docs/developer/devcontainer.md`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein Abnahmeszenario und
kein Meilenstein hängt an ihm. Angelegt am 2026-10-17 auf Wunsch des Nutzers vom selben Tag.

**Bezug:** [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität) nur insoweit, als die Host-Voraussetzungen des Repos berührt sind (neu: `ssh`, Entscheidung des Nutzers, §1); kein Anspruch des Produkts.
Vorgänger: `slice-harness-devcontainer-austausch-repo` und `slice-harness-devcontainer-doku` (beide `done/`); der Slice automatisiert deren Ablauf. Weiterer Vorgänger seit dem 2026-10-17: `slice-harness-devcontainer-u-boot-0-10` (stellt den Container auf `u-boot` v0.10.0 um und misst `make gates` im Container); dieser Slice baut auf dessen Ergebnis auf (§4). Keine aktive ADR; der Architect legt eine an (§6).

**Berührte Spec-Stellen:** neuer Vertrag eines Harness-Werkzeugs in `harness/vertraege/` (Nummer vergibt der Architect); Host-Voraussetzungen: [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Make-Ziel `container-einrichten` am Host richtet für einen laufenden Dev-Container (`docs/developer/devcontainer.md` §6 und §9) in einem Aufruf ein: Bare-Repo `/workspaces/austausch.git`, Einlieferung von `main`, Arbeitsklon `/workspaces/pgwire-recorder`, Identität des Klons, `make hooks-install` und `make traeger-fetch` im Klon. Arbeit im Dev-Container wird damit der Standardweg, und `main` am Host bleibt für den Stop-Hook frei.

**Sendung von `slice-harness-devcontainer-u-boot-0-10` („Arbeit im Container“, eingetragen vom Planner am 2026-10-17, Entscheidung des Nutzers vom selben Tag):** Lang laufende Slices laufen im Container, damit `main` am Host nicht über Stunden durch einen Implementer-Lauf gesperrt ist (`BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`). Dieser Slice trägt davon den Einrichtungsweg (das Ziel) und im Dokument den Ablauf eines Implementer-Laufs im Container-Klon: Arbeit im Klon, `make gates` dort, Ergebnis per `git fetch container` am Host. Der Geber liefert das Gerüst auf `u-boot` v0.10.0, den Schlüssel von Claude Code als Umgebungsvariable, die gemessene Dauer von `make gates` im Container und die Probe von `claude`; dieser Slice baut auf diesen Befunden auf und misst sie nicht neu. Hat der Geber die Messung nicht grün belegt, startet dieser Slice nicht (§4).

Ablauf des Ziels, in dieser Reihenfolge (Nutzerwunsch 2026-10-17):

1. Bare-Repo anlegen: `git init --bare --initial-branch=main /workspaces/austausch.git` im Container, per `ssh`.
2. `main` einliefern: Remote `container` (`ssh://vscode@127.0.0.1:${PGR_SSH_PORT:-2222}/workspaces/austausch.git`) und `core.sshCommand` im Host-Repo, dann `git push container main`.
3. Klon `/workspaces/pgwire-recorder` aus dem Bare-Repo anlegen.
4. Name und E-Mail des Klons **lokal** aus der Git-Konfiguration des Host-Repos setzen: nie `--global`, nie ein im Dokument oder Skript fest stehender Name; fehlt am Host einer der beiden Werte, bricht das Ziel ab, bevor es im Container etwas ändert (Reihenfolge: Randform in §6); Werte mit Anführungszeichen oder Apostroph kommen unverändert an.
5. `make hooks-install` und `make traeger-fetch` im Klon.

**Entscheidung des Nutzers (2026-10-17): `ssh` ist neue Voraussetzung am Host.** Auf die Frage, ob `ssh` als Voraussetzung für dieses Ziel in Ordnung ist, antwortete der Nutzer „ja“. Die Voraussetzungen des Repos sind bisher `git`, `docker` und `make` (`AGENTS.md` §3.1, [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität)). Der Architect trägt die Entscheidung in die Spezifikation des Vertrags und in die ADR ein (§3); `ssh` gilt nur für dieses Ziel, kein Gate und kein anderes Ziel braucht es.

**Kein Gate**, wie `hooks-install`: Das Ziel steht im Abschnitt *Werkzeuge* von `harness/README.md`, nicht in `GATE_CHECKS`. Seine Gegenprobe ist ein Gate (§2).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Umbau auf die neue Version von `u-boot` (Features `sshd`, `node`, `claude-code`) und keine Messung von `make gates` im Container: Das liefert der Vorgänger `slice-harness-devcontainer-u-boot-0-10` (Klasse 1, Adresse nimmt an: dessen Punkt 1 und 2). `Dockerfile`, `start.sh` und Compose gelten, wie `docs/developer/devcontainer.md` sie nach jenem Slice nennt; reicht der Inhalt für den Ablauf nicht, hält der Implementer an (§6).
- Keine Entscheidung zwischen `podman` und `docker` als `nestedRuntime` (offene Frage im Dokument, `slice-harness-devcontainer-u-boot-0-10`): Der Ablauf des Ziels setzt den Container voraus, wie er nach jenem Slice läuft.
- Kein Anlegen, Starten oder Neubau des Containers: Das Ziel setzt einen laufenden Container voraus (§3 des Dokuments, `u-boot up`) und fragt nur, ob er per `ssh` erreichbar ist. Container zu erzeugen wäre eine Entscheidung über Image, Volumes und Schlüssel, die nicht in ein Einrichtungsziel gehört.
- Keine Entscheidung zum Hostschlüssel (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`): Das Ziel nennt, was `ssh` ausgibt, und setzt keine Einstellung (`StrictHostKeyChecking`, `UserKnownHostsFile`) ohne Entscheidung des Nutzers (§6).
- Keine Docker-Daemon-Probe im Ziel selbst: Es prüft den Container per `ssh`, nicht per `docker`. Die Gegenprobe darf den Daemon nur nutzen, wenn kein `make gates` läuft (§3, *Docker-Last*).
- Kein Produktcode (Schicht-Abgrenzung): Es ändern sich Werkzeug, Gegenprobe und Dokumentation, nichts unter `cmd/` und `internal/`.
- Kein Eintrag in der Roadmap: Der Slice ist wellenlos.

## 2. Definition of Done

Liefer-Punkte (drei):

- [ ] **Ziel `container-einrichten`** (Fragment unter `harness/mk/` und Skript unter `tools/harness/`): erfüllt den Vertrag, den der Architect in `harness/vertraege/` festlegt (Ablauf aus §1, Randformen aus §6); kein Gate; in `harness/README.md` §Werkzeuge genannt.
- [ ] **Gegenprobe** `container-einrichten-gegenprobe` (Gate, hängt an `GATE_CHECKS`, steht in der Ausnahmeliste der schnellen Kette): gegen einen Scratch-Container (Präfix, Entfernen nur per Name, `docker rm -fv`, nichts Fremdes) je Randform aus §6 ein Fall und je Zusage und je Element einer Aufzählung eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die Gegenprobe steht im Gate-Index.
- [ ] **Dokument** `docs/developer/devcontainer.md` §6 und §9 nennen `make container-einrichten` als Weg und führen die Einzelbefehle nur, soweit das Dokument sie als Erklärung braucht; dazu der Ablauf eines Implementer-Laufs im Container-Klon (Sendung von `slice-harness-devcontainer-u-boot-0-10`, §1): Arbeit im Klon, `make gates` dort mit der Dauer, die jener Slice misst, Ergebnis per `git fetch container`. Jede Aussage ist durch einen Lauf belegt oder steht in §13.

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
| `harness/vertraege/SPEC-<NNN>-container-einrichten.md` | neu (Architect) | Vertrag des Werkzeugs: Ablauf, Randformen aus §6, Host-Voraussetzung `ssh`; Nummer vergibt der Architect |
| `docs/plan/adr/` (eine ADR) | neu (Architect) | Entscheidung: Ziel als Werkzeug, Host-Voraussetzung `ssh` (Entscheidung des Nutzers vom 2026-10-17); `Schärft:` zeigt auf den Vertrag (`AGENTS.md` §3.8) |
| `harness/mk/container-einrichten.mk` | neu | Ziel `container-einrichten` (kein Gate), Gegenprobe hängt an `GATE_CHECKS` |
| `tools/harness/container-einrichten.sh` | neu | Ablauf aus §1; Eingaben aus der Umgebung (`PGR_DEV`, `PGR_SSH_PORT`, Remote-Name), Fehlermeldungen mit Schritt |
| `tools/harness/container-einrichten-gegenprobe.sh`, `harness/mk/container-einrichten-gegenprobe.mk` | neu | Gegenprobe gegen Scratch-Container; Fälle und Mutationen aus §6 |
| `harness/mk/gates-schnell.mk` | update | Gegenprobe in die Ausnahmeliste der schnellen Kette (Docker; Pflege nach `SPEC-052`, Architect prüft Punkt der Gate-Kette; die Datei ändert zurzeit ein anderer Lauf, §6) |
| `harness/README.md` | update | Abschnitt *Werkzeuge* (`container-einrichten`) und *Sensors* (Gegenprobe) |
| `docs/developer/devcontainer.md` | update | §6 und §9 nennen das Ziel; §13 für Ungeprüftes |

- Das Skript trennt Prüfen und Handeln: erst alle Vorbedingungen am Host (Identität, `ssh` vorhanden, Schlüssel lesbar, Container per `ssh` erreichbar), dann die Schritte im Container. Die Reihenfolge der Gründe, aus denen es abbricht, legt der Architect fest (§6).
- Die Identität reist ohne Auswertung durch eine zweite Shell, soweit der Architect keine andere Form festlegt: Der Implementer entscheidet die Quoting-Form nicht, sondern setzt die des Vertrags um (§6, Sonderzeichen).
- **Docker-Last.** Die Gegenprobe startet Scratch-Container (Präfix, Entfernen nur per Name, `docker rm -fv` und Prüfung, dass keines der Volumes bleibt, wie `SPEC-057`). Sie läuft nie gleichzeitig mit `make gates`; die Vorbedingung `pgrep -fa "^make gates"` leer ist Teil der Gegenprobe (ob als Abbruch oder als Warten, legt der Vertrag fest). Ihre Laufzeit ist zu messen (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`).
- Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-17): **3 Liefer-Punkte** (Ziel, Gegenprobe, Dokument). **Schichten: 2**, nämlich Harness-Werkzeug (`tools/`, `harness/mk/`; Gegenprobe und Gate-Index gehören dazu, kein Produkt-Code) und Nutzer- und Wartungs-Doku (`docs/developer/`); `harness/README.md` ist Index desselben Werkzeugs, Planung und Vertrag sind keine Schicht. Nehmer von Sendungen sind nicht eingetragen: Dieser Slice hat keine Folge-Slices, `grep` nach seiner Kennung in den anderen Plänen findet ihn nur in `slice-harness-devcontainer-u-boot-0-10` (Vorgänger) und in der Roadmap. **Nachzählung der Sendung von `slice-harness-devcontainer-u-boot-0-10` (Planner, 2026-10-17):** siehe §8; die Sendung liegt im Dokument-Punkt (Liefer-Punkt 3), kein vierter Punkt und keine dritte Schicht.

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-harness-devcontainer-u-boot-0-10` liegt in `done/`, und seine Messung von `make gates` im Container ist grün (Entscheidung des Nutzers vom 2026-10-17: zuerst der Dev-Container, dann die lang laufenden Slices darin; Reihenfolge: nach `slice-harness-coverage`, hinter `slice-harness-devcontainer-u-boot-0-10`; die Release-Slices und die Closure von `welle-v1-abschluss` verlangen nur jenen in `done/`, diesen Slice nicht, Entscheidung des Nutzers vom 2026-10-17). Bis dahin bleibt der Slice in `open/`, weil der Vorgänger noch nichts liefert, worauf er baut; er geht mit der Priorisierung nach `next/` (`Verantwortlich:` setzt dieser Übergang), sobald jener Slice schließt. Vorbedingung des Starts außerdem: Vertrag und ADR des Architect liegen vor (`Accepted` ist die Entscheidung des Nutzers, keine Annahme des Planners).

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; der Vertrag und die ADR (Commits des Architect vor dem Start) liegen vor der Basis und zählen nicht. Schätzung (Planner, 2026-10-17):
Skript rund 90 bis 130 Zeilen Produktcode, Fragmente rund 50, Gegenprobe mit dem Doppelten des Produktcodes (Faktor aus `AGENTS.md` §3.14; Fall- und Mutationstabelle, Attrappe oder Scratch-Container-Rahmen) rund 300 bis 450, Dokument und Index rund 40 bis 70, dazu der Ablauf eines Implementer-Laufs im Container-Klon (Sendung von `slice-harness-devcontainer-u-boot-0-10`) rund 30 bis 60: **510 bis 760**. Der Architect bestätigt die Schätzung vor dem Code. Meldeschwelle 700, **Marke 900**; die obere Hälfte der Schätzung liegt über der Meldeschwelle, unter der Marke, ein Schnitt unterbleibt (Praxis vom 2026-10-16: geschnitten wird über der Marke).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder der Vertrag verlangt mehr als drei Liefer-Punkte oder mehr als zwei Schichten (zum Beispiel eine Änderung an `Dockerfile` oder `start.sh`), oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner (Vorschlag: Ziel mit Dokument zuerst, die Gegenprobe der Randformen als Folge-Slice mit eigener Kennung).
- `in-progress` → `open` (blockiert): Der Ablauf trägt nicht ohne Eingriff in das Image (zum Beispiel fehlt `git` im Container, Rechte an `/workspaces`), oder `ssh` am Host ist nicht zugesagt (Entscheidung des Nutzers wird zurückgenommen), oder das Gerüst, das `slice-harness-devcontainer-u-boot-0-10` hinterlässt, trägt den Ablauf nicht (der Befund geht an den Planner).

## 5. Closure-Trigger

DoD vollständig, `make gates` grün auf dem Endstand, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben. Das Mutationsprotokoll (Zusage · Mutation · roter Test, `AGENTS.md` §3.10) steht im Bericht des Implementers. Die Gegenprobe ist einmal gegen einen echten Scratch-Container grün gesehen, nicht nur gegen Attrappen.

## 6. Risiken und offene Punkte

Der Slice liefert einen **neuen Vertrag** (Werkzeug mit Zusagen); §3.12 gilt. Die Randformen stehen hier, entschieden werden sie **vor dem ersten Code-Commit** in `harness/vertraege/` (Architect); eine Form, die erst ein Commit mit Code oder Tests hier einträgt, ist nicht genannt. Der Implementer entscheidet keine still, er hält an und gibt sie dem Architect zurück. Vorläufiges Wissen aus dem Dokument steht dabei, nicht als Entscheidung.

**Randformen des Ablaufs** (Liste für den Architect; Spalte *Wo entschieden*: Fundstelle im Vertrag, vom Architect einzutragen):

| # | Randform | Vorwissen aus `docs/developer/devcontainer.md` (§6, §9) | Wo entschieden |
|---|---|---|---|
| 1 | Bare-Repo besteht schon | `git init --bare` auf vorhandenem Verzeichnis meldet `Reinitialized`, Exit 0; Frage: Wiederholung ist erlaubt oder bricht ab? Besteht es mit anderem `HEAD` (nicht `main`)? | offen |
| 2 | Klon besteht schon (leer oder nicht leer, anderer Remote) | `git clone` in nicht leeres Verzeichnis: Exit 128. Frage: überspringen, aktualisieren oder abbrechen; was geschieht mit uncommittetem Stand im Klon? | offen |
| 3 | Remote `container` besteht schon (gleiche URL, andere URL, anderer Port) | zweites `git remote add` bricht mit Exit 3 ab. Frage: bei gleicher URL weiter, bei anderer abbrechen oder setzen | offen |
| 4 | Container läuft nicht / `ssh` nicht erreichbar (Port, Schlüssel, Zeitgrenze) | — Frage: Meldung, Exit-Code, Zeitgrenze für die Verbindung (`ConnectTimeout`), vor jeder Änderung | offen |
| 5 | Identität am Host leer (nur Name, nur E-Mail, beide, nur Leerzeichen) | der Aufruf in §6 bricht ab, wenn eine fehlt. Frage: Reihenfolge gegen Randform 4 (welcher Grund zuerst), Meldung nennt den fehlenden Wert | offen |
| 6 | Identität mit Sonderzeichen: Anführungszeichen `"`, Apostroph `'`, Leerzeichen, `$`, Backtick, Backslash, Zeilenumbruch, Nicht-UTF-8 | `printf %q` am Host; Gegenprobe am Wert nötig. Frage: erlaubte Menge, Quoting-Form (die Shell im Container ist die Anmelde-Shell von `vscode`), ob Zeilenumbruch abgelehnt wird | offen |
| 7 | Identität nur lokal, nie global | Zusage: `--global` und `--system` werden nie gesetzt; Frage: gilt auch, wenn das Host-Repo die Werte nur aus `--global` bezieht (`git config user.name` liest alle Ebenen) | offen |
| 8 | Anderer Port, mehrere Instanzen (`PGR_SSH_PORT`, zweiter Remote-Name) | Dokument §10: je Instanz eigener Port und Remote (`container2`). Frage: wie wählt das Ziel Remote-Name und Bare-Repo-Pfad je Instanz; `PGR_SSH_PORT` leer, nicht numerisch, außerhalb 1 bis 65535; Instanz im Pfad | offen |
| 9 | Hostschlüssel unbekannt, geändert (Neubau) | Warnung am Host bekannt; keine Einstellung festgelegt (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`). Frage: Ziel setzt nichts und bricht ab, oder die Entscheidung des Nutzers liegt vor | offen (Entscheidung des Nutzers) |
| 10 | Netz fehlt bei `make traeger-fetch` im Klon | braucht Netz im Container; das Egress des Containers bleibt eingeschränkt (Entscheidung des Nutzers vom 2026-10-17), der Host des Klons muss in der `egress.allow`-Liste stehen, die `slice-harness-devcontainer-u-boot-0-10` misst und einträgt; Frage: Exit-Code des Ziels, ob die vorigen Schritte stehen bleiben (ja, kein Rückbau) und wie die Wiederholung aussieht | offen |
| 11 | Abbruch mitten im Ablauf (`SIGINT`, `SIGTERM`, Verbindungsabbruch) und Wiederholung | kein Rückbau; Frage: jeder Schritt wiederholbar (Randformen 1 bis 3), Meldung nennt den letzten Schritt | offen |
| 12 | `main` am Host fehlt, ist leer, steht nicht im Host-Repo oder nur auf anderem Zweig; Host-Repo im Zustand `detached HEAD` | Frage: einzuliefernder Zweig ist `main`, das Ziel hängt nicht am ausgecheckten Zweig; Meldung bei fehlendem `main` | offen |
| 13 | Push abgewiesen (Bare-Repo hat schon einen `main` mit anderer Historie) | `! [rejected]` bei nicht-fast-forward. Frage: kein `--force`, Meldung und Exit | offen |
| 14 | Ausgang mit Exit-Codes: Reihenfolge der Gründe, wenn mehrere zutreffen (Identität fehlt **und** Container aus) | Reihenfolge der Gründe bindet der Vertrag (`AGENTS.md` §3.12, Stelle einer Nachricht) | offen |
| 15 | Zugangsschlüssel von Claude Code im Container | Entscheidung des Nutzers vom 2026-10-17: Der Container muss den Schlüssel nicht besser schützen als der Host, wo die Datei jeder Benutzer des Hosts lesen kann; kein Sicherheitsargument gegen das Durchreichen. Entscheidung des Nutzers vom 2026-10-17 dazu: Das Token erneuert der Host, nicht der Container; danach wird der Container mit dem neuen Schlüssel neu gestartet. Der Container schreibt den Schlüssel nie, ein schreibgeschütztes Durchreichen genügt. Beide Wege müssen im Container funktionieren (Entscheidung des Nutzers vom 2026-10-17): API-Schlüssel und Abo-Anmeldung (Abo-Token, `CLAUDE_CODE_OAUTH_TOKEN` oder was `claude` für das Abo braucht); wie das Abo-Token durchgereicht wird, misst der Vorgänger, es ist nicht entschieden. Der Zugang kommt als Umgebungsvariable, ein Mount der Datei entfällt (Entscheidung des Nutzers vom 2026-10-17: „da reicht vielleicht die Env-Var“). Seit dem 2026-10-17 beantwortet der Vorgänger `slice-harness-devcontainer-u-boot-0-10` die Probe (dort Randformen 1 bis 4 und 7: woher die Variable kommt, welcher Name, was `claude` bei fehlendem oder abgelaufenem Token tut); dieser Slice übernimmt das Ergebnis aus dessen Dokument und entscheidet es nicht neu. Das erste Bereitstellen bleibt Handarbeit | Architect, Spezifikation des Werkzeugs |

**Aufrufarten von `make`** (Randform des Treibers, `AGENTS.md` §3.12 *Aufrufarten des Treibers*): Das Ziel ändert Zustand außerhalb des Arbeitsbaums (Container, Remote, Konfiguration des Host-Repos). Je Flag eine Entscheidung im Vertrag und je Flag ein Fall und eine Mutation in der Gegenprobe: `-n` (Rezept mit `$(MAKE)` darf nichts ausführen: kein `ssh`, kein Push, keine Konfiguration), `-q`, `-t` (Ziel ist `.PHONY`, trägt kein Zeitstempel-Ziel), `-i` (Fehler dürfen den Ablauf nicht überlesen: kein „Einrichtung abgeschlossen“ nach einem Fehler), `MAKEFLAGS` aus der Umgebung, auch in Kombination und mit Vorrang. Dazu `-k` und `-j` (parallel laufende Schritte dürfen die Reihenfolge nicht brechen).

**Weitere Risiken:**

- **Neue Host-Abhängigkeit `ssh`.** Entscheidung des Nutzers vom 2026-10-17 („ja“); Eintrag in Vertrag und ADR durch den Architect, kein Plan-Eintrag in `AGENTS.md` §3.1 und nicht in Lastenheft oder Gate-Index, solange der Architect nichts anderes festlegt. Fehlt `ssh` am Host, meldet das Ziel es vor jeder Änderung (Randform 4 in der Form „`ssh` nicht gefunden“, eigener Fall). — **Ausgang (Closure):** offen bis zur Closure; erwartet *entfallen* mit Verweis auf Vertrag und ADR.
- **Der volle Lauf von `make gates` im Container ist ungemessen** (`docs/developer/devcontainer.md` §13). Das Ziel macht die Arbeit im Container zum Standardweg; läuft dort `make gates` nicht (Docker-in-Docker, Laufzeit, Speicher, Volumes), bleibt der Standardweg ohne Gate-Nachweis und der Stop-Hook am Host blockiert dennoch. Seit dem 2026-10-17 misst `slice-harness-devcontainer-u-boot-0-10` den Lauf (dessen DoD Punkt 2, Eintrag mit der Kennung des Gebers dieses Risikos); dieser Slice misst ihn nicht neu und schreibt kein Ergebnis nach, das jener Slice nicht belegt. — **Ausgang (Closure):** entfallen, weil `slice-harness-devcontainer-u-boot-0-10` den Lauf misst und der Start dieses Slice seine grüne Messung verlangt (§4); trägt die Messung des Vorgängers nicht, ist der Slice nicht gestartet.
- **Gegenprobe stört einen anderen Lauf** (Docker-Daemon geteilt, `make gates` des Implementers eines anderen Slice) oder hinterlässt Volumes (`BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`). — **Ausgang (Closure):** entfallen mit Beleg, dass die Gegenprobe nur Objekte ihres Präfixes entfernt und nach Ende keines bleibt, oder eingetreten → Folge-Slice mit Kennung.
- **Gegenprobe wächst** (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`): zwei Zusagen je Aufzählungselement. — **Ausgang (Closure):** entfallen mit gemessener Laufzeit und Mutationsprotokoll.
- **Größe** über Schätzung wegen der Gegenprobe (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Faktor 2 eingerechnet, Rückführung §4). — **Ausgang (Closure):** entfallen mit gemessener Summe, oder eingetreten → `slice-harness-devcontainer-einrichten-gegenprobe` (Name vorläufig, vom Planner bei Rückführung vergeben).
- **Konflikt mit dem laufenden Lauf an `harness/mk/gates-schnell.mk`** (Stand 2026-10-17 ändert ein Implementer die Datei für `slice-harness-coverage`). — **Ausgang (Closure):** entfallen, wenn der Start nach dem Schließen von `slice-harness-coverage` liegt (der Vorgänger `slice-harness-devcontainer-u-boot-0-10` startet erst dahinter, Start-Trigger §4); sonst Abstimmung der Reihenfolge durch den Planner.
- **Die Sendung „Arbeit im Container“ wächst den Dokument-Punkt** (`BEO-REPO/slice-waechst-durch-uebernahmen`). — **Ausgang (Closure):** entfallen mit gemessener Summe nach der Metrik aus §4 (unter der Marke 900), oder eingetreten → `slice-harness-devcontainer-einrichten-gegenprobe` (Name vorläufig, Schnitt wie im Risiko „Größe“).
- **Das Dokument veraltet still gegen das Ziel** (kein Sensor) und Aussagen über die Probe hinaus (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). — **Ausgang (Closure):** weiter offen im selben Eintrag, oder eingetreten (Evidence dieses Slice).
- **Umbau des Dokuments verliert Vorbehalte der Vorfassung** (Sicherheits- und Isolationsaussagen aus §11, „Behandle den Container nicht als Sandbox“). — **Ausgang (Closure):** entfallen mit Vergleich der Vorbehalte der Vorfassung mit dem Diff im Review; ein Befund ist *eingetreten* und geht in `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.
- **Docker-Daemon-Probe nur ohne laufendes `make gates`** (Abgrenzung §1). — **Ausgang (Closure):** entfallen, wenn die Vorbedingung `pgrep -fa "^make gates"` vor jeder Probe leer war und das im Bericht steht.

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie ist so grob wie das Repo, trägt aber eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-17). Treffer für Werkzeug und Container: `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (Randform 9), `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Zählerstände liest der Implementer aus den Dateien in `evidence/` vor dem Start (kein Zähler wird hier gesetzt). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel; die verkörperten (§3.10 bis §3.14) gelten.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Nachzählung der Sendung von `slice-harness-devcontainer-u-boot-0-10` (Planner, 2026-10-17, im Commit, der die Sendung einträgt; `AGENTS.md` §3.13):** Liefer-Punkte 3 (Ziel, Gegenprobe, Dokument; die Sendung liegt im Dokument-Punkt, §2), Schichten 2 (Harness-Werkzeug, Nutzer- und Wartungs-Doku; das Gerüst des Vorgängers liegt außerhalb des Repos und ist keine Schicht dieses Slice). Schätzung in der Metrik von §4 von 480 bis 700 auf 510 bis 760; Meldeschwelle 700, Marke 900. Der Nehmer liegt nicht über der Größenregel; ein Schnitt unterbleibt, die Rückführung in §4 bleibt.
