# Die Liste grüner Mutanten eines Umbaus ist unvollständig

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Ein Slice baut Code ohne Verhaltensänderung um und fährt je umgebauter Funktion eine
Mutation (`.claude/commands/implement-slice.md` Schritt 19). Die grünen Mutanten, die er
in §7 mit Test-Idee und Grenze einordnet, sind richtig eingeordnet, aber nicht alle:
Review oder Verifikation finden an denselben Funktionen weitere grüne,
verhaltensändernde Mutanten, die in der Sammlung für den Planner fehlen. Die Pflicht
verlangt eine rote Mutation je Funktion, keine Suche nach den grünen.

Abgrenzung zu `BEO-REPO/negativtests-fehlen-bei-neuem-vertrag`: Dort fehlt der Test zu
einer Zusage eines neuen Vertrags. Hier liefert der Slice keinen Vertrag; die Lücke liegt
im Bestand, und was fehlt, ist ihre vollständige Übergabe.
