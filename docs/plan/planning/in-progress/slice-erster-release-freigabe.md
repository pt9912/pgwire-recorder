# Slice slice-erster-release-freigabe: Freigabe und erster echter Tag

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-erster-release.

**Bezug:** [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)

**Berührte Spec-Stellen:** `LH-FA-16.a` · `SPEC-031` · `SPEC-067` (Punkte 1 und 5: der Scan-Schritt zwischen `Image` und `Veröffentlichen` ändert den Text von `release.yml`; gelesen: Punkte 6, 8, 11, 12 zum ersten Tag; Sendungen aus `slice-erster-release-veroeffentlichung` und `slice-erster-release-veroeffentlichung-image`, §1)

**Verantwortlich:** pt9912 (Zweig impl-fr)

**Autor:** pt9912. **Datum:** 2026-10-09.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Der erste echte Tag ist nach einer Freigabe-Checkliste mit Beleg je Punkt freigegeben, sein Image ist vor dem Push gescannt, vor dem Tag liegt ein Freshness-Audit von Baseline und d-check vor, und der Release liegt in beiden Registries.

**Übernimmt:** `slice-erster-release-veroeffentlichung`, `slice-erster-release-veroeffentlichung-image` — aus `slice-erster-release-veroeffentlichung` nach dem Schnitt vom 2026-10-09 vor
dessen Start (Entscheidung des Nutzers; dort §1, *Abgegeben*, und §4) die Entscheidungen des
Nutzers vom 2026-10-09 (Geber `welle-erster-release`), die jener aus der Welle übernommen hatte:

3. Eine Freigabe-Checkliste in `docs/maintainer/releasing.md` mit einem Beleg je Punkt,
   Anti-Punkten (was ein Release nicht tun darf) und einer Incident-Klausel; der ausgefüllte
   Eintrag des ersten Releases steht in der Closure von welle-erster-release (dort §3).
4. Vor dem Tag ein Freshness-Audit der vendored Baseline und des gepinnten d-check.
5. Ein Scan des Images vor dem Push in eine Registry.

Aus `slice-doku-ist-stand` (dort §6, *Software-Version*): Mit dem ersten echten Tag trägt das
Benutzerhandbuch die veröffentlichte Version im Ist-Zustand — im Kopf *Software-Version* und
*Gültigkeitsbereich* die Ausgabe von `pgwire-recorder version` des veröffentlichten Binaries,
in §11 *Änderungshistorie* statt „Es gibt keine veröffentlichte Version.“ die Zeile dieses
Releases; die Stelle gehört zu DoD-Punkt 3.

Dazu der erste echte Tag: `slice-erster-release-veroeffentlichung` belegt die Pipeline für Tag und Binaries mit dem Trockenlauf (ein Probe-Tag ist nicht vorgesehen, Entscheidung des Nutzers vom 2026-10-19), `slice-erster-release-veroeffentlichung-image` ergänzt das Image; dieser Slice setzt den ersten stabilen Tag mit Checkliste und Scan. Der erste Lauf mit Tag ist die erste Außenwirkung in Registries und am GitHub-Release (`SPEC-067` Punkt 12): Was der Trockenlauf nicht belegt (Schreibrecht der Token, Sichtbarkeit des Pakets, Repository auf Docker Hub, Anlegen des GitHub-Release, Smoke aus den Registries), zeigt erst dieser Lauf. Die Randformen der drei Punkte zogen aus §6 jenes Slice hierher (§6 unten).

**Sendung aus `slice-erster-release-veroeffentlichung-image` — Einhängepunkt des Scans** (Planner, 2026-10-19, `AGENTS.md` §3.13, Schnitt von `slice-erster-release-veroeffentlichung`; `grep` nach der Kennung findet sie in §1, DoD 1 und §8): Der Scan (Punkt 5 oben) ist ein Schritt zwischen `Image` und `Veröffentlichen` in `release.yml` und liest das OCI-Layout, das jener Slice zwischen den Schritten auf der Platte lässt; der Index-Digest `D` bleibt unverändert, weil der Push genau dieses Layout veröffentlicht. Dieser Slice ändert dafür den Text von `release.yml` (Endfassung `1b` des Vertrags `SPEC-067`) und `SPEC-067` Punkte 1 und 5 (Vertragstext: Architect, vor dem Code) und trägt den Schritt in die Zeilenliste der Gegenprobe `release-gegenprobe` ein (Herkunft: dieser Slice, die Fälle von A und B bleiben unverändert).

**Sendung aus `slice-erster-release-veroeffentlichung` und `slice-erster-release-veroeffentlichung-image` — Bezug im Handbuch** (Planner, 2026-10-19, `AGENTS.md` §3.13 und §3.11): Handbuch §2 *Installation* bekommt die Abschnitte zum Bezug der Binaries vom GitHub-Release (Plattformen, Prüfung mit `SHA256SUMS`) und zu den Images aus `ghcr.io` und `docker.io` (Tags `<Version>` und `latest`) erst, wenn der erste echte Tag sie belegt (Ist-Zustand); `README.md` entsprechend. Die Stelle gehört zu DoD-Punkt 3, zusammen mit dem Kopf und §11 des Handbuchs, kein neuer Liefer-Punkt.

