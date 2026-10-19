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
| `tools/harness/release.sh` (Schritte `eingaben`, `binaries`, `github-release`) | neu | die Logik steht im Skript, nicht im YAML (Schnittstelle des Vertrags); B ergänzt die Schritte `image`, `veroeffentlichen`, `vergleichen`, `smoke`, `latest` |
| `VERSION` (Wurzel) | neu | Version der Software, eine Zeile `0.1.0-dev`, gegen den Tag geprüft; `make binaries VERSION=` trägt sie ins Binary |
| `tools/harness/release-gegenprobe.sh`, `harness/mk/release.mk`, Eintrag in `GEGENPROBEN_CONTAINER` | neu | Gegenprobe mit Attrappen für `make`, `sha256sum`, `gh` (und `docker`, `skopeo` für den Fall, dass der Trockenlauf sie nie ruft), ohne Netz: Fälle `p1-*`, `p2-*`, `p3-*` ohne (7), `p4-*`, `p9-*` ohne `p9-nicht-vor-registries`, `p10-*` ohne die Registries, `p11-*` für A, `p13-*`; je Zusage und je Element eine Mutation (`AGENTS.md` §3.10); die Zeilenliste von `release.yml` trägt die Herkunft A je Zeile |
| `harness/README.md` §Sensors | update | Zeile für `make release-gegenprobe` (Gegenprobe), Bindung `SPEC-067` und [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md) |
| `docs/maintainer/releasing.md` §3, §6, §9, §11 | update | vom beschriebenen auf das ausgeführte Verfahren bringen, soweit A es liefert: Versionsquelle (`VERSION`), was der Tag auslöst, Release-Notes, Fehlertabelle (Zeilen `eingaben`, `binaries`, `github-release`) und Wiederanlauf; auf dem Stand aufbauend, den `slice-v1-abschluss-container-doku` liefert. §8 (Kontrollen nach dem Push der Registries), §2 und §7 (Formel keine Lieferung des Release) ergänzt B beziehungsweise gehören nicht hierher |
| `docs/user/benutzerhandbuch.md`, `README.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung); nur Plattformen und `SHA256SUMS` der gebauten Binaries |

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
