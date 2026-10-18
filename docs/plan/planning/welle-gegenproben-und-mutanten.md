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

Die Tests und Gegenproben des Bestands halten, was sie zusagen: Die bekannten überlebenden Mutanten, die das Verhalten ändern, haben Tests (`slice-tests-ueberlebende-mutanten`, `slice-tests-ueberlebende-mutanten-driving`), ein Mutations-Gate meldet rot, wenn zu viele Mutanten die Tests überleben (`slice-harness-mutation`), und die Gegenproben der Gates sind geschärft: eine Mutation färbt nur Fälle ihrer Gruppe (`slice-harness-gegenprobe-matrix-p7-flanke`), die Kopie des Arbeitsbaums bricht bei laufenden Agenten nicht (`slice-harness-gegenprobe-baumkopie`), und Warten ohne Frist in Tests meldet der Lint (`slice-harness-lint-warten-ohne-frist`). Der Nutzer hat am 2026-10-18 entschieden, diese bis dahin wellenlosen Slices zu einer Welle zu gliedern. Es gibt kein Abnahmeszenario und keinen Meilenstein, der an ihr hängt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [welle-erster-release](welle-erster-release.md) und [welle-v1-nachtrag](welle-v1-nachtrag.md) sind `done` (Wellen vor Harness, Entscheidung des Nutzers vom 2026-10-08; die Harness-Reihe steht hinter beiden).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Das *Mehr* gegenüber den Slice-DoDs (**Empfehlung des Planners, vom Nutzer nicht entschieden**): Ein Lauf des Mutations-Gates über den Bestand, nachdem `slice-tests-ueberlebende-mutanten` und `slice-tests-ueberlebende-mutanten-driving` in `done/` liegen, nennt keinen der dort geführten Mutanten als überlebend; der Lauf steht in der Closure-Notiz. Die DoDs der Sammel-Slices belegen ihre Tests je für sich, die DoD von `slice-harness-mutation` die Schwelle, nicht ihr Zusammenspiel.
- Für `slice-harness-gegenprobe-matrix-p7-flanke`, `slice-harness-gegenprobe-baumkopie` und `slice-harness-lint-warten-ohne-frist` gibt es **kein gemeinsames *Mehr***: Ihre Bündelung ist Thema (Gegenproben und Güte der Tests), keine Closure-Bedingung. Verwirft der Nutzer das Kriterium oben, gilt das für die ganze Welle.
- `make gates` grün auf dem Endstand.
- Closure-Notiz in `welle-gegenproben-und-mutanten-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-harness-gegenprobe-matrix-p7-flanke | Eine Mutation der Gruppe `p7-sig` färbt nur Fälle ihrer Gruppe, oder die Voraussetzung dafür steht im Vertrag | kein Anspruch des Produkts |
| slice-tests-ueberlebende-mutanten | Tests für überlebende Mutanten, die das Verhalten ändern | [`LH-FA-02`](../../../spec/lastenheft.md#lh-fa-02--record-modus), [`LH-FA-10`](../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-05`](../../../spec/lastenheft.md#lh-qa-05--nachvollziehbare-fehler), [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) |
| slice-tests-ueberlebende-mutanten-driving | Tests für überlebende Mutanten im PGWire-Adapter und im Bootstrap | [`LH-FA-03`](../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-18`](../../../spec/lastenheft.md#lh-fa-18--extended-query-protocol), [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) |
| slice-harness-lint-warten-ohne-frist | Lint meldet Warten ohne Frist in Tests | kein Anspruch des Produkts |
| slice-harness-mutation | Mutationstests als Gate | kein Anspruch des Produkts |
| slice-harness-gegenprobe-baumkopie | Gegenproben kopieren den Arbeitsbaum nach `git ls-files` statt mit `tar` und `cp -R`, durch einen gemeinsamen Helfer | [`LH-QA-07`](../../../spec/lastenheft.md#lh-qa-07--prüfbarkeit-des-quellcodes) (nur Messung 1 und 3 der Gegenproben); kein Anspruch des Produkts |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: —
- Wird blockiert von: Welle [welle-v1-nachtrag](welle-v1-nachtrag.md) (Start-Trigger). Auf Slice-Ebene außerdem von `slice-harness-commit-struktur-id` aus [welle-harness-regeln](welle-harness-regeln.md): Dessen Start-Trigger verlangt `slice-harness-commit-struktur-id` in `done/` oder ausdrücklich zurückgestellt, und die Harness-Reihe führt ihn als ersten der übrigen Harness-Slices. Die Welle verschiebt diese Bindung nicht; sie liegt im Plan des Slice.
- Innerhalb der Welle: `slice-tests-ueberlebende-mutanten-driving` setzt `slice-tests-ueberlebende-mutanten` voraus (Sammelregel), `slice-harness-lint-warten-ohne-frist` den zweiten Sammel-Slice, `slice-harness-mutation` alle drei davor.

**Reihenfolge** (nach den Start-Triggern der Slices; WIP-Limit 1; die Welle ordnet nicht neu):

1. `slice-harness-gegenprobe-matrix-p7-flanke` — letzter der fünf wellenlosen Slices hinter dem Release, hinter `slice-harness-schaetzung-planner-faktor` aus [welle-harness-regeln](welle-harness-regeln.md).
2. `slice-tests-ueberlebende-mutanten`, danach `slice-tests-ueberlebende-mutanten-driving`, danach `slice-harness-lint-warten-ohne-frist`, danach `slice-harness-mutation` (Reihe der Harness-Slices, Entscheidung des Nutzers vom 2026-10-08).
3. `slice-harness-gegenprobe-baumkopie` — **Stelle nicht entschieden**: keine Abhängigkeit von einem anderen Slice; er ändert wie `slice-harness-lint-warten-ohne-frist` die Datei `tools/harness/lint-gegenprobe.sh`, welcher zuerst startet, entscheidet der Nutzer, der zweite rebased (§6 des Plans). Die Stelle hier ist ein Platzhalter, keine Entscheidung.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Neue Produkt-Anforderungen oder Produktcode — die Welle ändert Tests, Gates und Gegenproben, nicht das Produkt.
- Neue Tests für Mutanten, die erst die Messung des Mutations-Gates findet — eigene Arbeit mit eigenem Slice nach `slice-harness-mutation` (§1 dort).
- Mutationen in den Integrationstests — anderer Vorgang (§1 von `slice-harness-mutation`).
- Die Regel-Slices der Harness-Reihe — [welle-harness-regeln](welle-harness-regeln.md).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-gegenproben-und-mutanten-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
