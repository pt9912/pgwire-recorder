# Verifikation: `slice-release-v0-2-0` und `slice-go-toolchain-1-27-3`

**Rolle:** Verifier (DoD- und Entscheidungs-Konformität, Plan gegen Diff). **Datum:** 2026-10-24. **Stand:** `main` = `6630a708`, Zweig `ver-rel-v020` (von `main`), Arbeitsbaum `/Development/pgwire-recorder-verrel`. **Pläne:** [`slice-release-v0-2-0`](../plan/planning/in-progress/slice-release-v0-2-0.md) (Basis `193f0e0b`) und [`slice-go-toolchain-1-27-3`](../plan/planning/in-progress/slice-go-toolchain-1-27-3.md) (Basis `7cfad5ab`). **Gegenstand:** Eintrag [`freigaben/0.2.0.md`](../maintainer/freigaben/0.2.0.md), Beleg [`2026-10-24-nachweis-slice-release-v0-2-0.md`](2026-10-24-nachweis-slice-release-v0-2-0.md), Carveout [`CO-002`](../plan/carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md), Tag `v0.2.0`, Lauf `40074219261`. **Maßstab:** [`SPEC-067`](../../harness/vertraege/SPEC-067-release.md) Punkte 17, 18 und 20, `docs/maintainer/releasing.md` §10, `AGENTS.md` §3.9, §3.11, §3.14.

Der Review-Bericht `2026-10-24-review-slice-release-v0-2-0.md` wurde erst nach dem Urteil gelesen (Abschnitt *Abgleich mit dem Review* am Ende).

## Urteil

**Bestanden, mit Auflagen für die Closure.** Jede Zeile der Freigabe-Tabelle stimmt mit dem Original überein. Es gibt keine HIGH-Abweichung. Zwei MEDIUM-Befunde sind Closure-Bedingungen des Pin-Slice und keine Fehler der Lieferung. Beide Slices sind noch nicht geschlossen (`in-progress/`, §7 „noch nicht geschlossen“), die offenen Closure-Punkte stehen unten.

## Selbst gefahrene Sensoren (nicht übernommen)

| Sensor | Ergebnis |
|---|---|
| `make gates` an `6630a708` (sauberer Baum, `/Development/pgwire-recorder-verrel`; auf dem Daemon liefen zeitgleich fremde Container und ein weiterer `make gates`) | Exit 0, 2026-10-24 08:19:42Z bis 08:24:25Z (283 s), Last (1-Minuten-Mittel) 8,99 bei 20 Kernen um 08:20:49Z. `main` unterscheidet sich vom Tag-Commit `32d0d84a` nur in Planung, Eintrag, Handbuch, `README.md`, `releasing.md` und dem Beleg (`git diff --stat 32d0d84a main -- . ':!docs'` nennt allein `README.md`). |
| `GP_MUT='^ok$' make release-gegenprobe` (engster Fall, Nutzer-Regel) | Exit 0, `release-gegenprobe: grün (1 Läufe)`, 175 s, Last 15,90 bei 20 Kernen zu Beginn. Die Skripte sind an `main` und `32d0d84a` gleich. |
| `GP_MUT='^kein-treffer$' make test-container-gegenprobe` (engster Fall) | Exit 0, 27 s. Die Schlusszeile sagt „je Mutation genau die benannten Fälle rot“; gefahren ist allein die unveränderte Kopie, keine Mutation (Wortlaut wie im Eintrag Punkt 2). |
| Keine vollen Gegenproben, kein `make homebrew-nachweis` | nach Auftrag |

## Zeile für Zeile gegen die Originale

