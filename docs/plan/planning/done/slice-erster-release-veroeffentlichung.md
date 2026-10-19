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

**Sendung aus `slice-harness-ci-workflows` und `slice-harness-ci` — Pins und Härtung** (Planner, 2026-10-18, `AGENTS.md` §3.13; Auftrag des Nutzers vom 2026-10-18: die CI gehört vor diesen Slice): Die Datei der Pipeline, `.github/workflows/release.yml` (Name: Architect), fällt vom ersten Commit an unter `make workflows-check` und `make docs-check`. Sie trägt damit dieselben Konventionen wie `ci.yml`: jede fremde `uses:`-Referenz per vollem SHA mit Tag-Kommentar, `permissions: {}` oben und im Job nur die Rechte, die er braucht (in Fassung `1a`: `contents: write` für den GitHub-Release; `packages: write` für `ghcr.io` kommt mit B), `timeout-minutes`, `concurrency` mit der Gruppe `release-${{ github.event_name }}` (Lauf und Trockenlauf getrennt) und `cancel-in-progress: false` (ein **laufender** Tag-Lauf wird nie durch einen späteren abgebrochen; ein wartender kann von einem dritten ersetzt werden, Grenze in `SPEC-067` Punkt 13). `ci.yml` lässt Tags aus (`tags-ignore`), `release.yml` reagiert nur auf `v*`; ein Tag läuft also nur hier. Zugangsdaten nur als Geheimnisse von GitHub oder als `GITHUB_TOKEN`, nie im Repository. Welche Eigenschaften der Datei ein Test hält, legt `SPEC-067` Punkt 1 fest (Text der Datei, Zeilenliste in der Gegenprobe), die Prüfung der Form hält das Gate. Der **Träger** (GitHub Actions) ist auch für die Veröffentlichung entschieden ([ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md), Randform *Träger der Automatisierung*, §6). Voraussetzung ist beim Start: `.github/workflows/` liegt (aus `slice-harness-ci`).

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

- [x] **Liefer-Punkt 1, Eingaben** ([`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), `SPEC-067` Punkt 3 (1) bis (6), Schnittstelle): Vor jedem Schritt mit Außenwirkung prüft `release.sh eingaben`, in der Reihenfolge des Vertrags, dass Ereignis und Ref stimmen, der Tag die Form `v<SemVer>` hat, die Datei `VERSION` (neu, Wurzel, `0.1.0-dev`) ihre Form hat und dem Tag entspricht, die Änderungshistorie des Benutzerhandbuchs die Zeile der Version trägt und der Commit des Tags im Standard-Zweig liegt; sonst Exit 2, bevor etwas gebaut oder veröffentlicht wird (Test: die Fälle `p3-*` der Gegenprobe mit je einer Mutation, `AGENTS.md` §3.10). `make binaries` übernimmt die Version aus `VERSION=`, ohne Angabe bleibt es bei `dev`.
- [x] **Liefer-Punkt 2, Binaries am GitHub-Release** ([`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), `SPEC-067` Punkte 4 und 9): Die Binaries für Linux, macOS und Windows (`amd64`, `arm64`) hängen mit einer Datei `SHA256SUMS` am GitHub-Release (`LH-FA-19.a`); die Dateien am Release sind die aus `make binaries` mit `VERSION` aus dem Tag, genau sieben, und eine abweichende Summe am Release endet rot (Test mit Attrappen von `make` und `gh`: Fälle `p4-*` und `p9-*`; die Gleichheit zweier Builds belegt `make test-container`, [`SPEC-055`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md) Punkt 4, nicht dieser Slice). Die Release-Notes entstehen aus der neuen Zeile der Änderungshistorie des Handbuchs; ein `CHANGELOG.md` wird nicht geführt. Dass der GitHub-Release auf GitHub angelegt wird, belegt erst der Lauf für den ersten Tag (`SPEC-067` Punkt 12); der Trockenlauf legt keinen an. Handbuch §1 *Voraussetzungen* und §2 *Installation* und `README.md` sind auf das geprüft, was der Trockenlauf und `make binaries` belegen (Plattformen der Binaries, `SHA256SUMS`), und im Ist-Zustand des gebauten Binaries berichtigt, wo sie es nicht sind: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11); der Abschnitt zum Bezug vom GitHub-Release gehört `slice-erster-release-freigabe`.
- [x] **Liefer-Punkt 3, Trockenlauf, Wiederanlauf, Fehlertabelle** (`SPEC-067` Punkte 1a, 2, 11, 12): Die Datei `.github/workflows/release.yml` hat den Text der Fassung `1a`, reagiert nur auf `v*` und `workflow_dispatch`, und der Trockenlauf fährt `Eingaben` und `Binaries` ohne Außenwirkung (kein `gh`, kein GitHub-Release; Test: `p1-*`, `p2-*` mit Mutationen). Der Trockenlauf ist einmal auf dem Stand der Lieferung gefahren und in §7 belegt (Lauf-Kennung, Dauer, `SPEC-067` Punkt 12; Handlung des Betreibers oder des Nutzers, falls der Lauf Rechte am Repository braucht). `docs/maintainer/releasing.md` §11 trägt die Fehlertabelle mit den Zeilen `eingaben`, `binaries` und `github-release` und den Wiederanlauf (Punkt 8; die übrigen Zeilen liefert B); der Wiederanstoß eines roten Laufs auf GitHub ist hier nicht belegt, weil ein Trockenlauf grün gebaut ist und keinen roten Lauf erzeugt: Der Beleg ist an `slice-erster-release-freigabe` abgegeben (§6, dort §1), die Attrappen-Fälle `p11-wiederanlauf-*` belegen die Logik. Die Datei besteht `make workflows-check` und `make docs-check` und trägt die Pins, Rechte und Zeitgrenzen der Sendung aus `slice-harness-ci-workflows` und `slice-harness-ci` (§1); sie ist Akzeptanzkriterium dieses Punkts und von Punkt 2, kein vierter Liefer-Punkt.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.
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

- Der Trockenlauf braucht Rechte oder Einstellungen auf GitHub, die aus dem Klon nicht auditierbar sind (Actions im Repository erlaubt, `GITHUB_TOKEN` mit `contents: write`; `SPEC-067` Punkt 15) — **Ausgang:** entfallen (Grund: Der Trockenlauf lief am 2026-10-19 auf GitHub grün, Lauf `39208276024`, mit `permissions: {}` oben und `contents: write` im Job; Belege in §7. Was der Trockenlauf nicht berührt, ist der nächste Punkt).
- Die Zugangsdaten für `ghcr.io` und `docker.io` sind externe Betreiber-Handlungen — sie gehören B (dort §6); A liest nur `GITHUB_TOKEN`. **Ausgang:** entfallen für A (Grund: die Fassung `1a` hat keinen Schritt mit Registry; die Handlung steht bei `slice-erster-release-veroeffentlichung-image`).
- **Größe.** Mit den sieben Entscheidungen vom 2026-10-09 trug §2 sechs Liefer-Punkte, die
  Grenze ist drei. Entscheidung des Nutzers vom 2026-10-09: geschnitten nach dem Vorschlag des
  Planner (§4) — **Ausgang:** eingetreten, Folge-Slice `slice-erster-release-freigabe`. Der Schnitt vom 2026-10-19 (Größe 1970 bis 2570 gegen die Marke 1800) hat die Naht der Rückführung genutzt — **Ausgang:** eingetreten, Folge-Slice `slice-erster-release-veroeffentlichung-image`. Die Marke 1300 der Rückführung in §4 ist mit 1425 überschritten; die Entscheidung „weiter“ steht als datierte Zeile in §7, Review und Verifikation haben den Diff geprüft, ein weiterer Schnitt entfällt (Grund dort).
