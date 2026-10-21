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

**Entscheidung:** [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) (Accepted, Entscheidung des Nutzers vom 2026-10-20) · [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md) · [ADR-0053](../../adr/0053-homebrew-tap-in-der-release-pipeline.md). Der Nutzer will den nächsten stabilen Tag als `v0.1.1` (Entscheidung des Nutzers vom 2026-10-20); ein Ja zum Tag selbst liegt nicht vor. **Entscheidung des Nutzers vom 2026-10-20** (Angabe der Koordination, im Plan festgehalten bei der Closure von `slice-v1-abschluss-binary-namen`): Das Ja des Nutzers zum Tag und zur Veröffentlichung von `v0.1.1` schließt einen Folge-Tag `v0.1.2` ein, falls der erste Lauf einen Fehler im Skript oder im Vertrag zeigt, ohne neue Frage; Secrets und andere Freigaben bleiben beim Nutzer, und jeder weitere Tag danach braucht ein neues Ja.

**Verantwortlich:** pt9912
**Autor:** Planner. **Datum:** 2026-10-20.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründung je Punkt.

**Ziel:** Der stabile Tag `v0.1.1` ist nach der Freigabe-Checkliste (`docs/maintainer/releasing.md` §10, `SPEC-067` Punkt 17) mit Beleg je Punkt vorbereitet; setzt der Nutzer ihn nach seinem ausdrücklichen Ja, liegen am Release die sechs Binaries `pgwire-recorder-<os>-<arch>` ohne Version im Namen und `SHA256SUMS`, das Image in beiden Registries mit `0.1.1` und `latest`, und der Job `homebrew` hat die Formel im Tap versucht. Nach dem Lauf beschreiben Benutzerhandbuch und `README.md` den Bezug von `v0.1.1` mit den neuen Namen im Ist-Zustand. **Der Slice setzt keinen Tag und veröffentlicht nichts:** Er beschreibt die Vorbedingungen und den Lauf (Trockenlauf, dann der Tag), der Tag ist eine Handlung des Nutzers nach seinem ausdrücklichen Ja (`SPEC-067` Punkt 18); ein Implementer startet ihn nicht.

**Übernommen aus `slice-v1-abschluss-binary-namen`** (dort §1 und §6, `AGENTS.md` §3.13): der Abschnitt „Binary vom GitHub-Release“ und der Abschnitt zum Image in §2 *Installation* des Benutzerhandbuchs und die Zeile zu `v0.1.0` in `README.md` werden auf `v0.1.1` und die neuen Namen umgeschrieben, belegt durch den Lauf des Tags (`SPEC-067` Punkt 18; `BEO-REPO/folge-slice-liefert-handbuch-teil-nicht`). Dieser Slice ist die Adresse; `grep` nach `slice-v1-abschluss-binary-namen` findet die Sendung in §1, DoD 2 und §8.

**Bindung für `slice-erster-release-homebrew-nachweis`** (Geber: dieser Slice, `AGENTS.md` §3.13): Der Nachweis setzt einen veröffentlichten stabilen Tag voraus, dessen Binaries `pgwire-recorder-<os>-<arch>` heißen und dessen Lauf den Job `homebrew` hatte; dieser Slice liefert den Tag (nach dem Ja des Nutzers) und nennt im Eintrag der Freigabe das Ergebnis des Jobs `homebrew` (Lauf-Kennung, grün oder rot). Der Nachweis steht dort in §1, §4 und §6 mit der Kennung dieses Slice und als Start-Trigger.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Setzen des Tags und das Veröffentlichen — Handlung des Nutzers nach seinem ausdrücklichen Ja (Klasse 3, anderer Vorgang). Grund: `SPEC-067` Punkt 18 und `AGENTS.md`; ein Plan oder Slice ersetzt das Ja nicht.
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

