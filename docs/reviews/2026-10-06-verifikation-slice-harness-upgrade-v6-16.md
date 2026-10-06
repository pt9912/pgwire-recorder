# Verifikation: slice-harness-upgrade-v6-16 — 2026-10-06

**Rolle:** Verifier (Modul 11). Ich prüfe, ob der Slice Plan und DoD erfüllt, also ob richtig gebaut wurde. Ob das Richtige gebaut wurde, prüft der Validator. Den Diff gegen Entscheidungen und Hard Rules hat der Reviewer geprüft.

**Gegenstand:** Slice-Plan `docs/plan/planning/in-progress/slice-harness-upgrade-v6-16.md` (Kopf, §1 bis §8) gegen den Gesamt-Diff `53ea519..481ef52`. Darin: `01f9933` (d-check `v0.82.0`), `0fca02e` (Baseline `v6.16.0` statt `v6.13.0`), `aec4aab` (Abgleich und Regelblock-Tabelle), `01a45ae` (Skills, ADR-Index, Folge-Slice, Drift-Log), `117e23f` ([Review](2026-10-06-review-slice-harness-upgrade-v6-16.md), F-434 bis F-440) und `481ef52` (Nacharbeit zu F-434, F-435, F-436, F-439, ohne eigenes Review).

**Eingang:**

- die DoD-Liefer-Punkte, §6 *Randformen* und *Risiken*, §7 *Belege* des Plans und die Commit-Messages
- der Review-Report und die Nacharbeit `481ef52`
- die Entscheidungen des Nutzers laut Auftrag: Der Folge-Slice `slice-harness-gate-index-werkzeug-teil` bleibt als Slice. Kein Change Request an das Kurs-Repo. Die Gate-ADRs bleiben unverändert. Ältere Werkzeuge kommen erst mit einer Randform in Spezifikation §11.
- `harness/conventions.md`, `AGENTS.md` §1 und §3.12, `.claude/agents/architect.md`, `.claude/commands/implement-slice.md` und `plan-welle.md`, `.harness/skills/*.md`, `harness/sensors/*.md`, `docs/plan/adr/README.md`, der ADR-Kopf von [ADR-0025](../plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md), [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0029](../plan/adr/0029-benannte-welle-kennungen-im-commit-hook.md), [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) und [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md), Spezifikation §11
- alter Baum aus `git archive 53ea519 .harness/baseline/v6.13.0`, neuer Baum `.harness/baseline/v6.16.0/`

Ein eigener Bericht des Implementers lag mir nicht vor. Die Belege, die die DoD „im Bericht“ verlangt, stehen in §7 *Belege* des Plans. Jede Aussage darunter habe ich in diesem Kontext selbst nachgefahren (V-79). Bei meinem Start war der Arbeitsbaum sauber und HEAD `481ef52`. Der Commit dieses Berichts enthält nur diese Datei.

**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

**Wie die Proben gebaut wurden.** Jede Probe lief in einer frischen Kopie aus `git archive <stand>` unter `mut/m<N>/` im Scratchpad, ohne `cp -p` und nie im Repo. Das Asset habe ich mit `gh release download` in ein Scratchpad-Verzeichnis geladen. Die Kopien sind weggeräumt.

---

## 1. DoD-Liefer-Punkte

### Punkt 1: d-check `v0.82.0`. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| `DCHECK_IMAGE` nennt `v0.82.0`, `DCHECK_DIGEST` den Digest des Tags | `d-check.mk` Zeile 7 und 8. Drei unabhängige Quellen nennen `sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`: `docker image inspect … --format '{{json .RepoDigests}}'` (lokal), `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.82.0` (Registry) und `gh release view v0.82.0 -R pt9912/d-check` („Digest-Pin: `ghcr.io/pt9912/d-check@sha256:d28e9437…`“) | bestätigt |
| Neu erzeugt, genau vier Adaptionen | `diff` gegen `docker run --network none ghcr.io/pt9912/d-check:v0.82.0 --print-mk`: Kopfkommentar, `DCHECK_DIGEST` gepinnt, `doc-check` → `docs-check`, Muster von `doc-help` `^docs?-`. Rezepte gleich | bestätigt |
| Der Kopf nennt die `doc-help`-Adaption (§3.11) | „das Muster von doc-help auf '^docs?-' (listet docs-check mit)“; Mutation `m6` unten | bestätigt |
| `make docs-check` grün, neue Befunde gezählt und mit Ausgang | Im Gate-Lauf 281 Dateien und 0 Befunde. Probe `m9`: Stand `53ea519` unter `v0.82.0` hat 0 Befunde, der Bump allein bringt also keinen. Probe `m8`: `0fca02e` mit der alten Fassung von `MR-001` liefert genau den einen Befund `target-missing` aus §6. Der Ausgang ist nach §6 (a) behoben. `git diff --word-diff=porcelain` zeigt in `MR-001` nur den Wechsel vom Link zum Code-Span; Pfad mit `v6.13.0` und Anker sind gleich geblieben | bestätigt |
| `.d-check.yml` unverändert, keine Lockerung | `git diff 53ea519 481ef52 -- .d-check.yml` leer | bestätigt |

