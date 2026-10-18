# Slice slice-harness-devcontainer-einrichten-start: `container-start` startet den Container, prüft Port und Bereitschaft und fängt den Abbruch ab (Teil II von `SPEC-058`; Stopp und Neustart als zweiter Aufruf)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-dev-container (Entscheidung des Nutzers vom 2026-10-18; bis dahin ohne Welle). Angelegt am 2026-10-18 durch den zweiten Schnitt von `slice-harness-devcontainer-einrichten-geruest` (die Schätzung des Architect für das ganze Ziel nach `SPEC-058` lag bei 1000 bis 1450 Zeilen, über der Marke 900; Regel des Nutzers: schneiden, wenn die Schätzung die Marke selbst überschreitet). Der Schnitt liegt an der Naht des Vertrags: Teil II (Punkte 9 bis 12).

**Bezug:** [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität) nur insoweit, als die Host-Voraussetzungen des Repos berührt sind (entschieden mit [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md), im Geber); kein Anspruch des Produkts. [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md) (Accepted 2026-10-18).
Vorgänger in dieser Reihenfolge: `slice-harness-devcontainer-einrichten-geruest` (Geber; liefert Teil I und den Rahmen der Gegenprobe; dieser Slice setzt ihn in `done/` voraus), davor `slice-harness-devcontainer-austausch-repo`, `slice-harness-devcontainer-doku` und `slice-harness-devcontainer-u-boot-0-10` (alle `done/`; letzterer misst Neustart 1,7 bis 1,8 s und Erststart 65 s). Folge-Slices: `slice-harness-devcontainer-einrichten` (setzt den laufenden Container voraus), dann gegebenenfalls `slice-harness-devcontainer-einrichten-zugang`.