- [ ] [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität): **Freigabe und Tag.** Ein Release-Commit trägt `VERSION` mit `0.1.1` und die Zeile `0.1.1` der Änderungshistorie des Handbuchs (§11, Quelle der Release-Notizen, `releasing.md` §9). Der Eintrag der Freigabe (`docs/maintainer/freigaben/0.1.1.md`, eine Kopie der Tabelle aus `releasing.md` §10) trägt die Belege der Punkte 1 bis 6 für genau diesen Hash (`make gates`, Gegenproben-Container, `make gegenproben-host`, `grep`, Trockenlauf, Scan) und dazu das ausdrückliche Ja des Nutzers (Punkt 7, Wortlaut und Datum; **ohne es kein Tag**). Nach dem Tag des Nutzers trägt er die Belege der Punkte 8 bis 10 aus dem Lauf: Lauf-Kennung, Index-Digest je Registry gleich, Smoke, `latest`, die sieben Anhänge mit den Namen `pgwire-recorder-<os>-<arch>` (Ausgabe von `gh release view`), die Beobachtungen nach `releasing.md` §10 und, getrennt, das Ergebnis des Jobs `homebrew` (Lauf-Kennung, grün oder rot; Belege des Taps selbst gehören dem Nachweis, §1). Was der Lauf nicht zeigt, steht als Grenze im Eintrag, nie als Behauptung (`AGENTS.md` §3.11).
- [ ] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung): **Bezug nach dem Tag** (Punkt 11 der Checkliste, `SPEC-067` Punkt 18; Sendung aus `slice-v1-abschluss-binary-namen`, §1), in einem eigenen Commit nach dem Lauf, der nur `docs/user/`, `README.md`, `docs/maintainer/` und Planung ändert: Kopf des Handbuchs (*Software-Version*, *Gültigkeitsbereich*) nennt die Ausgabe von `version` des **veröffentlichten** Binaries; §2 *Installation* (Binary vom GitHub-Release, Image aus der Registry) nennt `v0.1.1`, die sechs Namen `pgwire-recorder-<os>-<arch>` ohne Version, die Prüfung mit `sha256sum -c` und die Tags `0.1.1` und `latest` mit dem Digest des Laufs — **nur mit Befehlen, die nach dem Tag einmal gelaufen sind** (Beleg im Eintrag: Befehl, Lauf, Ausgabe; die Beispiele liegen als Datei vor und laufen gegen das veröffentlichte Binary, `AGENTS.md` §3.11); `README.md` nennt die Version `0.1.1` mit den Adressen im Ist-Zustand; `docs/maintainer/releasing.md` §1 (*Stand*) beschreibt den Zustand nach dem Tag (Tag gesetzt, Läufe, Formel nach dem Ausgang des Jobs) im Ist-Zustand ohne Chronik. Der Abschnitt zu Homebrew bleibt dem Nachweis (§1).
- [ ] `make gates` grün; die Gegenproben am getaggten Commit sind Punkte 2 und 3 des Eintrags: voll, oder mit dem Carveout `CO-001` und den Ersatzbelegen aus dessen Text, wenn der Nutzer ihn aktiviert (Randform (m)); der Eintrag schreibt je Punkt, was gilt (Wortlaut im Carveout).
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

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `VERSION` | update | `0.1.1` im Release-Commit, der getaggt wird; ändert das Wurzelverzeichnis und ist damit nicht inert nach `SPEC-052` Punkt 12 (volle Kette) |
| `docs/user/benutzerhandbuch.md` §11 | update | Zeile `0.1.1` der Änderungshistorie im Release-Commit, vor dem Tag |
| `docs/maintainer/freigaben/0.1.1.md` | neu | Eintrag der Freigabe, Kopie der Tabelle aus `releasing.md` §10, Belege je Punkt |
| `docs/maintainer/releasing.md` §10 | update | Ort des Eintrags „für jeden späteren Release“, Zeile des Jobs `homebrew` unter Punkt 8 (Spiegel von `SPEC-067` Punkt 17), im Release-Commit oder davor |
| `docs/user/benutzerhandbuch.md` (Kopf, §2 *Installation*), `README.md` | update | Bezug nach dem Tag (`SPEC-067` Punkt 18): `v0.1.1`, neue Namen, Image-Tags, nur belegte Befehle |
| `docs/maintainer/releasing.md` §1 | update | *Stand* nach dem Tag |
| Lauf `release.yml` auf dem Tag (Trockenlauf davor) | Handlung | kein Dateiinhalt; der Tag ist Handlung des Nutzers, der Trockenlauf veröffentlicht nichts (`SPEC-067` Punkt 2) |

