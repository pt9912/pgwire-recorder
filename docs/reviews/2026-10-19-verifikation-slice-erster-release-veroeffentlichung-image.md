# Verifikation: slice-erster-release-veroeffentlichung-image

**Rolle:** Verifier. **Datum:** 2026-10-19. **Spanne:** Basis `5f68ebef` (Metrik-Basis `dc7f0a69`) bis `45ff2710` (Zweig `ver-rel-b` = `impl-rel-b` nach Review-Runde 1). **Gegenstand:** [slice-erster-release-veroeffentlichung-image](../plan/planning/done/welle-erster-release/slice-erster-release-veroeffentlichung-image.md) gegen DoD, [SPEC-067](../../harness/vertraege/SPEC-067-release.md) (Fassung 1b), [ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md) und die Hard Rules.

Der Review-Report ([2026-10-19-review-slice-erster-release-veroeffentlichung-image](2026-10-19-review-slice-erster-release-veroeffentlichung-image.md)) ist Eingang, nicht Beleg. Der Reviewer hat die Behebungen von Runde 1 nicht erneut gesehen (Beobachtung `nacharbeit-nach-review-ohne-erneuten-lauf`); ich habe F-1 bis F-8 selbst nachgefahren. Es gab keinen Registry-, Netz- oder GitHub-Aufruf, keinen Tag, keinen Push nach GitHub.

**Engste Sensoren in diesem Lauf** (Dev-Container pgr-dev3, Klon `/workspaces/pgr-dev3/pgwire-recorder`): `bash -n` der fünf Skripte (ok); `make docs-check` (1222 Dateien, 0 Befunde); `make workflows-check` (2 Dateien, kein Befund); `make kopf-check` (Exit 0); `GP_MUT='^ok$' bash tools/harness/release-gegenprobe.sh` (unveränderte Kopie, 317 Fälle, grün, 65 s) und 15 Brüche von Hand (unten), nacheinander, nie parallel. Nicht wiederholt: volle `make release-gegenprobe`, `make test-container-gegenprobe`, `make gates` (Belege des Implementers am Endstand in §7; `container-smoke.sh`, `run-container.sh`, `container-gegenprobe.sh` sind seit dem Lauf auf `a4a96c60` unverändert: `git diff a4a96c60..HEAD` über die drei Dateien ist leer).

## Bewusstes Brechen (Modul 11): Fix von Hand entfernt, Fall muss aus dem richtigen Grund rot werden

Jeder Bruch ist ein Textersatz in `tools/harness/release.sh` im Arbeitsbaum, danach `GP_MUT='^ok$'`, danach `git checkout`; `git status` danach sauber. Die roten Fälle prüfen die Sache (Argument am Aufruf, Ziel des Kopierbefehls, Text im Log), nicht nur den Exit.

| Finding | Bruch | rote Fälle | Urteil |
|---|---|---|---|
| F-1 | `--builder pgr-release` aus `bargs` entfernt | `p5-aufruf` (prüft das Argument am Aufruf des Baus) | behoben, richtig rot |
| F-2 | Bedingung `*-*` auf `*-rc*` verengt | `p8-vorabversion` (fährt `-rc.1`, `-beta.1`, `-alpha.1`) | behoben, richtig rot |
| F-3 | zusätzlicher `skopeo copy` auf `:${v%.*}` im Schritt `latest` | `p8-tag-satz` (jedes Ziel trägt `<Version>` oder `latest`, Inhalt beider Registries) | behoben, richtig rot |
| F-4a | Ausgabe der skopeo-Meldung in `veroeffentlichen` entfernt | `p6-inspect-meldung-im-log` | behoben, richtig rot |
| F-4b | dasselbe in `vergleichen` | `p6-inspect-meldung-im-log` | behoben, richtig rot |
| F-4c | Präfix `release: skopeo: ` zu `skopeo: ` | `p6-inspect-meldung-im-log` | behoben, richtig rot |
| F-5 | Wächter `[ -n "$soll" ]` entfernt | `p5-summe-eintrag-fehlt` | behoben, richtig rot |
| F-6 | Plattformprüfung nur für `ghcr.io` (entfällt für `docker.io`) | `p7-eine-plattform` (ein Bein je Registry, prüft die Meldung "hat nicht genau zwei Manifeste") | behoben, richtig rot; "kein äquivalenter Mutant" in §7 stimmt |
| G-1 | `smoke.ok` für `github-release` nicht verlangt | `p9-nicht-vor-registries`, `p11-alte-marke-b` | richtig rot |
| G-2 | `--no-creds` entfernt | `p7-anonym-docker`, `p7-anonym-ghcr`, `p7-gleich`, `p6-inspect-meldung-im-log` | richtig rot |
| G-3 | Anmeldung mit `--password <Token>` statt `--password-stdin` | `p10-stdin` | richtig rot |
| G-4 | `SOURCE_DATE_EPOCH` aus Autor-Zeit (`%at`) | `p5-epoch-committer` | richtig rot |
| G-5 | vorhandener Tag mit anderem Digest wird nicht rot | `p6-tag-vorhanden-anders`, `p10-nicht-im-log-image` | richtig rot |

