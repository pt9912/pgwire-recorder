# Welle welle-v1-nachtrag — v1-Nachtrag: Homebrew und CancelRequest ohne Schlüssel — Closure-Notiz

> **Zitier-Form** *(bleibt stehen — Norm, kein Ausfüll-Hinweis).* Dieses
> Artefakt friert ein; was es zitiert, bewegt sich weiter. Deshalb: **Kennung,
> nicht Adresse** — `slice-<Kennung>` statt seines Lifecycle-Pfads, `make <target>`
> statt eines Links auf die Sensor-Datei, eine Baseline-Stelle als
> `v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt> statt als Link
> (Baseline-Regelwerk `grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — beim Ausfüllen mit dem adoptierten Tag schreiben).

**Welle:** welle-v1-nachtrag
**Abschluss:** 2026-10-24
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3.

**Die Welle schließt mit einer Grenze, nicht mit erfülltem Closure-Trigger.** Der Trigger in §3 des Welle-Plans nennt Abnahmeszenario 11 „auf macOS und Linux mit dem echten Tap nachgewiesen“. Das ist nicht erreicht. Der Nutzer hat am 2026-10-24 entschieden (Option B, Wortlaut „Trigger unverändert, mit Grenze schließen“, so vom Koordinator übermittelt), den Trigger im Text zu belassen und mit der Grenze zu schließen; der Welle-Plan trägt die Entscheidung in §3.

| Kombination (Abnahmeszenario 11) | Stand bei Closure |
|---|---|
| Linux `amd64` | **nachgewiesen** mit der Formel `0.1.1` aus dem echten Tap (`slice-erster-release-homebrew-nachweis`, Werkzeug `make homebrew-nachweis`, Beleg in dessen §7, von der Verifikation mit eigenem Lauf reproduziert) |
| macOS `arm64` | **berichtet**, Stufe *Bericht* (`LH-FA-19.a`, `SPEC-067` Punkt 20 (6), `slice-homebrew-nachweis-macos-arm64`): Installation und `version`, Version `0.2.0`, Homebrew 7.0.10, Angabe des Nutzers vom 2026-10-24; **der Replay-Teil (Schritt (3), Paket-Test) dort ohne Beleg** |
| macOS `amd64` | **Grenze** (`LH-FA-19.a`), keine Angabe, kein Beleg |
| Linux `arm64` | **Grenze** (`LH-FA-19.a`), keine Angabe, kein Beleg |

**Offen, nicht erfüllt:** der Nachweis auf macOS (beide Architekturen) und auf Linux `arm64`, und der Replay-Teil auf macOS `arm64`. Der Bericht ist kein Nachweis und löst ihn nicht ab. Meilenstein M5 bleibt in der Roadmap `offen` (sein Trigger nennt Abnahmeszenario 11 als nachgewiesen).

Übrige Lieferung, 7 Slices in `done/`, alle als Datei dort (Abzählung `ls` und `grep` auf `Welle:`, Lauf des Planners am 2026-10-24):

- `slice-v1-abschluss-cancel-ohne-schluessel`: ein `CancelRequest` ohne Abbruchschlüssel wird in `record` und `replay` behandelt (`LH-FA-05.e`, `LH-FA-13`); Länge 12 ist `PGR-W3001` und endet ohne Antwort, Längen 13 bis 15 und die übrigen Ränder sind `PGR-E6001`, je Länge ein Test mit roter Mutation (Tabelle in dessen §7).
- `slice-v1-abschluss-binary-namen`: die sechs Binaries heißen `pgwire-recorder-<os>-<arch>` ohne Version im Namen, `SHA256SUMS` dazu (`LH-FA-01`, `LH-FA-19`; `SPEC-055` Punkt 3).
- `slice-v1-abschluss-homebrew`: das Verfahren der Formel (Job `homebrew` in `release.yml`, Vorlage, Probe-Tap als Bare-Repository in der Gegenprobe; `SPEC-067` Punkt 19).
- `slice-v1-abschluss-release-v0-1-1`: `VERSION`, Eintrag der Freigabe und Bezug im Handbuch für den stabilen Tag `v0.1.1`, den Claude am 2026-10-21 nach dem Ja des Nutzers setzte (Lauf `39551891874`).
- `slice-erster-release-homebrew-nachweis`: der Nachweis für Linux `amd64` (Tabelle oben).
- `slice-homebrew-nachweis-macos-arm64`: Doku der Stufe *Bericht* für macOS `arm64` (Tabelle oben).
- `slice-v1-abschluss-verbindungsende-upstream-extended`: das Ende einer Extended-Interaktion durch den Server ist gegen PostgreSQL 17.11 belegt, Fixtures mit Sync und mit Flush, vier Aufzeichnungs- und Replay-Tests und zwei ohne PostgreSQL; nur Tests und Testdaten (`LH-FA-25`, `LH-FA-18`); 280 Zeilen gegen Marke 1500.
- Nicht Mitglieder, vom Werkzeug dennoch mit archiviert (Archivierung unten): `slice-release-v0-2-0`, `slice-version-option`, `slice-play-feldvergleich`, `slice-go-toolchain-1-27-3`, `slice-harness-gegenproben-umzug-dev`, `slice-harness-regeln-rueckbau` (wellenlos seit der letzten Closure oder mit Nennung der Welle im Text, Zählung des Werkzeugs: 9 Mitglieder, 4 wellenlos).

