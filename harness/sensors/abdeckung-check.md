# `make abdeckung-check` — prüft, dass die Abdeckungstabellen den Abdeckungs-Deklarationen der Tests und der Gegenproben entsprechen

---

## Vertrag

Rot heißt: Eine Abdeckungs-Deklaration eines Tests oder des Kopfes einer Gegenprobe
oder ein Test verletzt
[`SPEC-048`](../vertraege/SPEC-048-abdeckung.md#spec-048--abdeckung-je-anforderung-und-pfad-abdeckung)
(Form der Deklaration, Paarung aus Anforderung und Pfad, Anforderung nicht im
Lastenheft, Test `TestE2E…` ohne Deklaration, geteilte Messung mit falscher Zahlform,
falschem Wertebereich oder im Widerspruch zu einer anderen Deklaration; Punkte 1 bis 3, 7
und 8), oder eine der fünf Tabellen `docs/user/abdeckung-*.md` entspricht nicht dem Stand,
den die Deklarationen ergeben (Punkt 4). Der Lauf schreibt nichts. Wie die Tabellen
entstehen, sagt der Vertrag; diese Datei entscheidet nichts davon.

## Grenze — was das Grün nicht abdeckt

1. **Ob eine Deklaration ihren Pfad belegt** — der Lauf liest Deklarationen, er führt
   keine Tests und keine Gegenproben aus; ob der Test das Akzeptanzkriterium trifft oder
   die Gegenprobe die deklarierte Messmethode, bleibt Urteil des Reviews. Permanent:
   Grenze von `SPEC-048`, akzeptiertes Negativ der Bindung.
2. **Ob die Tests grün laufen** — aus demselben Grund; das sagen `make test` und
   `make test-integration`.
3. **Ob `abdeckung-check` das Skript mit `--check` ruft und an der Kette der Gates hängt (ebenso `abdeckung-gegenprobe`; das Rezept in `harness/mk/abdeckung.mk` zeigt kein Fall), ob eine Gegenprobe an der Kette der Gates hängt, ob die Gesamtzahl n einer geteilten
   Messung stimmt und ob ein Teil zur richtigen Messmethode gehört** — der Lauf prüft nur,
   dass n für eine Anforderung einheitlich ist; das übrige bleibt Review. Permanent: Grenze
   von `SPEC-048`.

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`grep -c '^| \[' docs/user/abdeckung-gesamt.md` zählt die Anforderungen mit
mindestens einer Deklaration.

Der Lauf trägt keine Vollständigkeits-Zeile: Grün ist still.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | keine Fehlform, keine veraltete Tabelle |
| 1 | mindestens eine Fehlform oder mindestens eine veraltete Tabelle |
| 2 | das Lastenheft fehlt, das Lastenheft führt keine Überschrift `### LH-…` oder das Verzeichnis `docs/user/` fehlt (in dieser Reihenfolge); eine Meldung auf stderr, die nennt, was fehlt, nichts geprüft und nichts angelegt |

`make` meldet jeden Ausgang ungleich 0 als Fehler. Fehlform und veraltete Tabelle
teilen den Ausgang 1; auseinander hält sie die Zeile auf stderr: Eine Fehlform nennt
Datei und Zeile (bei einem `TestE2E…` ohne Deklaration die Zeile des Tests, bei einer Gegenprobe die der Deklaration), eine
veraltete Tabelle ihren Pfad unter `docs/user/`, je veralteter Tabelle eine Zeile
(Punkte 5 und 6 des Vertrags). Aus dem Rot führt bei einer Fehlform, die Deklaration oder den
Test zu berichtigen; bei einer veralteten Tabelle `make abdeckung`, das die
abweichenden Tabellen neu schreibt und je geschriebener eine Zeile auf stdout
meldet, und der Commit der geschriebenen Tabellen.

## Sperren

Keine: Der Lauf hat keine benannte Abbruch-Meldung, die vor der Prüfung greift.

## Bindung

[ADR-0028](../../docs/plan/adr/0028-abdeckung-je-anforderung-und-pfad.md) · Vertrag
[`SPEC-048`](../vertraege/SPEC-048-abdeckung.md#spec-048--abdeckung-je-anforderung-und-pfad-abdeckung)
