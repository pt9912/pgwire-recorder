# Welle welle-gegenproben-und-mutanten: Gegenproben der Gates und Mutationstests

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug und kein Abnahmeszenario; die Welle ist Harness-Arbeit (Messung 1 und 3 von [LH-QA-07](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) berührt nur `slice-harness-gegenprobe-baumkopie`).

**Verantwortlich:** pt9912. **Datum:** 2026-10-18.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die Tests und Gegenproben des Bestands halten, was sie zusagen: Die Kopie des Arbeitsbaums in den Gegenproben bricht bei laufenden Agenten nicht (`slice-harness-gegenprobe-baumkopie`), die bekannten überlebenden Mutanten, die das Verhalten ändern, haben Tests (`slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`), Warten ohne Frist in Tests meldet der Lint (`slice-harness-lint-warten-ohne-frist`), und ein Mutations-Gate meldet rot, wenn zu viele Mutanten die Tests überleben (`slice-harness-mutation`). Der Nutzer hat am 2026-10-18 entschieden, diese bis dahin wellenlosen Slices zu einer Welle zu gliedern, und am selben Tag, die Wellen 2 und 3 entlang der Harness-Reihe neu zu schneiden: Diese Welle ist die zweite der drei, hinter [welle-harness-regeln](welle-harness-regeln.md) und vor [welle-lh-links](welle-lh-links.md). Die Gegenproben verlassen am 2026-10-18 auf Entscheidung des Nutzers `make gates` (`slice-harness-gegenproben-aus-gates`, `-sammelziele`, `-fahrstelle` und `-fahrstelle-betrieb`, vor der Baumkopie, weil der Lauf zu lange dauert). Die Gegenprobe der Versionsmatrix (`slice-harness-gegenprobe-matrix-p7-flanke`) liegt seitdem in welle-harness-regeln. Auf Auftrag des Nutzers vom 2026-10-18 kommen `slice-harness-nachweis-uebernehmen` und `slice-harness-nachweis-uebernehmen-aufrufarten` dazu (geschnitten am 2026-10-18; Teil 2 trägt die Signale, die übrigen Aufrufarten und die Kopplung an den Hook) (der Nachweis des vollen Laufs reist aus einem Worktree oder Dev-Container in den Hauptbaum, damit der Stop-Hook dort nicht sperrt); beide sind auf Entscheidung des Nutzers vom 2026-10-18 aus der Reihenfolge der Wellen gelöst und starten nach dem Kettenumbau bzw. nach Teil 1. Es gibt kein Abnahmeszenario und keinen Meilenstein, der an ihr hängt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [welle-harness-regeln](welle-harness-regeln.md) ist `done` (Neuschnitt der Wellen 2 und 3, Entscheidung des Nutzers vom 2026-10-18; die Wellen laufen strikt nacheinander). **Ausgenommen** sind `slice-harness-gegenproben-aus-gates` (Start-Trigger: `slice-harness-devcontainer-einrichten-geruest` in `done/`), `slice-harness-gegenproben-aus-gates-sammelziele` (Start-Trigger: `slice-harness-gegenproben-aus-gates` in `done/`) `slice-harness-gegenproben-aus-gates-fahrstelle` (Start-Trigger: `slice-harness-gegenproben-aus-gates-sammelziele` in `done/`) `slice-harness-gegenproben-aus-gates-fahrstelle-umgebung` (Start-Trigger: `slice-harness-gegenproben-aus-gates-fahrstelle` in `done/`) und `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` (Start-Trigger: `slice-harness-gegenproben-aus-gates-fahrstelle-umgebung` in `done/`): Alle fünf sind auf Auftrag des Nutzers („Bitte sofort beheben“) aus der Reihenfolge der Wellen gelöst und starten vor dem Start der Welle (Drift-Log der Roadmap, 2026-10-18); dieser Trigger gilt für die übrigen Slices. Damit sind auch [welle-erster-release](done/welle-erster-release/welle-erster-release.md) und [welle-v1-nachtrag](welle-v1-nachtrag.md) `done` (Wellen vor Harness, Entscheidung des Nutzers vom 2026-10-08).
- `slice-harness-nachweis-uebernehmen` (Start-Trigger: `slice-harness-gegenproben-aus-gates` in `done/`, weil seine Gegenprobe an die Liste der Namensregel hängt) und `slice-harness-nachweis-uebernehmen-aufrufarten` (Start-Trigger: `slice-harness-nachweis-uebernehmen` in `done/`) sind ebenfalls aus der Reihenfolge der Wellen gelöst (Entscheidung des Nutzers vom 2026-10-18, so vom Koordinator übermittelt: aus den Wellen lösen, Start nach dem Kettenumbau; Drift-Log der Roadmap) und starten vor dem Start der Welle.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Das *Mehr* gegenüber den Slice-DoDs (Entscheidung des Nutzers vom 2026-10-18): Ein Lauf des Mutations-Gates nennt keinen der dort geführten Mutanten (`slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`) als überlebend; der Lauf steht in der Closure-Notiz. Die DoDs der Sammel-Slices belegen ihre Tests je für sich, die DoD von `slice-harness-mutation` die Schwelle, nicht ihr Zusammenspiel.
- Für `slice-harness-gegenprobe-baumkopie` und `slice-harness-lint-warten-ohne-frist` gibt es **kein weiteres gemeinsames *Mehr***: Ihre Bündelung ist Thema (Gegenproben und Güte der Tests), keine Closure-Bedingung.
- `make gates` grün auf dem Endstand, dazu ein grünes Ergebnis des Gegenproben-Containers (`make gegenproben`; fehlt es, fährt die Closure es selbst) und `make gegenproben-host` (Host) grün auf demselben Commit (Gegenproben verlassen `make gates` mit `slice-harness-gegenproben-aus-gates`, die Sammelziele kommen mit `slice-harness-gegenproben-aus-gates-sammelziele`; Namen laut `SPEC-060`).
- Closure-Notiz in `welle-gegenproben-und-mutanten-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-harness-gegenproben-aus-gates | Die Gegenproben und `a-check-negativ` verlassen `make gates` und hängen an zwei Listen (Kettenumbau) | [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (nur die Messungen, die Gegenproben tragen); kein Anspruch des Produkts |
| slice-harness-gegenproben-aus-gates-sammelziele | Zwei Sammelziele fahren die Gegenproben, `make gegenproben` im Gegenproben-Container, `make gegenproben-host` am Host | [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (nur die Messungen, die Gegenproben tragen); kein Anspruch des Produkts |
| slice-harness-gegenproben-aus-gates-fahrstelle | Hook für die Einlieferung in den Gegenproben-Container: `tools/harness/gegenproben-nach-einlieferung.sh` mit Gegenprobe | [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (nur die Messungen, die Gegenproben tragen); kein Anspruch des Produkts |
| slice-harness-gegenproben-aus-gates-fahrstelle-umgebung | Der Lauf des Hooks erbt keine Deskriptoren über 2 und ignoriert `SIGINT` nicht; `make` und seine Nachkommen ignorieren weder `SIGINT` noch `SIGQUIT`, mit Fällen und Mutationen in der Gegenprobe | [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (nur die Messungen, die Gegenproben tragen); kein Anspruch des Produkts |
| slice-harness-gegenproben-aus-gates-fahrstelle-betrieb | Pflicht und Betrieb der Gegenproben: Rollen-Anweisungen, Abschnitt über den Gegenproben-Container in `docs/developer/devcontainer.md` und die erste Fahrt | [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (nur die Messungen, die Gegenproben tragen); kein Anspruch des Produkts |
| slice-harness-nachweis-uebernehmen | Ein Werkzeug übernimmt den Nachweis des vollen Laufs aus einem Worktree oder Dev-Container, wenn die Dateiliste des Baums mit der des Nachweises übereinstimmt | [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (nur Messmethode (4), Werkzeuge der Prüfumgebung); kein Anspruch des Produkts |
| slice-harness-nachweis-uebernehmen-aufrufarten | Die Gegenprobe des Nachweis-Werkzeugs trägt die Signale, die übrigen Make-Flags, den Aufruf aus einem Unterverzeichnis und die Kopplung an den Hook | [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (nur Messmethode (4), Werkzeuge der Prüfumgebung); kein Anspruch des Produkts |
| slice-harness-gegenprobe-baumkopie | Gegenproben kopieren den Arbeitsbaum nach `git ls-files` statt mit `tar` und `cp -R`, durch einen gemeinsamen Helfer | [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (nur Messung 1 und 3 der Gegenproben); kein Anspruch des Produkts |
| slice-tests-ueberlebende-mutanten | Tests für überlebende Mutanten, die das Verhalten ändern | [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-10`](../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) |
| slice-tests-ueberlebende-mutanten-driving | Tests für überlebende Mutanten im PGWire-Adapter und im Bootstrap | [`LH-FA-03`](../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) |
| slice-harness-lint-warten-ohne-frist | Lint meldet Warten ohne Frist in Tests | kein Anspruch des Produkts |
| slice-harness-mutation | Mutationstests als Gate | kein Anspruch des Produkts |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-lh-links](welle-lh-links.md) (deren Start-Trigger verlangt diese Welle in `done/`); auf Slice-Ebene steht `slice-harness-lh-links-pflicht` hinter `slice-harness-mutation` (Drift-Log der Roadmap, 2026-10-09), das ergibt sich jetzt aus der Reihenfolge der Wellen.
- Wird blockiert von: Welle [welle-harness-regeln](welle-harness-regeln.md) (Start-Trigger; die Wellen laufen strikt nacheinander). Die Bindung auf Slice-Ebene bleibt im Plan von `slice-tests-ueberlebende-mutanten` stehen (er setzt `slice-harness-commit-struktur-id` voraus, den letzten Slice jener Welle) und ist mit dem Start dieser Welle erfüllt, außer der Slice dort ist ausdrücklich zurückgestellt.
- Innerhalb der Welle: `slice-tests-ueberlebende-mutanten-driving` setzt `slice-tests-ueberlebende-mutanten` voraus (Sammelregel), `slice-harness-lint-warten-ohne-frist` den zweiten Sammel-Slice, `slice-harness-mutation` alle drei davor; `slice-harness-gegenprobe-baumkopie` steht vor `slice-harness-lint-warten-ohne-frist` (beide ändern `tools/harness/lint-gegenprobe.sh`, die Baumkopie behebt die Flake der Läufe; Reihenfolge, kein Start-Trigger).

