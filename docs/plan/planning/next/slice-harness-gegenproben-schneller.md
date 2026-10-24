# Slice slice-harness-gegenproben-schneller: Die Gegenproben laufen deutlich schneller und räumen ihren Arbeitsordner je Mutation auf

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Entscheidung des Planners, 2026-10-20). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist, und keine Anforderung der Welle `welle-v1-nachtrag` hängt an ihm; würde er in sie aufgenommen, hinge deren Closure-Trigger an einem Werkzeug-Slice, den kein Slice der Welle voraussetzt. Angelegt bei der Closure von `slice-v1-abschluss-binary-namen` (Lerneintrag, benannte Lücke; Beobachtung des Verifiers F-V2 und des Reviews F-2) und nach dem Hinweis des Nutzers vom 2026-10-20, die Ursache der langen Läufe seien verschachtelte Schleifen mit vielen `sed`- und `awk`-Aufrufen; das ist die Vermutung des Nutzers, belegt ist sie nicht (§6). Der Slice liegt in `open/`: Er ist nicht priorisiert, und was der Architect nach dem Lesen der Skripte entscheidet, liegt nicht vor (§4 *Start*).

**Bezug:** kein Lastenheft-Bezug — der Slice ändert Werkzeuge des Harness, keinen Anspruch des Produkts. Bindung an das Register: `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` (die Größe „freier Platz des Datenträgers“, Beleg `evidence/slice-v1-abschluss-binary-namen.md`) und `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` (Laufzeit der Gegenproben). Entscheidungen: [ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md) (die Gegenproben stehen nicht in `GATE_CHECKS`; hier unverändert, §1) und [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) (Obergrenze der Zeit für jede Gegenprobe, Aufräumen je Mutation; Accepted, 2026-10-21).

**Berührte Spec-Stellen:** `SPEC-067` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-067-release.md); `make release-gegenprobe`) · `SPEC-055` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-055-container-und-binaries.md); `make test-container-gegenprobe`) · `SPEC-060` ([`harness/vertraege/`](../../../../harness/vertraege/SPEC-060-gegenproben.md); Regeln für Gegenproben). Die Zusagen zu Zeit und Arbeitsordner stehen in `SPEC-060` Punkt 16 (Entscheidung des Nutzers vom 2026-10-21; [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md), Accepted) und mit Verweis in `SPEC-067` Punkt 14 und `SPEC-055` Punkt 12.

**Verantwortlich:** pt9912
**Autor:** Planner. **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Gegenproben `tools/dev/harness/release-gegenprobe.sh` (2045 Zeilen, 14 Zeilen mit `sed`, 10 mit `awk`, 59 Schleifen) und `tools/dev/harness/container-gegenprobe.sh` (790 Zeilen) laufen deutlich schneller, bei gleichem Ergebnis, und der Arbeitsordner einer Mutation wird nach ihr entfernt statt am Ende des Laufs. Bisher wächst der Arbeitsordner von `make release-gegenprobe` um rund 40 MB je Mutation, über 5 GB bei 535 Mutationen, und wird erst am Ende freigegeben; ein voller Lauf braucht rund 6 GB und scheitert auf einem Host mit weniger freiem Platz als scheinbarer Befund (`slice-v1-abschluss-binary-namen` §7, Verifikation F-V1 und F-V2, Review F-2: sechs Mutationen `m-p16-*` rot im Lauf, allein gefahren grün). **Die Zahl (Entscheidung des Nutzers vom 2026-10-21: „die gegenproben werden zeitlich begrenzt. höchstens 10 minuten“; [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)):** Ein voller Lauf jeder Gegenprobe dauert höchstens **600 s**, Wanduhr auf ruhigem Host (`SPEC-060` Punkt 16 (1)); die Gegenprobe misst sich selbst und endet bei Überschreitung mit Exit 1. **Dieser Slice nimmt `release-gegenprobe` und `test-container-gegenprobe` auf**; die übrigen 17 Gegenproben der Container-Liste und `container-start-gegenprobe` stehen im Zustand `nicht gemessen` in der Tabelle von `SPEC-060` Punkt 16, mit der Adresse eines Folge-Slice (`slice-harness-gegenproben-grenze-rest`, angelegt vom Planner am 2026-10-21, §6 Rückgabe 4). Zum Vergleich der Stand `8dbe6a14` (§3, Messung): `release-gegenprobe` rund 25 bis 30 Minuten bei vier gleichzeitigen Läufen auf ruhigem Host und Stunden unter Last, `test-container-gegenprobe` 22 s und rund 1 s je Mutation. **Eine Mutation wird nicht gestrichen und nicht zusammengelegt** (§1 *Ausdrücklich NICHT*; die Option, je Mutation nur eine Auswahl der Fälle zu fahren, gehört dem Nutzer, §6 Rückgabe 3).