| Zeile des Eintrags | Original | Ergebnis |
|---|---|---|
| Tag `v0.2.0` | `git ls-remote origin`: Tag-Objekt `fa5de159…`, aufgelöst `32d0d84a924a1e67a44e9cee9441a69b22eb4fdd`; `32d0d84a` ist Vorfahr von `main`; `R` = `d720e44f` ist Vorfahr von `32d0d84a` | gleich |
| Unterschied `R`..`32d0d84a` | `git diff --numstat d720e44f 32d0d84a`: `Dockerfile` 1/1, `harness/mk/build.mk` 1/1, Plandatei 174/0 (der Eintrag nennt „Plandatei“ ohne Zahl) | gleich |
| 1. Gates | eigener Lauf (oben), Exit 0 an `6630a708` | bestätigt; (b) des Eintrags selbst nicht nachprüfbar, siehe F-3 |
| 2. Gegenproben nach `CO-002` | beide engsten Fälle grün (oben) | bestätigt, die Grenze (Mutationen und übrige Gegenproben nicht gelaufen) steht im Eintrag |
| 3. Host-Gegenproben nach `CO-002` | `git diff --numstat -M 35778a2f 32d0d84a -- <sechs Pfade>`: 1/1, 4/4, 1/1; `-U0` zeigt je sechs Zeilen wie im Eintrag; `harness/mk/container-start-gegenprobe.mk` Zeile 8 nennt `bash tools/dev/harness/container-start-gegenprobe.sh`, die Datei liegt dort | bestätigt; die Zeile `../..` zu `../../..` ist Code, siehe F-4 |
| 4. `VERSION`, Zeile | `git show 32d0d84a:VERSION` = `0.2.0`; Zeile 2307, drei Zellen | gleich |
| 5. Trockenlauf | `40072517299`: `workflow_dispatch`, `main`, `headSha` `32d0d84a…`, 07:01:38Z bis 07:04:36Z, `success`; Job `homebrew` übersprungen; Log: `release: eingaben: grün für v0.2.0 (Trockenlauf)`, Hinweis zur Zeile 0.2.0 ist die Zeile selbst, `Index-Digest sha256:1d501812…` | gleich |
| 6. Scan | `40072517299`: `release: scan: 0 Befund(e) über der Schwelle, 0 ohne Behebung` (07:04:32Z); `40074219261`: dieselbe Zeile (07:31:23Z) | gleich |
| 7. Ja des Nutzers | nicht prüfbar, die Wortlaute „release“ und „ja“ sind Wiedergabe der Koordination; der Eintrag sagt das selbst | als Grenze belassen, F-5 |
| 8. Lauf des Tags | `40074219261`: `push`, `headBranch` `v0.2.0`, `headSha` `32d0d84a…`, `run_attempt` 1, 07:27:44Z bis 07:32:11Z, `success`; Jobs `release` (07:27:46Z bis 07:32:00Z) und `homebrew` (07:32:02Z bis 07:32:11Z); Zeilen `eingaben`, `binaries`, `image`, `veroeffentlichen` (zweimal), `vergleichen` (Digest `1d501812…` in beiden), `smoke` (zweimal), `latest` (zweimal), `github-release angelegt`; null Zeilen `release: hinweis:` | gleich |
| Registries | `docker buildx imagetools inspect` ohne Anmeldung: `ghcr.io` und `docker.io`, je `0.2.0` und `latest`: `sha256:1d501812…`; `0.1.1` bleibt auf `624c099b…` | gleich |
| Job `homebrew` | `git ls-remote …homebrew-pgwire-recorder.git refs/heads/main` = `d96e29684bd0b9d70dcec87a60c0746230fcb41c`; Commit 2026-10-24 07:32:07Z, Nachricht `pgwire-recorder 0.2.0`; `gh secret list`: `HOMEBREW_TAP_GITHUB_TOKEN` 2026-10-20T07:41:39Z, unverändert | gleich |
| Formel im Tap | `version "0.2.0"`, vier Adressen `…/v0.2.0/pgwire-recorder-<os>-<arch>`; die Summen sind die der `SHA256SUMS` des Release (darwin-arm64 `7e43a883…`, darwin-amd64 `745313db…`, linux-arm64 `8e397889…`, linux-amd64 `4177b046…`) | gleich |
| 9. Anhänge | `gh release view v0.2.0`: kein Entwurf, keine Vorabversion, `publishedAt` 07:31:59Z, genau die sieben Namen | gleich |
| Binary | `SHA256SUMS` heruntergeladen; `sha256sum -c --ignore-missing`: `pgwire-recorder-linux-amd64: OK` (`4177b046…`); `version`, `--version`, `-v`: je `pgwire-recorder 0.2.0` | gleich |
| 10. Beobachtungen | zwei `release: skopeo:`-Zeilen mit `manifest unknown` (07:31:26Z, 07:31:33Z); Image `Created` 2026-10-24T07:01:05Z gleich Committer-Zeit von `32d0d84a` (`09:01:05+02:00`); `docker run` des Image gibt `pgwire-recorder 0.2.0` | gleich |
| 11. Bezug | `da5fe8b6` ändert allein `README.md` (3/3) und `docs/user/benutzerhandbuch.md` (15/15); Kopf, §2 und README nennen `0.2.0` und den Digest `1d501812…`, `0.1.1` kommt außer der Zeile der Änderungshistorie nirgends vor | gleich; inert nach `SPEC-052` Punkt 12 |
| Beleg-Datei | `sha256sum`: `9ca8a58d9b46213601164a78b177053265b3a28a944e4cf2f2fe7f051dc316bc`, gleich §7; Kopf (Datum, Stand, Plattform, Image mit Digest), Befehle (a) bis (e), Exit 0, eine Zeile `level=WARN code=PGR-W2001` und sonst `level=INFO`: erfüllt die Regel „welche Zeilen zählen“ aus `SPEC-067` Punkt 20 (3); eine gekürzte Stelle mit Marke und Grund; die Datei wurde seit dem Anlegen nicht geändert (`git log` nennt einen Commit) | gleich |
| Probe der Zeile 0.2.0 | am veröffentlichten Binary: `--help` nennt `--version, -v` und `--diff-all`; Replay-Server und `play` (selbst gefahren): `--compare-responses` gibt Exit 5 und eine Abweichung („Position 1, Art Befehl“); mit `--diff-all` Exit 5 und „Abweichungen: 2“ mit Position 1 und 2; `--diff-all` allein Exit 2 `PGR-E2001` | gleich dem Eintrag Punkt 4 |

