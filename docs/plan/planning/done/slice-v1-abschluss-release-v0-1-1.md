# Slice slice-v1-abschluss-release-v0-1-1: Release v0.1.1 (erster Tag mit Binaries ohne Version im Namen)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-v1-nachtrag.

**Bezug:** [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität)

**Berührte Spec-Stellen:** `SPEC-067` (Punkte 2, 3, 8, 11, 17, 18 und 19 (2), (7), [`harness/vertraege/`](../../../../harness/vertraege/SPEC-067-release.md); gelesen, der Slice ändert den Vertrag nicht) · `SPEC-055` (Punkt 3, Namen der Binaries) · `LH-FA-19.a` · `LH-FA-16.a` · `LH-FA-01` · `LH-QA-03`

**Entscheidung:** [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) (Accepted, Entscheidung des Nutzers vom 2026-10-20) · [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md) · [ADR-0053](../../adr/0053-homebrew-tap-in-der-release-pipeline.md). Der Nutzer will den nächsten stabilen Tag als `v0.1.1` (Entscheidung des Nutzers vom 2026-10-20). **Entscheidung des Nutzers vom 2026-10-21** (wörtlich): „ich gebe alles frei, damit wir jetzt endlich diesen brew release machen können“ und „bitte auch den tag fahren“. Das ist das ausdrückliche Ja zum Tag `v0.1.1` und zu seiner Veröffentlichung (Punkt 7 der Checkliste); **den Tag setzt Claude (die Koordination) selbst**, weder der Nutzer noch ein Implementer-Lauf (`SPEC-067` Punkt 18: „wer taggen darf“, ein Implementer startet ihn nicht). Nicht eingeschlossen sind Secrets (Geheimnisse beider Registries, `HOMEBREW_TAP_GITHUB_TOKEN`), die der Nutzer selbst setzt (`releasing.md` §4); der Wortlaut nennt sie nicht, und der Slice behauptet nicht, dass sie gesetzt sind. **Entscheidung des Nutzers vom 2026-10-20** (Angabe der Koordination, im Plan festgehalten bei der Closure von `slice-v1-abschluss-binary-namen`): Das Ja des Nutzers zum Tag und zur Veröffentlichung von `v0.1.1` schließt einen Folge-Tag `v0.1.2` ein, falls der erste Lauf einen Fehler im Skript oder im Vertrag zeigt, ohne neue Frage; Secrets und andere Freigaben bleiben beim Nutzer, und jeder weitere Tag danach braucht ein neues Ja.

**Verantwortlich:** pt9912 (Zweig impl-rel)
**Autor:** Planner. **Datum:** 2026-10-20; nachgezogen 2026-10-21.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründung je Punkt.

**Ziel:** Der stabile Tag `v0.1.1` ist nach der Freigabe-Checkliste (`docs/maintainer/releasing.md` §10, `SPEC-067` Punkt 17) mit Beleg je Punkt vorbereitet; setzt Claude ihn nach dem ausdrücklichen Ja des Nutzers vom 2026-10-21, liegen am Release die sechs Binaries `pgwire-recorder-<os>-<arch>` ohne Version im Namen und `SHA256SUMS`, das Image in beiden Registries mit `0.1.1` und `latest`, und der Job `homebrew` hat die Formel im Tap versucht. Nach dem Lauf beschreiben Benutzerhandbuch und `README.md` den Bezug von `v0.1.1` mit den neuen Namen im Ist-Zustand. **Der Implementer dieses Slice setzt keinen Tag und veröffentlicht nichts:** Er bereitet vor (Release-Commit, Eintrag der Punkte 1 bis 6, Trockenlauf) und übergibt; den Tag setzt Claude (die Koordination) nach dem Ja des Nutzers vom 2026-10-21 (`SPEC-067` Punkt 18), danach trägt der Implementer die Belege aus dem Lauf ein.

**Übernommen aus `slice-v1-abschluss-binary-namen`** (dort §1 und §6, `AGENTS.md` §3.13): der Abschnitt „Binary vom GitHub-Release“ und der Abschnitt zum Image in §2 *Installation* des Benutzerhandbuchs und die Zeile zu `v0.1.0` in `README.md` werden auf `v0.1.1` und die neuen Namen umgeschrieben, belegt durch den Lauf des Tags (`SPEC-067` Punkt 18; `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`). Dieser Slice ist die Adresse; `grep` nach `slice-v1-abschluss-binary-namen` findet die Sendung in §1, DoD 2 und §8.