**Bindung aus `slice-v1-abschluss-homebrew` (§3.13):** Dessen Closure weist diesem Slice zu, dass der volle Lauf von `make release-gegenprobe` am Endstand die rund 95 Mutationen `m-p19-*` fährt, die jener Slice nicht wiederholt hat; sie sind Teil des Laufs, den die DoD unten verlangt, und der Beleg der Ergebnisgleichheit (Deklarationszeilen) deckt sie. Keine eigene Lieferung, kein vierter Punkt (§8).

**Sendung aus dem Carveout `CO-001`** (Geber: [`CO-001`](../../carveouts/CO-001-gegenproben-am-getaggten-commit.md), Aktiv seit 2026-10-21; `AGENTS.md` §3.13): `CO-001` nennt diesen Slice als Folge-Slice für `release-gegenprobe` und `test-container-gegenprobe`. Dieser Slice nimmt an: Mit der Aufnahme beider in der Tabelle von `SPEC-060` Punkt 16 ist der Teil des Auflösungs-Triggers von `CO-001` erfüllt, der ihn betrifft; den Carveout selbst löst er nicht auf (er bleibt Aktiv, bis `slice-harness-gegenproben-grenze-rest` und die von ihm angelegten Slices die übrigen aufgenommen haben). Die DoD trägt die Zeile; `grep` nach `CO-001` findet die Sendung in §1, DoD und §8. **Sendung aus dem Carveout `CO-003`** (Geber: [`CO-003`](../../carveouts/CO-003-gegenproben-am-wellenende.md), Aktiv seit 2026-10-24; Pflicht am Wellenende): `CO-003` hat denselben Auflösungs-Trigger wie `CO-001` und nennt diesen Slice mit demselben Gegenstand; neue Arbeit entsteht daraus nicht. Die Closure nennt den Stand von `CO-003` neben dem von `CO-001`. **Ausgeschlossen für diesen Slice bleibt**, was `CO-001` und `CO-003` ihm nicht zuweisen: die übrigen Gegenproben (§1 *Ausdrücklich NICHT*, Adresse `slice-harness-gegenproben-grenze-rest`).

**Die Ursache ist gemessen (Architect, 2026-10-20; Zahlen in §3).** Die Vermutung des Nutzers (`sed`, `awk`) trifft nicht zu: beide sind 2,4 Prozent der rund 50 000 Prozess-Starts einer Mutation der Gruppe `f`. Die Zeit liegt in Prozess-Starts aus einer großen `bash` (Systemzeit 21 s von 33 s) und darin, dass jede Mutation die ganze Gruppe (bis 112 Fälle) fährt, deren Fälle `release.sh` mehrfach unter Attrappen starten (`bash`, `mkdir`, `rm`, `jq`, `sha256sum`, `tr`, `cat`, `date`, `tar`). Kopien des Arbeitsbaums und `git`-Repositories sind es nicht: Die Welt ist ein Repository mit drei leeren Commits, eine Kopie je Fall kostet unter 1 MB.

