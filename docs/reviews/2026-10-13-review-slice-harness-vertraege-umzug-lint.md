# Review slice-harness-vertraege-umzug-lint

**Rolle:** Reviewer. **Datum:** 2026-10-13.
**Gegenstand:** Diff `f4b6d6d..9a5c80d` (Move-Commit bis HEAD) gegen den Plan `docs/plan/planning/done/slice-harness-vertraege-umzug-lint.md`, [ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), [MR-002](../../harness/conventions.md#mr-002), `AGENTS.md` §3.3, §3.7 bis §3.14.
**Summary:** 0 HIGH · 0 MEDIUM · 2 LOW · 3 INFO.

## Findings

### F-1 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Nachzählen beim Eintragen; Eine Bindung ist eine Sendung)
- `pfad`: `docs/plan/planning/open/slice-harness-lint-warten-ohne-frist.md` · "Ort seit `slice-harness-vertraege-umzug-lint`"
- `befund`: Der Diff trägt den Ort des Vertrags mit der Kennung des Gebers in §1, §3 (nun zwei Zeilen: `harness/vertraege/SPEC-049-lint-profil.md` und `spec/spezifikation.md`) und §6 des Nehmers ein, die Bindung ist damit getragen. §8 des Nehmers nennt keine Zählung seiner Liefer-Punkte und Schichten mit dieser Änderung; das Überschreiten der Grenzen ist nicht zu erkennen (der Ort wechselt nur von `spec/` nach `harness/`), die Zählung steht aber nicht da. Als LOW geführt, weil kein Liefer-Punkt dazukommt und die Schichtzahl sich durch den Ortswechsel nicht erhöht; die Einordnung als Sendung im Sinn der Skill-Regel "Nehmer nicht nachgezählt" (MEDIUM) ist ein Urteil des Architect.
- `verifizierbar`: nein
- `klasse`: Ortsangabe im Nehmer ohne Nachzählung in §8

### F-2 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.9 (Der Slice-Plan folgt jeder Korrektur)
- `pfad`: `docs/plan/planning/done/slice-harness-vertraege-umzug-lint.md` · "`SPEC-049` (153 Zeilen, zehn Punkte)"
- `befund`: Kopf-Ziel (§1), DoD-Punkt 1 ("Zeilen 3012 bis 3164") und §3 ("Zeilen 3012 bis 3164") nennen 153 Zeilen bzw. bis 3164; die Datei trägt 152 Zeilen (3012 bis 3163, Zeile 3164 ist leer), und §7 belegt das mit 3012 bis 3163. Ebenso zählt §3 und §6 "17 Dateien" gegen 18 nachgezogene in §7. §7 benennt die Abweichung, §1 bis §3 und §6 ziehen sie nicht nach.
- `verifizierbar`: ja (`git show 231089e:spec/spezifikation.md | sed -n '3164p'` ist leer)
- `klasse`: Zahl im Plan nach Messung nicht nachgezogen

### F-3 INFO

- `quelle`: `AGENTS.md` §3.7 · `pfad`: `docs/plan/planning/next/slice-harness-vertraege-umzug-gate-kette.md` · "(und bis `slice-harness-vertraege-umzug-lint` für `SPEC-049`)"
- `befund`: Der Nehmer nennt den Zwischenstand mit zwei Orten; der Satz ist bis zur Closure dieses Slice wahr und danach Chronik. Er gehört dem Folge-Slice (Klasse 1, Sendung aus `slice-harness-vertraege-umzug-ordner`); Kopf und DoD des Nehmers sind dagegen bereits auf den Stand nach diesem Slice gefasst. Keine Aktion für diesen Diff.
- `klasse`: Zwischenstand im Nehmer, bis Geber schließt

### F-4 INFO

- `quelle`: `docs/plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md` (akzeptiertes Negativ) · `pfad`: `docs/plan/planning/done/slice-harness-lint.md` · "(`spezifikation.md` §11: Vertrag des Gates"
- `befund`: Fließtext ohne Link in Zeitdokumenten nennt weiter `spezifikation.md` §11 für `SPEC-049` (auch `slice-harness-blackbox-driven`, Reviews, Verifikationen). Das deckt die Regel "Links nur dort, wo sie sonst brechen"; kein Verstoß. Die Links in diesen Dateien zeigen alle auf den neuen Ort.
- `klasse`: Fließtext in Zeitdokument bleibt (akzeptiert)

### F-5 INFO

- `quelle`: `AGENTS.md` §3.11 · `pfad`: `docs/plan/planning/done/slice-harness-vertraege-umzug-lint.md` · "Rang- und Ortsaussagen in Prosa haben kein Gate"
- `befund`: Die Zusagen "kein Verweis der Spezifikation auf `harness/`" (nur Links, durch `matrix-forbidden` gefangen) und "kein Pfad in Prosa" tragen nur über die `grep`-Zeilen in §7; der Plan sagt das selbst und sagt keine Sperre zu. Kein Verstoß gegen §3.11.
- `klasse`: Zusage nur über `grep` belegt (benannt)

## Eigene Prüfung (nicht auf Zusage)

- Wortgleichheit `SPEC-049`: `git show 231089e:spec/spezifikation.md | sed -n '3012,3163p' | sed '1s/^### /# /' | diff - harness/vertraege/SPEC-049-lint-profil.md` leer. Zeile 3164 am Stand `231089e` ist leer, Zeile 3165 beginnt `SPEC-050`; die Datei enthält also den ganzen Vertrag.
- Neun Historie-Zeilen: Zeilen 3514 bis 3522 am Stand `231089e`, jede einzeln mit `grep -cxF` im `README.md` des Ordners gefunden (je 1). Reihenfolge nach Randform 2: "Abschnitt angelegt" zuerst, dann fünf Zeilen vom 2026-10-06, drei vom 2026-10-07, eine vom 2026-10-08, danach die Zeile von `SPEC-050`; das entspricht der Reihenfolge der Spezifikation. Die Zeile von `SPEC-049` in der Tabelle der Verträge ist vorhanden, der Link löst auf.
- Spezifikation: `git diff f4b6d6d..HEAD -- spec` enthält nur Löschungen (keine Einfügung); `SPEC-049` und seine Historie-Zeilen sind weg. Die Spezifikation führt keinen Link auf `harness/vertraege/`; die verbliebenen Nennungen von `harness/…` in `SPEC-052` sind Pfade des Gate-Vertrags, kein Verweis auf den Ordner.
- Verweise: Repo-weiter `grep` (ohne `.git`, `.harness`) auf `spezifikation.md#spec-049`, `spezifikation.md SPEC-049` und Fließtext-Formen mit `§11`/`spezifikation` neben `SPEC-049`. Treffer nur in Zeitdokumenten (F-4), im Beispiel-Anker von `SPEC-050` (Text, wortgleich, Entscheidung 3 des Ordner-Slice), in Plänen als Zitat der Greps und in `MR-002`/[ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) (beschreiben den Umzug). Kein Link auf den alten Anker mehr. Lebende Pläne in `open/`/`next/` mit `spezifikation.md §11` nennen den Ort *neuer* Kennungen (Sache von `slice-harness-vertraege-umzug-gate-kette`).
- Kommentare: `.golangci.yml` Zeile 2 bleibt eine Zeile, Pfad `harness/vertraege/SPEC-049-lint-profil.md`; `tools/harness/lint.sh` Zeile 5 und `tools/harness/lint-gegenprobe.sh` Zeile 4 ebenso. `harness/mk/lint.mk` nennt keinen Pfad, unverändert. Der Rest an `SPEC-049` in diesen Dateien ist Kennung ohne Pfad. Keine `//nolint`-Direktive, keine Regel im Profil geändert (§3.2, §3.6).
- `AGENTS.md` §3.12: nur noch `SPEC-052` in `spec/spezifikation.md` §11; `grep -n "SPEC-049" AGENTS.md` zeigt keine Zeile, die §11 als Ort nennt.
- Bindung §3.13: `slice-harness-lint-warten-ohne-frist` nennt in Kopf, §1, §3, §6 und DoD-Umgebung den neuen Ort, §1 mit der Kennung des Gebers; die Historie-Zeile von `SPEC-038` bleibt in `spezifikation.md` §12 (richtig getrennt). Nehmer nicht im Ausschluss, liegt in `open/`. Siehe F-1.
- §3.3: `f4b6d6d` ist der reine Move; die drei folgenden Commits enthalten keinen `git mv`/Rename.
- §3.5: In [ADR-0034](../plan/adr/0034-lint-gate-mit-solid-nahem-profil.md) sind nur vier Link-Pfade geändert, Entscheidung und Gründe unberührt; Nachzug sieht der Plan (§3) und [ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) vor.
- §3.10: Kein neuer Vertrag, kein neuer Code; die Mutationstabelle in §7 deckt die Zusagen mit Gate- oder `grep`-Bezug.
- Mutationen nachgefahren (danach `git checkout`, Baum sauber): (a) `X` vor Zeile 3 der Vertragsdatei: der Vergleich meldet `3c3`. (b) `harness/sensors/lint.md` Zeile 8 zurück auf `../../spec/spezifikation.md#spec-049…`: `make docs-check` meldet `anchor-missing`, 1 Befund. (c) Link auf `../harness/vertraege/README.md` am Ende von `spec/spezifikation.md`: `make docs-check` meldet `matrix-forbidden — Referenz spec-straten → aussen`, 1 Befund. Alle drei rot, wie in §7 behauptet.
- Größe (§3.14, Pfade des Plans): `git diff --numstat f4b6d6d..2b2a1cc` Summe 393 (+198, −195), bis HEAD 438 (+239, −199); Marke 600 nicht erreicht, innerhalb der Schätzung 370 bis 440. Liefer-Punkte 3, Schichten 2 (Spezifikation, Harness); das Zählen der Planung als keine Schicht folgt §3.13.
- §3.7: Neue Kommentare und README-Sätze beschreiben den Ist-Zustand (kein Konjunktiv, keine Chronik); "Ort seit `slice-…`" ist ein Herkunfts-Anker in der Form von §3.7.
- §3.9: Kopf (`Bezug`, `Berührte Spec-Stellen`), §1, §3, §6 und §7 folgen dem Diff bis auf F-2.

## Geprüft, ohne Befund

- `harness/vertraege/` (README und `SPEC-049-lint-profil.md`): ohne Befund.
- `spec/`: ohne Befund.
- `harness/README.md`, `harness/sensors/lint.md`: ohne Befund.
- `.golangci.yml`, `tools/harness/lint.sh`, `tools/harness/lint-gegenprobe.sh`: ohne Befund (nur Kommentarzeilen).
- `AGENTS.md`: ohne Befund.
- `docs/plan/adr/0034-lint-gate-mit-solid-nahem-profil.md`: ohne Befund (nur Link-Pfade).
- `docs/plan/planning/done/`, `docs/reviews/` (nur Link-Pfade): ohne Befund außer F-4.
