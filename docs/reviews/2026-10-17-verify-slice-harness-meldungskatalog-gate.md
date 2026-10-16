# Verifikation slice-harness-meldungskatalog-gate

**Rolle:** Verifier (Modul 11). **Datum:** 2026-10-17. **Gegenstand:** Diff `8d14a918..9c22cf8a`
(Code `0e581db3`, `2ab2ee55`; Plan `21e60054`, `9c22cf8a`; Review `f716b83c`). Plan:
`docs/plan/planning/in-progress/slice-harness-meldungskatalog-gate.md`. Bezug:
[`LH-QA-05`](../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md),
[`SPEC-054`](../../harness/vertraege/SPEC-054-meldungskatalog.md), `SPEC-052` Punkte 4, 13, 21.

## Ergebnis

**Bestätigt, mit zwei LOW-Befunden am Plan (V-1, V-2) und keinem am Code.** Alle Liefer-Punkte der DoD
sind gegen den Diff und gegen eigene Läufe belegt; offen sind nur die Closure-Pflichten, die der Plan
ausdrücklich der Closure zuweist.

## DoD gegen Diff und Lauf

| DoD-Punkt | Befund | Beleg |
|---|---|---|
| Vertrag `SPEC-054` vor dem ersten Code-Commit, entscheidet jede Randform aus §6 | erfüllt | liegt in `8d14a918` (Basis); die später nachgetragenen Punkte 1, 6, 8 (`db77c937`, `8871ebb8`) sind Architect-Entscheidungen auf Rückgabe, kein Implementer-Entscheid |
| Prüfskript bash/awk, am Fragment, an `GATE_CHECKS`, am Stand des Starts grün, `fehler.go`/`SPEC-034` unverändert | erfüllt | `make meldungskatalog-check` Exit 0 ohne Ausgabe; `git diff --stat 8d14a918..HEAD -- internal spec docs/user cmd test .claude .githooks Makefile` leer; `SPEC-052` unverändert |
| in schneller Kette; Gegenprobe in `GATE_SCHNELL_AUSNAHME` und `LANGSAM` | erfüllt | Diff je ein Wort in `gates-schnell.mk` und `gates-schnell-gegenprobe.sh`; `make gates-schnell-gegenprobe` grün (Deckung der Ausnahmeliste) |
| Gegenprobe je Punkt und je Element, Anschluss-Fall mit Makefile | erfüllt | `make meldungskatalog-check-gegenprobe` grün, Fälle `nr1`–`nr11`; keine Fälle übersprungen (Nutzer, nicht root; `grep -c übersprungen` = 0) |
| Tabelle Zusage · Mutation · roter Test in §7 | erfüllt | 139 Mutanten im Läufer nachgezählt (`len(M)` = 139), 137 rot, 2 äquivalent |
| Sensor-Datei, zwei Zeilen in `harness/README.md` §Sensors mit Bindung | erfüllt | `harness/sensors/meldungskatalog-check.md` (75 Zeilen), zwei Zeilen mit [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md) und `SPEC-054` |
| `make gates` grün | erfüllt, nicht wiederholt | `.harness/state/gates-passed.diffsha` 01:37, nach `2ab2ee55` (01:27); seither ändert nur `9c22cf8a` eine Datei unter `docs/plan/planning/` (Positivliste). Code-Stand identisch, Stop-Hook antwortet `approve` |
| Review durchgeführt | erfüllt | `docs/reviews/2026-10-17-review-slice-harness-meldungskatalog-gate.md`; F-1 bis F-3 nachgearbeitet (siehe unten) |
| Closure-Notiz, Lerneintrag, Register, Risiko-Ausgänge, Paarungen | offen (Closure) | §7 trägt „Wird bei Closure gefüllt"; DoD-Häkchen stehen noch nicht — kein Befund, Vorschlag unten |

**§8 vorhanden** (F-1 behoben): `fe68464c`-Fassung wiederhergestellt, drei „Nachgezählt"-Absätze, der letzte für die Nacharbeit
(drei Liefer-Punkte, eine Schicht, 607 gegen 900/750). Siehe V-1 zum Inhalt der Sichtung.

## Review-Befunde nachgeprüft

- F-1 (§8 gelöscht): behoben, Abschnitt vorhanden und nachgezählt.
- F-2 (vier Elemente ohne Mutation): behoben; gebrochen und rot gesehen, siehe nächster Abschnitt.
- F-3 (stille Randformen): Tabulator und CRLF stehen als akzeptierte Negative in `SPEC-054` (Punkt 9, Grenze); der Symlink ist
  eine laute Sperre (Punkt 8) mit Fällen `nr8-symlink-*`. Kein stilles Grün.

