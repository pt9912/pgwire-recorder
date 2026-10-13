# Review slice-harness-vertraege-umzug-ordner

**Rolle:** Reviewer. **Datum:** 2026-10-13.
**Gegenstand:** Diff `6459815..e77b87e` gegen den Plan `docs/plan/planning/done/slice-harness-vertraege-umzug-ordner.md`, [ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md), MR-002, `AGENTS.md` §3.3, §3.7 bis §3.14.
**Summary:** 0 HIGH · 1 MEDIUM · 1 LOW · 3 INFO.

## Findings

### F-1 MEDIUM

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 (Eine Bindung ist eine Sendung; Adresse nimmt an)
- `pfad`: `AGENTS.md` · "`SPEC-049` und `SPEC-052` in `spec/spezifikation.md` §11 (Ort seit"
- `befund`: Der Satz in §3.12 ist nur bis zum Umzug von `SPEC-049` wahr. `slice-harness-vertraege-umzug-lint` führt `AGENTS.md` weder in §1 noch in der DoD noch in §3 (keine Treffer) und nennt die Kennung des Gebers nicht dafür. `slice-harness-vertraege-umzug-gate-kette` deckt den Nachzug nur über einen allgemeinen Satz (DoD Liefer-Punkt 3, "`AGENTS.md` §3.12 … falls sie … nennen") und einen `grep`, ohne den Geber zu nennen. Szenario: Nach dem Lint-Slice nennt §3.12 `SPEC-049` weiter unter `spec/spezifikation.md` §11, und ein Architect oder Implementer sucht die Randform-Entscheidungen am falschen Ort. Nachzählung der Nehmer (§3.13) ist ebenfalls nicht eingetragen.
- `verifizierbar`: ja (`grep -n "AGENTS.md" docs/plan/planning/next/slice-harness-vertraege-umzug-lint.md` zeigt nur Regel-Verweise)
- `klasse`: Sendung an Folge-Slice ohne Eintrag beim Nehmer

### F-2 LOW

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.9 (Plan folgt der Korrektur)
- `pfad`: `docs/plan/planning/done/slice-harness-vertraege-umzug-ordner.md` · "`AGENTS.md` | update | §2 (Rang 2) und §3.12"
- `befund`: Der Zwischenstand "§3.12 nennt zwei Orte" steht in §3 des Plans, aber weder in §1 noch in §6 (Risiko "Zwischenstand mit zwei Orten" nennt nur die Suche nach `SPEC-049`/`SPEC-052`, nicht die Aussage in `AGENTS.md`). Der Gegenstand der Sendung (F-1) ist im Geber nicht als Sendung benannt.
- `verifizierbar`: nein
- `klasse`: Zwischenstand-Sendung nicht in §1/§6 des Gebers

### F-3 INFO

- `quelle`: `AGENTS.md` §3.5 · `pfad`: `docs/plan/adr/0035-struktur-kennungen-im-commit-hook-abgelehnt.md` · "Schärft:"
- `befund`: Nur Link-Pfade einer `Accepted`-ADR geändert, Entscheidung und Gründe unberührt; der Plan (§3) und [ADR-0040](../plan/adr/0040-vertraege-der-harness-werkzeuge-im-ordner-harness.md) sehen den Nachzug vor. Kein Verstoß, Hinweis für den Verifier.
- `klasse`: Link-Nachzug in Accepted-ADR

### F-4 INFO

- `quelle`: `AGENTS.md` §3.11 · `pfad`: `harness/vertraege/README.md` · "Eine Datei je Kennung"
- `befund`: Die Sätze zu Dateiform, Überschrift und Reihe der Kennungen sind Regeln an den Schreiber ohne Gate; sie sagen keine Prüfung zu. Ein Gate gibt es nur für Links und Anker (`docs-check`).
- `klasse`: Regel ohne Gate (zulässig als Regel an den Schreiber)

### F-5 INFO

- `quelle`: Maintainability · `pfad`: `docs/plan/planning/next/slice-harness-abdeckung-gate.md` · "in `harness/vertraege/` nur den Abdeckungs-Vertrag"
- `befund`: Die Schicht-Abgrenzung nennt nun `harness/vertraege/` neben "die Spezifikation" ausgeschlossen; Schichtzahl ändert sich nicht (Ortsangabe, kein Liefer-Punkt). Keine Aktion.
- `klasse`: Ortsangabe im Nehmer nachgezogen

## Eigene Prüfung (nicht auf Zusage)

- Wortgleichheit: Je Vertrag Text aus `git show 231089e:spec/spezifikation.md` (Überschrift bis zur nächsten `#`-Überschrift), `### ` zu `# `, Leerzeilen am Ende abgeschnitten, `diff` gegen die Datei in `harness/vertraege/`: alle drei leer.
- Historie: Die zwei Zeilen im README sind wortgleich zu den Zeilen 3524 und 3525 am Stand `231089e`.
- `spec/spezifikation.md`: keine Einfügung im Diff, 185 Zeilen entfernt; keine Kennung `SPEC-047`, `SPEC-048`, `SPEC-050` mehr in der Datei, kein Verweis auf `harness/`.
- Verweise: Repo-weiter `grep` (ohne `.git`, `.harness`) auf `spezifikation.md#spec-0(47|48|50)`, `spezifikation.md SPEC-0(47|48|50)`, "Abschnitt für Harness-Werkzeuge", `spezifikation.md … §11/§12` bei den drei Verträgen: außerhalb von Zeitdokumenten (`done/`, Reviews), MR-002 und den Plänen der Folge-Slices keine Treffer. Die übrigen Treffer (`slice-harness-mutation`, `-coverage`, `-meldungskatalog-gate`, `-lh-links-pflicht`, `-lint-warten-ohne-frist`) nennen §11 als Ort neuer Kennungen und sind im Plan §1 ausdrücklich dem Gate-Kette-Slice zugewiesen, der sie in DoD Liefer-Punkt 3 trägt.
- Rang 2: `harness/README.md` und `AGENTS.md` §2 nennen beide Orte, neun Ränge, Charakter-Spalte unverändert; keine Prüfzusage.
- Größe (Metrik §3.14, Pfade des Plans): 473 bis `0da0595`, 520 bis HEAD; Marke 750 nicht erreicht, innerhalb der Schätzung 450 bis 600 (der Beleg in §7 nennt 473 ohne den eigenen Eintrag).
- §3.3: Der Move-Commit `6459815` ist rein; die folgenden Commits enthalten keinen `git mv`.
- §3.10: kein neuer Vertrag, die vier Mutationen der Tabelle in §7 entsprechen den Zusagen mit Gate-Bezug.

## Geprüft, ohne Befund

- `harness/vertraege/` (README und drei Verträge): ohne Befund außer F-4.
- `spec/`: ohne Befund.
- `harness/README.md`, `harness/sensors/`: ohne Befund.
- `tools/` (nur Kommentare): ohne Befund; keine Logik geändert.
- `docs/plan/adr/`: ohne Befund (F-3 INFO).
- `docs/reviews/` (ein Verifikationsbericht, nur Links): ohne Befund.
- `docs/plan/planning/open/` und `next/` (zwei Nehmer-Pläne): ohne Befund (F-5 INFO).
- `internal/`, `.claude/`, `.githooks/`: nicht im Diff, kein Verweis auf die drei Verträge mit altem Pfad.
- HIGH-Liste (Spec-Stratum nennt Artefakt außerhalb, Suppression, Core-Reinheit, Kommentar-Klassen): ohne Befund.
