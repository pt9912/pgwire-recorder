# Slice slice-harness-ci-range-konfiguration: Die Range-Prüfungen wirken in diesem Repo — `.d-check.yml` führt `commits` und `vcs`, eine Gegenprobe zeigt es

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-erster-release, zweite Stelle der Reihenfolge, direkt hinter `slice-harness-ci-workflows` und vor `slice-harness-ci` (Schnitt aus `slice-harness-ci` am 2026-10-18: Empfehlung des Planners auf Anweisung der Koordination, vom Nutzer nicht bestätigt; Grund in §4 und im Drift-Log der Roadmap). Start-Trigger weicht von der Welle ab: `slice-harness-ci-workflows` und `slice-harness-gegenproben-aus-gates` in `done/` (§4).

**Bezug:** Kein Anspruch des Produkts (`LH-*`): der Slice ändert die Konfiguration eines Gates des Harness, nicht das Binary. Aktive ADRs: [ADR-0050](../../adr/0050-ci-auf-github-actions.md) (CI auf GitHub Actions; ihr Befund, dass `make doc-commits` und `make doc-immutable` ohne die Blöcke nichts prüfen, ist der Gegenstand), [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (Gegenproben in den Listen, nicht in der Kette), [ADR-0026](../../adr/0026-build-und-test-im-multistage-dockerfile.md) (Docker-only).

**Berührte Spec-Stellen:** `SPEC-064` Teil 3 (Punkte 10 bis 12, die Konfiguration, und Punkt 16, der Aufbau von `range-konfiguration-gegenprobe`; die Teilung der Gegenprobe in ein eigenes Ziel hat der Architect am 2026-10-18 in `SPEC-064` vorgenommen; Punkt 17, Zeilenende und Kodierung, gelesen) · `SPEC-060` Teil 1 (gelesen: Namens- und Zuordnungsregel, Liste `GEGENPROBEN_CONTAINER`) · [`SPEC-050`](../../../../harness/vertraege/SPEC-050-struktur-kennungen-im-commit.md) (gelesen: Kennungen im Commit, Strenge des Hooks) · `SPEC-063` (gelesen: Punkt 8, der Eintrag der Pins gehört `slice-harness-ci`).

**Verantwortlich:** pt9912 (Zweig impl-ci-range-konfiguration)
**Autor:** Planner. **Datum:** 2026-10-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make doc-commits` und `make doc-immutable` prüfen in diesem Repo etwas. Heute melden beide über jeder Range „0 Befund(e)“ mit Exit 0, weil `.d-check.yml` die Blöcke `commits` und `vcs` nicht führt; nur der Vorlauf `history-range-guard` fängt eine leere oder unauflösbare Range. Der Slice setzt die Blöcke nach `SPEC-064` Punkte 11 und 12 (Anwesenheit einer Kennung in jeder Nicht-Merge-Message einer Range; eine `Accepted`-ADR ändert sich in der Range nur im Abschnitt `Geschichte`, im Link-Ziel und im Übergang zu `Superseded by`), und eine Gegenprobe (Ziel: Name des Architect, in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS`) zeigt mit dem echten Image von d-check (`DCHECK_REF`) je Zusage und je Element ihrer Aufzählungen (Muster, Ausnahme, Abschnitt, Status-Zeile, Link-Ziel) einen Fall und eine Mutation von `.d-check.yml`, die nur ihn färbt. Der Wert steht einzeln: Auch ohne CI prüft ein Mensch oder ein Skript jede Range lokal mit `make doc-commits RANGE=…` und `make doc-immutable RANGE=…`, und die Konfiguration kann nicht still verschwinden.