**Reihenfolge** (nach den Start-Triggern der Slices; WIP-Limit 1; die Welle ordnet nicht neu; Reihenfolge der drei Wellen strikt: [welle-harness-regeln](welle-harness-regeln.md), diese, [welle-lh-links](welle-lh-links.md), Entscheidung des Nutzers vom 2026-10-18):

1. `slice-harness-gegenproben-aus-gates`, danach `slice-harness-gegenproben-aus-gates-sammelziele`, danach `slice-harness-gegenproben-aus-gates-fahrstelle`, danach `slice-harness-gegenproben-aus-gates-fahrstelle-umgebung` (Hook: Deskriptoren und Signale, Folge der ersten echten Fahrt vom 2026-10-18), danach `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` — erste bis fünfte Stelle, vor `slice-harness-gegenprobe-baumkopie` (Entscheidung des Nutzers vom 2026-10-18, sofort): Der Nutzer will die Gegenproben aus `make gates` nehmen, weil der Lauf zu lange dauert; der zweite liefert die Sammelziele, der dritte den Hook für die Einlieferung mit Gegenprobe, der vierte Pflicht, Beschreibung und die erste Fahrt, je direkt dahinter (Schnitt des ersten Slice und Schnitt des dritten vom 2026-10-18, Schätzung der ungeteilten Slice über der Marke). Alle vier stehen vor der Baumkopie, weil diese ihre Gegenproben dann nur noch in `make gegenproben` fährt (DoD und Closure-Trigger dort nachgezogen). Alle vier Slices sind aus der Reihenfolge der Wellen gelöst (Start-Trigger in §2) und laufen vor dem Start der Welle; die Reihenfolge hier gilt, sobald die Welle läuft, und ändert nichts an ihrer Stelle vor der Baumkopie.
   Sechste und siebte Stelle, hinter `slice-harness-gegenproben-aus-gates-fahrstelle-betrieb` und vor der Baumkopie: `slice-harness-nachweis-uebernehmen`, danach `slice-harness-nachweis-uebernehmen-aufrufarten` (Schnitt am 2026-10-18; beide aus der Reihenfolge der Wellen gelöst, Entscheidung des Nutzers vom 2026-10-18, Start-Trigger in §2). Er setzt `slice-harness-gegenproben-aus-gates` voraus; ob er vor `slice-harness-gegenproben-aus-gates-fahrstelle` rückt, hat der Nutzer nicht gesagt, und nichts bindet ihn daran (Plan von `slice-harness-nachweis-uebernehmen`, §6 Frage 1).
