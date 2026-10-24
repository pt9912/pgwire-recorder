# Slice slice-homebrew-nachweis-macos-arm64: Bericht der Installation aus dem Tap auf macOS (arm64)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag (Begründung in §1).

**Bezug:** [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew)

**Berührte Spec-Stellen:** `LH-FA-19.a` (geändert: Absatz *Bericht je Kombination*, Entscheidung des Nutzers vom 2026-10-24, im Commit `4279f2b3` geliefert, nicht Teil dieses Slice) · `SPEC-067` (Punkt 20 (6), geändert: Bericht als Stufe, im selben Commit)

**Verantwortlich:** pt9912 (Zweig impl-mac1)
**Autor:** pt9912. **Datum:** 2026-10-24.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Für die Kombination **macOS auf `arm64`** steht die Stufe *Bericht* nach `LH-FA-19.a` und `SPEC-067` Punkt 20 (6) in Handbuch, `README.md` und `docs/maintainer/releasing.md`: als **berichtet** (Version `0.2.0`, Homebrew 7.0.10), nie als gefahren oder belegt. Das Zitat der Nutzerangabe mit Datum steht wortgetreu in diesem §1 (Absatz unten); eine Beleg-Datei unter `docs/reviews/` und eine Summe entstehen nicht. Der Replay-Teil von Abnahmeszenario 11 (Schritt (3), `brew test`, Ausgabe der Installation, Exit-Codes) bleibt für die Kombination ohne Beleg; macOS auf `amd64` und Linux auf `arm64` bleiben Grenze.

**Zitat und Merkmale des Berichts (Nutzer, 2026-10-24, wortgetreu).** Der Nutzer hat in der Sitzung geschrieben: „ich habe es eben auf einem macos installiert und die version abgefragt“; auf die Frage nach `uname -sm` und `brew --version`: „Darwin arm64 -- Homebrew 7.0.10“; auf die Frage nach Befehl und Version: „\"brew install pt9912/pgwire-recorder/pgwire-recorder\" ist richtig -- \"pgwire-recorder 0.2.0\"“. Zur Stufe *Bericht*: „Ja, Bericht gilt (Recommended)“ (Entscheidung des Nutzers vom 2026-10-24). Die Merkmale nach `SPEC-067` Punkt 20 (6), je als Angabe des Nutzers: Datum 2026-10-24; `uname -s` Darwin, `uname -m` arm64; Homebrew 7.0.10; Befehl `brew install pt9912/pgwire-recorder/pgwire-recorder`; Zeile von `pgwire-recorder version`: `pgwire-recorder 0.2.0`. **Das ist ein Bericht, kein Nachweis:** Er trägt die Schritte (1) und (2) und die Zusage des Happy Path von `LH-FA-19`, als berichtet. Nicht getragen: Schritt (3) (`replay`-Lauf, Zeilen, Exit-Code), `brew test`, Ausgabe und Exit-Code der Installation. Geprüft ist nur die Form; dass die Angabe wahr ist, prüft niemand (akzeptiertes Negativ der Spezifikation).

**Wirkung** (`SPEC-067` Punkt 20 (6), *Wirkung*): allein das Wort „berichtet“ in Handbuch, `README.md` und `docs/maintainer/releasing.md`, mit Kombination und Version; „gefahren“ und „belegt“ bleiben dem Nachweis (Linux auf `amd64`). Freigabe-Eintrag, Abdeckungstabellen und `make homebrew-nachweis` bleiben unberührt. Ein Nachweis derselben Kombination löst den Bericht ab; ein Bericht zu anderer Version ersetzt ihn nicht.