**Herkunft (Schnitt vom 2026-10-18):** Die Blöcke waren eine Sendung von `slice-harness-ci-workflows` an `slice-harness-ci` (Liefer-Punkt 2 dort). `slice-harness-ci` lag mit der Schätzung des Architect (535 bis 735) nach der Messung von `slice-harness-ci-workflows` (899 gegen 545 bis 695, Faktor 1,3 bis 1,65; die Gegenprobe trug zwei Drittel) bei rund 700 bis 1210, Mitte 920, über der Marke 900. Die Konfiguration mit ihrer Gegenprobe ist die Naht: Sie hängt weder an `ci.yml` noch an `ci-range.sh` (`SPEC-064` Punkt 16 baut ihre Fälle in eigenen Temp-Repositories), und sie liefert allein Wert (siehe Ziel). `slice-harness-ci` bleibt der Nehmer der Pins (`workflows` in `.d-check.yml`) und aller Teile, die `ci.yml` und die Range betreffen (§4 dort).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `.github/workflows/ci.yml`, `tools/harness/ci-range.sh` und deren Gegenprobe — `slice-harness-ci` (Nehmer, nach diesem Slice); dieser Slice legt keine Workflow-Datei an und braucht keine.
- Der Eintrag des Moduls `workflows` in `.d-check.yml` (`modules:` und `workflows.dir`, `SPEC-063` Punkt 8) — `slice-harness-ci`: Das Modul meldet ohne Datei und ohne `uses:` einen Befund (leere Prüfmenge), der Eintrag kann nicht vor der ersten Workflow-Datei stehen.
- Der Bestand der Historie — die Prüfungen laufen über eine Range, nicht über den Bestand; ältere Commits ohne Kennung und ältere Änderungen an einer `Accepted`-ADR bleiben bewusst stehen (Bestand wird nicht geändert, `SPEC-064` Punkt 15; sonst meldete ein Sensor später gegen Altbestand, den niemand entschieden hat).
- Die Existenz der genannten Slice- oder Welle-Datei und die Ablehnung der Struktur-Kennungen (`SPEC-`, `ARC-`) in der Range — der lokale Hook prüft sie am Index (`SPEC-050`, [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md)); akzeptiertes Negativ nach `SPEC-064` Punkt 11, weil ein Commit der Range einen Slice nennen kann, der in einem späteren Commit derselben Range umbenannt wurde.
- Ein Lauf auf GitHub und der Branch-Schutz — `slice-harness-ci`; hier läuft nichts auf GitHub, und die Einstellung ist aus dem Klon nicht auditierbar.
- Änderungen an `d-check.mk`, `harness/mk/doc-gate.mk` oder `history-range-guard.sh` — das sind die Treiber der Prüfungen, nicht ihre Konfiguration; die Gegenprobe kopiert sie unverändert (Schicht-Abgrenzung: Konfiguration und Gegenprobe, nicht das Werkzeug).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Liefer-Punkt 1 (Konfiguration):** `.d-check.yml` trägt `commits.id-patterns` mit den vier Mustern und `commits.exempt-pattern` (`SPEC-064` Punkt 11) und `vcs.paths`, `vcs.immutable-when`, `vcs.exclude-sections`, `vcs.status-line`, `vcs.head-allow` und `vcs.ignore-link-targets` (`SPEC-064` Punkt 12), wörtlich; `make docs-check` bleibt ohne Befund. Belegt in §7 mit je einem Lauf von `make doc-commits` und `make doc-immutable` über die Range dieses Slice (Anzahl der Befunde, erwartet 0) und mit der Benennung der Commits einer Probe-Range, die die Prüfung rot färbt (im Temp-Repository der Gegenprobe, nicht im Repo).
- [ ] **Liefer-Punkt 2 (Gegenprobe):** Das Ziel `range-konfiguration-gegenprobe` (Fragment `harness/mk/range-konfiguration.mk`, Skript `tools/harness/range-konfiguration-gegenprobe.sh`, `SPEC-064` Punkt 16) fährt je Fall der Punkte 11 und 12 von `SPEC-064` (`p11-*`, `p12-*`) einen Temp-Baum mit echtem Image von d-check und je Zusage und Element eine Mutation von `.d-check.yml`, die nur die Fälle ihrer Gruppe färbt (`AGENTS.md` §3.10); die unveränderte Kopie grün; der Arbeitsbaum bleibt unberührt. Das Ziel steht in `GEGENPROBEN_CONTAINER` (Namens- und Zuordnungsregel von `SPEC-060`); `make gates-schnell-gegenprobe` bleibt grün. Aufrufarten von `make` haben hier keinen Fall (`SPEC-064` Punkt 14: kein Hook und kein Stand liest das Ergebnis).
- [ ] **Liefer-Punkt 3 (Gate-Index):** `harness/README.md` §Sensors führt die Zeile der Gegenprobe mit Vertrag und Bindung; die Prosa nennt die Regel an den Schreiber von `.d-check.yml` und zeigt auf `SPEC-064` und `.d-check.yml` als Quelle, nicht auf eine nachgeschriebene Aufzählung (`AGENTS.md` §3.11).
- [ ] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `harness/` und `.d-check.yml` ändern) und `make gates-schnell-gegenprobe` grün; die Gegenprobe der Konfiguration grün auf dem Endstand.
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
| `.d-check.yml` | update | Blöcke `commits` und `vcs` nach `SPEC-064` Punkte 11, 12 (Sendung aus `slice-harness-ci`, ursprünglich `slice-harness-ci-workflows`) |
| `harness/mk/range-konfiguration.mk` und `tools/harness/range-konfiguration-gegenprobe.sh` | neu | Gegenprobe der Konfiguration mit echtem Image von d-check; Eintrag in `GEGENPROBEN_CONTAINER` |
| `harness/vertraege/SPEC-064-…md` | update (Architect, erledigt 2026-10-18) | Teilung: Punkt 16 `range-konfiguration-gegenprobe`, Punkt 17 Zeilenende; nicht Gegenstand dieses Slice |
| `harness/README.md` §Sensors | update | Zeile der Gegenprobe; Prosa nach §3.11 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-ci-workflows` und `slice-harness-gegenproben-aus-gates` liegen in `done/` (jener als Schritt 1 der Reihenfolge, dieser für die Listen der Gegenproben), dazu der Start-Trigger von [welle-erster-release](../welle-erster-release.md), und der Architect hat `SPEC-064` geteilt (Punkt 16: die Gegenprobe der Konfiguration als eigenes Ziel; `slice-harness-ci` behält Punkt 14 mit den Fällen zu `ci.yml` und `ci-range.sh`) und jede Randform aus §6 entschieden (`AGENTS.md` §3.12). Am 2026-10-18 geprüft: beide Slices in `done/`, Teilung vorhanden.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: mehr als drei Liefer-Punkte, oder der Diff ist nicht in einer Review-Sitzung prüfbar — zurück zur Zerlegung. Marke (Metrik nach `AGENTS.md` §3.14): die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness .d-check.yml Makefile`, Tests und Gegenprobe eingeschlossen, liegt über 500 (Vertragsdateien des Architect unter `harness/vertraege` zählen mit, wie bei `slice-harness-ci-workflows`); Basis: Elternteil des Commits, der diese Datei nach `in-progress/` bewegt, hinter dem Commit des Architect `18efb7f6`, der `SPEC-064` teilt (Basis: `eac90f75`). Schätzung des Architect in dieser Metrik: 300 bis 405 (§8; sie ersetzt die frühere des Planners, 190 bis 330, die den Faktor der Gegenprobe zu niedrig ansetzte). Frühere Schätzung des Planners (2026-10-18, ohne Basis, mit dem gemessenen Faktor 1,3 bis 1,65 aus `slice-harness-ci-workflows` auf den Anteil der Konfiguration an der Schätzung des Architect für `slice-harness-ci`): `.d-check.yml` 20 bis 30, Fragment 20, Gegenprobe (17 Fälle der Punkte 11 und 12, rund 13 Mutationen, je Fall ein Git-Repository mit echtem Image) 100 bis 140, `harness/README.md` 5 bis 10, zusammen 145 bis 200, mit dem Faktor 190 bis 330. Die Marke liegt bei rund dem Anderthalbfachen der Obergrenze; Schnitt, falls sie eintritt: Fälle zu `commits` (Punkt 11) und Fälle zu `vcs` (Punkt 12).
- `in-progress` → `open`: das gepinnte d-check (`DCHECK_REF`) liest die Muster, die Ausnahme oder die Felder von `vcs` anders als `SPEC-064` Punkte 11 und 12 sie beschreiben (§6, Risiko 1); der Architect entscheidet dann neu.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, `make gates` und die Gegenprobe der Konfiguration grün auf dem Endstand, Review- und Verifikations-Report liegen vor, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Das gepinnte d-check liest die Konfiguration anders, als `SPEC-064` sie beschreibt** (Muster als regulärer Ausdruck, `exempt-pattern` auf dem Betreff, `exclude-sections` in der Bedeutung für `vcs`, `head-allow`). — **Ausgang (vorgesehen):** entfallen, wenn die Gegenprobe mit dem echten Image jeden Fall der Punkte 11 und 12 zeigt; sonst eingetreten → Rückführung nach `open` (§4) und Rückgabe an den Architect.
- **Ein Befund über den Bestand färbt die erste Range rot, ohne dass die Konfiguration falsch ist** (`SPEC-064` Punkt 15: die erste Probe auf GitHub prüft die Range des ersten Push). — **Ausgang (vorgesehen):** entfallen, wenn der Beleg in §7 die Range dieses Slice mit 0 Befunden nennt; der Bestand bleibt bewusst stehen (§1).
- **Die Gegenprobe braucht das Image von d-check** (Netz beim ersten Pull). — **Ausgang (vorgesehen):** entfallen, weil sie in `GEGENPROBEN_CONTAINER` steht und kein Teil von `make gates` ist.
- **Ein Löschen der Blöcke färbt die Gegenprobe, nicht die CI** (`SPEC-064` Punkt 15, akzeptiertes Negativ). — **Ausgang (vorgesehen):** entfallen mit dem Grund aus `SPEC-064` Punkt 15, solange `make gegenproben` am Ende jeder Welle fährt; sonst weiter offen → Register.
- **Die Größe liegt im Schätzbereich 300 bis 405 (Architect) gegen die Marke 500.** — **Ausgang (vorgesehen):** entfallen, wenn die Messung nach der Metrik unter 500 bleibt.