**Übergang für die Rollen (Empfehlung des Architect, bis der Slice geliefert ist):** Kein Implementer fährt einen vollen Lauf von `make release-gegenprobe` oder `make test-container-gegenprobe`; er fährt den engsten Fall mit `GP_MUT=<Muster>` (`AGENTS.md` §6 Punkt 6), der Verifier den vollen Lauf einmal am Endstand. Die Anweisung steht seit dem 2026-10-21 in `AGENTS.md` §6 Punkt 6 und in `.claude/commands/implement-slice.md` Schritt 15 (Vorschlag des Architect, vom Nutzer mit „ich gebe alles frei“ freigegeben, Planner); mit der Lieferung dieses Slice nimmt der Implementer sie aus beiden Stellen, weil der volle Lauf dann unter 600 s endet. Die Zeile je Gegenprobe in `harness/README.md` §Sensors bleibt bis zur Lieferung unverändert (§3.11: sie sagt nur zu, was ein Fall prüft).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Änderung dessen, was die Gegenproben prüfen: Ergebnis, Fälle, Mutationen und die Mengen der roten Fälle je Mutation bleiben gleich (Klasse 3, anderer Vorgang). Grund: Der Slice ändert, wie schnell geprüft wird, nicht, was; eine gestrichene oder zusammengelegte Mutation senkte die Prüfung, nicht die Laufzeit (`AGENTS.md` §3.6, Gates nicht ohne ADR lockern).
- Ein Nachtlauf der Gegenproben in der CI — ein anderer Vorgang (Klasse 3). Grund: Er wäre eine neue Außenwirkung des Harness und braucht eine ADR; die Slices von `welle-ci-nachtlaeufe` und `slice-harness-ci-gegenproben` (`open/`) tragen ihn, dieser Slice nimmt ihn nicht an.
- Die übrigen Gegenproben (`make hook-gegenprobe`, `make lint-gegenprobe`, `make workflows-check-gegenprobe` und die weiteren der Liste `GEGENPROBEN_CONTAINER`) — Bestand bleibt bewusst stehen (Klasse 2). Grund: Gemessen und gelesen sind die zwei größten; ob eine andere dieselbe Form trägt, zeigt das Lesen des Architect, und was er dort findet, geht in einen Folge-Slice, nicht in diesen.
- Die Aufnahme einer Gegenprobe in `GATE_CHECKS` oder `make gates` — Klasse 3 ([ADR-0048](../../adr/0048-gegenproben-verlassen-die-gate-kette.md)). Grund: Schneller macht sie nicht zum Gate der Kette.
- Produkt-Code (`cmd/`, `internal/`) und `Dockerfile` — Schicht-Abgrenzung (Klasse 4): Der Slice berührt allein Skripte unter `tools/harness/`, ihre Gegenproben und die Zeile in `harness/README.md`.
- Der Gegenstand `slice-harness-gegenprobe-baumkopie` (Kopie des Arbeitsbaums nach `git ls-files` durch einen gemeinsamen Helfer, `open/`) — Klasse 1 des anderen Slice, nicht dieses. Grund: Er trägt die Kopie des Baums für `a-check-negativ` und weitere; berührt der Architect hier dieselbe Stelle, nennt er das und gibt es an den Planner.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] `make release-gegenprobe` läuft höchstens 600 s (`SPEC-060` Punkt 16 (1)) mit unverändertem Ergebnis: die unveränderte Kopie grün, je Mutation genau die benannten Fälle rot. Beleg der Gleichheit: der Unterschied der Deklarationszeilen (Zeilen `mutr`, `mf`, `mh`, `mutv`, `mut1`, `mut13g`, die Fall-Listen `FAELLE_*`) zwischen Basis und Endstand ist leer (`git diff <Basis>..HEAD -- tools/dev/harness/release-gegenprobe.sh`, gefiltert auf diese Zeilen), und ein voller Lauf am Endstand ist grün; ein Vorher-Lauf ist nicht nötig, weil die Gegenprobe je Mutation die tatsächlichen mit den deklarierten roten Fällen vergleicht. Die Stichprobe (§3) sinkt von rund 66 s auf höchstens 22 s.
- [ ] Die Fälle und Mutationen von `SPEC-060` Punkt 16 (9) stehen in beiden Gegenproben, je Element eine Mutation, die nur ihren Fall färbt (`AGENTS.md` §3.10).
- [ ] `make test-container-gegenprobe` ebenso: höchstens 600 s, Ergebnis unverändert (Beleg wie oben).
- [ ] Die Tabelle in `SPEC-060` Punkt 16 führt beide Gegenproben als `aufgenommen`, und die Closure nennt, dass der Teil des Auflösungs-Triggers von `CO-001`, der diese beiden betrifft, erfüllt ist (Sendung `CO-001`; der Carveout bleibt Aktiv, bis die übrigen Zeilen der Tabelle aufgenommen sind, Nehmer `slice-harness-gegenproben-grenze-rest`; ebenso der Stand von `CO-003`, derselbe Trigger); der volle Lauf von `release-gegenprobe` am Endstand fährt die rund 95 Mutationen `m-p19-*` mit, die `slice-v1-abschluss-homebrew` nicht wiederholt hat (Bindung oben).
- [ ] Der Arbeitsordner je Mutation und je Fall wird nach ihrem Ende entfernt (grün wie rot; bei `SIGINT` und `SIGTERM` der ganze Arbeitsordner, wie bisher, `SPEC-060` Punkt 16 (5)); die Spitze des Platzes eines vollen Laufs von `make release-gegenprobe` liegt unter 100 MB (gemessen mit `du` im Abstand von einer Sekunde, im Bericht), und der volle Lauf in einem Zug ist am Endstand belegt (`slice-v1-abschluss-binary-namen`, F-V2).
- [ ] `make gates` grün; die Gegenproben, die der Slice ändert, am Endstand in einem Zug (`AGENTS.md` §6); die übrigen Gegenproben der Listen bleiben unberührt (§1).
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

