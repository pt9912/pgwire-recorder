# Slice slice-harness-devcontainer-doku: Dokument zum Arbeiten im Dev-Container (`docs/developer/devcontainer.md`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice; kein Abnahmeszenario und
kein Meilenstein hängt an ihm. Angelegt am 2026-10-17 auf Wunsch des Nutzers, ohne Voraussetzung
eines Schritts von welle-v1-abschluss oder welle-erster-release; er steht nach
`slice-v1-abschluss-container-doku` und vor den wellenlosen Slices vor dem Release oder
parallel dazu.

**Bezug:** kein Anspruch des Produkts — der Slice dokumentiert die Arbeitsumgebung der
Entwickler. Keine aktive ADR.

**Berührte Spec-Stellen:** —

**Verantwortlich:** pt9912

**Autor:** pt9912 (Entwurf: Planner). **Datum:** 2026-10-17.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Dokument `docs/developer/devcontainer.md` beschreibt, wie man mit einem
Dev-Container arbeitet, den das Werkzeug `u-boot` (`pt9912/u-boot`) erzeugt, vom Einrichten bis
zum Aufräumen, mit den Dateiinhalten, die außerhalb des Repos liegen.

Inhalt des Dokuments, je Punkt nur mit Probe am laufenden Container (§3.11):

- Voraussetzungen: Docker ≥ 24, Compose ≥ 2.20, `u-boot`.
- Einrichten: `u-boot init <name> --devcontainer --no-git`.
- Manuelle Ergänzungen **mit den Dateiinhalten**: `Dockerfile` (Docker Engine, sshd, Node, Claude
  Code), `start.sh` (`dockerd`, `authorized_keys`, `chown /workspaces`, `sshd -D`), Service `dev`
  in `compose.yaml` (`privileged`, `cpus`, `mem_limit`, Port `127.0.0.1:2222:22`, Volumes,
  Healthcheck).
- Schlüsselpaar erzeugen; Start, Stopp und Löschen (`u-boot up/down`, `docker compose`);
  Verbinden per SSH; Klon im Container.
- Claude Code starten und selbst anmelden (keine Zugangsdaten kopieren).
- Arbeiten mit `make` und den Gates im Container (eigener Docker-Daemon); gemessene Zeiten:
  `docs-check` 5 s, `gates-schnell` 201 s kalt; der volle `make gates` ist im Container nicht
  gemessen und steht als solcher da.
- Ergebnisse zurück ins Host-Repo: Branch `container/<thema>`, `git fetch` über SSH oder
  `git bundle`.
- Grenzen und Sicherheit (privilegiert, nur 127.0.0.1, nur Schlüssel, Ressourcen) und Aufräumen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Kein Umbau des Repos und kein Eingriff in Makefile, Gates oder Dockerfile des Produkts — der
  Slice beschreibt eine Umgebung um das Repo, nicht das Repo (Schicht-Abgrenzung).
- Kein `.devcontainer/` im Repo: Die Dateien stehen im Dokument, weil sie zum Versuch außerhalb
  des Repos gehören und ein Verzeichnis im Repo ein Gate und eine Pflege verlangte, die niemand
  zugesagt hat (anderer Vorgang).
- Keine Aussagen zu Windows- oder macOS-Hosts — nicht geprüft, und §3.11 erlaubt nur Belegtes.
- Kein Schritt, den der Implementer am laufenden Container nicht nachgeführt hat: Er ist als
  „nicht geprüft“ markiert oder fehlt (§3.11).
- Keine Änderung der Positivliste des Stop-Hooks um `docs/developer/` — das ist eine
  Regeländerung und Entscheidung des Nutzers (§6).

## 2. Definition of Done

- [ ] `docs/developer/devcontainer.md` liegt vor und deckt die Punkte aus §1 ab; jeder Schritt ist
      am laufenden Container nachgeführt oder als „nicht geprüft“ markiert, die Dateiinhalte
      stimmen mit denen des Versuchs überein.