**Sendung aus `slice-erster-release-veroeffentlichung` — Belege, die der Trockenlauf nicht erzeugt** (Planner, 2026-10-19, `AGENTS.md` §3.13, Closure jenes Slice; `grep` nach der Kennung findet sie in §1, DoD 3, §4, §6 und §8): Der Trockenlauf von A ist grün (Lauf `39208276024`, dort §7), belegt aber drei Dinge nicht: (a) den Wiederanstoß eines roten Laufs auf GitHub nach der Fehlertabelle von `docs/maintainer/releasing.md` §11 (ein Trockenlauf ist grün gebaut und erzeugt keinen roten Lauf); (b) den Wortlaut `release not found` von `gh release view` für einen nicht vorhandenen Release (`SPEC-067` Punkt 9, akzeptierte Grenze: weicht er ab, zieht der Implementer das Muster in `release.sh` und die Attrappe der Gegenprobe nach, mit Fall und Mutation); (c) das Anlegen des GitHub-Release mit `contents: write` (`SPEC-067` Punkt 12). Der Lauf des ersten Tags zeigt (b) und (c) im Log des Schrittes `GitHub-Release`; (a) nur, wenn er rot endet, dann steht der Wiederanstoß mit Lauf-Kennung in der Checkliste. Endet er grün, bleibt (a) unbelegt, und die Checkliste nennt das als Grenze, ohne es zu behaupten (`AGENTS.md` §3.11). Die Sendung gehört zum dritten Liefer-Punkt (Checkliste und erster Tag), kein neuer Liefer-Punkt.

**Sendung aus `slice-erster-release-veroeffentlichung-image` — Belege, die der Trockenlauf nicht erzeugt** (Planner, 2026-10-19, `AGENTS.md` §3.13, Closure jenes Slice; `grep` nach der Kennung findet sie in §1, DoD 3, §4, §6 und §8): Der Trockenlauf des Images ist grün (Lauf `39232775305`, dort §7) und belegt den Pin des Builders, die Form des Layouts bei zwei Plattformen, die gesetzten Geheimnisse und die Rechte der Datei; er belegt nicht, weil `Veröffentlichen`, `Vergleichen`, `Smoke`, `Latest` und `GitHub-Release` im Trockenlauf nicht laufen (`SPEC-067` Punkt 12): (a) die Wortlaute von `skopeo inspect` für einen nicht vorhandenen Tag in `ghcr.io` und `docker.io` (`manifest unknown`, `name unknown`, `not found`; akzeptierte Grenze, `SPEC-067` Punkt 6) und die Verfügbarkeit von `skopeo` und `gh` im Runner-Image; der Wortlaut steht im Log der Schritte `Veröffentlichen` und `Vergleichen` mit dem Präfix `release: skopeo: `, und weicht er ab, zieht der Implementer das Muster in `release.sh` und die Attrappe der Gegenprobe nach, mit Fall und Mutation je Wortlaut; (b) das Schreibrecht der Token in beiden Registries, die Sichtbarkeit des Pakets in `ghcr.io` (anonymer Zugriff, Punkt 7) und der Smoke aus beiden Registries mit der Referenz `<Registry>/pt9912/pgwire-recorder@<D>`; (c) die Wirkung von `SOURCE_DATE_EPOCH` ohne `ARG` in der Stufe `runtime` (Zeitstempel des Images; Grenze und keine Zusage, `SPEC-067` Punkte 5 und 15), die der Lauf des ersten Tags als Beobachtung nennt, nicht als Zusage; (d) der Wiederanstoß eines roten Laufs mit dem Image (Teil von (a) der Sendung aus `slice-erster-release-veroeffentlichung`, Fehlertabelle `releasing.md` §11). Die Sendung gehört zum dritten Liefer-Punkt (Checkliste und erster Tag), kein neuer Liefer-Punkt.