### Punkt 2: Baseline `v6.16.0` vendored. **Bestätigt.**

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Asset-sha256 vor dem Entpacken geprüft | Unabhängig nachgeladen: `sha256sum lab-regelwerk.zip` ergibt `feb4d7444c92ec4d11fcf88035ce2eaef48da64ae990eb8bbcf44550a5e87063` und ist gleich der Zeile im Release-Asset `SHA256SUMS` des Kurs-Repos (Rang 1 aus §6 *Netzzugriff*) | bestätigt |
| Vendored Baum gleich dem Asset | Asset frisch entpackt: oberste Ebene nur `regelwerk/` und `templates/`. `diff -r` gegen `.harness/baseline/v6.16.0/{regelwerk,templates}` ist leer | bestätigt |
| `SHA256SUMS` vollständig, in der Form des Skripts | 54 Zeilen, `sha256sum -c` OK. Die Pfadliste ist gleich `find regelwerk templates ! -type d \| LC_ALL=C sort`, und die Liste ist sortiert | bestätigt |
| `make baseline-verify` grün und nennt `v6.16.0` | „baseline-verify: v6.16.0 OK — 54 Dateien“, einzeln und im Gate-Lauf | bestätigt |
| An keinem Commit zwei Tag-Verzeichnisse | `git ls-tree` über alle sieben Stände von `53ea519` bis `481ef52`: je genau ein Tag, Wechsel in `0fca02e` | bestätigt |
| Stand und Datum in §Baseline, Asset-URL in §Adoptierte Konventions-Quellen und `AGENTS.md` §1, Stand-Zeile und Messzeile | `v6.16.0` / 2026-10-06. URL an beiden Orten. Stand-Zeile „Kurs-Welle 159 · 2026-10-06“ ist gleich `regelwerk/README.md`. Messzeile `harness/conventions.md:102` | bestätigt |
| `grep -rn "v6\.13\.0"` über die lebenden Harness-Dateien leer | Gelaufen über `AGENTS.md`, `README.md`, `harness/` ohne `harness/conventions/done/`, `.claude/`, `.harness/skills/`, `tools/`, `Makefile`, `*.mk`, `harness/mk/`, `spec/`, `.d-check.yml`, `.githooks/`, `docs/plan/adr/README.md` und `docs/plan/planning/observations/`: leer. Einziger Treffer unter `harness/conventions/done/` ist `MR-001`, Code-Span nach §6 | bestätigt |
| Kein Plan in `open/`, `next/`, `in-progress/` verweist auf einen Pfad unter `.harness/baseline/v6.13.0/` | Treffer gibt es nur im eigenen Plan. Mein Urteil zu F-440 steht in V-78: Keine der Stellen ist ein Verweis | bestätigt (Urteil) |
| `find . -path ./.git -prune -o -xtype l -print` leer | leer; beide Symlinks unter `.claude/rules/` zeigen auf `v6.16.0` und lösen auf | bestätigt |

### Punkt 3: Abgleich v6.13.0 → v6.16.0. **Bestätigt mit einer Lücke** (V-76).