Die Anforderungen (MUSS und SOLL) sind nach dem Stand der Slices umgesetzt; die Anforderung `LH-FA-19` ist es mit der Grenze oben (die Formel liegt im Tap und installiert unter Linux `amd64` und, berichtet, unter macOS `arm64`).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3.

- Der Schnitt des Nachweises nach Kombination: Ein Beleg für Linux `amd64` stand früh, die Grenze für den Rest blieb im Lastenheft sichtbar, statt die Welle auf einen fehlenden Mac warten zu lassen.
- Verfahren der Formel (Probe-Tap als Bare-Repository) und echter Tap getrennt: Das Verfahren war ohne Netz und ohne Tag prüfbar; was erst der Lauf des Tags zeigen konnte, fand er (Rechte des Tap-Geheimnisses, drei rote Versuche mit `error: 403`), und beides steht im Register.
- Die Rollen-Sequenz hielt: Review und Verifikation fanden je Slice Befunde, die der Implementer nicht sah (zum Beispiel zwei grüne verhaltensändernde Mutanten bei `CancelRequest`).

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3.

- **Der Closure-Trigger nannte eine Umgebung, für die kein Träger stand.** Bei der Eröffnung nannte §3 macOS und Linux, ohne zu sagen, wer auf einem Mac fährt; die Entscheidung dazu fiel erst 13 Tage nach der Eröffnung vom 2026-10-11 (Schritt 5 am 2026-10-21 mit der Grenze, Schritt 6 und die Entscheidung des Nutzers am 2026-10-24). Folge: die Stufe *Bericht* (`LH-FA-19.a`), und der Trigger wurde nicht geändert, sondern mit Grenze geschlossen.
- **Der letzte Slice der Welle entstand nach dem ersten Nachweis** (`slice-homebrew-nachweis-macos-arm64`, am 2026-10-24 auf Entscheidung des Nutzers); die Welle war damit nicht geschlossen, als der Nachweis geschlossen war.
- **Binaries ohne Version im Namen** (`slice-v1-abschluss-binary-namen`, Entscheidung des Nutzers vom 2026-10-20) kam in die Welle, weil Formel und Vorlage die Namen brauchen; der erste Tag mit dem Schema war `v0.1.1`, nicht `v0.1.0`.
- **Gegenproben am Wellenende nicht gefahren** — siehe *Gegenproben am Wellenende* unten.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 3 (Lese-Schritt: jeder Eintrag mit ≥ 3 Dateien unter `evidence/`).

**Lese-Schritt.** Eintrag mit ≥ 3 Belegen, dessen Zähler in dieser Welle wuchs (Zähler vom 2026-10-24):

- **Spec-Lücke** benannt und nur als Grenze geschlossen: Abnahmeszenario 11 (`LH-FA-19`) verlangt einen Nachweis auf macOS und Linux, das Projekt hat keinen Träger dafür (kein Mac, kein Apple-Silicon-Host unter den Host-Voraussetzungen `bash`, `git`, `docker`); der Beleg ist ein Lauf des Betreibers. Aufgelöst als Grenze über `LH-FA-19.a` und `SPEC-067` Punkt 20 (6) (Stufe *Bericht*, Entscheidung des Nutzers vom 2026-10-24). Sie ist **nicht erfüllt**: Der Bericht trägt Installation und `version`, nicht den Replay-Teil, und sagt für macOS `amd64` und Linux `arm64` nichts.
  Auslöser: `BEO-REPO/homebrew-nachweis-ohne-macos-und-arm64` (`slice-erster-release-homebrew-nachweis`, `slice-release-v0-2-0`, `slice-homebrew-nachweis-macos-arm64` — 3×; Stand `geplant`, Adresse `slice-homebrew-nachweis-macos-arm64`, der Rest ohne Adresse).