F-7 und F-8 sind Text, kein Code: siehe unten (V-1, V-2).

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg |
|---|---|---|
| Liefer-Punkt 1, Image bauen und prüfen (Layout mit einem Index und zwei Manifesten, `D`, Summe je Plattform, Geheimnisse vor dem ersten Push; Fälle `p5-*`, `p3-secret-*`) | Code-Hälfte bestätigt; Trockenlauf **offen** | Schritt `image` gelesen gegen Punkt 5: Builder gepinnt (`docker-container`, `--driver-opt image=`), Bau mit `--builder`, `type=oci`, `rewrite-timestamp`, `--provenance=false`, `--sbom=false`, Index-Typ, zwei Manifeste, Plattformen, Summe je `amd64`/`arm64` mit Wächter, `D` aus dem Index-Blob. Brüche F-1, F-5, G-4 rot. Fälle `p3-secret-*` und `p5-*` in der unveränderten Kopie grün. Handbuch und README: keine Änderung im Diff, kein Zielstand (Diffstat enthält keine der beiden Dateien). Lauf-Kennung und Dauer des Trockenlaufs fehlen in §7 (erwartet, siehe Offen). |
| Liefer-Punkt 2, Veröffentlichen, Vergleichen, Smoke (`ghcr.io` zuerst, nie überschreiben, drei Digests aus rohen Bytes, anonym, Plattformen, Smoke aus jeder Registry, `github-release` erst nach `smoke.ok`, `--password-stdin`) | bestätigt (Code, Gegenprobe); Außenwirkung **offen** | Brüche F-3, F-4, F-6, G-1, G-2, G-3, G-5 rot. `container-smoke.sh` als eigenes Skript, Aufruf mit `<Registry>/pt9912/pgwire-recorder@<D>`, zwei Läufe mit verschiedenen Namen. `SPEC-055`-Gate: Dateien seit dem Lauf des Implementers unverändert (Beleg dort 234 s, grün; nicht wiederholt). Keine Ausgabe eines Geheimniswerts: `printf` nach stdin, kein `set -x`, kein `env`, Fehlerpfade von `anmelden` und `zustand_tag` verwerfen die Meldung von skopeo (kein Wert im Log), `p10-nicht-im-log-image` grün. |
| Liefer-Punkt 3, `latest`, Fehlertabelle, Wiederanlauf (nur stabil, nur höchste stabile Version, nur nach Vergleich und Smoke; zweite Registry rot; `releasing.md` §11) | bestätigt | Brüche F-2, F-3 rot; `p8-*` (höchste stabile, `sort -V`, fremder Tag, Vorabversion zählt nicht, `p8-ohne-smoke`, `p8-nach-rotem-vergleich`, `p8-zweite-registry-faellt`) und `p11-*` in der unveränderten Kopie grün. `releasing.md` §11 nennt die Zeilen `veröffentlichen`, `vergleichen`, `smoke`, `latest`. Endfassung `1b` besteht `workflows-check` und `docs-check`. |
| `make gates`; `make release-gegenprobe`, `make test-container-gegenprobe`, `make workflows-check` grün | `workflows-check` von mir bestätigt; übrige Läufe Beleg des Implementers, nicht wiederholt | §7 Tabelle (398 Läufe, 721 s; 234 s; 279 s). Den Endstand des Codes (`9f506869`) nennt §7; der Kopf `45ff2710` ändert danach nur Plan §7 und §8 (`git diff 9f506869..45ff2710 --stat` über Code und Vertrag leer, siehe Offen/Prüfung unten). |
| Review durchgeführt, Report liegt vor | bestätigt | Report vom selben Tag, 0 HIGH, 4 MEDIUM, 4 LOW; Behebungen in Runde 1 von mir nachgefahren (oben). Eine zweite Review-Runde hat nicht stattgefunden (Behebung ohne erneuten Lauf des Reviewers; dieser Bericht ersetzt sie nur für F-1 bis F-8). |
| Closure-Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen | offen (erwartet) | Alle Häkchen leer; §6 trägt die Ausgänge als "offen bis Closure" bzw. "eingetreten, Folge-Slice `slice-erster-release-freigabe`"; §7 ist noch Beleg-Abschnitt, keine Closure-Notiz. |