**Bindung für `slice-erster-release-homebrew-nachweis`** (Geber: dieser Slice, `AGENTS.md` §3.13): Der Nachweis setzt einen veröffentlichten stabilen Tag voraus, dessen Binaries `pgwire-recorder-<os>-<arch>` heißen und dessen Lauf den Job `homebrew` hatte; dieser Slice liefert den Tag (nach dem Ja des Nutzers vom 2026-10-21, gesetzt von Claude) und nennt im Eintrag der Freigabe das Ergebnis des Jobs `homebrew` (Lauf-Kennung, grün oder rot). Der Nachweis steht dort in §1, §4 und §6 mit der Kennung dieses Slice und als Start-Trigger.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Setzen des Tags und das Veröffentlichen durch den Implementer-Lauf — Handlung von Claude (der Koordination) nach dem Ja des Nutzers vom 2026-10-21 (Klasse 3, anderer Vorgang). Grund: `SPEC-067` Punkt 18 („ein Implementer startet ihn nicht“); der Plan ersetzt weder das Ja noch die Handlung.
- Die Namen der Binaries in den Skripten, in `make binaries` und im Gate — `slice-v1-abschluss-binary-namen` (Klasse 1; liegt vor diesem Slice, Start-Trigger). Grund: eigener Liefergegenstand mit eigener Gegenprobe.
- Der Job `homebrew`, die Vorlage der Formel und ihre Fälle — `slice-v1-abschluss-homebrew` (Klasse 1; liegt vor diesem Slice, Start-Trigger). Grund: Der Tag trägt nur dann eine Formel, wenn der Job im getaggten Commit steht.
- Die Formel im echten Tap, der Abschnitt zu Homebrew in §2 *Installation* des Handbuchs, die Zeile zu Homebrew in `README.md`, die Belege (Rechte des Tap-Tokens, Annahme durch Homebrew, Binary auf macOS und Linux) und der Wiederanlauf des Jobs oder der Nachtrag von Hand — `slice-erster-release-homebrew-nachweis` (Klasse 1). Grund: Das Handbuch nennt nur Befehle, die einmal gelaufen sind (`AGENTS.md` §3.11); der Lauf mit dem echten Tap ist dessen Gegenstand. Dieser Slice nennt den Ausgang des Jobs im Eintrag, mehr nicht.
- Eine Korrektur an `release.sh`, `release.yml`, den Verträgen oder den Gegenproben, die der Lauf des Tags verlangt — ein Folge-Slice (Klasse 3). Grund: Ein Lauf führt die Fassung des getaggten Commits aus; ein Fehler im Skript oder Vertrag verlangt einen neuen, höheren Tag (`SPEC-067` Punkt 18, *Rot*), für `v0.1.2` ohne neue Frage (Entscheidung des Nutzers vom 2026-10-20, Kopf), für jeden weiteren mit neuem Ja; dieser Slice liefert keinen Code (Klasse 4: kein Produkt-Code, kein Werkzeug-Code, nur `VERSION`, Doku und Planung).
- Der Release `v0.1.0` — unverändert (Bestand bleibt, Klasse 2). Grund: Ein Tag wird nie verändert (`SPEC-067` Punkt 11); er behält die Namen mit Version und bekommt keine Formel ([ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) Punkt 3). Das Handbuch nennt nach diesem Slice nur noch `v0.1.1` (§6).
- Das Freshness-Audit — optional, kein Punkt der Checkliste, sperrt den Tag nie (`SPEC-067` Punkt 17); ein durchgeführtes steht als eigener Absatz im Eintrag, einen Sensor liefern die Slices von `welle-ci-nachtlaeufe`.
- Änderung des Lastenhefts — es legt weder Namen noch Version fest.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität): **Freigabe und Tag.** Ein Release-Commit trägt `VERSION` mit `0.1.1` und die Zeile `0.1.1` der Änderungshistorie des Handbuchs (§11, Quelle der Release-Notizen, `releasing.md` §9). Der Eintrag der Freigabe (`docs/maintainer/freigaben/0.1.1.md`, eine Kopie der Tabelle aus `releasing.md` §10) trägt die Belege der Punkte 1 bis 6 für genau diesen Hash (`make gates`, Gegenproben-Container, `make gegenproben-host`, `grep`, Trockenlauf, Scan; Punkte 2 und 3 nach `CO-001`, siehe unten) und dazu das ausdrückliche Ja des Nutzers (Punkt 7: der Wortlaut aus dem Kopf, „ich gebe alles frei, damit wir jetzt endlich diesen brew release machen können“ und „bitte auch den tag fahren“, Datum 2026-10-21; kein weiteres Ja wird erfunden). Nach dem Tag, den Claude setzt, trägt er die Belege der Punkte 8 bis 10 aus dem Lauf: Lauf-Kennung, Index-Digest je Registry gleich, Smoke, `latest`, die sieben Anhänge mit den Namen `pgwire-recorder-<os>-<arch>` (Ausgabe von `gh release view`), die Beobachtungen nach `releasing.md` §10 und, getrennt, das Ergebnis des Jobs `homebrew` (Lauf-Kennung, grün oder rot; Belege des Taps selbst gehören dem Nachweis, §1). Was der Lauf nicht zeigt, steht als Grenze im Eintrag, nie als Behauptung (`AGENTS.md` §3.11).
- [x] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung): **Bezug nach dem Tag** (Punkt 11 der Checkliste, `SPEC-067` Punkt 18; Sendung aus `slice-v1-abschluss-binary-namen`, §1), in einem eigenen Commit nach dem Lauf, der nur `docs/user/`, `README.md`, `docs/maintainer/` und Planung ändert: Kopf des Handbuchs (*Software-Version*, *Gültigkeitsbereich*) nennt die Ausgabe von `version` des **veröffentlichten** Binaries; §2 *Installation* (Binary vom GitHub-Release, Image aus der Registry) nennt `v0.1.1`, die sechs Namen `pgwire-recorder-<os>-<arch>` ohne Version, die Prüfung mit `sha256sum -c` und die Tags `0.1.1` und `latest` mit dem Digest des Laufs — **nur mit Befehlen, die nach dem Tag einmal gelaufen sind** (Beleg im Eintrag: Befehl, Lauf, Ausgabe; die Beispiele liegen als Datei vor und laufen gegen das veröffentlichte Binary, `AGENTS.md` §3.11); `README.md` nennt die Version `0.1.1` mit den Adressen im Ist-Zustand; `docs/maintainer/releasing.md` §1 (*Stand*) beschreibt den Zustand nach dem Tag (Tag gesetzt, Läufe, Formel nach dem Ausgang des Jobs) im Ist-Zustand ohne Chronik. Der Abschnitt zu Homebrew bleibt dem Nachweis (§1).
- [x] `make gates` grün; die Gegenproben am getaggten Commit sind Punkte 2 und 3 des Eintrags: nicht voll gefahren, sondern mit dem Carveout [`CO-001`](../../carveouts/CO-001-gegenproben-am-getaggten-commit.md) (Aktiv seit 2026-10-21) und den Ersatzbelegen aus dessen Text (Ablage des Gegenproben-Containers, unveränderte Kopie von `release-gegenprobe` und `test-container-gegenprobe` am getaggten Hash, Diff der Dateien von `container-start` für Punkt 3; Randform (m)); der Eintrag schreibt je Punkt den Wortlaut aus `CO-001` „Was der Eintrag schreibt“ mit Befehl, Hash, Datum, Dauer, Last und Ausgang; nie „grün“ ohne Beleg. Eine volle Fahrt ist erlaubt, wo sie in Minuten endet, aber nicht verlangt.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `VERSION` | update | `0.1.1` im Release-Commit, der getaggt wird; ändert das Wurzelverzeichnis und ist damit nicht inert nach `SPEC-052` Punkt 12 (volle Kette) |
| `docs/user/benutzerhandbuch.md` §11 | update | Zeile `0.1.1` der Änderungshistorie im Release-Commit, vor dem Tag |
| `docs/maintainer/freigaben/0.1.1.md` | neu | Eintrag der Freigabe, Kopie der Tabelle aus `releasing.md` §10, Belege je Punkt |
| `docs/maintainer/releasing.md` §10 | update | Ort des Eintrags „für jeden späteren Release“, Zeile des Jobs `homebrew` unter Punkt 8 (Spiegel von `SPEC-067` Punkt 17) und ein Satz unter der Tabelle, dass ein Punkt, den ein aktiver Carveout ausnimmt, im Beleg Carveout, Ersatz und Grenze nennt (Form: `CO-001` „Was der Eintrag schreibt“; Spiegel, der Vertrag bleibt beim Architect), im Release-Commit oder davor |
| `docs/user/benutzerhandbuch.md` (Kopf, §2 *Installation*), `README.md` | update | Bezug nach dem Tag (`SPEC-067` Punkt 18): `v0.1.1`, neue Namen, Image-Tags, nur belegte Befehle |
| `docs/maintainer/releasing.md` §1 | update | *Stand* nach dem Tag |
| Lauf `release.yml` auf dem Tag (Trockenlauf davor) | Handlung | kein Dateiinhalt; der Tag ist Handlung von Claude nach dem Ja des Nutzers, nicht des Implementers; der Trockenlauf veröffentlicht nichts (`SPEC-067` Punkt 2) |