**Bindung von `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`** (eingetragen vom Planner am 2026-10-18, Entscheidung des Nutzers vom selben Tag, `AGENTS.md` §3.13): Die Freigabe-Checkliste (DoD-Punkt 3, `docs/maintainer/releasing.md` §10) führt als Punkt mit Beleg: ein grünes Ergebnis des Gegenproben-Containers für den Commit, der getaggt wird (der `post-receive`-Hook im Bare-Repo von `pgr-gp` legt es ab; fehlt es oder ist es rot, fährt die Freigabe `make gegenproben` selbst), und `make gegenproben-host` am Host grün auf demselben Commit; die Gegenproben sind seit `slice-harness-gegenproben-aus-gates` nicht mehr Teil von `make gates`. Der Punkt gehört zum Liefer-Punkt der Checkliste, kein neuer Liefer-Punkt, keine neue Schicht. **Geber nach dem Schnitt vom 2026-10-18:** `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (Pflicht, Abschnitt im Dokument, erste Fahrt); der Hook, der das Ergebnis ablegt, liegt in `slice-harness-gegenproben-aus-gates-fahrstelle`. Die Bindung ruht auf beiden; `grep` nach `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` findet sie in §1, §6 und §8 dieses Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Prüfung von Tag und Versionsdatei, Binaries, `SHA256SUMS` und GitHub-Release, Trockenlauf — `slice-erster-release-veroeffentlichung`; Bau des Images, Veröffentlichung mit gleichem
  Index-Digest, `:latest`, Fehlertabelle und Wiederanlauf — `slice-erster-release-veroeffentlichung-image`;
  beide liegen vor diesem Slice, und der Scan setzt sich in deren Pipeline zwischen Bau und Push des Images.
- Die Homebrew-Formel im Tap und ihr Nachweis — `slice-erster-release-homebrew-nachweis` (`welle-v1-nachtrag`); er
  setzt den ersten stabilen Tag dieses Slice voraus.
- Ein Sensor für das Freshness-Audit — `harness/conventions.md` führt den Audit als Handlung,
  nicht als Sensor, weil er Netz über `bash`, `git` und `docker` hinaus braucht; Bestand bleibt
  bewusst stehen. Der Audit vor dem ersten Tag läuft von Hand (Entscheidung des Nutzers vom
  2026-10-18: der erste Release hat Vorrang vor den Nachtläufen). Einen Sensor liefern danach
  `slice-harness-ci-freshness` und `slice-harness-ci-freshness-images` in
  [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md) (Start: welle-erster-release in `done/`); dieser Slice
  baut keinen und liest keinen Nachtlauf als Beleg.
- Die Wahl der ersten Versionsnummer — wird beim Release festgelegt (Entscheidung des Nutzers
  vor dem Tag; dieser Slice trägt sie in die Versionsdatei und die Checkliste ein, wählt sie nicht).
- Konten, Namensräume, Zugangsdaten und Sichtbarkeit der Pakete in den Registries — Handlungen
  des Betreibers (`docs/maintainer/releasing.md` §4 Punkt 4); die Pipeline liest sie so, wie
  `slice-erster-release-veroeffentlichung` es festlegt.
- Ein Nachziehen von Baseline oder d-check, das das Audit findet — ein eigener Slice; dieser
  Slice dokumentiert den Befund (und entscheidet nach §6, ob er den Tag sperrt).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung): Die Pipeline scannt das Image jeder Plattform der Manifestliste vor dem Push, als Schritt zwischen `Image` und `Veröffentlichen`, der das Layout liest und `D` unverändert lässt (Einhängepunkt aus `slice-erster-release-veroeffentlichung-image`, `SPEC-067` Punkt 5; der Text von `release.yml` und `SPEC-067` Punkte 1 und 5 ändern sich mit); ein Befund über der Schwelle beendet den Lauf, bevor eine Registry etwas erhält (Punkt 5, Test mit einem Image über der Schwelle; Beleg in §7: Zusage · Mutation · roter Test, `AGENTS.md` §3.10). Benutzerhandbuch und `README.md` beschreiben, was dieser Slice liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung).
- [ ] `docs/maintainer/releasing.md` §5 beschreibt das Freshness-Audit der vendored Baseline und des gepinnten d-check vor dem Tag mit der Form seines Belegs; das Audit vor dem ersten echten Tag ist mit diesem Beleg ausgeführt (Punkt 4).
- [ ] [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew): `docs/maintainer/releasing.md` §10 trägt die Freigabe-Checkliste mit Beleg je Punkt, Anti-Punkten und Incident-Klausel (Punkt 3). Der erste echte, stabile Tag ist nach ihr freigegeben: Binaries mit `SHA256SUMS` am Release, Image in beiden Registries mit gleichem Index-Digest und `:latest`, Smoke aus Abnahmeszenario 9 bestanden. Der ausgefüllte Eintrag steht in der Closure von welle-erster-release. Die Checkliste und die Closure nennen aus dem Lauf des ersten Tags die Belege (b) und (c) und, falls der Lauf rot war, (a) (Sendung aus `slice-erster-release-veroeffentlichung`, §1; Beleg: Lauf-Kennung und Zeile im Log des Schrittes `GitHub-Release`, bei rotem Lauf der Wiederanstoß nach `releasing.md` §11). Aus `slice-erster-release-veroeffentlichung-image` (Sendung, §1) nennen Checkliste und Closure aus dem Lauf des ersten Tags: die Wortlaute von `skopeo` für einen nicht vorhandenen Tag und die Verfügbarkeit von `skopeo` und `gh` (Beleg: Zeilen `release: skopeo: ` im Log von `Veröffentlichen` und `Vergleichen`), das Schreibrecht der Token, die Sichtbarkeit des Pakets und den Smoke aus beiden Registries (Beleg: Lauf-Kennung, Digest je Registry), die Wirkung von `SOURCE_DATE_EPOCH` (Beleg: Zeitstempel des veröffentlichten Images, als Beobachtung); was der Lauf nicht zeigt, nennt die Checkliste als Grenze ohne Behauptung (`AGENTS.md` §3.11). Handbuch §2 *Installation* und `README.md` tragen die Abschnitte zum Bezug der Binaries vom GitHub-Release und der Images aus den Registries, belegt durch den ersten Tag (Sendung aus `slice-erster-release-veroeffentlichung` und `slice-erster-release-veroeffentlichung-image`, §1; Beleg: Bezug des veröffentlichten Binaries beziehungsweise Images und `grep` im Handbuch). Aus `slice-doku-ist-stand`: Das Benutzerhandbuch nennt im Kopf (*Software-Version*, *Gültigkeitsbereich*) die Version des veröffentlichten Binaries und trägt in §11 die Zeile des Releases an Stelle von „Es gibt keine veröffentlichte Version.“ (Beleg: `grep` im Handbuch und `version` des veröffentlichten Binaries).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Release-Automatisierung (`release.yml`, `release.sh`, Gegenprobe `release-gegenprobe`) | update | Scan des Images vor dem Push, mit Schwelle, zwischen `Image` und `Veröffentlichen` (Sendung, §1); das Muster der Wortlaute von `skopeo` und der Attrappe, falls der erste Tag abweichende zeigt (Sendung aus `slice-erster-release-veroeffentlichung-image`, §1) |
| `docs/maintainer/releasing.md` §5, §10 | update | Freshness-Audit vor dem Tag, Freigabe-Checkliste |
| Test der Pipeline | neu | Image über der Schwelle wird nicht veröffentlicht |
| `docs/plan/planning/welle-erster-release-results.md` | neu (bei Welle-Closure) | ausgefüllter Eintrag der Checkliste des ersten Releases |
| `docs/user/benutzerhandbuch.md`, `README.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung); Kopf *Software-Version* und §11 mit der veröffentlichten Version (aus `slice-doku-ist-stand`); Abschnitte zum Bezug der Binaries und Images (Sendung, §1) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-erster-release-veroeffentlichung` und `slice-erster-release-veroeffentlichung-image` liegen in `done/`
(Pipeline für Tag, Binaries und Image mit dem Trockenlauf belegt; ein Probe-Tag ist nicht vorgesehen, Entscheidung des Nutzers vom 2026-10-19). Schritt 6 der Reihenfolge in §5 von
[welle-erster-release](../welle-erster-release.md). Vor dem ersten Code-Commit entscheidet der
Architect die Randformen aus §6 (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Scan verlangt einen Umbau der
  Pipeline, der den Diff über eine Review-Sitzung hebt; Schnitt dann: *Scan* (DoD-Punkt 1)
  und *Freigabe mit erstem Tag* (DoD-Punkte 2 und 3). Marke (Metrik nach `AGENTS.md` §3.14): die
  Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'`
  liegt über 1200, Tests und Gegenprobe eingeschlossen; Basis: beim Anspruch (HEAD vor dem
  Anspruch, hier einzutragen). Schätzung des Planners in dieser Metrik (2026-10-17, ohne Basis
  gemessen): Scan in der Pipeline 120 bis 200, Gegenprobe des Scans mit Attrappe das Doppelte
  (240 bis 400), `docs/maintainer/releasing.md` §5 und §10 150 bis 250, Handbuch und `README.md`
  40 bis 80, zusammen 550 bis 930 (die Abschnitte zum Bezug im Handbuch und die Anpassung der Zeilenliste der Gegenprobe aus den Sendungen vom 2026-10-19 kommen hinzu: rund 60 bis 120, zusammen 610 bis 1050, Mitte 830, 69 % der Marke; Planner, Schätzung ohne Messung; die Sendung aus `slice-erster-release-veroeffentlichung` vom 2026-10-19 fügt höchstens das Nachziehen des Musters von `gh release view` mit Fall und Mutation hinzu, 20 bis 60, zusammen 630 bis 1110, Mitte 870, 73 % der Marke; die Sendung aus `slice-erster-release-veroeffentlichung-image` vom 2026-10-19 fügt höchstens das Nachziehen der Muster der Wortlaute von `skopeo` hinzu (Muster, Attrappe, Fall und Mutation je Wortlaut, 30 bis 80), zusammen 660 bis 1190, Mitte 925, 77 % der Marke, obere Grenze 99 %; der Messwert aus jenem Slice, Faktor Gegenprobe zu Produktcode rund 3,7 statt 2, ist in dieser Schätzung nicht eingerechnet, der Planner rechnet sie beim Anspruch damit nach); die Schätzung überschreitet die Marke nicht, ein Schnitt vor dem
  Start unterbleibt. Meldeschwelle 900: Wer sie überschreitet, meldet es vor dem nächsten
  Abschnitt, und die Entscheidung („weiter“ oder „schneiden“) steht als datierte Zeile in §7. Die
  Schätzung hängt an der Wahl des Scanners (§6) und am Zuschnitt der Pipeline durch den Vorgänger;
  der Planner rechnet sie nach dessen Closure nach, bevor der Slice den Anspruch bekommt.