## Plan-vs-Code-Diff

| Plan §3 Zeile | Code | Urteil |
|---|---|---|
| `release.yml` (`packages: write`, `env` der Geheimnisse, Schritte `Image`, `Veröffentlichen`, `Vergleichen`, `Smoke`, `Latest`) | 28 Zeilen hinzu, keine geändert; Geheimnisse nur in `Eingaben`, `Veröffentlichen`, `Latest`; Bedingung auf `push` bei den vier Schritten mit Außenwirkung; `permissions: {}` oben, Zeitgrenze 60, kein `environment:`, Pin von `actions/checkout` wie `ci.yml` | stimmt |
| `release.sh` (Schritte, `eingaben` (7)) | Schritte `image`, `veroeffentlichen`, `vergleichen`, `smoke`, `latest` und `pruefe_geheimnisse` (Reihenfolge der Gründe: (6) vor (7)) | stimmt; Plan §3 nennt als Zustandsdatei nur `vergleichen.ok`, der Code legt auch `smoke.ok` und `image/digest` an (V-2) |
| `container-smoke.sh` neu; `run-container.sh`, `container-gegenprobe.sh` | 94 neu; 5/52 und 30/29; Mutationen des Smokes auf das neue Skript gerichtet | stimmt |
| `release-gegenprobe.sh`, `release.mk`, `harness/README.md` | 623/50, 5/5, 1/1; gelöscht sind nur Kopfkommentar, Verallgemeinerungen von `p1-reihenfolge-*` und `p1-bedingung-*`, `p1-zusatz-packages-write` | stimmt; die Verallgemeinerung der A-Fälle `p1-reihenfolge-*` und `p1-bedingung-github-release` nennt §7 nicht (V-3) |
| `Dockerfile` keine Änderung | nicht im Diff | stimmt |
| `releasing.md` | 51/27 | stimmt |
| Handbuch, README keine Änderung | nicht im Diff | stimmt |

**Größe** (`AGENTS.md` §3.14): `git diff --numstat dc7f0a69..HEAD -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'` ergibt von mir nachgemessen **1225** (Einfügungen plus Löschungen), wie in §7. Marke 1700 nicht erreicht (72 %); die Meldeschwelle 1300 nicht erreicht. Die Schätzung des Architect (1150 bis 1500) trifft. Eine Meldung war nicht nötig.

## Findings

| ID | Schwere | Befund | Ort |
|---|---|---|---|
| V-1 | LOW | F-7 ist nur in `docs/maintainer/releasing.md` §6 berichtigt. DoD Liefer-Punkt 1 des Plans sagt weiter "das Binary `/pgwire-recorder` jedes Manifests hat die Summe". Geprüft ist ein zweiter Bau je Plattform (`type=local`, derselbe Builder), nicht das Layout, das `veroeffentlichen` kopiert; ob beide Bauten gleich sind, belegt erst ein Lauf mit echtem Builder. `SPEC-067` Punkt 5 nennt den Weg in einer Klammer (Export aus demselben Builder) und ist damit enger gefasst als die DoD. | Plan §2 Liefer-Punkt 1; `AGENTS.md` §3.11 |
| V-2 | LOW | Plan §3 (Zeile `release.sh`) nennt als Zustandsdatei nur `vergleichen.ok`; der Code legt zusätzlich `smoke.ok` und `image/digest` an und `latest`/`github-release` verlangen `smoke.ok`. §3 ist die einzige Stelle des Plans, die F-8 (Marke) nicht nachzieht; Vertrag, DoD, §6 und §7 sind einheitlich `smoke.ok`. | Plan §3; `AGENTS.md` §3.9 |
| V-3 | INFO | §7 sagt, die Fälle von A blieben unverändert; im Diff sind `p1-reihenfolge-*` und `p1-bedingung-github-release` von A in allgemeinere Fälle (je Schritt) umgebaut, `p1-zusatz-packages-write` entfällt (das nennt §7). Namen und rote Mengen bleiben, die Gegenprobe ist grün. | `release-gegenprobe.sh` |
| V-4 | INFO | `Smoke`, `Veröffentlichen`, `Vergleichen` und `Latest` laufen im Trockenlauf nicht (Bedingung auf `push`). Der Trockenlauf belegt also den Bau und das Layout, nicht den Smoke aus einer Registry; das ist in §6 (Risiko "Der Trockenlauf belegt die Außenwirkung nicht") und `SPEC-067` Punkt 12 so zugesagt. | `release.yml` |