Der Architect hat diese Tabelle nach dem Lesen der Skripte gefüllt (2026-10-20). **Messung** am Stand `8dbe6a14`, Gegenproben-Container, 20 Kerne, Last 4 bis 5, je ein Lauf mit `GP_MUT`:

- Eine Mutation der Gruppe `f` (`m-p16-pin-ohne-digest`, 112 Fälle): 30 s Wanduhr (12 s user, 21 s sys), 50 389 Prozess-Starts (`strace -f -e trace=execve`): `bash` 5244, `mkdir` 2024, `rm` 1594, `jq` 1452, `sha256sum` 1438, `tr` 1402, `cat` 1256, `grep` 1098, `date` 915, `tar` 894, `sed` 838, `cut` 836, `head` 724, `docker` (Attrappe) 609, `mktemp` 579, `awk` 342, `git` 239.
- Gruppe `b` (`m-p2-trockenlauf-image-ruft-skopeo`) rund 30 s, Gruppe `sh` (`m-p2-trockenlauf-ruft-skopeo`) 6 s, die unveränderte Kopie 104 s.
- Mutationen: 122 `mf` (Gruppe `f`), 65 `b`, 53 `sh`, 7 `p13`, 9 `gp13`; dazu die Mutationen des Textes (`p1`, rund 120, billig). Gesamt rund 6300 s Arbeit, bei `GP_PARALLEL=4` 25 bis 30 Minuten auf ruhigem Host.
- Platz: rund 32 MB nach 110 Fällen einer Mutation; der Ordner `t.<fall>` jedes Falls bleibt bis zum Ende des Laufs.
- `test-container-gegenprobe`: die unveränderte Kopie 22 s, zwei Mutationen der Gruppe `p7-cp` zusammen 2 s; sie läuft nacheinander mit einer Auswahl der Fälle (`nur`) je Mutation und räumt die Bäume `baum-<n>` erst am Ende.
- **Stichprobe** für Vorher und Nachher: `GP_PARALLEL=1 GP_MUT='^(m-p16-pin-ohne-digest|m-p2-trockenlauf-image-ruft-skopeo|m-p2-trockenlauf-ruft-skopeo)$' make release-gegenprobe` (ruhiger Host); heute 66 s, Ziel höchstens 22 s, Prozess-Starts der ersten Mutation von 50 000 auf höchstens 17 000.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/dev/harness/release-gegenprobe.sh` | update | Prozess-Starts je Fall senken (Ansatzpunkte, keine Vorschrift: `$(cat …)` durch Lesen in der `bash`; `rm -rf` und `mkdir -p` je Fall ohne Unterprozess, wo ein frischer Ordner schon leer ist; das Layout von `gp-layout` und die Antworten der Attrappe von `docker` je Variante einmal je Welt; ein `cp -a` je Fall nur für das, was der Fall ändert); Ordner je Fall und je Mutation entfernen; Zeitbudget und Ausgabe nach `SPEC-060` Punkt 16; Fälle und Mutationen dazu |
| `tools/dev/harness/container-gegenprobe.sh` | update | Bäume `baum-<n>` je Mutation entfernen; Zeitbudget und Ausgabe nach `SPEC-060` Punkt 16; Fälle und Mutationen dazu; Tempo nur, wenn der volle Lauf 300 s übersteigt (Messung im Bericht; geschätzt liegt er bei 4 bis 5 Minuten) |
| `harness/README.md` (§Sensors, Zeilen der beiden Gegenproben) | update, erst mit der Lieferung | nennt Obergrenze (600 s), Aufräumen je Mutation und `GP_BUDGET`; vorher nicht (§3.11) |
| `AGENTS.md` §6 Punkt 6, `.claude/commands/implement-slice.md` Schritt 15 | update, erst mit der Lieferung | nimmt den Übergangssatz zu `GP_MUT` und dem vollen Lauf des Verifiers heraus (Entfernen, kein Liefer-Punkt; `AGENTS.md` §3.9) |
| `harness/vertraege/` | durch den Architect, geschrieben (2026-10-20) | `SPEC-060` Punkt 16 (Zeit, Platz, Randformen, Fälle), Verweise in `SPEC-067` Punkt 14 und `SPEC-055` Punkt 12; [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) |

**Schätzung des Architect in der Metrik von §4** (Summe aus Einfügungen und Löschungen, Gegenprobe-Code ohne Faktor, weil der Code selbst Gegenprobe-Code ist; Vertragstext und ADR des Architect liegen vor der Basis und zählen nicht): `release-gegenprobe.sh` 200 bis 320 (Umbau von `lauf`, `kopie`, `attrappen`, Aufräumen, Budget 25, Fälle und Mutationen zu `SPEC-060` Punkt 16 (9) rund 120), `container-gegenprobe.sh` 100 bis 160, `harness/README.md` 10; zusammen 310 bis 490, mit Streuung durch Umstellungen bis 700. Unter der Marke von 1000.

Nicht in dieser Liste: Produkt-Code, `Dockerfile`, andere Gegenproben, CI (§1).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next`, erledigt am 2026-10-21): [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) ist angenommen (Entscheidung des Nutzers vom 2026-10-21: „die gegenproben werden zeitlich begrenzt. höchstens 10 minuten“), und der Slice ist eingeplant (Planner) mit `Verantwortlich: pt9912`. Reihenfolge: **hinter `slice-v1-abschluss-release-v0-1-1`**, der zuerst dran ist (Empfehlung des Planners; der Nutzer will zuerst den Release, Entscheidung des Nutzers vom 2026-10-21: „ich gebe alles frei, damit wir jetzt endlich diesen brew release machen können“); `slice-v1-abschluss-release-v0-1-1` setzt diesen Slice nicht voraus und er jenen nicht, die Punkte 2 und 3 der Freigabe laufen bis zu dessen Lieferung über `CO-001`. Eine Stellung zu `slice-erster-release-homebrew-nachweis` ist nicht entschieden, keiner der beiden setzt den anderen voraus. Das Lesen der zwei Skripte, Zielzahl, Randformen und Schnitt sind geschehen (Architect, 2026-10-20; §3, §6, §8); **Start von `next` nach `in-progress`:** `slice-v1-abschluss-release-v0-1-1` liegt nicht mehr in `in-progress/` (WIP-Limit 1 je Rolleninhaber), und vor dem ersten Code-Commit trägt der Architect `<Basis>` ein.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe aus Einfügungen und Löschungen von `git diff --numstat <Basis>..<letzter Commit des Slice> -- tools harness/README.md harness/vertraege` (Tests und Gegenproben eingeschlossen, `AGENTS.md` §3.14) **über 1000** liegt; `<Basis>` ist der Hash von `main` beim Übergang nach `in-progress/`, den der Architect in §3 einträgt. Die Schätzung des Architect in dieser Metrik steht in §3 (Faktor von Tabellen- und Gegenprobe-Code zum Produktcode rund zwei und mehr; hier ist der Code selbst Gegenprobe-Code, deshalb ohne Faktor, aber mit Zeilen, die ein Umbau verschiebt). Die Schätzung in §3 liegt bei 310 bis 490, höchstens 700; die Marke ist vor dem Code nicht absehbar überschritten, ein Schnitt ist nicht nötig (§8).
- `in-progress` → `next`: Die Stichprobe (§3) sinkt nach den Ansatzpunkten aus §3 nicht mindestens um den Faktor 2 (von 66 s auf höchstens 33 s), oder ein voller Lauf von `release-gegenprobe` überschreitet 600 s auch nach einem Lauf mit `GP_PARALLEL=8` auf ruhigem Host. Dann geht die Messung an den Architect und an den Nutzer (Wahl nach §6 Rückgabe 3), der Implementer grübelt nicht weiter. (Das Lesen hat gezeigt, dass die Laufzeit an den Skripten hängt, nicht an `docker` oder dem Datenträger.)
- `in-progress` → `open`: Der Nutzer entscheidet für eine der Lockerungen aus §6 Rückgabe 3; dann ändert sich `SPEC-060` Punkt 16 durch den Architect, und der Slice wird neu geschnitten.

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

