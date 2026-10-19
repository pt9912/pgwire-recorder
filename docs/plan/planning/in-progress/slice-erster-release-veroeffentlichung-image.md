# Slice slice-erster-release-veroeffentlichung-image: Veröffentlichung des Images

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-erster-release.

**Bezug:** [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus)

**Berührte Spec-Stellen:** `LH-FA-16.a` · `SPEC-031` · `SPEC-035` · `SPEC-055` (Punkt 7: der Smoke wird ausgelagert, der Text von `SPEC-055` bleibt) · `SPEC-063` und `SPEC-064` (gelesen: Gate der Workflows und Konventionen der CI) · `SPEC-067` (Fassung B: Punkte 1b, 2 (Ergänzung), 3 (7), 5, 6, 7, 8, 9 (Vorbedingung), 10 (Registries), 11 (Zeilen der Registries), 12 (Ergänzung), 13 bis 15 (Ergänzung); Zuordnung in der Tabelle des Vertrags)

**Verantwortlich:** pt9912 (Zweig impl-rel-b)

**Autor:** pt9912. **Datum:** 2026-10-19.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Tag `v<SemVer>` baut das Docker/OCI-Image einmal für `linux/amd64` und `linux/arm64` in ein OCI-Layout, veröffentlicht genau dieses Layout als **eine** Manifestliste mit gleichem Index-Digest in `ghcr.io/pt9912/pgwire-recorder` und `docker.io/pt9912/pgwire-recorder`, vergleicht die Digests aus den rohen Bytes des Index ohne Anmeldung, lässt das Image aus beiden Registries den Smoke aus Abnahmeszenario 9 bestehen und setzt `latest` nur nach einem stabilen Tag; ist der Vergleich rot, endet der Lauf, bevor `latest` oder der GitHub-Release entsteht. Der Trockenlauf (`workflow_dispatch`) baut und prüft das Image, ohne etwas in eine Registry zu schreiben. Das ist Fassung B des Vertrags [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md): `release.yml` in der Endfassung `1b`, die die Fassung von `slice-erster-release-veroeffentlichung` um Zeilen ergänzt, ohne eine ihrer Zeilen zu ändern oder zu streichen.

**Übernimmt:** `slice-erster-release-veroeffentlichung` — aus dem Schnitt vom 2026-10-19 (dort §1, *Abgegeben*, und §4; der Schnitt in A und B ist eine Empfehlung des Planners auf Anweisung der Koordination und vom Nutzer nicht bestätigt, die Teilung der Fassungen in [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) ist Entscheidung des Architect): das Image und alles, was an ihm hängt. Bindend für diesen Slice aus den Entscheidungen des Nutzers vom 2026-10-09 (Geber `welle-erster-release`, über jenen Slice):

1. `:latest` zeigt nur nach einem stabilen Tag auf das neue Image; eine Vorabversion ändert `:latest` nicht (`SPEC-067` Punkt 8).
2. Der Spiegel auf Docker Hub ist fail-closed: Beide Registries tragen denselben Index-Digest der Manifestliste, sonst ist der Release rot (`SPEC-067` Punkte 6 und 7).

