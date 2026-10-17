# Verifikationsbericht: slice-harness-devcontainer-u-boot-0-10 — 2026-10-17

**Rolle:** Verifier (Modul 11) · **Frage:** Bauen wir es richtig, gegen Plan und DoD?
**Gegenstand:** Diff `8070a14d..HEAD` (Zweig `dev-container-u-boot`, Commits `bf1d170e`, `641a5b14`): `docs/developer/devcontainer.md` (+297/−301) und §7 des Slice-Plans (+40). Das Gerüst `$PGR_DEV` liegt außerhalb des Repos und ist nur über das Dokument und §7 sichtbar.
**Eingang:** Plan §1 bis §7, Dokument Version 0.3. Kein Review-Report im Zweig gefunden (DoD Punkt „Review“ offen).

## Eigene Läufe

| Lauf | Ergebnis |
|---|---|
| `make docs-check` im Worktree | 1029 Dateien, 0 Befunde, Exit 0 |
| `make kopf-check` im Worktree | Exit 0 |
| `git diff --numstat 05fd1403..HEAD -- docs/developer` | 297 / 301, Summe **598** (stimmt mit §7) |
| `docker ps -a` / `docker volume ls` nach `probe` | kein Container, kein Volume (Scratch-Entfernen am Host bestätigt; §7 selbst enthält keinen Beleg dafür) |
| `make gates`, Docker-Proben im Container | nicht gefahren (Auftrag); gestützt auf Belege des Implementers |

## DoD-Konformität

| DoD-Punkt | Urteil | Beleg / Befund |
|---|---|---|
| 1 Gerüst auf v0.10.0 (Features, Zugang als Variable, beide Wege, Beleg) | teilweise | Build, Start, `ssh`, `claude --version`, Variablenlänge (23/28) belegt in §7. **Beide Wege „funktionieren“ ist nicht belegt**: nur Platzhalter (401 / Invalid API key); echter Login offen und in Dok. §7/§13 ehrlich als offen genannt. Gerüst-Dateien stehen im Dokument (§3.3 bis §3.5) |
| 2 Messung und Probe | erfüllt mit Lücken | (a) `make gates` Exit 0, 1701 s, `pgrep` leer vor und nach: belegt. (b) Fehlend/leer/ungültig, Platzhalter-Token, Neustart ohne Neubau: belegt. (c) Laufzeitfrage als Frage in §13: ja. Lücken siehe V-2 bis V-5 |
| 3 Dokument | erfüllt mit Befunden | alle geforderten Abschnitte da (§1, §2, §3 bis §5, §7, §8, §11 inkl. Leitplanke und `doctor`-Warnung, §13 mit `podman`/`docker`, Hostschlüssel, Umbau). Befunde V-1, V-2, V-6 bis V-9 |
| Gates: `docs-check`, `kopf-check` | erfüllt | eigene Läufe oben. `make gates` nicht gefahren; Diff liegt nur auf Positivliste (`docs/developer/`, `docs/plan/planning/`), Freigabe durch Stop-Hook plausibel |
| Review, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen | §7 „Closure-Notiz“ trägt `(wird bei der Closure geschrieben)`; Pflichten unten |

## Nicht Gemessenes: als offen gekennzeichnet?

| Punkt | Urteil |
|---|---|
| Echter Login (API-Schlüssel, Abo) | ja: Dok. §7 („lief nicht“, zwei Befehle für den Nutzer) und §13 |
| Abgelaufenes Token mitten in der Sitzung | ja: Dok. §7 und §13; ersetzt laut Plan durch Neustart-Messung |
| Interaktive Oberfläche von `claude` | ja: §7, §13 |
| `podman` | ja: §13 (Zitat aus Changelog, nicht gelaufen); `nestedRuntime: docker` gemessen |
| Login-Hosts im Browser | ja: §7, §11, §13 |
| Umbau des echten Containers, Hostschlüssel-Einstellung, Colima/Docker Desktop | ja: §13, §1 |