**Einträge über der Schwelle ohne zugewiesenen Ausgang — Zuweisung entscheidet der Nutzer, der Planner weist nichts zu** (Vorschläge unter *Offene Punkte beim Nutzer*):

- `BEO-REPO/liste-gruener-mutanten-unvollstaendig` (4×, Stand `offen`; in dieser Welle das vierte Auftreten durch `slice-v1-abschluss-cancel-ohne-schluessel`): „Grüne Mutanten: keiner“ stand zu Unrecht in §7, zwei grüne verhaltensändernde Mutanten fand erst die Nach-Verifikation. Ein Sensor dafür ist `slice-harness-mutation` (Anker im Register, kein Ausgang zugewiesen).
- `BEO-REPO/wiederanstoss-eines-roten-release-laufs-unbelegt` (3×, Stand `offen`; Zählung strittig, siehe seine `state.md`: die Datei `evidence/slice-release-v0-2-0.md` stammt von einem wellenlosen Slice, `slice-v1-abschluss-release-v0-1-1` aus dieser Welle trägt das zweite Auftreten).

**Klassen mit eigener Regel in `AGENTS.md`, in dieser Welle fortgeschrieben** (Zähler in Klammern, gesamt; die Regel gilt weiter, kein Absatz kommt hinzu, `AGENTS.md` §3.15): `zusage-im-kommentar-weiter-als-pruefung` (98; `slice-v1-abschluss-cancel-ohne-schluessel`, `slice-v1-abschluss-homebrew`, `slice-erster-release-homebrew-nachweis`, `slice-homebrew-nachweis-macos-arm64`), `plan-folgt-korrektur-nicht` (85; cancel-ohne-schluessel, release-v0-1-1, erster-release-homebrew-nachweis, macos-arm64), `negativtests-fehlen-bei-neuem-vertrag` (65; cancel-ohne-schluessel, homebrew), `spec-randform-erst-im-review-entschieden` (51; cancel-ohne-schluessel, homebrew, erster-release-homebrew-nachweis), `gruener-mutant-faelschlich-aequivalent` (11; cancel-ohne-schluessel).

**Rate nach `AGENTS.md` §3.15:** Seit `slice-harness-regeln-rueckbau` sind neun Slices geschlossen (`slice-v1-abschluss-verbindungsende-upstream-extended`, `slice-version-option`, `slice-play-feldvergleich`, `slice-play-feldvergleich-doku`, `slice-go-toolchain-1-27-3`, `slice-release-v0-2-0`, `slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`, `slice-homebrew-nachweis-macos-arm64`; Reihenfolge nach `git log` der Moves nach `done/`). Es sind weniger als zehn: **keine Rate und kein Vorschlag**; die Closure des zehnten Slice dahinter rechnet sie.

**Lerneintrag der Welle.** Geschärfte Regel: keine. Neuer Sensor: keiner (ein Sensor für einen Lauf auf einem Mac ist im Klon nicht möglich). Der Lerneintrag ist die benannte Spec-Lücke oben, mit einer Folgerung für künftige Wellen, die der Planner als Empfehlung führt (Entscheidung des Nutzers steht aus): Nennt ein Closure-Trigger eine Umgebung, die das Projekt nicht selbst fährt, nennt der Welle-Plan bei der Eröffnung, wer sie fährt und welche Stufe (Nachweis, Bericht, Grenze) den Trigger erfüllt. Hier fehlte das, und die Entscheidung fiel erst bei der Closure.

## Trigger-Audit