| Teilbehauptung | Beleg | Urteil |
|---|---|---|
| Jeder Regelblock steht in der Tabelle mit genau einem Wert | Gegen `ls .harness/baseline/v6.16.0/regelwerk/*.md` (26 Dateien) gezählt: Jede steht in der Tabelle. 23 Dateien haben eine Zeile, `grundlagen-harness-dateien.md`, `modul-02-…` und `modul-15-…` je zwei Zeilen nach Abschnitt. Das sind 29 Zeilen. Keine Zeile nennt eine Datei, die es im Baum nicht gibt. Die Werte stammen alle aus der geschlossenen Menge: 18 × *Träger kommt mit*, 8 × *liegt bei, nicht verdrahtet* und 3 × *kommt nicht mit*, jede davon mit Grund und Dauer | bestätigt |
| Keine Datei neu, umbenannt oder entfallen | `diff -r` alt gegen neu: dieselben 54 Dateien | bestätigt |
| Inhaltliches Delta wie in §7 *Belege* | `diff -r -I 'v6\.1[36]\.0'`: Inhalt ändert sich genau in `README.md`, `grundlagen-begriffe.md`, `grundlagen-harness-dateien.md`, `grundlagen-referenz-richtung.md`, `modul-03-spec.md`, `modul-10-review-harness.md`, `modul-13-quality-gates.md` und `modul-15-observability.md`, dazu in 14 Vorlagen | bestätigt |
| Delta gegen verkörperte Regeln gehalten, je Befund ein Ausgang | Stichprobe in Abschnitt 3. Befund (1) Reviewer-Regeln ist in beiden Skills nachgezogen, (3) Gate-Index-Teil ist an den Folge-Slice gegangen. Bei (4) und (5) ist „kein Befund“ belegt. Zu Befund (2) Werkzeug-Festlegungen im Technik-Stratum gibt §7 einen falschen Grund an, und ein Teilbefund hat keinen Ausgang (V-76) | **Lücke V-76** |
| Adaptions-Durchgang mit `MR-000` und `MR-001` | Keine aktive Adaption. Die Zeile `MR-001` in §Aufgelöste Adaptionen nennt `v6.16.0` und `grundlagen-referenz-richtung.md` §Spec-Straten. Den Abschnitt gibt es dort (Delta Zeile 289 bis 305). Die Datei in `done/` hat nur die Form-Reparatur bekommen | bestätigt |

### Übrige DoD-Punkte (konstant je Slice)

| Punkt | Stand |
|---|---|
| `make gates` grün | bestätigt, Abschnitt 5 |
| Review mit Report | liegt vor (F-434 bis F-440). Die Nacharbeit `481ef52` hat kein eigenes Review; ihre Wirkung prüft Abschnitt 2 |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | offen. §7 trägt Platzhalter, und die vier Risiken in §6 stehen auf „— (bei Closure)“. Das erledigt die Closure |

---

## 2. Review-Findings

| Finding | Status | Beleg |
|---|---|---|
| F-434 (MEDIUM, §3.12 lässt die ADR als Ort zu) | **behoben** | `AGENTS.md` §3.12 nennt nur die Spezifikation und die ADR „nur Entscheidung und Gründe … `Schärft:` … (§3.8)“. `architect.md` Zeile 25 („an dem Ort, den §3.12 nennt“), `implement-slice.md` und `plan-welle.md` zeigen nur auf §3.12, und keine lebende Datei unter `.claude/` und `.harness/skills/` nennt die ADR noch als Ort einer Randform (`grep "oder einer ADR"` leer). Plan §3 Zeile `AGENTS.md` ist nachgezogen |
| F-435 (LOW, Zelle mit zwei Werten) | **behoben** | zwei Zeilen; die zweite mit *kommt nicht mit*, Grund und Dauer („bis ai-harness-init den Teil schreibt“) |
| F-436 (LOW, verwaiste Klammer) | **behoben** | „adoptierter Stand (Stand-Zeile in … Kurs-Repo);“ |
| F-437 (INFO, §3.3) | ohne Aktion, richtig so | Hinweis an den Architect, kein Implementer-Auftrag |
| F-438 (INFO, Lesart unter `MR-000`, Folge-Slice) | **entschieden, nicht festgehalten** | Laut Auftrag hat der Nutzer entschieden, dass der Folge-Slice bleibt. Weder der Plan noch der Folge-Slice halten das fest (V-77) |
| F-439 (INFO, Report-Gerüst) | **behoben** | Der Skill nennt `.harness/baseline/<tag>/templates/docs/reviews/review-report.template.md`, und die Datei liegt im Baum |
| F-440 (INFO, an den Verifier) | **beurteilt** | V-78 |