- `in-progress` → `open` (blockiert — Carveout?): Der Scan meldet für das Basis-Image einen
  Befund über der Schwelle ohne verfügbare Behebung; dann zuerst die Entscheidung des Nutzers
  zur Schwelle oder ein Carveout, kein Tag.

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

- Der erste echte Tag ist nicht zurückzunehmen (Rollback-Regel in `docs/maintainer/releasing.md`
  §11); ein Fehler, den Checkliste oder Scan nicht fangen, ist veröffentlicht und nur durch
  einen höheren Tag zu ersetzen — **Ausgang:** offen bis Closure.
- Der Scan hängt an einer Schwachstellen-Datenbank aus dem Netz; derselbe Quellstand kann an
  zwei Tagen verschieden bewertet werden — **Ausgang:** offen bis Closure.
- Das Freshness-Audit kann einen neueren Stand von Baseline oder d-check finden; ob das den Tag
  sperrt, ist eine Randform (unten), und ein Nachziehen wäre ein eigener Slice — **Ausgang:**
  offen bis Closure.
- Der erste Tag setzt Außenwirkung in zwei Registries und am GitHub-Release; bricht die Pipeline
  dazwischen ab, liegt ein halber erster Release vor, und der Tag ist nicht zu verändern
  (`docs/maintainer/releasing.md` §11) — **Ausgang:** offen bis Closure; die Randform
  *Abbruch beim ersten Tag* (unten) fängt ihn.