- Carveouts: `CO-001` (Gegenproben am getaggten Commit für `v0.1.1`/`v0.1.2`) bleibt aktiv: sein Auflösungs-Trigger (Tabelle in `SPEC-060` Punkt 16, alle Gegenproben `aufgenommen`, jeder Nehmer in `done/`) ist nicht eingetreten, die Nehmer `slice-harness-gegenproben-schneller` (`next/`) und `slice-harness-gegenproben-grenze-rest` (`open/`) liegen nicht in `done/`; verlängert wird er nicht, er deckt nur Tags. `CO-002` ist verbraucht und liegt in `carveouts/done/`. `CO-003` (Gegenproben am Wellenende, Entscheidung des Nutzers vom 2026-10-24, Wortlaut „CO-003 anlegen (Recommended)“) deckt das Wellenende (siehe unten); sein Trigger ist der von `CO-001` und nicht eingetreten. 2 offen (`CO-001`, `CO-003`), keiner fällig.
- Bootstrap-aware Gates: keine Reifestufe im Repo, 0 offen (`harness/README.md` §Sensors führt keine Stufe).
- Entscheidungen/ADRs: Die Re-Evaluierungs-Trigger der berührten Entscheidungen wurden in dieser Closure nicht erneut einzeln gelesen. Gelesen ist die Aussage der Slices dieser Welle: [ADR-0053](../../adr/0053-homebrew-tap-in-der-release-pipeline.md) (Homebrew-Tap in der Release-Pipeline) und [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) (Binaries ohne Version) trugen die Slices, und die Slices melden keinen eingetretenen Trigger; ein Wächter dafür ist `make baseline-verify` nicht. **Urteil:** nicht eingetreten laut Slices, ungeprüft in dieser Closure (Stichprobe entfällt, Grund: Frist des Auftrags); Nachprüfung beim Nutzer oder in der Closure der nächsten Welle.
- Hard Rules (`AGENTS.md` §3.2 bis §3.15): keine hat einen Auflösungs-Trigger, der eintrat; keine wurde in dieser Closure gelockert oder entfernt.
- MR: `MR-002` (permanent bis das Kurs-Repo einen Ort nennt) — Trigger nicht eingetreten.

## Gegenproben am Wellenende

`SPEC-060` Punkt 14 und `AGENTS.md` §6 verlangen am Ende jeder Welle `make gegenproben` und `make gegenproben-host`. Diese Closure hat sie **nicht gefahren**: Der Auftrag der Koordination schloss volle Läufe aus (Entscheidung des Nutzers vom 2026-10-22 und `AGENTS.md` §6: Zeitbudget; `CO-001` und `CO-002` decken allein Tags). Diese Closure ändert keinen Gate-Code, kein Skript und keinen Vertrag (nur Planungsdokumente und Archiv); die Slices der Welle, die Gegenproben änderten (`slice-v1-abschluss-binary-namen`, `slice-v1-abschluss-homebrew`), fuhren sie in ihrem eigenen §7 (in Teilen mit `GP_MUT`, nicht voll). Das ist eine Abweichung von `SPEC-060` Punkt 14, gedeckt durch den Carveout `CO-003` (Aktiv seit 2026-10-24, Entscheidung des Nutzers vom selben Tag). Ersatz nach `CO-003`: `make gates` auf `main` (siehe *Verifikation*) und die Angabe der geänderten Gegenproben je Slice (oben); **Grenze:** die Mutationen der Gegenproben (Wirksamkeit) und die nicht geänderten Gegenproben am Hash sind nicht gelaufen.

## Bestand in `open/` und `next/` (Vorschlag, Urteil des Nutzers)

32 Dateien: 27 in `open/`, 5 in `next/`. Gruppiert nach Gegenstand; **kein Einzelurteil neu gelesen**, jede Gruppe ist *bestätigt* (bleibt, Priorität unverändert) als Vorschlag; konsolidiert und entfallen: keiner. Kein Slice im Bestand gehört zu dieser Welle (§4 des Welle-Plans führt keinen).

- Gegenproben und Zeit: `slice-harness-gegenproben-schneller`, `slice-harness-gegenproben-grenze-rest`, `slice-harness-gegenprobe-baumkopie`, `slice-harness-gegenprobe-matrix-p7-flanke`, `slice-harness-randform-umgebung-der-gegenprobe`, `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb`, `slice-test-unter-last-ursache` — Träger von `CO-001` bzw. Register-Ausgang `geplant`.
- Rollenlauf-Zeitgrenze: `slice-harness-rollenlauf-zeitgrenze`, `-haertung`, `-regel`, `-wiederholungsgrenze`.
- Regeln und Register-Ausgänge: `slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`, `slice-harness-regeln-gate-konfiguration`, `slice-harness-regeln-implementer-uebergabe`, `slice-harness-commit-struktur-id`, `slice-harness-lint-warten-ohne-frist`, `slice-harness-mutation`, `slice-harness-kommentare-beleg-ort-und-traeger-adr`.
- Doku-Gate und Pins: `slice-harness-lh-links-bestand`, `slice-harness-lh-links-pflicht`, `slice-harness-ci-freshness`, `slice-harness-ci-freshness-images`, `slice-harness-ci-gegenproben`.
- Nachweis und Werkzeuge: `slice-harness-nachweis-uebernehmen`, `slice-harness-nachweis-uebernehmen-aufrufarten`, `slice-harness-slice-mv-exakter-treffer`, `slice-harness-gate-index-werkzeug-teil`, `slice-harness-dev-kopie-und-pin`, `slice-harness-container-start-u-boot-0-13`.
- Dev-Container: `slice-harness-devcontainer-einrichten-doku`, `slice-harness-devcontainer-einrichten-zugang`.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register.

