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

**Berührte Spec-Stellen:** `SPEC-064` Teil 3 (Punkte 10 bis 12, die Konfiguration, und Punkt 16, der Aufbau von `range-konfiguration-gegenprobe`; die Teilung der Gegenprobe in ein eigenes Ziel hat der Architect am 2026-10-18 in `SPEC-064` vorgenommen; Punkt 18, Abbruch und Herunterfahren der Gegenprobe durch Ende, `SIGINT` und `SIGTERM`, vom Architect nach dem Review ergänzt und von diesem Slice geliefert; Punkt 17, Zeilenende und Kodierung, gelesen) · `SPEC-060` Teil 1 (gelesen: Namens- und Zuordnungsregel, Liste `GEGENPROBEN_CONTAINER`) · [`SPEC-050`](../../../../harness/vertraege/SPEC-050-struktur-kennungen-im-commit.md) (gelesen: Kennungen im Commit, Strenge des Hooks) · `SPEC-063` (gelesen: Punkt 8, der Eintrag der Pins gehört `slice-harness-ci`).

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

**Ziel:** `make doc-commits` und `make doc-immutable` prüfen in diesem Repo etwas. Heute melden beide über jeder Range „0 Befund(e)“ mit Exit 0, weil `.d-check.yml` die Blöcke `commits` und `vcs` nicht führt; nur der Vorlauf `history-range-guard` fängt eine leere oder unauflösbare Range. Der Slice setzt die Blöcke nach `SPEC-064` Punkte 11 und 12 (Anwesenheit einer Kennung in jeder Nicht-Merge-Message einer Range; eine `Accepted`-ADR ändert sich in der Range nur im Abschnitt `Geschichte`, im Link-Ziel und im Übergang zu `Superseded by`), und eine Gegenprobe (Ziel: Name des Architect, in `GEGENPROBEN_CONTAINER`, nicht in `GATE_CHECKS`) zeigt mit dem echten Image von d-check (`DCHECK_REF`) je Zusage und je Element ihrer Aufzählungen (Muster mit ihrer Form und ihrem Anker, Ausnahme, Abschnitt, Status-Zeile, Link-Ziel) einen Fall und eine Mutation von `.d-check.yml`, die nur ihn färbt. Die Gegenprobe räumt bei Ende, `SIGINT` und `SIGTERM` ihren Temp-Baum und beendet ihre Nachkommen (`SPEC-064` Punkt 18; Exit 130 und 143, Fälle `p18-*` mit Mutationen des Skripts), zählt einen Namen der Fallliste ohne Funktion als roten Fall und startet in ihrer Kopie keine zweite Kopie, die Mutationen des Abbruchs fährt. Die Zeichenklassen der Muster und die Leerzeichen der Ausnahme sind nicht zugesagt (Grenze in `SPEC-064` Punkt 11). Der Wert steht einzeln: Auch ohne CI prüft ein Mensch oder ein Skript jede Range lokal mit `make doc-commits RANGE=…` und `make doc-immutable RANGE=…`, und die Konfiguration kann nicht still verschwinden.

