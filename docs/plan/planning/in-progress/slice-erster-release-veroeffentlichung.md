# Slice slice-erster-release-veroeffentlichung: Veröffentlichung von Tag und Binaries

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-erster-release.

**Bezug:** [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)

**Berührte Spec-Stellen:** `LH-FA-19.a` · `SPEC-035` · `SPEC-055` (gelesen: Voraussetzung `slice-v1-abschluss-container`, `slice-v1-abschluss-container-doku`) · `SPEC-063` und `SPEC-064` (gelesen: Gate der Workflows und Konventionen der CI; Sendung aus `slice-harness-ci-workflows` und `slice-harness-ci`, §1 unten) · `SPEC-067` (Fassung A: Punkte 1a, 2, 3 (1) bis (6), 4, 9, 10, 11, 12, 13 bis 15, jeweils der Anteil von A; Zuordnung in der Tabelle des Vertrags)

**Verantwortlich:** pt9912 (Zweig impl-rel-a)
**Autor:** pt9912. **Datum:** 2026-10-03. **Geschnitten:** 2026-10-19 (siehe §4, *Schnitt vor dem Start, 2026-10-19*).

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Tag `v<SemVer>` wird gegen die Versionsdatei `VERSION` und die Zeile der Änderungshistorie des Handbuchs geprüft, bevor etwas gebaut wird; die Pipeline `release.yml` hängt die Binaries für die Zielplattformen mit `SHA256SUMS` aus `make binaries` an den GitHub-Release (die Reproduzierbarkeit der Binaries hält `make test-container`, kein zweiter Bau) und lässt sich nach einem Fehler wieder anstoßen. Dieselbe Datei läuft als Trockenlauf (`workflow_dispatch`) auf einem Zweig und tut dann alles außer dem, was nach außen wirkt; er belegt Eingaben, Binaries und Dauer auf GitHub. Ein Probe-Tag ist nicht vorgesehen (Entscheidung des Nutzers vom 2026-10-19, [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md)). Das Image ist nicht Teil dieses Slice (`slice-erster-release-veroeffentlichung-image`, unten).

**Schnitt vom 2026-10-19** (Empfehlung des Planners auf Anweisung der Koordination, vom Nutzer nicht bestätigt, nach dem Vertrag [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md), Zuordnung dort; die Marke und die Rückführung in §4 sahen diesen Schnitt vor): Dieser Slice ist Fassung A des Vertrags (Tag und Binaries); das Image (Bau, Veröffentlichen, Vergleichen, Smoke, `latest`) ist Fassung B in `slice-erster-release-veroeffentlichung-image`, das hinter diesem Slice liegt. Die Fassung `1a` von `release.yml` setzt nichts von B voraus; B hängt seine Zeilen an.

**Übernommen aus `welle-erster-release`** (Entscheidungen des Nutzers vom 2026-10-09 nach dem
Vergleich von `docs/maintainer/releasing.md` mit den Release-Dokumenten von a-check, d-check,
pg-change-feed und u-boot; `docs/maintainer/releasing.md` nennt diesen Slice dort, wo er
liefert); die Punkte 1 und 2 gingen mit dem Schnitt vom 2026-10-19 an `slice-erster-release-veroeffentlichung-image` (dort §1):