Der Zähler steht in [`../observations/`](../observations/). Was in dieser Welle **3×** erreicht hat, steht oben unter *Steering-Loop-Einträge*.

## Folge-Slices

Keiner neu. `slice-homebrew-nachweis-macos-arm64` ist geschlossen und zugleich die Adresse des Register-Ausgangs `geplant`; ein Nachweis auf macOS (Lauf des Betreibers, Ausgabe als Beleg-Datei nach `SPEC-067` Punkt 20) wäre ein neuer Slice, den der Planner nur anlegt, wenn der Nutzer eine Ausgabe liefert. Die übrigen Lücken (macOS `amd64`, Linux `arm64`) haben keine Adresse und werden nicht erfunden.

Paarungen: (a) *Anker* — der Lerneintrag verkörpert keine Regel mit `liegt in`; entfällt. (b) *Folge-Slice* — keiner genannt; der Bestand oben ist keine Sendung. (c) *Register* — `homebrew-nachweis-ohne-macos-und-arm64`, `liste-gruener-mutanten-unvollstaendig` und `wiederanstoss-eines-roten-release-laufs-unbelegt` existieren als Verzeichnisse mit nicht leerem `evidence/` (Lauf des Planners 2026-10-24: `ls`).

## Archivierung

- `make archive-welle WELLE=welle-v1-nachtrag` (Träger `.harness/state/bin/ai-harness-init`) hat 13 Slice-Dateien und den Welle-Plan nach `done/welle-v1-nachtrag/archiv.zip` (226195 Bytes) archiviert, 14 Stubs daneben, 0 Review-Reports; zwei Commits (reiner Move, dann Archiv, Stubs und Verweis-Nachzug). Die Zeile `Bezug` der Stubs schrieb `make archiv-bezug WELLE=welle-v1-nachtrag` (13 Stubs, eigener Commit); weitere Verweise auf den Welle-Plan zog der Planner nach (zehn Dateien, ein eigener Commit).
- **Nicht archiviert:** `slice-play-feldvergleich-doku` bleibt als Volltext in `done/`: Sein Feld `Welle:` lautet „— (wellenlos), wie `slice-play-feldvergleich`“, das Werkzeug zählte ihn weder als Mitglied noch als wellenlos. Feststellung für den Nutzer, kein Eingriff in das Werkzeug; die nächste Closure sammelt ihn, wenn sein Feld es zulässt.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur, Schritt 1.

- Alle sieben Slices der Welle liegen in `done/` (Abzählung oben).
- Trigger *Abnahmeszenario 11 auf macOS und Linux nachgewiesen*: **nicht erfüllt**; geschlossen mit Grenze (Entscheidung des Nutzers vom 2026-10-24, Tabelle oben).
- Mutation je Zusage von Abnahmeszenario 11: Das Szenario ist ein Lauf mit Ausgabe, kein Test; die Zusagen belegt die Ausgabe des Laufs in `slice-erster-release-homebrew-nachweis` §7, den die Verifikation jenes Slice mit eigenem Lauf reproduzierte. Mutationen tragen das Verfahren der Formel (`slice-v1-abschluss-homebrew` §7, Fälle `p19-*` der Gegenprobe der Veröffentlichung, in Teilen mit `GP_MUT` gefahren) und den Nachweis-Vertrag (`SPEC-067` Punkt 20). Ein eigener Verifier-Lauf dieser Closure fand nicht statt.
- `make gates` auf `main` am Stand `1a9115cb` (nach dem Fast-Forward-Merge der Closure, vor diesem Nachtrag): Exit 0, 4 min 43 s; `make docs-check` im Worktree am selben Inhalt: 0 Befunde (1473 Dateien); `make kopf-check` Exit 0.