Die Kennzeichnung ist sauber. Kein Satz im Dokument behauptet einen echten Login.

## Plan-vs-Code-Diff (Plan wird erst nach Review nachgezogen; hier nur festgehalten)

1. **Neustart:** Plan §1/§2 nennt `u-boot down`/`u-boot up`. Das Gerüst startet über `devcontainer up`; `u-boot up` endet mit `no service selected` (Exit 12). Das Dokument verwendet `docker stop` und `devcontainer up`.
2. **Neue Host-Voraussetzungen:** `devcontainer`-CLI (npm, Node am Host) und `mk-override.sh` (von Hand, nicht im Plan §3). Plan-Kopf `Bezug` sagt, es gelte nur `ssh` und `u-boot` zusätzlich. Die CLI und Node am Host fehlen dort.
3. **Gerüst-Dateien:** Plan §3 nennt `u-boot.yaml`, `Dockerfile`, `start.sh`, `compose.yaml`; tatsächlich `start.sh` entfällt, `compose.yaml` ohne Dienst und ungenutzt, `Dockerfile`, `docker-init.sh`, `egress-init.sh`, `devcontainer.json` von `u-boot` erzeugt, `mk-override.sh` neu.
4. **Zugang:** Plan nimmt an, `remoteEnv` nenne nur `ANTHROPIC_API_KEY`. Tatsächlich nennt `u-boot` auch `GIT_TOKEN` (Dok. §3.4); `CLAUDE_CODE_OAUTH_TOKEN` wird über `mk-override.sh` ergänzt. Plan-Frage „wie durchgereicht“ damit beantwortet.
5. **Befund mit Wirkung auf den Nehmer** `slice-harness-devcontainer-einrichten` (Quellenauflösung): `remoteEnv` greift nur bei `devcontainer exec`, nicht in `ssh`-Sitzungen und nicht bei `docker exec`. Eine Auflösung, die beim Start durchreicht, trägt in der `ssh`-Sitzung nicht; der Nehmer muss das in §1/§6 aufnehmen (AGENTS.md §3.13).
6. **Annahme „Token beim Start frisch“** (Plan §1, Dok. §8) und Dok. §7 („Variable wird bei jedem `devcontainer exec` neu gelesen, Neustart nicht nötig“) widersprechen sich leicht: die Begründung für „ein Slice pro Sitzung“ trägt für den `devcontainer exec`-Weg nicht mehr.
7. **Größe:** Summe 598 über Meldeschwelle 350 und Marke 450 (§4). §7 benennt das, die Entscheidung liegt beim Planner. Eine Rückführung ist nach §4 für „Summe über 450“ vorgesehen; der Schnitt (Punkte 1 und 3 zuerst, Messung als Folge-Slice) passt nicht mehr, weil die Messung bereits vorliegt. Planner entscheidet (Ausnahme wie bei `slice-harness-coverage`, Entscheidung des Nutzers).
8. **Randform 10** (`/dev/net/tun`, Host ohne Colima) ist im Dokument nur über §1/§13 (Linux-Host) beantwortet; `/dev/net/tun` kommt nicht vor.

## Befunde