## DoD und Entscheidungs-Konformität

`slice-release-v0-2-0`: Die vier Punkte der DoD vor den Closure-Pflichten (drei Liefer-Punkte und `make gates`) sind erfüllt. (1) Freigabe und Tag: Release-Commit `d720e44f` mit `VERSION` und der Zeile (drei Zellen), Eintrag mit elf Punkten und der Zeile des Jobs `homebrew`, alle Zahlen am Original bestätigt. (2) Bezug: `da5fe8b6`, nur belegte Befehle, Grenze für macOS, Windows und `arm64` im Handbuch genannt. (3) `releasing.md` in eigenem Commit `2b0869ff` mit Gate-Lauf danach, Version 0.17, Zeilen 0.15 bis 0.17. (4) `make gates` und `CO-002`-Ersatz: bestätigt (oben). Offen sind die Closure-Pflichten (Review-Report liegt vor, Closure-Notiz, Register, Risiko-Ausgänge, Verbrauch von `CO-002` samt Index und `git mv`). Alle sieben Risiken in §6 tragen „offen bis zur Closure“; das ist vor der Closure zulässig.

Entscheidungen: `SPEC-067` Punkt 18 („ein Implementer startet keinen Tag“) gehalten: Der Tag stammt von der Koordination, der Implementer hat nach Eintrag nur gelesen. „Jeder weitere Tag“: Der Bezug ersetzt den von `v0.1.1`, die Zeile der Änderungshistorie bleibt. Punkt 20 (6): Beleg als Datei mit Summe in §7. `releasing.md` §10: elf Punkte plus Zeile `homebrew`, Carveout im Beleg mit Ersatz, Hash, Datum, Dauer, Last, Ausgang und Grenze. `AGENTS.md` §3.11: Der Satz „die Formel trägt die höchste stabile Version“ ist enger gefasst (`3bc2f931`), die Formel selbst nennt die Adresse von `v0.2.0` (bestätigt). §3.14: Metrik am Endstand 100 Zeilen (Pfade `VERSION docs/user docs/maintainer README.md`, Basis `193f0e0b`), am Stand `da5fe8b6` 75 wie in §7; Marke 1500 nicht erreicht. §3.10: kein Code, Gate oder Vertrag geändert, `git diff` über `Dockerfile`, `harness/mk`, `tools`, `internal`, `cmd` ab `R` bis `main` zeigt allein die zwei Pin-Zeilen.

