# Review-Report: slice-harness-meldungskatalog-gate — 2026-10-17

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten); der Gegenstand sind Harness-Werkzeuge (Prüfskript, Gegenprobe, Fragment, Sensor-Datei), geprüft gegen den Vertrag `SPEC-054`.

**Gegenstand:** `slice-harness-meldungskatalog-gate` (Schritt 43, wellenlos), Diff `8d14a918..21e60054` (Pfade `tools/harness`, `harness`, dazu Plan und Roadmap); die Architect-Commits `8d14a918` (Vertrag `SPEC-054`) und `db77c937` (Entscheidung der acht Randformen) sind der Maßstab.

**Skill:** `.harness/skills/reviewer.md` @ c5c6ec08 · <!-- d-check:ignore (Skill-Pfad des Repos, Zitat des Laufs) -->
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-17

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-harness-meldungskatalog-gate` (§1, §3, §4, §6, §7 Belege des Implementers; §8 fehlt am Stand des Laufs, siehe F-1)
- `SPEC-054` (Vertrag, Punkte 1 bis 11 und Grenze), `SPEC-052` (Punkt 4 Ausnahmeliste, Punkt 13 Handbuch als Quelle, Punkt 21 `-i`), [ADR-0011](../plan/adr/0011-meldungscodes-praefix-pgr.md)
- `LH-QA-05`
- `AGENTS.md` (Hard Rules, §3.6, §3.7, §3.9, §3.10, §3.11, §3.12, §3.13, §3.14)
- Muster: `harness/sensors/lint.md`, `harness/sensors/kopf-check.md`
- Eigene Proben (bash ohne Docker, ein Lauf zur Zeit, als Nutzer ohne root): `make meldungskatalog-check` (Exit 0, ohne Ausgabe), `make meldungskatalog-check-gegenprobe` (grün), `make gates-schnell-gegenprobe` (grün); alle 124 Mutanten des Läufers aus §7 nachgefahren (121 rot, 3 grün: die zwei benannten äquivalenten und `8-wurzel-ohne-sperre`, siehe F-3); rund 20 eigene Mutanten und Randform-Bäume aus Kopien von `fehler.go`, `SPEC-034` und Handbuch; die Gegenprobe in einer Kopie ohne `spec/`, `docs/`, `cmd/`, `internal/` (grün)

---

## Findings

Jedes Finding folgt dem **§Output-Schema des Reviewer-Skills** — der
verbindlichen Single Source of Truth. Die Spalten unten sind nur
**gespiegelt** (Bequemlichkeit beim Ausfüllen), nicht neu definiert; bei
Abweichung gilt der Skill bzw. dessen Quelle
`v6.18.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill
— Tag einsetzen, denn diese Zeile wandert in den eingefrorenen Report.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | §8 (Sub-Area-Prüfungen, Modus-Begründung und die Nachzählung der Liefer-Punkte und Schichten) steht am Stand `21e60054` nicht mehr im Plan: Der Commit löscht den Abschnitt (rund 45 Zeilen im Diff `fe68464c..21e60054`, darunter die Sichtung des Registers und zwei Nachzählungen), §3 verweist noch auf ihn. Der Abschnitt entfällt nach der Baseline nie. | `v6.18.0` · `regelwerk/modul-05-planning-harness.md` §Zwei Schritte vor der Modus-Begründung; `AGENTS.md` §3.9, §3.13 (Nachzählen) | `docs/plan/planning/in-progress/slice-harness-meldungskatalog-gate.md` · „diese Liste liefert die **Pfad-Kandidaten** für §8“ (§3); `## 8.` kommt nicht mehr vor | ja — `grep -c '^## 8\.'` auf dem Plan liefert 0, auf `fe68464c` 1 | Abschnitt des Plans beim Umschreiben von §7 verloren |
| F-2 | MEDIUM | Zusage-Elemente von `SPEC-054` haben keine Mutation, die nur sie bricht; vom Review gefahren und grün: (a) Ende des Abschnitts bei jeder Zeile mit `#` — ein Ende nur bei `###` (K) bzw. `##` (H) bleibt grün, die Fälle nutzen nur `###`; (b) Leerraum in der ersten Zelle (nur Leerzeichen statt Leerzeichen und Tabulator zu trimmen bleibt grün) und die leere Zelle ``` `` ```; (c) eine Quelle, deren Zeilen alle `Form`-Befunde sind (Zählung erst nach der Formprüfung bleibt grün: Befund `Form` oder Sperre `ohne Code(s)` ist ungetestet); (d) `sort` in der Sperre „der erste Pfad in Bytefolge“ — ohne `sort` bleibt grün, weil die Reihenfolge von `find` im Fall zufällig stimmt (nur `sort -r` ist rot). Szenario für (a): Eine Umformung der Endeerkennung hält alle Fälle grün, während eine `## `-Zeile hinter `SPEC-034` fremde Tabellen in K zieht. | `AGENTS.md` §3.10 (Aufzählung und Ortsangabe zählen je Element), Register `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag` | `tools/harness/meldungskatalog-check-gegenprobe.sh` · `# --- Nr. 1: Quellen`; `tools/harness/meldungskatalog-check.sh` · `substr($0, 1, 1) == "#"` und `\| sort \| head -n 1` | ja — die Mutanten aus Kopien (`k-ende-nur-###`, `h-section-hash-ende-nur-###`, `trim-nur-space`, `len3`, `anzahl-nur-gueltige`, `unlesbar-sort-weg`) enden mit Exit 0 der Gegenprobe | Zusage-Element ohne Mutation, die nur es bricht |
| F-3 | MEDIUM | Randformen, die der Code still entscheidet und weder `SPEC-054` noch §6 nennen (verifiziert in Bäumen aus Kopien): ein Tabulator im Code-Text einer Zelle zerbricht die Ausgabe (`…benutzerhandbuch.md:Form: PGR-E2001: x`, Zeile und Code verschoben, Punkt 9); ein Handbuch mit CRLF ergibt die Sperre `Abschnitt ### Fehlercodes fehlt`; ein Symlink auf ein Verzeichnis `cmd` oder auf eine Go-Datei liefert Exit 0 trotz Literal `"PGR-` dahinter (Punkt 6); Leerraum und Tabulator um die Zelle sowie die Quelle nur aus `Form`-Zeilen sind entschieden, ohne dass ein Fall es hält (F-2). Das Register führt `BEO-REPO/harness-lesart-ohne-entscheidungsort` mit 2×; dies ist das dritte Mal. | `AGENTS.md` §3.12; `SPEC-054` Punkt 9 | `tools/harness/meldungskatalog-check.sh` · `-print0 2>/dev/null)` und `printf "%s\t%s\t%s\t%s\n"` | ja — Bäume aus Kopien von `fehler.go`, `SPEC-034` und Handbuch mit Tabulator im Code, CRLF und Symlink | Randform vom Code entschieden, in §6 und Vertrag nicht genannt |
| F-4 | LOW | Der Rot-Beleg des Falls `nr8-wurzel-fehlt` hängt am Aufrufort: Der Läufer aus §7 nennt `8-wurzel-ohne-sperre` rot (122 rot, 2 äquivalent), aber der Mutant (`cd` ohne Sperre) bleibt grün, sobald das Arbeitsverzeichnis keine Quellen trägt (gemessen aus einem Verzeichnis außerhalb des Repos; aus der Wurzel rot). Die Gegenprobe in einer Kopie ohne `spec/` und `docs/` (`SPEC-052` Punkt 13) fängt den Mutanten damit nicht. | `AGENTS.md` §3.10, §3.11 | `tools/harness/meldungskatalog-check-gegenprobe.sh` · `bau; sp nr8-wurzel-fehlt "$arbeit/nirgends" "$F fehlt"` | ja — `python3 -I` des Läufers aus §7, einmal aus der Wurzel, einmal aus einem leeren Verzeichnis | Rot-Beleg eines Mutanten hängt am Arbeitsverzeichnis |
| F-5 | LOW | Die Gegenprobe endet grün, wenn Fälle übersprungen sind: als root zehn Fälle (`nr8-unlesbar-*`, die beiden Reihenfolge-Fälle mit nicht lesbaren Pfaden), ohne die Locale `de_DE.UTF-8` `nr9-bytefolge`. Die Zeile in `harness/README.md` §Sensors („darunter die Reihenfolge der Sperren“) und die Sensor-Datei nennen das nicht; nur der Kopf der Gegenprobe tut es. | `AGENTS.md` §3.11 | `tools/harness/meldungskatalog-check-gegenprobe.sh` · `übersprungen — Rechte wirken nicht (root)`; `harness/README.md` · „darunter die Reihenfolge der Sperren“ | ja — `make meldungskatalog-check-gegenprobe` als root meldet Überspringen auf stderr bei Exit 0 | Zusage in Prosa weiter als die Prüfung unter root |
| F-6 | INFO | Der Vertrag sagt im Kopf, die Gegenprobe prüfe die Punkte 1 bis 10; `harness/README.md` nennt 1 bis 6 und 8 bis 10 (Punkt 7 hat keinen eigenen Fall, er fällt unter 4b), und die Sperre `internal/ fehlt` ist nicht erreichbar (offengelegt, Mutant `8-ohne-internal` äquivalent). Ablauf: Der Code-Commit `9b7ff671` steht vor der Entscheidung `db77c937`; der Implementer hat die acht Randformen zurückgegeben, die Entscheidungen sind umgesetzt. Zuständig: Architect (Wortlaut Kopf von `SPEC-054`). | `AGENTS.md` §3.11, §3.12 | `harness/vertraege/SPEC-054-meldungskatalog.md` · „Die Punkte 1 bis 10 prüft“ | nein — Wortlautfrage | Kopfsatz des Vertrags zählt Punkt 7 mit |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/meldungskatalog-check.sh` (Sperrgründe und Reihenfolge, Mengenregeln 4a bis 4d, Klassen 1 bis 6, Form, Duplikate, Sortierung, Literalregel, typisierte Konstante mit Punkt im Typnamen, Verzeichnis statt Datei, fehlende Wurzel, zwei Literale in einer Zeile, Zeile ohne schließendes `\|`, eingerückt = akzeptiertes Negativ) | geprüft, ohne Befund gegen `SPEC-054`; Findings F-2 bis F-4 betreffen Tests und nicht genannte Randformen |
| Mutationstabelle §7 (124 Mutanten, Läufer nachgefahren) | geprüft; 121 rot, 3 grün; die beiden Äquivalenz-Aussagen (`8-exit-awk-sperre`, `8-ohne-internal`) nennen Invariante und Ort und halten (awk-Ausgang wird auf Exit 2 abgebildet, `internal/` fehlt zuerst als `<C> fehlt`); der dritte grüne Mutant ist F-4 |
| Aufrufarten (`make -n/-q/-t/-i`, `MAKEFLAGS`, `SPEC-052` Punkt 21) | geprüft, ohne Befund: kein Stand, Fälle `nr10-*` und `nr11-make*`; Mutanten `10-makeflags-*` rot |
| Anschluss: `GATE_CHECKS`, `GATE_SCHNELL_AUSNAHME`, `LANGSAM`, `make gates-schnell-gegenprobe` | geprüft, ohne Befund; die Abweichung von der Vorgabe des Architect (`LANGSAM`) ist gemessen und bestätigt |
| Gegenprobe liest den Arbeitsbaum nicht (`SPEC-052` Punkt 13) | geprüft, ohne Befund: grün in einer Kopie nur mit `Makefile`, `harness/mk/` und beiden Skripten; gelesen werden Skript, `Makefile` und zwei Fragmente |
| Unverändert: `internal/hexagon/model/fehler.go`, `spec/`, `docs/user/`, Stop-Hook, Positivliste, Stop-Hook-Gegenprobe, `.claude/` | geprüft, ohne Befund (`git diff --stat` leer für diese Pfade) |
| `harness/sensors/meldungskatalog-check.md`, `harness/README.md` §Sensors (zwei Zeilen), `harness/mk/meldungskatalog-check.mk` | geprüft, ohne Befund außer F-5; kein Werkzeugname in `spec/architecture.md` |
| Messung (`git diff --numstat 8d14a918..21e60054 -- tools harness`) | geprüft, ohne Befund: 560 nachgerechnet (142 + 286 + 72 + 18 + 36 + 2 + 2 + 2), eigene 524; Marke 900 und Meldeschwelle 750 nicht erreicht |
| Gate am Stand des Starts (C 30, K 34, H 30) | geprüft, ohne Befund, grün ohne Änderung an `fehler.go` und `SPEC-034` (`AGENTS.md` §3.6: keine Stufung, keine Lockerung) |
| Plan §1, §3, §4, §6 gegen den Diff (`AGENTS.md` §3.9) | geprüft, ohne Befund außer F-1; §6 nennt Aufrufarten und Reihenfolgen; Adressen (`slice-v1-abschluss-container`) unverändert |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Abschnitt des Plans beim Umschreiben von §7 verloren · Zusage-Element ohne Mutation, die nur es bricht · Randform vom Code entschieden, in §6 und Vertrag nicht genannt · Rot-Beleg eines Mutanten hängt am Arbeitsverzeichnis · Zusage in Prosa weiter als die Prüfung unter root

## Verdikt

**Merge-blockierend:** ja — F-1 bis F-3 (MEDIUM) blockieren. Zuständigkeit: F-1 und F-2 Implementer; F-3 geht an den Architect (Randformen im Vertrag entscheiden, danach Fälle durch den Implementer, `AGENTS.md` §3.12); F-4 und F-5 Implementer; F-6 Architect.

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt: F-1); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