- Der Trockenlauf belegt die Außenwirkung nicht (Schreibrecht der Token, Sichtbarkeit des Pakets, Repository auf Docker Hub, Anlegen des GitHub-Release, Smoke aus den Registries; `SPEC-067` Punkt 12); der erste Lauf mit einem Tag ist deshalb die erste Außenwirkung (Sendung aus `slice-erster-release-veroeffentlichung` und `-image`, die Risiken dort §6 weisen hierher; der Trockenlauf von A ist grün, der Wiederanstoß eines roten Laufs, der Wortlaut von `gh release view` und das Anlegen des GitHub-Release bleiben für den Lauf des ersten Tags, §1; aus `slice-erster-release-veroeffentlichung-image`, Trockenlauf `39232775305` grün: die Wortlaute von `skopeo` und die Verfügbarkeit von `skopeo` und `gh`, das Schreibrecht der Token, die Sichtbarkeit des Pakets, der Smoke aus den Registries, die Wirkung von `SOURCE_DATE_EPOCH` und der Wiederanstoß mit dem Image, ebenfalls §1) — **Ausgang:** offen bis Closure; die Randform *Abbruch beim ersten Tag* fängt ihn.
- Der Nutzer gibt den Tag frei: Das Veröffentlichen ist eine eigene, ausdrückliche Entscheidung des
  Nutzers und kein Schritt, den der Implementer von sich aus tut — **Ausgang:** offen bis Closure.

**Randformen** (`AGENTS.md` §3.12) — **offen für den Architect, außer wo angegeben**, zu entscheiden vor dem
ersten Code-Commit in der Spezifikation oder einer ADR; sie zogen mit dem Schnitt vom
2026-10-09 aus §6 von `slice-erster-release-veroeffentlichung` hierher (der Abbruch beim ersten Tag ist seit dem 2026-10-19 in `SPEC-067` Punkt 11 entschieden, siehe unten):

- *Freigabe-Checkliste (Punkt 3)* — Form des Belegs je Punkt (Lauf-Kennung, Digest, Link);
  welche Anti-Punkte; Inhalt der Incident-Klausel (Auslöser, wer, wo festgehalten); wo der
  ausgefüllte Eintrag eines späteren Releases ohne Welle steht.
- *Freshness-Audit (Punkt 4)* — was als frisch gilt; ob ein neuerer Stand der Baseline oder
  von d-check den Tag sperrt oder nur dokumentiert wird; Beleg des Audits. Der Audit braucht
  Netz; `harness/conventions.md` führt ihn als Handlung, nicht als Sensor.
- *Abbruch beim ersten Tag und Wiederholung* — entschieden in `SPEC-067` Punkt 11 (Architect, 2026-10-19): ein roter Lauf wird mit demselben Tag wieder angestoßen (Re-run oder erneuter Lauf), solange kein Digest und kein Anhang abweicht; sonst verlangt der Fehler einen höheren Tag; ein Git-Tag wird nie verändert, ein Tag in einer Registry nie überschrieben. Offen für die Freigabe: wer den Lauf des ersten Tags anstößt und wo der Abbruch in der Checkliste festgehalten wird.
- *Reihenfolge von Checkliste, Versionszeile und Tag* — die Zeile der Änderungshistorie und die
  Versionsdatei stehen vor dem Tag (§5 von `releasing.md`), das Handbuch nennt im Kopf die
  Version, die das veröffentlichte Binary meldet; ob der Kopf nach dem Release die Version oder
  weiter `dev` für Builds aus dem Repository nennt (`make build` meldet `dev`), und in welchem Commit
  (vor dem Tag) das steht.
- *Form des Tags* — ob der erste echte Tag mit `v0.` oder `v1.` beginnt (Entscheidung
  des Nutzers, nicht des Architect). Ein Probe-Tag entfällt (Entscheidung des Nutzers vom 2026-10-19,
  `SPEC-067` Punkt 12); `latest` setzt der Lauf des ersten stabilen Tags nach grünem Vergleich und Smoke (`SPEC-067` Punkt 8).
