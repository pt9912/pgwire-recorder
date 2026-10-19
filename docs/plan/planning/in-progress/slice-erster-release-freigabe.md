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

**Handbuch und `README.md`:** Ihr Teil (DoD-Punkt 3: Kopf, §11, §2 *Installation*, Bezug) entsteht in einem eigenen Commit nach dem Lauf des Tags (`SPEC-067` Punkt 18), nicht in diesem Zweig; der Entwurf steht in §7.

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

- [ ] [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung): Die Pipeline scannt das Image jeder Plattform der Manifestliste vor dem Push, als Schritt zwischen `Image` und `Veröffentlichen`, der das Layout liest (je Plattform als Verzeichnis mit dem Manifest dieser Plattform, `SPEC-067` Punkt 16) und `D` unverändert lässt (Einhängepunkt aus `slice-erster-release-veroeffentlichung-image`, `SPEC-067` Punkt 5; der Text von `release.yml` und `SPEC-067` Punkte 1 und 5 ändern sich mit); ein Befund über der Schwelle beendet den Lauf, bevor eine Registry etwas erhält (Punkt 5, Test mit einem Image über der Schwelle; Beleg in §7: Zusage · Mutation · roter Test, `AGENTS.md` §3.10). Der Scan ändert kein Verhalten des Binaries; Handbuch und `README.md` tragen dazu nichts (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung).
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
| Release-Automatisierung (`release.yml`, `release.sh`, Gegenprobe `release-gegenprobe`) | update | Scan des Images vor dem Push, mit Schwelle, zwischen `Image` und `Veröffentlichen` (Sendung, §1), das Layout je Plattform als Verzeichnis (Trockenlauf `39263455592`); das Muster der Wortlaute von `skopeo` und der Attrappe, falls der erste Tag abweichende zeigt (Sendung aus `slice-erster-release-veroeffentlichung-image`, §1) |
| `docs/maintainer/releasing.md` §1, §2, §5, §6, §8, §10, §11 | update | Freshness-Audit vor dem Tag, Freigabe-Checkliste, Schritt `Scan` in der Beschreibung des Tags, Zeile `scan` der Fehlertabelle |
| `harness/README.md` (Zeile `make release-gegenprobe`) | update | die Gegenprobe prüft zusätzlich den Scan (`p16-*`, Werkzeuge `jq` und GNU `date`) |
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
  liegt über 1200, Tests und Gegenprobe eingeschlossen; Basis: 760c51d4 (HEAD vor dem
  Anspruch). Schätzung des Planners in dieser Metrik (2026-10-17, ohne Basis
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

- *Image-Scan:* Trivy als Image mit Digest-Pin (Vorbild d-check `0.74.0`), `--exit-code 0` und Auswertung des JSON; Schwelle `CRITICAL` und `HIGH` mit verfügbarer Behebung, was ohne Behebung gemeldet wird, zählt und sperrt nicht; **keine Ausnahmeliste**; je Plattform, das Layout je Plattform als Verzeichnis mit dem Manifest dieser Plattform (Trivy liest ein Layout nur als Verzeichnis und bei zwei Manifesten das erste, Trockenlauf `39263455592`, `SPEC-067` Punkt 16), die gescannte Architektur wird verglichen; Datenbank bei jedem Lauf frisch, älter als 72 Stunden ist Exit 2; Zeitgrenze `--timeout 10m`; Exit 0, 1, 2 mit den Stichwörtern `schwachstelle`, `scan-gescheitert`, `pruefmenge`, `datenbank`, `plattform`, `vorbedingung`; der Scan läuft auch im Trockenlauf; Marke `scan.ok`, `veröffentlichen` verlangt sie — Punkt 16.
- *Variante des Manifests (Randform, entschieden vom Architect in `SPEC-067` Punkt 16):* die Auswahl des Manifests je Plattform vergleicht `os` und `architecture`, ein `variant` wird ignoriert, zwei Manifeste derselben Plattform mit verschiedener Variante sind doppelt (`plattform`).
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

### Belege des Implementers

Stand: Kopf des Codes `de48c534` auf dem Zweig `impl-fr` (Basis `e2372b18`, Basis der Metrik `760c51d4`; Code `07250783`, `releasing.md` `8ea23258`, Erwartung einer Mutation von B nachgezogen `de48c534`); Läufe im Dev-Container pgr-dev3 (Docker im Container), nacheinander am Endstand. Während der Arbeit nur die engsten Sensoren: `bash -n`, Teilläufe der Gegenprobe mit `GP_MUT` und `GP_ENTDECKE` (die roten Fälle je Mutation gemessen, dann die Erwartung im Skript gesetzt und gegen die Vermutung gehalten), `make workflows-check`, `make docs-check`, zuletzt die Läufe unten je einmal am Endstand.

| Lauf | Ergebnis | Dauer |
|---|---|---|
| `make release-gegenprobe` | `release-gegenprobe: grün (491 Läufe)`: die unveränderte Kopie, die Mutationen von A und B (mit drei geänderten Erwartungen, oben) und 87 Mutationen `m-p16-*`, der Rekursionsschutz und die Probe der Fallliste | 1237 s |
| `make workflows-check` | Exit 0, `workflows-check: 2 Datei(en) geprüft, kein Befund` | 1 s |
| `make docs-check` | Exit 0, `d-check: 1245 Datei(en) geprüft, 0 Befund(e)` (auf `de48c534`, vor diesem Eintrag; nach ihm wiederholt, siehe unten) | 2 s |
| `make gates` | Exit 0 (darunter `test`, `test-integration`, `test-container`, `test-postgres-versionen`, `lint`, `kopf-check`, `workflows-check`, `docs-check`) | 266 s |

Nicht gelaufen, mit Grund: der Scanner selbst (kein Netz aus dem Container, `AGENTS.md` und Auftrag: keine Registry-Aufrufe); der Trockenlauf auf GitHub (nicht Aufgabe des Implementers: kein Push nach GitHub); die Auflösung des Digest von Trivy gegen den Docker Hub (die Messung des Hosts vom 2026-10-19 steht im Auftrag: `sha256:af6acf9a…855daa`, Tag 0.75.0, identisch auf Docker Hub und `ghcr.io`, veröffentlicht 2026-10-01; aus dem Container nicht nachgemessen); `make gates-schnell-gegenprobe` (seit `e2372b18` ändert sich weder `harness/mk/release.mk` noch eine Liste `GEGENPROBEN_CONTAINER`).

**Messung der Größe** (`AGENTS.md` §3.14): `git diff --numstat 760c51d4..de48c534 -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'` ergibt als Summe aus Einfügungen und Löschungen **527** (`release-gegenprobe.sh` 318, `releasing.md` 133, `release.sh` 72, `release.yml` 2, `harness/README.md` 2). Marke 1200 (44 %) und Meldeschwelle 900 sind nicht erreicht, die Schätzung des Architect (685 bis 1205) liegt darüber: `releasing.md` §5 und §10 blieben mit 133 unter 150 bis 250, und Handbuch und `README.md` (100 bis 200 in der Schätzung) liefert dieser Lauf nicht (Entwurf unten, nach dem Tag). Die Gegenprobe trägt 318 von 527 Zeilen, das 4,3-fache des Produktcodes (`release.sh` und `release.yml`, zusammen 74): Aufrufe gegen eine Attrappe, aber 76 Fälle mit je einem Mutanten und 87 Mutationen; der Faktor 2 des Architect gilt für diese Form nicht, 3,7 aus `slice-erster-release-veroeffentlichung` liegt näher.

**Mutationstabelle** (`AGENTS.md` §3.10): Zusage (Punkt von `SPEC-067`) · Mutation (Textersatz in einer Kopie von `release.sh`; Namen wie in `tools/harness/release-gegenprobe.sh`) · roter Fall. Je Mutation sind genau die genannten Fälle rot, jede einmal gesehen (`GP_ENTDECKE=1`), die unveränderte Kopie ist grün. Bei mehr als sechs roten Fällen nennt die Zeile die Zahl und drei Namen; die vollständige Menge steht in der Ausgabe von `GP_ENTDECKE=1`, im Skript als Liste oder als `ausser …`. Neu sind 99 Fälle (94 `p16-*` in der Gruppe `f`, `p16-werkzeug-jq` und `-date` in der Gruppe `gp13`, `p13-veroeffentlichen-ohne-image-bei-scan-ok`, `p11-alte-marke-image-b` und `p11-image-rot-loescht-digest-b` in der Gruppe `b`) und 107 Mutationen `m-p16-*` (nach der Behebungsrunde, unten); die Fälle von A und B sind unverändert, ihre Fixture `bild` und `latest_von` fahren den Schritt `scan` mit. Drei Erwartungen von Mutationen von B ändern sich, weil ein zweiter Wächter (`scan.ok`) hinter dem ersten steht: `m-p13-veroeffentlichen-ohne-image`, `m-p11-eingaben-loescht-image-nicht` und `m-p11-image-loescht-nicht` färben jetzt die neuen Fälle mit gefälschter Marke `scan.ok` (Zeilen unten), `m-p5-summe-amd64-nicht-geprueft` zusätzlich `p11-image-rot-loescht-digest-b`.

Zeilenliste von `release.yml` (Punkt 1, Herkunft F): Zeilen 40 und 41 des Textes (`- name: Scan`, `run: bash tools/harness/release.sh scan`); die Schleifen von A und B erzeugen die Fälle `p1-fehlt-40`, `p1-fehlt-41`, `p1-geaendert-40`, `p1-geaendert-41` und die Mutationen `m-p1-zeile-40-optional`, `-41-optional`, `-40-nicht-verglichen`, `-41-nicht-verglichen` (je genau ihr Fall rot), dazu die Paare vertauschter Schritte `image-scan` und `scan-veroeffentlichen` (Fälle `p1-reihenfolge-image-scan`, `p1-reihenfolge-scan-veroeffentlichen`, Mutationen `m-p1-reihenfolge-<Paar>-tauschbar`, je genau ihr Fall rot); `image-veroeffentlichen` tauscht jetzt die Stellen 2 und 4.

| Zusage (Punkt von SPEC-067) | Mutation | roter Fall |
|---|---|---|
| 16 Pin: `aquasec/trivy@sha256:<64 Hexziffern>` | `m-p16-pin-ohne-digest`: der Pin trägt einen Tag statt des Digests | `p16-pin-form` |
| 16 Pin: ein Platzhalter aus Nullen ist ungültig | `m-p16-pin-platzhalter`: der Digest besteht aus Nullen | `p16-pin-platzhalter` |
| 16 Pin: Tag als Kommentar hinter der Zeile | `m-p16-pin-ohne-tag-kommentar`: der Kommentar `# 0.75.0` fehlt | `p16-pin-form` |
| 16 E1: `docker run --rm` | `m-p16-e1-rm-entfaellt`: `--rm` entfällt | `p16-e1` |
| 16 E2: Verzeichnis der Plattform als Mount `-v <Verzeichnis>:/layout:ro` | `m-p16-e2-ro-entfaellt`: `:ro` entfällt | `p16-e2`, `p16-e2-quelle` |
| 16 E2/Layout je Plattform: die Quelle ist ein Verzeichnis im Temp-Verzeichnis des Schrittes, nicht das Verzeichnis von `image` | `m-p16-e2-quelle-image-verzeichnis`: der Mount zeigt auf `$zdir/image` | 90 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 E3: das Image des Pins | `m-p16-e3-anderes-image`: der Aufruf nennt `aquasec/trivy:latest` | `p16-e3` |
| 16 E4: Unterbefehl `image` | `m-p16-e4-unterbefehl-entfaellt`: `image` entfällt | `p16-e4` |
| 16 E4: `--input /layout` | `m-p16-e4-input-anderer-pfad`: `--input` nennt einen anderen Pfad | `p16-e4-input` |
| 16 E5: `--format json` | `m-p16-e5-format-entfaellt`: das Argument entfällt | `p16-e5` |
| 16 E5: `--format json` | `m-p16-e5-format-table`: `--format table` | `p16-e5` |
| 16 E6: `--exit-code 0` | `m-p16-e6-exit-code-entfaellt`: das Argument entfällt | `p16-e6` |
| 16 E6: `--exit-code 0` | `m-p16-e6-exit-code-1`: `--exit-code 1` | `p16-e6` |
| 16 E7: `--severity HIGH,CRITICAL` | `m-p16-e7-severity-entfaellt`: das Argument entfällt | `p16-e7` |
| 16 E7: nur `HIGH,CRITICAL` (verengt) | `m-p16-e7-nur-critical`: `--severity CRITICAL` | `p16-e7` |
| 16 E7: nur `HIGH,CRITICAL` (geweitet) | `m-p16-e7-mit-medium`: `--severity HIGH,CRITICAL,MEDIUM` | `p16-e7` |
| 16 E8: `--timeout 10m` | `m-p16-e8-timeout-entfaellt`: das Argument entfällt | `p16-e8` |
| 16 E8: `--timeout 10m` | `m-p16-e8-timeout-60m`: `--timeout 60m` | `p16-e8` |
| 16 E9: kein `--platform` (Trivy ignoriert es bei einem Index mit zwei Manifesten) | `m-p16-e9-platform-kommt-hinzu`: `--platform` kommt hinzu | `p16-e9-ohne-platform` |
| 16 E9: ein Lauf mit dem Verzeichnis von `linux/arm64` | `m-p16-e9-nur-amd64`: nur `linux/amd64` wird gescannt | 20 Fälle, darunter `p16-alle-befunde`, `p16-befund-nur-arm64`, `p16-blobs-inhalt` |
| 16 E9: ein Lauf mit dem Verzeichnis von `linux/amd64` | `m-p16-e9-nur-arm64`: nur `linux/arm64` wird gescannt | 55 Fälle, darunter `p16-alle-befunde`, `p16-auswahl-vor-docker`, `p16-befund-im-zweiten-ergebnis` |
| 16 E9: genau die zwei Plattformen | `m-p16-e9-dritte-plattform`: eine dritte Plattform wird gescannt | 58 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 E10: kein `--skip-db-update` | `m-p16-e10-skip-db-update`: `--skip-db-update` kommt hinzu | `p16-e10-skip-db-update` |
| 16 E10: kein `--skip-files` | `m-p16-e10-skip-files`: `--skip-files` kommt hinzu | `p16-e10-skip-files` |
| 16 E10: kein `--skip-dirs` | `m-p16-e10-skip-dirs`: `--skip-dirs` kommt hinzu | `p16-e10-skip-dirs` |
| 16 E10: kein `--ignorefile` | `m-p16-e10-ignorefile`: `--ignorefile` kommt hinzu | `p16-e10-ignorefile` |
| 16 E10: kein `--ignore-policy` | `m-p16-e10-ignore-policy`: `--ignore-policy` kommt hinzu | `p16-e10-ignore-policy` |
| 16 E10: kein `--vex` | `m-p16-e10-vex`: `--vex` kommt hinzu | `p16-e10-vex` |
| 16 E10: keine Ausnahmedatei im Mount (genau zwei Mounts) | `m-p16-e10-dritter-mount`: ein dritter Mount kommt hinzu | `p16-e10-dritter-mount` |
| 16 E13: `--user <uid>:<gid>` des Aufrufers | `m-p16-user-entfaellt`: `--user` entfällt | `p16-user` |
| 16 E12: `--cache-dir /cache` | `m-p16-cache-dir-entfaellt`: das Argument entfällt | `p16-cache-dir`, `p16-e8` |
| 16 E11: der Zwischenspeicher liegt im Temp-Verzeichnis des Schrittes, nie aus einem früheren Lauf | `m-p16-cache-ueber-laeufe`: der Cache liegt unter `$RUNNER_TEMP/release` | `p16-cache-ueber-laeufe` |
| 16 Schwelle: nur `CRITICAL` und `HIGH` (geweitet um MEDIUM) | `m-p16-schwelle-medium`: `MEDIUM` zählt | `p16-alle-befunde`, `p16-medium-mit-fix`, `p16-summe` |
| 16 Schwelle: alle anderen Schweregrade sperren nicht | `m-p16-schwelle-jede`: jeder Schweregrad zählt | `p16-alle-befunde`, `p16-low-mit-fix`, `p16-medium-mit-fix`, `p16-summe`, `p16-unknown-mit-fix` |
| 16 Schwelle: `HIGH` zählt | `m-p16-schwelle-nur-critical`: nur `CRITICAL` zählt | 10 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 Schwelle: `CRITICAL` zählt | `m-p16-schwelle-nur-high`: nur `HIGH` zählt | 8 Fälle, darunter `p16-alle-befunde`, `p16-critical-mit-fix`, `p16-critical-ohne-fix` |
| 16 Schwelle: Einträge ohne Behebung sperren nicht | `m-p16-ohne-fix-sperrt`: Einträge ohne Behebung zählen als Befund | `p16-alle-befunde`, `p16-critical-ohne-fix`, `p16-high-ohne-fix`, `p16-high-ohne-fix-leer`, `p16-high-ohne-fix-null`, `p16-summe` |
| 16 Behebung: fehlender Schlüssel und `null` sind ohne Behebung | `m-p16-fix-null-gilt-als-fix`: `FixedVersion != ""` (fehlend und `null` gelten als behoben) | `p16-alle-befunde`, `p16-critical-ohne-fix`, `p16-high-ohne-fix`, `p16-high-ohne-fix-null`, `p16-summe` |
| 16 Behebung: die leere Zeichenkette ist ohne Behebung | `m-p16-fix-leer-gilt-als-fix`: `FixedVersion != null` (leer gilt als behoben) | `p16-high-ohne-fix-leer`, `p16-summe` |
| 16 Summe: `<m>` zählt nur Einträge ohne Behebung über der Schwelle | `m-p16-ohne-fix-zaehlt-alle`: `<m>` zählt jeden Eintrag über der Schwelle | `p16-alle-befunde`, `p16-critical-mit-fix`, `p16-high-mit-fix`, `p16-summe` |
| 16 Auswertung: alle Befunde werden genannt (je Plattform) | `m-p16-nur-erster-befund`: nur der erste Befund einer Plattform | `p16-alle-befunde`, `p16-summe` |
| 16 Auswertung: alle Ergebnisse (`Results`) einer Plattform | `m-p16-nur-erstes-ergebnis`: nur `.Results[0]` | `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis` |
| 16 Auswertung: Befunde von `linux/arm64` zählen | `m-p16-arm64-nicht-ausgewertet`: die Befunde von `arm64` werden verworfen | `p16-alle-befunde`, `p16-befund-nur-arm64`, `p16-plattform-genannt` |
| 16 Auswertung: Befunde von `linux/amd64` zählen | `m-p16-amd64-nicht-ausgewertet`: die Befunde von `amd64` werden verworfen | 10 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 Meldung: der Befund nennt die Plattform | `m-p16-plattform-nicht-genannt`: die Plattform fehlt in der Zeile | 7 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 Meldung: die Summenzeile `<n> Befund(e) über der Schwelle, <m> ohne Behebung` | `m-p16-summe-entfaellt`: die Zeile entfällt | `p16-alle-befunde`, `p16-critical-mit-fix`, `p16-high-mit-fix`, `p16-summe` |
| 16 Meldung: `<n>` zählt nur Befunde mit Behebung | `m-p16-summe-zaehlt-unbehebbare-mit`: `<n>` addiert `<m>` | `p16-alle-befunde`, `p16-summe` |
| 16 Meldung: jeder Befund steht als `release: scan: schwachstelle: …` | `m-p16-befunde-nicht-genannt`: die Befundzeilen entfallen | 7 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 Auswertung: aus dem JSON, nicht aus einer Anzeige | `m-p16-auswertung-aus-anzeige`: die Befunde kommen aus einem `grep` über die Ausgabe | 13 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 Gescheitert: ein Exit von Trivy ungleich 0 ist `scan-gescheitert` | `m-p16-exit-gilt-als-gruen`: der Exit wird nicht geprüft | `p16-reihenfolge-exit2-vor-befund`, `p16-reihenfolge-gescheitert-vor-pruefmenge`, `p16-stichwort-scan-gescheitert`, `p16-trivy-exit`, `p16-trivy-exit-125`, `p16-trivy-exit-arm64` |
| 16 Gescheitert: jeder Exit ungleich 0 (nicht nur 1) | `m-p16-exit-nur-1`: nur Exit 1 gilt als gescheitert | `p16-trivy-exit-125` |
| 16 Gescheitert: Exit von Trivy für `linux/arm64` | `m-p16-exit-arm64-ignoriert`: der Exit wird nur für `amd64` geprüft | `p16-reihenfolge-exit2-vor-befund`, `p16-trivy-exit-arm64` |
| 16 Gescheitert: ein leeres Ergebnis ist kein Grün | `m-p16-pruefmenge-leer-gruen`: `length >= 0` | `p16-leeres-ergebnis`, `p16-leeres-ergebnis-arm64`, `p16-reihenfolge-pruefmenge-vor-datenbank`, `p16-stichwort-pruefmenge` |
| 16 Gescheitert: kein JSON oder kein Ergebnis ist `pruefmenge` | `m-p16-pruefmenge-entfernt`: die Prüfung entfällt | 8 Fälle, darunter `p16-ergebnis-kein-array`, `p16-kein-json`, `p16-kein-json-arm64` |
| 16 Gescheitert: Prüfmenge für `linux/arm64` | `m-p16-pruefmenge-arm64-ignoriert`: die Prüfung läuft nur für `amd64` | `p16-kein-json-arm64`, `p16-leeres-ergebnis-arm64` |
| 16 Gescheitert: Datenbank älter als 72 Stunden ist `datenbank` | `m-p16-db-nicht-geprueft`: der Vergleich mit der Grenze entfällt | `p16-db-alt`, `p16-db-alt-arm64`, `p16-db-downloadedat`, `p16-reihenfolge-amd64-vor-arm64`, `p16-reihenfolge-datenbank-vor-plattform`, `p16-stichwort-datenbank` |
| 16 Datenbank: Grenze 72 Stunden (geweitet) | `m-p16-db-7-tage`: die Grenze steht bei 7 Tagen | `p16-db-alt`, `p16-db-alt-arm64`, `p16-db-downloadedat`, `p16-reihenfolge-amd64-vor-arm64`, `p16-reihenfolge-datenbank-vor-plattform`, `p16-stichwort-datenbank` |
| 16 Datenbank: Grenze 72 Stunden (verengt) | `m-p16-db-70-stunden`: die Grenze steht bei 70 Stunden | `p16-db-frisch` |
| 16 Datenbank: Grenze 72 Stunden (geweitet) | `m-p16-db-74-stunden`: die Grenze steht bei 74 Stunden | `p16-db-alt`, `p16-db-alt-arm64`, `p16-db-downloadedat`, `p16-reihenfolge-amd64-vor-arm64`, `p16-reihenfolge-datenbank-vor-plattform`, `p16-stichwort-datenbank` |
| 16 Datenbank: das Alter ist das des Standes (`UpdatedAt`), nicht das der Abholung | `m-p16-db-downloadedat`: `DownloadedAt` | 7 Fälle, darunter `p16-db-alt`, `p16-db-alt-arm64`, `p16-db-downloadedat` |
| 16 Datenbank: das Alter wird nach jedem Lauf gelesen | `m-p16-db-arm64-nicht-geprueft`: nur nach dem Lauf für `amd64` | `p16-db-alt-arm64`, `p16-db-unlesbar-arm64` |
| 16 Datenbank: ein nicht lesbares Alter ist `datenbank` | `m-p16-db-zeitpunkt-ungeprueft`: die Prüfung der Form des Zeitpunkts entfällt | 7 Fälle, darunter `p16-db-unlesbar`, `p16-db-unlesbar-arm64`, `p16-db-unlesbar-jetzt` |
| 16 Datenbank: ein unlesbares Datum behauptet kein Alter | `m-p16-db-datum-ungeprueft`: der Fehler von `date` wird überlesen | `p16-db-unlesbar-kein-datum` |
| 16 Gescheitert: die gescannte Architektur ist die Plattform, sonst `plattform` | `m-p16-plattform-nicht-verglichen`: der Vergleich entfällt | `p16-plattform-abweichend`, `p16-plattform-abweichend-arm64`, `p16-plattform-fehlt`, `p16-stichwort-plattform` |
| 16 Plattform: Vergleich für `linux/arm64` | `m-p16-plattform-arm64-nicht-verglichen`: der Vergleich entfällt für `arm64` | `p16-plattform-abweichend-arm64` |
| 16 Plattform: gleich der Plattform (nicht nur vorhanden) | `m-p16-plattform-nur-nicht-leer`: `[ -n "$arch" ]` | `p16-plattform-abweichend`, `p16-plattform-abweichend-arm64`, `p16-stichwort-plattform` |
| 16 Vorbedingung: `image/digest` (Marke des Schrittes `image`) | `m-p16-image-marke-nicht-verlangt`: die Prüfung entfällt | `p16-ohne-image-mit-layout`, `p16-reihenfolge-digest-vor-layout` |
| 16 Vorbedingung: das Layout liegt da | `m-p16-layout-nicht-verlangt`: die Prüfung entfällt | `p16-ohne-layout` |
| 16 Vorbedingung: ohne Image kein Scan | `m-p16-vorbedingung-beide-entfernt`: beide Prüfungen entfallen | `p16-ohne-image`, `p16-ohne-image-mit-layout`, `p16-ohne-layout`, `p16-reihenfolge-digest-vor-layout`, `p16-reihenfolge-vorbedingung-vor-docker`, `p16-stichwort-vorbedingung` |
| 16 Zustand: `scan.ok` nur bei Grün | `m-p16-marke-bei-rot`: die Marke wird zu Beginn angelegt | 53 Fälle, darunter `p16-alle-befunde`, `p16-auswahl-vor-docker`, `p16-befund-im-zweiten-ergebnis` |
| 16 Zustand: `scan.ok` bei Grün | `m-p16-marke-entfaellt`: die Marke wird nicht angelegt | 14 Fälle, darunter `p16-critical-ohne-fix`, `p16-db-frisch`, `p16-db-zukunft-toleranz` |
| 16 Zustand: ein früherer Scan hinterlässt keine Marke | `m-p16-marke-loescht-entfaellt`: `scan` löscht die Marke nicht zu Beginn | `p16-marke-alt-auswahl`, `p16-marke-alt-entpacken`, `p16-marke-alt-index`, `p16-marke-loescht-entfaellt` |
| 16 Zustand: `eingaben` löscht `scan.ok` | `m-p16-eingaben-loescht-scan-nicht`: die Zeile entfällt | `p16-eingaben-loescht-marke` |
| 16 Zustand: `image` löscht `scan.ok` | `m-p16-image-loescht-scan-nicht`: die Zeile entfällt | `p16-image-loescht-marke` |
| 16 Zustand: `veroeffentlichen` verlangt `scan.ok` (Exit 2 vor `skopeo`) | `m-p16-veroeffentlichen-ohne-scan-marke`: die Prüfung entfällt | `p16-ueber-schwelle-keine-registry`, `p16-veroeffentlichen-ohne-scan` |
| 16 keine Außenwirkung: der Scan ruft `skopeo` nie | `m-p16-scan-ruft-skopeo`: `skopeo inspect` im Schritt | `p16-gruen`, `p16-trockenlauf` |
| 16 Trockenlauf: `skopeo` wird nie gerufen | `m-p16-trockenlauf-ruft-skopeo`: im Trockenlauf ruft der Scan `skopeo` | `p16-trockenlauf` |
| 16 Trockenlauf: der Scan läuft auch bei `workflow_dispatch` | `m-p16-scan-verlangt-push`: der Schritt lehnt ohne Push ab | `p16-trockenlauf` |
| 16 Stelle: `D` und das Layout bleiben unverändert | `m-p16-digest-veraendert`: der Schritt hängt Bytes an das Layout | `p16-d-unveraendert` |
| 16 Stichwort `schwachstelle` | `m-p16-stichwort-schwachstelle`: das Stichwort weicht ab | 7 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 Stichwort `scan-gescheitert` | `m-p16-stichwort-scan-gescheitert`: das Stichwort weicht ab | `p16-reihenfolge-exit2-vor-befund`, `p16-reihenfolge-gescheitert-vor-pruefmenge`, `p16-stichwort-scan-gescheitert`, `p16-trivy-exit`, `p16-trivy-exit-125`, `p16-trivy-exit-arm64` |
| 16 Stichwort `pruefmenge` | `m-p16-stichwort-pruefmenge`: das Stichwort weicht ab | 10 Fälle, darunter `p16-ergebnis-kein-array`, `p16-kein-json`, `p16-kein-json-arm64` |
| 16 Stichwort `datenbank` (zu alt) | `m-p16-stichwort-datenbank-alt`: das Stichwort weicht ab | `p16-db-alt`, `p16-db-alt-arm64`, `p16-db-downloadedat`, `p16-reihenfolge-amd64-vor-arm64`, `p16-reihenfolge-datenbank-vor-plattform`, `p16-stichwort-datenbank` |
| 16 Stichwort `datenbank` (nicht lesbar) | `m-p16-stichwort-datenbank-lesbar`: das Stichwort weicht ab | 7 Fälle, darunter `p16-db-unlesbar`, `p16-db-unlesbar-arm64`, `p16-db-unlesbar-jetzt` |
| 16 Stichwort `plattform` | `m-p16-stichwort-plattform`: das Stichwort weicht ab | `p16-plattform-abweichend`, `p16-plattform-abweichend-arm64`, `p16-plattform-fehlt`, `p16-stichwort-plattform` |
| 16 Exit 2 für einen Scan, der nicht lief | `m-p16-fehler-exit-1`: Exit 1 | 42 Fälle, darunter `p16-auswahl-vor-docker`, `p16-db-alt`, `p16-db-alt-arm64` |
| 16 Exit 1 für einen Befund | `m-p16-befund-exit-2`: Exit 2 | 10 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 13 (B) `veroeffentlichen` verlangt `image/digest`, auch mit `scan.ok` | `m-p13-veroeffentlichen-ohne-image`: `# P13veroeffentlichen`: die Prüfung entfällt | `p11-alte-marke-image-b`, `p11-image-rot-loescht-digest-b`, `p13-veroeffentlichen-ohne-image-bei-scan-ok` |
| 11 (B) `eingaben` löscht `image/`, auch mit `scan.ok` | `m-p11-eingaben-loescht-image-nicht`: `image/` bleibt | `p11-alte-marke-image-b` |
| 11 (B) `image` löscht das alte `image/`, auch mit `scan.ok` | `m-p11-image-loescht-nicht`: `rm -rf` entfällt | `p11-image-rot-loescht-digest-b` |
| 16 Datenbank: `UpdatedAt` mehr als 5 Minuten in der Zukunft ist `datenbank` | `m-p16-db-zukunft-gilt-als-frisch`: die Prüfung entfällt | `p16-db-zukunft` |
| 16 Datenbank: bis 5 Minuten Zukunft sind frisch | `m-p16-db-zukunft-ohne-toleranz`: jede Zukunft ist rot (`-ge 0`) | `p16-db-zukunft-toleranz` |
| 16 Datenbank: Toleranz genau 5 Minuten (verengt/geweitet) | `m-p16-db-zukunft-toleranz-30-minuten`: die Toleranz steht bei 30 Minuten | `p16-db-zukunft` |
| 16 Gescheitert: `Results` ist ein Array, sonst `pruefmenge` (nie grün) | `m-p16-ergebnis-array-ungeprueft`: `(type == "array") and` entfällt | `p16-ergebnis-kein-array` |
| 16 Gescheitert: eine nicht auswertbare Ausgabe (`Vulnerabilities` nicht iterierbar) ist `pruefmenge` | `m-p16-auswertung-fehler-ignoriert`: der Fehler von `jq` wird überlesen | `p16-vulnerabilities-nicht-iterierbar` |
| 16 Meldung: bei Grün steht die Summenzeile als eine Zeile auf stdout | `m-p16-gruen-zeile-fehlt`: die Zeile entfällt | 8 Fälle, darunter `p16-critical-ohne-fix`, `p16-gruen`, `p16-high-ohne-fix` |
| 16 Meldung: bei Grün steht die Summenzeile auf stdout | `m-p16-gruen-zeile-auf-stderr`: die Zeile geht auf stderr | 8 Fälle, darunter `p16-critical-ohne-fix`, `p16-gruen`, `p16-high-ohne-fix` |
| 16 Meldung: erst die Befunde, dann die Summe, auf stderr | `m-p16-summe-vor-befunden`: die Summe steht vor den Befunden | `p16-alle-befunde` |
| 16 E11: zweiter Mount `-v <Cache>:/cache` | `m-p16-cache-mount-entfaellt`: der Mount entfällt | 72 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 E11: der Zwischenspeicher ist schreibbar | `m-p16-cache-mount-ro`: `:ro` | 72 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 Exit 2: genau eine Meldung | `m-p16-meldung-doppelt`: die Meldung steht zweimal | 46 Fälle, darunter `p16-auswahl-vor-docker`, `p16-db-alt`, `p16-db-alt-arm64` |
| 16 Reihenfolge: die Vorbedingung vor dem ersten `docker` | `m-p16-reihenfolge-vorbedingung-nach-docker`: die Prüfungen stehen nach der Schleife | `p16-ohne-image`, `p16-ohne-image-mit-layout`, `p16-ohne-layout`, `p16-reihenfolge-digest-vor-layout`, `p16-reihenfolge-vorbedingung-vor-docker`, `p16-stichwort-vorbedingung` |
| 16 Reihenfolge: `image/digest` vor `image/layout.tar` | `m-p16-reihenfolge-layout-vor-digest`: die beiden Prüfungen sind vertauscht | `p16-reihenfolge-digest-vor-layout` |
| 16 Reihenfolge: `scan-gescheitert` vor `pruefmenge` | `m-p16-reihenfolge-pruefmenge-vor-exit`: die Prüfung der Prüfmenge steht davor | `p16-reihenfolge-gescheitert-vor-pruefmenge`, `p16-trivy-exit-125` |
| 16 Reihenfolge: `pruefmenge` vor `datenbank` | `m-p16-reihenfolge-datenbank-vor-pruefmenge`: die Prüfung der Datenbank steht davor | `p16-db-alt`, `p16-db-alt-arm64`, `p16-db-downloadedat`, `p16-reihenfolge-pruefmenge-vor-datenbank` |
| 16 Reihenfolge: `datenbank` vor `plattform` | `m-p16-reihenfolge-plattform-vor-datenbank`: die Prüfung der Plattform steht davor | `p16-reihenfolge-datenbank-vor-plattform` |
| 16 Reihenfolge: `linux/amd64` vor `linux/arm64` | `m-p16-reihenfolge-arm64-vor-amd64`: die Plattformen sind vertauscht | `p16-auswahl-vor-docker`, `p16-manifest-doppelt`, `p16-manifest-doppelt-variante`, `p16-manifest-fehlt`, `p16-reihenfolge-amd64-vor-arm64` |
| 16 Reihenfolge: ein Exit 2 geht den Befunden vor | `m-p16-reihenfolge-befund-vor-exit2`: ein Befund beendet den Schritt mit Exit 1 vor dem Lauf der zweiten Plattform | `p16-alle-befunde`, `p16-critical-mit-fix`, `p16-high-mit-fix`, `p16-plattform-genannt`, `p16-reihenfolge-exit2-vor-befund`, `p16-summe` |
| Gegenprobe: `jq` in der Werkzeugliste (Exit 2 „jq fehlt“) | `m-p16-werkzeug-jq-entfaellt`: `jq` entfällt aus der Liste | `p16-werkzeug-jq` |
| Gegenprobe: GNU `date` in der Werkzeugliste (Exit 2 „date fehlt“) | `m-p16-werkzeug-date-entfaellt`: `date` entfällt aus der Liste | `p16-werkzeug-date` |
| 16 E10: das Verzeichnis enthält genau `oci-layout`, `index.json`, `blobs` | `m-p16-e10-verzeichnis-layout-tar`: `layout.tar` liegt dabei | `p16-e10-verzeichnis-inhalt`, `p16-verzeichnis-amd64`, `p16-verzeichnis-arm64` |
| 16 E10: das Verzeichnis enthält genau `oci-layout`, `index.json`, `blobs` | `m-p16-e10-verzeichnis-weitere-datei`: eine weitere Datei liegt dabei | `p16-e10-verzeichnis-inhalt`, `p16-verzeichnis-amd64`, `p16-verzeichnis-arm64` |
| 16 Layout je Plattform: `blobs/` im Verzeichnis | `m-p16-verzeichnis-blobs-fehlt`: die Kopie der Blobs entfällt | `p16-blobs-inhalt`, `p16-e10-verzeichnis-inhalt`, `p16-verzeichnis-amd64`, `p16-verzeichnis-arm64` |
| 16 Layout je Plattform: `oci-layout` im Verzeichnis | `m-p16-verzeichnis-oci-layout-fehlt`: die Kopie entfällt | `p16-e10-verzeichnis-inhalt`, `p16-verzeichnis-amd64`, `p16-verzeichnis-arm64` |
| 16 Layout je Plattform: das Manifest der Plattform, nicht das erste | `m-p16-auswahl-erstes-manifest`: die Auswahl nimmt `.manifests[0]` | 68 Fälle, darunter `p16-alle-befunde`, `p16-auswahl-vor-docker`, `p16-befund-im-zweiten-ergebnis` |
| 16 Layout je Plattform: `platform.os` = `linux` | `m-p16-auswahl-os-nicht-gelesen`: `platform.os` wird nicht gelesen | `p16-index-fremdes-os` |
| 16 Layout je Plattform: das Manifest steht nicht mehr als einmal im Index (`plattform`) | `m-p16-auswahl-nicht-auf-genau-eines`: `-ge 1` statt genau eines | `p16-manifest-doppelt`, `p16-manifest-doppelt-variante` |
| 16 Layout je Plattform: das Manifest fehlt nicht im Index (`plattform`) | `m-p16-auswahl-nicht-auf-fehlend`: `-le 1` statt genau eines | `p16-auswahl-vor-docker`, `p16-manifest-fehlt` |
| 16 Layout je Plattform: die `index.json` nennt genau einen Eintrag | `m-p16-index-nennt-beide-manifeste`: sie nennt alle Manifeste des Index | 64 Fälle, darunter `p16-alle-befunde`, `p16-befund-im-zweiten-ergebnis`, `p16-befund-nur-amd64` |
| 16 Reihenfolge: die Auswahl des Manifests vor dem Aufruf von `docker` | `m-p16-auswahl-nach-docker`: die Auswahl steht nach dem Aufruf | 92 Fälle, darunter `p16-alle-befunde`, `p16-auswahl-vor-docker`, `p16-befund-im-zweiten-ergebnis` |
| 16 Layout je Plattform: ein nicht entpackbares Layout ist `scan-gescheitert` | `m-p16-entpacken-fehler-ignoriert`: der Exit von `tar` wird überlesen | `p16-layout-nicht-entpackbar`, `p16-marke-alt-entpacken` |
| 16 Layout je Plattform: ein nicht lesbarer Index oder Index-Blob ist `scan-gescheitert` | `m-p16-index-guard-entfernt`: die Prüfung entfällt | `p16-index-blob-nicht-lesbar`, `p16-index-nicht-lesbar`, `p16-marke-alt-index` |
| 16 Stichwort `scan-gescheitert` (Entpacken) | `m-p16-stichwort-entpacken`: das Stichwort weicht ab | `p16-layout-nicht-entpackbar`, `p16-marke-alt-entpacken` |
| 16 Stichwort `scan-gescheitert` (Index) | `m-p16-stichwort-index`: das Stichwort weicht ab | `p16-index-blob-nicht-lesbar`, `p16-index-nicht-lesbar`, `p16-marke-alt-index` |
| 16 Stichwort `plattform` (Auswahl) | `m-p16-stichwort-auswahl`: das Stichwort weicht ab | `p16-auswahl-vor-docker`, `p16-manifest-doppelt`, `p16-manifest-doppelt-variante`, `p16-manifest-fehlt`, `p16-marke-alt-auswahl` |
| 16 Zustand: `scan` löscht eine alte Marke vor dem ersten Exit 2 (Entpacken) | `m-p16-marke-loescht-nach-entpacken`: das Löschen steht hinter dem Entpacken | `p16-marke-alt-entpacken` |
| 16 Zustand: `scan` löscht eine alte Marke vor dem Lesen des Index | `m-p16-marke-loescht-nach-index`: das Löschen steht hinter dem Lesen des Index (färbt auch den Fall davor, weil dessen Fehlerstelle früher liegt) | `p16-marke-alt-entpacken`, `p16-marke-alt-index` |
| 16 Zustand: `scan` löscht eine alte Marke vor der Auswahl (`plattform`) | `m-p16-marke-loescht-nach-auswahl`: das Löschen steht hinter der Auswahl (färbt auch die beiden Fälle davor) | `p16-marke-alt-auswahl`, `p16-marke-alt-entpacken`, `p16-marke-alt-index` |
| 16 Layout je Plattform: `blobs/` trägt dieselben Blobs wie das Layout | `m-p16-blobs-leer`: `mkdir` statt `cp -al` | `p16-blobs-inhalt` |
| 16 Layout je Plattform: weitere Manifeste im Index (Attestation, anderes Betriebssystem) bleiben unbeachtet | `m-p16-index-mehr-als-zwei-abgelehnt`: ein Index mit mehr als zwei Manifesten wird abgelehnt | `p16-index-attestation`, `p16-index-fremdes-os`, `p16-manifest-doppelt-variante` |
| 16 Layout je Plattform: die Auswahl ignoriert `variant` | `m-p16-auswahl-variante-ausgeschlossen`: ein Eintrag mit `variant` wird ausgeschlossen | `p16-index-variante` |
| 16 Layout je Plattform: zwei Manifeste der Plattform mit verschiedener Variante sind doppelt (`plattform`) | `m-p16-auswahl-gleiche-variante-doppelt`: doppelt ist nur, was dieselbe Variante trägt | `p16-manifest-doppelt-variante` |
| 16 Layout je Plattform: das Layout wird einmal entpackt | `m-p16-entpacken-je-plattform`: es wird je Plattform erneut entpackt | `p16-entpacken-einmal` |

**Grüne Mutanten (Stand der ersten Fassung):** Die Aussage „keiner“ war falsch (Review F-2); sie ist in der Behebungsrunde unten ersetzt.

**Test-Ideen, die kein Fall tragen konnte (Grenze, an die Closure):** (a) dass Trivy die Schwachstellen-Datenbank von seinen Quellen holt und `db/metadata.json` im Cache ablegt, ist nicht gemessen (die Attrappe schreibt die Datei); (b) dass `--platform` mit `--input` einen Index im OCI-Tar auf ein Manifest einengt, ist nicht gemessen; (c) dass `Metadata.ImageConfig.architecture` im JSON steht, ist nicht gemessen. Alle drei zeigt der erste Trockenlauf, und jede Abweichung endet rot (`pruefmenge`, `datenbank`, `plattform`), nie grün.

**Rückgaben an den Architect** (*Randform · Frage*; `.claude/commands/implement-slice.md`, Randform-Rückgabe): Der Auftrag des Koordinators sah vor, offene Randformen im Bericht zurückzugeben und weiterzuarbeiten. Der Code trug deshalb eine **vorläufige Lesart** (Spalte *Im Code*), je an einer Stelle. **Entschieden** vom Architect in `SPEC-067` Punkt 11, 13 und 16 (Behebungsrunde unten); der Vertragstext ist nicht vom Implementer geändert.

| Randform | Im Code (vorläufig) | Frage |
|---|---|---|
| Woher kommt das Alter der Datenbank (Punkt 16 nennt keine Quelle; ob der JSON-Bericht von Trivy es trägt, ist nicht gemessen)? | `db_alter` liest `UpdatedAt` aus `db/metadata.json` eines Cache-Verzeichnisses, das der Schritt im Temp-Verzeichnis anlegt (zweiter Mount `-v <Cache>:/cache`, `--cache-dir /cache`, dazu `--user <uid>:<gid>` des Aufrufers, damit die Dateien des Containers dem Runner gehören und das Aufräumen gelingt); `UpdatedAt` (Stand der Daten), nicht `DownloadedAt` | Bestätigt? Punkt 16 E10 nennt nur den Mount des Layouts; der Cache-Mount, `--cache-dir` und `--user` sind drei Elemente mehr (je ein Fall und eine Mutation: `p16-user`, `p16-cache-dir`, `p16-cache-kein-zwischenspeicher`, `p16-e10-mount`). Alternative: `trivy version --format json` im selben Cache |
| Die Marke `image.ok` aus Punkt 16 gibt es nicht (Schritt `image` hinterlässt `image/digest` und `image/layout.tar`, Punkt 13 nennt keine `image.ok`) | `scan` verlangt beide (Exit 2, `release: vorbedingung: …` wie die übrigen Schritte); je Wächter ein Fall (`p16-ohne-image-mit-layout`, `p16-ohne-layout`) | Den Namen im Vertrag auf `image/digest` ändern, oder soll `image` eine Marke `image.ok` anlegen? |
| Form der Meldung bei Exit 2 | `release: scan: <Stichwort>: …` für `scan-gescheitert`, `pruefmenge`, `datenbank`, `plattform`; `release: vorbedingung: …` für die Vorbedingung | Bindend so? |
| Reihenfolge mehrerer Gründe (`AGENTS.md` §3.12): Punkt 16 ordnet sie nicht | je Plattform nacheinander (`amd64`, `arm64`), je Plattform `scan-gescheitert`, `pruefmenge`, `datenbank`, `plattform` (Reihenfolge der Aufzählung); ein Exit 2 einer späteren Plattform geht dem Exit 1 (Befunde der früheren) vor, weil die Befunde erst nach beiden Läufen genannt werden; **kein Fall** zur Reihenfolge | Binden (dann je Paar ein Fall) oder offen lassen? |
| `scan.ok` und Zustand früherer Läufe | `eingaben` und `image` löschen `scan.ok`, `scan` löscht es zu Beginn (je ein Fall und eine Mutation) | Punkt 11 und 13 sollen es nennen |
| Summenzeile bei Grün | `release: scan: 0 Befund(e) über der Schwelle, <m> ohne Behebung` auf stdout (Checkliste Punkt 7 verlangt „Zeile aus dem Log“) | Punkt 16 nennt die Summe nur bei Exit 1 |
| Tabelle in Punkt 11 ohne Zeile `scan` | `releasing.md` §11 führt die Zeile (nicht gelaufen: wiederholen; Befund: beheben, neuer Tag) | Die Tabelle des Vertrags ergänzen, sonst „deckt sich“ nicht |
| Ein Image ohne jedes Paket hätte kein `Results` | `pruefmenge` (Exit 2): eine leere Prüfmenge ist kein Grün; ob `distroless/static` Pakete im Bericht führt, ist nicht gemessen | Ist das gemeint? |
| Trivy nicht gemessen: `--platform` zusammen mit `--input` auf einem OCI-Tar mit Index, `Metadata.ImageConfig.architecture`, Betrieb als Nicht-root mit `--cache-dir`, Form von `db/metadata.json`, `Results` bei sauberem Image | wie oben beschrieben; fail-closed | Der erste Trockenlauf misst; wer entscheidet danach, wenn eine Annahme falsch war (Vertrag nachziehen und neuer Lauf)? |

**Entwurf für Handbuch und `README.md` (nicht geliefert; erst nach dem Lauf des Tags, `SPEC-067` Punkt 18):** Der Lauf liefert hier weder die Zeile in §11 noch den Kopf noch §2 *Installation*; sie hängen an `VERSION` und dem Release-Commit, den der Koordinator festlegt. Wortlaut als Entwurf; `<…>` ist der Wert, den der Lauf liefert, ein Befehl steht erst im Handbuch, wenn er nach dem Tag einmal gelaufen ist (Beleg: Befehl, Lauf, Ausgabe in §7 des Commits, der ihn einträgt).

- §11 *Änderungshistorie* (vor dem Tag, im Release-Commit; ersetzt „Es gibt keine veröffentlichte Version.“ durch die Tabelle mit der Kopfzeile `| Version | Datum | Änderung |`): `| 0.1.0 | <Datum des Tags> | Erste veröffentlichte Version: Aufzeichnen (record), Wiedergeben (replay) und Einspielen (play) von PostgreSQL-Verkehr; Binaries für Linux, macOS und Windows (amd64, arm64) und ein Image für linux/amd64 und linux/arm64. |` (drei Zellen, kein weiteres `|`; der Text ist die Notiz des GitHub-Release).
- Kopf (nach dem Lauf): die Zeile *Software-Version* nennt `0.1.0` mit dem Zusatz „Ausgabe von `pgwire-recorder version` des veröffentlichten Binaries“; *Gültigkeitsbereich* nennt das Binary in der Version 0.1.0 (Ausgabe von `pgwire-recorder version`: `pgwire-recorder 0.1.0`) und sagt, dass ein aus dem Repository gebautes Binary `pgwire-recorder dev` meldet. Beleg: `version` des veröffentlichten Binaries.
- §2 *Installation*, neue Abschnitte (nach dem Lauf, nur belegte Befehle): *Das Binary vom GitHub-Release* (`curl -fsSLO …/releases/download/v<Version>/pgwire-recorder_<Version>_<os>_<arch>` und `…/SHA256SUMS`, `sha256sum -c --ignore-missing SHA256SUMS`, die Plattformen; `darwin` und `windows` sind gebaut, nicht ausgeführt) und *Das Image aus einer Registry* (`docker pull ghcr.io/pt9912/pgwire-recorder:<Version>` und `docker.io/…`, `docker run --rm … version`, `:latest`, weil der Lauf es gesetzt und der Smoke es belegt hat). Beleg je Befehl: Lauf nach dem Tag.
- `README.md`: derselbe Bezug in zwei Sätzen, ohne Zielstand.

### Behebungsrunde 1 (Review `docs/reviews/2026-10-19-review-slice-erster-release-freigabe.md`, Entscheidungen des Architect)

Stand: Kopf des Codes `f7986fd7` auf `impl-fr` (Merge von `arch-fr` mit `SPEC-067` `5a1e2171`, Report `0a2de07c`, Code `4fd605f8`, `releasing.md` `598d1ff8`). Läufe nacheinander am Endstand:

| Lauf | Ergebnis | Dauer |
|---|---|---|
| `make release-gegenprobe` | `release-gegenprobe: grün (511 Läufe)` | 1419 s |
| `make workflows-check` | Exit 0, `workflows-check: 2 Datei(en) geprüft, kein Befund` | 2 s |
| `make docs-check` | Exit 0, `d-check: 1246 Datei(en) geprüft, 0 Befund(e)` | 2 s |
| `make gates` | Exit 0 | 296 s |

**Messung der Größe** (`AGENTS.md` §3.14): `git diff --numstat 760c51d4..f7986fd7 -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'` ergibt als Summe aus Einfügungen und Löschungen **617** (`release-gegenprobe.sh` 390, `releasing.md` 133, `release.sh` 72, `SPEC-067` 18 vom Architect, `release.yml` 2, `harness/README.md` 2); ohne den Vertragstext 599. Marke 1200 (51 %) und Meldeschwelle 900 sind nicht erreicht, eine Meldung war nicht nötig.

**Rückgaben der ersten Runde: entschieden** vom Architect in `SPEC-067` Punkte 11, 13 und 16 (`dddf7cac`, `8968ab86`, `5a1e2171`); die Spalte *Im Code* der Tabelle oben gilt damit als Entscheidung, nicht mehr als vorläufige Lesart. Ausnahmen, die der Code und die Gegenprobe nachziehen mussten: das Alter in der Zukunft (neu), die sieben Fälle der Reihenfolge, die Fälle `p16-stichwort-*`, `p16-cache-*`, `p16-db-downloadedat` und `p16-e10-dritter-mount` (Namen aus dem Vertrag; die Fälle `p16-cache-kein-zwischenspeicher` und `p16-cache-ueber-laeufe` sind jetzt zwei, `p16-marke-loescht-entfaellt` heißt der frühere `p16-alte-marke`, `p16-e10-mount` heißt `p16-e10-dritter-mount`). Offen bleibt nur, was der Vertrag selbst als nicht gemessen nennt (Punkt 16, *Nicht gemessene Annahmen* (a) bis (f)); die Messung ist der erste Trockenlauf.

| Review | Behebung |
|---|---|
| F-1 (MEDIUM) Randformen im Code | durch die Entscheidungen des Architect aufgelöst; Code geändert nur an `db_alter` (Zukunft) |
| F-2 (MEDIUM) drei Wächter ohne Fall | die beiden `jq`-Aufrufe der Auswertung sind einer (ein Wächter statt zweier, die sich gegenseitig deckten); je ein Fall und eine Mutation für `(type == "array") and` (`p16-ergebnis-kein-array`, `m-p16-ergebnis-array-ungeprueft`) und für den Wächter der Auswertung (`p16-vulnerabilities-nicht-iterierbar`, `m-p16-auswertung-fehler-ignoriert`); die Aussage „Grüne Mutanten: keiner“ der ersten Runde war falsch und ist ersetzt (unten) |
| F-3 (LOW) Werkzeugliste | `jq` und `date` in der Liste der Gegenprobe (Exit 2, „jq fehlt“); Fälle `p16-werkzeug-jq`, `p16-werkzeug-date`, Mutationen `m-p16-werkzeug-jq-entfaellt`, `m-p16-werkzeug-date-entfaellt` (Tabelle oben) |
| F-4 (LOW) Fehlertabelle | `releasing.md` §11 Zeile `scan`: bei `plattform`, `pruefmenge` oder gleichbleibender Meldung Korrektur und neuer, höherer Tag (Punkt 18 (b)); nur Netz, Zwischenspeicher und Datenbank: wiederholen |
| F-5 (INFO) Zukunft | `db_alter`: Alter unter −300 s ist `datenbank`; Fälle `p16-db-zukunft`, `p16-db-zukunft-toleranz`; Mutationen: Zukunft gilt als frisch, Toleranz fehlt, Toleranz 30 Minuten |
| F-6 (INFO) Name `pgr-gp` | `releasing.md` §10 Punkt 2: „Ablage des Gegenproben-Containers“ |
| F-7 (b) | §1 und DoD-Punkt 1 an den Stand angeglichen: der Handbuch- und README-Teil entsteht nach dem Tag (Punkt 18); Entwurf unten |

**Geprüft gegen den Vertragstext** (Punkt 16 nach `5a1e2171`): Die Summenzeile bei Grün steht auf stdout und ist genau eine Zeile: `p16-gruen` und die Fälle `p16-medium-mit-fix`, `-low-mit-fix`, `-unknown-mit-fix`, `p16-high-ohne-fix` (mit `<m>`), `-leer`, `-null`, `p16-critical-ohne-fix` prüfen stdout gleich der Zeile; Mutationen `m-p16-gruen-zeile-fehlt` und `m-p16-gruen-zeile-auf-stderr` (neu). Erst die Befunde, dann die Summe auf stderr: `p16-alle-befunde`, Mutation `m-p16-summe-vor-befunden` (neu). `p16-ohne-image` fing die Reihenfolge Digest vor Layout nicht (beide Wächter hätten denselben Exit): `p16-reihenfolge-digest-vor-layout` (neu, beide fehlen, die Meldung nennt den Digest). Bei Exit 2 genau eine Meldung (ohne Zeilen `release: trivy: `) und keine Befunde: Prüfung in jedem Fall mit Exit 2, Mutation `m-p16-meldung-doppelt` (neu). Die Namen `p16-cache-*`, `p16-user`, `p16-e10-dritter-mount`, `p16-eingaben-loescht-marke`, `p16-image-loescht-marke`, `p16-marke-loescht-entfaellt` stimmen mit dem Vertrag überein; `scan` löscht die Marke nach seinen Vorbedingungen (Zeile hinter `# P16layout`).

**Grüne Mutanten (Stand Runde 1; berichtigt in Runde 3):** keiner, und kein Mutant gilt als äquivalent. Wo zwei Wächter dieselbe Eingabe fingen, ist einer entfernt (die beiden `jq`-Aufrufe der Auswertung) oder jeder hat einen Fall, in dem nur er greift (Digest und Layout; `eingaben`, `image` und `scan` löschen `scan.ok`, jeder mit eigenem Fall).

### Behebungsrunde 2 (erster Trockenlauf, `SPEC-067` Punkt 16 nach `077da231`)

**Erster roter Trockenlauf:** GitHub-Lauf `39263455592` (2026-10-19). **Abweichung:** Annahme (a) des Vertrags war falsch. Trivy liest ein OCI-Layout nur als **Verzeichnis**, nicht als Tar, und bei einem Index mit zwei Manifesten liest es das erste und ignoriert `--platform` (am Host mit dem gepinnten Image gemessen: die Architektur im Bericht ist beidemal `amd64`). Der Scan war damit rot, bevor er eine Plattform bewertete. **Gemessen** sind jetzt (a), (b), (c), (d) und (f) (Vertrag, Punkt 16, *Annahmen über Trivy*); **offen** bleiben (e) (ein Image ohne Befund liefert `Results` mit einem Ziel; für ein Testimage auf `busybox` fehlte der Schlüssel) und (g) (ein Verzeichnis mit einer `index.json`, die nur ein Manifest nennt, wird gelesen und der Bericht trägt dessen Architektur) bis zum nächsten Trockenlauf; die Attrappen der Gegenprobe belegen beide nicht.

**Umsetzung** (Stand: Kopf des Codes `271ce087` auf `impl-fr2`, Basis `origin/main` `b307b020` mit `SPEC-067` `077da231`): `release.sh` entpackt `image/layout.tar` einmal in ein frisches Verzeichnis im Temp-Verzeichnis des Schrittes (Fehler: `scan-gescheitert`), liest `index.json` und den Index-Blob mit `jq` (nicht lesbar: `scan-gescheitert`) und legt je Plattform ein Verzeichnis an: `oci-layout`, `blobs/` (Hardlinks) und eine `index.json` mit genau dem unveränderten Deskriptor des Manifests (`os` = `linux`, passende `architecture`); fehlt er oder steht er mehrfach, endet der Schritt mit Exit 2 `plattform` vor dem Aufruf von `docker`. Der Aufruf mountet `-v <Verzeichnis>:/layout:ro`, nennt `image --input /layout` und kein `--platform`; die Architekturprüfung im Bericht bleibt. Das Layout aus `image` und `D` bleiben unberührt (`p16-d-unveraendert`). Die Attrappe von `docker` zeichnet beim Aufruf den Inhalt des Mount-Verzeichnisses auf.

**Neue Elemente: Zusage · Mutation · roter Fall** (die Zeilen stehen in der Mutationstabelle oben; die Fälle sind `p16-verzeichnis-amd64`, `p16-verzeichnis-arm64`, `p16-e10-verzeichnis-inhalt`, `p16-e9-ohne-platform`, `p16-manifest-fehlt`, `p16-manifest-doppelt`, `p16-index-attestation` (grün), `p16-index-fremdes-os` (grün, `windows/amd64` neben `linux/amd64`; Grenzfall zu `platform.os`), `p16-layout-nicht-entpackbar` (ein Tar, der bis auf ein abgeschnittenes letztes Mitglied vollständig ist: der Wächter von `tar` ist der einzige, der greift), `p16-index-nicht-lesbar`, `p16-index-blob-nicht-lesbar`, `p16-auswahl-vor-docker`; geändert `p16-e2`, `p16-e2-quelle`, `p16-e4-input`, `p16-e9-amd64`, `p16-e9-arm64`): Auswahl nimmt das erste Manifest; Auswahl prüft nicht auf genau eines (`-ge 1`) und nicht auf fehlend (`-le 1`); `index.json` nennt beide Manifeste; `blobs` bzw. `oci-layout` fehlt; `platform.os` nicht gelesen; `layout.tar` oder eine weitere Datei im Mount; `--platform` kommt hinzu; Auswahl nach dem Aufruf von `docker`; Prüfung von `tar` und von Index und Blob entfernt; drei Stichwörter der neuen Meldungen. Jede Mutation färbt genau die in der Tabelle genannten Fälle (`GP_ENTDECKE=1`); es gibt keinen grünen Mutanten (Stand dieser Runde; berichtigt in Runde 3, F-1 und F-6). Wo zwei Prüfungen dieselbe Eingabe fingen (Index und Blob, Entpacken und Index), ist eine entfernt (Index und Blob sind eine Prüfung, `p16-index-nicht-lesbar` und `p16-index-blob-nicht-lesbar` färben dieselbe Mutation) oder der Fall trifft nur die eine (`p16-layout-nicht-entpackbar`).

**Messung der Größe** (`AGENTS.md` §3.14, Basis `760c51d4`, bis `271ce087`): Summe aus Einfügungen und Löschungen von `git diff --numstat 760c51d4..271ce087 -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'`: **692** (`release-gegenprobe.sh` 440, `releasing.md` 136, `release.sh` 93, `SPEC-067` 19 vom Architect, `release.yml` 2, `harness/README.md` 2); ohne den Vertragstext 673. Marke 1200 (58 %) und Meldeschwelle 900 sind nicht erreicht.

| Lauf | Ergebnis | Dauer |
|---|---|---|
| `make release-gegenprobe` | `release-gegenprobe: grün (526 Läufe)` | 1780 s |
| `make workflows-check` | Exit 0, kein Befund | 1 s |
| `make docs-check` | Exit 0, `d-check: 1247 Datei(en) geprüft, 0 Befund(e)` | 2 s |
| `make gates` | Exit 0 | 280 s |

**Offene Rückgaben:** keine neue Randform. Offen sind nur (e) und (g) des Vertrags bis zum nächsten Trockenlauf ((g) ist seit `4d8fc29a` gemessen, siehe Runde 3).

### Behebungsrunde 3 (Review Runde 2, `docs/reviews/2026-10-19-review-slice-erster-release-freigabe-r2.md`; `SPEC-067` Punkt 16 nach `4d8fc29a`)

Stand: Kopf des Codes `c97ed3bd` auf `impl-fr2` (Merge von `arch-fr4`, `4d8fc29a`: die Auswahl vergleicht nur `os` und `architecture`, `variant` wird ignoriert, zwei `linux/arm64`-Manifeste mit verschiedener Variante sind doppelt; Annahme (g) ist gemessen, offen bleibt nur (e)).

| Review | Behebung |
|---|---|
| F-1 (MEDIUM) `m-p16-cache-ueber-laeufe` färbte keinen Fall | `p16-cache-ueber-laeufe` unterscheidet den Mount `:/cache` vom Mount `:/layout:ro` (Muster ` -v <Temp-Verzeichnis des Schrittes>/…:/cache `); die Mutation färbt wieder genau diesen Fall; die Aussage „keine grünen Mutanten“ der Runden 1 und 2 galt damit zu Unrecht und ist hier auf den Stand dieser Runde berichtigt |
| F-2 (LOW) alte Marke an den neuen Exit-2-Stellen | `p16-marke-alt-entpacken`, `-index`, `-auswahl` (alte `scan.ok` vorab, Erwartung: Marke weg); Mutationen: das Löschen steht hinter dem Entpacken (färbt nur den Fall des Entpackens), hinter dem Lesen des Index (färbt zusätzlich den Fall davor), hinter der Auswahl (zusätzlich beide davor): die Fehlerstellen liegen hintereinander, ein Löschen hinter einer Stelle lässt die Marke an allen früheren stehen, deshalb sind die Mengen ineinander enthalten |
| F-3 (LOW) Inhalt von `blobs/` | Die Attrappe schreibt als vierte Zeile je Aufruf die Summen der Dateien von `blobs/`; `p16-blobs-inhalt` vergleicht sie mit denen des entpackten Layouts; Mutation `mkdir` statt `cp -al` (`m-p16-blobs-leer`) |
| F-4 (LOW) Mutation zu `p16-index-attestation` | `m-p16-index-mehr-als-zwei-abgelehnt` (färbt `p16-index-attestation`, `p16-index-fremdes-os` und `p16-manifest-doppelt-variante`, die drei Einträge hat); Zeile in der Tabelle |
| F-5 (INFO) Variante | in §6 als Randform genannt (unten), entschieden vom Architect (`4d8fc29a`); Fälle `p16-index-variante` (grün) und `p16-manifest-doppelt-variante` (Exit 2 `plattform`, nur der Lauf für `amd64`); Mutationen: ein Eintrag mit `variant` wird ausgeschlossen (`m-p16-auswahl-variante-ausgeschlossen`), doppelt ist nur, was dieselbe Variante trägt (`m-p16-auswahl-gleiche-variante-doppelt`, färbt anders: nur den zweiten Fall) |
| F-6 (INFO) zwei beobachtbar äquivalente Mutanten | (1) Entpacken je Plattform: kein äquivalenter Mutant mehr, `p16-entpacken-einmal` zählt die Aufrufe von `tar -xf` an einer Attrappe, Mutation `m-p16-entpacken-je-plattform`; (2) `manifests[-1]` statt `manifests[0]` im äußeren Index ist äquivalent: **Invariante** das Layout trägt genau einen Eintrag im äußeren Index (Punkt 5); **Ort** der Schritt `image` (`# P5ein-index`, Fall `p5-layout-form`, Mutation `m-p5-ein-index-entfernt`), ein Layout mit zwei Einträgen erreicht `scan` nie; **Sonde** der Lauf der Gruppe `f` mit der Mutation `.manifests[-1].digest` an einer Wegwerfkopie der Gegenprobe (`GP_ENTDECKE=1`, Mutation `m-tmp-aequiv`): alle Fälle grün, die Mutation färbt nichts; trennen würde sie nur ein Layout mit zwei Einträgen und verschiedenen Indizes, und das hält `p5-layout-form` (Variante `zwei-eintraege`) in `image` auf, grün im Lauf unten. Der Ort ist Code, die Invariante ist nicht die Zusage dieses Slice |
| F-7 (INFO) Hardlinks | Annahme: Quelle (`$scratch/layout`) und Ziel (`$scratch/layout-<Architektur>`) liegen im selben Temp-Verzeichnis, also im selben Dateisystem; ohne Hardlinks endet der Lauf rot (Trivy liest das Verzeichnis nicht, `scan-gescheitert`), nie grün; kein eigener Fall |
| F-8 (INFO) Fehlertabelle | `releasing.md` §11 Zeile `scan`: `scan-gescheitert` für das Layout (nicht entpackbar, Index nicht lesbar) ist ein Fehler des Schrittes `image`: den Lauf ab `image` wiederholen; bleibt es, Korrektur und neuer Tag (Punkt 18 (b)) |

**Randform Variante des Manifests** (`AGENTS.md` §3.12, Entscheidung des Architect, `SPEC-067` Punkt 16): Trägt das Manifest einer Plattform ein `variant` (`arm64/v8`), wird es wie `arm64` gewählt; zwei Manifeste derselben Plattform mit verschiedener Variante sind doppelt (`plattform`). Das Produkt-Image trägt keine Variante; die Randform schützt gegen einen Index, dessen Builder sie setzt.

**Neue Zeilen der Mutationstabelle** stehen oben (Zeilen zu `m-p16-marke-loescht-nach-*`, `m-p16-blobs-leer`, `m-p16-index-mehr-als-zwei-abgelehnt`, `m-p16-auswahl-variante-ausgeschlossen`, `m-p16-auswahl-gleiche-variante-doppelt`, `m-p16-entpacken-je-plattform`).

**Messung der Größe** (`AGENTS.md` §3.14, Basis `760c51d4`, bis `c97ed3bd`): Summe aus Einfügungen und Löschungen von `git diff --numstat 760c51d4..c97ed3bd -- . ':(exclude)docs/plan' ':(exclude)docs/reviews'`: **723** (`release-gegenprobe.sh` 471, `releasing.md` 136, `release.sh` 93, `SPEC-067` 19 vom Architect, `release.yml` 2, `harness/README.md` 2); ohne den Vertragstext 704. Marke 1200 (60 %) und Meldeschwelle 900 sind nicht erreicht.

| Lauf | Ergebnis | Dauer |
|---|---|---|
| `make release-gegenprobe` | `release-gegenprobe: grün (534 Läufe)` | 2269 s (davon rund 5 Minuten Aufräumen der Temp-Verzeichnisse) |
| `make workflows-check` | Exit 0, kein Befund | 1 s |
| `make docs-check` | Exit 0, `d-check: 1248 Datei(en) geprüft, 0 Befund(e)` | 2 s |
| `make gates` | Exit 0 | 272 s |

**Offene Rückgaben:** keine. Offen ist nur (e) bis zum nächsten Trockenlauf.

**Frage an den Koordinator:** Wie entsteht der Release-Commit? Der Plan lässt es offen: `VERSION` (heute `0.1.0-dev`) und die Zeile der Änderungshistorie stehen im Commit, der getaggt wird (`SPEC-067` Punkt 17), und `make gates`, der Gegenproben-Container und `make gegenproben-host` laufen auf genau diesem Hash. Der Implementer legt ihn nicht an: `VERSION` bleibt `0.1.0-dev`, die Zeile fehlt. Ob der Commit nach Review und Verifikation auf dem Hauptzweig entsteht (dann ist die Closure dieses Slice nach dem Tag) und wer ihn schreibt, ist zu entscheiden; ebenso das Datum der Zeile.

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
