# Verifikation: slice-harness-vertraege-umzug-ordner

**Rolle:** Verifier. **Datum:** 2026-10-13. **Spanne:** `6459815` bis `1722af9` (HEAD). **Gegenstand:** [slice-harness-vertraege-umzug-ordner](../plan/planning/in-progress/slice-harness-vertraege-umzug-ordner.md) gegen DoD (§2), [ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), [MR-002](../../harness/conventions.md#mr-002) und Spezifikation. Review: [Review-Report](2026-10-13-review-slice-harness-vertraege-umzug-ordner.md).

Alle Belege unten sind in diesem Lauf selbst erzeugt, nicht aus §7 übernommen.

## DoD-Punkte

| DoD-Punkt | Urteil | Beleg (selbst gefahren) |
|---|---|---|
| Liefer-Punkt 1: Ordner, README, drei Dateien, wortgleich | bestätigt | Vergleich je Vertrag gegen `231089e` (Zeilen 2880-2953, 2954-3011, 3165-3215; `###` zu `#`, Schluss-Leerzeilen weg): `diff` leer, Exit 0 für alle drei. Historie-Zeilen 3524/3525 wortgleich im README. Einleitung weicht nur an den Stellen von Randform 4 ab. Reihe als Regel, nicht als Zahl (Randform 6). |
| Liefer-Punkt 2: Spezifikation | bestätigt | `git diff 231089e HEAD --stat -- spec`: 185 Löschungen, keine Einfügung; keine Kennung `SPEC-047/048/050` mehr in `spec/*.md`; `grep -c "harness/" spec/spezifikation.md` = 0; `make docs-check`: 711 Dateien, 0 Befunde. |
| Liefer-Punkt 3: Verweise | bestätigt | `grep -rn "spezifikation.md#spec-0\(47\|48\|50\)"`: keine Treffer. `grep "harness/\|spezifikation.md"` in `SPEC-047`: keine Treffer (MR-002 Adaption 5). Form ohne Link: Treffer nur für `SPEC-049`/`SPEC-052`, Zeitdokumente und Pläne (Randform 9/10). `harness/README.md` Rang 2 (beide Orte, neun Ränge) und Sensors-Links, `AGENTS.md` §2 und §3.12 gelesen und korrekt. |
| `make gates` grün | bestätigt | Eigener Lauf auf HEAD: Exit 0, `lint-gegenprobe: gruen`, `stop-hook-gegenprobe: gruen, 138 Fälle`. `make kopf-check`: Exit 0, keine Ausgabe. |
| Review durchgeführt | bestätigt | Report liegt vor; F-1/F-2 siehe unten. |
| Handbuch und README unberührt | bestätigt | `git diff 6459815..HEAD --stat`: weder `docs/user/` noch Wurzel-`README.md` im Diff. |
| Closure-Notiz mit Lerneintrag | offen (erwartet) | §7 trägt bisher nur die Belege des Implementers; Notiz und Lerneintrag entstehen bei der Closure. |
| Beobachtungs-Register fortgeschrieben | offen (erwartet) | Closure-Schritt; §7 hält fest: keine Evidence-Datei vorgeschlagen, Urteil bei der Closure. |
| Jedes Risiko aus §6 mit Ausgang | offen (erwartet) | Alle vier Risiken stehen noch auf "offen bis Closure". |
| Drei Paarungen | offen (erwartet) | Closure-Schritt. |

## Bewusstes Brechen (selbst nachgefahren, danach `git checkout`, Baum sauber)

| Zusage | Mutation | Ergebnis |
|---|---|---|
| wortgleich | Zeichen `x` vor Zeile 3 von `SPEC-048-abdeckung.md` | Vergleich meldet `3c3` mit beiden Zeilen: rot aus dem richtigen Grund |
| Link auf Vertrag zeigt auf neuen Ort | Link in `harness/sensors/kopf-check.md:9` auf `../../spec/spezifikation.md` zurückgesetzt | `docs-check` rot, 1 Befund. Abweichung zu §7: mein Pfad (`../../`) meldet `repo-escape`, nicht `anchor-missing`; rot ist es in beiden Fällen, die Zusage ist belegt |
| Spezifikation verweist nicht auf `harness/` | Link `../harness/vertraege/README.md` am Ende von `spec/spezifikation.md` | `matrix-forbidden: spec-straten -> aussen nicht erlaubt` |
| `SPEC-047` nennt den Ort nicht | Zeile `Ort: harness/vertraege` angehängt | `grep` meldet Zeile 74 |

## Review-Findings F-1 und F-2

Behoben. `slice-harness-vertraege-umzug-lint` und `slice-harness-vertraege-umzug-gate-kette` tragen die Sendung des Zwischenstands aus `AGENTS.md` §3.12 mit der Kennung dieses Slice in §1, DoD (Liefer-Punkt 3), §3, §6 und §8 (per `grep` gefunden); die Ausnahme für `SPEC-049` fällt im Lint-Slice, die für `SPEC-052` samt Zwischenstand im Gate-Kette-Slice. Der Nehmer-Plan wurde mit Zählung nachgezählt (3 Liefer-Punkte, 2 Schichten).

## Plan-vs-Code-Diff

Alle Zeilen aus §3 sind im Diff, keine Datei darüber hinaus außer dem Review-Report. `.d-check.yml`, `internal/`, `.claude/`, Roadmap unberührt (Roadmap erst bei Closure, wie §3 sagt). Kommentarzeilen in `tools/` ändern nur den Pfad. Größe in der Metrik von `AGENTS.md` §3.14 (Pfade des Plans): 614 bis HEAD (473 auf `0da0595` laut §7), Marke 750 nicht erreicht; Schätzung des Architect 480 bis 560 überschritten, aber unter der Marke (Review-Report, Nachzug F-1/F-2 und §7-Belege kamen hinzu).

## Befunde

- **HIGH/MEDIUM:** keine.
- **LOW (L-1):** §7 nennt für die Mutation "Link zurück" `anchor-missing`; die Mutation, die ein Nachfahrer naheliegend wählt, ergibt `repo-escape`. Beide rot; Zusage unberührt. Kein Nachzug nötig, wenn die Closure die Zeile beim Eintragen genauer fasst.
- **INFO (I-1):** Die Größe in §7 (473) steht vor dem Nachzug von F-1/F-2; die Closure kann 614 nennen.
- **INFO (I-2):** Die vier DoD-Punkte der Closure (Notiz, Register, Risiko-Ausgänge, Paarungen) sind erwartungsgemäß noch offen und gehören vor den `git mv` nach `done/`. Die vier Risiken aus §6 sind nach Befund absehbar *entfallen* (Beispiel-Anker bleibt Text; Zwischenstand abgegeben; `grep` gefahren; Schnitt-Grenze trug) bzw. für `BEO-REPO/schnitt-laesst-haelfte-an-der-grenze` zu urteilen; das ist Urteil der Closure.

## Fazit

Die drei Liefer-Punkte, `make gates`, Review und die Unberührtheit von Handbuch und README sind bestätigt. Offen sind nur die Closure-Pflichten. Freigabe zur Closure.