`slice-go-toolchain-1-27-3`: Liefer-Punkt 1 erfüllt (`[x]`): `Dockerfile:9` und `harness/mk/build.mk:15` tragen denselben Pin `golang:1.27.3-alpine@sha256:b8c9b82a…505e9c` (`git grep` an `32d0d84a`); der alte Digest `f92b6ef800e4` kommt außerhalb `docs/` nicht mehr vor; der Digest ist heute erneut gegen die Registry gelesen gleich. Gate-Punkt erfüllt (528 s an `143e9994` laut §7, mein Lauf an `main` Exit 0). Die Marke 1500 ist mit 4 Zeilen (`Dockerfile harness/mk`, Basis `7cfad5ab`) nicht berührt. Liefer-Punkt 2 (Scan) ist nicht abgehakt und in §7 des Pin-Slice nicht belegt: F-1.

## Bewusstes Brechen (Modul 11): „Scan-Befund weg“

Der Pin hat keinen Test und braucht keinen (`AGENTS.md` §3.10: Text, keine Zusage an einer Grenze). Die Behauptung „der Scan-Befund ist weg“ trägt der Vergleich zweier Läufe, deren Eingaben sich nur im Pin unterscheiden.

**Der rote Lauf war aus dem richtigen Grund rot.** Lauf `40071373625` (`workflow_dispatch`, `main`, `d720e44f`): `Eingaben`, `Binaries` und `Image` laufen grün (Index-Digest `3be2ee12…`), der Schritt `Scan` meldet (06:46:52Z) je Plattform `schwachstelle: linux/amd64 CVE-2026-111943 stdlib v1.27.2 -> 1.26.10, 1.27.3 HIGH` und dieselbe Zeile für `linux/arm64`, dann `2 Befund(e) über der Schwelle, 0 ohne Behebung` und `exit code 1`; `homebrew` übersprungen. Der Befund trägt die eingebaute Version (`v1.27.2`) und nennt `1.27.3` als Behebung; der Lauf ist also der Gegenfall, „alter Pin, derselbe Scanner, derselbe Tag“, aus der Standardbibliothek und nicht aus Datenbank, Netz oder Schwelle.

**Der grüne Lauf ist aus dem richtigen Grund grün.** Lauf `40072517299` (`32d0d84a`): Das Image-Log lädt `golang:1.27.3-alpine@sha256:b8c9b82a…` (13 Treffer auf `1.27.3`); der Index-Digest ist ein anderer (`1d501812…`, das Image hat sich also geändert); der Schritt `Scan` meldet keine Zeile `schwachstelle` und `0 Befund(e) über der Schwelle, 0 ohne Behebung`. Der Unterschied der beiden Stände ist `Dockerfile` (1 Zeile), `harness/mk/build.mk` (1 Zeile) und die Plandatei.

**Eigene Sonde am Ausgeliefertes** (nicht im Eintrag): Das Binary `pgwire-recorder-linux-amd64` von `v0.2.0` trägt die Zeichenkette `go1.27.3` (zweimal) und kein `go1.27.2`; das Binary von `v0.1.1` trägt `go1.27.2` (zweimal). Die Binaries entstehen aus demselben Pin wie das Image, der Wechsel der Toolchain ist damit am Ausgelieferten belegt, unabhängig vom Scanner.

**Grenzen des Belegs:** Es gibt keinen dritten Lauf mit altem Pin auf dem neuen Stand und keinen Test, der den Pin hält. Der rote Lauf ersetzt den ersten, aber die Datenbank des Scanners ist zwischen den Läufen (18 Minuten) nicht protokolliert; eine veränderte Datenbank kann den Befund nicht entfernen, weil `1.27.2` im Image nicht mehr vorkommt. Dass kein Sensor die zwei Pin-Stellen abgleicht, ist im Plan als „nicht übernommen“ geführt (§1, §6 (d)); beide Stellen tragen heute denselben Text (`git grep`).

## Befunde

### HIGH

Keine.

### MEDIUM