**Herkunft (Schnitt vom 2026-10-18):** Die Blöcke waren eine Sendung von `slice-harness-ci-workflows` an `slice-harness-ci` (Liefer-Punkt 2 dort). `slice-harness-ci` lag mit der Schätzung des Architect (535 bis 735) nach der Messung von `slice-harness-ci-workflows` (899 gegen 545 bis 695, Faktor 1,3 bis 1,65; die Gegenprobe trug zwei Drittel) bei rund 700 bis 1210, Mitte 920, über der Marke 900. Die Konfiguration mit ihrer Gegenprobe ist die Naht: Sie hängt weder an `ci.yml` noch an `ci-range.sh` (`SPEC-064` Punkt 16 baut ihre Fälle in eigenen Temp-Repositories), und sie liefert allein Wert (siehe Ziel). `slice-harness-ci` bleibt der Nehmer der Pins (`workflows` in `.d-check.yml`) und aller Teile, die `ci.yml` und die Range betreffen (§4 dort).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `.github/workflows/ci.yml`, `tools/harness/ci-range.sh` und deren Gegenprobe — `slice-harness-ci` (Nehmer, nach diesem Slice); dieser Slice legt keine Workflow-Datei an und braucht keine.
- Der Eintrag des Moduls `workflows` in `.d-check.yml` (`modules:` und `workflows.dir`, `SPEC-063` Punkt 8) — `slice-harness-ci`: Das Modul meldet ohne Datei und ohne `uses:` einen Befund (leere Prüfmenge), der Eintrag kann nicht vor der ersten Workflow-Datei stehen.
- Der Bestand der Historie — die Prüfungen laufen über eine Range, nicht über den Bestand; ältere Commits ohne Kennung und ältere Änderungen an einer `Accepted`-ADR bleiben bewusst stehen (Bestand wird nicht geändert, `SPEC-064` Punkt 15; sonst meldete ein Sensor später gegen Altbestand, den niemand entschieden hat).
- Die Existenz der genannten Slice- oder Welle-Datei und die Ablehnung der Struktur-Kennungen (`SPEC-`, `ARC-`) in der Range — der lokale Hook prüft sie am Index (`SPEC-050`, [ADR-0025](../../adr/0025-benannte-slice-kennungen-im-commit-hook.md)); akzeptiertes Negativ nach `SPEC-064` Punkt 11, weil ein Commit der Range einen Slice nennen kann, der in einem späteren Commit derselben Range umbenannt wurde.
- Ein Lauf auf GitHub und der Branch-Schutz — `slice-harness-ci`; hier läuft nichts auf GitHub, und die Einstellung ist aus dem Klon nicht auditierbar.
- Zeichenklassen der Muster (`[A-Z]`, `[a-z0-9]`) und Leerzeichen der Ausnahme `^(Merge |Revert )`, Umbenennung der verlinkten Datei, die Zeile `Merge …` nach einer Leerzeile, `SIGKILL` und `SIGHUP` — Grenzen in `SPEC-064` Punkte 11, 12 und 18, vom Architect gezogen (nach Review F-6, F-7 und Verifikation V-2, V-3); kein Fall und keine Mutation, ein Fall bräuchte eine Zusage, die der Vertrag nicht gibt (Klasse 2: bewusst stehen gelassen, mit Grund im Vertrag).
- Änderungen an `d-check.mk`, `harness/mk/doc-gate.mk` oder `history-range-guard.sh` — das sind die Treiber der Prüfungen, nicht ihre Konfiguration; die Gegenprobe kopiert sie unverändert (Schicht-Abgrenzung: Konfiguration und Gegenprobe, nicht das Werkzeug).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Liefer-Punkt 1 (Konfiguration):** `.d-check.yml` trägt `commits.id-patterns` mit den vier Mustern und `commits.exempt-pattern` (`SPEC-064` Punkt 11) und `vcs.paths`, `vcs.immutable-when`, `vcs.exclude-sections`, `vcs.status-line`, `vcs.head-allow` und `vcs.ignore-link-targets` (`SPEC-064` Punkt 12), wörtlich; `make docs-check` bleibt ohne Befund. Belegt in §7 mit je einem Lauf von `make doc-commits` und `make doc-immutable` über die Range dieses Slice (Anzahl der Befunde, erwartet 0) und mit der Benennung der Commits einer Probe-Range, die die Prüfung rot färbt (im Temp-Repository der Gegenprobe, nicht im Repo).
- [x] **Liefer-Punkt 2 (Gegenprobe):** Das Ziel `range-konfiguration-gegenprobe` (Fragment `harness/mk/range-konfiguration.mk`, Skript `tools/harness/range-konfiguration-gegenprobe.sh`, `SPEC-064` Punkte 16 und 18) fährt je Fall der Punkte 11, 12 und 18 von `SPEC-064` (`p11-*` einschließlich der `-knapp`-Fälle, `p12-*`, `p18-sigint`, `p18-sigterm`; zusammen 27) einen Temp-Baum mit echtem Image von d-check und je Zusage und Element eine Mutation von `.d-check.yml` (Form und Anker je Muster eingeschlossen), die nur die Fälle ihrer Gruppe färbt, und für den Abbruch je Signal drei Mutationen des Skripts (zusammen 41, `AGENTS.md` §3.10); dazu die Läufe `fall-ohne-funktion` (ein Name der Fallliste ohne Funktion ist rot) und `rekursion-schutz` (die Kopie des Skripts startet keine zweite Kopie bei `GP_MUT` auf `m-p18-*`); die unveränderte Kopie grün; der Arbeitsbaum bleibt unberührt; bei Ende, `SIGINT` und `SIGTERM` bleibt kein Temp-Baum und kein Nachkomme. Das Ziel steht in `GEGENPROBEN_CONTAINER` (Namens- und Zuordnungsregel von `SPEC-060`); `make gates-schnell-gegenprobe` bleibt grün. Aufrufarten von `make` haben hier keinen Fall (`SPEC-064` Punkt 14: kein Hook und kein Stand liest das Ergebnis).
- [x] **Liefer-Punkt 3 (Gate-Index):** `harness/README.md` §Sensors führt die Zeile der Gegenprobe mit Vertrag und Bindung; die Prosa nennt die Regel an den Schreiber von `.d-check.yml` und zeigt auf `SPEC-064` und `.d-check.yml` als Quelle, nicht auf eine nachgeschriebene Aufzählung (`AGENTS.md` §3.11).
- [x] `make gates` grün (Stop-Hook verlangt den vollen Lauf, weil `harness/` und `.d-check.yml` ändern) und `make gates-schnell-gegenprobe` grün; die Gegenprobe der Konfiguration grün auf dem Endstand. Belegt in §7: `make gates` Exit 0 (4 min 21 s) am Stand `c1d2b1f3`, dem letzten Code-Commit (danach nur §7 und Plan); der Host trägt den Lauf am Integrationsstand nach der Closure nach.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen“ in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.d-check.yml` | update | Blöcke `commits` und `vcs` nach `SPEC-064` Punkte 11, 12 (Sendung aus `slice-harness-ci`, ursprünglich `slice-harness-ci-workflows`) |
| `harness/mk/range-konfiguration.mk` und `tools/harness/range-konfiguration-gegenprobe.sh` | neu | Gegenprobe der Konfiguration mit echtem Image von d-check; 27 Fälle (17 `p11-*`, 8 `p12-*`, 2 `p18-*`), 41 Mutationen (von `.d-check.yml` und, für den Abbruch, des Skripts) und die Läufe `ok`, `fall-ohne-funktion` und `rekursion-schutz` (zusammen 44); Eintrag in `GEGENPROBEN_CONTAINER` |
| `harness/vertraege/SPEC-064-…md` | update (Architect, erledigt 2026-10-18) | Teilung: Punkt 16 `range-konfiguration-gegenprobe`, Punkt 17 Zeilenende; nicht Gegenstand dieses Slice |
| `harness/README.md` §Sensors | update | Zeile der Gegenprobe; Prosa nach §3.11 |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-harness-ci-workflows` und `slice-harness-gegenproben-aus-gates` liegen in `done/` (jener als Schritt 1 der Reihenfolge, dieser für die Listen der Gegenproben), dazu der Start-Trigger von [welle-erster-release](welle-erster-release.md), und der Architect hat `SPEC-064` geteilt (Punkt 16: die Gegenprobe der Konfiguration als eigenes Ziel; `slice-harness-ci` behält Punkt 14 mit den Fällen zu `ci.yml` und `ci-range.sh`) und jede Randform aus §6 entschieden (`AGENTS.md` §3.12). Am 2026-10-18 geprüft: beide Slices in `done/`, Teilung vorhanden.

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