**Welle:** `welle-v1-nachtrag` nimmt an (geprüft am 2026-10-24, `AGENTS.md` §3.13): Die Welle liegt flach unter `docs/plan/planning/`, nicht in `done/`. Ihr Closure-Trigger (§3) verlangt Abnahmeszenario 11 „auf macOS und Linux mit dem echten Tap nachgewiesen“; mit einem Bericht ist der macOS-Teil **nicht nachgewiesen**, der Replay-Teil ohne Beleg. Ob der Trigger gesenkt wird oder die Welle mit Grenze schließt, entscheidet der Nutzer (Optionen im Bericht des Planners vom 2026-10-24); dieser Plan ändert den Trigger nicht. Die Stellung des Slice in der Welle (§4 der Welle, Reihenfolge 6) bleibt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein Nachweis für macOS auf `arm64`, eine Beleg-Datei unter `docs/reviews/` oder eine Summe — der Bericht verlangt sie nicht (`SPEC-067` Punkt 20 (6), *Ort*); sie machten ihn nicht zum Nachweis. Der Lauf von `bash tools/harness/homebrew-nachweis.sh --schritte` auf dem Mac wäre ein anderer Vorgang (Klasse 3) und löste den Bericht ab; der Nutzer hat ihn am 2026-10-24 abgelehnt („Ich habe dir doch schon alles gegeben“).
- Eine Aussage zum `replay`-Lauf, zu `brew test`, zur Ausgabe oder zu Exit-Codes der Installation auf macOS — für sie liegt keine Angabe vor; die Grenze steht in `LH-FA-19.a`, der Bericht trägt sie nicht (kein Folge-Slice, solange keine Angabe kommt).
- Eine Zusage für macOS auf `amd64` und Linux auf `arm64` — keine Angabe; Grenze in `LH-FA-19.a`, je Kombination einzeln (Bestand, bewusst).
- Eine weitere Änderung von `LH-FA-19.a`, `SPEC-067` oder des Lastenhefts — die Stufe steht seit `4279f2b3`; der Slice liefert nur Doku. Das Lastenheft nennt keine Plattformgrenze und bleibt unberührt.
- Die Wörter „gefahren“, „belegt“, „geprüft“ für die Kombination in der Doku — die Wirkung des Berichts ist allein „berichtet“ (§3.11).
- Ein Gate, eine Gegenprobe oder eine Abdeckungs-Deklaration — `SPEC-067` Punkt 20 (7); die Abdeckungstabellen führen `LH-FA-19` nicht.
- Eine Änderung am Skript `tools/harness/homebrew-nachweis.sh` — Schicht-Abgrenzung: kein Werkzeug-Code.
- Ein Eintrag in `docs/maintainer/freigaben/0.2.0.md` — Belege des Taps gehören dem Nachweis (`SPEC-067` Punkt 17 (8)); der Bericht lässt ihn unberührt.
- Ein neuer Release oder Tag — anderer Vorgang.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Nutzer- und Wartungs-Doku (eine Schicht, `AGENTS.md` §3.13): Handbuch §2 *Installation* (Satz „Gefahren ist diese Installation …“, die Sätze zu den Binaries für macOS und `arm64`), `README.md` (Zeile zu Homebrew, Satz zu den Binaries) und `docs/maintainer/releasing.md` §1, §7, §10 (Satz zu `darwin` und `windows`) und Änderungshistorie nennen macOS auf `arm64`, Version `0.2.0`, als **berichtet** (nicht gefahren, nicht belegt, ohne Replay-Aussage) und führen macOS auf `amd64` und Linux auf `arm64` weiter als Grenze, Linux auf `amd64` weiter als gefahren (§3.11: jeder Satz gegen `LH-FA-19.a` und `SPEC-067` Punkt 20 (6) gelesen; das Handbuch ohne Verweis auf Spezifikation, Slice oder Review, ohne Chronik).
- [x] Zitat und Merkmale des Berichts stehen in §1 oder §7 dieses Plans (Datum, Wortlaut, `uname`, Homebrew, Befehl, Versionszeile); keine Beleg-Datei, keine Summe.
- [x] `make docs-check` und `make kopf-check` grün; `make gates` grün (kein Code berührt; der Stop-Hook liest Positivliste und schnelle Kette).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben: Der Eintrag `BEO-REPO/homebrew-nachweis-ohne-macos-und-arm64` bekommt seinen Ausgang (*geplant* mit dieser Kennung, bei Closure: Bericht für macOS/arm64 geliefert; offen bleiben der Replay-Teil dort und die zwei Kombinationen ohne Angabe; Urteil des Planners, vorgemerkt am 2026-10-24).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen (im Repo mit Wellen von der Welle-Closure geprüft).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/user/benutzerhandbuch.md` | update | Installation: macOS `arm64` berichtet (Version 0.2.0); nur was Punkt 20 (6) trägt (§3.11) |
| `README.md` | update | Zeile zu Homebrew und Satz zu den Binaries, wie das Handbuch |
| `docs/maintainer/releasing.md` | update | §1, §7 (Grenze „macOS“ wird „macOS `amd64`“, `arm64` berichtet), §10 (Satz zu `darwin` und `windows`), Änderungshistorie mit neuer Zeile und Kopf-Version |

- Eingabe: das Zitat in §1 (kein Lauf, keine Datei). Spezifikation und Lastenheft bleiben in diesem Slice unberührt (§1). Kein Eintrag in `docs/user/abdeckung-*.md`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): (1) Die Stufe *Bericht* steht auf `main` (`4279f2b3`, erfüllt); (2) der Architect hat sie geschrieben, die Randformen stehen in `SPEC-067` Punkt 20 (6) (erfüllt); (3) WIP-Limit frei. Die Ausgabe eines Skriptlaufs ist kein Start-Trigger mehr. `Verantwortlich:` trägt mit dem Anspruch Person und Zweig (`pt9912 (Zweig impl-mac1)`). Basis der Größenmessung ist der Commit des Übergangs nach `in-progress/`; sein Hash steht in der Marke unten.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `open` (blockiert): Ein Merkmal des Berichts fehlt oder stimmt nicht (Version der Formel im Tap anders als `0.2.0`, `uname` nennt eine andere Kombination) oder der Nutzer widerruft die Angabe — die Kombination bleibt Grenze; der Grund wird beim Übergang nachgetragen. Kein Carveout nötig (kein rotes Gate).
- `in-progress` → `next`: **zu groß**, wenn die Summe die **Marke 1500** erreicht (Metrik: `AGENTS.md` §3.14, dort der Befehl; Pfade des Plans: `docs/user docs/maintainer README.md`; Basis: der Commit des Übergangs nach `in-progress/`, Basis: `34a8cf11`). Die Marke ist keine Schätzung.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; das Zitat des Berichts steht in §1 (Merkmale) und wird in §7 genannt; Review und Verifikation (anderer Kontext, jede geänderte Zeile von Handbuch, `README.md` und `releasing.md` gegen `LH-FA-19.a` und `SPEC-067` Punkt 20 (6) gelesen: nur „berichtet“); `make docs-check`, `make kopf-check`, `make gates` grün; Closure-Notiz mit Lerneintrag. Danach schließt der Planner `welle-v1-nachtrag` erst nach der Entscheidung des Nutzers über den Closure-Trigger der Welle (§3: Abnahmeszenario 11 „auf macOS und Linux nachgewiesen“; mit einem Bericht für macOS/arm64 ist er ohne Änderung nicht erreicht).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12, Plattform; keine neue Wire-, Datei- oder Formatgrenze, kein neuer Vertrag; die Randformen des Berichts stehen in `SPEC-067` Punkt 20 (6) und sind dort entschieden):

- **Plattform:** berichtet ist macOS auf `arm64` (2026-10-24, Version 0.2.0, Homebrew 7.0.10); Linux auf `amd64`: Nachweis (`docs/reviews/2026-10-24-nachweis-slice-release-v0-2-0.md`); macOS auf `amd64` und Linux auf `arm64`: Grenze (`LH-FA-19.a`); Windows: nicht über Homebrew.
- **Merkmale:** Fehlt eines (Datum, `uname`, Homebrew-Version, Befehl, Versionszeile) oder stammt die Angabe nicht vom Nutzer, ist die Kombination Grenze (`SPEC-067` Punkt 20 (6), *Merkmale*); alle fünf stehen im Zitat in §1.
- **Version:** Der Bericht gilt für `0.2.0`; die Doku nennt `0.2.0`, kein „aktuelle Version“. Ein Bericht zu anderer Version ersetzt ihn nicht.
- **Ablösung:** Ein Nachweis derselben Kombination löst den Bericht ab; dann wird aus „berichtet“ „gefahren“ (anderer Slice).

Risiken:

- Die Doku sagt mehr, als der Bericht trägt („gefahren“, „belegt“, `replay` auf macOS, `brew test`). — **Ausgang:** entfallen mit Begründung, wenn Review und Verifikation jeden Satz gegen Punkt 20 (6) lesen (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`); sonst eingetreten: Befund im Review, im Slice behoben.
- Änderung an `LH-FA-19.a` oder `SPEC-067` nötig (Stufe unterhalb von (6)). — **Ausgang:** eingetreten (2026-10-24): Die Stufe *Bericht* ist mit Entscheidung des Nutzers vom 2026-10-24 geliefert (`4279f2b3`); dieser Plan ist nachgezogen (`AGENTS.md` §3.9), Zählung in §8 neu.
- Der Replay-Teil von Abnahmeszenario 11 bleibt auf macOS ohne Beleg, und der Closure-Trigger von `welle-v1-nachtrag` wird nicht erreicht. — **Ausgang:** eingetreten: Entscheidung des Nutzers über Trigger oder Grenze, Vorlage durch den Planner am 2026-10-24; kein Folge-Slice, solange der Nutzer keinen Nachweis liefert (sonst legt der Planner einen Slice mit Kennung an).
- Die Angabe des Nutzers erweist sich als unvollständig oder falsch (Formel im Tap in anderer Version, Tippfehler). — **Ausgang:** weiter offen: → `BEO-REPO/homebrew-nachweis-ohne-macos-und-arm64` im Register (Rückführung nach `open/`, §4); die Wahrheit der Angabe prüft niemand (akzeptiertes Negativ der Spezifikation).
- macOS hat Eigenheiten, die Linux nicht zeigt (bash 3.2, Gatekeeper). — **Ausgang:** entfallen mit Begründung: Der Bericht trägt nur Installation und Versionszeile, die der Nutzer auf dem Mac erhielt; für `replay` und `brew test` sagt der Slice nichts.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
Zeiger sind so zu schreiben, wie sie vom Ruheort `done/` auflösen.