## Bewusstes Brechen (eigene Stichprobe, 24 Mutanten, eigener Läufer)

Je Mutant eine Kopie von Skript, Gegenprobe, Makefile und Fragmenten, die Gegenprobe aus leerem Verzeichnis. Rot und damit gefangen:

- Symlink-Sperre: `-type l` aus dem `find` (rot: `nr8-symlink-datei`, `-erste-sortiert`, `-name`, `-startpunkt`, `-verzeichnis`); Sperrblock
  entfernt (dieselben); Reihenfolge (nicht lesbar vor Symlink) über die Tabelle des Plans, hier nicht getauscht.
- Die vier F-2-Elemente: Abschnittsende nur `###` (`nr1-k-ende-ebene-2`) bzw. nur `##` (`nr1-h-ende-ebene-1`); Zelle Tab vor/nach
  (`nr1-k-zelle-tab-vor`, `-nach`) und leer (`nr1-k-zelle-leer`); Zählung nach Formprüfung für C (`nr1-c-nur-form-zeilen`); `sort` der
  Bytefolge entfernt, nicht lesbar (`nr8-unlesbar-erste-sortiert`) und Symlink (`nr8-symlink-erste-sortiert`).
- Mengenregeln 4a, 4b (Endung `000` entfernt), 4c, 4d; Klasse 3 übersprungen (`nr5-klasse-3`); Endanker der Form (`nr2-form-zu-lang`);
  doppelt in K (`nr3-doppelt-k`).
- Literal `"PGR-` ohne Anführungszeichen gesucht (`nr6-*`, 8 Fälle rot); `_test.go` gelesen (`nr6-test`); Sperre in der Shell mit Exit 1
  (`nr8-*`).

**Überlebt, beide wie im Plan als äquivalent geführt:**

- `8-exit-awk-sperre` (awk endet bei Sperre mit 1): Invariante und Ort genannt (`if ! befunde="$(awk …)"; then exit 2`). Nachgefahren:
  jeder awk-Ausgang ungleich 0 wird dort zu Exit 2; der Mutant am Ort (`exit 1` im `then`) ist rot (`nr8-*`, 12 Fälle). Äquivalenz trägt.
- `8-ohne-internal` (`internal/ fehlt` entfällt): C liegt unter `internal/`, Punkt 8 prüft C zuerst; ein Baum ohne `internal/` endet
  vorher mit `<C> fehlt`. Keine Eingabe erreicht die Zeile (auch kein `internal` als Datei, `-f` auf C scheitert vorher). Trägt;
  die Aussage steht im Vertrag (Punkt 8) und in der Gegenprobe (Kopf). Register `gruener-mutant-faelschlich-aequivalent` (5×):
  beide Einordnungen halten der Nachfahrt stand, kein weiteres Auftreten.

## Sensoren (selbst gefahren, nicht aus der Behauptung)

| Lauf | Ergebnis |
|---|---|
| `make meldungskatalog-check` | Exit 0, keine Ausgabe (C 30, K 34, H 30 nach Vertrag Punkt 11) |
| `make meldungskatalog-check-gegenprobe` | grün, 2 s |
| `make gates-schnell-gegenprobe` | grün, Ausnahmeliste deckt `GATE_CHECKS` |
| `make gates-schnell` | Exit 0 |
| Stop-Hook `.claude/hooks/stop-require-gates.sh` | `approve` |
| `make docs-check` (vor dem Commit) | 939 Dateien, 0 Befunde |
| `make gates` | nicht wiederholt, Code-Stand identisch zum Nachweis (s. o.) |

**Gegenprobe liest den Arbeitsbaum nicht** (`SPEC-052` Punkt 13): Kopf der Gegenprobe nennt die Regel; sie baut jeden Baum selbst und
kopiert für Punkt 11 nur Makefile und Fragmente. Plan §7 belegt den Lauf aus einer Kopie ohne `spec/`, `docs/`, `cmd/`, `internal/`; ich habe diesen Lauf nicht wiederholt.

## Aussagen der Sensor-Datei (§3.11)