Kein HIGH, kein MEDIUM. F-1 bis F-6 behoben und mit Fix-Entzug rot gesehen; F-7 Text berichtigt in `releasing.md` (Rest V-1); F-8 in Vertrag, DoD, Code und §6 einheitlich (Rest V-2); F-9 bis F-11 unverändert INFO.

## Offen, nicht erledigt (ausdrücklich)

1. **Trockenlauf auf GitHub**: nicht gefahren. Er braucht den Merge nach `main` (`workflow_dispatch` nur für eine Datei auf dem Standard-Zweig). Die Secrets `DOCKERHUB_USERNAME` und `DOCKERHUB_TOKEN` sind laut Auftrag gesetzt und das Repository auf Docker Hub existiert; belegt ist das von mir nicht. Lauf-Kennung und Dauer fehlen in §7 und gehören in die Closure (DoD Liefer-Punkt 1).
2. **Pin `moby/buildkit` v0.33.1** (`sha256:cec9f139…3dea`): Messung des Architect; ich habe sie ohne Netz nicht wiederholt. Dass der Bau mit ihm läuft, zeigt erst der Trockenlauf.
3. **Layout bei zwei Plattformen** im Treiber `docker-container`: nicht gemessen (Messung nur für eine Plattform im Treiber `docker`); jede andere Form als ein Index-Blob mit zwei Manifesten ist rot (Grenze, `SPEC-067` Punkt 5).
4. **Wortlaute von `skopeo` und `gh`** (`manifest unknown`, `name unknown`, `not found`, `release not found`): nicht gemessen; im Dev-Container sind beide Werkzeuge nicht installiert, die Gegenprobe läuft gegen Attrappen. Seit F-4 steht der echte Wortlaut mit `release: skopeo: ` im Log des ersten Laufs.
5. **Wirkung von `SOURCE_DATE_EPOCH`** ohne `ARG` in der Stufe `runtime`: Grenze, keine Zusage (`SPEC-067` Punkte 5 und 15); der Trockenlauf zeigt sie.
6. Außenwirkung (Schreibrecht der Token, Sichtbarkeit des Pakets, Smoke aus den Registries, Anlegen des Releases): erst der Lauf für den ersten Tag (`slice-erster-release-freigabe`).
7. Volle `make release-gegenprobe`, `make test-container-gegenprobe` und `make gates` am Endstand: Belege des Implementers, von mir nicht wiederholt. Der Code nach `9f506869` ist bis `45ff2710` unverändert, soweit §7 es sagt; die Übereinstimmung von Kopf und Endstand habe ich über die unveränderte Kopie der Gegenprobe (grün auf `45ff2710`) gestützt, nicht über einen zweiten `make gates`.

## Urteil

DoD-Konformität: die Code- und Gegenprobe-Hälfte der Liefer-Punkte 1 bis 3 ist bestätigt; die Behebungen F-1 bis F-6 tragen (je Bruch aus dem richtigen Grund rot). Offen sind erwartungsgemäß die Außenwirkungs-Belege (Trockenlauf, erster Tag) und die Closure. Die Findings V-1 und V-2 sind Text-Nachzüge für den Implementer (Plan §2 und §3), kein Merge-Hindernis; V-3 und V-4 sind zur Kenntnis. Spec-Lücke oder Vertragsänderung ist nicht nötig.