**Randformen** (`AGENTS.md` §3.12; vom Architect vor dem ersten Code-Commit zu entscheiden und im Vertrag festzuhalten; der Planner nennt sie, entscheidet sie nicht):

| Randform | Frage | Entscheidung |
|---|---|---|
| (a) Arbeitsordner einer roten Mutation | Bleibt er für die Diagnose stehen (Pfad im Log) oder wird er auch dann entfernt? | **entschieden: entfernt auch bei rot**; ein bleibender Ordner gibt es nicht, der Arbeitsordner war schon bisher mit dem Lauf weg; `SPEC-060` Punkt 16 (5) |
| (b) Abbruch mitten in einer Mutation | `SIGINT` und `SIGTERM`: der Ordner der laufenden Mutation wird entfernt; die Fälle `p13-*` tragen das heute für den ganzen Lauf | **entschieden: der ganze Arbeitsordner, die laufenden Mutationen eingeschlossen, nach dem Beenden der Nachkommen (Bestand, Fälle `p13-*`); keine Zeile `Zeit` bei Abbruch**; `SPEC-060` Punkt 16 (2), (4), (5) |
| (c) Parallele Läufe (`GP_PARALLEL`) | Mehrere Mutationen zugleich: Platzbedarf und Entfernen je Ordner, kein gemeinsames Verzeichnis | **entschieden: je Mutation ein eigener Ordner, je Job sein Entfernen; Vorgabe von `GP_PARALLEL` darf von 4 auf höchstens 8 steigen, wenn zwei volle Läufe in Folge grün sind**; `test-container-gegenprobe` läuft nacheinander und bleibt es; `SPEC-060` Punkt 16 (5), (7) |
| (d) Aufrufarten außerhalb des Normalfalls | `GP_MUT`, `GP_ENTDECKE`, `GP_FALL` und die Aufrufe unter `make` mit `-n`, `-q`, `-t`, `-i` und `MAKEFLAGS` (`AGENTS.md` §3.12, Aufrufarten des Treibers): unverändert, mit je einem Fall | **entschieden: die vier Variablen unverändert; unter `make` kein Fall (kein Hook und kein Stand liest das Ergebnis, Bestand `SPEC-067` Punkt 14); neu `GP_BUDGET` (ganze Zahl, `0` aus, ungültig Exit 2 vor jedem Fall), Vorgabe der Obergrenze nur ohne `GP_MUT`; Reihenfolge der Gründe: ungültig, Signal, rote Fälle, Zeit**; `SPEC-060` Punkt 16 (3), (4), (7) |
| (e) Plattform | Gegenproben laufen im Gegenproben-Container unter Linux; macOS und Windows nicht zugesagt | **entschieden: Linux gemessen und zugesagt; macOS und Windows Grenze aus dem Kopf von `SPEC-060`; fehlt `/proc/loadavg`, steht `-` als Last**; `SPEC-060` Punkt 16 (2), (6) |
| (f) Umgebung der Fahrt | freier Platz, Last, root, Locale und Zeitzone (`BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts`): der Slice fährt die Gegenprobe vor der Übergabe einmal unter fremdem Benutzer, fremder Locale und unter Last | **entschieden: die Obergrenze gilt auf ruhigem Host (Last unter einem Fünftel der Kernzahl); eine Überschreitung unter Last ist rot ohne Befund am Code, die Zeile nennt die Last, wer sie liest, fährt auf ruhigem Host neu (keine Prüfung, die die Last erkennt und auslässt); die Fahrt unter fremdem Benutzer, fremder Locale und unter Last nennt die Zeit ohne Geltung der Grenze; Platz: Spitze unter 100 MB**; `SPEC-060` Punkt 16 (1), (3), (8) |