Nicht in dieser Liste: `tools/`, `.github/`, `Dockerfile`, `cmd/`, `internal/` und die Verträge (§1), der Abschnitt zu Homebrew in Handbuch und `README.md` (§1, Nachweis).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-binary-namen` und `slice-v1-abschluss-homebrew` liegen in `done/` (die Namen im getaggten Commit und der Job `homebrew` in `release.yml`; Reihenfolge in §5 von [welle-v1-nachtrag](../welle-v1-nachtrag.md)), [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) ist `Accepted` (Entscheidung des Nutzers vom 2026-10-20), und die Vorbedingungen des Betreibers stehen (`releasing.md` §4): die Geheimnisse für beide Registries und `HOMEBREW_TAP_GITHUB_TOKEN` sind gesetzt, das Repository des Taps `pt9912/homebrew-pgwire-recorder` besteht (Handlung des Betreibers; der Slice prüft, legt nichts an und liest keinen Wert). Das Ja zum Tag ist erteilt (Kopf, 2026-10-21) und **kein Start-Trigger**, sondern ein Schritt im Slice (DoD 1): Der Implementer arbeitet am Release-Commit und an den Belegen 1 bis 6 und übergibt; den Tag setzt Claude.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe aus Einfügungen und Löschungen von `git diff --numstat 0256ddff..<letzter Commit des Slice> -- VERSION docs/user docs/maintainer README.md` (Tests eingeschlossen, `AGENTS.md` §3.14; der Slice liefert keinen Test- oder Werkzeug-Code) **über 500** liegt. Schätzung des Planners in dieser Metrik, ohne Messung: Handbuch (Abschnitte zu Binary und Image umgeschrieben, ersetzte Zeilen zählen doppelt) rund 100, Kopf, §11 und `README.md` rund 20, `releasing.md` §1 rund 40, Eintrag der Freigabe rund 90 bis 140, zusammen 250 bis 300. Kein vorgemerkter Schnitt nötig; fiele der Slice darüber, wäre die Naht (a) Release-Commit, Eintrag bis Punkt 7 und Tag, (b) Bezug nach dem Tag.
- `in-progress` → `open`: Das Ja des Nutzers ist erteilt (Kopf); zurück geht der Slice, wenn Claude den Tag nicht setzt (etwa weil das Ja nicht mehr trägt) oder ein vorbereitender Beleg (Punkte 1 bis 6) wird rot und die Ursache liegt außerhalb dieses Slice (etwa ein roter Scan wegen eines Befunds ohne Behebung: Entscheidung des Nutzers zur Schwelle oder Carveout, kein Tag; `SPEC-067` Punkt 16). Ein Slice, der wartet, hält das WIP-Limit; deshalb geht er zurück nach `open/`, statt zu warten.
- `in-progress` → `next`: Der Lauf des Tags endet rot aus einem Grund im Skript, im Text der Datei oder im Vertrag (`SPEC-067` Punkt 18 (b)): Der Tag bleibt, die Korrektur ist ein Folge-Slice, ein neuer höherer Tag `v0.1.2` ist durch das Ja zu `v0.1.1` gedeckt (Entscheidung des Nutzers vom 2026-10-20, Kopf), jeder weitere verlangt ein neues Ja; dieser Slice schließt dann mit dem Eintrag für `v0.1.1` (Befund, Ersatz-Tag) nach der Incident-Klausel.

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

**Randformen** (`AGENTS.md` §3.12). Der Slice liefert keinen neuen Vertrag, sondern führt `SPEC-067` Punkt 17 und 18 ein zweites Mal aus; die Randformen stehen dort oder sind als Rückgabe benannt.

| Randform | Entscheidung | Stelle |
|---|---|---|
| (a) Version | `0.1.1`, ein stabiler Tag (kein `-`): `latest` wird gesetzt, der Job `homebrew` läuft; Tag `v0.1.1` gleich `v` plus Inhalt von `VERSION`, sonst endet `Eingaben` mit Exit 2 | Entscheidung des Nutzers vom 2026-10-20; `SPEC-067` Punkte 3, 8, 19 (2) |
| (b) Wer taggt, wann | **Claude (die Koordination), nach dem Ja des Nutzers vom 2026-10-21** („bitte auch den tag fahren“); der Implementer-Lauf setzt keinen Tag und startet keinen (`SPEC-067` Punkt 18); Übergabe: Implementer übergibt nach Belegen 1 bis 6 mit dem Hash des Release-Commits, Claude tagt genau diesen Hash | `SPEC-067` Punkt 18; `releasing.md` §10 |
| (c) Reihenfolge der Commits | (1) Release-Commit mit `VERSION` und der Zeile `0.1.1` (Zeile **vor** dem Tag, sie ist die Quelle der Notizen), (2) Tag auf diesem Commit, (3) Bezug-Commit nach dem Lauf; die Punkte 1 bis 3 der Checkliste gelten für den Hash von (1) | `SPEC-067` Punkte 17, 18; `releasing.md` §10 |
| (d) Rot nach dem Tag | Umgebung (Token, Schreibrecht, Sichtbarkeit, Netz): Wiederanstoß desselben Tags; Skript, Text der Datei oder Vertrag: Korrektur auf dem Hauptzweig und neuer höherer Tag mit neuem Ja; Digest oder Anhang abweichend: immer neuer Tag. **Das Ja zum ersten Tag schloss `v0.1.1` ein; das Ja zu `v0.1.1` schließt `v0.1.2` ein (Entscheidung des Nutzers vom 2026-10-20; das Ja vom 2026-10-21 ändert das nicht)** (Entscheidung des Nutzers vom 2026-10-20, Kopf) — jeder weitere Tag danach braucht ein neues Ja | `SPEC-067` Punkte 11, 18; Entscheidung des Nutzers im Kopf |
| (e) Job `homebrew` rot oder grün | Ein roter Job ändert weder Release noch Registries; der Slice nennt sein Ergebnis im Eintrag und gibt Wiederanlauf oder Nachtrag von Hand an den Nachweis (kein neuer Tag); Belege des Taps sind nicht Teil dieses Slice | `SPEC-067` Punkt 19 (1), (7); Nachweis |
| (f) `v0.1.0` und das alte Schema | `v0.1.0` bleibt mit Namen mit Version; das Handbuch nennt nach dem Bezug-Commit nur `v0.1.1` und die neuen Namen und sagt nichts über die alten (Ist-Zustand, `AGENTS.md` §3.11); der Schritt `homebrew` lehnt das alte Schema ab | `SPEC-055` Punkt 3; `SPEC-067` Punkt 19 (3); Rückgabe an den Architect unten |
| (g) `latest` | wandert von `0.1.0` auf `0.1.1`, wenn Version stabil, Vergleich der Digests beider Registries grün und `0.1.1` höher als der vorhandene `latest`; der Tag `0.1.0` des Images bleibt | `SPEC-067` Punkt 8 |
| (h) Plattform | Linux (`amd64`) im Lauf geprüft: Binary gestartet, Image gestartet; `darwin`, `windows` und `linux/arm64` gebaut, nicht ausgeführt: Grenze, die das Handbuch so nennt | `SPEC-067` Punkt 15; `SPEC-055` Punkt 13 (4) |
| (i) Aufrufarten, Abbruch | entfällt: Der Slice liefert kein Werkzeug und kein Ergebnis, das ein Hook liest; Treiber der Pipeline ist `release.sh` (nie unter `make`) | `SPEC-067` Punkt 13 |
| (j) Ort des Eintrags der Freigabe | **entschieden (Architect, 2026-10-21): `docs/maintainer/freigaben/0.1.1.md`**; der Vertrag nennt `freigaben/<Version>.md` jetzt für jeden späteren Release, mit oder ohne Welle | `SPEC-067` Punkt 17 |
| (k) Job `homebrew` in der Checkliste | **entschieden: Teil von Punkt 8**, getrennt darunter, mit Lauf-Kennung, grün oder rot, bei grün Hash des Commits im Tap (`git ls-remote`, Datum), bei rot Stichwort der Meldung; ein roter Job macht Punkt 8 nicht unerfüllt; Mindestmenge bleibt elf, keine Neunummerierung (`releasing.md` §10 zieht die Zeile von Punkt 8 mit, Ort der Doku: dieser Slice) | `SPEC-067` Punkt 17 |
| (l) Bezug des zweiten Tags | **entschieden: dieselbe Form wie beim ersten Tag; der Bezug ersetzt den des vorigen Tags und schweigt über frühere Namen** (Randform (f) bestätigt) | `SPEC-067` Punkt 18 (Absatz „Jeder weitere Tag“) |
| (m) Punkte 2 und 3 der Checkliste (Gegenproben am getaggten Commit) | **entschieden mit der Freigabe des Nutzers vom 2026-10-21** (Wortlaut im Kopf; der Entwurf lag ihm vor): Die Pflicht bleibt im Vertrag; für `v0.1.1` und `v0.1.2` ruht sie über den Carveout `CO-001` (Aktiv seit 2026-10-21, `docs/plan/carveouts/`), mit Ersatzbelegen in Minuten; Trigger und Folge-Slices dort. Offen bleibt allein, ob der Nutzer einen dritten Ersatzbeleg verlangt (`CO-001` „Was der Nutzer noch entscheiden kann“) | `docs/plan/carveouts/CO-001-gegenproben-am-getaggten-commit.md`; [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md); `AGENTS.md` §3.6 |

**Rückgaben an den Architect:** beantwortet am 2026-10-21 — (1) Ort: Randform (j); (2) Job `homebrew`: Randform (k); (3) zweiter Tag: Randform (l). Neu entschieden: (m).

**Entscheidung des Nutzers vom 2026-10-20 zu Randform (d)** (vorher als Frage geführt): Das Ja zu `v0.1.1` schließt einen Folge-Tag `v0.1.2` ein, falls der Lauf von `v0.1.1` einen Fehler im Skript oder im Vertrag zeigt; kein weiteres Ja dafür. Nicht eingeschlossen: Secrets, andere Freigaben, jeder Tag nach `v0.1.2`. Das Ja zum Tag `v0.1.1` selbst hat der Nutzer am 2026-10-21 gegeben (Kopf, Wortlaut); den Tag setzt Claude.

**Risiken:**

- Der Lauf des Tags ist die erste Außenwirkung des Jobs `homebrew` und des neuen Schemas auf GitHub (Rechte des Tap-Tokens, Repository des Taps, Wortlaut von `gh release download`, Annahme durch Homebrew); der Trockenlauf führt den Job nicht (`SPEC-067` Punkt 19 (1)). Rot dort ist kein Grund für einen neuen Tag (Randform (e)) — **Ausgang:** *eingetreten*, Folge-Slice `slice-erster-release-homebrew-nachweis`: Die Versuche 1 bis 3 des Jobs waren rot (`Permission … denied`, `error: 403`), grün wurde er im vierten nach zwei Eingriffen des Betreibers; was offen blieb (Annahme der Formel durch Homebrew, `brew tap` und `brew install`, Wortlaut von `gh release download`, Schreibrecht ohne Eingriff), steht im Nehmer in §1 und §6 mit der Kennung dieses Slice, eingetragen mit Nachzählung in dessen §8 (`AGENTS.md` §3.13). Nicht *weiter offen*: Das Risiko ist eingetreten und hat eine Adresse, der Zähler ist dafür nicht der Ort (Begründung in §7).
- Ein Fehler im Skript, der erst der Lauf des Tags zeigt, kostet einen Folge-Tag `v0.1.2`, den das Ja zu `v0.1.1` deckt (Entscheidung des Nutzers vom 2026-10-20), ein weiterer brauchte ein neues Ja — **Ausgang:** *entfallen*: Der Job `release` war im ersten Versuch grün und wurde in die Versuche 2 bis 4 übernommen; die Ursache der roten Versuche lag in der Umgebung (`SPEC-067` Punkt 18 (a)), nicht im Skript, im Text der Datei oder im Vertrag. Ein Ersatz-Tag `v0.1.2` war nicht nötig, `v0.1.1` bleibt gültig.
- Die Fundstellen des alten Schemas im Handbuch und in `README.md` (`pgwire-recorder_` mit Version) sind der Gegenstand des Bezug-Commits — **Ausgang:** *entfallen*: `grep -rn 'pgwire-recorder_' docs/user README.md` am Endstand ohne Treffer (Exit 1), von der Verifikation nachgefahren; `0.1.0` steht nur noch in der Zeile der Änderungshistorie, und die bleibt (`SPEC-067` Punkt 18).
- `docs/maintainer/releasing.md` §1 (*Stand*) beschreibt den Zustand von vor dem ersten Tag („Es ist noch kein Tag gesetzt worden“, „Eine Formel im Tap gibt es nicht“); der Planner fand das beim Lesen, ob der Stand an anderer Stelle nachgezogen wird, ist nicht geprüft — **Ausgang:** *entfallen*: Der Bezug-Commit `9879253a` hat §1 nachgezogen (die Sätze „Es ist noch kein Tag gesetzt worden“ und die drei „hatte noch keinen Lauf“ sind weg, an ihrer Stelle steht der Zustand nach `v0.1.1`); die Verifikation hat den Abschnitt gelesen.

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

### Closure (Planner, 2026-10-21)

**Zwei beobachtbare Kriterien.** (1) Der Tag `v0.1.1` liegt auf dem Release-Commit `3594b9f0`, der Lauf `39551891874` ist `success`, und der Eintrag `docs/maintainer/freigaben/0.1.1.md` trägt alle elf Punkte der Checkliste samt der Zeile des Jobs `homebrew`; die Verifikation vom 2026-10-21 hat jede überprüfbare Zahl, Zeile und Kennung am Original nachgefahren (`docs/reviews/2026-10-21-verifikation-slice-v1-abschluss-release-v0-1-1.md`). (2) `make gates` grün (Messung des Implementers; die Dateiliste des vollen Laufs deckt den Endstand bis auf den Plan-Pfad), `make docs-check` grün. Die Punkte 2 und 3 stehen nicht „grün“, sondern mit dem Carveout [`CO-001`](../../carveouts/CO-001-gegenproben-am-getaggten-commit.md) und seinen Ersatzbelegen — ein rotes Gate ist das nicht, eine nicht gefahrene Kette mit dokumentierter Ausnahme schon.

**Befunde und ihr Verbleib.** Review `ccc23dac` (1 HIGH, 4 MEDIUM, 2 LOW, 2 INFO): F-1 bis F-7 nachgearbeitet in `0c8a13d3`, von der Verifikation einzeln gegengeprüft; F-8 und F-9 sind Vertragswortlaut und liegen beim Architect (`SPEC-067` Punkt 17 (5) gegen Punkt 3 (5); Punkt 10 der Checkliste in `releasing.md` nennt den ersten Lauf, wo jeder Lauf gilt). Verifikation `e583cde2` (0 HIGH, 1 MEDIUM, 3 LOW, 3 INFO): F-V1, F-V3, F-V5 und F-V7 im Eintrag und F-V1 in §7 berichtigt (Commit `0ee9f164`), F-V2 und F-V4 unten in dieser Sektion, F-V6 im Nehmer (Commit `9318c4b7`). Der Grund für den fehlenden Ersatzbeleg 1 nannte das falsche Repository: Gemeint ist die Ablage des Gegenproben-Containers `pgr-gp` (`SPEC-061`, per `ssh`); den Container gibt es nicht mehr (`docker ps -a` nennt allein `pgr-dev2`), und `impl-rel` wurde nie dorthin eingeliefert — die Lücke ist real, ihre Begründung war es nicht.

**Lerneintrag — benannte Spec-Lücke** (`SPEC-067` Punkt 17 und `docs/maintainer/releasing.md` §4): Von den Geheimnissen des Betreibers prüft die Freigabe-Checkliste die **Existenz**, nicht das **Recht**, und nennt das nicht als Grenze. Der Trockenlauf kann es nicht zeigen, denn er führt den Job `homebrew` nicht aus (`SPEC-067` Punkt 19 (1)). Beide Vorbedingungen galten damit als erfüllt, und erst der Lauf des gesetzten Tags widerlegte eine: drei rote Versuche mit `error: 403`, zwei Eingriffe des Betreibers, der zweite erst wirksam. Die Lücke ist benannt, nicht geschlossen — der Wortlaut eines Vertrags gehört dem Architect, die Entscheidung dem Nutzer; der Vorschlag (Grenze bei Punkt 7 der Checkliste, je Schritt mit seiner Folge; teurere Alternative: ein schreibender Probeaufruf je Geheimnis im Trockenlauf, mit Außenwirkung) steht im Register unter `BEO-REPO/vorbedingung-geheimnis-nur-existenz-geprueft`. Kein `liegt in`: Dieser Slice verkörpert keine Regel.

**Risiko-Ausgänge (§6), vier Risiken, je einer:**

1. Erste Außenwirkung des Jobs `homebrew` und des neuen Schemas — **eingetreten**, Folge-Slice `slice-erster-release-homebrew-nachweis`. Das Risiko hat sich gezeigt (Versuche 1 bis 3 rot, `Permission … denied`, `error: 403`), und was es offen ließ — Annahme der Formel durch Homebrew, `brew tap` und `brew install` auf macOS und Linux, Wortlaut von `gh release download`, Schreibrecht ohne Eingriff des Betreibers — führt der Nehmer in §1 und §6 mit der Kennung dieses Slice (eingetragen mit Nachzählung, Commit `9318c4b7`, `AGENTS.md` §3.13). §6 hatte „weiter offen“ vorformuliert; der Ausgang *weiter offen* hängt ein Risiko an den Zähler, und das wäre falsch für eines, das eingetreten ist und eine Adresse hat (Baseline-Regelwerk `modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst).
2. Fehler im Skript mit Folge-Tag `v0.1.2` — **entfallen**. Der Job `release` war im ersten Versuch grün und wurde in die Versuche 2 bis 4 übernommen; die Ursache der roten Versuche lag in der Umgebung (`SPEC-067` Punkt 18 (a)), nicht im Skript, im Text der Datei oder im Vertrag. Ein Ersatz-Tag war nicht nötig, `v0.1.1` bleibt der gültige Tag. Das Ja zu `v0.1.2` bleibt ungenutzt und verfällt nicht (Entscheidung des Nutzers vom 2026-10-20, Kopf).
3. Fundstellen des alten Schemas in Handbuch und `README.md` — **entfallen**. `grep -rn 'pgwire-recorder_' docs/user README.md` am Endstand ohne Treffer (Exit 1), von der Verifikation nachgefahren; `0.1.0` steht nur noch in der Zeile der Änderungshistorie, und die bleibt (`SPEC-067` Punkt 18).
4. `docs/maintainer/releasing.md` §1 (*Stand*) auf dem Stand von vor dem ersten Tag — **entfallen**. Der Bezug-Commit `9879253a` hat ihn nachgezogen; die Sätze „Es ist noch kein Tag gesetzt worden“ und die drei „hatte noch keinen Lauf“ sind weg, an ihrer Stelle steht der Zustand nach `v0.1.1`. Die Verifikation hat den Abschnitt gelesen.

