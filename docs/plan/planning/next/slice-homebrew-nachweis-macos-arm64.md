# Slice slice-homebrew-nachweis-macos-arm64: Nachweis der Installation aus dem Tap auf macOS (arm64)

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

**Berührte Spec-Stellen:** `LH-FA-19.a` (gelesen, nicht geändert: der Text ist je Kombination allgemein gefasst) · `SPEC-067` (Punkt 20 (5) und (6), gelesen, nicht geändert)

**Verantwortlich:** pt9912
**Autor:** pt9912. **Datum:** 2026-10-24.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Abnahmeszenario 11 ist für die Kombination **macOS auf `arm64`** nach [`SPEC-067`](../../../../harness/vertraege/SPEC-067-release.md) Punkt 20 belegt: Der Nutzer fährt `bash tools/harness/homebrew-nachweis.sh --schritte` auf seinem Mac, die Ausgabe liegt ungekürzt als Datei unter `docs/reviews/`, und Handbuch, `README.md` und `docs/maintainer/releasing.md` nennen macOS auf `arm64` als gefahren, die übrigen Kombinationen weiter als Grenze.

**Anlass (Angabe der Koordination, kein Beleg).** Der Nutzer hat am 2026-10-24 in der Sitzung geschrieben: „ich habe es eben auf einem macos installiert und die version abgefragt“; auf die Frage nach `uname -sm` und `brew --version`: „Darwin arm64 -- Homebrew 7.0.10“; auf die Frage nach Befehl und Version: „\"brew install pt9912/pgwire-recorder/pgwire-recorder\" ist richtig -- \"pgwire-recorder 0.2.0\"“. Ausgabe und Exit-Code von `brew install` und `brew test` liegen nicht vor, ebenso nicht der Replay-Lauf. **Das ist nach `LH-FA-19.a` und `SPEC-067` Punkt 20 (6) kein Nachweis:** „Eine Angabe ohne die Ausgaben ist kein Nachweis“, und der Nachweis besteht aus drei Schritten (Installation, `version`, startender `replay`-Lauf). Belegt wäre allein die Versionszeile, und auch die nur als Angabe. Die Kombination gilt bis zu einer Datei mit den Ausgaben als nicht geprüft; dieselbe Lesart wandte `slice-erster-release-homebrew-nachweis` auf die macOS-Angabe vom 2026-10-21 an.

**Was für die Mac-Variante der Beleg-Datei gilt** (`SPEC-067` Punkt 20 (6), Muster `docs/reviews/2026-10-24-nachweis-slice-release-v0-2-0.md`): Datei `docs/reviews/<Datum>-nachweis-slice-homebrew-nachweis-macos-arm64.md` (Namensform: Datum, `nachweis`, Kennung des Slice; das Datum ist das des Laufs, nicht das des Plans); Kopf in Prosa mit Datum, Ausgabe von `uname -s` und `uname -m`, Stand des Skripts (Commit) und Tag; dazu die Ausgabe von `brew --version`. **Entfällt gegenüber Linux:** Tag und Digest des Image `homebrew/brew` (am Host läuft kein Container) und der Aufruf `make homebrew-nachweis`; der Aufruf ist `bash tools/harness/homebrew-nachweis.sh --schritte` am Host, auf macOS mit bash 3.2. **Gleich:** die Befehle (a) bis (e) mit voller Ausgabe in einem Codeblock, die Version aus (c), der Exit-Code des Skripts und des Replay-Laufs; gekürzt wird allein ein Geheimnis oder ein kurzlebiger Token (signierte Download-Adresse) mit der Marke `[gekürzt: <Grund>]`, beim Ablegen der Datei. **Was nicht gefordert werden kann:** ein Digest des Image, ein Lauf im Container, und eine Ausgabe, die der Nutzer nicht hat — fehlt (e) oder der Exit-Code, ist die Datei kein Nachweis, und die Kombination bleibt Grenze. **Herkunft:** Der Agent fährt `brew` nicht am Host (`AGENTS.md` §3.1, `SPEC-067` Punkt 20 (5)); der Kopf der Datei sagt, dass der Nutzer den Lauf fuhr und die Ausgabe einfügte, und der Review prüft Form und Stimmigkeit (Zeitstempel, `uname`, Version), nicht die Echtheit.

