# Slice slice-harness-d-check-v0-85: d-check v0.85.0, Freshness-Audit, Nachzählen im Review

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle. Die Closure-Bedingung ist die DoD dieses Slice. Eingesammelt wird er
von der nächsten Welle-Closure. Vorgezogen nach Entscheidung des Nutzers vom 2026-10-09
(d-check `v0.85.0` behebt zwei HIGH-CVEs im Image des Doku-Gates): nach
`slice-v1-abschluss-einspielen` und vor `slice-v1-abschluss-einspielen-laufsteuerung`
(WIP-Limit 1); Reihenfolge in §4 *Start*.

**Bezug:** [MR-000](../../../../harness/conventions.md#mr-000--baseline-aussage) (Baseline-Aussage: adoptierter Stand und Asset-Quelle — der Freshness-Audit prüft den Stand gegen die Release-Liste, ohne ihn zu heben). Keine Anforderung des Lastenhefts im Scope: Der Slice ändert die Prüfumgebung und einen Skill, nicht das Produkt.

**Berührte Spec-Stellen:** — (die Spezifikation nennt d-check nicht; die Werkzeugverträge im Abschnitt für Harness-Werkzeuge bleiben unverändert)

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

**Ziel:** Das Doku-Gate läuft auf d-check `v0.85.0` (Tag und Digest gepinnt, ohne die zwei HIGH-CVEs des Images `v0.82.0`), die Workflow-Commands sagen über das Doku-Gate nur zu, was die Gegenprobe zeigt, §7 hält das Ergebnis des Freshness-Audits der vendored Baseline gegen die Release-Liste des Kurs-Repos fest, und der Reviewer-Skill prüft das Nachzählen beim Eintragen einer Sendung nach `AGENTS.md` §3.13.

**Herkunft der drei Teile** (Entscheidungen des Nutzers vom 2026-10-09):

- **Pin:** `ghcr.io/pt9912/d-check:v0.85.0` mit Digest
  `sha256:c07f1fe6053b1f790c4a1e01a76bcf4d3a8ff85e6eb609fe1aaaaf6ab6f09abe`; ein Probelauf
  des Nutzers mit diesem Image fand im Arbeitsbaum 424 Dateien, 0 Befunde.
- **Freshness-Audit:** Das Kurs-Repo führt `v6.17.0`, das Repo ist auf `v6.16.0` gepinnt
  (`harness/conventions.md` §Baseline). Geprüft und berichtet wird hier, gehoben nicht.
- **Reviewer-Skill:** Die MEDIUM-Klasse *Adresse nimmt nicht an* in
  `.harness/skills/reviewer.md` nennt bisher nur die Annahme durch den Nehmer (§3.13, seit
  slice-lint-bestand-kern-driven), nicht das Nachzählen von Liefer-Punkten und Schichten
  des Nehmers mit der Sendung (§3.13 *Nachzählen beim Eintragen*, seit
  slice-v1-abschluss-einspielen; `BEO-REPO/slice-waechst-durch-uebernahmen`, verkörpert).
- **Berichtigung von DoD-Punkt 1** (Entscheidung des Nutzers vom 2026-10-09 nach den Belegen
  in §7): Die Teil-Zusage „bares `LH-`-Token ist rot“ war vor dem Code nicht erfüllbar, weil
  die erste Abgrenzung unten `.d-check.yml` unverändert lässt und das `ids`-Muster für
  Lastenheft-Kennungen dort nur als Kommentar steht. DoD-Punkt 1 nennt jetzt das aktive,
  gegenprobierte Verhalten; die Teil-Zusage geht an `slice-harness-lh-links-pflicht`. Der
  Block *Strenges Doc-Gate* in drei Workflow-Commands sagt dieselbe Link-Pflicht für `LH-`,
  `ADR-` und `MR-` zu (Nebenbefund des Implementers) und wird hier enger gefasst, nicht erst
  im Nehmer: Dieser startet nach welle-v1-abschluss, und bis dahin läse jeder Lauf von
  Planner, Implementer und Welle-Closure eine Prüfung, die es nicht gibt (`AGENTS.md`
  §3.11); die Korrektur ist Text in der Schicht, die dieser Slice ohnehin berührt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Neue Opt-in-Module von d-check aktivieren (die Zeile `modules:` in `.d-check.yml`) — ein
  anderer Vorgang: Ihre Aktivierung ist eine eigene Entscheidung (`harness/conventions.md`,
  Regelblock `modul-15-observability.md` §Doku-Konsistenz-Drift), und der Nutzer hat sie
  für diesen Slice ausgeschlossen. Die Zeile bleibt `[links, anchors, ids, matrix, spans]`.
- Ein bares `LH-`-Token rot machen (das `ids`-Muster für Lastenheft-Kennungen in
  `.d-check.yml` aktivieren, den Bestand verlinken, `MR-` gleich behandeln, den Block
  *Strenges Doc-Gate* auf das neue Verhalten erweitern) — ein Folge-Slice übernimmt es:
  `slice-harness-lh-links-pflicht` (`open/`, nimmt die Sendung in §1 und DoD mit dieser
  Kennung an). Es verschärft das Gate um gemessen 357 Befunde im lebenden Bestand und ist
  eine eigene Entscheidung; `.d-check.yml` bleibt hier unverändert.
- Das Upgrade der Baseline auf `v6.17.0` (vendored Baum, `SHA256SUMS`, Abgleich der
  verkörperten Regeln) — ein anderer Vorgang, eigener Slice nach Entscheidung des Nutzers;
  dieser Slice liefert nur den Befund des Audits in §7, den der Planner danach zuschneidet.
- `d-check.mk` neu erzeugen (`d-check --print-mk`) — Bestand bleibt bewusst stehen: Die
  Datei trägt Adaptionen ggü. der Ausgabe von d-check (Kopfkommentar: `docs-check` als
  Befund-Gate, Muster von `doc-help`, Digest-Pin), und `harness/mk/doc-gate.mk` bindet die
  Ziele `doc-immutable` und `doc-commits` mit Rezept vor. Umgepinnt werden nur
  `DCHECK_IMAGE` und `DCHECK_DIGEST`.
- Weitere Änderungen an `.harness/skills/reviewer.md` oder an `AGENTS.md` §3.13 — die Regel
  selbst ist verkörpert; dieser Slice trägt sie nur in die Klassifikation des Reviews.
- Produkt-Code, Tests unter `test/` und die Spezifikation — Schicht-Abgrenzung: Der Slice
  ändert einen Pin, einen Skill, einen Block in drei Workflow-Commands und diesen Plan.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **d-check `v0.85.0`:** `DCHECK_IMAGE` in `d-check.mk` nennt
      `ghcr.io/pt9912/d-check:v0.85.0`, `DCHECK_DIGEST` den Digest
      `sha256:c07f1fe6053b1f790c4a1e01a76bcf4d3a8ff85e6eb609fe1aaaaf6ab6f09abe`; der Digest
      gehört zum Tag (Beleg in §7: Abruf des Tags, Digest der Ausgabe); `make docs-check`
      meldet 0 Befunde. Gegenprobe in einer Kopie des Arbeitsbaums mit dem neuen Image: ein
      toter Anker, eine Kennung `ADR-` ohne Link, ein totes Linkziel und ein Link aus einem
      Spec-Stratum auf eine ADR (Matrix) sind je rot (Beleg in §7: Mutation · Befundzeile).
      Der Block *Strenges Doc-Gate* in `.claude/commands/implement-slice.md`,
      `.claude/commands/plan-welle.md` und `.claude/commands/close-welle.md` sagt über das
      Doku-Gate nur zu, was `.d-check.yml` aktiv prüft und eine Zeile der Gegenprobe in §7
      rot zeigt (heute: `ADR-` ohne Link überall, auch in Inline-Code; nicht: `LH-` ohne
      Link, `MR-` außerhalb der Spec-Dateien, Pfade in Inline-Code — das Modul `codepaths`
      ist nicht aktiv); was er Agenten weiter abverlangt (Kennungen als Links schreiben),
      steht als Regel, nicht als Befund des Gates (`AGENTS.md` §3.11).
- [ ] **Freshness-Audit:** §7 nennt die Release-Liste des Kurs-Repos (neuester Tag gegen den
      gepinnten `v6.16.0`) und für das Delta bis zum neuesten Tag je Eintrag des
      Adaptions-Blocks (`MR-000`, aktive und aufgelöste) den Ausgang nach
      Baseline-Regelwerk `modul-02-harness-bootstrap.md` §Freshness-Audit der vendored
      Baseline, dazu die Regelblöcke der Tabelle in `harness/conventions.md`, die das Delta
      umbenennt, hinzufügt oder wegnimmt; ohne Änderung am vendored Baum.
- [ ] **Reviewer-Skill:** `.harness/skills/reviewer.md` §Klassifikation (MEDIUM) führt das
      Nachzählen: Trägt der Diff eine Sendung in einen Nehmer ein, ohne dass dessen §8 die
      Zählung von Liefer-Punkten und Schichten mit der Sendung nennt, oder liegt der Nehmer
      damit über drei Punkten oder zwei Schichten, ist das mindestens MEDIUM gegen
      `AGENTS.md` §3.13 (seit slice-v1-abschluss-einspielen).
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
| `d-check.mk` | update | `DCHECK_IMAGE` auf `v0.85.0`, `DCHECK_DIGEST` auf den Digest des Tags; sonst unverändert (§1) |
| `.harness/skills/reviewer.md` | update | MEDIUM *Adresse nimmt nicht an* um das Nachzählen nach `AGENTS.md` §3.13 ergänzt |
| `.claude/commands/implement-slice.md`, `.claude/commands/plan-welle.md`, `.claude/commands/close-welle.md` | update | Block *Strenges Doc-Gate* auf das gegenprobierte Verhalten gefasst (DoD-Punkt 1, §1 *Berichtigung*): rot sind `ADR-` ohne Link (auch in Inline-Code), totes Linkziel, toter Anker, in `implement-slice.md` dazu der Verweis aus `spec/` auf ADR, Slice und `MR-`; nicht geprüft `LH-` ohne Link, `MR-` außerhalb von `spec/`, Pfade in Inline-Code. Unter `.claude/agents/` und `.harness/skills/` steht kein gleichlautender Block (§7) |
| dieser Plan, §7 | update | Belege des Pins (Digest, Gegenprobe) und Befund des Freshness-Audits |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen` liegt in `done/` (WIP-Limit
1). Vorgezogen vor `slice-v1-abschluss-einspielen-laufsteuerung`, Schritt 9 der Reihenfolge
in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md); danach geht die Welle dort weiter.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): `v0.85.0` meldet im Arbeitsbaum
  Befunde, die mehr als eine kleine Korrektur der Doku verlangen, oder `d-check.mk` muss
  neu erzeugt werden, weil sich Ziele oder Aufruf von d-check geändert haben (§1). Schnitt
  dann: der Pin als eigener Slice; Audit und Skill sind ohne ihn lieferbar.
- `in-progress` → `open` (blockiert — Carveout?): Der Tag `v0.85.0` ist nicht abrufbar, oder
  sein Digest ist nicht der genannte; dann zuerst die Klärung mit dem Nutzer, der Pin
  bleibt bis dahin auf `v0.82.0`.

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

**Randformen** (`AGENTS.md` §3.12): keine. Der Slice liefert keinen neuen Vertrag: Der Pin
ändert das Image eines bestehenden Gates, der Audit ist eine Handlung ohne Sensor
(`harness/conventions.md`, *Der mitgelieferte Baum altert still*), und die Klasse im
Reviewer-Skill trägt eine verkörperte Regel in die Klassifikation.

**Risiken:**

- `v0.85.0` wertet die aktiven Module anders aus als `v0.82.0`, und der Probelauf mit 0
  Befunden verdeckt eine Prüfung, die still schwächer geworden ist
  (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, 2×); die Gegenprobe in DoD-Punkt
  1 deckt vier Klassen, nicht alle (`spans` nicht) — **Ausgang:** offen bis Closure. Die
  Gegenprobe zeigte, dass schon der Plan das Gate weiter las, als es wirkt (bares
  `LH-`-Token, §7); ob das ein Beleg zu diesem Eintrag ist, entscheidet die Closure, die
  Teil-Zusage selbst ging an `slice-harness-lh-links-pflicht` (§1).
- Das Delta bis `v6.17.0` berührt eine verkörperte Regel dieses Repos, ohne dass der Audit es
  sieht, weil er nur Adaptions-Block und Regelblock-Tabelle durchgeht
  (`BEO-REPO/verkoerperte-regel-nach-baseline-sprung-nicht-nachgezogen`, 1×) — **Ausgang:**
  offen bis Closure; der Befund geht an den Upgrade-Slice, den der Planner nach diesem Slice
  anlegt.

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

**Pin (DoD-Punkt 1).** Am 2026-10-09:

- *Digest gehört zum Tag:* `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.85.0`
  meldet einen OCI-Index mit `Digest: sha256:c07f1fe6053b1f790c4a1e01a76bcf4d3a8ff85e6eb609fe1aaaaf6ab6f09abe`
  (Manifeste `linux/amd64` `sha256:eb73e50a…`, `linux/arm64` `sha256:69493255…`);
  `docker pull` desselben Tags meldet denselben Digest, `docker image inspect --format
  '{{json .RepoDigests}}'` ebenso. Gleich dem Digest der Release-Notiz.
- *Ziele und Aufruf unverändert:* `--print-mk` von `v0.82.0` (alter Digest) und `v0.85.0`
  unterscheiden sich nur in der Zeile `DCHECK_IMAGE`; die Rückführung nach §4 (`d-check.mk`
  neu erzeugen) tritt nicht ein. Geändert sind in `d-check.mk` nur `DCHECK_IMAGE` und
  `DCHECK_DIGEST`; der Kopfkommentar („auf den erzeugenden Image-Digest gepinnt“) bleibt wahr.
- *Weitere Nennungen des alten Pins:* `grep` nach `v0.82.0` und dem alten Digest über
  `*.md`, `*.mk`, `*.yml`, `*.sh` außerhalb von `.harness/baseline/` findet ihn nur in
  `d-check.mk` (umgepinnt), im Drift-Log der Roadmap (Zeile vom 2026-10-06) und in
  `done/slice-harness-upgrade-v6-16.md` — beides Zeitdokumente, bleiben stehen.
  `harness/README.md` und `harness/conventions.md` nennen keine Version, nur „Tag und Digest
  stehen in `d-check.mk`“.
- *`make docs-check` am Arbeitsbaum:* `d-check: 435 Datei(en) geprüft, 0 Befund(e)`.

**Gegenprobe** in Kopien des Arbeitsbaums (je Mutation eine frische Kopie ohne `.git` im
Scratchpad, angehängt an `harness/README.md`, sofern nicht anders genannt), je Kopie gegen
`v0.85.0` (`c07f1fe6…`) und zum Vergleich gegen `v0.82.0` (`d28e9437…`); unveränderte Kopie
unter beiden 0 Befunde, Exit 0.

| Zusage | Mutation | Befundzeile (`v0.85.0`; `v0.82.0` identisch) |
|---|---|---|
| toter Anker ist rot (`anchors`) | Link auf `conventions.md#gibt-es-nicht` | `harness/README.md:141 conventions.md#gibt-es-nicht anchor-missing`, Exit 1 |
| Kennung `ADR-` ohne Link ist rot (`ids`) | `ADR-` mit vier Ziffern (0001) blank im Fließtext | `harness/README.md:141` · die Kennung · `id-unlinked`, Exit 1 |
| totes Linkziel ist rot (`links`, Kontrolle) | Link auf `gibt-es-nicht.md` | `harness/README.md:141 gibt-es-nicht.md target-missing`, Exit 1 |
| Spec-Stratum nennt keine ADR (`matrix`, Kontrolle) | Link aus `spec/architecture.md` auf eine ADR | `spec/architecture.md:645 … matrix-forbidden Referenz spec-straten → adr`, Exit 1 |
| bares `LH-`-Token ist rot | `LH-FA-01` blank in `harness/README.md`, in `spec/spezifikation.md` und in diesem Plan (drei Kopien) | **kein Befund**, 0 Befunde, Exit 0 — unter `v0.85.0` **und** `v0.82.0` |

Das Modul `spans` ist nicht mutiert. **Befund zur DoD:** Ein bares `LH-`-Token ist unter
der Konfiguration dieses Repos kein Befund, unabhängig von der Version: Das `ids`-Muster
für die Lastenheft-Kennungen steht in `.d-check.yml` nur als Kommentar, aktiv ist allein das
Muster für `ADR-`. Dasselbe hält schon `done/slice-lastenheft-pruefbarkeit.md` §1 fest
(„einen Link verlangt es nur für ADR-Kennungen, nicht für Lastenheft-Kennungen“). Keine
Regression von `v0.85.0`; die Teil-Zusage „bares `LH-`-Token rot“ in DoD-Punkt 1 ist mit
`.d-check.yml` unverändert (§1) nicht erfüllbar. Der Implementer hakt DoD-Punkt 1 deshalb
nicht ab und gibt den Punkt an den Planner (DoD berichtigen, oder ein eigener Slice
aktiviert das Muster als Entscheidung über das Gate). Ebenso zu weit griff der Block
*Strenges Doc-Gate* in drei Commands (`AGENTS.md` §3.11); der Planner hat DoD-Punkt 1
berichtigt (Commit `dfb7e88`), der Block ist enger gefasst (Gegenprobe unten).

**Gegenprobe zum Block *Strenges Doc-Gate*** (Berichtigung von DoD-Punkt 1), nur gegen
`v0.85.0` (`c07f1fe6…`), am Stand `dfb7e88` mit den geänderten Commands im Arbeitsbaum. Je
Mutation eine frische Kopie im Scratchpad: jeder Eintrag der obersten Ebene außer `.git` mit
`cp -r` (Symlinks unter `.claude/rules/` bleiben Symlinks, ohne `-p`), danach gelöscht.
Unveränderte Kopie: `436 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

| Zusage im Block | Mutation | Befundzeile |
|---|---|---|
| `ADR-` ohne Link ist rot, auch in Inline-Code | `ADR-` mit vier Ziffern (0001) in Backticks, angehängt an `harness/README.md` | `harness/README.md:141` · die Kennung · `id-unlinked`, 1 Befund, Exit 1 |
| `MR-` in `spec/` ist rot, blank | `MR-001` im Fließtext, angehängt an `spec/spezifikation.md` | `spec/spezifikation.md:2668 MR-001 matrix-forbidden … spec-straten → adaptionsblock`, 1 Befund, Exit 1 |
| `MR-` in `spec/` ist rot, in Inline-Code | `` `MR-001` `` in Backticks, angehängt an `spec/spezifikation.md` | dieselbe Zeile, 1 Befund, Exit 1 |
| Slice aus `spec/` ist rot | `slice-harness-d-check-v0-85` blank, angehängt an `spec/architecture.md` | `spec/architecture.md:645 slice- matrix-forbidden … spec-straten → slice`, 1 Befund, Exit 1 |
| nicht zugesagt: `MR-` außerhalb von `spec/` | `MR-001` im Fließtext, angehängt an `harness/README.md` | **kein Befund**, 0 Befunde, Exit 0 |
| nicht zugesagt: Pfade in Inline-Code | `` `tools/gibt-es-nicht.sh` ``, angehängt an `harness/README.md` | **kein Befund**, 0 Befunde, Exit 0 (`codepaths` nicht in `modules:`) |

Toter Anker, totes Linkziel und ADR aus `spec/` stehen in der Tabelle davor; `LH-` ohne Link
ebenda (kein Befund). Ein Verweis aus `spec/` auf eine Welle und auf eine superseded ADR ist
nicht mutiert — der Block sagt beides nicht mehr zu, ebenso nicht die Ausnahme von
`docs/reviews/**` (`.d-check.yml` nimmt den Pfad nur aus der Status-Prüfung der Matrix aus).

*Fundstellen:* `grep -rn -i "Strenges Doc-Gate\|klickbare\|Anker-Link\|codepaths"` über
`.claude/`, `.harness/skills/` und `AGENTS.md` findet den Block nur in den drei Commands; dazu
in `.claude/commands/plan-welle.md` Schritt 8 „Kennungen als Anker-Links“ — eine Regel an den
Planner, keine Zusage über das Gate, bleibt. `.claude/agents/*.md` und die Skills führen
keine gleichlautende Stelle.

**Freshness-Audit (DoD-Punkt 2).** `gh release list -R pt9912/ai-harness-course` am
2026-10-09: neuester Tag `v6.17.0` (Latest, 2026-10-07), davor `v6.16.0` (2026-10-06, der
Pin). Delta: genau ein Release. Das Asset `lab-regelwerk.zip` von `v6.17.0` nur in den
Scratchpad geladen (`sha256 afe50df8…`, gleich dem `SHA256SUMS` des Releases) und mit
`diff -r` gegen `.harness/baseline/v6.16.0/` gehalten; am vendored Baum nichts geändert.

- *Dateien:* dieselben, keine neu, umbenannt oder entfallen; keine Überschrift in einer
  `.md` geändert.
- *Inhalt über Quell-Zeile und Tag hinaus:* `grundlagen-harness-dateien.md` und
  `modul-13-quality-gates.md` — die Disjunktheit der Teile eines Gate-Index (kein Target in
  zwei Teilen) prüft d-check nur mit eigenem Schalter (ab `v0.83.0`), sonst steht die Regel
  nur im Briefing; `templates/.d-check.yml` führt dazu auskommentiert
  `authority-disjoint: true`; `regelwerk/README.md` Stand-Zeile „Kurs-Welle 160 ·
  2026-10-07“; `templates/AGENTS.template.md` und `templates/harness/conventions.template.md`
  nur die Asset-URL mit dem Tag.
- *Regelblock-Tabelle in `harness/conventions.md`:* kein Regelblock umbenannt, hinzugefügt
  oder weggenommen. Inhaltlich berührt sind die Zeilen `grundlagen-harness-dateien.md`
  (Teil des Gate-Index eines Werkzeugs, *kommt nicht mit*) und `modul-13-quality-gates.md`
  (*Träger kommt mit*); beide Werte bleiben, denn dieses Repo führt keinen Gate-Index in
  Teilen über das Modul `targets` (`.d-check.yml` nennt weder `targets` noch `authority`).
- *Adaptions-Block, je Eintrag der Ausgang:* `MR-000` **bleibt gültig** — das Delta ändert
  weder Verzeichniskonvention, Lifecycle, Carveout-Disziplin noch ID-Schema. `MR-001`
  (aufgelöst) **bleibt gültig** als aufgelöster Eintrag — das Delta berührt
  `grundlagen-referenz-richtung.md` nur in der Quell-Zeile; kein Nachfolge-Eintrag. Aktive
  Adaptionen: keine.
- *Nicht gelaufen:* die Stichprobe gegen den Bestand (Modul 2, siebte Eigenschaft) — die DoD
  verlangt sie nicht; sie gehört zum Upgrade-Slice, den der Planner anlegt.

**Reviewer-Skill (DoD-Punkt 3).** `.harness/skills/reviewer.md` §Klassifikation, MEDIUM
*Nehmer nicht nachgezählt*. Ein Wächter mit Mutation ist das nicht: Der Skill ist
Urteilsgrundlage des Review, kein Sensor prüft ihn; belegbar ist nur der Text
(`grep -n "Nehmer nicht nachgezählt" .harness/skills/reviewer.md`).

**Läufe.** `make gates` am Stand `79d7bb8` (Pin und Skill committet): Exit 0; darin
`baseline-verify: v6.16.0 OK — 54 Dateien`, `d-check: 435 Datei(en) geprüft, 0 Befund(e)`
mit Digest `c07f1fe6…`, Integrationstests, alle Gegenproben grün. Der letzte Lauf nach
dem Commit dieser Belege steht im Bericht.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area, `*` (Kürzel
`REPO`, Greenfield, `harness/conventions.md`); dieser Slice berührt nur sie.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-REPO/` am Stand `5b7bd50` gesichtet (Zähler = Dateien
unter `evidence/`). Treffer:

- `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — ein neues Image des
  Doku-Gates kann eine Regel anders auswerten; daher die Gegenprobe in DoD-Punkt 1 und das
  erste Risiko in §6. Ein Beleg entsteht erst, wenn sie eine Abweichung zeigt; mit ihm
  stünde der Eintrag bei 3×.
- `BEO-REPO/gate-regel-ersetzt-statt-ergaenzt` (1×, offen) — der Slice ersetzt keine Regel,
  nur den Pin; kein Beleg erwartet.
- `BEO-REPO/verkoerperte-regel-nach-baseline-sprung-nicht-nachgezogen` (1×, offen) — betrifft
  den Upgrade-Slice, nicht diesen; hier nur als Risiko des Audits (§6).
- `BEO-REPO/slice-waechst-durch-uebernahmen` (3×, verkörpert) und
  `BEO-REPO/folge-slice-adresse-nimmt-nicht-an` (5×, verkörpert) — Gegenstand von DoD-Punkt
  3. Nachgezählt beim Anlegen: drei Liefer-Punkte, zwei Schichten (Gate-Konfiguration
  `d-check.mk`, Skill `.harness/skills/reviewer.md`); der Plan selbst zählt nicht.
  Nachgezählt bei der Berichtigung vom 2026-10-09: weiter drei Liefer-Punkte (die Fassung
  der Commands gehört zu DoD-Punkt 1, sie sagt das gegenprobierte Verhalten desselben Gates
  aus) und zwei Schichten (Gate-Konfiguration; Agenten-Anweisungen: Skill und Commands).
  Der Nehmer der Teil-Zusage, `slice-harness-lh-links-pflicht`, ist dort nachgezählt (§8).

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
