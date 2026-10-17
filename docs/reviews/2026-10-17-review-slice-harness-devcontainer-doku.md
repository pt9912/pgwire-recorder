# Review-Report: slice-harness-devcontainer-doku — 2026-10-17

**Review-Art:** Code (Dokument) — Review gegen Plan + Konventionen (`AGENTS.md` §3.7, §3.9, §3.11, §3.14).

**Gegenstand:** slice-harness-devcontainer-doku, Diff `bda6efef..1818f42a -- docs/developer docs/plan/planning/in-progress/slice-harness-devcontainer-doku.md` (Commits 0657b0e5, 531f357e, 1818f42a).

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec08
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-harness-devcontainer-doku` (§1, §3, §4, §6, §7 Belege des Implementers, §8)
- aktive ADRs: keine (Plan-Kopf: „Keine aktive ADR“)
- berührte `LH-*`-IDs: keine (Plan-Kopf `Berührte Spec-Stellen: —`)
- `AGENTS.md` (Hard Rules)
- Proben (lesend) am Versuchs-Container `pgr-dev` (SSH, `docker inspect`, `u-boot doctor`, `u-boot --help`)

---

## Findings

Spalten gespiegelt aus dem §Output-Schema des Reviewer-Skills; bei Abweichung gilt der Skill. Zuständigkeit: Implementer, soweit nicht anders genannt.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Das Dokument sagt „der Container kann nicht nach `/host-repo` schreiben“ und führt „Host-Repo: schreibgeschützt“ als Schutz auf. Der Benutzer `vscode` hat im Versuchs-Container `sudo` ohne Passwort, der Container hat alle Capabilities (`capsh`: `=ep`), und `/host-repo` ist das Host-Blockgerät `/dev/nvme0n1p3` mit `ro`. Ein Remount oder Schreiben über das Gerät ist damit für einen Prozess im Container möglich (nicht ausgeführt, um den Host nicht anzufassen); belegt ist nur, dass `touch` als `vscode` ohne `sudo` scheitert. Bei einer Claude-Code-Session im Container (§7) gilt „kann nicht“ so nicht. | `AGENTS.md` §3.11 (nur zusagen, was eine Probe belegt) | `docs/developer/devcontainer.md` · „der Container kann nicht nach `/host-repo` schreiben“ (§9); · „**Host-Repo:** schreibgeschützt als `/host-repo`.“ (§11) | ja — Probe `sudo -n true` und `mount` im Container; ein Remount-Versuch wäre die Gegenprobe | Sicherheitsaussage weiter als die Probe (Schutzzusage trotz privilegiertem Container mit `sudo`) |
| F-2 | LOW | Die Zeiten der Zeilen „Versuchs-Container“ in der Tabelle (5 s, 201 s, 1 s) hat der Implementer nicht selbst gemessen; §7 des Plans nennt sie als Angaben des Nutzers aus `scripts/nachweis2.sh` und `nachweis3.sh` unter `/Development/pgr-dev`. Im Dokument steht die Quelle nicht, die Zeilen lesen sich wie eigene Messungen. | `AGENTS.md` §3.11 (Messwerte als Messwerte mit Quelle) | `docs/developer/devcontainer.md` · „erster Lauf im Versuchs-Container; Probe im Klon eines zweiten, frischen Containers“ (§8) | nein — die Läufe sind nicht wiederholt | Messwert ohne Quellenangabe im Dokument |
| F-3 | LOW | §9 sagt zu, dass nach Variante A oder B „der Branch `container/<thema>` im Host-Repo“ liegt. Beide Wege hat der Implementer nur in ein Scratch-Klon bzw. leeres Repo geholt, nicht ins Host-Repo; §13 „Nicht geprüft“ nennt das nicht. Zudem fragt `ssh` bei der ersten Verbindung zu `127.0.0.1:2222` nach dem Hostschlüssel (`StrictHostKeyChecking ask`, kein Eintrag in `known_hosts`); das Dokument sagt dazu nichts. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „Danach liegt der Branch `container/<thema>` im Host-Repo“ (§9) | nein | „Nicht geprüft“ unvollständig (Fetch ins echte Host-Repo, Hostschlüssel-Frage bei Erstverbindung) |
| F-4 | LOW | §4 sagt „`docker compose up -d` baut das Image, wenn es fehlt, und startet ebenso.“ §7 belegt nur `u-boot up` mit Build; `docker compose up -d` wurde nur im Port-Kollisions-Fall (Start schlägt fehl) gefahren. `u-boot up --help` nennt zwar `docker compose up -d` als Mechanismus, ein Lauf mit Build ohne `u-boot` steht nicht in der Tabelle. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „`docker compose up -d` baut das Image, wenn es fehlt, und startet ebenso.“ (§4) | ja — Lauf von `docker compose up -d` mit fehlendem Image im Scratch-Verzeichnis | Zusage aus Hilfetext statt aus Probe |
| F-5 | LOW | Die Tabelle in §3.1 („Was `u-boot` erzeugt“) nennt `u-boot.yaml`, `.env.example`, `.gitignore`, `README.md`, `CHANGELOG.md`; die Ausgabe von `u-boot init` in §7 nennt zusätzlich die Verzeichnisse `docker/`, `scripts/`, `docs/`. Die Tabelle liest sich als vollständige Liste. | `AGENTS.md` §3.11 (Zahlwörter/Vollständigkeit) | `docs/developer/devcontainer.md` · „Was `u-boot` erzeugt und was von Hand dazukommt:“ (§3.1) | ja — `u-boot init` im leeren Verzeichnis, `ls` | Liste liest sich vollständig, ist es nicht |
| F-6 | INFO | Der Befehl `git config --global --add safe.directory /host-repo` steht im Codeblock von §6 als Schritt, der Absatz darunter sagt, er sei in der Probe nicht nötig gewesen. Beides ist ehrlich, der Schritt ist nicht probiert „mit Wirkung“. Zuständig: keine Aktion erwartet. | Maintainability | `docs/developer/devcontainer.md` · „`git config --global --add safe.directory /host-repo`“ (§6) | nein | — |
| F-7 | INFO | Das Schlüsselpaar wird ohne Passphrase erzeugt (`-N ''`) und gibt Zugang zu einem privilegierten Container mit `sudo`; das Dokument nennt die fehlende Passphrase nicht. Der Plan-Kopf `Bezug` sagt „kein Anspruch des Produkts“, der Übergangs-Commit 0657b0e5 nennt dennoch `LH-QA-03`; die Traceability-Regel ist über die Slice-Kennung erfüllt, der Verweis ist nur nicht im Kopf begründet. Zuständig: Planner (Kopf) / Implementer (Dokument). | Maintainability | `docs/developer/devcontainer.md` · „ssh-keygen -t ed25519 -N '' -C pgr-dev -f ssh/id_ed25519“ (§3.2) | nein | — |

## Prüfungen nach Auftrag (ohne Befund, soweit nicht oben)

1. **Codeblöcke:** `Dockerfile` (Zeilen 87–115), `start.sh` (125–134) und `compose.yaml` (148–187) stimmen per `diff` mit `/Development/pgr-dev/.devcontainer/Dockerfile`, `…/start.sh` und `…/compose.yaml` überein.
2. **Befehle gegen Container/Hilfe:** `u-boot doctor` (16 ok, Docker ≥ 24, Compose ≥ 2.20, Prüfliste wie in §12), `u-boot up --help` (Timeout 60 s, Exit 12), `u-boot add --help`, `docker info` (Name `pgr-dev`, 20 CPUs), `docker ps -a` im Container leer, `claude --version` 2.1.302, `sshd`-Konfiguration, einziger Schlüssel in `authorized_keys`, `/workspaces` enthält `pgwire-recorder` stimmen mit §7 des Plans und dem Dokument überein. Der Bild-/Volume-Platz (2,05 GB, 5,87 GB), die Zeiten 191 s und 44 s und der Probe-Container selbst sind nicht wiederholt (Probe-Container entfernt, Docker belegt); gelesen gegen §7, dort schlüssig, nicht unabhängig neu belegt. Hostschlüssel-Aussage: Die Schlüsseldatei des Containers ist älter als der Container und stammt aus dem Image, daher passt „nach Neubau des Images“ in §13.
3. **„Nicht geprüft“ (§13):** ehrlich in dem, was sie nennt; unvollständig laut F-3 (und die Ursache „Zeit eines Builds ohne Cache“ steht nur in §4, nicht in §13).
4. **Kennungen im Text:** Treffer von `LH-`/`SPEC-`/`ADR-`/`slice`/`welle` im Dokument nur in den Kommentaren des `u-boot`-Blocks (`LH-FA-DEV-004`, `LH-FA-ADD-*`, `LH-SA-FILE-002`), die `u-boot` selbst schreibt und die byte-gleich zum Versuch bleiben müssen; sie sind Kennungen eines fremden Projekts, keine dieses Repos. Der Plan führt für `docs/developer/` die Linie „wie Ist-Zustand ohne Chronik, keine Kennungen“ (§3 des Plans); eingehalten. Chronik: „Stand: 2026-10-17“ im Kopf ist ein Datum, keine Chronik.
5. **Sicherheitsaussagen:** Port nur `127.0.0.1` (`docker port`/Plan-Beleg `ss`), nur Schlüssel (`sshd_config.d/pgr-dev.conf`), kein Docker-Socket, der private Schlüssel „bleibt im Verzeichnis `ssh/` auf dem Host“ (außerhalb des Repos) belegt; Ausnahme `/host-repo`-Schutz laut F-1.
6. **Positivliste (§6 des Plans):** korrekt genannt; `docs/developer/` steht nicht auf der Positivliste in `harness/README.md` (`docs/plan/planning/`, `docs/reviews/`, `docs/user/`, Wurzel-`README.md`), der volle `make gates` ist vor dem Push nötig. Das Risiko trägt den Ausgang „weiter offen“ mit Zeiger auf die Entscheidung des Nutzers; bei der Closure ist es ins Register zu tragen.
7. **`make docs-check`:** nicht gefahren (Hintergrundlauf von `make gates`); §7 des Plans nennt Exit 0 (966 Dateien), nicht unabhängig bestätigt.
8. **§3.14 Größe:** `git diff --numstat bda6efef..HEAD -- docs/developer` = 385 Einfügungen, 0 Löschungen, Summe 385; unter Meldeschwelle 400 und Marke 500, wie in §7 des Plans.
9. **§3.9 Plan folgt Diff:** §1 und §3 nennen die Abschnitte und Messwerte (docs-check 5 s, gates-schnell 201 s), §8 zählt 1 Liefer-Punkt, 1 Schicht; kein Widerspruch zum Diff gefunden.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/developer/` (Codeblöcke, Befehle, Kennungen, Größe) | geprüft, Befunde F-1 bis F-6 |
| `docs/plan/planning/in-progress/slice-harness-devcontainer-doku.md` (§1, §3, §4, §6, §8) | geprüft, ohne Befund (F-7 betrifft nur den Kopf) |
| `internal/`, `cmd/`, `test/`, `harness/`, `.claude/` | nicht im Diff, nicht betrachtet |
| `spec/`, `docs/plan/adr/` | nicht im Diff; Stratum-Regel (§3.4) nicht berührt, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Sicherheitsaussage weiter als die Probe · Messwert ohne Quellenangabe im Dokument · „Nicht geprüft“ unvollständig · Zusage aus Hilfetext statt aus Probe · Liste liest sich vollständig, ist es nicht

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM) ist vor der Closure zu beheben oder durch den Architect zu entscheiden; F-2 bis F-5 blockieren nicht.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` ist die nächstliegende Beobachtung für die erste und vierte Klasse). Dieser Report ersetzt keine Verifikation.