**Rückgaben an den Architect** (nach dem Lesen, vor dem Code):

1. **Beantwortet (Architect, 2026-10-20):** Ort ist `SPEC-060` Punkt 16 (gemeinsam, weil die Zusage für beide Gegenproben dieselbe ist), mit Verweisen in `SPEC-067` Punkt 14 und `SPEC-055` Punkt 12; Randformen (a) bis (f) in der Tabelle.
2. **Beantwortet:** Messgröße ist die Wanduhr des vollen Laufs, gemessen vom Skript selbst, dazu die Stichprobe (§3) für Vorher und Nachher; Zahl 600 s für jede Gegenprobe, Stichprobe 66 s auf höchstens 22 s. Eine ADR ist nötig, weil die Zusage Rot bei Überschreitung trägt und das Streichen oder Lockern eine Lockerung wäre (`AGENTS.md` §3.6): [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md), Accepted (2026-10-21).
3. **Schnitt: keiner** (§8). **Offen beim Nutzer (Empfehlung des Architect):** Es wird keine Mutation gestrichen oder zusammengelegt (je Zusage und Element eine, `AGENTS.md` §3.10). Reichen die Kosten je Mutation für die Obergrenze nicht (§4, Rückführung), entscheidet der Nutzer zwischen (i) einer höheren Obergrenze und (ii) je Mutation nur den Fällen, die sie berühren kann (Auswahl nach den benannten roten Fällen und ihren Nachbarn; Faktor 10 und mehr, aber „genau die benannten Fälle rot“ gilt dann nur unter der Auswahl); beides ist eine Lockerung und braucht eine ADR.
4. **An den Planner — angelegt am 2026-10-21 (`slice-harness-gegenproben-grenze-rest`, `open/`; er nimmt in §1 und DoD an und trägt die Sendung):** Die übrigen 17 Gegenproben der Container-Liste und `container-start-gegenprobe` (Host) sind nicht gemessen und stehen in der Tabelle von `SPEC-060` Punkt 16 mit dem Nehmer `slice-harness-gegenproben-grenze-rest` (Vorschlag des Architect; der Slice ist anzulegen, und der Nehmer trägt die Sendung in §1 und DoD ein, bevor die Kennung gilt, §3.13; ohne vorhandenen Nehmer ist der Tabelleneintrag ein Fehler des Vertrags). Der Nehmer misst jede, nimmt die unter 600 s auf und schneidet die übrigen nach Gegenprobe. Größte Kandidaten nach Zeilen: `lint-gegenprobe` (1206), `stop-hook-gegenprobe` (1217), `container-start-gegenprobe` (886, nur Host), `gegenproben-nach-einlieferung-gegenprobe` (839), `ci-gegenprobe` (638). **Die Baumkopie** (`slice-harness-gegenprobe-baumkopie`) wird hier nicht berührt: Die Welt der Gegenprobe ist ein Repository mit drei leeren Commits, keine Kopie des Arbeitsbaums.