**Die drei Paarungen.**

- **Anker:** keiner, und das ist kein Versäumnis — der Slice verkörpert keine Regel (kein Diff an `AGENTS.md`, an einem Gate oder an einem Command), also steht kein `liegt in` und kein `seit slice-v1-abschluss-release-v0-1-1` irgendwo. Der Lerneintrag ist eine benannte Spec-Lücke.
- **Folge-Slice:** `slice-erster-release-homebrew-nachweis`, in `next/` und nicht geschlossen; er nimmt die Sendung in §1 („Wiederanlauf des Jobs, Nachtrag von Hand und alle Belege des Taps gehören diesem Slice“) und in §6 ausdrücklich an und nennt dort die Kennung dieses Slice — `grep -n slice-v1-abschluss-release-v0-1-1 docs/plan/planning/in-progress/slice-erster-release-homebrew-nachweis.md` findet sie in §1, §4, §6 und §8. Nachgezählt im selben Commit: Liefer-Punkte 2, Schichten zwei.
- **Register:** vier Einträge berührt (unten).

**Beobachtungs-Register** (`../observations/BEO-REPO/`):

| Eintrag | Was dieser Slice beiträgt |
|---|---|
| `wiederanstoss-eines-roten-release-laufs-unbelegt` | **Zweites Auftreten** (Beleg-Datei angelegt, Zähler 2, unter der Schwelle): Der Job `release` endete im ersten Versuch grün, der Wiederanstoß eines roten **Laufs** blieb wieder ungegangen. Die andere Hälfte ist seit diesem Lauf belegt — ein roter **Job** wurde dreimal wieder angestoßen, zweimal ohne Änderung an der Ursache; `state.md` grenzt die Beobachtung darauf ein. |
| `plan-folgt-korrektur-nicht` | **Wiederaufgetreten** (Zähler 76): Review F-2 (Tap-Belege im Eintrag, während §1 und §6 sie dem Nachweis zuwiesen) und die Verifikations-Befunde F-V2 und F-V4 (zwei Zahlen in §7 hinter dem Stand). Die Regel `AGENTS.md` §3.9 bleibt. |
| `vorbedingung-geheimnis-nur-existenz-geprueft` | **Neu, erstes Auftreten** — die Lücke des Lerneintrags; Abgrenzung zu `vertrag-nennt-aufruf-eines-fremdwerkzeugs-ungemessen` (dort eine messbare Annahme des Vertrags über ein Fremdwerkzeug) und zu `ci-ohne-branch-schutz` (dort sagt das Repo die fehlende Einstellung selbst) steht in `observation.md`. Ohne Adresse. |
| `test-faellt-unter-cpu-last-aus` | **Kein Auftreten, Gegenbeleg**, keine Beleg-Datei, Zähler bleibt 4: Drei Läufe von `make gates` in diesem Slice endeten im ersten Versuch grün, der eine bei Last 13,7 auf 20 Kernen; `TestE2EPlayVergleichAbweichung` fiel in keinem aus. Die Beobachtung ist damit nicht widerlegt, aber die Last allein löst sie nicht aus. |

