# Review-Report: slice-erster-release-veroeffentlichung-image — 2026-10-19

**Review-Art:** Code — gegen Plan, Vertrag `SPEC-067` (Fassung 1b), [ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md) und die Hard Rules (`AGENTS.md` §3.9 bis §3.14).

**Gegenstand:** slice-erster-release-veroeffentlichung-image, Diff `5f68ebef..94011311` (Zweig `rev-rel-b` = `impl-rel-b`, 9 Commits, Merge der Entscheidungen des Architect eingeschlossen); Metrik-Basis `dc7f0a69`.

**Skill:** `.harness/skills/reviewer.md` @ 229b5efb
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-19

> **Zitier-Form.** Kennung statt Adresse: `slice-<Kennung>`, `make <target>`, Baseline-Stellen als Tag + Pfad in Inline-Code. Das `pfad`-Feld zitiert den Stand des Laufs.

**Eingangs-Kontext:**

- Slice-Plan `slice-erster-release-veroeffentlichung-image` (§1, §3, §6, §7, §8, DoD)
- `SPEC-067` (Fassung 1b, alle Punkte), `SPEC-055` Punkt 7 (Auslagerung des Smokes), [ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md)
- `LH-FA-16`, `LH-QA-03`, `LH-QA-01`
- `AGENTS.md` (Hard Rules, insbesondere §3.9 bis §3.14)
- §7 des Plans (Belege des Implementers); der Bericht des Implementers lag nicht vor