2. `slice-harness-gegenprobe-baumkopie` — Stelle vor `slice-harness-lint-warten-ohne-frist` (Entscheidung des Nutzers vom 2026-10-18): Beide ändern `tools/harness/lint-gegenprobe.sh`, und die Baumkopie behebt die Flake der Läufe (`tar` bricht bei Änderungen im Arbeitsbaum). Vorbedingung des Starts ist der Vertrag des Architect (§4 des Slice).
3. `slice-tests-ueberlebende-mutanten`, danach `slice-tests-ueberlebende-mutanten-driving`, danach `slice-harness-lint-warten-ohne-frist`, danach `slice-harness-mutation` (Reihe der Harness-Slices, Entscheidung des Nutzers vom 2026-10-08).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Neue Produkt-Anforderungen oder Produktcode — die Welle ändert Tests, Gates und Gegenproben, nicht das Produkt.
- Neue Tests für Mutanten, die erst die Messung des Mutations-Gates findet — eigene Arbeit mit eigenem Slice nach `slice-harness-mutation` (§1 dort).
- Mutationen in den Integrationstests — anderer Vorgang (§1 von `slice-harness-mutation`).
- Die Regel-Slices der Harness-Reihe und die Gegenprobe der Versionsmatrix (`slice-harness-gegenprobe-matrix-p7-flanke`) — [welle-harness-regeln](welle-harness-regeln.md).
- Lastenheft- und MR-Kennungen als Links und deren Pflicht im Doku-Gate — [welle-lh-links](welle-lh-links.md).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-gegenproben-und-mutanten-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