**Was dieser Slice nicht entschieden hat** (für den Lese-Schritt der Welle-Closure): der Ausgang von `test-faellt-unter-cpu-last-aus` (4×, über der Schwelle, Zuweisung beim Nutzer), ein dritter Ersatzbeleg zu `CO-001` („Was der Nutzer noch entscheiden kann“) und der Vertragswortlaut zu F-8 und F-9 (Architect).

**Läufe bei der Closure** (Zweig `clo-rel`, Worktree `wt-clo2`, Dev-Container `pgr-dev2`, 2026-10-21): `make docs-check` grün (1322 Dateien, 0 Befunde) und `make kopf-check` grün (Exit 0, keine Ausgabe), gefahren im Arbeitsbaum mit allen Dateien dieser Closure, vor dem Commit; der Baum trug dabei auch den neuen Plan `slice-version-option` in `open/`, der in einem eigenen Commit liegt und nicht zu diesem Slice gehört. Die Zahl liegt über den 1315 am Stand `0c8a13d3`, weil fünf Register-Dateien und dieser Plan hinzukommen. Kein `make gates` und keine Gegenprobe (Entscheidung des Nutzers vom 2026-10-21; die Closure ändert nur `docs/plan/planning/` und `docs/maintainer/freigaben/0.1.1.md`, dieser allein in Zeilen des Eintrags, nicht in `VERSION`, Code, Gates oder Verträgen).