- **Das gepinnte d-check liest die Konfiguration anders, als `SPEC-064` sie beschreibt** (Muster als regulärer Ausdruck, `exempt-pattern` auf dem Betreff, `exclude-sections` in der Bedeutung für `vcs`, `head-allow`). — **Ausgang:** entfallen: Die Gegenprobe mit dem echten Image zeigt jeden Fall der Punkte 11 und 12 (27 Fälle, 41 Mutationen, `grün (44 Läufe)`; die Verifikation fuhr die Mutationen am Host nach), die Vorprüfung des Implementers fand alle acht Schlüssel im gepinnten Image; Rückführung nach `open` trat nicht ein. Die Fälle, in denen d-check anders auswertete, als es sich las, sind als Beobachtung eingetragen (`BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, 3×), das Risiko selbst ist nicht eingetreten.
- **Ein Befund über den Bestand färbt die erste Range rot, ohne dass die Konfiguration falsch ist** (`SPEC-064` Punkt 15: die erste Probe auf GitHub prüft die Range des ersten Push). — **Ausgang:** entfallen: `make doc-commits` und `make doc-immutable` melden über `368e3d43..f6e346ed` (Implementer) und über `368e3d43..d8086425` (Verifikation, frischer Klon) 0 Befunde, der Implementer auch über den Bestand (2076 Commits); von Hand gegengelesen für `368e3d43..7a460361`: kein Commit ohne die Kennung dieses Slice, keine Änderung an `docs/plan/adr/`. Der Bestand bleibt bewusst stehen (§1). Die Messung mit `make` über die Range am Integrationsstand trägt der Host nach der Closure nach (siehe DoD).
- **Die Gegenprobe braucht das Image von d-check** (Netz beim ersten Pull). — **Ausgang:** entfallen: Das Ziel steht in `GEGENPROBEN_CONTAINER` und ist kein Teil von `make gates`; `make gates-schnell-gegenprobe` bestätigt die Zuordnung (Implementer und Verifikation), `make gates` lief ohne die Gegenprobe grün.
- **Ein Löschen der Blöcke färbt die Gegenprobe, nicht die CI** (`SPEC-064` Punkt 15, akzeptiertes Negativ). — **Ausgang:** entfallen mit dem Grund aus `SPEC-064` Punkt 15: Die Bedingung steht in `AGENTS.md` §6 Schritt 6 („am Ende jeder Welle und vor dem Release fahren `make gegenproben` im Gegenproben-Container und `make gegenproben-host` am Host“), und `range-konfiguration-gegenprobe` steht in `GEGENPROBEN_CONTAINER`; ein Löschen der Blöcke färbt sie spätestens dort. Fiele die Fahrt am Ende der Welle weg, wäre das ein Bruch von `AGENTS.md` §6, kein offenes Risiko dieses Slice.
- **Die Größe liegt im Schätzbereich 300 bis 405 (Architect) gegen die Marke 500.** — **Ausgang:** entfallen: Gemessen 445 am letzten Commit des Slice `7a460361` (Summe aus Einfügungen und Löschungen, Basis `eac90f75`; Skript 398, `.d-check.yml` 19, Fragment 11, `SPEC-064` 16, `harness/README.md` 1; 433 am Stand `8d2b4ec3`, den die Verifikation am Kopf `d8086425` ebenso maß), unter der Marke 500 und der Meldeschwelle 450; die Schätzung des Architect (300 bis 405, nach dem Review 365 bis 415) ist um 40 bzw. 30 überschritten (`BEO-REPO/schaetzung-in-zeilen-unterschaetzt-tests`, 17×).

**Randformen** (`AGENTS.md` §3.12) — entschieden vom Architect in `SPEC-064` Teil 3 (Vertrag vom 2026-10-18); die Teilung ist erfolgt (Punkt 16 und 17). Vorbereitet vom Planner (Empfehlungen, kein Beschluss):

1. **Block fehlt.** Ein fehlender Block lässt das Modul inert (0 Befunde). **Entschieden:** `SPEC-064` Punkt 10 und die Mutation „Block fehlt“ in Punkten 11 und 12.
2. **Kennung je Muster.** Ein Fall je Muster (ADR, LH, MR, Slice, Welle); Merge- und Revert-Betreff; Kennung im Rumpf, im Kommentar. **Entschieden:** `SPEC-064` Punkt 11.
3. **ADR-Änderung je Form.** Kern geändert, nur `Geschichte`, `Superseded by`, `Rejected`, `Proposed`, gelöscht, Link-Ziel. **Entschieden:** `SPEC-064` Punkt 12.
4. **Plattform.** Docker unter Linux (`linux/amd64`); GNU-Werkzeuge in der Gegenprobe. **Entschieden:** `SPEC-064` Abschnitt Plattform (andere Plattformen nicht zugesagt).
5. **Aufrufarten von `make`.** Kein Hook und kein Stand liest das Ergebnis; kein Fall. **Entschieden:** `SPEC-064` Punkt 14 (Aufrufarten; für `range-konfiguration-gegenprobe` kein Fall, Punkt 16).
6. **Zeilenende und Kodierung der Commit-Messages und der ADR-Dateien** (CRLF in einer Datei der Range). Empfehlung: der Architect nennt die Grenze oder einen Fall, denn das Modul `vcs` vergleicht Zeilen. **Entschieden:** `SPEC-064` Punkt 17 (Grenze: LF und UTF-8 zugesagt, CRLF und andere Kodierung nicht; der Implementer misst eine Message mit `\r\n` und eine ADR mit umgestelltem Zeilenende einmal und trägt das Ergebnis in §7 ein, ohne Fall und ohne Zusage); `BEO-REPO/zeilenende-crlf-im-gelesenen-wert` führt die Beobachtung aus `slice-harness-ci-workflows`.
7. **Name, Fragment und Skript der Gegenprobe; Abgrenzung zu `ci-gegenprobe`.** **Entschieden:** `SPEC-064` Kopf („Zwei Slices, zwei Ziele“) und Punkt 16: Ziel `range-konfiguration-gegenprobe`, Fragment `harness/mk/range-konfiguration.mk`, Skript `tools/harness/range-konfiguration-gegenprobe.sh`, in `GEGENPROBEN_CONTAINER`; die Fälle `p11-*` und `p12-*` gehören hierher, keine `p1-*` bis `p9-*` und `p14-*`. Der Eintrag `workflows` in `.d-check.yml` gehört `slice-harness-ci`.

8. **Abbruch und Herunterfahren der Gegenprobe** (`SIGINT`, `SIGTERM`, Ende; Unterprozesse und Docker-Läufe). **Entschieden:** `SPEC-064` Punkt 18 (räumt ganz, beendet Unterprozesse zuerst, Exit 130 und 143; Fälle `p18-sigint`, `p18-sigterm`).
9. **Form der Muster, Ausnahme am Betreff, Struktur-Kennung allein, Umbenennung im Link-Ziel, Index.** **Entschieden:** `SPEC-064` Punkte 11 und 12 (Form je Muster mit `-knapp`-Fall und Mutationen; Fälle `p11-exempt-nicht-am-anfang`, `p11-struktur-allein`, `p12-index` in der Liste; Umbenennung und „Merge nach Leerzeile“ als Grenze).

10. **Zeichenklassen der Muster und Leerzeichen der Ausnahme** (Verifikation V-2, V-3). **Entschieden:** `SPEC-064` Punkt 11: nicht zugesagt (Grenze mit Grund); kein Fall, keine Mutation. Der Schwanz `(-[a-z0-9]+)*` des Slice-Musters ist für die Anwesenheit äquivalent (Invariante, Ort, Sonde stehen im Report der Verifikation und gehören in §7).

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

### Belege des Implementers

Stand: Code-Commits `6022bc9e`, `f6e346ed` und, nach dem Review, `8d2b4ec3` (Fälle `-knapp`, Mutationen Weitung und Anker, Abbruch nach Punkt 18; die Belege dieses Abschnitts gelten für diesen Stand, wo sie nichts anderes nennen); Zweig `impl-ci-range-konfiguration`, Basis der Metrik `eac90f75`. Die Gegenprobe läuft im Dev-Container pgr-dev3 (Benutzer `vscode`, uid 1000, Docker-in-Docker, d-check v0.86.1, Digest `3e0b9779…`).

**Vorprüfung (Risiko 1, Rückgabe an den Architect).** Das gepinnte Image kennt alle acht Schlüssel (`d-check --print-config`: `commits.id-patterns`, `commits.exempt-pattern`, `vcs.paths`, `vcs.immutable-when`, `vcs.exclude-sections`, `vcs.status-line`, `vcs.head-allow`, `vcs.ignore-link-targets`); Meldungsform und Exit-Codes stimmen mit dem Vertrag (Exit 1, Code `commit-untraceable` bzw. `core-drift-vcs`, Summenzeile auf stderr). Beide Module stehen nicht in `modules:`; die Targets schalten sie mit `--enable` ein. Keine Rückgabe.

**Befundzeilen vorher und nachher** (Temp-Repository: Basis-Commit mit `Accepted`-ADR, dann ein Commit `foo` ohne Kennung und ein Commit, der den Text der ADR ändert; `.d-check.yml` aus `368e3d43` bzw. vom Code-Commit):

| Aufruf | `.d-check.yml` ohne die Blöcke (`368e3d43`) | mit den Blöcken |
|---|---|---|
| `--enable commits --range base..HEAD` | `d-check: 1 Datei(en) geprüft, 0 Befund(e)`, Exit 0 | `<sha>:1	<sha>	commit-untraceable	foo`, 1 Befund(e), Exit 1 |
| `--enable vcs --range base..HEAD` | `0 Befund(e)`, Exit 0 | `docs/plan/adr/0001-eins.md:3 … core-drift-vcs	Core einer immutablen Datei hat sich über die Commit-Range geändert`, 1 Befund(e), Exit 1 |

**Range dieses Slice über das Repo** (nach dem letzten Commit): `make doc-commits RANGE=368e3d43..f6e346ed` und `make doc-immutable RANGE=368e3d43..f6e346ed` (2 Commits, die beiden Code-Commits): je `history-range-guard: … 2 Commit(s) — OK` und `d-check: 1174 Datei(en) geprüft, 0 Befund(e)`, Exit 0. Über den Bestand (`d027ecc5…..f6e346ed`, 2076 Commits) beide Targets ebenfalls 0 Befund(e), Exit 0: Die erste Range auf GitHub färbt der Bestand nicht (Risiko 2, soweit der Klon die ganze Historie hat). Commits einer Probe-Range, die die Prüfung rot färbt, stehen im Temp-Repository der Gegenprobe: `foo` (ohne Kennung, Fall `p11-ohne-kennung`) und der Commit, der in `0001-eins.md` den Text ändert (Fall `p12-kern`); DoD Liefer-Punkt 1 ist damit mit Befundzeilen vorher und nachher belegt (Tabelle oben).

**Fälle (27).** `p11-*` (17): `ohne-kennung`, `adr`, `lh`, `mr`, `slice`, `welle` (je mit der Kennung mitten in einer Zeile, `Siehe <Kennung> x`), `adr-knapp` (`ADR-123`), `lh-knapp` (`LH-FA-1`), `mr-knapp` (`MR-12`), `slice-knapp` (`slice- x`), `welle-knapp` (`welle- x`) (die fünf `-knapp` rot), `revert`, `merge-betreff`, `kennung-im-rumpf`, `kennung-im-kommentar`, `exempt-nicht-am-anfang`, `struktur-allein`; `p12-*` (8): `kern`, `geschichte`, `ersetzt`, `ersetzt-anders`, `vorschlag`, `geloescht`, `link-ziel`, `index`; `p18-*` (2): `sigint`, `sigterm`. Die Fallliste des Vertrags führt `p11-exempt-nicht-am-anfang`, `p11-struktur-allein` (mit Mutation) und `p12-index` (Beleg ohne eigene Mutation) seit der Entscheidung des Architect vom 2026-10-18. Die Fälle `p11-slice` und `p11-welle` nennen eine Slice- und eine Welle-Datei, die im Temp-Repository nicht existiert; sie sind zugleich der Fall zum akzeptierten Negativ „Existenz wird nicht geprüft“ (Punkt 11). Ein Fall besteht, wenn Exit-Code, stdout (ohne die Zeile des Vorlaufs `history-range-guard`) und stderr vollständig dem erwarteten Text entsprechen (SHAs, Dateizahl und Zeilennummer in `d-check.mk` normiert); Aufrufarten von `make` haben keinen Fall (Punkt 16). Ein Name der Fallliste ohne Funktion zählt als roter Fall (Lauf `fall-ohne-funktion` mit dem erfundenen Namen `p11-tippfehler`: rot ist genau er).

**Mutationen (41): Zusage · Mutation · roter Test.** Jede Mutation ist ein Textersatz in der Kopie von `.d-check.yml`; rot sind genau die benannten Fälle (die Gegenprobe verlangt diese Menge, mehr und weniger machen sie rot); mit `GP_ENTDECKE=1` gesehen, dann mit dem vollen Lauf (`43 Läufe`, grün: 41 Mutationen, die unveränderte Kopie mit allen 27 Fällen und der Lauf `fall-ohne-funktion`). Die Gruppe nennt, über welche Fälle die Mutation läuft; die beiden `Block fehlt` laufen über alle 25 Fälle von p11 und p12; die Mutationen des Skripts (`m-p18-*`) laufen über die beiden Fälle p18, ihr Textersatz gilt dem Skript statt `.d-check.yml` (kein p12-Fall wird von `commits` gefärbt und umgekehrt).

| Mutation | Zusage / Element | rote Fälle |
|---|---|---|
| `m-p11-block-fehlt` (alle) | Block `commits` wirkt (Punkt 10/11) | `p11-ohne-kennung`, `p11-kennung-im-kommentar`, `p11-exempt-nicht-am-anfang`, `p11-struktur-allein`, die fünf `-knapp`-Fälle (`p11-exempt-nicht-am-anfang` steht in der Aufzählung von Punkt 11 nicht, ist aber rot) |
| `m-p11-muster-adr-fehlt` | Muster `ADR-[0-9]{4}` | `p11-adr` |
| `m-p11-muster-lh-fehlt` | Muster `LH-[A-Z]{2}-[0-9]{2}` | `p11-lh` |
| `m-p11-muster-mr-fehlt` | Muster `MR-[0-9]{3}` | `p11-mr`, `p11-kennung-im-rumpf` (der Fall trägt `MR-001`) |
| `m-p11-muster-slice-welle-fehlt` | Muster `(slice\|welle)-…` | `p11-slice`, `p11-welle` |
| `m-p11-alternative-slice-fehlt` | Element `slice` der Alternative | `p11-slice` |
| `m-p11-alternative-welle-fehlt` | Element `welle` der Alternative | `p11-welle` |
| `m-p11-muster-adr-geweitet` | Form `ADR-` und genau vier Ziffern (`ADR-[0-9]`) | `p11-adr-knapp` |
| `m-p11-muster-lh-geweitet` | Form `LH-`, zwei Großbuchstaben, zwei Ziffern (`LH-[A-Z]{1,2}-[0-9]`) | `p11-lh-knapp` |
| `m-p11-muster-mr-geweitet` | Form `MR-` und drei Ziffern (`MR-[0-9]`) | `p11-mr-knapp` |
| `m-p11-muster-slice-welle-geweitet` | Form `slice-`/`welle-` mit mindestens einem Zeichen (`(slice\|welle)-`) | `p11-slice-knapp`, `p11-welle-knapp` |
| `m-p11-muster-adr-verankert` | Kennung an beliebiger Stelle der Zeile, Muster ADR (Anker `^`) | `p11-adr` |
| `m-p11-muster-lh-verankert` | dasselbe, Muster LH | `p11-lh` |
| `m-p11-muster-mr-verankert` | dasselbe, Muster MR | `p11-mr`, `p11-kennung-im-rumpf` (der Fall trägt `MR-001` in der dritten Zeile) |
| `m-p11-muster-slice-welle-verankert` | dasselbe, Muster slice/welle | `p11-slice`, `p11-welle` |
| `m-p11-muster-struktur` | `SPEC-…` zählt nicht als Kennung (Muster `SPEC-[0-9]{3}` ergänzt) | `p11-struktur-allein` |
| `m-p11-ausnahme-fehlt` | `exempt-pattern` | `p11-revert`, `p11-merge-betreff` |
| `m-p11-ausnahme-merge-fehlt` | Element `Merge ` | `p11-merge-betreff` |
| `m-p11-ausnahme-revert-fehlt` | Element `Revert ` | `p11-revert` |
| `m-p11-ausnahme-ohne-anker` | Anker `^` der Ausnahme | `p11-exempt-nicht-am-anfang` |
| `m-p12-block-fehlt` (alle) | Block `vcs` wirkt (Punkt 10/12) | `p12-kern`, `p12-ersetzt-anders`, `p12-geloescht` |
| `m-p12-paths-leer` | `paths` | `p12-kern`, `p12-ersetzt-anders`, `p12-geloescht` |
| `m-p12-paths-fehlt` | Eintrag `paths` | `p12-kern`, `p12-ersetzt-anders`, `p12-geloescht` |
| `m-p12-immutable-when-jeder-status` | nur `Accepted` ist unveränderlich | `p12-vorschlag` |
| `m-p12-immutable-when-fehlt` | Eintrag `immutable-when` (ohne ihn ist jede Datei aus `paths` unveränderlich) | alle acht |
| `m-p12-immutable-when-anderer-status` | Wert `Accepted` von `immutable-when` | `p12-kern`, `p12-ersetzt-anders`, `p12-geloescht` |
| `m-p12-abschnitt-fehlt` | `exclude-sections` | `p12-geschichte` |
| `m-p12-abschnitt-falsch` | Abschnitt `Geschichte` | `p12-geschichte` |
| `m-p12-status-zeile-fehlt` | Eintrag `status-line` | alle außer `p12-geloescht` |
| `m-p12-status-zeile-eng` | `status-line` zieht die Status-Zeile aus dem Core (Muster nur für `Accepted`) | `p12-ersetzt`, `p12-ersetzt-anders` |
| `m-p12-head-allow-fehlt` | `head-allow` | `p12-ersetzt-anders` |
| `m-p12-head-allow-eng` | Element `Superseded by ` | `p12-ersetzt` |
| `m-p12-head-allow-ohne-accepted` | Element `Accepted` | `p12-geschichte`, `p12-index`, `p12-link-ziel`, `p12-vorschlag` (jede Datei mit `Accepted` im Kopf, auch unverändert) |
| `m-p12-link-ziel-fehlt` | `ignore-link-targets` fehlt | `p12-link-ziel` |
| `m-p12-link-ziel-aus` | `ignore-link-targets: false` | `p12-link-ziel` |
| `m-p18-int-falle-fehlt` | Abbruch durch `SIGINT`: die Falle (Aufräumen mit Exit 130) | `p18-sigint` |
| `m-p18-int-kinder-bleiben` | `SIGINT`: Nachkommen werden vor dem Aufräumen beendet | `p18-sigint` |
| `m-p18-int-verzeichnis-zu-frueh` | `SIGINT`: Reihenfolge Nachkommen beenden, dann Verzeichnis entfernen (Mutant entfernt zuerst, beendet nach 5 s) | `p18-sigint` |
| `m-p18-term-falle-fehlt` | dasselbe für `SIGTERM` (Exit 143) | `p18-sigterm` |
| `m-p18-term-kinder-bleiben` | dasselbe | `p18-sigterm` |
| `m-p18-term-verzeichnis-zu-frueh` | dasselbe | `p18-sigterm` |

Kein Mutant blieb grün. `p12-index` ist ein Beleg ohne eigene Mutation (Vertrag Punkt 12: ein weiteres `paths` bricht `immutable-when` und `head-allow` zugleich und färbt alle Fälle; die frühere Mutation `m-p12-paths-index` entfällt). Zwei Formen sind Verhalten von d-check ohne Fall (Grenzen in Punkt 11 und 12): Eine Umbenennung der verlinkten Datei (anderer Name) bleibt `core-drift-vcs`, und die Zeile `Merge y` nach einer Leerzeile (Betreff `x`) ist rot; ich habe beides gemessen (Handbuch von d-check, `(?m)^(Merge |Revert )` färbt nichts), aber kein Fall und keine Zusage trägt es.

**Weitere Sonden (rot gesehen, nicht im Lauf).** (1) Ohne `chmod -R a+rX` und unter `umask 077` wird die unveränderte Kopie rot (alle Fälle von p11 und p12, `d-check: error: kein lesbares git-Repository unter /repo: … permission denied`); die Rechte der Temp-Bäume sind damit Bedingung. (2) `DOCKER=/tmp/dock.sh` (Hülle, die jeden Aufruf zählt): 21 Aufrufe (Vorlauf und 20 Fälle) — `DOCKER` wirkt auf `d-check.mk`; `DOCKER=/gibt/es/nicht` und ein Digest ohne Image enden mit Exit 2 und Meldung. (3) Als root (`sudo`, `LANG=de_DE.UTF-8` fällt auf `C` zurück, dort ist die Locale nicht installiert), als Benutzer mit `LC_ALL=en_US.utf8`, `TZ=Asia/Tokyo` und `umask 077`, mit `GP_PARALLEL=8` und mit dem Vorgabewert 4: jeweils `grün (30 Läufe)`.

**Sonden der Nacharbeit (rot gesehen).** (4) `declare -F`-Zeile der Fallliste in einer Kopie gestrichen, Lauf `fall-ohne-funktion`: `ROT — Fall ohne Funktion: erwartet rot: p11-tippfehler; rot sind: nichts`; mit der Zeile grün. (5) Bei den Mutationen des Abbruchs war die Röte zuerst vom Zeitverhalten abhängig (Docker-Läufe unter Last länger als die Wartezeit des Mutanten, zwei der sechs blieben abwechselnd grün); die Kopie des Skripts läuft deshalb mit einer Attrappe von `docker` (eine Sekunde, Exit 0), `faelle` legt das Verzeichnis des Laufs vor jedem Fall an, und die Prüfung wartet bis zu 8 s auf Reste; danach drei Läufe der sechs Mutationen hintereinander mit derselben roten Menge. (6) Ohne die Hülle `docker` mit `--mount` legte Docker bei einem Mutanten, dessen Nachkommen weiterliefen, die fehlenden Quellpfade als root an und der Aufräum-Lauf meldete `Permission denied`; die Kopie des Skripts bekommt deshalb `PATH` mit der Hülle der äußeren Gegenprobe nicht mehr, sondern die Attrappe als `DOCKER`.

**Messung Zeilenende und Kodierung (Randform 6/7, `SPEC-064` Punkt 17; keine Zusage, kein Fall).** Commit-Message mit `\r\n` (`git commit --cleanup=verbatim`): `x\r\n\r\nslice-abc\r\n` ergibt 0 Befunde (Exit 0); `x\r\n\r\nohne\r\n` ergibt 1 Befund `commit-untraceable` (Exit 1), die Meldung zeigt den Betreff als `x ` (das `\r` ist Teil des Werts, ohne sichtbares Zeichen). ADR mit umgestelltem Zeilenende: reines Umstellen LF→CRLF und CRLF→LF einer `Accepted`-ADR ergibt 0 Befunde (Exit 0), ebenso eine neue Zeile im Abschnitt `Geschichte` einer CRLF-Datei; Textänderung in einer CRLF-ADR (mit und ohne gleichzeitiges Umstellen) ergibt `core-drift-vcs` (Exit 1); Status `Rejected` in einer CRLF-ADR ergibt `unzulässiger Status-Übergang einer immutablen Datei` (Exit 1). In dieser Messung (eine ADR-Datei, ein Aufruf mit d-check v0.86.1) ergab die Umstellung des Zeilenendes keinen `core-drift-vcs`. Für andere Dateien oder Werkzeugstände ist das nicht zugesagt (SPEC-064 Punkt 17).

**Link-Ziel.** Der Fall `p12-link-ziel` zieht eine Datei in ein anderes Verzeichnis bei gleichem Dateinamen um (`open/s.md` → `in-progress/s.md`): grün. Die Umbenennung ist nach Punkt 12 (Grenze) Verhalten von d-check ohne Fall.

**Läufe.** Am Endstand der Nacharbeit (Code-Commit `8d2b4ec3`, seitdem nur Plan) im Dev-Container pgr-dev3, nacheinander, ohne parallele volle Läufe: `make range-konfiguration-gegenprobe` `grün (43 Läufe)` (2 min 32 s); `make gates-schnell-gegenprobe` grün (1 min 00 s); `make gates` Exit 0 (4 min 20 s; Nachweis geschrieben); `make docs-check` 0 Befunde. Die Läufe der Abschnitte oben (Stand `f6e346ed`: `make gates` 4 min 23 s, Gegenprobe als root und unter `umask 077`) gelten für den früheren Stand.

**Metrik gegen die Marke von §4** (`git diff --numstat eac90f75..8d2b4ec3 -- tools harness .d-check.yml Makefile`): Summe aus Einfügungen und Löschungen **433** (Gegenprobe 386, `.d-check.yml` 19, Fragment 11, `harness/README.md` 1, `SPEC-064` 16 = die Änderung des Architect `313d6459`, die mitzählt; ohne sie 417) gegen die Marke 500 und die Nachschätzung des Architect 365 bis 415: unter der Marke, die Schätzung um 2 Zeilen (ohne Vertrag) bzw. 18 (mit Vertrag) überschritten; keine Rückführung, die Meldeschwelle 450 nicht erreicht.

**Kommentare (Schritt 20).** Neu geschrieben: Kopf und Zeilen der Gegenprobe, Kommentar zu `.d-check.yml` (Zusage: Fehlt ein Block oder ein Eintrag, bei `paths` auch leer, wird `range-konfiguration-gegenprobe` rot; Test: `m-p11-block-fehlt`, `m-p12-block-fehlt`, die Mutationen „Muster fehlt“, `m-p11-ausnahme-fehlt`, `m-p12-paths-fehlt`, `m-p12-paths-leer`, `m-p12-immutable-when-fehlt`, `m-p12-abschnitt-fehlt`, `m-p12-status-zeile-fehlt`, `m-p12-head-allow-fehlt`, `m-p12-link-ziel-fehlt`), Kopf des Fragments. Jede Zusage dort hat einen Fall oder eine gesehene Sonde (oben); „Kein Fall hängt von Dateirechten ab“ ist durch den Lauf als root und unter `umask 077` belegt, „die Temp-Bäume sind lesbar … weil das Image als Nicht-root liest“ durch die Sonde (1).

**Material für die Closure (Beobachtungen).** Vier meiner Erwartungen vor dem ersten Lauf mit `GP_ENTDECKE=1` waren falsch, weil `vcs` eine Konfiguration anders auswertet, als sie sich liest (Treffer `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, §8): `head-allow` prüft jede unveränderliche Datei im Kopf, auch eine unveränderte (`m-p12-head-allow-ohne-accepted`, `m-p12-paths-index` färbten mehr Fälle als gedacht); `immutable-when` auf jeden Status machte die unveränderte `Proposed`-ADR in jedem Fall rot, bis `0002` nur noch im Fall `p12-vorschlag` angelegt wurde; ein `status-line`-Muster, das nichts trifft, und ein fehlender Eintrag `status-line` oder `immutable-when` färben fast alle oder alle Fälle (`vcs` schlägt dann fehl statt still zu schweigen). Dazu ein Fehler des Geschirrs: Die Hülle für `DOCKER` rief sich bei `DOCKER=docker` selbst auf (PATH-Rekursion) und lief ohne Ende, bis sie von Hand beendet wurde; sie führt `DOCKER` jetzt mit dem PATH des Aufrufers aus (`GP_PFAD`). 

