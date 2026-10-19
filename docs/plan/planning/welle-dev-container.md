# Welle welle-dev-container: Arbeiten im Dev-Container

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

Ein Implementer-Lauf kann im Dev-Container laufen, ohne dass der Arbeitsbaum von `main` am Host für den Stop-Hook gesperrt ist (`BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`): `make container-start` legt das Gerüst an und startet den Container, `make container-einrichten` richtet darin Bare-Repo, Einlieferung von `main`, Arbeitsklon und Zweig ein, und `docs/developer/devcontainer.md` beschreibt den Weg. Der Nutzer hat am 2026-10-18 entschieden, die bis dahin wellenlosen Dev-Container-Slices zu einer Welle zusammenzufassen. Es gibt kein Abnahmeszenario und keinen Meilenstein, der an dieser Welle hängt.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- `slice-harness-devcontainer-u-boot-0-10` liegt in `done/` (erfüllt).
- Die Welle hängt an keiner anderen Welle; [welle-erster-release](done/welle-erster-release.md) und die Release-Slices verlangen nur `slice-harness-devcontainer-u-boot-0-10`, nicht diese Welle (Entscheidung des Nutzers vom 2026-10-17).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`. Für `slice-harness-devcontainer-einrichten-zugang` gilt das auch mit §7 `Gegenstand: entfallen: <Grund>`, wenn der Nutzer den Zugang endgültig nicht verlangt (Liefer-Punkte leer).
- Das *Mehr* gegenüber den Slice-DoDs (**Empfehlung des Planners, vom Nutzer nicht entschieden**): Ein frischer Scratch-Container wird mit `make container-start` und danach `make container-einrichten` aufgebaut, im Klon läuft `make gates` grün, und ein Commit im Zweig des Klons kommt per `git fetch container` am Host an; die Läufe stehen in der Closure-Notiz. Die DoDs belegen Start und Einrichtung je für sich (Attrappen, je einmal ein Scratch-Container), nicht die Kette bis zum Ergebnis am Host. Verwirft der Nutzer das Kriterium, gibt es keines: Die Welle hat dann kein *Mehr* und die Bündelung ist Thema, nicht Closure-Bedingung.
- `make gates` grün.
- Closure-Notiz in `welle-dev-container-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-harness-devcontainer-einrichten-geruest | Make-Ziel `container-start` legt das Gerüst des Dev-Containers an und startet ihn (Teil I) | [LH-QA-03](../../../spec/lastenheft.md#lh-qa-03--portabilität) (nur Host-Voraussetzungen); kein Anspruch des Produkts |
| slice-harness-devcontainer-einrichten-start | `container-start` startet den Container, prüft Port und Bereitschaft und fängt den Abbruch ab (Teil II) | [LH-QA-03](../../../spec/lastenheft.md#lh-qa-03--portabilität) (nur Host-Voraussetzungen); kein Anspruch des Produkts |
| slice-harness-devcontainer-einrichten | Make-Ziel `container-einrichten` richtet Bare-Repo, Einlieferung, Klon und Zweig ein (`docs/developer/devcontainer.md`) | [LH-QA-03](../../../spec/lastenheft.md#lh-qa-03--portabilität) (neu: `ssh` am Host); kein Anspruch des Produkts |
| slice-harness-devcontainer-einrichten-zugang | `container-einrichten` löst den Zugang von Claude Code auf und reicht ihn an den Container (zurzeit nicht gebraucht, Entscheidung des Nutzers vom 2026-10-18) | kein Anspruch des Produkts |

`slice-harness-devcontainer-einrichten-start` trägt `Welle: welle-dev-container` (Entscheidung des Nutzers vom 2026-10-18), ebenso `slice-harness-devcontainer-einrichten-geruest`.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: —
- Wird blockiert von: —

**Reihenfolge** (Reihenfolge der Slices nach deren §4 *Start*; WIP-Limit 1):

1. `slice-harness-devcontainer-einrichten-geruest` — Gerüst und Teil I von `container-start`.
2. `slice-harness-devcontainer-einrichten-start` — Teil II: Start, Port, Bereitschaft, Abbruch; setzt 1 voraus.
3. `slice-harness-devcontainer-einrichten` — setzt den laufenden Container voraus (2).
4. `slice-harness-devcontainer-einrichten-zugang` — setzt 1 und 3 voraus **und** den Wunsch des Nutzers; bleibt bis dahin in `open/`, ohne Priorisierung.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Der Zugang von Claude Code im Container (Quellenauflösung, Login) — nur `slice-harness-devcontainer-einrichten-zugang`, und nur wenn der Nutzer ihn verlangt; die Hostsitzung treibt die Arbeit per `ssh` (Entscheidung des Nutzers vom 2026-10-18).
- Die Entscheidung zum Hostschlüssel nach einem Neubau (`BEO-REPO/hostschluessel-nach-neubau-ohne-entscheidung`) — Entscheidung des Nutzers, nicht dieser Welle.
- Dev-Container als Pflicht für Implementer-Läufe oder als Bedingung eines Gates — die Welle liefert den Weg, ändert keine Regel von `AGENTS.md` und keinen Start-Trigger eines Release-Slice.
- Der erste Release und seine Pipeline — [welle-erster-release](done/welle-erster-release.md).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-dev-container-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