**F-1 (`slice-go-toolchain-1-27-3`): Der Scan-Beleg fehlt in §7 des Slice, der ihn verlangt; Closure-Bedingung.** DoD 2 und §5 (b) verlangen Lauf-Kennung, Hash und Zeile „in §7“; §7 nennt nur: „DoD 2 (Scan) belegt die Koordination über den Trockenlauf; nicht abgehakt“. Die Belege stehen in `freigaben/0.2.0.md` Punkt 5 und 6 und in §6 des Release-Slice (`40072517299` auf `32d0d84a`, `40074219261`). Sie sind richtig (oben bestätigt). Vor der Closure trägt §7 des Pin-Slice die Kennungen `40071373625` (rot), `40072517299` (grün) und die Zeile, und DoD 2 wird abgehakt; erst dann darf die Datei nach `done/`. Außerdem fehlt noch `evidence/slice-go-toolchain-1-27-3.md` im Register (`ls` der Ablage: drei Dateien, nicht vier), DoD-Punkt der Closure.

**F-2 (`slice-release-v0-2-0`): Der Lerneintrag und die Ausgänge der Risiken stehen aus; die Closure ist vor allem Handlung des Planners.** Das ist kein Mangel der Lieferung. Die Verifikation bestätigt dafür die Tatsachen: Risiko 6 (Scan, Pin) ist *eingetreten* mit Ausgang `slice-go-toolchain-1-27-3` (Adresse existiert, liegt in `in-progress/`, nimmt den Befund an und führt ihn); Risiko 1 (Ersatzbelege) ist *entfallen*, die Ersatzbelege sind grün; Risiko 2 (Token) ist *entfallen* für `v0.2.0`: der Job `homebrew` schrieb im ersten Versuch, das Geheimnis ist unverändert seit 2026-10-20, der Beleg bleibt Existenz und der Lauf, und die Zählung von `BEO-REPO/vorbedingung-geheimnis-nur-existenz-geprueft` führt der Planner. Risiko 7 (Fehler im Skript, Folge-Tag) ist *entfallen*: Der Lauf war im ersten Versuch grün. Risiko 5 (`docs/maintainer/` nicht inert) ist *entfallen* durch den Schnitt (zwei Commits, bestätigt). Risiko 3 (Zeile sagt mehr als das Binary) ist *entfallen*: Die Probe am veröffentlichten Binary ist von mir wiederholt und gleich. Risiko 4 (Kopf des Handbuchs `0.1.1`) ist *entfallen* mit `da5fe8b6`.

### LOW

**F-3: Punkt 1 (b) des Eintrags ist nicht auflösbar.** Der Eintrag nennt für den Lauf an `32d0d84a` kein Log im Repo, keine Dauer und keine Last und stützt ihn auf (a) und den Diff `143e9994..32d0d84a` (nur Planung); das steht dort als Grenze. Mein Lauf an `6630a708` (283 s, Exit 0) deckt den Hash nicht selbst, deckt aber den Code (kein Unterschied außerhalb von Doku) und schließt die Lücke inhaltlich. Als Form bleibt: `releasing.md` §10 verlangt „Hash, Datum“, `CO-002` für die Punkte 2 und 3 mehr; für Punkt 1 gilt die Form erfüllt.

**F-4: Punkt 3, Zeile (2) ist Code, kein Pfadname.** `../..` wird `../../..` in `container-start-gegenprobe.sh` (Berechnung der Wurzel des Repos). Das Kriterium von `CO-002` („nur Pfadzeilen“) nennt eine Änderung der Logik ausdrücklich kein Ersatz, der Nutzer entscheidet vor dem Tag. Der Eintrag nennt die Zeile offen und gibt den Nutzer als Quelle der Zustimmung („ja“) an; die Wiedergabe ist von der Koordination und nicht nachprüfbar (wie Punkt 7). Die Sache selbst ist richtig: Der Umzug von `tools/harness/` nach `tools/dev/harness/` verlangt genau diese Stufe, und die Gegenprobe ist am Tag-Commit mit dem Fall `p16-keine-kette` lauffähig (`make -n` nennt den neuen Pfad). Dass der Ersatz „kein Lauf am Host“ ist, steht als Grenze im Eintrag.