**Engste Sensoren dieses Laufs** (die vollständigen Läufe des Implementers wurden nicht wiederholt): `bash -n` der fünf Skripte; `make workflows-check` und `make docs-check` grün; 27 zusätzliche Mutationen von `release.sh` in einer Wegwerf-Kopie der Gegenprobe (`GP_MUT`, `GP_ENTDECKE=1`; die Kopie ist entfernt, nichts davon liegt im Diff). Die Metrik 1187 (Summe aus Einfügungen und Löschungen, `git diff --numstat dc7f0a69..HEAD` ohne `docs/plan` und `docs/reviews`) ist nachgemessen und stimmt; Meldeschwelle 1300 und Marke 1700 sind nicht erreicht. Es gab keinen Registry- und keinen Netzaufruf.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Zusage „Bau mit dem gepinnten Builder“ hat keinen Fall für das Argument `--builder pgr-release` des Baus: Wird es entfernt, bleibt die Gegenprobe grün (Mutation `x-builder-arg` des Reviews: rote Fälle „nichts“). Der Bau liefe dann im Standard-Builder (Treiber `docker`), der Pin wäre wirkungslos und ein Bau für zwei Plattformen in ein OCI-Layout scheiterte oder liefe auf einer ungepinnten Engine; `p5-aufruf` prüft `--driver docker-container` und `--driver-opt image=` am Aufruf von `buildx create`, aber nicht, dass der Bau diesen Builder wählt. | `AGENTS.md` §3.10 (Je Zusage eine Mutation) | `tools/harness/release.sh` · `bargs=(buildx build --builder pgr-release` | ja — `make release-gegenprobe` mit der Mutation, die `--builder pgr-release` aus `bargs` entfernt | Argument eines Aufrufs ohne Fall (Mutation grün) |
| F-2 | MEDIUM | Die Zusage „nur stabile Versionen setzen `latest`“ (Kriterium „kein `-`“) ist nur mit einer Vorabversion der Form `1.2.0-rc.1` geprüft; ein auf `*-rc*` verengtes Muster bleibt grün (Mutation `x-vorab-rc`: rote Fälle „nichts“). Mit dieser Mutation setzte `v1.2.0-beta.1` in beiden öffentlichen Registries `latest`; das ist die Grenzmutation „Muster verengt“, die der Plan für Zusagen der Form „nur X“ als Lehre aus `slice-erster-release-veroeffentlichung` benennt und für die Bedingung der Vorabversion bei `latest` nicht durchführt. | `AGENTS.md` §3.10 (Zusagen „nur X“: Grenzmutation neben der Entfernung) | `tools/harness/release.sh` · `if [[ $v == *-* ]]; then echo "release: latest:` | ja — `make release-gegenprobe` mit der Mutation `*-*` → `*-rc*` | Zusage „nur X“ ohne Grenzmutation (Muster verengt) |
| F-3 | MEDIUM | Die Zusage „Tag-Satz nur `<Version>` und `latest`“ (Entscheidung des Nutzers vom 2026-10-19) ist nur für den Schritt `veroeffentlichen` mit Fall und Mutation belegt (`p6-tag-satz`, `m-p6-zusaetzlicher-tag`). Ein zusätzlicher Kopierbefehl im Schritt `latest` auf `<MAJOR>.<MINOR>` lässt jeden Fall grün (Mutation `y-latest-tag-satz-extra`: rote Fälle „nichts“); die Fälle `p8-*` prüfen nur die Ziele `:latest`. | `AGENTS.md` §3.10 (Zusage „nur X“, Ortsangabe zählt je Schritt) | `tools/harness/release.sh` · `echo "release: latest: $reg zeigt auf $d"` | ja — `make release-gegenprobe` mit einem zusätzlichen `skopeo copy` auf `:${v%.*}` im Schritt `latest` | Zusage „nur X“ ohne Grenzmutation (zweiter Ort) |
| F-4 | MEDIUM | Die Meldung von `skopeo inspect` wird in `veroeffentlichen` und `vergleichen` in eine Scratch-Datei geschrieben und bei den roten Ausgängen nicht ausgegeben (`rot "skopeo inspect für $reg endet mit Exit $rc"`, `rot "$reg ist ohne Anmeldung nicht lesbar (Exit $rc)…"`). `SPEC-067` Punkt 6 und 9 und der Plan nennen die Wortlaute als nicht gemessen und sagen, „der erste Lauf mit einem Tag zeigt sie“; im Log des Laufs steht der Wortlaut nicht, Betreiber und Architect sehen bei Abweichung nur den Exit-Code, und der Vertrag ist nicht nachzuziehen. `github-release` gibt `gh.err` dagegen aus. | `SPEC-067` Punkt 6 (Grenze) · Maintainability | `tools/harness/release.sh` · `rot "skopeo inspect für $reg endet mit Exit $rc"` | ja — Fall mit Attrappe, die einen unbekannten Wortlaut meldet, und Prüfung, dass er in stderr steht | Zusage („zeigt der Lauf“) ohne Träger im Log |
| F-5 | LOW | Der Wächter `[ -n "$soll" ]` in der Summenprüfung je Plattform hat keinen Fall (Mutation `x-summe-leer`: rote Fälle „nichts“). Ohne ihn gälten ein fehlendes Binary im Image und ein fehlender Eintrag in `SHA256SUMS` zusammen als gleich (beide leer); `P4liste` prüft die Dateien im Verzeichnis, nicht die Zeilen von `SHA256SUMS`. | `AGENTS.md` §3.10 | `tools/harness/release.sh` · `[ -n "$soll" ] && [ "$ist" = "$soll" ]` | ja — Fall mit `SHA256SUMS` ohne Zeile der Plattform und Image ohne Binary | Wächter ohne Fall (Mutation grün) |
| F-6 | LOW | Die Prüfung „Index nennt genau die zwei Plattformen“ im Schritt `vergleichen` hat keine Mutation je Registry: Entfällt sie nur für `docker.io`, bleibt alles grün (Mutation `x-plattformen-nur-ghcr`). Das Ergebnis ist über die folgende Gleichheit der Digests gedeckt, aber §7 des Plans führt weder diese Mutation noch eine Äquivalenz mit Invariante, Ort und Sonde; „Keine äquivalenten Mutanten“ gilt nur für die benannten Mutationen. | `AGENTS.md` §3.10 (Aufzählung zählt je Element) | `tools/harness/release.sh` · `index_pruefen "$scratch/roh-$reg" "der Index aus $reg" # P7plattformen` | ja — `make release-gegenprobe` mit der Mutation `[ "$reg" != docker.io ] \|\| continue` vor dem Aufruf | Element einer Aufzählung (Registry) ohne eigene Mutation |
| F-7 | LOW | „Das Binary `/pgwire-recorder` jedes Manifests hat die Summe …“ (`SPEC-067` Punkt 5, `docs/maintainer/releasing.md` §6 Schritt 3) liest weiter als die Prüfung: Geprüft wird ein zweiter Bau je Plattform mit `type=local` aus demselben Builder, nicht das Layout, das `veroeffentlichen` kopiert. Weicht das Binary im Layout vom Export ab, merkt es die Prüfung nicht. | `AGENTS.md` §3.11 (Zusage nur so weit wie die Prüfung) | `tools/harness/release.sh` · `--output "type=local,dest=$scratch/fs-$arch"` | nein — erst ein Lauf mit echtem Builder zeigt, ob beide Bauten gleich sind | Zusage weiter als Prüfung (anderes Artefakt geprüft) |
| F-8 | LOW | Plan DoD Liefer-Punkt 2 („`github-release` läuft erst nach grünem `vergleichen`“), `SPEC-067` Punkt 9 („Vorbedingung `vergleichen` von B“) und die Zuordnungstabelle des Vertrags sagen `vergleichen`; Code, Fall `p9-nicht-vor-registries`, Punkt 7 und `releasing.md` §6 verlangen `smoke.ok`. Die Marken-Aufzählung in Punkt 13 nennt `smoke.ok` nicht. | `AGENTS.md` §3.9 (Der Slice-Plan folgt jeder Korrektur) | `docs/plan/planning/in-progress/slice-erster-release-veroeffentlichung-image.md` · `A läuft erst nach grünem` | nein — Lesart von Plan und Vertrag; `make kopf-check` prüft sie nicht | Plan/Vertrag folgt dem Code nicht (Vorbedingung) |
| F-9 | INFO | Die Randformen Pin des Builders, Marken, Committer-Zeit, Wortlaute von „Tag fehlt“, Quelle von `POSTGRES_IMAGE`, Layout-Form und beide fehlenden Geheimnisse standen in §6 nicht, als die Arbeit begann; der erste Code-Commit `299c8953` trug sie als vorläufige Werte (Platzhalter-Pin aus Nullen, nur `manifest unknown`). Der Implementer hat sie offengelegt und zurückgegeben, der Architect hat sie in `SPEC-067` entschieden (`5cb3c3ce`), bevor der Code sie umsetzte (`8339d87b`). Kein Verstoß; die Liste des Architect vor dem Code war unvollständig. | `AGENTS.md` §3.12 | `docs/plan/planning/in-progress/slice-erster-release-veroeffentlichung-image.md` · `Entschieden vom Architect am 2026-10-19 nach der Rückgabe des Implementers` | nein | Randform nach der Rückgabe entschieden (Liste vor dem Code unvollständig) |
| F-10 | INFO | Sicherheit, ohne Befund gegen den Vertrag, zur Kenntnis: (a) `Eingaben` bekommt `DOCKERHUB_TOKEN` auch im Trockenlauf; ein `workflow_dispatch` auf einem Zweig mit veränderter `release.yml` liest die Geheimnisse (Folge der Entscheidung für Repository-Geheimnisse, [ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md), `SPEC-067` Punkt 15). (b) `packages: write` und `contents: write` gelten für den ganzen Job, nicht nur für die Schritte mit Anmeldung. (c) `container-smoke.sh` räumt nichts auf, und `release.sh` räumt Container und Netz des Smokes nicht; auf einem Runner folgenlos. | [ADR-0051](../plan/adr/0051-veroeffentlichung-auf-github-actions.md) · Maintainability | `.github/workflows/release.yml` · `RELEASE_DOCKERHUB_TOKEN: ${{ secrets.DOCKERHUB_TOKEN }}` | nein | Reichweite der Geheimnisse (entschieden) |
| F-11 | INFO | Nicht prüfbar ohne Netz und von mir nicht wiederholt: Digest und Tag des Pins `moby/buildkit` (Messung des Architect), das Layout eines Baus für zwei Plattformen im Treiber `docker-container`, die Wortlaute von `skopeo` und `gh`, die Wirkung von `SOURCE_DATE_EPOCH` ohne `ARG`. Der Vertrag führt sie als Grenze; belegt sie erst der Trockenlauf und der erste Tag. Der Trockenlauf auf GitHub ist nicht gefahren (§7). | `SPEC-067` Punkte 5, 6, 12 | `tools/harness/release.sh` · `BUILDKIT_IMAGE='moby/buildkit@sha256:cec9f139` | nein — Trockenlauf auf GitHub | Außenwirkung nur durch Lauf belegbar |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.github/workflows/release.yml` (Text 1b, Bedingungen auf `push`, Geheimnisse nur in `Eingaben`, `Veröffentlichen`, `Latest`; Pin von `actions/checkout` gleich `ci.yml`; `permissions: {}` oben; Zeitgrenze; kein `environment:`) | geprüft, ohne Befund; `make workflows-check` und `make docs-check` grün |
| `tools/harness/release.sh`, Sicherheit (Anmeldung per `--password-stdin` über `printf`, kein `set -x`, kein `env`, kein Wert in Meldungen; Außenwirkung nur bei Push eines Tags; Trockenlauf ohne `skopeo` und `login`) | geprüft, ohne Befund; Mutationen mit dem Token von `ghcr.io` in Meldung und Fehlerpfad werden rot (`y-ghcr-token-*`) |
| `tools/harness/release.sh`, fail-closed (drei Digests aus rohen Bytes, `--no-creds`, vorhandener Tag nie überschrieben, Marken je Schritt und Löschen beim Neubeginn, Wortlaute von „Tag fehlt“ mit Schreibweise, Reihenfolge `ghcr.io` vor `docker.io`, `latest` nur nach `smoke.ok`, höchste stabile Version, `sort -V`, Anker von `KERN`) | geprüft, ohne Befund; Mutationen `x-tag-mit-v`, `x-ghcr-anzeige`, `x-kern-anker`, `x-latest-immer`, `x-fehlt-re-weit`, `x-ghcr-creds-docker`, `x-latest-ohne-anmeldung`, `x-p7tag`, `y-secrets-nur-im-lauf`, `y-zahl-manifeste-mind`, `y-latest-ohne-vergleichen-marke` werden rot |
| `tools/harness/container-smoke.sh`, `run-container.sh`, `container-gegenprobe.sh` (Auslagerung; Text von `SPEC-055` Punkt 7 unverändert; Mutationen auf das neue Skript gerichtet, `p9-smoke-weiter` neu) | geprüft, ohne Befund; `make test-container-gegenprobe` nicht wiederholt (Beleg des Implementers in §7) |
| `tools/harness/release-gegenprobe.sh`, `harness/mk/release.mk`, `harness/README.md` (Zeilenliste Herkunft B, Fälle und Mutationstabelle; Zeile des Sensors) | geprüft; Befunde F-1 bis F-3, F-5, F-6; `bash -n` ohne Fehler; `make release-gegenprobe` nicht wiederholt (390 Läufe laut §7) |
| `docs/maintainer/releasing.md` (Stand 0.8, §1, §2, §4, §6, §8, §11) | geprüft, ohne Befund außer F-7 (Wortlaut „jede Plattform“); Verhaltenssätze haben die Belege der Tabelle in §7 |
| `docs/user/benutzerhandbuch.md`, `README.md` | geprüft, ohne Befund: keine Änderung im Diff, kein Zielstand, keine Chronik |
| Kommentare in Skripten und Fragment (`AGENTS.md` §3.7), Suppression (`nolint`), Moves und Inhalt im selben Commit (§3.3), ADR-Immutabilität, Core-Reinheit | geprüft, ohne Befund |
| Plan (§1 Abgrenzung, §3, §6, §8; Sendungen an `slice-erster-release-freigabe` mit Kennung des Gebers in dessen §1; Zählung Liefer-Punkte drei, Schichten zwei; Metrik 1187 nachgemessen) | geprüft; Befund F-8 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Argument eines Aufrufs ohne Fall (Mutation grün) · Zusage „nur X“ ohne Grenzmutation · Zusage („zeigt der Lauf“) ohne Träger im Log · Wächter ohne Fall (Mutation grün) · Element einer Aufzählung ohne eigene Mutation · Zusage weiter als Prüfung · Plan/Vertrag folgt dem Code nicht · Randform nach der Rückgabe entschieden

## Verdikt

**Merge-blockierend:** ja — vier MEDIUM (F-1 bis F-4), je mit einer in diesem Lauf gefahrenen Mutation oder einem Pfad im Code, der das Versagen zeigt; kein HIGH. F-5 bis F-8 sind LOW, F-9 bis F-11 INFO.

**Übergabe:** Findings gehen an den Implementer; F-8 berührt `SPEC-067` (Stratum 2) und geht zusätzlich an den Architect. Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation — DoD-Konformität prüft der Verifier.