### Belege des Implementers (Vorbereitung, 2026-10-21, Zweig impl-rel)

Der Slice liefert keinen Code und keinen Vertrag (§1); die Zusagen sind die Punkte der Checkliste. Der Hash des Tag-Kandidaten ist der Release-Commit `3594b9f0578dd8c844e9a9d3e5bc6f147bcd3aaa`; der Eintrag `docs/maintainer/freigaben/0.1.1.md` steht in einem späteren Commit (er kann den Hash seines eigenen Commits nicht nennen) und nennt die Belege je Punkt mit Befehl, Dauer und Last.

| Zusage | Beleg | Gegenbeispiel (rot gesehen) |
|---|---|---|
| `VERSION` ist `0.1.1` im getaggten Commit (Punkt 4) | `git show 3594b9f0:VERSION` gibt `0.1.1` | `release.sh eingaben` mit `VERSION` `0.1.2` und Zeile `0.1.1` meldet `keine Zeile für 0.1.2` (Hinweis im Trockenlauf; im Lauf des Tags wäre es Exit 2 `tag-ungleich`, nur über die Gegenprobe belegt, nicht hier gefahren) |
| Die Zeile `0.1.1` der Änderungshistorie steht im getaggten Commit mit nicht leerer dritter Zelle (Punkt 4, `releasing.md` §9) | Zeile 2174 des Handbuchs; `release.sh eingaben` im Trockenlauf-Modus gibt `release: hinweis: Änderungshistorie 0.1.1: …` und `eingaben: grün für v0.1.1` | dieselbe Probe mit entfernter Zeile gibt `die Änderungshistorie des Benutzerhandbuchs hat keine Zeile für 0.1.1` statt des Textes |
| `make gates` grün am Hash (Punkt 1) | Exit 0, 287 s, Last 13,7 | nicht gefahren: kein Mutant (Gate des Bestands, §1) |
| Punkt 2 nach `CO-001`: unveränderte Kopie von `release.yml` und `release.sh` grün | `GP_MUT='^ok$' make release-gegenprobe` Exit 0, 142 s, `grün (1 Läufe)` | nicht gefahren (Muster `^ok$` lässt nur die unveränderte Kopie laufen; Grenze steht im Eintrag) |
| Punkt 2 nach `CO-001`: unveränderte Kopie von `binaries` und `test-container` grün | `GP_MUT='^kein-treffer$' make test-container-gegenprobe` Exit 0, 22 s | nicht gefahren, Grenze im Eintrag |
| Punkt 3 nach `CO-001`: Dateien von `container-start` seit der letzten grünen Fahrt unverändert | `git diff --stat 35778a2f 3594b9f0 -- <fünf Dateien>` leer | ein Diff ungleich leer wäre die Grenze und der Halt; die letzte grüne Fahrt ist nur eine Angabe der Koordination (Eintrag von `welle-erster-release`), im Eintrag als Grenze genannt |
| Punkte 5 und 6 (Trockenlauf, Scan) | **nicht erfüllt**: `origin` des Worktrees ist der lokale Austausch, `gh workflow run` wurde dem Lauf verweigert; die Koordination fährt ihn | Teilbeleg lokal nur für den Schritt `eingaben` (siehe oben) |

**Nachtrag nach dem Tag (2026-10-21).** Der Halt ist aufgelöst: Die Koordination fuhr den Trockenlauf (`39550948115`, `workflow_dispatch` auf `main` = `762c05dd`, success, 2 min 50 s) und setzte den Tag `v0.1.1` auf `3594b9f0`; der Lauf des Tags ist `39551891874`. Die Punkte 5, 6 und 8 bis 11 des Eintrags tragen seitdem ihre Belege; der Implementer hat die Zeilen des Logs mit `gh run view <Kennung> --log` selbst gelesen, nicht übernommen. Der Job `homebrew` lief in derselben Ausführung viermal: die Versuche 1 bis 3 rot mit `tap-push` (05:55:27Z, 06:10:21Z, 06:12:45Z), der vierte grün (06:14:14Z); die conclusion des Laufs ist seitdem `success`. Die Zeitstempel je Versuch hat der Implementer mit `gh api …/attempts/<n>/jobs` selbst gelesen.