**Belege des Implementers (Zweig `impl-mac1`, Stand `a6bffa1a`):**

- `make docs-check`: 0 Befunde (1467 Dateien); `make kopf-check`: Exit 0.
- `make gates`: der erste volle Lauf war rot, `TestE2EPlayVergleichAlle/Umgebungsvariable` überschritt die Frist unter Last; der Diff berührt keinen Code, das Rot ist nicht reproduziert und die Ursache nicht untersucht (Register: `BEO-REPO/test-faellt-unter-cpu-last-aus`, siebtes Auftreten). Der zweite volle Lauf am selben Stand war grün. Der Endstand ist grün, ein Carveout entfällt.
- Marke 1500 (`AGENTS.md` §3.14, Basis `34a8cf11`, Pfade `docs/user docs/maintainer README.md`): Summe **41**.
- Review (`docs/reviews/2026-10-24-review-slice-homebrew-nachweis-macos-arm64.md`, Stand `34962d3c`): 0 HIGH, 0 MEDIUM, 3 LOW (F-1 bis F-3, behoben in `5dae30d9`), 3 INFO. Verifikation (`docs/reviews/2026-10-24-verifikation-slice-homebrew-nachweis-macos-arm64.md`, Stand `a6bffa1a`): konform. Validation: n/a — der Slice liefert Doku zu einem Bericht, kein neues Nutzerverhalten.
- Zitat und Merkmale des Berichts: §1 (Datum 2026-10-24, `Darwin arm64`, Homebrew 7.0.10, Befehl, Versionszeile `pgwire-recorder 0.2.0`); keine Beleg-Datei, keine Summe.

