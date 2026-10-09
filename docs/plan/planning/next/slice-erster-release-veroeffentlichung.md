# Slice slice-erster-release-veroeffentlichung: Veröffentlichung von Binaries und Images

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-erster-release.

**Bezug:** [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)

**Berührte Spec-Stellen:** `LH-FA-16.a` · `LH-FA-19.a` · `SPEC-031` · `SPEC-035`

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Tag `v<SemVer>` baut die Binaries für die Zielplattformen reproduzierbar und das Docker/OCI-Image und veröffentlicht das Image in `ghcr.io/pt9912/pgwire-recorder` und `docker.io/pt9912/pgwire-recorder`.

**Übernommen aus `welle-erster-release`** (Entscheidungen des Nutzers vom 2026-10-09 nach dem
Vergleich von `docs/maintainer/releasing.md` mit den Release-Dokumenten von a-check, d-check,
pg-change-feed und u-boot; `docs/maintainer/releasing.md` nennt diesen Slice dort, wo er
liefert):

1. `:latest` zeigt nur nach einem stabilen Tag auf das neue Image; eine Vorabversion ändert
   `:latest` nicht.
2. Der Spiegel auf Docker Hub ist fail-closed: Beide Registries tragen denselben Index-Digest
   der Manifestliste, sonst ist der Release rot.
3. Eine Freigabe-Checkliste in `docs/maintainer/releasing.md` mit einem Beleg je Punkt,
   Anti-Punkten (was ein Release nicht tun darf) und einer Incident-Klausel; der ausgefüllte
   Eintrag des ersten Releases steht in der Closure von welle-erster-release.
4. Vor dem Tag ein Freshness-Audit der vendored Baseline und des gepinnten d-check.
5. Ein Scan des Images vor dem Push in eine Registry.
6. Zusätzlich zum Tag eine Versionsdatei (etwa `docs/user/version.md`), die die Pipeline gegen
   den Tag prüft.
7. Die Prüfung des Tags auf die Form `v<SemVer>` vor jedem Push.

Die Randformen dieser Punkte sind offen und gehen vor dem Code an den Architect (§6). Mit
ihnen ist der Slice zu groß (§2, §4); der Schnitt ist vorgeschlagen und nicht entschieden (§6).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Homebrew-Formel und ihr Nachweis — `slice-erster-release-homebrew-nachweis`.
- Die Wahl der ersten Versionsnummer — wird beim Release festgelegt; die Versionsdatei
  (Punkt 6) trägt sie, wählt sie aber nicht.
- Reproduzierbarkeit des Images — `LH-QA-01` verlangt sie nur für das Binary; das Image wird
  über seinen Digest identifiziert (Punkt 2).


## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung): Der erste echte Tag veröffentlicht das Image mit `linux/amd64` und `linux/arm64` als eine Manifestliste in beiden Registries, und es besteht den Smoke aus Abnahmeszenario 9. Beide Registries tragen denselben Index-Digest, sonst endet der Lauf rot (Punkt 2); `:latest` wird nur bei einem stabilen Tag gesetzt, ein Probe-Tag lässt es unverändert (Punkt 1, Test mit Probe-Tag).
- [ ] Die Release-Notes des GitHub-Releases entstehen aus der neuen Zeile der Änderungshistorie des Handbuchs; ein `CHANGELOG.md` wird nicht geführt.
- [ ] [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew): Die Binaries für Linux, macOS und Windows (`amd64`, `arm64`) hängen mit einer Datei `SHA256SUMS` am Release (`LH-FA-19.a`); zwei Builds desselben Tags liefern dieselben Summen.
- [ ] Vor jedem Push prüft die Pipeline, dass der Tag die Form `v<SemVer>` hat und der Version in der Versionsdatei entspricht; sonst endet sie, bevor etwas gebaut oder veröffentlicht wird (Punkte 6 und 7, Test mit falschem Tag und abweichender Datei).
- [ ] Das Image wird vor dem Push gescannt; ein Befund über der Schwelle beendet den Lauf, bevor eine Registry etwas erhält (Punkt 5, Test mit einem Image über der Schwelle).
- [ ] `docs/maintainer/releasing.md` §10 trägt die Freigabe-Checkliste mit Beleg je Punkt, Anti-Punkten und Incident-Klausel; einer ihrer Punkte ist das Freshness-Audit der Baseline und von d-check vor dem Tag (Punkte 3 und 4). Der ausgefüllte Eintrag des ersten Releases steht in der Closure von welle-erster-release.
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
| Release-Verfahren (`docs/maintainer/releasing.md`) | update | vom beschriebenen auf das ausgeführte Verfahren bringen |
| Release-Automatisierung | neu | Prüfung von Tag und Versionsdatei, Build, Scan, Veröffentlichung mit Vergleich der Index-Digests, `:latest` nur bei stabilem Tag, Prüfsummen |
| Versionsdatei (Ort offen, §6) | neu | Version der Software, gegen den Tag geprüft |
| `docs/maintainer/releasing.md` §5, §6, §8, §10, §11 | update | Prüfungen vor dem Tag, Freigabe-Checkliste, Kontrollen nach dem Push, Wiederanlauf |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-v1-abschluss` ist `done`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: die Release-Automatisierung verlangt mehr als zwei Liefer-Punkte — zurück zur Zerlegung. Mit den Entscheidungen vom 2026-10-09 trägt §2 sechs Liefer-Punkte; der Slice ist schon vor dem Start zu groß und wird vor `next` → `in-progress` geschnitten (Vorschlag in §6, Entscheidung des Nutzers offen).
- `in-progress` → `open`: Zugangsdaten für die Registries fehlen — extern klären.


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

- Die Zugangsdaten für `ghcr.io` und `docker.io` sind externe Betreiber-Handlungen — **Ausgang:** offen bis Closure.
- **Größe.** Mit den sieben Entscheidungen vom 2026-10-09 trägt §2 sechs Liefer-Punkte, die
  Grenze ist drei. *Vorschlag des Planner vom 2026-10-09*, nicht angelegt: Dieser Slice behält
  die Pipeline — Tag und Versionsdatei (Punkte 6, 7), Binaries mit `SHA256SUMS` samt
  GitHub-Release mit Release-Notes, Image in beiden Registries mit gleichem Index-Digest und
  `:latest` nur bei stabilem Tag (Punkte 1, 2) —, nachgewiesen mit einem Probe-Tag statt dem
  ersten echten Tag. Ein neuer Slice (etwa `slice-erster-release-freigabe`) nimmt Scan (5),
  Freshness-Audit (4) und Freigabe-Checkliste (3) und setzt mit ihr den ersten echten Tag; er
  setzt diesen voraus, weil der Scan in dessen Pipeline läuft, und
  `slice-erster-release-homebrew-nachweis` setzt ihn voraus. **Entscheidung des Nutzers
  offen** — **Ausgang:** offen bis Closure.
- **`:latest` im Beispiel der Spezifikation.** Das Beispiel in `LH-FA-16.a` (Spezifikation,
  Technik-Stratum) nennt `pgwire-recorder:latest`, ohne Registry. Mit Punkt 1 gibt es
  `:latest` erst nach dem ersten stabilen Tag, und ohne Registry zieht das Beispiel kein
  veröffentlichtes Image. Das Lastenheft (`LH-FA-16`) sagt nichts über Tags; ein
  Change Request am Lastenheft ist danach nicht nötig. *Für den Architect:* ob das Beispiel
  geändert wird (etwa `ghcr.io/pt9912/pgwire-recorder:<version>`) und ob die Regel zu `:latest`
  als Zusage an Anwender in die Spezifikation gehört (`LH-FA-16.a` oder `SPEC-031`) — vor dem
  Code (`AGENTS.md` §3.12). Der Planner ändert keine Spezifikation — **Ausgang:** offen bis
  Closure.

**Randformen** (`AGENTS.md` §3.12) — **alle offen für den Architect**, zu entscheiden vor dem
ersten Code-Commit in der Spezifikation oder einer ADR:

- *`:latest` (Punkt 1)* — Vorabversion und Probe-Tag (unverändert, so entschieden); stabiler
  Tag mit kleinerer Version als der, auf den `:latest` zeigt (Patch eines älteren Zweigs);
  Zeitpunkt des Setzens, wenn eine Registry den Push schon hat und die andere scheitert.
- *Fail-closed (Punkt 2)* — womit der Index-Digest gelesen wird; ob der Spiegel dieselbe
  Manifestliste kopiert oder neu baut; was mit dem Versions-Tag geschieht, der in einer
  Registry schon liegt, wenn der Vergleich rot ist (Rollback-Regel: Tags werden nicht
  verändert); Wiederanlauf eines roten Laufs.
- *Freigabe-Checkliste (Punkt 3)* — Form des Belegs je Punkt (Lauf-Kennung, Digest, Link);
  welche Anti-Punkte; Inhalt der Incident-Klausel (Auslöser, wer, wo festgehalten); wo der
  ausgefüllte Eintrag eines späteren Releases ohne Welle steht.
- *Freshness-Audit (Punkt 4)* — was als frisch gilt; ob ein neuerer Stand der Baseline oder
  von d-check den Tag sperrt oder nur dokumentiert wird; Beleg des Audits. Der Audit braucht
  Netz; `harness/conventions.md` führt ihn als Handlung, nicht als Sensor.
- *Image-Scan (Punkt 5)* — Scanner und sein Pin; Schwelle (Schweregrad, nur behebbare
  Befunde?); Ausnahmen und ihr Ort (`AGENTS.md` §3.2, Suppression-Verbot); Scan je Plattform
  der Manifestliste; Netz für die Schwachstellen-Datenbank.
- *Versionsdatei (Punkt 6)* — Ort (`docs/user/version.md` oder anders) und Form (eine Zeile,
  SemVer ohne `v`?); ob das Binary die Version trägt und `pgwire-recorder version` sie aus der
  Datei oder dem Tag nimmt; Vorabversion in der Datei; Wirkung auf die Reproduzierbarkeit des
  Binaries (`LH-QA-01`).
- *SemVer-Prüfung (Punkt 7)* — „vor jedem Push“: vor dem Push eines Images in der Pipeline
  oder auch vor `git push` des Tags; Build-Metadaten (`+…`); führende Nullen; Form des
  Probe-Tags (`-probe.<N>`) als Vorabversion.

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

- `BEO-REPO/slice-waechst-durch-uebernahmen` (2×) — der Slice nimmt sieben Entscheidungen auf.
  Kein Beleg: Die Übernahmen stehen als eigene Liefer-Punkte in §2, die Zählung zeigt die
  Größe (sechs), und §6 schlägt den Schnitt vor; das Muster ist das Wachsen *ohne* mehr Punkte.
- `BEO-REPO/spec-randform-erst-im-review-entschieden` (16×, verkörpert in `AGENTS.md` §3.12) —
  die Randformen der neuen Punkte stehen offen in §6 und gehen vor dem Code an den Architect.
- `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (17×, `AGENTS.md` §3.10) — die
  Prüfungen der Pipeline (Tag, Versionsdatei, Scan, Digest-Vergleich) sind neue Verträge; je
  Zusage eine Mutation.

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