| Zusage | Beleg | Gegenbeispiel (rot gesehen) |
|---|---|---|
| Der Lauf des Tags stand auf dem Release-Commit | `gh run view 39551891874 --json headSha` gibt `3594b9f0578dd8c844e9a9d3e5bc6f147bcd3aaa`, Auslöser `push` von `refs/tags/v0.1.1` | nicht gefahren: der Hash ist eine Ablesung, kein Wächter |
| Index-Digest je Registry gleich (Punkt 8) | `release: vergleichen: sha256:624c09… in ghcr.io und docker.io`; dazu vier eigene `docker pull` (`:0.1.1` und `:latest`, beide Registries, leeres `DOCKER_CONFIG`), alle vier melden `sha256:624c099b…9b15` | der Schritt `vergleichen` ist fail-closed (`SPEC-067` Punkt 7); seine Mutanten fährt `make release-gegenprobe`, hier nicht erneut |
| Der GitHub-Release trägt genau sieben Anhänge mit den Namen ohne Version (Punkt 9) | `gh release view v0.1.1 --json assets` nennt die sechs `pgwire-recorder-<os>-<arch>` und `SHA256SUMS`, sonst nichts | `SHA256SUMS` des Release gegen die geladene Datei: `sha256sum -c` hätte eine Abweichung als `FAILED` gemeldet; gemeldet wurde `OK` |
| Job `homebrew` grün, im vierten Versuch (Zeile unter Punkt 8) | `release: homebrew: Formel 0.1.1 im Tap abgelegt` (Versuch 4, 06:14:20Z); Tap-Commit `6c4ac72843996939b4295b01aeb7cbe9e759a7ac` als Spitze von `main`, gelesen mit `git ls-remote https://github.com/pt9912/homebrew-pgwire-recorder.git refs/heads/main` (die Form, die `SPEC-067` Punkt 17 (8) nennt) | die Versuche 1 bis 3 desselben Jobs sind rot mit `tap-push: der Tap weist den Push ab (Exit 128)` — das ist die rot gesehene Gegenprobe zur Zusage „der Job legt die Formel im Tap ab“, und sie zeigt zugleich, dass ein Wiederanlauf ohne Änderung an der Umgebung dieselbe Meldung ergibt (Versuch 2 nach der ersten Berichtigung) |
| Handbuch §2: `sha256sum -c --ignore-missing SHA256SUMS` gibt `pgwire-recorder-linux-amd64: OK` | gelaufen gegen die veröffentlichte Datei, Exit 0 | ohne `--ignore-missing` endet derselbe Aufruf mit Exit 1 und `WARNING: 5 listed files could not be read` — der Satz des Handbuchs über „alle sechs Dateien“ ist damit belegt, nicht nur behauptet |
| Handbuch §2: `./pgwire-recorder-linux-amd64 version` gibt `pgwire-recorder 0.1.1` | gelaufen, Ausgabe wörtlich | ein Binary mit anderer Version gäbe eine andere Zeile; das Binary von `make build` meldet `dev` (Satz im Kopf) |
| Handbuch §2: beide Registries geben das Image ohne Anmeldung heraus | vier `docker pull` mit `DOCKER_CONFIG` auf ein leeres Verzeichnis, alle vier Exit 0 | mit Anmeldung liefe es ebenfalls; die Sonde ist genau das leere `DOCKER_CONFIG`, ohne das der Satz eine Annahme bliebe |
| Handbuch §2: `latest` zeigt in beiden Registries auf denselben Digest und gibt dieselbe Ausgabe | `docker image inspect` für `:latest` beider Registries gibt `sha256:624c099b…9b15`; `docker run … :latest version` gibt beide Male `pgwire-recorder 0.1.1` | ein `latest`, das auf `0.1.0` stünde, gäbe `pgwire-recorder 0.1.0`; der Schritt `latest` lief laut Log für beide Registries |
| Handbuch §2: gestartet ist allein `linux/amd64` | `uname -m` gibt `x86_64`, `docker image inspect --format '{{.Os}}/{{.Architecture}}'` gibt `linux/amd64` | kein Lauf für `linux/arm64`, `darwin`, `windows` — deshalb steht im Handbuch „gebaut, aber nicht ausgeführt“ |

**Nacharbeit zum Review vom 2026-10-21** (`docs/reviews/2026-10-21-review-slice-v1-abschluss-release-v0-1-1.md`,
1 HIGH, 4 MEDIUM, 2 LOW; F-8 und F-9 sind Vertragswortlaut und gehen an den Architect):
F-1 — Eintrag und §7 nennen jetzt alle vier Versuche des Jobs mit Zeitstempeln, die drei roten
eingeschlossen; selbst nachgefahren mit `gh api …/attempts/<n>/jobs` für n = 1 bis 4 und den Logs der
Versuche 2 und 3. F-2 — siehe den Absatz darunter. F-3 — die Berichtigung der Token-Rechte steht als
*Angabe der Koordination* mit Datum, daneben die eigene Messung `gh secret list` (Zeitstempel des
Geheimnisses unverändert). F-4 — der Grund für den fehlenden Ersatzbeleg 1 steht bei Punkt 2 des
Eintrags; er nannte zuerst das falsche Repository und ist bei der Closure berichtigt (unten, F-V1).
F-5 — Kopf von `releasing.md` auf Version 0.14 und Stand 2026-10-21. F-6 — der
Tap-Commit ist jetzt mit `git ls-remote … refs/heads/main` belegt, der Form des Vertrags. F-7 — die
Abweichung des Bezug-Commits (`docs/maintainer/` außerhalb der Positivliste) steht bei Punkt 11 des
Eintrags.

**Zu F-2 des Reviews (2026-10-21): die Tap-Belege bleiben beim Nehmer.** §1 und §6 Randform (e)
weisen die Belege des Taps `slice-erster-release-homebrew-nachweis` zu, und dessen §1 nimmt sie an
(„Wiederanlauf des Jobs, Nachtrag von Hand und alle Belege des Taps gehören diesem Slice“). Statt §1
und §6 zu weiten, sind die überschießenden Aussagen aus dem Eintrag und aus `releasing.md` §1
entfernt: Dort steht jetzt allein, was `SPEC-067` Punkt 17 (8) für diese Zeile verlangt —
Lauf-Kennung, grün, Hash des Tap-Commits mit Datum — plus die Grenze, dass Rechte, Inhalt und
Annahme der Formel dem Nachweis gehören. Der Plan ist damit nicht geändert, und keine Sendung
wandert in einen anderen Slice (`AGENTS.md` §3.13, nichts einzutragen).

**Messung des Implementers, die nicht Beleg dieses Slice ist** (Eingabe für
`slice-erster-release-homebrew-nachweis`, dort §1): Die Formel `Formula/pgwire-recorder.rb` im Tap
trug am 2026-10-21 `version "0.1.1"`, und ihre vier `sha256`-Werte waren zeichengleich mit den vier
`darwin`- und `linux`-Zeilen der geladenen `SHA256SUMS` des Release (sortierter `diff` leer); ihre
vier URLs zeigten auf die Anhänge von `v0.1.1`. Der Nachweis misst das gegen den Stand seines
eigenen Laufs, diese Zeile nimmt ihm nichts ab und ändert seinen Umfang nicht.

**Durchsicht nach Quantoren** (Schritt 17 des Commands; Merkmal der neuen Sätze: Namen ohne Version, ein Digest für beide Registries, `latest` auf demselben Digest). Treffer in `docs/user/benutzerhandbuch.md` und `README.md`:

| Fundstelle | Satz | Entscheidung |
|---|---|---|
| Handbuch §2, Zeile 124 | „Ohne `--ignore-missing` verlangt `sha256sum` alle sechs Dateien“ | unverändert, mit Probe belegt (Zeile oben) |
| Handbuch §2, Zeile 139 | „Beide Registries geben das Image ohne Anmeldung heraus“ | unverändert, mit Probe belegt (leeres `DOCKER_CONFIG`) |
| Handbuch §2, Zeile 140 | „`latest` zeigt in beiden Registries auf denselben Digest und gibt dieselbe Ausgabe“ | enger gefasst gegenüber `v0.1.0` („zeigt auf dieselbe Version“) und mit Probe belegt |
| Handbuch §3, Zeile 175 | „`make binaries` baut das Binary für jede der sechs Plattformen“ | unverändert mit Grund: Satz über den lokalen Bau, nicht über den Release; Gegenstand von `slice-v1-abschluss-binary-namen`, durch `make test-container` gedeckt |
| `README.md` | „sechs Binaries ohne Version im Namen und `SHA256SUMS`“ | neu, durch `gh release view` belegt (Punkt 9) |