**F-5: Punkt 7 (Ja des Nutzers) ist eine Angabe der Koordination.** Beide Wortlaute („release“, „ja“) und ihre Zeitpunkte sind im Repo nicht belegbar; der Eintrag sagt „Wiedergabe der Koordination, nicht nachgefahren“. Das ist die Form, die `SPEC-067` Punkt 20 (6) für eine Angabe ohne Datei vorsieht. Ein Beleg, den ein Verifier nachfahren könnte, gibt es nicht. Bei der Closure bleibt der Punkt eine Zusage, keine Messung; kein Befund gegen den Implementer.

**F-6: Reihenfolge Beleg-Datei und §7.** `SPEC-067` Punkt 20 (6): „Die Datei wird in dem Commit abgelegt, der §7 ihren Pfad nennt.“ Die Datei liegt in `8d1fbd5e`, die Zeile mit Pfad und Summe in §7 steht in `6630a708`, der Eintrag verweist seit `1d0519ef` auf sie. Der Bezug-Commit `da5fe8b6` (09:51:48+02:00) und der Abschnitt des Handbuchs zu Homebrew standen 16 Minuten ohne Datei auf `main` (Datei 10:08:45+02:00); der erste Lauf (07:42:09Z, 7 min 46 s) war nicht gesichert. Der Plan (§3, Reihenfolge Schritt 5) nannte die Datei vor dem Bezug-Commit. Behoben durch den zweiten Lauf und die enger gefasste Aussage; der Zustand ist heute richtig.

### INFO

**I-1:** `releasing.md` §5 Schritt 3 und §10 („getaggter Commit ist der Commit, der `VERSION` und die Zeile trägt“) sagen nicht, dass der Tag auf einem Nachfahren dieses Commits liegen darf (hier `32d0d84a` statt `R`). Der Satz bleibt wahr, weil `32d0d84a` beide trägt; der Eintrag und `CO-002` erklären die Abweichung. Ein Hinweis in `releasing.md` ist Sache des Planners, falls er die Form für künftige Tags festhalten will.

**I-2:** Der Index-Digest des Trockenlaufs auf `32d0d84a` und des Laufs des Tags ist gleich (`1d501812…`), und `Created` des Image gleich der Committer-Zeit: In diesem Fall war der Bau stabil. `SPEC-067` Punkt 5 sagt das nicht zu, und der Eintrag nennt es Beobachtung.

**I-3:** `BEO-REPO/wiederanstoss-eines-roten-release-laufs-unbelegt`: Der Lauf endete grün, der Wiederanstoß eines roten Laufs bleibt unbelegt (im Eintrag als Grenze). Die Zählung führt der Planner.

**I-4:** Der Pin-Plan nennt in §6 den Ausgang des vierten Risikos „weiter offen“ mit Register-Adresse `slice-harness-ci-freshness-images`; der Report des Reviewers deckt zwei Slices unter einem Dateinamen (auch dieser Bericht). Das ist keine Abweichung gegen eine Regel, nur gegen die Zuordnung „ein Bericht je Slice“.

## Plan gegen Diff

| Plan | Diff | Ergebnis |
|---|---|---|
| Release-Commit `R`: `VERSION`, Zeile 0.2.0 | `d720e44f`: `VERSION` und eine Zeile im Handbuch | gleich |
| Tag auf `R` | Tag auf `32d0d84a` (Korrektur dokumentiert in §3 und im Eintrag; Anlass der rote Scan `40071373625`) | Abweichung, begründet und nachgezogen (`AGENTS.md` §3.9) |
| Der Kopf des Handbuchs bleibt `0.1.1` bis nach dem Lauf | Release-Commit lässt den Kopf unberührt; `da5fe8b6` setzt ihn auf `0.2.0` | gleich |
| Bezug-Commit nur `docs/user/`, `README.md`, Planung | `da5fe8b6`: `README.md`, Handbuch | gleich |
| `releasing.md` in eigenem Commit | `2b0869ff`; `80e999fc` für F-5 des Reviews | gleich |
| `make homebrew-nachweis` vor dem Bezug-Commit, Beleg als Datei | Lauf 07:42Z unbelegt, Bezug-Commit 09:51+02:00 (07:51Z), Datei aus zweitem Lauf 08:06Z, Commit 10:08+02:00 | Abweichung der Reihenfolge, F-6 |
| Pin an zwei Stellen, Zählung zwei Zeilen je Stelle | `Dockerfile` 1/1, `build.mk` 1/1; Gesamt 4 Zeilen | gleich |
| Basis `193f0e0b` und `7cfad5ab` | beide sind der reine Move nach `in-progress/` | gleich |
| Marke 1500 | 100 und 4 | nicht berührt |