### Nacharbeit nach dem Review (`docs/reviews/2026-10-18-review-slice-harness-ci-range-konfiguration.md`)

Ausgang je Finding (die Entscheidung des Architect steht in `SPEC-064` Punkte 11, 12, 16, 18; Commit `313d6459` des Architect, Commit `8d2b4ec3` hier):

| Finding | Ausgang |
|---|---|
| F-1 MEDIUM Form der Muster, Anker je Muster | behoben: fünf Fälle `p11-*-knapp`, die Fälle `p11-adr` bis `p11-welle` tragen die Kennung mitten in der Zeile, je Muster eine Mutation Weitung (4) und Anker (4) statt `m-p11-muster-verankert`; Befund als Klasse: das Merkmal ist die Form des Musters (Länge der Ziffernfolge, Zahl der Buchstaben, Mindestlänge des Namens) und der Anker, je Ausprägung eine Mutation, nicht nur die gemeldete (Tabelle oben) |
| F-2 MEDIUM Abbruch hinterlässt Temp-Baum und Nachkommen | behoben nach Punkt 18: Fallen für `SIGINT` und `SIGTERM` (`abbruch_int`, `abbruch_term`: Nachkommen anhalten, beenden, abwarten; dann Verzeichnis entfernen; Exit 130 bzw. 143), `EXIT` räumt nur auf; Fälle `p18-sigint`, `p18-sigterm`, je drei Mutationen des Skripts (Tabelle oben); `SIGKILL` und `SIGHUP` nicht zugesagt, ein laufender Container von d-check endet von selbst (Punkt 18) |
| F-3 LOW `m-p12-paths-index` | gestrichen; `p12-index` ist Beleg ohne eigene Mutation (Punkt 12) |
| F-4 LOW Name ohne Funktion zählt als grün | behoben: `faelle` zählt einen Namen ohne Funktion als roten Fall; Lauf `fall-ohne-funktion` (Fall) und Sonde (4) (Mutation: Prüfung gestrichen, rot gesehen); die Mutation läuft nicht in der Gegenprobe, weil sie das Skript selbst ändern müsste und `m-p18-*` diese Form schon für das Skript trägt — Hinweis an den Reviewer |
| F-5 LOW Versionen im Kopf | Kopf berichtigt: nennt `wait -n` (bash), nicht mehr „bash ≥ 4“, „git ≥ 2.32“ und sagt, dass das Skript keine Version prüft; die Plattform steht in `SPEC-064` |
| F-6 INFO Vertragslücke Fallliste | vom Architect in die Liste aufgenommen (Punkte 11 und 12); keine Arbeit |
| F-7 INFO Umbenennung, „Merge nach Leerzeile“ | als Grenze in Punkt 11 und 12; ohne Fall (siehe oben, Absatz nach der Tabelle) |
| F-8 INFO Schlusssatz der CRLF-Messung | auf die Messung begrenzt, Wortlaut des Architect (Absatz oben) |
| F-9 INFO Plan-Nachzug | §3 nachgezogen (27 Fälle, 41 Mutationen); `make kopf-check` im `make gates` grün |
| F-10 INFO Worktree ohne `.git`-Verzeichnis | Eigenschaft der Treiber, keine Arbeit; die Läufe hier im Klon des Containers |