**Gegenstand:** geliefert (Liefer-Punkte der DoD abgehakt).

**Risiko-Ausgänge (aus §6):**

1. Die Doku sagt mehr, als der Bericht trägt — *entfallen:* Review und Verifikation lasen jeden geänderten Satz gegen `LH-FA-19.a` und `SPEC-067` Punkt 20 (6); `grep` nach „gefahren“, „belegt“ und „geprüft“ findet für macOS/`arm64` nur Verneinungen. Die drei LOW des Reviews betrafen Wortlaut, nicht diese Aussage.
2. Änderung an `LH-FA-19.a` oder `SPEC-067` nötig — *eingetreten:* geliefert mit Commit `4279f2b3` (Entscheidung des Nutzers vom 2026-10-24); der Plan ist nachgezogen.
3. Replay-Teil ohne Beleg, Closure-Trigger von `welle-v1-nachtrag` nicht erreicht — *eingetreten:* Entscheidung des Nutzers vom 2026-10-24, Option B „Trigger unverändert, mit Grenze schließen“ (so vom Koordinator übermittelt); die Welle schließt mit der Grenze für den Replay-Teil auf macOS `arm64`, für macOS auf `amd64` und Linux auf `arm64`; ihr Plan bleibt unverändert, ihre Closure ist nicht Teil dieses Slice. Kein Folge-Slice, solange der Nutzer keine Ausgabe liefert (§1, Abgrenzung); ein Nachweis löst den Bericht ab.
4. Die Angabe des Nutzers ist unvollständig oder falsch — *weiter offen:* → `BEO-REPO/homebrew-nachweis-ohne-macos-und-arm64` im Register; die Rückführung nach `open/` entfällt, da der Slice geschlossen ist, ein Widerruf wäre ein Auftreten dort.
5. macOS-Eigenheiten (bash 3.2, Gatekeeper) — *entfallen:* Der Bericht trägt nur Installation und Versionszeile; für `replay` und Paket-Test sagt der Slice nichts.