## Offene Closure-Punkte (für den Planner)

1. `slice-go-toolchain-1-27-3`: §7 mit den Scan-Belegen (F-1), DoD 2 abhaken, Datei `evidence/slice-go-toolchain-1-27-3.md` im Register, Lerneintrag, Risiko-Ausgänge, Paarungen.
2. `slice-release-v0-2-0`: Closure-Notiz mit Lerneintrag, Ausgänge der Risiken (Vorschläge unter F-2), Register (`vorbedingung-geheimnis-nur-existenz-geprueft`, `wiederanstoss-eines-roten-release-laufs-unbelegt`, `release-pipeline-pin-und-werkzeug-altern-ohne-melder`), Index von `docs/plan/carveouts/README.md` und `git mv` von `CO-002` nach `done/` im Commit nach dem Inhalt.
3. Der Punkt aus §1 „macOS und `arm64` nicht übernommen, keine Adresse“ geht an den Planner (Review F-3, im Slice vermerkt).

## Abgleich mit dem Review

Der Review-Bericht gilt dem Stand `511428c9`; seine Befunde F-1 bis F-8 hat die Nacharbeit (`90199364` bis `6630a708`) bearbeitet. Mein Urteil steht auf dem Stand `6630a708`.

| Review | Stand heute (nachgefahren) |
|---|---|
| F-1 (MEDIUM, Beleg-Datei fehlt) | behoben: Datei und Summe stimmen (siehe Tabelle oben); F-6 oben nennt die Reihenfolge der Commits als Rest |
| F-2 (MEDIUM, Plan folgt nicht) | behoben: §1, §3 (Tabelle, Punkt 11, Reihenfolge Schritt 5), §6 (h) und §7 nennen den Nachweis als Lieferung; die Aussage in §7 „§1 schloss ihn aus“ ist als Geschichte der Korrektur gekennzeichnet |
| F-3, F-4 (LOW, Adresse nimmt nicht an) | behoben: beide Pläne führen den Punkt als „nicht übernommen“, keine Sendung |
| F-5 (LOW, `releasing.md` §1 und §7) | behoben: beide Stellen nennen `0.1.1` und `0.2.0` (`80e999fc`) |
| F-6 (LOW, `--diff-all` am Binary) | behoben, und von mir am veröffentlichten Binary wiederholt: gleiche Ausgabe wie der Eintrag |
| F-7 (LOW, Hash in `CO-002`) | behoben als Tatsachenberichtigung im Kopf von `CO-002` (`c224f283`) |
| F-8 (INFO, Log nicht im Repo) | im Eintrag als nicht auflösbar genannt; von mir durch einen eigenen Lauf an `main` gedeckt (F-3 oben) |
| F-9, F-10 (INFO) | F-9 bleibt bis zur Closure offen (hier F-1); F-10 betrifft auch diesen Bericht |

**Was die Verifikation über den Review hinaus zeigt:** (a) die Probe am ausgelieferten Binary (`go1.27.3` statt `go1.27.2`) für die Behauptung „Scan-Befund weg“, die der Review nicht prüfte; (b) die Läufe von `make gates` und der zwei engsten Gegenproben; (c) die Reihenfolge von Beleg-Datei, §7 und Bezug-Commit (F-6 oben, im Review nicht genannt); (d) der Befund F-1 oben zum Scan-Beleg in §7 des Pin-Slice, den der Review als INFO (F-9) führt: Ich stufe ihn als MEDIUM, weil die DoD ihn wörtlich in §7 verlangt und die Datei ohne ihn nicht nach `done/` darf. Dagegen halte ich F-4 und F-5 oben als LOW fest, wo der Review die Zahlen und Zeilen (Punkt 3, Punkt 7) als ohne Befund führt: Beide betreffen Angaben, die kein Verifier nachfahren kann.