### Nacharbeit nach der Verifikation (`docs/reviews/2026-10-18-verify-slice-harness-ci-range-konfiguration.md`)

| Befund | Ausgang |
|---|---|
| V-1 MEDIUM Rekursion bei `GP_MUT=m-p18-*` | behoben: Die Kopie des Skripts startet mit `env -u GP_MUT -u GP_ENTDECKE -u GP_FALL` und `GP_PARALLEL=2`; `GP_INNEN=1` schaltet in jedem Lauf der Kopie die p18-Fälle (`faelle p18` liefert nichts) und die Mutationen `m-p18-*` (`mut18` kehrt sofort zurück) ab. Fall: Lauf `rekursion-schutz` (eine Kopie mit `GP_INNEN=1` und `GP_MUT=m-p18` fährt nichts, Exit 1, Meldung „keine Mutation und keine Kopie gefahren“, unter `timeout 60`). Mutation (Sonde, nicht im Lauf): die Zeile `[ -z "${GP_INNEN-}" ] \|\| return 0` in `mut18` gestrichen, dann endet der Lauf `rekursion-schutz` nicht mit Exit 1 (rot, die Meldung enthält die Läufe der inneren Kopie). |
| V-4 LOW Kommentar zu `GP_INNEN` | Kopf berichtigt (Wortlaut oben im Skript): nennt `env -u …`, `GP_PARALLEL=2` und die abgeschalteten Fälle und Mutationen |
| V-2 MEDIUM, V-3 INFO Zeichenklassen, Schwanz, Leerzeichen | keine Mutation, kein Fall: Der Architect hat Zeichenklasse der Muster und Leerzeichen der Ausnahme als Grenze in `SPEC-064` Punkt 11 genannt; zugesagt sind Stellung, Länge und Ziffernzahl. Kein Satz in Fragment, Skriptkopf oder README-Zeile sagt „je Zeichenklasse“ zu (durchgesehen; der Kopf nennt die Grenze jetzt). Der Schwanz `(-[a-z0-9]+)*` des Slice-Musters ist für die Anwesenheitsprüfung äquivalent: Invariante: ein Treffer des längeren Musters enthält einen Treffer des kürzeren; Ort: die Suche nach dem Muster in der Zeile, kein Ganzzeichen-Vergleich; Sonde (Verifikation, Mutant `v-schwanz`): kein Fall wird rot, Original und Mutant trennen sich nicht |

