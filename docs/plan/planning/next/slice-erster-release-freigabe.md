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

**Berührte Spec-Stellen:** `LH-FA-16.a` · `SPEC-031`

**Verantwortlich:** pt9912

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

**Übernimmt:** `slice-erster-release-veroeffentlichung` — nach dem Schnitt vom 2026-10-09 vor
dessen Start (Entscheidung des Nutzers; dort §1, *Abgegeben*, und §4) die Entscheidungen des
Nutzers vom 2026-10-09 (Geber `welle-erster-release`), die jener aus der Welle übernommen hatte:

3. Eine Freigabe-Checkliste in `docs/maintainer/releasing.md` mit einem Beleg je Punkt,
   Anti-Punkten (was ein Release nicht tun darf) und einer Incident-Klausel; der ausgefüllte
   Eintrag des ersten Releases steht in der Closure von welle-erster-release (dort §3).
4. Vor dem Tag ein Freshness-Audit der vendored Baseline und des gepinnten d-check.
5. Ein Scan des Images vor dem Push in eine Registry.

Dazu der erste echte Tag: Jener Slice belegt die Pipeline mit einem Probe-Tag, dieser setzt
den ersten stabilen Tag mit Checkliste und Scan. Die Randformen der drei Punkte zogen aus §6
jenes Slice hierher (§6 unten).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Prüfung von Tag und Versionsdatei, Build, `SHA256SUMS`, Veröffentlichung mit gleichem
  Index-Digest, `:latest`, Fehlertabelle und Wiederanlauf — `slice-erster-release-veroeffentlichung`;
  er liegt vor diesem Slice, und der Scan setzt sich in dessen Pipeline vor den Push.
- Die Homebrew-Formel im Tap und ihr Nachweis — `slice-erster-release-homebrew-nachweis`; er
  setzt den ersten stabilen Tag dieses Slice voraus.
- Ein Sensor für das Freshness-Audit — `harness/conventions.md` führt den Audit als Handlung,
  nicht als Sensor, weil er Netz über `bash`, `git` und `docker` hinaus braucht; Bestand bleibt
  bewusst stehen.
- Die Wahl der ersten Versionsnummer — wird beim Release festgelegt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung): Die Pipeline scannt das Image jeder Plattform der Manifestliste vor dem Push; ein Befund über der Schwelle beendet den Lauf, bevor eine Registry etwas erhält (Punkt 5, Test mit einem Image über der Schwelle; Beleg in §7: Zusage · Mutation · roter Test, `AGENTS.md` §3.10).
- [ ] `docs/maintainer/releasing.md` §5 beschreibt das Freshness-Audit der vendored Baseline und des gepinnten d-check vor dem Tag mit der Form seines Belegs; das Audit vor dem ersten echten Tag ist mit diesem Beleg ausgeführt (Punkt 4).
- [ ] [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew): `docs/maintainer/releasing.md` §10 trägt die Freigabe-Checkliste mit Beleg je Punkt, Anti-Punkten und Incident-Klausel (Punkt 3). Der erste echte, stabile Tag ist nach ihr freigegeben: Binaries mit `SHA256SUMS` am Release, Image in beiden Registries mit gleichem Index-Digest und `:latest`, Smoke aus Abnahmeszenario 9 bestanden. Der ausgefüllte Eintrag steht in der Closure von welle-erster-release.
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
| Release-Automatisierung | update | Scan des Images vor dem Push, mit Schwelle |
| `docs/maintainer/releasing.md` §5, §10 | update | Freshness-Audit vor dem Tag, Freigabe-Checkliste |
| Test der Pipeline | neu | Image über der Schwelle wird nicht veröffentlicht |
| `docs/plan/planning/welle-erster-release-results.md` | neu (bei Welle-Closure) | ausgefüllter Eintrag der Checkliste des ersten Releases |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-erster-release-veroeffentlichung` liegt in `done/`
(Pipeline mit Probe-Tag belegt). Schritt 2 der Reihenfolge in §5 von
[welle-erster-release](../welle-erster-release.md). Vor dem ersten Code-Commit entscheidet der
Architect die Randformen aus §6 (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Der Scan verlangt einen Umbau der
  Pipeline, der den Diff über eine Review-Sitzung hebt; Schnitt dann: *Scan* (DoD-Punkt 1)
  und *Freigabe mit erstem Tag* (DoD-Punkte 2 und 3).
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

**Randformen** (`AGENTS.md` §3.12) — **alle offen für den Architect**, zu entscheiden vor dem
ersten Code-Commit in der Spezifikation oder einer ADR; sie zogen mit dem Schnitt vom
2026-10-09 aus §6 von `slice-erster-release-veroeffentlichung` hierher:

- *Freigabe-Checkliste (Punkt 3)* — Form des Belegs je Punkt (Lauf-Kennung, Digest, Link);
  welche Anti-Punkte; Inhalt der Incident-Klausel (Auslöser, wer, wo festgehalten); wo der
  ausgefüllte Eintrag eines späteren Releases ohne Welle steht.
- *Freshness-Audit (Punkt 4)* — was als frisch gilt; ob ein neuerer Stand der Baseline oder
  von d-check den Tag sperrt oder nur dokumentiert wird; Beleg des Audits. Der Audit braucht
  Netz; `harness/conventions.md` führt ihn als Handlung, nicht als Sensor.
- *Image-Scan (Punkt 5)* — Scanner und sein Pin; Schwelle (Schweregrad, nur behebbare
  Befunde?); Ausnahmen und ihr Ort (`AGENTS.md` §3.2, Suppression-Verbot); Scan je Plattform
  der Manifestliste; Netz für die Schwachstellen-Datenbank; ob auch der Probe-Tag gescannt
  wird.

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

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