Dazu aus dem Schnitt: die Geheimnisse `DOCKERHUB_USERNAME` und `DOCKERHUB_TOKEN` und die Anmeldung (`SPEC-067` Punkte 3 (7) und 10), die Image-Zeilen der Fehlertabelle in `docs/maintainer/releasing.md` §11 und der Wiederanlauf (Punkt 11), die Image-Anteile von `docs/maintainer/releasing.md` (§2, §7, §8), und der Smoke als eigenes Skript `tools/harness/container-smoke.sh` ([`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 7), das `make test-container` und dieser Lauf gemeinsam rufen. Aus `slice-doku-ist-stand` (dort §6, *Installation*), soweit der Slice sie belegt: Handbuch §2 und `README.md` zu den Images, siehe DoD 1 und die Abgabe unten.

**Builder-Image** (Zuordnung aus `slice-harness-ci-freshness-images` §1, 2026-10-19): Das Image des Builders (Treiber `docker-container` von `buildx`, per Digest gepinnt, `SPEC-067` Punkt 5) gehört diesem Slice. Der Pin ist Teil von Liefer-Punkt 1 (Schritt `Image`), kein weiterer Punkt.

**Voraussetzung `slice-erster-release-veroeffentlichung`** (Geber der Fassung `1a` von `release.yml`, von `release.sh` mit den Schritten `eingaben`, `binaries`, `github-release`, von `VERSION`, der Gegenprobe `release-gegenprobe` und ihrer Zeilenliste): B hängt seine Zeilen an die Liste an und lässt die Fälle von A unverändert (`SPEC-067` Punkt 1); der Slice liegt deshalb in `done/`, bevor dieser beginnt. Dazu die Voraussetzungen jenes Slice (`slice-harness-ci-workflows`, `slice-harness-ci`, `slice-v1-abschluss-container` und `-container-doku`).

**Bindung aus `slice-erster-release-veroeffentlichung` — Pins und Härtung** (Planner, 2026-10-19): Die Endfassung `1b` fällt unter `make workflows-check` und `make docs-check`, trägt dieselben Pins, Rechte und Zeitgrenzen wie Fassung `1a` und `ci.yml` (Konventionen aus `slice-harness-ci-workflows` und `slice-harness-ci`), dazu `packages: write` im Job. Die Geheimnisse stehen nur in den Schritten, die sie lesen (`SPEC-067` Punkt 1, 10).

**Sendung an `slice-erster-release-freigabe` — Einhängepunkt des Scans** (Planner, 2026-10-19, `AGENTS.md` §3.13): Der Scan des Images läuft zwischen den Schritten `Image` und `Veröffentlichen` und liest das Layout aus Punkt 5; `D` bleibt unverändert. Die Freigabe ändert dazu den Text von `release.yml` und `SPEC-067` Punkte 1 und 5 (dort §1 *Übernimmt* und DoD 1, mit der Kennung dieses Slice). Dieser Slice hält die Stelle frei: Der Schritt `Image` endet mit dem Layout auf der Platte, der Schritt `Veröffentlichen` liest genau dieses Layout.

**Sendung an `slice-erster-release-freigabe` — Handbuch zum Bezug der Images** (Planner, 2026-10-19): Der Abschnitt des Handbuchs §2 *Installation* zu den Images aus `ghcr.io` und `docker.io` (Bezug, Tags `<Version>` und `latest`, Plattformen) steht erst, wenn der erste echte Tag ihn belegt (`AGENTS.md` §3.11); der Trockenlauf und die Gegenprobe belegen keinen Bezug aus einer Registry. Dieser Slice schreibt in Handbuch und `README.md` nur, was Gegenprobe, `make test-container` und Trockenlauf belegen (DoD 1).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Tag, `VERSION`, Eingaben bis auf die Geheimnisse, Binaries, GitHub-Release, der Trockenlauf selbst und die Fassung `1a` — `slice-erster-release-veroeffentlichung` (Voraussetzung, oben); dieser Slice ändert keine Zeile von A.
- Image-Scan, Freshness-Audit, Freigabe-Checkliste, der Bezug-Abschnitt des Handbuchs und der erste echte Tag — `slice-erster-release-freigabe` (Sendung oben): Der Scan hängt an der Wahl des Scanners, die Freigabe setzt diese Pipeline voraus.
- Die Homebrew-Formel und ihr Nachweis — `slice-erster-release-homebrew-nachweis` (`welle-v1-nachtrag`).
- Ein Probe-Tag — nicht vorgesehen (Entscheidung des Nutzers vom 2026-10-19, [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md)): Er legte ein Image einer Vorabversion in öffentliche Registries; der Trockenlauf baut und prüft das Image und schreibt nichts (`SPEC-067` Punkt 2).
- Konten, Namensräume, das Repository auf Docker Hub, die Sichtbarkeit und die Verknüpfung des Pakets in `ghcr.io` und das Hinterlegen der Geheimnisse im Repository — Handlungen des Betreibers (`docs/maintainer/releasing.md` §4 Punkt 4); dieser Slice liest sie nur, und die Prüfung der Folge (anonymer Zugriff) ist Punkt 7. Die Zugangsdaten stehen nie im Repository.
- Reproduzierbarkeit des Images — `LH-QA-01` verlangt sie nur für das Binary; das Image wird über seinen Digest identifiziert, `SOURCE_DATE_EPOCH` ist ein Versuch, kein Versprechen (`SPEC-067` Punkte 5 und 15).
- Eine Änderung des Textes von `SPEC-055` — der Smoke wird ausgelagert, sein Verhalten bleibt (Auslagerung, kein neuer Vertrag); `make test-container` und `make test-container-gegenprobe` bleiben grün und unverändert in ihren Fällen.
- Tags `<MAJOR>` und `<MAJOR>.<MINOR>`, Löschen oder Zurücksetzen eines Tags in einer Registry — nicht vorgesehen (Entscheidung des Nutzers vom 2026-10-19: Tag-Satz nur `<Version>` und `latest`; Tags werden nicht verändert).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1, Image bauen und prüfen** ([`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), `SPEC-067` Punkte 5 und 3 (7), 12 (Ergänzung)): Der Schritt `Image` baut für `linux/amd64` und `linux/arm64` in **ein** OCI-Layout mit genau einem Index und zwei Manifesten, berechnet den Index-Digest `D`, und das Binary `/pgwire-recorder` jedes Manifests hat die Summe des Eintrags für seine Plattform in `SHA256SUMS`; sonst Exit 1 mit Plattform. Fehlt eines der Geheimnisse `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`, endet `eingaben` vor dem ersten Push (Test: Fälle `p5-*`, `p3-secret-*` mit Mutationen, `AGENTS.md` §3.10). Der Trockenlauf fährt den Schritt `Image` einmal auf dem Stand der Lieferung und ist in §7 belegt (Lauf-Kennung, Dauer). Handbuch §2 *Container* und `README.md` sind im Ist-Zustand des gebauten Binaries berichtigt, soweit der Slice etwas liefert: Plattformen des Images, Prüfung des Binaries im Image; kein Bezug aus einer Registry (Abgabe an `slice-erster-release-freigabe`, §1), ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11).
- [ ] **Liefer-Punkt 2, Veröffentlichen, Vergleichen, Smoke** ([`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), `SPEC-067` Punkte 6, 7, 9 (Vorbedingung), 10 (Registries), [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 7): Der Push kopiert das Layout ohne Neubau in beide Registries (`ghcr.io` zuerst), ein vorhandener Tag wird nie überschrieben; der Vergleich der drei Digests `D`, `ghcr.io`, `docker.io` aus den rohen Bytes des Index, der anonyme Zugriff und die Plattformen sind fail-closed; der Smoke (`tools/harness/container-smoke.sh <Referenz>`, ausgelagert aus dem Smoke von `make test-container`, dessen Gate und Gegenprobe unverändert grün bleiben) besteht aus **jeder** Registry; der Schritt `github-release` von A läuft erst nach grünem `smoke` (Marke `smoke.ok`, `p9-nicht-vor-registries`). Anmeldung nur per `--password-stdin`, kein Wert eines Geheimnisses im Log (Test: Fälle `p6-*`, `p7-*`, `p10-*` ohne `GITHUB_TOKEN`-Anteil von A und `p9-nicht-vor-registries`, je Zusage und Element eine Mutation; `make test-container-gegenprobe` unverändert grün). Dass die Token schreiben dürfen, das Paket öffentlich ist und der Smoke aus den Registries besteht, belegt erst der Lauf für den ersten Tag (`SPEC-067` Punkt 12), nicht dieser Slice.
- [ ] **Liefer-Punkt 3, `latest`, Fehlertabelle, Wiederanlauf** ([`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), `SPEC-067` Punkte 8 und 11): `latest` wird nur gesetzt, wenn die Version stabil ist, Punkt 7 für beide Registries grün war und die Version höher ist als jede andere stabile Version unter den Tags `v*`; eine Vorabversion und der Patch eines älteren Zweigs lassen es unverändert; scheitert die zweite Registry, ist der Lauf rot (Test: Fälle `p8-*`, `p11-*` mit Mutationen). `docs/maintainer/releasing.md` §11 trägt die übrigen Zeilen der Fehlertabelle (`veröffentlichen`, `vergleichen`, `smoke`, `latest`) und den Wiederanlauf; §2, §7 und §8 sind für das Image auf das ausgeführte Verfahren gebracht (die Formel ist keine Lieferung des Release, Kontrollen nach dem Push), auf dem Stand aufbauend, den A liefert. Die Endfassung `1b` von `release.yml` besteht `make workflows-check` und `make docs-check`, die Gegenprobe hält die Zeilenliste mit der Herkunft B je Zeile; sie ist Akzeptanzkriterium der Punkte 1 bis 3, kein vierter Liefer-Punkt.
- [ ] `make gates` grün; `make release-gegenprobe`, `make test-container-gegenprobe` und `make workflows-check` (die betroffenen Gegenproben, `AGENTS.md` §6 Schritt 6) grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/release.yml` (Endfassung `1b`, `SPEC-067` Punkt 1) | update | fügt hinzu: `packages: write`, den `env`-Block der Geheimnisse in `Eingaben`, die Schritte `Image`, `Veröffentlichen`, `Vergleichen`, `Smoke`, `Latest`; Zeilen von A bleiben |
| `tools/harness/release.sh` (Schritte `image`, `veroeffentlichen`, `vergleichen`, `smoke`, `latest`; `eingaben` (7)) | update | Bau in ein OCI-Layout, `skopeo copy --all --preserve-digests`, Digest aus den rohen Bytes, `latest`, Zustandsdatei `vergleichen.ok` |
| `tools/harness/container-smoke.sh` | neu | Auszug des Smokes aus [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 7 (Image-Referenz als Parameter) |
| `tools/harness/run-container.sh`, `tools/harness/container-gegenprobe.sh` | update | `run-container.sh` ruft das ausgelagerte Skript statt eigenem Text (Punkt 7: Plattform des Docker-Servers und Tag bleiben dort); die Gegenprobe kopiert `container-smoke.sh` und richtet die Mutationen des Smokes (Punkt 7, `p2-dev`, `p9-volume-replay`) auf das neue Skript; Verhalten und Fälle von `SPEC-055` unverändert (`make test-container`, `make test-container-gegenprobe` grün). `harness/mk/container.mk` bleibt unverändert |
| `tools/harness/release-gegenprobe.sh`, `harness/mk/release.mk`, `harness/README.md` (Zeile von `make release-gegenprobe`) | update | Fälle `p2-*` (B), `p3-secret-*`, `p5-*`, `p6-*`, `p7-*`, `p8-*`, `p9-nicht-vor-registries`, `p10-*` (Registries), `p11-*` (B), Zeilenliste von `release.yml` mit Herkunft B; Attrappen von `docker` (`buildx`), `skopeo`, `gp-layout` und des Smokes; Fälle von A bleiben unverändert (die Fixture trägt die Geheimnisse und legt `vergleichen.ok` vor `github-release`; `p1-zusatz-packages-write` entfällt, weil Punkt 1b die Zeile verlangt); Kommentar und Hilfetext in `release.mk`, die Zeile in `harness/README.md` |
| `Dockerfile` | keine Änderung | `SOURCE_DATE_EPOCH` geht als Build-Argument in den Bau (`SPEC-067` Punkt 5), ohne `ARG` in der Stufe `runtime`; ob die Zeitstempel des Images damit gesetzt sind, zeigt erst der Trockenlauf |
| `docs/maintainer/releasing.md` §1, §2, §4, §6, §7 (Satz zur Formel in §2), §8, §11 | update | Image-Anteile des ausgeführten Verfahrens, übrige Zeilen der Fehlertabelle |
| `docs/user/benutzerhandbuch.md`, `README.md` | keine Änderung | Plattformen des Images und Prüfung des Binaries im Image stehen im Ist-Zustand (Handbuch §2 *Container*, `make test-container`); der Bau in ein OCI-Layout und die Veröffentlichung sind Verfahren der Maintainer (`docs/maintainer/releasing.md`), und der Bezug aus einer Registry gehört `slice-erster-release-freigabe` (§1) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-erster-release-veroeffentlichung` liegt in `done/` (Fassung `1a` von `release.yml`, `release.sh` mit den drei Schritten, `VERSION`, Gegenprobe mit Zeilenliste und der Trockenlauf belegt). Schritt 5 der Reihenfolge in §5 von [welle-erster-release](../welle-erster-release.md). Die Randformen aus §6 sind entschieden ([`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md), `AGENTS.md` §3.12); der Implementer hält an und gibt eine nicht genannte Randform an den Architect zurück.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar. Schnitt dann: *Bau des Images* (Liefer-Punkt 1) gegen *Veröffentlichen, Vergleichen, `latest`* (Punkte 2 und 3); die Endfassung `1b` müsste dann in zwei Stufen entstehen, was der Architect vorher in `SPEC-067` teilt. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'` liegt über **1700**, Tests und Gegenprobe eingeschlossen; Basis: dc7f0a69 (HEAD vor dem Anspruch). Schätzung des Architect in dieser Metrik (2026-10-19, nach `SPEC-067`, Fassung B): **1150 bis 1500** (Gegenprobe mit Attrappen für `docker buildx`, `skopeo`, `gh` das Doppelte des Produktcodes); Mitte 1325, 78 % der Marke. Meldeschwelle **1300**: Wer sie überschreitet, meldet es vor dem nächsten Abschnitt, und die Entscheidung („weiter“ oder „schneiden“) steht als datierte Zeile in §7. Faktor Tests zu Produktcode: 2, Beleg `slice-harness-ci-workflows` (Messung 899 gegen eine Schätzung von 545 bis 695, Drift-Log der Roadmap, 2026-10-18); das Drift-Log der Roadmap nennt für `slice-harness-ci` den Faktor 1,3 bis 1,65 der Unterschätzung des Architect, deshalb liegt die Obergrenze der Schätzung nah an der Marke.
- `in-progress` → `open` (blockiert — Carveout?): Der Trockenlauf scheitert an Einstellungen des Repositorys (Actions, Rechte), die der Betreiber setzt, oder `slice-erster-release-veroeffentlichung` liegt nicht in `done/` — extern klären beziehungsweise zurück an den Planner.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Zugangsdaten und Einstellungen für `ghcr.io` und `docker.io` (Geheimnisse `DOCKERHUB_USERNAME` und `DOCKERHUB_TOKEN`, Repository auf Docker Hub, Sichtbarkeit und Verknüpfung des Pakets in `ghcr.io`) sind externe Betreiber-Handlungen und kommen nur vom Nutzer — **Ausgang:** offen bis Closure.
- **Größe.** Schätzung 1150 bis 1500 gegen die Marke 1700 (Mitte 78 %): Der Slice liegt nah an der Marke, weil Bau, Veröffentlichung, Vergleich und Smoke je Zusage eine Mutation tragen — **Ausgang:** offen bis Closure.
- **Der Trockenlauf belegt die Außenwirkung nicht.** Schreibrecht der Token, Paket-Sichtbarkeit, Repository auf Docker Hub, Smoke aus den Registries und das Anlegen des GitHub-Release zeigt erst der Lauf für den ersten Tag (`SPEC-067` Punkt 12); ein Fehler dort ist nach Punkt 11 behebbar, ohne einen Tag zu verändern — **Ausgang:** eingetreten, Folge-Slice `slice-erster-release-freigabe` (dort §6, *Abbruch beim ersten Tag*).
- **Das Image ist nicht reproduzierbar zugesagt** (`SPEC-067` Punkte 5 und 15): Ein Wiederanlauf eines teilweise veröffentlichten Tags kann daran scheitern und endet dann rot mit dem Weg „neuer Tag“ — **Ausgang:** offen bis Closure.
- **Werkzeuge des Runners** (`skopeo`, `gh`, `buildx`) kommen in der Version des Runner-Images; eine Änderung kann das Verhalten ändern (`SPEC-067` Punkt 15). Der Trockenlauf belegt ihre Verfügbarkeit zum Zeitpunkt der Lieferung — **Ausgang:** offen bis Closure.
- **Auslagerung des Smokes.** Der Auszug aus dem Gate `make test-container` darf dessen Verhalten nicht ändern; das belegt `make test-container-gegenprobe` mit den unveränderten Fällen von [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) — **Ausgang:** offen bis Closure.
- **Pin des Builder-Images.** `moby/buildkit@sha256:cec9f139…3dea` (`# v0.33.1`) ist eine Messung des Architect vom 2026-10-19 und altert (`SPEC-067` Punkt 5): Der Pin ist eine Handlung bei der Pflege, kein Wert des Vertrags; ob der Trockenlauf mit ihm baut, zeigt erst der Lauf — **Ausgang:** offen bis Closure.
- **Der Scan** (`slice-erster-release-freigabe`) ändert `release.yml` und `SPEC-067` Punkte 1 und 5; hält dieser Slice die Stelle nicht frei (Layout auf der Platte zwischen `Image` und `Veröffentlichen`), muss er das Layout oder `D` ändern — **Ausgang:** eingetreten, Folge-Slice `slice-erster-release-freigabe` (Sendung in §1).

**Randformen** (`AGENTS.md` §3.12) — **entschieden vom Architect am 2026-10-19** (vom Nutzer am 2026-10-19 angenommen: [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md); Ort: [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md), Randform → Punkt in der Tabelle dort; Anteil von B):

- *Manifestliste und Image-Bau:* ein Bau in ein OCI-Layout, Kopie ohne Neubau in beide Registries; Neubau statt Binary aus `dist/`, Summe je Plattform geprüft; Scan zwischen Bau und Push — Punkt 5.
- *Registry-Namen und Tags:* `<Version>` ohne `v` und `latest`; kein `<MAJOR>.<MINOR>`; anonymer Zugriff wird geprüft — Punkt 6, 7.
- *Reihenfolge der Schritte:* `eingaben`, `binaries`, `image`, (Scan), `veröffentlichen`, `vergleichen`, `smoke`, `latest`, `github-release`; die erste Außenwirkung ist der Push nach `ghcr.io` — Punkt 2, 6.
- *Zugangsdaten:* `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN` als Repository-Geheimnisse (Entscheidung des Nutzers vom 2026-10-19), `GITHUB_TOKEN` mit `packages: write`; fehlendes Geheimnis ist ein Eingabefehler vor dem ersten Push; kein Wert im Log; Anmeldung nur per `--password-stdin` — Punkt 3 (7), 10.
- *`latest`:* nur stabil, nur die höchste stabile Version, nach grünem Vergleich und Smoke; Vorabversion und Patch eines älteren Zweigs ändern es nicht; Beispiel und Regel in `LH-FA-16.a` — Punkt 8.
- *Fail-closed:* Digest aus den rohen Bytes des Index, drei Werte gleich; derselbe Layout-Inhalt wird kopiert; vorhandener Tag mit anderem Digest ist rot, nie überschrieben — Punkt 6, 7.
- *Wiederanlauf und Fehlertabelle (B):* Zeilen `veröffentlichen`, `vergleichen`, `smoke`, `latest`; ein Git-Tag wird nie verändert — Punkt 11.
- *Trockenlauf (B):* Bau und Prüfung des Images, Geheimnisse gesetzt; keine Anmeldung, kein Push — Punkt 2, 12.
- *Plattform:* `ubuntu-24.04`, `linux/amd64` gestartet und geprüft, `linux/arm64` nach Struktur, Architektur und Summe — Kopf des Vertrags, Punkt 15 (Grenze); nach `AGENTS.md` §3.12 *Die Plattform gehört dazu*.
- *Aufrufarten von `make`, Abbruch, Reihenfolge der Gründe:* Punkt 13, 11 (`release.sh` läuft nie unter `make`).
- *Stelle des Scans:* zwischen `Image` und `Veröffentlichen`, der Scan liest das Layout, `D` bleibt unverändert — Punkt 5.

**Entschieden vom Architect am 2026-10-19 nach der Rückgabe des Implementers** ([`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkte 5, 6, 7, 8, 9, 11): (1) Pin des Builders `moby/buildkit@sha256:<Index-Digest von v0.33.1>` mit Tag-Kommentar `# v0.33.1` (gemessen 2026-10-19: `sha256:cec9f139f45e93c5c69c60f8b07cfad9f43f4ef6b6a6cd917527fea5ff2e3dea`), ein Platzhalter ist ungültig, die Gegenprobe prüft Form und Platzhalter; (2) `POSTGRES_IMAGE` aus `harness/mk/integration.mk`, fehlt es Exit 2 `vorbedingung`; (3) Vorbedingungen je Schritt als Marken, Exit 2; (4) Marke `smoke.ok`, `latest` und `github-release` verlangen sie (nicht nur der Workflow-Text); (5) `SOURCE_DATE_EPOCH` ist die Committer-Zeit; (6) „Tag fehlt“ bei `manifest unknown`, `name unknown`, `not found`, sonst rot, nicht gemessen (Grenze); (7) nur stabile SemVer-Tags zählen für `latest`, `latest` wird immer kopiert; (8) Form des Layouts bei zwei Plattformen: Grenze, der Trockenlauf zeigt sie, jede andere Form ist rot; (9) Wirkung von `SOURCE_DATE_EPOCH` ohne `ARG`: Grenze, keine Zusage; (10) Limit anonymer Pulls von PostgreSQL im Smoke: rot, kein Anmelden, Wiederanlauf; (11) beide Geheimnisse fehlend: eine Meldung nennt beide Namen.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

Wird bei Closure gefüllt (vor dem `git mv` nach `done/`).

### Belege des Implementers

Stand: Kopf des Codes `2ee49303` auf dem Zweig `impl-rel-b` (Basis `5f68ebef`; erste Fassung `299c8953`, Auszug des Smokes `f95a320e`, die Entscheidungen des Architect `5cb3c3ce` umgesetzt in `8339d87b`); Läufe im Dev-Container pgr-dev3 (Docker im Container), nacheinander am Endstand. `make test-container-gegenprobe` lief auf `a4a96c60` allein grün (234 s); seither ändert sich weder `container-smoke.sh` noch `run-container.sh` noch `container-gegenprobe.sh`. Ein früherer gleichzeitiger Lauf beider Gegenproben meldete dort die Signalfälle `p9-*-signal` rot (Last).

| Lauf | Ergebnis | Dauer |
|---|---|---|
| `make release-gegenprobe` | `release-gegenprobe: grün (390 Läufe)`: die unveränderte Kopie, 387 Mutationen, der Rekursionsschutz und die Probe der Fallliste; 314 Fälle (161 `p1-*`, 81 `p2-*` bis `p13-*` von A, 66 `p2-*` bis `p11-*` von B, 6 `p13-*` des Abbruchs) | 645 s |
| `make test-container-gegenprobe` | `container-gegenprobe: gruen — Punkte 1 bis 9, 11 und 12 …` (auf `a4a96c60`), Fälle von `SPEC-055` unverändert, Mutationen des Smokes auf `container-smoke.sh` gerichtet, neu `p9-smoke-weiter` | 234 s |
| `make workflows-check` | Exit 0, `workflows-check: 2 Datei(en) geprüft, kein Befund` | 1 s |
| `make docs-check` | Exit 0, `d-check: 1221 Datei(en) geprüft, 0 Befund(e)` | 2 s |
| `make gates` | Exit 0 (darunter `test-container`, `test-postgres-versionen`, `lint`, `kopf-check`, `docs-check`) | 296 s |

Während der Arbeit nur die engsten Sensoren: `bash -n`, Teilläufe der Gegenprobe mit `GP_MUT` und `GP_ENTDECKE` (die roten Fälle je Mutation gemessen, dann die Erwartung im Skript gesetzt und gegen die Vermutung gehalten; Abweichungen sind in der Tabelle, etwa die Fälle, die einen Fehler der Meldung prüfen), `make workflows-check`, `make docs-check`, zuletzt die Läufe oben je einmal am Endstand. `make gates-schnell-gegenprobe` ist nicht gefahren: Seit `5f68ebef` ändert sich weder ein Fragment noch eine Liste (`GEGENPROBEN_CONTAINER`).

**Messung der Größe** (`AGENTS.md` §3.14): `git diff --numstat dc7f0a69..2ee49303 -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'` ergibt als Summe aus Einfügungen und Löschungen **1187** (`release-gegenprobe.sh` 643, `release.sh` 188, `container-smoke.sh` 94, `container-gegenprobe.sh` 59, `releasing.md` 78, `run-container.sh` 57, `release.yml` 28, `SPEC-067` 28 aus den Entscheidungen des Architect, `release.mk` 10, `harness/README.md` 2); ohne den Vertragstext 1159. Die Commits zwischen `dc7f0a69` und `5f68ebef` ändern nur Pläne. Vor den Entscheidungen des Architect waren es 1094. Schätzung des Architect 1150 bis 1500, Meldeschwelle 1300 und Marke 1700 sind nicht erreicht. Die Gegenproben tragen 702 von 1187 Zeilen (Faktor zu `release.sh`, `container-smoke.sh`, `run-container.sh` und `release.yml`, zusammen 367, rund 1,9).

**Mutationstabelle** (`AGENTS.md` §3.10): Zusage · Mutation (Textersatz in einer Kopie von `release.sh` bzw. in `pruefe_release`; Namen wie in `tools/harness/release-gegenprobe.sh`) · roter Fall. Je Mutation sind genau die genannten Fälle rot, jede Mutation einmal gesehen (`GP_ENTDECKE=1`), die unveränderte Kopie ist grün. Bei mehr als sechs roten Fällen nennt die Zeile die Zahl und drei Namen; die vollständige Menge steht in der Ausgabe von `GP_ENTDECKE=1`. Die Mutationen von A (Tabelle in `done/slice-erster-release-veroeffentlichung.md` §7) sind unverändert und mit unveränderten roten Mengen grün; ihre Fälle laufen in der Gruppe `sh`, die der Mutationen von B in der Gruppe `b`.

Zeilenliste von `release.yml` (Punkt 1, Herkunft B): 28 der 66 Zeilen sind Herkunft B (Zeilen 18, 32 bis 34 und 38 bis 61 des Textes). Der Fall `p1-zusatz-packages-write` von A und seine Mutation entfallen in der Endfassung: Punkt 1b verlangt die Zeile `packages: write` (Herkunft B, Zeile 18), und die Mutation „Zeilen mit diesem Wort überlesen“ färbte mit der geforderten Zeile jeden Fall. Ein zusätzliches Vorkommen bricht die Zahl der Zeilen (`m-p1-zusatzzeilen-erlaubt`), ohne eigenen Fall.

| Zusage (Punkt von SPEC-067) | Mutation | roter Fall |
|---|---|---|
| 1: jede der 66 Zeilen steht in der Datei | Zeile n ist optional (`m-p1-zeile-n-optional`, n = 1 bis 66) | `p1-fehlt-n`; bei den fünf Zeilen mit `if:` (41, 49, 52, 55, 63) auch `p1-bedingung-<Schritt>` |
| 1: jede Zeile hat genau den Wortlaut | Inhalt der Zeile n wird nicht verglichen (`m-p1-zeile-n-nicht-verglichen`, n = 1 bis 66) | `p1-geaendert-n`; bei der Zeile des Pins auch die fünf `p1-pin-*` |
| 1: jeder Schritt mit Außenwirkung trägt die Bedingung auf `push` | die Zeile `if:` fehlt, je Schritt (`Veröffentlichen`, `Vergleichen`, `Smoke`, `Latest`, `GitHub-Release`) | `p1-bedingung-veroeffentlichen`, `-vergleichen`, `-smoke`, `-latest`, `-github-release` (Fälle) mit den Mutationen `m-p1-zeile-n-optional` für n = 41, 49, 52, 55, 63 |
| 1: Reihenfolge der Schritte, je Paar | vertauschte Schritte werden angenommen (`m-p1-reihenfolge-<Paar>-tauschbar`, Paare: `eingaben-binaries`, `binaries-release`, `binaries-image`, `image-veroeffentlichen`, `veroeffentlichen-vergleichen`, `vergleichen-smoke`, `smoke-latest`, `latest-release`) | `p1-reihenfolge-<Paar>` |
| 2: `latest` lehnt ohne Push ab | die Prüfung entfällt im Schritt (`m-p2-latest-ohne-ereignispruefung`) | `p2-aussenwirkung-ohne-push-b`, `p2-ref-kein-tag-b` |
| 2: `smoke` lehnt ohne Push ab | die Prüfung entfällt im Schritt (`m-p2-smoke-ohne-ereignispruefung`) | `p2-aussenwirkung-ohne-push-b`, `p2-ref-kein-tag-b` |
| 2: der Trockenlauf meldet sich nicht an | `image` ruft im Trockenlauf `docker login` (`m-p2-trockenlauf-image-meldet-an`) | `p2-trockenlauf-image` |
| 2: der Trockenlauf bis `image` ruft `skopeo` nie | `image` ruft im Trockenlauf `skopeo` (`m-p2-trockenlauf-image-ruft-skopeo`) | `p2-trockenlauf-image` |
| 2: `vergleichen` lehnt ohne Push ab | die Prüfung entfällt im Schritt (`m-p2-vergleichen-ohne-ereignispruefung`) | `p2-aussenwirkung-ohne-push-b`, `p2-ref-kein-tag-b` |
| 2: `veroeffentlichen` lehnt ohne Push von `refs/tags/v…` ab (Exit 2, vor skopeo, gh, Smoke) | die Prüfung des Ereignisses entfällt im Schritt (`m-p2-veroeffentlichen-ohne-ereignispruefung`) | `p2-aussenwirkung-ohne-push-b`, `p2-ref-kein-tag-b` |
| 3 (7): Prüfung der Geheimnisse (Lauf und Trockenlauf) | der Aufruf entfällt (`m-p3-7-geheimnisse-entfernt`) | 7 Fälle, darunter `p10-nicht-im-log-image`, `p3-schluesselwort-geheimnis`, `p3-secret-fehlen-beide` |
| 3: Reihenfolge der Gründe (6) vor (7) | (6) und (7) vertauscht (`m-p3-reihenfolge-6-7-vertauscht`) | `p3-reihenfolge-6-7` |
| 3: Meldung beginnt mit `release: geheimnis: ` | das Schlüsselwort weicht ab (`m-p3-schluesselwort-geheimnis-weicht-ab`) | `p3-schluesselwort-geheimnis`, `p3-secret-fehlen-beide`, `p3-secret-fehlt-token`, `p3-secret-fehlt-user`, `p3-secret-leer-token`, `p3-secret-leer-user` |
| 3 (7): fehlen beide Geheimnisse, nennt die Meldung beide Namen | nur der erste Name (`m-p3-secret-nur-der-erste`) | `p3-secret-fehlen-beide` |
| 3 (7): `DOCKERHUB_TOKEN` gesetzt | die Prüfung entfällt (`m-p3-secret-token-entfernt`) | `p3-schluesselwort-geheimnis`, `p3-secret-fehlen-beide`, `p3-secret-fehlt-token`, `p3-secret-leer-token` |
| 3 (7): `DOCKERHUB_TOKEN` nicht leer | ein leerer Wert gilt als gesetzt (`m-p3-secret-token-leer-gilt-als-gesetzt`) | `p3-secret-leer-token` |
| 3 (7): `DOCKERHUB_USERNAME` gesetzt | die Prüfung entfällt (`m-p3-secret-user-entfernt`) | `p10-nicht-im-log-image`, `p3-secret-fehlen-beide`, `p3-secret-fehlt-user`, `p3-secret-leer-user` |
| 3 (7): `DOCKERHUB_USERNAME` nicht leer | ein leerer Wert gilt als gesetzt (`m-p3-secret-user-leer-gilt-als-gesetzt`) | `p3-secret-leer-user` |
| 3 (7): die Meldung nennt nie einen Wert | die Meldung trägt den Wert des Tokens (`m-p3-secret-wert-in-meldung`) | `p10-nicht-im-log-image`, `p3-secret-fehlt-user`, `p3-secret-leer-user` |
| 5/11: ein scheiternder Bau ist rot | der Exit von `buildx build` wird überlesen (`m-p5-bau-fehler-ignoriert`) | `p5-bau-scheitert` |
| 5/11: ein scheiternder Builder ist rot | der Exit von `buildx create` wird überlesen (`m-p5-builder-fehler-ignoriert`) | `p5-bau-scheitert` |
| 5: Pin des Builders `moby/buildkit@sha256:<64 Hexziffern>` | der Pin trägt einen Tag statt eines Digests (`m-p5-builder-ohne-digest`) | `p5-builder-pin-form` |
| 5: der Pin trägt den Tag als Kommentar | der Kommentar `# v0.33.1` fehlt (`m-p5-builder-ohne-tag-kommentar`) | `p5-builder-pin-form` |
| 5: das Image des Builders ist der Pin aus `release.sh` | `--driver-opt image=…` entfällt (`m-p5-builder-pin-nicht-uebergeben`) | `p5-aufruf` |
| 5: ein Platzhalter aus Nullen ist ungültig | der Digest besteht aus Nullen (`m-p5-builder-platzhalter`) | `p5-builder-platzhalter` |
| 5: `D` ist der Hash des Index-Blobs | `D` aus dem Hash des Archivs (`m-p5-digest-aus-archiv`) | 25 Fälle, darunter `p10-stdin`, `p11-alte-marke-b`, `p11-wiederanlauf-ghcr-vorhanden` |
| 5: das Layout enthält genau einen Index | die Prüfung entfällt (`m-p5-ein-index-entfernt`) | `p5-layout-form` |
| 5: `SOURCE_DATE_EPOCH` ist die Committer-Zeit | die Autor-Zeit (`%at`) (`m-p5-epoch-autor-zeit`) | `p5-epoch-committer` |
| 5: der Eintrag des Layouts ist ein Index | die Prüfung der Art des Eintrags entfällt (`m-p5-index-typ-entfernt`) | `p5-layout-form` |
| 5: der Index nennt genau zwei Manifeste für die Plattformen | der Aufruf von `index_pruefen` im Schritt entfällt (`m-p5-manifeste-entfernt`) | `p5-plattform-fehlt-amd64`, `p5-plattform-fehlt-arm64`, `p5-zwei-manifeste` |
| 5: `type=oci` | `type=docker` (`m-p5-output-nicht-oci`) | `p5-aufruf`, `p5-epoch-committer` |
| 5: `linux/amd64` verlangt | nur `linux/arm64` verlangt (`m-p5-plattform-amd64-nicht-verlangt`) | `p5-plattform-fehlt-amd64` |
| 5: `linux/arm64` verlangt | nur `linux/amd64` verlangt (`m-p5-plattform-arm64-nicht-verlangt`) | `p5-plattform-fehlt-arm64` |
| 5: Bau für `linux/amd64,linux/arm64` | `--platform` entfällt (`m-p5-plattformen-entfallen`) | `p5-aufruf` |
| 5: `--provenance=false` | das Argument entfällt (`m-p5-provenance-entfaellt`) | `p5-aufruf` |
| 5: `rewrite-timestamp=true` | das Element der Ausgabe entfällt (`m-p5-rewrite-timestamp-entfaellt`) | `p5-aufruf` |
| 5: `--sbom=false` | das Argument entfällt (`m-p5-sbom-entfaellt`) | `p5-aufruf` |
| 5: `--build-arg SOURCE_DATE_EPOCH=<Commit-Zeit>` | das Argument entfällt (`m-p5-source-date-epoch-entfaellt`) | `p5-aufruf`, `p5-epoch-committer` |
| 5: Summe des Binaries für `linux/amd64` | die Prüfung entfällt für `amd64` (`m-p5-summe-amd64-nicht-geprueft`) | `p11-alte-marke-b`, `p5-summe-amd64` |
| 5: Summe des Binaries für `linux/arm64` | die Prüfung entfällt für `arm64` (`m-p5-summe-arm64-nicht-geprueft`) | `p5-summe-arm64` |
| 5: `--target runtime` | das Argument entfällt (`m-p5-target-entfaellt`) | `p5-aufruf` |
| 5: Treiber `docker-container` | Treiber `docker` (`m-p5-treiber-nicht-docker-container`) | `p5-aufruf` |
| 5: `--build-arg VERSION=<Version>` | das Argument entfällt (`m-p5-version-entfaellt`) | `p5-aufruf` |
| 5: genau zwei Manifeste | die Prüfung der Zahl entfällt (`m-p5-zahl-entfernt`) | `p5-zwei-manifeste`, `p7-eine-plattform` |
| 5: Commit-Zeit des getaggten Commits | die Zeit von `HEAD` (`m-p5-zeit-aus-head`) | `p5-aufruf` |
| 6: `--all` | das Argument entfällt (`m-p6-all-entfernt`) | `p6-all-und-digests` |
| 6: ein vorhandener Tag mit anderem Digest ist rot | die Meldung entfällt, der Lauf endet grün (`m-p6-anderer-digest-angenommen`) | `p10-nicht-im-log-image`, `p6-tag-vorhanden-anders` |
| 6: eine gescheiterte Anmeldung ist rot | der Exit von `skopeo login` wird überlesen (`m-p6-anmeldung-fehler-ignoriert`) | `p10-nicht-im-log-image`, `p6-anmeldung-scheitert` |
| 6: nur drei Meldungen heißen „Tag fehlt“, jede andere ist rot (Grenze: Muster leer) | jeder Fehler gilt als „fehlt“ (`m-p6-fehler-wie-fehlend`) | `p6-inspect-anderer-fehler`, `p7-anonym-docker`, `p7-anonym-ghcr` |
| 6/11: ein vorhandener Tag mit Digest `D` wird nicht erneut kopiert | kopiert trotzdem (`m-p6-gleicher-tag-kopiert`) | `p11-wiederanlauf-ghcr-vorhanden`, `p6-tag-vorhanden-gleich` |
| 6/11: ein vorhandener Tag mit Digest `D` wird übersprungen (nicht rot) | Exit 1 bei gleichem Digest (`m-p6-gleicher-tag-rot`) | `p11-wiederanlauf-ghcr-vorhanden`, `p6-tag-vorhanden-gleich` |
| 6: nur drei Meldungen heißen „Tag fehlt“ (Grenze) | Muster geweitet um `denied` (`m-p6-muster-geweitet`) | `p6-inspect-anderer-fehler` |
| 6: `--preserve-digests` | das Argument entfällt (`m-p6-preserve-digests-entfernt`) | `p6-all-und-digests` |
| 6: Reihenfolge `ghcr.io`, dann `docker.io` | Schleife `docker.io ghcr.io` (`m-p6-reihenfolge-vertauscht`) | `p6-reihenfolge`, `p6-zweite-registry-faellt` |
| 6: die Schreibweise der Meldung gilt nicht | `grep` ohne `-i` (`m-p6-schreibweise-beachtet`) | `p6-tag-fehlt-not-found` |
| 6: ein vorhandener Tag wird nie überschrieben | bei anderem Digest wird kopiert (`m-p6-vorhandener-tag-ueberschrieben`) | `p10-nicht-im-log-image`, `p6-tag-vorhanden-anders` |
| 6: `manifest unknown` heißt „Tag fehlt“ | der Wortlaut wird nicht erkannt (`m-p6-wortlaut-manifest-unknown-nicht-erkannt`) | 37 Fälle, darunter `p10-nicht-im-log-image`, `p10-stdin`, `p11-alte-marke-b` |
| 6: `name unknown` heißt „Tag fehlt“ | der Wortlaut wird nicht erkannt (`m-p6-wortlaut-name-unknown-nicht-erkannt`) | `p6-tag-fehlt-name-unknown` |
| 6: `not found` heißt „Tag fehlt“ | der Wortlaut wird nicht erkannt (`m-p6-wortlaut-not-found-nicht-erkannt`) | `p6-tag-fehlt-not-found` |
| 6: Tag-Satz nur `<Version>` | zusätzlicher Tag `<MAJOR>.<MINOR>` (`m-p6-zusaetzlicher-tag`) | 7 Fälle, darunter `p11-wiederanlauf-ghcr-vorhanden`, `p6-reihenfolge`, `p6-tag-fehlt-manifest-unknown` |
| 6: die Meldung der zweiten Registry nennt den Zustand | der Zustand entfällt (`m-p6-zustand-fehlt`) | `p6-zweite-registry-faellt` |
| 6: scheitert die zweite Registry, ist der Lauf rot | der Exit von `skopeo copy` wird überlesen (`m-p6-zweite-registry-ignoriert`) | `p10-nicht-im-log-image`, `p11-wiederanlauf-ghcr-vorhanden`, `p6-zweite-registry-faellt` |
| 7: Lesen ohne Anmeldung (`--no-creds`) | `--no-creds` entfällt (`m-p7-anonym-mit-anmeldung`) | `p7-anonym-docker`, `p7-anonym-ghcr`, `p7-gleich` |
| 7: Digest aus den rohen Bytes | der Digest aus der Anzeige `skopeo inspect` (`m-p7-digest-aus-anzeige`) | 21 Fälle, darunter `p10-stdin`, `p11-alte-marke-b`, `p11-wiederanlauf-latest` |
| 7: Digest `docker.io` gleich `D` | der Vergleich entfällt (`m-p7-docker-vergleich-entfernt`) | `p7-digest-docker` |
| 7: Digest `ghcr.io` gleich `D` | der Vergleich entfällt (`m-p7-ghcr-vergleich-entfernt`) | `p10-nicht-im-log-image`, `p11-alte-marke-b`, `p7-digest-ghcr` |
| 7: `vergleichen.ok` nur bei grünem Vergleich | die Marke entsteht vor den Prüfungen (`m-p7-marke-zu-frueh`) | 7 Fälle, darunter `p11-alte-marke-b`, `p7-anonym-docker`, `p7-anonym-ghcr` |
| 7: der Index der Registry nennt genau die zwei Plattformen | der Aufruf entfällt (`m-p7-plattformen-nicht-geprueft`) | `p7-eine-plattform` |
| 7: fehlt `POSTGRES_IMAGE` in `integration.mk`, Exit 2 `vorbedingung` | die Prüfung entfällt (`m-p7-postgres-quelle-fehlt-angenommen`) | `p7-postgres-fehlt` |
| 7: ein roter Smoke ist rot | der Exit des Smokes wird überlesen (`m-p7-smoke-fehler-ignoriert`) | `p7-smoke-docker`, `p7-smoke-ghcr`, `p7-smoke-ok-marke` |
| 7: die Läufe der zwei Smokes tragen verschiedene Namen | derselbe Name (`m-p7-smoke-lauf-gleich`) | `p7-smoke-aufruf` |
| 7: `smoke.ok` nur bei Grün | die Marke entsteht vor dem Smoke (`m-p7-smoke-marke-bei-rot`) | `p7-smoke-ok-marke` |
| 7: `smoke` legt bei Grün `smoke.ok` an | die Marke wird nicht angelegt (`m-p7-smoke-marke-entfaellt`) | 11 Fälle, darunter `p10-stdin`, `p11-wiederanlauf-latest`, `p7-smoke-ok-marke` |
| 7: Meldungen des Smokes beginnen mit `release: smoke: ` | anderes Präfix (`m-p7-smoke-name-weicht-ab`) | `p7-smoke-aufruf` |
| 7: Smoke aus jeder Registry | die Schleife läuft nur für `ghcr.io` (`m-p7-smoke-nur-eine-registry`) | `p7-smoke-aufruf`, `p7-smoke-docker` |
| 7: `POSTGRES_IMAGE` (aus `integration.mk`) geht an den Smoke | die Variable entfällt (`m-p7-smoke-postgres-image-fehlt`) | `p7-smoke-aufruf` |
| 7: Referenz `<Registry>/pt9912/pgwire-recorder@<D>` | Tag statt Digest (`m-p7-smoke-referenz-ohne-digest`) | `p7-smoke-aufruf` |
| 7: `VERSION` geht an den Smoke | die Variable entfällt (`m-p7-smoke-version-fehlt`) | `p7-smoke-aufruf` |
| 7: der Smoke bekommt ein vorhandenes Verzeichnis | ein nicht angelegtes Verzeichnis (`m-p7-smoke-werk-fehlt`) | `p7-smoke-aufruf` |
| 7: der lokale Digest `D` wird verglichen | `D` folgt dem Digest aus `ghcr.io` (`m-p7-vergleich-gegen-d-entfernt`) | `p7-digest-ghcr`, `p7-digest-lokal` |
| 8: `--all` beim Setzen von `latest` | das Argument entfällt (`m-p8-all-entfernt`) | `p8-stabil-hoechste` |
| 8: nur stabile SemVer-Tags zählen für `latest` | jeder Tag `v<Ziffer>…` zählt (`v1.2`) (`m-p8-fremder-tag-zaehlt`) | `p8-fremder-tag`, `p8-stabil-hoechste`, `p8-vorabversion-hoeher-zaehlt-nicht` |
| 8: nur die höchste stabile Version setzt `latest` | die Prüfung entfällt (`m-p8-hoechste-entfaellt`) | `p8-aelterer-zweig`, `p8-semver-numerisch` |
| 8: Vergleich nach SemVer (1.10.0 über 1.9.0) | Sortierung ohne `-V` (`m-p8-lexikographisch`) | `p8-semver-numerisch` |
| 8: `--preserve-digests` beim Setzen von `latest` | das Argument entfällt (`m-p8-preserve-digests-entfernt`) | `p8-stabil-hoechste` |
| 8: `latest` zeigt auf `D` (Quelle `@<D>`) | Quelle ist der Tag `<Version>` (`m-p8-quelle-tag`) | 8 Fälle, darunter `p10-stdin`, `p11-wiederanlauf-latest`, `p8-erster-stabiler` |
| 8: `ghcr.io` zuerst | Schleife `docker.io ghcr.io` (`m-p8-reihenfolge-vertauscht`) | `p8-stabil-hoechste`, `p8-zweite-registry-faellt` |
| 8: `latest` nur nach grünem Smoke (`smoke.ok`) | die Prüfung der Marke entfällt (`m-p8-smoke-marke-nicht-verlangt`) | `p11-reihenfolge-gruende-b`, `p8-nach-rotem-vergleich`, `p8-ohne-smoke` |
| 8: eine Vorabversion lässt `latest` unverändert | die Bedingung entfällt (`m-p8-vorabversion-setzt-latest`) | `p8-vorabversion` |
| 8: Vorabversionen und fremde Tags zählen nicht als höher | alle Tags `v*` zählen (`m-p8-vorabversionen-zaehlen`) | `p8-fremder-tag`, `p8-stabil-hoechste`, `p8-vorabversion-hoeher-zaehlt-nicht` |
| 8: die Meldung nennt, welche Registry auf welchen Digest zeigt | der Zustand entfällt (`m-p8-zustand-fehlt`) | `p8-zweite-registry-faellt` |
| 8: scheitert die zweite Registry, ist der Lauf rot | der Exit von `skopeo copy` wird überlesen (`m-p8-zweite-registry-ignoriert`) | `p11-wiederanlauf-latest`, `p8-zweite-registry-faellt` |
| 9: `github-release` läuft erst nach grünem Smoke (`smoke.ok`) | die Prüfung der Marke entfällt (`m-p9-smoke-marke-nicht-verlangt`) | `p11-alte-marke-b`, `p9-nicht-vor-registries` |
| 10: kein `env` | `env` im Fehlerfall (`m-p10-env-bei-fehler-registry`) | `p10-nicht-im-log-image` |
| 10: kein `set -x` | `set -x` im Skript (Gruppe b) (`m-p10-set-x-image`) | `p10-nicht-im-log-image`, `p3-schluesselwort-geheimnis`, `p3-secret-fehlen-beide`, `p3-secret-fehlt-user`, `p3-secret-leer-user` |
| 10: Anmeldung nur per `--password-stdin` | `--password <Token>` als Argument (`m-p10-token-als-argument`) | `p10-stdin` |
| 10: kein Token in der Ausgabe (Fehler im Vergleich der Tags) | das Skript gibt `$RELEASE_DOCKERHUB_TOKEN` aus (`m-p10-token-bei-fehler-registry`) | `p10-nicht-im-log-image` |
| 11: kein Schritt läuft mit den Dateien des Images eines früheren Laufs | `eingaben` löscht sie nicht (`m-p11-eingaben-loescht-image-nicht`) | `p11-alte-marke-b` |
| 11: kein Schritt läuft mit `smoke.ok` eines früheren Laufs | `eingaben` löscht die Marke nicht (`m-p11-eingaben-loescht-smoke-nicht`) | `p11-alte-marke-b` |
| 11: kein Schritt läuft mit `vergleichen.ok` eines früheren Laufs | `eingaben` löscht die Marke nicht (`m-p11-eingaben-loescht-vergleichen-nicht`) | `p11-alte-marke-b` |
| 11: `image` läuft nicht mit den Dateien eines früheren Laufs | `image` löscht sie nicht (`m-p11-image-loescht-nicht`) | `p11-alte-marke-b` |
| 11: `smoke` läuft nicht mit der Marke eines früheren Laufs | `smoke` löscht seine Marke nicht (`m-p11-smoke-loescht-marke-nicht`) | `p7-smoke-ok-marke` |
| 11: `vergleichen` läuft nicht mit der Marke eines früheren Laufs | `vergleichen` löscht seine Marke nicht (`m-p11-vergleichen-loescht-marke-nicht`) | `p11-alte-marke-b` |
| 13: `image` läuft nicht ohne grünes `binaries` | die Prüfung der Marke entfällt (`m-p13-image-ohne-binaries`) | `p11-reihenfolge-gruende-b`, `p5-ohne-binaries` |
| 13: `smoke` läuft nicht ohne grünes `vergleichen` | die Prüfung entfällt (`m-p13-smoke-ohne-vergleichen`) | `p11-alte-marke-b`, `p11-reihenfolge-gruende-b`, `p7-smoke-ohne-vergleichen` |
| 13: `vergleichen` läuft nicht ohne `image` | die Prüfung entfällt (`m-p13-vergleichen-ohne-image`) | `p11-reihenfolge-gruende-b` |
| 13: `veroeffentlichen` läuft nicht ohne `image` | die Prüfung entfällt (`m-p13-veroeffentlichen-ohne-image`) | `p11-alte-marke-b`, `p11-reihenfolge-gruende-b` |

**Fälle von `test-container-gegenprobe` mit neuem Ort:** die Mutationen zu Punkt 7 (`p7-intern`, `p7a-*`, `p7b-*`, `p7c-*`, `p7-exit-*`, `p7-terminal`, `p7-user*`, `p7-volume*`, `p7-mount`, `p7-cp-*`, `p7-rm-pg`), `p2-dev` und `p9-volume-replay` zielen auf `tools/harness/container-smoke.sh` statt `run-container.sh`, mit unveränderten roten Mengen; `p7-tag` zielt auf den Aufruf in `run-container.sh`. Neu: `p9-smoke-weiter` (`kill -TERM "$laeuft"` im Smoke entfällt) färbt `p9-sigINT-psql-signal` und `p9-sigTERM-psql-signal`.

**Keine äquivalenten Mutanten:** jede Mutation von B färbt mindestens einen Fall (die Ausgabe von `GP_ENTDECKE=1` nennt keine Mutation ohne roten Fall).

**Die Fälle prüfen auch die Meldung, nicht nur den Ausgang:** `p5-zwei-manifeste` und `p7-eine-plattform` (Zahl und Plattformen) und die Fälle `p7-digest-*` prüfen den Text, der den abweichenden Wert nennt; deshalb färbt `m-p5-zahl-entfernt` auch `p7-eine-plattform`, und `m-p7-vergleich-gegen-d-entfernt` auch `p7-digest-ghcr` (die Meldung nennt dann `docker.io`).

**Grenzmutationen** (Zusagen der Form „nur X“, Lehre aus `slice-erster-release-veroeffentlichung`): `m-p6-fehler-wie-fehlend`, `m-p6-muster-geweitet` (um `denied`), `m-p6-wortlaut-*-nicht-erkannt` (je Wortlaut) und `m-p6-schreibweise-beachtet` (nur drei Meldungen heißen „Tag fehlt“, ohne Beachtung der Schreibweise), `m-p8-vorabversionen-zaehlen` und `m-p8-fremder-tag-zaehlt` (nur stabile SemVer-Tags zählen), `m-p5-plattform-amd64-nicht-verlangt` und `-arm64-` (genau diese zwei Plattformen), `m-p6-vorhandener-tag-ueberschrieben` (nie überschreiben, gegen die Entfernung der Meldung `m-p6-anderer-digest-angenommen`).

**Messungen gegen den echten Docker** (lokal im Container, ohne Registry): Ein Bau mit `docker buildx build --provenance=false --sbom=false -o type=oci,dest=x.tar` für `FROM scratch` (eine Plattform, Treiber `docker`) liefert ein Tar mit `index.json`, `oci-layout` und `blobs/sha256/…`, ohne Präfix `./`; `index.json` ist kompakt (`"digest":"sha256:…"` und `"platform":{"architecture":"amd64","os":"linux"}` in dieser Folge). `release.sh` liest daraus mit denselben Ausdrücken (`manifeste`, `plattformen`, die Hexziffern des Digests) 1, `linux/amd64` und den Namen des Blobs. **Nicht gemessen** (Grenzen in `SPEC-067` Punkte 5 und 6): das Layout eines Baus für zwei Plattformen mit dem Treiber `docker-container` (jede andere Form als ein Index-Blob mit zwei Manifesten ist rot), die Wortlaute `manifest unknown`, `name unknown`, `not found` von `ghcr.io` und `docker.io`, die Wirkung von `SOURCE_DATE_EPOCH` ohne `ARG`, `skopeo` und `gh` im Dev-Container (nicht installiert).

**Trockenlauf auf GitHub:** nicht gefahren (der Auftrag verbietet Push und Netz). Der Pin des Builders ist gesetzt (`moby/buildkit@sha256:cec9f139…3dea`, `# v0.33.1`, Messung des Architect vom 2026-10-19); der Lauf mit Lauf-Kennung und Dauer kommt in die Closure.

**Handbuch und README:** keine Zeile neu oder geändert. Gelesen gegen `make test-container-gegenprobe` und das grüne `make test-container` in `make gates`: Plattformen des Images, Prüfung des Binaries im Image und der Absatz in `README.md` stehen im Ist-Zustand; der Bezug aus einer Registry gehört `slice-erster-release-freigabe` (§1).

**`releasing.md` (Maintainer-Doku), Sätze mit Verhaltensaussage und ihr Beleg:**

| Satz (Abschnitt) | Beleg |
|---|---|
| §6 Schritt 3: ein OCI-Layout mit genau einem Index und zwei Manifesten; die Summe des Binaries je Plattform | `p5-ein-index`, `p5-zwei-manifeste`, `p5-plattform-fehlt-*`, `p5-summe-amd64`, `p5-summe-arm64` |
| §6 Schritt 3: Treiber `docker-container`, das Image des Builders mit Digest | `p5-aufruf`, `p5-builder-pin-form`, `p5-builder-platzhalter` |
| §6 Schritt 4: dasselbe Layout ohne Neubau, `ghcr.io` zuerst, Tag `<SemVer>` ohne `v` | `p6-reihenfolge`, `p6-all-und-digests`, `p6-tag-satz` |
| §6 Schritt 4: gleicher Digest übersprungen, anderer Digest rot, nie überschrieben | `p6-tag-vorhanden-gleich`, `p6-tag-vorhanden-anders`, `p6-tag-fehlt-*`, `p6-inspect-anderer-fehler` |
| §6 Schritt 5: rohe Bytes ohne Anmeldung, drei Digests gleich, genau zwei Plattformen, Paket nicht öffentlich ist rot | `p7-gleich`, `p7-rohbytes`, `p7-digest-ghcr`, `-docker`, `-lokal`, `p7-eine-plattform`, `p7-anonym-ghcr`, `-docker` |
| §6 Schritt 6: Smoke aus jeder Registry mit `<Registry>/pt9912/pgwire-recorder@<D>` | `p7-smoke-aufruf`, `p7-smoke-ghcr`, `p7-smoke-docker` |
| §6 Schritt 7: `latest` auf `D`, `ghcr.io` zuerst, nur stabil und nur die höchste | `p8-stabil-hoechste`, `p8-vorabversion`, `p8-aelterer-zweig`, `p8-erster-stabiler`, `p8-vorabversion-hoeher-zaehlt-nicht`, `p8-semver-numerisch`, `p8-fremder-tag`, `p8-ohne-smoke` |
| §6 Schritt 8: der GitHub-Release läuft nur nach grünem `Smoke` | `p9-nicht-vor-registries` |
| §4 Punkt 4: fehlt oder ist leer ein Geheimnis, endet `Eingaben` mit Exit 2 | `p3-secret-fehlt-user`, `-token`, `p3-secret-leer-user`, `-token`, `p3-secret-fehlen-beide` |
| §6 Trockenlauf: `Image` ohne Anmeldung, ohne Veröffentlichung | `p2-trockenlauf-image` |
| §11 Tabelle: Zeilen `veröffentlichen`, `vergleichen, smoke`, `latest`, Wiederanlauf | `p6-zweite-registry-faellt`, `p6-tag-vorhanden-anders`, `p11-wiederanlauf-ghcr-vorhanden`, `p11-wiederanlauf-latest`, `p11-reihenfolge-gruende-b`, `p11-alte-marke-b` |
| §8: Smoke auf dem Docker-Server des Runners, `linux/arm64` nicht gestartet | kein Fall (Plattform des Runners, `SPEC-067` Kopf und Punkt 15); die Pipeline startet nur das Image, das Docker für den Server wählt |
| §1: der Scan fehlt, eine Formel gibt es nicht | Tatsache des Standes |

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`, Greenfield); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `docs/plan/planning/observations/BEO-REPO/` am Stand `52dc73fc` gesichtet (Zähler = Dateien unter `evidence/`). Treffer:

- `BEO-REPO/slice-waechst-durch-uebernahmen` (7×) — dieser Slice nimmt Gegenstand aus `slice-erster-release-veroeffentlichung` auf; die Aufnahme ist als Liefer-Punkte in §2 gezählt (drei), der Schnitt lief vor dem Start (Größe 1970 bis 2570 gegen 1800). Kein neuer Beleg.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (42×, verkörpert in `AGENTS.md` §3.12) — die Randformen stehen in §6 und sind in `SPEC-067` entschieden; der Implementer hält an und gibt zurück, was dort nicht steht.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (58×, `AGENTS.md` §3.10) und `BEO-REPO/gruener-mutant-faelschlich-aequivalent` (9×) — der Vertrag trägt Fälle und Mutationen je Zusage und je Element; ein als äquivalent geführter Mutant braucht Invariante, Ort und Sonde.
- `BEO-REPO/verhalten-nur-unter-linux-geprueft` (3×, verkörpert in `AGENTS.md` §3.12) — `linux/arm64` ist nach Struktur und Summe belegt, nicht ausgeführt, `darwin` und `windows` sind Sache von A (`SPEC-067` Kopf, Punkt 15); die Plattform steht in §6.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (13×, `AGENTS.md` §3.13) — die Sendungen an `slice-erster-release-freigabe` (Einhängepunkt, Handbuch zum Bezug) stehen im selben Commit dort in §1; `grep` nach dieser Kennung findet sie.
- `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (17×, `AGENTS.md` §3.14) — die Schätzung rechnet mit dem Faktor 2 der Gegenprobe (§4).

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Nachgezählt am 2026-10-19** (Planner, `AGENTS.md` §3.13, beim Anlegen im Commit, der den Schnitt vollzieht und die Sendungen einträgt): Liefer-Punkte drei (§2: Image bauen und prüfen mit Handbuch- und `README.md`-Teil | Veröffentlichen, Vergleichen, Smoke | `latest`, Fehlertabelle, Wiederanlauf); die Endfassung `1b` von `release.yml` ist Akzeptanzkriterium, kein vierter Punkt, ebenso der Auszug des Smokes (Teil von Punkt 2). Schichten zwei (Release-Automatisierung: Workflow, `tools/harness/`, `harness/mk/`, gegebenenfalls `Dockerfile`; Dokumentation: `docs/maintainer/`, `docs/user/`, `README.md`); die Gegenprobe zählt mit der Automatisierung. Der Eintrag der Sendungen an `slice-erster-release-freigabe` liegt dort; er nimmt aus diesem Slice nichts, B gibt sie ab.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Nachgezählt am 2026-10-19** (Implementer, nach den Entscheidungen des Architect `5cb3c3ce`): Liefer-Punkte drei, Schichten zwei, unverändert; die neuen Marken (`smoke.ok`), Fälle und der Pin liegen in der Release-Automatisierung und ihrer Gegenprobe.
