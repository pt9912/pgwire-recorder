# Rollen im selben Klon des Containers wechseln einander den Zweig

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Arbeiten mehrere Rollenläufe (Implementer, Reviewer, Verifier, Architect, Planner) im Dev-Container im selben Klon (`/workspaces/pgr-dev2/pgwire-recorder`), wechselt ein Lauf mit `git switch` oder `git checkout` den Zweig, auf dem ein anderer arbeitet oder liest. Der andere sieht Dateien eines fremden Standes oder committet auf den falschen Zweig; ein Lauf auf einem Stand, den er nicht gewählt hat, belegt etwas anderes, als sein Bericht sagt.

Abgrenzung zu `BEO-REPO/implementer-im-hauptbaum-sperrt-stop-hook`: Dort sperrt ein Lauf im Hauptbaum am Host den Stop-Hook; hier kollidieren Läufe im Klon des Containers, und kein Hook ist beteiligt. Abgrenzung zu `BEO-REPO/span-zuordnung-folgt-dem-hauptbaum-nicht-dem-klon`: Dort ordnet die Erfassung nach dem falschen Arbeitsverzeichnis zu; hier ändert sich der Inhalt des Arbeitsverzeichnisses.