- **Probe-Tag in öffentlichen Registries.** Entfällt mit der Entscheidung des Nutzers vom 2026-10-19: kein Probe-Tag, der Trockenlauf schreibt nichts in eine Registry und legt keinen GitHub-Release an (`SPEC-067` Punkt 2 und 12) — **Ausgang:** entfallen (Grund: die Entscheidung des Nutzers vom 2026-10-19). Was der Trockenlauf nicht belegt (Schreibrecht des Tokens, Anlegen des GitHub-Release, Sichtbarkeit des Pakets), zeigt erst der Lauf für den ersten Tag; das Risiko trägt `slice-erster-release-freigabe` (dort §6, *Abbruch beim ersten Tag*).
- **`:latest` im Beispiel der Spezifikation.** Entschieden vom Architect am 2026-10-19 (Beispiel mit Registry und Version, Regel zu `latest` im Text von `LH-FA-16.a`); die Zusage und ihre Prüfung (`SPEC-067` Punkt 8) liefert B — **Ausgang:** entfallen für A (Grund: Gegenstand von B, dort §6).
- **Die Bezugs-Abschnitte des Handbuchs.** Der Abschnitt zum Bezug der Binaries vom GitHub-Release steht erst nach dem ersten Tag (Ist-Zustand, `AGENTS.md` §3.11) und ist an `slice-erster-release-freigabe` gesendet (§1) — **Ausgang:** eingetreten, Folge-Slice `slice-erster-release-freigabe`.
- **Was der Trockenlauf nicht belegt** (Verifikation V-1, Review F-7): (a) der Wiederanstoß eines roten Laufs auf GitHub nach der Fehlertabelle von `docs/maintainer/releasing.md` §11, weil ein Trockenlauf grün gebaut ist; (b) der Wortlaut `release not found` von `gh release view` für einen nicht vorhandenen Release (akzeptierte Grenze, `SPEC-067` Punkt 9); (c) das Anlegen des GitHub-Release mit `contents: write` (`SPEC-067` Punkt 12). Der Lauf für den ersten echten Tag zeigt (b) und (c) und (a), falls er rot endet — **Ausgang:** eingetreten, Folge-Slice `slice-erster-release-freigabe` (Sendung des Planners vom 2026-10-19, dort §1, DoD 3, §6 und §8, `AGENTS.md` §3.13).

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

