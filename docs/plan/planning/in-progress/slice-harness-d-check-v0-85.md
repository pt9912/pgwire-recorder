# Slice slice-harness-d-check-v0-85: d-check v0.86.1, Freshness-Audit, Nachzählen im Review

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

**Ziel:** Das Doku-Gate läuft auf d-check `v0.86.1` (Tag und Digest gepinnt, ohne die zwei HIGH-CVEs des Images `v0.82.0`), die Workflow-Commands sagen über das Verhalten des Doku-Gates nichts zu, sondern nennen die Link-Regel für Kennungen und verweisen auf `.d-check.yml`, §7 hält das Ergebnis des Freshness-Audits der vendored Baseline gegen die Release-Liste des Kurs-Repos fest, und der Reviewer-Skill prüft das Nachzählen beim Eintragen einer Sendung nach `AGENTS.md` §3.13.

**Herkunft der drei Teile und der Berichtigung** (Entscheidungen des Nutzers vom 2026-10-09):

- **Pin:** `ghcr.io/pt9912/d-check:v0.85.0` mit Digest
  `sha256:c07f1fe6053b1f790c4a1e01a76bcf4d3a8ff85e6eb609fe1aaaaf6ab6f09abe`; ein Probelauf
  des Nutzers mit diesem Image fand im Arbeitsbaum 424 Dateien, 0 Befunde. Nach
  Entscheidung des Nutzers vom 2026-10-10 ist der Pin auf `ghcr.io/pt9912/d-check:v0.86.0`
  mit Digest `sha256:d90200e94db311a70f9b15045e1290feecde03e9185ad78fb1ebea265e6ee44b`
  gehoben; die Kennung dieses Slice bleibt. Nach Information des Nutzers vom 2026-10-10 ist
  der Pin auf `ghcr.io/pt9912/d-check:v0.86.1` mit Digest
  `sha256:3e0b9779a71e2fa942960e8b961428fba036455c535513d34a75bfb799ffce0e` gehoben, ein
  Sicherheitsrelease (`golang.org/x/net` `v0.60.0`).
- **Freshness-Audit:** Das Kurs-Repo führt `v6.17.0`, das Repo ist auf `v6.16.0` gepinnt
  (`harness/conventions.md` §Baseline). Geprüft und berichtet wird hier, gehoben nicht.
- **Reviewer-Skill:** Die MEDIUM-Klasse *Adresse nimmt nicht an* in
  `.harness/skills/reviewer.md` nennt bisher nur die Annahme durch den Nehmer (§3.13, seit
  slice-lint-bestand-kern-driven), nicht das Nachzählen von Liefer-Punkten und Schichten
  des Nehmers mit der Sendung (§3.13 *Nachzählen beim Eintragen*, seit
  slice-v1-abschluss-einspielen; `BEO-REPO/slice-waechst-durch-uebernahmen`, verkörpert).