Die Datei sagt zu: Ausgabe, Exit 0/1/2, Sperren, Reihenfolge, Grenze (Spalten, Fundorte, Verhalten, Handbuch, Codeblöcke). Jede
Zusage hat Fall oder ist als akzeptiertes Negativ im Vertrag benannt. Die Aussage zu übersprungenen Fällen („Grün heißt jeder Fall, der
lief; als root … laut übersprungen, Exit 0") stimmt mit den Zeilen 183, 193, 204 und 248–251 der Gegenprobe überein. Die README-Zeilen
nennen „1 bis 6 und 8 bis 10" wie der Vertrag; Punkt 11 ist der Makefile-Fall.

## Unveränderter Bestand

`git diff --stat 8d14a918..HEAD` berührt nicht: `internal/hexagon/model/fehler.go`, `spec/`, Handbuch, `.claude/` (Stop-Hook), `.githooks/`,
Positivliste, `SPEC-052`. Berührt sind nur Plan, Roadmap, Review, `harness/`, `tools/harness/`.

## Messung nachgerechnet

`git diff --numstat 8d14a918..2ab2ee55 -- tools harness` = **607** (148 + 308 + 75 + 18 + 52 + 2 + 2 + 2). Die 52 Zeilen `SPEC-054` sind
`39 + 13` aus den Architect-Commits; eigene Zeilen 555. Marke 900 und Meldeschwelle 750 nicht erreicht; die Schätzung 450 bis 750 trifft.

## Findings

| Nr. | Schwere | Befund | Zuständigkeit |
|---|---|---|---|
| V-1 | LOW | §8 *Offene Beobachtungen sichten* nennt Stand „gemergter Stand, 2026-10-07" und Zähler (13×, 17×, 11×, 3×), die das Register heute anders führt (Dateien in `evidence/`: 44, 68, 37, 13), und nennt `BEO-REPO/gruener-mutant-faelschlich-aequivalent` (5×) nicht, obwohl §7 zwei Mutanten als äquivalent einstuft. Der Plan wurde 2026-10-17 nachgezählt, die Sichtung nicht. | Implementer bei der Closure: Datum und Zähler berichtigen oder „Zähler siehe Register" schreiben, den fehlenden Eintrag aufnehmen und die Nachfahrt des Verifiers (oben) als Beleg nennen |
| V-2 | LOW | §6 und §7 nennen „Adresse aus `slice-v1-abschluss-container`" und die Folgelast als Risiken ohne Ausgang; der Nehmer nennt das Gate in §1 (`grep` findet es), die Tabellen `### Fehlercodes`/`### Warnungen` tragen dort aber keinen eigenen DoD-Punkt. Kein Mangel dieses Slice, aber der Ausgang *entfallen* braucht diese Begründung. | Implementer bei der Closure (Ausgänge, Vorschlag unten); Planner, falls der Container-Slice die Bindung als Liefer-Punkt führen soll |

Kein Befund am Code, am Vertrag oder an der Sensor-Datei.

## Vorschlag der Risiko-Ausgänge (zur Closure, nicht gesetzt)

1. **Bestand** → *entfallen*: Stand des Starts grün (`make meldungskatalog-check` Exit 0, von mir am Endstand wiederholt), `fehler.go`
   und `SPEC-034` unverändert, kein Rot ist aufgetreten.
2. **Folgelast für welle-v1-abschluss** → *entfallen* mit Begründung: Das Gate meldet jede Abweichung ab dem ersten Lauf in
   `make gates` und `make gates-schnell`; eine Folge-Slice-ID ist nicht nötig. Alternativ *weiter offen*, wenn der Nutzer die
   Pläne der Welle ändern lassen will.
3. **Adresse aus `slice-v1-abschluss-container`** → *eingetreten* nicht; *entfallen*: Das Gate hält beide Handbuch-Tabellen
   (Punkte 1, 4c, 4d), der Nehmer nennt es in §1.
4. **Lesart ohne Entscheidungsort** → *weiter offen* ins Register (`BEO-REPO/harness-lesart-ohne-entscheidungsort`, mit diesem Slice
   dritter Beleg: die Liste der Randformen entschied der Architect auf Rückgabe, nicht der Code). Das ist die Evidenz-Datei der
   Closure, kein Zähler.

Beobachtungen für das Register nach Plan §7: `gate-uebergeht-ablage-eintrag-still` (nicht lesbares Verzeichnis, jetzt Sperre),
Vorgabe des Architect ohne Messung (`LANGSAM`), Mutanten mit Aufbau aller Gutfälle.

## Offen für den Planner

Nichts, was den Schnitt oder die Entscheidungen berührt. Spec-Lücke: keine.
