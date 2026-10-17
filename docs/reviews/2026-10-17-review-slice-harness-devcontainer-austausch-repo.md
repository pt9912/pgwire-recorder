# Review-Report: slice-harness-devcontainer-austausch-repo — 2026-10-17

**Review-Art:** Code (Dokument-Diff) — gegen Plan + Konventionen (`AGENTS.md` §3.7, §3.9, §3.11, §3.13).

**Gegenstand:** `git diff 38a9b20a..HEAD` (Commits 9f8faef9, 1e70abc3), Zweig `slice/devcontainer-austausch-repo`; geändert sind `docs/developer/devcontainer.md` und der Slice-Plan.

**Skill:** `.harness/skills/reviewer.md` @ 1e70abc3 ·
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

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-harness-devcontainer-austausch-repo` (§1 bis §8, darin §7 Belege des Implementers)
- keine aktive ADR, kein `LH-*` (Kopf: ohne Bezug)
- `AGENTS.md` (Hard Rules), insbesondere §3.7, §3.9, §3.11, §3.13, §3.14
- Selbst gefahren: `grep` auf `host-repo`, `PGR_REPO`, `Variante`, `scp`, `bundle`, `/Development`, `vier Mounts` (alle leer), Quantoren-`grep`, `git diff --numstat`, `wc -l`, `make docs-check` im Worktree (`0 Befund(e)`). Kein Docker-Container gestartet; die Beobachtungen der Proben sind aus §7 des Plans gelesen, nicht wiederholt.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Aussage zum Host-Repo schließt aus „kein Mount“ („also“), dass im Container keine uncommitteten oder ignorierten Dateien des Hosts liegen. Dasselbe Dokument nennt den Container `privileged` mit `sudo` ohne Passwort; die Vorfassung trug dazu den Vorbehalt „gegen einen Prozess mit `sudo` ist es keine Zusage“, der entfallen ist. Ob ein Root-Prozess im Container über Gerätedateien des Hosts dessen Arbeitsbaum lesen kann, ist weder genannt noch geprüft (§13 führt es nicht). | `AGENTS.md` §3.11 (Sicherheitsaussage weiter als Prüfung) | `docs/developer/devcontainer.md` · „also liegen darin keine uncommitteten oder ignorierten Dateien des Hosts“ | nein — nur ein Lauf im privilegierten Container (Lesen eines Host-Geräts) würde es zeigen | Sicherheitsaussage ohne Vorbehalt zu privileged/sudo |
| F-2 | MEDIUM | §7 sagt „in keinem Volume (wo `claude` sie ablegt, ist nicht geprüft)“, §11 „liegen im Dateisystem des Containers (kein Volume)“, §4 „entfällt damit ebenfalls (… gefolgert)“; §13 führt „wo `claude` seine Anmeldedaten ablegt“ als nicht geprüft. Die Verneinung „kein Volume“ ist ohne Lauf, und der Satz widerspricht seiner Klammer; legt `claude` etwas unter `/workspaces` ab, überlebt es `u-boot down` und könnte mit dem Klon weitergegeben werden. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „liegen danach im Container, nicht im Host-Repo und in keinem Volume“ und „im Dateisystem des Containers (kein Volume)“ | nein — erst ein Lauf von `claude` mit Anmeldung | Verneinung ohne Lauf, Widerspruch zu §13 |
| F-3 | MEDIUM | §10 schließt mit „nehmen `-p "$PGR_SSH_PORT"` und `PGR_DEV` der Instanz“, aber die Variable steht in diesem Ablauf nur in der `.env` von Compose, nicht in der Shell. Der Befehl wäre `-p ""` (Fehler); in §5 und §9 gilt `${PGR_SSH_PORT:-2222}`, was für die zweite Instanz still auf Port 2222, also den ersten Container, zeigt (`git remote add container … :2222` am Host). Der Beispielblock selbst nutzt das Literal 2223, die Querverweis-Regel widerspricht ihm. | `AGENTS.md` §3.11 (Befehl nicht wie geschrieben ausführbar) | `docs/developer/devcontainer.md` · „nehmen `-p "$PGR_SSH_PORT"` und `PGR_DEV` der Instanz“ | ja — `bash -c 'echo ${PGR_SSH_PORT:-2222}'` ohne Export gibt 2222; die Probe in §7 nutzte feste Ports | Befehl im Querverweis nicht ausführbar wie geschrieben |
| F-4 | LOW | Der Block in §10 beginnt mit `cd "$PGR_DEV"` (Verzeichnis ohne Git); der folgende Block `git remote add container2 …` / `git config core.sshCommand …` läuft laut §9 „im Host-Repo“, das §10 nicht neu nennt. Im Verzeichnis `pgr-zwei` endet er mit „not a git repository“ (oder setzt `core.sshCommand` in ein umgebendes Repo). | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „git remote add container2 ssh://vscode@127.0.0.1:2223“ | ja — Block hintereinander im Verzeichnis `pgr-zwei` ausführen | Arbeitsverzeichnis des Codeblocks nicht genannt |
| F-5 | LOW | §4 nennt für den Host-Schlüssel nach `down`/`up` „(gemessen, `git fetch` ohne Warnung)“. Die Probe in §7 belegt den Fingerabdruck und einen **Push** ohne Warnung, nicht den Fetch. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „(gemessen, `git fetch` ohne Warnung)“ | ja — Fetch nach `down`/`up` | Beleg nennt anderen Befehl als die Aussage |
| F-6 | LOW | §3.5 sagt, der Block weiche von „der Datei ab, die `u-boot` mit `init` erzeugt“, im Schlüsselpfad. Zwei Absätze davor steht, `init` schreibe `services: {}`; die erzeugte Datei hat keinen Schlüsselpfad. Gemeint ist der Messort. | `AGENTS.md` §3.7/§3.11 | `docs/developer/devcontainer.md` · „von der Datei ab, die `u-boot` mit `init` erzeugt“ | nein — Lesefehler im Text | Abschnitte widersprechen sich |
| F-7 | LOW | Die Tabelle in §7 des Plans deckt die geänderten Aussagen; mehrere Aussagen des Dokuments haben dort keine Zeile: §8 (Zeiten 191 s, 1 s „in einem frischen Container“, „Probe“), §11 Platz („im frischen Container der Probe gleich“), §5 (`ss -ltn`, `docker port` mit Vorgabe 2222), §12 (iptables, Permission denied). Sie stammen aus der Probe des Vorgängers; „die Probe“ meint im Dokument nun auch den Scratch-Container `pgr-aus`, der diese Läufe nicht fuhr. Die DoD verlangt „jede Aussage“ belegt oder in §13. | `AGENTS.md` §3.11; Plan §2 DoD | `docs/developer/devcontainer.md` · „Die zweite stammt aus der Probe in einem frischen Container“ | nein — Zuordnung Aussage · Lauf ist Urteil | Aussage ohne Zeile in der Probentabelle; „Probe“ mehrdeutig |
| F-8 | LOW | Die Quantoren-Durchsicht in §7 zitiert „Jeder `ssh`-Aufruf von `git`“ (§9) und „Er bekommt nur, was du … einlieferst“ (§11) als eng gefasst; beide Sätze stehen im Dokument nicht (§9: „Ein `git ls-remote` auf eine `ssh://`-URL …“, §11: „Das Bare-Repo enthält, was du … einlieferst“). Der Plan folgt dem Diff an dieser Stelle nicht. | `AGENTS.md` §3.9 | `docs/plan/planning/in-progress/slice-harness-devcontainer-austausch-repo.md` · „Er bekommt nur, was du … einlieferst“ (§11)“ | ja — `grep` nach beiden Sätzen im Dokument leer | Plan zitiert Satz, der im Diff nicht steht |
| F-9 | LOW | §9 sagt für den Fall „beide Seiten schreiben auf `main`“: „im Container holt `git pull` den Stand“. Der Container arbeitet auf `thema`, das in §6 ohne Upstream gepusht wird (`git push origin thema`); die Probe in §7 zog `git pull` nur auf `main`. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „im Container holt `git pull` den Stand“ | ja — `git pull` auf `thema` ohne Upstream | Aussage über den Branch hinaus, den die Probe lief |
| F-10 | INFO | §4 „Stände, die nur dort lagen, auch“ trägt im Plan die Einordnung „gefolgert“, im Dokument steht keine Kennzeichnung; es folgt aus `Removed`. Für Rolle Implementer/Planner. | `AGENTS.md` §3.11 | `docs/developer/devcontainer.md` · „Stände, die nur dort lagen, auch“ | nein | Einordnung im Plan, nicht im Dokument |
| F-11 | INFO | Zitierte Meldungen sind teils in der Locale des Probe-Hosts (Deutsch: „externes Repository container existiert bereits“, „Bitte stellen Sie sicher …“), teils englisch aus dem Container; ein Host mit anderer Locale sieht andere Texte, Exit-Codes (3, 128, 255) bleiben. | Maintainability | `docs/developer/devcontainer.md` · „externes Repository container existiert bereits“ | nein | Locale-abhängige Meldung als Beleg |
| F-12 | INFO | Der Satz in §1 „Das Repo liegt als Klon in einem Volume des Containers, der aus einem Bare-Repo im selben Volume entsteht“ ist doppeldeutig (Container oder Klon entsteht aus dem Bare-Repo). | Maintainability | `docs/developer/devcontainer.md` · „Containers, der aus einem Bare-Repo im selben Volume entsteht“ | nein | Bezug eines Relativsatzes unklar |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Reste `/host-repo`, `PGR_REPO`, Variante A/B, `scp`, `bundle`, „vier Mounts“ (`grep`) | geprüft, ohne Befund (leer); auch `/Development` und feste Home-Pfade außerhalb des Containers kommen nicht vor, Host-Pfade nur als `PGR_DEV` |
| Querverweise `§1` bis `§13`, Anker, `make docs-check` | geprüft, ohne Befund (`0 Befund(e)`); Verweise auf §3.2 bis §3.5, §6, §9, §10, §11, §12 passen zu Inhalt und Reihenfolge |
| Compose-Block, Dockerfile-Block, `start.sh` gegen Text (drei Mounts, Port, Healthcheck, `chown /workspaces`) | geprüft, ohne Befund; Blöcke unverändert zur Vorfassung, Plan §7 belegt `diff` leer und drei Mounts |
| Größe nach Plan §4 und `AGENTS.md` §3.14 | geprüft, ohne Befund: `git diff --numstat 38a9b20a..HEAD -- docs/developer` = 151 + 100 = 251 (Basis `92cfaf23` gleiche Zahl), Marke 450; Dokument 549 Zeilen, Grenze 560 |
| Plan §3.9 (Kopf, §1, §3, §6) und §3.13 (Nehmer, Zählung §8) | geprüft; Kopf, §1, §3, §6 folgen dem Diff (bis auf F-8); keine Sendung an einen anderen Slice, Zählung in §8 vorhanden |
| Zielbild §10: zwei `-i` in `core.sshCommand`, Remote-Namen, Ports | geprüft, ohne Befund bis auf F-3 und F-4; Probe in §7 deckt „einer allein abgewiesen, beide laufen“ |
| Portabilität (kein `/Development`, nur `PGR_DEV`-Form) | geprüft, ohne Befund |
| Sicherheit: `sudo`, SSH-Port, Schlüssel ohne Passphrase, Loopback, Review vor Merge | geprüft; Aussagen zu `sudo`, Loopback und Schlüssel tragen und sind als gefolgert gekennzeichnet; Befund nur F-1 und F-2 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 6 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Sicherheitsaussage ohne Vorbehalt zu privileged/sudo · Verneinung ohne Lauf, Widerspruch zu §13 · Befehl im Querverweis nicht ausführbar wie geschrieben · Arbeitsverzeichnis des Codeblocks nicht genannt · Beleg nennt anderen Befehl als die Aussage · Abschnitte widersprechen sich · Aussage ohne Zeile in der Probentabelle · Plan zitiert Satz, der im Diff nicht steht

## Verdikt

**Merge-blockierend:** ja — F-1 bis F-3 (MEDIUM) blockieren; die LOW- und INFO-Findings gehen mit.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler (Klasse „Verneinung ohne Lauf“ gehört zum Muster `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, Vorgänger zweimal). Dieser Report selbst
ist ein **Lauf-Beleg**. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
