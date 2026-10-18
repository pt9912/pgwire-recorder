# Ein Implementer-Commit überschreibt Abschnitte des Slice-Plans

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Der Plan eines Slice gehört dem Planner; der Implementer trägt in ihn nur §7 (Belege) ein. Ein Commit des Implementers ersetzte §4 bis §7 des Plans durch seine Belege: Closure-Trigger, Risiken mit ihren Ausgängen, Randformen und die vorab benannten Rückführungs-Bedingungen waren weg, und DoD und §1 verwiesen weiter auf §6. Die Belege, auf die der Plan sich berief, standen nicht in der Datei. Review (F-1, F-2, F-11) und Verifikation (V-1, V-2) fanden es; der Implementer stellte den Plan auf den Stand von `main` her und ließ nur §7 stehen.

Abgrenzung zu `BEO-REPO/plan-folgt-korrektur-nicht`: Dort folgt der Plan dem Code nicht; hier entfernt der Code-Lauf Plantext, den niemand zu ändern bat. Abgrenzung zu `BEO-REPO/implementer-bericht-erreicht-pruefer-nicht`: Dort fehlt der Bericht beim Prüfer; hier liegt er an der falschen Stelle und verdrängt den Plan.