- [ ] `make kopf-check` und `make docs-check` grün; `make gates` grün (Stop-Hook, §6).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder in §7 „keine Beobachtung“ notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/developer/devcontainer.md` | neu | Liefer-Punkt: das Dokument nach §1 |

- Der Implementer führt die Schritte am laufenden Container unter `/Development/pgr-dev` nach
  (nur lesend und per SSH; den Container nicht verändern) und übernimmt die Dateiinhalte von dort.
- Das Dokument ist Ist-Zustand ohne Chronik (§3.11); es nennt keine Kennungen von Lastenheft,
  Spezifikation, ADRs, Slices oder Wellen. Nennt es doch eine ADR-Kennung, ist sie ein Link
  (`.d-check.yml`: `ids` mit `link-policy: always`).
- `.d-check.yml` schließt `docs/developer/` nicht aus (`scan.roots: ["."]`, `ignore` ohne den
  Pfad): Links und Anker im Dokument werden geprüft. `harness/README.md` braucht keinen Eintrag
  (kein neues Target, keine Source-Precedence-Zeile).

## 4. Trigger

**Start** (`next` → `in-progress`): ohne Voraussetzung. Der Anspruch ist ein reiner `git mv` auf
dem Hauptzweig vor der Arbeit.

**Größe.** Metrik nach `AGENTS.md` §3.14: Summe aus Einfügungen und Löschungen von
`git diff --numstat <Basis>..<letzter Commit des Slice> -- docs/developer`. Basis ist der
Elternteil des Commits, der diese Datei nach `in-progress/` bewegt. Schätzung 150 bis 300,
Meldeschwelle 400, Marke 500.

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): Die Summe nach der Metrik liegt über 500, oder die Schritte
  verlangen mehr als dieses eine Dokument.
- `in-progress` → `open` (blockiert): Der Container unter `/Development/pgr-dev` läuft nicht
  mehr und die Schritte lassen sich nicht nachführen, ohne ihn neu zu bauen.

## 5. Closure-Trigger

DoD vollständig, `make kopf-check`, `make docs-check` und `make gates` grün, Review- und
Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Stop-Hook und Positivliste.** `docs/developer/` steht nicht auf der Positivliste des
  Stop-Hooks (`docs/plan/planning/`, `docs/reviews/`, `docs/user/`, Wurzel-`README.md`); der Hook
  verlangt für diesen Slice einen vollen `make gates` vor dem Push. Dieser läuft ohnehin mit dem
  nächsten Code-Slice; die Positivliste ändert dieser Slice nicht (Regeländerung, Entscheidung
  des Nutzers). — **Ausgang:** weiter offen: → Entscheidung des Nutzers über `docs/developer/`
  auf der Positivliste; bis dahin trägt der volle Lauf.
- **Das Dokument beschreibt einen Versuch außerhalb des Repos.** Wird der Container umgebaut,
  veraltet das Dokument still; kein Sensor prüft es gegen den Container. — **Ausgang:** weiter
  offen: → BEO-REPO/zusage-im-kommentar-weiter-als-pruefung im Register, falls der Fall bei der
  Closure auftritt; sonst entfallen mit Begründung.
- **Nicht nachführbare Schritte** (etwa der volle `make gates` im Container, Zeit nicht
  gemessen). — **Ausgang:** entfallen: Das Dokument markiert sie als „nicht geprüft“ (§1).
- **Randformen:** kein neuer Vertrag, keine Randformen im Sinn von `AGENTS.md` §3.12.

## 7. Closure-Notiz

*(Wird bei der Closure gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Zählung nach `AGENTS.md` §3.13:** 1 Liefer-Punkt (das Dokument), 1 Schicht (Nutzer- und
Wartungs-Doku; `docs/developer/` zählt wie `docs/user/` und `docs/maintainer/` als
Dokumentations-Schicht, die Planung zählt nicht). Unter den Grenzen von drei Punkten und zwei
Schichten.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist nur die Default-Sub-Area `*` (Kürzel
`REPO`, GF) aus `harness/conventions.md`; sie erfüllt die Schwelle als Default für das gesamte Repo.

**Vorgelagert — offene Beobachtungen sichten:** Register unter `../observations/` durchgegangen
für `REPO`: der einschlägige Eintrag ist `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`
(Zusage über die Prüfung hinaus); das Dokument sagt nur Nachgeführtes zu (§1), der Zähler-Stand
ist bei der Closure aus den Dateien zu lesen.

Alle berührten Sub-Areas GF.
