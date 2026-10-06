# `make abdeckung-check` — prüft, dass die Abdeckungstabellen den Abdeckungs-Deklarationen der Tests entsprechen

---

## Vertrag

Rot heißt: Eine Abdeckungs-Deklaration oder ein Test verletzt
[`SPEC-048`](../../spec/spezifikation.md#spec-048--abdeckung-je-anforderung-und-pfad-abdeckung)
(Form der Deklaration, Paarung aus Anforderung und Pfad, Anforderung nicht im
Lastenheft, Test `TestE2E…` ohne Deklaration; Punkte 1 bis 3), oder eine der vier
Tabellen `docs/user/abdeckung-*.md` entspricht nicht dem Stand, den die Deklarationen
ergeben (Punkt 4). Der Lauf schreibt nichts. Wie die Tabellen entstehen, sagt der
Vertrag; diese Datei entscheidet nichts davon.

## Grenze — was das Grün nicht abdeckt

1. **Ob eine Deklaration ihren Pfad belegt** — der Lauf liest Deklarationen, er führt
   keine Tests aus; ob der Test das Akzeptanzkriterium trifft, bleibt Urteil des
   Reviews. Permanent: Grenze von `SPEC-048`, akzeptiertes Negativ der Bindung.
2. **Ob die Tests grün laufen** — aus demselben Grund; das sagen `make test` und
   `make test-integration`.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`grep -c '^| \[' docs/user/abdeckung-gesamt.md` zählt die Anforderungen mit
mindestens einer Deklaration.

Der Lauf trägt keine Vollständigkeits-Zeile: Grün ist still.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | keine Fehlform, keine veraltete Tabelle |
| 1 | mindestens eine Fehlform oder mindestens eine veraltete Tabelle; ohne jede Zeile auf stderr, wenn das Lastenheft keine Überschrift `### LH-…` führt (Grenze des Vertrags) |
| 2 | das Lastenheft fehlt; eine Meldung auf stderr, nichts geprüft |

`make` meldet jeden Ausgang ungleich 0 als Fehler. Fehlform und veraltete Tabelle
teilen den Ausgang 1; auseinander hält sie die Zeile auf stderr: Eine Fehlform nennt
Testdatei und Zeile (bei einem `TestE2E…` ohne Deklaration die Zeile des Tests), eine
veraltete Tabelle ihren Pfad unter `docs/user/`, je veralteter Tabelle eine Zeile
(Punkte 5 und 6 des Vertrags). Einen Ausgang 1 ohne Zeile gibt es, wenn das Lastenheft
keine Anforderung führt. Aus dem Rot führt bei einer Fehlform, die Deklaration oder den
Test zu berichtigen; bei einer veralteten Tabelle `make abdeckung`, das die
abweichenden Tabellen neu schreibt und je geschriebener eine Zeile auf stdout
meldet, und der Commit der geschriebenen Tabellen.

## Sperren

Keine: Der Lauf hat keine benannte Abbruch-Meldung, die vor der Prüfung greift.

## Bindung

[ADR-0028](../../docs/plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) · Vertrag
[`SPEC-048`](../../spec/spezifikation.md#spec-048--abdeckung-je-anforderung-und-pfad-abdeckung)
