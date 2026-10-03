# Welle welle-v1-abschluss: v1-Abschluss: Sessions, Protokollrand, Betrieb, Container

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-<Kennung>-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld. **Geplante Wellen bekommen noch keine
Datei:** Sie stehen in der Roadmap unter *Nächste Wellen* und nirgends sonst —
zwei Positionen, nicht drei.

**Zielmeilenstein:** M3.

**Verantwortlich:** pt9912. **Datum:** 2026-10-03.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Alle v1-Anforderungen sind umgesetzt: parallele Sessions im Record, definierter Protokollrand, kontrolliertes Herunterfahren mit atomarem Schreiben, nicht-interaktive Konfiguration und ein Container-Image mit Betriebsdokumentation. Die Abnahmeszenarien 1 bis 10, 12 und 13 des Lastenhefts sind nachweisbar, Szenario 11 (Homebrew) wird mit dem ersten stabilen Release nachgewiesen, und alle Anforderungen (MUSS und SOLL) sind umgesetzt ([`LH-FA-12`](../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration)).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle [welle-replay-semantik](welle-replay-semantik.md) done.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Die Abnahmeszenarien 1 bis 10, 12 und 13 des Lastenhefts laufen automatisiert; Szenario 11 ist keine Bedingung der Welle und wird mit dem ersten stabilen Release nachgewiesen — das *Mehr* gegenüber den Slice-DoDs.
- Container-Image und Binary sind reproduzierbar gebaut.
- `make gates` grün.
- Closure-Notiz in `welle-v1-abschluss-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-v1-abschluss-sessions | Mehrere Sessions und Verbindungsfehler | [`LH-FA-12`](../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) |
| slice-v1-abschluss-protokollrand | Protokollversion und Protokollrand | [`LH-FA-05`](../../../spec/lastenheft.md#lh-fa-05--simple-query-protocol) |
| slice-v1-abschluss-betrieb | Signale, Schreiben und Konfiguration | [`LH-FA-07`](../../../spec/lastenheft.md#lh-fa-07--persistente-recordings), [`LH-FA-08`](../../../spec/lastenheft.md#lh-fa-08--auswahl-eines-recordings), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus), [`LH-FA-15`](../../../spec/lastenheft.md#lh-fa-15--ci-eignung), [`LH-FA-17`](../../../spec/lastenheft.md#lh-fa-17--maschinenlesbare-konfiguration) |
| slice-v1-abschluss-homebrew | Homebrew-Bereitstellung | [`LH-FA-19`](../../../spec/lastenheft.md#lh-fa-19--bereitstellung-über-homebrew) |
| slice-v1-abschluss-einspielen | Einspielen einer Aufzeichnung, auch zeitgetreu | [`LH-FA-20`](../../../spec/lastenheft.md#lh-fa-20--einspielen-einer-aufzeichnung), [`LH-FA-21`](../../../spec/lastenheft.md#lh-fa-21--zeitgetreues-einspielen) |
| slice-v1-abschluss-container | Container-Image und Betriebsdokumentation | [`LH-FA-16`](../../../spec/lastenheft.md#lh-fa-16--container-eignung), [`LH-QA-03`](../../../spec/lastenheft.md#lh-qa-03--portabilität), [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: —
- Wird blockiert von: Welle [welle-replay-semantik](welle-replay-semantik.md).
- Innerhalb der Welle: `slice-v1-abschluss-container` setzt die übrigen Slices außer `slice-v1-abschluss-homebrew` voraus (Image und Doku bilden den Endstand ab); `slice-v1-abschluss-homebrew` setzt `slice-v1-abschluss-container` voraus (Release-Artefakte).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- TLS-Terminierung, `COPY`, Replikationsprotokoll — nicht Teil des Produkts in dieser Welle.
- Gate für Code-Tabelle und Katalog — Folge-Slice nach `slice-replay-semantik-meldungscodes`.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-v1-abschluss-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
