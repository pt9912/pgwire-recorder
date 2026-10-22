# Welle welle-harness-regeln: Regeln und Werkzeuge des Harness nach v1

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** kein Meilenstein-Bezug und kein Abnahmeszenario; die Welle ist Harness-Arbeit.

**Verantwortlich:** pt9912. **Datum:** 2026-10-18.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Die Regeln und Werkzeuge des Harness, die aus dem Bau von v1 offen sind, sind geschärft: Die Regel-Slices tragen ihre Register-Einträge in `AGENTS.md`, `.claude/agents/` und `.claude/commands/` (`slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`, `slice-harness-schaetzung-planner-faktor`), `slice-mv` findet die Quelle am exakten Namen (`slice-harness-slice-mv-exakter-treffer`), die Gegenprobe der Versionsmatrix färbt mit einer Mutation der Gruppe `p7-sig` nur Fälle ihrer Gruppe (`slice-harness-gegenprobe-matrix-p7-flanke`), und der Commit-Träger lehnt Struktur-Kennungen ab (`slice-harness-commit-struktur-id`). Der Nutzer hat am 2026-10-18 entschieden, diese bis dahin wellenlosen Slices zu einer Welle zu gliedern, und am selben Tag, die Wellen 2 und 3 entlang der Harness-Reihe (Entscheidung des Nutzers vom 2026-10-08) neu zu schneiden, damit sie strikt nacheinander laufen: Diese Welle ist die erste der drei (danach [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md), danach [welle-lh-links](welle-lh-links.md)). Es gibt kein Abnahmeszenario und keinen Meilenstein, der an ihr hängt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [welle-erster-release](done/welle-erster-release/welle-erster-release.md) und [welle-v1-nachtrag](welle-v1-nachtrag.md) sind `done` (Wellen vor Harness, Entscheidung des Nutzers vom 2026-10-08 und vom 2026-10-16: die Harness-Reihe steht hinter beiden).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Das *Mehr* gegenüber den Slice-DoDs (Entscheidung des Nutzers vom 2026-10-18): `make docs-check` und `make gates` sind auf dem Endstand der Welle grün, mit allen Regeln eingetragen. Die DoDs der Slices belegen ihren Lauf je für sich auf dem Stand, den sie kennen, nicht den Stand, den die übrigen fünf danach schreiben.
- Für die Regel-Slices, `slice-harness-slice-mv-exakter-treffer`, `slice-harness-gegenprobe-matrix-p7-flanke` und `slice-harness-commit-struktur-id` gibt es **kein weiteres gemeinsames *Mehr***: Ihre Bündelung ist Thema (Regeln und Werkzeuge des Harness), keine Closure-Bedingung.
- Closure-Notiz in `welle-harness-regeln-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-harness-slice-mv-exakter-treffer | `slice-mv` findet die Quelle am exakten Namen, und eine Gegenprobe hält es | kein Anspruch des Produkts |
| slice-harness-regeln-nach-v1-abschluss | Regeln aus dem Verbindungsende-Block (Abdeckungstabellen zählen nicht als Schicht, Wortlaut-Suche und Arbeitsbaum vor E2E-Mutanten, offene Fragen zu §3.10) | kein Anspruch des Produkts; Register-Einträge in deren Bezug |
| slice-harness-regeln-spezifikation-satzabgleich | Der Architect gleicht vor dem Code die Sätze einer Spezifikationsstelle ab | kein Anspruch des Produkts; `BEO-REPO/spezifikation-widerspricht-sich-im-wortlaut` |
| slice-harness-groessen-metrik | `AGENTS.md` §3.14 trägt den berichtigten Befehl der Größen-Metrik, und [ADR-0056](../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) wird angenommen | kein Anspruch des Produkts; `BEO-REPO/groessenprobe-ohne-messbaren-stand` |
| slice-harness-schaetzung-planner-faktor | Der Planner schätzt in der Metrik von §3.14, mit benanntem Faktor, und nennt die Meldeschwelle | kein Anspruch des Produkts; Register-Einträge in deren Bezug |
| slice-harness-gegenprobe-matrix-p7-flanke | Eine Mutation der Gruppe `p7-sig` färbt nur Fälle ihrer Gruppe, oder die Voraussetzung dafür steht im Vertrag | kein Anspruch des Produkts |
| slice-harness-commit-struktur-id | Commit-Träger lehnt Struktur-Kennungen ab | kein Anspruch des Produkts |
| slice-harness-rollenlauf-zeitgrenze-regel | Die Regel der Stunde steht im Briefing und in den sechs Rollen: Ein Rollenlauf hält bei einer Stunde an und berichtet, die Stunde ist zweite Rückführungs-Marke, `.claude/agents/` steht in der Liste von `AGENTS.md` §3.9 | kein Anspruch des Produkts; `BEO-REPO/rollenlauf-ohne-zeitgrenze` in seinem Bezug |
| slice-harness-rollenlauf-zeitgrenze | Ein Rollenlauf wird nach einer Stunde angehalten: Wächter im Werkzeug-Aufruf, Stufe 1, die Stunde wirkt | kein Anspruch des Produkts; weiterer Beleg in `BEO-REPO/rollenlauf-ohne-zeitgrenze` |
| slice-harness-rollenlauf-zeitgrenze-haertung | Der Wächter der Stunde nennt seine Lage und ist gehärtet (Zeitform, Reihenfolge der Gründe, Selbstauskunft, Umgebung), Stufe 2 | kein Anspruch des Produkts; `BEO-REPO/rollenlauf-ohne-zeitgrenze` in seinem Bezug |
| slice-harness-rollenlauf-wiederholungsgrenze | Nach dem ersten vollen Lauf einer Gegenprobe je Rollenlauf sperrt der Wächter einen weiteren ohne `GP_NUR` oder `GP_MUT` | kein Anspruch des Produkts; Anlass in `BEO-REPO/gate-laufzeit-waechst-mit-gegenprobe` |
| slice-harness-randform-umgebung-der-gegenprobe | Die Umgebung einer Gegenprobe ist eine Randform (`AGENTS.md` §3.12), und der Implementer fährt sie vor der Übergabe einmal unter fremder Locale, Zeitzone und fremdem Benutzer | kein Anspruch des Produkts; `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` in seinem Bezug |
| slice-harness-regeln-implementer-uebergabe | Ein Folgelauf nach Review oder Verifikation, wenn die Nacharbeit neue ausführbare Zeilen trägt, und der Ort des Belegs eines Nachweises in den Anweisungen | kein Anspruch des Produkts; `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` in seinem Bezug |

`slice-harness-rollenlauf-zeitgrenze` ist am 2026-10-21 in `open/` angelegt (Anlass: der Implementer-Lauf von `slice-harness-devcontainer-einrichten` lief 4,34 h, ohne dass eine Marke feuerte) und am **2026-10-22 in vier Slices geschnitten**, nachdem der Architect Vertrag (`SPEC-070`) und Entscheidung ([ADR-0058](../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md)) geschrieben hatte: Seine Schätzung (800 bis 1400 in der Metrik von [ADR-0056](../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md), Gegenprobe als Treiber) lag über der Marke 900. **Entscheidung des Nutzers vom 2026-10-21** ist die **Sache** — die Grenze von einer Stunde je Rollenlauf („Wir haben rollen (siehe .claude/agents) und commands -- max. 1 Stunde"), Anhalten mit Bericht statt Abbruch und der Mechanismus (Wächter im Werkzeug-Aufruf, reines bash+awk, kein Hintergrundprozess, Span-Strom als Quelle). **Entscheidungen des Nutzers vom 2026-10-22** („deine empfehlungen" und, wörtlich, „2. ja"; je vom Koordinator übermittelt): `AGENTS.md` §3.9 führt `.claude/agents/` auf, die Stunde als zweite Rückführungs-Marke steht in der Regel (`slice-harness-rollenlauf-zeitgrenze-regel`), der Wächter deckt nur Rollenläufe, und [ADR-0058](../adr/0058-zeitgrenze-eines-rollenlaufs-als-wachter-im-werkzeug-aufruf.md) ist vom Nutzer angenommen und `Accepted`. **Empfehlung des Planners**, nicht entschieden, sind Welle und Stelle in der Reihenfolge (§5) und der Schnitt: Die Regel steht allein (`-regel`, klein); der Wächter zerfällt nach **Wirkung** (`slice-harness-rollenlauf-zeitgrenze`, Stufe 1) und **Härtung und Selbstauskunft** (`-haertung`, Stufe 2), beide unter 70 % der Marke; die Wiederholungs-Grenze für volle Gegenproben hat ihre Adresse (`-wiederholungsgrenze`). **Voraussetzung der zwei Stufen:** Der Architect ordnet die Punkte von `SPEC-070` den zwei Lieferungen zu (der Vorschlag steht in §6 von `slice-harness-rollenlauf-zeitgrenze`, Randform 7).

`slice-harness-randform-umgebung-der-gegenprobe` ist am 2026-10-22 in `open/` angelegt (Anlass: `BEO-REPO/gegenprobe-haengt-an-umgebung-des-hosts` steht bei 10× ohne Regel). **Entscheidung des Nutzers vom 2026-10-22** ist die **Sache** — „deine empfehlungen“ auf die Vorschläge des Planners im Register, so vom Koordinator übermittelt: Der Slice entsteht, mit der Randform in `AGENTS.md` §3.12 und der Fahrt vor der Übergabe. **Empfehlung des Planners**, nicht entschieden, sind Welle und Stelle in der Reihenfolge (§5), der Schnitt (die Fahrt unter Last und ein Sensor bleiben draußen, dessen §1) und der Wortlaut, den der Implementer dem Nutzer vorlegt.

`slice-harness-regeln-implementer-uebergabe` ist am 2026-10-22 in `open/` angelegt (Anlass: `BEO-REPO/nacharbeit-nach-review-ohne-erneuten-lauf` steht bei 3× ohne Regel, und kein vorhandener Slice nahm sie ohne Verstoß gegen den Schnitt nach Lieferwert an). **Entscheidung des Nutzers vom 2026-10-22** ist die **Sache** („deine empfehlungen", so vom Koordinator übermittelt): Der Slice entsteht, mit der Regel zum Folgelauf nach der Nacharbeit; dazu nimmt er den Ort des Belegs eines Nachweises in den Anweisungen an (Sendung von `slice-harness-devcontainer-einrichten`). **Empfehlung des Planners**, nicht entschieden, sind Welle, Stelle (§5) und Wortlaut. Die Regel aus `BEO-REPO/implementer-ueberschreibt-plan-abschnitte` nimmt er **nicht** an (Kandidat, die Zuweisung durch den Nutzer ist offen; die Größe lässt den Platz).

`slice-harness-groessen-metrik` trägt `Welle: welle-harness-regeln` nach **Entscheidung des Nutzers vom 2026-10-21** (wörtlich, vom Koordinator übermittelt: „ok deine empfehlugnen“ auf eine Vorlage mit sechs Punkten, deren erster Welle und Stelle dieses Slice war); seine Stelle in der Reihenfolge — vor `slice-harness-schaetzung-planner-faktor` — ist mit derselben Entscheidung gesetzt und steht in §5. Bis dahin waren Welle und Stelle Empfehlung des Planners (`BEO-REPO/empfehlung-als-nutzerentscheidung-bezeichnet`).

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md) (deren Start-Trigger verlangt diese Welle in `done/`). Die Bindung auf Slice-Ebene `slice-tests-ueberlebende-mutanten` setzt `slice-harness-commit-struktur-id` voraus ergibt sich damit aus der Reihenfolge der Wellen; sie bleibt im Plan des Slice stehen und ändert sich nicht.
- Wird blockiert von: Welle [welle-v1-nachtrag](welle-v1-nachtrag.md) (Start-Trigger).
- Innerhalb der Welle: keine Abhängigkeit zwischen den Slices außer der Reihenfolge unten.

**Reihenfolge** (nach den Start-Triggern der Slices; WIP-Limit 1; die Welle ordnet nicht neu; Reihenfolge der drei Wellen strikt: diese, [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md), [welle-lh-links](welle-lh-links.md), Entscheidung des Nutzers vom 2026-10-18):

1. Die Reihe der fünf wellenlosen Slices hinter dem Release: `slice-harness-slice-mv-exakter-treffer`, `slice-harness-regeln-nach-v1-abschluss`, `slice-harness-regeln-spezifikation-satzabgleich`, `slice-harness-schaetzung-planner-faktor`, `slice-harness-gegenprobe-matrix-p7-flanke`.
2. `slice-harness-groessen-metrik` steht **vor** `slice-harness-schaetzung-planner-faktor` (Entscheidung des Nutzers vom 2026-10-21, §4): Beide schreiben denselben Absatz `AGENTS.md` §3.14, und die DoD von `slice-harness-schaetzung-planner-faktor` verlangt einen Diff ohne Löschzeile in der Regel selbst, nimmt die Berichtigung also nicht an, solange sie nicht eingetragen ist (`AGENTS.md` §3.13, [ADR-0056](../adr/0056-groessen-metrik-misst-den-gelieferten-diff.md) *Folgepflicht*). Die Reihe aus Punkt 1 bleibt im Übrigen unverändert.
3. `slice-harness-commit-struktur-id` — erster der übrigen Harness-Slices.
4. `slice-harness-rollenlauf-zeitgrenze-regel` steht **hinter** `slice-harness-groessen-metrik` (**Empfehlung des Planners vom 2026-10-22**, nicht entschieden): Der Metrik-Slice schreibt `AGENTS.md` §3.14 neu, und die neue harte Regel entsteht unmittelbar dahinter — ein gleichzeitiger Lauf erzeugt einen Konflikt statt einer Regel. Die Regel an den Leser gilt ab ihrer Lieferung und wartet nicht auf den Wächter. Seine Stelle gegenüber den übrigen Slices der Reihe ist nicht festgelegt.
5. `slice-harness-rollenlauf-zeitgrenze` (der Wächter, Stufe 1) steht **hinter** `slice-harness-rollenlauf-zeitgrenze-regel` (**Empfehlung des Planners vom 2026-10-22**; kein Start-Trigger): Sein Start-Trigger verlangt, dass Vertrag und Entscheidung des Architect vorliegen und dass der Architect die Punkte von `SPEC-070` den zwei Stufen zugeordnet hat (dessen §4).
6. `slice-harness-rollenlauf-zeitgrenze-haertung` (Stufe 2) steht **direkt hinter** `slice-harness-rollenlauf-zeitgrenze` (**Empfehlung des Planners vom 2026-10-22**): Beide ändern Hook und Gegenprobe; Start-Trigger ist die erste Stufe in `done/`.
7. `slice-harness-rollenlauf-wiederholungsgrenze` steht **hinter** `slice-harness-rollenlauf-zeitgrenze-haertung` (**Empfehlung des Planners vom 2026-10-22**): Sie erweitert denselben Hook und erbt Strom, Anker und Fail-Richtung; ihr Start-Trigger verlangt beide Stufen in `done/` und die Randformen des Architect.
8. `slice-harness-randform-umgebung-der-gegenprobe` steht **hinter** `slice-harness-regeln-spezifikation-satzabgleich` (**Empfehlung des Planners vom 2026-10-22**, nicht entschieden): Beide schreiben `AGENTS.md` §3.12, und der spätere rechnet den Wortlaut des ersten ein. Seine Stelle gegenüber den übrigen Slices der Reihe ist nicht festgelegt; ein Start-Trigger außer dem Anspruch durch den Nutzer hat er nicht.
9. `slice-harness-regeln-implementer-uebergabe` steht **hinter** `slice-harness-randform-umgebung-der-gegenprobe` (**Empfehlung des Planners vom 2026-10-22**, nicht entschieden): Beide ändern `.claude/commands/implement-slice.md`; seine Regeln für `implementer.md` und `verifier.md` rechnen den Wortlaut von `slice-harness-rollenlauf-zeitgrenze-regel` ein, der dieselben Dateien ändert.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Das Produkt — die Welle ändert Regeln, Werkzeuge und Dokumente des Harness, keinen Produktcode.
- Tests und Gegenproben des Bestands (außer der Gegenprobe der Versionsmatrix) sowie das Mutations-Gate — [welle-gegenproben-und-mutanten](welle-gegenproben-und-mutanten.md).
- Lastenheft- und MR-Kennungen als Links und deren Pflicht im Doku-Gate — [welle-lh-links](welle-lh-links.md).
- Der Gate-Index um den Teil des Werkzeugs (`slice-harness-gate-index-werkzeug-teil`) — der Slice ist wellenlos: Sein Auslöser ist ein Lauf von ai-harness-init, und eine Welle, die auf Fremdes wartet, schließt nie (Entscheidung des Nutzers vom 2026-10-18).
- Der Wortlaut der Regeln steht beim Nutzer: Die Regel-Slices legen ihn dem Nutzer vor und committen `AGENTS.md` und die Commands erst nach der Bestätigung (Plan der Slices); die Welle ändert das nicht.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-harness-regeln-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
