# Der Nachweis der Installation aus dem Tap belegt nur Linux auf `amd64`

**Sub-Area:** `*` (gesamtes Repo, Kürzel `REPO`)

`make homebrew-nachweis` fährt die Schritte des Nachweises in einem Linux-Container (Image `homebrew/brew`); macOS (`amd64` und `arm64`) und Linux auf `arm64` sind in `LH-FA-19.a` Grenze, das Handbuch nennt nur die belegte Kombination. Die Formel verweist für alle vier Kombinationen auf Binaries, die der Lauf der Pipeline baut, aber nur eine ist je als Installation belegt. Die Angabe des Nutzers, die Installation habe auf macOS funktioniert, ist eine Angabe, kein Beleg (`SPEC-067` Punkt 20 (6)). Kein offener Slice nimmt den Beleg für die übrigen Kombinationen an.