**Berührte Spec-Stellen:** [`SPEC-058`](../../../../harness/vertraege/SPEC-058-container-start.md) Teil II (Punkte 9 bis 12) und, soweit sie Teil II betreffen, die Punkte 1, 4, 13, 14, 15, 16 und 17; [`SPEC-052`](../../../../harness/vertraege/SPEC-052-gate-kette.md) (Punkt 8: gleichzeitige Läufe); [`SPEC-057`](../../../../harness/vertraege/SPEC-057-integration-aufraeumen.md) (Punkte 3, 4 und 7: Aufräumen, Fallen); [LH-QA-03](../../../../spec/lastenheft.md#lh-qa-03--portabilität).

**Verantwortlich:** —

**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** `make container-start` startet nach dem Gerüst aus `slice-harness-devcontainer-einrichten-geruest` den Container und meldet ihn bereit: Schritt `zustand` liest Container und Port (Punkte 9 und 10), `start` führt `devcontainer up --workspace-folder . --override-config .devcontainer/override.json` aus oder lässt einen laufenden Container stehen, `warten` probt `ssh` als `vscode` mit `$PGR_DEV/ssh/id_ed25519` auf `127.0.0.1:<PGR_SSH_PORT>` bis `PGR_WARTEN` Sekunden (Punkt 11), `SIGINT` und `SIGTERM` beenden den Lauf mit Schritt in der Meldung (Punkt 12). Stopp ist `docker stop <Name>`, Neustart der zweite Aufruf (Punkt 9). Dazu die Gruppe Docker der Gegenprobe (Container des Präfixes `pgr-cs-`, Aufräumen nach `SPEC-057`) und das Dokument `docs/developer/devcontainer.md` §2 bis §5.

**Annahme (Entscheidung des Nutzers vom 2026-10-17, „ja“): ein Slice pro Container-Sitzung.** Ein Slice läuft in einer Sitzung, danach Neustart (`docker stop`, dann der zweite Aufruf von `container-start`; `u-boot up` startet das gemessene Gerüst nicht, `no service selected`). Die Volumes `<Projektname>-docker-<ID>` und `<Projektname>-workspace-<ID>` bleiben über den Neustart erhalten, der Neustart ist billig (1,7 bis 1,8 s bis `ssh` und `dockerd` antworten, gemessen von `slice-harness-devcontainer-u-boot-0-10`). Das Ziel ist deshalb wiederholbar (Punkt 9) und baut ein vorhandenes Gerüst nicht neu; die Namen der Volumes tragen die `devcontainerId` und ändern sich, wenn das Verzeichnis des Gerüsts wechselt.

**Sendung des Gebers `slice-harness-devcontainer-einrichten-geruest` (eingetragen vom Planner am 2026-10-18, `AGENTS.md` §3.13):** Dieser Slice nimmt an, was der Geber ausschließt (dessen §1, Klasse 1): `SPEC-058` Punkte 9 bis 12, die Teile der Punkte 13, 14, 15 und 17, die Container betreffen (Nachbarpaare 5 und 6 der Reihenfolge, Erfolgsfall und Fehlerfälle von Teil II unter den Aufrufarten, Gruppe Docker mit den Regeln (b) bis (d)), und das Dokument (§2 bis §5). Er erweitert Skript, Fragment und Gegenprobe des Gebers. Er setzt voraus, was der Geber zusagt: die Schritte bis `override` in der Reihenfolge von Punkt 4, die Exit-Codes 2 bis 5, `$PGR_DEV` mit Gerüst, Schlüsselpaar und `mk-override.sh`, die Eingaben mit Vorgaben und Formen (`PGR_NAME`, `PGR_SSH_PORT`, `PGR_WARTEN`, `PGR_CPUS`, `PGR_MEM`) und den Rahmen der Gegenprobe (Attrappen mit Protokoll, Ersetzen genau einer Stelle, Aufräumen bei Ende, `SIGINT`, `SIGTERM`, Präfix `pgr-cs-`). Ändert der Geber einen dieser Punkte, folgt dieser Slice in §1, §3 und §6. Die Entscheidung des Architect zum Zwischenstand nach `override` (Z1 im Geber) ist mit diesem Slice erledigt: Das Ziel endet mit `bereit`.

**Bindung an `slice-harness-devcontainer-einrichten` (Sendung, eingetragen vom Planner am 2026-10-18, `AGENTS.md` §3.13):** Der Nehmer setzt voraus, was dieser Slice zusagt: ein Container, der per `ssh` als `vscode` mit `$PGR_DEV/ssh/id_ed25519` auf `127.0.0.1:${PGR_SSH_PORT:-2222}` antwortet (Punkt 11), mit dem Namen `PGR_NAME` (Vorgabe `pgr-dev`), und ein zweiter Aufruf, der einen laufenden Container stehen lässt (Punkt 9). Dass `dockerd` im Container bereit ist, sagt dieser Slice nicht zu (Punkt 17 (4)); der Nehmer prüft es selbst. Ändert dieser Slice einen dieser Punkte, folgt der Nehmer in §1, §3 und §6.

**Sendung des Gebers `slice-harness-devcontainer-u-boot-0-10` („Arbeit im Container“, eingetragen vom Planner am 2026-10-17, Entscheidung des Nutzers vom selben Tag)** teilt sich auf: Den Weg bis zum laufenden Container tragen `slice-harness-devcontainer-einrichten-geruest` (Gerüst) und dieser Slice (Start), den Weg im Container `slice-harness-devcontainer-einrichten`.

**Kein Gate**, wie im Geber: Das Ziel steht im Abschnitt *Werkzeuge* von `harness/README.md`; die Gegenprobe ist ein Gate in der Ausnahmeliste der schnellen Kette.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Gerüst, kein Schlüsselpaar, kein `mk-override.sh`, keine Prüfung von Eingaben und Voraussetzungen (Punkte 2 bis 8): `slice-harness-devcontainer-einrichten-geruest` (Klasse 1, Geber; dieser Slice setzt dessen Lieferung voraus).
- Kein Bare-Repo, keine Einlieferung, kein Klon, keine Identität, kein Zweig, kein Remote und kein `core.sshCommand` im Host-Repo: `slice-harness-devcontainer-einrichten` (Klasse 1; dessen §1 nimmt die Schritte an).
- Kein Zugang von Claude Code im Container (Quellenauflösung, `devcontainer exec`, Variablen `ANTHROPIC_API_KEY`, `CLAUDE_CODE_OAUTH_TOKEN`, Dokument §7): `slice-harness-devcontainer-einrichten-zugang` (Klasse 1; Entscheidung des Nutzers vom 2026-10-18, der Zugang wird nicht gebraucht). `devcontainer up` löst den Verweis `${localEnv:…}` in `override.json` aus der Umgebung des Hosts auf; ob dabei ein Wert in der Ausgabe erscheint, ist eine Randform (Z2, §6).
- Kein Löschen von Container, Volumes oder Image (Dokument §4, *Container entfernen*, *Alles löschen*, *Image löschen*) und kein Rückbau nach Abbruch (`SPEC-058` Punkte 1 und 12): Handlung des Nutzers; die Gegenprobe entfernt nur Objekte ihres Präfixes per Name (`docker rm -fv`).
- Kein Umbau bei Änderung von `egress.allow`: Der Nutzer ruft `u-boot generate devcontainer` und entfernt den Container von Hand (`SPEC-058` Punkt 9); ein Neubau ändert den Hostschlüssel (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`).
- Keine Entscheidung zum Hostschlüssel über die Probe hinaus: Die zwei Optionen gelten nur für die Probe, das Ziel setzt nie `known_hosts`, `~/.ssh/config` oder `core.sshCommand` (`SPEC-058` Punkt 11, Entscheidung des Nutzers mit [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md)).
- Keine Prüfung, ob `dockerd` im Container bereit ist, und nichts über das, was im Container liegt (`SPEC-058` Punkt 17 (4)): akzeptiertes Negativ; die Bereitschaft ist die Anmeldung als `vscode`.
- Keine Entscheidung zwischen `podman` und `docker` als `nestedRuntime` (`BEO-REPO/nested-runtime-podman-oder-docker-offen`).
- Kein Produktcode (Schicht-Abgrenzung): Es ändern sich Werkzeug, Gegenprobe, Index und Dokument, nichts unter `cmd/` und `internal/`.
- Kein Eintrag in der Roadmap: Der Slice ist wellenlos. Die Release-Slices und die Closure von `welle-v1-abschluss` hängen nicht an ihm (Entscheidung des Nutzers vom 2026-10-17).

## 2. Definition of Done

Liefer-Punkte (drei):

- [ ] **Ziel `container-start`, Teil II** (`tools/harness/container-start.sh`, Fragment des Gebers): erfüllt `SPEC-058` Punkte 9 bis 12 mit den Teilen der Punkte 13 und 14, die Container betreffen (Schritte `zustand` (Container, Port), `start`, `warten`; Exit-Codes 4 und 5 für Teil II, 130 und 143; Meldung `container-start: bereit <Name> 127.0.0.1:<Port>`); `harness/README.md` nennt das Ziel mit dem vollen Verhalten.
- [ ] **Gegenprobe**, Teil II (Erweiterung von `container-start-gegenprobe`): Attrappen von `docker`, `devcontainer` und `ssh` für die Fälle `p9-*`, `p11-*`, `p12-*` und die Fälle von Teil II unter `p13-*` und `p14-*`; Gruppe Docker nach `SPEC-058` Punkt 15 (Container des Präfixes `pgr-cs-`, Port frei und nie 2222, `PGR_CPUS=2`, `PGR_MEM=2g`, Entfernen per Name mit `docker rm -fv` und Prüfung, dass keines der Volumes bleibt, Aufräumen bei Ende, `SIGINT` und `SIGTERM`); je Randform aus §6 ein Fall und je Zusage und je Element einer Aufzählung eine Mutation, die nur ihn rot färbt (`AGENTS.md` §3.10); die Laufzeit der Gegenprobe ist gemessen und im Bericht genannt.
- [ ] **Dokument** `docs/developer/devcontainer.md` §3 bis §5 nennen `make container-start` als Weg und führen die Einzelbefehle nur, soweit das Dokument sie als Erklärung braucht; die Host-Voraussetzungen in §2 sind nachgeführt. Jede Aussage ist durch einen Lauf belegt oder steht in §13. §6 bis §9 gehören `slice-harness-devcontainer-einrichten`, §7 `slice-harness-devcontainer-einrichten-zugang`.

Gates und Closure:

- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `tools/` und `harness/` ändern). Die Regel, nicht gleichzeitig mit dem Lauf eines anderen Slice auf demselben Daemon zu laufen und nie gegen den Container, in dem gearbeitet wird, richtet sich an den Aufrufer (`SPEC-058` Punkt 15).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/vertraege/SPEC-058-container-start.md`, `docs/plan/adr/0045-container-start-als-werkzeug-im-repo.md` | liegen vor (Architect) | Vertrag und Entscheidung; zählen nicht zur Größe |
| `tools/harness/container-start.sh` | update | Schritte `zustand` (Container, Port), `start`, `warten`, Falle für `SIGINT` und `SIGTERM` (Punkte 9 bis 12) |
| `tools/harness/container-start-gegenprobe.sh` | update | Attrappen von `docker`, `devcontainer`, `ssh`; Fälle und Mutationen von Teil II; Gruppe Docker mit Aufräumen |
| `harness/mk/container-start.mk`, `harness/mk/container-start-gegenprobe.mk` | unverändert erwartet | Rezept und Gate-Anbindung stehen im Geber; ändert sich ein Rezept, zählt es mit |
| `harness/README.md` | update | Zeilen *Werkzeuge* und *Sensors* nennen das volle Verhalten und die Fälle gegen reales Docker |
| `docs/developer/devcontainer.md` | update | §2 und §3 bis §5 nennen das Ziel; §13 für Ungeprüftes |

- **Vor dem Code zu messen (Implementer):** das Label `devcontainer.local_folder` (`SPEC-058` Punkt 9): Setzt die `devcontainer`-CLI es und mit welchem Wert (`$PGR_DEV` wörtlich)? Der Implementer misst es an einem Scratch-Container des Präfixes, bevor er den Zustand des Containers schreibt; trägt die Annahme nicht, hält er an und gibt Punkt 9 dem Architect zurück.
- **Docker-Last.** Die Gruppe Docker startet höchstens einen Container des Gerüsts zugleich (Präfix `pgr-cs-`, `SPEC-058` Punkt 15). Eine Vorbedingung „kein `make gates` läuft“ gibt es nicht: Die Gegenprobe läuft selbst in `make gates`, und gleichzeitige Läufe auf einem Daemon sind nicht zugesagt (`SPEC-052` Punkt 8); die Regel richtet sich an den Aufrufer. Ein Erststart baut ein Image (Dokument §4: 65 s bei warmem Cache); die Laufzeit der Gegenprobe ist zu messen (`BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`).
- **Arbeit des Implementers** wie in `slice-harness-devcontainer-einrichten-geruest` §3 (Container `pgr-dev2`, `ssh` und Befehl, Klon, Zweig `slice-harness-devcontainer-einrichten-start`). Die Gruppe Docker legt Container an: Sie läuft dort, wo Docker, `u-boot`, `devcontainer` und `node` vorliegen (Befund des Gebers, Punkt 4 dort), sonst am Host; nie `pgr-dev2` anfassen, nie dessen Port 2223 belegen.
- Zählung nach `AGENTS.md` §3.13 siehe §8.

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-harness-devcontainer-einrichten-geruest` liegt in `done/` (Reihenfolge: Geber, dieser Slice, `slice-harness-devcontainer-einrichten`, dann gegebenenfalls `slice-harness-devcontainer-einrichten-zugang`). Der Slice liegt in `open/`, bis der Nutzer ihn priorisiert und den Anspruch trägt; `Verantwortlich:` setzt dieser Übergang.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt; Vertrag und ADR zählen nicht.

**Schätzung (Planner, 2026-10-18; Metrik wie oben, Gegenprobe mit dem Doppelten des Produktcodes, Faktor aus `AGENTS.md` §3.14; Grundlage: die Schätzung des Architect für das ganze Ziel von 1000 bis 1450, aufgeteilt an der Naht von `SPEC-058`):**

| Teil | Zeilen |
|---|---|
| Skript: `zustand` (Container, Port), `start`, `warten`, Falle | 95 bis 150 |
| Gegenprobe (Attrappen von `docker`, `devcontainer`, `ssh`, Fälle von Teil II, Gruppe Docker mit Aufräumen; mehr als das Doppelte des Skripts, weil die Fälle gegen reales Docker Rahmen brauchen) | 225 bis 330 |
| Index in `harness/README.md` | 0 bis 10 |
| Dokument (§2 bis §5 gekürzt, Löschungen zählen) | 80 bis 130 |
| **Zusammen** | **400 bis 620** |

Meldeschwelle 700, **Marke 900**; die ganze Schätzung liegt unter der Meldeschwelle. Der Geber `slice-harness-devcontainer-einrichten-geruest` schätzt 610 bis 845; beide Teile zusammen liegen bei 1010 bis 1465, im Spektrum der Schätzung des Architect für das ganze Ziel.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 900, oder der Vertrag verlangt mehr als drei Liefer-Punkte oder mehr als zwei Schichten, oder die Randformen aus §6 lassen sich nicht in einer Review-Sitzung prüfen; dann zerlegt der Planner (Vorschlag: Gruppe Docker der Gegenprobe getrennt vom Rest).
- `in-progress` → `open` (blockiert): Das Label `devcontainer.local_folder` trägt die Annahme von Punkt 9 nicht (Rückgabe an den Architect), oder die Probe der Bereitschaft trägt im Container nicht (Befund an den Planner), oder der Geber liefert nicht, was §1 voraussetzt.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün auf dem Endstand, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben. Das Mutationsprotokoll (Zusage · Mutation · roter Test, `AGENTS.md` §3.10) steht im Bericht des Implementers. Die Fälle der Gruppe Docker sind einmal gegen einen echten Scratch-Container grün gesehen, nicht nur gegen Attrappen; die Laufzeit der Gegenprobe ist gemessen.

## 6. Risiken und offene Punkte

Der Slice liefert die zweite Hälfte eines **neuen Vertrags**; §3.12 gilt. Die Randformen stehen hier mit ihrer Fundstelle in `SPEC-058`; sie sind entschieden, bevor der erste Code-Commit steht (Architect). Eine Form, die erst ein Commit mit Code oder Tests hier einträgt, ist nicht genannt. Der Implementer entscheidet keine still, er hält an und gibt sie dem Architect zurück.

**Entschieden nach dem Schnitt (Architect, 2026-10-18):** **Z2** Das Skript entfernt `CLAUDE_CODE_OAUTH_TOKEN`, `ANTHROPIC_API_KEY` und `GIT_TOKEN` zu Beginn aus seiner Umgebung, `devcontainer up` läuft ohne sie; die CLI löst den Verweis `${localEnv:…}` damit zu einem leeren Wert auf, und kein Wert gelangt in die Ausgabe von `up` oder `docker inspect` (`SPEC-058` Punkt 8, Umgebung der Kinder; Punkt 9, Fall `p9-umgebung`). Eine Messung ist nicht nötig. **Z3** Punkt 12 (Abbruch) einschließlich des Falls im Schritt `geruest` und der Wiederholung liegt vollständig in diesem Slice. **Z4** Gruppe Docker: dieser Slice trägt (b) bis (d) und die Regeln für Container, Image und Volumes (`SPEC-058` Punkt 15). **Z1** Mit diesem Slice entfallen der Zwischenstand (Exit 6, Zeile ‚Teil II fehlt‘) und sein Fall `p13-zwischenstand` (`SPEC-058` Punkt 13).

**Randformen von Teil II** (Spalte *Wo entschieden*: Fundstelle in `SPEC-058`; die Randformen von Teil I stehen in `slice-harness-devcontainer-einrichten-geruest` §6):

| # | Randform | Wo entschieden |
|---|---|---|
| A4 | Port `PGR_SSH_PORT` belegt: Verbindungsaufbau zu `127.0.0.1:<Port>` gelingt, Container des Gerüsts läuft nicht (kein Container oder gestoppt) → Exit 4 mit Port; bei laufendem Container keine Prüfung; Port auf anderer Adresse oder ohne Listener wird nicht erkannt (Docker meldet in `start`, Exit 5) | Punkt 10, Punkt 17 (5) |
| A5 | Zweiter Aufruf: Container läuft (kein `up`, kein `override`, weiter zu `warten`); besteht und läuft nicht (`exited`, `created`: `override` und `up`); besteht nicht (`up`, Volumes bleiben) | Punkt 9 |
| A5b | Kollisionen eines vorhandenen Containers: Name `PGR_NAME` ohne das Label (Kollision), Label mit anderem Namen, Label und Name mit veröffentlichtem Port 2222 ungleich `127.0.0.1:<PGR_SSH_PORT>` → je Exit 4, die Meldung nennt Name oder Port; `PGR_CPUS` und `PGR_MEM` des vorhandenen Containers werden nicht verglichen | Punkt 9, Punkt 17 (3) |
| A5c | Label `devcontainer.local_folder` ist das Merkmal des Containers „des Gerüsts“; Annahme vor dem Code zu messen | Punkt 9 (Implementer misst, Rückgabe bei Abweichung) |
| A5d | `devcontainer up` gilt als erfolgreich bei Exit 0 und Ausgabezeile `"outcome":"success"`; Exit 0 ohne diese Zeile und Exit 1 enden mit Exit 5, ohne Rückbau; ein Container im Zustand `created` ist beim nächsten Aufruf „besteht, läuft nicht“; ein ungültiger Wert von `PGR_CPUS` oder `PGR_MEM` endet als Fehler von Docker in `start` | Punkt 9, Punkt 17 (3) |
| A8 | Bereitschaft: `ssh` als `vscode` mit dem Schlüssel des Gerüsts, `true` mit Exit 0; Aufruf mit den genannten Optionen; Probe jede Sekunde bis `PGR_WARTEN`; Erfolg im ersten Versuch, nach mehreren, nie (Exit 5, Meldung nennt `warten` und die letzte Ausgabe von `ssh`, Dauer nicht unter `PGR_WARTEN`); Container bleibt stehen; `dockerd` gehört nicht zur Bereitschaft | Punkt 11, Punkt 17 (4) |
| A9 | Hostschlüssel unbekannt oder geändert: `StrictHostKeyChecking=no` und `UserKnownHostsFile=/dev/null` gelten nur für die Probe; das Ziel schreibt nie `known_hosts`, `~/.ssh/config`, `core.sshCommand` | Punkt 11 (Entscheidung des Nutzers mit [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md)), Punkt 17 (7) |
| A10 | Abbruch mitten im Ablauf: `SIGINT` (130) und `SIGTERM` (143) beenden das Kind (`u-boot`, `devcontainer`, `ssh`), Meldung `container-start: <Schritt>: abgebrochen`, kein Rückbau; Wiederholung setzt fort; akzeptierte Negative (Dateien von `u-boot init` ohne `u-boot.yaml`, Lücke zwischen den Schlüsseldateien); `SIGKILL` fängt nichts auf | Punkt 12 |
| A11 | Reihenfolge der Gründe, Nachbarpaare mit Container: Schlüsselmangel vor Kollision des Containers; Kollision vor belegtem Port; Fehler eines Schrittes (5) vor Signal (130, 143) | Punkt 13 |
| A12 | Stopp und Neustart: `docker stop` dann Aufruf (dieselbe ID, läuft); `docker rm -fv` auf den Container und Aufruf (neue ID, dieselben Volume-Namen); kein zweites Ziel; kein Neubau nach Änderung von `egress.allow` | Punkt 9 (E6) |
| A23 | Der Container wird nie gestoppt, entfernt oder abgebaut (kein `docker stop`, `rm`, `rmi`, `volume rm`, `prune`); bis einschließlich `zustand` nur `ps`, `inspect`, `port`, `info`, `version` | Punkt 1 (Fall `p1-keine-zerstoerung`), Punkt 4 |
| A24 | Gleichzeitige Aufrufe mit demselben `PGR_DEV` sind nicht zugesagt | Punkt 14, Punkt 17 (6) |
| A25 | Gruppe Docker: Präfix `pgr-cs-` für Container, Verzeichnis und Image; Port frei, nie 2222; kein Name ohne Präfix wird angefasst (Abbruch, wenn ein Name aus der Umgebung das Präfix nicht trägt); Entfernen nur per Name mit Prüfung der Volumes, Rest gemeldet und entfernt; Aufräumen bei Ende, `SIGINT`, `SIGTERM`; Mutanten ersetzen genau eine Stelle | Punkt 15 |

**Aufrufarten von `make`** (Randform des Treibers, `AGENTS.md` §3.12 *Aufrufarten des Treibers*; Fundstelle `SPEC-058` Punkt 14): Die Entscheidungen je Flag stehen im Vertrag, die Fälle `p14-*` für Rezept und Teil I im Geber. Dieser Slice ergänzt, was erst mit Container gilt: unter `-n`, `-q`, `-t` kein Aufruf von `docker` und `devcontainer` (Protokoll der Attrappen leer), unter `-i` ein Fehler in `start` oder `warten` ohne Zeile `bereit` (Exit 5, Meldung des Schrittes), unter `-k` und `-j` ein Erfolgsfall von Teil II mit `bereit`; Mutation: ein Rezept mit `$(MAKE)` oder `+` bleibt im Geber.

**Weitere Risiken:**

- **Gegenprobe stört einen anderen Lauf** (Docker-Daemon geteilt, `make gates` eines anderen Slice) oder den Container, in dem der Implementer arbeitet (`pgr-dev2`), oder hinterlässt Container, Images oder Volumes (`BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`). Die Regel „nie gleichzeitig mit dem Lauf eines anderen Slice und nie gegen den Container, in dem gearbeitet wird“ richtet sich an den Aufrufer (`SPEC-058` Punkt 15, `SPEC-052` Punkt 8); die Gegenprobe prüft sie nicht. — **Ausgang (Closure):** entfallen mit Beleg, dass die Gegenprobe nur Objekte ihres Präfixes entfernt und nach Ende, `SIGINT` und `SIGTERM` keines bleibt, oder eingetreten → Folge-Slice mit Kennung.
- **Gegenprobe wächst und dauert** (Image-Bau, Erststart je Lauf; `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`). — **Ausgang (Closure):** entfallen mit gemessener Laufzeit und Mutationsprotokoll, oder weiter offen → Register, wenn die Laufzeit den vollen Lauf dominiert (Re-Evaluierungs-Trigger von [ADR-0045](../../adr/0045-container-start-als-werkzeug-im-repo.md)).
- **Das Label trägt die Annahme von Punkt 9 nicht** (`devcontainer.local_folder` fehlt oder hat eine andere Form). — **Ausgang (Closure):** entfallen, wenn der Implementer das Label vor dem Code misst und der Bericht die Messung nennt; sonst eingetreten → Rückgabe an den Architect, Rückführung nach `open/` (§4).
- **Größe über Schätzung wegen der Gegenprobe** (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`; Faktor eingerechnet, Rückführung §4). — **Ausgang (Closure):** entfallen mit gemessener Summe unter 900, oder eingetreten → Folge-Slice mit Kennung (vom Planner bei Rückführung vergeben).
- **Konflikt an Skript, Gegenprobe, `harness/README.md`** mit dem Geber und mit `slice-harness-devcontainer-einrichten`. — **Ausgang (Closure):** entfallen, weil die Reihenfolge diesen Slice hinter den Geber und den Nehmer hinter diesen Slice setzt (Start-Trigger).
- **Das Dokument veraltet still gegen das Ziel** (kein Sensor; `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). — **Ausgang (Closure):** weiter offen im selben Eintrag, oder eingetreten (Evidence dieses Slice).
- **Umbau des Dokuments verliert Vorbehalte der Vorfassung** (Sicherheitsaussagen aus §11, „Behandle den Container nicht als Sandbox“). — **Ausgang (Closure):** entfallen mit Vergleich der Vorbehalte der Vorfassung mit dem Diff im Review; ein Befund ist *eingetreten* und geht in `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`.
- **Hostschlüssel ändert sich nach Neubau** (A9). — **Ausgang (Closure):** weiter offen → `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung` (Evidence dieses Slice).

## 7. Closure-Notiz

*(wird bei der Closure geschrieben)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`, `BEO-REPO/nested-runtime-podman-oder-docker-offen`, `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe`, `BEO-REPO/integrationstest-hinterlaesst-anonyme-volumes`, `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`; Zählerstände liest der Implementer aus den Dateien in `evidence/` vor dem Start (kein Zähler wird hier gesetzt). Keine Klasse erreicht mit diesem Plan 3× neu ohne Regel; die verkörperten (§3.10 bis §3.14) gelten.

alle berührten Sub-Areas GF (`REPO`: Spezifikation und Entscheidung gehen dem Code voraus, der Vertrag des Architect steht vor dem Code).

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der die Sendung des Gebers einträgt):** Liefer-Punkte 3 (Ziel Teil II, Gegenprobe Teil II, Dokument). Schichten 2 (Harness-Werkzeug mit Gegenprobe und Index in `harness/README.md`; Nutzer- und Wartungs-Doku `docs/developer/`; Planung und Vertrag sind keine Schicht). Schätzung in der Metrik von §4: 400 bis 620, **unter der Marke 900** und unter der Meldeschwelle 700. Der Nehmer liegt nicht über der Größenregel; ein Schnitt unterbleibt. Dieser Slice ist Geber der Bindung an `slice-harness-devcontainer-einrichten`; `grep` nach der Kennung dieses Slice findet ihn im Geber `slice-harness-devcontainer-einrichten-geruest` (§1, §4, §6, §8) und im Nehmer `slice-harness-devcontainer-einrichten` (§1, §4).