**Randformen** (`AGENTS.md` §3.12) — entschieden vom Architect in `SPEC-064` Teil 3 (Vertrag vom 2026-10-18); die Teilung ist erfolgt (Punkt 16 und 17). Vorbereitet vom Planner (Empfehlungen, kein Beschluss):

1. **Block fehlt.** Ein fehlender Block lässt das Modul inert (0 Befunde). **Entschieden:** `SPEC-064` Punkt 10 und die Mutation „Block fehlt“ in Punkten 11 und 12.
2. **Kennung je Muster.** Ein Fall je Muster (ADR, LH, MR, Slice, Welle); Merge- und Revert-Betreff; Kennung im Rumpf, im Kommentar. **Entschieden:** `SPEC-064` Punkt 11.
3. **ADR-Änderung je Form.** Kern geändert, nur `Geschichte`, `Superseded by`, `Rejected`, `Proposed`, gelöscht, Link-Ziel. **Entschieden:** `SPEC-064` Punkt 12.
4. **Plattform.** Docker unter Linux (`linux/amd64`); GNU-Werkzeuge in der Gegenprobe. **Entschieden:** `SPEC-064` Abschnitt Plattform (andere Plattformen nicht zugesagt).
5. **Aufrufarten von `make`.** Kein Hook und kein Stand liest das Ergebnis; kein Fall. **Entschieden:** `SPEC-064` Punkt 14 (Aufrufarten; für `range-konfiguration-gegenprobe` kein Fall, Punkt 16).
6. **Zeilenende und Kodierung der Commit-Messages und der ADR-Dateien** (CRLF in einer Datei der Range). Empfehlung: der Architect nennt die Grenze oder einen Fall, denn das Modul `vcs` vergleicht Zeilen. **Entschieden:** `SPEC-064` Punkt 17 (Grenze: LF und UTF-8 zugesagt, CRLF und andere Kodierung nicht; der Implementer misst eine Message mit `\r\n` und eine ADR mit umgestelltem Zeilenende einmal und trägt das Ergebnis in §7 ein, ohne Fall und ohne Zusage); `BEO-REPO/zeilenende-crlf-im-gelesenen-wert` führt die Beobachtung aus `slice-harness-ci-workflows`.
7. **Name, Fragment und Skript der Gegenprobe; Abgrenzung zu `ci-gegenprobe`.** **Entschieden:** `SPEC-064` Kopf („Zwei Slices, zwei Ziele“) und Punkt 16: Ziel `range-konfiguration-gegenprobe`, Fragment `harness/mk/range-konfiguration.mk`, Skript `tools/harness/range-konfiguration-gegenprobe.sh`, in `GEGENPROBEN_CONTAINER`; die Fälle `p11-*` und `p12-*` gehören hierher, keine `p1-*` bis `p9-*` und `p14-*`. Der Eintrag `workflows` in `.d-check.yml` gehört `slice-harness-ci`.