| ID | Schwere | Befund |
|---|---|---|
| V-1 | MEDIUM | Speicher widersprüchlich: §7 nennt für den Lauf mit Cache aus Fehlläufen **5,424 GiB**, Dok. §8 (Tabelle, Zeile 2) lässt den Wert weg, §8/§11 sagen „höchstens 4,2 GiB“. Die Aussage „höchstens“ ist durch §7 widerlegt. Gegen `AGENTS.md` §3.11 (Zusage weiter als Messung). |
| V-2 | MEDIUM | Zwei Läufe mit identischen 1701 s: Das Dokument nennt die Gleichheit „Zufall“ (§8). Das ist eine ungeprüfte Behauptung (nur §7 „zuvor“ vs. „frisches Volume“, Reihenfolge der Läufe unklar). Zeit zweier Läufe sollte aus den Rohdaten bestätigt oder die Zeile gestrichen werden. |
| V-3 | MEDIUM | Der Zugangs-Pfad mit `GIT_TOKEN` (aus `remoteEnv`, von `u-boot` erzeugt) kommt in §7, §11 (Zugangsdaten) und in der Sichtbarkeitsprobe nicht vor. Wo er sichtbar ist, ist ungemessen und ungenannt. |
| V-4 | LOW | Randform 3: Plan nennt Kindprozesse und Shell-History; die Sichtbarkeitsprobe deckt `ps`, `environ`, `inspect`, `logs` ab. Kindprozesse und History fehlen im Dokument (kein Vermerk in §13). Die Aussage „nicht in Dateien unter `/home/vscode` gefunden“ hat keine Zeile in §7. |
| V-5 | LOW | Nicht in der Tabelle §7 belegt, aber im Dokument als „gemessen“ genannt: `ss -ltn`, `/etc/shadow`, `sudo -n -l`, `AcceptEnv`, `tty` `/dev/pts/0`, `cgroup memory.max`, `docker info` (20 CPUs), Mounts (drei), Eigentümer von `/workspaces`, `git version 2.55.0`, 1,7 und 7,9 s (§7 nennt 1,8 und 6,4). Der Plan verlangt „Beleg in §7 als Tabelle“; die Tabelle ist eine Verdichtung. Ich konnte diese Werte ohne Docker-Läufe nicht prüfen. |
| V-6 | LOW | Dok. §7: „das Token wurde von `api.anthropic.com` geprüft“ ist aus der 401-Meldung geschlossen; die Namensauflösung ist in §7 nur global („im Log nur Namen der Liste“) belegt. |
| V-7 | LOW | Dok. §2 verweist für „Docker Engine ≥ 24“ auf `u-boot doctor` und §12; §12 sagt dazu nichts. Der Verweis läuft ins Leere, die Aussage ist ungemessen. |
| V-8 | LOW | §7 des Plans trägt keinen Beleg zum Entfernen des Scratch-Containers samt Volumes, obwohl der Closure-Trigger (§5) das im Bericht verlangt. Mein Lauf: nichts mehr vorhanden. |
| V-9 | INFO | Plan-Kopf `Bezug` und Dok. Voraussetzungen stimmen für die Host-Anforderungen nicht überein (Plan-vs-Code Punkt 2). |

Keine DoD-Verletzung im Sinn von „behauptet, aber nicht gemessen“ beim Login; V-1 ist die nächste Stufe (Zahl im Dokument widerlegt durch §7).

## Pflichten für den Planner / Implementer vor Closure

1. V-1 beheben (Speicher „höchstens 5,4 GiB“ oder Messung erklären), V-2, V-3 klären; V-4 bis V-8 im Dokument eingrenzen oder in §7 nachtragen.
2. Plan nach dem Review nachziehen: §1/§2/§3 (Punkte 1 bis 4, 8 oben), Kopf `Bezug`, Nachzählung `slice-harness-devcontainer-einrichten` (Punkt 5: `remoteEnv` nur bei `devcontainer exec`) im selben Commit (§3.13).
3. Größe: 598 über Marke 450 — Entscheidung des Planners; Ausgang für Risiko „Größe“ ist *eingetreten* oder Abweichung nach Entscheidung des Nutzers.
4. Risiko-Ausgänge: `make gates` im Container (entfallen mit 1701 s, grün), Laufzeitfrage (weiter offen → `BEO-REPO/nested-runtime-podman-oder-docker-offen`), Hostschlüssel (weiter offen, Befund: Fingerabdruck ändert sich bei Neubau), Messung neben anderem Lauf (entfallen, `pgrep` leer), Vorbehalte der Vorfassung (Vergleich im Review nötig; §11-Vorbehalte „keine Sandbox“ erhalten), Start der Release-Slices, Stop-Hook im Hauptbaum (weiter offen).
5. Review-Report nachreichen, Closure-Notiz mit Lerneintrag, Register, Paarungen, dann reiner `git mv` nach `done/`.