---

## 3. Plan gegen Code und Stichprobe des Abgleichs

- **Kopf:** `Bezug` (MR-000) und `Berührte Spec-Stellen: —` stimmen. Spezifikation und Lastenheft sind im Diff unverändert, und `make kopf-check` ist im Gate-Lauf grün.
- **§1:** Die Abgrenzung hält. Im Diff gibt es keine `*.go`-Datei und keine Änderung an `Dockerfile`, `tools/` oder `.d-check.yml`. Eingefrorene Dokumente sind bis auf die Form-Reparatur in `MR-001` unberührt.
- **§3:** Jede Datei im Diff außerhalb des Baums steht in einer Zeile von §3. Der Review-Report gehört zur DoD. Die Zeilen `harness/README.md` und `harness/sensors/*.md` sind „prüfen“, und beide Dateien sind unverändert. Die Commit-Folge (1) bis (4) entspricht dem *Ansatz*.
- **§6 *Randformen*:** Jede Entscheidung ist so umgesetzt, wie sie dasteht (Abschnitt 1). §6 hält auch das Negativ fest, dass kein Gate einen hängenden Symlink sieht. Probe `m5` zeigt es: `make docs-check` meldet 0 Befunde. Der Ort der Randformen ist dieser Plan und nicht die Spezifikation. Das widerspricht dem neuen §3.12 nicht, denn §3.12 bindet nur einen Slice, der einen neuen Vertrag liefert, und dieser liefert keinen.
- **Stichprobe Delta → verkörperte Regeln** (über die Befundliste in §7 hinaus nachgelesen):
  - `modul-10` „Zuordnung zur Beobachtung“ statt `BEO-<NNN>`: Keine lebende Datei unter `.harness/skills/`, `.claude/` und `harness/` nennt `BEO-<NNN>` oder `observations.md`. Ohne Befund.
  - `grundlagen-begriffe.md` und `gate.template.md`: Eine Sensor-Datei verlinkt die Spec-Kennung und trägt keine Schwelle und keine Randform. `harness/sensors/kopf-check.md` und `abdeckung-check.md` verlinken `SPEC-047` bzw. `SPEC-048` und nennen deren Grenze nur als Zeiger. Ohne Befund.
  - `harness/README.template.md`: „Spec-Kennung“ ist jetzt in der Bindung zulässig, und die Zeilen von `kopf-check` und `abdeckung-check` nennen sie in der Spalte Vertrag. Ohne Befund. Was den Teil des Werkzeugs betrifft, liegt beim Folge-Slice.
  - `spezifikation.template.md` §7 *Festlegungen der Harness-Werkzeuge*: Spezifikation §11 trägt diese Rolle. Sie führt nur `SPEC-047` bis `SPEC-049`. Für Commit-Hook, `baseline-verify` und `a-check-negativ` gibt es keine Festlegung in §11; siehe V-76.
  - `adr/NNNN-titel.template.md` und `adr/README.template.md`: „Seine ADR schärft sie.“ Der ADR-Index nennt die Regel jetzt in *Konventionen*. [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) und [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) tragen aber `**Schärft:** —`, obwohl `SPEC-048` und `SPEC-047` ihre Spec-Stellen sind. Nur [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) schärft `SPEC-049`. Siehe V-76.

---

## 4. Mutationen (§3.10, dort, wo eine Zusage prüfbar ist)

Der Slice legt keinen neuen Vertrag an. Geprüft sind deshalb die Zusagen, auf die sich die DoD stützt.