- *Image-Scan (Punkt 5)* — Scanner und sein Pin; Schwelle (Schweregrad, nur behebbare
  Befunde?); Ausnahmen und ihr Ort (`AGENTS.md` §3.2, Suppression-Verbot); Scan je Plattform
  der Manifestliste; Netz für die Schwachstellen-Datenbank; Stelle: zwischen `Image` und
  `Veröffentlichen` (`SPEC-067` Punkt 5, Sendung in §1).

- **Gegenproben nicht Teil von `make gates`** (Bindung aus `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`): Ein grünes `make gates` vor dem Tag sagt nichts über die Gegenproben; ohne den Punkt in der Checkliste bliebe eine Regression im Werkzeug unbemerkt. — **Ausgang (vorgesehen):** entfallen, wenn die Checkliste den Punkt mit Beleg trägt; sonst eingetreten → Folge-Slice mit Kennung vom Planner. Reihenfolge: Die Welle dieses Slice liegt vor den Harness-Wellen; trägt `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` nach, wenn dieser Slice schon `done/` ist, ist das Risiko eingetreten.

**Entschieden vom Architect am 2026-10-19** (Ort: [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkte 16 bis 18; ADR: [ADR-0052](../../adr/0052-scan-des-images-vor-dem-push.md), Proposed, die Annahme entscheidet der Nutzer; der Start des Slice setzt sie voraus, weil eine Schwelle eine Entscheidung ist, `AGENTS.md` §3.6):

- *Image-Scan:* Trivy als Image mit Digest-Pin (Vorbild d-check `0.74.0`), `--exit-code 0` und Auswertung des JSON; Schwelle `CRITICAL` und `HIGH` mit verfügbarer Behebung, was ohne Behebung gemeldet wird, zählt und sperrt nicht; **keine Ausnahmeliste**; je Plattform, die gescannte Architektur wird verglichen; Datenbank bei jedem Lauf frisch, älter als 72 Stunden ist Exit 2; Zeitgrenze `--timeout 10m`; Exit 0, 1, 2 mit den Stichwörtern `schwachstelle`, `scan-gescheitert`, `pruefmenge`, `datenbank`, `plattform`, `vorbedingung`; der Scan läuft auch im Trockenlauf; Marke `scan.ok`, `veröffentlichen` verlangt sie — Punkt 16.
- *Freshness-Audit:* Handlung von Hand, Tabelle als Beleg (gepinnter Stand, neuester stabiler Release der Liste, Datum der Abfrage höchstens sieben Tage vor dem Tag, Ergebnis); ein neuerer Stand sperrt den Tag nicht von selbst, er wird dokumentiert; bei einer Sicherheitskorrektur in den Notizen entscheidet der Nutzer vor dem Tag — Punkt 17.
- *Freigabe-Checkliste:* Tabelle `Punkt | Form des Belegs | Beleg` in `releasing.md` §10, zwölf Punkte, Anti-Punkte, Incident-Klausel (neuer, höherer Tag, der fehlerhafte bleibt); Ort des Eintrags: für den ersten Release die Closure von `welle-erster-release`, später `docs/maintainer/freigaben/<Version>.md`; Belege der Gegenproben: die Dateien `gegenproben-*-<Hash>` der Ablage von `pgr-gp` und die Zusammenfassung von `make gegenproben-host`, beide für den Hash, der `VERSION` trägt und getaggt wird — Punkt 17. Kein Fall (Handlung).
- *Sendungen von A und B:* der Lauf des ersten Tags belegt sie im Log (Wortlaute `skopeo` und `gh`, Verfügbarkeit, Schreibrecht, Sichtbarkeit, Smoke aus den Registries, Zeitstempel als Beobachtung); „rot“ danach: Fehler der Umgebung ⇒ Wiederanstoß desselben Tags, Fehler im Skript oder Vertrag ⇒ neuer, höherer Tag (der Lauf führt die Fassung des Tag-Commits aus); der Nutzer wird gefragt, ob `v0.1.1` abgesprochen ist — Punkt 18.
- *Handbuch und README:* Zeile in §11 vor dem Tag (Quelle der Notizen); Kopf *Software-Version* und §2 Installation in einem Commit nach dem Lauf, nur mit Befehlen, die nach dem Tag einmal gelaufen sind (Beleg in §7); `latest` wird genannt, weil der Lauf ihn gesetzt hat — Punkt 18.
- *Form des Tags:* `v0.1.0` (Entscheidung des Nutzers vom 2026-10-18); *Reihenfolge Checkliste, Versionszeile und Tag:* `VERSION` und Zeile im Commit, der getaggt wird; die Gegenproben laufen auf genau diesem Hash — Punkt 17.

**Schätzung des Architect nach den Verträgen** (2026-10-19; Metrik von §4; Faktor 2,0 für Gegenproben mit Attrappen, 3,7 für zeilenweise gehaltenen Text, beide aus A und B gemessen): Schritt `scan` in `release.sh` 120 bis 170, vier Zeilen in `release.yml` (die Zeilenliste der Gegenprobe wächst um sie, mit Faktor 3,7 rund 15), Gegenprobe des Scans (rund 28 Fälle, rund 30 Mutationen, Attrappe von `docker` und `jq`-Eingaben) 300 bis 450, `releasing.md` §5 und §10 150 bis 250, Handbuch und `README.md` 100 bis 200 (Bezug und Kopf, nach dem Tag), das Nachziehen der Muster der Wortlaute von `skopeo` und `gh` 0 bis 120 (nur wenn sie abweichen); zusammen **685 bis 1205**, Mitte um 945, 79 % der Marke 1200; die Meldeschwelle 900 wird in der oberen Hälfte erreicht. Schnitt, falls sie eintritt: *Scan* (Punkt 16 mit seiner Gegenprobe) und *Checkliste mit erstem Tag* (Punkte 17 und 18); der Planner rechnet beim Anspruch mit der Basis nach.

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

- **Belege zur DoD (Implementer):** <…>
- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*`
(Kürzel `REPO`, Greenfield); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `3f7f52d` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/verkoerperte-regel-nach-baseline-sprung-nicht-nachgezogen` (1×) — das
  Freshness-Audit (Punkt 4) kann einen Baseline-Sprung auslösen; der Sprung selbst ist nicht
  Teil dieses Slice.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (16×, verkörpert in `AGENTS.md` §3.12) —
  die Randformen stehen offen in §6 und gehen vor dem Code an den Architect.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (17×, `AGENTS.md` §3.10) — der Scan ist ein
  neuer Vertrag; je Zusage eine Mutation.
- `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (3×, §3.13) — der Geber nennt diesen Slice in
  §1 unter *Abgegeben*, dieser ihn unter *Übernimmt*; `slice-erster-release-homebrew-nachweis`
  setzt ihn in §4 voraus, alles im selben Commit.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

Nachgesichtet am 2026-10-17 (Planner, Vorbereitung des Anspruchs): `BEO-REPO/postgres-versionsliste-altert` (2×, offen) betrifft den Smoke aus Abnahmeszenario 9 mit `POSTGRES_IMAGE` nur mittelbar; der Release ändert die Liste nicht. `BEO-REPO/verhalten-nur-unter-linux-geprueft` trifft die Aussage des Handbuchs zu den Plattformen: Binaries für `darwin` und `windows` sind nur gebaut, nicht gelaufen (`SPEC-055` Punkt 13, Grenze 4); die Plattform-Zeile in §1 des Handbuchs sagt nur das zu.

Nachgezählt beim Eintragen der Sendung aus `slice-doku-ist-stand` (2026-10-10, `AGENTS.md` §3.13): Kopf und §11 des Handbuchs liegen in den zwei Dateien, die §3 schon führt, und in DoD-Punkt 3; kein neuer Liefer-Punkt. Schichten nach der Zählweise von `slice-harness-d-check-v0-85` §8 (Nutzer- und Wartungs-Doku `docs/user/` und `docs/maintainer/` sind eine Schicht): Release-Automatisierung und Dokumentation, zwei; die Sendung fügt keine hinzu.

Nachgezählt am 2026-10-17 (Planner, `AGENTS.md` §3.13; keine neue Sendung, Stand nach den Schnitten von `slice-v1-abschluss-container` und `-container-doku`): Liefer-Punkte drei (§2: Scan mit Handbuch- und `README.md`-Teil | Freshness-Audit in `releasing.md` §5 | Checkliste und erster Tag mit Handbuch-Kopf), Schichten zwei. Die Wahl des Scanners (§6) kann `Dockerfile` oder `harness/` als Bereich öffnen; der Architect zählt mit der Entscheidung nach, und läge der Slice darüber, schneidet der Planner entlang der Linie aus §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit, der die Bindung aus `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` einträgt): Der Punkt gehört in die Checkliste des dritten Liefer-Punkts (DoD 3); Liefer-Punkte weiter drei, Schichten weiter zwei (Release-Automatisierung, Dokumentation).

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit, der die Ergänzung zur Bindung aus `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` einträgt): Die Fassung des Punktes ändert nur den Beleg (Ergebnis des Gegenproben-Containers statt eigener Fahrt); Liefer-Punkte weiter drei, Schichten weiter zwei.

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit, der die Bindung aus `slice-harness-ci-freshness` einträgt): Die Bindung ändert den Beleg von Punkt 4 (Ergebnis des Nachtlaufs statt eigener Handlung) und den Start-Trigger; sie ist Teil des zweiten Liefer-Punkts (Freshness-Audit in `releasing.md` §5), kein neuer. Liefer-Punkte weiter drei, Schichten weiter zwei (Release-Automatisierung, Dokumentation).

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit, der `slice-harness-gegenproben-aus-gates-fahrstelle` schneidet): Die Bindung aus `slice-harness-gegenproben-aus-gates-fahrstelle` wechselt den Geber zu `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (Pflicht, Abschnitt im Dokument, erste Fahrt); der Inhalt des Punktes der Checkliste ändert sich nicht. Der Punkt gehört weiter zur Checkliste des dritten Liefer-Punkts (DoD 3); Liefer-Punkte weiter drei, Schichten weiter zwei (Release-Automatisierung, Dokumentation).

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-18, im Commit, der die Bindung aus `slice-harness-ci-freshness` zurücknimmt; Entscheidung des Nutzers vom 2026-10-18: der erste Release hat Vorrang vor den Nachtläufen): Der Audit vor dem Tag (Punkt 4) läuft wieder von Hand, der Nachtlauf ist weder Beleg noch Start-Trigger, `SPEC-066` ist aus dem Kopf entfernt, der Ausschluss eines Sensors im Abschnitt *Ausdrücklich NICHT* gilt wieder; der Sensor folgt in [welle-ci-nachtlaeufe](../welle-ci-nachtlaeufe.md). Die Rücknahme ändert nur den Beleg von Punkt 4 und den Start-Trigger, sie ist Teil des zweiten Liefer-Punkts (Freshness-Audit in `releasing.md` §5), kein neuer. Liefer-Punkte weiter drei, Schichten weiter zwei (Release-Automatisierung, Dokumentation); Schätzung und Marke in §4 unverändert (Marke 1200, Meldeschwelle 900).

