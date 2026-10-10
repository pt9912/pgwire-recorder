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

**Ziel:** Ein Tag `v<SemVer>` wird gegen die Versionsdatei geprüft, baut die Binaries für die Zielplattformen reproduzierbar mit `SHA256SUMS` und das Docker/OCI-Image, veröffentlicht das Image mit gleichem Index-Digest in `ghcr.io/pt9912/pgwire-recorder` und `docker.io/pt9912/pgwire-recorder` und lässt sich nach einem Fehler wieder anstoßen; belegt mit einem Probe-Tag.

**Übernommen aus `welle-erster-release`** (Entscheidungen des Nutzers vom 2026-10-09 nach dem
Vergleich von `docs/maintainer/releasing.md` mit den Release-Dokumenten von a-check, d-check,
pg-change-feed und u-boot; `docs/maintainer/releasing.md` nennt diesen Slice dort, wo er
liefert):

1. `:latest` zeigt nur nach einem stabilen Tag auf das neue Image; eine Vorabversion ändert
   `:latest` nicht.
2. Der Spiegel auf Docker Hub ist fail-closed: Beide Registries tragen denselben Index-Digest
   der Manifestliste, sonst ist der Release rot.
6. Zusätzlich zum Tag eine Versionsdatei (etwa `docs/user/version.md`), die die Pipeline gegen
   den Tag prüft.
7. Die Prüfung des Tags auf die Form `v<SemVer>` vor jedem Push.
8. Wiederanlauf und Fehlertabelle: `docs/maintainer/releasing.md` §11 nennt je Schritt der
   Pipeline, was ein Fehler dort hinterlässt und wie der Lauf wieder angestoßen wird
   (Entscheidung des Nutzers vom 2026-10-09 beim Schnitt unten).

Die Randformen dieser Punkte sind offen und gehen vor dem Code an den Architect (§6).

**Abgegeben** an `slice-erster-release-freigabe` (dort §1, *Übernimmt*, mit der Kennung dieses
Slice), nach dem Schnitt vom 2026-10-09 vor dem Start (Entscheidung des Nutzers, §4): die
Entscheidungen 3 (Freigabe-Checkliste), 4 (Freshness-Audit) und 5 (Image-Scan vor dem Push)
und der erste echte Tag. Dieser Slice belegt die Pipeline mit einem Probe-Tag.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die Homebrew-Formel und ihr Nachweis — `slice-erster-release-homebrew-nachweis`.
- Image-Scan, Freshness-Audit, Freigabe-Checkliste und der erste echte Tag —
  `slice-erster-release-freigabe` (oben, *Abgegeben*); er setzt diese Pipeline voraus, weil der
  Scan in ihr vor dem Push läuft.
- Die Wahl der ersten Versionsnummer — wird beim Release festgelegt; die Versionsdatei
  (Punkt 6) trägt sie, wählt sie aber nicht.
- Reproduzierbarkeit des Images — `LH-QA-01` verlangt sie nur für das Binary; das Image wird
  über seinen Digest identifiziert (Punkt 2).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung): Ein Probe-Tag veröffentlicht das Image mit `linux/amd64` und `linux/arm64` als eine Manifestliste in beiden Registries, und es besteht den Smoke aus Abnahmeszenario 9. Beide Registries tragen denselben Index-Digest, sonst endet der Lauf rot (Punkt 2); `:latest` wird nur bei einem stabilen Tag gesetzt, der Probe-Tag lässt es unverändert (Punkt 1). Beleg in §7: Lauf-Kennung und Digest des Probe-Tags (Test). Benutzerhandbuch und `README.md` beschreiben, was dieser Slice liefert, im Ist-Zustand des gebauten Binaries: ohne Chronik, ohne Zielstand, im Handbuch ohne Verweis auf Spezifikation, ADRs, Slices oder Reviews (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung). Aus `slice-doku-ist-stand` (dort §6, *Installation*): Der Slice liefert im Handbuch §2 *Installation* die Abschnitte zu den Release-Binaries (Plattformen, `SHA256SUMS`) und zu den Images aus `ghcr.io` und `docker.io` sowie die Plattform-Zeile in §1 *Voraussetzungen*, jeweils nur, soweit der Probe-Tag sie belegt; ein Abschnitt, den erst ein echter Tag belegt, bleibt aus.