**Läufe nach der Verifikation** (Stand `c1d2b1f3`, nacheinander): `make range-konfiguration-gegenprobe` `grün (44 Läufe)` (2 min 28 s; 41 Mutationen, `ok`, `fall-ohne-funktion`, `rekursion-schutz`); `make gates-schnell-gegenprobe` grün (59 s); `make gates` Exit 0 (4 min 21 s); `make docs-check` 0 Befunde vor jedem Commit.

**Messung V-1** (Container pgr-dev3, jede Zeile unter `timeout`, Prozesse mit `ps -e` gezählt, alle 0,3 s abgetastet): `GP_MUT=m-p18-int-falle-fehlt GP_ENTDECKE=1 timeout 100 bash tools/harness/range-konfiguration-gegenprobe.sh`: `m-p18-int-falle-fehlt => p18-sigint` und `grün (1 Läufe)`; Prozesse vorher 16, Spitze 53, nachher 16 (die Kopie fährt ihre Läufe mit der Attrappe von `docker`; es bleibt kein Prozess des Skripts zurück). Vorher, im Stand `d8086425`, lief dieselbe Zeile ohne Ende.

### Closure-Notiz

- **Was hat funktioniert:** `make doc-commits` und `make doc-immutable` prüfen in diesem Repo etwas: Über dasselbe Temp-Repository melden sie ohne die Blöcke 0 Befunde und mit ihnen je einen (Tabelle *Befundzeilen vorher und nachher*). Die Randformen standen vor dem Code in `SPEC-064` (Punkte 10 bis 12, 16, 17), die Rückgaben während der Arbeit gingen an den Architect (Form der Muster, Fallliste, Abbruch, Grenzen: Punkte 11, 12, 18). Die Gegenprobe führt mit dem echten Image 27 Fälle und 41 Mutationen; Review (0 HIGH, 2 MEDIUM) und Verifikation (alle Mutationen am Host nachgefahren, frischer Klon mit Probe-Commits) fanden keinen Fehler in der Konfiguration, nur Lücken der Gegenprobe und einen Fehler ihres Geschirrs (V-1).
- **Was ging anders als geplant:** (1) Der Implementer erwartete vier rote Mengen falsch, weil `vcs` die Konfiguration anders auswertet, als sie sich liest (Treffer, siehe Register unten). (2) Die Gegenprobe wuchs durch das Review: Form und Anker der Muster (F-1), Abbruch (F-2), Rekursionsschutz (V-1); gemessen 445 gegen die Schätzung 300 bis 405 (Marke 500 nicht erreicht). (3) Der Plan folgte der Nacharbeit nicht: Kopf, §1 und DoD Liefer-Punkt 2 trugen Punkt 18, die Form der Muster und zwei Läufe nicht (P-1), weil die Nacharbeit nur §3 und §7 ändern durfte; bei der Closure nachgezogen. (4) Die Belege der Range-Prüfungen ließen sich im Worktree des Reviews nicht wiederholen (F-10).
- **Steering-Loop-Eintrag:** *neuer Sensor*. Das Ziel `make range-konfiguration-gegenprobe` — liegt in `harness/mk/range-konfiguration.mk` (`GEGENPROBEN_CONTAINER`) und `tools/harness/range-konfiguration-gegenprobe.sh`, Vertrag `SPEC-064` Punkte 11, 12, 16 und 18, Zeile mit dem Anker `· seit slice-harness-ci-range-konfiguration` in `harness/README.md` §Sensors. Es fängt die Klasse „ein Block der Range-Prüfungen fehlt, ist leer oder zu weit gefasst“: Damit sind `make doc-commits` und `make doc-immutable` nicht mehr über nichts grün; ohne die Gegenprobe färbte das Löschen der Blöcke keinen Lauf. *Benannte Spec-Lücke:* keine offen; die Lücken des Reviews und der Verifikation (Form der Muster, Fallliste, Abbruch, Zeichenklassen) hat der Architect in `SPEC-064` Punkte 11, 12 und 18 geschlossen oder als Grenze genannt. *Geschärfte Regel:* nicht in diesem Slice; die Klasse „Konfiguration wirkt anders als gelesen“ steht bei 3× und hat den Ausgang `geplant` (`slice-harness-regeln-gate-konfiguration`, Vorschlag des Planners, unten). Auslöser: `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen` (`slice-walking-skeleton-build-gates`, `slice-harness-lint-werkzeug`, `slice-harness-ci-range-konfiguration` — 3×).
- **Beobachtungs-Register (`../observations/`):** je ein Beleg `evidence/slice-harness-ci-range-konfiguration.md` in `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen/` (Zähler 3, Stand `offen` → `geplant`, Kennung `slice-harness-regeln-gate-konfiguration`), `gegenprobe-haengt-an-umgebung-des-hosts/` (Zähler 7; Ausgang bleibt beim Nutzer), `gegenprobe-startet-sich-ueber-einen-fall-selbst/` (V-1; Zähler 2, `offen`), `negativtests-fehlen-bei-neuem-vertrag/` (F-1, V-2; 57), `spec-randform-erst-im-review-entschieden/` (F-2; 42), `zusage-im-kommentar-weiter-als-pruefung/` (F-5, F-8, V-4; 84), `plan-folgt-korrektur-nicht/` (P-1, F-9; 69), `schaetzung-in-zeilen-unterschaetzt-tests/` (17, `geplant` → `slice-harness-schaetzung-planner-faktor`) und `zeilenende-crlf-im-gelesenen-wert/` (Messung; 2, `offen`); neu mit einem Beleg: `BEO-REPO/range-pruefung-laeuft-nicht-im-worktree/` (F-10; 1×, `offen`). Vermerk ohne Beleg in `gruener-mutant-faelschlich-aequivalent` (V-3: der Schwanz des Musters nennt Invariante, Ort und Sonde, die Verifikation bestätigte ihn; Zähler 5).
- **Folge-Slices:** `slice-harness-regeln-gate-konfiguration` (Regel zu `BEO-REPO/gate-konfiguration-wirkt-anders-als-gelesen`, ohne Welle — ist eine Datei in `open/`; **Vorschlag des Planners**, vom Nutzer nicht bestätigt, nimmt die Sendung des Registers in §1 an) und `slice-harness-ci` (der Nachfolger in `welle-erster-release`, dessen Start-Trigger „`slice-harness-ci-range-konfiguration` in `done/`“ mit diesem Slice erfüllt ist — ist eine Datei in `open/`; `grep -n slice-harness-ci-range-konfiguration` findet die Abgrenzung in dessen §1, Start, §3 und §8).
- **Ausgang je Finding** (Review F-1 bis F-10, Verifikation V-1 bis V-4 und P-1; Maßnahmen und Commits in den Abschnitten *Nacharbeit* oben):

