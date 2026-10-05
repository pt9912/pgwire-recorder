# Welle welle-replay-semantik: Replay-Semantik: Mismatch, Fehlerreplay, Meldungscodes

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

Replay und Record verhalten sich bei Abweichungen und Fehlern wie spezifiziert: Mismatches werden eindeutig diagnostiziert, PostgreSQL-Fehlerantworten werden aufgezeichnet und reproduziert, und Fehler tragen stabile Meldungscodes. Das sind Abnahmeszenario 4 und 6 des Lastenhefts ([`LH-FA-10`](../../../spec/lastenheft.md#lh-fa-10--abweichende-anfrage), [`LH-FA-11`](../../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers)).

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Welle [welle-extended-query](welle-extended-query.md) done.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle Slices der Welle liegen in `done/`.
- Abnahmeszenario 4 (Replay-Abweichung) und 6 (Fehlerreplay) sind für Simple Query end-to-end automatisiert nachgewiesen (Extended: `welle-extended-query`); der Exit-Code beim Herunterfahren folgt in `welle-v1-abschluss` — das *Mehr* gegenüber den Slice-DoDs.
- `make gates` grün.
- Closure-Notiz in `welle-replay-semantik-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-replay-semantik-mismatch | Nicht verbrauchte Interaktionen als Fehler | [`LH-FA-03`](../../../spec/lastenheft.md#lh-fa-03--replay-modus), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) |
| slice-replay-semantik-fehlerreplay | Fehlerantworten, Resultsets und Transaktionen | [`LH-FA-11`](../../../spec/lastenheft.md#lh-fa-11--fehler-des-postgresql-servers), [`LH-FA-06`](../../../spec/lastenheft.md#lh-fa-06--aufzeichnung-von-anfragen-und-antworten), [`LH-FA-12`](../../../spec/lastenheft.md#lh-fa-12--geordnete-interaktionen) |
| slice-replay-semantik-meldungscodes | Meldungscodes, Fehlertext und Logging | [`LH-FA-14`](../../../spec/lastenheft.md#lh-fa-14--diagnoseausgaben), [`LH-FA-13`](../../../spec/lastenheft.md#lh-fa-13--prozessbeendigung-und-fehlerstatus) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: Welle [welle-v1-abschluss](welle-v1-abschluss.md) (Sessions und Betrieb bauen auf Fehlerebenen und Meldungscodes auf).
- Wird blockiert von: Welle [welle-extended-query](welle-extended-query.md).
- Innerhalb der Welle: `slice-replay-semantik-meldungscodes` setzt `slice-replay-semantik-mismatch` voraus.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- Mehrere Sessions, Signalbehandlung, Container — Gegenstand von welle-v1-abschluss.
- Konfigurationsoptionen über `--listen`, `--upstream`, `--input`, `--output` hinaus — Gegenstand von welle-v1-abschluss.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-replay-semantik-results.md` als Geschwister im Ruheort `done/` (entsteht bei Closure).
Zähler: Beobachtungs-Register `../observations/README.md` (eine Ebene über dem Ruheort).