**Welle:** `welle-v1-nachtrag` nimmt an (geprüft am 2026-10-24, `AGENTS.md` §3.13): Die Welle liegt flach unter `docs/plan/planning/`, nicht in `done/`; ihr Closure-Trigger (§3) verlangt Abnahmeszenario 11 „auf macOS und Linux mit dem echten Tap“ — dieser Slice liefert das *Mehr* gegenüber den Slice-DoDs, das der Absatz *Stand von Schritt 5* in §5 der Welle offen ließ; §6 der Welle schließt macOS nicht aus („weitere Plattformen“ meint Plattformen jenseits von macOS, Linux und Windows). Wellenlos wäre er falsch, weil die Welle sonst mit der Grenze für macOS schlösse und ihr Trigger leer bliebe. Die Umplanung steht im Drift-Log der Roadmap und in §4, §5 der Welle; Meilenstein M5 unverändert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Eine Zusage für macOS auf `amd64` und für Linux auf `arm64` — für sie liegt keine Ausgabe vor; die Grenze steht in `LH-FA-19.a` und bleibt (Bestand, bewusst; die Zusage gilt je Kombination einzeln).
- Eine Änderung von `LH-FA-19.a`, `SPEC-067` Punkt 20 oder des Lastenhefts — der Text der Spezifikation nennt keine Kombination als belegt, sondern verlangt „Nachweis je Kombination mit Ausgaben“; der belegte Stand steht im Handbuch, in `README.md` und in `releasing.md`. Das Lastenheft nennt keine Plattformgrenze (`LH-FA-19` und Abnahmeszenario 11 sagen „macOS oder Linux“); eine Änderung hätte nichts zu ändern. (Regel „Versionierung des Lastenhefts“ in `harness/conventions.md` geprüft: Status Draft, frei änderbar, kein Bump; er entfiele auch ohne Draft, weil keine Aussage sich ändert.)
- Ein Teilbeleg unterhalb von Punkt 20 (6) (nur Versionszeile, Angabe zitiert, Handbuch „Installation und Version auf macOS arm64 berichtet“) — es wäre eine Lockerung der Regel „Eine Angabe ohne die Ausgaben ist kein Nachweis“, die der Nutzer am 2026-10-21 so gesetzt hat; Änderung von `LH-FA-19.a` und `SPEC-067`, nur mit Entscheidung des Nutzers und des Architect, und nicht Gegenstand dieses Slice (Klasse 3, anderer Vorgang).
- Ein Gate, eine Gegenprobe oder eine Abdeckungs-Deklaration — `SPEC-067` Punkt 20 (7): das Ziel urteilt nicht; die Abdeckungstabellen führen `LH-FA-19` nicht (`grep` leer).
- Eine Änderung am Skript `tools/harness/homebrew-nachweis.sh` — Schicht-Abgrenzung: kein Produkt- und kein Werkzeug-Code; findet der Lauf auf macOS einen Fehler im Skript oder eine Zeile, die Punkt 20 (3) nicht zulässt, geht er als Randform-Rückgabe an den Architect und die Änderung an einen Folge-Slice.
- Ein Eintrag in `docs/maintainer/freigaben/0.2.0.md` — der Eintrag gehört dem Release; Belege des Taps gehören dem Nachweis (`SPEC-067` Punkt 17 (8), Punkt 20 (6)).
- Ein neuer Release oder Tag — anderer Vorgang; der Nachweis läuft gegen die Formel, die im Tap liegt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), Abnahmeszenario 11, macOS/arm64: Die Beleg-Datei `docs/reviews/<Datum>-nachweis-slice-homebrew-nachweis-macos-arm64.md` liegt nach Form und Inhalt von §1 (`SPEC-067` Punkt 20 (6)): Kopf, Befehle (a) bis (e) ungekürzt, Version `pgwire-recorder <Version>` aus (c), Exit-Code des Replay-Laufs 0, eine Zeile `level=INFO` und höchstens die eine Zeile `PGR-W2001`; §7 nennt Pfad und Summe (`sha256sum`).
- [ ] Nutzer- und Wartungs-Doku (eine Schicht, `AGENTS.md` §3.13): Handbuch §2 *Installation* (Satz „Gefahren ist diese Installation …“, die Sätze zu den Binaries für macOS und `arm64` in der Voraussetzung und beim Binary vom Release), `README.md` (Zeile zu Homebrew, Satz zu den fünf Binaries) und `docs/maintainer/releasing.md` §1, §7 und Änderungshistorie nennen macOS auf `arm64` als gefahren, in der Version des Belegs, und führen macOS auf `amd64` und Linux auf `arm64` weiter als Grenze (§3.11: jeder Satz gegen die Datei gelesen; das Handbuch ohne Verweis auf Spezifikation, Slice oder Review, ohne Chronik).
- [ ] `make docs-check` und `make kopf-check` grün; `make gates` grün (kein Code berührt; der Stop-Hook liest Positivliste und schnelle Kette).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: Der Eintrag `BEO-REPO/homebrew-nachweis-ohne-macos-und-arm64` bekommt seinen Ausgang (*geplant* mit dieser Kennung, bei Closure bei gelieferter Kombination macOS/arm64 mit Rest-Grenze für zwei Kombinationen; Urteil des Planners).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen (im Repo mit Wellen von der Welle-Closure geprüft).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/reviews/<Datum>-nachweis-slice-homebrew-nachweis-macos-arm64.md` | neu | Beleg nach `SPEC-067` Punkt 20 (6); wird in dem Commit abgelegt, der §7 den Pfad nennt, und danach nicht geändert |
| `docs/user/benutzerhandbuch.md` | update | Installation: macOS `arm64` gefahren; nur was die Datei belegt (§3.11) |
| `README.md` | update | Zeile zu Homebrew und Satz zu den Binaries, wie das Handbuch |
| `docs/maintainer/releasing.md` | update | §1, §7 (Grenze „macOS“ wird „macOS `amd64`“), Änderungshistorie mit neuer Zeile und Kopf-Version |

- Eingabe: Die Ausgabe des Laufs kommt vom Nutzer (er fährt das Skript auf seinem Mac); der Implementer legt sie ab und fährt keinen Teil davon. Ohne diese Eingabe beginnt der Slice nicht (§4).
- Spezifikation und Lastenheft bleiben unberührt (§1). Kein Eintrag in `docs/user/abdeckung-*.md`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): (1) Die **Ausgabe** des Skriptlaufs auf dem Mac liegt vor (Nutzer, vollständig: Befehle (a) bis (e), `uname -s`, `uname -m`, `brew --version`, Exit-Codes); (2) der Architect hat den Plan gelesen und bestätigt, dass keine Änderung an `LH-FA-19.a`, `SPEC-067` und Lastenheft nötig ist, und die Randformen von §6 geprüft (`AGENTS.md` §3.12), vor dem Anspruch; (3) WIP-Limit frei. `Verantwortlich:` trägt dann Person und Zweig (`pt9912 (Zweig impl-mac1)`). Den Hash des Übergangs trägt §4 in der Marke unten ein.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `open` (blockiert): Der Nutzer liefert keine vollständige Ausgabe, oder sie zeigt eine Zeile, die `SPEC-067` Punkt 20 (3) nicht zulässt, oder ein Schritt scheitert — die Kombination bleibt Grenze; der Grund wird beim Übergang nachgetragen. Kein Carveout nötig (kein rotes Gate).
- `in-progress` → `next`: **zu groß**, wenn die Summe die **Marke 1500** erreicht (Metrik: `AGENTS.md` §3.14, dort der Befehl; Pfade des Plans: `docs/user docs/maintainer README.md`; Basis: der Commit des Übergangs nach `in-progress/`, Hash beim Übergang eintragen: `—`). Erwartet ist eine Größenordnung unter 60 Zeilen; die Marke ist keine Schätzung.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig; die Beleg-Datei liegt mit Pfad und Summe in §7; Review und Verifikation (anderer Kontext, jede geänderte Zeile von Handbuch, `README.md` und `releasing.md` gegen die Datei gelesen); `make docs-check`, `make kopf-check`, `make gates` grün; Closure-Notiz mit Lerneintrag; danach schließt der Planner `welle-v1-nachtrag` (Closure-Trigger der Welle: Abnahmeszenario 11 auf macOS und Linux; `make gates`; Closure-Notiz).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12, Plattform; keine neue Wire-, Datei- oder Formatgrenze und kein neuer Vertrag, deshalb nur Plattform und Eingabe):

- **Plattform:** Geprüft wird macOS auf `arm64` (Ausgabe des Nutzers, `uname -s` Darwin, `uname -m` arm64, Homebrew 7.0.10 laut Angabe, bash 3.2 laut Skript-Kopf). Linux auf `amd64`: belegt (`docs/reviews/2026-10-24-nachweis-slice-release-v0-2-0.md`). macOS auf `amd64` und Linux auf `arm64`: Grenze, entschieden in `LH-FA-19.a` (Zusage je Kombination); Windows: nicht über Homebrew (`LH-FA-19` Negative).
- **Fehlende Ausgabe:** Fehlt einer der Befehle (a) bis (e), der Exit-Code oder `uname`, ist die Datei kein Nachweis (Rückführung in §4); eine Teilausgabe wird nicht abgelegt und nicht als Nachweis geführt.
- **Zeile außerhalb von Punkt 20 (3)** (`ERROR`, anderes `WARN`, zweite `PGR-W2001`, Zeile ohne `level=`): Der Nachweis gilt als nicht erbracht, Randform-Rückgabe an den Architect.
- **Version:** Die Formel des Taps trägt die höchste stabile Version zum Zeitpunkt des Laufs; liegt sie über `0.2.0`, nennen Datei und Doku die Version des Laufs, nicht `0.2.0`.
- **Ablage:** Die Datei wird einmal abgelegt und nicht geändert; ein neuer Lauf ist eine neue Datei (`SPEC-067` Punkt 20 (6)).

Risiken:

- Die Ausgabe des Nutzers lässt sich nicht vollständig beschaffen (Terminal-Verlauf verloren); die Angabe vom 2026-10-24 trägt dann nicht. — **Ausgang:** weiter offen: → `BEO-REPO/homebrew-nachweis-ohne-macos-und-arm64` im Register (der Slice geht dann nach `open/`, nicht nach `done/`; der Ausgang entsteht bei Closure, falls es so kommt).
- macOS hat Eigenheiten, die Linux nicht zeigt (bash 3.2, `brew test` als Entwicklerbefehl, `127.0.0.1:55432` belegt, Gatekeeper für ein nicht signiertes Binary, das nicht über `curl` kommt). — **Ausgang:** entfallen mit Begründung, wenn die Ausgabe grün ist (der Lauf zeigt es); sonst eingetreten: Folge-Slice mit Kennung, die der Planner bei Eintritt anlegt.
- Das Handbuch sagt mehr, als die Datei belegt (zum Beispiel „macOS“ statt „macOS auf `arm64`“, oder die Binaries für macOS als geprüft statt allein das Binary für `arm64`). — **Ausgang:** entfallen mit Begründung, wenn Review und Verifikation jeden Satz gegen die Datei lesen (`BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`); sonst eingetreten: Befund im Review, im Slice behoben.
- Der Architect hält eine Änderung an `LH-FA-19.a` oder `SPEC-067` für nötig (zum Beispiel eine Stufe unterhalb von (6)). — **Ausgang:** eingetreten: der Plan wird vor dem Anspruch nachgezogen (`AGENTS.md` §3.9), mit Entscheidung des Nutzers; Zählung in §8 neu (dann Spezifikation als zweite Schicht neben der Doku).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
Zeiger sind so zu schreiben, wie sie vom Ruheort `done/` auflösen.

(entsteht bei Closure)

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt ist allein `REPO` (`*`, `harness/conventions.md`), Greenfield; sie ist nicht ausdifferenziert und braucht es hier nicht (Doku und Beleg-Datei).

**Vorgelagert — offene Beobachtungen sichten** (Register gelesen am 2026-10-24, gemergter Stand): Treffer: `BEO-REPO/homebrew-nachweis-ohne-macos-und-arm64` (2 Belege, `Stand: offen`, Ausgang nicht zugewiesen; dieser Slice ist die Adresse des Vorschlags (b) und verändert den Zähler nicht, ein Eintrag entsteht erst bei Closure); `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung` (97 Belege; Regel `AGENTS.md` §3.11, deren Prüfung der Slice an Handbuch, README und `releasing.md` trägt); `BEO-REPO/verhalten-nur-unter-linux-geprueft` (3); `BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet` (1; hier gilt: die Angabe vom 2026-10-24 ist Angabe, die Entscheidung über Doku und Slice ist das „ok“ des Nutzers). Zu `AGENTS.md` §3.15: keine Regel wird geändert oder erweitert.

Alle berührten Sub-Areas GF.

**Zählung nach `AGENTS.md` §3.13** (Planner, 2026-10-24, beim Anlegen): **Liefer-Punkte 2** (Beleg-Datei; Doku-Update als ein Punkt über drei Dateien). Kein Punkt für Spezifikation und Lastenheft, weil beide unberührt bleiben (§1). **Schichten 1** (Nutzer- und Wartungs-Doku; Beleg-Datei und Plan zählen nicht). Sendungen an Nehmer: keine; `welle-v1-nachtrag` nimmt den Slice auf (Eintrag in §4 und §5 der Welle im selben Zug). Über der Grenze (3 Punkte, 2 Schichten) wäre der Slice erst, wenn der Architect eine Änderung der Spezifikation verlangt (§6, letztes Risiko): dann 3 Punkte, 2 Schichten, noch innerhalb.