**Fragen an den Nutzer** (jede mit Empfehlung des Planners; eine Empfehlung ist keine Entscheidung):

1. **Schnitt.** Die Teilung aus `slice-harness-ci` ist die Empfehlung des Planners auf Anweisung der Koordination. Alternative: ungeteilt lassen und bei der Marke von `slice-harness-ci` zurückführen (§4 dort). Empfehlung: geteilt, weil die Schätzung des ungeteilten Slice mit dem gemessenen Faktor in der Mitte über der Marke liegt.

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

**Vorgelagert — offene Beobachtungen sichten:** `docs/plan/planning/observations/BEO-REPO/` durchgesehen (Stand 2026-10-18). Treffer: `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` (`AGENTS.md` §3.10) und `BEO-REPO/spec-randform-erst-im-review-entschieden` (`AGENTS.md` §3.12) — die Gegenprobe ist ein neuer Vertrag, je Zusage und Element eine Mutation; `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (2×, offen) — genau der Gegenstand: eine Konfiguration liest sich als Zusage, das Werkzeug wertet sie anders aus; die Gegenprobe mit echtem Image ist das Gegenmittel, ein drittes Auftreten hier wäre Schwelle; `BEO-REPO/werkzeug-commit-ohne-zugelassene-kennung` und `BEO-REPO/commit-nennt-struktur-kennung` (Kennungen im Commit); `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (6×, offen, ohne Ausgang) — die Gegenprobe läuft im Container als Benutzer ohne root, mit Locale und Zeitzone des Containers und am Host; der Implementer fährt sie unter `LANG=de_DE.UTF-8`, als root und als Benutzer. Kein Eintrag erreicht mit diesem Plan die Schwelle 3× neu.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13 (Planner, 2026-10-18, im Commit, der den Plan anlegt):** Der Slice trägt eine Sendung: die Blöcke `commits` und `vcs` in `.d-check.yml` (Geber: `slice-harness-ci-workflows` an `slice-harness-ci`, von dort mit diesem Schnitt an diesen Slice; `grep -n slice-harness-ci-workflows` und `grep -n slice-harness-ci ` finden sie in §1, §3 und DoD Liefer-Punkt 1). Liefer-Punkte 3 (Konfiguration | Gegenprobe | Gate-Index). Schichten 1 (Harness: `.d-check.yml`, `harness/mk`, `tools/harness`, `harness/README.md`); Planung, Vertrag und ADR sind keine Schicht, Wartungs-Doku und Handbuch bleiben unberührt. Schätzung 300 bis 405 (Architect, unten), unter der Marke 500.

**Schätzung des Architect nach dem Schnitt (2026-10-18; Metrik von §4; `SPEC-064` vor dem Anspruch geschrieben, nicht gezählt; Faktor der Gegenprobe aus `slice-harness-ci-workflows`, gemessen 899 gegen geschätzt 545 bis 695):** `.d-check.yml` 25 bis 35, Fragment 20, Gegenprobe 250 bis 340 (Aufbau der Git-Repositories und Aufruf von d-check 70, 17 Fälle, rund 13 Mutationen, Kopf), `harness/README.md` 5 bis 10, zusammen 300 bis 405, unter der Marke 500.
