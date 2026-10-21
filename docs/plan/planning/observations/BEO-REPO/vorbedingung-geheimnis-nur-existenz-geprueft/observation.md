# Eine Vorbedingung gilt als erfüllt, weil das Geheimnis existiert — sein Recht zeigt erst der echte Lauf

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

Die Freigabe-Checkliste (`SPEC-067` Punkt 17) und die Vorbedingungen des Betreibers (`docs/maintainer/releasing.md` §4) verlangen, dass die Geheimnisse gesetzt sind; geprüft wird dabei, dass sie **existieren** (`gh secret list`), nicht, dass ihr **Recht** für den Schritt reicht, der sie benutzt. Der Trockenlauf zeigt es nicht: Er führt den Job, der das Geheimnis braucht, gar nicht aus (`SPEC-067` Punkt 19 (1)), und auch sonst nennt kein Punkt der Checkliste eine Grenze dazu. Plan, Eintrag und Koordination führen die Vorbedingung deshalb als erfüllt, und erst der Lauf des Tags widerlegt es — mit Außenwirkung, nach dem der Tag schon steht.

Abgrenzung zu `BEO-REPO/vertrag-nennt-aufruf-eines-fremdwerkzeugs-ungemessen`: Dort trägt der Vertrag eine ungemessene Annahme über Aufruf und Wirkung eines Fremdwerkzeugs, die am Host in Minuten messbar gewesen wäre; hier ist die Annahme eine Berechtigung in einem fremden Konto, die der Betreiber hält und die niemand im Repo messen kann — die Lücke ist, dass die Checkliste sie weder prüft noch als Grenze nennt. Abgrenzung zu `BEO-REPO/ci-ohne-branch-schutz`: Dort fehlt eine Einstellung auf GitHub und das Repo sagt das auch; hier behauptet die Checkliste eine erfüllte Vorbedingung.