| Zusage | Mutation | Lauf | Urteil |
|---|---|---|---|
| `baseline-verify` erkennt eine geänderte Datei | `m1`: eine Zeile an `regelwerk/modul-03-spec.md` | „weicht von SHA256SUMS ab“, Exit 1 | rot aus dem richtigen Grund |
| erkennt eine ungelistete Datei | `m2`: `templates/extra.md` | „Dateibestand … weicht ab“, `> templates/extra.md`, Exit 1 | rot aus dem richtigen Grund |
| endet bei zwei Tag-Verzeichnissen rot (§6 *Alte Baseline*) | `m3`: `v6.13.0` aus `53ea519` daneben | „mehr als ein <tag>-Verzeichnis“, Exit 1 | rot; die Begründung für Entscheidung (a) trägt |
| belegt nur Stabilität, nicht Herkunft (§6, Entscheidung (1)) | `m4`: Datei und ihre Zeile zusammen entfernt | „v6.16.0 OK — 53 Dateien“, Exit 0 | grün wie zugesagt; die Herkunft hängt am Asset-sha256, das ich selbst nachgeprüft habe |
| `find -xtype l` fängt einen hängenden Symlink | `m5`: `.claude/rules/modul-05-planning-harness.md` zeigt auf `v6.13.0` | `find` meldet die Datei; `make docs-check` 281 Dateien, 0 Befunde | Prüfung trägt; das Negativ „kein Gate sieht ihn“ ist bestätigt |
| `doc-help` listet `docs-check` (Kopf von `d-check.mk`) | `m0` Grundlauf, danach `m6` mit Muster `^doc-` | Grundlauf: Zeile `docs-check` vorhanden; `m6`: 0 Zeilen `docs-check` | rot aus dem richtigen Grund |
| `docs-check` mit `v0.82.0` prüft die aktiven Module (Risiko *Konfigurations-Drift*) | `m7`: kaputter Link, kaputter Anker und eine unverlinkte ADR-Kennung in `harness/README.md` | `target-missing`, `anchor-missing`, `id-unlinked`; d-check Exit 1, `make` Exit 2 | `links`, `anchors` und `ids` wirken. `matrix` und `spans` habe ich nicht selbst geprobt; den Beleg nennt §7 |
| d-check-Befunde aus dem Tausch (§6 *Was am alten Pfad hängt*) | `m8`: `0fca02e` mit `MR-001` von `53ea519` | genau 1 Befund `target-missing`, `MR-001:21` | wie entschieden |

---

## 5. Befunde

