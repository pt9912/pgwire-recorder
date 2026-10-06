# `make kopf-check` — prüft, dass der Kopf jedes lebenden Slice-Plans die Kennungen aus §1 und §2 führt

---

## Vertrag

Rot heißt: Ein Slice-Plan in `docs/plan/planning/open/`, `docs/plan/planning/next/`
oder `docs/plan/planning/in-progress/` verletzt
[`SPEC-047`](../../spec/spezifikation.md#spec-047--kopf-der-pläne-kopf-check): §1 oder
§2 nennt eine Kennung aus Lastenheft, Spezifikation oder Sicht, die keines der beiden
Kopf-Felder `Bezug` und `Berührte Spec-Stellen` führt, oder dem Plan fehlt ein
Kopf-Feld oder einer der beiden Abschnitte, oder der Plan ist nicht lesbar. Was als
Plan, als Kennung, als Bereich, als Kopf und als Abschnitt gilt, sagen die Punkte 1
bis 7 des Vertrags; diese Datei entscheidet nichts davon.

## Grenze — was das Grün nicht abdeckt

Die Grenze von `SPEC-047`; jede ist ein akzeptiertes Negativ der Bindung und
permanent, solange sie gilt.

1. **Existenz der Kennung** — ein Kopf, der eine Kennung führt, die es nicht gibt, ist
   grün; ob sie existiert, sagt der Vergleich nicht.
2. **Feldwahl** — in welchem der beiden Kopf-Felder eine Kennung steht, ist ohne
   Belang.
3. **Kennungen außerhalb von §1 und §2** — was §3 und folgende, andere Kopf-Felder
   oder der Text vor dem Kopf nennen, prüft der Lauf nicht.
4. **Codeblöcke** — eine Zeile `## ` in einem Codeblock beendet den Abschnitt; was
   danach im Abschnitt steht, ist nicht geprüft.

Was das Grün bei einem nicht lesbaren Lifecycle-Verzeichnis oder bei einem Symlink
als Plan bedeutet, sagt der Vertrag nicht; das ist offen
(`BEO-REPO/gate-uebergeht-ablage-eintrag-still`).

**Wie groß der Ausschnitt ist, sagt das Kommando, nicht diese Datei:**
`find docs/plan/planning/open docs/plan/planning/next docs/plan/planning/in-progress -mindepth 1 -maxdepth 1 -type f -name 'slice-*.md' | wc -l`.

Der Lauf trägt keine Vollständigkeits-Zeile: Grün ist still.

## Ausgabe und Ausgänge

Je Befund eine Zeile auf stderr, `kopf-check: <pfad>: <abschnitt>: <befund>`, sortiert
nach Pfad, Abschnitt und Befund (Punkt 8 des Vertrags). Ausgang der Prüfung:

| Exit | Bedeutung |
|---|---|
| 0 | kein Befund; keine Ausgabe |
| 1 | mindestens ein Befund; alle Pläne sind gelesen, jede Befund-Zeile steht da |
| 2 | die Ablage `docs/plan/planning/` fehlt; nichts geprüft |

`make` meldet jeden Ausgang ungleich 0 als Fehler. Auseinander hält die beiden
Ursachen die Zeile auf stderr: Befund-Zeilen nennen einen Plan und einen Abschnitt,
der Abbruch nennt `docs/plan/planning fehlt`. Aus dem Rot führt bei einem Befund
`<kennung> fehlt im Kopf`, die Kennung in `Bezug` oder `Berührte Spec-Stellen`
nachzutragen oder die Nennung in §1 oder §2 zu berichtigen (`AGENTS.md` §3.9).

## Sperren

- `docs/plan/planning fehlt` — die Ablage fehlt unter der Wurzel, an der der Lauf
  startet → `make kopf-check` in der Wurzel des Repos starten.

## Bindung

[ADR-0032](../../docs/plan/adr/0032-kopf-sensor-fuer-slice-plaene.md) · Vertrag
[`SPEC-047`](../../spec/spezifikation.md#spec-047--kopf-der-pläne-kopf-check)