**Risiken:**

- Die Ansatzpunkte aus §3 senken die Prozess-Starts weniger als geschätzt (die Ursache ist gemessen, die Wirkung der Mittel nicht) — **Ausgang:** offen bis Closure (entfallen: Stichprobe und voller Lauf am Endstand innerhalb der Zahlen; eingetreten: Rückführung nach §4 und Entscheidung des Nutzers nach Rückgabe 3).
- Das Ergebnis ändert sich unbemerkt: Ein schnellerer Fall färbt eine andere Menge roter Fälle oder bleibt grün, wo er rot sein müsste — **Ausgang:** offen bis Closure (Gegenmittel: Vorher-Nachher-Vergleich der roten Mengen aller Mutationen mit `GP_ENTDECKE=1`, einmal am Endstand in einem Zug; entfallen, wenn der Vergleich leer ist).
- Zeitmessungen schwanken mit der Last des Hosts; ein voller Lauf kann unter Last die Obergrenze reißen und rot enden (`SPEC-060` Punkt 16 (3), akzeptiert) — **Ausgang:** offen bis Closure (`BEO-REPO/lastabhaengige-zaehlung`: weiter offen, wenn eine Zahl der DoD an der Last hängt).
- `GP_PARALLEL` über 4 lässt Fälle mit Wartezeiten und `pgrep` ausfallen (`BEO-REPO/test-faellt-unter-cpu-last-aus`) — **Ausgang:** offen bis Closure (entfallen, wenn die Vorgabe bei 4 bleibt oder zwei volle Läufe in Folge grün sind).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen. Treffer: `gegenprobe-haengt-an-umgebung-des-hosts` (9×, Größe „freier Platz“ neu; der Ausgang liegt beim Nutzer, dieser Slice ist ein Teil der Antwort auf diese Größe, nicht die Zuweisung), `gate-laufzeit-waechst-mit-gegenprobe` (verkörpert; die Laufzeit der Gegenproben außerhalb von `make gates` ist hier der Gegenstand), `lastabhaengige-zaehlung` (Zeitmessung, §6), `gegenprobe-temp-praefix-geteilt` (Arbeitsordner, Präfix; der Architect liest den Eintrag vor dem Entfernen je Mutation). Der Stand der Zähler ist beim Start zu lesen; erreicht einer mit diesem Slice drei, braucht er einen Folge-Slice.

