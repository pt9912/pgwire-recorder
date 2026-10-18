# ADR-0045: Start des Dev-Containers als Werkzeug im Repo mit Gegenprobe und eigenen Host-Voraussetzungen

**Status:** Proposed

**Datum:** 2026-10-18

**Autor:** Architect; die Annahme ist Entscheidung des Nutzers und steht aus. Alles unter *Entscheidung* ist bis dahin eine Empfehlung.

**Bezug:** [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes), [ADR-0026](0026-build-und-test-im-multistage-dockerfile.md), [ADR-0038](0038-gate-kette-voll-und-schnell.md), [ADR-0040](0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md)

**Schärft:** [`SPEC-058`](../../../harness/vertraege/SPEC-058-container-start.md)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`docs/developer/devcontainer.md` beschreibt, wie ein Dev-Container mit eigenem Docker-Daemon aus `u-boot` und der `devcontainer`-CLI entsteht: neun Befehle, ein Schlüsselpaar, ein Skript für die Zusätze, dann Start und Warten. Das Dokument ist von Hand zu befolgen und veraltet still gegen das, was funktioniert (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`). Ein Make-Ziel kann den Weg in einem Aufruf gehen und eine Gegenprobe belegen; es schreibt aber Verzeichnisse außerhalb des Repos, startet Container und braucht Werkzeuge, die das Repo bisher nicht voraussetzt (`git`, `docker`, `make`, `AGENTS.md` §3.1): `u-boot`, die `devcontainer`-CLI, `node`, `ssh`, `ssh-keygen`.

## Entscheidung

Empfehlung: **Option A.** Ein Make-Ziel `container-start` im Repo (Fragment unter `harness/mk/`, Skript unter `tools/harness/`), kein Gate, im Abschnitt *Werkzeuge* von `harness/README.md`. Seine Gegenprobe `container-start-gegenprobe` ist ein Gate in der Ausnahmeliste der schnellen Kette ([ADR-0038](0038-gate-kette-voll-und-schnell.md)) und fährt gegen Scratch-Container eines eigenen Namenspräfixes. Die Host-Voraussetzungen `u-boot` ab 0.10.0, `devcontainer`, `node`, `ssh` und `ssh-keygen` gelten nur für dieses Ziel und seine Gegenprobe; kein anderes Ziel und kein anderes Gate braucht sie. Stoppen ist `docker stop`, Neustart der zweite Aufruf; ein zweites Ziel gibt es nicht. Die Einzelregeln stehen in [`SPEC-058`](../../../harness/vertraege/SPEC-058-container-start.md).

Was davon der Nutzer schon entschieden hat: `ssh` als Voraussetzung am Host (2026-10-17). Was nach Empfehlung gilt, ohne dass der Nutzer es entschieden hätte: `u-boot`, `devcontainer`, `node`, `ssh-keygen`, der Ort im Repo, das Fehlen eines zweiten Ziels und die Probe der Bereitschaft mit Optionen, die den Hostschlüssel nur für die Probe außer Kraft setzen (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — Ziel im Repo, Gegenprobe als Gate** | Vertrag, Gegenprobe und Dokument ändern sich zusammen; der Gate-Lauf zeigt einen Bruch; derselbe Weg wie `hooks-install` | fünf Host-Werkzeuge mehr für dieses eine Ziel; die Gegenprobe startet Container und verlängert den vollen Lauf |
| B — Host-Werkzeug außerhalb des Repos | keine neuen Voraussetzungen im Repo | kein Gate, kein Vertrag am Ort der Werkzeuge ([ADR-0040](0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md)); das Dokument verwiese auf etwas, das niemand prüft |
| C — Dokument allein, Handarbeit | nichts zu pflegen | der Weg bleibt ungeprüft, jede Sitzung wiederholt neun Befehle und ihre Fehler |

## Konsequenzen

- Positiv: ein Aufruf bis zum erreichbaren Container; Randformen (vorhandenes Verzeichnis, belegter Port, zweiter Aufruf, Abbruch) sind entschieden und geprüft; der Nehmer des Slice-Plans (Bare-Repo, Klon) setzt auf zugesagte Größen auf.
- Negativ: Die Gegenprobe braucht Docker, `u-boot` und `devcontainer` auf dem Host des Gate-Laufs; ohne sie ist `make gates` dort rot. Die Laufzeit des vollen Laufs wächst um die Zeit des Erststarts eines Gerüsts.
- Akzeptiertes Negativ: Das Ziel läuft nur auf einem Linux-Host zugesagt; macOS und Windows sind nicht gemessen. Grund: kein solcher Host steht bereit, und das Ziel lehnt sie nicht ab.
- Akzeptiertes Negativ: Die Bereitschaft prüft die Anmeldung als `vscode`, nicht `dockerd` im Container. Grund: `devcontainer up` endet nach dem `postStartCommand`, der den Daemon wartet.
- Folgepflicht: `harness/README.md` §Werkzeuge nennt `container-start`, §Sensors die Gegenprobe mit Bindung an diese ADR; streicht der Nutzer eine Voraussetzung, entfällt der Automatismus für das Gerüst, und der Vertrag ist durch eine Folge-ADR neu zu fassen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| bash, Attrappen, Docker | Eingaben, Voraussetzungen, Zustand, Gerüst, Start, Bereitschaft, Abbruch, Reihenfolge der Gründe und Aufrufarten werden an Fällen und Mutationen geprüft | `make container-start-gegenprobe` |

## Re-Evaluierungs-Trigger

Wenn der Nutzer eine der Host-Voraussetzungen streicht, wenn `u-boot` den Start selbst übernimmt, wenn `podman` als `nestedRuntime` gewählt wird, wenn die Gegenprobe den vollen Lauf dominiert oder wenn ein Host außer Linux zur Nutzung gehört.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-18 | Proposed (Empfehlung des Architect) | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