- **Berichtigung von DoD-Punkt 1** (Entscheidung des Nutzers vom 2026-10-09 nach den Belegen
  in §7): Die Teil-Zusage „bares `LH-`-Token ist rot“ war vor dem Code nicht erfüllbar, weil
  der Plan `.d-check.yml` nicht vorsah — die Schicht-Abgrenzung unten nannte beim Anlegen
  einen Pin, einen Skill und diesen Plan, §3 nannte die Datei nicht — und das `ids`-Muster
  für Lastenheft-Kennungen dort nur als Kommentar steht. DoD-Punkt 1 nennt jetzt das aktive,
  gegenprobierte Verhalten; die Teil-Zusage geht an `slice-harness-lh-links-pflicht`. Der
  Block *Strenges Doc-Gate* in drei Workflow-Commands sagt dieselbe Link-Pflicht für `LH-`,
  `ADR-` und `MR-` zu (Nebenbefund des Implementers) und wird hier enger gefasst, nicht erst
  im Nehmer: Dieser startet nach welle-v1-abschluss, und bis dahin läse jeder Lauf von
  Planner, Implementer und Welle-Closure eine Prüfung, die es nicht gibt (`AGENTS.md`
  §3.11); die Korrektur ist Text in der Schicht, die dieser Slice ohnehin berührt. Nach Review
  F-555, F-557 und F-561 nimmt der Block umzäunte Code-Blöcke ausdrücklich aus, nennt für
  eine ADR-Kennung ohne Link in `spec/` den gemeldeten Code `id-unlinked` und sagt in
  `implement-slice.md` den Link auf eine superseded ADR samt Ausweg zu; je Zusage eine Zeile
  der Gegenprobe in §7. Nach Verifikation V-134 und V-138 (Entscheidung des Koordinators vom
  2026-10-09): Das Aufzählen der Ausnahmen konvergiert nicht — jede Runde fand eine weitere
  Form. Der Block sagt deshalb eine kleine, geschlossene Menge zu und zählt keine Ausnahme mehr
  auf: rot sind die Kernfälle, die die Gegenprobe zeigt, eingegrenzt nach Ort und Form; alles
  andere sagt er nicht zu und verweist auf `.d-check.yml`. Die gemessenen Ausnahmen bleiben in
  §7 als Messbefund, nicht als Zusage; ebenso der Link auf eine superseded ADR (F-561), den der
  Block nicht mehr zusagt. Nach Verifikation V-141 bis V-144 (Entscheidung des Nutzers vom
  2026-10-10): Auch die geschlossene Menge hielt nicht — Formen über einen Zeilenumbruch und
  eingerückte Überschriften blieben grün. Der Block sagt deshalb kein rot/grün-Verhalten des
  Gates mehr zu. Er nennt die Regel, Kennungen (`LH-`, `ADR-`, `MR-`) als Anker-Links zu
  schreiben, dass `LH-` ohne Link noch nicht erzwungen ist (`slice-harness-lh-links-pflicht`),
  und dass `.d-check.yml` festlegt, was das Gate prüft, und `make docs-check` maßgeblich ist.
  Alle Gegenproben zum Block bleiben in §7 als Messbefund (Stand `v0.85.0`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Neue Opt-in-Module von d-check aktivieren (die Zeile `modules:` in `.d-check.yml`) — ein
  anderer Vorgang: Ihre Aktivierung ist eine eigene Entscheidung (`harness/conventions.md`,
  Regelblock `modul-15-observability.md` §Doku-Konsistenz-Drift), und der Nutzer hat sie
  für diesen Slice ausgeschlossen. Die Zeile bleibt `[links, anchors, ids, matrix, spans]`.
- Ein bares `LH-`-Token rot machen (das `ids`-Muster für Lastenheft-Kennungen in
  `.d-check.yml` aktivieren, den Bestand verlinken, `MR-` gleich behandeln, den Block
  *Strenges Doc-Gate* auf das neue Verhalten erweitern) — ein Folge-Slice übernimmt es:
  `slice-harness-lh-links-pflicht` (`open/`, nimmt die Sendung in §1 und DoD mit dieser
  Kennung an; den Bestand in ADRs und in Nutzer- und Wartungs-Doku hat er nach Review F-556
  an `slice-harness-lh-links-bestand` abgeschnitten, der vor ihm startet). Es verschärft das Gate um gemessen 357 Befunde im lebenden Bestand und ist
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
  ändert einen Pin, einen Skill, einen Block in drei Workflow-Commands und diesen Plan; als
  Planung dazu die Nehmer-Pläne `open/slice-harness-lh-links-pflicht.md` und
  `open/slice-harness-lh-links-bestand.md` und zwei Zeilen im Drift-Log der Roadmap
  (2026-10-09, Anlage und Schnitt von `slice-harness-lh-links-pflicht`).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **d-check `v0.86.1`:** `DCHECK_IMAGE` in `d-check.mk` nennt
      `ghcr.io/pt9912/d-check:v0.86.1`, `DCHECK_DIGEST` den Digest
      `sha256:3e0b9779a71e2fa942960e8b961428fba036455c535513d34a75bfb799ffce0e`; der Digest
      gehört zum Tag (Beleg in §7: Abruf des Tags, Digest der Ausgabe); `make docs-check`
      meldet 0 Befunde. Gegenprobe in einer Kopie des Arbeitsbaums mit dem neuen Image: ein
      toter Anker, eine Kennung `ADR-` ohne Link, ein totes Linkziel und ein Link aus einem
      Spec-Stratum auf eine ADR (Matrix) sind je rot (Beleg in §7: Mutation · Befundzeile).
      Der Block *Strenges Doc-Gate* in `.claude/commands/implement-slice.md`,
      `.claude/commands/plan-welle.md` und `.claude/commands/close-welle.md` sagt kein
      rot/grün-Verhalten des Gates zu (Entscheidung des Nutzers vom 2026-10-10 nach
      Verifikation V-141 bis V-144); er nennt nur die Regel, Kennungen (`LH-`, `ADR-`, `MR-`)
      als Anker-Links zu schreiben, dass das Gate `LH-` ohne Link noch nicht erzwingt
      (`slice-harness-lh-links-pflicht`), und dass `.d-check.yml` festlegt, was das Gate prüft,
      und `make docs-check` maßgeblich ist.
- [ ] **Freshness-Audit:** §7 nennt die Release-Liste des Kurs-Repos (neuester Tag gegen den
      gepinnten `v6.16.0`) und für das Delta bis zum neuesten Tag je Eintrag des
      Adaptions-Blocks (`MR-000`, aktive und aufgelöste) den Ausgang nach
      Baseline-Regelwerk `modul-02-harness-bootstrap.md` §Freshness-Audit der vendored
      Baseline, dazu die Regelblöcke der Tabelle in `harness/conventions.md`, die das Delta
      umbenennt, hinzufügt oder wegnimmt; ohne Änderung am vendored Baum.
- [ ] **Reviewer-Skill:** `.harness/skills/reviewer.md` §Klassifikation (MEDIUM) führt das
      Nachzählen: Trägt der Diff eine Sendung in einen Nehmer ein, ohne dass dessen §8 im
      selben Commit die Zählung von Liefer-Punkten und Schichten mit der Sendung nennt, oder liegt der Nehmer
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
| `d-check.mk` | update | `DCHECK_IMAGE` auf `v0.86.1` (zuerst `v0.85.0`, dann `v0.86.0`, beide Hebungen am 2026-10-10), `DCHECK_DIGEST` auf den Digest des Tags; sonst unverändert (§1) |
| `.harness/skills/reviewer.md` | update | neue MEDIUM-Klasse *Nehmer nicht nachgezählt* nach `AGENTS.md` §3.13 *Nachzählen beim Eintragen*, mit der Bedingung „im selben Commit“ (Review F-560); die Klasse *Adresse nimmt nicht an* bleibt unverändert |
| `.claude/commands/implement-slice.md`, `.claude/commands/plan-welle.md`, `.claude/commands/close-welle.md` | update | Block *Strenges Doc-Gate* ohne Zusage über das Verhalten des Gates (DoD-Punkt 1, §1 *Berichtigung*, Verifikation V-141 bis V-144): die Regel, Kennungen als Anker-Links zu schreiben, `LH-` noch nicht erzwungen, `.d-check.yml` legt fest, `make docs-check` ist maßgeblich. Unter `.claude/agents/` und `.harness/skills/` steht kein gleichlautender Block (§7) |
| dieser Plan, §7 | update | Belege des Pins (Digest, Gegenprobe) und Befund des Freshness-Audits |
| `docs/plan/planning/open/slice-harness-lh-links-pflicht.md`, `docs/plan/planning/open/slice-harness-lh-links-bestand.md` | neu (Planung) | Nehmer der Teil-Zusage „bares `LH-`-Token ist rot“ und dessen abgeschnittener Bestand (§1, §3.13 mit Nachzählen in §8 dort); in `slice-harness-lh-links-pflicht` §6 dazu die Ausnahmen des `ids`-Moduls aus der Gegenprobe (Verifikation V-137) und die von d-check ausgelassenen Verzeichnisse (Verifikation V-139) |
| `docs/plan/planning/in-progress/roadmap.md` | update (Planung) | zwei Zeilen im Drift-Log vom 2026-10-09: Anlage und Schnitt von `slice-harness-lh-links-pflicht` |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-einspielen` liegt in `done/` (WIP-Limit
1). Vorgezogen vor `slice-v1-abschluss-einspielen-laufsteuerung`, Schritt 9 der Reihenfolge
in §5 von [welle-v1-abschluss](../welle-v1-abschluss.md); danach geht die Welle dort weiter.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der gepinnte Tag (`v0.86.1`) meldet im Arbeitsbaum
  Befunde, die mehr als eine kleine Korrektur der Doku verlangen, oder `d-check.mk` muss
  neu erzeugt werden, weil sich Ziele oder Aufruf von d-check geändert haben (§1). Schnitt
  dann: der Pin als eigener Slice; Audit und Skill sind ohne ihn lieferbar.
- `in-progress` → `open` (blockiert — Carveout?): Der gepinnte Tag (`v0.86.1`) ist nicht abrufbar, oder
  sein Digest ist nicht der genannte; dann zuerst die Klärung mit dem Nutzer, der Pin
  bleibt bis dahin auf dem vorigen Tag.

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

- `v0.86.1` (davor `v0.86.0` und `v0.85.0`) wertet die aktiven Module anders aus als `v0.82.0`, und der Probelauf mit 0
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

**Pin auf `v0.86.1` (DoD-Punkt 1, Information des Nutzers vom 2026-10-10: Sicherheitsrelease,
`golang.org/x/net` `v0.60.0`, laut Changelog keine Verhaltensänderung).** Am 2026-10-10:

- *Digest gehört zum Tag:* `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.86.1`
  meldet einen OCI-Index mit `Digest: sha256:3e0b9779a71e2fa942960e8b961428fba036455c535513d34a75bfb799ffce0e`
  (Manifeste `linux/amd64` `sha256:1d3b0a2b…`, `linux/arm64` `sha256:01e31f8f…`);
  `docker pull` desselben Tags, danach `docker image inspect --format '{{json .RepoDigests}}'`,
  meldet den Index-Digest (daneben den des Manifests `linux/amd64`). Gleich dem Digest der
  Release-Notiz.
- *Ziele und Aufruf unverändert:* `--print-mk` von `v0.86.0` (`d90200e9…`) und `v0.86.1`
  (`3e0b9779…`), je 76 Zeilen, unterscheiden sich nur in der Zeile `DCHECK_IMAGE` (`diff`,
  eine Zeile); die Rückführung nach §4 tritt nicht ein. Geändert sind in `d-check.mk` nur
  `DCHECK_IMAGE` und `DCHECK_DIGEST`; `.d-check.yml` ist unverändert.
- *Weitere Nennungen des alten Pins:* `grep` nach `v0.86.0` und `d90200e9` außerhalb von
  `.git`, `.harness/` und `docs/reviews/` findet ihn nur in diesem Plan (Herkunft und
  Belege mit Stand), nicht mehr in `d-check.mk`.
- *`make docs-check` am Arbeitsbaum:* `d-check: 439 Datei(en) geprüft, 0 Befund(e)` mit
  Digest `3e0b9779…`.

Gegenprobe gegen `v0.86.1` (`3e0b9779…`) in Kopien des Arbeitsbaums: je Mutation eine frische
Kopie im Scratchpad, jeder Eintrag der obersten Ebene außer `.git` mit `cp -r` ohne `-p`,
Mutation per `printf` angehängt, Lauf `docker run --rm --network none` mit dem Digest, danach
nur diese Kopie gelöscht. Unveränderte Kopie: `439 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

| Zusage | Mutation | Befundzeile (`v0.86.1`) |
|---|---|---|
| toter Anker ist rot (`anchors`) | Link auf `conventions.md#gibt-es-nicht`, an `harness/README.md` | `harness/README.md:141 conventions.md#gibt-es-nicht anchor-missing`, 1 Befund, Exit 1 |
| Kennung `ADR-` ohne Link ist rot (`ids`) | `ADR-` mit vier Ziffern (0001) blank im Fließtext, an `harness/README.md` | `harness/README.md:141` · die Kennung · `id-unlinked`, 1 Befund, Exit 1 |
| totes Linkziel ist rot (`links`) | Link auf `gibt-es-nicht.md`, an `harness/README.md` | `harness/README.md:141 gibt-es-nicht.md target-missing`, 1 Befund, Exit 1 |
| Spec-Stratum nennt keine ADR (`matrix`) | Link aus `spec/architecture.md` auf ADR 0001 | `spec/architecture.md:645 … matrix-forbidden Referenz spec-straten → adr`, 1 Befund, Exit 1 |
| Kontrolle: bares `LH-`-Token bleibt grün (nicht erzwungen, `slice-harness-lh-links-pflicht`) | `LH-FA-01` blank, angehängt an `harness/README.md`, an `spec/spezifikation.md` und an diesen Plan (drei Kopien) | kein Befund, je `439 Datei(en) geprüft, 0 Befund(e)`, Exit 0 — wie unter `v0.85.0` und `v0.82.0` |

**Pin auf `v0.86.0` (Stand 2026-10-10, abgelöst durch den Pin auf `v0.86.1` oben; DoD-Punkt 1,
Entscheidung des Nutzers vom 2026-10-10).** Am 2026-10-10:

- *Digest gehört zum Tag:* `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.86.0`
  meldet einen OCI-Index mit `Digest: sha256:d90200e94db311a70f9b15045e1290feecde03e9185ad78fb1ebea265e6ee44b`
  (Manifeste `linux/amd64` `sha256:7fbd5a20…`, `linux/arm64` `sha256:59e52bac…`);
  `docker pull` desselben Tags, danach `docker image inspect --format '{{json .RepoDigests}}'`,
  meldet denselben Digest. Gleich dem Digest der Release-Notiz.
- *Ziele und Aufruf unverändert:* `--print-mk` von `v0.85.0` (`c07f1fe6…`) und `v0.86.0`
  (`d90200e9…`) unterscheiden sich nur in der Zeile `DCHECK_IMAGE` (`diff`, eine Zeile); die
  Rückführung nach §4 tritt nicht ein. Geändert sind in `d-check.mk` nur `DCHECK_IMAGE` und
  `DCHECK_DIGEST`.
- *Changelog `0.86.0`:* geändert ist `reviews.match: name`, neu `skip-allows-empty`; beide
  wirken hier nicht — `.d-check.yml` führt weder das Modul `reviews` (`modules: [links,
  anchors, ids, matrix, spans]`) noch einen Schlüssel `reviews` oder `skip-allows-empty`
  (`grep`).
- *Weitere Nennungen des alten Pins:* `grep` nach `v0.85.0` und `c07f1fe6` außerhalb von
  `.git`, `.harness/` und `docs/reviews/` findet ihn in `d-check.mk` nicht mehr, sonst nur in
  diesem Plan (Belege und Herkunft), im Drift-Log der Roadmap (Zeile vom 2026-10-09, Anlass
  der Anlage) und in `open/slice-harness-lh-links-pflicht.md` (Messungen unter `v0.85.0`,
  mit Version genannt) — Zeitdokumente bzw. Messbefunde mit Stand, bleiben stehen. Das
  Heben des Pins ist keine Umplanung; das Drift-Log bekommt keine Zeile.
- *`make docs-check` am Arbeitsbaum:* `d-check: 439 Datei(en) geprüft, 0 Befund(e)` mit
  Digest `d90200e9…`.

Gegenprobe gegen `v0.86.0` (`d90200e9…`) in Kopien des Arbeitsbaums: je Mutation eine frische
Kopie im Scratchpad, jeder Eintrag der obersten Ebene außer `.git` mit `cp -r` ohne `-p`,
Mutation per `printf` angehängt, Lauf `docker run --rm --network none` mit dem Digest, danach
nur diese Kopie gelöscht. Unveränderte Kopie: `439 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

| Zusage | Mutation | Befundzeile (`v0.86.0`) |
|---|---|---|
| toter Anker ist rot (`anchors`) | Link auf `conventions.md#gibt-es-nicht`, an `harness/README.md` | `harness/README.md:141 conventions.md#gibt-es-nicht anchor-missing`, 1 Befund, Exit 1 |
| Kennung `ADR-` ohne Link ist rot (`ids`) | `ADR-` mit vier Ziffern (0001) blank im Fließtext, an `harness/README.md` | `harness/README.md:141` · die Kennung · `id-unlinked`, 1 Befund, Exit 1 |
| totes Linkziel ist rot (`links`) | Link auf `gibt-es-nicht.md`, an `harness/README.md` | `harness/README.md:141 gibt-es-nicht.md target-missing`, 1 Befund, Exit 1 |
| Spec-Stratum nennt keine ADR (`matrix`) | Link aus `spec/architecture.md` auf ADR 0001 | `spec/architecture.md:645 … matrix-forbidden Referenz spec-straten → adr`, 1 Befund, Exit 1 |

**Pin auf `v0.85.0` (Stand 2026-10-09, abgelöst durch den Pin auf `v0.86.0` oben).** Am 2026-10-09:

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

**Messbefund zum Block *Strenges Doc-Gate*, Stand `v0.85.0`.** Seit der Entscheidung des
Nutzers vom 2026-10-10 sagt der Block kein rot/grün-Verhalten des Gates zu (DoD-Punkt 1). Die
folgenden Tabellen messen frühere Fassungen des Blocks, die solches Verhalten zusagten; sie
bleiben als Messbefund über d-check `v0.85.0` unter `.d-check.yml` stehen, nicht als Zusage,
und sind gegen `v0.86.0` und `v0.86.1` nicht neu gefahren.

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
ebenda (kein Befund). Ein Verweis aus `spec/` auf eine Welle ist nicht mutiert — der Block
sagt ihn nicht zu. Den Link auf eine superseded ADR sagt der Block seit der Nacharbeit
unten wieder zu.

**Nacharbeit nach Review F-555, F-557 und F-561** (Block in den drei Commands und DoD-Punkt
1), am Stand `6db272b` mit den geänderten Commands im Arbeitsbaum, nur gegen `v0.85.0`
(`c07f1fe6…`). Kopien wie oben (je Mutation frisch, `cp -r` ohne `-p`, ohne `.git`, danach
gelöscht); unveränderte Kopie: `438 Datei(en) geprüft, 0 Befund(e)`, Exit 0.

| Zusage im Block | Mutation | Befundzeile |
|---|---|---|
| nicht zugesagt: Kennung in umzäuntem Code-Block (F-555) | `ADR-` mit vier Ziffern (0001) in einem Block mit drei Backticks, angehängt an `harness/README.md` | **kein Befund**, 0 Befunde, Exit 0 |
| dasselbe, Zaun `~~~` | wie davor, Zaun mit drei Tilden | **kein Befund**, 0 Befunde, Exit 0 |
| dasselbe, auch in `spec/` | im Block mit drei Backticks je eine Kopie: `ADR-` (0001) in `spec/architecture.md`, `MR-001` in `spec/spezifikation.md`, `slice-harness-d-check-v0-85` in `spec/architecture.md` | je **kein Befund**, 0 Befunde, Exit 0 |
| `ADR-` ohne Link ist rot in Inline-Code (Gegenstück) | `ADR-` (0001) in Backticks, angehängt an `harness/README.md` | `harness/README.md:141` · die Kennung · `id-unlinked`, 1 Befund, Exit 1 |
| ADR-Kennung ohne Link in `spec/` meldet `id-unlinked` (F-557) | `ADR-` (0001) in Backticks, angehängt an `spec/architecture.md` | `spec/architecture.md:645` · die Kennung · `id-unlinked`, 1 Befund, Exit 1 — kein `matrix-forbidden` |
| Slice-Kennung in `spec/` ist rot, in Inline-Code | `slice-harness-d-check-v0-85` in Backticks, angehängt an `spec/architecture.md` | `spec/architecture.md:645 slice- matrix-forbidden … spec-straten → slice`, 1 Befund, Exit 1 |
| Link auf eine superseded ADR ist rot (F-561) | Link auf die ADR 0021, angehängt an `harness/README.md` | `harness/README.md:141 ../docs/plan/adr/0021-vergleichsregeln-beim-einspielen.md matrix-inactive … Superseded by` (Link auf 0023), 1 Befund, Exit 1 |
| Ausweg: auf die ersetzende ADR verlinken | Link auf die ADR 0023, angehängt an `harness/README.md` | **kein Befund**, 0 Befunde, Exit 0 |
| Ausnahme `docs/reviews/` | Link auf die ADR 0021, angehängt an den Review-Report dieses Slice | **kein Befund**, 0 Befunde, Exit 0 |

Zur Ausweg-Zeile: Die ADR 0022 ist selbst superseded (ein Link auf sie meldet ebenfalls
`matrix-inactive` mit Verweis auf 0023); der Block sagt deshalb „die ersetzende ADR“, nicht
„die nächste“. Die Ausnahme des ADR-Index ist nicht mutiert, belegt ist sie am Bestand:
`docs/plan/adr/README.md` verlinkt die ADR 0021 (Zeile 27), und die unveränderte Kopie meldet
0 Befunde. Nicht zugesagt: eine ADR mit Status `deprecated` (das Repo führt keine);
ein blanker Verweis auf eine superseded ADR meldet `id-unlinked`, nicht `matrix-inactive`
(eigene Kopie, 1 Befund, Exit 1), das ist die Zeile `ADR-` ohne Link.

**Nacharbeit nach Verifikation V-134** (Block in den drei Commands und DoD-Punkt 1), am Stand
`9293191` mit den geänderten Commands im Arbeitsbaum, nur gegen `v0.85.0` (`c07f1fe6…`). Die
Ausnahmen sind aus `.d-check.yml` abgeleitet (`scan.ignore`, Ziel des `ids`-Musters,
`exempt-paths`, `exclude-sections`, Zeilen-Marker aus den Kommentaren dort) und je
Markdown-Form (Überschrift, Code-Block, Inline-Code, Tabelle, Zitat, Link- und Bildtext,
Fußnote, HTML-Kommentar und -Block) am Verhalten von d-check gemessen. Kopien wie oben: je
Mutation frisch im Scratchpad, jeder Eintrag der obersten Ebene außer `.git` mit `cp -r` ohne
`-p`, Mutation per `printf` angehängt (an `harness/README.md`, sofern nicht anders genannt),
Lauf `docker run --rm --network none` mit dem Digest aus `d-check.mk`, danach nur diese Kopie
gelöscht. Unveränderte Kopie: `439 Datei(en) geprüft, 0 Befund(e)`, Exit 0. Rot heißt hier
1 Befund, Exit 1; grün 0 Befunde, Exit 0.

| Zusage oder Ausnahme im Block | Mutation | Ergebnis |
|---|---|---|
| rot: `ADR-` ohne Link im Absatz | `ADR-` (0001) blank in einem Absatz | `harness/README.md:141` · die Kennung · `id-unlinked`, rot |
| rot: im Listenpunkt | dasselbe als `- …` | dieselbe Zeile, rot |
| rot: in der Tabellenzelle | Tabelle mit der Kennung in einer Zelle | `harness/README.md:143` · `id-unlinked`, rot |
| rot: im Zitat | dasselbe als `> …` | `harness/README.md:141` · `id-unlinked`, rot |
| rot: in Inline-Code darin | die Kennung in Backticks | `harness/README.md:141` · `id-unlinked`, rot |
| rot: in `spec/` ebenso `id-unlinked` | die Kennung in Backticks, an `spec/architecture.md` | `spec/architecture.md:645` · `id-unlinked`, rot |
| rot: gescannt sind auch die Commands | die Kennung blank, an `.claude/commands/plan-welle.md` | `.claude/commands/plan-welle.md:119` · `id-unlinked`, rot |
| nicht: `#`-Überschrift | `## Probe` mit der Kennung | grün |
| Gegenkontrolle: nur `#`-Überschriften | dieselbe Überschrift in Setext-Form (Unterstrich `---`) | `harness/README.md:141` · `id-unlinked`, rot |
| nicht: Text eines Links, beliebiges Ziel | die Kennung als Link-Text, Ziel `conventions.md` | grün |
| nicht: Text eines Bildes | die Kennung als Alt-Text, Ziel `conventions.md` | grün |
| nicht: Zeile mit `<!-- d-check:ignore … -->` | die Kennung blank, Marker in derselben Zeile | grün |
| nicht: Datei unter `docs/plan/adr/`, ADR | die Kennung blank, an ADR 0002 | grün |
| nicht: dasselbe, ADR-Index | die Kennung blank, an `docs/plan/adr/README.md` | grün |
| nicht: dasselbe, jede Datei darunter | neue Datei `docs/plan/adr/x/notiz.md` mit der Kennung blank | grün (440 Dateien geprüft) |
| nicht gescannt: `.harness/**`, auch die Skills | die Kennung blank, an `.harness/skills/reviewer.md` | grün |
| nicht gescannt: `**/*.template.md` | neue Datei `harness/probe.template.md` mit der Kennung blank | grün, 439 Dateien geprüft |
| nicht gescannt: `.tmp/**` | neue Datei `.tmp/probe.md` mit der Kennung blank | grün, 439 Dateien geprüft |
| nicht: umzäunter Code-Block, drei Backticks | die Kennung in einem Zaun | grün |
| nicht: dasselbe, `~~~` | die Kennung in einem Tilden-Zaun | grün |
| rot: totes relatives Linkziel | Link auf `gibt-es-nicht.md` | `harness/README.md:141 gibt-es-nicht.md target-missing`, rot |
| rot: toter Anker | Link auf `conventions.md#gibt-es-nicht` | `… anchor-missing`, rot |
| nicht: externe URL | Link auf eine `https://`-Adresse, die es nicht gibt | grün |
| nicht: HTML-Link | `<a href>` auf `gibt-es-nicht.md` | grün |
| nicht: Link in Inline-Code | der Link auf `gibt-es-nicht.md` in Backticks | grün |
| nicht: Link im Zaun | derselbe Link in einem Zaun, je mit drei Backticks und mit `~~~` | grün, grün |
| rot: Link aus `spec/` auf eine ADR | Link auf ADR 0001, an `spec/architecture.md` | `spec/architecture.md:645 … matrix-forbidden Referenz spec-straten → adr`, rot |
| rot: Slice-Kennung in `spec/`, blank | `slice-harness-d-check-v0-85` blank, an `spec/architecture.md` | `spec/architecture.md:645 slice- matrix-forbidden`, rot |
| rot: dasselbe in Inline-Code | in Backticks | dieselbe Zeile, rot |
| rot: `MR-` in `spec/`, blank | `MR-001` blank, an `spec/spezifikation.md` | `spec/spezifikation.md:2668 MR-001 matrix-forbidden`, rot |
| rot: dasselbe in Inline-Code | in Backticks | dieselbe Zeile, rot |
| nicht: Abschnitt `Geschichte`, Slice-Kennung | `## Geschichte`, darunter die Slice-Kennung blank, an `spec/architecture.md` | grün |
| nicht: dasselbe, `MR-` | `## Geschichte`, darunter `MR-001`, an `spec/spezifikation.md` | grün |
| nicht: dasselbe, Link auf eine ADR | `## Geschichte`, darunter Link auf ADR 0001, an `spec/architecture.md` | grün |
| Gegenkontrolle: der Abschnitt endet | nach `## Geschichte` ein `## Danach` mit der Slice-Kennung | `spec/architecture.md:651 slice- matrix-forbidden`, rot |
| Gegenkontrolle: „genau `Geschichte`“ | `## 9. Geschichte`, darunter die Slice-Kennung | `spec/architecture.md:647 slice- matrix-forbidden`, rot |
| nicht: Slice- oder `MR-`-Kennung in einer Zeile mit `<!-- d-check:status-provenance -->` | die Slice-Kennung blank mit dem Marker, an `spec/architecture.md` | grün |
| Gegenkontrolle: der Marker nimmt keinen Link aus | Link auf ADR 0001 mit dem Marker, an `spec/architecture.md` | `… matrix-forbidden Referenz spec-straten → adr`, rot |
| nicht: Zaun in `spec/` | die Slice-Kennung in einem Zaun, an `spec/architecture.md`, je mit drei Backticks und mit `~~~`; dazu Link auf ADR 0001 in einem Tilden-Zaun | grün, grün, grün |
| rot: Link auf eine superseded ADR | Link auf ADR 0021 | `harness/README.md:141 … 0021-… matrix-inactive`, rot |
| Ausweg: die ersetzende ADR | Link auf ADR 0023 | grün |
| nicht: `docs/reviews/` | Link auf ADR 0021, an den Review-Report dieses Slice | grün |
| nicht: ADR-Index | Link auf ADR 0021, an `docs/plan/adr/README.md` | grün |
| nicht: aus einer ADR | Link auf ADR 0021, an ADR 0010 | grün |
| Gegenkontrolle: aus einer ADR, nicht aus dem Verzeichnis | Link auf ADR 0021 in neuer Datei `docs/plan/adr/notiz.md` | `docs/plan/adr/notiz.md:3 … matrix-inactive`, rot |
| nicht: Abschnitt `Geschichte` | `## Geschichte`, darunter Link auf ADR 0021 | grün |
| nicht: Zaun | Link auf ADR 0021 in einem Zaun, je mit drei Backticks und mit `~~~` | grün, grün |
| nicht: `LH-` ohne Link | `LH-FA-01` blank | grün |
| nicht: `MR-` außerhalb von `spec/` | `MR-001` blank | grün |
| nicht: Pfade in Inline-Code | `` `tools/gibt-es-nicht.sh` `` | grün |

Was der Block nicht nennt und rot ist, bleibt ungenannt (er sagt weniger zu, als das Gate
prüft): Setext-Überschrift (oben), Fußnote, HTML-Kommentar und -Block, Titel einer
Link-Referenz, Autolink, eingerückter Code-Block (je die `ADR-`-Kennung blank, `id-unlinked`,
eigene Kopien); totes Linkziel in einer Überschrift, Tabelle, Bild, Link-Referenz,
HTML-Kommentar, in `docs/plan/adr/` und in einem Abschnitt `Geschichte` (`target-missing`);
Slice-Kennung in `spec/` in Überschrift, Tabelle, HTML-Kommentar, eingerücktem Code-Block
(`matrix-forbidden`); `## Geschichte` nimmt `id-unlinked` nicht aus (Verifikation, Abschnitt 1).
`<!-- d-check:ignore … -->` nimmt weder `target-missing`, `anchor-missing`, `matrix-forbidden`
noch `matrix-inactive` aus, `<!-- d-check:status-provenance -->` weder `id-unlinked` noch
`matrix-inactive` (je eigene Kopie, rot); der Block nennt jeden Marker nur bei dem Befund, den
er ausnimmt. Eine Datei, die nicht auf `.md` endet, zählt d-check nicht (`harness/probe.txt`:
439 Dateien, grün).

*Fundstellen:* `grep -rn -i "Strenges Doc-Gate\|klickbare\|Anker-Link\|codepaths"` über
`.claude/`, `.harness/skills/` und `AGENTS.md` findet den Block nur in den drei Commands; dazu
in `.claude/commands/plan-welle.md` Schritt 8 „Kennungen als Anker-Links“ — eine Regel an den
Planner, keine Zusage über das Gate, bleibt. `.claude/agents/*.md` und die Skills führen
keine gleichlautende Stelle.

**Nacharbeit nach Verifikation V-138 und V-139** (Block in den drei Commands auf eine
geschlossene Zusage gefasst, DoD-Punkt 1 und §1 *Berichtigung*), am Stand `a08d0b9` mit den
geänderten Commands und diesem Plan im Arbeitsbaum, nur gegen `v0.85.0` (`c07f1fe6…`). Die
Tabellen davor messen frühere Fassungen des Blocks; ihre Ausnahmen gelten seither als
Messbefund, nicht als Zusage. Kopien wie oben: je Mutation frisch im Scratchpad, jeder Eintrag
der obersten Ebene außer `.git` mit `cp -r` ohne `-p`, Mutation per `printf` angehängt bzw. als
neue Datei, Lauf `docker run --rm --network none` mit dem Digest aus `d-check.mk`, danach nur
diese Kopie gelöscht; 96 Kopien, je sechs parallel. Unveränderte Kopie: `439 Datei(en)
geprüft, 0 Befund(e)`, Exit 0. Rot heißt mindestens 1 Befund, Exit 1; grün 0 Befunde, Exit 0.
*Die Kennung* ist `ADR-` mit vier Ziffern (0001), *die Slice-Kennung* die dieses Slice, *ein
toter Link* ein Inline-Link mit Text `x`; *Absatz* heißt angehängt als eigener Absatz.

| Zusage im Block | Mutation | Ergebnis |
|---|---|---|
| rot: die Kennung als Wort im Absatz, je Ort | angehängt an `AGENTS.md`, `spec/lastenheft.md`, `spec/spezifikation.md`, `spec/architecture.md`, `harness/README.md`, `harness/conventions.md`, `docs/plan/planning/README.md`, `welle-v1-abschluss.md`, je einen Plan in `open/`, `next/`, `in-progress/` (dieser), `in-progress/roadmap.md`, je einen Slice und eine Results-Notiz in `done/` (14 Kopien) | je `id-unlinked` an der angehängten Zeile, rot |
| dasselbe in einer neuen Datei | `harness/probe.md` und `open/slice-probe.md` neu (2 Kopien) | je `id-unlinked`, 440 Dateien, rot |
| rot: im Listenpunkt | `- `, `1. `, verschachtelt unter `- `, in `harness/README.md`; `- ` in einem Slice in `done/` (4 Kopien) | je `id-unlinked`, rot |
| rot: in Inline-Code | im Absatz und im Listenpunkt von `harness/README.md`, im Absatz von `spec/architecture.md` (3 Kopien) | je `id-unlinked`, rot |
| rot: als Wort, mit Satzzeichen und Auszeichnung | in Klammern mit Komma, fett, auf der Folgezeile eines Absatzes, nach einem Link mit derselben Kennung als Text, in einer Zeile mit `<br>`, als Pfadteil einer nackten URL, eine superseded ADR (0021) blank (7 Kopien) | je `id-unlinked`, rot |
| rot: toter Link auf eine `.md` | aus `harness/README.md`: `gibt-es-nicht.md` im Absatz und im Listenpunkt, mit Titel, `./…`, `../…`, `sensors/…`, `Conventions.md` (Groß-/Kleinschreibung), Leerzeichen als `%20`, Ziel in `vendor/`, `docs/vendor/`, `.harness/`, `.tmp/`, `*.template.md`; aus `AGENTS.md`, einem Slice in `done/`, `spec/architecture.md` (auch unter `## Geschichte`) (17 Kopien) | je `target-missing`, aus `spec/` dazu `matrix-forbidden … → aussen`, rot |
| rot: dasselbe außerhalb des Repos | `../../../gibt-es-nicht.md` | `repo-escape`, rot |
| rot: toter Anker in einer `.md` der genannten Orte | aus `harness/README.md` auf `conventions.md`, `AGENTS.md`, `spec/lastenheft.md`, `planning/README.md`, `in-progress/roadmap.md`, einen Slice in `done/`, dazu ein Anker in Großschreibung eines vorhandenen Slugs; aus einem Plan in `open/` auf diesen Plan; aus `spec/spezifikation.md` auf `lastenheft.md` (9 Kopien) | je `anchor-missing`, rot |
| rot (`implement-slice.md`): Slice-Kennung in `spec/` | blank im Absatz je Datei unter `spec/`, in Inline-Code, im Listenpunkt; nach einem beendeten Abschnitt `Geschichte` (`## Danach`); unter `## 9. Geschichte` (7 Kopien) | je `matrix-forbidden … → slice`, rot |
| rot (`implement-slice.md`): `MR-` mit drei Ziffern in `spec/` | `MR-001` blank je Datei unter `spec/`, in Inline-Code im Listenpunkt, fett (5 Kopien) | je `matrix-forbidden … → adaptionsblock`, rot |
| rot (`implement-slice.md`): Link aus `spec/` auf eine ADR | Link auf ADR 0001 je Datei unter `spec/`, mit vorhandenem Anker, Link auf ADR 0021 im Listenpunkt (5 Kopien) | je `matrix-forbidden … → adr`, rot |

**Abgleich mit den Formen aus V-134 und V-138** — keine fällt unter die Zusage; je eine Kopie,
gemessen zur Kontrolle:

| Form | Ergebnis | Was der Wortlaut ausnimmt |
|---|---|---|
| die Kennung im Titel eines Links | grün | „nicht in einem Link oder Bild“ |
| die Kennung als Text eines Bildes; die Slice-Kennung als Link-Text und im Link-Titel in `spec/`; `MR-001` als Link-Text in `spec/` | je grün | dasselbe |
| `<a href>` auf `gibt-es-nicht.md`; `<a href>` aus `spec/` auf ADR 0001 | je grün | „Markdown-Link der Form `[Text](Pfad)`“ |
| externe URL auf ein totes Ziel; externe URL aus `spec/` auf eine ADR | je grün | „mit relativem Pfad“ |
| Anker in ein Ziel ohne `.md` (`../go.mod#L1`) | grün | „auf eine `.md`-Datei“ |
| toter Anker in `a.txt`; Anker in einen Slice unter `done/<welle>/` | je rot | nicht zugesagt (Ziel ohne `.md`, kein genannter Ort) — das Gate prüft mehr |
| `#`-Überschrift mit der Kennung | grün | „in einem Absatz oder Listenpunkt“ |
| die Kennung im umzäunten Code-Block | grün | dasselbe |
| toter Link in Inline-Code | grün | „nicht in Inline-Code“ |
| Zeile mit `<!-- d-check:ignore … -->`; Slice-Kennung in `spec/` mit `<!-- d-check:status-provenance -->` | je grün | „in einer Zeile ohne HTML-Kommentar“ |
| `## Geschichte` in `spec/architecture.md` mit der Slice-Kennung | grün | „außerhalb eines Abschnitts mit der Überschrift `Geschichte`“ |
| die Kennung in ADR 0002, in `.harness/skills/reviewer.md`, in `harness/vendor/probe.md`, `docs/vendor/probe.md` (V-139) | je grün, 439 Dateien | Orte: `AGENTS.md` und `.md` **direkt** unter den genannten Verzeichnissen |
| Link auf eine superseded ADR außerhalb von `spec/` (`matrix-inactive`) | rot (Tabelle zur Nacharbeit nach Review F-561) | nicht mehr zugesagt |

*Fundstellen:* `grep -rn "Strenges Doc-Gate"` über `.claude/`, `.harness/skills/` und
`AGENTS.md` findet den Block nur in den drei Commands.

**Messbefund der Verifikation V-141 bis V-144** (Schlussprüfung `8d79094`, Stand `v0.85.0`,
gemessen im Verifikationsbericht): grün blieben ein Markdown-Link, dessen Text über einen
Zeilenumbruch reicht (totes Ziel, toter Anker, Link aus `spec/` auf eine ADR; V-141), die
Kennung `ADR-` in Inline-Code über einen Zeilenumbruch (V-142), Abschnitt und Anker nach einer
Zeile `    ## …` in einem eingerückten Code-Block (V-143), ein Anker, den es nur als `<a id>`
in einem HTML-Kommentar gibt (V-144). Grenzen des Gates, keine Zusage des Blocks mehr.

**Block ohne Zusage (Entscheidung vom 2026-10-10).** Der Block in den drei Commands nennt die
Regel, Kennungen (`LH-`, `ADR-`, `MR-`) als Anker-Links zu schreiben, dass `LH-` ohne Link
noch nicht erzwungen ist (`slice-harness-lh-links-pflicht`), und dass `.d-check.yml` festlegt,
was das Gate prüft, und `make docs-check` maßgeblich ist. Er sagt kein Verhalten des Gates zu,
deshalb gibt es zu ihm keine Mutation (`AGENTS.md` §3.10 greift nur auf Zusagen); die Pflicht
für `ADR-` belegt die Gegenprobe der Klasse `ids` oben. Inline-Code im Block steht je in einer
Zeile. *Fundstellen:* `grep -rn "Strenges Doc-Gate"` über `.claude/`, `.harness/skills/` und
`AGENTS.md` findet den Block nur in den drei Commands, dort nach dem Titel je mit demselben
Wortlaut.

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
mit Digest `c07f1fe6…`, Integrationstests, alle Gegenproben grün.

Letzter Lauf (Verifikation V-136): `make gates` am Stand `6565310` (Nacharbeit nach V-134 bis
V-137 committet), Arbeitsbaum sauber: Exit 0; darin `baseline-verify: v6.16.0 OK — 54 Dateien`,
`d-check: 439 Datei(en) geprüft, 0 Befund(e)`, `run-integration-tests: gruen`, je `gruen`
`a-check-negativ`, `commit-msg-gegenprobe`, `abdeckung-gegenprobe`, `kopf-check-gegenprobe`,
`lint-gegenprobe`. Der Commit danach trägt nur diesen Absatz.

Letzter Lauf (Nacharbeit nach V-138 und V-139): `make gates` am Stand `b5e4437`, Arbeitsbaum
sauber: Exit 0; darin `baseline-verify: v6.16.0 OK — 54 Dateien`, `d-check: 439 Datei(en)
geprüft, 0 Befund(e)`, `run-integration-tests: gruen`, je `gruen` `a-check-negativ`,
`commit-msg-gegenprobe`, `abdeckung-gegenprobe`, `kopf-check-gegenprobe`, `lint-gegenprobe`.
Der Commit danach trägt nur diesen Absatz.

Letzter Lauf (Pin auf `v0.86.0`, Block ohne Zusage): `make gates` am Stand `fda566b`,
Arbeitsbaum sauber: Exit 0; darin `baseline-verify: v6.16.0 OK — 54 Dateien`, `d-check: 439
Datei(en) geprüft, 0 Befund(e)` mit Digest `d90200e9…`, `run-integration-tests: gruen`, je
`gruen` `a-check-negativ`, `commit-msg-gegenprobe`, `abdeckung-gegenprobe`,
`kopf-check-gegenprobe`, `lint-gegenprobe`. Der Commit danach trägt nur diesen Absatz.

Letzter Lauf (Pin auf `v0.86.1`): `make gates` am Stand `006c559`, Arbeitsbaum sauber: Exit 0;
darin `baseline-verify: v6.16.0 OK — 54 Dateien`, `d-check: 439 Datei(en) geprüft, 0
Befund(e)` mit Digest `3e0b9779…`, `run-integration-tests: gruen`, je `gruen`
`a-check-negativ`, `commit-msg-gegenprobe`, `abdeckung-gegenprobe`, `kopf-check-gegenprobe`,
`lint-gegenprobe`. Der Commit danach trägt nur diesen Absatz.

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
  3. Nachgezählt beim Anlegen: drei Liefer-Punkte; der Plan selbst zählt nicht.
  Nachgezählt bei der Berichtigung vom 2026-10-09: weiter drei Liefer-Punkte (die Fassung
  der Commands gehört zu DoD-Punkt 1, sie sagt das gegenprobierte Verhalten desselben Gates
  aus). Schichten nach der Schichtteilung unten, neu am 2026-10-09 (Review F-556): eine,
  der Harness (`d-check.mk`, Skill, Commands); die Zählung bei der Berichtigung trennte
  Gate-Konfiguration und Agenten-Anweisungen und kam auf zwei. Der Nehmer der
  Teil-Zusage, `slice-harness-lh-links-pflicht`, ist dort nach derselben Teilung
  nachgezählt und geschnitten; der abgeschnittene Teil `slice-harness-lh-links-bestand`
  ebenso (je §8).

Keiner der Einträge erreicht mit diesem Plan die Schwelle 3× neu.

**Schichtteilung** (gleich in `slice-harness-lh-links-pflicht`, `slice-harness-lh-links-bestand`
und diesem Slice; festgelegt am 2026-10-09 nach Review F-556): Eine Schicht im Sinn von
Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice ist ein Bereich mit
eigenem Maßstab im Review — die Spezifikation (`spec/`) · die Entscheidungen
(`docs/plan/adr/`) · der Harness (Gate-Konfiguration und Gate-Werkzeuge, `Dockerfile`,
Agenten-Anweisungen: `AGENTS.md`, Skills, Commands, `harness/`) · die Nutzer- und
Wartungs-Doku (`docs/user/`, `docs/maintainer/`) · im Produkt-Code je Schicht des Hexagons,
ebenso deren Tests. So zählten `slice-lastenheft-pruefbarkeit` (Harness und Lastenheft mit
Spezifikation als verschiedene Schichten) und `slice-harness-lint-werkzeug` (Harness mit
`Dockerfile` als eine Schicht, Tests je Schicht des Hexagons). Keine Schicht ist die
Planung (Pläne, Roadmap, Register): In ihr schreibt jeder Slice (eigener Plan, Sendungen,
Drift-Log); was sie an Umfang trägt, misst das dritte Kriterium, eine Review-Sitzung.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
