# ADR-0032: Kopf-Sensor für Slice-Pläne

**Status:** Accepted

**Datum:** 2026-10-05

**Autor:** pt9912

**Bezug:** [`LH-QA-04`](../../../spec/lastenheft.md#lh-qa-04--automatisierbarkeit)

**Schärft:** —

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

`AGENTS.md` §3.9 verlangt, dass der Kopf eines Slice-Plans (`Bezug`, `Berührte Spec-Stellen`) jeder Korrektur folgt. Die Prosa-Regel hat nicht getragen (`BEO-REPO/plan-folgt-korrektur-nicht`). Ein Gate prüft die maschinell entscheidbare Hälfte: Jede Kennung aus Lastenheft, Spezifikation oder Sicht, die §1 oder §2 nennt, steht im Kopf. Für Harness-Werkzeuge gibt es keine Spezifikation; nach `AGENTS.md` §3.12 sind die Randformen dieses Vertrags darum hier entschieden, vor dem ersten Code-Commit. Den Vertrag des Gates nennt `harness/README.md` §Sensors; die Begründung und den Bestand führt der Slice-Plan, der das Gate liefert.

## Entscheidung

Wir wählen einen Abgleich zweier Kennungs-Mengen ohne Urteil über Inhalt: Was §1 oder §2 eines lebenden Slice-Plans an Kennungen nennt, muss die Vereinigung der beiden Kopf-Felder enthalten. Im Einzelnen:

1. **Gegenstand.** Geprüft werden die Slice-Pläne (Dateiname mit Präfix `slice` und Bindestrich, Endung `.md`) in `open/`, `next/` und `in-progress/` unter `docs/plan/planning/`. Nicht geprüft werden `done/` (flach und archiviert, Stubs eingeschlossen), Welle-Pläne, Roadmap und Register.
2. **Kennungen.** Gezählt werden `LH-XX-NN`, `LH-XX-NN.x` (x ein Kleinbuchstabe), `SPEC-NNN` und `ARC-NNN` in genau dieser Schreibweise, als ganzes Wort. Markdown-Auszeichnung ist ohne Belang: Code-Span, Linktext, Hervorhebung und Codeblock zählen gleich. Link-Anker sind klein geschrieben und darum keine Nennung. Eine andere Schreibweise (klein, ohne Punkt, ohne führende Null) ist keine Kennung, weder im Kopf noch in §1 oder §2.
3. **Bereich.** `SPEC-NNN bis SPEC-MMM` (ebenso `ARC`, Backticks erlaubt, NNN < MMM) steht im Kopf wie in §1 und §2 für jede Kennung dazwischen. Andere Bereichsformen kennt der Sensor nicht.
4. **Gleichheit ist exakt.** Eine Unterkennung im Kopf deckt ihre Hauptkennung nicht, eine Hauptkennung keine Unterkennung.
5. **Kopf.** Der Kopf sind die Absätze, die mit `**Bezug:**` und `**Berührte Spec-Stellen:**` beginnen, vor der ersten Überschrift `## `, je bis zur nächsten Leerzeile. Geprüft wird gegen die Vereinigung beider; in welchem Feld eine Kennung steht, prüft der Sensor nicht. `—` ist die leere Menge, Mehrfachnennung ist ohne Belang, eine Kopf-Kennung ohne Nennung in §1 oder §2 ist kein Befund.
6. **§1 und §2.** Ein Abschnitt reicht von der Zeile `## 1.` bzw. `## 2.` bis zur nächsten Zeile `## `; erkannt wird er an der Nummer, nicht am Titel. Ausgenommen ist der Absatz, der mit `Regeln dieser Sektion` beginnt, bis zur nächsten Leerzeile. Alles übrige zählt, auch Abgrenzung, Herkunft und Bereits-Geliefertes; eine Markierung für eine Nennung ohne Anspruch gibt es nicht.
7. **Formfehler.** Fehlt einem geprüften Plan eines der beiden Kopf-Felder oder einer der beiden Abschnitte, ist das ein Befund.
8. **Ausgabe.** Je Befund eine Zeile auf stderr mit Repo-relativem Pfad, Abschnitt und Kennung bzw. dem fehlenden Feld oder Abschnitt; eine Kennung je Abschnitt einmal; sortiert nach Pfad, Abschnitt, Kennung. Der Sensor liest alle Pläne und endet dann mit Exit 1 bei mindestens einem Befund, mit Exit 0 ohne Befund, mit Exit 2, wenn die Ablage fehlt.
9. **Start ohne Stufung.** Das Gate ist vom ersten Lauf an voll scharf; was der Bestand meldet, berichtigt der Slice selbst.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Prosa-Regel §3.9 allein | kein Werkzeug | sechsmal wieder aufgetreten |
| B — Markierung für Nennungen ohne Anspruch (Abgrenzung, Herkunft) | Kopf führt nur, was der Slice beansprucht | neue Syntax; ob eine Nennung Anspruch ist, bleibt Urteil und wandert in die Markierung |
| C — Kennung je Klasse im zugehörigen Feld (`LH-XX-NN` unter `Bezug`) | strengere Form | zwingt Abgrenzungs-Nennungen unter `Bezug`, als diente der Slice ihnen; die Feldwahl ist Urteil |
| D — Hierarchischer Abgleich von Haupt- und Unterkennung | weniger Kopf-Pflege | zwei verschiedene Stellen gälten als eine; zusätzliche Regel |
| E — `done/` mitprüfen | ein Plan mehr im Blick | geschlossene Pläne würden bei jeder Schärfung des Sensors rot und müssten nachträglich editiert werden |
| F — gestufter Start (warnend, dann rot) | kein Bestandsaufwand | eine Lockerung mit Hochschalt-Trigger für einen Bestand, den der Slice selbst berichtigen kann |
| **G — exakter Mengen-Abgleich gegen die Vereinigung, lebende Pläne, voll scharf** | urteilsfrei; keine neue Syntax; geschlossene Pläne bleiben eingefroren | eine Abgrenzung kostet eine Kopfzeile |

## Konsequenzen

- Positiv: Ein Kopf, der §1 oder §2 nicht folgt, fällt im Gate-Lauf auf, nicht im Review.
- Negativ (akzeptiert): Der Sensor prüft weder, ob eine Kennung existiert, noch, ob sie im passenden Feld steht, noch Kennungen außerhalb von §1 und §2. Codeblöcke werden nicht gesondert verfolgt; eine Zeile `## ` in einem Codeblock beendete den Abschnitt. Ein Plan, der nie in einem geprüften Verzeichnis stand, wird nie geprüft; nach der Lifecycle-Regel gibt es ihn nicht.
- Folgepflicht: Eine Gegenprobe hält jede Nummer oben, die eine Mutation fangen kann; `harness/README.md` §Sensors nennt beide Ziele und bindet diese ADR.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Skript | Kennungen aus §1 und §2 lebender Slice-Pläne stehen im Kopf | Gate-Ziel, genannt in `harness/README.md` §Sensors |
| Gegenprobe | Fehlformen je Nummer werden abgelehnt, ein vollständiger Kopf angenommen | Gegenproben-Ziel, genannt in `harness/README.md` §Sensors |

## Re-Evaluierungs-Trigger

Wenn die Vorlage des Slice-Plans Kopf-Felder oder Abschnittsnummern ändert, wenn das Lastenheft eine andere Form der Unterkennung einführt, oder wenn Abgrenzungs-Nennungen den Kopf so füllen, dass er nicht mehr lesbar ist.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-05 | Proposed | — |
| 2026-10-05 | Accepted | — |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**. Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit `Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md` §Hard Rule für Accepted-ADRs).