**Zählung** (Planner, 2026-10-20, `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 3** (1: `release-gegenprobe.sh` schneller; 2: `container-gegenprobe.sh` schneller; 3: Arbeitsordner je Mutation entfernt, in beiden Skripten). **Schichten: eine** — Harness-Werkzeug (Skripte, Zeilen in `harness/README.md`); Verträge liegen vor dem Code beim Architect und zählen nicht als Liefer-Punkt. Es sind noch keine Sendungen anderer Slices eingetragen. Größe: Marke in §4.

**Nachzählung** (Architect, 2026-10-20, nach dem Lesen): **Liefer-Punkte 3, Schichten eine, kein Schnitt.** Das Zeitbudget (Messung, Ausgabe, Rot) und der Arbeitsordner sind je Skript ein Liefer-Punkt zusammen mit dem Tempo des Skripts, nicht ein vierter: Beide Skripte tragen sie inline (kein gemeinsamer Helfer; ein solcher wäre `slice-harness-gegenprobe-baumkopie` Klasse 1 nahe und ein dritter Ort). Der Punkt 3 des Planners (Arbeitsordner) geht dabei in 1 und 2 auf; der dritte Liefer-Punkt ist die Fall-Menge zu `SPEC-060` Punkt 16 (9) in beiden Gegenproben. Die Skripte sind in einer Review-Sitzung prüfbar, wenn die Umbauten je Skript klein bleiben (Schätzung §3: 310 bis 490); wächst ein Skript über die Marke, schneidet der Planner nach Skript. Die Zusagen aus `SPEC-060` Punkt 16 stehen vor dem Code im Vertrag; Sendungen an andere Slices trägt dieser Commit nicht ein (die übrigen Gegenproben: §6 Rückgabe 4, Adresse beim Planner).

**Nachzählung der Bindung aus `slice-v1-abschluss-homebrew`** (Architect, 2026-10-21, §3.13): Die rund 95 nicht wiederholten Mutationen `m-p19-*` sind keine neue Lieferung, sondern Teil des vollen Laufs von `release-gegenprobe` am Endstand, den DoD 1 ohnehin verlangt (es gibt keinen Vorher-Lauf). Liefer-Punkte bleiben 3, Schichten eine; die Marke von §4 ist unberührt (sie zählt Diff-Zeilen, der Lauf ändert keine). Die Messung in §3 stammt vom Stand vor Homebrew und zählt sie nicht: Die Gruppe `h` (Fälle `p19-*`, rund 95 Mutationen `mh`) kommt zu den 122 + 65 + 53 hinzu und ist nicht gemessen; der Implementer misst einen Lauf mit `GP_MUT='^m-p19-'` am Anfang (Stichprobe der Gruppe `h`, je Mutation) und meldet die Zahl dem Architect, bevor er umbaut, weil sie die Schätzung der Arbeit (rund 6300 s) erhöht.

**Nachzählung des Übergangssatzes** (Planner, 2026-10-21, `AGENTS.md` §3.13): Das Herausnehmen des Satzes aus `AGENTS.md` §6 und `implement-slice.md` ist Teil des Abschlusses der Liefer-Punkte 1 und 2 (Harness-Artefakte, dieselbe Schicht wie die Skripte), kein vierter Punkt; Liefer-Punkte bleiben 3, Schichten eine, die Schätzung wächst um unter 10 Zeilen.

**Nachzählung der Sendung aus `CO-001`** (Planner, 2026-10-21, `AGENTS.md` §3.13): Die Zeile in DoD 4 (Stand des Auflösungs-Triggers in der Closure) ist ein Closure-Text, kein Liefer-Punkt, und berührt keine Schicht. **Liefer-Punkte bleiben 3, Schichten eine;** die Marke und die Schätzung aus §3 und §4 sind unberührt (die Zeile ändert keine Diff-Zeilen von `tools`, `harness/README.md` oder `harness/vertraege`).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Verträge und Konventionen des Harness gehen den Skripten voraus).

**Nachgezählt** (`AGENTS.md` §3.13, Planner, 2026-10-22, im Commit, der die Ortsangabe aus `slice-harness-gegenproben-umzug-dev` einträgt): Die Gegenproben liegen seit jenem Slice unter `tools/dev/<bereich>/`, die Pfade der Gegenproben in diesem Plan sind nachgezogen (`tools/dev/` ist wie `tools/harness/` Harness-Werkzeug, keine neue Schicht). Ortsangabe der Gegenprobe geändert, Liefer-Punkte (3) und Schichten (1) unverändert.

**Nachzählung der Sendung aus `CO-003`** (Planner, 2026-10-24, `AGENTS.md` §3.13, im Commit, der sie einträgt): Die Sendung steht in §1 und DoD 4 und ist wie die aus `CO-001` ein Closure-Text ohne Diff-Zeilen. **Liefer-Punkte bleiben 3, Schichten eine;** Marke und Schätzung unberührt.
