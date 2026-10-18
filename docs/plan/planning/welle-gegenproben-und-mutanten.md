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

Die Tests und Gegenproben des Bestands halten, was sie zusagen: Die Kopie des Arbeitsbaums in den Gegenproben bricht bei laufenden Agenten nicht (`slice-harness-gegenprobe-baumkopie`), die bekannten überlebenden Mutanten, die das Verhalten ändern, haben Tests (`slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`), Warten ohne Frist in Tests meldet der Lint (`slice-harness-lint-warten-ohne-frist`), und ein Mutations-Gate meldet rot, wenn zu viele Mutanten die Tests überleben (`slice-harness-mutation`). Der Nutzer hat am 2026-10-18 entschieden, diese bis dahin wellenlosen Slices zu einer Welle zu gliedern, und am selben Tag, die Wellen 2 und 3 entlang der Harness-Reihe neu zu schneiden: Diese Welle ist die zweite der drei, hinter [welle-harness-regeln](welle-harness-regeln.md) und vor [welle-lh-links](welle-lh-links.md). Die Gegenproben verlassen am 2026-10-18 auf Entscheidung des Nutzers `make gates` (`slice-harness-gegenproben-aus-gates` und `-fahrstelle`, vor der Baumkopie, weil der Lauf zu lange dauert). Die Gegenprobe der Versionsmatrix (`slice-harness-gegenprobe-matrix-p7-flanke`) liegt seitdem in welle-harness-regeln. Es gibt kein Abnahmeszenario und keinen Meilenstein, der an ihr hängt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [welle-harness-regeln](welle-harness-regeln.md) ist `done` (Neuschnitt der Wellen 2 und 3, Entscheidung des Nutzers vom 2026-10-18; die Wellen laufen strikt nacheinander). Damit sind auch [welle-erster-release](welle-erster-release.md) und [welle-v1-nachtrag](welle-v1-nachtrag.md) `done` (Wellen vor Harness, Entscheidung des Nutzers vom 2026-10-08).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Das *Mehr* gegenüber den Slice-DoDs (Entscheidung des Nutzers vom 2026-10-18): Ein Lauf des Mutations-Gates nennt keinen der dort geführten Mutanten (`slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`) als überlebend; der Lauf steht in der Closure-Notiz. Die DoDs der Sammel-Slices belegen ihre Tests je für sich, die DoD von `slice-harness-mutation` die Schwelle, nicht ihr Zusammenspiel.
- Für `slice-harness-gegenprobe-baumkopie` und `slice-harness-lint-warten-ohne-frist` gibt es **kein weiteres gemeinsames *Mehr***: Ihre Bündelung ist Thema (Gegenproben und Güte der Tests), keine Closure-Bedingung.
- `make gates` grün auf dem Endstand, dazu `make gegenproben` (Gegenproben-Container) und `make gegenproben-host` (Host) grün auf demselben Commit (Gegenproben verlassen `make gates` mit `slice-harness-gegenproben-aus-gates`; Namen laut Vertrag dort).
- Closure-Notiz in `welle-gegenproben-und-mutanten-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-harness-gegenproben-aus-gates | Die Gegenproben und `a-check-negativ` verlassen `make gates`; zwei Sammelziele fahren sie, eines im Gegenproben-Container, eines am Host | [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (nur die Messungen, die Gegenproben tragen); kein Anspruch des Produkts |
| slice-harness-gegenproben-aus-gates-fahrstelle | Pflicht und Fahrstelle der Gegenproben: Rollen-Anweisungen und der Gegenproben-Container in `docs/developer/devcontainer.md` | kein Anspruch des Produkts |
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

1. `slice-harness-gegenproben-aus-gates`, danach `slice-harness-gegenproben-aus-gates-fahrstelle` — erste und zweite Stelle, vor `slice-harness-gegenprobe-baumkopie` (Entscheidung des Nutzers vom 2026-10-18, sofort): Der Nutzer will die Gegenproben aus `make gates` nehmen, weil der Lauf zu lange dauert; der zweite Slice liefert Pflicht und Fahrstelle direkt dahinter. Beide stehen vor der Baumkopie, weil diese ihre Gegenproben dann nur noch in `make gegenproben` fährt (DoD und Closure-Trigger dort nachgezogen). Hinweis: Die Welle startet erst hinter `welle-harness-regeln`; „sofort“ verlangt, dass der Nutzer den Slice aus der Reihenfolge der Wellen löst (Frage 1 in §6 des Slice, Drift-Log der Roadmap).
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