**Nachgezählt am 2026-10-19** (Planner, `AGENTS.md` §3.13, im Commit, der die Sendungen aus `slice-erster-release-veroeffentlichung` und `slice-erster-release-veroeffentlichung-image` einträgt: Einhängepunkt des Scans zwischen `Image` und `Veröffentlichen` samt Änderung von `release.yml`, Zeilenliste der Gegenprobe und `SPEC-067` Punkte 1 und 5; Abschnitte des Handbuchs zum Bezug der Binaries und Images): Der Einhängepunkt gehört zum ersten Liefer-Punkt (Scan, DoD 1), der Handbuch-Bezug zum dritten (Checkliste und erster Tag mit Handbuch-Kopf, DoD 3); Liefer-Punkte weiter drei (Scan mit Handbuch- und `README.md`-Teil | Freshness-Audit in `releasing.md` §5 | Checkliste und erster Tag mit Handbuch), Schichten weiter zwei (Release-Automatisierung, Dokumentation). Die Zählung liegt unter der Grenze; die Schätzung steigt um 60 bis 120 (§4), die Marke 1200 und die Meldeschwelle 900 bleiben. Der Start verlangt jetzt zusätzlich `slice-erster-release-veroeffentlichung-image` in `done/` (§4).

**Nachgezählt am 2026-10-19** (Planner, `AGENTS.md` §3.13, im Commit, der die Sendung aus `slice-erster-release-veroeffentlichung` einträgt: die Belege des Wiederanstoßes, des Wortlauts von `gh release view` und des Anlegens des GitHub-Release aus dem Lauf des ersten Tags): Die Sendung ist Beleg und Bedingung im dritten Liefer-Punkt (Checkliste und erster Tag mit Handbuch, DoD 3); ein Nachziehen des Musters in `release.sh`, falls der Wortlaut abweicht, liegt in der Release-Automatisierung, die der erste Punkt schon führt. Liefer-Punkte weiter drei (Scan mit Handbuch- und `README.md`-Teil | Freshness-Audit in `releasing.md` §5 | Checkliste und erster Tag mit Handbuch), Schichten weiter zwei (Release-Automatisierung, Dokumentation). Die Schätzung steigt um höchstens 60 auf 630 bis 1110 (§4), die Marke 1200 bleibt, die obere Grenze liegt über der Meldeschwelle 900; unter der Marke, ein Schnitt vor dem Start unterbleibt. Der Messwert aus `slice-erster-release-veroeffentlichung` (1425 gegen 820 bis 1070) macht die Nachrechnung beim Anspruch zur Bedingung, nicht zur Option.