| ID | Klasse | Befund | Pfad | Beleg |
|---|---|---|---|---|
| V-76 | DoD-Lücke (Ausgang eines Abgleich-Befunds falsch begründet bzw. fehlend) | §7 begründet bei Befund (2), warum die Gate-ADRs unverändert bleiben, mit „angenommene Gate-ADRs, deren `Schärft:` eine Anforderung nennt“. Das stimmt nicht. Bei [ADR-0025](../plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md), [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md), [ADR-0029](../plan/adr/0029-benannte-welle-kennungen-im-commit-hook.md) und [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) steht `Schärft: —`. `LH-QA-04` steht in ihrem `Bezug:`. [ADR-0028](../plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) und [ADR-0032](../plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) haben mit `SPEC-048` und `SPEC-047` eine Spec-Stelle, die sie nach dem neuen Satz im ADR-Index schärfen müssten, aber nicht nennen. Dass sie so bleiben, hat der Nutzer entschieden, und das ist zulässig. §7 nennt die Entscheidung aber nicht und gibt stattdessen einen falschen Grund. Hinzu kommt ein Teilbefund ohne Ausgang: Für die älteren Werkzeuge (Commit-Hook nach [ADR-0025](../plan/adr/0025-benannte-slice-kennungen-im-commit-hook.md) und [ADR-0029](../plan/adr/0029-benannte-welle-kennungen-im-commit-hook.md), `baseline-verify`, `a-check-negativ`, `hook-gegenprobe`) hat §11 keine Festlegung. Ihre Regeln stehen in ADR, Skriptkopf oder `harness/README.md`, also an Orten, die `v6.16.0` ausschließt. Laut Auftrag hat der Nutzer entschieden, dass sie erst mit einer Randform nach §11 kommen. §7 nennt weder den Teilbefund noch diese Entscheidung. Die DoD verlangt für jeden Befund genau einen Ausgang aus *behoben · Folge-Slice · Register*. Für „aufgeschoben bis zu einer Randform“ muss der Planner einen dieser drei Ausgänge wählen, wahrscheinlich *Register*, oder den Teilbefund mit Grund als „kein Befund“ führen, wie bei (4) und (5) | `docs/plan/planning/in-progress/slice-harness-upgrade-v6-16.md` · „angenommene Gate-ADRs, deren `Schärft:` eine Anforderung nennt“ | `grep -n "Schärft" docs/plan/adr/00{25,28,29,32}*.md`; Spezifikation §11 führt nur `SPEC-047` bis `SPEC-049` |
| V-77 | Entscheidung nicht festgehalten | Der Nutzer hat entschieden, dass `slice-harness-gate-index-werkzeug-teil` trotz F-438 ein Slice bleibt, mit einem Start, der an einem externen Lauf hängt. Diese Entscheidung steht weder im Plan noch im Folge-Slice. Gleiches gilt für die Lesart, auf der `MR-000` bis dahin ruht („die eine Datei, die die Regel für diesen Fall zulässt“). Ohne den Eintrag stellt die nächste Rolle die Frage aus F-438 erneut | `docs/plan/planning/open/slice-harness-gate-index-werkzeug-teil.md` · „ist das die eine Datei, die die Regel für diesen Fall zulässt“ | `grep "Entscheidung des Nutzers"` im Folge-Slice leer |
| V-78 | Urteil zu F-440 (kein Befund) | Der eigene Plan nennt `.harness/baseline/v6.13.0/` an sieben Stellen. Zeile 108 ist die DoD selbst, die das Muster definiert. Zeile 146 nennt das entfernte Verzeichnis als Gegenstand in §3. Zeile 241 ist `git show <vorher>:.harness/baseline/v6.13.0/…`, eine Adresse auf eine Revision, die über `git` auflöst. Zeilen 251 und 255 beschreiben den Zustand vor dem Tausch, Zeile 363 das Risiko und Zeile 501 das, was der `grep` fängt. Keine dieser Stellen schickt einen Leser in den Arbeitsbaum. Alle nennen den Pfad als Gegenstand des Entfernens und fallen damit unter „den alten Stand als Wort nennen, dieser eingeschlossen“. Dieser Teil der DoD ist erfüllt. Grenze: Der `grep` trennt Wort und Verweis nicht, die Prüfung bleibt Urteil. Der Folge-Slice (Zeile 15) und das Drift-Log (Zeile 108) nennen nur das Wort | `docs/plan/planning/in-progress/slice-harness-upgrade-v6-16.md` · „duldet ein Tag-Verzeichnis zur Zeit“ | `grep -rn "baseline/v6\.13\.0" docs/plan/planning/{open,next,in-progress}` |
| V-79 | Beleg fehlt im Eingang (Hinweis) | Die DoD verlangt Digest-Herkunft, Asset-sha256 und die Liste der Blöcke „im Bericht“. Einen Bericht des Implementers gab es nicht, die Werte stehen in §7 *Belege*. Alle habe ich unabhängig bestätigt (Abschnitt 1). Für die Closure reicht §7, wenn die Closure-Notiz sagt, dass er der Bericht ist | DoD Punkt 1 bis 3 | — |

---

## 6. Gate-Ergebnis

`make gates` am Stand `481ef52` im Arbeitsbaum: **grün** (Exit 0). Gelaufen sind laut `GATE_CHECKS` `abdeckung-check`, `abdeckung-gegenprobe`, `a-check` (0 Befunde), `a-check-negativ`, `baseline-verify` (`v6.16.0`, 54 Dateien), `build`, `test`, `docs-check` (281 Dateien, 0 Befunde, Image per Digest `sha256:d28e9437…`), `hook-gegenprobe`, `test-integration`, `kopf-check` und `kopf-check-gegenprobe`. Danach war der Arbeitsbaum unverändert.

**Urteil:** Liefer-Punkt 1 und 2 sind erfüllt. Pin, Digest, Asset-Herkunft, `SHA256SUMS`, `baseline-verify`, `grep` und Symlinks habe ich selbst nachgeprüft, die tragenden Zusagen mit Mutationen. Auch Liefer-Punkt 3 ist erfüllt, mit Ausnahme von V-76: Ein Ausgang in §7 ist falsch begründet, ein Teilbefund hat keinen Ausgang, und zwei Entscheidungen des Nutzers sind nicht festgehalten (V-76, V-77). Bevor der Slice nach `done/` geht, zieht der Planner §7 nach. Code und Baum brauchen keine Änderung. Die Nacharbeit zu F-434 bis F-436 und F-439 trägt.
