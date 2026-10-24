# Slice slice-release-v0-2-0: Release v0.2.0 (zweiter stabiler Tag mit `--version`, `-v` und `play --diff-all`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle (Entscheidung des Planners, 2026-10-24). Der Slice hat keine Closure-Bedingung, die von seiner DoD verschieden ist, und keine Anforderung einer Welle hängt an ihm.

**Bezug:** [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-24`](../../../../spec/lastenheft.md#lh-fa-24--vergleich-der-antworten-beim-einspielen), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität)

**Berührte Spec-Stellen:** `SPEC-067` (Punkte 2, 3, 8, 11, 17, 18 und 19, [`harness/vertraege/`](../../../../harness/vertraege/SPEC-067-release.md); gelesen, der Slice ändert den Vertrag nicht) · `SPEC-060` (Punkt 14, Pflicht zum vollen Lauf der Gegenproben am getaggten Commit; gelesen) · `SPEC-055` (Punkt 3, Namen der Binaries; gelesen) · `SPEC-052` (Punkt 12, Positivliste; gelesen) · `LH-FA-19.a` · `LH-FA-16.a` · `LH-FA-01` · `LH-FA-24` · `LH-QA-03`

**Entscheidung:** [ADR-0051](../../adr/0051-veroeffentlichung-auf-github-actions.md) · [ADR-0053](../../adr/0053-homebrew-tap-in-der-release-pipeline.md) · [ADR-0054](../../adr/0054-binaries-ohne-version-im-dateinamen.md) · [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md). **Entscheidung des Nutzers vom 2026-10-24:** Version **v0.2.0** und das Ja zum Tag (Wortlaut „release“, Punkt 7 der Checkliste, `SPEC-067` Punkt 18); den Tag setzt die Koordination erst nach grünem Trockenlauf und nochmaliger Bestätigung. Auf die Frage nach den Entscheidungen zu [`CO-002`](../../carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md) antwortete der Nutzer wörtlich „deine empfehlungen“: neues `CO-002` (kein verlängertes `CO-001`), kein dritter Ersatzbeleg, kein Folge-Tag, damit Annahme; der Carveout ist Aktiv (nur der engste Fall je Gegenprobe mit `GP_MUT`, die Grenze im Eintrag genannt). Die Wiedergabe der Wortlaute stammt von der Koordination und ist nicht nachgefahren. Secrets setzt allein der Nutzer; sie sind laut Nutzer seit `v0.1.1` hinterlegt, ob sie noch tragen (Ablaufdatum, Rechte), zeigt erst der Lauf (§6).

**Verantwortlich:** pt9912 (Zweig impl-rel-v020b)
**Autor:** Planner. **Datum:** 2026-10-24.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründung je Punkt.

**Ziel:** Der stabile Tag `v0.2.0` ist nach der Freigabe-Checkliste (`docs/maintainer/releasing.md` §10, `SPEC-067` Punkt 17) mit Beleg je Punkt vorbereitet; setzt die Koordination ihn nach dem Ja des Nutzers, liegen am Release die sechs Binaries `pgwire-recorder-<os>-<arch>` und `SHA256SUMS`, das Image in beiden Registries mit `0.2.0` und `latest`, und der Job `homebrew` hat die Formel `0.2.0` im Tap versucht. Mit der Closure ist der Carveout [`CO-002`](../../carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md) verbraucht (Auflösungs-Trigger, Punkt 2 und 3 der Freigabe laufen nach ihm). Nach dem Lauf beschreiben Benutzerhandbuch und `README.md` den Bezug von `v0.2.0` im Ist-Zustand und nennen nur Befehle, die nach dem Tag einmal gelaufen sind. Seit `v0.1.1` ausgeliefert sind `--version` und `-v` (`LH-FA-01`) und `play --diff-all` (`LH-FA-24`) samt Handbuch; das Upstream-Ende bei Extended Query (`slice-v1-abschluss-verbindungsende-upstream-extended`) lieferte Tests und Testdaten ohne Produkt-Code und kommt in der Zeile der Änderungshistorie nicht vor.

**Wer was tut** (Rollen, `SPEC-067` Punkt 18: „ein Implementer startet keinen Tag“):

| Rolle | Tut | Tut nicht |
|---|---|---|
| Koordination (Claude) | merged den Release-Commit per ff auf `main` und pusht, startet den Trockenlauf (`workflow_dispatch`) auf genau diesem Stand, setzt und pusht den Tag nach dem Ja, beobachtet den Lauf, gibt dem Implementer die Lauf-Kennung | schreibt keinen Eintrag, keinen Carveout |
| Architect | hat `CO-002` festgelegt (§6, Rückgabe 1, erledigt) und prüft die Randformen in §6 | ändert keinen Vertrag ohne Anlass; nimmt den Carveout nicht an |
| Nutzer | hat `CO-002` angenommen (Wortlaut im Carveout), entscheidet die offenen Punkte in §6, setzt und pflegt Secrets, nennt sein Ja zum Tag im Wortlaut | — |
| Implementer | Release-Commit (`VERSION`, Zeile der Änderungshistorie), Belege 1 bis 6 mit Befehl, Hash, Datum, Dauer, Last, Eintrag `docs/maintainer/freigaben/0.2.0.md`, nach dem Lauf die Belege 8 bis 10 aus dem Log und der Bezug-Commit (Punkt 11) | setzt keinen Tag, startet keinen Lauf auf GitHub, fasst kein Geheimnis an, ändert weder Code noch Vertrag noch Gate |
| Reviewer, Verifier | lesen den Eintrag gegen die Originale (`gh`, `git`, Binary) | — |

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Das Setzen des Tags, das Anstoßen des Trockenlaufs und der Push nach `origin` durch den Implementer-Lauf — Handlung der Koordination nach dem Ja des Nutzers (Klasse 3, anderer Vorgang). Grund: `SPEC-067` Punkt 18; der Plan ersetzt weder das Ja noch die Handlung (wie bei `v0.1.1`: der Implementer hatte keinen Zugriff auf `gh workflow run`).
- Die Änderung von [`CO-002`](../../carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md) und von `CO-001` — der Carveout ist Aktiv und gehört dem Nutzer und dem Architect (Klasse 3, anderer Vorgang). Grund: Eine Ausnahme von einer Pflicht der Checkliste (`AGENTS.md` §3.6) ist eine Entscheidung des Nutzers und kein Planinhalt; dieser Slice nutzt sie für die Punkte 2 und 3 und **verbraucht** sie mit seiner Closure (Trigger von `CO-002`: der Tag `v0.2.0` gesetzt und der Eintrag geschlossen, oder der Slice ohne Tag zurückgeführt; danach `git mv` nach `done/`). `CO-001` bleibt Aktiv als Träger der Ursache.
- Die Verkürzung der Gegenproben (`slice-harness-gegenproben-schneller`, `slice-harness-gegenproben-grenze-rest`, beide Klasse 1) — sie lösen den Carveout auf; dieser Slice wartet nicht auf sie und liefert sie nicht. Grund: Die Zeitgrenze von [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md) macht den vollen Lauf der Punkte 2 und 3 heute unmöglich; ein Tag, der darauf wartete, käme nicht.
- Eine Korrektur an `release.sh`, `release.yml`, `SPEC-067`, den Gegenproben oder der Vorlage der Formel, die der Lauf verlangt — Folge-Slice und höherer Tag (Klasse 3, `SPEC-067` Punkt 18 „Rot“). Grund: Ein Lauf führt die Fassung des getaggten Commits aus; dieser Slice liefert keinen Werkzeug- und keinen Produkt-Code (Klasse 4: nur `VERSION`, Doku und Planung).
- Die Belege des Taps für macOS und `arm64` und das Schreibrecht des Tokens ohne Eingriff — bleiben bei `slice-erster-release-homebrew-nachweis` (Klasse 2, Bestand: dort geliefert, für `0.1.1` unter Linux/amd64). Grund: Dieser Slice nennt nur das Ergebnis des Jobs `homebrew` für `0.2.0` im Eintrag (Zeile unter Punkt 8); den Nachweis für `0.2.0` unter Linux/`amd64` hat der Implementer nach dem Tag mit `make homebrew-nachweis` gefahren (Punkt 11, damit der Abschnitt im Handbuch nur Belegtes nennt), und das Handbuch nennt nur Befehle nach einem Lauf (`AGENTS.md` §3.11).
- Das Freshness-Audit — optional, kein Punkt der Checkliste, sperrt den Tag nie (`SPEC-067` Punkt 17). Grund: Niemand hat es angefordert; ein durchgeführtes stünde als eigener Absatz im Eintrag.
- Der Release `v0.1.1` und `v0.1.0` — unverändert (Klasse 2). Grund: Ein Tag wird nie verändert (`SPEC-067` Punkt 11); der Bezug **ersetzt** den des vorigen Tags und schweigt über frühere Namen (`SPEC-067` Punkt 18, „Jeder weitere Tag“), die Zeilen der Änderungshistorie bleiben.
- Änderung des Lastenhefts — es legt weder Version noch Namen fest.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew), [`LH-QA-03`](../../../../spec/lastenheft.md#lh-qa-03--portabilität): **Freigabe und Tag.** Ein Release-Commit trägt `VERSION` mit `0.2.0` und die Zeile `0.2.0` der Änderungshistorie des Handbuchs (drei Zellen, dritte nicht leer, kein weiteres `\|`, `releasing.md` §9). Der Eintrag `docs/maintainer/freigaben/0.2.0.md` (Kopie der Tabelle aus `releasing.md` §10, elf Punkte, dazu die Zeile des Jobs `homebrew` unter Punkt 8) trägt für genau diesen Hash die Belege 1 bis 6 (Punkte 2 und 3 nach [`CO-002`](../../carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md), Aktiv seit 2026-10-24) und das Ja des Nutzers im Wortlaut mit Datum (Punkt 7; der Implementer erfindet keines). Nach dem Tag, den die Koordination setzt, trägt der Implementer die Belege 8 bis 10 aus dem Log (`gh run view <Kennung> --log`, selbst gelesen, nicht übernommen): Lauf-Kennung, Index-Digest je Registry gleich, Smoke aus beiden Registries, `latest`, die sieben Anhänge, die Beobachtungen und, getrennt, der Job `homebrew` mit seinem Ausgang. Was der Lauf nicht zeigt, steht als Grenze im Eintrag, nie als Behauptung (`AGENTS.md` §3.11).
- [ ] [`LH-FA-01`](../../../../spec/lastenheft.md#lh-fa-01--kommandozeilenanwendung), [`LH-FA-16`](../../../../spec/lastenheft.md#lh-fa-16--container-eignung): **Bezug nach dem Tag** (Punkt 11, `SPEC-067` Punkt 18), in einem eigenen Commit nach dem Lauf, der nur `docs/user/`, `README.md` und Planung ändert (inert nach `SPEC-052` Punkt 12): Kopf des Handbuchs (*Software-Version*, *Gültigkeitsbereich*) nennt die Ausgabe von `version` des **veröffentlichten** Binaries (`pgwire-recorder 0.2.0`); §2 *Installation* (Binary vom GitHub-Release, Image aus beiden Registries) nennt `v0.2.0`, die sechs Namen ohne Version, die Prüfung mit `sha256sum -c`, die Tags `0.2.0` und `latest` mit dem Digest des Laufs — **nur mit Befehlen, die nach dem Tag einmal gelaufen sind**, belegt im Eintrag mit Befehl, Lauf, Ausgabe; `README.md` nennt `0.2.0` mit den Adressen im Ist-Zustand und schweigt über `0.1.1` (Zeile in der Änderungshistorie bleibt). Der Abschnitt zu Homebrew im Handbuch nennt für `0.2.0` nur, was `make homebrew-nachweis` nach dem Tag gezeigt hat (Linux/`amd64`); für macOS und `arm64` bleibt die Grenze.
- [ ] [`LH-FA-19`](../../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew): **Stand der Wartungs-Doku** (nicht inert, daher in einem eigenen Commit mit `make gates` danach): `docs/maintainer/releasing.md` §1 (*Stand*: jüngster Tag, Digest, Ausgang des Jobs `homebrew`, Verweis auf `freigaben/0.2.0.md`), Kopf (Version, Stand) und eine Zeile der Änderungshistorie des Dokuments. Dieser Punkt ist von Punkt 2 getrennt, damit der Bezug-Commit inert bleibt; bei `v0.1.1` stand er im Bezug-Commit und musste im Eintrag als Abweichung genannt werden.
- [ ] `make gates` grün auf dem getaggten Commit (Punkt 1 des Eintrags); die Gegenproben sind Punkte 2 und 3 des Eintrags: nicht voll gefahren, sondern nach [`CO-002`](../../carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md) mit den Ersatzbelegen aus dessen Text; der Eintrag schreibt je Punkt den Wortlaut des Carveouts „Was der Eintrag schreibt“ mit Befehl, Hash, Datum, Dauer, Last und Ausgang; nie „grün“ ohne Beleg.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure verbraucht den Carveout [`CO-002`](../../carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md) (Trigger: der Tag `v0.2.0` gesetzt und `freigaben/0.2.0.md` geschlossen, oder der Slice ohne Tag zurückgeführt, Grund in §7): Der Eintrag im Index `docs/plan/carveouts/README.md` zieht im Inhalts-Commit vor dem `git mv` nach `done/` nach.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung angefallen" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

**Stand vor dem Tag** (Planner, 2026-10-24, am Stand `0fe0e023`): `VERSION` ist `0.1.1`; die Änderungshistorie des Handbuchs endet bei `0.1.1`, eine Zeile `0.2.0` fehlt (ohne sie endet `Eingaben` im Lauf mit Exit 2, im Trockenlauf ist es ein Hinweis, der Punkt 5 der Checkliste verletzt); der Kopf des Handbuchs nennt `0.1.1` als *Software-Version* und im *Gültigkeitsbereich*; ein Tag `v0.2.0` gibt es nicht; `release.sh`, `release.yml` und die Vorlage der Formel sind seit `v0.1.1` unverändert bis auf eine Zeile Kommentar in `release.sh` (Pfad der Gegenprobe, Commit `3eca59e0`).

**Korrektur des Tag-Commits (2026-10-24, nach dem Trockenlauf `40071373625`):** Der Release-Commit `R` = `d720e44f` trägt den Tag nicht. Der Scan im Trockenlauf auf `R` war rot (`CVE-2026-111943`, Go-Standardbibliothek 1.27.2); `slice-go-toolchain-1-27-3` hob den Pin auf 1.27.3 (Commit `143e9994`), und `main` = `32d0d84a` enthält `R` als Vorfahr und den Pin. **Der Tag `v0.2.0` kommt auf `32d0d84a`.** Wo dieser Abschnitt, die Tabellen darunter und §6 für Punkt 1 bis 6, die Reihenfolge und die Tag-Handlung „`R`“ sagen, gilt der Hash `32d0d84a`; `R` bleibt der Commit mit `VERSION` und der Zeile `0.2.0`. Die Punkte 1 bis 6 sind an `32d0d84a` neu belegt (Eintrag; Punkt 1 mit zwei Läufen, Punkte 2 und 3 nach `CO-002` neu gefahren, Punkt 5 und 6 mit dem Trockenlauf `40072517299` auf `32d0d84a`). Der Trockenlauf steht damit auf dem Hash des Tags, `main` ist gleich diesem Hash (Randform (c)).

**Zeile `0.2.0` (Vorschlag des Planners; der Implementer prüft jeden Satz gegen das gebaute Binary, `AGENTS.md` §3.11):** Version `0.2.0`, Datum des Release-Commits, dritte Zelle: „`pgwire-recorder --version` und `-v` geben die Version aus wie `version`; `play --diff-all` meldet mit `--compare-responses` alle Abweichungen einer Interaktion statt der ersten“. Die Zeile hat genau drei Zellen und kein weiteres Pipe-Zeichen, auch nicht maskiert (`releasing.md` §9).

**Der Kopf des Handbuchs bleibt bei `0.1.1` bis nach dem Lauf** (Entscheidung des Planners, Grund: `SPEC-067` Punkt 18 lässt den Kopf die Ausgabe des **veröffentlichten** Binaries nennen, und vor dem Tag gibt es kein solches; eine Software-Version `0.2.0` im Kopf behauptete ein Binary, das noch nicht existiert, `AGENTS.md` §3.11). Der Release-Commit ändert den Kopf deshalb nicht, nur die Zeile der Änderungshistorie. Zwischen Release-Commit und Bezug-Commit beschreibt das Handbuch Verhalten, das erst `0.2.0` hat (`--diff-all`, `--version`), unter einem Kopf, der `0.1.1` nennt; das ist heute auf `main` schon so und endet mit dem Bezug-Commit (§6, Risiko 4).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `VERSION` | update | `0.2.0` im Release-Commit, der getaggt wird; ändert das Wurzelverzeichnis und ist damit nicht inert nach `SPEC-052` Punkt 12 (volle Kette) |
| `docs/user/benutzerhandbuch.md` §11 | update | Zeile `0.2.0` der Änderungshistorie im Release-Commit, vor dem Tag; Text und Maßgabe über der Tabelle |
| `docs/maintainer/freigaben/0.2.0.md` | neu | Eintrag der Freigabe, Kopie der Tabelle aus `releasing.md` §10, Belege je Punkt; Form wie `freigaben/0.1.1.md`; steht in einem Commit nach dem Release-Commit (er kann dessen Hash nicht tragen) |
| `docs/user/benutzerhandbuch.md` (Kopf, §2 *Installation*), `README.md` | update | Bezug nach dem Tag (Punkt 11), ersetzt den von `v0.1.1`; nur belegte Befehle |
| `docs/maintainer/releasing.md` (§1 *Stand*, Kopf, Änderungshistorie) | update | Zustand nach dem Tag, eigener Commit (DoD 3) |
| Lauf `release.yml` auf dem Tag (Trockenlauf davor) | Handlung | kein Dateiinhalt; Handlung der Koordination, nicht des Implementers; der Trockenlauf veröffentlicht nichts (`SPEC-067` Punkt 2) |

**Die elf Punkte der Checkliste, je Schritt, wer, womit** (Reihenfolge der Handlungen siehe darunter):

| Punkt | Wer | Schritt |
|---|---|---|
| 1. `make gates` | Implementer | `make gates` am Hash `R` des Release-Commits, sauberer Baum (`git status --short` leer), im Dev-Container oder Worktree; Hash, Datum, Exit, Dauer, Last (1-Minuten-Mittel, Kernzahl) im Eintrag |
| 2. Gegenproben-Container | Implementer, nach [`CO-002`](../../carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md) | **Ersatz (Festlegung im Carveout):** am Hash `R`, je mit Datum, Dauer, Last (1-Minuten-Mittel, Kernzahl), Ausgang: `GP_MUT='^ok$' make release-gegenprobe` (unveränderte Kopie; Muster wählt allein sie, ein Muster ohne Treffer endet mit Exit 1) und `GP_MUT='^kein-treffer$' make test-container-gegenprobe` (allein die unveränderte Kopie, keine Mutation gefahren). Der Ausgang der zweiten ist der Exit-Code und „unveränderte Kopie grün; keine Mutation gefahren“; ihre Schlusszeile ist kein Beleg (`AGENTS.md` §3.11). Ersatz 1 von `CO-001` (Ablage des Gegenproben-Containers) entfällt, den Container `pgr-gp` gibt es nicht mehr (Eintrag `0.1.1`, Punkt 2). Grenze im Eintrag: die Mutationen der Gegenproben und die übrigen Gegenproben am Hash sind nicht gelaufen |
| 3. `make gegenproben-host` | Implementer, nach [`CO-002`](../../carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md) | **Ersatz (Festlegung im Carveout):** `git diff --stat -M 35778a2f R -- tools/harness/container-start.sh tools/harness/container-start-gegenprobe.sh tools/dev/harness/container-start-gegenprobe.sh harness/mk/container-start.mk harness/mk/container-start-gegenprobe.mk harness/vertraege/SPEC-058-container-start.md` (letzte belegte grüne Fahrt `35778a2f`, 2026-10-19, Angabe ohne Log). Der Eintrag nennt die Ausgaben von `git diff -M --numstat` (Befehl wie oben) und `git diff -M -U0`. **Kriterium: nur Pfadzeilen.** Jede geänderte Zeile ist ein Pfad oder ein Kommentar auf einen Pfad; jede andere Änderung (Logik, Fall, Vertrag) ist kein Ersatz, der Nutzer entscheidet vor dem Tag. Dazu `make -n container-start-gegenprobe` am Hash `R`: die Ausgabe nennt `bash tools/dev/harness/container-start-gegenprobe.sh`, und die Datei liegt dort. Grenze im Eintrag: ein Lauf am Host am Hash `R` steht aus; der Diff belegt, dass sich Werkzeug und Gegenprobe nicht inhaltlich änderten, nicht, dass die Fahrt heute grün wäre |
| 4. `VERSION` und Zeile | Implementer | `git show R:VERSION` gibt `0.2.0`; `git show R:docs/user/benutzerhandbuch.md \| grep -n '^\| 0\.2\.0 \|'` findet die Zeile |
| 5. Trockenlauf | Koordination fährt, Implementer trägt ein | `workflow_dispatch` auf `main` **an dem Stand, auf dem `main` gleich `R` ist** (siehe Reihenfolge); Lauf-Kennung, Dauer, Zeile `release: eingaben: grün für v0.2.0 (Trockenlauf)`, die Zeile der Änderungshistorie ohne Hinweis auf eine fehlende Zeile; der Eintrag nennt, auf welchem Hash der Lauf stand (bei `v0.1.1` stand er auf einem späteren Stand als der Tag, das war eine Grenze, hier soll es keine sein) |
| 6. Scan | Implementer liest | `release: scan: 0 Befund(e) über der Schwelle, <m> ohne Behebung` im Log des Trockenlaufs und des Laufs des Tags; ein Befund über der Schwelle ist rot, kein Tag, Entscheidung des Nutzers (`SPEC-067` Punkt 16, keine Ausnahme) |
| 7. Ja des Nutzers | Nutzer, Koordination nennt | Wortlaut „release“, Datum 2026-10-24, und Version `v0.2.0`: **Entscheidung des Nutzers** vom 2026-10-24 (Kopf); die Wiedergabe stammt von der Koordination und ist nicht nachgefahren. Der Tag folgt erst nach grünem Trockenlauf und nochmaliger Bestätigung durch den Nutzer; sie trägt der Eintrag mit Datum, sobald sie vorliegt |
| 8. Lauf des Tags | Koordination setzt den Tag, Implementer liest das Log | Lauf-Kennung, Index-Digest je Registry gleich, Smoke, `latest` auf `D`, Zeilen wörtlich; darunter, getrennt, der Job `homebrew` (Lauf-Kennung, grün oder rot; bei grün Hash und Datum aus `git ls-remote https://github.com/pt9912/homebrew-pgwire-recorder.git refs/heads/main`; bei rot das Stichwort `release: homebrew: <Stichwort>:`); ein roter Job lässt Punkt 8 erfüllt, Behebung ist der Wiederanlauf oder der Nachtrag von Hand (`releasing.md` §7, §11), kein neuer Tag |
| 9. Anhänge | Implementer | `gh release view v0.2.0 --json assets`: genau sieben Dateien (sechs Binaries ohne Version im Namen, `SHA256SUMS`) |
| 10. Beobachtungen des Laufs | Implementer | je Beobachtung aus `releasing.md` §10 Lauf-Kennung und Zeile oder Grenze; **neu gegenüber dem ersten Tag:** `skopeo`-Wortlaut und `gh release view` für einen nicht vorhandenen Tag sind bei `v0.1.1` belegt und werden nur auf Abweichung gelesen; der Wiederanstoß eines roten **Laufs** bleibt Grenze, wenn der Lauf grün endet |
| 11. Bezug nach dem Tag | Implementer | DoD 2; jeder Befehl läuft nach dem Tag gegen die **veröffentlichten** Artefakte (Linux/amd64: `curl`, `sha256sum -c --ignore-missing`, `./pgwire-recorder-linux-amd64 version` gibt `pgwire-recorder 0.2.0`; `docker pull` ohne Anmeldung in beiden Registries für `0.2.0` und `latest`, gleicher Digest, `docker run … version`); zusätzlich ein Beleg, dass `--version` und `-v` des **veröffentlichten** Binaries dieselbe Zeile geben (`LH-FA-01`; sonst nennt das Handbuch es nicht); `darwin`, `windows` und `linux/arm64` sind gebaut, nicht ausgeführt, das Handbuch sagt nur das |

**Reihenfolge der Handlungen** (jede mit Rolle; Abweichungen nennt der Eintrag):

1. Architect: Festlegung von `CO-002` (§6, Rückgabe 1, erledigt). Nutzer: hat ihn am 2026-10-24 angenommen („deine empfehlungen“).
2. Implementer: Release-Commit `R` (`VERSION`, Zeile `0.2.0`), `make gates` an `R`, die Ersatzbelege 2 und 3 an `R`. Danach in getrennten Commits der Eintrag bis Punkt 7 (Punkte 8 bis 11 leer).
3. Koordination: `R` per ff auf `main` und `origin`, so dass `R` die Spitze von `main` ist; Trockenlauf auf diesem Stand. `Eingaben` verlangt, dass der Commit des Tags im Standard-Zweig liegt (`releasing.md` §6, Schritt 1); der Eintrag und die Planung folgen **hinter** `R` und dürfen den Tag nicht tragen.
4. Koordination: grüner Trockenlauf, Ja des Nutzers liegt vor und ist nach dem Trockenlauf nochmals bestätigt, `CO-002` ist Aktiv: `git tag v0.2.0 R`, `git push origin v0.2.0`. Roter Trockenlauf: kein Tag (§4).
5. Implementer: Belege 8 bis 10 aus dem Log, danach der Bezug-Commit (Punkt 11), danach der Commit zu `releasing.md`; Review, Verifikation, Closure.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Die Entscheidung des Nutzers vom 2026-10-24 liegt vor (Kopf); `slice-play-feldvergleich-doku` und `slice-version-option` liegen in `done/`; der Architect hat `CO-002` festgelegt, und der Nutzer hat ihn am 2026-10-24 angenommen (Aktiv). Die Vorbedingungen des Betreibers (`releasing.md` §4: Secrets beider Registries, `HOMEBREW_TAP_GITHUB_TOKEN`, Repository des Taps) sind Handlung des Nutzers; der Slice prüft ihre Existenz, nicht ihr Recht, und liest keinen Wert. WIP-Limit 1 je Rolleninhaber (`Verantwortlich`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next`: **zu groß**, wenn die Summe die **Marke 1500** erreicht (Metrik: `AGENTS.md` §3.14, dort der Befehl; Pfade des Plans: `VERSION docs/user docs/maintainer README.md`; Basis: der Commit des Übergangs nach `in-progress/`, Hash beim Übergang eintragen: ``193f0e0b``). Schnitt wäre (a) Release-Commit, Eintrag bis Punkt 7 und Tag, (b) Bezug nach dem Tag. Basis der Metrik bleibt `193f0e0b`.
- `in-progress` → `open`: (i) ein Ersatz aus `CO-002` ist vom Implementer nicht in Minuten zu fahren, oder die Ersatzbelege 2 oder 3 sind rot: Die Punkte 2 und 3 bleiben unerfüllt, kein Tag; ein Slice, der wartet, hält das WIP-Limit, deshalb geht er nach `open/`. (ii) Trockenlauf oder Scan sind rot aus einem Grund außerhalb dieses Slice (Befund mit Behebung, Datenbank, Netz): Entscheidung des Nutzers, kein Tag (`SPEC-067` Punkt 16). (iii) Die Koordination setzt den Tag nicht (das Ja trägt nicht mehr).
- `in-progress` → `next`: Der Lauf des Tags endet rot aus einem Grund im Skript, im Text der Datei oder im Vertrag (`SPEC-067` Punkt 18 (b)): Der Tag bleibt, die Korrektur ist ein Folge-Slice, ein weiterer Tag `v0.2.1` braucht ein **neues Ja des Nutzers** (das Ja vom 2026-10-24 nennt nur `v0.2.0`; ein Folge-Tag ist nicht eingeschlossen, anders als bei `v0.1.1`); dieser Slice schließt dann mit dem Eintrag für `v0.2.0` nach der Incident-Klausel.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig, Review-Report liegt vor, Closure-Notiz mit Lerneintrag geschrieben. Beobachtbar: der Lauf des Tags ist `success`, und der Eintrag `freigaben/0.2.0.md` trägt alle elf Punkte samt Zeile des Jobs `homebrew`, jede Zahl am Original nachgefahren.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

**Randformen** (`AGENTS.md` §3.12). Der Slice liefert keinen neuen Vertrag, sondern führt `SPEC-067` Punkt 17 und 18 ein drittes Mal aus (`v0.1.0`, `v0.1.1`); die Randformen stehen dort oder sind als Rückgabe benannt. Wo sie sich gegenüber `slice-v1-abschluss-release-v0-1-1` ändern, steht es hier.

| Randform | Entscheidung | Stelle |
|---|---|---|
| (a) Version | `0.2.0`, stabil (kein `-`): `latest` wandert von `0.1.1` auf `0.2.0`, der Job `homebrew` läuft; Tag `v0.2.0` gleich `v` plus Inhalt von `VERSION`, sonst endet `Eingaben` mit Exit 2 | Entscheidung des Nutzers vom 2026-10-24 (Kopf); `SPEC-067` Punkte 3, 8, 19 (2) |
| (b) Wer taggt, wann | Koordination, nach dem Ja („release“); der Implementer setzt keinen Tag, startet keinen Lauf und pusht nicht; Übergabe mit dem Hash `32d0d84a` (Korrektur in §3), und die Koordination taggt genau ihn | `SPEC-067` Punkt 18; `releasing.md` §10 |
| (c) Reihenfolge der Commits | **Der Tag-Commit (`32d0d84a`, `R` plus Pin des Go-Builders) ist die Spitze von `main`, wenn der Trockenlauf und der Tag laufen** (Korrektur in §3: `R` trug ihn wegen des Scans nicht); danach Eintrag, Planung, Bezug-Commit, Commit zu `releasing.md`. Neu gegenüber `v0.1.1`: Dort lief der Trockenlauf auf einem späteren Stand als der Tag (Grenze im Eintrag), hier steht er auf dem Hash des Tags (Lauf `40072517299`) | `SPEC-067` Punkte 17, 18; `releasing.md` §5, §10 |
| (d) Rot nach dem Tag | Umgebung (Token, Schreibrecht, Sichtbarkeit, Netz): Wiederanstoß desselben Tags, Handlung des Betreibers; Skript, Text der Datei oder Vertrag: Korrektur auf `main`, ein weiterer Tag braucht ein **neues Ja des Nutzers** und einen neuen Carveout oder eine Verlängerung von `CO-002` durch ihn; Digest oder Anhang abweichend: immer neuer Tag. **Folge-Tag: nein, entschieden** (Entscheidung des Nutzers vom 2026-10-24, „deine empfehlungen“: `CO-002` deckt allein `v0.2.0`, `v0.2.1` ist nicht eingeschlossen; anders als bei `v0.1.1`, dessen Ja `v0.1.2` einschloss) | `SPEC-067` Punkte 11, 18; `CO-002` Geltungsbereich |
| (e) Job `homebrew` | Ein roter Job ändert weder Release noch Registries; der Slice nennt sein Ergebnis im Eintrag, Wiederanlauf oder Nachtrag von Hand (`releasing.md` §7) sind kein neuer Tag und Handlung des Betreibers | `SPEC-067` Punkt 19 (1), (7) |
| (f) Alte Tags | `v0.1.1` und `v0.1.0` bleiben; der Bezug nennt nur `v0.2.0` | `SPEC-067` Punkt 18, „Jeder weitere Tag“ |
| (g) `latest` | wandert auf `0.2.0`, wenn der Vergleich der Digests beider Registries grün und `0.2.0` höher ist als jede andere stabile Version unter den Tags `v*`; der Image-Tag `0.1.1` bleibt | `SPEC-067` Punkt 8 |
| (h) Plattform | Linux/amd64 im Lauf und im Bezug geprüft (Binary und Image gestartet); `darwin`, `windows`, `linux/arm64` gebaut, nicht ausgeführt: Grenze, die das Handbuch so nennt; die Angabe des Nutzers, `brew install` auf macOS habe für `0.1.1` funktioniert, ist eine Angabe, kein Beleg (`SPEC-067` Punkt 20 (6)) und kommt nicht in das Handbuch | `SPEC-067` Punkt 15; `LH-FA-19.a`; `LH-QA-03` |
| (i) Aufrufarten, Abbruch | entfällt: kein Werkzeug, kein Ergebnis, das ein Hook liest; Treiber ist `release.sh`, nie unter `make` | `SPEC-067` Punkt 13 |
| (j) Ort des Eintrags | `docs/maintainer/freigaben/0.2.0.md` | `SPEC-067` Punkt 17 |
| (k) Punkte 2 und 3 | nach [`CO-002`](../../carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md) (Aktiv seit 2026-10-24): Ersatz in §3, Punkte 2 und 3; die Pflicht bleibt im Vertrag (`SPEC-060` Punkt 14). `CO-001` trägt `v0.1.1` und `v0.1.2`, nennt `v0.2.0` nicht und bleibt unverändert als Träger der Ursache | `docs/plan/carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md`; `AGENTS.md` §3.6 |
| (l) Der Kopf des Handbuchs | bleibt `0.1.1` bis zum Bezug-Commit (§3) | `SPEC-067` Punkt 18 |

**Rückgabe 1 an den Architect (erledigt):** `CO-001` war für `v0.2.0` nicht tragfähig. Der Architect hat `CO-002` festgelegt, der Nutzer hat ihn am 2026-10-24 angenommen („deine empfehlungen“): neues `CO-002`, Ersatz 1 von `CO-001` entfällt, Ersatz 2 mit den Mustern aus §3 Punkt 2, Ersatz 3 mit dem Kriterium „nur Pfadzeilen“ (§3 Punkt 3), kein dritter Ersatzbeleg, kein Folge-Tag. Den Auflösungs-Trigger von `CO-002` führt dieser Slice (Closure verbraucht ihn, DoD); die Nehmer der Ursache (`slice-harness-gegenproben-schneller`, `slice-harness-gegenproben-grenze-rest`) trägt allein `CO-001`, `CO-002` sendet ihnen nichts.

**Entscheidungen des Nutzers vom 2026-10-24:** Form, Wortlaut, kein dritter Ersatzbeleg und kein Folge-Tag sind mit „deine empfehlungen“ entschieden (Kopf); das Ja zum Tag trägt den Wortlaut „release“. Offen bleibt allein die nochmalige Bestätigung vor dem Tag (Koordination).

**Risiken:**

- Die Ersatzbelege der Punkte 2 und 3 aus `CO-002` sind rot oder nicht in Minuten zu fahren (der volle Lauf der Gegenproben übersteigt die Zeitgrenze, [ADR-0055](../../adr/0055-gegenproben-haben-ein-zeitbudget.md)); der Diff des Ersatzes 3 ist nicht leer, trägt nach dem Kriterium des Carveouts aber nur Pfadzeilen — **Ausgang:** offen bis zur Closure.
- Das Tap-Token oder ein Registry-Geheimnis ist seit `v0.1.1` abgelaufen oder in seinem Recht geändert (ein Fine-grained Token trägt ein Ablaufdatum, `releasing.md` §4); die Checkliste prüft nur die Existenz (`BEO-REPO/vorbedingung-geheimnis-nur-existenz-geprueft`, Zähler 2, mit diesem Lauf als drittem Auftreten möglich) — **Ausgang:** offen bis zur Closure.
- Die Zeile `0.2.0` der Änderungshistorie sagt mehr, als das Binary kann (`AGENTS.md` §3.11) — **Ausgang:** offen; Gegenmittel: jeder Satz der Zeile wird mit einem Lauf des gebauten Binaries belegt (`--version`, `-v`, `play --diff-all` mit und ohne `--compare-responses`) und der Beleg steht im Eintrag.
- Der Handbuch-Kopf nennt `0.1.1`, während der Text Verhalten von `0.2.0` beschreibt, bis der Bezug-Commit liegt; ein Leser, der das Binary `0.1.1` hat, findet `--diff-all` und `--version` nicht — **Ausgang:** offen; der Zustand besteht auf `main` schon heute und endet mit dem Bezug-Commit; wird der Tag nicht gesetzt (Rückführung), bleibt er bestehen und geht an den Planner.
- `docs/maintainer/` steht nicht auf der Positivliste von `SPEC-052` Punkt 12: Ein Commit an `releasing.md` löst die volle Kette aus; bei `v0.1.1` stand die Änderung im Bezug-Commit und musste als Abweichung genannt werden — **Ausgang:** entfallen durch den Schnitt in DoD 2 und 3 (zwei Commits), die Zählung dazu steht in §8.
- Der Scan ist rot durch einen alternden Pin (der Trockenlauf auf `R`, Lauf `40071373625`, `CVE-2026-111943`, Go-Standardbibliothek 1.27.2) — **Ausgang:** eingetreten → `slice-go-toolchain-1-27-3` (hob den Pin auf 1.27.3, Commit `143e9994`; der Trockenlauf `40072517299` auf `32d0d84a` ist grün: `release: scan: 0 Befund(e) über der Schwelle, 0 ohne Behebung`); die Beobachtung führt `BEO-REPO/release-pipeline-pin-und-werkzeug-altern-ohne-melder`, deren Datei der Pin-Slice anlegt. Rückführung (ii) aus §4 ist nicht ausgelöst: der Grund lag außerhalb dieses Slice, ist behoben, und der Slice bleibt in `in-progress/`.
- Ein Fehler im Skript, den erst der Lauf des Tags zeigt, kostet einen Folge-Tag, den das Ja nicht deckt (Randform (d)) — **Ausgang:** offen bis zur Closure.

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

(noch nicht geschlossen)

**Belege des Implementers (Stand nach dem Tag, 2026-10-24).** Die Punkte 7 bis 11 stehen in [`freigaben/0.2.0.md`](../../../maintainer/freigaben/0.2.0.md); hier nur Lauf und Messung. Ja des Nutzers zum Tag, nochmals bestätigt am 2026-10-24, Wortlaut „ja“ (Punkt 7). Tag `v0.2.0` auf `32d0d84a` (Tag-Objekt `fa5de159`), Lauf `40074219261` (`push`, 2026-10-24 07:27:44Z bis 07:32:11Z) `success` im ersten Versuch: Job `release` 4 min 14 s, Job `homebrew` 9 s, Formel `0.2.0` im Tap (`d96e2968`, 07:32:07Z); Index-Digest `sha256:1d50181235d6a0317aaf376c33c1f17b421fb25036dea3324946a652dd2ceee4` in beiden Registries gleich und gleich dem des Trockenlaufs, Scan 0/0, `latest` auf demselben Digest, sieben Anhänge, veröffentlicht 07:31:59Z. Bezug (Punkt 11): `da5fe8b6` ändert nur `docs/user/benutzerhandbuch.md` und `README.md`; Befehle am veröffentlichten Binary und an den Images gelaufen (`version`, `--version`, `-v` je `pgwire-recorder 0.2.0`; `sha256sum -c --ignore-missing` OK; `docker pull` ohne Anmeldung, vier Referenzen, ein Digest) und `make homebrew-nachweis` Exit 0 (7 min 46 s, `brew install` aus dem Tap gibt `0.2.0`). Größenmetrik (`AGENTS.md` §3.14, Basis `193f0e0b`, Pfade `VERSION docs/user docs/maintainer README.md`) am committeten Stand `da5fe8b6`: 75 Zeilen (Marke 1500). Der Commit zu `releasing.md` (DoD 3) und `make gates` am Endstand folgen unten. **Plan folgt Korrektur (§3.9):** §3 sagt „Der Kopf des Handbuchs bleibt bei `0.1.1` bis nach dem Lauf“ und die Zeile zu Punkt 11 „nur Befehle, die gelaufen sind“; beides hält der Bezug-Commit ein. §1 schloss einen neuen Homebrew-Nachweis für `0.2.0` aus („anderer Vorgang“); der Implementer hat `make homebrew-nachweis` trotzdem gefahren, weil der Abschnitt zu Homebrew im Handbuch sonst `0.1.1` als Ergebnis eines Laufs nannte, den es für `0.2.0` nicht gab (Auftrag der Koordination: Abschnitt nur mit Belegtem). Der Beleg steht in Punkt 11; der Nachweis bleibt ein Werkzeug, kein Gate, und der Slice `slice-erster-release-homebrew-nachweis` bleibt für die übrigen Plattformen zuständig.

**Belege des Implementers (Stand vor dem Tag, 2026-10-24).** Die Belege der Punkte 1 bis 7 der Checkliste stehen in [`freigaben/0.2.0.md`](../../../maintainer/freigaben/0.2.0.md), nicht hier. Tag-Commit `32d0d84a` (Korrektur in §3), enthält `R` = `d720e44f` (`VERSION` `0.2.0`, Zeile `0.2.0`) und den Pin `143e9994`; Basis der Größenmetrik `193f0e0b`, gemessen am committeten Stand vor dem Eintrag: 3 Zeilen (Marke 1500). Läufe: `make gates` Exit 0 an `143e9994` (528 s) und an `32d0d84a`; `GP_MUT='^ok$' make release-gegenprobe` Exit 0 (160 s) und `GP_MUT='^kein-treffer$' make test-container-gegenprobe` Exit 0 (26 s) an `32d0d84a` nach `CO-002`; Messung des Diffs der sechs Pfade von `35778a2f` bis `32d0d84a`: 3 Dateien, 6/6 Zeilen, gleich dem Befund an `R`, davon eine Zeile Code der Pfadberechnung (`../..` wird `../../..` in `container-start-gegenprobe.sh`), Urteil des Nutzers dazu: am 2026-10-24 als Pfadzeile gedeckt akzeptiert, Wortlaut „ja“; Trockenlauf `40072517299` auf `32d0d84a` grün (Scan 0/0). Mutationen: dieser Slice ändert weder Code noch Gate noch Vertrag (§1); es gibt keine Zusage an einer Wire-, Datei- oder Formatgrenze und keinen Wächter, für den eine Mutation zu fahren wäre (`AGENTS.md` §3.10).

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `REPO` (gesamtes Repo, `harness/conventions.md`); die Schwelle ≥ 2 von 3 Achsen ist nicht berührt.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-REPO/` durchgegangen (2026-10-24, gemergter Stand). Treffer: `vorbedingung-geheimnis-nur-existenz-geprueft` (Zähler 2, zwei Dateien unter `evidence/`; erreicht der Lauf des Tags ein drittes Auftreten, ist der Eintrag eine Lücke und braucht einen Folge-Slice, §6 Risiko 2), `wiederanstoss-eines-roten-release-laufs-unbelegt` (Zähler 2; ein grüner Lauf des Tags ist ein drittes Auftreten dieser Grenze, wenn der Wiederanstoß wieder nicht gefahren wird: die Zählung beim Start lesen), `plan-folgt-korrektur-nicht` (`AGENTS.md` §3.9; der Plan folgt jeder Korrektur, §7 nennt die Stände). Der Zähler folgt aus den Dateien; der Planner setzt keinen.

**Zählung** (Planner, 2026-10-24, nachgezählt mit der Eintragung von [`CO-002`](../../carveouts/CO-002-gegenproben-am-getaggten-commit-v0-2-0.md), `AGENTS.md` §3.13 und §3.14): **Liefer-Punkte: 3** (1: Freigabe und Tag mit `VERSION`, Zeile der Änderungshistorie und Eintrag; 2: Bezug nach dem Tag in Handbuch und `README.md`; 3: Stand in `releasing.md`). **Schichten: zwei** — Release (Datei `VERSION`, Lauf der Pipeline als Handlung der Koordination) und Nutzer- und Wartungs-Doku (`docs/user/`, `docs/maintainer/`, `README.md`); Planung zählt nicht, Produkt- und Werkzeug-Code berührt der Slice nicht (§1). **Sendung `CO-002`:** Die Kennung und der Trigger „Closure verbraucht den Carveout“ stehen in §1 und DoD; sie sind eine Closure-Pflicht (Index nachziehen, `git mv` des Carveouts), kein Liefer-Punkt, und berühren keine dritte Schicht (`docs/plan/carveouts/` ist Planung). Die Ersatzbelege 2 und 3 sind Teil von Liefer-Punkt 1 (Belege des Eintrags), nicht ein vierter. Mit der Sendung bleiben es drei Punkte und zwei Schichten, die Grenzen sind gehalten. Gegenüber `slice-v1-abschluss-release-v0-1-1` (zwei Punkte) ist `releasing.md` ein eigener Punkt, damit der Bezug-Commit inert bleibt. **Nehmer der Ursache:** `slice-harness-gegenproben-schneller` und `slice-harness-gegenproben-grenze-rest` werden nicht angefasst, `CO-002` sendet ihnen nichts (sie nehmen `CO-001` an).

**Nachzählung der Korrektur** (Implementer, 2026-10-24, `AGENTS.md` §3.9 und §3.13): Die Korrektur des Tag-Commits (§3, §6) fügt weder einen Liefer-Punkt noch eine Schicht hinzu: Es bleiben drei Punkte und zwei Schichten. Der Pin-Slice `slice-go-toolchain-1-27-3` ist Ausgang eines Risikos (§6), keine Sendung an einen Nehmer; er liegt in `in-progress/` und nennt diesen Slice in seinem §1 als Halter von Tag und Eintrag.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (`SPEC-067` und die Entscheidungen gehen dem Eintrag voraus).