**Entschieden vom Architect am 2026-10-19 nach der Rückgabe des Implementers** ([`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 3, 9, 13): Der Trockenlauf wertet die Zeile der Änderungshistorie aus, eine fehlende Zeile ist dort nur ein Hinweis (Handbuch bleibt Ist-Zustand), eine leere dritte Zelle ein Fehler; die Freigabe fährt den Trockenlauf mit Zeile und `VERSION` des Tags einmal vor dem Tag. Zustand der Schritte in `$RUNNER_TEMP/release` mit Marken. Schlüsselwörter der Meldungen zugesagt, übriger Wortlaut nicht; Vorabversions-Bezeichner nach SemVer 2.0 (`-01` abgelehnt); `gh` mit `--repo`, ohne `GH_TOKEN` Exit 2, nur `release not found` heißt „nicht vorhanden“.

**Entschieden vom Architect am 2026-10-19 nach dem Review** ([`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 1, 3, 9, 13): Die Gruppe von `concurrency` heißt `release-${{ github.event_name }}` (ein Trockenlauf verdrängt keinen Lauf); „ein Tag wird nie abgebrochen“ gilt nur für den laufenden Lauf, ein wartender kann von einem dritten ersetzt werden (Grenze, Weg: Re-run in der Oberfläche); die Zelle der Änderungshistorie hat kein `|` (auch nicht maskiert), sonst Eingabefehler `notes`. Der Implementer ändert die Zeilen `group:` in Text 1a und in den Fällen `p1-*` sowie die Prüfung der Zelle.

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

Stand: Kopf des Codes `39d964e9` auf dem Zweig `impl-rel-a` (Basis 11a2ba91, nach den Merges der Entscheidungen des Architect `5f9fddb5` und `eec9c839` und des Review-Reports); Lauf im Dev-Container pgr-dev3 (Docker im Container), die Läufe nacheinander am Endstand.

| Lauf | Ergebnis | Dauer |
|---|---|---|
| `make release-gegenprobe` | `release-gegenprobe: grün (226 Läufe)`: die unveränderte Kopie, 223 Mutationen, der Rekursionsschutz und die Probe der Fallliste; 183 Fälle (96 `p1-*`, 81 `p2-*` bis `p13-*` ohne den Abbruch, 6 `p13-*` des Abbruchs) | 234 s |
| `make workflows-check` | Exit 0, `workflows-check: 2 Datei(en) geprüft, kein Befund` | 1 s |
| `make docs-check` | Exit 0, `d-check: 1208 Datei(en) geprüft, 0 Befund(e)` | 2 s |
| `make gates` | Exit 0 (darunter `test-container`, `test-postgres-versionen`, `kopf-check`) | 285 s |

Der Lauf von `make gates-schnell-gegenprobe` vom ersten Stand (Exit 0, 63 s) ist nicht wiederholt: Seit ihm ändert sich weder ein Fragment noch die Liste `GEGENPROBEN_CONTAINER`. Während der Arbeit nur die engsten Sensoren: `bash -n`, Teilläufe der Gegenprobe mit `GP_MUT` und `GP_ENTDECKE` (die roten Fälle je Mutation gemessen, dann die Erwartung im Skript gesetzt), `make docs-check`.

**Messung der Größe** (`AGENTS.md` §3.14): `git diff --numstat 11a2ba91..39d964e9 -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'` ergibt als Summe aus Einfügungen und Löschungen **1425** (`release-gegenprobe.sh` 1037, `release.sh` 231, `releasing.md` 83, `release.yml` 38, `release.mk` 12, `SPEC-067` 22 aus den Entscheidungen des Architect, `harness/README.md` 1, `VERSION` 1); ohne den Vertragstext 1403. Schätzung des Architect 820 bis 1070; Meldeschwelle 1000 und Marke 1300 sind überschritten. 2026-10-19 (Implementer): Der Bericht an die Koordination nennt die Zahl; die Entscheidung „weiter“ steht in der Closure unten.

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
| 2/3 (1): Ereignis `push` | die Prüfung des Ereignisses entfällt (`m-p2-ereignis-nicht-geprueft`) | `p2-aussenwirkung-ohne-push`, `p3-ereignis-fehlt`, `p3-ereignis-unbekannt`, `p3-schluesselwort-ereignis` |
| 2/3 (1): Ref `refs/tags/v…` | die Prüfung des Refs entfällt (`m-p2-ref-nicht-geprueft`) | `p2-ref-kein-tag`, `p3-ohne-v`, `p3-reihenfolge-1-2` |
| 2: Ref `refs/tags/v…` mit `v` | Muster `refs/tags/*` (`m-p2-ref-ohne-v-praefix`) | `p3-ohne-v`, `p3-reihenfolge-1-2` |
| 2/3 (1): `workflow_dispatch` ist der Trockenlauf | Trockenlauf wird wie ein Lauf geprüft (`m-p2-trockenlauf-wie-lauf`) | 9 Fälle, darunter `p2-trockenlauf-ohne-aussenwirkung`, `p3-trockenlauf-gueltig`, `p3-trockenlauf-im-standard-zweig` |
| 3: Prüfung (1) Ereignis und Ref | Aufruf entfernt (`m-p3-1-ereignis-entfernt`) | `p2-ref-kein-tag`, `p3-ereignis-fehlt`, `p3-ereignis-unbekannt`, `p3-ohne-v`, `p3-reihenfolge-1-2`, `p3-schluesselwort-ereignis` |
| 3: Prüfung (2) Form des Tags | Aufruf entfernt (`m-p3-2-tag-form-entfernt`) | 13 Fälle, darunter `p3-build-metadaten`, `p3-fuehrende-null`, `p3-fuehrende-null-minor` |
| 3: Prüfung (3) Form von `VERSION` | Aufruf entfernt (`m-p3-3-version-datei-entfernt`) | 56 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-reihenfolge-gruende` |
| 3: Prüfung (4) Tag gleich `v` plus `VERSION` | Aufruf entfernt (`m-p3-4-tag-gleich-entfernt`) | `p3-reihenfolge-4-5`, `p3-schluesselwort-tag-ungleich`, `p3-tag-ungleich-datei`, `p3-tag-ungleich-gross`, `p3-tag-ungleich-tag-kuerzer`, `p3-tag-ungleich-tag-laenger` |
| 3: Prüfung (5) Zeile der Änderungshistorie | Aufruf entfernt (`m-p3-5-notes-entfernt`) | 10 Fälle, darunter `p3-notes-abschnitt`, `p3-notes-abschnitt-ende`, `p3-notes-fehlen` |
| 3: Prüfung (6) Commit im Standard-Zweig | Aufruf entfernt (`m-p3-6-standard-zweig-entfernt`) | `p3-nicht-im-standard-zweig`, `p3-nicht-im-standard-zweig-vorwaerts`, `p3-schluesselwort-standard-zweig`, `p3-trockenlauf-im-standard-zweig`, `p3-trockenlauf-nicht-im-standard-zweig` |
| 3: Reihenfolge der Gründe (1) vor (2) | (1) und (2) vertauscht (`m-p3-reihenfolge-1-2-vertauscht`) | `p2-ref-kein-tag`, `p3-ohne-v`, `p3-reihenfolge-1-2` |
| 3: Reihenfolge der Gründe (2) vor (3) | (2) und (3) vertauscht (`m-p3-reihenfolge-2-3-vertauscht`) | 12 Fälle, darunter `p3-build-metadaten`, `p3-fuehrende-null`, `p3-fuehrende-null-minor` |
| 3: Reihenfolge der Gründe (3) vor (4) | (3) und (4) vertauscht (`m-p3-reihenfolge-3-4-vertauscht`) | 56 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-reihenfolge-gruende` |
| 3: Reihenfolge der Gründe (4) vor (5) | (4) und (5) vertauscht (`m-p3-reihenfolge-4-5-vertauscht`) | `p3-reihenfolge-4-5`, `p3-schluesselwort-tag-ungleich`, `p3-tag-ungleich-datei`, `p3-tag-ungleich-gross`, `p3-tag-ungleich-tag-kuerzer`, `p3-tag-ungleich-tag-laenger` |
| 3: Reihenfolge der Gründe (5) vor (6) | (5) und (6) vertauscht (`m-p3-reihenfolge-5-6-vertauscht`) | `p3-reihenfolge-5-6` |
| 3 (2): Tag hat die Form `v<SemVer>` | der Vergleich mit der Form entfällt (Länge bleibt) (`m-p3-form-entfernt`) | 12 Fälle, darunter `p3-build-metadaten`, `p3-fuehrende-null`, `p3-fuehrende-null-minor` |
| 3 (2): Länge bis 100 Zeichen | die Prüfung entfällt (`m-p3-laenge-entfernt`) | `p3-zu-lang` |
| 3 (2): 101 Zeichen sind zu lang | Grenze 101 (`m-p3-laenge-101`) | `p3-zu-lang` |
| 3 (2): 100 Zeichen sind zugelassen | Grenze 99 (`m-p3-laenge-99`) | `p3-laenge-100` |
| 3 (2): MAJOR ohne führende Null | `[0-9]+` für MAJOR (`m-p3-null-major`) | `p3-fuehrende-null`, `p3-reihenfolge-2-3` |
| 3 (2): MINOR ohne führende Null | `[0-9]+` für MINOR (`m-p3-null-minor`) | `p3-fuehrende-null-minor` |
| 3 (2): PATCH ohne führende Null | `[0-9]+` für PATCH (`m-p3-null-patch`) | `p3-fuehrende-null-patch` |
| 3 (2): `MAJOR.MINOR.PATCH`, drei Teile | PATCH optional (`m-p3-kern-zwei-teile`) | `p3-schluesselwort-tag-form`, `p3-zwei-teile` |
| 3 (2): genau drei Teile | vierter Teil zugelassen (`m-p3-kern-vier-teile`) | `p3-vier-teile` |
| 3 (2): Vorabversion nach `-` nicht leer | `v1.2.3-` zugelassen (`m-p3-vor-leer-erlaubt`) | `p3-leere-vorabversion` |
| 3 (2): Bezeichner der Vorabversion nicht leer | leerer Bezeichner zwischen Punkten zugelassen (`m-p3-vor-leerer-bezeichner`) | `p3-leerer-bezeichner` |
| 3 (2): keine Build-Metadaten | `+…` zugelassen (`m-p3-vor-build-metadaten`) | `p3-build-metadaten` |
| 3 (2): Bezeichner durch `.` getrennt | nur ein Bezeichner zugelassen (`m-p3-vor-ein-bezeichner`) | `p3-tag-gueltig`, `p3-tag-ungleich-gross`, `p3-tag-ungleich-tag-kuerzer`, `p3-tag-ungleich-tag-laenger`, `p9-vorabversion-prerelease` |
| 3 (2): numerischer Bezeichner ohne führende Null | jede Zeichenfolge zugelassen (`m-p3-bez-fuehrende-null`) | `p3-vorabversion-fuehrende-null` |
| 3 (2): Zeichen `0-9A-Za-z-` | `_` zugelassen (`m-p3-bez-zeichen`) | `p3-ungueltiges-zeichen` |
| 3 (2): Bezeichner mit Buchstaben | nur Zahlen zugelassen (`m-p3-bez-nur-zahlen`) | 10 Fälle, darunter `p2-trockenlauf-ohne-aussenwirkung`, `p3-laenge-100`, `p3-tag-gueltig` |
| 3 (2): Ziffern in alphanumerischen Bezeichnern | Ziffern nicht zugelassen (`m-p3-bez-ohne-ziffern`) | `p3-tag-gueltig`, `p3-vorabversion-fuehrende-null` |
| 3 (2): Großbuchstaben | nicht zugelassen (`m-p3-bez-ohne-grossbuchstaben`) | `p3-tag-gueltig`, `p3-tag-ungleich-gross` |
| 3 (2): Kleinbuchstaben | nicht zugelassen (`m-p3-bez-ohne-kleinbuchstaben`) | 10 Fälle, darunter `p2-trockenlauf-ohne-aussenwirkung`, `p3-laenge-100`, `p3-tag-gueltig` |
| 3 (2): Bindestrich im Bezeichner | nicht zugelassen (`m-p3-bez-ohne-bindestrich`) | `p3-tag-gueltig` |
| 3 (2): Präfix `v` | `v` optional (äquivalent, siehe unten) (`m-p3-form-v-optional`) | keiner (äquivalent) |
| 3 (3): Form von `VERSION` | der Vergleich mit der Form entfällt (`m-p3-vform-entfernt`) | 10 Fälle, darunter `p3-reihenfolge-3-4`, `p3-schluesselwort-version-datei`, `p3-trockenlauf-version-datei-form` |
| 3 (3): LF am Ende | Zeilenende optional (`m-p3-vdatei-kein-lf`) | `p3-schluesselwort-version-datei`, `p3-trockenlauf-version-datei-form`, `p3-version-datei-form-kein-lf` |
| 3 (3): LF, nicht CRLF | CR vor LF zugelassen (`m-p3-vdatei-crlf`) | `p3-reihenfolge-3-4`, `p3-version-datei-form-crlf` |
| 3 (3): keine Leerzeichen | Leerzeichen vorn und hinten zugelassen (`m-p3-vdatei-leerzeichen`) | `p3-version-datei-form-leerzeichen-hinten`, `p3-version-datei-form-leerzeichen-vorn` |
| 3 (3): `<SemVer>` ohne `v` | `v` zugelassen (`m-p3-vdatei-v`) | `p3-version-datei-form-v` |
| 3 (3): genau eine Zeile | mehrere Zeilen zugelassen (`m-p3-vdatei-mehrere-zeilen`) | `p3-version-datei-form-zwei-zeilen` |
| 3 (4): Vergleich `v` plus `VERSION` | Vergleich ohne `v` (`m-p3-gleich-ohne-v`) | 46 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-reihenfolge-gruende` |
| 3 (4): Zeichen für Zeichen | Groß- und Kleinschrift gleich (`m-p3-gleich-gross`) | `p3-tag-ungleich-gross` |
| 3 (4): Zeichen für Zeichen | Tag darf `v<VERSION>` fortsetzen (`m-p3-gleich-praefix-tag-laenger`) | `p3-tag-ungleich-tag-laenger` |
| 3 (4): Zeichen für Zeichen | Tag darf Anfang von `v<VERSION>` sein (`m-p3-gleich-praefix-tag-kuerzer`) | `p3-tag-ungleich-tag-kuerzer` |
| 3 (5): Tabelle im Abschnitt `Änderungshistorie` | jede Tabelle der Datei zählt (`m-p3-notes-abschnitt-ignoriert`) | `p3-notes-abschnitt`, `p3-notes-abschnitt-ende`, `p3-notes-ueberschrift-endet`, `p9-notes-aus-abschnitt` |
| 3 (5): Überschrift endet auf `Änderungshistorie` | Überschrift enthält das Wort (`m-p3-notes-ueberschrift-ende-ignoriert`) | `p3-notes-ueberschrift-endet` |
| 3 (5): nur die Tabelle dieses Abschnitts | Abschnitt endet nicht an der nächsten Überschrift (`m-p3-notes-abschnitt-ende-ignoriert`) | `p3-notes-abschnitt-ende` |
| 3 (5)/9: Zeile mit erster Zelle gleich der Version | Anfang der ersten Zelle genügt (`m-p3-notes-praefix`) | `p9-notes-zeile-der-version` |
| 9: Text ist die dritte Zelle | zweite Zelle (`m-p3-notes-zelle-datum`) | `p3-notes-zelle-leer`, `p3-trockenlauf-notes-vorhanden`, `p3-trockenlauf-notes-zelle-leer`, `p9-notes-aus-abschnitt`, `p9-notes-aus-zelle`, `p9-notes-zeile-der-version` |
| 3 (6): Commit liegt im Standard-Zweig | `merge-base` ohne `--is-ancestor` (`m-p3-anc-ohne-is-ancestor`) | `p3-nicht-im-standard-zweig`, `p3-nicht-im-standard-zweig-vorwaerts`, `p3-schluesselwort-standard-zweig`, `p3-trockenlauf-nicht-im-standard-zweig` |
| 3 (6): Commit liegt im Standard-Zweig | Argumente vertauscht (`m-p3-anc-umgekehrt`) | `p3-nicht-im-standard-zweig`, `p3-nicht-im-standard-zweig-vorwaerts`, `p3-schluesselwort-standard-zweig`, `p3-trockenlauf-nicht-im-standard-zweig` |
| 3 (6): Standard-Zweig aus `RELEASE_DEFAULT_BRANCH` | fest `origin/main` (`m-p3-standard-zweig-main`) | `p3-standard-zweig-aus-umgebung` |
| 3 (6): beim Lauf Eingabefehler | wie im Trockenlauf nur ein Hinweis (`m-p3-lauf-nur-hinweis`) | `p3-nicht-im-standard-zweig`, `p3-nicht-im-standard-zweig-vorwaerts`, `p3-schluesselwort-standard-zweig` |
| 3 (6): im Trockenlauf kein Fehler | Trockenlauf wird wie der Lauf behandelt (`m-p3-trockenlauf-fehler`) | `p3-trockenlauf-im-standard-zweig`, `p3-trockenlauf-nicht-im-standard-zweig` |
| 3 (6): der Hinweis nennt die Antwort | immer „liegt im“ (`m-p3-hinweis-immer-ja`) | `p3-trockenlauf-nicht-im-standard-zweig` |
| 3 (6): im Trockenlauf steht die Antwort als Hinweis | kein Hinweis (`m-p3-hinweis-fehlt`) | `p3-trockenlauf-im-standard-zweig`, `p3-trockenlauf-nicht-im-standard-zweig` |
| 2: Tag des Trockenlaufs ist `v<VERSION>` | Tag aus dem Ref (`m-p3-trockenlauf-tag-aus-ref`) | 9 Fälle, darunter `p2-trockenlauf-ohne-aussenwirkung`, `p3-trockenlauf-gueltig`, `p3-trockenlauf-im-standard-zweig` |
| 4: Aufruf mit `VERSION=<Version>` | `VERSION=` fehlt im Aufruf (`m-p4-version-fehlt`) | 20 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
| 4: `DIST` außerhalb des Arbeitsbaums | `DIST=dist` (`m-p4-dist-im-baum`) | 20 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
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
| 9/11: vorhandener Release wird ergänzt, nicht neu angelegt | `create` auch bei vorhandenem Release (`m-p9-create-trotz-vorhandenem-release`) | `p10-nicht-im-log`, `p11-wiederanlauf-release-teilweise`, `p11-wiederanlauf-release-vollstaendig`, `p9-release-existiert-abweichend`, `p9-release-existiert-gleich`, `p9-repo-argument` |
| 9/11: nur fehlende Anhänge werden hochgeladen | alle Anhänge (`m-p9-upload-alle`) | `p11-wiederanlauf-release-teilweise`, `p11-wiederanlauf-release-vollstaendig`, `p9-release-existiert-gleich` |
| 9: Fehler von `gh release view` ist rot | jeder Fehler gilt als „nicht gefunden“ (`m-p9-view-anderer-fehler-wie-fehlend`) | `p10-nicht-im-log`, `p9-view-anderer-fehler` |
| 3 (5): leere dritte Zelle ist im Lauf ein Fehler | im Lauf angenommen (`m-p3-zelle-leer-angenommen-lauf`) | `p3-notes-zelle-leer` |
| 3 (5): leere dritte Zelle ist im Trockenlauf ein Fehler | im Trockenlauf angenommen (`m-p3-zelle-leer-angenommen-trockenlauf`) | `p3-trockenlauf-notes-zelle-leer` |
| 3 (5): leere dritte Zelle gilt wie eine fehlende Zeile | die Prüfung entfällt (`m-p3-zelle-leer-angenommen`) | `p3-notes-zelle-leer`, `p3-trockenlauf-notes-zelle-leer` |
| 3 (5): der Trockenlauf zeigt den Text der vorhandenen Zeile als Hinweis | der Hinweis entfällt (`m-p3-trockenlauf-zeile-nicht-ausgewertet`) | `p3-trockenlauf-notes-vorhanden` |
| 3 (5): im Trockenlauf ist eine fehlende Zeile ein Hinweis | sie ist ein Fehler (`m-p3-trockenlauf-fehlende-zeile-fehler`) | `p3-trockenlauf-notes-fehlen` |
| 3 (5): im Lauf ist eine fehlende Zeile ein Fehler `notes` | die eigene Prüfung entfällt (äquivalent, siehe unten) (`m-p3-lauf-fehlende-zeile-ohne-eigene-pruefung`) | keiner (äquivalent) |
| 3: Meldung beginnt mit `release: ereignis: ` | das Schlüsselwort weicht ab (`m-p3-schluesselwort-ereignis-weicht-ab`) | 7 Fälle, darunter `p2-aussenwirkung-ohne-push`, `p2-ref-kein-tag`, `p3-ereignis-fehlt` |
| 3: Meldung beginnt mit `release: tag-form: ` | das Schlüsselwort weicht ab (`m-p3-schluesselwort-tag-form-weicht-ab`) | 13 Fälle, darunter `p3-build-metadaten`, `p3-fuehrende-null`, `p3-fuehrende-null-minor` |
| 3: Meldung beginnt mit `release: version-datei: ` | das Schlüsselwort weicht ab (`m-p3-schluesselwort-version-datei-weicht-ab`) | 10 Fälle, darunter `p3-reihenfolge-3-4`, `p3-schluesselwort-version-datei`, `p3-trockenlauf-version-datei-form` |
| 3: Meldung beginnt mit `release: tag-ungleich: ` | das Schlüsselwort weicht ab (`m-p3-schluesselwort-tag-ungleich-weicht-ab`) | `p3-reihenfolge-4-5`, `p3-schluesselwort-tag-ungleich`, `p3-tag-ungleich-datei`, `p3-tag-ungleich-gross`, `p3-tag-ungleich-tag-kuerzer`, `p3-tag-ungleich-tag-laenger` |
| 3: Meldung beginnt mit `release: notes: ` | das Schlüsselwort weicht ab (`m-p3-schluesselwort-notes-weicht-ab`) | 8 Fälle, darunter `p3-notes-abschnitt`, `p3-notes-abschnitt-ende`, `p3-notes-fehlen` |
| 3: Meldung beginnt mit `release: standard-zweig: ` | das Schlüsselwort weicht ab (`m-p3-schluesselwort-standard-zweig-weicht-ab`) | `p3-nicht-im-standard-zweig`, `p3-nicht-im-standard-zweig-vorwaerts`, `p3-schluesselwort-standard-zweig` |
| 9: jeder Aufruf von `gh` mit `--repo` | `--repo` fehlt bei `release view` (`m-p9-repo-fehlt-view`) | `p9-repo-argument` |
| 9: jeder Aufruf von `gh` mit `--repo` | `--repo` fehlt bei `release create` (`m-p9-repo-fehlt-create`) | `p9-repo-argument` |
| 9: jeder Aufruf von `gh` mit `--repo` | `--repo` fehlt bei `release download` (`m-p9-repo-fehlt-download`) | `p9-repo-argument` |
| 9: jeder Aufruf von `gh` mit `--repo` | `--repo` fehlt bei `release upload` (`m-p9-repo-fehlt-upload`) | `p9-repo-argument` |
| 9: ohne `GH_TOKEN` Exit 2 vor dem ersten Aufruf von `gh` | die Prüfung entfällt (`m-p9-gh-ohne-token-pruefung`) | `p9-gh-token-fehlt` |
| 9: nur `release not found` heißt „nicht vorhanden“ (F-1) | Muster geweitet auf `not found` ohne Schreibweise (`m-p9-view-muster-geweitet`) | `p9-view-anderer-fehler` |
| 9: Anhänge genau die sieben Dateien (F-4) | `darwin/amd64` fehlt in `erwartete_namen` (`m-p9-name-fehlt-darwin_amd64`) | 20 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
| 9: Anhänge genau die sieben Dateien (F-4) | `darwin/arm64` fehlt in `erwartete_namen` (`m-p9-name-fehlt-darwin_arm64`) | 20 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
| 9: Anhänge genau die sieben Dateien (F-4) | `linux/amd64` fehlt in `erwartete_namen` (`m-p9-name-fehlt-linux_amd64`) | 20 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
| 9: Anhänge genau die sieben Dateien (F-4) | `linux/arm64` fehlt in `erwartete_namen` (`m-p9-name-fehlt-linux_arm64`) | 21 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
| 9: Anhänge genau die sieben Dateien (F-4) | `windows/amd64` fehlt in `erwartete_namen` (`m-p9-name-fehlt-windows_amd64`) | 20 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
| 9: Anhänge genau die sieben Dateien (F-4) | `windows/arm64` fehlt in `erwartete_namen` (`m-p9-name-fehlt-windows_arm64`) | 20 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
| 9: Anhänge genau die sieben Dateien (F-4) | `SHA256SUMS` fehlt in `erwartete_namen` (`m-p9-name-fehlt-sha256sums`) | 20 Fälle, darunter `p10-nicht-im-log`, `p11-alte-marke`, `p11-wiederanlauf-release-teilweise` |
| 3 (5)/9: Zeile mit genau drei Zellen, ein `|` im Text ist ein Fehler `notes` (F-5) | die Zeile mit weiterem `|` wird angenommen (`m-p3-pipe-abgeschnitten-angenommen`) | `p3-notes-pipe` |
| 10: kein Wert des Tokens in der Ausgabe (Fehler im Vergleich) | das Skript gibt `$GH_TOKEN` aus (`m-p10-token-bei-fehler`) | `p10-nicht-im-log` |
| 10: kein Wert des Tokens in der Ausgabe (Fehler der Summe) | das Skript gibt `$GH_TOKEN` aus (`m-p10-token-bei-fehler-binaries`) | `p10-nicht-im-log` |
| 10: kein `set -x` | `set -x` im Skript (`m-p10-set-x`) | 7 Fälle, darunter `p10-nicht-im-log`, `p3-schluesselwort-ereignis`, `p3-schluesselwort-notes` |
| 10: kein `env` | `env` im Fehlerfall (`m-p10-env-bei-fehler`) | `p10-nicht-im-log` |
| 11: `binaries` läuft nicht nach rotem `eingaben` | die Marke wird nicht geprüft (`m-p11-binaries-ohne-eingaben-marke`) | `p11-alte-marke`, `p11-reihenfolge-gruende`, `p13-binaries-ohne-eingaben` |
| 11: `github-release` läuft nicht nach rotem `binaries` | die Marke wird nicht geprüft (`m-p11-release-ohne-binaries-marke`) | `p11-alte-marke`, `p11-reihenfolge-gruende`, `p13-release-ohne-binaries` |
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

**Findings der Review-Runde 1** (Report `2026-10-19-review-slice-erster-release-veroeffentlichung`; Behebung als Klasse, `implement-slice.md` Schritt 19):

| Finding | Behebung | Beleg |
|---|---|---|
| F-1 (MEDIUM): Muster „nicht vorhanden“ nur an `release not found` | Fall `p9-view-anderer-fehler` hat ein zweites Bein: `gh` meldet `HTTP 404: Not Found`, Exit 1, kein `create`; Mutation `m-p9-view-muster-geweitet` (`grep -qi 'not found'`) | rot: `p9-view-anderer-fehler` |
| F-2 (MEDIUM): gleichzeitige Läufe | Entscheidung des Architect `eec9c839`; `release.yml` mit `group: release-${{ github.event_name }}` (Zeile des Textes in `p1-*`, Zeile 9, mit den Mutationen aller Zeilen), Plan §1 auf den **laufenden** Lauf eingegrenzt | `p1-gleich` grün, `p1-fehlt-9`/`p1-geaendert-9` rot unter den Zeilenmutationen |
| F-3 (LOW): Probe-Tag in `releasing.md` | §3 sagt, dass die Begriffe nur das Verfahren der Formel (§7) betreffen und die Veröffentlichung der Binaries (§6) den Trockenlauf statt eines Probe-Tags hat | `make docs-check` 0 Befunde |
| F-4 (LOW): Aufzählung der Anhänge ohne Mutation je Element | sieben Mutationen `m-p9-name-fehlt-<Element>` (sechs Binaries, `SHA256SUMS`); jede färbt die Fälle, die `binaries` grün fahren, `linux_arm64` zusätzlich `p4-datei-fehlt`; auf einzelne Fälle lässt sich das nicht eingrenzen, weil `erwartete_namen` Bau und Release gemeinsam benutzen | die Mengen in der Tabelle |
| F-5 (LOW): Zelle mit `|` | Entscheidung des Architect `eec9c839` (genau drei Zellen); `notiz_der_version` meldet `X` bei einer Zeile ohne genau drei Zellen, `pruefe_notes` und `github-release` lehnen ab (`notes`, Exit 2); Fall `p3-notes-pipe` (roh und als `\|`, Lauf und Trockenlauf); Mutation `m-p3-pipe-abgeschnitten-angenommen` | rot: `p3-notes-pipe` |
| F-6, F-7 (INFO) | Größe: oben, Entscheidung steht aus; Trockenlauf auf GitHub und Wortlaut von `gh`: unverändert offen | — |

**Zwei grüne Mutanten, äquivalent** (`AGENTS.md` §3.10):

1. `m-p3-form-v-optional` (die Form des Tags verlangt das `v` nicht mehr) färbt keinen Fall. (a) Invariante: Wenn `pruefe_tag_form` läuft, beginnt der Tag mit `v`. (b) Ort: beim Lauf `aussenwirkung_erlaubt` über `pruefe_ereignis` (Ref `refs/tags/v*`, Exit 2 sonst), beim Trockenlauf `tag_ableiten` (`v` plus erste Zeile von `VERSION`); beide sind Zusage dieses Slice und haben eigene Mutationen, die Fälle rot färben (`m-p2-ref-ohne-v-praefix`, `m-p3-trockenlauf-tag-aus-ref`). (c) Sonde: Eingaben `push` mit Ref `refs/tags/1.2.3`, `VERSION` `1.2.3`, Zeile der Version im Handbuch, gegen das Original und den Mutanten: beide Exit 2 mit dem Schlüsselwort `ereignis`. Gegen-Sonde: Ohne die Prüfung des Refs (Zeile `# P2ref` entfernt) endet das Original mit `tag-form` und der Mutant mit `tag-ungleich`, beide Exit 2; die Sonde trennt also Original und Mutant am Schlüsselwort, sobald die Invariante nicht mehr gilt.
2. `m-p3-lauf-fehlende-zeile-ohne-eigene-pruefung` (im Lauf entfällt der Fehler `notes` für eine fehlende Zeile) färbt keinen Fall. (a) Invariante: Im Lauf gilt eine leere dritte Zelle wie eine fehlende Zeile (SPEC-067 Punkt 3 (5)), beides ist Exit 2 mit dem Schlüsselwort `notes`. (b) Ort: die Prüfung der Zelle `# P3zelle` in `pruefe_notes`, die für `z` leer (keine Zeile) ebenso greift; sie ist Zusage dieses Slice (`m-p3-zelle-leer-angenommen-lauf` färbt `p3-notes-zelle-leer`). (c) Sonde: Handbuch ohne Zeile der Version, Lauf: Original `release: notes: die Änderungshistorie … keine Zeile für 1.2.3`, Exit 2; Mutant `release: notes: die Zeile für 1.2.3 hat eine leere dritte Zelle`, Exit 2 (anderer Text, gleiches Schlüsselwort und Exit-Code). Gegen-Sonde: Ohne beide Prüfungen (`# P3notes` und `# P3zelle` entfernt) endet der Lauf mit Exit 0, die Sonde trennt also, sobald die Invariante nicht mehr gilt. Aufgerufen am 2026-10-19 mit `bash <Kopie von release.sh> eingaben` im Dev-Container, Umgebung `RELEASE_EVENT=push`, `RELEASE_REF=refs/tags/v1.2.3`.

**Trockenlauf auf GitHub:** vom Implementer nicht gefahren (sein Auftrag verbietet einen Push auf GitHub); gefahren am 2026-10-19, Beleg in der Closure unten. Gemessen vom Implementer auf dem Stand der Lieferung mit den Eingaben eines Trockenlaufs (`RELEASE_EVENT=workflow_dispatch`, `VERSION` `0.1.0-dev`): `bash tools/harness/release.sh eingaben` endet mit Exit 0 und den Hinweisen `release: hinweis: die Änderungshistorie des Benutzerhandbuchs hat keine Zeile für 0.1.0-dev` und (auf dem Zweig des Klons) `release: hinweis: der Commit liegt nicht im Standard-Zweig origin/main`.

**Handbuch und README:** keine Zeile neu oder geändert. Gelesen gegen `make binaries` und `make test-container` (im grünen `make gates`): Voraussetzungen (sechs Plattformen, nur `linux/amd64` ausgeführt), Installation (`make binaries`, Namen mit `VERSION`, `SHA256SUMS`, `dist/` leer oder fehlend) und der Absatz von `README.md` stehen schon im Ist-Zustand; der Abschnitt zum Bezug vom GitHub-Release gehört `slice-erster-release-freigabe` (§1).

### Closure (Planner, 2026-10-19)

Regeln: Baseline-Regelwerk `modul-05-planning-harness.md` §Closure- und Lerneintrag-Regeln und `modul-06-roadmap.md`. Gelesen: Review-Report (`docs/reviews/2026-10-19-review-slice-erster-release-veroeffentlichung.md`), Verifikationsbericht (`docs/reviews/2026-10-19-verifikation-slice-erster-release-veroeffentlichung.md`), die Läufe auf GitHub unten (Angaben der Koordination, vom Planner nicht selbst nachgefragt). Der Planner fuhr keine Gegenprobe und kein `make gates`; seit `39d964e9` ändern sich im Slice nur Pläne, Reports und `docs/developer/` (Beleg des Implementers gilt), und die CI auf `main` fuhr `make gates` auf dem Stand unten grün.

**Stand.** Code des Slice `39d964e9`; die Läufe unten liegen auf `main` bei `e3ad97e6` (Merge der Fassung A mit den Entscheidungen des Architect und den Reports).

**Läufe auf GitHub** (Closure-Trigger §5; Kontrolle `gh run view <id>`):

| Lauf | Anlass · Stand | Ergebnis | Dauer | Schritte |
|---|---|---|---|---|
| `39208276024` | `release.yml`, `workflow_dispatch` auf `main` · `e3ad97e6` | success | 146 s (2026-10-19 06:09:18Z bis 06:11:44Z) | `Eingaben` (Hinweise: keine Zeile für `0.1.0-dev` in der Änderungshistorie des Handbuchs; der Commit liegt im Standard-Zweig `origin/main`), `Binaries` (sieben Dateien), `GitHub-Release` übersprungen (nur bei `push`) |
| CI auf `main` | `ci.yml` · `e3ad97e6` | success | nicht erhoben | `make gates` auf dem Runner |

Der Trockenlauf belegt: die Rechte der Datei (`permissions: {}` oben, `contents: write` nur im Job) reichen für `Eingaben` und `Binaries`, die Dauer liegt bei rund zweieinhalb Minuten, der Trockenlauf legt keinen GitHub-Release an. Die Hinweise sind die zugesagten (`SPEC-067` Punkt 3 (5) und (6), Punkt 12: Zeile fehlt, Trockenlauf nur Hinweis).

**Nicht auf GitHub belegt** (an `slice-erster-release-freigabe` abgegeben, §6): der Wiederanstoß eines roten Laufs (V-1 b), der Wortlaut `release not found` von `gh release view` und das Anlegen des GitHub-Release.

**Metrik gegen die Marke von §4:** gemessen 1425 (§7 *Messung der Größe*, vom Verifier nachgemessen: 1425 mit beiden Basen), Marke 1300 und Meldeschwelle 1000 überschritten, Schätzung des Architect 820 bis 1070. 2026-10-19 (Planner, auf Anweisung der Koordination): **weiter, kein Schnitt.** Gründe: Die Gegenprobe trägt 1037 von 1425 Zeilen (73 %), `release.sh`, `release.yml` und `release.mk` zusammen 281; ein Schnitt nach der Lieferung bringt nichts, weil die Gegenprobe nicht in einen Teil ohne Code zerfällt. Die Review-Nacharbeit (F-1 bis F-5) und die Entscheidungen des Architect machten rund 120 Zeilen aus. Review und Verifikation haben den Diff geprüft, ohne dass die Prüfbarkeit in einer Sitzung zum Befund wurde (Rückführung `in-progress` → `next` nicht nötig). Messwert für `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`: Faktor Gegenprobe zu Produktcode hier rund 3,7 (1037 zu 281), nicht 2.

**Findings und Ausgänge.** Review: F-1 (MEDIUM, Muster `release not found` ohne Grenzmutation), F-3, F-4, F-5 (LOW): behoben (Tabelle der Findings oben); F-2 (MEDIUM, Gleichzeitigkeit): vom Architect entschieden (`eec9c839`, Gruppe `release-${{ github.event_name }}`, Grenze in `SPEC-067` Punkt 13); F-6 (Größe): oben; F-7 (INFO): Trockenlauf oben, Wortlaut von `gh` abgegeben. Verifikation: V-1 (Belege auf GitHub): der Trockenlauf oben, der Rest abgegeben; V-2 (Größe): oben; V-3 (Behebungen vom Reviewer nicht erneut gesehen, nur Stichprobe des Verifiers): kein zweiter Review-Durchgang angeordnet (Entscheidung der Koordination), als Auftreten in `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` eingetragen; V-4 (`dist` im zweiten Aufruf von Hand) und V-5, V-6 (INFO): keine Arbeit, V-4 in der Pipeline nicht erreichbar (jeder Re-run läuft auf einem neuen Runner) und nicht zugesagt.

**Lerneintrag — Form: benannte Lücken und Vorschläge** (eine geschärfte Regel trägt dieser Slice nicht ein, ein Herkunfts-Anker `liegt in` entfällt, weil nichts verkörpert wurde; Regeländerungen in `AGENTS.md` und `.claude/` brauchen das Ja des Nutzers und stehen hier als Vorschlag).

- *Lücke 1 (benannt, Planner):* Die DoD von Liefer-Punkt 3 verlangte einen Beleg („ein roter Lauf ist nach der Tabelle wieder angestoßen“), den das gewählte Belegverfahren nicht erzeugen kann: Ein Trockenlauf ist grün gebaut. Der Plan nannte den Beleg und den Trockenlauf im selben Punkt, ohne beide zu prüfen; erst die Verifikation stellte es fest (V-1). Register: `BEO-REPO/dod-beleg-vom-belegverfahren-nicht-erzeugbar` (neu, 1×). Vorschlag an den Nutzer, nicht eingetragen: Beim Schnitt prüft der Planner zu jedem DoD-Beleg auf einer fremden Plattform, welcher Lauf ihn erzeugt; erzeugt das gewählte Verfahren ihn nicht, geht er als Sendung (`AGENTS.md` §3.13) an den Slice, dessen Lauf es tut (`.claude/commands/plan-welle.md`, Schritt zur DoD).
- *Lücke 2 (Vorschlag, nicht eingetragen):* Die Regel „Aufzählung und Ortsangabe zählen je Element“ (`AGENTS.md` §3.10) stand im Repo, und F-4 traf sie am Skript der Pipeline dennoch (sechs Binaries und `SHA256SUMS` ohne Mutation je Element); F-1 zeigt eine Form, die die Regel nicht nennt: die **Eingrenzung** („nur `release not found` heißt nicht vorhanden“) braucht eine Mutation, die das Muster weitet, nicht nur eine, die die Prüfung entfernt. Vorschlag: ein Satz in `AGENTS.md` §3.10, dass eine Zusage der Form „nur X“ eine Grenzmutation (Muster geweitet) trägt. Register: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (Zähler 59).
- *Lücke 3 (benannt, beim Architect erledigt):* Eine Zusage über das Verhalten der Plattform („ein Tag wird nie durch einen späteren abgebrochen“, `concurrency`) war weiter gefasst als die dokumentierte Semantik (F-2); der Architect hat sie auf den laufenden Lauf eingegrenzt und die Grenze festgehalten. Ob das Verhalten stimmt, zeigt nur ein Lauf auf GitHub. Register: `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (86).
- *Was trug:* Der Trockenlauf statt eines Probe-Tags belegte die Rechte und die Dauer in 146 s ohne Außenwirkung; die Randform-Rückgabe (`AGENTS.md` §3.12) lief zweimal an den Architect, bevor Code entschieden wurde, und die Entscheidungen (Marken, Zeile der Änderungshistorie, `gh`) standen vor dem Review. Die Meldeschwelle neben der Marke (1000 neben 1300) wurde gemeldet und die Entscheidung als datierte Zeile festgehalten (`BEO-REPO/meldeschwelle-ohne-artefakt`, Gegenbeleg, Zähler bleibt).
- *Was ging anders als geplant:* Die Größe (1425 gegen 820 bis 1070): die Zeilenliste von `release.yml` (38 Zeilen, je Zeile ein Fall „fehlt“ und ein Fall „geändert“ mit je einer Mutation) und die Mutation je Element der SemVer-Form und der Anhänge machten die Gegenprobe zum Hauptteil; der Faktor 2 des Architect (`AGENTS.md` §3.14) gilt für diese Form nicht. Vorschlag für `slice-harness-schaetzung-planner-faktor` (Messwert, keine Sendung eingetragen): für eine Gegenprobe, die den Text einer Datei zeilenweise hält, mit dem Dreifachen bis Vierfachen des Produktcodes rechnen.

**Register (Paarung 3).** Neu: `BEO-REPO/dod-beleg-vom-belegverfahren-nicht-erzeugbar` (1×, offen). Evidence: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (Zähler 59, F-1, F-4), `BEO-REPO/spec-randform-erst-im-review-entschieden` (43, F-5), `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (86, F-2), `BEO-REPO/plan-folgt-korrektur-nicht` (71, F-3: `releasing.md` §3 nannte den Probe-Tag nach der Entscheidung vom 2026-10-19 weiter als Verfahren), `BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests` (18), `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` (3, V-3; **Schwelle erreicht**: Zuweisung `verkörpert` oder `geplant` entscheidet der Nutzer, der Planner weist nichts zu). Vermerke ohne Auftreten: `BEO-REPO/meldeschwelle-ohne-artefakt` (Gegenbeleg), `BEO-REPO/groesse-erst-in-der-pruefung-des-architect-sichtbar` (die Größe war vor dem Anspruch sichtbar, der Schnitt vom 2026-10-19 lief davor), `BEO-REPO/gruener-mutant-faelschlich-aequivalent` (die zwei äquivalenten Mutanten trugen, das Review fuhr die Sonde nach). Die Einträge mit Evidence liegen bis auf den letzten und den neuen schon über der Schwelle (verkörpert bzw. `geplant`); der Lese-Schritt gehört der Welle-Closure.

**Paarungen.** Anker: keiner (keine verkörperte Regel). Folge-Slice: `slice-erster-release-freigabe` nimmt die Sendung an (dort §1 *Sendung aus `slice-erster-release-veroeffentlichung` — Belege, die der Trockenlauf nicht erzeugt*, DoD 3, §6, §8 mit Nachzählung; `grep` nach der Kennung dieses Slice findet sie); es liegt in `next/` und ist kein geschlossener Slice. Register: oben. Die Wellen-Zugehörigkeit bleibt (welle-erster-release, Schritt 4); eine Umplanung fällt nicht an, das Drift-Log trägt keine Zeile.

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

**Nachgezählt am 2026-10-19** (Planner, Closure): Die Abgabe des Rests der Belege an `slice-erster-release-freigabe` nimmt aus diesem Slice, sie fügt nichts hinzu; Liefer-Punkte drei, Schichten zwei, wie geliefert. Register am Stand der Closure nachgesichtet: `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` trifft die Abgabe; der Nehmer nimmt an (Zählung in dessen §8).
