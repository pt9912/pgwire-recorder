# Review-Report: slice-harness-randformen-vor-code — 2026-10-05

**Review-Art:** Prozess-Texte (Hard Rule, Rollen-Typ, Commands, Register), geprüft gegen Plan, Entscheidungen und Hard Rules (Modul 10). Die DoD prüft dieses Review nicht; das ist Aufgabe des Verifiers.

**Gegenstand:** `git show 8fc60d7` für `slice-harness-randformen-vor-code`: `AGENTS.md` §3.12, `.claude/commands/plan-welle.md` Schritt 6, `.claude/agents/architect.md`, `.claude/commands/implement-slice.md` Schritt 13, Slice-Plan §1/§3/§6 und `BEO-REPO/spec-randform-erst-im-review-entschieden/state.md`. Schwerpunkte laut Auftrag: Wirkung an den drei Stellen (Planung, Architect, Implementer) samt gemeldeter Lücken, Widerspruch zu bestehenden Regeln oder Command-Schritten, Register-Eintrag, Plan folgt Diff, Proportionalität.

**Skill:** `.harness/skills/reviewer.md` @ `8fc60d7`
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-harness-randformen-vor-code.md` (vollständig)
- `AGENTS.md` §1, §3.3, §3.7, §3.9 bis §3.12
- `.claude/agents/architect.md`, `.claude/agents/implementer.md`, `.claude/commands/implement-slice.md` (Schritte 1 bis 25), `.claude/commands/plan-welle.md` (Schritte 4 bis 8)
- Baseline-Regelwerk v6.13.0: `modul-05-planning-harness.md` (§Lifecycle, §Offene Risiken, §Ziel-Form: Slice), `modul-06-roadmap.md` §Das Beobachtungs-Register, `grundlagen-traceability.md` §Herkunfts-Anker; Vorlage `slice.template.md` §4, §6, §7; Modul 8 laut Rollen-Typ
- `docs/plan/planning/observations/README.md`; Register `BEO-REPO/spec-randform-erst-im-review-entschieden` (4 Evidence-Dateien), zum Vergleich `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`, `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`, `BEO-REPO/plan-folgt-korrektur-nicht` samt ihrer Commits `071b94e`, `908908f`
- Vorheriger Report: [Review zu `slice-extended-query-lebendpruefung`](2026-10-05-review-slice-extended-query-lebendpruefung.md); höchste vergebene Nummer im Repo vor diesem Lauf F-363

**Ausgeführte Läufe im Repo:** `make docs-check` nach dem Anlegen dieser Datei. `make gates` lief nicht. Vor dem Anlegen war der Arbeitsbaum unverändert (`git status --short` leer). Keine Mutationen: Der Diff enthält keinen Code und keinen Wächter.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-364 | MEDIUM | Die Rückgabe-Kante des Implementers hat weder Empfänger noch Übergabe-Artefakt. Schritt 13 verweist mit „Rücksprung unten“ auf die Plan-Defekt-Rücksprungkanten. Diese führen von einem roten Sensor oder Gate zurück zu Schritt 13 („den Plan verfeinern“), bleiben also in der Implementer-Rolle. Zusammen mit §3.9, nach dem der Implementer §6 selbst nachzieht, kann er die Randform selbst in §6 entscheiden und weitercoden. Das ist dieselbe stille Entscheidung, nur in Prosa statt im Code. §3.12 sagt „gibt er sie an den Plan zurück“: Der Plan ist ein Dokument, keine Rolle. Offen bleibt auch, ob „wartest auf die Entscheidung“ ein Lifecycle-Rücksprung `in-progress → open` (blockiert, Schritt 11, Vorlage §4) ist oder ein Anhalten im Lauf. Der Kopf desselben Commands verlangt, dass keine Rolle ohne Übergabe-Artefakt rückwärts springt. Einen Widerspruch zu Vorlage §4 gibt es nicht, die Kante ist dort nur nicht verortet. | Modul 8 (kein Pfeil ohne benennbares Artefakt); `implement-slice.md` Kopf; `AGENTS.md` §3.12, §3.9; Reviewer-Skill MEDIUM „unklare Fehlerbehandlung am Rand des Spec-Bereichs“ (sinngemäß: Behandlung der unbekannten Randform) | `.claude/commands/implement-slice.md:10`, `:112`–`:115`, `:123`–`:126`; `AGENTS.md:232`–`:233` | nein (Urteil) | Rückgabe-Kante ohne Empfänger und Artefakt |
| F-365 | LOW | `state.md` nennt vier Zielorte und schreibt dahinter „Herkunfts-Anker `seit slice-harness-randformen-vor-code`“. Den Anker trägt nur `AGENTS.md` §3.12. Schritt 6 von `plan-welle.md`, `architect.md` und Schritt 13 von `implement-slice.md` tragen ihn nicht. `architect.md` steht ohne Suffix; dann muss der Anker nach der Anker-Paarung irgendwo in der Datei stehen. Bei den beiden vorhandenen Einträgen mit Command-Zielort tragen die Schritte den Anker (`implement-slice.md:151`, `:163`, `seit welle-extended-query`). Übernimmt §7 die vier Orte als `liegt in`, schlägt die Paarung bei der Closure für drei davon fehl. Die Form „Zielort und Anker auf einer Zeile“ entspricht sonst den Vorbildern. | `grundlagen-traceability.md` §Herkunfts-Anker (Sensor 1, „bei einer Datei ohne Suffix irgendwo in ihr“); Register-README §Die drei Ausgänge | `docs/plan/planning/observations/BEO-REPO/spec-randform-erst-im-review-entschieden/state.md:3`; `.claude/commands/plan-welle.md:69`–`:71`; `.claude/agents/architect.md:23`–`:26`; `.claude/commands/implement-slice.md:112` | ja (Lesen; Anker-Paarung bei Closure) | Zielort im Register ohne Herkunfts-Anker |
| F-366 | LOW | Der Register-Stand steht vor Review, Verifikation und Closure auf `verkörpert`, und der Commit liegt auf dem Hauptzweig. Das Regelwerk und Schritt 25 von `implement-slice.md` binden das Schreiben ins Register an die Slice-Closure („hängt an der Closure, nicht an der Implementation“). Die Vorbilder setzten `verkörpert` im Closure-Commit (`071b94e`, `908908f`). Der Anker `seit slice-<Kennung>` löst über §7 in `done/` auf, und §7 ist leer. Wer jetzt das Register sichtet, liest die Regel als verkörpert, solange sie noch ungeprüft ist; F-364 bis F-368 können Zielorte noch verschieben. Der Implementer hat sich an den Plan gehalten: Die DoD (Punkt 3) führt den Stand als Liefer-Punkt. Der Konflikt liegt zwischen Plan und Command, nicht im Diff allein. | Modul 6 §Das Beobachtungs-Register („Wer schreibt“); `implement-slice.md` Schritt 25; `grundlagen-traceability.md` §Herkunfts-Anker | `docs/plan/planning/observations/BEO-REPO/spec-randform-erst-im-review-entschieden/state.md:1`; `docs/plan/planning/in-progress/slice-harness-randformen-vor-code.md:85`–`:88`; `.claude/commands/implement-slice.md:197`–`:198` | ja (Lesen) | Register-Stand vor der Closure gesetzt |
| F-367 | LOW | Risiko 3 und §3.12 sagen zu, der Architect prüfe die Liste vor dem Code. Für einen wellenlosen Slice ruft aber kein Command-Schritt den Architect auf. `implement-slice.md` nennt die Rollen-Sequenz nur im Kopf, und `plan-welle.md` Schritt 6 erreicht diesen Slice nicht. Für wellenlose Slices fehlt damit die Planungs-Hälfte. Es bleiben die Rückgabe in Schritt 13 und das Review. Gerade wellenlose Harness-Slices liefern über die Liste aus §3.10 neue Verträge, nämlich Gates (siehe F-369). Die Zusage im Plan reicht weiter als der Ablauf, der sie trägt. | `AGENTS.md` §3.11 (gilt auch für Plan-Zeilen); Plan §6 Risiko 3 | `docs/plan/planning/in-progress/slice-harness-randformen-vor-code.md:152`–`:154`; `AGENTS.md:232`; `.claude/commands/implement-slice.md:5`–`:6` | ja (Lesen) | Zusage im Plan weiter als der Ablauf, der sie trägt |
| F-368 | LOW | Im Satz „Fehlt eine, die der Vertrag offensichtlich hat, ergänzt der Plan sie“ ist „der Plan“ Subjekt, kein Handelnder. Der Architect schreibt Entscheidungen, keine Pläne. Sein Ausgang ist der bestätigte Bezug oder ein Folge-Vorschlag, und keines von beiden nennt die fehlende Randform. Außerdem spricht der Typ von „Frage an den Nutzer“, §3.12 aber von einer „Entscheidung des Nutzers, die §6 festhält“. Die zweite Formulierung driftet also schon im ersten Commit von der Regel ab (AGENTS.md §1: Pointer statt zweiter Formulierung). Dass „offensichtlich“ ein Urteil ist, ist kein Befund: Plan §1 schließt einen Sensor mit genau dieser Begründung aus. | `AGENTS.md` §1, §3.12; Rollen-Typ Eingang/Ausgang; Maintainability | `.claude/agents/architect.md:10`, `:24`–`:26`; `AGENTS.md:230`–`:231` | nein (Urteil) | Handlung ohne Handelnden im Rollen-Typ |
| F-369 | LOW | Das Ziel in §1 grenzt „neuer Vertrag“ auf „(Format, Leser, Protokollrand, Diagnose, Option)“ ein. §3.12 bindet den Begriff über „(§3.10)“ an dessen Liste, und die nennt zusätzlich „Gate“. Risiko 2 erwähnt die Bindung an §3.10, nicht aber diese Ausweitung. §3.9 greift hier nicht im Wortlaut, weil der Slice weder Code noch Gates, Tests oder die Spezifikation ändert. Der Plan nennt aber einen kleineren Geltungsbereich als die Regel. | Maintainability; Plan §1, §6 | `docs/plan/planning/in-progress/slice-harness-randformen-vor-code.md:37`–`:40`, `:149`–`:151`; `AGENTS.md:197`–`:198`, `:227` | ja (Lesen) | Plan-Ziel enger als die gelieferte Regel |
| F-370 | INFO | Hinweis an den Architect. §3.12 legt entschiedene Randformen in §6 ab, dem Abschnitt „Risiken und offene Punkte“. Dort braucht nach Modul 5 jeder Eintrag bei der Closure einen von drei Ausgängen. Eine entschiedene Randform ist kein Risiko und ginge dort als „entfallen“ mit Begründung hinaus. Hält nur §6 eine Entscheidung des Nutzers fest, liegt eine Verhaltensentscheidung zu einem Vertrag in der zeitlichen Schicht. Mit der Archivierung des Slice wandert sie ins Archiv, und der Verifier misst gegen Spec und DoD. Das Beispiel „Richtig“ der Regel legt sie in die Spezifikation. Ob der dritte Ort trägt, entscheidet der Architect. | Modul 5 §Offene Risiken; Source Precedence | `AGENTS.md:230`–`:231`, `:237`–`:238` | nein (Urteil) | — (kein Fehlermuster) |
| F-371 | INFO | Die unveränderte Reviewer-Skill genügt für die Review-Hälfte. Ihr Kontext-Eingang umfasst die Hard Rules von `AGENTS.md`, und ihr HIGH-Anker „ADR-Verstoß (Layer, Tool, Hard Rule)“ fasst einen Verstoß gegen §3.12. Der Ausschluss in Plan §1 trägt. Folge für die Proportionalität: Jede Randform, die ein Diff entscheidet, ohne dass §6 sie nennt, ist nach Skill-Wortlaut HIGH. Nach Häufigkeit und Wirkung unterscheidet der Anker nicht. Das ist ein Hinweis für den Steering Loop, keine Aktion in diesem Slice. | Reviewer-Skill §Klassifikation; `AGENTS.md` §3.12 Begründung | `.harness/skills/reviewer.md`; `AGENTS.md:241`–`:243` | nein | — (kein Fehlermuster) |

## Antwort auf die Schwerpunkte

1. **Wirkung an den drei Stellen.**
   - *Planung:* `plan-welle.md` Schritt 6 trägt die Regel für Slices einer Welle, mit einem Zeiger statt einer zweiten Formulierung. Die gemeldete Lücke (Risiko 3) ist echt und größer, als der Plan sagt. Für wellenlose Slices fehlt nicht nur Schritt 6, sondern auch ein Auslöser für den Architect (F-367).
   - *Architect:* Prüfung und Entscheidung vor dem Code stehen im Kontext-Zuschnitt. Ergänzen kann eine fehlende Randform aber niemand, weil Handelnder und Artefakt fehlen (F-368). „Offensichtlich“ ist als Urteil zulässig.
   - *Implementer:* Schritt 13 erfasst genannte und nicht genannte Randformen, auch beim Implementieren. Die Rückgabe-Kante läuft allerdings in die Implementer-interne Schleife (F-364). Diese Stelle trägt die Regel am stärksten, weil sie auch wellenlose Slices erreicht, und ist zugleich am schwächsten verortet.
   - *Review:* Die unveränderte Skill genügt (F-371).
2. **Widersprüche.** Kein harter Widerspruch. §3.12 und Vorlage §4 passen zusammen: Die Rückführungen in §4 betreffen den Lifecycle, die Randform-Rückgabe den Plan. Ob ein Warten auf den Nutzer „blockiert“ heißt, bleibt aber offen (F-364). Spannung besteht zum Command-Kopf, nach dem kein Rücksprung ohne Übergabe-Artefakt erfolgt (F-364), und zur Closure-Logik von §6 (F-370). Zu §3.9 bis §3.11 gibt es keinen Konflikt. §3.12 verwendet die Vertragsliste aus §3.10 wieder, statt eine eigene zu führen (F-369 betrifft nur den Plan).
3. **Register-Eintrag.** Die Form stimmt: Stand, Zielort und Anker auf einer Zeile, wie bei den Vorbildern. Gegen die Evidenz ist der Wert zulässig (4×, Schwelle 3). Inhaltlich tragen drei der vier Zielorte den Anker nicht (F-365). Das Setzen vor der Closure widerspricht Modul 6 und Schritt 25 von `implement-slice.md`, verlangt es aber die DoD des Plans. Es gehört in die Closure (F-366).
4. **Plan folgt Diff.** §1 (Abgrenzung Register-Stand und übrige Rollen-Typen), §3 (Schritt- und Abschnittsangaben) und §6 (Risiko 2 ergänzt, Risiko 3 neu) folgen dem Diff. Kopf `Bezug: —` und `Berührte Spec-Stellen: —` stimmen. Abweichungen: der Geltungsbereich in §1 (F-369) und die Zusage in Risiko 3 (F-367). Die DoD nennt „Schritt 12/13“, geändert ist nur 13. Das prüft der Verifier.
5. **Proportionalität.** Kein HIGH. Kein Befund verletzt eine Hard Rule, eine ADR oder §3.3. Das einzige MEDIUM (F-364) betrifft die Stelle, an der die Regel ohne Empfänger leerläuft. Die übrigen Befunde sind Form oder Reichweite.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `AGENTS.md` §3.12 — Form | geprüft. Überschrift mit Anker `(seit slice-…)` wie §3.9 bis §3.11. Falsch/Richtig/Begründung im Indikativ, Retirement-Check vorhanden. Befunde F-364, F-369, F-370. |
| `AGENTS.md` §1 — Pointer statt Duplikation | geprüft. `plan-welle.md` zeigt nur; `architect.md` und Schritt 13 formulieren Teile neu, im Rahmen des Vorbilds Schritt 19/20. Befund F-368. |
| Hard Rule 3.7 — Kommentare/Zustandsfelder | geprüft, ohne Befund. Neue Texte beschreiben den Ist-Zustand, keine Chronik. `state.md` nennt Stand und Anker, keine Geschichte. |
| Hard Rule 3.3 — Move und Inhalt getrennt | geprüft, ohne Befund. Kein Move in `8fc60d7`; die Moves `8f0e389`, `837aaca` sind rein. |
| Hard Rule 3.4, 3.5, 3.6, 3.8 | geprüft, ohne Befund. Weder Spec noch ADR noch Gate im Diff. |
| `.claude/agents/` | geprüft. Typname `architect` unverändert, die Rollen-Achse bleibt besetzt. Befund F-368. |
| `.claude/commands/` | geprüft. Befunde F-364, F-365, F-367. |
| `docs/plan/planning/in-progress/` — Plan | geprüft. Befunde F-367, F-369. |
| `docs/plan/planning/observations/` | geprüft. Befunde F-365, F-366. |
| `.harness/skills/reviewer.md`, `.harness/baseline/` | geprüft, ohne Befund. Unverändert, wie Plan §1 begründet. Hinweis F-371. |
| Commit-Message `8fc60d7` | geprüft, ohne Befund. Nennt `slice-harness-randformen-vor-code` und `MR-000`. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 5 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:**

- Rückgabe-Kante ohne Empfänger und Artefakt (F-364)
- Zielort im Register ohne Herkunfts-Anker (F-365)
- Register-Stand vor der Closure gesetzt (F-366)
- Zusage im Plan weiter als der Ablauf, der sie trägt (F-367; Register `BEO-REPO/zusage-im-kommentar-weiter-als-pruefung`)
- Handlung ohne Handelnden im Rollen-Typ (F-368)
- Plan-Ziel enger als die gelieferte Regel (F-369)

**Summary-Zeile:** 0 HIGH · 1 MEDIUM · 5 LOW · 2 INFO — Die Regel wirkt an der Implementer-Stelle am weitesten, hat dort aber keinen Empfänger (F-364). Für wellenlose Slices fehlt der Auslöser des Architect (F-367). Die Zielorte im Register tragen den Anker nicht, und der Stand ist vor der Closure gesetzt (F-365, F-366).