- [ ] [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-01`](../../../../spec/lastenheft.md#lh-qa-01--determinismus), [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew): Die Binaries für Linux, macOS und Windows (`amd64`, `arm64`) hängen mit einer Datei `SHA256SUMS` am GitHub-Release des Probe-Tags (`LH-FA-19.a`); zwei Builds desselben Tags liefern dieselben Summen. Die Release-Notes entstehen aus der neuen Zeile der Änderungshistorie des Handbuchs; ein `CHANGELOG.md` wird nicht geführt.
- [ ] Vor jedem Push prüft die Pipeline, dass der Tag die Form `v<SemVer>` hat und der Version in der Versionsdatei entspricht; sonst endet sie, bevor etwas gebaut oder veröffentlicht wird (Punkte 6 und 7, Test mit falschem Tag und abweichender Datei). `docs/maintainer/releasing.md` §11 trägt die Fehlertabelle je Schritt der Pipeline und den Wiederanlauf; ein roter Probe-Lauf ist nach ihr wieder angestoßen (Punkt 8, Beleg in §7).
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
| Release-Automatisierung | neu | Prüfung von Tag und Versionsdatei, Build, Veröffentlichung mit Vergleich der Index-Digests, `:latest` nur bei stabilem Tag, Prüfsummen; der Scan vor dem Push kommt mit `slice-erster-release-freigabe` dazu |
| Versionsdatei (Ort offen, §6) | neu | Version der Software, gegen den Tag geprüft |
| `docs/maintainer/releasing.md` §3, §6, §8, §9, §11 | update | Versionsquelle, was der Tag auslöst, Kontrollen nach dem Push, Release-Notes, Fehlertabelle und Wiederanlauf |
| `docs/user/benutzerhandbuch.md`, `README.md` | update | Ist-Zustand des gelieferten Verhaltens (`AGENTS.md` §3.11, seit slice-v1-abschluss-einspielen-laufsteuerung) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-v1-abschluss` ist `done`, und
`slice-harness-abdeckung-gate` und `slice-harness-coverage` liegen in `done/` (Start-Trigger von
welle-erster-release). Schritt 1 der Reihenfolge in §5 von
[welle-erster-release](../welle-erster-release.md). Vor dem ersten Code-Commit entscheidet der
Architect die Randformen aus §6 (`AGENTS.md` §3.12).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: die Release-Automatisierung verlangt mehr als drei Liefer-Punkte, oder der Diff ist nicht in einer Review-Sitzung prüfbar — zurück zur Zerlegung. Schnitt dann: *Tag und Binaries* (DoD-Punkte 2 und 3) und *Image* (DoD-Punkt 1).

**Schnitt vor dem Start, 2026-10-09:** Mit den sieben Entscheidungen vom 2026-10-09 trug §2 sechs
Liefer-Punkte. Der Nutzer entschied am 2026-10-09 den Schnitt nach dem Vorschlag in §6: Die
Pipeline bleibt hier, dazu Wiederanlauf und Fehlertabelle; Scan, Audit, Checkliste und der erste
echte Tag gehen an `slice-erster-release-freigabe`. Der Slice lag in `next/` und hat nicht
begonnen; ein Übergang fällt nicht an.
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
- **Größe.** Mit den sieben Entscheidungen vom 2026-10-09 trug §2 sechs Liefer-Punkte, die
  Grenze ist drei. Entscheidung des Nutzers vom 2026-10-09: geschnitten nach dem Vorschlag des
  Planner (§4) — **Ausgang:** eingetreten, Folge-Slice `slice-erster-release-freigabe`.
- **Probe-Tag in öffentlichen Registries.** Der Beleg mit einem Probe-Tag legt ein Image einer
  Vorabversion in `ghcr.io` und `docker.io` ab, das kein Release ist; ob und wie es dort bleibt,
  ist eine Randform (unten) — **Ausgang:** offen bis Closure.
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
- *Versionsdatei (Punkt 6)* — Ort (`docs/user/version.md` oder anders) und Form (eine Zeile,
  SemVer ohne `v`?); ob das Binary die Version trägt und `pgwire-recorder version` sie aus der
  Datei oder dem Tag nimmt; Vorabversion in der Datei; Wirkung auf die Reproduzierbarkeit des
  Binaries (`LH-QA-01`).
- *SemVer-Prüfung (Punkt 7)* — „vor jedem Push“: vor dem Push eines Images in der Pipeline
  oder auch vor `git push` des Tags; Build-Metadaten (`+…`); führende Nullen; Form des
  Probe-Tags (`-probe.<N>`) als Vorabversion.
- *Wiederanlauf und Fehlertabelle (Punkt 8)* — die Schritte der Pipeline, in denen ein Fehler
  etwas hinterlässt (GitHub-Release angelegt, eine Registry beschrieben, die andere nicht);
  ob ein Wiederanlauf denselben Tag neu fährt oder einen höheren verlangt (Rollback-Regel:
  Tags werden nicht verändert); wer ihn anstößt; was mit einem halb veröffentlichten
  Probe-Tag geschieht.
- *Probe-Tag in den Registries* — ob das Image eines Probe-Tags in `ghcr.io` und `docker.io`
  bleibt, gelöscht wird oder nur in einen eigenen Namensraum geht; ob der GitHub-Release des
  Probe-Tags als Vorabversion stehen bleibt.

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

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
