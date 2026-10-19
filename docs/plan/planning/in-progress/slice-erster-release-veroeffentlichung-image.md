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
- [ ] **Liefer-Punkt 2, Veröffentlichen, Vergleichen, Smoke** ([`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), `SPEC-067` Punkte 6, 7, 9 (Vorbedingung), 10 (Registries), [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 7): Der Push kopiert das Layout ohne Neubau in beide Registries (`ghcr.io` zuerst), ein vorhandener Tag wird nie überschrieben; der Vergleich der drei Digests `D`, `ghcr.io`, `docker.io` aus den rohen Bytes des Index, der anonyme Zugriff und die Plattformen sind fail-closed; der Smoke (`tools/harness/container-smoke.sh <Referenz>`, ausgelagert aus dem Smoke von `make test-container`, dessen Gate und Gegenprobe unverändert grün bleiben) besteht aus **jeder** Registry; der Schritt `github-release` von A läuft erst nach grünem `vergleichen` (`p9-nicht-vor-registries`). Anmeldung nur per `--password-stdin`, kein Wert eines Geheimnisses im Log (Test: Fälle `p6-*`, `p7-*`, `p10-*` ohne `GITHUB_TOKEN`-Anteil von A und `p9-nicht-vor-registries`, je Zusage und Element eine Mutation; `make test-container-gegenprobe` unverändert grün). Dass die Token schreiben dürfen, das Paket öffentlich ist und der Smoke aus den Registries besteht, belegt erst der Lauf für den ersten Tag (`SPEC-067` Punkt 12), nicht dieser Slice.
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
| `tools/harness/run-container.sh`, `harness/mk/container.mk`, `tools/harness/container-gegenprobe.sh` | update | rufen das ausgelagerte Skript statt eigenem Text; Verhalten und Fälle von `SPEC-055` unverändert (`make test-container`, `make test-container-gegenprobe` grün) |
| `tools/harness/release-gegenprobe.sh` | update | Fälle `p3-secret-*`, `p5-*`, `p6-*`, `p7-*`, `p8-*`, `p9-nicht-vor-registries`, `p10-*` (Registries), `p11-*` (B), Zeilenliste von `release.yml` mit Herkunft B; Attrappen von `docker` (`buildx`), `skopeo`, `gh`; Fälle von A bleiben unverändert (nur die Fixture um `vergleichen.ok` ergänzt) |
| `Dockerfile` | update, falls nötig | `ARG SOURCE_DATE_EPOCH` für die Stufe `runtime`, wenn der Bau der Zeitstempel sie braucht; der Architect entscheidet es mit `SPEC-067` Punkt 5 und zählt nach (§8) |
| `docs/maintainer/releasing.md` §2, §7, §8, §11 | update | Image-Anteile des ausgeführten Verfahrens, übrige Zeilen der Fehlertabelle |
| `docs/user/benutzerhandbuch.md`, `README.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11); nur, was Gegenprobe und Trockenlauf belegen |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-erster-release-veroeffentlichung` liegt in `done/` (Fassung `1a` von `release.yml`, `release.sh` mit den drei Schritten, `VERSION`, Gegenprobe mit Zeilenliste und der Trockenlauf belegt). Schritt 5 der Reihenfolge in §5 von [welle-erster-release](../welle-erster-release.md). Die Randformen aus §6 sind entschieden ([`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md), `AGENTS.md` §3.12); der Implementer hält an und gibt eine nicht genannte Randform an den Architect zurück.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Diff ist nicht in einer Review-Sitzung prüfbar. Schnitt dann: *Bau des Images* (Liefer-Punkt 1) gegen *Veröffentlichen, Vergleichen, `latest`* (Punkte 2 und 3); die Endfassung `1b` müsste dann in zwei Stufen entstehen, was der Architect vorher in `SPEC-067` teilt. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'` liegt über **1700**, Tests und Gegenprobe eingeschlossen; Basis: beim Anspruch (HEAD vor dem Anspruch, hier einzutragen). Schätzung des Architect in dieser Metrik (2026-10-19, nach `SPEC-067`, Fassung B): **1150 bis 1500** (Gegenprobe mit Attrappen für `docker buildx`, `skopeo`, `gh` das Doppelte des Produktcodes); Mitte 1325, 78 % der Marke. Meldeschwelle **1300**: Wer sie überschreitet, meldet es vor dem nächsten Abschnitt, und die Entscheidung („weiter“ oder „schneiden“) steht als datierte Zeile in §7. Faktor Tests zu Produktcode: 2, Beleg `slice-harness-ci-workflows` (Messung 899 gegen eine Schätzung von 545 bis 695, Drift-Log der Roadmap, 2026-10-18); das Drift-Log der Roadmap nennt für `slice-harness-ci` den Faktor 1,3 bis 1,65 der Unterschätzung des Architect, deshalb liegt die Obergrenze der Schätzung nah an der Marke.
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