**Nachgezählt am 2026-10-19** (Planner, `AGENTS.md` §3.13, im Commit, der die Sendung aus `slice-erster-release-veroeffentlichung-image` einträgt: die Belege, die der Trockenlauf des Images nicht erzeugt): Die Sendung ist Beleg und Bedingung im dritten Liefer-Punkt (Checkliste und erster Tag mit Handbuch, DoD 3); ein Nachziehen der Muster der Wortlaute von `skopeo`, falls sie abweichen, liegt in der Release-Automatisierung, die der erste Punkt schon führt. Liefer-Punkte weiter drei (Scan mit Handbuch- und `README.md`-Teil | Freshness-Audit in `releasing.md` §5 | Checkliste und erster Tag mit Handbuch), Schichten weiter zwei (Release-Automatisierung, Dokumentation). Die Schätzung steigt um höchstens 80 auf 660 bis 1190 (§4), obere Grenze 99 % der Marke 1200 und über der Meldeschwelle 900; unter der Marke, ein Schnitt vor dem Start unterbleibt, die Nachrechnung beim Anspruch bleibt Bedingung (der Faktor Gegenprobe zu Produktcode liegt je nach Form der Gegenprobe bei 2 bis 3,7: Attrappen gegen Aufrufe 2, Text einer Datei zeilenweise 3,7; die Gegenprobe des Scans hält Aufrufe gegen eine Attrappe).
