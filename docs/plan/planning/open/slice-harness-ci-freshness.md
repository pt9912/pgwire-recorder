# Slice slice-harness-ci-freshness: Der Nachtlauf gegen upstream — `upstream-drift.yml` und das Audit der Baseline und von d-check

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-erster-release, vierte Stelle der Reihenfolge, direkt hinter `slice-harness-ci-gegenproben` und vor `slice-harness-ci-freshness-images` und `slice-erster-release-veroeffentlichung` (Auftrag des Nutzers vom 2026-10-18; Schnitt und Reihenfolge darin: Empfehlung des Planners, siehe `slice-harness-ci-workflows` §6 Frage 1). Start-Trigger weicht von der Welle ab: `slice-harness-ci` in `done/` (die Konventionen der Workflows und `make workflows-check` stehen dort).

**Bezug:** Kein Anspruch des Produkts (`LH-*`): der Slice macht den Alterungs-Audit der gepinnten Fremdbestände zu einem Lauf mit Netz außerhalb der Gates. Aktive ADRs: [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Docker-only; Werkzeuge laufen in gepinnten Images). Die ADR dieses Slice (Nachtlauf gegen upstream: wann rot, wann Meldung) schreibt der Architect vor dem Code (Entscheidung des Nutzers).

**Berührte Spec-Stellen:** `SPEC-066` (neu; Nummer vergibt der Architect, nächste freie; Vertrag in `harness/vertraege/`, vor dem ersten Code-Commit; `slice-harness-ci-freshness-images` ergänzt ihn)  · `SPEC-060` (gelesen: Namens- und Zuordnungsregel der Listen und Ausnahmeliste) · `SPEC-063` und `SPEC-064` (gelesen: Gate der Workflows, Rechte und Pins der CI) · Baseline-Regelwerk `modul-02-harness-bootstrap.md` §Freshness-Audit der vendored Baseline (gelesen: die sieben Eigenschaften des Audits) · `harness/conventions.md` §Der mitgelieferte Baum altert still (wird nachgezogen, Frage 3).

**Verantwortlich:** —
**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `.github/workflows/upstream-drift.yml` (nächtlich per `schedule`, von Hand per `workflow_dispatch`) fragt zwei Make-Ziele, die **nicht** in `make gates` stehen, ob upstream einen neueren Stand führt als der Pin dieses Repos: die vendored **Baseline** (die Release-*Liste* des Kurs-Repos gegen den Tag im Verzeichnisnamen unter `.harness/baseline/`, die geschuldete Handlung aus `harness/conventions.md`, hier der Ort mit Netz) und **d-check** (die Release-Liste von `pt9912/d-check` gegen den Tag in `d-check.mk`). Ein neuerer Stand ist ein roter Schritt, ein Netzfehler ein benanntes Überspringen mit Zählung (Randform 4), jede Achse ein eigener Schritt mit `if: always()`, damit eine rote die andere nicht verdeckt. Das Skript `tools/harness/pin-freshness.sh` (Name: Architect) kennt in diesem Slice eine Art der Frage — „gibt es einen neueren Tag in der Release-Liste von `<Besitzer>/<Repo>`“ — und holt sie mit `curl` **in einem per Digest gepinnten Image**, nicht auf dem Host (`AGENTS.md` §3.1; Vorbild `tools/harness/traeger-fetch.sh`). Die Gegenprobe zeigt Fälle und Mutationen ohne Netz.