**Risiko-Belege aus §6** (die Ausgänge trägt die Closure, nicht dieser Abschnitt):

- `grep -rn 'pgwire-recorder_' docs/user README.md` am Endstand: kein Treffer (Exit 1).
- `docs/maintainer/releasing.md` §1 *Stand*: nachgezogen; der Satz „Es ist noch kein Tag gesetzt worden“ und die drei Sätze „hatte noch keinen Lauf“ sind weg, an ihrer Stelle steht der Zustand nach `v0.1.1`: Release, Image in beiden Registries und die Formel `0.1.1` im Tap.
- Der Lauf des Tags war im Job `release` grün; ein Fehler im Skript oder Vertrag trat nicht auf, ein Folge-Tag `v0.1.2` ist nicht nötig.

**Nicht belegt (Grenzen, die der Slice stehen lässt):** die Belege des Taps — Rechte des Tokens, Inhalt und Annahme der Formel, `brew tap` und `brew install` unter Linux und macOS —, alle Gegenstand von `slice-erster-release-homebrew-nachweis`; *welches* Recht der Betreiber zwischen den Versuchen änderte (Angabe, aus Repo und Lauf nicht auditierbar, `SPEC-067` Punkt 15); der Wiederanstoß eines roten **Laufs** (nur der eines roten **Jobs** ist belegt, dreimal); der Ersatzbeleg 1 aus `CO-001` (den Gegenproben-Container `pgr-gp` und damit seine Ablage gibt es nicht mehr, Grund im Eintrag bei Punkt 2); ein Lauf von `make gegenproben-host` am getaggten Hash (Carveout `CO-001`); die Mutationen der beiden gefahrenen Gegenproben; `darwin`, `windows` und `linux/arm64` als ausgeführte Plattformen.

**Größenmarke** (§4, Metrik: Summe aus Einfügungen und Löschungen von `git diff --numstat 0256ddff..<letzter Commit> -- VERSION docs/user docs/maintainer README.md`), am Endstand: 227, unter 500. Gemessen an fünf Ständen — `a3ad8222`, `56c46178`, `ccc23dac`, `0c8a13d3` (vier davon die Verifikation, F-V4) und `0ee9f164` (Berichtigungen der Closure) — jedes Mal 227; die Zahl ist stabil, weil die späteren Commits Zeilen in Dateien ändern, die in diesem Bereich neu sind oder ersetzt werden. Der Klammersatz davor behauptete, spätere Commits änderten nur `docs/plan/planning/`; `0c8a13d3` berührt zwei der vier Pfade (F-V4), die Closure-Berichtigung einen. Die Basis `0256ddff` enthält Änderungen aus `slice-v1-abschluss-binary-namen` und `…-cancel-ohne-schluessel`, die nicht zu diesem Slice gehören.

**Läufe am Endstand:** `make gates` Exit 0 im ersten Versuch, 265 s, 2026-10-21, Last 4,6 bei 20 Kernen, kein `FAIL` im Log (der lastempfindliche Test `TestE2EPlayVergleichAbweichung` fiel in keinem Lauf aus, keine Wiederholung) · `make docs-check` grün (**1315** Dateien, 0 Befunde, am Stand `0c8a13d3`; die Zahl 1314 stammte von einem Stand vor dem Review-Report, F-V2 — Ergebnis unverändert grün, die Läufe der Closure stehen oben). Dieselbe Kette war auf den beiden Ständen davor ebenso grün (272 s vor den Belegen des Jobs `homebrew`, 271 s vor der Nacharbeit zum Review). Nach dem letzten Gate-Lauf ändert allein diese Zeile in `docs/plan/planning/` (inert nach `SPEC-052` Punkt 12). Keine Gegenprobe am Endstand (Entscheidung des Nutzers vom 2026-10-21; die beiden kurzen Läufe stehen bei den Punkten 2 und 3 des Eintrags).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Das Repo deklariert eine Sub-Area für das gesamte Repo (`harness/conventions.md`, Kürzel `REPO`); der Slice berührt sie, die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen. Treffer: `folge-slice-liefert-handbuch-teil-nicht` (der Handbuch-Teil aus `slice-v1-abschluss-binary-namen` hat hier seine Adresse mit Kennung, §1; der Stand des Zählers ist beim Start zu lesen, erreicht er mit diesem Slice drei, braucht er einen Folge-Slice), `folge-slice-adresse-nimmt-nicht-an` (die Sendungen in §1 stehen in den Nehmern, `AGENTS.md` §3.13), `verhalten-nur-unter-linux-geprueft` (hier als Randform (h) entschieden), `empfehlung-als-nutzerentscheidung-bezeichnet` (Version und Schema sind Entscheidungen des Nutzers vom 2026-10-20; das Folge-Ja zu `v0.1.2` ist seit dem 2026-10-20 eine Entscheidung des Nutzers (Kopf), das Ja zum Tag `v0.1.1` ist seit dem 2026-10-21 eine Entscheidung des Nutzers im Wortlaut des Kopfes, ohne Dehnung über `v0.1.1` und `v0.1.2` hinaus).

**Zählung** (Planner, 2026-10-20, `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 2** (1: Freigabe und Tag, mit `VERSION`, Zeile der Änderungshistorie und Eintrag; 2: Bezug nach dem Tag in Handbuch, `README.md` und `releasing.md` §1). **Schichten: zwei** — Release (Datei `VERSION`, Lauf der Pipeline als Handlung von Claude nach dem Ja des Nutzers) und Nutzer- und Wartungs-Doku (`docs/user/`, `docs/maintainer/`, `README.md`); Planung zählt nicht, Produkt- und Werkzeug-Code berührt der Slice nicht (§1). Mit der Sendung aus `slice-v1-abschluss-binary-namen` (Handbuch- und README-Teil) nachgezählt: Der Teil liegt im zweiten Liefer-Punkt, kein neuer Punkt. Größe: Schätzung und Marke in §4.

**Nachzählung der Änderung vom 2026-10-21** (Planner, `AGENTS.md` §3.13): Das Ja des Nutzers, die Handlung von Claude als Tag-Setzer, der Carveout `CO-001` (Ersatzbelege statt der vollen Läufe der Punkte 2 und 3) und der Satz zu Carveouts in `releasing.md` §10 sind Inhalt der beiden bestehenden Liefer-Punkte (Eintrag der Freigabe; Spiegel in `releasing.md`), kein dritter. **Liefer-Punkte bleiben 2, Schichten zwei** (Release, Nutzer- und Wartungs-Doku). Die Schätzung der Marke in §4 (250 bis 300) steigt um den Satz in `releasing.md` und die Zeilen zu `CO-001` im Eintrag, rund 30 bis 50, bleibt weit unter 500.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält Spezifikation und Entscheidungen vor der Doku; `SPEC-067` und [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) gehen dem Eintrag voraus).