| Finding | Ausgang |
|---|---|
| F-1 MEDIUM, V-2 MEDIUM | entfallen: Form und Anker je Muster haben Fälle `-knapp` und je acht Mutationen (`8d2b4ec3`); die Zeichenklassen hat der Architect als Grenze in `SPEC-064` Punkt 11 genannt (`816da910`), kein Satz der Gegenprobe sagt sie zu. Klasse im Register (`BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`) |
| F-2 MEDIUM | entfallen: Punkt 18 und Fälle `p18-*` mit je drei Mutationen des Skripts (`8d2b4ec3`); `SIGKILL` und `SIGHUP` nicht zugesagt (Grenze in Punkt 18). Klasse im Register (`BEO-REPO/spec-randform-erst-im-review-entschieden`) |
| F-3 LOW | entfallen: `m-p12-paths-index` gestrichen, `p12-index` ist Beleg ohne eigene Mutation (Punkt 12) |
| F-4 LOW | entfallen: ein Name ohne Funktion ist ein roter Fall, Lauf `fall-ohne-funktion` (`8d2b4ec3`); die Mutation der Prüfung läuft nicht in der Gegenprobe, die Sonde (4) hat sie rot gesehen |
| F-5 LOW | entfallen: Kopf berichtigt (`8d2b4ec3`). Klasse im Register (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`) |
| F-6 INFO, F-7 INFO | entfallen: Der Architect hat die Fallliste ergänzt und die Umbenennung und „Merge nach Leerzeile“ als Grenze genannt (`313d6459`); keine Arbeit |
| F-8 INFO | entfallen: Schlusssatz der CRLF-Messung auf die Messung begrenzt (`00e072bf`); die Messung ist im Register (`BEO-REPO/zeilenende-crlf-im-gelesenen-wert`, 2×) |
| F-9 INFO, P-1 LOW | entfallen: §3 mit der Nacharbeit nachgezogen (`00e072bf`), Kopf, §1 und DoD Liefer-Punkt 2 bei der Closure (Punkt 18, Form der Muster, die Läufe `fall-ohne-funktion` und `rekursion-schutz`); `make kopf-check` und `make docs-check` grün. Klasse im Register (`BEO-REPO/plan-folgt-korrektur-nicht`) |
| F-10 INFO | weiter offen → `BEO-REPO/range-pruefung-laeuft-nicht-im-worktree` (1×, neu): `d-check.mk` bindet nur das Arbeitsverzeichnis ein, in einem Worktree ist `.git` eine Datei. Keine Adresse; der Nutzer entscheidet, ob `d-check.mk` den Fall trägt oder die Rollen den frischen Klon als Voraussetzung nennen |
| V-1 MEDIUM, V-4 LOW | entfallen: `env -u`, `GP_INNEN=1` in der Kopie, Lauf `rekursion-schutz`, Kopf berichtigt (`c1d2b1f3`); Messung der Prozesszahl vorher 16, Spitze 53, nachher 16. Klasse im Register (`BEO-REPO/gegenprobe-startet-sich-ueber-einen-fall-selbst`, 2×) |
| V-3 INFO | entfallen: Der Schwanz `(-[a-z0-9]+)*` ist für die Anwesenheit äquivalent (Invariante, Ort und Sonde in *Nacharbeit nach der Verifikation* oben, von der Verifikation bestätigt); das Leerzeichen der Ausnahme ist Grenze in Punkt 11 |

- **Risiken aus §6:** alle fünf *entfallen* mit Begründung, siehe §6. Risiko 2 trägt die Messung des Implementers (`368e3d43..f6e346ed`) und der Verifikation (`368e3d43..d8086425`, frischer Klon); Risiko 4 trägt `AGENTS.md` §6 Schritt 6.
- **Drei Paarungen:** Anker: Der Sensor liegt in `harness/mk/range-konfiguration.mk` und `tools/harness/range-konfiguration-gegenprobe.sh` mit dem Anker `· seit slice-harness-ci-range-konfiguration` in `harness/README.md` §Sensors; eine Regel in `AGENTS.md` wird nicht verkörpert, `liegt in` dafür entfällt. Folge-Slice: `slice-harness-regeln-gate-konfiguration` und `slice-harness-ci` liegen in `open/` (Annahme geprüft nach `AGENTS.md` §3.13: der erste nimmt die Sendung des Registers in §1 und DoD an und zählt in §8 nach; der zweite nennt die Abgrenzung zu diesem Slice in §1 und die Pins als eigene Sendung von `slice-harness-ci-workflows`). Register: die Einträge oben, jedes Verzeichnis mit nicht leerem `evidence/`.
- **Auf den Host wartet:** (1) `make gates` am Integrationsstand (der Beleg hier gilt für den letzten Code-Commit `c1d2b1f3`, 4 min 21 s, Exit 0); (2) `make doc-commits` und `make doc-immutable` über die Range am Integrationsstand (in diesem Worktree nicht fahrbar, F-10); (3) `make gegenproben` am Ende der Welle (Risiko 4); (4) die Entscheidung des Nutzers zum Folge-Slice `slice-harness-regeln-gate-konfiguration` und zu Frage 1 aus §6 (Schnitt aus `slice-harness-ci`, bisher Empfehlung des Planners).

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

**Nachschätzung des Architect (2026-10-18, nach der Prüfung; Metrik von §4):** Die Zusagen aus `SPEC-064` Punkte 11, 12 und 18 fügen 5 Fälle `-knapp`, 8 Mutationen (Weitung und Anker je Muster), 2 Fälle mit 3 Mutationen für den Abbruch und eine Falle hinzu; der Implementer hat 285 gemessen. Geschätzt kommen 80 bis 130 Zeilen hinzu, zusammen 365 bis 415, unter der Marke 500.