1. *(an `slice-erster-release-veroeffentlichung-image`)* `:latest` zeigt nur nach einem stabilen Tag auf das neue Image.
2. *(an `slice-erster-release-veroeffentlichung-image`)* Der Spiegel auf Docker Hub ist fail-closed.
6. Zusätzlich zum Tag eine Versionsdatei, die die Pipeline gegen den Tag prüft: `VERSION` im Wurzelverzeichnis, eine Zeile `<SemVer>` ohne `v` ([`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md), Schnittstelle).
7. Die Prüfung des Tags auf die Form `v<SemVer>` vor jedem Schritt mit Außenwirkung.
8. Wiederanlauf und Fehlertabelle: `docs/maintainer/releasing.md` §11 nennt je Schritt der
   Pipeline, was ein Fehler dort hinterlässt und wie der Lauf wieder angestoßen wird
   (Entscheidung des Nutzers vom 2026-10-09 beim Schnitt unten); dieser Slice trägt die Zeilen
   `eingaben`, `binaries` und `github-release`, B die übrigen.

Die Randformen dieser Punkte sind entschieden: [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) (Architect, 2026-10-19, vom Nutzer mit [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md) angenommen; §6).

**Voraussetzung `slice-v1-abschluss-container` und `slice-v1-abschluss-container-doku`** (Hinweis des Architect, 2026-10-17; `-container-doku` seit dem Schnitt vom 2026-10-17 (Handbuch, README, `docs/maintainer/releasing.md`): die Pipeline und das Release-Dokument bauen auf deren Stand auf; keine Sendung, kein Eintrag in §8): Die Binaries entstehen mit `make binaries` (`dist/` mit `SHA256SUMS`, [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md)); die Pipeline hängt `dist/` an und baut kein zweites Mal, die Reproduzierbarkeit (zwei Builds, dieselben Summen) hält dort `make test-container`.

**Sendung aus `slice-harness-ci-workflows` und `slice-harness-ci` — Pins und Härtung** (Planner, 2026-10-18, `AGENTS.md` §3.13; Auftrag des Nutzers vom 2026-10-18: die CI gehört vor diesen Slice): Die Datei der Pipeline, `.github/workflows/release.yml` (Name: Architect), fällt vom ersten Commit an unter `make workflows-check` und `make docs-check`. Sie trägt damit dieselben Konventionen wie `ci.yml`: jede fremde `uses:`-Referenz per vollem SHA mit Tag-Kommentar, `permissions: {}` oben und im Job nur die Rechte, die er braucht (in Fassung `1a`: `contents: write` für den GitHub-Release; `packages: write` für `ghcr.io` kommt mit B), `timeout-minutes`, `concurrency` mit `cancel-in-progress: false` (ein Tag wird nie durch einen späteren abgebrochen). `ci.yml` lässt Tags aus (`tags-ignore`), `release.yml` reagiert nur auf `v*`; ein Tag läuft also nur hier. Zugangsdaten nur als Geheimnisse von GitHub oder als `GITHUB_TOKEN`, nie im Repository. Welche Eigenschaften der Datei ein Test hält, legt `SPEC-067` Punkt 1 fest (Text der Datei, Zeilenliste in der Gegenprobe), die Prüfung der Form hält das Gate. Der **Träger** (GitHub Actions) ist auch für die Veröffentlichung entschieden ([ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md), Randform *Träger der Automatisierung*, §6). Voraussetzung ist beim Start: `.github/workflows/` liegt (aus `slice-harness-ci`).

**Abgegeben** an `slice-erster-release-freigabe` (dort §1, *Übernimmt*, mit der Kennung dieses
Slice), nach dem Schnitt vom 2026-10-09 vor dem Start (Entscheidung des Nutzers, §4): die
Entscheidungen 3 (Freigabe-Checkliste), 4 (Freshness-Audit) und 5 (Image-Scan vor dem Push)
und der erste echte Tag. Dieser Slice belegt die Pipeline mit dem Trockenlauf (statt eines Probe-Tags, Entscheidung des Nutzers vom 2026-10-19).

**Abgegeben** an `slice-erster-release-veroeffentlichung-image` (dort §1, *Übernimmt*, mit der Kennung dieses
Slice), Schnitt vom 2026-10-19: das Image (Bau beider Plattformen in ein OCI-Layout, Veröffentlichen, Vergleichen der Index-Digests, Smoke aus den Registries, `latest`), die Geheimnisse der Registries und ihre Anmeldung, die Image-Anteile von `docs/maintainer/releasing.md`, Handbuch und `README.md` und der Auszug des Smokes aus `SPEC-055` Punkt 7 in `tools/harness/container-smoke.sh`.

**Abgegeben** an `slice-erster-release-freigabe` (dort §1 und DoD 3, mit der Kennung dieses Slice), Schnitt vom 2026-10-19: Der Abschnitt des Handbuchs §2 *Installation* zum Bezug der Binaries vom GitHub-Release (Plattformen, Prüfung mit `SHA256SUMS`) steht erst, wenn der erste echte Tag ihn belegt (`AGENTS.md` §3.11, Ist-Zustand); vorher gibt es keinen Release, aus dem man sie beziehen kann.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Homebrew-Formel und ihr Nachweis — `slice-erster-release-homebrew-nachweis` (`welle-v1-nachtrag`).
- Das Image: Bau, Veröffentlichen in `ghcr.io` und `docker.io`, Vergleich der Digests, Smoke aus den Registries, `latest`, die Geheimnisse `DOCKERHUB_USERNAME` und `DOCKERHUB_TOKEN` und das Recht `packages: write` — `slice-erster-release-veroeffentlichung-image` (oben, *Abgegeben*); Fassung `1a` von `release.yml` hat keinen Schritt mit Registry und liest kein Geheimnis außer `GITHUB_TOKEN`, sodass dieser Slice ohne B lieferbar ist.
- Image-Scan, Freshness-Audit, Freigabe-Checkliste und der erste echte Tag —
  `slice-erster-release-freigabe` (oben, *Abgegeben*); er setzt diese Pipeline und B voraus, weil der
  Scan in ihr zwischen Bau und Push des Images läuft.
- Ein Probe-Tag — nicht vorgesehen (Entscheidung des Nutzers vom 2026-10-19): Er legte ein Image und einen GitHub-Release einer Vorabversion ab, die kein Release sind; der Trockenlauf belegt, was sich ohne Außenwirkung belegen lässt, und die erste Außenwirkung ist der Lauf für den ersten echten Tag (`SPEC-067` Punkt 12).
- Die Wahl der ersten Versionsnummer — wird beim Release festgelegt; die Versionsdatei
  (Punkt 6) trägt sie, wählt sie aber nicht (Anfangswert `0.1.0-dev`).
- Konten, Namensräume und Zugangsdaten der Registries, die Sichtbarkeit der Pakete und das
  Hinterlegen der Zugangsdaten im Repository — Handlungen des Betreibers (`docs/maintainer/releasing.md`
  §4 Punkt 4); dieser Slice liest nur `GITHUB_TOKEN`. Die Zugangsdaten stehen nie im Repository.
- Ein zweiter Bau der Binaries oder ein Vergleich zweier Builds in der Pipeline — das Gate
  `make test-container` (Geber `slice-v1-abschluss-container`) hält die Zusage; die Pipeline
  ruft `make binaries` und baut nichts ein zweites Mal.
- Der Abschnitt des Handbuchs zum Bezug vom GitHub-Release — `slice-erster-release-freigabe` (oben, *Abgegeben*): Er beschreibt einen Release, den es erst nach dem ersten Tag gibt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1, Eingaben** ([`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), `SPEC-067` Punkt 3 (1) bis (6), Schnittstelle): Vor jedem Schritt mit Außenwirkung prüft `release.sh eingaben`, in der Reihenfolge des Vertrags, dass Ereignis und Ref stimmen, der Tag die Form `v<SemVer>` hat, die Datei `VERSION` (neu, Wurzel, `0.1.0-dev`) ihre Form hat und dem Tag entspricht, die Änderungshistorie des Benutzerhandbuchs die Zeile der Version trägt und der Commit des Tags im Standard-Zweig liegt; sonst Exit 2, bevor etwas gebaut oder veröffentlicht wird (Test: die Fälle `p3-*` der Gegenprobe mit je einer Mutation, `AGENTS.md` §3.10). `make binaries` übernimmt die Version aus `VERSION=`, ohne Angabe bleibt es bei `dev`.
- [ ] **Liefer-Punkt 2, Binaries am GitHub-Release** ([`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), `SPEC-067` Punkte 4 und 9): Die Binaries für Linux, macOS und Windows (`amd64`, `arm64`) hängen mit einer Datei `SHA256SUMS` am GitHub-Release (`LH-FA-19.a`); die Dateien am Release sind die aus `make binaries` mit `VERSION` aus dem Tag, genau sieben, und eine abweichende Summe am Release endet rot (Test mit Attrappen von `make` und `gh`: Fälle `p4-*` und `p9-*`; die Gleichheit zweier Builds belegt `make test-container`, [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 4, nicht dieser Slice). Die Release-Notes entstehen aus der neuen Zeile der Änderungshistorie des Handbuchs; ein `CHANGELOG.md` wird nicht geführt. Dass der GitHub-Release auf GitHub angelegt wird, belegt erst der Lauf für den ersten Tag (`SPEC-067` Punkt 12); der Trockenlauf legt keinen an. Handbuch §1 *Voraussetzungen* und §2 *Installation* und `README.md` sind auf das geprüft, was der Trockenlauf und `make binaries` belegen (Plattformen der Binaries, `SHA256SUMS`), und im Ist-Zustand des gebauten Binaries berichtigt, wo sie es nicht sind: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11); der Abschnitt zum Bezug vom GitHub-Release gehört `slice-erster-release-freigabe`.
- [ ] **Liefer-Punkt 3, Trockenlauf, Wiederanlauf, Fehlertabelle** (`SPEC-067` Punkte 1a, 2, 11, 12): Die Datei `.github/workflows/release.yml` hat den Text der Fassung `1a`, reagiert nur auf `v*` und `workflow_dispatch`, und der Trockenlauf fährt `Eingaben` und `Binaries` ohne Außenwirkung (kein `gh`, kein GitHub-Release; Test: `p1-*`, `p2-*` mit Mutationen). Der Trockenlauf ist einmal auf dem Stand der Lieferung gefahren und in §7 belegt (Lauf-Kennung, Dauer, `SPEC-067` Punkt 12; Handlung des Betreibers oder des Nutzers, falls der Lauf Rechte am Repository braucht). `docs/maintainer/releasing.md` §11 trägt die Fehlertabelle mit den Zeilen `eingaben`, `binaries` und `github-release` und den Wiederanlauf (Punkt 8; die übrigen Zeilen liefert B); ein roter Lauf ist nach ihr wieder angestoßen (Beleg in §7). Die Datei besteht `make workflows-check` und `make docs-check` und trägt die Pins, Rechte und Zeitgrenzen der Sendung aus `slice-harness-ci-workflows` und `slice-harness-ci` (§1); sie ist Akzeptanzkriterium dieses Punkts und von Punkt 2, kein vierter Liefer-Punkt.
- [ ] `make gates` grün.
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
| `.github/workflows/release.yml` (Fassung `1a`, `SPEC-067` Punkt 1) | neu | nur `v*` und `workflow_dispatch`; Schritte `Eingaben`, `Binaries`, `GitHub-Release` (nur bei `push`); Pins und Rechte wie `ci.yml` |
| `tools/harness/release.sh` (Schritte `eingaben`, `binaries`, `github-release`) | neu | die Logik steht im Skript, nicht im YAML (Schnittstelle des Vertrags); der Zustand zwischen den Schritten liegt in `$RUNNER_TEMP/release` (Marken `eingaben.ok` und `binaries.ok`, Verzeichnis `dist/`), jeder Schritt prüft die Marke seines Vorgängers; B ergänzt die Schritte `image`, `veroeffentlichen`, `vergleichen`, `smoke`, `latest` und hängt seine Marken an |
| `VERSION` (Wurzel) | neu | Version der Software, eine Zeile `0.1.0-dev`, gegen den Tag geprüft; `make binaries VERSION=` trägt sie ins Binary |
| `tools/harness/release-gegenprobe.sh`, `harness/mk/release.mk`, Eintrag in `GEGENPROBEN_CONTAINER` | neu | Gegenprobe mit Attrappen für `make`, `gh`, `docker` und `skopeo` (`git`, `sha256sum`, `awk`, `cmp` sind die echten Werkzeuge), ohne Netz: Fälle `p1-*`, `p2-*`, `p3-*` ohne (7), `p4-*`, `p9-*` ohne `p9-nicht-vor-registries`, `p10-*` ohne die Registries, `p11-*` für A, `p13-*`; je Zusage und je Element eine Mutation (`AGENTS.md` §3.10); die Zeilenliste von `release.yml` trägt die Herkunft A je Zeile |
| `harness/README.md` §Sensors | update | Zeile für `make release-gegenprobe` (Gegenprobe), Bindung `SPEC-067` und [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md) |
| `docs/maintainer/releasing.md` §1, §2, §3, §5, §6, §9, §11 | update | vom beschriebenen auf das ausgeführte Verfahren bringen, soweit A es liefert: Stand, Versionsquelle (`VERSION`), was der Tag auslöst, Release-Notes, Fehlertabelle (Zeilen `eingaben`, `binaries`, `github-release`) und Wiederanlauf; auf dem Stand aufbauend, den `slice-v1-abschluss-container-doku` liefert. §8 (Kontrollen nach dem Push der Registries), §2 und §7 (Formel keine Lieferung des Release) ergänzt B beziehungsweise gehören nicht hierher |
| `docs/user/benutzerhandbuch.md`, `README.md` | geprüft, unverändert | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung); die Sätze über Plattformen und `SHA256SUMS` der gebauten Binaries stehen schon im Ist-Zustand, §7 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-v1-abschluss` ist `done`, und
`slice-harness-abdeckung-gate` und `slice-harness-coverage` liegen in `done/` (Start-Trigger von
welle-erster-release), dazu `slice-harness-devcontainer-u-boot-0-10`
(wellenlos, nicht `slice-harness-devcontainer-einrichten`; Entscheidung des Nutzers vom 2026-10-17, Drift-Log der Roadmap). `slice-harness-ci-workflows` und `slice-harness-ci` liegen in `done/` (Sendung in §1; Auftrag des Nutzers vom 2026-10-18). Schritt 4 der Reihenfolge in §5 von
[welle-erster-release](../welle-erster-release.md) (A; B folgt als Schritt 5, `slice-erster-release-freigabe` als Schritt 6). Die Randformen aus §6 sind entschieden ([`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md), `AGENTS.md` §3.12); der Implementer hält an und gibt eine nicht genannte Randform an den Architect zurück.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: die Release-Automatisierung von A verlangt mehr als drei Liefer-Punkte, oder der Diff ist nicht in einer Review-Sitzung prüfbar — zurück zur Zerlegung. Ein weiterer Schnitt wäre *Eingaben* (Liefer-Punkt 1) gegen *Binaries und Trockenlauf* (Punkte 2 und 3); die Fassung `1a` von `release.yml` müsste dann beide tragen, und der Architect teilte die Zeilenliste. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'` liegt über **1300**, Tests und Gegenprobe eingeschlossen; Basis: 11a2ba91 (HEAD vor dem Anspruch). Schätzung des Architect in dieser Metrik (2026-10-19, nach `SPEC-067`, Fassung A): **820 bis 1070** (Gegenprobe mit Attrappen das Doppelte des Produktcodes; `release.sh` mit drei Schritten, `release.yml`, `VERSION`, `releasing.md`, Handbuch, `README.md`); Mitte 945, 73 % der Marke. Meldeschwelle **1000**: Wer sie überschreitet, meldet es vor dem nächsten Abschnitt, und die Entscheidung („weiter“ oder „schneiden“) steht als datierte Zeile in §7. Faktor Tests zu Produktcode: 2 (die Gegenprobe trägt zwei Drittel des Diffs), Beleg `slice-harness-ci-workflows` (Messung 899 gegen eine Schätzung von 545 bis 695, Drift-Log der Roadmap, 2026-10-18); eine eigene Messung für diese Gegenprobe gibt es nicht.
- `in-progress` → `open`: Das Repository gibt dem Trockenlauf keine Rechte (der Lauf auf GitHub scheitert an Einstellungen des Betreibers) — extern klären; ebenso, wenn `slice-v1-abschluss-container` oder `slice-v1-abschluss-container-doku` nicht liegen, was sie liefern (Voraussetzung in §1).

**Schnitt vor dem Start, 2026-10-19** (Planner; Anlass: Größe 1970 bis 2570 des ungeteilten Slice gegen die Marke 1800 der Fassung davor, nach dem Vertrag `SPEC-067`): Der Slice ist an der Naht *Tag und Binaries* gegen *Image* geteilt, wie die Rückführung der Fassung davor sie vorsah (Liefer-Punkte 2 und 3 dort gegen Punkt 1). A behält Eingaben, Binaries, GitHub-Release, `VERSION` und den Trockenlauf (Schätzung 820 bis 1070, Marke 1300); B ist `slice-erster-release-veroeffentlichung-image` (Schätzung 1150 bis 1500, Marke 1700, `next/`, Start: dieser Slice in `done/`). Der Slice lag in `next/` und hat nicht begonnen; ein Übergang fällt nicht an.

**Schnitt vor dem Start, 2026-10-09:** Mit den sieben Entscheidungen vom 2026-10-09 trug §2 sechs
Liefer-Punkte. Der Nutzer entschied am 2026-10-09 den Schnitt nach dem Vorschlag in §6: Die
Pipeline bleibt hier, dazu Wiederanlauf und Fehlertabelle; Scan, Audit, Checkliste und der erste
echte Tag gehen an `slice-erster-release-freigabe`. Der Slice lag in `next/` und hat nicht
begonnen; ein Übergang fällt nicht an.

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

- Der Trockenlauf braucht Rechte oder Einstellungen auf GitHub, die aus dem Klon nicht auditierbar sind (Actions im Repository erlaubt, `GITHUB_TOKEN` mit `contents: write`; `SPEC-067` Punkt 15) — **Ausgang:** offen bis Closure.
- Die Zugangsdaten für `ghcr.io` und `docker.io` sind externe Betreiber-Handlungen — sie gehören B (dort §6); A liest nur `GITHUB_TOKEN`. **Ausgang:** entfallen für A (Grund: die Fassung `1a` hat keinen Schritt mit Registry; die Handlung steht bei `slice-erster-release-veroeffentlichung-image`).
- **Größe.** Mit den sieben Entscheidungen vom 2026-10-09 trug §2 sechs Liefer-Punkte, die
  Grenze ist drei. Entscheidung des Nutzers vom 2026-10-09: geschnitten nach dem Vorschlag des
  Planner (§4) — **Ausgang:** eingetreten, Folge-Slice `slice-erster-release-freigabe`. Der Schnitt vom 2026-10-19 (Größe 1970 bis 2570 gegen die Marke 1800) hat die Naht der Rückführung genutzt — **Ausgang:** eingetreten, Folge-Slice `slice-erster-release-veroeffentlichung-image`.
- **Probe-Tag in öffentlichen Registries.** Entfällt mit der Entscheidung des Nutzers vom 2026-10-19: kein Probe-Tag, der Trockenlauf schreibt nichts in eine Registry und legt keinen GitHub-Release an (`SPEC-067` Punkt 2 und 12) — **Ausgang:** entfallen (Grund: die Entscheidung des Nutzers vom 2026-10-19). Was der Trockenlauf nicht belegt (Schreibrecht des Tokens, Anlegen des GitHub-Release, Sichtbarkeit des Pakets), zeigt erst der Lauf für den ersten Tag; das Risiko trägt `slice-erster-release-freigabe` (dort §6, *Abbruch beim ersten Tag*).
- **`:latest` im Beispiel der Spezifikation.** Entschieden vom Architect am 2026-10-19 (Beispiel mit Registry und Version, Regel zu `latest` im Text von `LH-FA-16.a`); die Zusage und ihre Prüfung (`SPEC-067` Punkt 8) liefert B — **Ausgang:** entfallen für A (Grund: Gegenstand von B, dort §6).
- **Die Bezugs-Abschnitte des Handbuchs.** Der Abschnitt zum Bezug der Binaries vom GitHub-Release steht erst nach dem ersten Tag (Ist-Zustand, `AGENTS.md` §3.11) und ist an `slice-erster-release-freigabe` gesendet (§1) — **Ausgang:** eingetreten, Folge-Slice `slice-erster-release-freigabe`.

**Randformen** (`AGENTS.md` §3.12) — **entschieden vom Architect am 2026-10-19** (vom Nutzer am 2026-10-19 angenommen: [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md); Ort: [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md), Randform → Punkt in der Tabelle dort; Anteil von A, was B zugeordnet ist, steht in `slice-erster-release-veroeffentlichung-image` §6):

- *Träger der Automatisierung:* GitHub Actions, `release.yml`, nur Tags `v*` und `workflow_dispatch`, wörtlicher Text (Fassung `1a`); Pins und Rechte wie `ci.yml` — Punkt 1, [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md).
- *Reihenfolge der Schritte:* `eingaben`, `binaries`, `github-release` (A); die erste Außenwirkung von A ist der GitHub-Release, mit B der Push nach `ghcr.io` — Punkt 2.
- *Anhänge am GitHub-Release:* sechs Binaries und `SHA256SUMS`, sonst nichts; `prerelease` bei `-…`; Text aus der Zeile der Änderungshistorie, fehlt sie, Eingabefehler — Punkt 3 (5), 9.
- *Zugangsdaten (A):* `GITHUB_TOKEN`; kein Wert im Log, kein `set -x` — Punkt 10.
- *Versionsdatei:* `VERSION` im Wurzelverzeichnis, eine Zeile `<SemVer>` ohne `v`, Anfangswert `0.1.0-dev`; das Binary trägt sie über `VERSION=` des Baus, `make binaries` ohne Angabe bleibt bei `dev` — Schnittstelle, Punkt 3, 4.
- *SemVer-Prüfung:* in der Pipeline vor jedem Schritt mit Außenwirkung; Build-Metadaten abgelehnt; führende Nullen abgelehnt — Punkt 3. Ein Prüfen vor `git push` des Tags gibt es nicht (der Hook liegt nicht im Repo).
- *Wiederanlauf und Fehlertabelle:* Punkt 11 (Zeilen `eingaben`, `binaries`, `github-release`); ein Git-Tag wird nie verändert.
- *Probe-Tag:* nicht vorgesehen (Entscheidung des Nutzers vom 2026-10-19); Trockenlauf auf einem Zweig statt dessen — Punkt 2, 12.
- *Plattform, Aufrufarten von `make`, Abbruch, Reihenfolge der Gründe:* Kopf, Punkt 3, 11, 13, 15.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
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

Stand: Kopf der Lieferung ohne die Plan-Commits `4178d430` auf dem Zweig `impl-rel-a` (Basis 11a2ba91); Lauf im Dev-Container pgr-dev3 (Docker im Container), die drei Läufe nacheinander am Endstand.

| Lauf | Ergebnis | Dauer |
|---|---|---|
| `make release-gegenprobe` | `release-gegenprobe: grün (200 Läufe)`: die unveränderte Kopie, 197 Mutationen, der Rekursionsschutz und die Probe der Fallliste; 169 Fälle (96 `p1-*`, 67 `p2-*` bis `p11-*`, 6 `p13-*`) | 161 s |
| `make gates-schnell-gegenprobe` | Exit 0 (Namens- und Zuordnungsregel der Gegenproben mit `release-gegenprobe` in `GEGENPROBEN_CONTAINER`) | 63 s |
| `make gates` | Exit 0 (darunter `workflows-check: 2 Datei(en) geprüft, kein Befund`, `test-container`, `test-postgres-versionen`, `docs-check`, `kopf-check`) | 508 s |

Während der Arbeit nur die engsten Sensoren: `bash -n`, Teilläufe der Gegenprobe mit `GP_MUT` und `GP_ENTDECKE` (die roten Fälle je Mutation gemessen, dann die Erwartung im Skript gesetzt), `make workflows-check` und `make docs-check`.

**Messung der Größe** (`AGENTS.md` §3.14): `git diff --numstat 11a2ba91..4178d430 -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'` ergibt als Summe aus Einfügungen und Löschungen **1303** (`release-gegenprobe.sh` 957, `release.sh` 218, `releasing.md` 76, `release.yml` 38, `release.mk` 12, `harness/README.md` 1, `VERSION` 1). Schätzung des Architect 820 bis 1070; die Gegenprobe trägt 73 % (197 Mutationen mit je einer benannten Menge roter Fälle, 169 Fälle). Meldeschwelle 1000 und Marke 1300 sind überschritten. 2026-10-19 (Implementer): Entscheidung „weiter“ oder „schneiden“ steht aus, der Bericht an die Koordination nennt die Zahl.

**Mutationstabelle** (`AGENTS.md` §3.10): Zusage · Mutation (Textersatz in einer Kopie von `release.sh` bzw. in `pruefe_release`; Namen wie in `tools/harness/release-gegenprobe.sh`) · roter Fall. Je Mutation sind genau die genannten Fälle rot, jede Mutation einmal gesehen (`GP_ENTDECKE=1`), die unveränderte Kopie ist grün. Bei mehr als sechs roten Fällen nennt die Zeile die Zahl und drei Namen; die vollständige Menge steht im Skript.

| Zusage (Punkt von SPEC-067) | Mutation | roter Fall |
|---|---|---|
| 1: jede der 38 Zeilen des Textes steht in der Datei | Zeile n ist optional (`m-p1-zeile-n-optional`, n = 1 bis 38) | `p1-fehlt-n`; bei der Zeile mit `if:` auch `p1-bedingung-github-release` |
| 1: jede Zeile hat genau den Wortlaut | Inhalt der Zeile n wird nicht verglichen (`m-p1-zeile-n-nicht-verglichen`, n = 1 bis 38) | `p1-geaendert-n`; bei der Zeile des Pins auch die fünf `p1-pin-*` |
| 1: kein `pull_request_target` | eine Zeile mit `pull_request_target:` wird überlesen (`m-p1-zusatz-pull-request-target-erlaubt`) | `p1-zusatz-pull-request-target` |
| 1: kein `secrets.ACTIONS_RUNNER_DEBUG` | eine Zeile mit diesem Geheimnis wird überlesen (`m-p1-zusatz-secrets-debug-erlaubt`) | `p1-zusatz-secrets-debug` |
| 1: kein weiterer Schritt mit `uses:` | eine Zeile `actions/cache@` wird überlesen (`m-p1-zusatz-schritt-cache-erlaubt`) | `p1-zusatz-schritt-cache` |
| 1: kein `environment:` | eine Zeile `environment:` wird überlesen (`m-p1-zusatz-environment-erlaubt`) | `p1-zusatz-environment` |
| 1: kein `continue-on-error: true` | diese Zeile wird überlesen (`m-p1-zusatz-continue-on-error-erlaubt`) | `p1-zusatz-continue-on-error` |
| 1/11: kein Schritt läuft nach einem roten (`always()`) | eine Zeile mit `always()` wird überlesen (`m-p1-zusatz-always-erlaubt`) | `p1-zusatz-always` |
| 1a: kein `packages: write` | eine Zeile `packages: write` wird überlesen (`m-p1-zusatz-packages-write-erlaubt`) | `p1-zusatz-packages-write` |
| 1: die Datei hat nur die Zeilen des Textes | die Zahl der Zeilen wird nicht verglichen (`m-p1-zusatzzeilen-erlaubt`) | 7 Fälle, darunter `p1-zusatz-always`, `p1-zusatz-continue-on-error`, `p1-zusatz-environment` |
| 1: Pin aus 40 Hexziffern in Kleinschrift | Großbuchstaben zugelassen (`m-p1-pin-grossbuchstaben`) | `p1-pin-gross` |
| 1: Pin aus genau 40 Hexziffern | 39 Ziffern zugelassen (`m-p1-pin-39-ziffern`) | `p1-pin-kurz` |
| 1: Pin mit Kommentar `# <Tag>` | Kommentar optional (`m-p1-pin-ohne-kommentar`) | `p1-pin-ohne-tag` |
| 1: der Tag-Kommentar ist nicht leer | leerer Kommentar zugelassen (`m-p1-pin-kommentar-ohne-text`) | `p1-pin-tag-leer` |
| 1: der Pin gilt `actions/checkout` | eine andere Aktion zugelassen (`m-p1-pin-andere-action`) | `p1-pin-andere-action` |
| 1: Kommentarzeilen zählen nicht | Kommentarzeilen werden verglichen (`m-p1-kommentarzeilen-nicht-ignoriert`) | `p1-kommentarzeile`, `p1-kommentarzeile-eingerueckt` |
| 1: Leerzeilen zählen nicht | Leerzeilen werden verglichen (`m-p1-leerzeilen-nicht-ignoriert`) | `p1-leerzeile` |
| 1: Kommentarzeile beginnt mit `#` | `#` irgendwo in der Zeile genügt (`m-p1-kommentar-irgendwo-in-der-zeile`) | `p1-gleich`, `p1-kommentarzeile`, `p1-kommentarzeile-eingerueckt`, `p1-leerzeile`, `p1-text-gleich` |
| 1: Kommentarzeile darf eingerückt sein | nur `#` in Spalte 1 zählt (`m-p1-kommentar-nur-am-zeilenanfang`) | `p1-kommentarzeile-eingerueckt` |
| 1: Reihenfolge `Eingaben`, `Binaries` | vertauschte Schritte werden angenommen (`m-p1-reihenfolge-eingaben-binaries-tauschbar`) | `p1-reihenfolge-eingaben-binaries` |
| 1: Reihenfolge `Binaries`, `GitHub-Release` | vertauschte Schritte werden angenommen (`m-p1-reihenfolge-binaries-release-tauschbar`) | `p1-reihenfolge-binaries-release` |
| 2: der Trockenlauf ruft weder `skopeo` noch `gh` | `binaries` ruft im Trockenlauf `skopeo` (`m-p2-trockenlauf-ruft-skopeo`) | `p2-trockenlauf-ohne-aussenwirkung` |
| 2: `github-release` lehnt ohne Push von `refs/tags/v…` ab (je Schritt mit Außenwirkung) | der Aufruf der Prüfung entfällt im Schritt (`m-p2-github-release-ohne-ereignispruefung`) | `p2-aussenwirkung-ohne-push`, `p2-ref-kein-tag` |
| 2/3 (1): Ereignis `push` | die Prüfung des Ereignisses entfällt (`m-p2-ereignis-nicht-geprueft`) | `p2-aussenwirkung-ohne-push`, `p3-ereignis-fehlt`, `p3-ereignis-unbekannt` |
| 2/3 (1): Ref `refs/tags/v…` | die Prüfung des Refs entfällt (`m-p2-ref-nicht-geprueft`) | `p2-ref-kein-tag`, `p3-ohne-v`, `p3-reihenfolge-1-2` |
| 2: Ref `refs/tags/v…` mit `v` | Muster `refs/tags/*` (`m-p2-ref-ohne-v-praefix`) | `p3-ohne-v`, `p3-reihenfolge-1-2` |
| 2/3 (1): `workflow_dispatch` ist der Trockenlauf | Trockenlauf wird wie ein Lauf geprüft (`m-p2-trockenlauf-wie-lauf`) | 7 Fälle, darunter `p2-trockenlauf-ohne-aussenwirkung`, `p3-trockenlauf-gueltig`, `p3-trockenlauf-im-standard-zweig` |
| 3: Prüfung (1) Ereignis und Ref | Aufruf entfernt (`m-p3-1-ereignis-entfernt`) | `p2-ref-kein-tag`, `p3-ereignis-fehlt`, `p3-ereignis-unbekannt`, `p3-ohne-v`, `p3-reihenfolge-1-2` |
| 3: Prüfung (2) Form des Tags | Aufruf entfernt (`m-p3-2-tag-form-entfernt`) | 12 Fälle, darunter `p3-build-metadaten`, `p3-fuehrende-null`, `p3-fuehrende-null-minor` |
| 3: Prüfung (3) Form von `VERSION` | Aufruf entfernt (`m-p3-3-version-datei-entfernt`) | 46 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-reihenfolge-gruende` |
| 3: Prüfung (4) Tag gleich `v` plus `VERSION` | Aufruf entfernt (`m-p3-4-tag-gleich-entfernt`) | `p3-reihenfolge-4-5`, `p3-tag-ungleich-datei`, `p3-tag-ungleich-gross`, `p3-tag-ungleich-tag-kuerzer`, `p3-tag-ungleich-tag-laenger` |
| 3: Prüfung (5) Zeile der Änderungshistorie | Aufruf entfernt (`m-p3-5-notes-entfernt`) | `p3-notes-abschnitt`, `p3-notes-abschnitt-ende`, `p3-notes-fehlen`, `p3-notes-ueberschrift-endet`, `p3-reihenfolge-5-6`, `p3-trockenlauf-notes-fehlen` |
| 3: Prüfung (6) Commit im Standard-Zweig | Aufruf entfernt (`m-p3-6-standard-zweig-entfernt`) | `p3-nicht-im-standard-zweig`, `p3-nicht-im-standard-zweig-vorwaerts`, `p3-trockenlauf-im-standard-zweig`, `p3-trockenlauf-nicht-im-standard-zweig` |
| 3: Reihenfolge der Gründe (1) vor (2) | (1) und (2) vertauscht (`m-p3-reihenfolge-1-2-vertauscht`) | `p2-ref-kein-tag`, `p3-ohne-v`, `p3-reihenfolge-1-2` |
| 3: Reihenfolge der Gründe (2) vor (3) | (2) und (3) vertauscht (`m-p3-reihenfolge-2-3-vertauscht`) | 11 Fälle, darunter `p3-build-metadaten`, `p3-fuehrende-null`, `p3-fuehrende-null-minor` |
| 3: Reihenfolge der Gründe (3) vor (4) | (3) und (4) vertauscht (`m-p3-reihenfolge-3-4-vertauscht`) | 46 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-reihenfolge-gruende` |
| 3: Reihenfolge der Gründe (4) vor (5) | (4) und (5) vertauscht (`m-p3-reihenfolge-4-5-vertauscht`) | `p3-reihenfolge-4-5`, `p3-tag-ungleich-datei`, `p3-tag-ungleich-gross`, `p3-tag-ungleich-tag-kuerzer`, `p3-tag-ungleich-tag-laenger` |
| 3: Reihenfolge der Gründe (5) vor (6) | (5) und (6) vertauscht (`m-p3-reihenfolge-5-6-vertauscht`) | `p3-reihenfolge-5-6` |
| 3 (2): Tag hat die Form `v<SemVer>` | der Vergleich mit der Form entfällt (Länge bleibt) (`m-p3-form-entfernt`) | 11 Fälle, darunter `p3-build-metadaten`, `p3-fuehrende-null`, `p3-fuehrende-null-minor` |
| 3 (2): Länge bis 100 Zeichen | die Prüfung entfällt (`m-p3-laenge-entfernt`) | `p3-zu-lang` |
| 3 (2): 101 Zeichen sind zu lang | Grenze 101 (`m-p3-laenge-101`) | `p3-zu-lang` |
| 3 (2): 100 Zeichen sind zugelassen | Grenze 99 (`m-p3-laenge-99`) | `p3-laenge-100` |
| 3 (2): MAJOR ohne führende Null | `[0-9]+` für MAJOR (`m-p3-null-major`) | `p3-fuehrende-null`, `p3-reihenfolge-2-3` |
| 3 (2): MINOR ohne führende Null | `[0-9]+` für MINOR (`m-p3-null-minor`) | `p3-fuehrende-null-minor` |
| 3 (2): PATCH ohne führende Null | `[0-9]+` für PATCH (`m-p3-null-patch`) | `p3-fuehrende-null-patch` |
| 3 (2): `MAJOR.MINOR.PATCH`, drei Teile | PATCH optional (`m-p3-kern-zwei-teile`) | `p3-zwei-teile` |
| 3 (2): genau drei Teile | vierter Teil zugelassen (`m-p3-kern-vier-teile`) | `p3-vier-teile` |
| 3 (2): Vorabversion nach `-` nicht leer | `v1.2.3-` zugelassen (`m-p3-vor-leer-erlaubt`) | `p3-leere-vorabversion` |
| 3 (2): Bezeichner der Vorabversion nicht leer | leerer Bezeichner zwischen Punkten zugelassen (`m-p3-vor-leerer-bezeichner`) | `p3-leerer-bezeichner` |
| 3 (2): keine Build-Metadaten | `+…` zugelassen (`m-p3-vor-build-metadaten`) | `p3-build-metadaten` |
| 3 (2): Bezeichner durch `.` getrennt | nur ein Bezeichner zugelassen (`m-p3-vor-ein-bezeichner`) | `p3-tag-gueltig`, `p3-tag-ungleich-gross`, `p3-tag-ungleich-tag-kuerzer`, `p3-tag-ungleich-tag-laenger`, `p9-vorabversion-prerelease` |
| 3 (2): numerischer Bezeichner ohne führende Null | jede Zeichenfolge zugelassen (`m-p3-bez-fuehrende-null`) | `p3-vorabversion-fuehrende-null` |
| 3 (2): Zeichen `0-9A-Za-z-` | `_` zugelassen (`m-p3-bez-zeichen`) | `p3-ungueltiges-zeichen` |
| 3 (2): Bezeichner mit Buchstaben | nur Zahlen zugelassen (`m-p3-bez-nur-zahlen`) | 9 Fälle, darunter `p2-trockenlauf-ohne-aussenwirkung`, `p3-laenge-100`, `p3-tag-gueltig` |
| 3 (2): Ziffern in alphanumerischen Bezeichnern | Ziffern nicht zugelassen (`m-p3-bez-ohne-ziffern`) | `p3-tag-gueltig` |
| 3 (2): Großbuchstaben | nicht zugelassen (`m-p3-bez-ohne-grossbuchstaben`) | `p3-tag-gueltig`, `p3-tag-ungleich-gross` |
| 3 (2): Kleinbuchstaben | nicht zugelassen (`m-p3-bez-ohne-kleinbuchstaben`) | 9 Fälle, darunter `p2-trockenlauf-ohne-aussenwirkung`, `p3-laenge-100`, `p3-tag-gueltig` |
| 3 (2): Bindestrich im Bezeichner | nicht zugelassen (`m-p3-bez-ohne-bindestrich`) | `p3-tag-gueltig` |
| 3 (2): Präfix `v` | `v` optional (äquivalent, siehe unten) (`m-p3-form-v-optional`) | keiner (äquivalent) |
| 3 (3): Form von `VERSION` | der Vergleich mit der Form entfällt (`m-p3-vform-entfernt`) | 9 Fälle, darunter `p3-reihenfolge-3-4`, `p3-trockenlauf-version-datei-form`, `p3-version-datei-form-crlf` |
| 3 (3): LF am Ende | Zeilenende optional (`m-p3-vdatei-kein-lf`) | `p3-trockenlauf-version-datei-form`, `p3-version-datei-form-kein-lf` |
| 3 (3): LF, nicht CRLF | CR vor LF zugelassen (`m-p3-vdatei-crlf`) | `p3-reihenfolge-3-4`, `p3-version-datei-form-crlf` |
| 3 (3): keine Leerzeichen | Leerzeichen vorn und hinten zugelassen (`m-p3-vdatei-leerzeichen`) | `p3-version-datei-form-leerzeichen-hinten`, `p3-version-datei-form-leerzeichen-vorn` |
| 3 (3): `<SemVer>` ohne `v` | `v` zugelassen (`m-p3-vdatei-v`) | `p3-version-datei-form-v` |
| 3 (3): genau eine Zeile | mehrere Zeilen zugelassen (`m-p3-vdatei-mehrere-zeilen`) | `p3-version-datei-form-zwei-zeilen` |
| 3 (4): Vergleich `v` plus `VERSION` | Vergleich ohne `v` (`m-p3-gleich-ohne-v`) | 37 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-reihenfolge-gruende` |
| 3 (4): Zeichen für Zeichen | Groß- und Kleinschrift gleich (`m-p3-gleich-gross`) | `p3-tag-ungleich-gross` |
| 3 (4): Zeichen für Zeichen | Tag darf `v<VERSION>` fortsetzen (`m-p3-gleich-praefix-tag-laenger`) | `p3-tag-ungleich-tag-laenger` |
| 3 (4): Zeichen für Zeichen | Tag darf Anfang von `v<VERSION>` sein (`m-p3-gleich-praefix-tag-kuerzer`) | `p3-tag-ungleich-tag-kuerzer` |
| 3 (5): Tabelle im Abschnitt `Änderungshistorie` | jede Tabelle der Datei zählt (`m-p3-notes-abschnitt-ignoriert`) | `p3-notes-abschnitt`, `p3-notes-abschnitt-ende`, `p3-notes-ueberschrift-endet`, `p9-notes-aus-abschnitt` |
| 3 (5): Überschrift endet auf `Änderungshistorie` | Überschrift enthält das Wort (`m-p3-notes-ueberschrift-ende-ignoriert`) | `p3-notes-ueberschrift-endet` |
| 3 (5): nur die Tabelle dieses Abschnitts | Abschnitt endet nicht an der nächsten Überschrift (`m-p3-notes-abschnitt-ende-ignoriert`) | `p3-notes-abschnitt-ende` |
| 3 (5)/9: Zeile mit erster Zelle gleich der Version | Anfang der ersten Zelle genügt (`m-p3-notes-praefix`) | `p9-notes-zeile-der-version` |
| 9: Text ist die dritte Zelle | zweite Zelle (`m-p3-notes-zelle-datum`) | `p9-notes-aus-abschnitt`, `p9-notes-aus-zelle`, `p9-notes-zeile-der-version` |
| 3 (6): Commit liegt im Standard-Zweig | `merge-base` ohne `--is-ancestor` (`m-p3-anc-ohne-is-ancestor`) | `p3-nicht-im-standard-zweig`, `p3-nicht-im-standard-zweig-vorwaerts`, `p3-trockenlauf-nicht-im-standard-zweig` |
| 3 (6): Commit liegt im Standard-Zweig | Argumente vertauscht (`m-p3-anc-umgekehrt`) | `p3-nicht-im-standard-zweig`, `p3-nicht-im-standard-zweig-vorwaerts`, `p3-trockenlauf-nicht-im-standard-zweig` |
| 3 (6): Standard-Zweig aus `RELEASE_DEFAULT_BRANCH` | fest `origin/main` (`m-p3-standard-zweig-main`) | `p3-standard-zweig-aus-umgebung` |
| 3 (6): beim Lauf Eingabefehler | wie im Trockenlauf nur ein Hinweis (`m-p3-lauf-nur-hinweis`) | `p3-nicht-im-standard-zweig`, `p3-nicht-im-standard-zweig-vorwaerts` |
| 3 (6): im Trockenlauf kein Fehler | Trockenlauf wird wie der Lauf behandelt (`m-p3-trockenlauf-fehler`) | `p3-trockenlauf-im-standard-zweig`, `p3-trockenlauf-nicht-im-standard-zweig` |
| 3 (6): der Hinweis nennt die Antwort | immer „liegt im“ (`m-p3-hinweis-immer-ja`) | `p3-trockenlauf-nicht-im-standard-zweig` |
| 3 (6): im Trockenlauf steht die Antwort als Hinweis | kein Hinweis (`m-p3-hinweis-fehlt`) | `p3-trockenlauf-im-standard-zweig`, `p3-trockenlauf-nicht-im-standard-zweig` |
| 2: Tag des Trockenlaufs ist `v<VERSION>` | Tag aus dem Ref (`m-p3-trockenlauf-tag-aus-ref`) | 7 Fälle, darunter `p2-trockenlauf-ohne-aussenwirkung`, `p3-trockenlauf-gueltig`, `p3-trockenlauf-im-standard-zweig` |
| 4: Aufruf mit `VERSION=<Version>` | `VERSION=` fehlt im Aufruf (`m-p4-version-fehlt`) | 18 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
| 4: `DIST` außerhalb des Arbeitsbaums | `DIST=dist` (`m-p4-dist-im-baum`) | 18 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
| 4: `DIST` ist neu | das Skript legt das Verzeichnis vorher an (`m-p4-dist-vorher-angelegt`) | `p4-aufruf` |
| 4: `sha256sum -c` grün | die Prüfung entfällt (`m-p4-summe-entfernt`) | `p10-nicht-im-log`, `p11-alte-marke`, `p11-reihenfolge-gruende`, `p4-summe-falsch` |
| 4: genau die sechs Dateien plus `SHA256SUMS` | die Dateiliste wird nicht geprüft (`m-p4-liste-entfernt`) | `p4-datei-fehlt`, `p4-datei-zu-viel` |
| 4/11: ein Bau, der scheitert, ist rot | Exit von `make` wird überlesen (`m-p4-make-fehler-ignoriert`) | `p4-make-scheitert` |
| 9: Anhänge genau die sieben Dateien | ein achter Anhang (`m-p9-achter-anhang`) | `p11-wiederanlauf-release-teilweise`, `p11-wiederanlauf-release-vollstaendig`, `p9-anhaenge`, `p9-release-existiert-gleich` |
| 9: `--prerelease` genau bei einer Vorabversion | immer (`m-p9-prerelease-immer`) | `p9-stabil-ohne-prerelease` |
| 9: `--prerelease` genau bei einer Vorabversion | nie (`m-p9-prerelease-nie`) | `p9-vorabversion-prerelease` |
| 9: `--verify-tag` | fehlt (`m-p9-ohne-verify-tag`) | `p9-aufruf` |
| 9: `--title v<Version>` | fehlt (`m-p9-ohne-title`) | `p9-aufruf` |
| 9: `--notes-file` | fehlt (`m-p9-ohne-notes-file`) | `p9-aufruf`, `p9-notes-aus-abschnitt`, `p9-notes-aus-zelle`, `p9-notes-zeile-der-version` |
| 9: nicht als Entwurf | `--draft` (`m-p9-entwurf`) | `p9-aufruf` |
| 9: abweichende `SHA256SUMS` am Release ist rot | der Vergleich entfällt (`m-p9-vergleich-entfernt`) | `p10-nicht-im-log`, `p9-release-existiert-abweichend` |
| 9: `SHA256SUMS` wird nie überschrieben | Hochladen mit `--clobber` statt Fehler (`m-p9-ueberschreiben-statt-vergleichen`) | `p10-nicht-im-log`, `p9-release-existiert-abweichend` |
| 9/11: fehlende Anhänge ohne Überschreiben | `--clobber` beim Hochladen (`m-p9-upload-clobber`) | `p11-wiederanlauf-release-teilweise`, `p9-release-existiert-gleich` |
| 9/11: vorhandener Release wird ergänzt, nicht neu angelegt | `create` auch bei vorhandenem Release (`m-p9-create-trotz-vorhandenem-release`) | `p10-nicht-im-log`, `p11-wiederanlauf-release-teilweise`, `p11-wiederanlauf-release-vollstaendig`, `p9-release-existiert-abweichend`, `p9-release-existiert-gleich` |
| 9/11: nur fehlende Anhänge werden hochgeladen | alle Anhänge (`m-p9-upload-alle`) | `p11-wiederanlauf-release-teilweise`, `p11-wiederanlauf-release-vollstaendig`, `p9-release-existiert-gleich` |
| 9: Fehler von `gh release view` ist rot | jeder Fehler gilt als „nicht gefunden“ (`m-p9-view-fehler-wie-fehlend`) | `p10-nicht-im-log`, `p9-view-fehler` |
| 10: kein Wert des Tokens in der Ausgabe (Fehler im Vergleich) | das Skript gibt `$GH_TOKEN` aus (`m-p10-token-bei-fehler`) | `p10-nicht-im-log` |
| 10: kein Wert des Tokens in der Ausgabe (Fehler der Summe) | das Skript gibt `$GH_TOKEN` aus (`m-p10-token-bei-fehler-binaries`) | `p10-nicht-im-log` |
| 10: kein `set -x` | `set -x` im Skript (`m-p10-set-x`) | `p10-nicht-im-log` |
| 10: kein `env` | `env` im Fehlerfall (`m-p10-env-bei-fehler`) | `p10-nicht-im-log` |
| 11: `binaries` läuft nicht nach rotem `eingaben` | die Marke wird nicht geprüft (`m-p11-binaries-ohne-eingaben-marke`) | `p11-alte-marke`, `p11-reihenfolge-gruende` |
| 11: `github-release` läuft nicht nach rotem `binaries` | die Marke wird nicht geprüft (`m-p11-release-ohne-binaries-marke`) | `p11-alte-marke`, `p11-reihenfolge-gruende` |
| 11: kein Schritt läuft mit der Marke eines früheren Laufs | `eingaben` löscht die Marken nicht (`m-p11-eingaben-loescht-marken-nicht`) | `p11-alte-marke` |
| 11: kein Schritt läuft mit der Marke eines früheren Laufs | `binaries` löscht seine Marke nicht (`m-p11-binaries-loescht-marke-nicht`) | `p11-alte-marke` |
| 13: `release.sh` räumt bei Ende | die Falle für EXIT fehlt (`m-p13-ende-falle-fehlt`) | `p13-ende` |
| 13: `release.sh` räumt bei `SIGINT` | die Falle fehlt (`m-p13-int-falle-fehlt`) | `p13-sigint` |
| 13: `release.sh` räumt bei `SIGTERM` | die Falle fehlt (`m-p13-term-falle-fehlt`) | `p13-sigterm` |
| 13: Unterprozesse zuerst beendet (`SIGINT`) | sie werden nicht beendet (`m-p13-int-kinder-bleiben`) | `p13-sigint` |
| 13: Unterprozesse zuerst beendet (`SIGTERM`) | sie werden nicht beendet (`m-p13-term-kinder-bleiben`) | `p13-sigterm` |
| 13: erst Unterprozesse, dann Verzeichnis (`SIGINT`) | Verzeichnis vor dem Ende der Unterprozesse entfernt (`m-p13-int-verzeichnis-zu-frueh`) | `p13-sigint` |
| 13: erst Unterprozesse, dann Verzeichnis (`SIGTERM`) | Verzeichnis vor dem Ende der Unterprozesse entfernt (`m-p13-term-verzeichnis-zu-frueh`) | `p13-sigterm` |
| 13: die Gegenprobe räumt bei Ende | die Falle für EXIT fehlt (`m-p13-gegenprobe-ende-falle-fehlt`) | `p13-gegenprobe-ende` |
| 13: die Gegenprobe räumt bei `SIGINT` | die Falle fehlt (`m-p13-gegenprobe-int-falle-fehlt`) | `p13-gegenprobe-sigint` |
| 13: die Gegenprobe räumt bei `SIGTERM` | die Falle fehlt (`m-p13-gegenprobe-term-falle-fehlt`) | `p13-gegenprobe-sigterm` |
| 13: Unterprozesse der Gegenprobe zuerst beendet (`SIGINT`) | sie werden nicht beendet (`m-p13-gegenprobe-int-kinder-bleiben`) | `p13-gegenprobe-sigint` |
| 13: Unterprozesse der Gegenprobe zuerst beendet (`SIGTERM`) | sie werden nicht beendet (`m-p13-gegenprobe-term-kinder-bleiben`) | `p13-gegenprobe-sigterm` |
| 13: erst Unterprozesse, dann Verzeichnis (Gegenprobe, `SIGINT`) | Verzeichnis zu früh entfernt (`m-p13-gegenprobe-int-verzeichnis-zu-frueh`) | `p13-gegenprobe-sigint` |
| 13: erst Unterprozesse, dann Verzeichnis (Gegenprobe, `SIGTERM`) | Verzeichnis zu früh entfernt (`m-p13-gegenprobe-term-verzeichnis-zu-frueh`) | `p13-gegenprobe-sigterm` |

**Ein grüner Mutant, äquivalent** (`AGENTS.md` §3.10): `m-p3-form-v-optional` (die Form des Tags verlangt das `v` nicht mehr) färbt keinen Fall. (a) Invariante: Wenn `pruefe_tag_form` läuft, beginnt der Tag mit `v`. (b) Ort: beim Lauf `aussenwirkung_erlaubt` über `pruefe_ereignis` (Ref `refs/tags/v*`, Exit 2 sonst), beim Trockenlauf `tag_ableiten` (`v` plus erste Zeile von `VERSION`); beide sind Zusage dieses Slice und eigene Mutationen (`m-p2-ref-ohne-v-praefix`, `m-p3-trockenlauf-tag-aus-ref`) färben Fälle. (c) Sonde: Eingaben `push` mit Ref `refs/tags/1.2.3`, `VERSION` `1.2.3`, Zeile der Version im Handbuch, gegen das Original und den Mutanten: beide Exit 2 mit dem Kennwort `ereignis`. Gegenprobe der Sonde: Ohne die Prüfung des Refs (Invariante aufgehoben) endet das Original mit `tag-form` und der Mutant mit `tag-ungleich`, beide Exit 2; die Sonde trennt also Original und Mutant am Kennwort, sobald die Invariante nicht mehr gilt. Gemessen am 2026-10-19 mit `bash <Kopie von release.sh> eingaben` im Dev-Container, Umgebung `RELEASE_EVENT=push`, `RELEASE_REF=refs/tags/1.2.3`: Original Kennwort `ereignis`, Exit 2 · Mutant `ereignis`, Exit 2 · Original ohne die Zeile `# P2ref` `tag-form`, Exit 2 · Mutant ohne die Zeile `# P2ref` `tag-ungleich`, Exit 2.

**Trockenlauf auf GitHub:** nicht gefahren (der Auftrag des Implementers verbietet einen Push auf GitHub); der Beleg der DoD (Lauf-Kennung, Dauer, `SPEC-067` Punkt 12) steht aus. Gemessen auf dem Stand der Lieferung mit den Eingaben eines Trockenlaufs (`RELEASE_EVENT=workflow_dispatch`, `VERSION` `0.1.0-dev`): `bash tools/harness/release.sh eingaben` endet mit `release: eingaben: notes: die Änderungshistorie des Benutzerhandbuchs hat keine Zeile mit Text für 0.1.0-dev`, Exit 2, weil das Handbuch (letztes Kapitel: „Es gibt keine veröffentlichte Version.“) keine Tabelle führt.

**Handbuch und README:** keine Zeile neu oder geändert. Gelesen gegen `make binaries` und `make test-container` (im grünen `make gates`): Voraussetzungen (sechs Plattformen, nur `linux/amd64` ausgeführt), Installation (`make binaries`, Namen mit `VERSION`, `SHA256SUMS`, `dist/` leer oder fehlend) und der Absatz von `README.md` stehen schon im Ist-Zustand; der Abschnitt zum Bezug vom GitHub-Release gehört `slice-erster-release-freigabe` (§1).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `0676302` nachgesichtet, mit der Übernahme
der Entscheidungen vom 2026-10-09. Treffer:

- `BEO-REPO/slice-waechst-durch-uebernahmen` (2×) — der Slice nahm sieben Entscheidungen auf.
  Kein Beleg: Die Übernahmen standen als eigene Liefer-Punkte in §2, die Zählung zeigte die
  Größe (sechs), und der Schnitt vom 2026-10-09 lief vor dem Start; das Muster ist das Wachsen
  *ohne* mehr Punkte.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (16×, verkörpert in `AGENTS.md` §3.12) —
  die Randformen der neuen Punkte stehen offen in §6 und gehen vor dem Code an den Architect.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (17×, `AGENTS.md` §3.10) — die
  Prüfungen der Pipeline (Tag, Versionsdatei, Digest-Vergleich) sind neue Verträge; je
  Zusage eine Mutation.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Nachgezählt beim Eintragen der Schichtzählung aus `slice-doku-ist-stand` (2026-10-10, Entscheidung des Nutzers, `AGENTS.md` §3.13): Die Dokumentation ist eine Schicht; der Handbuch- und README-Teil liegt im ersten Liefer-Punkt und in denselben zwei Dateien, kein neuer Liefer-Punkt. Schichten: zwei Schichten (Release-Automatisierung, Dokumentation).

**Nachgezählt am 2026-10-17** (Planner, Vorbereitung des Anspruchs, `AGENTS.md` §3.13): Seit dem Eintrag von `slice-v1-abschluss-container` (Voraussetzung, keine Sendung) und der Metrik in §4 hat sich nichts an den Liefer-Punkten geändert. Liefer-Punkte: drei (§2: Image mit Manifestliste und Digest-Vergleich | Binaries mit `SHA256SUMS` am GitHub-Release | Prüfung von Tag und Versionsdatei samt Fehlertabelle und Wiederanlauf); Handbuch- und `README.md`-Teil liegt im ersten Punkt. Schichten: zwei (Release-Automatisierung; Dokumentation: `docs/user/`, `docs/maintainer/`, `README.md`). Die Randformen aus §6 (Träger, Reihenfolge, Manifestliste) können einen dritten Bereich öffnen (Dockerfile, wenn das Image das Binary aus `dist/` übernimmt); der Architect zählt es mit der Entscheidung nach, und läge der Slice darüber, schneidet der Planner (§4).

**Nachgezählt am 2026-10-18** (Planner, `AGENTS.md` §3.13, im Commit, der die Sendung aus `slice-harness-ci-workflows` und `slice-harness-ci` einträgt): Die Sendung ist ein Akzeptanzkriterium der Datei der Pipeline (Pins, Rechte, Zeitgrenze, `concurrency`), die Punkt 1 und 2 schon liefern, und der Zusatz zu Punkt 3 nennt es; Liefer-Punkte weiter drei (Image mit Manifestliste und Digest-Vergleich | Binaries mit `SHA256SUMS` am GitHub-Release | Prüfung von Tag und Versionsdatei samt Fehlertabelle und Wiederanlauf). Schichten weiter zwei (Release-Automatisierung; Dokumentation). Die Größe ändert sich nicht (die Prüfung der Form trägt `make workflows-check`, die Datei besteht sie oder nicht); die Schätzung in §4 bleibt, der Architect zählt mit der Entscheidung der Randformen nach.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit, der die Reihenfolge der Welle auf den kritischen Pfad kürzt; Entscheidung des Nutzers vom 2026-10-18: der erste Release hat Vorrang vor den Nachtläufen): Die Sendung aus `slice-harness-ci-workflows` und `slice-harness-ci` (Pins und Härtung) bleibt unverändert; der Slice steht jetzt an dritter Stelle der Welle (Schritt 3) und setzt weiter nur diese beiden voraus. Die Nachtläufe (`slice-harness-ci-gegenproben`, `-ci-freshness`, `-ci-freshness-images`) liegen in [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md) hinter dieser Welle und sind keine Voraussetzung. Der Inhalt der Sendung ändert sich nicht; Liefer-Punkte weiter drei, Schichten weiter zwei (Release-Automatisierung, Dokumentation); Schätzung und Marke in §4 unverändert.

**Nachgezählt am 2026-10-19** (Planner, `AGENTS.md` §3.13, im Commit, der den Schnitt in A und B vollzieht und die Sendungen an `slice-erster-release-veroeffentlichung-image` und `slice-erster-release-freigabe` einträgt; Basis: `SPEC-067` in der Fassung vom 2026-10-19, Zuordnung der Punkte in der Tabelle des Vertrags): Liefer-Punkte drei (§2: Eingaben samt `VERSION` | Binaries am GitHub-Release samt Handbuch- und `README.md`-Teil | Trockenlauf, Fehlertabelle und Wiederanlauf); `release.yml` Fassung `1a` ist Akzeptanzkriterium von Punkt 2 und 3, kein vierter Punkt. Schichten zwei (Release-Automatisierung: Workflow, `tools/harness/`, `harness/mk/`, `VERSION`; Dokumentation: `docs/maintainer/`, `docs/user/`, `README.md`); die Gegenprobe zählt mit der Automatisierung. Die Abgaben an B und an `slice-erster-release-freigabe` nehmen aus diesem Slice, sie fügen nichts hinzu; die Marke 1300 (Schätzung 820 bis 1070, Mitte 73 % der Marke) steht in §4. Der ungeteilte Slice (Image mit Manifestliste und Digest-Vergleich als erster Punkt) lag bei 1970 bis 2570 gegen 1800 und ist damit nicht mehr zu führen.