**Lerneintrag (benannte Spec-Lücke, als Grenze geschlossen):** Die Aussage, die Angabe des Nutzers trage als Nachweis, war falsch: `SPEC-067` Punkt 20 verlangte die Ausgabe des Laufs, und eine Angabe ohne Ausgabe ist ein Bericht. Die Stufe *Bericht* (`LH-FA-19.a`, `SPEC-067` Punkt 20 (6), `4279f2b3`) schließt die Lücke nur als Grenze: Sie trägt die Schritte (1) und (2) als berichtet, nicht den Replay-Teil, sie löst den Nachweis nicht ab, und ihre Wirkung ist allein das Wort „berichtet“ mit Kombination und Version. Offen bleibt die Lücke für den Replay-Teil auf macOS `arm64` und für zwei Kombinationen ohne Angabe.

**Register (Schritt 25):** `BEO-REPO/homebrew-nachweis-ohne-macos-und-arm64` Zähler 3, Ausgang `geplant` mit dieser Kennung (der Rest bleibt Grenze, siehe seine `state.md`); `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` 98 (Review F-1, F-3); `BEO-REPO/plan-folgt-korrektur-nicht` 85 (Review F-2); `BEO-REPO/test-faellt-unter-cpu-last-aus` 7 (erster roter Lauf). Review F-4 bis F-6 (INFO) sind keine Pflicht und führen keinen Beleg.

**Lese-Schritt und Rate (`AGENTS.md` §3.15, ohne Welle bei der Slice-Closure):** Neu mit einer Datei unter `evidence/` sind `zusage-im-kommentar-weiter-als-pruefung` (§3.11) und `plan-folgt-korrektur-nicht` (§3.9). Seit `slice-harness-regeln-rueckbau` sind mit diesem Slice neun geschlossen (`slice-v1-abschluss-verbindungsende-upstream-extended`, `slice-version-option`, `slice-play-feldvergleich`, `slice-play-feldvergleich-doku`, `slice-go-toolchain-1-27-3`, `slice-release-v0-2-0`, `slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`, dieser; Reihenfolge nach `git log` der Moves nach `done/`): **keine Rate und kein Vorschlag**, es sind weniger als zehn. Die Rate rechnet die Closure des zehnten Slice dahinter. Neu bei 3 oder mehr ist allein `homebrew-nachweis-ohne-macos-und-arm64` (Ausgang oben); sie hat keine Regel in `AGENTS.md`.

**Paarungen:** Anker — keine Regel verkörpert, entfällt. Folge-Slice — keiner zugewiesen (Risiko 3 ohne Folge-Slice, Begründung dort; der Register-Ausgang `geplant` zeigt auf diesen Slice selbst). Register — fortgeschrieben. Die Prüfung der drei Paarungen im Repo mit Wellen trägt die Welle-Closure von `welle-v1-nachtrag`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist allein `REPO` (`*`, `harness/conventions.md`), Greenfield; sie ist nicht ausdifferenziert und braucht es hier nicht (Doku).

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-24, gemergter Stand): Treffer: `BEO-REPO/homebrew-nachweis-ohne-macos-und-arm64` (2 Belege, `Stand: offen`, Ausgang nicht zugewiesen; dieser Slice ist die Adresse des Vorschlags (b) und verändert den Zähler nicht, ein Eintrag entsteht erst bei Closure); `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (97 Belege; Regel `AGENTS.md` §3.11, deren Prüfung der Slice an Handbuch, README und `releasing.md` trägt); `BEO-REPO/verhalten-nur-unter-linux-geprueft` (3); `BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet` (1; hier gilt: Die Stufe *Bericht* ist Entscheidung des Nutzers, Wortlaut in §1). Zu `AGENTS.md` §3.15: keine Regel wird geändert oder erweitert.

Alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13** (Planner, 2026-10-24, nachgezogen nach der Stufe *Bericht*): **Liefer-Punkte 2** (Doku-Update als ein Punkt über drei Dateien; Zitat und Merkmale in §1 als zweiter Punkt). Die Änderung von `LH-FA-19.a` und `SPEC-067` ist im Commit `4279f2b3` geliefert und kein Punkt dieses Slice. **Schichten 1** (Nutzer- und Wartungs-Doku; Plan und Roadmap zählen nicht). Sendungen an Nehmer: keine; `welle-v1-nachtrag` nimmt den Slice auf (Reihenfolge 6), die Frage des Closure-Triggers geht als Vorlage an den Nutzer, nicht als Eintrag.