**Das Verhältnis zum Audit der Baseline.** Der Audit hat nach dem Regelwerk Eigenschaften, die dieser Slice einlöst: ein beobachtbarer Auslöser (der Lauf, nicht ein Kalender), Netz außerhalb der Gates, die Release-Liste statt des Assets. Den Durchgang durch die Adaptions-Liste (`harness/conventions.md` §Aktive Adaptionen) macht weiter ein Mensch beim Baseline-Sprung; ein roter Lauf ist der Anlass, nicht der Durchgang.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Aktualität der Images (Tag und Digest in `Dockerfile` und `harness/mk/*.mk`: golang, distroless, golangci-lint, a-check, d-migrate, postgres, curl) — `slice-harness-ci-freshness-images` (Folge-Slice, direkt dahinter, setzt diesen voraus und erweitert das Skript um zwei Arten der Frage); hier genügt eine Art, und die Größe bliebe sonst über der Marke.
- Aktualität der Action-SHA und der Werkzeug-Versionen im Workflow — kein Slice nimmt sie an (`slice-harness-ci-workflows` §6 Frage 3); die Form der Pins hält `make workflows-check`.
- Der Vergleich des Assets des gepinnten Tags (Hash gegen `SHA256SUMS`, „Tag upstream verschoben“) — anderer Vorgang: Die Integrität der Arbeitskopie hält `make baseline-verify` lokal; die Frage „wurde der gepinnte Tag nachträglich neu gebaut“ ist eine zweite Messung mit eigenem Exit (Frage 4).
- Ein Baseline-Sprung, der Durchgang durch die Adaptionen und das Nachziehen verkörperter Regeln — Handlung eines Menschen nach dem Befund, ein eigener Slice je Sprung; dieser Slice liefert den Anlass.
- Ein Issue oder eine Nachricht bei Rot — `issues: write` wäre ein Recht mehr und ein zweiter Mechanismus; Rot in der Actions-Übersicht trägt (Randform 5).
- Eine Aufnahme in `make gates` oder in die Kette — der Lauf braucht Netz, und ein Upstream-Ausfall dürfte die Gates nicht färben; `harness/conventions.md` begründet die Trennung, der Slice zieht dort nur den Satz nach, dass der Sensor außerhalb der Gates mitkommt (Frage 3).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Skript und Ziele):** `tools/harness/pin-freshness.sh` mit den Zielen `make baseline-freshness` und `make d-check-freshness` (Namen: `SPEC-066`) liest den Pin aus seiner Quelle, fragt die Release-Liste in einem gepinnten Image, vergleicht gleich/ungleich nach der Normalisierung des Vertrags und endet nach `SPEC-066` mit je einem Exit für aktuell, veraltet und übersprungen; keines der Ziele steht in `GATE_CHECKS`, beide stehen in der Ausnahmeliste der schnellen Kette, wo die Namensregel von `SPEC-060` es verlangt; jede Entscheidung je Randform aus §6 steht vor dem ersten Code-Commit in `SPEC-066`.
- [ ] **Liefer-Punkt 2 (Workflow):** `.github/workflows/upstream-drift.yml` mit `schedule` und `workflow_dispatch`, `permissions: {}` oben, `contents: read` im Job, Zeitgrenze, Pins nach `SPEC-063`, ein Schritt je Achse mit `if: always()`; besteht `make workflows-check` und `make docs-check`. Belegt in §7 mit der Nummer eines grünen Laufs auf GitHub, der Dauer und dem Lauf eines Probe-Pins (ein absichtlich veralteter Pin in einem Probe-Zweig färbt genau seinen Schritt rot).
- [ ] **Liefer-Punkt 3 (Gegenprobe und Index):** `make pin-freshness-gegenprobe` (Name: Architect) zeigt am Skript mit einer Attrappe von `docker` und von `curl` ohne Netz je Zusage von `SPEC-066` und je Element ihrer Aufzählungen (Pin-Quellen, Formen des Tags, Antworten der Quelle, Exit-Arten) einen Fall und eine Mutation, die nur ihn färbt (`AGENTS.md` §3.10), an Kopien von `upstream-drift.yml` die Zusagen der Datei, dazu `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` je ein Fall; das Ziel steht in `GEGENPROBEN_CONTAINER`. `harness/README.md` §Sensors führt die drei Zeilen (zwei Ziele als Werkzeuge ohne Gate-Bindung, Gegenprobe) mit Vertrag und Grenze; `harness/conventions.md` trägt den Satz zum Sensor nach (Frage 3).
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `harness/` und `tools/` ändern) und `make gates-schnell-gegenprobe` grün; `make pin-freshness-gegenprobe` grün auf dem Endstand.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/pin-freshness.sh`, `harness/mk/freshness.mk` (Namen: Architect) | neu | Skript mit einer Art der Frage (Release-Liste); zwei Ziele, kein Gate |
| `.github/workflows/upstream-drift.yml` | neu | Nachtlauf, ein Schritt je Achse; Vorbild d-check `upstream-drift.yml` |
| Gegenprobe (`harness/mk/…`, `tools/harness/…`; Namen: Architect) | neu | Fälle und Mutationen an Skript und Datei; Eintrag in `GEGENPROBEN_CONTAINER` |
| `harness/vertraege/SPEC-066-…md` | neu (Architect) | Vertrag; vor dem Code |
| `harness/README.md` §Sensors, `harness/conventions.md` | update | Zeilen; der Satz zum Sensor außerhalb der Gates |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-ci` liegt in `done/`, dazu der Start-Trigger von [welle-erster-release](../welle-erster-release.md), und der Architect hat `SPEC-066` und die ADR geschrieben, jede Randform aus §6 entschieden (`AGENTS.md` §3.12). Schritt 4 der Reihenfolge in §5 der Welle.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: mehr als drei Liefer-Punkte, oder der Diff ist nicht in einer Review-Sitzung prüfbar — zurück zur Zerlegung. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- .github tools harness Makefile`, Tests und Gegenprobe eingeschlossen, liegt über 900; Basis: beim Anspruch (HEAD vor dem Anspruch, hier einzutragen). Schätzung des Planners in dieser Metrik (2026-10-18, ohne Basis, nach dem Maß des Skripts `pin-freshness.sh` von d-check mit 193 Zeilen für drei Arten): Skript mit einer Art und Fragment 110 bis 160, Workflow 90 bis 120, Gegenprobe das Doppelte des Produktcodes (220 bis 320), `SPEC-066` 110 bis 160, `harness/README.md` und `harness/conventions.md` 25 bis 40, zusammen 555 bis 800. Schnitt, falls die Marke eintritt: *Baseline* und *d-check* (zwei Ziele) — ein Ziel je Slice.
- `in-progress` → `open`: das gepinnte Image für `curl` ist nicht zu beschaffen, oder der Nutzer entscheidet in Frage 1 gegen jedes Überspringen bei Netzfehlern.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, ein Lauf von Hand (`workflow_dispatch`) auf GitHub für den Endstand mit beiden Achsen grün und ein Lauf mit Probe-Pin rot als Beleg in §7, `make gates` grün, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Ein Wächter, den man wegen Fremdstörungen abschaltet, ist schlechter als keiner; ein Wächter, der bei Netzausfall grün bleibt, behauptet Deckung.** Das Überspringen ist fail-open und nennt sich (Randform 4); ein Netz, das nie antwortet, wäre dauerhaft grün. — **Ausgang (vorgesehen):** weiter offen → neues `BEO-REPO/nachtlauf-uebersprungen-sieht-gruen-aus` bei der Closure; die Zusammenfassung zählt die übersprungenen Teile, der Vertrag nennt die Grenze.
- **Ein Nachtlauf, in den niemand sieht, ist ein verwaister Sensor; ein dauerhaftes Rot verliert das Signal für die nächste Achse; `schedule` wird nach 60 Tagen ohne Aktivität abgeschaltet.** — **Ausgang (vorgesehen):** weiter offen → `BEO-REPO/nachtlauf-ohne-leser` (angelegt bei der Closure von `slice-harness-ci-gegenproben` oder hier, je nachdem, welcher Slice zuerst schließt; ein Beleg je Vorgang).
- **`releases/latest` kann eine ältere Version nennen als der Pin** (Patch auf einem alten Zweig, nachträglich veröffentlicht): „ungleich“ ist dann nicht „veraltet“. — **Ausgang (vorgesehen):** entfallen, wenn der Vertrag das Wort wählt, das die Messung trägt („anders“ statt „veraltet“) und die Gegenprobe den Fall zeigt (Randform 3).
- **Ein neuer Baseline-Tag ist kein Auftrag.** Der Durchgang durch die Adaptionen kostet; ein Dauerrot ohne Handlung wiederholt das Problem. — **Ausgang (vorgesehen):** weiter offen → dieselbe Beobachtung wie Risiko 2; das Verfahren nach dem Befund ist Handlung eines Menschen und steht nicht in diesem Slice.

**Randformen** (`AGENTS.md` §3.12) — **alle offen für den Architect**, zu entscheiden vor dem ersten Code-Commit in `SPEC-066` (Technik-Stratum, `harness/vertraege/`), mit ADR für Entscheidung und Gründe. Vorbereitet vom Planner am 2026-10-18 (Empfehlungen, kein Beschluss):

1. **Quellen der Pins.** Baseline: der Verzeichnisname unter `.harness/baseline/` (genau eines, keines, mehrere; Form `v<N>.<M>.<P>`; Abgleich mit der Stand-Zeile in `harness/conventions.md`). d-check: `DCHECK_IMAGE` in `d-check.mk` (Zeile mit `?=`, Form `…:v<SemVer>`, Zeile fehlt, mehrere, Tag ohne `v`). Empfehlung: jede Abweichung der Form ist ein Aufruffehler (Exit 2), kein Überspringen; ein Fall je Form.
2. **Quelle der Frage.** Weiterleitung von `…/releases/latest` (kein API, kein Token, kein `jq`) auf `…/releases/tag/<x>`; Endet die Weiterleitung nicht auf `/tag/…` (Repo ohne Release), ist das ein eigener Fall. Empfehlung nach dem Muster von d-check; Fälle: Weiterleitung auf Tag, ohne Tag, 404, 403, 429, 5xx, Zeitüberschreitung, leere Antwort.
3. **Vergleich und Wort.** Gleich/ungleich nach Normalisierung (führendes `v` auf beiden Seiten gestrichen); keine Reihenfolge von Versionen. „Ungleich“ heißt „anders“, nicht „veraltet“ (Risiko 3). Vorabversionen: `latest` übergeht sie; ein Pin auf eine Vorabversion ist gleich/ungleich zum Tag der letzten Hauptversion. Ein Fall je Form.
4. **Was bricht den Lauf und was meldet.** Rot (Exit ungleich 0): der Tag ist anders; die Quelle antwortet 404 (der Besitzer oder das Repo ist weg). Überspringen (Exit 0 mit `SKIP: <Teil>: <Grund>` und Zählung in der Zusammenfassung): Netzfehler, Zeitüberschreitung, Rate-Limit (403/429), 5xx, leere oder unlesbare Antwort. Empfehlung nach d-check: fail-open für Fremdstörung; jeder Fall eine Zeile. Ein Fehler des Aufrufs (Pin nicht lesbar, fehlendes Werkzeug) ist Exit 2 und rot.
5. **Meldung.** Rot in der Actions-Übersicht (kein Issue, kein `issues: write`, keine Geheimnisse); Benachrichtigung ist Einstellung auf GitHub. Der Vertrag nennt die Grenze (Dauerrot; Abschaltung von `schedule` nach 60 Tagen ohne Aktivität).
6. **Rate-Limit und Netz.** Kein Token, kein API-Aufruf (die Obergrenze für anonyme API-Aufrufe der geteilten Adresse trifft nur diese); `curl` mit Obergrenze für Zeit und Weiterleitungen (`--max-time`, `--max-redirs`), kein automatisches Wiederholen. Nicht erreichbares Image für `curl` (Pull-Fehler) ist ein Aufruffehler, nicht ein Überspringen, weil der Sensor dann nie gefragt hat.
7. **Docker-only.** `curl` läuft in einem Image mit Digest; der Pin des Images hat eine Stelle (der Pin von `tools/harness/traeger-fetch.sh` oder ein gemeinsamer). Empfehlung: ein gemeinsamer Pin, eine Gegenprobe hält beide Stellen gegeneinander, sonst driftet ein Zweitbestand. Netzzugriff: das Image bekommt Netz (anders als die netzlosen Gates); der Mount ist schreibgeschützt oder entfällt.
8. **Aufrufarten des Treibers** (`AGENTS.md` §3.12): `-n`, `-q`, `-t`, `-i`, `MAKEFLAGS` aus der Umgebung; die Ziele schreiben keinen Stand, das Rezept trägt kein `$(MAKE)`; `-i` verschluckt einen Befund (kein Fehler des Skripts ändert sich); Empfehlung: `-i` ist im Workflow nicht gesetzt, ein Fall hält es. Ein Fall je Flag.
9. **Workflow.** Eine Datei, ein Job, ein Schritt je Achse mit `if: always()` (der erste ohne); Reihenfolge: Baseline, d-check; `concurrency` mit Gruppe, `cancel-in-progress: false`; Uhrzeit abseits der vollen Stunde (versetzt zu `gegenproben.yml`); Zeitgrenze 15 Minuten (Vorbild); `workflow_dispatch` ohne Eingaben.
10. **Gleichzeitige Läufe und `-j`.** Die Ziele sind nicht Teil von `GATE_CHECKS`; zwei gleichzeitige Aufrufe teilen nichts (kein Stand, kein Name); Aufruf-Voraussetzung: unabhängig, kein Fall nötig außer „kein gemeinsamer Zustand“.
11. **Ort der Namen.** Wo steht der Name der Ziele für die Ausnahmeliste der schnellen Kette und die Zuordnung von `SPEC-060`: im Fragment; die Deckungs-Gegenprobe der Kette zählt sie mit (Folge von `slice-harness-gegenproben-aus-gates`).

**Fragen an den Nutzer** (jede mit Empfehlung des Planners; eine Empfehlung ist keine Entscheidung):

1. **Netzfehler und Rate-Limit: überspringen oder rot?** Empfehlung: überspringen, nennen und zählen (fail-open nach d-check), 404 und „anders“ rot. Gegenposition: Alles rot, was nicht antwortet; dann färbt jede Störung von GitHub den Lauf.
2. **Meldung: nur Rot oder zusätzlich ein Issue?** Empfehlung: nur Rot.
3. **`harness/conventions.md` nachziehen.** Dort steht, ein Sensor für den Audit komme nicht mit, weil er Netz über `bash`, `git` und `docker` hinaus braucht („Dauerhaft, solange diese drei die einzigen Host-Abhängigkeiten sind“). Der Sensor läuft in einem Image, braucht am Host nur `docker` und gehört nicht zu `make gates`; die Aussage stimmt dann nicht mehr. Empfehlung: den Satz im Abschnitt *Der mitgelieferte Baum altert still* und die Zeile der Tabelle zu `modul-02-harness-bootstrap.md` §Freshness-Audit ändern auf „Träger kommt mit, außerhalb der Gates“, mit dem Grund. Das Dokument gehört dem Repo; die Entscheidung ist deine.
4. **Der Vergleich des Assets** (Tag nachträglich neu gebaut, Exit 4 bei d-check). Empfehlung: nein, `make baseline-verify` hält die Integrität lokal; anderer Vorgang.
5. **Uhrzeit des Nachtlaufs.** Empfehlung: 01:00 UTC wie d-check, `gegenproben.yml` drei bis vier Stunden später.

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
Slice-Typ) und die **vier Pflichtkriterien**.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist die Default-Sub-Area `REPO` (`harness/conventions.md` §Modus-Deklaration pro Sub-Area), die einzige deklarierte; sie trägt eine einheitliche Greenfield-Deklaration, daher keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (`AGENTS.md` §3.10) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (`AGENTS.md` §3.12) — Skript und Datei sind ein neuer Vertrag; `BEO-REPO/postgres-versionsliste-altert` und `BEO-REPO/verkoerperte-regel-nach-baseline-sprung-nicht-nachgezogen` — beide sind der Anlass dieses Slice (Pins und Baseline altern still); `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` — der Sensor ergänzt die Aussage in `harness/conventions.md`, er ersetzt sie nicht still. Kein Eintrag erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Der Slice trägt keine Sendung eines Gebers (`SPEC-063` und `SPEC-064` sind gelesen). Die Zusage an den Nehmer `slice-erster-release-freigabe` (Beleg des Audits vor dem Tag) und an `slice-harness-ci-freshness-images` (Erweiterung des Skripts) tragen diese Nehmer in ihrem §1 mit der Kennung dieses Slice. Liefer-Punkte 3 (Skript und Ziele | Workflow | Gegenprobe samt Index). Schichten 2: Workflows und Harness (`tools/harness`, `harness/mk`, `harness/vertraege`, `harness/README.md`, `harness/conventions.md`); Planung, Vertrag und ADR sind keine Schicht; Wartungs-Doku und Handbuch bleiben unberührt. Schätzung 555 bis 800 (§4), unter der Marke.