Nicht in dieser Liste: `tools/`, `.github/`, `Dockerfile`, `cmd/`, `internal/` und die Verträge (§1), der Abschnitt zu Homebrew in Handbuch und `README.md` (§1, Nachweis).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-v1-abschluss-binary-namen` und `slice-v1-abschluss-homebrew` liegen in `done/` (die Namen im getaggten Commit und der Job `homebrew` in `release.yml`; Reihenfolge in §5 von [welle-v1-nachtrag](../welle-v1-nachtrag.md)), [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) ist `Accepted` (Entscheidung des Nutzers vom 2026-10-20), und die Vorbedingungen des Betreibers stehen (`releasing.md` §4): die Geheimnisse für beide Registries und `HOMEBREW_TAP_GITHUB_TOKEN` sind gesetzt, das Repository des Taps `pt9912/homebrew-pgwire-recorder` besteht (Handlung des Betreibers; der Slice prüft, legt nichts an und liest keinen Wert). Das Ja zum Tag ist **kein Start-Trigger**, sondern ein Schritt im Slice (DoD 1): Bis dahin arbeitet der Slice nur am Release-Commit und an den Belegen 1 bis 6.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe aus Einfügungen und Löschungen von `git diff --numstat 0256ddff..<letzter Commit des Slice> -- VERSION docs/user docs/maintainer README.md` (Tests eingeschlossen, `AGENTS.md` §3.14; der Slice liefert keinen Test- oder Werkzeug-Code) **über 500** liegt. Schätzung des Planners in dieser Metrik, ohne Messung: Handbuch (Abschnitte zu Binary und Image umgeschrieben, ersetzte Zeilen zählen doppelt) rund 100, Kopf, §11 und `README.md` rund 20, `releasing.md` §1 rund 40, Eintrag der Freigabe rund 90 bis 140, zusammen 250 bis 300. Kein vorgemerkter Schnitt nötig; fiele der Slice darüber, wäre die Naht (a) Release-Commit, Eintrag bis Punkt 7 und Tag, (b) Bezug nach dem Tag.
- `in-progress` → `open`: Das ausdrückliche Ja des Nutzers zum Tag bleibt aus, oder ein vorbereitender Beleg (Punkte 1 bis 6) wird rot und die Ursache liegt außerhalb dieses Slice (etwa ein roter Scan wegen eines Befunds ohne Behebung: Entscheidung des Nutzers zur Schwelle oder Carveout, kein Tag; `SPEC-067` Punkt 16). Ein Slice, der auf das Ja wartet, hält das WIP-Limit; deshalb geht er zurück nach `open/`, statt zu warten.
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
| (b) Wer taggt, wann | der Nutzer, nach seinem ausdrücklichen Ja; der Slice setzt keinen Tag, ein Implementer startet ihn nicht; ohne Ja kein Tag (Rückführung nach `open`, §4) | `SPEC-067` Punkt 18; `releasing.md` §10 |
| (c) Reihenfolge der Commits | (1) Release-Commit mit `VERSION` und der Zeile `0.1.1` (Zeile **vor** dem Tag, sie ist die Quelle der Notizen), (2) Tag auf diesem Commit, (3) Bezug-Commit nach dem Lauf; die Punkte 1 bis 3 der Checkliste gelten für den Hash von (1) | `SPEC-067` Punkte 17, 18; `releasing.md` §10 |
| (d) Rot nach dem Tag | Umgebung (Token, Schreibrecht, Sichtbarkeit, Netz): Wiederanstoß desselben Tags; Skript, Text der Datei oder Vertrag: Korrektur auf dem Hauptzweig und neuer höherer Tag mit neuem Ja; Digest oder Anhang abweichend: immer neuer Tag. **Das Ja zum ersten Tag schloss `v0.1.1` ein; das Ja zu `v0.1.1` schließt `v0.1.2` ein** (Entscheidung des Nutzers vom 2026-10-20, Kopf) — jeder weitere Tag danach braucht ein neues Ja | `SPEC-067` Punkte 11, 18; Entscheidung des Nutzers im Kopf |
| (e) Job `homebrew` rot oder grün | Ein roter Job ändert weder Release noch Registries; der Slice nennt sein Ergebnis im Eintrag und gibt Wiederanlauf oder Nachtrag von Hand an den Nachweis (kein neuer Tag); Belege des Taps sind nicht Teil dieses Slice | `SPEC-067` Punkt 19 (1), (7); Nachweis |
| (f) `v0.1.0` und das alte Schema | `v0.1.0` bleibt mit Namen mit Version; das Handbuch nennt nach dem Bezug-Commit nur `v0.1.1` und die neuen Namen und sagt nichts über die alten (Ist-Zustand, `AGENTS.md` §3.11); der Schritt `homebrew` lehnt das alte Schema ab | `SPEC-055` Punkt 3; `SPEC-067` Punkt 19 (3); Rückgabe an den Architect unten |
| (g) `latest` | wandert von `0.1.0` auf `0.1.1`, wenn Version stabil, Vergleich der Digests beider Registries grün und `0.1.1` höher als der vorhandene `latest`; der Tag `0.1.0` des Images bleibt | `SPEC-067` Punkt 8 |
| (h) Plattform | Linux (`amd64`) im Lauf geprüft: Binary gestartet, Image gestartet; `darwin`, `windows` und `linux/arm64` gebaut, nicht ausgeführt: Grenze, die das Handbuch so nennt | `SPEC-067` Punkt 15; `SPEC-055` Punkt 13 (4) |
| (i) Aufrufarten, Abbruch | entfällt: Der Slice liefert kein Werkzeug und kein Ergebnis, das ein Hook liest; Treiber der Pipeline ist `release.sh` (nie unter `make`) | `SPEC-067` Punkt 13 |
| (j) Ort des Eintrags der Freigabe | **entschieden (Architect, 2026-10-21): `docs/maintainer/freigaben/0.1.1.md`**; der Vertrag nennt `freigaben/<Version>.md` jetzt für jeden späteren Release, mit oder ohne Welle | `SPEC-067` Punkt 17 |
| (k) Job `homebrew` in der Checkliste | **entschieden: Teil von Punkt 8**, getrennt darunter, mit Lauf-Kennung, grün oder rot, bei grün Hash des Commits im Tap (`git ls-remote`, Datum), bei rot Stichwort der Meldung; ein roter Job macht Punkt 8 nicht unerfüllt; Mindestmenge bleibt elf, keine Neunummerierung (`releasing.md` §10 zieht die Zeile von Punkt 8 mit, Ort der Doku: dieser Slice) | `SPEC-067` Punkt 17 |
| (l) Bezug des zweiten Tags | **entschieden: dieselbe Form wie beim ersten Tag; der Bezug ersetzt den des vorigen Tags und schweigt über frühere Namen** (Randform (f) bestätigt) | `SPEC-067` Punkt 18 (Absatz „Jeder weitere Tag“) |
| (m) Punkte 2 und 3 der Checkliste (Gegenproben am getaggten Commit) | **Empfehlung des Architect, Entscheidung des Nutzers steht aus:** Die Pflicht bleibt im Vertrag; für `v0.1.1` ruht sie über den Carveout `CO-001` (Entwurf, `docs/plan/carveouts/`), mit Ersatzbelegen in Minuten; Trigger und Folge-Slices dort | `docs/plan/carveouts/CO-001-gegenproben-am-getaggten-commit.md`; [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md); `AGENTS.md` §3.6 |

**Rückgaben an den Architect:** beantwortet am 2026-10-21 — (1) Ort: Randform (j); (2) Job `homebrew`: Randform (k); (3) zweiter Tag: Randform (l). Neu entschieden: (m).

**Entscheidung des Nutzers vom 2026-10-20 zu Randform (d)** (vorher als Frage geführt): Das Ja zu `v0.1.1` schließt einen Folge-Tag `v0.1.2` ein, falls der Lauf von `v0.1.1` einen Fehler im Skript oder im Vertrag zeigt; kein weiteres Ja dafür. Nicht eingeschlossen: Secrets, andere Freigaben, jeder Tag nach `v0.1.2`. Das Ja zum Tag `v0.1.1` selbst gibt der Nutzer separat (DoD 1, noch nicht gegeben).

**Risiken:**

- Der Lauf des Tags ist die erste Außenwirkung des Jobs `homebrew` und des neuen Schemas auf GitHub (Rechte des Tap-Tokens, Repository des Taps, Wortlaut von `gh release download`, Annahme durch Homebrew); der Trockenlauf führt den Job nicht (`SPEC-067` Punkt 19 (1)). Rot dort ist kein Grund für einen neuen Tag (Randform (e)) — **Ausgang:** weiter offen, übergeben an `slice-erster-release-homebrew-nachweis` (dort §6).
- Ein Fehler im Skript, der erst der Lauf des Tags zeigt, kostet einen Folge-Tag `v0.1.2`, den das Ja zu `v0.1.1` deckt (Entscheidung des Nutzers vom 2026-10-20), ein weiterer brauchte ein neues Ja — **Ausgang:** offen bis Closure (eingetreten: Kennung des Folge-Slice und des Ersatz-Tags; entfallen: der Lauf war grün).
- Die Fundstellen des alten Schemas im Handbuch und in `README.md` (`pgwire-recorder_` mit Version) sind der Gegenstand des Bezug-Commits — **Ausgang:** offen bis Closure (`grep -rn 'pgwire-recorder_' docs/user README.md` am Endstand ohne Treffer).
- `docs/maintainer/releasing.md` §1 (*Stand*) beschreibt den Zustand von vor dem ersten Tag („Es ist noch kein Tag gesetzt worden“, „Eine Formel im Tap gibt es nicht“); der Planner fand das beim Lesen, ob der Stand an anderer Stelle nachgezogen wird, ist nicht geprüft — **Ausgang:** offen bis Closure (DoD 2 trägt den Stand nach dem Tag).

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

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen. Treffer: `folge-slice-liefert-handbuch-teil-nicht` (der Handbuch-Teil aus `slice-v1-abschluss-binary-namen` hat hier seine Adresse mit Kennung, §1; der Stand des Zählers ist beim Start zu lesen, erreicht er mit diesem Slice drei, braucht er einen Folge-Slice), `folge-slice-adresse-nimmt-nicht-an` (die Sendungen in §1 stehen in den Nehmern, `AGENTS.md` §3.13), `verhalten-nur-unter-linux-geprueft` (hier als Randform (h) entschieden), `empfehlung-als-nutzerentscheidung-bezeichnet` (Version und Schema sind Entscheidungen des Nutzers vom 2026-10-20; das Folge-Ja zu `v0.1.2` ist seit dem 2026-10-20 eine Entscheidung des Nutzers (Kopf), das Ja zum Tag `v0.1.1` selbst nicht und bleibt ein Schritt der DoD).

**Zählung** (Planner, 2026-10-20, `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 2** (1: Freigabe und Tag, mit `VERSION`, Zeile der Änderungshistorie und Eintrag; 2: Bezug nach dem Tag in Handbuch, `README.md` und `releasing.md` §1). **Schichten: zwei** — Release (Datei `VERSION`, Lauf der Pipeline als Handlung des Nutzers) und Nutzer- und Wartungs-Doku (`docs/user/`, `docs/maintainer/`, `README.md`); Planung zählt nicht, Produkt- und Werkzeug-Code berührt der Slice nicht (§1). Mit der Sendung aus `slice-v1-abschluss-binary-namen` (Handbuch- und README-Teil) nachgezählt: Der Teil liegt im zweiten Liefer-Punkt, kein neuer Punkt. Größe: Schätzung und Marke in §4.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (das Repo enthält Spezifikation und Entscheidungen vor der Doku; `SPEC-067` und [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) gehen dem Eintrag voraus).
